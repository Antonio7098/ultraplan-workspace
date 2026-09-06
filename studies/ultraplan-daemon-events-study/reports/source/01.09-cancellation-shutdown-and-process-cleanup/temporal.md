# Source Analysis: temporal

## Dimension 01.09: Cancellation, Shutdown, and Process Cleanup

### Source Info

| Field | Value |
|-------|-------|
| Name | temporal |
| Path | `studies/ultraplan-daemon-events-study/sources/temporal` |
| Language / Stack | Go monorepo (fx DI, gRPC: frontend/history/matching/worker, event-sourced mutable state + transfer/timer queues) |
| Analyzed | 2026-09-03 |

## Summary

Temporal implements cancellation as a durable, cooperative, idempotent control request, not an interrupt. Client `RequestCancelWorkflowExecution` (`service/frontend/workflow_handler.go:2265`) is forwarded to history (`service/history/handler.go:934`, `service/history/history_engine.go:679`) and persisted as `EVENT_TYPE_WORKFLOW_EXECUTION_CANCEL_REQUESTED` via `service/history/api/requestcancelworkflow/api.go:82` / `service/history/workflow/mutable_state_impl.go:5091-5118` (sets `executionInfo.CancelRequested=true` at `service/history/workflow/mutable_state_impl.go:5123`). It is durable without an attached worker: the write also schedules a new workflow task (`service/history/api/requestcancelworkflow/api.go:86`) so a later poller observes it. Runtime acknowledgement is a separate SDK command `COMMAND_TYPE_CANCEL_WORKFLOW_EXECUTION` handled at `service/history/api/respondworkflowtaskcompleted/workflow_task_completed_handler.go:951-979`, closing the workflow as `EVENT_TYPE_WORKFLOW_EXECUTION_CANCELED` (`service/history/historybuilder/event_factory.go:600-612`, `service/history/workflow/mutable_state_impl.go:5127-5167`). Propagation to activities is via heartbeat-pollable `CancelRequested` flag (`service/history/api/recordactivitytaskheartbeat/api.go:77`) plus eager `WorkerCommand_CancelActivity` (`service/history/api/respondworkflowtaskcompleted/workflow_task_completed_handler.go:693-760`); to child workflows via transfer-queue `CancelExecutionTask` re-invoking history cancel RPC (`service/history/transfer_queue_active_task_executor.go:1612-1640`). Daemon shutdown is drain-then-`GracefulStop` per service in reverse `initOrder` (`temporal/server_impl.go:47-53`, `temporal/server_impl.go:109-124`) bounded by `serviceStopTimeout=5m` (`temporal/server.go:12`, `temporal/fx.go:371`), triggered only by SIGINT/SIGTERM (`temporal/interrupt.go:9-21`). There is no OS process-group/cgroup kill or SIGKILL escalation in the server binary (verified by empty grep for `Setpgid|SIGKILL|cgroup`). Cleanup-after-cancel uses SDK-side detached scopes, in-server only as detached `context.Context` copies (`service/history/shard/context_impl.go:2408-2433`).

## Rating

**7 / 10** — Clear model with tests, explicit interfaces, and operational safeguards.

Rationale: durable cancel intent + idempotent ack + race arbitration (`IsWorkflowExecutionRunning` / `IsCancelRequested` guards) + transfer-queue child fan-out + heartbeat propagation are all implemented, interfaced, and covered by functional tests (`tests/cancel_workflow_test.go:33-670`, `service/history/transfer_queue_active_task_executor_test.go:1575-1750`). Deducted because (a) daemon shutdown drain defaults are `0s` (`common/dynamicconfig/constants.go:869-873,1425-1429,1918-1922`), (b) there is no process-tree/group/cgroup termination or kill-escalation — inapplicable to Temporal's goroutine model but a gap against this dimension's OS-process expectations, (c) worker service `Stop()` has no `GracefulStop`-fallback timer unlike history/matching (`service/worker/service.go:305-323` vs `service/history/service.go:164-168`).

## Evidence Collected

Every entry MUST include a file path with line numbers. Format: `path/to/file.ts:NN`.

