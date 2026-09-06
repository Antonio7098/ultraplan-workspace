# Sprint index: Cancellation and concurrency

> Project: `aren-phase-01-execution-lifecycle`
> Sprint: `02-cancellation-and-concurrency`
> Purpose: select the contracts, evidence, accepted project conclusions, Sprint 1 decisions, area templates, and review protocols needed to prove cancellation and observation under adversarial concurrency.
> Inputs Used: `projects/aren-phase-01-execution-lifecycle/project-index.md`, `projects/aren-phase-01-execution-lifecycle/project-reasoning/index.md`, `projects/aren-phase-01-execution-lifecycle/roadmap.md`, the accepted project-reasoning outputs, and completed Sprint 1 planning and review artifacts.

This index is prepared in advance. Before Sprint 2 starts, confirm that every named Sprint 1 artifact exists and update any selection affected by the realized design.

## Sprint Scope

- **Sprint Goal:** Prove that the Sprint 1 lifecycle remains truthful, deterministic, observable, and leak-resistant when explicit cancellation, parent cancellation, replayable observers, and adversarial concurrency are added.
- **Planned Output:** Cancellation acceptance and disposition, deterministic terminal resolution, replayable multi-observer history, bounded resource ownership, a diagnostic CLI, race and stress evidence, the final lifecycle contract, and the Phase 1 review.
- **Depends On:** Accepted project reasoning plus completed and reviewed Sprint 1 planning, implementation, tests, and decisions.
- **Non-Goals:** Providers, model messages, token streaming, tools, subprocesses, persistence, restart recovery, pause and resume, durable workflows, remote execution, daemon hosting, network APIs, production telemetry exporters, and general executor abstractions.

## Source Project Index

- `projects/aren-phase-01-execution-lifecycle/project-index.md` is the authoritative catalog.
- `projects/aren-phase-01-execution-lifecycle/project-reasoning/index.md` is the authoritative project-reasoning manifest.

## Selected Contracts

| Contract | Why Selected |
|---|---|
| Architecture | Governs retained Sprint 1 ownership, any justified architecture delta, dependency direction, and capability boundaries. |
| Errors | Governs cancellation causes, returned errors, panic, cleanup failures, terminal classification, and diagnostic identity. |
| Observability | Governs cancellation events, canonical history, sequence ordering, replay, observer delivery, and correlation. |
| Testing | Requires deterministic collisions, negative controls, race runs, repeated stress, leak checks, CLI evidence, and acceptance traceability. |
| Documentation | Requires final lifecycle semantics, superseding decisions, public behavior, and Phase 1 evidence to be recorded. |
| CLI Surface | Governs `aren dev run` diagnostics, injected IO, exit behavior, and command-level evidence. |
| Workflows | Applies only local cancellation, terminal-state, concurrency, and lifecycle-history clauses. Persistence and durable replay remain excluded. |
| Performance | Governs bounded observer resources, shutdown latency, leak resistance, and stress limits. |

## Selected Evidence Reports

| Report | Path | Covers |
|---|---|---|
| Control-Flow Ownership | `studies/agent-harness-study/reports/final/01.02-control-flow-ownership.md` | Cancellation acceptance, runtime authority, and terminal transition ownership. |
| Termination And Loop Bounds | `studies/agent-harness-study/reports/final/01.04-termination-and-loop-bounds.md` | Stop reasons and runtime-owned termination, excluding agent-loop budgets. |
| Pause / Resume / Interrupt | `studies/agent-harness-study/reports/final/01.05-pause-resume-interrupt-semantics.md` | Request, acceptance, acknowledgement, and interruption distinctions. |
| Concurrency And Parallel Advancement | `studies/agent-harness-study/reports/final/01.07-concurrency-and-parallel-advancement.md` | Single-writer commitment, ownership trees, cancellation propagation, and sibling behavior. |
| Delivery Guarantees And Idempotency | `studies/agent-harness-study/reports/final/01.09-delivery-guarantees-and-idempotency.md` | Canonical recording, replay duplication, delivery limits, and exactly-once warnings. |
| Replay And Determinism | `studies/agent-harness-study/reports/final/01.10-replay-and-determinism.md` | Sequence replay, deterministic reconstruction boundaries, and ordering. |
| Mutation Discipline And State Transitions | `studies/agent-harness-study/reports/final/02.04-mutation-discipline-and-state-transitions.md` | Transition gates, concurrent mutation, and coherent publication. |
| Completion And Finalization Semantics | `studies/agent-harness-study/reports/final/03.09-completion-and-finalization-semantics.md` | Completion, cancellation, panic, cleanup, and terminal information preservation. |
| Go Lifecycle Transition Ownership And Terminal Arbitration | `studies/aren-go-runtime-study/reports/final/01.01-lifecycle-transition-ownership-and-terminal-arbitration.md` | Direct Go evidence for one completion slot, terminal arbitration, and publication ordering. |
| Go Cancellation Goroutine Ownership And Cleanup | `studies/aren-go-runtime-study/reports/final/01.02-cancellation-goroutine-ownership-and-cleanup.md` | Direct Go evidence for cancellation causes, joins, cleanup, idempotency, and leak risks. |
| Go Adversarial Concurrency And Failure Verification | `studies/aren-go-runtime-study/reports/final/01.03-adversarial-concurrency-and-failure-verification.md` | Deterministic schedules, property checks, transition matrices, race execution, stress, and failure reproduction. |
| Go Ordered Observation Live Streaming And Backpressure | `studies/aren-go-runtime-study/reports/final/01.04-ordered-observation-live-streaming-and-backpressure.md` | Sequence cursors, replay and live handoff, observer isolation, terminal delivery, gaps, and backpressure. |
| Error Handling | `studies/go-cli-study/reports/final/05-error-handling.md` | Cause preservation, stable diagnostics, and exit mapping. |
| IO Abstraction | `studies/go-cli-study/reports/final/06-io-abstraction.md` | Injectable streams and deterministic diagnostic command tests. |
| State And Context | `studies/go-cli-study/reports/final/07-state-context.md` | Parent lifetime, cancellation causes, wait contexts, and application state. |
| Concurrency | `studies/go-cli-study/reports/final/08-concurrency.md` | Goroutine ownership, safe cancellation, and concurrent callers. |
| Testing Strategy | `studies/go-cli-study/reports/final/11-testing-strategy.md` | Unit, integration, command, race, stress, and fixture boundaries. |
| Philosophy | `studies/go-cli-study/reports/final/15-philosophy.md` | Scope control and resistance to speculative runtime generality. |

