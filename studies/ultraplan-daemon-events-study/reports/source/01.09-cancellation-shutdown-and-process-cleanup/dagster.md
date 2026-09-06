# Source Analysis: dagster

## Cancellation, Shutdown, and Process Cleanup

### Source Info

| Field | Value |
|-------|-------|
| Name | dagster |
| Path | `studies/ultraplan-daemon-events-study/sources/dagster` |
| Language / Stack | Python (dagster core + dagster-graphql + dagster-shared, gRPC code servers, K8s/Docker/ECS/Celery launchers) |
| Analyzed | 2026-09-03 |

## Summary

Dagster models cancellation as a durable persisted intent (`CANCELING` run status backed by a `RUN_CANCELING` event) plus cooperative runtime acknowledgement via `SIGINT`/interrupt polling. Client path is `GraphQL terminateRun → RunCoordinator.cancel_run → RunLauncher.terminate → gRPC CancelExecution → multiprocessing.Event.set → termination-thread sends SIGINT`. A polling `kill-on-cancel` thread covers in-process runs; multiprocess/step workers poll a per-step `term_event`. The monitoring daemon bounds cleanup with `start_timeout`/`cancel_timeout`/`max_runtime` escalation that force-marks runs `CANCELED`/`FAILED` even if compute leaks. Daemon and gRPC-server shutdown is drain-oriented (`ShutdownServer → shutdown_once_executions_finish → grace period → stop`), not cancel-oriented; code-server processes are detached by default (`wait_for_local_processes_on_shutdown=False`). There is no process-group/cgroup kill and no SIGKILL escalation in the hot path — termination is one SIGINT with no retry.

## Rating

**7/10 — Clear model with tests, explicit interfaces, and operational safeguards, but fragile at the OS process-tree boundary.**

Rationale: durable `CANCELING` intent, explicit `terminate`/`CancelExecution`/`CanCancelExecution` interfaces, per-layer polling/propagation, daemon-enforced cancel/start timeouts, and cancellation-specific tests are all implemented. Downgraded from 9 because (a) runtime kill is cooperative SIGINT-only with no escalation deadline, (b) grandchild/orphan processes can escape (no `killpg`/`cgroup`, plain `Popen`/`multiprocess.Process` on Unix), (c) forced terminal marking explicitly warns resources "may not have been fully cleaned up," and (d) cleanup can hang indefinitely when timeouts are disabled (`cancel_timeout_seconds=0` skips enforcement).

## Evidence Collected

Every entry MUST include a file path with line numbers. Format: `path/to/file.ts:NN`.

