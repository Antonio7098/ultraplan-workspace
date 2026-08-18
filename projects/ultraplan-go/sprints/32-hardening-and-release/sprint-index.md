# Sprint Index: Local Web Hardening and Observable-Product Release

> Project: `ultraplan-go`
> Sprint: `32-hardening-and-release`
> Purpose: selected context for this sprint. Must be a subset of `projects/ultraplan-go/project-index.md`.
> **Inputs Used:** `projects/ultraplan-go/project-index.md`, `projects/ultraplan-go/roadmap.md`, `projects/ultraplan-go/sprints/32-hardening-and-release/requirements.md`, `projects/ultraplan-go/docs/ARCHITECTURE.md`, `projects/ultraplan-go/docs/PRD.md`, `projects/ultraplan-go/docs/TRD.md`

This document selects what must be read, distilled, reasoned through, or checked for this sprint. It does not make implementation decisions. All selections must come from the project index; no items may be included that are not listed in the project index.

## Sprint Scope

- **Sprint Goal:** Turn the Sprint 30-31 local browser implementation into a secure, accessible, compatibility-controlled, recoverable, and release-ready interface while proving that new workflow stages can use the shared application capability without route-specific workflow logic.
- **Planned Output:** Stable `/api/v1` routes and typed envelopes; hardened HTTP, operation, security, SSE, recovery, and server-lifecycle behavior; reusable web application and operation capabilities; layered templates, CSS, and dependency-free progressive enhancement; compatibility, security, lifecycle, concurrency, accessibility, integration, and capability tests; complete local-web documentation; a release checklist; and gated real-runtime and smoke-harness browser evidence.
- **Depends On:** Sprint 30's implemented local-web foundation and review findings; Sprint 31's planning artifacts and eight completed execute tasks, including guarded operations, SSE, mutation exclusion, cancellation, and shutdown behavior; Sprint 28's applicable review-to-smoke findings; and the Phase 4 requirements in the project Architecture, PRD, TRD, and roadmap.
- **Non-Goals:** New planning, execution, review, smoke, QA, repair, retrieval, or `code-context` stages; hosted, LAN/public, remote, multi-user, account, authorization, tenancy, collaboration, or remote-worker behavior; a client router, frontend framework/build pipeline, authoritative browser store, service worker, WebSockets, or agent chat; a database, web-owned durable state, detached worker, operation queue, or second scheduler; arbitrary browser file editing or automatic repair; issue tracking; Git mutation; and Product Phase 5 content identity, provenance, retrieval, SQLite, knowledge graph, cloud, or Aren capabilities.

## Source Project Index

- `projects/ultraplan-go/project-index.md` is the authoritative source. Any file or item referenced below must appear there.

## Selected Contracts

Each contract applies as a flat whole to this sprint. All paths must appear in the project index's "Active Contract Pool" table.

| Contract | Why Selected |
| --- | --- |
| Architecture | Governs `internal/web -> internal/app` dependency direction, explicit interface composition, product ownership of workflow and durable recovery rules, presentation hierarchy, and the prohibition on route-specific workflow logic or alternate durable state. |
| CLI Surface | Applies to release stability and documentation for `ultraplan serve`, its loopback options, lifecycle, diagnostics, build behavior, and agreement with browser and TUI projections. |
| Configuration | Applies to documented local-web settings and precedence, loopback restrictions, bounded operation and stream limits, safe diagnostics, and secret redaction. |
| Documentation | Applies to the local-web, configuration, CLI, user, recovery, architecture, troubleshooting, accessibility, packaging, and release-checklist outputs and their agreement with implemented routes, fields, limits, states, and commands. |
| Errors | Applies to compatibility-controlled typed API errors, methods and status mappings, safe HTML/JSON/SSE projection, lifecycle and recovery diagnostics, and stable machine-readable error codes. |
| LLM Evaluation / Cost / Safety | Applies to safe gated runtime-backed browser operations, bounded evaluation behavior, runtime/model metadata, confirmation scope, and protection from prompts, provider payloads, stderr, and secrets. |
| LLM Runtime | Applies because representative browser planning and study operations must preserve shared runtime boundaries, progress, cancellation, process cleanup, safe results, and durable reconciliation. |
| Observability | Applies to truthful operation status, bounded safe events, lifecycle and shutdown visibility, cross-surface state agreement, recovery guidance, redacted diagnostics, and release evidence. |
| Performance | Applies to bounded operations, preparations, events, payloads, subscribers, streams, retention, heartbeat, lifetime, polling, request bodies, cleanup waits, and slow-client behavior. |
| Persistence And Migrations | Applies because product-owned workspace and run state remain authoritative, terminal and cleanup-uncertain outcomes must be durable and atomic, restart reconciliation must be conservative, and the web layer must not add persistent state. |
| Security | Applies to loopback, Host/Origin, CSRF, session, CSP and response headers, body limits, hostile input and Markdown, path containment, forged references, secret redaction, subprocess safety, and release audit requirements. |
| Testing | Applies to deterministic compatibility fixtures, `httptest`, fake runtime and harness integration, accessibility assertions, race/leak/concurrency coverage, packaging checks, cross-surface agreement, and gated real-system evidence. |
| Workflows | Applies to guarded preparation and execution, mutation exclusion, progress, cancellation, exact-once shutdown cleanup, durable reconciliation, shared operation capabilities, and recovery without a web-owned workflow engine. |

## Selected Evidence Reports

