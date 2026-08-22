# Source Analysis: crewai

## Dimension 13.01: Error Taxonomy

### Source Info

| Field | Value |
|-------|-------|
| Name | crewai |
| Path | `studies/agent-harness-study/sources/crewai` |
| Language / Stack | Python monorepo (`lib/crewai` core framework, `lib/crewai-core`, `lib/crewai-files`, `lib/crewai-tools`, `lib/cli`, `lib/devtools`) |
| Analyzed | 2026-08-21 |

All citations below are workspace-relative paths rooted at `studies/agent-harness-study/sources/crewai`.

## Summary

CrewAI does not maintain a single, central error taxonomy. Instead, error classification is **federated per subsystem**, and the maturity varies sharply between them:

1. **A2A protocol layer (most mature):** a full JSON-RPC-style taxonomy with an `IntEnum` of error codes, a dataclass exception hierarchy, default message registry, and explicit retry/client-error predicates (`lib/crewai/src/crewai/a2a/errors.py:25-503`).
2. **File upload subsystem (mature):** a transient/permanent classification hierarchy whose types directly drive an exponential-backoff retry loop vs. stop decision (`lib/crewai-files/src/crewai_files/processing/exceptions.py:86-103`, `lib/crewai-files/src/crewai_files/resolution/resolver.py:337-398`).
3. **Core agent loop (ad-hoc at the edges):** dispatch is by exception type for a few known classes (`OutputParserError`, context-length errors), but provider errors are classified by **module-name prefix string matching** (`e.__class__.__module__.startswith("litellm")`) and context-length detection by **error-message substring matching** against a hardcoded phrase list (`lib/crewai/src/crewai/utilities/exceptions/context_window_exceeding_exception.py:4-13`, `lib/crewai/src/crewai/agents/crew_agent_executor.py:444-458`).
4. **Event bus (observability taxonomy):** errors are re-classified *by source* as typed events — LLM call failed, agent execution error, tool usage/validation/selection/execution errors, task failed, crew kickoff failed (`lib/crewai/src/crewai/events/types/llm_events.py:117`, `tool_usage_events.py:72-97`, `crew_events.py:51`).

Tool and validation failures are generally not raised to callers; they are converted into natural-language error strings (i18n templates) fed back into the model's context so the agent can self-correct, with bounded retries.

**Rating: 6 / 10.**

## Rating

Rating: 6/10

| Score | Meaning |
|-------|---------|
| **6** | Present but inconsistent: strong typed taxonomies in A2A and file handling with retry-routing logic and tests for some paths, but the core agent loop relies on fragile string/module-name matching, there is no central exceptions module (the aggregator is empty), and no dedicated documentation page for error categories. |

Rationale: The A2A and crewai-files subsystems alone would score 7–8 (clear model, explicit interfaces, operational safeguards). The core harness drags this down: classification of provider/context errors in the main execution path is message-substring and module-prefix based, which is exactly the "fragile" failure mode the rubric flags.

## Evidence Collected

Every entry includes a file path with line numbers, relative to `studies/agent-harness-study/sources/crewai`.

