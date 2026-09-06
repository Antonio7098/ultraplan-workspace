# Source Analysis: nomad

## 01.09 Resource Accounting, Overload, and Bounded Work

### Source Info

| Field | Value |
|-------|-------|
| Name | nomad |
| Path | `studies/aren-go-runtime-study/sources/nomad` |
| Language / Stack | Go (server+client, Raft, yamux, go-memdb) |
| Analyzed | 2026-08-29 |

## Summary

Nomad implements a multi-layered resource governance stack: in-memory scheduling queues (`EvalBroker`, `BlockedEvals`, `PlanQueue`) with priority-heap admission and per-job serialization; per-node alloc fitting via `AllocsFit` reconciling CPU/memory/cores/devices/network against `NodeResources - ReservedResources`; client-side GC with three triggered thresholds (disk, inode, alloc count); bounded worker pools (`EvaluatePool`, `AllocGarbageCollector` semaphore), connection pools with idle-stream caps and single-flight dial coalescing, and explicit pipeline backpressure between scheduler, leader snapshot, and Raft apply. The design is strong for disk and CPU reservations and for preventing alloc-dir exhaustion, but weak on global memory backpressure, hard caps for total queued evals/plans, and true per-tenant isolation — one job/namespace without quota can still flood the broker, and the server has no OOM-prevention circuit breaker beyond a 1 MB Raft warn and best-effort metric-map truncation.

## Rating

**6/10** — Concrete boundedness exists where failure is most visible (disk GC, alloc fitting, plan pipeline, RPC/stream limits, pool idle caps), with explicit admission→reservation→release and compensation on failure. Missing or soft: no global queue-depth or memory-pressure quotas, no pre-OOM load shedding for the broker/planner heap, metrics are rich but mostly counters/gauges without alert thresholds wired to admission, and saturation benchmarks exist only for the scheduler micro-benchmark, not end-to-end bursts or long-lived session retention.

## Evidence Collected

Every entry MUST include a file path with line numbers. Format: `path/to/file.ts:NN`.

