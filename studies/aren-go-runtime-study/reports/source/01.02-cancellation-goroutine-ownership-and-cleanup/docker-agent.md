# Source Analysis: docker-agent

## 01.02 Cancellation, Goroutine Ownership, and Cleanup

### Source Info

| Field | Value |
|-------|-------|
| Name | docker-agent |
| Path | `studies/aren-go-runtime-study/sources/docker-agent` |
| Language / Stack | Go (module `github.com/docker/docker-agent`), Bubble Tea TUI, MCP SDK, OTel |
| Analyzed | 2026-08-26 |

## Summary

docker-agent implements a layered cancellation architecture with a deliberate, documented split between **request-scoped** and **connection-scoped** lifetimes. The root context comes from `signal.NotifyContext(SIGINT, SIGTERM)` (`main.go:16`); each CLI/TUI turn derives a per-message `context.WithCancel` (`pkg/cli/runner.go:88`, `pkg/tui/page/chat/chat.go:1301`), which flows into `RunStream`. The runtime itself detaches its long-lived base context at construction (`context.WithoutCancel`, `pkg/runtime/runtime.go:679`) so toolset connections survive individual requests.

Goroutine ownership follows three patterns: (1) **channel-close contracts** — every `RunStream` goroutine guarantees its events channel closes on all exit paths via a deferred `finalizeEventChannel` (`pkg/runtime/loop.go:335-337`), making `for range` the authoritative completion signal; (2) **supervisor-owned watchers** — one watcher goroutine per MCP/LSP connection that "must outlive ctx; the only way to stop it is Stop" (`pkg/tools/lifecycle/supervisor.go:289`); (3) **deliberately detached work** — background agents (`context.WithCancel(context.WithoutCancel(ctx))`, `pkg/tools/builtin/agent/agent.go:337`), background shell jobs (`pkg/tools/builtin/backgroundjobs/backgroundjobs.go:272`), and event-log pumps (`pkg/server/session_manager.go:239-246`), each with an explicit stop path (`Handler.StopAll`, `ToolSet.Stop`, `dropEventLog`).