| Area | Evidence | File:Line |
|------|----------|-----------|
| Cancellation state/events | `CANCELING`/`CANCELED` statuses; `IN_PROGRESS=[STARTING,STARTED,CANCELING]`; `FINISHED=[SUCCESS,FAILURE,CANCELED]`; `CANCELABLE=[STARTED,QUEUED]` only | `python_modules/dagster/dagster/_core/storage/dagster_run.py:80-123` |
| Cancellation state/events | `EVENT_TYPE_TO_PIPELINE_RUN_STATUS` maps `RUN_CANCELING→CANCELING`, `RUN_CANCELED→CANCELED`; alias `PIPELINE_CANCELING=RUN_CANCELING` | `python_modules/dagster/dagster/_core/events/__init__.py:157-270` |
| Cancellation state/events | `report_run_canceling` emits `PIPELINE_CANCELING` event; `report_run_canceled`/`report_run_failed` emit terminal events with "from outside the execution context" messages | `python_modules/dagster/dagster/_core/instance/methods/event_methods.py:374-442` |
| Cancellation state/events | `handle_run_event` translates terminal events to run status, sets `end_time` on `PIPELINE_CANCELED/FAILURE/SUCCESS`; status column authoritative over body on concurrent writes | `python_modules/dagster/dagster/_core/storage/runs/sql_run_storage.py:173-268` |
| Client request → durable intent | `terminate_pipeline_execution` GraphQL resolver: permission check, `CANCELABLE` gate, `MARK_AS_CANCELED_IMMEDIATELY` force path vs `run_coordinator.cancel_run` path | `python_modules/dagster-graphql/dagster_graphql/implementation/execution/__init__.py:91-164` |
| Client request → durable intent | Force-mark path `_force_mark_as_canceled` reloads record, calls `report_run_canceled` with "forcibly marked as canceled ... may not have been fully cleaned up" | `python_modules/dagster-graphql/dagster_graphql/implementation/execution/__init__.py:72-88` |
| Client request → durable intent | `QueuedRunCoordinator.cancel_run`: `QUEUED→report_run_canceling→report_run_canceled` synchronously (no worker needed); else delegates to `run_launcher.terminate`; notes race `issues/3323` | `python_modules/dagster/dagster/_core/run_coordinator/queued_run_coordinator.py:327-341` |
| Client request → durable intent | `DefaultRunLauncher.terminate`: returns `False` if no instance/run/finished; else `report_run_canceling` (durable) then `CancelExecution` RPC; `False` + engine event if no gRPC client | `python_modules/dagster/dagster/_core/launcher/default_run_launcher.py:153-184` |
| Context/token propagation | `CancelExecution` RPC sets per-run `multiprocessing.Event` + records `_termination_times[run_id]` under `_execution_lock`; `CanCancelExecution` returns `in _executions and not terminated` | `python_modules/dagster/dagster/_grpc/server.py:1289-1336` |
| Context/token propagation | Run worker subprocess started as `multiprocessing.Process(target=start_run_in_subprocess, args=[..., termination_event])`; registered in `_executions`/`_termination_events` maps | `python_modules/dagster/dagster/_grpc/server.py:1380-1400` |
| Context/token propagation | Subprocess `_run_in_subprocess` calls `start_termination_thread(termination_event, done_event)`; on exit sets both events to stop watcher | `python_modules/dagster/dagster/_grpc/impl.py:188-263` |
| Context/token propagation | `_termination_handler`: `should_stop.wait(); if not done: send_interrupt()` (SIGINT to self); `start_termination_thread` spawns daemon `termination-handler` thread | `python_modules/dagster/dagster/_utils/__init__.py:353-382` |
| Context/token propagation | `send_interrupt`: Windows `thread.interrupt_main()`, Unix `os.kill(pid, SIGINT)` — single signal, no escalation | `python_modules/libraries/dagster-shared/dagster_shared/ipc.py:64-71` |
| Context/token propagation | In-process run path `start_run_cancellation_thread` polls `cancellation_thread_poll_interval_seconds` (default 10s) for `CANCELING/CANCELED/FAILURE` then `send_interrupt()` | `python_modules/dagster/dagster/_core/execution/run_cancellation_thread.py:9-44` |
| Context/token propagation | Multiprocess parent on `check_for_interrupts()`: emits "forwarding to active child processes", `mark_interrupted()`, sets each live step's `term_event`; child has own `start_termination_thread` | `python_modules/dagster/dagster/_core/executor/multiprocess.py:227-240` |
| Context/token propagation | Step child `MultiprocessExecutorChildProcessCommand.execute` starts termination thread around `execute_plan_iterator`; sets `done+term` in `finally` | `python_modules/dagster/dagster/_core/executor/multiprocess.py:76-111` |
| Context/token propagation | Interrupt checkpoints: `capture_interrupts()` at process entry (`start_run_in_subprocess`, `ExecuteRunWithPlanIterable`, `dagster api` commands); `raise_execution_interrupts()` during step bodies; `check_for_interrupts/pop_captured_interrupt` polling | `python_modules/dagster/dagster/_grpc/impl.py:266-276`, `python_modules/dagster/dagster/_core/execution/api.py:862-867`, `python_modules/dagster/dagster/_core/execution/plan/active.py:548`, `python_modules/dagster/dagster/_utils/interrupts.py:36-98` |
| Shutdown handlers | `DagsterDaemonController.__exit__`: logs reason, sets `_daemon_shutdown_event`, `thread.join(timeout=30)` per daemon, logs "did not shut down gracefully" but does not kill | `python_modules/dagster/dagster/_daemon/controller.py:323-344` |
| Shutdown handlers | Daemon loop `check_daemon_loop` wraps sleep in `raise_interrupts_as(KeyboardInterrupt)`; sensor/scheduler ticks accept `user_interrupted` and resume partial ticks | `python_modules/dagster/dagster/_daemon/controller.py:293-321`, `python_modules/dagster/dagster/_daemon/sensor.py:167-192` |
| Shutdown handlers | gRPC `ShutdownServer` sets `_shutdown_once_executions_finish_event`; `StartRun` rejected after shutdown; cleanup thread drains `_executions` then sets `_server_termination_event` | `python_modules/dagster/dagster/_grpc/server.py:1266-1287`, `python_modules/dagster/dagster/_grpc/server.py:1338-1350`, `python_modules/dagster/dagster/_grpc/server.py:584-611` |
| Shutdown handlers | `server_termination_target`: `server.stop(grace=shutdown_grace_period)` then `wait(grace+5)`; warns "did not shut down cleanly" — bounded drain | `python_modules/dagster/dagster/_grpc/server.py:1477-1489` |
| Shutdown handlers | Shutdown grace defaults to `max(grpc, schedule, sensor timeouts)`, overridable via `DAGSTER_GRPC_SHUTDOWN_GRACE_PERIOD` | `python_modules/dagster/dagster/_grpc/utils.py:133-147` |
| Shutdown handlers | `GrpcServerRegistry.__exit__`: `shutdown_all_processes()` (sends `ShutdownServer` RPC) then optionally `wait_for_processes()` only if `wait_for_local_processes_on_shutdown=True` (default `False`) | `python_modules/dagster/dagster/_core/remote_representation/grpc_server_registry.py:264-295`, `python_modules/dagster/dagster/_core/instance/methods/settings_methods.py:131-132` |
| Shutdown handlers | Run-worker CLI `_shutdown_threads`: `cancellation_shutdown.set(); join(timeout=15)` else engine event "did not shutdown gracefully"; always emits "Process for run exited" | `python_modules/dagster/dagster/_cli/api.py:211-240` |
| Signal/process-group/cgroup code | `setup_interrupt_handlers` maps `SIGTERM→SIGINT` handler (k8s); no `setpgrp`/`setsid`/`killpg`/cgroup anywhere in core execution path | `python_modules/dagster/dagster/_utils/interrupts.py:16-28` |
| Signal/process-group/cgroup code | `open_ipc_subprocess` sets `CREATE_NEW_PROCESS_GROUP` only on Windows; Unix gets plain `Popen` with no session/process-group isolation | `python_modules/libraries/dagster-shared/dagster_shared/ipc.py:20-39` |
| Signal/process-group/cgroup code | `interrupt_ipc_subprocess` sends single `SIGINT`/`CTRL_BREAK`; `interrupt_then_kill_ipc_subprocess` (SIGINT then `kill` after 10s) exists but no callers in hot cancel path | `python_modules/libraries/dagster-shared/dagster_shared/ipc.py:42-53` |
| Signal/process-group/cgroup code | `get_terminate_signal` returns `SIGKILL` (posix) but is only used for crash-explanation messaging, not for escalation | `python_modules/dagster/dagster/_utils/__init__.py:614-639` |
| Kill escalation deadlines | `monitor_canceling_run`: if `now - RUN_CANCELING_event >= cancel_timeout_seconds` (default 180s) → `report_run_canceled` with timeout message + debug info; skipped entirely if `cancel_timeout_seconds==0` | `python_modules/dagster/dagster/_daemon/monitoring/run_monitoring.py:70-105`, `python_modules/dagster/dagster/_core/instance/methods/settings_methods.py:107-108`, `python_modules/dagster/dagster/_daemon/monitoring/run_monitoring.py:208-213` |
| Kill escalation deadlines | `monitor_starting_run`: `STARTING/NOT_STARTED` older than `start_timeout_seconds` (default 180s) → `report_run_failed(START_TIMEOUT)` | `python_modules/dagster/dagster/_daemon/monitoring/run_monitoring.py:28-67` |
| Kill escalation deadlines | `check_run_timeout` (`max_runtime_seconds` tag or instance default): `report_run_canceling` → `launcher.terminate` → `report_run_failed` → unconditional `_force_mark_as_failed` ("may not have been fully cleaned up") | `python_modules/dagster/dagster/_daemon/monitoring/run_monitoring.py:226-289` |
| Container launchers | K8s `terminate`: `report_run_canceling` then `delete_job`; Docker `terminate`: `report_run_canceling` then `container.stop(timeout=stop_timeout)`; both return `False` if `is_finished` | `python_modules/libraries/dagster-k8s/dagster_k8s/launcher.py:352-378`, `python_modules/libraries/dagster-docker/dagster_docker/docker_run_launcher.py:206-229` |
| Cleanup/finalization tests | `test_kill_on_cancel_sends_interrupt_for_terminal_statuses` covers `CANCELING/CANCELED/FAILURE` → `send_interrupt` | `python_modules/dagster/dagster_tests/execution_tests/misc_execution_tests/test_run_cancellation_thread.py:19-36` |
| Cleanup/finalization tests | `test_monitor_canceling`: fresh `CANCELING` stays, 1000s-old `CANCELING` → `CANCELED` | `python_modules/dagster/dagster_tests/daemon_tests/test_monitoring_daemon.py:256-286` |
| Cleanup/finalization tests | `test_terminated_run`: `launcher.terminate` → poll → `CANCELED`; second `terminate` returns `False`; `PIPELINE_CANCELED` event asserted | `python_modules/dagster/dagster_tests/launcher_tests/test_default_run_launcher.py:500-532` |
| Cleanup/finalization tests | `MockRunLauncher.terminate` + `monitor_canceling_run` coverage in monitoring daemon tests | `python_modules/dagster/dagster_tests/daemon_tests/test_monitoring_daemon.py:70-80` |

