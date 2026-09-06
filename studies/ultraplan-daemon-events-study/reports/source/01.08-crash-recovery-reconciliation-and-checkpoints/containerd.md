# Source Analysis: containerd

## Crash Recovery, Reconciliation, and Checkpoints

### Source Info

| Field | Value |
|-------|-------|
| Name | containerd |
| Path | `studies/ultraplan-daemon-events-study/sources/containerd` |
| Language / Stack | Go (daemon + per-container shim processes, bbolt `meta.db`, ttrpc/gRPC shim API, CRIU/CRI plugin) |
| Analyzed | 2026-09-03 |

## Summary

containerd survives daemon crashes by keeping execution in out-of-process shims and all intent in durable stores, then re-attaching on startup. The startup reconciliation has three independent layers: (1) `core/runtime/v2.ShimManager.LoadExistingShims` re-discovers bundles under the state dir and reconnects to live shims via per-bundle `bootstrap.json`, reaping empty/leaked shims; (2) the CRI plugin `recover()` rebuilds in-memory sandbox/container stores from `meta.db` containers plus per-container `status` JSON checkpoints, corrects them against live shim task status, and garbage-collects orphan ID dirs; (3) the `plugins/restart` monitor polls desired-vs-actual task state every 10s and starts/stops tasks to enforce `no/always/on-failure/unless-stopped` policies. Durable checkpoints are: bbolt `meta.db` (containers, images, leases, snapshots), shim bundle dir (`bootstrap.json`, `address`, `shim-binary-path`, `rootfs`), CRI `RootDir/StateDir/{sandboxes,containers}/<id>/status` JSON, and explicit CRIU task checkpoints (`rootfs-diff.tar` + CRIU image dir) for clone/migrate. Arbitrary in-memory execution never resumes — only re-attachable shim tasks and checkpointed metadata resume; interrupted pulls/transfers, execs, and in-flight creates are retried or failed by the caller, never transparently continued. SIGKILL during an external side effect is not decided safely: there is no write-ahead log or idempotency token for side effects, so restart either deletes a `Created` task and leaves retry to kubelet/client, or maps a vanished `Running` task to `UNKNOWN` exit 255 and requires manual/kubelet retry.

## Rating

**7 / 10 — Clear model with tests, explicit interfaces, and operational safeguards.**

Rationale: shim reconnect + CRI recovery + restart-policy reconciliation are explicit, bounded (load/cleanup/shutdown timeouts), tested (`restart_test.go`, `recover_test.go`, `status_test.go`), and clean up orphans (dead-shim events, bundle deletion, work-dir GC, CRI ID-dir GC, mount-manager BoltDB). Downgraded from 9 because retry-safety is best-effort: poll-based restart monitor has a documented `on-failure` race when the task was deleted, CRI recovery explicitly assumes no adds/touches while down, `NoSync` mode trades durability for throughput, and there is no transactional exactly-once guard for external side effects.

## Evidence Collected

Every entry MUST include a file path with line numbers. Format: `path/to/file.ts:NN`.

