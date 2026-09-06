# Source Analysis: containerd

## Dimension 01.09: Cancellation, Shutdown, and Process Cleanup

### Source Info

| Field | Value |
|-------|-------|
| Name | containerd |
| Path | `studies/ultraplan-daemon-events-study/sources/containerd` |
| Language / Stack | Go (daemon + per-task shim processes, gRPC + ttrpc, OCI runc, BoltDB metadata) |
| Analyzed | 2026-09-03 |

## Summary

containerd has no durable cancellation primitive. "Cancellation" is an ephemeral OS signal delivered via `Kill` RPCs (`sources/containerd/plugins/services/tasks/local.go:473`, `sources/containerd/cmd/containerd-shim-runc-v2/task/service.go:493`, `sources/containerd/cmd/containerd-shim-runc-v2/runc/container.go:413`) that fan out to `runc kill [--all]` (`sources/containerd/cmd/containerd-shim-runc-v2/process/init.go:367`). There is no persisted cancel intent, no cancel queue, and no poll loop: if the shim is unreachable the `Kill` fails (`NotFound`/`Unavailable`/ttrpc-closed) and the intent is lost. gRPC context cancellation is transport-level only — it aborts the in-flight RPC, not the task.

What containerd does well is bounded cleanup after the parent context is gone. The codebase systematically detaches cleanup from the caller's context with `context.WithoutCancel` plus a fixed timeout: generic `cleanup.Do` (10s, `sources/containerd/internal/cleanup/context.go:27`), shim `cleanupTimeout` (5s) and `shutdownTimeout` (3s) (`sources/containerd/core/runtime/v2/shim.go:62-77`), `cleanupShimTask` / `cleanupShim` / mount deactivate paths (`sources/containerd/core/runtime/v2/shim.go:200`, `sources/containerd/core/runtime/v2/shim_manager.go:462`, `sources/containerd/core/runtime/v2/task_manager.go:193`). Daemon shutdown is detach-and-close, not drain-or-cancel: `SIGTERM`/`SIGINT` triggers `cancel()` + `Server.Stop()` which closes plugins LIFO (`sources/containerd/cmd/containerd/command/main_unix.go:37`, `sources/containerd/cmd/containerd/server/server.go:282`); running shims/containers survive because shims are separate `Setpgid` processes (`sources/containerd/cmd/containerd-shim-runc-v2/manager/manager_linux.go:110`, `sources/containerd/pkg/shim/util_unix.go:53`). Automatic `SIGTERM`→`SIGKILL` escalation exists only in the CRI `StopContainer` path (`sources/containerd/internal/cri/server/container_stop.go:165`), not in the raw tasks API. Process-group/cgroup containment is partial: image verifiers get their own process group with negative-PID group kill (`sources/containerd/pkg/imageverifier/bindir/processes_unix.go:35`), shim spawn uses `Setpgid` (`sources/containerd/pkg/shim/util_unix.go:53`), and container-wide kill depends on the caller passing `All=true` through to `runc kill --all` (`sources/containerd/cmd/containerd-shim-runc-v2/process/init.go:375`).

## Rating

**6 / 10** — Present but inconsistent and fragile on durability.

Rationale: bounded, detached cleanup with explicit timeouts is applied consistently (generic helper + shim/mount/shutdown paths) and is covered by focused tests (`sources/containerd/internal/cleanup/context_test.go:27`, `sources/containerd/internal/cri/server/container_stop_test.go:69`, `sources/containerd/pkg/shim/shim_windows_test.go:574`). But cancellation itself is non-durable (ephemeral signal, lost when no shim is attached), daemon shutdown does not drain or cancel work (detach semantics), kill escalation is layered inconsistently (CRI only), and completion-vs-cancel races are resolved by documented duplicate-event tolerance rather than arbitration. That maps to the 4–6 band ("present but inconsistent, weakly documented, or fragile"), at its top end because the timeout/detach machinery is explicit and tested.

## Evidence Collected

Every entry MUST include a file path with line numbers. Format: `path/to/file.ts:NN`.

