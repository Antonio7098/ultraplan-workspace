# Source Analysis: docker-agent

## 01.07 Subprocess Supervision and Process-Tree Containment

### Source Info

| Field | Value |
|-------|-------|
| Name | docker-agent |
| Path | `studies/aren-go-runtime-study/sources/docker-agent` |
| Language / Stack | Go 1.24+ / os/exec, Docker sandbox backend |
| Analyzed | 2026-05-13 |

## Summary

docker-agent implements two primary local-process runners — the synchronous `shell` tool (`pkg/tools/builtin/shell/shell.go:32`) and the asynchronous `background_jobs` tool (`pkg/tools/builtin/backgroundjobs/backgroundjobs.go:54`) — both built on `os/exec` with platform-split process-group containment, manual timeout/cancellation, and `WaitDelay` pipe-draining. A third ad-hoc runner (`script_shell`, `pkg/tools/builtin/shell/script_shell.go:240`) and hook commands (`pkg/hooks/handler.go:224`) and several auxiliary `exec.CommandContext` call sites provide no process-group containment at all. The design correctly mitigates the classic Unix pipe-deadlock for backgrounded grandchildren via `WaitDelay`, bounds output only for background jobs (10 MB), and leaves synchronous shell output unbounded. Unix uses `Setpgid:true` + `kill(-pid, SIGTERM)` with SIGKILL escalation; Windows uses `JOB_OBJECT_LIMIT_KILL_ON_JOB_CLOSE` job objects with handle-closure. No cgroup/namespace isolation is applied at the tool level; Docker sandboxing exists only as an outer execution wrapper (`pkg/sandbox/sandbox.go:28`) and is not a subprocess-containment guarantee for local tools.

## Rating

**5 / 10**

Rationale: The two main tool runners demonstrate competent OS-level containment (process groups / job objects), disciplined reaping, pipe-drain handling, and distinct status classification, with regression tests for orphaned-grandchild hangs. Points deducted for: unbounded output on the synchronous shell path (memory-exhaustion vector), absence of any containment in `script_shell`/`hooks`/`MCP stdio`/`environment` runners, inherent races in post-Start job-object assignment and SIGTERM-only escalation, no `Pdeathsig`/`Cancel`/`WaitDelay` hardening for non-shell paths, and no cgroup/namespace/resource bounding with explicit non-guarantees only partially documented.

## Evidence Collected

Every entry MUST include a file path with line numbers. Format: `path/to/file.ts:NN`.

