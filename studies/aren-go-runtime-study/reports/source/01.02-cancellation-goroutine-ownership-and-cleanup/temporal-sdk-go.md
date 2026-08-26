# Source Analysis: temporal-sdk-go

## 01.02 Cancellation, Goroutine Ownership, and Cleanup

### Source Info

| Field | Value |
|-------|-------|
| Name | temporal-sdk-go |
| Path | `studies/aren-go-runtime-study/sources/temporal-sdk-go` |
| Language / Stack | Go (gRPC, protobuf; `golang.org/x/time/rate` limiters, `x/sync`-style WaitGroup patterns, Go 1.24 `testing/synctest` in tests) |
| Analyzed | 2026-08-26 |

## Summary

The Temporal Go SDK implements a layered ownership tree for cancellation and shutdown. A user-facing `AggregatedWorker` owns per-task-type workers (`workflowWorker`, `activityWorker`, `sessionWorker`, `nexusWorker`), each of which wraps a `baseWorker` that owns poller goroutines, a task dispatcher, an eager dispatcher, and slot-reservation goroutines — all tracked on two WaitGroups (`stopWG`, `pollerWG`) and stopped by closing one `stopCh` (`internal/internal_worker_base.go:226-256`). Shutdown is staged top-down: `AggregatedWorker.Stop()` sets `noRepoll`, closes its own `stopC`, sends a best-effort `ShutdownWorker` RPC, then stops child workers sequentially (`internal/internal_worker.go:1614-1676`); the workflow worker additionally closes a *second* stop channel only after the main workflow worker drains so pending local activities can finish (`internal/internal_worker.go:389-397`, `466-472`).

Cancellation is cooperative and cause-typed. Activity contexts derive from a user-suppliable root context wrapped with `context.WithCancelCause` at worker construction (`internal/internal_worker.go:2291-2295`), are registered by task token in an `activityCancellationCallbacks` map (`internal/internal_task_handlers.go:2126-2151`), and are cancelled from three sources: server cancel delivered through heartbeat responses, worker commands (pause/reset) delivered over a control queue, and worker shutdown published as `ErrWorkerShutdown` cause after the graceful-wait window (`internal/internal_worker_base.go:921-923`). The activity handler distinguishes these causes explicitly via `context.Cause` (`internal/internal_task_handlers.go:2480-2485`). Workflow coroutines get their own `WithCancel` root whose cancel handler is idempotent by contract (`internal/internal_workflow.go:566-574`), and dispatcher teardown deliberately exits stuck coroutines in an unwaited goroutine (`internal/internal_workflow.go:1365-1381`).

Shutdown boundedness is engineered rather than incidental: long polls are capped at 70 s (`internal/internal_task_pollers.go:37-39`), drain-on-shutdown mode keeps already-polled tasks flowing to the dispatcher while it discards late arrivals after rate-limit failure (`internal/internal_worker_base.go:729-749`, `835-847`), completion RPCs run on detached background contexts so results still publish during stop (`internal/internal_task_pollers.go:1520`), and the whole wait is bounded by `WorkerStopTimeout` (`default 0s`, `internal/worker.go:291-294`). The behavior is backed by an unusually strong deterministic test battery using `testing/synctest` plus `goleak` sweeps (`internal/internal_worker_base_test.go:665-1127`, `test/integration_test.go:109-134`). Weaknesses: `AggregatedWorker.Stop()` is not safe against truly concurrent callers (check-then-close without a lock can panic on double channel close), several goroutines are intentionally detached (bounded but unwaited), and the local-activity tunnel uses a 100,000-slot buffered channel.

## Rating

