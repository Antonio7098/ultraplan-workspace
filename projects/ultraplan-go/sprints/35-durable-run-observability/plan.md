# Sprint Plan: Durable Run Identity and Cross-Surface Observability

> Project: `ultraplan-go`
> Sprint: `35-durable-run-observability`
> Source: `reasoning.md`
> **Inputs Used:** `projects/ultraplan-go/project-index.md`, `projects/ultraplan-go/roadmap.md`, `projects/ultraplan-go/docs/ARCHITECTURE.md`, `projects/ultraplan-go/docs/PRD.md`, `projects/ultraplan-go/docs/TRD.md`, `projects/ultraplan-go/sprints/35-durable-run-observability/requirements.md`, `projects/ultraplan-go/sprints/35-durable-run-observability/code-context.md`, `projects/ultraplan-go/sprints/35-durable-run-observability/sprint-index.md`, `projects/ultraplan-go/sprints/35-durable-run-observability/technical-handbook.md`, `projects/ultraplan-go/sprints/35-durable-run-observability/reasoning/api-design.md`, `projects/ultraplan-go/sprints/35-durable-run-observability/reasoning/architecture.md`, `projects/ultraplan-go/sprints/35-durable-run-observability/reasoning/frontend.md`, `projects/ultraplan-go/sprints/35-durable-run-observability/reasoning.md`, `../ultraplan-go/go.mod`, `../ultraplan-go/internal/app/app.go`, `../ultraplan-go/internal/app/surfaces.go`, `../ultraplan-go/internal/platform/config/config.go`

This plan executes `reasoning.md`. It does not reopen the selected package, SQLite store, same-host topology, ownership, lifecycle, API, presentation, telemetry, or failure-policy decisions.

The selected evidence reports were not reread individually because `technical-handbook.md`, the three area-reasoning documents, and final `reasoning.md` contain the decision-relevant findings and no unresolved decision required deeper report evidence. Live implementation files were inspected only to confirm composition, configuration, dependency, command, and documentation seams; `code-context.md` remains the broader implementation inventory.

## Reasoning Source

- **Sprint Reasoning:** `reasoning.md`
- **Sprint Index:** `sprint-index.md`
- **Technical Handbook:** `technical-handbook.md`
- **Area Reasoning:** `reasoning/api-design.md`, `reasoning/architecture.md`, `reasoning/frontend.md`

## Sprint Status

- **Status:** `in progress`
- **Owner:** `manual execution agent`
- **Start Date:** `2026-08-21`
- **Completion Date:** `pending`

## Decisions To Execute

| Decision | Source Section | Execution Implication |
| --- | --- | --- |
| Same-Host Run-Control Product Boundary | `reasoning.md#decision-1-same-host-run-control-product-boundary` | Add focused `internal/runcontrol`, compose it through `internal/app`, support multiple local processes over one canonical local-filesystem workspace, and keep web/TUI/product/runtime authorities separate. |
| Direct Multi-Process SQLite Repository | `reasoning.md#decision-2-direct-multi-process-sqlite-repository` | Persist operational truth in private `.ultraplan/run-control.db` with direct fenced writers, SQLite transactions, WAL, `synchronous=FULL`, foreign keys, a 5-second busy timeout, short transactions, and no daemon or in-memory authority fallback. |
| Stable Identity, Acceptance, and Authority Hierarchy | `reasoning.md#decision-3-stable-identity-acceptance-and-authority-hierarchy` | Assign opaque 128-bit run/attempt IDs and durably accept and claim each asynchronous execution before product or child work starts; retain typed, safe product/runtime/process correlations without conflating authorities. |
| Fenced Ownership, Lease, and Conservative Reconciliation | `reasoning.md#decision-4-fenced-ownership-lease-and-conservative-reconciliation` | Implement fixed owner tick, heartbeat, lease, scan, and grace timings; verify process birth; fence every owner write; never adopt work or infer success during reconciliation. |
| Durable Cancellation and One Operational Terminal Winner | `reasoning.md#decision-5-durable-cancellation-and-one-operational-terminal-winner` | Persist cancellation before routing, acknowledge under the current fence, propagate through existing context/runtime cleanup, and arbitrate all terminal proposals through one immutable compare-and-set. |
| Sanitized Ordered Journal, Replay, and Backpressure | `reasoning.md#decision-6-sanitized-ordered-journal-replay-and-backpressure` | Redact and bound events before append, allocate sequence transactionally, publish only committed events, provide at-least-once cursor replay, coalesce only equivalent high-frequency progress, and disconnect slow subscribers. |
| Explicit Snapshot, Retention, Quota, and Failure Policy | `reasoning.md#decision-7-explicit-snapshot-retention-quota-and-failure-policy` | Preserve current snapshots and replay boundaries while enforcing per-run limits, timed compaction, tombstones, soft/hard quotas, reserved headroom, and fail-closed or fail-visible persistence behavior. |
| Additive Run API, Compatibility, Authorization, and Migration | `reasoning.md#decision-8-additive-run-api-compatibility-authorization-and-migration` | Add canonical run CLI/API resources, retain exact operation compatibility projections, separate workspace reads from fresh mutation authority, provide typed gaps/tombstones/legacy recovery, and implement backed-up one-way migrations with stop-and-restore rollback. |
| Shared Cross-Surface Presentation | `reasoning.md#decision-9-shared-cross-surface-presentation` | Add server-rendered workspace run list/detail pages and shared CLI/TUI/browser vocabulary; keep lifecycle, liveness, product status, cancellation, gaps, and uncertainty distinct without client-side authority. |
| Local Telemetry, Support, and Release Verification | `reasoning.md#decision-10-local-telemetry-support-and-release-verification` | Correlate safe logs, health, local metrics, diagnostics, and bounded support export; defer OpenTelemetry and broad metrics; require deterministic, race, process, browser, documentation, review, build, and gated real-runtime evidence. |

## Requirements / Contracts To Satisfy

`RO-1` through `RO-11` and `AC-1` through `AC-13` use the numbering defined in `reasoning.md#sprint-purpose`.

