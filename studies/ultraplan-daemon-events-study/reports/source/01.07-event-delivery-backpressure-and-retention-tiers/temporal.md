# Source Analysis: temporal

## Event Delivery, Backpressure, and Retention Tiers

### Source Info

| Field | Value |
|-------|-------|
| Name | temporal |
| Path | `studies/ultraplan-daemon-events-study/sources/temporal` |
| Language / Stack | Go (history/frontend/matching/worker services, Cassandra/SQL persistence, gRPC long-poll) |
| Analyzed | 2026-09-03 |

## Summary

Temporal splits events into two explicit tiers: a durable history tier (workflow history nodes written transactionally through shard-gated `AppendHistoryNodes`, plus durable transfer/timer/replication/visibility task queues) and ephemeral notification tiers (history-event fan-out `service/history/events/notifier.go:25,120` with 1000-slot ingress queue + per-subscriber size-1 channels, and generic long-poll pub/sub `service/history/notification/notifier.go:15-25` with buffered-1 drop-newest semantics). Hot-path durable writes happen inside the mutable-state transaction before fan-out (`service/history/workflow/transaction_impl.go:678,718` notifies only after commit), so subscribers never block persistence. High-volume pressure is shed by dropping/coalescing notifications, `TrySubmit`-gated queue readers, and monitor/mitigator slice-shrinking, while storage is bounded by per-namespace retention timers (`service/history/workflow/task_generator.go:341-367`), a 4-stage ordered delete protocol (`service/history/shard/context_impl.go:910-935`), reference-counted branch pruning (`common/persistence/history_manager.go:131-211`), history size/count forced-termination limits (`service/history/workflow/context.go:1466-1522`), and buffered-event caps (`service/history/workflow/mutable_state_impl.go:8977-8981`).

## Rating

**8/10 — Clear model with tests, explicit interfaces, and operational safeguards.**

Rationale: durability vs. ephemerality is explicit in types (`HistoryMutation`, `EventStore` db vs. mem buffers, `DeleteHistoryEventTask`), drop policies are documented in code comments with counters/gauges, queue backpressure has a closed sensor→alert→mitigator loop with rate limiters, and retention/archival/deletion is a staged, version-checked, jittered protocol with metrics. Downgraded from 9 because (a) notification drops are counter-only with no replay sequence, (b) retention delete across clusters/visibility/history is eventually consistent with orphan-branch GC as backstop, and (c) buffered-event and history-size limits fail by rejecting/terminating work rather than shedding load gracefully.

## Evidence Collected

Every entry MUST include a file path with line numbers. Format: `path/to/file.ts:NN`.

