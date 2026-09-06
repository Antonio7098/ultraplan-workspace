# Source Analysis: containerd

## Event Delivery, Backpressure, and Retention Tiers

### Source Info

| Field | Value |
|-------|-------|
| Name | containerd |
| Path | `studies/ultraplan-daemon-events-study/sources/containerd` |
| Language / Stack | Go (daemon + shims, gRPC + ttrpc, BoltDB metadata, CRI plugin) |
| Analyzed | 2026-09-03 |

## Summary

containerd treats core daemon events as an ephemeral, in-memory broadcast with no persistence tier. The central `Exchange` (`core/events/exchange/exchange.go`) fans out `Envelope`s synchronously through `docker/go-events` `Broadcaster` → per-subscriber unbounded `Queue` → unbuffered `Channel(0)` → per-subscription forwarding goroutine → gRPC `srv.Send`. There is no WAL, no replay, no retention/compaction/pruning, no batching, and no drop/coalesce/sample policy on this path: a slow subscriber grows its own unbounded `list.List` queue instead of stalling publishers.

Two peripheral tiers do have explicit bounds. The shim→daemon forwarder (`pkg/shim/publisher.go`) has a bounded retry channel (`queueSize=2048`) with eviction after `maxRequeue=5` attempts (drop + log). The CRI `GetContainerEvents` path (`internal/eventq/eventq.go` + `internal/cri/server/service.go`) has a time-bounded hold (`discardAfter=5m`), a per-subscriber buffer (`cap 100`), a discard callback, and the only drop metric (`container_events_dropped`). A third queue, the CRI `EventMonitor` backoff (`internal/cri/server/events/events.go`), is unbounded per container/sandbox key. Durable lifecycle facts (image/content/snapshot rows in BoltDB, task state in the shim) are committed before events are published, so the fact survives but the notification does not.

## Rating

**4/10 — Present but inconsistent, weakly documented, fragile.**

Rationale: one peripheral queue (CRI `eventq`, with tests and a drop counter) models retention explicitly, and the shim publisher models bounded retry/drop explicitly. But the primary lifecycle bus (task exit/delete/OOM, image/snapshot/content notifications) has no durability tier, no replay, no memory bound, no lag/queue-depth observability, and silently drops or OOMs under the exact stress in the dimension prompt (never-reading subscriber + flood). Failure handling is split across three unrelated queue implementations with no shared policy.

## Evidence Collected

Every entry MUST include a file path with line numbers. Format: `path/to/file.ts:NN`.