**8 / 10.** Rubric alignment: full marks for context propagation depth (cancel causes threaded from server events to user code and back), ownership-tree clarity (documented `stopC` creation/closure contracts on nearly every struct), bounded shutdown (timeouts at every layer), and leak/shutdown test coverage (`synctest`-based race tests, goleak in CI suites). Deductions for: (1) non-atomic idempotency of the top-level `Stop()` — the exact question "how is shutdown made idempotent under concurrent callers" is answered best-effort, not atomically; (2) a handful of detached goroutines whose lifetime is only implicitly bounded (local-activity execution after early return, legacy-path poll goroutines, dispatcher coroutine-exit loop); (3) effectively unbounded buffering in the LA tunnel; (4) default `WorkerStopTimeout = 0` silently collapsing the graceful-drain window that `activity.GetWorkerStopChannel`'s documentation promises.

## Evidence Collected

Every entry cites file paths relative to the source root `studies/aren-go-runtime-study/sources/temporal-sdk-go/`.

| Area | Evidence | File:Line |
|------|----------|-----------|
| baseWorker ownership fields | `stopCh` created in `newBaseWorker`, closed by `Stop()`; `stopWG`/`pollerWG`; two limiter contexts with distinct lifecycles (task-limiter stays live during drain) | internal/internal_worker_base.go:226-239 |
| Poller goroutine starts | One goroutine per fixed poller or per autoscaling runner, each registered on both WaitGroups | internal/internal_worker_base.go:447-458 |
| Dispatcher/eager dispatcher goroutines + queue-close ordering | Closer goroutine waits for `pollerWG` then closes `taskQueueCh` so the dispatcher can drain deterministically | internal/internal_worker_base.go:467-478 |
| Drain-mode dispatcher | Ranges `taskQueueCh` until close; releases slots and discards remaining tasks once limiter is cancelled | internal/internal_worker_base.go:729-749 |
| Legacy dispatcher | Selects on `stopCh` vs `taskQueueCh`, releases slot when stopping | internal/internal_worker_base.go:751-769 |
| Eager dispatcher drains on stop | Drains buffered eager tasks before exiting on `stopCh` | internal/internal_worker_base.go:778-793 |
| Slot reservation tracked | `reserveSlotAsync` runs under `stopWG.Go`, releases reserved permit if receiver context dies | internal/internal_worker_base.go:623-648 |
| Poll send/close race handling | In normal mode send races against `stopCh`; in drain mode blocking send relies on documented dispatcher liveness invariant | internal/internal_worker_base.go:835-847 |
| `baseWorker.Stop()` sequence | Close `stopCh` → cancel limiter ctx → close session token bucket → bounded `awaitWaitGroup(stopWG)` → cancel task-limiter ctx → `backgroundContextCancel(ErrWorkerShutdown)` | internal/internal_worker_base.go:899-926 |
| Bounded WaitGroup wait helper | `awaitWaitGroup` spawns a waiter goroutine; leaks it past deadline until WG drains | internal/internal_utils.go:161-181 |
| AggregatedWorker shutdown ordering | `stopPolling()` first (prevents re-poll race vs ShutdownWorker RPC), then `close(stopC)`, then RPC, then child workers in order workflow→activity→session→nexus | internal/internal_worker.go:1613-1676 |
| Two-stage workflow-worker stop | `close(ww.stopC)` → main worker Stop → `close(localActivityStopC)` → LA worker Stop, so pending local activities finish first | internal/internal_worker.go:465-472 |
| Local activity tunnel channels | `taskCh` buffered 100000, unbuffered `resultCh`, both guarded by `stopCh` in send/receive | internal/internal_task_pollers.go:260-284 |
| doneCh escape hatch for LA results | LA worker send on `laResultCh` selects against workflow-task `doneCh`, closed by defer in `processWorkflowTask` | internal/internal_task_pollers.go:956-967; internal/internal_task_pollers.go:461-465 |
| Long-poll bound & drain behavior | `pollTaskServiceTimeOut = 70s`; drain mode blocks on poll completion instead of cancelling | internal/internal_task_pollers.go:37-39, 305-343 |
| Detached LA execution goroutine | If ctx.Done wins before `doneCh`, `executeLocalActivityTask` returns while the activity goroutine still runs; result dropped | internal/internal_task_pollers.go:1059-1093 |
| Activity root context wiring | User `BackgroundActivityContext` (or Background) wrapped in `WithCancelCause`; cancel func handed to activity workers | internal/internal_worker.go:2289-2295, 2415-2416 |
| Per-activity cancel registration | `canCtx := WithCancelCause(rootCtx)`; cancel registered/unregistered by task token | internal/internal_task_handlers.go:2389-2399 |
| Server-cancel path via heartbeat | Heartbeat error classes call `i.cancelHandler(err)` (the cancel-cause func) to deliver cancellation into the activity context | internal/internal_task_handlers.go:2277-2302 |
| Worker-command cancel path (pause/reset) | Shared heartbeat manager polls control queue and calls `activityCancellationCallbacks.cancel(taskToken)` on the registered per-token cancel funcs | internal/internal_worker_heartbeat.go:307-335; internal/internal_task_handlers.go:2142-2151 |
| Cause-based discrimination | `isActivityCanceled := ctx.Err() == context.Canceled && IsCanceledError(context.Cause(ctx))`; comment names `ErrWorkerShutdown`/`ErrActivityPaused` as non-server causes | internal/internal_task_handlers.go:2480-2485 |
| Public shutdown-channel contract | `activity.GetWorkerStopChannel` documents close→graceful wait→context-cancel→exit order | internal/activity.go:274-282 |
| Heartbeat batching lifecycle goroutine | Batch flusher exits on `hbBatchEndTimer`, `workerStopChannel`, or `closeCh`; known lost-progress race documented inline | internal/internal_task_handlers.go:2200-2235 |
| Invoker close flushes heartbeats | `Close` closes `closeCh` and optionally flushes last details on a timeout-bounded context | internal/internal_task_handlers.go:2313-2324 |
| Completion publication survives stop | `reportActivityComplete(context.Background(), ...)` — respond RPCs rooted at Background, bounded by gRPC retry policy/default timeout | internal/internal_task_pollers.go:1520, 1536-1568 |
| Workflow coroutine cancellation | `d.cancel()` behind `RegisterCancelHandler`; documented as safe to call repeatedly | internal/internal_workflow.go:566-574 |
| Cancel event delivery | History event `WorkflowExecutionCancelRequested` invokes registered handler | internal/internal_event_handlers.go:1692-1693 |
| Dispatcher close idempotent + detached | `closed` flag guards re-entry; coroutine exit loop runs in an unwaited goroutine because stuck coroutines/blocking defers cannot be waited on | internal/internal_workflow.go:1365-1381 |
| Sticky-cache eviction cleanup | LRU `RemovedFunc` calls `onEviction` from cache's own goroutine; locks context and closes event handler | internal/internal_worker_cache.go:80-85; internal/internal_task_handlers.go:659-672 |
| Session token bucket wakeup | `waitForAvailableToken` parks on cond var; `baseWorker.Stop` closes bucket so parked pollers report closed | internal/session.go:506-523; internal/internal_worker_base.go:906-908 |
| Session creation cleanup ordering | Creation activity completes session on ctx.Done and distinguishes worker-shutdown cancel from server cancel | internal/session.go:424-436 |
| Session completion idempotent | `CompleteSession` checks map membership before closing `doneChan` | internal/session.go:584-592 |
| Heartbeat manager join | `stop()` CAS-idempotent; cancels heartbeat ctx, closes `stopC`, blocks on `stoppedC` | internal/internal_worker_heartbeat.go:371-381 |
| Worker-commands poller joins | `run` defers `heartbeatCancel()` and `<-workerCommandsDone` | internal/internal_worker_heartbeat.go:176-186 |
| Start memoization | `sync.OnceValue` makes duplicate `Start()` calls return consistent result | internal/internal_worker.go:1380-1385, 2665 |
| Fatal-error self-stop | `fatalErrorCallback` records first error under lock, notifies plugin, stops worker unless already stopped | internal/internal_worker.go:2330-2353 |
| Client connection close refcounted | `Close()` decrements `unclosedClients`; closes shared gRPC conn only at zero, then pins counter to max for idempotency | internal/internal_workflow_client.go:1776-1797 |
| Partial-start rollback | Activity-worker start failure stops workflow worker; session-worker start failure stops workflow+activity workers | internal/internal_worker.go:1463-1501 |
| Partial-start rollback (session worker) | Mirrored rollback inside `sessionWorker.Start` itself | internal/internal_worker.go:566-578 |
| Deterministic shutdown tests | Task-not-dropped, autoscaling-not-dropped, drain respects dispatch rate, drain-after-cancel release, legacy no-processing, stop-timeout bound, prompt legacy return, session-token park | internal/internal_worker_base_test.go:665-1127, 1367 |
| Slot-supplier cancellation tests | Reserve honors cancelled contexts including previously-missing retry-loop cancellation point | internal/resource_tuner_cancel_test.go:33-100 |
| Leak detection | goleak sweep with retry window ignoring known coroutine yield frame; dedicated child-workflow leak test | test/integration_test.go:109-134; test/child_workflow_leak_test.go:24-42 |
| Integration shutdown tests | Activity/LA cancel from worker shutdown, LA no-heartbeat-after-stop, complete-within-graceful-shutdown, active timer+activity workflows during shutdown, session shutdown w/ poll-complete | test/integration_test.go:8316-8600, 2234 |
| Interrupt helper goroutine | `InterruptCh` spawns a signal-forwarding goroutine that never exits if no signal arrives | internal/internal_utils.go:183-196 |

