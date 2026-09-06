# Source Analysis: buildkit

## Leases, Heartbeats, Fencing, and Stale-Worker Rejection

### Source Info

| Field | Value |
|-------|-------|
| Name | buildkit |
| Path | `studies/ultraplan-daemon-events-study/sources/buildkit` |
| Language / Stack | Go — single-daemon solver + containerd leases/content + BoltDB (bbolt) + gRPC control plane |
| Analyzed | 2026-09-02 |

## Summary

BuildKit has no distributed worker-fencing subsystem. Execution ownership is process-local: a `Solver` holds an in-memory `jobs` map and `actives` content-addressed vertex graph (`sources/buildkit/solver/jobs.go:42-44`), dispatching work through a single scheduler and coalescing duplicate work with `flightcontrol.Group` (`sources/buildkit/solver/jobs.go:992-997`). Leases exist only as containerd GC-protection (`sources/buildkit/util/leaseutil/manager.go:13-39`, `sources/buildkit/vendor/github.com/containerd/containerd/v2/core/leases/lease.go:43-47`) with a `containerd.io/gc.expire` timestamp label (`sources/buildkit/vendor/github.com/containerd/containerd/v2/core/leases/lease.go:93-102`) and a `buildkit/lease.temporary` label (`sources/buildkit/util/leaseutil/manager.go:83-89`), not as ownership leases. There is no heartbeat loop, no lease renewal, no epoch/generation CAS on job or cache mutations, and the `BuildHistoryRecord.Generation` field (`sources/buildkit/api/services/control/control.proto:215` / `sources/buildkit/solver/llbsolver/history/buildhistory.go:335`) is an unconditional counter, not a fencing token. Stale-worker rejection is therefore absent; the daemon instead prevents duplicate execution by collapsing it (`sharedOp.gExecRes.Do(ctx,"exec",…)` at `sources/buildkit/solver/jobs.go:1219-1220`) and relies on single-writer + context cancellation. Pausing an old worker and resuming after its job is discarded does not corrupt authoritative history via a checked fence, but can overwrite history last-write-wins and can race with GC via un-renewed lease expiry.

## Rating

**3 / 10 — Absent / unsafe for distributed ownership**

Lease is a GC root, not an ownership lease. No heartbeat, no renewal, no expiry-driven revocation, no fencing generation, no CAS ownership token on any product-state or artifact mutation, and no stale-writer test. The `Generation` bump is gratuitous. The design is intentional for a single-daemon build cache, but is unsafe if BuildKit were used as a multi-worker orchestrator.

Rationale maps to rubric 1-3 (absent, implicit, ad-hoc, unsafe): ownership is implicit in the in-memory `Solver.mu` lock (`sources/buildkit/solver/jobs.go:687-712`), expiry rules are clock-based labels never revalidated before writes, fencing is absent, takeover has no epoch, and late-completion behavior is "last writer wins" or "cached first writer wins" depending on layer.

## Evidence Collected

Every entry MUST include a file path with line numbers. Format: `path/to/file.ts:NN`.