| Area | Evidence | File:Line |
|------|----------|-----------|
| EvalBroker queue | `EvalBroker` struct: `pending map[NamespacedID]PendingEvaluations`, `ready map[string]ReadyEvaluations`, `unack map[string]*unackEval`, `timeWait`, `delayHeap *DelayHeap`, `delayedEvalsUpdateCh chan struct{1}`, metrics maps `enqueuedTime/dequeuedTime` capped 10k | `nomad/eval_broker.go:53-122` |
| Eval admission dedup | `processEnqueue` checks `b.evals[eval.ID]` dedup; `enqueueLocked` enforces 1 ready eval per job `jobEvals[nsID]` else pushes to `pending` heap | `nomad/eval_broker.go:254-274`, `nomad/eval_broker.go:320-349` |
| Delayed/wait queue | `Wait>0` → `time.AfterFunc` + `timeWait` map; `WaitUntil` → `delayHeap.Push` + `runDelayedEvalsWatcher` signals `delayedEvalsUpdateCh` buffered 1 | `nomad/eval_broker.go:278-293`, `nomad/eval_broker.go:886-918` |
| Priority dequeue + fairness | `scanForSchedulers` picks highest `Priority` then `CreateIndex`, random among ties | `nomad/eval_broker.go:435-488` |
| At-least-once + nack | `dequeueForSched` arms `time.AfterFunc(nackTimeout)` → `Nack`, `unack` map, `deliveryLimit` → `_failed` queue; `nackReenqueueDelay` compounds `(prev-1)*subsequent` | `nomad/eval_broker.go:493-525`, `nomad/eval_broker.go:678-737` |
| EnqueuedTime truncation | `len(enqueuedTime)<10_000` and `len(dequeuedTime)<10_000` before storing timings — prevents metric-map OOM | `nomad/eval_broker.go:337-338`, `nomad/eval_broker.go:404-406` |
| BlockedEvals secondary queue | `BlockedEvals.captured/escaped/jobs/unblockIndexes`, `capacityChangeCh chan *capacityUpdate, 8096` comment “ensure FSM doesn't block → back-pressure on Raft” | `nomad/blocked_evals.go:17-21`, `nomad/blocked_evals.go:36-78` |
| One-blocked-eval per job | `processBlockJobDuplicate` keeps newest by `max(CreateIndex,SnapshotIndex)`, older → `duplicates` batch | `nomad/blocked_evals.go:298-343` |
| Unblock index retention bound | `pruneInterval=5m`, `pruneThreshold=15m`, `pruneUnblockIndexes` evicts old `unblockIndexes` | `nomad/blocked_evals.go:24-28`, `nomad/blocked_evals.go:926-935` |
| PlanQueue commit queue | `PlanQueue.ready PendingPlans` heap ordered `Priority desc → enqueueTime FIFO`, `waitCh chan struct{1}`, `Depth` gauge | `nomad/plan_queue.go:33-40`, `nomad/plan_queue.go:99-127` |
| Plan pipeline depth | `planApply` creates `planIndexCh chan uint64,maxPipelineDepth` (6-8), `inFlightPlans` tracking, `snapshotMinIndex` 10s timeout as backpressure; drains via non-blocking `planIndexCh` | `nomad/plan_apply.go:100-106`, `nomad/plan_apply.go:140-164`, `nomad/plan_apply.go:212-219`, `nomad/plan_apply.go:279-300` |
| EvaluatePool bounded workers | `EvaluatePool workers=int, req/res chan evaluateRequest/Result buf=64`, `NewEvaluatePool(poolSize,64)` where `poolSize=NumCPU/2`, `SetSize` resize, `Shutdown` closes workers | `nomad/plan_apply_pool.go:11-26`, `nomad/plan_apply_pool.go:42-55`, `nomad/plan_apply.go:120-126` |
| Alloc fitting reservation | `AllocsFit(node, allocs, nil, true)` computes `used=ComparableResources`, checks `NodeMaxAllocs`, reserved cores overlap, `available=NodeResources-ReservedResources`, `Superset(used)`, `NetworkIndex`, `DeviceAccounter` | `nomad/structs/funcs.go:142-232` |
| MaxParallel admission knob | `MaxParallel` validated `<0`, enforced `scheduler_system.go:697-702`, `reconcile_cluster.go:798-810` — caps rolling updates/preemptions | `nomad/structs/structs.go:5417-5425`, `scheduler/scheduler_system.go:697-702`, `scheduler/reconciler/reconcile_cluster.go:798-810` |
| Client GC retained-history bound | `GCConfig {MaxAllocs, DiskUsageThreshold, InodeUsageThreshold, Interval, ParallelDestroys}`, defaults `2,80,70,50` | `client/config/config.go:234-253`, `client/config/config.go:951-954` |
| GC semaphore | `AllocGarbageCollector.destroyCh chan struct{} capacity=ParallelDestroys`, acquire `<-destroyCh` before `Destroy` and release after | `client/gc.go:54-55`, `client/gc.go:83`, `client/gc.go:175-196` |
| GC priority queue | `IndexedGCAllocPQ heap Less=timestamp Before`, FIFO oldest-first, `Pop`→`destroyAllocRunner` goroutine | `client/gc.go:240-349` |
| GC trigger logic | `keepUsageBelowThreshold` polls `Collect()`, checks `UsedPercent>80`, `InodesUsedPercent>70`, `liveAllocs>50` (`2x` before WARN), pops until under threshold | `client/gc.go:120-169` |
| ConnPool bounded streams | `ConnPool.maxStreams`, `Conn.clients *list.List`, `returnClient` keeps `Len()<maxStreams` else `Close()`, `Shrink()` on yamux stream | `helper/pool/pool.go:180-193`, `helper/pool/pool.go:122-139` |
| ConnPool idle eviction | `reap()` every 1s, `now.Sub(*lastUsed)<maxTime` skip, `refCount>0` skip else `Close()+delete(pool)` | `helper/pool/pool.go:549-583` |
| Single-flight dial throttle | `limiter map[string]chan struct{}`, first thread dials `getNewConn`, others block on `wait` chan, `delete(limiter)`+`close(wait)` on done | `helper/pool/pool.go:317-385` |
| RPC connection limits | `rpcHandler.connLimiter/streamLimiter *connlimit.Limiter`, `MaxConnsPerClientIP`, `Accept` wraps with `connlimit.Wrap`, streaming has lower limit `connLimit - LimitsNonStreamingConnsPerClient` | `nomad/rpc.go:51-63`, `nomad/rpc.go:96-105`, `nomad/rpc.go:221-234`, `nomad/rpc.go:393-401` |
| HTTP conn limiter | `connLimiter(connLimit)` → `rate.NewLimiter(10,100)` + `connlimit.NewLimiter` | `command/agent/http.go:322-330` |
| Rate limiters (golang.org/x/time/rate) | `volume_watcher limiter 1/sec burst 3`, `deploymentwatcher 100 qps burst 100`, `encrypter keyring replication`, `drainer queryLimiter`, `core_sched 100 burst 100`, `leader replicationRateLimit` | `nomad/volumewatcher/volume_watcher.go:46-69`, `nomad/deploymentwatcher/deployments_watcher.go:114`, `nomad/encrypter.go:1135`, `nomad/drainer/drainer.go:182`, `nomad/core_sched.go:1324` |
| Worker dequeue backoff | `dequeueTimeout=500ms`, `raftSyncLimit=5s`, `backoffLimitSlow=10s`, `helper.Backoff(base,limit,failures)`, `snapshotMinIndex` timeout → `sendNack` + 50s slow catchup | `nomad/worker.go:27-50`, `nomad/worker.go:892-915`, `nomad/worker.go:428-444` |
| Raft enqueue/size bound | `enqueueLimit=30s`, `raftWarnSize=1MB`, `raftApplyFuture` warns on large entry | `nomad/rpc.go:39-45`, `nomad/rpc.go:789-802` |
| User lookup cache TTL | `cacheTTL=1h`, `failureTTL=1m`, `TTL` eviction not LRU, `sync.Mutex` serialized | `helper/users/cache.go:15-17`, `helper/users/cache.go:51-87` |
| Cgroup reservation partitioning | `Init` creates `share`/`reserve` cgroup partitions, sets `cpuset.cpus/mems`, `subtree_control` v1/v2 | `client/lib/cgroupslib/init.go:26-193` |
| Scheduler benchmarks | `BenchmarkServiceScheduler` with harness, `BenchmarkServiceStack_With_ComputedClass` 5000 nodes, `BenchmarkHTTPRequests`, `BigBenchmarkJob` | `scheduler/benchmarks/benchmarks_test.go:74-116`, `scheduler/feasible/stack_test.go:17-27`, `command/agent/http_test.go:46`, `nomad/mock/job.go:627` |

