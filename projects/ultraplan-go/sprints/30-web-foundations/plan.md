# Sprint Plan: Local Web Foundation and Read-Only Dashboard

> Project: `ultraplan-go`
> Sprint: `30-web-foundations`
> Source: `reasoning.md`
> **Inputs Used:** `projects/ultraplan-go/project-index.md`, `projects/ultraplan-go/sprints/30-web-foundations/requirements.md`, `projects/ultraplan-go/docs/ARCHITECTURE.md`, `projects/ultraplan-go/docs/PRD.md`, `projects/ultraplan-go/docs/TRD.md`, `projects/ultraplan-go/sprints/30-web-foundations/sprint-index.md`, `projects/ultraplan-go/sprints/30-web-foundations/technical-handbook.md`, `projects/ultraplan-go/sprints/30-web-foundations/reasoning/api-design.md`, `projects/ultraplan-go/sprints/30-web-foundations/reasoning/architecture.md`, `projects/ultraplan-go/sprints/30-web-foundations/reasoning/frontend.md`, `projects/ultraplan-go/sprints/30-web-foundations/reasoning.md`, `builtin:templates/sprint-plan.md`

This plan executes `reasoning.md`. It does not invent architecture, scope, or decisions. The selected reports were not reopened because `technical-handbook.md` and the completed reasoning documents contain sufficient evidence and concrete source references for every planned decision; no selected evidence was otherwise omitted.

## Reasoning Source

- **Sprint Reasoning:** `reasoning.md`
- **Sprint Index:** `sprint-index.md`
- **Technical Handbook:** `technical-handbook.md`
- **Area Reasoning:** `reasoning/api-design.md`, `reasoning/architecture.md`, `reasoning/frontend.md`

## Sprint Status

- **Status:** `implementation complete; ready for review`
- **Owner:** `Codex execution agent`
- **Start Date:** `2026-07-29`
- **Completion Date:** `2026-07-29`

Implementation must not start until the owner confirms that Product Phase 3 is complete or that its exceptions are recorded. A failed readiness check is a blocker, not permission to duplicate missing Phase 3 behavior in `internal/web`.

## Decisions To Execute

| Decision | Source Section | Execution Implication |
| --- | --- | --- |
| Explicit Surface Composition And Query-Only App Boundary | `reasoning.md#decision-1-explicit-surface-composition-and-query-only-app-boundary` | Keep `cmd/ultraplan` construction-only; inject independent TUI/web runners; expose a cohesive query-only app facade; permit `internal/web -> internal/app` only. |
| Serve Command, Immutable Configuration, And Owned Lifecycle | `reasoning.md#decision-2-serve-command-immutable-configuration-and-owned-lifecycle` | Add lazy `serve` dispatch, source-aware immutable config, numeric loopback validation, owned listener/goroutines, fixed timeouts/concurrency, graceful shutdown, and warning-only browser launch. |
| Read-Only HTML And Versioned JSON Route Contract | `reasoning.md#decision-3-read-only-html-and-versioned-json-route-contract` | Register only the decided `GET`/implicit `HEAD` HTML, static, and `/api/v1` resources; map explicit DTO/envelope/error contracts; never fall API routes through to HTML. |
| Bounded Opaque Artifact Preview | `reasoning.md#decision-4-bounded-opaque-artifact-preview` | Issue and resolve opaque app-owned references; enforce containment, allowlisting, and a 256 KiB bound; render Markdown/JSON as escaped source without direct web filesystem access. |
| Local HTTP Security, Error Projection, And Diagnostics | `reasoning.md#decision-5-local-http-security-error-projection-and-diagnostics` | Enforce exact Host/Origin and no-CORS policy, request limits, headers, typed safe errors, server request IDs, and redacted structured diagnostics. |
| Server-Rendered Accessible Progressive Frontend | `reasoning.md#decision-6-server-rendered-accessible-progressive-frontend` | Embed complete `html/template` pages, one stylesheet, and at most one small dependency-free script; preserve no-JavaScript, accessibility, responsive reflow, and hostile-content safety. |
| Live, Bounded, Sequential State Projection And Truthful Health | `reasoning.md#decision-7-live-bounded-sequential-state-projection-and-truthful-health` | Read authoritative state per request, aggregate sequentially, cap collections at 200 with truncation metadata, avoid caches/watchers/polling/fan-out, and keep health cheap and truthful. |
| Layered Verification, Documentation, And Existing-Surface Protection | `reasoning.md#decision-8-layered-verification-documentation-and-existing-surface-protection` | Combine focused fakes, `httptest`, listener/command paths, race/build checks, manual accessibility evidence, three documentation changes, and required review protocols while protecting CLI/TUI behavior. |

## Decision Execution Records

