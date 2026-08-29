# Sprint Index: Bounded Manual and Automatic Repair

> Project: `ultraplan-go`
> Sprint: `38-bounded-repair`
> Purpose: selected context for this sprint. Must be a subset of `projects/ultraplan-go/project-index.md`.
> **Inputs Used:** `projects/ultraplan-go/project-index.md`, `projects/ultraplan-go/sprints/38-bounded-repair/requirements.md`, `projects/ultraplan-go/roadmap.md`, `projects/ultraplan-go/docs/ARCHITECTURE.md`, `projects/ultraplan-go/docs/PRD.md`, `projects/ultraplan-go/docs/TRD.md`

This document selects what must be read, distilled, reasoned through, or checked for this sprint. It does not make implementation decisions. All selections come from the project index.

## Sprint Scope

- **Sprint Goal:** Add a governed repair phase that freezes one current adjudicated QA issue, requires explicit confirmation, permits one bounded production mutation, progressively reverifies the result, and exposes optional automatic repair only after the same manual protocol has current proof.
- **Planned Output:** Frozen repair packets and single-use confirmations; isolated bounded proposals and product-owned production mutation; progressive reverification and cleanup; strict private repair state and bounded flow summaries; lower-only automatic limits and closed outcomes; shared durable prepare, start, status, resume, cancel, and recovery operations; consistent CLI, JSON, TUI, and browser controls; documentation, tests, and real-runtime manual and automatic evidence gates.
- **Depends On:** Sprint 37 current evidence-producing QA, adjudicated repair-eligible issues, containing smoke parity, and real-runtime evidence; Sprint 36 deterministic QA maps and theory ownership; Sprint 35 durable run control and writer fencing; current isolation, process cleanup, target identity, mutation locking, shared app use cases, and server shutdown behavior. Automatic mode also depends on a current qualifying manual repair proof produced through the shared protocol.
- **Non-Goals:** Repair of unadjudicated or multiple issues, test or expectation weakening, requirement or evidence changes, speculative scope expansion, broad refactoring, dependency upgrades, migrations, product redesign, unbounded autonomy, general issue tracking, cross-project or remote repair, hosted collaboration, content identity or retrieval, alternate product persistence, and Git mutation.

## Source Project Index

- `projects/ultraplan-go/project-index.md` is the authoritative source. Every selected contract, report, reasoning template, and protocol below appears there.

## Selected Contracts

Each contract applies as a flat whole to this sprint. All paths appear in the project index's "Active Contract Pool" table.

| Contract | Why Selected |
| --- | --- |
| Architecture | Governs sprint ownership of repair semantics, platform ownership of generic isolation and process mechanics, shared app boundaries, and separation of Conformance Review, QA, repair, and interface adapters. |
| CLI Surface | Governs packet preview, explicit confirmation, manual and automatic commands, bounded queries, cancellation, recovery, stable text and JSON output, and exit behavior for every terminal outcome. |
| Configuration | Governs safe defaults, immutable maxima, lower-only workspace and environment limits, complete effective-source reporting, validation, and redaction. |
| Documentation | Governs architecture, CLI, user, schema, recovery, local-web, and release documentation for manual proof, automatic admission, bounded operation, and escalation. |
| Errors | Governs typed and actionable admission, stale-input, conflict, persistence, runtime, cancellation, cleanup, recovery, and terminal-outcome failures without false success. |
| LLM Evaluation / Cost / Safety | Governs bounded proposal generation, distrust of model claims, manual-first proof, cost and turn limits, safe retained evidence, and automatic stop conditions. |
| LLM Runtime | Governs Agentwrap/OpenCode proposal execution, permissions, cancellation, structured results, runtime identity, and cleanup without granting direct production write or verdict authority. |
| Observability | Governs durable run correlation, writer ownership, bounded progress and evidence, cancellation and shutdown truth, replay, diagnostics, and cross-surface agreement. |
| Performance | Governs finite cycles, commands, turns, output, patch size, changed files and bytes, wall-clock time, retention, cleanup, queries, and automatic stagnation checks. |
| Persistence And Migrations | Governs strict versioned repair records, private permissions, digest-bound immutable evidence, atomic publication, fencing, rollback, retention, compatibility, resume, and recovery. |
| Security | Governs explicit mutation authority, path containment, protected-file denial, symlink and hard-link defense, target identity checks, isolated proposals, secret redaction, browser guards, and fail-closed scope enforcement. |
| Testing | Governs deterministic, integration, race, failure-injection, cancellation, recovery, drift, cleanup, boundedness, hostile-input, parity, and gated real-runtime evidence. |
| Workflows | Governs durable acceptance, confirmation, single ownership, manual-first sequencing, bounded automatic cycles, progressive gates, cancellation, resume, recovery, cleanup, and exactly one terminal result. |

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
| Architecture | `projects/ultraplan-go/sprints/38-bounded-repair/reasoning/architecture.md` | Reason through repair authority, packet and state ownership, isolated proposal and production apply boundaries, writer fencing, progressive reverification, manual-proof publication, automatic reuse, cleanup, and terminal arbitration. |
| API Design | `projects/ultraplan-go/sprints/38-bounded-repair/reasoning/api-design.md` | Reason through prepare and confirmation binding, operation kinds, bounded packet and cycle queries, manual and automatic requests, cancellation, resume, recovery, stable outcomes, stale requests, and interface compatibility. |
| Frontend | `projects/ultraplan-go/sprints/38-bounded-repair/reasoning/frontend.md` | Reason through packet review, separate confirmation, manual-first gating, cycle and gate progress, scope and cleanup facts, hostile-text bounds, no-JavaScript operation, reconnect, cancellation, outcomes, blockers, and recovery guidance. |

