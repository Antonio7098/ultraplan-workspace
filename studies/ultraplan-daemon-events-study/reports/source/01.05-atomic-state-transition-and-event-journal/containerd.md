# Source Analysis: containerd

## Dimension 01.05: Atomic State Transition and Event Journal

### Source Info

| Field | Value |
|-------|-------|
| Name | containerd |
| Path | `studies/ultraplan-daemon-events-study/sources/containerd` |
| Language / Stack | Go 1.24+, go.etcd.io/bbolt (vendored), gRPC/TTRPC, github.com/docker/go-events |
| Analyzed | 2026-09-02 |

## Summary

containerd implements a strongly consistent, single-writer authoritative state store in `core/metadata/db.go:84` (`DB` wrapping `*bolt.DB`/`Transactor` at `core/metadata/bolt.go:30`) with full namespace isolation under the `v1/<namespace>/<object>/<key>` schema (`core/metadata/buckets.go:1`). Every meaningful mutation (containers, images, content blobs/ingests, snapshots, leases, sandboxes, namespaces) runs inside one short, bounded `bbolt.Tx` via `view`/`update` helpers (`core/metadata/bolt.go:37`, `core/metadata/bolt.go:47`) that support context-bound transaction reuse (`core/metadata/boltutil/context.go:31`). Schema ownership and migrations are explicit (`core/metadata/migrations.go:42`, `core/metadata/db.go:156`) and the daemon deliberately rejects `WithEventsPublisher` inside a transaction (`core/metadata/db.go:288`).

The operational event system (`core/events/exchange/exchange.go:36`, `core/events/events.go:67`) is a separate, purely ephemeral in-memory broadcaster (`github.com/docker/go-events`) with **no durable journal, no transactional outbox, and no shared transaction/connection owner** with the state commit. All state stores publish *after* their `Update` has committed (`core/metadata/images.go:168`, `core/metadata/content.go:235`, `plugins/services/containers/local.go:139`) and `GarbageCollect` even publishes asynchronously (`core/metadata/db.go:441`). A crash between the bolt commit and the `Publish`/`Forward` call yields a committed-but-unobserved state. The design trades durability and exactly-once observability for simplicity and bounded bolt write latency; the vendor explicitly documents `NoSync` as a crash-safety tradeoff (`plugins/metadata/plugin.go:66`).

There is no optimistic revision field. Concurrent writers serialize on bbolt's single writer plus a `wlock` RWMutex that blocks writers during GC mark (`core/metadata/db.go:89`); last-write-wins is the conflict policy.

## Rating

**5 / 10 — Present but inconsistent / fragile**

**Rationale:** State-transition atomicity is strong: short CoW transactions, composable `WithTransaction` reuse, explicit `dirty`/GC dirty tracking, and versioned migrations with tests. The companion requirement — atomically appending a durable operational event in the same transaction/connection — is intentionally absent. `DB.Publisher` returns `nil` inside a transaction (`core/metadata/db.go:289`) and every call site publishes after commit with best-effort, lossy broadcast; shim forwarding adds only an in-memory requeue (`pkg/shim/publisher.go:102`). No `ExpectedVersion`/`FailedPrecondition` (revision) guard exists. This satisfies "avoid full event sourcing" but violates the dimension's core atomicity question. The model is coherent and well-known, not ad-hoc, yet fragile for observers that require exactly-once or at-least-once delivery.

## Evidence Collected

Every entry MUST include a file path with line numbers. Format: `path/to/file.ts:NN`.

