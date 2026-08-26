# Source Analysis: docker-agent

## 01.04 Ordered Observation, Live Streaming, and Backpressure

### Source Info

| Field | Value |
|-------|-------|
| Name | docker-agent |
| Path | `studies/aren-go-runtime-study/sources/docker-agent` |
| Language / Stack | Go (bubbletea TUI, Echo HTTP server, SSE transport) |
| Analyzed | 2026-08-26 |

## Summary

docker-agent implements a three-tier event pipeline with deliberately different delivery semantics per tier. At the **producer tier**, `LocalRuntime.RunStream` returns a buffered channel (128) fed by a blocking `channelSink` so back-pressure propagates from consumer all the way back to the model-stream reader (`pkg/runtime/defaults.go:17`, `pkg/runtime/event_sink.go:35-38`). At the **process tier**, the TUI App runs a single fan-out goroutine that throttles (50 ms), merges adjacent deltas, and scatters to subscribers non-blockingly, preferring irrecoverable turn-boundary events over droppable content deltas on overflow (`pkg/app/app.go:956-1006`, `pkg/app/app.go:1591-1698`). At the **transport tier**, a per-session `eventLog` assigns monotonic sequence numbers, retains a 1024-event ring buffer for replay, drops (disconnects) slow SSE clients instead of ever blocking the pump, and closes with a sequenced terminal `session_exited` event that lets clients distinguish "session gone" from "connection dropped" (`pkg/server/eventlog.go:42-129`).

The result is that execution can be delayed by design at the runtime boundary (bounded back-pressure rather than unbounded buffering), but no observer — TUI subscriber or SSE client — can stall another observer or the agent loop itself. Reconnect recovers full semantics inside the buffer window via sequence-number replay and falls back to an explicit snapshot-plus-tail protocol (with a `gap` marker) when the resume point was evicted.

## Rating

**8/10**

Rationale: This is one of the most complete ordered-observation implementations in a Go agent codebase this reviewer has examined: per-session monotonic sequencing (`pkg/server/eventlog.go:91-93`), exactly-once replay/live registration under a single lock with a stress test proving it (`pkg/server/eventlog_test.go:152-180`), a documented and tested terminal-truth contract (`pkg/runtime/loop.go:187-200`, `pkg/runtime/elicitation_test.go:251-279`), priority-aware drop policy (`pkg/app/fanout_test.go:19-52`), and heartbeats to detect hung transports (`pkg/server/server.go:670-713`). Points withheld: (1) runtime observers run synchronously on the forwarder goroutine and can back-pressure the whole stream by contract (`pkg/runtime/observer.go:14-18`) with no timeout or isolation mechanism; (2) the remote client SSE consumer has no reconnect/resume logic at all — it is best-effort display data only (`pkg/runtime/client.go:496-499`); (3) ring eviction is an O(capacity) copy-down per append once full (`pkg/server/eventlog.go:96-101`).

## Evidence Collected

Every entry cites file path with line numbers.