| Area | Evidence | File:Line |
|------|----------|-----------|
| Cancellation entry: tasks Kill RPC | `local.Kill` resolves task (or exec sub-process) and calls `p.Kill(ctx, r.Signal, r.All)` — pure signal forwarding, no persisted intent, no polling loop | sources/containerd/plugins/services/tasks/local.go:473 |
| Cancellation entry: shim Kill RPC | `service.Kill` looks up container and delegates to `container.Kill(ctx, r)`; returns immediately after signal delivery | sources/containerd/cmd/containerd-shim-runc-v2/task/service.go:493 |
| Cancellation entry: container dispatch | `Container.Kill` resolves `ExecID` process and calls `p.Kill(ctx, r.Signal, r.All)` | sources/containerd/cmd/containerd-shim-runc-v2/runc/container.go:413 |
| Cancellation terminal: init kill | `Init.kill` calls `runtime.Kill(ctx, id, signal, {All})` and maps errors via `checkKillError` | sources/containerd/cmd/containerd-shim-runc-v2/process/init.go:367 |
| Cancellation terminal: group kill | `Init.KillAll` forces `SIGKILL` with `All:true` through runc | sources/containerd/cmd/containerd-shim-runc-v2/process/init.go:375 |
| Kill error normalization (idempotent cancel) | `checkKillError` maps "already finished / container not running / no such process / ESRCH" to `ErrNotFound`, "does not exist" to `ErrNotFound` | sources/containerd/cmd/containerd-shim-runc-v2/process/utils.go:115 |
| State-gated kill (no queue) | Per-state `Kill` impls: `runningState`/`pausedState` forward to `kill`; `stoppedState.Kill` still calls `kill` (ESRCH→NotFound); `deletedState`/`createdState` variants gate behavior — kill against a non-running task is an error, not queued intent | sources/containerd/cmd/containerd-shim-runc-v2/process/init_state.go:272 |
| Kill escalation with deadline (CRI only) | `stopContainer`: send stop signal once (CAS-guarded), wait `timeout` via `sigTermCtx`, on expiry send `SIGKILL`, then bounded `waitContainerStop` | sources/containerd/internal/cri/server/container_stop.go:165 |
| Single-flight stop signal | `atomic.CompareAndSwapUint32(IsStopSignaledWithTimeout, 0, 1)` prevents double SIGTERM on concurrent stops | sources/containerd/internal/cri/server/container_stop.go:205 |
| Parent-cancel short-circuit in escalation | If caller `ctx.Err() != nil` after the TERM wait, return immediately instead of escalating to KILL | sources/containerd/internal/cri/server/container_stop.go:224 |
| Completion/cancel race retry | `stopContainerRetryOnConnectionClosed`: up to 3 retries with 100ms·i² backoff when shim ttrpc is closed (self-exit vs stop race) | sources/containerd/internal/cri/server/container_stop.go:87 |
| Wait bounded by context | `waitContainerStop` selects on `ctx.Done()` vs `container.Stopped()` — cancellation aborts the wait, does not force the kill | sources/containerd/internal/cri/server/container_stop.go:253 |
| Wait cancellation test | `waitContainerStop` tested with pre-cancelled context and timeout contexts | sources/containerd/internal/cri/server/container_stop_test.go:69 |
| Daemon signal policy | Unix handled signals `SIGTERM/SIGINT/SIGUSR1/SIGPIPE`; `SIGUSR1` dumps stacks, `SIGPIPE` ignored, default path notifies stopping, cancels root ctx, calls `server.Stop()`, closes `done` | sources/containerd/cmd/containerd/command/main_unix.go:30 |
| Daemon boot vs signal race | Server construction runs in goroutine so `SIGTERM` during BoltDB open / init is still honored; `select ctx.Done()` vs `chsrv` in two places | sources/containerd/cmd/containerd/command/main.go:219 |
| Daemon stop = LIFO plugin close | `Server.Stop` iterates plugins reverse, `Close()`s each `io.Closer`, logs and continues on error; no timeout, no drain, no task cancel | sources/containerd/cmd/containerd/server/server.go:282 |
| Shutdown service (shim-side) | `WithShutdown` context is immune to parent cancel; `Shutdown()` runs callbacks concurrently via `errgroup` on a detached `Background+30s` context, records first error or `ErrShutdown` | sources/containerd/pkg/shutdown/shutdown.go:51 |
| Shutdown callback registration | `RegisterCallback` appends under mutex; callbacks run only on explicit `Shutdown()` | sources/containerd/pkg/shutdown/shutdown.go:108 |
| Shim shutdown guard | `service.Shutdown` is a no-op while `len(containers) > 0`; otherwise fires `shutdown.Shutdown()` | sources/containerd/cmd/containerd-shim-runc-v2/task/service.go:608 |
| Generic detached cleanup (10s) | `cleanup.Do` wraps fn in `WithoutCancel+WithTimeout(10s)`, preserving values but clearing cancellation | sources/containerd/internal/cleanup/context.go:27 |
| Generic cleanup test | `TestDo` verifies values survive, cancel is cleared, re-cancel inside fn works, and default timeout fires (~10s context) | sources/containerd/internal/cleanup/context_test.go:27 |
| Shim cleanup timeouts | `cleanupTimeout=5s`, `shutdownTimeout=3s`, `loadTimeout=5s` registered in `pkg/timeout` | sources/containerd/core/runtime/v2/shim.go:62 |
| Dead-shim cleanup path | `cleanupAfterDeadShim` uses `timeout.WithContext(ctx, cleanupTimeout)`, calls `binaryCall.Delete`, publishes synthetic `TaskExit(255)` + `TaskDelete` when no response | sources/containerd/core/runtime/v2/shim.go:144 |
| Failed-start cleanup detached | `cleanupShimTask` detaches with `WithoutCancel+cleanupTimeout`, deletes, then re-arms a fresh timeout for `Shutdown` if the first budget expired; notes ttrpc vs gRPC error flattening | sources/containerd/core/runtime/v2/shim.go:200 |
| Shim delete→shutdown→remove ordering | `shimTask.delete`: `Delete` RPC, `removeTask` only on nil error, `waitShutdown` (3s), `ShimInstance.Delete`, final `removeTask`; shutdown failure only logged ("shim might be leaked") | sources/containerd/core/runtime/v2/shim.go:582 |
| Manager disconnect cleanup | Shim `onClose` callback calls `cleanupAfterDeadShim(WithoutCancel(ctx),…)` then removes from map | sources/containerd/core/runtime/v2/shim_manager.go:333 |
| Manager error-path cleanup | `cleanupShim` detaches with `WithoutCancel+cleanupTimeout` before `shim.Delete` | sources/containerd/core/runtime/v2/shim_manager.go:462 |
| Mount deactivate detached | `TaskManager.Create` failure path deactivates mounts under `WithoutCancel+cleanupTimeout` | sources/containerd/core/runtime/v2/task_manager.go:193 |
| Shim spawn isolation | Shim child launched with `Setpgid:true` so it survives daemon death and forms its own group | sources/containerd/cmd/containerd-shim-runc-v2/manager/manager_linux.go:110 |
| Shim util isolation | `getSysProcAttr` returns `Setpgid:true` for shim-spawned helpers | sources/containerd/pkg/shim/util_unix.go:53 |
| Verifier group kill | Verifier command assigned new process group; `Cancel()` sends `SIGKILL` to negative PID (whole group) without killing containerd itself | sources/containerd/pkg/imageverifier/bindir/processes_unix.go:35 |
| Forced manager teardown | `manager.Stop` force-deletes runc container (`Force:true`), unmounts rootfs, reads init pidfile, synthesizes `128+SIGKILL` exit | sources/containerd/cmd/containerd-shim-runc-v2/manager/manager_linux.go:300 |
| Transfer streaming respects cancel | Streaming reader/writer/select loops abort on `ctx.Done()` — client disconnect stops the transfer leg in flight | sources/containerd/core/transfer/streaming/stream.go:51 |
| Shutdown/reap semantics test | `reap` returns `context.Canceled` on ctx cancel but survives `SIGTERM`; `handleExitSignals` runs shutdown callback only on signal, not on ctx cancel | sources/containerd/pkg/shim/shim_windows_test.go:574 |
| Sandbox stop cleanup | `controllerLocal.Shutdown` calls `ShutdownSandbox` then `cleanupShim` on the sandbox shim | sources/containerd/plugins/sandbox/controller.go:106 |

