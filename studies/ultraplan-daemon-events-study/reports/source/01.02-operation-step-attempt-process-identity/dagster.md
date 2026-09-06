# Source Analysis: dagster

## Operation, Step, Attempt, and Process Identity

### Source Info

| Field | Value |
|-------|-------|
| Name | dagster |
| Path | `studies/ultraplan-daemon-events-study/sources/dagster` |
| Language / Stack | Python (dagster core), SQLite/Postgres run storage, event log |
| Analyzed | 2026-09-02 |

## Summary

Dagster models durable intent as `DagsterRun` (persisted via `RunStorage` under a UUID `run_id`) with an accompanying `ExecutionPlanSnapshot` that checkpoints the intended step DAG and `KnownExecutionState`. The smallest checkpointable unit is the execution step (`ExecutionStep`, keyed by `step_key`/`StepHandle`), checkpointed via `EventLogEntry` stream (`STEP_START`/`STEP_SUCCESS`/`STEP_FAILURE`/`STEP_SKIPPED`/`STEP_OUTPUT`). Two retry layers exist: intra-run step retries driven by `RetryPolicy`/`RetryState` that reuse the same `run_id` and increment a per-step attempt counter, and inter-run re-executions that mint a new `run_id` with immutable `parent_run_id`/`root_run_id` linkage and a `PastExecutionState` chain. Plans and known state are serdes-serialized and can span multiple OS processes (in-process, multiprocess, k8s/celery executors) while all events remain unambiguously keyed by `(run_id, step_key, step_handle)`. The design is mature and test-covered but step-level retry lacks an external attempt fence token, leaving late-arriving attempt 1 outputs indistinguishable except by timestamp ordering.

## Rating

**8 / 10** — Clear, formally-typed hierarchy (run → plan → step → attempt) with persisted snapshots, explicit parent/root lineage, tag-correlated run groups, and process-spanning reconstruction. Intrusive tests exercise both step-policy retries and run-level resume-from-failure. Deductions for: (a) step retries reuse `run_id+step_key` without a monotonic attempt ID on the event envelope, so a late attempt-1 handler cannot be cryptographically rejected; (b) `PreviousRetryState` is in-memory until snapshotted and not itself durably keyed per attempt row — failure modes rely on event replay heuristics.

## Evidence Collected

Every entry MUST include a file path with line numbers. Format: `path/to/file.ts:NN`.