| Decision | Requirements Satisfied | Evidence Used | Trade-Off Accepted | Alternative Rejected | Risk / Follow-Up |
| --- | --- | --- | --- | --- | --- |
| 1. Explicit composition and query-only boundary | `RO-4`-`RO-7`, `AC-7`, `AC-12`, `AC-13`, `C-1`-`C-5` | Architecture contract and docs; handbook reports `01`, `02`, `03`, `06`, `12`, `13`; architecture reasoning | Visible runner/facade wiring and result mapping | App-constructed web cycle, mutable global registration, web-to-product imports, broad service injection, interface-per-route fragmentation | Prevent facade and composition-root growth; Sprint 31 requires a separate command capability. |
| 2. Serve command and lifecycle | `RO-4`, `RO-6`, `RO-8`, `RO-9`, `AC-1`-`AC-5`, `AC-13`, `AC-14`, `C-3`, `C-4`, `C-6`, `C-10` | Configuration, Errors, Observability, Security, CLI Surface, Performance, Testing; handbook reports `03`, `04`, `07`, `08`, `10`, `11`, `14` | Restart for config changes, fixed limits, no `localhost`, warning-only launcher failure | Non-loopback/hostname bind, hot reload, eager web setup, detached lifecycle, fatal launcher failure | Measure limits; preserve validation-before-listen and owned cleanup. |
| 3. HTML and `/api/v1` contract | `RO-10`-`RO-12`, `AC-6`, `AC-7`, `AC-9`, `AC-12`, `C-1`, `C-2`, `C-5`, `C-9` | API reasoning; Architecture, Errors, Observability, Security, Testing, Documentation; handbook reports `01`, `02`, `03`, `05`, `11`, `12` | DTO mapping and internal `v1` compatibility burden; whole-query failures | RPC query endpoint, direct model serialization, API-to-HTML fallback, operations/confirmation/SSE routes | Breaking DTO changes require a version or coordinated migration; do not claim public API support. |
| 4. Artifact preview | `RO-13`, `RO-14`, `RO-18`, `AC-6`, `AC-10`-`AC-12`, `C-1`, `C-2`, `C-5`, `C-7`, `C-10` | API/architecture/frontend reasoning; Security, Errors, Performance, Testing; handbook reports `06`, `13`, `14` | Source-only preview, opaque stale refs, fixed truncation | Caller paths, lexical cleaning alone, web filesystem reads, whole-file reads, rich Markdown | Keep symlink/encoding matrix strong; larger or rich previews need new evidence and security reasoning. |
| 5. HTTP security and diagnostics | `RO-12`, `RO-14`, `RO-17`, `AC-3`, `AC-9`-`AC-11`, `C-5`-`C-10` | Security, Errors, Observability, Performance, Testing; handbook reports `05`, `10`, `13` | Strict aliases/proxies, reduced response/log detail, no session/auth state | Trusting loopback, permissive CORS, caller request IDs, raw/string-parsed errors, hosted auth, differentiated path errors | This is local single-user exposure reduction, not isolation from a compromised account; Sprint 31 needs separate CSRF/session reasoning. |
| 6. Server-rendered frontend | `RO-15`, `RO-19`, `AC-6`, `AC-8`, `AC-10`-`AC-12`, `C-2`, `C-5`, `C-7`, `C-8`, `C-10` | Frontend reasoning; Architecture, Security, Performance, Testing, Documentation; handbook reports `01`, `02`, `03`, `06`, `08`, `10`, `11`, `13`, `14` | Less interactivity and buffered rendering | React/Vue/Vite, client routing/state, polling, custom widgets, third-party/inline assets, trusted workspace markup | Manual browser/accessibility review remains required; automation may be added during hardening. |
| 7. Live bounded state and health | `AC-4`, `AC-6`, `AC-14`, `C-2`, `C-4`, `C-10` | Architecture, Observability, Performance, Testing; handbook reports `07`, `08`, `14` | Repeated reads, sequential latency, fail-whole aggregates, visible truncation | Cache/snapshot/watcher, polling, unbounded results, speculative fan-out/partial panels, deep health | Add cache, pagination, fan-out, or partial results only after measurement and an explicit contract. |
| 8. Verification and docs | `RO-8`, `RO-16`-`RO-23`, `AC-1`, `AC-13`, `AC-14`, `C-10` | Testing, Documentation, CLI Surface, Architecture, Performance; handbook reports `06`, `11`; selected review protocols | Broad deterministic suite and selective golden maintenance | Handler-only tests, all-page goldens, real runtime/smoke dependencies, manual-only checks, changed CLI/TUI output | Record browser/accessibility evidence and protocol results; do not let test seams proliferate beyond behavior boundaries. |

## Requirements / Contracts To Satisfy

Requirement identifiers follow the ordering defined in `requirements.md` and explained in `reasoning.md`.

