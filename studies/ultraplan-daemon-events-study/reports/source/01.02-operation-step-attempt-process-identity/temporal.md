# Source Analysis: temporal

## Operation, Step, Attempt, and Process Identity

### Source Info

| Field | Value |
|-------|-------|
| Name | temporal |
| Path | `studies/ultraplan-daemon-events-study/sources/temporal` |
| Language / Stack | Go 1.24 (`go.mod:3`), protobuf, gRPC, persistence (SQL/Cassandra), history/matching/frontend services, CHASM state machines |
| Analyzed | 2026-09-02 |

## Summary

Temporal models durable execution as a 4-level identity lattice: **Workflow (durable intent) → Run (single attempt of that intent) → Step (history-event-addressed unit: WorkflowTask/Activity/Timer/ChildExecution/Signal/etc.) → Attempt (retry instance within a step)**. Each level has explicit IDs, monotonic counters, and fencing tokens, making the counterfactual “attempt 2 succeeds after attempt 1 returns late” decidable by several independent coordinates.

The top-level intention is a **WorkflowExecution** identified by `(NamespaceId, WorkflowId)` (`proto/internal/temporal/server/api/persistence/v1/executions.proto:57-58`). A single workflow may spawn a chain of **Runs** each with a new `RunId` (`proto/internal/temporal/server/api/persistence/v1/executions.proto:382-384`), carrying `Attempt` (`executions.proto:122`), `FirstExecutionRunId` / `OriginalExecutionRunId` (`executions.proto:149,270`), and `RootWorkflowId/RootRunId` (`executions.proto:227-228`) plus full `Parent` fields (`executions.proto:59-62`, `service/history/workflow/retry.go:240-258`). Retries (`SetupNewWorkflowForRetryOrCron` at `service/history/workflow/retry.go:271`) increment `previousExecutionInfo.Attempt + 1` (`retry.go:271`) and write it on the new `WorkflowExecutionStarted` event via `req.Attempt` (`retry.go:294`), while the old run's `NewExecutionRunId`/`SuccessorRunId` (`executions.proto:157,289`) chains it.

Within a Run, the **checkpointable sub-units** are history-event-identified entities stored in `MutableStateImpl` maps (`service/history/workflow/mutable_state_impl.go:128-155`): `pendingActivityInfoIDs` by `ScheduledEventId`, `pendingTimerInfoIDs` by `TimerId`, `pendingChildExecutionInfoIDs` by `InitiatedEventId`, plus `SignalInfo`, `RequestCancelInfo`, and CHASM `sub_state_machines_by_type` (`executions.proto:212`). Each entity’s state is journaled as a `HistoryEvent` in the branch identified by `VersionHistories` / `branch_token` (`executions.proto:147`, `mutable_state_impl.go:1011-1028`) and its transition is stamped with `VersionedTransition` (`executions.proto:208`) and `VersionHistoryItem (eventId, version)` (`proto/internal/temporal/server/api/history/v1/message.proto:19-21`).

**Attempt identity** is distinct at every level: workflow-level `WorkflowExecutionInfo.Attempt` (`executions.proto:122`), per-workflow-task `WorkflowTaskAttempt` (`executions.proto:85`) + `WorkflowTaskAttemptsSinceLastSuccess` (`executions.proto:114`) + `WorkflowTaskStamp` (`executions.proto:108`), per-activity `ActivityInfo.Attempt` (`executions.proto:615`) + `Stamp` (`executions.proto:673`), and per-callback/Nexus `Attempt` (`executions.proto:875,924`). Tokens that cross process boundaries encode these: `tokenspb.Task { Attempt, ActivityAttemptStamp }` (`proto/internal/temporal/server/api/token/v1/message.proto:44,55`) created via `tasktoken.NewWorkflowTaskToken(... attempt ...)` (`common/tasktoken/token.go:9-30`) and `NewActivityTaskToken(... attempt, activityAttemptStamp ...)` (`common/tasktoken/token.go:33-61`). History validates them with `Stamp` checks (`service/history/api/recordworkflowtaskstarted/api.go:78`, `service/history/api/recordactivitytaskstarted/api.go:97,181`, `service/history/api/isactivitytaskvalid/api.go:60`, `service/history/api/isworkflowtaskvalid/api.go:63`), so a stale completion with `attempt=1` after `attempt=2` is rejected as `ObsoleteMatchingTask` / `Activity task rejected; stamp has changed`.

