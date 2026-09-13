# Sprint index: Core lifecycle

> Project: `aren-phase-01-execution-lifecycle`
> Sprint: `01-core-lifecycle`
> Purpose: select the contracts, evidence, project conclusions, area templates, and review protocols needed to establish one coherent in-process lifecycle.
> Inputs Used: `projects/aren-phase-01-execution-lifecycle/project-index.md`, `projects/aren-phase-01-execution-lifecycle/project-reasoning/index.md`, `projects/aren-phase-01-execution-lifecycle/roadmap.md`, and the accepted project-reasoning outputs.

This index is ready for use after project reasoning receives a passing review. It does not waive the project reasoning admission gate.

## Sprint Scope

- **Sprint Goal:** Establish a complete non-cancellation lifecycle for one supervised in-process work function.
- **Planned Output:** A buildable Aren Go runtime foundation that creates and starts a run, resolves success, returned failure, or panic, publishes one coherent terminal outcome, retains basic lifecycle history, and supports multiple waiters.
- **Depends On:** Accepted and current Phase 1 project reasoning.
- **Non-Goals:** Caller cancellation, cancellation disposition, parent-cancellation arbitration, multiple live observers, replay cursors, abandoned-observer hardening, persistence, providers, tools, subprocesses, workflows, daemon hosting, remote execution, and a general executor interface.

## Source Project Index

- `projects/aren-phase-01-execution-lifecycle/project-index.md` is the authoritative catalog.
- `projects/aren-phase-01-execution-lifecycle/project-reasoning/index.md` is the authoritative project-reasoning manifest.

## Selected Contracts

| Contract | Why Selected |
|---|---|
| Aren Development Doctrine | Keeps Sprint 1 focused on one runnable lifecycle slice, earned abstractions, strict exclusions, and real behaviour before framework machinery. |
| Aren Runtime Architecture | Governs per-run authority, package/public boundaries, thin entrypoints, one mutation owner, and concrete-before-interface design. |
| Aren Execution Lifecycle | Governs identity, legal transitions, success/error/panic resolution, one terminal commitment, timing, coherent publication, and invariant visibility. |
| Aren Cancellation And Lifetimes | Selected narrowly for invocation-carrier ownership, context/resource lifetime, release paths, and truthful boundaries; caller cancellation semantics remain Sprint 2. |
| Aren Events Observation And Waiting | Applies the Sprint 1 subset for canonical retained history, sequence identity, defensive event publication, and reusable multi-waiter outcome observation. Live cursors and abandonment remain deferred. |
| Aren Verification | Governs independent transition/outcome oracles, controlled schedules, negative controls, race execution, mutation attacks, release evidence, and real runtime proof. |
| Aren Observability | Requires the first lifecycle slice to be explainable from structured canonical evidence, keeps diagnostics distinct from truth, and requires observation to remain passive without introducing a telemetry platform. |
| Aren Performance Engineering | Applies qualitative Phase 1 invariants now: no accidental global serialization, owned runtime-task growth, explainable retention, correctness-before-speed, and evidence before optimization. Sprint 3 owns the formal method and baseline. |
| Documentation | Requires stable lifecycle behavior and realized decisions to remain inspectable and promotable after review. |

## Selected Evidence Reports

