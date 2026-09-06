# Source Analysis: dagster

## Crash Recovery, Reconciliation, and Checkpoints

### Source Info

| Field | Value |
|-------|-------|
| Name | dagster |
| Path | `studies/ultraplan-daemon-events-study/sources/dagster` |
| Language / Stack | Python (dagster core + daemon), SQL/In-memory run + event-log storage |
| Analyzed | 2026-09-03 |

## Summary

Dagster persists all recovery-relevant state in the instance databases (run storage, event log storage, tick/backfill storage) so runs, ticks, and backfills survive client, daemon, and host crashes. There is no boot-time reconciliation sweep: recovery is lazy and polling-driven. The run-monitoring daemon each iteration collects all `STARTING`/`STARTED`/`CANCELING`/`NOT_STARTED` runs and applies per-status rules — start/cancel timeouts, dead-worker detection with bounded `resume_run` retries, and forced failure. In-memory execution never resumes mid-call; recovery is either (a) same-`run_id` resume by launching a new run worker (only for launchers with `supports_resume_run`/`supports_check_run_worker_health`), or (b) a new run linked via `parent_run_id`/`root_run_id` that skips completed steps using event-log-derived `KnownExecutionState`. Sensor/asset-daemon ticks and asset-backfill iterations are checkpointed in the DB and resumed idempotently after interruption. Orphan cleanup is best-effort: `terminate()` via gRPC, otherwise force-mark-failed with an explicit "resources may not have been cleaned up" caveat.

## Rating

**7/10** — Clear model with explicit interfaces (`RunLauncher.resume_run` / `check_run_worker_health`, `KnownExecutionState`, tick/backfill checkpoints), operational safeguards (timeouts, max resume attempts, failure resubmission caps), and crash/recovery tests (`test_monitoring_daemon.py`, asset-daemon failure-recovery tests). Deducted because the OSS `DefaultRunLauncher` supports neither health checks nor resume (dead workers are just marked failed), there is no startup reconciliation pass (stale `STARTED` runs wait for the next monitoring poll), and orphan resource reclamation is explicitly non-guaranteed.

## Evidence Collected

Every entry MUST include a file path with line numbers. Format: `path/to/file.ts:NN`.

