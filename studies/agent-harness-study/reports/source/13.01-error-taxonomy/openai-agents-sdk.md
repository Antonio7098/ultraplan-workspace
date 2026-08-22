# Source Analysis: openai-agents-sdk

## Dimension 13.01: Error Taxonomy

### Source Info

| Field | Value |
|-------|-------|
| Name | openai-agents-sdk |
| Path | `studies/agent-harness-study/sources/openai-agents-sdk` |
| Language / Stack | Python (Agents SDK for LLM agents; OpenAI Responses/Chat Completions APIs) |
| Analyzed | 2026-08-21 |

## Summary

The SDK classifies errors along three complementary axes rather than one global enum. (1) The core run-loop taxonomy is a Python class hierarchy rooted in `AgentsException` (`src/agents/exceptions.py:46`) with leaf classes that map cleanly to source categories: model (`ModelBehaviorError` at `src/agents/exceptions.py:66`, `ModelRefusalError` at `src/agents/exceptions.py:78`), user/developer misuse (`UserError` at `src/agents/exceptions.py:89`), policy/guardrails (four tripwire exceptions at `src/agents/exceptions.py:121-174`), tool/timeout (`ToolTimeoutError` at `src/agents/exceptions.py:109`, `MCPToolCancellationError` at `src/agents/exceptions.py:99`), and orchestration (`MaxTurnsExceeded` at `src/agents/exceptions.py:56`). (2) The sandbox subsystem has a far richer machine-readable taxonomy: a stable string `ErrorCode` enum (`src/agents/sandbox/errors.py:12-55`), an operation literal (`OpName`, `src/agents/sandbox/errors.py:58-73`), category base classes (`ConfigurationError`, `SandboxRuntimeError`, `ArtifactError`, `SnapshotError` at `src/agents/sandbox/errors.py:111-124`), and an explicit `retryable: bool | None` field (`src/agents/sandbox/errors.py:95`) that directly answers retry-vs-stop. (3) Provider/model errors are not wrapped in the SDK hierarchy; instead the retry layer normalizes raw OpenAI-SDK exceptions into `ModelRetryNormalizedError` facts (`status_code`, `error_code`, `is_network_error`, `is_timeout`, `retry_after` — `src/agents/retry.py:50-60`) and composable `RetryPolicy` objects (`src/agents/retry.py:231-361`) decide retry behavior from those facts.

Errors are actively used for routing: `RunErrorHandlers` is a dict keyed by error kind (`"max_turns"`, `"model_refusal"`) with isinstance-based dispatch (`src/agents/run_error_handlers.py:50-54`, `src/agents/run_internal/error_handlers.py:137-140`); guardrail tripwires terminate runs; tool errors are by default converted to model-visible strings so the run continues (`src/agents/tool.py:1609-1618`); and every `AgentsException` gets a `RunErrorDetails` snapshot attached for observability (`src/agents/run.py:1505-1517`, `src/agents/run_internal/run_loop.py:1175-1187`). The taxonomy is documented in `docs/running_agents.md:553-564` and exercised by dedicated tests (`tests/test_exception_exports.py`, `tests/test_run_internal_error_handlers.py`, `tests/models/test_model_retry.py` with 55 tests, `tests/sandbox/test_errors.py`).

## Rating

**7 / 10.** The core model is clear, tested, documented, and dispatch is isinstance-based so new exception types can be added without breaking handlers. The sandbox taxonomy is genuinely mature (stable machine-readable codes, per-error retryable classification, trace/event surfacing). It falls short of 9-10 because of concrete fragility and gaps: `terminal_metadata_for_exception` classifies exceptions by matching class *name strings* instead of isinstance checks (`src/agents/sandbox/memory/rollouts.py:183-188`); provider errors propagate raw (e.g., `ResponsesWebSocketError` extends `RuntimeError`, outside the SDK hierarchy — `src/agents/models/openai_responses.py:314`); the public docs' Exceptions overview omits `ModelRefusalError`, `MCPToolCancellationError`, and both tool-guardrail tripwire types (`docs/running_agents.md:553-564`); and core `AgentsException` classes carry no retryable/escalate metadata — that capability exists only in the sandbox and retry-policy subsystems.

