# Source Analysis: openai-agents-sdk

## Dimension 01.16: Result and Outcome Publication Contract

### Source Info

| Field | Value |
|-------|-------|
| Name | openai-agents-sdk |
| Path | `studies/agent-harness-study/sources/openai-agents-sdk` |
| Language / Stack | Python / asyncio, Pydantic, OpenAI Responses API |
| Analyzed | 2026-09-01 |

## Summary

OpenAI Agents SDK exposes two terminal-outcome types (`RunResult`, `RunResultStreaming`) built atop plain `asyncio` coroutines, not a first-class `Future`/`Promise` combinator. The non-streaming path returns `RunResult` directly from `Runner.run()`; callers `await` the coroutine and receive an immutable-by-convention dataclass. The streaming path returns a `RunResultStreaming` object immediately and publishes outcome asynchronously via an `asyncio.Queue[StreamEvent|QueueCompleteSentinel]` plus a background `run_loop_task`. Final commitment is gated by queue-drain and task settlement (`stream_events()` joins `run_loop_task`, `_event_queue.join()`, `_stream_consumers_stopped`). Errors are surfaced two ways: raised synchronously from `Runner.run()` with an attached `RunErrorDetails` (`AgentsException.run_data`), or lazily from `stream_events()` via `_stored_exception` / `run_loop_task.exception()` after `QueueCompleteSentinel`. Value ownership is shallow-copy based (`copy_input_items`, `list()` snapshots) with no freeze/ownership transfer; `final_output` and `new_items` remain mutable aliases. No shared-future fanout, no weak-promise coalescing, and no cancellation token that pauses waiting without stopping the run itself.

## Rating

**5 / 10**

Strong on typed terminal-outcome separation and diagnostic preservation (`RunErrorDetails`, `run_data`, `run_loop_exception`), weak on publication contract rigor: no explicit immutable publication barrier, no copied/frozen outcome, no stable concurrent/late join semantics beyond ad-hoc `stream_events()` consumer counting, and no waiter-vs-run cancellation isolation (only run-level `cancel(immediate|after_turn)`).

## Evidence Collected

Every entry MUST include a file path with line numbers. Format: `path/to/file.ts:NN`.