| Area | Evidence | File:Line |
|------|----------|-----------|
| Authoritative store — DB struct & version | `schemaVersion="v1"`, `dbVersion=4`, `DB { db Transactor; wlock; dirty; dirtySS; dirtyCS }` | `core/metadata/db.go:42` |
| Authoritative store — schema docs | `v1/<namespace>/<object>/<key>` layout, blob/ingest/snapshot/container buckets | `core/metadata/buckets.go:1` |
| Authoritative store — bucket keys | `bucketKeyVersion`, `bucketKeyDBVersion`, per-object buckets | `core/metadata/buckets.go:146` |
| Transaction abstraction | `Transactor { View, Update }` and `view`/`update` reusing `boltutil.Transaction(ctx)` | `core/metadata/bolt.go:30` |
| Transaction reuse | `WithTransaction`/`Transaction` context key storing `*bolt.Tx` | `core/metadata/boltutil/context.go:25` |
| DB View/Update & GC gate | `View` delegate, `Update` takes `wlock.RLock`, fires `mutationCallbacks` after commit | `core/metadata/db.go:267` |
| DB gate for events | `Publisher(ctx)` returns nil if inside transaction: "Do not publish events within a transaction" | `core/metadata/db.go:288` |
| DB open & durability knobs | `bolt.Open(path,0644,&options)` with `NoFreelistSync=true`, `NoSync`/`NoGrowSync`, `NoStatistics=true`, `Timeout=0` (block forever) | `plugins/metadata/plugin.go:132` |
| BoltConfig | `BoltConfig { ContentSharingPolicy; NoSync bool }` with crash-loss warning | `plugins/metadata/plugin.go:49` |
| Migration ownership | `migrations = []migration{ v1.1 addChildLinks, v1.2 migrateIngests, v1.3 noOp, v1.4 migrateSandboxes }` | `core/metadata/migrations.go:42` |
| Migration driver | `Init` walks `migrations[i:]` inside one `db.Update`, writes `bucketKeyDBVersion` varint | `core/metadata/db.go:156` |
| GC & WAL/checkpoint | No WAL segment; bbolt CoW + freelist; snapshot metastore `MetaStore { dbL sync.Mutex; db *bolt.DB; opts bolt.Options { NoStatistics } }` | `core/snapshots/storage/metastore.go:69` |
| Event Envelope & interfaces | `Envelope { Timestamp, Namespace, Topic, Event }`, `Publisher { Publish }`, `Forwarder { Forward }`, `Subscriber { Subscribe }` | `core/events/events.go:26` |
| In-memory broadcaster | `Exchange { broadcaster *goevents.Broadcaster }`, `Forward`/`Publish` via `broadcaster.Write`, `Subscribe` via `NewChannel+NewQueue+NewFilter` | `core/events/exchange/exchange.go:36` |
| State then event — images | `update` tx `130-164`, then `Publisher.Publish "/images/create"` outside tx (returns publish error though DB committed) | `core/metadata/images.go:124` |
| State then event — images update/delete | Same pattern for `/images/update` and `/images/delete` (`dirty.Add(1)` before publish) | `core/metadata/images.go:180` , `core/metadata/images.go:282` |
| State then event — content delete | `update` deletes blob bucket + `dirtyCS`, then `Publish "/content/delete"` | `core/metadata/content.go:204` |
| State then event — content commit | `namespacedWriter.Commit`: `pre-Sync` outside lock, `update` inside lock, then `Publish "/content/create"` after | `core/metadata/content.go:582` |
| State then event — snapshots | `Prepare` createSnapshot then `Publish "/snapshot/prepare"`; `Commit` does `Snapshotter.Commit` *inside* `update` then `Publish "/snapshot/commit"` | `core/metadata/snapshot.go:274` , `core/metadata/snapshot.go:520` |
| State then event — containers service | `withStoreUpdate` → `Store.Create` then `publisher.Publish "/containers/create"`; same for update/delete | `plugins/services/containers/local.go:122` |
| State then event — GC sweep | Collect `events []namespacedEvent` inside `db.Update` sweep (comment "queue event to publish after successful commit"), then `wg.Go(m.publishEvents)` async after commit | `core/metadata/db.go:395` |
| GC async publish | `publishEvents` uses `context.Background()` + `namespaces.WithNamespace`, logs but does not fail GC on publish error | `core/metadata/db.go:348` |
| Shim forward requeue | `RemoteEventsPublisher { requeue chan 2048 }`, `Publish` → `Forward` RPC, `queue(i)` on error, `processQueue` retry `1s*count` max 5 | `pkg/shim/publisher.go:80` |
| Revision / OCC — absent | `Update` reads, mutates, `UpdatedAt=Now()`, no input-revision check; only `ErrAlreadyExists`/`ErrNotFound` | `core/metadata/containers.go:174` |
| Revision — images same | Field-mask replace or per-field patch, no version guard; delete conditional on digest returns `ErrNotFound` not `FailedPrecondition` | `core/metadata/images.go:180` |
| Revision — timestamps | `WriteTimestamps`/`ReadTimestamps` on `createdat`/`updatedat` keys, never compared | `core/metadata/boltutil/helpers.go:110` |
| Writer contention | `wlock RWMutex`: GC holds `Lock` 383-489, writers hold `RLock`; `Metastore.dbL sync.Mutex` serializes `Begin` | `core/metadata/db.go:89` , `core/snapshots/storage/metastore.go:72` |
| GC resources & leases | `ResourceContent/Snapshot/Container/Image/Lease`, lease `AddResource/DeleteResource` all inside same `update` as object create (atomic link) | `core/metadata/gc.go:35` , `core/metadata/leases.go:38` |
| Proto event topics | `/containers/*`, `/images/*`, `/content/*`, `/snapshot/*`, `/namespaces/*`, `/tasks/*`, `task.proto` etc. — payload types only, not storage | `api/events/container.proto:1` , `core/runtime/events.go:24` |
| No durable journal — search | No bucket for events, no WAL file, no `eventstore`/`Journal` in core/pkg; `internal/eventq` is ephemeral `discardAfter` queue | `internal/eventq/eventq.go:56` |

