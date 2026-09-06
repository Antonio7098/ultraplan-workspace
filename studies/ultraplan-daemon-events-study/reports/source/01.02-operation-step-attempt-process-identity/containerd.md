# Source Analysis: containerd

## 01.02 Operation, Step, Attempt, and Process Identity

### Source Info

| Field | Value |
|-------|-------|
| Name | containerd |
| Path | `studies/ultraplan-daemon-events-study/sources/containerd` |
| Language / Stack | Go (daemon + shims), BoltDB metadata, gRPC/TTRPC, OCI runtime-spec |
| Analyzed | 2026-09-02 |

## Summary

containerd implements a two-tier durable hierarchy (Namespace → Container/Sandbox) persisted in BoltDB, with a separate ephemeral runtime hierarchy (Task → Exec Process → OS PID) managed by per-container shim OS processes. The durable intention is the `Container`/`Sandbox` record; checkpointable progress is at the container-metadata transaction and, optionally, a CRIU task checkpoint; there is no first-class operation/step/attempt abstraction. Retries do not get distinct identities—recreating with the same ID yields `AlreadyExists` and shim publisher retries reuse the same envelope. One logical operation (container lifecycle) intentionally spans multiple gRPC calls and survives shim/daemon restarts via bundle-on-disk and `bootstrap.json`. Events can be attributed via `namespace+topic+container_id+exec_id+pid+timestamp`, but lack a global operation or correlation/causation ID to disambiguate late/duplicate attempts.

## Rating

**5/10 — Present but inconsistent / fragile**

Rationale: Container/Sandbox and namespace identity are mature, validated, and persisted transactionally with tests. Task/Process identity and event routing are explicit and observable. However containerd has no durable operation or step ledger, no attempt identifier, and no idempotency key. Retry is caller-responsible with `ErrAlreadyExists` as the only guard; late arrivals (duplicate TaskExit/Delete) are explicitly documented as possible and delegated to callers (`core/runtime/v2/shim.go:606-611`). Tracing exists out-of-band via OTel but does not appear in the event envelope.

## Evidence Collected

Every entry MUST include a file path with line numbers. Format: `path/to/file.ts:NN`.

