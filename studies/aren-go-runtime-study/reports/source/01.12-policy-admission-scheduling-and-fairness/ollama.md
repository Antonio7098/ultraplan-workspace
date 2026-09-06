# Source Analysis: ollama

## 01.12 Policy, Admission, Scheduling, and Fairness

### Source Info

| Field | Value |
|-------|-------|
| Name | ollama |
| Path | `studies/aren-go-runtime-study/sources/ollama` |
| Language / Stack | Go (Gin HTTP server, channel-based scheduler, llama-server subprocess per model) |
| Analyzed | 2026-08-30 |

## Summary

Ollama implements a single-node model-residency scheduler, not a cluster scheduler. All governance is concentrated in `server/sched.go` + `server/routes.go` around one scarce resource (VRAM / loaded-model slots). Admission is a two-stage gate — synchronous policy checks in `scheduleRunner` then a bounded channel enqueue in `getRunner` — with an immediate `503 ErrMaxQueue` on overflow rather than delayed queuing. Feasibility is predictive: `load` estimates VRAM via `PredictServerVRAM` against `free*80/100` headroom, selects placement via `selectLlamaServerPlacement`, and only then spawns `llama-server`; correction happens after failure through one-shot OOM retry (`reduceAutoNumCtxForLoadOOM` stepping `32768→4096` then `evictAllAndWait`). Dispatch is serialized through a single `processPending` goroutine and `activeLoading` single-slot, with `processCompleted` owning ref-count, TTL timers, and VRAM-reclamation. There is no priority, no per-tenant quota, no fair-share, and no workspace isolation — FIFO channel order + `ByDurationAndName` LRU eviction is the only ordering. Every effect is gated: no inference runs without an acquired `runnerRef` (`refCount++` + `finishedReqCh` release), and every rejection is observable as an HTTP status (400/404/503/500/499). The design is correct for a desktop model server but leaves Phase 13 gaps: unbounded waiting inside `processPending`, invisible queue depth, and client-owned retry after `ErrMaxQueue`.

## Rating

**4 / 10** — Bounded admission with explicit rejection and pre-effect VRAM checks, but no priority, fairness, starvation protection, per-tenant limits, or waiting-time bounds. Single-threaded pending loop provides ordering but not SLOs; OOM retry is bounded to one attempt; staleness detection is limited to `needsReload` + `Ping`. Useful as a model-residency pattern, insufficient as a governance substrate.

Rationale: strong on "no start before commit" (capability check → options commit → VRAM predict → runner pin), explicit `ErrMaxQueue` → `503` denial, and guaranteed ref-count release; weak on all fairness/tenant/quota dimensions, no queue-deadline or `Retry-After`, no starvation timer, and predictive rather than authoritative reservation.

## Evidence Collected

Every entry MUST include a file path with line numbers. Format: `path/to/file.ts:NN`.