| Area | Evidence | File:Line |
|------|----------|-----------|
| Durable run model | `class DagsterRun` with fields `run_id`, `job_name`, `status`, `tags`, `root_run_id`, `parent_run_id`, `job_snapshot_id`, `execution_plan_snapshot_id`, `step_keys_to_execute` | `python_modules/dagster/dagster/_core/storage/dagster_run.py:236` |
| Run status enum | `QUEUED, NOT_STARTED, STARTING, STARTED, SUCCESS, FAILURE, CANCELING, CANCELED` + `IN_PROGRESS_RUN_STATUSES` | `python_modules/dagster/dagster/_core/storage/dagster_run.py:54` |
| Run parent/root invariant | `Must set both root_run_id and parent_run_id when creating a PipelineRun that belongs to a run group` and `root_run_id and parent_run_id` both-required check | `python_modules/dagster/dagster/_core/storage/dagster_run.py:311` and `python_modules/dagster/dagster/_core/instance/runs/run_domain.py:143` |
| Run group query | `get_run_group(run_id)` derives `root_run_id` then queries `RunTagsTable where key==ROOT_RUN_ID_TAG` to return `(root_run_id, [root_run, *group])` | `python_modules/dagster/dagster/_core/storage/runs/sql_run_storage.py:555` |
| Tag lineage | `PARENT_RUN_ID_TAG = "dagster/parent_run_id"` and `ROOT_RUN_ID_TAG = "dagster/root_run_id"` | `python_modules/dagster/dagster/_core/storage/tags.py:34` |
| Reexecution linkage | `root_run_id=parent_dagster_run.root_run_id or parent_dagster_run.run_id, parent_run_id=parent_dagster_run.run_id` on `create_run_for_job` in `_reexecute_job` | `python_modules/dagster/dagster/_core/execution/api.py:569` |
| Execution plan checkpoint | `ExecutionPlanSnapshot` with `steps`, `step_keys_to_execute`, `initial_known_state`, `job_snapshot_id`, `snapshot_version=1`, `executor_name` and `create_execution_plan_snapshot_id` | `python_modules/dagster/dagster/_core/snap/execution_plan_snapshot.py:44` and `:35` |
| Smallest checkpointable unit | `ExecutionPlan` built from `ExecutionStep`/`UnresolvedMappedExecutionStep`, addressed by `step.key`/`StepHandle`; `step_keys_to_execute` subset plan; `get_executable_step_deps` drives topological execution | `python_modules/dagster/dagster/_core/execution/plan/plan.py:664` and `:745` |
| Step handle taxonomy | `StepHandle`, `UnresolvedStepHandle`, `ResolvedFromDynamicStepHandle` and `StepHandleTypes` union | `python_modules/dagster/dagster/_core/execution/plan/plan.py:85` |
| Step event identity | `EventLogEntry` carries `run_id: str` and `step_key: str\|None`, plus nested `dagster_event.step_handle` to disambiguate dynamic resolves | `python_modules/dagster/dagster/_core/events/log.py:31` and `python_modules/dagster/dagster/_core/events/__init__.py:476` and `:579` |
| Per-step attempt state | `KnownExecutionState.previous_retry_attempts: Mapping[str,int]` and `RetryState._attempts` with `get_attempt_count`/`mark_attempt`/`snapshot_attempts` | `python_modules/dagster/dagster/_core/execution/plan/state.py:101` and `python_modules/dagster/dagster/_core/execution/retries.py:65` |
| Intra-run step retry | `ActiveExecution.mark_up_for_retry` enqueues to `_waiting_to_retry` or `_pending` depending on `RetryMode` and calls `_retry_state.mark_attempt(step_key)`; also `handle_event` dispatches `is_step_up_for_retry` | `python_modules/dagster/dagster/_core/execution/plan/active.py:550` and `:606` |
| Retry mode | `RetryMode.ENABLED/DISABLED/DEFERRED` with `for_inner_plan()` deferring retries to orchestrator (multiprocess) | `python_modules/dagster/dagster/_core/execution/retries.py:30` |
| Step restart event | `STEP_RESTARTED` emitted when `previous_attempt_count>0` on `core_dagster_event_sequence_for_step` and via `DagsterEvent.step_restarted_event` | `python_modules/dagster/dagster/_core/execution/plan/execute_step.py:448` and `python_modules/dagster/dagster/_core/events/__init__.py:1089` |
| Retry policy declaration | `RetryPolicy(max_retries, delay, backoff, jitter)` with `calculate_delay` | `python_modules/dagster/dagster/_core/definitions/policy.py:34` |
| Resume-from-failure reexecution | `KnownExecutionState.build_resume_retry_reexecution` validates parent `FAILURE/CANCELED` and calls `_derive_state_from_logs` to compute `steps_to_retry` and `known_state` | `python_modules/dagster/dagster/_core/execution/plan/state.py:170` |
| Past execution chaining | `PastExecutionState(run_id, produced_outputs, parent_state)` forms linked list across re-executions; `is_resume_retry` checks `RESUME_RETRY_TAG == "true"` | `python_modules/dagster/dagster/_core/execution/plan/state.py:32` and `python_modules/dagster/dagster/_core/storage/dagster_run.py:551` |
| Deriving retry set from logs | `_derive_state_of_past_run` scans `instance.all_logs(parent_run_id, of_type={STEP_START,STEP_FAILURE,STEP_SUCCESS,STEP_OUTPUT,STEP_SKIPPED})` to populate `to_retry`, `dynamic_mappings`, `output_set` | `python_modules/dagster/dagster/_core/execution/plan/state.py:225` |
| Auto-run retries | `auto_reexecution_should_retry_run` counts `len(run_group)-1 <= max_retries` via `instance.get_run_group(run.run_id)` and guards `RETRY_NUMBER_TAG`, `WILL_RETRY_TAG`, `AUTO_RETRY_RUN_ID_TAG` | `python_modules/dagster/dagster/_core/execution/retries.py:83` |
| Step vs run ownership | `DagsterRun` owns lifecycle `status`; `ActiveExecution` owns transient step state (`_success/_failed/_skipped/_in_flight/_waiting_to_retry`) and persists via `get_known_state()` returning `KnownExecutionState(previous_retry_attempts=snapshot)` | `python_modules/dagster/dagster/_core/execution/plan/active.py:117` and `:671` |
| Multi-process spanning | `ExecutionPlan.rebuild_from_snapshot` and `snapshot_from_execution_plan` enable ship-to-worker; `can_isolate_steps` gate; `step_handle_for_single_step_plans` hack for per-step subprocess | `python_modules/dagster/dagster/_core/execution/plan/plan.py:1027` and `:1161` |
| External process / orchestration tags | `dagster_execution_info` propagates `dagster/run-id`, `dagster/job`, `dagster/partition`, `dagster/code-location`; `RunLauncher.launch_run(LaunchRunContext)` and `ResumeRunContext.resume_attempt_number`, `CheckRunHealthResult.run_worker_id` | `python_modules/dagster/dagster/_core/storage/dagster_run.py:452` and `python_modules/dagster/dagster/_core/launcher/base.py:15` and `:26` and `:47` |
| Per-step stats / attempt count surface | `RunStepKeyStatsSnapshot.attempts` and `attempts_list` built from `STEP_START`/`STEP_SUCCESS`/`STEP_FAILURE` event stream in `build_run_step_stats_from_events` | `python_modules/dagster/dagster/_core/execution/stats.py:146` and `:255` |
| Schema lineage | `RunsTable.c.run_id` PK, `RunTagsTable.c.run_id FK → runs.run_id ON DELETE CASCADE`, indices on `run_tags(run_id)` | `python_modules/dagster/dagster/_core/storage/runs/schema.py:17` and `:75` and `:168` |

