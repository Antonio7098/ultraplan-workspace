# Sprint Index: Requirements-Driven Performance Stage

> Project: `ultraplan-go`
> Sprint: `39-performance-stage`
> Purpose: selected context for this sprint. Must be a subset of `projects/ultraplan-go/project-index.md`.
> **Inputs Used:** `projects/ultraplan-go/project-index.md`, `projects/ultraplan-go/sprints/39-performance-stage/requirements.md`, `projects/ultraplan-go/roadmap.md`, `projects/ultraplan-go/docs/ARCHITECTURE.md`, `projects/ultraplan-go/docs/PRD.md`, `projects/ultraplan-go/docs/TRD.md`

This document selects what must be read, distilled, reasoned through, or checked for this sprint. It does not make implementation decisions. All selections come from the project index.

## Sprint Scope

- **Sprint Goal:** Add an explicitly enabled performance verification phase after execute and before Conformance Review. It must turn every current requirements-owned target into frozen benchmarks, qualified baseline and candidate measurements, bounded isolated optimization attempts, and product-derived verdicts without weakening correctness or target authority.
- **Planned Output:** Strict project activation and requirements-target contracts; immutable target packets and benchmark manifests; qualified measurements and bounded isolated optimization; private versioned attempt evidence and canonical `performance.md`; additive verification ordering and freshness; shared durable prepare, dry-run, start, status, resume, cancel, recover, and result operations; consistent CLI, JSON, TUI, and browser views; documentation, compatibility tests, race and failure tests, and gated real-runtime evidence.
- **Depends On:** Sprint 38 manual repair passing its end-to-end exit gate and bounded automatic repair retaining fixed limits; Sprint 37 isolated evidence production and smoke compatibility; Sprint 36 verification-phase separation; Sprint 35 durable ownership, writer fencing, cancellation, terminal arbitration, and reconciliation; current execute worktree identity, process execution, disposable-copy isolation, product-owned mutation, and shared application boundaries.
- **Non-Goals:** Implicit or always-on benchmarking, target values outside current sprint requirements, invented numeric thresholds, unrestricted commands or environment forwarding, benchmark rewriting after baseline, correctness or acceptance weakening, indefinite profiling or optimization, distributed benchmarking, historical leaderboards, replacement of Conformance Review or QA, content identity or retrieval, alternate product persistence, cloud or Aren integration, general issue tracking, and Git mutation.

## Source Project Index

- `projects/ultraplan-go/project-index.md` is the authoritative source. Every selected contract, report, reasoning template, and protocol below appears there.

## Selected Contracts

Each contract applies as a flat whole to this sprint. All paths appear in the project index's "Active Contract Pool" table.

| Contract | Why Selected |
| --- | --- |
| Architecture | Governs project and sprint ownership, verification-phase separation, platform process boundaries, shared application interfaces, and the distinction between detailed evidence and bounded flow summaries. |
| CLI Surface | Governs performance command shape, dry-run and lifecycle controls, help, stable text and JSON results, cancellation, recovery, and outcome-specific exit behavior. |
| Configuration | Governs stage-specific model selection, finite defaults, immutable product maxima, lower-only operational limits, effective-source reporting, validation, and redaction. |
| Documentation | Governs architecture, CLI, user, schema, recovery, local-web, and release documentation for activation, targets, outcomes, freshness, cancellation, and recovery. |
| Errors | Governs typed and actionable policy, target, admission, parser, measurement, drift, correctness, persistence, cancellation, cleanup, and recovery failures without false success. |
| LLM Evaluation / Cost / Safety | Governs bounded benchmark authoring and optimization proposals, distrust of model verdicts, safe evidence retention, cost limits, and gated real-runtime evaluation. |
| LLM Runtime | Governs Agentwrap/OpenCode model execution, permissions, cancellation, structured output, runtime identity, and cleanup without granting target, command, comparison, or mutation authority. |
| Observability | Governs durable run correlation, bounded progress and evidence pointers, command and environment identity, cancellation truth, diagnostics, replay, and cross-interface agreement. |
| Performance | Governs benchmark cost, bounded commands and workers, sample and output limits, variance reporting, startup and repository-scan behavior, retention, and finite optimization. |
| Persistence And Migrations | Governs strict versioned state, immutable digest-bound attempt records, private permissions, atomic publication, writer fencing, retention, compatibility, resume, and conservative recovery. |
| Security | Governs sole target authority, protected paths, explicit argv, environment allowlists, isolation, symlink and hard-link defense, source identity, patch containment, secret redaction, and fail-closed cleanup. |
| Testing | Governs parser and comparison fixtures, integration and adversarial coverage, failure injection, race and cancellation tests, compatibility, interface parity, and gated dogfood. |
| Workflows | Governs admission, frozen inputs, durable acceptance, single ownership, ordered measurement and optimization, bounded convergence, cancellation, resume, recovery, and exactly one terminal result. |

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
| Architecture | `projects/ultraplan-go/sprints/39-performance-stage/reasoning/architecture.md` | Reason through activation and target authority, verification ordering, benchmark and implementation isolation, frozen identities, state ownership, promotion boundaries, freshness, and terminal arbitration. |
| API Design | `projects/ultraplan-go/sprints/39-performance-stage/reasoning/api-design.md` | Reason through prepare, dry-run, start, status, resume, cancel, recover, result, and bounded evidence queries; stable outcomes; confirmation; stale requests; and interface compatibility. |
| Frontend | `projects/ultraplan-go/sprints/39-performance-stage/reasoning/frontend.md` | Reason through policy and target presentation, preparation and progress, per-target verdicts, blockers, stale state, cancellation, recovery, hostile evidence bounds, accessibility, and no-JavaScript behavior. |