## Answers to Dimension Questions

1. **Is client disconnect different from explicit cancellation? — Yes. Disconnect does nothing; only explicit termination cancels.**
   Evidence: the only state-changing entry points are `terminate_pipeline_execution` (`python_modules/dagster-graphql/dagster_graphql/implementation/execution/__init__.py:91-164`), `RunCoordinator.cancel_run` (`python_modules/dagster/dagster/_core/run_coordinator/queued_run_coordinator.py:327-341`), `RunLauncher.terminate` (`python_modules/dagster/dagster/_core/launcher/default_run_launcher.py:153-184`), and `CancelExecution` (`python_modules/dagster/dagster/_grpc/server.py:1289-1317`). GraphQL log subscription `gen_events_for_run` only `watch_event_logs`/`end_watch_event_logs` (`python_modules/dagster-graphql/dagster_graphql/implementation/execution/__init__.py:222-309`) with no status write. gRPC server heartbeat loss triggers server self-shutdown (`python_modules/dagster/dagster/_grpc/server.py:564-574`), not run cancellation — orphaned runs are later marked `FAILED` by `_check_for_orphaned_runs` (`python_modules/dagster/dagster/_grpc/server.py:584-605`), never `CANCELED`.

2. **Is cancellation durable if no worker is currently attached? — Yes for intent; runtime acknowledgement still needs a worker or the monitoring daemon.**
   Evidence: `report_run_canceling` persists a `PIPELINE_CANCELING` event → `CANCELING` status via `handle_run_event` (`python_modules/dagster/dagster/_core/instance/methods/event_methods.py:374-390`, `python_modules/dagster/dagster/_core/storage/runs/sql_run_storage.py:173-193`). `QueuedRunCoordinator.cancel_run` handles `QUEUED` with no worker at all (`python_modules/dagster/dagster/_core/run_coordinator/queued_run_coordinator.py:327-339`). `DefaultRunLauncher.terminate` writes `CANCELING` before attempting the RPC and returns `False` with an engine event if no gRPC client exists (`python_modules/dagster/dagster/_core/launcher/default_run_launcher.py:161-175`) — intent survives, and `monitor_canceling_run` will force `CANCELED` after `cancel_timeout_seconds` (`python_modules/dagster/dagster/_daemon/monitoring/run_monitoring.py:70-105`).

