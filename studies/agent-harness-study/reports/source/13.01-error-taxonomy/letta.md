# Source Analysis: letta

## 13.01 — Error Taxonomy

### Source Info

| Field | Value |
|-------|-------|
| Name | letta |
| Path | `studies/agent-harness-study/sources/letta` |
| Language / Stack | Python (FastAPI, SQLAlchemy, Pydantic, asyncpg) |
| Analyzed | 2026-08-21 |

## Summary

Letta implements a layered, source-oriented error taxonomy. The core is a single module, `letta/errors.py`, defining an `ErrorCode` enum of 11 machine-readable categories (NOT_FOUND, UNAUTHENTICATED, PERMISSION_DENIED, INVALID_ARGUMENT, INTERNAL_SERVER_ERROR, CONTEXT_WINDOW_EXCEEDED, RATE_LIMIT_EXCEEDED, TIMEOUT, CONFLICT, EXPIRED, PAYMENT_REQUIRED) plus ~58 exception classes grouped by origin: provider/model errors (`LLMError` hierarchy), MCP/tool-server errors (`LettaMCPError` hierarchy), database/infrastructure errors (`letta/orm/errors.py`), message-validation errors (`LettaMessageError` family), and request/conflict errors (`ConversationBusyError`, `PendingApprovalError`, etc.).

The taxonomy is actively used for routing at four distinct dispatch points: (1) per-provider `handle_llm_error()` methods translate vendor SDK exceptions into the common taxonomy (`letta/llm_api/llm_client_base.py:369-397`, implemented in `openai_client.py:1216-1427`, `anthropic_client.py:965-1130`, `google_vertex_client.py:895-1030`); (2) the agent loop dispatches on exception type to decide fallback vs retry vs stop (`letta/agents/letta_agent_v3.py:1177-1218`); (3) FastAPI exception handlers map each class to a specific HTTP status code (~51 registrations in `letta/server/rest_api/app.py:544-779`); (4) the streaming service converts typed exceptions into SSE error events with stable `error_type` strings (`letta/services/streaming_service.py:639-776`). Retryability is answerable from the type: rate-limit/server-error/overloaded trigger backoff or model-fallback, connection errors are retried transiently, context-window errors trigger compaction retries, and everything else stops the run with a classified `StopReasonType`.

Weaknesses: classification knowledge is duplicated across mechanisms (exception classes, string heuristics, SSE `error_type` strings, stop reasons), several legacy exceptions carry no `ErrorCode`, and cross-provider detection of context-window/billing failures relies on best-effort message-string matching.

## Rating

**8 / 10**

Rationale: The taxonomy is explicit, centralized, and heavily used for handling decisions — provider mapping is an enforced interface method (`handle_llm_error`, abstract at `letta/llm_api/llm_client_base.py:369`), the agent loop makes retry/fallback/stop decisions purely on error type (`letta/agents/letta_agent_v3.py:1183-1226`), HTTP status codes and SSE error events derive deterministically from classes, and regression tests pin the provider→taxonomy mappings (`tests/adapters/test_letta_llm_stream_adapter_error_handling.py:179-296`). It loses points for: scattered duplication of category strings between layers (e.g., `"llm_timeout"` hard-coded in both `app.py:698` and `streaming_service.py:646`), best-effort string heuristics as a load-bearing classification mechanism (`letta/llm_api/error_utils.py:8-41`), MCP expected-error detection by class-name string comparison (`letta/services/tool_executor/mcp_tool_executor.py:17-19`), and a long tail of exceptions without `ErrorCode` values (`LLMJSONParsingError` at `letta/errors.py:331`, `LocalLLMError` at `letta/errors.py:338`, export/import exceptions not even deriving from `LettaError` at `letta/errors.py:442,474`).

## Evidence Collected