| Area | Evidence | File:Line |
|------|----------|-----------|
| Lease / lock schema | `type Lease struct { ID, CreatedAt, Labels }` and `type Resource { ID, Type }` — generic containerd GC handle, no owner field | `sources/buildkit/vendor/github.com/containerd/containerd/v2/core/leases/lease.go:42-54` |
| Lease manager interface | `Manager { Create, Delete, List, AddResource, DeleteResource, ListResources }` — no Renew, no Heartbeat, no CAS | `sources/buildkit/vendor/github.com/containerd/containerd/v2/core/leases/lease.go:32-39` |
| Lease expiry rule | `WithExpiration(d)` writes `Labels["containerd.io/gc.expire"] = now+d RFC3339` — wall-clock label, evaluated by containerd GC, never checked before writes | `sources/buildkit/vendor/github.com/containerd/containerd/v2/core/leases/lease.go:93-102` |
| BuildKit wrapper lease creation | `NewLease` does `lm.Create( WithRandomID, WithExpiration(1h), opts...)` — 1h default expiry, no renewal loop | `sources/buildkit/util/leaseutil/manager.go:31-39` |
| Temporary marker | `MakeTemporary` adds `Labels["buildkit/lease.temporary"]=now RFC3339Nano` — GC hint only | `sources/buildkit/util/leaseutil/manager.go:83-89` |
| WithLease helper | `WithLease(ctx, ls, opts…)` reuses existing lease if `leases.FromContext(ctx)` present, else creates one; returns `Delete` closer, no heartbeat goroutine | `sources/buildkit/util/leaseutil/manager.go:13-29` |
| LeaseRef.Adopt | Copies resources to new lease via `ListResources` + `AddResource` loop, then async `go Discard` — no atomic move, no fencing, racy if source deleted concurrently | `sources/buildkit/util/leaseutil/manager.go:54-81` |
| History lease namespace | `h.hLeaseManager = opt.LeaseManager.WithNamespace(ns+"_history")`, `h.hContentStore = opt.ContentStore.WithNamespace(ns+"_history")` — isolated GC domain for history blobs | `sources/buildkit/solver/llbsolver/history/buildhistory.go:98-105` |
| History lease ID | `func leaseID(id string) { return "ref_" + id }` — deterministic, reused across updates | `sources/buildkit/solver/llbsolver/history/buildhistory.go:279-281` |
| History blob short-lease | `OpenBlobWriter` creates `WithRandomID, WithExpiration(5m), MakeTemporary` then `leases.WithLease(ctx,l.ID)` for content ingest | `sources/buildkit/solver/llbsolver/history/buildhistory.go:567-579` |
| Solver lease for exports | `solver.Solve` creates single `leaseutil.WithLease(ctx, lm, MakeTemporary)` spanning export+cache-export+finalize to prevent GC between create and push | `sources/buildkit/solver/llbsolver/solver.go:381-386` |
| Generation field definition | `int32 Generation = 12` on `BuildHistoryRecord` | `sources/buildkit/api/services/control/control.proto:215` |
| Generation increment | `br.Generation++` unconditionally after `upt(&br)` inside `UpdateRef`, no expected-value check, no CAS predicate | `sources/buildkit/solver/llbsolver/history/buildhistory.go:335` |
| Conditional-update absence (history) | `h.update()` does `h.hLeaseManager.Create(WithID(h.leaseID(rec.Ref)))` with `ErrAlreadyExists` → `l = Lease{ID:…}` fallback, then `b.Put(ref, marshaled)` — unconditional BoltDB put, `created` delete-on-error cleanup only (`if err!=nil && created {Delete}`) | `sources/buildkit/solver/llbsolver/history/buildhistory.go:418-486` |
| Job ownership claim | `NewJob(id)` locks `jl.mu`, checks `if _,ok:=jl.jobs[id]; ok {err exists}`, else inserts with `identity.NewID()` uniqueID; `Discard()` sleeps 10s then deletes `jl.jobs[id]` | `sources/buildkit/solver/jobs.go:687-712`, `sources/buildkit/solver/jobs.go:890-897` |
| Job Get wait | `Get(id)` with 6s timeout + `updateCond.Wait()` polling — blocking wait for insertion, not lease/heartbeat validation | `sources/buildkit/solver/jobs.go:717-743` |
| Content-addressed sharing (anti-fencing) | `loadUnlocked` keys `actives[dgst]` by `Vertex.Digest()` and merges jobs into same `state`; `dgstWithoutCache` variant for `IgnoreCache` | `sources/buildkit/solver/jobs.go:530-663` |
| Execution coalescing | `sharedOp { gDigest, gCacheRes, gExecRes flightcontrol.Group }` + `gExecRes.Do(ctx,"exec", func…)` guarantees exactly one `op.Exec` per digest; `execDone`+`execRes` cache prevents re-exec | `sources/buildkit/solver/jobs.go:992-997`, `sources/buildkit/solver/jobs.go:1219-1285` |
| Edge merging (ownership alias) | `solver scheduler edge.index.LoadOrStore` + `edge.takeOwnership(old)` increments `releaserCount` and aliases `old.owner=e` | `sources/buildkit/solver/edge.go:124-128`, `sources/buildkit/solver/scheduler.go:156` (grep evidence) |
| Session liveness | `session.Manager` tracks `sessions map[string]*client` under `mu`, `handleConn` inserts then `<-c.ctx.Done()` then deletes; no heartbeat, no TTL, liveness = gRPC stream open | `sources/buildkit/session/manager.go:29-145` |
| GC loops (no heartbeat) | `history.Queue` starts `gc()` every 120s and `clearOrphans()` on start; `Controller` throttles GC `Throttle(1m)` and `ReleaseUnreferenced 5m` — periodic janitor, not per-lease heartbeat | `sources/buildkit/solver/llbsolver/history/buildhistory.go:131-137`, `sources/buildkit/control/control.go:143-145` |
| History GC rule | `gc()` requires both `len(records) >= MaxEntries(50)` and `now - MaxAge(48h) > CompletedAt` then `delete` | `sources/buildkit/solver/llbsolver/history/buildhistory.go:152-198` |
| Deleted-ref refcount | `delete()` defers actual BoltDB delete if `h.refs[ref] >0` (Listen refcount) via `deleted` map | `sources/buildkit/solver/llbsolver/history/buildhistory.go:241-265` |
| UpdateRef pin path | `Controller.UpdateBuildHistory` → `history.UpdateRef` → `Generation++` → `update()` → unconditional `Put`; `Pinned` toggle is only guarded field, no token | `sources/buildkit/control/control.go:360-368`, `sources/buildkit/solver/llbsolver/history/buildhistory.go:309-348` |
| Late trace update | `history.go:318-327 UpdateRef(context.TODO(), id, func(rec) {rec.Trace=desc})` runs asynchronously 3s after STARTED, with `context.TODO()` — no lease check, last-write-wins after `Update COMPLETE` already wrote `Generation=N+1` | `sources/buildkit/solver/llbsolver/history.go:285-331` |
| Cache export lease | `runCacheExporters` + `runExporters` share caller-provided `ctx` with `leaseutil.WithLease` — resolver via `leases.FromContext` inside content store; no fencing before `CacheManager.Save` | `sources/buildkit/solver/llbsolver/solver.go:389-430`, `sources/buildkit/cache/remote.go:40` |
| Tests for stale writers | `grep -r "stale"` across `sources/buildkit/solver` returns only `cachemiss_stale_dep_test.go` (cache-key staleness, not worker fencing); `grep -r "fenc\|heartbeat"` returns zero non-vendor hits | `sources/buildkit/solver/cachemiss_stale_dep_test.go:9`, `sources/buildkit/util/leaseutil/manager.go:13` (negative evidence) |
| Takeover/reconciliation | `clearOrphans` scans `h.hLeaseManager.ListResources` per record, deletes records with 0 resources; not an ownership takeover, just orphan blob GC | `sources/buildkit/solver/llbsolver/history/buildhistory.go:200-239` |

