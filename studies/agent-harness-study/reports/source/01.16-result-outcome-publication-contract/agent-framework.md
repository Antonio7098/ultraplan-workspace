# Source Analysis: agent-framework

## Dimension 01.16: Result and Outcome Publication Contract

### Source Info

| Field | Value |
|-------|-------|
| Name | agent-framework |
| Path | `studies/agent-harness-study/sources/agent-framework` |
| Language / Stack | Python (core `agent-framework-core`), dotnet, TypeScript; analysis focused on `python/packages/core` |
| Analyzed | 2026-09-01 |

## Summary

`agent-framework` implements two parallel runtime publication contracts: *Agent* (`AgentResponse`/`AgentResponseUpdate`/`ResponseStream`) and *Workflow* (`WorkflowRunResult`/`WorkflowEvent`/`WorkflowRunState`/`WorkflowErrorDetails`). Both are event-sourced: work functions (`Executor` handlers, `@workflow` steps, chat clients) produce data-plane emissions (`ctx.yield_output`, `ctx.send_message`, tool results) that the engine re-labels into caller-facing `output`/`intermediate` events via `OutputDesignation`, collects in a `WorkflowRunResult` list subclass, and seals with a terminal `WorkflowRunState` status event. Failure bypasses the regular result: executors emit `executor_failed`, the workflow drains pending events, emits `failed+FAILED`, updates `_status`, and re-raises the original exception — so no combined `value+error` object is ever returned. Ownership is reference-based with selective `deepcopy` at yield/invoke boundaries, not frozen; concurrent execution on the same instance is rejected (`_is_running`), and streaming uses a single-consumer `ResponseStream` whose `get_final_response()` finalizer runs only after the event queue is drained and cleanup hooks fire, providing visibility-before-release but no broadcast to multiple waiters.

## Rating

**5/10**

Rationale: Publication pipeline is explicit and well-typed (presence of `WorkflowRunResult`, `WorkflowErrorDetails`, `ResponseStream` finalizers, `OutputDesignation` taxonomies are concrete), trace/error diagnostics are retained, and the `IDLE/FAILED/IDLE_WITH_PENDING_REQUESTS` state machine is consistently exercised in tests. Downgraded for: (1) no value/error union type — partial outputs before a crash are observable only in the streaming event prefix, non-streaming callers receive only the exception; (2) no first-class cancellation (`CANCELLED` state exists but is never produced by the core workflow engine — only `WorkflowExecutor` mentions it); (3) ownership not hardened (`data` aliased, no freeze, mutable `WorkflowRunResult` is a mutable `list`); (4) no concurrent/late-waiter broadcast primitive — reuse is forbidden, not coordinated; (5) waiter cancellation is not separable from run cancellation.

## Evidence Collected

Every entry includes a file path with line numbers.