| Contract / Requirement ID | Required Behavior | Evidence Planned |
| --- | --- | --- |
| Architecture; RO-1 through RO-7; AC-10 | Run control is a focused product capability outside adapters; product modules and agentwrap retain their authorities. | Import-boundary tests, constructor/composition tests, and an authority trace from acceptance through terminal projection. |
| RO-1, RO-2; AC-1 | Durable run/attempt identity and first owner claim exist before any asynchronous product or child execution; required persistence failure starts nothing. | Acceptance transaction tests, command-path spies, crash-boundary tests, and runtime child-start assertions. |
| CLI Surface; RO-3, RO-4, RO-7; AC-2, AC-3, AC-7 | CLI, JSON, TUI, dashboard, list, and detail expose the same workspace lifecycle and recovery facts. | Cross-surface fixtures and a two-CLI/two-server process integration scenario. |
| RO-2, RO-4, RO-5; AC-4 through AC-6 | Sanitized committed events replay by durable sequence after session/server/transport changes, with typed cursor gaps and current snapshots. | Concurrent append, replay, restart, cursor, compaction, slow-subscriber, and browser cursor-handoff tests. |
| Workflows; Security; RO-1, RO-6, RO-7; AC-8, AC-9 | Leases, process birth, fencing, reconciliation, cancellation, and terminal arbitration are conservative and race-safe. | Fake-clock/process-probe tests, stale-writer rejection, PID-reuse tests, duplicate cancellation, owner-kill, and terminal race tests. |
| Persistence And Migrations; RO-2, RO-8; AC-5, AC-6, AC-12 | SQLite mutations are durable and atomic; migration, backup, integrity, compatibility, tombstone, restore, and rollback behavior is explicit. | Real SQLite migration/restore fixtures, WAL/locking checks, unsupported-schema tests, corruption/disk/permission injection, and legacy URL contracts. |
| Security; LLM Runtime; LLM Evaluation / Cost / Safety; RO-1, RO-2, RO-9; AC-11 | Raw prompts, provider payloads, secrets, unsafe paths, and unbounded output never enter persistence or presentation by default. | Adversarial redaction, omission, path, strict DTO, HTML escaping, logs, diagnostics, and support-bundle tests. |
| Configuration; Performance; RO-2, RO-3, RO-5, RO-9; AC-12 | Retention, replay, compaction, polling, subscriber, reconciliation, and quota work is bounded with validated safe settings. | Boundary/config precedence tests, quota/headroom tests, bounded-query/DOM/replay checks, and measured local multi-observer behavior. |
| Errors; Observability; RO-4 through RO-9 | Store, gap, tombstone, owner, cancellation, migration, and reconciliation failures remain typed, safe, correlated, and actionable. | Exact API/CLI fixtures, degraded health scenarios, structured-log fields, local metrics, diagnostics JSON, and support export. |
| Testing; Documentation; RO-10, RO-11; AC-13 | Normal/race/build/browser suites pass and gated real-runtime dogfood proves cross-surface observation or reports a truthful blocker; documentation covers operation and recovery. | Required commands, Architecture Review, Sprint Review, Deep Smoke Sprint evidence, and updated governed/user/operator documentation. |

## Execution Order And Release Rule

Tasks are ordered so repository and lifecycle invariants exist before adapters consume them. Intermediate commits may compile behind internal composition changes, but no releasable path may start asynchronous runtime work unless durable acceptance and the first owner claim are active for every CLI, TUI, and web entry point. If implementation evidence contradicts a final reasoning decision, stop, record the deviation in `reasoning.md` through the governed reasoning stage, and regenerate this plan before continuing.

## Tasks

- [x] **Task 1: Establish The Run-Control Domain And SQLite Contract**
  > Executes: Decisions 1-3; Architecture; Persistence And Migrations; RO-1, RO-2; AC-1, AC-10
  - [x] Add `internal/runcontrol` with explicit run, attempt, target, lifecycle, liveness, cancellation, terminal, event, omission, correlation, record-state, snapshot, query, health, and typed-error models; keep the package free of web, TUI, sprint, study, and runtime imports.
  - [x] Implement validation and generation for opaque `run_<base32-128-bit>` and `att_<base32-128-bit>` identifiers without encoded path, time, PID, owner, provider, or authority data.
  - [x] Define narrow repository/service interfaces and injected clock, owner/process probe, logger, and optional notifier seams; use context only for cancellation/deadlines.
  - [x] Add and pin a reviewed `modernc.org/sqlite` version compatible with Go 1.26 and the supported release platforms; record binary/build impact in release evidence.
  - [x] Implement `.ultraplan/run-control.db` creation with directory mode `0700`, database mode `0600`, WAL, `synchronous=FULL`, foreign keys, 5-second busy timeout, and bounded connection/transaction behavior.
  - [x] Create transactional `runs`, `attempts`, immutable `events`, `operation_aliases`, and bounded `reconciliation_log` schema records and indexes for active queries, `(run_id, sequence)` replay, owner control, retention, and reconciliation.
  - [x] Prove snapshot/lifecycle invariants, unique identity, concurrent writers, monotonic sequence allocation, immutable terminal compare-and-set, and private permissions against real temporary SQLite databases.
  - [x] **Stop condition:** local Linux locking, WAL/FULL policy, transaction isolation, private permissions, reopen integrity, pure-Go Linux/macOS amd64/arm64 builds, and package race tests passed; app integration may proceed without fallback.

- [x] **Task 2: Implement Schema Migration, Configuration, Backup, And Restore**
  > Executes: Decisions 2, 7, 8; Configuration; Persistence And Migrations; RO-2, RO-8; AC-12
  - [x] Add schema versioning with SQLite `user_version` plus an application schema record, unsupported-newer-schema rejection, and an advisory workspace migration lock shared by local UltraPlan processes.
  - [x] Before migration, checkpoint WAL, create a timestamped bounded backup, apply one-way transactional migrations, and run integrity checks before enabling acceptance or mutation.
  - [x] Extend existing configuration precedence and source reporting with only full-history duration, tombstone duration, and workspace quota; preserve fixed lease, heartbeat, event-size, replay, and per-run limits.
  - [x] Apply defaults of 7-day full terminal history, 30-day tombstones, 496 MiB soft quota, 512 MiB hard quota, and 16 MiB reserved headroom; validate minimum 1-hour full history, 24-hour tombstones, 64 MiB hard quota, and 16 MiB headroom.
  - [x] Reject unknown/unsafe combinations and redact configuration sources/values in diagnostics; update config text/JSON fixtures without weakening existing precedence.
  - [x] Implement tested restore instructions and fixtures that stop all UltraPlan processes and restore both the matching binary and pre-migration database backup; do not dual-write to old/new stores.
  - [x] Leave existing flow, execute, study, session, lock, cleanup, artifact, Git, and smoke records untouched; an initial repository is empty and never fabricates historical runs.

