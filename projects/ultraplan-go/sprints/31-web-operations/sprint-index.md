# Sprint Index: Guarded Web Operations and SSE Progress

> Project: `ultraplan-go`
> Sprint: `31-web-operations`
> Purpose: selected context for this sprint. Must be a subset of `projects/ultraplan-go/project-index.md`.
> **Inputs Used:** `projects/ultraplan-go/project-index.md`, `projects/ultraplan-go/sprints/31-web-operations/requirements.md`, `projects/ultraplan-go/docs/ARCHITECTURE.md`, `projects/ultraplan-go/docs/PRD.md`, `projects/ultraplan-go/docs/TRD.md`

This document selects what must be read, distilled, reasoned through, or checked for this sprint. It does not make implementation decisions. All selections must come from the project index; no items may be included that are not listed in the project index.

## Sprint Scope

- **Sprint Goal:** Expose existing guarded validation, planning, execute, review, smoke, verify, study run-loop, and cancellation operations through the local web surface with server-issued confirmations and bounded SSE progress, without adding a second workflow engine or durable web-owned state.
- **Planned Output:** Typed app operation use cases; product-owned sprint mutation locks; a bounded ephemeral web operation hub; confirmed operation, status/result, SSE, and cancellation routes and handlers; browser security extensions; server-rendered operation views with minimal JavaScript; deterministic app, sprint-lock, route, operation, SSE, security, and template tests; and local-web, CLI-reference, and architecture documentation updates.
- **Depends On:** Sprint 30's `ultraplan serve`, loopback security, `/api/v1` envelopes, embedded browser assets, artifact previews, and explicit interface composition; Sprints 24-25 guarded TUI operation patterns; Sprint 23 execute state and cancellation; Sprints 26-29 review, smoke, and verify use cases; existing study run-loop locking and cancellation; and the Phase 4 rules in the project Architecture, PRD, and TRD.
- **Non-Goals:** New product workflows or browser-specific workflow semantics; arbitrary browser editing; hosted, remote, LAN/public, multi-user, account, permission, tenant, collaboration, or remote-worker behavior; WebSockets, terminal transport, live agent chat, or remote shell control; durable web state, a database, or a web job queue; new execute, review, smoke, planning, or study run-loop semantics; general-purpose issue tracking; automatic fixes or Git mutation; a frontend framework or build pipeline; and Sprint 32 browser compatibility, accessibility-release, packaging, and gated real-runtime hardening.

## Source Project Index

- `projects/ultraplan-go/project-index.md` is the authoritative source. Any file or item referenced below must appear there.

## Selected Contracts

Each contract applies as a flat whole to this sprint. All paths must appear in the project index's "Active Contract Pool" table.

| Contract | Why Selected |
| --- | --- |
| Architecture | Governs `internal/web -> internal/app` dependency direction, product ownership of workflow semantics and mutation locks, explicit interface composition, and separation of ephemeral transport state from durable product state. |
| Errors | Applies to typed app errors, structured HTTP error envelopes, conflict and cancellation classification, safe result projection, and actionable diagnostics without leaking provider payloads, stderr, secrets, or absolute paths. |
| Configuration | Applies to runtime/model and harness summaries, operation preflight, bounded server settings, same-origin session protection, config precedence, and redaction. |
| Observability | Applies to truthful typed progress, operation identifiers, lifecycle and cancellation visibility, safe diagnostics, monotonic event ordering, and recovery guidance based on durable state. |
| Security | Applies to confirmation binding and expiry, CSRF/session protection, Host/Origin checks, path containment, request and stream bounds, redaction, subprocess safety, and the loopback-only trust boundary. |
| Testing | Applies to deterministic app, lock, HTTP, SSE, security, template, cancellation, shutdown, race, and cross-surface agreement coverage using fakes and temporary workspaces. |
| Documentation | Applies to local-web operation guidance, CLI/API affordance documentation, architecture updates, trust limits, cancellation, reconnect, and durable-state recovery behavior. |
| LLM Runtime | Applies because browser-exposed runtime-backed planning, execute, review, and study operations must preserve the existing generic runtime boundary, context cancellation, safe event projection, and provider/model separation. |
| LLM Evaluation / Cost / Safety | Applies to confirmation-time runtime/model summaries, safe runtime-backed operation exposure, bounded evaluation behavior, metadata handling, and protection from unsafe payload disclosure. |
| Workflows | Applies to existing planning, execute, review, smoke, verify, and study run-loop orchestration, resumability, cancellation, mutation conflicts, and the prohibition on a second web-owned workflow engine. |
| Performance | Applies to bounded active operations, event buffers, subscribers, retention, request bodies, streams, heartbeats, and non-blocking slow or disconnected subscribers. |
| Persistence And Migrations | Applies because workspace artifacts and product run state remain authoritative, product locks and state writes must remain safe, and the ephemeral web hub must not introduce a durable state format. |

## Selected Evidence Reports

Copied from the project index's "Available Evidence Reports" table. These tell the technical handbook which reports to read – the project index is the authoritative source.

