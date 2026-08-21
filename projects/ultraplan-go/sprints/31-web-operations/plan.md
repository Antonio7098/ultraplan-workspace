# Sprint Plan: Guarded Web Operations and SSE Progress

> Project: `ultraplan-go`
> Sprint: `31-web-operations`
> Source: `reasoning.md`
> **Inputs Used:** `projects/ultraplan-go/project-index.md`, `projects/ultraplan-go/roadmap.md`, `projects/ultraplan-go/sprints/31-web-operations/requirements.md`, `projects/ultraplan-go/docs/ARCHITECTURE.md`, `projects/ultraplan-go/docs/PRD.md`, `projects/ultraplan-go/docs/TRD.md`, `projects/ultraplan-go/sprints/31-web-operations/sprint-index.md`, `projects/ultraplan-go/sprints/31-web-operations/technical-handbook.md`, `projects/ultraplan-go/sprints/31-web-operations/reasoning/api-design.md`, `projects/ultraplan-go/sprints/31-web-operations/reasoning/architecture.md`, `projects/ultraplan-go/sprints/31-web-operations/reasoning/frontend.md`, `projects/ultraplan-go/sprints/31-web-operations/reasoning.md`, `../ultraplan-go/docs/plans/server-shutdown-run-cancellation-contract.md`, `../ultraplan-go/internal/app/operations.go`, `../ultraplan-go/internal/app/app.go`, `../ultraplan-go/internal/app/sprint_commands.go`, `../ultraplan-go/internal/app/study_commands.go`, `../ultraplan-go/internal/app/tui_commands.go`, `../ultraplan-go/internal/app/web_usecases.go`, `../ultraplan-go/internal/sprint/service.go`, `../ultraplan-go/internal/sprint/verification_lock.go`, `../ultraplan-go/internal/study/locks.go`, `../ultraplan-go/internal/study/run_loop.go`, `../ultraplan-go/internal/web/routes.go`, `../ultraplan-go/internal/web/handlers.go`, `../ultraplan-go/internal/web/security.go`, `../ultraplan-go/internal/web/server.go`, `../ultraplan-go/cmd/ultraplan/main.go`

This plan executes `reasoning.md`. It does not reopen architecture, scope, routes, states, limits, lock domains, confirmation semantics, cancellation ownership, persistence, or frontend technology.

## Reasoning Source

- **Sprint Reasoning:** `reasoning.md`
- **Sprint Index:** `sprint-index.md`
- **Technical Handbook:** `technical-handbook.md`
- **Area Reasoning:** `reasoning/api-design.md`, `reasoning/architecture.md`, `reasoning/frontend.md`
- **Normative Shutdown Contract:** `../ultraplan-go/docs/plans/server-shutdown-run-cancellation-contract.md`
- **Evidence Scope:** The selected final study reports were not reopened because `technical-handbook.md`, the three area reasoning documents, and `reasoning.md` contain sufficient cited evidence for every implementation decision. Implementation files were inspected only to map decisions to current seams.

## Sprint Status

- **Status:** `partial — implementation and verification complete; shutdown-deadline deviation and governed review remain`
- **Owner:** `implementation agent`
- **Start Date:** `2026-08-16`
- **Completion Date:** `2026-08-16`

## Decisions To Execute

