# Sprint Requirements: Local Web Hardening and Observable-Product Release

> Project: `ultraplan-go`
> Sprint: `32-hardening-and-release`
> Purpose: the authoritative, human-readable sprint contract. All other sprint artifacts must satisfy these requirements.

## Sprint Goal

Turn the Sprint 30-31 local browser implementation into a secure, accessible, compatibility-controlled, recoverable, and release-ready interface while proving that new workflow stages can use the shared application capability without route-specific workflow logic.

## Required Outputs

| Output | Path | Description |
| --- | --- | --- |
| Stable API routes and envelopes | `../ultraplan-go/internal/web/routes.go` | Preserve and stabilize the documented `/api/v1` routes, success envelopes, typed error envelopes, methods, and status mappings. |
| HTTP and page handlers | `../ultraplan-go/internal/web/handlers.go` | Map validated HTTP requests and typed app results into compatibility-controlled JSON DTOs and explicit page view models. |
| Operation HTTP handlers | `../ultraplan-go/internal/web/operation_handlers.go` | Harden guarded operation preparation, start, status, cancellation, recovery, and SSE transport mapping. |
| Browser security policy | `../ultraplan-go/internal/web/security.go` | Enforce loopback, Host/Origin, CSRF, session, CSP/header, body-limit, containment, and redaction policy. |
| Operation and SSE hub | `../ultraplan-go/internal/web/operations.go` | Maintain bounded operation, event, subscriber, retention, cancellation, and terminal-result behavior without durable web-owned state. |
| Server lifecycle | `../ultraplan-go/internal/web/server.go` | Implement draining, exact-once cancellation, bounded cleanup, truthful durable reconciliation, and graceful shutdown ordering. |
| Shared web application capability | `../ultraplan-go/internal/app/web_usecases.go` | Expose status, artifacts, commands, progress, cancellation, and recovery through typed shared application abstractions. |
| Shared operation capability | `../ultraplan-go/internal/app/operations.go` | Keep operation normalization, confirmation scope, execution, progress, and safe results reusable by CLI, TUI, and web. |
| Template primitives | `../ultraplan-go/internal/web/templates/primitives/primitives.html` | Define namespaced presentation-only atoms used by higher template layers. |
| Template components | `../ultraplan-go/internal/web/templates/components/components.html` | Define namespaced reusable UI compositions for navigation, state, findings, progress, artifacts, and errors. |
| Template layouts | `../ultraplan-go/internal/web/templates/layouts/layouts.html` | Define namespaced shared page shells and arrangements composed only from primitives and components. |
| Route pages | `../ultraplan-go/internal/web/templates/pages/pages.html` | Define namespaced route-level pages using explicit typed view models and downward-only template composition. |
| CSS tokens | `../ultraplan-go/internal/web/static/css/tokens.css` | Define the browser presentation design tokens. |
| CSS base | `../ultraplan-go/internal/web/static/css/base.css` | Define document defaults, typography, focus foundations, and reduced-motion behavior. |
| CSS primitives | `../ultraplan-go/internal/web/static/css/primitives.css` | Style primitive controls and content atoms. |
| CSS components | `../ultraplan-go/internal/web/static/css/components.css` | Style reusable status, progress, finding, artifact, confirmation, and error components. |
| CSS layouts | `../ultraplan-go/internal/web/static/css/layouts.css` | Provide responsive dashboard and detail-page layouts. |
| CSS utilities | `../ultraplan-go/internal/web/static/css/utilities.css` | Provide narrowly scoped presentation utilities without encoding product state. |
| Baseline browser enhancement | `../ultraplan-go/internal/web/static/js/app.js` | Provide dependency-free progressive enhancement without a client router or authoritative store. |
| Guarded-operation enhancement | `../ultraplan-go/internal/web/static/js/operations.js` | Enhance prepare, confirm, start, cancellation, result, and durable-refresh behavior. |
| SSE enhancement | `../ultraplan-go/internal/web/static/js/sse.js` | Handle bounded SSE connection, replay, reconnect, gaps, terminal events, and fallback refresh. |
| API compatibility fixtures | `../ultraplan-go/internal/web/api_compatibility_test.go` | Freeze documented `/api/v1` response fields, envelopes, methods, status codes, and typed error codes. |
| Security regression tests | `../ultraplan-go/internal/web/security_test.go` | Cover Host/Origin, CSRF, session, CSP/headers, request limits, hostile input, paths, Markdown, and redaction. |
| Lifecycle and shutdown tests | `../ultraplan-go/internal/web/server_test.go` | Cover draining, active-operation cancellation, bounded cleanup, terminal persistence, restart reconciliation, and disconnect isolation. |
| Operation and concurrency tests | `../ultraplan-go/internal/web/operations_test.go` | Cover bounds, races, leaks, slow subscribers, retention, exact-once cancellation, and terminal arbitration. |
| SSE protocol tests | `../ultraplan-go/internal/web/sse_test.go` | Cover event ordering, IDs, replay, rollover, heartbeat, reconnect, slow clients, disconnects, and terminal flush. |
| Template and accessibility tests | `../ultraplan-go/internal/web/templates_test.go` | Cover hierarchy, namespaced definitions, complete no-JavaScript rendering, escaping, keyboard semantics, focus, labels, live regions, and embedded paths. |
| Representative browser integration tests | `../ultraplan-go/internal/web/integration_test.go` | Exercise substantial study and sprint workflows over temporary workspaces with fake runtime and harness dependencies. |
| Interface capability tests | `../ultraplan-go/internal/app/web_usecases_test.go` | Prove shared app abstractions expose status, artifacts, operations, progress, cancellation, and recovery without route-specific workflow code. |
| Local web guide | `../ultraplan-go/docs/local-web.md` | Document startup, browser workflows, API/SSE contracts, limits, security, accessibility, shutdown, recovery, and troubleshooting. |
| Configuration reference | `../ultraplan-go/docs/configuration.md` | Document local-web configuration, precedence, loopback restrictions, limits, and safe diagnostics. |
| CLI reference | `../ultraplan-go/docs/cli-reference.md` | Document `ultraplan serve`, flags, lifecycle, browser/API scope, and CLI/TUI/web state agreement. |
| User guide | `../ultraplan-go/docs/user-guide.md` | Document supported browser inspection and guarded operation journeys. |
| Recovery guide | `../ultraplan-go/docs/recovery.md` | Document reconnect, cancellation, interrupted work, cleanup uncertainty, stale locks, restart reconciliation, and durable-state recovery. |
| Architecture guide | `../ultraplan-go/docs/architecture.md` | Document web-to-app dependency direction, presentation hierarchy, operation ownership, durable truth, and extensibility boundary. |
| Release checklist | `../ultraplan-go/docs/release-checklist.md` | Record deterministic, race, build, packaging, security, accessibility, fake integration, and gated real-system release checks. |

