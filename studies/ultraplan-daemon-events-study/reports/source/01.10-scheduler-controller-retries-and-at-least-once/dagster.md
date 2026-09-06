# Source Analysis: dagster

## Dimension 01.10: Scheduler Controller, Retries, and At-Least-Once Semantics

### Source Info

| Field | Value |
|-------|-------|
| Name | dagster |
| Path | `studies/ultraplan-daemon-events-study/sources/dagster` |
| Language / Stack | Python (dagster daemon, scheduler, sensor, run coordinator, execution engine) |
| Analyzed | 2026-09-03 |

## Summary

Dagster implements a reconciliation-style daemon fleet (scheduler, sensor, queued-run-coordinator, backfill, asset, monitoring daemons) rather than a single scheduler controller. Each daemon is an infinite generator loop (`DagsterDaemon.run_daemon_loop`) with heartbeat-based liveness, interval jitter, and error-ring capture. Desired state comes from code locations (schedules/sensors/jobs); observed state comes from the instance DB (ticks, runs, heartbeats). Scheduling is tick-based (`TickStatus STARTED/SKIPPED/SUCCESS/FAILURE` with per-tick `failure_count` and cross-tick `consecutive_failure_count`); dequeue is queue-based with `max_concurrent_runs`, tag limits, op-concurrency pools, and stable priority sort. Retries are explicitly stacked across five layers — op (`RetryPolicy`/`RetryRequested`), run auto-reexecution (`dagster/max_retries`), dequeue re-queue (`max_user_code_failure_retries`), tick redo (`max_tick_retries`), and run-health resume (`max_resume_run_attempts`) — each with its own budget. Failure classification distinguishes user-code-unreachable (infinite/pause-and-retry) from framework errors (permanent). Duplicate execution is assumed: scheduler dedups on `(scheduled_execution_time, run_key, repo selector)`, sensors on `run_key`, dequeue on re-fetch of `QUEUED/STARTING` status, and crash-recovery tests assert single-run resumption. A documented race (scheduler crash between launch and executor marking STARTED) can double-launch.

## Rating

**8 / 10** — Clear model with tests, explicit interfaces, and operational safeguards.

Rationale: controller loops, capacity models, priority ordering, layered retry budgets, backoff/jitter, and idempotency keys are all implemented (not just documented) with dedicated failure-recovery tests (scheduler/sensor crash matrices, queued-coordinator priority/retry tests). Deducted from 9–10 because retries are stacked rather than centralized (budgets interact, e.g. manual retries counted in run-group race), backoff utility has explicitly no jitter, and one at-least-once race is acknowledged unfixed in-test comment.

## Evidence Collected

Every entry MUST include a file path with line numbers. Format: `path/to/file.ts:NN`.