- [x] **Task 3: Make Durable Acceptance Unavoidable In App Composition**
  > Executes: Decisions 1-3; Architecture; LLM Runtime; Workflows; RO-1 through RO-3; AC-1 through AC-3, AC-10
  - [x] Add app-level run-control use cases and explicit composition in `cmd/ultraplan`/`internal/app`, opening one workspace repository per process and injecting it into CLI, TUI, and web capabilities without globals or adapter cycles.
  - [x] Inventory every asynchronous/runtime-backed command path, including sprint planning stages, execute/review/smoke/verify work and study start/resume operations; make all paths enter one acceptance wrapper rather than relying only on the current web hub.
  - [x] Preserve current preparation, confirmation, canonical-request, governed-input fingerprint, product-lock, and product-state checks; insert durable acceptance after valid preparation/confirmation and before goroutine creation or product/runtime invocation.
  - [x] Persist the accepted snapshot and first fenced owner claim before starting work. If either required write fails, return a typed persistence failure and prove no product operation, subprocess, runtime, or agentwrap child began.
  - [x] Record operation kind/target and safe product, stage/task, runtime, agentwrap, provider-session, process, and external-harness correlations only when observed; never fabricate missing values or persist prompt/provider-native content.
  - [x] Preserve product modules as the authority for locks, artifacts, task status, stage outcome, retry/resume, and validation. Project product status separately from the operational lifecycle.
  - [x] Store only a non-reversible confirmation-token digest for retained start deduplication; retrying the same canonical request through another server returns the same run, while conflicting reuse returns typed `idempotency_conflict`.
  - [x] Add command-path tests for every inventoried entry and an import/authority test that prevents adapter or run-control ownership drift.
  - [x] **Stop condition:** no surface integration task closes until a coverage table proves that every asynchronous entry uses durable acceptance and no bypass can start work when the store is unavailable.

- [x] **Task 4: Add Fenced Owner Control, Liveness, And Reconciliation**
  > Executes: Decision 4; Workflows; Security; Observability; RO-1, RO-6; AC-8, AC-9
  - [x] Generate a random owner identity per process lifetime and capture safe host digest, boot identity, PID, and exact process-birth token through platform-specific injected probes; uncertainty must never authorize signaling.
  - [x] Run the owner control loop at 1 second, commit heartbeats every 5 seconds, use 15-second leases, scan every 10 seconds in long-lived web/TUI and active CLI owners, and wait 45 seconds after lease expiry before a terminal reconciliation proposal.
  - [x] Use SQLite time for lease comparisons and require run ID, attempt ID, owner ID, and repository-allocated fencing generation on every heartbeat, append, cancellation acknowledgement, progress/snapshot, and terminal mutation.
  - [x] Run one bounded startup reconciliation pass before accepting runtime work, then periodic bounded idempotent scans without leader election; expose backlog and decision evidence.
  - [x] Classify unexpired ownership as active, expired ownership with exact live process birth as stalled/suspect, and missing/mismatched ownership through grace as interrupted or cleanup-uncertain candidates; never infer success from lease, PID, artifact, lock, or stage status.
  - [x] Do not adopt orphaned workers. Correlate product-specific reconciliation results without rewriting product state or making them operational success proof.
  - [x] Add fake-clock, clock-jump, owner-stall, owner-kill, PID-reuse, process-birth mismatch, stale-fence, repeated-reconciler, and unsupported-probe tests on Linux and macOS-supported paths.

- [x] **Task 5: Persist Sanitized Ordered Events Before Delivery**
  > Executes: Decision 6; Security; Performance; Observability; LLM Runtime; RO-2, RO-5, RO-9; AC-4, AC-6, AC-11, AC-12
  - [x] Define the canonical event envelope and allowed event types with stable run ID, decimal sequence, commit time, attempt/stage/task correlation, safe payload, and omission metadata; preserve known/unknown usage semantics.
  - [x] Centralize allowlisting, path suppression, diagnostic redaction, raw omission, and size projection before persistence, logs, subscribers, JSON, HTML, TUI, metrics, or support export.
  - [x] Allocate each per-run sequence and update the durable snapshot in the append transaction; invoke transport callbacks/notifiers only after commit so every visible event is replayable.
  - [x] Cap encoded events at 16 KiB and replace oversize detail with a safe warning/omission event; never persist full prompts, raw provider payloads, credentials, arbitrary stdout/stderr, or unrestricted paths.
  - [x] Persist at most one equivalent high-frequency progress update per run per 250 ms and record omitted count/time range; never coalesce lifecycle, warning, finding, artifact, cancellation, recovery, or terminal events.
  - [x] Implement indexed replay and polling at 250 ms while catching up and up to 1 second while idle; treat in-process notification only as an optimization and delivery as at-least-once with `(run_id, sequence)` deduplication.
  - [x] Bound replay to 512 historical events per SSE connection and bound app/web/TUI delivery queues; disconnect slow consumers so they resume from durable cursors without blocking execution or append.
  - [x] Test concurrent ordering, commit-before-delivery, duplicate producer/network delivery, polling across two processes, oversize replacement, hostile content, coalescing facts, slow consumers, replay catch-up, and dropped local TUI delivery.