| Area | Evidence | File:Line |
|------|----------|-----------|
| Startup shim reconciliation entrypoint | `LoadExistingShims` walks `<stateDir>/<namespace>` then calls `loadShims` + `cleanupWorkDirs` per namespace | `sources/containerd/core/runtime/v2/shim_load.go:39-64` |
| Per-shim load with bounded budget | `loadShim` wraps whole load in `timeout.WithContext(ctx, loadTimeout)` (5s, `loadTimeout`/`cleanupTimeout`/`shutdownTimeout` constants); resolves runtime from `shim-binary-path` file or metadata store, then `loadShimTask` | `sources/containerd/core/runtime/v2/shim_load.go:128-187`, `sources/containerd/core/runtime/v2/shim.go:62-77` |
| Shim reconnect transport | `loadShim` re-opens shim log pipe with `AnonReconnectDialer` (fail-fast) and reconnects via `bootstrap.json` params (`ttrpc+unix` / `grpc+vsock`), `makeConnection` dials with 10s unix pre-dial to avoid stalling restart on dead shims | `sources/containerd/core/runtime/v2/shim.go:79-142`, `sources/containerd/core/runtime/v2/shim.go:339-398` |
| Durable shim identity | `writeBootstrapParams` atomically writes `bootstrap.json`; `restoreBootstrapParams` reads it or migrates legacy `address` file; `shim-binary-path` preserves custom runtime across restarts | `sources/containerd/core/runtime/v2/shim.go:293-331`, `sources/containerd/core/runtime/v2/shim_manager.go:350-379`, `sources/containerd/core/runtime/v2/shim_load.go:142-146` |
| Leaked-shim takeover rule | After load, if not a sandbox (`sandboxStore.Get` → NotFound) and `Pids` empty/NotFound → `cleanupShimTask`; otherwise `shims.Add`. Comment cites containerd#6860 (task never created due to crash) | `sources/containerd/core/runtime/v2/shim_load.go:189-226` |
| Dead-shim orphan publication | `cleanupAfterDeadShim` (5s budget) calls `binary.Delete`, synthesizes `TaskExit(exit 255 if no response)+TaskDelete` events so subscribers observe the loss | `sources/containerd/core/runtime/v2/shim.go:144-187` |
| Failed-start reaping | `cleanupShimTask` detaches context (`WithoutCancel`), bounds delete+shutdown by `cleanupTimeout`, closes client and leaves bundle for caller to remove | `sources/containerd/core/runtime/v2/shim.go:189-222` |
| Unmounted-rootfs + bundle GC on load failure | Empty bundle dir deleted (fast path); unresolvable runtime → `bundle.Delete()`; load failure → `bundle.Delete()` + error log, never blocks other namespaces | `sources/containerd/core/runtime/v2/shim_load.go:96-125`, `sources/containerd/core/runtime/v2/shim_load.go:161-166` |
| Stale work-dir GC after reboot | `cleanupWorkDirs` removes `<rootDir>/<ns>/<id>` dirs with no loaded shim — covers `/run` wiped on reboot while persistent work dir remains | `sources/containerd/core/runtime/v2/shim_load.go:264-296` |
| CRI startup reconciliation loop | `recover()` lists sandbox containers + SandboxStore entries, `RecoverContainer` per sandbox (errgroup), re-arms `WaitSandbox` monitors, then same for workload containers + `CheckImages`, then 4-way orphan dir sweep | `sources/containerd/internal/cri/server/restart.go:54-249` |
| CRI recovery preconditions (fragility contract) | Comment block: files/checkpoints MUST NOT be touched, containers SHOULD NOT be added, tasks SHOULD NOT be created/started while CRI plugin down | `sources/containerd/internal/cri/server/restart.go:45-53` |
| CRI container checkpoint schema | `Status{Pid, CreatedAt, StartedAt, FinishedAt, ExitCode, Reason, Message, Resources}` JSON `versionedStatus{Version:"v1"}` at `<root>/status`; `Starting/Removing/Unknown` explicitly `json:"-"` (not checkpointed); `StoreStatus/LoadStatus` via `continuity.AtomicWriteFile` | `sources/containerd/internal/cri/store/container/status.go:63-101`, `sources/containerd/internal/cri/store/container/status.go:164-195` |
| CRI status write discipline | `Update` = memory only; `UpdateSync` = memory + atomic checkpoint; tests assert non-sync updates do not persist | `sources/containerd/internal/cri/store/container/status.go:270-299`, `sources/containerd/internal/cri/store/container/status_test.go:152-205` |
| CRI interrupted-state mapping | Task NotFound + checkpoint CREATED → recreate IO (assumes start was interrupted, stays CREATED); + RUNNING → `FinishedAt=now, exit 255/Unknown`; task Created but checkpoint not CREATED → error; task Running but checkpoint EXITED → error; task Stopped → adopt `ExitStatus/ExitTime` then `Delete(WithProcessKill)` | `sources/containerd/internal/cri/server/restart.go:342-413` |
| CRI sandbox interrupted-state mapping | Task NotFound → NOTREADY; Running → re-`Wait` and READY; non-Running → `Delete` + NOTREADY; kubelet expected to stop NOTREADY sandboxes | `sources/containerd/internal/cri/server/podsandbox/recover.go:99-143`, `sources/containerd/internal/cri/server/podsandbox/recover.go:177-180` |
| CRI unknown-state sentinel | `unknownContainerStatus(){ExitCode:255, Reason:"Unknown", Unknown:true}`; `State()==CONTAINER_UNKNOWN` when `Unknown` | `sources/containerd/internal/cri/server/helpers.go:249-266`, `sources/containerd/internal/cri/store/container/status.go:103-118` |
| CRI orphan resource discovery | `cleanupOrphanedIDDirs`: list `<RootDir|StateDir>/{sandboxes,containers}`, delete ID dirs with no matching containerd container; warns on non-dir files, best-effort per-dir remove | `sources/containerd/internal/cri/server/restart.go:217-247`, `sources/containerd/internal/cri/server/restart.go:449-476` |
| Restart-policy reconciliation | `monitor.run` ticks `interval` (default 10s, TOML `plugins.restart.interval`); `reconcile` lists namespaces, `monitor()` diffs `containerd.io/restart.status` label vs live task status, `Reconcile()` gates on `policy/count/explicitly-stopped` | `sources/containerd/plugins/restart/monitor.go:37-64`, `sources/containerd/plugins/restart/monitor.go:90-129`, `sources/containerd/plugins/restart/monitor.go:131-194`, `sources/containerd/core/runtime/restart/restart.go:117-143` |
| Restart takeover semantics | `startChange.apply`: bump `restart.count` label, `killTask(SIGKILL+Delete)`, `NewTask+Start`; `stopChange.apply`: `killTask` only. Documented race: deleted task yields empty status → `on-failure` reconcile error | `sources/containerd/plugins/restart/change.go:39-103`, `sources/containerd/plugins/restart/monitor.go:164-176` |
| Durable metadata store | bbolt `meta.db` at `<root>/meta.db` (`NoFreelistSync`, flock timeout, `NoStatistics`); `no_sync` opts (`NoSync+NoGrowSync`) documented as crash-data-loss risk | `sources/containerd/plugins/metadata/plugin.go:67-74`, `sources/containerd/plugins/metadata/plugin.go:132-184` |
| Mount-manager crash tracking | Activation state in BoltDB + target dir; documented "Crash recovery: mounts tracked and cleaned up after daemon restart", GC + lease integration | `sources/containerd/docs/mounts.md:318-325` |
| Explicit task checkpoints (CRIU) | `Task.Checkpoint(path)` → `CheckpointTaskRequest`; `Create(... Checkpoint, RestoreFromPath)` unpacks sibling `rootfs-diff.tar` between Create and Start; `ctr checkpoint/restore` and CRI `checkpoint --export` flows | `sources/containerd/core/runtime/v2/shim.go:838-848`, `sources/containerd/core/runtime/v2/shim.go:643-717`, `sources/containerd/plugins/services/tasks/local.go:178-221`, `sources/containerd/cmd/ctr/commands/containers/checkpoint.go:29-96`, `sources/containerd/cmd/ctr/commands/containers/restore.go:34-89` |
| Checkpoint-gated restore events | `runtime.TaskCheckpointedEventTopic` SHOULD fire on checkpoint (runtime-v2 contract) | `sources/containerd/docs/runtime-v2.md:398` |
| Crash/restart tests | `TestRestartPolicyReconcile` table-tests policy×status; `TestRecoverContainer` fakes Running/error/deleted tasks; `status_test.go` covers checkpoint vs memory-only updates and corrupt decode | `sources/containerd/core/runtime/restart/restart_test.go:109-218`, `sources/containerd/internal/cri/server/podsandbox/recover_test.go:216-405`, `sources/containerd/internal/cri/store/container/status_test.go:63-205` |