| Area | Evidence | File:Line |
|------|----------|-----------|
| Work-function return types | `Agent` tool functions return `str|dict|object` via `FunctionTool`; tool failure mapped to `FunctionToolResult` text | `src/agents/tool.py:93-160`, `src/agents/run_internal/tool_execution.py:1-200` |
| Invocation API | `Runner.run(starting_agent, input, ...) -> RunResult`; `Runner.run_sync` wraps `run`; `Runner.run_streamed` -> `RunResultStreaming` | `src/agents/run.py:255-272`, `src/agents/run.py:362-376`, `src/agents/run.py:465-479` |
| Terminal outcome types | `RunResultBase` abstract base; `RunResult` (non-streaming) with `final_output: Any, _last_agent, raw_responses, new_items`; `RunResultStreaming` with `final_output: Any\|None` until complete, `is_complete: bool`, `raw_responses`, `new_items` | `src/agents/result.py:308-410`, `src/agents/result.py:483-590`, `src/agents/result.py:596-691` |
| Wait/join/future types | No `Future`/`Promise` handle; `RunResultStreaming.run_loop_task: asyncio.Task`, `stream_events() -> AsyncIterator[StreamEvent]`, `run_loop_exception` property | `src/agents/result.py:642-647`, `src/agents/result.py:883-916`, `src/agents/result.py:793-817` |
| Result+error combinations | `MaxTurnsExceeded`, `InputGuardrailTripwireTriggered`, `OutputGuardrailTripwireTriggered`, `ModelBehaviorError`, `ModelRefusalError` all carry optional `run_data: RunErrorDetails`; raised vs stored | `src/agents/exceptions.py:413-432`, `src/agents/exceptions.py:434-572`, `src/agents/result.py:1048-1091` |
| Outcome construction | `_finalize_result()` sets `_starting_agent_for_state`, `_trace_state`, `_generated_prompt_cache_key`, `sandbox_runtime.apply_result_metadata`; streaming `_finalize_streamed_final_output()` publishes `final_output`, pushes `QueueCompleteSentinel`, sets `is_complete` | `src/agents/run.py:828-851`, `src/agents/run_internal/run_loop.py:491-649` |
| Publication ordering | `stream_events()` guarantees queue drain before raising stored exception: checks `_stored_exception` + `should_drain_queued_events`, waits `event_queue.join()` + `stream_consumers_stopped`, `_drain_event_queue()` in `finally` | `src/agents/result.py:883-1026`, `src/agents/result.py:867-882`, `src/agents/result.py:1121-1139` |
| Concurrent/repeated/late access | Consumer registration with `active_stream_consumers` counter, `_stream_consumers_stopped: asyncio.Event`, `register_current_consumer`/`unregister_consumer`; second consumer shares same queue (first-come drains) | `src/agents/result.py:893-935`, `src/agents/result.py:654-657` |
| Value ownership / copying | `copy_input_items()` used for every `original_input` storage; `list(generated_items)`/`list(session_items)` snapshots; `to_state()` via `_populate_state_from_result` deep-copies pending fields; no freeze; `release_agents()` weakref dance | `src/agents/run_internal/items.py:1-60`, `src/agents/run.py:770-775`, `src/agents/result.py:108-173`, `src/agents/result.py:385-412` |
| Error cause/panic/stack preservation | `RunErrorDetails(input, new_items, raw_responses, last_agent, context_wrapper, guardrail_results)` built via `build_run_error_data`; `_create_error_details()`; `_is_error_data_redacted` redaction path detaches traceback to avoid leaking payload data | `src/agents/run_internal/error_handlers.py:113-136`, `src/agents/result.py:1027-1046`, `src/agents/exceptions.py:32-218` |
| Cancellation separation | `RunResultStreaming.cancel(mode="immediate"|"after_turn")`; `immediate` cancels tasks + `is_complete=True` + `QueueCompleteSentinel`; `after_turn` sets flag, streaming loop checks `_cancel_mode` and stops before next turn | `src/agents/result.py:819-865`, `src/agents/run_internal/run_loop.py:357-367`, `src/agents/run_internal/run_loop.py:1470-1476` |
| Diagnostic retention tests | `test_run_error_includes_data` asserts `exc.value.run_data is RunErrorDetails` with `raw_responses`, `new_items`; `test_streamed_run_error_includes_data` same via `stream_events()` | `tests/test_run_error_details.py:12-48` |
| Max-turns / handler exclusivity | `MaxTurnsExceeded` raised when `current_turn > max_turns`; `resolve_run_error_handler_result` can convert to `RunResult` with validated `final_output` instead of exception | `src/agents/run.py:1477-1575`, `src/agents/run_internal/run_loop.py:1559-1673`, `src/agents/run_internal/error_handlers.py:226-262` |
| Streaming cancellation tests | `test_streaming_guardrail_style_cancel_after_threshold` (immediate leaves `final_output is None`), `test_streaming_cancel_after_turn_allows_turn_completion` (final_output populated) | `tests/test_example_workflows.py:391-441` |
| Interruption terminal outcomes | `NextStepInterruption` stored as `interruptions: list[ToolApprovalItem]` on both `RunResult` and `RunResultStreaming`; `to_state()` preserves `interruption` via `RunState`; not an exception but terminal outcome | `src/agents/result.py:517-518`, `src/agents/result.py:651-652`, `src/agents/run_state.py:845-850` |

## Answers to Dimension Questions

### 1. Is a work return distinct from the runtime's terminal outcome?

