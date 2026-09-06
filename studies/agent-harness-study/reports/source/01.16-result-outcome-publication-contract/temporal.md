# Source Analysis: temporal

## Result and Outcome Publication Contract

### Source Info

| Field | Value |
|-------|-------|
| Name | temporal |
| Path | `studies/agent-harness-study/sources/temporal` |
| Language / Stack | Go 1.26.4 / go.temporal.io/server, protobuf, gRPC |
| Analyzed | 2026-09-01 |

## Summary

Temporal Server turns a single SDK command into one irreversible terminal history event and one `WorkflowExecutionState/COMPLETED × WorkflowExecutionStatus={COMPLETED|FAILED|CANCELED|TERMINATED|TIMED_OUT|CONTINUED_AS_NEW}` row. The work return (`CompleteWorkflowExecutionCommandAttributes.Result: Payloads`, `FailWorkflowExecutionCommandAttributes.Failure: Failure`, `CancelWorkflowExecutionCommandAttributes.Details: Payloads`) is distinct from the runtime publication contract (history event + mutable-state transition + close-tasks). Immutability is enforced by an append-only `HistoryBuilder.EventStore` (`service/history/historybuilder/history_builder.go:26-31,78`) chained into a transactional `WorkflowContext.Lock` (`service/history/workflow/context.go:43,128-137`) with `DBRecordVersion` + `RangeID` OCC. Concurrent/late visibility is served via two publish primitives: generic `common/future.FutureImpl[T]` with atomic `pending→setting→ready` (`common/future/future_impl.go:11-19,66-84`) for in-memory waiters, and durable replay via `GetWorkflowExecutionHistory`/`GetOrPollWorkflowMutableState` long-poll (`service/history/api/getworkflowexecutionhistory/api.go:220-279`, `service/history/api/get_workflow_util.go:128-258`). Outcomes are frozen protobufs; waiter release is gated on `effect.Controller.OnAfterCommit` (`service/history/workflow/update/update.go:276,658,790`) after the history batch and `ExecutionInfo.CloseTime/CompletionEventBatchId` are durable.

## Rating

**8 / 10**

Rationale: Publication is strongly typed, exclusive, and irreversibly ordered. `setStateStatus` (`service/history/workflow/mutable_state_state_status.go:16-125`) plus `ValidateCreate/UpdateWorkflowStateStatus` (`common/persistence/workflow_state_status_validator.go:30-85`) rejects illegal `State×Status` combos; `ValidateCommandSequence` (`service/history/api/command_attr_validator.go:637-678`) enforces “close command must be last, alone.” Six terminal builders (`service/history/historybuilder/history_builder.go:424-466,603-609`) map 1:1 to status values and all go through `checkMutability` → `hBuilder.add` → `UpdateWorkflowStateStatus(COMPLETED, …)` → `GenerateWorkflowCloseTasks` → transactional commit (`service/history/workflow/mutable_state_impl.go:4941-5160,5649-5674`). In-memory fan-out uses `FutureImpl` that guarantees single `Set`, broadcast via `close(readyCh)` (`common/future/future_impl.go:83`), and tested for 1024 concurrent `Get`s (`common/future/future_test.go:95-125`). Late observers read the same immutable history event via pagination. Deductions: –1 for no defensive `bytes.Clone` on `Payloads` after publication (assignment in `service/history/historybuilder/event_factory.go:344,361` relies on protobuf immutability by convention), –1 for `FutureImpl.Set` panic-on-double-set (`common/future/future_impl.go:72-78`) rather than idempotent idempotency at caller unless `SetIfNotReady` is used.

## Evidence Collected

Every entry MUST include a file path with line numbers. Format: `path/to/file.ts:NN`.