| Contract / Requirement ID | Required Behavior | Evidence Planned |
| --- | --- | --- |
| Architecture; `RO-4`-`RO-7`; `AC-7`, `AC-12`, `AC-13`; `C-1`-`C-5` | Explicit `cmd -> app+tui+web`, `web -> app`; query-only app facade; no global runner, product import, CLI parsing, mutation capability, or duplicate state. | Import/dependency review, composition/use-case tests, full test/build, Architecture Review. |
| Configuration; CLI Surface; `AC-1`-`AC-4`; `C-3`, `C-6` | `serve --help`, existing workspace discovery, explicit `--listen` and `--open-browser`, numeric loopback only, immutable source-aware config, validation before runner/listener. | Help, precedence, workspace, bind matrix, runner-not-called, startup/exit tests. |
| Errors; `AC-1`, `AC-9`, `AC-11`; `C-2`, `C-9` | Preserve typed identity and separately map CLI, HTML, and the fixed API status/code set without raw causes or string parsing. | Error table tests, response/log leak assertions, startup and route tests. |
| Observability; `AC-4`, `AC-6`, `AC-11`; `C-2`, `C-4` | Truthful health, server-generated request IDs, safe stable request/lifecycle diagnostics, and separate diagnostic/output channels. | Health `200`/`503`, request ID, structured capture, redaction, shutdown tests. |
| Security; `AC-3`, `AC-10`, `AC-11`; `C-5`-`C-9` | Loopback, exact Host/Origin, absent-Origin rules, no CORS, bodies/targets/identifiers/query limits, security headers, opaque contained previews, safe rendering. | Host/Origin and limit matrices, path/symlink/hostile-content tests, header assertions. |
| Performance; `RO-9`, `RO-13`; `AC-4`, `AC-10`, `AC-14`; `C-10` | Lazy startup; fixed HTTP/shutdown/concurrency limits; 200-entry and 256 KiB bounds; cancellation; sequential reads; no recursive scans/cache/polling/fan-out. | Boundary/truncation/cancellation tests, diagnostics, race suite, implementation review. |
| Testing; `RO-8`, `RO-16`-`RO-20`; `AC-14`; `C-10` | Deterministic app/web fakes, `httptest`, actual command/listener paths, semantic assertions, narrow envelope goldens, race and build checks. | Focused suites and all commands in Verification Commands. |
| Documentation; `RO-21`-`RO-23`; `AC-1`, `AC-8` | Document command, local/read-only trust boundary, state authority, bounds, previews, shutdown, troubleshooting, internal API status, and package direction. | Help assertions and review of all three implementation-repository docs. |
| `AC-5` | Context cancellation or signal performs bounded graceful shutdown without blocked handlers or leaked goroutines. | Controlled listener/in-flight handler cancellation tests and race suite. |
| `AC-6`, `AC-7` | Pages inspect workspace, projects, sprints, studies, validations, flow/run/review/smoke status, and artifacts exclusively through typed app query results. | Fake app result tests and temporary-workspace agreement tests. |
| `AC-8`; `C-7`, `C-8` | Single-binary `html/template`, embedded CSS and minimal JS; no Node, Vite, separate asset server, database, frontend framework, or build step. | Embed/template tests, dependency/build inspection, successful Go build. |
| `AC-9`; `C-9` | Stable `{data,meta}` success and `{error,meta}` failure objects; unknown `/api/` and unsupported versions always return JSON. | Route/envelope/content-type/cache tests and narrow normalized goldens. |
| `AC-10` | Only app-issued Markdown/JSON refs resolve; previews are contained, allowlisted, 256 KiB bounded, and inert. | App artifact and web rendering/security tests. |
| `AC-13` | Existing CLI and TUI behavior, output, dispatch, startup cost, and exit classes remain unchanged except explicit composition. | Focused dispatch/TUI tests and full regression suite. |
| `NG-1`-`NG-10`; `C-5`, `C-11` | No operations, mutation, runtime/smoke invocation, confirmations, cancellation API, SSE/WebSockets, remote/auth/multi-user service, arbitrary editing, alternate persistence, frontend build stack, issue tracking, fixes, or Git mutation. | Route/UI/capability/import absence assertions and scope review. |

## Fixed Implementation Contract

- `serve` defaults to `--listen 127.0.0.1:8080`; `--open-browser` defaults to false; the existing global `--workspace` behavior remains authoritative.
- Listen values must be numeric loopback IP literals with explicit valid ports. `127.0.0.1` and bracketed `::1` forms are valid; `localhost`, malformed values, and non-loopback addresses fail before runner/listener use.
- HTTP limits are 5 seconds header-read, 15 seconds request-read, 30 seconds write, 60 seconds idle, 10 seconds graceful shutdown, and 32 in-flight requests.
- Request limits are 64 KiB declared body, 8 KiB request target, and 128 bytes per decoded identifier/reference. Sprint 30 routes reject bodies, duplicate/unknown query parameters, malformed identifiers, controls, and separators.
- Collections and findings return at most 200 entries with `returned_count`, `total_count`, and `truncated`. Artifact reads return at most 256 KiB plus one byte for truncation detection and report total/returned bytes.
- JSON success is `{data,meta}` with `api_version`, server-generated `request_id`, and `generated_at`. JSON failure is `{error:{code,message,details?},meta}`. Public codes are only `invalid_request`, `request_rejected`, `not_found`, `method_not_allowed`, `request_too_large`, `internal_error`, and `unavailable` with the status mapping in Decision 5.
- Exact Origin, when present, must match the effective loopback origin. Absent Origin is accepted only after Host validation for navigation and local `GET`/`HEAD`; `Origin: null`, malformed, cross-origin, and non-loopback origins are rejected. No permissive CORS response is emitted.
- All responses receive the applicable restrictive CSP, `X-Content-Type-Options: nosniff`, `Referrer-Policy: no-referrer`, frame denial, and `Cache-Control: no-store` policy.
- Health is a cheap server-readiness and configured-workspace-query check only: `200`/`ok` or `503`/`unavailable`. It must not scan collections/artifacts or invoke runtime, provider, review, or smoke checks.

## Route Inventory

