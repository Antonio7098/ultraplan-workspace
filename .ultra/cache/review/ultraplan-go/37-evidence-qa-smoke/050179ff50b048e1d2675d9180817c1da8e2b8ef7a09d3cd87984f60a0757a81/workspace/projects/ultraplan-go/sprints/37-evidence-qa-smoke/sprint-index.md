# Sprint Index: Evidence-Producing QA and Smoke Integration

> Project: `ultraplan-go`
> Sprint: `37-evidence-qa-smoke`
> Purpose: selected context for this sprint. Must be a subset of `projects/ultraplan-go/project-index.md`.
> **Inputs Used:** `projects/ultraplan-go/project-index.md`, `projects/ultraplan-go/sprints/37-evidence-qa-smoke/requirements.md`, `projects/ultraplan-go/roadmap.md`, `projects/ultraplan-go/docs/ARCHITECTURE.md`, `projects/ultraplan-go/docs/PRD.md`, `projects/ultraplan-go/docs/TRD.md`

This document selects what must be read, distilled, reasoned through, or checked for this sprint. It does not make implementation decisions. All selections come from the project index.

## Sprint Scope

- **Sprint Goal:** Add isolated, evidence-producing QA with product-owned adjudication, canonical QA reporting, and smoke execution through QA while preserving target immutability and every existing smoke safety and compatibility guarantee.
- **Planned Output:** Validated per-attempt isolation, frozen evidence plans, durable generated evidence, global adjudication, bounded promoted issues, canonical `qa.md`, smoke-as-QA parity, and consistent CLI, JSON, TUI, browser, durable-run, documentation, and recovery behavior.
- **Depends On:** Sprint 36 current acceptable Conformance Review and required smoke evidence; Sprint 35 durable run control; the existing review and smoke implementations; the cataloged `ultraplan-go-smoke` harness and manifest.
- **Non-Goals:** Production repair or patch application, permanent regression-test promotion, a general issue tracker, replacement of `review` or `smoke` compatibility, planning-stage expansion, content identity or retrieval, hosted or remote operation, automatic Git mutation, and alternate product persistence.

## Source Project Index

- `projects/ultraplan-go/project-index.md` — authoritative source. Every file or item referenced below appears there.

## Selected Contracts

Each contract applies as a flat whole to this sprint. All paths appear in the project index's "Active Contract Pool" table.

| Contract | Why Selected |
| --- | --- |
| Architecture | Governs sprint ownership of evidence semantics and adjudication, platform ownership of generic isolation mechanics, and adapter dependency direction. |
| CLI Surface | Governs `qa`, `qa --suite smoke`, compatibility commands, flags, help, exit behavior, and stable text and JSON output. |
| Configuration | Governs validated finite QA limits, runtime settings, environment allowlists, effective-source reporting, and redaction. |
| Documentation | Governs the required CLI, architecture, workflow, browser, recovery, schema, and release documentation. |
| Errors | Governs fail-closed error classification and actionable blocked, stale, invalid, cancelled, cleanup-uncertain, and recovery outcomes. |
| LLM Evaluation / Cost / Safety | Governs distrust of model claims, bounded investigator and adjudicator use, safety policy, usage records, and gated dogfood evaluation. |
| LLM Runtime | Governs Agentwrap/OpenCode requests, permissions, cancellation, structured output, and runtime identity without moving product decisions into the runtime adapter. |
| Observability | Governs durable run correlation, bounded progress, evidence diagnostics, cancellation truth, replay, and agreement across local interfaces. |
| Performance | Governs finite attempts, commands, output, storage, duration, retries, follow-up, and sequential writable execution until isolation is proven. |
| Persistence And Migrations | Governs schema versions, private state, digest links, atomic dependency-ordered publication, stale-writer fencing, compatibility, and recovery. |
| Security | Governs target immutability, path and symlink safety, process containment, default-deny permissions, browser protections, and secret-safe evidence. |
| Testing | Governs deterministic fixtures, fault injection, race coverage, smoke parity, cross-surface agreement, and gated real-repository checks. |
| Workflows | Governs bounded orchestration, cancellation, resume, recovery, durable ownership, terminal arbitration, and preservation of completed evidence. |

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
| Architecture | `projects/ultraplan-go/sprints/37-evidence-qa-smoke/reasoning/architecture.md` | Resolve product/platform ownership, isolation boundaries, state authority, publication order, smoke reuse, and compatibility without adding a general sandbox or workflow engine. |
| API Design | `projects/ultraplan-go/sprints/37-evidence-qa-smoke/reasoning/api-design.md` | Resolve additive app, CLI/JSON, browser request, durable-run, cancellation, recovery, and compatibility contracts for QA and smoke-suite operation. |
| Frontend | `projects/ultraplan-go/sprints/37-evidence-qa-smoke/reasoning/frontend.md` | Resolve bounded evidence, adjudication, issues, smoke status, blockers, cancellation, hostile content, accessibility, and no-JavaScript presentation across TUI and browser views. |

