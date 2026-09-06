# Source Analysis: buildkit

## Atomic State Transition and Event Journal

### Source Info

| Field | Value |
|-------|-------|
| Name | buildkit |
| Path | `studies/ultraplan-daemon-events-study/sources/buildkit` |
| Language / Stack | Go — bbolt (go.etcd.io/bbolt v1.5.0) + containerd leases/content + in-memory solver + gRPC control plane |
| Analyzed | 2026-09-02 |

## Summary

BuildKit explicitly separates **authoritative current-state** (in-memory `solver.Solver.jobs`/`actives` graph and `cache.metadata.Store` / `bboltcachestorage.Store` buckets) from an **operational observation journal** (`history.Queue` backed by `history.db` `_records` bucket + content blobs). It is not event sourcing: history records are written *after* the build completes as a sidecar for `buildctl --history` and UI replay, not as the source of truth for scheduling. Within each store, single-bucket mutations are correctly bounded in one `bolt.Tx` (`Transactor.Update`), and `history.Queue.update` bundles `b.Put` + `lease.Create` + `AddResource` for every descriptor in the same transaction. Cross-boundary atomicity is absent: solver job creation, progress streaming, blob ingest, and history publication each own their own transaction/connection. STARTED events are purely in-memory pubsub sends, COMPLETE events publish **after** the bolt commit, blob writers commit to the content store in a temporary lease before the history record references them, and the solver's live graph is never journaled. There is no optimistic revision check on `BuildHistoryRecord.Generation`, no WAL/checkpoint beyond bbolt's mmap/freelist, and `cache.db` is opened `NoSync:true` while `history.db` relies on `SafeOpen`'s corrupt-and-reset fallback that silently discards history.

## Rating

**3 / 10 — Absent / ad-hoc / unsafe**

The history queue demonstrates a correct *within-bucket* atomic pattern (`update` at `solver/llbsolver/history/buildhistory.go:418`), and cache storage operations are short bounded transactions. As a system-level atomic state+event journal, however, BuildKit is absent by design: authoritative solver state is RAM-only with no durable journal, the event log is an audit sidecar with lossy in-memory fan-out, state and event never share a transaction owner, revisions are not fenced, and crash recovery is best-effort (10 s job GC, 5 min temp leases, `fallbackOpen` DB reset). It trades durability for single-node throughput and simplicity, which is consistent with `flock`-guarded single-writer deployment but fails the dimension's crash-between-mutation-and-publication test.

## Evidence Collected

Every entry MUST include a file path with line numbers. Format: `path/to/file.ts:NN`.