| Area | Evidence | File:Line |
|------|----------|-----------|
| Work-function return type (complete) | `CompleteWorkflowExecutionCommandAttributes{Result *commonpb.Payloads}` validated (empty allowed) | `service/history/api/command_attr_validator.go:219-228` |
| Work-function return type (fail) | `FailWorkflowExecutionCommandAttributes{Failure *failurepb.Failure}` required non-nil | `service/history/api/command_attr_validator.go:230-242` |
| Work-function return type (cancel) | `CancelWorkflowExecutionCommandAttributes{Details *commonpb.Payloads}` (optional) | `service/history/api/command_attr_validator.go:244-253` |
| Close-command exclusivity | `ValidateCommandSequence` enforces close command last and alone, `COMPLETE|FAIL|CANCEL|CONTINUE_AS_NEW` close set | `service/history/api/command_attr_validator.go:637-678` |
| Terminal-outcome type | `enumspb.WorkflowExecutionStatus={RUNNING,COMPLETED,FAILED,CANCELED,TERMINATED,CONTINUED_AS_NEW,TIMED_OUT,PAUSED}` + `enumsspb.WorkflowExecutionState={CREATED,RUNNING,COMPLETED,ZOMBIE,CORRUPTED}` | `common/persistence/workflow_state_status_validator.go:10-27` |
| Terminal-outcome construction | `setStateStatus` transition matrix: CREATED→RUNNING/COMPLETED/ZOMBIE, RUNNING→RUNNING/COMPLETED/ZOMBIE, COMPLETED→COMPLETED same status only | `service/history/workflow/mutable_state_state_status.go:16-125` |
| Persistence guard | `ValidateCreateWorkflowStateStatus`/`ValidateUpdateWorkflowStateStatus` called pre-DB; `COMPLETED` rejects `RUNNING|PAUSED` | `common/persistence/workflow_state_status_validator.go:30-85` |
| Mutable-state publication | `UpdateWorkflowStateStatus` short-circuits identity, sets `executionStateUpdated` flag, delegates to `setStateStatus` | `service/history/workflow/mutable_state_impl.go:7392-7406` (referenced `service/history/workflow/mutable_state_impl.go:4966-4982` apply paths) |
| Outcome builders (6) | `AddCompletedWorkflowEvent`, `AddFailWorkflowEvent`, `AddTimeoutWorkflowEvent`, `AddWorkflowExecutionTerminatedEvent`, `AddWorkflowExecutionCanceledEvent` | `service/history/historybuilder/history_builder.go:424-432,434-448,450-456,458-466,603-609` |
| Outcome apply (immutable) | `ApplyWorkflowExecutionCompletedEvent` → `UpdateWorkflowStateStatus(COMPLETED,COMPLETED)`, sets `CompletionEventBatchId`, `CloseTime`, `ClearStickyTaskQueue`, `writeEventToCache`, `processCloseCallbacks` | `service/history/workflow/mutable_state_impl.go:4941-4982` (`4966-4982` apply) |
| Outcome apply failed | `ApplyWorkflowExecutionFailedEvent` → `COMPLETED/FAILED`, `NewExecutionRunId` from failure attrs, branches on `RetryState!=IN_PROGRESS` | `service/history/workflow/mutable_state_impl.go:5010-5033` |
| Outcome apply canceled | `ApplyWorkflowExecutionCanceledEvent` → `COMPLETED/CANCELED` | `service/history/workflow/mutable_state_impl.go:5120-5160` (`5144-5160` apply) |
| Outcome apply timeout | `ApplyWorkflowExecutionTimedoutEvent` → `COMPLETED/TIMED_OUT` | `service/history/workflow/mutable_state_impl.go:5035-5082` (`5059-5082` apply) |
| Outcome apply terminated | `AddWorkflowExecutionTerminatedEvent` → `hBuilder.Add…`, `Apply…`, `GenerateWorkflowCloseTasks(deleteAfterTerminate)` | `service/history/workflow/mutable_state_impl.go:5649-5674` |
| Event-factory copy | `CreateCompletedWorkflowEvent{Result: command.Result}`, `CreateFailWorkflowEvent{Failure: command.Failure}` — direct assignment, no clone, relies on proto immutability | `service/history/historybuilder/event_factory.go:335-349,351-367` |
| Cancel-fact vs failure | `AddWorkflowExecutionCancelRequestedEvent` sets `executionInfo.CancelRequested=true` (non-terminal) vs `Canceled` (terminal) | `service/history/workflow/mutable_state_impl.go:5084-5118` |
| History immutability | `HistoryBuilderState{Mutable,Immutable,Sealed}`; `workflowFinished` flag; `memEventsBatches` append-only; `IsDirty` guard | `service/history/historybuilder/history_builder.go:26-31,78` |
| Transactional publication | `checkMutability` guard at start of each `Add*` + `PrioritySemaphore(1)` per execution (`workflow.ContextImpl.lock`) | `service/history/workflow/mutable_state_impl.go:4946`, `service/history/workflow/context.go:43,116,128-137` |
| Future contract | `Future[T]{Get(ctx) (T,error), Ready() bool}`, `FutureImpl[T]{status atomic pending→setting→ready, readyCh chan struct{}}` | `common/future/future.go:5-10`, `common/future/future_impl.go:22-29,11-19` |
| Future publication | `Set` CAS `pending→setting`, assign `value/err`, CAS `setting→ready`, `close(readyCh)`; panic on double-Set | `common/future/future_impl.go:66-84` |
| Future idempotent variant | `SetIfNotReady` returns false instead of panic | `common/future/future_impl.go:87-104` |
| Ready future (late waiter) | `ReadyFutureImpl[T]` always `Ready()=true`, `Get` returns immediate `value,err` without channel | `common/future/ready_future_impl.go:14-32` |
| Concurrent waiter tests | 1024 parallel `Get` blocked until `Set`, then all observe same value; `Ready→Get` with canceled ctx still returns value | `common/future/future_test.go:95-125,127-170,172-183` |
| Late-waiter / repeated Get | `TestSetReadyGet_*` verifies `Ready()=true` then `Get` returns same committed facts without blocking | `common/future/future_test.go:127-170` |
| Context-cancel wait | `Get` selects `readyCh` vs `ctx.Done()`; canceled ctx returns `ctx.Err()` only if not ready, otherwise committed outcome wins (`GetIfReady` → `errorFutureNotReady`) | `common/future/future_impl.go:42-56,58-64` |
| Update as second publication example | `Update{accepted Future[*Failure], outcome Future[*Outcome]}`; `WaitLifecycleStage` selects outcome vs accepted vs admitted by `Ready()` before blocking | `service/history/workflow/update/update.go:60-62,156-253` |
| Waiter release ordering | Futures completed inside `effect.Controller.OnAfterCommit` after `EventStore.Add*` succeeds; rollback restores prior state via `OnAfterRollback` | `service/history/workflow/update/update.go:276-313,658-703,790-815` |
| Durable late-consumer path | `GetOrPollWorkflowMutableState` long-poll: `WatchHistoryEvent` → `<-channel` or `longPollCtx.Done()` or `ctx.Done()`, rechecks `expectedNextEventID < response.NextEventId` to drop stale notifs | `service/history/api/get_workflow_util.go:128-258` (as documented in `01.08` report), `service/history/historybuilder/history_builder.go:458-466` close tasks |
| History long-poll | `GetWorkflowExecutionHistory` long-poll branch `WaitNewEvent` re-queries mutable state when `IsWorkflowRunning && token==nil` | `service/history/api/getworkflowexecutionhistory/api.go:220-279` (evidence via reports) |
| Diagnostic preservation | `Failure` fields propagated: `WorkflowTaskFailedEvent{ Failure *Failure}`, `WorkflowExecutionFailedEvent{ Failure}`, `ActivityTaskFailedEvent{ Failure, RetryState}` | `service/history/historybuilder/event_factory.go:202-228,297-315,351-367` |
| Completion metadata | `ExecutionInfo.CompletionEventBatchId`, `CloseTime`, `NewExecutionRunId` set atomically on each `Apply*` | `service/history/workflow/mutable_state_impl.go:4976-4978,5020-5022,5054-5069,5154-5156` |
| Termination with payload copy | `WorkflowExecutionTerminatedEventAttributes{Reason string, Details Payloads, Identity string}` | `service/history/historybuilder/event_factory.go:383-399` |