## Answers to Dimension Questions

**Who is responsible for waiting for every goroutine started on behalf of a run?**
The `baseWorker` is the single join point: every poller, autoscaling poll goroutine, autoscaler ticker loop, dispatcher, eager dispatcher, task processor, and slot reservation is counted on `bw.stopWG` (`internal/internal_worker_base.go:449-478`, `611-619`, `628`, `694-695`), and `baseWorker.Stop()` waits on it under `options.stopTimeout` (`internal/internal_worker_base.go:910-918`). A secondary `pollerWG` exists purely to know when all pollers have exited so `taskQueueCh` can be closed, letting the dispatcher's range loop terminate (`internal/internal_worker_base.go:467-472`). Exceptions that are *not* awaited: the local-activity execution goroutine spawned per LA task when context-done wins the select (`internal/internal_task_pollers.go:1059-1083`), the legacy-path poll RPC goroutine abandoned after `doPoll` returns `errStop` (`internal/internal_task_pollers.go:320-342`), the heartbeat-batching goroutine (exits on timer/stopCh/closeCh, `internal/internal_task_handlers.go:2207-2234`), and the dispatcher's coroutine-exit goroutine (`internal/internal_workflow.go:1376-1380`). Each is bounded by some external condition (activity deadline, gRPC timeout, channel close), but none is joined.

