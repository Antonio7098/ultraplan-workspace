# Source Analysis: containerd

## Leases, Heartbeats, Fencing, and Stale-Worker Rejection

### Source Info

| Field | Value |
|-------|-------|
| Name | containerd |
| Path | `studies/ultraplan-daemon-events-study/sources/containerd` |
| Language / Stack | Go / bbolt (single-node daemon + shim processes, gRPC/TTRPC) |
| Analyzed | 2026-09-02 |

## Summary

containerd has no execution-ownership lease, heartbeat, or fencing subsystem in the dimension's sense. Its `Lease` (`core/leases/lease.go:42`) is a garbage-collection root, not a worker-ownership token. Ownership of long-lived work is split: durable intent lives as Bolt transactions for `Container` (`core/metadata/containers.go:142`) and filesystem `Bundle` / `bootstrap.json` (`core/runtime/v2/bundle.go:121`, `core/runtime/v2/shim.go:294`) with re-attachment on restart, while ephemeral execution lives in an out-of-process `shim` (`core/runtime/v2/shim.go:228`) that forwards `TaskExit` events via `exchange.Exchange` (`core/events/exchange/exchange.go:81`) and `pkg/shim/publisher.go:102` with at-most-once requeue. There is no heartbeat loop, no renewal, no generation/epoch/CAS on any mutation, and no staleness check before authoritative writes. Lease expiry is a wall-clock `containerd.io/gc.expire` label (`core/leases/lease.go:99`) evaluated lazily during `gcContext.scanRoots` (`core/metadata/gc.go:571`); it prevents future collection of leased content, not a stale worker's transaction. The documented contract for the counterfactual is explicit: `cleanupAfterDeadShim` notes "moby should handle duplicate events" and "should not rely on only one exit" (`core/runtime/v2/shim.go:605-611`). Pausing an old worker until its lease is stolen (or its shim dies) and then letting it finish does not corrupt via a checked fence — it double-publishes or last-writer-wins, and callers are responsible for idempotence.

## Rating

**3/10 — Absent for execution fencing; present only as GC protection (ad-hoc, unsafe for multi-worker orchestration).**

Lease/GC protection is clear, tested, and observable, but the dimension's required properties — claimed/renewed ownership lease, heartbeat cadence, fencing generation/CAS checked on product-state and artifact mutations, takeover with stale-completion rejection, and PID-rotation-safe identity — are all absent. The `Lease` struct has no revision, `Manager` has no `Renew` (`core/leases/lease.go:32`), `Update` has no expected-value predicate, and task `PID uint32` (`core/runtime/v2/shim.go:173`) recycles without a nonce. The system instead multiplexes safety into Bolt transactions + single-writer + shim restart reattachment, which is coherent for a single-daemon container supervisor but fails the "pause-old-worker-until-stolen" test as a correctness proof would.

## Evidence Collected

Every entry MUST include a file path with line numbers. Format: `path/to/file.ts:NN`.

