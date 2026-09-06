# Source Analysis: nomad

## 01.12 Policy, Admission, Scheduling, and Fairness

### Source Info

| Field | Value |
|-------|-------|
| Name | nomad |
| Path | `studies/aren-go-runtime-study/sources/nomad` |
| Language / Stack | Go (Raft-replicated scheduler, memdb state, heap queues) |
| Analyzed | 2026-08-30 |

## Summary

Nomad enforces policy before effects via a synchronous admission chain at `Job.Register`/`Job.Validate` followed by asynchronous scheduling through `EvalBroker` → `Worker` → `Scheduler` → `PlanQueue` → Raft-serial `planApply`. Every state-mutating effect is gated: jobs must pass mutators/validators + Sentinel/ACL before an eval is created, and allocs materialize only after `evaluatePlan` and Raft `UpsertPlanResults`. Queuing is priority-ordered per scheduler type with per-job serialization (one ready + heap-pending), Nack-driven retries with compounding delays, and a dedicated `BlockedEvals` tracker that parks unschedulable evals until capacity/quota unblocks appear. Fairness is minimal: strict priority + random tie-break, no tenant quotas in OSS (enterprise stub), no queue-depth admission, no starvation-bound timers. Stale-plan handling is explicit via `RefreshIndex`/`SnapshotIndex` checks and optimistic pipelined apply with snapshot invalidation.

## Rating

**6 / 10** — Mature admission → broker → blocked → plan pipeline with observable pre-effect decisions and robust blocked-eval correctness, but fairness is priority-only (random shuffle is the only anti-starvation), bounded waiting is partial (delays + deliveryLimit, no global depth limit), and quota/resource governance is enterprise-gated and invisible in OSS code.

Rationale: strong on “no start before commit”, explicit Nack/retry ownership, and blocked-state lifecycle; weak on per-tenant isolation, queue bounds, and bounded wait visibility.

## Evidence Collected

Every entry MUST include a file path with line numbers. Format: `path/to/file.ts:NN`.