| Area | Evidence | File:Line |
|------|----------|-----------|
| Run lifecycle states | `DagsterRunStatus` enum (`QUEUED/NOT_STARTED/STARTING/STARTED/SUCCESS/FAILURE/CANCELING/CANCELED`); `IN_PROGRESS_RUN_STATUSES = [STARTING, STARTED, CANCELING]`; cancelable-status comment | `python_modules/dagster/dagster/_core/storage/dagster_run.py:55-123` |
| Durable run record fields | `DagsterRun` carries `parent_run_id`, `root_run_id`, `step_keys_to_execute`, `execution_plan_snapshot_id`, tags; serdes history documents renames/additions | `python_modules/dagster/dagster/_core/storage/dagster_run.py:251-349` |
| Plan snapshot durability | `add_execution_plan_snapshot` / `get_execution_plan_snapshot` content-addressed via `create_execution_plan_snapshot_id`; SQL + in-memory impls | `python_modules/dagster/dagster/_core/storage/runs/sql_run_storage.py:631-647`, `python_modules/dagster/dagster/_core/storage/runs/base.py:270-301` |
| Monitoring sweep (reconciliation loop) | `execute_run_monitoring_iteration` fetches all `IN_PROGRESS + CANCELING + NOT_STARTED` run records each iteration and dispatches to `monitor_starting_run` / `monitor_started_run` / `monitor_canceling_run` | `python_modules/dagster/dagster/_daemon/monitoring/run_monitoring.py:174-223` |
| Start/cancel timeouts | `monitor_starting_run` fails runs stuck in `STARTING`/`NOT_STARTED` past `run_monitoring_start_timeout_seconds` with `RunFailureReason.START_TIMEOUT`; `monitor_canceling_run` force-cancels past `run_monitoring_cancel_timeout_seconds` | `python_modules/dagster/dagster/_daemon/monitoring/run_monitoring.py:28-105` |
| Dead-worker resume loop | `monitor_started_run`: `check_run_worker_health` → if unhealthy and `num_prev_attempts < max_resume_run_attempts`, emit `RESUME_RUN_LOG_MESSAGE` engine event and `instance.resume_run(...)`; else `report_run_failed` | `python_modules/dagster/dagster/_daemon/monitoring/run_monitoring.py:113-171` |
| Resume-attempt counting | `count_resume_run_attempts` counts `ENGINE_EVENT`s with message `RESUME_RUN_LOG_MESSAGE` — the attempt counter is itself a durable event-log checkpoint | `python_modules/dagster/dagster/_daemon/monitoring/run_monitoring.py:25-110` |
| Resume gate | `run_will_resume` returns False unless monitoring enabled and attempts < max; `count_resume_run_attempts` delegates to monitoring module | `python_modules/dagster/dagster/_core/instance/methods/run_launcher_methods.py:196-209` |
| Resume defaults/off | `run_monitoring_max_resume_run_attempts` defaults to 0; `_initialize_run_monitoring` invariants that resume-capable launchers are required when attempts > 0 | `python_modules/dagster/dagster/_core/instance/methods/settings_methods.py:135-217` |
| Launcher capability interface | `RunLauncher.supports_check_run_worker_health` / `check_run_worker_health` / `supports_resume_run` / `resume_run`; base impls return False / raise; `WorkerStatus` enum (`RUNNING/NOT_FOUND/FAILED/SUCCESS/UNKNOWN`) | `python_modules/dagster/dagster/_core/launcher/base.py:38-117` |
| OSS default launcher gap | `DefaultRunLauncher` implements only `launch_run`/`terminate`/`join`; no health/resume overrides, so it inherits "not supported" and dead workers are failed, not resumed | `python_modules/dagster/dagster/_core/launcher/default_run_launcher.py:98-212` |
| Interrupt-vs-resume branch | `job_execution_iterator` finally-block: on `KeyboardInterrupt`/interrupted error, emits resume-note engine event if `run_will_resume`, else `job_failure` with `RunFailureReason.UNEXPECTED_TERMINATION`; special-cases `CANCELING`/`FAILURE` reloads | `python_modules/dagster/dagster/_core/execution/api.py:757-799` |
| Step-level checkpoint derivation | `KnownExecutionState.build_resume_retry_reexecution` / `build_for_reexecution` derive `steps_to_retry`, dynamic mappings, produced-output sets from durable step events (`STEP_START/FAILURE/SUCCESS/OUTPUT/SKIPPED`); interrupted = started-but-no-terminal-event; downstream-of-retry re-queued | `python_modules/dagster/dagster/_core/execution/plan/state.py:169-380` |
| Step-executor resume protocol | `StepDelegatingExecutor`: on `resume_from_failure` replays prior events, `rebuild_from_events`, health-checks in-flight steps and relaunches only unhealthy ones; on interrupt, forwards `terminate_step` unless `run_will_resume` (then leaves steps alone for the new worker) | `python_modules/dagster/dagster/_core/executor/step_delegating/step_delegating_executor.py:187-292` |
| Daemon-crash-safe auto-retry | `retry_run` checks for an already-created child run in the run group (`get_automatically_retried_run_if_exists`) and resubmits instead of duplicating when the daemon died between create and submit; tags `AUTO_RETRY_RUN_ID_TAG` | `python_modules/dagster/dagster/_daemon/auto_run_reexecution/auto_run_reexecution.py:214-229` |
| Tick interruption resume | Sensor daemon reuses prior `STARTED` tick if interrupted part-way (`has_unrequested_runs` and age ≤ `MAX_TIME_TO_RESUME_TICK_SECONDS = 24h`), else moves dangling `STARTED` tick to `SKIPPED`; `FAILURE` ticks retried up to `MAX_FAILURE_RESUBMISSION_RETRIES = 1` | `python_modules/dagster/dagster/_daemon/sensor.py:74-83`, `python_modules/dagster/dagster/_daemon/sensor.py:530-569` |
| Backfill iteration checkpoint | Asset backfill persists `submitting_run_requests`/`reserved_run_ids`; on restart reconstructs in-progress iteration and re-submits; partition backfills checkpoint on `last_submitted_partition_name` | `python_modules/dagster/dagster/_core/execution/asset_backfill.py:1052-1062`, `python_modules/dagster/dagster/_core/execution/job_backfill.py:107-153` |
| Queue-dequeue idempotency | `_dequeue_run` re-checks run is still `QUEUED` before launching; on user-code-unreachable errors re-enqueues with bounded `max_user_code_failure_retries`, else marks failed | `python_modules/dagster/dagster/_daemon/run_coordinator/queued_run_coordinator_daemon.py:363-491` |
| Orphan cleanup limit | `check_run_timeout` calls `run_launcher.terminate` then unconditionally `_force_mark_as_failed`, whose message states "computational resources created by the run may not have been fully cleaned up" | `python_modules/dagster/dagster/_daemon/monitoring/run_monitoring.py:226-289` |
| Termination path | `DefaultRunLauncher.terminate` reports `RUN_CANCELING` and sends `CancelExecutionRequest` over the stored gRPC client; returns False when no client/finished run | `python_modules/dagster/dagster/_core/launcher/default_run_launcher.py:153-184` |
| Daemon liveness (not run recovery) | Interval daemons post `DaemonHeartbeat`s; controller tracks `last_healthy_heartbeat_times` with heartbeat interval + tolerance; duplicate-daemon detection via `daemon_id` mismatch | `python_modules/dagster/dagster/_daemon/daemon.py:147-194`, `python_modules/dagster/dagster/_daemon/controller.py:213-264` |
| Crash/recovery tests | `MockRunLauncher` with `supports_resume_run`/`supports_check_run_worker_health` drives start-timeout, resume-attempt counting (1..3 then fail), and cancel-timeout cases | `python_modules/dagster/dagster_tests/daemon_tests/test_monitoring_daemon.py:37-110`, `python_modules/dagster/dagster_tests/daemon_tests/test_monitoring_daemon.py:224-332` |
| Tick failure-recovery tests | Asset-daemon tests advance time past/before `MAX_TIME_TO_RESUME_TICK_SECONDS` to assert give-up vs resume | `python_modules/dagster/dagster_tests/declarative_automation_tests/daemon_tests/test_asset_daemon_failure_recovery.py:100-135` |