| Area | Evidence | File:Line |
|------|----------|-----------|
| Lease schema | `type Lease struct { ID string; CreatedAt time.Time; Labels map[string]string }` — no `Generation`, `Epoch`, `Revision`, or expiry field; expiry is a label | `core/leases/lease.go:42-47` |
| Lease manager interface | `Manager { Create, Delete, List, AddResource, DeleteResource, ListResources }` — no `Renew`, `Heartbeat`, `KeepAlive`, `CAS` | `core/leases/lease.go:32-39` |
| Lease gRPC / proto | `message Lease { string id; Timestamp created_at; map labels }` — same minimal fields | `api/services/leases/v1/leases.proto:51-57` |
| Lease creation | `CreateBucket(id)` inside `update()` txn; `ErrBucketExists → ErrAlreadyExists`; `MarshalBinary(time.Now().UTC())` → `bucketKeyCreatedAt` | `core/metadata/leases.go:67-96` |
| Expiration label writer | `WithExpiration(d) Labels["containerd.io/gc.expire"] = time.Now().Add(d).Format(RFC3339)` — wall-clock, client-computed | `core/leases/lease.go:93-102` |
| Expiration reader (lease) | `expThreshold:=time.Now()` at scan start; `time.Parse(RFC3339, labelGCExpire)`; `if expThreshold.After(exp) { skip lease root }` | `core/metadata/gc.go:501-581` |
| Expiration reader (image) | `isExpiredImage` same RFC3339 parse; invalid value logged and treated as non-expired | `core/metadata/gc.go:1136-1150` |
| Expiration (ingest) | `writeExpireAt(now+24h)` when `leased==false`; `readExpireAt` checked in `scanRoots` with `expThreshold.After` | `core/metadata/content.go:458-462`, `core/metadata/gc.go:708-721` |
| Label-based GC handler | `labelGCExpire` constant `[]byte("containerd.io/gc.expire")` — expected RFC3339, not monotonic clock | `core/metadata/gc.go:89-95` |
| Flat lease variant | `labelGCFlat` distinguishes `resourceContentFlat/resourceSnapshotFlat` but still not a fencing token | `core/metadata/gc.go:97-102`, `core/metadata/gc.go:595-612` |
| GC mark & sweep | `DB.GarbageCollect` → `wlock.Lock` → `startGCContext` → `getMarked` → `gc.Tricolor` / `ConcurrentMark` → `scanAll` with `c.remove` | `core/metadata/db.go:383-489`, `pkg/gc/gc.go:64-194` |
| Lease context propagation | `WithLease` stores `leaseKey` in context + `metadata.Pairs("containerd-lease", id)` for gRPC; `FromContext` / `fromGRPCHeader` | `core/leases/context.go:23-40`, `core/leases/grpc.go:30-57` |
| Client helper (no heartbeat) | `WithLease(ctx, opts…)` defaults to `WithRandomID + WithExpiration(24h)`; creates lease, stores in ctx, returns `Delete` closer — no ticker, no renewal | `client/lease.go:27-53` |
| Lease ID mint | `WithRandomID` → `t.Nanosecond + base64(rand 3 bytes)` — random, not monotonic fencing nonce | `core/leases/id.go:27-34` |
| Snapshot/image lease helpers | `addSnapshotLease / addContentLease / addIngestLease` gate on `leases.FromContext(ctx)` and `getBucket(...Leases, lid)` existence, not liveness | `core/metadata/leases.go:337-509` |
| Container create (dedup) | `CreateBucket(id)` inside `Update` → `ErrAlreadyExists` as the only idempotency gate | `core/metadata/containers.go:148-154` |
| Container update (no CAS) | `view readContainer → field-mask merge → writeContainer` inside `update()` txn; `UpdatedAt:=Now().UTC()`; no `if updated.Revision != expected` check; `CreatedAt` restored | `core/metadata/containers.go:194-279` |
| Container delete (dirty) | `DeleteBucket(id)` + `dirty.Add(1)` — no generation check | `core/metadata/containers.go:304-320` |
| DB transaction & locking | `Transactor { View, Update }`, `view()/update()` delegate to `boltutil.Transaction(ctx)` if present; `DB.wlock sync.RWMutex` blocks writers during GC mark, not a distributed lease | `core/metadata/bolt.go:29-55`, `core/metadata/db.go:89-108`, `core/metadata/db.go:271-284` |
| Content ingest Writer | `Writer()` does `createIngestBucket` + `addIngestLease`; if unleased writes `expireAt +24h`; shared-content fast path via `addContentLease`; `commit` does `createBlobBucket → ErrAlreadyExists` | `core/metadata/content.go:373-497`, `core/metadata/content.go:696-718` |
| Shim / Task identity | `ShimInstance { ID(), Namespace(), Bundle(), Client(), Delete(), Endpoint() }`; `shim { bundle, client any, address, version }`; `Endpoint()` → `(address, version)` from `bootstrap.json` | `core/runtime/v2/shim.go:228-477` |
| Reattachment (no generation) | `loadShim` → `restoreBootstrapParams(bundle.Path)` → `makeConnection(address, protocol)`; `ShimManager.Start` writes `bootstrap.json` then `loadShim`; no epoch handshake | `core/runtime/v2/shim.go:79-142`, `core/runtime/v2/shim_manager.go:353-382` |
| PID-based task identity | `shimTask.PID() → task.Connect → TaskPid uint32`; `TaskExit{ ContainerID, ID, Pid uint32, ExitStatus, ExitedAt}` | `core/runtime/v2/shim.go:571-580`, `core/runtime/v2/shim.go:173-180` |
| Shim liveness & takeover | `ShimManager.startShim` registers `onClose → cleanupAfterDeadShim(WithoutCancel, id, shims, events, binary)` + `shims.Delete(id)`; `grpcDialContext` watches `connectivity.Ready/Idle/Shutdown` | `core/runtime/v2/shim_manager.go:333-342`, `core/runtime/v2/shim.go:377-427` |
| Duplicate-event contract | `cleanupAfterDeadShim` publishes `TaskExit`+`TaskDelete` with fallback `Status 255`; `shimTask.delete` comment "If shimErr == nil removeTask else allow callback to re-publish — moby should handle duplicate events" + TODO #4769 | `core/runtime/v2/shim.go:144-188`, `core/runtime/v2/shim.go:596-631` |
| Duplicate-event tolerance in publisher | `RemoteEventsPublisher` reuses identical `Envelope` with `maxRequeue=5, Sleep(1*count)` on failure | `pkg/shim/publisher.go:40-127` (traced via `core/runtime/v2/shim.go` publisher path) |
| Event envelope (no fence) | `Envelope{ Timestamp time.Time; Namespace string; Topic string; Event typeurl.Any }` — timestamp is wall-clock at publish (`time.Now()` in exchange), no `generation`/`attempt`/`stamp` | `core/events/events.go:27-32`, `core/events/exchange/exchange.go:81-99` |
| Lease tests (no stale-writer) | `TestLeases` / `TestLeasesList` verify `Create/List/Delete` + filter; `TestLeaseResource` verifies `AddResource/ListResources/DeleteResource` for 4 types; no heartbeat, no GC-expiry fencing, no CAS test | `core/metadata/leases_test.go:32-418`, `core/leases/lease_test.go:27-87` |
| Integration lease test | `TestLeaseResources` verifies that `AddResource(snapshot)` keeps snapshot alive after `Image.Delete`, and `Delete→SynchronousDelete` GC removes it — GC liveness, not fencing | `integration/client/lease_test.go:31-142` |
| Negative evidence: heartbeat | `grep -R heartbeat\|Heartbeat` across `sources/containerd/core` (excluding vendor) returns 0 hits | `sources/containerd` (negative evidence) |
| Negative evidence: fencing CAS | `grep -R generation\|epoch\|revision\|compareAndSwap\|WithRevision\|CAS` on non-vendor `core/` returns only `dbVersion` migration, no per-object revision | `sources/containerd` (negative evidence) |