One logical workflow routinely spans many processes: a Workflow Task is dispatched through `matching` TaskQueue partitions (`proto/internal/temporal/server/api/taskqueue/v1/message.proto:15-80`) to any worker polling (`Matching long-poll`), while Activities are dispatched independently to the activity queue; sidecar determinism is enforced server-side. External process identity is recorded as `StartedIdentity`, `TaskQueue`, `WorkerDeploymentVersion` (`executions.proto:692`, `activity.go:107-108`).

## Rating

**9 / 10 — Mature, durable, observable, extensible, and proven under failure or scale**

Rationale: Four-level hierarchy is explicit in protobuf, mutable-state structs, and history events, with monotonic attempt counters, stamps/clocks for fencing, immutable event attribution (`(namespace, workflowId, runId, scheduledEventId, attempt, stamp, version, branchToken)`), and exhaustive stale-task rejection paths tested in `history_engine_test.go` and `mutable_state_impl_test.go`. Deductions only for the historical dual-track (VersionHistories vs. TransitionHistory) that temporarily duplicates versioning checks and the `use_compatible_version` / `assigned_build_id` deprecation surface (`executions.proto:633-650`).

## Evidence Collected

Every entry MUST include a file path with line numbers. Format: `path/to/file.ts:NN`.

| Area | Evidence | File:Line |
|------|----------|-----------|
| Run/job/task/attempt models — workflow-level durable intent | `WorkflowExecutionInfo.namespace_id/workflow_id` are the primary durable key; every other field hangs off it | `proto/internal/temporal/server/api/persistence/v1/executions.proto:57-58` |
| Run/job/task/attempt models — run identity | `WorkflowExecutionState.create_request_id/run_id/state/status/start_time/request_ids/first_execution_run_id` — one row per attempt (Run) of the workflow | `proto/internal/temporal/server/api/persistence/v1/executions.proto:381-400` |
| Run/job/task/attempt models — workflow retry chain | `attempt` + `has_retry_policy/retry_*` on workflow info; `new_execution_run_id` / `successor_run_id` chain continuations; `first_execution_run_id` / `original_execution_run_id` preserve origin | `proto/internal/temporal/server/api/persistence/v1/executions.proto:122,36-39,157,289,149,270` |
| Run/job/task/attempt models — root/chain identity | `root_workflow_id/root_run_id` documented as “same as parent’s root” so whole tree shares a root | `proto/internal/temporal/server/api/persistence/v1/executions.proto:227-228` |
| Run/job/task/attempt models — parent nesting | `parent_namespace_id/parent_workflow_id/parent_run_id/parent_initiated_id/parent_clock/parent_initiated_version` | `proto/internal/temporal/server/api/persistence/v1/executions.proto:59-62,62,160-161` |
| Run/job/task/attempt models — child workflow step | `ChildExecutionInfo { version/initiated_event_batch_id/started_event_id/started_workflow_id/started_run_id/create_request_id/namespace/workflow_type_name/parent_close_policy/... priority }` keyed by `initiated_event_id` | `proto/internal/temporal/server/api/persistence/v1/executions.proto:761-779` |
| Run/job/task/attempt models — workflow task step | `workflow_task_version/scheduled_event_id/started_event_id/timeout/attempt/stamp/suggest_continue_as_new/history_size_bytes` | `proto/internal/temporal/server/api/persistence/v1/executions.proto:81-96,108-114` |
| Run/job/task/attempt models — activity step | `ActivityInfo { version/scheduled_event_batch_id/scheduled_time/started_event_id/request_id/activity_id/activity_type/attempt/task_queue/started_identity/has_retry_policy/... stamp/paused/priority/start_version/worker_control_task_queue/started_clock }` keyed by `scheduled_event_id` | `proto/internal/temporal/server/api/persistence/v1/executions.proto:592-748` |
| Run/job/task/attempt models — timer step | `TimerInfo { version/started_event_id/expiry_time/task_status/timer_id/last_update_versioned_transition }` keyed by `timer_id` | `proto/internal/temporal/server/api/persistence/v1/executions.proto:750-760` |
| Run/job/task/attempt models — signal/cancel steps | `SignalInfo { version/initiated_event_batch_id/request_id/initiated_event_id }` and `RequestCancelInfo { version/initiated_event_batch_id/cancel_request_id/initiated_event_id }` — each by `initiated_event_id` | `proto/internal/temporal/server/api/persistence/v1/executions.proto:790-802,781-789` |
| Run/job/task/attempt models — Nexus/callback steps | `NexusOperationInfo { endpoint/service/operation/request_id/state/attempt/last_attempt_failure/next_attempt_schedule_time/... schedule_to_close_timeout }` and `CallbackInfo { state/attempt/request_id ... }` | `proto/internal/temporal/server/api/persistence/v1/executions.proto:888-950,855-886` |
| Run/job/task/attempt models — mutable-state holder | `MutableStateImpl { pendingActivityInfoIDs/updateActivityInfos/deleteActivityInfos, pendingTimerInfoIDs, pendingChildExecutionInfoIDs, pendingSignalInfoIDs, pendingSignalRequestedIDs, executionInfo, executionState, hBuilder, chasmTree, stateMachineNode, InsertTasks ... }` | `service/history/workflow/mutable_state_impl.go:128-286` |
| Run/job/task/attempt models — initialization | `NewMutableState` initializes `WorkflowTaskAttempt: 1` and `WorkflowTask* = EmptyEventID/EmptyVersion` | `service/history/workflow/mutable_state_impl.go:382-392` |
| Run/job/task/attempt models — CHASM extensible steps | `sub_state_machines_by_type map<string,StateMachineMap>` plus `state_machine_timers` and `sub_state_machine_tombstone_batches` — HSM/CHASM steps are first-class | `proto/internal/temporal/server/api/persistence/v1/executions.proto:212,235,253` |
| Database schemas — shard / execution columns | `ShardInfo { shard_id/range_id/owner/stolen_since_renew/update_time/queue_states }` and `WorkflowMutableState = ExecutionInfo + ExecutionState + ActivityInfos + TimerInfos + ChildExecutionInfos + SignalInfos + RequestCancelInfos + ChasmNodes + SignalRequestedIds + Checksum + NextEventId + BufferedEvents` | `proto/internal/temporal/server/api/persistence/v1/executions.proto:33-53`, `service/history/workflow/mutable_state_impl.go:984-998` |
| Database schemas — transition/branch identity | `VersionedTransition { namespace_failover_version, transition_count }` + `transition_history[]` and `VersionHistories { current_version_history_index, histories[] }` with `VersionHistory { branch_token, items[] }` and `VersionHistoryItem { event_id, version }` | `proto/internal/temporal/server/api/persistence/v1/executions.proto:208`, `proto/internal/temporal/server/api/history/v1/message.proto:19-34` |
| Parent/root ID fields | Workflow parent/root propagation via `ParentExecutionInfo { namespaceId, namespace, execution{workflowId,runId}, initiatedId, clock, initiatedVersion }` and `RootExecutionInfo { execution }` built in `SetupNewWorkflowForRetryOrCron` | `proto/internal/temporal/server/api/workflow/v1/message.proto:5-20`, `service/history/workflow/retry.go:240-258` |
| Parent/root ID fields — runtime | `GetWorkflowKey() -> (namespaceId, workflowId, runId)` is the canonical lookup key; `GetCurrentBranchToken()` + `getCurrentBranchTokenAndEventVersion(eventId)` resolve history branch per event version | `service/history/workflow/mutable_state_impl.go:1003-1029` |
| Process/task handles — task queues | `TaskQueuePartition { task_queue, task_queue_type, partition_id { normal_partition_id/sticky_name } }` and `TaskForwardInfo/TaskVersionDirective` for versioned routing | `proto/internal/temporal/server/api/taskqueue/v1/message.proto:60-92,100-120` |
| Process/task handles — worker polling identity | `PhysicalTaskQueueInfo { pollers[] PollerInfo, internal_task_queue_status, task_queue_stats }` tracks poller build-ids / identities per queue | `proto/internal/temporal/server/api/taskqueue/v1/message.proto:40-58` |
| Process/task handles — task token | `tokenspb.Task { namespace_id/workflow_id/run_id/scheduled_event_id/attempt/activity_id/activity_type/clock/started_event_id/version/started_time/start_version/component_ref/activity_attempt_stamp }` — complete fencing token | `proto/internal/temporal/server/api/token/v1/message.proto:39-56` |
| Process/task handles — token factories | `NewWorkflowTaskToken(namespaceID, workflowID, runID, scheduledEventID, startedEventID, startedTime, attempt, clock, version)` and `NewActivityTaskToken(... attempt, activityAttemptStamp ...)` | `common/tasktoken/token.go:9-31,33-61` |
| Process/task handles — activity dispatcher identity | `ActivityInfo.task_queue/started_identity/last_worker_deployment_version/priority/worker_control_task_queue/started_clock` — which worker/process is bound | `proto/internal/temporal/server/api/persistence/v1/executions.proto:629-748` |
| Correlation and causation metadata | `WorkflowExecutionState.request_ids map<string, RequestIDInfo { event_type, event_id, attach_time }>` — every creation/attach is deduped by request ID | `proto/internal/temporal/server/api/persistence/v1/executions.proto:394,402-408` |
| Correlation and causation metadata — links | `Callback { Nexus{url,header} | HSM{namespace_id/workflow_id/run_id/ref/method}, links[] }` + `HSMCompletionCallbackArg { namespace_id/workflow_id/run_id, last_event }` and `Link.WorkflowEvent { namespace/workflowId/runId, reference { EventRef{eventId,eventType} | RequestIdRef{requestId,eventType} } }` | `proto/internal/temporal/server/api/persistence/v1/executions.proto:810-854`, `service/history/workflow/mutable_state_impl.go:799-813,886-912` |
| Correlation and causation metadata — causation clock | `clockspb.VectorClock` on `parent_clock`, `ActivityInfo.started_clock`, `ChildExecutionInfo.clock`, `TimerTaskInfo.branch_token/version/schedule_attempt/event_id`, `NexusOperationInfo.request_id/scheduled_event_token` | `proto/internal/temporal/server/api/persistence/v1/executions.proto:62,746-747,775,503-541,889-904` |
| Retry-attempt creation code — workflow retry | `attempt = previousExecutionInfo.Attempt + 1` when `CONTINUE_AS_NEW_INITIATOR_RETRY`, else `1`; passed as `historyservice.StartWorkflowExecutionRequest.Attempt` and written to `WorkflowExecutionStarted` + first WFT | `service/history/workflow/retry.go:271,294,292-330` |
| Retry-attempt creation code — workflow task | `WorkflowTaskAttempt` incremented on schedule/failure/timeout; `HasStartedWorkflowTask()`, `workflowTaskManager`, `AddFirstWorkflowTaskScheduled(parentClock, event, false)`; transient vs. speculative distinction via `WorkflowTaskAttempt > 1` (`ndc_task_util.go:272`) | `service/history/workflow/mutable_state_impl.go:3394,2408,9574-9575`, `service/history/workflow/retry.go:326-330`, `service/history/ndc_task_util.go:272` |
| Retry-attempt creation code — activity | `UpdateActivityInfoForRetries(ai, version, attempt, failure, nextScheduledTime, isStampEnabled)` sets `ai.Attempt = attempt`, clears `StartedEventId/RequestId/StartedTime/StartedClock`, masks `TimerTaskStatus`, and bumps `Stamp` when `attempt > previousAttempt` | `service/history/workflow/activity.go:59-90` |
| Retry-attempt creation code — activity reschedule | `GetNextScheduledTime` computes `ExponentialBackoffAlgorithm(initial, coeff, Attempt)` bounded by `RetryMaximumInterval`; `ResetActivity` resets `Attempt=1` and may regenerate retry task | `service/history/workflow/activity.go:222-243,315-378` |
| Retry-attempt creation code — activity visible attempt | `GetPendingActivityInfo` exposes `Attempt / NextAttemptScheduleTime / CurrentRetryInterval / LastFailure / PauseInfo` for UI/CLI | `service/history/workflow/activity.go:92-220` |

