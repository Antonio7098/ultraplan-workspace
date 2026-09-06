# Source Analysis: nats-server

## 01.04 Ordered Observation, Live Streaming, and Backpressure

### Source Info

| Field | Value |
|-------|-------|
| Name | nats-server |
| Path | `studies/aren-go-runtime-study/sources/nats-server` |
| Language / Stack | Go (NATS server v2.15.0-dev, JetStream, Raft, file/mem stores) |
| Analyzed | 2026-08-29 |

## Summary

nats-server decouples ordered durable storage from live observation across two layers: core NATS per-client fan-out (`server/client.go`) and JetStream ordered log + consumer delivery (`server/consumer.go`, `server/memstore.go`/`filestore.go`, `server/stream.go`). Ordering is total per-stream via monotonically assigned `LastSeq`/`FirstSeq` (`server/memstore.go:288-293`, `server/stream.go:state()`) and per-consumer replay via `sseq`/`dseq`/`asflr`/`adflr` cursors (`server/consumer.go:450-453`, `server/consumer.go:6267-6326`). Live streaming is pull/push over `jsOutQ/ipQueue` (`server/stream.go:8497-8512`) and `loopAndGatherMsgs` (`server/consumer.go:5343-5608`) with isolated `client.out.nb/pb` per observer (`server/client.go:349-361`). Backpressure for core clients is bounded `MAX_PENDING_SIZE=64MB` (`server/const.go:102`) + `WriteDeadline` (`server/const.go:132`) leading to `SlowConsumerPendingBytes`/`SlowConsumerWriteDeadline` disconnect or `Retry` for routes/leafs/gateways (`server/client.go:736-752`, `2598-2615`, `1947-2005`). JetStream never drops retained facts; it stalls delivery on `MaxAckPending` (`server/consumer.go:4921`, `4982`, `5418-5421`) and on byte-based flow-control window `maxpb/pbytes/fcid` (`server/consumer.go:568-569`, `5824-5842`, `5912-5921`). Reconnect is split: ephemeral core subscriptions are best-effort with no cursor, while durable JetStream consumers persist `sseq/asflr` in the consumer store and resume exactly from the next sequence, including redelivery queues (`server/consumer.go:6014-6045`).

## Rating

**6 / 10**

**Rationale:** Per-stream total ordering, per-consumer monotonic cursors, per-observer `nb/pb` isolation, and two-tier backpressure (disconnect for core, flow-control + `MaxAckPending` stall for JetStream) are production-hardened and well-tested. Durable facts are retained in the store and infinitely replayable; lossy live deltas are pushed through isolated outbound buffers with explicit `MaxPending`/`WriteDeadline` bounds. Points deducted because (1) core NATS has no durable/live distinction — there is no coalescing or best-effort delta path with dropped-event accounting, only `slowConsumers` counters; (2) producer isolation is incomplete — a slow consumer creates a `stc` stall gate (`server/client.go:3622-3624`) and `stalledWait` (`server/client.go:3711-3749`) that throttles fast producers up to `stallTotalAllowed=10ms` (`server/client.go:127`), coupling one slow observer to producers; (3) there is no single "exactly one terminal fact" delivery primitive; stream `LastSeq` advancement and consumer `EOF`/`408`/`404` headers are stream-level, not a per-execution terminal guarantee; (4) reconnect for non-durable observers is new subscription from `DeliverNew`/`DeliverLast` etc., not cursor-based resumption.

## Evidence Collected

Every entry MUST include a file path with line numbers. Format: `path/to/file.ts:NN`.

