# Source Analysis: buildkit

## Durable Submission and Idempotent Acceptance

### Source Info

| Field | Value |
|-------|-------|
| Name | buildkit |
| Path | `studies/ultraplan-daemon-events-study/sources/buildkit` |
| Language / Stack | Go / gRPC (control.proto), bbolt, containerd leases/content, solver + history queue |
| Analyzed | 2026-09-02 |

## Summary

BuildKit's build submission is a synchronous unary gRPC `Control.Solve` that treats the client-provided `Ref` as a transient in-memory job identifier, not a durable idempotency key. Acceptance is a `Solver.NewJob` map insert under a mutex (`solver/jobs.go:687`), followed by a best-effort `history.Queue.Update(STARTED)` BoltDB write (`solver/llbsolver/history.go:49`) — the two are not atomic and no transaction couples them to the gRPC acknowledgment. The `Solve` RPC blocks until the full build (frontend solve, exporter, cache export) completes, so there is no prepare/start split and no durable queue ahead of execution. Duplicate `Ref` while a job is live is rejected with `job ID %s exists`; after the 10 s post-`Discard` grace the same `Ref` is silently accepted as a new build with no fingerprint or conflict check. Vertex-level execution is heavily deduplicated via `flightcontrol`/`sharedOp` and content-addressed cache, but job-level submission has no idempotency table, no fingerprint, no retention policy for deduplication state beyond the ephemeral `jobs` map, and no client retry/resume protocol. History records provide durable observability after the fact (`MaxAge 48h / MaxEntries 50`) but do not gate acknowledgement and do not enable retry to return the original operation.

## Rating

**2 / 10 — Absent, implicit, ad-hoc, unsafe for durable/idempotent acceptance**

Rationale: No durable acceptance before acknowledgement; acknowledgement (in-memory `NewJob` success) precedes durable history commit, with no WAL or transaction. No idempotency key semantics, fingerprinting, or deduplication store. Duplicate submission handling is a bare existence check with no identical-vs-conflicting-input distinction and no retry-to-resume. Retention of submission state is a 10 s in-memory grace. The mature vertex-level cache/merge/flightcontrol machinery does not extend to the submission boundary.

## Evidence Collected

Every entry MUST include a file path with line numbers. Format: `path/to/file.ts:NN`.