| Area | Evidence | File:Line |
|------|----------|-----------|
| Bounded admission channel | `InitScheduler` sizes `pendingReqCh`, `finishedReqCh`, `expiredCh`, `unloadedCh` to `envconfig.MaxQueue()` (default 512) | `studies/aren-go-runtime-study/sources/ollama/server/sched.go:90-96` |
| MaxQueue default | `MaxQueue = Uint("OLLAMA_MAX_QUEUE", 512)` | `studies/aren-go-runtime-study/sources/ollama/envconfig/config.go:279` |
| Queue overload fast-fail | `select { case s.pendingReqCh <- req: default: req.errCh <- ErrMaxQueue }` with `ErrMaxQueue = errors.New("server busy... maximum pending requests exceeded")` | `studies/aren-go-runtime-study/sources/ollama/server/sched.go:88` , `studies/aren-go-runtime-study/sources/ollama/server/sched.go:204-208` |
| MaxRunners bound | `maxRunners := envconfig.MaxRunners()` then `if maxRunners>0 && loadedCount>=int(maxRunners) { runnerToExpire = s.findRunnerToUnload() }`; fallback `defaultModelsPerGPU=3` and auto `maxRunners = defaultModelsPerGPU * max(len(gpus),1)` | `studies/aren-go-runtime-study/sources/ollama/server/sched.go:226` , `studies/aren-go-runtime-study/sources/ollama/server/sched.go:265-268` , `studies/aren-go-runtime-study/sources/ollama/server/sched.go:86` , `studies/aren-go-runtime-study/sources/ollama/server/sched.go:280-291` |
| Policy gate before queue | `scheduleRunner` checks `model==nil`, `mllama` projector guard, `model.CheckCapabilities(caps...)`, then `usesAutomaticNumCtx/NumBatch` + `modelOptions` before `getRunner` | `studies/aren-go-runtime-study/sources/ollama/server/routes.go:202-232` |
| Fast-path reuse avoids queue | `getRunner` locks `loaded[modelKey]`, `if runner!=nil && !needsReload { req.useLoadedRunner(runner, finishedReqCh) } else { pendingReqCh <- req }` | `studies/aren-go-runtime-study/sources/ollama/server/sched.go:196-211` |
| Pending serialization | `processPending` single loop, `pending.schedAttempts++`, `if pending.ctx.Err()!=nil { continue }`, then inner `for` selecting expire/reload/load/evict | `studies/aren-go-runtime-study/sources/ollama/server/sched.go:225-245` |
| Stale-plan / needs reload | `runner.needsReload(ctx, req)` deep-equals adapters/projectors, compares `Runner` opts with `numCtxAuto/numBatchAuto/useMMapAuto` exemptions, checks `contextShift`, `Ping(ctx)` with 10s/2m timeout | `studies/aren-go-runtime-study/sources/ollama/server/sched.go:1381-1431` |
| Predictive feasibility (preflight) | `predictedForLoad > freeMemory*80/100 → return true` (signal eviction) else fits; logs `predicted/available/gpu_free/system_free` | `studies/aren-go-runtime-study/sources/ollama/server/sched.go:546-570` |
| Placement | `selectLlamaServerPlacement`: compact onto single GPU if `predictedVRAM <= 80% free` and `!SchedSpread`, else best backend group by `availableMemoryForLoad`; explicit `MainGPU` path | `studies/aren-go-runtime-study/sources/ollama/server/sched.go:972-1031` |
| Automatic context/batch correction | `reduceAutoNumCtxForLoadOOM` steps `nextLowerAutoNumCtx: 32768→4096` then reapplies `applyAutomaticGenerationBatch` (512/1024/2048 with headroom) | `studies/aren-go-runtime-study/sources/ollama/server/sched.go:757-782` , `studies/aren-go-runtime-study/sources/ollama/server/sched.go:807-835` |
| Single-slot loading | `loadedMu` + `activeLoading llm.LlamaServer` only one load at a time; `if llama==nil { newServerFn } else { panic if modelPath mismatch }` | `studies/aren-go-runtime-study/sources/ollama/server/sched.go:66-74` , `studies/aren-go-runtime-study/sources/ollama/server/sched.go:517-606` |
| Load → WaitUntilRunning | `gpuIDs, err := llama.Load(ctx, systemInfo, loadGpus, requireFull)` handles `ErrLoadRequiredFull`, `IsOutOfMemory` with `oomRetryAttempted` guard; then `llama.WaitUntilRunning` in goroutine with `s.expiredCh <- runner` on fail | `studies/aren-go-runtime-study/sources/ollama/server/sched.go:627-674` , `studies/aren-go-runtime-study/sources/ollama/server/sched.go:730-752` |
| Reservation via refCount | `useLoadedRunner: runner.refCount++`, stop `expireTimer`, `successCh <- runner`, `go { <-pending.ctx.Done(); finished <- pending }`; `processCompleted` decrements, arms `time.AfterFunc(sessionDuration)` or resets | `studies/aren-go-runtime-study/sources/ollama/server/sched.go:475-492` , `studies/aren-go-runtime-study/sources/ollama/server/sched.go:367-413` |
| Resource release | `expiredCh` handler checks `refCount>0` requeue, deletes `loaded[modelKey]`, calls `waitForVRAMRecovery` then `runner.unload()` + `unloadedCh <- {}` | `studies/aren-go-runtime-study/sources/ollama/server/sched.go:414-467` |
| VRAM reclamation wait | `waitForVRAMRecovery` polls `getGpuFn` every 250ms, succeeds on `freeBefore→freeNow > 75% vramSize`, timeout `s.waitForRecovery` (default 5s) then proceeds anyway | `studies/aren-go-runtime-study/sources/ollama/server/sched.go:1441-1499` |
| OOM retry ownership | `oomRetryAttempted bool` + `evictAllAndWait(ctx, pendingKey)` drain `unloadedCh` per runner; `expireRunnersForRuntimeOOM` expires all on `IsOutOfMemory` during `Completion/Embedding` | `studies/aren-go-runtime-study/sources/ollama/server/sched.go:38-40` , `studies/aren-go-runtime-study/sources/ollama/server/sched.go:1586-1661` |
| Victim selection | `findRunnerToUnload` sorts `ByDurationAndName` (sessionDuration numeric then modelKey lex), prefers `refCount==0` idle then shortest duration | `studies/aren-go-runtime-study/sources/ollama/server/sched.go:1663-1693` , `studies/aren-go-runtime-study/sources/ollama/server/sched.go:1557-1578` |
| Observable rejection | `handleScheduleError` maps `errCapabilities/errRequired→400`, `context.Canceled→499`, `ErrMaxQueue→503 ServiceUnavailable`, `os.ErrNotExist→404`, else `500`; `scheduleRunner` also maps `errCapabilityCompletion→400` | `studies/aren-go-runtime-study/sources/ollama/server/routes.go:3107-3120` , `studies/aren-go-runtime-study/sources/ollama/server/routes.go:467-473` |
| Expire via API | `POST /api/generate` with `prompt=="" && keep_alive==0` → `s.sched.expireRunner(m)` returns `DoneReason: "unload"`; same for `ChatHandler` `len(Messages)==0` | `studies/aren-go-runtime-study/sources/ollama/server/routes.go:403-414` , `studies/aren-go-runtime-study/sources/ollama/server/routes.go:2504-2515` |
| Per-runner concurrency | `sem *semaphore.Weighted` sized `numParallel` (arch-gated to 1 for `mllama/qwen3vl/etc`); `Completion` does `sem.Acquire(ctx,1)` | `studies/aren-go-runtime-study/sources/ollama/llm/llama_server.go:960` , `studies/aren-go-runtime-study/sources/ollama/llm/llama_server.go:1536-1542` , `studies/aren-go-runtime-study/sources/ollama/server/sched.go:497-510` |
| Load tests showing queue full | `TestSchedGetRunner` sets `OLLAMA_MAX_QUEUE=1`, second `getRunner` gets `ErrMaxQueue: server busy` | `studies/aren-go-runtime-study/sources/ollama/server/sched_test.go:520-534` |
| VRAM predictive tests | `TestSchedLlamaServerEvictsWhenVRAMInsufficient` (free=0 → `needEvict=true`), `FitsAlongside` (free 20 GiB → false), parallel context accounting | `studies/aren-go-runtime-study/sources/ollama/server/sched_test.go:1221-1372` |
| Cancellation while queued | `TestSchedAlreadyCanceled` pre-cancelled `dctx` → `pendingReqCh` drained with no `successCh/errCh` | `studies/aren-go-runtime-study/sources/ollama/server/sched_test.go:1191-1207` |

