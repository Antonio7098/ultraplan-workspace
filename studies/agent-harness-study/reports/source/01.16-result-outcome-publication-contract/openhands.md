# Source Analysis: openhands

## Result and Outcome Publication Contract

### Source Info

| Field | Value |
|-------|-------|
| Name | openhands |
| Path | `studies/agent-harness-study/sources/openhands` |
| Language / Stack | Python (Pydantic SDK, FIFOLock, EventLog FileStore, FastAPI/REST+WebSocket for remote) |
| Analyzed | 2026-09-01 |

## Summary

OpenHands has no typed `Future[T]`/`Promise`/`Result<T,E>` for an agent run. The **work return** is an open-ended `Event` stream (`ActionEvent`/`ObservationEvent`/`AgentErrorEvent`/`ConversationErrorEvent` + `MessageEvent`) and the **terminal outcome** is a single mutable enum field `ConversationState.execution_status: ConversationExecutionStatus` (`_sdk_inspect/sdk/conversation/state.py:36-73`). The `Agent` writes its "result" by emitting a `FinishAction`/`FinishObservation` via `FinishTool` (`_sdk_inspect/sdk/tool/builtins/finish.py:21-69`), but the runtime does not capture its payload as a typed value — the finish message is just another observation whose side-effect is setting `execution_status = FINISHED` inside `_ActionBatch.finalize` (`_sdk_inspect/sdk/agent/agent.py:206-237`, `_sdk_inspect/sdk/conversation/impl/local_conversation.py:468-472`). Failures are split across two exclusive channels: per-tool `AgentErrorEvent` (recoverable, fed back to LLM) and conversation-level `ConversationErrorEvent` + `execution_status=ERROR` (terminal, with `ConversationRunError` wrapping the original exception via `raise ... from e` – ` _sdk_inspect/sdk/conversation/impl/local_conversation.py:873-888`). There is no explicit run-handle; `LocalConversation.run()` is a synchronous blocking loop holding `FIFOLock` (`_sdk_inspect/sdk/conversation/state.py:218-220`, `_sdk_inspect/sdk/conversation/fifo_lock.py:32-180`) and `RemoteConversation.run(blocking=True)` simulates a join via a `queue.Queue` + REST polling fallback (`_sdk_inspect/sdk/conversation/impl/remote_conversation.py:976-1127`). Immutability is provided by `Event` being `frozen=True` (`_sdk_inspect/sdk/event/base.py:14`) but no deep-freeze on observation payloads, and publication ordering relies on emit-then-mark-finished inside the same lock acquisition (`_sdk_inspect/sdk/agent/agent.py:462-473`). No concurrent fan-out handle, no waiter-only cancellation, and no idempotent repeated `get_result()` API were found.

## Rating

**4/10**

Rationale: Work-return vs terminal-outcome is informally separate (event stream vs `execution_status`) but not typed; `FinishTool` payload is not surfaced as a `RunResult` value. Error/cancel/fault combinations collapse to a single `ERROR` with string `code`/`detail` and lose traceback beyond `str(e)`. Publication is ordered for the single `run()` caller (FIFOLock + emit-before-status), but there is no stable multi-waiter contract, no repeated/late `Future` observation, no copy/freeze on the terminal value, and cancellation of a waiter cannot be distinguished from cancelling the run.

## Evidence Collected

Every entry MUST include a file path with line numbers. Format: `path/to/file.ts:NN`.