| Area | Evidence | File:Line |
|------|----------|-----------|
| Event creation & sequencing | `storeRawMsg` enforces `seq == LastSeq+1` else `ErrSequenceMismatch`, updates `FirstSeq/LastSeq/FirstTime/LastTime`, per-subject `SimpleState.First/Last/Msgs` | `server/memstore.go:288-324` |
| Ordered store state | `StreamState.FirstSeq/LastSeq/Msgs/Bytes`, `FastState` | `server/memstore.go:36-54` + `server/stream.go:9161-9176` |
| Consumer ordered cursors | `sseq` next stream seq, `dseq` next consumer seq, `asflr/adflr` ack floors, `reconcileStateWithStream` rewinding floors on data loss | `server/consumer.go:450-453`, `6267-6326`, `6328-6444` |
| Replay & starting seq | `selectStartingSeqNo` handles `DeliverAll/Last/LastPerSubject/ByStartSeq/ByStartTime/DeliverNew`, clamps to `FirstSeq/LastSeq` | `server/consumer.go:6328-6444` |
| `LastPerSubject` skip list | `lastSeqSkipList.resume/seqs` + `MultiLastSeqs` | `server/consumer.go:6255-6280`, `server/memstore.go:903-962` |
| Fan-out | `deliverMsg` per `subscription.client`, `processMsgResults` fan-out via `SublistResult` + `routeTarget`, per-client `queueOutbound` | `server/client.go:3754-3988` |
| Per-observer buffering | `outbound.nb net.Buffers`, `pb int64`, `wnb`, `mp`, `wdl`, `stc chan struct{}`; pooled buffers `nbPoolSmall/Medium/Large` | `server/client.go:349-361`, `363-408` |
| Global pending bound | `MAX_PENDING_SIZE = 64 MiB`, `MAX_PENDING` default, per-client `out.mp = opts.MaxPending` | `server/const.go:101-102`, `server/client.go:753` |
| Write deadline bound | `DEFAULT_FLUSH_DEADLINE=10s`, per-kind `WriteDeadline`, `nc.SetWriteDeadline(time.Now().Add(wdl))` loop | `server/const.go:132`, `server/client.go:1842-1856`, `opts.go:6128-6129` |
| Pending-bytes slow consumer | `if pb > mp` → decrement, `slowConsumers++`, `markConnAsClosed(SlowConsumerPendingBytes)` | `server/client.go:2598-2615` |
| Write-deadline slow consumer | `handleWriteTimeout` increments `slowConsumers`, `scStats`, logs `WriteDeadline exceeded`, `markConnAsClosed(SlowConsumerWriteDeadline)` unless `wtp==Retry` for routes/leafs | `server/client.go:1947-2005` |
| Write timeout policy | `WriteTimeoutPolicyClose` for clients vs `Retry` for `ROUTER/LEAF/GATEWAY` | `server/client.go:736-752` |
| Producer stall gate | `if pb > mp*3/4 → stc=make(chan)`, `stalledWait` sleeps `2ms..5ms` up to `10ms` total per readLoop, closed on recovery `close(stc)` | `server/client.go:3622-3624`, `3711-3749`, `1934-1937`, `127` |
| Flush path ordering & framing | `flushOutbound` copies `nb→wnb`, compress/frame outside lock, `WriteTo` with deadline, returns `n/attempted` accounting, signals `flushSignal` if `pb>0` | `server/client.go:1718-1945` |
| JetStream durable facts | `memStore`/`filestore` retain all messages under `StreamConfig.Retention/MaxMsgs/MaxBytes`, `purge/compact/truncate` are explicit | `server/memstore.go:1241-1273`, `1591-1675`, `1812-1860` |
| Live consumer delivery | `loopAndGatherMsgs`: if push inactive or `pbytes>maxpb` or pull `waiting.isEmpty` → `waitForMsgs`; `getNextMsg` → `deliverMsg`; heartbeats via `hbc` | `server/consumer.go:5343-5608` |
| Flow-control window | `pblimit = 32 MiB` (`JsFlowControlMaxPending`), `maxpb = pblimit/16`, `needFlowControl` at `pbytes > maxpb/2`, `sendFlowControl` with `fcid/fcsz`, receiver `processFlowControl` halves window or doubles | `server/consumer.go:601-602`, `1845`, `5824-5842`, `5912-5921`, `5844-5868` |
| MaxAckPending backpressure | `errMaxAckPending`, `getNextMsg` returns it when `len(pending) >= maxp`; `loopAndGatherMsgs` stalls on that error; `checkPending` redeliveries respect `AckWait/BackOff` | `server/consumer.go:4921`, `4982`, `5441`, `6070-6186` |
| Exactly-once pending tracking | `trackPending` map `pending[sseq]=Pending{dseq,timestamp}`, `updateDelivered`, `rdq/rdqi` redelivery queue, `ackWait` timers | `server/consumer.go:5889-6067` |
| Pull-mode waiting queue | `waitingRequest{reply,n,b,expires,hbt}`, `processWaiting` emits `408 Request Timeout` or `404 No Messages`, heartbeats `100 Idle Heartbeat` with `JSLastConsumerSeq/JSLastStreamSeq` | `server/consumer.go:5038-5137`, `5610-5621` |
| Rate limiting | `o.rlimit.ReserveN(now, sz)` delay before `deliverMsg` | `server/consumer.go:5524-5539` |
| Ordered stream → observer bridge | `jsOutQ` wraps `ipQueue[*jsPubMsg]` with `push`/`pop`/`recycle`, `internalLoop` drains `outq/msgs/gets/ackq` via single `createInternalJetStreamClient` | `server/stream.go:8497-8512`, `8570-8701` |
| Bounded internal queues | `ipQueue` supports `ipqLimitByLen`/`ipqLimitBySize` with `errIPQLenLimitReached`/`errIPQSizeLimitReached`; JS API queues use `queueLimit=10_000` and drain on overflow | `server/ipqueue.go:68-82`, `114-189`, `server/jetstream_api.go:940-956` |
| Loss accounting | `slowConsumers int64`, `scStats.clients/routes/gateways/leafs`, `acc.stats.slowConsumers`, varz `SlowConsumers` | `server/server.go:421`, `client.go:1974-1993`, `2606-2612`, `monitor.go:1277`, `server/events.go:959` |
| Reconnect / resume | Durable consumer `sseq = state.LastSeq+1` or filtered `LastSeq`, `asflr/adflr` floors persisted via `store.SetStarting`; ephemeral pull uses `NextSeq`/`sseq` in memory only | `server/consumer.go:6328-6442`, `6247-6253` |
| No cross-observer blocking (claim vs reality) | `flushClients` with `budget` and `fsp` accounting attempts to avoid blocking readLoop on flush, but `stc` still stalls producers | `server/client.go:1418-1448`, `3680-3749` |

