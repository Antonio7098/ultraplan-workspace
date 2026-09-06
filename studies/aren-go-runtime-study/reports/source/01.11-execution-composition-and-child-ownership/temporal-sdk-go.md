# Source Analysis: temporal-sdk-go

## 01.11 Execution Composition and Child Ownership

### Source Info

| Field | Value |
|-------|-------|
| Name | temporal-sdk-go |
| Path | `studies/aren-go-runtime-study/sources/temporal-sdk-go` |
| Language / Stack | Go (deterministic workflow SDK, gRPC, protobuf) |
| Analyzed | 2026-08-30 |

## Summary

temporal-sdk-go models child ownership at the durable workflow level, not as an in-process supervisor tree. Parents are Temporal workflows; children are separately persisted executions created via `ExecuteChildWorkflow`. Ownership is expressed through server-enforced `ParentClosePolicy` (`ABANDON`/`TERMINATE`/`REQUEST_CANCEL`) and `WaitForCancellation`, plus client-side deterministic concurrency primitives (`Go`/`GoNamed` coroutines multiplexed on `dispatcherImpl`, `Future`/`Channel`/`Selector`, `WaitGroup`, `Semaphore`). There is no built-in structured-concurrency composition library for fail-fast parallel, collect-all with ordered results, or bounded pools — the SDK provides bases (futures, selectors, mutex/semaphore) and leaves aggregation, cancellation, and error policies to workflow code. Parent terminal semantics are eventually consistent: the parent completes (`Complete` → `syncWorkflowDefinition.Execute` result) and the server applies the close policy to children; local coroutines are torn down via `dispatcher.Close` with `Goexit`. Sibling failure has no implicit propagation; sibling cancellation requires an explicit shared `WithCancel` context. Result ordering is caller-defined (iteration order over futures vs. completion order via `Selector`).

## Rating

**6 / 10**

Rationale: identity, lifecycle, and parent-close semantics are well-defined and durable (protobuf command state machine + server policy). Deterministic scheduling (`dispatcherImpl.ExecuteUntilAllBlocked`) and panic-to-failure conversion preserve executor semantics and observability (stack traces, `workflowPanicError`). However, composition helpers are minimal: no first-class `Group`/`Pool` with fail-fast, collect-all, ordered results, or bounded-parallel abstractions; bounded concurrency relies on manual `Semaphore`; test suite shows patterns but not SDK-provided helpers. Sibling cancellation and result determinism are explicit rather than structured, increasing hand-roll risk.

## Evidence Collected

Every entry MUST include a file path with line numbers. Format: `path/to/file.ts:NN`.

