# Sprint Index: Read-only QA decomposition and synthesis

> Project: `ultraplan-go`
> Sprint: `36-read-only-qa`
> Purpose: selected context for this sprint. Must be a subset of `projects/ultraplan-go/project-index.md`.
> **Inputs Used:** `projects/ultraplan-go/project-index.md`, `projects/ultraplan-go/sprints/36-read-only-qa/requirements.md`, `projects/ultraplan-go/roadmap.md`, `projects/ultraplan-go/docs/ARCHITECTURE.md`, `projects/ultraplan-go/docs/PRD.md`, `projects/ultraplan-go/docs/TRD.md`

This document selects what must be read, distilled, reasoned through, or checked for this sprint. It does not make implementation decisions. All selections come from the project index.

## Sprint Scope

- **Sprint Goal:** Add a deterministic, durable, cross-surface read-only QA phase that maps changed behavior into bounded verification surfaces, records resumable theory outcomes, and synthesizes cross-shard findings without generating checks, promoting issues, repairing code, or mutating the target repository.
- **Planned Output:** A verification-phase and QA domain model; separate schema-versioned QA state; deterministic mapping, read-only investigation, and bounded synthesis services; shared CLI, JSON, TUI, and browser use cases and presentation; tests, schemas, workflow documentation, recovery guidance, and release checks.
- **Depends On:** Sprint 35 durable run observability and its Conformance Review and smoke evidence; the delivered grounded-planning context; current execute and Conformance Review outputs; the existing Agentwrap runtime and shared application operation boundaries.
- **Non-Goals:** Generated checks or fixtures, smoke-as-QA integration, issue promotion, repair eligibility or execution, `qa.md`, changes to existing `review.md` or `smoke.md` compatibility, alternate artifact persistence, content identity, retrieval, cloud authority, hosted collaboration, and automatic Git mutation.

## Source Project Index

- `projects/ultraplan-go/project-index.md` - authoritative source. Every selected contract, report, reasoning template, and protocol below appears there.

## Selected Contracts

Each contract applies as a flat whole to this sprint. All paths appear in the project index's "Active Contract Pool" table.

| Contract | Why Selected |
| --- | --- |
| Architecture | Governs sprint ownership of QA semantics, product and platform separation, shared app use cases, and thin CLI, TUI, and web adapters. |
| Errors | Governs fail-closed schema, fingerprint, permission, cancellation, recovery, and command error behavior. |
| Configuration | Governs any QA runtime and budget settings, validation, model selection, environment forwarding, and redaction. |
| Observability | Governs durable run correlation, bounded progress, diagnostics, cancellation state, recovery, and cross-surface reporting. |
| Security | Governs read-only runtime permissions, path containment, command allowlisting, hostile content handling, redaction, and target identity checks. |
| Testing | Governs deterministic fixtures, fake runtimes, fault injection, race coverage, adapter parity, and gated real-runtime checks. |
| Documentation | Governs CLI, architecture, workflow, browser, recovery, schema, and release documentation for the QA phase. |
| CLI Surface | Governs command help, aliases, validation, text and JSON output, exit behavior, and compatibility for `review` and Conformance Review terminology. |
| LLM Runtime | Governs bounded investigator and synthesizer execution through Agentwrap without a competing runtime contract. |
| LLM Evaluation / Cost / Safety | Governs bounded investigation, explicit budgets, read-only safety, theory validation, and truthful blocked outcomes. |
| Workflows | Governs shard orchestration, durable acceptance, cancellation, resume, recovery, retries, and fingerprint invalidation. |
| Performance | Governs finite shard, context, command, output, wall-clock, concurrency, progress, and rendering limits. |
| Persistence And Migrations | Governs versioned QA schemas, atomic writes, compatibility, reconciliation, recovery, and separation from `flow-state.json` detail. |

## Selected Evidence Reports

Copied from the project index's "Available Evidence Reports" table. These tell the technical handbook which reports to read. The project index remains authoritative.

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

All paths appear in the project index's "Available Reasoning Templates" table.