## Answers to Dimension Questions

**Can one slow observer delay execution or other observers?**
No for other observers, partially yes for producers. Core NATS isolates outbound state per `client` (`server/client.go:349-361`): each `queueOutbound` appends to that client's `nb/pb` only (`server/client.go:2566-2596`) and `flushOutbound` operates on a `wnb` working copy outside the producer's lock (`server/client.go:1764-1783`). Other clients' `pb` are untouched. However, when a client's `pb > 3/4 mp` a per-client stall gate `stc` is created (`server/client.go:3622-3624`) and any producer that delivers to that client via `deliverMsg` blocks in `stalledWait` for `2-5ms` per call, up to `10ms` total per read loop (`server/client.go:127`, `3711-3749`). This throttles fast fan-in producers but does not delay other consumers' already-queued `flushOutbound` writes (driven by separate `writeLoop`s per client, `server/client.go:1368-1413`). JetStream consumers are further isolated: `loopAndGatherMsgs` blocks only that consumer when `pbytes > maxpb` or `errMaxAckPending` (`server/consumer.go:5418-5442`); the stream's `internalLoop` (`server/stream.go:8570-8701`) and store ingest continue.

**What is bounded, dropped, coalesced, retained, or replayed under pressure?**
- **Bounded:** Per-client `pb` ≤ `mp` (`MAX_PENDING_SIZE` 64 MiB, `server/const.go:101-102`, `server/client.go:753`). Per-write `nbMaxVectorSize=1024` (~64 MiB per `WriteTo` batch, `server/client.go:363`, `1830-1842`). JetStream consumer bytes `pbytes` ≤ `maxpb` (`JsFlowControlMaxPending=32 MiB` window, `server/consumer.go:601-602`, `5688-5699`) and counts `pending` ≤ `MaxAckPending` (default 1000, `server/consumer.go:604`, `684-705`). Internal `ipQueue`/`jsOutQ` can be bounded by `mlen`/`msz` (`server/ipqueue.go:114-128`); API ingress bounded at `JSDefaultRequestQueueLimit=10_000` (`server/jetstream_api.go:403-405`).
- **Dropped:** Core NATS drops by disconnecting the slow observer: `SlowConsumerPendingBytes` when `pb>mp` (`server/client.go:2606-2614`) and `SlowConsumerWriteDeadline` when `WriteTo` times out (`server/client.go:1947-2004`). No message is dropped silently for other observers. JetStream store never drops due to consumer slowness; it stalls.
- **Coalesced:** No evidence of delta coalescing for observation. Per-subject `MaxMsgsPer` will replace oldest subject message (`server/memstore.go:333-337`) but that is retention, not a live observer optimization. Search for coalesce delta yielded no relevant code.
- **Retained:** All stream messages are retained in `memStore.msgs`/`filestore` until explicit `RemoveMsg/Purge/Compact` or limits/age (`server/memstore.go:1591-1675`). Consumer pending `map[seq]Pending` retains delivery state until ack (`server/consumer.go:5924-5956`).
- **Replayed:** Fully replayable from any `sseq` via `LoadNextMsg`/`LoadNextMsgMulti`/`FilteredState`/`NumPending` (`server/memstore.go:601-818`, `server/consumer.go:4914-5036`). Missed sequences while stalled are not lost for durable consumers — `rdq/rdqi` redelivers on `AckWait/BackOff` (`server/consumer.go:6070-6186`).

