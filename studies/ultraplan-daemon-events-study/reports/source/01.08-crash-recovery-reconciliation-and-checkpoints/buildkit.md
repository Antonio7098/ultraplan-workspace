# Source Analysis: buildkit

## Dimension 01.08: Crash Recovery, Reconciliation, and Checkpoints

### Source Info

| Field | Value |
|-------|-------|
| Name | buildkit |
| Path | `studies/ultraplan-daemon-events-study/sources/buildkit` |
| Language / Stack | Go / buildkitd daemon + bbolt (cache.db, history.db, containerdmeta.db) + containerd leases/snapshotter/content store |
| Analyzed | 2026-09-03 |

## Summary

BuildKit has **no job-level crash recovery**: solver jobs, active vertices/edges, and in-flight build records are purely in-memory, so a daemon crash or host reboot abandons all non-terminal work with no startup reconciler, no interrupted-state machine, and no automatic new attempt. What *is* durable is **completed work only**, stored content-addressed in bbolt files under the daemon root: `cache.db` (solver cache keys/links/results), the cache metadata store (per-ref `committed`/`deleted`/`equalMutable`/snapshot IDs), and `history.db` (only `COMPLETE` build records; `STARTED` events live solely in an in-memory `active` map). Recovery therefore means **rebuild, not resume**: the client re-issues `Solve` and completed vertices are skipped via cache hits, while the interrupted vertex re-executes from scratch.

Around that model there is deliberate, tested crash-window hardening at the cache layer: `getRecord` repairs half-finalized refs (snapshot committed but `equalMutable` not cleared), finishes deleting records marked `deleted` before a crash, and drops mutable records whose snapshot is gone; prune marks `deleted` *before* removing data; `init` walks all metadata on startup and purges unloadable entries; `boltutil.SafeOpen` renames a corrupt DB to `.bak` and starts empty; worker startup sweeps temporary leases; and the history queue deletes completed records whose blobs are missing (`clearOrphans`). Exec cancellation is prompt (ctx-cancel → SIGKILL), but no daemon-restart sweep of orphaned runc containers was found — that is delegated to containerd/runc. A SIGKILL during an external side effect (e.g. registry push in an exporter finalizer) has no BuildKit-level idempotency guard or exactly-once record; retry safety comes only from content-addressed registry semantics, and a killed push simply re-runs on the next manual `Solve`.

## Rating

**6/10** — Present but inconsistent/fragile. Completed-work durability is explicit, well-schematized, and covered by a dedicated crash-simulation test, plus operational safeguards (deleted-first prune, temp-lease sweep, corrupt-DB fallback, orphan-record purge). But the dimension's core — startup reconciliation of non-terminal work, resumable attempts, takeover rules — is absent by design: interrupted builds vanish silently, resume granularity is whole completed vertices only, retry is always manual, and durability is weakened by `NoSync: true` plus fire-and-forget async cleanup. No evidence found of crash/restart integration tests beyond the single cache-manager unit test.

## Evidence Collected

Every entry MUST include a file path with line numbers. Format: `path/to/file.ts:NN`.