| Area | Evidence | File:Line |
|------|----------|-----------|
| Error code enum | `ErrorCode` with 11 categories incl. CONTEXT_WINDOW_EXCEEDED, RATE_LIMIT_EXCEEDED, TIMEOUT, PAYMENT_REQUIRED | `letta/errors.py:11-24` |
| Base error class | `LettaError(message, code, details)` with BYOK-aware `__str__` | `letta/errors.py:27-45` |
| Provider/model family | `LLMError` tree: Connection, RateLimit, BadRequest, InsufficientCredits, Authentication, PermissionDenied, NotFound, UnprocessableEntity, ServerError, EmptyResponse, Timeout, ProviderOverloaded | `letta/errors.py:259-314` |
| Empty-response subclass rationale | Docstring explains subclassing to preserve retry behavior while enabling special handling | `letta/errors.py:300-306` |
| Tool/MCP server family | `LettaMCPError`, `LettaInvalidMCPSchemaError`, `LettaMCPConnectionError` (INTERNAL_SERVER_ERROR), `LettaMCPTimeoutError` (TIMEOUT) | `letta/errors.py:208-237` |
| Infrastructure/DB family | `NoResultFound`, `MalformedIdError`, unique/FK constraint violations, `DatabaseLockNotAvailableError` (PG 55P03), `DatabaseTimeoutError`, `DatabaseDeadlockError` (PG 40P01) | `letta/orm/errors.py:1-38` |
| Context-window family | `ContextWindowExceededError` → `SystemPromptTokenExceededError`, coded CONTEXT_WINDOW_EXCEEDED | `letta/errors.py:352-371` |
| Concurrency/policy family | `PendingApprovalError`, `ConcurrentUpdateError`, `ConversationBusyError`, `MemoryRepoBusyError`, all CONFLICT-coded with structured `details.error_code` | `letta/errors.py:48-123` |
| Billing family | Org-level `InsufficientCreditsError`; provider-level `LLMInsufficientCreditsError` (PAYMENT_REQUIRED) | `letta/errors.py:478-485`, `letta/errors.py:275-276` |
| Provider mapping contract | Abstract `handle_llm_error(e, llm_config)` documented as "Maps provider-specific errors to common LLMError types" | `letta/llm_api/llm_client_base.py:368-397` |
| OpenAI mapping | Timeout/APIConnection/RateLimit/BadRequest(HTML vs context_length_exceeded)/Auth/Permission/NotFound/Unprocessable/status>=500 split; 402 & 413 special-cased | `letta/llm_api/openai_client.py:1216-1427` |
| Anthropic mapping | Same shape incl. overloaded→`LLMProviderOverloaded` and 413→`ContextWindowExceededError` | `letta/llm_api/anthropic_client.py:965-1130` |
| Google Vertex mapping | google.genai ClientError/ServerError mapped incl. token-limit 400s→`ContextWindowExceededError` | `letta/llm_api/google_vertex_client.py:895-1030` |
| String-heuristic classifiers | `is_context_window_overflow_message` (5 provider phrasings), `is_insufficient_credits_message` (quota/billing phrasings); docstrings admit "best-effort" | `letta/llm_api/error_utils.py:8-41` |
| Agent-loop dispatch on type | `except ValueError` / `LLMEmptyResponseError` → invalid_llm_response; `(LLMRateLimitError, LLMServerError, LLMProviderOverloaded)` → circuit-breaker fallback route else raise; generic `LLMError` → llm_api_error; `ContextWindowExceededError` → compaction retry up to max_summarizer_retries | `letta/agents/letta_agent_v3.py:1177-1226` |
| Circuit-breaker success/failure recording | `record_success` after clean step; `record_failure` before fallback switch | `letta/agents/letta_agent_v3.py:1171-1174`, `letta/agents/letta_agent_v3.py:1190` |
| Streaming dispatch | Typed excepts for LLMTimeoutError/LLMRateLimitError/LLMAuthenticationError/LLMEmptyResponseError/LLMError/SystemPromptTokenExceededError produce SSE `event: error` with fixed `error_type` strings; catch-all → "internal_error" + Sentry | `letta/services/streaming_service.py:639-776` |
| Stream-completeness watchdog | Missing terminal event synthesizes `error_type="stream_incomplete"` and marks run failed | `letta/services/streaming_service.py:601-629`, `letta/services/streaming_service.py:783-799` |
| Stop-reason classification | `StopReasonType.run_status` property partitions reasons into completed / failed / cancelled | `letta/schemas/letta_stop_reason.py:24-49` |
| Credit-gated loop exit | Pending credit task failure sets `stop_reason=insufficient_credits` and breaks the step loop | `letta/agents/letta_agent_v3.py:333-338` |
| HTTP status-code routing | ~51 handlers: 400/404/408/409/410/415/422/500/502/503/504 groups; deadlock & lock handlers add `Retry-After: 1` headers | `letta/server/rest_api/app.py:556-649` |
| LLM-specific HTTP semantics | LLMTimeout→504, LLMRateLimit→429 (BYOK-aware message), LLMInsufficientCredits→402, LLMAuthentication→401, MCPConnection→502, LLMBadRequest→400 | `letta/server/rest_api/app.py:692-779` |
| DB error enrichment | sqlalchemy base wraps timeouts/deadlocks into typed errors with `original_exception` preserved | `letta/orm/sqlalchemy_base.py:84-106`, `letta/orm/sqlalchemy_base.py:916` |
| Transient-retry set | `_RETRYABLE_ERRORS = (httpx.ReadError, WriteError, ConnectError, RemoteProtocolError, LLMConnectionError)` with exponential backoff (2^attempt) | `letta/llm_api/chatgpt_oauth_client.py:100-102`, `letta/llm_api/chatgpt_oauth_client.py:411-435` |
| Legacy backoff limiter | `retry_with_exponential_backoff` raises `RateLimitExceededError(max_retries)` after N 429s; non-rate-limit codes logged as `llm_non_retryable_error` | `letta/llm_api/llm_api_tools.py:38-118` |
| Summarizer fallback | Overload/rate-limit reclassified via `handle_llm_error`, then provider-specific fallback (Anthropic→Bedrock, ZAI→Baseten), original error chained on fallback failure | `letta/services/summarizer/summarizer.py:750-816` |
| Tool-error feedback channel | Tool failures become `ToolExecutionResult(status="error")` fed back to the model rather than raised; friendly-message formatting with char limit | `letta/services/tool_executor/sandbox_tool_executor.py:181-193`, `letta/utils.py:1091-1097`, `letta/schemas/tool_execution_result.py:9` |
| MCP expected-error allowlist | Class-name string check `{"McpError", "ToolError"}` (deliberately import-free) distinguishes user-facing tool errors from infra failures; ExceptionGroup unwrapping | `letta/services/tool_executor/mcp_tool_executor.py:17-19`, `letta/services/tool_executor/mcp_tool_executor.py:64-88` |
| Client-facing error schema | `LettaErrorMessage(run_id, error_type, message, detail, seq_id)` exposed in OpenAPI components | `letta/schemas/letta_message.py:386-403`, `fern/openapi.json` (schemas: LettaErrorMessage) |
| Mapping tests | Regression tests: anthropic streaming APIStatusError→LLMServerError; 413→ContextWindowExceededError; httpx read/write→LLMConnectionError; Google token-limit 400s→ContextWindowExceededError; credit strings/402→LLMInsufficientCreditsError; negative case non-credit→LLMBadRequestError; empty-stream→LLMEmptyResponseError (LET-7679) | `tests/adapters/test_letta_llm_stream_adapter_error_handling.py:29-363` |