| Area | Evidence | File:Line |
|------|----------|-----------|
| Transaction abstraction | `type Transactor interface { View; Update }` — thin wrapper over `*bolt.Tx` | `sources/buildkit/util/db/transactor.go:7-11` |
| Bolt open helper | `bolt.Open` forwarded without sync/WAL options | `sources/buildkit/util/db/boltutil/db.go:10-15` |
| Corrupt-and-reset open | `SafeOpen` → `fallbackOpen` renames corrupted DB to `*.bak` and creates empty DB; history silently lost | `sources/buildkit/util/db/boltutil/safe_open.go:17-48` |
| Daemon DB ownership | `historyDB` opened in `newController` via `SafeOpen(history.db, FreelistMapType)`; `cacheStorage` via `bboltcachestorage.NewStore(cache.db, NoSync:true)`; single lease manager from default worker `w.LeaseManager()` | `sources/buildkit/cmd/buildkitd/main.go:897-905`, `sources/buildkit/cmd/buildkitd/main.go:959-960` |
| Single-writer guard | `flock.New(lockPath).TryLock()` prevents multi-daemon writers; effective fencing is OS file lock | `sources/buildkit/cmd/buildkitd/main.go:386-397` |
| Authoritative solver state | `Solver.jobs map[string]*Job` + `actives map[digest]*state` protected only by `sync.RWMutex`; no bolt backing | `sources/buildkit/solver/jobs.go:42-44` |
| Cache storage buckets | `_result, _links, _byresult, _backlinks` created in one `db.Update`; `AddResult` writes `_links+_result+_byresult` atomically | `sources/buildkit/solver/bboltcachestorage/storage.go:16-25`, `sources/buildkit/solver/bboltcachestorage/storage.go:37-43`, `sources/buildkit/solver/bboltcachestorage/storage.go:144-172` |
| Cache NoSync | `NewStore` opens `cache.db` with `NoSync:true, FreelistMapType` — trades durability for speed | `sources/buildkit/solver/bboltcachestorage/storage.go:28-31` |
| Metadata store buckets | `_main, _index, _external` — `StorageItem.Queue` collects funcs, `Commit` applies all in one `s.db.Update` | `sources/buildkit/cache/metadata/metadata.go:20-23`, `sources/buildkit/cache/metadata/metadata.go:343-364` |
| History authoritative store | `recordsBucket="_records", versionBucket="_version"`; `Queue.active` is duplicate in-memory map mirroring `_records` | `sources/buildkit/solver/llbsolver/history/buildhistory.go:38-40`, `sources/buildkit/solver/llbsolver/history/buildhistory.go:57-58` |
| History schema/migration | `NewQueue` checks `versionBucket` version>1; otherwise calls `migrateV2()` isolating records to `ns+"_history"` | `sources/buildkit/solver/llbsolver/history/buildhistory.go:107-129`, `sources/buildkit/solver/llbsolver/history/migrate.go:18` |
| History bounded tx | `update()` does `DB.Update` → `Create/leaseID` → `addResource(Logs,Trace,ExternalError,Result...)` → `b.Put`; all in one tx; `AddResource` failures roll back `b.Put` via tx abort + lease cleanup defer | `sources/buildkit/solver/llbsolver/history/buildhistory.go:418-487` |
| History Generation | `UpdateRef` loads record via `View`, increments `br.Generation++`, then calls `update()` with `b.Put` — no expected-generation predicate, last-write-wins | `sources/buildkit/solver/llbsolver/history/buildhistory.go:309-347` |
| History update API | File `solver/llbsolver/history.go:29-56` STARTED via `history.Update(STARTED)` (memory-only); COMPLETE via `history.Update(COMPLETE)` after build (calls `update()` tx) | `sources/buildkit/solver/llbsolver/history.go:49-52`, `sources/buildkit/solver/llbsolver/history.go:271-273` |
| State+event split | `Queue.Update` for STARTED: `h.active[id]=record; ps.Send(e)` — no DB write. For COMPLETE: `ps.Send` **after** `h.update()` tx commits | `sources/buildkit/solver/llbsolver/history/buildhistory.go:538-550` |
| Pubsub (event journal) | In-memory `pubsub[T]` with buffered `chan(32)` and async `go c.send(v)`; `Send` never blocks tx, lossy on `close`, not durable | `sources/buildkit/solver/llbsolver/history/pubsub.go:5-28`, `sources/buildkit/solver/llbsolver/history/pubsub.go:49-54` |
| Blob vs record split | `OpenBlobWriter` creates 5-min temporary lease then `content.OpenWriter`; `Writer.Commit` commits blob; caller later calls `ImportStatus`/`update` in separate tx to attach descriptor to persistent lease | `sources/buildkit/solver/llbsolver/history/buildhistory.go:567-635`, `sources/buildkit/solver/llbsolver/history/buildhistory.go:679-768` |
| Trace late attach | Trace blobs committed after `BuildHistoryEventType_COMPLETE` via separate `UpdateRef` tx that does `b.Put` for `rec.Trace` | `sources/buildkit/solver/llbsolver/history.go:285-331` |
| Solver job lifecycle | `NewJob` inserts into `jobs` map under mutex; `Discard` sleeps 10 s before deleting, no bolt tx | `sources/buildkit/solver/jobs.go:687-715`, `sources/buildkit/solver/jobs.go:859-898` |
| Solver progress | `Master progress.MultiWriter` fan-out to `allPw`; jobs added via `connectProgressFromState` — in-memory only | `sources/buildkit/solver/jobs.go:59-62`, `sources/buildkit/solver/jobs.go:665-685` |
| Listen replay | `Listen` first replays `active` map, then `DB.View` `_records` bucket for COMPLETE events, then subscribes to live `pubsub` | `sources/buildkit/solver/llbsolver/history/buildhistory.go:770-887` |
| Delete + GC | `delete()` sends `DELETED` via pubsub then `DB.Update` deletes bucket + lease; `gc()` on 120 s tick deletes only if both `MaxEntries` and `MaxAge` exceeded | `sources/buildkit/solver/llbsolver/history/buildhistory.go:241-266`, `sources/buildkit/solver/llbsolver/history/buildhistory.go:152-198` |
| Orphan cleanup | `clearOrphans` on startup deletes `_records` entries with empty `ListResources(leaseID)` — compensates missing atomic blob+record | `sources/buildkit/solver/llbsolver/history/buildhistory.go:200-239` |