| Decision | Source Section | Requirement | Evidence And Trade-Off | Rejected Alternative | Execution Implication / Risk |
| --- | --- | --- | --- | --- | --- |
| One Closed App Operation Capability And Thin Web Adapter | `reasoning.md#decision-1-one-closed-app-operation-capability-and-thin-web-adapter` | `AC-03`, `AC-16`, `AC-17`; `C-01`, `C-02` | Thin-adapter and constructor-injection evidence supports extending the existing typed app capability; explicit variants add wiring but keep authority auditable. | Direct product imports, CLI parsing/subprocesses, app job manager, generic registry. | Evolve `internal/app/operations.go`; converge CLI/TUI/web on it; keep web transport-only. Stop if an existing operation cannot be shared without changing workflow semantics. |
| Binding Prepare/Start Confirmation And Versioned Operation API | `reasoning.md#decision-2-binding-preparestart-confirmation-and-versioned-operation-api` | `AC-01`-`AC-04`, `AC-14`; `C-04`, `C-10` | Post-resolution validation and explicit trust-boundary evidence justify a two-minute, session-bound, single-use preparation record; repeated normalization costs work but keeps consent current. | Direct start, browser-authored summaries, stateless replayable tokens, preparation locks, web queue. | Implement the fixed `/api/v1/operations` methods, strict DTOs, current SHA-256 fingerprint checks, stable errors, and no accepted-work queue. |
| Product-Owned Mutation Exclusion And Conservative Recovery | `reasoning.md#decision-3-product-owned-mutation-exclusion-and-conservative-recovery` | `AC-09`, `AC-11`-`AC-13`; `C-05`, `C-09` | Typed lock and bounded-cleanup evidence supports one per-sprint mutation lease; reduced same-sprint concurrency is accepted for state safety. | HTTP-only locking, global mutex, optimistic overlap, PID-only stale unlock. | Add `internal/sprint/locks.go`; reuse existing study locking; hold ownership through durable reconciliation or explicit uncertainty. |
| Server-Owned Lifecycle, Exact-Once Cancellation, And Ordered Shutdown | `reasoning.md#decision-4-server-owned-lifecycle-exact-once-cancellation-and-ordered-shutdown` | `AC-06`, `AC-09`-`AC-12`; `C-07`-`C-09` | Root-context and bounded-shutdown evidence supports one owner and one cleanup path; shutdown may wait to its deadline and report uncertainty. | Disconnect cancellation, fire-and-forget work, fail-fast global cleanup, early lock release, unbounded background cleanup. | Add draining, exact-once cancellation, terminal arbitration, deadline escalation, terminal-before-stream-close ordering, and startup reconciliation. |
| Bounded Ephemeral Hub And Progress-Only SSE | `reasoning.md#decision-5-bounded-ephemeral-hub-and-progress-only-sse` | `AC-04`, `AC-05`, `AC-07`, `AC-08`; `C-03`, `C-06`, `C-07` | Bounded pub/sub evidence supports producer-independent queues and explicit replay gaps; transient progress may be lost and must defer to durable state. | Blocking or unbounded subscribers, silent drops, durable event log, WebSockets, global event IDs. | Implement the exact limits, event names, per-operation IDs, heartbeat, replay-gap recovery, retention, and load-shedding decisions. |
| Stable Safe Error, Result, Security, And Observability Projections | `reasoning.md#decision-6-stable-safe-error-result-security-and-observability-projections` | `AC-02`, `AC-07`, `AC-14`, `AC-16`; `C-01`, `C-04`, `C-10` | Typed-error, redaction, and stable-field evidence supports allowlisted DTO duplication; compatibility maintenance is accepted to prevent disclosure and drift. | Error-string matching, raw model serialization, encoding-only redaction, accepted failures as transport failures, permissive CORS. | Project safe data before retention; preserve Sprint 30 security; separate accepted outcomes, API errors, SSE, and logs while correlating IDs. |
| Server-Rendered Operation Views With Narrow Progressive Enhancement | `reasoning.md#decision-7-server-rendered-operation-views-with-narrow-progressive-enhancement` | `AC-01`, `AC-04`, `AC-07`, `AC-08`, `AC-15`; `C-07`, `C-11` | Safe projection and bounded rendering evidence supports server-owned views plus narrow JavaScript; live cancellation requires JavaScript. | SPA/client store, generic command console, persistent browser state, unbounded console, unload cancellation, POST method override. | Extend owning detail pages, render all lifecycle/recovery states, use one bounded `EventSource`, and provide CLI cancellation guidance without JavaScript. |
| Deterministic Cross-Surface Verification, Required Reviews, And Documentation | `reasoning.md#decision-8-deterministic-cross-surface-verification-required-reviews-and-documentation` | `AC-18`; `C-12`, `C-13`; all test/documentation `RO-*` | Behavior-first fakes and selective goldens make lifecycle races reproducible; broad deterministic test infrastructure is accepted. | Live providers in normal tests, all-golden coverage, private map/goroutine assertions, unit-only proof. | Produce focused package, race, full-suite, build, documentation, Architecture Review, and Sprint Review evidence. |

## Requirements / Contracts To Satisfy

