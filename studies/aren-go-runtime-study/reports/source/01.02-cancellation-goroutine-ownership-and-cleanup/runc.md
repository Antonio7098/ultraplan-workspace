# Source Analysis: runc

## 01.02 Cancellation, Goroutine Ownership, and Cleanup

### Source Info

| Field | Value |
|-------|-------|
| Name | runc |
| Path | `studies/aren-go-runtime-study/sources/runc` |
| Language / Stack | Go 1.26 (`go.mod`), Linux-only runtime; urfave/cli v3; containerd/console; CRIU for checkpoint/restore |
| Analyzed | 2026-08-26 |
| Version analyzed | `1.5.0-rc.1+dev`, commit `9674194` (`sources/runc/VERSION`) |

Citation convention: all paths below are relative to `studies/aren-go-runtime-study/`, e.g. `sources/runc/signals.go:18`. Vendored dependencies live inside the source tree and are cited as `sources/runc/vendor/...`.

## Summary

runc is a short-lived CLI that supervises one container process per invocation, so its cancellation model is **signal-driven, not context-driven**. The root context is `context.Background()` (`sources/runc/main.go:184`) and every cli `Action` discards its context parameter (e.g. `sources/runc/run.go:77`, `sources/runc/create.go:63`). The only genuine `context.WithTimeout` usage in the codebase is a mount-source worker goroutine bounded to 1 minute (`sources/runc/libcontainer/process_linux.go:692`). Termination of the supervisor is instead negotiated through Unix mechanisms: signals are forwarded from the `runc run` parent to the container init (`sources/runc/signals.go:93-98`), child exits are reaped via a SIGCHLD-driven `wait4(-1)` loop (`sources/runc/signals.go:106-126`), and the terminal publication is the return value of the forward loop converted into `os.Exit(status)` (`sources/runc/run.go:81-87`).

Goroutine ownership is deliberately asymmetric: three IO-copy goroutines are tracked with a `sync.WaitGroup` and joined in an ordered `tty.Close()` (`sources/runc/tty.go:170-187`), but at least five goroutine families are intentionally *detached* — the epoller loop (`sources/runc/tty.go:128`), stdin copiers (`sources/runc/tty.go:56-59`, `sources/runc/tty.go:129`), the interrupt hard-exit handler (`sources/runc/tty.go:145-151`), the events stats ticker (`sources/runc/events.go:77-86`), and cgroup memory-event watchers (`sources/runc/libcontainer/notify_linux.go:41-59`). These are bounded by process lifetime, which is acceptable for a per-command CLI but would be a liability if the same patterns were embedded in a long-lived supervisor.

The strongest engineering is in **cleanup ordering and truthful completion**: `initProcess.start` has a deferred failure path that detects OOM kills before tearing down the cgroup, SIGKILLs and waits for init, and destroys managers while preserving the primary error (`sources/runc/libcontainer/process_linux.go:789-821`); log-forwarder fatal messages from `runc init` are merged into the returned error rather than overwriting it (`sources/runc/libcontainer/container_linux.go:406-421`, errors combined with `%w; %w`); and `Destroy()` is a state-machine operation guarded by a mutex, refusing running/paused containers and reconciling in-memory state against live kernel facts (PID start-time checks defeat PID reuse) before removing resources (`sources/runc/libcontainer/state_linux.go:39-67`, `sources/runc/libcontainer/container_linux.go:919-952`). There are no goroutine-leak tests (no goleak dependency); leak coverage is fd-based (`sources/runc/libcontainer/integration/exec_test.go:1693-1778`) plus behavioral kill tests for shared-PID-namespace leftover processes (`sources/runc/libcontainer/integration/exec_test.go:1371-1427`, `sources/runc/tests/integration/delete.bats:13-67`).

## Rating

**6 / 10.**

Rationale against the dimension's focus areas:

- **Cancellation propagation (weak):** no caller-to-operation cancellation contract; contexts are created only as `context.Background()`/`context.WithTimeout` inside one helper (`sources/runc/main.go:184`, `sources/runc/libcontainer/process_linux.go:879,692`). All "cancellation" is signal forwarding.
- **Goroutine ownership (mixed):** explicit tracked ownership for console/pipe IO (`sources/runc/tty.go:22,60-62,130-131`) and synchronous join for hook/BPF/mount workers (`sources/runc/libcontainer/configs/config.go:626-632`, `sources/runc/libcontainer/seccomp/patchbpf/enosys_linux.go:158-170`, `sources/runc/libcontainer/process_linux.go:751`), but several untracked senders can block forever on unbuffered channels (`sources/runc/notify_socket.go:133`, `sources/runc/libcontainer/notify_v2_linux.go:67`) and `Epoller.Close()` is never called anywhere, leaving the epoller goroutine parked in `epoll_wait(-1)` (`sources/runc/tty.go:128`, `sources/runc/tty.go:178`; confirmed by grep — only `CloseConsole` is used).
- **Bounded shutdown (good):** hook timeouts (`sources/runc/libcontainer/configs/config.go:620-633`), 1-minute mount-worker deadline (`sources/runc/libcontainer/process_linux.go:692`), 30s sd_notify barrier deadline (`sources/runc/notify_socket.go:227`), ~10 s bounded poll in `runc delete --force` (`sources/runc/delete.go:17-26`), pidfd-accelerated fifo wait (`sources/runc/libcontainer/container_linux.go:284-303`).
- **Cleanup ordering & failure preservation (strong):** ordered LIFO defers in `runner.run` (`sources/runc/utils_linux.go:225-231,278,301-306`), ordered `tty.Close` (`sources/runc/tty.go:172-186`), state-machine destroy with thaw-before-kill and kill-before-cgroup-removal ordering (`sources/runc/libcontainer/state_linux.go:160-163,186-194`, `sources/runc/libcontainer/init_linux.go:684-720`).
- **Truthful completion (strong):** exit status published only after pipe flush (`sources/runc/signals.go:81-85`), init fatals merged into primary error (`sources/runc/libcontainer/container_linux.go:410-419`), empty-exec-fifo reported as "init did not signal execve readiness" (`sources/runc/libcontainer/container_linux.go:335-337`).

Not a 7–8 because cooperative cancellation plumbing is essentially absent, multiple channel sends are unguarded against consumer disappearance, shutdown idempotency re-runs poststop hooks on repeat `Destroy()` (`sources/runc/libcontainer/state_linux.go:64-66,104-106`), and there is no goroutine-leak test discipline.

## Evidence Collected

Every entry cites `path:line` relative to `studies/aren-go-runtime-study/`.