| Area | Evidence | File:Line |
|------|----------|-----------|
| Submit/start handler (gRPC entry) | `func (c *Controller) Solve(ctx context.Context, req *SolveRequest) (*SolveResponse, error)` — validates, resolves exporters/cache, delegates to `c.solver.Solve(ctx, req.Ref, ...)` and blocks until build completes | `control/control.go:409-592` |
| Solve proto (no idempotency fields) | `message SolveRequest { string Ref = 1; pb.Definition Definition = 2; string Session = 5; ... }` and `service Control { rpc Solve(SolveRequest) returns (SolveResponse); }` — unary, no idempotency token, no fingerprint, no dedup metadata | `api/services/control/control.proto:14-18,61-82` |
| Acceptance = in-memory map insert | `func (jl *Solver) NewJob(id string) (*Job, error)` checks `if _, ok := jl.jobs[id]; ok { return nil, errors.Errorf("job ID %s exists", id) }` under `jl.mu.Lock`, creates `progress` ctx and inserts to `jl.jobs` — no DB, no hash | `solver/jobs.go:687-715` |
| No transaction coupling acceptance to durability | `recordBuildHistory` called *after* `NewJob`, does `s.history.Update(ctx, {Type: STARTED, Record: rec})` in separate BoltDB `Update`; failure returns `Internal` but `NewJob` already succeeded; `Solve` `defer j.Discard()` still runs | `solver/llbsolver/solver.go:210-292` and `solver/llbsolver/history.go:49-57` |
| Durable history write (STARTED) | `func (h *Queue) Update(...){ h.mu.Lock(); if e.Type==STARTED { h.active[e.Record.Ref]=e.Record; h.ps.Send(e)}; if COMPLETE { delete(h.active,...); h.update(ctx,e.Record) } }` and `h.update` does `DB.Update` + `leaseManager.Create` + `addResource` + `b.Put` | `solver/llbsolver/history/buildhistory.go:531-551,418-487` |
| Current-only idempotency store is RAM | `type Solver struct { mu sync.RWMutex; jobs map[string]*Job; actives map[digest.Digest]*state }` — no bbolt bucket for jobs; `actives` content-addressed per vertex digest, not per `Ref` | `solver/jobs.go:41-46` |
| Discard / retention (submission dedup window) | `func (j *Job) Discard() error { ...; go func(){ time.Sleep(10*time.Second); delete(j.list.jobs, j.id) }() }` — 10 s grace keeps `Ref` reserved; after that any `Ref` reusable | `solver/jobs.go:859-898` |
| History retention (observability, not dedup) | `Default MaxAge 48h, MaxEntries 50` `go func(){ for { h.gc(); time.Sleep(120*time.Second) } }()`; GC requires `len(records) >= MaxEntries` AND `now-MaxAge > CompletedAt`; pinned never deleted | `solver/llbsolver/history/buildhistory.go:81-86,131-198` |
| Duplicate Ref handling (same existence check) | Identical error for any duplicate active `Ref`, no payload comparison: `return nil, errors.Errorf("job ID %s exists", id)`; llbsolver bubbles as `return nil, err` from `s.solver.NewJob(id)` | `solver/jobs.go:692` and `solver/llbsolver/solver.go:210-213` |
| Cache export duplicate rejection (not submission idempotency) | `findDuplicateCacheOptions` / `duplicate cache exports` — dedupes `CacheOptionsEntry` by `ref` or hash, not `SolveRequest` idempotency | `control/control.go:465-474,778-806` |
| Client Ref generation (no stable idempotency key) | `ref := identity.NewID(); if opt.Ref != "" { ref = opt.Ref }` per `Client.solve`; `identity.NewID()` is random; `session.NewSession` also random | `client/solve.go:107-110,126` |
| Client submit is single shot, no retry | `resp, err := c.ControlClient().Solve(ctx, sopt)` once in `errgroup`; error wrapped as `failed to solve`; no retry loop, no idempotency header | `client/solve.go:337-343,399` |
| Validation around acceptance boundary | `translateLegacySolveRequest`, `compat.ValidateCompatibilityVersion`, `entitlements.WhiteList`, `validateSourcePolicy`, exporter/session resolve — no request fingerprint, no `Definition` digest check against existing `Ref` | `control/control.go:412-446,452-482` and `solver/llbsolver/solver.go:243-264` |
| Status streaming: history vs live | `func (s *Solver) Status(ctx, id, ch)` tries `s.history.Status` (BoltDB+content blob) first; on `os.ErrNotExist` falls back to `s.solver.Get(id)` live progress | `solver/llbsolver/solver.go:465-481` |
| Vertex-level dedup (contrasts with job-level absence) | `sharedOp.gCacheRes / gDigest / gExecRes flightcontrol.Group` dedupes `CacheMap` and `Exec` per vertex digest; `jobs` still per-vertex shared `state.jobs` map | `solver/jobs.go:995-998,1137-1208,1219-1280` and `solver/jobs.go:52-53` |
| History finalizer / trace commit | `AcquireFinalizer` / `Finalize` + async `UpdateRef` for trace blob after `COMPLETE` — proves history commit is deferred, not on acceptance | `solver/llbsolver/history/buildhistory.go:489-530` and `solver/llbsolver/history.go:263-285,299-332` |

## Answers to Dimension Questions

### 1. Can acknowledgement happen before durable acceptance?

**Yes — acknowledgement (in-memory job creation success) always precedes durable acceptance, and they are not atomic.**

