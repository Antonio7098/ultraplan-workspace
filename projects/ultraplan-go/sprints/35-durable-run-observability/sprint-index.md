# Sprint Index: Durable Run Identity and Cross-Surface Observability

> Project: `ultraplan-go`
> Sprint: `35-durable-run-observability`
> Purpose: selected context for this sprint. Must be a subset of `projects/ultraplan-go/project-index.md`.
> **Inputs Used:** `projects/ultraplan-go/project-index.md`, `projects/ultraplan-go/sprints/35-durable-run-observability/requirements.md`, `projects/ultraplan-go/roadmap.md`, `projects/ultraplan-go/docs/ARCHITECTURE.md`, `projects/ultraplan-go/docs/PRD.md`, `projects/ultraplan-go/docs/TRD.md`

This document selects what must be read, distilled, reasoned through, or checked for this sprint. It does not make implementation decisions. All selections must come from the project index; no items may be included that are not listed in the project index.

## Sprint Scope

- **Sprint Goal:** Make every accepted runtime-backed UltraPlan execution a durable, workspace-visible run that can be discovered, inspected, followed, cancelled when authorized, and reconciled consistently from CLI, TUI, and any supported local web-server instance attached to the same workspace.
- **Planned Output:** A durable operational run identity and lifecycle model; durable sanitized event history and replay; workspace-wide active-run and detail projections; ownership, liveness, cancellation, terminal arbitration, and reconciliation behavior; compatibility and migration handling; correlated telemetry and support diagnostics; deterministic fault-injection, race, integration, browser, build, and gated real-runtime evidence; and aligned architecture, product, technical, CLI/API, local-web, operations, and recovery documentation.
- **Depends On:** Sprint 31 guarded operations, SSE, cancellation, locking, reconnect, and shutdown behavior; Sprint 32 browser hardening, API compatibility, recovery, and interface-state agreement fixtures; Sprint 34 shared-context real-runtime dogfood evidence; current implementation-repository operational plans; current agentwrap supervision and observability capabilities; and existing sprint, study, flow, execute, checkpoint, lock, cleanup, and smoke-evidence formats.
- **Non-Goals:** Hosted SaaS, public network exposure, team accounts, multi-user authorization, or an unproven remote-worker protocol; replacement of canonical Markdown, flow state, Git/source workspaces, or external smoke evidence; technology selection by preference alone; persistence of unredacted provider streams, full prompts, arbitrary output, or unlimited history; a general scheduler, workflow engine, distributed queue, or issue tracker; automatic repair or Git mutation; content identity, retrieval, knowledge-graph expansion, or alternate authored-product persistence; and browser disconnect as cancellation.

## Source Project Index

- `projects/ultraplan-go/project-index.md` is the authoritative source. Any file or item referenced below must appear there.

## Selected Contracts

Each contract applies as a flat whole to this sprint. All paths must appear in the project index's "Active Contract Pool" table.

| Contract | Why Selected |
| --- | --- |
| Architecture | Governs ownership of the shared execution-control boundary, package and dependency direction, separation from `internal/web`, product-module authority, agentwrap integration, and avoidance of a generic workflow engine. |
| CLI Surface | Governs workspace-wide run discovery, inspection, cancellation, diagnostics, text/JSON agreement, stable commands, exit behavior, and compatibility across local surfaces. |
| Configuration | Governs local-first topology and operational settings, precedence, validation, bounded limits, health behavior, safe diagnostics, and secret redaction. |
| Documentation | Governs the required architecture, PRD, TRD, CLI/API, local-web, operations, recovery, retention, compatibility, and diagnostic guidance. |
| Errors | Governs typed persistence, replay-gap, stale-owner, cancellation, migration, and recovery failures with preserved causes and actionable safe guidance. |
| LLM Evaluation / Cost / Safety | Governs safe runtime correlation, bounded retained metadata, cost and usage truthfulness, redaction, and exclusion of unsafe provider-native content. |
| LLM Runtime | Governs agentwrap as the runtime-supervision boundary, run/session/attempt correlation, canonical event consumption, cancellation propagation, cleanup, and durable observing integration. |
| Observability | Governs stable run, attempt, stage, task, runtime, owner, sequence, lifecycle, logging, metrics, health, diagnostics, and support-bundle correlation. |
| Performance | Governs bounded journals, retention, compaction, replay work, disk use, subscribers, backpressure, reconciliation, active-run queries, and multi-process behavior. |
| Persistence And Migrations | Governs fail-closed acceptance, atomic durable state, concurrent access, schema evolution, legacy operation links, tombstones, compatibility, recovery, and rollback. |
| Security | Governs loopback and session boundaries, separate read and mutation authorization, path safety, process identity, redaction-before-persistence, and safe local diagnostics. |
| Testing | Governs deterministic lifecycle and repository tests, race and fault injection, cross-process/server/browser integration, gated real-runtime dogfood, and build verification. |
| Workflows | Governs accepted-run lifecycle, attempts, cancellation, leases, reconciliation, terminal arbitration, resumability, and coexistence with product-owned locks and outcomes. |

## Selected Evidence Reports

Copied from the project index's "Available Evidence Reports" table. These tell the technical handbook which reports to read; the project index is the authoritative source.

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
| `09-terminal-ux` | `studies/go-cli-study/reports/final/09-terminal-ux.md` | Terminal output, progress indicators, human UX, color and formatting |
| `10-logging-observability` | `studies/go-cli-study/reports/final/10-logging-observability.md` | Logs, diagnostics, structured events, observability patterns |
| `11-testing-strategy` | `studies/go-cli-study/reports/final/11-testing-strategy.md` | Unit, integration, fixture, command-level tests, coverage strategy |
| `12-extensibility` | `studies/go-cli-study/reports/final/12-extensibility.md` | Extension points, plugin architecture, package boundaries, API design |
| `13-security` | `studies/go-cli-study/reports/final/13-security.md` | Path safety, secrets, command injection risks, sandboxing |
| `14-performance` | `studies/go-cli-study/reports/final/14-performance.md` | Startup latency, large repos, memory management, performance |

