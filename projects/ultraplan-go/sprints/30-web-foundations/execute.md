# Sprint Execution: Local Web Foundation and Read-Only Dashboard

> Project: `ultraplan-go`
> Sprint: `30-web-foundations`
> Status: `implementation complete; ready for review`
> Executed: `2026-07-29`
> Execution evidence reconciled: `2026-08-15`

## Target

- Repository: `/home/antonioborgerees/coding/ultraplan-go`
- Starting revision: `f8446a6` (`feat: complete phase 3 release hardening`)
- Starting state: clean at the recorded revision
- Ending state: Sprint 30 implementation changes present and uncommitted

## Readiness

Phase 3 readiness was confirmed before implementation. Sprint 29 recorded complete deterministic test, race, build, vet, and formatting gates, with the user-directed exception that live OpenCode and external-harness dogfood were omitted. Sprint 30 did not duplicate or alter Phase 3 workflow behavior.

## Implemented

- Added explicit CLI, TUI, and web surface composition without mutable global runner registration.
- Added typed, query-only application use cases for dashboard, project, sprint, study, validation, artifact preview, and health projections.
- Added `ultraplan serve` parsing, workspace/config preflight, loopback-only listen validation, optional browser launch, cancellation handling, and graceful shutdown.
- Added a bounded local HTTP lifecycle with fixed timeouts and request concurrency limits.
- Added server-rendered read-only HTML routes and versioned `/api/v1` JSON routes with stable success and structured error envelopes.
- Added exact Host and Origin enforcement, no-CORS defaults, request limits, security headers, request IDs, redacted diagnostics, and safe error projection.
- Added opaque, app-owned artifact references with containment, allowlisting, symlink-escape protection, and a 256 KiB preview bound.
- Added embedded Go templates, one stylesheet, and a small dependency-free progressive-refresh script with no frontend build pipeline.
- Added deterministic app and web tests covering composition, lifecycle, routing, API behavior, security, artifact safety, hostile content, templates, accessibility semantics, and state freshness.
- Added local web documentation and updated the CLI reference and implementation architecture documentation.

The browser surface remains read-only. It does not expose workflow mutation, runtime execution, review or smoke invocation, cancellation endpoints, SSE, WebSockets, remote binding, authentication, browser-owned durable state, arbitrary editing, or Git operations.

## Changed Paths

```text
cmd/ultraplan/main.go
docs/architecture.md
docs/cli-reference.md
docs/local-web.md
internal/app/app.go
internal/app/serve_commands.go
internal/app/serve_commands_test.go
internal/app/sprint_usecases.go
internal/app/surfaces.go
internal/app/usecases.go
internal/app/web_usecases.go
internal/app/web_usecases_test.go
internal/sprint/handbook.go
internal/sprint/handbook_test.go
internal/sprint/service.go
internal/web/artifacts.go
internal/web/artifacts_test.go
internal/web/handlers.go
internal/web/routes.go
internal/web/routes_test.go
internal/web/security.go
internal/web/security_test.go
internal/web/server.go
internal/web/server_test.go
internal/web/static/app.css
internal/web/static/app.js
internal/web/templates/artifact.html
internal/web/templates/dashboard.html
internal/web/templates/error.html
internal/web/templates/project.html
internal/web/templates/projects.html
internal/web/templates/shell.html
internal/web/templates/sprint.html
internal/web/templates/studies.html
internal/web/templates/study.html
internal/web/templates_test.go
internal/web/test_fakes_test.go
```

The `internal/sprint/handbook*` paths were present in the reconciled working tree even though the original execution run-state file did not list them. Review should determine whether they belong to this sprint or are adjacent user work; this execution summary does not claim ownership beyond recording the observed tree.

## Verification

All implementation verification was run from `/home/antonioborgerees/coding/ultraplan-go`.

| Command or check | Result |
| --- | --- |
| `go test ./internal/app -run 'Serve|Surface|Dispatch|TUI|Web'` | pass |
| `go test ./internal/web -run 'Server|Lifecycle|Listen|Shutdown|Health'` | pass |
| `go test ./internal/web -run 'Route|API|Method|Head|Envelope'` | pass |
| `go test ./internal/web -run 'Security|Host|Origin|Limit|Header|Redact|RequestID'` | pass |
| `go test ./internal/web -run 'Artifact|Preview|Traversal|Symlink|Hostile|Truncat'` | pass |
| `go test ./internal/web -run 'Template|Render|Accessibility|Static'` | pass |
| `go test ./internal/web -run ...` | pass; this is the plan's shorthand for the five exact focused web commands listed above |
| `go test ./...` | pass |
| `go test -race ./...` | pass |
| `go build ./cmd/ultraplan` | pass |
| `go vet ./...` | pass |
| `git diff --check` | pass |
| `go run ./cmd/ultraplan serve --help` | pass |
| `go list -deps ./internal/web` | pass; direct imports contain only `internal/app` and the standard library |
| Live loopback health, dashboard, cross-origin rejection, and SIGTERM shutdown check | pass |

The full test, race, build, vet, and diff checks were independently reconfirmed against the current implementation tree on `2026-08-15`.

## Execution Evidence

- The server served the real planning workspace on `127.0.0.1:18080`.
- `/api/v1/health` returned `200` with healthy status and the dashboard returned semantic HTML.
- A cross-origin health request returned `403`.
- `SIGTERM` produced a clean shutdown.
- Direct `internal/web` imports contain only `internal/app` and the Go standard library.
- No operation routes or frontend/runtime capabilities excluded by the sprint were found.
- Static accessibility review covered no-JavaScript behavior, keyboard/focus visibility, semantic landmarks, narrow reflow, local overflow, zoom, reduced motion, and hostile-content escaping.

## Deviations, Deferrals, and Risks

- No implementation deviation from `reasoning.md` was recorded.
- Browser automation remains deferred to Sprint 32; Sprint 30 uses semantic route/template tests and manual accessibility evidence.
- Guarded operations, confirmation, explicit cancellation endpoints, and SSE remain deferred to Sprint 31.
- Remote access, authentication, pagination, caching, rich Markdown, and a frontend framework remain outside this sprint.
- The implementation is still uncommitted, so review must inspect the current working-tree diff rather than a stable commit.
- The historical `.run-state.json` predates the current task-record schema. It remains valid legacy completion evidence through the CLI's compatibility path but should not be represented as a newly generated governed runtime execution.

## Completion

All ten plan tasks and the plan's evidence checklist are complete. No implementation blocker is recorded. The execution stage is complete and the next governed stage is Sprint Review, which must generate `review.md` before the sprint can receive final sign-off.
