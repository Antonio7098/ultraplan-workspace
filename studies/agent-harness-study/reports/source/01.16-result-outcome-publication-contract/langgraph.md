# Source Analysis: langgraph

## Result and Outcome Publication Contract

### Source Info

| Field | Value |
|-------|-------|
| Name | langgraph |
| Path | `studies/agent-harness-study/sources/langgraph` |
| Language / Stack | Python (Pregel BSP, langchain_core Runnable) |
| Analyzed | 2026-09-01 |

## Summary

LangGraph does not expose a classic `Future`/`Promise` run-handle. A graph *work return* is a node’s direct Python return (dict update, `Command`, `Send`, `interrupt()` value, `entrypoint.final`) that is translated inside the tick into `task.writes` (`RETURN`, channel writes) and then committed via `PregelLoop.put_writes`. The *terminal outcome* is a separate, higher-level construct derived only after the `PregelLoop` exits: `loop.output = read_channels(channels, output_keys)` plus an orthogonal `interrupts` tuple. `Pregel.invoke`/`ainvoke` synthesizes that terminal outcome by draining `stream(..., stream_mode="values")` (v1 merges `__interrupt__` into the dict, v2 returns `GraphOutput(value, interrupts)`). Failures are never returned as values; they are appended as `(ERROR, exc)` writes, optionally routed to a node-level error handler, and otherwise re-raised through `PregelRunner._panic_or_proceed`/`BackgroundExecutor.__exit__`. Immutability is not enforced on the returned value (aliased channel storage) but durability ordering is enforced via `_delta_write_futs` / `_checkpointer_put_after_previous` gating; `durability="sync"` blocks `invoke` until the checkpoint is durable, while `async`/`exit` return as soon as the in-memory loop is done. There is no multi-waiter fan-out handle; concurrency exists only inside a superstep (multiple `PregelExecutableTask`s racing via `FuturesDict`), and late observation is via `get_state()`/`get_state_history()` on the checkpointer.

## Rating

**5/10**

Rationale: The work-return → terminal-outcome mapping is well-structured (explicit `RETURN`/`ERROR`/`INTERRUPT` reserved keys, distinct `GraphOutput` for v2, single stable `loop.output` published in `PregelLoop._suppress_interrupt`), and error-cause preservation plus checkpoint durability ordering is rigorous. However LangGraph provides no first-class wait/join handle for concurrent/repeated/late waiters (one `SyncQueue` per `stream()` call, `invoke` blocks the caller thread), returns mutable aliased state with no freeze/copy, and leaves `durability="async"`/`"exit"` publication racy from an external reader’s perspective. Tests exercise recovery and error-handler routing but not fan-out or cancellation-without-run-cancellation semantics.

## Evidence Collected

Every entry MUST include a file path with line numbers. Format: `path/to/file.ts:NN`.