| Area | Evidence | File:Line |
|------|----------|-----------|
| Admission mutators | `NewJobEndpoints` wires `mutators: [jobCanonicalizer, jobVaultHook, jobConsulHook, jobConnectHook, jobExposeCheckHook, jobImpliedConstraints, jobNodePoolMutatingHook, jobImplicitIdentitiesHook, jobNumaHook]` | `studies/aren-go-runtime-study/sources/nomad/nomad/job_endpoint.go:66-76` |
| Admission validators | `validators: [jobConnectHook, jobExposeCheckHook, jobVaultHook, jobConsulHook, jobNamespaceConstraintCheckHook, jobNodePoolValidatingHook, &jobValidate, &memoryOversubscriptionValidate, jobNumaHook, &jobSchedHook]` | `studies/aren-go-runtime-study/sources/nomad/nomad/job_endpoint.go:77-88` |
| Admission chain order | `admissionControllers` runs mutators then validators; `doRegister` calls `admissionControllers` before eval creation, then `enforceSubmitJob` (Sentinel CE) | `studies/aren-go-runtime-study/sources/nomad/nomad/job_endpoint_hooks.go:178-192` and `studies/aren-go-runtime-study/sources/nomad/nomad/job_endpoint.go:132-249` |
| Pre-register policy denial | ACL checks (`AllowNsOpAnyOf submitJob`, host volume, CSI), `PolicyOverride` sentinel override check, `registrationsAreAllowed` load-shedding, `validateJobUpdate`, `JobMaxCount`/`JobMaxPriority` checks | `studies/aren-go-runtime-study/sources/nomad/nomad/job_endpoint.go:115-211,382-396,517-551,1030-1050` and `studies/aren-go-runtime-study/sources/nomad/nomad/job_endpoint_hooks.go:517-551` |
| Load-shedding flag | `SchedulerConfiguration.RejectJobRegistration bool` gate | `studies/aren-go-runtime-study/sources/nomad/nomad/structs/operator.go:230-232` |
| Load-shedding enforcement | `registrationsAreAllowed` returns false when `RejectJobRegistration` set unless management ACL | `studies/aren-go-runtime-study/sources/nomad/nomad/job_endpoint.go:1383-1396` |
| Job → eval creation | `doRegister` builds `structs.Evaluation{Priority: job.Priority (or EvalPriority override), Type: job.Type, TriggeredBy: JobRegister}` only if not periodic/parameterized; otherwise raft `JobRegisterRequest` | `studies/aren-go-runtime-study/sources/nomad/nomad/job_endpoint.go:304-322,358-372` |
| Eval broker core | `EvalBroker` struct: `pending map[NamespacedID]PendingEvaluations`, `ready map[string]ReadyEvaluations`, `unack`, `requeue`, `delayHeap`, `enqueuedTime/dequeuedTime` capped 10k | `studies/aren-go-runtime-study/sources/nomad/nomad/eval_broker.go:53-122,337-339,405` |
| Enqueue dedup + waitUntil | `processEnqueue`: dedups `evals[ID]`, handles `Wait>0` via `time.AfterFunc`, `WaitUntil` via `delayHeap`, else `enqueueLocked` | `studies/aren-go-runtime-study/sources/nomad/nomad/eval_broker.go:254-296` |
| Per-job serialization | `enqueueLocked` checks `jobEvals[NamespacedID]`; first ID goes to `ready`, subsequent go to `pending` heap; `Ack` promotes next pending via `heap.Pop` and `MarkForCancel` | `studies/aren-go-runtime-study/sources/nomad/nomad/eval_broker.go:320-378,644-667` |
| Priority ordering | `ReadyEvaluations.Less` flipped heap: higher `Priority` wins, tie `CreateIndex`; `PendingEvaluations.Less` by priority then `ModifyIndex` | `studies/aren-go-runtime-study/sources/nomad/nomad/eval_broker.go:1040-1091` |
| Cross-scheduler fairness | `scanForSchedulers` collects highest-priority across schedulers, random `rand.Intn(n)` on tie | `studies/aren-go-runtime-study/sources/nomad/nomad/eval_broker.go:435-488` |
| Dequeue blocking | `Dequeue` + `waitForSchedulers` with per-scheduler `waiting chan` and goroutine fan-out | `studies/aren-go-runtime-study/sources/nomad/nomad/eval_broker.go:385-566` |
| Ack | Stops nack timer, checks `deliveryLimit`, cleans `unack/evals/jobEvals`, promotes pending, handles `requeue[token]`, updates stats, collects `cancelable` | `studies/aren-go-runtime-study/sources/nomad/nomad/eval_broker.go:599-675` |
| Nack & retry delay | `Nack` stops timer, re-enqueues to `failedQueue` if `dequeus>=deliveryLimit` else `Wait= nackReenqueueDelay(dequeues)` with `initialNackDelay` then compounding `subsequentNackDelay` | `studies/aren-go-runtime-study/sources/nomad/nomad/eval_broker.go:678-738` |
| Nack timeout pause | `PauseNackTimeout`/`ResumeNackTimeout` around `Plan.Submit` queue | `studies/aren-go-runtime-study/sources/nomad/nomad/eval_broker.go:740-772` and `studies/aren-go-runtime-study/sources/nomad/nomad/plan_endpoint.go:70-73` |
| Delivery limit & failed queue | `failedQueue="_failed"` sentinel, `deliveryLimit` param to `NewEvalBroker` | `studies/aren-go-runtime-study/sources/nomad/nomad/eval_broker.go:29,53-55` |
| Time metrics bound | `enqueuedTime`/`dequeuedTime` capped 10k entries | `studies/aren-go-runtime-study/sources/nomad/nomad/eval_broker.go:337-339,405` |
| Blocked evals tracker | `BlockedEvals` struct: `captured`, `escaped`, `system`, `jobs NamespacedID->ID`, `unblockIndexes`, channels `capacityChangeCh` (buf 8096), `unblockIndexes` pruning | `studies/aren-go-runtime-study/sources/nomad/nomad/blocked_evals.go:17-29,31-95,926-930` |
| Block path | `Block`/`Reblock` enqueue async `capacityUpdate` via `watchCapacity`; `processBlock` dedups per-job, checks `missedUnblock`, sets `ClassEligibility`, `EscapedComputedClass`, `QuotaLimitReached` | `studies/aren-go-runtime-study/sources/nomad/nomad/blocked_evals.go:176-287,289-343,359-417` |
| Missed unblock | `missedUnblock` compares `SnapshotIndex` vs `unblockIndexes[index]` for quota, class, missing resources, escaped max | `studies/aren-go-runtime-study/sources/nomad/nomad/blocked_evals.go:359-417` |
| Unblock triggers | `Unblock(computedClass)`, `UnblockQuota`, `UnblockClassAndQuota`, `UnblockNode`, `UnblockNonNodeResources` store index and send `capacityUpdate` | `studies/aren-go-runtime-study/sources/nomad/nomad/blocked_evals.go:484-644` |
| Unblock fan-out | `unblock()` unblocks escaped always on class change, resource-missing escapes, quota/class-eligible captured, node-specific system evals, then `EnqueueAll` | `studies/aren-go-runtime-study/sources/nomad/nomad/blocked_evals.go:674-759` |
| FSM→blocked wiring | Node ready/eligibility updates call `blockedEvals.Unblock*`; alloc terminal calls `UnblockClassAndQuota`; eval upsert routes to `Enqueue`/`Block`/`Untrack` | `studies/aren-go-runtime-study/sources/nomad/nomad/fsm.go:458,510-516,544-549,612-619,969-984,1067-1075` |
| Worker dequeue → scheduler | `worker.run` loops `dequeueEvaluation` → `snapshotMinIndex(waitIndex)` → `invokeScheduler` → `sendAck` on success else `sendNack` | `studies/aren-go-runtime-study/sources/nomad/nomad/worker.go:398-472,476-529` |
| Worker backoff & version gate | `dequeueEvaluation` backoff on `scheduler version mismatch` with `backoffSchedulerVersionMismatch=30s`; `backoffErr` exponential | `studies/aren-go-runtime-study/sources/nomad/nomad/worker.go:39-45,508-519,875-911` |
| Plan submission | `Worker.SubmitPlan` sets `EvalToken`, `SnapshotIndex`, `NormalizeAllocations`, RPC `Plan.Submit` with retry on `No cluster leader`/`plan queue is disabled` + backoff | `studies/aren-go-runtime-study/sources/nomad/nomad/worker.go:650-726,875-885` |
| Plan queue | `PlanQueue`: priority heap by `Plan.Priority` then `enqueueTime`, `Depth` stat, `Flush` on disable, bounded `waitCh` | `studies/aren-go-runtime-study/sources/nomad/nomad/plan_queue.go:32-40,99-186,232-238` |
| Plan apply pipeline | `planner.planApply` pipelined 6-8 Raft applies: drains `planIndexCh`, tracks `prevPlanResultIndex`/`inFlightPlans`, `snapshotMinIndex(max(prev, plan.SnapshotIndex))`, `evaluatePlan` | `studies/aren-go-runtime-study/sources/nomad/nomad/plan_apply.go:100-243,279-300` |
| Plan evaluation & stale check | `evaluatePlan`: `evaluatePlanQuota`→`evaluatePlanPlacements`; `RefreshIndex = max(RefreshIndex, AllocIndex)` for stale | `studies/aren-go-runtime-study/sources/nomad/nomad/plan_apply.go:542-560,482-512` |
| Quota stub OSS | `evaluatePlanQuota` returns false in OSS; ent overrides quota enforcement | `studies/aren-go-runtime-study/sources/nomad/nomad/plan_apply_ce.go:27-30` |
| Feasibility & placement | `GenericScheduler.process` → `computeJobAllocs` → `feasible.NewEvalContext`/`GenericStack` → `reconciler.NewAllocReconciler` → `computePlacements` with `Select` + preemption fallback | `studies/aren-go-runtime-study/sources/nomad/scheduler/generic_sched.go:204-330,485-722,910-935` |
| Feasibility wrapper eligibility cache | `FeasibilityWrapper.Next` caches `EvalComputedClassIneligible/Eligible/Escaped` to skip checks; tracks `QuotaLimitReached`, `MissingResources` | `studies/aren-go-runtime-study/sources/nomad/scheduler/feasible/feasible.go:1418-1554` |
| Resources fit check at apply | `evaluateNodePlan` calls `AllocsFit(node, proposed, nil, true)`; handles disconnected/down nodes | `studies/aren-go-runtime-study/sources/nomad/nomad/plan_apply.go:776-844` |
| Shuffling anti-starvation | `ShuffleNodes` seeds with `plan.EvalID` last 8 bytes xor `index` for Fisher-Yates before `StaticIterator` | `studies/aren-go-runtime-study/sources/nomad/scheduler/feasible/feasible.go:131-165` |
| Scheduler retry ownership | `GenericScheduler.Process` with `retryMax(limit, s.process, progressMade)`; on `SetStatusError` creates `blockedEval` with `EvalTriggerMaxPlans` and `setStatus` | `studies/aren-go-runtime-study/sources/nomad/scheduler/generic_sched.go:137-177,181-200` |
| Blocked eval creation | `createBlockedEval(planFailure)` stores `ClassEligibility`, `Escaped`, `QuotaLimitReached`, `MissingResources` | `studies/aren-go-runtime-study/sources/nomad/scheduler/generic_sched.go:181-200` |
| Delayed rescheduling | `GenericScheduler.process` creates `followUpEvals` with `WaitUntil` for reschedule tracker; `EvalBroker` delays via `delayHeap` | `studies/aren-go-runtime-study/sources/nomad/scheduler/generic_sched.go:272-286` and `studies/aren-go-runtime-study/sources/nomad/nomad/eval_broker.go:283-293,886-936` |
| Observable decisions | `SetStatus` via `ReblockEval`/`UpdateEval`/`CreateEval` → Raft `EvalUpdateRequestType`; eval `Status` (`pending/blocked/complete/failed`) + `FailedTGAllocs` + `AllocMetric.FilteredNodes` + `PlanAnnotations` | `studies/aren-go-runtime-study/sources/nomad/nomad/worker.go:728-870` and `studies/aren-go-runtime-study/sources/nomad/scheduler/generic_sched.go:161-177` |
| Observability metrics | `BrokerStats` (`TotalReady/Unacked/Pending/Waiting/Cancelable/ByScheduler/DelayedEvals`), `BlockedStats` (`TotalBlocked/Escaped/QuotaLimit`), `QueueStats.Depth`, `EmitStats` gauges | `studies/aren-go-runtime-study/sources/nomad/nomad/eval_broker.go:1019-1033,984-1016` and `studies/aren-go-runtime-study/sources/nomad/nomad/blocked_evals.go:867-905` and `studies/aren-go-runtime-study/sources/nomad/nomad/blocked_evals_stats.go:23-70` and `studies/aren-go-runtime-study/sources/nomad/nomad/plan_queue.go:201-217` |
| Resource release | `handleJobDeregister` with `DeleteJobTxn`/`UpsertJob Stopped` + `PeriodicLaunch` cleanup; GC `CoreScheduler.jobGC/nodeGC/evalGC`; `applyAllocClientUpdate` unblocks on terminal | `studies/aren-go-runtime-study/sources/nomad/nomad/fsm.go:868-937,1067` and `studies/aren-go-runtime-study/sources/nomad/nomad/core_sched.go:94-212` |
| Saturation test | BlockedEvals tests for `Quota` block/unblock, `IncidentalQuota`, `ImmediateUnblock` (missed), `Untrack` | `studies/aren-go-runtime-study/sources/nomad/nomad/blocked_evals_test.go:70-78,285-331,357,500-519,596-626` |
| Eval broker saturation | `EvalBroker` tests for dequeue priority etc. (exists but not fully enumerated) | `studies/aren-go-runtime-study/sources/nomad/nomad/eval_broker_test.go` (referenced) |