| Area | Evidence | File:Line |
|------|----------|-----------|
| Root context | `app.Run(context.Background(), os.Args)`; all Action signatures discard ctx | `sources/runc/main.go:184`; `sources/runc/run.go:77`; `sources/runc/create.go:63`; `sources/runc/events.go:35` |
| Only real context use | `context.WithTimeout(ctx, 1*time.Minute)` + `context.AfterFunc(ctx, ...)` closing request channel | `sources/runc/libcontainer/process_linux.go:692-693` |
| Mount worker goroutine | Locked OS thread (no unlock → thread destroyed), setns into container mntns, select on requests vs ctx.Done | `sources/runc/libcontainer/process_linux.go:695-748` |
| Worker lifecycle | `goCreateMountSources(context.Background())` + `defer cancel()` in `initProcess.start` | `sources/runc/libcontainer/process_linux.go:877-885` |
| Signal handler setup | Buffer-2048 signal chan; SIGCHLD registered eagerly; all signals registered on background goroutine, published via chan | `sources/runc/signals.go:18-38` |
| Forward loop termination | On reap of pid1: `process.Wait()` then return status ("flush pipes" comment) | `sources/runc/signals.go:70-87` |
| Signal forwarding default | Non-SIGCHLD/SIGWINCH/SIGURG forwarded via `process.Signal(s)` | `sources/runc/signals.go:93-99` |
| Reaping | `wait4(-1, WNOHANG)` loop until ECHILD/0 | `sources/runc/signals.go:106-126` |
| Supervisor orchestration | `runner.run`: deferred destroy when non-detach or error; deferred SIGKILL+Wait terminate-on-error | `sources/runc/utils_linux.go:223-231,301-306` |
| Subreaper ownership | `SetSubreaper(1)` before signal-handler registration unless `--no-subreaper` | `sources/runc/utils_linux.go:264-269` |
| Handler publication wait | `handler := <-handlerCh` before `forward(process, tty)`; not created at all when detached | `sources/runc/utils_linux.go:270-273,325-326` |
| tty ownership struct | `closers`, `postStart`, `sync.WaitGroup`, `consoleC` fields | `sources/runc/tty.go:16-24` |
| Tracked IO goroutines | `t.wg.Add(2)` + stdout/stderr `copyIO`; `t.wg.Add(1)` + epoll console stdout copy | `sources/runc/tty.go:60-62,130-131` |
| Untracked stdin copier | `io.Copy(i.Stdin, os.Stdin)` without wg accounting | `sources/runc/tty.go:56-59` |
| Untracked epoller/stdin-console | `go func() { _ = epoller.Wait() }()`, `io.Copy(epollConsole, os.Stdin)` | `sources/runc/tty.go:128-129` |
| Hard-exit escape hatch | `handleInterrupt` goroutine: on SIGINT reset console then `os.Exit(0)`, bypassing all defers | `sources/runc/tty.go:137,145-151` |
| Console teardown order | postStart close → `console.Shutdown(epoller.CloseConsole)` → `wg.Wait()` → closers → host console reset | `sources/runc/tty.go:170-187` |
| Console shutdown mechanism | `Shutdown` broadcasts read/write cond vars, marks closed; Read exits when `closed && n==0`; Epoller.Wait loops on `EpollWait(-1)`; `Epoller.Close` never called by runc | `sources/runc/vendor/github.com/containerd/console/console_linux.go:257-267,196,109-137,155-161` |
| Console setup handoff | `consoleC` buffered(1); recvtty runs in its own goroutine; consumed by `tty.waitConsole()` | `sources/runc/utils_linux.go:116-119`; `sources/runc/tty.go:153-158`; `sources/runc/utils_linux.go:307-309` |
| Log forwarder goroutine | `ForwardLogs` scans pipe to EOF, closes pipe, sends+close buffered(1) done chan | `sources/runc/libcontainer/logs/logs.go:17-52` |
| Log error merging | deferred `<-logsDone` merges `runc init` fatal logs into retErr without dropping primary error | `sources/runc/libcontainer/container_linux.go:406-421` |
| comm fd hygiene | `closeChild`/`closeParent`; `logPipeParent` intentionally kept alive for ForwardLogs | `sources/runc/libcontainer/process_linux.go:106-116` |
| init start failure path | OOM count read *before* kill; `terminate()` = Kill+Wait; `manager.Destroy` + intelRdt destroy; cleanup errors downgraded to warnings | `sources/runc/libcontainer/process_linux.go:789-821` |
| Benign terminate errors | `ignoreTerminateErrors` filters ExitError/ErrProcessDone/"Wait was already called" | `sources/runc/libcontainer/container_linux.go:1209-1228` |
| Init sync protocol | bootstrap copy → getChildPid → waitForChildExit (replaces cmd.Process) → config JSON → parseSync loop → SHUT_WR → procReady check | `sources/runc/libcontainer/process_linux.go:848-867,904-1063`; `sources/runc/libcontainer/process_linux.go:654-672` |
| Crash-safety state write | State saved *before* `procRun` sync so `runc delete --force` can clean leaked stage-2 init (explicit comment) | `sources/runc/libcontainer/process_linux.go:997-1010` |
| Atomic state save | temp file + `os.Rename` onto state.json | `sources/runc/libcontainer/container_linux.go:882-906` |
| Exec fifo wait | `handleFifo`: O_NONBLOCK open, pidfd+fifo single poll (no timeout), fallback 100 ms polling with zombie check; fifo removed after read | `sources/runc/libcontainer/container_linux.go:247-327,262` |
| Empty fifo truthfulness | "exec fifo is empty: container init did not signal execve readiness" | `sources/runc/libcontainer/container_linux.go:329-339` |
| Destroy API contract | Doc: running containers must first be stopped using Signal; paused must be resumed; "No error ... if already destroyed" claim | `sources/runc/libcontainer/container_linux.go:785-802` |
| Destroy implementation | Mutex-guarded dispatch to `state.destroy()` | `sources/runc/libcontainer/container_linux.go:795-802` |
| destroy() ordering | shared-pidns SIGKILL-all → cgroup destroy → IntelRDT destroy → RemoveAll(stateDir) → nil initProcess → poststop hooks → stoppedState | `sources/runc/libcontainer/state_linux.go:39-67` |
| Per-state guards | running refuses with ErrRunning; paused thaws first or ErrPaused; created SIGKILLs init first; loaded refreshes state then re-dispatches; stopped re-runs full destroy | `sources/runc/libcontainer/state_linux.go:134-139,186-194,160-163,240-245,104-106` |
| Live state reconciliation | `refreshState` transitions based on freezer state, liveness, exec-fifo presence | `sources/runc/libcontainer/container_linux.go:919-936` |
| PID-reuse defense | `hasInit` compares `/proc/<pid>` StartTime and rejects Zombie/Dead; `signalInit` refuses to signal non-running container | `sources/runc/libcontainer/container_linux.go:939-952,473-477` |
| Whole-cgroup kill | SIGKILL with shared PID ns → `signalAllProcesses`; cgroup.kill fast path; freeze→enumerate→kill→thaw | `sources/runc/libcontainer/container_linux.go:452-467`; `sources/runc/libcontainer/init_linux.go:684-720` |
| Thaw-before-SIGKILL | frozen cgroup thawed so SIGKILL lands (cgroup v1) | `sources/runc/libcontainer/container_linux.go:481-488` |
| Bounded forced delete | SIGKILL then 100×100 ms Signal(0) polling, else "container init still running" | `sources/runc/delete.go:17-26` |
| Hook timeout | waiter goroutine + timer; on timeout Kill then `<-errC` join | `sources/runc/libcontainer/configs/config.go:612-633` |
| BPF pipe copier | errChan buffered(1)+closed; writer closed for EOF; joined synchronously | `sources/runc/libcontainer/seccomp/patchbpf/enosys_linux.go:156-171` |
| CPU-affinity starter | dedicated LockOSThread goroutine for cmd.Start; deliberate no UnlockOSThread | `sources/runc/libcontainer/process_linux.go:214-239` |
| CRIU child mgmt | criuServer closed right after Start so CRIU crash cannot hang runc; deferred Wait on error path; CloseWrite+Wait on success | `sources/runc/libcontainer/criu_linux.go:924-952,1059-1065` |
| Restore process swap | `post-restore` notification replaces `cmd.Process` with restored init; `restoredProcess.wait` has open TODO about waiting actual process | `sources/runc/libcontainer/criu_linux.go:1145-1159`; `sources/runc/libcontainer/restored_process.go:47-58` |
| sd_notify reader | reader goroutine sends READY line on unbuffered chan; main loop exits on pid1 death via 100 ms ticker | `sources/runc/notify_socket.go:117-151` |
| Detached notify run | non-detach: synchronous `run(os.Getpid())` then detached `go run(0)` (pid 0 stat fails → self-terminates) | `sources/runc/notify_socket.go:154-168` |
| Events fan-out | encoder via `WaitGroup.Go`; untracked ticker goroutine using `time.Tick` sending to stats chan; main loop ends when OOM watcher closes | `sources/runc/events.go:54-110` |
| Memory event watchers v1/v2 | eventfd/inotify goroutines close fds+channel on return; termination tied to cgroup destruction/populated==0; unguarded unbuffered sends | `sources/runc/libcontainer/notify_linux.go:40-59`; `sources/runc/libcontainer/notify_v2_linux.go:31-77` |
| Leak tests (fd only) | TestFdLeaks/TestFdLeaksSystemd diff /proc/self/fd before/after Run with GC disabled and whitelist | `sources/runc/libcontainer/integration/exec_test.go:1693-1778` |
| Leftover-process kill tests | TestHostPidnsInitKill/TestSharedPidnsInitKill require both processes dead after Signal(SIGKILL) | `sources/runc/libcontainer/integration/exec_test.go:1371-1427` |
| Forwarding stop test | `TestLogForwardingStopsAfterClosingTheWriter` bounds forwarding stop with 10 s deadline | `sources/runc/libcontainer/logs/logs_linux_test.go:53-68` |
| State transition unit tests | stopped/paused/restored/running/created transition tests incl. invalid-transition errors | `sources/runc/libcontainer/state_linux_test.go:52-110` |
| Integration delete semantics | host-pidns test kills init externally, requires leftover processes killed and container removed | `sources/runc/tests/integration/delete.bats:13-67` |
| No goleak | `grep -rn goleak` over go.mod/go.sum and non-vendor sources returns nothing | (search boundary stated in Questions/Gaps) |