**What does the API claim when cancellation is requested but work has not stopped?**
Three layers of claims. (1) Public doc: after `stopC` closes, the worker "will wait until the worker stop timeout finishes. After the timeout hits, the worker will cancel the activity context and then exit" (`internal/activity.go:274-277`) — i.e., between request and cancel there is a grace window owned by `WorkerStopTimeout` (default `0s`, `internal/worker.go:291-294`). (2) Cause truthfulness: activities can distinguish *why* they were cancelled — server cancel vs `ErrWorkerShutdown` vs pause/reset — via `IsCanceledError(context.Cause(ctx))` (`internal/internal_task_handlers.go:2480-2485`). (3) Completion semantics: a task polled during shutdown in drain mode is guaranteed dispatched and processed rather than dropped, enforced by tests (`internal/internal_worker_base_test.go:662-736`); activity completion/failure publication intentionally ignores shutdown by using background-rooted RPCs (`internal/internal_task_pollers.go:1520`). If the graceful window expires, `Stop()` logs "Worker graceful stop timed out" and proceeds to cancel anyway — work may still be running when the root context dies.

**Can cleanup block forever or overwrite the primary failure?**
Blocking is bounded everywhere we traced: long polls ≤70 s (`internal/internal_task_pollers.go:39`), respond RPCs ≤ default 10 s gRPC timeout with retry policy (`internal/internal_utils.go:110-141`), `sendShutdownWorkerRPC` uses the same builder and errors are logged, never propagated (`internal/internal_worker.go:1692-1762`), session token parking ends on bucket close, and the final WaitGroup wait is capped by `WorkerStopTimeout`. Two soft spots: `awaitWaitGroup`'s waiter goroutine lives past its deadline until the WG actually drains (`internal/internal_utils.go:164-181`), and heartbeat flush inside `temporalInvoker.Close` runs synchronously on the caller's ctx with `recordTimeout` bounds (`internal/internal_task_handlers.go:2313-2324`). Overwriting the primary failure is actively prevented: the fatal-error callback stores only the first error under `fatalErrLock` (`internal/internal_worker.go:2332-2353`), `Run` returns that stored error rather than a shutdown artifact (`internal/internal_worker.go:1604-1609`), and a heartbeat visitor failure is carried as the *cause* of context cancellation so the handler skips double-reporting a response (`internal/internal_task_handlers.go:2153-2159`, `2473-2478`).