| Area | Evidence | File:Line |
|------|----------|-----------|
| Controller loop | `DagsterDaemon.run_daemon_loop`: `while not shutdown.is_set()`, heartbeat check, `next(generator)`, on Exception capture + `wait(error_sleep)` + rebuild generator; StopIteration treated as error | `python_modules/dagster/dagster/_daemon/daemon.py:88` |
| Controller error sleep | `DAGSTER_DAEMON_CORE_LOOP_EXCEPTION_SLEEP_INTERVAL=5`, error ring `deque(maxlen=5)` | `python_modules/dagster/dagster/_daemon/daemon.py:42` |
| Interval scheduling | `IntervalDaemon.core_loop`: `startup_jitter + interval+random.uniform(0,interval_jitter)`, START/END_SPAN markers, `run_iteration()` abstract | `python_modules/dagster/dagster/_daemon/daemon.py:227` |
| Daemon fleet wiring | `SchedulerDaemon/SensorDaemon/BackfillDaemon/MonitoringDaemon` delegate to `execute_scheduler_iteration_loop`, `execute_sensor_iteration_loop`, `execute_backfill_iteration_loop`, `execute_run_monitoring_iteration` | `python_modules/dagster/dagster/_daemon/daemon.py:268` |
| Controller supervision | `DagsterDaemonController`: one thread per daemon `Thread(target=daemon.run_daemon_loop, daemon=True)`; `check_daemon_loop` sleeps `THREAD_CHECK_INTERVAL=5`, checks thread liveness, workspace freshness (`RELOAD=60`, tolerance 300), heartbeats every 60s | `python_modules/dagster/dagster/_daemon/controller.py:133` |
| Heartbeat / split-brain detect | `_check_add_heartbeat`: expire errors older than error_interval, skip write if `daemon_skip_heartbeats_without_errors`, log error if `last_stored_heartbeat.daemon_id != daemon_uuid` (multiple daemons unsupported) | `python_modules/dagster/dagster/_daemon/daemon.py:147` |
| Scheduler tick loop | `execute_scheduler_iteration_loop`: infinite `while True`, next iteration aligned to minute boundary (min cron granularity), on error `min(start+5s, next)` | `python_modules/dagster/dagster/_scheduler/scheduler.py:194` |
| Desired-vs-observed | Observed=`instance.all_instigator_state(SCHEDULE)`; desired=`workspace.get_code_location_entries()->get_schedules()`; create missing `InstigatorState(DECLARED_IN_CODE, ScheduleInstigatorData(cron,now))`; delete orphans immediately if DECLARED_IN_CODE else 12h grace `RETAIN_ORPHANED_STATE_INTERVAL_SECONDS`, skip error locations | `python_modules/dagster/dagster/_scheduler/scheduler.py:291` |
| Tick state | `_ScheduleLaunchContext.update_state()` writes `TickStatus` + `failure_count/consecutive_failure_count`; reset to 0 on SKIPPED/SUCCESS; purge ticks on exit | `python_modules/dagster/dagster/_scheduler/scheduler.py:80` |
| Tick enum | `TickStatus{STARTED,SKIPPED,SUCCESS,FAILURE}`; `InstigatorTick.with_status` sets `end_timestamp=None` while STARTED | `python_modules/dagster/dagster/_core/scheduler/instigation.py:305` |
| Tick history counters | `TickData(failure_count, consecutive_failure_count)` per-tick vs cross-tick | `python_modules/dagster/dagster/_core/scheduler/instigation.py:605` |
| Capacity model | `_get_runs_to_dequeue`: `max_runs_to_launch=max_concurrent_runs-len(in_progress)`; `-1` disables; paginate `get_runs(QUEUED, ascending=True, limit=100)`; `batch[:max_runs_to_launch]` | `python_modules/dagster/dagster/_daemon/run_coordinator/queued_run_coordinator_daemon.py:192` |
| Priority / fairness | `_priority_sort`: `int(tags[dagster/priority]\|\|0)`, `sorted(reverse=True)`, stable → FIFO within priority | `python_modules/dagster/dagster/_daemon/run_coordinator/queued_run_coordinator_daemon.py:345` |
| Tag + pool quotas | `TagConcurrencyLimitsCounter` + `GlobalOpConcurrencyLimitsCounter(op_concurrency_slot_buffer, pool_granularity)`; `is_blocked()` drops from batch; in-batch accounting via `update_counters_with_launched_item()`; fallback to unblocked on counter-init failure | `python_modules/dagster/dagster/_daemon/run_coordinator/queued_run_coordinator_daemon.py:269` |
| Queue config | `RunQueueConfig(max_concurrent_runs, max_user_code_failure_retries=0, user_code_failure_retry_delay=60, should_block_op_concurrency_limited_runs, op_concurrency_slot_buffer, max_concurrent_runs_all_branch_deployments)`; `QueuedRunCoordinator(max_concurrent_runs=10, dequeue_interval=5)` | `python_modules/dagster/dagster/_core/run_coordinator/queued_run_coordinator.py:109` |
| Instance concurrency config | `concurrency.runs.max_concurrent_runs`, `tag_concurrency_limits[{key,value,limit,applyLimitPerUniqueValue}]`, `pools.{granularity:OP\|RUN, default_limit, op_granularity_run_buffer}`, `run_retries{enabled,max_retries,retry_on_asset_or_op_failure}` | `python_modules/dagster/dagster/_core/instance/config.py:473` |
| Starvation avoidance | `shuffled_round_robin_by_key()`: shuffle within/across groups + round-robin interleave; used for scheduler/sensor/asset fan-out to thread pool | `python_modules/dagster/dagster/_daemon/utils.py:17` |
| Tick concurrency guard | Only one tick per schedule in flight (`scheduler_run_futures`); catch-up clamp `tick_times[-max_catchup_runs:]`; non-partitioned no catch-up | `python_modules/dagster/dagster/_scheduler/scheduler.py:418` |
| Backfill ordering | One backfill per job in flight; jobs `sorted(backfill_timestamp)` | `python_modules/dagster/dagster/_daemon/backfill.py:142` |
| Backfill failure taxonomy | `_is_retryable_backfill_error`: retry `UserCodeUnreachable, CodeLocationLoadError, RunAlreadyExists`; `not retry DagsterError\|CheckError` (framework invariant = permanent) | `python_modules/dagster/dagster/_daemon/backfill.py:155` |
| Scheduler failure taxonomy | `UserCodeUnreachable\|CodeLocationLoadError → FAILURE` without incrementing `failure_count` (retry forever); else `failure_count+1, consecutive+1` | `python_modules/dagster/dagster/_scheduler/scheduler.py:702` |
| Run health taxonomy | `WorkerStatus{RUNNING,NOT_FOUND,FAILED,SUCCESS,UNKNOWN}`, `CheckRunHealthResult(transient)`; `RunFailureReason{STEP_FAILURE,…,START_TIMEOUT}`; monitoring resumes vs fails | `python_modules/dagster/dagster/_core/launcher/base.py:38` |
| Run conflict taxonomy | `DagsterRunAlreadyExists`, `DagsterRunConflict`, `DagsterRunNotFoundError` | `python_modules/dagster/dagster/_core/errors.py:562` |
| Op retry policy | `RetryPolicy(max_retries=1, delay, Backoff{LINEAR:attempt*delay, EXPONENTIAL:(2^attempt-1)*delay}, Jitter{FULL:random()*delay, PLUS_MINUS:delay±delay})`, `calculate_delay(attempt_num,…)` | `python_modules/dagster/dagster/_core/definitions/policy.py:35` |
| Op retry signal | `RetryRequested(max_retries=1, seconds_to_wait)`; `RetryRequestedFromPolicy` derived from policy | `python_modules/dagster/dagster/_core/definitions/events.py:735` |
| Op error boundary owns retries | `op_execution_error_boundary`: direct `RetryRequested` escalates; `Failure(allow_retries=False)` bypasses policy; else `RetryRequestedFromPolicy(calculate_delay(prev+1))`; interrupts also respect policy | `python_modules/dagster/dagster/_core/execution/plan/utils.py:54` |
| Step retry enforcement | On `RetryRequested`: if `retry_mode.disabled → STEP_FAILURE`; elif `prev>=max_retries → STEP_FAILURE (+MaxRetriesExceeded)` else `STEP_UP_FOR_RETRY(seconds_to_wait)`; counted via `RetryState` | `python_modules/dagster/dagster/_core/execution/plan/execute_plan.py:246` |
| Retry mode layering | `RetryMode{ENABLED,DISABLED,DEFERRED}` (+`for_inner_plan() ENABLED→DEFERRED` for orchestrator engines); `get_retries_config() Selector{enabled,disabled}` | `python_modules/dagster/dagster/_core/execution/retries.py:30` |
| Run-level retry budget | `auto_reexecution_should_retry_run`: must be FAILURE; max from `dagster/max_retries` tag else `instance.run_retries_max_retries`; budget `len(run_group)-1 <= max_retries`; documents parallel-manual extra-retry race | `python_modules/dagster/dagster/_core/execution/retries.py:83` |
| Auto-reexecution daemon | `should_retry()` via `dagster/will_retry`; reuse via `AUTO_RETRY_RUN_ID_TAG` or `parent+dagster/retry_number`; `retry_run()` creates `FROM_FAILURE` with `RETRY_NUMBER=len(group)`; safe to call multiple times | `python_modules/dagster/dagster/_daemon/auto_run_reexecution/auto_run_reexecution.py:106` |
| Dequeue retry budget | Only `UserCodeUnreachable\|CodeLocationLoadError` re-queues (`location_timeouts[loc]=now+retry_delay` + `PIPELINE_ENQUEUED`); budget `len(enqueue_events)-1 >= max_user_code_failure_retries → report_run_failed`; else unrecoverable → fail immediately | `python_modules/dagster/dagster/_daemon/run_coordinator/queued_run_coordinator_daemon.py:426` |
| Tick retry budget | Redo `STARTED` or `FAILURE and failure_count<=max_tick_retries`; sensor `MAX_FAILURE_RESUBMISSION_RETRIES=1`; asset `auto_materialize_max_tick_retries=3`; backfill `failure_count<max` only if REQUESTED+retryable | `python_modules/dagster/dagster/_scheduler/scheduler.py:572` |
| Infra backoff (no jitter) | `exponential_delay_generator(base 0.1×2.0)`, `BACKOFF_MAX_RETRIES=4`, `backoff()/async_backoff()` explicitly "doesn't implement any jitter" | `python_modules/dagster/dagster/_utils/backoff.py:11` |
| Daemon jitter | `IntervalDaemon(interval+random.uniform(0,jitter))`; scheduler checkpoint `random.randint(0,JITTER)` hourly | `python_modules/dagster/dagster/_daemon/daemon.py:249` |
| Run-health resume budget | `count_resume_run_attempts()` counts ENGINE_EVENT `RESUME_RUN`; budget `num<max_resume_run_attempts → resume_run(attempt+1)` else `report_run_failed`; timeouts `start_timeout/cancel_timeout/max_runtime` | `python_modules/dagster/dagster/_daemon/monitoring/run_monitoring.py:108` |
| Scheduler idempotency key | `_get_existing_run_for_request`: match on `tags_for_schedule + scheduled_execution_time + run_key`, repo-selector filter; `_submit_run_request`: if existing and `status!=NOT_STARTED` return existing (crash before SUCCESS) else reuse NOT_STARTED | `python_modules/dagster/dagster/_scheduler/scheduler.py:966` |
| Scheduler create+submit split | `_create_scheduler_run` tags `RUN_KEY_TAG`/`SCHEDULED_EXECUTION_TIME_TAG`; `submit_run()` wrapped so "created successfully but failed to launch" is retryable | `python_modules/dagster/dagster/_scheduler/scheduler.py:1004` |
| Sensor idempotency key | `fetch_existing_runs()` per-run_key serial `get_runs(tags={run_key})` (join perf note), filter by `SENSOR_NAME_TAG`+repo selector; `_get_or_create_sensor_run`: `status!=NOT_STARTED → SkippedSensorRun` (crash before tick update), same-tick dedup map | `python_modules/dagster/dagster/_daemon/sensor.py:1292` |
| Dequeue at-least-once guard | Re-fetch `get_run_by_id`, skip if `!=QUEUED`; emit `PIPELINE_STARTING`; re-fetch, ignore if advanced past `QUEUED/STARTING` before `launch_run` | `python_modules/dagster/dagster/_daemon/run_coordinator/queued_run_coordinator_daemon.py:372` |
| Duplicate-run tolerance | `gen_ignore_duplicate_run_worker()` — ignore duplicate run started elsewhere | `python_modules/dagster/dagster/_core/execution/api.py:79` |
| Crash test: scheduler | Crash matrix (`RUN_CREATED`/`TICK_CREATED`/`RUN_ADDED`/`TICK_SUCCESS`, SIGKILL/SIGTERM/exception): asserts single tick STARTED + single run, re-run resumes existing run to SUCCESS; documents unfixed race if crash lands between launch and executor STARTED | `python_modules/dagster/dagster_tests/scheduler_tests/test_scheduler_failure_recovery.py:155` |
| Crash test: sensor | Crash matrix (`RUN_CREATED`/`RUN_LAUNCHED`): asserts run_key `only_once` resumption, single tick stays STARTED then SUCCESS, next tick SKIPPED with no new runs | `python_modules/dagster/dagster_tests/daemon_sensor_tests/test_sensor_failure_recovery.py:160` |
| Ordering test | Queued-coordinator tests with `PRIORITY_TAG` `-1/3/foobar/-100/100` assert priority launch order; FIFO-within-priority via stable sort | `python_modules/dagster/dagster_tests/daemon_tests/test_queued_run_coordinator_daemon.py:320` |
| Dequeue retry test | `max_user_code_failure_retries=2/1` parametrized (threads on/off): "gives up after N retries", good runs from same location still dequeue | `python_modules/dagster/dagster_tests/daemon_tests/test_queued_run_coordinator_daemon.py:851` |

