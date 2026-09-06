# Source Analysis: containerd

## Dimension 01.06: Event Envelope, Ordering, Causation, and Replay

### Source Info

| Field | Value |
|-------|-------|
| Name | containerd |
| Path | `studies/ultraplan-daemon-events-study/sources/containerd` |
| Language / Stack | Go / gRPC + TTRPC + protobuf (google.protobuf.Any + typeurl), docker/go-events Broadcaster |
| Analyzed | 2026-09-02 |

## Summary

containerd implements a minimal, live-broadcast event envelope over in-memory fan-out, not a durable log. The envelope (`api/types/event.proto:27-33`, `core/events/events.go:27-32`) carries `timestamp`, `namespace`, `topic`, `event (Any)` with no sequence, cursor, ID, or causation fields. Publishing stamps `time.Now().UTC()` at the exchange (`core/events/exchange/exchange.go:99`) and validates `namespace`/`topic` (`core/events/exchange/exchange.go:224-237`), then writes to `docker/go-events` `Broadcaster` (`core/events/exchange/exchange.go:42-44`, `vendor/github.com/docker/go-events/broadcast.go:47-54`). Subscription is a live-only `Subscribe(ctx, filters...) -> (<-chan Envelope, <-chan error)` (`core/events/events.go:78-80`) multiplexed via `Broadcaster.Add(dst)` + `Queue`/`Channel` per subscriber (`core/events/exchange/exchange.go:128-198`). The gRPC/TTRPC surfaces (`api/services/events/v1/events.proto:27-49`, `plugins/services/events/service.go:101-120`, `plugins/services/events/ttrpc.go:36-42`, `client/events.go:80-125`, `core/events/proxy/remote_events.go:99-147`) expose only `SubscribeRequest { repeated string filters }` — no offset, sequence, or replay cursor. `internal/eventq/eventq.go:56-140` provides a short-time bounded `discardAfter` buffer for CRI use, but the core exchange has no persistence. Events are typed past-tense facts (`ContainerCreate`, `TaskExit`, etc. in `api/events/*.proto`) versioned only implicitly by protobuf evolution with `Any` type URLs. No causation/correlation metadata exists; consumers infer cause from `topic` + payload + `namespace`.

## Rating

**4 / 10 — Present but inconsistent, weakly documented, fragile for durable contracts.**

Rationale: Typed envelope and topic validation exist with tests, but there is no sequence/cursor, no durable ordering, no gap-free resume, no causation, and no schema-version contract. Operationally the system behaves as best-effort in-memory pub/sub: disconnect mid-stream demonstrably loses events with no detection mechanism.

## Evidence Collected

Every entry MUST include a file path with line numbers. Format: `path/to/file.ts:NN`.