| Contract / Requirement ID | Required Behavior | Evidence Planned |
| --- | --- | --- |
| `AC-01`, `AC-02`, `C-04` | Side-effect-free prepare returns normalized scope, relative paths, mutation class, effective runtime/model or harness summary, prerequisites, current fingerprint, two-minute expiry, and a bound single-use token; start revalidates all inputs. | Per-variant canonicalization/fingerprint fixtures; no-write/no-lock/no-runtime assertions; expiry, mismatch, replay, stale, and session tests. |
| `AC-03`, `AC-16`, `AC-17`, `C-01`, `C-02` | Only existing allowlisted operations are exposed through one typed app capability; web imports no product/runtime/process/CLI handler and existing CLI/TUI behavior remains supported. | App dispatch tests, CLI/TUI/web agreement fixtures, compile/import review, full regression suite. |
| `AC-04`, `AC-07`, `AC-08`, `C-07` | POST starts commands, DELETE requests cancellation, and SSE observes only; monotonic typed events, heartbeat, reconnect, gaps, slow readers, and browser disconnect cannot control work. | Route/method tests and focused SSE concurrent publication, replay, rollover, heartbeat, slow-reader, disconnect, and terminal-flush tests. |
| `AC-05`, `C-03`, `C-06` | The operation hub is bounded, ephemeral, has no queue or durable schema, and recovers through product state after eviction/restart. | Capacity, aggregate-bound, retention, eviction, restart, and no-web-store tests. |
| `AC-06`, `AC-09`-`AC-12`, `C-08`, `C-09` | The server owns each accepted operation; explicit/shutdown cancellation reaches canonical cleanup; draining, bounded shutdown, escalation, lock ordering, and conservative startup recovery remain truthful. | Exact-once/race tests, multi-operation shutdown barriers, nested runtime/harness/process fakes, durable state fixtures, stale-lock restart tests. |
| `AC-13`, `C-05` | Sprint flow/execute/review/smoke/verify share product-owned per-sprint exclusion and study run-loop retains its existing product lock. | Lock acquire/release/conflict tests across CLI/TUI/web, cleanup barriers, study coexistence, actionable safe diagnostics. |
| `AC-14`, `C-10` | Stable API and accepted-operation outcomes are safe, typed, correlated, redacted before retention, loopback-only, same-origin, CSRF/session protected, body-limited, timeout-bounded, and containment-safe. | Error matrix and compatibility fixtures; hostile secret/path/stderr/provider tests across app, hub, JSON, SSE, HTML, and logs; Sprint 30 security regression tests. |
| `AC-15`, `C-11` | Server-rendered pages and narrow JavaScript show confirmation, progress, findings, all terminal states, shutdown, gaps, connection loss, and durable recovery without client-owned truth. | Template/static semantic tests, representative goldens, hostile rendering, keyboard/focus/live-region/reduced-motion checks, bounded DOM and one-stream tests. |
| `AC-18`, `C-12`, `C-13` | Deterministic focused/full/race/build verification passes; no automatic Git mutation occurs. | Recorded focused tests, `go test ./...`, `go test -race ./...`, `go build ./cmd/ultraplan`, review protocol evidence. |
| Required documentation outputs | `docs/local-web.md`, `docs/cli-reference.md`, and `docs/architecture.md` describe the implemented API, limits, ownership, shutdown, trust, reconnect, cancellation, and recovery contracts. | Documentation diff checked against public routes, event/error/state names, numerical limits, and both required review protocols. |

## Tasks

- [x] **Task 1: Stabilize The Closed Shared App Operation Capability**
  > Executes: `Decision 1`, `Decision 2`, `Decision 6`; `AC-01`-`AC-03`, `AC-16`, `AC-17`
  - [x] Evolve `../ultraplan-go/internal/app/operations.go` into the separate closed capability for validation, prompt preview, dry run, sprint flow, execute, review, smoke, verify, and study run-loop without adding web/session/SSE concepts.
  - [x] Define explicit per-kind request normalization, affected relative paths, mutation class, prerequisites, effective runtime/model or harness summary, governed-input inventory, SHA-256 fingerprint, safe event projection, terminal result, durable refresh guidance, and typed error mapping.
  - [x] Keep preparation side-effect-free and make synchronous execution re-normalize and compare the expected fingerprint before mutation, runtime, or harness work.
  - [x] Refactor `internal/app/sprint_commands.go`, `study_commands.go`, and `tui_commands.go` only as needed so CLI and TUI exercise the same app paths and retain existing flags, output, confirmation, and automation semantics.
  - [x] Add `internal/app/web_operations_test.go` with per-variant normalization/fingerprint mutation fixtures, preparation side-effect assertions, safe projection/redaction, nested cancellation, and cross-surface agreement.
  - [x] **Stop condition:** Record a deviation before continuing if an existing operation cannot accept caller context/progress or cannot be shared without changing its workflow semantics.

- [x] **Task 2: Add Product-Owned Sprint Mutation Exclusion And Recovery Truth**
  > Executes: `Decision 3`, `Decision 4`; `AC-09`-`AC-13`, `C-05`, `C-08`, `C-09`
  - [x] Add `../ultraplan-go/internal/sprint/locks.go`, building on the existing verification-lock behavior, so sprint flow, execute, review, smoke, and verify share one exclusive per-sprint mutation lease across processes and surfaces.
  - [x] Acquire after start/fingerprint revalidation and before mutation or external execution; release only after durable terminal reconciliation or product-owned `cleanup_uncertain` recording.
  - [x] Preserve `internal/study/locks.go` as the study run-loop authority and map both lock domains to typed, actionable, safely projected conflicts without exposing lock internals.
  - [x] Extend only owning sprint/study state and compatibility validation needed to represent cancellation, interruption, cleanup uncertainty, and stale-running recovery; never add a web-owned recovery file or infer success from process absence/artifact presence.
  - [x] Add `internal/sprint/locks_test.go` and focused recovery fixtures for acquire/release, cross-surface conflict, cleanup-barrier ordering, cancellation/shutdown, stale owner, process evidence, compatibility, and study-lock non-overlap.
  - [x] **Stop condition:** If an existing durable schema cannot truthfully represent interruption/uncertainty compatibly, record the owner-specific schema question and migration impact before implementation proceeds.

