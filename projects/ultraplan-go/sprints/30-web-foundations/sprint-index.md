# Sprint Index: Local Web Foundation and Read-Only Dashboard

> Project: `ultraplan-go`
> Sprint: `30-web-foundations`
> Purpose: selected context for this sprint. Must be a subset of `projects/ultraplan-go/project-index.md`.
> **Inputs Used:** `projects/ultraplan-go/project-index.md`, `projects/ultraplan-go/roadmap.md`, `projects/ultraplan-go/sprints/30-web-foundations/requirements.md`, `projects/ultraplan-go/docs/ARCHITECTURE.md`, `projects/ultraplan-go/docs/PRD.md`, `projects/ultraplan-go/docs/TRD.md`

This document selects what must be read, distilled, reasoned through, or checked for this sprint. It does not make implementation decisions. All selections must come from the project index; no items may be included that are not listed in the project index.

## Sprint Scope

- **Sprint Goal:** Add `ultraplan serve` as a loopback-only Go HTTP server that serves a read-only browser dashboard over existing typed app use cases without changing the governed workflow, workspace state model, CLI behavior, or TUI behavior.
- **Planned Output:** Explicit CLI/TUI/web composition; a documented `ultraplan serve` command; typed read-only web app use cases; loopback HTTP lifecycle, HTML and `/api/v1` routes, security middleware, bounded artifact previews, embedded templates and static assets; focused app/web tests; and local-web, CLI-reference, and architecture documentation.
- **Depends On:** Product Phase 3 release completion or recorded exceptions; Sprints 24-25 shared local-interface foundations; Sprint 26 and follow-on Phase 3 review/smoke status surfaces; existing typed `internal/app` query/error surfaces; and the Phase 4 rules in the project Architecture, PRD, and TRD.
- **Non-Goals:** Browser-triggered validation or workflow operations; confirmations, operation handles, cancellation endpoints, or SSE progress; browser editing; runtime or smoke-harness invocation; hosted, LAN/public, remote, multi-user, account, permission, tenant, collaboration, or remote-worker behavior; WebSockets or terminal transport; browser-owned durable state or a database; frontend frameworks or build pipelines; general-purpose issue tracking; automatic fixes; and Git mutation.

## Source Project Index

- `projects/ultraplan-go/project-index.md` is the authoritative source. Any file or item referenced below must appear there.

## Selected Contracts

Each contract applies as a flat whole to this sprint. All paths must appear in the project index's "Active Contract Pool" table.

| Contract | Why Selected |
| --- | --- |
| Architecture | Governs explicit interface composition, `internal/web -> internal/app` dependency direction, thin transport boundaries, package-cycle avoidance, and authoritative product state ownership. |
| Errors | Applies to serve-command exit mapping, structured API errors, safe HTTP error projection, method/not-found handling, startup failures, and actionable diagnostics. |
| Configuration | Applies to workspace discovery, command preflight, explicit listen/browser options, loopback-address validation, precedence, and redaction. |
| Observability | Applies to truthful health and status projections, request/server diagnostics, redaction, lifecycle visibility, and read-only display of existing run and flow state. |
| Security | Applies to loopback-only binding, Host/Origin and same-origin policy, request limits, security headers, workspace path containment, artifact allowlisting, hostile-content handling, and secret redaction. |
| Testing | Applies to deterministic app and HTTP fakes, `httptest`, route/security/template/artifact coverage, lifecycle and shutdown tests, race tests, and CLI/TUI non-regression checks. |
| Documentation | Applies to serve help, local-web usage and troubleshooting, CLI reference updates, architecture boundary documentation, and accurate read-only limitations. |
| CLI Surface | Applies to the `serve` command shape, flags, help, workspace selection, output discipline, shutdown behavior, and exit codes without changing existing CLI behavior. |
| Performance | Applies to bounded HTTP timeouts and request sizes, server lifecycle, embedded-asset startup, and bounded artifact reads without accidental large workspace scans. |

## Selected Evidence Reports

Copied from the project index's "Available Evidence Reports" table. These tell the technical handbook which reports to read – the project index is the authoritative source.

| Report | Path | Covers |
| --- | --- | --- |
| `01-project-structure` | `studies/go-cli-study/reports/final/01-project-structure.md` | Project layout, cmd/internal/pkg, dependency direction, thin entrypoints |
| `02-command-architecture` | `studies/go-cli-study/reports/final/02-command-architecture.md` | Command routing, flags, help text, command organization, shell completion |
| `03-dependency-injection` | `studies/go-cli-study/reports/final/03-dependency-injection.md` | Dependency construction, seams, testability, constructor patterns |
| `04-configuration-management` | `studies/go-cli-study/reports/final/04-configuration-management.md` | Config loading, precedence, environment variables, paths |
| `05-error-handling` | `studies/go-cli-study/reports/final/05-error-handling.md` | Error wrapping, classification, user-facing diagnostics, sentinel errors |
| `06-io-abstraction` | `studies/go-cli-study/reports/final/06-io-abstraction.md` | Filesystem/stdin/stdout abstraction, test seams, interface design |
| `07-state-context` | `studies/go-cli-study/reports/final/07-state-context.md` | Context propagation, app state, cancellation |
| `08-concurrency` | `studies/go-cli-study/reports/final/08-concurrency.md` | Worker pools, cancellation, parallel execution, goroutine management |
| `10-logging-observability` | `studies/go-cli-study/reports/final/10-logging-observability.md` | Logs, diagnostics, structured events, observability patterns |
| `11-testing-strategy` | `studies/go-cli-study/reports/final/11-testing-strategy.md` | Unit, integration, fixture, command-level tests, coverage strategy |
| `12-extensibility` | `studies/go-cli-study/reports/final/12-extensibility.md` | Extension points, plugin architecture, package boundaries, API design |
| `13-security` | `studies/go-cli-study/reports/final/13-security.md` | Path safety, secrets, command injection risks, sandboxing |
| `14-performance` | `studies/go-cli-study/reports/final/14-performance.md` | Startup latency, large repos, memory management, performance |

