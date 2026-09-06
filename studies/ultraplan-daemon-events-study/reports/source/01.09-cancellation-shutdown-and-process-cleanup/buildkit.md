# Source Analysis: buildkit

## Dimension 01.09: Cancellation, Shutdown, and Process Cleanup

### Source Info

| Field | Value |
|-------|-------|
| Name | buildkit |
| Path | `studies/ultraplan-daemon-events-study/sources/buildkit` |
| Language / Stack | Go (gRPC daemon `buildkitd`, containerd/runc executors, LLB solver) |
| Analyzed | 2026-09-03 |

## Summary

BuildKit implements cancellation as volatile gRPC `context.Context` cancellation propagated through the solver pipe/scheduler into executor kill paths — there is no durable cancel intent object or cancel RPC. The scheduler (`solver/scheduler.go`, `solver/edge.go`, `solver/internal/pipe/pipe.go`) arbitrates cancel-vs-complete by marking pipe status `Canceled` and refusing to cache or commit canceled errors, so a canceled job does not poison shared vertices. Executors translate `ctx.Done` into `SIGKILL` with bounded kill RPC timeouts (`7s` runc, `10s` containerd, `10s` git) but then wait unboundedly for process exit. Daemon shutdown is two-stage drain: first `SIGTERM/SIGINT` cancels `appcontext.Context()` while `server.GracefulStop()` drains gRPC, second signal trips `appcontext.Shutdown()`, third calls `Fatal`; history queue defers listener close until finalizers drain. Post-cancel cleanup systematically uses `context.WithoutCancel` for releases, deletes, leases, and history recording, so finalization is not starved by the parent cancel.

## Rating

**7 / 10** — Clear context-propagation model with scheduler-level cancel arbitration, per-executor SIGKILL escalation with timeouts, two-stage daemon shutdown, and systematic `WithoutCancel` cleanup contexts, all covered by scheduler cancel tests and gateway exec-cancel integration tests. Withheld from 9–10 because cancellation is not durable (pure transport cancel, no persisted intent), daemon shutdown is drain-only with no cancel/interrupt of in-flight builds, and kill paths can block indefinitely after the kill timeout waiting for `p.ended`.

## Evidence Collected

Every entry MUST include a file path with line numbers. Format: `path/to/file.ts:NN`.