| Area | Evidence | File:Line |
|------|----------|-----------|
| Work-function return types | `Agent.step(conversation, on_event, on_token)` returns `None`; tool executors return `Observation`, not `Result` — `ToolExecutor.__call__(action, conversation) -> ObservationT` | `_sdk_inspect/sdk/agent/agent.py:475-603`, `_sdk_inspect/sdk/tool/tool.py:130-180` |
| Terminal outcome type | `ConversationExecutionStatus` enum with `IDLE/RUNNING/PAUSED/WAITING_FOR_CONFIRMATION/FINISHED/ERROR/STUCK/DELETING` and `is_terminal() -> FINISHED|ERROR|STUCK` | `_sdk_inspect/sdk/conversation/state.py:36-73` |
| Terminal-outcome field | `ConversationState.execution_status: ConversationExecutionStatus` plus `max_iterations`, `blocked_actions/messages`, `last_user_message_id` | `_sdk_inspect/sdk/conversation/state.py:109-165` |
| Work-return payload (finish) | `FinishAction(message: str)` → `FinishObservation` via `FinishExecutor`; `FinishTool` is `ToolDefinition[FinishAction, FinishObservation]` with `readOnlyHint=True` | `_sdk_inspect/sdk/tool/builtins/finish.py:21-115` |
| Finish truncation | `_ActionBatch._truncate_at_finish` discards tool calls after `FinishTool.name` and logs warning | `_sdk_inspect/sdk/agent/agent.py:128-154` |
| Finish → terminal transition | `_ActionBatch.finalize` calls `mark_finished` (`state.execution_status = FINISHED`) only if `has_finish` and last action not blocked; iterative-refinement hook may inject `MessageEvent` instead | `_sdk_inspect/sdk/agent/agent.py:206-237` |
| Finish emission ordering | `_execute_actions` does `batch = prepare(); batch.emit(on_event); batch.finalize(...)` with `mark_finished = setattr(state, "execution_status", FINISHED)` | `_sdk_inspect/sdk/agent/agent.py:447-473` |
| Success observation | `ObservationEvent(observation: Observation, action_id, tool_call_id)` + `to_llm_message() -> Message(role="tool")` | `_sdk_inspect/sdk/event/llm_convertible/observation.py:31-72` |
| Per-tool failure | `AgentErrorEvent(error: str, tool_name, tool_call_id, source="agent")` created in `ParallelToolExecutor._run_safe` and `Agent._execute_action_event`; fed back to LLM via `to_llm_message` | `_sdk_inspect/sdk/event/llm_convertible/observation.py:123-168`, `_sdk_inspect/sdk/agent/parallel_executor.py:88-135`, `_sdk_inspect/sdk/agent/agent.py:945-960` |
| Conversation-level failure | `ConversationErrorEvent(code: str, detail: str, source="environment")` not LLM-convertible, signals terminal `ERROR` | `_sdk_inspect/sdk/event/conversation_error.py:10-35` |
| Runtime-fault wrapping | `ConversationRunError(conversation_id, original_exception, persistence_dir)` raised with `raise ... from e`, chain preserves `__cause__`; detail line includes log path | `_sdk_inspect/sdk/conversation/exceptions.py:24-75`, `_sdk_inspect/sdk/conversation/impl/local_conversation.py:885-888` |
| Publication gating (max iterations) | Run loop checks `iteration >= max_iteration_per_run`: if already `FINISHED` preserves it, else sets `ERROR` and emits `ConversationErrorEvent(code="MaxIterationsReached")` | `_sdk_inspect/sdk/conversation/impl/local_conversation.py:850-872` |
| Publication gating (unhandled exception) | `except Exception as e: state.execution_status=ERROR; on_event(ConversationErrorEvent(code=e.__class__.__name__, detail=str(e))); raise ConversationRunError(...)` | `_sdk_inspect/sdk/conversation/impl/local_conversation.py:873-888` |
| Stuck as terminal | `StuckDetector.is_stuck()` → `state.execution_status = STUCK` (also `is_terminal()`) | `_sdk_inspect/sdk/conversation/impl/local_conversation.py:812-820`, `_sdk_inspect/sdk/conversation/state.py:61-72` |
| Wait / join APIs (local) | `LocalConversation.run() -> None` is **synchronous blocking**; no `Future`, `Promise`, `join()`, `wait()` found; `send_message` also synchronous | `_sdk_inspect/sdk/conversation/impl/local_conversation.py:745-888`, `_sdk_inspect/sdk/conversation/base.py:130-160` |
| Wait / join APIs (remote) | `RemoteConversation.run(blocking: bool, poll_interval, timeout)` + `_wait_for_run_completion` draining `queue.Queue` of terminal status events + fallback REST `GET /conversations/{id}` polling; `blocking=False` returns immediately | `_sdk_inspect/sdk/conversation/impl/remote_conversation.py:976-1127` |
| Remote waiter notification | `queue.Queue` `_terminal_status_queue` fed by `WebSocket` callback `run_complete_callback`; `queue.get(timeout=poll_interval)` vs `Empty` fallback | `_sdk_inspect/sdk/conversation/impl/remote_conversation.py:844-888`, `1072-1090` |
| Concurrent / late access primitive | `ConversationState._lock: FIFOLock` (reentrant, FIFO `Condition` queue) plus `with self._state:` context manager on every `send_message`/`run` mutation | `_sdk_inspect/sdk/conversation/state.py:218-230`, `_sdk_inspect/sdk/conversation/state.py:540-560`, `_sdk_inspect/sdk/conversation/fifo_lock.py:32-110` |
| Repeated access | `EventLog.__getitem__`/`__iter__` backed by `FileStore` files (`BASE_STATE` + `events/<id>.json`); stale index triggers `rebuild_from_disk` | `_sdk_inspect/sdk/conversation/event_store.py:40-180` |
| Value ownership / copying | `Event.model_config = ConfigDict(frozen=True)`; `LocalConversation.fork()` deep-copies events via `event.model_copy(deep=True)` and `copy.deepcopy(agent_state)`; `ask_agent` does `model_copy(deep=True)` for LLM | `_sdk_inspect/sdk/event/base.py:14`, `_sdk_inspect/sdk/conversation/impl/local_conversation.py:344-406`, `_sdk_inspect/sdk/conversation/impl/local_conversation.py:1062-1065` |
| Alias / mutation after publication | `EventLog.append` writes `event.model_dump_json` to disk and caches id->idx; `__iter__` yields `Event.model_validate_json(txt)` fresh copy, not alias, but observation payloads are `BaseModel` without freeze guarantee on internal dicts | `_sdk_inspect/sdk/conversation/event_store.py:130-180` |
| Autosave publication | `ConversationState.__setattr__` auto-persists `BASE_STATE` on any public field change when `_autosave_enabled` and `_fs` present, plus fires `_on_state_change -> ConversationStateUpdateEvent` | `_sdk_inspect/sdk/conversation/state.py:380-445` |
| Error cause / stack retention | `_execute_action_event` logs with `exc_info=True` but `AgentErrorEvent` stores only `str(e)`; `ConversationErrorEvent` stores `code`/`detail` strings, no traceback field | `_sdk_inspect/sdk/agent/parallel_executor.py:116-123`, `_sdk_inspect/sdk/agent/agent.py:945-955` |
| Parallel result/error interleaving | `ParallelToolExecutor.execute_batch` with `ThreadPoolExecutor` + `ResourceLockManager.lock(*keys)` per `DeclaredResources`; exceptions mapped to `AgentErrorEvent` per task, preserving input order `return [future.result() ...]` | `_sdk_inspect/sdk/agent/parallel_executor.py:34-135` |
| Late-waiter via persistence | `ConversationState.create` resume path reads `BASE_STATE` JSON and reattaches `EventLog` from `FileStore`; caller can reopen same `persistence_dir` to see terminal `execution_status` + event history | `_sdk_inspect/sdk/conversation/state.py:295-380` |
| Lock documentation warning | `EventLog` doc notes `LocalFileStore` `flock()` "does NOT work reliably on NFS mounts" | `_sdk_inspect/sdk/conversation/event_store.py:28-35` |

