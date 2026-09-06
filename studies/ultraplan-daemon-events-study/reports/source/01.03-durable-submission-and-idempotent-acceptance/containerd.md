# Source Analysis: containerd

## 01.03 Durable Submission and Idempotent Acceptance

### Source Info

| Field | Value |
|-------|-------|
| Name | containerd |
| Path | `studies/ultraplan-daemon-events-study/sources/containerd` |
| Language / Stack | Go / gRPC + bbolt (go.etcd.io/bbolt) + filesystem (content blobs, snapshot storage) + ttrpc/shim v2 runtime |
| Analyzed | 2026-09-02 |

## Summary

containerd has **no generic idempotency-key / request-fingerprint layer**. Deduplication is conflated with resource identity: the *canonical key* IS the idempotency key. Container `ID` (validated via `pkg/identifiers`), image `Name`, snapshot `key` (`snapshotter/key`), content `digest` (`digest.Digest`), ingest `ref`, lease `ID`, and task `ContainerID` each act as the deduplication key. Submission becomes durable only after a **bbolt `Update` transaction commits** (`core/metadata/db.go:272`, `core/metadata/bolt.go:47`); gRPC responses are sent *after* commit, satisfying durable-before-acknowledgement for metadata objects. Content ingest is hybrid: resumable files under `root/ingest/<hash>` (`plugins/content/local/store.go:546`) plus a metadata bucket (`core/metadata/content.go:427`), with final commit atomically creating a `blob/<digest>` bucket. Task `Create`/`Start` are outliers: they deduplicate via shim runtime state (`plugins/services/tasks/local.go:270-276`) rather than bbolt, so acknowledgement can race with shim durability. No time-based idempotency-record retention exists; deduplication lives as long as the object lives, and is reclaimed by GC/leak-cleanup rather than TTL.

## Rating

**5 / 10 — Present but inconsistent, fragile under retry/crash**

**Rationale:** bbolt-backed `Create` paths for containers/images/snapshots/content are correctly ordered (validate → `db.Update` → commit → reply) and deterministically return `ErrAlreadyExists` on duplicate (`core/metadata/containers.go:148`, `core/metadata/images.go:145`, `core/metadata/snapshot.go:344`, `core/metadata/content.go:698`, `plugins/content/local/store.go:542`). Duplicate/conflict tests exist (`core/metadata/containers_test.go`, `core/metadata/content_test.go:129`, `plugins/content/local/store_test.go:218`, `core/mount/manager/manager_test.go:341`). However: (a) there is no cross-cutting idempotency token — retry with different payload under the same key is simply rejected without conflict diagnosis; (b) task creation has no transactional durability and `Start` is not idempotent; (c) content `Prepare` uses a two-phase commit that can leave a backend snapshot without a metadata entry if the second tx fails (`core/metadata/snapshot.go:508`); (d) event publish occurs *after* commit and can fail post-durability (`plugins/services/containers/local.go:139`, `core/metadata/snapshot.go:280`, `core/metadata/images.go:168`, `core/metadata/content.go:631`), leaving caller uncertain; (e) retention is implicit ( GC or explicit `Delete`), not a bounded idempotency window.

## Evidence Collected

Every entry MUST include a file path with line numbers. Format: `path/to/file.ts:NN`.

