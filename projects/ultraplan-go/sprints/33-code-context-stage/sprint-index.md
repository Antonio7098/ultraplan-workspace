# Sprint Index: Code-Context Stage Vertical Slice

> Project: `ultraplan-go`
> Sprint: `33-code-context-stage`
> Purpose: selected context for this sprint. Must be a subset of `projects/ultraplan-go/project-index.md`.
> **Inputs Used:** `projects/ultraplan-go/project-index.md`, `projects/ultraplan-go/sprints/33-code-context-stage/requirements.md`, `projects/ultraplan-go/roadmap.md`, `projects/ultraplan-go/docs/ARCHITECTURE.md`, `projects/ultraplan-go/docs/PRD.md`, `projects/ultraplan-go/docs/TRD.md`

This document selects what must be read, distilled, reasoned through, or checked for this sprint. It does not make implementation decisions. All selections must come from the project index - no items may be included that are not listed in the project index.

## Sprint Scope

- **Sprint Goal:** Make `code-context` a fully operational sprint stage immediately after requirements that can inspect the resolved implementation repository, atomically produce and validate the single authoritative `code-context.md` artifact, and expose truthful operation and recovery behavior through existing shared application and web boundaries.
- **Planned Output:** A vertical slice covering the code-context stage domain and artifact model, flow and persisted-state compatibility, runtime-backed generation and structural validation, embedded defaults, stage-specific configuration, CLI and shared application wiring, browser presentation through generic operations, and focused deterministic tests.
- **Depends On:** Sprint 32's shared application and web extensibility baseline; the existing sprint runtime, flow-state, defaults, configuration, operation, and web-operation implementations; and the project implementation target cataloged by the project index.
- **Non-Goals:** Downstream shared-prefix injection, a manual code-context skill, broad documentation and real-runtime release dogfood reserved for Sprint 34; repository indexing, RAG, embeddings, caches, parallel context manifests, automatic staleness or amendment; expanded QA or repair, alternate persistence, knowledge graphs, cloud or Aren integration; and new generic stage, workflow, plugin, or web-product frameworks.

## Source Project Index

- `projects/ultraplan-go/project-index.md` - authoritative source. Any file or item referenced below must appear there.

## Selected Contracts

Each contract applies as a flat whole to this sprint. All paths must appear in the project index's "Active Contract Pool" table.

| Contract | Why Selected |
| --- | --- |
| Architecture | Governs sprint ownership of stage semantics, product/platform dependency direction, shared app use cases, and the prohibition on route-specific web workflow logic. |
| CLI Surface | Applies to prompt, validation, flow, help, status, dry-run, and stable JSON behavior for the new stage. |
| Configuration | Governs stage-specific model and variant fallback, source tracking, validation, and safe configuration projection. |
| Documentation | Applies to executable help and synchronized embedded prompt/template defaults while broad documentation remains deferred. |
| Errors | Governs actionable validation, runtime, cancellation, compatibility, and recovery diagnostics with stable classifications. |
| LLM Evaluation / Cost / Safety | Applies to safe runtime metadata, output validation, permission boundaries, and the rule that runtime exit success is insufficient. |
| LLM Runtime | Governs generic runtime execution, context propagation, model/variant selection, expected output handling, and adapter boundaries. |
| Observability | Applies to truthful status, progress, findings, cancellation, recovery, and bounded safe metadata across CLI and browser projections. |
| Persistence And Migrations | Governs atomic artifact and flow-state writes plus deterministic compatibility for persisted pre-code-context state. |
| Security | Applies to contained repository reads, single-artifact writes, safe previews, redaction, and prohibition of source or Git mutation. |
| Testing | Governs deterministic fake-runtime, fixture, command, app, web, compatibility, cancellation, and atomic-write coverage. |
| Workflows | Applies to stage ordering, prerequisites, cumulative dispatch, validation gating, reruns, cancellation, and truthful terminal transitions. |

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

## Selected Reasoning Templates

All paths must appear in the project index's "Available Reasoning Templates" table.