**How is shutdown made idempotent under concurrent callers?**
Unevenly. Idempotent pieces: `dispatcherImpl.Close` uses a `closed` flag (`internal/internal_workflow.go:1365-1372`); `sharedNamespaceWorker.stop` uses `started.CompareAndSwap` (`internal/internal_worker_heartbeat.go:371-374`); session `CompleteSession` checks map membership before closing the channel (`internal/session.go:584-592`); client `Close` pins its refcount to max after first decrement (`internal/internal_workflow_client.go:1781-1790`); `Start` is memoized with `sync.OnceValue` (`internal/internal_worker.go:2665`). But the top-level `AggregatedWorker.Stop()` guards `close(aw.stopC)` with only a non-atomic select check ("Only attempt stop if we haven't attempted before", `internal/internal_worker.go:1614-1640`) — no mutex or Once — so two simultaneous callers can both pass the check and the second `close` panics. The same check-then-act pattern appears in `fatalErrorCallback` (`internal/internal_worker.go:2346-2351`). No test exercises concurrent `Stop()` (searches for concurrent/double-stop tests across `internal/*_test.go` and `test/*.go` found nothing).

## Architectural Decisions

1. **One owner, two WaitGroups, one close-ordering goroutine.** All runtime goroutines hang off `baseWorker`; `pollerWG` feeds a closer goroutine that closes `taskQueueCh` so the dispatcher can distinguish "no more tasks" from "blocked" and drain cleanly (`internal/internal_worker_base.go:467-478`). This converts shutdown from a race into a phase transition.
2. **Channel-close as the universal stop primitive, with documented creation/closure sites.** Nearly every struct comments where its `stopC`/`stopCh` is created and closed (e.g., `internal/internal_worker.go:91-98`, `110-112`, `1279-1282`; `internal/internal_nexus_worker.go:22-24`; `internal/internal_worker_heartbeat.go:164-169`). Closure order encodes semantics: workflow worker fully stops before `localActivityStopC` closes so pending LAs complete (`internal/internal_worker.go:389-397`).
3. **Cause-carrying cancellation.** `context.WithCancelCause` at both worker root (`internal/internal_worker.go:2295`) and per-task level (`internal/internal_task_handlers.go:2394`) lets one context serve server cancels, shutdown, pause, and reset without sentinel channels.
4. **Two shutdown modes selected by server capability.** `workerPollCompleteOnShutdown` switches pollers/dispatcher between legacy immediate-abandon and drain-to-completion behavior, including keeping the task dispatch limiter alive during drain (`internal/internal_worker_base.go:235-239`, `326-333`, `835-847`).
5. **Detached-but-bounded teardown for untrustworthy code.** Coroutine exit on dispatcher close runs unwaited because stuck coroutines and blocking defers would deadlock the stopper (`internal/internal_workflow.go:1373-1380`); symmetrically, LA execution is abandoned when its context fires first (`internal/internal_task_pollers.go:1061-1073`).
6. **Result publication decoupled from lifecycle context.** Respond-completion/failed/canceled RPCs use `context.Background()` roots so shutdown never eats a finished result (`internal/internal_task_pollers.go:1520`, `1536-1568`), mirroring the outbound storage visitor using `backgroundContext` so a cancelled activity can still upload payloads (`internal/internal_task_handlers.go:2530-2537`).