## Answers to Dimension Questions

### 1. Can two workers believe they own the same work?

**Yes — there is no single-owner claim.** Leases are additive GC references, not exclusive locks. Any client can `WithLease(ctx, id) → AddResource(lease, {type: snapshots/overlayfs, id: <key>})` (`core/metadata/leases.go:186`, `core/metadata/leases.go:337`) without acquiring ownership; multiple leases can reference the same `content/snapshot/image/ingest` concurrently (`core/metadata/gc.go:594-658` emits all non-expired lease roots). Container creation is the only mutual-exclusion point (`CreateBucket` → `ErrAlreadyExists` at `core/metadata/containers.go:148`) and it deduplicates intent, not ongoing execution ownership. Task execution enters via `ShimManager.Start → NSMap.Add` (`core/runtime/v2/shim_manager.go:285`) — a process-local map keyed by `id`, not a distributed lease. `ShimManager` never validates a fencing token before `task.Start/Create/Kill` (`core/runtime/v2/shim.go:643-754`); the shim itself is the worker. Two containerd daemons on different hosts would happily `loadShim` the same bundle path if the underlying storage were shared — which it is not (bbolt is local). Within one daemon, two goroutines with the same lease ID can both hold `WithLease(ctx, id)` contexts and both `Update` the same container; the last `update()` txn wins (`core/metadata/containers.go:269`).