## Answers to Dimension Questions

**1. What survives a client crash, daemon crash, and host reboot?**
- Client crash: everything — daemon, shims, `meta.db`, bundle dirs, CRI status files are daemon-side. Client holds no durable state; a new client re-lists via API.
- Daemon crash (host up, `/run` intact): survives = bbolt `meta.db` (`sources/containerd/plugins/metadata/plugin.go:169-184`), shim bundle dirs + `bootstrap.json`/`address` (`sources/containerd/core/runtime/v2/shim.go:293-331`), live shim processes (out-of-process, re-dialed at `sources/containerd/core/runtime/v2/shim.go:79-142`), CRI `status` JSON files (`sources/containerd/internal/cri/store/container/status.go:164-195`). Does not survive = in-memory CRI stores, exit monitors, restart-monitor tick state, exec processes attached to dead shims (re-published as synthetic exit/delete at `sources/containerd/core/runtime/v2/shim.go:144-187`).
- Host reboot (`/run` wiped): shim processes die; bundles under persistent state dir are re-scanned but fail to dial → deleted (`sources/containerd/core/runtime/v2/shim_load.go:96-125`); leftover `<rootDir>/<ns>/<id>` work dirs GC'd (`sources/containerd/core/runtime/v2/shim_load.go:264-296`); CRI tasks vanish → containers map to `UNKNOWN`/255 or `NOTREADY`, images re-checked, orphan ID dirs swept (`sources/containerd/internal/cri/server/restart.go:208-247`). CRIU checkpoint tarballs survive only if exported to content store/registry (`sources/containerd/docs/features.md:121-137`).