- [x] **Task 3: Implement Binding Session And Confirmation Policy**
  > Executes: `Decision 2`, `Decision 6`; `AC-01`, `AC-02`, `AC-14`, `C-04`, `C-10`
  - [x] Extend `../ultraplan-go/internal/web/security.go` with a bounded preparation store, per-process secret/session ownership, independent CSRF validation, opaque tokens, single-use arbitration, injected clock/ID sources, and two-minute expiry.
  - [x] Bind each record to session, canonical normalized request, current fingerprint, and expiry; consume on mismatch/staleness/expiry or successful creation as decided, while preserving a token across pre-creation capacity/draining rejection.
  - [x] Preserve Sprint 30 loopback, exact Host/Origin, session cookie, security headers, containment, request timeout, and 64 KiB body policy on all new routes.
  - [x] Test strict JSON and unknown fields, caller-authored authority fields, CSRF/session separation, expiry, mismatch, stale input, replay, response-loss replay protection, preparation capacity, and hostile token/log fixtures.

- [x] **Task 4: Build The Bounded Ephemeral Operation Hub**
  > Executes: `Decision 4`, `Decision 5`, `Decision 6`; `AC-05`-`AC-12`, `C-03`, `C-06`-`C-09`
  - [x] Add `../ultraplan-go/internal/web/operations.go` with one server-derived root context, canonical exact-once cancel function, owner goroutine, terminal arbitration point, safe state/result, subscriber registry, timestamps, and short eviction lifecycle per accepted operation.
  - [x] Enforce no queue and the exact defaults/upper bounds: 8 active operations; 128 preparations/2 minutes; 256 events and 256 KiB per operation; 16 KiB encoded event; 256 KiB terminal result; 8 subscribers per operation; 32 streams server-wide; 32 queued events per subscriber; 10-minute terminal retention; 15-second heartbeat; 30-minute stream lifetime.
  - [x] Serialize per-operation decimal event IDs and retain only already-redacted `snapshot`, `progress`, `warning`, `finding`, `artifact`, `cancel_requested`, `recovery_required`, and `terminal` projections.
  - [x] Make producer publication non-blocking; disconnect a full subscriber queue; keep operation, request, subscriber, and bounded-cleanup contexts distinct.
  - [x] Add `internal/web/operations_test.go` for capacity, aggregate bounds, exact-once cancellation, terminal races, retention/eviction, restart loss, slow/abandoned subscribers, cleanup, and absence of durable hub state.

- [x] **Task 5: Add Versioned Operation Routes, Results, Errors, And SSE**
  > Executes: `Decision 2`, `Decision 5`, `Decision 6`; `AC-01`-`AC-08`, `AC-14`, `AC-16`
  - [x] Extend `../ultraplan-go/internal/web/routes.go` and `handlers.go` with exactly `POST /api/v1/operations/prepare`, `POST /api/v1/operations`, `GET /api/v1/operations/{id}`, `GET /api/v1/operations/{id}/events`, and `DELETE /api/v1/operations/{id}`.
  - [x] Return `200` preparation, `202` creation with `Location`, `200` retained status/result, idempotent `202`/`200` cancellation, structured unknown-route/method responses with `Allow`, `429` capacity, and `503` draining behavior.
  - [x] Keep accepted product validation/runtime/timeout/cancellation/interruption/cleanup outcomes inside a `200` operation document; map pre-acceptance failures to the stable codes decided in `reasoning/api-design.md` using `errors.Is/As`, never strings.
  - [x] Implement SSE initial snapshot, retained replay after `Last-Event-ID`, explicit replay-gap `recovery_required` plus snapshot, heartbeat comments, terminal flush, session recheck, and bounded closure without command semantics.
  - [x] Add or extend `routes_test.go` and add `sse_test.go` for methods, envelopes, strict DTOs, compatibility with Sprint 30, every state/error code, concurrent monotonic IDs, replay/rollover, heartbeat, payload bounds, slow readers, disconnect isolation, and terminal ordering.