Flow: `Control.Solve` → `solver.NewJob(id)` (`solver/jobs.go:687`) acquires `jl.mu`, checks existence, inserts into `jl.jobs`, broadcasts — this is the acknowledgement gate; the gRPC handler only returns error if this fails. Durable state is the `BuildHistoryRecord` written later by `llbsolver.recordBuildHistory` (`solver/llbsolver/history.go:49`) which does `h.update` → BoltDB `DB.Update` + lease/content writes (`solver/llbsolver/history/buildhistory.go:418-485`). The two are sequential, not in a single transaction, and `recordBuildHistory` is invoked *after* `NewJob` (`solver/llbsolver/solver.go:210-285`). A crash after `NewJob` returns but before `h.update` commits leaves an in-memory active vertex graph with no persisted `STARTED` record; `clearOrphans` on restart will delete orphan BoltDB records with missing blobs, not recover the job (`solver/llbsolver/history/buildhistory.go:200-239`). Further, `Solve` is unary and blocks until export completes — there is no intermediate ack that client can use to poll; the only durable artifact is the `COMPLETE` record written via `s.history.Update` in the finalizer (`solver/llbsolver/history.go:271-278`), executed with `context.WithoutCancel`.

### 2. What happens if the same key is retried with identical input?

**Hard failure if job still tracked; silent new build if grace expired — no idempotent return of original operation.**

Active window: `NewJob` rejects any duplicate `Ref` present in `jl.jobs` regardless of payload: `errors.Errorf("job ID %s exists", id)` (`solver/jobs.go:692`). `llbsolver.Solve` bubbles this directly (`solver/llbsolver/solver.go:210-213`), which `Control.Solve` returns as gRPC error; `client/solve.go:337-339` wraps as `failed to solve`. There is no check that `Definition`/`FrontendAttrs`/`Exporters` match, and no path that returns the original `SolveResponse` or joins the existing execution. `Solver.Get` (`solver/jobs.go:717-743`) waits up to 6 s for a job ID to appear for `Status` but is not used for `Solve` deduplication.

Post-completion: `j.Discard()` removes job-vertex links immediately but retains `jl.jobs[id]` for 10 s (`solver/jobs.go:890-895`). Within that window retry still gets `exists`. After deletion, the same `Ref` with identical `Definition` is accepted as a brand-new job; solver will recreate `state` entries keyed by vertex digest (`solver/jobs.go:545-594`) which may hit content-addressed cache and thus run faster, but it is not recognized as a retry — new `uniqueID`, new `ch` for `ImportStatus`, new history `STARTED` event.

### 3. What happens if the same key is reused with different input?

**Identical behavior to Q2 — no conflict detection, no input comparison.**

The existence check is the only gate (`solver/jobs.go:692`). No canonical request fingerprint is computed; `Definition.Def` (LLB DAG), `Frontend`, `FrontendAttrs`, `FrontendInputs`, `CacheImports/Exports`, `Entitlements`, `SourcePolicy`, `Exporters` are never hashed or compared to the original job's values stored in `BuildHistoryRecord` (`solver/llbsolver/history.go:35-47`) or `state.vtx`. If the prior job is still active (including 10 s grace), different input is also rejected with the same `job ID exists` error — the caller cannot distinguish conflict-vs-duplicate via error code (plain `errors.Errorf`, not `errdefs` with status code). If the prior job has been discarded (>10 s), the new different input is accepted and executed as a completely new build; history will contain two separate `BuildHistoryRecord`s with the same `Ref` but different `Generation`? Actually `UpdateRef` increments `Generation` on mutation (`solver/llbsolver/history/buildhistory.go:335`), but a new `STARTED` does not increment — it overwrites `h.active[Ref]` (`buildhistory.go:539`). If a `COMPLETE` record already exists in BoltDB with that `Ref`, the next build's `COMPLETE` will overwrite it via `b.Put([]byte(rec.Ref), dt)` (`buildhistory.go:485`), losing the prior result.

Operationally this means `Ref` is not safe to reuse as an idempotency key with mutation: BuildKit documents no constraint and provides no `400 Conflict` for differing payloads.

### 4. How long is deduplication state retained?

**Two distinct lifetimes, neither intended as idempotency retention:**

