# Sprint Plan: Local Web Hardening and Observable-Product Release

> Project: `ultraplan-go`
> Sprint: `32-hardening-and-release`
> Source: `reasoning.md`
> **Inputs Used:** `projects/ultraplan-go/project-index.md`, `projects/ultraplan-go/sprints/32-hardening-and-release/requirements.md`, `projects/ultraplan-go/docs/ARCHITECTURE.md`, `projects/ultraplan-go/docs/PRD.md`, `projects/ultraplan-go/docs/TRD.md`, `projects/ultraplan-go/sprints/32-hardening-and-release/sprint-index.md`, `projects/ultraplan-go/sprints/32-hardening-and-release/technical-handbook.md`, `projects/ultraplan-go/sprints/32-hardening-and-release/reasoning/api-design.md`, `projects/ultraplan-go/sprints/32-hardening-and-release/reasoning/architecture.md`, `projects/ultraplan-go/sprints/32-hardening-and-release/reasoning/frontend.md`, `projects/ultraplan-go/sprints/32-hardening-and-release/reasoning.md`

This plan executes `reasoning.md`. It does not invent architecture, scope, or decisions. Requirements references use the `AC-*`, `C-*`, and `OUT-*` labels assigned by `reasoning.md`; implementation must preserve the underlying requirement text when an assigned label is ambiguous.

## Reasoning Source

- **Sprint Reasoning:** `reasoning.md`
- **Sprint Index:** `sprint-index.md`
- **Technical Handbook:** `technical-handbook.md`
- **Area Reasoning:** `reasoning/api-design.md`, `reasoning/architecture.md`, `reasoning/frontend.md`

## Sprint Status

- **Status:** `implemented; downstream review, smoke, and manual browser audit pending`
- **Owner:** `Codex manual execution`
- **Start Date:** `2026-08-18`
- **Completion Date:** `2026-08-18` (execute stage)

## Decisions To Execute

| Decision | Source Section | Execution Implication |
| --- | --- | --- |
| Enforce A Thin Web Boundary And Shared Capability Model | `reasoning.md#decision-1-enforce-a-thin-web-boundary-and-shared-capability-model` | Inventory boundary violations first; move workflow and durable-state interpretation behind cohesive typed app capabilities; keep `internal/web` limited to app and standard-library dependencies; prove future-stage exposure without route branching or a plugin registry. |
| Freeze And Safely Map The Existing `/api/v1` Contract | `reasoning.md#decision-2-freeze-and-safely-map-the-existing-apiv1-contract` | Establish the documented Sprint 30-31 baseline before DTO changes; preserve routes, methods, envelopes, fields, nullability, statuses, codes, content types, and cache rules; map typed app errors and make start/cancel retries safe. |
| Validate A Fail-Closed Local Security And Resource Policy Before Serving | `reasoning.md#decision-3-validate-a-fail-closed-local-security-and-resource-policy-before-serving` | Merge configuration in the established precedence order, validate bind/security/resource invariants before listening, enforce browser/path/content boundaries, and apply type-level plus final-projection redaction. |
| Keep Operations And SSE Bounded, Ephemeral, And Retry-Safe | `reasoning.md#decision-4-keep-operations-and-sse-bounded-ephemeral-and-retry-safe` | Bound confirmations, operations, events, results, subscribers, streams, retention, heartbeat, lifetime, polling, and cleanup; atomically consume confirmations and publish deduplicated operations; isolate slow/disconnected subscribers and require durable refresh after gaps or expiry. |
| Make Shutdown And Restart Reconciliation Truthful And Product-Owned | `reasoning.md#decision-5-make-shutdown-and-restart-reconciliation-truthful-and-product-owned` | Centralize terminal arbitration and exact-once cancellation; drain before rejecting mutations; perform bounded cleanup outside hub locks; persist an authoritative terminal or owner-specific uncertainty before closing; reconcile restart state conservatively. |
| Ship A Server-First, Accessible, Embedded Presentation System | `reasoning.md#decision-6-ship-a-server-first-accessible-embedded-presentation-system` | Parse and validate one embedded namespaced template tree; use typed view models and downward composition; provide complete no-JavaScript behavior, layered CSS, disposable enhancement state, accessibility, and outside-tree packaging. |
| Gate Release With Layered Deterministic, Review, Documentation, And Real-System Evidence | `reasoning.md#decision-7-gate-release-with-layered-deterministic-review-documentation-and-real-system-evidence` | Keep test layers distinct, reconcile docs with implementation and fixtures, run required build/race/review protocols, and record unavailable real-runtime or harness prerequisites as blocked rather than passed. |

## Requirements / Contracts To Satisfy

