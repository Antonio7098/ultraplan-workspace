# Source Analysis: crush

## 01.01 Lifecycle Transition Ownership and Terminal Arbitration

### Source Info

| Field | Value |
|-------|-------|
| Name | crush |
| Path | `studies/aren-go-runtime-study/sources/crush` |
| Language / Stack | Go 1.26 (`github.com/charmbracelet/crush`, v0.91.1); fantasy LLM abstraction, Bubble Tea v2 TUI, SQLite/sqlc persistence, in-process pubsub brokers |
| Analyzed | 2026-08-26 |

Citation convention: all paths below are workspace-relative under the source root `studies/aren-go-runtime-study/sources/crush/`. Example: `internal/agent/agent.go:566` denotes `studies/aren-go-runtime-study/sources/crush/internal/agent/agent.go:566`.

### Isolation note (Aren citations)

The dimension lists two Aren requirement documents as citations (`Aren/docs/phase-1-prd/02-lifecycle-requirements.md`, `01-product-definition-and-scope.md`). These live outside the selected source directory and the sandbox permission policy denies external-directory reads, so they were not inspected. Aren's contract is therefore framed only from the injected dimension text itself; every code claim below traces exclusively into crush.

## Summary

Crush gives one component — `sessionAgent` (`internal/agent/agent.go:169`) — ownership of the accepted → (cancel-on-entry | queued | active) dispatch decision and of terminal publication for a session run. There is no explicit run state machine: "running" is the presence of a per-session entry in an `activeRequests` map (`internal/agent/agent.go:185,2048-2051`), and "terminal" is a single `notify.RunComplete` event plus a persisted `FinishReason*` on the assistant message. The dispatch decision is serialized by a lazily-created per-session mutex (`dispatchMu`, `internal/agent/agent.go:191,588-589`) held only across the handoff, never across DB or LLM I/O, which makes two concurrent runs on one session impossible (proven by `TestRun_ConcurrentInProcessDispatchStartsOneRun`, `internal/agent/dispatch_race_test.go:113-155`).

Terminal publication is funneled through one emit path, `publishRunComplete` (`internal/agent/agent.go:539-548`), reached from four producers that are each guarded to fire at most once per RunID: the streaming defer (`internal/agent/agent.go:749-781`), the cancel-on-entry early return (`internal/agent/agent.go:591-618`), the queue-handoff branch (`internal/agent/agent.go:1282-1325`), and two upstream dedup layers — the coordinator's retry-coalescing hook (`internal/agent/coordinator.go:281-336`) and the dispatcher's marker-gated fallback (`internal/backend/agent.go:109-121` with `internal/agent/run_marker.go:31-52`). Cancellation requests are represented separately from confirmed termination: `Cancel` invokes the context cancel func but deliberately keeps the `activeRequests` entry until the goroutine finishes (`internal/agent/agent.go:1966-1973`), records a sequence-bounded pending-cancel mark for accepted-but-not-yet-active runs (`internal/agent/agent.go:1988-2003`), and termination is confirmed only when `Run` returns, persists `FinishReasonCanceled` (`internal/agent/agent.go:501-528,1159-1160`), and derives `RunComplete.Cancelled` from the actual error (`internal/agent/agent.go:767-772`). The design is unusually test-rich: every arbitration guarantee named above has a dedicated regression test that observes published events, not internal bookkeeping.

The main gaps: panic containment is boundary-specific (an unrecovered panic inside stream callbacks would crash the process), the `Summarize` path registers its active request outside `dispatchMu` and contains a vestigial `sessionID+"-summarize"` cancel-key lookup that can never match, and `PublishMustDeliver` still permits drops after a 50 ms per-subscriber timeout.

## Rating

**8 / 10**

Rationale against the dimension's four questions:

1. *Conflicting terminal outcomes* — strongly prevented and tested (per-session mutex dispatch gate, exactly-once emit path, layered dedup). Deducted points because correctness rests on a chain of cooperating guards spread across three packages rather than one enforced choke point.
2. *Cancellation vs termination* — cleanly separated with three distinct representations (context cancel, pending mark, confirmed outcome) and regression tests for poison-prevention (`internal/agent/dispatch_cancel_test.go:326-447`).
3. *Coherent visibility of state/outcome/timing/event* — flush-before-publish ordering plus embedded final text reconcile out-of-order observers (`internal/agent/agent.go:735-781`; consumer side `internal/cmd/run.go:384-429`), but cross-broker ordering is not guaranteed and must-deliver can still drop.
4. *Adoptable size* — the core arbitration kit is a few hundred lines of dependency-free primitives; however several supporting behaviors are entangled with crush-specific concerns (queue folding, provider auth retries).

## Evidence Collected

Every entry cites files under `studies/aren-go-runtime-study/sources/crush/`.

| Area | Evidence | File:Line |
|------|----------|-----------|
| Run creation (fire-and-forget dispatch) | `SendMessage` validates, calls `BeginAccepted` synchronously, adds a `runWG` ticket under `runMu` with a closing check, then `go b.runAgent(...)`; errors reach observers via events, not return value | internal/backend/agent.go:30-61 |
| Run creation → running | Coordinator `run()` gates on readiness (`readyWg.Wait`), refreshes models, threads caller RunID from context onto the call | internal/agent/coordinator.go:223-313 |
| Running state representation | `activeRequests *csync.Map[string,*activeCancel]`; `IsSessionBusy` = key present; no status enum anywhere on `Session` | internal/agent/agent.go:185, 2048-2051; internal/session/session.go:50-63 |
| Dispatch decision gate | Per-session `dispatchMu` serializes accepted → (cancel-on-entry \| queued \| active); busy-check + active-registration made atomic; comment explains two callers cannot both pass | internal/agent/agent.go:184-222, 580-650 |
| Accept reservation handle | `AcceptedRun` owns one reservation; `Close` idempotent via `atomic.Bool` CAS; `acceptedRuns>0` means dispatched-but-not-active; separate `acceptedMu` avoids lock re-entry deadlock | internal/agent/agent.go:264-337 |
| Cancel high-water mark | Monotonic `acceptSeqGen`; Cancel raises per-session `cancelMark` to latest seq so one cancel covers all then-accepted handles while later accepts are never poisoned | internal/agent/agent.go:1981-2003, 217-222 |
| Terminal publication owner | Single emit path `publishRunComplete`: honors per-call `OnComplete` hook else falls back to `PublishMustDeliver` on the RunComplete broker; "the single emit path shared by the streaming defer and the cancel-on-entry early return" | internal/agent/agent.go:530-548 |
| Coherent commit (state+outcome+timing+event) | Defer first does `FlushAll` of debounced message writes on a detached bounded context, then publishes RunComplete carrying final MessageID+Text; comment documents flush-then-publish order and out-of-order reconciliation | internal/agent/agent.go:729-781 |
| Cancelled outcome derivation | `complete.Cancelled = errors.Is(retErr, context.Canceled)` or `ctx.Err() != nil`; Error and Cancelled may overlap when cancel triggers downstream error | internal/agent/agent.go:767-772; internal/agent/notify/notify.go:52-78 |
| Persisted terminal marker | `persistCanceledTurn` writes assistant message with `FinishReasonCanceled` using `context.WithoutCancel` + 5s timeout so shutdown cannot drop it | internal/agent/agent.go:501-528 |
| ABA-safe cleanup | Deferred `activeRequests.CompareAndDelete(sessionID, ac)` removes only the finishing run's own entry; primitive documented in csync | internal/agent/agent.go:160-167, 652-657; internal/csync/maps.go:63-80 |
| Busy-state ordering vs notify | Active request deleted before agent-finished notification so TUI pollers never see stale busy state | internal/agent/agent.go:1209-1224 |
| Queue handoff atomicity | Handoff to queued prompt happens under `dispatchMu`: pending-cancel filtering, fresh accept reservation before unlock, recursion into `Run` closes dequeue→re-register window | internal/agent/agent.go:1226-1326 |
| Duplicate-suppression (recursion) | `skipRunComplete` flag stops outer defer from publishing when handing off; outer turn publishes its own RunComplete iff its RunID differs from all queued prompts | internal/agent/agent.go:1282-1325 |
| Hook stripping on enqueue | `enqueueCall` strips `OnComplete` (caller's coalesce scope has ended) and `Accepted`, preserves `acceptSeq` and RunID | internal/agent/agent.go:356-379 |
| Queue drain atomicity | `drainQueueForStep` evaluates cancel coverage under the dispatch mutex; canceled RunID-bearing drops returned for explicit terminal publish | internal/agent/agent.go:381-424 |
| Dropped-queued-prompt terminals | `publishCanceledQueueDrops` emits cancelled RunComplete for dropped prompts carrying a RunID using detached bounded context | internal/agent/agent.go:426-471 |
| Coordinator retry coalescing | Per-attempt `OnComplete` overwrites `latest`; publish exactly once after retries resolve; `MarkRunCompletePublished(ctx)` signals dispatcher not to duplicate | internal/agent/coordinator.go:271-336 |
| Dispatcher fallback dedup | `runAgent` skips fallback if `msg.RunID == ""` or `agent.RunCompletePublished(ctx)`; context.Canceled produces no error terminal | internal/backend/agent.go:87-122; internal/agent/run_marker.go:16-52 |
| RunID correlator | `WithRunID`/`RunIDFromContext` context helpers carry per-request correlator to the terminal event | internal/agent/runid.go:14-33 |
| Delivery semantics split | `Publish` lossy non-blocking for streaming deltas; `PublishMustDeliver` bounded-blocking (50 ms/subscriber default) for terminal events; drop counters exposed | internal/pubsub/broker.go:4-23, 34-46, 159-236 |
| Debounced message store contract | `Update` buffers or sync-flushes on terminal detection; `FlushAll` waits out in-flight timer writes so shutdown/session-switch see durable state | internal/message/message.go:215-302, 304-339 |
| Consumer-side correlation | `crush run` exits only on RunComplete whose RunID matches (else SessionID fallback); reconciles stdout against embedded Text backstop | internal/cmd/run.go:384-455 |
| Panic containment (HTTP) | `recoverHandler` middleware logs stack and returns 500; re-panics `http.ErrAbortHandler` | internal/server/recover.go:9-38 |
| Panic containment (goroutines) | MCP client init goroutine converts panic to error state; shell tool recovers; TUI event pump recovers via `log.RecoverPanic`; no recover wraps `sessionAgent.Run` itself (searched `recover()` repo-wide: 5 sites listed) | internal/agent/tools/mcp/init.go:601-625; internal/shell/shell.go:262; internal/shell/run.go:63; internal/app/app.go:707; internal/log/log.go:66-89 |
| Shutdown ordering | Shutdown cancels all agents first, then FlushAll with 5s bound, then parallel cleanup — agents finish writing state before DB close | internal/app/app.go:743-761 |
| Test: concurrent dispatch | Burst of 8 concurrent `Run` calls; concurrency-probe model asserts maxSeen==1 stream in flight | internal/agent/dispatch_race_test.go:101-155 |
| Test: cancel-on-entry end-to-end | Accepted-cancel race driven through real coordinator via gated wrapper; asserts persisted `FinishReasonCanceled` turn, not stubs | internal/backend/accepted_run_integration_test.go:65-131 |
| Test: both cancel branches fire | Active run + accepted follow-up: single Cancel fires active cancel func AND records pending mark | internal/agent/dispatch_cancel_test.go:68-88 |
| Test: idle cancel no-op | Idle Escape records no mark; next prompt runs normally | internal/agent/dispatch_cancel_test.go:326-359 |
| Test: seq poisoning | Prompt accepted after a cancel (higher seq) runs normally; earlier siblings still cancel-on-entry; mark clears after all covered handles resolve | internal/agent/dispatch_cancel_test.go:361-447 |
| Test: exactly-one terminal event | Cancelled queued RunID prompt publishes exactly one cancelled RunComplete; second event fails the test | internal/agent/run_complete_test.go:223-284 |
| Test: queued prompt own lifecycle | Gated model proves queued RunID prompt runs recursively and both turns publish distinct RunCompletes | internal/agent/queued_runid_test.go:72-179 |
| Test: dispatcher fallback trio | Pre-run error publishes terminal; coordinator-published suppresses duplicate; cancellation publishes no error terminal | internal/backend/agent_runcomplete_test.go:76-163 |

## Answers to Dimension Questions

**Q1: Can two goroutines publish conflicting terminal outcomes, and what prevents it?**

No — by construction, though it takes a guard chain rather than a single choke point. First, only one run can be active per session: the busy check and `activeRequests.Set` are performed atomically inside `dispatchMu` (`internal/agent/agent.go:586-650`), so a second concurrent dispatcher queues instead of streaming; the regression test proves max-one in-flight stream under an 8-goroutine burst (`internal/agent/dispatch_race_test.go:113-155`). Second, all terminal publications route through `publishRunComplete` (`internal/agent/agent.go:539-548`), and the four upstream producers are each prevented from double-firing: the queue-recursion defer via `skipRunComplete` (`internal/agent/agent.go:1282-1286`), the coordinator's unauthorized→retry chain via overwriting coalescing closure and single post-retry publish (`internal/agent/coordinator.go:281-336`), and the backend dispatcher's fallback via an atomic `runCompleteMarker` carried in context (`internal/agent/run_marker.go:39-52`; checked at `internal/backend/agent.go:112`). Enqueued copies have their `OnComplete` hook stripped so a drained turn cannot funnel its terminal event into a dead closure (`internal/agent/agent.go:363-379`; test `internal/agent/run_complete_test.go:21-59`). Residual honesty: these are conventions enforced by review and tests, not by a type system or a single owner object.

**Q2: Is a cancellation request represented separately from confirmed termination?**

Yes, in three distinct representations. (a) The request: `Cancel` invokes the stored cancel func but intentionally leaves the `activeRequests` entry installed so `IsBusy()` remains true until the goroutine fully completes, including error-path DB writes (`internal/agent/agent.go:1966-1973`). (b) Requests racing dispatch: a pending-cancel high-water mark over monotonically increasing accept sequences covers exactly the handles accepted-but-not-yet-active at cancel time (`internal/agent/agent.go:1988-2003`, semantics at 483-499); an idle Escape records nothing and poisons nothing (`internal/agent/dispatch_cancel_test.go:326-359`). (c) Confirmed termination: only when `Run` returns does the deferred publisher set `Cancelled` from the observed error or `ctx.Err()` (`internal/agent/agent.go:767-772`) and persist `FinishReasonCanceled` on the assistant message (`internal/agent/agent.go:526,1159-1160`). Even a cancel landing in the window between registration and assistant-message creation is persisted via the `currentAssistant == nil` recovery branch (`internal/agent/agent.go:1075-1089`) and proven against real machinery in `internal/backend/accepted_run_integration_test.go:80-131`.

**Q3: Can observers see a terminal state before its outcome or terminal event is available?**

Ordering is engineered but not globally guaranteed. Within a run, the terminal event is published only after `FlushAll` drains debounced message writes, and the event embeds the final MessageID and full text as a reconciliation backstop for subscribers that observe message events out of order across brokers (`internal/agent/agent.go:735-781`); the CLI consumes this contract explicitly (`internal/cmd/run.go:409-428`). Busy-state is released before the agent-finished notification so pollers never see stale running state (`internal/agent/agent.go:1209-1214`). Two residual windows remain by design: message updates and RunComplete travel on different brokers whose fan-in is not serialized (acknowledged in the code comments cited above), and `PublishMustDeliver` gives up per subscriber after a 50 ms bounded block, incrementing a drop counter and delegating recovery to the subscriber (`internal/pubsub/broker.go:190-236`). So: a well-behaved observer cannot see a terminal event ahead of its committed outcome, but can in principle miss the terminal event entirely and must recover.

**Q4: Which parts can Aren adopt without importing a framework-sized lifecycle model?**

Highly adoptable, dependency-free pieces (all small enough to lift nearly verbatim):
- Per-session lazy mutex registry with guarded creation (`internal/agent/agent.go:339-354`, ~15 lines).
- `activeCancel` pointer-identity CAS cleanup to stop deferred cleanup deleting a newer run's entry (`internal/agent/agent.go:160-167,652-657`; `internal/csync/maps.go:63-80`).
- Accept counter + monotonic sequence + cancel high-water mark (~70 lines total, `internal/agent/agent.go:264-337,1981-2003`) — the core answer to cancel-vs-dispatch races.
- RunID correlator and run-complete marker as unexported context keys (`internal/agent/runid.go`; `internal/agent/run_marker.go`) — zero-signature-change dedup between layers.
- The 7-field `RunComplete` payload with Error/Cancelled distinction (`internal/agent/notify/notify.go:71-78`).
- Split delivery semantics: lossy for deltas, bounded-blocking for terminal events (`internal/pubsub/broker.go`).

What would over-generalise Aren's contract if copied wholesale: the queue fold-vs-run split keyed on RunID presence exists to keep the `crush run` CLI from hanging and encodes crush's client topology (`internal/agent/agent.go:381-424`); the coordinator's OnComplete coalescing is shaped around provider 401-retry chains (`internal/agent/coordinator.go:266-336`); and crush's implicit, memory-only run states (maps instead of a persisted run row) assume a long-lived interactive process — Aren's Phase 1 contract, per the dimension, wants state/outcome/timing committed coherently, which suggests a persisted transition record rather than crush's ephemeral maps. Abstractions smaller than Aren's contract: `runCompleteMarker` (one bool) suffices only because crush has exactly two terminal publishers; a system with more producers needs the generalised version of the same idea.

## Architectural Decisions

1. **Implicit lifecycle state in concurrent maps, not an enum.** Running = `activeRequests` entry (`internal/agent/agent.go:2048-2051`); dispatched-not-active = `acceptedRuns` count (`internal/agent/agent.go:192-197`); pending cancel = `cancelMark` (`internal/agent/agent.go:198-206`). Nothing about run phase is persisted; durable history lives only as message `FinishReason`s (`internal/message/content.go`, used at `internal/agent/agent.go:1019,1160`).
2. **Lock the handoff, never the work.** `dispatchMu` is held only across the dispatch decision and queue handoff; all DB/LLM I/O happens unlocked ("The lock is held only during the brief handoff (no DB or LLM I/O under the lock)", `internal/agent/agent.go:187-190`).
3. **Exactly-once terminal event as an explicit product contract.** Comments repeatedly frame a missing/duplicate RunComplete as a user-visible hang for `crush run` (`internal/agent/agent.go:596-604,426-434`; `internal/backend/agent.go:69-80`), driving the single-emit-path design.
4. **Sequence-numbered cancellation coverage.** One Cancel covers all handles accepted before it and none accepted after, via a max-raised high-water mark (`internal/agent/agent.go:1988-2003`); idempotent repeated cancels (`internal/agent/accepted_run_test.go:100-114`).
5. **Detached, bounded contexts for terminal writes.** Persistence of cancelled turns, final-state cleanup, flushes, and dropped-prompt publishes all use `context.WithoutCancel` + timeouts so shutdown cannot strand them (`internal/agent/agent.go:509,754,1097,446`).
6. **Layered dedup over signature changes.** Rather than changing interfaces, dedup state travels in context values (RunID, marker) (`internal/agent/runid.go:14-22`; `internal/agent/run_marker.go:31-43`).

## Notable Patterns

- **Cancel-on-entry**: a pre-registered cancel converts an incoming run into a persisted cancelled turn without ever touching the model (`internal/agent/agent.go:591-618`), including its own terminal publish since the streaming defer is not yet installed.
- **ABA-safe deferred cleanup** with pointer-identity compare-and-delete (`internal/agent/agent.go:653-657`).
- **Coalescing callback injection** to collapse multi-attempt provider retries into one visible terminal event (`internal/agent/coordinator.go:285-288`).
- **Probe-model race tests**: a fake model tracks max concurrent streams via CAS loop, converting a locking bug into a numeric assertion (`internal/agent/dispatch_race_test.go:19-67,150-151`).
- **Observable-event assertions**: terminal-event tests read the published broker channel and fail on a second event, matching what external callers actually rely on (`internal/agent/run_complete_test.go:228-245`).
- **Gated production machinery in integration tests**: a wrapping coordinator parks `RunAccepted` so a cancel deterministically lands in the accepted-but-not-active window against the real implementation (`internal/backend/accepted_run_integration_test.go:27-63`).

## Tradeoffs

- **Convention over enforcement**: the exactly-once guarantee spans three packages connected by hooks, flags, and context markers; a new producer that bypasses `publishRunComplete` would silently break it. No type-level or runtime invariant prevents this.
- **Lossy-by-design streaming**: ordinary `Publish` drops under contention with only a counter and warning (`internal/pubsub/broker.go:159-188`) — right for token deltas, but it means intermediate state has no delivery guarantee.
- **Bounded-blocking terminals still lossy**: `PublishMustDeliver` drops after 50 ms per subscriber (`internal/pubsub/broker.go:42-46,222-235`); hang-free clients are achieved by embedding the outcome in the event and expecting subscriber-side recovery, not by guaranteed delivery.
- **Memory-only run state**: fast and simple, but a process restart loses all in-flight runs; recovery is deferred to the next prompt, which synthesizes tool results for orphaned tool calls so the conversation isn't bricked (`internal/agent/agent.go:1655-1687`).
- **Queueing beats rejecting**: busy sessions enqueue follow-ups rather than failing with `ErrSessionBusy` (that error is only returned by `Summarize`, `internal/agent/agent.go:1330-1332`); simpler for users, but it makes the fold-vs-sequence logic (and its RunID special cases) necessary complexity.
- **Polling shutdown wait**: `CancelAll` sleeps in 200 ms increments up to 5 s waiting for quiescence (`internal/agent/agent.go:2026-2034`) — best-effort, not a barrier.

## Failure Modes / Edge Cases

- **Unrecovered panic inside stream callbacks**: `sessionAgent.Run` installs no `recover`; repo-wide search found recovery only at HTTP middleware, MCP init goroutines, shell execution, and the TUI pump (`internal/server/recover.go:18`, `internal/agent/tools/mcp/init.go:613`, `internal/shell/shell.go:262`, `internal/shell/run.go:63`, `internal/app/app.go:707`). A panic raised by a fantasy stream callback or an `OnXxx` handler propagates through `Run`, skipping the terminal-publish defer's remaining work only up to the defer (defers do run; the process still dies afterward since nothing above catches it). Work errors, runtime failures, and panics are therefore NOT uniformly distinct across the goroutine boundary — panics escape.
- **Dead cancel key for summarize**: `Cancel` probes `activeRequests.Get(sessionID+"-summarize")` (`internal/agent/agent.go:1976`) but `Summarize` registers under plain `sessionID` (`internal/agent/agent.go:1355`); the branch can never match. Harmless today only because the plain-key entry is cancelled by the first branch.
- **Summarize dispatch not serialized**: `Summarize` checks busy then sets `activeRequests` without holding `dispatchMu` (`internal/agent/agent.go:1330-1357`), leaving a narrow window where a concurrent Cancel misses the cancel func; its queued-handoff (`internal/agent/agent.go:1464-1477`) also lacks the guarded handoff `Run` performs at `internal/agent/agent.go:1238-1314`. Secondary paths do not meet the same arbitration bar as the main path.
- **Untracked title goroutine**: `GenerateTitle` is launched detached with `WithoutCancel` (`internal/agent/agent.go:703-710`) and is not counted in any WaitGroup; its fallback rename (`internal/agent/agent.go:1736-1744`) could execute after shutdown's FlushAll/DB close (inferred risk from ordering at `internal/app/app.go:743-761`, not an observed failure).
- **Stale-mark hygiene depends on cooperation**: marks clear at three places (end of accepted batch `internal/agent/agent.go:330-334`, normal completion `internal/agent/agent.go:1265-1279`, drain `internal/agent/agent.go:403-424`); tests cover staleness (`internal/agent/dispatch_cancel_test.go:173-201`) but the invariant is maintained by hand.
- **Idle-cancel poisoning was a real bug class**: the sequence-mark design exists specifically because an earlier counted-pending-cancel scheme let later accepts consume earlier cancellations (`internal/agent/dispatch_cancel_test.go:361-370` documents the reviewer finding).

## Future Considerations

For Aren Phase 1 (lifecycle ownership before any provider/tool/persistence broadening):

1. **Make the terminal choke point structural.** Adopt `publishRunComplete`'s idea but as the only constructor of terminal transitions (e.g., a `Terminal` value type that must be minted once), removing reliance on flags (`skipRunComplete`) and hook stripping.
2. **Lift the accept-counter + sequence high-water mark** essentially verbatim for Aren's cancellation races (`internal/agent/agent.go:264-337,1988-2003`); add the property tests from `dispatch_cancel_test.go`.
3. **Persist the transition record.** Crush's ephemeral maps satisfy an interactive TUI but leave restart semantics undefined; Aren's requirement that state, outcome, timing, and notification commit coherently argues for writing the transition (with timestamps) in the same critical section that flips observable state.
4. **Decide the panic contract explicitly.** Either wrap task bodies with recover→error conversion (as MCP init does at `internal/agent/tools/mcp/init.go:613`) or declare panics fatal; don't inherit crush's accidental mix.
5. **Generalise the marker** (one atomic bool) into a producer-count-aware terminal ledger if Aren expects more than two terminal publishers.

## Questions / Gaps

- **Aren requirement docs unread**: `Aren/docs/phase-1-prd/02-lifecycle-requirements.md` and `01-product-definition-and-scope.md` were not accessible from this task's sandbox (external-directory denial). Alignment claims above are relative to the dimension text only.
- **Fantasy internals out of scope**: whether `fantasy.Agent.Stream` internally recovers panics or serializes callbacks could change the panic assessment; the vendored library is outside the source boundary, so this is unresolved (search boundary: entire `sources/crush` tree for `recover()`).
- **TUI consumption path not traced**: how the Bubble Tea UI reacts to `RunComplete`/busy-state transitions (e.g., `internal/ui/...`) was sampled only via `app.Subscribe` (`internal/app/app.go:706-736`); a UI-level study would need its own pass.
- **No evidence found** of any persisted, queryable run/status record: searched `internal/db/sql` schema surface indirectly through `internal/session/session.go:50-63` (no status field) and grep for status-like enums in the agent package (none beyond message `FinishReason` constants). If Aren requires durable run rows, crush offers no template for it.
- **`clearPendingCancel` appears test-only** (defined `internal/agent/agent.go:473-481`; production callers none found via grep excluding `_test.go`), suggesting either dead code or an unfinished intent to expose mark-clearing to embedders.

---

Generated by `01.01-lifecycle-transition-ownership-and-terminal-arbitration` against `crush`.