## Answers to Dimension Questions

### Can a request start before its policy and resources are committed?

No — every inference path is gated before effect, but the resource commit is predictive not authoritative.

Flow `trace`:
1. **Admission/policy** — HTTP handler parses `modelRef` then calls `scheduleRunner` (`server/routes.go:466`, `634`, `854`, `2634`). `scheduleRunner` (`server/routes.go:202-232`) denies synchronously: nil model → `errRequired`; vision guard; `CheckCapabilities` → `errCapabilities`; `modelOptions(WithEmbeddingBatchDefault)` merges model+request opts. Capability denial maps to `400 Bad Request` via `handleScheduleError` (`server/routes.go:3109-3110`), not queued.
2. **Queue ownership** — `getRunner` (`server/sched.go:169-211`) fast-paths a loaded runner (`useLoadedRunner`) else tries `pendingReqCh <- req` with `default: ErrMaxQueue→503` (`server/sched.go:204-208`). No request enters `processPending` without either immediate reuse or bounded enqueue.
3. **Feasibility before spawn** — `processPending` (`server/sched.go:225-360`) holds the pending request blocked on a single goroutine, refreshes `gpus` + `systemInfo`, calls `updateFreeSpace` (`server/sched.go:295`) to reconcile actual `VRAMByGPU`, then `load`’s preflight (`server/sched.go:547-570`) compares `predictedVRAM + batchSurcharge` vs `free*80/100`. If over threshold and `requireFull`, it returns `true` to evict rather than spawn.
4. **Reservation** — `load` holds `activeLoading` (`server/sched.go:517-523`) guaranteeing one spawn at a time, spawns via `newServerFn`, then `llama.Load` + `WaitUntilRunning` (`server/sched.go:627-737`). Only after `WaitUntilRunning` returns does `load` insert `runnerRef{refCount:1}` into `loaded` (`server/sched.go:718-728`) and signal `successCh`. `Completion`/`Chat` then run under `sem.Acquire(ctx,1)` (`llm/llama_server.go:1536`) and a pinned `refCount` (`server/sched.go:475-492`). The model-level effect (token generation) cannot start without that pin.
5. **Release** — `go { <-pending.ctx.Done(); finished <- pending }` (`server/sched.go:487-491`) guarantees `refCount--` in `processCompleted` (`server/sched.go:383-413`) and timer/expiry handling, even on client cancel.