| Area | Evidence | File:Line |
|------|----------|-----------|
| Cancel entry / transport | `Solve` is a unary gRPC call taking caller `ctx`; no `Cancel` RPC exists — cancel is `ctx` cancellation. `Status` is a separate streaming RPC on same `Ref`. | `studies/ultraplan-daemon-events-study/sources/buildkit/control/control.go:409` |
| Session disconnect wiring | `Session` hijacks stream, derives `WithCancelCause(stream.Context())`, cancels only on hijack `closeCh`; session cancel does not itself cancel solver jobs. | `studies/ultraplan-daemon-events-study/sources/buildkit/control/control.go:627-642` |
| Scheduler cancel injection | `scheduler.build` wraps caller ctx with `WithCancelCause`, goroutine calls `p.Receiver.Cancel()` on `<-ctx.Done()`. | `studies/ultraplan-daemon-events-study/sources/buildkit/solver/scheduler.go:229-235` |
| Pipe cancel primitive | `receiver.Cancel()` sets `req.Canceled=true` idempotently; `Pipe.Sender.Request().Canceled` drives `cancel()` of function pipes. | `studies/ultraplan-daemon-events-study/sources/buildkit/solver/internal/pipe/pipe.go:182-187` |
| Cancel-vs-complete arbitration | `finishIncoming`: if request canceled and edge has no error, synthesize `context.Canceled`. | `studies/ultraplan-daemon-events-study/sources/buildkit/solver/edge.go:178-186` |
| Canceled errors not committed | `processCacheMapReq`, `processExecReq`, `processDepReq`, `processDepSlowCacheReq` all guard with `!upt.Status().Canceled` before assigning `e.err`. | `studies/ultraplan-daemon-events-study/sources/buildkit/solver/edge.go:569-576`, `studies/ultraplan-daemon-events-study/sources/buildkit/solver/edge.go:621-632`, `studies/ultraplan-daemon-events-study/sources/buildkit/solver/edge.go:650-657`, `studies/ultraplan-daemon-events-study/sources/buildkit/solver/edge.go:700-706` |
| Canceled exec not cached | `sharedOp.Exec/CacheMap/CalcSlowCache`: on error with `ctx.Done` + `IsCanceled`, set `complete=false`, `releaseError(err)`, wrap `context.Cause(ctx)` — result not stored in `execDone/cacheDone`. | `studies/ultraplan-daemon-events-study/sources/buildkit/solver/jobs.go:1248-1260`, `studies/ultraplan-daemon-events-study/sources/buildkit/solver/jobs.go:1170-1180`, `studies/ultraplan-daemon-events-study/sources/buildkit/solver/jobs.go:1105-1113` |
| Cancel classifier | `IsCanceled` checks `errors.Is(err, context.Canceled)`, gRPC `codes.Canceled`, plus `context.Cause(ctx)` with `EOF`/`"context canceled"` substring fallback for containerd-typed errors. | `studies/ultraplan-daemon-events-study/sources/buildkit/solver/errdefs/context.go:12-27` |
| Cancel-vs-exit race | `exitError`: on non-zero/unknown exit, `select { case <-ctx.Done(): wrap with context cause; default: return stack error }` — cancel wins if both ready. | `studies/ultraplan-daemon-events-study/sources/buildkit/executor/runcexecutor/executor.go:429-435` |
| runc kill path | `runcProcessHandle`: detached `runcCtx=Background`; goroutine on callers `ctx.Done` calls `doKillProc`; `Kill` sends `SIGKILL` via `runc kill` (run) or pidfile+`SIGKILL` (exec) with 10s fail-safe timeout. | `studies/ultraplan-daemon-events-study/sources/buildkit/executor/runcexecutor/executor.go:653-687`, `studies/ultraplan-daemon-events-study/sources/buildkit/executor/runcexecutor/executor.go:581-597` |
| runc kill escalation deadline | `doKillProc`: kill attempt uses `context.WithoutCancel(ctx)` + 7s timeout; on success breaks; on failure retries every 50ms until `p.ended`; after success waits unboundedly (`<-p.ended` after 50ms warn). | `studies/ultraplan-daemon-events-study/sources/buildkit/executor/runcexecutor/executor.go:780-821` |
| runc process-group / death signal | runc `Setpgid:true`; `PdeathSignal=SIGKILL` with explicit `// this can still leak the process` comment. SIGKILL signals are rerouted to `killer.Kill`, never directly to runc monitor. | `studies/ultraplan-daemon-events-study/sources/buildkit/executor/runcexecutor/executor.go:125-132`, `studies/ultraplan-daemon-events-study/sources/buildkit/executor/runcexecutor/executor_linux.go:24-27`, `studies/ultraplan-daemon-events-study/sources/buildkit/executor/runcexecutor/executor.go:762-768` |
| containerd kill path | On `ctxDone`: fresh `Background` + 10s timeout kill ctx, `p.Kill(SIGKILL)`, `io.Cancel()`; on `killCtxDone` timeout returns `"failed to kill process on cancel"`; exit path wraps with `context.Cause(ctx)` if canceled. | `studies/ultraplan-daemon-events-study/sources/buildkit/executor/containerdexecutor/executor.go:429-483` |
| containerd cleanup without cancel | Container/task delete use `context.WithoutCancel(ctx)` + `WithProcessKill` so teardown survives parent cancel. | `studies/ultraplan-daemon-events-study/sources/buildkit/executor/containerdexecutor/executor.go:223-247` |
| Host-process group kill | Git fetch on Linux: `Setpgid:true, Pdeathsig:SIGTERM`; on `ctx.Done` kill negative pgid with `SIGTERM`, escalate to `SIGKILL` after 10s. | `studies/ultraplan-daemon-events-study/sources/buildkit/source/git/source_linux.go:37-63` |
| Daemon signal policy | `appcontext`: 1st SIG cancels `appContext`, 2nd cancels `shutdownContext`, 3rd `Fatal`. Two static contexts shared process-wide. | `studies/ultraplan-daemon-events-study/sources/buildkit/util/appcontext/appcontext.go:33-65` |
| Daemon shutdown drain | `buildkitd` main: `ctx=appcontext.Context()`; on `errCh` or `ctx.Done`, log `stopping server`, `server.GracefulStop()` (drain, no `Stop()`/interrupt path); telemetry closers run under `appcontext.Shutdown()` + 5s timeout. | `studies/ultraplan-daemon-events-study/sources/buildkit/cmd/buildkitd/main.go:267-268`, `studies/ultraplan-daemon-events-study/sources/buildkit/cmd/buildkitd/main.go:452-483` |
| History drain on graceful stop | History queue watches `GracefulStop`; closes pubsub only when `finalizers==0 && active==0`, else defers close to finalizer completion. | `studies/ultraplan-daemon-events-study/sources/buildkit/solver/llbsolver/history/buildhistory.go:139-147`, `studies/ultraplan-daemon-events-study/sources/buildkit/solver/llbsolver/history/buildhistory.go:500-513` |
| Controller close deadlines | `Controller.Close` joins HistoryDB/Worker/CacheStore/Solver closes; only trace forwarder has a deadline (5s `WithTimeoutCause`). | `studies/ultraplan-daemon-events-study/sources/buildkit/control/control.go:162-186` |
| Cleanup contexts | `WithoutCancel` used for mount-stub cleaner, cache/snapshot releases, lease deletes, history `rec()`, exporter `done()`, gateway container `Release`. Examples across cache, solver, frontend, exporter. | `studies/ultraplan-daemon-events-study/sources/buildkit/executor/runcexecutor/executor.go:280`, `studies/ultraplan-daemon-events-study/sources/buildkit/cache/manager.go:165-180`, `studies/ultraplan-daemon-events-study/sources/buildkit/solver/llbsolver/solver.go:290-292`, `studies/ultraplan-daemon-events-study/sources/buildkit/solver/llbsolver/solver.go:385-387` |
| Scheduler cancel tests | `TestSingleCancelCache`, `TestSingleCancelExec`, `TestSingleCancelParallel` assert canceled builds return `context.Canceled` and parallel sibling continues. | `studies/ultraplan-daemon-events-study/sources/buildkit/solver/scheduler_test.go:472-513`, `studies/ultraplan-daemon-events-study/sources/buildkit/solver/scheduler_test.go:514-555`, `studies/ultraplan-daemon-events-study/sources/buildkit/solver/scheduler_test.go:557-603` |
| Exec cancel integration tests | Gateway tests: ctx-cancel during `sh` exec yields `context.Canceled`; `Release` terminates `sleep 10` pids within 10s; SIGKILL delivery test exists. | `studies/ultraplan-daemon-events-study/sources/buildkit/client/gateway_container_exec_test.go:73-98`, `studies/ultraplan-daemon-events-study/sources/buildkit/client/gateway_container_exec_test.go:106-176`, `studies/ultraplan-daemon-events-study/sources/buildkit/client/gateway_container_exec_test.go:270-277` |
| Client cancel plumbing | `client/solve.go`: separate `solveCtx`/`statusContext`; solve cancel on gateway failure after 5s fallback; status stream ignores `context.Canceled`/inactivity-timeout errors; session closed after solve. | `studies/ultraplan-daemon-events-study/sources/buildkit/client/solve.go:258-286`, `studies/ultraplan-daemon-events-study/sources/buildkit/client/solve.go:347-385` |