| Area | Evidence | File:Line |
|------|----------|-----------|
| Parent-child identity & options | `ChildWorkflowOptions` fields: `WorkflowID`, `TaskQueue`, `ParentClosePolicy`, `WaitForCancellation`, `RetryPolicy`, `TypedSearchAttributes` | `internal/workflow.go:441-569` |
| WorkflowInfo parent linkage | `WorkflowInfo` parent/root execution pointers, namespace, attempt, run IDs | `internal/workflow.go:1512-1583` |
| Child future dual-result | `childWorkflowFutureImpl` holds `decodeFutureImpl` (result) + `executionFuture` (start) | `internal/internal_workflow.go:242-245` |
| Child execution API | `ExecuteChildWorkflow` checks `ctx.Err`, generates `childWorkflowID`, builds `ExecuteWorkflowParams` with per-child `dataConverter`/`failureConverter` | `internal/workflow.go:1386-1510` |
| Server command & ownership | `ExecuteChildWorkflow` builds `StartChildWorkflowExecutionCommandAttributes` with `ParentClosePolicy`, `RetryPolicy`, `Header`, `Memo` | `internal/internal_event_handlers.go:565-648` |
| Command state machine | `childWorkflowCommandStateMachine` state transitions, `cancel()` moves to `CanceledAfterStarted` and `moveCommandToBack` | `internal/internal_command_state_machine.go:79-83,394-404,712-784` |
| Coroutine dispatcher ownership | `dispatcherImpl` sequences coroutines, `NewCoroutine`, `ExecuteUntilAllBlocked`, `Close` drains with `Goexit` | `internal/internal_workflow.go:163-180,1247-1339,1365-1381` |
| Context & cancellation tree | `cancelCtx` children map, `propagateCancel`, `cancel` closes `done` channel and recursively cancels children | `internal/context.go:181-325` |
| Disconnected context escape hatch | `NewDisconnectedContext` does not propagate parent cancellation (cleanup pattern) | `internal/context.go:197-200` |
| Future composition | `futureImpl.Get` blocks on channel, `Set`/`Chain`/`GetAsync` used by `Selector` | `internal/internal_workflow.go:54-60,359-473` |
| Selector composition | `selectorImpl.Select` multiplexes `Channel` and `Future` cases, deterministic dispatch | `internal/internal_workflow.go:130-145,1433-1573` |
| WaitGroup collect-all | `waitGroupImpl` with `future`/`settable`, `Go`/`Wait`, panic on reuse/negative counter | `internal/internal_workflow.go:62-68` ; `internal/workflow.go:303-321` |
| WaitGroup test - collect-all | `waitGroupWorkflowTest` launches N `ExecuteChildWorkflow` via `WaitGroup.Go` then `Wait` aggregates sum | `internal/internal_workflow_test.go:1089-1118` |
| Semaphore bounded parallelism | `semaphoreImpl.Acquire` via `Await(size-cur>=n)`, `Release` panics on over-release | `internal/internal_workflow.go:75-79,2045-2070` |
| Sync child handling | `scheduledChildWorkflow.handle` / `handleFailedToStart` guard double-handle panic | `internal/internal_event_handlers.go:303-318` |
| Test env parent-close policy | `handleParentClosePolicy` ABANDON noop, TERMINATE `Complete(Terminated)`, REQUEST_CANCEL `cancelWorkflow` | `internal/internal_workflow_testsuite.go:1222-1237` |
| Should-stop event loop | `shouldStopEventLoop` keeps loop alive while ABANDON/REQUEST_CANCEL children unhandled | `internal/internal_workflow_testsuite.go:989-1006` |
| Panic propagation | `coroutineState.run` recovers panic → `workflowPanicError`, `dispatcher.ExecuteUntilAllBlocked` returns it, `executeDispatcher` `Complete(nil, panicErr)` | `internal/internal_workflow.go:1234-1245,1318-1321,692-698` |
| Panic test | `TestWorkflowPanic` asserts `PanicError` stack trace contains workflow func | `internal/internal_workflow_test.go:243-261` |
| Deterministic parallel example | `splitJoinActivityWorkflow` uses `Go` + `Channel` + `Selector` to join two activities | `internal/internal_workflow_test.go:153-204` |
| Child start vs result separation | `GetChildWorkflowExecution().Get` for handle, future `.Get` for result | `internal/internal_workflow.go:475-489` ; `internal/workflow.go:394-411` |
| Interceptor preserves semantics | `WorkflowOutboundInterceptor.ExecuteChildWorkflow` delegates to `Next`, not type-erased to generic executor | `internal/interceptor.go:241-243` ; `internal/interceptor_base.go:237-243` |
| Disconnected context usage advised | Workflow doc example for cleanup via `NewDisconnectedContext` after cancellation | `workflow/context.go:55-66` |

## Answers to Dimension Questions

**When is a parent allowed to become terminal relative to child cleanup?**

Parent terminal is gated by local coroutine completion (`syncWorkflowDefinition.Execute` sets `workflowResult` then `executeDispatcher` calls `env.Complete` only when root result pointer non-nil `internal/internal_workflow.go:699-734`), but not by child termination. Server applies `ParentClosePolicy` after parent `Complete`: `ABANDON` lets child outlive parent (no close action) `internal/internal_workflow_testsuite.go:1228-1229`; `TERMINATE` force-terminates child `internal/internal_workflow_testsuite.go:1230-1231`; `REQUEST_CANCEL` requests cancellation `internal/internal_workflow_testsuite.go:1232-1233`. Locally, `dispatcher.Close` `internal/internal_workflow.go:1365-1381` `Goexit`s remaining coroutines; remaining blocked `Await`/`Receive` are abandoned. `WaitForCancellation` (`ChildWorkflowOptions.WaitForCancellation` `internal/workflow.go:484-487`) controls whether parent waits for child cancel confirmation. So parent can be terminal while children are still running only under `ABANDON`; otherwise parent completion triggers child cancel/terminate via commands (`childWorkflowCommandStateMachine.cancel` `internal/internal_command_state_machine.go:770-784`).