| Area | Evidence | File:Line |
|------|----------|-----------|
| Core fan-out entry point | `Exchange` wraps `goevents.Broadcaster`; `Publish` validates topic/namespace, marshals Any, stamps `time.Now().UTC()`, then `broadcaster.Write` | `studies/ultraplan-daemon-events-study/sources/containerd/core/events/exchange/exchange.go:36-45`, `:80-119` |
| Core forward path | `Forward` validates envelope (namespace, topic, non-zero timestamp) then same `broadcaster.Write` | `studies/ultraplan-daemon-events-study/sources/containerd/core/events/exchange/exchange.go:55-75`, `:224-238` |
| Subscriber registry + per-sub buffers | `Subscribe` builds `Channel(0)` (unbuffered) → `Queue(channel)` → optional `Filter(queue)`; registers `dst` via `broadcaster.Add`; forwarding goroutine blocks on `evch <- env` or `ctx.Done` | `studies/ultraplan-daemon-events-study/sources/containerd/core/events/exchange/exchange.go:128-199` |
| Broadcaster is synchronous, per-sink | Single `run()` loop pulls from unbuffered `b.events` and calls each `sink.Write` inline; a closed sink is removed, other write errors only log `broadcaster: dropping event` | `studies/ultraplan-daemon-events-study/sources/containerd/vendor/github.com/docker/go-events/broadcast.go:44-54`, `:105-161` |
| Per-subscriber queue is unbounded | `Queue.Write` only appends to `list.List` + signals; `run()`/`next()` blocks on cond; drop only logged at Debug on downstream `Write` error | `studies/ultraplan-daemon-events-study/sources/containerd/vendor/github.com/docker/go-events/queue.go:10-12`, `:35-47`, `:66-88` |
| Downstream channel is rendezvous | `Channel.Write` is `select { ch.C <- event \| <-closed }` on a `buffer`-sized chan; Exchange passes `0`, so `Queue.run` blocks until the subscribe goroutine drains `channel.C` | `studies/ultraplan-daemon-events-study/sources/containerd/vendor/github.com/docker/go-events/channel.go:20-42` |
| gRPC subscribe blocks on slow UI | `service.Subscribe` loops `ev := <-eventq` then `srv.Send(toProto(ev))`; a slow `Send` backpressures `evch`, hence `channel.C`, hence that subscriber's `Queue` list | `studies/ultraplan-daemon-events-study/sources/containerd/plugins/services/events/service.go:101-120` |
| Client/proxy subscribe same shape | `client/eventRemote.Subscribe` and `core/events/proxy` use unbuffered `evq`, `session.Recv()` loop, blocking `evq <- envelope` select vs `ctx.Done` — no buffer, no drop policy | `studies/ultraplan-daemon-events-study/sources/containerd/client/events.go:80-125`, `studies/ultraplan-daemon-events-study/sources/containerd/core/events/proxy/remote_events.go:99-147`, `:183-231` |
| Shim→daemon bounded retry + drop | `queueSize=2048`, `maxRequeue=5`; `processQueue` evicts with `log.Errorf("evicting %s ...")` after 5 tries; `queue()` requeues via `time.Sleep(1s*count)` in a new goroutine with blocking `requeue <- i` | `studies/ultraplan-daemon-events-study/sources/containerd/pkg/shim/publisher.go:35-38`, `:102-124` |
| Shim publish blocks hot path | `Publish` builds envelope, then synchronous `forwardRequest` with 5s timeout; on failure it `queue(i)` and returns the error; reconnect path retries once | `studies/ultraplan-daemon-events-study/sources/containerd/pkg/shim/publisher.go:127-152`, `:154-188` |
| Shim task terminal facts use that path | Dead-shim cleanup publishes `TaskExit` + `TaskDelete` and ignores errors; shim `oomEvent` publishes `TaskOOM` and only logs on error | `studies/ultraplan-daemon-events-study/sources/containerd/core/runtime/v2/shim.go:144-187`, `studies/ultraplan-daemon-events-study/sources/containerd/cmd/containerd-shim-runc-v2/task/service.go:698-705` |
| Durable write precedes event, then couples | Bolt `Update`/commit happens first; `publisher.Publish` is called after; on publish error the API returns the error even though the row is committed (image create/update, content commit, snapshot prepare/commit) | `studies/ultraplan-daemon-events-study/sources/containerd/core/metadata/images.go:164-175`, `studies/ultraplan-daemon-events-study/sources/containerd/core/metadata/content.go:624-641`, `studies/ultraplan-daemon-events-study/sources/containerd/core/metadata/snapshot.go:274-291` |
| GC events are best-effort | `publishEvents` logs at Debug and continues on `Publish` failure; only `ImageDelete`/`SnapshotRemove` topics are emitted | `studies/ultraplan-daemon-events-study/sources/containerd/core/metadata/db.go:348-368` |
| OOM publishers are fire-and-log | cgroupsv1/v2 watchers and daemon cgroup metrics publish `TaskOOM` and only `log ... publish OOM event` on error | `studies/ultraplan-daemon-events-study/sources/containerd/pkg/oom/v2/v2.go:73-78`, `studies/ultraplan-daemon-events-study/sources/containerd/core/metrics/cgroups/v1/cgroups.go:90` |
| CRI retention tier (time-bounded) | `eventq.New[T](discardAfter, discardFn)`: with zero subscribers, events accumulate in `discardQueue` with `discardAt=now+discardAfter`; late subscriber replays the held window; expiry calls `discardFn`; `Subscribe()` chan cap is `100` | `studies/ultraplan-daemon-events-study/sources/containerd/internal/eventq/eventq.go:53-63`, `:90-130`, `:154-163` |
| CRI retention configured at 5m + metric | `containerEventsQ = eventq.New[...](5*time.Minute, ... containerEventsDroppedCount.Inc() ...)`; only drop observability in repo for events | `studies/ultraplan-daemon-events-study/sources/containerd/internal/cri/server/service.go:236-244`, `studies/ultraplan-daemon-events-study/sources/containerd/internal/cri/server/metrics.go:39`, `:66` |
| CRI event production is non-blocking | `generateAndSendContainerEvent` ends with `containerEventsQ.Send(&event)`; `Send` blocks only on unbuffered `events` chan vs shutdown, not on subscribers | `studies/ultraplan-daemon-events-study/sources/containerd/internal/cri/server/helpers.go:369-377`, `studies/ultraplan-daemon-events-study/sources/containerd/internal/eventq/eventq.go:147-152` |
| CRI EventMonitor backoff is unbounded | `backOff.queuePool map[string]*backOffQueue`, `events []any` appended in `enBackOff`/`reBackOff` with exponential `1s→5m` ticker `1s`; no length cap, no drop, no metric | `studies/ultraplan-daemon-events-study/sources/containerd/internal/cri/server/events/events.go:54-71`, `:126-186`, `:231-260` |
| Windows shim log drop policy (explicit) | `reconnectingLogWriter.Write` drops while no reader attached and fakes `len(p), nil` on short write/conn loss so "logging never backpressures the shim"; documented as intentional discard, not buffered | `studies/ultraplan-daemon-events-study/sources/containerd/pkg/shim/shim_windows.go:211-218`, `:244-271` |
| Topic taxonomy (lifecycle set) | Task topics `/tasks/create|start|oom|exit|delete|exec-added|exec-started|paused|resumed|checkpointed`; `GetTopic` maps proto type → topic | `studies/ultraplan-daemon-events-study/sources/containerd/core/runtime/events.go:24-47`, `:51-77` |
| Tests cover fan-out/filter, not pressure | `TestExchangeBasic` (2 subs, 3 events), `TestExchangeFilters`, `TestExchangeValidateTopic`; `eventq` tests cover replay within window, expiry, shutdown discard — no slow-consumer, flood, or memory-bound test | `studies/ultraplan-daemon-events-study/sources/containerd/core/events/exchange/exchange_test.go:34-113`, `:115-271`, `studies/ultraplan-daemon-events-study/sources/containerd/internal/eventq/eventq_test.go:26-124` |