## Answers to Dimension Questions

### Who is responsible for waiting for every goroutine started on behalf of a run?

Nobody globally; responsibility is split per-subsystem and partly delegated to process exit:

- The `tty.WaitGroup` covers exactly three goroutines: stdout/stderr pipe copiers and the epoll-console stdout copier (`sources/runc/tty.go:60-62,130-131`), joined in `tty.Close()` after console shutdown (`sources/runc/tty.go:180`).
- Synchronous helpers join their workers inline: hook runner joins via `<-errC` even on timeout (`sources/runc/libcontainer/configs/config.go:626-632`), BPF disassembler joins its pipe copier (`sources/runc/libcontainer/seccomp/patchbpf/enosys_linux.go:170`), and `goCreateMountSources` blocks until the worker's setup phase succeeds or fails (`sources/runc/libcontainer/process_linux.go:750-755`).
- Deliberately unowned/detached: epoller.Wait loop (`sources/runc/tty.go:128`), stdin copiers (`sources/runc/tty.go:57,129`), `handleInterrupt` (`sources/runc/tty.go:137`), notify-socket readers (`sources/runc/notify_socket.go:121`), events ticker (`sources/runc/events.go:77`), and memory-event watchers until cgroup death (`sources/runc/libcontainer/notify_linux.go:41`). For child *processes*, the CLI makes itself subreaper by default so reparented children remain reapable (`sources/runc/utils_linux.go:264-269`), and the library documents that a SIGKILL into a shared-PID-namespace container obliges the libcontainer user to implement a proper child reaper (`sources/runc/libcontainer/container_linux.go:433-437`). The poststart-hook failure path explicitly notes "We're still init's parent so wait is required" and waits (`sources/runc/libcontainer/container_linux.go:346-355`).

