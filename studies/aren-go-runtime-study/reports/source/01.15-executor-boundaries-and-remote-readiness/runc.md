# Source Analysis: runc

## 01.15 Executor Boundaries and Remote Readiness

### Source Info

| Field | Value |
|-------|-------|
| Name | runc |
| Path | `studies/aren-go-runtime-study/sources/runc` |
| Language / Stack | Go (low-level Linux container runtime, OCI-compliant) |
| Analyzed | 2026-08-30 |

## Summary

runc is a deliberately local-first container runtime whose entire model assumes the supervising process (the CLI) and the executed init process share a host, a PID tree, a kernel, and (optionally) filesystem namespaces. The "executor boundary" in runc is the Go interface between a libcontainer `Container` (`sources/runc/libcontainer/container.go:14-58`, concrete impl `sources/runc/libcontainer/container_linux.go:33-46`) and a `parentProcess` (`sources/runc/libcontainer/process_linux.go:38-57`). All non-local work is implemented at the kernel level (namespaces, cgroups, mount namespaces, pidfd, mount-source threads) or via CRIU's own checkpoint/restore machinery (which itself is cross-process and even cross-host for the lazy-page-server case). The CLI binds a single container to a single host root (`sources/runc/utils_linux.go:32-39`, `sources/runc/libcontainer/factory_linux.go:35-103`); nothing about the architecture is designed around submission to a remote worker. The closest distributed mechanism runc actually exposes is CRIU's `--page-server` (`sources/runc/checkpoint.go:42-44`, `sources/runc/checkpoint.go:146-160`, `sources/runc/libcontainer/criu_linux.go:394-400`) for cross-host checkpoint streaming, which still presupposes a trusted operator who already has both ends reachable. Seccomp listener sockets (`sources/runc/libcontainer/configs/config.go:46-47`, `sources/runc/libcontainer/process_linux.go:533-574`) cross a Unix-domain boundary to an external agent, but they are an out-of-band signaling channel, not an executor.

Because of this, almost every "remote" question the dimension asks about is answered by runc with a structural refusal: leases, heartbeats, fencing, duplicate-dispatch protection, and lost-worker reconciliation simply do not exist because runc is not a worker manager. The narrow executor boundary it does have is therefore a useful negative example for Aren: it shows what a deliberately non-remote, in-process executor looks like, and where the seams are that a future remote implementation would have to add.

## Rating

**4 / 10** for the dimension's own scoring rubric (remote/distributed readiness).

Rationale: runc gets full marks for a clean, type-specific executor boundary that preserves local success, cancellation, signal, cleanup, and identity semantics (`sources/runc/libcontainer/container_linux.go:33-46`, `sources/runc/libcontainer/process_linux.go:38-57`, `sources/runc/libcontainer/process.go:14-153`). It deliberately does not pretend to be remote. However, every dimension question that probes leases, fencing, duplicate dispatch, stale completion, or remote acknowledgement reduces to "the host IS the worker; the kernel enforces identity via `InitProcessStartTime`" (`sources/runc/libcontainer/container_linux.go:948-951`). There is no abstraction layer for "the worker might not be there", so a 4 (rather than 1 or 2) reflects the strong local contract and the deliberate scope choice, while honestly marking the absence of distributed mechanics.

## Evidence Collected

Every entry includes a file path with line numbers. Format: `path/to/file.ts:NN`.