## Selected Reasoning Templates

| Template | Output Path | Why Selected |
|---|---|---|
| Aren Cancellation And Terminal Resolution | `projects/aren-phase-01-execution-lifecycle/sprints/02-cancellation-and-concurrency/reasoning/cancellation-and-terminal-resolution.md` | Sprint 2 must settle request acceptance, causes, acknowledgement, dispositions, parent cancellation, and the terminal truth table. |
| Aren Concurrency Interleavings And Goroutine Ownership | `projects/aren-phase-01-execution-lifecycle/sprints/02-cancellation-and-concurrency/reasoning/concurrency-interleavings-and-goroutine-ownership.md` | Sprint 2 must attack the realized design with cancellation, watcher, waiter, observer, shutdown, and inspection races. |
| Aren Observers Replay And Delivery | `projects/aren-phase-01-execution-lifecycle/sprints/02-cancellation-and-concurrency/reasoning/observers-replay-and-delivery.md` | Sprint 2 must settle cursor meaning, replay and live handoff, terminal completion, observer isolation, abandonment, and delivery claims. |
| Aren Control Observation And API Authority | `projects/aren-phase-01-execution-lifecycle/sprints/02-cancellation-and-concurrency/reasoning/control-observation-and-api-authority.md` | Sprint 2 must separate read, wait, observe, and cancel capabilities without exposing transition authority. |
| Aren Race Stress And Leak Verification | `projects/aren-phase-01-execution-lifecycle/sprints/02-cancellation-and-concurrency/reasoning/race-stress-and-leak-verification.md` | Sprint 2 must define controlled collisions, repeated stress, negative controls, leak proof, CLI evidence, and the Phase 1 exit suite. |

`Aren Sprint 2 Architecture Delta` is intentionally not selected now. Select it only if completed Sprint 1 code or review evidence proves that local extension cannot preserve package ownership or dependency direction.

## Selected Project Reasoning

| Document | Path | Why Selected |
|---|---|---|
| Project synthesis | `projects/aren-phase-01-execution-lifecycle/project-reasoning/reasoning.md` | Supplies the accepted Phase 1 constraints and conclusions that Sprint 2 must preserve or explicitly supersede. |
| Evidence assessment | `projects/aren-phase-01-execution-lifecycle/project-reasoning/areas/evidence-assessment.md` | Preserves corpus limitations, contradictions, and negative-transfer warnings while Sprint 2 uses broader concurrency evidence. |
| Lifecycle authority | `projects/aren-phase-01-execution-lifecycle/project-reasoning/areas/lifecycle-authority-and-atomic-publication.md` | Keeps cancellation and observers behind the same Aren-owned lifecycle authority. |
| Outcomes and terminal resolution | `projects/aren-phase-01-execution-lifecycle/project-reasoning/areas/outcomes-failures-and-terminal-resolution.md` | Governs deterministic interpretation of completion, error, panic, and accepted cancellation facts. |
| Cancellation and cleanup | `projects/aren-phase-01-execution-lifecycle/project-reasoning/areas/cancellation-goroutine-ownership-and-cleanup.md` | Directly governs cancellation truthfulness, cause retention, goroutine ownership, joins, and cleanup. |
| Events and observation | `projects/aren-phase-01-execution-lifecycle/project-reasoning/areas/events-observation-waiting-and-replay.md` | Directly governs canonical history, replay, handoff, observers, waiting, and delivery limits. |
| Verification and Go correctness | `projects/aren-phase-01-execution-lifecycle/project-reasoning/areas/verification-and-go-correctness.md` | Governs adversarial schedules, race, stress, negative-control, leak, and Phase 1 exit evidence. |