### What does the API claim when cancellation is requested but work has not stopped?

`Container.Signal`/`Signal(SIGKILL)` returns immediately after the kill syscall succeeds; it performs no wait and publishes nothing about completion (`sources/runc/libcontainer/container_linux.go:438-490`). Completion truth is recovered separately: `hasInit()` re-checks liveness with start-time comparison on every state query (`sources/runc/libcontainer/container_linux.go:939-952`), and `Destroy()` refuses running containers with `ErrRunning` / paused with `ErrPaused` (`sources/runc/libcontainer/state_linux.go:134-139,186-194`) — i.e., the API claims *nothing* has stopped until verified. The CLI layers a bounded wait: `runc delete --force` polls `Signal(0)` for up to ~10 s and reports `"container init still running"` otherwise (`sources/runc/delete.go:17-26`). The `run/create` foreground path is truthful in the other direction: it only publishes the exit status after `process.Wait()` has flushed Go-side pipe handling (`sources/runc/signals.go:81-85`).

### Can cleanup block forever or overwrite the primary failure?

Blocking: `tty.Close()`'s `wg.Wait()` can block until container stdio pipes reach EOF (stdout/stderr copiers exit only when writers close, `sources/runc/tty.go:26-30,180`) — this is intentional drain-after-death but is unbounded in time. The exec-fifo fast path polls with `Poll(-1)` (no timeout) but is woken by pidfd exit (`sources/runc/libcontainer/container_linux.go:284-303`); the pre-pidfd fallback relies on 100 ms liveness polling (`sources/runc/libcontainer/container_linux.go:306-327`). `sdNotifyBarrier` is capped at 30 s and treats timeout as non-fatal (`sources/runc/notify_socket.go:227-241`). Overwriting: protected — init fatal logs are appended to, not substituted for, the primary error (`sources/runc/libcontainer/container_linux.go:410-419` uses `%w; %w` joining); cleanup failures during `initProcess.start` teardown become warnings while retErr survives (`sources/runc/libcontainer/process_linux.go:811-819`); conversely, in `destroy()` a cgroup-destroy failure *aborts* the remaining steps (state dir left behind, hooks skipped — `sources/runc/libcontainer/state_linux.go:52-62`), while a poststop-hook failure is returned but the container still transitions to stopped (`sources/runc/libcontainer/state_linux.go:64-66`).