## Answers to Dimension Questions

**Which resources are bounded per run, per workspace, and process-wide?**

- **Per-job (per-run) bounded:** One ready eval per `NamespacedID` (`jobID+namespace`) serialized via `jobEvals` + `pending` heap (`nomad/eval_broker.go:327-349`); at most one `BlockedEvals` entry per job with duplicate cancellation (`nomad/blocked_evals.go:298-343`). Rolling parallelism bounded by `UpdateStrategy.MaxParallel` / `TG.Update.MaxParallel` (`nomad/structs/structs.go:5417`, `scheduler/scheduler_system.go:697`, `scheduler/reconciler/reconcile_cluster.go:798`) — controls how many allocs/placements proceed concurrently during a deployment, not total job size.
- **Per-node (workspace-equivalent) bounded:** Node alloc count via `Node.NodeMaxAllocs` checked in `AllocsFit` (`nomad/structs/funcs.go:145-149`); per-task CPU/memory `Resources.CPU/Memory/Cores` accounted in `ComparableResources` and validated against `NodeResources - ReservedResources` (`nomad/structs/funcs.go:142-198`). Disk/inode thresholds trigger GC on that client (`client/gc.go:141-153`). Cgroup partitions `share` vs `reserve` enforce `ReservedCores` (`client/lib/cgroupslib/init.go:26-193`).
- **Process-wide bounded:** `PlanQueue.Depth` heap (priority) but unbounded length except memory; `EvalBroker` ready/pending/unack/waiting all unbounded except metric-map truncation at 10k (`nomad/eval_broker.go:337`); `BlockedEvals.capacityChangeCh` buffered `8096` (`nomad/blocked_evals.go:18-21`); ConnPool: one `yamux.Session` per address, `maxStreams` idle cached streams per conn, `maxTime` TTL (`helper/pool/pool.go:187-193`, `549-583`); RPC `connLimiter/streamLimiter` per client IP (`nomad/rpc.go:96-105`); Raft pipeline `maxPipelineDepth` 6-8 (`nomad/plan_apply.go:105`); `EvaluatePool` workers `NumCPU/2` with 64-buffered chans (`nomad/plan_apply_pool.go:15`, `nomad/plan_apply.go:121`); GC destroys `ParallelDestroys=2` (`client/config/config.go:951`). Quotas provide multi-tenant limits (`nomad/structs/eval.go:200`, `nomad/blocked_evals.go:371-383`) but are opt-in.
- **Not bounded:** Total evaluations/plans in flight have no hard max count — growth is limited only by Raft throughput and memory. At-least-once redelivery can re-enqueue indefinitely until `deliveryLimit`.