3. **Can cleanup hang indefinitely? — Yes, in three configurations; bounded by default.**
   Evidence: (a) If `cancel_timeout_seconds=0`, `execute_run_monitoring_iteration` skips `monitor_canceling_run` entirely (`python_modules/dagster/dagster/_daemon/monitoring/run_monitoring.py:208-213`), so a worker ignoring SIGINT stays `CANCELING` forever. (b) Cooperative SIGINT only: `_termination_handler`/`send_interrupt` send one interrupt with no retry or SIGKILL (`python_modules/dagster/dagster/_utils/__init__.py:353-360`, `python_modules/libraries/dagster-shared/dagster_shared/ipc.py:64-71`); user code swallowing `KeyboardInterrupt`/`DagsterExecutionInterruptedError` or blocked in non-interruptible C calls never yields. (c) Daemon shutdown is best-effort: `thread.join(timeout=30)` then log-and-continue (`python_modules/dagster/dagster/_daemon/controller.py:338-343`); code-server processes are detached unless `wait_for_local_processes_on_shutdown=True` (`python_modules/dagster/dagster/_core/remote_representation/grpc_server_registry.py:271-272`). Bounded by default via 180s cancel/start timeouts and `server.stop(grace)` (`python_modules/dagster/dagster/_grpc/server.py:1477-1489`).

