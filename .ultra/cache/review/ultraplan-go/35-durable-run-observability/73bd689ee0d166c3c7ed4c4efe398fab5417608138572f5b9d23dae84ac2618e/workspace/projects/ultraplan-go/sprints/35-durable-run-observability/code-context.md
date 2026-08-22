# Sprint Code Context

## Sprint Scope

Sprint 35 replaces the web server's session-scoped, in-memory operation identity and event history with a durable workspace-wide run control plane. The relevant implementation spans operation acceptance and execution, CLI/TUI/web adapters, agentwrap correlation, lifecycle and cancellation semantics, product-owned locks and checkpoints, reconciliation, API compatibility, SSE replay, redaction, bounded delivery, shutdown behavior, and browser projections.

The existing sprint, study, flow, execute, review, smoke, and runtime records remain authoritative for their product-specific outcomes. The new operational model must correlate with those records without treating them as a substitute run repository or inferring successful execution from artifacts.

## Inspected Repository Areas

- `internal/app`: shared operation contracts, preparation fingerprints, runtime-backed dispatch, web composition, reconciliation, and cleanup-uncertainty routing.
- `internal/web`: server lifecycle, in-memory operation hub, HTTP and SSE handlers, routes, session authorization, browser projections, diagnostics, compatibility fixtures, and concurrency tests.
- `internal/tui`: direct preparation/execution and lossy in-process progress delivery.
- `internal/platform/runtime`: agentwrap supervision, run/session/attempt identity, canonical event mapping, cancellation, bounded event retention, raw-payload omission, and diagnostic redaction.
- `internal/sprint`: mutation locking, execute run state, planning-stage session checkpoints, runtime correlation, cleanup uncertainty, and interrupted-work reconciliation.
- `internal/study`: durable run/task state, process lock ownership, PID-based liveness and cancellation, cleanup uncertainty, and interrupted-run reconciliation.
- Public API schema, route, lifecycle, SSE, cancellation, capacity, shutdown, and compatibility tests.

## Selected Source References

### Shared Operation Boundary

- **Path:** `internal/app/operations.go`
- **Lines:** `22-145`
- **Symbol:** `OperationalUseCases`, `WebOperations`, `OperationCleanupRecorder`, `OperationReconciler`, `OperationRequest`, `OperationEvent`, `OperationResult`
- **Rationale:** Defines the current adapter-neutral operation vocabulary and lifecycle payloads. It has no durable run or attempt identity, owner/lease data, event sequence, heartbeat, or repository capability, making it the central compatibility boundary for introducing shared run observability without moving workflow semantics into `internal/web`.

- **Path:** `internal/app/operations.go`
- **Lines:** `149-262`
- **Symbol:** `dashboardUseCases.PrepareOperation`
- **Rationale:** Establishes confirmation, mutation classification, runtime detection, governed inputs, durable refresh links, and pre-execution fingerprints. Durable run acceptance must fit between successful confirmation and child execution, preserve stale-input protection, and avoid advertising acceptance before required persistence succeeds.

- **Path:** `internal/app/operations.go`
- **Lines:** `265-390`
- **Symbol:** `dashboardUseCases.RunOperation`
- **Rationale:** Shows the common execution entry point currently used by web and TUI, including re-preparation, fingerprint validation, initial/terminal events, runtime-free operations, and delegation of runtime-backed work. This is the immediate boundary where durable run context and correlation must reach product services.

- **Path:** `internal/app/operation_runner.go`
- **Lines:** `15-139`
- **Symbol:** `sharedOperationRunner`
- **Rationale:** Enumerates all runtime-backed product operations and translates product progress into shared events. It demonstrates that workflow completion remains owned by sprint/study services and identifies where run, stage, task, and runtime-attempt correlation must be propagated consistently.

### Current Web Authority And Lifecycle

- **Path:** `internal/web/operations.go`
- **Lines:** `18-143`
- **Symbol:** operation limits, `operationDocument`, `operationRecord`, `operationHub`
- **Rationale:** Captures the existing bounded but entirely process-local authority: records, sessions, events, subscribers, counters, active count, and cancellation functions all live in one mutex-protected map. These types and limits are the baseline to replace or demote to transport-only state.

