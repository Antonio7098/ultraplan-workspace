# Source Analysis: nats-server

## 01.09 Resource Accounting, Overload, and Bounded Work

### Source Info

| Field | Value |
|-------|-------|
| Name | nats-server |
| Path | `studies/aren-go-runtime-study/sources/nats-server` |
| Language / Stack | Go 1.25 / NATS Server 2.15.0-dev |
| Analyzed | 2026-08-29 |

## Summary

nats-server is a process-wide broker, not a per-run runner. It bounds every hot resource explicitly: per-client outbound bytes, per-server connection/subscription/payload limits, JetStream reserved memory/store at stream-create time plus atomic actual-usage counters, disk I/O concurrency semaphores, and intra-process queues with length/size caps. Overload is handled by hard disconnect (slow-consumer, max-connections, max-subs, max-payload) or brief fast-producer stalling before disconnect. Accounting is reconciled (reserved vs. actual store scans) and always released on close/delete, but there is no per-stream goroutine pool, no memory-pressure shedding before OOM, and no burst/queue-deadline or load-shedding for JetStream API beyond fixed queue-length limits. This maps well to Aren Phase 13 needs for quota enforcement, but leaves gaps on soft degradation and fairness under starvation.

## Rating

**6 / 10** — Strong explicit per-client and per-stream quotas plus reservation accounting, but weak on proactive memory-pressure handling, fairness isolation, and test coverage of sustained/burst overload; most load is tested via disconnected benchmarks rather than saturation tests, and several queues remain unbounded by default.

## Evidence Collected

Every entry MUST include a file path with line numbers. Format: `path/to/file.ts:NN`.