**What happens before memory pressure becomes an out-of-memory crash?**

- Limited defenses. Client side: disk/inode alloc GC mitigates disk exhaustion, not heap. Server side: `enqueuedTime/dequeuedTime` maps capped at 10k (`nomad/eval_broker.go:337`), metrics sampling avoids unbounded tag cardinality via `ParentIDFromJobID` deduplication, Raft entry size warns at 1 MB (`nomad/rpc.go:39`) but does not reject, `snapshotMinIndex` 10s timeout sheds via blocking schedulers (`nomad/plan_apply.go:286-297`). `BlockedEvals` pruning (`pruneInterval 5m`, `pruneThreshold 15m`) and `helper/users/cache.go:15` TTL (1h) bound auxiliary maps. HTTP and pool reapers evict idle conns. However, there is no global heap-pressure listener, no `cgroup` memory limit enforcement on the Nomad process itself, no queue-depth shedding (plan/eval enqueue never rejects due to depth), and no explicit `runtime.MemStats` → load-shedding. Under sustained eval spam the broker’s `ready+pending+unack+evals` maps grow until GC or OOM.

**Are estimates reconciled with actual use and always released on failure?**

- **Reservation vs actual:** `AllocsFit` is the reconciliation point: it sums `AllocatedResources.Comparable()` over non-terminal allocs plus proposed placements, subtracts `ReservedResources`, and fails the plan if not a superset (`nomad/structs/funcs.go:142-198`). Preemptions remove allocations before re-adding (`nomad/plan_apply.go:825-839`). Terminal allocs are skipped (`ClientTerminalStatus`). Partial commits set `RefreshIndex` forcing schedulers to re-snapshot (`nomad/plan_apply.go:709-718`), and optimistic snapshot `UpsertPlanResults` is only kept while `inFlightPlans==0` else invalidated (`nomad/plan_apply.go:157-163`).
- **Release on failure:** `PlanQueue.Flush` and `EvalBroker.flush` on leadership loss error pending futures (`nomad/plan_queue.go:168-186`, `nomad/eval_broker.go:822-865`), `Nack`/`Ack` handlers remove `unack`, `evals`, `jobEvals`, promote one `pending` entry and batch-cancel the rest via `MarkForCancel` (`nomad/eval_broker.go:599-674`, `1116-1129`). `BlockedEvals.untrackImpl` and `UnblockFailed` release jobs on successful eval (`nomad/blocked_evals.go:445-482`, `772-805`). GC `destroyCh` semaphore is always released `<-destroyCh` after `ar.Destroy()` (`client/gc.go:192-195`) and shutdown paths guard `select <-shutdownCh` before acquire. Conn `refCount`/`shouldClose` ensures pooled session closed only at zero refs (`helper/pool/pool.go:73-79`, `450-462`).

**Can one run monopolise workers, model memory, event buffers, or disk?**

- **Workers:** Scheduler workers dequeue fairly: `scanForSchedulers` randomly picks among schedulers at equal highest priority (`nomad/eval_broker.go:484-487`), but priority starvation is possible — a tight loop of Priority=100 evals starves lower priorities; per-job serialization prevents one job from holding multiple ready slots, but it can continuously re-enqueue via `Reblock` and stay head-of-queue. `EvaluatePool` is shared process-wide, not per-job, so a large plan fan-out (many `nodeID`s) can occupy all `NumCPU/2` workers until `AllocsFit` completes, delaying other plans.
- **Model memory:** No model residency concept; closest is alloc memory reservations enforced by `AllocsFit` and cgroups. Without quota, a single job with high `count` can reserve all allocatable memory across nodes.
- **Event buffers:** `BlockedEvals.capacityChangeCh` 8096 and broker `waiting chan 1` are the largest buffers; flooding node updates (via FSM) enqueues unblocks that are drained by single `watchCapacity` goroutine (`nomad/blocked_evals.go:648-672`) — 8096 depth is large but not back-pressured to callers except via `select <-done` non-blocking send semantics (`nomad/blocked_evals.go:192-195`) which drops on shutdown only, not overload. `PlanQueue.waitCh` coalesces signals (size 1) so bursts coalesce but not drop work.
- **Disk:** Single alloc cannot monopolize disk beyond its `ephemeral_disk` reservation? Actually disk not accounted in `AllocsFit` except via GC — `GCDiskUsageThreshold 80%` and `GCInode 70%` plus `GCMaxAllocs 50` (`client/config/config.go:951`) trigger FIFO GC oldest first (`client/gc.go:161-167`), limiting alloc-dir retention. But log rotation (`client/logmon/logging/rotator_test.go:278`) and artifact/host-volume usage are not bounded per-job.

