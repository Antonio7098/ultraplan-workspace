# Sprint Requirements: Local Web Foundation and Read-Only Dashboard

> Project: `ultraplan-go`
> Sprint: `30-web-foundations`
> Purpose: the authoritative, human-readable sprint contract. All other sprint artifacts must satisfy these requirements.

## Sprint Goal

Add `ultraplan serve` as a loopback-only Go HTTP server that serves a read-only browser dashboard over existing typed app use cases without changing the governed workflow, workspace state model, CLI behavior, or TUI behavior.

## Required Outputs

| Output | Path | Description |
| --- | --- | --- |
| Sprint requirements | `projects/ultraplan-go/sprints/30-web-foundations/requirements.md` | This sprint contract derived from the project index, roadmap, PRD, TRD, architecture, and prior sprint reviews. |
| Architecture reasoning | `projects/ultraplan-go/sprints/30-web-foundations/reasoning/architecture.md` | Decision record for explicit CLI/TUI/web composition, `internal/web` boundaries, package dependency direction, and state ownership. |
| API design reasoning | `projects/ultraplan-go/sprints/30-web-foundations/reasoning/api-design.md` | Decision record for read-only HTML routes, `/api/v1` resources, structured JSON errors, health response, and artifact-preview semantics. |
| Frontend reasoning | `projects/ultraplan-go/sprints/30-web-foundations/reasoning/frontend.md` | Decision record for server-rendered templates, progressive enhancement, accessibility, hostile-content handling, and no frontend build pipeline. |
| Serve command wiring | `../ultraplan-go/cmd/ultraplan/main.go` | Constructs the web runner explicitly alongside CLI and TUI dependencies without package-global mutable runner registration. |
| App surface composition | `../ultraplan-go/internal/app/surfaces.go` | Defines explicit local-interface runner dependencies for TUI and web composition. |
| Web app use cases | `../ultraplan-go/internal/app/web_usecases.go` | Provides typed read-only dashboard, detail, validation, artifact-preview, health, and error-projection use cases used by HTTP handlers. |
| Serve command adapter | `../ultraplan-go/internal/app/serve_commands.go` | Parses `ultraplan serve` flags, resolves workspace/config, starts the web server, maps errors to exit codes, and preserves CLI output discipline. |
| Serve command tests | `../ultraplan-go/internal/app/serve_commands_test.go` | Covers serve help, loopback address validation, workspace resolution, startup failure, shutdown path, and non-regression of existing CLI/TUI command wiring. |
| Web server lifecycle | `../ultraplan-go/internal/web/server.go` | Owns loopback `net/http` lifecycle, graceful shutdown, timeouts, optional browser opening hook, and server health wiring. |
| Web routes | `../ultraplan-go/internal/web/routes.go` | Registers HTML routes, static asset routes, `/api/v1` read-only endpoints, health endpoint, and API-not-found handling. |
| Web handlers | `../ultraplan-go/internal/web/handlers.go` | Maps HTTP requests to app use cases and renders HTML or JSON without importing product modules or CLI handlers. |
| Web security middleware | `../ultraplan-go/internal/web/security.go` | Enforces loopback/Host/Origin policy, same-origin defaults, request limits, security headers, redaction, and safe error projection for the read-only foundation. |
| Artifact preview support | `../ultraplan-go/internal/web/artifacts.go` | Implements bounded, allowlisted Markdown/JSON artifact preview with workspace containment and hostile HTML/script safety. |
| Web templates | `../ultraplan-go/internal/web/templates/*.html` | Embedded `html/template` pages and partials for dashboard, project, sprint, study, validation, artifact preview, and error pages. |
| Web static assets | `../ultraplan-go/internal/web/static/*` | Embedded CSS and minimal JavaScript for navigation, progressive refresh, and non-operational browser behavior. |
| Web lifecycle tests | `../ultraplan-go/internal/web/server_test.go` | Uses `httptest` or equivalent fakes to cover startup, shutdown, loopback binding rules, timeouts, and health behavior. |
| Web route/API tests | `../ultraplan-go/internal/web/routes_test.go` | Covers HTML routes, `/api/v1` route shapes, unknown `/api/` JSON errors, method handling, and structured error envelopes. |
| Web security tests | `../ultraplan-go/internal/web/security_test.go` | Covers Host/Origin rejection, request-size limits, security headers, redaction, path containment, and unsupported bind addresses. |
| Artifact preview tests | `../ultraplan-go/internal/web/artifacts_test.go` | Covers allowlisted paths, escaping paths, unsupported extensions, size bounds, Markdown escaping, JSON preview, and hostile content. |
| Template tests | `../ultraplan-go/internal/web/templates_test.go` | Parses embedded templates and verifies representative dashboard/detail/error rendering with fake app data. |
| Local web user documentation | `../ultraplan-go/docs/local-web.md` | Explains `ultraplan serve`, loopback-only scope, browser dashboard, read-only limits, artifact previews, shutdown, and troubleshooting. |
| CLI reference update | `../ultraplan-go/docs/cli-reference.md` | Documents `ultraplan serve`, flags, loopback binding, exit behavior, and relationship to CLI/TUI. |
| Architecture documentation update | `../ultraplan-go/docs/architecture.md` | Documents the local web adapter boundary, `internal/web -> internal/app` dependency rule, and single-binary asset embedding. |