### How is shutdown made idempotent under concurrent callers?

All mutating container entry points take `c.m` (`Destroy` at `sources/runc/libcontainer/container_linux.go:795-797`, `Signal` at 438-441). Idempotency comes from the state machine plus live reconciliation: every decision re-derives reality via `refreshState()` (freezer state, liveness, exec-fifo presence — `sources/runc/libcontainer/container_linux.go:919-936`) and `loadedState.destroy` re-dispatches after refresh (`sources/runc/libcontainer/state_linux.go:240-245`). A second `Destroy()` on a stopped container re-runs the full `destroy()` (`sources/runc/libcontainer/state_linux.go:104-106`): `RemoveAll` tolerates a missing dir (`sources/runc/libcontainer/state_linux.go:60-62`), matching the doc claim "No error is returned if the container is already destroyed" (`sources/runc/libcontainer/container_linux.go:790-791`) — but it will also **re-run poststop hooks**, which is a real idempotency gap (hooks are not marked as executed). Cross-process races (e.g., concurrent `runc delete` from another invocation) are arbitrated by filesystem state: atomic rename of state.json (`sources/runc/libcontainer/container_linux.go:882-906`) and state persisted before the `procRun` sync precisely so a killed creator leaves a deletable record (`sources/runc/libcontainer/process_linux.go:997-1010`).

## Architectural Decisions