| Area | Evidence | File:Line |
|------|----------|-----------|
| Cancellation ingress | `WorkflowHandler.RequestCancelWorkflowExecution` validates and forwards to `historyClient.RequestCancelWorkflowExecution` | `service/frontend/workflow_handler.go:2261-2293` |
| Cancellation transport proto | `RequestCancelWorkflowExecutionRequest{namespace_id, cancel_request, external_initiated_event_id, external_workflow_execution, child_workflow_only}` | `proto/internal/temporal/server/api/historyservice/v1/request_response.proto:566-576` |
| History handler → engine | `Handler.RequestCancelWorkflowExecution` resolves shard + engine; `historyEngineImpl.RequestCancelWorkflowExecution` invokes `requestcancelworkflow.Invoke` | `service/history/handler.go:934-963`, `service/history/history_engine.go:679-684` |
| Durable intent write | `requestcancelworkflow.Invoke`: noop-if-closed (`IsWorkflowExecutionRunning`), chain guard (`FirstExecutionRunId`), parent guard (`childWorkflowOnly`), idempotent noop (`IsCancelRequested`), else `AddWorkflowExecutionCancelRequestedEvent` + `UpdateWorkflowWithNewWorkflowTask` | `service/history/api/requestcancelworkflow/api.go:14-96` |
| Mutable-state intent | `AddWorkflowExecutionCancelRequestedEvent` (dup reject, builder + apply, persist `CancelRequestId`); `Apply...: CancelRequested=true`; `IsCancelRequested()` | `service/history/workflow/mutable_state_impl.go:5091-5118`, `service/history/workflow/mutable_state_impl.go:5120-5125`, `service/history/workflow/mutable_state_impl.go:2510-2511` |
| CancelRequested event factory | `CreateWorkflowExecutionCancelRequestedEvent`: type 20, `Cause=Reason`, `Identity`, `ExternalInitiatedEventId`, `ExternalWorkflowExecution`, `Links` | `service/history/historybuilder/event_factory.go:586-598` |
| Runtime ack (cancel command) | `COMMAND_TYPE_CANCEL_WORKFLOW_EXECUTION → handleCommandCancelWorkflow`: fail on buffered events, validate attrs, drop if closed, `AddWorkflowExecutionCanceledEvent` | `service/history/api/respondworkflowtaskcompleted/workflow_task_completed_handler.go:303-304`, `service/history/api/respondworkflowtaskcompleted/workflow_task_completed_handler.go:951-979` |
| Terminal close | `AddWorkflowExecutionCanceledEvent` + `Apply...: STATE_COMPLETED/STATUS_CANCELED, ClearStickyTaskQueue, processCloseCallbacks, GenerateWorkflowCloseTasks`; factory type 21 | `service/history/workflow/mutable_state_impl.go:5127-5167`, `service/history/historybuilder/event_factory.go:600-612` |
| Activity propagation | `handleCommandRequestCancelActivity`: `AddActivityTaskCancelRequestedEvent`, immediate-canceled if not started, else enqueue `WorkerCommand_CancelActivity` via `flushWorkerCommandsTasks` | `service/history/api/respondworkflowtaskcompleted/workflow_task_completed_handler.go:680-779` |
| Activity cancel flag | `Add/ApplyActivityTaskCancelRequestedEvent (ai.CancelRequested=true)`, `Add/ApplyActivityTaskCanceledEvent` (requires flag, `DeleteActivity`) | `service/history/workflow/mutable_state_impl.go:4736-4779`, `service/history/workflow/mutable_state_impl.go:4860-4946` |
| Activity heartbeat poll | `cancelRequested=ai.CancelRequested` returned in `RecordActivityTaskHeartbeatResponse` | `service/history/api/recordactivitytaskheartbeat/api.go:77-101` |
| Activity ack | `RespondActivityTaskCanceled` requires `ai.CancelRequested` else `ErrActivityTaskNotCancelRequested`, then `AddActivityTaskCanceledEvent` + new workflow task | `service/history/api/respondactivitytaskcanceled/api.go:86-108` |
| Child/external propagation | `handleCommandRequestCancelExternalWorkflow` → `AddRequestCancelExternalWorkflowExecutionInitiatedEvent` + `GenerateRequestCancelExternalTasks` | `service/history/api/respondworkflowtaskcompleted/workflow_task_completed_handler.go:981-1019`, `service/history/workflow/mutable_state_impl.go:5169-5192` |
| Transfer-queue child delivery | `processCancelExecution` → `requestCancelExternalExecution` rebuilds nested `historyservice.RequestCancelWorkflowExecutionRequest` and calls `historyRawClient.RequestCancelWorkflowExecution`; source records `ExternalWorkflowExecutionCancelRequested` (24) or `...Failed` (23) | `service/history/transfer_queue_active_task_executor.go:532-634`, `service/history/transfer_queue_active_task_executor.go:1612-1640` |
| Parent-close-policy cancel | `PARENT_CLOSE_POLICY_REQUEST_CANCEL → RequestCancelWorkflowExecution{FirstExecutionRunId, ExternalWorkflowExecution:Parent, ChildWorkflowOnly}` | `service/worker/parentclosepolicy/workflow.go:122-135` |
| Context/token propagation | Workflow-task token (`clock`), activity-task token, `RequestId/CancelRequestId` dedupe, `FirstExecutionRunId/ChildWorkflowOnly/ExternalInitiatedEventId` attribution | `common/tasktoken/token.go:9-61`, `proto/internal/temporal/server/api/token/v1/message.proto:39`, `service/history/api/requestcancelworkflow/api.go:26-34` |
| Signal handling | `InterruptCh`: `signal.Notify(c, os.Interrupt, SIGTERM)` only; `InterruptOn` → `blockingStart`; `ServerFx.Start` blocks then `Stop()` | `temporal/interrupt.go:9-21`, `temporal/server_option.go:76-82`, `temporal/fx.go:349-363`, `cmd/server/main.go:226-243` |
| Shutdown order/timeout | `initOrder{matching:1,history:2,frontend:3,worker:4}`, stop reversed; `serviceStartTimeout=15s`, `serviceStopTimeout=5m`; `ServicesMetadata.Stop` wraps `WithTimeout(ctx,5m)` | `temporal/server_impl.go:47-53`, `temporal/server_impl.go:109-146`, `temporal/server.go:11-12`, `temporal/fx.go:370-377` |
| Per-service stop | Frontend: drain sleep + `GracefulStop` + `AfterFunc(requestDrainTime→Stop)`; history/matching: drain sleep + `GracefulStop` + `AfterFunc(2s→Stop)`; worker: scanner/per-ns/worker-manager/parent-close stops + `GracefulStop`, no fallback timer | `service/frontend/service.go:548-595`, `service/history/service.go:117-176`, `service/matching/service.go:87-125`, `service/worker/service.go:305-323` |
| Drain deadline defaults | `FrontendShutdownDrainDuration 0s`, `FrontendShutdownFailHealthCheckDuration 0s`, `MatchingShutdownDrainDuration 0s`, `HistoryShutdownDrainDuration 0s` | `common/dynamicconfig/constants.go:869-878`, `common/dynamicconfig/constants.go:1425-1429`, `common/dynamicconfig/constants.go:1918-1922` |
| Cleanup contexts | `newDetachedContext: CopyContextValues + WithTimeout` ("won't be affected if base is cancelled"); `verifychildworkflowcompletionrecorded` detached resend because "gRPC cancels when handler returns"; long-poll `ctx.Done` fallback | `service/history/shard/context_impl.go:2408-2433`, `service/history/api/verifychildworkflowcompletionrecorded/api.go:159-162`, `service/history/api/get_workflow_util.go:248-258` |
| Detached workflow cleanup (SDK-side) | In-server only `workflow.NewDisconnectedContext` in migration handover/lifecycle workflows; `DetachedCancellationScope` lives in SDK repo, not vendored | `service/worker/migration/handover_workflow.go:141`, `service/worker/migration/workflow_lifecycle_events.go:82` |
| Process-tree termination | No evidence found: grep `Setpgid|Kill\(|SIGKILL|cgroup|process group` returns zero hits in `sources/temporal/*.go` | `(no file — absence verified)` |
| Race arbitration | Cancel-vs-close: cancel noops on closed run while terminate errors; same-WFT `Complete/Fail/Cancel` first-writer-wins (`MultipleCompletionCommandsCounter`); `Start+RequestCancel` child same WFT rejected | `service/history/api/requestcancelworkflow/api.go:46-53`, `service/history/api/terminateworkflow/api.go:47-49`, `service/history/api/respondworkflowtaskcompleted/workflow_task_completed_handler.go:855-979`, `service/history/api/command_attr_validator.go:260-295` |
| Cancel tests | `TestExternalRequestCancelWorkflowExecution` (idempotent double-cancel, exact history), `..._TargetRunning/TargetFinished/TargetNotFound`, `TestImmediateChildCancellation_WorkflowTaskFailed`, `TestCancelWhilePaused`, worker-command cancel dispatch | `tests/cancel_workflow_test.go:33-670`, `tests/pause_workflow_execution_test.go:1959-2030`, `tests/worker_commands_task_test.go:614-874` |
| Transfer/persistence cancel tests | `TestProcessCancelExecution_{Success,Failure,TargetNamespaceNotFound,Duplication}`, `NewHistoryExecutionRequestCancelSuite`, `TestTransferRequestCancelTask`, `TestRefreshRequestCancelExternalTasks` | `service/history/transfer_queue_active_task_executor_test.go:1575-1750`, `common/persistence/sql/sqlplugin/tests/history_execution_request_cancel.go:31`, `common/persistence/serialization/task_serializers_test.go:89`, `service/history/workflow/task_refresher_test.go:1148` |
| Terminate vs Reset vs Cancel | Terminate force-closes (`FORCE_CLOSE_COMMAND` + `Terminated`); Reset replays to `FinishEventId-1` with new run ID + reapply policy | `service/history/workflow/util.go:98-125`, `service/history/api/resetworkflow/api.go:28-184`, `service/frontend/workflow_handler.go:2261-2467` |