| Area | Evidence | File:Line |
|------|----------|-----------|
| Work-function return types | `PregelNode.bound: Runnable` + `DEFAULT_BOUND` lambda passthrough; `NodeBuilder.do(node)` coerces any `RunnableLike` | `libs/langgraph/langgraph/pregel/_read.py:97-234`, `libs/langgraph/langgraph/pregel/main.py:303-314` |
| Work-function wrapping (functional API) | `@task` returns `SyncAsyncFuture[T]`; `@entrypoint` maps `entrypoint.final` via `ChannelWrite(END)+ChannelWrite(PREVIOUS)` and `LastValue` channels | `libs/langgraph/langgraph/func/__init__.py:59-94`, `libs/langgraph/langgraph/func/__init__.py:546-598` |
| Call-primitive Future | `call()` → `SyncAsyncFuture[T] extends concurrent.futures.Future[T]` with `__await__`; impl fetched from `config[CONF][CONFIG_KEY_CALL]` | `libs/langgraph/langgraph/pregel/_call.py:253-298` |
| Task runnable invocation | `run_with_retry` / `arun_with_retry` clears `task.writes` then `task.proc.invoke/ainvoke` | `libs/langgraph/langgraph/pregel/_retry.py:615-617`, `libs/langgraph/langgraph/pregel/_retry.py:738-744` |
| Task write vocabulary (success/failure/cancel/interrupt) | Reserved keys `RETURN`, `ERROR`, `ERROR_SOURCE_NODE`, `INTERRUPT`, `RESUME`, `NO_WRITES`, `TASKS` | `libs/langgraph/langgraph/_internal/_constants.py:7-24` |
| Commit path (success) | `PregelRunner.commit` on success appends `NO_WRITES` if empty, emits `node_finished`, then `put_writes(task.id, task.writes)` | `libs/langgraph/langgraph/pregel/_runner.py:604-613` |
| Commit path (error) | `commit` on exception appends `(ERROR, exc)` and optionally `(ERROR_SOURCE_NODE, task.name)` then `put_writes` | `libs/langgraph/langgraph/pregel/_runner.py:596-603` |
| Commit path (interrupt) | `GraphInterrupt` saves `(INTERRUPT, ...)` (+ `RESUME`) via `put_writes` | `libs/langgraph/langgraph/pregel/_runner.py:584-591` |
| Commit path (cancel) | `asyncio.CancelledError` saved as `(ERROR, CancelledError)` | `libs/langgraph/langgraph/pregel/_runner.py:579-583` |
| Error-handler routing gate | `_should_route_to_error_handler` + `SKIP_RERAISE_SET`/`_handled_exception_ids` prevents fatal re-raise for handled errors | `libs/langgraph/langgraph/pregel/_runner.py:171-175`, `libs/langgraph/langgraph/pregel/_runner.py:301-323` |
| Terminal outcome construction | `SyncPregelLoop._suppress_interrupt` assigns `self.output = read_channels(self.channels, self.output_keys)` on any exit (interrupt or normal) | `libs/langgraph/langgraph/pregel/_loop.py:1349-1355` |
| Terminal outcome publication (sync stream) | `Pregel.stream` enters `SyncPregelLoop` as context manager, runs `while loop.tick(): runner.tick(...); loop.after_tick()`, then checks `loop.status` and calls `run_manager.on_chain_end(loop.output)` | `libs/langgraph/langgraph/pregel/main.py:2934-3032` |
| Terminal outcome publication (invoke v1) | `invoke` drains `stream` with `stream_mode=["updates","values"]`, collects `latest` payload and `interrupts` list, then `return {**latest, INTERRUPT: interrupts}` or `latest` | `libs/langgraph/langgraph/pregel/main.py:3927-3972` |
| Terminal outcome publication (invoke v2) | `invoke` with `version="v2"` yields `GraphOutput(value=latest, interrupts=tuple(interrupts))` via `types.GraphOutput` frozen dataclass | `libs/langgraph/langgraph/types.py:368-378`, `libs/langgraph/langgraph/pregel/main.py:3904-3963` |
| Output channel reading | `read_channels(channels, output_keys)` returns `channels[chan].get()` directly (no copy); empty channel yields `None`/`{}` | `libs/langgraph/langgraph/pregel/_io.py:23-53` |
| Updates vs values mapping | `map_output_updates` filters out `ERROR`/`INTERRUPT` tasks; `map_output_values` emits only when `updated_channels ∩ output_keys != ∅` | `libs/langgraph/langgraph/pregel/_io.py:100-174` |
| Channel immutability promise | BSP model comment: "channel updates from step N are only visible in step N+1, guaranteed immutable for duration of step" | `libs/langgraph/langgraph/pregel/main.py:2974-2978` |
| Durability-gated publication | `_checkpointer_put_after_previous(prev, ...)` waits on `_delta_write_futs` and `prev.result()` before `checkpointer.put`/`aput`; `durability="sync"` awaits `loop._put_checkpoint_fut.result()` before returning | `libs/langgraph/langgraph/pregel/_loop.py:1507-1524`, `libs/langgraph/langgraph/pregel/_loop.py:1759-1778`, `libs/langgraph/langgraph/pregel/main.py:3002-3003` |
| Exit-mode accumulator | `_exit_delta_writes: list[tuple[int,str,str,Any]]` captures non-snapshot delta writes; `_put_exit_delta_writes` stages stub + writes before final checkpoint | `libs/langgraph/langgraph/pregel/_loop.py:213-222`, `libs/langgraph/langgraph/pregel/_loop.py:1201-1280` |
| Wait/Join handle absence (per-invoke queue) | Each `stream()` allocates `stream = SyncQueue()`/`AsyncQueue()` local to the call; no shared handle is returned | `libs/langgraph/langgraph/pregel/main.py:2763`, `libs/langgraph/langgraph/pregel/main.py:3165-3170` |
| Internal fan-out (FuturesDict) | `FuturesDict` maps `Future -> PregelExecutableTask`, wakes `event` when counter 0 or `_should_stop_others` true; `_should_stop_others` cancels siblings on first non-handled failure (exception of `GraphInterrupt` not considered failure) | `libs/langgraph/langgraph/pregel/_runner.py:75-132`, `libs/langgraph/langgraph/pregel/_runner.py:616-634` |
| Panic ordering | `_panic_or_proceed` cancels inflight futures, collects `GraphInterrupt`s into one, re-raises first non-handled error only after done callbacks have fired (`futures.event.wait()`) | `libs/langgraph/langgraph/pregel/_runner.py:649-697` |
| Background executor publication gate | `BackgroundExecutor.__exit__` cancels `__cancel_on_exit__` futures, waits `concurrent.futures.wait(pending)`, then re-raises first `__reraise_on_exit__` error; `AsyncBackgroundExecutor.__aexit__` analogous | `libs/langgraph/langgraph/pregel/_executor.py:93-121`, `libs/langgraph/langgraph/pregel/_executor.py:186-211` |
| Chained-future ordering guarantee | `_call`/`_acall` wrap via `chain_future` so `FuturesDict.on_done`/`commit` fires before the downstream future resolves, ensuring stream order | `libs/langgraph/langgraph/pregel/_runner.py:783-786`, `libs/langgraph/langgraph/pregel/_runner.py:809-841` |
| Cause/stack preservation | `run_with_retry`/`arun_with_retry` call `exc.add_note(f"During task with name '{name}' and id '{id}'")` (py3.11+); `_panic_or_proceed` trims only `EXCLUDED_FRAME_FNAMES` then re-raises original `exc` | `libs/langgraph/langgraph/pregel/_retry.py:642-643`, `libs/langgraph/langgraph/pregel/_retry.py:797-798`, `libs/langgraph/langgraph/pregel/_runner.py:350-358` |
| NodeError diagnostic | `NodeError(node, error)` injected via `CONFIG_KEY_NODE_ERROR`/`CONFIG_KEY_SCRATCHPAD` into error-handler signature | `libs/langgraph/langgraph/errors.py:148-165`, `libs/langgraph/langgraph/_internal/_constants.py:77-80` |
| Cancellation discrimination | `_is_user_raised_cancelled()` checks `asyncio.Task.cancelling()==0` to map user `CancelledError` → `NodeCancelledError`; framework cancel (cancelling>=1) propagates silently | `libs/langgraph/langgraph/pregel/_retry.py:315-335`, `libs/langgraph/langgraph/pregel/_retry.py:785-793` |
| Late-waiter API (checkpointer) | `get_state` / `get_state_history` / `_prepare_state_snapshot` read persisted `CheckpointTuple`; `StateSnapshot.values/next/tasks/interrupts` is the late-observable outcome | `libs/langgraph/langgraph/pregel/main.py:1391-1433`, `libs/langgraph/langgraph/pregel/main.py:1144-1265`, `libs/langgraph/langgraph/types.py:643-661` |
| Copy vs alias in checkpoint | `_put_checkpoint` calls `copy_checkpoint(checkpoint)` before `checkpointer.put(..., copy_checkpoint(checkpoint), ...)` — in-memory dict is copied for the saver, but `loop.output` alias is not | `libs/langgraph/langgraph/pregel/_loop.py:1186` |
| Tests: recovery shows exclusivity | `test_checkpoint_recovery` asserts failed run raises, checkpoint retains input values, `state.tasks[0].error` stringified, retry with durability succeeds | `libs/langgraph/tests/test_pregel.py:5387-5447` |
| Tests: concurrent independent runs | `test_concurrent_execution_thread_safety` runs 10 threads each `graph.invoke({"counter":0})` and asserts independent `{"counter":1}` — no shared outcome | `libs/langgraph/tests/test_pregel.py:5349-5385` |
| Tests: waiter-cancel without run-cancel | `test_cancel_ainvoke_with_async_node` vs `test_cancel_ainvoke_with_sync_node` — async node cancellation propagates, sync node thread orphan continues | `libs/langgraph/tests/test_runtime.py:735-820` |
| Tests: GraphOutput v2 | `test_invoke_v2_graph_output_with_interrupts` asserts `isinstance(result, GraphOutput)` and `result.value`/`result.interrupts` | `libs/langgraph/tests/test_stream_events_v3.py:375-410`, `libs/langgraph/tests/test_stream_events_v3.py:698-773` |