## Prior Decisions To Carry Forward

All decision paths must appear in the project index's "Prior Decisions" table. The project index currently lists no prior decisions.

| Decision | Path | Constraint For This Sprint |
| --- | --- | --- |
| None cataloged | Not applicable | Carry dependencies and compatibility requirements from the governed requirements and project documents, but do not invent a prior-decision artifact. |

## Required Review Protocols

All paths appear in the project index's "Review Protocols" table.

| Protocol | Path | Required Evidence |
| --- | --- | --- |
| Architecture Review | `system/protocols/architecture-review-protocol.md` | Package ownership, dependency direction, isolation and process seams, state authority, smoke reuse, and absence of parallel workflow or persistence authorities. |
| Sprint Review | `system/protocols/review-sprint-protocol.md` | Requirement coverage, target non-mutation, adjudication authority, schema and publication safety, compatibility, cross-surface agreement, tests, docs, and scope exclusions. |
| Deep Smoke Sprint | `system/protocols/deep-smoke-sprint-protocol.md` | Gated real-repository isolation evidence, adjudication audit, target identity preservation, cleanup proof, durable run correlation, and identical smoke authority through `smoke` and `qa --suite smoke`. |

## Excluded Context

| Context | Reason Excluded | Revisit If |
| --- | --- | --- |
| Implementation execution | This index selects Sprint 37 context and does not execute `plan.md` or implementation tasks. | The governed sprint reaches the execute stage after reasoning and planning complete. |
| Independent smoke investigation implementation | Smoke is in scope only as a QA executor that reuses the existing manifest-driven smoke path. A second discovery, selection, invocation, evidence, or verdict implementation is prohibited. | A later governed sprint explicitly changes smoke compatibility after parity and migration evidence. |
| Conformance Review automation changes | Sprint 37 consumes the current independent Conformance Review as an assessment input and must not rewrite its evidence, automation, verdict, or `review.md`. | A later sprint explicitly scopes changes to Conformance Review while preserving QA independence. |
| General-purpose issue tracking | Sprint 37 stores only bounded, adjudicated, evidence-backed QA issue records. Assignment, scheduling, remote synchronization, and project-management behavior remain out of scope. | A later product phase establishes a separate issue-management requirement and authority. |
| Git mutation | Investigators, checks, adapters, and QA orchestration must not add, commit, reset, clean, branch, merge, or otherwise mutate Git state. | A later explicitly governed publication or repair workflow owns a specific Git operation. |
| Production repair and generated-patch application | Generated patches remain evidence or regression candidates. Sprint 38 owns confirmed production repair and reverification. | Sprint 38 begins after Sprint 37 isolation, adjudication, and smoke-parity gates pass. |
| Permanent regression-test promotion | QA may classify a generated check as a regression candidate but cannot install it in the target repository. | A later governed repair or maintenance sprint explicitly promotes the candidate. |
| Content identity, retrieval, and knowledge graphs | Verification-scoped IDs and direct governed context are sufficient for this sprint. Global identity, provenance, retrieval, embeddings, and graph projections remain gated after QA and repair dogfood. | Post-Sprint-39 evidence justifies the content-contract gate. |
| Alternate product persistence | Detailed QA remains in sprint-owned versioned filesystem state, and Sprint 35 run-control persistence remains operational only. | A later authority decision proves a need for another authored-artifact persistence mode. |
| Hosted service and remote workers | QA remains local, same-host, and available through the existing CLI, TUI, and loopback browser boundaries. | A later product phase defines remote identity, authorization, authority, and worker protocols. |
| General sandbox, scheduler, broker, or workflow engine | Only reusable copy and process mechanics belong in the platform process package; QA policy and outcomes remain sprint-owned. | Multiple concrete non-QA consumers prove a smaller reusable capability is needed. |

## Next Artifacts

- `technical-handbook.md` reads from the evidence reports listed above.
- `reasoning/*.md` captures area-specific reasoning.
- `reasoning.md` makes final sprint decisions.
- `plan.md` executes `reasoning.md`.
- `review.md` runs the selected review protocols against implementation.