## Answers to Dimension Questions

### 1. Can two workers believe they own the same work?

**Yes — by design, and they are not fenced; they are collapsed.**

BuildKit is single-daemon; "workers" are in-process `Solver` goroutines. Ownership is claimed per `Job` via `Solver.NewJob(id)` (`sources/buildkit/solver/jobs.go:687-712`) keyed by caller-supplied `SolveRequest.Ref` (`sources/buildkit/solver/llbsolver/solver.go:210`, `sources/buildkit/control/control.go:575`). Two concurrent `Solve` calls with the same `Ref` cannot coexist — the second gets `job ID exists` error — but two different `Ref`s whose LLB graphs hash to the same `Vertex.Digest()` intentionally share a `state` object (`sources/buildkit/solver/jobs.go:560-594`) and coalesce execution to a single `Exec` via `sharedOp.gExecRes.Do(ctx,"exec", …)` (`sources/buildkit/solver/jobs.go:1219-1220`).

Edge merging further aliases ownership: `edge.takeOwnership` (`sources/buildkit/solver/edge.go:124-128`) sets `old.owner=e` and `releaserCount` bookkeeping so that `release()` on either edge touches the merged target. This is not fencing; it is deduplication. The `uniqueID` (`sources/buildkit/solver/jobs.go:707`) assigned per `Job` is never checked before mutations — it is used only for provenance tagging (`sources/buildkit/solver/jobs.go:913`). There is no epoch/generation exchanged between workers, no compare-and-swap on `actives`, and no distributed lock. The effect is that concurrent solves free-ride on the first execution's result rather than fencing it.