## Answers to Dimension Questions

### 1. Is a work return distinct from the runtime's terminal outcome?
Yes, strictly distinct.

- **Work return** = whatever the node callable returns. For `StateGraph`, typically `dict[str, Any]` that is fed to `ChannelWrite` writers (`libs/langgraph/langgraph/pregel/_read.py:223-234` composes `bound` + `writers`). For functional API, `entrypoint.final` splitting via `ChannelWriteEntry(END)` vs `PREVIOUS` (`libs/langgraph/langgraph/func/__init__.py:583-589`). The raw value is also materialized inside `PregelRunner._call` handling of `RETURN` channel (`libs/langgraph/langgraph/pregel/_runner.py:735-756` captures `__return__`).
- **Terminal outcome** = `PregelLoop.output` / `GraphOutput`. It is not the node’s return. `PregelLoop._suppress_interrupt` computes `self.output = read_channels(self.channels, self.output_keys)` (`libs/langgraph/langgraph/pregel/_loop.py:1350-1355`) after the superstep loop drains. `Pregel.invoke` does not return the node’s value at all; it consumes the `stream(..., stream_mode="values")` iterator, keeps `latest` and `interrupts`, and reconstructs (`libs/langgraph/langgraph/pregel/main.py:3900-3972`):
  - v1: `latest` dict + optional `{"__interrupt__": interrupts}`,
  - v2: `GraphOutput(value=latest, interrupts=tuple(interrupts))` (`libs/langgraph/langgraph/types.py:368-378`).
- Evidence of separation: `map_output_updates` explicitly skips `ERROR`/`INTERRUPT` writes (`libs/langgraph/langgraph/pregel/_io.py:128-129`), while `map_output_values` reflects channel state, not individual writes. A node returning `{"foo": 1}` may trigger zero `output_keys` changes (if reducer drops it), in which case no `values` chunk is emitted and `latest` stays whatever the previous checkpoint held.

### 2. Which result, failure, cancellation, and runtime-fault combinations are valid?
LangGraph’s valid terminal combinations are enumerated by `PregelRunner.commit` and `PregelLoop._suppress_interrupt`:

| Combination | Valid? | How observed |
|-------------|--------|--------------|
| **Success value only** | Yes | `commit` path `else:` writes `(NO_WRITES, None)` if empty then `put_writes` (`libs/langgraph/langgraph/pregel/_runner.py:604-613`); caller gets `latest` dict / `GraphOutput.value` |
| **Success value + Interrupts** | Yes | `GraphInterrupt` → `(INTERRUPT, interrupts)` plus `(RESUME)`; `invoke` merges (`libs/langgraph/langgraph/pregel/main.py:3961-3969` v1 dict merge, `libs/langgraph/langgraph/types.py:369` v2 fields) |
| **Failure (no value)** | Yes | `commit` on non-`GraphInterrupt` exception → `(ERROR, exc)` + optionally `(ERROR_SOURCE_NODE, task.name)` (`libs/langgraph/langgraph/pregel/_runner.py:596-603`); `_panic_or_proceed` re-raises `exc` to the `invoke` caller (`libs/langgraph/langgraph/pregel/_runner.py:657-687`) |
| **Failure routed to handler → Success** | Yes | `node_error_handler_map` + `schedule_error_handler`/`aschedule_error_handler` creates handler task; original `exc` id added to `_handled_exception_ids` and `SKIP_RERAISE_SET` so not re-raised (`libs/langgraph/langgraph/pregel/_runner.py:228-233`, `libs/langgraph/langgraph/pregel/_runner.py:297-323`, `libs/langgraph/langgraph/pregel/_loop.py:1549-1584`) |
| **Cancellation (framework)** | Treated as silent tear-down | `asyncio.CancelledError` in `commit` → `(ERROR, CancelledError)` (`libs/langgraph/langgraph/pregel/_runner.py:579-583`); `_exception()` maps cancelled future to `CancelledError` but `_should_stop_others` returns `False` for handled ids (`libs/langgraph/langgraph/pregel/_runner.py:616-624`) |
| **Cancellation (user-raised)** | Converted to failure | `_is_user_raised_cancelled()==True` ( `cancelling()==0`) → `NodeCancelledError` (`libs/langgraph/langgraph/pregel/_retry.py:315-335`, `libs/langgraph/langgraph/pregel/_retry.py:790-792`) then panic |
| **Value + Error same task** | **Not valid (exclusive)** | `run_with_retry` clears `task.writes` each attempt (`libs/langgraph/langgraph/pregel/_retry.py:615`, `738`); `commit` branches `if exception:` else success — never both |
| **Timeout** | Valid as `NodeTimeoutError` (a failure subtype) | `_arun_with_timeout` raises `NodeTimeoutError(kind="idle"/"run")` (`libs/langgraph/langgraph/pregel/_retry.py:496-501`, `libs/langgraph/langgraph/errors.py:190-241`) |
| **Runtime fault (recursion, drain, invalid update)** | Valid | `loop.status == "out_of_steps"` → `GraphRecursionError` (`libs/langgraph/langgraph/pregel/main.py:3017-3026`); `loop.status == "draining"` → `GraphDrained` (`libs/langgraph/langgraph/pregel/main.py:3027-3030`); `InvalidUpdateError` from channel validation |

No ambiguous `value+error` tuple is ever returned; the channel write set for a task is exclusive.

### 3. What happens when work returns both a value and an error?
This cannot happen by construction.

- Per-attempt atomicity: `task.writes.clear()` at attempt start (`libs/langgraph/langgraph/pregel/_retry.py:615`, `738`). On success, writers populate channel writes + optional `RETURN` (`libs/langgraph/langgraph/pregel/_write.py` via `ChannelWrite`). On exception, `commit` appends only `(ERROR, exc)` (and `ERROR_SOURCE_NODE` if routed) (`libs/langgraph/langgraph/pregel/_runner.py:596-603`). `map_output_updates` even filters out tasks whose first write is `ERROR`/`INTERRUPT` (`libs/langgraph/langgraph/pregel/_io.py:128-129`).
- If a node tried to both `return {"x": 1}` and `raise`, Python’s control flow makes it one or the other; the `try: task.proc.invoke(...)` in `run_with_retry` either returns (success path) or raises (error path) — never both (`libs/langgraph/langgraph/pregel/_retry.py:617`, `735-752`).
- Tests for handler routing confirm the exclusivity: a handled error does not also emit the failed node’s partial dict; only the handler’s output becomes visible (`libs/langgraph/tests/test_retry.py:2008-2170` graph error handler tests).

If a user wants “partial success,” they must model it explicitly (e.g., return `{"warnings": [...], "status": "partial"}`) — the runtime treats error channels as a disjoint outcome.

### 4. Can a terminal outcome expose a partial or mutable value?
**Partial: No via `invoke`, Yes via `get_state`:**  
`invoke`/`ainvoke` either returns the final `read_channels` snapshot or raises. It never returns a “half-written” superstep because `after_tick` applies writes atomically via `apply_writes` then `updated_channels = apply_writes(...)` and only then `_put_checkpoint` (`libs/langgraph/langgraph/pregel/_loop.py:676-706`). However, `StateSnapshot.values` (`libs/langgraph/langgraph/types.py:643-647`) can be read late via `get_state()` mid-run or after a failure and will reflect whatever channels were durably checkpointed so far (`libs/langgraph/langgraph/pregel/main.py:1391-1433`).

**Mutable: Yes — no freeze/copy.**  
Channels store the *live* object: `LastValue`, `BinaryOperatorAggregate`, `Topic`, `DeltaChannel` all `get()` the stored reference directly (`libs/langgraph/langgraph/channels/base.py:xx`, `libs/langgraph/langgraph/pregel/_io.py:30`). `read_channels` returns that reference (`libs/langgraph/langgraph/pregel/_io.py:45-53`), and `invoke` assigns `latest = payload` without deepcopy (`libs/langgraph/langgraph/pregel/main.py:3957`). Checkpoint persistence does `copy_checkpoint` for the saver (`libs/langgraph/langgraph/pregel/_loop.py:1186`) but the caller’s returned dict aliases the in-memory channel storage. Mutating the returned dict will not corrupt the persisted checkpoint, but it is still observable mutation from the caller’s perspective. No `freeze`, `deepcopy`, or `MappingProxy` is applied; even `DeltaChannel` message-id assignment mutates the original object in place for ID stability (`libs/langgraph/langgraph/pregel/_loop.py:456-458`).