| Contract / Requirement ID | Required Behavior | Evidence Planned |
| --- | --- | --- |
| Architecture; `AC-08`-`AC-11`; `C-01`; app/web required outputs | `internal/web` is a thin adapter over shared typed app capabilities, imports no product/runtime/process/CLI packages, and uses explicit view models and downward templates. | Import-boundary check, handler delegation tests, fake future-stage capability fixture, template parse/render tests, Architecture Review. |
| CLI Surface; `AC-01`, `AC-02`, `AC-06`, `AC-22` | Preserve `ultraplan serve`, numeric-loopback behavior, single-binary operation, and CLI/TUI/web state agreement. | Help/config/build checks, cross-surface fixtures, outside-tree binary launch, synchronized CLI documentation. |
| Configuration; `AC-02`, `AC-16`; `C-04`, `C-10` | Apply defaults, workspace config, environment, and flags in order; reject unsafe or incoherent security/resource values before serving. | Table-driven precedence/default/source tests and startup-failure tests for every invalid bound or combination. |
| Errors; `AC-03`, `AC-04`, `AC-15` | Preserve stable typed API codes and status mappings; unknown API paths and unsupported methods remain JSON; no raw causes leak to any projection. | Exhaustive route/method fixtures, `errors.Is`/`errors.As` mapping tests, hostile error and forbidden-value corpus. |
| Security; `AC-02`, `AC-14`, `AC-15`; `C-03`, `C-05` | Enforce numeric loopback, Host/Origin, CSRF/session, CSP/headers, body/time/path/reference/Markdown controls, cache policy, and layered redaction. | IPv4/IPv6 and browser-security matrix, path/reference/Markdown tests, captured HTML/JSON/SSE/log/retention scans, manual audit. |
| Workflows; `AC-07`, `AC-08`, `AC-17`-`AC-20`; `C-06`-`C-09` | Use shared guarded operations, product-owned mutation exclusion and recovery, exact-once cancellation, disconnect-is-subscription behavior, and truthful terminal arbitration. | Representative study/sprint journeys, deterministic cancellation/terminal races, shutdown/restart reconciliation tests. |
| Persistence And Migrations; `AC-06`, `AC-18`, `AC-19`; `C-02`, `C-08` | Keep workspace/run state authoritative and atomic; add no web durability; record terminal or cleanup uncertainty before shutdown releases ownership. | Atomic-state assertions, forced-interruption and abrupt-restart fixtures, proof that process/artifact presence does not imply success. |
| Observability; `AC-06`, `AC-15`, `AC-18`, `AC-23` | Emit safe structured request, operation, event-gap, cancellation, cleanup, and recovery diagnostics without retaining forbidden payloads. | Captured diagnostic field assertions and full forbidden-value scans across responses, events, results, and logs. |
| Performance; `AC-05`, `AC-16`, `AC-17`, `AC-20`; `C-04` | Bound every request, operation, confirmation, payload, event, subscriber, stream, retention, polling, and cleanup path; isolate slow clients. | Capacity/overflow/lifetime/slow-client tests, race suite, goroutine/subscriber cleanup checks, representative workload validation. |
| Testing; `AC-01`, `AC-03`-`AC-24`; `C-11`-`C-12` | Use deterministic normal tests with temporary workspaces and fakes, separate compatibility/security/lifecycle/presentation/integration/package layers, and gate real dependencies. | Full and focused Go test suites, race run, build, package run, review records, and real-system pass or truthful blocked evidence. |
| Documentation; `AC-22`, `AC-23`; all documentation outputs | Document the actual command, routes, fields, limits, security scope, accessibility, shutdown/recovery states, packaging, and troubleshooting. | Documentation-to-fixture/config checklist covering every required guide and public contract. |
| LLM Runtime; LLM Evaluation / Cost / Safety; `AC-07`, `AC-15`, `AC-23`; `C-06`, `C-12` | Use shared app/runtime ownership, show safe runtime/model metadata before confirmation, preserve bounded cancellation/cleanup, and do not expose raw provider material. | Fake runtime/harness browser journeys, metadata/redaction assertions, gated real runtime and smoke-harness operation or blocked record. |

## Planning Boundaries And Open Questions

- The existing documented Sprint 30-31 API baseline is not enumerated in the planning artifacts. Task 1 must produce a reviewed route/field/method/status/code/cache inventory from current implementation and public documentation before any compatibility fixture or DTO change is accepted. An unexplained mismatch stops dependent API work.
- Current route-specific workflow logic and shutdown-reconciliation code locations are intentionally unresolved. Task 1 inventories them; discovery may resize Tasks 2 and 6 but cannot change the decided dependency direction or product ownership.
- Static asset URLs must use revalidation unless Task 1 proves that a URL is content-addressed. No immutable cache policy may be assigned based on filename convention alone.
- Concrete resource defaults are not selected by the evidence documents. Task 4 must preserve existing documented values where safe, test them against representative workflows, and record any proposed change as an explicit compatibility/configuration decision before implementation continues.
- `reasoning.md` assigns `C-01` through `C-14`, while the current `requirements.md` Constraints section contains fewer individually listed bullets. Implementation and review must trace labels to the quoted requirement behavior and record any label normalization rather than silently guessing.
- Real runtime and smoke-harness availability is environmental. Missing executable, provider/model, credentials, network, or harness prerequisites must be recorded as `blocked` with the exact missing prerequisite; deterministic release checks remain mandatory.

## Tasks

- [x] **Task 1: Inventory The Compatibility Baseline, Boundaries, And Release Inputs**
  > Executes: Decisions 1, 2, and 7; Architecture, Errors, Configuration, Documentation; `AC-03`-`AC-06`, `AC-08`-`AC-11`, `AC-22`, `AC-23`
  - [ ] Inventory every existing HTML, static, `/api/v1`, and SSE route in `internal/web`, including methods, request DTOs, success/error envelopes, field types/nullability, status/code mapping, content type, cache policy, confirmation/cancellation behavior, and documented examples.
  - [ ] Compare the inventory with Sprint 30-31 documentation and tests; classify each mismatch as implementation defect, documentation defect, or explicit unresolved compatibility question. Do not create or regenerate compatibility goldens while any mismatch lacks a rationale.
  - [ ] Inventory imports and direct behavior in `internal/web`: product/runtime/process/CLI dependencies, filesystem/state interpretation, stage-name branching, prompt/verdict/recovery logic, executable construction, and mutation-lock ownership.
  - [ ] Inventory operation ownership, goroutine/timer/event-source launch sites, locks and acquisition order, cancellation entry points, terminal writers, cleanup/reconciliation paths, and startup stale-state handling.
  - [ ] Inventory templates, template definitions/dependencies, page fallback forms, CSS/JavaScript assets, static URL versioning, source-relative asset access, and accessibility states.
  - [ ] Inventory existing web configuration sources/defaults and all body, request, operation, preparation, event, result, subscriber, stream, retention, heartbeat, lifetime, polling, and cleanup limits.
  - [ ] Record the baseline manifest and implementation findings in test fixtures or implementation-owned documentation selected during execution; do not add a new durable product registry, web database, or planning artifact.
  - [ ] **Dependency:** none; this is the first implementation task.
  - [ ] **Evidence:** reviewed baseline matrix, boundary/import inventory, lifecycle ownership/lock-order inventory, asset/cache inventory, and configuration-limit inventory referenced by subsequent task changes.
  - [ ] **Stop condition:** stop DTO, fixture, cache, or default-limit changes if current implementation and documented behavior disagree without an explicit compatibility decision.