## Answers to Dimension Questions

**1. Is a work return distinct from the runtime's terminal outcome?**

Yes — strongly separated. The SDK work return is a per-command payload: `CompleteWorkflowExecutionCommandAttributes.Result` (`service/history/api/command_attr_validator.go:219-228`) is a `*commonpb.Payloads`, `FailWorkflowExecutionCommandAttributes.Failure` is a `*failurepb.Failure` (`service/history/api/command_attr_validator.go:230-242`), `CancelWorkflowExecutionCommandAttributes.Details` is optional `Payloads`. The runtime terminal outcome is not that value; it is the pair `(WorkflowExecutionState=COMPLETED, WorkflowExecutionStatus)` plus the appended history event (`EVENT_TYPE_WORKFLOW_EXECUTION_COMPLETED|FAILED|CANCELED|TERMINATED|TIMED_OUT|CONTINUED_AS_NEW`) and persisted `ExecutionInfo{CloseTime, CompletionEventBatchId, NewExecutionRunId}` (`service/history/workflow/mutable_state_impl.go:4966-4982,5010-5033,5144-5160`). `ValidateCommandSequence` (`service/history/api/command_attr_validator.go:637-678`) ensures exactly one close command is the last command, so the work return cannot be confused with a mid-workflow `RecordMarker` or `Activity` result. Even after the SDK returns, the server may rewrite the outcome (e.g., `AddTimeoutWorkflowEvent` via timer, `forceTerminateWorkflow` via size limits in `service/history/workflow/context.go:1406-1446`) without SDK cooperation.