## Answers to Dimension Questions

**1. Are errors classified by source?**
Yes, explicitly and at multiple granularities. Top-level sources get dedicated families: provider/model (`LLMError` subtree, `letta/errors.py:259-314`), tool/MCP servers (`LettaMCPError` subtree plus the runtime expected-error allowlist, `letta/errors.py:208-237` and `letta/services/tool_executor/mcp_tool_executor.py:17-19`), validation (`LettaInvalidArgumentError`, `LettaMessageError` family, pydantic `ValidationError` handler at `letta/server/rest_api/app.py:602`), policy/permissions (`BedrockPermissionError` at `letta/errors.py:317`, org credits at `letta/services/credit_verification_service.py:27`), context (`ContextWindowExceededError` family at `letta/errors.py:352-371`), infrastructure (`letta/orm/errors.py:17-38`), timeout (TIMEOUT-coded classes at `letta/errors.py:232-237,309-310`), and concurrency/user-conflict (CONFLICT-coded classes at `letta/errors.py:48-123`). Within the provider family, a second axis records whether the failure is BYOK (`details["is_byok"]`, threaded through every mapper, e.g., `letta/llm_api/openai_client.py:1240-1242`) so messaging and billing responsibility can be attributed.

**2. Is the taxonomy used for handling?**
Yes — it is the primary control input. The agent loop's retry/fallback/stop decision is a pure `isinstance` dispatch (`letta/agents/letta_agent_v3.py:1183-1226`): rate-limit/server/overloaded → record failure with the routing client and switch to the configured fallback model, otherwise fail with `llm_api_error`; empty-response/ValueError → fail fast as `invalid_llm_response`; context-window → in-loop compaction retry (`letta/agents/letta_agent_v3.py:1218-1234`). The ChatGPT-OAuth client retries only its declared `_RETRYABLE_ERRORS` tuple with exponential backoff (`letta/llm_api/chatgpt_oauth_client.py:411-435`). The REST layer turns classes into precise status codes including `Retry-After` headers for lock/deadlock conflicts (`letta/server/rest_api/app.py:613-641`) and 504/429/402/401 for LLM timeout/rate-limit/credits/auth (`letta/server/rest_api/app.py:692-754`). So the answer to "can you tell from the error type whether to retry, escalate, or stop?" is largely yes — though the retryability knowledge lives at the catch sites (tuples and isinstance chains), not as a property on the errors themselves.