| Area | Evidence | File:Line |
|------|----------|-----------|
| Process creation – shell (manual cancel) | `exec.Command` (not `CommandContext`) with `SysProcAttr` from `platformSpecificSysProcAttr()` and `WaitDelay=500ms`; timeout via `context.WithTimeout`, `Start`+`Wait` in goroutine, `kill` on timeout/cancel | `pkg/tools/builtin/shell/shell.go:188-233` |
| Process creation – background jobs | `exec.Command(... )` with `SysProcAttr`, `WaitDelay=1s`, goroutine `monitorJob` waiting on `cmd.Wait()` | `pkg/tools/builtin/backgroundjobs/backgroundjobs.go:235-272`, `283-326` |
| Process creation – script_shell (no supervision) | `exec.CommandContext(ctx, shell, ...)` + `cmd.Run()` — no `SysProcAttr`, no `WaitDelay`, no process-group, no output bound | `pkg/tools/builtin/shell/script_shell.go:240-293` |
| Process creation – hooks (no supervision) | `exec.CommandContext(ctx, h.shell, ...)` + `cmd.Run()` with only stdout/stderr buffers, no `SysProcAttr`, no `WaitDelay` | `pkg/hooks/handler.go:224-259` |
| Process creation – MCP stdio | `exec.CommandContext(ctx, c.command, ...)` delegated to `mcp.CommandTransport` — no `SysProcAttr` set at this layer | `pkg/tools/mcp/stdio.go:62-67` |
| Environment construction – shell | `toolsetEnv` prepends `os.Environ()` so host env inherited, toolset env wins via last-wins dedup; `askpass` env injected separately | `pkg/tools/builtin/shell/shell.go:277-287`, `484` |
| Environment construction – script_shell | Per-call env clone from `t.env` (or `os.Environ()` if nil), then `toolConfig.Env` (sorted keys, `${env.X}` expansion) then filtered `params` keys, NUL check | `pkg/tools/builtin/shell/script_shell.go:246-283` |
| Working directory | `checkWorkDir` pre-validates cwd existence/type; `resolveWorkDir` joins relative `cwd` against `h.workingDir`; `cmd.Dir=cwd` set explicitly | `pkg/tools/builtin/shell/shell.go:166-168`, `195`, `316-342` |
| Unix process group | `SysProcAttr{Setpgid:true}`, `kill` does `syscall.Kill(-pid, SIGTERM)` | `pkg/tools/builtin/shell/cmd_unix.go:14-25`, `pkg/tools/builtin/backgroundjobs/cmd_unix.go:12-21` |
| Windows job object | `CreateJobObject` + `SetInformationJobObject(JOB_OBJECT_LIMIT_KILL_ON_JOB_CLOSE)` + `OpenProcess` + `AssignProcessToJobObject`; `kill` closes handles then `proc.Kill()` | `pkg/tools/builtin/shell/cmd_windows.go:20-73`, `pkg/tools/builtin/backgroundjobs/cmd_windows.go:20-70` |
| JS/WASM stub | `platformSpecificSysProcAttr() -> nil`, `kill -> error "not supported"` | `pkg/tools/builtin/shell/cmd_js.go:15-24` |
| Signal escalation | On `timeoutCtx.Done()`: `kill(pg)` then `select` on `done` with `3s` grace, then `cmd.Process.Kill()` (SIGKILL) for shell; `2s` grace in `reapSpawnedChild` | `pkg/tools/builtin/shell/shell.go:221-233`, `244-261` |
| Startup-failure reaping | `createProcessGroup` failure → `reapSpawnedChild(cmd, pg)` (SIGTERM then 2s → SIGKILL + `Wait`) | `pkg/tools/builtin/shell/shell.go:207-213`, `244-261`; `pkg/tools/builtin/backgroundjobs/backgroundjobs.go:261-265`, `459-476` |
| Pipe draining / WaitDelay – shell | `cmd.WaitDelay = 500ms` with comment explaining OS-pipe + copy-goroutine deadlock when grandchild inherits fds (e.g. `docker run &`) | `pkg/tools/builtin/shell/shell.go:173-198` |
| Pipe draining – background jobs | `cmd.WaitDelay = 1s`; comment documents truncation/SIGPIPE tradeoff; `monitorJob` treats `exec.ErrWaitDelay` as success if `ExitCode==0` | `pkg/tools/builtin/backgroundjobs/backgroundjobs.go:41-52`, `239`, `291-301` |
| Pipe draining – eval container | `cmd.Cancel = SIGINT` (proxy through Docker CLI), `cmd.WaitDelay=10s`, `StdoutPipe`/`StderrPipe` + `bufio.Scanner` with 1M→10M buffer, `io.ReadAll(stderr)` goroutine | `pkg/evaluation/eval.go:471-502` |
| Output bounding – background jobs | `maxBackgroundJobOutputBytes=10MB`, `limitedWriter` that reports `len(p)` but writes only `remaining` bytes; truncation notice on render | `pkg/tools/builtin/backgroundjobs/backgroundjobs.go:39`, `99-118`, `377-382` |
| Output bounding – shell (unbounded) | `commandOutput` with `bytes.Buffer` + mutex, `Write` always appends all `p`; no size cap; only formatting via `strings.TrimSpace` + `cmp.Or("<no output>")` | `pkg/tools/builtin/shell/shell.go:62-91`, `345-360` |
| Output bounding – eval transcript | `response[:500]+"...(truncated)"` per tool response; scanner buffer limit 10MB; termination-field rune limit 400 | `pkg/evaluation/eval.go:616-617`, `500` |
| Exit classification – shell | `formatCommandOutput` distinguishes `timeoutCtx.Err()!=nil` → `ctx.Err()!=nil` ? "cancelled" : "timed out", else `cmdErr != nil` ? "Error executing..." | `pkg/tools/builtin/shell/shell.go:344-360` |
| Exit classification – background | `monitorJob`: `exec.ErrWaitDelay` specially handled vs `*exec.ExitError` → stored `exitCode`/`err`/newStatus; `statusStopped` via `CompareAndSwap` in `StopBackgroundJob` | `pkg/tools/builtin/backgroundjobs/backgroundjobs.go:286-412`, `425-440` |
| Exit classification – hooks | `*exec.ExitError` → structured `ExitCode`; other error → `ExitCode=-1` + bubble up | `pkg/hooks/handler.go:248-257` |
| Containment non-guarantee docs | Hook doc: spawned hooks "may keep running as orphaned process, or may be torn down... depending on your environment rather than on anything docker-agent guarantees" | `docs/guides/headless/index.md:87` (and `pkg/tools/builtin/backgroundjobs/backgroundjobs.go:43-50` note about redirection need) |
| Privilege boundary | `sudoAskpass` + `SUDO_ASKPASS` wrapper, Unix socket token auth, only added when `commandInvokesSudo` and `posixShellForFunc`; no general privilege drop | `pkg/tools/builtin/shell/askpass.go:58-63`, `406-486` |
| Tests – orphaned grandchild pipe hang | `TestShellTool_BackgroundedChildDoesNotBlockReturn` (`sleep 30 &` must return <5s) and detached `setsid sleep 30 &` variant; background-jobs equivalent | `pkg/tools/builtin/shell/shell_test.go:334-397`, `pkg/tools/builtin/backgroundjobs/backgroundjobs_test.go:427-471` |
| Tests – reaping | `TestReapSpawnedChild` verifies `ProcessState != nil` after kill+Wait with 3s deadline | `pkg/tools/builtin/shell/shell_test.go:399-451` |
| Tests – timeout | `TestBackgroundJobsTool_WaitBackgroundJob_Timeout` (`sleep 30`, 1s wait → "Timed out") and `ContextCancelled` | `pkg/tools/builtin/backgroundjobs/backgroundjobs_test.go:169-223` |
| Lint guardrail | `ConstructorCommandExec` lint forbids `os/exec.Command/CommandContext` and `Start/Run/Output` inside `New*` constructors | `lint/constructor_command_exec.go:10-54` |