- **Submission deduplication (the only true dedup):** `jl.jobs` map entry lifetime = job execution time + 10 s grace (`solver/jobs.go:890-895`). Comment: *don't clear job right away. there might still be a status request coming to read progress*. No persistence; restart clears all. No configuration knob.
- **Observability history (not deduplication):** BoltDB `recordsBucket` retention governed by `config.HistoryConfig` defaults `MaxAge 48h`, `MaxEntries 50` (`solver/llbsolver/history/buildhistory.go:82-84`), GC every 120 s (`buildhistory.go:131-137`) deleting only when *both* `len(records) >= MaxEntries` and record is older than `MaxAge` (`buildhistory.go:175-196`). Pinned records (`Pinned=true`) are excluded. Deleted refs are announced via `BuildHistoryEventType_DELETED` and `deleted` map with refcounting (`buildhistory.go:58-60,241-266,770-870`). These records are not consulted during `NewJob`; they cannot prevent duplicate work.

Thus idempotency-style deduplication window is effectively *while job is running + 10 s*, not hours.

### 5. Can a network/client failure accidentally start duplicate work?

**Yes, on both sides:**

- **Server still runs on client disconnect, but client retry creates duplicate if it reuses `Ref`:** `Solve` `ctx` is the gRPC stream context (`control/control.go:409`, `solver/llbsolver/solver.go:195`). Cancellation (client timeout/disconnect) cancels `br.Solve` and exporter via `ctx.Done()` (`solver/llbsolver/solver.go:295-310`, `solver/jobs.go:1251-1260` handling `IsCanceled`). However history finalizer still runs with `context.WithoutCancel` to record failure (`solver/llbsolver/history.go:59,74`). If the client retries with the *same* `Ref` before 10 s grace, it gets `job ID exists` and must not retry with same `Ref` — but the default client path generates a fresh `Ref = identity.NewID()` per `Client.solve` (`client/solve.go:107`), so automatic retry would *not* reuse `Ref` and would start duplicate work with a new `Ref` (new vertex graph, though cache may dedup vertices). If caller *does* supply stable `opt.Ref` (used e.g. by `client/client_control_test.go:68` for history tests), retry after cancel will either collide or, after grace, start an independent build doing redundant execution despite cache sharing.
- **No at-most-once guard:** There is no request fingerprint, no deduplication store, no transactional outbox, and no idempotency token. The `gateway` path `gatewayForwarder.RegisterBuild(ctx, id, fwd)` (`solver/llbsolver/solver.go:280`) registers concurrently but does not dedupe — it is for forwarding, not submission idempotency. Vertex-level `flightcontrol.Group` (`solver/jobs.go:995-998`) prevents duplicate *vertex* execution within a job graph, but not duplicate *job* submission.
- **Crash after commit but before response:** Since `Solve` is synchronous, the commit point is the final `SolveResponse` write. If daemon crashes after `COMPLETE` history update (`history.go:271`) but before gRPC response flushes, client retry with same `Ref` cannot return the original operation after 10 s — it will start new work. Within 10 s it gets `exists`, not the prior response. `ListenBuildHistory` with `Ref` filter + `EarlyExit` could be polled to fetch `COMPLETE` record (`solver/llbsolver/history/buildhistory.go:770-887`), but `Client.solve` does not implement this resume logic (`client/solve.go:369-397` only tails `Status`).

## Architectural Decisions

- **Synchronous unary `Solve` over persistent queue:** Chose blocking `Solve` (`control/control.go:409`) that holds the gRPC handler for the entire build, delegating concurrency to in-memory `Solver.jobs` and vertex scheduler, instead of a durable queue/prepare-commit pattern. Simpler, no need for persistence of queued builds, but sacrifices at-most-once and resume.
  Tradeoff: low operational complexity vs no durability across restarts and no client timeout resilience.

- **`Ref` as ephemeral job ID, not idempotency key:** `Ref` is `identity.NewID()` per solve (`client/solve.go:107`), stored only in `Solver.jobs` map (`solver/jobs.go:400-403`). No hashing of `Definition` or `FrontendAttrs` into a fingerprint; `BuildHistoryRecord.Ref` is just the same string for observability (`solver/llbsolver/history.go:36`). Enables concurrent builds with distinct IDs but provides no deduplication semantics.