## Prior Decisions To Carry Forward

All decision paths must appear in the project index's "Prior Decisions" table. The project index currently lists no prior decisions.

| Decision | Path | Constraint For This Sprint |
| --- | --- | --- |
| None cataloged | Not applicable | Carry dependencies, authority boundaries, and compatibility requirements from the governed requirements and project documents, but do not invent a prior-decision artifact. |

## Required Review Protocols

All paths appear in the project index's "Review Protocols" table.

| Protocol | Path | Required Evidence |
| --- | --- | --- |
| Architecture Review | `system/protocols/architecture-review-protocol.md` | Repair remains a sprint-owned verification phase; runtime and adapters cannot bypass packet, confirmation, isolation, scope, apply, reverification, cleanup, persistence, writer-fence, or terminal-result authority. |
| Sprint Review | `system/protocols/review-sprint-protocol.md` | Requirement and selected-contract coverage for manual-first delivery, protected paths, strict state, lower-only limits, all outcomes and stop conditions, cancellation and recovery, cross-surface parity, docs, tests, race checks, and build evidence. |
| Deep Smoke Sprint | `system/protocols/deep-smoke-sprint-protocol.md` | Real-runtime evidence for one manual repair through full reverification and cleanup on every interface, followed only with current proof by a bounded automatic stop path whose consumed limits survive resume or restart. |

## Excluded Context

| Context | Reason Excluded | Revisit If |
| --- | --- | --- |
| Implementation execution during index creation | This artifact selects repair context. It does not execute `plan.md`, mutate production, generate proposals, or choose implementation details. | Validated reasoning and plan artifacts authorize governed execute work. |
| Independent smoke investigation implementation | Repair may run the current containing smoke suite through progressive reverification, but it cannot add a second smoke discovery, execution, evidence, or verdict path. | A later governed sprint explicitly changes smoke authority after compatibility evidence. |
| Conformance Review automation replacement | Repair may request and record a focused delta through the independent current review capability. It cannot edit `review.md`, change review verdict rules, or infer global success. | A later sprint explicitly scopes review changes while preserving independent authority. |
| General-purpose issue tracking | Repair consumes exactly one current adjudicated QA issue. Assignment, scheduling, backlog management, remote synchronization, and repair of caller-supplied issues remain out of scope. | A later product phase defines separate issue-management authority and requirements. |
| Git mutation | Repair cannot add, commit, push, branch, merge, reset, checkout, stash, clean, edit hooks, or mutate the index. Product code applies only the confirmed production patch through its own non-Git boundary. | A later explicit requirement governs a specific Git operation. |
| Test, requirement, acceptance, evidence, or governed-input mutation | Repair must stop rather than weaken checks, rewrite expectations, update baselines, delete evidence, alter QA adjudication, or expand authority. | A newly adjudicated issue and separately governed sprint explicitly authorize a different asset class. |
| Multi-issue repair, broad refactoring, and product redesign | One packet authorizes one issue and a finite production path set. Scope growth, new issue classes, design decisions, upgrades, and migrations require escalation. | New adjudication and confirmation establish a separate bounded repair scope. |
| Unbounded or silently enabled automatic repair | Automatic mode is opt-in per run, requires current manual proof, reuses the manual protocol, and stops at lower-only product limits. | Later real dogfood supports a separately governed change to automation policy. |
| Content identity, retrieval, and knowledge graphs | Verification-scoped identifiers and digest-bound references remain sufficient. Global identity, provenance, retrieval, embeddings, and graph projections stay gated until QA and repair dogfood. | Post-Sprint-39 evidence justifies the content-contract gate. |
| Alternate product persistence | Repair details remain in strict sprint-owned filesystem records, while Sprint 35 storage remains authoritative only for operational run facts. | A later authority decision proves and governs another authored-artifact storage mode. |
| Hosted service, remote workers, and cross-project repair | Repair remains local, same-host, single-target, and available through the current CLI, TUI, and loopback browser boundaries. | A later product phase defines remote identity, authorization, ownership, cleanup, and artifact authority. |

## Next Artifacts

- `technical-handbook.md` reads from the evidence reports listed above.
- `reasoning/*.md` captures area-specific reasoning.
- `reasoning.md` makes final sprint decisions.
- `plan.md` executes `reasoning.md`.
- `review.md` runs the selected review protocols against implementation.