## Answers to Dimension Questions

### 1. Is a work return distinct from the runtime's terminal outcome?
**Partially – distinct concepts, weakly typed.** The work-return is the tool-layer `Observation` produced by any `ToolExecutor` (`_sdk_inspect/sdk/tool/tool.py:130`) and surfaced as `ObservationEvent`/`AgentErrorEvent` in the event log. The terminal outcome is **not** that observation but the `ConversationExecutionStatus` enum value on `ConversationState` (`_sdk_inspect/sdk/conversation/state.py:36-73`). The only tool whose execution mutates the terminal outcome is `FinishTool`: `_ActionBatch.finalize` translates the presence of a `FinishAction` into `execution_status = FINISHED` (`_sdk_inspect/sdk/agent/agent.py:468-472`), while all other tools leave the outcome unchanged. However there is no separate typed `RunResult<T>`; the agent's `step()` returns `None` (`_sdk_inspect/sdk/agent/agent.py:475`) and the finish *message* is buried inside `FinishAction.message` (`_sdk_inspect/sdk/tool/builtins/finish.py:22`) and surfaced only as the `FinishObservation` text, not as a `value` field on a terminal handle. Reconstruction of "what was the result" requires scanning `EventLog` for `FinishAction.message`.

### 2. Which result, failure, cancellation, and runtime-fault combinations are valid?
Valid terminal combinations via `execution_status`:

| Combination | Valid? | Mechanism |
|-------------|--------|-----------|
| **Success (value only)** | Yes | `FinishTool` → `FINISHED`; `FinishObservation.from_text(action.message)` ; `ConversationState.last_user_message_id` style tracking irrelevant (`_sdk_inspect/sdk/tool/builtins/finish.py:63-69`, `_sdk_inspect/sdk/agent/agent.py:237`) |
| **Per-tool failure (recoverable, run continues)** | Yes | Tool exception → `AgentErrorEvent(error=str(e), tool_call_id)` ; `execution_status` stays `RUNNING`; event fed back to LLM as `Message(role="tool")` (`_sdk_inspect/sdk/agent/parallel_executor.py:88-135`, `_sdk_inspect/sdk/event/llm_convertible/observation.py:123-168`) |
| **Conversation failure (terminal)** | Yes | Unhandled exception in `run()` loop or `max_iterations` → `execution_status=ERROR` + `ConversationErrorEvent(code=e.__class__.__name__, detail=str(e))` + `ConversationRunError` raised to caller (`_sdk_inspect/sdk/conversation/impl/local_conversation.py:850-888`) |
| **Stuck** | Yes | `StuckDetector.is_stuck()` → `STUCK` (also `is_terminal()`) (`_sdk_inspect/sdk/conversation/impl/local_conversation.py:812-820`) |
| **Interrupted / paused** | Yes, but not terminal | `PAUSED` / `WAITING_FOR_CONFIRMATION` are valid non-terminal states; they deliberately do **not** satisfy `is_terminal()` (`_sdk_inspect/sdk/conversation/state.py:61-72`, `_sdk_inspect/sdk/conversation/impl/local_conversation.py:777-830`) |
| **Cancellation** | No explicit `CANCELLED` state | `CancelledError`-like concept is mapped to `ERROR` or `FINISHED` via tool error; `asyncio.CancelledError` handling is only in `ACPAgent._cancel_inflight_tool_calls` (`_sdk_inspect/sdk/agent/acp_agent.py:1168`) – no `ConversationExecutionStatus.CANCELLED` |
| **Value + Error same task** | No (exclusive) | `_ActionBatch` splits per-tool outcomes; each tool call yields either `ObservationEvent` or `AgentErrorEvent` never both (`_sdk_inspect/sdk/agent/parallel_executor.py:88-108`, `_sdk_inspect/sdk/event/llm_convertible/observation.py:31-168`) |
| **Timeout** | Valid as `ERROR` subtype | `max_iteration_per_run` → `ConversationErrorEvent(code="MaxIterationsReached")` (`_sdk_inspect/sdk/conversation/impl/local_conversation.py:859-871`); LLM retry helpers (`_sdk_inspect/sdk/llm/utils/retry_mixin.py:40-100`) also map timeout to `ERROR` |

