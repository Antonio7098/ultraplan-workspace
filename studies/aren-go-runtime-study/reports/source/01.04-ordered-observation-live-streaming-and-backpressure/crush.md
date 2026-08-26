# Source Analysis: crush

## 01.04 Ordered Observation, Live Streaming, and Backpressure

### Source Info

| Field | Value |
|-------|-------|
| Name | crush |
| Path | `studies/aren-go-runtime-study/sources/crush` |
| Language / Stack | Go (Charm stack: bubbletea v2 TUI, fantasy LLM abstraction, SQLite/sqlc persistence, HTTP SSE client/server split) |
| Analyzed | 2026-08-26 |

## Summary

Crush implements ordered observation through a layered pipeline: per-domain generic pub/sub brokers (`internal/pubsub/broker.go`) feed a single fan-in broker of `tea.Msg` in `internal/app/app.go`, which feeds either the local TUI via non-blocking `program.Send` or an SSE endpoint (`GET /v1/workspaces/{id}/events`) for remote clients. The design explicitly separates two delivery classes: **lossy, non-blocking** `Publish` for high-frequency streaming state (token deltas), and **bounded-blocking** `PublishMustDeliver` (50 ms per-subscriber timeout) for terminal facts (`RunComplete`, tool results, permission requests). Deltas are further coalesced by a write-behind buffer in the message service (33 ms debounce) that classifies each snapshot as structural/terminal vs. progress via `shouldFlushNow`, flushing terminal state synchronously and publishing it on the must-deliver path.

Exactly-once terminal truth is engineered rather than assumed: the coordinator coalesces retry-attempt completions into a single `RunComplete` published once via must-deliver; the agent's deferred publish flushes all debounced deltas *before* emitting `RunComplete` and embeds the final text in the event so clients can reconcile out-of-order delivery; a dispatcher fallback emits an errored `RunComplete` only when the coordinator died before publishing, deduplicated by a context marker. Slow observers are isolated by bounded per-subscriber buffers with drop counters; a stalled SSE client loses events for itself only. Reconnect is deliberately best-effort: there are no event IDs, no cursors, and no replay buffer — clients reconnect with capped exponential backoff and recover semantics by re-fetching full state, never by resuming the stream. Observers cannot control execution because runs are bound to the workspace context, not to any HTTP request or subscriber lifetime.

## Rating

**8/10**

Rationale: This is an unusually disciplined implementation of exactly this dimension's concerns. The dual publish semantics are documented at the code level and enforced by type-checked interfaces (`pubsub.Publisher[T]`, `internal/pubsub/events.go:63-66`); drop accounting exists for both classes; the durable-fact vs. lossy-delta boundary is explicit in `shouldFlushNow` (`internal/message/message.go:419-447`); terminal-event guarantees are covered by end-to-end tests including out-of-order reconciliation (`internal/cmd/run_stream_test.go:89-103`) and cross-client cancellation (`internal/server/e2e_agent_test.go:324-375`). The score is held below 9–10 because (a) reconnect recovers by full resync only — dropped events during a stream gap are unaccounted to the observer, and the documented policy is "gone" (`internal/workspace/workspace.go:44-48`); (b) no global sequence numbers exist, so cross-broker ordering is best-effort by construction (acknowledged at `internal/agent/agent.go:747-748`); and (c) the drop counters are exported but never consumed outside tests, so saturation is observable only in logs.

## Evidence Collected

Every entry includes a file path with line numbers.