| Area | Evidence | File:Line |
|------|----------|-----------|
| Work-function return types | `Executor` handlers annotated `async def handle(self, message: str, ctx: WorkflowContext[int,str]) -> None` — work returns `None`, output via `ctx.send_message`/`ctx.yield_output` not return value | `python/packages/core/agent_framework/_workflows/_executor.py:38-65` |
| Work-function return types | `@handler` decorator validates `message_param` + `WorkflowContext[OutT,W_OutT]` and populates `output_types`/`workflow_output_types` | `python/packages/core/agent_framework/_workflows/_executor.py:694-761` |
| Work-function return types | Functional `@workflow` handler returns value → workflow engine re-emits as `WorkflowEvent("output", ...)` only if non-`None` | `python/packages/core/agent_framework/_workflows/_functional.py:982-987` |
| Workflow invocation handle types | `Workflow.run(message, stream=False)` → `Awaitable[WorkflowRunResult]`; `stream=True` → `ResponseStream[WorkflowEvent, WorkflowRunResult]`; mutual-exclusion validation for `message/responses/checkpoint_id` | `python/packages/core/agent_framework/_workflows/_workflow.py:674-771` |
| Workflow invocation handle types | Functional workflow overload `run(...) -> WorkflowRunResult | ResponseStream` with same `message: Any | None, *, stream: bool` shape | `python/packages/core/agent_framework/_workflows/_functional.py:770-840` |
| Workflow invocation handle types | `WorkflowAgent.run(messages, stream=False)->AgentResponse` vs `stream=True->ResponseStream[AgentResponseUpdate,AgentResponse]` wrapping workflow events | `python/packages/core/agent_framework/_workflows/_agent.py:135-226` |
| Agent run handle APIs | `BaseAgent / RawAgent / Agent.run(... stream:bool)` → `Awaitable[AgentResponse] | ResponseStream[AgentResponseUpdate,AgentResponse]` via `ResponseStream.from_awaitable` | `python/packages/core/agent_framework/_agents.py:854-977` |
| Terminal outcome type | `WorkflowRunResult(list[WorkflowEvent])` with `status_timeline()`, `get_outputs()`, `get_intermediate_outputs()`, `get_final_state()->WorkflowRunState` | `python/packages/core/agent_framework/_workflows/_workflow.py:101-165` |
| Terminal outcome type | `AgentResponse` (non-streaming) with `messages: list[Message], response_id, usage_details, value (structured), raw_representation` and `AgentResponseUpdate` streaming chunk | `python/packages/core/agent_framework/_types.py:2536-2833` |
| Terminal outcome type | `WorkflowRunState` enum: `STARTED, IN_PROGRESS, IN_PROGRESS_PENDING_REQUESTS, IDLE, IDLE_WITH_PENDING_REQUESTS, FAILED, CANCELLED` | `python/packages/core/agent_framework/_workflows/_events.py:58-68` |
| Wait/join/future/promise API | `ResponseStream` generic `ResponseStream[UpdateT,FinalT]` with `__aiter__`, `__anext__`, `get_final_response()->FinalT`, `map`, `with_finalizer`, `with_transform_hook`, `with_cleanup_hook` | `python/packages/core/agent_framework/_types.py:2939-3331` |
| Wait/join/future/promise API | Workflow `.status: WorkflowRunState` mirrors last status event, updated in lockstep in `_run_workflow_with_tracing` | `python/packages/core/agent_framework/_workflows/_workflow.py:368-378,526-603` |
| Wait/join run-handle | `AgentExecutorResponse` used inside `WorkflowExecutor` for concurrent sub-workflow fan-out | `python/packages/orchestrations/agent_framework_orchestrations/_concurrent.py:9-95` (not core) |
| Result+error resolution | `_run_workflow_with_tracing` drains pending events, creates `WorkflowErrorDetails.from_exception(exc)`, yields `failed` + `FAILED` status, then `raise` (no `WorkflowRunResult` returned on failure) | `python/packages/core/agent_framework/_workflows/_workflow.py:606-628` |
| Result+error resolution | Functional workflow same pattern: `except Exception: yield WorkflowEvent.failed(details); yield status(FAILED); raise` | `python/packages/core/agent_framework/_workflows/_functional.py:1051-1064` |
| Result+error resolution | Executor-level failure emits `executor_failed` before propagate, allowing runner to drain it before `failed` | `python/packages/core/agent_framework/_workflows/_executor.py:280-286` |
| Result+error resolution | Executor failure drained by runner before raising: `if has_events: drain_events` inside `except` | `python/packages/core/agent_framework/_workflows/_runner.py:123-130` |
| Outcome immutability gate | `State.commit()` at superstep boundary only; per-run iteration count reset only for `not is_continuation`; checkpoints commit before creation | `python/packages/core/agent_framework/_workflows/_runner.py:141-144,221-224` ; `python/packages/core/agent_framework/_workflows/_workflow.py:570-571` |
| Outcome immutability gate | Workflow final status decided after convergence: `saw_request` → `IDLE_WITH_PENDING_REQUESTS` else `IDLE` | `python/packages/core/agent_framework/_workflows/_workflow.py:594-603` |
| Publication code | `ResponseStream._finalizer` (`_finalize_events` → `WorkflowRunResult(filtered, status_events)`) filters `started`/`status` unless `include_status_events` | `python/packages/core/agent_framework/_workflows/_workflow.py:838-863` |
| Publication code | `ResponseStream.get_final_response()` consumes stream if `_consumed==False`, runs `cleanup_hooks`, then `finalizer`, caching `_final_result`/`_finalized` | `python/packages/core/agent_framework/_types.py:3165-3271` |
| Publication code | `Workflow._run_cleanup` clears runtime checkpoint storage and resets `_is_running` flag (called via `cleanup_hooks`) | `python/packages/core/agent_framework/_workflows/_workflow.py:832-837` |
| Copying/freezing/aliasing | `Executor.execute` captures `copy.deepcopy(message)` for `executor_invoked` event | `python/packages/core/agent_framework/_workflows/_executor.py:274-277` |
| Copying/freezing/aliasing | `WorkflowContext.yield_output` → `self._yielded_outputs.append(copy.deepcopy(output))` but `WorkflowEvent(..., data=output)` keeps alias | `python/packages/core/agent_framework/_workflows/_workflow_context.py:360-370` |
| Copying/freezing/aliasing | `InProcRunnerContext.drain_messages` → `copy(self._messages)`; `State.export_state` used for checkpoint deep serialization | `python/packages/core/agent_framework/_workflows/_runner_context.py:307-310` |
| Copying/freezing/aliasing | `CheckpointDatabase.save` → `copy.deepcopy(checkpoint)` | `python/packages/core/agent_framework/_workflows/_checkpoint.py:201` |
| Copying/freezing/aliasing | `Content.__deepcopy__` preserves `raw_representation` by reference, deep-copies rest | `python/packages/core/agent_framework/_types.py:576-591` |
| Concurrent/repeated/late access | `Workflow._ensure_not_running()` → `raise RuntimeError("Workflow is already running. Concurrent executions are not allowed.")`; functional counterpart same | `python/packages/core/agent_framework/_workflows/_workflow.py:379-383` ; `python/packages/core/agent_framework/_workflows/_functional.py:1229-1231` |
| Concurrent/repeated/late access | `Runner.run_until_convergence` `if self._running: raise WorkflowRunnerException` | `python/packages/core/agent_framework/_workflows/_runner.py:80-82` |
| Concurrent/repeated/late access | `ResponseStream` single-consumer: `_consumed`/`_finalized` flags, `_updates` list; repeated `get_final_response()` returns cached `_final_result` without re-consuming | `python/packages/core/agent_framework/_types.py:2967-2971,3248-3271` |
| Cancel vs wait cancel | `Runner.run_until_convergence` `except CancelledError: iteration_task.cancel(); suppress; raise` propagates caller cancellation to iteration task | `python/packages/core/agent_framework/_workflows/_runner.py:115-120` |
| Cancel vs wait cancel | `WorkflowExecutor` handles `CANCELLED` state but core engine never sets `CANCELLED` — only `FAILED`/`IDLE*` | `python/packages/core/agent_framework/_workflows/_workflow_executor.py:627-630` |
| Diagnostic retention | `WorkflowErrorDetails.from_exception` captures `error_type`, `message`, `traceback` via `traceback.format_exception`, optional `executor_id`, `extra` | `python/packages/core/agent_framework/_workflows/_events.py:71-99` |
| Diagnostic retention | `Message`/`Content` `raw_representation` kept shallow in deepcopy by `Content._SHALLOW_COPY_FIELDS` | `python/packages/core/agent_framework/_types.py:536-591` |
| Diagnostic retention | `WorkflowEvent.request_info` stores `request_type`/`response_type` via `serialize_type` for later `ValueError` on mismatch in `_send_responses_internal` | `python/packages/core/agent_framework/_workflows/_events.py:296-313,412-448` ; `python/packages/core/agent_framework/_workflows/_workflow.py:936-960` |
| Concurrency tests | `test_functional_workflow.py::test_step_failure_emits_executor_failed` asserts `executor_failed` present before exception propagation | `python/packages/core/tests/workflow/test_functional_workflow.py:320-335` |
| Stable access tests | `test_workflow_failure_emits_failed_status` asserts `failed` event + `FAILED` status before `RuntimeError` re-raise | `python/packages/core/tests/workflow/test_functional_workflow.py:336-351` |
| Publication ordering | `test_full_conversation.py` / `test_strict_mode_event_labeling.py` exercise `get_outputs()`, `get_intermediate_outputs()`, `status_timeline()` | `python/packages/core/tests/workflow/test_strict_mode_event_labeling.py:92-113` etc. |
| Outcome not tested: late waiter | No evidence of tests for two concurrent `await workflow.run(...)` or two concurrent `await stream.get_final_response()` on same stream | Search `tests/workflow` for `asyncio.gather.*run` — No evidence found |