**2. Can arbitrary in-memory execution resume, or only checkpointed work?**
- Only checkpointed/re-attachable work. Running tasks resume by re-attaching to the surviving shim (`loadShimTask` + `PID` connectivity probe at `sources/containerd/core/runtime/v2/shim_load.go:228-262`); `Wait` channels are re-armed (`sources/containerd/internal/cri/server/restart.go:162-173`, `sources/containerd/internal/cri/server/podsandbox/recover.go:159-165`). Anything without a live shim — in-flight `Create/Start`, execs, streaming IO, pull/transfer progress — is not resumed; the model is re-attach-or-declare-dead (`exit 255/Unknown/NOTREADY`) and let the owner retry. Explicit CRIU checkpoints are the only whole-process memory snapshots, and they require an explicit `Checkpoint` call and `restore` path (`sources/containerd/core/runtime/v2/shim.go:838-848`).

**3. When is a new attempt started automatically?**
- Two automatic paths: (a) restart-monitor: if `containerd.io/restart.status=running` and live status diverges and `Reconcile(policy,status,labels)` returns true (`always`, `on-failure` with non-zero exit and under max-retries, `unless-stopped` without explicit-stop), it kills any remnant and `NewTask+Start`s, incrementing `restart.count` (`sources/containerd/plugins/restart/monitor.go:167-186`, `sources/containerd/core/runtime/restart/restart.go:117-143`, `sources/containerd/plugins/restart/change.go:45-79`). (b) Dead-shim path publishes exit/delete so orchestrators (kubelet) start replacements; CRI itself does not auto-start workload containers — it converges to `NOTREADY/UNKNOWN` and waits for kubelet (`sources/containerd/internal/cri/server/podsandbox/recover.go:177-180`). Sandbox `READY` tasks only get their exit monitor re-armed, not restarted.

**4. When is human/manual retry required?**
- Whenever reconciliation maps to a terminal-unknown state: `UNKNOWN` container status (`sources/containerd/internal/cri/server/helpers.go:256-266`), `NOTREADY` sandbox, `StateUnknown` on load errors (`sources/containerd/internal/cri/server/podsandbox/recover.go:145-147`), `on-failure` with exhausted `max-retries` or unparseable policy/count (`sources/containerd/core/runtime/restart/restart.go:124-135`), `unless-stopped` after explicit stop, `no` policy, and any pull/push/apply that died mid-stream (transfer progress is in-memory; resume is a fresh idempotent retry by digest, not a continuation — `sources/containerd/core/transfer/local/pull.go:122-184`). Corrupt `status` JSON → falls back to unknown status with a warning (`sources/containerd/internal/cri/server/restart.go:288-293`); missing metadata extension → container fails to load and is skipped, requiring operator cleanup (`sources/containerd/internal/cri/server/restart.go:274-286`).

**5. How are orphaned resources discovered?**
- Four complementary sweeps, all list-and-diff: (a) shim layer: empty bundle dirs, undialable shims, `shouldCleanupShim` (non-sandbox + no pids), and post-reboot work dirs without a loaded shim (`sources/containerd/core/runtime/v2/shim_load.go:96-125`, `sources/containerd/core/runtime/v2/shim_load.go:197-226`, `sources/containerd/core/runtime/v2/shim_load.go:264-296`); (b) dead-shim callback: `binary.Delete` + synthetic exit/delete events + removal from task list (`sources/containerd/core/runtime/v2/shim.go:144-187`); (c) CRI layer: 4-way `cleanupOrphanedIDDirs` diffing filesystem ID dirs against live container lists (`sources/containerd/internal/cri/server/restart.go:217-247`); (d) storage layer: mount-manager BoltDB activation tracking (`sources/containerd/docs/mounts.md:318-325`), lease-expiry GC (`containerd.io/gc.expire` at `sources/containerd/core/leases/lease.go:93-99`), and content/snapshot GC via metadata DB.