| Area | Evidence | File:Line |
|------|----------|-----------|
| Broker core | Generic `Broker[T]` with map-of-channels fan-out under `sync.RWMutex`; per-subscriber buffered channel sized from `channelBufferSize` | internal/pubsub/broker.go:48-70 |
| Buffer bound | `bufferSize = 4096`, documented as sized to cover one full streamed assistant turn even under TUI render stalls | internal/pubsub/broker.go:34-40 |
| Lossy publish | Non-blocking send; on full buffer drops per-subscriber, increments counter, logs warning | internal/pubsub/broker.go:165-188 |
| Must-deliver publish | Fast-path non-blocking, slow-path blocking send bounded by per-subscriber timeout; drop counted + error logged; publisher never blocks indefinitely | internal/pubsub/broker.go:201-236 |
| Timeout bound | `defaultMustDeliverTimeout = 50ms`, overridable via `SetMustDeliverTimeout` (tests) | internal/pubsub/broker.go:42-45,72-83 |
| Drop counters | `DropCount()` / `MustDeliverDropCount()` exposed "so callers can surface saturation in telemetry" | internal/pubsub/broker.go:22-23,146-157 |
| Subscriber removal | Context-cancel goroutine deletes + closes sub channel under lock; `Shutdown` closes all subs and marks broker done | internal/pubsub/broker.go:85-102,120-135 |
| Event envelope | `Event[T]{Type,Payload}` lifecycle verbs created/updated/deleted; discriminated wire envelope `Payload{Type,PayloadType,json.RawMessage}` | internal/pubsub/events.go:8-12,35-54 |
| Publisher interface | `Publisher[T]` bakes both semantics into the contract: lossy `Publish` + bounded `PublishMustDeliver` | internal/pubsub/events.go:56-67 |
| Fan-in wiring | `setupEvents` fans ~12 domain brokers into one `app.events` broker of `tea.Msg`; terminal-class services use the must-deliver variant | internal/app/app.go:585-609 |
| Fan-in lossy hop | `setupSubscriber` republishes upstream events with lossy `Publish` | internal/app/app.go:611-634 |
| Fan-in must-deliver hop | `setupSubscriberMustDeliver` used for permissions/questions/run-completions; comment explains dropping RunComplete would hang `crush run` | internal/app/app.go:636-666 |
| Delta coalescing | Message service buffers updates behind 33ms debounce; terminal updates bypass debounce and flush synchronously | internal/message/message.go:16-21,215-274 |
| Terminal classification | `shouldFlushNow`: message finished, tool-call set grew, tool call finished, or reasoning ended ⇒ structural flush | internal/message/message.go:414-447 |
| Terminal publish path | Flushed snapshots classified against last-flushed baseline; terminal uses `PublishMustDeliver`, progress uses lossy `Publish` | internal/message/message.go:376-383 |
| Write serialization | Per-ID `flushing` flag serializes concurrent flushers "so SQL writes never reorder"; sync caller spins on 1ms yields until in-flight write lands | internal/message/message.go:71-102,304-331 |
| Detached timer flush | Debounce timer flush runs on `context.Background()` so a canceled stream cannot strand buffered writes | internal/message/message.go:264-271 |
| Streaming callbacks | `OnTextDelta`/`OnReasoningDelta` append to the message and call `messages.Update` per token — producer never touches subscribers directly | internal/agent/agent.go:880-918 |
| RunComplete emit point | Deferred publish after run exit: flush-all first (5s bounded, detached ctx), then build `RunComplete{MessageID,Text}`; Text embedded to reconcile out-of-order observers | internal/agent/agent.go:735-781 |
| Exactly-once coordinator | Per-run `OnComplete` closure overwrites `latest`; publishes once via `PublishMustDeliver` after retries resolve; sets `MarkRunCompletePublished` | internal/agent/coordinator.go:272-336 |
| Dispatcher fallback | `runAgent` emits errored RunComplete only when `msg.RunID != "" && !RunCompletePublished(ctx)`; cancel errors produce none | internal/backend/agent.go:87-122 |
| Marker plumbing | `MarkRunCompletePublished`/`RunCompletePublished` context marker prevents duplicate terminal events | internal/agent/run_marker.go:27-47 |
| Cancel-on-entry terminal | Immediately-canceled accepted run still persists turn and publishes cancelled RunComplete before installing the defer | internal/agent/agent.go:591-618 |
| Queued-drop terminal | Cancel clears queue and emits exactly-one cancelled RunComplete per dropped RunID prompt (tested end-to-end) | internal/agent/run_complete_test.go:255-284,293-321 |
| Fire-and-forget dispatch | `SendMessage` validates, accepts, spawns goroutine bound to workspace ctx, returns immediately (HTTP 202) | internal/backend/agent.go:30-61 |
| Observer ≠ owner | Run lifetime owned by workspace, not caller request; e2e test proves expired POST context doesn't kill run | internal/server/e2e_agent_test.go:532-570 |
| SSE subscribe-before-attach | Handler subscribes to broker before `AttachClient` so presence implies delivery; headers flushed immediately so quiet workspaces don't block the client handshake | internal/server/proto.go:281-337 |
| SSE wrapping | `wrapEvent` converts typed events to discriminated envelopes; unrecognized/unrepresentable types dropped (logged), never fabricated | internal/server/events.go:27-170,185-200 |
| Client SSE reader | Buffered chan(100); parse loop decodes envelopes by `PayloadType`; malformed lines logged and skipped | internal/client/proto.go:113-199 |
| Client send backpressure | `sendEvent` blocks only on ctx.Done — a slow consumer stalls its own reader goroutine, not the server | internal/client/proto.go:268-278 |
| Reconnect loop | Capped exponential backoff 250ms→10s; 404 triggers workspace re-registration; any stream close treated as degraded requiring resync | internal/workspace/client_workspace.go:71-75,827-896 |
| Workspace recovery | `recoverWorkspace` re-creates workspace from cached snapshot, adopts new ID, re-inits coder agent | internal/workspace/client_workspace.go:899-942 |
| Resync-not-replay | `ErrStreamClosed` doc: events published while away "are lost for good"; recovery re-asserts session + tells UI to reload | internal/workspace/workspace.go:44-49; internal/workspace/client_workspace.go:961-977 |
| UI resync rationale | "events published while the stream was down are gone" → recovered connection reloads open session | internal/ui/model/ui.go:1509-1512 |
| TUI consumption | Local app subscribes and forwards via `program.Send` (non-blocking); panic recovery quits gracefully | internal/app/app.go:706-736 |
| CLI terminal wait | `crush run` exits only on RunComplete matching its own minted RunID; suppresses live message events; reconciles stdout from embedded Text if out of order | internal/cmd/run.go:258-333; internal/cmd/run_stream_test.go:83-132 |
| Multi-client fan-out | E2E: same message reaches two independent SSE clients; killing A's stream does not break B | internal/server/e2e_test.go:335-425,516-560 |
| Cross-client cancel semantics | Cancel by client B surfaces FinishReasonCanceled message to A with no AgentEvent error, tested on both streams | internal/server/e2e_agent_test.go:319-375 |
| Ordering convergence test | Coalesced vs zero-debounce update sequences converge to identical final DB state | internal/message/message_test.go:272-315 |
| Backpressure unit tests | Saturated subscriber increments DropCount; must-deliver honors timeout then drops; active reader receives all 10 events with zero drops | internal/message/message_test.go:335-422 |
| Terminal-over-lossy priority | Test proves must-deliver publish overtakes a saturated channel where lossy publish would drop | internal/message/message_test.go:424-461; internal/agent/run_complete_test.go:200-221 |