### 2. Does lease expiry alone prevent stale writes?

**No. Leases do not gate writes, and expiry is not polled before mutations.**

A lease in BuildKit is a containerd GC lease (`sources/buildkit/vendor/github.com/containerd/containerd/v2/core/leases/lease.go:42-47`) whose only runtime effect is `Labels["containerd.io/gc.expire"]` (`sources/buildkit/vendor/github.com/containerd/containerd/v2/core/leases/lease.go:93-102`). The GC may delete content not reachable from a non-expired lease, but nothing in `history.Queue.update` (`sources/buildkit/solver/llbsolver/history/buildhistory.go:418-486`) or `solver.Solver.Solve` checks lease existence/expiry before `b.Put` or `CacheManager.Save`. 

Leases are created with fixed wall-clock expiries — 1h via `NewLease` (`sources/buildkit/util/leaseutil/manager.go:32`), 5m via `OpenBlobWriter` (`sources/buildkit/solver/llbsolver/history/buildhistory.go:568`), 24h in vendor client (`sources/buildkit/vendor/github.com/containerd/containerd/v2/client/lease.go:40-41`) — and never refreshed. There is no heartbeat goroutine in `WithLease` (`sources/buildkit/util/leaseutil/manager.go:13-29`): it returns a `Delete` closer, not a renewal ticker. A build running longer than the lease's wall-clock expiry would still hold `leases.WithLease(ctx, id)` in its context, but the underlying lease row could be GC'd out-of-band; conversely, holding the lease context after logical ownership loss does not prevent `b.Put` in `UpdateRef` (`sources/buildkit/solver/llbsolver/history/buildhistory.go:335-485`) because that path creates or reuses the `ref_<Ref>` lease without checking the caller's original temporary lease. Expiry is advisory for GC, not a revocation fence for mutations.

### 3. Is fencing checked by product-state and artifact mutations, not only event writes?

**No fencing is checked anywhere.**

Search for fencing generations, CAS, or ownership tokens across non-history state:

* `solver/jobs.go` job insertion is a mutex-guarded map check, not a CAS on `Generation`.
* `solver/edge.go` `Cache()` writes via `CacheManager.Save` (`sources/buildkit/solver/edge.go:998-1002`) with no generation; `CacheMap` and `CalcSlowCache` are memoised via `flightcontrol.Group` not via versioned records.
* `cache/manager_test.go` and `cache/remote.go` use `leaseutil.WithLease` for content GC only.
* `exporter/containerimage/export.go:242`, `exporter/oci/export.go:156` acquire a temporary lease for blob creation but never verify it before `content.Writer.Commit`.
* `solver/llbsolver/history.go:270-276` writes `BuildHistoryEventType_COMPLETE` with fully populated `CompletedAt/Logs/Results` then the async trace writer at `history.go:318-327` overwrites `Trace` via a separate `UpdateRef` with `context.TODO()` — neither checks `Generation`.

The sole `Generation` field (`sources/buildkit/api/services/control/control.proto:215`) is incremented unconditionally in `UpdateRef` (`sources/buildkit/solver/llbsolver/history/buildhistory.go:335`) after reading under `h.mu`, then written with `b.Put`. There is no `if Generation==expected` predicate on the `DB.Update` transaction, so concurrent callers both read `G`, both write `G+1`, and last writer silently overwrites. History, cache, and content are all unfenced.

### 4. How is process identity distinguished from reusable PIDs?

**It is not. Process identity is a caller-supplied `Ref` plus an ephemeral random `uniqueID`; neither forms a fencing token.**

* `SolveRequest.Ref` (`sources/buildkit/api/services/control/control.proto:62`) is caller-supplied (or server-generated if empty) and identifies a build invocation. Reuse after `Job.Discard` (10s delayed delete at `sources/buildkit/solver/jobs.go:890-897`) is allowed; the same string can be re-inserted with a new `Job` object.
* `Job.uniqueID = identity.NewID()` (`sources/buildkit/solver/jobs.go:707`) is an opaque random ID stored per `Job`, used for provenance grouping (`sources/buildkit/solver/jobs.go:913-914`) and never presented to clients as a lease token.
* There is no daemon `serverID`, no worker `instanceID`, and no `attempt` ordinal attached to requests. Session identity `SessionID` (`sources/buildkit/solver/jobs.go:387`, `sources/buildkit/session/manager.go:29-33`) is per-client gRPC stream, not per-attempt.
* OS PIDs are irrelevant — execution runs inside the daemon's `Executor` (runc/containerd shims) but the `Process` handle is not minted into any durable record.

