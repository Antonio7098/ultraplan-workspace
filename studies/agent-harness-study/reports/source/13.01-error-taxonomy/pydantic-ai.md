# Source Analysis: pydantic-ai

## Dimension 13.01 — Error Taxonomy

### Source Info

| Field | Value |
|-------|-------|
| Name | pydantic-ai |
| Path | `studies/agent-harness-study/sources/pydantic-ai` |
| Language / Stack | Python (pydantic-core, anyio, httpx; graph engine in `pydantic_graph`) |
| Analyzed | 2026-08-21 |

## Summary

Pydantic AI implements an explicit, class-hierarchy-based error taxonomy centered on a single canonical module, `pydantic_ai_slim/pydantic_ai/exceptions.py`, with all public exceptions exported from the package root (`pydantic_ai_slim/pydantic_ai/__init__.py:32-46`). Errors are classified by **source category** through both inheritance and payload fields:

- **Provider/API errors**: `ModelAPIError` (carries `model_name`) and its subclass `ModelHTTPError` (adds `status_code` and `body`) — `pydantic_ai_slim/pydantic_ai/exceptions.py:236-266`. Every model adapter normalizes SDK exceptions into these types (e.g., `pydantic_ai_slim/pydantic_ai/models/openai.py:189-192`, `models/anthropic.py:250-253`, `models/bedrock.py:123-124`, `models/groq.py:84-87`, `models/xai.py:86-87`, `models/cohere.py:203-204`, `models/google.py:836-841`, `models/mistral.py:109-110`).
- **Model-behavior errors**: `AgentRunError` base (`exceptions.py:181`) → `UnexpectedModelBehavior` (`exceptions.py:203`) with subclasses `ContentFilterError` (`exceptions.py:232`) and `IncompleteToolCall` (`exceptions.py:308`).
- **Policy errors**: `UsageLimitExceeded` (`exceptions.py:195`) and `ConcurrencyLimitExceeded` (`exceptions.py:199`).
- **User/developer errors**: `UserError` (`exceptions.py:159`), plus graph-setup analogues in `pydantic_graph/pydantic_graph/exceptions.py:7-48`.
- **Tool/validation errors**: Pydantic `ValidationError` and user-raised `ModelRetry` are wrapped into the internal `ToolRetryError` carrying a model-visible `RetryPromptPart` (`exceptions.py:273-291`; wrapping at `pydantic_ai_slim/pydantic_ai/tool_manager.py:184-191`).
- **Control-flow exceptions** (not failures): `ModelRetry` (`exceptions.py:40`), `CallDeferred` (`exceptions.py:80`), `ApprovalRequired` (`exceptions.py:98`), and hook-bypass signals `SkipModelRequest` / `SkipToolValidation` / `SkipToolExecution` (`exceptions.py:116-156`).
- **Aggregation**: `FallbackExceptionGroup(ExceptionGroup)` (`exceptions.py:269`) groups all fallback-model failures, including a distinct `ResponseRejected` marker for rejected responses (`models/fallback.py:43-47`).

The taxonomy is actively used for routing: `except` clauses dispatch to retry-node construction, deferral collection, skip-shortcuts, or propagation; `FallbackModel.fallback_on` routes by exception type with `(ModelAPIError,)` as default (`models/fallback.py:91`); capability hook `on_model_request_error` is an explicit retry/suppress/propagate decision point (`capabilities/abstract.py:596-617`). The answer to "can you tell from the error type whether to retry, escalate, or stop?" is largely **yes**: retryable-model-feedback types (`ModelRetry`, `ToolRetryError`) convert into `RetryPromptPart` messages; provider errors escalate to fallback models; policy and limit errors stop the run.

## Rating

**8 / 10** — A clear source-oriented hierarchy with explicit dispatch, dedicated tests (hashability, pickle round-trips), observability semantics per category, durable-execution serialization, and configurable type-based fallback routing. Falls short of 9–10 because there is no consolidated taxonomy documentation table, tool timeouts collapse into `ModelRetry` losing a distinct timeout category, `ToolRetryError` is not publicly exported despite being central to dispatch, and `FallbackExceptionGroup` breaks plain `except ModelAPIError` compatibility (a documented but real sharp edge).

## Evidence Collected

Every entry includes a file path with line numbers relative to the source root `studies/agent-harness-study/sources/pydantic-ai/`.