## Answers to Dimension Questions

**Does cancelling the parent guarantee that descendants stop?**

No. Guarantee is best-effort and platform-conditional.

* Unix (`pkg/tools/builtin/shell/cmd_unix.go:14-25`): `Setpgid:true` puts the direct shell child in a new process group; `kill` sends `SIGTERM` to `-pid` (the whole group). Children that were `fork`ed before cancellation and share the pgid receive SIGTERM. However: (a) only SIGTERM is sent initially, so a process ignoring SIGTERM survives until the 3s (shell, `shell.go:229`) / 2s (reap) grace expires and `SIGKILL` is sent — but only to the *direct* `cmd.Process`, not the group again. (b) A `setsid`-detached grandchild is in a different session/pgid and is *not* reached by `kill(-pid)`; the regression test `TestShellTool_DetachedBackgroundedChildDoesNotBlockReturn` (`pkg/tools/builtin/shell/shell_test.go:362-397`) explicitly documents that "the process-group kill fallback ... cannot reach it". In that case cancellation only closes pipes via `WaitDelay` and lets the grandchild leak. (c) `script_shell` and `hooks` have no pgid at all, so their descendants are never signalled.

* Windows (`pkg/tools/builtin/shell/cmd_windows.go:20-73`): A `JobObject` with `JOB_OBJECT_LIMIT_KILL_ON_JOB_CLOSE` is created *after* `Start`; if the child spawns grandchildren before `AssignProcessToJobObject` succeeds they may escape the job. The kill path closes job+process handles to trigger kernel kill, then `proc.Kill()` as fallback. Nested jobs (child already in a job) cause `AssignProcessToJobObject` to fail and `reapSpawnedChild` is invoked.

* `background_jobs.Stop` (`backgroundjobs.go:425-440`) and `ToolSet.Stop` (`backgroundjobs.go:619-627`) iterate all `statusRunning` jobs and `kill` each group/job, but again SIGTERM-only, no second-wave group SIGKILL.

* Hooks run under `exec.CommandContext(ctx, ...)` (`pkg/hooks/handler.go:224-225`) so parent context cancellation kills only the direct shell via Go's internal signal, not descendants (no `Setpgid`).

Explicit non-guarantee: hook docs state orphaned hook processes "may keep running as an orphaned process, or may be torn down by whatever supervises the job ... depending on your environment rather than on anything docker-agent guarantees" (`docs/guides/headless/index.md:87`).

**Can unread output deadlock process completion or exhaust memory?**

Both vectors exist but are partially mitigated.