Because identities are neither attached to mutations nor validated on writes, a rebooted daemon with the same `Ref` cannot be distinguished from the prior holder. `core/leases/lease.go:42` `Lease.CreatedAt` is the closest monotonic marker, but BuildKit never exposes it as a fence.

### 5. What happens to a late completion from a superseded attempt?

**There is no "superseded attempt" entity; behavior depends on layer:**

* **Scheduler / exec:** No supersession — `sharedOp.gExecRes.Do` (`sources/buildkit/solver/jobs.go:1220`) ensures the first `Exec` result is cached in `execRes/execDone` (`sources/buildkit/solver/jobs.go:1004-1015`). A resumed "old worker" thread blocked in `op.Exec` would either return before `execDone` flips and its result would be adopted, or after `execDone` its result would be discarded by the `flightcontrol` group (second `Do` call returns the cached value, not the late one). Late completions cannot overwrite; they are deduped, not fenced via generation.

* **History / authoritative state:** Last-write-wins with silent overwrite. Example: `recordBuildHistory` writes `BuildHistoryEventType_COMPLETE` with `Generation` tentatively at `G+1` inside `s.history.Update` (`sources/buildkit/solver/llbsolver/history.go:271-276`). The async trace finalizer then calls `UpdateRef` (`sources/buildkit/solver/llbsolver/history.go:318-327`) which reads `G+1`, bumps to `G+2`, and `b.Put`s. If an operator manually mutated the record between those two writes (e.g., `UpdateBuildHistory Pinned=true` at `sources/buildkit/control/control.go:360-368` via the same `UpdateRef` path), the trace update would overwrite that pin change at `G+2` without conflict detection. More critically, `Job.Discard` does not cancel in-flight `Update`/`UpdateRef` calls — they carry `context.WithoutCancel(ctx)` (`sources/buildkit/util/leaseutil/manager.go:51`, `sources/buildkit/solver/llbsolver/history/buildhistory.go:441`). A late `Solve` retry that recreates the same `Ref` after `Discard` could interleave `history.Update` writes for the new build with the old build's still-running trace finalizer, both operating on `ref_<Ref>` (`sources/buildkit/solver/llbsolver/history/buildhistory.go:279`).

A concrete thought experiment (pause-until-stolen):

> Pause a `Build` goroutine after it has acquired the scheduler's `execRes` slot but before `CacheManager.Save` (`sources/buildkit/solver/edge.go:998`). Let a concurrent solver or GC clean the intermediate content lease (5m blob lease at `sources/buildkit/solver/llbsolver/history/buildhistory.go:568` expires), or let the job be `Discard`ed (`sources/buildkit/solver/jobs.go:890`). Resume the old worker. It will proceed to `b.Put` / `Save` without re-reading any ownership token and will produce a `NewCachedResult` that is returned to whichever caller is still waiting on `gExecRes.Do`. If the original caller already timed out and a retry created a new `Job` with different `SessionID`, the late result still mutates the shared content store and cache records attached to the old context, potentially persisting blobs under an orphaned lease.

No test pins this: `grep stale` hits only cache-key staleness (`sources/buildkit/solver/cachemiss_stale_dep_test.go:9`) and `grep generation` hits only the proto display helper (`sources/buildkit/cmd/buildctl/debug/histories.go:91`).

## Architectural Decisions

* **Single-daemon, single-writer, collapse-duplicates:** `Solver` is an in-process DAG scheduler keyed by `Vertex.Digest()` (`sources/buildkit/solver/jobs.go:530-663`). All `Job`s multiplex onto shared `state` objects. This decision maximizes cache reuse and eliminates distributed fencing entirely — acceptable because the daemon is the sole mutator of its BoltDB and content store.