| Area | Evidence | File:Line |
|------|----------|-----------|
| Error type definitions (canonical module) | `__all__` lists 17 public exceptions; single module owns taxonomy | `pydantic_ai_slim/pydantic_ai/exceptions.py:19-37` |
| Retry control-flow type | `ModelRetry` with message + pydantic serialization schema (`kind: 'model-retry'`) | `pydantic_ai_slim/pydantic_ai/exceptions.py:40-77` |
| Deferral control-flow types | `CallDeferred`, `ApprovalRequired` carry optional metadata | `pydantic_ai_slim/pydantic_ai/exceptions.py:80-113` |
| Hook-skip types | `SkipModelRequest`, `SkipToolValidation`, `SkipToolExecution` | `pydantic_ai_slim/pydantic_ai/exceptions.py:116-156` |
| User-error category | `UserError(RuntimeError)` — "usage mistake by the application developer" | `pydantic_ai_slim/pydantic_ai/exceptions.py:159-167` |
| Run-error base | `AgentRunError(RuntimeError)` base for run failures | `pydantic_ai_slim/pydantic_ai/exceptions.py:181-192` |
| Policy category | `UsageLimitExceeded(AgentRunError)`, `ConcurrencyLimitExceeded(AgentRunError)` | `pydantic_ai_slim/pydantic_ai/exceptions.py:195-200` |
| Provider category | `ModelAPIError` (has `model_name`); `ModelHTTPError(ModelAPIError)` adds `status_code`, `body` | `pydantic_ai_slim/pydantic_ai/exceptions.py:236-266` |
| Model-behavior category | `UnexpectedModelBehavior` keeps `body`; subclasses `ContentFilterError`, `IncompleteToolCall` | `pydantic_ai_slim/pydantic_ai/exceptions.py:203-233,308-309` |
| Fallback aggregation | `FallbackExceptionGroup(ExceptionGroup[Any])`; `ResponseRejected` for rejected responses | `pydantic_ai_slim/pydantic_ai/exceptions.py:269-270`; `pydantic_ai_slim/pydantic_ai/models/fallback.py:43-47` |
| Tool-retry wrapper | `ToolRetryError` wraps `RetryPromptPart`; formats validation details for the model | `pydantic_ai_slim/pydantic_ai/exceptions.py:273-305` |
| Provider normalization (OpenAI) | SDK errors → `ModelHTTPError` (4xx/5xx) or `ModelAPIError` | `pydantic_ai_slim/pydantic_ai/models/openai.py:189-192,1940-1943` |
| Provider normalization (Anthropic/Bedrock/Groq/xAI/Cohere/Mistral/Gemini/Google/OpenRouter/HuggingFace) | same pattern across adapters | `pydantic_ai_slim/pydantic_ai/models/anthropic.py:250-253`; `models/bedrock.py:123-124`; `models/groq.py:84-87`; `models/xai.py:86-87`; `models/cohere.py:203-204`; `models/mistral.py:109-110`; `models/gemini.py:277`; `models/google.py:836-841,1509-1514`; `models/openrouter.py:976-992,1156`; `models/huggingface.py:83` |
| Embeddings use same taxonomy | `ModelAPIError` raised from embedding client | `pydantic_ai_slim/pydantic_ai/embeddings/voyageai.py:172` |
| Policy enforcement raising typed errors | `UsageLimits.check_requests/check_tokens/check_tool_calls` raise `UsageLimitExceeded` | `pydantic_ai_slim/pydantic_ai/usage.py:380-418` |
| Token-limit check wired into run loop | `usage_limits.check_tokens(ctx.state.usage)` after each response | `pydantic_ai_slim/pydantic_ai/_agent_graph.py:1023-1024` |
| Validation→retry wrapping | `ValidationError`/`ModelRetry` → `ToolRetryError(RetryPromptPart)` | `pydantic_ai_slim/pydantic_ai/tool_manager.py:184-191,282,341-344,462` |
| Per-tool retry budget → stop | `_check_max_retries` raises `UnexpectedModelBehavior('Tool ... exceeded max retries ...')` chained via `from error` | `pydantic_ai_slim/pydantic_ai/tool_manager.py:177-181` |
| Global output-retry budget → stop | `consume_output_retry` raises `UnexpectedModelBehavior('Exceeded maximum output retries (N)')` | `pydantic_ai_slim/pydantic_ai/_agent_graph.py:162-179` |
| Truncation-specific error | `check_incomplete_tool_call` raises `IncompleteToolCall` with remediation hint | `pydantic_ai_slim/pydantic_ai/_agent_graph.py:146-160` |
| Model-retry routing in agent loop | `except exceptions.ModelRetry` → `_build_retry_node` builds `ModelRequestNode` with `RetryPromptPart` | `pydantic_ai_slim/pydantic_ai/_agent_graph.py:703-714,747-753,1027-1040` |
| Output-validation retry routing | `except ToolRetryError` → append `e.tool_retry` as next request part | `pydantic_ai_slim/pydantic_ai/_agent_graph.py:1145-1149,1243-1249,1782-1784` |
| Deferral routing | `except CallDeferred` → external queue; `except ApprovalRequired` → approvals queue | `pydantic_ai_slim/pydantic_ai/_agent_graph.py:1889-1894` |
| Skip short-circuits | `except SkipModelRequest` uses provided response instead of model call | `pydantic_ai_slim/pydantic_ai/_agent_graph.py:617,788-794` |
| Hook error-routing contract | `on_model_request_error`: raise = propagate, return `ModelResponse` = suppress, raise `ModelRetry` = retry; skipped for `SkipModelRequest`/`ModelRetry` | `pydantic_ai_slim/pydantic_ai/capabilities/abstract.py:596-617` |
| Type-based fallback routing | `fallback_on` accepts exception types/handlers/response handlers; default `(ModelAPIError,)`; group raise when exhausted | `pydantic_ai_slim/pydantic_ai/models/fallback.py:33-41,91,118-148,157-165,225-244,317-327` |
| Output-pipeline wrap logic | `except ToolRetryError: raise` (already wrapped); `except (ValidationError, ModelRetry)` → wrap as retry prompt | `pydantic_ai_slim/pydantic_ai/_output.py:172-179,207-210,219-226` |
| Streaming output failure classification | `UnexpectedModelBehavior('Output validation failed during streaming...')` chained from `ValidationError`/`ModelRetry` | `pydantic_ai_slim/pydantic_ai/result.py:245-309` |
| Timeout handling | tool timeout → `ModelRetry(f'Timed out after {timeout} seconds.')` (category collapsed into model-retry) | `pydantic_ai_slim/pydantic_ai/toolsets/function.py:651-656` |
| Observability per category | `CallDeferred`/`ApprovalRequired` recorded as control flow (span attr, not ERROR unless v<5); `ToolRetryError` recorded as ERROR span with retry prompt attribute | `pydantic_ai_slim/pydantic_ai/capabilities/instrumentation.py:353-383` |
| Durable-execution serialization | Temporal toolset serializes/deserializes `ModelRetry`, `ApprovalRequired`, `CallDeferred` | `pydantic_ai_slim/pydantic_ai/durable_exec/temporal/_toolset.py:53-79,126-141` |
| CLI surface | `except UserError` for friendly CLI error rendering | `pydantic_ai_slim/pydantic_ai/_cli/__init__.py:296` |
| Model-visible retry message part | `RetryPromptPart` documents six trigger sources (validation, tool retry, unknown tool, text-instead-of-structured, etc.) | `pydantic_ai_slim/pydantic_ai/messages.py:1407-1444` |
| Tests: hashability | parametrized over 11 exception classes | `tests/test_exceptions.py:29-70` |
| Tests: pickle round-trip | attribute preservation incl. `status_code`, `body`, `metadata`, `tool_retry` | `tests/test_exceptions.py:73-140` |
| Tests: retry-prompt formatting | serialized ErrorDetails formatting edge cases | `tests/test_exceptions.py:143-183` |
| Docs: API reference page | mkdocstrings page for `pydantic_ai.exceptions` | `docs/api/exceptions.md:1-3` |
| Docs: model-errors guide | `capture_run_messages` + `__cause__` inspection pattern | `docs/agent.md:1203-1231` |
| Docs: fallback & exception-group caveats | default `ModelAPIError` trigger; `except*` guidance for `FallbackExceptionGroup` | `docs/models/overview.md:180,329-346,465-499` |
| Docs: retry budgets | per-tool vs output budget semantics | `docs/tools-advanced.md:476-492`; `docs/output.md:555` |
| Docs: transport-layer HTTP retries | separate tenacity-based layer distinct from taxonomy | `docs/retries.md:1-9,328-342` |