| Area | Evidence | File:Line |
|------|----------|-----------|
| Event envelope — proto wire | `message Envelope { timestamp=1, namespace=2, topic=3, Any event=4 }` | `api/types/event.proto:27-33` |
| Event envelope — generated | `type Envelope struct {Timestamp *timestamppb.Timestamp; Namespace string; Topic string; Event *anypb.Any}` | `api/types/event.pb.go:40-49` |
| Event envelope — internal | `type Envelope struct { Timestamp time.Time; Namespace string; Topic string; Event typeurl.Any }` + `Field([]string)` adaptor | `core/events/events.go:27-62` |
| Envelope conversion | `toProto` / `fromProto` mapping `Timestamp<->protobuf`, `Event<->typeurl.MarshalProto` | `plugins/services/events/service.go:122-138` |
| Topic validation | `validateTopic: must start with '/' , len>1, components identifiers.Validate` | `core/events/exchange/exchange.go:201-222` |
| Envelope validation | `validateEnvelope: namespace valid, topic valid, Timestamp non-zero for Forward` | `core/events/exchange/exchange.go:224-237` |
| Publish path | `Publish(ctx, topic, event) { NamespaceRequired(ctx); validateTopic; MarshalAny; Timestamp=Now.UTC(); broadcaster.Write }` | `core/events/exchange/exchange.go:80-119` |
| Forward path | `Forward(ctx, envelope) { validateEnvelope; broadcaster.Write }` deferred logging | `core/events/exchange/exchange.go:55-75` |
| Broadcaster impl | `type Broadcaster { sinks []Sink; events chan Event; adds/removes }`, `Write(event) { select events<-event vs closed }` | `vendor/github.com/docker/go-events/broadcast.go:13-54` |
| Queue impl | `type Queue { dst Sink; events *list.List; cond }`, `Write` pushes + `cond.Signal`, `run` drains via `next()` | `vendor/github.com/docker/go-events/queue.go:13-111` |
| Subscriber API — core | `type Subscriber interface { Subscribe(ctx context.Context, filters ...string) (<-chan *Envelope, <-chan error) }` | `core/events/events.go:78-80` |
| Subscribe impl | Creates `Channel(0)`, `Queue(Channel)`, optional `NewFilter`, `broadcaster.Add(dst)`, goroutine `case ev:=<-channel.C: evch<-env`, cleanup `channel.Close(); queue.Close(); Remove` | `core/events/exchange/exchange.go:128-199` |
| gRPC service definition | `service Events { rpc Publish; rpc Forward; rpc Subscribe(SubscribeRequest) returns (stream Envelope) }` | `api/services/events/v1/events.proto:27-49` |
| SubscribeRequest | `message SubscribeRequest { repeated string filters=1; }` — no cursor/sequence | `api/services/events/v1/events.proto:60-62` |
| TTRPC mirror | Same `Forward` + `SubscribeRequest` definitions (TTRPC variant) | `api/services/ttrpc/events/v1/events.proto:27-36` (via generated `api/services/ttrpc/events/v1/events.pb.go:40-45`) |
| gRPC service handler | `Subscribe(req, srv) { events.Subscribe(ctx,filters); for select ev<-eventq: srv.Send(toProto(ev)) }` | `plugins/services/events/service.go:101-120` |
| TTRPC handler | `fromTProto` mirrors `fromProto` | `plugins/services/events/ttrpc.go:36-50` |
| Client Subscribe | `eventRemote.Subscribe { client.Subscribe(ctx, req) ; go { for ev,err:=session.Recv() } }` no replay param | `client/events.go:80-125` |
| Proxy Subscribe (reuse) | `grpcEventsProxy.Subscribe` and `ttrpcEventsProxy.Subscribe` — identical pattern, `io.EOF` handling | `core/events/proxy/remote_events.go:99-147` and `183-231` |
| Topic constants | Task topics `/tasks/create`, `/tasks/start`, `/tasks/exit` etc., fallback `/tasks/?` | `core/runtime/events.go:24-47` |
| Type→topic map | `GetTopic(e any) string { switch e.(type) case *TaskCreate: return TaskCreateEventTopic ...}` | `core/runtime/events.go:51-77` |
| Event payload examples | `ContainerCreate { id, image, Runtime{ name, Any options}}`, `ContainerDelete {id}`, `TaskExit {container_id, id, pid, exit_status, exited_at}` | `api/events/container.proto:27-46` and `api/events/task.proto:59-65` |
| Other event families | `ImageCreate/Update/Delete`, `SnapshotPrepare/Commit/Remove`, `NamespaceCreate`, `SandboxExit`, `ContentCreate/Delete` — all past-tense facts, no version | `api/events/image.proto:26-37`, `api/events/snapshot.proto:26-41`, `api/events/namespace.proto:26-38`, `api/events/sandbox.proto:25-37`, `api/events/content.proto:26-34` |
| Filter adaptation | `adapt(ev) -> Envelope.Field(fieldpath)` or stub `("",false)`, then `filters.ParseAll + filter.Match(adapt(gev))` | `core/events/exchange/exchange.go:240-248` and `147-157` |
| Field selection | `Envelope.Field(["namespace"])`, `["topic"]`, `["event", ...decoded.Field]` | `core/events/events.go:36-62` |
| Short-buffer queue (not core) | `EventQueue[T]{ events, subscriberC, shutdownC } New[T](discardAfter, discardFn) { discardQueue, discardTime, subscribers }` holds undelivered events for `discardAfter` then calls `discardFn` | `internal/eventq/eventq.go:24-163` |
| Queue tests | `TestMissedEvents` shows buffering 2 events before subscriber joins still delivered if within `discardAfter`; `TestSubscribersDifferentTime`; `TestDiscardedAfterTime` proves expiry → loss | `internal/eventq/eventq_test.go:53-103` |
| Exchange tests | `TestExchangeBasic`: two subscribers publish 3 `ContainerCreate`, both receive; `TestExchangeFilters`: topic/regex/field filtering | `core/events/exchange/exchange_test.go:34-271` |
| CTR consumer | `client.EventService().Subscribe(ctx, filters) -> for e:=<-eventsCh: Marshal(v) -> fmt.Println(timestamp, namespace, topic, json)` | `cmd/ctr/commands/events/events.go:44-77` |
| CRI consumer + backoff | `EventMonitor.Subscribe(subscriber filters) { ch, errCh = subscriber.Subscribe }`, handles only `TaskOOM/SandboxExit/Image* /TaskExit`, else `unsupported event`, retries via `backOffQueue` | `internal/cri/server/events/events.go:86-114` and `116-186` |
| Publisher interface | `Publisher { Publish(ctx, topic string, event Event) error }`, `Forwarder { Forward(ctx, *Envelope) }` — no sequence return | `core/events/events.go:68-76` |
| Absence: sequence/cursor | Grep `sequence\|cursor\|replay\|ordering` across `core/events` and `api/` yields zero envelope/cursor fields; only unrelated bbolt/CSI hits | (negative evidence — searched `core/events/*:1-248`, `api/types/event.proto:1-33`, `api/services/events/v1/events.proto:1-62`) |
| Absence: causation | Grep `causation\|correlation\|trace_id\|traceID` yields no envelope metadata; only `pkg/tracing` unrelated | (negative evidence — searched whole repo, no `Envelope.causation`) |