- **Path:** `internal/web/operations.go`
- **Lines:** `146-269`
- **Symbol:** `operationHub.startConfirmed`, `operationHub.run`, `operationHub.finish`
- **Rationale:** Defines current acceptance ordering and terminal arbitration. An operation is inserted, counted, and exposed before a goroutine invokes product execution, but nothing is persisted; terminal state is first-writer-wins only within one process. This is the critical fail-closed acceptance and single-terminal-outcome seam.

- **Path:** `internal/web/operations.go`
- **Lines:** `271-385`
- **Symbol:** `operationHub.status`, `operationHub.activeOperations`, `operationHub.cancelOperation`, `operationHub.subscribe`
- **Rationale:** Demonstrates the principal sprint defect: reads, lists, cancellation, and event subscription require the originating browser session and the same in-memory hub. Active counts exclude CLI/TUI and other server processes, while cancellation can only reach a locally retained context.

- **Path:** `internal/web/operations.go`
- **Lines:** `388-428`
- **Symbol:** `operationHub.appendEventLocked`
- **Rationale:** Current event sequencing, size bounds, retention, redaction projection, and slow-subscriber disconnection are implemented during in-memory fan-out. This is the baseline for durable monotonic sequencing, explicit replay gaps, write-before-publish ordering, and retaining bounded transport buffers rather than transport authority.

- **Path:** `internal/web/operations.go`
- **Lines:** `439-535`
- **Symbol:** `operationHub.drainAndWait`, `persistCleanupUncertain`, `markCleanupUncertain`, `reapLocked`
- **Rationale:** Defines shutdown cancellation, deadline handling, product-specific uncertainty persistence, and ten-minute terminal eviction. It exposes race and migration concerns because the operation record itself disappears on restart or expiry even when a separate cleanup marker survives.

- **Path:** `internal/web/operations.go`
- **Lines:** `537-629`
- **Symbol:** `projectOperationResult`, `terminalOperationState`, `safeProjectedText`, `parseEventID`
- **Rationale:** Provides current terminal-state vocabulary, result bounding, basic secret/path redaction, and decimal cursor parsing. Durable recording and support diagnostics must preserve or strengthen these safety limits before persistence and fan-out.

### HTTP, SSE, And Compatibility

- **Path:** `internal/web/operation_handlers.go`
- **Lines:** `120-239`
- **Symbol:** operation start, status, active-list, cancel, and event handlers
- **Rationale:** Defines current API behavior and ordering: start returns `202` and a process-local URL, all observation and mutation use session identity, and SSE replays only hub memory. These handlers are the main interface adapters that must resolve durable runs while keeping commands separate from observation.

- **Path:** `internal/web/operation_handlers.go`
- **Lines:** `307-334`
- **Symbol:** `handleHTMLOperationStatus`, `handleHTMLOperationCancel`
- **Rationale:** Contains the user-visible `Operation not retained` dead end required by the sprint to eliminate. Both inspection and cancellation currently map missing or cross-session records to an undifferentiated `404`.

- **Path:** `internal/web/operation_handlers.go`
- **Lines:** `555-653`
- **Symbol:** `writeOperationError`, `safeOperationCause`, `logOperation`
- **Rationale:** Defines typed API errors, recovery hints, and structured operation logging. It is the compatibility and telemetry seam for typed retention gaps, storage failures, stale ownership, cancellation uncertainty, and shared run correlation.

- **Path:** `internal/web/routes.go`
- **Lines:** `244-350`
- **Symbol:** `allowedMethods`, `matchRoute`
- **Rationale:** Freezes the current `/api/v1/operations` and `/operations/{id}` route shapes and methods. Durable inspection, listing, replay, cancellation, and legacy-link resolution must account for these public URLs.

- **Path:** `internal/web/operations_contract_test.go`
- **Lines:** `82-189`
- **Symbol:** lifecycle, SSE, and error compatibility contract tests
- **Rationale:** Records the browser-consumed lifecycle states, stable SSE event names, terminal classification, durable-refresh fields, and operation error codes. Public model changes require coordinated server, browser, fixture, and documentation migration.