## Answers to Dimension Questions

**1. Is client disconnect different from explicit cancellation?**
Yes, structurally. There is no explicit `CancelBuild` API; both cases are transport `ctx` cancellations but on different RPCs with different effects. Canceling the unary `Solve` ctx (`studies/ultraplan-daemon-events-study/sources/buildkit/control/control.go:409`) propagates via `scheduler.build` → `Receiver.Cancel()` (`studies/ultraplan-daemon-events-study/sources/buildkit/solver/scheduler.go:229-235`) and kills executors. Dropping only the `Status` stream (`studies/ultraplan-daemon-events-study/sources/buildkit/control/control.go:594-625`) just ends that `solver.Status` subscription; the build continues. Dropping the `Session` hijacked conn (`studies/ultraplan-daemon-events-study/sources/buildkit/control/control.go:627-642`) cancels only the session handler context. Client `solve.go` reinforces the split with independent `solveCtx` vs `statusContext` plus a 5s inactivity killer for the status stream (`studies/ultraplan-daemon-events-study/sources/buildkit/client/solve.go:258-286`).

**2. Is cancellation durable if no worker is currently attached?**
No. No evidence found of a persisted cancel intent, cancel journal, tombstone, or ref-indexed cancel state. Cancellation lives entirely in the in-memory pipe graph (`Request().Canceled` in `studies/ultraplan-daemon-events-study/sources/buildkit/solver/internal/pipe/pipe.go:60-73`, `studies/ultraplan-daemon-events-study/sources/buildkit/solver/edge.go:735-750`). `Solver.Get(id)` only looks up live `jl.jobs` with a 6s wait (`studies/ultraplan-daemon-events-study/sources/buildkit/solver/jobs.go:717-744`); history `Finalize/Delete` (`studies/ultraplan-daemon-events-study/sources/buildkit/control/control.go:348-368`) operate on records, not on cancel signals. A cancel issued with no attached `Solve` ctx is a no-op; re-attaching to the same `Ref` starts fresh scheduling.

