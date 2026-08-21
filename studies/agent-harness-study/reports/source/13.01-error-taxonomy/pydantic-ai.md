# Source Analysis: pydantic-ai

## 13.01 Error Taxonomy

### Source Info

| Field | Value |
|-------|-------|
| Name | pydantic-ai |
| Path | `studies/agent-harness-study/sources/pydantic-ai` |
| Language / Stack | Python 3.10–3.13 (uv workspace, pydantic-ai-slim core + pydantic-graph + pydantic-evals) |
| Analyzed | 2026-08-21 |

## Summary

pydantic-ai has a deliberately small, deeply structured exception taxonomy that is the spine of its control-flow / error-routing model. All public exception classes live in `pydantic_ai_slim/pydantic_ai/exceptions.py` (16 named classes in `__all__` plus `ToolRetryError`). Errors are **classified by source category** (model behavior, model infrastructure, validation, tool, policy/limits, user/developer, fallback, hook timeout, control-flow skip/defer/approval) via a shallow inheritance hierarchy rooted at `RuntimeError` for run-impacting errors and at `Exception` for control-flow signals. The taxonomy is actively used for routing in the agent graph (`pydantic_ai_slim/pydantic_ai/_agent_graph.py`), the tool manager (`pydantic_ai_slim/pydantic_ai/tool_manager.py`), and the output pipeline (`pydantic_ai_slim/pydantic_ai/_output.py`), and is complemented by an extensibility seam: capability hooks (`on_*_error`) that take a typed exception and may recover the run, retry, or re-raise. The taxonomy is documented via per-class docstrings and tested for hashability and pickling (`tests/test_exceptions.py`), but the dispatcher behavior is verified indirectly through VCR-based tests rather than dedicated dispatcher tests.

## Rating

**Score: 7/10 — Clear model with tests, explicit interfaces, and operational safeguards.**

Rationale:
- A canonical, narrow file (`pydantic_ai_slim/pydantic_ai/exceptions.py:19-309`) defines a clear hierarchy aligned to source categories; categories map 1:1 to dispatchers (model → `_map_api_errors`, tool → `tool_manager.py`, validation → `_output.py`, policy → `usage.py`/`concurrency.py`).
- Routing is centralized in a few well-marked sites: `_agent_graph.py:617, 700-743, 788, 823-829, 1247-1249, 1647-1654, 1782-1784, 2023-2024`; `tool_manager.py:282-285, 330-350, 458-465, 571-574, 665-669, 734-736, 760-765`; `_output.py:107-226`.
- `ModelRetry` and `ToolRetryError` make "retry" an explicit, inspectable signal carrying a `RetryPromptPart`, not a side-channel.
- Tests exist for the public surface (`tests/test_exceptions.py:58, 118, 129, 143, 150, 163`); per-class docstrings document intent.
- The hook contracts (`capabilities/abstract.py:443-924`) define a typed error-handling surface a user can extend.
- Extensibility is **open** for user error types (`UserError`, `AgentRunError`, `ModelAPIError` are subclass-friendly `RuntimeError` roots), and **closed** at dispatchers (no `assert_never`-style exhaustiveness over the exception union).

Why not higher:
- `UserError`, `AgentRunError`, `MCPError` extend `RuntimeError`, so `except Exception:` does NOT catch them — a documented but easy-to-miss footgun (`pydantic_ai_slim/pydantic_ai/exceptions.py:159, 181`; `pydantic_ai_slim/pydantic_ai/mcp.py:143`).
- `ToolRetryError` and `ContentFilterError` are not re-exported from the top-level `pydantic_ai/__init__.py` `__all__` (lines 191-207), creating an asymmetric public API (`pydantic_ai_slim/pydantic_ai/exceptions.py:34-36`).
- There is **no framework-level retry policy keyed on `ModelHTTPError.status_code`**; retries happen either at the HTTP transport layer (`pydantic_ai_slim/pydantic_ai/retries.py:117, 215, 312`) or via `FallbackModel` (`pydantic_ai_slim/pydantic_ai/models/fallback.py:91, 231-269`), never inside the agent graph.
- The two near-identical "wrap validation error as retry prompt" helpers (`pydantic_ai_slim/pydantic_ai/_output.py:121` and `pydantic_ai_slim/pydantic_ai/tool_manager.py:184`) duplicate logic that could drift.
- The `_check_max_retries` boundary in `tool_manager.py:177-181` uses equality (`==`) while `_agent_graph.py:177-181` uses `>`; two boundary conventions coexist.
- No central "convert any error to model-readable" function; retry-prompt formatting is scattered.
- `BaseException` sweeping at `_agent_graph.py:83, 684, 729, 1950`, `run.py:217`, `agent/__init__.py:1594, 1635, 1661, 1676` is correct for cancellation but catches `SystemExit` too.

## Evidence Collected