**Does one child failure cancel siblings, and is that policy explicit?**

No implicit sibling cancellation. Each `ExecuteChildWorkflow` registers independent `scheduledChildWorkflow` state machines `internal/internal_event_handlers.go:637-643` and separate futures `internal/workflow.go:1400-1405`. Failure of one future does not auto-cancel siblings. Explicit policy requires shared `Context` from `WithCancel` `internal/context.go:179-183` and `propagateCancel` `internal/context.go:210-230`; sibling `ctx.Done` channels close only if they share parent `cancelCtx`. SDK tests show fail-fast must be hand-coded (e.g., select on error channel, then `cancel()`); `WaitGroup` tests deliberately do not cancel on error `internal/internal_workflow_test.go:1101-1103` (err overwritten, `Wait` just waits). This explicitness is documented but no `Group` helper enforces it.

**Are results ordered by submission, start, or completion?**

SDK primitives impose no ordering. `ExecuteChildWorkflow` returns a `ChildWorkflowFuture` per submission `internal/workflow.go:1399-1405`; `WaitGroup` example appends results as coroutines complete `internal/internal_workflow_test.go:1100-1102` (`results = append(results, result)` inside concurrent `Go`), yielding completion-order but race-prone. Ordered-by-submission requires iterating an array of futures sequentially (`for i, f := range futures { f.Get(ctx,&out[i]) }`), which is a user pattern, not an SDK helper. `Selector.AddFuture` `internal/internal_workflow.go:1407-1414` fires one ready future per `Select` in registration-order scan `internal/internal_workflow.go:1444-1558`, but when multiple futures become ready simultaneously, first case wins — still caller-defined. No `Gather` utility returns submission-ordered slices.

**Can a child outlive or mutate state after its parent is terminal?**

Yes, but only via server durability, not in-process mutation. With `ParentClosePolicy==ABANDON`, child execution persists on the server after parent `Complete` (test env `shouldStopEventLoop` keeps parent env alive until abandon children acknowledged, but real server allows true outlive `internal/internal_workflow_testsuite.go:998-1003`). Child cannot mutate parent workflow state after parent terminal because parent dispatcher is closed `internal/internal_workflow.go:667-670,1365-1381` and `Complete` is terminal (`isWorkflowCompleted` guard `internal/internal_workflow_testsuite.go:1143-1147`); channels/selectors are owned by closed dispatcher and `getState` panics outside workflow context `internal/internal_workflow.go:740-758`. Local activity/child callbacks after parent close are dropped via `Goexit` or ignored if `already completed`. Nexus/child signal paths are severed (`SignalChildWorkflow` requires live `GetChildWorkflowExecution` `internal/internal_workflow.go:479-489`).

## Architectural Decisions

- **Durable parent-child via protobuf command machine vs. in-memory supervisor:** `commandsHelper` + `childWorkflowCommandStateMachine` `internal/internal_command_state_machine.go:394-404,712-730` persists `StartChildWorkflowExecution` commands; history replay reconstructs ownership, surviving worker crashes. Tradeoff: local ownership is weak (coroutines die with parent), durability pushes lifecycle to server.
- **Deterministic cooperative dispatcher instead of OS threads:** `dispatcherImpl` `internal/internal_workflow.go:163-180,1281-1339` drives coroutines via `aboutToBlock`/`unblock` channels, `ExecuteUntilAllBlocked` loop, deadlock detection `internal/internal_workflow.go:1156-1186`. Ensures replay determinism but forbids native `go`/`chan`/`select` inside workflows (`assertNotInReadOnlyState` `internal/internal_workflow.go:760-768`).
- **Dual-future child handle (execution vs result):** `childWorkflowFutureImpl` `internal/internal_workflow.go:242-245,475-489` separates start confirmation (`GetChildWorkflowExecution`) from completion, allowing signal/cancel after start without blocking on result.
- **Context tree with explicit children map:** `cancelCtx` `internal/context.go:268-325` mirrors stdlib cancellation propagation for workflow `Context`, enabling `WithCancel`/`NewDisconnectedContext` patterns for cleanup `workflow/context.go:55-66`.
- **Policy-driven parent close, not reference counting:** `ParentClosePolicy` enum on `ChildWorkflowOptions` `internal/workflow.go:539-543` and `scheduledChildWorkflow.waitForCancellation` `internal/internal_event_handlers.go:80-87` delegate close semantics to server; test env mirrors via `handleParentClosePolicy` `internal/internal_workflow_testsuite.go:1222-1237`.