## Evidence Collected

| Area | Evidence | File:Line |
|------|----------|-----------|
| Error base class | `AgentsException` is the documented base for all SDK exceptions; carries `run_data` slot | `src/agents/exceptions.py:46-53` |
| Model-source errors | `ModelBehaviorError` (unexpected model output, malformed JSON), `ModelRefusalError` (carries `refusal` text) | `src/agents/exceptions.py:66-75`, `src/agents/exceptions.py:78-86` |
| User/developer errors | `UserError` — "raised when the user makes an error using the SDK" | `src/agents/exceptions.py:89-96` |
| Policy/guardrail errors | `InputGuardrailTripwireTriggered`, `OutputGuardrailTripwireTriggered`, `ToolInputGuardrailTripwireTriggered`, `ToolOutputGuardrailTripwireTriggered`, each carrying the guardrail result | `src/agents/exceptions.py:121-174` |
| Tool/timeout errors | `ToolTimeoutError` carries `tool_name` + `timeout_seconds`; `MCPToolCancellationError` for cancelled MCP calls | `src/agents/exceptions.py:99-118` |
| Error context payload | `RunErrorDetails` dataclass (input, new_items, raw_responses, last_agent, guardrail results) attached to any `AgentsException` on failure | `src/agents/exceptions.py:30-43`; `src/agents/run.py:1505-1517` |
| Streamed-path parity | Same `RunErrorDetails` attachment in the streaming loop's `except AgentsException` branch | `src/agents/run_internal/run_loop.py:1175-1187` |
| Handler registry keyed by kind | `RunErrorHandlers` TypedDict with `max_turns` / `model_refusal` keys | `src/agents/run_error_handlers.py:50-54` |
| Dispatch by error source | `resolve_run_error_handler_result` uses `isinstance(error, ModelRefusalError)` to pick the handler key | `src/agents/run_internal/error_handlers.py:137-140` |
| Max-turns raise + handler resolution | Non-streamed and streamed paths both build `RunErrorData` and resolve handlers before raising `MaxTurnsExceeded` | `src/agents/run.py:1058-1081`; `src/agents/run_internal/run_loop.py:881-912` |
| Refusal detection → typed error | Model refusal extracted from message items, routed to handler or raised as `ModelRefusalError` | `src/agents/run_internal/turn_resolution.py:764-789` |
| Guardrail raise sites | Input/output tripwires raised at guardrail execution and at multiple stream-loop points | `src/agents/run_internal/guardrails.py:139`, `src/agents/run_internal/guardrails.py:174`; `src/agents/run_internal/run_loop.py:711`, `:983`, `:1217`, `:1266`, `:1278` |
| Tool guardrail raise sites | Tool input/output guardrail tripwires raised during tool execution | `src/agents/run_internal/tool_execution.py:2364`, `src/agents/run_internal/tool_execution.py:2401` |
| MCP cancellation classification | Cancelled MCP tool call converted to `MCPToolCancellationError` | `src/agents/mcp/util.py:698-702` |
| Timeout taxonomy with behavior switch | `asyncio.TimeoutError` → `ToolTimeoutError`; `timeout_behavior="raise_exception"` fails the run, default `"error_as_result"` returns a model-visible message | `src/agents/tool.py:1823-1847`; behavior enum at `src/agents/tool.py:439-447` |
| Tool errors as model-visible results | `default_tool_error_function` returns a retryable-in-place message; JSON decode errors get a distinct message | `src/agents/tool.py:1609-1618` |
| Tool error kinds for formatter | `ToolErrorFormatterArgs.kind` is a closed `Literal["approval_rejected", "tool_not_found"]` | `src/agents/run_config.py:70-92` |
| Tool-not-found policy | `tool_not_found_behavior` default `"raise_error"` (`ModelBehaviorError`), opt-in `"return_error_to_model"` | `src/agents/run_config.py:332`; docs `docs/running_agents.md:158-159`, `:196` |
| Sandbox machine-readable codes | `ErrorCode(str, Enum)` — 27 stable codes (exec, workspace, git, mount, snapshot, config) | `src/agents/sandbox/errors.py:12-55` |
| Sandbox structured error | `SandboxError` carries `message`, `error_code`, `op`, `context`, `cause`, `retryable: bool \| None` | `src/agents/sandbox/errors.py:77-108` |
| Retryable inheritance | `__post_init__` propagates `retryable` from a chained `SandboxError` cause | `src/agents/sandbox/errors.py:97-99` |
| Retryable surfaced operationally | Trace span data and session events expose `error_retryable` | `src/agents/sandbox/session/sandbox_session.py:354-356`, `:483`; `src/agents/sandbox/session/events.py:73` |
| Provider error normalization | `ModelRetryNormalizedError` facts: `status_code`, `error_code`, `request_id`, `retry_after`, `is_abort`, `is_network_error`, `is_timeout` | `src/agents/retry.py:50-60` |
| Retry decision object | `RetryDecision` with `retry`, `delay`, `reason`, internal `_hard_veto` and `_approves_replay` flags | `src/agents/retry.py:114-122` |
| Composable retry policies | `retry_policies.never/provider_suggested/network_error/retry_after/http_status/all/any` — classification-driven retry routing | `src/agents/retry.py:231-361` |
| OpenAI provider advice | `get_openai_retry_advice` honors `x-should-retry` header, classifies `APIConnectionError`/`APITimeoutError` via the error chain, retries 408/409/429/5xx | `src/agents/models/_openai_retry.py:125-210` |
| Normalization merge point | `_normalize_retry_error` merges provider advice into normalized facts; `_evaluate_retry` applies abort/unsafe-replay veto | `src/agents/run_internal/model_retry.py:205-242`, `:370-392` |
| Terminal-state classification | Rollout metadata maps exception to `max_turns_exceeded`/`guardrail_tripped`/`cancelled`/`failed` — by class-name string matching | `src/agents/sandbox/memory/rollouts.py:175-196` |
| Tracing integration | Errors attached to spans as `SpanError`; helper decides which types skip generic span errors | `src/agents/util/_error_tracing.py:7-16`; `src/agents/run_internal/run_loop.py:261-265` |
| Trace redaction | Tool errors redacted in traces unless `trace_include_sensitive_data` | `src/agents/util/_tool_errors.py:3-8` |
| Public exports | All 11 core exception classes re-exported from `agents` | `src/agents/__init__.py:20-33` |
| Docs: exceptions overview | Documented taxonomy: `AgentsException`, `MaxTurnsExceeded`, `ModelBehaviorError`, `ToolTimeoutError`, `UserError`, input/output tripwires | `docs/running_agents.md:553-564` |
| Docs: error handlers | `error_handlers` documented as "dict keyed by error kind" with worked examples | `docs/running_agents.md:468-526` |
| Docs: automatic retries | `conversation_locked` errors auto-retried with backoff | `docs/running_agents.md:420` |
| Tests: export completeness | `test_all_public_exception_classes_are_reexported` guards the public error surface | `tests/test_exception_exports.py:16` |
| Tests: handler dispatch | Async handlers, validation failures, wrapped payloads | `tests/test_run_internal_error_handlers.py:87`, `:75`, `:44` |
| Tests: retry classification | 55 tests covering retry policy composition, stateful replay safety, provider-hint interactions | `tests/models/test_model_retry.py:156-483` |
| Tests: error context | Run and streamed runs attach error data | `tests/test_run_error_details.py:12`, `:31` |
| Tests: sandbox taxonomy | `tests/sandbox/test_errors.py` (5 tests) exercises error construction/codes | `tests/sandbox/test_errors.py` |