| Area | Evidence | File:Line |
|------|----------|-----------|
| Durable vs. buffered event split (EventStore struct) | `EventStore{dbBufferBatch, memBufferBatch, memEventsBatches, memLatestBatch}` — sealed batches are persistable, buffer batches carry placeholder `BufferedEventID` | `service/history/historybuilder/event_store.go:15-39` |
| Buffer-routing policy | `add()` routes by `bufferEvent(type)`; `bufferEvent()` allowlist: state-change / WFT / command-derived / update events return `false` (immediately durable), `PAUSED/UNPAUSED/TIME_SKIPPING_TRANSITIONED`/default return `true` | `service/history/historybuilder/event_store.go:85-106,297-360` |
| Buffer promotion to durable | `FlushBufferToCurrentBatch()` concats db+mem buffers, `reorderBuffer()`, allocates sequential EventIDs, appends to latest batch; no-op if empty, dropped if workflow finished | `service/history/historybuilder/event_store.go:161-200` |
| Batch sealing + task-ID assignment | `Finish()/FlushAndCreateNewBatch()` rolls `memLatestBatch` into `memEventsBatches`, calls `taskIDGenerator` per event; `appendToLatestBatch()` rolls batch when `proto.Size` exceeds `maxEventBatchSizeInBytes()` | `service/history/historybuilder/event_store.go:120-132,202-283` |
| Buffer flush gates | `FlushBufferedEvents()` gated on `!HasStartedWorkflowTask()`; WFT-schedule path flushes before `WorkflowTaskScheduled` to avoid transient/speculative IDs | `service/history/workflow/mutable_state_impl.go:1140-1145`, `service/history/workflow/workflow_task_state_machine.go:327-333` |
| Durable write validation + batching | `serializeAppendHistoryNodesRequest()` rejects empty events, `EventId<1`, mixed versions, non-continuous IDs, enforces `transactionSizeLimit()`; `AppendHistoryNodes/AppendRawHistoryNodes` write `InternalAppendHistoryNodesRequest{BranchToken, Node{NodeID, Events, PrevTxn, Txn}}` | `common/persistence/history_manager.go:325-411,482-515` |
| Shard-gated durable append (hot path) | `ShardContext.AppendHistoryEvents()` checks `errorByState()`, injects ShardID, calls execution manager, emits per-namespace `HistorySize` metric + size-threshold warn | `service/history/shard/context_impl.go:871-908` |
| Storage durability (Cassandra/SQL) | Cassandra `AppendHistoryNodes`: new branch = `LoggedBatch(InsertTree+UpsertNode)`, else `UpsertHistoryNode`; SQL counterpart identical interface; wrapped by rate-limited/retry/metric clients | `common/persistence/cassandra/history_store.go:64-108`, `common/persistence/persistence_rate_limited_clients.go:762`, `common/persistence/persistence_retryable_clients.go:446`, `common/persistence/persistence_metric_clients.go:913` |
| Post-commit fan-out (durable first) | `transaction_impl.go` calls `engine.NotifyNewHistoryEvent(events.NewNotification(...))` only after commit at two commit sites — subscribers never block the write | `service/history/workflow/transaction_impl.go:678,718` |
| History-event notifier: ingress drop + per-subscriber coalesce | `eventsChan` cap 1000 (`eventsChanSize`); `enqueueHistoryEventNotification` sets timestamp, non-blocking send else `HistoryEventNotificationFailDeliveryCount++`; single-goroutine dequeue emits `InFlightMessageGauge` + `QueueingLatency`; `dispatchHistoryEventNotification` per-workflow fan-out with `select ch<-e default` on size-1 subscriber chans ("should NOT happen" = drop under load) | `service/history/events/notifier.go:25,120,180-229` |
| Generic long-poll pub/sub: drop-newest + subscriber cap | `Notify` fans out synchronously on publisher goroutine with non-blocking buffered-1 sends (full = drop new, old wins); `Watch` caps `maxSubscribersPerKey`, excess returns `ResourceExhausted/CONCURRENT_LIMIT` | `service/history/notification/notifier.go:15-25,83-108,131-148` |
| Time-skipping waiter bound | `maxFastForwardWaitersPerExecution=5` via `NewPubSubNotifier`; documented drop edge falls back to `POLL_TIMEOUT` re-poll | `service/history/notification/timeskipping_notifier.go:28-34,81-87` |
| Long-poll slow-consumer release | `GetOrPollWorkflowMutableState`: `WatchHistoryEvent` + re-check loop with `WithDeadlineBuffer(LongPollExpirationInterval)`, ignores stale versions, branch-change → `CurrentBranchChanged`; history API re-queries on empty token race, nils token on non-running to terminate | `service/history/api/get_workflow_util.go:127-134,190-260`, `service/history/api/getworkflowexecutionhistory/api.go:220-258,405-407,500-503` |
| Durable queue backpressure (non-blocking submit) | Reader `loadAndSubmitTasks`: `ratelimiter.Wait()` then `if !scheduler.TrySubmit(exe)` shed + `verifyPendingTaskSize (pending < MaxPendingTasksCount)` else `Pause(PollBackoffInterval)`; `TrySubmit` contract never blocks; Nack path falls back to `rescheduler.Add` | `service/history/queues/reader.go:384-403,427-463,527-533`, `service/history/queues/scheduler.go:46`, `service/history/queues/rescheduler.go:231` |
| Queue monitor → alert → mitigator loop | `SetSlicePendingTaskCount/SetSliceCount/SetReaderWatermark` aggregate pending/slices, emit `Alert{QueuePendingTaskCount/SliceCount/ReaderStuck}` (deduped, `alertChSize=10` non-blocking); `handleAlert` → `mitigator.Mitigate` (shrink/split-and-clear + 10s `Pause`) + checkpoint + `notifyReaders` | `service/history/queues/monitor.go:133-153,167-219,243-263`, `service/history/queues/alerts.go:7-39`, `service/history/queues/action_pending_task_count.go:52-72,159-202`, `service/history/queues/mitigator.go:58-93`, `service/history/queues/queue_base.go:433-445` |
| In-memory scheduled queue handoff | `newTaskCh` unbuffered — producer blocks until `processQueueLoop` receives; `!TrySubmit` → `go{ Visibility+=1s workerBusyRescheduleDelay; Add() }` retry in new goroutine (no drop); factory uses `FIFOScheduler{QueueSize:0}` to push backpressure to queue | `service/history/queues/memory_scheduled_queue.go:27,54,101-103,152-159`, `service/history/memory_scheduled_queue_factory.go:61-67` |
| Buffered-event growth bound | `BufferSizeAcceptable()` rejects tx when `NumBufferedEvents > MaximumBufferedEventsBatch()` or bytes > `MaximumBufferedEventsSizeInBytes()`; batch roll at `MaximumEventBatchSizeInBytes`; history/queue configs in history config + dynamicconfig constants | `service/history/workflow/mutable_state_impl.go:8977-8981`, `service/history/configs/config.go:210-211,215,681-682,686` |
| History size/count circuit breaker | `maxHistorySizeExceeded/maxHistoryCountExceeded` compare `ExecutionStats.HistorySize` / `NextEventID-1` against per-namespace warn+error limits; error + running → `forceTerminateWorkflow(HistorySize/CountExceedsLimit)` | `service/history/workflow/context.go:1466-1522` |
| Retention timer generation | `GenerateDeleteHistoryEventTask(closeTime)`: `retention=namespace.Retention()` (fallback default), `deleteTime=closeTime+retention+FullJitter(RetentionTimerJitterDuration)`, emits in-memory `DeleteHistoryEventTask{WorkflowKey, VisibilityTimestamp:deleteTime, Version:closeVersion, BranchToken}` (caller must `AddTasks`) | `service/history/workflow/task_generator.go:282-295,341-367` |
| Retention execution + resurrection guard | `executeDeleteHistoryEventTask`: skips if state != COMPLETED, `CheckTaskVersion`, delegates to `deleteManager.DeleteWorkflowExecutionByRetention`; `NotFound` → `deleteHistoryBranch` direct | `service/history/timer_queue_task_executor_base.go:81-161` |
| Ordered retention delete protocol | `ShardContext.DeleteWorkflowExecution` 4 stages: 1 visibility+replication, 2 current pointer, 3 mutable state, 4 history branch (must be last; stage-3 loss orphans branch for background GC); `DeleteWorkflowExecutionByRetention` marks replication-processed so clusters delete independently | `service/history/shard/context_impl.go:910-935`, `service/history/deletemanager/delete_manager.go:139-216` |
| Branch pruning / compaction (ref-counted) | `DeleteHistoryBranch` loads full tree, builds `usedBranches` from siblings/ancestors, deletes bottom-up via `InternalDeleteHistoryBranchRange` to spare shared fork ancestors; `TrimHistoryBranch` paged-reads + selective `DeleteHistoryNodes` (best-effort, log-only failure) | `common/persistence/history_manager.go:131-224`, `common/persistence/cassandra/history_store.go:111-137,268-286` |
| Archival-then-delete chain | `Archive()` rate-limits (`WaitN` else `ResourceExhausted`), fans out history+visibility archivers, `multierr.Combine`; archival queue executor archives then regenerates delete task; drops task on `ErrMutableStateIsNil`/`ErrWorkflowExecutionIsStillRunning` | `service/history/archival/archiver.go:107-174`, `service/history/archival_queue_task_executor.go:96-109,222-256` |
| Metrics: lag / depth / WAL-equivalent growth | `history_event_notification_{queueing_latency,fanout_latency,inflight_message_gauge,fail_delivery_count}`, `shardinfo_{immediate,scheduled}_queue_lag/backlog_age`, `pending_tasks/task_batch_complete/task_rescheduler_pending/task_scheduler_throttled/queue_reader/slice_count/actions/alert_shadow` | `common/metrics/metric_defs.go:811-814,826-837,949-962`, `service/history/queues/queue_base.go:331-346`, `service/history/events/notifier.go:202-229` |
| Configs bounding growth | `EventsShardLevelCacheMaxSizeBytes 512KB` / `EventsHostLevelCacheMaxSizeBytes 256MB` + `EventsCacheTTL` (LRU, single-event `PageSize:1` fetch); `QueueMaxPredicateSize`, `OutboundQueue*`, per-queue `TaskBatchSize/MaxPollRPS/Interval`, replication `EnableReplicationTaskBatching/ReplicationStreamReadBufferSize` + stream sender QPS flow controller | `service/history/configs/config.go:94-98,123,151-196,249-255,331-340`, `service/history/events/cache.go:56-64,84-89,185`, `service/history/replication/stream_sender_flow_controller.go:45-59` |