| Area | Evidence | File:Line |
|------|----------|-----------|
| Container type (the executor object the CLI drives) | `Container` struct with single-host fields: state dir, cgroup manager, intelrdt, init process, sync mutex, created timestamp, exec fifo | `sources/runc/libcontainer/container_linux.go:33-46` |
| Platform-neutral Container base | `Status` enum (`Created`, `Running`, `Paused`, `Stopped`); `BaseState` with `ID`, `InitProcessPid`, `InitProcessStartTime`, `Created`, `Config` | `sources/runc/libcontainer/container.go:14-58` |
| Container lifecycle API (CLI-facing surface) | `Start`, `Run`, `Exec`, `Signal`, `Pause`, `Resume`, `Destroy`, `Set`, `Status`, `State`, `OCIState`, `Processes`, `Stats`, `NotifyOOM`, `NotifyMemoryPressure` | `sources/runc/libcontainer/container_linux.go:93-228,438-490,806-869` |
| Factory / Load (constructor vs reattach) | `Create(root,id,config)` validates id, mkdirs state, makes cgroup manager; `Load(root,id)` reads `state.json` and binds a `nonChildProcess` to a process whose parent is *not* the current runc | `sources/runc/libcontainer/factory_linux.go:35-148` |
| Process interface (inner executor boundary) | `parentProcess` Go interface: `pid`, `start`, `terminate`, `wait`, `startTime`, `signal`, `externalDescriptors`, `setExternalDescriptors`, `forwardChildLogs` | `sources/runc/libcontainer/process_linux.go:38-57` |
| Process user-facing wrapper | `Process` struct with `Wait`/`Pid`/`Signal`/`closeClonedExes` and inner `processOperations`; `IO` type for stdio endpoints | `sources/runc/libcontainer/process.go:14-169` |
| Process implementations | `containerProcess` (fork+exec child), `setnsProcess` (join existing namespaces), `initProcess` (CLONE_NEWNS/PID/...), `restoredProcess`, `nonChildProcess` (Load path) | `sources/runc/libcontainer/process_linux.go:118-723`, `sources/runc/libcontainer/restored_process.go:11-127` |
| Acknowledgement between parent and init | `syncSocket` (SOCK_SEQPACKET), `processComm` with `initSock` (bootstrap config and final PID JSON), `syncSock` (control events), `logPipe` (log forwarding) | `sources/runc/libcontainer/process_linux.go:59-116`, `sources/runc/libcontainer/sync_unix.go:18-85` |
| Sync protocol (mount requests, seccomp, hooks, ready, run) | `syncType` constants: `procError`, `procReady`, `procRun`, `procHooks`, `procHooksDone`, `procMountPlease`, `procMountFd`, `procSeccomp`, `procSeccompDone`; pair-based request/response with optional FD via SCM_RIGHTS | `sources/runc/libcontainer/sync.go:41-203` |
| Ready / run handshake with init | After bootstrap and config delivery, parent calls `setupRlimits` then `writeSync(pipe, procRun)`; init waits for `procRun` before exec'ing the user process | `sources/runc/libcontainer/process_linux.go:512-525`, `sources/runc/libcontainer/init_linux.go:422-431` |
| Identity across parent and child | `pid` struct over `initSock`: `Pid` (final child), `PidFirstChild` (intermediate); `startTime` from `/proc/PID/stat` used to detect PID reuse | `sources/runc/libcontainer/init_linux.go:40-43`, `sources/runc/libcontainer/process_linux.go:133-136`, `sources/runc/libcontainer/container_linux.go:939-952` |
| Cancellation contract | `Container.Signal` delegates to init; SIGKILL on a private-PID container triggers `signalAllProcesses` over cgroup, with explicit rootless fallback to `signalInit` | `sources/runc/libcontainer/container_linux.go:438-490` |
| Setup & "execution accepted" semantics | `handleFifo` waits on the exec fifo with `pidfd_open` (kernel 5.3+) or polling fallback; distinguishes "init wrote before dying" vs "died before writing" via EOF on read | `sources/runc/libcontainer/container_linux.go:247-339` |
| Streaming/IO contract | `tty` struct + `setupIO` (PTY vs stdio), `inheritStdio`, `setupProcessPipes`; pty master FD passed over `_LIBCONTAINER_CONSOLE` extra-fd slot | `sources/runc/tty.go:16-194`, `sources/runc/utils_linux.go:100-159`, `sources/runc/libcontainer/process_linux.go:592-620` |
| Artifact transfer (mounts) | `goCreateMountSources` goroutine that detaches via `CLONE_FS`, joins container mntns, services `procMountPlease` requests, returns open_tree-style fds via `procMountFd` | `sources/runc/libcontainer/process_linux.go:682-776` |
| Artifact transfer (seccomp agent) | `procSeccomp` sync message carries child FD; parent uses `pidfd_getfd` to clone the seccomp FD into its own process, then `cmsg.SendRawFd` over a Unix socket to a listener | `sources/runc/libcontainer/process_linux.go:532-574,940-982`, `sources/runc/libcontainer/process_linux.go:1106-1142` |
| Notify socket (sd_notify relay) | `notifySocket` binds `unixgram`, bind-mounts into container at `/run/notify`, forwards READY/MAINPID and performs `sd_notify_barrier` with a 30s deadline | `sources/runc/notify_socket.go:22-246` |
| Cleanup contract | `destroy()` tears down cgroup + intelrdt + state dir, then runs `Poststop` hooks; paused containers must Resume before Destroy; `restoredState` checks for `checkpoint` dir | `sources/runc/libcontainer/state_linux.go:39-67,134-222` |
| Lost-process detection | `hasInit()` re-`stat`s the PID and compares `StartTime` against `initProcessStartTime`; treats `Zombie`/`Dead` as "no init" | `sources/runc/libcontainer/container_linux.go:939-952` |
| State refresh on reattach | `refreshState` reconciles live kernel state (cgroup freezer, exec fifo presence, hasInit) against the in-memory `containerState` (stopped/created/running/paused/restored/loaded) | `sources/runc/libcontainer/container_linux.go:919-936`, `sources/runc/libcontainer/state_linux.go:84-245` |
| State persistence (durable state on disk) | `state.json` written atomically via `CreateTemp` + `Rename`; `InitProcessStartTime` recorded so a reloaded `nonChildProcess` cannot be fooled by PID reuse | `sources/runc/libcontainer/container_linux.go:871-906`, `sources/runc/libcontainer/factory_linux.go:124-148` |
| CLI lifecycle wiring | `runner.run` drives: `newProcess`, `setupIO`, optional `setupPidfdSocket`, `Container.Start`/`Run`/`Restore`, console wait, pid-file, notify forward, signal handler forwarding | `sources/runc/utils_linux.go:223-327` |
| Action type (run/create/restore) | `CtAct` enum selects the `Container.*` method; `init=true` only for first run | `sources/runc/utils_linux.go:373-379,381-432` |
| CLI state commands | `state` writes `state.json` view; `list` enumerates `--root`; `events` polls stats and OOM channel | `sources/runc/state.go:13-60`, `sources/runc/list.go:45-183`, `sources/runc/events.go:21-113` |
| Kill / Delete | `kill` parses signal, calls `Container.Signal`; `delete` checks `Status` and either `Destroy` (stopped) or `killContainer` (force / created); `killContainer` polls SIGKILL via signal(0) | `sources/runc/kill.go:15-67`, `sources/runc/delete.go:17-92` |
| Subreaper + signal forwarding | `signalHandler.forward` forwards SIGWINCH, reaps via `SIGCHLD`, calls `process.Wait` on init exit, returns exit code | `sources/runc/signals.go:47-127` |
| Pause / Resume | `Pause`/`Resume` only on cgroup freezer; refuses on bad state via typed errors | `sources/runc/libcontainer/container_linux.go:806-845`, `sources/runc/libcontainer/state_linux.go:165-194` |
| Checkpoint / Restore (CRIU) | CLI builds `CriuOpts`; libcontainer passes to CRIU via go-criu RPC; `--page-server ADDRESS:PORT` for cross-host lazy page streaming | `sources/runc/checkpoint.go:118-182`, `sources/runc/libcontainer/criu_linux.go:394-400`, `sources/runc/restore.go:107-128` |
| Error vocabulary | `ErrExist`, `ErrInvalidID`, `ErrNotExist`, `ErrPaused`, `ErrRunning`, `ErrNotRunning`, `ErrNotPaused`, `ErrCgroupNotExist` — locally meaningful only | `sources/runc/libcontainer/error.go:1-14` |