**3. Are error categories documented?**
Partially. Documentation is embedded in code: docstrings explain intent (e.g., why `LLMEmptyResponseError` subclasses `LLMServerError` to preserve retry behavior, `letta/errors.py:300-306`; the "best-effort" caveat on heuristic classifiers, `letta/llm_api/error_utils.py:8-14`; the BYOK-mapping contract on the abstract mapper, `letta/llm_api/llm_client_base.py:369-380`). The public API surface documents only `HTTPValidationError` and `LettaErrorMessage` schemas (`fern/openapi.json` components); there is no published catalog of `ErrorCode` values, `error_type` strings, or their HTTP-status semantics. No dedicated error-taxonomy documentation exists in the repo (searched `fern/`, README, and docs paths; only validation-error references found).

**4. Can new error types be added without breaking existing handling?**
Mostly yes, with caveats. New `LettaError` subclasses automatically fall through to sensible defaults: unhandled provider exceptions funnel into `LLMError`/generic 400-classification via the mapper chain ending in the base fallback (`letta/llm_api/openai_client.py:1427`, `letta/llm_api/llm_client_base.py:397`), unknown stream errors become `internal_error` + Sentry capture (`letta/services/streaming_service.py:758-776`), and unknown stop reasons raise loudly rather than silently misclassifying runs (`letta/schemas/letta_stop_reason.py:48-49`). Adding a new `ErrorCode` enum member is additive. Two friction points: FastAPI requires an explicit `add_exception_handler` registration per class to escape the generic handler (51 manual registrations, `letta/server/rest_api/app.py:556-611`), and the `StopReasonType.run_status` property must be exhaustively updated or it raises `ValueError("Unknown StopReasonType")`. Also note the name-string-based MCP allowlist (`mcp_tool_executor.py:76`) would silently misclassify any renamed third-party error class.

## Architectural Decisions

1. **Single central taxonomy module with per-source subtrees.** All agent-facing exceptions live in `letta/errors.py` under one `LettaError` root carrying `(code, details)`; infrastructure errors that must not depend on the core live separately in `letta/orm/errors.py:1-38`. This gives one import point for dispatch code while keeping ORM-layer errors dependency-free.

2. **Provider normalization behind an enforced interface.** Every LLM client must implement `handle_llm_error` translating vendor exceptions into the shared taxonomy (`letta/llm_api/llm_client_base.py:368-397`). This keeps downstream consumers (agent loop, summarizer, streaming) provider-agnostic — they never see an `openai.RateLimitError`.