Every entry includes a file path with line numbers, format `path/to/file.ts:NN`. All paths are workspace-relative under `studies/agent-harness-study/sources/pydantic-ai/`.

| Area | Evidence | File:Line |
|------|----------|-----------|
| Canonical exception file | All 17 public/internal exception classes defined here | `pydantic_ai_slim/pydantic_ai/exceptions.py:19-309` |
| `__all__` (exceptions module) | `ModelRetry, CallDeferred, ApprovalRequired, SkipModelRequest, SkipToolValidation, SkipToolExecution, UserError, UndrainedPendingMessagesError, AgentRunError, UnexpectedModelBehavior, UsageLimitExceeded, ConcurrencyLimitExceeded, ModelAPIError, ModelHTTPError, ContentFilterError, IncompleteToolCall, FallbackExceptionGroup` | `pydantic_ai_slim/pydantic_ai/exceptions.py:19-37` |
| Top-level re-export asymmetry | `ToolRetryError` and `ContentFilterError` NOT re-exported by top-level `__all__` (lines 191-207) | `pydantic_ai_slim/pydantic_ai/__init__.py:191-207` vs `pydantic_ai_slim/pydantic_ai/exceptions.py:34-36` |
| Model-behavior control signal | `ModelRetry` with `__get_pydantic_core_schema__` for (de)serialization | `pydantic_ai_slim/pydantic_ai/exceptions.py:40-77` |
| Deferred/approval/skip control signals | `CallDeferred`, `ApprovalRequired`, `SkipModelRequest`, `SkipToolValidation`, `SkipToolExecution` | `pydantic_ai_slim/pydantic_ai/exceptions.py:80-156` |
| User/developer error root | `UserError(RuntimeError)` and `UndrainedPendingMessagesError(UserError)` | `pydantic_ai_slim/pydantic_ai/exceptions.py:159-178` |
| Agent-run error root | `AgentRunError(RuntimeError)` | `pydantic_ai_slim/pydantic_ai/exceptions.py:181-192` |
| Policy/limits errors | `UsageLimitExceeded`, `ConcurrencyLimitExceeded` (both `AgentRunError`) | `pydantic_ai_slim/pydantic_ai/exceptions.py:195-200` |
| Model-behavior error | `UnexpectedModelBehavior(AgentRunBehavior)` with `body` field | `pydantic_ai_slim/pydantic_ai/exceptions.py:203-229` |
| Model-behavior subclasses | `ContentFilterError(UnexpectedModelBehavior)`, `IncompleteToolCall(UnexpectedModelBehavior)` | `pydantic_ai_slim/pydantic_ai/exceptions.py:232-233, 308-309` |
| Model infrastructure errors | `ModelAPIError(AgentRunError)`, `ModelHTTPError(ModelAPIError)` | `pydantic_ai_slim/pydantic_ai/exceptions.py:236-266` |
| Fallback group | `FallbackExceptionGroup(ExceptionGroup[Any])` | `pydantic_ai_slim/pydantic_ai/exceptions.py:269-270` |
| Internal retry-prompt carrier | `ToolRetryError(Exception)` with `tool_retry: RetryPromptPart` | `pydantic_ai_slim/pydantic_ai/exceptions.py:273-286` |
| Validation formatting helper | `_format_error_details` (manual, not `ValidationError.from_exception_data`) | `pydantic_ai_slim/pydantic_ai/exceptions.py:288-305` |
| MCP-specific runtime error | `MCPError(RuntimeError)` outside main taxonomy | `pydantic_ai_slim/pydantic_ai/mcp.py:143` |
| Hook timeout error | `HookTimeoutError(TimeoutError)` | `pydantic_ai_slim/pydantic_ai/capabilities/hooks.py:64` |
| Internal fallback rejection | `ResponseRejected` (private) | `pydantic_ai_slim/pydantic_ai/models/fallback.py:43-47` |
| Per-provider error mapping — OpenAI | `_map_api_errors` context manager (HTTPError vs APIError split) | `pydantic_ai_slim/pydantic_ai/models/openai.py:184` |
| Per-provider error mapping — Anthropic | `_map_api_errors` | `pydantic_ai_slim/pydantic_ai/models/anthropic.py:245` |
| Per-provider error mapping — Mistral | `_map_api_errors` (uses `SDKError`) | `pydantic_ai_slim/pydantic_ai/models/mistral.py:104` |
| Per-provider error mapping — Groq / Cohere / Bedrock / xAI / HuggingFace | `_map_api_errors` context managers | `pydantic_ai_slim/pydantic_ai/models/groq.py:79`, `cohere.py:62`, `bedrock.py:117`, `xai.py:79`, `huggingface.py:79` |
| Provider-specific re-raise — Google inline | `except errors.APIError` (different shape) | `pydantic_ai_slim/pydantic_ai/models/google.py:834, 1507-1514` |
| Provider-specific re-raise — Gemini / OpenRouter inline | `except ... → ModelHTTPError/ModelAPIError` | `pydantic_ai_slim/pydantic_ai/models/gemini.py:276`, `pydantic_ai_slim/pydantic_ai/models/openrouter.py:976, 992, 1156` |
| Status-code gate inside provider | xAI maps `grpc.StatusCode` → HTTP code | `pydantic_ai_slim/pydantic_ai/models/xai.py:79-87, 94, 96` |
| Skip dispatch in agent graph | `SkipModelRequest` from `_prepare_request` → empty stream + `_finish_handling` | `pydantic_ai_slim/pydantic_ai/_agent_graph.py:617` |
| Skip dispatch in model request | `SkipModelRequest` from `_prepare_request` → skip model call entirely | `pydantic_ai_slim/pydantic_ai/_agent_graph.py:788` |
| Wrap-task recovery dispatcher | `wrap_task.result()` → `_resolve_wrap_result` (decides Skip / ModelRetry / on_model_request_error) | `pydantic_ai_slim/pydantic_ai/_agent_graph.py:700-743, 999-1013` |
| Wrap-model-request routing | `SkipModelRequest` → use response; `ModelRetry` → re-raise; `Exception` → `on_model_request_error` capability; retry via `_build_retry_node` | `pydantic_ai_slim/pydantic_ai/_agent_graph.py:823-829, 1027-1040` |
| Output-retry budget | `if output_retries_used > max_output_retries: raise UnexpectedModelBehavior` | `pydantic_ai_slim/pydantic_ai/_agent_graph.py:177-181` |
| Output tool retry (per-call) | `consume_output_retry` returns new node carrying retry prompt | `pydantic_ai_slim/pydantic_ai/_agent_graph.py:162-179` |
| Output tool retry — first dispatch | `ToolRetryError` → emit `ModelRequestNode(retry prompt)` | `pydantic_ai_slim/pydantic_ai/_agent_graph.py:1145-1149, 1247-1249` |
| Output tool retry — second dispatch | increment `output_retries_used`, append `e.tool_retry` to `output_parts` | `pydantic_ai_slim/pydantic_ai/_agent_graph.py:1647-1654` |
| Deferred-tool retry | `ToolRetryError` from invalidated deferred-tool call | `pydantic_ai_slim/pydantic_ai/_agent_graph.py:1782-1784` |
| Deferred-result resolver | `ToolRetryError` raised inside `_call_tool` (deferred-result resolution) → returns `tool_retry` part | `pydantic_ai_slim/pydantic_ai/_agent_graph.py:2023-2024` |
| Incomplete tool call upgrade | `args_as_dict(raise_if_invalid=True)` failures re-raised as `IncompleteToolCall` | `pydantic_ai_slim/pydantic_ai/_agent_graph.py:157-160` |
| Cancellation handling | `BaseException` catch with re-raise of original node error | `pydantic_ai_slim/pydantic_ai/_agent_graph.py:83-86, 684, 729, 1950` |
| Tool validate error routing | `(ValidationError, ModelRetry)` → `on_tool_validate_error` | `pydantic_ai_slim/pydantic_ai/tool_manager.py:282-285` |
| Tool execute central dispatch | control flow (re-raise), `ModelRetry` (propagate), `Exception` → `on_tool_execute_error`, outer `(ValidationError, ModelRetry)` → `_wrap_error_as_retry` | `pydantic_ai_slim/pydantic_ai/tool_manager.py:330-350` |
| Tool validation dispatcher | `SkipToolValidation` → use `e.validated_args`; `(ValidationError, ModelRetry)` → `_make_validation_failure` | `pydantic_ai_slim/pydantic_ai/tool_manager.py:458-465` |
| Output-tool validate dispatch | `(ToolRetryError, ValidationError, ModelRetry)` → `_make_validation_failure` | `pydantic_ai_slim/pydantic_ai/tool_manager.py:571-574` |
| Output-process hooks dispatch | `ToolRetryError` → `_check_max_retries` then re-raise | `pydantic_ai_slim/pydantic_ai/tool_manager.py:665-669` |
| Skip tool execution escape hatch | `SkipToolExecution` → use `e.result` | `pydantic_ai_slim/pydantic_ai/tool_manager.py:734-736` |
| Raw tool body ModelRetry | `ModelRetry` from tool body → `_wrap_error_as_retry` (becomes `ToolRetryError`) | `pydantic_ai_slim/pydantic_ai/tool_manager.py:760-765` |
| Per-tool retry budget boundary | `if retries[name] == max_retries: raise UnexpectedModelBehavior` | `pydantic_ai_slim/pydantic_ai/tool_manager.py:177-181` |
| Retry counter increment | `self.ctx.retries[failed_tool_name] + 1` | `pydantic_ai_slim/pydantic_ai/tool_manager.py:127, 348, 407, 763` |
| Validation→retry-prompt converter (tool layer) | `_wrap_error_as_retry` | `pydantic_ai_slim/pydantic_ai/tool_manager.py:184-191` |
| Validation→retry-prompt converter (output layer) | `_make_retry_prompt` | `pydantic_ai_slim/pydantic_ai/_output.py:121-129` |
| Output validate hooks dispatch | `run_output_validate_hooks` (`ToolRetryError` pass-through; `(ValidationError, ModelRetry)` → convert if `wrap_validation_errors=True`) | `pydantic_ai_slim/pydantic_ai/_output.py:107-179` |
| Output process hooks dispatch | `run_output_process_hooks` (`ToolRetryError` / `ModelRetry` pass-through; `Exception` → `on_output_process_error`; outer `(ValidationError, ModelRetry)` → `_make_retry_prompt`) | `pydantic_ai_slim/pydantic_ai/_output.py:182-226` |
| Output validator with hooks | `run_output_with_hooks` (`ModelRetry` → `ToolRetryError` if `wrap_validation_errors=True`) | `pydantic_ai_slim/pydantic_ai/_output.py:392-404` |
| Streaming partial-validation handling | `stream_output` swallows `(ValidationError, ModelRetry)` | `pydantic_ai_slim/pydantic_ai/result.py:90-93` |
| Streaming short-circuit on failure | `run_stream` raises `UnexpectedModelBehavior('...retries are not supported...')` | `pydantic_ai_slim/pydantic_ai/result.py:304-309` |
| Run-level error capture | `ErrorMarker` from graph iteration; pending-messages check raises `UndrainedPendingMessagesError` | `pydantic_ai_slim/pydantic_ai/run.py:135, 217, 226, 250` |
| Node wrap recovery | `except Exception` from `wrap_node_run` → `on_node_run_error` capability | `pydantic_ai_slim/pydantic_ai/run.py:296-301` |
| Run-level recovery chain | `wrap_run` → `on_run_error` → `after_run` (with `BaseException` capture + re-raise) | `pydantic_ai_slim/pydantic_ai/agent/__init__.py:1594, 1629-1691` |
| Hook contracts (typed error surface) | `on_run_error`, `on_node_run_error`, `on_model_request_error`, `on_tool_validate_error`, `on_tool_execute_error`, `on_output_validate_error`, `on_output_process_error` | `pydantic_ai_slim/pydantic_ai/capabilities/hooks.py:891, 932, 990, 1047, 1106` |
| Hook protocol defaults (re-raise) | All `on_*_error` hooks default to "do nothing" → re-raise | `pydantic_ai_slim/pydantic_ai/capabilities/abstract.py:443-924` (e.g. 463, 596, 617, 684, 742) |
| `on_model_request_error` ModelRetry contract | "raise `ModelRetry` to retry the request" — only this hook may | `pydantic_ai_slim/pydantic_ai/capabilities/abstract.py:611-612` |
| Fallback model dispatch | `_should_fallback`; default `fallback_on=(ModelAPIError,)`; final `FallbackExceptionGroup` | `pydantic_ai_slim/pydantic_ai/models/fallback.py:91, 231-269, 317-327` |
| Usage-limit enforcement | `UsageLimitExceeded` raise sites (request_tokens / response_tokens / etc.) | `pydantic_ai_slim/pydantic_ai/usage.py:383, 387, 393, 401, 405, 411, 418` |
| Concurrency-limit enforcement | `ConcurrencyLimitExceeded` raise site | `pydantic_ai_slim/pydantic_ai/concurrency.py:193` |
| Approval toolset | `ApprovalRequired` emitted on approval-required tools | `pydantic_ai_slim/pydantic_ai/toolsets/approval_required.py:30` |
| Temporal durable-exec bridge | `ApprovalRequired` / `CallDeferred` re-raised for durable workflows | `pydantic_ai_slim/pydantic_ai/durable_exec/temporal/_toolset.py:137, 139` |
| HTTP retry transport (tenacity) | `TenacityTransport`, `AsyncTenacityTransport`, `wait_retry_after` | `pydantic_ai_slim/pydantic_ai/retries.py:117, 215, 312` |
| Tests — hashability | `test_exceptions_hashable` for 11 classes | `tests/test_exceptions.py:58-71` |
| Tests — pickle round-trip | `test_exceptions_pickle_round_trip`, `test_tool_retry_error_pickle_round_trip` | `tests/test_exceptions.py:118-141` |
| Tests — `ToolRetryError` formatting | string content, ErrorDetails, value_error ctx stripping | `tests/test_exceptions.py:143-183` |
| Example of docs link from code | `capture_run_messages` docstring → `../agent.md#model-errors` | `pydantic_ai_slim/pydantic_ai/_agent_graph.py:2068` |