## Answers to Dimension Questions

1. **Which events must never be dropped?**
   Workflow history nodes (`common/persistence/history_manager.go:482-515` → Cassandra/SQL `AppendHistoryNodes`) and durable queue tasks (transfer/timer/replication/visibility/archival via `service/history/queues/` + `service/history/queue_factory_base.go:79-106`). They are written inside the shard-gated transaction (`service/history/shard/context_impl.go:871-908`) before any notification (`service/history/workflow/transaction_impl.go:678,718`). Close lifecycle is additionally pinned by a `DeleteHistoryEventTask` carrying `BranchToken+closeVersion` (`service/history/workflow/task_generator.go:341-367`). No drop path exists for these — pressure surfaces as transaction rejection (`transactionSizeLimit`, `BufferSizeAcceptable`, history-size termination) rather than silent loss.

2. **Which events are ephemeral or short-lived?**
   History-event notifications (`service/history/events/notifier.go:202-229`: 1000-slot `eventsChan` ingress drop + size-1 per-subscriber coalesce with `FailDeliveryCount`), generic long-poll wakeups (`service/history/notification/notifier.go:131-148`: buffered-1 drop-newest, old wins), time-skipping fast-forward signals (5 waiters max, documented `POLL_TIMEOUT` re-poll at `service/history/notification/timeskipping_notifier.go:28-34`), and in-queue scheduling retries (`service/history/queues/memory_scheduled_queue.go:152-159` delayed re-`Add`, `service/history/queues/rescheduler.go:231` Nack buffer). All require the waiter to re-read authoritative state on wake/timeout.