## Answers to Dimension Questions

### 1. Are errors classified by source?

Yes. The hierarchy in `pydantic_ai_slim/pydantic_ai/exceptions.py:159-266` encodes source categories directly:

- **user**: `UserError` ("caused by a usage mistake by the application developer", `exceptions.py:159`)
- **provider/infrastructure**: `ModelAPIError` / `ModelHTTPError` with `model_name`, `status_code`, `body` (`exceptions.py:236-266`), raised uniformly by every model adapter after catching SDK exceptions (`models/openai.py:189-192`, `models/anthropic.py:250-253`, et al.)
- **model behavior**: `UnexpectedModelBehavior` and subclasses `ContentFilterError`, `IncompleteToolCall` (`exceptions.py:203-232,308`)
- **policy**: `UsageLimitExceeded`, `ConcurrencyLimitExceeded` (`exceptions.py:195-200`)
- **tool/validation**: Pydantic `ValidationError` normalized into `ToolRetryError` carrying `RetryPromptPart` (`tool_manager.py:184-191`)
- **control flow (not errors)**: `ModelRetry`, `CallDeferred`, `ApprovalRequired`, `Skip*` (`exceptions.py:40-156`) — instrumentation explicitly treats deferrals as non-errors (`capabilities/instrumentation.py:358-361`).