## Answers to Dimension Questions

### Can one slow observer delay execution or other observers?

No for execution; mostly no for other observers, with one bounded exception.

- Execution is structurally immune: agent runs execute on goroutines bound to the workspace context (`internal/backend/agent.go:59,91-95`), stream provider deltas into the message service's in-memory buffer (`internal/agent/agent.go:884-917`), and never read from any subscriber channel. An SSE client that stops reading TCP merely stops draining its own subscription; the server handler blocks on its channel read (`internal/server/proto.go:314-335`), the 4096-slot buffer fills, and subsequent publishes drop for that subscriber only (`internal/pubsub/broker.go:177-187`). E2E test `TestE2E_KillingClientASSEDoesNotBreakClientB` (`internal/server/e2e_test.go:519-560`) proves client isolation.
- Other observers: `Publish` iterates subscribers under `RLock` doing only non-blocking sends, so a slow subscriber costs nothing beyond a map iteration. The exception is `PublishMustDeliver`, which can block up to 50 ms per slow subscriber per event while holding the read lock (`internal/pubsub/broker.go:214-235`). Because the lock is a read lock, other publishers proceed concurrently; the worst case is added latency on that publish call, not a stall. The timeout bounds total head-of-line delay.
- The TUI consumes via `program.Send`, which is asynchronous (`internal/app/app.go:733`), so render slowness does not back-propagate into the fan-in goroutine.