## Answers to Dimension Questions

1. **Which events must never be dropped?** No evidence found that any event class is guaranteed never-dropped. `TaskExit`/`TaskDelete` (`studies/ultraplan-daemon-events-study/sources/containerd/core/runtime/v2/shim.go:173-186`), `TaskOOM` (`studies/ultraplan-daemon-events-study/sources/containerd/pkg/oom/v2/v2.go:73-78`), and image/content/snapshot mutations (`studies/ultraplan-daemon-events-study/sources/containerd/core/metadata/images.go:168-175`, `studies/ultraplan-daemon-events-study/sources/containerd/core/metadata/content.go:631-639`, `studies/ultraplan-daemon-events-study/sources/containerd/core/metadata/snapshot.go:280-288`) are the closest to "must not lose", but every path can lose them: the core `Exchange` drops only via undeliverable-sink logging (`studies/ultraplan-daemon-events-study/sources/containerd/vendor/github.com/docker/go-events/broadcast.go:119-129`), the shim forwarder evicts after 5 retries (`studies/ultraplan-daemon-events-study/sources/containerd/pkg/shim/publisher.go:102-108`), and GC notifications are debug-log-and-continue (`studies/ultraplan-daemon-events-study/sources/containerd/core/metadata/db.go:363-365`). Durability lives in BoltDB rows and shim task state, not in the event bus.

2. **Which events are ephemeral or short-lived?** All `Exchange` envelopes are ephemeral: in-memory only, no log, no replay for subscribers that attach late (`studies/ultraplan-daemon-events-study/sources/containerd/core/events/exchange/exchange.go:128-199`). CRI `ContainerEventResponse`s are short-lived with a 5-minute hold when nobody listens, then passed to a discard callback (`studies/ultraplan-daemon-events-study/sources/containerd/internal/eventq/eventq.go:90-130`, `studies/ultraplan-daemon-events-study/sources/containerd/internal/cri/server/service.go:236-244`). Windows shim log bytes with no reader are explicitly ephemeral/discarded (`studies/ultraplan-daemon-events-study/sources/containerd/pkg/shim/shim_windows.go:244-271`).