## Answers to Dimension Questions

### 1. Can state commit without its corresponding durable event?

**Yes — by design and in the common case.** Every mutation path commits the `bbolt.Tx` first and only then calls `Publisher.Publish` outside the transaction. The canonical pattern is visible in `core/metadata/images.go:124` (Create), `core/metadata/images.go:180` (Update), `core/metadata/images.go:282` (Delete), `core/metadata/content.go:204` (content Delete), `core/metadata/snapshot.go:274` (Prepare), and `plugins/services/containers/local.go:122` (service Create). `DB.Publisher` actively prevents in-transaction publishing by returning `nil` when `boltutil.Transaction(ctx)` is present (`core/metadata/db.go:288`). If the publisher is nil, unreachable, or the broadcaster has no subscribers, the bolt commit has already succeeded and there is no durable event record to replay. `GarbageCollect` makes this even looser: events collected inside the sweep transaction (`core/metadata/db.go:395`) are flushed asynchronously via `wg.Go(m.publishEvents)` (`core/metadata/db.go:441`) with failure only logged (`core/metadata/db.go:363`).

A crash injected between the `db.Update` return and the subsequent `Publish` yields the observable inconsistency: metadata readers (View/List) see the new state while event subscribers never receive the envelope. Because the exchange is in-memory `goevents.Broadcaster` (`core/events/exchange/exchange.go:36`) with no backing `meta.db` bucket, the lost event cannot be reconstructed without rescanning state.

The exception path is `core/metadata/snapshot.go:546` `Commit`, where `s.Snapshotter.Commit(ctx, nameKey, bkey)` runs *inside* the bolt transaction to avoid "keys becoming out of sync" (`core/metadata/snapshot.go:639`). If that inner commit succeeds but the outer `db.Update` later fails, `core/metadata/snapshot.go:660` logs "transaction failed after commit, snapshot should be removed" — an admitted out-of-sync artifact requiring manual cleanup. This is state-backend coupling, not event atomicity.

### 2. Can an event commit for a state transition that later rolls back?

**No durable event can be rolled back because there is no durable event to begin with; in-memory events cannot precede commit by construction.** The framework forbids publishing inside the transaction (`core/metadata/db.go:289` guard), so an event never commits before the state transaction. The closest risk is pre-sync side effects: `core/metadata/content.go:592` performs `fp.Sync()` outside the `Update` lock, and `core/metadata/snapshot.go:642` notes the backend `Commit` happens inside the bolt tx. In those cases the *external* backend (content store, overlay snapshotter) may have mutated before the bolt transaction rolls back, but the *event* still has not been published — the violation is state-backend divergence, not phantom event.

Shim → daemon `Forward` is similarly post-operation (`pkg/shim/publisher.go:127`) and only retries on failure (`pkg/shim/publisher.go:102`) with eviction after 5 attempts. No path publishes an event and then rolls back the state transaction.