3. **Can a slow UI stall execution?**
   No. Durable appends do not wait for subscribers: `Notify` is non-blocking (`service/history/notification/notifier.go:137-144`), history dispatch drops rather than blocks (`service/history/events/notifier.go:180-212`), long-polls expire via `LongPollExpirationInterval` (`service/history/api/get_workflow_util.go:127-134`), queue readers shed via `TrySubmit` + `MaxPendingTasksCount` pause (`service/history/queues/reader.go:527-533`), and queue schedulers are `TrySubmit`-only (`service/history/queues/scheduler.go:46`). A stalled UI only loses its own wakeups and re-polls. The one blocking handoff is internal: unbuffered `newTaskCh` in `memory_scheduled_queue.go:101-103` blocks the producer goroutine until the loop receives — but that is shard-internal, and busy-scheduler overflow is re-scheduled in a fresh goroutine, not dropped.

4. **How are lifecycle and terminal facts flushed under pressure?**
   Terminal facts take the durable path, never the notification path: close writes seal `EventStore` batches with task IDs (`service/history/historybuilder/event_store.go:202-283`), commit via `AppendHistoryNodes`, then emit `DeleteHistoryEventTask` with jittered `VisibilityTimestamp` (`service/history/workflow/task_generator.go:341-367`). Under pressure, `BufferSizeAcceptable` (`service/history/workflow/mutable_state_impl.go:8977-8981`), `transactionSizeLimit` (`common/persistence/history_manager.go:325-411`), and history size/count error limits (`service/history/workflow/context.go:1466-1522`) reject or force-terminate the transaction instead of flushing partial state. Retention execution is version-checked and resurrection-guarded (`service/history/timer_queue_task_executor_base.go:81-161`), archival is rate-limited before delete (`service/history/archival/archiver.go:107-174`), and the 4-stage delete order guarantees a workflow is never left accessible-but-historyless (`service/history/shard/context_impl.go:910-935`).