## Answers to Dimension Questions

### 1. What is the durable user intention?

The durable intention is a **Workflow Execution** named by `(NamespaceId, WorkflowId)` (`proto/internal/temporal/server/api/persistence/v1/executions.proto:57-58`). Workflow is the only entity that survives indefinitely from the user’s perspective; all other entities (Runs, Activities, Timers) are subordinate. Idempotent creation is keyed by `WorkflowExecutionState.create_request_id` (`executions.proto:382`) and retained in `request_ids` (`executions.proto:394`) so a retried `StartWorkflowExecutionRequest.RequestId` maps to the same `RunId` without creating a duplicate workflow. User-visible continuation intents (retry, cron, continue-as-new) produce the next Run in the chain but keep the same `WorkflowId` and a new `RunId`.

### 2. What is the smallest checkpointable unit?

The smallest unit is a **single HistoryEvent** within the ordered event journal of one Run, identified by `(RunId, EventId, Version, BranchToken, VersionedTransition)` (`proto/internal/temporal/server/api/history/v1/message.proto:19-22`, `executions.proto:147,192-208`). At the mutable-state layer the coarser but operationally relevant checkpoint is a **state-machine transition**: one transaction that mutates a single subordinate map entry (e.g., one `ActivityInfo` at `scheduledEventId`, one `TimerInfo` at `timerId`, one `ChildExecutionInfo` at `initiatedEventId`, one CHASM node) and appends one or more history events via `hBuilder` (`mutable_state_impl.go:414-423`). Activity attempts are sub-checkpoints: `LastAttemptCompleteTime` / `LastHeartbeatDetails` are persisted per attempt (`executions.proto:667-668`) so a crash after `ActivityTaskFailed` can reschedule the next attempt without losing the prior outcome. WorkflowTask completion is the workflow-code checkpoint: commands emitted by the worker become history events atomically with mutable-state updates in `InsertTasks` / `BestEffortDeleteTasks` (`mutable_state_impl.go:229-256`).