## Answers to Dimension Questions

### 1. Is a work return distinct from the runtime's terminal outcome?

**Yes, strictly distinct.** Every workload handler returns `None` (`Executor` handlers are `-> None` by contract `python/packages/core/agent_framework/_workflows/_executor.py:219-295`). Output is via side-effect methods: `ctx.send_message()` for inter-executor messages, `ctx.yield_output()` for workflow-level outputs (`python/packages/core/agent_framework/_workflows/_workflow_context.py:308-371`), while `AgentResponse` is produced downstream from chat client results. The *terminal outcome* is a different type family: `WorkflowRunResult` (a `list[WorkflowEvent]` with `status_timeline()`/`get_final_state()` `python/packages/core/agent_framework/_workflows/_workflow.py:101-165`) for workflows, or `AgentResponse`/`ChatResponse` (`python/packages/core/agent_framework/_types.py:2536-2644,1979-2200`) for agents, sealed with a `WorkflowRunState` (`python/packages/core/agent_framework/_workflows/_events.py:58-68`). `workflow.status` mirrors the last status event to expose the terminal fact outside the event payload (`python/packages/core/agent_framework/_workflows/_workflow.py:368-378`). Functional workflows blur the line slightly by returning a value from the user function, but `_run_core` immediately converts it to an `output` event (`python/packages/core/agent_framework/_workflows/_functional.py:985-987`), so the consumer still observes an event, not a raw return.