Timeout is the one dimension-category handled implicitly: tool timeouts are converted to `ModelRetry` with a textual message rather than a dedicated type (`toolsets/function.py:655-656`). Context-window/token-limit truncation *does* get a dedicated type (`IncompleteToolCall`, `_agent_graph.py:146-160`).

### 2. Is the taxonomy used for handling?

Extensively. Dispatch sites route on exact exception types:

- **Retry**: `except ModelRetry` → `_build_retry_node` constructs the next `ModelRequestNode` with a `RetryPromptPart` (`_agent_graph.py:1027-1040`); `except ToolRetryError` feeds `e.tool_retry` back as the next request (`_agent_graph.py:1782-1784`); validation failures wrap into `ToolRetryError` (`tool_manager.py:571-574`).
- **Defer/suspend**: `except CallDeferred` / `except ApprovalRequired` classify calls into external vs approval queues (`_agent_graph.py:1889-1894`).
- **Escalate/fallback**: `FallbackModel._should_fallback` runs registered exception-type handlers against each failure; default is `(ModelAPIError,)` so any provider API error escalates to the next model (`models/fallback.py:91,157-165,225-244`). Custom handlers can match arbitrary types or even inspect responses.
- **Stop**: policy errors propagate (`usage.py:380-418` via `_agent_graph.py:1023-1024`); exhausted retry budgets stop the run with `UnexpectedModelBehavior` carrying the original error as `__cause__` (`tool_manager.py:177-181`, `_agent_graph.py:162-179`).
- **Hook-level decision point**: `on_model_request_error` formalizes the retry/suppress/propagate choice for any `Exception` (`capabilities/abstract.py:596-617`), giving users a single interception point before errors leave the model-request node.

So yes — the error type effectively encodes the disposition: model-feedback types retry, provider types escalate to fallback, policy/behavior types stop.

### 3. Are error categories documented?

Documented, though scattered rather than consolidated:

- Every class has a docstring rendered into the API reference (`docs/api/exceptions.md:1-3` via mkdocstrings).
- `docs/agent.md:1203-1231` documents `UnexpectedModelBehavior` diagnosis via `capture_run_messages` and `__cause__`.
- `docs/models/overview.md:180,276-330,465-499` documents fallback triggers, `FallbackExceptionGroup` contents, and the `except*` caveat for ExceptionGroups.
- `docs/tools-advanced.md:476-492` documents per-tool retry limits and the resulting `UnexpectedModelBehavior` message; `docs/output.md:555` documents the output retry budget.
- `RetryPromptPart` enumerates its six trigger sources (`messages.py:1410-1419`).
- There is **no single taxonomy reference table** mapping source → type → default disposition; understanding requires reading several pages.

### 4. Can new error types be added without breaking existing handling?

Largely yes. The mechanism is ordinary Python exception subclassing plus `isinstance`-based handlers:

- `FallbackOn` accepts any `type[Exception]` tuple, handler callables, or mixed sequences (`models/fallback.py:33-41`), so new types integrate with fallback routing without library changes.
- `on_model_request_error` receives generic `Exception` (`capabilities/abstract.py:601`), so custom capabilities can react to novel errors.
- New exceptions don't require registry updates; unhandled types simply propagate, which is safe-by-default.
- Serialization support must be added manually per class (`__reduce__` methods at `exceptions.py:94-95,222-223,246-247,265-266,285-286`; pydantic schema at `exceptions.py:61-77`), and `tests/test_exceptions.py:73-126` pins pickle round-tripping — extension requires following this convention but nothing breaks if a new type skips it until needed.
- Risk: broad catches like `except BaseException` span recording (`capabilities/instrumentation.py:380-383`) handle unknowns generically, and the graph's exhaustive `assert_never` pattern applies to message parts (`_agent_graph.py:1202-1203`), not exceptions, so no exhaustiveness gate blocks additions.

## Architectural Decisions

1. **Single canonical exceptions module with explicit `__all__`** (`exceptions.py:19-37`, re-exported at `__init__.py:32-46`): one import point for the entire taxonomy; internal-only helpers like `ToolRetryError` stay out of `__all__`.
2. **Inheritance-based source classification instead of an enum/kind field**: categories form a tree rooted at `Exception`/`RuntimeError`/`AgentRunError`. Payload fields (`model_name`, `status_code`, `body`) carry machine-readable context where needed. Tradeoff: no uniform `.kind` discriminator (only `RetryPromptPart.part_kind='retry-prompt'` at `messages.py:1443` and `ModelRetry`'s `'model-retry'` serialization kind at `exceptions.py:67`).
3. **Control flow expressed as exceptions**: `ModelRetry`, `CallDeferred`, `ApprovalRequired`, and `Skip*` reuse the exception channel for non-failure signaling (`exceptions.py:40-156`). This makes hook composition simple (raise to divert) at the cost of blurring "error" vs "signal" — mitigated by explicit comments ("Control flow, not error" at `_output.py:208`) and instrumentation special-casing (`capabilities/instrumentation.py:353-371`).
4. **Errors converted to model-visible messages at a boundary**: `ValidationError`/`ModelRetry` never reach the user raw mid-run; they become `RetryPromptPart` content inside `ToolRetryError` (`tool_manager.py:184-191`), keeping the LLM conversation self-healing while budgets last.
5. **Two independent retry layers**: taxonomy-driven model retries (per-tool counters + global output budget) versus transport-driven HTTP retries via tenacity (`docs/retries.md:1-9`), deliberately kept separate — the skill docs warn not to conflate them (`pydantic_ai_slim/pydantic_ai/.agents/skills/building-pydantic-ai-agents/references/TOOLS-ADVANCED.md:110-112`).
6. **Type-configurable escalation**: `FallbackModel.fallback_on` turns taxonomy membership into a declarative routing policy (`models/fallback.py:91,118-148`).

## Notable Patterns

- **Adapter-normalized provider errors**: every `models/{provider}.py` maps SDK exceptions to `ModelHTTPError`/`ModelAPIError` at one choke point (e.g., `models/openai.py:189-192`), so downstream code never sees vendor exception types. This follows the repo's stated rule of raising explicit errors rather than letting SDK types leak (`models/AGENTS.md` Error Handling rules).
- **Cause chaining for diagnosability**: terminal errors chain their origin — `_check_max_retries(...) from error` (`tool_manager.py:181`), streaming validation failure `from e` (`result.py:304-308`) — enabling the documented `e.__cause__` inspection pattern (`docs/agent.md:1229-1230`).
- **Category-aware telemetry**: span status differs by category — deferrals set attributes instead of ERROR status; retries record the exact prompt the model will see (`capabilities/instrumentation.py:353-379`).
- **Durable-execution round-tripping**: exceptions crossing the Temporal workflow boundary are reified as dataclasses and reconstructed (`durable_exec/temporal/_toolset.py:53-79,130-141`), proving the taxonomy survives serialization — rare among agent frameworks.
- **Budget accounting tied to taxonomy**: consuming a retry unit distinguishes who erred (tool vs output validator vs empty response) and only then decides continue-vs-stop (`_agent_graph.py:162-179,1145-1149`).

## Tradeoffs

- **Exceptions-as-control-flow** simplifies hooks but means `try/except Exception` in user code around `agent.run()` can accidentally intercept deferral signals; the library compensates by documenting which types are control flow (`logfire.md:333` notes spans deliberately don't mark them ERROR).
- **`FallbackExceptionGroup` breaks flat catches**: after exhaustion, individual `ModelAPIError`s are wrapped in an `ExceptionGroup`, so pre-existing `except ModelAPIError` handlers miss them; users must adopt `except*` or unwrap manually (`docs/models/overview.md:465-499`). This is documented but remains the taxonomy's sharpest edge.
- **Timeout category loss**: converting `TimeoutError` to `ModelRetry` (`toolsets/function.py:655-656`) means programmatic handling cannot distinguish "tool timed out" from "model asked to correct args" except by string matching.
- **Hierarchy depth vs catch breadth**: catching `AgentRunError` gets you usage + behavior + API errors but not `ModelRetry` (deliberately outside that tree), forcing users to know which subtree they care about.
- **No numeric/machine-readable severity or retryability hint on provider errors**: `ModelHTTPError.status_code` exists, but 429-vs-500 retryability decisions live in the tenacity transport config, not the taxonomy (`docs/retries.md:36-56`).

## Failure Modes / Edge Cases

- **Retry-budget exhaustion surfaces as model-behavior error**: when a tool exceeds `max_retries`, the run stops with `UnexpectedModelBehavior('Tool {name!r} exceeded max retries count of {N}')` and the triggering `ModelRetry`/validation error as `__cause__` (`tool_manager.py:177-181`) — operators must inspect the cause to find the root category.
- **Streaming cannot retry**: output validation failures during streaming become immediate `UnexpectedModelBehavior` because retries are unsupported in `run_stream()` (`result.py:304-309`); partial-output mode is the escape hatch.
- **Truncated tool calls detected late-but-deterministically**: `finish_reason == 'length'` with invalid final tool-call args triggers `IncompleteToolCall` with an actionable remediation message about `max_tokens` (`_agent_graph.py:146-160`).
- **Empty fallback handler set**: `FallbackModel` rejects an empty effective `fallback_on` with `UserError` explaining the consequence (`models/fallback.py:143-148`), preventing silently-no-fallback misconfiguration.
- **All-models-failure aggregation**: when every fallback candidate fails, all exceptions plus any rejected responses are combined into one `FallbackExceptionGroup` (`models/fallback.py:317-327`), preserving full history for diagnosis at the cost of catch simplicity.
- **Deferral without output-type support**: deferred calls arriving when `DeferredToolRequests` isn't among output types produce a prescriptive `UserError` (`result.py:271-275`, `_agent_graph.py:1850-1855`).

## Future Considerations

- Add a stable, machine-readable category accessor (e.g., `error.kind`) or a documented taxonomy table mapping source → type → default disposition (retry / defer / fallback / stop), consolidating what is currently spread across `docs/agent.md`, `docs/models/overview.md`, `docs/tools-advanced.md`, and docstrings.
- Introduce a dedicated `ToolTimeoutError` (or add `timeout=True` metadata to the retry prompt) so timeouts remain distinguishable programmatically while still feeding the model a retry.
- Export `ToolRetryError` (or a public alias) since it is central to dispatch and appears in documented `__cause__` chains.
- Consider attaching parsed provider error codes / retryability hints to `ModelAPIError` so transport-independent consumers can make retry decisions without tenacity coupling.
- Provide an `except*`-friendly helper for unwrapping `FallbackExceptionGroup` given the documented catch-compatibility pitfall.

## Questions / Gaps

- No evidence found of a centralized error-code registry or structured logging schema keyed by taxonomy category; search covered `pydantic_ai_slim/pydantic_ai/**` for `class .*Error|class .*Retry` definitions and `docs/**` for taxonomy references — categorization lives solely in the class hierarchy.
- Whether `ContentFilterError`'s empty-response-with-metadata recovery path (raising from `details.get('finish_reason')` at `_agent_graph.py:1112-1130`) is exercised for all providers could not be confirmed from a single test file; coverage appears distributed across provider VCR tests (e.g., `tests/models/`).
- Rate-limit (429) handling has no taxonomy-level special-casing (no `RateLimitError`); it relies on `status_code` inspection by user code or transport config — intentional per the two-layers design, but worth noting as a gap versus some peer frameworks.

---

Generated by Dimension 13.01 (Error Taxonomy) against `pydantic-ai`.