No explicit `Cancelled` or `Aborted` terminal; user "cancel" is approximated by `pause()` or `reject_pending_actions()`.

### 3. What happens when work returns both a value and an error?
This is **excluded by construction at the tool level and by the `_ActionBatch` aggregation at the batch level**. Each `ActionEvent` execution via `ParallelToolExecutor._run_safe` takes the single `try: tool(action)` branch: either it returns an `Observation` which is wrapped as `ObservationEvent`, or it raises and is converted to a single `AgentErrorEvent` (`_sdk_inspect/sdk/agent/parallel_executor.py:88-135`; `_sdk_inspect/sdk/agent/agent.py:905-965`). A tool cannot return both because its signature is `-> Observation` with exceptions as the error channel. The special case `FinishTool` is truncated to a single finish (`_sdk_inspect/sdk/agent/agent.py:128-154`), so a batch containing a value tool succeeded + an error tool failed results in **two separate events** in the log (one `ObservationEvent`, one `AgentErrorEvent`) but the conversation `execution_status` reflects only the terminal signal (often still `RUNNING` because the error was per-tool). There is no `Result(value, error)` tuple. At the conversation level, a run that both emitted a `FINISHED` and raised an exception preserves the earlier status: the `except Exception` at `_sdk_inspect/sdk/conversation/impl/local_conversation.py:873` overwrites to `ERROR`, but the `max_iterations` guard explicitly avoids overwriting `FINISHED` (`_sdk_inspect/sdk/conversation/impl/local_conversation.py:854-858`), showing the design intent is "error wins unless already finished".

### 4. Can a terminal outcome expose a partial or mutable value?
**Partial: No via terminal status, Yes via late `EventLog` read.** The terminal outcome itself (`execution_status`) is atomic (single enum write under `FIFOLock` + autosave in `__setattr__` – `_sdk_inspect/sdk/conversation/state.py:380-420`). However the *payload* that gives meaning to `FINISHED` (the `FinishAction.message` and any preceding `ObservationEvent`s) is assembled incrementally: `EventLog.append` is called per emitted event inside `_ActionBatch.emit` (`_sdk_inspect/sdk/agent/agent.py:187-204`), so a reader polling `execution_status` between emits could see `RUNNING` with a partial log, and after `FINISHED` could still observe that `AgentErrorEvent`s coexist with the finish. There is no transactional "publish all events then flip status" – the status flip `mark_finished()` happens *after* `emit` in `_execute_actions` (`_sdk_inspect/sdk/agent/agent.py:462-473`), but within the same `agent.step` call which is itself inside a `with self._state:` block only at the outer `run()` loop (`_sdk_inspect/sdk/conversation/impl/local_conversation.py:773`), so batch-internal races exist when `max_workers>1` (events emitted serially in action order via `emit` even though execution was parallel – `_sdk_inspect/sdk/agent/parallel_executor.py:55-72` ensures order restoration).

**Mutable: Partially mitigated but not strongly guaranteed.** `Event` is `frozen=True` (`_sdk_inspect/sdk/event/base.py:14`) so fields cannot be reassigned after `__init__`. The log on disk is append-only JSON lines via `FileStore.write` (`_sdk_inspect/sdk/conversation/event_store.py:130-180`) and reads via `model_validate_json` produce a fresh instance per access (`_sdk_inspect/sdk/conversation/event_store.py:110-145`), i.e., copy-on-read. But the in-memory `Observation` objects held by `ObservationEvent.observation` are Pydantic models whose `.to_llm_content` may contain mutable `list[TextContent | ImageContent]` (`_sdk_inspect/sdk/tool/schema.py:1-50`), and no `deepcopy`/`freeze` is applied on `append` – `append` stores `event.model_dump_json` (`_sdk_inspect/sdk/conversation/event_store.py:147`), which serializes but does not freeze the caller-held reference. A caller retaining a pre-append reference to the same `ObservationEvent` could still mutate its Python dict fields before the disk write (though `frozen=True` blocks attribute reassignment, it does not freeze nested dicts if the observation carries `dict[str, Any]`). `fork()` explicitly does `model_copy(deep=True)` (`_sdk_inspect/sdk/conversation/impl/local_conversation.py:386`) to avoid this, indicating the general path does not.