### 2. Which result, failure, cancellation, and runtime-fault combinations are valid?

* **Success (IDLE):** `WorkflowRunResult` with zero or more `output`/`intermediate` events plus terminal `status(IDLE)` (`python/packages/core/agent_framework/_workflows/_workflow.py:599-603`). Hitless normal path.
* **Success-with-HITL (IDLE_WITH_PENDING_REQUESTS):** same `WorkflowRunResult` but at least one `request_info` event observed, final state `IDLE_WITH_PENDING_REQUESTS` (`python/packages/core/agent_framework/_workflows/_workflow.py:594-598`, `python/packages/core/agent_framework/_workflows/_functional.py:1012-1015`). Caller resumes via `run(responses=...)`.
* **Executor failure (FAILED):** Per-executor `executor_failed` event with `WorkflowErrorDetails` (`python/packages/core/agent_framework/_workflows/_executor.py:281-284`), followed by workflow-level `failed` event + `status(FAILED)` (`python/packages/core/agent_framework/_workflows/_workflow.py:612-619`). Runner convergence failure also produces `WorkflowConvergenceException` (`python/packages/core/agent_framework/_workflows/_runner.py:152-153`).
* **Runtime faults (checkpoint/storage):** `WorkflowCheckpointException` variants for missing/incompatible checkpoints (`python/packages/core/agent_framework/_workflows/_runner.py:270-300`, `python/packages/exceptions.py:257-260`).
* **Cancellation (CANCELLED):** Symbol exists (`_events.py:67`) but no core engine path sets it; only `WorkflowExecutor` switches on `CANCELLED` (`python/packages/core/agent_framework/_workflows/_workflow_executor.py:627-630`). The real cancellation primitive is `asyncio.CancelledError` propagated from the awaiter through `Runner.run_until_convergence` (`python/packages/core/agent_framework/_workflows/_runner.py:115-120`), which cancels the in-flight iteration task. So the framework distinguishes language-level cancellation from a committed `FAILED` or `CANCELLED` terminal state, but the latter is effectively unused.
* **Invalid combos:** `value+signed failed status` is not returned as a single object. On success you get `WorkflowRunResult`; on failure the call *raises* (`raise` after yielding `failed` in `_run_workflow_with_tracing`), never returns a result object wrapping an error.

### 3. What happens when work returns both a value and an error?

No union type exists. The engine obeys Python exception semantics: whichever path wins, the other is discarded as a return value. If handlers emitted outputs before raising, those outputs are still queued in `InProcRunnerContext._event_queue` and are drained and yielded *before* the terminal failure (`python/packages/core/agent_framework/_workflows/_workflow.py:607-610` drains `context.drain_events()`; `python/packages/core/agent_framework/_workflows/_runner.py:126-130` does the same). Streaming consumers will observe `[output[], executor_failed, failed, FAILED]` before the exception propagates (verified by `test_step_failure_emits_executor_failed` and `test_workflow_failure_emits_failed_status` at `tests/workflow/test_functional_workflow.py:320-351`). Non-streaming callers (`await workflow.run(...)`) receive only the raised exception, not a `WorkflowRunResult` containing the prefix outputs — `_run_core` is inside a `ResponseStream` finalizer that never calls `_finalize_events` when an exception is raised, because `get_final_response` is skipped in favor of exception propagation (`python/packages/core/agent_framework/_types.py:3104-3140` records `_stream_error`). So "value+error" is observable only transiently on the stream, not as a persisted outcome.

### 4. Can a terminal outcome expose a partial or mutable value?

**Partial: partially yes.** Because failures still publish prefix outputs on the streaming prefix, late or replay consumers that capture events before the exception could interpret partial results as terminal. Non-streaming path hides this. Intermediate vs terminal labeling (`OutputDesignation` `python/packages/core/agent_framework/_workflows/_workflow.py:170-204`) can also hide payloads (`classify` returning `None` drops the event entirely in `WorkflowContext.yield_output` `python/packages/core/agent_framework/_workflows/_workflow_context.py:364-366`), so not all work is exposed terminally.

**Mutable: yes.** Publication does *not* freeze:

- `yield_output` deepcopies into `_yielded_outputs` (`python/packages/core/agent_framework/_workflows/_workflow_context.py:362`) but the emitted `WorkflowEvent.data` retains the original reference (`WorkflowEvent("output", executor_id, data=output)` line 369). The `WorkflowRunResult` stores those event objects directly (`python/packages/core/agent_framework/_workflows/_workflow.py:122-123,848-862`), and `get_outputs()` returns `event.data` references (`line 131`). A caller mutating a mutable output (e.g., `dict`, `list`, `Message.contents`) mutates the published outcome seen by any late reader holding the same `WorkflowRunResult` object.
- For agents, `AgentResponse.messages` is a plain mutable list (`python/packages/core/agent_framework/_types.py:2617-2630`).
- Checkpoint serialization deepcopies only when persisting via storage (`python/packages/core/agent_framework/_workflows/_checkpoint.py:201`), not for in-memory publication.
- `Content.__deepcopy__` intentionally keeps `raw_representation` by reference (`python/packages/core/agent_framework/_types.py:576-591`).

Thus after `get_final_response()` returns, the outcome is reference-shared and mutably observable; no `FrozenList`, `MappingProxy`, or `deepcopy` on `WorkflowRunResult` construction.

### 5. Do concurrent, repeated, and late waiters observe the same committed facts?

* **Concurrent:** No — concurrent calls are rejected, not coordinated. `Workflow._ensure_not_running` (`python/packages/core/agent_framework/_workflows/_workflow.py:379-383`), functional `Workflow._is_running` guard (`_functional.py:1229-1231`), and `Runner._running` (`_runner.py:80-82`) each raise `RuntimeError`/`WorkflowRunnerException` immediately. No queue, no join, no shared future.
* **Repeated:** Non-streaming `await workflow.run(...)` returns a fresh `WorkflowRunResult` instance per invocation; the workflow's `State` and `_status` persist across runs (iteration count is reset but `State` is *not* cleared except at restore, `python/packages/core/agent_framework/_workflows/_workflow.py:545-547`). Calling `get_final_state()` on an earlier result remains stable; calling `workflow.status` reflects the newest terminal state, so older result objects don't track later runs — they are snapshots.
* **Late:** `WorkflowRunResult` is a synchronous list, so any holder can call `get_outputs()`/`status_timeline()` idempotently; `ResponseStream.get_final_response()` after consumption returns the cached `_final_result` without re-running (`python/packages/core/agent_framework/_types.py:3248-3271`). But "late waiter joins after completion without polling" (a future-like `wait()`/`join()`) does not exist: there is no `workflow.wait()` API; you must retain the previously returned result or replay via checkpoint. No evidence of a `Future`/`Promise` handle.
* **Streaming late observer:** If two coroutines iterate the same `ResponseStream`, they share a single `_iterator` (`python/packages/core/agent_framework/_types.py:3113-3115`); after `_consumed=True`, further `__anext__` immediately runs finalization and raises `StopAsyncIteration`. So late iteration yields no events but can still `await stream.get_final_response()` to retrieve the cached final value.

### 6. Can waiting be cancelled without cancelling the run itself?

**No.** Cancellation is coupled. `Runner.run_until_convergence` explicitly catches `asyncio.CancelledError` on the polling loop, cancels the iteration task, suppresses its completion, and re-raises (`python/packages/core/agent_framework/_workflows/_runner.py:115-120`). The caller task that `await`s `workflow.run()` or iterates `run(stream=True)` is indistinguishable from the executor task — cancelling the awaiter cancels the run. No `cancel_wait=False` flag, no separate `CancellationToken` parameter, no `asyncio.shield`. MCP tool layer (`python/packages/core/agent_framework/_mcp.py:275-288`) does expose `cancel_remote_task_on_local_cancellation`, but that still treats local cancellation as run cancellation. There is no API to time out a waiter while leaving the workflow running in the background (background running via `BackgroundAgentsProvider` is a different abstraction that yields `TaskInfo` polling, not waiter cancellation).

### 7. Who owns values and diagnostics after publication?