## Answers to Dimension Questions

### Can a request start before its policy and resources are committed?

No. Admission is synchronous on the `Job.Register` RPC path before any eval or alloc exists. `Job.doRegister` (`studies/aren-go-runtime-study/sources/nomad/nomad/job_endpoint.go:132-136`) runs `admissionControllers` (mutators then validators from `nomad/job_endpoint.go:60-89`, `nomad/job_endpoint_hooks.go:178-213`) plus ACL/ Sentinel (`nomad/job_endpoint.go:242-256`), volume/host/CSI permission checks (`nomad/job_endpoint.go:160-211`), `validateJobUpdate`/`JobMaxCount`/`JobMaxPriority` (`nomad/job_endpoint_hooks.go:518-551`), and load-shed `registrationsAreAllowed` (`nomad/job_endpoint.go:1383-1396`). Only after success is an `Evaluation{Status:pending}` created (`nomad/job_endpoint.go:310-322`) and committed via Raft. Resources are not reserved until `planApply.evaluatePlan` checks quota (`nomad/plan_apply.go:542-557` — OSS stub false, ent enforces) and `AllocsFit` per node (`nomad/plan_apply.go:784-843`), then Raft `ApplyPlanResults` atomically creates allocs (`nomad/plan_apply.go:389-400`). Allocs begin pending; node eligibility disconnected/down rejections guarantee no placement without capacity. The one gap is optimistic snapshot: schedulers run against `StateSnapshot` at `waitIndex` (`nomad/worker.go:429-460`) and may place based on slightly stale state, but the plan is never committed without re-validation against a `SnapshotMinIndex(max(prevPlanIndex, plan.SnapshotIndex))` (`nomad/plan_apply.go:182-189`) and `evaluateNodePlan` fit check.