| Report | Path | Covers |
|---|---|---|
| Execution Model Taxonomy | `studies/agent-harness-study/reports/final/01.01-execution-model-taxonomy.md` | Bounds Aren's supervised unit and separates it from agent turns and workflow tasks. |
| Control-Flow Ownership | `studies/agent-harness-study/reports/final/01.02-control-flow-ownership.md` | Runtime authority, transition ownership, and control overrides. |
| Step / Turn / Task Atomicity | `studies/agent-harness-study/reports/final/01.03-step-turn-task-atomicity.md` | Coherent publication and the difference between semantic and durable atomicity. |
| State Taxonomy And Ownership | `studies/agent-harness-study/reports/final/02.01-state-taxonomy-and-ownership.md` | Lifecycle, outcome, event, timing, and observer-state ownership. |
| Mutation Discipline And State Transitions | `studies/agent-harness-study/reports/final/02.04-mutation-discipline-and-state-transitions.md` | Single mutation entry points, transition guards, and publication boundaries. |
| Completion And Finalization Semantics | `studies/agent-harness-study/reports/final/03.09-completion-and-finalization-semantics.md` | Success, failure, panic, completion, and information preservation. |
| Go Lifecycle Transition Ownership And Terminal Arbitration | `studies/aren-go-runtime-study/reports/final/01.01-lifecycle-transition-ownership-and-terminal-arbitration.md` | Direct Go mechanisms and source references for ownership, one completion slot, terminal arbitration, and publication order. |
| Go Adversarial Concurrency And Failure Verification | `studies/aren-go-runtime-study/reports/final/01.03-adversarial-concurrency-and-failure-verification.md` | Transition matrices, controlled races, negative controls, race execution, and regression strategy. |
| Project Structure | `studies/go-cli-study/reports/final/01-project-structure.md` | Thin entrypoints, package boundaries, and dependency direction. |
| Error Handling | `studies/go-cli-study/reports/final/05-error-handling.md` | Cause preservation, classification, and actionable diagnostics. |
| State And Context | `studies/go-cli-study/reports/final/07-state-context.md` | Context propagation and separation from owned runtime state. |
| Concurrency | `studies/go-cli-study/reports/final/08-concurrency.md` | Goroutine ownership, cancellation-safe structure, and concurrent API behavior. |
| Testing Strategy | `studies/go-cli-study/reports/final/11-testing-strategy.md` | Unit, command, integration, race, and fixture boundaries. |
| Philosophy | `studies/go-cli-study/reports/final/15-philosophy.md` | Scope control and deliberate rejection of unearned complexity. |

## Selected Reasoning Templates

| Template | Output Path | Why Selected |
|---|---|---|
| Aren Runtime Architecture And Authority | `projects/aren-phase-01-execution-lifecycle/sprints/01-core-lifecycle/reasoning/runtime-architecture-and-authority.md` | Sprint 1 must decide package ownership, dependency construction, runtime capabilities, and the work invocation boundary. |
| Aren Lifecycle Transition And Publication | `projects/aren-phase-01-execution-lifecycle/sprints/01-core-lifecycle/reasoning/lifecycle-transition-and-publication.md` | Sprint 1 must settle states, legal transitions, the linearization point, and coherent terminal publication. |
| Aren Outcome And Failure Contract | `projects/aren-phase-01-execution-lifecycle/sprints/01-core-lifecycle/reasoning/outcome-and-failure-contract.md` | Sprint 1 must settle result, returned error, panic, timing, invariant failure, and immutable outcome construction. |
| Aren Concurrency Waiters And Immutability | `projects/aren-phase-01-execution-lifecycle/sprints/01-core-lifecycle/reasoning/concurrency-waiters-and-immutability.md` | Sprint 1 must settle single-writer synchronization, concurrent inspection, multiple waiters, goroutine ownership, and defensive publication. |
| Aren Event History And Observation | `projects/aren-phase-01-execution-lifecycle/sprints/01-core-lifecycle/reasoning/event-history-and-observation.md` | Sprint 1 must settle the minimum canonical history and sequence semantics without predesigning Sprint 2 observers. |
| Aren Lifecycle Verification Model | `projects/aren-phase-01-execution-lifecycle/sprints/01-core-lifecycle/reasoning/lifecycle-verification-model.md` | Sprint 1 needs an independent state model, invariant catalog, controlled schedules, negative controls, and race evidence. |

## Selected Project Reasoning