### 3. Does each retry get a distinct identity?

**Yes, at all three levels.**

- **Workflow retry:** new `RunId = uuid.NewString()` (`retry.go:294-createRequest.RequestId`), new `Attempt = previous.Attempt + 1` (`retry.go:271`), new `WorkflowExecutionStarted` event carrying that attempt (`retry.go:294,330-350`). The old run stores `NewExecutionRunId / SuccessorRunId` pointing forward, the new run stores `FirstExecutionRunId / OriginalExecutionRunId` pointing back (`executions.proto:55,99,157,289,399`).
- **Workflow-task retry:** `WorkflowExecutionInfo.workflow_task_attempt` (`executions.proto:85`) increments per schedule; `workflow_task_stamp` (`executions.proto:108`) disambiguates spurious completions. History writes `Attempt` into `WorkflowTaskScheduledEventAttributes` via `historybuilder.CreateWorkflowTaskScheduledEvent(taskQueue, timeout, attempt, scheduleTime)` (`service/history/historybuilder/event_factory.go:109-117`, `history_builder.go:192-195`).
- **Activity retry:** `ActivityInfo.attempt` (`executions.proto:615`) is set to the caller-supplied `attempt` (`activity.go:68`), not auto-incremented, so server controls the count; `Stamp` (`executions.proto:673`) is bumped iff `attempt > previousAttempt` (`activity.go:87-89`). The next attempt’s `scheduledTime` is computed from retry policy (`activity.go:230-240`) and dispatched as a new matching task. Paused activities are excluded from retry (`activityPendingRetry` returns false if `Paused`, `activity.go:41-45`). Reset APIs can force `Attempt=1` (`activity.go:319,403`).

