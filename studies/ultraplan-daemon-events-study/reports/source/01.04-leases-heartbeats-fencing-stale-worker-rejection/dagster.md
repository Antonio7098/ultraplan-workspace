# Source Analysis: dagster

## Leases, Heartbeats, Fencing, and Stale-Worker Rejection

### Source Info

| Field | Value |
|-------|-------|
| Name | dagster |
| Path | `studies/ultraplan-daemon-events-study/sources/dagster` |
| Language / Stack | Python / SQLAlchemy (Dagster OSS) |
| Analyzed | 2026-09-02 |

## Summary

Dagster's daemon subsystem implements **observability-only heartbeats**, not ownership leases. Each daemon process generates a random UUID (`python_modules/dagster/dagster/_daemon/controller.py:156`) and periodically overwrites a single row per `daemon_type` in `daemon_heartbeats` (`python_modules/dagster/dagster/_core/storage/runs/schema.py:94`). There is no lease table, TTL enforcement, generation/epoch, CAS, or conditional-write fencing on any authoritative mutation (ticks, runs, sensor cursor, backfills). A stale worker that regains CPU after its replacement has started can unconditionally overwrite state. The only staleness mitigation is a 24-hour heuristic that converts abandoned `STARTED` ticks to `SKIPPED` and an in-memory per-process guard against concurrent evaluation of the same sensor.

## Rating

**2 / 10 — Absent / unsafe.** Heartbeats exist for health display; lease acquisition, renewal, fencing tokens, and stale-writer rejection are absent. The system explicitly warns that multiple daemons are "not supported" but does not enforce single ownership.

*Rationale:* Per rubric 1-3 = "Absent, implicit, ad-hoc, or unsafe." No lease/lock schema, no heartbeat expiry gate, no generation field, no conditional update, no fencing on product-state writes, single integration-free warning at `python_modules/dagster/dagster/_daemon/daemon.py:182`. Failure-safe answer to "Pause an old worker until its lease is stolen, then let it finish. Can it corrupt state? — Yes."

## Evidence Collected

Every entry MUST include a file path with line numbers. Format: `path/to/file.ts:NN`.