| Area | Evidence | File:Line |
|------|----------|-----------|
| Per-client outbound bound | `MAX_PENDING_SIZE = 64 MiB` default per-client pending; `Options.MaxPending:458` configurable via `max_pending:1358-1359`; applied `c.out.mp = opts.MaxPending:753` `server/client.go:753` | `server/const.go:101-102`, `server/opts.go:457`, `server/client.go:753` |
| Slow-consumer pending check | `queueOutbound` adds to `c.out.pb` then if `c.out.pb > c.out.mp` rolls back `c.out.pb -= len(data)` and counts `atomic.AddInt64(&srv.slowConsumers,1)` / `scStats.clients` / `acc.stats.slowConsumers++` then `markConnAsClosed(SlowConsumerPendingBytes):2614` | `server/client.go:2598-2614` |
| Slow-consumer write deadline | `c.out.wdl` snapshot from `Options.WriteDeadline` or per-kind `Cluster/Gateway/LeafNode.WriteDeadline:726-733`; `flushOutbound` sets `nc.SetWriteDeadline(time.Now().Add(wdl)):1842` and on timeout `markConnAsClosed(SlowConsumerWriteDeadline):1965,1999` | `server/client.go:726-733,1842,1965` |
| Fast-producer stall | Stall channel created at >75 % of max-pending `if c.out.pb > mp/4*3 && stc==nil: make(chan struct{}):3622-3624`; `stalledWait` uses `stallClientMinDuration 2ms / Max 5ms / Total 10ms:125-127` and closes stall chan when caught up `c.out.stc != nil && (n==attempted \|\| pb < mp/4*3) close:1934` | `server/client.go:125-127,2622-2624,3711-3749,1928-1934` |
| WriteDeadline budget & metrics | `DEFAULT_FLUSH_DEADLINE 10s:132`, `maxBufSize 64k:113`, `maxFlushPending 10:115`; `stalls int64:424` `atomic.AddInt64(&srv.stalls,1):3731-3732`; `slowConsumers int64:421` per-kind `scStats clients/routes/leafs/gateways:427-433` | `server/const.go:132`, `server/client.go:113-115`, `server/server.go:412-440,3872-3902` |
| Max connections / Max subs / Max payload | `DEFAULT_MAX_CONNECTIONS 64k:105`, `Options.MaxConn:432`, enforce `if len(clients) >= MaxConn -> MaxConnectionsExceeded:3381,2542`; `MaxSubs:433` enforce `closeConnection(MaxSubscriptionsExceeded):992`; `MAX_PAYLOAD_SIZE 1MB:94` enforce `maxPayloadViolation:2552-2555` | `server/const.go:94,105`, `server/server.go:3381`, `server/client.go:992,2552` |
| Disk I/O semaphore | `diskIOSemaphore{ch chan struct{}}:27-33` capacity `[minConcurrentIOs 4 .. maxConcurrentIOs 8192]:22-23`, `defaultConcurrentIOs 4096:21` created `newDiskIOSemaphore(opts.JetStreamConcurrentIOs):785`; `acquire()<-ch / release() ch<-:52-81` with wait metrics `waiters/waits/waitNanos/maxWaitNanos:29-32` | `server/dios.go:21-85`, `server/server.go:785` |
| Sync-request semaphore | `syncOutSem chan struct{}` cap `maxConcurrentSyncRequests` filled in `NewServer:799`; `Server.syncOutSem:369` used for JetStream catchup | `server/server.go:369,798-800` |
| Global catchup bytes | `gcbOut / gcbOutMax:363-364` with `gcbKick chan struct{}:366` to kick stalled catchup | `server/server.go:362-366` |
| Intra-process queue bounds | Generic `ipQueue[T]` with `mlen`/`msz` limits, `calc func(e T) uint64:39-44`, `push` returns `errIPQLenLimitReached:117` / `errIPQSizeLimitReached:123`, `pushMany` atomic revert `:153-177`, `size/mutex/inprogress atomic:323-325` | `server/ipqueue.go:26-86,114-189` |
| ipQueue callers (bounded API) | `jsAPIRoutedReqs`, `jsAPIRoutedInfoReqs`, `delayedAPIResponses:372-376`, `sendq:28-30`, `recvq/recvqp:133-134` — most created via `newIPQueue(s,name):87` without `ipqLimitByLen/Size` → unbounded by default | `server/server.go:372-376`, `server/events.go:133-134`, `server/sendq.go:28-36` |
| SendQ loop | `sendq.internalLoop` dequeues `q.pop()` in dedicated goroutine `internalLoop:41-100` reusing internal system client | `server/sendq.go:41-101` |
| Closed-clients ring buffer | `closed *closedRingBuffer:214`, `newClosedRingBuffer(opts.MaxClosedClients):896`, `DEFAULT_MAX_CLOSED_CLIENTS 10000:192` prevents unbounded history | `server/server.go:214,896`, `server/const.go:192` |
| Outbound buffer pool | `nbPoolSmall 512 / Medium 4k / Large 64k:365-388` via `sync.Pool`, `nbPoolGet/Put:393-423` recycles fixed-cap buffers | `server/client.go:365-423` |
| JetStream reserved accounting | `jetStream.memReserved/storeReserved:115-118`, `reserveStreamResources += cfg.MaxBytes:2718-2723`, `releaseStreamResources -= cfg.MaxBytes:2738-2744`; per-`jsAccount.usage[tier].local.mem/store` and `js.memUsed/storeUsed atomic:2291-2299` | `server/jetstream.go:115-118,2710-2751` |
| JetStream admission vs actual | `checkAllLimits(..., limits)` tests `memReserved+totalMaxMemory > MaxMemory:2678` and `storeReserved > MaxStore:2683`; `updateUsage(delta)` atomic add `memUsed/storeUsed:2291-2299`; `checkAndSyncUsage` scans stores `FastState(&state):2240` and repairs drift `delta = total - usage.local.mem:2250-2263` | `server/jetstream.go:2675-2708,2228-2272,2274-2306` |
| Dynamic sizing & thresholds | `dynJetStreamConfig` 75 % of `sysmem.Memory():2795` or `GOMEMLIMIT:2797-2801`, default `JetStreamMaxMemDefault 256MB:2759` and `MaxStore 1TB:2757`; disk `diskAvailable` scaled to 75 % free:2787,568 | `server/jetstream.go:2764-2808,2757-2759`, `server/sysmem/mem_linux.go:20-27` |
| TLS rate limiting | `rateCounter{limit,count,blocked}:21-28`, `newRateCounter(opts.TLSRateLimit):803`, `allow() 1/s window:37-55`, `logRejectedTLSConns countBlocked():998-1004` | `server/rate_counter.go:21-56`, `server/server.go:347,802-803,990-1004` |
| Consumer flow-control cap | `JsFlowControlMaxPending 32 MiB:602`, `setMaxPendingBytes:5631`, queue `nextMsgReqs *ipQueue:476` / `ackMsgs *ipQueue:550` | `server/consumer.go:601-602,476,550` |
| Monitoring (actionable) | `Varz.MaxConn,MaxSubs:1235-1236`, `SlowConsumers int64:226,380`, `StalledClients:384`, `ipQueuesz:1202-1212`, `JetStreamStats ReservedMemory/ReservedStore:64-65`, `JetStreamTier Limits:89`, `GOMEMLIMIT in Varz:1262` | `server/monitor.go:1235-1236,1262`, `server/events.go:226,380-384`, `server/jetstream.go:62-89` |
| Benchmarks (not saturation) | `BenchmarkCoreFanOut/FanIn/RequestReply:32,101,260,397`, `Benchmark_FileStore*:8623-9137`, `Benchmark_WS_*:5082-5382` — throughput benchmarks with ignored slow-consumer errors, no stated hardware, retention, or steady-state/burst/long-session coverage | `server/core_benchmarks_test.go:32-427`, `server/filestore_test.go:8623-9137` |
| Config validation of bounds | `validateOptions` rejects `MaxPayload > MaxPending:1163-1165`, validates cluster pool size / pinned accounts:1105-1145 | `server/server.go:1158-1166` |