Thus attempt 2 and attempt 1 have different `RunId` (workflow), different `ScheduledEventId+Attempt+Stamp` (activity/WFT), and different `request_id / branch_token / version history item`.

### 4. Can one logical operation span multiple runtime calls or OS processes?

**Yes, by design.**

- A single Workflow (`WorkflowId`) spans many **Runs** if it retries / continues-as-new; each Run is a distinct mutable-state row but logically one workflow.
- A single Run spans many **Workflow Task** RPCs: each `PollWorkflowTaskQueue` (matching) → `RespondWorkflowTaskCompleted` (history) cycle runs on an arbitrary worker process chosen by `TaskQueuePartition` / `PhysicalTaskQueueInfo` (`taskqueue/message.proto:60-80,40-58`). The `Task` token (`token/message.proto:39-56`) names the exact workflow task so the correct Run can be targeted without session affinity.
- A single Run fans out to many **Activity** RPCs on potentially different task queues and deployments; each activity’s `StartedIdentity / LastWorkerDeploymentVersion / WorkerControlTaskQueue` (`executions.proto:629,692,732`) records which worker holds that attempt.
- Server internals also span processes: the frontend validates `RequestId`, history owns the mutable state shard (`ShardInfo.owner/range_id` `executions.proto:33-36`), matching owns the backlog, and `TimerTaskInfo / TransferTaskInfo / ReplicationTaskInfo / CallbackInfo / NexusOperationInfo` each flow through their own queue category (`executions.proto:411-589`). All of them carry `(namespace_id, workflow_id, run_id, version, clock, branch_token)` so any host can continue.