| Methods | Route | Required Result |
| --- | --- | --- |
| `GET`, `HEAD` | `/` | Complete server-rendered workspace dashboard. |
| `GET`, `HEAD` | `/projects` | Bounded project list. |
| `GET`, `HEAD` | `/projects/{project}` | Project detail and bounded sprint summaries. |
| `GET`, `HEAD` | `/projects/{project}/sprints/{sprint}` | Sprint planning/flow/validation/review/smoke inspection. |
| `GET`, `HEAD` | `/studies` | Bounded study list. |
| `GET`, `HEAD` | `/studies/{study}` | Study status/run-state/validation inspection. |
| `GET`, `HEAD` | `/artifacts/{ref}` | Bounded escaped Markdown/JSON source preview. |
| `GET`, `HEAD` | `/api/v1/dashboard` | Dashboard DTO and bounded summaries. |
| `GET`, `HEAD` | `/api/v1/projects` | Bounded project DTO collection. |
| `GET`, `HEAD` | `/api/v1/projects/{project}` | Project DTO and bounded sprint summaries. |
| `GET`, `HEAD` | `/api/v1/projects/{project}/sprints/{sprint}` | Sprint status DTO. |
| `GET`, `HEAD` | `/api/v1/studies` | Bounded study DTO collection. |
| `GET`, `HEAD` | `/api/v1/studies/{study}` | Study status DTO. |
| `GET`, `HEAD` | `/api/v1/validations?scope=...&ref=...` | Existing findings only; never starts validation. |
| `GET`, `HEAD` | `/api/v1/artifacts/{ref}` | Bounded preview DTO. |
| `GET`, `HEAD` | `/api/v1/health` | Cheap readiness DTO. |
| `GET`, `HEAD` | `/static/...` | Embedded first-party assets only. |

Every other method for a known resource returns `405` with `Allow: GET, HEAD`. Unknown `/api/` paths and unsupported API versions return JSON errors; unknown browser paths render the safe HTML error page. No operation, confirmation, mutation, cancellation, SSE, or WebSocket route may be registered.

## Tasks

- [x] **Task 1: Confirm Phase 3 And Existing-Surface Readiness**
  > Executes: `Decision 8`; Dependencies; `AC-13`; Phase 3 readiness assumption
  - [x] Record whether the Phase 3 release gate is complete or identify its recorded exceptions before changing the implementation repository.
  - [x] Inventory the existing app query/error, workspace discovery/config precedence, CLI dispatch, and TUI runner seams needed by Sprint 30; record missing read projections as app work, not web-owned behavior.
  - [x] Identify representative existing CLI/TUI tests and outputs that must remain unchanged, including help/version paths that must not initialize web resources.
  - [x] Stop and record a blocker if required status/query behavior cannot be exposed without changing governed workflow semantics.

- [x] **Task 2: Establish Explicit Surface Composition And Query-Only App Use Cases**
  > Executes: `Decision 1`, `Decision 7`; `RO-4`, `RO-5`; `AC-6`, `AC-7`, `AC-12`, `AC-13`; `C-1`-`C-5`
  - [x] Refactor `cmd/ultraplan/main.go` and `internal/app/surfaces.go` so `cmd` explicitly constructs independent TUI and lazy web runner dependencies without globals, registration callbacks, service locators, or context-carried services.
  - [x] Add cohesive typed read-only app requests/results in `internal/app/web_usecases.go` for dashboard, project, sprint, study, validation summary, artifact reference/preview, and health projections.
  - [x] Keep aggregation, canonical identifiers, opaque artifact-reference issuance/resolution, containment, allowlisting, typed error identity, deterministic ordering, 200-entry limits, and truncation counts in app/product ownership.
  - [x] Propagate request context through every query and keep aggregate reads sequential with coherent whole-query failure.
  - [x] Add app tests proving fresh repeated reads do not mutate state, bounds at zero/exact/over-limit, cancellation, typed errors, opaque-reference behavior, cheap health, and agreement with existing surface semantics.
  - [x] Add composition/non-regression tests proving web facilities are lazy, runners are independently replaceable, and existing CLI/TUI paths retain behavior and output.

- [x] **Task 3: Add Serve Parsing, Preflight, And Lifecycle**
  > Executes: `Decision 2`, `Decision 5`; `RO-6`, `RO-8`, `RO-9`; `AC-1`-`AC-5`, `AC-13`; `C-3`, `C-4`, `C-6`, `C-10`
  - [x] Implement `internal/app/serve_commands.go` with documented `--listen`, `--open-browser`, existing global workspace behavior, existing output discipline, typed startup/cancellation errors, and existing exit classes.
  - [x] Preserve built-in, workspace, environment, explicitly set flag precedence; produce one immutable effective serve config after workspace/config validation.
  - [x] Validate numeric IPv4/IPv6 loopback and explicit ports before runner invocation, then recheck bind policy before listener acquisition; never silently choose a different port after collision.
  - [x] Implement `internal/web/server.go` with the fixed HTTP timeouts, 32-request bound, root/request cancellation, owned listener and goroutines, and 10-second graceful shutdown.
  - [x] Inject the listener and optional browser launcher seams where deterministic lifecycle tests require them; launch only after successful listen and report launch failure as a redacted warning without failing a healthy server.
  - [x] Test help, discovery/precedence, valid/rejected listen matrices, runner-not-called preflight, listen/start failure, canonical URL, launcher warning, cancellation, in-flight completion/cancellation, cleanup, and existing command/TUI dispatch.