## Answers to Dimension Questions

### 1. What ordering is actually guaranteed?

Best-effort FIFO per publisher thread via synchronous in-memory channels, not a durable total order.

- `Exchange.Publish` and `Forward` call `broadcaster.Write(envelope)` (`core/events/exchange/exchange.go:74`, `118`) which sends on `Broadcaster.events chan Event` (`vendor/github.com/docker/go-events/broadcast.go:47-50`). The `run` loop (`broadcast.go:116-129`) iterates `for _, sink := range b.sinks { sink.Write(event) }` serially. Each subscriber's `Queue.Write` appends to `list.List` and `cond.Signal` (`queue.go:35-47`), drained in insertion order by `Queue.run -> next()` (`queue.go:66-87`, `93-111`). This preserves publication order per broadcaster instance as long as sinks do not error.
- There is **no** sequence number, per-subject monotonic counter, or term. No persistence means ordering is only as reliable as the process lifetime. Concurrent `Publish` calls race on channel send; no linearizable commit point. `Task` topics (`core/runtime/events.go:24-47`) imply logical ordering (e.g., `TaskCreate` before `TaskExit`) but it is not enforced by the bus — producers must publish in correct order.
- Slow consumers: `Exchange.Subscribe`'s forwarder does `select { case evch<-env: case <-ctx.Done(): break }` (`exchange.go:179-183`). The goroutine blocks on `evch` send until the consumer receives; upstream `Queue` will buffer unboundedly in its `list.List`, but if the consumer never reads, the whole pipeline stalls, not dropped — until context cancel. Conversely, broadcaster dropping on `ErrSinkClosed` logs at `Debug` (`queue.go:82-85`, `broadcast.go:126-127`) and removes sink.
- Operational reality: restarts wipe history; multiple containerd instances have independent buses; no cross-node ordering.

### 2. Can a reconnecting client resume without gaps or duplicates?

**No — disconnect guarantees gap, and duplicate semantics are undefined.**