**3. Can cleanup hang indefinitely?**
Yes, in one bounded-kill/unbounded-wait spot. `doKillProc` bounds the kill RPC itself (7s `WithoutCancel` ctx, 50ms retry loop in `studies/ultraplan-daemon-events-study/sources/buildkit/executor/runcexecutor/executor.go:780-807`) but after a successful SIGKILL it does `select { case <-p.ended: ... case <-time.After(50ms): warn }; <-p.ended` (`studies/ultraplan-daemon-events-study/sources/buildkit/executor/runcexecutor/executor.go:810-821`) with no further deadline — a wedged runc monitor blocks `Run` return forever. Counter-evidence: containerd path fails fast with `"failed to kill process on cancel"` on 10s kill timeout (`studies/ultraplan-daemon-events-study/sources/buildkit/executor/containerdexecutor/executor.go:476-482`); git path escalates TERM→KILL at 10s (`studies/ultraplan-daemon-events-study/sources/buildkit/source/git/source_linux.go:48-56`); most resource finalizers use `WithoutCancel` so they are not starved by the cancel itself.

**4. Can child/grandchild processes escape termination?**
Mostly no for container work, possibly yes for daemon-side host processes. Container kills go through `runc kill SIGKILL` to the container init (`studies/ultraplan-daemon-events-study/sources/buildkit/executor/runcexecutor/executor.go:595-597`) or containerd `Task.Kill(SIGKILL)` (`studies/ultraplan-daemon-events-study/sources/buildkit/executor/containerdexecutor/executor.go:440`) — kernel delivers SIGKILL unblockably to that namespace scope, and cgroup membership (`CgroupsPath` recorded in `studies/ultraplan-daemon-events-study/sources/buildkit/executor/runcexecutor/executor.go:349-357`) is used for OOM accounting, not for kill fallback. No cgroup freezer/kill or pid-namespace iteration code was found. Daemon-spawned host binaries that ignore process-group discipline could escape: only the git fetcher was found to set `Setpgid+Pdeathsig` and negative-pgid kill (`studies/ultraplan-daemon-events-study/sources/buildkit/source/git/source_linux.go:37-63`); the runc shim itself notes `PdeathSignal=SIGKILL // this can still leak the process` (`studies/ultraplan-daemon-events-study/sources/buildkit/executor/runcexecutor/executor_linux.go:26`). No global process-reaper or cgroup-sweep on shutdown was found.