**2. Which result, failure, cancellation, and runtime-fault combinations are valid?**

Exactly one status per closed execution, enforced by `setStateStatus` (`service/history/workflow/mutable_state_state_status.go:61-64: RUNNING→COMPLETED rejects RUNNING|PAUSED`) and `ValidateUpdateWorkflowStateStatus` (`common/persistence/workflow_state_status_validator.go:70-84`):

| Terminal status | Legal store | Source event | Payload slot |
|---|---|---|---|
| `COMPLETED` | success | `WorkflowExecutionCompletedEvent{Result Payloads}` (`service/history/historybuilder/event_factory.go:335-349`) |
| `FAILED` | failure | `WorkflowExecutionFailedEvent{Failure, RetryState, NewExecutionRunId}` (`service/history/historybuilder/event_factory.go:351-367`) |
| `CANCELED` | cancellation fact | `WorkflowExecutionCanceledEvent{Details Payloads}` (`service/history/historybuilder/event_factory.go:600-612`) + prior `CancelRequested` flag (`service/history/workflow/mutable_state_impl.go:5113-5118`) |
| `TERMINATED` | runtime fault/operator action | `WorkflowExecutionTerminatedEvent{Reason string, Details Payloads}` (`service/history/historybuilder/event_factory.go:383-399`) |
| `TIMED_OUT` | runtime fault (deadline) | `WorkflowExecutionTimedOutEvent{RetryState, NewExecutionRunId}` (`service/history/historybuilder/event_factory.go:369-381`) |
| `CONTINUED_AS_NEW` | successful handoff | `WorkflowExecutionContinuedAsNewEvent{NewExecutionRunId, Input, Header, RetryPolicy…}` (`service/history/historybuilder/event_factory.go:476-506`) |
| `PAUSED`/`RUNNING` | not terminal; never under `STATE_COMPLETED` per `setStateStatus:74-85,`83-85`` | 

No combination of value+error in one outcome: `Complete` carries only `Result`, `Fail` only `Failure`, `Cancel` only `Details`, `Terminated` only `Reason+Details`. A workflow with a retry policy that fails with `RETRY_STATE_IN_PROGRESS` writes `Failed` but then `processUpdateCloseCallbacks` vs `processCloseCallbacks` (`service/history/workflow/mutable_state_impl.go:5026-5032`) indicates a new run will inherit callbacks — the failure is still the terminal fact for that run.

**3. What happens when work returns both a value and an error?**

Impossible at the server contract. `ValidateCompleteWorkflowExecutionAttributes` accepts any `Result` (including nil/empty) and does not inspect failure (`service/history/api/command_attr_validator.go:219-228`), while `ValidateFailWorkflowExecutionAttributes` requires `Failure != nil` and rejects nil (`service/history/api/command_attr_validator.go:238-240`). `ValidateCommandSequence` (`service/history/api/command_attr_validator.go:637-678`) guarantees at most one close command per `RespondWorkflowTaskCompleted`. The worker cannot send both `CompleteWorkflowExecution` and `FailWorkflowExecution` in the same batch; the server would reject the second as “must be last command”. If a workflow function in the SDK returns `(value, error)` with both non-nil, the SDK chooses one translation (error wins → `Fail`), so only one event is ever appended. No server code path merges both into a single `Completed+Failure` outcome — `CreateCompletedWorkflowEvent` writes only `Result` (`service/history/historybuilder/event_factory.go:342-346`), `CreateFailWorkflowEvent` writes only `Failure` (`service/history/historybuilder/event_factory.go:359-363`).

**4. Can a terminal outcome expose a partial or mutable value?**

No partial variant exists. Each terminal builder writes a single atomic batch via `hBuilder.add(event)` (`service/history/historybuilder/history_builder.go:424-432` pattern) and `Apply*` sets `COMPLETED/<status>` in one step; there is no “partial result” field. `CompletionEventBatchId` is written atomically with the event (`service/history/workflow/mutable_state_impl.go:4976,5020,5069`). Payloads are protobuf `repeated bytes` which are immutable once persisted to `history` table and returned via `GetWorkflowExecutionHistory` as copies; `CreateCompletedWorkflowEvent` assigns `command.Result` by pointer (`service/history/historybuilder/event_factory.go:344`) without defensive `Clone`, so aliasing is possible only if the caller retains and mutates the original slice before commit — but after `hBuilder.add` the `HistoryEvent` is owned by the `EventStore` and subsequent reads deserialize a new instance from storage, so post-publication mutation cannot affect committed facts. `GenerateWorkflowCloseTasks` is called *after* the status transition but *before* the DB transaction commits (`service/history/workflow/mutable_state_impl.go:4999-5005`), so visibility tasks never race ahead of the close event.