## Acceptance Criteria

- [ ] `go test ./...`, `go test -race ./...`, and `go build ./cmd/ultraplan` pass from `../ultraplan-go`.
- [ ] `ultraplan serve` binds only to numeric loopback addresses and runs from the single Go binary without Node.js, Vite, a separate frontend process, or a database.
- [ ] Compatibility fixtures cover every documented `/api/v1` success envelope, error envelope, route method, status mapping, and stable machine-readable error code introduced through Sprint 32.
- [ ] Unknown `/api/` routes and unsupported methods return structured JSON errors and never fall through to HTML pages.
- [ ] HTML, static assets, and API responses have documented and tested cache policies; any polling or refetch used to observe CLI/TUI changes is bounded and cannot become product state.
- [ ] Browser, CLI, and TUI return equivalent durable state, readiness, verdicts, artifacts, terminal outcomes, and next actions for the representative fixtures.
- [ ] One substantial study workflow and one substantial sprint workflow can be prepared, confirmed, started, observed, cancelled or completed, and recovered through the browser using shared app operations.
- [ ] The interface capability test demonstrates that adding a stage can expose status, artifacts, commands, progress, cancellation, and recovery without adding route-specific workflow orchestration.
- [ ] `internal/web` imports `internal/app` and standard-library packages only; it does not import product, runtime, process, or CLI-handler packages or invoke `ultraplan` as a subprocess.
- [ ] Templates are organized into primitives, components, layouts, and pages with namespaced definitions and downward-only composition, and invalid or duplicate definitions fail server startup.
- [ ] Handlers construct explicit typed view models; templates perform no filesystem reads, app calls, HTTP validation, durable mutation, or product-state interpretation.
- [ ] Representative pages remain complete and operable without JavaScript; JavaScript is dependency-free progressive enhancement with no client router, persistent workflow store, or authoritative state.
- [ ] Dashboard, detail, confirmation, progress, finding, artifact, and error views pass keyboard-navigation, focus visibility, semantic landmark, control-label, error-association, live-region, color-independent-state, reduced-motion, zoom, and narrow-layout checks.
- [ ] Host/Origin, CSRF, session, CSP/security-header, request-smuggling/body-limit, hostile-Markdown, path-containment, forged-reference, and secret-redaction regression tests pass.
- [ ] No token, cookie, CSRF value, prompt, unsafe path, raw provider payload, raw stderr, or representative secret fixture appears in JSON, SSE, HTML, retained events, terminal results, or captured diagnostics.
- [ ] Operation and SSE tests prove all configured operation, preparation, event, payload, result, subscriber, stream, retention, heartbeat, and lifetime bounds under concurrent and slow-client scenarios.
- [ ] Browser refresh, tab close, navigation, SSE loss, or a slow subscriber does not cancel, complete, fail, or block product work.
- [ ] Graceful shutdown enters draining, rejects new mutations, requests canonical cancellation exactly once for every server-owned active operation, waits for bounded cleanup outside hub locks, and closes HTTP/SSE only after a durable terminal or cleanup-uncertain outcome is recorded.
- [ ] The Sprint 31 gap is closed: shutdown deadline exhaustion persists owner-specific durable `cleanup_uncertain` or equivalent truthful uncertainty before exit, and restart reconciliation does not infer success from process absence or artifact presence.
- [ ] Race, leak, cancellation, reconciliation, forced-interruption, restart, replay-gap, and slow-client tests pass without orphaned goroutines, locks, subscribers, or detached operations.
- [ ] Templates, CSS, and JavaScript are embedded and resolvable from the built binary, and packaging checks prove no source-tree asset dependency at runtime.
- [ ] Local-web user, configuration, API, security, recovery, troubleshooting, architecture, accessibility, and packaging documentation matches the implemented routes, fields, limits, states, and commands.
- [ ] Gated browser evidence exercises one real runtime-backed operation and one real smoke-harness-backed operation when the required environment is available; unavailable prerequisites produce a documented blocked result, never a pass.
- [ ] Sprint review reports no unresolved high-severity local-web security, path, redaction, concurrency, cancellation, shutdown, API compatibility, or accessibility finding.