| Area | Evidence | File:Line |
|------|----------|-----------|
| Solver jobs are in-memory only; no persisted attempt state | `Solver` holds `jobs map[string]*Job` and `actives map[digest.Digest]*state`, all in-process; `NewJob` just allocates structs and a progress context | `sources/buildkit/solver/jobs.go:43-44`, `sources/buildkit/solver/jobs.go:687-715` |
| Job discard is local cleanup, not a durable state transition | `Discard` removes the job from in-memory `actives`, runs `releasers`, and deletes the job entry after a 10s delay — nothing is checkpointed | `sources/buildkit/solver/jobs.go:859-898` |
| No retry/backoff in scheduler or edges | `edge.markFailed` terminally records the error for that attempt; repo-wide search finds no attempt counter, backoff, or auto-retry in `solver/scheduler.go`, `solver/edge.go`, `solver/jobs.go` | `sources/buildkit/solver/edge.go:374-379` |
| Durable solver-result checkpoints (completed cache keys) | `bboltcachestorage.NewStore` opens `cache.db` with `NoSync: true` and buckets `_result`, `_links`, `_byresult`, `_backlinks`; `AddResult`/`AddLink` persist completed step results keyed by cache ID | `sources/buildkit/solver/bboltcachestorage/storage.go:27-48`, `sources/buildkit/solver/bboltcachestorage/storage.go:144-172` |
| Daemon wires durable stores from disk root | `main.go` opens `cache.db` via `bboltcachestorage.NewStore` and `history.db` via `boltutil.SafeOpen`, both under `cfg.Root`, and injects them into the controller | `sources/buildkit/cmd/buildkitd/main.go:897-908`, `sources/buildkit/cmd/buildkitd/main.go:947-966` |
| Per-ref checkpoint schema (committed/deleted/equalMutable/snapshot IDs) | Metadata keys `snapshot.committed`, `cache.deleted`, `cache.equalMutable`, `cache.snapshot`, `cache.blob`, `cache.chainID`, `cache.blobChainID` are the durable checkpoint fields | `sources/buildkit/cache/metadata.go:15-40` |
| Crash-window repair on record load | `getRecord`: missing equal-mutable + existing snapshot → `clearEqualMutable` and reuse ("there may have been a crash while finalizing"); `deleted` flag → finish `remove`; mutable with missing snapshot → purge record | `sources/buildkit/cache/manager.go:440-506` |
| Startup walk purges unloadable snapshots | `cacheManager.init` iterates `MetadataStore.All()`, and for any record `getRecord` cannot load, logs at debug, `Clear`s metadata and deletes the `ID` and `ID-variants` leases | `sources/buildkit/cache/manager.go:326-341` |
| Deleted-first prune ordering for crash safety | Prune path sets `dead = true`, `queueDeleted()` + `commitMetadata()` *before* releasing locks and removing data, so a crash mid-prune resumes deletion on next load | `sources/buildkit/cache/manager.go:1216-1233` |
| Finalize ordering (lease → snapshot commit → clear equalMutable → commit metadata) | `finalize` creates lease, adds snapshot resource, `Snapshotter.Commit`, then clears `equalMutable` and commits metadata; the window between snapshot commit and metadata commit is exactly what `getRecord` repairs | `sources/buildkit/cache/refs.go:1498-1544` |
| Crash-simulation test for half-finalized ref | `TestLoadHalfFinalizedRef` commits the snapshot without clearing metadata, reopens the manager, and asserts the immutable ref remains usable while the stale mutable is gone | `sources/buildkit/cache/manager_test.go:2309-2389` |
| Corrupt-DB fallback (cache + history) | `SafeOpen` recovers from open panics/errors by renaming the DB to `<path>.<id>.bak` and opening fresh; log message explicitly attributes corruption to buildkitd crash/SIGKILL; data loss accepted for disposable DBs | `sources/buildkit/util/db/boltutil/safe_open.go:17-49` |
| History persists only COMPLETE records | `Queue.Update` keeps `STARTED` events in the in-memory `active` map and only writes to the `_records` bucket on `COMPLETE`; `Status`/`Listen` replay comes from bolt, actives from memory | `sources/buildkit/solver/llbsolver/history/buildhistory.go:531-551` |
| Orphan history-record purge on startup | `clearOrphans` lists records whose history-namespace lease has no content resources ("missing blobs") and deletes them; launched in a background goroutine at `NewQueue` alongside periodic `gc` | `sources/buildkit/solver/llbsolver/history/buildhistory.go:131-137`, `sources/buildkit/solver/llbsolver/history/buildhistory.go:200-239` |
| Temp-lease sweep on worker startup (orphan content/snapshots) | `NewWorker` lists leases labelled `buildkit/lease.temporary` and deletes them, releasing content/snapshots held only by crashed builds; `MakeTemporary` marks such leases | `sources/buildkit/worker/base/worker.go:222-228`, `sources/buildkit/util/leaseutil/manager.go:83-90` |
| Sessions do not survive daemon restart | `getSalt` generates a fresh random salt per process ("unique component per daemon restart to avoid persistent keys"), invalidating prior token-authority signatures | `sources/buildkit/session/auth/auth.go:19-29` |
| Exec cancellation is SIGKILL, in-process only | `runcProcessHandle` documents "on ctx.Done the in-container process will receive a SIGKILL"; `procKiller.Kill` sends SIGKILL; Linux sets `PdeathSignal = SIGKILL`; no restart-time container sweep found in BuildKit code | `sources/buildkit/executor/runcexecutor/executor.go:647-653`, `sources/buildkit/executor/runcexecutor/executor.go:577-630`, `sources/buildkit/executor/runcexecutor/executor_linux.go:26` |
| Export phase uses one temp lease, no idempotency record | `llbsolver.Solve` creates a single temporary lease spanning exporters/cache-export/provenance/finalize so blobs are not GC'd mid-push; no exactly-once token or partial-push journal found | `sources/buildkit/solver/llbsolver/solver.go:371-387` |
| Failed builds recorded only if the daemon lives to finish them | History completion handler imports the error into a blob (`ImportError`), sets `ExternalError`/`Error` on the record, then persists a `COMPLETE` event — a SIGKILL before this point leaves no record at all | `sources/buildkit/solver/llbsolver/history.go:251-261`, `sources/buildkit/solver/llbsolver/history.go:271-278` |