## Selected Project Reasoning

No project reasoning document is selected. Project reasoning is optional, and project status reports no accepted current project synthesis.

| Document | Path | Why Selected |
| --- | --- | --- |

## Prior Decisions To Carry Forward

All decision paths must appear in the project index's "Prior Decisions" table. The project index currently lists no prior decisions.

| Decision | Path | Constraint For This Sprint |
| --- | --- | --- |
| None cataloged | Not applicable | Carry dependencies, authority boundaries, compatibility requirements, and evidence gates from the governed requirements and project documents, but do not invent a prior-decision artifact. |

## Required Review Protocols

All paths appear in the project index's "Review Protocols" table.

| Protocol | Path | Required Evidence |
| --- | --- | --- |
| Architecture Review | `system/protocols/architecture-review-protocol.md` | Performance remains a sprint-owned verification phase; project policy controls activation only; platform and runtime boundaries cannot bypass targets, frozen benchmarks, qualification, correctness, isolation, promotion, persistence, freshness, or terminal authority. |
| Sprint Review | `system/protocols/review-sprint-protocol.md` | Requirement and selected-contract coverage for disabled compatibility, strict target parsing, admission, frozen identities, qualification, bounded optimization, protected paths, outcomes, freshness, cancellation, recovery, interface parity, docs, tests, race checks, and build evidence. |
| Deep Smoke Sprint | `system/protocols/deep-smoke-sprint-protocol.md` | Gated real-runtime proof with required absolute and baseline-relative targets plus a report or baseline target, including a rejected correctness regression, contained improvement or truthful non-acceptance, stable benchmarks, cancellation, recovery, and cross-interface agreement. |

## Excluded Context

| Context | Reason Excluded | Revisit If |
| --- | --- | --- |
| Implementation execution during index creation | This artifact selects performance context. It does not execute `plan.md`, author benchmarks, run measurements, propose patches, or choose implementation details. | Validated reasoning and plan artifacts authorize governed execute work. |
| Independent smoke investigation implementation | Performance precedes Conformance Review and QA. It must not add a second smoke discovery, execution, evidence, or verdict path. | A later governed sprint explicitly changes smoke authority after compatibility evidence. |
| Conformance Review automation changes | Performance gates the current Conformance Review and can make later evidence stale, but it cannot rewrite review prompts, evidence, verdict rules, or `review.md`. | A later sprint explicitly scopes review changes while preserving independent authority. |
| General-purpose issue tracking | Performance records bounded target misses, blockers, hypotheses, and proposals. It does not add assignment, scheduling, backlog management, or remote synchronization. | A later product phase defines separate issue-management authority and requirements. |
| Git mutation | Performance runtimes and commands cannot add, commit, push, branch, merge, reset, checkout, stash, clean, edit hooks, or mutate the Git index. | A later explicit requirement governs a specific Git operation. |
| Target sources outside requirements | Project policy, configuration, environment, flags, runtime output, benchmark files, stored baselines, and prior artifacts cannot add, remove, replace, or relax target rows or values. | A governed edit changes current sprint requirements before a new attempt starts. |
| Arbitrary benchmark commands, parsers, or environments | Product code owns explicit command descriptors, parser versions, working directories, environment allowlists, timeouts, output bounds, and cleanup policy. | A separate governed requirement expands a versioned product-owned descriptor or parser set. |
| Benchmark or correctness weakening | Frozen benchmarks, tests, fixtures, correctness commands, sample policy, acceptance criteria, and prior evidence cannot change to manufacture a pass. | A separate governed change occurs before a new target packet and attempt are frozen. |
| Unbounded optimization and broad redesign | Each cycle addresses one target-linked miss under finite attempts, commands, changed paths and bytes, output, cleanup, and wall time. Distributed search, migrations, upgrades, and unrelated refactors remain out of scope. | Later dogfood proves a separately governed bounded expansion is safe and necessary. |
| Historical performance service and fleet scheduling | Sprint 39 produces current sprint evidence, not a performance database, cloud runner, cross-project leaderboard, distributed benchmark service, or fleet scheduler. | A later product phase proves a concrete need and defines authority, identity, retention, and trust boundaries. |
| Content identity, retrieval, and knowledge graphs | Digest-bound performance records remain sufficient. Global identity, provenance, retrieval, embeddings, and graph projections stay gated until verification dogfood. | Post-Sprint-40 evidence justifies the content-contract gate. |
| Alternate product persistence | Performance details remain in strict sprint-owned filesystem records, while Sprint 35 SQLite remains authoritative only for operational run facts. | A later authority decision proves and governs another authored-artifact storage mode. |
| Hosted service, remote workers, and Aren integration | Performance remains local and same-host through the current CLI, TUI, and loopback browser boundaries. It does not change Aren's phases or activate performance work for Aren. | A later product phase defines remote identity, authorization, execution, cleanup, and artifact authority. |

## Next Artifacts

- `technical-handbook.md` reads from the evidence reports listed above.
- `reasoning/*.md` captures area-specific reasoning.
- `reasoning.md` makes final sprint decisions.
- `plan.md` executes `reasoning.md`.
- `review.md` runs the selected review protocols against implementation.