## Answers to Dimension Questions

1. **Is client disconnect different from explicit cancellation?** No meaningful difference at the daemon layer. Both surface as gRPC/ttrpc context cancellation (`ctx.Done()`/`ctx.Err()`). Transfers and waits abort the in-flight leg (`sources/containerd/core/transfer/streaming/stream.go:51`, `sources/containerd/internal/cri/server/container_stop.go:253`), but neither stops the underlying container: a disconnected `Kill` RPC that already delivered its signal still killed; a disconnected `StopContainer` before signal delivery leaves the container running. There is no separate "disconnect" vs "cancel" intent type — verified by absence of any disconnect-specific handler in `plugins/services/tasks/local.go:473` and `cmd/containerd-shim-runc-v2/task/service.go:493` (searched Kill/Delete paths; no disconnect branch found).

2. **Is cancellation durable if no worker is currently attached?** No. `Kill` requires resolving a live task/shim (`sources/containerd/plugins/services/tasks/local.go:473`, `sources/containerd/cmd/containerd-shim-runc-v2/task/service.go:493`). If the shim is gone, the call fails (`NotFound` via `sources/containerd/cmd/containerd-shim-runc-v2/process/utils.go:115`, ttrpc-closed, or `Unavailable`) and the intent is dropped. The only durability-adjacent mechanism is the dead-shim synthesizer (`sources/containerd/core/runtime/v2/shim.go:144`), which publishes exit/delete events after the fact — it does not replay a missed kill. No evidence found of a cancel journal, pending-intent table, or poll loop.