## Answers to Dimension Questions

1. **Can the same logical step execute more than once? — Yes, explicitly at-least-once.** Scheduler `_submit_run_request` (`python_modules/dagster/dagster/_scheduler/scheduler.py:783`), sensor `_get_or_create_sensor_run` (`python_modules/dagster/dagster/_daemon/sensor.py:1339`), and dequeue `_dequeue_run` (`python_modules/dagster/dagster/_daemon/run_coordinator/queued_run_coordinator_daemon.py:363`) all re-fetch and resume rather than assume exactly-once. Op steps re-execute via `STEP_UP_FOR_RETRY` (`python_modules/dagster/dagster/_core/execution/plan/execute_plan.py:286`); runs re-execute via auto-reexecution (`python_modules/dagster/dagster/_daemon/auto_run_reexecution/auto_run_reexecution.py:232`); run workers resume via monitoring (`python_modules/dagster/dagster/_daemon/monitoring/run_monitoring.py:134`). Idempotency is pushed to the caller: schedule/sensor dedup keys and op-level idempotent logic.
2. **Is retry policy centralized or stacked across layers? — Stacked across five layers, each owning its budget.** Op (`RetryPolicy` in `python_modules/dagster/dagster/_core/definitions/policy.py:35`, enforced in `plan/utils.py:54` + `plan/execute_plan.py:246`), run auto-reexecution (`retries.py:83`, daemon in `auto_run_reexecution.py:106`), dequeue re-queue (`queued_run_coordinator_daemon.py:426`), tick redo (`scheduler.py:572`, sensor `sensor.py:81`, asset-daemon, backfill `backfill.py:206`), run-health resume (`run_monitoring.py:108`). No single retry manager; `RetryMode.DEFERRED` (`retries.py:30`) is the only cross-layer coordination (inner plans defer to orchestrator engine).
3. **How does the scheduler recover capacity after crashes? — By recomputing from DB state each iteration; no leased capacity to recover.** In-progress set is re-queried (`_get_in_progress_run_records`, `queued_run_coordinator_daemon.py:342`, `IN_PROGRESS_RUN_STATUSES`); `max_runs_to_launch = max - len(in_progress)` (`queued_run_coordinator_daemon.py:210`) can even go negative (launched-bypassing-queue) and just yields empty batch. Ticks in `STARTED` or `FAILURE with failure_count<=max_tick_retries` are redone from persisted `last_iteration_timestamp` (`scheduler.py:572`); hourly checkpoint `_write_and_get_next_checkpoint_timestamp` (`scheduler.py:1077`) bounds re-scan. Location-level pause (`_location_timeouts`, `queued_run_coordinator_daemon.py:430`) expires by wall clock, so a crash clears the in-memory pause (fail-open). Heartbeat duplicate-daemon detection (`daemon.py:176`) prevents two controllers from double-spending capacity.
4. **Are scheduled and started events semantically distinct? — Yes.** `NOT_STARTED` (created, idempotency key exists, not yet submitted) vs `QUEUED` (`PIPELINE_ENQUEUED`) vs `STARTING` (`PIPELINE_STARTING` emitted pre-launch in `queued_run_coordinator_daemon.py:400`) vs `STARTED` (executor-owned) vs tick `STARTED` (claim, `instigation.py:306`, `end_timestamp=None`) vs tick `SUCCESS/FAILURE/SKIPPED` (terminal). Transitions are guarded: submit only if `status != FAILURE` (`scheduler.py:829`); dequeue skips if `!=QUEUED` (`queued_run_coordinator_daemon.py:380`) and aborts if advanced past `QUEUED/STARTING` (`queued_run_coordinator_daemon.py:420`); sensor returns `SkippedSensorRun` if existing run already left `NOT_STARTED` (`sensor.py:1355`).
5. **Can priority/fairness decisions be tested deterministically? — Yes for priority; fairness is probabilistic by design.** `_priority_sort` is a pure function (`int(tag)||0`, stable descending sort, `queued_run_coordinator_daemon.py:345`) covered by deterministic tests with tags `-100/-1/3/100/foobar` (`test_queued_run_coordinator_daemon.py:320`). Tag/pool limit blocking is deterministic given fixtures. Cross-location fairness via `shuffled_round_robin_by_key` (`_daemon/utils.py:17`) intentionally shuffles, so only statistical starvation-freedom holds, not deterministic order.
6. **Crash after external work succeeds but before acknowledgement — is retry expected, and is the side effect safe? — Retry is expected; safety depends on the layer.** Scheduler/sensor crash after run creation but before tick SUCCESS re-runs the tick and reuses the existing run via idempotency key (`scheduler.py:783`, `sensor.py:1355`) — safe (no second run). Op retry after a side effect re-executes user code (`execute_plan.py:286`) — safe only if the op is idempotent, which Dagster requires but does not enforce. The one unsafe window is scheduler crash after `launch_run` succeeds but before the executor flips `NOT_STARTED→STARTED`: the recovery test comment states the same run could be launched twice and only works around it with `wait_for_all_runs_to_start` (`test_scheduler_failure_recovery.py:189`).