## Answers to Dimension Questions

1. **What survives a client crash, daemon crash, and host reboot?** Everything in instance storage: `DagsterRun` rows (status, `parent/root_run_id`, `step_keys_to_execute`, `execution_plan_snapshot_id`, tags), execution-plan/job snapshots, the append-only event log (step + engine events including `RESUME_RUN_LOG_MESSAGE` attempt markers), daemon ticks/instigator state, and backfill records (`submitting_run_requests`, `reserved_run_ids`, `last_submitted_partition_name`). In-memory state — the executor's `ActiveExecution`, child processes, gRPC worker handles, `DefaultRunLauncher._run_ids` (`default_run_launcher.py:39`), daemon thread pools — does not survive.
2. **Can arbitrary in-memory execution resume, or only checkpointed work?** Only checkpointed work. There is no process snapshot/restore. Same-`run_id` resume means launching a brand-new run worker that replays durable events (`KnownExecutionState` from logs, `step_delegating_executor.py:187-253`); anything not yet emitted as a step event is re-executed. Cross-run recovery (`reexecute_from_failure`, auto-retry) creates a new run and skips already-succeeded steps.
3. **When is a new attempt started automatically?** (a) Monitoring daemon: dead `STARTED` worker → `resume_run` automatically while attempts < `max_resume_run_attempts` (`run_monitoring.py:134-148`). (b) Auto-run-reexecution daemon: terminal failed runs matching policy get a new child run via `retry_run` (`auto_run_reexecution.py:149-258`). (c) Queue daemon: dequeue failures due to unreachable code servers are re-enqueued up to `max_user_code_failure_retries` (`queued_run_coordinator_daemon.py:426-477`). (d) Sensor/asset ticks: interrupted `STARTED` ticks resumed within 24h, failed ticks resubmitted once.
4. **When is human/manual retry required?** When resume budget is exhausted (run marked failed with "surpassed the configured maximum attempts", `run_monitoring.py:150-166`); when the launcher lacks `supports_resume_run` (notably the OSS `DefaultRunLauncher` — invariant in `settings_methods.py:211-217` forces `max_resume_run_attempts=0`, so every dead worker just fails); when start/cancel timeouts fire (`START_TIMEOUT`, forced cancel); when auto-retry preconditions fail (missing origin/repo/job, deleted root run, `retry_on_asset_or_op_failure` opt-out, `auto_run_reexecution.py:31-96`); and for step-level business failures, which need explicit "re-execute from failure" (new run via `build_resume_retry_reexecution`, which rejects non-`FAILURE`/`CANCELED` parents at `state.py:174-177`).
5. **How are orphaned resources discovered?** Polling, not event-driven: the monitoring iteration lists all non-terminal runs and calls launcher `check_run_worker_health` (`WorkerStatus NOT_FOUND/FAILED` ⇒ dead). Step-delegating executors additionally poll `check_step_health`. Discovery therefore depends on launcher support and poll interval (`run_monitoring_poll_interval_seconds`, default 120s); cleanup is `terminate()` (gRPC cancel) with a force-fail fallback that openly may leak resources.

