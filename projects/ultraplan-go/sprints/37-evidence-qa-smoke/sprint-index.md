# Sprint Index: Evidence-producing QA and smoke integration

> Project: `ultraplan-go`
> Sprint: `37-evidence-qa-smoke`
> Purpose: selected context for this sprint. Must be a subset of `projects/ultraplan-go/project-index.md`.
> **Inputs Used:** `projects/ultraplan-go/project-index.md`, `projects/ultraplan-go/sprints/37-evidence-qa-smoke/requirements.md`, `projects/ultraplan-go/roadmap.md`, `projects/ultraplan-go/docs/ARCHITECTURE.md`, `projects/ultraplan-go/docs/PRD.md`, `projects/ultraplan-go/docs/TRD.md`

This document selects what must be read, distilled, reasoned through, or checked for this sprint. It does not make implementation decisions. All selections come from the project index.

## Sprint Scope

- **Sprint Goal:** Add isolated evidence-producing QA with product-owned adjudication, canonical QA reporting, and smoke execution through QA while preserving target immutability and every existing smoke safety and compatibility guarantee.
- **Planned Output:** Product-neutral isolation and process mechanics; sprint-owned writable investigation, evidence, adjudication, bounded issue, assessment, and versioned verification state; canonical `qa.md`; smoke-as-QA compatibility; shared CLI, JSON, TUI, and browser operations; deterministic, fault-injection, race, parity, and gated dogfood tests; and aligned command, architecture, workflow, browser, recovery, schema, and release documentation.
- **Depends On:** Sprint 36's current acceptable Conformance Review and required smoke evidence proving deterministic mapping, changed-path coverage, read-only investigation, cancellation, resume, fingerprint invalidation, synthesis, and target non-mutation; Sprint 35 durable run control; the delivered grounded-planning context; current governed inputs and implementation identity; the existing review and smoke implementation; Agentwrap and the platform process boundary; and the cataloged `ultraplan-go-smoke` harness and manifest.
- **Non-Goals:** Production repair or generated-patch application; permanent promotion of generated checks into product tests; replacement or removal of `review`, `review.md`, `smoke`, or `smoke.md`; QA as a planning stage; a general issue tracker, workflow engine, scheduler, broker, plugin system, or remote worker protocol; content identity, provenance, retrieval, alternate authored-artifact persistence, hosted service, cloud authority, or multi-user collaboration; and automatic Git mutation.

## Source Project Index

- `projects/ultraplan-go/project-index.md` - authoritative source. Every selected contract, report, reasoning template, and protocol below appears there.

## Selected Contracts

Each contract applies as a flat whole to this sprint. All paths appear in the project index's "Active Contract Pool" table.

| Contract | Why Selected |
| --- | --- |
| Architecture | Governs sprint ownership of QA policy and evidence semantics, product-neutral platform isolation, shared app use cases, adapter dependency direction, and reuse of the existing smoke path. |
| Errors | Governs fail-closed isolation, identity, containment, cleanup, validation, cancellation, stale-state, and recovery failures with actionable diagnostics. |
| Configuration | Governs finite QA limits, timeout and concurrency settings, environment forwarding, runtime selection, validation, and redaction. |
| Observability | Governs durable run correlation, bounded progress and evidence, cancellation and cleanup facts, diagnostics, replay, recovery, and cross-surface agreement. |
| Security | Governs default-deny investigator permissions, target immutability, path and symlink containment, process isolation, browser request guards, hostile content, and secret redaction. |
| Testing | Governs deterministic fixtures, fake runtimes and processes, fault injection, race coverage, smoke parity, adapter agreement, and gated real-repository evidence. |
| Documentation | Governs command, architecture, workflow, browser, recovery, schema, compatibility, and release documentation for evidence-producing QA. |
| CLI Surface | Governs `qa`, `qa --suite smoke`, focused shard, status, cancellation, resume, compatibility aliases, stable text and JSON output, help, and exit behavior. |
| LLM Runtime | Governs bounded investigator and adjudicator requests through Agentwrap, frozen context, structured output validation, permissions, cancellation, and runtime metadata. |
| LLM Evaluation / Cost / Safety | Governs untrusted model output, evidence quality checks, bounded cost and retries, safety policy, adjudication discipline, and truthful blocked outcomes. |
| Workflows | Governs isolated attempt orchestration, durable acceptance, fencing, cancellation, resume, recovery, retries, publication ordering, and smoke-suite delegation. |
| Performance | Governs finite workspace, attempt, command, evidence, output, storage, duration, retry, follow-up, rendering, and concurrency limits. |
| Persistence And Migrations | Governs versioned verification state, atomic publication, digest links, stale-writer fencing, compatibility, invalidation, recovery, and bounded `flow-state.json` projections. |

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
| Architecture | `projects/ultraplan-go/sprints/37-evidence-qa-smoke/reasoning/architecture.md` | Reason through isolation ownership, target and workspace identity, evidence authority, adjudication, verification-state publication, durable run reuse, smoke delegation, and package dependency direction. |
| API Design | `projects/ultraplan-go/sprints/37-evidence-qa-smoke/reasoning/api-design.md` | Reason through typed app results, QA and smoke-suite requests, versioned JSON, guarded starts, cancellation, durable progress, reconnect, compatibility aliases, and recovery errors. |
| Frontend | `projects/ultraplan-go/sprints/37-evidence-qa-smoke/reasoning/frontend.md` | Reason through bounded evidence and issue presentation, canonical assessment, smoke-suite status, hostile text, no-JavaScript behavior, progressive enhancement, accessibility, cancellation, and recovery. |