| Area | Evidence | File:Line |
|------|----------|-----------|
| Submit/start handlers — containers | `local.Create` wraps `Store.Create` in `withStoreUpdate` (bbolt `Update`) then publishes `/containers/create` only on success | `plugins/services/containers/local.go:122-151` |
| Submit — container store transaction | `Create` validates, calls `update(ctx, s.db, ...)` which opens `bkt.CreateBucket([]byte(container.ID))`; on `ErrBucketExists` maps to `errdefs.ErrAlreadyExists` before any reply | `core/metadata/containers.go:122-172` |
| Submit — images | `imageStore.Create` creates bucket under `getImagesBucket`; `ErrBucketExists` → `ErrAlreadyExists` inside same `update` tx; publish after | `core/metadata/images.go:124-178` |
| Submit — snapshots Prepare | `snapshotter.Prepare` delegates to `createSnapshot`; first tx checks `bkt.Bucket(key)` and `Bucket(target)` → `ErrAlreadyExists` + `addSnapshotLease`; second tx after backend `Prepare` creates bucket | `core/metadata/snapshot.go:274-518` |
| Submit — snapshots Commit | `Commit` creates `name` bucket inside tx; `ErrBucketExists` → `ErrAlreadyExists`; commits backend `Snapshotter.Commit` inside same tx | `core/metadata/snapshot.go:520-681` |
| Submit — content Writer (prepare) | `contentStore.Writer` checks `getBlobBucket(tx,ns,dgst) != nil` → `exists=true` (outside commit) then inside tx creates `ingest/<ref>` bucket; if shared it skips backend write; if digest already present returns `AlreadyExists` after tx | `core/metadata/content.go:373-497` |
| Submit — content Commit (durable accept) | `namespacedWriter.commit` creates `blob/<digest>` bucket; `ErrBucketExists` → `ErrAlreadyExists`; `Commit` outer `update` deletes ingest bucket and adds content lease atomically; publishes `/content/create` only after success | `core/metadata/content.go:582-642` , `core/metadata/content.go:651-718` |
| Submit — content local store | `store.writer` checks `os.Stat(blobPath(expected)) == nil` else `ErrAlreadyExists` (`plugins/content/local/store.go:536-544`); `Writer` holds per-ref `tryLock`/`unlock` ensuring single writer per ref | `plugins/content/local/store.go:475-499`, `plugins/content/local/store.go:531-629` |
| Submit — content start boundary | `store.Writer` requires non-empty `Ref` (`ErrInvalidArgument`), acquires lock, creates `ingest/<hash>/ref|data|startedat|updatedat` files transactionally | `plugins/content/local/store.go:484-606`, `plugins/content/local/writer.go:153` |
| Submit — tasks Create | `local.Create` resolves container, checks `rtime.Get(ctx, ContainerID)` — if found `ErrAlreadyExists`, else `rtime.Create`; no bbolt tx | `plugins/services/tasks/local.go:171-293` |
| Submit — tasks Start | `local.Start` fetches task and calls `p.Start(ctx)` each call; no idempotency guard — second Start fails with runtime error not `AlreadyExists` | `plugins/services/tasks/local.go:295-316` |
| Idempotency tables | bbolt buckets are the dedup store: `v1/<ns>/containers/<id>` (`core/metadata/buckets.go:242`), `v1/<ns>/images/<name>` (`core/metadata/buckets.go:229`), `v1/<ns>/snapshots/<snapshotter>/<key>` (`core/metadata/buckets.go:253`), `v1/<ns>/content/blob/<digest>` (`core/metadata/buckets.go:269`), `v1/<ns>/content/ingests/<ref>` (`core/metadata/buckets.go:289`), `v1/<ns>/leases/<id>` (`core/metadata/buckets.go:212`) | `core/metadata/buckets.go:146-317` |
| Canonical fingerprinting | Container ID validated via `identifiers.Validate` (`core/metadata/containers.go:324`); content via `digest.Digest.Validate()` and `blobPath` (`plugins/content/local/store.go:646`); snapshot key via `createKey(sid, ns, key)` with monotonic `NextSequence()` (`core/metadata/snapshot.go:360`) | `core/metadata/containers.go:323-326`, `plugins/content/local/store.go:646-652`, `core/metadata/snapshot.go:60-63` |
| Transactions around acceptance | `DB.Update` acquires `wlock.RLock` then `bbolt.Update` (`core/metadata/db.go:272-284`); `update` helper reuses tx from context if present (`core/metadata/bolt.go:47-55`) | `core/metadata/db.go:271-284`, `core/metadata/bolt.go:47-55` |
| Duplicate/conflict tests — containers | Validation and `AlreadyExists` asserted via `store.Create`/`Update`/`Delete` workflows; error unwrapping with `errors.Is(err, ErrAlreadyExists/ErrInvalidArgument)` | `core/metadata/containers_test.go:174-676` |
| Duplicate — content local | `TestContentWriter` writes blob, second `Commit` with same digest expects `ErrAlreadyExists` (`plugins/content/local/store_test.go:212-223`); `store_test.go:169` asserts second `Writer` with same ref fails | `plugins/content/local/store_test.go:143-224` |
| Duplicate — metadata content leases | `TestContentLeased` writes blob under lease, then second `Writer` with same digest under different lease expects `ErrAlreadyExists` | `core/metadata/content_test.go:97-142` |
| Duplicate — snapshots | `target snapshot %q: ErrAlreadyExists` and `snapshot %q: ErrAlreadyExists` paths inside first tx; backend `AlreadyExists` handling walks committed snapshots to deduplicate `target` | `core/metadata/snapshot.go:334-430` |
| Duplicate — mount manager | `TestActivateAlreadyExists` — second `Activate` with same mount name returns `ErrAlreadyExists` | `core/mount/manager/manager_test.go:341-363` |
| Client retry behavior | Client pullers and mount consumers treat `ErrAlreadyExists` as benign (`client/pull.go:309`, `client/container.go:457`, `core/unpack/unpacker.go:419`, `pkg/rootfs/apply.go:73`); `remote snapshotter` doc notes client calls `Stat` after `Prepare` returns `ErrAlreadyExists` | `client/pull.go:309`, `client/container.go:457`, `core/unpack/unpacker.go:419-471`, `docs/snapshotters/remote-snapshotter.md:118-126` |
| Leases (explicit retain) | `leaseManager.Create` creates lease bucket; duplicate ID → `ErrAlreadyExists` (`core/metadata/leases.go:67-78`); `addSnapshotLease`/`addContentLease`/`addIngestLease` gate on `leases.FromContext` | `core/metadata/leases.go:51-102`, `core/metadata/leases.go:337-456` |
| Content end-to-end test suite | `ContentSuite`, `ContentCrossNSSharedSuite`, `ContentCrossNSIsolatedSuite` exercise Writer/Commit/Abort across shared vs isolated policy | `core/metadata/content_test.go:84-95`, `plugins/content/local/store_test.go:96-103` |
| Proto API — no idempotency key | `CreateContainerRequest`, `CreateTaskRequest`, content `Write` proto have no `idempotency_key` / `request_id` field; dedup key is the resource field itself (`id`, `name`, `digest`, `ref`) | `api/services/containers/v1/containers.proto:147-149`, `api/services/tasks/v1/tasks.proto` (CreateTaskRequest ContainerID), `api/services/content/v1/content.proto:261` |