## Answers to Dimension Questions

**1. Are errors classified by source?**
Yes, by class hierarchy and by structured fields. Core: `ModelBehaviorError`/`ModelRefusalError` = model; `UserError` = developer misuse; guardrail tripwires = policy; `ToolTimeoutError`/`MCPToolCancellationError` = tool; `MaxTurnsExceeded` = orchestration budget (`src/agents/exceptions.py:56-174`). Sandbox: category base classes (`ConfigurationError` = user config, `SandboxRuntimeError` = infrastructure, `ArtifactError`/`SnapshotError` = materialization/persistence) plus a stable `ErrorCode` enum (`src/agents/sandbox/errors.py:111-128`). Provider/infrastructure/timeout are classified by *normalization facts* (`is_network_error`, `is_timeout`, `status_code`) rather than exception types (`src/agents/retry.py:50-60`, `src/agents/models/_openai_retry.py:137-142`). One gap: provider errors themselves are not wrapped into the `AgentsException` hierarchy — raw OpenAI SDK exceptions propagate to the caller, and `ResponsesWebSocketError` extends `RuntimeError` (`src/agents/models/openai_responses.py:314`).

**2. Is the taxonomy used for handling?**
Yes, extensively. isinstance dispatch selects error handlers (`src/agents/run_internal/error_handlers.py:137-140`); tripwire types stop the run at dedicated raise sites (`src/agents/run_internal/guardrails.py:139`, `:174`); tool errors are converted to model-visible strings so runs continue by default, with an explicit opt-in to raise (`src/agents/tool.py:1834-1847`); retry policies branch on normalized facts (`is_network_error or is_timeout` at `src/agents/retry.py:260-268`; status-code sets at `src/agents/retry.py:285-296`); and sandbox `retryable` is surfaced into trace/event data for operational decisions (`src/agents/sandbox/session/sandbox_session.py:354-356`). The core question — "can you tell from the error type whether to retry, escalate, or stop?" — is answered *yes* for sandbox errors (`retryable` field), *yes* for provider errors (retry policies + `x-should-retry`/`replay_safety`), and *partially* for core run errors: `MaxTurnsExceeded`/`ModelRefusalError` route to escalation handlers, tripwires stop, but ordinary `ModelBehaviorError`/`UserError` carry no retry/escalate metadata.