Correlation across processes uses `NamespaceId/WorkflowId/RunId` plus `ScheduledEventId` / `InitiatedEventId` and `Link` causation fields; ordering within a Run is established by `EventId` + `VersionHistories` + `VersionedTransition` (`executions.proto:147,208`, `history/message.proto:19-22`).

### 5. Can events be unambiguously attributed to the right entity?

**Yes.** Four coordinates disambiguate:

1. **Workflow/Run:** `(NamespaceId, WorkflowId, RunId)` — appears on every persistence row and token (`token/message.proto:40-42`, `executions.proto:57-58,383`).
2. **Step:** `(EventId, EventType)` or the pending-map key (`ScheduledEventId` for activities, `InitiatedEventId` for child workflows, `TimerId` for timers) plus the `HistoryEvent`’s `EventId/Version` and `RequestIdRef/EventRef` link (`token/message.proto:39-56`, `executions.proto:762-779`).
3. **Attempt:** `Attempt (int32)` on `WorkflowExecutionInfo`, `WorkflowTaskScheduled`, `ActivityInfo`, `CallbackInfo`, `NexusOperationInfo` (`executions.proto:122,85,615,875,924`); `ActivityAttemptStamp` on the token (`token/message.proto:55`) and `WorkflowExecutionInfo.workflow_task_stamp` (`executions.proto:108`) break ties when attempt counters could overlap.
4. **Causation/version fence:** `VectorClock`, `VersionHistories / VersionHistoryItem`, `VersionedTransition`, `BranchToken`, and `Stamp` together invalidate stale readers/writers after failover or reset (`executions.proto:62,208,147`, `history/message.proto:19-22`, `mutable_state_impl.go:1003-1029`).

Concrete quarantine for the counterfactual: if workflow retry attempt 2 succeeds while attempt 1’s late response arrives, the system has **at least four independent rejectors**:

- **RunId mismatch:** attempt 1’s token / request targets `RunId_A`; the current run is `RunId_B`; `GetAndUpdateWorkflowWithNew` loads by `WorkflowKey(namespaceId, workflowId, runId)` (`isactivitytaskvalid/api.go:24`, `mutable_state_impl.go:1003`) and fails the lookup.
- **Workflow `Attempt` field:** even if WorkflowId-only indexing were used, `executionInfo.Attempt` (`executions.proto:122`) is 2 on the new run but the stale completion still carries `Attempt=1` (`retry.go:294`); dedup via `request_ids` (`executions.proto:394`) or `first_execution_run_id` chain resolves to the wrong Run.
- **Activity/WFT Stamp+Attempt:** `isActivityTaskValid` requires `ai.StartedEventId == EmptyEventID && ai.Stamp == request.Stamp` (`isactivitytaskvalid/api.go:60`); `recordactivitytaskstarted/api.go:97,181` and `recordworkflowtaskstarted/api.go:78` reject with `Activity task with this stamp not found` / `ObsoleteMatchingTask (stamp mismatch)`. A retry increments `Stamp` (`activity.go:87-89`), so late attempt 1 vs. attempt 2 fails the equality.
- **Version/Clock:** `VersionHistoryItem.version` + `VersionedTransition.transition_count` on the late event’s branch does not match the current branch’s `GetCurrentBranchToken()` / `CurrentVersionedTransition` (`mutable_state_impl.go:1011-1029`, `executions.proto:208`), so the write is discarded as referencing stale state, and `CheckSum` verification can also catch it (`mutable_state_impl.go:557-570`, `executions.proto:804-807`).

## Architectural Decisions