- [x] **Task 4: Implement Read-Only Routing, DTOs, And Error Projection**
  > Executes: `Decision 3`, `Decision 5`; `RO-10`-`RO-12`; `AC-6`, `AC-7`, `AC-9`; `C-1`, `C-2`, `C-5`, `C-9`
  - [x] Register the exact Route Inventory in `internal/web/routes.go`, including embedded static assets, implicit `HEAD`, `405`/`Allow`, JSON API not-found/version handling, and HTML browser not-found handling.
  - [x] Implement thin handlers in `internal/web/handlers.go` that validate transport inputs, call only typed app queries, and independently map app results to explicit HTML view models and route-specific JSON DTOs.
  - [x] Implement success/error metadata, optional-field omission, empty-array behavior, stable field names, content types, `no-store`, generated timestamps/request IDs, collection metadata, and the fixed typed error mapping.
  - [x] Reject bodies and invalid request targets, identifiers, references, and validation queries at the decided boundaries; do not add pagination or undocumented parameters.
  - [x] Add table-driven tests for every route, method, `HEAD`, `Allow`, envelope, error code, content type, optional/null behavior, counts/truncation, unsupported version, unknown `/api/`, and absent operation routes.
  - [x] Use normalized goldens only for stable JSON success/error envelope shapes; use semantic assertions for IDs, times, paths, counts, and page content.

- [x] **Task 5: Enforce Local HTTP Security And Redacted Diagnostics**
  > Executes: `Decision 5`; `RO-14`, `RO-17`; `AC-3`, `AC-9`-`AC-11`; `C-5`-`C-10`
  - [x] Implement centralized exact Host authority and Origin policy for IPv4 and bracketed IPv6 listeners, including allowed absent Origin and rejected null/malformed/cross-origin/non-loopback cases with no CORS relaxation.
  - [x] Apply request concurrency/body/target/identifier/query limits and the decided security headers to HTML, JSON, static, and error responses.
  - [x] Generate request IDs server-side and emit normalized route/method/status/duration/response-byte/error-code plus lifecycle/shutdown diagnostics on the diagnostic channel.
  - [x] Ensure normal diagnostics and all responses omit raw URLs/query values, absolute paths, Host/Origin values, artifact content, environment, secrets, provider payloads, raw stderr, and internal causes.
  - [x] Add Host/Origin/no-CORS, body/target/query/identifier boundary, header, request-ID, all-status mapping, redaction/leak, cancellation, and race-oriented tests.

- [x] **Task 6: Implement Opaque Bounded Artifact Preview**
  > Executes: `Decision 4`, `Decision 5`; `RO-13`, `RO-14`, `RO-18`; `AC-10`-`AC-12`; `C-1`, `C-2`, `C-5`, `C-7`, `C-10`
  - [x] Complete app-owned opaque URL-safe reference issuance and canonical resolution against authoritative workspace state, including stale/forged, traversal, encoded traversal, absolute, separator/control, symlink escape, and non-allowlisted rejection.
  - [x] Allow only governed Markdown and JSON artifact classes and read at most 256 KiB plus one byte for truncation detection; return safe relative display metadata and total/returned/truncated fields.
  - [x] Keep `internal/web/artifacts.go` limited to validating the app result's media/size contract and mapping source content; do not inject or use an arbitrary filesystem reader in web.
  - [x] Render Markdown and JSON as escaped source in labelled `<pre><code>`; optionally include parsed JSON in API output only when it parses within the same bound; never use workspace data as `template.HTML`, `innerHTML`, scripts, or executable links.
  - [x] Collapse all invalid/stale/unsupported/escaping references to safe `404 not_found` and avoid existence disclosure.
  - [x] Test round trips, stale/forged refs, traversal/encoding/symlink/extension cases, exact/over-limit bytes, invalid JSON, hostile HTML/script/Markdown/URLs, escaping, `nosniff`, metadata, and direct-web-filesystem absence.

- [x] **Task 7: Build The Embedded Accessible Progressive Frontend**
  > Executes: `Decision 6`; `RO-15`, `RO-19`; `AC-6`, `AC-8`, `AC-10`-`AC-12`; `C-2`, `C-5`, `C-7`, `C-8`, `C-10`
  - [x] Add embedded templates for the shared shell, dashboard, project list/detail, sprint detail, study list/detail, validation presentation, artifact preview, health, empty, not-found, and error states.
  - [x] Parse templates once during serve startup and fail startup on parse error; buffer page rendering before headers so failures produce complete safe responses.
  - [x] Implement the semantic shell with skip link, unique title/`h1`, header, primary navigation, workspace context, breadcrumbs, `main`, footer/server status, and only earned repeated partials.
  - [x] Add one embedded stylesheet for high-contrast documentation-console presentation, visible focus, textual statuses, single-column narrow reflow, local table/code overflow, touch/zoom support, and reduced motion.
  - [x] If enhancement is needed, keep it to at most one small external embedded dependency-free script for explicit refresh with polite status, stale/failure messaging, focus preservation, and ordinary reload fallback; add no polling, router, store, cache, storage, cookies, service worker, operations, or SSE.
  - [x] Add template/route tests for representative normal, unknown, empty, error, and truncation states; hostile names/content/URLs; title/headings/landmarks/skip link/navigation/labels/tables/status/live regions; and absence of inline/third-party assets or operational controls.
  - [x] Record desktop, narrow viewport, keyboard, visible-focus, no-JavaScript, local-overflow, 200-percent zoom, and reduced-motion review evidence.

- [x] **Task 8: Complete Live Projection And Truthful Health Coverage**
  > Executes: `Decision 7`; `AC-4`, `AC-6`, `AC-14`; `C-2`, `C-4`, `C-10`
  - [x] Verify every page/API request reads current app-owned state, propagates cancellation, returns visible bounds, and never installs cache, watcher, snapshot, hidden preload, recursive source scan, browser persistence, background refresh, or handler fan-out.
  - [x] Implement `/api/v1/health` as server readiness plus lightweight configured-workspace query availability only, returning exactly `200`/`ok` or `503`/`unavailable` without deep checks or disclosure.
  - [x] Add repeated-read freshness/non-mutation, zero/exact/over-limit collection, fail-whole aggregate, request-cancellation, health/no-deep-scan, and duration-diagnostic tests.

