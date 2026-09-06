# Source Analysis: buildkit

## Operation, Step, Attempt, and Process Identity

### Source Info

| Field | Value |
|-------|-------|
| Name | buildkit |
| Path | `studies/ultraplan-daemon-events-study/sources/buildkit` |
| Language / Stack | Go (gRPC, containerd, BoltDB, OpenTelemetry) |
| Analyzed | 2026-09-02 |

## Summary

BuildKit models durable execution as a content-addressed DAG, not as a classical operation/attempt hierarchy. The user intent is a `SolveRequest` identified by a client-supplied `Ref` (fallback `identity.NewID()`), scoped to a `SessionID` (`session.Group`). The solver (`solver.Solver`) creates a `Job` per `Ref` (`solver/jobs.go:373`) that references a shared graph of `state`/`edge` nodes keyed by `Vertex.Digest()` (`solver/types.go:17`, `solver/jobs.go:43-78`). The smallest checkpointable unit is an `edge` (`solver/edge.go:41-76`) with states `initial → cache-fast → cache-slow → complete` (`solver/edge.go:16`), whose result is a `CacheKey`/`CacheRecord` persisted via `CacheManager`/`CacheKeyStorage` and content leases. There is no first-class attempt/retry object: concurrent callers for the same vertex are coalesced by `flightcontrol.Group` inside `sharedOp` (`solver/jobs.go:995`, `solver/jobs.go:1219`) and by the scheduler's `edge` merge (`solver/scheduler.go:288`, `solver/edge.go:124`). Failures set `edge.err` and propagate via pipe cancellation (`solver/edge.go:651`), not via a new attempt ID. Process identity is ephemeral (`executor.ProcessInfo` / container task ID passed to `worker.Worker.Executor.Run`), and correlation relies on `Ref`, `SessionID`, `vertex digest`, `progress.Writer` ID, and OTel `trace.Span`/`MultiSpan` (`solver/jobs.go:61`, `solver/jobs.go:1374`), but late/duplicate executions cannot be quarantined by an attempt token.

## Rating

**6 / 10 — Present but inconsistent, weakly documented, fragile**

Rationale: Vertex/job/history layer is explicit and tested (scheduler 3000+ line test suite, cache-key depot, history lease protection). However the dimension's core requirement — distinct, durable attempt identity that survives retries and lets the system reject a late attempt 1 after attempt 2 succeeds — is absent. All retries collapse to the same vertex digest / flightcontrol key (`"exec"`, `solver/jobs.go:1219`); there is no monotonic attempt counter, no idempotency key, and progress events are multiplexed to every job sharing the vertex via `MultiWriter`/`MultiSpan`, which intentionally conflates attribution for deduplication efficiency.

## Evidence Collected

Every entry MUST include a file path with line numbers. Format: `path/to/file.ts:NN`.