| Area | Evidence | File:Line |
|------|----------|-----------|
| Lease/lock schema | No lease or lock table exists. Only `daemon_heartbeats` with `unique daemon_type` and no holder-TTL columns. | `python_modules/dagster/dagster/_core/storage/runs/schema.py:94-107` |
| Heartbeat schema & storage | `DaemonHeartbeat(timestamp, daemon_type, daemon_id, errors)` NamedTuple; serialized to `body` column. Table enforces one row per `daemon_type`. | `python_modules/dagster/dagster/_daemon/types.py:25-52`, `python_modules/dagster/dagster/_core/storage/runs/schema.py:94-107` |
| Heartbeat write path | `instance.add_daemon_heartbeat(DaemonHeartbeat(curr_time.timestamp(), daemon_type, daemon_uuid))` — last-write-wins insert-on-conflict fallback to update `WHERE daemon_type =` with no CAS. | `python_modules/dagster/dagster/_daemon/daemon.py:194-201`, `python_modules/dagster/dagster/_core/storage/runs/sql_run_storage.py:852-873` |
| Heartbeat cadence & expiry rules | `DEFAULT_HEARTBEAT_INTERVAL_SECONDS = 30`, `DEFAULT_DAEMON_HEARTBEAT_TOLERANCE_SECONDS = 1800` (30 min). Health = `curr_time <= heartbeat_timestamp + interval + tolerance` plus null-errors check. Only used for display/health RPC. | `python_modules/dagster/dagster/_daemon/controller.py:36-45`, `python_modules/dagster/dagster/_daemon/controller.py:452-468` |
| Heartbeat interval check inside daemon thread | `_check_add_heartbeat` throttles on local `self._last_heartbeat_time` and optionally skips when `daemon_skip_heartbeats_without_errors` and no errors. | `python_modules/dagster/dagster/_daemon/daemon.py:147-192` |
| Duplicate-daemon detection | Logs `Another %s daemon is still sending heartbeats ... Last heartbeat daemon id: %s, Current daemon_id: %s` but continues writing; no abort, no lease revocation. | `python_modules/dagster/dagster/_daemon/daemon.py:176-191` |
| Process identity | Per-process UUID `self._daemon_uuid = str(uuid.uuid4())` shared across all daemon threads; stored as `daemon_id` (not PID). Distinguishes restarts but never consulted for fencing. | `python_modules/dagster/dagster/_daemon/controller.py:156`, `python_modules/dagster/dagster/_daemon/types.py:32` |
| Tick model — no generation/epoch | `TickData`/`InstigatorTick` carry `tick_id` (auto-increment PK), `status`, `timestamp`, `run_ids`, `cursor`, `failure_count` etc. No `generation`, `epoch`, `owner_id`, `lease_expiry`. | `python_modules/dagster/dagster/_core/scheduler/instigation.py:590-720` |
| Tick create path — unconditional | `create_tick` inserts with no ownership check; on PK conflict raises generic invariant. | `python_modules/dagster/dagster/_core/storage/schedules/sql_schedule_storage.py:394-417` |
| Tick update path — no CAS/fencing | `update_tick` does `UPDATE ... WHERE id == tick_id SET status/body/timestamp` with no version/owner predicate, no `WHERE version = expected`. | `python_modules/dagster/dagster/_core/storage/schedules/sql_schedule_storage.py:419-436` |
| Sensor tick resume heuristic | `MAX_TIME_TO_RESUME_TICK_SECONDS = 24h`; if prior `STARTED` tick older than 24h or has no unsubmitted runs it is forced to `SKIPPED` via `instance.update_tick(most_recent_tick.with_status(SKIPPED))`; otherwise it is resumed. Pure time heuristic, not lease. | `python_modules/dagster/dagster/_daemon/sensor.py:74-77`, `python_modules/dagster/dagster/_daemon/sensor.py:538-555` |
| Sensor state cursor write — no fencing | `SensorLaunchContext._write` unconditionally `instance.update_tick(self._tick)` then `instance.update_instigator_state(state.with_data(SensorInstigatorData(... cursor ...)))` — late worker can overwrite `cursor` and `last_tick_success_timestamp`. | `python_modules/dagster/dagster/_daemon/sensor.py:232-282` |
| Sensor thread deduplication — in-memory only | `if selector_id in sensor_tick_futures and not future.done(): continue` prevents duplicate ticks *within single process*; not distributed. | `python_modules/dagster/dagster/_daemon/sensor.py:474-492` |
| Concurrency-slot fencing (unrelated) | `_claim_concurrency_slot` uses `SELECT ... FOR UPDATE SKIP LOCKED` semaphore on `concurrency_slots` — fences step execution, not daemon ownership. | `python_modules/dagster/dagster/_core/storage/event_log/sql_event_log.py:2772-2794` |
| Health check loop — warning only | `check_daemon_heartbeats` collects `is_daemon_healthy==False` daemons and emits `logger.warning(message)`. `check_daemon_threads` only checks local thread liveness; `_daemon_heartbeat_health` falls back to `last_healthy_heartbeat_times` local dict on exception. | `python_modules/dagster/dagster/_daemon/controller.py:261-273`, `python_modules/dagster/dagster/_daemon/controller.py:213-245` |
| Run-key idempotence (partial late-result mitigation) | `fetch_existing_runs` + `fetch_existing_runs_by_key` + `_get_or_create_sensor_run` deduplicate by `RUN_KEY_TAG` + repository selector; `SkippedSensorRun` returned if run already `!= NOT_STARTED`. Not fencing — merely avoids duplicate `NOT_STARTED` runs. | `python_modules/dagster/dagster/_daemon/sensor.py:1292-1336`, `python_modules/dagster/dagster/_daemon/sensor.py:1339-1375` |
| Tests for heartbeats | Storage heartbeat tests only verify insert/get/wipe and serialization backcompat, not concurrent writers or stale rejection. No test for stale worker corrupting tick. | `python_modules/dagster/dagster_tests/storage_tests/utils/run_storage.py:1328-1370`, `python_modules/dagster/dagster_tests/daemon_tests/test_types.py:1-27` |
| Daemon skip-heartbeats optimization | `daemon_skip_heartbeats_without_errors == False` by default; comment notes cloud enables it to reduce DB writes — confirms heartbeats are optional telemetry there. | `python_modules/dagster/dagster/_core/instance/methods/daemon_methods.py:122-126` |

## Answers to Dimension Questions

### 1. Can two workers believe they own the same work?