**How is exactly one terminal fact delivered after partial streaming?**
No first-class terminal fact primitive exists for core NATS streams. The closest is JetStream's `StreamState.LastSeq/FirstSeq` versus consumer `sseq-1` reaching `LastSeq` (EOF, `ErrStoreEOF`, `server/consumer.go:5438-5442`, `5672-5683`) plus per-request pull terminators: `408 Request Timeout` and `404 No Messages` headers on `waitingRequest` expiry (`server/consumer.go:5084-5095`), and `100 Idle Heartbeat` (`server/consumer.go:5610-5620`). Stream-level advisories (`$JS.EVENT.ADVISORY.STREAM.QUORUM_LOST`, etc., `server/jetstream_api.go:318-343`) are separate from consumer delivery. There is no server-guaranteed single delivery of a final execution summary distinct from the ordered log tail; consumers detect completion by observing that `pending==0` and `sseq > LastSeq`, but under truncation/purge races `checkAckFloor`/`reconcileStateWithStream` may move floors (`server/consumer.go:5153-5247`, `6267-6326`). So exactly-once terminal delivery is not enforced by the protocol — idempotent client handling is required.

**Does reconnect recover semantics or merely resume best-effort display data?**
Both, partitioned by durability. **Ephemeral/core observers** resume best-effort: subscriptions are ephemeral (`server/client.go:642-661`), closed connections drop `nb/wnb` (`server/client.go:1734-1738`), and reconnect creates new `sid`s with no cursor — continuation depends on client-chosen `DeliverPolicy` (`DeliverNew/All/Last/LastPerSubject/ByStartSeq`, `server/consumer.go:313-340`). **Durable JetStream consumers** recover semantics: `ConsumerConfig.Durable` (`server/consumer.go:90`) persists `sseq/dseq/asflr/adflr/pending/rdc` in the consumer store (`consumerStore` via `writeStoreState`, `readStoredState`), so after disconnect the server resumes from `sseq` exactly, redelivers unacked `pending` via `rdq` (`server/consumer.go:6014-6045`), and `processNextMsgRequest` can request batch from any sequence. The envelope distinguishes headers (`JSLastConsumerSeq`, `JSLastStreamSeq`, `JSConsumerStalled`, `JSMsgSize`, `server/consumer.go:5612-5620`, `5728-5752`) but does not tag "fact vs delta" in core; JetStream's separation is implicit: retained log (facts) vs transient `outq` push.