## Answers to Dimension Questions

1. **What survives a client crash, daemon crash, and host reboot?**
   - *Client crash*: daemon-side state persists — completed cache entries in `cache.db` (`sources/buildkit/solver/bboltcachestorage/storage.go:144-172`), committed ref metadata (`sources/buildkit/cache/metadata.go:15-40`), completed history records in `history.db` (`sources/buildkit/solver/llbsolver/history/buildhistory.go:531-551`). The in-flight build's gRPC/session context is canceled, so its running execs are SIGKILLed (`sources/buildkit/executor/runcexecutor/executor.go:647-653`) and its `STARTED` history entry (memory-only `active` map) is dropped.
   - *Daemon crash*: same durable set as above, since all three bolt files live under `cfg.Root` (`sources/buildkit/cmd/buildkitd/main.go:897-908`). All solver jobs/actives (`sources/buildkit/solver/jobs.go:43-44`), progress pipes, and `active` history entries are lost. On next start, `cacheManager.init` purges unloadable refs (`sources/buildkit/cache/manager.go:326-341`), temp leases are swept (`sources/buildkit/worker/base/worker.go:222-228`), and corrupt bolt files are renamed to `.bak` with fresh empty DBs (`sources/buildkit/util/db/boltutil/safe_open.go:37-49`) — i.e. crash recovery can silently discard history/cache.
   - *Host reboot*: identical to daemon crash for BuildKit's own state (files on disk), except running containers/execs are gone at the OS level and temp-lease sweep plus containerd snapshotter/content-store GC reclaim their resources. No evidence found of a BuildKit-level reboot reconciler beyond the startup paths above.
2. **Can arbitrary in-memory execution resume, or only checkpointed work?** Only checkpointed (completed) work. There is no execution checkpoint, journal, or snapshot of running vertices — resume granularity is "whole completed vertex via content-addressed cache hit on re-`Solve`". A vertex interrupted mid-exec re-runs from scratch; mid-vertex progress (partial snapshots, pipes, scheduler edge state) is not persisted anywhere.
3. **When is a new attempt started automatically?** Never, at the daemon level. No automatic retry, backoff, or takeover rule was found: `markFailed` is terminal for the attempt (`sources/buildkit/solver/edge.go:374-379`), and interrupted builds are not re-queued on startup. The only "automatic" behavior is cache-hit skipping of already-completed steps when the *client* starts a new `Solve`.
4. **When is human/manual retry required?** Always, for anything not already durably cached: any build that fails, is canceled by client disconnect, or is interrupted by daemon crash/reboot must be re-submitted by the client. The completed history record (with `Error`/`ExternalError`, `sources/buildkit/solver/llbsolver/history.go:251-261`) is the observability artifact the operator uses to decide to retry.
5. **How are orphaned resources discovered?** Four mechanisms, all heuristic/sweep-based, none transactional: (a) `cacheManager.init` load-and-purge walk (`sources/buildkit/cache/manager.go:326-341`); (b) lazy per-record repair in `getRecord` via `committed`/`deleted`/`equalMutable`/snapshot-existence checks (`sources/buildkit/cache/manager.go:440-506`); (c) worker-startup deletion of `buildkit/lease.temporary` leases (`sources/buildkit/worker/base/worker.go:222-228`); (d) history `clearOrphans` deleting `COMPLETE` records whose lease holds no blobs (`sources/buildkit/solver/llbsolver/history/buildhistory.go:200-239`). Orphaned *running* runc containers after a daemon SIGKILL have no BuildKit-side discovery found — containerd/runc reaping is the implicit backstop (gap).