- [x] **Task 9: Update User And Architecture Documentation**
  > Executes: `Decision 8`; `RO-21`-`RO-23`; `AC-1`, `AC-8`
  - [x] Add `../ultraplan-go/docs/local-web.md` covering startup inside/with an explicit workspace, exact flags/defaults, canonical loopback forms, read-only scope, local trust limits, explicit refresh, state authority, list/preview bounds, source-preview behavior, launcher warnings, shutdown, and troubleshooting.
  - [x] Update `../ultraplan-go/docs/cli-reference.md` with `serve` help, options, workspace/config behavior, loopback rejection, startup/cancellation/exit behavior, and CLI/TUI relationship.
  - [x] Update `../ultraplan-go/docs/architecture.md` with explicit composition, `internal/web -> internal/app`, query-only capability, transport/product state ownership, artifact resolution, embedded single-binary assets, and Sprint 31 deferrals.
  - [x] Document `/api/v1` as compatibility-controlled for the bundled browser, not a promised public integration API, and document health as readiness rather than complete product validity.
  - [x] Confirm docs do not imply `localhost`, remote access, rich Markdown, automatic polling, operations, SSE, auth, database state, or a frontend build pipeline.

- [x] **Task 10: Run Layered Verification And Prepare Review Evidence**
  > Executes: `Decision 8`; all acceptance criteria; `RO-16`-`RO-23`; `C-10`
  - [x] Run the focused app composition/serve/query suites and focused web lifecycle/route/security/artifact/template suites; retain command and result evidence.
  - [x] Run the full deterministic test, race, and binary-build commands from `../ultraplan-go`; normal verification must not require a browser process, provider/runtime, smoke harness, Node.js, database, frontend server, or network service.
  - [x] Inspect imports/dependencies, route inventory, assets/build files, and implementation diff for forbidden capabilities, scope, state, dependencies, unbounded work, and CLI/TUI regressions.
  - [x] Assemble the manual accessibility/responsive evidence and documentation/help evidence with the test/build results.
  - [x] Supply the implementation diff and evidence to `system/protocols/architecture-review-protocol.md` and `system/protocols/review-sprint-protocol.md`; record deviations from `reasoning.md` before implementation continues rather than silently changing this plan.

## Evidence Checklist

- [x] Tests prove all required behavior, including exact boundaries and rejected-scope absence.
- [x] Runtime or diagnostic evidence exists for startup, request, health, launcher warning, cancellation, and shutdown behavior.
- [x] Documentation and `serve --help` updates are complete and mutually consistent.
- [x] Deviations from `reasoning.md` are recorded before implementation continues.
- [x] Architecture Review evidence covers composition, imports, ownership, state, capability restriction, and existing-surface protection.
- [x] Sprint Review inputs cover all selected contracts, handbook guidance, decisions, plan tasks, verification evidence, findings, and final verdict.
- [x] Accessibility/responsive evidence covers no-JavaScript, keyboard/focus, semantics, narrow reflow, local overflow, zoom, reduced motion, and hostile-content inertness.
- [x] API evidence covers every route/method, exact envelopes/error codes, unknown API behavior, stable DTO fields, bounds, `HEAD`, content types, headers, and absent operation routes.
- [x] Security evidence covers listen, Host/Origin/no-CORS, input bounds, containment/symlinks, redaction, request IDs, safe errors, and response headers.
- [x] Existing CLI/TUI help, output, dispatch, initialization, and exit behavior remain supported.

## Verification Commands

Run every command from the target implementation repository `../ultraplan-go` unless the check says otherwise.

| Check | Command | Expected Result |
| --- | --- | --- |
| App serve/composition regression | `go test ./internal/app -run 'Serve|Surface|Dispatch|TUI|Web'` | Help, preflight, runner/query composition, cancellation, exit mapping, and CLI/TUI non-regression tests pass. |
| Web lifecycle | `go test ./internal/web -run 'Server|Lifecycle|Listen|Shutdown|Health'` | Loopback, listener, timeout, health, cancellation, graceful shutdown, and cleanup tests pass. |
| Web routes/API | `go test ./internal/web -run 'Route|API|Method|Head|Envelope'` | Exact route/method/HEAD/Allow/envelope/unknown-API behavior passes. |
| Web security | `go test ./internal/web -run 'Security|Host|Origin|Limit|Header|Redact|RequestID'` | Host/Origin, no-CORS, limits, headers, safe errors, IDs, and leak assertions pass. |
| Artifact preview | `go test ./internal/web -run 'Artifact|Preview|Traversal|Symlink|Hostile|Truncat'` | Opaque refs, containment, allowlist, bounds, escaping, and hostile-content tests pass. |
| Templates/frontend | `go test ./internal/web -run 'Template|Render|Accessibility|Static'` | Embedded parsing, semantic rendering, escaping, states, and asset-policy tests pass. |
| Full deterministic suite | `go test ./...` | All packages pass without external runtime, smoke harness, browser, Node, database, or network dependencies. |
| Race and lifecycle safety | `go test -race ./...` | Suite passes with no detected races or lifecycle leaks. |
| Single-binary build | `go build ./cmd/ultraplan` | Go binary builds with embedded templates/assets and no frontend build step. |
| Web dependency direction | `go list -deps ./internal/web` | Review confirms `internal/app` plus standard/support dependencies only, with no direct product/runtime/process/CLI-handler dependency. |
| Serve help | `go run ./cmd/ultraplan serve --help` | Help accurately documents listen/open-browser/workspace and shutdown behavior without starting a listener. |