## Answers to Dimension Questions

**1. Can state commit without its corresponding durable event?**
Yes — by design, and in both directions. Solver authoritative state (`Solver.jobs`/`actives` at `solver/jobs.go:42`) commits in RAM with no durable event at all (`NewJob` at `solver/jobs.go:687`). History STARTED is the opposite: `Queue.Update` for `STARTED` at `solver/llbsolver/history/buildhistory.go:538-541` does `h.active[id]=rec; ps.Send(e)` with no bolt write, so the event exists only in memory and past-subscribers' replay cannot recover it after a crash — observers must poll `DB.View` which lacks active builds. History COMPLETE is the intended durable event: `history.go:271-273` calls `Queue.Update(COMPLETE)` which runs `update()`'s single `DB.Update` at `buildhistory.go:419` (`b.Put` + `AddResource` for logs/trace/result). That tx commits the state, but the live notification `ps.Send` at `buildhistory.go:548` happens after commit, not inside it — a crash between the two leaves a committed record with no live delivery (recovered only on next `Listen` DB scan).

**2. Can an event commit for a state transition that later rolls back?**
No durable event can outlive a rolled-back transition within a single store, but cross-store rollback does leave orphan events. Within `update()` (`buildhistory.go:418-487`) the `b.Put` and all `AddResource` calls share the same `bolt.Tx`; if any `AddResource` fails the tx aborts and the record is not visible, and the deferred lease cleanup at `buildhistory.go:439-443` removes the half-created persistent lease. Across stores, however, `Writer.Commit` at `buildhistory.go:618-635` commits the blob to `content.Store` under a temporary lease before the history record that references it is created. If `update()`'s tx later fails or the process crashes in between, the content blob remains durably committed with no referencing `_records` entry. `clearOrphans` at `buildhistory.go:214-216` and 5-min `MakeTemporary` expiry at `util/leaseutil/manager.go:83-88` mitigate this, but the observable window is `content.Info` exists without history record. The solver's progress events (`progress.MultiWriter` at `solver/jobs.go:59`) are never rolled back at all — they are fire-and-forget `pw.Write` calls at `solver/jobs.go:1369-1388`.

**3. Is the event log an audit/observation journal or the sole source of truth?**
Audit/observation journal. `BuildHistoryRecord` at `api/services/control/control.proto:203-225` carries `Ref, Frontend, CreatedAt/CompletedAt, Logs, Trace, Result, Error, Generation` as post-facto observation; the solver never reads history to decide scheduling — it reads `actives`/`jobs` and `CacheManager`/`CacheKeyStorage`. `Status` at `solver/llbsolver/solver.go:465-481` first tries `history.Status` (blob replay) and on `ErrNotExist` falls back to live `s.solver.Get(id)` — the live graph is the source of truth while the build is active. History's own `active` map at `buildhistory.go:57` is a non-durable mirror of in-flight builds that is rebuilt only via new `Update(STARTED)` calls, not from bucket contents. GC (`buildhistory.go:152`) and `clearOrphans` (`buildhistory.go:200`) are janitors over an audit log, not compaction of an event-sourced entity.