* **Leases as GC roots, not locks:** Choosing containerd's lease primitive (`sources/buildkit/vendor/github.com/containerd/containerd/v2/core/leases/lease.go:32-39`) repurposed for GC protection is lightweight and crash-safe (orphaned leases eventually expire), but the 1h/5m wall-clock expiries (`sources/buildkit/util/leaseutil/manager.go:32`, `sources/buildkit/solver/llbsolver/history/buildhistory.go:568`) mean long builds are never re-fenced; they rely on the daemon process staying alive to keep `WithLease` context alive, not on lease renewal.

* **No heartbeat / no reverse channel:** `session.Manager` (`sources/buildkit/session/manager.go:100-145`) uses gRPC stream completion (`<-c.ctx.Done()`) for liveness, with no periodic `Ping`. `history.Queue` (`sources/buildkit/solver/llbsolver/history/buildhistory.go:131-137`) uses a passive 120s GC ticker. This avoids complexity and clock-sync problems at the cost of no timely revocation.

* **Unconditional `Generation` bump:** Adding `Generation int32` (`sources/buildkit/api/services/control/control.proto:215`) without a CAS check suggests an intent to support optimistic concurrency later, but the increment at `sources/buildkit/solver/llbsolver/history/buildhistory.go:335` is currently decorative — it provides observability for clients polling `ListenBuildHistory` but no safety.

* **Async finalizer with `context.TODO()`:** The trace/provenance writer runs after the build's `context` is cancelled (`sources/buildkit/solver/llbsolver/history.go:286-332`) with `context.TODO()`/`WithTimeout(300s)`. This ensures traces are not lost on cancellation, but removes any association with the original lease/attempt and makes fencing impossible.

## Notable Patterns

* **Content-addressed dedup > fencing:** Every `Vertex` is canonicalized by `digest.FromBytes(...inputs + IgnoreCache flag)` (`sources/buildkit/solver/jobs.go:547-575`) and `state.LoadOrStore`-like insertion into `actives`. Duplicate subgraphs are physically merged (`edge.takeOwnership` at `sources/buildkit/solver/edge.go:124-128`); this is the inverse of fencing — intentionally reusing the first writer's result.

* **`flightcontrol.Group` as execution barrier:** `gExecRes.Do(ctx,"exec",…)` (`sources/buildkit/solver/jobs.go:1220`) plus `execDone`/`cacheDone` flags serialize `CacheMap`/`Exec` to at-most-one execution, sidestepping fencing by ensuring there is never a second concurrent executor for the same digest.

* **`ref_<Ref>` lease overlay for history:** `history.Queue.update` creates or reuses `ref_<Ref>` (`sources/buildkit/solver/llbsolver/history/buildhistory.go:429-435`) on every `BuildHistoryEventType_COMPLETE` (`sources/buildkit/solver/llbsolver/history.go:271`), attaching all `logs/trace/result` blobs as resources (`sources/buildkit/solver/llbsolver/history/buildhistory.go:445-483`). `clearOrphans` reconciles by inspecting `ListResources` length (`sources/buildkit/solver/llbsolver/history/buildhistory.go:214-217`) rather than generation.

* **Delayed `Discard` + `updateCond` wake:** `Job.Discard` retains `jobs[id]` for 10s (`sources/buildkit/solver/jobs.go:890-897`) to drain `Status` listeners, while `Get(id)` polls with `updateCond.Wait()` and 6s timeout (`sources/buildkit/solver/jobs.go:717-743`). This window absorbs stale `ListenBuildHistory` subscribers without any lease fencing.

* **Negative pattern — no stale-writer test:** `grep -r "heartbeat\|fenc\|CAS\|compareAndSwap"` in non-vendor yields zero fencing tests; `solver/cachemiss_stale_dep_test.go:9` tests cache staleness due to shared dep completion, not worker fencing.

## Tradeoffs

* **Reuse wins, retry quarantine loses:** Collapsing `flightcontrol.Group` and content hashing gives excellent cache hit rates and cross-build sharing, but sacrifices the ability to quarantine a suspect execution — there is no `(Ref,attempt)` key to force re-execution without mutating inputs (`IgnoreCache` is the only escape hatch at `sources/buildkit/solver/jobs.go:563-565`).