- **Content-addressed vertex cache as implicit dedup:** `state` keyed by `v.Digest()` and `dgstWithoutCache` variant (`solver/jobs.go:545-571`), `actives` map sharing vertices across jobs, `sharedOp` flightcontrol for `CacheMap`/`Exec` (`solver/jobs.go:1137-1286`). This gives powerful reuse and at-most-once *vertex* execution, but masks absence of job-level idempotency — two jobs with same DAG will share underlying `state` until one `Discard` prunes it.

- **BoltDB history as post-hoc observability, not commit log:** `history.Queue` with `recordsBucket` and `contentStore` lease-protected blobs (`buildhistory.go:418-487`, `history.go:112-178`) is written on `STARTED`/`COMPLETE`. Not consulted for admission; GC is age+count based (`buildhistory.go:152-198`). Decouples durability from execution path, simplifying hot path at cost of no transactional guarantee.

- **10 s job map grace for status tailing, not for idempotency:** `Job.Discard` delayed delete (`solver/jobs.go:890-895`) allows `Status` RPC (`solver/llbsolver/solver.go:465-481` falling back to `s.solver.Get`) to succeed briefly after build finishes; not sized or documented as idempotency window.

## Notable Patterns

- **Flightcontrol per-vertex, not per-request:** `flightcontrol.Group[digest.Digest]` / `Group[[]*CacheMap]` / `Group[*execRes]` in `sharedOp` (`solver/jobs.go:995-998`) coalesces concurrent identical vertex evaluations — a mature pattern for cache/merge but scoped inside a build, not across builds.

- **Lease + content store for history blobs:** `OpenBlobWriter` creates ephemeral lease `history_migration_*` and content writer (`buildhistory.go:567-593`, `history.go:97-178`); `addResource` migrates orphan blobs (`buildhistory.go:283-307`). Finalizer `AcquireFinalizer` with `trigger`/`done` channels (`buildhistory.go:489-530`) coordinates async trace commit — a careful lease lifecycle pattern reused from containerd.

- **Progress multi-writer fanout:** `state.mpw progress.MultiWriter` plus `allPw` set (`solver/jobs.go:58-61,665-685`) and `notifyStarted` (`solver/jobs.go:1368-1388`) fan out vertex progress to all contributing jobs; enables `Status` streaming without durable log.

- **CacheOptions duplicate guard (narrow):** `findDuplicateCacheOptions` (`control/control.go:778-806`) hashes registry refs or full `Type+Attrs` via `hashstructure` (`control/control.go:808-824`) — a rare explicit dedup, but only for cache exporter config and it *rejects* rather than dedupes.

- **Gateway forwarder registration race comment:** `RegisterBuild` before `recordBuildHistory` due to 3 s timeout (`solver/llbsolver/solver.go:276-282`) — indicates awareness of blocking history may starve other paths, yet no fix for atomicity.

## Tradeoffs

- **Fast in-memory admission vs crash safety:** `NewJob` is a mutex + map insert (`solver/jobs.go:688-715`), O(1) and no I/O. Avoids BoltDB transaction on hot path, but any crash loses all in-flight jobs and their vertex graphs; no replay.
- **Synchronous response vs resumability:** Blocking `Solve` gives simple client semantics (one RPC, one response) but forces client to hold connection for minutes/hours; timeout forces retry without resume. An async prepare/commit + poll model would enable durable acceptance and retry-to-fetch, at cost of API complexity.
- **Random `Ref` per build vs stable fingerprint:** Random avoids collisions and simplifies concurrent builds, but eliminates natural idempotency (same inputs produce different `Ref`). Fingerprinting `Definition` (LLB digest) would enable dedup but would require centralizing `Ref` generation and handling `different input, same key` conflicts.
- **Vertex sharing vs isolation:** `loadUnlocked` sharing `actives[dgst]` across jobs (`solver/jobs.go:530-663`) maximizes cache reuse and memory efficiency; however `addJobs` propagation (`solver/jobs.go:257-296`) and `deleteIfUnreferenced` (`solver/jobs.go:747-783`) couple jobs' lifetimes, making job-level GC subtle and reliant on correct `parents`/`childVtx` bookkeeping.
- **History as best-effort vs WAL:** Writing `STARTED` then later `COMPLETE` with `Generation` bump (`buildhistory.go:335`) and async trace (`history.go:299-332`) keeps build hot path fast; failure to write history is `Internal` but does not roll back build. A WAL would need fsync per build start.