- Live-only contract: `Subscribe(ctx, filters...) (<-chan *Envelope, <-chan error)` (`core/events/events.go:78-80`) has no `since`, `cursor`, `offset`, or `resumeToken`. `SubscribeRequest` carries only `repeated string filters` (`api/services/events/v1/events.proto:60-62`; TTRPC mirror `api/services/ttrpc/events/v1/events.pb.go:40-45`). Both gRPC and TTRPC handlers create a fresh `Channel(0)+Queue` and `broadcaster.Add(dst)` (`exchange.go:130-160`), then stream only events arriving after `Add`. Nothing is replayed.
- `client/events.go:88-125` and `core/events/proxy/remote_events.go:99-147` bridge the stream as `for ev,err:=session.Recv() { evq<-envelope }`. They detect `io.EOF`/`err != nil` and forward to `errq` — but there is no reconnection or catch-up logic.
- The disconnect mid-stream test: `TestExchangeBasic` (`exchange_test.go:34-113`) requires subscribers to be created *before* publishing; no test publishes while a subscriber is away and then re-subscribes. `internal/eventq/eventq_test.go:53-103` proves the only buffering is `EventQueue[T]` with a `discardAfter` window (e.g., `time.Second` or `3600*time.Second`). After `discardAfter`, `discardFn` drops events (`eventq.go:113-129`), and late subscribers miss them (contrast `TestMissedEvents` vs `TestDiscardedAfterTime:85-102`). Crucially, the core `Exchange` does **not** use `eventq` — it uses `docker/go-events` with zero retention (`NewQueue(channel)` no TTL). So the answer is stronger: even brief disconnects lose data.
- Deduplication: envelope has no `id`/`sequence`, so a reconnecting consumer cannot detect duplicates even if it attempted at-least-once replay via external storage. The consumer would need its own idempotency key derived from `topic`+`timestamp`+payload, which is unreliable (timestamp is `time.Now().UTC()` at publish, not monotonic — `exchange.go:99`).

### 3. Are events named as immutable past-tense facts?

**Mostly yes at the payload level; envelope/topic naming is present-tense path; enforcement is conventional, not contractual.**

- Payload types are past-tense: `ContainerCreate`, `ContainerUpdate`, `ContainerDelete` (`api/events/container.proto:27-46`); `TaskCreate`, `TaskStart`, `TaskExit`, `TaskOOM`, `TaskPaused`, `TaskResumed`, `TaskCheckpointed` (`api/events/task.proto:28-93`); `ImageCreate/Update/Delete` (`image.proto:26-37`); `SnapshotPrepare/Commit/Remove` (`snapshot.proto:26-41`); `NamespaceCreate/Update/Delete` (`namespace.proto:26-38`); `ContentCreate/Delete` (`content.proto:26-34`); `SandboxCreate/Start/Exit` (`sandbox.proto:25-37`). The fieldpath generation (`container.proto:26`) enables filter `event.id=="qwer"` (`exchange_test.go:206-210`), reinforcing fact-like queries.
- Topics use noun/verb path: `/tasks/create` (`core/runtime/events.go:26-46`), `/containers/create` (`exchange_test.go:134`). These are dispatched as literal strings validated by `validateTopic` (`exchange.go:201-222`). The mapping `GetTopic` (`runtime/events.go:51-77`) is a static switch.
- Immutability: events are pure data structs with no mutation semantics in transport; but nothing prevents a publisher from emitting contradictory facts (e.g., duplicate `TaskDelete` for same container). Versioning is absent (`typeurl.Any` type URL is the only discriminator). The system treats events as notifications, not append-only facts with retention.

### 4. Can the consumer tell why an event happened?

**Barely — only via `topic` + payload + `namespace`/`timestamp`; no causation/correlation chain.**

- Envelope provides `Timestamp` (`api/types/event.proto:29`), `Namespace` (`30`), `Topic` (`31`), `Event Any` (`32`). Example CTR prints exactly these four (`cmd/ctr/commands/events/events.go:67-71`). `Namespace` is derived from `namespaces.NamespaceRequired(ctx)` at publish (`exchange.go:86-88`), so consumers can scope by tenant/k8s (`internal/cri/server/events/events.go:135` filter `!=constants.K8sContainerdNamespace`). `Topic` hints intent (`/tasks/exit` vs `/tasks/create`).
- No `causationId`, `correlationId`, `parentId`, `actor`, `operationId`, or OpenTelemetry `traceId` in envelope. Tracing (`pkg/tracing/log.go:81` injects `trace_id` into logs) is not propagated in events. The `TaskExit` payload includes `exit_status, exited_at, pid` (`task.proto:59-65`) but not which client requested deletion or which controller triggered it. `ContainerCreate.Runtime.options Any` (`container.proto:31`) may carry opaque config, but not provenance.
- `Forward` (`exchange.go:55-75`) allows forwarding on behalf of another namespace/publisher, preserving `Timestamp` if set, but does not add `forwardedBy` metadata. Consumers must correlate externally (e.g., join with Kubernetes audit logs).
- Filterability helps diagnose: `topic=="/containers/create"`, `event.id`, regex `/containers/*` (`exchange_test.go:170-194`) can isolate a subject, but not causal chain.