## Architectural Decisions

- **Per-client outbound buffers with pooled slices** (`server/client.go:349-423`, `2560-2625`): each observer owns `nb`/`pb`; `nbPoolGet/Put` avoids allocations. Enables per-observer isolation and bounded memory.
- **Two-path backpressure split** (core vs JetStream): core disconnects slow clients immediately on `pb>mp` or write deadline; JetStream stalls consumers via `MaxAckPending` and byte `FlowControl` instead of disconnecting, because the log is durable (`server/client.go:2598-2615` vs `server/consumer.go:5418-5442`, `5824-5842`).
- **Separate `writeLoop` per connection** (`server/client.go:1368-1413`): `flushOutbound` runs outside producer lock with `wnb` copy and deadline-bounded `WriteTo` loop, preventing producer goroutine from doing blocking I/O directly (producer only does `queueOutbound` + `flushSignal`).
- **Ordered log with stream-level sequencing, consumer-level cursors** (`server/memstore.go:288-324`, `server/consumer.go:450-453`, `6328-6444`): total order at the stream, independent ordered view per consumer. Allows reconnect at arbitrary `seq`/time without global coordination.
- **Internal message bus via `ipQueue[*jsPubMsg]/jsOutQ`** (`server/stream.go:8497-8512`, `8570-8701`, `server/ipqueue.go:26-142`): decouples ingestion (`msgs` queue), direct gets (`gets`), outbound sends (`outq`), and acks (`ackq`) into separate channels processed by a single `internalLoop` client, preserving ordering per queue while avoiding cross-queue head-of-line blocking.
- **`WriteTimeoutPolicy` typed by kind** (`server/client.go:239-257`, `736-752`): clients `Close`, routes/leafs/gateways `Retry` (mark `isSlowConsumer` but keep connection). Avoids cascading mesh partitions on transient slowness.
- **JetStream flow-control adaptive window** (`server/consumer.go:5854-5867`): `maxpb` doubles on each `processFlowControl` ack until `pblimit`; avoids static window underutilization.
- **API ingress shedding via `JSDefaultRequestQueueLimit`** (`server/jetstream_api.go:940-956`): `apiDispatch` pushes to `jsAPIRoutedReqs` and drains/drops when `pending >= limit`, publishes `JSAdvisoryAPILimitReached`. Protects execution from observer query overload.

## Notable Patterns

- **Isolation via per-observer `out` struct** — no shared outbound buffer; fan-out loop calls `queueOutbound` per `subscription.client` (`server/client.go:3976-3995` shows per-client framing).
- **Working-copy flush** — `nb → wnb` swap under lock, I/O outside lock, partial-write `wnb` retention (`server/client.go:1764-1792`, `1886-1893`) preserves order across partial writes without head-of-line stall.
- **Stall gate throttling** — `stc` channel closed on recovery; producers wait on channel with timeout rather than busy-loop (`server/client.go:3711-3749`).
- **Heartbeat as liveness, not data** — idle heartbeats carry `JSLastConsumerSeq/JSLastStreamSeq` and optionally `JSConsumerStalled` (`server/consumer.go:5610-5621`) separately from payload; WebSocket/MQTT callbacks are distinct handlers (`server/mqtt.go:5084-5370`).
- **Cursor-based replay** — `LoadNextMsg(filter,withHeaders,seq,sm)` + `FilteredState` compute next applicable sequence efficiently using `SubjectTree` (`server/memstore.go:601-818`). Enables gap-aware replay (`consumer.go:5025-5034` updates `skipped`).
- **Pooled `jsPubMsg` / `StoreMsg`** — `getJSPubMsgFromPool`/`returnToPool`, `nbPool` pools (`server/stream.go:8468-8494`, `server/client.go:369-408`) bound allocations under high-volume deltas.
- **Bounded shedding at edge** — `ipQueue.push` returns `errIPQLenLimitReached`/`errIPQSizeLimitReached` (`server/ipqueue.go:117-127`) and `jetstream_api.go` drops API requests over limit rather than blocking store ingestion.