## Answers to Dimension Questions

**1. What is the durable user intention?**
The durable intention is a `DagsterRun` (`python_modules/dagster/dagster/_core/storage/dagster_run.py:236`). Creation goes through `instance.create_run_for_job` / `run_domain.create_run` (`python_modules/dagster/dagster/_core/instance/runs/run_domain.py:86`, `:977`) which atomically persists a `DagsterRun` + `ExecutionPlanSnapshot` + job snapshot before any compute starts. Fields `run_id` (UUID via `make_new_run_id`), `job_name`, `run_config`, `asset_selection`/`op_selection`/`resolved_op_selection`, `tags`, and `step_keys_to_execute` fully capture intent. Status lifecycle (`dagster_run.py:54` `QUEUED → NOT_STARTED → STARTING → STARTED → SUCCESS/FAILURE/CANCELED`) is durably stored in `RunStorage` (SQL row `RunsTable.c.run_id` at `runs/schema.py:17`). Survives process restarts; re-launch requires only `run_id`.

**2. What is the smallest checkpointable unit?**
The step. Physically `ExecutionStep` identified by `step.key` / `StepHandle` (`plan/plan.py:317`, `handle.py` referenced there) and logically persisted as `ExecutionStepSnap.key` (`snap/execution_plan_snapshot.py:159`). Checkpointing is event-sourced: each step emits `STEP_START` (`events/__init__.py:1081`), followed by zero or more `STEP_OUTPUT`, then terminal `STEP_SUCCESS` / `STEP_FAILURE` / `STEP_SKIPPED` / `STEP_UP_FOR_RETRY`. `EventLogStorage` durably stores these as `EventLogEntry(run_id, step_key)` (`events/log.py:38`). `KnownExecutionState.ready_outputs: set[StepOutputHandle>` and `dynamic_mappings` (`plan/state.py:110`, `:108`) accumulate which `StepOutputHandle`s have been materialized; they are snapshotted into `ExecutionPlanSnapshot.initial_known_state` (`snap/execution_plan_snapshot.py:52`) so a re-executed plan can skip already-successful outputs without re-running them (`active.py:269` `_has_produced_output` consults `parent_state.produced_outputs`). Stats layer reconstructs per-step attempt history by replaying the event log (`execution/stats.py:217`).

**3. Does each retry get a distinct identity?**
Depends on retry tier:
- *Step-level retry (within one run, governed by `RetryPolicy`)* — **No new identity.** The run keeps its `run_id`; the step keeps `step_key`. The attempt is modeled only as an in-memory counter `RetryState._attempts[step_key]` (`retries.py:67`) and as a derived `previous_retry_attempts` map in `KnownExecutionState` (`state.py:106`). Events `STEP_UP_FOR_RETRY` → delay → `STEP_RESTARTED` (`events/__init__.py:1089`, `plan/execute_step.py:448`) annotate ordering but the event envelope has no `attempt_id` or fencing nonce. `result.retry_attempts_for_node` (`execution/execution_result.py:232`) is computed by counting these events after the fact.
- *Run-level retry (reexecution/resume-from-failure or auto-reexecution)* — **Yes, distinct identity.** `create_reexecuted_run` / `_reexecute_job` (`instance/runs/run_domain.py:446`, `execution/api.py:508`) mints `DagsterRun(run_id=make_new_run_id(), root_run_id=…, parent_run_id=parent.run_id)`. Lineage is persisted both in columns `root_run_id`/`parent_run_id` and mirrored to tags `ROOT_RUN_ID_TAG`/`PARENT_RUN_ID_TAG` (`storage/tags.py:34`). Each retry chain can be enumerated via `get_run_group` (`storage/runs/sql_run_storage.py:555`). Tags `RESUME_RETRY_TAG`, `RETRY_NUMBER_TAG`, `WILL_RETRY_TAG`, `AUTO_RETRY_RUN_ID_TAG` (`storage/tags.py:42`, `:61`) mark its kind.