3. **Can a slow UI stall execution?** No global stall was found, but isolation is by unbounded buffering, not by policy. `srv.Send` blocking (`studies/ultraplan-daemon-events-study/sources/containerd/plugins/services/events/service.go:108-111`) stalls only that subscription's forwarding goroutine → that subscriber's `Channel(0)` (`studies/ultraplan-daemon-events-study/sources/containerd/vendor/github.com/docker/go-events/channel.go:35-42`) → that subscriber's `Queue` list (`studies/ultraplan-daemon-events-study/sources/containerd/vendor/github.com/docker/go-events/queue.go:66-88`). `Broadcaster.run` stays live because `Queue.Write` is a mutex append (`studies/ultraplan-daemon-events-study/sources/containerd/vendor/github.com/docker/go-events/queue.go:35-47`). The one synchronous coupling against the hot path is metadata: `Publish` after commit runs inline in the mutating RPC and its error fails the RPC (`studies/ultraplan-daemon-events-study/sources/containerd/core/metadata/images.go:168-175`), so a wedged broadcaster (unbuffered `b.events` in `studies/ultraplan-daemon-events-study/sources/containerd/vendor/github.com/docker/go-events/broadcast.go:47-54`) can add latency to metadata writes; and shim `Publish` blocks up to 5s per forward (`studies/ultraplan-daemon-events-study/sources/containerd/pkg/shim/publisher.go:154-160`).

4. **How are lifecycle and terminal facts flushed under pressure?** They are not flushed to durable storage as events at all. Ordering is commit-then-notify for metadata (fact durable in Bolt, notification best-effort), and publish-and-retry (max 5, ~1+2+3+4+5s delays, buffer 2048) for shim→daemon (`studies/ultraplan-daemon-events-study/sources/containerd/pkg/shim/publisher.go:102-124`). Under sustained daemon outage the shim evicts with only an error log, so `TaskExit`/`TaskOOM` can vanish while the container state itself remains queryable. CRI compensates for handler failure with an unbounded per-ID backoff requeue (`studies/ultraplan-daemon-events-study/sources/containerd/internal/cri/server/events/events.go:144-152`, `:161-171`), i.e. retry without a flush deadline.

5. **What bounds storage and in-memory growth?** Almost nothing on the core path: no byte/count cap, no TTL, no compaction, no sampling. Bounds found: (a) shim `requeue` channel cap 2048 + 5-retry eviction (`studies/ultraplan-daemon-events-study/sources/containerd/pkg/shim/publisher.go:35-38`); (b) CRI `eventq` time bound 5m + per-subscriber chan cap 100 (`studies/ultraplan-daemon-events-study/sources/containerd/internal/eventq/eventq.go:154-163`, `studies/ultraplan-daemon-events-study/sources/containerd/internal/cri/server/service.go:237`); (c) implicit OS/socket buffers on gRPC streams. Unbounded: per-subscriber `Queue` `list.List`, CRI backoff `events []any` per ID, broadcaster sink slice. No evidence found of retention/cleanup jobs, WAL/log growth metrics, or lag/queue-depth metrics for the core bus; the sole event-drop metric is `container_events_dropped` for the CRI queue.

> Stress answer: attach a subscriber that never reads while runtime output floods the system. The daemon does not deadlock: publishers keep succeeding (fast `Queue.Write` appends), the dead subscriber's `Queue` list grows without bound until the daemon OOMs or the subscriber cancels; `Broadcaster.Add/Remove` for other subscribers still cycle through the same run loop so churn gets slower as the sink slice grows. What remains durable is whatever was committed to Bolt (images/content/snapshots) and shim-side task state — but any envelope published while the slow subscriber lagged, and any envelope published before a new subscriber attached, is unrecoverable (no replay). The CRI `GetContainerEvents` consumer is the exception: it replays up to the last 5 minutes on subscribe, then drops with a counter.

## Architectural Decisions