* **Wall-clock expiry vs active renewal:** Static `WithExpiration` labels (`sources/buildkit/util/leaseutil/manager.go:32`) avoid background heartbeat traffic and work without clock sync, but fail open if the daemon's wall clock skews or a build outlives the label. Containerd GC's `gcContext` scan interval (not shown in BuildKit, vendor-controlled) is the de-facto heartbeat, with indeterminate staleness.

* **Last-write-wins vs abort-late-writer:** `UpdateRef`'s `b.Put` (`sources/buildkit/solver/llbsolver/history/buildhistory.go:485`) never returns `Conflict`. Simpler and never blocks the build pipeline, but violates the dimension prompt's desired guarantee that a superseded worker cannot corrupt state.

* **In-memory ownership vs durable ownership:** `Solver.mu` (`sources/buildkit/solver/jobs.go:42`) is fast and crash-cleared on daemon restart (all `actives` are rebuilt from LLB), but durability of in-flight work is lost — there is no WAL to replay and no takeover path beyond `clearOrphans` janitor.

## Failure Modes / Edge Cases

* **Clock-skewed lease expiry hides blobs:** Worker A's 1h lease and 5m blob lease (`sources/buildkit/util/leaseutil/manager.go:32`, `sources/buildkit/solver/llbsolver/history/buildhistory.go:568`) are RFC3339 wall-clock labels (`sources/buildkit/vendor/github.com/containerd/containerd/v2/core/leases/lease.go:99`). If the daemon clock jumps forward, GC may delete the lease before `Writer.Commit` (`sources/buildkit/solver/llbsolver/history/buildhistory.go:621`) commits content; retry will re-create lease but content is lost, requiring `migrateBlobV2` recovery (`sources/buildkit/solver/llbsolver/history/buildhistory.go:287-303`) which synthesizes a new lease only if blob still exists.

* **Overwriting trace after completion:** As described, `recordBuildHistory`'s async `UpdateRef(trace)` (`sources/buildkit/solver/llbsolver/history.go:318-327`) uses `context.TODO()` and unconditional `Generation++`. If the same `Ref` is reused quickly after `Discard`, the trace finalizer of build N+0 can overwrite `Generation` of build N+1's `COMPLETE` record, or attach trace blobs to a lease that now belongs to a different logical build (same `ref_<Ref>` string).

* **No heartbeat → silent orphan on daemon kill -9:** In-flight `WithLease(ctx, lm, MakeTemporary)` contexts die with the process; temporary leases remain until `WithExpiration` label expires (1h default) or until operator Prune/GC runs. `clearOrphans` (`sources/buildkit/solver/llbsolver/history/buildhistory.go:200-239`) only prunes `BuildHistoryRecord`s with 0 resources, not orphaned cache/content leases, which linger until containerd's GC runs.

* **Session disconnect does not cancel exec:** `session.Manager.handleConn` deletes the session on `<-c.ctx.Done()` (`sources/buildkit/session/manager.go:135-145`), but `solver/jobs.go:1219-1236 Exec` context is the build's `ctx` from `Control.Solve`, not `session.Group`'s context. A client that drops its session leaves `gExecRes.Do` running; no heartbeat revokes it. The build completes or fails on its own timeout, not on session loss detection.

* **Concurrent `UpdateRef` loses updates:** Two goroutines calling `UpdateRef` simultaneously (e.g., `pinned→true` via API at `sources/buildkit/control/control.go:360` and trace commit) both `View` same `Generation=G`, both `Put` `G+1`, but with different fields. BoltDB transaction serializes writes, so the second overwrites the first's fields. No `ErrConflict`, no retry.

* **Reused `Ref` after `Discard` window:** After `Job.Discard`'s 10s grace (`sources/buildkit/solver/jobs.go:892`), a new `NewJob(sameRef)` succeeds. Listeners subscribed via `history.Queue.Listen` with `req.Ref==oldRef` still increment `refs[ref]` (`sources/buildkit/solver/llbsolver/history/buildhistory.go:783-794`). The new build's `h.update` call shares `ref_<Ref>` lease ID with the old build's still-pending async writer, causing resource attachment races.

## Future Considerations