## Answers to Dimension Questions

**1. Can acknowledgement happen before durable acceptance?**
No for bbolt-backed objects: `plugins/services/containers/local.go:122-151` and `core/metadata/containers.go:142-167` commit via `db.Update` (which is `bbolt.Tx.Commit` under `wlock.RLock` at `core/metadata/db.go:272`) *before* populating `resp` and returning; callers receive gRPC reply only after `update` returns `nil`. Same ordering holds for images (`core/metadata/images.go:130-177`), snapshots (`core/metadata/snapshot.go:327-373` second `update` before return), content commit (`core/metadata/content.go:604-629`), and leases (`core/metadata/leases.go:67-100`). However acknowledgement *can* appear before durability for tasks: `plugins/services/tasks/local.go:277-293` returns after `rtime.Create` and `monitor.Monitor` succeed without a bbolt transaction; task state lives in the shim (filesystem + memory) and is only eventually re-synced via `v2Runtime.Tasks(..., true)` on daemon restart (`plugins/services/tasks/local.go:143-149`). Content ingest also has a hybrid case: the ingest data file (`plugins/content/local/store.go:608-629`) is written to disk before the metadata bucket is committed; a crash after file sync but before commit leaves an orphaned `ingest/<hash>` directory cleaned by `contentStore.garbageCollect` (`core/metadata/content.go:844-960`). Event publish (`publisher.Publish`) intentionally runs *after* commit via `DB.Publisher` check at `core/metadata/db.go:288-292` and can fail after durability (`plugins/services/containers/local.go:139` returns the publish error even though the container is already durable), forcing retry uncertainty.