Thus intra-run retries are attempt-number-only, inter-run retries are fresh run rows.

**4. Can one logical operation span multiple runtime calls or OS processes?**
Yes. A logical run spans the launcher phase (`RunLauncher.launch_run(LaunchRunContext)` at `launcher/base.py:62`) and the executor phase. `ExecutionPlanSnapshot` is serdes-whitelisted and `can_reconstruct_plan` (`snap/execution_plan_snapshot.py:137`) — it is persisted alongside the run (`run_domain.py:350` `execution_plan_snapshot_id`) and reconstructed in any worker via `ExecutionPlan.rebuild_from_snapshot` (`plan/plan.py:1027`). `KnownExecutionState` (including `previous_retry_attempts`, `dynamic_mappings`, `parent_state`) is shipped with the snapshot and handed to each `StepExecutionContext` (`context/system.py:864` reads attempt count from known state). `can_isolate_steps` (`plan/plan.py:1161`) guarantees outputs use persistent IO so isolated step processes can hand off artifacts. Control plane examples: `multiprocess` executor forks subprocesses (see `plan.py:951` comment about single-step sub-plans), and `RunLauncher` abstractions for ECS/K8s delegate whole runs to remote processes, coordinated via `ResumeRunContext.resume_attempt_number` and health checks `check_run_worker_health` (`launcher/base.py:98`). The only shared durable state is `RunStorage` + `EventLogStorage`; there is no in-memory sticky session requirement.

**5. Can events be unambiguously attributed to the right entity?**
Within the persisted model, yes — modulo the late-attempt race:
- Every `EventLogEntry` is keyed by `(run_id, step_key, dagster_event.step_handle)` (`events/log.py:38`, `events/__init__.py:476`). Readers filter by `run_id` (`storage/event_log/base.py:226` `get_logs_for_run(run_id)` and `run_storage` per-run queries). Step outputs are keyed by `StepOutputHandle(step_key, output_name, mapping_key)` (`plan/outputs.py` via `StepOutputHandle` refs) so fan-in and dynamic `ResolvedFromDynamicStepHandle` resolves are unambiguous.
- Run-group membership is unambiguous via `root_run_id` column + `ROOT_RUN_ID_TAG` on the tag table; `get_run_group` returns exactly the set belonging to one logical operation (`sql_run_storage.py:555`).
- Parent vs child runs are unambiguous: `parent_run_id` and the `PastExecutionState` chain (`plan/state.py:32`) link re-executions; `is_resume_retry` (`dagster_run.py:551`) disambiguates.
- Gap: for intra-run retries, all attempts share `run_id+step_key` with no per-attempt envelope ID. Attribution of `STEP_OUTPUT` or `STEP_SUCCESS` to attempt N relies on temporal ordering between `STEP_UP_FOR_RETRY` and `STEP_RESTARTED`. If attempt 1 is late (e.g., stuck subprocess finally emits `STEP_SUCCESS` after attempt 2 already emitted `STEP_RESTARTED` + `STEP_SUCCESS`), the event log will contain two `STEP_SUCCESS` for the same `step_key` under the same `run_id` — the executor's state machine will treat the second as unknown/late but there is no fencing check that discards the stale attempt based on an `attempt_id`. The `previous_attempt_count` is a derived counter, not a causal token passed to the handler. For run-level retries the new `run_id` fully fences, so late events for the old `run_id` are naturally quarantined by `run_id` filter.

## Architectural Decisions