**5. Do concurrent, repeated, and late waiters observe the same committed facts?**

Yes, by two layers.

*In-memory (`FutureImpl`):* `pending→setting→ready` CAS (`common/future/future_impl.go:72-83`) guarantees single publication; `close(readyCh)` broadcasts to all waiters. `TestSetGetReady_Parallel` with 1024 concurrent `Get` blocked on `readyCh` asserts all receive identical `value,err` (`common/future/future_test.go:95-125`); `TestSetReadyGet_Parallel` asserts late waiters that poll `Ready()` then `Get` also see the same facts (`common/future/future_test.go:137-170`). Repeated `Get` after ready is idempotent — `Ready()` shortcut returns `value,err` without re-blocking (`common/future/future_impl.go:44-47`).

*Durable (workflow close):* Late consumers via `GetWorkflowExecutionHistory` (paginated, with `CompletionEventBatchId` pointer) or `GetOrPollWorkflowMutableState` (`service/history/api/get_workflow_util.go:128-258`) read the same committed `HistoryEvent`; history is append-only and `COMPLETED` state forbids mutation (`service/history/workflow/mutable_state_state_status.go:74-92: COMPLETED→COMPLETED requires same status else error`). `ReadyFutureImpl` in `service/history/workflow/update/registry.go:218-224,455-480` loads already-completed updates from storage as `newCompleted(id, NewReadyFuture(outcome, err))`, so a poller that arrives after completion immediately returns the stored `Outcome` without waiting.

**6. Can waiting be cancelled without cancelling the run itself?**

Yes — waiter cancellation is decoupled from execution cancellation.

- Generic future: `FutureImpl.Get(ctx)` selects on `ctx.Done()` (`common/future/future_impl.go:49-55`); if the caller cancels context before `Set`, `Get` returns `ctx.Err()` with zero value and `Ready()` stays `false`. `TestGetWhenContextCanceled` shows that if `Ready()==true`, `Get` ignores the canceled context and returns the committed outcome (`common/future/future_test.go:172-183`). Later `Get` with fresh context will still observe the outcome.
- History long-poll: `GetOrPollWorkflowMutableState` takes distinct `ctx` (request lifetime) and `longPollCtx` (soft timeout); expiry returns empty response, not workflow cancel (`service/history/api/get_workflow_util.go:128-258` pattern, `service/history/api/getworkflowexecutionhistory/api.go:220-279`).
- Update poll: `PollWorkflowExecutionUpdate` → `Update.WaitLifecycleStage(ctx, waitStage, softTimeout)` (`service/history/api/pollupdate/api.go:19-82`, `service/history/workflow/update/update.go:156-253`) uses `stCtx` with `softTimeout`; soft timeout returns `statusAdmitted/Accepted` with empty outcome, hard `ctx` cancellation returns `ctx.Err()` without aborting the update. Aborting the update (`Abort(AbortedByServerErr)` → `registryClearedErr`) is retried by caller (`service/history/workflow/update/update.go:238-240`), not a workflow cancel.
- Workflow cancel is an explicit separate command: `RequestCancelWorkflowExecution` sets `CancelRequested=true` non-terminal (`service/history/workflow/mutable_state_impl.go:5084-5118`), while `CancelWorkflowExecution` writes `Canceled` terminal (`service/history/workflow/mutable_state_impl.go:5120-5160`).

**7. Who owns values and diagnostics after publication?**

The **history event** owns them durably; callers receive copies.