*Deadlock:* Go's `exec.Cmd` with `Stdout/Stderr` not an `*os.File` creates OS pipes with copy goroutines; `Wait()` blocks until pipes see EOF. A backgrounded grandchild inheriting the pipe (e.g. `sleep 10 &`, `docker run &`) previously caused tool hangs until the configured 30s timeout (`pkg/tools/builtin/shell/shell.go:173-186` comment). Fix is `WaitDelay` (shell 500 ms, background 1 s, eval 10 s) which force-closes pipes and makes `Wait()` return `exec.ErrWaitDelay`. Tests for this are `TestShellTool_BackgroundedChildDoesNotBlockReturn` and `TestBackgroundJobsTool_BackgroundedChildDoesNotBlockReturn`. `WaitDelay` trades correctness for liveness: still-running grandchild output is truncated and next write gets SIGPIPE (documented at `pkg/tools/builtin/backgroundjobs/backgroundjobs.go:41-50`). Pipelines not using `WaitDelay` (`script_shell.go:284-288`, `hooks/handler.go:242-243`) retain the deadlock risk — mitigated only by using in-memory `bytes.Buffer` writers rather than pipes, but `script_shell` uses `commandOutput` which still relies on pipe plumbing under `exec`.

*Exhaustion:* `shell` tool's `commandOutput.buf` (`pkg/tools/builtin/shell/shell.go:62-91`) is an unbounded `bytes.Buffer` — every `Write(p)` appends `p[:]` with no cap. A `yes | head` style or large `cat` can grow buffer to O(output size) and OOM the agent. `background_jobs` is the exception: `limitedWriter` (`pkg/tools/builtin/backgroundjobs/backgroundjobs.go:99-118`) caps at `maxBackgroundJobOutputBytes = 10 MB` (`:39`), reporting `len(p)` to the child but dropping excess, with `[Output truncated at 10MB limit]` notice. `script_shell` and `hooks` are unbounded (`bytes.Buffer` at `script_shell.go:284`, `hooks/handler.go:242`). Stdio MCP (`stdio.go:62`) streams with no explicit byte bound (handled by SDK). Eval scanner is capped at 10 MB line buffer (`eval.go:500`).

**Which containment claims hold on each supported operating system?**

| OS | Mechanism | Claim & evidence | Limitation |
|----|-----------|------------------|------------|
| Linux / macOS (`!windows && !js`) | Process groups (`Setpgid:true` + `kill(-pid,SIGTERM)`) | Descendants sharing the group receive SIGTERM on timeout/cancel/stop (`pkg/tools/builtin/shell/cmd_unix.go:12-25`) | `setsid`/daemonized children escape; SIGTERM-only initial signal; no `Pdeathsig`, no cgroup/namespace, no fd closure beyond pipes |
| Windows | Job Objects (`JOB_OBJECT_LIMIT_KILL_ON_JOB_CLOSE`) | All processes assigned to job are killed when last job handle closed (`pkg/tools/builtin/shell/cmd_windows.go:26-35`) | Race between `Start` and `AssignProcessToJobObject`; nested-job children fail assignment; requires `OpenProcess` handle; no effect on already-detached processes |
| js/wasm | No-op | `platformSpecificSysProcAttr() -> nil`, `kill -> error` (`cmd_js.go:15-24`, `backgroundjobs/cmd_js.go:13`) | No process supervision exists in browser; claim is explicit non-support |
| All (outer) | Docker sandbox | `pkg/sandbox/sandbox.go:28-170` creates/ensures a Docker sandbox (`docker sandbox create`) and `BuildExecCmd` mounts workspaces; used only when `--sandbox` is active | Not a per-tool subprocess boundary; `shell`/`background_jobs` still run as host processes unless the whole agent is inside the sandbox container; no cgroup/namespace claims for individual `exec` calls |

No evidence found for explicit cgroup, namespace, `unshare`, `chroot`, or capability-dropping per tool invocation; resource/containment boundaries are absent and not documented as non-guarantees except for the hook note above.

**How are exit status, timeout, cancellation, setup failure, and cleanup failure distinguished?**

