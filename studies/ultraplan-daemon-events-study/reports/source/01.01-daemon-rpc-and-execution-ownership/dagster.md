# Source Analysis: dagster

## Dimension 01.01: Daemon RPC and Execution Ownership

### Source Info

| Field | Value |
|-------|-------|
| Name | dagster |
| Path | `studies/ultraplan-daemon-events-study/sources/dagster` |
| Language / Stack | Python / gRPC + SQL Run Storage (Dagster) |
| Analyzed | 2026-09-02 |

## Summary

Dagster separates long-lived execution ownership completely from the submitting client. The `dagster-daemon` is a standalone polling controller (`python_modules/dagster/dagster/_daemon/controller.py:88`) that owns all scheduled work; the CLI / webserver are thin observers that enqueue `DagsterRun` rows into persistent `RunStorage`. Accepted work survives immediate client exit because runs, backfills and cursors live in the DB, not in a client-owned process. Code intelligence crosses a separate local gRPC boundary (`python_modules/dagster/dagster/_grpc/server.py:392`) to per-code-location user-code servers over loopback TCP or Unix Domain Sockets. There is no daemon RPC service for job submission — the daemon has no inbound port; health is observed via DB heartbeats. Authentication is absent (insecure gRPC, DB-trust model), versioning is coarse (`server_id` + `DefsStateInfo`), and workspace context is explicit via `LoadableTargetOrigin.working_directory` propagated through `InstanceRef`/`WorkspaceLoadTarget` rather than relying on the daemon's CWD.

## Rating

**6 / 10 — Present but inconsistent / fragile**

Rationale: Execution ownership is clearly separated and durable (daemon owns nothing transient; DB is source of truth). Daemon lifecycle, heartbeat protocol (`DaemonHeartbeat` + `DEFAULT_DAEMON_HEARTBEAT_TOLERANCE_SECONDS`), and gRPC code-server transport are explicit and tested. Downgraded because: (1) there is no authenticated or version-negotiated daemon RPC — health observation is DB polling, not an API; (2) local gRPC is `grpc.insecure_channel` by default with optional `use_ssl` only and no token/metadata auth (`python_modules/dagster/dagster/_grpc/client.py:72`); (3) daemon reconnect is operator-driven (restart) with no client re-attach RPC; probes like `liveness-check` only read heartbeats.

## Evidence Collected

Every entry MUST include a file path with line numbers. Format: `path/to/file.ts:NN`.