- [/] **Task 6: Integrate Ordered Server Shutdown And Startup Reconciliation** — implemented except the recorded owner-specific durable deadline-exhaustion write. — Deferred: owner-specific durable cleanup-uncertain persistence after an exhausted shutdown deadline requires a separately governed app/product capability; existing cancellation, process escalation, ephemeral uncertainty reporting, and startup reconciliation are complete.
  > Executes: `Decision 3`, `Decision 4`; `AC-06`, `AC-09`-`AC-12`, `C-08`, `C-09`
  - [x] Extend `../ultraplan-go/internal/web/server.go` and `../ultraplan-go/cmd/ultraplan/main.go` so the hub and app capability are explicitly constructed and server startup runs conservative product recovery before accepting mutations.
  - [x] On graceful shutdown, enter draining first, reject preparation/start, snapshot owners under lock, cancel each exactly once with `server_shutdown` outside locks, and keep bounded status/SSE reads available during cleanup.
  - [ ] Waits outside hub locks for canonical cancellation and existing platform process-tree escalation; the hub records `cleanup_uncertain` at deadline and startup reconciles dead owners, but owner-specific durable uncertainty is not written before exit when the deadline is already exhausted.
  - [x] Publish terminal outcomes, release product locks only through reconciliation, close SSE subscribers, then stop HTTP and return; do not use unbounded `context.Background()` cleanup or detached work.
  - [x] Add deterministic server/hub integration tests for idle and multi-operation shutdown, draining rejection, nested cancellation, process escalation, exact-once races, lock ordering, terminal-before-close, bounded deadline, and stale-running/lock restart reconciliation.

- [x] **Task 7: Add Server-Rendered Operation Views And Narrow Enhancement**
  > Executes: `Decision 7`; `AC-01`, `AC-07`, `AC-08`, `AC-15`, `C-07`, `C-11`
  - [x] Add allowlisted operation entry forms to the owning project, sprint, and study templates and explicit handler-built view models for normalized confirmation scope, relative paths, mutation/runtime/harness summary, prerequisites, and expiry.
  - [x] Add operation components/pages for `accepted`, `running`, `cancelling`, `succeeded`, `failed`, `cancelled`, `interrupted`, and `cleanup_uncertain`, including `server_shutdown`, connection loss, incomplete history, findings/artifacts, and durable recovery actions.
  - [x] Extend embedded dependency-free JavaScript for the same prepare/start DTOs, one `EventSource`, duplicate suppression, bounded reconnect/fallback, bounded DOM updates, and explicit `DELETE` cancellation; store no token, event history, result, or workflow truth persistently.
  - [x] Keep preparation/start/status/recovery server-rendered without JavaScript; provide explicit CLI cancellation guidance rather than a POST method override.
  - [x] Extend CSS/templates/static tests for escaped text and allowlisted URLs, native controls, focus/error behavior, live regions, reduced motion, color-independent states, bounded timeline, one stream, all lifecycle/reason/error fixtures, and enhanced/non-enhanced DTO parity.

- [x] **Task 8: Complete Safe Observability And Cross-Boundary Security Evidence**
  > Executes: `Decision 5`, `Decision 6`; `AC-07`, `AC-14`, `AC-16`, `C-01`, `C-06`, `C-10`
  - [x] Correlate request ID, operation ID, safe product run/task ID, operation kind/scope, state transition, cancellation reason, duration, event sequence, subscriber outcome, cleanup outcome, and typed error code without logging tokens, CSRF/cookies, raw payloads, unsafe paths, prompts, stderr, or provider data.
  - [x] Add local counters for starts, active/terminal outcomes, rejection/cancellation reasons, active streams, slow subscribers, replay gaps, projection drops, and shutdown cleanup; do not add a network metrics endpoint.
  - [x] Run hostile fixtures through app event/result mapping, hub retention/replay, JSON, SSE, HTML, and captured logs to prove redaction occurs before retention.
  - [x] Inspect `internal/web` imports and composition to prove it has no direct `internal/study`, `internal/project`, `internal/sprint`, runtime/process, or CLI-handler dependency and no alternate workflow engine.

- [x] **Task 9: Update Public And Architecture Documentation**
  > Executes: `Decision 8`; documentation `RO-*`
  - [x] Update `../ultraplan-go/docs/local-web.md` with operation kinds, prepare/confirm/start/status/events/cancel flow, exact limits, states/events/errors, browser disconnect, explicit cancellation, draining/shutdown, gaps/eviction, local trust limits, and durable recovery.
  - [x] Update `../ultraplan-go/docs/cli-reference.md` with web/API affordances and server-owned shutdown semantics without changing CLI automation behavior.
  - [x] Update `../ultraplan-go/docs/architecture.md` with the closed app capability, product lock/recovery ownership, ephemeral hub, context/cancellation ordering, progress-only SSE, and prohibition on detached/web-owned workflow state.
  - [x] Keep Sprint 32 compatibility, release accessibility audit, packaging, and gated real-runtime/browser hardening explicitly deferred.

