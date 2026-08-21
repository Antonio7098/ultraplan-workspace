# Sprint Index: Shared Context Integration and Grounded-Planning Release

> Project: `ultraplan-go`
> Sprint: `34-shared-context`
> Purpose: selected context for this sprint. Must be a subset of `projects/ultraplan-go/project-index.md`.
> **Inputs Used:** `projects/ultraplan-go/project-index.md`, `projects/ultraplan-go/sprints/34-shared-context/requirements.md`, `projects/ultraplan-go/roadmap.md`, `projects/ultraplan-go/docs/ARCHITECTURE.md`, `projects/ultraplan-go/docs/PRD.md`, `projects/ultraplan-go/docs/TRD.md`

This document selects what must be read, distilled, reasoned through, or checked for this sprint. It does not make implementation decisions. All selections must come from the project index — no items may be included that are not listed in the project index.

## Sprint Scope

- **Sprint Goal:** Reuse the stored, validated `requirements.md` and reference-only `code-context.md` unchanged through one stable shared prompt prefix, resolve its repository-relative references into transient prompt evidence, and apply that prefix across every downstream agent-backed sprint operation.
- **Planned Output:** A grounded-planning release comprising shared sprint-context rendering across planning, execute, Conformance Review, and smoke authoring; exact-once cumulative flow integration; a manual code-context skill; deterministic conformance tests; aligned user and architecture documentation; and truthful gated dogfood evidence.
- **Depends On:** Sprint 33's accepted code-context stage, validation, artifact, compatibility, cancellation, and atomic-rerun behavior; Sprint 32's shared application/interface boundary and truthful gated-evidence standard; the configured implementation repository and an available gated runtime for real-repository dogfood.
- **Non-Goals:** Repository indexing, retrieval, RAG, embeddings, cache ownership or provider cache-hit guarantees, a parallel context manifest, automatic staleness or amendment, content identity or provenance, QA or repair, a generic stage/workflow/plugin framework, route-specific browser workflow semantics, hosted or multi-user service, remote workers, issue tracking, and automatic Git mutation.

## Source Project Index

- `projects/ultraplan-go/project-index.md` — authoritative source. Any file or item referenced below must appear there.

## Selected Contracts

Each contract applies as a flat whole to this sprint. All paths must appear in the project index's "Active Contract Pool" table.

| Contract | Why Selected |
| --- | --- |
| Architecture | Governs `internal/sprint` ownership of shared composition, dependency direction, shared application boundaries, and the required separation from generic runtime infrastructure. |
| Errors | Governs actionable failures for source-reference resolution, validation, cancellation, unavailable dogfood prerequisites, and downstream prompt construction. |
| Observability | Governs truthful runtime and flow diagnostics while keeping timestamps, run IDs, attempts, sessions, and other dynamic execution data outside the stable common prefix. |
| Security | Governs implementation-repository containment, safe reference resolution, bounded diagnostics, redaction, permission boundaries, and restricted mutation. |
| Testing | Governs deterministic fake-runtime fixtures, exact-byte assertions, cross-stage integration coverage, regression protection, race checks, and truthful gated tests. |
| Documentation | Governs the required README, CLI, user, architecture, recovery, planning-smoke, stage-skill, local-web, and generated-workspace documentation alignment. |
| CLI Surface | Governs cumulative flow behavior, code-context command documentation, status behavior, and manual skill delegation to the canonical CLI operation. |
| LLM Runtime | Governs deterministic prompt composition, runtime request boundaries, agent-backed downstream coverage, permissions, and generic agentwrap integration. |
| LLM Evaluation / Cost / Safety | Governs bounded prompt evidence, safe runtime use, truthful gated dogfood reporting, and avoidance of unsupported provider cache claims. |
| Workflows | Governs canonical stage ordering, exact-once code-context execution, cancellation and recovery, state transitions, and integration across downstream agent-backed operations. |
| Performance | Governs prompt-budget bounds and deterministic source-evidence injection without repository-wide indexing or accidental unbounded input growth. |
| Persistence And Migrations | Governs exact stored artifact reuse, atomic code-context replacement, durable flow-state compatibility, and preservation of the last valid artifact after unsuccessful reruns. |

## Selected Evidence Reports

Copied from the project index's "Available Evidence Reports" table. These tell the technical handbook which reports to read – the project index is the authoritative source.

