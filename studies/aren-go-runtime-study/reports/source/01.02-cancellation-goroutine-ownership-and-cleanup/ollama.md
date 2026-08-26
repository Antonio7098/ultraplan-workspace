# Source Analysis: ollama

## 01.02 Cancellation, Goroutine Ownership, and Cleanup

### Source Info

| Field | Value |
|-------|-------|
| Name | ollama |
| Path | `studies/aren-go-runtime-study/sources/ollama` |
| Language / Stack | Go (gin HTTP server, `os/exec`-managed llama-server/MLX subprocesses, errgroup, golang.org/x/sync/semaphore) |
| Analyzed | 2026-08-26 |

Citation convention: paths are relative to the source root (`studies/aren-go-runtime-study/sources/ollama/`). All line numbers verified against the working tree.

## Summary

Ollama is a long-lived model server whose runtime centers on three ownership layers: (1) a process-level lifecycle in `server.Serve` that wires a signal handler to an HTTP server close, a scheduler context cancel, and runner unloading (`server/routes.go:2002-2032`); (2) a two-goroutine scheduler (`processPending`, `processCompleted`) that owns `runnerRef` entries and arbitrates load/unload through buffered channels (`server/sched.go:213-223`, `server/sched.go:60-81`); and (3) per-runner subprocess wrappers (`llm.llamaServerRunner`, `x/mlxrunner.Client`) that each own an `exec.Cmd`, a reaper goroutine, and a `done` channel (`llm/llama_server.go:1005-1030`, `x/mlxrunner/client.go:398-402`).

Context propagation is generally disciplined: request contexts flow from gin handlers into the scheduler (`server/routes.go:223` → `server/sched.go:169-211`) and into subprocess health polling and completion streaming (`llm/llama_server.go:1344-1346`, `llm/llama_server.go:1627-1671`). Cancellation during model load is handled explicitly (`server/sched.go:237-240`, `llm/llama_server.go:1299-1302`), and unload is made idempotent by PID matching plus map-presence checks under two locks with a documented lock-ordering rule (`server/sched.go:429-451`, `server/sched.go:1753-1754`).

The main weaknesses are at the edges of streaming publication and terminal shutdown: streaming handlers pair an unbuffered progress/generation channel with a producer goroutine whose only cancellation lever is context — but a client disconnect stops the consumer before the producer, leaving it parked on a channel send inside the completion callback (`server/routes.go:649-792`, `llm/llama_server.go:1708`); SIGTERM handling uses `srvr.Close()` rather than graceful drain and never joins the scheduler loops (`server/routes.go:2021-2030`, `server/routes.go:2060-2067`); and the llama-server kill path waits unboundedly for process reaping after SIGKILL with no escalation timeout, unlike the MLX client which bounds it at 5s (`llm/llama_server.go:2573-2589` vs `x/mlxrunner/client.go:146-150`).

## Rating

**6 / 10.**

Rationale: strong context plumbing into blocking operations, genuinely idempotent unload bookkeeping, refcounted detached downloads with cancel-on-last-release, one exemplary dedicated-thread worker package, and at least one real goroutine-leak regression test. Held back by: producer-side channel sends with no consumer-abandonment escape in all streaming handlers, abrupt `Close()` on SIGTERM with no wait-group join of background loops, an unbounded final reap wait on the primary inference path, dropped cleanup errors, and no repo-wide leak detection (no goleak; searched `goleak|VerifyTestMain` across all Go files — no matches).

## Evidence Collected