* **If BuildKit remains single-daemon:** No fencing overhaul needed, but harden the current accidental fences: (a) turn `BuildHistoryRecord.Generation` into a real CAS by changing `UpdateRef` to `if br.Generation != expected { return ErrConflict }` and exposing `expectedGeneration` on `UpdateBuildHistoryRequest` (`sources/buildkit/api/services/control/control.proto:227`); (b) renew 5m blob leases for long trace captures (e.g., ticker every 2m inside `OpenBlobWriter`'s `Writer.Commit` path at `sources/buildkit/solver/llbsolver/history/buildhistory.go:618-635`); (c) add a stale-writer smoke test that sleeps an `Exec` goroutine for >`WithExpiration(0)` and asserts second waiter gets cached result via `gExecRes`. These mirror containerd's `WithExpiration(0)` trick in `cache/manager_test.go:676`.

* **If BuildKit grows a multi-daemon executor pool:** Introduce a first-class `AttemptID` on `SolveRequest` (alongside `Ref` at `sources/buildkit/api/services/control/control.proto:62`) and propagate it into `Job.uniqueID` + `BuildHistoryRecord.Generation` CAS, plus a lease renewal loop in `leaseutil.Manager` (heartbeat every `expiration/3`). Mutations to `b.Put`, `CacheManager.Save` (`sources/buildkit/solver/edge.go:998`), and `content.Writer.Commit` would then predicate on `(Ref, AttemptID, Generation)`. Add `heartbeat` RPC to `Control` service for active builds, with `ShardOwnershipLost`-style redirect on lease expiration (containerd's `caching_redirector` analogue).

* **Artifact fencing:** Lease `Adopt` (`sources/buildkit/util/leaseutil/manager.go:54-81`) should be transactional — copy-then-delete under a single BoltDB txn or use containerd's `SynchronousDelete` option (`sources/buildkit/vendor/github.com/containerd/containerd/v2/core/leases/lease.go:61-68`) to avoid the current window where resources are reachable from two leases.

* **Observability:** Emit `buildkit_build_generation` gauge and `buildkit_lease_expiry_seconds` histogram from `history.Queue.update` and `leaseutil.NewLease`; log `stamp mismatch` on rejected `UpdateRef` once CAS exists, analogous to Temporal's `StampNotFound`.

## Questions / Gaps

* `BuildHistoryRecord.Generation` was added for BoltDB migration (`sources/buildkit/solver/llbsolver/history/buildhistory.go:108-129` v2 migration) but never wired to a CAS check. Was the intention to add compare-and-swap and then abandoned, or is increment-only deliberately for UI ordering? History search shows no design doc linking `Generation` to fencing.

* What is the intended GC contract for `buildkit/lease.temporary` (`sources/buildkit/util/leaseutil/manager.go:87`) vs `containerd.io/gc.expire` (`sources/buildkit/vendor/github.com/containerd/containerd/v2/core/leases/lease.go:99`)? BuildKit labels both on the same lease; containerd GC's label priority and parsing (RFC3339 vs RFC3339Nano) not verified.

* `LeaseRef.Adopt` async `go Discard` (`sources/buildkit/util/leaseutil/manager.go:79`) after `AddResource` — if `AddResource` partially fails, resources are left on both leases. No test covers partial failure.

* No heartbeat constant or `DeamonHeartbeat` analogue found in BuildKit (zero vendor-excluded `heartbeat` hits). Is there a plan to reuse containerd's lease heartbeat service (`core/leases` `WithLabel` + expiry refresh) or is the daemon intentionally heartbeat-free?

* `sources/buildkit/solver/jobs.go:1219` `flightcontrol.Group` coalescing is the de-facto fencing for exec. Does this hold under memory pressure when `state.Release()` (`sources/buildkit/solver/jobs.go:315-330`) evicts `actives` entries while an exec is still blocked? `deleteIfUnreferenced` (`sources/buildkit/solver/jobs.go:745-783`) checks `len(st.jobs)==0 && len(st.parents)==0` — but a paused exec holds no extra ref beyond `st.jobs`, so eviction seems prevented, yet not explicitly tested with a pause-until-stolen scenario.

---

Generated by `dimensions/01.04-leases-heartbeats-fencing-stale-worker-rejection.md` against `buildkit`.