| Template | Output Path | Why Selected |
| --- | --- | --- |
| Architecture | `projects/ultraplan-go/sprints/33-code-context-stage/reasoning/architecture.md` | Reason through stage ownership, flow and compatibility integration, runtime boundaries, artifact persistence, and shared app/web dependency direction. |
| API Design | `projects/ultraplan-go/sprints/33-code-context-stage/reasoning/api-design.md` | Reason through stable CLI/JSON projections and reuse of generic operation, cancellation, status, finding, and recovery capabilities. |
| Frontend | `projects/ultraplan-go/sprints/33-code-context-stage/reasoning/frontend.md` | Reason through code-context readiness, findings, artifact access, progress, and explicit rerun presentation within the established server-rendered browser hierarchy. |

## Prior Decisions To Carry Forward

All decision paths must appear in the project index's "Prior Decisions" table.

| Decision | Path | Constraint For This Sprint |
| --- | --- | --- |
| None cataloged | Not applicable | The project index states that no prior decisions are cataloged; no prior-decision path can be selected or invented. Requirements carry the applicable Sprint 32 shared-boundary, redaction, compatibility, cancellation, and fake-first constraints. |

## Required Review Protocols

All paths must appear in the project index's "Review Protocols" table.

| Protocol | Path | Required Evidence |
| --- | --- | --- |
| Architecture Review | `system/protocols/architecture-review-protocol.md` | Evidence that stage semantics remain in `internal/sprint`, runtime infrastructure remains generic, CLI and web use shared app boundaries, and no parallel workflow or persistence abstraction is introduced. |
| Sprint Review | `system/protocols/review-sprint-protocol.md` | Evidence that required outputs and acceptance criteria are implemented, focused tests and release commands pass, deferred scope is absent, and repository/output mutation boundaries hold. |

## Excluded Context

| Context | Reason Excluded | Revisit If |
| --- | --- | --- |
| Implementation execution | This sprint adds code-context generation, not execute-stage plan-task behavior or implementation mutation by the code-context runtime. | A later sprint explicitly changes execute semantics or grants the stage implementation-write scope. |
| Smoke investigation | Real-runtime requirements-to-plan dogfood and Phase 5 release smoke are Sprint 34 scope; this sprint requires deterministic fake-runtime evidence. | Sprint 34 begins or a gated environment is explicitly approved as supplemental evidence. |
| Review automation | Existing conformance-review behavior is not being redesigned; only compatibility with generic shared boundaries is relevant. | Requirements explicitly change review orchestration, reviewers, verdicts, or review artifacts. |
| Issue tracking | General-purpose issue records, assignment, scheduling, and synchronization remain deferred product scope. | A future roadmap gate explicitly promotes issue-management behavior. |
| Git mutation | The code-context stage is read-only toward the implementation repository and must not alter Git state. | A future explicit contract and sprint requirement authorizes a narrowly scoped Git operation. |
| Downstream shared context injection | Exact requirements/code-context prefix reuse across sprint index, handbook, reasoning, plan, execute, review, and smoke is reserved for Sprint 34. | Sprint 34 shared-context integration starts. |
| Manual code-context skill | Creating or materializing the manually invoked stage skill is reserved for Sprint 34. | Sprint 34 shared-context integration starts. |
| Broad documentation and release dogfood | This vertical slice covers executable help and defaults only; broad guides, recovery documentation, generated-workspace docs, and real-runtime release proof are deferred. | Sprint 34 release work starts. |
| Repository indexing, retrieval, and caching | The artifact is a curated Markdown context pack, not an index, RAG/embedding system, cache subsystem, cache key, provider cache-control feature, or parallel JSON manifest. | A later evidence gate explicitly promotes retrieval or caching work. |
| Automatic staleness and context amendments | Source-change detection, automatic refresh, and amendment workflows are excluded from the initial stage. | A later sprint defines explicit freshness and amendment semantics. |
| Content identity, expanded QA/repair, alternate persistence, graph, cloud, and Aren | These capabilities are gated beyond the Sprint 34 grounded-planning release and would expand this vertical slice. | The roadmap promotes the corresponding evidence-gated product phase. |
| New generic frameworks | Generic stage frameworks, plugin registries, workflow engines, repository-index packages, web-specific product services, and package-global mutable runner registration are non-goals. | Repeated implemented use cases establish a concrete need and a later sprint selects the abstraction. |

## Next Artifacts

- `technical-handbook.md` reads from the evidence reports listed above.
- `reasoning/*.md` captures area-specific reasoning.
- `reasoning.md` makes final sprint decisions.
- `plan.md` executes `reasoning.md`.
- `review.md` runs the selected review protocols against implementation.