## Tradeoffs

- **Hard disconnect vs. stall:** Core's `MaxPending` disconnect (`server/client.go:2606-2614`) is simple and prevents unbounded memory, but loses exactly-once semantics for that observer; JetStream's stall (`pbytes>maxpb`/`MaxAckPending`) preserves facts at cost of increased latency and consumer timer state (`ptmr`, `server/consumer.go:5952-5955`).
- **Fast-producer stall couples producer to one slow consumer:** `stalledWait` throttles the producer's `readLoop` (`server/client.go:3711-3749`, `1450-1455`) up to 10 ms. Protects slow consumers from overflow but violates ideal observer-does-not-control-execution; under fan-in 1: N a single slow client throttles all producers routed through that read loop iteration.
- **Per-connection `writeLoop` vs. shared flush pool:** Dedicated goroutine + `sg sync.Cond` (`server/client.go:355-357`, `4102-4109`) gives isolation but costs goroutine per connection (scales with `DEFAULT_MAX_CONNECTIONS=64k`, `server/const.go:105`).
- **Fixed `MAX_PENDING_SIZE=64MiB` everywhere:** Uniform bound is predictable but not adaptive to message size or consumer priority; `opts.MaxPending` is tunable (`server/opts.go:6128-6129`) but still global default.
- **Durable vs. ephemeral split leaks to API:** Clients must choose durable names and `AckPolicy` to get replay guarantees; misconfigured ephemeral pull with `no_wait` will see `404/408` terminators (`server/consumer.go:5084-5095`) that look like terminal facts but are not.
- **Flow-control window doubling:** Adaptive `maxpb *=2` (`server/consumer.go:5854-5856`) improves throughput on fast consumers but can overshoot memory on bursty consumers before feedback.

## Failure Modes / Edge Cases

- **Pending-bytes explosion under fan-out:** One producer publishing 1 MiB payloads to 1k subscriptions can queue `1k * 1 MiB = 1 GiB` transiently across `nb`s before any `pb>mp` check fires per client. The check is per enqueue (`server/client.go:2600`), so detection is per client, not global — server OOM possible before disconnect fires if `MaxPending` is large.
- **WriteDeadline starved by large vector:** `flushOutbound` loops `wnb.WriteTo` in `nbMaxVectorSize` batches each with fresh `SetWriteDeadline(time.Now().Add(wdl))` (`server/client.go:1842-1843`). A 64 MiB buffer with 10s deadline can take `N*deadline` across batches if kernel is slow but not timing out per batch, delaying detection.
- **Stall gate leak under high fan-in:** `stc` created at 75% threshold (`server/client.go:3622`) but only closed when `n==attempted` or `pb < mp*3/4` (`server/client.go:1934-1937`). If `pb` hovers at 76-99% with successful partial writes, gate never closes, causing repeated `stalledWait` timeouts for many producers.
- **JetStream backpressure deadlock with `MaxAckPending=0` + `FlowControl`:** `config.MaxAckPending <=0` with `FlowControl` returns `JSConsumerAckFCRequiresMaxAckPendingError` (`server/consumer.go:795-796`), but if limits allow flow control without heartbeat, advisory `FlowControlNeedsHeartbeats` (`server/consumer.go:981-982`) can cause silent stall with no heartbeat to unblock.
- **Redelivery queue unbounded growth:** `rdq`/`rdqi` appends every expired pending (`server/consumer.go:6153-6155`) without size cap; on prolonged consumer outage pending can grow to millions, `checkPending` sorts `expired` slice (`server/consumer.go:6154`) causing GC pressure.
- **Lost terminal detection on clustered failover:** `reconcileStateWithStream` (`server/consumer.go:6267-6326`) resets `asflr` to `LastSeq` after stream truncation; a consumer that had seen `sseq==LastSeq+1` before failover may see `sseq` rewound by one, violating exactly-once terminal delivery expectation.
- **API queue drain discards facts:** `apiDispatch` drains entire `jsAPIRoutedReqs` when `pending >= limit` (`server/jetstream_api.go:943-945`) — concurrent `CONSUMER.INFO`/`STREAM.INFO` bursts cause all pending requests to be dropped, not just excess; callers get `JSClusterNotAvail` or timeout with no dropped-event accounting.
- **Per-subject limits silently drop old facts:** `enforcePerSubjectLimit` (`server/memstore.go:1231-1245`) deletes oldest subject message when `MaxMsgsPer` exceeded, without notifying live observers beyond `StreamState` shift — observers doing `DeliverAll` will skip sequences via `FilteredState` gap, not via explicit tombstone.