### 5. Do concurrent, repeated, and late waiters observe the same committed facts?
**Single waiter – Yes. Concurrent fan-out – No handle to test. Repeated & late – Yes via persistence, with caveats.**

- **Single waiter (the `run()` caller):** `with self._state:` + `FIFOLock` ensures the loop sees a consistent snapshot; the comment at `_sdk_inspect/sdk/conversation/impl/local_conversation.py:836-842` explicitly preserves `FINISHED` if a concurrent `send_message` raced to set `IDLE`, showing awareness of the race. For `RemoteConversation`, `_wait_for_run_completion` drains the `queue.Queue` first (WebSocket push) then falls back to `GET /conversations/{id}` polling, and on fallback does `_state.refresh_from_server()` + `reconcile()` to align (`_sdk_inspect/sdk/conversation/impl/remote_conversation.py:1086-1125`).

- **Concurrent waiters on same run:** Not supported. `LocalConversation` has no `Future`/`share()` handle; two threads calling `run()` concurrently contend on the same `FIFOLock` and the second will either spin in `with self._state:` or observe the already-terminal status and break (`_sdk_inspect/sdk/conversation/impl/local_conversation.py:777-809`). `RemoteConversation` has one `queue.Queue` and one `WebSocket` connection per instance; cloning the instance is required for a second waiter, and there is no coordination to give both waiters the same terminal `execution_status` atomically. No test asserts fan-out.

- **Repeated waiters (idempotent poll):** A caller may re-read `state.execution_status` after `run()` returns or after `close()`; the value is stable because only `send_message` resets `FINISHED`→`IDLE` (`_sdk_inspect/sdk/conversation/impl/local_conversation.py:704-710`) and no other path mutates a terminal state except `pause`/`error` overwrites. But there is no `get_result()` cache – repeated calls must re-scan `EventLog` to rebuild the logical result (`ConversationState.get_unmatched_actions` pattern – `_sdk_inspect/sdk/conversation/state.py:471-530`).

- **Late waiters:** Supported via persistence. `ConversationState.create` resume (`_sdk_inspect/sdk/conversation/state.py:320-360`) re-reads `BASE_STATE` and rebuilds `EventLog` from `FileStore`, so a new `LocalConversation` with the same `persistence_dir` + `conversation_id` observes the identical `execution_status` and event history. The `EventLog` `reconcile()` path (`_sdk_inspect/sdk/conversation/event_store.py:40-180`) syncs `disk_length` against in-memory `self._length` if another process wrote while blocked on `fs.lock`. However the staleness window is non-zero: `EventLog.__getitem__` rebuilds on `KeyError` but concurrent readers outside a `with self._state:` block can observe torn reads (documented NFS caveat).

### 6. Can waiting be cancelled without cancelling the run itself?
**No.** There is no `WaitHandle.cancel()` distinct from run control. The only cancellation-adjacent operations are:

- `LocalConversation.pause()` sets `execution_status = PAUSED` and emits `PauseEvent` (`_sdk_inspect/sdk/conversation/impl/local_conversation.py:927-951`) – this **pauses the run itself**, not just a waiter. The `run()` loop checks `if execution_status in (PAUSED, STUCK): break` as the next iteration's gate (`_sdk_inspect/sdk/conversation/impl/local_conversation.py:777-781`). Thus waiter cancellation == run cancellation.

- `RemoteConversation.pause()` POSTs to `/conversations/{id}/pause` (`_sdk_inspect/sdk/conversation/impl/remote_conversation.py:1233-1238`) – same semantic.

- `reject_pending_actions()` transitions `WAITING_FOR_CONFIRMATION -> IDLE` and emits `UserRejectObservation`s (`_sdk_inspect/sdk/conversation/impl/local_conversation.py:896-925`) – it cancels *pending tool approvals*, not a waiter.

- `ACPAgent._cancel_inflight_tool_calls` (`_sdk_inspect/sdk/agent/acp_agent.py:1168`) cancels child tool futures but is invoked only when the run itself is being torn down (`close()`/`run` exit).

A thread blocked in `RemoteConversation._wait_for_run_completion` on `queue.get(timeout=poll_interval)` (`_sdk_inspect/sdk/conversation/impl/remote_conversation.py:1076`) can only be freed by the terminal status arriving or timeout expiration as `ConversationRunError`; there is no `cancel(token)` that unblocks the waiter while leaving `execution_status=RUNNING` on the server.