3. **Errors as data payloads, not just signals.** Structured `details` dicts carry `error_code`, `pending_request_id`, `run_id`, `lock_holder_token`, `max_retries`, `is_byok` etc. (`letta/errors.py:55,68,77,101-104,118-122`), enabling the HTTP/SSE layers to enrich responses without re-parsing messages (`app.py:586-593` surfaces `run_id`; `app.py:707-713` branches messaging on `is_byok`).

4. **Dual-channel failure reporting: raise internally, feed back externally.** Tool failures do not propagate as exceptions; they become `ToolExecutionResult(status="error")` returned to the model as conversation content (`letta/services/tool_executor/mcp_tool_executor.py:83-86`, `sandbox_tool_executor.py:187-193`), while platform failures raise through the taxonomy to terminate runs with typed stop reasons.

5. **Classification layered over unreliable provider signals.** Because providers phrase similar failures differently, Letta adds a string-heuristic layer (`letta/llm_api/error_utils.py:16-23,33-41`) used alongside SDK exception types, acknowledging in docstrings that this is best-effort and centralized precisely so all layers behave consistently.

## Notable Patterns

- **Coded enum + details dict**: every `CONFLICT` subclass embeds a machine-readable `details["error_code"]` like `PENDING_APPROVAL`, `CONVERSATION_BUSY`, `MEMORY_REPO_BUSY` (`letta/errors.py:55,102,119`), giving two orthogonal identifiers per error.
- **Retry-preserving subclassing**: `LLMEmptyResponseError(LLMServerError)` inherits retry behavior but is independently catchable for request-modification strategies (`letta/errors.py:300-306`).
- **Exhaustive status-partition property**: `StopReasonType.run_status` centralizes the reason→RunStatus mapping so no call site re-derives completed/failed/cancelled (`letta/schemas/letta_stop_reason.py:24-49`).
- **Retry-After headers on contention**: database lock/deadlock handlers return 409 with `Retry-After: 1`, converting infra contention into client-actionable backoff (`letta/server/rest_api/app.py:613-641`).
- **Stream watchdog**: streams lacking a terminal event are retroactively failed with synthesized `stream_incomplete` errors, preventing hung clients and stuck runs (`letta/services/streaming_service.py:610-629,783-799`).
- **Regression-test discipline for mappings**: provider→taxonomy conversions are pinned by named regression tests tied to incident IDs (LET-7679) (`tests/adapters/test_letta_llm_stream_adapter_error_handling.py:300-309`).
- **Observability hooks at classification points**: OpenRouter upstream errors logged with searchable `[OPENROUTER_PROVIDER_ERROR]` tags including extracted upstream provider and status (`letta/llm_api/openai_client.py:1222-1230`).

## Tradeoffs

- **String heuristics vs precision**: matching phrases like `"exceeds the context window"` across providers catches cases SDK types miss (streaming `openai.APIError`, HTML-wrapped ALB errors), but breaks if vendors reword messages; the code accepts this explicitly ("best-effort", `letta/llm_api/error_utils.py:9-14`).
- **Name-based MCP allowlist vs import isolation**: checking `exception_to_check.__class__.__name__ in {"McpError", "ToolError"}` avoids importing fastmcp/mcp packages (`mcp_tool_executor.py:17-19`) at the cost of fragility against renames or new expected-error classes.
- **Centralization vs coupling**: one `errors.py` means every subsystem imports core; the ORM family deliberately stays separate (`letta/orm/errors.py`), but export/import exceptions bypass `LettaError` entirely (`letta/errors.py:442,474`), creating a second-class tier without codes.
- **Catch-site-encoded retry policy vs self-describing errors**: retryability is expressed where errors are caught (tuples/isinstance chains), which keeps error classes simple but means a new error type can default to "stop" unless someone edits dispatch sites.
- **Duplication of category strings**: `error_type` literals ("llm_timeout", "llm_rate_limit", "llm_authentication", ...) appear in both the REST handlers (`app.py:692-779`) and streaming service (`streaming_service.py:639-706`) as independent literals rather than shared constants — drift-prone but currently consistent.