**2. What happens if the same key is retried with identical input?**
Deterministic `AlreadyExists`. Container `Create` second call hits `bkt.CreateBucket` → `errbolt.ErrBucketExists` → `errdefs.ErrAlreadyExists` at `core/metadata/containers.go:149-152`; service maps via `errgrpc.ToGRPC` to `codes.AlreadyExists` (`plugins/services/containers/local.go:137`). Image `Create` at `core/metadata/images.go:145-150`, snapshot Prepare at `core/metadata/snapshot.go:344-350` or `458-463`, content `Writer` at `core/metadata/content.go:482` or `commit` at `core/metadata/content.go:698-700`, and content local `Writer` at `plugins/content/local/store.go:542` behave identically (`TestContentWriter:218-222`, `TestContentLeased:129-133`). The error is idempotent-safe: docs and callers explicitly handle it as success-path (`docs/snapshotters/remote-snapshotter.md:118-126`, `client/pull.go:309` ignores `AlreadyExists`, `core/unpack/unpacker.go:419`). Task `Create` also returns `AlreadyExists` if `rtime.Get` finds the task (`plugins/services/tasks/local.go:271-276`), but only within the lifetime of the shim.

**3. What happens if the same key is reused with different input?**
Still `AlreadyExists`, with **no diff or conditional update**. The store does not compare payloads; the bucket existence check is the sole gate. For containers, attempting `Create` with same `ID` but different `Spec`, `Image`, `SnapshotKey`, `Runtime` or `Labels` at `core/metadata/containers.go:148-153` returns the same `AlreadyExists` as an identical retry — callers cannot distinguish stale retry from conflicting intent without a subsequent `Get`. Content has a nuanced variant: if `WithDescriptor` digest differs, `Writer` at `core/metadata/content.go:398-411` checks `getBlobBucket` for the *new* digest, so a different blob under the same `ref` is allowed; but if the *digest* is the same and committed, any different `Size`/`MediaType` is still `AlreadyExists` (`core/metadata/content.go:698`). Snapshot `Prepare` with same `key` but different `parent` at `core/metadata/snapshot.go:344` also returns `AlreadyExists` and leaks no information about the existing parent; the `Commit` path at `core/metadata/content.go:651-675` validates `size`/`expected` only when a writer exists, otherwise requires exact match to `nw.desc` or returns `ErrFailedPrecondition`. The API provides `Update` with field masks for mutation, but `Create` remains unconditional overwrite-reject, so a client that intended to change fields must `Get` → `Update` rather than blind retry.

**4. How long is deduplication state retained?**
Indefinitely — until explicit `Delete` or GC. Containers, images, snapshots, and committed content blobs have no TTL; their dedup buckets persist as live objects (`core/metadata/buckets.go:75-114`). Retention ends only via `Delete` (`core/metadata/containers.go:287-321`, `core/metadata/images.go:282-345`, `core/metadata/snapshot.go:683-746`, `core/metadata/content.go:204-243`) or lease expiry + GC. Leased objects are pinned via `leases` bucket (`core/metadata/leases.go:186-249`); unleased but unreachable objects become GC roots only if they lose GC references (`core/metadata/gc.go:495-828`). Content ingests have a soft 24h `expireAt` set when not leased (`core/metadata/content.go:458-462`, `plugins/content/local/store.go:593-605`) and are eventually removed by `contentStore.garbageCollect` walking `ingestSeen` (`core/metadata/content.go:856-960`) and `store.WalkStatusRefs` (`plugins/content/local/store.go:364-394`). Snapshots left without metadata after a failed second tx are cleaned by `snapshotter.garbageCollect` walking backend and deleting unseen keys (`core/metadata/snapshot.go:862-941`). There is no bounded idempotency window table that auto-expires successful `Create` records — reuse of a resource name will conflict forever.