**Yes, clearly distinct.** Worker-level returns (function tool return values, `custom_output_extractor` strings, `AgentOutputSchema.validate_json` dicts) are intermediates (`FunctionToolResult.output: str`, `ToolCallOutputItem`) collected in `SingleStepResult.new_step_items` and `turn_result.tool_output` (`src/agents/run_internal/turn_resolution.py:857-875`, `src/agents/run_internal/run_steps.py`). The runtime's terminal outcome is a separate framework type: `RunResult.final_output: Any` or `RunResultStreaming.final_output` plus `raw_responses`, `new_items`, `input_guardrail_results`, etc. (`src/agents/result.py:483-517`, `src/agents/result.py:596-620`). `check_for_final_output_from_tools` (`src/agents/run_internal/turn_resolution.py:753-781`) and `execute_final_output` (`src/agents/run_internal/turn_resolution.py:392-423`) explicitly convert tool outputs into the terminal `NextStepFinalOutput` only when `tool_use_behavior` demands it. Non-streaming callers never see raw tool returns directly; they `await Runner.run()` for the aggregated `RunResult` (`src/agents/run.py:255-343`). Evidence gap: no dedicated `WorkReturn<T>` generic; distinction is by convention, not by type-level separation.

### 2. Which result, failure, cancellation, and runtime-fault combinations are valid?

Valid combinations enumerated:

| Combination | Representation | Evidence |
|-------------|---------------|----------|
| Success | `RunResult(final_output, new_items, raw_responses)` with `interruptions=[]`, `is_complete=True` via queue sentinel | `src/agents/run.py:2007-2034`, `src/agents/result.py:622-623` |
| Max turns exhausted | `MaxTurnsExceeded(message).run_data: RunErrorDetails` raised, or handler-produced `RunResult(final_output=handler_output)` if `error_handlers["max_turns"]` returns `RunErrorHandlerResult` | `src/agents/run.py:1478-1575`, `src/agents/run_internal/error_handlers.py:226-262` |
| Input guardrail tripwire | `InputGuardrailTripwireTriggered(guardrail_result).run_data` raised immediately | `src/agents/exceptions.py:519-529`, `src/agents/run.py:995-1014` |
| Output guardrail tripwire | `OutputGuardrailTripwireTriggered(guardrail_result).run_data` raised after `run_output_guardrails` | `src/agents/exceptions.py:532-542`, `src/agents/run_internal/run_loop.py:509-575` |
| Interruption (HITL) | `RunResult(interruptions=[ToolApprovalItem], final_output=prior)` with `RunState._current_step=NextStepInterruption`; not exceptional | `src/agents/result.py:517`, `src/agents/run.py:2035-2099` |
| Model behavior / refusal | `ModelBehaviorError` / `ModelRefusalError` with optional handler conversion to `RunResult` | `src/agents/exceptions.py:454-475`, `src/agents/run_internal/turn_resolution.py:958-998` |
| Runtime fault (non-Agents exception) | propagated unchanged, `attach_generic_agent_error` annotates span, no `run_data` | `src/agents/run.py:2134-2157`, `src/agents/run_internal/run_loop.py:1984-1999` |
| Cancellation | `asyncio.CancelledError` (or `_RedactedExceptionCancellationError` when data-redacted); streaming `cancel()` produces `is_complete=True` with `final_output=None` (immediate) or current-turn final_output (after_turn) | `src/agents/exceptions.py:40-41`, `src/agents/result.py:819-865` |

No evidence of simultaneous value+error (no `Result[Value, Error]` dual). Framework uses exclusivity: success publishes `final_output`; failure raises with attached partial history, never both.

### 3. What happens when work returns both a value and an error?