## Failure Modes / Edge Cases

- **Duplicate `Ref` with racing `NewJob`:** Two concurrent `Solve` with same `Ref` — mutex serializes, one gets `exists`, other proceeds. Caller retried due to timeout may see success-after-failure non-deterministically. No idempotent winner semantics.
- **Reconnect after 10 s grace:** Client timed out after `Solve` accepted but before `SolveResponse` flushed, retries with same `Ref` after 10 s — new job starts, duplicates work, may overwrite prior `COMPLETE` record (`b.Put` same key, `buildhistory.go:485`).
- **Daemon crash between `NewJob` and `Update(STARTED)`:** In-memory job exists, but `h.active` map and BoltDB lack record. Restart clears `jl.jobs`; orphan vertices GC'd; history shows no record — caller cannot discover outcome via `ListenBuildHistory` (`buildhistory.go:770-847` early exit).
- **BoltDB `Update` failure on `STARTED`:** `recordBuildHistory` returns `Internal` (`history.go:52-56`) and `stopTrace` is cleaned; but job already created and `Discard` will still run — leaks job that will timeout? Actually `Solve` returns error without ever running `br.Solve`, but job remains until `Discard` (deferred in `llbsolver.Solve`, `solver.go:215`). Client sees error, may retry, hitting `exists` until discard.
- **`Ref` collision across restarts:** Since `Ref` is random, collision probability negligible, but user-supplied `opt.Ref` stable across retries could collide with a new random `Ref` after restart — no namespacing.
- **Vertex `IgnoreCache` digest split:** `dgstWithoutCache` branch (`solver/jobs.go:547-565`) creates separate `state` for same logical vertex with/without cache ignore. Two jobs with same `Ref` but differing `IgnoreCache` option share no state, but dedup check does not distinguish — one blocked by the other's unrelated `Ref` reservation.
- **History GC while `Listen` active:** `h.delete` checks `refs` refcount (`buildhistory.go:241-246`) and defers actual `b.Delete` until `Listen` with `Ref` filter releases (`buildhistory.go:786-795,883-795`). However `Find` via `solver.Get` for `Status` does not hold `refs`, so `Status` after GC may return `os.ErrNotExist` even though live job still exists — `Solver.Status` fallback handles via `history.Status` then `solver.Get` (`solver/llbsolver/solver.go:465-481`) but a race window exists.
- **Session `Ref` vs Build `Ref` aliasing:** `Job.SessionID = sessionID` (`solver/llbsolver/solver.go:266`) and `session.NewGroup(req.Session)` for cache resolvers (`control/control.go:483`) use same session ID as `SolveRequest.Session` (`control.proto:69`). Duplicate session IDs across builds do not cause `job ID exists` — separate namespaces, but confusion could cause wrong content store lease association.
- **No validation of `Definition` digest vs `Ref`:** Caller can send empty `Definition` with gateway path (`solver/llbsolver/solver.go:273-274`) — duplicate `Ref` with different gateway state shares no check, `RegisterBuild` may overwrite prior `fwd`.

## Future Considerations