## Architectural Decisions

| Decision | Why made / Consequence | Evidence |
|----------|------------------------|----------|
| Per-job ready serialization + pending heap | Avoids thundering herd for hot jobs, preserves FIFO by `ModifyIndex` within priority, makes cancellation cheap via `MarkForCancel` | `nomad/eval_broker.go:320-349`, `nomad/eval_broker.go:1116-1129` |
| Separate Delayed queue (`Wait`/`WaitUntil`) + watcher goroutine | Supports cron/periodic and dedup without busy-polling; `DelayHeap` with `delayedEvalsUpdateCh:1` coalesces wakeups | `nomad/eval_broker.go:278-305`, `nomad/eval_broker.go:886-918` |
| `BlockedEvals` as secondary buffer decoupled from Raft FSM via 8096 channel | Prevents Raft apply-thread blocking; capacity changes processed async, with duplicate collapsing | `nomad/blocked_evals.go:17-21`, `nomad/blocked_evals.go:298-343`, `nomad/blocked_evals.go:648-672` |
| Optimistic pipelined plan applier (6-8 depth) with snapshot invalidation | Overlaps Raft commit latency with next plan evaluation; trades staleness (higher rejection rate under load) for throughput | `nomad/plan_apply.go:70-127`, `nomad/plan_apply.go:140-164`, `nomad/plan_apply.go:212-241` |
| `AllocsFit` centralizes all resource checks | Single superset check covers CPU/memory/cores/network/devices/volumes/max-allocs; schedulers rely on it for correctness | `nomad/structs/funcs.go:142-232` |
| `MaxParallel` as only rolling-update throttle | Gives operators explicit concurrency knob; nil deployment case allows `MaxParallel` placements (`reconcile_cluster.go:798`) | `nomad/structs/structs.go:5417`, `scheduler/scheduler_system.go:697-702` |
| GC semaphore `ParallelDestroys=2` + FIFO GC heap | Protects client from IO storm during bulk GC while favoring oldest terminated allocs for prompt disk reclamation | `client/gc.go:54-55`, `client/gc.go:83`, `client/gc.go:240-315` |
| ConnPool per-address session + bounded idle streams + reaper | Reuses TCP/yamux sessions, caps idle memory via `Shrink()`, TTL eviction, single-flight dial avoids SYN flood | `helper/pool/pool.go:187-196`, `helper/pool/pool.go:122-139`, `helper/pool/pool.go:549-583`, `helper/pool/pool.go:317-385` |
| Per-IP RPC `connlimit` + lower streaming limit | Prevents one client (e.g., misbehaving node) from starving Raft/RPC; streaming capped lower to reserve non-streaming conns | `nomad/rpc.go:51-63`, `nomad/rpc.go:96-105`, `nomad/rpc.go:221-234` |

## Notable Patterns