> **SIGKILL the execution owner during an external side effect. On restart, how does the system decide whether retry is safe?** It doesn't decide — there is no safety check. Exporter finalizers (e.g. registry push) run after solve under a single temp lease (`sources/buildkit/solver/llbsolver/solver.go:371-387`) with no idempotency token, no side-effect journal, and no partial-push record. A SIGKILL mid-push leaves no trace in `history.db` (the `COMPLETE` record is only written if the daemon survives to the completion handler, `sources/buildkit/solver/llbsolver/history.go:271-278`). Retry safety is entirely emergent: pushes are content-addressed by digest, so re-running the build re-pushes identical blobs idempotently at the registry layer. Anything non-idempotent at the sink (e.g. a mutating tag push that succeeded before the kill) could be duplicated or left half-applied with no BuildKit-level detection — no evidence found of fencing, deduplication keys, or exactly-once export semantics.

## Architectural Decisions

- **Rebuild-from-cache instead of resume-from-journal.** The core decision is that the content-addressed result cache *is* the checkpoint mechanism. This trades fine-grained recovery for simplicity: no WAL, no attempt log, no fencing — crash consistency reduces to making each cache commit atomic-ish and repairing torn commits on load (`sources/buildkit/cache/manager.go:440-506`, `sources/buildkit/cache/refs.go:1498-1544`).
- **Deleted-first prune ordering.** Marking `cache.deleted` in metadata before touching data (`sources/buildkit/cache/manager.go:1219-1228`) makes GC crash-safe at the cost of a metadata write per prune batch.
- **Disposable-DB philosophy for cache/history.** `SafeOpen` prefers availability over durability: a corrupt `cache.db`/`history.db` is quarantined to `.bak` and replaced with an empty DB (`sources/buildkit/util/db/boltutil/safe_open.go:37-49`). Combined with `NoSync: true` on `cache.db` (`sources/buildkit/solver/bboltcachestorage/storage.go:27-31`), this explicitly trades committed-but-unflushed results for startup reliability and write throughput.
- **Leases as the GC root, swept on start.** Temporary-lease sweeping (`sources/buildkit/worker/base/worker.go:222-228`) plus lease-backed history blobs (`sources/buildkit/solver/llbsolver/history/buildhistory.go:429-487`) unify orphan detection around the containerd lease manager rather than a BuildKit-owned reconciler.
- **History as observability, not recovery log.** Only terminal records are persisted; `STARTED` is ephemeral (`sources/buildkit/solver/llbsolver/history/buildhistory.go:538-551`). This keeps the history DB small and simple but forfeits any "interrupted builds" view after restart.

## Notable Patterns

- **Torn-commit repair idiom**: `Snapshotter.Commit` → (crash window) → `clearEqualMutable` + `commitMetadata`; the reader (`getRecord`) detects and heals the torn state. Same shape as write-ahead-mark patterns in log-structured systems, but implemented ad hoc per record type.
- **Startup-sweep pattern**: `init` walk, temp-lease list-and-delete, and `clearOrphans` are all "list everything, delete what looks dead" passes rather than event-driven reconciliation.
- **Single temp lease spanning a multi-phase commit** (solve → export → cache-export → provenance → finalize, `sources/buildkit/solver/llbsolver/solver.go:371-387`) to hold GC at bay across phases that each assume lease protection individually.
- **Crash-attributing log message as operational interface**: the `SafeOpen` fallback log explicitly names crash/SIGKILL as the likely cause and points at the issue tracker (`sources/buildkit/util/db/boltutil/safe_open.go:39-41`).

## Tradeoffs

- `NoSync: true` on `cache.db` buys write throughput but a host crash can lose recently committed results with no detection (subsequent builds just recompute — safe but wasteful).
- `.bak`-and-restart on corruption buys availability but silently discards history and cache indexes; an operator investigating a crash may find the evidence quarantined beside an empty DB.
- Memory-only `active` history map keeps the hot path lock-simple but means interrupted builds are invisible post-restart — no "was running, now unknown" state for operators or automation to key on.
- Whole-vertex re-execution keeps the executor stateless across restarts but makes long single-RUN builds pay full redo cost after any late crash; no intra-op checkpointing (e.g. CRIU paths exist in vendored go-runc but are not wired into the exec path — no non-vendor references found).
- Async fire-and-forget `clearOrphans`/GC goroutine (`sources/buildkit/solver/llbsolver/history/buildhistory.go:131-137`) avoids slowing startup but races with early reads and swallows errors (returns `error` that the goroutine ignores).