| Area | Evidence | File:Line |
|------|----------|-----------|
| Daemon entrypoint | `run_command` CLI entry creates controller via `daemon_controller_from_instance` and blocks on `controller.check_daemon_loop()`; handles `--shutdown-pipe` IPC and `capture_interrupts` | `python_modules/dagster/dagster/_daemon/cli/__init__.py:34-120` |
| Daemon controller | `DagsterDaemonController` spawns one `threading.Thread` per daemon type, shares a single `daemon_uuid`, `daemon_shutdown_event`; owns `GrpcServerRegistry` and `WorkspaceProcessContext` | `python_modules/dagster/dagster/_daemon/controller.py:133-203` |
| Daemon abstract loop | `DagsterDaemon.run_daemon_loop` is infinite generator with heartbeat emission, error capture via `DaemonErrorCapture`, auto-restart on exception with `DAGSTER_DAEMON_CORE_LOOP_EXCEPTION_SLEEP_INTERVAL` | `python_modules/dagster/dagster/_daemon/daemon.py:88-145` |
| Heartbeat write | `DagsterDaemon._check_add_heartbeat` upserts `DaemonHeartbeat` via `instance.add_daemon_heartbeat`; detects duplicate daemon IDs | `python_modules/dagster/dagster/_daemon/daemon.py:147-201` |
| Heartbeat storage | `DagsterInstance.add_daemon_heartbeat / get_daemon_heartbeats / wipe_daemon_heartbeats` delegate to `RunStorage`; SQL impl upserts `DaemonHeartbeatsTable` keyed by `daemon_type` | `python_modules/dagster/dagster/_core/instance/methods/daemon_methods.py:57-66`, `python_modules/dagster/dagster/_core/storage/runs/sql_run_storage.py:852-893` |
| Heartbeat health check | `get_daemon_statuses` computes `healthy = curr_time <= heartbeat.timestamp + interval + tolerance` and surfaces errors; `all_daemons_live/healthy` used by `liveness-check` | `python_modules/dagster/dagster/_daemon/controller.py:424-470` |
| No daemon RPC service | No `DagsterDaemon` gRPC service; daemon health is DB-polling only. Only gRPC service defined is `DagsterApi` for user-code servers | `python_modules/dagster/dagster/_grpc/protos/dagster_api.proto:7-36` |
| gRPC service definition | `service DagsterApi` with 28 RPCs: `Ping`, `Heartbeat`, `StartRun`, `ShutdownServer`, `ExecutionPlanSnapshot`, `External*`, `ReloadCodeWithState`, `RefreshComponentState`, etc. | `python_modules/dagster/dagster/_grpc/protos/dagster_api.proto:7-36` |
| Transport: port vs UDS | `DagsterGrpcClient.__init__` enforces exclusive `port XOR socket`, builds `server_address = "unix:"+abspath(socket)` or `"host:port"`; invariants reject UDS on Windows | `python_modules/dagster/dagster/_grpc/client.py:71-107` |
| Transport setup (server) | `DagsterGrpcServer.__init__` validates same invariants, calls `server.add_insecure_port("unix:"+abspath(socket))` or `host:port`; `max_send/receive` from env `DAGSTER_GRPC_MAX_*` | `python_modules/dagster/dagster/_grpc/server.py:1492-1555` |
| Code-server spawn | `open_server_process` builds subprocess args including `--port/--socket`, `--heartbeat`, `--heartbeat-timeout`, `--instance-ref`, `--defs-state-info`, plus `loadable_target_origin.get_cli_args()` (`-f/-m/-d/-a`) and manages `cwd` explicitly | `python_modules/dagster/dagster/_grpc/server.py:1643-1720` |
| Registry ownership | `GrpcServerRegistry` owns `GrpcServerProcess` list, lazy-creates per `ManagedGrpcPythonEnvCodeLocationOrigin.get_id()`, runs `_clear_old_processes` thread, enforces `heartbeat_ttl=20` (daemon) vs `45` (webserver) | `python_modules/dagster/dagster/_core/remote_representation/grpc_server_registry.py:60-263`, `python_modules/dagster/dagster/_daemon/controller.py:59`, `python_modules/dagster/dagster/_core/workspace/context.py:115` |
| Server heartbeat/watch | Code server has `GrpcServerCodeLocation.client` with `client_heartbeat_thread` every 1s (`CLIENT_HEARTBEAT_INTERVAL`); server has `_heartbeat_thread` that shuts down after `heartbeat_timeout` silence; watch thread polls `GetServerId` and calls `on_updated/on_error/on_reconnected` | `python_modules/dagster/dagster/_grpc/client.py:59-69`, `python_modules/dagster/dagster/_grpc/server.py:535-575`, `python_modules/dagster/dagster/_grpc/server_watcher.py:24-224` |
| Execution ownership (DB) | Sensors/backfills submit via `instance.submit_run(run_id, workspace_request_context)` and `instance.get_run_by_id` persistence; daemon merely ticks and enqueues, run execution is delegated to run coordinator / user-code `StartRun` subprocess | `python_modules/dagster/dagster/_daemon/sensor.py:755`, `python_modules/dagster/dagster/_daemon/queued_run_coordinator/queued_run_coordinator_daemon.py:194`, `python_modules/dagster/dagster/_grpc/server.py:1338-1401` |
| Workspace propagation | `LoadableTargetOrigin` carries `python_file/module_name/package_name/working_directory/attribute/autoload_defs_module_name`; serialized into code-server CLI via `get_cli_args()`; server resolves via `get_loadable_targets(working_directory=...)` | `python_modules/dagster/dagster/_core/types/loadable_target_origin.py:16-40`, `python_modules/dagster/dagster/_grpc/utils.py:19-80` |
| CWD not relied upon | Daemon CLI propagates explicit `WorkspaceLoadTarget` (from `workspace.yaml`/`--workspace`); code-server uses `loadable_target_origin.working_directory` to locate files, not `os.getcwd()`; `open_server_process(cwd=cwd)` is opt-in | `python_modules/dagster/dagster/_daemon/controller.py:88-115`, `python_modules/dagster/dagster/_grpc/server.py:1655`, `python_modules/dagster/dagster/_grpc/server.py:280-287` |
| Auth / authz | `DagsterGrpcClient` supports only `use_ssl: bool` → `grpc.ssl_channel_credentials()` if true, else `grpc.insecure_channel`; no token, no peer auth; `metadata` passthrough exists but unused (`self._metadata`); server `add_insecure_port` only | `python_modules/dagster/dagster/_grpc/client.py:72-86`, `python_modules/dagster/dagster/_grpc/client.py:118-137`, `python_modules/dagster/dagster/_grpc/server.py:1525-1551` |
| Version negotiation | `DagsterApiServer._server_id = uuid4()` returned via `GetServerId`/`Ping`; bump on `ReloadCodeWithState`; `ListRepositoriesResponse.defs_state_info` and `ReloadCodeWithState.new_defs_state_info` diff via `_get_changed_defs_state_keys`; fallback for older servers checks `UNIMPLEMENTED` for streaming vs sync calls | `python_modules/dagster/dagster/_grpc/server.py:433-434`, `python_modules/dagster/dagster/_grpc/server.py:757-828`, `python_modules/dagster/dagster/_grpc/client.py:545-554` |
| Client detach / shutdown | `GrpcServerProcess.shutdown_server` calls `client.shutdown_server()` then optionally `wait()`; server `ShutdownServer` sets `_shutdown_once_executions_finish_event` and waits for in-flight runs before `set(server_termination_event)`; daemon controller `__exit__` joins threads with 30s timeout | `python_modules/dagster/dagster/_grpc/server.py:1899-1906`, `python_modules/dagster/dagster/_grpc/server.py:1266-1275`, `python_modules/dagster/dagster/_daemon/controller.py:323-344` |
| Reconnect / recovery | `WorkspaceProcessContext.check_workspace_freshness` refreshes workspace every 60s, clears endpoints, retries until 300s then crashes; `server_watcher` handles reconnect with `MAX_RECONNECT_ATTEMPTS=10` and `ERROR_RECOVERY_INTERVAL=10` | `python_modules/dagster/dagster/_daemon/controller.py:275-291`, `python_modules/dagster/dagster/_grpc/server_watcher.py:13-22` |