| Area | Evidence | File:Line |
|------|----------|-----------|
| Run/job/task/attempt models | `Job` struct: `id`, `SessionID`, `uniqueID`, `startedTime`, `completedTime`, `span` — job is the durable solve instance; no Attempt struct exists | `solver/jobs.go:373-389` |
| Run/job/task/attempt models | `Solver.jobs map[string]*Job` + `actives map[digest.Digest]*state` — shared graph: jobs multiplex over content-addressed states | `solver/jobs.go:41-45` |
| Run/job/task/attempt models | `Vertex` interface (`Digest()`, `Sys()`, `Options()`, `Inputs()`, `Name()`) — smallest logical step definition | `solver/types.go:17-36` |
| Run/job/task/attempt models | `Edge` struct wrapping `Vertex` + `Index` — per-output task node | `solver/types.go:44-47` |
| Run/job/task/attempt models | `edge` internal state machine (`edgeStatusInitial/CacheFast/CacheSlow/Complete` + `result`, `cacheMap`, `keys`, `cacheRecords`) — checkpointable execution unit | `solver/edge.go:14-76` |
| Run/job/task/attempt models | `sharedOp` with `gDigest`, `gCacheRes`, `gExecRes flightcontrol.Group` — coalesces concurrent executions of same vertex to one `Exec` call | `solver/jobs.go:992-1006`, `solver/jobs.go:1219` |
| Run/job/task/attempt models | `Op` interface (`CacheMap`, `Exec`, `Acquire`) — worker-bound execution handler, resolved per vertex | `solver/types.go:173-183` |
| Run/job/task/attempt models | `llbsolver.Solver.Solve` creates job via `s.solver.NewJob(id)` where `id` is `SolveRequest.Ref`; defers `j.Discard()`; uses `j.UniqueID()` for provenance | `solver/llbsolver/solver.go:210-214`, `solver/llbsolver/solver.go:449` |
| Run/job/task/attempt models | `Controller.Solve` receives `SolveRequest.Ref`, increments `buildCount`, delegates to `llbsolver.Solver.Solve` | `control/control.go:409-412`, `control/control.go:575` |
| Database schemas | `BuildHistoryRecord` proto: `Ref`, `Frontend`, `CreatedAt`, `CompletedAt`, `Result/Results`, `Logs`, `Trace`, `Generation`, `pinned` — durable build record | `api/services/control/control.proto:203-225` |
| Database schemas | `bbolt` bucket `_records` storing `BuildHistoryRecord` VT-marshaled; bucket `_version` | `solver/llbsolver/history/buildhistory.go:38-39` |
| Database schemas | `CacheKey` (`ID`, `digest`, `vtx`, `output`, `deps [][]CacheKeyWithSelector`, `ids map[*cacheManager]string`) — content-addressed cache identity | `solver/cachekey.go:35-48` |
| Database schemas | `CacheRecord` (`ID`, `Size`, `CreatedAt`, `Priority`, `cacheManager`, `key`) — stored result handle | `solver/types.go:257-265` |
| Database schemas | `History.Queue.update` persists record + adds `content` resources under lease `ref_<Ref>` in isolated containerd NS `<ns>_history` | `solver/llbsolver/history/buildhistory.go:418-487`, `solver/llbsolver/history/buildhistory.go:104-105` |
| Parent/root ID fields | `state.parents map[digest.Digest]struct{}` + `childVtx map[digest.Digest]struct{}` — DAG parent links by vertex digest | `solver/jobs.go:52-56` |
| Parent/root ID fields | `state.jobs map[*Job]struct{}` — many jobs share same state (vertex deduplication) | `solver/jobs.go:52-53` |
| Parent/root ID fields | `provenanceBridge` / `subBuilder` + `prost` history: exporter `id := exporterVertexID(j.SessionID, i)` creates synthetic vertices for export/finalize | `solver/llbsolver/solver.go:415` |
| Parent/root ID fields | `llbBridgeForwarder` gateway: `workerRefByID map[string]*worker.WorkerRef` — maps leased refs back to originating build | `frontend/gateway/gateway.go:439`, `frontend/gateway/gateway.go:547` |
| Process/task handles | `executor.Executor.Run(ctx, id string, rootfs Mount, mounts []Mount, process ProcessInfo, started chan)` — `id` is task identity (container ID) | `worker/base/worker.go:80`, `worker/base/worker.go:376` (impl `proxyPolicyExecutor`) |
| Process/task handles | `cacheResultStorage.getWorkerRef(id)` parses `id` as `<workerID>/<refID>` to reconstruct `WorkerRef` | `worker/cacheresult.go:43-47` |
| Process/task handles | `WorkerRef.ID() => "<workerID>/<refID>"` delegates to `ImmutableRef.ID()` — result handle doubles as process artifact handle | `worker/result.go:21-26` |
| Correlation and causation metadata | `SolveRequest.Ref` (client id or `identity.NewID()`) + `SolveRequest.Session = s.ID()` — correlated across `Control` gRPC and solver | `client/solve.go:107-109`, `client/solve.go:322` |
| Correlation and causation metadata | `Job.SessionID` + `Job.uniqueID = identity.NewID()` — session group vs internal provenance invocation ID (`pr.RunDetails.Metadata.InvocationID = j.UniqueID()`) | `solver/jobs.go:387-388`, `solver/llbsolver/provenance.go:449` |
| Correlation and causation metadata | `session.Group` propagation: `state.Session()`, `ResolverCache`, `sessionGroup.NextSession()` tracing session lineage through parents | `solver/jobs.go:80-82`, `solver/jobs.go:142-196` |
| Correlation and causation metadata | `client.Vertex` (`Digest`, `Name`, `Inputs[]digest`, `ProgressGroup`) + `progress.Writer.Write(id, Vertex)` where `id := identity.NewID()` per state transition (`notifyStarted`) | `solver/jobs.go:1325-1332`, `solver/jobs.go:1368-1387` |
| Correlation and causation metadata | `tracing.MultiSpan` + `execSpan trace.Span` per `state`; `progress.MultiWriter` with `WithMetadata("vertex", dgst)` — causation via OTel span context | `solver/jobs.go:61-62`, `solver/jobs.go:583-584`, `solver/jobs.go:1042` |
| Correlation and causation metadata | `control.proto:StatusResponse { repeated Vertex, VertexStatus, VertexLog, VertexWarning }` — status stream keyed by `Vertex.digest` + `VertexStatus.ID`/`vertex` | `api/services/control/control.proto:120-165` |
| Retry-attempt creation code | No retry loop / attempt counter found; `solver` search for `retry/Retry` returns only unrelated `attemptUnpackDockerCompatibility` | `solver/pb/ops.proto:356` (negative evidence) |
| Retry-attempt creation code | `sharedOp.Exec` uses `gExecRes.Do(ctx, "exec", ...)` with `execDone/execRes/execErr` memoization — second caller waits on same `flightcontrol.call`, not a new attempt | `solver/jobs.go:1219-1230`, `solver/jobs.go:1004-1015` |
| Retry-attempt creation code | `edge` failure handling: `processUpdate` sets `e.err`, `respondToIncoming` finishes senders with `e.err`; cancellation via `pipe.Receiver.Cancel()` — no creation of `attempt2` edge | `solver/edge.go:650-656`, `solver/edge.go:622-632` |
| Retry-attempt creation code | `flightcontrol.Group.Do` implements singleflight with `errRetry` backoff but dedupes by key without generating distinct attempt IDs | `util/flightcontrol/flightcontrol.go:34-79` |