3. **Can cleanup hang indefinitely?** By design, no on the paths that matter — every detached cleanup carries a fixed budget: 10s generic (`sources/containerd/internal/cleanup/context.go:27`), 5s shim cleanup / 3s shim shutdown (`sources/containerd/core/runtime/v2/shim.go:62`), 30s shutdown-callback fan-out (`sources/containerd/pkg/shutdown/shutdown.go:79`). `waitTimeout` for IO also bounds waits (`sources/containerd/cmd/containerd-shim-runc-v2/process/utils.go:157`). Two exceptions: `Server.Stop` plugin `Close()` loop has no timeout (`sources/containerd/cmd/containerd/server/server.go:282`) — a wedged plugin `Close` stalls daemon exit; and shim `delete` logs-and-continues on shutdown failure rather than retrying (`sources/containerd/core/runtime/v2/shim.go:616`), which bounds time at the cost of a possible leaked shim.

4. **Can child/grandchild processes escape termination?** Yes, depending on which API is used. Raw `Kill` with `All=false` signals only the init/exec PID (`sources/containerd/cmd/containerd-shim-runc-v2/process/init.go:367`); full-tree kill requires `All=true` → `runc kill --all` (cgroup-based, `sources/containerd/cmd/containerd-shim-runc-v2/process/init.go:375`). Processes that leave the container cgroup, ignore/mask the signal, or are in uninterruptible sleep survive a single `Kill` — the shim performs no automatic escalation. Only CRI `StopContainer` escalates `TERM→KILL` after its grace timeout (`sources/containerd/internal/cri/server/container_stop.go:165`). Daemon-level isolation (`Setpgid`, `sources/containerd/cmd/containerd-shim-runc-v2/manager/manager_linux.go:110`) protects the shim from the daemon's signals but does not contain container descendants; only the verifier path explicitly group-kills (`sources/containerd/pkg/imageverifier/bindir/processes_unix.go:40`).

5. **How are cancellation and completion races arbitrated?** Ad-hoc, per layer, favoring idempotence + duplicate tolerance over exactly-once arbitration: (a) CRI stop retries 3× on ttrpc-closed to cover self-exit-vs-stop (`sources/containerd/core/runtime/v2/shim.go:87`-adjacent, `sources/containerd/internal/cri/server/container_stop.go:87`); (b) `checkKillError` converts post-exit kills to `NotFound` so kill-after-exit is benign (`sources/containerd/cmd/containerd-shim-runc-v2/process/utils.go:115`); (c) `shimTask.delete` only removes the task record on nil error but always re-removes at the end, and explicitly documents that duplicate `TaskExit`/`TaskDelete` events are possible and downstream (moby) must tolerate them (`sources/containerd/core/runtime/v2/shim.go:596`); (d) single-flight stop-signal CAS prevents double-TERM but does not order TERM vs natural exit (`sources/containerd/internal/cri/server/container_stop.go:205`). No central race arbiter or fencing token was found.

**Terminal scenario (cancel lands exactly as work completes; child ignores first signal):** the `Kill(SIGTERM, All=false)` either returns success (signal delivered to an already-exiting PID — harmless) or `NotFound` (process already reaped — benign per `sources/containerd/cmd/containerd-shim-runc-v2/process/utils.go:115`). Because the raw tasks API does no escalation, a child ignoring `SIGTERM` keeps running until someone issues a second `Kill(SIGKILL, All=true)` or deletes the task; only a CRI `StopContainer` caller would get automatic escalation after its `timeout` (`sources/containerd/internal/cri/server/container_stop.go:217`). Completion wins the event race by tolerance, not arbitration: the exit path publishes `TaskExit`/`TaskDelete` (or synthetic 255 on dead shim, `sources/containerd/core/runtime/v2/shim.go:144`), the late cancel resolves to `NotFound`, and any duplicate exit notifications are the consumer's problem (`sources/containerd/core/runtime/v2/shim.go:596`). Cleanup still runs bounded: delete/shutdown under 5s/3s budgets (`sources/containerd/core/runtime/v2/shim.go:200`), mount deactivate detached (`sources/containerd/core/runtime/v2/task_manager.go:193`).