* `shell.RunShell` (`pkg/tools/builtin/shell/shell.go:136-238`, `344-360`): Two-layer `select` on `timeoutCtx.Done()` vs `done<-cmd.Wait()`. `formatCommandOutput` checks `timeoutCtx.Err()!=nil` first; if also `ctx.Err()!=nil` → `"Command cancelled"` (user/parent cancellation), else → `"Command timed out after %v\nOutput: ..."` (deadline exceeded). If context not cancelled, `cmdErr != nil` → `"Error executing command: %s\nOutput: %s"` (non-zero exit or signal). Success → raw output trimmed, with `"<no output>"` fallback. Exit code itself is embedded only via `err.Error()` string, not a structured field. The synchronous timeout param (`RunShellArgs.Timeout` int seconds, default 30s, `shell.go:141-147`) and per-call `cwd` are stamped on OTel span.

* `background_jobs`: `monitorJob` (`backgroundjobs.go:283-326`) classifies `cmd.Wait()` error: if `errors.Is(err, exec.ErrWaitDelay)` → treat as success using `ProcessState.ExitCode()` (grandchild pipe case); else if `*exec.ExitError` → `exitCode = ExitCode()`, `statusFailed`; else `exitCode=-1`, `statusFailed`. Success with `exitCode==0` → `statusCompleted` else `statusFailed`. `StopBackgroundJob` atomically swaps `statusRunning→statusStopped`; subsequent `monitorJob` sees CAS failure and returns without overwriting status, so stopped is distinguishable from failed/completed. `WaitBackgroundJob` (`:397-423`) distinguishes completion vs timeout (timer channel) vs context cancellation, each with distinct header strings. View/render includes explicit `Status: running/completed/stopped/failed` strings.

* `script_shell` (`script_shell.go:288-293`): Minimal — `cmd.Run()` error → `ResultError("Error executing command '%s': %s\nOutput: %s")`, else `ResultSuccess(output)`. No timeout/cancel distinction (relies on `CommandContext(ctx,...)` propagation, which returns `context.Canceled/DeadlineExceeded` as `err`).

* Hooks (`hooks/handler.go:246-257`): `*exec.ExitError` → `HandlerResult{ExitCode: exitErr.ExitCode()}` with `err=nil` (structured); other error (e.g. binary missing, start failure) → `ExitCode=-1`, returned error bubbles up to fail-closed. Distinguishes run failure from spawn failure.

* Setup failures (e.g. `createProcessGroup` after `Start`): Both `shell.go:207-213` and `backgroundjobs.go:261-265` treat as setup error — `reapSpawnedChild` is called to kill+wait the already-started child (SIGTERM → 2s → SIGKILL), then return `ResultError("Error creating process group: ...")`. Child is not leaked. `checkWorkDir` preflight mirrors this for cwd (`shell.go:166-168`).

* Cleanup failures: `kill` errors are logged/returned as strings where visible: `shell` timeout path ignores `_ = kill(...)` (`shell.go:223`); `StopBackgroundJob` returns `"Job %s marked as stopped, but error killing process: %s"` on `kill` error (`backgroundjobs.go:437`); Windows `kill` closes handles first then `proc.Kill()` and returns its error. `reapSpawnedChild` escalation failures are ignored (`_ = kill`, `_ = cmd.Wait()`).

## Architectural Decisions

| Decision | Rationale (inferred from code/comments) | Consequence |
|----------|----------------------------------------|-------------|
| Manual cancellation with `exec.Command` + explicit `kill` instead of `exec.CommandContext` for `shell` (`pkg/tools/builtin/shell/shell.go:188-191` comment) | Keep cancellation, timeout, process-group kill, and `WaitDelay` in one place with escalation logic | Allows SIGTERM→SIGKILL staging and group kill; diverges from simpler `CommandContext` used elsewhere, creating inconsistency |
| `WaitDelay` (500 ms shell, 1 s background, 10 s eval) to force-close pipes | Backgrounded `&` grandchildren inherit pipe fds and block `Wait()` for the full timeout (`pkg/tools/builtin/shell/shell.go:173-186`) | Prevents HANG regression (tested); truncates/drops grandchild output and SIGPIPE-kills it on next write — documented tradeoff |
| Unix `Setpgid:true` + `Kill(-pid)` vs Windows `JobObject` | Portable tree-kill without `cgroups`/platform-specific APIs | Unix cannot reach `setsid`-detached trees; Windows has post-Start assignment race; js no containment |
| Bounded output only for `background_jobs` (10 MB `limitedWriter`) | Background jobs may run long (servers) and OOM agent if unbounded; shell calls assumed short | Leaves `shell`, `script_shell`, `hooks`, `MCP` unbounded — inconsistent exhaustion boundary |
| Preflight `checkWorkDir` before `Start` (`pkg/tools/builtin/shell/shell.go:316-331`) | Avoid misleading `fork/exec <shell>: no such file` when `SysProcAttr` forces `chdir` failure to be misattributed to shell binary | Improves error UX but TOCTOU remains (dir may vanish between check and `Start`) |
| `toolsetEnv` always `append(os.Environ(), env...)` with last-wins dedup | Spawns inherit host environment while toolset overrides win; avoids stripping PATH/etc. (`pkg/tools/builtin/shell/shell.go:281-286`, `script_shell.go:244-249`) | Implicit full-env inheritance complicates hermeticity; no allow-list filtering |
| Lint rule `ConstructorCommandExec` | Prevent hiding `exec` side effects in `New*` constructors | Enforces explicit operational boundaries for process lifecycle |