- **Bounded semaphore via buffered channel:** `destroyCh chan struct{} capacity=ParallelDestroys` (`client/gc.go:55`) and `connLimiter` via `go-connlimit` provide explicit parallelism caps without mutex.
- **Coalescing channel (size 1):** `waiting[sched] chan struct{1}`, `PlanQueue.waitCh`, `delayedEvalsUpdateCh`, `triggerCh` all use `select { case ch<-struct{}{}: default: }` to collapse bursts into one wake-up (`nomad/eval_broker.go:373-377`, `nomad/plan_queue.go:122-125`, `nomad/eval_broker.go:288-291`, `client/gc.go:110-114`).
- **Single-flight dial:** `limiter map[string]chan struct{}` pattern ensures only leader thread dials per address; followers wait on close (`helper/pool/pool.go:317-385`).
- **Two-level priority heaps:** `ReadyEvaluations` (priority → CreateIndex) and `PendingEvaluations` (priority → ModifyIndex) implement staged admission (`nomad/eval_broker.go:1043-1091`).
- **TTL-based bounded caches without LRU:** `helper/users/cache.go:15-17` bounds by time, not size — simple but can grow unbounded if unique usernames unbounded (unlikely but possible).
- **Rate-limited blocking queries:** `drainer`, `deploymentwatcher`, `volumewatcher` all inject `rate.Limiter` into watches to prevent query storms (`nomad/drainer/drainer.go:182`, `nomad/deploymentwatcher/deployments_watcher.go:114`, `nomad/volumewatcher/volume_watcher.go:69`).
- **Exponential backoff with jitter:** Worker uses `helper.Backoff` with fast baseline 20ms → 10s slow for general errors, 30s flat for version mismatch (`nomad/worker.go:46-50`, `nomad/worker.go:892-904`).

## Tradeoffs

- **Throughput vs staleness:** Pipelining 6-8 plans improves Raft batching but optimistic snapshot may be 6-8 plans behind, increasing `AllocsFit` rejections and `RefreshIndex` retries (`nomad/plan_apply.go:223-227` comment explicitly calls out trade-off).
- **Deduplication vs fairness:** Per-job single-ready slot prevents one job from flooding ready queue but can starve other jobs if high-priority job continuously re-blocks; low-priority jobs never get dequeued while high-priority jobs churn (`nomad/eval_broker.go:341-349`).
- **Large `unblockBuffer` (8096) vs memory:** Decoupling FSM from scheduler avoids Raft backpressure at cost of holding up to 8096 `capacityUpdate` objects (each with eval refs) in heap during burst.
- **Metric-map truncation (10k) vs observability:** Capping `enqueuedTime/dequeuedTime` prevents OOM but drops `wait_time/process_time/response_time` samples for evals beyond 10k, biasing tail latency metrics (`nomad/eval_broker.go:337-338`).
- **Fixed GC thresholds (80%/70%/50 allocs) vs workload heterogeneity:** One-size-fits-all; nodes with large alloc-dir disks hit threshold late, small disks trigger GC too aggressively; no per-job disk quota.
- **Global `EvaluatePool` vs per-plan isolation:** Shared `NumCPU/2` pool maximizes utilization but allows one large plan (many nodes) to occupy all workers, head-of-line blocking smaller plans.
- **Per-IP conn limits vs NAT:** `connlimit` per client IP fails open under NAT where many nodes share IP; could reject legitimate burst.

## Failure Modes / Edge Cases

- **Broker heap OOM under eval flood:** No max-depth reject policy; attacker or bug submitting 1M evals can grow `ready+pending` heaps until process OOM — metric-map truncation doesn’t bound heap length. `EvalsTotalWaiting` metric may climb without alert threshold. Mitigation: none except Raft `enqueueLimit 30s` (`nomad/rpc.go:45`) and operator quota.
- **BlockedEvals leak if `unblock` never fires:** Evaluations blocked on infeasible constraints (e.g., unsatisfiable `distinct_hosts`) remain in `captured/escaped` forever; pruner only evicts `unblockIndexes`, not blocked evals themselves (`nomad/blocked_evals.go:924-935` vs `nomad/blocked_evals.go:908-921`). Requires manual `EvalDelete` or job deregister to clear via `Untrack`.
- **GC livelock when `destroyAllocRunner` hangs:** `destroyCh` slot held until `ar.DestroyCh()` signals (`client/gc.go:187-195`); if task driver hangs on destroy, semaphore slot leaks (capacity effectively 1) and GC stalls. `shutdownCh` path releases via `select` but not on driver hang.
- **Pool conn leak on `refCount` race:** `reap` skips `refCount>0` conns (`helper/pool/pool.go:569`), but if `markForUse`/`releaseUse` races with `ShouldClose` flag, a conn with `refCount==0 && shouldClose==1` is only closed in `releaseUse` — a path missing in `reap` if closed concurrently may double-close or leave idle conn past `maxTime`.
- **Plan pipeline snapshot staleness cascade:** Under heavy rejection, each partial commit forces `RefreshIndex !=0` → `snapshotMinIndex` 10s block (`nomad/plan_apply.go:709-718`, `279-300`); if Raft is slow, this amplifies contention, causing more `snapshotMinIndex` timeouts that `sendNack` and requeue, increasing broker load.
- **Worker `raftSyncLimit` vs `slowServerSyncLimit` gap:** Fast limit 5s Nacks; slow 50s warn-only path (`nomad/worker.go:441-444`) does not re-Nack — eval remains outstanding with timer running, risking `Nack` race after scheduler already started work.
- **TTL cache unbounded growth under unique keys:** `helper/users/cache.go:51-88` has no size bound; enumerating unique usernames (e.g., via CSI plugins) could grow maps unbounded within 1h window.
- **Rate limiter burst exhaustion:** `rate.NewLimiter(1,3)` (`nomad/volumewatcher/volume_watcher.go:69`) allows 3 bursts then throttles to 1/sec; under rapid volume churn, GC claims queue builds without bound (no drop policy).