| Area | Evidence | File:Line |
|------|----------|-----------|
| Durable intent: Container model | `Container` struct defines durable intent: `ID`, `Labels`, `Image`, `Runtime`, `Spec`, `SnapshotKey/Snapshotter`, `SandboxID`, timestamps | `core/containers/containers.go:30-84` |
| Durable intent: Sandbox model | `Sandbox` struct with `ID`, `Runtime`, `Spec`, `Sandboxer`, labels/extensions | `core/sandbox/store.go:29-46` |
| Durable intent: Bundle on disk | `Bundle{ID, Path, Namespace}` with `Path = <state>/<ns>/<id>` and work dir symlink | `core/runtime/v2/bundle.go:121-129` |
| Metadata persistence | Container Create writes `CreatedAt/UpdatedAt` in single Bolt transaction; duplicate bucket yields `ErrAlreadyExists` | `core/metadata/containers.go:142-158` |
| Metadata persistence: Update | Field-mask Update with `UpdatedAt = Now().UTC()` inside Bolt transaction | `core/metadata/containers.go:174-280` |
| Namespace as root isolation | `WithNamespace` stores namespace in context + gRPC/TTRPC headers; `NamespaceRequired` validates with `identifiers.Validate` | `pkg/namespaces/context.go:38-43` and `pkg/namespaces/context.go:69-77` |
| Identifier validation | Container/Task/Sandbox IDs validated via `identifiers.Validate` (alphanum + `._-`, max 76) | `pkg/identifiers/validate.go:52-64` |
| Runtime hierarchy: Task interface | `Task` embeds `Process` + `PID()`, `Namespace()`, `Exec()`, `Pids()`, `Checkpoint()` | `core/runtime/task.go:63-86` |
| Runtime hierarchy: Process | `Process` with `ID()`, `State()`, `Kill/Start/Wait` ; `ExecProcess` adds `Delete` | `core/runtime/task.go:35-60` |
| Runtime hierarchy: State ownership | `State{Status, Pid, ExitStatus, ExitedAt, Stdin/Stdout/Stderr}` is per-process result | `core/runtime/task.go:119-134` |
| Process persistence proto | `Process{container_id, id, pid, status, exit_status, exited_at}` — exec `id` is task id for init | `api/types/task/task.proto:35-46` and `api/types/task/task.pb.go:98-113` |
| Parent/root: Sandbox linkage | `Container.SandboxID` optional immutable parent; `Runtime.CreateOpts.SandboxID` passed to shim | `core/containers/containers.go:80-83` and `core/runtime/runtime.go:54-55` |
| Parent/root: Bundle namespace | `NewBundle`/`LoadBundle` derive path `<root>/<namespace>/<id>` | `core/runtime/v2/bundle.go:52-60` |
| Parent/root: Lease context | `WithLease(ctx,lid)` + gRPC header propagation; `FromContext` retrieves lease ID | `core/leases/context.go:24-39` |
| Lease model | `Lease{ID, CreatedAt, Labels}`, resources `content/ingest/snapshot/image` | `core/leases/lease.go:42-53` and `core/metadata/leases.go:184-328` |
| Correlation: Event envelope | `Envelope{Timestamp, Namespace, Topic, Event: typeurl.Any}` — gRPC proto mirrors with `google.protobuf.Any` | `core/events/events.go:27-32` and `api/types/event.proto:27-33` |
| Correlation: Publisher contract | `Publisher.Publish(ctx, topic, event)` derives namespace from context, sets `time.Now().UTC()` | `core/events/exchange/exchange.go:81-102` |
| Correlation: Shim publisher | `RemoteEventsPublisher.Publish` marshals `Envelope{Timestamp, Namespace, Topic, Event}` and forwards via TTRPC with requeue on failure | `pkg/shim/publisher.go:127-152` |
| Causation: Topic taxonomy | Constants `/tasks/create`, `/tasks/start`, `/tasks/exit`, `/tasks/delete`, etc. and `GetTopic()` dispatch | `core/runtime/events.go:24-77` |
| Event payload attribution | `TaskExit{container_id, id, pid, exit_status, exited_at}` ; `TaskDelete` similar ; exec `id` disambiguates init vs exec | `api/events/task.proto:59-65` and `api/events/task.proto:42-50` |
| Event validation & filtering | `validateEnvelope` checks namespace+topic+timestamp; `Subscribe` with `filters.ParseAll` and adapter matching | `core/events/exchange/exchange.go:224-238` and `core/events/exchange/exchange.go:128-158` |
| Process handle: Shim instance | `ShimInstance{ID(), Namespace(), Bundle(), Client(), Endpoint(): (address, version)}` — shim is concrete OS process | `core/runtime/v2/shim.go:228-245` |
| Process handle: Shim process file | Shim bundle holds `bootstrap.json` with `address, protocol, version` for reconnection after restart | `core/runtime/v2/shim.go:294-331` |
| Multi-process span | `ShimManager.Start` either reuses sandbox shim or forks new shim binary; one container lifecycle spans multiple `Start/Kill/Wait/Delete/Checkpoint` calls | `core/runtime/v2/shim_manager.go:208-307` and `core/runtime/v2/task_manager.go:159-260` |
| Retry: no attempt ID | Publisher `item{ev, ctx, count}` with `maxRequeue=5` increments count but reuses same `Envelope` (no attempt ID) | `pkg/shim/publisher.go:40-43` and `pkg/shim/publisher.go:102-124` |
| Retry: AlreadyExists guard | `CreateBucket` returns `ErrAlreadyExists` on duplicate container ID; no CAS token | `core/metadata/containers.go:148-151` |
| Duplicate event caveat | Comment acknowledges duplicate exit/delete possible from `cleanupAfterDeadShim` + ttrpc-on-close; caller must handle | `core/runtime/v2/shim.go:605-611` |
| Checkpoint | `Task.Checkpoint(ctx, path, opts)` and `TaskCheckpointed{container_id, checkpoint}` event | `core/runtime/task.go:78-79` and `api/events/task.proto:90-93` |
| Shim recovery | `loadShim` + `restoreBootstrapParams` reload shim after daemon restart; `LoadExistingShims` on manager init | `core/runtime/v2/shim.go:79-142` and `core/runtime/v2/task_manager.go:108-110` |
| GC lease tracking | Lease-referenced resources (content/snapshot) GC via `gcContext` scan; no operation ledger | `core/metadata/db.go:383-489` |