## Notable Patterns

- **Registry/callback map keyed by task token** maps async server decisions onto live contexts: `register`/`delete` returns an unregister closure (`internal/internal_task_handlers.go:2130-2139`), reused by both the activity worker and the shared namespace heartbeat worker (`internal/internal_worker_heartbeat.go:76`).
- **Escape-hatch selects on sends**: LA result delivery picks between `laResultCh` and the workflow task's `doneCh` (`internal/internal_task_pollers.go:960-966`); poll→queue send picks between `taskQueueCh` and `stopCh` outside drain mode (`internal/internal_worker_base.go:841-846`).
- **Cond-var bucket with broadcast close** so blocked pollers wake on shutdown and return "closed" instead of waiting for sessions (`internal/session.go:506-523`).
- **Virtual-time deterministic tests**: `synctest.Test` verifies Stop timing exactly (e.g., Stop returns after precisely `stopTimeout`, not the 70 s poll bound, `internal/internal_worker_base_test.go:1084-1091`) and that no task is dropped mid-shutdown (`internal/internal_worker_base_test.go:722-734`).
- **Goleak with targeted ignore**: the suite tolerates the known parked-coroutine frame while failing on anything else, retried up to a minute (`test/integration_test.go:114-133`).
- **Suppressed-error logging near shutdown**: poll errors are not logged after `stopCh` closes and server-graceful-stop Unavailable errors are filtered for internal workers (`internal/internal_worker_base.go:851-884`).

## Tradeoffs

- **Grace period vs. immediacy.** Default `WorkerStopTimeout=0` means `awaitWaitGroup` effectively doesn't wait, so the documented graceful window collapses and `ErrWorkerShutdown` cancellation lands almost immediately after `stopC` closes (`internal/worker.go:291-294`; `internal/internal_worker_base.go:913-923`). Users must opt in to grace; the SDK chooses fast, loud shutdown by default.
- **Drain vs. responsiveness.** Drain mode blocks `doPoll` on the in-flight long poll (up to 70 s) to let the server complete polls gracefully (`internal/internal_task_pollers.go:326-332`), trading slower `Stop()` for no dropped tasks and no canceled-poll noise; legacy mode returns instantly but abandons the poll goroutine.
- **Detached coroutine exit vs. data integrity.** Exiting stuck coroutines unwaited guarantees the dispatcher closes but means deferred user cleanup in workflows may interleave arbitrarily with eviction (`internal/internal_workflow.go:1373-1380`).
- **Big buffer vs. backpressure.** `laTunnel.taskCh` (100k) plus `pushEagerTask`'s "should always be non-blocking" assumption (`internal/internal_task_pollers.go:262`; `internal/internal_worker_base.go:685-688`) favor throughput and simplicity over strict backpressure.
- **Refcount-by-atomics client close** avoids mutex cost and double-close but makes `Close` ordering-sensitive: ignored when siblings remain open (`internal/internal_workflow_client.go:1776-1797`).

## Failure Modes / Edge Cases