## Answers to Dimension Questions

**Which resources are bounded per run, per workspace, and process-wide?**
- nats-server has no run/workspace abstraction; bounds are **per-client**, **per-account**, and **process-wide (server)**. Per-client: `MaxPending (64 MiB default):101-102,753`, `MaxPayload 1 MiB:94`, write deadline `10 s:132` plus kind-specific overrides; per-account: `MaxConn`, `max_connections` in JWT/limits (`server/accounts.go:3502, server/jwt_test.go:1388`), `MaxSubs:433`, consumer `MaxAckPending` and `MaxRequestBatch` capped by `JSLimitOpts:375-384,858-866`; process-wide: `MaxConn 64k:105` / `MaxClosedClients 10k:192`, JetStream `MaxMemory` (~75 % sysmem or GOMEMLIMIT):2795-2804 and `MaxStore` (75 % free disk or 1 TB):2787, `diskIOSemaphore` default 4096 slots capped 4-8192:21-42, `syncOutSem` maxConcurrentSyncRequests:369, `gcbOutMax` global catchup bytes:364, route `PoolSize:159,972-978`. `ipQueue` supports per-queue `MaxLen / MaxSize:75-82` but most internal queues (API, sendq, recvq) are instantiated without limits and thus unbounded by default.

**What happens before memory pressure becomes an out-of-memory crash?**
- Very little soft shedding. Client-side backpressure is the only guard: when `out.pb > 75 % MaxPending` a stall channel is created `:2622`, producers block `2-5 ms` (total 10 ms per readLoop invocation) `:125-127,3711-3749`; beyond `MaxPending` the client is hard-disconnected as `SlowConsumerPendingBytes:2614`, and write-deadline exceeded disconnects as `SlowConsumerWriteDeadline:1965`. JetStream admission checks `reserved + MaxBytes > config limit` before stream creation `:2678-2704`, but there is no heap/GC-pressure listener, no `MaxBytes` auto-shrink, no queue shedding, and no broadcast slow-down before hitting the hard limits; the server relies on `GOMEMLIMIT` detection at startup `:2797` rather than continuous reclamation. Disk pressure is similarly only bounded by `diskAvailable` at config time, not by runtime free-space polling (the 75 % scaling happens once at `dynJetStreamConfig:2787`).

**Are estimates reconciled with actual use and always released on failure?**
- **Reserved vs actual is reconciled.** `reserveStreamResources` adds `cfg.MaxBytes` to `memReserved/storeReserved:2720-2722` on stream create; `releaseStreamResources` subtracts on delete `:2741-2743`. Actual usage is incremented via `updateUsage(delta)` from storage layers `:2291-2299` maintaining `memUsed/storeUsed` atomically and per-account `usage[tier].local:2292`. Drift is corrected by `checkAndSyncUsage` which iterates every matching `StreamStore.FastState:2240`, compares `total` to `usage.local` and applies `delta` to `usage` and `js.memUsed/storeUsed:2248-2264` with warning logs. For outbound buffers the optimistic-then-revert pattern `pb += len(data)` → if `pb > mp` then `pb -= len(data):2603` guarantees no leak on slow-consumer rejection (verified by `TestPBNotIncreasedOnMaxPending:2503` and `TestClientMaxPending`). Channel/resource releases are also guaranteed: `pb -= n:1921` after flush, `nbPoolPut` on close `:1735-1737`, `dios.release() <-ch:80`, `ipQueue.recycle` returns slices under cap `:273-278`. Failure paths for JetStream still call `releaseStreamResources` via stream delete, but if a replicated stream fails mid-create before reservation is committed, no reservation was taken, so there is no double-release risk; idempotent `queueOutbound` early-return on `isClosed:2562` prevents double-count.