5. **What bounds storage and in-memory growth?**
   Storage: per-namespace retention + jitter (`service/history/workflow/task_generator.go:341-367`), archival-then-delete (`service/history/archival_queue_task_executor.go:222-256`), ref-counted branch delete + selective trim (`common/persistence/history_manager.go:131-224`), history size/count forced termination (`service/history/workflow/context.go:1466-1522`), event cache caps 512KB shard / 256MB host + TTL (`service/history/configs/config.go:94-98`). In-memory: history notifier 1000 + per-subscriber 1 (`service/history/events/notifier.go:25,131`), pub/sub buffered-1 + per-key waiter cap (`service/history/notification/notifier.go:84,93-96`), `MaxPendingTasksCount` scaled by `0.8^readerID` (`service/history/queues/queue_base.go:151-159`), `QueueMaxPredicateSize` (`service/history/configs/config.go:123`), per-queue batch/poll/RPS limits and replication stream read-buffer + QPS limiters (`service/history/configs/config.go:151-196,331-340`, `service/history/replication/stream_sender_flow_controller.go:45-59`).

## Architectural Decisions

- **Durable-first, notify-after-commit.** Transaction commit precedes `NotifyNewHistoryEvent` (`service/history/workflow/transaction_impl.go:678,718`), so notification loss never implies data loss.
- **Buffered vs. immediate event classes.** Command/state-change events bypass the buffer for ordering (`service/history/historybuilder/event_store.go:297-360`); signals/markers/etc. buffer until the next WFT boundary (`service/history/workflow/mutable_state_impl.go:1140-1145`).
- **Lossy wakeups + authoritative re-read.** Both notifiers document best-effort delivery and push recovery to the waiter (re-read mutable state, re-poll on timeout) — `service/history/notification/notifier.go:15-25`, `service/history/events/notifier.go:202-229`.
- **Sensor→alert→mitigator queue control.** Pending-count/slice-count/reader-stuck sensors feed a central mitigator that shrinks, splits-and-clears, and pauses readers (`service/history/queues/monitor.go:133-263`, `service/history/queues/action_pending_task_count.go:159-202`).
- **Jittered, versioned, staged retention.** Retention timers carry full jitter, close-version checks, and a 4-stage visibility→current→mutable→branch delete order (`service/history/workflow/task_generator.go:341-367`, `service/history/shard/context_impl.go:910-935`).
- **Ref-counted history branches.** Fork-shared ancestor nodes are spared during branch delete; trim is a separate best-effort path (`common/persistence/history_manager.go:131-224`).

## Notable Patterns

- **Non-blocking fan-out with counters:** `select ch<-v default: metric++` at both ingress (1000) and per-subscriber (1) stages of the history notifier.
- **Old-wins coalescing:** full buffered-1 channels keep the stale value and drop the newer one, forcing waiters to re-read — the opposite of drop-oldest, chosen because any wake triggers a full state reload.
- **Per-key waiter admission control:** `maxSubscribersPerKey` → `ResourceExhausted` instead of unbounded waiter growth.
- **`TrySubmit`-only scheduling:** no scheduler in the history queue stack blocks; overflow stays in the slice or moves to the rescheduler.
- **Decay-scaled reader quotas:** higher `readerID` (overflow groups) gets `0.8^readerID` of the critical pending budget, containing fan-out cost of slow groups.
- **Single-goroutine notifier pump with latency/depth gauges:** dequeue loop emits `InFlightMessageGauge` and `QueueingLatency` on every event.

## Tradeoffs