Cleanup is bounded almost everywhere (process-group SIGTERM→SIGKILL escalation with a 3s grace in `pkg/tools/builtin/shell/shell.go:221-232`; a 30s budget around `StopToolSets` in `cmd/root/run.go:1381-1387`), shutdown is idempotent under concurrent callers (CAS state machines, `sync.Once`, published-stop-request handshakes), and cleanup ordering is specified so teardown cannot overwrite or mask the primary outcome. The codebase ties many of these behaviors to regression tests named after issue numbers (#2872, #3069/#3070, #3200, #3584, #4001, #4004).

## Rating

**9 / 10**

Rationale:
- Every claim about lifetime is enforced by code and usually by a test: channel-close terminal contracts (`pkg/runtime/loop.go:187-200`), supervisor stop/reap semantics with 25+ dedicated tests (`pkg/tools/lifecycle/supervisor_test.go:190-820`), close-vs-send race elimination under RWMutex with race-detector tests (`pkg/runtime/elicitation.go:123-142`).
- Cancellation propagation reaches the leaves that matter: OS processes via process groups (`pkg/tools/builtin/shell/cmd_unix.go:14-25`), HTTP via `CancelCauseFunc` closing the TCP body (`pkg/runtime/streaming.go:60-63,350-351`), blocked user prompts via `select { resume | ctx.Done() }` (`pkg/runtime/toolexec/dispatcher.go:835-845`).
- Detachment decisions are explicit and justified in comments (e.g., why MCP must detach and how opt-in cancellation is restored, `pkg/tools/mcp/cancellable_parent.go:14-59`).
- Deductions: no process-wide goroutine-leak harness (no goleak anywhere; verified by search), a handful of fire-and-forget goroutines whose reaping depends on callers following documented contracts rather than type-enforced ownership, and `restoreAndClose` can block for the lifetime of a parked sender's ctx (`pkg/runtime/elicitation.go:128-136`, accepted trade-off).

## Evidence Collected

Every entry cites paths relative to `studies/aren-go-runtime-study/sources/docker-agent`.

| Area | Evidence | File:Line |
|------|----------|-----------|
| Root signal handling | `signal.NotifyContext(context.Background(), SIGINT, SIGTERM)`; cancel called in both exit branches | main.go:16-26 |
| Per-turn CLI cancel | CLI run wraps ctx in `WithCancel`; Ctrl+C abort path calls `cancel()`; cancelled confirmations skip Resume | pkg/cli/runner.go:88-89,174-177,191-194 |
| Runtime base ctx detachment | LocalRuntime stores `func() context.Context { return context.WithoutCancel(ctx) }` used for out-of-stream ops (e.g., toolsChanged emit) | pkg/runtime/runtime.go:679,1545 |
| RunStream goroutine + channel contract | RunStream spawns loop goroutine returning `<-chan Event`; deferred `finalizeEventChannel` guarantees close on every exit incl. early returns | pkg/runtime/loop.go:235-258,327-337 |
| Terminal-signal ordering | StreamStopped emitted non-blockingly before session-end hooks; channel close is the authoritative "all cleanup done" signal (#3070/#3074 rationale) | pkg/runtime/loop.go:187-211 |
| Non-cancellable cleanup hooks | `session_end` hooks run under `context.WithoutCancel(ctx)` so they fire after Ctrl+C | pkg/runtime/loop.go:213-215 |
| turn_end on all exits | Deferred dispatch runs even on panic/cancel; endReason re-checked against `ctx.Err()` | pkg/runtime/loop.go:750-764 |
| Observer forwarder goroutine | `observe` spawns forwarder closing `out` when inner drains; observers run synchronously (back-pressure documented) | pkg/runtime/observer.go:66-81,14-24 |
| Elicitation close/send race fix | bridge `send` holds RLock across send w/ recover; `restoreAndClose` takes write lock before `close(current)` (#3069) | pkg/runtime/elicitation.go:104-121,137-142 |
| Elicitation waiter terminal-state machine | resolve vs cancel atomic CAS; ctx.Done branch drains an already-won response instead of discarding it (#3584) | pkg/runtime/elicitation.go:144-171,547-570 |
| Detached bridge delivery | elicitation request mirrored to stream channel from a detached goroutine bounded by ctx | pkg/runtime/elicitation.go:532-545 |
| Headless fast-decline | non-interactive sessions decline elicitations instead of parking forever (#3200) | pkg/runtime/elicitation.go:473-487; pkg/runtime/loop.go:283-296 |
| Batch tool-call cancellation | `batchCtx = WithCancelCause(ctx)`; first canceled/stopped call cancels siblings via `stopOnce`; cause distinguishes user-cancel vs hook-stop messages | pkg/runtime/toolexec/dispatcher.go:240-276,582-598 |
| Confirmation wait cancellation | `askUser` selects on resume channel vs `ctx.Done()`; serialized by `confirmationMu` so only one prompt waits | pkg/runtime/toolexec/dispatcher.go:809-845 |
| Cancel-as-Ok translation | tool errors where `errors.Is(err, context.Canceled)` become user-cancellation results (span Ok) sent back to model next turn | pkg/runtime/toolexec/dispatcher.go:1111-1131 |
| Unbounded parallel fan-out | `concurrent.MapSlice` spawns one goroutine per element, no bound (documented as intended for small fan-outs) | pkg/concurrent/mapslice.go:5-22 |
| Shell timeout + group kill | manual `timeoutCtx`; on expiry SIGTERM to process group, 3s grace then SIGKILL, wait for pipe goroutines | pkg/tools/builtin/shell/shell.go:188-238 |
| Grandchild pipe drain bound | `cmd.WaitDelay = 500ms` so backgrounded grandchildren inheriting pipes can't block Wait forever | pkg/tools/builtin/shell/shell.go:173-197 |
| Process groups (unix) | `Setpgid: true`; `kill(-pid, SIGTERM)` hits the whole tree | pkg/tools/builtin/shell/cmd_unix.go:10-25 |
| Failed-setup reap | `reapSpawnedChild` kills and waits (2s grace → SIGKILL) when post-start setup fails | pkg/tools/builtin/shell/shell.go:240-261 |
| Background jobs detached monitor | `go h.monitorJob(context.WithoutCancel(ctx), ...)` — jobs intentionally outlive request; Stop kills all running job groups | pkg/tools/builtin/backgroundjobs/backgroundjobs.go:235,272,619-627 |
| Background jobs wait tool | `wait_background_job` selects on done/timeout/ctx — cancellable wait with explicit cancel message | pkg/tools/builtin/backgroundjobs/backgroundjobs.go:397-423 |
| Background agents ownership | taskCtx = WithCancel(WithoutCancel(ctx)); explicit stop via HandleStop CAS; `StopAll` cancels all + `wg.Wait()` | pkg/tools/builtin/agent/agent.go:333-356,496-515 |
| Runtime.Close wiring | `LocalRuntime.Close()` = `bgAgents.StopAll()`; explicitly does NOT close embedder-supplied session store | pkg/runtime/runtime.go:202-208,1480-1483 |
| Close-doesn't-close-shared-store test | regression test pinning that Close never closes an external SQLite store (#2872) | pkg/runtime/close_test.go:27-55 |
| Supervisor watcher lifetime | watcher runs `context.WithoutCancel(ctx)`; only `Stop` ends it; `watchDone` closed on watcher exit | pkg/tools/lifecycle/supervisor.go:287-294,512-522 |
| Supervisor idempotent Stop | `stopping` flag; second caller waits on watchDone bounded by ctx; pending connect reaped in background | pkg/tools/lifecycle/supervisor.go:404-452 |
| In-flight connect adoption | startup-timeout leaves Connect running in `pendingConnect`; next Start adopts it, Stop reaps it; never two overlapping Connects | pkg/tools/lifecycle/supervisor.go:302-402 |
| StartupTimeout via timer not ctx | timer chosen because MCP connector detaches ctx — deadline would be stripped; wedged Connect intentionally never superseded to avoid orphaned handshakes/subprocess leaks | pkg/tools/lifecycle/supervisor.go:111-117,312-330 |
| Restart policy + backoff | RestartOnFailure/Never/Always; max attempts default 5; permanent errors (auth) stop retries | pkg/tools/lifecycle/supervisor.go:572-590,595-667 |
| Bounded start attempts | `TryStartWithTimeout` races TryStart against timeout; abandoned attempt keeps single-flight lock so later calls join rather than race | pkg/tools/startable.go:376-420 |
| Pending-stop handshake | `StopIfStarted` publishes request before waiting on lock; lock holder's `unlock` settles it under `WithoutCancel` (#4001 wedged-start case) | pkg/tools/startable.go:548-586,601-646 |
| MCP connection detach + opt-in parent | `ctx = WithoutCancel(ctx)` then stash original as value; OAuth select observes parent's Done to avoid stuck stream lock/409 busy | pkg/tools/mcp/mcp.go:530-546; pkg/tools/mcp/cancellable_parent.go:14-59 |
| Per-call MCP timeout | callCtx child of caller ctx (never replaces detached session ctx); timeout classified as `ErrCallTimeout` only when outer ctx still live | pkg/tools/mcp/mcp.go:768-784,808-830 |
| LSP process lifetime | `processCtx = WithCancel(WithoutCancel(callerCtx))`; stderr drain goroutine bound to it; handshake abort closes stdin on caller cancel | pkg/tools/builtin/lsp/lsp_lifecycle.go:78-147 |
| LSP Wait/Close single-flight | `sync.Once` + pre-allocated waitDone channel; Close idempotent via atomic.Bool; signal-exit after close treated clean | pkg/tools/builtin/lsp/lsp_lifecycle.go:210-255 |
| Scheduler fully detached loop | Start builds its own `context.Background()` cancel; Stop cancels and waits on WaitGroup | pkg/tools/builtin/scheduler/scheduler.go:107-125 |
| Session manager pump lifetime | RegisterEventSource pump on `context.WithCancel(context.Background())`; DeleteSession tombstones via `dropEventLog` so no pump/log resurrects | pkg/server/session_manager.go:227-247,300-313 |
| RunSession stream goroutine defers | defer order: Unlock before close(streamChan) so a drained consumer cannot spuriously see ErrSessionBusy | pkg/server/session_manager.go:1032-1040 |
| DeleteSession reaping | cancels runtime cancel, moves entry to deletedSessions; background poller TryLocks streaming mutex up to 5 min; WaitStopped offers bounded client wait | pkg/server/session_manager.go:853-888,891-919 |
| Attached-runtime cancel discipline | DELETE cancels attach-lifetime cancel, never the in-flight stream cancel (done != nil guard) | pkg/server/session_manager.go:977-983,1246-1250 |
| Recall-run detached persistence | recall stream persists final session under `context.WithoutCancel(ctx)` so Ctrl+C doesn't lose history write | pkg/server/session_manager.go:1270-1277 |
| SSE idle-timeout shutdown | reader goroutine + idleTimer; stall fires `cancelStream(errStreamIdle)` whose CancelCauseFunc closes the TCP connection to unblock the reader | pkg/runtime/streaming.go:19-39,73-103,349-351 |
| Graceful model-cancel classification | `handleStreamError` treats `context.Canceled` as graceful fatal stop, not retryable error | pkg/runtime/loop_steps.go:168-173 |
| Toolset startup fan-out | EmitStartupInfo starts toolsets concurrently (WaitGroup + per-toolset result channels), dependents in second wave; per-toolset timeouts skip wedged starts | pkg/runtime/runtime.go:1747-1786,1803-1838 |
| Bounded tool listing | `listToolsWithTimeout` orphans the Tools() goroutine safely via buffered channel | pkg/runtime/runtime.go:1925-1958 |
| Hook execution bounds | each hook runs under `context.WithTimeout(hook.GetTimeout())`; timeout/cancel normalized to fail-closed error | pkg/hooks/executor.go:336-364 |
| Pause gate | pauseCh closed on resume wakes all waiters; waitIfPaused selects on ctx so cancel breaks pause | pkg/runtime/pause.go:14-48 |
| TUI per-message cancel | chatPage stores `msgCancel`; new message/retry cancels previous stream context first | pkg/tui/page/chat/chat.go:1293-1301,1349-1356 |
| App drain-after-cancel | after ctx cancel the App keeps draining but forwards only StreamStopped (under WithoutCancel) so supervisor marks idle | pkg/app/app.go:754-779 |
| TUI supervisor cleanup outside lock | CloseSession cancels runner ctx, removes maps, runs cleanup in goroutine outside mutex | pkg/tui/service/supervisor/supervisor.go:425-471 |
| Backend cleanup chain | localBackend cleanup `sync.Once`: stopToolSets (30s) then rt.Close; backend.Close owns shared store; worktree cleanup uses fresh `WithoutCancel` ctx | cmd/root/backend.go:110-133,173-182; cmd/root/run.go:573-585,1381-1387 |
| Team-level stop aggregation | Team.StopToolSets iterates agents, continuing across failures (pinned by test) | pkg/team/team.go:158-165; pkg/team/team_test.go:39-53 |
| Cancellation-aware fs walks | CollectFiles/DirectoryTree return `context.Canceled`/`DeadlineExceeded` promptly (dedicated tests) | pkg/fsx/collect_cancellation_test.go:15-87 |
| Server graceful shutdown pattern | MCP serve: Serve in errCh goroutine; on ctx.Done Shutdown under `WithTimeout(WithoutCancel(ctx), 5s)` | pkg/mcp/server.go:117-135 |
| Supervisor concurrency tests | StopIdempotent, StopConcurrent, StopWaitsForWatcher, adopt/reap late connects, permanent-error no-retry | pkg/tools/lifecycle/supervisor_test.go:190-280,401-508,561-637,781-820 |

## Answers to Dimension Questions

**Who is responsible for waiting for every goroutine started on behalf of a run?**
It depends on the class of goroutine, and the responsibility is assigned explicitly:
- *Stream goroutines* are waited on implicitly by consumers draining the returned channel; the runtime guarantees termination (close) via `finalizeEventChannel` on every path (`pkg/runtime/loop.go:230-337`). The API-server wrapper adds ordered defers so unlock precedes channel close (`pkg/server/session_manager.go:1032-1040`).
- *Background agent tasks* are owned by `agenttool.Handler`: `runtime.Close()` → `StopAll()` cancels all tasks and blocks on `wg.Wait()` (`pkg/runtime/runtime.go:1480-1483`, `pkg/tools/builtin/agent/agent.go:505-515`).
- *Long-lived toolsets* are waited for by the embedder: `cleanup()` runs `StopToolSets` under a 30s budget then `rt.Close()` (`cmd/root/backend.go:122-133`); the supervisor's Stop blocks until the watcher exits (`pkg/tools/lifecycle/supervisor.go:431-440`).
- *Deliberately unwaited* goroutines exist and are documented: abandoned `TryStartWithTimeout` attempts keep holding the single-flight lock by design (`pkg/tools/startable.go:389-395`), `listToolsWithTimeout` orphans exit later via buffered sends (`pkg/runtime/runtime.go:1929-1931`), and event pumps live until `DeleteSession` tombstones them (`pkg/server/session_manager.go:305-313`). There is no central registry or leak detector enforcing these contracts — correctness rests on documentation plus targeted tests.

**What does the API claim when cancellation is requested but work has not stopped?**
The terminal publication is two-tiered: a best-effort `StreamStoppedEvent` carrying reason `"canceled"` (emitted non-blocking, droppable) followed by the guaranteed events-channel close; consumers are told to treat only the close as authoritative (`pkg/runtime/loop.go:187-211`). For sessions deleted while still streaming, `DeleteSession` returns immediately and offers `?wait=true` / `WaitStopped` to poll the streaming mutex up to a caller-supplied timeout; otherwise a background reaper polls for at most 5 minutes (`pkg/server/session_manager.go:853-919`). Mid-flight tool calls synthesize "canceled" error responses into the conversation so the model sees consistent tool_result pairs next turn (`pkg/runtime/toolexec/dispatcher.go:582-598,1111-1121`).

**Can cleanup block forever or overwrite the primary failure?**
Blocking is bounded by construction: shell teardown escalates SIGTERM→SIGKILL after 3s and always joins the Wait goroutine (`pkg/tools/builtin/shell/shell.go:221-232`); supervisor Stop honors caller ctx (`pkg/tools/lifecycle/supervisor.go:442-452`); hook/session-end cleanup runs under `WithoutCancel` but each hook is individually timeout-bounded (`pkg/runtime/loop.go:213-215`, `pkg/hooks/executor.go:336-364`). One documented unbounded-ish case: `restoreAndClose` waits on the write lock while a parked bridge send holds the read lock; that send is ctx-bounded and dispatched from a detached goroutine, so the wait ends when the sender's ctx does — accepted over a send-on-closed panic (`pkg/runtime/elicitation.go:128-142`). Overwriting the primary failure is guarded: cleanup errors are logged, not propagated over run errors (`cmd/root/backend.go:124-133`), `Supervisor.Stop` returns closeErr only when ctx is still live (`pkg/tools/lifecycle/supervisor.go:431-439`), and the stream end-reason defaults to normal but is overridden by a `ctx.Err()` check so cancellation isn't misreported as success (`pkg/runtime/loop.go:750-757`).

**How is shutdown made idempotent under concurrent callers?**
Multiple complementary mechanisms: `Supervisor.Stop` uses a `stopping` flag so concurrent callers just wait on `watchDone` (`pkg/tools/lifecycle/supervisor.go:406-412`; tested at `pkg/tools/lifecycle/supervisor_test.go:421-433,802-820`); `StartableToolSet.StopIfStarted` publishes a stop request that exactly one holder claims and settles at unlock (`pkg/tools/startable.go:562-586,621-646`); `lspSession.Close` is idempotent via `closed atomic.Bool` and `sync.Once` wait (`pkg/tools/builtin/lsp/lsp_lifecycle.go:220-255`); background-task stop transitions are CAS-guarded (`pkg/tools/builtin/agent/agent.go:496-500`); backend cleanup is `sync.Once` (`cmd/root/backend.go:122-133`); and `signal.NotifyContext`'s cancel is safe to call in both success and failure branches (`main.go:16-27`).

## Architectural Decisions

1. **Two-lifetime context model.** Request contexts drive turns; connection contexts are detached at the boundary (`pkg/tools/mcp/mcp.go:544-546`, `pkg/tools/builtin/lsp/lsp_lifecycle.go:84`, `pkg/runtime/runtime.go:679`). Where an operation inside a detached region needs user cancellation, the parent is stashed as a context value for opt-in observation rather than weakening the detachment (`pkg/tools/mcp/cancellable_parent.go:14-59`).
2. **Channel close as the sole authoritative completion signal.** Events may be dropped under backpressure; close never is (`pkg/runtime/loop.go:187-211`). This makes `for range RunStream(...)` a correct join point everywhere.
3. **Centralize reconnect/restart in one Supervisor** shared by stdio MCP, remote MCP, and LSP, parameterized by Connector + Policy, with adoption/reaping of timed-out connects (`pkg/tools/lifecycle/supervisor.go:153-207,302-402`).
4. **Bounded abandonment over unbounded waiting for hostile dependencies.** Wedged toolset starts/listings are raced against timers and abandoned while holding locks so subsequent operations join instead of racing (`pkg/tools/startable.go:383-420`, `pkg/runtime/runtime.go:1932-1958`).
5. **Cancellation is conversational, not silent.** User cancels produce explicit tool_result errors ("canceled by the user") fed back to the model, keeping the message invariant intact (`pkg/runtime/toolexec/dispatcher.go:240-244,1111-1131`).
6. **Cleanup hooks survive cancellation.** `turn_end`/`session_end` dispatch under `context.WithoutCancel` so audit/teardown hooks run after Ctrl+C (`pkg/runtime/loop.go:213-215,758-762`).
7. **Explicit ownership registry for detached sub-work**: live-session registry with unregister-then-drain so queued compaction requests can't be stranded (`pkg/runtime/live_sessions.go:62-83,121-142`).

## Notable Patterns

- **Published-stop-request handshake** (`pkg/tools/startable.go:563-586,625-646`): publish intent before contending for a lock; whoever releases the lock settles it — converts "can't acquire lock during shutdown" into guaranteed eventual stop.
- **Terminal-state CAS for waiter resolution** (`pkg/runtime/elicitation.go:144-171`): resolves response-vs-cancel races without losing answered responses.
- **Adopt-or-reap pending work** (`pkg/tools/lifecycle/supervisor.go:302-402`): a single `adopted bool` guarantees a late Connect result is either handed to exactly one Start or closed by Stop.
- **Grace-period kill ladder** for processes: SIGTERM group → timed wait → SIGKILL → unconditional Wait (`pkg/tools/builtin/shell/shell.go:221-232`, `reapSpawnedChild` at 240-261; same shape in `pkg/tools/builtin/backgroundjobs/backgroundjobs.go:459-476`).
- **WaitDelay for grandchild pipes** (`pkg/tools/builtin/shell/shell.go:173-197`): pragmatic bound on Go's exec pipe-copy goroutines.
- **Ordered defers encode shutdown ordering** (`pkg/server/session_manager.go:1032-1040`, `pkg/runtime/loop.go:339-351`): LIFO registration documented so compaction drains and plan subscriptions release before the channel closes.
- **Issue-number-pinned regression tests** for lifecycle bugs (#2872 `pkg/runtime/close_test.go:27-55`; #3069/#3070 in `pkg/runtime/elicitation.go:123-136` and `pkg/runtime/loop.go:196-199`; #3584 suite `pkg/runtime/elicitation_concurrency_test.go:27-113`; #4001 in `pkg/tools/startable.go:556-561`).

## Tradeoffs

- **Detached-by-default connections**: maximizes session reuse but requires every interactive operation inside the detached region to remember to observe the stashed cancellable parent; a new operation that forgets will hang until the connection dies (`pkg/tools/mcp/cancellable_parent.go:31-49` documents the exact incident that motivated the pattern).
- **Abandoned-but-running starts**: `TryStartWithTimeout` timeouts leave goroutines and subprocesses starting in the background; memory-safe (buffered channels, lock handoff) but resource usage is unbounded until the underlying attempt returns (`pkg/tools/startable.go:389-395`).
- **Unbounded fan-out parallelism** in tool-call batches and hook dispatch — fine for typical batch sizes, a footgun for pathological models issuing hundreds of parallel calls (`pkg/concurrent/mapslice.go:8-11`).
- **restoreAndClose blocking on parked senders** trades potential teardown delay for panic-freedom (`pkg/runtime/elicitation.go:128-136`, explicitly marked "do not 'fix' by dropping the lock").
- **Polling-based session-stop detection** (`TryLock` tickers at 100ms/50ms with hard deadlines) instead of condition variables — simpler and deadlock-free, at the cost of latency and CPU noise (`pkg/server/session_manager.go:861-881,891-919`).
- **No global leak enforcement**: relies on review and targeted tests; a future detached goroutine without a stop path would not be caught mechanically.

## Failure Modes / Edge Cases

- **SIGTERM-ignoring children**: covered by the 3s escalation ladder (`pkg/tools/builtin/shell/shell.go:226-232`).
- **Grandchild processes holding pipes**: bounded by `WaitDelay` (500ms), deliberately letting grandchildren keep running while the tool returns (`pkg/tools/builtin/shell/shell.go:176-186`).
- **Half-open TCP to model gateway**: 5-minute SSE idle timeout with `cancelStream(errStreamIdle)` closing the transport (`pkg/runtime/streaming.go:19-39,349-351`).
- **Concurrent RunStreams stealing the single elicitation slot**: reliable sink delivery made primary, bridge demoted to best-effort mirror on a detached, ctx-bounded goroutine (#3584) (`pkg/runtime/elicitation.go:67-76,527-545`).
- **DeleteSession racing lazy event-log creation**: tombstone map + lock serialization ensures no log/pump is created for a deleted session (`pkg/server/session_manager.go:220-232,263-272`).
- **Partial startup then crash**: toolset start failures latch partial composites started so healthy subsets remain usable and failed ones retry (`pkg/tools/startable.go:484-507`); team-level stop continues across agent failures (`pkg/team/team.go:158-165`).
- **Pause held while user quits**: `waitIfPaused` selects on ctx, so Ctrl+C exits a paused loop rather than deadlocking (`pkg/runtime/loop.go:451-461`, `pkg/runtime/pause.go:35-48`).
- **OAuth flow vs user abort mid-flow**: without the cancellable-parent pattern this deadlocked the streaming lock and 409'd the next message; now the flow's select watches the parent Done (`pkg/tools/mcp/cancellable_parent.go:31-41`).

## Future Considerations

- Introduce `goleak` (or equivalent) in package test mains for the concurrency-heavy packages (`pkg/runtime`, `pkg/tools/lifecycle`, `pkg/server`) to mechanically catch regressions in the many documented-but-unenforced goroutine contracts.
- Bound `concurrent.MapSlice` parallelism (semaphore) or document per-call-site maxima, since tool batches originate from model output.
- Give background shell jobs the same explicit owner/wait story as background agents: their `Stop` exists (`pkg/tools/builtin/backgroundjobs/backgroundjobs.go:619-627`) but nothing in the runtime's `Close` path enumerates them independently of generic `StopToolSets` ordering; a crash between `rt.Close()` and `StopToolSets` (ordered `stopToolSets` first today, `cmd/root/backend.go:124-131`) is currently safe only because of that order.
- Reconsider polling reapers (`deletedSessions` 100ms ticker, 5-minute cap, `pkg/server/session_manager.go:861-881`) in favor of a done-channel on the stream goroutine for precise, immediate reclamation.
- The `restoreAndClose` wait could take a bounded timeout with a last-resort recover-based close if teardown latency under a wedged consumer ever matters in production.

## Questions / Gaps

- No evidence found of a repo-wide goroutine-leak test harness: searched `goleak` (no matches) and `NumGoroutine` in tests (no matches). Leak prevention is per-feature testing only.
- The dimension citation references Aren PRD docs (`Aren/docs/phase-1-prd/02-lifecycle-requirements.md`, phased roadmap). Per source-isolation rules this study inspected only `sources/docker-agent`, so mapping to those requirement documents was out of scope and was not performed.
- `pkg/board/`, `pkg/chatserver/`, and ACP adapters contain additional goroutines (`pkg/board/controller.go:127`, `pkg/chatserver/server.go:193`) that were sampled only lightly; their shutdown paths follow the same `Shutdown(bounded WithoutCancel ctx)` pattern observed elsewhere but were not exhaustively traced.
- Whether any Windows-specific gap exists in process-group kill semantics was not verified beyond confirming platform-split files exist (`pkg/tools/builtin/shell/cmd_windows.go`, `cmd_unix.go`); only the Unix path was read in full.

---

Generated by `01.02-cancellation-goroutine-ownership-and-cleanup` against `docker-agent`.