### 2. Does lease expiry alone prevent stale writes?

**No.** Expiry (`core/leases/lease.go:99`) is evaluated only during `gcContext.scanRoots` (`core/metadata/gc.go:571`) and `isExpiredImage` (`core/metadata/gc.go:1136`) when building the reachable set for the next `GC`. A stale worker that retains its `context.WithValue(leaseKey, id)` can still open a Bolt `Update` transaction and mutate containers, snapshots, or content — `containerStore.Update` (`core/metadata/containers.go:194`), `snapshotter` methods, and `contentStore.Writer.commit` (`core/metadata/content.go:651`) never check whether the lease in the context is still listed as a `ResourceLease` root or whether label `containerd.io/gc.expire` is in the future. The content path is the closest to "expiry-guarded": an unleased ingest gets `expireAt = now+24h` (`core/metadata/content.go:458`) and `scanRoots` silently skips expired `ingest` roots (`core/metadata/gc.go:716`), and `addContentLease` (`core/metadata/leases.go:386`) fails with `ErrNotFound` if the lease bucket is gone — but that error is surfaced only to callers that explicitly call the leased path; the non-leased content path proceeds. There is no `LeaseRevoked` error on subsequent Bolt writes.

Clock dependence is fragile: expiry uses `time.Now().UTC().Format(RFC3339)` at lease creation (`core/leases/lease.go:99`) and `Parse(RFC3339, exp)` at GC (`core/metadata/gc.go:572`); invalid labels are logged and ignored (`core/metadata/gc.go:574`), i.e., fail-open to "non-expired lease." GC runs under `wlock` (`core/metadata/db.go:384`) but its `expThreshold` is captured once at `scanRoots` entry (`core/metadata/gc.go:501`), so a writer that starts before GC and commits after is not revoked.

### 3. Is fencing checked by product-state and artifact mutations, not only event writes?

**No fencing is checked anywhere.** The dimension asks for conditional mutation on product state and on artifact stores; containerd has neither:

* **Product state (Bolt):** `containerStore.Update` (`core/metadata/containers.go:194`), `image` stores, and `snapshotter` metadata writes all do read-modify-write inside a single `Update` txn with no expected-generation predicate. `writeContainer` (`core/metadata/containers.go:412`) unconditionally overwrites `Labels/Spec/Extensions/Runtime.Options/SnapshotKey` + `UpdatedAt`. The only conflict signal is at `Create` time.

* **Artifacts (content/snapshots):** `contentStore.Writer` → `commit` (`core/metadata/content.go:651`) uses `createBlobBucket` (`core/metadata/content.go:696`) which fails with `ErrAlreadyExists` if that digest already exists — a content-addressed dedup, not a worker fence. `snapshotter.Prepare/Commit` similarly operate on `BucketKeyObjectSnapshots/<snapshottner>/<key>` without a generation token.

* **Events:** `events.Exchange.Publish` (`core/events/exchange/exchange.go`) and `shimTask.delete → events.Publish(TaskExit)` (`core/runtime/v2/shim.go:173`) emit `Envelope{Timestamp, Namespace, Topic, Event}` (`core/events/events.go:27`) with no `generation/stamp/epoch`. Late events are not rejected by the event bus; they are forwarded to all subscribers and de-duplicated (or not) by the consumer.

### 4. How is process identity distinguished from reusable PIDs?