- [x] **Task 2: Enforce Shared App Capabilities And The Thin Web Boundary**
  > Executes: Decision 1; Architecture, Workflows, CLI Surface, LLM Runtime, Testing; `AC-06`-`AC-11`, `C-01`, `C-06`; `OUT-Web application capability`, `OUT-Shared operation capability`
  - [ ] In `internal/app/web_usecases.go` and `internal/app/operations.go`, expose cohesive typed capabilities for status/readiness, findings, artifacts, available commands, normalized preparation, execution, progress, cancellation, terminal results, and durable recovery/next action.
  - [ ] Move each Task 1 workflow, filesystem/state interpretation, runtime/process access, stage branch, mutation-lock, terminal reconciliation, and recovery violation from `internal/web` into its owning app or product module without generalizing transport DTOs into domain types.
  - [ ] Keep app results surface-neutral and keep web request/response DTOs and page view models inside `internal/web`; use small consumer-oriented interfaces only where multiple surfaces or deterministic fakes require substitution.
  - [ ] Wire capabilities explicitly at `cmd/ultraplan`; add no package-global runner, reflection container, context service locator, dynamic production registry, or broad `WebService`/manager object.
  - [ ] Add `internal/app/web_usecases_test.go` coverage with a fake future-stage capability that exposes the complete fixed vocabulary through existing generic web behavior without a stage-specific route branch or production registration mechanism.
  - [ ] Add a direct import-boundary check and handler delegation tests proving `internal/web` imports only `internal/app` and standard-library packages and never invokes the `ultraplan` binary.
  - [ ] **Dependency:** Task 1 boundary inventory and compatibility baseline.
  - [ ] **Evidence:** focused app tests, import-boundary output, handler tests, fake-stage capability fixture, and cross-surface typed-result comparisons.
  - [ ] **Stop condition:** stop and record a plan deviation if a required use case cannot fit cohesive shared capabilities without a broad service locator, plugin system, or route-owned orchestration.

- [x] **Task 3: Freeze The `/api/v1` Contract And Typed Transport Mapping**
  > Executes: Decision 2; Errors, Security, Testing, Documentation, Observability; `AC-03`-`AC-05`, `AC-15`; `OUT-Stable API routes and envelopes`, `OUT-HTTP and page handlers`, `OUT-API compatibility fixtures`
  - [ ] Stabilize route and method registration in `internal/web/routes.go`, including structured JSON errors for unknown `/api/` paths and wrong methods with no HTML fallback.
  - [ ] Map validated transport requests and typed app results in `internal/web/handlers.go` and `internal/web/operation_handlers.go`; preserve the reviewed baseline fields, meanings, types, nullability, status mappings, stable codes, and accepted asynchronous statuses.
  - [ ] Define a small caller-actionable typed error mapping using `errors.Is`/`errors.As`; prevent raw causes, error strings, app/domain structs, and internal fields from defining the wire contract.
  - [ ] Apply `no-store` to state-bearing pages, API responses, confirmations, operation projections, errors, and SSE; apply static revalidation unless Task 1 proved a content-addressed URL.
  - [ ] Make operation start retry-safe through confirmation-to-operation deduplication and cancellation idempotent without allowing late cancellation to replace an authoritative terminal result.
  - [ ] Build `internal/web/api_compatibility_test.go` as an exhaustive reviewed route/method matrix that fails on field addition/omission, type/nullability, status, code, content type, and cache-policy changes; include empty/null cases, wrong methods, and unknown API paths.
  - [ ] Require a written contract rationale for every fixture change and pair exact fixtures with semantic tests for deduplication, cancellation, safe errors, operation expiry, and durable refresh.
  - [ ] **Dependency:** Tasks 1 and 2.
  - [ ] **Evidence:** reviewed compatibility fixtures and semantic API tests with an explicit rationale for every baseline difference.
  - [ ] **Stop condition:** any unexplained difference from the Task 1 baseline blocks release and must not be blessed by golden regeneration.

- [x] **Task 4: Fail Closed On Configuration, Browser Security, Paths, And Redaction**
  > Executes: Decision 3; Configuration, Security, Performance, Observability, LLM Evaluation / Cost / Safety; `AC-02`, `AC-14`-`AC-16`, `C-03`-`C-05`, `C-10`; `OUT-Browser security policy`
  - [ ] Resolve built-in defaults, workspace configuration, environment, and command flags into one immutable effective server configuration while preserving source/precedence diagnostics without secrets.
  - [ ] Validate numeric IPv4/IPv6 loopback binding and every positive/coherent request, operation, preparation, event, result, subscriber, stream, retention, heartbeat, lifetime, polling, and cleanup bound before opening the listener.
  - [ ] Harden `internal/web/security.go` for strict Host and exact same-origin Origin checks, CSRF and per-process session behavior, safe cookie attributes for actual loopback HTTP, CSP/security headers, strict content types, body and timeout limits, and ambiguous-header rejection.
  - [ ] Reject malformed and duplicate-significant JSON, unsupported operation kinds, caller-supplied internals, raw paths, executable strings, forged references, path escapes, and oversized preview/content input; resolve opaque or allowlisted references through app/product authority.
  - [ ] Render Markdown through the approved safe path with executable HTML/scripts disabled and keep browser-visible values escaped and allowlisted.
  - [ ] Apply secret-safe values where owned inward and final projection-level redaction independently to HTML, JSON, SSE, retained events, terminal results, and diagnostics; use constant-time comparison for bearer-like session/confirmation secrets.
  - [ ] Add `internal/web/security_test.go` coverage for IPv4/IPv6, Host/Origin, CSRF/session, CSP/headers, cache, body/smuggling cases, hostile JSON/Markdown/HTML/URLs, path/reference forgery, and a forbidden corpus containing tokens, cookies, CSRF values, confirmations, prompts, unsafe paths, provider payloads, stderr, arguments, and representative secrets.
  - [ ] Verify invalid or incoherent configuration prevents startup with safe field-level diagnostics rather than warnings, silent correction, or fallback.
  - [ ] **Dependency:** Task 1 configuration inventory and Task 2 app authority boundaries; coordinate wire errors with Task 3.
  - [ ] **Evidence:** configuration precedence/startup suite, complete security regression matrix, and forbidden-value scan across all projections and captures.
  - [ ] **Stop condition:** do not weaken a hard bound or fail-closed check to make a representative workflow pass; record and review a proposed safe default change first.