## Notable Patterns

- **WaitGroup collect-all:** `workflow.NewWaitGroup` + `wg.Go(ctx, fn)` + `wg.Wait(ctx)` `internal/workflow.go:303-321`, impl via `futureImpl`+`Settable` `internal/internal_workflow.go:62-68`; tests `internal/internal_workflow_test.go:1089-1132,1190-1218` show multi-wait and concurrent Wait panic.
- **Selector as structured join:** `NewSelector(ctx).AddReceive(ch, ...).AddFuture(f, ...).AddDefault(...).Select(ctx)` `internal/workflow.go:267-296`, dispatch in `selectorImpl.Select` `internal/internal_workflow.go:1433-1573`; enables dynamic fan-in without global bus.
- **Per-invocation data converters:** child workflow `ExecuteWorkflowParams.dataConverter`/`failureConverter` `internal/internal_workflow.go:228-229` captured from `WithDataConverter` context and propagated to `scheduledChildWorkflow` `internal/internal_event_handlers.go:629-643`, preserving executor-specific serialization without erasing concrete workflow type.
- **Eager coroutine priority:** `dispatcherImpl.NewCoroutine(..., highPriority)` `internal/internal_workflow.go:1247-1255,1257-1273` and `newEagerCoroutines` list ensures update handlers run before root, preserving deterministic ordering.

## Tradeoffs

- **Durability over local supervision:** Strength: child survives worker failure via server history; weakness: no in-process supervision tree — parent cannot inspect child stack or apply local backpressure; debugging relies on server history/replay tests `test/replaytests/workflows.go:677`.
- **Minimal composition primitives:** Strength: small API surface, workflow code stays explicit; weakness: every team re-implements `gather`, `first-error-wins`, `bounded pool` with subtle ordering/cancellation bugs (e.g., concurrent `results` append race in `internal/internal_workflow_test.go:1102` without mutex is safe only because coroutines are cooperatively scheduled, not preemptive — non-obvious).
- **Explicit cancellation:** Strength: precise control (`WithCancel` + `RequestCancelChildWorkflow` `internal/internal_event_handlers.go:417-420`); weakness: forgetting to share cancel context yields leaked children (ABANDON) or orphaned work.
- **Channel/Selector determinism vs ergonomics:** Channels are not Go native, `Select` is SDK-provided `internal/internal_workflow.go:1397-1431`; `Send`/`Receive` block via `yield` `internal/internal_workflow.go:1074-1088,814-856`, requiring careful flag handling (`SDKFlagBlockedSelectorSignalReceive` `internal/internal_workflow.go:1455-1466`) to avoid signal loss — leaky internal complexity exposed to user via subtle `HasPending` semantics.

## Failure Modes / Edge Cases