| Area | Evidence | File:Line |
|------|----------|-----------|
| Root context tree | `Serve` creates root ctx + child `schedCtx`; both canceled by the signal goroutine | server/routes.go:2002-2004, 2021-2030 |
| Scheduler loop ownership | `Run` spawns exactly two loops; both exit on `<-ctx.Done()`; doc comment says "will shutdown when ctx is done" | server/sched.go:213-223, 228-232, 369-373 |
| Request→scheduler context flow | `LlmRequest.ctx` stored at enqueue; skipped if already canceled | server/sched.go:30, 169-195, 237-240 |
| Runner allocation select | Handler selects on `runnerCh`/`errCh`, both cap-1 so sends never block producers | server/routes.go:223-231, server/sched.go:189-190 |
| Load-time cancellation | `WaitUntilRunning` returns wrapped ctx.Err on client close; failure publishes to `errCh` then `expiredCh` | llm/llama_server.go:1299-1302, server/sched.go:731-738 |
| Streaming cancellation check | Completion scanner loop checks `ctx.Done()` once per token before parsing | llm/llama_server.go:1667-1671 |
| Concurrency bounding per runner | `semaphore.Weighted(numParallel)` acquired/released around Completion | llm/llama_server.go:1536-1542, 948-963 |
| Subprocess reaper | Detached goroutine runs `cmd.Wait()`, stores `doneErr`, closes `done` | llm/llama_server.go:1018-1027 |
| Close path (llama-server) | `stopProcess`: skip if already exited → SIGKILL → unbounded `<-s.done` wait | llm/llama_server.go:2569-2589 |
| Close path (MLX) | SIGINT → 5s bounded wait → SIGKILL escalation | x/mlxrunner/client.go:137-154 |
| Unload idempotency | Duplicate expired events deduped by map presence; orphaned runners shut down via PID mismatch | server/sched.go:429-451 |
| Double-load safeguard | Same-key reload force-unloads stale runner ("Shouldn't happen, but safeguard against leaking a runner") | server/sched.go:718-725 |
| Lock ordering rule | "refMu must not be acquired while holding loadedMu: the expiration path locks them in the opposite order" | server/sched.go:1753-1754 |
| Keepalive expiry | `time.AfterFunc` per idle runner; timer stopped/reset under `refMu` | server/sched.go:393-410, 479-482 |
| VRAM-recovery gate | Bounded poll (default 5s) on `context.WithTimeout(context.Background(), …)`; unload waits on result before publishing `unloadedCh` | server/sched.go:1441-1499, 453-466 |
| Evict-all OOM retry | `evictAllAndWait` counts exact unload events, aborts on ctx.Done | server/sched.go:1586-1630, 316-326 |
| Detached download | `go download.Run(context.Background(), …)` survives client disconnect; refcount release cancels via `CancelFunc` | server/download.go:512-513, 436-444, 218-220 |
| Detached upload | Same pattern: Background ctx + acquire/release + CancelFunc | server/upload.go:400-401, 309-315, 131 |
| Parallel part I/O | `errgroup.WithContext` + `SetLimit` for download parts; stall watchdog as second errgroup goroutine | server/download.go:277-309, 333-384 |
| Cache hydration workers | Bounded worker pool over jobs channel, ctx-checked send loop, WaitGroup joined | server/model_show_cache.go:339-384 |
| Detached cache refresh | `context.WithoutCancel` + fire-and-forget refresh goroutine guarded by an in-flight flag | server/model_show_cache.go:250-274 |
| Startup hydration | One-shot goroutines for list/show caches started off the serve path | server/model_list_cache.go:65-83, server/model_show_cache.go:126-131 |
| GPU validation fan-out | WaitGroup + shared 30s timeout ctx for per-device probe goroutines | discover/runner.go:132-180 |
| Dedicated OS thread | MLX thread worker: LockOSThread never unlocked, stop-job pattern, CAS-idempotent Stop, panic capture with stack | x/internal/mlxthread/thread.go:46-66, 95-132, 134-155, 181-191 |
| Process-group kill (agent bash tool) | `exec.CommandContext` + `Setpgid` + group SIGKILL + `WaitDelay` + explicit ErrWaitDelay branch | agent/tools/bash_unix.go:24-49, agent/tools/bash.go:82-97, 130-133 |
| Graceful proxy shutdown (contrast) | Claude gateway: `sync.Once` shutdown flag, `server.Shutdown(ctx)`, CloseIdleConnections | internal/proxy/claude_desktop.go:210-236 |
| Readiness-check goroutine | Single in-flight readiness probe; callers select on ctx vs done | internal/proxy/claude_desktop.go:442-462, 464-511 |
| CLI cancel wiring | SIGINT channel → ctx cancel for interactive chat | cmd/cmd.go:1770-1779 |
| App restart flow | Server run under child ctx; restart = cancel, wait done, re-spawn | app/cmd/app/app.go:242-266 |
| Idempotent UI progress stop | `stopOnce` + closed done channel; spinner same | progress/progress.go:45-65, progress/spinner.go:91-99 |
| Goroutine-leak regression test | `TestStopReapsGoroutines` asserts NumGoroutine returns to baseline after 50 start/stop cycles | progress/progress_test.go:49-73 |
| Scheduler shutdown tests | Expire→unload cycle including loop termination; already-canceled request; unload-all | server/sched_test.go:632-714, 1191-1207, 1158-1179 |
| Load stall tests | Timeout when load stalls; timeout extended on log activity | llm/llama_server_test.go:1075, 1106 |
| Stream-close ordering test | Done callback deferred until response body closed (`TestLlamaServerCompletionDoneCallbackAfterStreamClosed`) | llm/llama_server_test.go:2714-2754 |
| Queue bounds | All four scheduler channels sized by `OLLAMA_MAX_QUEUE` (default 512); ErrMaxQueue on overflow | envconfig/config.go:279, server/sched.go:88-96, 204-208 |