## Answers to Dimension Questions

**What remains common across executor types, and what stays deliberately type-specific?**

Common to every executor flavor (init, setns, restored, nonChild) is the `parentProcess` interface (`sources/runc/libcontainer/process_linux.go:38-57`) plus the `Process` user struct (`sources/runc/libcontainer/process.go:25-126`) which exposes `Wait`, `Pid`, `Signal`, and FD management uniformly. Common also: lifecycle methods on the container (`Start`/`Run`/`Exec`/`Signal`/`Pause`/`Resume`/`Destroy`/`Set`/`Status`/`State`) at `sources/runc/libcontainer/container_linux.go:93-228,438-490,806-845`.

Type-specific:

- `initProcess.start()` does the bootstrap dance (apply cgroups, reset CPU affinity, copy bootstrap data, get child PID, set up network, send `initConfig`, parse sync, await `procReady` before sending `procRun`) — `sources/runc/libcontainer/process_linux.go:778-1063`.
- `setnsProcess.start()` instead `clone3(CLONE_INTO_CGROUP)` joins an existing container, then `execSetns()` reads the joined PID over `initSock` — `sources/runc/libcontainer/process_linux.go:459-629`.
- `restoredProcess.start()` returns `errors.New("restored process cannot be started")` — `sources/runc/libcontainer/restored_process.go:31-33`.
- `nonChildProcess` (the Load path) refuses `start`, `terminate`, `wait` outright — `sources/runc/libcontainer/restored_process.go:89-103`.
- `containerProcess.signal` uses `os.Process.Signal`; init's path also may go through `signalAllProcesses` over the cgroup if SIGKILL on a non-private-PID container — `sources/runc/libcontainer/container_linux.go:444-490`.