- **Broadcast, not log.** `Exchange` + `docker/go-events` (`studies/ultraplan-daemon-events-study/sources/containerd/core/events/exchange/exchange.go:36-45`) chooses fire-and-forget fan-out over a persistent log. Consequence: ordering per subscriber is FIFO via its queue, but there is no cross-restart continuity.
- **Decouple via unbounded queue.** Each subscriber gets its own `Queue` (`studies/ultraplan-daemon-events-study/sources/containerd/core/events/exchange/exchange.go:132-134`) so one slow reader cannot block `Broadcaster.run`. Consequence: backpressure becomes memory pressure.
- **Rendezvous handoff at the tail.** `NewChannel(0)` (`studies/ultraplan-daemon-events-study/sources/containerd/core/events/exchange/exchange.go:132`) forces `Queue.run` to rendezvous with the forwarding goroutine, which in turn rendezvouses with `srv.Send`. Consequence: the queue depth signal is hidden inside `go-events`, invisible to operators.
- **Commit-then-notify in metadata.** Bolt transaction first, `Publish` second (`studies/ultraplan-daemon-events-study/sources/containerd/core/metadata/images.go:164-175`). Consequence: facts are durable, notifications are not, yet publish errors still fail the RPC — a confusing half-coupling.
- **Retry-then-drop at the trust boundary.** Shim→daemon uses 5s-timeout forwards with delayed requeue and eviction (`studies/ultraplan-daemon-events-study/sources/containerd/pkg/shim/publisher.go:102-152`). Consequence: transient daemon restarts are bridged (~15s window), longer outages lose terminal events.
- **Separate time-boxed bus for CRI streaming.** `eventq` with 5m hold + discard hook + counter (`studies/ultraplan-daemon-events-study/sources/containerd/internal/cri/server/service.go:236-244`). Consequence: kubelet-style watchers get a replay window the native event API lacks.

## Notable Patterns

- **Three queue dialects, no shared contract.** `go-events Queue` (unbounded, debug-log drop), `shim requeue` (bounded 2048, 5-retry drop), `eventq` (time-bounded, counted drop), plus CRI `backOff` (unbounded, silent growth). Each has its own tests or none.
- **Filter at the sink.** Topic/field filtering wraps the queue (`studies/ultraplan-daemon-events-study/sources/containerd/core/events/exchange/exchange.go:147-158`), so non-matching events never enter the per-subscriber queue — the only volume-reduction mechanism on the core path.
- **Log-and-continue for derived events.** OOM (`studies/ultraplan-daemon-events-study/sources/containerd/pkg/oom/v2/v2.go:76-78`), GC (`studies/ultraplan-daemon-events-study/sources/containerd/core/metadata/db.go:363-365`), shim OOM (`studies/ultraplan-daemon-events-study/sources/containerd/cmd/containerd-shim-runc-v2/task/service.go:702-704`) never propagate publish failures to state transitions.
- **Fail-the-RPC for primary mutations.** Image/content/snapshot publishes return errors to the caller post-commit — the inverse pattern, coupling API latency/success to bus health.
- **Drop-logs-never-block-logs on Windows.** `reconnectingLogWriter` (`studies/ultraplan-daemon-events-study/sources/containerd/pkg/shim/shim_windows.go:244-271`) is the clearest backpressure statement in the tree: prefer data loss over blocking the runtime.

## Tradeoffs

- **Isolation vs. memory safety.** Per-subscriber unbounded queues isolate publishers from slow consumers but trade a stall for an OOM. A bounded queue with drop + counter (like `eventq`) would invert the tradeoff.
- **Availability vs. consistency on publish.** Post-commit publish keeps write availability (Bolt never waits for subscribers) but permits committed-without-notified states; returning the publish error does not roll back, so callers cannot distinguish the two.
- **Transient resilience vs. terminal loss.** Shim retry bridging (~15s, 2048 slots) covers daemon restarts but guarantees nothing; terminal `TaskExit` has the same policy as any other forward.
- **Observability concentrated in one corner.** The CRI queue has a counter and log line per discard; the far higher-volume core bus has only trace logs per publish/forward (`studies/ultraplan-daemon-events-study/sources/containerd/core/events/exchange/exchange.go:60-72`, `:104-116`) — nothing to alert on lag or drops.
- **No batching anywhere.** Every envelope is dispatched singly through mutex + cond + channel hops; high-volume progress/output deltas would pay per-event overhead with no coalescing knob.