**Yes.** Two `DagsterDaemonController` processes can start simultaneously, each with its own `_daemon_uuid` (`python_modules/dagster/dagster/_daemon/controller.py:156`). Each spawns one thread per daemon type (`python_modules/dagster/dagster/_daemon/controller.py:189-202`) sharing that UUID. Both threads call `_check_add_heartbeat` (`python_modules/dagster/dagster/_daemon/daemon.py:105-111`) — the only cross-process check is reading `get_daemon_heartbeats()[daemon_type]` (`python_modules/dagster/dagster/_daemon/daemon.py:176`) and logging if `daemon_id` differs (`python_modules/dagster/dagster/_daemon/daemon.py:182-190`). No lock is acquired, no lease row is inserted with `IF NOT EXISTS`, and the write proceeds regardless (`python_modules/dagster/dagster/_daemon/daemon.py:194-201`). For ticks, `create_tick`/`update_tick` (`python_modules/dagster/dagster/_core/storage/schedules/sql_schedule_storage.py:394-436`) have no owner predicate; any daemon can create or mutate any sensor's tick. The in-memory `sensor_tick_futures` guard (`python_modules/dagster/dagster/_daemon/sensor.py:474-478`) only prevents one sensor tick at a time *within a single process*.

### 2. Does lease expiry alone prevent stale writes?

**There is no lease, so no.** Health expiry (`healthy = curr_time <= timestamp + 30 + 1800`, `python_modules/dagster/dagster/_daemon/controller.py:455-458`) is computed only for `get_daemon_statuses`/`all_daemons_healthy` and the controller's warning loop (`python_modules/dagster/dagster/_daemon/controller.py:261-273`). It is never consulted before `update_tick` (`python_modules/dagster/dagster/_core/storage/schedules/sql_schedule_storage.py:419`), `update_instigator_state` (`python_modules/dagster/dagster/_daemon/sensor.py:267`), or `create_run` (`python_modules/dagster/dagster/_daemon/sensor.py:1419`). A paused worker whose heartbeat has aged out can still successfully call those storage methods on resume because storage does not validate the caller's liveness.

### 3. Is fencing checked by product-state and artifact mutations, not only event writes?

**No fencing anywhere.** Surveyed mutations:
- Schedule/sensor tick updates: unconditional `UPDATE WHERE id` (`python_modules/dagster/dagster/_core/storage/schedules/sql_schedule_storage.py:432`).
- Sensor cursor / `SensorInstigatorData` write: unconditional `update_instigator_state` (`python_modules/dagster/dagster/_daemon/sensor.py:267`).
- Run creation: `instance.create_run` (`python_modules/dagster/dagster/_daemon/sensor.py:1419`) checks run-key dedup (`python_modules/dagster/dagster/_daemon/sensor.py:1351-1364`) but not daemon generation.
- Backfill/asset-daemon evaluations: `add_auto_materialize_asset_evaluations` (`python_modules/dagster/dagster/_core/storage/schedules/sql_schedule_storage.py:478-514`) upserts on `(evaluation_id, asset_key)` with no owner.
- Daemon heartbeats themselves: last-write-wins (`python_modules/dagster/dagster/_core/storage/runs/sql_run_storage.py:852-873`).
Event writes via `report_runless_asset_event` (`python_modules/dagster/dagster/_daemon/sensor.py:866`) similarly have no fencing token. No file contains `generation`, `epoch`, `fencing_token`, or `compare_and_swap` / `CAS` in a daemon ownership context (verified via grep of `dagster/_daemon` and `dagster/_core/storage`).

### 4. How is process identity distinguished from reusable PIDs?

Via a random UUID, not a PID. `DagsterDaemonController.__init__` sets `self._daemon_uuid = str(uuid.uuid4())` (`python_modules/dagster/dagster/_daemon/controller.py:156`) and passes it to every `daemon.run_daemon_loop` thread (`python_modules/dagster/dagster/_daemon/controller.py:193`). `DaemonHeartbeat.daemon_id` (`python_modules/dagster/dagster/_daemon/types.py:32`) persists this UUID in the `daemon_id` column (`python_modules/dagster/dagster/_core/storage/runs/schema.py:104`) and `INSERT/UPDATE` body (`python_modules/dagster/dagster/_core/storage/runs/sql_run_storage.py:860`). This correctly distinguishes process incarnations (unlike reusable OS PIDs) but is only used for the diagnostic log at `python_modules/dagster/dagster/_daemon/daemon.py:182-190`; it is not a fencing token checked on writes. No `pid`, `hostname`, or `epoch` column exists.

### 5. What happens to a late completion from a superseded attempt?

**It can corrupt authoritative state.** Concrete path for sensor:

1. Old daemon creates tick N `STARTED` (`python_modules/dagster/dagster/_daemon/sensor.py:572-582`), reserves `reserved_run_ids` (`python_modules/dagster/dagster/_daemon/sensor.py:1105-1111`), then pauses.
2. New daemon starts, observes tick N `STARTED` and if `now - timestamp > 86400` moves it to `SKIPPED` (`python_modules/dagster/dagster/_daemon/sensor.py:553-555`) or retries on `FAILURE` (`python_modules/dagster/dagster/_daemon/sensor.py:562-569`). Otherwise it may resume the same `tick_id`. In any case, new daemon may create tick N+1 with correct cursor.
3. Old daemon resumes, finishes evaluation, and in `SensorLaunchContext.__exit__` calls `_write` (`python_modules/dagster/dagster/_daemon/sensor.py:332-333`) which does `update_tick` (`python_modules/dagster/dagster/_daemon/sensor.py:237`) setting `SUCCESS` and then `update_instigator_state` with `cursor = self._tick.cursor` (`python_modules/dagster/dagster/_daemon/sensor.py:260-262`) and `last_tick_success_timestamp = now` (`python_modules/dagster/dagster/_daemon/sensor.py:278`). This overwrites the newer cursor/timestamp.

Mitigations that partially help but do not fence:
- `run_key` dedup in `fetch_existing_runs` (`python_modules/dagster/dagster/_daemon/sensor.py:1292`) prevents duplicate launched runs; late `create_run` with same `run_key` returns `SkippedSensorRun` (`python_modules/dagster/dagster/_daemon/sensor.py:1358`).
- `with_status` on tick moves `STARTED→SKIPPED` for abandoned ticks (`python_modules/dagster/dagster/_daemon/sensor.py:554`), but a tick resumed before 24h will be double-submitted.

No test covers "stale worker resumes after lease stolen" — grep for `stale`, `superseded`, `generation`, `fence` in daemon tests returns no fencing tests.

## Architectural Decisions

| Decision | Evidence | Consequence |
|----------|----------|-------------|
| Heartbeats as health telemetry, not leases | `DaemonHeartbeatsTable` unique on `daemon_type` only (`python_modules/dagster/dagster/_core/storage/runs/schema.py:103`); write is UPSERT without CAS (`python_modules/dagster/dagster/_core/storage/runs/sql_run_storage.py:852-873`) | No mutual exclusion; multiple daemons silently race. |
| Single-daemon deployment assumption | Log string "You likely have multiple daemon processes ... which is not supported" (`python_modules/dagster/dagster/_daemon/daemon.py:183-184`); controller spawns threads not processes (`python_modules/dagster/dagster/_daemon/controller.py:189-202`) | Safety relies on operations, not code. |
| Tick is the unit of work, with time-based resume window | `MAX_TIME_TO_RESUME_TICK_SECONDS = 86400` (`python_modules/dagster/dagster/_daemon/sensor.py:77`); `_get_evaluation_tick` resume vs `SKIPPED` branches (`python_modules/dagster/dagster/_daemon/sensor.py:538-555`) | Late worker within 24h can double-submit `reserved_run_ids`. |
| In-memory thread deduplication instead of distributed lock | `sensor_tick_futures[selector_id]` futures dict (`python_modules/dagster/dagster/_daemon/sensor.py:361-492`) | Only guards within one OS process. |
| UUID for daemon identity, but not as fencing token | `DagsterDaemonController._daemon_uuid` random UUID (`python_modules/dagster/dagster/_daemon/controller.py:156`), stored as `daemon_id` (`python_modules/dagster/dagster/_daemon/types.py:32`) but only logged (`python_modules/dagster/dagster/_daemon/daemon.py:182-190`) | Solves PID reuse correctly, but misses opportunity to fence. |
| `run_key` idempotence for sensor runs | `fetch_existing_runs` keyed by `RUN_KEY_TAG` + repo selector (`python_modules/dagster/dagster/_daemon/sensor.py:1292-1336`) + `_get_or_create_sensor_run` early-return (`python_modules/dagster/dagster/_daemon/sensor.py:1354-1364`) | Prevents duplicate NOT_STARTED runs; does not prevent stale cursor tick mutation. |
| Configurable skip-heartbeats optimization | `daemon_skip_heartbeats_without_errors` default `False` (`python_modules/dagster/dagster/_core/instance/methods/daemon_methods.py:122-126`) | In Dagster Cloud heartbeats may not even be written, reinforcing they are not leases. |