**Not a supported dual-return.** Tools are Python callables returning `str`; exceptions become tool-error text via `failure_error_function`/`ToolErrorFormatter` or `normalize_shell_output` (`src/agents/run.py:52-58`, `src/agents/run_internal/tool_execution.py`). The runner never expects `(value, error)` tuples. If a function tool raises, `execute_function_tool_calls` captures `Exception` and produces a `ToolCallOutputItem` containing error text (`An error occurred...`), and the turn proceeds as normal tool results (`tests/test_example_workflows.py:938-1011` shows cancelled nested tool preserving sibling output). If `failure_error_function=None` is explicitly set, `CancelledError` propagates as cancellation instead of stringification (`tests/test_example_workflows.py:1138-1192`). Partial return (tool succeeded + guardrail tripwire in same turn) resolves by raising the tripwire; already-completed tool outputs are still appended to `new_items` and `run_data.new_items` (`src/agents/run_internal/run_loop.py:453-477`). No ambiguity: error always wins via exception path; value path only when no tripwire/cancellation.

### 4. Can a terminal outcome expose a partial or mutable value?

**Yes, partially mutable.** `RunResult.new_items: list[RunItem]` and `RunResultStreaming.new_items` are plain mutable lists holding mutable `RunItem` objects with `agent: Agent` references and `raw_item: dict|BaseModel` (`src/agents/items.py:500-600`, `src/agents/result.py:315-316`). Callers can mutate `result.new_items.append(...)` or `result.new_items[0].raw_item["content"]=...` without copy barrier. `final_output: Any` is stored by reference, not deep-copied or frozen (`src/agents/result.py:326-324`, `src/agents/result.py:616`). The framework defensively copies *internally* (`copy_input_items`, `list(generated_items)` at `src/agents/run.py:770-775`, `src/agents/result.py:123-124`) but publishes aliases without `copy.deepcopy` or `MappingProxyType` freezing. Mitigation: `RunState` snapshots deep-copy lists (`src/agents/result.py:135`, `src/agents/run_state.py:925-929`), and `RunResultBase.to_input_list()` re-materializes via `ItemHelpers.input_to_new_input_list` (`src/agents/result.py:432-454`). Partial values are possible transiently: during streaming `final_output is None` until terminal (`src/agents/result.py:616`), but external code observing `result.final_output` mid-stream sees `None`, not a placeholder. After completion, `final_output` is fully populated (`src/agents/run_internal/run_loop.py:644`). No evidence of defensive exposure tests for post-publication mutation.

### 5. Do concurrent, repeated, and late waiters observe the same committed facts?

**No strong guarantee; late/repeated semantics are fragile.**

- **Non-streaming**: no multi-waiter handle. `await Runner.run()` returns a fresh `RunResult` object; repeated reads of local variable return same reference (`result is result` identity). No shared future to join twice.
- **Streaming**: `stream_events()` is the only waiter mechanism. Implementation tracks concurrent consumers via `active_stream_consumers: int` and `stream_consumers_stopped: Event` (`src/agents/result.py:654-657`). Concurrent callers share a single `asyncio.Queue` (`_event_queue`). The loop yields each queue item once (`await _event_queue.get()` then `task_done()`) (`src/agents/result.py:958-982`), so two slow consumers race: first `get()` consumes event, second blocks until next event - not broadcast. This means concurrent waiters do **not** observe identical event sequences; events are load-balanced, not replicated. Test suite contains no test asserting two simultaneous `stream_events()` consumers see same events (`tests/test_example_workflows.py:391-441` only tests single consumer).
- **Repeated** (re-entering `stream_events()` after exhaustion): `is_complete` flag causes immediate `break` if queue empty (`src/agents/result.py:954-955`), so second iteration yields no events but does not replay. Not idempotent replay.
- **Late** (start consuming after loop finishes): if caller starts after `QueueCompleteSentinel` already enqueued and drained by `_drain_event_queue` (`src/agents/result.py:1012-1014`), sentinel is gone; late consumer sees empty queue and hangs waiting for `get()` unless `is_complete` was already True. `_waiting_on_event_queue` sentinel push in `cancel("immediate")` (`src/agents/result.py:857`) mitigates but not general. No `run_loop_exception` broadcast until `stream_events()` is entered and `_check_errors()` polled.