- [x] **Task 6: Route Durable Cancellation And Arbitrate Terminal Outcomes**
  > Executes: Decision 5; Workflows; LLM Runtime; Errors; Security; RO-6, RO-7; AC-8, AC-9
  - [x] Implement cancellation as a transactionally persisted command with canonical trusted reasons, request timestamp/state, owner routing, acknowledgement, and uncertainty facts; do not use SSE or browser disconnect as a command channel.
  - [x] Have the current fenced owner poll on its 1-second control tick, durably acknowledge once, then cancel its owning context so existing product and agentwrap cleanup paths remain authoritative.
  - [x] Make duplicate requests idempotent and make already-terminal requests return the immutable winner without a second signal, cancellation proposal, or lifecycle rewrite.
  - [x] Route completion, failure, timeout, cancellation, interruption, cleanup uncertainty, server shutdown, persistence degradation, and reconciliation proposals through one immutable terminal compare-and-set.
  - [x] Preserve the truth that completion may win after cancellation was requested; represent unreachable owners and unproven cleanup as uncertain rather than signalling an unverified PID or claiming cancellation.
  - [x] Change server shutdown to drain new starts, persist `server_shutdown` cancellation only for workers owned by that server, wait for bounded cleanup/terminal persistence, and retain uncertainty when cleanup cannot be proven.
  - [x] Test first/duplicate/terminal cancel, persistence failure, owner acknowledgement, unreachable/stale owner, PID reuse, user interrupt, timeout, shutdown, and all cancellation/completion/failure/reconciliation races under the race detector.

- [x] **Task 7: Enforce Retention, Compaction, Quota, And Persistence Degradation**
  > Executes: Decision 7; Configuration; Performance; Persistence And Migrations; Errors; RO-2, RO-5, RO-9; AC-6, AC-12
  - [x] Preserve an independent current snapshot, last sequence, oldest retained sequence, history-complete flag, omission totals, and retention state for full, compacted, and tombstone records.
  - [x] Enforce 4,096 retained events and 16 MiB event bytes per run; compact expired terminal progress first, then convert expired full records to tombstones, then remove expired tombstones while preserving required warnings/findings/artifact references, cancellation facts, terminal facts, and explicit gaps.
  - [x] Start background compaction at 80 percent of quota, reject new acceptance at soft quota, reserve 16 MiB for active heartbeat/cancel/recovery/terminal writes, cancel owned active work at hard quota, and never delete an active snapshot.
  - [x] Checkpoint WAL and vacuum incrementally outside acceptance/append transactions; bound compaction and reconciliation batches so foreground work remains observable.
  - [x] Retry event append for at most 5 seconds before cancelling owned work, cancel if heartbeat cannot commit before lease expiry, and retry terminal commit on a cleanup context for 30 seconds before later conservative reconciliation.
  - [x] Refuse acceptance/mutation on corruption or unsupported schema, preserve evidence for support, and never report stale in-memory success when repository reads fail.
  - [x] Test deterministic compaction order, replay lower-bound advances, active preservation, tombstone expiry, soft/hard quota/headroom, total disk exhaustion, busy timeout, permission loss, corruption, and failure at acceptance/append/heartbeat/cancel/terminal boundaries.

- [x] **Task 8: Add Shared Run Queries And CLI Commands**
  > Executes: Decisions 8-10; CLI Surface; Errors; Observability; RO-3, RO-4, RO-7, RO-9; AC-2, AC-3, AC-7
  - [x] Expose app capabilities for bounded list, detail, replay/follow, cancellation, health/diagnostics, and active projection using run-control DTOs rather than web or TUI types.
  - [x] Add `ultraplan run list`, `show`, `follow`, `cancel`, and `diagnostics` with documented flags, help, exit/error mapping, stable text, and JSON output matching canonical lifecycle/liveness/cursor/recovery facts.
  - [x] Support list filters and opaque pagination with newest-first canonical order, default limit 50, maximum 200, and active lifecycle defined only as accepted/queued/running/cancelling.
  - [x] Ensure interrupting `run follow` stops observation only; cancellation remains the explicit `run cancel` command and follows local workspace/OS authority checks.
  - [x] Add an explicit bounded support-export mode through the diagnostics surface, with allowlisted contents, private output, deterministic size limits, and no prompts/provider payloads/source content/credentials/unsafe paths/arbitrary output.
  - [x] Integrate run-control preflight and degraded health without blocking unrelated runtime-free commands unless their requested run-control read/mutation requires the store.
  - [x] Add strict help, text, JSON, pagination/filter, unknown-value, error, diagnostics, support, and CLI/app agreement tests.

- [x] **Task 9: Add Canonical Run HTTP/SSE APIs And Preserve Operation Compatibility**
  > Executes: Decisions 6, 8; Errors; Security; Persistence And Migrations; RO-3 through RO-8; AC-2 through AC-7, AC-9
  - [x] Add `GET /api/v1/runs`, `GET /api/v1/runs/{id}`, `GET /api/v1/runs/{id}/events`, and authorized idempotent `DELETE /api/v1/runs/{id}`, with HEAD behavior matching existing router policy and existing v1 success/error envelopes.
  - [x] Implement explicit canonical DTOs for full, compacted, tombstone, lifecycle, liveness, attempts, progress, cursor boundaries, cancellation, immutable terminal, diagnostics, safe correlations, product-status projection, and recovery links.
  - [x] Accept `Last-Event-ID` or `after`, reject conflicting or malformed cursors, return typed `cursor_ahead`, and return JSON `409 replay_gap` with requested/oldest/last boundaries, reason, current snapshot, and recovery choices before opening SSE.
  - [x] Poll/replay committed SQLite events and keep only bounded transient connection queues in `internal/web`; heartbeat comments remain unsequenced transport data and another server must continue from the same durable cursor.
  - [x] Make safe run/operation reads workspace-visible across browser sessions while retaining loopback, Host/Origin, same-origin, path, no-CORS, and output-redaction protections; require a current session plus CSRF for browser cancellation.
  - [x] Keep `POST /api/v1/operations` as the confirmed start command, return `202` only after durable acceptance/claim, use the run ID as the operation ID, preserve the exact frozen JSON document, and add a canonical run `Link` header.
  - [x] Project active list/detail/events/cancel from the durable repository into the exact existing operation field order, lifecycle mappings, event names, route/method matrix, and error envelope; remove originating-session read filtering without weakening mutation authorization.
  - [x] Resolve durable operation aliases, compacted runs, tombstones, and recognized pre-durable `op_*` IDs; return typed `410 legacy_operation_not_retained` and recovery guidance instead of fabricating history or showing generic `Operation not retained`.
  - [x] Extend health additively with a safe `run_control` summary and map store/gap/tombstone/owner/cancel/migration/reconciliation failures to stable codes without substring classification.
  - [x] Update route, exact old-schema, strict old-client, embedded-browser bundle, session rotation, CSRF/Origin, cross-server dedupe, cursor, tombstone, legacy, unavailable-store, and SSE compatibility fixtures.