## Selected Reasoning Templates

All paths must appear in the project index's "Available Reasoning Templates" table.

| Template | Output Path | Why Selected |
| --- | --- | --- |
| Architecture | `projects/ultraplan-go/sprints/30-web-foundations/reasoning/architecture.md` | Reason through explicit CLI/TUI/web composition, `internal/web` ownership, dependency direction, app use-case boundaries, and product versus ephemeral HTTP state. |
| API Design | `projects/ultraplan-go/sprints/30-web-foundations/reasoning/api-design.md` | Reason through read-only HTML and `/api/v1` resources, methods, compatibility boundaries, structured errors, health data, and artifact-preview semantics. |
| Frontend | `projects/ultraplan-go/sprints/30-web-foundations/reasoning/frontend.md` | Reason through server-rendered browser flows, progressive enhancement, accessibility, state placement, hostile-content rendering, and the no-build-pipeline constraint. |

## Prior Decisions To Carry Forward

No prior decisions are cataloged in the project index's "Prior Decisions" table. Phase 3 readiness and the Sprint 24-26 shared app/TUI/review foundations remain dependencies through the sprint requirements, but no decision path is selectable here.

## Required Review Protocols

All paths must appear in the project index's "Review Protocols" table.

| Protocol | Path | Required Evidence |
| --- | --- | --- |
| Architecture Review | `system/protocols/architecture-review-protocol.md` | Evidence that composition is explicit and package-cycle-free, `internal/web` depends only on typed `internal/app` surfaces, product modules retain workflow/state ownership, and CLI/TUI behavior remains intact. |
| Sprint Review | `system/protocols/review-sprint-protocol.md` | Current `review.md` evidence covering selected contracts, handbook guidance, sprint decisions, plan execution, route/security/lifecycle/template/artifact tests, documentation, findings, and final verdict. |

## Excluded Context

| Context | Reason Excluded | Revisit If |
| --- | --- | --- |
| Implementation execution | The index selects planning and review context only; browser-triggered execute behavior and implementation execution itself are outside the read-only web foundation. | `plan.md` is complete and the governed execute stage begins, or a later sprint explicitly scopes browser execution operations. |
| Smoke investigation | Sprint 30 may display existing smoke status read-only, but it does not invoke the harness, investigate live failures, maintain harness evidence, or change smoke semantics. | A later sprint explicitly requires web smoke operations or a selected deep-smoke protocol requires real-system evidence. |
| Review automation | Sprint 30 may display existing review status read-only, but it does not run reviewers, alter review orchestration, or change review verdict semantics. | A later sprint explicitly brings guarded browser review operations into scope. |
| Issue tracking | General-purpose issue creation, assignment, scheduling, synchronization, and project-management behavior are explicit non-goals. | A later project phase explicitly introduces issue-management requirements and catalog context. |
| Git mutation | Automatic or browser-triggered add, commit, push, branch, merge, reset, checkout, or other Git mutation remains prohibited. | A future sprint explicitly scopes governed Git operations under newly selected requirements and contracts. |
| Guarded web operations and SSE | Confirmation tokens, operation start/status, explicit cancellation, progress streaming, reconnect behavior, and product mutation locking belong to Sprint 31 rather than the read-only foundation. | Sprint 31 begins after the read-only routes, app queries, lifecycle, and security foundation are accepted. |
| Hosted and remote web service behavior | The server is loopback-only; hosted SaaS, LAN/public binding, accounts, teams, permissions, tenants, collaboration, remote workers, and remote synchronization are excluded. | A later product phase explicitly defines and governs a remote or multi-user security model. |
| Browser-owned state and arbitrary editing | Workspace artifacts and product-owned run state remain authoritative; the browser cannot edit arbitrary files or introduce database-backed or web-specific durable product state. | Product requirements explicitly add a governed editing or alternate-persistence capability. |
| Frontend framework or build pipeline | The initial UI uses embedded Go templates, CSS, and minimal JavaScript and must not require Node.js, Vite, a separate frontend process, or an asset build step. | Demonstrated client-side complexity leads a later sprint to select and govern a richer frontend architecture. |
| WebSockets and terminal/session transport | Bidirectional streams, agent chat, terminal transport, and remote process control are unrelated to the read-only dashboard and are explicit non-goals. | A later phase explicitly brings interactive bidirectional transport into scope with security and lifecycle requirements. |

## Next Artifacts

- `technical-handbook.md` reads from the evidence reports listed above.
- `reasoning/*.md` captures area-specific reasoning.
- `reasoning.md` makes final sprint decisions.
- `plan.md` executes `reasoning.md`.
- `review.md` runs the selected review protocols against implementation.