### 7. Who owns values and diagnostics after publication?
- **Values (observations, finish message):** The runtime owns the durable copy (serialized JSON files per event under `events/<uuid>.json` via `EventLog.append` – `_sdk_inspect/sdk/conversation/event_store.py:147-153`). The caller receives **aliased Pydantic objects on read** (`Event.__iter__` → `model_validate_json` fresh instance) but the underlying `Observation` payload may still alias nested structures if the observer cached an earlier reference. No `Buffer` transfer or `freeze()` – ownership is copy-by-serialization for persistence, share-by-reference in-memory until `fork()` which uses `copy.deepcopy` (`_sdk_inspect/sdk/conversation/impl/local_conversation.py:394-395`). There is no explicit `close()` required to release value ownership; `LocalConversation.close()` closes executors, not values (`_sdk_inspect/sdk/conversation/impl/local_conversation.py:976-1012`).

- **Diagnostics:** 
  - Per-tool diagnostics live as `AgentErrorEvent.error: str` (`_sdk_inspect/sdk/event/llm_convertible/observation.py:146`) with no traceback. 
  - Terminal diagnostics live as `ConversationErrorEvent(code, detail)` (`_sdk_inspect/sdk/event/conversation_error.py:17-32`) where `code = e.__class__.__name__` and `detail = str(e)` (`_sdk_inspect/sdk/conversation/impl/local_conversation.py:877-882`). 
  - The original exception is owned only transiently by the `ConversationRunError.original_exception` field (`_sdk_inspect/sdk/conversation/exceptions.py:32-45`) and chained via `raise ... from e`; after `run()` returns via exception the caller owns it but a late poller that only reads `EventLog` sees the stringified copy.
  - `FIFOLock` and `FileStore` retain no diagnostic ownership; logs concern `logger.debug/warning` with `exc_info=True` only at emit time (`_sdk_inspect/sdk/agent/parallel_executor.py:116-123`, `_sdk_inspect/sdk/agent/agent.py:945-955`).

After publication, diagnostics are immutable strings on disk; values are JSON-serialized, so their runtime type is lost (observation is re-hydrated as the discriminated union variant via `Event.model_validate_json` – `_sdk_inspect/sdk/conversation/event_store.py:130-145`).

### 8. Does waiter release occur only after the complete outcome is visible?
**For the local synchronous caller, yes. For remote/polling/late readers, only eventually.**

- **Local `run()`:** The sequence is `ParallelToolExecutor.execute_batch` → `batch.emit` (all `ObservationEvent`/`AgentErrorEvent` appended and `_save_base_state` via autosave) → `batch.finalize` → `mark_finished` → `state.execution_status = FINISHED` (`_sdk_inspect/sdk/agent/agent.py:462-473` and `_sdk_inspect/sdk/conversation/state.py:380-420` autosave). The `run()` loop then requires one more `with self._state:` iteration to observe `FINISHED` at the top-of-loop guard before `break` (`_sdk_inspect/sdk/conversation/impl/local_conversation.py:783-809`). So the blocking caller cannot return until both the events and the `FINISHED` flag are durable in `BASE_STATE`. The `FIFOLock` ensures no interleaving waiter sees `FINISHED` without the preceding events.

- **Remote `run(blocking=True)`:** The preferred path waits on `queue.get(timeout=poll_interval)` which is signaled only after the WebSocket delivers `ConversationStateUpdateEvent(key="execution_status", value=terminal)` (`_sdk_inspect/sdk/conversation/impl/remote_conversation.py:1072-1088` + `_sdk_inspect/sdk/conversation/state.py:400-420` callback). That terminal status is sent **after** the server has appended the final events (server-side equivalent of `batch.emit`). However the fallback path of 3 consecutive `GET` polls (`_sdk_inspect/sdk/conversation/impl/remote_conversation.py:1106-1127`) can fire if WebSocket is delayed, and between the third `200` and the subsequent `refresh_from_server()` + `reconcile()` there is a window where the poller has declared completion but the final event reconciliation has not yet run. The code explicitly orders `refresh_from_server()` before `reconcile()` and before `return`, mitigating but not eliminating the gap where a *different* poller could observe `execution_status=FINISHED` via REST without all events.

- **Durability caveat:** `EventLog.append` uses `FileStore.lock(LOCK_TIMEOUT_SECONDS=30)` (`_sdk_inspect/sdk/conversation/event_store.py:133-153`) – a crash between final event write and `BASE_STATE` write could leave `execution_status` stale while events are durable, or vice versa, since they are separate files. No single atomic commit transaction is used.

## Architectural Decisions