## Answers to Dimension Questions

1. **Is client disconnect different from explicit cancellation?** Yes. Disconnect is a transport event with no workflow effect: long-poll helpers return gracefully on `ctx.Done` (`service/history/api/get_workflow_util.go:248-258`, `service/history/api/polltimeskipping/api.go:155-164`), and async post-handler work is explicitly detached because "gRPC cancels when handler returns" (`service/history/api/verifychildworkflowcompletionrecorded/api.go:159-162`). Explicit cancellation is a durable history write (`service/history/api/requestcancelworkflow/api.go:82`) that survives disconnect and requires SDK ack to close. No `Disconnect` symbol exists in the daemon shutdown path.
2. **Is cancellation durable if no worker is currently attached?** Yes. `requestcancelworkflow.Invoke` persists `CancelRequested` and unconditionally returns `UpdateWorkflowWithNewWorkflowTask` (`service/history/api/requestcancelworkflow/api.go:86`), scheduling a workflow task into matching that waits for a future poller. Paused-workflow test proves durability across quiescence: cancel while paused writes `CANCEL_REQUESTED` with no task dispatched until `Unpause` (`tests/pause_workflow_execution_test.go:1997-2030`).
3. **Can cleanup hang indefinitely?** Bounded at the daemon layer, cooperative at the workflow layer. Daemon stops are bounded by `serviceStopTimeout=5m` (`temporal/server.go:12`, `temporal/fx.go:371`), gRPC `AfterFunc(2s→Stop)` fallbacks for history/matching (`service/history/service.go:164-168`), drain sleeps (`service/history/service.go:117-176`), and `WithTimeout` cleanup contexts (`service/history/shard/context_impl.go:2408-2433`, `temporal/fx.go:1075-1123`). Workflow-level cleanup can wait indefinitely if the worker never returns the `CancelWorkflowExecution` command — cancellation is cooperative, and the workflow stays open in `CancelRequested` state until the SDK acks (`service/history/api/respondworkflowtaskcompleted/workflow_task_completed_handler.go:951-979`). Compensation itself is an SDK pattern with no server timeout primitive (no saga primitive found).
4. **Can child/grandchild processes escape termination?** No OS processes exist to escape — Temporal manages goroutines/workflows, not a process tree (zero hits for process-group/cgroup/kill). At the workflow-tree level, escape is possible by design: child cancel is a best-effort transfer task (`service/history/transfer_queue_active_task_executor.go:1612-1640`) that fails open to `REQUEST_CANCEL_EXTERNAL...FAILED` on not-found/namespace errors (`service/history/transfer_queue_active_task_executor.go:603-622`, `tests/cancel_workflow_test.go:404-456`); a finished target drops the cancel (`tests/cancel_workflow_test.go:267-365`); `Start+RequestCancel` in the same WFT is rejected rather than forced (`service/history/api/command_attr_validator.go:260-295`). Grandchild propagation depends on each intermediate workflow's own cancel handler.
5. **How are cancellation and completion races arbitrated?** First-writer-wins on mutable state plus idempotent-cancel semantics. Late cancel to a closed run is a success-noop (`service/history/api/requestcancelworkflow/api.go:46-53`), whereas late terminate errors (`service/history/api/terminateworkflow/api.go:47-49`). Duplicate cancel is a noop (`service/history/api/requestcancelworkflow/api.go:73-80`). Competing `Complete/Fail/Cancel` commands in one workflow task are arbitrated in `workflow_task_completed_handler.go:855-979` — buffered-event check first, then `!IsWorkflowExecutionRunning` → metric + warn + drop ("just pick first one"). `checkMutability` (`service/history/workflow/mutable_state_impl.go:9116-9130`) rejects post-close writes with `ErrWorkflowFinished`.