**5. Can a network/client failure accidentally start duplicate work?**
For metadata objects, no — duplicate work is prevented because the second attempt sees `AlreadyExists` rather than creating a second resource. The risk is instead **ambiguous success**: if the client times out after the bbolt commit but before the gRPC response arrives, it has no idempotency token to safely query the outcome without knowing the key; retry *will* return `AlreadyExists`, which is sufficient to infer success if the retry uses the same key/digest, but the client must implement that mapping itself (e.g., containerd clients treat `AlreadyExists` as benign at `client/pull.go:309`, `client/container.go:457`). For content ingest, a timeout after `Writer.Write` but before `Commit` leaves a resumable ingest directory (`plugins/content/local/store.go:546-580` supports `resumeStatus` via `io.CopyBuffer` over existing `data` at `plugins/content/local/store.go:501-529`); retry with same `ref` resumes from `Offset` (`plugins/content/local/store.go:571-577`), not duplicate. For tasks, yes — divergent: if `Create` succeeds in the shim but the response is lost, the client retry will receive `AlreadyExists` and can call `Get`/`Start` to proceed; however if the task was only half-created (e.g., `rtime.Create` succeeded but `monitor.Monitor` or `c.PID` at `plugins/services/tasks/local.go:282-288` failed) the shim may hold a stale task that must be `Delete`d before retry succeeds. More subtly, `snapshotter.createSnapshot` at `core/metadata/snapshot.go:374-515` creates backend state between two metadata transactions; a crash between them can orphan a backend snapshot (cleaned only after `Remove` or GC), and retry will hit the first-tx `AlreadyExists` path that adds a lease but does not re-prepare mounts, leaving the caller without mounts and needing `Stat`/`Mounts` to recover. No `at-most-once` queue semantics exist; all safety derives from the idempotency of the canonical key.

> Crash or disconnect after commit but before the start response reaches the client. Can retry return the original operation?
For metadata (containers/images/snapshots/content blobs): **yes** — the original operation is already durable in bbolt (or filesystem for content) and a retry with the same key/digest returns `AlreadyExists` which the client can interpret as “already accepted”; a follow-up `Get`/`Info`/`Stat` returns the committed object (`core/metadata/containers.go:55-79`, `core/metadata/content.go:72-92`, `core/metadata/snapshot.go:101-143`). For tasks: **partially** — retry of `Create` returns `AlreadyExists` and `Get` will return the task if the shim survived, but `Start` after a lost `Start` response requires re-`Start` (which may return `FailedPrecondition` if already running) and `Wait`/`PID` must be re-queried; there is no returned `operation` handle to poll.

## Architectural Decisions

* **bbolt as the single durable commit point** (`core/metadata/db.go:84-88`, `core/metadata/bolt.go:37-55`). All namespace-scoped metadata is written via `update(ctx, db, fn)` which reuses an outer transaction if present, ensuring that validation, bucket creation, timestamp/label writes, and lease attachment happen atomically. Decision trades write concurrency (single `wlock.RLock` + bbolt `RW` tx) for strong consistency and crash safety (bbolt `fsync` on commit). Evidence: `DB.Update` at `core/metadata/db.go:272` holds `wlock.RLock` for duration.

* **Resource identity == idempotency key**. No auxiliary `request_id` table; `container.ID`, `image.Name`, `digest`, `snapshot key`, `lease ID` serve as fingerprint. Keeps API surface minimal and aligns with OCI/container semantics where names/digests are content-addressed. Consequence: cross-object idempotency (e.g., “create container with digest X”) is not deduplicated unless name matches.

* **Content-addressed dedup for blobs, reference-addressed for ingests**. Committed blobs dedup by `digest` (`core/metadata/content.go:698` bucket `blob/<digest>`), ingests dedup by `ref` (`core/metadata/content.go:427` bucket `ingests/<ref>`). Shared vs isolated namespace policy (`DB.WithPolicyIsolated` at `core/metadata/db.go:61`) decides whether an `expected` digest can alias existing content across namespaces (`core/metadata/content.go:414-423`).