The pattern is: a thin uniform seam (`parentProcess`) over a thick, intentionally-non-uniform body. runc is honest that "init" and "exec into an existing container" are different beasts with different ack semantics, even though both produce a `*Process`.

**When is submitted work durably accepted, and how can a caller know?**

The caller can know submitted work is durably accepted only at one precise point per process type, and that point is below runc's CLI surface — at the `parentProcess.start()` boundary:

- `initProcess.start` blocks in `parseSync` until `procReady` arrives, then writes `procRun` after applying rlimits and stamping `created` time. Until `procReady`, the child has not yet completed namespace setup; until `procRun`, the child has not `execve`'d. Only after `parseSync` returns is the container considered "running" — `sources/runc/libcontainer/process_linux.go:778-1063` (especially `:908-1063`), `sources/runc/libcontainer/init_linux.go:422-431`.
- `setnsProcess.start` waits for `procReady` and rlimits apply, then sends `procRun`. The setns child is "accepted" only after its intermediate C process exits successfully and a PID is decoded from `initSock` — `sources/runc/libcontainer/process_linux.go:459-629`.
- Durability on disk is recorded via `Container.updateState` -> `saveState`, which atomically `Rename`s `state.json` only after `procReady` — `sources/runc/libcontainer/container_linux.go:871-906`.

There is no separate "submitted" vs "accepted" verb at the CLI. The CLI's success return on `runc start` means "the init reached `procReady` and was sent `procRun`". For non-detached work, `signalHandler.forward` additionally waits for `SIGCHLD` on init before returning the exit code (`sources/runc/signals.go:53-102`), so the caller observes a final completion. Anything weaker (a "submission accepted" ticket before init is even ready) is *not* a runc concept.

**How are stale workers prevented from publishing authoritative completion?**

Runc has no concept of a "worker" that can become stale, because the supervising process IS the parent of the init process. Three mechanisms prevent false completion anyway, and they are the closest analog runc has to "stale worker prevention":

1. PID + start-time tuple identity. `Container.initProcessStartTime` is captured at ready time (`sources/runc/libcontainer/process_linux.go:1010`) and `Container.hasInit` only returns true if `/proc/PID/stat` shows the same `StartTime` and a non-Zombie/Dead state — `sources/runc/libcontainer/container_linux.go:939-952`. This blocks PID-reuse spoofing when `Load` reattaches.
2. Cgroup as authoritative liveness ground. `Destroy` requires `cgroupManager.Destroy()` to succeed (`sources/runc/libcontainer/state_linux.go:52-54`); `signalAllProcesses` walks the cgroup for SIGKILL (`sources/runc/libcontainer/container_linux.go:452-466`); `notifyOnOOM` reads the cgroup event file (`sources/runc/libcontainer/container_linux.go:847-869`). If the cgroup is gone, the container is gone.
3. State file on disk as durable receipt. `state.json` records `InitProcessPid` and `InitProcessStartTime` (`sources/runc/libcontainer/container_linux.go:871-906`, `sources/runc/libcontainer/factory_linux.go:150-172`). A future runc process opening the same ID must match the start-time tuple to claim ownership (`sources/runc/libcontainer/container_linux.go:939-952`).

None of this is a lease, a heartbeat, or a fencing token — these are kernel-bound identity primitives. If you imagine runc as a worker manager with one worker per host, that worker is the kernel, and runc has no protocol to evict it.