- **Single `DagsterRun` row as source of truth** (`storage/dagster_run.py:236`). All run-level metadata (config, selection, snapshots) is anchored to an immutable `run_id`. Tradeoff: simple linearizable writes via `run_storage.add_run`, but run row bloats with snapshot IDs; deletes cascade via FK (`runs/schema.py:75`).
- **Dual lineage encoding (columns + tags)**. `root_run_id`/`parent_run_id` are first-class columns for invariants and `get_run_group` resolution, but lineage is also mirrored to `dagster/parent_run_id` / `dagster/root_run_id` tags for SQL queryability without joins on the runs table (`sql_run_storage.py:578`). Decision preserves backward compat (tags existed before columns) and enables backfill asset slicing.
- **Event-sourced checkpoint via `ExecutionPlanSnapshot` + `KnownExecutionState`**. Rather than checkpointing arbitrary closure state, Dagster snapshots the declarative plan plus a minimal known state (`previous_retry_attempts`, `dynamic_mappings`, `ready_outputs`). Workers are stateless — they rebuild the `ExecutionPlan` from the snapshot (`plan/plan.py:1027`). Keeps retry/resume logic pure-functional.
- **Two-tier retry separation**. `RetryPolicy` (op-level, `policy.py:34`) for transient step failures vs. `ReexecutionOptions` / `auto_reexecution_should_retry_run` (run-level) for FAILED runs. Explicit `RetryMode.DEFERRED` (`retries.py:30`) delegates retry responsibility to the outer orchestrator so multiprocess/k8s executors do not double-retry.
- **Serdes-minimal `PastExecutionState` chain** (`plan/state.py:32`). Stores only `produced_outputs: set[StepOutputHandle>` per ancestor, not full event history, sufficient to decide skip vs re-run while keeping snapshot small.

## Notable Patterns

- **Immutable run group**: root → parent → child chain via `KnownExecutionState.parent_state` linked list; `get_run_group` reconstructs group by tag lookup, not recursive parent pointers, so orphaned middle runs are still discoverable via root tag.
- **Dynamic step resolution via `ResolvedFromDynamicStepHandle`**: unresolved `UnresolvedStepHandle` plus mapping key resolves to distinct `ResolvedFromDynamicStepHandle`s whose `to_key()` forms the event partition key; tracking dicts in `_derive_state_of_past_run` correctly coalesce unresolved vs resolved handles (`state.py:194`).
- **Unified event envelope**: `EventLogEntry` plus `DagsterEvent` with `DagsterEventType` enum and `event_specific_data` union — single storage table serves logs, metrics, and audit, enabling `build_run_step_stats_from_events` to synthesize `attempts` counts post-hoc (`stats.py:217`).
- **Tag-based correlation as API**: consumption filters (`RunsFilter`) accept arbitrary `tags` map (`dagster_run.py:612`); lineage tags are just another tag query, so UIs and sensors reuse the same path.
- **Deferred retry contract**: `ActiveExecution.mark_up_for_retry` checks `RetryMode` — if `ENABLED`, re-queues `_pending`; if `DEFERRED`, marks `_abandoned` and leaves retry to the daemon that will mint a new run. Pattern isolates retry policy from executor topology.

## Tradeoffs

- **Intra-run retry simplicity vs. safety**: Keeping retry within same `run_id` avoids run proliferation and keeps one event timeline per logical operation. Cost: no attempt fence ID, so idempotency depends on worker liveness; late outputs are not automatically rejected.
- **Snapshot-everything vs. snapshot-economy**: `ExecutionPlanSnapshot` snapshots every step+input+output eagerly, which makes any executor restart trivial, but snapshot size grows with plan width (hundreds of steps) and is duplicated per re-execution even when only a few steps are retried.
- **Tags for lineage**: Tag table scan for `get_run_group` is simple and index-friendly (`idx_run_tags_run_idx` at `runs/schema.py:168`) but tag writes are separate statement batch after run insert (`sql_run_storage.py:154`), creating a small window where lineage query is incomplete if crash occurs mid-create; column values mitigate but tag query path still races.
- **Event-order-dependent attempt accounting**: `RetryState` replay from `previous_retry_attempts` requires `KnownExecutionState` to be plumbed correctly; if an executor fails to persist `get_known_state()` (`active.py:671`) before crash, attempt count is lost and replay will undercount, potentially exceeding `max_retries`.
- **`PastExecutionState` minimalism**: Only `StepOutputHandle` set is retained per ancestor, so rich output metadata must be re-fetched from event logs if needed — trades snapshot size for extra log scans during resume.

## Failure Modes / Edge Cases