| Area | Evidence | File:Line |
|------|----------|-----------|
| Context/model error type | `LLMContextLengthExceededError` + `CONTEXT_LIMIT_ERRORS` phrase list used to classify provider messages | `studies/agent-harness-study/sources/crewai/lib/crewai/src/crewai/utilities/exceptions/context_window_exceeding_exception.py:4-16` |
| String-based classifier | `_is_context_limit_error()` matches lowercase substrings like `"maximum context length"`, `"context_length_exceeded"` | `studies/agent-harness-study/sources/crewai/lib/crewai/src/crewai/utilities/exceptions/context_window_exceeding_exception.py:32-44` |
| A2A error code enum | `A2AErrorCode(IntEnum)` with JSON-RPC standard (-32700..-32603), A2A-specific (-32001..-32007), CrewAI extensions (-32009..-32018) | `studies/agent-harness-study/sources/crewai/lib/crewai/src/crewai/a2a/errors.py:25-98` |
| A2A exception hierarchy | Dataclass hierarchy rooted at `A2AError` with per-code subclasses carrying structured fields (`task_id`, `retry_after`, `required_scope`) | `studies/agent-harness-study/sources/crewai/lib/crewai/src/crewai/a2a/errors.py:127-162, 352-365, 369-384` |
| A2A retry predicate | `is_retryable_error()` → INTERNAL_ERROR, RATE_LIMIT_EXCEEDED, TASK_TIMEOUT are retryable | `studies/agent-harness-study/sources/crewai/lib/crewai/src/crewai/a2a/errors.py:464-478` |
| A2A client-vs-server split | `is_client_error()` enumerates request-caused codes (parse, invalid params, not-found, unsupported version) | `studies/agent-harness-study/sources/crewai/lib/crewai/src/crewai/a2a/errors.py:481-503` |
| File error taxonomy | `FileProcessingError` tree: validation (`FileTooLargeError`, `UnsupportedFileTypeError`), dependency, transient/permanent, upload variants via multiple inheritance | `studies/agent-harness-study/sources/crewai/lib/crewai-files/src/crewai_files/processing/exceptions.py:4-103` |
| Upload error classifier | `classify_upload_error()` maps exception type name (`RateLimit`, `APIConnection`, `Authentication`, `BadRequest`) and `status_code` (>=500 or 429) to Transient vs Permanent | `studies/agent-harness-study/sources/crewai/lib/crewai-files/src/crewai_files/processing/exceptions.py:106-135` |
| Taxonomy-driven retry routing | `_upload_with_retry`: `PermanentUploadError` → stop immediately; `TransientUploadError`/unknown → exponential backoff up to `UPLOAD_MAX_RETRIES`; error kind recorded in metrics metadata | `studies/agent-harness-study/sources/crewai/lib/crewai-files/src/crewai_files/resolution/resolver.py:361-398` |
| Provider-specific classifiers | `_classify_gemini_error` (rate limit → transient) and `_classify_s3_error` (throttling → transient) | `studies/agent-harness-study/sources/crewai/lib/crewai-files/src/crewai_files/uploaders/gemini.py:48-64`, `bedrock.py:28-49` |
| Tool error types | `ToolUsageError` and `ToolUsageLimitExceededError` | `studies/agent-harness-study/sources/crewai/lib/crewai/src/crewai/tools/tool_usage.py:68-73`, `structured_tool.py:111-112` |
| Tool error → model feedback loop | Exceptions caught, formatted via i18n template, retried until `_max_parsing_attempts` (3, or 2 for large OpenAI models), then returned as `ToolUsageError` string to the LLM ("Moving on then") | `studies/agent-harness-study/sources/crewai/lib/crewai/src/crewai/tools/tool_usage.py:100-118, 422-465` |
| i18n error-message taxonomy | `"errors"` block keyed by category: `wrong_tool_name`, `tool_arguments_error`, `tool_usage_exception`, `task_repeated_usage`, `validation_error`, `force_final_answer` | `studies/agent-harness-study/sources/crewai/lib/crewai/src/crewai/translations/en.json:45-56` |
| Agent-loop error dispatch | `except OutputParserError` → re-prompt handler; generic `Exception` → litellm-module check → context-length check → unknown-error log + raise | `studies/agent-harness-study/sources/crewai/lib/crewai/src/crewai/agents/crew_agent_executor.py:434-458` |
| Module-prefix provider check | `e.__class__.__module__.startswith("litellm")` treats any litellm-raised error as passthrough (no agent-level retry) | `studies/agent-harness-study/sources/crewai/lib/crewai/src/crewai/agents/crew_agent_executor.py:445-446`, also `agent/core.py:695-704` |
| Context-window escalation policy | `handle_context_length()`: summarize messages if `respect_context_window`, else `SystemExit` with user guidance | `studies/agent-harness-study/sources/crewai/lib/crewai/src/crewai/utilities/agent_utils.py:712-749` |
| Agent retry policy | `_check_execution_error()`: litellm errors re-raised immediately; otherwise retried until `max_retry_limit` exceeded, then emit `AgentExecutionErrorEvent` and raise | `studies/agent-harness-study/sources/crewai/lib/crewai/src/crewai/agent/core.py:685-717` |
| Provider exception mapping (OpenAI) | `NotFoundError` → `ValueError`, `APIConnectionError` → `ConnectionError`, generic → context-length check or re-raise; each branch emits `LLMCallFailedEvent` first | `studies/agent-harness-study/sources/crewai/lib/crewai/src/crewai/llms/providers/openai/completion.py:923-947` |
| Provider transport retries delegated | `max_retries: int = 2` forwarded to the OpenAI SDK client config rather than handled in CrewAI code | `studies/agent-harness-study/sources/crewai/lib/crewai/src/crewai/llms/providers/openai/completion.py:209, 342-343` |
| BaseLLM interface contract | Docstrings declare implementations should raise `ValueError` (invalid format), `TimeoutError` (request timeout), `RuntimeError` (other) | `studies/agent-harness-study/sources/crewai/lib/crewai/src/crewai/llms/base_llm.py:316-318, 353-355` |
| Validation result type | `GuardrailResult` model with `success`/`result`/`error` and mutual-exclusivity validator | `studies/agent-harness-study/sources/crewai/lib/crewai/src/crewai/utilities/guardrail.py:60-103` |
| Guardrail retry loop | Validation failure becomes retry context via `validation_error` i18n template; after `guardrail_max_retries`, raises plain `Exception` | `studies/agent-harness-study/sources/crewai/lib/crewai/src/crewai/task.py:1292-1324` |
| Rate limiting | `RPMController.check_or_wait()` blocks (sleeps) before requests instead of classifying a rate-limit error | `studies/agent-harness-study/sources/crewai/lib/crewai/src/crewai/utilities/rpm_controller.py:38-64` |
| Event-bus source taxonomy | Typed error events per source: `LLMCallFailedEvent(error: str)`, `AgentExecutionErrorEvent`, `TaskFailedEvent`, `CrewKickoffFailedEvent` | `studies/agent-harness-study/sources/crewai/lib/crewai/src/crewai/events/types/llm_events.py:117-120`, `agent_events.py:52-67`, `task_events.py:48-51`, `crew_events.py:51-54` |
| Fine-grained tool error events | `ToolUsageErrorEvent`, `ToolValidateInputErrorEvent`, `ToolSelectionErrorEvent`, `ToolExecutionErrorEvent` distinguish tool-failure subcategories | `studies/agent-harness-study/sources/crewai/lib/crewai/src/crewai/events/types/tool_usage_events.py:72-97` |
| Storage/auth/tool-lib errors | `DatabaseOperationError`, `AgentRepositoryError`, `AuthError` (core), `BedrockError` tree with `BedrockValidationError` (tools lib) | `studies/agent-harness-study/sources/crewai/lib/crewai/src/crewai/utilities/errors.py:10-25, 55-56`, `crewai-core/src/crewai_core/auth/token.py:8`, `crewai-tools/src/crewai_tools/aws/bedrock/exceptions.py:4-17` |
| Empty central aggregator | `utilities/exceptions/__init__.py` contains only a docstring — no central export of the taxonomy | `studies/agent-harness-study/sources/crewai/lib/crewai/src/crewai/utilities/exceptions/__init__.py:1` |
| Tests: context-window conversion | `test_context_window_exceeded_error_handling` asserts litellm `ContextWindowExceededError` converts to `LLMContextLengthExceededError` (streaming and non-streaming) | `studies/agent-harness-study/sources/crewai/lib/crewai/tests/test_llm.py:367-401` |
| Docs: retry & context knobs | `max_retry_limit` documented (default 2); `respect_context_window` and automatic context-window management documented | `studies/agent-harness-study/sources/crewai/docs/edge/en/concepts/agents.mdx:60-61, 374` |
| Docs: protocol error codes | A2A module docstring documents JSON-RPC code ranges and extension ranges | `studies/agent-harness-study/sources/crewai/lib/crewai/src/crewai/a2a/errors.py:1-10` |