- **Durability over latency on the write path:** every history append pays serialization + persistence round-trip inside the shard transaction; notifications add no latency because they are post-commit and lossy.
- **Drop-newest vs. drop-oldest:** keeping the stale wakeup is cheaper (one re-read covers all missed versions) but a long-stalled waiter can act on arbitrarily old version hints until it re-checks (`service/history/api/get_workflow_util.go:219` stale-version ignore).
- **Reject/terminate instead of shed:** buffered-event caps and history size/count limits protect storage by failing transactions or killing workflows — safe but user-visible, unlike silent sampling.
- **Unbuffered internal handoff + async retry:** `memory_scheduled_queue.go:101-103` risks producer parking, mitigated by spawning the busy-retry in a new goroutine (goroutine churn under sustained saturation).
- **Eventual retention consistency:** per-cluster independent deletes + `MarkProcessed(StageReplication)` avoid cross-cluster coordination but leave windows where one cluster retains what another deleted; orphan-branch GC is the backstop.

## Failure Modes / Edge Cases

- **Never-reading subscriber under flood (dimension probe):** ingress `eventsChan` (1000) fills first → `HistoryEventNotificationFailDeliveryCount` increments and new notifications are dropped; per-subscriber size-1 channels coalesce so the parked subscriber holds one stale event while execution continues committing durable history. First degradation: notification loss counter + rising `InFlightMessageGauge`/`QueueingLatency`; durable history, queue tasks, and retention timers remain intact. No evidence of execution stall was found — all fan-out sends are non-blocking.
- **Notifier pump stall/crash:** single `dequeueHistoryEventNotifications` goroutine is the sole consumer; if it blocks in `dispatchHistoryEventNotification` (it does not block on sends, but `GetAndDo` takes map locks), `eventsChan` fills and everything drops to counters.
- **Retention resurrection race:** close-then-reopen with same version could race the timer; guarded by `state != COMPLETED` check + `CheckTaskVersion` (`service/history/timer_queue_task_executor_base.go:81-148`), with `NotFound` falling back to branch-only delete.
- **Stage-3 delete loss orphans history branch:** explicitly documented at `service/history/shard/context_impl.go:910-935` — orphan relies on background GC rather than inline recovery.
- **Trim failure is log-only:** `DeleteHistoryNodes` best-effort path can leave trimmed nodes behind with no retry counter beyond logs (see `common/persistence/history_manager.go:214ff` + `execution_manager_test.go:146-224` expectations).
- **Archival rate-exhaustion blocks delete chain:** `Archive()` returns `ResourceExhausted` when the limiter is saturated, delaying the subsequent `GenerateDeleteHistoryEventTask` and extending storage retention under archival pressure.

## Future Considerations

- Add notification sequence numbers so waiters can detect how many versions were coalesced/dropped between wakeups (currently only stale-version ignore exists).
- Expose per-workflow fan-out depth and per-key waiter-count gauges for the generic pub/sub (currently only history-notifier gauges/counters exist).
- Consider drop-oldest-with-latest-version-hint for the history notifier so a stalled subscriber wakes with the freshest version instead of the oldest.
- Add retention-delete lag/overdue metrics (visibility timestamp vs. actual delete time) and archival-queue depth to complement existing queue lag gauges.
- Evaluate backpressure propagation from `BufferSizeAcceptable` failures to the frontend (throttle vs. reject) to avoid client retry storms when buffered-event caps trip.

## Questions / Gaps

- No evidence found for history-node compaction beyond branch delete/trim — searched `common/persistence/`, `service/history/historybuilder/`, `schema/` for compaction jobs; only ref-counted delete and selective node delete exist.
- No evidence found for a spill-to-disk or persistent overflow for the 1000-slot history notification queue — searched `service/history/events/`; overflow is counter-only drop.
- No evidence found for visibility-store retention enforcement details beyond timer-driven workflow deletion — visibility delete is stage 1 of `DeleteWorkflowExecution` but its own TTL/GC policy was not traced in this pass.
- Open: what are the production default values for `MaximumBufferedEventsBatch/SizeInBytes`, `PendingTasksCriticalCount`, and `RetentionTimerJitterDuration`? Config wires dynamicconfig keys (`service/history/configs/config.go:608,681-682`) but defaults live in `common/dynamicconfig/constants.go` and were not fully enumerated here.
- Open: does the replication stream flow controller (`service/history/replication/stream_sender_flow_controller.go:45-59`) apply per-shard or per-cluster QPS, and how does it interact with `ReplicationStreamReadBufferSize` under catch-up floods?

---
Generated by `Dimension 01.07: Event Delivery, Backpressure, and Retention Tiers` against `temporal`.