| Decision | Evidence | Consequence |
|----------|----------|-------------|
| **WorkflowId invariant, RunId per attempt** — logical intent stable (`WorkflowId`) while each retry/run gets a UUID `RunId` and monotonic `Attempt` | `executions.proto:57-58,122,382-383`, `retry.go:271,294`, `executions.proto:149,270,157,289` | Clients reference intent by `WorkflowId`; history can GC or query any attempt by `RunId` without ambiguity; retries never mutate the same row. |
| **Event-Id-addressed pending maps in MutableState** — every step keyed by the `ScheduledEventId`/`InitiatedEventId` assigned at scheduling time | `mutable_state_impl.go:130-155`, `executions.proto:592-802` | A step’s identity is stable before execution starts; cancellation, heartbeating, and retry target the exact event id. |
| **Monotonic Attempt + Stamp pair for fencing** — activities carry both `Attempt` and `Stamp`; workflow tasks carry `WorkflowTaskAttempt` + `WorkflowTaskStamp` | `executions.proto:615,673,85,108`, `activity.go:67-89`, `token/message.proto:44,55`, `token.go:33-61` | `Attempt` is user-observable (backoff, query); `Stamp` is the anti-ABA token that invalidates previously dispatched tasks without changing the user-visible count. |
| **VersionHistories + VersionedTransition dual fence** — branch token `+ (eventId, version)` for history forks and `VersionedTransition (failoverVersion, transitionCount)` for state-based tasks | `executions.proto:147,208`, `history/message.proto:19-34`, `mutable_state_impl.go:1003-1029,1202-1214` | Failover, NDC replication, and reset can create divergent histories; transition history subsumes the event-counter fence for HSM tasks. |
| **RequestId map on the execution row** — `WorkflowExecutionState.request_ids` keeps every creation/attach request id with its event type/id | `executions.proto:394,402-408` | Idempotent `StartWorkflowExecution` and CHASM `UpdateComponent` dedup without a separate table; request-id lookups stay local to the execution row. |
| **Tokenized cross-process dispatch** — `(namespace, workflow, run, scheduledEventId, attempt, stamp, clock, version)` packed into `tokenspb.Task` and handed to matching/workers | `token/message.proto:39-56`, `token.go:9-61` | No sticky affinity required; any frontend/matching/worker can resume; stale tokens are self-describing and fail closed. |

## Notable Patterns

- **Attempt-aware started-state clearing:** `ClearActivityStartedState` zeros `StartedEventId/StartVersion/RequestId/StartedTime/StartedClock` (`activity.go:51-57`) before reuse, so old attempt’s in-flight worker cannot “claim” the new attempt’s slot.
- **Transient vs. speculative WFT:** `IsWorkflowTaskTransient()` depends on `WorkflowTaskAttempt > 1` (`mutable_state_impl.go:2408`, `ndc_task_util.go:272`); transient WFT is not persisted until completion, speculative until validated, reducing churn on repeated failures.
- **Pause as first-class state, not a status bit:** `ActivityInfo.paused + PauseInfo { Manual{identity,reason} | ruleId, request_id }` (`executions.proto:676,698-720`) and `activityPendingRetry` short-circuit (`activity.go:41-45`) freeze retry dispatch without dropping the step.
- **Links for causation without session stickiness:** `Callback.links[]` / `HSMCompletionCallbackArg.last_event` (`executions.proto:841,855-853`) plus `Link.WorkflowEvent_RequestIdRef` vs. `EventRef` (`mutable_state_impl.go:799-912`) let callbacks be completed even after reset and by a different process, because the link names the cause event, not the process.
- **Stamp-gated task validation pattern** repeated for both WFT and Activity: `isActivityTaskValid` (`isactivitytaskvalid/api.go:60`) and `isWorkflowTaskValid` (`isworkflowtaskvalid/api.go:63`) compare the incoming token’s `Stamp` against mutable state; on mismatch return non-retryable fencing errors surfaced to `serviceerror.NewObsoleteMatchingTask`.

## Tradeoffs

- **Three versioning columns for backward compatibility:** both `VersionHistories` and `TransitionHistory` must be maintained and reconciled (`mutable_state_impl.go:1148-1170,1216-1244`). Enables rolling-upgrade but adds sanity-check burden and double-branch validation.
- **Wide mutable-state row:** `WorkflowMutableState` proto carries six maps + `ChasmNodes + BufferedEvents + Checksum` (`mutable_state_impl.go:984-998`). Avoids joins, but large workflows (many pending activities) pay read-modify-write cost on every transaction and can hit blob-size limits (`ExecutionStats.history_size` `executions.proto:370`).
- **Stamp is invisible to SDKs:** callers see `Attempt` in `PendingActivityInfo` (`activity.go:120`) but not `Stamp`; debugging a late-attribute requires server logs that include `ActivityAttemptStamp`. Prevents SDK misuse but slightly obscures fencing diagnostics.
- **Sticky queues as an optimization, not identity:** `StickyTaskQueue` (`executions.proto:118-121`) is an affinity hint, not part of the identity lattice; clearing it loses no correctness but can increase latency after worker loss.
- **Root/Parent denormalization:** `ParentWorkflowId/RunId/RootWorkflowId/RootRunId` are copied onto every child/continuation (`executions.proto:59-62,227-228`, `retry.go:240-258`). Speeds up tree queries, but updates are not transactional across the whole tree—consistent view requires chasing `ParentExecutionInfo` links.