- [x] **Task 5: Bound Confirmations, Operations, Events, Subscribers, And SSE Recovery**
  > Executes: Decision 4; Performance, Workflows, Observability, Persistence And Migrations, LLM Runtime; `AC-05`, `AC-07`, `AC-16`, `AC-17`, `AC-20`, `C-02`, `C-04`, `C-06`, `C-07`; `OUT-Operation HTTP handlers`, `OUT-Operation and SSE hub`, `OUT-SSE enhancement`
  - [ ] In `internal/web/operations.go`, retain only bounded operation IDs, safe normalized summaries, confirmation/deduplication records, canonical cancel functions, redacted events, subscriber queues, and short-lived terminal projections; add no durable web-owned state.
  - [ ] Bind short-lived confirmations to normalized request, scope, mutation class, governed-input fingerprint, and session; define one critical ordering for expiry/staleness checks, confirmation consumption, capacity reservation, deduplication record creation, and operation publication.
  - [ ] Give every operation, goroutine, timer, event source, subscriber, reconnect/poll loop, and cleanup path one explicit owner, hard bound, and terminal cleanup path.
  - [ ] Assign monotonically increasing per-operation event IDs, bounded replay, stable heartbeat comments, explicit non-terminal replay-gap signals, bounded terminal delivery, and durable HTTP refresh after gap, expiry, terminal, or reconnect exhaustion.
  - [ ] Evict slow subscribers without blocking product execution, persistence, terminal arbitration, or shutdown; browser disconnect, tab close, refresh, navigation, and stream expiry cancel subscriptions only.
  - [ ] Keep SSE observation-only and keep preparation/start/cancel as HTTP commands; ensure query parameters and event cursors cannot mutate work.
  - [ ] Add `internal/web/operations_test.go` and `internal/web/sse_test.go` coverage using deterministic clocks, IDs, barriers, and subscriber controls for bounds, expiry, deduplication, ordering, replay, rollover, gaps, heartbeat, slow clients, disconnects, terminal flush, capacity errors, cancellation races, and cleanup.
  - [ ] Run focused race and leak checks against concurrent starts, duplicate submissions, cancellations, event publication, subscriber eviction, operation expiry, and terminal arbitration.
  - [ ] **Dependency:** Tasks 2-4.
  - [ ] **Evidence:** deterministic operation/SSE suites, focused race output, and no-leak/no-backpressure assertions.
  - [ ] **Stop condition:** stop if any subscriber path can block app work or if retry safety cannot be proven under concurrent duplicate start and shutdown.

- [x] **Task 6: Close The Shutdown And Restart Reconciliation Gap**
  > Executes: Decision 5; Architecture, Workflows, Persistence And Migrations, Errors, Observability, LLM Runtime; `AC-18`-`AC-20`, `C-06`, `C-08`, `C-09`; `OUT-Server lifecycle`, `OUT-Shared operation capability`
  - [ ] Propagate one server-owned operation work context through app, product, runtime, and process work; prohibit request/SSE ownership and fresh detached backgrounds in normal work.
  - [ ] Centralize terminal arbitration for completion, failure, timeout, user cancellation, and shutdown so one authoritative terminal wins and canonical cancellation is invoked at most once with an explicit reason.
  - [ ] Define and test lock ordering across hub state, product mutation locks, persistence, subscribers, process cleanup, and reconciliation; snapshot under lock and perform no app callback, I/O, send, wait, or cancellation invocation while holding hub locks.
  - [ ] In `internal/web/server.go`, enter draining, reject new preparations/starts/mutations with the stable typed response, snapshot active cancellations, request `server_shutdown` cancellation exactly once, and wait only for bounded cleanup.
  - [ ] Use a distinct purpose-specific bounded cleanup context for process-tree cleanup, product-owned durable reconciliation, and lock ownership resolution; do not use it to continue normal product work.
  - [ ] Persist an authoritative completion/failure/cancellation/interruption or owner-specific `cleanup_uncertain` equivalent before closing HTTP/SSE or releasing unresolved ownership; preserve an already-authoritative result when cancellation arrives late.
  - [ ] Reconcile stale active state and locks at startup through owning product modules; process absence, HTTP completion, SSE delivery, and artifact presence remain evidence only and never imply success.
  - [ ] Add `internal/web/server_test.go` and supporting app/product tests for multiple active operations, draining rejection, exact-once reasoned cancellation, deadline exhaustion, terminal races, no-I/O-under-lock, process-tree cleanup metadata, atomic terminal/uncertainty writes, forced interruption, abrupt restart, stale locks, and orphan detection.
  - [ ] **Dependency:** Tasks 2, 4, and 5; this task owns the Sprint 31 carry-forward gap.
  - [ ] **Evidence:** lifecycle/reconciliation suite, durable state fixtures for clean and uncertain outcomes, focused race/leak checks, and safe captured lifecycle diagnostics.
  - [ ] **Stop condition:** release remains blocked if deadline exhaustion can exit without durable owner-specific uncertainty or if restart can infer success from process/artifact observations.

- [x] **Task 7: Build The Namespaced Server-Rendered Template Hierarchy**
  > Executes: Decision 6; Architecture, Security, Documentation, Testing; `AC-10`-`AC-15`, `AC-21`; `OUT-Template primitives`, `OUT-Template components`, `OUT-Template layouts`, `OUT-Route pages`
  - [ ] Establish behavior-preserving render assertions for current pages and no-JavaScript actions before moving markup into the required template layers.
  - [ ] Build `templates/primitives/primitives.html`, `templates/components/components.html`, `templates/layouts/layouts.html`, and `templates/pages/pages.html` with stable names such as `primitive/*`, `component/*`, `layout/*`, and `page/*` and strict `page -> layout -> component -> primitive` dependencies.
  - [ ] Parse one embedded template tree at startup and fail on missing, duplicate, invalid, cyclic/upward, or route-unresolvable definitions rather than relying on parse order.
  - [ ] Construct explicit route/page view models in handlers; pass no app/domain objects or `map[string]any`; keep template functions presentation-only and free of filesystem, app, validation, mutation, and workflow interpretation.
  - [ ] Render complete dashboard, detail, finding, artifact, empty, confirmation, progress, terminal, error, stale/gap, draining, cancellation, and recovery states without JavaScript.
  - [ ] Use semantic landmarks, one descriptive page heading, native links/buttons, labels and error associations, table/list semantics, stable focus targets, restrained status/alert regions, and text-labelled non-color state distinctions.
  - [ ] Add `internal/web/templates_test.go` coverage for names/dependencies/startup failures, typed models, every representative page/state, no-JavaScript forms/actions, safe Markdown, hostile text/URL/attribute escaping, labels/landmarks/headings/status regions, and embedded paths.
  - [ ] **Dependency:** Tasks 2-4 for typed results, stable routes, and safe projection rules.
  - [ ] **Evidence:** parse/startup tests, render fixtures with reviewed changes, semantic accessibility assertions, hostile-content tests, and no-JavaScript journey coverage.
  - [ ] **Stop condition:** do not complete extraction if baseline forms/routes/fields or no-JavaScript behavior regress without an explicit compatibility rationale.