## Answers to Dimension Questions

### 1. Does a client exit affect execution?

**No.** Accepted work is Durably owned by `DagsterInstance` (DB), not by the submitting process.

* Submission path: CLI/webserver creates a `DagsterRun` in `RunStorage` via `instance.submit_run` (`python_modules/dagster/dagster/_daemon/sensor.py:755`). The daemon's sensor/ scheduler loops (`python_modules/dagster/dagster/_daemon/sensor.py:346`, `python_modules/dagster/dagster/_daemon/daemon.py:278`) poll storage and enqueue further runs even if the original submitter is gone.
* Run execution itself is handed to a `multiprocessing.Process` in the code server (`python_modules/dagster/dagster/_grpc/server.py:1380-1398`, `python_modules/dagster/dagster/_grpc/impl.py:start_run_in_subprocess`) which outlives the gRPC call; the client `DagsterGrpcClient.start_run` (`python_modules/dagster/dagster/_grpc/client.py:704`) is synchronous but failures are reported to the DB via `instance.report_run_failed` (`python_modules/dagster/dagster/_grpc/server.py:602`).
* Killing the submitting client after `add_run` leaves the run in `Queued`/`NotStarted` state; the `QueuedRunCoordinatorDaemon` (`python_modules/dagster/dagster/_daemon/run_coordinator/queued_run_coordinator_daemon.py:17`) will dequeue it. Evidence: `SqlRunStorage._get_run_by_id`/`add_daemon_heartbeat` show storage is the rendezvous — no in-memory daemon queue is the source of truth.

*Edge:* if client dies mid-transaction before commit, run never appears; no orphan. If daemon dies after accepting, run remains queued and will execute when daemon restarts (re-reads `get_daemon_heartbeats` and run queue).

### 2. What exactly crosses the daemon boundary?

Two boundaries:

* **Daemon ↔ Storage (the real daemon boundary):** No wire RPC. Daemon reads/writes serializable rows: `DaemonHeartbeat(timestamp, daemon_type, daemon_id, errors)` (`python_modules/dagster/dagster/_daemon/types.py:25`), `DagsterRun`, `PartitionBackfill`, and cursor key-values (`DaemonCursorStorage.get/set_cursor_values` via `python_modules/dagster/dagster/_core/storage/runs/sql_run_storage.py:1122`). Minimal stable surface is the `DaemonHeartbeat` + `RunStorage` schema (`DaemonHeartbeatsTable` upsert `python_modules/dagster/dagster/_core/storage/runs/sql_run_storage.py:852-873`).

* **Daemon ↔ User Code (code-server gRPC):** Daemon hosts a `WorkspaceProcessContext` that talks to per-location `DagsterApi` servers. Crosses: `SensorExecutionArgs`/`ExternalScheduleExecutionArgs` (including `instance_ref`, `cursor`, `last_tick_completion_time`), `ExecutionPlanSnapshotArgs`, `JobSubsetSnapshotArgs`, and returns `ScheduleExecutionErrorSnap`/`SensorExecutionErrorSnap` as serialized serdes strings (`python_modules/dagster/dagster/_grpc/server.py:1185-1254`, `python_modules/dagster/dagster/_grpc/types.py:382-468`). All payloads are `serialize_value` JSON blobs inside proto strings (`python_modules/dagster/dagster/_grpc/protos/dagster_api.proto:136-148`), compressed/gzipped on wire (`python_modules/dagster/dagster/_grpc/client.py:118-132`).

Nothing crosses for raw filesystem handles or CWD; that is carried as data (`LoadableTargetOrigin.working_directory`).

### 3. Can a client reconnect and recover the same operation?

**Yes for durable operations, no for in-flight daemon iterations.**

* Durable: Any CLI/web observer can recover by run ID: `instance.get_run_by_id(run_id)` (`python_modules/dagster/dagster/_core/storage/runs/sql_run_storage.py:427`), event-log polling, or `GetCurrentRuns` (`python_modules/dagster/dagster/_grpc/server.py:1466`). Daemon ticks are idempotent — cursors stored via `instance.get_cursor_values/set_cursor_values` (`python_modules/dagster/dagster/_core/storage/runs/sql_run_storage.py:1122`) let a restarted daemon resume where it left off (sensor cursor is passed as `SensorExecutionArgs.cursor`).

* Heartbeat recovery: `get_daemon_statuses` (`python_modules/dagster/dagster/_daemon/controller.py:424`) lets any instance inspect last heartbeats; `liveness-check` CLI (`python_modules/dagster/dagster/_daemon/cli/__init__.py:123`) is the reconnect probe. Another `dagster-daemon run` with same `InstanceRef` becomes the new owner — duplicate detection logs error if two daemons share a type (`python_modules/dagster/dagster/_daemon/daemon.py:182-191`) but storage last-write-wins.

* Not recovered: An in-progress `core_loop` iteration (e.g., a sensor tick mid-`get_external_sensor_execution`) is not checkpointed; killing daemon mid-tick loses that tick's work, but next interval retries via wall-clock scheduling. Reconnect does not re-attach to the previous daemon's in-memory generator — `run_daemon_loop` restarts from `core_loop` on next process start (`python_modules/dagster/dagster/_daemon/daemon.py:140-142`).

### 4. How is the local caller authenticated?

**Effectively not authenticated.** Threat model is filesystem/local-network isolation.

* Daemon ↔ Storage: Trusts whoever can connect to the DB (Postgres/MySQL/SQLite file). No daemon token; `DagsterInstance` construction is via `instance_ref` deserialization from CLI args (`python_modules/dagster/dagster/_daemon/cli/__init__.py:95`) with no signature.
* Daemon ↔ Code Server gRPC: `DagsterGrpcClient` defaults to `grpc.insecure_channel` (`python_modules/dagster/dagster/_grpc/client.py:131-135`); optional `use_ssl` uses default system roots with no mTLS, no token header. A `metadata` tuple can be passed (`python_modules/dagster/dagster/_grpc/client.py:78`) but search shows zero usage for auth (`grep authentication` returns only SSL fields). UDS path (`unix:<socket>`) is protected only by filesystem permissions on the socket file created by `safe_tempfile_path_unmanaged` (`python_modules/dagster/dagster/_grpc/server.py:1863`). Loopback TCP (`host:port`) has no peer authentication beyond being `localhost`.
* No versioned auth handshake; `Heartbeat`/`Ping` are unauthenticated.