| Template | Output Path | Why Selected |
| --- | --- | --- |
| Architecture | `projects/ultraplan-go/sprints/36-read-only-qa/reasoning/architecture.md` | Reason through verification-phase ownership, detailed-state authority, runtime boundaries, durable run reuse, and adapter dependency direction. |
| API Design | `projects/ultraplan-go/sprints/36-read-only-qa/reasoning/api-design.md` | Reason through typed app results, versioned HTTP requests, operation preparation, cancellation, reconnect, and JSON compatibility. |
| Frontend | `projects/ultraplan-go/sprints/36-read-only-qa/reasoning/frontend.md` | Reason through TUI and browser QA views, no-JavaScript state, progressive enhancement, bounded rendering, accessibility, and recovery. |

## Prior Decisions To Carry Forward

The project index catalogs no prior decisions. Sprint dependencies and compatibility constraints remain governed by the sprint requirements and project documents rather than an uncataloged decision entry.

## Required Review Protocols

All paths appear in the project index's "Review Protocols" table.

| Protocol | Path | Required Evidence |
| --- | --- | --- |
| Architecture Review | `system/protocols/architecture-review-protocol.md` | Verification phases remain independent of planning stages; sprint, app, runtime, run-control, TUI, and web ownership and import boundaries remain intact. |
| Sprint Review | `system/protocols/review-sprint-protocol.md` | The completed implementation satisfies the selected contracts, deterministic mapping, read-only enforcement, state durability, synthesis bounds, compatibility, and cross-surface parity. |
| Deep Smoke Sprint | `system/protocols/deep-smoke-sprint-protocol.md` | Gated real-runtime evidence covers map stability, non-mutation, durable progress, cancellation, restart recovery, and truthful blocking when prerequisites are unavailable. |

## Excluded Context

| Context | Reason Excluded | Revisit If |
| --- | --- | --- |
| Implementation execution by QA investigators | Sprint 36 investigators inspect and run approved existing non-mutating checks only. They do not implement fixes, tests, fixtures, or probes. | Sprint 37 defines isolated evidence production, or Sprint 38 supplies a frozen adjudicated repair scope. |
| Smoke investigation as QA | Existing smoke compatibility remains, but smoke-as-QA integration and evidence-producing investigation belong to Sprint 37. | Sprint 37 starts after read-only mapping, cancellation, resume, invalidation, and synthesis pass their gate. |
| New review automation | The existing `review` capability and `review.md` remain compatible and are labeled Conformance Review. QA consumes its findings but does not create a second review implementation or alter its verdict. | A later governed requirement changes the independent Conformance Review capability. |
| Issue tracking and promotion | QA retains theories and outcomes but does not create issue records, repair eligibility, or a general-purpose issue tracker. | Sprint 37 defines central evidence adjudication and frozen issue packets. |
| Git mutation | Automatic or investigator-driven Git mutation is prohibited. Before-and-after identity checks must expose any target drift. | A later explicit governance change authorizes a separately bounded Git workflow. |
| Generated tests, fixtures, probes, and smoke scenarios | This sprint is read-only and may only run approved existing checks. | Sprint 37 proves writable isolation and cleanup semantics. |
| Repair and repair loops | Production repair, manual repair, automatic repair, issue packets, and mutation authorized by theory outcomes are outside this sprint. | Sprint 38 begins with adjudicated issue packets and enforced repair scope. |
| Canonical `qa.md` and verdict promotion | Detailed QA state is verdict-neutral and cannot replace or upgrade Conformance Review. | A later sprint defines evidence adjudication and canonical QA reporting. |
| Alternate product persistence, content identity, retrieval, graph, cloud, and collaboration | Product Phase 5 keeps current artifact authority and defers these systems until later evidence gates. | Post-Sprint-39 evidence justifies a separately governed capability. |

## Next Artifacts

- `technical-handbook.md` reads from the evidence reports listed above.
- `reasoning/*.md` captures area-specific reasoning.
- `reasoning.md` makes final sprint decisions.
- `plan.md` executes `reasoning.md`.
- `review.md` runs the selected review protocols against implementation.