## Answers to Dimension Questions

### 1. What is the durable user intention?

The durable intention is a **`Container`** (`core/containers/containers.go:30`) and, when used, a **`Sandbox`** (`core/sandbox/store.go:29`). Both are namespaced (`pkg/namespaces/context.go:38`) and persisted atomically in BoltDB (`core/metadata/containers.go:142-172`, `core/metadata/db.go:84-116`). Fields: user-chosen `ID` (validated `pkg/identifiers/validate.go:52`), `Runtime.Name`, immutable snapshotter/spec pointers, mutable labels/image/spec. A `Bundle` on filesystem (`core/runtime/v2/bundle.go:121-129`) mirrors the intention as `<state>/<namespace>/<id>` for the shim to consume. Leases (`core/leases/lease.go:42`) represent a weaker, GC-protecting intention for in-flight pulls/pushes but expire.

Content/images/snapshots are not themselves intentions—they are resources referenced by containers and protected by leases (`core/metadata/leases.go:337-364`).

Inferred intent (not durable): a desired task state (create→start→pause→exit→delete) is not stored as a separate durable operation record; it is implied by the existence of container metadata plus the current shim Task `State` (`core/runtime/task.go:119`).

### 2. What is the smallest checkpointable unit?

* **Metadata layer:** single container/sandbox create or field-mask update is checkpointed as one BoltDB write transaction (`core/metadata/containers.go:142-172`, `core/metadata/containers.go:174-280`). No multi-step transaction spans containers; `update`/`view` helpers lock via `wlock` (`core/metadata/db.go:271-284`).
* **Filesystem layer:** `Bundle` creation (`core/runtime/v2/bundle.go:47-118`) is checkpointed as mkdir + `work` symlink + spec write; `atomicDelete` (`core/runtime/v2/bundle.go:159-169`) is two-phase rename.
* **Runtime layer:** the smallest runtime checkpoint is a CRIU checkpoint (`core/runtime/task.go:78-79`, `core/runtime/v2/shim.go:838-848`) producing a directory + `TaskCheckpointed` event (`api/events/task.proto:90-93`). The runtime `State` itself (`core/runtime/task.go:119`) is not checkpointed durably—only held in the shim and queried via `State(ctx)` (`core/runtime/v2/shim.go:882-905`).

There is no operation log that checkpoints progress through create→start→wait as one logical operation.

### 3. Does each retry get a distinct identity?

**No.** containerd has no operation/attempt ledger.

* Container ID reuse is rejected, not versioned: `bkt.CreateBucket([]byte(container.ID))` returns `ErrAlreadyExists` (`core/metadata/containers.go:148-151`). No `attempt`, `generation`, `retryCount`, or CAS token is stored.
* Image pull/transfer progress (`core/transfer/transfer.go:171-178`) is callback-based; retries re-enter `Transfer` with the same lease ID (`core/transfer/local/transfer.go:164`).
* The only retry counter in-tree is `pkg/shim/publisher.go:40-43` `item.count` with `maxRequeue=5` (`pkg/shim/publisher.go:37`) which re-queues the *identical* `Envelope` (`pkg/shim/publisher.go:102-124`, `pkg/shim/publisher.go:127-152`). The `Envelope` has no `attempt_id`, `message_id`, or dedup key; `Timestamp` is set once at publish (`core/events/exchange/exchange.go:99`).
* Lease GC also has no attempt: `AddResource`/`DeleteResource` are idempotent puts/deletes (`core/metadata/leases.go:186-249`).