### 5. How does a schema evolve without breaking old clients?

**Via protobuf3 + `google.protobuf.Any` + `typeurl` — additive-safe, but no explicit versioning policy.**

- Wire evolution: `Envelope.event` is `google.protobuf.Any` (`event.proto:32`) and `PublishRequest.event` is also `Any` (`events.proto:53`). Generated Go uses `anypb.Any` / `typeurl.Any` (`event.pb.go:48`, `events.pb.go:47`, `events.go:31`). Protobuf3 semantics preserve unknown fields (`protoimpl.UnknownFields` in generated structs), so adding new fields to a payload (e.g., `TaskExit.exited_at`) is forward compatible. Removing/renaming breaks.
- Type discrimination: the `Any.type_url` (e.g., `types.containerd.io/opencontainers/...`, actual URLs derived from `typeurl.MarshalAny`) functions as a schema identifier. Clients must register types via `typeurl.Register` (see `client/client.go:87-91`, `core/runtime/typeurl.go:31-35`) otherwise `MarshalAny` fails with `ErrNotFound`. However, `api/events` generated protos do not call `Register` in-package; registration is centralized in `typeurl` init via generated type URLs, so adding a new event type requires both proto definition and Go registration, but old binaries simply see unknown `type_url` and `UnmarshalAny` will error (`events.go:49-51`, `cri/server/events/events.go:92-94`).
- No schema registry, version header, or compatibility test harness. Options like `fieldpath_all = true` (`container.proto:26`) expose fields for filtering; adding a filterable field is safe, but changing field numbers/types is breaking. No `api_version` in envelope, so clients cannot negotiate. Duration: additive changes (new event types, new optional fields) are safe due to `Any`; rework of topic strings would break subscribers — topics validated strictly (`exchange.go:201-222`) but not versioned.

## Architectural Decisions

| Decision | Evidence | Consequence |
|----------|----------|-------------|
| In-memory Broadcast + Queue per subscriber instead of log | `Exchange.broadcaster *goevents.Broadcaster` (`exchange.go:37`), `NewBroadcaster()` (`exchange.go:42`), `NewChannel(0)+NewQueue(channel)` (`exchange.go:132-134`) | Simple, low-latency fan-out, but no durability/replay; restart = loss |
| Envelope = Timestamp + Namespace + Topic + Any | `api/types/event.proto:27-33`, `core/events/events.go:27-32`, `plugins/services/events/service.go:122-129` | Minimal, cross-language via protobuf, but no id/sequence/correlation |
| Topic as hierarchical path with strict validation | `validateTopic` (`exchange.go:201-222`), `Task*EventTopic = "/tasks/..."` (`core/runtime/events.go:24-46`) | Enables cheap string filtering, but no versioning/isolation |
| Filters are optional `filters.ParseAll` with `goevents.NewFilter` | `exchange.go:147-157`, `adapt` (`240-247`), tested in `exchange_test.go:115-271` | Flexible client-side predicate without server state |
| Timestamp assigned at publish, zero-check only on Forward | `envelope.Timestamp=time.Now().UTC()` (`exchange.go:99`), `validateEnvelope.Timestamp.IsZero` (`233-235`) | Wall-clock ordering only; clock skew directly perturbs order |
| Split Publish vs Forward | `Publisher.Publish` introspects namespace (`86-88`), `Forwarder.Forward` preserves envelope (`55-75`); gRPC exposes both (`events.proto:27-39`) | Enables shim forwarding (`cmd/containerd-shim-runc-v2/task/service.go:811-817`) and cross-namespace proxy without losing original timestamp |
| gRPC streaming Subscribe with TTRPC parity | `service.go:101-120`, `ttrpc.go:36-42`, `core/events/proxy/remote_events.go:99-231` | Unified API for CRI and external clients; but both lack resume |
| Protobuf Any + typeurl for extensibility | `api/types/event.pb.go:48`, `core/events/events.go:31`, `client/client.go:87-91` | Decouples envelope from payload schema; requires explicit registration |