## Non-Goals

- Adding `code-context` or any other new planning, execution, review, smoke, QA, repair, or retrieval stage.
- Hosted SaaS, LAN/public binding, remote access, remote workers, accounts, authentication services, authorization roles, tenancy, teams, or multi-user collaboration.
- A client-side router, framework, package manager, build pipeline, authoritative browser store, service worker, WebSocket protocol, or bidirectional agent chat.
- A database, alternate durable web state, detached worker, operation queue, or second workflow scheduler.
- Browser editing of arbitrary workspace files or automatic product/test repair.
- General-purpose issue tracking, assignment, scheduling, remote issue synchronization, or project-management features.
- Automatic Git add, commit, push, branch, merge, reset, or other Git-state mutation.
- Product Phase 5 content identity, provenance, retrieval, SQLite, knowledge graph, cloud, or Aren capabilities.

## Constraints

- `internal/web` may depend only on typed `internal/app` abstractions and the Go standard library; product workflow and durable recovery rules remain in their owning modules.
- Workspace artifacts and product-owned run state are authoritative; browser state, operation handles, event buffers, subscribers, and confirmations are bounded and ephemeral.
- The server must remain loopback-only, same-origin, fail-closed, path-contained, body/time/stream bounded, and secret-redacting across HTML, JSON, SSE, logs, and retained operation data.
- HTTP POST/DELETE carries commands and cancellation; SSE is one-way progress only and cannot mutate state or determine task success.
- Server-owned operations must share product-owned mutation exclusion, canonical cancellation, process-tree cleanup, and durable reconciliation with CLI and TUI operations.
- Browser disconnect cancels only the subscription; explicit cancellation or server shutdown is required to cancel product work.
- Server shutdown must not detach active work or release product locks before durable terminal reconciliation or truthful cleanup uncertainty is recorded.
- Templates must use `html/template`, stable namespaced definitions, explicit typed view models, and downward-only `page -> layout -> component -> primitive` composition.
- CSS must remain layered as tokens, base, primitives, components, layouts, and utilities; JavaScript must remain dependency-free and split by narrow progressive-enhancement capability.
- Normal tests must use `httptest`, temporary workspaces, fake app/runtime/harness dependencies, deterministic clocks/IDs/barriers, and no live provider requirement.
- Real runtime and smoke-harness checks must remain gated and must report missing prerequisites as blocked rather than skipped success.
- No automatic Git mutation is permitted, and release verification must preserve unrelated working-tree changes.