## Notable Patterns

- **Warning-only multi-writer detection:** Heartbeat mismatch emits `logger.error` but proceeds — classic missing-enforcement anti-pattern (`python_modules/dagster/dagster/_daemon/daemon.py:176-201`).
- **Last-write-wins tick update:** No `version`/`if-match` header; compare to patterns that use `UPDATE ... WHERE version = :expected RETURNING` (`python_modules/dagster/dagster/_core/storage/schedules/sql_schedule_storage.py:419-436`).
- **Span-based liveness	yield:** Daemon loop yields `SpanMarker.START_SPAN/END_SPAN` so controller can interleave heartbeat checks (`python_modules/dagster/dagster/_daemon/daemon.py:98-145`, `python_modules/dagster/dagster/_daemon/sensor.py:368-395`).
- **Two-phase run submission with reservation:** `reserved_run_ids` written to `tick_data` via `set_run_requests` (`python_modules/dagster/dagster/_daemon/sensor.py:1109-1111`) before `submit_run` (`python_modules/dagster/dagster/_daemon/sensor.py:755`); survivors resume via `unsubmitted_run_ids_with_requests` (`python_modules/dagster/dagster/_core/scheduler/instigation.py:560-568`). Provides at-least-once run idempotence, not exactly-once fencing.
- **Monotonic `failure_count` / `consecutive_failure_count` on ticks:** Tracked in `TickData` (`python_modules/dagster/dagster/_core/scheduler/instigation.py:604-616`) and reset on `SKIPPED/SUCCESS` (`python_modules/dagster/dagster/_daemon/sensor.py:160-162`).

## Tradeoffs

| Tradeoff | Pro | Con |
|----------|-----|-----|
| Single-row per daemon_type, no lease table | Trivial schema, tiny DB load, works with any SQL backend | Zero safety under split-brain; requires external orchestration to ensure single daemon. |
| Long 30-min tolerance (`python_modules/dagster/dagster/_daemon/controller.py:36`) | Avoids flapping health status under GC / DB contention | Stale daemon considered healthy for 30 min; no fast failover. |
| 24h resume window (`python_modules/dagster/dagster/_daemon/sensor.py:77`) | Recovers from short daemon outage without losing work | Old `STARTED` tick <24h resumed concurrently by new daemon can duplicate runs. |
| In-memory futures dedup (`python_modules/dagster/dagster/_daemon/sensor.py:475`) | No DB lock contention, low latency | No cross-process mutual exclusion; horizontal scaling unsafe. |
| Run-key idempotence vs fencing token | Cheap, user-visible (`run_key` tag), no schema change | Late tick can still corrupt cursor/state after dedup handles runs. |
| Heartbeat errors field aggregated over `DAEMON_HEARTBEAT_ERROR_LIMIT=5` (`python_modules/dagster/dagster/_daemon/daemon.py:42`) + `error_interval_seconds=300` (`python_modules/dagster/dagster/_daemon/controller.py:45`) | Surfaces recent failures in UI | Error window logic (`python_modules/dagster/dagster/_daemon/daemon.py:154-161`) can suppress heartbeats entirely when errors age out and `daemon_skip_heartbeats_without_errors`. |

## Failure Modes / Edge Cases