* **Values:** Ownership transfers to the consumer by reference, with caller responsible for defensive copying. As above, `WorkflowRunResult` holds aliases; `WorkflowContext.yield_output` deepcopies only for executor-bookkeeping (`_yielded_outputs`) but not for the published event. Checkpoint persistence performs a defensive `deepcopy` on the stored copy (`_checkpoint.py:201`), but the in-memory published copy is shared. No ownership transfer via move semantics or `copy.deepcopy` at publication time, except for `Message` invoke capture (`copy.deepcopy(message)` in `_executor.py:276`) and diagnostics buckets (`_yielded_outputs` copy).
* **Diagnostics:** `WorkflowErrorDetails` (`python/packages/core/agent_framework/_workflows/_events.py:71-99`) owns a `traceback: str|None` formed by `format_exception` at `from_exception`. Callers receive `failed` and `executor_failed` events carrying `details` (`WorkflowEvent.failed(details)` line 266, `executor_failed` line 335). The original `Exception` instance is also re-raised, preserving the Python cause chain (`raise` without `from None` in `python/packages/core/agent_framework/_workflows/_workflow.py:628`), so callers can inspect `__cause__`/`__context__`. However `WorkflowErrorDetails` stores only the formatted string, not the `BaseException` object itself (nor `exc.__traceback__` as an object), so sideband causality must be recovered from the raised exception, not from the stored `failed` event alone. `raw_representation` fields on `Message`/`Content`/`AgentResponse` are shallow-retained as additional diagnostics (`python/packages/core/agent_framework/_types.py:576-591,2230-2243`).
* **Ownership boundary for state:** `State.export_state()` is used as checkpoint payload but caller `WorkflowRunResult` does not own `State` — state survives across `run()` calls in the `Workflow` instance (`python/packages/core/agent_framework/_workflows/_workflow.py:349-357`).

### 8. Does waiter release occur only after the complete outcome is visible?

**Mostly, with a caveat for failures.** For success:

- Non-streaming path: `ResponseStream.get_final_response()` drives `async for _ in self:` until `StopAsyncIteration`, at which point `_run_cleanup` (resetting `_is_running`, clearing runtime storage) runs *before* `finalizer` (`python/packages/core/agent_framework/_types.py:3104-3141` in `__anext__` cleanup, and `3165-3271` for finalizer ordering). The caller `await workflow.run(...)` (`python/packages/core/agent_framework/_workflows/_workflow.py:770` delegates to `response_stream.get_final_response()`) resumes only after `WorkflowRunResult` is materialized. Visibility precedes release.
- Streaming path: `Runner.run_until_convergence` yields `superstep_started`/`completed` delimiters around each iteration and drains events (`_runner.py:101,146,135-137`), but `_run_workflow_with_tracing` yields terminal `status(IDLE/IDLE_WITH_PENDING_REQUESTS)` *before* `span.add_event(WORKFLOW_COMPLETED)` (`_workflow.py:594-605`). The cleanup hook (`_run_cleanup`) is not run until the stream is fully consumed (`python/packages/core/agent_framework/_types.py:3118-3125`), so a consumer that stops iterating early (break without draining) would leave `_is_running=True` and miss the tail. Well-behaved consumers that exhaust the stream see complete visibility before `_is_running` reset.

**Failure caveat:** On failure, `_run_workflow_with_tracing` yields `failed` then `FAILED` status (`python/packages/core/agent_framework/_workflows/_workflow.py:613-619`) and re-raises. `ResponseStream.__anext__` captures `_stream_error`, runs cleanup, and re-raises (`_types.py:3122-3128`), so a consumer observing via `async for` sees the failed events *before* the exception propagates. A non-streaming consumer sees only the exception (no `WorkflowRunResult`), but failure events were already visible to the tracing/spans instrumentation. So waiter release (exception) occurs *after* failure diagnostics are enqueued, satisfying "complete outcome visible" for observability consumers, but not as a returned value for direct awaiters.

## Architectural Decisions

| Decision | Evidence | Effect |
|----------|----------|--------|
| Separate work side-effects from outcome types | Handlers `-> None` plus `WorkflowContext.yield_output/send_message` (`_executor.py:38-65`, `_workflow_context.py:308-371`) | Makes publication contract explicit; engine controls labeling via `OutputDesignation` |
| Reify terminal outcome as `WorkflowRunResult` list+timeline plus `WorkflowRunState` status events | `WorkflowRunResult` definition `_workflow.py:101-165`, `WorkflowRunState` `_events.py:58-68` | Idempotent snapshot for late readers; separates data-plane `output` from control-plane `status/failed` |
| Event identity via event queue, not Future/Promise | `InProcRunnerContext._event_queue = asyncio.Queue` `_runner_context.py:289`, `has_events/drain_events/next_event` `_runner_context.py:315-341` | Enables streaming interleaving with convergence loop; requires single-consumer discipline |
| Generic `ResponseStream[UpdateT,FinalT]` with finalizer/result-hooks/cleanup-hooks | `_types.py:2939-3331` | Uniform publication for Agents, Workflows, Chat clients; guarantees `finalizer()` ordering after drain |
| Defensive `OutputDesignation` at build time (allow/hide by executor ID) | `Workflow._output_designation` `_workflow.py:336-348`, `OutputDesignation.classify` `171-203` | Lets orchestration filter noise; prevents leaks of internal yields |
| Eager `_is_running` flag instead of lock/queue | `_workflow.py:379-383`; `_runner.py:80-82` | Fail-fast on concurrent reuse; avoids torn outcomes |
| `capture_exception` + `WorkflowErrorDetails` string trace | `_agents.py: ` via `WorkflowErrorDetails.from_exception` `_events.py:81-99` + `raise` `_workflow.py:628` | Preserves human-readable diagnostics alongside machine `error_type` without retaining frame objects |