Commitment barrier exists: `_wait_for_turn_event_consumption` joins `_event_queue.join()` before starting next turn (`src/agents/result.py:867-882`), so turn N events are fully consumed before turn N+1 begins, but this only applies to currently-registered consumers.

### 6. Can waiting be cancelled without cancelling the run itself?

**No.** Cancellation is attached to the run, not the waiter. `RunResultStreaming.cancel(mode)` cancels the background `run_loop_task` and input-guardrail tasks (`src/agents/result.py:1092-1101`), marks `is_complete`, and enqueues `QueueCompleteSentinel`. There is no `wait_token` or `WaitHandle.cancel(wait=True, run=False)`. `stream_events()` `except asyncio.CancelledError` path calls `self.cancel()` then re-raises (`src/agents/result.py:960-963`), ensuring waiter cancellation kills the run. Timeout mitigation is at model-call level (`ModelSettings.timeout`, `ToolSettings.timeout` producing `ModelTimeoutError`/`ToolTimeoutError` (`src/agents/exceptions.py:477-516`)) but those are failures, not waiter-isolated cancellations. No evidence of `asyncio.shield` or `wait_for(..., shield_run=True)` pattern.

### 7. Who owns values and diagnostics after publication?

**Ownership model: snapshot + alias, with explicit weakref release for agents.**

- `RunResult.new_items`, `raw_responses`, `input_guardrail_results` are owned by `RunResult` instance but alias the internal working lists (`RunResult(new_items=session_items)` reuses `session_items` list directly `src/agents/run.py:1377-1380`; only streaming deep-copies via `list()` snapshots at publish points).
- `final_output` is owned by caller after return; SDK retains no reference besides `RunResult.final_output` attribute. No `copy.deepcopy` on publish; caller can mutate string/dict without affecting reruns except when handler validation re-parses (`src/agents/run_internal/error_handlers.py:170-205`).
- `RunErrorDetails` (`src/agents/exceptions.py:413-432`) holds snapshots `list(new_items)`, `list(raw_responses)` taken at failure boundary (`src/agents/run_internal/error_handlers.py:129-130`, `src/agents/result.py:1036-1045`), detached from subsequent mutation.
- Agent references: `RunResult._last_agent_ref: weakref.ReferenceType` plus `__del__` + `release_agents()` drop strong references (`src/agents/result.py:385-412`). After `release_agents()`, `last_agent` property may raise `AgentsException` if GC'd (`src/agents/result.py:522-532`).
- Diagnostics (`SpanError`, `Trace`, `Usage`) tied to `RunContextWrapper.usage` which remains mutable on `RunResult.context_wrapper` (`src/agents/run_context.py`). No transfer of ownership; caller can read but also mutate `context_wrapper.usage`.
- No evidence of arena/copy-on-write: caller-visible `result.to_input_list()` recomputes from `new_items` each call (`src/agents/result.py:432-454`), not memoizing an owned buffer.

### 8. Does waiter release occur only after the complete outcome is visible?

**Mostly yes, with one observable gap for early `InputGuardrailTripwireTriggered` in streaming.**

Non-streaming: trivially yes - `await runner._run_impl()` only returns after `return _finalize_result(result)` at `src/agents/run.py:1403`, `src/agents/run.py:2034`, guaranteeing all `new_items`, `raw_responses`, `input_guardrail_results` etc. are populated before caller unblocks. Exception path populates `exc.run_data` before `raise` (`src/agents/run.py:2146-2156`).

Streaming: commitment barrier is:

1. `_event_queue.put_nowait(QueueCompleteSentinel())` **after** setting `final_output`/`is_complete`/`output_guardrail_results` (`src/agents/run_internal/run_loop.py:644-649`).
2. `stream_events()` `finally` awaits `run_loop_task` via `_await_task_safely` and re-checks `_check_errors()` before draining (`src/agents/result.py:996-1001`), ensuring `stored_exception` (including `MaxTurnsExceeded`, `ModelBehaviorError`, guardrail trips) is visible as exception raise **after** queue exhaustion (`src/agents/result.py:1017-1025`).
3. `_wait_for_turn_event_consumption` (`src/agents/result.py:867-882`) blocks turn advancement until `queue.join()` confirms consumers processed events.

Visibility after error is enforced: `_check_errors()` populates `_stored_exception` from `run_loop_task.exception()`, `input_guardrails_task`, `max_turns` (`src/agents/result.py:1048-1091`), and `stream_events()` raises it only after draining sentinel or checking queue emptiness with guard `should_drain_queued_events` (`src/agents/result.py:941-952`). Tests confirm `async for event in result.stream_events(): pass` precedes `exc.value.run_data` inspection (`tests/test_run_error_details.py:40-48`).

Gap: For `InputGuardrailTripwireTriggered` that fires before `start_streaming` persists session, the non-streaming path raises synchronously with `persist_session_items_for_guardrail_trip` having already run (`src/agents/run.py:995-1014`); streaming path similarly enqueues sentinel then raises (`src/agents/run_internal/run_loop.py:1222-1240`). In both, diagnostics are visible. No gap observed for final `RunResult` waiter; gap only for potential late consumer that never calls `stream_events()` (polling `run_loop_exception` property is opt-in, `src/agents/result.py:793-817`).

## Architectural Decisions

- **Dataclass value-object outcomes over Future combinators** (`src/agents/result.py:308-596`). Chosen for ergonomic `await Runner.run()` and Pydantic compatibility (`__get_pydantic_core_schema__` at `src/agents/result.py:370-378`). Tradeoff: no built-in multi-join, timeout-composition, or `asyncio.Future` cancellation propagation.
- **Queue + sentinel streaming** (`src/agents/result.py:635-640`, `src/agents/run_internal/run_loop.py:868-887`). Models OpenAI `Response` event stream (`RawResponsesStreamEvent` at `src/agents/stream_events.py`). Sentinel pattern avoids `None` ambiguity but requires manual drain ordering.
- **Error attachment via `AgentsException.run_data`** (`src/agents/exceptions.py:434-442`) rather than tuple `Result[Outcome|Error]`. Keeps `try/except` idiomatic Python but forces callers to inspect `exc.run_data` for partial history; no typed `Result` envelope.
- **Weakref agent indirection** (`src/agents/result.py:484-532`, `src/agents/result.py:621-726`). Solved graph-retention leaks when `RunResult` outlives `Agent`; `release_agents()` eagerly drops ref.
- **Deep `RunState` checkpoint copy** (`src/agents/run_state.py:925-964`, `src/agents/result.py:108-173`) for pause/resume, but shallow publish for prompt path. Architectural split: durable cross-process boundary is strict, in-process outcome is permissive.
- **Redacted-error boundary** (`src/agents/exceptions.py:159-410`, `src/agents/run.py:341-359`). Prevents model/tool payload leaking via traceback/cause chain; sanitizes `__traceback__`, `__cause__`, `__context__`, args clearing before re-raise. Influences outcome publication: data-redacted failures still carry sanitized `run_data` where safe, otherwise minimal message.

## Notable Patterns