**5. How are cancellation and completion races arbitrated?**
At three layers, cancel is checked but completion wins if already committed. (a) Scheduler: `respondToIncoming`/`finishIncoming` map canceled-but-errorless requests to `context.Canceled` while completed edges complete remaining waiters (`studies/ultraplan-daemon-events-study/sources/buildkit/solver/edge.go:178-186`, `studies/ultraplan-daemon-events-study/sources/buildkit/solver/edge.go:729-804`). (b) Shared-op cache: canceled errors are never stored (`complete=false` + `releaseError` in `studies/ultraplan-daemon-events-study/sources/buildkit/solver/jobs.go:1248-1260`), so a loser does not poison the winner; non-canceled errors are committed. (c) Executor exit: `exitError`/`statusCh` select prefers the already-observed exit unless `ctx.Done` is also ready, in which case the error is wrapped with `context.Cause(ctx)` (`studies/ultraplan-daemon-events-study/sources/buildkit/executor/runcexecutor/executor.go:429-435`, `studies/ultraplan-daemon-events-study/sources/buildkit/executor/containerdexecutor/executor.go:471-475`). Prompt-scenario answer: cancel arriving exactly as work completes yields either success (if edge/result committed first) or `context.Canceled`/`codes.Canceled` (if cancel won the select); a child ignoring the first signal still dies because the only signal ever sent on cancel is unblockable `SIGKILL` (`studies/ultraplan-daemon-events-study/sources/buildkit/executor/runcexecutor/executor.go:581-597`), with the caveat that the daemon then waits unboundedly for the reaped exit as in Q3.

## Architectural Decisions

- **Transport-cancel as the only cancel primitive.** No cancel endpoint/state machine; every layer (scheduler pipe → op exec → runc/containerd kill) keys off `ctx.Done`/`Receiver.Cancel`. Keeps the API small at the cost of durability (Q2).
- **Detached kill contexts.** Both executors mint a fresh `Background + timeout` kill ctx via `WithoutCancel` (`studies/ultraplan-daemon-events-study/sources/buildkit/executor/runcexecutor/executor.go:785`, `studies/ultraplan-daemon-events-study/sources/buildkit/executor/containerdexecutor/executor.go:437-438`) so the kill RPC itself cannot be canceled by the cancel it serves.
- **SIGKILL-only cancel semantics.** No graceful SIGTERM-then-wait for container processes; cancel means immediate SIGKILL. User `Signal` messages (SIGHUP/SIGINT) travel a separate `handleSignals`/`Kill(eventCtx, sig)` path (`studies/ultraplan-daemon-events-study/sources/buildkit/executor/runcexecutor/executor.go:749-776`).
- **Drain-biased daemon shutdown.** `GracefulStop` + history finalizer drain (`studies/ultraplan-daemon-events-study/sources/buildkit/cmd/buildkitd/main.go:465`, `studies/ultraplan-daemon-events-study/sources/buildkit/solver/llbsolver/history/buildhistory.go:139-147`) rather than cancel/interrupt in-flight builds; escalation is operator-driven (2nd/3rd signal).
- **`WithoutCancel` finalization discipline.** Releases, lease deletes, history records, exporter completions, and network-namespace teardown all run detached, distinguishing cleanup from cancellation pervasively.

## Notable Patterns