## Acceptance Criteria

- [ ] `ultraplan serve --help` documents the serve command, explicit loopback listen address option, optional browser-opening option, workspace selection behavior, and shutdown behavior.
- [ ] `ultraplan serve` starts from inside a workspace and with `--workspace <path>` using the same workspace discovery and config validation rules as existing CLI/TUI entry points.
- [ ] The server accepts only loopback listen addresses such as `127.0.0.1` or `::1`; non-loopback bind requests fail before serving requests.
- [ ] The server shuts down cleanly on context cancellation or process signal without leaving request handlers blocked.
- [ ] Browser HTML pages can inspect workspace overview, projects, project sprints, studies, validation summaries, current run or flow state, and bounded artifact previews through shared app use cases.
- [ ] Project, sprint, and study detail pages are backed by typed app query results, not by CLI subprocesses, CLI command handlers, stdout parsing, or direct imports of `internal/project`, `internal/sprint`, or `internal/study` from `internal/web`.
- [ ] The UI and server run from the Go binary using `html/template`, embedded CSS, and minimal JavaScript; there is no Node.js, Vite, separate asset server, database, React app, or frontend build step.
- [ ] Initial `/api/v1` read-only JSON endpoints return stable JSON objects for dashboard/detail/validation/artifact-preview/health data, and unknown `/api/` routes return structured JSON errors rather than HTML.
- [ ] Artifact preview accepts only allowlisted Markdown and JSON workspace artifacts, rejects escaping or unsupported paths, enforces size bounds, and never executes workspace HTML or scripts.
- [ ] Host and Origin checks, same-origin defaults, request limits, redaction, security headers, safe error projection, and path-containment rules are enforced by HTTP tests.
- [ ] `internal/web` owns only HTTP lifecycle, transport DTOs, templates/static assets, security middleware, and browser rendering concerns; it does not own product state machines, validators, prompts, runtime execution, smoke harness invocation, or durable persistence.
- [ ] Existing CLI and TUI behavior remain supported and unchanged except for the explicit composition refactor needed to add the web runner.
- [ ] Normal verification passes with deterministic fakes: `go test ./...`, `go test -race ./...`, `go build ./cmd/ultraplan`, focused `internal/app` serve-command tests, and focused `internal/web` route/security/template/artifact tests.

## Non-Goals

- Guarded browser operations, workflow mutation, runtime-backed actions, confirmation tokens for starting work, explicit cancellation endpoints, and SSE progress are not included; they belong to Sprint 31.
- Browser editing of arbitrary workspace files or generated artifacts is not included.
- Hosted SaaS, LAN/public binding, remote access, accounts, authentication teams, permissions, tenants, collaboration, remote workers, and remote workspace synchronization are not included.
- WebSockets, bidirectional terminal/session transport, agent chat, and remote process control are not included.
- A database-backed state model, browser-owned durable state, or web-specific product state is not included.
- General-purpose issue tracking, issue assignment, project-management features, automatic product fixes, and automatic Git mutation are not included.
- A frontend framework, Node.js runtime, Vite, separate frontend process, or asset build pipeline is not included.
- Changes to review, smoke, execute, study run-loop, or planning-stage semantics are not included except where read-only status/query use cases are needed for the dashboard.

## Constraints