## Future Considerations

- Implement hard caps + explicit shed metrics for `EvalBroker` and `PlanQueue` depth (e.g., `maxBrokerReady=50k` with `eval_broker.shed_total` counter) and test with bursty workload harness similar to `scheduler/benchmarks/benchmarks_test.go:74` but with stated hardware, repetitions, and retained heap measurement — currently missing steady-state/burst/long-lived session saturation tests.
- Add memory-pressure hook (read `runtime.MemStats` / cgroup `memory.pressure`) that triggers sampled load shedding before OOM, wiring to `BrokerStats.TotalWaiting` and `QueueStats.Depth` gauge thresholds rather than only disk/inode.
- Per-tenant resource accounting: enforce nominal quotas by default; extend `quota` admission (`nomad/blocked_evals.go:371`) to cover CPU/memory reservations, not just allocation count, and expose `blocked_evals.total_quota_limit` as alertable SLO.
- Scope `EvaluatePool` per-scheduler or use weighted fair queuing so one large plan cannot monopolize `NumCPU/2` workers; alternatively bound per-plan node fan-out.
- Replace infinite retry for missing unblocks with TTL on blocked evals (e.g., 24h as done for `batch_eval_gc_threshold` 24h) to auto-expire infeasible evals and prevent indefinite retention.
- Add benchmarks that retain memory across iterations (e.g., 10k jobs × 100 allocs sustained for 30m) to measure retained-history growth of `enqueuedTime`, `unblockIndexes`, and ConnPool idle streams — currently `BenchmarkServiceScheduler` does not assert heap or goroutine counts.

## Questions / Gaps

- **No evidence of process-wide memory or goroutine budget:** Searched `scheduler/`, `nomad/`, `client/`, `helper/` for `MemStats`, `GOMAXPROCS`, `runtime`, `semaphore`, `errgroup` — only per-subsystem caps found; no global guardian that sheds before OOM. Confirmed absence via grep `Memory.*Threshold|mem.*pressure|runtime.Mem` returning only `client/hoststats` collectors, not admission.
- **Long-lived session retention not tested:** Grep `Benchmark.*` yielded only scheduler micro and HTTP benchmarks (`scheduler/benchmarks/benchmarks_test.go:30`, `command/agent/http_test.go:46`); no end-to-end soak test with 24h allocations or blocked-eval retention curve.
- **Per-run cost accounting absent:** No per-evaluation/plan cost metric (CPU time or allocs evaluated) — `metrics.MeasureSince` emits latency but not cost attribution to job/namespace for chargeback (`nomad/plan_apply.go:526-528`, `nomad/worker.go:609-611`). `AllocGarbageCollector` reason string is human-readable but not a counter label for per-workspace cost.
- **Fairness across namespaces under evaluation flood not verified:** No test asserts that 1 noisy namespace cannot starve another at equal priority; `TestBlockedEvals`, `TestEvalBroker` cover duplicate collapse but not cross-tenant fairness under burst.
- **Disk accounting excludes host volumes and logs:** `GCDiskUsageThreshold` monitors `AllocDirStats` only (`client/gc.go:134-143`); host volumes (`exclusiveHostVolumeClaims` in `AllocsFit:179-185`) and log rotator outputs are separate resources with independent (or missing) quotas.

---

Generated by `01.09 Resource Accounting, Overload, and Bounded Work` against `nomad`.