**SIGKILL the execution owner during an external side effect. On restart, how does the system decide whether retry is safe?**
- It does not — there is no commit marker, write-ahead log, or idempotency key for side effects, so safety is delegated to the caller and to content addressing. Evidence: the `Created`-task window explicitly deletes the task and keeps the container `CREATED`, pushing the decision back to the starter (`sources/containerd/internal/cri/server/restart.go:366-375`); the `Running`-task-vanished window fabricates `exit 255/Unknown` rather than claiming success or retrying (`sources/containerd/internal/cri/server/restart.go:354-362`); the restart monitor only looks at exit status + retry count, not at whether the side effect landed, and its own comment admits the deleted-task status race (`sources/containerd/plugins/restart/monitor.go:164-176`); CRI recovery refuses to reason about concurrent writers at all (`sources/containerd/internal/cri/server/restart.go:45-53`). Idempotent operations (content pull/push by digest, snapshot prepare) are safe to retry by construction; non-idempotent ones (network push that timed out after server applied, container start that sent `Start` to the shim before crashing) require human/orchestrator judgment — kubelet semantics (stop `NOTREADY`, recreate) are the de-facto retry arbiter, not containerd.

## Architectural Decisions

- **Shims outlive the daemon.** Execution lives in per-container shim processes; the daemon only holds a dial address (`bootstrap.json`) plus bundle paths. Restart is re-attachment, not replay (`sources/containerd/core/runtime/v2/shim.go:79-142`).
- **Two stores, two truths, task wins on conflict.** Desired/config truth in `meta.db` + CRI status files; actual truth from live shim `Task.Status`. Reconciliation corrects checkpoint toward task, never the reverse, except to synthesize `Unknown/255` when the task is gone (`sources/containerd/internal/cri/server/restart.go:342-413`).
- **Atomic file checkpoints, not DB writes, for hot status.** CRI status uses `continuity.AtomicWriteFile` + `Update` (memory) vs `UpdateSync` (durable) split to keep the fast path cheap (`sources/containerd/internal/cri/store/container/status.go:164-195`, `sources/containerd/internal/cri/store/container/status.go:270-299`).
- **Bounded, non-blocking startup.** Every shim load/delete/shutdown is wrapped in 3–5s timeouts; one wedged shim logs and is skipped/deleted instead of stalling the daemon (`sources/containerd/core/runtime/v2/shim.go:62-77`, `sources/containerd/core/runtime/v2/shim_load.go:128-140`).
- **Poll, don't watch, for restart policy.** The restart monitor is a 10s poll loop over labels vs task status, not an event-driven controller (`sources/containerd/plugins/restart/monitor.go:90-100`). Simple and crash-immune (no subscription to lose) at the cost of up-to-10s latency and races.
- **Explicit CRIU checkpoints, not continuous journaling.** Whole-container memory snapshots exist only when requested, stored as content/registry artifacts with a sibling `rootfs-diff.tar` applied between `Create` and `Start` (`sources/containerd/core/runtime/v2/shim.go:643-717`).

## Notable Patterns

- **Reconnect-dialer split:** `AnonDialer` (retry, for fresh start) vs `AnonReconnectDialer` (fail-fast, for reload) — prevents restart from hanging on dead sockets (`sources/containerd/core/runtime/v2/shim.go:336-338`).
- **Synthetic terminal events:** dead shims produce `TaskExit(255)+TaskDelete` so downstream (moby/kubelet) converges without special crash handling (`sources/containerd/core/runtime/v2/shim.go:173-186`).
- **Label-as-desired-state:** restart intent (`status/policy/count/explicitly-stopped/loguri`) lives in container labels, making it durable in `meta.db` and visible to any client (`sources/containerd/core/runtime/restart/restart.go:44-56`).
- **Versioned JSON status with downgrade tolerance:** `versionedStatus{Version:"v1"}` + `decode` switch anticipates upgrade skew (`sources/containerd/internal/cri/store/container/status.go:128-141`).
- **Best-effort, per-item isolation in recovery:** per-container 10s `loadContainerTimeout` + errgroup means one wedged container cannot fail the whole recovery (`sources/containerd/internal/cri/server/restart.go:262-270`).

## Tradeoffs