## Answers to Dimension Questions

**1. Are errors classified by source?**
Partially. There is no single source-oriented taxonomy; instead each subsystem defines its own:
- *Model/provider*: `LLMContextLengthExceededError` plus provider-specific mapping in completion handlers (`openai/completion.py:923-947`); litellm-originated errors identified only by module-name prefix (`agent/core.py:695`).
- *Tool*: `ToolUsageError` (`tools/tool_usage.py:68`) and four distinct tool error event types (`events/types/tool_usage_events.py:72-97`).
- *Validation*: `GuardrailResult.error` (`utilities/guardrail.py:77-79`), `FileValidationError` subtree (`crewai-files/.../exceptions.py:18-63`), `BedrockValidationError` (`crewai-tools/.../bedrock/exceptions.py:16`).
- *Policy/auth*: `AuthError` (`crewai-core/src/crewai_core/auth/token.py:8`), `AuthenticationRequiredError` / `AuthorizationFailedError` / `RateLimitExceededError` (`a2a/errors.py:330-365`).
- *Context*: `LLMContextLengthExceededError`.
- *Infrastructure*: `DatabaseOperationError` (`utilities/errors.py:10`), `TransientFileError` (`crewai-files/.../exceptions.py:86`).
- *Timeout*: `TaskTimeoutError` and `A2APollingTimeoutError` (`a2a/errors.py:369-384, 21-22`); `TimeoutError` convention at the `BaseLLM` interface.
- *User*: No dedicated user-error category was found; user-caused problems surface as validation errors or generic `ValueError`s.
The event bus is where source classification is most systematic — every error event is namespaced by emitting subsystem.