## Answers to Dimension Questions

1. **Are errors classified by source?**
   Yes. The taxonomy in `pydantic_ai_slim/pydantic_ai/exceptions.py:19-309` explicitly partitions errors into:
   - **Model behavior**: `UnexpectedModelBehavior` (line 203) with `ContentFilterError` (line 232) and `IncompleteToolCall` (line 308) subclasses
   - **Model infrastructure**: `ModelAPIError` (line 236) with `ModelHTTPError` (line 250) subclass
   - **Validation/tool/output retry**: `ModelRetry` (line 40) and the internal `ToolRetryError` (line 273)
   - **Approval/deferral**: `CallDeferred` (line 80), `ApprovalRequired` (line 98)
   - **Skip/control-flow**: `SkipModelRequest` (line 116), `SkipToolValidation` (line 133), `SkipToolExecution` (line 146)
   - **Policy/limits**: `UsageLimitExceeded` (line 195), `ConcurrencyLimitExceeded` (line 199)
   - **User/developer**: `UserError` (line 159), `UndrainedPendingMessagesError` (line 170)
   - **Run-level base**: `AgentRunError` (line 181)
   - **Fallback**: `FallbackExceptionGroup` (line 269)
   - **Hook timeout**: `HookTimeoutError` (`capabilities/hooks.py:64`)
   - **MCP**: `MCPError` (`mcp.py:143`)