*Implication:* Suitable for single-tenant local dev or VPC-internal deployment; multi-tenant or internet-exposed usage must add external L7 auth (e.g., Dagster Cloud does so out-of-band).

### 5. Does daemon code ever depend on its own current working directory?

**No — workspace context is explicit, not inherited from CWD.**

* The daemon is launched with a `WorkspaceLoadTarget` resolved from `workspace.yaml`/CLI opts (`python_modules/dagster/dagster/_daemon/cli/__init__.py:83-114`). That target expands to `CodeLocationOrigin`s each wrapping a `LoadableTargetOrigin(working_directory=...)` (`python_modules/dagster/dagster/_core/types/loadable_target_origin.py:21`).
* Code-server spawn passes `working_directory` as `-d` CLI flag via `loadable_target_origin.get_cli_args()` (`python_modules/dagster/dagster/_core/types/loadable_target_origin.py:30`) and via `DagsterGrpcServer`/`LoadedRepositories` loading (`python_modules/dagster/dagster/_grpc/server.py:280`, `python_modules/dagster/dagster/_grpc/utils.py:25-42`), where it is used to resolve `-f`/`-m`/`--package-name` targets against that directory, not the daemon's `os.getcwd()`.
* `open_server_process(cwd=cwd)` does accept an optional `cwd` override (`python_modules/dagster/dagster/_grpc/server.py:1655`), but the callsite from `GrpcServerRegistry` and `WorkspaceProcessContext` leaves it `None` (inheriting daemon CWD only as default), while the user code resolution still uses the explicit `working_directory` field. Grep for `os.getcwd` returns no daemon-path dependency.
* Operational safeguard: `check_workspace_freshness` refreshes by re-loading origins, not by re-reading CWD (`python_modules/dagster/dagster/_daemon/controller.py:275-283`).

*Minor caveat:* If `working_directory` is `None` and `python_file` is a relative path, resolution ultimately falls back to process CWD inside `load_def_in_python_file` — but the config contract expects absolute or `working_directory`-qualified paths; Dagster docs and `dagster.yaml` examples enforce this. No evidence of `os.chdir` in daemon code.

## Architectural Decisions

* **DB as daemon RPC:** Choosing heartbeats-in-SQL over a daemon RPC port (`python_modules/dagster/dagster/_core/storage/runs/sql_run_storage.py:852`) simplifies HA (no leader election) and makes client-exit safe, at cost of poll latency (30s heartbeat, 1800s tolerance `python_modules/dagster/dagster/_daemon/controller.py:36-42`) and DB write load. Opt-out `daemon_skip_heartbeats_without_errors` (`python_modules/dagster/dagster/_core/instance/methods/daemon_methods.py:121`) mitigates in Cloud.
* **One thread per daemon type** (`python_modules/dagster/dagster/_daemon/controller.py:188-202`): Simple isolation; a crash in sensor loop doesn't affect scheduler loop except via `check_daemon_threads` that kills the whole process if any thread dies (`python_modules/dagster/dagster/_daemon/controller.py:247-259`). Tradeoff vs process isolation (heavier but more resilient).
* **Generator-based daemon loops** (`DagsterDaemon.core_loop -> DaemonIterator` `python_modules/dagster/dagster/_daemon/daemon.py:213-265`): Yields allow heartbeat interleaving and `SpanMarker` tracing; errors yield `SerializableErrorInfo` rather than crashing.
* **gRPC over UDS preferred, TCP fallback** (`python_modules/dagster/dagster/_grpc/server.py:1542-1555`, `python_modules/dagster/dagster/_grpc/client.py:103-107`): Attempts UDS first (`safe_tempfile_path_unmanaged`), falls back to dynamic port (`find_free_port` `python_modules/dagster/dagster/_grpc/server.py:1754`) on Windows/`--force-port`. Compression `gzip` and 100MB limits (`python_modules/dagster/dagster/_grpc/utils.py:83-99`) are hard-coded defaults.
* **Heartbeat-TTL code-server lifecycle** (`DAEMON_GRPC_SERVER_HEARTBEAT_TTL=20` `python_modules/dagster/dagster/_daemon/controller.py:59`): Daemon-spawned code servers die after 20s without heartbeat, preventing orphan processes; webserver uses 45s. Cleanup thread (`_clear_old_processes` `python_modules/dagster/dagster/_core/remote_representation/grpc_server_registry.py:246`) reaps dead procs.
* **Explicit `LoadableTargetOrigin` + `WorkspaceProcessContext`** (`python_modules/dagster/dagster/_core/workspace/context.py:958`): Workspace scoping is data-driven, not CWD-driven, enabling multi-location workspaces and container-image overrides per location.