**Scenario — cancel arrives exactly as work completes and a child process ignores the first signal:** the terminal outcome is whichever close event (`Canceled` vs `Completed`/`Failed`/`Terminated`) wins the `IsWorkflowExecutionRunning` race; the loser is dropped (late cancel → success-noop per `service/history/api/requestcancelworkflow/api.go:46-53`; same-task multi-complete → first-wins per `workflow_task_completed_handler.go:855-979`). Cleanup for the child is one transfer-queue `CancelExecution` attempt (`service/history/transfer_queue_active_task_executor.go:1612-1640`) — there is no signal-escalation loop ("ignores the first signal" has no server-side retry/escalation; activity path similarly sends one `WorkerCommand_CancelActivity` at `workflow_task_completed_handler.go:751-760` and otherwise relies on heartbeat `CancelRequested` polling at `recordactivitytaskheartbeat/api.go:77`). If the child already finished, the cancel is dropped and recorded as `...FAILED` on the source (`event_factory.go:640-668`, `tests/cancel_workflow_test.go:404-456`).

## Architectural Decisions

- **Cancellation as durable event + scheduled workflow task** (`service/history/api/requestcancelworkflow/api.go:82-86`): decouples intent from delivery; no worker need be attached. Tradeoff: extra workflow-task churn per cancel.
- **Cooperative ack via SDK command** (`workflow_task_completed_handler.go:951-979`): workflow code observes `CancelRequested` and decides cleanup/compensation. Tradeoff: no forced preemption; cleanup duration is worker-controlled.
- **Idempotent, success-on-closed cancel** (`api.go:46-80`) vs **strict terminate** (`terminateworkflow/api.go:47-49`): cancel is safe to retry; terminate is the force-close escape hatch.
- **Transfer-queue fan-out for cross-workflow cancel** (`transfer_queue_active_task_executor.go:1612-1640`) with `RequestId` dedupe and `ChildWorkflowOnly` guard (`api.go:64-71`): at-least-once, parent-authenticated delivery. Tradeoff: async, fails open.
- **Heartbeat + eager worker-command dual channel for activities** (`recordactivitytaskheartbeat/api.go:77`, `workflow_task_completed_handler.go:731-760`): fast path plus poll fallback. Tradeoff: pre-feature/backoff activities skip eager cancel (`workflow_task_completed_handler.go:718-728`).
- **Reverse-`initOrder` drain-then-`GracefulStop`** (`temporal/server_impl.go:47-53,109-124`): worker stops before frontend, matching/history drain membership first (`SetDraining/EvictSelf`). Tradeoff: correctness depends on non-zero drain durations that default to `0s`.
- **Detached contexts for post-cancel/cleanup persistence** (`service/history/shard/context_impl.go:2408-2433`): shard-ownership and resend work survive caller cancellation. Tradeoff: must pair with explicit timeouts to avoid leaks.