2. **Is the taxonomy used for handling?**
   Yes. Three dispatch sites route based on the taxonomy:
   - **Agent graph** (`pydantic_ai_slim/pydantic_ai/_agent_graph.py:617, 700-743, 788, 823-829, 1145-1149, 1247-1249, 1647-1654, 1782-1784, 2023-2024`): converts `ModelRetry` → retry node, `ToolRetryError` → retry prompt part, `SkipModelRequest` → empty stream, `IncompleteToolCall` re-raise from `args_as_dict`.
   - **Tool manager** (`pydantic_ai_slim/pydantic_ai/tool_manager.py:177-181, 282-285, 330-350, 458-465, 571-574, 665-669, 734-736, 760-765`): centralizes `(ValidationError, ModelRetry) → ToolRetryError` conversion via `_wrap_error_as_retry` (line 184-191); `SkipTool*` and `CallDeferred`/`ApprovalRequired` are re-raised as control flow; per-tool retry budget via `_check_max_retries`.
   - **Output pipeline** (`pydantic_ai_slim/pydantic_ai/_output.py:107-179, 182-226, 392-404`): `run_output_validate_hooks`, `run_output_process_hooks`, and `run_output_with_hooks` all dispatch on `(ValidationError, ModelRetry, ToolRetryError, Exception)` and convert per the `wrap_validation_errors` policy.
   - The per-provider `_map_api_errors` context managers (`models/openai.py:184`, `models/anthropic.py:245`, `models/mistral.py:104`, `models/groq.py:79`, `models/cohere.py:62`, `models/bedrock.py:117`, `models/xai.py:79`, `models/huggingface.py:79`, inline in `models/google.py:834, 1507-1514`, `models/gemini.py:276`, `models/openrouter.py:976, 992, 1156`) convert provider SDK exceptions into `ModelHTTPError` (status >= 400) or `ModelAPIError` (connection/other) — the only "classify by source" step at the model boundary.
   - `FallbackModel` (`pydantic_ai_slim/pydantic_ai/models/fallback.py:91, 231-269, 317-327`) uses `fallback_on=(ModelAPIError,)` by default to decide retry-vs-bubble across a chain.