- `internal/web` may depend on `internal/app` use-case interfaces and plain result types only; it must not import `internal/study`, `internal/project`, `internal/sprint`, runtime adapters, process adapters, or CLI command packages.
- HTTP handlers must call typed app use cases and must not invoke `ultraplan` as a subprocess, call CLI handlers, parse CLI output, or duplicate product validation/workflow logic.
- `cmd/ultraplan` must explicitly construct web dependencies; adding package-global mutable runner registration is prohibited.
- Workspace artifacts and product-owned run state remain authoritative across browser refresh, disconnect, and server restart; Sprint 30 may retain only ephemeral HTTP/request state.
- Sprint 30 is read-only from the browser: HTTP routes must not start, cancel, mutate, review, smoke, execute, or run runtime-backed workflows.
- The server must bind only to loopback in Phase 4 and must reject non-loopback bind addresses before accepting requests.
- Browser output must be rendered with Go `html/template`; untrusted Markdown or artifact content must be escaped or sanitized so embedded HTML/scripts do not execute.
- Static assets must be embedded in the Go binary and must not require a separate runtime, asset server, or build step.
- API routes under `/api/v1` must use structured JSON success/error responses; unknown `/api/` paths must not fall through to HTML rendering.
- Tests must use `httptest`, fake app use cases, fake runtimes, and fake harness data by default; real runtime or real smoke harness behavior is not required for Sprint 30 verification.
- Automatic Git add, commit, push, branch, merge, reset, or other Git mutation remains prohibited.

## Dependencies

| Prior Sprint / Output | Required For | Notes |
| --- | --- | --- |
| Product Phase 3 release gate or recorded exceptions | Sprint 30 readiness | Roadmap requires Phase 3 completion or explicit exception recording before Phase 4 starts. |
| Sprints 24-25 TUI foundation and operational controls | Shared local-interface architecture | The web surface must follow the established pattern of shared app use cases instead of CLI stdout integration. |
| Sprint 26 automated review use cases and follow-on Phase 3 review/smoke flow | Read-only verification status display | The dashboard must inspect current review/smoke state without changing review or smoke semantics. |
| `internal/app` typed query/error surfaces | HTTP handlers and dashboard data | Web handlers depend on app use cases for workspace, project, sprint, study, validation, artifact-preview, and status data. |
| `projects/ultraplan-go/docs/ARCHITECTURE.md` Phase 4 rules | Package layout and dependency direction | Defines `internal/web` ownership, dependency rules, and explicit composition requirements. |
| `projects/ultraplan-go/docs/TRD.md` section 7.5 | Server, browser UI, security, and testing requirements | Sprint 30 implements the read-only foundation subset of the Phase 4 local HTTP requirements. |
| `projects/ultraplan-go/docs/PRD.md` browser scenarios | Product behavior and non-goals | Confirms loopback-only browser inspection and defers hosted/multi-user/browser-owned state. |

## Review Expectations

| What | How Verified |
| --- | --- |
| Serve command exists and is documented | Inspect `ultraplan serve --help`, `../ultraplan-go/internal/app/serve_commands.go`, and `../ultraplan-go/docs/cli-reference.md`. |
| Loopback-only lifecycle works | Run focused `internal/web` lifecycle tests and review bind-address validation in `server.go`. |
| Explicit composition is package-cycle-free | Inspect `cmd/ultraplan`, `internal/app/surfaces.go`, imports, and `go test ./...` output. |
| Web boundary stays transport-only | Review `internal/web` imports and handlers to confirm dependency only on `internal/app` types plus standard/support packages. |
| Dashboard and detail pages use shared app use cases | Review handler-to-use-case calls and fake app tests for workspace, project, sprint, study, validation, and artifact-preview data. |
| `/api/v1` read-only API shape is stable | Run `internal/web/routes_test.go` and inspect JSON fixtures/assertions for success and error envelopes. |
| Unknown API routes return JSON errors | Run route tests covering unknown `/api/` paths and verify no HTML fallback. |
| Artifact previews are bounded and safe | Run `internal/web/artifacts_test.go` and inspect allowlist, containment, size, Markdown escaping, and hostile-content cases. |
| Security foundation is enforced | Run `internal/web/security_test.go` for Host/Origin, request limit, security headers, redaction, and path-containment cases. |
| No frontend runtime dependency exists | Inspect embedded templates/static assets, build files, and `go build ./cmd/ultraplan`; verify no Node/Vite/frontend process is required. |
| CLI and TUI regressions are avoided | Run `go test ./...`, focused `internal/app` tests, and review serve composition changes for unchanged existing command behavior. |
| Sprint verification passes | Run `go test ./...`, `go test -race ./...`, and `go build ./cmd/ultraplan` from `../ultraplan-go`. |