### What prevents starvation when long and short runs share capacity?

Only weak mechanisms. Priority is the sole ordering primitive: `ReadyEvaluations` heap sorts by `Priority` desc then `CreateIndex` FIFO (`nomad/eval_broker.go:1040-1048`), `PendingEvaluations` by priority then `ModifyIndex` (`nomad/eval_broker.go:1086-1091`), and `PendingPlans` by priority then `enqueueTime` (`nomad/plan_queue.go:232-238`). `scanForSchedulers` picks highest priority across scheduler types and randomizes ties via `rand.Intn(n)` (`nomad/eval_broker.go:474-488`); node selection shuffles with eval-ID-seeded Fisher-Yates (`scheduler/feasible/feasible.go:131-152`). Per-job serialization (one ready + pending heap, `nomad/eval_broker.go:320-350`) prevents a single job monopolizing the ready queue but does not bound cross-job starvation. There is no aging, no deadline, no weighted fair queue, no back-pressure on short vs long task groups; long service batches and short system jobs compete purely on priority. Blocked evals avoid busy-loop for unsatisfiable jobs, but satisfiable low-priority jobs can starve indefinitely behind high-priority churn. Nack delays (`initialNackDelay`, `subsequentNackDelay` compounding, `nomad/eval_broker.go:728-738`) provide transient backoff, not fairness.