## Architectural Decisions

- **Tick as the unit of scheduling record** (`TickStatus` + `TickData.failure_count/consecutive_failure_count`, `python_modules/dagster/dagster/_core/scheduler/instigation.py:305`) — separates "claimed work" (STARTED) from outcome (SUCCESS/FAILURE/SKIPPED) and enables tick-granular redo budgets distinct from run-granular retries.
- **User-code-unreachable is never a tick failure** (`scheduler.py:702`, `backfill.py:155`, dequeue `queued_run_coordinator_daemon.py:426`) — treated as infrastructure outage: pause location / retry forever rather than consume failure budget. Framework `DagsterError/CheckError` is the opposite: permanent, non-retryable.
- **Idempotency keys in tags, not a separate dedup table** (`RUN_KEY_TAG`, `SCHEDULED_EXECUTION_TIME_TAG`, `SENSOR_NAME_TAG` in `python_modules/dagster/dagster/_core/storage/tags.py:52`) — makes dedup queryable via existing `get_runs(RunsFilter(tags=…))` at the cost of serial per-key fetches (perf comment in `sensor.py:1303`).
- **Create-then-submit split** (`_create_scheduler_run` + `submit_run`, `scheduler.py:1004`/`scheduler.py:831`) — creates a `NOT_STARTED` anchor row so a crash before acknowledgement is resumable without duplicating the logical request.
- **Layered budgets instead of a global retry policy** — each layer (op delay/backoff/jitter, run-group count, dequeue enqueue-event count, tick failure_count, resume ENGINE_EVENT count) counts with its own observable signal, avoiding cross-layer clock coupling but allowing budget stacking.