## Notable Patterns

- `GetAndUpdateWorkflowWithNew` + `UpdateWorkflowAction{Noop, CreateWorkflowTask}` lease pattern for all cancel/terminate/reset mutations (`service/history/api/requestcancelworkflow/api.go:36-91`, `service/history/api/workflow_lease.go:19-29`).
- `IsWorkflowExecutionRunning` / `IsCancelRequested` / `checkMutability` guard triad (`mutable_state_impl.go:2510-2511,5099,9116-9130`).
- `hBuilder` event-factory + `Apply...` rebuilder symmetry for every cancel event (`historybuilder/event_factory.go:534-693`, `mutable_state_rebuilder.go:524-591`).
- `RequestId`/`CancelRequestId` end-to-end idempotency across client → history → child (`api.go:26-34`, `transfer_queue_active_task_executor.go:1620-1636`, `mutable_state_impl.go:5116`).
- fx `StartStopHook`/`StopHook` cascade for deterministic daemon teardown (`temporal/fx.go:941-946`, `common/resource/fx.go:274-289`, `common/persistence/client/fx.go:222-239`).
- Membership `SetDraining` → sleep → `EvictSelf` → `GracefulStop` → `Stop()` fallback sequence (`service/frontend/service.go:548-595`, `service/history/service.go:117-176`).