### Who owns retry after a rejected or stale scheduling plan?

Worker owns transport retry; scheduler owns semantic retry; broker/plan-queue own queuing. Flow:
1. Worker `invokeScheduler` error → `sendNack(eval,token)` (`nomad/worker.go:462-466`); `Plan.Submit` errors `No cluster leader`/`plan queue is disabled` → `shouldResubmit` + `backoffErr` loop (`nomad/worker.go:875-885,689-700`).
2. Inside `GenericScheduler.Process`, `retryMax(limit, s.process, progressMade)` (`scheduler/generic_sched.go:138-160`) re-invokes `process()` up to 5 (service) /2 (batch) attempts while `progressMade(planResult)`.
3. If `evaluatePlan` detects quota exceed or node over-commit, it returns `PlanResult{RefreshIndex: max(allocIndex, nodeIndex)}` (`nomad/plan_apply.go:545-556,709-718`); `ComputePlacements` returns `!fullCommit` → `missing state refresh after partial commit` → outer retryMax re-processes with new snapshot (`scheduler/generic_sched.go:307-322`). The stale `eval.Status==blocked && failedTGAllocs!=0` path reblocks (`scheduler/generic_sched.go:163-172`).
4. Persistent failure creates a blocked eval via `createBlockedEval` (`scheduler/generic_sched.go:181-200`) with `TriggeredBy: MaxPlans` or `FailedPlacements`, carrying `ClassEligibility/Escaped/QuotaLimitReached`. The FSM routes it to `BlockedEvals.Block/Reblock` (`nomad/fsm.go:969-984`). Retry is then event-driven: `Unblock*` events (node ready, capacity change, quota update, alloc terminal) re-enqueue via `EvalBroker.EnqueueAll` (`nomad/blocked_evals.go:752-758`). The broker’s `Nack` path also requeues with delay or fails to `_failed` queue after `deliveryLimit`.