- **Late-arriving step output after retry**: Worker A (attempt 1) hangs, executor starts attempt 2 after `STEP_UP_FOR_RETRY` + `STEP_RESTARTED`; worker A's delayed `STEP_OUTPUT`/`STEP_SUCCESS` will be persisted under same `(run_id, step_key)` with no attempt token. `ActiveExecution.handle_event` will `mark_success` on the first success and move step to `_success`; a duplicate success event will call `_mark_complete` on a step no longer in `_in_flight` and raise `Invariant` — or, if timing allows, corrupt `_step_outputs` set. No explicit `attempt_id` comparison exists (`active.py:579`, `retries.py:65` — counter only). Mitigation: use run-level reexecution (new `run_id`) for fence-critical jobs.
- **Unknown-state abandonment**: If a step worker dies without emitting a terminal event, `ActiveExecution.verify_complete` (`active.py:627`) moves the step to `mark_unknown_state` → `_unknown_state` → raises `DagsterUnknownStepStateError` on context exit (`active.py:172`). This intentionally fails the run rather than silently succeeding, but leaves manual intervention to retry.
- **Lost `KnownExecutionState` on worker crash**: `get_known_state` is only materialized at handoff boundaries. If the orchestrator process crashes before persisting handoff, dynamic mappings and retry counts revert to the last persisted snapshot, causing duplicate dynamic step expansions or extra retries beyond `max_retries` (noted in `retries.py:102` extra-retry race).
- **Run-group expansion race**: `auto_reexecution_should_retry_run` notes manual retries launched concurrently with auto-retry counting can cause one extra run (`retries.py:102` comment). `get_run_group` read is non-atomic with the subsequent `will_retry` tag write.
- **Tag-only lineage orphan**: If `RunStorage.delete_run` cascades deletes from `RunTagsTable` (`runs/sqlite/sqlite_run_storage.py:172`), deleting the root run orphans group enumeration for remaining children — `get_run_group` then raises `DagsterRunNotFoundError` (`sql_run_storage.py:567`).
- **Step key non-uniqueness after dynamic resolution**: Improper handling of `ResolvedFromDynamicStepHandle` vs `StepHandle` identity can collapse multiple dynamic instances into one tracking key; `_update_tracking_dict` (`state.py:207`) is the subtle correct path — using `unresolved_form.to_key()` for resolved handles.

## Future Considerations

- **Add explicit attempt token to `EventLogEntry` / `DagsterEvent`**: Introduce `attempt: int` or `attempt_id: UUID` on `StepStartData`/`StepOutputData`, passed to the handler and validated on write (compare-and-set). This would let storage reject late attempt-1 events cleanly and give a direct answer to the dimension's fencing question. Must be serdes-compatible additive field.
- **Durable `RetryState` row**: Persist `previous_retry_attempts` as its own table or as part of `ExecutionPlanSnapshot` versioning so crash recovery accurately preserves count without relying on log replay.
- **Unified run-group table**: Materialize run groups in a dedicated `run_groups` table instead of tag-scans to eliminate the tag-write window and accelerate `get_run_group` under large backfill scales (currently linear in `run_tags` join).
- **Attempt-scoped log partitioning**: Store per-attempt event sub-streams (e.g., `run_id:attempt_N`) to make `GET logs for run/step/attempt` O(attempt) instead of filtering by timestamp slices.
- **Cross-process attempt clock**: Adopt monotonic attempt fencing via `DagsterRun.execution_plan_snapshot_id` version bump on each retry so stale workers' writes target a stale snapshot ID that storage can reject.

## Questions / Gaps

- **No evidence of per-attempt process ID binding**: Searched `launcher/base.py:47` `CheckRunHealthResult.run_worker_id` and `storage/tags.py:80` `RUN_WORKER_ID_TAG`; these exist for health checks but are not joined to step-attempt events, so we cannot trace “which OS pid produced which attempt” from the event log alone — flagged as missing forensic link.
- **No distributed lease / fencing token on step retry**: Verified `retries.py:65` `RetryState` only holds count, not token. If `attempt 2 succeeds after attempt 1 returns late`, the only IDs available to quarantine attempt 1 are `(run_id, step_key, timestamp ordering)` — no cryptographic `attempt_id`. Confirmed by absence of `attempt_id` in `EventLogEntry` (`log.py:31`) and `DagsterEvent` (`__init__.py`).
- **Unclear atomicity of run creation + snapshot persistence**: `run_domain.py:257` `add_run` then optional `_log_asset_planned_events` is two commits; not confirmed whether `add_run` itself is transactional across `runs` + `run_tags` + `snapshots` in Postgres vs SQLite adapters — requires deeper adapter verification.

---

Generated by `dimension 01.02` against `dagster`.