- **Polling simplicity vs latency/ races:** 10s restart interval survives subscription loss but delays restarts and hits the documented deleted-task/`on-failure` race (`sources/containerd/plugins/restart/monitor.go:90-100`, `sources/containerd/plugins/restart/monitor.go:164-176`).
- **fsync vs throughput:** `no_sync` (`NoSync+NoGrowSync`) boosts write throughput at explicit crash-data-loss risk; default is durable (`sources/containerd/plugins/metadata/plugin.go:67-74`, `sources/containerd/plugins/metadata/plugin.go:160-165`).
- **Strict preconditions vs robustness:** CRI recovery is simple because it assumes no adds/touches while down; violation returns load errors and skips containers rather than merging (`sources/containerd/internal/cri/server/restart.go:45-53`, `sources/containerd/internal/cri/server/restart.go:180-204`).
- **Fail-fast dial vs slow discovery:** 10s unix pre-dial + 5s load budget keeps startup fast but may reap a slow-starting shim as leaked (`sources/containerd/core/runtime/v2/shim.go:392-398`, `sources/containerd/core/runtime/v2/shim_load.go:197-209`).
- **CRIU power vs trust surface:** checkpoint archives carry annotations/image refs/paths that must be treated as untrusted (recent CVEs CRI-005); durability and portability come with re-validation cost (`sources/containerd/docs/security/THREAT_MODEL.md:428-436`).

## Failure Modes / Edge Cases

- Daemon SIGKILL between shim `Start` and status `UpdateSync` → checkpoint says CREATED/RUNNING-stale, task says otherwise; code patches `StartedAt/Pid` to now or maps to `Unknown/255` — the true start time and side-effect outcome are lost (`sources/containerd/internal/cri/server/restart.go:379-403`).
- Daemon SIGKILL during `NewTask+Start` in `startChange.apply` → retry on next 10s tick bumps `restart.count` again; no exactly-once guard (`sources/containerd/plugins/restart/change.go:64-78`).
- `meta.db` corruption on freelist read → mitigated by `NoFreelistSync=true` but not eliminated; `NoSync` mode can lose recent writes (`sources/containerd/plugins/metadata/plugin.go:132-137`, `sources/containerd/plugins/metadata/plugin.go:160-165`).
- Host reboot wipes shim sockets: every task becomes orphan; CRI converges to NOTREADY/UNKNOWN and depends on kubelet to delete/recreate — containerd never recreates workload containers itself (`sources/containerd/core/runtime/v2/shim_load.go:96-125`, `sources/containerd/internal/cri/server/podsandbox/recover.go:119-141`).
- Missing `ContainerMetadataExtension` or undecodable status → container skipped with error log; orphan dirs for it are then deleted as if the container never existed (`sources/containerd/internal/cri/server/restart.go:274-293`).
- Sandbox metadata without `SandboxStore` extension (pre-#11612 upgrade path) → leaked sandbox record explicitly deleted to allow startup (`sources/containerd/internal/cri/server/restart.go:107-121`).

## Future Considerations

- Add a per-operation intent log (or at minimum a `StartedAt`-with-`TaskCreateAttemptID`) so the `Created`-vs-`Running` crash window can distinguish "never started" from "started, effect unknown" without guessing.
- Consider event-driven restart (subscribe to `TaskExit/Delete`) with poll as fallback, to cut 10s latency and close the deleted-task/`on-failure` race.
- Expose recovery metrics (shims re-attached vs reaped, CRI UNKNOWN/NOTREADY counts, orphan dirs removed) — currently only logs.
- Document the `NoSync` durability contract and the "do not touch while down" CRI assumption in operator docs, not just code comments.

## Questions / Gaps

- No evidence found for transfer/pull resume across daemon restart: progress trackers are in-memory (`sources/containerd/core/transfer/local/progress.go:83-166`); content chunks already committed are reusable by digest, but whether an interrupted `Pull` auto-resumes or requires a new call was not traced beyond the handler wiring — searched `core/transfer/*` for `resume/checkpoint/recover`.
- No evidence found for a daemon-level crash/restart integration test (kill -9 daemon, restart, assert re-attachment): searched `integration/*restart*`, `shim_load*test*`, `recover_test.go` (unit fakes only).
- Lease-expiry interaction with crash windows (does a lease expire while the daemon is down and cause GC of still-needed content?) not determined — `WithExpiration` format found (`sources/containerd/core/leases/lease.go:93-99`) but GC scheduling vs wall-clock downtime not traced.

---
Generated by `Dimension 01.08: Crash Recovery, Reconciliation, and Checkpoints` against `containerd`.