### Are queue depth and waiting time bounded and visible to callers?

Partially. Waiting time has bounded delays: `Wait` (immediate) and `WaitUntil` (absolute time via `delayHeap` + `runDelayedEvalsWatcher`, `nomad/eval_broker.go:283-293,886-917`), Nack requeue delays (`initialNackDelay` then `prevDequeues*subsequentNackDelay`, `nomad/eval_broker.go:728-738`), and per-eval `NackTimeout` with pause/resume around plan queue (`nomad/plan_endpoint.go:70-73`, `nomad/eval_broker.go:740-772`). `deliveryLimit` caps retries before moving to `_failed` (`nomad/eval_broker.go:709-711`). Visibility is strong: `BrokerStats{TotalReady,TotalUnacked,TotalPending,TotalWaiting,TotalCancelable,DelayedEvals,ByScheduler}` (`nomad/eval_broker.go:1019-1027`, `984-1016`), `BlockedStats` with `TotalQuotaLimit` etc. (`nomad/blocked_evals_stats.go:23-70`, `nomad/blocked_evals.go:867-905`), `QueueStats.Depth` (`nomad/plan_queue.go:220-217`), broker timing histograms (`wait_time/process_time/response_time`, `nomad/eval_broker.go:408-413,798`). Queue depth itself is **not bounded**: `ready` heap, `pending` maps, `PlanQueue.ready`, `blockedEvals.captured/escaped` and `capacityChangeCh` buffer 8096 (`nomad/blocked_evals.go:21`) grow unbounded; no admission rejection on full queue. Callers do not receive a 429 queue-full; they see `EvalStatus pending/blocked` via `Eval.GetEval/List` and `PlanResult.RefreshIndex` forcing refresh, but not an explicit depth limit error. The 10k cap on `enqueuedTime/dequeuedTime` maps (`nomad/eval_broker.go:337`) limits metric memory but not queue memory.

## Architectural Decisions

| Decision | Why chosen | Evidence |
|---------|-----------|----------|
| Admission as in-process mutator/validator chain before Raft | Fail-fast without replication cost; enforces policy before eval creation | `nomad/job_endpoint_hooks.go:178-213`, `nomad/job_endpoint.go:132-256` |
| Priority heap + per-job ready slot + pending heap | Serializes concurrent updates to same job, preserves priority across scheduler types | `nomad/eval_broker.go:320-350,644-667,1040-1091` |
| Separate `BlockedEvals` tracker keyed by computed node class / quota / escaped / nodeID | Avoids busy polling; unblocks are O(groups) on capacity events | `nomad/blocked_evals.go:31-95,674-759` |
| Optimistic pipelined `planApply` with `inFlightPlans` and `SnapshotMinIndex` | Overlaps Raft latency, throughput 6-8 in flight, consistency via max(prevPlanIndex, planSnapshotIndex) | `nomad/plan_apply.go:100-189,279-300` |
| Stale-plan rejection via `RefreshIndex` | Schedulers retry against fresher snapshot; guards double-placement | `nomad/plan_apply.go:554-556,709-723,510-512` |
| Ack/Nack token protocol with deliveryLimit | At-least-once delivery without eval loss; failed queue after limit | `nomad/eval_broker.go:493-525,599-724` |
| Sentinel/Quota enterprise-gated | OSS runs without multi-tenancy cost; stubs return no-op | `nomad/plan_apply_ce.go:27-30`, `nomad/job_endpoint_ce.go:12` |
| Node shuffling seeded by eval ID | Deterministic yet distributed placement, precludes thundering herd | `scheduler/feasible/feasible.go:131-165` |
| PauseEvalBroker via `SchedulerConfiguration.PauseEvalBroker` | Operator pause for debugging/delete without losing evals | `nomad/structs/operator.go:235-239`, `nomad/eval_endpoint.go:161-165` |

## Notable Patterns