**Which distributed mechanisms should Aren explicitly avoid until a measured requirement appears?**

1. **Any form of lease or heartbeat over an executor.** runc relies on cgroup-freezer, `/proc` stats, and PID start-time for liveness. Adding leases before there's a measured "executor disappeared" failure mode would be premature.
2. **Application-level fencing tokens.** runc's start-time check at `sources/runc/libcontainer/container_linux.go:948-951` is already a kernel-mediated fencing check (PID reuse defense). A higher-level fencing layer would duplicate this without buying anything until the worker really does live in a separate trust domain.
3. **Cross-host CRIU as a general remote executor.** The `--page-server ADDRESS:PORT` flag at `sources/runc/checkpoint.go:42-44` is CRIU's TCP page streaming, not a generic remote executor — runc already passes it through transparently (`sources/runc/libcontainer/criu_linux.go:394-400`). It is one-way (sender pushes pages to a peer running `criu page-server`). Aren should not generalize from this to "runc is remotely executable"; the supervision protocol is local SOCK_SEQPACKET.
4. **Generic RPC submission.** Sockets are used by runc only for sd_notify (`sources/runc/notify_socket.go:107-152`), seccomp agent FD hand-off (`sources/runc/libcontainer/process_linux.go:1119-1142`), and seccomp listener FD delivery (`sources/runc/libcontainer/process_linux.go:533-574`). None of them are control-plane submission. Aren should keep remote-execution submission off the executor boundary until it has a concrete load shape that local cgroup+namespace can't satisfy.
5. **Duplicate-dispatch protection.** runc does not have one because the CLI invocation is one-shot and the kernel enforces single-PID-1-per-container via clone-flags. Aren should only add a duplicate-dispatch guard if and when the same task ID can be delivered twice from a queue.

## Architectural Decisions

- **Local-only by design.** There is no daemon, no listening socket, no client/server split. The CLI's life IS the container's life in the supervising process (`sources/runc/main.go:88-186`). When you `runc run` (foreground), exit codes are propagated via `os.Exit(status)` (`sources/runc/run.go:77-88`); when you `runc start` detached, the supervising runc exits and a separate `runc` invocation must be issued to talk to the container again (`sources/runc/start.go:13-57`, `sources/runc/utils_linux.go:223-327`).
- **Identity = (PID, start time).** Not a name, not a UUID, not a generated token. Just the kernel's clock-tick-accurate identity (`sources/runc/libcontainer/container_linux.go:939-952`, `sources/runc/libcontainer/container_linux.go:992-1010`). Reattach reads `state.json` and trusts only this tuple.
- **Acknowledge via sequence-numbered sync messages on a single SOCK_SEQPACKET channel.** No retries, no acknowledgements-of-acknowledgements, no timeouts on the sync socket itself. Errors propagate as `procError` carrying an `initError` JSON payload (`sources/runc/libcontainer/sync.go:41-202`, `sources/runc/libcontainer/init_linux.go:140-148`).
- **State machine as explicit struct-per-state.** `stoppedState`, `runningState`, `createdState`, `pausedState`, `restoredState`, `loadedState` each enforce legal transitions (`sources/runc/libcontainer/state_linux.go:84-245`). Invalid transitions return `stateTransitionError`.
- **Filesystem-as-state.** Every container has a state dir (`sources/runc/libcontainer/factory_linux.go:91-103`), an `exec.fifo` (`sources/runc/libcontainer/container_linux.go:492-539`), and a `notify/` dir for sd_notify. State durability is literally `Rename(state-tmp, state.json)` (`sources/runc/libcontainer/container_linux.go:882-906`).
- **Conservative use of new kernel features with explicit fallbacks.** `pidfd_open` is preferred for exec fifo wait, but a polling fallback exists (`sources/runc/libcontainer/container_linux.go:271-327`). `CLONE_INTO_CGROUP` is used for `runc exec` cgroup join with a fallback to writing `cgroup.procs` (`sources/runc/libcontainer/process_linux.go:368-457`).
- **Hard separation of three execution flavors.** init (CLONE_NEW*), setns (join existing), restored (CRIU). They share `parentProcess` but each owns its own `start()` because their ack semantics genuinely differ.
- **CRIU is the only "could be remote" mechanism, and only for memory pages.** Even then runc just hands the address/port to CRIU (`sources/runc/libcontainer/criu_linux.go:394-400`); runc itself does not implement TCP page streaming.