| Decision | Rationale / Evidence | Consequence |
|----------|----------------------|-------------|
| **Enum terminal state instead of typed `Future[T]`** (`_sdk_inspect/sdk/conversation/state.py:36-73`) | Simplest persistence model: single JSON key `execution_status` in `BASE_STATE` readable by any restart; matches UI polling over REST | Loses value typing; callers must synthesize result by scanning `EventLog`; no `await` composition or typed error unions |
| **`FinishTool` as the success carrier** (`_sdk_inspect/sdk/tool/builtins/finish.py:21-115`) | Reuses LLM tool-call channel; no extra RPC endpoint; message is human-readable | Success value is unstructured `str`; no schema, no machine-readable payload; truncation logic silently discards sibling tool calls |
| **Two-tier errors: `AgentErrorEvent` (recoverable) vs `ConversationErrorEvent` (terminal)** (`_sdk_inspect/sdk/event/llm_convertible/observation.py:123-168`, `_sdk_inspect/sdk/event/conversation_error.py:10-35`) | Allows agent self-correction loop: per-tool errors returned as `Message(role="tool")` for next `step` | Splits error handling burden: callers must inspect two streams; terminal error detail is `str(e)` only |
| **FIFOLock + autosave on `__setattr__`** (`_sdk_inspect/sdk/conversation/state.py:380-445`) | Guarantees "overwrite → persist → notify WebSocket" inside the lock; WebSocket subscribers see ordered `ConversationStateUpdateEvent`s | Contention on single lock for all mutations; lock does not protect cross-file atomicity (BASE_STATE vs event files) |
| **`EventLog` file-per-event append-only store** (`_sdk_inspect/sdk/conversation/event_store.py:40-180`) | Crash-safe, process-safe via `FileStore.lock`, resumable via `create(persistence_dir)` replay | No index for outcome lookup; publication requires scanning entire log; NFS `flock` unreliable (doc warning) |
| **Parallel executor restores input order** (`_sdk_inspect/sdk/agent/parallel_executor.py:55-72` + `_sdk_inspect/sdk/agent/agent.py:198-204`) | Deterministic LLM history: `on_event` called in original LLM tool-call order even when threads race | Adds latency (collect-then-emit barrier); a slow tool delays publication of an already-finished sibling's observation |
| **Remote polling fallback after 3 terminal `GET`s** (`_sdk_inspect/sdk/conversation/impl/remote_conversation.py:1059-1127`) | Tolerates WebSocket drops; still provides bounded waiter release without infinite hang | Weakens "complete outcome before release": REST may report `FINISHED` before final `reconcile()` has fetched all events |

## Notable Patterns

- **Truncate-at-finish barrier:** `_ActionBatch._truncate_at_finish` (`_sdk_inspect/sdk/agent/agent.py:128-154`) enforces that at most one terminal signal per LLM response, preventing the "both value and continue" paradox.
- **Blocked-reason partitioning:** `pop_blocked_action` check during `prepare` (`_sdk_inspect/sdk/agent/agent.py:170-176` + `_sdk_inspect/sdk/conversation/state.py:451-459`) converts disallowed actions into `UserRejectObservation(rejection_source="hook")` without executing the tool, keeping the terminal outcome exclusive.
- **Chain-then-commit preserve-FINISHED guard:** `max_iterations` and the `except Exception` block both special-case `if execution_status == FINISHED: break/preserve` (`_sdk_inspect/sdk/conversation/impl/local_conversation.py:854-858`, `836-853` comment) to avoid overwriting an already-published success with an `ERROR`.
- **Drain-stale-queue before run:** `RemoteConversation.run` drains `queue.get_nowait()` in a loop before triggering `POST /run` (`_sdk_inspect/sdk/conversation/impl/remote_conversation.py:995-1001`) so a late waiter from a prior run does not spuriously resolve the new run.
- **Copy-on-fork as ownership isolation:** `LocalConversation.fork` deep-copies events + `agent_state` + stats (`_sdk_inspect/sdk/conversation/impl/local_conversation.py:344-406`) – the only place where explicit ownership transfer is implemented.

## Tradeoffs

| Tradeoff | Pro | Con |
|----------|-----|-----|
| **Synchronous `run()->None` with enum status vs `Future[Outcome]`** | No async plumbing, easy to test sequentially, trivial persistence | No composition, no cancellation token, no multi-waiter fan-out; callers must poll `state.execution_status` or `EventLog` |
| **Frozen `Event` but mutable payload** (`_sdk_inspect/sdk/event/base.py:14`) | Prevents field reassignment bugs; cheap copy-on-read via JSON | Nested dict/list inside `Observation` still mutable; no `deepcopy` on `append` so aliasing hazard remains |
| **File-per-event durability vs atomic commit** | Simple I/O, survives process death, resumable | `BASE_STATE` + event files not atomic; crash between writes can leave divergent terminal flag vs events |
| **Ordered emit barrier in parallel executor** | Deterministic LLM message construction (`events_to_messages` grouping by `llm_response_id` – `_sdk_inspect/sdk/event/base.py:75-130`) | Throughput limited; all parallel results withheld until slowest finishes |
| **Stringified error diagnostics** (`code=str(type)`, `detail=str(e)`) | Language-agnostic REST payload; no pickle concerns | Type fidelity and traceback lost for late readers; only `exc_info=True` log lines retain stack |
| **Single `FIFOLock` for all state** | FIFO fairness prevents starvation (`_sdk_inspect/sdk/conversation/fifo_lock.py:32-80`) | Coarse contention: `send_message` during `run` blocks on lock held by `agent.step` LLM call; no read/write split |

## Failure Modes / Edge Cases