## Answers to Dimension Questions

### 1. What is the durable user intention?

The durable intention is a **Solve** keyed by `SolveRequest.Ref` (`api/services/control/control.proto:62`, `control/control.go:409`). Client generates `Ref = identity.NewID()` unless caller supplies `SolveOpt.Ref` (`client/solve.go:107-110`). It is persisted in `BuildHistoryRecord.Ref` (`api/services/control/control.proto:204`) in BoltDB (`_records` bucket, `solver/llbsolver/history/buildhistory.go:420-485`) and protected by a containerd lease `ref_<Ref>` in namespace `<ns>_history` (`solver/llbsolver/history/buildhistory.go:279`, `solver/llbsolver/history/buildhistory.go:104`). `Ref` survives `Job.Discard()` — the solver keeps jobs 10s after discard for late status polling (`solver/jobs.go:890-896`) and the history queue retains the record until GC (48h / 50 entries default, `solver/llbsolver/history/buildhistory.go:82-83`). `Session` (`SolveRequest.Session`) is a secondary ephemeral correlation for file/content sync, not the durable id.

### 2. What is the smallest checkpointable unit?

The smallest unit is an **`edge` (vertex + output index)** (`solver/edge.go:41`). The logical step is a `Vertex` (`solver/types.go:17`) whose identity is `Vertex.Digest()` (content hash of LLB definition + inputs). The scheduler creates one `edge` per `(Vertex, Index)` (`solver/jobs.go:217-218`). Checkpoint progression: `edge.status` moves `initial → cacheFast → cacheSlow → complete` (`solver/edge.go:16`), and each state caches `CacheMap`/`CacheKey`/`CacheRecord` via `CacheManager.Query/Records/Load/Save` (`solver/types.go:283-297`) backed by `CacheKeyStorage` + `containerd` content. Intermediate results are `CachedResult`/`SharedCachedResult` with lease-protected `WorkerRef`s (`solver/edge.go:108-114`, `worker/result.go:16-26`). Resumption Granularity: if a build is re-issued with the same LLB graph, cache hits stop re-execution at the edge level (`solver/edge.go:922-943` vs `solver/edge.go:974-1025`).

### 3. Does each retry get a distinct identity?

**No.** There is no retry/attempt entity and no attempt counter. Evidence of absence: `grep` for `retry/Retry/attempt` under `solver/` finds only `attemptUnpackDockerCompatibility` (`solver/pb/ops.proto:356`); no `Attempt` type, no `retryCount` field on `Job`/`state`/`edge`/`CacheKey`. Concurrency deduplication uses `flightcontrol.Group` keyed by static strings `"exec"` and `"cachemap-<idx>"` (`solver/jobs.go:1219`, `solver/jobs.go:1146`) — all callers join the same `call[T]` (`util/flightcontrol/flightcontrol.go:62-79`). On failure `sharedOp.execErr` is memoized (`solver/jobs.go:1006`) and future callers immediately receive the same error (or `errRetry` after cleanup, `util/flightcontrol/flightcontrol.go:138-139`), without a new diagnostic attempt ID. Edge merging (`solver/scheduler.go:288-328`, `solver/edge.go:124-128`) collapses duplicate vertices to a single `edge` with `releaserCount` ownership, intentionally losing per-requestor identity.