### 3. Is the event log an audit/observation journal or the sole source of truth?

**Audit/observation journal only; the sole source of truth is the current-state `meta.db` (`core/metadata/db.go:84`) plus the snapshot/content backends.** All read paths (`Get`/`List`/`Walk`) use `view` on bolt (`core/metadata/containers.go:55`, `core/metadata/images.go:50`, `core/metadata/buckets.go:183`) or backend stores; none replay an event log. The `Exchange` (`core/events/exchange/exchange.go:41`) keeps zero history — `Subscribe` creates a fresh `goevents.NewChannel(0)` + `goevents.NewQueue` (`core/events/exchange/exchange.go:132`) and only delivers events published after subscription. Restart clears all in-flight events. Proto definitions (`api/events/*.proto`) describe payloads for notification (e.g., `ContainerCreate`, `ImageDelete`, `SnapshotCommit`) not for rebuilding state. Containerd's documentation and bucket schema (`core/metadata/buckets.go:1`) explicitly frame the store as "fully namespaced" current-state storage, not an append-only log.

### 4. Are transactions short and bounded?

**Yes, generally short and bounded by convention and locking.** The bolt-options comments require backend snapshotters to "commit fast and reliably to prevent metadata store locking and minimizing rollbacks" (`core/metadata/snapshot.go:638`). Held state includes a single bucket create/read/write plus lease linkage (`core/metadata/images.go:140`, `core/metadata/leases.go:337`) and timestamp writes (`core/metadata/boltutil/helpers.go:110`). No transaction spans RPC streams or blob I/O: `core/metadata/content.go:582` does `Writer.Commit` with a pre-sync outside the `l.RLock` + `update` critical section, and `plugins/metadata/plugin.go:132` sets `NoStatistics` and `NoFreelistSync` to reduce contention. `Update` itself only holds `wlock.RLock` (`core/metadata/db.go:272`) while `GarbageCollect` holds `wlock.Lock` (`core/metadata/db.go:384`) to prevent writers during mark-and-sweep, but mark (`View`) and sweep (`Update`) are themselves one transaction each. The longest transaction is the multi-bucket migrations inside `Init` (`core/metadata/db.go:162`) and the `Commit` path that nests `Snapshotter.Commit` inside the bolt tx (`core/metadata/snapshot.go:642`), both bounded to metadata writes.

### 5. How are stale state revisions rejected?

**They are not — there is no optimistic revision / `ExpectedVersion` / `FailedPrecondition` guard.** Searches for `revision`/`Revision` hit only `version.Revision` (build SHA) and kernel strings; `Expected`/`ETag`/`optimistic` return zero hits outside vendor. `containerStore.Update` (`core/metadata/containers.go:174`) reads via `readContainer` (`core/metadata/containers.go:205`), applies a field-mask (`core/metadata/containers.go:230`), and blindly overwrites `UpdatedAt = time.Now().UTC()` (`core/metadata/containers.go:270`) with no comparison to an input timestamp or sequence. `ImageStore.Update` (`core/metadata/images.go:180`) behaves identically, preserving only `CreatedAt`. All stores return at most `ErrAlreadyExists` on `CreateBucket` race (`core/metadata/containers.go:150`, `core/metadata/images.go:150`) and `ErrNotFound` on missing bucket; no `ErrConflict`. The one conditional check — `Delete` with `options.Target.Digest` mismatch (`core/metadata/images.go:316`) — returns `ErrNotFound` as an ad-hoc CAS, not a first-class revision field. Writer serialization relies entirely on bbolt's single writer + `wlock` (`core/metadata/db.go:89`) and context-tx reuse (`core/metadata/bolt.go:37`); concurrent `Update` callers execute sequentially, last write wins.

## Architectural Decisions