## Notable Patterns

- **Single-consumer stream with cached final value** — `ResponseStream` holds `_updates`, `_final_result`, `_finalized` (`_types.py:2967-2981`) enabling `await stream.get_final_response()` after iterator exhaustion to return the cached outcome. Mirrors the classic "channel then promise" but without multi-waiter broadcast.
- **Poll-and-yield convergence loop** — `Runner.run_until_convergence` drives `_run_iteration()` as an `asyncio.create_task` and polls `await wait_for(next_event, timeout=0.05)` (`_runner.py:105-114`), giving near-real-time streaming while preserving ordered delivery per edge runner.
- **Framework- vs executor-originated events** — `ContextVar[WorkflowEventSource]` `_event_origin_context` `_events.py:38-56` and `_framework_event_origin()` guard prevent executors from forging `started/status/failed/output` events; violations become `warning` (`_workflow_context.py:373-391`).
- **Request/response HITL as first-class event** — `ctx.request_info()` posts `request_info` event correlated by `request_id` (`_workflow_context.py:393-424`), `InProcRunnerContext.add_request_info_event`/`send_request_info_response` maintain `_pending_request_info_events` map (`_runner_context.py:446-494`), resumed via `run(responses=...)` or `WorkflowAgent._extract_function_responses` (`_agent.py:716-737`).

## Tradeoffs

| Tradeoff | Pro | Con | File:Line |
|----------|-----|-----|------------|
| Exclusive execution via exception vs queued/joined | Simple invariant, no torn intermediate state | Legitimate concurrent use cases must build their own orchestration; no waiters coordination | `_workflow.py:381-383` |
| Stream-then-finalize wrapper (`ResponseStream`) vs `asyncio.Future` handle | Works for both streaming & non-streaming call sites, carries transform/result hooks | Only single consumer, not thread-safe, late reader must reuse same object — no broadcast `Future` | `_types.py:2939-2980` |
| Return `WorkflowRunResult` on success, `raise` on failure | Idiomatic Python, callers can distinguish via `try/except` | No unified discriminated `Result|Error` value; partial outputs lost to non-streaming callers | `_workflow.py:606-628` |
| Reference-share published values (select `deepcopy`) vs deep freeze/clone | Zero-copy fast path, preserves non-serializable `raw_representation` | Caller can mutate published outcome; not safe for sharing across waiters/threads | `_workflow_context.py:362-369` |
| Status event stream vs single status field | Allows timeline reasoning, works with streaming inspection | Callers must scan timeline or call `get_final_state()` which `raise RuntimeError` if none emitted (`_workflow.py:152-160`) | `_workflow.py:151-161` |
| `CANCELLED` state reserved but not set by engine | Leaves space for external orchestrator cancellation | Confusing for users who observe `WorkflowRunState.CANCELLED` in public enum | `_events.py:67`, `_workflow_executor.py:627` |

## Failure Modes / Edge Cases