**3. Are error categories documented?**
Yes, but incompletely. `docs/running_agents.md:553-564` documents the base class, `MaxTurnsExceeded`, `ModelBehaviorError` (with malformed-JSON and tool-misuse sub-cases), `ToolTimeoutError`, `UserError`, and both agent-level tripwires. Error handlers are documented as "dict keyed by error kind" with examples (`docs/running_agents.md:468-526`), and tool error kinds (`"approval_rejected"`, `"tool_not_found"`) at `docs/running_agents.md:218`. Omissions: `ModelRefusalError`, `MCPToolCancellationError`, `ToolInputGuardrailTripwireTriggered`, and `ToolOutputGuardrailTripwireTriggered` are absent from the Exceptions overview (though `ModelRefusalError` is covered in the handlers section), and the sandbox `ErrorCode` enum has no dedicated doc page.

**4. Can new error types be added without breaking existing handling?**
Mostly yes. All core dispatch is isinstance-based with fallback behavior: an unknown `AgentsException` subclass still gets `RunErrorDetails` attached (`src/agents/run.py:1507`) and simply finds no handler (`error_handlers.get(...)` returns `None`, `src/agents/run_internal/error_handlers.py:141-142`). The sandbox hierarchy is open for subclassing since `SandboxError` is a dataclass base accepting any `ErrorCode`/`op`. Two closed seams: `RunErrorHandlers` is a `TypedDict` with exactly two keys (`src/agents/run_error_handlers.py:50-54`) and `ToolErrorFormatterArgs.kind` is a closed `Literal` (`src/agents/run_config.py:73`) — new kinds there require SDK changes. Also, the string-matching classifier (`"Guardrail" in exc_name`, `src/agents/sandbox/memory/rollouts.py:185`) would misclassify any new exception whose class name happens to contain "Guardrail" without behaving like one, and misses subclasses only if they rename.

## Architectural Decisions