## Architectural Decisions

- **Event log as the checkpoint, not a snapshot store.** Resume state (`steps_to_retry`, dynamic mappings, ready outputs) is recomputed from immutable step events at resume time (`state.py:225-380`) rather than serializing executor memory. Trades extra log scanning for a single source of truth that also powers UI, re-execution, and audit.
- **Lazy polling reconciliation instead of a startup sweep.** No "on boot, reconcile non-terminal runs" pass was found; the monitoring daemon's periodic `get_run_records(IN_PROGRESS...)` query is the reconciler. Simple and stateless across daemon restarts, but detection latency equals the poll interval and stale `STARTED` rows linger until the next iteration.
- **Capability-gated resume.** Health-checking and resume are opt-in launcher capabilities with config validation (`settings_methods.py:207-217`). Keeps core portable across executors (in-process, multiprocess, K8s/Docker in libraries) at the cost that the default OSS path silently degrades to fail-only.
- **Attempt budget recorded as engine events.** Resume attempts are counted from the event log itself (`count_resume_run_attempts`), so the budget survives daemon crashes without a separate counter table — but it also means deleting/compacting events would corrupt the budget.
- **Idempotent daemon writes for crash windows.** Dequeue, auto-retry, tick-submit, and backfill-submit paths all re-read-then-write (run-still-`QUEUED` check, existing-child-run check, unsubmitted-run-ids check, `submitting_run_requests` replay) so a daemon dying mid-write converges instead of duplicating work.

## Notable Patterns

- **Interrupted-step inference:** any step with a start but no terminal event (success/failure/skip) is classified interrupted and retried (`state.py:297-304`, `338-339`).
- **Dangling-tick promotion to SKIPPED** rather than leaving `STARTED` rows forever (`sensor.py:550-555`).
- **Bounded retry everywhere:** resume attempts, failure resubmission (`MAX_FAILURE_RESUBMISSION_RETRIES=1`), user-code dequeue retries, 24h tick-resume horizon.
- **Status-change recheck before acting** (`monitor_started_run` reloads the run and aborts if status moved, `run_monitoring.py:125-133`) to avoid racing the worker.