## Notable Patterns

* **Platform-split supervision files:** `cmd_unix.go` / `cmd_windows.go` / `cmd_js.go` triplets for both `shell` and `background_jobs`, selected by build tags (`//go:build !windows && !js`). Clean compile-time abstraction but duplication (≈ identical job-object code in two packages).
* **Graceful escalation ladder:** `SIGTERM` (group/job) → timed wait (2–3 s) → `SIGKILL`/`Kill()` → `Wait()`. Seen in `shell.go:221-233`, `reapSpawnedChild:248-260`, `backgroundjobs.go:459-476`, and eval's `Cancel: SIGINT` + `WaitDelay:10s` (`eval.go:476-479`).
* **`limitedWriter` as bounded pipe:** Mutex-shared between writer and readers, always returns `len(p)` to child while dropping excess (`backgroundjobs.go:108-118`). Prevents child blocking on full buffer.
* **`commandOutput` live streaming:** `Write` locks, appends to `buf`, then `emit(Runtime.EmitOutput)` with captured `ctx` (`shell.go:70-85`). Streams output to runtime while buffering for final result — but without bound.
* **Deterministic shell resolution:** `shellpath.DetectShell` prefers `SHELL`/`/bin/sh` on Unix, `pwsh`/`powershell` via `LookPath` then `ComSpec`/`SystemRoot\System32\cmd.exe` on Windows (`shellpath.go:40-81`), mitigating PATH hijacking (CWE-426).
* **Setup-failure reaping pattern:** `createProcessGroup` error → immediate `reapSpawnedChild` ensures no zombie or pipe leak even on partial startup failure.
* **Eval's cancellation proxy:** `cmd.Cancel = Signal(os.Interrupt)` instead of default `SIGKILL` because Docker CLI only proxies `SIGINT` to container (`eval.go:471-478`).

## Tradeoffs

* **WaitDelay vs. completeness:** Low delays (500 ms–1 s) guarantee responsiveness when shell backgrounds a server, but sacrifice grandchild output and may SIGPIPE-kill it. Higher delay (eval 10 s) preserves logs longer. No knob per call.
* **SIGTERM vs. SIGKILL:** Initial SIGTERM is graceful (allows cleanup handlers) but ignorer processes stall up to 3 s; escalation kills only the direct pid on Unix, not the group again, so ignoring members of the group survive.
* **Full-env inheritance vs. hermeticity:** `append(os.Environ(), ...)` keeps tools functional (PATH, HOME, etc.) but leaks secrets/unintended vars into every child; default-deny filtering only applied to `script_shell` declared args (`script_shell.go:271-273` rejects undeclared keys like `LD_PRELOAD`).
* **Bounded background vs. unbounded shell:** 10 MB bound protects long-lived jobs but shell remains unbounded — an attacker-influenced `cmd` (e.g. `cat huge.bin`) can exhaust agent memory.
* **Job objects vs. process groups:** Windows job kill is atomic for assigned members but assignment race window allows escape; Unix pgid kill is race-free (pgid set by kernel at fork when `Setpgid:true`) but scoped only to pgid members.
* **No cgroup/namespace isolation vs. simplicity:** Avoids root/CapSysAdmin complexity but means CPU/memory/IO/filesystem are not contained; rely on outer Docker sandbox when `--sandbox` is used.

## Failure Modes / Edge Cases