- **Path:** `internal/web/api_compatibility_test.go`
- **Lines:** `11-73`
- **Symbol:** `TestAPICompatibilityRouteMethodMatrix`, `TestAPICompatibilityTransportSchemas`
- **Rationale:** Explicitly freezes operation routes and JSON field ordering/schema for v1 clients. It is the strongest repository evidence that run API evolution must be additive, versioned, or accompanied by an explicit compatibility rationale.

- **Path:** `internal/web/operations_test.go`
- **Lines:** `124-268`
- **Symbol:** operation lifecycle, shutdown, capacity, replay, and subscriber tests
- **Rationale:** Verifies current session isolation, idempotent local cancellation, shutdown cleanup markers, active-operation limits, bounded queues, and slow-subscriber eviction. These tests identify behavior to preserve where appropriate and assumptions that must change for workspace-wide visibility.

### Server Lifecycle And Reconciliation

- **Path:** `internal/web/server.go`
- **Lines:** `38-151`
- **Symbol:** `Run`
- **Rationale:** Shows startup reconciliation occurs once before handler creation, every server creates its own operation root and hub, and server shutdown cancels all work it accepted. It is central to multiple-server topology, periodic reconciliation, observer/owner separation, and deciding whether owner exit interrupts or cancels work.

- **Path:** `internal/app/web_usecases.go`
- **Lines:** `195-258`
- **Symbol:** `RecordOperationCleanupUncertain`, `ReconcileOperations`
- **Rationale:** Routes shutdown uncertainty into product-owned sprint/study files and scans all products at web startup. This is existing cross-product recovery composition, but it has no workspace run repository, periodic backlog, or per-run diagnostics.

### Other Local Surfaces

- **Path:** `internal/tui/app.go`
- **Lines:** `210-253`
- **Symbol:** `confirmationCmd`, `operationCmd`, `waitOperationEvent`
- **Rationale:** The TUI invokes the same app operation methods directly but keeps identity and progress only in its model and a lossy channel. It identifies the integration point needed for TUI-started work to receive durable IDs and appear in workspace projections.

- **Path:** `internal/web/static/app.js`
- **Lines:** `236-264`
- **Symbol:** active-operation refresh
- **Rationale:** The browser top bar fetches `/api/v1/operations` alongside dashboard data, so its count inherits the hub's session-local scope. This is the concrete browser projection that must consume workspace lifecycle state rather than planning-stage status or page-local state.

- **Path:** `internal/web/static/app.js`
- **Lines:** `716-856`
- **Symbol:** operation SSE, start, reconnect, and cancel client logic
- **Rationale:** Shows the browser's EventSource lifecycle, terminal handling, preparation/start sequence, and cancellation requests. It is required context for cursor resume, durable snapshots, typed replay gaps, refresh behavior, and legacy API compatibility.

### Runtime Supervision And Correlation

- **Path:** `internal/platform/runtime/runtime.go`
- **Lines:** `16-105`
- **Symbol:** `Request`, `Result`, `Event`, `EventStats`
- **Rationale:** Defines UltraPlan's agentwrap boundary and existing trace, session, provider, run, event, timestamp, attempt, and omission metadata. A product run model should correlate these identities rather than duplicate provider/process supervision or persist raw provider payloads.

- **Path:** `internal/platform/runtime/runtime.go`
- **Lines:** `246-337`
- **Symbol:** `Adapter.StartRun`
- **Rationale:** Shows child creation, canonical event consumption, callback timing, cancellation propagation to agentwrap, bounded wait after cancellation, and result arbitration. Durable run acceptance must occur before this method starts a child, while cancellation coordination ultimately needs to reach this owner path.

- **Path:** `internal/platform/runtime/runtime.go`
- **Lines:** `339-435`
- **Symbol:** `eventCollection`
- **Rationale:** Implements bounded in-memory retention and explicit dropped-event statistics for runtime results. It supplies useful precedent for bounded storage and observable compaction while confirming that current runtime retention cannot serve cross-process replay.