- `Payloads` (`Result`, `Details`) and `Failure` are stored inside a `HistoryEvent` which is serialized to the `history` table and cached via `writeEventToCache` (`service/history/workflow/mutable_state_impl.go:4980,5024,5158`). Readers get deserialized copies via `GetWorkflowExecutionHistory`/`DescribeWorkflowExecution` (`service/frontend/workflow_handler.go:960-1000` as referenced).
- The `ExecutionInfo` mirrors only the batch pointer and timestamps (`CompletionEventBatchId`, `CloseTime`, `NewExecutionRunId` per `service/history/workflow/mutable_state_impl.go:4976-4978`), not the payload bytes themselves.
- No `bytes.Clone` or `proto.Clone` is performed in `Create*WorkflowEvent` (`service/history/historybuilder/event_factory.go:344,361`), so ownership transfer is by convention: callers must not mutate the input `Payloads`/`Failure` after handing it to `Add*` (the builder owns the reference until persistence). After persistence, the storage layer is the sole owner; in-memory `Update.outcome/accepted` futures hold a single copy and `Get` returns that pointer without copy-on-write — callers must treat it as read-only (no `DeepCopy` API exists). Diagnostics (`Failure.message`, `stack_trace`, `cause *Failure`, `RetryState`) are retained verbatim in the event (`service/history/historybuilder/event_factory.go:359-363`).

**8. Does waiter release occur only after the complete outcome is visible?**

Yes — release is deferred to `OnAfterCommit`.

- For updates: `onAcceptanceMsg` sets state `ProvisionallyAccepted` then registers `OnAfterCommit` to `Set(nil,nil)` on `accepted` future only after `AddWorkflowExecutionUpdateAcceptedEvent` is durable (`service/history/workflow/update/update.go:658-695`). `onResponseMsg` and `abort` follow the same pattern (`service/history/workflow/update/update.go:790-807,276-306`) with ordering invariant: when acceptance+completion happen in one WFT, `accepted` is set *after* `outcome` (`stateProvisionallyCompletedAfterAccepted` gate at `service/history/workflow/update/update.go:675-691,796-806`).
- For workflow close: `AddCompletedWorkflowEvent` does `checkMutability → hBuilder.Add… → Apply… (UpdateWorkflowStateStatus, set CloseTime/BatchId, writeEventToCache) → GenerateWorkflowCloseTasks` *before* returning (`service/history/workflow/mutable_state_impl.go:4941-4964`); the enclosing `Context.UpdateWorkflowExecutionAsActive` commits via `NewTransaction(...).UpdateWorkflowExecution` (`service/history/workflow/context.go:884-897`). `writeEventToCache` makes the event visible to in-process readers only after the `EventStore` mutation, and external pollers that use `GetOrPollWorkflowMutableState` loop until `NextEventId` advances (`service/history/api/get_workflow_util.go:218-221` check), so they never see a pre-close snapshot.

## Architectural Decisions

- **Single-command close with sequence validator.** `ValidateCommandSequence` (`service/history/api/command_attr_validator.go:637-678`) treats `Complete|Fail|Cancel|ContinueAsNew` as “must be last” — simplifies the publication contract to one terminal event per batch; tradeoff is no mid-batch `Complete` with follow-up cleanup commands.
- **Two-level state/status split.** `WorkflowExecutionState` (lifecycle) vs `WorkflowExecutionStatus` (terminal reason) with a 4×8 matrix (`service/history/workflow/mutable_state_state_status.go:16-125`) gives persistence validators (`common/persistence/workflow_state_status_validator.go:56-85`) a fail-closed guard independent of history events. Cost is a matrix that must be kept consistent between `setStateStatus` and `ValidateUpdateWorkflowStateStatus`.
- **Append-only EventStore as source of truth.** Invalidation on `CORRUPTED` and `IsDirty` gating (`service/history/workflow/mutable_state_impl.go:4947,5001-5006`) means visibility, replication, and archival all derive from the same serialized history — strong auditability, but forces `CompletionEventBatchId` indirection for loading close payload lazily.
- **CAS-fenced future + channel broadcast.** `FutureImpl` uses `pending→setting→ready` with two `CompareAndSwap` to prevent data race between mutating `value/err` and publishing `readyCh` close (`common/future/future_impl.go:11-19,72-83`). `SetIfNotReady` provides non-panic arbitration for races (e.g., shard `future` controllers). Tradeoff is global promise style rather than Go2+ `context.AfterFunc`.
- **Effect-controller commit barrier.** Update lifecycle uses `effect.Controller{OnAfterCommit, OnAfterRollback}` to decouple history durability from in-memory future completion (`service/history/workflow/update/update.go:276-313`). This prevents “waiter unblocked but event not persisted” anomalies at cost of a provisional-state enum (`ProvisionallyAccepted/Completed/Aborted`).

## Notable Patterns