**Implication for the rubric question:** if attempt 2 succeeds after attempt 1 returns late, there is no ID to quarantine attempt 1. The `container_id`+`pid`+`exit_status` in `TaskExit` (`api/events/task.proto:59-65`) refer to the same logical container; the system relies on callers to handle duplicate `TaskExit`/`TaskDelete` (`core/runtime/v2/shim.go:605-611`).

### 4. Can one logical operation span multiple runtime calls or OS processes?

**Yes, by design.**

* Lifecycle: `Create` (`core/runtime/v2/task_manager.go:159` → `core/runtime/v2/shim.go:643`) → `Start` (`core/runtime/v2/shim.go:737`) → `Wait` (`core/runtime/v2/shim.go:820`) → `Delete` (`core/runtime/v2/task_manager.go:290`) are separate gRPC/TTRPC calls against the same `container_id` (which equals task ID).
* OS process boundary: each container/sandbox runs in its own shim OS process (`containerd-shim-runc-v2`) launched via `shimBinary(...).Start` (`core/runtime/v2/shim_manager.go:326-347`). The shim exposes the Task service via TTRPC/gRPC socket (`core/runtime/v2/shim.go:339-372`). `ShimInstance.Endpoint()` (`core/runtime/v2/shim.go:464-465`) returns `(address, version)` to reconnect.
* Survivability: `writeBootstrapParams`/`readBootstrapParams`/`restoreBootstrapParams` persist connection info in `bootstrap.json` (`core/runtime/v2/shim.go:294-331`, `core/runtime/v2/shim_manager.go:353-381`). On daemon restart, `LoadExistingShims` (`core/runtime/v2/task_manager.go:108`) and `loadShim` (`core/runtime/v2/shim.go:79`) re-attach to existing shim processes. `Bundle.Path` includes namespace (`core/runtime/v2/bundle.go:40-41`, `58-59`), so identity survives daemon process boundaries.
* Sandbox grouping: when `SandboxID` is set, multiple containers can share one shim process (`core/runtime/v2/shim_manager.go:212-290`, `core/runtime/v2/shim.go:257-268`), meaning one OS process backs multiple logical container operations.

### 5. Can events be unambiguously attributed to the right entity?

**Partially — to container/exec/pid/namespace, but not to operation/attempt.**

* Envelope attribution: `Envelope{Namespace, Topic, Timestamp, Event}` (`core/events/events.go:27-32`). `Namespace` is injected from context (`core/events/exchange/exchange.go:86-89`, `pkg/shim/publisher.go:128-131`); `Topic` is derived from event type (`core/runtime/events.go:51-77`). Filtering via `filters.ParseAll` on `Field()` (`core/events/events.go:36-62`, `core/events/exchange/exchange.go:155-157`).
* Payload attribution: task events carry `container_id` (always), `id` (exec ID; empty = init), `pid`, `exit_status`, `exited_at` (`api/events/task.proto:59-65`, `api/events/task.proto:42-50`, `api/types/task/task.proto:35-46`). Sandbox events carry `sandbox_id` (`api/events/sandbox.proto:26-34`). `Process` adds `Stdin/Stdout/Stderr/Terminal` (`api/types/task/task.pb.go:98-113`).
* Process handle attribution: host PID is returned via `Pid()` (`core/runtime/task.go:67`, `core/runtime/v2/shim.go:571-579`) and included in `Pids()` (`core/runtime/v2/shim.go:780-795`) and `TaskExit`’s `Pid` (`core/runtime/v2/shim.go:173-179`).
* **Gaps:**
  * No global `operation_id`, `correlation_id`, `causation_id`, or `attempt_id` in `Envelope` (`api/types/event.proto:27-33`, `core/events/events.go:27-32`).
  * No tracing span is serialized into the envelope; OTel is configured via `plugins.TracingProcessorPlugin` (`plugins/types.go:58`) but not carried as envelope field.
  * Duplicate suppression is explicitly not guaranteed: `cleanupAfterDeadShim` publishes `TaskExit`+`TaskDelete` even though the shim may also have, and the comment says “moby should handle the duplicate events” and “should not rely on only one exit event” (`core/runtime/v2/shim.go:605-611`, `core/runtime/v2/shim.go:596-604`).
  * A shim publisher `queue` retries with backoff but can evict after `maxRequeue` without NACK (`pkg/shim/publisher.go:102-108`), creating attribution gaps under partition.