**It is not distinguished with a fencing nonce.** The shim process identity is `(Namespace, Bundle.Path, ShimInstance.ID)` (`core/runtime/v2/shim.go:228`, `core/runtime/v2/bundle.go:121`) plus the out-of-process PID obtained via `task.Connect → TaskPid uint32` (`core/runtime/v2/shim.go:571`). `TaskExit` and `TaskDelete` carry the same raw `Pid uint32` (`core/runtime/v2/shim.go:173`, `api/events/task.proto:59` via `pkg/shim/publisher.go`). There is no `boot_id`, `creation_generation`, or `stamp` equivalent to Temporal's `ActivityAttemptStamp`. A killed shim's PID can be recycled by the OS; a late `Wait` reply from that PID is still accepted as `TaskExit.Pid = oldPid` (`core/runtime/v2/shim.go:820`). `ShimManager` stores shims in `runtime.NSMap[ShimInstance]` (`core/runtime/v2/shim_manager.go:163`) indexed by container `id`, not by `(id, generation)`. Reattachment uses `bootstrap.json { Version, Protocol, Address }` (`core/runtime/v2/shim.go:294`) to reconnect to the same Unix socket address after daemon restart, not to prove it is the same OS incarnation.

### 5. What happens to a late completion from a superseded attempt?

**It is published as a duplicate `TaskExit` and the caller is expected to tolerate it.** The comment at `core/runtime/v2/shim.go:605-611` is authoritative:

> "NOTE: If the shim has been killed and ttrpc connection has been closed, the shimErr will not be nil. For this case, the event subscriber, like moby/moby, might have received the exit or delete events. Just in case, we should allow ttrpc-callback-on-close to send the exit and delete events again... TODO: It's hard to guarantee that the event is unique and sent only once. The moby/moby should not rely on that assumption that there is only one exit event."

Concretely, a shim that is load-shedded or whose `onClose` fired (`core/runtime/v2/shim_manager.go:333-341`) publishes `TaskExit{ExitStatus:255, ExitedAt:now}` via `cleanupAfterDeadShim` (`core/runtime/v2/shim.go:170-179`). If the original `shimTask.delete → task.Delete` (`core/runtime/v2/shim.go:582`) eventually returns, it also calls `removeTask` and publishes (or swallows, depending on `shimErr`). The two publishes race; both carry the same `ContainerID/ID` and the same raw `Pid` (or `255`), but no `attempt_id` to disambiguate. `RemoteEventsPublisher` retries the same envelope up to 5 times (`pkg/shim/publisher.go:102`). The runtime itself never rejects the second publish; the containerd API surface delegates filtering to callers. For Bolt-backed state, a late `Commit` on content (`core/metadata/content.go:582`) after the "new owner" already committed the same digest returns `ErrAlreadyExists` (content-addressed); for snapshots/containers it blindly overwrites (last-writer-wins).

## Architectural Decisions

* **GC leases, not ownership leases** (`core/leases/lease.go:42`, `core/metadata/gc.go:89`): containerd chose to implement `Lease` as a label-gated reference list whose only evaluation point is the periodic `GarbageCollect` mark (`core/metadata/db.go:383`). This keeps the hot path (Bolt `Update`) uncontended and makes "holding something alive across a daemon pull/unpack" a first-class operation (`client/lease.go:27`), but it means a daemon operator cannot say "worker A held lease L for task T and now worker B holds it." There is no transaction between `Lease.Create` and `Lease.AddResource` that guarantees the resource was not collected in between — GC holds `wlock` only during mark (`core/metadata/db.go:384`).

* **Single-node bbolt + shim process delegation** (`core/metadata/db.go:84`, `core/runtime/v2/shim.go:228`): By hosting authoritative state in a single-process BoltDB and delegating long-lived `Task` lifetime to an OS shim re-parented on daemon restart (`core/runtime/v2/shim_manager.go:353`, `core/runtime/v2/shim.go:79`), containerd avoids distributed coordination entirely. The tradeoff is that the daemon cannot lose and regain ownership of a task without local restart; there is no cross-host fencing to reason about. The `wlock` (`core/metadata/db.go:89`) and per-store `sync.RWMutex` (`core/metadata/content.go:46`, `core/metadata/snapshot.go:48`) are the only mutual exclusion.