- [/] **Task 8: Layer CSS And Add Disposable Progressive Enhancement** — Deferred: implementation and deterministic semantic checks are complete, but manual keyboard, assistive-announcement, color, motion, 200% zoom, and narrow-reflow checks require an installed interactive browser; no Chromium or Firefox executable is available in this environment.
  > Executes: Decisions 4 and 6; Security, Performance, Testing, Documentation; `AC-11`-`AC-17`, `AC-20`, `AC-21`; all CSS/JavaScript required outputs
  - [ ] Implement `tokens.css`, `base.css`, `primitives.css`, `components.css`, `layouts.css`, and narrowly scoped `utilities.css`; keep product-state meaning in semantic component styles rather than utilities.
  - [ ] Provide visible tokenized focus, text and control contrast, non-color state cues, reduced-motion behavior, adequate targets, mobile-first reflow, deliberate long-content/table handling, and operability at narrow widths and 200% zoom.
  - [ ] Keep `app.js`, `operations.js`, and `sse.js` dependency-free and capability-focused with no router, persistent store, service worker, framework, build pipeline, or authoritative client state.
  - [ ] Enhance prepare/confirm/start/cancel and SSE only after the server-rendered baseline; keep one abort owner for requests, streams, retries, and polling and stop all loops on navigation, terminal state, expiry, or lifetime exhaustion.
  - [ ] Preserve the last truthful server-rendered state while pending; manage stable accessible labels and deterministic focus transitions for validation, confirmation, start, cancellation, terminal refresh, and retry failure.
  - [ ] Coalesce progress into one polite live region, announce immediate failures as alerts and terminal state once, and keep verbose events/heartbeats/reconnect attempts from flooding assistive technology.
  - [ ] On replay gap, operation expiry, terminal event, reconnect exhaustion, or page restoration, display a non-terminal transition when needed and refresh operation plus durable product state before claiming an outcome.
  - [ ] Add deterministic browser/JavaScript fixtures for duplicate suppression, pending controls, stale confirmation, explicit cancellation, event ordering/replay/gaps, bounded reconnect/polling, focus races, terminal refresh, disconnect isolation, expiry, errors, and no-JavaScript fallback.
  - [ ] **Dependency:** Tasks 3, 5, and 7.
  - [ ] **Evidence:** browser fixtures, CSS/render review matrix, timer/EventSource cleanup checks, accessibility semantic tests, and recorded manual keyboard/focus/announcement/color/motion/zoom/reflow checks.
  - [ ] **Stop condition:** do not use client state to repair server retry, terminal, or recovery correctness; move correctness failures back to the owning server/app task.

- [x] **Task 9: Prove Representative Workflows And Cross-Surface Agreement**
  > Executes: Decisions 1, 4, 5, and 7; CLI Surface, Workflows, Testing, LLM Runtime; `AC-06`-`AC-08`, `AC-17`-`AC-20`; `OUT-Representative browser integration tests`, `OUT-Interface capability tests`
  - [ ] Build `internal/web/integration_test.go` around real app use cases, temporary workspaces, `httptest`, deterministic clocks/IDs, and fake runtime/process/smoke-harness dependencies.
  - [ ] Exercise one substantial study journey and one substantial sprint journey through inspect, prepare, safe confirmation details, start, progress, cancellation or completion, durable terminal state, refresh, reconnect, operation expiry, and recovery.
  - [ ] Compare app, CLI, TUI, HTML, and JSON projections for workspace identity, readiness, findings, artifacts, verdicts, terminal outcomes, next actions, blocked states, and cleanup uncertainty using the same fixtures.
  - [ ] Prove product-owned mutation exclusion across CLI, TUI, and web and verify conflicts are typed/actionable rather than accidental concurrent execution.
  - [ ] Include browser refresh, navigation, tab/SSE disconnect, slow subscriber, forced interruption, and server restart cases proving they neither determine product success nor create detached/orphaned work.
  - [ ] Exercise safe runtime/model metadata and confirmation scope without exposing prompts, provider payloads, stderr, raw paths, arguments, sessions, or secrets.
  - [ ] **Dependency:** Tasks 2-8.
  - [ ] **Evidence:** deterministic representative study/sprint integration results, cross-surface comparison assertions, mutation-conflict proof, and forbidden-value scan.
  - [ ] **Stop condition:** any surface disagreement in durable state, terminal outcome, readiness, artifacts, verdict, or next action blocks documentation and release sign-off.