- **Finish with sibling errors:** A batch containing `FinishAction` + a failing tool in the same LLM response truncates the non-finish tools only if they appear *after* finish (`_sdk_inspect/sdk/agent/agent.py:136-148`). If the error precedes finish, both an `AgentErrorEvent` and the `FINISHED` status are published – a "partial success" visible only to event scanners, not to status pollers.
- **Max-iterations override race:** If the agent finishes on the final allowed iteration, the loop preserves `FINISHED` (`_sdk_inspect/sdk/conversation/impl/local_conversation.py:854-858`); any other terminal (`STUCK`/`ERROR`/`PAUSED`) is overwritten by `MaxIterationsReached` `ERROR`, potentially masking the true cause.
- **NFS `flock` failure:** `EventLog` doc warns `flock() does NOT work reliably on NFS` (`_sdk_inspect/sdk/conversation/event_store.py:30`) – late waiters on shared storage can observe torn `EventLog` or duplicate `event_id` writes despite `ValueError` check (`_sdk_inspect/sdk/conversation/event_store.py:143-148`).
- **Canceled waiter orphan:** A thread calling `RemoteConversation._wait_for_run_completion` that times out raises `ConversationRunError(TimeoutError)` (`_sdk_inspect/sdk/conversation/impl/remote_conversation.py:1062-1070`) but the server-side run continues; subsequent `send_message` resets `FINISHED->IDLE` (`_sdk_inspect/sdk/conversation/impl/local_conversation.py:704-710`) creating ambiguous "was the timeout a failure?"
- **Stale `EventLog` index:** `EventLog.__getitem__` rebuilds from disk on `KeyError` (`_sdk_inspect/sdk/conversation/event_store.py:95-110`) but concurrent `__iter__` does not lock – a late reader iterating during a concurrent `append` may skip or duplicate the newly-appended event.
- **Blocked-action leak across resume:** `blocked_actions`/`blocked_messages` persist in `BASE_STATE`; a `pop_blocked_action` that is not consumed before crash leaves the entry durable and will be retried on resume, potentially double-rejecting an action.

## Future Considerations

- Introduce a typed `RunOutcome` / `RunHandle` that surfaces `FinishAction.message` as a structured value and `AgentErrorEvent` tracebacks alongside `execution_status`, e.g. `handle = conversation.run_async() -> Future[RunOutcome]` with `outcome.value`, `outcome.error`, `outcome.events`, closing the work-return vs terminal-outcome gap.
- Add **deep freeze / copy-on-append** for `ObservationEvent` payloads (e.g. `model_copy(deep=True)` before `FileStore.write`) or at least a `MappingProxy`-wrapped view for returned `EventLog` snapshots, preventing post-publication mutation.
- Preserve **cause chain and traceback** in `ConversationErrorEvent` (`traceback.format_exc()` + `__cause__`) so remote/late readers can reconstruct diagnostics without relying on server logs.
- Implement **waiter-only cancellation** (`handle.cancel(waiter_only=True)` vs `conversation.cancel(run=True)`) and distinguish `CANCELLED` in `is_terminal()` so pause/cancel semantics do not conflate.
- Provide **idempotent repeated `get_outcome()`** that returns the same `RunOutcome` object (or a `Future` already resolved) without rescanning the log, enabling true concurrent/late waiter fan-out tested via `Future.share()`.
- Make publication **atomic** by writing events + `BASE_STATE` under a single `FileStore.lock` transaction or by staging a `FINISHED` marker file only after all event files are fsync'd, formalizing "outcome visible only after complete facts".

## Questions / Gaps

- No evidence of a `wait(timeout) -> Outcome` or `join()` contract for multiple concurrent waiters; searches for `wait_until_ready`, `wait_all`, `wait_for_sandbox` in `_sdk_inspect/sdk/conversation/` show only internal progress polling, not run-handle fan-out. Confirm intention is poll-loop vs future-share.
- How should an SDK caller retrieve the **machine-readable finish value** without parsing the human `FinishAction.message` string? No structured `result` field on `ConversationState` or `ConversationStats` was found.
- Is the `FINISHED` → `IDLE` reset on next `send_message` (`_sdk_inspect/sdk/conversation/impl/local_conversation.py:704-710`) intended to allow **re-entering** a terminal conversation without explicit `create/fork`? Interaction with `is_terminal()` observers that already cached `FINISHED` is untested.
- What is the expected **exactly-once** guarantee for `AgentErrorEvent` vs `ConversationErrorEvent` when `ParallelToolExecutor` races a tool timeout with the run's `max_iterations` error? No concurrent error-composition test found.
- Does `EventLog.reconcile()` (`_sdk_inspect/sdk/conversation/impl/remote_conversation.py:1124`) have a bounded wait for NFS-durable visibility, or can a late waiter miss events that arrived between `FINISHED` poll and `reconcile`?

---

Generated by `dimension 01.16: Result and Outcome Publication Contract` against `openhands`.