4. **Can child/grandchild processes escape termination? — Yes.**
   Evidence: No `setsid`/`setpgrp`/`killpg`/cgroup logic in the cancel path; `open_ipc_subprocess` isolates process groups only on Windows (`python_modules/libraries/dagster-shared/dagster_shared/ipc.py:20-39`). Unix run workers are bare `multiprocessing.Process` (`python_modules/dagster/dagster/_grpc/server.py:1382-1390`) and step workers are bare `multiproc_ctx.Process` (`python_modules/dagster/dagster/_core/executor/child_process_executor.py:151-154`). Parent-to-child signalling is `term_event.set()` + child-side SIGINT to its own main thread — grandchildren spawned by user code (shell-outs, `Popen`, Spark/Dask workers) receive nothing. Multiprocess parent on interrupt deletes `processes[key]` entries after setting events without `terminate()`/`kill()` (`python_modules/dagster/dagster/_core/executor/multiprocess.py:236-240`). Container launchers are the exception: K8s `delete_job` and Docker `container.stop` reclaim the whole container (`python_modules/libraries/dagster-k8s/dagster_k8s/launcher.py:352-378`, `python_modules/libraries/dagster-docker/dagster_docker/docker_run_launcher.py:206-229`).

5. **How are cancellation and completion races arbitrated? — By reload-and-branch on persisted status plus `is_finished` gates; last terminal event wins in storage.**
   Evidence: `execute_run_iterator`'s `finally` reloads the run: `CANCELING→job_canceled`, already-`CANCELED`→benign engine event, `FAILURE`→interrupted-but-already-failed event, else `job_failure(UNEXPECTED_TERMINATION)` (`python_modules/dagster/dagster/_core/execution/api.py:766-799`); same reload pattern for exception and step-failure paths (`python_modules/dagster/dagster/_core/execution/api.py:800-831`). Entry gates refuse to act on finished runs: `terminate` returns `False` if `is_finished` (`python_modules/dagster/dagster/_core/launcher/default_run_launcher.py:161-163`), GraphQL returns `TerminateRunFailure` if `is_finished` or not `CANCELABLE` (`python_modules/dagster-graphql/dagster_graphql/implementation/execution/__init__.py:133-141`). `CanCancelExecution` is `False` once the run leaves `_executions` (`python_modules/dagster/dagster/_grpc/server.py:1319-1336`, `python_modules/dagster/dagster/_grpc/server.py:1437-1441`). Monitoring loop re-reads status to avoid acting on a concurrently changed run (`python_modules/dagster/dagster/_daemon/monitoring/run_monitoring.py:125-133`). Storage note warns the status column wins over a stale body on concurrent writes (`python_modules/dagster/dagster/_core/storage/runs/sql_run_storage.py:261-267`).