**Can one run monopolise workers, model memory, event buffers, or disk?**
- **Yes — partially.** A single fast publisher fanning into one slow client will fill that client's `c.out.pb` up to 64 MiB then be disconnected; aggregate fan-in is briefly throttled via the 10 ms stall window, but there is no per-account global pending cap, so many slow clients can together consume `O(N * 64 MiB)`. JetStream is protected: a single account/stream cannot exceed `MaxBytes` or tier `MaxMemory/MaxStore` because admission sums `memReserved/storeReserved` plus requested `totalMaxMemory:2678-2705`; however, within those caps one large `MemoryStorage` stream (or many small ones up to the account limit) can pin most of `MaxMemory` (75 % RAM) because there is no per-stream fair-share beyond `MaxBytesRequired:79,2412`. Disk I/O is fairly shared via the bounded semaphore (4096 slots) but disk *space* can be monopolised by a single account up to its tier limit. `ipQueues` used for API/routed requests are unbounded, so a burst of `$JS.API.>` requests from one account can grow those queues without shedding until `wouldExceedLimits:2388-2402` rejects later streams — event buffers are not per-sender capped. Goroutine amplification is limited: `grWG` tracks loops, `grTmpClients` holds handshaking conns `:251,899`, but no worker pool token is per-run isolated.

## Architectural Decisions

- **Per-client byte cap + hard kill instead of load shedding.** Choosing `MAX_PENDING_SIZE 64 MiB:102` and `queueOutbound` pending-bytes check `:2600-2614` gives deterministic OOM isolation per client at the cost of abrupt disconnects; the alternative (global backpressure / shedding with deadlines) would complicate client contracts. The 75 % stall gate `:2622` is a deliberate micro-backpressure (2-5 ms) to absorb fan-in bursts without yet killing the consumer.

- **Reservation accounting dual-track (reserved vs used).** JetStream distinguishes `memReserved/storeReserved` (sum of `MaxBytes`) from `memUsed/storeUsed` (actual `state.Bytes`) `:2376-2399,2628-2629`. This allows admission before data exists and reconciles drift via periodic store scans `:2236-2264`, trading a small consistency window for lock-free `isClusteredNoLock` checks `:2282`.

- **Channel semaphore for blocking I/O.** `diskIOSemaphore{ch chan struct{}}:27` caps concurrent file ops to avoid exhausting OS threads (`issue #2742:25-26`), sized `4..8192` and defaulting to `4096:21` rather than `NumCPU`; metrics `waiters/waits/maxWaitNanos:29-32` expose saturation without shedding.

- **Generic ipQueue with optional caps, but caps off by default.** `ipqueue.go:26-82` provides `ipqLimitByLen/Size` with fail-fast `errIPQLenLimitReached:84`, yet server instantiates `jsAPIRoutedReqs/sendq/recvq` without limits `:372,36,133`. This favours throughput/deadlock avoidance over boundedness for internal control planes.

- **Process-wide config-time sizing (75 % heuristic).** `dynJetStreamConfig:2764-2808` sets `MaxMemory = 3/4 * sysmem.Memory()` (or `GOMEMLIMIT` if lower) and `MaxStore = diskAvailable()` scaled to 75 % free. This is simple and prevents unbounded defaults, but does not adapt to runtime memory pressure or container cgroup updates.

- **TLS handshake rate limiting as token bucket over 1 s.** `rateCounter{limit, interval=1s:33}` with `allow()` counting and `countBlocked:58-64`, logged via `logRejectedTLSConns:990-1004`. This protects accept path without a global connection queue.

## Notable Patterns

- **Optimistic accounting with immediate revert on breach** — `c.out.pb += len(data)` before limit check, then `pb -= len(data)` on slow consumer `:2603` avoids TOCTOU between check and enqueue; pattern repeated for `ipQueue` `sz` in `pushMany` revert `:153-177`.