- **Tokenized at-least-once with pause**: `Dequeue → token + NackTimer → OutstandingReset → Ack/Nack/{Pause,Resume}NackTimeout` (`nomad/eval_broker.go:493-772`, `nomad/worker.go:650-673`). Raft-index coupling via `plan.SnapshotIndex`/`waitIndex` ensures scheduler view includes dequeued eval.
- **Blocked-as-cache**: `processBlock` deduplicates per `NamespacedID` (`nomad/blocked_evals.go:289-343`, `298-342`), preferring higher `max(CreateIndex,SnapshotIndex)`, and stores `ClassEligibility`/`Escaped`/`QuotaLimitReached`/`MissingNonNodeResources` for precise unblock filtering.
- **Missed-unblock correctness**: compares `eval.SnapshotIndex` to `unblockIndexes[class/quota/resource].index` to immediately re-enqueue evals that missed a capacity event while scheduling (`nomad/blocked_evals.go:359-417` → `EnqueueAll`).
- **Feasibility memoization**: `FeasibilityWrapper` caches `EvalComputedClass{Eligible,Ineligible,Escaped}` and skips checkers on hot path (`scheduler/feasible/feasible.go:1450-1542`).
- **Preemption as second pass**: `selectNextOption` first without preemption, then with `Preempt=true` if `PreemptionConfig` enabled (`scheduler/generic_sched.go:910-935`).
- **Destructive placement atomicity**: `reconciler` defers stopping previous alloc until new alloc placed (`scheduler/generic_sched.go:573-575,689-717`).

## Tradeoffs

- **Priority-only vs fair sharing**: Simplifies reasoning, favors latency for high priority; risks indefinite starvation for low priority batch/large jobs. No WFQ, leasing, or runtime quotas in OSS.
- **Optimistic pipeline vs consistency**: Throughput gain (batched Raft) trades snapshot staleness → higher `RejectedNodes`/`RefreshIndex` churn under heavy load. `planApply` tolerates via retryMax; operators see info logs for plan rejections (`nomad/plan_apply.go:616-618`).
- **Unbounded queues vs availability**: Ready/pending/plan queues never reject; favors correctness under burst but risks memory pressure (only `unblockBuffer=8096`, `enqueuedTime 10k` bounded). Load-shedding is coarse (`RejectJobRegistration` bool).
- **Per-job serialization vs parallelism**: Prevents duplicate schedulers stepping on same job, but serializes scaling/deployment updates for large jobs.
- **BlockedEvals class-coarse unblocking**: Unblocking by computed class wakes all escaped evals (`nomad/blocked_evals.go:693-699`) → thundering herd on large clusters vs fine-grained feasibility miss.
- **OS-level resource fit at apply time only**: Rich feasibility filters during scheduling, final `AllocsFit` gate at plan apply (`nomad/plan_apply.go:842`); over-commit impossible but scheduler may waste work on stale data.
- **Enterprise quota absent locally**: Nomad OSS cannot enforce per-namespace CPU/mem reservations; relies on node `AllocsFit` and deployment-concurrency; multi-tenant fairness must be external (node pools only).

## Failure Modes / Edge Cases