- **Idempotent admission via `Update.Admit` deduplication.** `FindOrCreate` + `Admit` returns nil if `state != Created` (`service/history/workflow/update/update.go:320-325`), giving inherent `updateID` deduplication without extra storage round-trip, reinforced by `TryResurrect` (`service/history/workflow/update/registry.go:238-281`).
- **Provisional → committed state bridge for atomic commit.** Both `MutableState` (`checkMutability` at `service/history/workflow/mutable_state_impl.go:4947…`) and `Update` state machines buffer effects and flip on commit (`service/history/workflow/update/update.go:354-368`), mirroring DB transaction staging.
- **PrioritySemaphore per execution.** `workflow.ContextImpl.lock = NewPrioritySemaphore(1)` (`service/history/workflow/context.go:43,116`) serializes all transactions on one run, so parallel `RespondWorkflowTaskCompleted` calls for the same execution queue behind cache `GetWorkflowLease(..., PriorityHigh)` (`service/history/api/pollupdate/api.go:29-43`).
- **Long-poll drop-on-backpressure notifier.** `events.Notifier` uses non-blocking `dispatchHistoryEventNotification` with `default:` drop (`service/history/events/notifier.go:191-198` per reports) + stale-check in `GetOrPollWorkflowMutableState` (`service/history/api/get_workflow_util.go:218-221`), favoring liveness over lossless wakeups.

## Tradeoffs

- **Execution-local future vs durable visibility:** `FutureImpl` is pure in-memory, lost on failover/restart. Durable `HistoryEvent` survives via `CompletionEventBatchId`. System keeps both: `FutureImpl` for low-latency in-process `PollWorkflowExecutionUpdate` (`service/history/workflow/update/update.go:156-253`), history for cross-shard/cluster recovery. Cost is dual source of truth that must be reconciled via `NewReadyFuture` on reload (`service/history/workflow/update/registry.go:455-480`).
- **Assignment without Clone for payloads:** `event_factory.go:344 Result: command.Result` shares underlying `[]byte` slices until serialization. Saves allocation on hot path, but relies on callers not retaining alias; a defensive `proto.Clone` would cost ~2× for large results (embed document pattern). Chosen trade favors throughput.
- **Panic vs soft-fail on double Set:** `FutureImpl.Set` panics (`common/future/future_impl.go:77`) catches bugs fast in tests, but production code must remember to call `SetIfNotReady` in contender paths (e.g., shard controller `common/future` usage). Alternative would be silent ignore, hiding races.
- **Close-command last constraint:** Prevents atomic “complete + schedule cleanup child” in one batch, forcing a preceding batch. Simplifies `WorkflowCloseTasks` generation (`service/history/workflow/mutable_state_impl.go:4999-5005`) but requires two round trips for some SDK patterns.
- **One task per run for `GenerateWorkflowCloseTasks`:** Includes `VisibilityCloseTask`, `ArchivalTask`, `CloseCallbacks` in one sweep; guarantees “outcome visible before listing side effects” vs eventual visibility store lag (`service/history/workflow/mutable_state_impl.go:4999-5005,5166-5170`).

## Failure Modes / Edge Cases

- **Double-close attempt:** `checkMutability` fails second `AddCompletedWorkflowEvent` if execution already `COMPLETED`; `setStateStatus` (`service/history/workflow/mutable_state_state_status.go:74-85`) returns `serviceerror.Internalf("unable to change workflow state…")`; caller bubbles to workflow task failure with cause `WORKFLOW_TASK_FAILED_CAUSE_*` rather than silent second outcome.
- **CancelRequested without Cancel:** `CancelRequested=true` (`service/history/workflow/mutable_state_impl.go:5116`) is non-terminal; if worker never issues `CancelWorkflowExecution`, workflow remains `RUNNING` until timeout/termination — client poll on outcome blocks until soft timeout then returns `ADMITTED` stage with empty outcome (`service/history/workflow/update/update.go:245-252`).
- **Workflow size-triggered termination:** `forceTerminateWorkflow` (`service/history/workflow/context.go:1406-1446`) may *overwrite* the SDK-intended close outcome by clearing `MutableState` and applying `TerminateWorkflow` under shard context. The SDK close payload is discarded; `Terminate` details= nil and reason `FailureReasonHistorySizeExceedsLimit`. No torn partial outcome because `Clear()` drops the pending transaction before `LoadMutableState`.
- **Future double-set race:** Concurrent callers calling `Set` simultaneously panic on second CAS (`common/future/future_impl.go:72-78`). In `Update.abort` the `accepted` future is set only if state was pre-accepted (`service/history/workflow/update/update.go:302-305`), else leaves as is — but simultaneous `onResponseMsg` + `abort` in same WFT rely on `stateProvisionallyCompletedAfterAccepted` ordering guard (`service/history/workflow/update/update.go:285-298`).
- **Backpressured notifier drop:** `dispatchHistoryEventNotification` default branch drops wakeup if subscriber channel buffer full (`service/history/events/notifier.go:191-198`); stale-recheck in `GetOrPollWorkflowMutableState` recovers by looping on `expectedNextEventID < response.NextEventId` (`service/history/api/get_workflow_util.go:218-258`), but spurious wakeups can delay long-poll by one interval.
- **Mutable-state size vs payload size:** Very large `Result` payloads (e.g., >2MB) pass `ValidateComplete` (no size check) but fail later at `AddTasks` history-size enforcement in `Context.enforceHistorySizeCheck` (`service/history/workflow/context.go:1289-1302`) → `forceTerminate` with `ErrHistorySizeExceedsLimit`, converting SUCCESS into TERMINATED — a subtle violation where size policy masquerades as failure reason.
- **ContinuedAsNew with `NewExecutionRunId` empty:** If `newExecutionRunID == ""` (no retry), `ExecutionInfo.NewExecutionRunId` is set to `""` (`service/history/workflow/mutable_state_impl.go:4977`) and `GenerateWorkflowCloseTasks` treats empty as non-chained run; visibility `NewRunID` search attribute absent, so chaining via visibility fails even though history chain exists.