Gap: the commit is *optimistic*. `PredictServerVRAM` and `availableMemoryForPlacement` are estimates; `load` still spawns then checks `ErrLoadRequiredFull` / `IsOutOfMemory` and retries once (`oomRetryAttempted`, `reduceAutoNumCtxForLoadOOM`, `evictAllAndWait` at `server/sched.go:649-669`). Success is only confirmed after `WaitUntilRunning` log parsing (`memoryParsingWriter` in `llm/llama_server.go`). A mispredicted spawn has observable cost (process start/stop) but never delivers tokens before commit — the failure is returned via `errCh` → `500`.

### What prevents starvation when long and short runs share capacity?

**Nothing.** Ollama has no priority, no fair-share, no aging, and no short-vs-long classification.

- **Queue order** is strict FIFO of the Go channel `pendingReqCh` (`server/sched.go:233`), which is drained by a single `processPending` goroutine. All pending requests block head-of-line behind the current one’s load/eviction/recovery sequence — including `waitForVRAMRecovery` (up to 5 s, `server/sched.go:1445`) and sequential `evictAllAndWait` draining `unloadedCh` per runner (`server/sched.go:1619-1628`). A long-running inference holding `refCount>0` pins its runner; `findRunnerToUnload` (`server/sched.go:1663-1693`) prefers idle runners but, if none idle, returns the shortest `sessionDuration` runner even though that runner is busy — the caller then marks `sessionDuration=0` and blocks on `unloadedCh` (`server/sched.go:338-358`), i.e., waits for that runner’s client to finish. There is no preemption.
- **Eviction order** `ByDurationAndName` (`server/sched.go:1557-1578`) sorts by `sessionDuration` numeric (treating negative `KeepAlive` as huge via `uint64` cast) then `modelKey` lex. It does not consider request age, expected duration, or historical share. The comment at `server/sched.go:1576` notes "In the future we can enhance the algorithm to be smarter".
- **Per-runner concurrency** is `numParallel` (`server/sched.go:497-510`, `llm/llama_server.go:960`) — embedding forced 1 (`server/sched.go:501-502`), some arches forced 1 — but this is intra-model parallelism, not inter-tenant fairness.
- No `priority` field exists on `LlmRequest` (`server/sched.go:29-58`), no `tenant/workspace` field, no quota map, and no tests for starvation (`sched_test.go` searched: 0 hits for `fair|starv|priority`).

Consequence: a burst of long `num_ctx=262144` loads can occupy all `MaxRunners` slots with large `KeepAlive` (including `KeepAlive=-1 → MaxInt64` at `envconfig/config.go:140`) and starve short `chat` requests indefinitely; the short requests pile on `pendingReqCh` until `ErrMaxQueue` turns starvation into visible `503`.

### Who owns retry after a rejected or stale scheduling plan?

**Split ownership, single server retry, client owns the rest.**

| Outcome | Who retries | Mechanism | Evidence |
|---------|------------|-----------|----------|
| **Queue full** (`ErrMaxQueue`) | **Client** | Server returns `503 ServiceUnavailable` via `handleScheduleError` (`server/routes.go:3113-3114`); no requeue. `pendingReqCh` is drop-on-full (`select default` at `server/sched.go:204`). Test asserts `server busy` (`server/sched_test.go:535`). | `server/sched.go:88,204-208` , `server/routes.go:3113-3114` |
| **Stale plan (needs reload)** | **Server inside `processPending`** | `runner.needsReload` (`server/sched.go:1381-1431`) returning true sets `runnerToExpire` → expire → wait `unloadedCh` → loop → `load` new runner. No Nack to caller. `activeLoading` panic guard at `server/sched.go:599-605` prevents concurrent stale load. | `server/sched.go:255-264,338-358,1381-1431` |
| **Feasibility miss (predicted overflow)** | **Server** | `load` preflight `predicted>free*80/100` → `return true` → caller evicts one runner and loops (`server/sched.go:547-560,304-328`). No client involvement. | `server/sched.go:546-560` |
| **Load OOM / requireFull** | **Server once** | `load` catches `ErrLoadRequiredFull` → return true (evict); `IsOutOfMemory` → if `reduceAutoNumCtxForLoadOOM` succeeds (step 32768→4096) set `oomRetryAttempted=true` and `return true`; else if `otherLoaded && !oomRetryAttempted` evict-all-and-retry once (`server/sched.go:627-669`). Guard `oomRetryAttempted` prevents loop (`server/sched.go:39-40`). `evictAllAndWait` drains all `unloadedCh` (`server/sched.go:1590-1630`). | `server/sched.go:627-669,757-782,1586-1630` |
| **Runtime OOM during inference** | **Server expires, client retries** | `expireRunnersForRuntimeOOM` (`server/sched.go:1632-1661`) on `Completion/Embedding/Chat` error matching `IsOutOfMemory`/`outOfMemorySubstrings` expires all runners (`sessionDuration=0`, `expiredCh <- runner` if idle). Next request will reload with corrected `auto` ctx/batch. No policy-recheck or backoff header. | `server/sched.go:1632-1661` , `server/routes.go:736,1001,2900,3027` , `llm/status.go` |
| **Transient spawn wait failure** | **Server** | `load`’s trailing goroutine on `WaitUntilRunning` error sends `errCh` and `expiredCh <- runner` (`server/sched.go:731-737`); that runner refCount stays 0 and is reclaimed via `processCompleted`. | `server/sched.go:730-740` |
| **Context cancel while queued** | **Server drops** | `processPending` checks `if pending.ctx.Err()!=nil { continue }` (`server/sched.go:237-240`); `TestSchedAlreadyCanceled` confirms dropped without `errCh/successCh` (`server/sched_test.go:1191-1207`). | `server/sched.go:237-240` |
| **MMProj OOM** | **Server once** | `llm/llama_server.go:1086-1115` `retryWithMMProjCPUOffload` retries once with `forceNoMMProjOffload` if `IsOutOfMemory` and projector offload was enabled. | `llm/llama_server.go:1086-1115` |