- [x] **Task 10: Embed, Package, And Reconcile Public Documentation**
  > Executes: Decisions 2, 3, 6, and 7; Documentation, CLI Surface, Security, Testing; `AC-02`, `AC-05`, `AC-21`-`AC-23`; all documentation outputs
  - [ ] Embed and resolve all templates, CSS, and JavaScript in the built binary; remove source-tree fallback behavior that could hide a missing embedded asset.
  - [ ] Build the binary and launch it from a temporary directory outside the source checkout with an empty browser/cache context; request every documented page, API capability, and static asset and verify no Node.js, Vite, package manager, CDN, database, source asset tree, or separate frontend process is required.
  - [ ] Update `docs/local-web.md` with startup, route/API/SSE contracts, cache behavior, guarded workflows, limits, security scope, accessibility behavior, shutdown, recovery, and troubleshooting.
  - [ ] Update `docs/configuration.md` with every web setting, source precedence, numeric-loopback restriction, defaults, valid ranges/combinations, fail-closed behavior, and redacted diagnostics.
  - [ ] Update `docs/cli-reference.md` and `docs/user-guide.md` with `ultraplan serve`, flags, lifecycle, supported browser inspection/operation journeys, and CLI/TUI/web agreement.
  - [ ] Update `docs/recovery.md` with reconnect, cancellation, interruption, stale locks, restart reconciliation, cleanup uncertainty, operation expiry, replay gaps, and durable recovery actions.
  - [ ] Update `docs/architecture.md` with web-to-app direction, capability vocabulary, presentation hierarchy, operation ownership, durable truth, shutdown ordering, and future-stage extensibility boundary.
  - [ ] Update `docs/release-checklist.md` with deterministic, race/leak, build, outside-tree packaging, API compatibility, security/redaction, manual accessibility, fake integration, architecture/sprint review, and gated real-system checks.
  - [ ] Reconcile every documented route, method, field, code, state, bound, flag, cache rule, recovery action, and security limitation against Task 3 fixtures and Task 4 effective configuration; do not document intended behavior that tests do not prove.
  - [ ] **Dependency:** Tasks 3-9.
  - [ ] **Evidence:** outside-tree packaging transcript, complete asset/route request matrix, and documentation-to-fixture/config reconciliation checklist.
  - [ ] **Stop condition:** package or documentation sign-off is blocked by any source-tree asset dependency or mismatch with compatibility fixtures/effective configuration.

- [/] **Task 11: Run Deterministic Release Gates And Resolve Applicable Findings** — Deferred: focused/full/race/vet/build/package gates pass; Architecture Review and Sprint Review are intentionally owned by the downstream `review` stage and were not generated during execute.
  > Executes: Decision 7; Testing, Security, Performance, Documentation, Observability; `AC-01`, `AC-03`-`AC-24`, `C-11`-`C-14`; all test and release outputs
  - [ ] Run focused `internal/app` and `internal/web` suites first so capability, API, security, operation/SSE, lifecycle, template, and integration failures remain localized.
  - [ ] Run the full deterministic suite, race suite, and production binary build from `../ultraplan-go` without mutating Git state or disturbing unrelated worktree changes.
  - [ ] Run outside-tree packaging, forbidden-value, goroutine/subscriber/lock leak, and documentation reconciliation checks and retain their reviewable evidence.
  - [ ] Execute the selected Architecture Review protocol against the dependency boundary, ownership, capability model, template hierarchy, ephemeral operation state, and lack of alternate persistence or route-specific workflow logic.
  - [ ] Execute the selected Sprint Review protocol against all selected contracts, reasoning decisions, plan tasks, implementation diff, deterministic evidence, documentation, and release checklist.
  - [ ] Resolve all applicable high-severity local-web security, path, redaction, concurrency, cancellation, shutdown, API compatibility, and accessibility findings; do not waive them through fixture updates or narrowed test scope.
  - [ ] Record every implementation deviation from `reasoning.md` before continuing and route architecture/scope changes back through governed planning rather than silently accepting them.
  - [ ] **Dependency:** Tasks 1-10 complete.
  - [ ] **Evidence:** focused/full/race/build/package results, Architecture Review evidence, current Sprint Review evidence, finding dispositions, and no unresolved applicable high-severity finding.
  - [ ] **Stop condition:** any deterministic gate failure or unresolved applicable high-severity finding blocks gated release evidence and sprint completion.

- [/] **Task 12: Capture Gated Real-System Evidence And Final Release State** — Deferred: the required current review verdict does not yet exist, so the review-gated real runtime and `ultraplan-go-smoke` harness operations must run in the downstream `smoke` stage; execute does not create `review.md` or `smoke.md`.
  > Executes: Decision 7; LLM Runtime, LLM Evaluation / Cost / Safety, Testing, Observability, Documentation; `AC-07`, `AC-15`, `AC-23`, `AC-24`; gated release evidence
  - [ ] After deterministic and review gates pass, use the selected Deep Smoke Sprint protocol and cataloged `ultraplan-go-smoke` harness to exercise one real runtime-backed browser operation and one real smoke-harness-backed browser operation.
  - [ ] Verify safe confirmation metadata, bounded progress, explicit cancellation or completion, durable terminal/recovery state, redacted evidence, and browser/CLI/TUI agreement for the gated operations.
  - [ ] Link harness run/evidence identifiers from the governed smoke output without copying raw harness evidence into the sprint or creating alternate issue artifacts.
  - [ ] If any prerequisite is unavailable, record `blocked` with the exact executable, provider/model, credential, network, browser, or harness prerequisite and retain deterministic release results; never report skipped/unavailable work as passed.
  - [ ] Re-run only the deterministic checks affected by any gated finding, then repeat required review/smoke gates according to their protocols; do not treat a narrow rerun as a substitute for the mandatory release matrix.
  - [ ] Complete the release checklist only when all deterministic checks, documentation, reviews, packaging, and gated pass-or-blocked records are current and mutually consistent.
  - [ ] **Dependency:** Task 11 passes with no unresolved applicable high-severity finding.
  - [ ] **Evidence:** real runtime and harness run links with pass evidence, or a truthful blocked record naming exact prerequisites; finalized release checklist.
  - [ ] **Stop condition:** a failed real operation is a failure to investigate, not a blocked prerequisite or pass; completion requires current evidence and truthful classification.

## Evidence Checklist