| Decision | Evidence | Effect | Tradeoff |
|----------|----------|--------|----------|
| **Single bbolt file `meta.db` as authoritative store** | `plugins/metadata/plugin.go:169` `filepath.Join(root,"meta.db")`, `bolt.Open(path,0644,&options)` `plugins/metadata/plugin.go:184`, `core/metadata/db.go:84` | One ordered, namespaced current-state DB with atomic `Update` for images/containers/snapshots/content/leases together | Single-writer lock, GC pauses writers, no horizontal sharding; `no_sync` toggles crash safety |
| **Two-tier bolt options for durability vs perf** | `BoltConfig.NoSync` doc `plugins/metadata/plugin.go:67`, `NoFreelistSync=true` `plugins/metadata/plugin.go:137`, `NoSync+NoGrowSync` `plugins/metadata/plugin.go:161`, `core/snapshots/storage/metastore.go:90` `NoStatistics=true` | Reduces `fsync`/`flock` overhead; freelist kept in memory to avoid corruption (`bbolt#1,#6`) | `NoSync` admits torn pages on power loss; durable events still absent regardless |
| **Explicit `Publisher == nil` inside tx** | `DB.Publisher` checks `boltutil.Transaction` `core/metadata/db.go:288`, `update` rejects `!Writable` `core/metadata/bolt.go:51` | Forces callers to publish outside transaction; prevents holding bolt lock during broadcast | Guarantees state-without-event window; error from `Publish` surfaces despite committed state |
| **Best-effort in-memory `go-events` broadcast** | `Exchange.broadcaster` `core/events/exchange/exchange.go:36`, `Publish` builds `Envelope{time.Now(), namespace, topic, MarshalAny}` `core/events/exchange/exchange.go:99`, `Forward` vs `Publish` `core/events/events.go:72` | Subscribers get near-real-time notifications with filter support (`goevents.NewFilter` `core/events/exchange/exchange.go:155`) | Ephemeral, lossy, no replay, no persistence, late subscribers miss history |
| **Async GC event emission** | `publishEvents` with `context.Background` `core/metadata/db.go:348`, `GarbageCollect` `wg.Go(m.publishEvents)` `core/metadata/db.go:441` | Keeps sweep `Update` (`core/metadata/db.go:396`) short; snapshot/content backend cleanup runs in parallel via `wg` `core/metadata/db.go:452` | GC-delete observed out-of-order or not at all if publish fails |
| **Lease linkage inside same tx as object** | `addImageLease` inside `Create` `core/metadata/images.go:140`, `addSnapshotLease`/`removeSnapshotLease` inside `Commit` `core/metadata/snapshot.go:553`, `core/metadata/leases.go:337` | Referential integrity: GC can safely `scanRoots` + `references` + `scanAll` + `remove` (`core/metadata/gc.go:495`, `core/metadata/gc.go:832`, `core/metadata/gc.go:936`, `core/metadata/gc.go:1034`) without orphans | Tightly couples lease lifecycle to caller ctx; no lease without namespace |
| **GC guarded by `wlock` RWMutex + dirty flags** | `wlock sync.RWMutex` `core/metadata/db.go:93`, `dirty atomic.Uint32` `core/metadata/db.go:99`, `dirtySS/dirtyCS` `core/metadata/db.go:105`, `mutationCallbacks` `core/metadata/db.go:108` | Writers block only during `getMarked` + sweep; `dirty` signals scheduler (`plugins/gc/scheduler.go:34` thresholds) | Mutation burst can starve writer or delay reclamation; scheduler's adaptive interval is heuristic |

## Notable Patterns

- **Context-bound transaction reuse (composable writes).** `boltutil.WithTransaction` (`core/metadata/boltutil/context.go:31`) lets `withStore`/`withStoreUpdate` (`plugins/services/containers/local.go:207`) batch lease+object mutations in one `db.Update` without nested transactions. Seen in `containers.go:63`, `images.go:58`, `content.go:79`, `snapshot.go:88`, `namespaces.go:28`, `leases.go:59` — all call `view`/`update` which reuse the context tx. Enables atomic lease linking without exposing raw `bolt.Tx` to callers.

- **Field-mask guarded partial updates.** `Update` methods accept `fieldpaths ...string` and only mutate listed fields (`core/metadata/containers.go:230` for `labels/spec/extensions/image/snapshotkey`; `core/metadata/images.go:209` for `labels/target/annotations`). If no fieldpaths, `Update` either whitelists allowed full-replace fields or rejects immutable ones (`Snapshotter`, `Runtime.Name` `core/metadata/containers.go:218`). Preserves operator intent without full-object overwrite.