## Notable Patterns

- **Live pub/sub with filter middleware**: `Broadcaster -> Queue -> Filter -> Channel` chain in `Exchange.Subscribe` (`exchange.go:130-160`) — classic docker `go-events` pipeline.
- **Adaptor/fieldpath for predicate pushdown**: `Envelope.Field(["event","id"])` delegates to decoded `typeurl.Any` via generated `Field` methods (`events.go:36-62` combined with `container_fieldpath.pb.go:6` generation via `cmd/protoc-gen-go-fieldpath/generator.go:84-89`).
- **Graceful subscriber teardown**: `closeAll` closure removes sink and closes `errq` (`exchange.go:137-142`), goroutine defers it and propagates `ctx.Err()` except `context.Canceled` (`189-195`). Mirrored in proxy/client (`client/events.go:98-122`).
- **Dual API (gRPC + TTRPC) with single backend**: Both services delegate to same `exchange.Exchange` (`plugins/services/events/service.go:66-73`, `ttrpc.go:32-34`) and share `toProto/fromProto` (`service.go:122-138`).
- **Generic buffered queue for CRI backoff**: `internal/eventq/eventq.go:56-140` is a generic `EventQueue[T]` with `discardAfter` TTL — used for transient replay before expiry; pattern not leveraged by core events.
- **Conventional topic registry**: `core/runtime/events.go:51-77` centralizes topic assignment; unknown types fall back to `/tasks/?`.

## Tradeoffs

- **Simplicity vs durability**: In-memory broadcast is easy to reason about and fast, but trades away exactly-once/resume semantics required for controllers.
- **Loose coupling via Any vs discoverability**: `Any` allows payloads to evolve without envelope changes, but consumers must know `type_url`s and handle `UnmarshalAny` errors; no self-describing schema.
- **Unbounded Queue vs bounded memory**: `list.List` in `queue.go:15,43` is unbounded — protects against burst loss but risks OOM if a subscriber stalls while publisher is hot.
- **Wall-clock timestamp vs sequence**: `time.Now().UTC()` (`exchange.go:99`) avoids coordinator, but inherits NTP skew and duplicates on equal timestamps; a monotonic sequence would enable idempotency/replay.
- **Filter pushdown vs server state**: Evaluating `filters.Match(adapt(gev))` inside the broadcaster (`exchange.go:155-157`) avoids per-client goroutines enumerating all events, but Regex/field filters are executed synchronously on the broadcast hot path.

## Failure Modes / Edge Cases

- **Lost events on disconnect/restart**: No cursor → any event published between `Remove` and next `Add` is irrevocably lost. Even with `eventq`, `discardTime` expiry (`eventq.go:113-129`) drops buffered events and calls `discardFn`. Tests explicitly demonstrate `TestDiscardedAfterTime` (`eventq_test.go:85-103`) where late subscriber misses `[1,2]`.
- **Clock skew / duplicate timestamps**: Two publishes in same nanosecond get distinct but equal logical moments; consumers sorting by `Timestamp` observe non-deterministic order and cannot deduplicate.
- **Slow consumer blocking**: `evch <- env` (`exchange.go:180`) blocks on slow consumer; upstream `Queue` grows unbounded. Alternative risk: if consumer context cancels mid-send, loop breaks without draining, and `errq` reports `context.Canceled` filtered out (`192`), causing silent termination vs explicit error in `client/events.go:115-119`.
- **Filter parse failure is immediate error**: `filters.ParseAll` error sends to `errq` then `closeAll` (`exchange.go:148-153`). Caller sees error channel but may not correlate to specific filter string.
- **Invalid envelope on Forward**: Missing `Namespace`/`Topic`/`Timestamp` returns `ErrInvalidArgument` via `errgrpc.ToGRPC` (`service.go:93-96`, `ttrpc.go:37-39`); callers must handle gRPC status, but CTR tool just warns and `continue` (`events.go:58-65`).
- **TypeURL mismatch**: If publisher and subscriber run different containerd builds (e.g., new `TaskCheckpointed` unknown to old subscriber), `typeurl.UnmarshalAny` fails (`events.go:48-50`, `cri/server/events/events.go:92-94`), currently downgraded to `log.Warn` or `continue`, not an error envelope — effectively silent loss of that fact.
- **Namespace isolation bypass**: Default `Subscribe` receives all namespaces unless filtered (`events.proto:44-45` comment); forgetting `namespace==default` filter leaks cross-tenant events.
- **Broadcaster shutdown race**: `Broadcaster.Close` waits on `b.closed` (`broadcast.go:93-100`); `Exchange` never exposes `Close`, but process exit races with `Write` may return `ErrSinkClosed` (`broadcast.go:51`), which publisher surfaces as generic error.
- **Shim forwarding loop risk**: `shim/publisher.Publish` (`task/service.go:811-817`) forwarding via same bus without loop detection could duplicate or amplify events if shims also subscribe.