| Document | Path | Why Selected |
|---|---|---|
| Project synthesis | `projects/aren-phase-01-execution-lifecycle/project-reasoning/reasoning.md` | Supplies the accepted Phase 1 constraints and conclusions that all sprint decisions must obey. |
| Evidence assessment | `projects/aren-phase-01-execution-lifecycle/project-reasoning/areas/00-evidence-map.md` | Preserves evidence quality, applicability limits, contradictions, and negative-transfer warnings used by the handbook. |
| Lifecycle authority | `projects/aren-phase-01-execution-lifecycle/project-reasoning/areas/01-lifecycle-authority-and-atomic-publication.md` | Directly governs Sprint 1 transition ownership and coherent publication. |
| Outcomes and terminal resolution | `projects/aren-phase-01-execution-lifecycle/project-reasoning/areas/02-outcomes-failures-and-terminal-resolution.md` | Directly governs success, returned failure, panic, and terminal outcome construction. |
| Cancellation and cleanup | `projects/aren-phase-01-execution-lifecycle/project-reasoning/areas/03-cancellation-goroutine-ownership-and-cleanup.md` | Supplies the goroutine ownership and cleanup constraints Sprint 1 must preserve even though caller cancellation is deferred. |
| Events and observation | `projects/aren-phase-01-execution-lifecycle/project-reasoning/areas/04-events-observation-waiting-and-replay.md` | Governs canonical history and the Sprint 1 seams that must remain safe for Sprint 2 extension. |
| Verification and Go correctness | `projects/aren-phase-01-execution-lifecycle/project-reasoning/areas/05-verification-and-go-correctness.md` | Governs invariant, schedule, race, negative-control, and evidence obligations. |

## Prior Decisions To Carry Forward

| Decision | Path | Constraint For This Sprint |
|---|---|---|
| Phase 1 product contract | `projects/aren-phase-01-execution-lifecycle/docs/PRD.md` | Product behavior and exclusions are fixed; detailed Go design remains open to sprint reasoning. |
| Go language decision | `projects/aren-phase-01-execution-lifecycle/docs/final-language-decision.md` | Implementation, tests, concurrency reasoning, and public contracts must follow the accepted Go rules. |
| Accepted project reasoning | `projects/aren-phase-01-execution-lifecycle/project-reasoning/reasoning.md` | Sprint reasoning may narrow implementation choices but cannot silently contradict an accepted project conclusion. |

## Required Review Protocols

| Protocol | Path | Required Evidence |
|---|---|---|
| Architecture Review | `system/protocols/architecture-review-protocol.md` | Package ownership, dependency direction, lifecycle authority, synchronization boundaries, and public/internal capability review. |
| Sprint Review | `system/protocols/review-sprint-protocol.md` | Acceptance-criterion traceability, implementation conformance, test evidence, and unresolved findings. |

## Excluded Context

| Context | Reason Excluded | Revisit If |
|---|---|---|
| Implementation execution during planning | The planning flow must settle requirements, evidence, decisions, and tasks before source mutation starts. | `plan.md` validates and execution begins. |
| Smoke investigation during planning | Sprint 1 has no implemented runtime to investigate during index, handbook, or reasoning stages. | Execution completes and a selected smoke protocol has a real target. |
| Review automation during planning | Review must examine the realized implementation and current evidence rather than predicted code. | Execution produces reviewable changes. |
| Issue tracking mutation | Findings belong in sprint review artifacts before any external issue workflow is considered. | Review identifies a deferred issue that needs separate ownership. |
| Git mutation during planning | Planning artifacts do not authorize branch, commit, merge, or repository-history changes. | The governed execute or merge stage owns the operation. |
| Caller cancellation and replay APIs | These are Sprint 2 concerns. Speculative Sprint 1 interfaces would constrain the harder concurrency design before evidence exists. | Sprint 2 requirements and realized Sprint 1 behavior are available. |
| Providers, tools, subprocesses, persistence, workflows, daemon hosting, and remote execution | These belong to later Aren phases and would distort the in-memory lifecycle proof. | The owning Aren phase begins. |

## Next Artifacts

- `technical-handbook.md` distills the selected reports and records how accepted project reasoning applies.
- `reasoning/*.md` resolves the six selected Sprint 1 decision areas.
- `reasoning.md` makes the final Sprint 1 decisions.
- `plan.md` implements only those decisions and records the project-reasoning digest lock after validation.
- `review.md` runs the selected protocols against the realized implementation.