**2. Is the taxonomy used for handling?**
Yes in two subsystems, partially elsewhere. A2A exposes `is_retryable_error()` / `is_client_error()` for programmatic retry-vs-stop decisions (`a2a/errors.py:464-503`). File uploads route on the transient/permanent distinction inside the retry loop itself — permanent errors stop, transient errors back off exponentially (`resolver.py:375-393`). In the core agent loop, handling is driven by exception type for `OutputParserError` and context-length errors, but by fragile heuristics for everything else: litellm module-prefix matching decides passthrough vs retry (`crew_agent_executor.py:444-446`), and substring matching against eight hardcoded phrases decides whether a provider error is a context overflow (`context_window_exceeding_exception.py:4-13`). Tool errors bypass Python exception flow entirely: they become English sentences fed back to the model, bounded by attempt counters (`tools/tool_usage.py:422-465`).

**3. Are error categories documented?**
Weakly. Documentation exists as scattered docstrings: the BaseLLM contract specifies ValueError/TimeoutError/RuntimeError semantics (`llms/base_llm.py:316-318, 353-355`), the A2A module documents its JSON-RPC code ranges (`a2a/errors.py:1-10`), and every exception class carries a docstring. User-facing docs cover the operational knobs (`max_retry_limit` default 2 and `respect_context_window` summarization, `docs/edge/en/concepts/agents.mdx:60-61, 374`) but I found no page that documents the error taxonomy itself. Searches across `docs/edge/en/**/*.mdx` for "error handling" surfaced only per-tool pages (e.g., `docs/edge/en/concepts/tools.mdx`), none describing error categories.

**4. Can new error types be added without breaking existing handling?**
Mostly yes. Dispatch is by concrete-type checks and fall-through chains: a new exception that is not recognized falls into generic branches (log + raise, or unknown-error bucket in the retry loop, `resolver.py:384-386`). The A2A enum reserves the -32768..-32100 range for implementation-defined codes (`a2a/errors.py:9, 69`), making protocol-code extension safe. Two caveats: (a) new LLM providers must be audited against the hardcoded `CONTEXT_LIMIT_ERRORS` phrases — a provider that words context exhaustion differently will be misclassified as a generic error and escalate instead of summarize; (b) because `utilities/exceptions/__init__.py` exports nothing (`utilities/exceptions/__init__.py:1`), consumers cannot rely on one import surface, so new types tend to be discovered by grep rather than by API.