## Tradeoffs

- Durability vs latency: every cancel pays a history write plus a workflow-task round trip, even for already-cancelled/closed runs (mitigated by noop fast paths).
- Cooperativeness vs boundedness: SDK-controlled cleanup enables compensation patterns but means `CancelRequested` can linger if workers stall; server provides no cancel deadline.
- Async child fan-out vs atomicity: transfer-queue delivery scales across shards/namespaces but introduces races (target-finished drop) and partial fan-out.
- Drain-by-default-off (`0s` defaults): fast restarts out of the box, but rolling deploys must explicitly configure `ShutdownDrainDuration` family or lose in-flight graceful behavior.
- No OS-level kill chain: correct for a goroutine-per-task server, but operators expecting cgroup/process-group guarantees must enforce them outside the binary (container runtime).

## Failure Modes / Edge Cases

- Cancel lands after close → silent success-noop; callers cannot distinguish "cancelled" from "already finished" without reading history (`api.go:46-53`).
- `FirstExecutionRunId` mismatch → `ErrWorkflowExecutionNotFound` (`api.go:60-62`); wrong-parent child-only cancel → `ErrWorkflowParent` (`api.go:64-71`).
- Buffered events at ack time → workflow task failed `UNHANDLED_COMMAND`, cancel deferred (`workflow_task_completed_handler.go:955`).
- Activity never heartbeats and eager command unsupported (backoff/pre-feature) → cancel visible only on next state transition (`workflow_task_completed_handler.go:718-728`).
- Child namespace deleted/unreachable → `...FAILED` event on source, no retry beyond queue retry policy (`transfer_queue_active_task_executor.go:603-622`, `..._test.go:1693-1750`).
- Daemon SIGKILL or `Stop()`-fallback path → in-flight `GracefulStop` abandoned; history/matching fall back after 2s but worker has no fallback (`service/worker/service.go:305-323`).
- `0s` drain defaults + immediate membership evict → pollers may observe `NOT_SERVING` before draining; long-poll clients get graceful `ctx.Done` returns rather than errors (`get_workflow_util.go:248-258`).
- Detached contexts without timeouts would leak; in-tree uses pair them with `WithTimeout` (1s OTEL, 30s namespace init, 5m service stop).

## Future Considerations

- Add a server-side cancel-to-close watchdog (optional `CancelRequested`-age timeout emitting `TimedOut` or force-`Terminated`) for stalled cooperative cleanup.
- Document and lint `ShutdownDrainDuration` non-zero in production overlays; consider raising non-zero defaults or warning on `0s` with multiple services per process.
- Add a `Stop()`-fallback timer to `service/worker/service.go:305-323` symmetric with history/matching 2s `AfterFunc`.
- Expose cancel-fan-out observability (transfer `CancelExecution` retry/failed counters, `CancelRequested`-to-`Canceled` latency histogram) tied to `MultipleCompletionCommandsCounter` race metric.
- Clarify client API docs: cancel-success-on-closed vs terminate-error-on-closed, and `FirstExecutionRunId` chain semantics.

## Questions / Gaps

- No evidence found for OS process-group/cgroup termination or kill-escalation deadlines in the server binary (search boundary: `Setpgid|Kill\(|SIGKILL|cgroup` over `sources/temporal/**/*.go` returned zero hits; only `SIGINT/SIGTERM` handled in `temporal/interrupt.go:9-21`). Operator-level containment presumably delegated to container runtime — not verified in-repo.
- No evidence found for a bounded cleanup deadline after workflow cancellation (no `CancelRequested`-TTL or activity-cancel timeout in `mutable_state_impl.go` or dynamicconfig); cleanup boundedness rests on SDK cooperation.
- `CancelRequestId` is persisted in `executionInfo` (`mutable_state_impl.go:5116`) but not in the history event payload (`event_factory.go:586-598`) — replication/visibility of the dedupe key beyond the issuing shard not fully traced.
- SDK `DetachedCancellationScope` semantics referenced by in-server `NewDisconnectedContext` comments but defined in the SDK repo outside this source boundary.

---
Generated by `Dimension 01.09: Cancellation, Shutdown, and Process Cleanup` against `temporal`.