There is no `Ack/Nack` + `BlockedEvals`-like staging. The only server-side requeue is the OOM/evict loop inside `processPending`; all other failures are terminal for that attempt and surface as `500` (or `400/503/499`). No `Retry-After` header is set.

### Are queue depth and waiting time bounded and visible to callers?

**Depth bounded, time unbounded, visibility weak.**

- **Depth bound:** `pendingReqCh` capacity `= OLLAMA_MAX_QUEUE` (default 512 at `envconfig/config.go:279`) enforced as drop, not block. `finishedReqCh`, `expiredCh`, `unloadedCh` are identically sized (`server/sched.go:91-96`), so internal signals cannot deadlock on full queue either. `scheduleRunner`’s handlers surface `ErrMaxQueue` synchronously before any timer, because `getRunner` uses non-blocking send (`server/sched.go:204-208`). No caller ever waits unbound in a queue send — but they may wait unbound *after* enqueue.
- **Waiting time unbounded:** Once enqueued, a `pendingReqCh` entry is processed FIFO by a single goroutine (`server/sched.go:233`). There is no per-request deadline, no `Wait`/`WaitUntil`, no Nack timeout, no maximum dwell. The wait is the sum of: prior pending items × (each `load`/`WaitUntilRunning` up to `OLLAMA_LOAD_TIMEOUT` default 5 m at `envconfig/config.go:150` + `waitForVRAMRecovery` up to 5 s at `server/sched.go:1465` + `evictAllAndWait` draining N runners). `processPending` has no aging boost. Callers experience this as HTTP handler blocking on `select { case runner=<-runnerCh ; case err=<-errCh }` (`server/routes.go:223-229`) until `ctx` cancels. Cancellation is cooperative: `processPending` only checks `pending.ctx.Err()` once at dequeue (`server/sched.go:237`), not while blocked on `unloadedCh` (`server/sched.go:350-358`) — a cancelled waiter still holds its slot until the prior unload completes.
- **Visibility:** Weak.
  - Queue depth is not exposed via any metric or header. `PsHandler` (`server/routes.go:2271-2303`) lists *loaded* models with `ExpiresAt` but not `pendingReqCh` len. No `/api/stats`, no Prometheus gauge for `pendingDepth`, no `BrokerStats`-like emission.
  - Waiting time has no histogram; only `TotalDuration`/`LoadDuration` returned on successful `GenerateResponse`/`ChatResponse` (`server/routes.go:709-710,2803-2804`).
  - Rejection is visible as `503` with body `{"error":"server busy, please try again. maximum pending requests exceeded"}` (`server/routes.go:3114`), but without `Retry-After`. No advisory for approaching capacity.
  - Tests exercise only the binary `queue full` case (`server/sched_test.go:520-535`) and never assert dwell time or depth metrics.

## Architectural Decisions