* **Two-phase snapshot creation with lease side-car**. First tx reserves metadata + lease (`addSnapshotLease` at `core/metadata/leases.go:337-364`), backend `Prepare` runs outside tx, second tx finalizes metadata. Lease ensures GC cannot collect snapshot between phases. Design isolates slow snapshot driver I/O from bbolt lock (lock released between phases) at cost of a window where metadata and backend can diverge.

* **Lease-based GC roots instead of refcounting**. All durability retention is modeled through `content`, `snapshot`, `ingest`, `image`, `container` `gc.Node` graph (`core/metadata/gc.go:36-54`) and lease buckets; unreferenced objects are removed by mark-and-sweep (`core/metadata/db.go:383-489`). Avoids explicit reference counts, leverages bbolt cursor scans.

* **Event publishing intentionally outside tx** (`core/metadata/db.go:288-298` returns `nil` publisher inside transaction; `plugins/services/containers/local.go:139` publishes after `withStoreUpdate`). Prevents tx bloat and deadlocks with event bus, but introduces post-commit publish failure mode.

## Notable Patterns

* **Bucket-existence-as-conflict** pattern repeated everywhere: `bkt.CreateBucket(key)` with `errbolt.ErrBucketExists` mapped to `errdefs.ErrAlreadyExists`. Observed in 8 locations: containers (`core/metadata/containers.go:149`), images (`core/metadata/images.go:146`), snapshots (`core/metadata/snapshot.go:345`, `458`), content blob (`core/metadata/content.go:698`), ingests (`core/metadata/content.go:427`), leases (`core/metadata/leases.go:76`), local content store (`plugins/content/local/store.go:542`), tasks service (`plugins/services/tasks/local.go:275`).

* **Resumable writer with hash replay**. `store.resumeStatus` at `plugins/content/local/store.go:501-529` reopens `data` file, re-hashes via `io.CopyBuffer` into `digester` to reconstruct `Offset`/`Status`; client can resume after crash by re-opening same `ref`. Metadata wrapper `contentStore.Status` at `core/metadata/content.go:306-330` resolves `ref → bref` (hashed ingest root) via bbolt.

* **Hashed ingest paths**. `ingestRoot(ref)` hashes `ref` with `digest.FromString(ref)` to keep paths constant length (`plugins/content/local/store.go:654-659`); `core/metadata/content.go:444` also uses `createKey(sid, ns, ref)` to generate unique `bref`. Prevents filesystem path traversal from user-controlled refs.

* **Shared-namespace optimization**. `isSharedContent` scans all namespaces for `labels.LabelSharedNamespace=="true"` plus existing blob (`core/metadata/content.go:758-776`); when `shared==true`, `Writer` avoids backend allocation and merely records `expected` digest (`core/metadata/content.go:464-467`), making cross-namespace ingest almost free.

* **Orphan cleanup via walk-and-diff GC**. Both snapshot (`core/metadata/snapshot.go:862-1003` `walkTree`/`pruneBranch`) and content (`core/metadata/content.go:844-960`) build `seen` maps from metadata, then delete backend objects not in `seen` after `GCStats` phase.

## Tradeoffs

* **Strong durability, limited concurrency.** bbolt gives ACID persistence with minimal ops overhead, but `Update` serializes writers behind a single file lock; high-throughput parallel `Create` (e.g., CRI pulling many images) contends on `wlock`. Alternative like PostgreSQL MVCC would scale but complicate embedded-daemon deployment.

* **Simplicity of identity-dedup vs. ergonomics.** Reusing resource names as idempotency keys avoids a second table and aligns with `digest` content addressing, but forces clients to generate unique IDs themselves (e.g., container `ID` must be UUID) — a collision on ID is indistinguishable from a retry. No `Idempotency-Key` header or `request_id` decoupled from resource name.