Thus `namespace + container_id + exec_id + pid + topic + timestamp` unambiguously names the durable entity and process instance, but cannot quarantine a stale attempt of the same logical operation.

## Architectural Decisions

* **Namespace as mandatory shard key** — every mutating API requires `NamespaceRequired(ctx)` (`core/metadata/containers.go:56-58`, `pkg/namespaces/context.go:69-77`) and metadata buckets are nested under version → namespace (`core/metadata/buckets.go:241-249`). Evidence: gRPC/ttrpc header propagation `pkg/namespaces/grpc.go:32-44` and `pkg/namespaces/ttrpc.go`. Tradeoff: strong isolation without cross-namespace transactions; callers must thread context correctly.
* **Container/Sandbox as separate durable types with immutable ID + runtime** — `validateContainer` enforces `ID` validation and `Runtime.Name` non-empty plus `Spec` required (`core/metadata/containers.go:323-354`); `Update` forbids changing `Snapshotter/Runtime.Name` without fieldpaths (`core/metadata/containers.go:218-224`). Decision keeps execution intent stable while allowing mutable labels/spec.
* **Task as ephemeral, shim-backed runtime object** — `TaskManager.Create` creates a `Bundle` then `ShimManager.Start` forks/links a shim (`core/runtime/v2/task_manager.go:159-260`, `core/runtime/v2/shim_manager.go:208-307`); task state is held only in shim memory and queried via `State()` (`core/runtime/v2/shim.go:882-905`). Decision favors crash isolation (one process per container/sandbox) over in-daemon state.
* **Bootstrap-params file as reconnection checkpoint** — `bootstrap.json` persists `(address, protocol, version)` (`core/runtime/v2/shim.go:294-331`) and `restoreBootstrapParams` migrates legacy `address` file (`core/runtime/v2/shim_manager.go:353-381`). Decision enables daemon restarts without losing tasks.
* **Event system as typeurl-typed pubsub with topic derived from type** — `Publish` marshals via `typeurl.MarshalAny` (`core/events/exchange/exchange.go:94-102`), `GetTopic` switches on `*events.Task*` type (`core/runtime/events.go:51-77`), `Forward` validates envelope (`core/events/exchange/exchange.go:224-238`). Decision decouples publishers from subscribers but avoids envelope-level correlation IDs.
* **Leases for GC-protection instead of operation ledger** — `leases.Manager` with `WithLease(ctx, lid)` propagated via gRPC headers (`core/leases/context.go:24-30`, `core/metadata/leases.go:337-364`) protects content/snapshots from GC (`core/metadata/db.go:383-489`). Decision optimizes for pull/push streaming, not workflow orchestration.

## Notable Patterns