- **Single-node channel scheduler instead of heap/raft broker.** Keep all placement in-process (`processPending`/`processCompleted` two goroutines at `server/sched.go:213-222`) with channels as the admission and lifecycle bus. Tradeoff: trivial correctness, no persistence, no multi-host plan. Alternatives like Nomad’s `EvalBroker` + `BlockedEvals` + `PlanQueue` survive restarts and schedule across nodes but pay Raft and heap-bookkeeping cost.
- **Drop-on-full (`select default`) over block-with-timeout.** Return `ErrMaxQueue→503` immediately (`server/sched.go:204-208`) rather than blocking the HTTP handler. Rationale: bound memory, avoid handler goroutine pile-up. Cost: client must implement backoff; no server-side `429 Retry-After` or queuing SLO (compare Nomad’s `Wait`/`WaitUntil` + `delayHeap`).
- **Serialized loading (`activeLoading` single slot).** Only one `llama-server` spawn at a time (`server/sched.go:69-74,517-606`) with `schedAttempts` per request (`server/sched.go:36`). Simplifies `updateFreeSpace` accounting and avoids concurrent VRAM overcommit. Cost: head-of-line blocking; throughput capped at one load at a time regardless of GPU count.
- **Predict-then-correct VRAM admission.** Estimate via `PredictServerVRAM` + `availableMemoryForPlacement` at `80%` headroom (`server/sched.go:547-550`), select placement via `selectLlamaServerPlacement` (`server/sched.go:972-1031`), correct via log-parsed `VRAMByGPU`/`memoryParsingWriter` and 75% reclamation wait (`server/sched.go:1486`). Rationale: driver-reported free memory lags; prediction is the only pre-spawn signal. Cost: misprediction triggers process restart + eviction churn.
- **Automatic `num_ctx`/`num_batch`/`useMMap` reduction as built-in retry.** `applyAutomaticGenerationBatch` (`server/sched.go:807-835`), `applyLlamaServerMmapDefaults`/`disableMmapForHostPressure` (`server/sched.go:1138-1160,1200-1241`), and `reduceAutoNumCtxForLoadOOM` (`server/sched.go:757-783`) encode policy that only applies when caller left the knob auto (`numCtxAuto`/`numBatchAuto`/`useMMapAuto` at `server/routes.go:176-197`). Explicit values are never coerced. Tradeoff: succeeds for default configs, fails closed for explicit large requests.
- **Ref-count + timer as residency lease.** `runnerRef{refMu, refCount, expireTimer, expiresAt, sessionDuration}` (`server/sched.go:1338-1364`) with `useLoadedRunner` bump and `finishedReqCh` drop (`server/sched.go:475-492`). `KeepAlive` from `OLLAMA_KEEP_ALIVE` (`envconfig/config.go:129-144`) including `KeepAlive=-1 → MaxInt64` (infinite). `findRunnerToUnload` prefers short TTL (`server/sched.go:1664-1692`). Tradeoff: simple TTL without heartbeat; infinite TTL can pin memory.
- **Bounded retry via `oomRetryAttempted`.** One `reduceAutoNumCtx` step then one `evictAllAndWait` (`server/sched.go:649-669`, `1586-1630`). No exponential backoff, no jitter. Compare Nomad’s `initialNackDelay`/`subsequentNackDelay` compounding — Ollama chose minimal retry because each retry is a heavyweight process restart.

## Notable Patterns

- **Channel trio as state machine bus.** `pendingReqCh` (admission), `finishedReqCh` (release), `expiredCh`/`unloadedCh` (reclamation) at `server/sched.go:61-65` partition ownership: `processPending` owns `pending→expired`, `processCompleted` owns `finished/expired→unloaded`, eliminating lock sharing except `loaded` map.
- **`updateFreeSpace` estimate reconciliation.** Before each feasibility check, subtract `llama.VRAMByGPU` per device from reported `FreeMemory` (`server/sched.go:1295-1335`) to counter stale driver reporting. Mirrors `nimbus` VRAM accounting but without a persistent ledger.
- **`needsReload` auto-vs-explicit exemption.** `if runner.numCtxAuto && req.numCtxAuto { optsNew.NumCtx = optsExisting.NumCtx }` (`server/sched.go:1399-1403`) and analogous for `NumBatch`/`UseMMap` (`server/sched.go:1402-1407`) avoids reload churn when server auto-tuned values drift.
- **Single-projector offload guard.** `mmprojMemoryRequirement` + `shouldDisableMMProjOffload` (`llm/llama_server.go:663-691,735-774`) plus one-shot CPU offload retry (`llm/llama_server.go:1086-1115`) — projector VRAM is reserved via `LLAMA_ARG_FIT_TARGET` pad (`llm/llama_server.go:694-732`) rather than bare prediction.
- **Http handler as admission façade.** `scheduleRunner` in `server/routes.go:202-232` concentrates all policy decisions and maps every scheduler error to a typed HTTP status in `handleScheduleError` (`server/routes.go:3107-3120`), so the scheduler itself never knows about HTTP.
- **Stale-GPU guard via `waitForVRAMRecovery`.** Poll `getGpuFn` post-unload until `75%` reclaimed or 5 s timeout (`server/sched.go:1441-1499`) before admitting next load — trades throughput for avoiding immediate re-OOM on driver lag.