## Dependencies

| Prior Sprint / Output | Required For | Notes |
| --- | --- | --- |
| Sprint 30 implementation and `projects/ultraplan-go/sprints/30-web-foundations/review.md` | Read-only server, browser dashboard, security baseline, and prior findings | Preserve the web-to-app boundary and close the remaining documentation-verification findings. |
| Sprint 31 planning artifacts through `projects/ultraplan-go/sprints/31-web-operations/plan.md` | Guarded operations, SSE, mutation exclusion, and shutdown contract | Planning content defines the operation and lifecycle behavior hardened here. |
| `projects/ultraplan-go/sprints/31-web-operations/execute.md` | Carry-forward implementation state | Eight tasks are complete; owner-specific shutdown uncertainty and exhaustive compatibility/review evidence remain mandatory Sprint 32 work. |
| `projects/ultraplan-go/sprints/28-review-to-smoke-flow/review.md` | Release redaction, observability, recovery, and documentation audit | Recheck still-applicable blocker/high findings rather than assuming later interface work resolved them. |
| `projects/ultraplan-go/docs/PRD.md` | Product behavior and Phase 4 definition of done | Browser operation, shutdown, recovery, security, and interface-agreement requirements are authoritative. |
| `projects/ultraplan-go/docs/TRD.md` | HTTP/SSE, security, lifecycle, template, and testing constraints | Sections 7.5 and 18A define the binding local-web technical contract. |
| `projects/ultraplan-go/docs/ARCHITECTURE.md` | Package ownership and presentation hierarchy | `internal/web` remains an interface adapter over shared app use cases. |
| `/home/antonioborgerees/coding/ultraplan/ultraplan-go-smoke/` | Gated real browser-operation evidence | Use only through the cataloged harness protocol; detailed evidence remains in the harness. |

## Review Expectations

| What | How Verified |
| --- | --- |
| API compatibility | Golden/fixture comparison of every documented `/api/v1` route, method, envelope, field, status, and typed error code. |
| Cross-surface agreement | Integration tests over identical temporary workspaces comparing app, CLI, TUI, HTML, and JSON projections. |
| Shared capability extensibility | A fake stage-capability fixture proves status, artifact, command, progress, cancellation, and recovery exposure without route-specific workflow code. |
| Architecture boundaries | Import inspection and code review confirm `internal/web -> internal/app` only and no web-owned product workflow or durable state. |
| Presentation hierarchy | Startup parse tests and representative render tests verify namespaced primitives, components, layouts, pages, explicit view models, and downward-only composition. |
| Accessibility | Automated semantic/render assertions plus keyboard-only, focus, reduced-motion, zoom, reflow, and color-independent-state review of representative pages. |
| Browser security | Host/Origin, CSRF, session, CSP/header, request-limit, hostile-Markdown, path escape, forged reference, and redaction test suites plus manual audit. |
| Operation and SSE safety | Deterministic concurrency, race, leak, slow-subscriber, replay, rollover, heartbeat, disconnect, cancellation, and terminal-order tests. |
| Shutdown and recovery | Multi-operation shutdown and abrupt-restart tests prove draining, exact-once cancellation, bounded cleanup, durable uncertainty, lock ordering, and conservative reconciliation. |
| Packaging | Build and launch the binary outside the source asset tree; verify all templates and static assets load and no frontend runtime is required. |
| Documentation | Review each public route, field, state, bound, flag, recovery action, and security limitation against implementation and compatibility fixtures. |
| Release quality | Run `go test ./...`, `go test -race ./...`, `go build ./cmd/ultraplan`, focused web/app tests, governed sprint review, and gated real-system checks where available. |