**4. Are transactions short and bounded?**
Inside each store, yes; end-to-end, no. `bboltcachestorage` ops (`AddResult` at `storage.go:144`, `AddLink` at `storage.go:314`, `Release` at `storage.go:200`) marshal small JSON and update at most 3 buckets within one `Update`. `cache.metadata.StorageItem.Commit` at `metadata.go:349-364` flushes a queued slice of small `Put`/`Index` ops in one `Update`. `history.update` at `buildhistory.go:418` loops over descriptors but each `AddResource` is a small lease-DB write; the bolt tx remains short. The unbounded work is explicitly pulled **outside** bolt txs: `OpenBlobWriter`/`ImportStatus`/`ImportError` (`buildhistory.go:567-768`) buffer up to 32 KiB frames and commit blobs under a separate content-store tx before the history tx; `recordBuildHistory` at `history.go:80-243` fans out provenance, status, and exporter descriptors via `errgroup` then calls `history.Update` once; trace attachment at `history.go:285-332` runs up to 3 s after `COMPLETE`. The solver's hot path (`jobs.go:1219` `gExecRes.Do`, `edge.go:124` `LoadOrStore`) holds only `sync.Mutex`, never a bolt tx, keeping writer contention bounded.

**5. How are stale state revisions rejected?**
They are not. `BuildHistoryRecord.Generation` at `control.proto:215` is incremented unconditionally in `UpdateRef` at `buildhistory.go:335` (`br.Generation++`) without an `expectedGeneration` argument or `WHERE generation=?` predicate; `h.update` at `buildhistory.go:485` does `b.Put` unconditionally. Concurrent `UpdateRef` callers serialize only via `h.mu` (`buildhistory.go:310-312`) in this single-daemon process — the last writer wins. Solver edges have no revision field; `sharedOp.gExecRes` flight control at `solver/jobs.go:992` and `state.jobs` fan-out at `solver/jobs.go:52` deduplicate by `Vertex.Digest()` content hash, not by version. Cache metadata `GetAndSetValue` at `metadata.go:432-444` does read-modify-write inside one `Update` but reloads `s.values[key]` from the in-memory `values` map populated at `newStorageItem` load time, not from a bolt `CompareAndSwap`; two concurrent `Update` transactions on the same `id` would both succeed, second overwriting first. The only rejection is identity collision on `Solver.NewJob` at `jobs.go:691` (`job ID exists` error).

## Architectural Decisions

* **RAM-first solver, bolt-second history.** `Solver` at `solver/jobs.go:41-46` keeps `jobs`+`actives` in sync maps for sub-millisecond scheduling; history is opt-out via `Internal` builds (`history.go:239`) and replayed only for `Status`/`Listen`. Decision trades single-node durability for cache-hit merging (`sharedOp` flight control) and avoids bolt writer contention on the scheduling hot path.
* **Single bolt tx per logical object.** `bboltcachestorage.AddResult` (`storage.go:144`) and `history.update` (`buildhistory.go:418`) each bundle all bucket writes + lease resources in one `DB.Update`. Lease and content references are co-committed with the record, preventing intra-object torn writes.
* **Content-then-record split with temporary lease.** `OpenBlobWriter` (`buildhistory.go:567`) creates a `MakeTemporary` lease (RFC3339 label at `leaseutil/manager.go:83`) with 5-min expiry, writes blob outside any history tx, then `update` re-parents the blob to the persistent `ref_<Ref>` lease. This releases the bolt writer during slow I/O while keeping blobs GC-safe.
* **Best-effort history with corrupt-and-reset.** `SafeOpen` at `boltutil/safe_open.go:17` treats a corrupted `history.db`/`cache.db` as disposable — backup + empty DB — rather than failing the daemon. History loss is acceptable; build scheduling continuity is not.
* **Pubsub as non-durable fan-out.** `pubsub[T]` at `pubsub.go:22` uses per-subscriber buffered channels and `go c.send` to isolate bolt txs from gRPC streams; `Listen` at `buildhistory.go:770` compensates by replaying `active` + `_records` before subscribing.
* **Monotonic `Generation` without fencing.** Increment-only counter at `buildhistory.go:335` supports UI ordering/counting (`prunehistories.go:114`) but not CAS — chosen because single-writer `flock` (`main.go:386`) makes optimistic concurrency unnecessary for the intended deployment model.

## Notable Patterns