## Tradeoffs

- **No exactly-once for side effects.** A SIGKILLed worker mid-side-effect leaves an interrupted (not failed) step that will be re-run on resume/re-execution; safety depends on user idempotency, `RetryPolicy`, and IO-manager memoization — the framework does not decide safety, it only bounds attempts and skips provably completed steps.
- **Default deployment fails instead of resumes.** Operators must configure a health/resume-capable launcher plus `run_monitoring.enabled/max_resume_run_attempts`; otherwise crash recovery is "mark failed, retry manually."
- **Poll latency vs load.** 120s default monitoring interval and per-run `get_run_worker_debug_info` calls trade detection speed against DB/RPC load; `start_timeout_seconds`/`cancel_timeout_seconds` (180s defaults) add further delay before stuck runs resolve.
- **Force-fail may leak.** `_force_mark_as_failed` prioritizes state convergence ("don't leave runs stuck") over resource hygiene, explicitly warning cleanup may be incomplete.

## Failure Modes / Edge Cases

- Daemon down longer than `MAX_TIME_TO_RESUME_TICK_SECONDS` (24h): old interrupted ticks are skipped, not resumed — scheduled work for that window is silently dropped to the next tick.
- Daemon dies between auto-retry run creation and submit: handled (resubmit existing `NOT_STARTED` child, `auto_run_reexecution.py:220-229`); but if the child run row is deleted, the daemon retries again (`consume_new_runs_for_automatic_reexecution` docstring, lines 266-271) — potential duplicate lineage.
- `terminate()` with no reachable gRPC client returns False; timeout path still force-fails, leaking the worker process/container.
- Multiprocess executor children (`multiprocess.py:232-392`) are local OS processes: host reboot orphans nothing locally but the run row stays `STARTED` until monitoring fails it — no cross-host adoption.
- Monitoring disabled (`run_monitoring.enabled=False`, the default): `run_will_resume` is always False, so every unexpected interruption becomes `UNEXPECTED_TERMINATION` failure with no automatic recovery.

## Future Considerations

- Add a startup reconciliation pass on daemon boot (verify all non-terminal runs immediately instead of waiting for the next poll; fail or resume stale `STARTING` rows left by a long outage).
- Implement `supports_check_run_worker_health`/`supports_resume_run` for the default OSS launcher (or document the fail-only default prominently with a startup warning).
- Record resume-attempt budget in run storage rather than by counting log messages, so event retention/compaction cannot reset budgets.
- Emit orphan-resource metrics/events on `_force_mark_as_failed` (currently just a log string) so operators can alert on leaks.
- Consider a retry-safety contract (effect idempotency declarations per op) so the SIGKILL-during-side-effect decision can be policy-driven rather than attempt-count-driven.

## Questions / Gaps

- No evidence found of a dedicated crash/restart integration test that SIGKILLs a real worker and asserts resume (searched `dagster_tests/daemon_tests`, `core_tests/instance_tests/test_instance.py:751`; only mock-launcher unit tests and tick time-travel tests found). Real-launcher resume is implemented outside `python_modules/dagster/dagster` (K8s/Docker/ECS launchers in `python_modules/libraries`), which was not traced per single-dimension scope — their `check_run_worker_health` fidelity is unverified here.
- No evidence found for event-log retention interplay with `KnownExecutionState` derivation or resume-attempt counting (what happens if old step events are purged).
- No evidence found for run-storage backend differences (SQLite vs Postgres/MySQL) affecting checkpoint durability semantics.

---
Generated by `Dimension 01.08: Crash Recovery, Reconciliation, and Checkpoints` against `dagster`.