* **Wall-clock label expiry** (`core/leases/lease.go:99`, `core/metadata/gc.go:571`, `core/metadata/content.go:458`): Using an RFC3339 label instead of a server-authoritative `expires_at` column keeps lease management stateless and debuggable via `ctr leases` (`cmd/ctr/commands/leases/leases.go`), at the cost of clock-skew vulnerability and the need for periodic GC polling to enforce expiry. The ingest `expireAt` at `core/metadata/content.go:458` is a `MarshalBinary(time.Time)` stored alongside the ingest bucket, but it is also checked only at `scanRoots`.

* **Duplicate-tolerant events, not exactly-once** (`core/events/events.go:27`, `core/runtime/v2/shim.go:605`, `pkg/shim/publisher.go:40`): containerd explicitly documents that `TaskExit` may be delivered twice; idempotence is pushed to consumers (moby/kubelet). This is coherent for a container runtime where the caller can query `Task.State` (`core/runtime/v2/shim.go:882`) reflexively, but is the inverse of a fencing design where the second publish would be rejected via `generation mismatch`.

## Notable Patterns

* **Label-ref callback graph for GC** (`core/metadata/gc.go:223-446`): `startGCContext` builds `labelHandlers` for `containerd.io/gc.ref.*`, `gc.bref.*`, `gc.cond.*`, `gc.expire` and sorts them by `bytes.Compare` (`core/metadata/gc.go:437`), then `scanRoots` streams `gc.Node` roots through a channel and resolves back-references via `c.references` → `gc.Tricolor`. This is the central, well-tested correctness mechanism for the lease feature; it is also the only place lease expiry is enforced.

* **`WithLease` context + header propagation** (`core/leases/context.go:23`, `core/leases/grpc.go:30`, `client/lease.go:50`): The canonical way to associate work with a lease is to put `WithLease(ctx, id)` in the context before `Writer()`/`Prepare()`; server-side helpers read it via `FromContext` falling back to `metadata.FromIncomingContext` (`core/leases/grpc.go:45`). This is a voluntary, caller-driven association, not an ownership check.

* **`bootstrap.json` reattachment** (`core/runtime/v2/shim.go:294`, `core/runtime/v2/shim_manager.go:353`): The durable handle for a task is its bundle directory and its `bootstrap.json` (`{Version, Protocol, Address}`); `loadShim` + `makeConnection` with `AnonReconnectDialer` can reattach to a still-running shim after daemon crash, preserving `Task` lifetime across daemon restarts. No lease renewal is required for that survival.

* **`ttrpc.WithOnClose` + `grpcDialContext` state watcher** (`core/runtime/v2/shim.go:339`, `core/runtime/v2/shim.go:377`): `ShimManager.startShim` registers `cleanupAfterDeadShim` on ttrpc close and polls `connectivity.Ready → Idle/Shutdown` for gRPC shims — this is the closest thing to a heartbeat: a liveness callback, but synchronous and not periodic, and without propagation of a generation token.

## Tradeoffs

| Decision | Benefit | Cost |
|----------|---------|------|
| GC lease via label + periodic `Tricolor` (`core/metadata/gc.go:64`, `core/metadata/gc.go:501`) | Zero contention on writes; leases survive daemon crash without heartbeat; simple to inspect | Expiry enforced only at GC cadence; wall-clock skew can extend or shorten real lifetime; stale writers not revoked between GC cycles |
| `WithExpiration` client-computed wall-clock (`core/leases/lease.go:99`) | Stateless, debuggable, no server clock to trust | Clock-skewed client can produce far-future or already-expired leases; server logs and ignores invalid RFC3339 (`core/metadata/gc.go:574`) fail-open |
| Bolt `Update` per RPC, no revision (`core/metadata/bolt.go:47`, `core/metadata/containers.go:194`) | Simple transactions, no ABA handling; `wlock` only blocks GC, not readers | Last-writer-wins on `Update`; no way for a new owner to reject a stale write with `expected Generation` |
| PID-as-identity (`core/runtime/v2/shim.go:571`, `api/events/task.proto:59`) | Native OS identity, trivial payload | Recyclable, not fencing; late or spoofed `TaskExit` with reused PID is indistinguishable after shim socket is reused |
| Duplicate-tolerant events (`core/runtime/v2/shim.go:605`) | Minimal correctness obligation on runtime; robust to shim-socket races | Consumers must implement their own `Seen(TaskExit)` dedup; no journal ordering beyond `Envelope.Timestamp` (`core/events/events.go:28`) |