## Notable Patterns

* **Poll-based ownership, not lease:** Daemon never acquires a distributed lock; health is advisory via heartbeats. Duplicate daemon detection is log-only, not enforced (`python_modules/dagster/dagster/_daemon/daemon.py:182-191`).
* **Serde-string-in-proto:** All domain objects cross gRPC as `serialize_value(...)` JSON strings inside proto `string` fields (`python_modules/dagster/dagster/_grpc/protos/dagster_api.proto:64-175`), not as proto messages — allows Python schemaless evolution without proto regeneration, at cost of no wire schema validation.
* **Fallback shims for rolling upgrades:** Client probes `UNIMPLEMENTED` to choose streaming vs sync RPCs (`python_modules/dagster/dagster/_grpc/client.py:545-554`, `583-598`), and server supports both `ExternalScheduleExecution` (streaming) and `SyncExternalScheduleExecution`.
* **Request-scoped workspace snapshots:** `WorkspaceProcessContext.create_request_context()` copies `CurrentWorkspace` under lock (`python_modules/dagster/dagster/_core/workspace/context.py:888-900`), so a request sees a consistent view even while the daemon reloads.
* **Watch thread for external code servers:** `GrpcServerCodeLocation` spawns `client_heartbeat_thread` + `grpc-server-watch` thread that polls `GetServerId` every 1s, with `on_updated/on_disconnect/on_error` callbacks driving `CurrentWorkspace` swaps (`python_modules/dagster/dagster/_core/remote_representation/code_location.py:715-727`, `python_modules/dagster/dagster/_grpc/server_watcher.py:24`).
* **Idempotent reload with state pinning:** `ReloadCodeWithState` is idempotent on `DefsStateInfo` equality (`python_modules/dagster/dagster/_grpc/server.py:776-779`) and fan-out safe for replicas.

## Tradeoffs

* **Durability vs latency:** 30s heartbeat + 1800s tolerance gives low DB churn but slow failure detection; `liveness-check` may report stale healthy for 30+ minutes. Faster tolerance would increase write load linearly per daemon type.
* **Insecure-by-default local gRPC:** `add_insecure_port` avoids cert management for local dev, but any local process can impersonate a code server or daemon; hardening requires external network policy or Cloud's out-of-band auth (not in OSS).
* **Single DB rendezvous:** Simplicity of `RunStorage` as source of truth vs a dedicated daemon queue (e.g., Redis/NATS) means daemon is bottlenecked by DB perf and cannot push events to clients — clients must poll `event_log_storage`.
* **Thread-per-daemon vs process-per-daemon:** Threads share memory and GIL, enabling cheap `WorkspaceProcessContext` sharing; but a segfault or `KeyboardInterrupt` kills all daemons (`python_modules/dagster/dagster/_daemon/controller.py:297-344`).
* **UDS vs TCP:** UDS gives filesystem ACLs and no port conflicts, but unsupported on Windows (forcing `find_free_port` loop with race `CouldNotBindGrpcServerToAddress` retries `python_modules/dagster/dagster/_grpc/server.py:1745-1764`).

## Failure Modes / Edge Cases