- [/] **Task 10: Run Deterministic Verification And Required Reviews** — deterministic verification and manual architecture review are complete; governed independent Sprint Review remains downstream. — Deferred: the independent governed Sprint Review is the downstream review stage and must not be executed or fabricated within execute; deterministic verification and the manual architecture review are complete.
  > Executes: `Decision 8`; `AC-18`, `C-12`, `C-13`
  - [x] Run the focused app/TUI, sprint/study lock/recovery, and web operation/SSE/security/template suites with deterministic fakes, barriers, clocks, IDs, and temporary workspaces.
  - [x] Run package race tests, the full test suite, full race suite, and binary build; record commands and outcomes in execution evidence.
  - [x] Apply `system/protocols/architecture-review-protocol.md` to dependency direction, ownership, bounds, cancellation, recovery, and absence of a second workflow engine.
  - [ ] Apply `system/protocols/review-sprint-protocol.md` in the governed downstream review stage. Execute did not create or replace `review.md`, as required by the execute-stage mutation boundary.
  - [x] Record any deviation from `reasoning.md` before further implementation rather than silently changing routes, limits, lock domains, states, ownership, or scope.

## Evidence Checklist

- [x] Tests prove per-kind normalized confirmation and shared CLI/TUI/web behavior.
- [ ] Tests prove product lock ownership, exact-once cancellation, ordered bounded shutdown, platform process-tree escalation, and conservative restart recovery; owner-specific durable `cleanup_uncertain` persistence before an exhausted server deadline remains the deviation recorded below.
- [x] Tests prove all hub, payload, result, subscriber, stream, retention, request, and DOM bounds.
- [ ] JSON/SSE/template fixtures cover the implemented methods and representative lifecycle/error paths; exhaustive every-code compatibility coverage remains for governed review/Sprint 32 hardening.
- [x] Host/Origin/session/CSRF/body/header/path/redaction evidence covers all mutation and stream routes.
- [x] Runtime/diagnostic evidence contains safe correlated fields and local counters without unsafe retained data.
- [x] Documentation updates are complete and match implemented values.
- [x] Deviations from `reasoning.md` are recorded before implementation continues.
- [ ] Architecture protocol was applied manually during execute; the independent governed Sprint Review remains the next stage and must produce current `review.md` evidence.

## Verification Commands

Run from `../ultraplan-go`.

| Check | Command | Expected Result |
| --- | --- | --- |
| App and cross-surface operations | `go test ./internal/app ./internal/tui` | Shared operation normalization, dispatch, projection, cancellation, and CLI/TUI behavior pass. |
| Product locks and recovery | `go test ./internal/sprint ./internal/study` | Per-sprint/study exclusion, cleanup ordering, state compatibility, and recovery fixtures pass. |
| Web API, hub, SSE, security, and UI | `go test ./internal/web` | Route, confirmation, lifecycle, SSE, shutdown, redaction, template, and Sprint 30 regression tests pass. |
| Focused race gate | `go test -race ./internal/app ./internal/tui ./internal/sprint ./internal/study ./internal/web` | No data races in operation ownership, locks, subscribers, state, or projections. |
| Full suite | `go test ./...` | All repository packages pass with deterministic normal dependencies. |
| Full race suite | `go test -race ./...` | All repository packages pass under the race detector. |
| Binary build | `go build ./cmd/ultraplan` | The `ultraplan` binary builds with embedded web assets and explicit composition. |
| Web dependency review | `go list -deps ./internal/web` | Review shows app/standard dependencies only and no direct product/runtime/process/CLI-handler import. |

## Assumptions, Open Questions, And Stop Conditions

| Item | Type | Required Handling |
| --- | --- | --- |
| Existing app/product use cases propagate context and expose sufficient progress/results. | Assumption | Verify in Task 1; refactor command glue inward without duplicating semantics. Stop and record a deviation if a boundary cannot propagate cancellation. |
| Existing durable states can represent cancellation/interruption or accept a compatible owner-specific extension. | Assumption | Verify before schema edits in Task 2; add no web recovery state. Record migration/compatibility impact if false. |
| Preparation dependencies can inspect configuration and prerequisites without mutation. | Assumption | Prove no lock, write, runtime, harness, process, or capacity reservation occurs; move unsafe checks to start. |
| Exact internal Go type names may follow current conventions. | Open implementation detail | Naming may vary, but package ownership, closed request union, synchronous app execution, safe event sink, lock domain, API, and lifecycle decisions are fixed. |
| `requirements.md` contains 18 acceptance criteria and 13 constraints, while `reasoning.md` line 10 says `C-01` through `C-15`. | Traceability note | Use only the 13 constraints actually present in authoritative `requirements.md`; do not invent `C-14` or `C-15`. This does not alter any decided behavior. |
| A required route, state, event, limit, lock domain, or terminal representation proves incompatible with current product state. | Stop condition | Record the evidence and deviation/open question before changing the decided contract or continuing dependent tasks. |

## Risks And Blockers