## Architectural Decisions

- **Ephemeral signals, not durable cancel records.** `Kill`/`Delete`/`Shutdown` are stateless RPCs over volatile shim connections (`sources/containerd/plugins/services/tasks/local.go:473`, `sources/containerd/cmd/containerd-shim-runc-v2/task/service.go:493`). Tradeoff: simple, no cancel journal to reconcile — but intent is lost across shim/daemon restarts.
- **Shims outlive the daemon (detach on shutdown).** `SIGTERM` stops plugins, not workloads (`sources/containerd/cmd/containerd/command/main_unix.go:37`, `sources/containerd/cmd/containerd/server/server.go:282`); shims are separate `Setpgid` processes (`sources/containerd/cmd/containerd-shim-runc-v2/manager/manager_linux.go:110`) reloaded on restart. Tradeoff: fast, non-disruptive daemon restarts — but shutdown provides no workload quiescence guarantee.
- **Detach cleanup from caller fate.** The `WithoutCancel + fixed timeout` idiom is the standard firewall between request lifetime and resource lifetime (`sources/containerd/internal/cleanup/context.go:27`, `sources/containerd/core/runtime/v2/shim.go:200`, `sources/containerd/core/runtime/v2/shim_manager.go:462`). Tradeoff: cleanup survives client disconnect — but a too-short budget (3–5s) can leave a shim leaked (explicitly logged, `sources/containerd/core/runtime/v2/shim.go:616`).
- **Escalation pushed to the highest layer.** Only CRI `StopContainer` implements `TERM→KILL` with a caller-supplied grace period (`sources/containerd/internal/cri/server/container_stop.go:165`); the tasks/shim/runc layers deliver exactly one signal per call. Tradeoff: policy lives with the orchestrator — but raw API users must reimplement escalation themselves.
- **Duplicate-tolerant events instead of exactly-once arbitration.** Exit/delete races are resolved by allowing redelivery and requiring consumers to dedupe (`sources/containerd/core/runtime/v2/shim.go:596`).

## Notable Patterns

- **WithoutCancel-then-timeout** as the universal cleanup wrapper: generic helper (`sources/containerd/internal/cleanup/context.go:27`), shim paths (`sources/containerd/core/runtime/v2/shim.go:200`), manager paths (`sources/containerd/core/runtime/v2/shim_manager.go:462`), mount paths (`sources/containerd/core/runtime/v2/task_manager.go:193`).
- **Error normalization to idempotent outcomes:** `checkKillError` ESRCH→`NotFound` (`sources/containerd/cmd/containerd-shim-runc-v2/process/utils.go:115`); `Delete` on missing task → `NotFound` (`sources/containerd/plugins/services/tasks/local.go:318`).
- **State-machine gating of destructive ops:** per-state `Kill`/`Delete` in `sources/containerd/cmd/containerd-shim-runc-v2/process/init_state.go:272` (e.g. `runningState.Delete` refuses, `stoppedState.Delete` proceeds).
- **Negative-PID group kill for helper trees:** verifier `Cancel` kills the whole process group (`sources/containerd/pkg/imageverifier/bindir/processes_unix.go:40`) — the one place group termination is unconditional.
- **Boot-path signal safety:** server construction off the signal path with dual `select ctx.Done()` gates (`sources/containerd/cmd/containerd/command/main.go:219`).

## Tradeoffs

- **Speed of daemon restart vs workload safety:** detach-on-shutdown keeps containers running across daemon upgrades, but `SIGTERM` to containerd is not a workload drain — orchestrators must stop containers first.
- **Bounded cleanup vs completeness:** 3–10s budgets guarantee the caller returns, but slow unmounts, wedged runtimes, or unresponsive shims produce leaks-by-log (`sources/containerd/core/runtime/v2/shim.go:616`) rather than blocking.
- **One-signal-per-RPC vs automatic escalation:** predictable low-level semantics, but every higher-level caller (CRI did, others must) owns the TERM→KILL timer and the `All` flag decision.
- **`All=true` power vs blast radius:** cgroup-wide kill closes the escape hatch for descendants but can kill execs the caller did not intend; default `All=false` is surgical but escapable.
- **Shutdown-callback parallelism (`errgroup`, 30s) vs ordering:** concurrent plugin/shim callbacks (`sources/containerd/pkg/shutdown/shutdown.go:79`) shut down fast but provide no dependency order; first error wins, rest are dropped.