* **Two-phase snapshot commit reduces lock time but risks divergence.** Holding bbolt lock while calling `Snapshotter.Prepare` (potentially remounting or contacting remote snapshotter) would stall all metadata writes; splitting tx avoids that (`core/metadata/snapshot.go:367-383` outside tx). However crash between backend success and second tx leaves orphaned backend snapshot; recovery depends on GC rather than immediate consistency.

* **Post-commit event publish decouples concerns at cost of at-most-once notification.** Publishing inside tx would couple bbolt commit latency to event subscribers; outside keeps tx fast (`core/metadata/containers.go:171` return before publish). On publish failure the handler returns error to client despite durability — client sees failure but object exists, requiring `AlreadyExists` handling.

* **Filesystem + bbolt hybrid for content creates split durability.** Data path (`plugins/content/local/store.go:608` `OpenFile` + `Seek`) and metadata path (`core/metadata/content.go:427` ingest bucket) are not updated atomically; `garbageCollect` repairs, but a crash can leave either side orphaned.

* **Task runtime not transactional.** Keeping task state in shim improves performance and isolation from metadata DB, but sacrifices the uniform `ErrAlreadyExists` + durability guarantees enjoyed by other objects; operators lose the ability to `List` tasks after daemon loss until shim reconciliation.

## Failure Modes / Edge Cases

* **Lost response after commit.** If gRPC response is dropped after `bbolt.Commit` and before TCP ack, client retry with same key receives `AlreadyExists`; naïve clients that treat that as fatal abort lose idempotency. Correct handling is to treat `AlreadyExists` as success and `Get` the object — pattern proven in `client/pull.go:309` and `core/unpack/unpacker.go:419` but not enforced by API docs (`api/services/containers/v1/containers.proto:48-49` has no retry guidance).

* **Reuse of key with different payload silently rejected.** No `409 Conflict` body distinguishing “identical retry” vs “semantic conflict”; a second `Create` with different `Spec` is indistinguishable from a retry, so a bug that reuses IDs corrupts intent without warning.

* **Partial snapshot create.** If second `update` at `core/metadata/snapshot.go:445-505` fails (e.g., `parent` parent bucket missing due to concurrent `Remove`), `rerr` cleanup removes backend snapshot at `core/metadata/snapshot.go:509-512`, but if `Remove` itself fails, backend leaks until GC. Concurrent `Prepare`/`Commit` with same `key` can also race: first tx reserves bucket, second caller sees `AlreadyExists` and only adds lease, never getting mounts.

* **Content ingest lost-Ref after crash.** If daemon crashes after `store.writer` creates `ingest/<hash>/data` but before `contentStore.Writer` commits the `ingests/<ref>` bucket (`core/metadata/content.go:427-463`), restart sees file but no metadata entry; `ListStatuses` at `core/metadata/content.go:245-292` misses it, and GC via `WalkStatusRefs` (`plugins/content/local/store.go:364-394`) will `Abort` the file. Resuming with same `ref` will recreate rather than resume.

* **Event publish failure after durable commit.** `publisher.Publish` error causes handler to return error (`plugins/services/containers/local.go:147`, `core/metadata/snapshot.go:286`) despite durable acceptance; if the handler’s gRPC interceptor maps publish errors to retryable codes, a client retry will hit `AlreadyExists` and may incorrectly treat the whole flow as failed.

* **Digest collision / AlreadyExists on shared store.** `store.writer` at `plugins/content/local/store.go:542` checks local filesystem for `expected` digest existence; if two namespaces push same digest concurrently, second gets `AlreadyExists` even though its `ref` differs. For isolated policy, `contentStore.Writer` at `core/metadata/content.go:414-423` only checks shared availability via `isSharedContent` scan, which is O(namespaces) and racy if GC deletes blob between check and commit.

* **Lease expiration races.** Ingest `expireAt` at `core/metadata/content.go:458` is `24h`; a client that writes slowly over a lease but loses network may have ingest aborted by GC before commit, returning `NotFound` on `Status`/`Commit`. Task `State` fetch at `plugins/services/tasks/local.go:367-377` races with `Delete` and returns `NotFound`/`Unavailable` which `getTasksMetrics` swallows (`core/mount/manager` etc.), hiding failures.