- **Hierarchy over enum for the core SDK.** Errors are Python exception subclasses of `AgentsException` (`src/agents/exceptions.py:46`) rather than a coded enum, enabling `except` clauses and isinstance dispatch, but sacrificing machine-readable codes at the core level.
- **Structured codes where machines need them.** The sandbox subsystem deliberately chose `ErrorCode(str, Enum)` with a documented contract ("Stable, machine-readable error codes", `src/agents/sandbox/errors.py:13`) plus `op` and `context` fields — an explicit design for programmatic handling (`src/agents/sandbox/errors.py:78-88`).
- **Normalize, don't wrap, provider errors.** Rather than wrapping OpenAI SDK exceptions, the retry layer extracts facts into `ModelRetryNormalizedError` and lets composable policies decide (`src/agents/retry.py:50-60`, `src/agents/retry.py:231-361`). This keeps the raw exception (and its chain) intact for the caller while enabling classification (`src/agents/models/_openai_retry.py:14-20` walks the error chain).
- **Recoverable-by-default tool errors.** Tool failures and timeouts default to `error_as_result` behavior — the model sees the error string and the run continues — with `raise_exception` as an explicit opt-in (`src/agents/tool.py:439-447`, `:1834`).
- **Error context as a first-class payload.** `RunErrorDetails`/`RunErrorData` snapshots (input, items, responses, last agent, guardrail results) are attached to every `AgentsException` on both streamed and non-streamed paths (`src/agents/run.py:1505-1517`, `src/agents/run_internal/run_loop.py:1178-1186`), giving handlers and loggers full state for post-mortems.
- **Closed registries for handler keys.** `RunErrorHandlers` and the tool-error formatter `kind` are closed unions, keeping the handler surface small and typed (`src/agents/run_error_handlers.py:53-54`, `src/agents/run_config.py:73`).

## Notable Patterns

- **Dual-path parity for error handling.** Max-turns handling and `RunErrorDetails` attachment are duplicated across `run.py` and `run_internal/run_loop.py` streaming paths (`src/agents/run.py:1066-1081` vs `src/agents/run_internal/run_loop.py:889-912`), consistent with the repo's stated streaming/non-streaming alignment rule (`AGENTS.md`, Agents Core Runtime Guidelines).
- **Retryable as a tri-state.** `retryable: bool | None` where `None` means "the SDK cannot safely classify the error" (`src/agents/sandbox/errors.py:86-87`), and cause-chain propagation of the flag (`src/agents/sandbox/errors.py:97-99`) — an honest encoding of classification uncertainty.
- **Replay-safety as a distinct retry concern.** `ModelRetryAdvice.replay_safety` ("safe"/"unsafe") and `RetryDecision._approves_replay` separate "can we retry" from "is replaying this request side-effect-safe" (`src/agents/retry.py:92-100`, `:146-210`), with stateful requests (previous_response_id/conversation_id) treated specially (`src/agents/models/_openai_retry.py:121-122`, `:208-210`).
- **Error-kind-specific trace treatment.** `_should_attach_generic_agent_error` skips span-error attachment for `ModelBehaviorError` and tripwires because those already have dedicated span errors (`src/agents/run_internal/run_loop.py:261-265`).
- **Sensitive-data redaction in errors.** Tool error strings are redacted in traces unless `trace_include_sensitive_data` is set (`src/agents/util/_tool_errors.py:6-8`).

## Tradeoffs

- **Class hierarchy vs machine-readable codes.** Callers can `except ModelBehaviorError` but cannot switch on a stable code for core errors; only the sandbox subsystem offers that (`src/agents/sandbox/errors.py:12-55`). A user wanting programmatic classification of core errors must rely on class identity.
- **Raw provider exceptions leak through.** Not wrapping OpenAI SDK errors keeps fidelity but means callers catching `AgentsException` will not catch provider failures (`src/agents/run.py:1507` only decorates `AgentsException` with `run_data`), and provider error behavior varies by backend (e.g., `ResponsesWebSocketError` is a `RuntimeError`, `src/agents/models/openai_responses.py:314`).
- **Closed handler registry vs extensibility.** The two-key `RunErrorHandlers` TypedDict is simple and type-safe but every new error kind (e.g., tool guardrail tripwires, which are *not* handler-routable) requires SDK-side changes (`src/agents/run_error_handlers.py:50-54`).
- **String-based classification is cheap but brittle.** `terminal_metadata_for_exception` avoids an import cycle by matching `type(exc).__name__` strings (`src/agents/sandbox/memory/rollouts.py:183-188`), trading robustness for decoupling.
- **Recoverable-by-default tools hide failures.** `error_as_result` keeps runs alive but means tool crashes surface only as model-visible text unless the user opts into exceptions or inspects traces (`src/agents/tool.py:439-447`).