## Tradeoffs

- **Throughput vs determinism.** Single `processPending` loop gives linearizable pending order with no lock contention but serializes all loads; multi-GPU hosts cannot load two models concurrently even when on disjoint devices.
- **Drop vs queue.** `ErrMaxQueue→503` bounds handler goroutines and memory but pushes backoff to clients with no `Retry-After`; a blocking queue with deadline would improve tail latency at cost of handler retention and timeout bookkeeping.
- **Prediction headroom vs utilization.** `80%` threshold (`server/sched.go:550,1065`) plus `60%/75%` batch headroom (`server/sched.go:887-895`) plus `mmapHostPressureHeadroom = max(8 GiB, total/10)` (`server/sched.go:1255-1260`) leaves slack to absorb driver inaccuracy but wastes VRAM on pessimistic fits; no adaptive tuning from observed `memoryParsingWriter` values.
- **Infinite KeepAlive vs eviction fairness.** `KeepAlive=-1`/`MaxInt64` (`envconfig/config.go:140`) is a legitimate pin; `ByDurationAndName` sorting treats it as `uint64` huge so it evicts last (`server/sched.go:1563`), but a single pinned large model can still starve all subsequent loads with no quota override.
- **Single OOM retry vs robust backoff.** One `reduceAutoNumCtx` step (`server/sched.go:917-926` → 32768→4096) covers the common auto-context overcommit; explicit large `num_ctx` following a genuine OOM gets no reduction (correctly fails fast) but genuine transient driver stalls have no jittered retry.
- **LLM semaphore per runner vs global concurrency cap.** `sem *semaphore.Weighted` sized `numParallel` (`llm/llama_server.go:960,1536`) isolates per-model parallelism but there is no global `MaxQueue`-like semaphore across models; OLLAMA_NUM_PARALLEL is per-model (`envconfig/config.go:275`).

## Failure Modes / Edge Cases

- **Pending head-of-line block while `unloadedCh` waits.** `processPending`’s `select { case <-s.unloadedCh: continue }` (`server/sched.go:350-358`) is not `ctx`-interruptible per pending request; a cancelled waiter still occupies the head slot until the in-flight unload + `waitForVRAMRecovery` finishes, delaying all later pending.
- **`expiredCh` send under `refMu` with bounded buffer.** `expireRunner` (`server/sched.go:1726-1727`), `evictAllAndWait` (`server/sched.go:1614`), and `processPending`’s expire (`server/sched.go:346`) all do `s.expiredCh <- runner` while holding `runner.refMu` and relying on `expiredCh` capacity 512. If `processCompleted` stalls in `waitForVRAMRecovery` (`server/sched.go:458-464`), `expiredCh` can fill and those sends block holding `refMu`, stalling `processPending` and API `expireRunner` calls (`server/routes.go:405`).
- **Duplicate/orphaned expired events.** `processCompleted` handles `runner.pid != loaded[modelKey].pid` by closing orphan and not deleting map or signaling `unloadedCh` (`server/sched.go:441-451`), and handles nil entry as duplicate (`server/sched.go:432-440`). Rapid cancel + reload races can still enqueue two `expiredCh` events for same key; de-dup relies on PID, not monotonic token.
- **`waitForVRAMRecovery` cargo-culting on integrated/Metal.** Bypass check at `server/sched.go:1445-1450` skips recovery for `len(gpus)==0 || !discreteGPUs || Metal single`, but mixed integrated+discrete lists still wait; `hasDiscreteGPU` is per-list, not per `runner.gpus` device.
- **`needsReload` Ping timeout masks transient blip as reload.** `Ping` uses 10 s (or 2 m if `loading`) context (`server/sched.go:1386-1421`); a single failed `GET /health` (e.g., `ServerStatusNoSlotsAvailable` transient) forces evict+reload of a healthy runner.
- **Automatic reduction only on first OOM.** `reduceAutoNumCtxForLoadOOM` requires `req.numCtxAuto==true` (`server/sched.go:758-760`); an auto-context request that OOMs once and is retried after context reduction to 4096 that still OOMs gets `errCh <- err` with no further step-down (`server/sched.go:671-672`), even though 0 could be valid (tests show next lower after 4096 returns `ok=false` at `server/sched.go:917-926`).
- **Queue-depth metric absence.** No `BrokerStats`/`QueueStats.Depth` equivalent; operators cannot alert on `pendingReqCh` approaching 512 until clients see `503`.
- **No per-workspace or authz isolation.** `Server` (`server/routes.go:97-103`) has no tenant ID; `auth.go` only covers registry pull tokens (`server/auth.go:53-99`), not inference admission; a single tenant can fill `pendingReqCh` and `MaxRunners` and deny others.
- **No starvation test.** Searched `sched_test.go`: no test asserts FIFO starvation, priority inversion, or long-vs-short ordering; only `TestSchedGetRunner` shows `503`.