## Future Considerations

* **Introduce request fingerprint table.** Add optional `X-Request-ID` / `Idempotency-Key` header handling in gRPC interceptors, persisting `{key, digest(input), response, expiry}` in a dedicated bbolt bucket (e.g., `v1/<ns>/idempotency/<key>`). Return stored response on `IsAlreadyExists` within TTL, and surface `409` with `existing` payload when digest mismatches — would disambiguate retry vs conflict for containers/images.

* **Make snapshot Prepare single-transaction or WAL.** Either hold metadata lock across backend `Prepare` with a timeout, or write a WAL entry (`prepare WAL` → backend → commit → delete WAL) replayed on restart to reclaim orphaned backend snapshots immediately rather than via periodic GC.

* **Unify task persistence.** Record task creation intent in bbolt (e.g., `v1/<ns>/tasks/<id>` with `CreatedAt`/`CreatedBy`) before `rtime.Create`, so `Create` becomes durable and `AlreadyExists` survives daemon restart; shim reconciliation can then re-attach.

* **Atomic content commit.** Bundle ingest file `fsync` and metadata commit under a single `Link`/`Rename` + bbolt batch tx, or use `bbolt` DB file on same filesystem with `Sync` before commit to avoid split-brain orphans; expose `Commit` with `Expected` digest validation even for shared path.

* **Bounded idempotency retention.** For resources that are intentionally short-lived (e.g., CRI sandbox containers), add `labelGCExpire` or explicit `Retention` on dedup buckets so retries beyond window succeed with fresh object rather than perpetual `AlreadyExists`; integrate with existing GC label handlers at `core/metadata/gc.go:223-399`.

* **Observability for duplicate vs conflict.** Emit metrics/logs distinguishing `AlreadyExists (identical)` vs `AlreadyExists (payload differs)` and map the latter to `gRPC AlreadyExists` with richer `status.Details` (via `errdefs` → `errgrpc`); add client SDK helper `IsIdempotentRetry(err)`.

## Questions / Gaps

* No evidence found of a central `Idempotency` or `Deduplication` service, table, or interceptor across the daemon. Search over `core/`, `plugins/services/`, `api/` for `idempot`, `dedupl`, `fingerprint`, `requestID` yielded only incidental mentions (mostly vendor or CRI `idempotent` verb comments at `vendor/k8s.io/cri-api/pkg/apis/runtime/v1/api.proto:35`).

* No explicit retention policy / TTL for idempotency records beyond implicit object lifetime and 24h ingest `expireAt`. Searched `core/metadata/buckets.go`, `core/metadata/db.go`, `core/metadata/gc.go` — no bucket with TTL other than `expireAt` for ingests.

* Prepare/Start boundary validation inspected: container `validateContainer` at `core/metadata/containers.go:323-354` checks `ID`, `Runtime.Name`, `Spec` non-nil, `SnapshotKey→Snapshotter` coherence, then `update` does existence check; snapshot `validateSnapshot` at `core/metadata/snapshot.go:851-859` only validates labels; task `Create` does no container spec re-validation beyond `getContainer` lookup (`plugins/services/tasks/local.go:172`). No cross-boundary revalidation (e.g., snapshot key still exists at task start) — task start only checks container exists.

* Duplicate/conflict tests exist but are per-resource, not cross-resource nor via gRPC e2e retry harness. No test found that simulates “client timeout after commit, retry with same key”.

* Client retry behavior is caller-side best-effort (treat `AlreadyExists` as success). No daemon-provided retry middleware or protobuf idempotency level is set (`api/services/containers/v1/containers.proto` has no `idempotency_level` option). Not clear if containerd intends gRPC `Idempotent` retry policy at infrastructure layer.

---

Generated by `01.03-durable-submission-and-idempotent-acceptance.md` against `containerd`.