- **Path:** `internal/platform/runtime/runtime.go`
- **Lines:** `473-548`
- **Symbol:** `toAgentwrapRequest`, `mapResult`, `mapEvent`
- **Rationale:** Maps UltraPlan metadata into agentwrap and returns agentwrap run/session/event identities while deliberately omitting raw event content. This is the safe correlation and redaction boundary for durable operational events.

- **Path:** `internal/platform/runtime/runtime.go`
- **Lines:** `584-629`
- **Symbol:** `mapError`, `mapSDKError`, `redactDiagnosticString`
- **Rationale:** Defines provider error sanitization before errors cross the runtime boundary. Durable diagnostics should consume these safe mapped forms and must not persist native response bodies, headers, credentials, or unrestricted metadata.

### Existing Durable Product State

- **Path:** `internal/sprint/execute_state.go`
- **Lines:** `14-146`
- **Symbol:** execute run-state load/save and legacy handling
- **Rationale:** Demonstrates schema validation, legacy terminal compatibility, temp-file write, file sync, atomic rename, and directory sync for sprint-owned execute state. It is evidence for durability expectations and a migration boundary that must remain authoritative for execute tasks.

- **Path:** `internal/sprint/execute_state.go`
- **Lines:** `149-238`
- **Symbol:** `ValidateExecuteRunState`, `isTerminalExecuteStatus`
- **Rationale:** Defines execute task lifecycle invariants, timestamps, diagnostics, evidence safety, and terminal statuses. Shared runs may project this progress but must not redefine task success or rewrite canonical execute evidence.

- **Path:** `internal/sprint/session_state.go`
- **Lines:** `19-147`
- **Symbol:** planning-stage session checkpoints and `startPlanningStageRun`
- **Rationale:** Persists provider session identity by sprint stage and resumes compatible interrupted sessions. This checkpoint is separate from product run identity and requires an explicit correlation/migration story rather than replacement.

- **Path:** `internal/sprint/locks.go`
- **Lines:** `13-121`
- **Symbol:** mutation lease context and `ReconcileInterruptedMutation`
- **Rationale:** Shows the existing cross-process sprint mutation lock and conservative startup reconciliation: live lock holders are not rewritten, while orphaned running execute tasks and expired flow attempts become interrupted evidence. The shared run layer must coexist with this product-owned exclusion and terminal authority.

- **Path:** `internal/study/run_state_domain.go`
- **Lines:** `5-118`
- **Symbol:** `RunState`, `TaskState`, `StatusSummary`, `LockInfo`
- **Rationale:** Defines the existing durable study run ID, task states, runtime metadata, timestamps, active counts, and PID lock information. It is a distinct product run model whose identifiers and lifecycle must be correlated without being mistaken for the workspace-wide operational run.

- **Path:** `internal/study/run_state.go`
- **Lines:** `15-100`
- **Symbol:** `NewRunState`, `ReconcileRunState`
- **Rationale:** Shows current study run ID generation and task-graph reconciliation. Timestamp-derived IDs and one current state per study are insufficient as general workspace run identity but are compatibility-critical persisted behavior.

- **Path:** `internal/study/locks.go`
- **Lines:** `15-159`
- **Symbol:** `AcquireRunLoopLock`, `RunLoopActive`, `CancelRunLoop`
- **Rationale:** Implements cross-process ownership and cancellation with an exclusive file, PID liveness probe, acquisition timestamp, and `SIGINT`. It is the only existing cross-process cancellation mechanism, but PID-only checks are vulnerable to PID reuse and provide no renewable lease or fencing token.

- **Path:** `internal/study/cleanup_uncertain.go`
- **Lines:** `19-149`
- **Symbol:** `CleanupUncertainRecord`, `RecordCleanupUncertain`, `ReconcileInterruptedRun`
- **Rationale:** Persists shutdown uncertainty and conservatively converts active study tasks to interrupted cancellation evidence only after acquiring product ownership. It is important precedent for fail-visible recovery, while also exposing limited owner identity and one-shot reconciliation.