- [x] **Task 10: Build Server-Rendered Run List And Detail Pages**
  > Executes: Decision 9; CLI Surface; Documentation; Errors; Security; Performance; RO-3 through RO-5; AC-2 through AC-7
  - [x] Add canonical `/runs` and `/runs/{run_id}` HTML routes and resolve/redirect `/operations/{id}` to run detail, tombstone, or precise legacy recovery.
  - [x] Add explicit view models and namespaced primitives/components/layouts/pages for lifecycle/liveness badges, active summary, filters, run rows/cards, identity, progress, attempts, bounded timeline, diagnostics/recovery, cancellation, gaps, tombstones, and unknown states.
  - [x] Render complete no-JavaScript list/detail snapshots. Render at most 200 events initially, page earlier retained events explicitly, and retain at most 500 timeline rows in the enhanced DOM without confusing presentation pruning with repository loss.
  - [x] Render lifecycle separately from liveness and operational result separately from product status; preserve unknown totals/usage and show stalled, interrupted, cleanup-uncertain, store-unavailable, cancellation-uncertain, compacted, tombstone, and legacy states truthfully.
  - [x] Update the top bar from workspace active lifecycle on load/focus and every 5 seconds while visible; pause while hidden and retain the last count with unavailable/stale labeling rather than replacing failure with zero.
  - [x] Hand off the server-rendered last sequence to dependency-free JavaScript, deduplicate by run/sequence, follow committed SSE events, preflight/refetch on connection errors, and expose connecting/live/reconnecting/gap/capacity/store/terminal states.
  - [x] Keep cancellation an authorized explicit action with confirmation and focus recovery; never optimistically set cancelled, cancel on navigation/disconnect, or add run-level retry semantics.
  - [x] Implement URL-backed filters, mobile cards below 720 px, zoom/narrow layouts, keyboard/focus behavior, throttled polite live summaries, non-color cues, reduced motion, safe text-only event insertion, and no-JS navigation/forms.
  - [x] Add server-rendering, no-JS, hostile-content, responsive, accessibility, bounded-DOM, hidden-tab, cursor-handoff, repeated 512-event catch-up, gap, terminal-refresh, cancellation-race, and two-server browser-engine tests.

- [x] **Task 11: Make TUI Observation Durable And Cross-Surface Consistent**
  > Executes: Decisions 3, 6, 9; CLI Surface; Performance; RO-3 through RO-5, RO-7; AC-2, AC-3, AC-6, AC-7
  - [x] Add workspace run list/detail/follow/cancel views/actions over shared app capabilities and the same lifecycle, liveness, cancellation, gap, omission, product-status, and recovery vocabulary as CLI/JSON/browser.
  - [x] Replace the current model-owned operation identity/terminal assumptions with the accepted durable run ID and current snapshot; keep bounded Bubble Tea message delivery as presentation state only.
  - [x] Track the last committed sequence, deduplicate replay, and refresh from the durable repository after a full/dropped local channel rather than silently losing progress.
  - [x] Preserve explicit guarded cancellation and ensure quitting/hiding a run view cancels observation only unless the user confirms the cancellation command.
  - [x] Add TUI model/view tests for active/terminal/stalled/uncertain/gap/tombstone/unknown states, dropped delivery recovery, duplicate events, session-independent detail, and label/result agreement fixtures shared with app/web.

- [x] **Task 12: Add Correlated Telemetry, Diagnostics, And Support Evidence**
  > Executes: Decision 10; Observability; Errors; Security; RO-9; AC-11, AC-12
  - [x] Standardize safe structured fields for request, run, attempt, operation/target, owner/fence, process, runtime/agentwrap, sequence, lifecycle transition, cancellation, reconciliation, persistence operation, and terminal winner across app/run-control/web logs.
  - [x] Add bounded local counters/histograms for acceptance/append/terminal latency/failure, active/stalled runs, lease renewal, reconciliation backlog/age, retention/compaction, replay gaps, subscriber lag/drop, and cancellation routing; never use run/attempt/target IDs as metric labels.
  - [x] Report repository/schema/WAL/quota/compaction health, stale ownership, reconciliation backlog, persistence degradation, and cancellation uncertainty through app health and `run diagnostics` text/JSON.
  - [x] Generate a bounded redacted support bundle containing allowlisted snapshots, event headers/omission facts, health, safe config-source facts, reconciliation evidence, and logs; test size, permissions, escaping, and redaction across every included record.
  - [x] Defer OpenTelemetry export and a broad Prometheus endpoint; keep the field vocabulary exporter-neutral.

- [x] **Task 13: Build The Deterministic Failure And Cross-Process Verification Matrix**
  > Executes: Decisions 1-10; Testing; RO-10; AC-1 through AC-13
  - [x] Add deterministic fake clock, ID source, owner/process probe, repository failure hooks, notifier/subscriber speed controls, and fake product/runtime operations without weakening real SQLite contract coverage.
  - [x] Add helper-process integration tests over one temporary local workspace for concurrent CLI owners and two loopback web observers, including session rotation, observer restart, retained replay followed by new events, and identical active/detail projections.
  - [x] Kill owners before/after acceptance, claim, append, heartbeat, cancellation acknowledgement, and terminal commits; distinguish observer loss from owner loss and assert grace, fencing, interruption/uncertainty, and no adoption.
  - [x] Exercise stale leases, clock jumps, exact PID reuse, stale fence writes, duplicate cancellation, cancel/completion/failure/timeout/reconciliation races, server shutdown, slow subscribers, replay capacity, compaction, quota, disk full, permission loss, corruption, busy locks, migration, backup, restore, and unsupported schema.
  - [x] Prove redaction and bounds against hostile prompts/payload-shaped content, credentials, absolute paths, oversized output, malformed runtime events, unsafe diagnostics, HTML/JSON/SSE injection, metrics, and support export.
  - [x] Prove existing sprint/study locks, flow state, execute state, session checkpoints, cleanup markers, artifacts, Git state, and smoke evidence remain authoritative and are never imported as invented run history or rewritten by run control.
  - [x] Run package, full, and race suites repeatedly enough to detect leaked goroutines, stale-writer success, flaky timing assumptions, incompatible DTOs, and polling/subscriber races.
  - [x] **Stop condition:** any timing-only, PID-only, in-memory-only, provider-required normal test, or test that infers success from artifacts must be replaced with deterministic authority/failure evidence before release work begins.