* **Transactional outbox avoided.** Unlike a transactional outbox, BuildKit publishes history events *after* commit (`buildhistory.go:544-548`) rather than atomically appending to an outbox table polled by a relay. `Listen`'s DB replay (`buildhistory.go:823`) is the ad-hoc compensating mechanism.
* **Queue-and-commit for metadata.** `StorageItem.Queue(fn)` at `metadata.go:343` + `Commit()` at `metadata.go:349` batches idempotent index updates into one bolt tx — a local analog of optimistic batching without version checks.
* **Lease-based GC preconditions.** All history artifacts are referenced via `leaseID="ref_<Ref>"` at `buildhistory.go:280-281` in namespace `ns+"_history"` (`buildhistory.go:104-105`); `clearOrphans` (`buildhistory.go:214`) and `content GC` only delete unleased blobs, preserving the invariant "no GC while referenced."
* **Bounded freelist + map freelist.** Both `history.db` (`main.go:903`) and `cache.db` (`storage.go:28`) open with `FreelistMapType`; `cache.db` adds `NoSync:true` for write throughput — two-tier durability tuning.
* **Finalizer + graceful stop.** `AcquireFinalizer`/`Finalize` at `buildhistory.go:489-529` gates `ps.Close()` until active builds finish, preventing event loss on shutdown without making history synchronous.

## Tradeoffs

* **Throughput over durability.** `cache.db` `NoSync:true` (`storage.go:28`) avoids fsync on every cache link; a crash can lose recent cache graph edges, recomputed on next build. History uses synchronous bbolt (no `NoSync`) but accepts 10 s job GC and eventual `clearOrphans` — stale reads possible during the window after crash.
* **Simplicity over linearizability.** No `expectedGeneration`/`CAS` on `UpdateRef` (`buildhistory.go:335`) and no fencing in `Solver.load` avoids contention and version plumbing. Concurrent writers (only possible if `flock` is bypassed) produce last-write-wins with silent overwrite.
* **Isolation over latency.** Splitting blob ingest from history tx prevents holding the global bbolt write lock during I/O, but creates a two-phase window requiring compensating `clearOrphans` and temp-lease expiry — orphan GC is O(n) scan on startup and every 120 s.
* **Observability over consistency.** `Generation` and `NumCached/CompletedSteps` are populated after `eg.Wait()` in `history.go:196-200` — accurate for display but not usable as a consistency token.
* **Reset over repair.** `fallbackOpen` (`safe_open.go:38-44`) prefers daemon availability over history retention; operational cost is silent audit-log loss rather than manual recovery.

## Failure Modes / Edge Cases

* **Crash between blob Commit and history Put.** Leaves content blob + temp lease with no `_records` entry. Observable as `content.Info` exists but `history.Status(ref)` returns `ErrNotExist`; `Listen` never emits `COMPLETE`. Mitigated by temp-lease 5-min expiry and `clearOrphans` deleting leased-or-orphaned records on next start.
* **Crash between bolt Put and pubsub Send (COMPLETE).** Record is durable at `buildhistory.go:485` but live `Listen` subscribers miss the event (channel buffered 32, async send dropped on `done`). New subscribers recover via `DB.View` replay at `buildhistory.go:823-843`, but in-flight gRPC streams observe a gap.
* **Crash after STARTED Send before any durable write.** `active` map lost (`buildhistory.go:539` never persisted). Restarted daemon has no record of in-flight build; client `Solve` context canceled, build aborted with no history entry to query. `Solver.jobs` GC after 10 s (`jobs.go:890-896`) also drops the job client-side.
* **bbolt corruption → silent history wipe.** `SafeOpen` at `safe_open.go:38` renames corrupted `history.db` to `*.bak` and opens empty DB. All `COMPLETE`/`DELETED` events since last backup disappear; `Containerd` backend content blobs remain but lose their `ref_<Ref>` lease reference, becoming GC-eligible.
* **NoSync cache loss.** Power loss with `cache.db` `NoSync:true` can corrupt or lose recent `AddResult`/`AddLink` entries. Detectable as `ErrNotFound` on next `Load`/`WalkLinks`; scheduler falls back to re-execution rather than cache hit — no error, just performance regression.
* **Stale revision overwrite.** Two concurrent `UpdateRef` (e.g., `UpdateBuildHistory Pinned` at `control.go:360` + trace late attach at `history.go:318`) race on `br.Generation++`; second `b.Put` silently overwrites first's `Pinned` or `Trace` field. No `ErrConflict` surfaced to caller.
* **Writer starvation.** bbolt's single writer means a long `update()` that iterates many descriptors can block `CacheManager.AddResult` and vice versa. BuildKit bounds each tx to at most `Result.Results + Attestations + Logs/Trace` descriptors, but large provenance (in-toto with many layers) can still spike tx duration.
* **Backpressure drop.** `pubsub.Send` at `pubsub.go:23-26` launches a goroutine per subscriber and non-blockingly tries `ch <- v`; if subscriber is slow and buffer (32) fills, the message is dropped for that subscriber with no retry or persistence.