## Failure Modes / Edge Cases

- **Late attempt completion after retry success:** Quarantined by `RunId` → `WorkflowKey` lookup, `Attempt`/`FirstExecutionRunId` chain, `Stamp+Attempt` equality in token validation, and `VersionedTransition` staleness check (see attribution above). The late writer receives `NotFound` / `ObsoleteMatchingTask` / `Activity task with this stamp not found` and is asked to discard.
- **Split-brain branch divergence (NDC / failover):** `VersionHistoryItem.version` and `branch_token` disambiguate histories; `VersionedTransition` disambiguates HSM tasks. `SanitizeMutableState` + `DiscardUnknownProto` (`mutable_state_impl.go:618-621`) protect against schema skew when replaying from a peer.
- **Activity heartbeat arriving after pause/reset:** `ai.ResetHeartbeats` + `ActivityReset` flags cause `LastHeartbeatDetails=nil` on next attempt (`activity.go:80-83,322-323,405-407`); the heartbeat processor ignores updates for a `StartedEventId==EmptyEventID` (`isactivitytaskvalid/api.go:60`).
- **Workflow task attempt overflow / tight loop:** `ContinueAsNewMinBackoff(durationpb.New(backoffInterval))` (`retry.go:292`) enforces a floor on retry backoff; `WorkflowTaskAttemptsSinceLastSuccess` (`executions.proto:114`) surfaces stuck-loop detection to `ReportedProblems` search attribute.
- **Tombstone exhaustion on heavily retried workflows:** `sub_state_machine_tombstone_batches` (`executions.proto:253`) accumulates deleted CHASM nodes; `totalTombstones` (`mutable_state_impl.go:170`) is tracked so compaction can page them out—without compaction, history payload growth can breach limits.
- **Request-ID map growth:** `WorkflowExecutionState.request_ids` is unbounded in theory; server sweeps by `history.maximumRequestIDsPerExecution` and `requestIDMaxAge` (referenced in `executions.proto:394` comment) to evict oldest CHASM request IDs by `attach_time`.

## Future Considerations

- Retire the `VersionHistories` column once `TransitionHistory` is proven to fully subsume branch fencing, removing the `shouldVerifyChecksum` dual-path in `mutable_state_impl.go:557-570`.
- Promote `ActivityAttemptStamp`/`WorkflowTaskStamp` to the visibility / `PendingActivityInfo` API so operators can trace which attempt a poller holds without inspecting matching logs.
- Persist the full `tokenspb.Task` blob on the execution row (as an alternative to `StartedClock` + reconstruction) to make token revalidation immune to future format drift (noted in `executions.proto:743` comment trade-off).
- Consider a compact `children_initialized_post_reset_point` index for very wide trees to avoid scanning the entire `RootExecutionInfo` chain when re-attaching signals after reset.

## Questions / Gaps

- No direct inspection of the Cassandra/SQL DDL for the `executions` table was performed beyond the proto column model; e.g., secondary indexes or TTL handling for tombstone batches were not traced.
- The exact `IsActivityTaskValid` → `MaxRetries` backoff caps for `NexusOperationInfo` / `CallbackInfo` are configured via dynamic config (`service/history/hsm/callbacks/config.go:26,32`) whose wiring into `Matching` polling was not traced end-to-end.
- Whether `WorkflowIdReusePolicy` enforcement consults `FirstExecutionRunId` vs. `RequestIds` in all legacy paths was not fully traced; search boundary covered `service/history/api/startworkflow` and `create_workflow_util.go` only.

---

Generated by `dimensions/01.02-operation-step-attempt-process-identity.md` against `temporal`.