* **Detached grandchild leak:** `setsid sleep 30 &` survives SIGTERM group kill; `WaitDelay` only unblocks `Wait()`, not kill. Test marks this as expected behavior (`shell_test.go:362-397` comment: "process-group kill fallback ... cannot reach it").
* **Windows assignment race:** Child spawns fast grandchildren between `cmd.Start()` and `AssignProcessToJobObject` → those grandchildren escape job and survive parent close.
* **Nested job failure:** Child already belongs to a job (e.g. launched under another job-aware supervisor) → `AssignProcessToJobObject` fails, `reapSpawnedChild` kills the direct child but error returned to caller as "Error creating process group".
* **TOCTOU on working dir:** `checkWorkDir` stat race; dir deleted after check → `chdir` fails inside child, surfaced as generic `Start`/`Run` error.
* **Unbounded shell buffer OOM:** `bytes.Buffer` grows without limit; no scanner split or streaming cap. Mitigated only by caller's timeout.
* **Pipe inherit stall without WaitDelay:** Any `script_shell` or `hooks` command that backgrounds a child remains hung until child dies or parent timeout, because they lack `WaitDelay`.
* **Signal masking:** Processes blocking/ignoring SIGTERM survive first kill; must wait for escalation path. Eval intentionally sends SIGINT instead of SIGTERM because Docker CLI filters SIGKILL/SIGTERM.
* **JS/WASM silent no-op:** `kill` returns error but callers ignore it (`_ = kill`), so background jobs on WASM appear to start then never terminate on stop.
* **Descriptor leak on early return:** `reapSpawnedChild` failure to wait after `SIGKILL` would leak zombie; code always `Wait()` even after `Kill()` in 2 s fallback.

## Future Considerations

* Unify supervision: extract process-group/job logic into a shared `pkg/subprocess` helper used by `shell`, `background_jobs`, `script_shell`, `hooks`, `MCP stdio`, and `environment/cmd_provider` to remove inconsistencies (4 of 6 exec sites currently unsupervised).
* Bound `shell` output similarly to `background_jobs` (configurable `maxOutputBytes`, default e.g. 2–5 MB) with prefix+truncation notice, or stream to temp file instead of `bytes.Buffer`.
* Replace unbounded `bytes.Buffer` + `Write` emit with an `io.Pipe` + scanner or ring buffer to bound memory while still streaming `EmitOutput`.
* Harden Unix termination: on escalation, re-send `Kill(-pgid)` before `Kill(pid)`; optionally set `Pdeathsig: SIGTERM` via `SysProcAttr` so orphaned children die if agent crashes (currently missing).
* Adopt `cmd.Cancel` + `cmd.WaitDelay` (Go 1.20+) uniformly: set `Cancel: func() error { return syscall.Kill(-pid, SIGTERM) }` so `context.WithTimeout` automatically triggers group kill, collapsing duplicated select logic and aligning with eval's pattern.
* Document containment non-guarantees explicitly for each platform (Unix pgid limits, Windows assignment race, js no-op, absence of cgroup/namespace) in tool `Instructions()` — currently only hooks mention it.
* Evaluate optional `cgroup`/`resource` limits for shell tool when running on Linux (e.g. via `systemd-run --scope` or lightweight cgroup wrapper) for Phase 17 broader execution types.

## Questions / Gaps

* No evidence found for cross-platform integration tests that inspect surviving descendants via `ps`/`/proc` or `job` enumeration — only timing regression tests for `WaitDelay`. Search scope: `pkg/tools/builtin/shell/*_test.go`, `backgroundjobs/*_test.go`, `e2e/*`.
* No evidence found for descriptor/resource limits (`RLIMIT_NOFILE`, `RLIMIT_AS`, `Setrlimit`) or fd-close-on-exec hardening beyond pipes.
* `suggest`/`evaluation` container path (`eval.go:390-537`) uses Docker's `--init` + `--privileged` but no per-evaluation cgroup/namespace assertion; whether evaluation depends on Docker sandbox for containment vs. host execution is not stated in code.
* `script_shell` execution bypasses all supervision — is this intentional for "script" toolsets (trusted admin) or an oversight? No comment or test addresses it.
* Privilege boundary for `askpass` grants a Unix socket at `0700` in temp dir with 32-byte hex token; no audit of token entropy or socket permission on shared `/tmp` with `noexec`/`sticky` concerns.

---

Generated by `01.07-subprocess-supervision-and-process-tree-containment` against `docker-agent`.