### What is bounded, dropped, coalesced, retained, or replayed under pressure?

- **Bounded**: per-subscriber channel capacity 4096 (`internal/pubsub/broker.go:40`); must-deliver wait 50 ms/subscriber/event (`broker.go:45`); client-side SSE decode buffer 100 events (`internal/client/proto.go:115`); post-run delta flush 5 s (`internal/agent/agent.go:754`); debounce window 33 ms (`internal/message/message.go:21`).
- **Dropped**: any event hitting a full subscriber channel under lossy `Publish` (counted, `internal/pubsub/broker.go:184-186`); must-deliver events after the timeout (counted, `broker.go:228-230`); unsupported MCP event types at the SSE boundary rather than being mis-typed (`internal/server/events.go:41-48`); pending coalesced state on message delete (`internal/message/message.go:147-156`).
- **Coalesced**: streaming text/reasoning/tool-input deltas within the debounce window collapse into one SQL write + one pubsub event (`internal/message/message.go:16-21,261-273`); retry-attempt RunCompletions coalesce to the final outcome (`internal/agent/coordinator.go:281-288`).
- **Retained**: authoritative conversation state lives in SQLite, not in the event log — the durability story is "state store + lossy stream," so nothing event-shaped is retained for replay.
- **Replayed**: nothing. There is no replay mechanism anywhere in the stream path; recovery is a full state refetch (see next answer).

### How is exactly one terminal fact delivered after partial streaming?

Through three cooperating mechanisms:

1. **Single emit point**: every exit branch of `sessionAgent.Run` funnels through `publishRunComplete` (`internal/agent/agent.go:530-548`), including the cancel-on-entry early return that fires before the deferred publisher is installed (`agent.go:591-618`).
2. **Exactly-once across retries and layers**: the coordinator replaces per-attempt emissions with a coalesced final publish and stamps a context marker (`internal/agent/coordinator.go:281-336`); the backend dispatcher consults `RunCompletePublished(ctx)` before emitting its own fallback for pre-run failures (`internal/backend/agent.go:109-121`, marker at `internal/agent/run_marker.go:27-47`). Tests assert exactly one cancelled RunComplete and no second event (`internal/agent/run_complete_test.go:228-245`).
3. **Ordering and reconciliation**: the deferred publisher drains all debounced deltas (bounded 5 s flush) *before* publishing RunComplete, and embeds `MessageID` + full `Text` in the payload so clients that observe events out of order (the fan-in does not serialize across brokers — acknowledged at `internal/agent/agent.go:747-748`) reconstruct the final answer. `crush run` demonstrates the consumer side: it exits only on a RunID-matching RunComplete and prints the embedded text when message events lag (`internal/cmd/run.go:319-333`; `internal/cmd/run_stream_test.go:89-103`). Delivery itself is must-deliver, but the design honestly documents that even must-deliver can drop after 50 ms, making recovery the subscriber's job (`internal/pubsub/broker.go:190-200`) — the embedded-text payload is that recovery.

### Does reconnect recover semantics or merely resume best-effort display data?

It recovers semantics by abandoning the stream entirely: reconnect is best-effort transport, semantic recovery is explicit resync. There are no SSE event IDs, no `Last-Event-ID` handling, no cursor parameter, and no replay buffer — searched for `Last-Event-ID`, `cursor`, `replay`, and sequence fields across `internal/server`, `internal/client`, `internal/proto`, and `internal/pubsub`; only per-test monotonic run IDs exist (`internal/server/e2e_agent_test.go:55-64`). On any stream close, the client marks the link degraded, retries with capped exponential backoff (250 ms→10 s), re-registers the whole workspace on 404 (`internal/workspace/client_workspace.go:821-896,899-942`), and on success re-asserts its current session and sends `ConnectionRecovered`, which makes the UI reload the open session because "events published while the stream was down are gone" (`internal/ui/model/ui.go:1509-1512`; `client_workspace.go:961-977`). So the system's answer is deliberate: the event stream carries only display-grade deltas, and correctness always comes from the SQLite state plus a refetch — reconnect never fabricates continuity.