- **Pipe `Canceled` bit + `finishIncoming` synthesis:** uniform vocabulary (`pipe.go:42`, `pipe.go:182-187`, `edge.go:178-186`) reused by scheduler, cache, and exec layers.
- **Don't-cache-the-cancel:** `complete=false` + `releaseError` idiom in all three `sharedOp` paths prevents canceled flight-control winners from poisoning shared graph state.
- **Kill-handle split (`monitorProcess` vs `killer`):** signals for terminal UX go to the runc monitor; SIGKILL always goes to the in-container init via `runc kill`/pidfile.
- **Three-strike signal escalation:** app-level analog of TERM→KILL for the whole daemon in `appcontext.go:51-63`.
- **Status-stream liveness isolation:** client inactivity timer kills only the status subscription, never the solve (`client/solve.go:264-281`).

## Tradeoffs

- **Simplicity (no durable cancel) vs operability:** no cancel store to reconcile, but orchestrators cannot "cancel by ref" after disconnect; must keep the `Solve` ctx open or kill the whole daemon.
- **SIGKILL immediacy vs graceful checkpointing:** guarantees un-ignorable termination and simple race logic, but denies checkpoint/flush hooks; slow-exit warning (`executor.go:815`) is observability without remedy.
- **Drain vs interrupt on shutdown:** in-flight builds survive `SIGTERM` (good for CI correctness) but a wedged build can delay `GracefulStop` indefinitely — no build-level deadline ties shutdown to a bound.
- **Substring-based cancel detection (`IsCanceled`):** pragmatic interop with untyped containerd/gRPC-stream errors at the cost of fragility (relies on `"EOF"`/`"context canceled"` text).

## Failure Modes / Edge Cases

- **Unbounded `p.ended` wait after SIGKILL** (`studies/ultraplan-daemon-events-study/sources/buildkit/executor/runcexecutor/executor.go:818-821`): runc-monitor hang → executor `Run` never returns → edge never completes → daemon shutdown stalls in `GracefulStop`.
- **Runc shim leak on daemon death:** acknowledged in `executor_linux.go:26` — `PdeathSignal` helps but `can still leak`.
- **Kill-path `context.Cause` shadowing:** `doKillProc` returns `context.Cause(ctx)` when kill RPC fails under cancel (`executor.go:786-790`); original kill error is logged but lost to caller, complicating kill-failure diagnosis.
- **No cgroup kill fallback:** OOM/cgroup code is read-only accounting (`executor_linux.go:208-247`, `executor/resources/monitor.go:19-70`); a container runtime that fails `runc kill` has no freezer-based second line.
- **History pubsub close races shutdown:** close is deferred to finalizer goroutines (`buildhistory.go:500-513`); listeners on a draining daemon may hang until the last build finalizes with no timeout.

## Future Considerations

- Add a bounded `Wait` after the post-kill warn (e.g. 30s) that escalates to container `Delete --force` / cgroup-freezer kill and surfaces `failed to reap after SIGKILL` instead of blocking forever.
- Introduce a first-class `Cancel(ref)` control API that records intent in the history/job table so cancel survives client disconnect and re-attach.
- Give `GracefulStop` a `--shutdown-timeout` that transitions drain → cancel (call `Receiver.Cancel` on remaining pipes) → `server.Stop()`, mirroring the existing 1-2-3 signal ladder with time instead of signals.
- Replace substring cancel detection with typed `ExitError`/`CancelError` propagation from containerd shims.

## Questions / Gaps

- No evidence found for cgroup/cgroupns-based kill sweeps (searched `Kill`, `cgroup`, `Setpgid`, `Pdeathsig` across repo; only accounting + git/runc snippets above). Confirm intended reliance on runtime `kill` vs kernel cgroup pressure.
- No evidence found for daemon-interrupt (cancel-in-flight) mode on shutdown — only `GracefulStop` + signal ladder. Open: is stalling shutdown on wedged builds accepted behavior?
- No evidence found for durable cancel replay after `buildkitd` restart (searched `jobs.go`, `history/`, `control.go`). Job table is in-memory (`jl.jobs`/`jl.actives`); restart drops all cancel state by design — flag if orchestrator layer assumes otherwise.

---
Generated by `Dimension 01.09: Cancellation, Shutdown, and Process Cleanup` against `buildkit`.