| Area | Evidence | File:Line |
|------|----------|-----------|
| Bounded producer channel | Every runtime Event channel uses `defaultEventChannelCapacity = 128`, set "small enough to keep memory pressure bounded when consumers are slow" | pkg/runtime/defaults.go:9-17 |
| Blocking sink preserves back-pressure | `channelSink.Emit` sends blocking; send-on-closed recovered; doc states back-pressure is intentional | pkg/runtime/event_sink.go:18-38 |
| Non-blocking variant for outliving goroutines | `nonBlockingChannelSink` drops on full buffer; reserved for long-lived goroutines like the RAG watcher and teardown | pkg/runtime/event_sink.go:40-68 |
| Terminal truth: channel close, not StreamStopped | `finalizeEventChannel` emits StreamStopped best-effort/non-blocking (#3070 deadlock fix); close is authoritative | pkg/runtime/loop.go:177-211 |
| StreamStopped ordering pinned by test | `TestLocalRuntime_FinalizeEventChannelStreamStoppedIsLastBeforeClose`; also deadlock test with full buffer + gone consumer | pkg/runtime/elicitation_test.go:212-241, 251-279 |
| StreamStopped reason classification | Reason carries turnEndReason classification (normal/error/canceled) | pkg/runtime/event.go:580-606 |
| Observer contract: synchronous, back-pressuring | Observers invoked synchronously on forward goroutine; "a slow observer back-pressures both downstream observers and the consumer"; must fan out privately | pkg/runtime/observer.go:14-29, 66-81 |
| Model delta emission (high-volume producer) | Per-token `AgentChoice` / `AgentChoiceReasoning` / partial tool-call args emitted from stream loop; idle-timeout cancels stalled model streams | pkg/runtime/streaming.go:240-264, 312-335, 337-355 |
| Delta-only partial tool calls | Only newly received argument bytes sent, not re-transmitting accumulated payload each token | pkg/runtime/streaming.go:243-263 |
| App fan-out: single goroutine, non-blocking scatter | `startFanOut` clones subscriber list under lock, sends with select/default; slow subscriber drops | pkg/app/app.go:956-988 |
| Turn-boundary priority on overflow | `isTurnBoundaryEvent` (StreamStarted/Stopped/UserMessage/Error/Paused/SessionTitle) evicts oldest pending delta then retries once; still wedged → drop | pkg/app/app.go:963-1005 |
| Priority drop test | `TestFanOut_TurnBoundaryEventEvictsPendingDelta` and `TestFanOut_DroppableEventIsDroppedOnOverflow` | pkg/app/fanout_test.go:19-52, 57-86 |
| Throttling + coalescing of deltas | 50 ms throttle window; `mergeEvents` merges adjacent AgentChoice/Reasoning/PartialToolCall/ToolCallOutput runs, O(N) via strings.Builder | pkg/app/app.go:146, 1591-1698 |
| Subscriber lifecycle | `SubscribeWith` creates chan(1024), registers/unregisters under mutex, fan-out started once per App | pkg/app/app.go:911-944 |
| Blocking forward into app bus with ctx escape | `sendEvent` selects on channel vs ctx.Done to avoid blocking when TUI stopped reading; StreamStopped forwarded even after ctx cancel via `context.WithoutCancel` | pkg/app/app.go:697-704, 754-767 |
| Per-session sequenced log | `eventLog`: mutex-guarded `seq uint64`, ring buffer (default 1024), listeners map | pkg/server/eventlog.go:42-73 |
| Slow-client drop, never block pump | `appendLocked` delivers with select/default; full listener buffer → channel closed, listener removed | pkg/server/eventlog.go:91-111 |
| Exactly-once replay/live handoff | Backlog snapshot + listener registration under single lock; "every event is seen exactly once, in order"; proven by `TestEventLog_RegistrationIsGapless` (50 race iterations) | pkg/server/eventlog.go:146-160; pkg/server/eventlog_test.go:149-180 |
| Gap marker protocol | Resume point evicted → seq-0 `gapEvent` sent first, client re-snapshots then re-tails; `gapped()` detects since+1 < oldest buffered | pkg/server/eventlog.go:19-26, 204-235 |
| Terminal session_exited contract | Sequenced, replayed terminal event distinguishes permanent end from transport drop; appended by idempotent `close(reason)` | pkg/server/eventlog.go:28-40, 113-129 |
| Slow subscriber dropped test | `TestEventLog_SlowSubscriberIsDropped`: stalled listener removed after 20 appends into capacity-4 buffer | pkg/server/eventlog_test.go:247-291 |
| Late subscriber sees terminal | `TestEventLog_LateSubscriberSeesTerminalEvent`: post-close connect still replays session_exited | pkg/server/eventlog_test.go:230-245 |
| SSE reconnect cursor | `?since=` query param or standard `Last-Event-ID` header parsed by `parseSinceParam`; each frame tagged `id: <seq>` | pkg/server/server.go:652-662, 715-750 |
| Heartbeat liveness contract | ": ping" comment every 15 s (default), writes serialized with events under `writeMu` to prevent interleaved frames | pkg/server/server.go:39-43, 649, 686-713 |
| Snapshot advertises tail position | `GetSessionSnapshot` returns `LastEventSeq` so gap-recovering clients know where to resume tailing | pkg/server/session_manager.go:512-557 |
| Pump wiring for attached sessions | `RegisterEventSource` starts pump goroutine feeding `log.append`, buffers even with zero clients, tombstone-guarded against deleted sessions | pkg/server/session_manager.go:201-247, 300-313 |
| Lazy logs for server-created sessions | `ensureEventLog` creates ring-only log on first out-of-band event; cancel closes log delivering session_exited | pkg/server/session_manager.go:255-298 |
| Sub-session event forwarding | `runForwarding` drains child stream into parent sink, keeps draining after first ErrorEvent so transcript stays balanced | pkg/runtime/agent_delegation.go:340-412 |
| Session-scoped filtering support | `SessionScoped` interface lets observers/persistence drop sub-session events without re-implementing checks | pkg/runtime/event.go:18-23, 398-401 |
| Nested stream accounting in UI | TUI tracks `streamDepth`/`agentStack`; duplicate stops at depth 0 cannot desync slices | pkg/tui/page/chat/runtime_events.go:289-298, 347-364 |
| Attention-event queueing per tab | Supervisor queues PendingEvents in arrival order (single-slot overwrite bug #3584 fixed) | pkg/tui/service/supervisor/supervisor.go:22-31 |
| Remote SSE client bounds | Scanner raised to 16 MiB/line; scanner errors surfaced as Error events, not silent close | pkg/runtime/client.go:397-454; pkg/runtime/defaults.go:19-28 |
| Remote client lacks reconnect | `StreamSessionEvents` "closed when ctx cancelled, max duration reached, or server closes" — no resume cursor used | pkg/runtime/client.go:496-499 |
| Control-plane bridge | `registerAppEventSource` pumps App fan-out (runtime.Events only) into session eventLog for HTTP consumers | cmd/root/run_listen.go:52-62 |
| External event hooks never block bus | `--on-event` shell commands run async off `SubscribeWith`; failures logged, never block | cmd/root/run_event_hooks.go:33-60 |

## Answers to Dimension Questions

### Can one slow observer delay execution or other observers?

Yes at two layers, no at the others — each layer's answer is explicit and tested:

- **Runtime observers: yes, by documented contract.** Observers are invoked synchronously, in registration order, on the goroutine that forwards events to the consumer; the docs state plainly that "a slow observer therefore back-pressures both downstream observers and the consumer" and instruct long-running work to fan out to a private goroutine (`pkg/runtime/observer.go:14-18`, dispatch loop at `pkg/runtime/observer.go:71-79`). There is no timeout, panic isolation, or unregister mechanism — "the runtime cannot recover from a misbehaving observer" (`pkg/runtime/observer.go:20-24`).
- **Consumer channel: yes up to 128 events, then it stalls the producer — intentionally.** The runtime emits through a blocking sink (`pkg/runtime/event_sink.go:35-38`) into a 128-slot channel (`pkg/runtime/defaults.go:9-17`). A wedged TUI/API consumer eventually blocks the model-stream read loop; this trades unbounded memory for paused execution. The exception is teardown: `StreamStopped` emission is non-blocking so a gone consumer cannot deadlock shutdown (#3070, `pkg/runtime/loop.go:196-211`, test at `pkg/runtime/elicitation_test.go:212-241`).
- **App fan-out: no.** One slow subscriber gets events dropped (or its oldest pending delta evicted for boundary events); the fan-out goroutine never blocks on any subscriber (`pkg/app/app.go:963-985`). Other subscribers are unaffected because dispatch iterates a cloned list and each send is select/default.
- **SSE eventLog: no.** The pump never blocks on a client; a client whose per-listener buffer (capacity-sized) fills is dropped entirely and its stream ends, forcing reconnect-with-replay (`pkg/server/eventlog.go:103-110`, test `pkg/server/eventlog_test.go:247-291`). The pump goroutine feeding the log is likewise decoupled from any HTTP writer.

### What is bounded, dropped, coalesced, retained, or replayed under pressure?

- **Bounded:** runtime channel 128 (`pkg/runtime/defaults.go:17`); app bus 128 (`pkg/app/app.go:145`); throttle output 128 (`pkg/app/app.go:1593`); subscriber channels 1024 (`pkg/app/app.go:932`); eventLog ring 1024 with per-listener channel equal to capacity (`pkg/server/eventlog.go:11, 158`); SSE line size 16 MiB (`pkg/runtime/defaults.go:19-28`).
- **Dropped:** non-boundary messages on full subscriber buffers (`pkg/app/app.go:966-970`); boundary events only after evict-and-once-retry fails (`pkg/app/app.go:971-984`); entire listeners on SSE overflow (`pkg/server/eventlog.go:106-109`); StreamStopped itself if the consumer vanished (`pkg/runtime/loop.go:208-211`); events for deleted sessions (`pkg/server/session_manager.go:315-322`).
- **Coalesced:** consecutive text/reasoning deltas, partial tool-call argument chunks, and tool-output chunks merged within a 50 ms window into single events (`pkg/app/app.go:1641-1698`) — lossy for display only.
- **Retained/replayed:** the last 1024 sequenced events per session, replayed to reconnectors newer than their cursor (`pkg/server/eventlog.go:88-101, 185-202`). Durable facts beyond the window live in the session store; the snapshot endpoint hands back message state plus `LastEventSeq` (`pkg/server/session_manager.go:512-557`).

### How is exactly one terminal fact delivered after partial streaming?

Two cooperating contracts:

1. **In-process (RunStream consumers):** the guaranteed terminal signal is *channel close*, not the `StreamStoppedEvent`. StreamStopped is emitted best-effort before session-end hooks, pinned as the final event before close by test, but explicitly droppable under a full buffer (`pkg/runtime/loop.go:187-200`; `pkg/runtime/elicitation_test.go:243-279`; doc on the event type at `pkg/runtime/event.go:580-588`). Its `Reason` field classifies normal/error/canceled so consumers need not reverse-engineer cause.
2. **Transport (SSE):** the terminal fact *is* durable and sequenced. `eventLog.close(reason)` appends `session_exited` as a numbered event before disconnecting listeners, so live, reconnecting, and late-joining clients all see it exactly once via the replay machinery (`pkg/server/eventlog.go:113-129`; tests at `pkg/server/eventlog_test.go:182-228, 230-245`). The decision rule is documented on both the type and the handler: received `session_exited` → stop forever; stream closed *without* it → transport drop → reconnect with `Last-Event-ID` (`pkg/server/eventlog.go:31-36`; `pkg/server/server.go:664-674`).

Exactly-once across the replay/live boundary holds because backlog snapshotting and listener registration occur under one lock hold (`pkg/server/eventlog.go:149-160`), stress-tested with concurrent appends over 50 iterations (`pkg/server/eventlog_test.go:152-180`).

### Does reconnect recover semantics or merely resume best-effort display data?

It recovers **semantics** within the retention window and degrades honestly outside it:

- Sequence numbers are assigned under the same lock as buffering, so replay is gapless and ordered (`pkg/server/eventlog.go:91-111`); a caught-up client gets no spurious gap (`pkg/server/eventlog_test.go:128-147`).
- If the resume point was evicted, the client receives an explicit `gap` control event (seq 0, excluded from the stream's id space, `pkg/server/server.go:722-726`) telling it to re-snapshot and continue from the snapshot's `LastEventSeq` (`pkg/server/eventlog.go:19-23`; `pkg/server/session_manager.go:537, 554`). This is semantic recovery with declared loss, not silent best-effort resume.
- Liveness is separately recoverable: heartbeats let a client distinguish quiet-session from hung-transport and reconnect proactively (`pkg/server/server.go:670-674`).
- Caveat: the shipped Go remote client does not exercise any of this — its SSE reader neither sends `Last-Event-ID` nor reconnects; the stream simply ends (`pkg/runtime/client.go:496-499`). The reconnect protocol exists for external EventSource-style consumers; the built-in remote client treats the stream as display data and relies on fresh snapshots elsewhere. Also note the App-level fan-out sits *upstream* of the control-plane event source (`cmd/root/run_listen.go:52-62`), so events dropped by the App's lossy fan-out are unrecoverable by SSE replay — replay guarantees apply only to what reached the log.

## Architectural Decisions

1. **Blocking back-pressure at the runtime core, isolation at the edges.** The runtime→consumer hop uses blocking sends into a small buffer so a genuinely stuck consumer halts production instead of leaking memory (`pkg/runtime/event_sink.go:18-22`, `pkg/runtime/defaults.go:9-17`), while every multi-consumer hop (App fan-out, SSE log) is non-blocking with drop policies (`pkg/app/app.go:947-949`, `pkg/server/eventlog.go:75-79`). Execution is protected by boundedness, not by infinite buffering.
2. **Channel close as authoritative terminal, events as advisory.** Documented after deadlock #3070 and ordering debate #3074: best-effort typed stop event for UX immediacy, channel close as the guarantee (`pkg/runtime/loop.go:187-200`). The SSE side inverts this: there the terminal is a real sequenced event because the connection close is ambiguous.
3. **Class-based drop policy instead of uniform loss.** Events are partitioned into irreplaceable turn-boundary facts and superseded-able deltas; overflow sacrifices deltas first (`pkg/app/app.go:990-1006`). The same philosophy appears upstream in the runtime: partial tool calls ship only new argument bytes (`pkg/runtime/streaming.go:243-247`).
4. **Lossy coalescing at the presentation boundary only.** Throttle+merge lives in the App (TUI process), not in the runtime or the eventLog — the durable sequenced record stays complete; only UI fan-out collapses deltas (`pkg/app/app.go:1591-1598` vs `pkg/server/eventlog.go:91-101`).
5. **Snapshot-plus-cursor resync as the gap fallback.** Rather than growing the buffer unboundedly, the design accepts eviction and pairs a `gap` marker with a snapshot API that advertises its own tail seq (`pkg/server/eventlog.go:19-23`, `pkg/server/session_manager.go:537-554`).
6. **Tombstoned lifecycle for logs.** Deleted sessions can never have logs resurrected by stale closures; registration and teardown serialize under a dedicated leaf lock with defined lock order against `sm.mux` (`pkg/server/session_manager.go:72-79, 220-232, 300-313`).

## Notable Patterns

- **Single-lock exactly-once subscription:** the backlog-copy + register-listener critical section is the load-bearing trick for gapless reconnects (`pkg/server/eventlog.go:150-160`), validated by a repeated-race test (`pkg/server/eventlog_test.go:155-179`).
- **Evict-to-admit priority queues:** on overflow, drop-oldest-of-lower-class to admit higher class, retry once, then give up gracefully (`pkg/app/app.go:971-984`).
- **Per-connection control events outside the sequence space:** `gapEvent` rides the same pipe with seq 0 and no SSE `id:` so client cursors never ingest it (`pkg/server/eventlog.go:19-26`, `pkg/server/server.go:722-726`).
- **Write serialization for mixed-writer streams:** heartbeat goroutine and event sender share a `writeMu` so SSE frames never interleave mid-frame (`pkg/server/server.go:688-712`).
- **Capability interfaces for dedup:** `elicitationSinkMirror` lets consumers detect double-delivery paths and skip mirror copies, preventing exactly-twice delivery of elicitations (#3584) (`pkg/server/session_manager.go:338-349`, `pkg/app/app.go:713-739`).
- **Depth-counted nested streams in the UI:** `streamDepth`/`agentStack` pairing with pop-only-at-depth>0 guards against desync on duplicate stops (`pkg/tui/page/chat/runtime_events.go:347-364`).
- **Attention replay for unfocused tabs:** supervisor queues attention-demanding events per tab in arrival order rather than dropping or showing only the latest (`pkg/tui/service/supervisor/supervisor.go:22-31`).

## Tradeoffs

- **Back-pressure vs liveness at the runtime hop.** Blocking Emit means a dead consumer freezes the agent mid-turn (until context cancellation). Chosen over silent event loss where correctness matters; mitigated by non-blocking teardown (`pkg/runtime/event_sink.go:18-22`, `pkg/runtime/loop.go:196-211`).
- **Synchronous observers vs simplicity.** No observer isolation machinery (timeouts, panics-as-data, per-observer goroutines); the contract pushes concurrency onto observer authors (`pkg/runtime/observer.go:14-24`). Cheap fast-path, foot-gun for slow paths.
- **Ring buffer copy-down vs true circular indexing.** Eviction copies the surviving suffix (`copy(l.buf, ...)`), making steady-state append O(capacity) ≈ O(1024) memmove per event once full (`pkg/server/eventlog.go:96-101`). Simple and allocation-free, but a hot-path cost a head/tail-indexed ring would avoid.
- **Replay fidelity limited by fixed window.** 1024 events covers snapshot-to-connect races comfortably but not long disconnects; recovery outside the window requires a full snapshot round-trip, which re-downloads message state the client may already have (`pkg/server/eventlog.go:8-11`).
- **Upstream lossy fan-out undermines downstream replay.** Because the `--listen` control plane taps the App's throttled/dropping fan-out as its event source (`cmd/root/run_listen.go:52-62`), a dropped-at-fan-out event never reaches the sequenced log, so SSE replay cannot restore it. Boundary events are protected precisely because they matter for replay, but e.g. ToolCallEvents are droppable pre-log.
- **Built-in remote client forgoes the protocol it serves.** The server implements Last-Event-ID, gaps, heartbeats, and terminal markers; the Go client consumes none of them (`pkg/runtime/client.go:396-457, 496-499`).

## Failure Modes / Edge Cases

- **Wedged subscriber still loses boundary events** after evict-and-single-retry fails — logged, accepted, unrecoverable via that subscriber's stream (`pkg/app/app.go:971-984`, test asserts survival only in the recoverable case at `pkg/app/fanout_test.go:44-51`).
- **Slow observer deadlock-by-config:** a user-supplied observer doing synchronous network I/O stalls every event and ultimately the agent loop; nothing enforces the "fan out privately" guidance (`pkg/runtime/observer.go:16-18`).
- **Model stream stall** is bounded by an idle timeout that cancels the underlying HTTP request and surfaces an error rather than hanging forever (`pkg/runtime/streaming.go:342-355`).
- **Oversized SSE line** previously ended the run silently via `bufio.ErrTooLong`; now raised cap plus explicit Error event on scanner failure (`pkg/runtime/client.go:403-407, 446-453`).
- **Interleaved SSE writers corrupt framing** if heartbeats weren't serialized with payloads; handled under `writeMu` (`pkg/server/server.go:688-691`).
- **Post-deletion resurrection races:** stale elicitation sinks or late registrations cannot recreate logs for deleted sessions — creation/registration is tombstone-gated under `eventLogsMu`; appends to closed logs are no-ops (`pkg/server/session_manager.go:220-232, 263-272`; `pkg/server/eventlog.go:83-85`).
- **Duplicate StreamStopped at depth 0** cannot desync the TUI's stack tracking (pop guarded by depth check, `pkg/tui/page/chat/runtime_events.go:356-364`).
- **Concurrent title generation vs re-emission** snapshotted under session lock before the pump goroutine launches (`pkg/server/session_manager.go:1026-1049`).

## Future Considerations

- **Observer isolation:** a timeout- or goroutine-per-observer wrapper option would make the documented "must not block" advice enforceable rather than aspirational.
- **True circular buffer** for the eventLog to make full-ring appends O(1).
- **Remote client parity:** teaching the Go client to send `Last-Event-ID`/`since`, honor `session_exited`/`gap`, and auto-reconnect would upgrade remote sessions from best-effort display to the same semantic recovery external SSE clients get.
- **Pre-loss fan-out tap:** sourcing the control-plane event log from the runtime observer chain (or the RunStream channel) instead of the App's lossy fan-out would extend replay guarantees to all event classes.
- **Back-pressure telemetry:** drops currently surface only as warnings (`pkg/app/app.go:968, 982`); counters exposed via status endpoints would make sustained-slow-consumer conditions observable operationally.

## Questions / Gaps

- **No evidence found for adaptive/size-aware retention:** the eventLog window counts events, not bytes; whether a session emitting huge tool outputs fits 1024 events within memory expectations was not analyzed (no byte-budget code located in `pkg/server/eventlog.go`).
- **No evidence found for ACP/a2a reconnect cursors:** this study scoped the HTTP/SSE and in-process paths; whether the ACP front-end (`pkg/acp/agent.go`) offers equivalent resume semantics was not traced and is unverified here.
- **Unquantified throttle latency impact:** the fixed 50 ms window (`pkg/app/app.go:146`) trades latency for merge efficiency; no measurement or config knob for it was found in the source.
- **Turn-boundary completeness unverified for future event types:** `isTurnBoundaryEvent` enumerates six concrete types (`pkg/app/app.go:994-1005`); a new lifecycle event added to `pkg/runtime/event.go` would silently default to droppable unless the switch is updated — no compile-time linkage between the two lists was found.

---

Generated by dimension `01.04 Ordered Observation, Live Streaming, and Backpressure` against `docker-agent`.