- **Concurrent `Stop()` panic window**: double `close(aw.stopC)` possible; no lock/Once (`internal/internal_worker.go:1616-1640`). Also reachable indirectly if `fatalErrorCallback` races a user-initiated stop.
- **Heartbeat batching progress-loss race**: acknowledged in-code — user `Heartbeat` can grab the lock before the batch-flusher goroutine, dropping the latest details (`internal/internal_task_handlers.go:2227-2232`).
- **Poller/ShutdownWorker RPC race**: a natural long-poll completion just after `stopC` closes could trigger a re-poll before `ShutdownWorker` is sent; mitigated by calling `stopPolling()` (`noRepoll`) *before* closing `stopC` (`internal/internal_worker.go:1621-1638`; `internal/internal_worker_base.go:510-512`).
- **Slot leak on reservation handoff**: `reserveSlotAsync` releases the permit if the receiver ctx died between reserve and handoff (`internal/internal_worker_base.go:642-646`); pollTask releases any unused permit via defer (`internal/internal_worker_base.go:799-803`).
- **Dispatcher discard under drain**: once the dispatch limiter is cancelled mid-drain, remaining queued tasks are released-and-dropped, so drain is best-effort, not a guarantee (`internal/internal_worker_base.go:735-744`).
- **LA result orphaning**: workflow task completing/timing out while LA still runs leaves the goroutine executing with a dead workflow; result later dropped via `doneCh` branch (`internal/internal_task_pollers.go:960-966`).
- **Sticky-cache forced eviction mid-flight**: LRU eviction of a live workflow destroys its coroutines from the cache goroutine and counts the metric (`internal/internal_task_handlers.go:663-671`); recovery relies on replay from history.
- **Session worker partial start**: if the resource-specific activity worker fails to start, `creationWorker.Stop()` rolls back (`internal/internal_worker.go:566-578`).
- **Interrupt helper leak**: `InterruptCh` goroutine persists for process lifetime even if the returned channel is discarded (`internal/internal_utils.go:189-193`).

## Future Considerations

- Make `AggregatedWorker.Stop()` atomic (mutex or `sync.Once`) and add a concurrency test; this is small, contained work in `internal/internal_worker.go:1614-1676`.
- Reconsider default `WorkerStopTimeout=0`: either document loudly next to `activity.GetWorkerStopChannel` or give drain-mode workers an implicit floor so the promised grace window exists by default.
- Bound or track the detached LA execution goroutine (e.g., register on a per-execution WaitGroup drained by the LA worker stop) to make "who waits" answer unconditional.
- Right-size `localActivityTunnel.taskCh` (100k) or convert to a semaphore-backed pending set to restore backpressure.
- Extend `synctest` coverage to the aggregated layer (ordering of `ShutdownWorker` RPC vs poller exit, session worker rollback) since current virtual-time tests stop at `baseWorker`.

## Questions / Gaps

- **No evidence found for concurrent-Stop testing or protection** beyond the racy select; searched `internal/*_test.go` and `test/*.go` for concurrent/double-stop scenarios (patterns: `concurrent.*Stop`, `Stop.*twice`, `double.*[Ss]top`) — nothing.
- **No evidence found** that `awaitWaitGroup`'s post-timeout waiter goroutine is ever reclaimed by design intent; it exits only when the WaitGroup eventually drains. Harmless but undocumented.
- The Aren requirement docs referenced by the dimension (`Aren/docs/phase-1-prd/02-lifecycle-requirements.md`, phased roadmap) were **not read**: source-isolation rules restrict analysis to the selected source directory, so SDK behavior was mapped without cross-checking those requirement texts.
- Whether the heartbeat-flush race (`internal/internal_task_handlers.go:2227-2232`) causes observable workflow-visible detail loss in production was not verifiable from code alone; no regression test pins the current behavior.
- Eager-dispatch path during shutdown drain (`pushEagerTask` while dispatcher draining) appears safe due to the 2000-buffer channel (`internal/internal_worker_base.go:392-394`), but no test exercises eager arrival during drain; behavior inferred, not demonstrated.

---

Generated by `01.02 Cancellation, Goroutine Ownership, and Cleanup` against `temporal-sdk-go`.