## Notable Patterns

- **Generator-as-daemon contract** (`DaemonIterator`, `daemon.py:60`): infinite generator that must periodically `yield` for heartbeat; errors yielded into a bounded ring and surfaced via heartbeat instead of crashing the process.
- **Re-fetch before mutate** (dequeue `queued_run_coordinator_daemon.py:373`/`queued_run_coordinator_daemon.py:407`/`queued_run_coordinator_daemon.py:418`, monitoring `run_monitoring.py:125`): every state transition re-reads the run row to tolerate concurrent daemon/duplicate-daemon/operator action — the core at-least-once guard.
- **Declarative delay modulation** (`RetryPolicy.calculate_delay`, `policy.py:87`): backoff × jitter computed without running user code, so `seconds_to_wait` is known at `RetryRequested` time and testable.
- **Stable-sort priority with FIFO preservation** (`queued_run_coordinator_daemon.py:353` comment "sorted is stable, so fifo is maintained") — minimal mechanism achieving priority + fairness within a priority band.
- **Crash-injection via debug flags** (`check_for_debug_crash(…,"RUN_CREATED"/"RUN_ADDED"/"TICK_SUCCESS")`, `scheduler.py:825`): deterministic at-least-once tests by killing the scheduler at each acknowledgement boundary.

## Tradeoffs