> **Scenario answer (cancel races completion; child ignores first SIGINT):** terminal outcome is whichever terminal event commits first. If the success event commits first, later `terminate`/`CancelExecution` no-ops (`is_finished→False`, `CanCancel→False`) and the run stays `SUCCESS`. If `CANCELING` commits first, the worker's `finally` emits `job_canceled` (or a "cleaned up after forcibly marked as canceled" engine event if the daemon already forced `CANCELED`). Cleanup behavior: exactly one SIGINT is delivered; an ignoring child is **not** re-signalled or SIGKILLed — the run sits in `CANCELING` until the monitoring daemon's `cancel_timeout_seconds` (default 180s) force-marks it `CANCELED`/`FAILED` with an explicit "computational resources ... may not have been fully cleaned up" warning, potentially orphaning the still-running child (reported as `FAILED` later via orphan detection if its parent process exits).

## Architectural Decisions

- **Durable intent + cooperative interrupt:** `CANCELING` is a first-class persisted status; runtime effect is always a polled event + SIGINT, never a synchronous kill. Traced from `event_methods.py:374-390` through `server.py:1300-1303` to `ipc.py:64-71`.
- **Two polling loops instead of one push:** in-process runs poll the DB (`run_cancellation_thread.py:9-30`, default 10s); out-of-process runs use an in-memory `multiprocessing.Event` pushed over gRPC (`server.py:1380-1400`). Keeps cancellation working across process boundaries without a persistent connection.
- **Daemon as backstop, not actor:** the monitoring daemon never sends SIGINT itself; it only `terminate()`s once (max-runtime path) and otherwise force-marks terminal status after timeouts (`run_monitoring.py:70-105,226-289`). Actual interruption is always worker-local.
- **Drain-on-shutdown everywhere:** gRPC servers, daemon controller, and code-server registry all prefer draining in-flight work over cancelling it (`server.py:1477-1489`, `controller.py:323-344`, `grpc_server_registry.py:264-295`).
- **Launcher-per-environment termination:** each `RunLauncher` owns the kill primitive (gRPC event, K8s job delete, Docker stop), but all share the `report_run_canceling`-first convention.

## Notable Patterns

- **Cancel-before-RPC convention:** every `terminate` writes `CANCELING` before touching compute (`default_run_launcher.py:165`, `queued_run_coordinator.py:334-338`, k8s `launcher.py:359`, docker `docker_run_launcher.py:212`), making intent durable even when the kill RPC fails.
- **Done-event handshake:** every `start_termination_thread` site pairs `termination_event` with a `done_event` and sets both on exit (`impl.py:195-226`, `multiprocess.py:108-111`) to avoid interrupting a process that already finished.
- **`capture_interrupts` / `raise_interrupts_as` sandwich:** defects SIGINT into a flag at process boundaries, re-raises as `DagsterExecutionInterruptedError` only inside step bodies (`interrupts.py:36-98`, `api.py:862-867`), so cleanup code outside checkpoints is not torn down mid-write.
- **Reload-before-decide:** cancellation/completion arbitration always re-reads the run from storage before emitting the terminal event (`api.py:767-831`, `run_monitoring.py:125-133`, `execution/__init__.py:78-86`).

## Tradeoffs

- **Cooperative vs preemptive kill:** single-SIGINT preserves event-log integrity and lets `finally` blocks flush, at the cost of no guarantee against ignore/blocked-SIGINT user code. No SIGKILL escalation exists on the path (only unused `interrupt_then_kill_ipc_subprocess`, `ipc.py:47-52`).
- **Poll latency vs DB load:** 10s default `cancellation_thread_poll_interval_seconds` and 120s `run_monitoring_poll_interval_seconds` bound responsiveness; faster cancellation costs more run-store reads.
- **Detach-by-default code servers:** `wait_for_local_processes_on_shutdown=False` keeps daemon shutdown fast but leaves user-code servers to heartbeat-expiry self-shutdown (`server.py:564-574`).
- **Narrow `CANCELABLE` set (`STARTED, QUEUED`):** avoids acting on `STARTING/NOT_STARTED` where control has been ceded to the worker, at the cost of a window where UI termination is rejected with `TerminateRunFailure` even though the run is not finished (`dagster_run.py:118-123`).
- **Timeout-forced terminal states trade correctness for liveness:** after 180s the system declares `CANCELED`/`FAILED` while compute may still run, explicitly accepting orphaned resources over a stuck `CANCELING` row.