## Future Considerations

* **Transactional outbox table** for history: append `BuildHistoryEvent` to a `_outbox` bucket inside the same `DB.Update` that does `b.Put`, then have `Listen` tail the outbox (or a background relay) — makes crash-between-commit-and-publish impossible without changing the publish-after-commit pattern.
* **Add `expectedGeneration` to `UpdateRef`** — change signature to `UpdateRef(ctx, ref string, expectedGen int32, fn) error` and return `codes.Aborted` on mismatch, mirroring containerd/temporal conditional writes; surface via `control.UpdateBuildHistory` for UI pin/trace races.
* **Unify durability flags** — document and make `NoSync` configurable for `cache.db`; consider `NoFreelistSync:false` + `InitialMmapSize` tuning or switching history to WAL mode (bbolt `NoSync:false` is already WAL-like via `freelistMaps` but not exposed).
* **Make STARTED durable** as `PENDING` record — `Queue.Update(STARTED)` should `CreateBucket+Put` a pending record with `Generation=0` inside tx, then `COMPLETE` does `CAS` update; this gives `Listen` durable replay of in-flight builds across restarts and enables takeover by a new daemon after `flock` handover.
* **Bounded blob ingest** — replace `clearOrphans` O(n) scan (which calls `ListResources` per record at `buildhistory.go:214`) with a `temp_lease` bucket indexed by expiry, scanned incrementally; or make `Writer.Commit` directly attach to persistent lease inside the same `leases.Manager` tx if containerd lease manager supports nested tx.
* **Preserve `Generation` for GC decisions** — use `Generation` + `CompletedAt` as composite cursor for `gc()` pagination instead of in-memory sort at `buildhistory.go:181`, avoiding full table load when `MaxEntries` is large.

## Questions / Gaps

* `w.LeaseManager()` is obtained from the default worker at `main.go:959`; workers backed by `containerd` vs `runc` use different `LeaseManager` implementations — does the lease namespace `ns+"_history"` (`buildhistory.go:104`) share the same underlying bbolt file as the worker's content leases, and is `AddResource` inside `history.update`'s `DB.Update` re-entrant on the same `bolt.Tx` or a separate containerd `metadata` tx? No evidence found of nested tx handling — search of `util/leaseutil/manager.go` shows it delegates to `namespaces.WithNamespace` then `manager.Create/AddResource` without exposing tx ownership.
* No load or chaos tests were found proving `SafeOpen` fallback does not race with concurrent `flock` acquisition on restart; `bolt.Options.Timeout` is not set for `historyDB`/`cacheDB` opens — what is the blocking behavior if the previous daemon held the file lock?
* `Generation` field appears only in `prunehistories.go:114` display and `buildhistory.go:335` increment — no test asserts generation monotonicity or conflict detection; search of `*_test.go` returned no `Generation` coverage.
* WAL/checkpoint: bbolt `FreelistMapType` vs `FreelistArrayType` choice is not justified in comments; `NoSync:true` on `cache.db` vs synchronous `history.db` is unexplained in docs — confirm intended durability contract for cache vs history.
* Trace late-attach at `history.go:318` uses `context.TODO()` after the build context is canceled — if that `UpdateRef` fails, the trace blob's temp lease is leaked until `fallbackOpen` cleanup or manual GC; no metric/alert covers this path.

---
Generated by `dimensions/01.05-atomic-state-transition-and-event-journal.md` against `buildkit`.