Focused test regexes may be adjusted to the final test names, but the named behavior and files in `requirements.md` must remain covered and the full commands are mandatory.

## Risks And Blockers

| Risk / Blocker | Source | Mitigation | Status |
| --- | --- | --- | --- |
| Phase 3 is incomplete or exceptions are not recorded. | `reasoning.md#assumptions-and-risks` | Confirm readiness before implementation; block rather than duplicate Phase 3 behavior. | closed: Sprint 29 records the user-directed live-dogfood exception and passing fake-backed release gates |
| Required typed app projections are missing. | `reasoning.md#assumptions-and-risks` | Add typed app queries while preserving product ownership; never parse CLI output or import product modules from web. | closed: query-only facade added in `internal/app` |
| Default port `127.0.0.1:8080` is occupied. | `reasoning.md#assumptions-and-risks` | Return actionable startup failure and document `--listen`; never silently choose another port. | mitigated and documented |
| Local pages/processes probe disclosed state. | `reasoning.md#assumptions-and-risks` | Exact loopback/Host/Origin, no CORS, opaque refs, headers, safe errors, redaction, and documented local trust model. | mitigated and tested |
| Strict IP-literal policy surprises `localhost` users. | `reasoning.md#assumptions-and-risks` | Print/open canonical bound URL and document accepted forms; do not relax without new security reasoning. | mitigated and documented |
| Live refresh amplifies reads or sequential panels are slow. | `reasoning.md#trade-off-and-debt-analysis` | Fixed bounds, explicit refresh, cancellation, duration diagnostics, no recursive source scans; measure before cache/fan-out. | mitigated; measurement carried forward |
| Opaque references become stale. | `reasoning.md#assumptions-and-risks` | Return safe 404 and require refresh; never redirect to another artifact. | mitigated and documented |
| Symlink/encoding bypass discloses files. | `reasoning.md#assumptions-and-risks` | App-owned canonical resolution and exhaustive traversal/symlink tests; no broad web reader. | mitigated and tested |
| `/api/v1` DTOs drift with browser changes. | `reasoning.md#trade-off-and-debt-analysis` | Explicit DTOs, semantic schema tests, narrow envelope goldens, compatibility review. | mitigated; compatibility duty carried forward |
| Escaping is bypassed through trusted types or DOM APIs. | `reasoning.md#assumptions-and-risks` | Prohibit workspace `template.HTML`/`innerHTML`, use CSP and hostile-content tests. | mitigated and tested |
| Composition changes regress CLI/TUI behavior or startup. | `reasoning.md#assumptions-and-risks` | Lazy construction, command-path non-regression tests, Architecture Review. | closed by focused and full regression suites |
| Listener, handler, or launcher work outlives shutdown. | `reasoning.md#assumptions-and-risks` | Explicit ownership, propagated context, 10-second bound, joined cleanup, lifecycle/race tests. | mitigated; lifecycle and race suites pass |
| Fixed 200-entry/256 KiB bounds do not fit real workspaces. | `reasoning.md#trade-off-and-debt-analysis` | Expose/document truncation and bytes/counts; tune only from measured evidence with compatibility review. | accepted foundation trade-off |
| Broader PRD/TRD operation and SSE scope leaks into Sprint 30. | `reasoning.md#assumptions-and-risks` | Query-only injection and explicit absence tests for POST/DELETE/operations/confirmations/SSE/UI controls. | closed by capability/import/route review |
| Manual browser checks miss browser-specific defects. | `reasoning.md#potential-technical-debt` | Record keyboard/responsive/zoom/no-JavaScript evidence; consider browser automation in Sprint 32 hardening. | carried forward to Sprint 32; semantic/static review complete |

## Assumptions And Open Questions

| Item | Disposition |
| --- | --- |
| Phase 3 completion or recorded exceptions | Confirmed through Sprint 29's explicit user-directed live-dogfood exception and passing fake-backed release commands. |
| Existing app/product packages can expose required read projections without semantic changes | Confirmed; the new facade composes existing results and adds a non-persisting sprint status option without changing CLI/TUI behavior. |
| `127.0.0.1:8080` is an acceptable default | Adopted by Decision 2; collisions are explicit startup errors. |
| Implementation owner and dates | Assigned and recorded in Sprint Status. |
| Public API support, pagination, caching, partial aggregates, configurable bounds, rich Markdown, remote serving, auth, operations, confirmation, cancellation, and SSE | Deliberately deferred, not implementation questions for Sprint 30. |

No technical design question remains open for Sprint 30. If implementation evidence invalidates a decided constant or boundary, stop and record a deviation/recommendation for reasoning review rather than choosing a replacement ad hoc.

## Deferred Scope And Revisit Triggers