## Relationships

- Web and TUI adapters call `app.PrepareOperation` and `app.RunOperation`; runtime-backed operations delegate through `sharedOperationRunner` into sprint or study services.
- Sprint and study services own workflow locks, task states, artifacts, checkpoints, and product terminal outcomes.
- Product services call `internal/platform/runtime`, which starts and supervises agentwrap runs and maps safe canonical events and identities back to UltraPlan.
- The web operation hub currently wraps the app boundary with a second, browser-specific lifecycle. It generates an `op_*` ID, retains events and cancellation context in memory, filters all reads by browser session, and disappears with the server.
- SSE is fed directly from hub mutation under the same mutex; its sequence is monotonic only for one retained in-memory record and cannot survive process changes.
- Web startup reconciles product-owned state before constructing a fresh hub. Shutdown cancels locally owned contexts and may write separate cleanup-uncertainty markers, but neither path preserves the operation document or journal.
- The dashboard top bar reads the same session-local operation endpoint, while CLI/TUI executions bypass the hub entirely. This explains the cross-surface count and discovery failures.
- Agentwrap already owns provider/process supervision and exposes run, session, event, attempt, usage, cleanup, and safe raw-omission facts. UltraPlan's durable model should correlate these with product/stage/task identity while leaving supervision at that boundary.
- Existing API fixtures freeze v1 routes, JSON schemas, lifecycle names, SSE event names, and error codes; compatibility work must be deliberate and test-backed.

## Constraints

- `internal/web` may retain HTTP/SSE connection state and bounded subscriber buffers, but durable identity, lifecycle arbitration, discovery, leases, and cancellation coordination must live outside it.
- Acceptance persistence must complete before `sharedOperationRunner` or `runtime.Adapter.StartRun` can start child work and before HTTP/TUI/CLI surfaces advertise the run.
- Product locks and state remain authoritative for workflow mutation and outcomes. Reconciliation must not infer success from missing processes, existing artifacts, or stale planning-stage status.
- Operation visibility must no longer depend on browser session identity; mutation authorization and CSRF/session checks remain separate concerns.
- Agentwrap run IDs and provider session IDs are correlation identities, not substitutes for stable workspace product-run and attempt identities.
- Events must be sanitized and bounded before durable append and fan-out. Existing raw omission, diagnostic redaction, path suppression, result limits, and slow-subscriber isolation are minimum safety behavior.
- Terminal arbitration must cover completion, failure, cancellation, timeout, interruption, cleanup uncertainty, shutdown, and reconciliation races across processes, not only one hub mutex.
- PID existence alone is not sufficient liveness because the study lock has no process-birth identity, heartbeat, lease expiry, or fencing token.
- Existing `.flow-state.json`, execute `.run-state.json`, study run state/history, stage sessions, cleanup markers, locks, smoke evidence, and operation URLs require explicit compatibility handling.
- Public `/api/v1` route and JSON changes are constrained by compatibility fixtures and the dependency-free embedded browser client.
- Loopback-only local-web policy remains the current security boundary; expanded host topology is not established by the implementation.

## Open Questions

- No existing package provides a workspace-wide durable operational run repository; the appropriate package boundary and storage primitive remain unresolved by source.
- The implementation does not establish whether multiple local processes may safely write one prospective repository directly or require a coordinator.
- Existing sprint mutation locks and study locks have different ownership formats and semantics; the source does not define a common lease, fencing, heartbeat, or process-birth contract.
- CLI command paths do not currently pass a shared operation/run context through every runtime-backed entry point; the exact composition seam needs confirmation during reasoning.
- Agentwrap capabilities visible through this repository do not show a durable UltraPlan event store or a cross-process cancellation registry; dependency capabilities should be verified before selecting coordination.
- Retention duration, compaction format, disk quota, replay-gap response shape, metrics surface, support-bundle format, and optional tracing export are not selected in current implementation.
- Legacy operation IDs have no persisted mapping today, so recovery after restart can only use product scope links, not reconstruct historical operation identity.