### 5. Do concurrent, repeated, and late waiters observe the same committed facts?
**Concurrent waiters on the same run: Not supported / No evidence.**

- There is no run-handle `Future` that can be `.wait()`-ed by multiple consumers. Each `stream()` call creates its own `SyncQueue`/`AsyncQueue` and `StreamProtocol` (`libs/langgraph/langgraph/pregel/main.py:2763`, `3165`). `FuturesDict` and `BackgroundExecutor` coordinate *internal* node futures, not external run consumers (`libs/langgraph/langgraph/pregel/_runner.py:75-132`, `libs/langgraph/langgraph/pregel/_executor.py:40-121`).
- `test_concurrent_execution_thread_safety` proves isolation, not sharing: 10 threads each call `graph.invoke({"counter":0})` and each gets independent `{"counter":1}` (`libs/langgraph/tests/test_pregel.py:5349-5384`). No test forks two waiters on the same `ainvoke` task.
- Implied behavior: second concurrent `invoke` on same `thread_id` would race on checkpointer `get_tuple`/`put` without isolation (no CAS test found).

**Repeated waiters (idempotent poll after completion): Partially yes via checkpoint.**

- Re-invoking `invoke` with same `config` and `None` input is defined as *resume*, not re-read (`libs/langgraph/langgraph/pregel/_loop.py:852-860` `input is None or Command` → `is_resuming=True`). To re-observe the prior terminal output without re-executing, the user must not call `invoke` but `get_state`/`get_state_history` which re-reads the persisted `CheckpointTuple` (`libs/langgraph/langgraph/pregel/main.py:1427-1433`). That late read is stable (checkpoint id is content-addressed).

**Late waiters: Yes, via checkpointer.**

- `get_state` returns a `StateSnapshot` whose `values` is `read_channels(channels_from_checkpoint(...))` and whose `tasks`/`interrupts` are derived from `pending_writes` (`libs/langgraph/langgraph/pregel/main.py:1144-1265`). Because `_checkpointer_put_after_previous` ensures the checkpoint is ordered after its delta writes (`libs/langgraph/langgraph/pregel/_loop.py:1515-1524`), any waiter that loads via `checkpointer.get_tuple` sees a consistent snapshot.

Missing: no test asserts two `AsyncQueue` consumers drain the same `stream()` output.

### 6. Can waiting be cancelled without cancelling the run itself?
**No uniform API; behavior diverges by executor and node type.**

- LangGraph has no explicit `waiter.cancel()` handle (there is no `join()`/`future` handle per run). The only cancellation primitive is `asyncio.Task.cancel()` on the `ainvoke`/`astream` asyncio task, or `RunControl.request_drain()` for cooperative superstep-boundary drain (`libs/langgraph/langgraph/runtime.py:RunControl`, `libs/langgraph/tests/test_runtime.py:582-666`).
- `test_cancel_ainvoke_with_async_node` (`libs/langgraph/tests/test_runtime.py:735-786`) shows cancelling the `ainvoke` task *does* cancel the run: `CancelledError` is delivered at the `await asyncio.sleep` point inside the async node and the graph stops (`"async_node:finished" not in timeline`). The node is on the event-loop thread, so cancellation propagates.
- `test_cancel_ainvoke_with_sync_node` (`libs/langgraph/tests/test_runtime.py:790-820`) shows the opposite for sync nodes in `ainvoke`: the node runs in a thread via `run_in_executor` (`libs/langgraph/langgraph/_internal/_runnable.py:run_in_executor`), the `asyncio.Task` cancel disconnects the `Future` but the thread orphan keeps running to completion. From the caller’s perspective, waiting was cancelled but the run was *not*.
- `RunControl.request_drain("sigterm")` tested in `test_external_drain_concurrent_sync/async` (`libs/langgraph/tests/test_runtime.py:582-666`) is cooperative: it sets `loop.status = "draining"` at the next `tick()` boundary and raises `GraphDrained`, still checkpointing (`libs/langgraph/langgraph/pregel/_loop.py:650-652`, `libs/langgraph/langgraph/pregel/main.py:3027-3030`). That is a *run* cancellation, not a waiter-only cancellation.
- Therefore: for async runs, waiter cancel ≈ run cancel; for sync-in-async runs, waiter cancel detaches but leaves an orphan; for sync `invoke` there is no async cancellation mechanism at all.

No test exercises “subscribe to stream, unsubscribe one consumer while others keep receiving.”