## Failure Modes / Edge Cases

- **SIGINT swallowed:** op catching `KeyboardInterrupt` or running non-interruptible native code stays alive; run pinned in `CANCELING` until daemon timeout; then orphaned (double-billing / zombie container).
- **Grandchild leak:** shell-outs, local subprocess pools, or external systems (Spark, Dask, DB transactions) are never signalled; no `atexit`/context-manager guarantee spans process death. Only container launchers (K8s/Docker) contain this.
- **Cancel lands during `STARTING`:** `CANCELABLE` excludes `STARTING`, so GraphQL rejects; worker will still transition to `STARTED` and only then become cancelable — user must retry termination.
- **Orphaned run process:** if the gRPC execution process dies without terminal events, `_check_for_orphaned_runs` marks `FAILED` (not `CANCELED`), so a run cancelled at the exact crash moment surfaces as failure (`server.py:584-605`).
- **Daemon thread join expiry:** `join(timeout=30)` failure only logs; a wedged sensor/scheduler thread can block daemon exit indefinitely while holding heartbeats stale.
- **Status-column/body skew:** concurrent `handle_run_event` + `add_run_tags` can leave the serialized body stale; readers must use the status column (`sql_run_storage.py:261-267`).
- **Known race on dequeue-cancel:** `QueuedRunCoordinator.cancel_run` documents a dequeuer-vs-canceller race (`queued_run_coordinator.py:331`), mitigated only by the monitoring loop's status recheck.

## Future Considerations

- Add bounded SIGINT→SIGTERM→SIGKILL escalation with per-launcher deadlines (wire the existing `interrupt_then_kill_ipc_subprocess` 10s pattern into `DefaultRunLauncher.terminate` and multiprocess `term_events`).
- Adopt process-group/cgroup containment for local/multiprocess executors (`setsid` + `killpg` on Unix, job objects on Windows) so grandchildren cannot escape; track spawned PIDs for audit.
- Make `cancel_timeout_seconds=0` mean "never force-mark" explicitly in docs, or forbid it — currently it silently disables the only backstop against `CANCELING` leaks.
- Widen `CANCELABLE_RUN_STATUSES` or queue `CANCELING` intent for `STARTING` runs so UI termination during launch does not require a retry.
- Emit `termination_times`/escalation-attempt counters as run events/metrics for observability (currently only in-memory `_termination_times`, never logged).

## Questions / Gaps

- No evidence found for client-disconnect-triggered cancellation (searched GraphQL subscriptions, `watch_event_logs`, gRPC heartbeat handling) — disconnect appears to be a no-op for runs. Confirm whether webserver/proxy teardown paths ever call `terminate`.
- No evidence found for cleanup-context deadlines inside user code (e.g., a bounded `finally`/teardown budget after interrupt) — searched `active.py`, `api.py`, `plan/utils.py`, executor paths. Appears unbounded by design.
- No evidence found for cgroup/process-group termination or grandchild reaping in OSS launchers — `killpg`, `setsid`, `start_new_session`, `cgroup` returned no hits in `dagster/` core. ECS/Celery launchers not inspected beyond `terminate` signature existence.
- Open: does `report_run_canceling`'s second call overwrite the original `RUN_CANCELING` timestamp used by `monitor_canceling_run` for timeout math? `get_logs_for_run(limit=1, ascending=False)` reads the latest event (`run_monitoring.py:76-86`), so repeated `terminate` calls may reset the 180s clock — worth a targeted test.

---

Generated by `Dimension 01.09: Cancellation, Shutdown, and Process Cleanup` against `dagster`.