Copied from the project index's "Available Evidence Reports" table. These tell the technical handbook which reports to read - the project index is the authoritative source.

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
| Architecture | `projects/ultraplan-go/sprints/32-hardening-and-release/reasoning/architecture.md` | Reason through web-to-app boundaries, shared stage capabilities, durable ownership, operation and shutdown lifecycle, presentation layering, recovery, and release packaging without route-specific workflow logic. |
| API Design | `projects/ultraplan-go/sprints/32-hardening-and-release/reasoning/api-design.md` | Reason through compatibility-controlled `/api/v1` routes, methods, success and error envelopes, status codes, cache behavior, confirmation, cancellation, SSE replay and terminal semantics, and safe recovery projections. |
| Frontend | `projects/ultraplan-go/sprints/32-hardening-and-release/reasoning/frontend.md` | Reason through namespaced template composition, typed view models, layered CSS, no-JavaScript completeness, progressive enhancement, accessibility, responsive behavior, operation recovery, and embedded assets. |

## Prior Decisions To Carry Forward

No prior decisions are cataloged in the project index's "Prior Decisions" table. Sprint 30's accepted web-to-app and read-only foundation, Sprint 31's guarded-operation and SSE design plus completed execute work, and Sprint 28's still-applicable verification findings remain explicit dependencies through the sprint requirements, but no decision path is selectable here.

## Required Review Protocols

All paths must appear in the project index's "Review Protocols" table.

| Protocol | Path | Required Evidence |
| --- | --- | --- |
| Architecture Review | `system/protocols/architecture-review-protocol.md` | Evidence that `internal/web` remains an interface adapter over typed app capabilities, workflow and durable recovery remain product-owned, templates compose downward from pages to primitives through explicit view models, operation state stays bounded and ephemeral, and no route-specific workflow logic or alternate persistence is introduced. |
| Sprint Review | `system/protocols/review-sprint-protocol.md` | Current `review.md` evidence covering every selected contract and handbook concern, API compatibility, cross-surface agreement, security and redaction, accessibility, concurrency and leaks, shutdown and recovery, packaging, documentation, release commands, findings, and a verdict with no unresolved applicable high-severity local-web issue. |
| Deep Smoke Sprint | `system/protocols/deep-smoke-sprint-protocol.md` | Gated evidence for one real runtime-backed browser operation and one real smoke-harness-backed browser operation, with linked harness evidence when prerequisites exist and an explicit blocked result rather than a pass when they do not. |

## Excluded Context

| Context | Reason Excluded | Revisit If |
| --- | --- | --- |
| Implementation execution during index creation | This artifact selects context and makes no implementation decisions or source changes; governed implementation belongs to the later reasoning, plan, and execute stages. | The planning chain reaches a validated `plan.md` and the governed execute stage begins. |
| Smoke investigation and harness maintenance | Gated real smoke-harness evidence is required, but live failure investigation, authored harness-suite maintenance, issue-record mutation, and new smoke semantics are not Sprint 32 product scope. | The selected Deep Smoke Sprint protocol exposes a blocker requiring separately governed harness work or a later sprint changes smoke behavior. |
| Review automation changes | Sprint review is required as release evidence, but reviewer orchestration, contract resolution, deterministic verdict rules, and `review.md` semantics are not being redesigned. | A later sprint explicitly changes the review stage rather than reviewing this implementation. |
| Issue tracking | General-purpose issue creation, assignment, scheduling, synchronization, and project-management behavior remain explicit non-goals; existing review and smoke findings remain evidence only. | A later product phase explicitly introduces issue-management requirements and catalog context. |
| Git mutation | Automatic add, commit, push, branch, merge, reset, checkout, or other Git-state mutation remains prohibited, and release verification must preserve unrelated working-tree changes. | A future sprint explicitly scopes governed Git operations under newly selected requirements and contracts. |
| New workflow stages | `code-context` and any new planning, execution, review, smoke, QA, repair, or retrieval stage are outside the web hardening release; this sprint proves only that shared app capabilities can expose a future stage. | Sprint 33 begins or another later requirement explicitly adds a stage. |
| Hosted and remote web behavior | The server remains numeric-loopback-only and same-origin; hosted SaaS, LAN/public binding, remote access, accounts, authorization, tenants, teams, collaboration, and remote workers are excluded. | A later product phase defines and governs a remote or multi-user security and authority model. |
| Browser-owned state and arbitrary editing | Workspace artifacts and product-owned run state remain authoritative; the browser cannot edit arbitrary files, own a workflow store, add a database, queue detached work, or become a second scheduler. | Product requirements explicitly add governed editing, durable workers, or an alternate persistence authority. |
| Rich client and bidirectional transport | Client routers, frontend frameworks, package managers, build pipelines, service workers, authoritative client stores, WebSockets, agent chat, terminal transport, and remote shell control are non-goals. | Demonstrated client-side or bidirectional interaction requirements lead a later sprint to select and govern them. |
| Product Phase 5 and later architecture | Content identity, provenance, retrieval, SQLite, knowledge graphs, cloud, Aren, empirical QA, and automatic repair are beyond Sprint 32 and its observable-product release gate. | A later roadmap sprint explicitly selects one of these evidence-gated capabilities. |

## Next Artifacts

- `technical-handbook.md` reads from the evidence reports listed above.
- `reasoning/*.md` captures area-specific reasoning.
- `reasoning.md` makes final sprint decisions.
- `plan.md` executes `reasoning.md`.
- `review.md` runs the selected review protocols against implementation.