* **Dual-bucket metadata layout with boltutil helpers** — `writeContainer`/`readContainer` + `boltutil.WriteAny/ReadAny/WriteLabels` (`core/metadata/containers.go:412-458`, `core/metadata/containers.go:356-410`) exhibit consistent serialization pattern across containers, sandboxes, images, leases.
* **NSMap for runtime tasks** — `runtime.NSMap[ShimInstance]` keyed by `(namespace, id)` (`core/runtime/v2/shim_manager.go:163`, `core/runtime/v2/shim_manager.go:305`). Pattern mirrors namespace isolation in memory.
* **Downgrade-on-NotImplemented** — `TaskManager.Create` retries with `Downgrade()` on `ErrNotImplemented` for rolling upgrades (`core/runtime/v2/task_manager.go:231-250`, `core/runtime/v2/shim.go:247-261`).
* **Shim-side at-least-once publisher with bounded requeue** — `RemoteEventsPublisher.processQueue` with `maxRequeue=5` and `time.Sleep(1*count)` (`pkg/shim/publisher.go:102-124`) is a simple temporal retry without persistence.
* **Exchange subscriber filter as adapter** — `Exchange.Subscribe` wraps `goevents` Broadcaster with `Filter` matcher (`core/events/exchange/exchange.go:146-158`) and per-subscriber `Channel(0)` queue (`core/events/exchange/exchange.go:128-134`), consistent across metadata GC events (`core/metadata/db.go:343-368`).

## Tradeoffs

* **Durability vs. operation observability:** container/sandbox records are strongly durable (Bolt transaction + `dirty` flag for GC), but the operation that created/started them has no ledger. Operators can list containers after a crash (`core/metadata/containers.go:81-121`) but cannot enumerate in-flight operations or their attempts.
* **Process isolation vs. attribution complexity:** one shim process per sandbox/container isolates failures (a crash publishes `TaskExit` with synthetic exit 255, `core/runtime/v2/shim.go:170-172`) but creates duplicate-event risk (two publishers for same `container_id` may emit same `TaskExit`/`TaskDelete`).
* **At-least-once shim events vs. exactly-once needs:** requeue (`pkg/shim/publisher.go:117-124`) and `maxRequeue` eviction (`pkg/shim/publisher.go:104-108`) favor liveness over durability; no DLQ or persistent queue. Good for best-effort observability, insufficient for workflow engines needing quarantine.
* **Lease expiry via label vs. explicit timeout:** `WithExpiration` encodes `containerd.io/gc.expire = RFC3339` label (`core/leases/lease.go:94-102`); GC interprets labels externally. Cheap and decentralized but not strongly synchronized with operation timeouts.
* **Synchronous Bolt transactions vs. async shim calls:** daemon metadata is CP-consistent inside Bolt; shim calls are async RPCs with 5-10s dial timeouts (`core/runtime/v2/shim.go:65-77`, `pkg/shim/publisher.go:157`). A container may be durably created but its task may still be `CreatedStatus` until shim responds—window where operation spans processes.

## Failure Modes / Edge Cases

* **Late duplicate TaskExit/Delete:** `cleanupAfterDeadShim` publishes synthetically if shim dies (`core/runtime/v2/shim.go:144-187`), while the shim’s publisher may also publish. Comment explicitly warns callers to tolerate duplicates (`core/runtime/v2/shim.go:606-611` → refs `containerd/containerd#4769`). Without an attempt ID, a stale exit (e.g., pid recycled) could be misattributed by a naive consumer filtering only on `container_id`.
* **Namespace mismatch drops events:** `Publish` errors `failed publishing event: namespace is required` (`core/events/exchange/exchange.go:86-89`); `validateEnvelope` rejects missing namespace (`core/events/exchange/exchange.go:225-227`). A shim publishing with empty namespace would fail and be re-queued up to 5 times then dropped (`pkg/shim/publisher.go:102-108`).
* **Bundle leak on failed Create:** `TaskManager.Create` defers `bundle.Delete()` on error (`core/runtime/v2/task_manager.go:164-167`); `cleanupShimTask` (`core/runtime/v2/shim.go:200-222`) separately handles post-start failures. If both paths fail concurrently (e.g., daemon crash mid-create), bundle may remain until GC or next `LoadExistingShims`.
* **PID recycling not distinguished:** `Process{pid uint32}` (`api/types/task/task.proto:38`, `core/runtime/task.go:122`) is host PID; after shim restart the same `container_id` could get a different host PID. No boot-id or PID namespace identifier is included, so pid alone is not stable across attempts.
* **Requeue loss:** shim publisher’s in-memory `requeue` channel (`pkg/shim/publisher.go:73`) of size 2048 (`pkg/shim/publisher.go:36`) loses buffered events if shim crashes before forward succeeds (publisher `Close` just closes channel, `pkg/shim/publisher.go:93-99`).
* **Downgrade race:** concurrent container creates through upgraded daemon talking to old shim can trigger version downgrade (`core/runtime/v2/task_manager.go:237-250`); downgrade is per-shim instance (`core/runtime/v2/shim.go:468-475`) and not persisted beyond `bootstrap.json`.
* **Checkpoint path not transactional with metadata:** `Checkpoint` writes to a user-provided `path` (`core/runtime/v2/shim.go:838-844`) outside Bolt; failure leaves partial directory without metadata rollback.