- **Stall channel as binary backpressure signal** — `c.out.stc = make(chan struct{})` created at 75 % occupancy `:2623` and `close(stc)` when `pb < 75 %` `:1934`; producers `select { case <-stall; case <-timer }` `:3743-3745`, a Go-idiomatic alternative to condition-variable broadcast.

- **Sync-pool buffer tiering** — Three fixed pools `512/4k/64k:365-388` with `nbPoolGet` best-fit `:393-402`; keeps `net.Buffers` allocation off the heap under fan-out (`server/client.go:349-361`) at cost of holding `cap` memory.

- **Ring buffer for closed-client history** — `closedRingBuffer` with `MaxClosedClients 10k:192` gives bounded retention for monitoring (`/connz?type=closed`) without growing with connection churn.

- **Atomic 64-bit counters with 32-bit safe alignment comment** — `jetStream` leads with `apiInflight/memMax/storeMax` `:110-119` plus note "first because of atomics on 32bit systems" for correct alignment, accessed via `atomic.LoadInt64/AddInt64:2394,2294`.

- **Rate-limited warning coalescing** — `rerrMu/rerrLast:344-345`, `rateLimitLogging map + chan:358-359`, and `RateLimitWarnf` calls `:7378,7646` avoid log flooding during resource-exceeded storms.

## Tradeoffs

- **Hard disconnect vs graceful shedding.** Killing slow consumers at `64 MiB` or `10 s` write deadline is simple and isolates the server, but shifts load-shedding burden to clients and causes thundering-herd reconnects under burst. A soft shed (drop oldest messages, apply queue deadline, return 503) would preserve connections but requires per-subject policy.

- **Reserved MaxBytes vs actual Bytes.** Reservation prevents overcommit at create time (`releaseStreamResources:2733`) but can block creation when reserved capacity is high while actual usage is low (e.g., many streams with `MaxBytes=1 TB` but near-empty). The opposite — no reservation — would admit more streams but risk OOM later.

- **Bounded semaphore (4096) vs throughput.** `dios` cap of 4096 concurrent file ops is far higher than old `NumCPU` heuristic `:45-49`, improving SSD parallelism at cost of higher file-descriptor and goroutine pressure; `min 4 / max 8192:22-23` bounds still allow tuning via `JetStreamConcurrentIOs` but not dynamically.

- **Unbounded internal queues for low latency vs overload safety.** Leaving `jsAPIRoutedReqs`, `sendq`, `recvq` unbounded avoids spurious `409` rejections under spike, but under sustained API flood they can grow to O(memory). Adding `ipqLimitByLen` would bound them but needs a defined shed policy (which `$JS.API.LimitReached` advisory `:377` hints at but does not enforce).

- **Static 75 % sizing vs cgroup-aware.** Using `sysmem.Memory()/4*3:2801` and `diskAvailable 75 %` at startup is cheap; it does not track `GOMEMLIMIT` changes or container memory limits after start, unlike a periodic `sysmem`/`cgroup` poll.

- **Micro-stall (10 ms total) vs global throttling.** Per-producer stall `2-5 ms, total 10 ms:125-127` smooths fan-in without global impact, but under many producers the aggregate stall budget (`producer.in.tst`) is per-producer, so aggregate load is not bounded.

## Failure Modes / Edge Cases

- **Oversized MaxBytes reservation blocks cluster.** Requesting many streams each with large `MaxBytes` sums `memReserved/storeReserved` and `wouldExceedLimits:2388` rejects further streams even when `memUsed` is small — deceptive "resource limits exceeded" until streams are deleted or `MaxBytes` is lowered (no automatic defragment).

- **Pending-bytes revert relies on single-thread lock.** `queueOutbound` assumes `c.mu` is held; if future code calls it without lock the `pb += / pb -=` revert `:2603` races with `flushOutbound pb -= n:1921`, causing under/over count and spurious slow-consumer or missed detection.

- **Store-drift correction is lazy.** `checkAndSyncUsage` runs on `go jsa.checkAndSyncUsage:2312` only when `local.mem <0` or periodically (1500 ms tick `:2316`); a leak that keeps `usage.local` positive but diverges from `store.FastState` can persist minutes before correction, allowing `wouldExceedLimits` to admit beyond real usage.

- **Disk full after dynamic sizing.** `diskAvailable` is probed at config time `:2787`; if the volume fills later (external writer, log rotation) JetStream writes fail with file errors (`server/filestore.go:424`) not with admission rejection; no proactive free-space shedding.