| Scenario | Behavior | Visibility | Risk |
|----------|----------|------------|------|
| Second `run()` while _is_running | `RuntimeError("Workflow is already running")` immediately, no queuing | Exception synchronous at entry | Misuse surfaces early; not retriable as a join |
| `await workflow.run()` cancelled mid-superstep | `CancelledError` propagates to `iteration_task` which is `cancel()`ed and suppressed (`_runner.py:115-120`), then re-raised to awaiter; no `FAILED`/`CANCELLED` status emitted | Caller gets `CancelledError`; observers on streaming prefix see events up to cancel point then `CancelledError` | Run left in indeterminate `IN_PROGRESS` if not retried via checkpoint; state not rolled back |
| Executor throws inside handler | `executor_failed` emitted with `WorkflowErrorDetails`, runner drains it, workflow emits `failed+FAILED`, then `raise original exc`; non-streaming caller sees only exception | `executor_failed.details.traceback` string retained; status timeline shows `FAILED` but only to streaming consumers that drained before exception | Diagnostics split: string traceback in one place, exception object in another |
| `max_iterations` exceeded with messages remaining | `WorkflowConvergenceException` raised (`_runner.py:152-153`), outer handler treats as failure → `failed/FAILED` as above | Exception message includes count; no partial result returned | Debugging non-convergence requires inspecting trace logs, not result payload |
| Checkpoint `graph_signature_hash` mismatch on `apply_checkpoint` | `WorkflowCheckpointException` raise (`_runner.py:275-279`) | Message names drift | Resume impossible without rebuilding original graph |
| `ctx.yield_output` with hidden executor | Silently dropped (`classify` returns `None`, early `return` `_workflow_context.py:364-366`) | No warning, no error | Silent data loss; caller may wrongly believe yield succeeded |
| `ctx.yield_output` deepcopies but event carries alias | Post-yield mutation of `output` may alter event payload but not `_yielded_outputs` bookkeeping | Inconsistent view between executor_completed bookkeeping and caller-visible result | Hard-to-reproduce divergence if handlers reuse mutable buffers |
| `missing status event → get_final_state()` | `RuntimeError("Final state is unknown because no status event was emitted...")` (`_workflow.py:152-160`) | Guard | Indicates non-standard entry point that bypassed `_run_workflow_with_tracing` |
| `responses` for unknown/non-pending `request_id` | `ValueError(f"Response provided for unknown request ID")` in `_send_responses_internal` (`_workflow.py:946`) | Early | Fails fast; types narrowed via `is_instance_of` checks |
| Streaming consumer short-circuit (break without drain) | `_is_running` never cleared until `__anext__` cleanup (`_types.py:3118-3125`) via `cleanup_hooks`; earlier fix required `await stream._run_cleanup()` | Subsequent calls fail "already running" | Leak if consumer forgets to `async for` to completion or `await get_final_response()` |

## Future Considerations

* Introduce a discriminated `TerminalOutcome[WorkflowRunResult, WorkflowErrorDetails]` that can be returned without raising, preserving prefix outputs on failures for non-streaming callers (currently only streaming prefix exposes them `tests/workflow/test_functional_workflow.py:320-351`).
* Add a `WorkflowRunHandle` / `Job` future with `wait()`, `join()`, `cancel(wait_only: bool)` and multi-waiter broadcast (e.g., `asyncio.Future` fan-out) instead of the single-consumer `ResponseStream`; document wait-vs-run cancellation separation.
* Materialize `CANCELLED` transition by wiring `asyncio.CancelledError` → `status(CANCELLED)` + `WorkflowCheckpointException` or explicit `workflow.cancel()` API; today `_events.py:67` is dead for workflows.
* Harden ownership: `deepcopy` event data at publication (`WorkflowEvent(data=copy.deepcopy(output))`) or return `Sequence` view with `frozen` data classes; consider `freeze()` helper for `Message.contents` lists to prevent post-publication mutation (parallel to `python/packages/core/agent_framework/_workflows/_executor.py:274-276` capture for invoke).
* Document and test publication ordering guarantee explicitly ("waiter release only after complete outcome visible") as a property test over `workflow.status` vs `WorkflowRunResult` vs `status_timeline()`, across both streaming and non-streaming paths.
* Expose waiter multiplicity tests: two concurrent `asyncio.gather(workflow.run(...), workflow.run(...))` expecting distinct errors; two `await stream.get_final_response()` from separate tasks on same stream expecting same cached identity.

## Questions / Gaps

* **No direct evidence of multi-waiter or late-waiter broadcast tests.** Grep of `tests/workflow` finds no `asyncio.gather` over a single `workflow.run`; search boundary was `rg` across `tests/` and `_types.py`. Reported as `No clear evidence found`.
* **Whether .NET `CANCELLED` path is implemented.** Only `WorkflowExecutor` switches on `CANCELLED` (`_workflow_executor.py:627-630`); core engine never sets it. Dotnet hosts may differ but were not analyzed within isolation rules.
* **`WorkflowRunState.CANCELLED` vs `asyncio.CancelledError` causality not documented.** Exception `raise` path preserves chain but `WorkflowErrorDetails` doesn't capture `BaseExceptionGroup` or nested `__notes__`; not verified for Python 3.11+ `BaseExceptionGroup`.
* **Alias-free guarantee for structured `value` on `AgentResponse.value`.** Lazy parsing (`_value_parsed` flag `_types.py:2650-2667`) may alias echoed text; deep diagnostic absent.
* **Distributed ownership with durable-task checkpoint store concurrency.** `python/packages/durabletask` hosts an external checkpoint store; its concurrency contract beyond `WorkflowCheckpointException` was out of scope for this isolated source analysis.

---

Generated by `Dimension 01.16: Result and Outcome Publication Contract` against `agent-framework`.