- **Sentinel-terminated async iterator** (`QueueCompleteSentinel` at `src/agents/run_internal/run_steps.py`, produced in `src/agents/run_internal/run_loop.py:354`, consumed at `src/agents/result.py:967-977`). Allows distinguishing successful stream end from cancellation.
- **Consumer-counted rendezvous** (`_active_stream_consumers`, `_stream_consumers_stopped` at `src/agents/result.py:654-657`, registration at `src/agents/result.py:893-933`). Lightweight alternative to `asyncio.Barrier`.
- **Policy-switch cancellation** — `immediate` vs `after_turn` (`src/agents/result.py:819-865`). Single flag `_cancel_mode` polled at turn boundaries (`src/agents/run_internal/run_loop.py:357-367`). No per-waiter cancel scope.
- **Error-last publication** (`_check_errors` at `src/agents/result.py:1048-1091` polled each loop iteration and after queue drain). Mirrors Go's `err` after `close(ch)` pattern.
- **Guardrail-results-as-published-sidecar** — `output_guardrail_results` accumulated in `_finalize_streamed_final_output` before raising tripwire, ensuring observer sees guardrail metadata even on failure (`src/agents/run_internal/run_loop.py:528-575`).
- **Handler-shaped recovery** (`RunErrorHandlers: dict[RunErrorHandlerKind, Callable]` at `src/agents/run_error_handlers.py`, resolved at `src/agents/run_internal/error_handlers.py:226-262`). Turns terminal faults (`max_turns`, `model_refusal`, `invalid_final_output`) into success outcomes when handler returns `RunErrorHandlerResult(final_output)`. Streaming variant mirrors via `_finalize_streamed_final_output` with `on_persisted_after_guardrails` callback.

## Tradeoffs

- **Simplicity vs. fan-out**: Returning `RunResult` directly avoids future-graph complexity but leaves concurrent waiter semantics undefined (queue load-balancing bug surface).
- **Mutability vs. copy cost**: Not deep-copying `final_output`/`new_items` keeps hot path allocation low (runs with thousands of `RunItem`); risk is caller accidentally mutating published outcome corrupts `to_state()` replay (`src/agents/result.py:542-589` clones only at checkpoint boundary).
- **Exception-carried diagnostics vs. envelope**: Attaching `run_data` to exception preserves idiomatic `pytest.raises` checks (`tests/test_run_error_details.py:23-27`) but forces `except BaseException` handling for cancellation redaction gymnastics (`src/agents/run.py:201-215`, `src/agents/exceptions.py:242-288`).
- **Queue-based stream vs. async generator directly yielding model events**: Queue decouples producer (run loop) and consumer (caller) enabling backpressure via `qsize` and `join`; cost is extra allocation and sentinel drain logic that must stay exception-safe, as seen in `_safe_redacted_persistence_error` exception-group reconstruction (`src/agents/run_internal/run_loop.py:652-741`).
- **Redaction vs. debuggability**: Stripping tracebacks on data-redacted errors (`_detach_data_redacted_error_traceback` at `src/agents/exceptions.py:149-156`) prevents payload exfiltration but makes production diagnosis harder; tracing span errors are kept generic (`GENERIC_AGENT_ERROR_MESSAGE` at `src/agents/run_internal/error_handlers.py:43`).

## Failure Modes / Edge Cases

- **Concurrent `stream_events()` duplication loss**: Two `async for` loops over same `RunResultStreaming` race on `queue.get()`; each event delivered to exactly one consumer, so aggregation tests expecting both see full history will flake. No mitigation (`src/agents/result.py:958-982`).
- **Late consumer deadlock**: If caller defers `stream_events()` until after `_drain_event_queue()` cleared sentinel (`src/agents/result.py:1012-1015`), `await _event_queue.get()` blocks forever unless `is_complete` already True. Workaround is to check `is_complete`/`run_loop_exception` first.
- **Cancelled waiter leak**: `CancelledError` during `stream_events()` cancels sibling tasks but `run_loop_task` traceback is detached if redacted, leaving original context dropped (`src/agents/result.py:960-963`, `src/agents/exceptions.py:281-288`).
- **Max-turns handler structured-output validation failure**: `validate_handler_final_output` raises `UserError` (`src/agents/run_internal/error_handlers.py:198-205`) which is then caught as generic `AgentsException` and published via `stored_exception`; caller sees handler fault, not original max-turns, which can mask limit breach.
- **Mutable `final_output` aliasing across nested agent tools**: `Agent.as_tool()` failure propagates as string output (`tests/test_example_workflows.py:990-1011`); if outer handler mutates dict `final_output` in place, inner `nested_state` checkpoint copy (`src/agents/run_state.py:925-964`) may observe mutated value due to shared `context` wrapper reference.
- **Session persistence failure masking terminal outcome**: `_safe_redacted_persistence_error` (`src/agents/run_internal/run_loop.py:676-741`) wraps persistence exception groups with sanitized message; if both persistence and guardrail fail, `redacted_persistence_error` takes precedence over guardrail-failure tripwire, surprising caller.
- **Interruptions swallowed on immediate cancel**: `cancel("immediate")` sets `is_complete=True` and drains queues (`src/agents/result.py:850-859`) before `_check_errors` can populate `interruptions`; approval items from current turn are lost.