- **Stale snapshot double placement**: Scheduler sees index N, plan applied after node allocation freed → `evaluateNodePlan` would fit but later Raft may overlap. Mitigated by `RefreshIndex` forcing retry and `AllocsSubset` check for in-place updates (`nomad/plan_apply.go:816-818`). Unbounded in-flight pipeline increases window.
- **Duplicate blocked evals**: Concurrent evals for same job coalesced to one (`processBlockJobDuplicate`, `nomad/blocked_evals.go:298-343`) and remainder sent to `duplicates` → `Cancelable` raft batch; bug leaves `jobs` map entry → `logger.Error "existing blocked evaluation is neither tracked as captured or escaped"` (`nomad/blocked_evals.go:319`).
- **Lost unblock due to broker disabled / flush**: `SetEnabled(false)` flushes broker (`nomad/eval_broker.go:200-211`) and blocked (`nomad/blocked_evals.go:844-865`) clearing `ready/pending/unack/captured/escaped/jobs/unblockIndexes`; in-flight evals Nacked via timer expiry. Leadership flaps close `waiting` channels (`nomad/eval_broker.go:824-827`) causing `waitForSchedulers` to retry.
- **Nack timer race**: `Ack` checks `!NackTimer.Stop()` → `Ack'd after Nack timer expiration` error (`nomad/eval_broker.go:620-622`), then still cleans up but leaves requeue logic; `OutstandingReset` returns `ErrNackTimeoutReached` if timer already fired.
- **Escaped class pressure**: Constraint outside computed class (e.g., `distinct_property` on meta) forces `EscapedComputedClass=true` → every node-class update wakes it (`nomad/blocked_evals.go:693-699`), causing O(escaped) scan.
- **Missing Raft index wait timeout**: `worker.snapshotMinIndex` with `raftSyncLimit=5s` → Nack and 10× slow sync wait (`nomad/worker.go:430-444`, `50`); slow follower behind leader repeatedly Nacks without progress until catchup.
- **JobMaxCount bypass via multiregion**: `multiregionRegister/Drop` paths (`nomad/job_endpoint.go:291-391`) replicate count checks per region but no cross-region aggregate enforcement.
- **Plan queue flushed under backpressure**: `Enqueue` returns `plan queue is disabled` (`nomad/plan_queue.go:105`) → worker `shouldResubmit` retries with slow backoff; unbounded retry can stall evaluation.
- **Prune window**: `unblockIndexes` pruned at `5m interval, 15m threshold` (`nomad/blocked_evals.go:23-29,907-935`); eval with `SnapshotIndex` older than pruned window loses missed-unblock detection → may stay blocked longer needing next capacity event.

## Future Considerations

- **Weighted fair queuing / tenant quotas**: Introduce `BlockedStats.TotalQuotaLimit` enforcement locally; surface quota as first-classAdmission (OSS `evaluatePlanQuota` currently no-op). Consider node-pool-aware `ReadyEvaluations` sharding.
- **Bounded admission on queue depth**: Add configurable `maxReady/maxPending/maxPlanDepth` with 429 `ErrQueueFull` observed via `BrokerStats`/`QueueStats`; expose to `Job.Register` as `JobModifyIndex` retry.
- **Aging / starvation timer**: Promote long-waiting evaluations (e.g., exponentially boost priority after `wait_time > threshold`) leveraging existing `enqueuedTime` metrics (`nomad/eval_broker.go:403-413`).
- **Fine-grained unblock indexing**: Replace class-level wake-all for escaped evals with inverted index on constraint keys to avoid O(escaped) fan-out.
- **Observable queuing depth API**: Expose `TotalReady/Pending/Waiting/Depth` plus per-job `PendingEvaluations` length to callers (currently only metrics/logs) for SRE dashboards and client-side backoff.

## Questions / Gaps

- **Per-workspace / per-tenant limits in sibling source?** Workspace config provider and policy paths were not inspected per isolation rule; Nomad OSS shows only `NodePool` scheduler overrides (`nomad/structs/operator.go:278-294`) and enterprise quota pointers in blocked stats. No workspace-scoped rate limits found; searched `quota` in `nomad/*` (found only blocked/quota stats and CE stub) and `policy` in admission scope.
- **Exact deliveryLimit, nackTimeout, initialNackDelay values?** Defaults are injected from `ServerConfig` (`EvalNackTimeout`, `EvalNackInitialDelay`, `EvalDeliveryLimit` not grepped due to scope); assumption is config-driven.
- **Cancellation while queued semantics beyond Ack?** `DropWaiting` (`nomad/eval_broker.go:804-818`) handles follow-up reconnect eval cancellation (`nomad/fsm.go:1082-1094`); user-initiated eval deletion while `unack` requires broker paused (`nomad/eval_endpoint.go:484`). No direct `Cancel(evalID)` while pending/ready beyond `BlockedEvals.GetDuplicates`/`Untrack`.
- **Resource reservations vs allocations?** Nomad has no `reservation` primitive distinct from alloc; capacity is checked via `AllocsFit` including proposed allocs, not held locks. Equivalent to immediate commit.
- **System scheduler fairness**: `system` jobs use `systemEvals` map (`nomad/blocked_evals.go:64-65,280-282,744-749`) keyed by node; fairness across system evaluators not benchmarked.

---

Generated by `01.12 Policy, Admission, Scheduling, and Fairness` against `nomad`.