### 7. Who owns values and diagnostics after publication?
- **Values:** Caller owns the aliased object. There is no transfer of ownership, no `Buffer` move, no `freeze()`. The runtime retains its own copy only inside the checkpoint saver (serialized via `serde` / `copy_checkpoint`). The `SyncQueue` holding stream chunks is drained and discarded per call (`libs/langgraph/langgraph/pregel/main.py:2763`, `2989-3015`). Mutating the returned dict does not affect the persisted checkpoint but does affect any other in-process reference that aliasing holds (no defensive copy found in `read_channels` or `_output`).
- **Diagnostics (errors/interrupts):** Owned similarly by checkpointer + caller:
  - In-memory: `task.writes` holds the live `BaseException` instance (`libs/langgraph/langgraph/pregel/_runner.py:597`), `GraphInterrupt.args[0]` holds live `Interrupt` objects (`libs/langgraph/langgraph/errors.py:106`), `NodeError` wraps a reference to that exception (`libs/langgraph/langgraph/errors.py:148-165`). Handler receives it via `CONFIG_KEY_NODE_ERROR` injection.
  - Persisted: `checkpointer.put_writes(..., (ERROR, exc))` serializes if the saver supports it (e.g., msgpack with allowlist via `_serde.build_serde_allowlist` `libs/langgraph/langgraph/func/__init__.py:610-618`); `StateSnapshot.tasks[i].error` and `CheckpointTask.error` are stringified (`"RuntimeError('...')"`) for late readers (`libs/langgraph/langgraph/types.py:180-202`, `libs/langgraph/tests/test_pregel.py:5427`). Cause chain (`__cause__`, `__notes__`) is preserved in-memory via `add_note`, but the persisted string loses type fidelity.
- **Lifetime:** No `close()` required. Values are GC’d with the graph instance; checkpoint values respect the saver’s TTL/GC (e.g., `InMemorySaver`, `AsyncPostgresSaver`).

### 8. Does waiter release occur only after the complete outcome is visible?
Gate depends on `Durability`:

- **`durability="sync"` — Yes, strongly gated.** `stream()`’s `SyncPregelLoop` waits `loop._put_checkpoint_fut.result()` before exiting the `with` block (`libs/langgraph/langgraph/pregel/main.py:3002-3003`), which itself had waited on `_delta_write_futs` and `prev` (`libs/langgraph/langgraph/pregel/_loop.py:1515-1524`). The subsequent `yield from _output(stream.get, ...)` drains remaining chunks before `run_manager.on_chain_end(loop.output)` (`libs/langgraph/langgraph/pregel/main.py:3006-3032`). The caller’s `invoke` returns only after this.
- **`durability="async"` — No, waiter may be released before durability.** `BackgroundExecutor.submit` for `put` is fire-and-forget with ordering via `prev` but not awaited before return except for the final `_suppress_interrupt` exit-mode flush (`libs/langgraph/langgraph/pregel/_loop.py:1301-1311` only for `exit` mode). `loop._put_checkpoint_fut` is not awaited. A concurrent reader calling `get_state` immediately after `invoke` returns may see stale checkpoint (the saver’s `put` is still in the thread pool). `_checkpointer_put_after_previous` comment explicitly calls this out: "if there's a previous checkpoint save in progress, wait for it ensuring checkpointers receive checkpoints in order" — ordering, not immediacy.
- **`durability="exit"` — Batched visibility.** Intermediate `after_tick` checkpoints are skipped (`_put_checkpoint` early returns if `durability=="exit"` unless `exiting` true `libs/langgraph/langgraph/pregel/_loop.py:1116-1117`); `_exit_delta_writes` are staged and `put` happens in `AsyncExitStack.__aexit__`/`ExitStack.__exit__` after `loop.output` is set (`libs/langgraph/langgraph/pregel/_loop.py:1301-1311`, `libs/langgraph/langgraph/pregel/_executor.py:99-121`). Waiter release (`__exit__` return) happens after that final flush, so the checkpointer becomes visible atomically with output — but intermediate supersteps are invisible.

Stream ordering is separately guaranteed by `chain_future` so `commit`/`_emit("updates")` happens-before the task future resolves (`libs/langgraph/langgraph/pregel/_runner.py:783-786`), ensuring a consumer sees `updates` before the dependent future’s result. However checkpoint visibility is orthogonal.

## Architectural Decisions