- **Stacked retries can multiply**: an op with `max_retries=3` inside a run with `max_retries=3` inside a tick redo can execute the op up to 4×4 times; run-group counting even documents an "extra retry" when manual and automatic retries race (`retries.py:101`). Explicit and bounded per layer, but no global ceiling.
- **Fail-open location pause**: `_location_timeouts` is in-memory only, so a daemon crash clears backoff for a sick code location and the next iteration hammers it again (vs. persisting circuit-breaker state). Trades availability (self-healing after crash) against thundering-herd on a flapping code server.
- **Serial per-run-key dedup queries** (`sensor.py:1306`): chosen because the planner handles `runs/run_tags` joins poorly with `IN` clauses — trades N round trips for predictable latency; noted as intentional perf workaround.
- **Infra `backoff()` has no jitter** (self-documented, `backoff.py:31`) while op `RetryPolicy` and daemon intervals do — fine for single-threaded RPC retry, risky if many daemons retry the same store simultaneously.
- **Orphan-state grace (12h) vs immediate delete for DECLARED_IN_CODE** (`scheduler.py:356`): trades stale-tick accumulation against accidental mass-cancellation during code-location outages (error locations skipped entirely).

## Failure Modes / Edge Cases

- **Double-launch race** (acknowledged, `test_scheduler_failure_recovery.py:189`): crash between external launch success and executor `STARTED` flip can launch the same run twice; mitigated in tests by waiting, not fixed in implementation.
- **Parallel manual + auto retry overrun** (`retries.py:101`): run-group size check is read-then-act without a lock; one extra automatic retry beyond `max_retries` is accepted as unlikely.
- **Duplicate daemon processes** (detected via heartbeat `daemon_id` mismatch, `daemon.py:176`): only logged, not fenced — two schedulers can still double-submit until an operator kills one; idempotency keys make this safe for schedules/sensors but op side effects may still duplicate.
- **Counter-init failure fails open** (`queued_run_coordinator_daemon.py:293`): if pool/concurrency keys cannot load, all op-concurrency blocking is skipped for the iteration — capacity guard silently disabled rather than stalling the queue.
- **Invalid `dagster/max_retries` tag** (`retries.py:135`): warns and disables auto-retry for that run instead of falling back to instance default — a typo silently removes the safety net.
- **Non-partitioned catch-up dropped** (`scheduler.py:637`): more than one missed tick collapses to the latest — trades exactly-once-per-cron-fire for bounded backlog; partitioned schedules instead clamp to `max_catchup_runs`.