## Answers to Dimension Questions

**Who is responsible for waiting for every goroutine started on behalf of a run?**
No single owner. The scheduler's two loops are owned by `schedCtx` but nothing joins them: after `srvr.Serve` returns `ErrServerClosed`, `Serve` waits only for the root ctx and returns nil without confirming `processPending`/`processCompleted` have exited (`server/routes.go:2060-2067`). Per-request watcher goroutines (`useLoadedRunner`'s ctx-done publisher at `server/sched.go:487-491`, the load goroutine at `server/sched.go:731-752`) terminate via their channel sends completing or ctx firing, not via any WaitGroup. The streaming producer goroutines in `GenerateHandler`/`PullHandler`/`PushHandler`/`CreateHandler` have no owner at all if their consumer abandons the channel (see Failure Modes). The only place with true join semantics is the bounded hydration pool (`server/model_show_cache.go:349-384`) and GPU discovery (`discover/runner.go:180`).

**What does the API claim when cancellation is requested but work has not stopped?**
Nothing explicit — there is no "canceling…" state or acknowledgment protocol. For generation, the stream simply ends mid-flight; the runner stays loaded and its keepalive timer decides actual unload later (`server/sched.go:393-410`). For loads, a client disconnect logs "client connection closed before llama-server finished loading, aborting load" and returns a wrapped ctx.Err phrased as a *timeout* ("timed out waiting for llama-server to start: %w") — conflating client cancel with stall timeout (`llm/llama_server.go:1300-1302`). Truthful failure publication does exist for crashes: the reaper converts the last stderr message into `doneErr`, surfaced as "llama-server process has terminated: <exit>: <msg>" (`llm/llama_server.go:1021-1025`, 1303-1325).

**Can cleanup block forever or overwrite the primary failure?**
Block forever: yes, on the primary path. `stopProcess` issues SIGKILL then waits unboundedly on `<-s.done` (`llm/llama_server.go:2579-2585`); a llama-server stuck in an uninterruptible GPU-driver teardown keeps `unloadAllRunners` blocked while it holds `loadedMu` (`server/sched.go:1695-1711`), wedging the whole scheduler and, because the signal goroutine itself called it, making further Ctrl-C inert (buffered size-1 channel, reader exited — `server/routes.go:2022-2030`). The MLX client shows the intended pattern: escalate to Kill after 5s (`x/mlxrunner/client.go:146-150`). Overwrite primary failure: no — cleanup errors are *dropped* rather than propagated: `runnerRef.unload` ignores `llama.Close()`'s error (`server/sched.go:1372-1374`) and `unloadAllRunners` ignores all Close results (`server/sched.go:1705-1710`), so they can't mask the request error pushed to `errCh` first (`server/sched.go:671`), but the losses are silent.

**How is shutdown made idempotent under concurrent callers?**
Layered guards: duplicate expired events are detected by checking whether the runner is still in `s.loaded` and comparing PIDs, with distinct handling for the "orphaned runner" case (`server/sched.go:429-451`); a same-key reload force-unloads any stale runner first ("safeguard against leaking a runner", `server/sched.go:718-725`); `stopProcess` short-circuits when `ProcessState != nil` (`llm/llama_server.go:2575-2577`); MLX `Close` nil-guards the cmd (`x/mlxrunner/client.go:142`); `ClaudeDesktop.Close` and the UI progress/spinner use `sync.Once` (`internal/proxy/claude_desktop.go:224`, `progress/progress.go:47`, `progress/spinner.go:98`); the mlxthread worker uses `stopping.CompareAndSwap` and treats post-stop enqueues as `ErrStopped` (`x/internal/mlxthread/thread.go:98-105`, 163-165). The signal-handler goroutine in `Serve` is one-shot by construction but has no second-signal escape hatch.

## Architectural Decisions