## Notable Patterns

- **Two-socket communication with init: `initSock` (bootstrap config + final PID) and `syncSock` (control messages + FDs).** Distinct from `logPipe` for log streaming (`sources/runc/libcontainer/process_linux.go:59-116`). This separation matters because the sync socket is SOCK_SEQPACKET (preserves message boundaries) while init data is a stream (`sources/runc/libcontainer/sync_unix.go:18-85`).
- **Atomic state persistence.** Temp file + rename (`sources/runc/libcontainer/container_linux.go:882-906`) so a crash never leaves a partial `state.json`.
- **Closure-delivered cleanup.** Every goroutine that may leak a thread uses `runtime.LockOSThread` without `UnlockOSThread` to let the runtime reap the thread (`sources/runc/libcontainer/process_linux.go:695-700`, `sources/runc/libcontainer/process_linux.go:232-237`). This is unusual but well documented in the comments.
- **Deferred init/setup ordering.** Init's `start()` does cgroup-apply *before* sync handshake so a child can never escape the cgroup (`sources/runc/libcontainer/process_linux.go:823-839`).
- **Identity verification on every operation that touches the init process.** `signalInit` first calls `hasInit()` to refuse killing a non-running container (`sources/runc/libcontainer/container_linux.go:473-490`), and `Destroy` requires `hasInit()==false` from `runningState.destroy` (`sources/runc/libcontainer/state_linux.go:134-139`).
- **FD handoff via SCM_RIGHTS.** Sync messages may carry an FD (`sources/runc/libcontainer/sync.go:90-110`), and the seccomp agent handoff re-uses the same primitive (`sources/runc/libcontainer/process_linux.go:1106-1142`).

## Tradeoffs

- **No out-of-process supervision of the init process.** If the supervising `runc runc` process dies while init is running, no one is left to reap or signal. runc mitigates by recommending `--detach` + an external supervisor (`sources/runc/run.go:46-50`), but there is no built-in replacement supervisor. For Aren this means: a future "executor" type that wants this property needs to design for it up front, because runc does not.
- **Cgroup is the only authoritative liveness source.** When the cgroup vanishes (e.g., rootless without delegated cgroup), runc cannot reliably kill processes (`sources/runc/libcontainer/container_linux.go:453-458`). The design deliberately refuses to fake this.
- **CRIU page server is one-way.** CRIU's `--page-server` lets a checkpoint stream pages to a TCP peer, but runc does not implement restoration-from-remote — that's CRIU's job. Treat it as a transport, not as a remote execution API.
- **`Load` reattachment assumes the kernel still knows about the init process.** If init was reaped elsewhere, `Load` succeeds but every operation on it will fail with `ErrNotRunning` or "process not found". runc does not maintain a tombstone beyond `state.json`.
- **All container lifetime assumptions depend on shared host root and shared filesystem.** No remote executor could reuse this layer without replicating `/run/runc/<id>/state.json`, the `exec.fifo`, and the `notify/` socket. That is a known constraint, not an oversight.

## Failure Modes / Edge Cases

- **Init dies before `procReady`.** Parent's `parseSync` returns; `ierr` becomes non-nil; `containerProcess.terminate()` runs and `nonChildProcess` may never be set (`sources/runc/libcontainer/process_linux.go:778-1063`). runc `runc delete --force` can clean this up via `killContainer` (`sources/runc/delete.go:17-26`).
- **Exec fifo created but init never wrote.** `readFromExecFifo` returns the explicit error `exec fifo is empty: container init did not signal execve readiness (process died before writing, or fifo already consumed)` (`sources/runc/libcontainer/container_linux.go:329-339`).
- **Init process reaped by someone else (PID reuse).** `hasInit` returns false because `stat.StartTime != initProcessStartTime` (`sources/runc/libcontainer/container_linux.go:946-951`). All subsequent calls fail with `ErrNotRunning`.
- **Paused cgroup + SIGKILL.** SIGKILL on a paused container thaws first then kills — `signalInit` calls `cgroupManager.Freeze(Thawed)` if paused (`sources/runc/libcontainer/container_linux.go:481-488`).
- **Shared PID namespace.** `signal(SIGKILL)` falls back to `signalAllProcesses` because the kernel will not auto-kill peers in a shared PID ns; destroy also walks the cgroup (`sources/runc/libcontainer/container_linux.go:452-466`, `sources/runc/libcontainer/state_linux.go:48-51`).
- **Rootless without delegated cgroups.** `ErrRootless` is ignored only if the container has a private PID namespace; otherwise the warning at `sources/runc/libcontainer/process_linux.go:831-834` says this will become an error.
- **State.json written but parent died before `procRun`.** The init process may become orphaned; only `runc delete/stop` can clean it (and only if it can find the state dir). The comment block at `sources/runc/libcontainer/process_linux.go:996-1010` documents this.
- **Notify socket — supervisor does not support `sd_notify_barrier`.** The 30-second timeout warns and continues (`sources/runc/notify_socket.go:236-240`). runc treats this as best-effort.
- **CRIU remote page server unreachable.** runc surfaces the CRIU error verbatim; there is no retry and no fencing.