### 4. Can one logical operation span multiple runtime calls or OS processes?

**Yes — by design — but without attempt-scoped fencing.** One `Solve (Ref)` fans out to:
- multiple gRPC `Control.Status` streams reading `progress.MultiReader` (`solver/jobs.go:695-698`, `control/control.go:594-624`);
- scheduler `build(ctx, Edge)` which creates an `edgePipe` and signals the scheduler loop (`solver/scheduler.go:210-243`);
- gateway frontend forwarding that spawns a container task via `Executor.Run` with ephemeral `id` and `ProcessInfo` (`worker/base/worker.go:376`, `docs/dev/request-lifecycle.md:218`);
- export/finalize as synthetic vertices (`solver/llbsolver/solver.go:415-420`, `solver/llbsolver/solver.go:504-519`).
Correlation uses `Ref`, `SessionID`, `vertex digest`, and OTel span propagation (`state.mspan`, `state.execSpan`, `solver/jobs.go:61-62`, `solver/llbsolver/solver.go:504`). Cross-process lifetime is managed by leases (`leaseutil.WithLease`, `solver/llbsolver/solver.go:381`) and content GC, but there is no per-attempt execution token that would let the controller reject a stale container exit after a re-execution.

### 5. Can events be unambiguously attributed to the right entity?

**No — attribution conflates at the shared-vertex layer by design.** Progress and tracing intentionally multiplex:
- `state.mpw = progress.NewMultiWriter(WithMetadata("vertex", dgst))` (`solver/jobs.go:583`) and `state.mspan = tracing.NewMultiSpan()` (`solver/jobs.go:584`) aggregate all jobs sharing the vertex; `connectProgressFromState` fans out `pw.Write(identity.NewID(), clientVertex)` to every `Job.pw` sharing the state (`solver/jobs.go:665-685`), so status events (`api/services/control/control.proto:120-147`) lack a single-owning `Job.Ref` or attempt id.
- Scheduler `edge` status pushes `edgeState` via `pipe.Pipe` to all `incoming` senders (`solver/edge.go:731-804`), meaning multiple logical operations waiting on the same digest receive the identical `edgeState` update.
- History import (`Queue.ImportStatus`, `solver/llbsolver/history/buildhistory.go:679`) collapses per-vertex events into aggregate counts (`NumCachedSteps`, `NumCompletedSteps`, `NumTotalSteps`) without retaining per-job attribution.
- OTel spans are named by `vertex.Name()/Digest` (`solver/jobs.go:1042`, `solver/jobs.go:1240`) not by `(Ref, attempt)`, so two concurrent builds sharing a vertex generate indistinguishable spans.

**Counterfactual:** If attempt 2 succeeds after attempt 1 returns late, the IDs available are `vertex digest` (identical), `CacheKey.ID` (identical on success), and `progress id := identity.NewID()` per `notifyStarted` (`solver/jobs.go:1374`) which is a transient progress write ID, not a fencing token. No monotonic `attempt-id` or `execution-lease` distinguishes the straggler; `flightcontrol` would have already coalesced attempt 1 and 2, and the late error would either be dropped (already `execDone`) or overwrite shared `execErr` — there is no quarantine path.

## Architectural Decisions