- [x] **Task 14: Complete Documentation, Reviews, Build, And Gated Dogfood**
  > Executes: Decision 10; Documentation; Testing; RO-10, RO-11; AC-13
  - [x] Update governed Architecture, PRD, and TRD through their authorized workflow with the selected `internal/runcontrol`, same-host/local-filesystem boundary, direct SQLite writers, no-adoption rule, exact lifecycle/liveness/retention/failure policies, and product-authority separation.
  - [x] Update implementation `docs/architecture.md`, `docs/cli-reference.md`, `docs/local-web.md`, `docs/configuration.md`, `docs/recovery.md`, `docs/web-compatibility-baseline.md`, `docs/user-guide.md`, and `docs/release-checklist.md` for run commands/APIs/pages, security, cursor/gap/tombstone behavior, telemetry, support export, migration, backup/restore, rollback, disk recovery, and unsupported multi-host topology.
  - [x] Document the immutable terminal/cancellation race model, owner lease/fencing/process-birth diagnosis, active-count definition, event sampling/retention bounds, stale/store-unavailable UI, and separate operational/product outcomes with concrete recovery commands.
  - [x] Record the pinned SQLite dependency, supported platform/filesystem evidence, measured polling/write/build impact, binary size, migration backup location/limits, and all intentional compatibility behavior.
  - [x] Run the normal tests, race tests, focused process/browser/migration suites, CLI build, and documentation checks; retain command output or CI links as execution evidence.
  - [/] Run Architecture Review and Sprint Review with the authority trace, failure matrix, API compatibility, security/redaction, product-state separation, and verification results. Deferred by explicit user direction on 2026-08-21; review preflight was ready, but no review artifact was produced.
  - [/] After a current acceptable review, run the Deep Smoke Sprint protocol with one real CLI runtime run observed/replayed from two local servers and either completed or explicitly cancelled/reconciled; unavailable runtime/browser/platform prerequisites produce a truthful blocked result, never a pass. Deferred by explicit user direction on 2026-08-21; no smoke command or artifact was produced.
  - [x] Do not promote later content, QA/repair, retrieval, authored-product persistence, graph, cloud, daemon, broker, remote-worker, frontend-framework, WebSocket, or OpenTelemetry work from this sprint.

## Asynchronous Entry-Point Coverage

| Entry point | Durable boundary before work | Failure evidence |
| --- | --- | --- |
| CLI sprint `flow`, `verify`, `execute`, `review`, and `smoke` | Synchronous app command acceptance/claim wraps the product call; every runtime child is additionally decorated. | `TestSprintFlowNonDryRunUsesConfiguredRuntime` proves a top-level `sprint-flow` run and zero additional runtime calls when the repository cannot open; controlled-runtime closed-store coverage proves child fail-closed behavior. |
| CLI study `run-loop`, `run-all`, `run`, and `synthesize` | Synchronous app command acceptance/claim wraps service/runtime initialization and supplies the owner-cancelled context; individual runtime calls remain decorated. | App package execution and controlled-runtime persistence-failure tests exercise the shared boundary. |
| Web confirmed operation | The hub consumes the current confirmation and calls the durable manager before record insertion or goroutine creation. | Durable manager tests close the repository before acceptance; hub tests prove start ordering and cross-server alias reuse. |
| TUI confirmed operation | The Tea update accepts/claims before returning the operation command that creates execution work. | TUI model tests exercise durable identity, explicit cancellation, and hiding/quitting as observation-only behavior. |
| Runtime/agentwrap child | `controlledRuntime.StartRun` accepts/claims before invoking the wrapped runtime and commits sanitized events before callbacks. | `TestControlledRuntimeDoesNotStartWhenAcceptancePersistenceFails` and commit-before-delivery tests use a child spy. |

Direct CLI invocations intentionally use a fresh run rather than retaining an
idempotency alias. Confirmed web starts retain only the SHA-256 digest of the
canonical request/fingerprint token and resolve that alias through SQLite
across local server processes.

## Evidence Checklist

- [x] Every asynchronous CLI, TUI, and web entry point has durable acceptance/claim evidence and no-child-on-persistence-failure coverage.
- [x] Tests prove stable IDs, lifecycle invariants, immutable terminal arbitration, lease/fence behavior, process-birth safety, and conservative reconciliation.
- [x] Real SQLite tests prove concurrent writers, sequence ordering, transaction durability, WAL/locking, permissions, migration, backup, integrity, restore, and rollback behavior.
- [x] Event evidence proves pre-append redaction/bounds, commit-before-delivery, at-least-once deduplication, replay/cursor gaps, coalescing facts, slow-subscriber recovery, and bounded storage.
- [x] Cross-process evidence proves two CLI runs and two local servers share active counts, inspection, replay, live follow, cancellation, and terminal results.
- [x] Compatibility evidence freezes old operation routes, DTO field order, lifecycle mappings, SSE names, error envelope, strict clients, and browser bundle while validating canonical run APIs.
- [x] Security evidence proves workspace-readable sanitized history, fresh browser cancellation authority, loopback/Host/Origin/CSRF policy, private storage, and no unsafe signaling or content exposure.
- [x] Browser/TUI evidence covers no-JS and enhanced states, responsive/accessibility behavior, bounded rendering, session/observer restart, unknown values, tombstones, gaps, and uncertainty.
- [x] Persistence fault evidence covers acceptance, claim, append, heartbeat, cancellation, terminal, quota, disk-full, permissions, busy locks, corruption, and unsupported schema without silent loss.
- [x] Product-authority tests prove run control does not replace or rewrite Markdown, flow/execute/study state, locks, checkpoints, artifacts, Git/source, or smoke evidence.
- [x] Structured logs, additive health, local metrics, diagnostics JSON, reconciliation records, and bounded redacted support export exist and share safe correlation.
- [x] Documentation, migration/rollback/recovery runbooks, topology limitations, API/CLI references, and release notes are complete.
- [/] Architecture Review, Sprint Review, and Deep Smoke Sprint protocol evidence is current; explicitly deferred by user direction, with no review or smoke artifact produced.
- [x] Deviations from `reasoning.md` are recorded through governed reasoning and planning before implementation continues; no implementation deviation required a reasoning change.