- [x] Import and delegation evidence proves `internal/web -> internal/app + stdlib` only and no CLI subprocess, product/runtime/process import, or route-owned workflow logic.
- [x] The fake future-stage fixture proves status, artifacts, commands, progress, cancellation, terminal results, and recovery without a route-specific branch or production plugin registry.
- [x] Reviewed compatibility fixtures freeze every documented `/api/v1` route, method, envelope, field, type/nullability, status, code, content type, and cache rule.
- [x] Semantic API tests prove unknown-route/method JSON behavior, typed safe errors, retry-safe start, idempotent cancellation, operation expiry, and durable refresh.
- [x] Configuration tests prove precedence, documented defaults, coherent bounds, numeric IPv4/IPv6 loopback, and fail-closed startup.
- [x] Security tests cover Host/Origin, CSRF/session, CSP/headers, request limits, malformed/duplicate-significant JSON, paths/references, hostile Markdown/content, and cache behavior.
- [x] No token, cookie, CSRF/session/confirmation value, prompt, unsafe path, provider payload, raw stderr, executable argument, or representative secret appears in HTML, JSON, SSE, retained events, terminal results, logs, or diagnostics.
- [x] Operation/SSE tests prove all configured bounds, confirmation binding/expiry/deduplication, event IDs/order/replay/gaps/rollover, heartbeat, slow-client isolation, disconnect behavior, terminal arbitration, and cleanup.
- [x] Lifecycle tests prove draining, mutation rejection, exact-once reasoned cancellation, no I/O/waits under hub locks, bounded cleanup, durable terminal/uncertainty, stale-lock handling, forced interruption, and conservative restart reconciliation.
- [x] Race and leak evidence shows no orphaned goroutine, timer, event source, subscriber, operation, process tree, or product lock.
- [x] Template tests prove one embedded namespaced tree, downward-only composition, explicit typed view models, startup failures for invalid definitions, escaping, no-JavaScript completeness, and embedded path resolution.
- [ ] Accessibility evidence covers deterministic semantics plus manual keyboard, visible focus, announcement timing, color independence, reduced motion, 200% zoom, text enlargement, and narrow reflow on representative pages.
- [x] Representative study and sprint browser journeys pass over temporary workspaces and fakes, including prepare, confirm, start, observe, cancel/complete, refresh/reconnect, restart, and recovery.
- [x] App, CLI, TUI, HTML, and JSON agree on durable state, readiness, findings, artifacts, verdicts, terminal outcomes, blocked/uncertain states, and next actions.
- [x] The built binary serves every page and asset outside the source tree with no frontend runtime, CDN, database, or separate process.
- [x] Local-web, configuration, CLI, user, recovery, architecture, troubleshooting, accessibility, packaging, and release documentation matches fixtures and effective configuration.
- [x] `go test ./...`, `go test -race ./...`, and `go build ./cmd/ultraplan` pass from `../ultraplan-go`.
- [ ] Architecture Review and Sprint Review evidence is current and contains no unresolved applicable high-severity local-web finding.
- [ ] Gated real-runtime and smoke-harness evidence passes, or unavailable prerequisites are recorded as blocked with exact reasons.
- [x] Deviations from `reasoning.md` are recorded before implementation continues.

## Verification Commands

Run commands from `../ultraplan-go` unless the row says otherwise. Commands are planned release checks and are not executed by this planning stage.

| Check | Command | Expected Result |
| --- | --- | --- |
| Shared app capability tests | `go test ./internal/app -run 'TestWeb|TestOperation|TestCapability'` | Shared typed capabilities and fake future-stage exposure pass without route-specific workflow orchestration. |
| Focused web tests | `go test ./internal/web` | API, security, operation/SSE, lifecycle, templates, and integration tests pass deterministically. |
| API compatibility | `go test ./internal/web -run 'TestAPICompatibility|TestRoute|TestMethod|TestErrorMapping|TestCache'` | Every frozen route/method/envelope/status/code/content-type/cache fixture and semantic mapping passes. |
| Security and redaction | `go test ./internal/web -run 'TestSecurity|TestHost|TestOrigin|TestCSRF|TestSession|TestCSP|TestBody|TestPath|TestMarkdown|TestRedaction'` | Fail-closed browser/path/content controls pass and forbidden values are absent from all captures. |
| Operations and SSE | `go test ./internal/web -run 'TestOperation|TestConfirmation|TestSSE|TestSubscriber|TestReplay|TestDisconnect'` | Bounds, deduplication, event protocol, gaps, slow clients, disconnect isolation, and cleanup pass. |
| Lifecycle and restart | `go test ./internal/web -run 'TestServer|TestShutdown|TestDraining|TestCleanup|TestRestart|TestReconcile'` | Exact-once cancellation, bounded cleanup, durable uncertainty, and conservative restart behavior pass. |
| Templates and accessibility semantics | `go test ./internal/web -run 'TestTemplate|TestRender|TestNoJavaScript|TestAccessibility|TestEmbedded'` | Namespaced downward composition, startup validation, complete pages, escaping, semantics, and embedded paths pass. |
| Representative integration | `go test ./internal/web -run 'TestIntegration|TestStudyWorkflow|TestSprintWorkflow|TestSurfaceAgreement'` | Study/sprint journeys and app/CLI/TUI/HTML/JSON agreement pass over temporary workspaces and fakes. |
| Web race pressure | `go test -race ./internal/web ./internal/app` | No operation, subscriber, cancellation, terminal, shutdown, or reconciliation race is reported. |
| Full deterministic suite | `go test ./...` | All package tests pass without live provider or external harness requirements. |
| Full race suite | `go test -race ./...` | Full repository passes race detection with no leak/race failure. |
| Production build | `go build ./cmd/ultraplan` | The single `ultraplan` binary builds with embedded presentation assets. |
| Serve help and flags | `go run ./cmd/ultraplan serve --help` | Documented loopback, lifecycle, browser, and limit flags are present and stable with safe help text. |
| Packaging check | `go test ./internal/web -run 'TestPackagedBinary|TestEmbeddedAssetsOutsideSourceTree'` | A built binary launched outside the checkout serves all pages/assets without source-tree or frontend-runtime dependencies. |
| Plan conformance review | `ultraplan sprint ultraplan-go 32-hardening-and-release review` | Current review covers selected contracts/protocols and reports no unresolved applicable high-severity finding. |
| Gated deep smoke | `ultraplan sprint ultraplan-go 32-hardening-and-release smoke` | One runtime-backed and one harness-backed browser operation pass, or exact missing prerequisites are reported as blocked rather than passed. |

## Risks And Blockers