> **Can you tell from the error type whether to retry, escalate, or stop?**
> Yes for A2A calls (`is_retryable_error`, `a2a/errors.py:464-478`) and file uploads (`TransientUploadError` → retry with backoff; `PermanentUploadError` → stop; `resolver.py:375-393`). Partially for the agent loop: parse and generic errors retry up to `max_retry_limit` / `_max_parsing_attempts`, litellm errors escalate immediately, unresumable context overflows stop via `SystemExit` (`agent_utils.py:747-749`). For raw provider faults other than context length, no — the type alone does not say whether the fault is retryable; transport retries are delegated to provider SDKs (`openai/completion.py:209`).

## Architectural Decisions

1. **Federated taxonomies instead of a central one.** Each package owns its error model (`a2a/errors.py`, `utilities/errors.py`, `crewai_files/processing/exceptions.py`, `crewai_tools/aws/bedrock/exceptions.py`). The intended central module `utilities/exceptions/` holds only one exception type and an empty `__init__` (`lib/crewai/src/crewai/utilities/exceptions/__init__.py:1`), signaling the center was never consolidated.
2. **Errors-as-feedback for the reasoning loop.** Tool, parsing, and guardrail failures are deliberately converted to strings injected into the model context (i18n templates at `lib/crewai/src/crewai/translations/en.json:45-56`; wiring at `tools/tool_usage.py:428-437` and `task.py:1310-1324`) so the LLM self-corrects, rather than propagating exceptions to the caller. This is a harness design choice: Python-level taxonomy matters mainly at subsystem boundaries.
3. **Retry policy concentrated at three layers:** provider SDK (`max_retries=2` passed through, `openai/completion.py:209`), agent (`max_retry_limit`, `agent/core.py:707-717`), and bounded self-correction loops (tool attempts, guardrail retries, RPM blocking in `utilities/rpm_controller.py:38-64`).
4. **Observability taxonomy via typed events.** Rather than enriching exceptions with source metadata, CrewAI emits source-typed error events (`AgentExecutionErrorEvent`, `LLMCallFailedEvent`, etc.) before re-raising (`crew.py:1054-1063`; `agent/core.py:696-704`), keeping exceptions simple and pushing classification to consumers of the event bus.
5. **Dataclass-based wire-format errors for protocols.** A2A errors double as JSON-RPC payloads (`to_dict`/`to_response`, `a2a/errors.py:146-162`), so the taxonomy is shaped by the wire protocol rather than by internal handling needs.

## Notable Patterns

- **Classifier functions attached to the exception**: `LLMContextLengthExceededError._is_context_limit_error` is a static method reused both to build and to detect the error (`utilities/exceptions/context_window_exceeding_exception.py:32-44`; called from `utilities/agent_utils.py:698-709`).
- **Multiple-inheritance mixins for orthogonal dimensions**: `TransientUploadError(UploadError, TransientFileError)` composes operation × permanence axes so catch sites can target either axis (`crewai-files/processing/exceptions.py:98-103`).
- **Classification recorded in telemetry**: the retry loop tags each failure `transient` / `permanent` / `unknown` into metrics metadata (`resolver.py:376, 382, 385`), making the taxonomy observable in operations.
- **Per-provider classifier hooks**: Gemini and Bedrock uploaders own their provider→taxonomy mapping functions (`uploaders/gemini.py:48-64`, `uploaders/bedrock.py:28-49`), keeping SDK specifics out of shared code.
- **Model-dependent tuning of error tolerance**: `_max_parsing_attempts` drops from 3 to 2 when the function-calling LLM is a "bigger" OpenAI model (`tools/tool_usage.py:44-62, 114-119`).
- **Passthrough allow-list for foreign exceptions**: `_passthrough_exceptions` tuple plus litellm module-prefix check let third-party errors skip agent-level retries (`agent/core.py:133, 695-706`).

## Tradeoffs