## Verification Commands

Commands run from `../ultraplan-go` unless noted otherwise. Focused test names are implementation deliverables; they must exist and pass before their task closes.

| Check | Command | Expected Result |
| --- | --- | --- |
| Run-control unit and SQLite contracts | `go test ./internal/runcontrol` | Domain, repository, lifecycle, lease, event, retention, migration, diagnostics, and fault contracts pass with real temporary SQLite where required. |
| Process integration | `go test ./internal/runcontrol -run '^TestProcess'` | Concurrent owners/observers, owner death, process identity, fencing, replay, quota, and restore scenarios pass deterministically. |
| App and CLI acceptance/surface contracts | `go test ./internal/app` | Every asynchronous command is durably accepted, CLI run commands agree with app DTOs, and persistence failure starts no work. |
| Web API, compatibility, SSE, security, and rendering | `go test ./internal/web` | Canonical run routes/pages and frozen operation contracts pass, including session rotation, gaps, tombstones, redaction, and two-server observation. |
| Browser-engine run scenarios | `go test ./internal/web -run '^TestBrowserRun'` | Responsive/accessibility, no-JS/enhanced replay, cancellation, reconnect, bounded DOM, and cross-server browser scenarios pass. |
| TUI projection and replay | `go test ./internal/tui` | Durable identity, cursor recovery, bounded delivery, cancellation, and cross-surface labels agree. |
| Migration, backup, restore, and rollback fixtures | `go test ./internal/runcontrol -run 'Test(Migration|Backup|Restore|Rollback)'` | Supported upgrades succeed atomically; unsupported/corrupt states fail visibly; restore returns the matching schema/binary state. |
| Full deterministic suite | `go test ./...` | All normal tests pass without real-provider requirements. |
| Full race suite | `go test -race ./...` | Concurrent writers, cancellation, terminal, reconciliation, polling, and subscriber paths pass without data races or stale-writer success. |
| CLI build | `go build ./cmd/ultraplan` | The single CLI binary builds with the pinned pure-Go SQLite dependency. |
| CLI run help and JSON smoke | `go run ./cmd/ultraplan run --help` | Documented run subcommands and flags render without opening a runtime. |
| Architecture review | `ultraplan sprint ultraplan-go 35-durable-run-observability review` | Review applies the selected Architecture and Sprint Review protocols and records a current acceptable result or actionable findings. |
| Gated real-runtime dogfood | `ultraplan sprint ultraplan-go 35-durable-run-observability smoke` | After an acceptable review, external harness evidence proves CLI-to-two-server observation/replay/recovery, or records a truthful blocked result. |

## Risks And Blockers

| Risk / Blocker | Source | Mitigation | Status |
| --- | --- | --- | --- |
| Workspace filesystem does not provide trustworthy SQLite WAL/locking semantics. | `reasoning.md#assumptions-and-risks` | Preflight/health rejection, real local-filesystem contract tests, documented same-host boundary, and stop rather than fallback. | open |
| `modernc.org/sqlite` build, FULL-sync, migration, or platform behavior is unsuitable. | `reasoning.md#assumptions-and-risks` | Pin/review the driver; test Linux/macOS builds, concurrency, integrity, corruption, backup, and restore before integration closes. | open |
| A runtime-backed command bypasses shared app acceptance. | `reasoning.md#assumptions-and-risks` | Maintain an explicit entry-point inventory and command-path no-child tests; block release until coverage is complete. | open |
| Process-birth identity is unavailable or uncertain on a supported host. | `reasoning.md#assumptions-and-risks` | Keep liveness/cancellation conservative, never signal by PID alone, expose uncertainty, and add platform-specific probes/tests. | open |
| Polling across servers/tabs creates contention or unacceptable latency. | `reasoning.md#assumptions-and-risks` | Indexed bounded reads, 250 ms catch-up/1-second idle backoff, 5-second top-bar polling, subscriber limits, metrics, and dogfood measurements. | open |
| Full filesystem exhaustion consumes reserved headroom and prevents terminal persistence. | `reasoning.md#assumptions-and-risks` | Quota admission, reserved writes, bounded retries, owner cancellation, conservative reconciliation, and an external disk-recovery runbook. | open |
| Mid-run persistence loss cancels otherwise healthy expensive work. | `reasoning.md#assumptions-and-risks` | Five-second append retry, explicit diagnostics, durable omission/failure facts where possible, and product-owned resume without invented success. | accepted |
| Workspace-readable sanitized metadata is still locally sensitive. | `reasoning.md#assumptions-and-risks` | Preserve loopback/same-origin policy, field allowlists, relative paths, private DB permissions, pre-storage redaction, and fresh mutation auth. | open |
| Canonical run and frozen operation projections drift. | `reasoning.md#assumptions-and-risks` | Central mapping tables, strict old-client/browser fixtures, explicit terminal flag, and CLI/TUI/web agreement tests. | open |
| Tombstone/dedup expiry surprises old links or retried starts. | `reasoning.md#assumptions-and-risks` | Typed `410`, explicit retention boundaries, run/product recovery links, and no fabricated history. | accepted |
| Operational and product terminal outcomes disagree after a fault. | `reasoning.md#assumptions-and-risks` | Present both authorities separately, preserve immutable operational facts, and route repair/resume through product commands. | accepted |
| Fixed timings or retention defaults produce false stalls or insufficient diagnosis. | `reasoning.md#assumptions-and-risks` | Fake-clock tests, 45-second grace, omission facts, metrics, and evidence-gated future tuning without adding unsafe configuration now. | open |
| Browser-engine or real-runtime prerequisites are unavailable. | `reasoning.md#expected-evidence` | Keep normal tests provider-free; make browser/runtime gates explicit and record blocked evidence rather than weakening or skipping acceptance claims. | open |

## Review Inputs

Review should use:

- `sprint-index.md`
- `technical-handbook.md`
- `reasoning/api-design.md`
- `reasoning/architecture.md`
- `reasoning/frontend.md`
- `reasoning.md`
- this `plan.md`
- implementation diff
- verification evidence
- migration/backup/restore and support-bundle fixtures
- process/browser/gated-runtime evidence
- `system/protocols/architecture-review-protocol.md`
- `system/protocols/review-sprint-protocol.md`
- `system/protocols/deep-smoke-sprint-protocol.md`