## Future Considerations

- **Defensive payload freeze.** Add `payload.Clone(payloads)` or `proto.Clone(Failure)` at `EventFactory.Create*` boundary (`service/history/historybuilder/event_factory.go:344,361`) behind a build tag to catch aliasing bugs in tests while keeping hot path zero-copy in prod.
- **Idempotent close builder.** Replace `FutureImpl.Set` panic with `SetIfNotReady` at all `Update` completion points (`service/history/workflow/update/update.go:286,695,795,759`) to tolerate double-commit under at-least-once `OnAfterCommit` replay; currently relies on state guard to prevent double call.
- **Unified outcome future for workflow close.** Expose a `WorkflowCloseFuture` (analogous to `Update.outcome`) from `WorkflowContext` so frontend long-polls (`GetWorkflowExecutionHistory WaitNewEvent`) can await `CloseTime != nil` via the same `future` primitive instead of polling `GetOrPollWorkflowMutableState`, reducing tail latency.
- **Close-command coupling.** Allow `CompleteWorkflowExecution` to carry `UpsertSearchAttributes` or `ModifyWorkflowProperties` as atomic sub-attributes (not separate commands) so SDK can publish result + visibility in one batch without violating `ValidateCommandSequence` (`service/history/api/command_attr_validator.go:637-678`).
- **Diagnostic backfill for terminated runs.** Attach `Failure.stack_trace` and `cause` chain to `WorkflowExecutionTerminatedEvent` even when `forceTerminateWorkflow` path is taken (`service/history/workflow/context.go:1438-1445` currently nil `Details`), improving post-mortem without extra lookup.

## Questions / Gaps

- **No evidence for client-side wait-cancel without server cancel distinct metric:** Grep for `ResourceExhaustedCauseTag("...")` on long-poll client timeouts found only update metrics (`WorkflowExecutionUpdateClientTimeout` at `service/history/workflow/update/update.go:182`) — no similar counter for history long-poll client cancel, so cost of many idle watchers is not observable.
- **CHASM vs legacy path divergence:** `Context.ForceTerminateWorkflow` branches to `chasm.Terminate` vs `TerminateWorkflow` (`service/history/workflow/context.go:1430-1445`); no test found proving both paths preserve `Failure.cause` identically — inferred but not verified.
- **Payload size limits:** Search for `MaxPayloadSize` or `BlobSizeLimitError` in `service/history/api/command_attr_validator.go` returned only timeout-length checks; actual blob limit enforcement point for workflow results could not be located — may live in `common/payload` or `service/frontend` handler, leaving “what size result fails vs terminates” ambiguous.
- **No in-repo evidence for concurrent waiters beyond `FutureImpl`:** History-level concurrent pollers (multiple frontend instances watching same `RunID`) rely on `ShardController` future coordination (`service/history/shard/controller_impl.go:13`) but no test exercising multi-shard `GetOrPollWorkflowMutableState` fanout was found.

---

Generated by `Dimension 01.16: Result and Outcome Publication Contract` against `temporal`.