## Failure Modes / Edge Cases

- **Kill during shim disconnect:** signal lost; caller sees ttrpc-closed/`Unavailable`. CRI stop retries 3× (`sources/containerd/internal/cri/server/container_stop.go:87`); raw `Kill` does not retry — caller must poll `State`/`Wait`.
- **Child ignores SIGTERM via raw API:** container lingers indefinitely; no in-shim watchdog escalates. Mitigation exists only in CRI grace path (`sources/containerd/internal/cri/server/container_stop.go:217`).
- **Grandchild escapes cgroup (`All=false` or cgroup-evasive fork):** survives `Kill`; reclaim requires host-level `runc kill --all`, cgroup removal, or `manager.Stop` force-delete (`sources/containerd/cmd/containerd-shim-runc-v2/manager/manager_linux.go:300`).
- **Cleanup budget expiry:** `waitShutdown` failure leaves shim process + socket + bundle; only a log line records it (`sources/containerd/core/runtime/v2/shim.go:616`). Stale sockets are partially mitigated by `RemoveSocket`/`cleanupSockets` (`sources/containerd/pkg/shim/util_unix.go:303`).
- **Daemon `Stop` hang:** a blocking plugin `Close()` stalls `SIGTERM` handling with no timeout (`sources/containerd/cmd/containerd/server/server.go:282`); recovery requires `SIGKILL` of the daemon (shims survive, per `Setpgid` design).
- **Double-stop stampede:** mitigated for CRI by the `IsStopSignaledWithTimeout` CAS (`sources/containerd/internal/cri/server/container_stop.go:205`); raw concurrent `Kill` calls each deliver a signal (harmless but noisy).
- **Exit/cancel duplicate delivery:** consumers receive two `TaskExit` events in the self-exit-vs-delete window; documented as consumer's burden (`sources/containerd/core/runtime/v2/shim.go:596`).

## Future Considerations

- Add an escalation helper at the tasks/shim layer (e.g. `Kill` with `SIGTERM` + deadline + automatic `SIGKILL All=true`) so non-CRI clients get the same guarantee currently confined to `sources/containerd/internal/cri/server/container_stop.go:165`.
- Consider a persisted or at least in-memory pending-kill record replayed on shim reconnect, to close the "signal lost while detached" gap in `sources/containerd/plugins/services/tasks/local.go:473`.
- Bound `Server.Stop` plugin closes with a per-plugin timeout (mirroring the 30s `errgroup` pattern in `sources/containerd/pkg/shutdown/shutdown.go:79`) so daemon `SIGTERM` cannot hang forever.
- Extend the negative-PID group-kill pattern from `sources/containerd/pkg/imageverifier/bindir/processes_unix.go:40` (or cgroup `kill`) to container teardown as a backstop for cgroup escapees, with explicit opt-in to control blast radius.
- Replace duplicate-tolerant exit delivery (`sources/containerd/core/runtime/v2/shim.go:596`) with idempotency keys or sequence numbers if exactly-once consumers emerge.

## Questions / Gaps

- No evidence found for a daemon-wide task drain policy on `SIGTERM` (searched `cmd/containerd/command`, `cmd/containerd/server/server.go:282`): is detach-always the intended contract, or is a `--drain-timeout` planned?
- No evidence found for cgroup ` freezer`/`kill` fallback inside the runc shim kill path beyond `All → runc kill --all` (`sources/containerd/cmd/containerd-shim-runc-v2/process/init.go:367`): which runc versions guarantee `--all` semantics, and what happens on cgroup v1 vs v2 when it is unavailable?
- No evidence found for cancel propagation into content fetch/unpack beyond transport abort (`sources/containerd/core/transfer/streaming/stream.go:51`): does cancelling an image pull leave partial content/locks requiring GC, and which component owns that reclamation?
- The 3s `shutdownTimeout` vs 5s `cleanupTimeout` relationship (`sources/containerd/core/runtime/v2/shim.go:62`) is unexplained: is 3s sufficient for a loaded shim to exit, and what telemetry counts `failed to shutdown shim task` leaks?

---

Generated by `Dimension 01.09: Cancellation, Shutdown, and Process Cleanup` against `containerd`.