## Failure Modes / Edge Cases

- **Heuristic misses**: a provider introducing new wording for context-window exhaustion would classify as `LLMBadRequestError` (400/INVALID_ARGUMENT) instead of triggering compaction retry — the exact scenario the comment at `letta/llm_api/openai_client.py:1314-1321` describes working around for streamed `APIError`s.
- **Silent misclassification risk in MCP executor**: unexpected errors re-raise (good), but if fastmcp renames `ToolError`, user-facing tool errors would crash steps instead of returning friendly messages to the model (`mcp_tool_executor.py:76-88`).
- **Uncoded legacy errors**: `LocalLLMError`, `LocalLLMConnectionError`, `LLMJSONParsingError`, `BedrockError` construct `LettaError` without a code, so `__str__` degrades to bare message and clients cannot branch on category (`letta/errors.py:331-349`).
- **Broad 400 net**: registering `ValueError` to the 400 handler (`app.py:564`) means internal programming errors can surface as user-caused bad requests.
- **Unknown stop reason raises at terminal-report time**: an added-but-unmapped `StopReasonType` fails inside `run_status` (`letta_stop_reason.py:48-49`), i.e., late, during run finalization.
- **ExceptionGroup edge**: only single-member ExceptionGroups are unwrapped before the MCP allowlist check (`mcp_tool_executor.py:70-74`); multi-error groups skip classification and re-raise.
- **Fallback loops bounded implicitly**: summarizer guards recursion with `is_fallback` flag and chains original error on fallback failure (`summarizer.py:770-775,814-816`); agent-loop fallback switches config persistently for subsequent retries (`letta_agent_v3.py:1201-1211`), relying on the outer max-attempts bound.

## Future Considerations

- Encode retryability/fallback eligibility as metadata on the error classes (e.g., a `retryable` attribute or registry) so dispatch sites consume a declaration instead of maintaining parallel isinstance tuples (`letta/agents/letta_agent_v3.py:1183`, `chatgpt_oauth_client.py:102`).
- Hoist the duplicated SSE/REST `error_type` string literals into shared constants co-located with `ErrorCode` (`app.py:692-779`, `streaming_service.py:639-706`).
- Backfill `ErrorCode` onto uncoded families (`LLMJSONParsingError`, `LocalLLM*`, Bedrock) and migrate export/import exceptions onto the `LettaError` root (`letta/errors.py:442-471`).
- Publish the taxonomy (ErrorCode ↔ HTTP status ↔ SSE error_type ↔ stop_reason matrix) in `fern/openapi.json` / docs; today only `HTTPValidationError` and `LettaErrorMessage` are schema-documented.
- Replace the MCP class-name allowlist with a capability probe or versioned protocol marker when import isolation can be relaxed (`mcp_tool_executor.py:17-19`).
- Add telemetry dimensions for taxonomy hits/misses of the heuristic classifiers to detect provider wording drift early (`error_utils.py` has no metric emission).

## Questions / Gaps

- No evidence found of documentation enumerating `ErrorCode` values for external API consumers beyond inline docstrings; searched `fern/openapi.json` components, repo-root markdown docs, and `letta/client/`.
- No test coverage found for the FastAPI exception-handler table itself (status-code assertions per class); the handler registrations (`app.py:556-779`) appear verified only indirectly. Searched `sources/letta/tests` for `ErrorCode`, `add_exception_handler`, and status-code-per-error patterns; only the adapter-level mapping tests exist (`tests/adapters/test_letta_llm_stream_adapter_error_handling.py`).
- Whether the legacy synchronous retry wrapper (`retry_with_exponential_backoff`, `letta/llm_api/llm_api_tools.py:38-118`, whose inner `wrapper` begins with a bare `pass` statement at line 50-51 yet still executes the retry loop) is still reachable in production paths could not be fully confirmed; it decorates `create` at line 122, suggesting partial dead code worth auditing.

---

Generated by dimension 13.01 (Error Taxonomy) against `letta`.