1. **Signals instead of contexts as the cancellation backbone.** The supervisor registers SIGCHLD eagerly and defers registering all other signals to a background goroutine for startup-latency reasons (issue #5208 referenced) (`sources/runc/signals.go:23-36`); the forward loop is the single place where "cancel" (forwarded signal) meets "done" (reaped exit) (`sources/runc/signals.go:65-100`).
2. **Per-invocation process model as leak containment.** Because each command exits after publishing status via `os.Exit(status)` (`sources/runc/run.go:81-87`, `sources/runc/create.go:67-73`, `sources/runc/restore.go:120-127`), detached goroutines are implicitly reclaimed by the OS. This is why untracked loops (epoller, tickers) are tolerated.
3. **State machine as the idempotency mechanism.** `containerState.transition/destroy/status` with per-state guards (`sources/runc/libcontainer/state_linux.go:33-245`) turns repeated/concurrent shutdown calls into validated transitions rather than ad-hoc flags.
4. **Kernel-assisted termination.** Termination prefers kernel primitives over cooperative draining: `cgroup.kill` write, freeze-enumerate-kill-thaw fallback (`sources/runc/libcontainer/init_linux.go:684-720`), pidfd to make fifo waits race-free (`sources/runc/libcontainer/container_linux.go:271-280`).
5. **Two-phase IO ownership.** fds handed to the container are `postStart` closers closed right after spawn; runc-owned fds are `closers` closed only after IO drains (`sources/runc/tty.go:46-55,160-187`).
6. **Locked-OS-thread goroutines as namespace tools.** Goroutines that must setns pin their thread and intentionally omit `UnlockOSThread` so the runtime destroys the thread afterward (`sources/runc/libcontainer/process_linux.go:696-700,232-235`).

## Notable Patterns

- **Ordered LIFO defer chain in `runner.run`**: destroy-container → terminate-on-error → tty.Close → pidfd conn close, each conditional on failure mode (`sources/runc/utils_linux.go:225-231,278,280-286,301-306`).
- **Buffered result channels sized 1** so producers never block if consumers gave up: `consoleC` (`sources/runc/utils_linux.go:116`), `ForwardLogs` done (`sources/runc/libcontainer/logs/logs.go:18`), BPF errChan (`sources/runc/libcontainer/seccomp/patchbpf/enosys_linux.go:157`).
- **Deadline-with-join pattern** for external commands: timer fires → Kill → still receive from the waiter channel (`sources/runc/libcontainer/configs/config.go:626-632`).
- **Anti-hang fd closing**: closing the CRIU transport server side immediately after spawn guarantees EOF-driven failure if CRIU dies (`sources/runc/libcontainer/criu_linux.go:938-939`).
- **Self-terminating watchers**: memory-event goroutines detect cgroup destruction (control path gone / populated==0), close their fds, and close the channel to signal consumers (`sources/runc/libcontainer/notify_linux.go:52-56`; `sources/runc/libcontainer/notify_v2_linux.go:69-73`); the events main loop treats channel close as end-of-stream and sets the receiver to nil (`sources/runc/events.go:93-101`).
- **Failure forensics before teardown**: reading OOMKillCount *before* killing init because systemd may remove the cgroup otherwise, then wrapping retErr with an OOM hint (`sources/runc/libcontainer/process_linux.go:791-809`).

## Tradeoffs

- **Simplicity of process-scoped lifetime vs. embedder burden:** libcontainer pushes reaping/waiting responsibility onto library users (subreaper doc, `sources/runc/libcontainer/container_linux.go:433-437`); fine for Docker/containerd-style supervisors, but the library itself provides no Join()/Shutdown() handle aggregating its goroutines.
- **Signal-buffer sizing (2048) vs. unbounded registration cost:** large chan avoids missed signals but the two-stage Notify exists purely to keep expensive `signal.Notify(all)` off the startup path (`sources/runc/signals.go:21-29`).
- **Drain-then-close vs. bounded shutdown:** `wg.Wait()` before closing runc-owned fds guarantees no lost output but means total runtime is bounded by slowest pipe consumer, not by any policy (`sources/runc/tty.go:177-183`).
- **Polling fallback vs. kernel requirement:** pre-5.3 kernels get a 100 ms polling loop for exec-fifo readiness — portable but adds latency and wakeups (`sources/runc/libcontainer/container_linux.go:306-327`).
- **Hard-exit on SIGINT for attached ttys:** restoring terminal state instantly at the cost of skipping container destroy/IO drain (`sources/runc/tty.go:145-151`).

## Failure Modes / Edge Cases

- **Send-side goroutine leaks (process-lifetime-bounded):**
  - notify-socket reader blocked forever on `fileChan <- line` if pid1 dies before READY (`sources/runc/notify_socket.go:130-147`).
  - mount-source worker blocked on unbuffered `responseCh <- response` if the requestFn caller times out on the paired receive (`sources/runc/libcontainer/process_linux.go:738` vs guarded receive at `767-768`) — the ctx guard covers send and receive in requestFn (`759-772`) but not the worker's reply.
  - events ticker goroutine parks on `stats <- s` after the main loop exits (`sources/runc/events.go:77-86`).
  - memory-event watchers park on `ch <- struct{}{}` if a consumer stops receiving before cgroup destruction (`sources/runc/libcontainer/notify_linux.go:57`; `sources/runc/libcontainer/notify_v2_linux.go:67`).
- **Epoller never closed:** `tty.Close` deregisters the console fd via `Shutdown(CloseConsole)` but nothing ever calls `Epoller.Close()`, so the `epoller.Wait()` goroutine remains blocked in `EpollWait(-1)` until exit (`sources/runc/tty.go:128,177-179`; `sources/runc/vendor/github.com/containerd/console/console_linux.go:109-137,155-161`).
- **`os.Exit` bypasses cleanup:** attached-tty SIGINT resets the console then exits without destroying the container or joining anything (`sources/runc/tty.go:145-151`) — callers relying on `runc run` defers will not observe them.
- **Partial-startup leftovers acknowledged in-code:** if the creator is killed between `procRun` sync and state visibility, stage-2 `runc init` could leak; mitigated by persisting state before the sync (`sources/runc/libcontainer/process_linux.go:997-1010`) and by `delete --force`'s kill-all-then-destroy path including stopped containers with shared PID ns (`sources/runc/delete.go:72-88`; `sources/runc/libcontainer/state_linux.go:44-51`).
- **Restore wait gap:** `restoredProcess.wait` waits on the *criu* command, not the restored init — an admitted TODO (`sources/runc/libcontainer/restored_process.go:47-49`); `nonChildProcess` cannot be waited or terminated at all (`sources/runc/libcontainer/restored_process.go:97-103`).
- **Rootless shared-pidns hole:** without cgroup access nor private pidns, killing all container processes is impossible; runc warns strongly and falls back to signaling init only (`sources/runc/libcontainer/process_linux.go:826-839`; `sources/runc/libcontainer/container_linux.go:454-459`).
- **destroy abort leaves partial state:** cgroup destroy failure skips state-dir removal and poststop hooks; retry works because stoppedState re-runs everything, but hooks then fire twice across attempts (`sources/runc/libcontainer/state_linux.go:52-66,104-106`).

## Future Considerations

- Introduce a per-run cancellation scope (one `context.Context` threaded from the CLI action into libcontainer) so long-running operations (exec-fifo wait, criuSwrk message loop, tty drain) can be cancelled cooperatively rather than by signal side-channels; the `goCreateMountSources` design (`sources/runc/libcontainer/process_linux.go:682-776`) is the in-house template to generalize.
- Guard every producer send with the consumer-lifetime select (or size channels to bound) to convert the four send-leak sites above into benign drops.
- Add `t.epoller.Close()` to the `tty.Close` sequence (after `Shutdown`, before/after `wg.Wait`) so the epoller goroutine actually terminates on kernels where `EpollWait` would otherwise block indefinitely.
- Make poststop-hook execution recorded-in-state so repeated `Destroy()` attempts do not double-fire hooks (`sources/runc/libcontainer/state_linux.go:64-66,104-106`).
- Adopt goleak (or an equivalent goroutine census) in unit tests alongside the existing fd-leak harness (`sources/runc/libcontainer/integration/exec_test.go:1693-1778`); today goroutine leaks are invisible to CI.
- Replace `time.Tick` with `time.NewTicker` + Stop and track the events goroutine with the existing `WaitGroup` (`sources/runc/events.go:77-86`).

## Questions / Gaps

- **No goroutine-leak test infrastructure found.** Searches performed: `grep -rln goleak` across the repo (only vendor hits absent; zero matches in project code, `go.mod`, `go.sum`), `grep -rn "func Test.*(Ll)eak"` (only `TestFdLeaks*`). Conclusion: goroutine lifetime is validated only indirectly (behavioral kill/fifo/log-stop tests listed above).
- **Detach-mode signal-handler behavior:** when `--detach` is set, `newSignalHandler` is never called, so no SIGCHLD reaping loop exists in the detaching parent; whether any reparented grandchildren need reaping depends on the caller's subreaper configuration — no code in-tree handles this beyond `SetSubreaper` (`sources/runc/utils_linux.go:264-273`); evidence of intended contract is documentation-level only.
- **Whether the unguarded `responseCh` send can practically deadlock:** it requires the 1-minute ctx to expire mid-request while init keeps issuing mount requests; the worker would then be stuck at `sources/runc/libcontainer/process_linux.go:738` until process exit. No test exercises ctx expiry of the mount worker (searched `_test.go` files for `goCreateMountSources` — no direct unit test found).
- **Cross-process concurrent `delete --force` vs. `run` teardown:** mutual exclusion across separate runc invocations relies on filesystem/state reconciliation; no lock file or flock was found around `Destroy` within this source tree (search boundary: `sources/runc/libcontainer/*.go`, top-level `*.go`).
- The dimension brief references Aren requirement docs (`../../../../Aren/docs/phase-1-prd/02-lifecycle-requirements.md`, phased roadmap); these lie outside the selected source directory and were therefore not opened per source-isolation rules — mapping to Aren phases is limited to what the injected prompt states.

---

Generated by `01.02-cancellation-goroutine-ownership-and-cleanup` against `runc`.