## Failure Modes / Edge Cases

- **SIGKILL between snapshot commit and metadata commit**: handled — `getRecord` clears `equalMutable` and reuses the ref (`sources/buildkit/cache/manager.go:455-464`); covered by `TestLoadHalfFinalizedRef` (`sources/buildkit/cache/manager_test.go:2309-2389`).
- **Crash mid-prune**: handled — `deleted` marker causes completion of removal on next access (`sources/buildkit/cache/manager.go:486-492`).
- **Corrupt bolt file after SIGKILL mid-write**: "handled" by amputation — DB reset to empty with `.bak` quarantine (`sources/buildkit/util/db/boltutil/safe_open.go:37-49`); all cache indexes or all history lost.
- **Crash leaving mutable with no snapshot**: handled — record purged (`sources/buildkit/cache/manager.go:494-506`).
- **Orphaned running containers after daemon SIGKILL**: no BuildKit-side handling found; reliance on containerd/runc. Stale container IDs could collide or leak until external GC.
- **SIGKILL during registry push / cache export**: no handling found — partial external side effects are invisible to the next attempt; safety depends entirely on sink idempotency.
- **History blob missing but record present** (e.g. content GC raced the record): `clearOrphans` deletes the whole record with only a warn log (`sources/buildkit/solver/llbsolver/history/buildhistory.go:232-236`) — observability data loss without operator recourse.
- **`init` purge logging at debug only** (`sources/buildkit/cache/manager.go:334`): post-crash ref loss is nearly silent in default logs.

## Future Considerations

- Persist `STARTED` history records (or a minimal interrupted-attempt journal) so post-restart tooling can distinguish "never completed" from "never started" and drive manual or policy-based retry.
- Add a startup reconciliation pass that enumerates `actives`-equivalent durable state (leases, uncommitted mutables, running container IDs) and either resumes, fences, or explicitly fails them with a recorded terminal event.
- Make `clearOrphans`/`init` purges synchronous-before-serve (or at least error-surfaced with metrics) rather than fire-and-forget, and promote purge logs to warn with counts.
- Consider `NoSync: false` (or fdatasync on commit paths) for `history.db`, which is low-write-volume but high-investigative-value, while keeping `NoSync` for `cache.db`.
- Introduce idempotency keys or a two-phase export record (intent → finalize) for registry/image pushes so a SIGKILL mid-push is detectable and safely retryable at the BuildKit layer.
- Add a crash/restart integration test (kill daemon mid-build, restart, assert cache reuse + no leaked leases/snapshots + history state) to sit alongside the single unit-level `TestLoadHalfFinalizedRef`.

## Questions / Gaps

- No evidence found for how (or whether) BuildKit reaps runc containers that outlive a daemon SIGKILL — searched `executor/runcexecutor`, `worker/runc`, `worker/base` for delete/sweep/monitor paths; only in-process ctx-cancel killing found. Needs runtime verification against a live daemon.
- No evidence found of snapshotter-level recovery (e.g. overlayfs upperdir orphaned by a killed exec) beyond lease-based GC — containerd snapshotter internals are out of scope of this source tree.
- No evidence found for client-side (buildctl) resume/retry semantics after daemon restart — whether buildctl retries `Solve` automatically or surfaces the disconnect directly was not traced (client package not examined in depth).
- `clearOrphans` uses `proto.Unmarshal` while writers use `MarshalVT` (`sources/buildkit/solver/llbsolver/history/buildhistory.go:209-213` vs `:424`); wire compatibility is assumed, not verified here.
- The `GracefulStop` path only drains history finalizers/pubsub (`sources/buildkit/solver/llbsolver/history/buildhistory.go:139-147`); no graceful-drain of running solver jobs was found — SIGTERM behavior of in-flight builds is unconfirmed.

---

Generated by `Dimension 01.08: Crash Recovery, Reconciliation, and Checkpoints` against `buildkit`.