| Decision | Evidence | Consequence |
|----------|----------|-------------|
| **Content-addressed vertex as primary identity** (`Digest()` of LLB op + inputs) | `solver/types.go:20`, `solver/jobs.go:545`, `solver/jobs.go:584-594` | Achieves aggressive deduplication across concurrent builds and persistent caching; but collapses distinct user operations that happen to compute the same content into one execution, erasing per-operation provenance. |
| **Job = ephemeral solve instance keyed by opaque `Ref`** (`Solver.jobs map[string]*Job`, `BuildHistoryRecord.Ref`) | `solver/jobs.go:41-43`, `solver/jobs.go:687-698`, `api/services/control/control.proto:62` | Lightweight checkpointing (bolt + lease) without a job state machine; requires caller to generate and remember `Ref`, no server-generated attempt sequence. |
| **Shared `state`/`edge` graph multiplexed across jobs** (`state.jobs map[*Job]struct{}`, `state.parents/childVtx`) | `solver/jobs.go:52-78`, `solver/jobs.go:257-296` | Enables zero-copy cache and parallel builds; progress/trace attribution necessarily fans out, violating unambiguous event routing. |
| **Flightcontrol singleflight for Op execution** (`sharedOp.gExecRes Group`) | `solver/jobs.go:992-997`, `solver/jobs.go:1219-1279` | Prevents thundering herd and duplicate resource acquisition; but serializes retries behind the same memo slot, with no attempt-scoped idempotency key. |
| **Lease-protected history in isolated containerd namespace** (`<ns>_history`, lease `ref_<id>`, `AddResource` for each blob) | `solver/llbsolver/history/buildhistory.go:104-105`, `solver/llbsolver/history/buildhistory.go:279-306` | Durable, GC-safe record of `Ref` with rich metadata; still lacks `attempt` child table or slot for per-retry output/error. |
| **Progress via `MultiWriter`/`MultiSpan` + `identity.NewID()` per `Vertex` write** | `solver/jobs.go:583-584`, `solver/jobs.go:1374`, `solver/jobs.go:675` | Simple live streaming to all waiters; writer ID is ephemeral and per-transition, not a durable causal edge. |
| **Deterministic merge on `CacheKey` collision** (`edge.index.LoadOrStore(k, e)` + `mergeTo`) | `solver/scheduler.go:156-177`, `solver/scheduler.go:289-329` | Further deduplicates identical subgraphs dynamically; `secondaryExporters` propagation is lossy (TODO cache provider merge absent). |

## Notable Patterns

- **Hierarchical DAG with flightcontrol coalescing:** `Scheduler.build → pipe → edge.unpark → op.CacheMap/Exec` with nested `Group.Do` at `sharedOp` layer (`solver/scheduler.go:210`, `solver/edge.go:331`, `solver/jobs.go:1137`).
- **Three-phase cache negotiation:** `cache-fast` (fast keys), `cache-slow` (content-digest via `ResultBasedCacheFunc`), `complete` (execute) — observable in `edgeStatusType` and `desiredStateDep` (`solver/edge.go:16-25`, `solver/edge.go:841-878`).
- **Lease-scoped artifact lifetime:** Single lease spans artifact creation + `finalize` (push) to prevent GC of blobs between export and push (`solver/llbsolver/solver.go:371-387`).
- **Synthetic export vertices:** `inBuilderContext` creates `client.Vertex{Digest: hash(id), Name: exporterName}` for exporter/finalize steps (`solver/llbsolver/solver.go:504-520`, `solver/llbsolver/solver.go:415`).
- **Idempotency by content, not by operation:** Cache hit is `CacheManager.Query(deps, inputIndex, digest, outputIndex)` (`solver/cachemanager.go:65`), keyed by op digest + selector chain, not by `(Ref, attempt)`.

## Tradeoffs

| Tradeoff | Pro | Con |
|----------|-----|-----|
| Maximize sharing (vertex digest) | High cache hit rate across concurrent/sequential builds; minimal storage | Loses per-build attribution; late duplicate results cannot be fenced by attempt |
| No attempt entity | Simpler state machine (4 statuses, no retry supervisor); less persistence overhead | No way to distinguish retry 1 vs 2; observability collapses; chaos tests must must invent external retry |
| Single memoized `sharedOp.execRes` | Exactly-once resource acquisition per buildkit daemon lifetime for a given digest | Stale failure poisons subsequent builds until `sharedOp` evicted via `deleteIfUnreferenced`; no per-attempt isolation |
| MultiWriter progress fan-out | Live `Status` streaming to many clients with no extra work | Event stream is the union of all sharing jobs; cannot audit which `Ref` triggered a specific `VertexLog`/`VertexWarning` |
| BoltDB + bbolt + content leases for history | Durable, file-level transactional history with GC and `pinned` flag | No secondary index on `SessionID`/`InvocationID`; history `Listen` must scan all records (`solver/llbsolver/history/buildhistory.go:821-841`) |

## Failure Modes / Edge Cases