## Prior Decisions To Carry Forward

| Decision or artifact | Path | Constraint For This Sprint |
|---|---|---|
| Sprint 1 requirements | `projects/aren-phase-01-execution-lifecycle/sprints/01-core-lifecycle/requirements.md` | Preserve the delivered core behavior unless Sprint 2 explicitly supersedes a decision with new evidence. |
| Sprint 1 technical handbook | `projects/aren-phase-01-execution-lifecycle/sprints/01-core-lifecycle/technical-handbook.md` | Reuse applicable evidence conclusions without repeating Sprint 1 analysis. |
| Sprint 1 area reasoning | `projects/aren-phase-01-execution-lifecycle/sprints/01-core-lifecycle/reasoning/` | Classify each relevant decision as preserved, extended, superseded, or unaffected. |
| Sprint 1 decision synthesis | `projects/aren-phase-01-execution-lifecycle/sprints/01-core-lifecycle/reasoning.md` | Treat this as the authoritative Sprint 1 decision set. Silence does not supersede it. |
| Sprint 1 plan | `projects/aren-phase-01-execution-lifecycle/sprints/01-core-lifecycle/plan.md` | Use it to understand intended implementation boundaries, not as authority over realized code. |
| Sprint 1 review | `projects/aren-phase-01-execution-lifecycle/sprints/01-core-lifecycle/review.md` | Carry unresolved findings and confirmed implementation facts into Sprint 2 reasoning. |

Before Sprint 2 starts, replace any missing artifact reference with the completed canonical path and record why an expected artifact is absent.

## Required Review Protocols

| Protocol | Path | Required Evidence |
|---|---|---|
| Architecture Review | `system/protocols/architecture-review-protocol.md` | Preserved ownership, justified architecture changes, synchronization boundaries, capability separation, and public API review. |
| Sprint Review | `system/protocols/review-sprint-protocol.md` | Acceptance traceability, Sprint 1 decision treatment, implementation conformance, race and stress results, and unresolved findings. |
| Deep Smoke Sprint | `system/protocols/deep-smoke-sprint-protocol.md` | Real diagnostic CLI scenarios for success, failure, panic, cancellation, parent cancellation, ignored cancellation, and collision behavior. |

## Excluded Context

| Context | Reason Excluded | Revisit If |
|---|---|---|
| Implementation execution during planning | Sprint 2 planning must settle decision changes and evidence obligations before source mutation begins. | `plan.md` validates and execution begins. |
| Smoke investigation during planning | Smoke must exercise the realized runtime and diagnostic CLI, not a predicted design. | Execution produces runnable commands. |
| Review automation during planning | Review must compare current code and evidence with accepted Sprint 1 and Sprint 2 decisions. | Execution produces reviewable changes. |
| Issue tracking mutation | Review findings remain sprint evidence until a deferred item receives separate scope and ownership. | Review identifies a finding that belongs outside the Phase 1 closure work. |
| Git mutation during planning | Planning artifacts do not authorize branch, commit, merge, or history changes. | The governed execute or merge stage owns the operation. |
| Unconditional architecture reopening | Cancellation and observers should extend the Sprint 1 owner locally unless concrete implementation evidence proves the architecture cannot preserve its invariants. | The architecture-delta selection gate is met. |
| Durable replay and exactly-once delivery | Phase 1 retains in-memory history and permits repeat delivery through replay. It does not promise persistence or exactly-once observer receipt. | Aren begins persistence and restart recovery. |
| Providers, tools, subprocesses, workflows, daemon hosting, network APIs, and remote execution | These belong to later Aren phases and must not enter through generic interfaces or speculative seams. | The owning Aren phase begins. |

## Next Artifacts

- `technical-handbook.md` distills the selected reports, accepted project reasoning, and applicable Sprint 1 evidence.
- `reasoning/*.md` resolves the five selected Sprint 2 decision areas.
- `reasoning.md` classifies relevant Sprint 1 decisions and makes the final Sprint 2 decisions.
- `plan.md` implements only those decisions and records the project-reasoning digest lock after validation.
- `review.md` and `smoke.md` run the selected protocols against the realized runtime and diagnostic CLI.