## Future Considerations

- **Add monotone sequence + cursor**: Extend `Envelope` with `seq uint64` (per-exchange or per-namespace), persist head in bbolt/bolt, and expose `SubscribeRequest { cursor? uint64 }` with server-side replay from `seq+1` (bounded buffer or spill to disk). Would enable gap detection and idempotent consumers.
- **Introduce event IDs and deduplication**: `id` (UUIDv7) + `causationId`/`correlationId` (trace propagation) in envelope. Populate from `context` `trace.SpanContext` (`pkg/tracing/log.go:81`) so CRI (`events/events.go:134`) can tie `ImageCreate` to image pull operation.
- **Persisted topic log for critical domains**: At least for `tasks`/`containers` state transitions, project a write-ahead log (BoltDB or SQLite) with TTL and snapshot, allowing late-joining controllers (CRI `EventMonitor` `backOff` `132-186`) to reconstruct missed `TaskExit` without exponential backoff guesswork.
- **Schema versioning policy**: Add `envelope.schema_version` or standardize type URLs with API version component (e.g., `containerd.events.v1.TaskExit`) and CI compat checks (compare `api/events/*.proto` descriptors). Document additive vs breaking changes.
- **Bounded backpressure with metrics**: Export `queue.go` depth, `broadcaster` drop counters, and `exchange.go:179` block duration histogram; switch `evch` to buffered + drop policy (`NewChannel(100)`) with explicit `ErrDropped` envelope instead of unbounded memory.
- **Unify buffering**: Promote `internal/eventq.EventQueue` TTL pattern (`eventq.go:56-99`) to core Exchange as optional `SubscribeWithHistory(discardAfter)` or adopt a single abstraction to avoid two queue implementations.

## Questions / Gaps

- **No evidence of sequence/cursor/replay tests or config**: Searched `core/events/*`, `api/types/event.proto`, `api/services/events/v1/events.proto`, `vendor/docker/go-events/*` — no cursor field. `TODO(stevvooe)` (`exchange.go:172`) hints incompleteness.
- **Ordering guarantee undocumented**: No doc states FIFO vs best-effort. `exchange_test.go` does not test concurrent publish ordering or cross-namespace interleaving.
- **Schema compatibility contract not stated**: No `CONTRIBUTING.md` or `docs/` reference to event evolution rules; `api/events/doc.go:17-18` is placeholder. Topic naming convention not formalized outside `validateTopic`.
- **Causation story missing**: Could `namespace` + `topic` + `container_id` suffice for all auditing needs? No design doc ties events to Kubernetes/OTel trace. Checked `internal/cri/server/events/events.go:98-112` — only extracts `ContainerID`/`SandboxID`/`Name`.
- **Operational safeguards unverified**: No metrics, no alert on dropped events (`queue.go:82` is `Debug`), no chaos test for disconnect-reconnect. The `EventMonitor` backoff (`events.go:144-151`) is per-`id` string — does it conflate different event types for same id?
- **Open question: Should `Forward` preserve original `Timestamp` or overwrite?** Current `validateEnvelope` (`233`) requires non-zero but otherwise honors caller value (`ttrpc.go:44-50`). Clock skew for forwarded shim events (`task/service.go:699-700` `TaskOOM`) is client-dependent.

---

Generated by `01.06-event-envelope-ordering-causation-and-replay` against `containerd`.