| Risk / Blocker | Source | Mitigation | Status |
| --- | --- | --- | --- |
| Fingerprint omits a material option, governed input, runtime/model, or harness identity. | `reasoning.md#assumptions-and-risks` | Canonical request, effective runtime summary, `ultraplan.yml`, governed project/sprint/study inputs, and smoke preparation scope are hashed; mutation tests cover governed plan changes. | `mitigated` |
| Completion, timeout, explicit cancel, and shutdown race to conflicting outcomes. | `reasoning.md#assumptions-and-risks` | One hub terminal arbitration point, exact-once cancellation, deterministic barriers, and focused/full race tests. | `mitigated` |
| Cleanup releases a lock early or blocks server exit. | `reasoning.md#assumptions-and-risks` | Product leases span composite workflows; server waits outside hub locks and marks ephemeral cleanup uncertainty at deadline. Owner-specific durable uncertainty before exit is the deviation below. | `open deviation` |
| Unsafe data reaches memory before redaction. | `reasoning.md#assumptions-and-risks` | Safe projections redact before hub retention; hostile token/path/authorization fixtures pass. | `mitigated` |
| Fixed limits multiply into excessive memory or goroutines. | `reasoning.md#assumptions-and-risks` | Exact aggregate/owner/subscriber/event/result limits and local atomic counters are implemented; tuning remains measurement-led. | `mitigated` |
| Browser mistakes a gap, eviction, disconnect, or cancellation request for a terminal result. | `reasoning.md#assumptions-and-risks` | Distinct lifecycle/events, `recovery_required`, durable refresh links, explicit DELETE cancellation, and no-JavaScript status pages are implemented. | `mitigated` |
| Forced termination leaves live descendants or stale ownership. | `reasoning.md#assumptions-and-risks` | Existing platform runner owns and escalates process groups; server startup reconciles dead-owner execute/review/smoke attempts and preserves live leases. | `mitigated with deviation below` |
| No-JavaScript cancellation lacks browser parity because cancellation must use DELETE. | `reasoning.md#trade-off-and-debt-analysis` | Preserve status/recovery pages and provide explicit CLI cancellation guidance; do not add POST override. | `accepted` |
| Local-only deployment is mistaken for trusted input. | `reasoning.md#assumptions-and-risks` | Exact Host/Origin, signed session, independent CSRF, strict JSON/form binding, body/path/confirmation/redaction checks apply to commands and streams. | `mitigated` |
| Existing broad app capability grows into a god service. | `reasoning.md#potential-technical-debt` | Closed variants, one shared app runner, private product dependencies, and no registry or generic workflow engine. | `mitigated` |

## Recorded Deviations And Follow-Up

| Deviation | Reason / Impact | Required Follow-Up |
| --- | --- | --- |
| If the 10-second graceful-shutdown deadline expires while an app operation still has not returned, the hub publishes `cleanup_uncertain`, closes bounded transport state, and the next startup reconciles dead-owner execute/review/smoke state. It does not currently force an owner-specific durable `cleanup_uncertain` write before process exit. | The existing app operation result does not expose a safe owner-specific emergency reconciliation handle that can write while the canonical operation still owns its product lease. Writing from `internal/web` or racing the live product owner would violate the decided ownership/lock boundary. Normal cancellation already reaches runtime/harness process-group cleanup; this gap is the deadline-exhaustion edge. | The governed review must classify severity. A follow-up app/product capability should atomically persist owner-specific cleanup uncertainty under the existing lease before transport closure, with deterministic deadline-race tests. |
| Exhaustive compatibility fixtures for every stable error/event/state combination were not added; representative route, lifecycle, replay, SSE, security, redaction, template, and no-JavaScript paths are covered. | Core behavior and all verification gates pass, but the planned compatibility matrix is broader than the focused fixtures delivered in this direct execute run. | Governed review and Sprint 32 compatibility hardening must add the complete matrix without changing `/api/v1` meanings. |

## Deferred Scope And Revisit Triggers

| Deferred Item | Revisit Trigger |
| --- | --- |
| Durable workers, detached operations, durable idempotency, or durable event history | A separately governed/versioned architecture defines leases, heartbeats, worker identity, ownership transfer, authentication, reconciliation, and replay. |
| Finer sprint mutation lock classes | A reviewed read/write matrix proves safe concurrency and includes deadlock/recovery analysis. |
| Adaptive limits, buffer pooling, or metrics endpoint | Sprint 32 measurements or an explicit observability requirement justify the additional contract. |
| Browser/API compatibility hardening, release accessibility audit, real-runtime/browser evidence, and packaging | Sprint 32 begins or deterministic Sprint 31 proof cannot be completed without a gated boundary check. |
| Frontend framework, client router/store, service worker, WebSocket, terminal/chat transport | Demonstrated requirements and interaction complexity justify a separately reasoned architecture. |
| Hosted, remote, multi-user, arbitrary editing, issue tracking, automatic fixes, or Git mutation | A later product phase explicitly selects and governs the capability. |