## Future Considerations

For Aren to extend runc's executor model without breaking its invariants, the seams to widen are:

1. **Identity tuple.** Promote (PID, start-time) into a richer executor handle (e.g., a `WorkerLease` with a kernel-derived nonce). This is the seam where leases and fencing would attach.
2. **`parentProcess.start()`.** Today it is a blocking in-process function. The natural place to introduce a non-blocking "submit, await ack later" is here, but only if Aren actually needs an executor that outlives the supervising runc.
3. **`Container.Status` and `refreshState`.** Currently uses cgroup freezer + exec fifo presence + PID stat. A remote executor would need to replace this with a pollable remote status source.
4. **`Container.Save`/`Load`.** Already durable on disk (`sources/runc/libcontainer/container_linux.go:871-906`); the natural extension is to persist a lease token alongside `InitProcessStartTime`.
5. **`destroy()`.** Local-only; assumes the same kernel can SIGKILL via cgroup. A remote executor needs a control channel to ask the worker to kill.
6. **CRIU `--page-server`.** Already a TCP-aware checkpoint path. Any future "remote cold start" would naturally compose here, but only as a CRIU-mediated flow, not a runc-native one.
7. **Typed errors.** The current error vocabulary (`ErrExist`, `ErrRunning`, `ErrNotRunning`, `ErrPaused`, `ErrNotPaused`, `ErrCgroupNotExist`) is host-local. A remote layer would need additional types such as `ErrWorkerLost`, `ErrStaleCompletion`, `ErrDuplicateDispatch`, and they would need to be checkable from afar.

## Questions / Gaps

- runc ships no integration tests for a remote/distributed executor (no test creates a `Container` whose `parentProcess` lives off-host). Search boundary: `tests/integration/`, `tests/cmd/`. Result: only local patterns, plus the `seccompagent` and `pidfd-kill` Unix-socket helpers (`sources/runc/tests/cmd/seccompagent/seccompagent.go`, `sources/runc/tests/cmd/pidfd-kill/pidfd-kill.go`), which exercise FD handoff but not remote submission.
- There is no concept of "submission acknowledgement" distinct from "process started". The closest evidence is the `procReady` -> `procRun` pair (`sources/runc/libcontainer/init_linux.go:422-431`), and these are local-only sync socket events.
- runc has no fencing token; identity is implicit in (PID, start-time). Any remote executor design that wishes to be Aren-compatible must decide whether to reuse the start-time approach or invent its own fencing token.
- CRIU's `--page-server` is the only built-in feature that touches a remote peer (`sources/runc/checkpoint.go:42-44`). It is unclear from runc's source whether runc has ever asserted that the page server is on a different host — the CLI only requires `ADDRESS:PORT` (`sources/runc/checkpoint.go:146-160`).
- `Container.Set` updates resources of a running container (`sources/runc/libcontainer/container_linux.go:168-201`). It assumes the cgroup is local; no path here is meaningful to a remote executor without an explicit control plane.
- No evidence found for any explicit "duplicate-dispatch" guard in the supervisor path. The CLI's idempotency is one-shot and assumes the kernel serializes.
- No evidence found for any "lost worker reconciliation" code path. The closest analog is `Load` + `refreshState`, both of which assume the worker (the kernel process) is still around.

---

Generated by `01.15-executor-boundaries-and-remote-readiness` against `runc`.