## Future Considerations

- Fence (not just log) duplicate daemon IDs via a DB lease so the double-launch race and double-spend of `max_concurrent_runs` are structurally impossible.
- Close the launch-vs-STARTED race by flipping run status inside the same transaction as launch submission, or by making `launch_run` itself idempotent on run ID.
- Add jitter option to `dagster._utils.backoff.backoff` (currently documented as absent) for multi-daemon store contention paths.
- Consider a global retry-ceiling or retry-budget dashboard spanning op/run/tick/dequeue layers, since stacked budgets are individually bounded but jointly multiplicative.
- Persist code-location circuit-breaker (`_location_timeouts`) so crash recovery does not immediately re-hammer an unreachable code server.

## Questions / Gaps

- No evidence found for throttle-specific (429/rate-limit) classification distinct from unreachable: searched `scheduler.py`, `backfill.py:155`, `queued_run_coordinator_daemon.py:426`, `launcher/base.py:38` — only reachable/unreachable, permanent/transient, conflict (`RunAlreadyExists` treated as retryable in backfill but as dedup-hit in scheduler) are classified.
- No evidence found for deterministic tests of `shuffled_round_robin_by_key` fairness bounds — only priority-order tests located; the shuffle path appears intentionally non-deterministic.
- Exact default of scheduler `max_tick_retries` not resolved in this pass (threaded through `daemon.py:291` from scheduler settings); tick-redo predicate (`scheduler.py:572`) confirmed but numeric default lives in scheduler-instance config outside the sampled files.

---

Generated by `Dimension 01.10: Scheduler Controller, Retries, and At-Least-Once Semantics` against `dagster`.