- **Two-phase snapshot create.** `createSnapshot` pattern (`core/metadata/snapshot.go:274` via `plugins/services/containers/local.go:207` helpers) first reserves `NextSequence` + `createKey` inside a tx, then calls the backend `Snapshotter.Prepare`, then writes final bucket fields inside a second tx — avoids holding the global bolt lock during overlay work.

- **Backend commit inside metadata transaction (intentional coupling).** `snapshotter.Commit` calls `s.Snapshotter.Commit` inside `db.Update` (`core/metadata/snapshot.go:642`) with comment "prevent metadata store locking and minimizing rollbacks" but notes risk of out-of-sync if bolt later fails (`core/metadata/snapshot.go:639`). Mirrors the classic "two resources, one transaction" hazard the dimension warns about.

- **Ephemeral `eventq` for CRI with discard.** `internal/eventq/eventq.go:56` `EventQueue[T]` buffers `discardQueue` and drops after `discardAfter` if no subscriber — documents the project's comfort with lossy delivery already seen in `Exchange`.

## Tradeoffs

- **Bounded write latency vs at-least-once observability.** Keeping bolt transactions short by deferring `Publish` keeps `Update` latency predictable and avoids deadlocking the single writer on slow subscribers, but observers can miss committed state with no way to recover except full rescan.

- **Simple current-state model vs powered remediation.** Avoiding an outbox/journal or change-data-capture avoids schema bloat and migration cost (only `bucketKeyVersion`/`bucketKeyDBVersion` `core/metadata/buckets.go:147` + `core/metadata/migrations.go:42`), at the cost of no transactional exactly-once consumer and no built-in audit trail beyond bolt file.

- **Shared vs isolated content policy.** `BoltConfig.ContentSharingPolicy` (`plugins/metadata/plugin.go:65`) defaults to `shared` so blobs pulled in one namespace become visible everywhere (fewer pulls, weaker isolation); `isolated` requires re-proving content via full ingest (`core/metadata/content.go:49`).

- **Durability knob.** `NoSync` improves write throughput 10×+ by skipping `fdatasync` (`plugins/metadata/plugin.go:68`) but widens the crash window for both state and already-published events. Default is safe (`NoSync:false` `plugins/metadata/plugin.go:104`).

- **Last-write-wins vs OCC.** No revision check avoids client-supplied `If-Match` plumbing and keeps API surface small (`containers.Store` `core/containers/containers.go:93`), but silent lost-update is possible for concurrent `Update` racing on `labels` or `target`. The bbolt serializer hides conflicts rather than surfacing them.

## Failure Modes / Edge Cases

- **Crash between commit and publish (primary dimension failure).** After `db.Update` returns (`core/metadata/images.go:164`, `core/metadata/content.go:231`, `plugins/services/containers/local.go:137`) and before `Publisher.Publish`, `meta.db` reflects the new object while no `Envelope` was written. Restart sees state but not event; consumers relying on `Subscribe` (`core/events/exchange/exchange.go:128`) to drive downstream work (e.g., scheduler, CRI) drift. **No crash-recovery replay exists.**

- **Publish failure after commit returns error though state persists.** Each post-commit publish returns the publish error to the caller (`core/metadata/images.go:172`, `plugins/services/containers/local.go:146`). The gRPC client sees `Unavailable/Internal` yet the object exists — retrying `Create` yields `ErrAlreadyExists` (`core/metadata/images.go:150`). Clients must handle create-then-publish-failed as success.

- **GC's async fire-and-forget can hide deletes.** `GarbageCollect` queues `namespacedEvent` inside sweep (`core/metadata/db.go:410`) then `publishEvents` in a goroutine (`core/metadata/db.go:441`). If the publisher is nil or topic unhandled, `publishEvents` just `Debug` logs (`core/metadata/db.go:363`). A watcher can hold a lease on deleted content without ever learning it was removed.

- **Overlay/backend out-of-sync after nested commit.** `snapshotter.Commit` (`core/metadata/snapshot.go:642`) commits to overlayfs/btrfs/zfs inside bolt's `Update`. Success + later bolt failure leaves a dangling backend snapshot logged at `core/metadata/snapshot.go:661` requiring manual removal. The reverse (bolt success but backend publish timeout) is not detected.