## Future Considerations

- Add per-consumer `dropped_msgs/dropped_bytes` counters (mirroring `pending_size` in `monitor.go:3064`) and expose in consumer `INFO` and heartbeat headers to give lossy observers explicit truncated accounting.
- Introduce coalescing for high-volume progress deltas (e.g., `JSConsumerStalled` heartbeats already exist; extend to value deltas) to bound `outq` under token-stream-like workloads without disconnecting.
- Replace global `stallTotalAllowed=10ms` (`server/client.go:127`) with per-producer token-bucket and weighted fair scheduling, so one slow consumer cannot throttle all producers on that read loop.
- Provide a first-class terminal frame (e.g., `NATS/1.0 200 Terminal` with final `JSLastStreamSeq` and `checksum`) distinct from `404/408` pull terminators and stream advisories, delivered exactly once via durable consumer `pending` semantics.
- Make durable/live distinction explicit in `ConsumerConfig`: add `RetentionPolicy`/`ReplayPolicy` presets for "facts" (durable, acked) vs. "deltas" (best-effort, lossy) with separate `MaxPending` windows.
- Expand `ipQueue` metrics (`inprogress`, `size`, `len`) to varz and add slow-consumer-style shedding policy per `jsOutQ` with configurable drop vs. block.
- Fix typo at `server/consumer.go:3340` comment ("superseeds" → "supersedes") to keep doc search reliable.

## Questions / Gaps

- No evidence of intra-message ordering test for partial `flushOutbound` writes interleaved with concurrent `queueOutbound` — how does `wnb` partial framing guarantee MQ header order when `c.out.nb` is modified outside the lock? (`server/client.go:1884-1886` note says `wnb` framing is left alone, but WebSocket `wsCollapsePtoNB` path not traced.)
- Search for `SSE`/`stream parser` yielded no observer SSE transport in this server repo; any SSE gap must be verified outside `nats-server` (likely CLI/client side). Searched `server/**/*.go` for `SSE|EventSource|text/event-stream` — no matches.
- Pending bytes semantics differ between `client.out.pb` (sum of `len(data)` enqueued) and `consumer.pbytes` (sum of `pmsg.size()` including subj/reply/hdr). Are they ever reconciled? `consumer.go:5768-5770` adds `psz` (which includes `dsubj/reply/hdr/msg`) while stream accounting uses raw msg size.
- Reconnect cursor for durable consumers under `DeliverLastPerSubject` with `MaxMsgsPer=1` uses `allLastSeqs` shortcut (`server/memstore.go:853-962`); behavior after stream purge of that subject's last seq during reconnect not documented.
- Exactly-one terminal fact guarantee cannot be verified from implementation alone; requires integration test that subscribes, streams `N` messages, terminates stream, and asserts single terminal delivery after partial `flushOutbound` failures. No such test found in `server/jetstream_*_test.go` with that name — search for `Slow.*Recover` shows only `Slow Consumer Recovered` log (`server/client.go:1940`), not guarantee.

---
Generated by `01.04-ordered-observation-live-streaming-and-backpressure` against `nats-server`.