## Prior Decisions To Carry Forward

The project index catalogs no prior decisions. Sprint 35 durable run control, Sprint 36 read-only QA, grounded-planning context, and existing review and smoke compatibility remain explicit dependencies through the sprint requirements and project documents rather than uncataloged decision entries.

## Required Review Protocols

All paths appear in the project index's "Review Protocols" table.

| Protocol | Path | Required Evidence |
| --- | --- | --- |
| Architecture Review | `system/protocols/architecture-review-protocol.md` | Evidence that QA semantics stay in `internal/sprint`, generic isolation stays product-neutral, web depends only on app use cases, run control remains authoritative for execution lifecycle, detailed QA state stays outside `flow-state.json`, and smoke uses its existing authority. |
| Sprint Review | `system/protocols/review-sprint-protocol.md` | Current evidence for isolation, non-mutation, adjudication, issue bounds, canonical assessment, atomic state, cancellation, recovery, security, documentation, compatibility, cross-surface agreement, and required test, race, vet, build, and diff checks. |
| Deep Smoke Sprint | `system/protocols/deep-smoke-sprint-protocol.md` | Gated real-repository evidence that generated checks remain contained, target identity stays unchanged, workspaces clean up, adjudication records a rejection or promotion audit, and `smoke` matches `qa --suite smoke` in authority and guarantees. |

## Excluded Context

| Context | Reason Excluded | Revisit If |
| --- | --- | --- |
| Implementation execution | QA may execute bounded evidence plans, but it cannot implement or repair production changes, apply generated patches, or mutate the target checkout. | Sprint 38 authorizes one frozen adjudicated issue packet and explicit repair scope. |
| Independent smoke investigation | Smoke execution is in scope only through the existing manifest-driven smoke authority, wrapped as the `smoke` QA suite. A second discovery, selection, invocation, evidence, or verdict path is excluded. | Parity evidence exposes a separately governed defect that cannot be fixed through the existing smoke path. |
| New review automation | The current Conformance Review is an independent read-only assessment input. QA cannot replace it, change its verdict, or create a second reviewer authority. | A later requirement explicitly changes the Conformance Review capability. |
| General-purpose issue tracking | Bounded adjudicated QA issue records are in scope, but assignment, scheduling, mutable workflows, remote synchronization, and project-management behavior are not. | A later roadmap gate explicitly promotes issue-management behavior. |
| Git mutation | QA cannot add, commit, push, branch, merge, reset, clean, or otherwise mutate Git state, and it cannot assume the target is clean or Git-backed. | A later explicit requirement and selected contracts authorize a specific Git workflow. |
| Production repair and convergence | Repair packets, manual repair, automatic repair, production mutation, progressive reverification, and convergence decisions belong to Sprint 38. | Sprint 37 passes isolation, evidence, adjudication, issue, and smoke-parity gates. |
| Permanent test promotion | Generated tests, fixtures, probes, and smoke scenarios remain evidence or regression candidates and are not installed into production or harness test suites by QA. | A later governed repair or maintenance sprint explicitly accepts a regression candidate. |
| Alternate smoke evidence authority | Raw smoke JSON, stdout, stderr, per-test artifacts, and harness issue files remain under manifest-declared external harness roots. UltraPlan stores validated links and bounded summaries only. | A later compatibility and authority decision explicitly migrates the harness contract. |
| Alternate product persistence, content identity, retrieval, graph, cloud, and collaboration | Product Phase 5 keeps current artifact authority and defers these systems until post-Sprint-39 evidence gates. | Post-Sprint-39 evidence justifies a separately governed capability. |
| Generic sandbox, workflow, or worker framework | The sprint needs reusable local isolation mechanics, not a plugin system, scheduler, broker, daemon, remote worker protocol, or general workflow engine. | Repeated concrete product needs justify a separately scoped shared capability. |

## Next Artifacts

- `technical-handbook.md` reads from the evidence reports listed above.
- `reasoning/*.md` captures area-specific reasoning.
- `reasoning.md` makes final sprint decisions.
- `plan.md` executes `reasoning.md`.
- `review.md` runs the selected review protocols against implementation.