1. **Single-threaded arbitration through channels, not shared queues.** All load/unload decisions funnel through `processPending`/`processCompleted` exchanging four maxQueue-sized channels; refcount mutations happen only in these loops or under `refMu` (`server/sched.go:60-81`, 367-470). This trades throughput for a serialized, auditable state machine.
2. **Cancellation is cooperative end-to-end, with hard kill as backstop.** Client ctx → semaphore acquire → HTTP call → per-token checks (`llm/llama_server.go:1536-1542`, 1627, 1667-1671); the subprocess itself only ever receives SIGKILL on teardown (`llm/llama_server.go:2579`). There is no attempt to gracefully interrupt llama-server inference (e.g., no SIGINT escalation like MLX's Close).
3. **Detached-by-design background work with reference-counted cancellation.** Pull/push intentionally survive requester disconnect (`//nolint:contextcheck` + `context.Background()`, `server/download.go:512-513`, `server/upload.go:400-401`); lifetime is tied to watchers via `acquire/release → CancelFunc` (`server/download.go:436-444`).
4. **VRAM convergence gating.** Unload ordering is: snapshot GPUs → `runner.unload()` → delete from map → block until recovery-or-timeout → publish `unloadedCh` (`server/sched.go:452-467`). This encodes the empirical fact that driver free-memory reporting lags process exit (`server/sched.go:1433-1440`).
5. **Lock hierarchy as documentation.** The loadedMu-before-refMu ordering constraint is stated at the snapshot helper that must respect it (`server/sched.go:1753-1754`) — an implicit invariant enforced by review, not by tooling.
6. **Newer code converges on stricter ownership.** The MLX thread package formalizes single-owner execution with queue-only cancellation semantics ("Once the worker accepts a job, the job runs until fn returns", `x/internal/mlxthread/thread.go:70-72`), and the Claude gateway proxy uses graceful `Shutdown(ctx)` + once-shutdown — both noticeably more rigorous than the core serve path.

## Notable Patterns

- **Cap-1 result channels** (`successCh`, `errCh`) so late failures after the caller gives up never block producers (`server/sched.go:189-190`).
- **Self-requeueing retry goroutine**: expired event with positive refCount sleeps 10ms and re-enqueues itself (`server/sched.go:417-426`) — simple bounded backpressure via the channel buffer.
- **Stall detection via atomic activity timestamp**: log-parsing writer updates `loadActivity`; `WaitUntilRunning` extends its deadline only on observed progress (`llm/llama_server.go:1160-1188`, 1281-1337).
- **Panic transport across goroutine boundary**: mlxthread wraps recovered panics with the worker stack so re-panic prints the original fault site (`x/internal/mlxthread/thread.go:33-44`, 181-191).
- **Process-group hygiene for agent tools**: `Setpgid` + negative-PID SIGKILL kills grandchildren spawned by shell scripts, with `cmd.WaitDelay` catching pipe-holders (`agent/tools/bash_unix.go:35-49`, `agent/tools/bash.go:94`).
- **Leak regression testing with `runtime.NumGoroutine`**: pins the fix for spinners that "previously ran until the process exited" (`progress/progress_test.go:49-73`).

## Tradeoffs

- **Abrupt SIGTERM (`srvr.Close()`) over graceful drain** (`server/routes.go:2026`): fast teardown of in-flight streams, but clients get connection resets instead of a final error frame, and there is no drain window for half-written NDJSON.
- **Serialized unload behind `loadedMu`** (`server/sched.go:1695-1711`): simple consistency; one slow-killing subprocess freezes eviction decisions for every other model.
- **Unbuffered progress channels** (`ch := make(chan any)`): natural backpressure against slow clients, but couples producer lifetime to consumer politeness (see Failure Modes).
- **Background-context detachment for transfers**: pull resilience across client churn, at the cost of disk/network work that no request can see or cancel except by starting-and-dropping a watcher.
- **Dropping cleanup errors**: keeps the primary failure truthful but hides resource-leak signals (e.g., repeated Kill failures would be invisible).

## Failure Modes / Edge Cases

- **Producer-parked goroutine leak on client disconnect (streaming)**: `GenerateHandler` spawns a producer that sends each token to an unbuffered `ch` from inside the `Completion` callback (`server/routes.go:649-656`, 704/715/727/733/739/741/746); the only consumer is `streamResponse(c, ch)` (`server/routes.go:792`). On client-gone the gin stream loop stops reading (gin v1.10.0 `Context.Stream` selects on the writer's CloseNotify — inferred from the pinned dependency in `go.mod:8`; gin source not present in this checkout), the handler returns, and the producer stays parked on `ch <- res` inside `fn` (`llm/llama_server.go:1708`). Because `Completion` cannot return, its `defer s.sem.Release(1)` never runs (`llm/llama_server.go:1542`); with `OLLAMA_NUM_PARALLEL` defaulting to 1 (`envconfig/config.go:275`), one abandoned connection pins the model's sole completion slot until the runner is evicted/reloaded. The ctx.Done check exists only at the top of the scanner loop (`llm/llama_server.go:1668-1671`), which is unreachable while blocked in `fn`. The same producer/consumer shape exists in `PullHandler` (`server/routes.go:1133-1160`), `PushHandler` (`1185-1215`), and `CreateHandler` (`server/create.go:115-120`). Non-streaming paths drain fully via `for rr := range ch` (`server/routes.go:755-789`) and `waitForStream` (`2070-2095`), so they are safe.
- **Unbounded final reap wait**: SIGKILL + `<-s.done` with no deadline (`llm/llama_server.go:2579-2585`) — a D-state subprocess hangs shutdown permanently and, via the signal goroutine, disables further SIGINT handling.
- **HTTP-handler blocking on scheduler channels**: `expireRunner` (unload requests) sends to `expiredCh` while holding `refMu` (`server/sched.go:1719-1729`, invoked from `server/routes.go:405`); if the 512-slot buffer fills while `processCompleted` is blocked waiting on VRAM recovery (`server/sched.go:463`), API calls stall behind it. Same send-under-lock shape in `processPending` (`server/sched.go:338-348`).
- **Partial-startup failure**: load crash after spawn is recovered via the evict-all-and-retry path with a one-shot latch (`oomRetryAttempted`, `server/sched.go:38-40`, 641-668), tested in `server/sched_test.go:1745-1901`. A second persistent crash fails fast to `errCh` — no infinite retry.
- **Repeated cancellation**: double expired events are absorbed (`server/sched.go:429-441`); a canceled pending request is silently skipped without publishing to its channels, relying on cap-1 buffers and the handler's select (`server/sched.go:237-240`, tested at `server/sched_test.go:1191-1207`).
- **Token-repeat runaway**: generation aborts with `ctx.Err()` after 30 identical tokens (`llm/llama_server.go:1690-1701`) — returns ctx.Err even when the context is live, mislabeling the cause.
- **Stream EOF mid-generation**: unexpected EOF triggers immediate `s.Close()` of the runner and surfaces the last stderr message (`llm/llama_server.go:1754-1767`) — aggressive but prevents serving from a corrupted process.

## Future Considerations

- Make streaming producers abandonment-safe: select on a done channel alongside every `ch <-` send, or give the producer ownership of a `context.WithCancel` it checks in `fn`; alternatively size `ch` to decouple callback return from consumer liveness. This also fixes the semaphore-slot pinning.
- Add an escalation timeout to `stopProcess` (mirror MLX's 5s SIGINT→SIGKILL ladder) and/or move runner closing out from under `loadedMu`.
- Replace `srvr.Close()` with `Shutdown(ctx)` plus a bounded drain, and add a `WaitGroup` (or `errgroup`) joining `processPending`/`processCompleted` before `Serve` returns nil.
- Introduce `goleak` (or a NumGoroutine harness generalized from `progress/progress_test.go:53-73`) for the scheduler and route packages; today leak coverage exists only for the UI progress widget.
- Second-signal forced exit in `Serve` (e.g., second SIGINT ⇒ `os.Exit(1)`) to escape wedged cleanup.
- Encode the loadedMu/refMu ordering in a lint rule or wrapper rather than a comment (`server/sched.go:1753-1754`).

## Questions / Gaps

- **Does gin v1.10.0's `c.Stream` actually abandon reads on client-gone?** The module cache was unavailable in this environment, so the producer-park scenario rests on repo-side evidence (unbuffered channel, synchronous callback send, sem released only at Completion return) plus library knowledge of gin's CloseNotify-based Stream loop. A targeted experiment (disconnect mid-stream, observe NumGoroutine and `/api/ps` slot occupancy) would confirm severity. Searched: `go.sum` pins gin v1.10.0; no vendored copy exists in-tree.
- **No evidence found for graceful llama-server interruption**: searched for `Signal(`, `SIGINT`, `Interrupt` across `llm/` — only the MLX client sends a pre-kill signal. Whether llama.cpp requires SIGKILL-only teardown (e.g., CUDA destructor hangs) is undocumented in-repo.
- **What happens to in-flight `finishedReqCh` events during shutdown?** After `schedDone()`, `processCompleted` exits on ctx; queued finished/expired events are dropped and `refCount` bookkeeping dies with the process. No evidence this matters beyond process exit, but it means shutdown-time unload does not follow the normal expire path (`server/sched.go:371-373` vs `unloadAllRunners` bypassing it entirely at `1695-1711`).
- **Aren PRD alignment**: the dimension cites `../../../../Aren/docs/phase-1-prd/02-lifecycle-requirements.md`, which lies outside the isolated source directory and was therefore not read (source-isolation rule). Contract-mapping to Aren requirements could not be performed from this study alone.

---

Generated by `01.02-cancellation-goroutine-ownership-and-cleanup` against `ollama`.