## Failure Modes / Edge Cases

* **Pause-old-worker-until-stolen counterfactual — decidedly unsafe.** Pause a shim goroutine after it has done meaningful work but before it publishes `TaskExit`. Let `ShimManager` detect `ttrpc.ErrClosed` via `onClose`, run `cleanupAfterDeadShim` publishing `ExitStatus 255` (`core/runtime/v2/shim.go:170`), and have a "new" task with the same `id` be `Create`d (which succeeds only after the first `DeleteBucket` at `core/metadata/containers.go:310`/`core/runtime/v2/shim_manager.go:475` — otherwise `ErrAlreadyExists`). Resume the old worker: if its `task.Delete` succeeds, it will publish a second `TaskExit` with the real exit code, potentially overwriting the `255` seen by controllers; `core/runtime/v2/shim.go:612-614` explicitly notes the shimErr path that allows the duplicate. If the old worker also committed a Bolt mutation (e.g., `contentStore.commit` at `core/metadata/content.go:696` via `createBlobBucket`), it will win or fail only on digest collision, not on lease ownership.

* **GC-while-writing race.** `DB.GarbageCollect` captures roots under `wlock.Lock` and `db.View` (`core/metadata/db.go:394`), but lease expiry (`core/metadata/gc.go:501`) and ingest `expireAt` (`core/metadata/gc.go:716`) are evaluated from that snapshot. A writer that has already entered `Update` (`core/metadata/db.go:272` acquires `wlock.RLock`) before GC's `wlock.Lock` can still `addContentLease` for a content blob whose lease just became unreachable but not yet swept — GC will then sweep the newly referenced blob only if the `marked` set was already built.

* **Unleased ingest expiration is advisory.** `writeExpireAt(now+24h)` (`core/metadata/content.go:458`) with `removeIngestLease` on delete (`core/metadata/content.go:620`) means an ingest created without a lease but abandoned before `Commit` stays live for up to 24 h (validated by the `expThreshold.After` check at `core/metadata/gc.go:716`). A stale ingest writer that resumes after 25 h will find its backing ingest bucket already deleted by `GC.remove` (`core/metadata/gc.go:1088`) but its `cs.Store.Writer` handle (`core/metadata/content.go:475` via `plugins/content/local`) may still hold a filesystem temporary file.

* **Invalid or skewed expiry never fences.** An ingest or lease label with unparsable RFC3339 is logged at `Info` and treated as non-expired (`core/metadata/gc.go:574`, `core/metadata/gc.go:1143`). A crashed client that writes `labels["containerd.io/gc.expire"]="bad"` permanently shields that resource from collection.

* **ORDS content path can outlive its metadata.** `contentStore.garbageCollect` (`core/metadata/content.go:845`) scans `blob`/`ingests` buckets for `contentSeen/ingestSeen`, then walks the backend `Store` to decide deletion. A late `Commit` that writes under `db.Update` (`core/metadata/content.go:604`) after `contentSeen` was snapshotted but before `Store.Delete` runs will be observed as a backend delete of a just-committed blob, dependent on filesystem/backend ordering.

* **Shim socket reuse / address race.** `grpcDialContext` probes `net.DialTimeout("unix", address, 10s)` (`core/runtime/v2/shim.go:394`) to decide if the socket is alive; a quickly restarted shim that reused the same `Address` (e.g., after `sandbox_shim` path in `core/runtime/v2/shim_manager.go:242`) would pass the dial check but may be a new OS incarnation with the same `id`. No epoch distinguishes them.

## Future Considerations