## Architectural Decisions

1. **Two-tier delivery semantics baked into the publisher interface** (`internal/pubsub/events.go:63-66`). Rather than a configurable knob, the lossy/must-deliver distinction is part of the `Publisher[T]` contract with package-level documentation explaining why (`internal/pubsub/broker.go:4-24`). Every call site declares which class its event belongs to, making the durable-vs-lossy taxonomy greppable.
2. **State store as source of truth, stream as notification** (`internal/message/message.go:33-45`; `internal/workspace/workspace.go:44-48`). Events are advisory; clients reconcile by refetching. This removes the need for sequencing, replay, and cursor machinery entirely — a major simplification versus a durable event log, paid for in resync cost after every disconnect.
3. **Central fan-in with per-service delivery class** (`internal/app/app.go:585-609`). Twelve domain brokers merge into one `tea.Msg` broker; the choice of lossy vs. must-deliver fan-in per service encodes which domains carry terminal facts (permissions, questions, run-completions) vs. progress (sessions, messages, history, LSP, MCP).
4. **Terminal truth via correlators + markers, not ordering** (`internal/agent/coordinator.go:289-336`; `internal/backend/agent.go:92-121`). RunIDs correlate terminal events to specific waits; a context marker deduplicates fallback emission. Correctness does not depend on the stream preserving order.
5. **Runs owned by workspace, not by requests** (`internal/backend/agent.go:16-61`). Dispatch is fire-and-forget with accept reservations (`BeginAccepted`/`AcceptedRun`) so cancellation racing dispatch remains coherent (`internal/agent/agent.go:580-650`).

## Notable Patterns

- **Write-behind coalescing with structural preemption**: the debounce buffer treats terminal transitions as barriers that stop the timer and flush synchronously (`internal/message/message.go:250-259`), with a baseline comparison against `lastFlushed` so classification is stable even when the flush itself races new deltas (`message.go:343-351`).
- **Flush-before-terminal-publish ordering**: the agent guarantees deltas precede RunComplete in the buffer, then embeds the payload anyway for clients that see inversion (`internal/agent/agent.go:735-748`) — belt-and-suspenders with an honest comment about why.
- **Subscribe-before-attach handshake** on the SSE endpoint so client presence (used for busy-state computation) implies event delivery has started (`internal/server/proto.go:288-302`), plus immediate header flush to unblock the client's initial RoundTrip (`proto.go:307-312`).
- **Drop observability by construction**: both drop paths increment atomic counters and log at distinct severities (warn for lossy, error for must-deliver) (`internal/pubsub/broker.go:184-186,228-230`).
- **Detached contexts for must-land writes**: canceled-turn persistence and the final flush use `context.WithoutCancel` + short timeouts so shutdown can't silently eat a turn's tail (`internal/agent/agent.go:508-510,750-755`; timer flush at `internal/message/message.go:266-270`).

## Tradeoffs

- **Lossiness as a feature vs. silent gaps**: dropping intermediate deltas is correct for rendering, but a dropped *structural* event that misclassifies as progress would leave stale UI until the next flush. Mitigation: terminal detection is conservative (any tool-call growth/finish flushes immediately, `internal/message/message.go:419-447`).
- **No sequencing vs. simplicity**: per-channel FIFO gives per-subscriber order, but the multi-broker fan-in means interleaving across services is undefined. The team accepted this and compensated with payload self-sufficiency (embedded text, finish reasons) instead of adding sequence numbers.
- **Bounded blocking vs. guaranteed delivery**: must-deliver can still drop after 50 ms; the code says so plainly and pushes recovery onto subscribers (`internal/pubsub/broker.go:198-200`). For `crush run` this is acceptable because the waiter also holds the run's result via state; a pure-log-semantics consumer would find this insufficient.
- **Resync-on-reconnect vs. bandwidth**: every transient blip causes a full session reload (`internal/ui/model/ui.go:1509-1512`). Simple and correct, potentially heavy on long sessions over flaky links.
- **Global LSP broker**: `lspBroker` is a package-level singleton (`internal/app/lsp_events.go:41`), trading test isolation and multi-workspace scoping for convenience.