| Report | Path | Covers |
| --- | --- | --- |
| `01-project-structure` | `studies/go-cli-study/reports/final/01-project-structure.md` | Project layout, cmd/internal/pkg, dependency direction, thin entrypoints |
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
| Architecture | `projects/ultraplan-go/sprints/31-web-operations/reasoning/architecture.md` | Reason through operation capabilities, app/web boundaries, product-owned mutation locks, cancellation ownership, hub ephemerality, and recovery through durable product state. |
| API Design | `projects/ultraplan-go/sprints/31-web-operations/reasoning/api-design.md` | Reason through prepare/start/status/result/cancel resources, confirmation binding, SSE event semantics, error envelopes, methods, reconnect behavior, and compatibility constraints. |
| Frontend | `projects/ultraplan-go/sprints/31-web-operations/reasoning/frontend.md` | Reason through server-rendered confirmation, progress, result, finding, failure, cancellation, stale/reconnect, and recovery views with minimal progressive enhancement. |

## Prior Decisions To Carry Forward

No prior decisions are cataloged in the project index's "Prior Decisions" table. Sprint 30's accepted local-web foundation and the existing execute, TUI, review, smoke, verify, and study operation behavior remain dependencies through the sprint requirements, but no decision path is selectable here.

## Required Review Protocols

All paths must appear in the project index's "Review Protocols" table.

| Protocol | Path | Required Evidence |
| --- | --- | --- |
| Architecture Review | `system/protocols/architecture-review-protocol.md` | Evidence that `internal/web` remains transport-only, handlers use typed app operations, product modules own locks and workflow state, operation-hub state is bounded and ephemeral, cancellation ownership is correct, and no package cycle or alternate workflow engine is introduced. |
| Sprint Review | `system/protocols/review-sprint-protocol.md` | Current `review.md` evidence covering selected contracts, handbook guidance, sprint decisions, plan execution, confirmation and API behavior, SSE bounds and reconnect, cancellation, mutation conflicts, security and redaction, cross-surface compatibility, documentation, verification results, findings, and final verdict. |

## Excluded Context

| Context | Reason Excluded | Revisit If |
| --- | --- | --- |
| New implementation execution semantics | Sprint 31 exposes the existing execute use case through a guarded browser operation but does not change plan-task extraction, runtime execution, durable execute state, retries, completion rules, or implementation semantics. | A later sprint explicitly changes the execute workflow rather than adding another interface over it. |
| Smoke investigation | The browser may start and observe the existing smoke use case, but live failure investigation, harness test/issue maintenance, new smoke selection semantics, and new external evidence behavior are outside this sprint; normal proof uses fake harnesses. | A selected Deep Smoke Sprint protocol or later hardening sprint requires gated real-harness investigation or changes smoke semantics. |
| Review automation changes | The browser may start and observe the existing automated review use case, but reviewer orchestration, contract resolution, deterministic checks, verdict rules, and `review.md` semantics are unchanged. | A later sprint explicitly changes review automation rather than adding another interface over it. |
| Issue tracking | General-purpose issue creation, assignment, scheduling, synchronization, and project-management behavior remain explicit non-goals; existing smoke issue evidence is only presented through current app results where applicable. | A later product phase explicitly introduces issue-management requirements and catalog context. |
| Git mutation | Automatic or browser-triggered add, commit, push, branch, merge, reset, checkout, or other Git mutation remains prohibited. | A future sprint explicitly scopes governed Git operations under newly selected requirements and contracts. |
| New product workflows or browser-owned workflow state | The browser may expose only current app use cases; it must not add stages, duplicate validators or verdicts, persist a durable job queue, or become a second workflow engine. | Product requirements explicitly introduce a new product workflow or persistence authority. |
| Hosted and remote web behavior | The server remains loopback-only and same-origin; hosted SaaS, LAN/public binding, accounts, teams, permissions, tenants, collaboration, remote workers, and remote synchronization are excluded. | A later product phase defines and governs a remote or multi-user security model. |
| Arbitrary browser editing | Browser editing of workspace files, generated artifacts, project docs, plans, review findings, and smoke issue records is not part of guarded operation exposure. | A later sprint defines allowlisted editing, validation, conflict, and persistence behavior. |
| WebSockets and interactive process transport | SSE is one-way progress only; bidirectional terminal/session transport, live agent chat, and browser-driven remote shell or process control are explicit non-goals. | A later phase explicitly brings bidirectional transport into scope with security and lifecycle requirements. |
| Frontend framework and asset pipeline | The UI remains Go `html/template` with embedded CSS and minimal dependency-free JavaScript, without Node.js, Vite, a separate frontend process, service workers, or a build pipeline. | Demonstrated client-side complexity leads a later sprint to select and govern a richer frontend architecture. |
| Sprint 32 release hardening | Browser-level compatibility hardening, accessibility release audit, gated real-runtime browser evidence, recovery hardening beyond Sprint 31 acceptance, and release packaging are deferred. | Sprint 32 begins or one item becomes necessary to prove Sprint 31 behavior with deterministic fakes. |

## Next Artifacts

- `technical-handbook.md` reads from the evidence reports listed above.
- `reasoning/*.md` captures area-specific reasoning.
- `reasoning.md` makes final sprint decisions.
- `plan.md` executes `reasoning.md`.
- `review.md` runs the selected review protocols against implementation.