- **Add idempotency contract to `SolveRequest`:** Introduce `IdempotencyKey` (or define `Ref` as such) plus `RequestFingerprint` = hash(`Definition`+`Frontend`+`FrontendAttrs`+`FrontendInputs`+`Cache`+`Exporters`+`Entitlements`). On `NewJob`, if `jobs[key]` exists, compare fingerprint: identical → return existing operation handle or join; different → return `AlreadyExists` with `FAILED_PRECONDITION`/conflict details (grpc `codes.AlreadyExists` with `errdefs`). Requires persisting fingerprint alongside `jobs` in BoltDB for restart recovery, or documenting that idempotency is only within process lifetime.
- **Durable acceptance before ACK:** Wrap `NewJob` + `history.Update(STARTED)` + `Status` progress pipe creation in a single BoltDB transaction (or write-ahead log) before returning success to gRPC. On crash-restart, replay `STARTED` records to rebuild `jobs` map or mark as failed. Currently `recordBuildHistory` is post-accept; inversion needs careful ordering with `gatewayForwarder.RegisterBuild` timeout.
- **Resume/return-original-operation RPC:** Add `Control.GetOperation(Ref)` or idempotent `Solve` that returns existing `SolveResponse`/error if `COMPLETE` record exists, enabling `Crash after commit but before start response` test. Leverage existing `history.Status` blob replay (`buildhistory.go:351-416`) and `BuildHistoryRecord.ExporterResponse` (`history.go:68-72`) to reconstruct response without re-executing. Client `solve` would then implement retry with backoff, checking `GetOperation` after `DeadlineExceeded`.
- **Configurable retention for dedup state:** Separate `IdempotencyRetention` (e.g., 24 h) from `HistoryConfig.MaxAge/MaxEntries`; store `idempotency` bucket keyed by `(key, fingerprint) -> response/error + expiry`; GC independently. Current 10 s job grace is insufficient; 48 h / 50 entries history not suitable as dedup store due to count threshold.
- **Distinguish transient vs durable errors for retry:** Classify `recordBuildHistory` `Internal` vs `exists` vs build `Err` via `grpcerrors.AsGRPCStatus` (`buildhistory.go:637-641`); document safe retry predicates so callers do not retry non-idempotent exporter mutations (registry push already has `retryhandler` but not at job level).
- **Test coverage:** Add integration test parallel to `TestJobsIntegration` (`solver/jobs_test.go:22`) that: (a) calls `Solve` with explicit `Ref`, kills context after `NewJob` but before response, retries with same `Ref` + same/different `Definition`, asserts idempotent return vs conflict error; (b) fills `recordsBucket` to exceed `MaxEntries+MaxAge` and verifies dedup window still holds if separate; (c) crashes daemon (or simulates `NewJob`+`Update` failure) and asserts `ListenBuildHistory` recovers.

## Questions / Gaps

- No evidence found of a durable idempotency table or fingerprint for `SolveRequest` submissions — searched `solver/jobs.go`, `control/control.go`, `solver/llbsolver/solver.go`, `history/*.go`, `client/solve.go`, and `control.proto` for `idempot`, `fingerprint`, `dedupl`, `hash`, `Ref` usage beyond existence check.
- No evidence found of transactional coupling between `NewJob` and history commit — `NewJob` uses `sync.RWMutex`, history uses `bolt.Tx`+`leaseManager`; no shared transaction.
- No evidence found of duplicate/conflict tests for `Ref` reuse with differing inputs — `control/control_test.go` only tests `findDuplicateCacheOptions` (`control/control_test.go:10-41`), not job submission; `solver/scheduler_test.go` and `solver/jobs_test.go` use distinct `job0`/`job1` IDs without collision cases; `client/client_control_test.go:63-120` tests delete/list history, not idempotency.
- No evidence found of client retry behavior for `Solve` — `client/solve.go:337-343` shows single `Solve` RPC; no `retryhandler` wrapping at this layer (retryhandler only for registry/content fetch `util/resolver/retryhandler/retry.go:1`, `util/pull/pull.go:153`).
- No evidence found of validation/revalidation around a prepare/start boundary — BuildKit has no prepare phase; `Control.Solve` is monolithic. The closest boundaries are `Frontend` input validation and `Exporter` resolution (`control/control.go:420-464`), but none re-validate after durable write.
- Gaps: Is `Ref` intended to be idempotency key or just a tracing correlation? Docs (`docs/dev/solver.md:178`, `control.proto:62`) describe `Ref` as build reference identifier, not idempotency token. Need maintainer clarification on intended retry contract.

---

Generated by `Dimension 01.03: Durable Submission and Idempotent Acceptance` against `buildkit`.