* **Least-invasive hardening if containerd remains single-daemon:** Add an optional `Revision` field (monotonic `uint64`) on `Container` backed by `bucketKeyRevision` and enforce `UpdateIfMatch(revision)` as a Bolt CAS (compare inside the same `Update` txn that currently does `readContainer` at `core/metadata/containers.go:205`). Expose it as `expected_revision` on `containers.Update` and on `Task` mutations. This mirrors BuildKit's unused `Generation` (`api/services/control/control.proto:215`) but makes it checked. Lease automation would stay unchanged; the fence moves to the mutating RPC, not the lease.

* **Lease-heartbeat for cross-host future:** If containerd is ever run over shared storage, promote `WithExpiration` (`core/leases/lease.go:93`) to a server-authoritative `Renew(ctx, Lease, TTL)` that bumps `Labels["containerd.io/gc.expire"]` inside an `Update` txn and fails if `CreatedAt` does not match (compare-and-swap on label map). Run a client-side ticker at `TTL/3` inside `client.WithLease` (currently at `client/lease.go:27` the closer just deletes) and a server-side "shim-must-renew" check before `task.TaskService` mutations (`core/runtime/v2/shim.go:643` onward). Document that wall-clock labels are advisory unless accompanied by periodic renewal.

* **Stale-PID fencing via shim nonce:** Mint a per-`Start` random `shim_generation` stored in `bootstrap.json` (`core/runtime/v2/shim.go:294`) and propagated into each `TaskExit` (`core/runtime/v2/shim.go:173`); require readers to match `(ContainerID, generation)` instead of `Pid`. This avoids boot-id dependence while still defeating PID recycling and the rapid-restart socket-reuse race (`core/runtime/v2/shim_manager.go:333`).

* **Duplicate-event idempotence key:** Include a server-generated `TaskExit.Sequence` or `EventEnvelope.id` monotonic per `id` (beyond `Envelope.Timestamp` at `core/events/events.go:28`) so consumers can `dedupe(TaskExit.ContainerID, Sequence)` pithily instead of relying on "moby should handle" prose (`core/runtime/v2/shim.go:610`).

* **Observability:** Emit `containerd_gc_lease_expiry_seconds_remaining` and `containerd_gc_invalid_expiry_total` from `scanRoots` (`core/metadata/gc.go:571`), and `containerd_task_duplicate_exit_total` from both branches of `shimTask.delete` (`core/runtime/v2/shim.go:612`).

## Questions / Gaps

| # | Question | Why it matters | Where to look |
|---|----------|----------------|---------------|
| 1 | Is delegation of TaskExit deduplication to consumers ("moby should handle duplicates" at `core/runtime/v2/shim.go:610`) a deliberate, supported contract or a known gap (#4769) intended to be closed with a fencing token? | Determines whether containerd should remain without stale-worker rejection or add a generation/CAS. | `core/runtime/v2/shim.go:606-612` + GitHub issue #4769 reference |
| 2 | Should `isExpiredImage` fail-closed on parse error rather than fail-open (`log Info + return false` at `core/metadata/gc.go:1143`)? | Fail-open leaks storage when a bad label is written; fail-closed would surface bugs faster. | `core/metadata/gc.go:1136-1150` vs `core/metadata/gc.go:572-576` (lease path also fail-open) |
| 3 | Is the single-node BoltDB explicitly out of scope for distributed fencing, or is multi-instance shared storage (e.g., overlay over shared Bolt via SAN) contemplated, which would force heartbeats/fencing? | Defines whether the 1-3 rating is by design or a gap. | `core/metadata/db.go:84` (`schemaVersion/Transactor` comment), `core/runtime/v2/shim_manager.go:163` (`NSMap`) |
| 4 | Why does `client.WithLease` default to `WithExpiration(24h)` (`client/lease.go:41`) while `core/metadata/content.go:458` uses `24h` for ingest and `core/metadata/gc.go` has no default lease TTL — are these durations coordinated or independent? | Uncoordinated TTLs risk a transfer lease expiring before ingests do, orphaning content. | `client/lease.go:39-42`, `core/metadata/content.go:458`, `core/leases/lease.go:93` |

---
Generated by `dimensions/01.04-leases-heartbeats-fencing-stale-worker-rejection.md` against `containerd`.