## Review Inputs

Review should use:

- `sprint-index.md`
- `technical-handbook.md`
- `reasoning/api-design.md`, `reasoning/architecture.md`, and `reasoning/frontend.md`
- `reasoning.md`
- this `plan.md`
- implementation diff
- focused/full/race/build verification evidence
- documentation diff
- `system/protocols/architecture-review-protocol.md`
- `system/protocols/review-sprint-protocol.md`

## Execution Log

| Date / Step | Action | Evidence / Notes |
| --- | --- | --- |
| `2026-08-16 / planning` | Created the Sprint 31 implementation plan from validated governed inputs and mapped it to current implementation seams. | Plan stage only; no implementation, smoke, review automation, issue tracking, Git mutation, or run-state artifact was executed. |
| `2026-08-16 / app and product ownership` | Added normalized, fingerprint-bound app preparation; one shared TUI/web runner; per-sprint composite mutation leases; and dead-owner execute/review/smoke startup reconciliation. | `internal/app/operations.go`, `operation_runner.go`, `web_operations_test.go`, `internal/sprint/locks.go`, `locks_test.go`; no Git mutation. |
| `2026-08-16 / web operations` | Added session/CSRF-bound preparations, bounded operation hub, versioned prepare/start/status/SSE/DELETE routes, exact-once cancellation, draining shutdown, safe projections/counters/log correlation, and progressive plus server-rendered no-JavaScript operation views. | `internal/web/operations.go`, `operation_handlers.go`, security/routes/server/templates/static files and focused tests. Exact limits match `reasoning/api-design.md`. |
| `2026-08-16 / documentation` | Updated local-web, CLI reference, and architecture documentation for operations, trust, bounds, SSE, cancellation, shutdown, and durable recovery. | `docs/local-web.md`, `docs/cli-reference.md`, `docs/architecture.md`. Sprint 32 release hardening remains deferred. |
| `2026-08-16 / verification` | Ran focused package tests, focused race tests, full suite, full race suite, safe binary build, diff check, and direct dependency import review. | Passing commands: `go test ./internal/app ./internal/tui ./internal/sprint ./internal/study ./internal/web`; `go test -race ./internal/app ./internal/tui ./internal/sprint ./internal/study ./internal/web`; `go test ./...`; `go test -race ./...`; `go build -o /tmp/ultraplan-sprint31 ./cmd/ultraplan`; `git diff --check`; `go list -f '{{join .Imports "\n"}}' ./internal/web` (only stdlib plus `internal/app`). The build used `/tmp` to preserve the pre-existing ignored `ultraplan` binary. |
| `2026-08-16 / architecture review` | Applied the Architecture Review protocol manually. | Verdict: acceptable fit with justified lifecycle complexity; transport-only direct imports and product state ownership preserved. Required finding: owner-specific durable cleanup-uncertain persistence at an exhausted shutdown deadline remains unresolved. The independent governed Sprint Review is the next stage and was not generated during execute. |
| `2026-08-16 / live UI verification` | Started the fresh binary against the real workspace and fetched the Sprint 31 detail page over HTTP. Added conservative compatibility handling for legacy terminal run summaries and version-1 map-based flow state encountered during startup reconciliation. | Live response contained `Sprint operations`, execute/review/smoke/verify controls, and the operation status console. Legacy files remain byte-for-byte untouched; focused, full, and full-race suites pass. |
| `2026-08-16 / browser asset fix` | Diagnosed Brave CSS/JS requests returning `403` because opaque subresource origins were subjected to the read/API Origin policy. Exempted only known GET/HEAD `/static/` assets from Origin validation while retaining exact loopback Host enforcement and unchanged API/mutation Origin/CSRF policy. | Live checks: CSS and JS with `Origin: null` return `200`; API health with the same origin remains `403`. Focused web/app/sprint and web race tests pass. |

## Completion Criteria

- [x] All tasks are complete or explicitly deferred with requirement and risk impact.
- [x] Every required output path in `requirements.md` exists and is covered by evidence.
- [x] Verification commands were run successfully or any deferral is documented with a blocker and owner.
- [ ] Evidence satisfies all expectations except the two recorded deviations (deadline-exhaustion durable uncertainty and exhaustive compatibility fixtures); governed review must assess them.
- [x] CLI, TUI, and web agree on operation preparation, execution, cancellation, terminal truth, conflict, and recovery behavior.
- [x] `internal/web` remains transport-only and no durable web state, queue, detached work, new workflow semantics, or prohibited scope was introduced.
- [x] `review.md` can evaluate conformance without guessing intent.