1. **Split-brain duplicate execution:** Two daemons both pass `is_under_min_interval` (`python_modules/dagster/dagster/_daemon/sensor.py:1259-1289`) due to stale `last_tick_start_timestamp` and concurrently create tick N. Both reserve different `reserved_run_ids` for same logical work, launching duplicate runs (mitigated only if `run_key` overlaps).
2. **Stale cursor overwrite:** Old worker resumes after new worker's successful tick and overwrites `SensorInstigatorData.cursor` (`python_modules/dagster/dagster/_daemon/sensor.py:260`) causing sensor to re-read already-consumed data and either duplicate or skip future events.
3. **Heartbeat DB partition:** If run storage is unavailable, `_check_add_heartbeat` swallows exception and logs (`python_modules/dagster/dagster/_daemon/daemon.py:112-116`); `get_daemon_heartbeats` may throw, causing `_daemon_heartbeat_health` fallback to `last_healthy_heartbeat_times` (`python_modules/dagster/dagster/_daemon/controller.py:239-245`) — health appears warm despite DB outage.
4. **Clock skew:** Heartbeat expiry compares `get_current_timestamp()` (wall clock) against stored `timestamp` (`python_modules/dagster/dagster/_daemon/controller.py:432-458`). No monotonic clock or bounded-drift handling; skew can mark healthy daemons unhealthy or vice versa.
5. **Daemon restart identity reuse:** On restart a new UUID is generated, immediately overwriting the prior row (`python_modules/dagster/dagster/_core/storage/runs/sql_run_storage.py:865-872`). There is no wait for prior lease TTL to expire, so any in-flight tick from prior incarnation that resumes can race with new incarnation's tick.
6. **Thread-pool starvation / hanging sensor:** `check_daemon_heartbeats` warns after 60s interval (`python_modules/dagster/dagster/_daemon/controller.py:49-320`) but does not kill or fence hung threads; `check_daemon_threads` only restarts when thread dies (`python_modules/dagster/dagster/_daemon/controller.py:247-259`), not when stuck.
7. **Missing stale-writer test coverage:** No daemon test asserts that an expired worker's `update_tick` or `update_instigator_state` is rejected; existing heartbeat tests (`python_modules/dagster/dagster_tests/storage_tests/utils/run_storage.py:1328` and `python_modules/dagster/dagster_tests/daemon_tests/test_types.py:1`) only verify serialization.

## Future Considerations

- **Lease table with holder+fencing token:** Add `daemon_leases(daemon_type PK, holder_id, fencing_token BIGSERIAL, expires_at)` with `INSERT ... ON CONFLICT ... WHERE expires_at < now()` and require `fencing_token` in `WHERE` clause of all `update_tick`/`update_instigator_state` statements (CAS pattern). Return `LeaseLost` on mismatch — standard fencing approach proven in Temporal/Chubby.
- **Generation/epoch on ticks and sensor state:** Add `generation` integer to `TickData` and `SensorInstigatorData` (`python_modules/dagster/dagster/_core/scheduler/instigation.py:590`), increment on each lease acquisition, and gate `update_tick` with `WHERE generation = :expected`.
- **Enforce single writer:** Change `DaemonHeartbeatsTable` to include lease expiry and make daemon loop abort on heartbeat mismatch rather than warn (`python_modules/dagster/dagster/_daemon/daemon.py:182`), or exit process for orchestrator to restart.
- **Monotonic expiry & bounded clock:** Use DB server time (`now()`) for lease expiry and enforce with DB transaction time, removing wall-clock skew sensitivity seen in `python_modules/dagster/dagster/_daemon/controller.py:455`.
- **Late-result shedding:** Make `SensorLaunchContext._write` (`python_modules/dagster/dagster/_daemon/sensor.py:232`) conditional on still holding the lease; on `LeaseLost` discard late cursor advancement and mark tick `SKIPPED`.
- **Observability:** Emit metrics/logs for lease steal, stale-write rejection, and fencing-token monotonicity; surface in `DaemonStatus` (`python_modules/dagster/dagster/_daemon/types.py:55`).

## Questions / Gaps

| Question | Search boundary | Finding |
|----------|-----------------|---------|
| Is there any hidden fencing in `DagsterInstance` wrappers? | Grepped `python_modules/dagster/dagster/_core/instance/**` for `fence`, `lease`, `generation`, `epoch`, `CAS`, `compare_and_swap` | No evidence found; only `daemon_methods.py` and `run_storage` heartbeat delegation. |
| Could `GlobalOpConcurrency` provide fencing? | Inspected `python_modules/dagster/dagster/_core/execution/plan/instance_concurrency_context.py:138` and `python_modules/dagster/dagster/_core/storage/event_log/sql_event_log.py:2541-2794` | Fences step slots with `FOR UPDATE SKIP LOCKED`, not daemon ownership. |
| Is there a distributed lock in `GrpcServerRegistry` TTL? | `python_modules/dagster/dagster/_daemon/controller.py:59` defines `DAEMON_GRPC_SERVER_HEARTBEAT_TTL=20` for code servers | TTL is for ephemeral gRPC servers, unrelated to daemon lease; not checked before tick mutation. |
| Are ticks protected by DB constraints? | `python_modules/dagster/dagster/_core/storage/schedules/schema.py` (tick schema) inspection | No unique constraint beyond PK; no owner FK to heartbeat. `python_modules/dagster/dagster/_core/storage/schedules/sql_schedule_storage.py:419` confirms unconditional update. |

---

Generated by `Dimension 01.04: Leases, Heartbeats, Fencing, and Stale-Worker Rejection` against `dagster`.