- **Unbounded ipQueue under API flood.** Without `msz/mlen` the API queues can exhaust heap; `push` returns `errIPQSizeLimitReached:85` only if a limit was set, so the default path has no backpressure signal to the caller — the caller would need to poll `len()/size()` out-of-band.

- **MaxPending/MemMax race on reload.** `reload.go:756-757` stores new `memMax/storeMax` via `atomic.StoreInt64` under `js.mu`; concurrent `reserveStreamResources:2717` and `wouldExceedLimits:2394` read atomically, so a reload to a smaller limit can transiently allow a stream whose reservation already passed the old limit but now exceeds the new one (no retroactive eviction).

- **TLS rate limit window edge.** `rateCounter` resets `count` when `now.After(end):42-44` using wall clock; a monotonic skew or 1 s burst at boundary can allow `2*limit-1` handshakes across the window.

- **Stall total budget bypass via multiple producers.** `stalledWait` increments `producer.in.tst` per-producer `:3748`, not per-consumer; 100 producers each with `10 ms` budget can collectively stall `1 s` on the same slow consumer without hitting `stallTotalAllowed` globally.

## Future Considerations

- **Add soft pre-OOM shedding.** Poll `sysmem.Memory`, `runtime.MemStats.HeapInuse`, and `debug.SetMemoryLimit`/`GOMEMLIMIT` periodically; before `heap > 85 % limit` shed low-priority work (trim `closed` buffer, drop non-essential advisories, reject new `MaxPending` growth with `423 MaxPendingExceeded` instead of waiting for 64 MiB).

- **Bound all ipQueues and expose shed metrics.** Construct `jsAPIRoutedReqs` with `ipqLimitByLen(10000)` + `ipqSizeCalculation`, `sendq` with size cap, and surface `ipQueuesz:1202` counters (`len`, `size`, `inProgress:323`) as actionable gauges with alerts, not vanity totals.

- **Per-account pending cap and fairness.** Introduce `Account.Limits.MaxPendingTotal` and `MaxOutstandingAPI` to prevent one account's slow consumers or API burst from consuming `N*64 MiB` — extend `queueOutbound` to check account-wide `acc.pendingBytes` with token-bucket.

- **Make disk reservation adaptive.** Re-probe `diskAvailable` on a timer (e.g., 30 s) and adjust `storeMax` down if free space < 10 %, emitting `JSStorageResourcesExceeded` advisory early, rather than only at startup.

- **SLO benchmarks for overload.** Add saturation tests: steady 50 k msg/s fan-out, 10x burst 1 M msg/s for 30 s, and 7-day long-lived session with retained streams; record `pb`, `stalls`, `slowConsumers`, `dios waitNanos/maxWaitNanos`, `heap` with stated hardware (CPU/RAM/disk) and repetitions — current benchmarks `core_benchmarks_test.go:32` and `filestore benchmarks:8623` are throughput-only and do not assert bounded memory.

- **Queue deadlines & load shedding policy for JetStream API.** Define `queueDeadline` (e.g., 500 ms) and `loadShedding` (return `503` when `q.len > 0.8*mlen`) for `$JS.API.>` handling, mirroring the client stall but at API layer.

## Questions / Gaps

- No evidence found for per-stream goroutine/pool isolation — `server.go:249-253` tracks `grWG/grTmpClients` globally but no bounded worker pool per stream/consumer; search for `workerPool`, `semaphore` (outside `dios`/`syncOutSem`/`rateCounter`) returned nothing.

- No evidence found for model/subprocess/stream memory eviction before OOM — `jetstream.go:2388-2403` only checks limits at create/publish time, not background reclamation; `memstore.go:35-75` holds `map[seq]*StoreMsg:39` unbounded except by `MaxBytes/MaxMsgs`.

- Benchmark hardware/workload retention not stated — `core_benchmarks_test.go:101-260` and `filestore_test.go:8623` do not log CPU, RAM, disk, or retained memory; need reproduction manifest.

- `ipQueuesz` endpoint exists `:1202` but no test asserts that queue depth stays bounded under burst; need saturation test harness.

- Admission vs actual usage reconciliation window: `usageTick 1500 ms:2316` and `minUsageUpdateWindow 250 ms:2334` — what is the worst-case drift allowed before `wouldExceedLimits` blocks? Not documented.

---

Generated by `01.09 Resource Accounting, Overload, and Bounded Work` against `nats-server`.