## Future Considerations

- **Introduce `Future`-backed `RunHandle`** with `result()` / `exception()` / `add_done_callback()` and `wait(timeout=...)` separate from `run.cancel()`. Would give deterministic late-join (`handle.result()` returns cached `RunResult` after completion) and repeated read idempotency, without changing `Runner.run` await shorthand.
- **Freeze or copy-on-read `RunResult`**: `__post_init__` could `object.__setattr__(self, "new_items", tuple(mapped))` or expose `Sequence[RunItem]` view; `final_output` guarded by `copy.deepcopy` for dict/list types, or provide `result.copy()` helper. Cost minimal vs. correctness gain.
- **Broadcast queue for streaming**: replace `asyncio.Queue` with `broadcast` (e.g., `asyncio.Queue` per consumer or `pubsub` fan-out) so concurrent `stream_events()` consumers each observe identical `StreamEvent` sequence. Add `result.subscribe()` returning independent async iterator.
- **Waiter-isolated cancellation token**: `result.wait(timeout=...)` that raises `asyncio.TimeoutError` without cancelling `run_loop_task`, using `asyncio.wait_for` with `shield`.
- **Explicit `Outcome` envelope**: `ResultKind = Literal["success","interruption","failure"]` discriminator on `RunResult`, making `interruptions` vs `final_output` vs `run_data` exclusivity machine-checked and preventing partial-mutable confusion.
- **Stabilize publication barrier**: define `RunResult.is_published: bool` set only after `_finalize_result` copies and before return; for streaming, guarantee `run_loop_exception` and `is_complete` reflect outcome without requiring consumer to consume queue.

## Questions / Gaps

- **No evidence of late-waiter guarantee tests**: searched `tests/test_agent_runner_streamed.py`, `tests/test_example_workflows.py`, `tests/test_run_state.py` — none assert that a second `async for event in result.stream_events()` after first exhaustion replays same events, nor that concurrent consumers see identical sets. Gap acknowledged as design, not just missing tests.
- **Ownership/copy after publication**: Searched `src/agents/result.py`, `src/agents/run_internal/items.py`, `src/agents/util/_pretty_print.py` — no `__copy__`/`__deepcopy__` on `RunResult`, no frozen dataclass (`frozen=False` default for `RunResultBase` at `src/agents/result.py:308`). Hence mutation after publication remains possible; no defensive test found.
- **Result+error exclusivity contract not documented**: `pyproject.toml`, `README.md`, `docs/` describe guardrails and error handlers separately but do not state that `RunResult.final_output` and `raising MaxTurnsExceeded` are mutually exclusive, nor how `include_in_history=False` affects `new_items` visibility on handler-success. Inferred from `src/agents/run_internal/error_handlers.py:226-262` only.
- **Runtime fault vs. `BaseException` narrow**: `src/agents/run_internal/run_loop.py:1984-1999` re-raises `BaseException` directly without attaching `run_data`; interaction between `KeyboardInterrupt` arriving during streaming and `stream_events()` consumer's `CancelledError` handling has no dedicated test.

---

Generated by `Dimension 01.16: Result and Outcome Publication Contract` against `openai-agents-sdk`.