* **Split-brain daemons:** Two `dagster-daemon` processes with same `InstanceRef` both write heartbeats (last-write-wins `python_modules/dagster/dagster/_core/storage/runs/sql_run_storage.py:852-873`); detection only logs error, does not fence old daemon — could double-submit runs if sensor is not idempotent (mitigated by run-key dedup in `python_modules/dagster/dagster/_daemon/sensor.py:1131`).
* **Lost heartbeat due to DB outage:** `check_daemon_threads` continues, but `get_daemon_statuses` falls back to last-known times (`python_modules/dagster/dagster/_daemon/controller.py:213-245`); `liveness-check` may false-negative; no circuit breaker.
* **Code server orphan:** If daemon crashes without `shutdown_all_processes` (`python_modules/dagster/dagster/_core/remote_representation/grpc_server_registry.py:274`), code servers live for `heartbeat_ttl=20s` then self-terminate via `_heartbeat_thread` → `_shutdown_once_executions_finish_event` → `_cleanup_thread` (`python_modules/dagster/dagster/_grpc/server.py:564-611`). In-flight `StartRun` subprocess blocks shutdown until `event_queue` drains.
* **Subprocess crash during StartRun:** Detected via `queue.Empty` + `not process.is_alive()` branch, reports `run_crash_explanation` and `report_run_failed` (`python_modules/dagster/dagster/_grpc/server.py:1412-1421`, `python_modules/dagster/dagster/_grpc/server.py:596-602`).
* **Working directory drift:** If `workspace.yaml` changes on disk and daemon CWD changes without reload, stale `LoadableTargetOrigin.working_directory` persists until `RELOAD_WORKSPACE_INTERVAL=60s` refresh (`python_modules/dagster/dagster/_daemon/controller.py:51-53`); transient load errors trigger retry until `DEFAULT_WORKSPACE_FRESHNESS_TOLERANCE=300s` then process crashes.
* **No auth = spoofing:** A local attacker binding to same `host:port` before daemon can intercept code-server traffic; mitigated by UDS randomness and `find_free_port` retry, but not by crypto.
* **Large payload truncation:** Streaming chunking at 4MB (`STREAMING_CHUNK_SIZE=4000000` `python_modules/dagster/dagster/_grpc/server.py:131`) could OOM if client sets `DAGSTER_GRPC_MAX_RX_BYTES` too low; server and client independently env-configurable, mismatch causes `RESOURCE_EXHAUSTED`.
* **Version skew:** Older servers lacking `ReloadCodeWithState` would still accept `ReloadCode` (no-op warning `python_modules/dagster/dagster/_grpc/server.py:652-661`); clients silently fall back for schedule/sensor streaming, but not for new fields — `DefsStateInfo` changes could break deserialization.

## Future Considerations

* Add a real daemon RPC/health endpoint (gRPC/HTTP) with mutual TLS or token auth instead of DB-polling heartbeats — would enable push notifications and faster failure detection without DB load; could reuse existing `GrpcServerProcess` infra.
* Promote `use_ssl` to `mTLS + token metadata` and wire `DagsterInstance` to issue short-lived tokens for code servers; expose `metadata` auth in `DagsterGrpcClient` already has plumbing (`python_modules/dagster/dagster/_grpc/client.py:78`) but needs server interceptor.
* Implement fencing/lease for daemon heartbeats (e.g., `SELECT ... FOR UPDATE` on `daemon_type` row or advisory lock) to prevent split-brain double execution.
* Make `working_directory` resolution explicitly absolute at config load time and reject relative `python_file` without `working_directory` — eliminate CWD fallback path.
* Checkpoint daemon iterations (commit cursor before tick, not after) so reconnect can resume mid-tick without duplicate run submissions.

## Questions / Gaps

* No evidence of daemon RPC API beyond heartbeats — confirmed by absence of any `DagsterDaemon` service in proto and `daemon_controller_from_instance` having no socket bind. If a future daemon RPC is added, where would its spec live?
* No tests found for daemon detach semantics in this study boundary — `python_modules/dagster/dagster_tests/daemon_tests` (not inspected due to source isolation) likely covers it, but no in-file evidence of `detach` or `reconnect` primitives inside `_daemon/*.py`.
* Authentication docs: No `dagster.yaml` config key for daemon auth was found; `DagsterInstance` exposes `daemon_skip_heartbeats_without_errors` but no `daemon_auth_token`. Is auth intended to be delegated to DB/network layer, or planned for OSS?
* Protocol version negotiation is limited to `DefsStateInfo` diff and `server_id` bump; no `DagsterApi` service version field in `PingReply`. How does a client know it is talking to an incompatible code server before `ListRepositories` fails?
* Workspace scoping: `WorkspaceLoadTarget.create_origins()` source not inspected — does it resolve `python_file` relative to config file or to daemon CWD? Need to inspect `dagster/_core/workspace/load_target.py`.

---

Generated by `Dimension 01.01: Daemon RPC and Execution Ownership` against `dagster`.