## Selected Reasoning Templates

All paths must appear in the project index's "Available Reasoning Templates" table.

| Template | Output Path | Why Selected |
| --- | --- | --- |
| Architecture | `projects/ultraplan-go/sprints/35-durable-run-observability/reasoning/architecture.md` | Reason through execution-control ownership, durable representation, coordinator and writer topology, authority boundaries, identity hierarchy, leases and fencing, reconciliation, cancellation routing, retention, and failure behavior. |
| API Design | `projects/ultraplan-go/sprints/35-durable-run-observability/reasoning/api-design.md` | Reason through durable run resources, operation compatibility, additive or versioned JSON changes, cursor replay and typed gaps, idempotent cancellation, authorization, tombstones, and recovery errors. |
| Frontend | `projects/ultraplan-go/sprints/35-durable-run-observability/reasoning/frontend.md` | Reason through workspace-wide running counts and lists, stable run detail, replay and reconnect states, liveness and uncertainty presentation, cancellation controls, recovery guidance, and CLI/TUI/browser agreement. |

## Prior Decisions To Carry Forward

No prior decisions are cataloged in the project index's "Prior Decisions" table. Sprint 31's guarded operation and SSE behavior, Sprint 32's compatibility and recovery constraints, and Sprint 34's cross-surface dogfood findings remain explicit dependencies through the sprint requirements, but no decision path is selectable here.

## Required Review Protocols

All paths must appear in the project index's "Review Protocols" table.

| Protocol | Path | Required Evidence |
| --- | --- | --- |
| Architecture Review | `system/protocols/architecture-review-protocol.md` | Evidence that durable run authority is outside interface adapters, product modules retain workflow and artifact authority, agentwrap retains runtime supervision, stale writers are fenced, lifecycle and terminal transitions have one authority, and the selected topology and store satisfy concurrency and recovery constraints. |
| Sprint Review | `system/protocols/review-sprint-protocol.md` | Current conformance evidence tracing acceptance through lease, attempts, ordered events, projections, cancellation, reconciliation, and terminal arbitration; proving cross-surface agreement, compatibility, redaction, bounded retention, persistence failures, race coverage, documentation, and required build/test results. |
| Deep Smoke Sprint | `system/protocols/deep-smoke-sprint-protocol.md` | Gated real-runtime evidence that CLI-started work is visible and inspectable across supported local surfaces and server instances, retained events replay after observer changes, cancellation and recovery remain truthful, and unavailable prerequisites produce a blocked result rather than a pass. |

## Excluded Context

| Context | Reason Excluded | Revisit If |
| --- | --- | --- |
| Implementation execution during index creation | This artifact selects planning context and makes no implementation or storage decisions; implementation belongs to the later reasoning, plan, and execute stages. | Validated reasoning and plan artifacts authorize governed execute work. |
| Smoke investigation during index creation | The sprint requires fault-injection and gated dogfood evidence, but live failure investigation and harness mutation are not part of selecting context. | A current review and authorized smoke stage require real-system evidence or expose a harness blocker. |
| Review automation changes | Sprint review is required evidence, but redesigning reviewer orchestration, verdict rules, or the review stage is outside the durable-run scope. | Durable-run integration proves that existing review observation cannot satisfy the requirements without a separately governed review change. |
| Issue tracking | General-purpose issue creation, assignment, scheduling, synchronization, and project-management behavior are explicit non-goals. | A later roadmap gate explicitly promotes issue-management behavior. |
| Git mutation | Automatic add, commit, push, branch, merge, reset, checkout, or other Git-state mutation is prohibited. | A later explicit requirement and selected contracts govern a specific Git operation. |
| Hosted, public, and multi-user operation | Hosted SaaS, public/LAN exposure, accounts, teams, tenants, and multi-user authorization are outside the local-first security boundary. | A later product phase defines a remote authority, identity, and security model. |
| General scheduling and remote-worker protocols | The sprint records and controls accepted UltraPlan executions; it does not build a workflow engine, distributed queue, broker, or remote-worker protocol by default. | Supported-topology reasoning proves a narrowly scoped coordination protocol is necessary for acceptance. |
| Alternate authored-product persistence | Durable operational run state must not replace Markdown, flow state, package-owned task state, Git/source state, or external smoke evidence. | A later persistence authority gate explicitly selects and migrates authored product storage. |
| Content identity, QA/repair expansion, retrieval, and knowledge graphs | These remain gated after durable-run dogfood and are not needed to repair execution identity and observation. | A later roadmap evidence gate promotes one of these capabilities. |
| Preselected infrastructure | SQLite, Postgres, a daemon, WebSockets, OpenTelemetry, brokers, filesystem journals, and coordinator topology are not selected by this index. | Area and final reasoning compare repository evidence and failure behavior and make the required decision. |
| Unsafe or unbounded event retention | Full prompts, raw provider payloads, credentials, unsafe paths, arbitrary stdout/stderr, and unlimited history are excluded from default persistence and streaming. | A separately authorized diagnostic mode defines explicit security, consent, retention, and redaction controls. |

## Next Artifacts

- `technical-handbook.md` reads from the evidence reports listed above.
- `reasoning/*.md` captures area-specific reasoning.
- `reasoning.md` makes final sprint decisions.
- `plan.md` executes `reasoning.md`.
- `review.md` runs the selected review protocols against implementation.