| Deferred Scope | Revisit Trigger |
| --- | --- |
| Guarded validation/workflow operations, confirmation, cancellation, SSE | Sprint 31 begins after the read-only foundation is accepted and defines separate command/security/lifecycle capabilities. |
| Authentication, TLS, remote/LAN bind, users/tenants | A governed remote-service phase supplies a threat and identity model. |
| Cache, snapshot, watcher, polling | Measured read amplification or latency exceeds targets with explicit staleness/invalidation semantics. |
| Pagination/search/filtering | Representative workspaces repeatedly exceed the 200 useful-entry bound. |
| App-level fan-out or partial panel errors | Measurements show independent slow panels and a per-section error/cancellation contract is approved. |
| Rich Markdown | A parser/sanitizer and hostile-content contract are selected through frontend/security reasoning. |
| Frontend framework/build pipeline | Demonstrated client complexity justifies a governed architecture change. |
| Public local API guarantees | Product explicitly commits to external clients, compatibility policy, pagination, and security support. |
| Browser automation | Sprint 32 hardening or an observed accessibility/browser regression warrants it. |

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

## Execution Log

| Date / Step | Action | Evidence / Notes |
| --- | --- | --- |
| 2026-07-29 / readiness | Assigned execution owner and confirmed Phase 3 readiness. | Sprint 29 `execute.md` records implementation complete, passing test/race/build/vet gates, and the explicit user-directed exception that live OpenCode/external-harness dogfood was omitted. |
| 2026-07-29 / implementation | Executed Tasks 2-9 in dependency order. | Added explicit surface composition, typed web queries, serve lifecycle, exact HTML/API routes, security middleware, opaque previews, embedded frontend, focused tests, and three implementation-repository docs. |
| 2026-07-29 / verification | Executed Task 10 and prepared review handoff. | Focused suites, full tests, race, build, vet, diff/whitespace, dependency/capability audit, help, and live loopback health/dashboard/origin/shutdown checks passed. |

## Execution Evidence

- Focused app: `go test ./internal/app -run 'Serve|Surface|Dispatch|TUI|Web'` — pass.
- Focused web lifecycle/routes/security/artifact/templates: all six planned `go test ./internal/web -run ...` commands — pass.
- Full deterministic suite: `go test ./...` — pass.
- Race suite: `go test -race ./...` — pass.
- Binary: `go build ./cmd/ultraplan` — pass; generated binary moved outside the repository.
- Static checks: `go vet ./...`, `git diff --check`, and new-file trailing-whitespace scan — pass.
- Help: `go run ./cmd/ultraplan serve --help` — pass and documents workspace, listen, browser opening, read-only scope, and shutdown.
- Dependency review: direct `internal/web` imports contain only `internal/app` plus the standard library; no direct product/runtime/process/CLI-handler dependency exists.
- Capability review: no POST/PUT/PATCH/DELETE operation route, confirmation, cancellation endpoint, SSE, WebSocket, runtime/process call, direct web filesystem read, browser storage, polling, Node/Vite/build pipeline, or third-party asset exists.
- Live command path: the built binary served the real planning workspace on `127.0.0.1:18080`; health returned `200`/`ok`, the dashboard returned semantic HTML, cross-origin health returned `403`, SIGTERM shut down cleanly, and normalized lifecycle/request diagnostics were recorded under `/tmp/ultraplan-s30-live.LlLoDR`.
- Accessibility/responsive static review: complete HTML without JavaScript; unique title/`h1`, skip link, primary navigation, breadcrumbs, `main`, footer/status, labelled tables/code and live region; visible `:focus-visible`; single-column narrow media query; local table/code overflow; relative sizing suitable for 200-percent zoom; reduced-motion query; hostile-content escaping tests. Browser-specific automation remains the explicitly accepted Sprint 32 follow-up.
- Architecture Review protocol result: good fit; explicit composition, query-only capability, transport/product ownership, non-mutating web reads, bounded lifecycle, typed boundary errors, and existing CLI/TUI protection conform to Decisions 1-8. Complexity increased as required for a third interface and remains localized.
- Sprint Review protocol handoff: governed inputs, implementation scope, selected contracts/protocols, plan execution, command evidence, deviations, deferrals, and risk dispositions are now present. Formal automated `review.md` generation remains the next governed workflow stage and was not substituted by execute-stage code.

## Deviations And Deferrals

- No implementation deviation from `reasoning.md` was required.
- Browser automation was not introduced because Sprint 30 permits semantic/manual evidence and explicitly assigns representative browser automation to Sprint 32 hardening. Residual browser-specific layout risk is carried forward.
- Guarded operations, confirmations, cancellation endpoints, SSE, remote/auth behavior, pagination, caching, rich Markdown, and frontend frameworks remain deferred exactly as planned.

## Review And Sign-Off

- **Implementation verdict:** complete; all required implementation and verification work passed.
- **Architecture review:** good fit; no blocker/high-severity finding identified in the protocol walkthrough.
- **Sprint review readiness:** ready for the governed review stage with current implementation diff and evidence.
- **Blockers:** none.
- **Completion time:** `2026-07-29T07:30:40Z`.

## Completion Criteria

- [x] All tasks are complete or explicitly deferred with requirement and reasoning impact recorded.
- [x] Every acceptance criterion and selected contract has implementation and evidence.
- [x] Verification commands were run successfully or deferrals are documented as blockers.
- [x] Expected evidence from `reasoning.md` is present and reviewable.
- [x] Exact routes, limits, errors, security policy, state ownership, frontend behavior, and scope exclusions conform to Decisions 1-8.
- [x] Documentation, help, and implementation agree on behavior and limitations.
- [x] Architecture Review and Sprint Review protocols have the required inputs and evidence.
- [x] `review.md` can evaluate conformance without guessing intent.