## Future Considerations

- **Add waiting-time bound and `Retry-After`.** Wrap `pendingReqCh` entries with `enqueueAt` and reject via `select { case <-time.After(maxWait): errCh <- ErrMaxQueue }` or set `Retry-After: <sec>` on `503` at `server/routes.go:3113-3114`. Expose `QueueStats.Depth` and `wait_time` histogram akin to Nomad’s `BrokerStats` for alerting.
- **Bounded admission on VRAM depth.** Surface `availableMemoryForPlacement`/`predictedVRAM` as preflight denial with structured `InsufficientVRAM` error (currently only logs at `server/sched.go:551-569`) and return `429` with estimate, so clients need not poll via `503`.
- **Priority + aging queue.** Extend `LlmRequest` (`server/sched.go:29-36`) with `priority int` and replace FIFO channel with priority heap (tie-break by enqueue time); add aging boost after `wait > threshold` to mitigate starvation — leverage existing `schedAttempts` counter as age signal.
- **Per-workspace / per-tenant quotas.** Propagate caller identity from `Server` (add auth middleware, map token→tenant) to `getRunner`, maintain `map[tenant]int` for `pendingCount` and `loadedCount` and deny with `429 quota exceeded` before `pendingReqCh` admission.
- **Starvation-aware victim selection.** Enhance `findRunnerToUnload` (`server/sched.go:1663-1693`) to consider waiting set’s total predicted VRAM and prefer victim that actually frees enough for the head waiter, or timeout-bound wait so short jobs don’t wait behind one long unload.
- **Unblock cancellation while waiting for `unloadedCh`.** Select `case <-pending.ctx.Done(): continue` alongside `case <-s.unloadedCh` inside `processPending`’s eviction wait (`server/sched.go:350-358` / `1590-1630`) to free the head slot when the waiter gives up.
- **Ack token + Nack timeout analogous to Nomad’s `token`/`NackTimer`.** Issue a dequeue token at `processPending` pickup and add `Ack`/`Nack` to close the window where a cancelled request still occupies the head until unload finishes; cures the pending-blame window found in failure modes.
- **Soak benchmark for VRAM churn.** Add hardware-pinned saturation test that asserts `pendingDepth` bounded, heap bounded, and VRAM reclaimed to `≥75%` within budget over burst / steady-state, covering the `waitForVRAMRecovery` timeout path.

## Questions / Gaps

- **What is the intended SLO for queue waiting vs fast-fail?** `envconfig.MaxQueue` is configurable but no doc states whether `512` intends burst absorption or strict shedding. Asked because `pendingReqCh` is FIFO with no deadline — the SLO determines whether to add `Wait`/`WaitUntil` or keep drop.
- **Is `MaxRunners` per host or per GPU intended to isolate tenants?** `OLLAMA_MAX_LOADED_MODELS` (`envconfig/config.go:277`) defaults to `3 * gpuCount` but has no tenant scoping. No evidence of per-tenant eviction preference was found; searched `quota`, `tenant`, `workspace`, `namespace` across `server/sched.go` and `server/routes.go` — only `workspace` hits were UI remotes, not admission.
- **Should `ContextLength` env-tier default be considered part of policy?** `totalVRAM >=47 GiB → 262144` (`server/routes.go:2050-2057`) sets the effective `numCtxAuto` that later drives `reduceAutoNumCtxForLoadOOM`; it is now a hidden policy knob not traced to admission tests.
- **No stale-plan fencing token.** Search for `SnapshotIndex`/`RefreshIndex`/`token` (Nomad-style) returned 0 in `server/`; Ollama uses `pid` + map key equality only (`server/sched.go:441`). Whether a stale `expiredCh` can unblock a newer pending awaits a fencing token is unanswered.
- **Workspace-scoped policy paths were not inspected per isolation rule;** the study was restricted to `studies/aren-go-runtime-study/sources/ollama` — sibling-source workspace config providers, `envconfig` remote allow-list (`Remotes` at `envconfig/config.go:166-175`), and auth registry flows (`server/auth.go`, `server/cloud_proxy.go`) were only inspected where they intersect inference admission (found no per-workspace quota).

---

Generated by `dimensions/01.12-policy-admission-scheduling-and-fairness.md` against `ollama`.