- **Duplicate child WorkflowID:** `commandsHelper.startChildWorkflowExecution` returns `childWorkflowExistsWithId` `internal/internal_event_handlers.go:621-627`; SDK surfaces as `ChildWorkflowExecutionAlreadyStartedError` `internal/workflow.go:2047-2074` via `failure_converter.go:285`. Caller must handle idempotency or use generated ID `internal/workflow.go:1418-1423`.
- **Negative WaitGroup counter / re-use before Wait returns:** panics `negative WaitGroup counter` `internal/internal_workflow_test.go:1257` and `WaitGroup is reused before previous Wait has returned` `internal/internal_workflow_test.go:1277`, routed as `workflowPanicError` to workflow failure `internal/internal_workflow.go:1319-1321`.
- **Double-handled activity/child/timer:** `scheduledChildWorkflow.handle` panics if `handled` `internal/internal_event_handlers.go:303-318`; `activityCommandStateMachine.handle` etc. similar `internal/internal_command_state_machine.go:288+`, surfacing nondeterminism as `[TMPRL1100]` history mismatch.
- **Context cancellation race:** child cancellation handler registered only after `startedCallback` `internal/workflow.go:1490-1503`; cancel before start is handled via immediate `ctx.Err()` check `internal/workflow.go:1407-1412` and via `canceledBeforeSent` state `internal/internal_command_state_machine.go:516-531`.
- **Panic in child coroutine:** any coroutine panic captured as `workflowPanicError` `internal/internal_workflow.go:1238-1241` and fails parent workflow deterministically, not isolated to child; stray yields during panic unwinding panic with `yield during panic unwinding` `internal/internal_workflow.go:1077-1083`.
- **Parent completes while Selector blocked:** `dispatcher.Close` closes coroutine via `Goexit` `internal/internal_workflow.go:1193-1219`; blocked `Select`/`Receive` abandons `blockedReceives` list; signal `SendAsync` after close would exceed buffer `workflow.go:928-939` if not guarded by policy.
- **Bounded parallelism via Semaphore over-release:** `semaphoreImpl.Release` panics if `cur<0` `internal/internal_workflow.go:2065-2070`; no timeout on `Acquire` except `Await` cancellation `internal/internal_workflow.go:2045-2054`.

## Future Considerations

- Provide structured `Group`/`Pool` helpers preserving concrete workflow-type results: `Group.Go(fn) → Future[T]`, `Group.Wait()` fail-fast vs `Gather()` ordered-by-submission, `Pool(n).Go(...)` bounded via `Semaphore` internally. Keep interceptor chain intact (wrap `WorkflowOutboundInterceptor` without erasing `dataConverter`/`failureConverter` per child).
- Add deterministic result-aggregation type `Gather[T any](ctx Context, futures []Future) ([]T, []error)` with explicit ordering policy (submission vs completion) and documented error handling (first error vs collect-all).
- Expose `ParentClosePolicy` lint/validation at `WithChildWorkflowOptions` time and document test-env `shouldStopEventLoop` ABANDON nuance to prevent accidental outlive leaks.
- Strengthen panic isolation option: child coroutine panic optionally surfaces as `ChildWorkflowExecutionError` (with `PanicError` cause) instead of failing parent, behind opt-in flag similar to `WaitForCancellation`.
- Add history-aware bounded-pool example using `Semaphore` + `Selector` pattern to replace unsynchronized `append` in `waitGroupWorkflowTest` `internal/internal_workflow_test.go:1102`.

## Questions / Gaps

- No evidence of built-in bounded-pool composition helper — searched `internal/workflow.go`, `internal/internal_workflow.go`, `internal/internal_workflow_test.go` for `Pool`/`LimitedGroup`/`errgroup` — only `Semaphore` primitive `internal/internal_workflow.go:75-79`; bounded pools are user-assembled. Gap: SDK leaves bounded parallelism ergonomics unaddressed.
- Synchronization of `results` slice in `waitGroupWorkflowTest` `internal/internal_workflow_test.go:1097-1103` appears racy under real goroutines but safe under cooperative dispatcher — not documented; verification requires reading dispatcher scheduling proof.
- Exact moment parent is considered terminal on server vs local `Complete` + `dispatcher.Close` ordering — server `ParentClosePolicy.TERMINATE` enforcement vs local `executeDispatcher` `internal/internal_workflow.go:692-734` interaction not fully observable without integration test history inspection (`test/integration_test.go:1412-1429` covers success but not mid-flight cancel).
- `WaitForCancellation bool` semantics beyond test assertion (`internal/internal_workflow_test.go:1219-1232` uses it but does not assert server behavior) — missing unit test for `waitForCancellation` signal path.
- Whether child events remain attributable without global bus: SDK uses per-workflow `commandsHelper` queues `internal/internal_command_state_machine.go:151-169` and per-signal channel map `internal/internal_workflow.go:209-211`, but no explicit per-child event stream abstraction — attribution is via `WorkflowExecution` IDs in history only.

---

Generated by `01.11-execution-composition-and-child-ownership` against `temporal-sdk-go`.