| Mode | Symptom | Current Guard | Gap |
|------|---------|---------------|-----|
| Straggler completion after merged edge succeeded | Late `Exec` error arrives after `execDone=true` | `sharedOp.gExecRes` returns memoized success, late goroutine's result discarded via `hasActiveOutgoing` cancellation (`solver/scheduler.go:184-189`) | No `attemptId` to explicitly quarantine the late record; caller cannot tell if its own invocation succeeded |
| Forgotten `Ref` collision | `Solver.NewJob(id)` errors `job ID %s exists` (`solver/jobs.go:691-692`) | Caller-supplied `Ref` via `client/solve.go:107`; random 129-bit `NewID()` (`identity/randomid.go:43-51`) | Client that retries with same `Ref` after `Discard()` + 10s window gets `UnknownJobError` via `Get` timeout (`solver/jobs.go:718-743`) — not a clean retry flow |
| Cache poisoning on slow-cache failure | `CalcSlowCache` stores `slowCacheErr` persistently (`solver/jobs.go:1115-1122`) | Future `Query` will reuse the stored error | No attempt-scoped invalidation; a transient I/O error during content hashing blocks the vertex until daemon restart or state eviction |
| Progress leak on job discard | `Job.Discard` removes `j.pw` from `state.allPw` but edge may still hold `mpw` reference (`solver/jobs.go:879`) | 10s grace before `jobs[id]` deleted (`solver/jobs.go:891-895`) | Late status poll may read stale `mspan`/`mpw` with no fence |
| Lost trace on missing `Recorder` | `recordBuildHistory` warns `no trace recorder found, skipping` (`solver/llbsolver/history.go:281`) | Trace saving runs in background goroutine with 3s `ready` timeout (`solver/llbsolver/history.go:289-293`) | Trace retrieval is eventually consistent; `History.Queue.Status` may return before `Trace` descriptor committed (`solver/llbsolver/history/buildhistory.go:770`) |
| Edge merge cycle / deadlock | Scheduler bug leaves `openIncoming` without `openOutgoing` | Deadlock detector in `scheduler.dispatch` emits `markFailed` with diagnostic error (`solver/scheduler.go:184-189`, `solver/edge.go:374-379`) | No attempt isolation to retry the failed subgraph alone |

## Future Considerations

- **Introduce `Attempt` child entity under `Job`**: `Job { id, attempts []Attempt }` where `Attempt { id: identity.NewID(), seq int, state, result, startedAt, completedAt }`. Persist in bbolt as bucket `attempts_<Ref>` and content-lease per attempt; status stream includes `attemptId`.
- **Fencing token for execution:** pass `attemptID` as `executor.ProcessInfo` label + OTel baggage; `WorkerRef.ID()` could embed `attemptID` for GC; scheduler flightcontrol key becomes `(vertexDigest, attemptID)` when `forceAttempt` flag set, allowing true retry without coalescing.
- **Per-job progress routing:** replace `MultiWriter` fan-out with per-job `progressWriter` keyed by `(Ref, vertexDigest)` or add `attemptId` to `progress.Progress.meta` so `Status` consumers can filter.
- **Idempotency header:** accept `Idempotency-Key` in `SolveRequest` (like cloud APIs) mapped to `Ref:attemptSeq`; server returns same `Ref` on duplicate key without creating new build, and distinct key creates quarantinable attempt.
- **History schema v4:** add `attempts[]` to `BuildHistoryRecord` with per-attempt `error/status/logs/trace`; maintain `Generation` per attempt for causal ordering.

## Questions / Gaps

| Gap | Search Boundary | Why It Matters |
|-----|----------------|--------------|
| Is `Ref` uniqueness enforced durably across restarts or only in-memory `Solver.jobs`? | Searched `solver/jobs.go:687` (`jobs map`) and `history/buildhistory.go:538-549` (`active` vs `recordsBucket`); no cross-restart `Ref` index found | If a client retries after daemon restart with same `Ref`, whether it collides or orphans prior history affects recoverability |
| Does gateway `LLBBridgeForwarder` solve (`serveLLBBridgeForwarder`) create its own job/attempt IDs? | Read `frontend/gateway/gateway.go:453-706`; it reuses outer `sid` and `sessionID`, no new `Ref` observed | Determines if frontend-driven sub-solves have independent identity or share fate with parent `Ref` |
| Are container task IDs (`Executor.Run id`) recorded in history/provenance for post-mortem correlation? | Grepped `worker/`, `executor/` for `ProcessInfo` logging; task `id` not linked to `HistoryRecord` | Without task ID persistence, correlating host-side cgroup/executor logs to a vertex requires external join |
| Retry semantics for `Acquire` / resource exhaustion — does `Acquire` ever retry with backoff? | Inspected `sharedOp.Exec` (`solver/jobs.go:1227-1230`) — `Acquire` failure returns immediately | If worker resource contention is transient, lack of attempt-scoped retry forces whole build to fail |

---

Generated by `01.02-operation-step-attempt-process-identity` against `buildkit`.