## Failure Modes / Edge Cases

- **Slow SSE client event starvation**: a client that stops reading accumulates up to 4096 queued events server-side, then loses everything newer; since there is no gap signal, the client cannot know it missed events unless the stream also closes. A mid-stream stall without disconnect is therefore invisible to the client. (Mechanism: `internal/pubsub/broker.go:177-187`; no gap detection exists in `internal/client/proto.go:132-199`.)
- **Must-deliver timeout under sustained saturation**: if all subscribers stall >50 ms, RunComplete drops after logging; a waiting `crush run` would hang. The dispatcher's fallback covers pre-run failures but not a dropped post-run RunComplete — recovery relies on the documented-but-unimplemented "re-fetch on the next session-visible event" (`internal/pubsub/broker.go:198-200`).
- **Fan-in amplification**: a must-deliver chain crosses two brokers (service → app.events), each with its own timeout; worst case adds latency and a second drop opportunity per hop (`internal/app/app.go:643-666`).
- **Cancel/dispatch races**: extensively enumerated and handled — cancel-on-entry, cancel between active-set and assistant creation, cancel of queued prompts, stale cancel marks cleared on normal completion (`internal/agent/agent.go:591-650`; tests at `internal/agent/dispatch_cancel_test.go:73-370`, `internal/server/e2e_agent_test.go:377-530`).
- **Delete during debounce window**: pending buffered state is dropped so a deleted row is never resurrected (`internal/message/message.go:147-156`, tested at `internal/message/message_test.go:317-333`).
- **Retry content reset**: on provider retry, streamed partial content is reset so the retried response doesn't concatenate with the failed attempt's fragments (`internal/agent/agent.go:931-941`).

## Future Considerations

- **Wire the drop counters into telemetry**: `DropCount`/`MustDeliverDropCount` have no production readers (searched all non-test call sites); surfacing them in `/v1/health` or logs-on-threshold would make saturation actionable instead of log-noise-only (`internal/pubsub/broker.go:146-157`).
- **Cheap gap signaling**: an occasional monotonically increasing sequence number in the SSE envelope (or a periodic heartbeat carrying it) would let clients detect mid-stream loss without a full replay protocol — a middle ground between today's "resync everything" and a durable log.
- **Scoped LSP broker**: moving the package-level `lspBroker` into per-app instances would align it with the rest of the broker lifecycle management (`internal/app/lsp_events.go:41`).
- **Pending-map pruning in message service**: entries in `service.pending` persist for the process lifetime after their final flush; harmless per entry, but a long-running server accrues one map entry per message ever updated (`internal/message/message.go:104-111`).

## Questions / Gaps

- **Is mid-stream (non-closing) event loss detectable by any client?** No evidence found. Searches for `Last-Event-ID`, `event id`, `cursor`, `gap`, and sequence numbering across `internal/server`, `internal/client`, `internal/proto`, and `internal/workspace` returned nothing beyond per-test run IDs. The only loss trigger for resync is full stream closure (`internal/workspace/workspace.go:44-48`).
- **What consumes drop counters operationally?** No evidence found outside tests; the doc comment's promise of telemetry surfacing (`internal/pubsub/broker.go:22-23`) appears unrealized in this tree.
- **Does any consumer rely on cross-service event ordering?** No evidence found of ordering assumptions; the codebase explicitly disclaims it (`internal/agent/agent.go:746-748`) and compensates with payload reconciliation, suggesting the question is settled by design rather than by mechanism.
- Note on scope: dimension citations reference Aren Phase 1/Phase 4 docs outside this source directory; per source-isolation rules those files were not accessed. All evidence above is drawn solely from `studies/aren-go-runtime-study/sources/crush`.

---

Generated by `01.04 Ordered Observation, Live Streaming, and Backpressure` against `crush`.