## Execution Log

| Date / Step | Action | Evidence / Notes |
| --- | --- | --- |
| 2026-08-21 / planning | Created the implementation plan from validated requirements, index, handbook, area reasoning, final reasoning, project docs/roadmap, code context, and focused live composition/config/dependency inspection. | Plan only; no implementation, review, smoke, issue, Git, or run-state work was executed. |
| 2026-08-21 / Task 1 | Added the adapter-independent `internal/runcontrol` domain and direct multi-process SQLite repository contract; pinned `modernc.org/sqlite v1.57.0`. | Real temporary-database tests prove opaque identity, private modes, WAL/FULL/foreign-key/5-second busy policy, schema/indexes, immutable events, concurrent sequence allocation, lock isolation, reopen integrity, and one terminal winner. `go test ./...`, package race, vet, and Linux/macOS amd64/arm64 cross-builds pass. Current CLI binary grew 3,480 bytes (26,708,820 to 26,712,300) from module graph changes before app wiring; cross-built runcontrol test binaries are 11.8-12.3 MiB and final linked CLI impact must be remeasured after Task 3. Sprint 34 real-runtime dogfood remains a user-accepted missing prior-evidence limitation and is not represented as pass evidence. |
| 2026-08-21 / Task 2 | Added one-way schema migration/version checks, advisory migration exclusion, WAL checkpoint, bounded timestamped private backups, integrity validation, and a stopped-process restore fixture; extended effective config/source reporting with the three selected retention/quota settings. | Migration, backup, restore, lock-contention, unsupported-schema, configuration precedence, redaction, and unsafe-bound tests pass. Soft quota and reserved headroom remain fixed derivations of the configured hard quota rather than additional settings. The new repository remains empty and does not inspect or rewrite product-owned records. |
| 2026-08-21 / Tasks 3-7 | Composed one repository/owner per process; added direct CLI, web, TUI, and runtime-child acceptance; fenced control/cancellation/reconciliation; sanitized commit-before-delivery events; immutable terminal arbitration; retention/compaction/quota and conservative persistence degradation. | Direct sprint command coverage proves the top-level durable run and fail-closed runtime start. Fake-clock/birth-token/fence/terminal tests, hostile/oversize event tests, helper-process concurrent sequence allocation and pre-claim owner-exit reconciliation, deterministic retention/quota tests, stale migration-lock recovery, and corrupt-database evidence preservation pass. An acceptance left unclaimed is interrupted after the 45-second grace with no fabricated attempt/PID evidence. Linux uses `/proc` birth/boot identity; Darwin cross-builds use `kern.proc.pid` and `kern.boottime`. The exhaustive fault matrix remains open. |
| 2026-08-21 / Tasks 8-12 | Added shared list/detail/replay/cancel/health use cases; `ultraplan run` commands and private support export; canonical run JSON/SSE and no-JS pages; cross-session/cross-server durable operation compatibility; bounded JavaScript replay; and durable TUI list/detail/cancel views. | App/web/TUI package suites, strict API compatibility fixtures, hostile HTML escaping, canonical cursor errors, cross-server missing-hub projection, JavaScript syntax validation, focused `TestBrowserRun` coverage, TUI lifecycle/liveness/retention fixtures, and repository-recovery polling pass. Support evidence now includes allowlisted config source classes and sanitized reconciliation decisions. Production-wide log plumbing and the complete browser/TUI matrices remain open. |
| 2026-08-21 / Verification checkpoint | Ran normal and race suites, vet, focused process/migration/browser suites, CLI help/build, Linux arm64 and Darwin amd64/arm64 CLI builds, and Darwin amd64/arm64 run-control test cross-builds. | `go test ./...`, `go test -race ./...`, `go vet ./...`, focused commands, JavaScript syntax validation, diff whitespace validation, and all builds passed after the final implementation edits. Current Linux CLI size is 33,596,517 bytes versus the pre-app-wiring 26,712,300-byte checkpoint (+6,884,217 bytes, predominantly the pure-Go SQLite link). |
| 2026-08-21 / review-readiness checkpoint | Added `execute.md` as a truthful in-progress execution record and ran the execution validator, plan validator, and review dry-run only. | `validate plan` passes. Post-checkpoint `validate execute` rejects already-checked executable tasks by design. Review dry-run now blocks only because unresolved top-level tasks lack authorized deferred outcomes. No `review.md`, `smoke.md`, issue, Git, or external-runtime artifact was created. |
| 2026-08-21 / completion | Completed the ten reopened implementation tasks, reconciled all 14 top-level tasks to complete, and reran normal/race/focused/build/documentation checks. | Plan validation and review dry-run preflight pass. The user explicitly directed that actual review and smoke not run; those two gates are marked `[/]`, and no `review.md`, `smoke.md`, or pass verdict is claimed. |

## Completion Criteria

- [x] All tasks are complete or explicitly deferred through a governed reasoning/plan update.
- [x] Every asynchronous entry point accepts and claims a durable run before work starts, and required persistence failure starts no child.
- [x] CLI, JSON, TUI, dashboard, run list/detail, and two local servers agree on active, lifecycle, liveness, cursor/gap, cancellation, terminal, and recovery facts.
- [x] Fencing, process birth, reconciliation, cancellation, terminal races, retention, quota, corruption, migration, and rollback satisfy the complete deterministic failure matrix.
- [x] Frozen operation compatibility and additive canonical run APIs/pages are test-backed across session/observer changes and legacy/tombstone states.
- [x] Redaction, bounds, private storage, fresh mutation authorization, telemetry, diagnostics, and support export satisfy the security and observability contracts.
- [x] Existing product artifacts, locks, state, checkpoints, Git/source, and smoke evidence remain separate and authoritative for their owned concerns.
- [x] `go test ./...`, `go test -race ./...`, focused process/browser/migration suites, and `go build ./cmd/ultraplan` pass or have explicit governed blockers.
- [x] Documentation and migration/rollback/recovery guidance are complete.
- [/] Architecture Review and Sprint Review are current and acceptable, and gated Deep Smoke Sprint evidence passes or records a truthful external-prerequisite blocker. Explicitly deferred by user direction; no result is claimed.
- [/] `review.md` can evaluate conformance without guessing intent. No `review.md` was produced because review execution was explicitly declined.