## Future Considerations

* **Add envelope-level correlation/causation IDs** — extend `api/types/event.proto:27-33` `Envelope` with `operation_id`, `attempt_id`, `parent_envelope_id` or propagate OTel `trace_id/span_id` (already configured via `TracingProcessorPlugin` at `plugins/types.go:58`). Would answer the rubric’s “which IDs quarantine attempt 1?” without breaking `filters` (add new fieldpath handling in `core/events/events.go:36-62`).
* **Introduce a durable operation/step ledger** — e.g., `operation` bucket parallel to `containers` with `(operation_id, container_id, attempt, status, retries)` and FK to lease. Currently `core/metadata/db.go:84-116` tracks dirty but not operations; a ledger would enable at-least-once→exactly-once elevation for callers like Kubelet/BuildKit.
* **Make shim event publisher persistent or sequence-numbered** — add monotonic `seq` per `(namespace, container_id)` in `bootstrap.json` or Bolt, included in `TaskExit`/`TaskDelete`. Today `Envelope.Timestamp` (`core/events/events.go:28`) is wall-clock (`time.Now().UTC()` at `core/events/exchange/exchange.go:99`) and not monotonic.
* **Adopt CAS/version for container updates** — store `UpdatedAt` version or `revision` (already has `CreatedAt/UpdatedAt`) and require it on `Update` (`core/metadata/containers.go:174`) to detect stale retries.
* **Expose shim Pid+BootID correlation** — include shim’s `ShimPid`/`TaskPid` (`core/runtime/v2/shim.go:603-604` usage in runc service) and host boot ID in `ProcessInfo.Info` (`api/types/task/task.pb.go:217-228`) for stronger process-instance disambiguation.

## Questions / Gaps

* **No evidence of operation/step definitions** — searched `core/*`, `api/*`, `pkg/shim/*` for `Attempt`, `Operation`, `Step`, `Job` models; only `Transfer`’s `ProgressEvent` and lease lifetimes were found. Containerd does not model multi-step user workflows (expected for a runtime).
* **No causation metadata in events** — inspected `core/events/events.go:27`, `api/types/event.proto:27`, `pkg/shim/publisher.go:136-142`, `core/events/exchange/exchange.go:80-119`; no parent/correlation field was located. If a higher-level orchestrator injects it via `extensions` or labels, it is not validated or indexed.
* **External-process identifier scope unclear for shim-grouped containers** — when `SandboxID` sharing is used (`core/runtime/v2/shim_manager.go:212-252`), one shim PID multiplexes several containers; an event’s `container_id` disambiguates logically but `pid` in `TaskExit` refers to the container’s init PID, not the shim PID. Need integration test to confirm exec `id` scoping under grouped mode.
* **Durability of in-flight transfer as “step”** — `core/transfer` defines `Transferrer` but not a durably-tracked step; progress is callback (`core/transfer/transfer.go:155-165`). Whether lease + content ingest qualifies as checkpointable step requires follow-up with pull/push flow tracing.

---

Generated by `01.02-operation-step-attempt-process-identity` against `containerd`.