## Failure Modes / Edge Cases

- **Never-reading subscriber → daemon OOM.** `Queue` list grows per published envelope; no cap, no eviction, no metric. First degradation is memory, then GC pressure, then `Broadcaster.run` slowdown for everyone.
- **Daemon down > retry window → terminal event loss.** Shim evicts after `count > maxRequeue` (`studies/ultraplan-daemon-events-study/sources/containerd/pkg/shim/publisher.go:104-108`); `cleanupAfterDeadShim` ignores publish errors (`studies/ultraplan-daemon-events-study/sources/containerd/core/runtime/v2/shim.go:173-186`), so exit/delete facts exist only in shim state.
- **Requeue-channel saturation parks goroutines.** `l.requeue <- i` inside `queue()`'s spawned goroutine (`studies/ultraplan-daemon-events-study/sources/containerd/pkg/shim/publisher.go:117-124`) blocks that goroutine when 2048 are pending; a flood of failures leaks sleeping goroutines.
- **CRI backoff hyper-accumulation.** Handler outage appends every event for an affected ID forever (`studies/ultraplan-daemon-events-study/sources/containerd/internal/cri/server/events/events.go:232-241`); recovery replays the whole slice inline on the ticker goroutine, delaying other IDs.
- **`eventq` replay stampede.** A subscriber attaching after an idle burst receives the entire 5-minute backlog synchronously in `subscriberC` handling (`studies/ultraplan-daemon-events-study/sources/containerd/internal/eventq/eventq.go:99-112`), blocking the central `eventq` loop for all other producers/subscribers meanwhile.
- **Filter-parse failure closes subscription.** `Subscribe` with bad filters sends one error and closes (`studies/ultraplan-daemon-events-study/sources/containerd/core/events/exchange/exchange.go:147-153`); clients get no partial stream.
- **No post-restart catch-up.** Daemon or `Exchange` restart drops all in-flight envelopes; new subscribers start from now. No evidence of snapshot/replay markers.

## Future Considerations

- Put one explicit policy on the core bus: bounded per-subscriber queue (e.g. N envelopes) with drop-oldest + `events_dropped_total{topic, reason}` counter and trace log, mirroring the `eventq` discard hook.
- Separate tiers by topic class: durable control topics (`/tasks/exit`, `/tasks/delete`, `/tasks/oom`) onto a persistent or at-least-replayed channel; keep progress/log deltas on the ephemeral broadcast with sampling/coalescing.
- Export queue depth and publish latency: per-subscriber depth gauge, broadcaster dispatch latency histogram, shim `requeue` depth gauge — currently only `container_events_dropped` exists.
- Decouple metadata RPC success from notification success: publish asynchronously or return success-with-warning so commit and notify states are distinguishable; document that subscribers must reconcile via List/Status, not rely on exactly-once delivery.
- Cap CRI backoff queues (max events per ID + drop policy + metric) or spill to disk; add a test with a failing handler under flood.

## Questions / Gaps

- No evidence found for maximum envelope size, per-topic rate limits, or payload sampling — searched `core/events`, `plugins/services/events`, `pkg/shim/publisher.go`, `internal/eventq`, `internal/cri/server/events`.
- No evidence found for retention/compaction/pruning jobs or WAL for events — the bus is memory-only; Bolt stores facts, not envelopes.
- No evidence found for lag/queue-depth/WAL-growth dashboards or alerts for the core bus — only `container_events_dropped` for the CRI queue.
- Unknown: expected client recovery contract after `Subscribe` error/EOF (reconcile via snapshot API vs. resume token) — no resume token or cursor exists in `api/services/events/v1`, `api/services/ttrpc/events/v1`, or `client/events.go`.
- Unknown: whether `ctr events` / CRI watchers apply any client-side buffering or coalescing under flood — both read-and-print/send directly with no buffer.

---

Generated by `Dimension 01.07: Event Delivery, Backpressure, and Retention Tiers` against `containerd`.