| Report | Path | Covers |
| --- | --- | --- |
| `01-project-structure` | `studies/go-cli-study/reports/final/01-project-structure.md` | Project layout, cmd/internal/pkg, dependency direction, thin entrypoints |
| `02-command-architecture` | `studies/go-cli-study/reports/final/02-command-architecture.md` | Command routing, flags, help text, command organization, shell completion |
| `03-dependency-injection` | `studies/go-cli-study/reports/final/03-dependency-injection.md` | Dependency construction, seams, testability, constructor patterns |
| `05-error-handling` | `studies/go-cli-study/reports/final/05-error-handling.md` | Error wrapping, classification, user-facing diagnostics, sentinel errors |
| `06-io-abstraction` | `studies/go-cli-study/reports/final/06-io-abstraction.md` | Filesystem/stdin/stdout abstraction, test seams, interface design |
| `07-state-context` | `studies/go-cli-study/reports/final/07-state-context.md` | Context propagation, app state, cancellation |
| `10-logging-observability` | `studies/go-cli-study/reports/final/10-logging-observability.md` | Logs, diagnostics, structured events, observability patterns |
| `11-testing-strategy` | `studies/go-cli-study/reports/final/11-testing-strategy.md` | Unit, integration, fixture, command-level tests, coverage strategy |
| `13-security` | `studies/go-cli-study/reports/final/13-security.md` | Path safety, secrets, command injection risks, sandboxing |
| `14-performance` | `studies/go-cli-study/reports/final/14-performance.md` | Startup latency, large repos, memory management, performance |
| `15-philosophy` | `studies/go-cli-study/reports/final/15-philosophy.md` | Cross-cutting design philosophy, tradeoffs, maintainability |

## Selected Reasoning Templates

All paths must appear in the project index's "Available Reasoning Templates" table.

| Template | Output Path | Why Selected |
| --- | --- | --- |
| Architecture | `projects/ultraplan-go/sprints/34-shared-context/reasoning/architecture.md` | Resolve shared-renderer ownership, composition boundaries, downstream integration seams, source-evidence handling, and separation between sprint semantics and generic runtime infrastructure. |

## Prior Decisions To Carry Forward

All decision paths must appear in the project index's "Prior Decisions" table.

| Decision | Path | Constraint For This Sprint |
| --- | --- | --- |
| None cataloged | None | The project index states that no prior decisions are cataloged, so no unindexed decision artifact is selected; the Sprint 33 and Sprint 32 carry-forward constraints remain explicit sprint dependencies above. |

## Required Review Protocols

All paths must appear in the project index's "Review Protocols" table.

| Protocol | Path | Required Evidence |
| --- | --- | --- |
| Architecture Review | `system/protocols/architecture-review-protocol.md` | Evidence that `internal/sprint` owns one deterministic shared renderer and product-stage composition, generic runtime code remains free of sprint semantics, source resolution is contained, and all interfaces retain shared application boundaries. |
| Sprint Review | `system/protocols/review-sprint-protocol.md` | Evidence that required integrations, exact-prefix and exact-byte fixtures, flow sequencing, cancellation and rerun regressions, documentation, release checks, and truthful gated dogfood satisfy the sprint requirements. |

## Excluded Context

| Context | Reason Excluded | Revisit If |
| --- | --- | --- |
| Implementation execution | This index selects planning context and does not run implementation tasks or mutate the implementation repository; execute prompt integration is planned scope, not permission for this artifact stage to execute it. | Validated reasoning and plan artifacts authorize the execute stage. |
| Smoke investigation | Shared-prefix integration for agent-backed smoke authoring and deterministic smoke tests are in scope, but live smoke investigation and harness mutation are not sprint-index work or part of the requirements-to-plan dogfood. | A current acceptable review and an explicitly authorized smoke stage require runtime-facing evidence. |
| Review automation | Existing Conformance Review requests must receive the common prefix, but new reviewer fan-out, verdict ownership, protocols, or review workflow behavior are outside this sprint. | Requirements explicitly call for changing review orchestration rather than preserving it. |
| Issue tracking | General-purpose issue records, assignment, scheduling, synchronization, and project-management behavior are deferred. | A later roadmap sprint explicitly promotes an issue-management capability. |
| Git mutation | Automatic add, commit, push, branch, merge, reset, or other Git-state mutation is prohibited. | A later explicit requirement and contract-approved workflow brings a specific Git operation into scope. |
| Repository indexing, retrieval, RAG, and embeddings | The context pack remains a curated reference artifact and downstream agents retain live repository access; no retrieval subsystem is needed for shared prompt reuse. | A later evidence gate promotes measured retrieval work. |
| Cache systems and provider cache guarantees | The sprint creates a stable prefix but neither owns cache keys nor claims or measures provider cache hits. | A later requirement explicitly introduces measured cache behavior with a supported provider contract. |
| Parallel manifests, staleness, provenance, and QA/repair | The release preserves one Markdown context pack and does not add identity, source-change detection, amendment, empirical QA, adjudication, or repair workflows. | Post-Sprint-34 evidence gates promote one of these capabilities. |
| New generic frameworks or runtime contracts | The sprint integrates existing product stages and agentwrap boundaries rather than adding a workflow engine, plugin system, repository-index package, or competing runtime contract. | Repeated concrete implementations demonstrate a separately approved shared abstraction. |
| Hosted, remote, and multi-user behavior | The release remains local-first and does not add cloud authority, remote binding, accounts, collaboration, remote workers, or Aren integration. | A later roadmap authority and cloud gate is explicitly satisfied. |

## Next Artifacts

- `technical-handbook.md` reads from the evidence reports listed above.
- `reasoning/*.md` captures area-specific reasoning.
- `reasoning.md` makes final sprint decisions.
- `plan.md` executes `reasoning.md`.
- `review.md` runs the selected review protocols against implementation.