- **Single-writer flock contention + infinite timeout.** `plugins/metadata/plugin.go:44` sets `timeout.Set(boltOpenTimeout,0)` — `bolt.Open` blocks forever on `flock(2)`. A stuck previous writer holds daemon start; only a 10s warning goroutine (`plugins/metadata/plugin.go:172`) fires.

- **Shim publisher overflow.** `pkg/shim/publisher.go:59` `requeue chan 2048`; `processQueue` retries `5` times at `1s*count` (`pkg/shim/publisher.go:102`). Task exits/oom bursts >2048 or sustained daemon unavailability evict events silently (`pkg/shim/publisher.go:104`).

- **Lost-update without detection.** Two concurrent `Update` with overlapping `fieldpaths` serialize via bbolt. The later timestamp wins (`core/metadata/containers.go:270`). No `412 Precondition Failed` or vector clock surfaces the conflict to the caller.

- **Durability under `NoSync`.** If `cfg.NoSync=true` (`plugins/metadata/plugin.go:160`), an OS crash can lose the last committed transaction plus any in-memory events, with no WAL to recover.

## Future Considerations

- **Transactional outbox in `meta.db`.** Add an `outbox` bucket under `v1/<namespace>/outbox/<seq>` written inside the same `Update` that mutates the object, with a background forwarder that `View`-scans and publishes via `Exchange`/`Forward`, then deletes outbox entries in a second `Update`. Preserve short-tx discipline by batching outbox drain. This yields at-least-once without full event sourcing.

- **Sequence/revision column for images/containers.** Store `revision uint64` next to `createdat/updatedat` (`core/metadata/buckets.go:67`, `core/metadata/boltutil/helpers.go:110`) incremented on each `Update`, return it to clients, accept `ExpectedRevision` field-mask (`fieldpath.proto`-style) and fail with `FailedPrecondition` on mismatch — surfaces lost-update today hidden by last-write-wins. Existing `Target.Digest` conditional on delete (`core/metadata/images.go:316`) can graduate to generic revision.

- **GC outbox unify.** Merge GC's `namespacedEvent` queue (`core/metadata/db.go:395`) into the same outbox so snapshot/content deletions share the durable path, rather than special-casing async `publishEvents` (`core/metadata/db.go:348`).

- **Bolt V2 small-transaction belt.** Move `Snapshotter.Commit` out of the bolt `Update` (`core/metadata/snapshot.go:642`) and instead do bolt-reserve → backend commit → bolt-confirm with a pending marker bucket, reducing lock hold under overlay latency while keeping atomic confirmation via outbox/completion flag.

- **Bounded journal retention for watchers.** Keep last N `Envelope` per namespace in a capped `events` bucket (cursor-scannable) so late subscribers can replay from `SinceTimestamp` instead of requiring full List. Complements rather than replaces current-state store.

## Questions / Gaps

- **No evidence of durable event bucket or WAL.** Greps for `eventstore`, `Journal`, `outbox`, `WAL`, `checkpoint` (DB) found only CRIU/media-type hits or vendor `bbolt` internal constants; `internal/eventq/eventq.go:24` confirms in-memory buffering with discard. Search boundary: full workspace (`sources/containerd`), including `api/types/event.proto:27` envelope as payload-only.

- **No `ExpectedVersion`/revision integration tests.** `core/metadata/db_test.go` and `core/metadata/gc_test.go` cover GC labels and lease expiry, but no test asserts OCC rejection or outbox replay; confirms revision absence is intentional.

- **Content lease back-pressure under mixed policy.** How `shared` vs `isolated` interacts with `NoSync` crash replay and lease expiry GC for pulls that died mid-ingest — docs at `plugins/metadata/plugin.go:49` and `core/metadata/content.go:49` leave expiry handling implicit.

- **Sizing of shim requeue vs Exchange buffer.** `Exchange` uses unbounded `goevents.NewChannel(0)` (`core/events/exchange/exchange.go:132`) while shim caps at 2048 (`pkg/shim/publisher.go:80`) — no benchmark found for burst task events (OOM storm) starving watchers.

---

Generated by `Dimension 01.05: Atomic State Transition and Event Journal` against `containerd`.