| Risk / Blocker | Source | Mitigation | Status |
| --- | --- | --- | --- |
| The compatibility baseline is documented but not enumerated in planning inputs. | `reasoning.md#assumptions-and-risks`; `reasoning/api-design.md#risks` | Task 1 inventories implementation/docs first; unexplained differences block fixture creation and DTO changes. | open |
| Existing route-specific workflow or direct state interpretation may make the boundary refactor larger than estimated. | `reasoning.md#assumptions-and-risks`; `reasoning/architecture.md#risks` | Move each violation to its owning app/product module before hardening; resize tasks without changing the decided architecture. | open |
| Shared capability vocabulary could become a god interface or duplicate transport DTOs. | `reasoning.md#potential-technical-debt`; `reasoning/architecture.md#risks` | Use cohesive consumer-oriented capabilities and concrete app results; prove the seam with the fake-stage fixture and reject broad service containers. | open |
| Confirmation consumption, deduplication, capacity reservation, publication, and shutdown can race. | `reasoning.md#assumptions-and-risks` | Define one critical ordering and prove it with deterministic barriers plus race tests. | open |
| Completion, failure, timeout, user cancellation, and shutdown can produce conflicting terminal state. | `reasoning.md#assumptions-and-risks` | Use one terminal arbiter, exact-once canonical cancellation, and authoritative-result preservation tests. | open |
| Hub, product locks, persistence, subscribers, and cleanup can deadlock. | `reasoning.md#assumptions-and-risks`; `reasoning/architecture.md#risks` | Document lock order; snapshot under lock; perform no callbacks, I/O, sends, cancellations, or waits under hub locks. | open |
| Cleanup deadline exhaustion can still produce false success or detached ownership. | `reasoning.md#decision-5-make-shutdown-and-restart-reconciliation-truthful-and-product-owned` | Persist owner-specific uncertainty before closing/releasing ownership; make failure to do so a release blocker. | open |
| Replay gaps or expired operation handles may appear terminal to users. | `reasoning.md#assumptions-and-risks`; `reasoning/frontend.md#risks` | Use explicit non-terminal gap/expiry states and mandatory durable HTTP refresh before outcome claims. | open |
| Resource defaults may reject substantial work or weaken memory/shutdown guarantees. | `reasoning.md#assumptions-and-risks` | Preserve safe existing values, test representative workflows, document/tune explicitly, and retain hard caps. | open |
| Static asset URLs may not be content-addressed. | `reasoning.md#assumptions-and-risks`; `reasoning/api-design.md#risks` | Use revalidation unless Task 1 proves content addressing; test cache policy in compatibility and packaging checks. | open |
| Template migration may alter routes, forms, fields, or no-JavaScript behavior. | `reasoning/frontend.md#risks` | Add behavior/semantic fixtures before extraction and require rationale for every fixture difference. | open |
| Static accessibility tests cannot prove focus, announcements, color independence, zoom, reflow, or motion behavior. | `reasoning.md#assumptions-and-risks`; `reasoning/frontend.md#risks` | Record mandatory manual representative-page checks in the release checklist. | open |
| Imperative enhancement may leak timers, streams, focus updates, or duplicate requests. | `reasoning.md#potential-technical-debt`; `reasoning/frontend.md#risks` | Assign one abort owner, enforce retry/lifetime limits, and test navigation/terminal cleanup and focus races. | open |
| Rich diagnostics or browser projections may leak sensitive runtime/browser values. | `reasoning.md#assumptions-and-risks` | Allowlist fields, redact at type and projection boundaries, avoid raw bodies/payloads/stderr, and scan a forbidden corpus everywhere. | open |
| Source-tree paths or browser cache can hide missing embedded assets. | `reasoning/architecture.md#risks`; `reasoning/frontend.md#risks` | Build and launch outside the checkout with empty cache and request every referenced page/asset. | open |
| Constraint labels in `reasoning.md` do not map one-to-one to the visible requirements bullets. | `reasoning.md` and `requirements.md` | Trace implementation/review to quoted behavior; record any normalized labels instead of guessing or dropping a contract. | open |
| Real runtime or smoke-harness prerequisites may be unavailable. | `reasoning.md#assumptions-and-risks` | Record exact missing prerequisites as blocked; never convert unavailable/skipped evidence into pass. | open |
| Unrelated working-tree changes may exist during release verification. | `reasoning.md#assumptions-and-risks` | Run non-destructive checks only; perform no automatic Git mutation or cleanup. | open |

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
- `system/protocols/architecture-review-protocol.md`
- `system/protocols/review-sprint-protocol.md`
- `system/protocols/deep-smoke-sprint-protocol.md`

## Execution Log

| Date / Step | Action | Evidence / Notes |
| --- | --- | --- |
| `2026-08-18 / baseline` | Inventoried routes, methods, envelopes, cache behavior, package imports, operation ownership, assets, and limits before compatibility changes. | `docs/web-compatibility-baseline.md`; import/API compatibility tests. |
| `2026-08-18 / implementation` | Hardened the web/app boundary, origin and framing policy, immutable server bounds, retry-safe operation starts, namespaced templates, layered embedded assets, progressive enhancement, integration/package tests, and public docs. | Target implementation diff and focused `internal/app`/`internal/web` suites. |
| `2026-08-18 / release gates` | Ran all focused matrix rows, full deterministic tests, full race tests, vet, production build, serve help, package launch outside checkout, and `git diff --check`; all passed. | Commands in Verification Commands; build output directed to `/tmp/ultraplan-sprint32`. |
| `2026-08-18 / deviations and deferrals` | Kept resource caps fixed and documented instead of inventing unsupported configuration; preserved adjacent template files while validating namespaced downward composition and adding the required layer manifests. Manual browser accessibility is blocked by no installed Chromium/Firefox. Review and smoke remain downstream-owned. | `docs/web-compatibility-baseline.md`; Tasks 8, 11, and 12 deferral rationales. |

## Completion Criteria

- [x] All tasks are complete or explicitly deferred with requirement, risk, and release impact recorded.
- [x] Verification commands were run or environmental deferrals are documented as blocked with exact prerequisites.
- [x] Evidence satisfies every execute-stage expectation in `reasoning.md`, including the normal-versus-gated distinction.
- [x] The Sprint 31 shutdown gap is covered by durable owner-specific cleanup uncertainty on deadline exhaustion and conservative restart reconciliation tests.
- [x] Compatibility, security, concurrency, cancellation, shutdown, recovery, presentation, deterministic accessibility, packaging, and documentation evidence is current and mutually consistent.
- [ ] Architecture Review and Sprint Review can evaluate conformance without guessing intent and report no unresolved applicable high-severity local-web finding.
- [ ] Gated real-system evidence passes or is truthfully blocked; a failed operation is not relabeled as unavailable.
- [ ] `review.md` can evaluate conformance without guessing intent.