| Decision | Rationale / Evidence | Consequence |
|----------|----------------------|-------------|
| **BSP superstep + immutable-during-step channels** (`libs/langgraph/langgraph/pregel/main.py:446-508`, `2974-2978`) | BSP gives deterministic ordering without distributed locks; channel values frozen for duration of step simplifies concurrent node reasoning | Terminal outcome is the fixed point when no triggers remain; partial updates are never exposed mid-step |
| **Single `loop.output` + `GraphOutput` wrapper instead of per-node `Future`** (`libs/langgraph/langgraph/pregel/_loop.py:1350-1355`, `libs/langgraph/langgraph/types.py:368-378`) | Treats the whole graph as the unit of execution (agent-centric) rather than a task DAG of futures | No fan-out/composition via `await task`; composition is via `Command`/`Send`/`call()` PUSH machinery (`libs/langgraph/langgraph/pregel/_runner.py:700-786`) |
| **`call()` returns `SyncAsyncFuture` that chains via `chain_future`** (`libs/langgraph/langgraph/pregel/_call.py:253-298`, `libs/langgraph/langgraph/pregel/_runner.py:783-786`) | Lets functional tasks `await` subtasks while retaining retry/timeout/checkpoint semantics and stream ordering | Still not a public run-handle; futures are scoped to the parent task’s scratchpad and cleaned up by `PregelRunner` |
| **Durability modes (`sync`/`async`/`exit`)** (`libs/langgraph/langgraph/types.py:87-93`, `libs/langgraph/langgraph/pregel/main.py:3878-3888`) | Lets latency vs safety be tuned; `exit` avoids per-step `put` cost for high-throughput agents | Visibility guarantee weakens as above; `async`/`exit` leave a window where `invoke` return is not yet readable via `get_state` |
| **`_handled_exception_ids` + `SKIP_RERAISE_SET` suppression for error handlers** (`libs/langgraph/langgraph/pregel/_runner.py:70-72`, `170-175`, `301-304`) | Avoids double panic when a handled error’s future still fails; error-handler writes marked durable before handler starts (`libs/langgraph/langgraph/pregel/_loop.py:1549-1584`) | Handler becomes the *new* terminal path; original failure stringified in checkpoint but not re-raised |
| **Copy-on-persist, alias-on-return** (`libs/langgraph/langgraph/pregel/_loop.py:1186` `copy_checkpoint`, `libs/langgraph/langgraph/pregel/_io.py:30`) | Cheapest path: caller gets zero-copy view, saver gets isolated snapshot | Mutability hazard noted in Q4; no `freeze()` contract |
| **`FuturesDict` + `_should_stop_others` + `_panic_or_proceed`** (`libs/langgraph/langgraph/pregel/_runner.py:75-132`, `616-634`, `649-697`) | Gives sibling cancellation on first peer failure without complex distributed consensus | Waiter sees exactly one aggregated `GraphInterrupt` or re-raised exception; inflight futures cancelled |

## Notable Patterns

- **Exclusive `ERROR`/`INTERRUPT` vs value writes:** `commit` if/elif chain and `map_output_updates` filter enforce that a task either contributes channel updates or signals a control channel, never both (`libs/langgraph/langgraph/pregel/_runner.py:579-613`, `libs/langgraph/langgraph/pregel/_io.py:128-129`).
- **Drain-then-persist:** `_delta_write_futs` drained into local list under lock then `wait`/`gather` before each `checkpointer.put` (`libs/langgraph/langgraph/pregel/_loop.py:1515-1517`, `1767-1771`). The `prev` future chain ensures strict checkpoint order even when `submit` is concurrent.
- **Two-level suppress:** `SyncPregelLoop._suppress_interrupt` (`libs/langgraph/langgraph/pregel/_loop.py:1313-1352`) and `BackgroundExecutor.__exit__`/`AsyncBackgroundExecutor.__aexit__` (`libs/langgraph/langgraph/pregel/_executor.py:93-211`) both swallow `GraphBubbleUp` so interrupts never surface as task failures.
- **Heartbeat-gated idle timeout:** `_TimedAttemptScope` wraps `config[CONFIG_KEY_SEND/STREAM/CALL/RUNTIME]` and `wait_for_idle_timeout` (`libs/langgraph/langgraph/pregel/_retry.py:127-223`) provides cooperative `idle_timeout` vs wall-clock `run_timeout` distinction.
- **Versioned stream coercion:** `version="v2"` remaps `__interrupt__` out of the value dict into `GraphOutput.interrupts` (`libs/langgraph/langgraph/pregel/main.py:3962-3969`, `4236-4243`), while v1 retains dict merge for backward compat.

## Tradeoffs

| Tradeoff | Pro | Con |
|----------|-----|-----|
| **Graph-level terminal outcome only** (no per-node Future handle) | Simple mental model; checkpoint is the single source of truth | Cannot `await` a specific node without coupling via `call()`; cannot fan-out multiple waiters on same run without external checkpointer polling |
| **Mutable alias return** | Zero-copy, fast for large state dicts/messages | Caller can corrupt subsequent reads if they mutate the returned dict; no `copy.deepcopy` protection |
| **`durability` knob** | Lets high-throughput agents use `exit` for minimal I/O | Weakens publication contract: `async`/`exit` expose a window where `invoke`’s return is not yet durable; late readers must know to wait or use `sync` |
| **Bulk-synchronous checkpoint per superstep** | All tasks in a step commit before next tick; strong step isolation | Cost of channel snapshot per step;DeltaChannel snapshot deferred via `delta_channels_to_snapshot` counter to amortize but still per-step versions |
| **Error-handler suppression via weak set** | Prevents panic on recovered runs; handler runs as a fresh task with its own retry/timeout | Original exception type lost to string in persisted snapshot; debugging must reconstruct from in-memory `NodeError` |
| **Queue-per-call streaming** | Backpressure via `SyncQueue._count`/`AsyncQueue`; `wait()` rescheduled on each drain (`libs/langgraph/langgraph/pregel/main.py:2951-2972`) | No natural “late join” to an in-flight stream; new `stream()` call starts a new loop execution, not a replay |

## Failure Modes / Edge Cases