## Failure Modes / Edge Cases

- **Name-based misclassification.** Any exception class whose name contains "Guardrail" (even unrelated) is logged as `guardrail_tripped` in rollouts; a renamed or subclassed `MaxTurnsExceeded` would fall to generic `failed` (`src/agents/sandbox/memory/rollouts.py:183-190`).
- **Unhandled tool-guardrail tripwires.** `ToolInputGuardrailTripwireTriggered`/`ToolOutputGuardrailTripwireTriggered` (`src/agents/exceptions.py:147-174`) are not routable via `error_handlers`, so a tripwire mid-run always fails the run with no recovery hook (raise sites: `src/agents/run_internal/tool_execution.py:2364`, `:2401`).
- **Handler output validation converts failure mode.** A handler returning unserializable output raises `UserError`, masking the original `MaxTurnsExceeded`/`ModelRefusalError` (`src/agents/run_internal/error_handlers.py:100-107`); invalid handler result dicts also become `UserError` (`:160-164`).
- **Timeout races.** On timeout, if the tool task already completed, its real result/exception is returned instead of `ToolTimeoutError` (`src/agents/tool.py:1824-1828`) — correct, but means timeout classification is best-effort under races.
- **Retry classification uncertainty is explicit but caller-visible.** Sandbox errors constructed with `retryable=None` (e.g., `GitCloneError`, `src/agents/sandbox/errors.py:709`) push the retry decision to consumers, who must handle the tri-state.
- **Provider retry hints can conflict with policies.** Provider advice can hard-veto a user policy (`_with_hard_veto`, `src/agents/retry.py:201-203`, applied in `all`/`any` composition at `:302-358`), so a policy that "wants" a retry may be overridden — documented behavior, but a surprising failure mode for policy authors.

## Future Considerations

- Wrap or register provider-originated errors (e.g., make transport errors `AgentsException` subclasses or expose a `provider_error` marker) so a single `except AgentsException` catches the full SDK failure surface.
- Add a stable machine-readable code (sandbox-style `ErrorCode`) to core `AgentsException` classes, or a `retryable`/`escalate` hint, so retry-vs-stop decisions don't require per-type isinstance knowledge.
- Replace name-string matching in `terminal_metadata_for_exception` with isinstance checks or an explicit registration map (`src/agents/sandbox/memory/rollouts.py:175-196`).
- Extend `RunErrorHandlers` (or introduce a generic fallback handler key) to cover tool guardrail tripwires and other currently non-recoverable error kinds.
- Complete the docs Exceptions overview to include `ModelRefusalError`, `MCPToolCancellationError`, and the tool-guardrail tripwire types (`docs/running_agents.md:553-564`).

## Questions / Gaps

- No evidence found of a single unified error-taxonomy document or enum covering all subsystems (core, sandbox, voice, realtime); each subsystem defines its own. Searched `src/agents/**/errors.py`, `exceptions.py`, `docs/`, and `mkdocs.yml` nav.
- No evidence found of retry/escalation policy application to non-provider errors (e.g., `ModelBehaviorError` is never retried automatically); the retry machinery in `src/agents/retry.py` and `src/agents/run_internal/model_retry.py` only engages around model calls.
- No evidence found that the sandbox `ErrorCode` enum is versioned or documented for external stability beyond its docstring ("Stable, machine-readable error codes", `src/agents/sandbox/errors.py:13`); searched `docs/` for sandbox error documentation.
- The `voice` and `realtime` subsystems have minimal taxonomy (`STTWebsocketConnectionError` at `src/agents/voice/exceptions.py:4`; `RealtimeModelErrorEvent`/`RealtimeModelExceptionEvent` at `src/agents/realtime/model_events.py:12`, `:159`); whether realtime errors integrate with the retry taxonomy was not investigated in depth (out of scope for this dimension's evidence boundary).

---

Generated by `13.01-error-taxonomy` against `openai-agents-sdk`.