- **String-matching classification (fragility over coverage).** Matching eight English phrases covers major providers today (`context_window_exceeding_exception.py:4-13`) with zero dependency on SDK exception types, but silently degrades when providers reword messages or new providers appear; misclassification escalates instead of summarizes.
- **Module-prefix provider detection (coupling).** `startswith("litellm")` ties retry policy to a specific provider library's packaging (`crew_agent_executor.py:445`); swapping the LLM stack changes retry semantics invisibly.
- **Errors-as-strings (recoverability vs information loss).** Feeding formatted text to the model maximizes self-correction ability but discards the typed structure — downstream code cannot programmatically distinguish `wrong_tool_name` from `tool_usage_exception` except by parsing the sentence.
- **Generic exceptions on exhaustion.** Guardrail retry exhaustion raises bare `Exception` (`task.py:1298-1301`), forcing callers to string-match if they want to handle it specially.
- **Delegated transport retries (simplicity vs observability).** Letting the OpenAI SDK retry internally (`openai/completion.py:209`) keeps CrewAI simple, but those retries are invisible to CrewAI telemetry and cannot be cancelled by higher-level policies.

## Failure Modes / Edge Cases

- **Misclassified provider error**: a non-OpenAI provider phrasing context exhaustion differently than the eight phrases yields a generic-exception path: logged, re-raised, counted against `max_retry_limit`, potentially killing the run instead of triggering summarization (`crew_agent_executor.py:444-458`).
- **Silent swallowing of tool failures**: after `_max_parsing_attempts`, the error becomes a "Moving on then" string and execution continues with degraded capability (`tools/tool_usage.py:426-437`); only verbose mode and events reveal it.
- **Blocking rate limiter**: `RPMController._wait_for_next_minute` sleeps 60s inline (`rpm_controller.py:73-75`) — under async execution this stalls the loop rather than surfacing a retryable rate-limit condition.
- **Unknown uploader errors retried as transient**: any unrecognized exception in `_upload_with_retry` is treated as retryable and burns all backoff attempts before returning `None` (`resolver.py:384-393`), e.g., a deterministic bug in the uploader wastes `UPLOAD_MAX_RETRIES` cycles.
- **Empty-taxonomy import trap**: code importing error types from `crewai.utilities.exceptions` gets only what individual modules leak; the near-empty `__init__` invites inconsistent import paths across the codebase.

## Future Considerations

- Introduce a shared base (e.g., `CrewAIError` with a `source: Literal["model","provider","tool","validation","policy","context","user","infrastructure","timeout"]` field) exported from `utilities/exceptions/__init__.py`, and migrate subsystem hierarchies onto it incrementally.
- Replace message-substring context detection with SDK-typed checks where available (the OpenAI provider already catches typed SDK errors before falling back to the substring check, `openai/completion.py:937-940`); keep the phrase list as last-resort fallback only.
- Replace the litellm module-prefix test with an explicit configurable passthrough exception list (`_passthrough_exceptions` already exists as a seam at `agent/core.py:133` but defaults to empty).
- Publish a docs page enumerating error categories and their handling semantics, mirroring what `docs/edge/en/concepts/agents.mdx:60-61` does for retry knobs.
- Add unit tests for `is_retryable_error` / `is_client_error` and for `classify_upload_error` status-code boundaries (429 vs 500 vs 4xx); no tests were found exercising these predicates (`lib/crewai/tests/a2a/` contains integration tests only).

## Questions / Gaps

- **No evidence found** for a documented, repo-level error taxonomy (searched `docs/edge/**` for "error handling", "retry"; inspected `utilities/exceptions/`; checked AGENTS.md). Documentation of categories is limited to inline docstrings and parameter tables.
- **No evidence found** of tests covering `A2AErrorCode` retryability helpers or `classify_upload_error`; `rg` over `lib/crewai/tests/` and `lib/crewai-files/tests/` for `is_retryable_error` and `TransientUploadError` returned no test-file hits (only production code).
- Whether the `user` error-source category is intentionally absent (folded into validation) or merely unrealized could not be determined from the code.
- The `_passthrough_exceptions` tuple is always empty in-repo (`agent/core.py:133`); no configuration path populating it was found, so its intended extension mechanism appears unused or reserved for external embedders.
- Cross-library consistency (e.g., whether `crewai-tools`' `BedrockValidationError` is ever caught by `lib/crewai` handling code) was not traced end-to-end; within the single-source boundary it appears the tools-lib exceptions propagate uncaught to users.

---

Generated by `dimensions/13.01-error-taxonomy.md` against `crewai`.