- **Double `END` writes / conflicting channel updates:** `InvalidUpdateError` raised during `apply_writes` if two nodes write incompatible values to the same `LastValue` channel in one superstep (covered in channel tests `libs/langgraph/tests/test_channels.py`).
- **Recursion limit exhausted:** `loop.status = "out_of_steps"` raises `GraphRecursionError` after `invoke` has already consumed stream chunks (`libs/langgraph/langgraph/pregel/main.py:3017-3026`). Caller receives an exception *instead* of `latest`; prior `latest` is discarded.
- **Sibling failure cancels peers:** `_panic_or_proceed` cancels inflight futures (`libs/langgraph/langgraph/pregel/_runner.py:678-680`). If a peer was mid-checkpoint `put_writes`, its `_delta_write_futs` may still be pending; next checkpoint waits on it, so cancellation does not produce a torn checkpoint.
- **`durability="async"` read-after-write race:** `invoke` may return before `checkpointer.get_tuple(config)` reflects the new checkpoint. A following `get_state` (late waiter) may see the previous checkpoint. The mitigation is `durability="sync"` which `await`s the future (`libs/langgraph/langgraph/pregel/main.py:3002-3003`).
- **Time-travel resume vs interrupt resume:** `_first` distinguishes `is_time_traveling` (`CONFIG_KEY_CHECKPOINT_ID` explicitly set) from normal resume to drop stale `RESUME` writes and clear `INTERRUPT` markers (`libs/langgraph/langgraph/pregel/_loop.py:862-890`). If the caller confuses the two, handler may not re-fire.
- **Cancelled `ainvoke` with sync node orphan:** The background thread continues, may still call `put_writes`/`put` after the `AsyncBackgroundExecutor` was cancelled, potentially persisting state for a run the caller considers cancelled (`libs/langgraph/tests/test_runtime.py:790-820`, `libs/langgraph/langgraph/pregel/_executor.py:196-200` sentinel cancellation).
- **Partial checkpoint on `GraphDrained`:** `draining` status still yields `loop.output = read_channels(...)` (`libs/langgraph/langgraph/pregel/_loop.py:1355`) and raises `GraphDrained` to caller, but checkpoint is already staged. Resumption via `invoke(None, config)` is valid (`libs/langgraph/langgraph/pregel/_loop.py:649-652`).
- **Mutable returned state used as next input:** If caller mutates the returned dict and passes it again as `input`, channel applying will treat it as new writes, not as the persisted snapshot — potential double-counting for `operator.add` reducers.

## Future Considerations

- Expose an explicit **run-handle Future** (`RunHandle[T]`) that multiple waiters can `await`/`share()` and that maps to `Future[GraphOutput]` semantics with `cancel(waiter_only=True)` vs `cancel(run=True)` flags, closing the gap in Q6.
- **Freeze/copy on return:** Wrap the terminal `read_channels` result in a `MappingProxy` or `deepcopy` (or `pydantic` `model_copy`) behind an opt-in flag, or at least document that the returned `GraphOutput.value` aliases in-memory channel storage.
- **Stable publication waiter API:** Add `await checkpointer.wait_for_checkpoint(config, checkpoint_id=loop.checkpoint["id"])` so late readers can block until `durability="async"`’s pending `put` drains, formalizing the visibility invariant currently hidden in `_checkpointer_put_after_previous`.
- **Type-preserving error persistence:** Store `error.__class__`/`traceback` alongside stringified message in `CheckpointTask.error` (e.g., `{"type":"RuntimeError","message":"...","notes":[...]}`) so `get_state_history` diagnostics do not lose fidelity.
- **Per-waiter stream replay:** Allow `stream(resume_from=checkpoint_id)` that replays `SyncQueue` history to late consumers without restarting the loop, enabling legitimate multi-waiter fan-out.
- **Cancellation semantics split:** Make `ainvoke`’s sync-node orphan behavior explicit — either propagate a cooperative `RunControl` flag into the worker thread so it can checkpoint and exit cleanly, or document/return a `CancelledError` handle the caller can use to `await thread_future`.

## Questions / Gaps

- No evidence of a **public `wait()`/`join()` contract** for exhaustive fan-out: searches for `waiter`, `join`, `Future[GraphOutput]` in `libs/langgraph/langgraph/pregel/*.py` and `libs/langgraph/tests/*.py` return only internal `get_waiter`/`FuturesDict` machinery (`libs/langgraph/langgraph/pregel/_runner.py:183-257`). Confirm intention is checkpointer-polling rather than future-sharing.
- **Repeated idempotent `invoke`:** Is there a guarantee that `invoke(input, config)` called twice with identical input and same `thread_id` is deduped via cache/checkpoint? `CachePolicy` tests (`libs/langgraph/tests/test_retry.py`) show per-node caching but not graph-level idempotency.
- **Late-waiter stability test:** No test asserts `concurrent.futures.wait` on two threads both calling `graph.invoke` on the same `thread_id` observe the same `checkpoint_id`. The only concurrency test is independent threads (`libs/langgraph/tests/test_pregel.py:5349-5385`).
- **Value + error via `Command(update=..., resume=...)`?** `map_command` yields both `RESUME` and update writes together (`libs/langgraph/langgraph/pregel/_io.py:56-78`). Does a `Command` that both resumes an interrupt and updates state count as a valid `value+interrupt-resume` combination at the terminal level? Search boundary: looked at `libs/langgraph/tests/test_interruption.py` — no explicit `Command` value+error case.
- **Ownership of `DeltaChannel` sub-snapshot values:** `counters_since_delta_snapshot` bookkeeping (`libs/langgraph/langgraph/pregel/_loop.py:1094-1142`) determines when delta history is folded into the checkpoint. The exact point where the in-memory delta list becomes “immutable and safe to publish” for concurrent subgraph walks is not directly tested.

---

Generated by `dimension 01.16: Result and Outcome Publication Contract` against `langgraph`.