3. **Are error categories documented?**
   Yes, but mostly **inline**. Each exception class in `pydantic_ai_slim/pydantic_ai/exceptions.py` carries a one- or multi-line docstring (lines 40-46, 80-88, 98-106, 116-124, 133-137, 146-150, 159-160, 170-178, 181-182, 195, 199-200, 203-204, 232-233, 236-237, 250-251, 269-270, 273-274, 308-309). The `_format_error_details` helper carries an explanatory comment at lines 292-295. The agent graph references user-facing docs from `_agent_graph.py:2068` (`../agent.md#model-errors`). The hook contracts at `pydantic_ai_slim/pydantic_ai/capabilities/abstract.py:443-924` document each `on_*_error` recovery semantics in their Protocol methods. **No central taxonomy diagram or "errors" doc page was found in the slim package** — the user-facing documentation lives under `docs/` (outside the searched slim scope).

4. **Can new error types be added without breaking existing handling?**
   Mostly yes, but with caveats:
   - **Open** at the roots: `UserError`, `AgentRunError`, `ModelAPIError`, `UnexpectedModelBehavior`, `ModelRetry`, `MCPError` are all `RuntimeError`/`Exception`-rooted; user code can subclass without breaking existing `except` clauses (Python's polymorphic `except` will match).
   - **Closed at dispatchers**: dispatchers use **explicit isinstance chains and multi-`except` lists**, not `assert_never` over an exception union (verified — no `assert_never` over exception classes was found). Adding a new subclass of `AgentRunError` will be propagated as-is; adding a new control-flow signal (e.g. a third `Skip*` exception) requires editing the matching `except` in `_agent_graph.py`, `tool_manager.py`, or `_output.py`. No exhaustive checking will warn at type-check time.
   - The capability hook surface (`on_*_error`) is **typed `error: Exception` (or narrower)**, so a new error class with no `except` clause will fall through to the default re-raise — not break, but also not be handled.
   - Two near-duplicate conversion helpers (`pydantic_ai_slim/pydantic_ai/tool_manager.py:184-191` and `pydantic_ai_slim/pydantic_ai/_output.py:121-129`) must be updated in tandem if `RetryPromptPart` shape changes.

## Architectural Decisions

1. **Single canonical exceptions file** (`pydantic_ai_slim/pydantic_ai/exceptions.py:1-309`). Keeping all exceptions in one place is unusual for a large framework but pays off for discoverability and stable public API. Trade-off: every new exception must be added here.

2. **`RuntimeError` as the root for user-facing errors** (`pydantic_ai_slim/pydantic_ai/exceptions.py:159, 181`). Forces users to either `except UserError` / `except AgentRunError` explicitly or use `except RuntimeError` (which is too broad). The docstrings are clear that these are "your fault" / "framework problem", not control flow. Cost: a `try/except Exception:` does NOT catch `UserError` — see `MCPError` at `mcp.py:143` for a third sibling.

3. **`Exception` (not `RuntimeError`) for control-flow signals** (`pydantic_ai_slim/pydantic_ai/exceptions.py:40, 80, 98, 116, 133, 146, 273`). `ModelRetry`, `CallDeferred`, `ApprovalRequired`, `SkipModelRequest`, `SkipToolValidation`, `SkipToolExecution`, `ToolRetryError` are `Exception` subclasses and are used as carrier/control-flow types, not errors. The agent graph catches them by `except` clauses; user code is expected to let them propagate.

4. **Provider boundary = single `_map_api_errors` per provider**, with status-code split into `ModelHTTPError` (≥400) vs `ModelAPIError` (other). The decision to **not** keep provider-specific subclasses (`OpenAIHTTPError`, etc.) keeps the public taxonomy small and lets `FallbackModel` work uniformly across providers — at the cost of providers sharing the same exception type even when their semantics differ.

5. **`ModelRetry` as the explicit retry signal**. Raising `ModelRetry` from a tool/validator is the only way to "ask the model to try again"; the framework converts it to a `RetryPromptPart` (`pydantic_ai_slim/pydantic_ai/_agent_graph.py:1027-1040`, `pydantic_ai_slim/pydantic_ai/tool_manager.py:184-191`). This makes "retry" a value, not a side channel, and is reflected in the `__get_pydantic_core_schema__` at `exceptions.py:62-77` (enables cross-process serialization).

6. **`ToolRetryError` as internal carrier vs `ModelRetry` as public signal**. `ToolRetryError` (`exceptions.py:273`) is an `Exception`-rooted value type carrying a `RetryPromptPart`; it's an implementation detail of the tool/output dispatch chain. It is intentionally NOT re-exported in `pydantic_ai/__init__.py` `__all__` (lines 191-207), so user code uses `ModelRetry` and the framework handles the rest.

7. **Capability hooks as the extensibility seam** (`pydantic_ai_slim/pydantic_ai/capabilities/hooks.py:891, 932, 990, 1047, 1106`). Each `on_*_error` takes a typed exception and may return a recovery value. This is where pydantic-ai's design diverges from "exceptions are handled in user code" — the framework promotes errors into a typed data model and offers explicit recovery points.

8. **Fallback as a model-level concern, not an agent-level concern** (`pydantic_ai_slim/pydantic_ai/models/fallback.py:91, 231-269, 317-327`). `FallbackModel.fallback_on=(ModelAPIError,)` is the default; an exhausted fallback raises a `FallbackExceptionGroup` carrying every individual error. This places the "model A failed, try model B" decision at the model abstraction, not at the agent-graph level.

9. **`ModelHTTPError` carries `(status_code, model_name, body)` as a typed struct** (`pydantic_ai_slim/pydantic_ai/exceptions.py:250-266`) rather than a stringly-typed message. Cost: status-code-aware retry logic could easily be implemented but isn't at the framework level.

10. **HTTP retry is opt-in via `TenacityTransport`** (`pydantic_ai_slim/pydantic_ai/retries.py:117, 215, 312`). Status-code-driven retry (`wait_retry_after`) is implemented but only at the transport layer. The framework does NOT silently retry `ModelHTTPError`.

## Notable Patterns

- **Carrier exception as value**: `ModelRetry`, `ToolRetryError`, `Skip*`, `CallDeferred`, `ApprovalRequired` are all `Exception` subclasses used as typed data carriers between graph nodes. They are caught and re-raised within the dispatchers (`pydantic_ai_slim/pydantic_ai/_agent_graph.py:617, 823-829`, `pydantic_ai_slim/pydantic_ai/tool_manager.py:330-350, 458-465, 734-736`) and never surface to user code.
- **Context-manager error mapping**: `_map_api_errors(model_name)` (`pydantic_ai_slim/pydantic_ai/models/openai.py:184` and 7 sibling providers) wraps provider SDK calls and translates native exceptions into the framework's taxonomy.
- **Recovery chain**: `wrap_*` → `on_*_error` → `after_*` for each capability layer (`pydantic_ai_slim/pydantic_ai/agent/__init__.py:1629-1691`). At the top of the chain, `wrap_run` → `on_run_error` → `after_run` is the user's last recovery point.
- **Polymorphic shape**: `ModelHTTPError.__init__(status_code, model_name, body=None)` with a derived message string; `UnexpectedModelBehavior.__init__(message, body=None)` with auto-JSON-pretty-print of body (`pydantic_ai_slim/pydantic_ai/exceptions.py:211-220`).
- **`__reduce__` for pickling**: every exception with non-trivial state implements `__reduce__` (`pydantic_ai_slim/pydantic_ai/exceptions.py:94-95, 112-113, 222-223, 246-247, 265-266, 285-286`). This is what enables durable-exec / cross-process resume (`durable_exec/temporal/_toolset.py:137, 139`) and is verified by `tests/test_exceptions.py:118-141`.
- **Hashable exceptions**: every named exception overrides `__hash__` and `__eq__` (`pydantic_ai_slim/pydantic_ai/exceptions.py:55-59, 58-59`) so they can be used as dict keys / set members; verified at `tests/test_exceptions.py:58-71`.
- **Stricter per-tool retry boundary (`==`) vs looser per-output budget (`>`)**: two boundary conventions coexist in the same dispatch chain (`pydantic_ai_slim/pydantic_ai/tool_manager.py:177-181` vs `pydantic_ai_slim/pydantic_ai/_agent_graph.py:177-181`).
- **Capability defaults are do-nothing re-raise** (`pydantic_ai_slim/pydantic_ai/capabilities/abstract.py:443, 463, 596, 617, 684, 742`). Hooks opt in by overriding; the framework's default is identical to "no hook registered".

## Tradeoffs

| Decision | Pro | Con |
|---|---|---|
| Single canonical exceptions file | Discoverable, stable API, easy to review | Every new feature touches one file; can become a god-module |
| `UserError`/`AgentRunError` rooted at `RuntimeError` | Distinguishes framework failures from generic exceptions; matches Python stdlib `RuntimeError` semantics | `except Exception:` does NOT catch them; user code must know the type or use `except RuntimeError:` |
| No per-provider error subclasses | Uniform `FallbackModel` semantics across providers; small public surface | Provider-specific retry/escalation logic cannot be expressed in the type system; must inspect `ModelHTTPError.status_code` |
| No framework-level retry on `ModelHTTPError.status_code` | Predictable — `run()` raises what the provider raised | Users must wire `TenacityTransport` or `FallbackModel` themselves; no defaults |
| `ToolRetryError`/`ContentFilterError` not in top-level `__all__` | Keeps "internal carrier" vs "public signal" distinction | Asymmetric; users discover the omission only when an `import` fails |
| `assert_never` is **not** used over the exception union | Adding a new error class does not require touching every dispatcher | Adding a new control-flow signal CAN silently fail to be caught; tests are the only safety net |
| `BaseException` sweeps in graph cleanup (`_agent_graph.py:83`, `run.py:217`, `agent/__init__.py:1594, 1635, 1661, 1676`) | Correct for `CancelledError` | Also catches `SystemExit`; reasoning surface is large |
| Two near-duplicate "wrap validation as retry prompt" helpers (`tool_manager.py:184-191`, `_output.py:121-129`) | Each lives next to its caller; locality of reasoning | Drift risk if `RetryPromptPart` shape changes |

## Failure Modes / Edge Cases

- **Streaming partial validation**: `result.py:90-93` swallows `(ValidationError, ModelRetry)` silently — user must rely on the final `get_output()` to detect failure. `result.py:304-309` explicitly raises `UnexpectedModelBehavior('...retries are not supported in run_stream()')` — the only place the framework advertises that streaming does not retry.
- **`_check_max_retries` boundary uses equality, not `>=`** (`tool_manager.py:177-181`). If a future refactor changes to `>=`, the per-tool budget breaks without a test catching it.
- **`BaseException` in `agent/__init__.py:1635-1643, 1660-1691`**: catches `KeyboardInterrupt`/`SystemExit`/`asyncio.CancelledError`; the code re-raises the original node error, but if the wrap-recovery itself raises, the new exception is propagated. Two comment blocks document the intent (`agent/__init__.py:1636-1638, 1664-1666`).
- **`undrained pending messages`** (`run.py:226`): bare `async for node in agent_run` skips drain of `'when_idle'` messages; users see `UndrainedPendingMessagesError` instead of their messages being processed.
- **`models/groq.py:221` catches `ModelHTTPError`** — the only provider that does so inside the slim package. Behavior differs from the other `_map_api_errors` providers; intent not documented in source.
- **`models/openai.py:386` Azure 400 special-case**: `if system == 'azure' and e.status_code == 400 ...` — could fire after `_map_api_errors` already raised `ModelHTTPError` at line 188; the actual code path is provider-specific and not centralized.
- **`MCPError` is not exported** in `pydantic_ai/__init__.py` `__all__` but is a `RuntimeError` sibling to `UserError`; users importing `from pydantic_ai import MCPError` will fail.
- **`tool_manager.py:629-646` and `tool_manager.py:864`** document the two-stage `wrap_validation_errors` contract — `ToolRetryError` at one boundary may be raw `ValidationError` at another; tests must follow this carefully.
- **`HookTimeoutError` only wraps user-supplied hook timeouts** (`capabilities/hooks.py:64, 243`); not raised by the agent graph for general infrastructure timeouts (those come through `ModelAPIError`/`ModelHTTPError`).

## Future Considerations

- **Consolidate `tool_manager.py:184-191` (`_wrap_error_as_retry`) and `_output.py:121-129` (`_make_retry_prompt`)** into a single helper in `exceptions.py` — these implement the same `ValidationError | ModelRetry → ToolRetryError(RetryPromptPart)` conversion.
- **Add `assert_never` exhaustiveness over an explicit `AgentRunError` union** in `_agent_graph.py:823-829` and `tool_manager.py:330-350`. Today, adding a new error subclass is silent at type-check time; today dispatchers use multi-`except` lists.
- **Re-export `ToolRetryError` and `ContentFilterError`** from `pydantic_ai/__init__.py` `__all__` (lines 191-207), or explicitly document why they are excluded.
- **Unify per-tool (`==`) and per-output (`>`) retry boundaries** (`tool_manager.py:177-181` vs `_agent_graph.py:177-181`).
- **Consider an opt-in framework-level status-code retry policy** based on `ModelHTTPError.status_code` (e.g., 429/503 → retry with backoff inside the agent graph). Today this only exists at the transport layer (`retries.py`).
- **Add a central "errors" doc page** in `docs/` (outside the slim scope), referenced from `_agent_graph.py:2068`'s `../agent.md#model-errors`. Today the taxonomy is documented inline in `exceptions.py` and scattered through `capabilities/abstract.py`; no central reference exists.
- **Add a dedicated dispatcher-behavior test file** (analogous to `tests/test_exceptions.py` for construction). Today dispatcher behavior is verified through VCR-based integration tests in `tests/test_agent.py` and `tests/models/`.
- **Make `MCPError` and `HookTimeoutError` discoverable** from the top-level `pydantic_ai/__init__.py` `__all__`.
- **Document the `RuntimeError` parent implication** for `UserError`, `AgentRunError`, `MCPError` in the exceptions module docstring or in a top-level `errors.md` doc — this is a documented but easy-to-miss footgun.

## Questions / Gaps

- **Where is the user-facing "errors" documentation page?** `_agent_graph.py:2068` references `../agent.md#model-errors`, but no `errors.md` was found under `pydantic_ai_slim/pydantic_ai/`. The user-facing docs live in `docs/` (outside the slim scope).
- **Is the `models/groq.py:221` `except ModelHTTPError` intentional, and what is its dispatch logic?** The other providers catch native SDK errors inside `_map_api_errors` and never see `ModelHTTPError` again. No clear evidence found in the slim package for groq's intent.
- **What is the relationship between `models/openai.py:386` Azure 400 special-case and the `_map_api_errors` raise at line 188?** No clear evidence found that the two paths don't overlap.
- **Why are `ToolRetryError` and `ContentFilterError` deliberately omitted from the top-level `__all__`?** The asymmetry is implicit; no source comment explains the exclusion.
- **Is there a framework-wide retry policy keyed on `ModelHTTPError.status_code` planned or in-progress?** No evidence found in the slim package; the hook contract at `capabilities/abstract.py:611-612` permits `on_model_request_error` to raise `ModelRetry`, but there is no automatic policy.
- **What is the dispatch policy for `HookTimeoutError`?** `capabilities/hooks.py:243` raises it; the consumer at the agent-graph level was not explicitly identified in the slim package. No clear evidence found for a dedicated `except HookTimeoutError` site.
