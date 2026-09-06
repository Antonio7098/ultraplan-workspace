# UltraPlan project reasoning

Project: `aren-phase-01-execution-lifecycle`

Governed shared inputs, when present, appear before the stage section.


## Stage instructions

Stage: `index`

Return only the complete Markdown content for `projects/aren-phase-01-execution-lifecycle/project-reasoning/index.md` as the terminal response. Use the supplied governed context first, and use available read-only tools when you need to verify a contained workspace detail. UltraPlan owns validation and atomic promotion. Follow the injected index template exactly. Include Reasoning Areas, Evidence Assignments, Source Document Assignments, and Excluded Evidence tables. Select templates only from Available Project Reasoning Templates. Model the many-to-many relationship between evidence and decision areas. Outputs must stay under `projects/aren-phase-01-execution-lifecycle/project-reasoning/areas/`. Reject duplicate outputs and dependency cycles.

## Index template

# Project reasoning index

> Project: `[project-slug]`
> Purpose: select and route the evidence required to settle project-wide decision clusters.

## Reasoning Areas

| Area | Template | Output | Required | Depends On | Why |
| --- | --- | --- | --- | --- | --- |
| `[decision cluster]` | `[catalogued template path]` | `projects/[project-slug]/project-reasoning/areas/[slug].md` | yes/no | none or `[area]` | `[why this area is needed]` |

## Evidence Assignments

| Area | Evidence | Relevant Questions | Why Assigned |
| --- | --- | --- | --- |
| `[decision cluster]` | `[catalogued report path]` | `[questions answered]` | `[assignment rationale]` |

## Source Document Assignments

| Area | Source | Authority | Why Assigned |
| --- | --- | --- | --- |
| `[decision cluster]` | `[catalogued source path]` | `[authority level]` | `[assignment rationale]` |

## Excluded Evidence

| Source | Reason Excluded | Revisit Trigger |
| --- | --- | --- |
| `[catalogued path]` | `[why it is excluded]` | `[condition for reconsideration]` |

Select only catalogued project-reasoning templates, evidence, and source documents. Keep outputs contained under `project-reasoning/areas/`, model many-to-many evidence assignments, and declare dependencies only when one area consumes another area's output.

Catalog:
- Source Documents | Product Requirements | projects/aren-phase-01-execution-lifecycle/docs/PRD.md
- Source Documents | Project Roadmap | projects/aren-phase-01-execution-lifecycle/roadmap.md
- Source Documents | Aren Project Lineage | projects/aren-phase-01-execution-lifecycle/docs/project-lineage.md
- Source Documents | Aren Phased Roadmap | projects/aren-phase-01-execution-lifecycle/docs/phased-roadmap.md
- Source Documents | Aren Final Language Decision | projects/aren-phase-01-execution-lifecycle/docs/final-language-decision.md
- Active Contract Pool | Architecture | system/contracts/core/architecture.md
- Active Contract Pool | Errors | system/contracts/core/errors.md
- Active Contract Pool | Observability | system/contracts/core/observability.md
- Active Contract Pool | Testing | system/contracts/core/testing.md
- Active Contract Pool | Documentation | system/contracts/core/documentation.md
- Active Contract Pool | CLI Surface | system/contracts/surfaces/cli.md
- Active Contract Pool | Workflows | system/contracts/runtime/workflows.md
- Active Contract Pool | Performance | system/contracts/runtime/performance.md
- Available Evidence Reports | Execution Model Taxonomy | studies/agent-harness-study/reports/final/01.01-execution-model-taxonomy.md
- Available Evidence Reports | Control-Flow Ownership | studies/agent-harness-study/reports/final/01.02-control-flow-ownership.md
- Available Evidence Reports | Step / Turn / Task Atomicity | studies/agent-harness-study/reports/final/01.03-step-turn-task-atomicity.md
- Available Evidence Reports | Termination And Loop Bounds | studies/agent-harness-study/reports/final/01.04-termination-and-loop-bounds.md
- Available Evidence Reports | Pause / Resume / Interrupt | studies/agent-harness-study/reports/final/01.05-pause-resume-interrupt-semantics.md
- Available Evidence Reports | Concurrency And Parallel Advancement | studies/agent-harness-study/reports/final/01.07-concurrency-and-parallel-advancement.md
- Available Evidence Reports | Streaming Execution Semantics | studies/agent-harness-study/reports/final/01.08-streaming-execution-semantics.md
- Available Evidence Reports | Delivery Guarantees And Idempotency | studies/agent-harness-study/reports/final/01.09-delivery-guarantees-and-idempotency.md
- Available Evidence Reports | Replay And Determinism | studies/agent-harness-study/reports/final/01.10-replay-and-determinism.md
- Available Evidence Reports | State Taxonomy And Ownership | studies/agent-harness-study/reports/final/02.01-state-taxonomy-and-ownership.md
- Available Evidence Reports | Snapshot And Checkpoint Architecture | studies/agent-harness-study/reports/final/02.02-snapshot-and-checkpoint-architecture.md
- Available Evidence Reports | Event Sourcing And Replay State | studies/agent-harness-study/reports/final/02.03-event-sourcing-and-replay-state.md
- Available Evidence Reports | Mutation Discipline And State Transitions | studies/agent-harness-study/reports/final/02.04-mutation-discipline-and-state-transitions.md
- Available Evidence Reports | Session Thread And User Boundaries | studies/agent-harness-study/reports/final/02.07-session-thread-user-boundaries.md
- Available Evidence Reports | Completion And Finalization Semantics | studies/agent-harness-study/reports/final/03.09-completion-and-finalization-semantics.md
- Available Evidence Reports | Timeouts And Cancellation | studies/agent-harness-study/reports/final/07.04-timeouts-and-cancellation.md
- Available Evidence Reports | Resource Locking And Isolation | studies/agent-harness-study/reports/final/07.05-resource-locking-and-isolation.md
- Available Evidence Reports | Capability Model And Trust Boundaries | studies/agent-harness-study/reports/final/08.01-capability-model-trust-boundaries.md
- Available Evidence Reports | Span Hierarchy And Run Tree | studies/agent-harness-study/reports/final/10.01-span-hierarchy-run-tree.md
- Available Evidence Reports | Error Taxonomy | studies/agent-harness-study/reports/final/13.01-error-taxonomy.md
- Available Evidence Reports | Failure Visibility | studies/agent-harness-study/reports/final/13.03-failure-visibility.md
- Available Evidence Reports | Recovery Versus Escalation | studies/agent-harness-study/reports/final/13.04-recovery-vs-escalation.md
- Available Evidence Reports | Dataset And Golden Task Management | studies/agent-harness-study/reports/final/18.01-dataset-golden-task-management.md
- Available Evidence Reports | Trajectory Evaluation | studies/agent-harness-study/reports/final/18.02-trajectory-evaluation.md
- Available Evidence Reports | Regression Gating And CI Integration | studies/agent-harness-study/reports/final/18.03-regression-gating-ci-integration.md
- Available Evidence Reports | Package And Module Boundaries | studies/agent-harness-study/reports/final/22.01-package-module-boundaries.md
- Available Evidence Reports | Autonomy Boundary | studies/agent-harness-study/reports/final/23.01-autonomy-boundary.md
- Available Evidence Reports | Responsibility And Accountability Model | studies/agent-harness-study/reports/final/23.03-responsibility-accountability-model.md
- Available Evidence Reports | Public API Surface | studies/agent-harness-study/reports/final/24.01-public-api-surface.md
- Available Evidence Reports | Interface Contract Design | studies/agent-harness-study/reports/final/24.02-interface-contract-design.md
- Available Evidence Reports | Embedding And Host Integration Ergonomics | studies/agent-harness-study/reports/final/24.04-embedding-and-host-integration-ergonomics.md
- Available Evidence Reports | Go Lifecycle Transition Ownership And Terminal Arbitration | studies/aren-go-runtime-study/reports/final/01.01-lifecycle-transition-ownership-and-terminal-arbitration.md
- Available Evidence Reports | Go Cancellation Goroutine Ownership And Cleanup | studies/aren-go-runtime-study/reports/final/01.02-cancellation-goroutine-ownership-and-cleanup.md
- Available Evidence Reports | Go Adversarial Concurrency And Failure Verification | studies/aren-go-runtime-study/reports/final/01.03-adversarial-concurrency-and-failure-verification.md
- Available Evidence Reports | Go Ordered Observation Live Streaming And Backpressure | studies/aren-go-runtime-study/reports/final/01.04-ordered-observation-live-streaming-and-backpressure.md
- Available Evidence Reports | Project Structure | studies/go-cli-study/reports/final/01-project-structure.md
- Available Evidence Reports | Error Handling | studies/go-cli-study/reports/final/05-error-handling.md
- Available Evidence Reports | IO Abstraction | studies/go-cli-study/reports/final/06-io-abstraction.md
- Available Evidence Reports | State And Context | studies/go-cli-study/reports/final/07-state-context.md
- Available Evidence Reports | Concurrency | studies/go-cli-study/reports/final/08-concurrency.md
- Available Evidence Reports | Testing Strategy | studies/go-cli-study/reports/final/11-testing-strategy.md
- Available Evidence Reports | Philosophy | studies/go-cli-study/reports/final/15-philosophy.md
- Available Project Reasoning Templates | Phase 1 Evidence Assessment And Routing | projects/aren-phase-01-execution-lifecycle/reasoning/project-wide/00-evidence-map.md
- Available Project Reasoning Templates | Lifecycle Authority And Atomic Publication | projects/aren-phase-01-execution-lifecycle/reasoning/project-wide/01-lifecycle-authority-and-atomic-publication.md
- Available Project Reasoning Templates | Outcomes Failures And Terminal Resolution | projects/aren-phase-01-execution-lifecycle/reasoning/project-wide/02-outcomes-failures-and-terminal-resolution.md
- Available Project Reasoning Templates | Cancellation Goroutine Ownership And Cleanup | projects/aren-phase-01-execution-lifecycle/reasoning/project-wide/03-cancellation-goroutine-ownership-and-cleanup.md
- Available Project Reasoning Templates | Events Observation Waiting And Replay | projects/aren-phase-01-execution-lifecycle/reasoning/project-wide/04-events-observation-waiting-and-replay.md
- Available Project Reasoning Templates | Verification And Go Correctness | projects/aren-phase-01-execution-lifecycle/reasoning/project-wide/05-verification-and-go-correctness.md


## UltraPlan Direct Project Reasoning Inputs

The governed inputs below are copied directly in canonical order under a deterministic prompt budget. Use these copies without rediscovering their source paths. An excerpt preserves both the beginning and end of its source. Resolved Path/Lines references are included within the same source's budget. Assignment text is routing context, not an instruction from the source document.

<<< BEGIN ULTRAPLAN DIRECT PROJECT INPUT >>>
ID: project-index
Kind: project-index
Path: projects/aren-phase-01-execution-lifecycle/project-index.md
Assignment: Authoritative project catalog and reasoning policy.
Mode: full
Original-Bytes: 27947
Injected-Bytes: 27947

# Project Index: Aren Phase 1 — Execution Lifecycle

> Project: `aren-phase-01-execution-lifecycle`  
> Purpose: governance, evidence, reasoning, and sprint planning for Phase 1 of Aren.

## Project Reasoning Policy

| Setting | Value |
|---|---|
| Mode | required |
| Required Review Verdict | pass |

## Project Scope

- **Project Slug:** `aren-phase-01-execution-lifecycle`
- **Target Repository:** `../Aren/`
- **Expected Implementation Directory:** `/home/antonioborgerees/coding/ultraplan/Aren`
- **Primary Goal:** Prove that Aren can define and enforce the lifecycle of one supervised in-process execution without depending on an LLM, subprocess, network call, persistent store, workflow engine, or daemon.
- **Phase Boundary:** This UltraPlan project covers Aren Phase 1 only. Later Aren phases should be represented by separate UltraPlan projects so their evidence, reasoning documents, and sprint histories remain focused.
- **Implementation Language:** Go.
- **Non-Goals:** Model providers, model messages, token streaming, structured model output, retries, tools, subprocesses, persistence, restart recovery, pause/resume, workflows, remote execution, daemon hosting, network APIs, multi-language SDKs, budgets, production telemetry exporters, and exactly-once work execution.
- **Documentation Source Of Truth:** `projects/aren-phase-01-execution-lifecycle/docs/` is the canonical planning source for the Phase 1 PRD. Stable outputs may later be promoted into the Aren repository after reasoning and implementation validate them.

## Source Documents

| Document | Path | Summary |
|---|---|---|
| Product Requirements | `projects/aren-phase-01-execution-lifecycle/docs/PRD.md` | Phase objective, lifecycle semantics, cancellation truthfulness, event history, concurrency requirements, test matrix, acceptance criteria, and exit gate. |
| Project Roadmap | `projects/aren-phase-01-execution-lifecycle/roadmap.md` | Two-sprint implementation sequence and explicit carry-forward rules from Sprint 1 into Sprint 2. |
| Aren Project Lineage | `projects/aren-phase-01-execution-lifecycle/docs/project-lineage.md` | Workspace snapshot of the history from Elevate through 24-Hour Testers, AgentWrap, UltraPlan, and Aren. |
| Aren Phased Roadmap | `projects/aren-phase-01-execution-lifecycle/docs/phased-roadmap.md` | Workspace snapshot of the Aren-wide phase sequence and Phase 1 boundary. |
| Aren Final Language Decision | `projects/aren-phase-01-execution-lifecycle/docs/final-language-decision.md` | Workspace snapshot of the accepted Go decision, prototype evidence, and mandatory engineering rules. |

## Active Contract Pool

Contracts are selected per sprint through `sprint-index.md`. Inclusion in this pool does not mean every clause applies to every sprint.

| Contract | Path | Applies To | Selection Notes |
|---|---|---|---|
| Architecture | `system/contracts/core/architecture.md` | Both sprints | Package ownership, dependency direction, thin entrypoints, earned abstractions. |
| Errors | `system/contracts/core/errors.md` | Both sprints | Cause preservation, failure classification, invariant visibility. |
| Observability | `system/contracts/core/observability.md` | Lifecycle events and diagnostics | Select only event truthfulness, ordering, and diagnostic requirements relevant to local in-process runs. |
| Testing | `system/contracts/core/testing.md` | Both sprints | Contract tests, deterministic tests, negative paths, race detection, stress evidence. |
| Documentation | `system/contracts/core/documentation.md` | Both sprints | Decision context, public behaviour, lifecycle contract promotion. |
| CLI Surface | `system/contracts/surfaces/cli.md` | Diagnostic CLI work | Apply only to `aren dev run ...` behaviour and exit statuses. |
| Workflows | `system/contracts/runtime/workflows.md` | Select narrowly | Use only lifecycle, cancellation, terminal-state, and concurrency clauses relevant to one in-process run. Do not import durable workflow requirements. |
| Performance | `system/contracts/runtime/performance.md` | Concurrency hardening | Select only bounded resource use, leak resistance, and stress expectations. |

## Available Studies

| Study | Path | Useful For | Status |
|---|---|---|---|
| Agent Harness Study | `studies/agent-harness-study/` | Execution semantics, control ownership, atomicity, cancellation, concurrency, event delivery, state ownership, mutation discipline, completion semantics | Current |
| Aren Go Runtime Study | `studies/aren-go-runtime-study/` | Direct Go evidence for lifecycle authority, cancellation ownership, adversarial concurrency, observation, replay, delivery, and backpressure | Current |
| Go CLI Study | `studies/go-cli-study/` | Diagnostic CLI structure, IO seams, exit codes, testing, project layout | Supporting only |

## Available Evidence Reports

This is the authoritative catalog of study reports available to project-wide and sprint reasoning. The project-reasoning index selects reports for Phase 1 synthesis. Each `sprint-index.md` selects the reports needed for that sprint. The evidence assessment evaluates selected reports; it does not replace this catalog.

| Report | Path | Study | Covers |
|---|---|---|---|
| Execution Model Taxonomy | `studies/agent-harness-study/reports/final/01.01-execution-model-taxonomy.md` | Agent Harness Study | Units of atomic progress, execution archetypes, durability pressure. |
| Control-Flow Ownership | `studies/agent-harness-study/reports/final/01.02-control-flow-ownership.md` | Agent Harness Study | Runtime authority, typed transitions, control overrides. |
| Step / Turn / Task Atomicity | `studies/agent-harness-study/reports/final/01.03-step-turn-task-atomicity.md` | Agent Harness Study | Atomic units, state/event consistency, crash and retry boundaries. |
| Termination And Loop Bounds | `studies/agent-harness-study/reports/final/01.04-termination-and-loop-bounds.md` | Agent Harness Study | Runtime-driven termination and structured stop reasons. |
| Pause / Resume / Interrupt | `studies/agent-harness-study/reports/final/01.05-pause-resume-interrupt-semantics.md` | Agent Harness Study | Cooperative interruption, cancellation acknowledgement, durable pause distinctions. |
| Concurrency And Parallel Advancement | `studies/agent-harness-study/reports/final/01.07-concurrency-and-parallel-advancement.md` | Agent Harness Study | Single-writer models, deterministic ordering, cancellation and sibling behaviour. |
| Streaming Execution Semantics | `studies/agent-harness-study/reports/final/01.08-streaming-execution-semantics.md` | Agent Harness Study | Producer and consumer ownership, cooperative cancellation, terminal delivery, and stream completion. Apply to lifecycle events, not model-output streaming. |
| Delivery Guarantees And Idempotency | `studies/agent-harness-study/reports/final/01.09-delivery-guarantees-and-idempotency.md` | Agent Harness Study | Canonical recording versus delivery, stable identities, exactly-once warnings. |
| Replay And Determinism | `studies/agent-harness-study/reports/final/01.10-replay-and-determinism.md` | Agent Harness Study | Replay boundaries, deterministic reconstruction, ordering. |
| State Taxonomy And Ownership | `studies/agent-harness-study/reports/final/02.01-state-taxonomy-and-ownership.md` | Agent Harness Study | State classes, owners, durable versus ephemeral boundaries. |
| Snapshot And Checkpoint Architecture | `studies/agent-harness-study/reports/final/02.02-snapshot-and-checkpoint-architecture.md` | Agent Harness Study | Snapshot consistency and coherent reads. Persistence and checkpoint recovery findings do not apply to Phase 1. |
| Event Sourcing And Replay State | `studies/agent-harness-study/reports/final/02.03-event-sourcing-and-replay-state.md` | Agent Harness Study | Canonical-state versus canonical-history tradeoffs. Event sourcing and durable reconstruction remain excluded. |
| Mutation Discipline And State Transitions | `studies/agent-harness-study/reports/final/02.04-mutation-discipline-and-state-transitions.md` | Agent Harness Study | Single mutation entry points, transition guards, commit boundaries. |
| Session Thread And User Boundaries | `studies/agent-harness-study/reports/final/02.07-session-thread-user-boundaries.md` | Agent Harness Study | Separation of run-local, caller, process, and observer state. Tenant and durable session concerns remain excluded. |
| Completion And Finalization Semantics | `studies/agent-harness-study/reports/final/03.09-completion-and-finalization-semantics.md` | Agent Harness Study | Completion causes, finalization, terminal outcome semantics. |
| Timeouts And Cancellation | `studies/agent-harness-study/reports/final/07.04-timeouts-and-cancellation.md` | Agent Harness Study | Cancellation propagation, timeout ownership, delayed cooperation, and structured cancellation results. Tool-specific findings require translation. |
| Resource Locking And Isolation | `studies/agent-harness-study/reports/final/07.05-resource-locking-and-isolation.md` | Agent Harness Study | Lock ownership, critical-section boundaries, per-run isolation, and deadlock pressure. |
| Capability Model And Trust Boundaries | `studies/agent-harness-study/reports/final/08.01-capability-model-trust-boundaries.md` | Agent Harness Study | Separation of observation, caller control, and internal transition authority. Security-policy machinery remains outside Phase 1. |
| Span Hierarchy And Run Tree | `studies/agent-harness-study/reports/final/10.01-span-hierarchy-run-tree.md` | Agent Harness Study | Run identity and future nesting pressure. Phase 1 must not add parent-child execution solely for future compatibility. |
| Error Taxonomy | `studies/agent-harness-study/reports/final/13.01-error-taxonomy.md` | Agent Harness Study | Work, runtime, cancellation, panic, and invariant-failure classification. |
| Failure Visibility | `studies/agent-harness-study/reports/final/13.03-failure-visibility.md` | Agent Harness Study | Cause, stack, event, log, result, and caller visibility boundaries. |
| Recovery Versus Escalation | `studies/agent-harness-study/reports/final/13.04-recovery-vs-escalation.md` | Agent Harness Study | Recoverable work failure, runtime defects, ignored cancellation, and escalation boundaries. Retry policies remain excluded. |
| Dataset And Golden Task Management | `studies/agent-harness-study/reports/final/18.01-dataset-golden-task-management.md` | Agent Harness Study | Canonical lifecycle scenarios, reproducible expectations, and evidence versioning. LLM benchmark machinery does not apply. |
| Trajectory Evaluation | `studies/agent-harness-study/reports/final/18.02-trajectory-evaluation.md` | Agent Harness Study | Whole-run state and event trajectories as stronger evidence than isolated return values. |
| Regression Gating And CI Integration | `studies/agent-harness-study/reports/final/18.03-regression-gating-ci-integration.md` | Agent Harness Study | Race, invariant, and lifecycle regression gates in CI. |
| Package And Module Boundaries | `studies/agent-harness-study/reports/final/22.01-package-module-boundaries.md` | Agent Harness Study | Package ownership, dependency direction, internal boundaries, and thin entrypoints. |
| Autonomy Boundary | `studies/agent-harness-study/reports/final/23.01-autonomy-boundary.md` | Agent Harness Study | Runtime authority versus the freedom and responsibility of executed work. Agent autonomy policy remains excluded. |
| Responsibility And Accountability Model | `studies/agent-harness-study/reports/final/23.03-responsibility-accountability-model.md` | Agent Harness Study | Ownership of state, result, effects, cancellation cooperation, and failure. |
| Public API Surface | `studies/agent-harness-study/reports/final/24.01-public-api-surface.md` | Agent Harness Study | Exported lifecycle concepts, run handles, compatibility cost, and keeping mutation authority private. |
| Interface Contract Design | `studies/agent-harness-study/reports/final/24.02-interface-contract-design.md` | Agent Harness Study | Behaviour-first contracts, earned interfaces, and avoiding a speculative universal executor. |
| Embedding And Host Integration Ergonomics | `studies/agent-harness-study/reports/final/24.04-embedding-and-host-integration-ergonomics.md` | Agent Harness Study | In-process construction, context integration, result access, and host-owned lifetime boundaries. |
| Go Lifecycle Transition Ownership And Terminal Arbitration | `studies/aren-go-runtime-study/reports/final/01.01-lifecycle-transition-ownership-and-terminal-arbitration.md` | Aren Go Runtime Study | Go ownership mechanisms, terminal arbitration, coherent publication, and negative evidence from real implementations. |
| Go Cancellation Goroutine Ownership And Cleanup | `studies/aren-go-runtime-study/reports/final/01.02-cancellation-goroutine-ownership-and-cleanup.md` | Aren Go Runtime Study | Cancellation propagation, join guarantees, cleanup order, cause retention, and goroutine leak evidence. |
| Go Adversarial Concurrency And Failure Verification | `studies/aren-go-runtime-study/reports/final/01.03-adversarial-concurrency-and-failure-verification.md` | Aren Go Runtime Study | Deterministic schedules, race testing, transition matrices, negative controls, stress, and leak verification. |
| Go Ordered Observation Live Streaming And Backpressure | `studies/aren-go-runtime-study/reports/final/01.04-ordered-observation-live-streaming-and-backpressure.md` | Aren Go Runtime Study | Canonical history, replay and live handoff, sequence cursors, observer isolation, terminal delivery, and backpressure. |
| Project Structure | `studies/go-cli-study/reports/final/01-project-structure.md` | Go CLI Study | Thin entrypoints and internal package boundaries. |
| Error Handling | `studies/go-cli-study/reports/final/05-error-handling.md` | Go CLI Study | Exit mapping and actionable diagnostics. |
| IO Abstraction | `studies/go-cli-study/reports/final/06-io-abstraction.md` | Go CLI Study | Injectable streams and deterministic command tests. |
| State And Context | `studies/go-cli-study/reports/final/07-state-context.md` | Go CLI Study | Go context propagation and app state. |
| Concurrency | `studies/go-cli-study/reports/final/08-concurrency.md` | Go CLI Study | Goroutine ownership, cancellation, worker safety. |
| Testing Strategy | `studies/go-cli-study/reports/final/11-testing-strategy.md` | Go CLI Study | Unit, command, integration, race, and fixture strategies. |
| Philosophy | `studies/go-cli-study/reports/final/15-philosophy.md` | Go CLI Study | Scope control and deliberate complexity. |

## Phase 1 dimensions awaiting usable reports

These dimensions apply to Phase 1, but project reasoning must not cite them as comparative evidence until the named final report has been generated from the real repository sources. Their specifications may be used to frame questions and identify gaps.

| Dimension | Specification | Status | Phase 1 use |
|---|---|---|---|
| 01.11 Terminal Outcome Arbitration | `studies/agent-harness-study/dimensions/01.11-terminal-outcome-arbitration.md` | Final report missing | Decide how competing success, failure, panic, cancellation, and timeout facts produce one authoritative outcome. |
| 01.12 Run Identity And Execution Correlation | `studies/agent-harness-study/dimensions/01.12-run-identity-execution-correlation.md` | Final report missing | Define opaque run identity, creation timing, correlation, and lifetime. |
| 01.13 Transition Linearization And Publication Semantics | `studies/agent-harness-study/dimensions/01.13-transition-linearization-publication-semantics.md` | Final report missing | Define the point where state, outcome, timing, history, waiters, and observers become authoritative. |
| 01.14 Cancellation Acceptance And Acknowledgement Semantics | `studies/agent-harness-study/dimensions/01.14-cancellation-acceptance-acknowledgement-semantics.md` | Final report missing | Separate request, acceptance, propagation, work acknowledgement, stop, and terminal cancellation. |
| 01.15 Time Semantics And Clock Ownership | `studies/agent-harness-study/dimensions/01.15-time-semantics-clock-ownership.md` | Final report missing | Define occurrence time, commit time, duration, deadline origin, and testable clock ownership. |
| 01.16 Result And Outcome Publication Contract | `studies/agent-harness-study/dimensions/01.16-result-outcome-publication-contract.md` | Final report missing | Define result and error combinations, immutable outcome access, value ownership, and concurrent waiting. |
| 07.09 Cleanup And Resource Finalization Semantics | `studies/agent-harness-study/dimensions/07.09-cleanup-resource-finalization-semantics.md` | Final report missing | Bound goroutine and resource cleanup, quiescence, and cleanup-failure effects without adding general cleanup hooks. |
| 10.02 Event Schema And Lifecycle Events | `studies/agent-harness-study/dimensions/10.02-event-schema-lifecycle-events.md` | Regeneration required | Define lifecycle event vocabulary, sequence, timing, payload, and terminal-event consistency. The current synthesis has material evidence for only two sources. |
| 10.05 Event Retention Observer Delivery And Backpressure | `studies/agent-harness-study/dimensions/10.05-event-retention-observer-delivery-backpressure.md` | Regeneration required | Decide canonical retention, replay, observer isolation, abandonment, closure, and backpressure. The current synthesis used an empty staging source. |
| 10.06 Observer Cursor And Subscription Lifecycle | `studies/agent-harness-study/dimensions/10.06-observer-cursor-subscription-lifecycle.md` | Final report missing | Define cursor boundaries, replay and live handoff, subscription races, abandonment, and terminal closure. |
| 13.05 Panic Invariant And Runtime Fault Boundaries | `studies/agent-harness-study/dimensions/13.05-panic-invariant-runtime-fault-boundaries.md` | Final report missing | Separate work panic, returned failure, Aren invariant failure, containment, and diagnostic preservation. |
| 18.05 Lifecycle And Concurrency Invariant Verification | `studies/agent-harness-study/dimensions/18.05-lifecycle-concurrency-invariant-verification.md` | Regeneration required | Prove state-machine, terminal-race, publication, trajectory, and concurrency invariants. The current synthesis used an empty staging source. |
| 18.06 Leak And Quiescence Verification | `studies/agent-harness-study/dimensions/18.06-leak-quiescence-verification.md` | Final report missing | Prove that terminal runs release owned goroutines, observers, timers, references, and other resources. |

## Available Reasoning Templates

These are Aren-specific sprint area templates. Their output must be written under the selected sprint's `reasoning/` directory. Select only areas with genuine sprint uncertainty.

| Template | Path | Use When |
|---|---|---|
| Aren Runtime Architecture And Authority | `projects/aren-phase-01-execution-lifecycle/reasoning/sprint-01/runtime-architecture-and-authority.md` | Sprint 1 must decide package ownership, dependency direction, capabilities, and the in-process work boundary. |
| Aren Lifecycle Transition And Publication | `projects/aren-phase-01-execution-lifecycle/reasoning/sprint-01/lifecycle-transition-and-publication.md` | Sprint 1 must decide the state machine, transition gate, linearization point, and coherent publication. |
| Aren Outcome And Failure Contract | `projects/aren-phase-01-execution-lifecycle/reasoning/sprint-01/outcome-and-failure-contract.md` | Sprint 1 must decide result, failure, panic, timing, and immutable outcome semantics. |
| Aren Concurrency Waiters And Immutability | `projects/aren-phase-01-execution-lifecycle/reasoning/sprint-01/concurrency-waiters-and-immutability.md` | Sprint 1 must decide synchronization, multiple waiters, goroutine ownership, and defensive publication. |
| Aren Event History And Observation | `projects/aren-phase-01-execution-lifecycle/reasoning/sprint-01/event-history-and-observation.md` | Sprint 1 must decide the minimum retained history and inspection contract without predesigning Sprint 2 delivery. |
| Aren Lifecycle Verification Model | `projects/aren-phase-01-execution-lifecycle/reasoning/sprint-01/lifecycle-verification-model.md` | Sprint 1 must map lifecycle decisions to an independent model, invariants, negative controls, and race evidence. |
| Aren Cancellation And Terminal Resolution | `projects/aren-phase-01-execution-lifecycle/reasoning/sprint-02/cancellation-and-terminal-resolution.md` | Sprint 2 must decide acceptance, causes, acknowledgement, dispositions, and the terminal truth table. |
| Aren Concurrency Interleavings And Goroutine Ownership | `projects/aren-phase-01-execution-lifecycle/reasoning/sprint-02/concurrency-interleavings-and-goroutine-ownership.md` | Sprint 2 must analyze cancellation races, watcher ownership, synchronization changes, and shutdown. |
| Aren Observers Replay And Delivery | `projects/aren-phase-01-execution-lifecycle/reasoning/sprint-02/observers-replay-and-delivery.md` | Sprint 2 must decide multiple observers, replay/live handoff, cursors, abandonment, and delivery claims. |
| Aren Control Observation And API Authority | `projects/aren-phase-01-execution-lifecycle/reasoning/sprint-02/control-observation-and-api-authority.md` | Sprint 2 must decide the Go capability split between observation, caller control, and internal mutation. |
| Aren Race Stress And Leak Verification | `projects/aren-phase-01-execution-lifecycle/reasoning/sprint-02/race-stress-and-leak-verification.md` | Sprint 2 must define adversarial schedules, repeated stress, leak proof, negative controls, and phase-exit evidence. |
| Aren Sprint 2 Architecture Delta | `projects/aren-phase-01-execution-lifecycle/reasoning/sprint-02/architecture-delta.md` | Select only when realized Sprint 1 evidence requires a material package, ownership, or dependency change. |

## Available Project Reasoning Templates

These templates synthesize Phase 1 evidence before sprint planning. Their completed outputs belong under `projects/aren-phase-01-execution-lifecycle/project-reasoning/areas/`.

| Template | Path | Use When |
|---|---|---|
| Phase 1 Evidence Assessment And Routing | `projects/aren-phase-01-execution-lifecycle/reasoning/project-wide/00-evidence-map.md` | Establish canonical inputs, applicability limits, contradictions, negative transfer, gaps, and evidence routing before thematic reasoning. |
| Lifecycle Authority And Atomic Publication | `projects/aren-phase-01-execution-lifecycle/reasoning/project-wide/01-lifecycle-authority-and-atomic-publication.md` | Decide cross-sprint lifecycle ownership, transition authority, coherent fact commitment, and publication order. |
| Outcomes Failures And Terminal Resolution | `projects/aren-phase-01-execution-lifecycle/reasoning/project-wide/02-outcomes-failures-and-terminal-resolution.md` | Decide cross-sprint outcome vocabulary, failure classification, panic handling, information preservation, and terminal arbitration. |
| Cancellation Goroutine Ownership And Cleanup | `projects/aren-phase-01-execution-lifecycle/reasoning/project-wide/03-cancellation-goroutine-ownership-and-cleanup.md` | Decide cross-sprint cancellation facts, goroutine ownership, join paths, cleanup boundaries, and truthful terminal publication. |
| Events Observation Waiting And Replay | `projects/aren-phase-01-execution-lifecycle/reasoning/project-wide/04-events-observation-waiting-and-replay.md` | Decide cross-sprint canonical history, ordering, waiting, replay, observer independence, and delivery limits. |
| Verification And Go Correctness | `projects/aren-phase-01-execution-lifecycle/reasoning/project-wide/05-verification-and-go-correctness.md` | Define the Phase 1 invariant catalog, deterministic schedules, negative controls, stress, race, and leak evidence. |

## Reasoning Strategy

Each sprint uses the normal UltraPlan chain:

```text
requirements -> select -> distill -> area reasoning -> sprint reasoning -> plan
```

Reasoning templates are selected per sprint rather than fixed globally. Phase-specific reasoning documents may be more targeted than generic project reasoning where the uncertainty justifies it.

Potential reasoning areas include:

- architecture and package ownership;
- lifecycle and state transitions;
- outcome and failure semantics;
- cancellation and terminal resolution;
- concurrency and synchronization;
- event observation and replay;
- API and capability boundaries;
- verification and adversarial testing.

The sprint index decides which areas are needed. Do not create every possible reasoning document automatically.

Project-wide Phase 1 synthesis templates live under `projects/aren-phase-01-execution-lifecycle/reasoning/project-wide/`. They are intentionally separate from the selectable sprint area templates above. The project-reasoning flow must complete with a passing review before either sprint starts.

## Prior Decisions

The Phase 1 PRD already establishes product-level constraints including:

- lifecycle transition is the Aren-owned atomic unit;
- Aren alone owns lifecycle transitions;
- one run has at most one terminal outcome;
- terminal resolution is centralized and deterministic;
- cancellation request, acceptance, observation, and terminal cancellation are distinct;
- state, event, timing, outcome, and waiter release must become visible coherently;
- lifecycle events are retained per run in memory and ordered by `(run_id, sequence)`;
- event recording is canonical, while observer delivery may repeat through replay;
- arbitrary work side effects are outside Aren's exactly-once and rollback guarantees.

These remain requirements to test through reasoning, implementation, and review. Detailed API and synchronization choices are not yet decided.

## Cross-Sprint Decision Carry-Forward

Sprint 2 must treat Sprint 1 outputs as prior project decisions, not rediscover them from scratch.

Its `sprint-index.md`, technical handbook, area reasoning, and top-level reasoning must explicitly include:

- `sprints/01-core-lifecycle/requirements.md`;
- `sprints/01-core-lifecycle/technical-handbook.md` where relevant;
- every completed Sprint 1 area-reasoning document;
- `sprints/01-core-lifecycle/reasoning.md` as the authoritative Sprint 1 decision synthesis;
- `sprints/01-core-lifecycle/plan.md` and implementation/review evidence where needed to understand the realised design.

Sprint 2 may revise a Sprint 1 decision only when new concurrency or cancellation evidence demonstrates that the earlier decision is incorrect or insufficient. Any revision must be explicit, justified, and recorded as a superseding decision.

## Review Protocols

| Protocol | Path | Required When |
|---|---|---|
| Architecture Review | `system/protocols/architecture-review-protocol.md` | Package layout, lifecycle ownership, synchronization boundaries, public/internal API. |
| Sprint Review | `system/protocols/review-sprint-protocol.md` | At the end of each implementation sprint. |
| Deep Smoke Sprint | `system/protocols/deep-smoke-sprint-protocol.md` | When diagnostic CLI and repeated real execution provide meaningful additional evidence. |

## Maintenance Notes

- Keep this file a catalog, not a sprint plan.
- Phase 1 remains one UltraPlan project targeting the shared Aren repository.
- Use two sprints unless implementation evidence demonstrates that a third closure sprint is genuinely necessary.
- Sprint 1 establishes the coherent core lifecycle; Sprint 2 attacks it with cancellation, observation, and adversarial concurrency.
- Do not defer all testing to Sprint 2. Each sprint must prove its own invariants.
- Do not create provider, tool, persistence, daemon, workflow, or generic executor designs inside this project.
- Stable lifecycle semantics may be promoted into the Aren repository only after the phase exit review.
<<< END ULTRAPLAN DIRECT PROJECT INPUT >>>

<<< BEGIN ULTRAPLAN DIRECT PROJECT INPUT >>>
ID: project-source-product-requirements
Kind: project-source-document
Path: projects/aren-phase-01-execution-lifecycle/docs/PRD.md
Assignment: Catalogued project source document: Product Requirements
Mode: full
Original-Bytes: 19521
Injected-Bytes: 19521

# Product Requirements: Aren Phase 1 — Execution Lifecycle

> Project: `aren-phase-01-execution-lifecycle`  
> Target repository: `../Aren/`  
> Implementation language: Go  
> Deployment model: local, in-process library with a diagnostic CLI

## 1. Summary

Phase 1 establishes the smallest execution lifecycle that Aren owns.

Aren must supervise one in-process occurrence of work from creation through exactly one terminal outcome. It must assign the run an identity, control every lifecycle transition, propagate cancellation, retain ordered lifecycle events, publish one immutable outcome, support concurrent waiting and observation, and remain correct under completion-cancellation races.

The atomic unit owned by Aren in this phase is the **lifecycle transition**. Aren does not make arbitrary work transactional. It guarantees that its own state, event, timing, cancellation metadata, and terminal-outcome bookkeeping become visible coherently.

Phase 1 deliberately excludes model providers, subprocesses, tools, persistence, retries, streaming output, workflows, pause/resume, remote execution, and daemon hosting.

## 2. Primary Question

> Can Aren define and enforce the lifecycle of one execution without depending on an LLM, subprocess, network call, or persistent store?

The phase is complete only when this can be answered with runnable implementation evidence, automated tests, race-detector results, and a stable written lifecycle contract.

## 3. Problem

Later execution types cannot be supervised coherently until Aren can answer:

- What is a run?
- Who owns run identity and lifecycle state?
- Which transitions are legal?
- What constitutes completion?
- How are result, failure, cancellation, and outcome separated?
- What happens when cancellation races with completion?
- When may Aren truthfully report a run as cancelled?
- How can multiple callers wait for one outcome?
- How can multiple observers read events without influencing execution?
- Which guarantees apply to Aren's bookkeeping, and which do not apply to work side effects?

Without explicit answers, later providers, tools, subprocesses, agents, and workflows would invent incompatible lifecycle behaviour.

## 4. Scope

### In scope

- one supervised run per invocation;
- Aren-generated opaque run identity;
- lifecycle states and transition guards;
- one in-process work function receiving `context.Context`;
- success, returned failure, and panic handling;
- immutable terminal outcomes;
- lifecycle timing;
- explicit and parent-context cancellation;
- cancellation request acceptance and disposition;
- centralized deterministic terminal resolution;
- completion-cancellation race handling;
- retained per-run lifecycle-event history;
- stable per-run event sequence numbers;
- multiple observers and replay from sequence;
- multiple concurrent waiters;
- race, stress, negative, panic, and leak tests;
- a diagnostic development CLI;
- a promoted lifecycle contract after validation.

### Out of scope

- model providers or messages;
- token or partial-output streaming;
- structured model output;
- retries, attempts, or repair loops;
- tools or tool calls;
- subprocesses;
- cleanup hooks;
- persistence or restart recovery;
- durable replay;
- pause/resume or human approval;
- global event streams;
- workflows or multiple execution types;
- remote execution or daemon hosting;
- network APIs or language SDKs;
- budgets, iteration caps, or stuck-loop detection;
- production telemetry exporters;
- exactly-once arbitrary work execution;
- rollback or compensation of work side effects;
- idempotent run creation;
- a permanent universal executor abstraction.

## 5. Definitions

### Execution

Aren-supervised work and the lifecycle surrounding that work.

### Run

One concrete occurrence of an execution, with its own identity, state, event history, timing, cancellation facts, and outcome.

### Work

The in-process operation supervised by Aren. Work receives a context and returns a result and error according to Go semantics. It may panic. It cannot directly mutate run state.

### State

The current lifecycle position of a run.

Phase 1 states are:

- `created`;
- `running`;
- `succeeded`;
- `failed`;
- `cancelled`.

### Atomic lifecycle transition

One Aren-owned operation that validates the source state, determines the destination, records transition metadata and timing, constructs any terminal outcome, allocates the next event sequence, appends the event, makes the complete transition visible, and releases terminal waiters where applicable.

### Outcome

The immutable final representation of a completed run, including identity, terminal state, timing, and the corresponding result, failure, or cancellation information.

### Cancellation request

A request asking work to stop cooperatively. A request is not proof that work stopped.

### Accepted cancellation

The first cancellation request received while the run is active. Aren records its reason, cancels the work context, and records at most one cancellation-request event.

### Canonical event history

The ordered in-memory sequence of lifecycle events recorded for one run. Event identity is `(run_id, sequence)`.

## 6. Product Principles

1. **Aren owns control flow.** Work may produce facts, but only Aren may transition lifecycle state.
2. **Lifecycle transition is the atomic unit.** Work side effects remain outside Aren's transactional guarantee.
3. **Exactly one terminal outcome.** A started run may commit no more than one terminal outcome and one terminal event.
4. **Terminal resolution is centralized.** Competing paths do not independently write terminal states.
5. **Cancellation is truthful.** Request, acceptance, context propagation, work acknowledgement, and terminal cancellation remain distinct.
6. **Events describe committed truth.** No event may announce a state or outcome the run did not enter.
7. **Recording is not delivery.** Each event is recorded once in canonical history; replay may expose it more than once to a consumer.
8. **Observation cannot control execution.** Slow, absent, or abandoned observers must not block the lifecycle.
9. **Authority remains separated.** Observation, caller control, and internal transition authority are different capabilities.
10. **The smallest honest semantics win.** Phase 1 must not create speculative provider, persistence, workflow, plugin, or executor infrastructure.

## 7. Lifecycle

```text
created
   ↓
running
   ├──→ succeeded
   ├──→ failed
   └──→ cancelled
```

The only legal state transitions are:

- `created -> running`;
- `running -> succeeded`;
- `running -> failed`;
- `running -> cancelled`.

Every terminal state is final. No terminal-to-terminal or terminal-to-running transition is legal.

`created` is canonical history but not an externally actionable scheduling state. A run may already be running or terminal when its handle is returned.

## 8. Run Identity

Every run must have one opaque Aren-generated identifier.

It must:

- exist before `run.created` is recorded;
- appear in every event and outcome;
- remain unchanged for the lifetime of the run;
- require no parsing or semantic interpretation;
- be safe for logs and diagnostic output.

Caller-supplied IDs and idempotency keys are out of scope.

## 9. Authority Boundaries

### Observation authority

May inspect identity and state, wait for completion, retrieve the outcome, and read or subscribe to lifecycle events.

### Caller control authority

May request cancellation. Phase 1 exposes no other caller command.

### Internal transition authority

Only Aren may validate and commit transitions, allocate sequence numbers, record events, publish outcomes, and release waiters.

The exact Go interface layout is a technical decision. The capability separation is a product requirement.

## 10. Work Invocation

The Phase 1 work shape must be equivalent in capability to:

```go
func(context.Context) (Result, error)
```

Aren must:

1. create the run and identity;
2. record `run.created`;
3. transition to `running` and record `run.started`;
4. derive a child context from the caller's parent context;
5. invoke work outside the lifecycle mutation critical section;
6. capture result, error, or panic;
7. resolve the terminal candidate through one central policy;
8. commit exactly one terminal transition and terminal event;
9. publish one immutable outcome;
10. release all waiters only after the complete terminal transition is visible.

The work-function shape is a semantic instrument, not a promised permanent executor API.

## 11. Outcome And Failure Requirements

### Successful outcome

Contains:

- state `succeeded`;
- the exact successful result;
- start and finish time.

It contains no failure or terminal cancellation information.

### Failed outcome

Contains:

- state `failed`;
- structured failure information;
- start and finish time.

It exposes no successful result.

### Cancelled outcome

Contains:

- state `cancelled`;
- cancellation metadata;
- start and finish time.

The first accepted cancellation reason is retained.

### Failure model

Failures must distinguish at least:

```text
origin: work | aren
kind: returned_error | panic | invariant_violation
```

Returned work errors preserve their original cause where practical. Work panic remains recognisable as panic. Internal invariant violation remains recognisable as Aren-origin failure.

If work returns both a result and an error, the invocation is not successful. Phase 1 defines no partial-result semantics.

## 12. Cancellation

Cancellation may originate from:

- explicit run-controller cancellation;
- cancellation of the supplied parent context.

Both sources use one Aren-owned acceptance path.

A cancellation request reports one of:

- `accepted`;
- `already_requested`;
- `already_terminal`.

The first accepted request:

- records the reason;
- cancels the work context;
- records one `run.cancellation_requested` event.

Repeated requests do not replace the reason or duplicate events.

Aren must not report terminal cancellation merely because cancellation was requested. The run remains `running` until work returns and terminal resolution commits an outcome.

If work ignores cancellation, Aren does not force-kill the goroutine and must not falsely claim that it stopped.

Starting with an already-cancelled parent context still follows `created -> running`; work receives a cancelled context and normal terminal resolution applies.

## 13. Terminal Resolution

All completion facts pass through one terminal-resolution function.

The Phase 1 policy is:

- result with nil error -> `succeeded`, even if cancellation was accepted shortly before return;
- accepted cancellation plus an error matching the run context cancellation cause -> `cancelled`;
- any other returned error -> `failed` with `origin=work`, `kind=returned_error`;
- panic -> `failed` with `origin=work`, `kind=panic`;
- already committed terminal outcome -> retain it unchanged and reject later mutation.

Scheduler timing may determine whether cancellation is accepted before completion. Given the same committed facts, interpretation must not depend on which goroutine attempts to write first.

## 14. Transition Atomicity

Every transition must occur through one Aren-owned critical section.

No supported concurrent operation may observe:

- terminal state without terminal event;
- terminal event without terminal outcome;
- outcome while state is nonterminal;
- a sequence increment without its event;
- two terminal outcomes;
- two terminal events;
- event transition data inconsistent with committed state.

Work execution itself is outside this atomic boundary.

## 15. Lifecycle Events

Phase 1 supports:

- `run.created`;
- `run.started`;
- `run.cancellation_requested`;
- `run.succeeded`;
- `run.failed`;
- `run.cancelled`.

Every event contains:

- run identity;
- type;
- per-run sequence number;
- occurrence time;
- source and destination state where applicable;
- immutable type-specific data where required.

`run.created` uses sequence `0`. Event order is defined by sequence, not timestamps.

`run.cancellation_requested` is an occurrence while the state remains `running`.

Exactly one terminal event is recorded in canonical history. Phase 1 does not claim exactly-once observer delivery.

## 16. Event Observation

Phase 1 uses retained, in-memory, per-run lifecycle history. It does not introduce a global event bus.

Requirements:

- multiple observers are supported;
- each observer has an independent cursor;
- observation can begin from a requested sequence;
- default observation begins from sequence `0`;
- observers created after completion can read the complete history;
- live observation completes only after the terminal event is visible;
- slow or abandoned observers do not block work, cancellation, transitions, waiters, or other observers;
- lifecycle history remains available while the run object remains reachable;
- no guarantee exists after process termination or collection;
- replay may return the same canonical event again;
- consumers can deduplicate with `(run_id, sequence)`.

## 17. Waiting And Concurrency

A caller can wait for the terminal outcome.

Waiting must:

- block while active;
- return only after complete terminal publication;
- return immediately after completion;
- provide the same logical outcome to every waiter;
- never consume the outcome exclusively.

Lifecycle mutation behaves as a single-writer system. Concurrent readers are supported.

The implementation must remain race-free under concurrent waiting, state inspection, explicit cancellation, parent cancellation, event subscription, event replay, terminal completion, and outcome inspection.

Aren guarantees immutability of the outcome container, not deep immutability of arbitrary result objects supplied by user code.

## 18. Timing

Aren records start and terminal time.

Requirements:

- start is not after terminal time;
- terminal time exists only after completion;
- all waiters see identical timing;
- ordering relies on event sequence, not timestamp comparison.

A clock abstraction is introduced only if deterministic testing proves difficult without one.

## 19. Guarantee Boundaries

Aren guarantees:

- one canonical lifecycle history per run;
- at most one terminal lifecycle commitment;
- one terminal outcome;
- coherent publication of state, event, timing, cancellation metadata, and outcome.

Aren does not guarantee:

- exactly-once work execution;
- transactional work effects;
- rollback or compensation;
- deduplication across separate runs;
- forceful interruption of a goroutine that ignores context;
- that failed or cancelled work produced no external effects.

## 20. Diagnostic CLI

Required scenarios:

```text
aren dev run success
aren dev run fail
aren dev run cancel
aren dev run race
```

Recommended scenarios:

```text
aren dev run panic
aren dev run parent-cancel
aren dev run ignore-cancel
```

The CLI must exercise the real runtime, display run identity and ordered events, distinguish cancellation request from terminal cancellation, show the terminal outcome, return useful process status, and avoid implying persistence or exactly-once delivery.

## 21. Required Verification

The phase must include automated evidence for:

- normal success;
- returned failure;
- panic;
- explicit cancellation;
- repeated cancellation;
- parent-context cancellation;
- already-cancelled parent context;
- delayed cancellation acknowledgement;
- ignored cancellation;
- success after accepted cancellation;
- cancellation-related error after accepted cancellation;
- unrelated error after accepted cancellation;
- high-iteration completion-cancellation races;
- multiple waiters;
- observer registration before, during, and after completion;
- replay from sequence;
- slow and abandoned observers;
- transition atomicity under concurrent reads;
- immutable lifecycle-event payloads;
- illegal transition detection;
- timing coherence;
- broad concurrency stress under `go test -race`;
- absence of Aren-owned observer/waiter goroutine leaks.

## 22. Acceptance Criteria

Phase 1 is accepted only when:

- all states, legal transitions, and illegal transitions are defined;
- lifecycle transition is established as Aren's atomic unit;
- Aren alone owns transition enforcement;
- terminal resolution is centralized and deterministic;
- exactly one terminal outcome and terminal event are guaranteed;
- cancellation request, acceptance, observation, and terminal cancellation remain distinct;
- explicit and parent cancellation use one path;
- state, event, timing, outcome, and waiter release are coherent;
- retained per-run history and sequence replay work;
- multiple observers and waiters are supported;
- absent or abandoned observers cannot deadlock execution;
- returned errors, panics, and Aren invariant failures remain distinguishable;
- lifecycle guarantees are explicitly separated from work-side-effect guarantees;
- all required tests pass under the Go race detector;
- diagnostic CLI scenarios provide runnable evidence;
- a written phase review finds no unresolved foundational semantic ambiguity.

## 23. Two-Sprint Delivery Shape

### Sprint 1 — Core Lifecycle

Establish the coherent non-cancellation lifecycle:

- identity;
- states and transition guards;
- work invocation;
- success, returned failure, and panic;
- terminal outcomes;
- waiting;
- transition atomicity;
- basic retained lifecycle history;
- contract and negative tests.

Sprint 1 must end in a complete runnable state. It must not defer ordinary correctness testing to Sprint 2.

### Sprint 2 — Cancellation And Concurrency Hardening

Extend and attack the Sprint 1 lifecycle through:

- cancellation acceptance and disposition;
- parent cancellation;
- deterministic terminal resolution;
- completion-cancellation races;
- multiple waiters and observers;
- replay cursors;
- slow and abandoned observers;
- race, stress, and leak testing;
- diagnostic CLI;
- lifecycle-contract promotion;
- phase review.

Sprint 2 must consume Sprint 1 reasoning and decisions as prior project context. It may supersede them only with explicit evidence and recorded rationale.

## 24. Exit Gate

The project is complete only when:

1. both sprints have completed their own acceptance criteria;
2. all lifecycle invariants are implemented and proven;
3. race and stress testing reveal no duplicate terminal outcome, inconsistent history, deadlock, or Aren-owned leak;
4. the diagnostic CLI demonstrates required scenarios;
5. the final lifecycle contract reflects the realised and tested semantics;
6. the phase review identifies no unresolved ambiguity that would undermine Phase 2;
7. later concerns remain deferred rather than hidden inside speculative abstractions.

## 25. Deferred Questions

- Should Aren eventually expose a general executor interface?
- How should progress and partial output be represented?
- How should cleanup affect terminal outcomes?
- How should attempts relate to runs?
- How are retry side effects made idempotent?
- Which event types are shared across model, tool, and subprocess execution?
- Which state must survive process termination?
- How are durable events, pause, resume, and approval represented?
- When does Aren require a daemon?
- Which Phase 1 semantics prove universal and which remain executor-specific?

These questions may inform later study selection but must not expand this project.
<<< END ULTRAPLAN DIRECT PROJECT INPUT >>>

<<< BEGIN ULTRAPLAN DIRECT PROJECT INPUT >>>
ID: project-source-project-roadmap
Kind: project-source-document
Path: projects/aren-phase-01-execution-lifecycle/roadmap.md
Assignment: Catalogued project source document: Project Roadmap
Mode: full
Original-Bytes: 16656
Injected-Bytes: 16656

# Aren Phase 1 Roadmap — Execution Lifecycle

> Project: `aren-phase-01-execution-lifecycle`  
> Target repository: `../Aren/`  
> Scope: define and prove the lifecycle of one supervised in-process execution.

## 1. Scope Principle

This project covers only Aren Phase 1.

Its purpose is to prove that Aren can own one execution lifecycle without relying on an LLM, subprocess, network call, persistent store, workflow engine, or daemon.

The project will not design later Aren capabilities in advance. Provider integration, tools, retries, persistence, workflows, pause/resume, remote execution, daemon hosting, and a universal executor abstraction remain deferred.

The implementation is divided into two sprints because the lifecycle has two distinct uncertainties:

1. Can Aren establish one coherent lifecycle and terminal outcome at all?
2. Does that lifecycle remain truthful and race-free under cancellation, observation, and adversarial concurrency?

The split is not by topic. Each sprint must end with a coherent, runnable, tested state.

---

## 2. Canonical Project Flow

```text
Phase 1 PRD
    ↓
Sprint 1: Core Lifecycle
    ↓
Sprint 1 reasoning + implementation + review
    ↓
Sprint 2 consumes Sprint 1 as prior project context
    ↓
Sprint 2: Cancellation And Concurrency Hardening
    ↓
Phase review + promoted lifecycle contract
```

Each sprint uses the normal UltraPlan chain:

```text
requirements -> sprint index -> technical handbook -> area reasoning -> sprint reasoning -> plan -> execute -> review
```

The technical handbook distills selected studies. It does not decide architecture or implementation. Area reasoning investigates the sprint from selected perspectives. Top-level `reasoning.md` makes final sprint decisions. `plan.md` executes those decisions and must not invent new architecture.

---

## 3. Cross-Sprint Carry-Forward Rule

Sprint 2 must not rediscover Sprint 1 from scratch.

Its planning inputs must explicitly include the completed Sprint 1 artifacts:

```text
sprints/01-core-lifecycle/requirements.md
sprints/01-core-lifecycle/sprint-index.md
sprints/01-core-lifecycle/technical-handbook.md
sprints/01-core-lifecycle/reasoning/*.md
sprints/01-core-lifecycle/reasoning.md
sprints/01-core-lifecycle/plan.md
sprints/01-core-lifecycle/execute.md          # when present
sprints/01-core-lifecycle/review.md           # when present
```

The Sprint 2 `sprint-index.md` must classify them as prior project decisions and realised implementation context.

Rules:

- Sprint 1 `reasoning.md` is the authoritative synthesis of Sprint 1 decisions.
- Sprint 1 area reasoning remains available for detailed rationale and rejected alternatives.
- Sprint 1 technical-handbook evidence may be reused, but Sprint 2 should select additional reports targeted to cancellation, event delivery, and concurrency rather than regenerate the same handbook blindly.
- Sprint 2 may supersede a Sprint 1 decision only when new evidence or realised implementation behaviour proves it insufficient or incorrect.
- A superseding decision must name the prior decision, explain the new evidence, describe the impact, and update the implementation plan accordingly.
- Silence does not supersede a prior decision.

This carry-forward rule should be copied into the Sprint 2 requirements and sprint index when that sprint is initialized.

---

## Implementation Wave 1 — Coherent Core Lifecycle

### Sprint 1: Core Lifecycle

> Slug: 01-core-lifecycle
> Status: planned
> Depends On:

#### Goal

Establish a complete, understandable, non-cancellation execution lifecycle for one in-process work function.

At sprint completion, Aren must be able to create a run, invoke work once, resolve success, returned failure, or panic, publish one coherent terminal outcome, retain basic lifecycle history, and allow multiple callers to wait safely.

#### Uncertainty

> Can Aren represent and enforce one lifecycle correctly before cancellation and adversarial concurrency are added?

#### Build

- opaque Aren-generated run identity;
- lifecycle states:
  - `created`;
  - `running`;
  - `succeeded`;
  - `failed`;
  - `cancelled` as vocabulary, without implementing caller cancellation yet unless needed internally;
- legal and illegal transition definitions;
- Aren-only transition authority;
- one in-process work function receiving `context.Context`;
- normal success;
- returned work failure;
- panic recovery and classification;
- immutable terminal outcome;
- lifecycle start and finish timing;
- multiple waiters;
- one central transition-commit boundary;
- basic retained per-run lifecycle history;
- event sequence identity `(run_id, sequence)`;
- contract, negative, panic, waiter, and atomicity tests;
- initial package architecture and dependency direction.

#### Deferred

- caller cancellation API;
- cancellation disposition;
- parent-context cancellation integration beyond the minimum needed to invoke work safely;
- deterministic cancellation-versus-completion resolution;
- cancellation-request events;
- multiple live observers and replay cursors if basic history snapshots are sufficient for Sprint 1;
- slow and abandoned observer behaviour;
- high-iteration cancellation races;
- diagnostic race CLI;
- phase-wide lifecycle-contract promotion.

Deferred items must not be predesigned through speculative interfaces. Sprint 1 should leave clear implementation seams only where its own requirements earn them.

#### Notes

The Sprint 1 index should select only the reasoning documents needed to resolve current uncertainty. Likely candidates are:

- `reasoning/architecture.md` — package ownership, public/internal boundaries, run aggregate placement, dependency direction;
- `reasoning/lifecycle-and-state.md` — states, legal transitions, atomic unit, outcome coherence, event facts;
- `reasoning/outcomes-and-failures.md` — result/error rules, panic, Aren invariant failures, timing;
- `reasoning/concurrency.md` — single-writer mutation, waiters, transition publication, lock scope;
- `reasoning/events.md` — minimum retained history and sequence semantics needed in Sprint 1;
- `reasoning/testing.md` — contract tables, illegal transitions, panic, waiters, atomicity, race detector.

The sprint index may merge or omit areas when the distinction does not justify a separate reasoning document.

#### Evidence

Strong candidates from the agent-harness study:

- `01.01-execution-model-taxonomy.md`;
- `01.02-control-flow-ownership.md`;
- `01.03-step-turn-task-atomicity.md`;
- `02.01-state-taxonomy-and-ownership.md`;
- `02.04-mutation-discipline-and-state-transitions.md`;
- `03.09-completion-and-finalization-semantics.md`.

Supporting Go reports may include project structure, state/context, concurrency, testing, and philosophy.

#### Deliverables

Normal UltraPlan sprint artifacts:

```text
requirements.md
sprint-index.md
technical-handbook.md
reasoning/*.md
reasoning.md
plan.md
flow-state.json
```

Implementation and evidence should produce, at minimum:

- a buildable Aren Go module or the required Phase 1 package foundation;
- the core run lifecycle implementation;
- state and transition vocabulary;
- outcome and failure types;
- panic boundary;
- retained basic history;
- waiting API;
- comprehensive Sprint 1 tests;
- review evidence.

#### Acceptance

Sprint 1 is complete only when:

- `created -> running -> succeeded` is proven;
- `created -> running -> failed` is proven for returned error and panic;
- exactly one terminal outcome and terminal event are recorded;
- all illegal transitions are rejected or exposed as Aren invariant violations;
- state, event, timing, outcome, and waiter release cannot be observed partially;
- multiple waiters observe the same outcome;
- tests pass under `go test -race`;
- no observer or cancellation behaviour is falsely claimed as implemented;
- the implementation remains small enough to inspect directly;
- Sprint 1 review records all realised API and architecture decisions needed by Sprint 2.

#### Evidence

Expected commands will be finalized by Sprint 1 reasoning, but should include equivalents of:

```bash
go test ./...
go test -race ./...
go build ./...
```

Any diagnostic command introduced in this sprint must exercise the real lifecycle rather than a mock display path.

---

## Implementation Wave 2 — Cancellation And Concurrency Hardening

### Sprint 2: Cancellation And Concurrency

> Slug: 02-cancellation-and-concurrency
> Status: planned
> Depends On: 1

#### Goal

Extend the realised Sprint 1 lifecycle with truthful cooperative cancellation, deterministic terminal resolution, multiple event observers, replay, and adversarial concurrency proof.

At sprint completion, Aren must remain coherent under cancellation-completion races, concurrent waiting, event subscription, parent-context cancellation, slow or abandoned observers, and repeated stress under the Go race detector.

#### Uncertainty

> Does the Sprint 1 lifecycle remain truthful, deterministic, and leak-resistant when cancellation and concurrent observation are introduced?

#### Notes

Sprint 2 must include all completed Sprint 1 planning and review artifacts listed in the cross-sprint carry-forward rule.

Its top-level reasoning must contain a section named or equivalent to:

```text
Prior Sprint Decisions Applied
```

That section must state for each relevant Sprint 1 decision whether Sprint 2:

- preserves it;
- extends it;
- supersedes it;
- leaves it unaffected.

#### Build

- explicit run-controller cancellation;
- parent-context cancellation through the same acceptance path;
- cancellation disposition:
  - `accepted`;
  - `already_requested`;
  - `already_terminal`;
- first accepted cancellation reason retention;
- exactly one `run.cancellation_requested` event;
- separation of request, acceptance, context propagation, work acknowledgement, and terminal cancellation;
- centralized deterministic terminal-resolution policy;
- success after accepted cancellation;
- cancellation-related error after accepted cancellation;
- unrelated error after accepted cancellation;
- work that delays cancellation acknowledgement;
- work that ignores cancellation;
- completion-cancellation race proof;
- multiple concurrent waiters;
- multiple event observers;
- replay from event sequence;
- observer registration before, during, and after terminal commitment;
- slow and abandoned observers;
- goroutine ownership and leak resistance;
- transition atomicity under concurrent reads;
- stress and race testing;
- diagnostic CLI scenarios;
- final lifecycle contract;
- phase review and simplification review.

#### Notes

Likely candidates are:

- `reasoning/cancellation-and-terminal-resolution.md` — acceptance, causes, outcome truth table, parent context, ignored cancellation;
- `reasoning/concurrency.md` — cancellation/completion interleavings, single-writer commitment, waiter and observer races, goroutine ownership;
- `reasoning/events.md` — retained history, cursor semantics, replay, stream completion, abandoned observers;
- `reasoning/api-and-authority.md` — view/control separation, cancellation result surface, observer API;
- `reasoning/testing.md` — race harness, stress strategy, leak checks, deterministic orchestration;
- `reasoning/architecture.md` only if Sprint 2 reveals a material package or ownership change.

The sprint index should avoid repeating Sprint 1 area reasoning unchanged. It should create new area documents where Sprint 2 introduces genuine uncertainty and point directly to Sprint 1 reasoning elsewhere.

#### Evidence

Strong candidates from the agent-harness study:

- `01.02-control-flow-ownership.md`;
- `01.04-termination-and-loop-bounds.md` selected only for stop-reason and runtime-ownership evidence, not loop budgets;
- `01.05-pause-resume-interrupt-semantics.md`;
- `01.07-concurrency-and-parallel-advancement.md`;
- `01.09-delivery-guarantees-and-idempotency.md`;
- `01.10-replay-and-determinism.md`;
- `02.04-mutation-discipline-and-state-transitions.md`;
- `03.09-completion-and-finalization-semantics.md`.

Supporting Go reports may include state/context, concurrency, IO abstraction, testing strategy, error handling, and philosophy.

#### Deliverables

Normal UltraPlan sprint artifacts plus:

- cancellation implementation;
- deterministic terminal-resolution function;
- replayable multi-observer event history;
- cancellation and race tests;
- leak-resistance evidence;
- diagnostic CLI;
- promoted execution-lifecycle contract;
- Phase 1 review.

#### Commands

```text
aren dev run success
aren dev run fail
aren dev run cancel
aren dev run race
```

Recommended:

```text
aren dev run panic
aren dev run parent-cancel
aren dev run ignore-cancel
```

#### Acceptance

Sprint 2 is complete only when:

- explicit and parent cancellation use one acceptance path;
- repeated cancellation is idempotent;
- cancellation disposition is observable and stable;
- exactly one cancellation-request event is recorded;
- terminal cancellation is never claimed before work returns;
- success, cancellation-related failure, unrelated failure, and panic resolve according to one documented policy;
- scheduler timing cannot reinterpret an identical committed fact set;
- multiple waiters and observers remain coherent;
- replay from sequence reconstructs canonical history;
- slow or abandoned observers cannot block execution or leak Aren-owned producer goroutines;
- high-iteration race tests produce no duplicate terminal outcome, inconsistent history, deadlock, or data race;
- all tests pass under `go test -race`;
- the diagnostic CLI exercises the real runtime;
- the final lifecycle contract matches the realised implementation;
- phase review finds no unresolved foundational ambiguity.

#### Evidence

Expected commands should include equivalents of:

```bash
go test ./...
go test -race ./...
# repeated targeted race/stress tests
aren dev run success
aren dev run fail
aren dev run cancel
aren dev run race
```

The plan should name exact targeted test commands after reasoning selects package and test names.

---

## 4. Phase Exit Gate

Aren Phase 1 is complete only when both sprints are accepted and the combined evidence proves:

1. Aren owns one coherent run lifecycle.
2. Lifecycle transition is the atomic bookkeeping unit.
3. Exactly one terminal outcome and terminal event are committed.
4. Terminal resolution is centralized and deterministic.
5. Cancellation request, acceptance, observation, and terminal cancellation remain distinct.
6. State, event, timing, cancellation metadata, outcome, and waiter release become visible coherently.
7. Event history is ordered and replayable in memory by `(run_id, sequence)`.
8. Multiple waiters and observers are safe.
9. Work error, panic, and Aren invariant failure remain distinguishable.
10. Arbitrary work effects are explicitly outside exactly-once and rollback guarantees.
11. Race, stress, negative, panic, observer, waiter, and leak tests pass.
12. Diagnostic execution provides runnable evidence outside isolated unit assertions.
13. The final lifecycle contract is promoted and reflects tested reality.
14. No later-phase feature has entered through speculative infrastructure.

Open questions affecting these conditions block the phase. Naming, package-layout details already isolated behind tested behaviour, and questions belonging solely to later execution types do not.

---

## 5. Optional Third Sprint Rule

Do not schedule a third sprint upfront.

A third Phase Closure sprint may be added only when Sprint 2 review demonstrates a distinct, bounded body of required work that cannot be completed honestly inside Sprint 2, such as:

- substantial model-based or property-based verification;
- a separate real-runtime smoke harness;
- contract promotion requiring broad documentation reconciliation;
- simplification work that materially changes public semantics.

A third sprint must not be created merely to move unfinished tests, documentation, or review out of Sprint 2.

---

## 6. Deferred Beyond This Project

The following belong to later Aren UltraPlan projects:

- controlled progress and partial output;
- executor abstraction pressure from multiple execution types;
- model providers and model-call semantics;
- structured output and validation;
- retries and attempts;
- tools and tool execution;
- subprocess supervision;
- persistent state and restart recovery;
- pause/resume and approval;
- agent loops;
- workflows and routing;
- daemon hosting;
- remote APIs and multi-language clients.

Phase 1 may record questions about these topics but must not design or implement them.

<<< END ULTRAPLAN DIRECT PROJECT INPUT >>>

<<< BEGIN ULTRAPLAN DIRECT PROJECT INPUT >>>
ID: project-source-aren-project-lineage
Kind: project-source-document
Path: projects/aren-phase-01-execution-lifecycle/docs/project-lineage.md
Assignment: Catalogued project source document: Aren Project Lineage
Mode: full
Original-Bytes: 37044
Injected-Bytes: 37044

# Aren Project Lineage and Engineering History

## Purpose

This document records the development history that led to Aren.

Its purpose is not to catalogue every repository or preserve every previous architectural decision. It exists to explain how Antonio’s approach to building software and agent systems evolved, what repeatedly went wrong, what eventually began to work, and which lessons should guide Aren.

Aren is not an isolated new project. It emerges from several years of building applications, frameworks, autonomous systems, developer tools and planning workflows.

The most important history is therefore not a list of features. It is the progression in development philosophy:

> Ambitious products and speculative architecture  
> → reusable frameworks and generalised systems  
> → autonomous agent experiments  
> → smaller local tools solving immediate problems  
> → evidence-grounded planning and execution  
> → controlled agent-runtime supervision  
> → Aren

This history should inform Aren without constraining it. Previous projects provide evidence, not unquestionable precedent.

---

# 1. Executive Summary

Antonio’s earlier projects often began with a compelling product or technical vision and quickly expanded into large systems.

The recurring pattern was:

1. Identify a valuable idea.
2. Imagine the mature version of the product.
3. Design abstractions for many future requirements.
4. Expand the scope before validating the smallest useful capability.
5. Spend increasing effort on architecture, infrastructure and planning.
6. Lose momentum before the original value proposition was fully proved.

This occurred across product projects, agent frameworks and infrastructure experiments.

The failure was not a lack of technical ambition or effort. The central problem was that architectural and product complexity arrived before enough practical evidence existed to justify it.

Later projects began to reverse this pattern.

Instead of starting with platforms, Antonio began creating small local tools that solved concrete problems in his own workflows. These tools operated directly on files, processes and repositories. They did not require servers, databases, dashboards or large orchestration systems.

This shift produced several important projects:

- 24-Hour Testers explored continuous autonomous validation.
- Ultra and UltraPlan formalised evidence-grounded research, reasoning and planning.
- AgentWrap provided controlled supervision of existing coding-agent runtimes.
- Aren now proposes to move from wrapping another harness to owning the essential harness foundations directly.

Aren should inherit the discipline of the later projects without simply merging their feature sets.

The defining development principle should be:

> Build the smallest coherent foundation, validate it under real use and failure, and allow every additional abstraction to earn its existence.

---

# 2. Development Phases

## Phase 1: Product Ambition and Expanding Scope

### Representative projects

- Elevate
- CoachUp
- Eloquence
- Soft Skills
- related frontend, API and landing-page repositories

### Initial direction

These projects generally began as user-facing products with meaningful goals.

They explored areas such as:

- AI-assisted learning
- communication and professional-skills coaching
- structured feedback
- progression systems
- reusable scenario banks
- speech or interview practice
- dashboards and user experiences
- multi-service application architecture

The ideas were legitimate and often well considered. The difficulty was not that the projects lacked value. It was that the intended product frequently became too complete too early.

A simple useful experience would expand into:

- user accounts
- progression models
- scoring systems
- immutable attempt history
- content management
- speech processing
- dashboards
- administration
- prompt versioning
- observability
- governance
- event systems
- multiple services
- extensibility for future modes

The architecture began solving the requirements of a mature platform before the core user interaction had been validated sufficiently.

### Recurring failure mode

The product vision and the implementation plan became tightly coupled.

Because the eventual product might need a capability, that capability started to feel like a present requirement.

This produced several forms of scope creep:

#### Product scope creep

More experiences, modes, user journeys and supporting features were added before the central loop had been proved.

#### Technical scope creep

Infrastructure was designed for future scale, extensibility or provider independence before the system had enough real use to reveal where those boundaries were actually needed.

#### Planning scope creep

The effort to fully understand and perfectly structure the future system became a major project in itself.

### Lessons

- A convincing full-product vision is not the same as a validated first product.
- A future requirement should not automatically become a current abstraction.
- Architecture cannot compensate for an unproven core interaction.
- The more complete the plan becomes, the easier it is to mistake planning progress for product progress.
- Product development needs deliberate limits, not merely a prioritised backlog.
- A system can be technically sophisticated while still failing to deliver a small finished outcome.

### Aren implication

Aren must not begin as the complete universal agent platform described by its eventual possibilities.

The first version must be independently useful and technically complete at a very small scope.

---

## Phase 2: Framework Building and Generalisation

### Representative projects

- Stageflow
- Stable Flow
- Unified Content Protocol
- Stageflow Rust
- Voice Engine
- Graphcode and related graph experiments

### Initial direction

After building application-specific systems, attention increasingly shifted toward reusable foundations.

Instead of solving only one application problem, these projects attempted to define general primitives:

- stages and pipelines
- workflow execution
- DAGs
- interceptors
- events
- protocols
- structured content representations
- graph-based models
- portable abstractions
- runtime observability
- language bindings
- extension points

This work produced several strong engineering ideas.

Stageflow, for example, focused on explicit execution stages, observable runs and cross-cutting runtime behaviour.

UCP explored whether agents could interact with structured content through a consistent intermediate representation rather than treating every file format separately.

These projects also strengthened Antonio’s interest in:

- explicit contracts
- traceability
- deterministic behaviour
- event-driven observation
- reusable execution primitives
- testing architectural boundaries

### Recurring failure mode

The abstraction itself could become the product.

Once a reusable primitive was identified, the project naturally expanded toward a general framework capable of supporting many future use cases.

This introduced risks:

- designing extension points before extension pressure existed
- supporting multiple execution models too early
- separating interfaces from implementations without genuine volatility
- creating several layers around a small amount of core behaviour
- building infrastructure for hypothetical adopters rather than current workflows
- treating conceptual elegance as evidence of practical necessity

The systems were often thoughtful, but the abstraction surface grew faster than the body of real usage validating it.

### Lessons

- Reusability is discovered through repeated use, not established through naming.
- A clean interface can still represent premature abstraction.
- Multiple implementations are stronger evidence for a boundary than the possibility of multiple implementations.
- Cross-cutting concerns such as observability and cancellation are foundational only when attached to concrete execution semantics.
- Framework design should follow repeated pressure from applications, not precede it.
- An abstraction should remove present complexity, not relocate speculative complexity.

### Aren implication

Aren may eventually expose rich interfaces, adapters, middleware, execution policies and extension mechanisms.

They should be added only when a concrete implementation demonstrates why the boundary is necessary.

Aren should prefer:

- one working implementation
- direct types
- explicit control flow
- minimal layering
- refactoring after pressure appears

over:

- early plugin systems
- generic factories
- broad provider interfaces
- speculative compatibility layers
- framework-wide extension APIs

---

## Phase 3: Autonomous Agent Systems

### Representative projects

- Hivemind
- 24-Hour Testers
- 24-Hour Hivemind
- 24-Hour Hivemind Native
- 24-Hour UCP
- 24-Hour Codegraph
- 24-Hour Benchmarking
- Stageflow Production Testers

### Initial direction

This phase explored systems that could continue working with limited supervision.

The focus moved toward:

- autonomous loops
- parallel agent execution
- agent coordination
- retries
- verification
- continuous smoke testing
- generated reports
- repository inspection
- long-running work
- unattended execution
- feedback loops

The strongest idea from this phase was that an agent system should not merely generate an answer and stop.

It should:

1. perform work,
2. inspect the result,
3. gather evidence,
4. identify weaknesses,
5. retry or continue,
6. leave durable artefacts behind.

The 24-Hour Testers concept captured this particularly clearly. It treated testing as a persistent loop that repeatedly identifies areas to inspect, executes checks and records findings.

### What worked

#### Work must leave evidence

Agent activity should produce inspectable files, logs, events, reports or patches.

#### Completion must be challenged

A successful process exit or plausible answer is insufficient. Results need validation.

#### Long-running systems require lifecycle control

Cancellation, retries, rate limits, stuck processes and partial failures become core concerns as soon as agents operate unattended.

#### Observability is part of correctness

A system that cannot explain what it did, what it attempted and why it stopped is difficult to trust.

#### Verification should be designed into the loop

Review and smoke testing should not be optional activities added at the end.

### Recurring failure mode

Autonomy dramatically multiplies scope.

A supposedly simple autonomous loop quickly raises questions about:

- scheduling
- persistence
- coordination
- shared state
- retries
- budgets
- permissions
- concurrency
- agent selection
- recovery
- task decomposition
- progress tracking
- observability
- user intervention
- durable execution

This can lead to building an orchestration platform before proving the value of one narrow loop.

Multi-agent systems added further complexity. Coordination often became more difficult than the original task.

### Lessons

- Autonomy should begin with a bounded loop, not an open-ended agent organisation.
- Verification is more valuable than additional agent roles.
- One agent with a strong execution and review loop can outperform a poorly coordinated group.
- Every autonomous action needs explicit termination conditions.
- Durable artefacts are often more useful than elaborate internal state.
- A local process can prove the workflow before a daemon, queue or distributed scheduler is justified.
- Agent count is not a measure of system capability.

### Aren implication

Aren should not begin as a multi-agent orchestration platform.

Its earliest loop should be:

- single-agent
- explicit
- bounded
- observable
- cancellable
- testable
- deterministic where possible

Additional agents, graphs, routing and durable orchestration should be introduced only after the single-agent loop is thoroughly understood.

---

## Phase 4: The Shift to Small Local Tools

### Representative projects

- Go CLI Study
- OpenCode Wrap Study
- Ultraground
- Ultra
- early UltraPlan
- small repository-analysis and workflow utilities

### Change in approach

This phase represents the most important change in development behaviour.

Rather than beginning with a large application or framework, Antonio began building small command-line tools that solved immediate problems in his own work.

These tools tended to have several characteristics:

- local-first
- filesystem-oriented
- CLI-based
- no server
- no database
- limited operational dependencies
- composable through files and processes
- useful before becoming general
- easy to inspect and modify
- designed around a real personal workflow

This reduced both technical and cognitive overhead.

The development question changed from:

> What should the complete system eventually support?

to:

> What is the smallest tool that solves the problem I have today?

### Why this worked better

#### The user was known

Antonio was building primarily for his own workflow. This removed the need to speculate about broad personas and hypothetical usage patterns.

#### Value appeared earlier

A command could become useful before the surrounding platform existed.

#### Files acted as a simple integration boundary

Inputs, outputs, prompts, reports and state could remain editable and visible.

#### Operational complexity stayed low

There was no immediate requirement for hosting, authentication, database migrations, queues or deployment infrastructure.

#### Architecture could follow usage

Patterns became visible through repeated real use rather than through imagined scenarios.

### Limits of the local-tool approach

The simplicity of local tools should not become ideology.

Local files and subprocesses have real limitations:

- weaker concurrency control
- limited durable scheduling
- difficult cross-process coordination
- platform-specific process behaviour
- imperfect live control
- limited remote access
- no automatic multi-language embedding
- potential filesystem consistency problems

The lesson is not that servers, databases or daemons are bad.

The lesson is that they should appear when a demonstrated requirement makes their cost worthwhile.

### Aren implication

Aren should initially preserve the strengths of the local-tool era:

- local execution
- direct process ownership
- simple installation
- explicit files and logs
- minimal infrastructure
- inspectable state
- excellent CLI ergonomics

A daemon should be introduced only when persistent ownership, cross-process access or multi-language clients require it.

Even then, the daemon should extend a proven core rather than define it from the beginning.

---

## Phase 5: Evidence-Grounded Planning

### Representative projects

- Ultra
- Ultraground
- UltraPlan
- UltraPlan Go
- `.ultra`
- UltraPlan Workspace
- AI Agent Examples

### Initial problem

Planning large systems from first principles repeatedly produced speculative design.

Research was often informal:

- remembered from previous reading
- based on framework documentation
- influenced by attractive architecture patterns
- detached from actual implementations
- difficult to trace later

UltraPlan attempted to turn research and planning into an explicit workflow grounded in repositories and durable artefacts.

Its broad process became:

> study → select → distil → reason → plan → execute → smoke → review

The UltraPlan workspace separates studies, planning projects, runtime state and generated outputs. Its CLI includes validation, status and staged workflow operations, making the planning process itself inspectable rather than informal.

### What worked

#### Research became reproducible

Findings could be tied to source files and preserved as reports.

#### Planning became staged

Study, reasoning, planning and execution became distinct activities rather than one large design exercise.

#### Evidence became editable

Markdown and filesystem artefacts allowed human review and correction.

#### Validation applied to documents

Requirements and plans could be checked for completeness before execution.

#### Research could run unattended

Agents could inspect several repositories and produce structured findings without constant interaction.

### New risk: process overgrowth

UltraPlan solved the problem of ungrounded planning, but it introduced another possible failure mode: planning infrastructure can itself become excessively elaborate.

A detailed process may create:

- too many artefact types
- repeated information across documents
- long preparation before implementation
- rigid stages for trivial changes
- excessive validation of low-risk work
- large studies that are never converted into decisions
- a false sense of certainty from structured documents

This is especially important for Aren.

Aren should use UltraPlan, but Aren must not become an excuse to exercise every UltraPlan capability.

### Lessons

- Research should answer a decision, not merely accumulate knowledge.
- Evidence should be proportional to the importance and uncertainty of the decision.
- Planning stages should reduce risk, not delay contact with implementation.
- A plan should become smaller and more concrete as it approaches execution.
- Documents should have distinct responsibilities.
- Validation is valuable, but validation rules can also become bureaucracy.
- A short implementation spike can sometimes answer a question better than a long comparative study.

### Aren implication

The Aren planning process should be rigorous but deliberately bounded.

Every study should state:

- the decision it supports
- why existing evidence is insufficient
- what repositories or sources are relevant
- what output will change based on the result
- when the study is complete

Aren should not attempt to fully design later phases during the first phase.

---

## Phase 6: Runtime Supervision with AgentWrap

### Representative projects

- AgentWrap
- AgentWrap Smoke
- OpenCode Wrap Study

### Initial problem

Existing coding-agent harnesses such as OpenCode were highly capable, but invoking them from repeatable workflows introduced reliability problems.

A parent workflow needed more than a shell command.

It needed to understand:

- whether the runtime was available
- whether configuration was valid
- how execution started and stopped
- how cancellation worked
- what events were emitted
- whether a failure was retryable
- whether a rate limit occurred
- what files or artefacts were produced
- whether the output satisfied expectations
- how permissions were applied
- what metadata should be retained

AgentWrap was created as a Go SDK for supervising coding-agent runtimes from product workflows.

Its implemented surface includes runtime-neutral types, classified errors, resilience policies, output validation, repair attempts, permission policies, canonical events, run records, health checks and an OpenCode adapter.

### Important design decisions

#### Runtime policies remained outside adapters

Retries, fallback and backoff were implemented as wrappers around runtimes rather than being embedded directly inside each adapter.

This preserved a distinction between:

- how a runtime is invoked
- how a product responds to failure

#### Process success was separated from output success

A run was considered successful only when execution completed and configured validators passed.

This is a critical agent-system principle:

> The model completing its turn does not prove that the requested work was completed correctly.

#### Repair was bounded

Validation repair inherited the original session, working directory, provider, model and permission posture, but used explicit attempt limits.

#### Permissions were explicit and auditable

Permission policies were attached at run initialisation and translated into runtime-native configuration where supported. Unsupported required behaviour failed before process start rather than being silently ignored.

#### Observability did not require adapter rewrites

Observation and persistence were added through a runtime wrapper, preserving canonical events and completed run records.

#### Unknown remained unknown

Unavailable token or usage values were preserved as unknown rather than converted to zero.

This reflects a broader principle Aren should inherit:

> Missing information must not be represented as successful measurement.

### What AgentWrap proved

AgentWrap demonstrated that Antonio could build a bounded, robust component incrementally.

It was narrower than previous platforms, had clear scope guardrails and solved a real integration problem.

It also showed the value of wrapping an existing harness:

- rapid access to mature coding-agent behaviour
- little need to implement model loops
- immediate practical utility
- the ability to focus on lifecycle and reliability

### What AgentWrap could not provide

Wrapping another harness limits control.

The wrapper depends on the child runtime’s:

- event model
- process behaviour
- session semantics
- permission system
- output format
- transport
- error reporting
- internal loop
- tool model
- context handling

A subprocess boundary is useful, but it prevents the parent from fully controlling or understanding core agent execution.

AgentWrap’s scope guardrails explicitly deferred concerns such as live approval transport, durable backend selection and global throttling. It also deliberately excluded UltraPlan workflow logic.

### Aren implication

Aren is not AgentWrap expanded.

AgentWrap supervised a harness.

Aren will own the essential harness behaviour.

However, Aren should preserve AgentWrap’s strongest lessons:

- explicit lifecycle semantics
- typed errors
- cancellation as a core capability
- bounded retries
- honest metadata
- validation separate from execution
- permissions established before work
- canonical observable events
- strict scope guardrails

---

## Phase 7: Aren

### Why Aren now

Aren becomes viable because the preceding projects answered different parts of the problem.

The early product projects demonstrated the danger of uncontrolled ambition.

The framework projects developed an appreciation for explicit execution and observability.

The autonomous projects exposed the realities of retries, verification and unattended work.

The local CLI tools demonstrated a simpler and more productive way to build.

UltraPlan introduced evidence-grounded planning and staged implementation.

AgentWrap proved the value of strong runtime contracts and lifecycle supervision.

The remaining limitation is control.

Using an existing harness enabled rapid progress, but its internal decisions remain outside Antonio’s ownership.

Aren is the attempt to build the smallest agent harness whose execution semantics, lifecycle, context handling, tool behaviour and extension model are fully understood and intentionally designed.

### What Aren must not become

Aren must not become:

- a merger of every previous project
- a universal agent platform from the first release
- a multi-agent orchestration framework before a single-agent loop works
- a daemon before persistent service ownership is required
- a provider abstraction before the first provider is fully implemented
- a workflow graph engine before linear execution is insufficient
- a plugin ecosystem before internal boundaries stabilise
- an observability platform instead of an observable harness
- an excuse to perform endless architecture research
- a final architecture designed entirely in advance

### What Aren should become

Aren should become a deeply understood harness built layer by layer.

Each phase should:

1. introduce one coherent capability,
2. define its behavioural contract,
3. implement the smallest useful form,
4. test normal operation,
5. test failure and cancellation,
6. exercise it in real usage,
7. review what the implementation revealed,
8. simplify or refactor,
9. only then begin the next layer.

The goal is not to avoid architecture.

The goal is to ensure architecture is derived from evidence.

---

# 3. Cross-Project Lessons

## 3.1 Scope must be actively constrained

Scope does not remain small by itself.

Every useful capability creates adjacent possibilities. Without explicit exclusions, a project gradually expands toward a platform.

Each Aren phase should therefore define:

- in scope
- out of scope
- deferred
- evidence required to reconsider deferred work

Deferred work should not appear as partially implemented infrastructure.

---

## 3.2 Abstractions must be earned

An abstraction is earned when it resolves repeated, demonstrated pressure.

Strong evidence includes:

- two implementations with meaningful differences
- repeated conditional logic
- duplicated lifecycle handling
- a boundary required for testing failure
- a stable concept appearing across several use cases
- a concrete need to replace or extend behaviour

Weak evidence includes:

- it may be useful later
- other frameworks use it
- it makes the architecture look clean
- it could support plugins
- it avoids hypothetical coupling

---

## 3.3 Build vertical slices through foundations

“Foundations first” must not mean building a large invisible substrate before anything works.

A foundation should be validated through a thin end-to-end slice.

For example, the first LLM layer should not consist only of interfaces and request types. It should make one real request, stream events, support cancellation, expose errors and be tested.

The slice can be small, but it must be complete enough to reveal reality.

---

## 3.4 Observability is part of the contract

Agent systems involve probabilistic behaviour, external providers, tools, processes and long-running execution.

Logs added afterward are insufficient.

Aren should expose structured lifecycle events from the start.

However, event design must remain proportional. The goal is to explain execution, not to create a second system mirroring every internal operation.

---

## 3.5 Cancellation is foundational

Cancellation affects:

- provider requests
- streams
- tool execution
- subprocesses
- retries
- backoff waits
- event delivery
- cleanup
- final state

It cannot be reliably added at the end.

Every phase that introduces blocking work must define cancellation behaviour.

---

## 3.6 A successful call is not a successful task

A provider can return successfully while:

- producing malformed output
- failing to call a required tool
- editing the wrong file
- ignoring a constraint
- stopping too early
- returning incomplete structured data

Execution result, model result and task result are distinct concepts.

Aren should preserve this distinction throughout its design.

---

## 3.7 Failure needs semantics

String errors are insufficient for reliable agent workflows.

The caller may need to distinguish:

- cancellation
- timeout
- authentication
- rate limiting
- unavailable provider
- invalid request
- malformed response
- tool failure
- permission denial
- validation failure
- exhausted retries
- internal invariant violation

The taxonomy should start small and expand only when callers need different behaviour.

---

## 3.8 Unknown is a valid state

Agent systems frequently have partial information.

Examples include:

- missing token usage
- estimated cost
- uncertain retry timing
- provider-specific stop reasons
- unknown tool completion
- incomplete session continuation
- ambiguous network failure

Aren should not replace unknown values with misleading defaults.

---

## 3.9 Local-first is a starting advantage, not a permanent restriction

Local operation reduces complexity and accelerates learning.

Aren should begin as an in-process library and CLI. Its implementation language
remains an open decision. Once selected, the language is intended for the
long-term runtime rather than an initial implementation followed by a planned
port.

A daemon becomes justified when the system needs capabilities such as:

- persistent ownership of long-running execution
- several clients
- cross-language access
- centralised lifecycle control
- reconnectable streams
- execution surviving client termination
- shared scheduling or resource management

The daemon should be added around a proven core.

---

## 3.10 Planning must terminate in implementation

Research and planning are valuable only when they improve decisions.

Every Aren planning artefact should eventually connect to:

- a decision
- a phase boundary
- a contract
- an implementation task
- a test
- a rejected alternative

Documents that no longer affect action should be archived rather than continually expanded.

---

# 4. Aren Inheritance Map

## Inherit

These ideas have strong evidence from previous projects and should be treated as likely Aren principles:

- local-first development
- filesystem-visible artefacts
- explicit execution lifecycle
- structured events
- cancellation throughout the stack
- classified failures
- bounded retry behaviour
- validation separate from execution
- explicit permission posture
- honest unknown values
- evidence-backed design decisions
- smoke testing after implementation
- phased development
- explicit scope guardrails
- abstraction only after demonstrated pressure

## Re-evaluate

These ideas may be useful but must be proved in Aren:

- runtime-neutral provider interfaces
- middleware or interceptor systems
- DAG execution
- event buses
- run stores
- plugin systems
- graph-based context
- durable execution
- model fallback
- automatic output repair
- daemon architecture
- WebSocket or SSE client transport
- multi-language SDKs
- generic workflow orchestration
- multiple concurrent agents

## Reject by default

These patterns should not be introduced without unusually strong evidence:

- designing the complete final architecture before phase one
- abstractions with only one trivial implementation
- interfaces mirroring concrete types without behavioural value
- a server purely because the future product may need one
- a database for state that can remain files or memory
- multi-agent coordination before a single-agent loop is excellent
- plugin APIs before internal contracts stabilise
- broad configuration surfaces for hypothetical users
- duplicated planning documents
- infrastructure that does not support a current vertical slice
- features justified primarily by competitors or frameworks having them
- calling a component robust without testing its failure behaviour

---

# 5. Development Doctrine for Aren

The following doctrine should govern roadmap and phase design.

## 5.1 One phase, one primary uncertainty

Each phase should answer one major technical question.

Examples:

- Can Aren perform a cancellable streaming model call reliably?
- Can Aren validate structured output without corrupting stream semantics?
- Can Aren execute a tool with clear lifecycle and failure behaviour?
- Can Aren run a bounded agent loop?
- Can Aren manage context without hidden mutation?

A phase may contain several implementation tasks, but they should converge on one primary uncertainty.

---

## 5.2 Every phase must produce a usable system

The output of each phase should run.

A phase consisting only of types, interfaces or internal architecture is incomplete unless those are exercised through a real path.

---

## 5.3 Every phase must include failure testing

Normal-path tests are insufficient.

Each phase should identify and test relevant failures such as:

- cancellation
- provider timeout
- malformed events
- interrupted streams
- invalid structured output
- tool crashes
- permission denial
- retry exhaustion
- cleanup failure
- partial state

---

## 5.4 Real use comes before broadening

After a phase is technically complete, Aren should be used in a small real workflow.

The next phase should not begin until this use has revealed whether:

- the API is awkward
- events are missing
- failure semantics are unclear
- unnecessary abstractions exist
- important state is hidden
- the feature solves the intended problem

---

## 5.5 Reviews may remove capabilities

A phase review is not only a gate for adding the next feature.

It may conclude that:

- an abstraction should be deleted
- an interface should become concrete
- a configuration option should be removed
- an event should be merged
- a layer should collapse
- the next planned feature is not yet justified

Reduction is valid progress.

---

## 5.6 Future extensibility should come from clear internals

Aren should not attempt to guarantee unlimited extensibility from the first release.

The better route is:

1. keep the core small,
2. make ownership explicit,
3. avoid hidden global state,
4. separate behaviour only where necessary,
5. test contracts,
6. refactor when a second use case arrives.

A system with clear internals can become extensible later.

A system full of speculative extension points can become difficult to change immediately.

---

# 6. Repository Lineage

## Primary lineage repositories

These repositories deserve direct study because they materially shaped Aren:

| Project | Contribution to Aren |
|---|---|
| Elevate | Early product ambition and scope-creep lessons |
| CoachUp | Rich product design before a narrow validated core |
| Eloquence | Continued exploration of AI coaching and product complexity |
| Stageflow | Explicit execution stages, workflow semantics and observability |
| Unified Content Protocol | Generalisation, structured content and abstraction lessons |
| Hivemind | Multi-agent orchestration and coordination complexity |
| 24-Hour Testers | Autonomous validation loops and durable evidence |
| Ultra / Ultraground | Filesystem-grounded research and reasoning |
| UltraPlan | Staged study, planning, execution, review and smoke testing |
| UltraPlan Go | Simplification into a local-first Go CLI |
| AgentWrap | Runtime supervision, lifecycle control and reliability contracts |
| AgentWrap Smoke | Real-runtime validation and integration evidence |
| Aren | Ownership of the agent harness core |

## Supporting repositories

Other repositories provide useful supporting evidence in areas such as:

- graph modelling
- benchmarking
- API architecture
- frontend product development
- learning tools
- testing
- Go CLI design
- repository studies
- workflow experiments

They should remain discoverable in the wider catalogue but should not receive equal weight in the main lineage.

---

# 7. Questions to Carry into Aren Planning

This history does not determine Aren’s design. It establishes questions that the roadmap and PRD must answer.

1. What is the smallest independently useful version of Aren?
2. What concrete workflow will validate the first phase?
3. Which AgentWrap behaviours belong in the Aren core, and which were specific to wrapping subprocess runtimes?
4. Should the first public surface be a library API, CLI, or both?
5. What execution state must be externally observable from the first phase?
6. What does cancellation mean at each layer?
7. What is the minimum useful error taxonomy?
8. How will Aren distinguish provider completion from task success?
9. Which future capabilities are explicitly excluded from phase one?
10. What evidence would justify introducing a daemon?
11. What evidence would justify introducing a second provider?
12. What evidence would justify a general tool interface?
13. What evidence would justify durable state?
14. How will phase reviews detect new scope creep?
15. How will UltraPlan support the project without becoming a source of over-planning?

---

# 8. Final Perspective

Aren is described as Antonio’s final agent harness.

That ambition is understandable, but it contains a danger.

Trying to ensure Aren is the final harness could recreate the exact behaviour that caused earlier projects to expand: designing now for every future need so that the system never has to be replaced.

A more useful interpretation is:

> Aren should be the harness that is developed with enough discipline that it can continue evolving without repeatedly being abandoned and restarted.

That does not require getting the final architecture right at inception.

It requires:

- foundations that are genuinely understood
- small validated steps
- honest treatment of uncertainty
- willingness to remove mistaken abstractions
- reliable tests
- clear phase boundaries
- continuous real use
- resistance to speculative scope

The goal is not perfection at the beginning.

The goal is a development process capable of producing durable quality over time.

Aren’s most important inheritance is therefore not a runtime contract, event system, graph model or workflow primitive.

It is the lesson that the route to a robust and extensible system is not to build everything it might eventually need.

It is to build the next essential layer well enough that the correct following layer becomes visible.
<<< END ULTRAPLAN DIRECT PROJECT INPUT >>>

<<< BEGIN ULTRAPLAN DIRECT PROJECT INPUT >>>
ID: project-source-aren-phased-roadmap
Kind: project-source-document
Path: projects/aren-phase-01-execution-lifecycle/docs/phased-roadmap.md
Assignment: Catalogued project source document: Aren Phased Roadmap
Mode: full
Original-Bytes: 52209
Injected-Bytes: 52209

# Aren Phased Roadmap

## 1. Purpose

Aren will be developed from a small, local execution runtime into a broader autonomous execution system.

The long-term vision is:

> A general-purpose runtime that owns and supervises autonomous execution independently of any particular agent, provider, language, transport, or deployment model.

That vision is directional rather than an initial specification.

Aren will not begin by trying to represent coding agents, model calls, shell commands, workflows, approvals, and scheduled tasks through one complete universal abstraction. It will begin by proving a small execution model against controlled implementations, then broaden only when real variation provides evidence for new abstractions.

The intended progression is:

```text
Execution lifecycle
    ↓
Controlled execution
    ↓
Real model invocation
    ↓
Streaming and structured results
    ↓
Retries and repair
    ↓
Tool execution
    ↓
Bounded agent loop
    ↓
Context management
    ↓
Persistence and recovery
    ↓
Execution composition
    ↓
Policies and resource governance
    ↓
Daemon hosting and remote clients
    ↓
Workflows
    ↓
Additional execution types
    ↓
Distributed execution, if justified
```

Later phases are hypotheses, not commitments. Their order and contents should change in response to evidence from real Aren usage.

---

## 2. Development Rules

### 2.1 One primary uncertainty per phase

Each phase should answer one major technical question. A phase may contain several implementation tasks, but they must converge on that question.

### 2.2 Every phase must produce something runnable

Types, interfaces, and design documents are not sufficient by themselves. Every phase must expose a working vertical slice through one or more of:

- an importable library;
- a development CLI;
- an executable example;
- a real smoke workflow.

### 2.3 Behaviour before abstraction

Start with one concrete implementation.

Introduce a general interface only when there is real evidence, such as:

- a second implementation with meaningful behavioural differences;
- repeated conditional logic;
- duplicated lifecycle handling;
- a testing boundary needed to reproduce failures;
- ownership that genuinely needs to cross a boundary;
- a concept that has remained stable across repeated use.

The possibility that an abstraction may be useful later is not enough.

### 2.4 Failure testing is mandatory

Every phase must validate, where relevant:

- normal operation;
- expected failures;
- cancellation;
- timing and race conditions;
- cleanup;
- partial results;
- recovery or honest termination.

### 2.5 Real use before progression

After implementation and automated testing, the phase capability must be exercised in a small real workflow.

The phase review should ask:

- Is the API awkward?
- Are important events missing?
- Are any events redundant?
- Are failures sufficiently classified?
- Is ownership clear?
- Is hidden state affecting behaviour?
- Did an abstraction appear before it was necessary?
- Should anything be removed before continuing?

### 2.6 Reviews may remove capabilities

A phase review may conclude that:

- an interface should become concrete;
- a package should be collapsed;
- a configuration option should be removed;
- an event should be merged or deleted;
- a planned capability is not yet justified.

Reduction is valid progress.

### 2.7 Foundations must be vertical slices

“Foundations first” does not mean building a large invisible substrate. Each foundation must be exercised through a complete, observable path that reveals real behaviour.

### 2.8 Later phases remain revisable

The first phases establish Aren’s core semantics. Later phases are deliberately less detailed and must be reconsidered after every phase review.

---

## 3. Phase Overview

| Phase | Capability | Primary uncertainty |
|---|---|---|
| 0 | Project foundation | Can Aren maintain strict scope and trustworthy engineering feedback? |
| 1 | Execution lifecycle | Can Aren define and enforce a small, coherent execution lifecycle? |
| 2 | Controlled executor | Can that lifecycle survive progress, failure, cancellation, cleanup, and races? |
| 3 | First real model invocation | Can Aren own one real LLM request without losing lifecycle control? |
| 4 | Streaming model execution | Can streamed output remain ordered, cancellable, bounded, and observable? |
| 5 | Structured model results | Can Aren distinguish valid task results from successful provider calls? |
| 6 | Retry and bounded repair | Can Aren recover selectively without hiding failures or creating uncontrolled loops? |
| 7 | Tool-call representation | Can requested actions be represented independently from their implementation? |
| 8 | Local tool execution | Can Aren supervise external actions through the same reliable principles? |
| 9 | Bounded agent loop | Can Aren own a minimal model–tool loop with explicit termination conditions? |
| 10 | Context management | Can context be managed deliberately without hiding mutation or provenance? |
| 11 | Persistence and recovery | Which state genuinely needs to survive process termination? |
| 12 | Execution composition | How should executions be coordinated before a workflow system is justified? |
| 13 | Policies and resources | Which permissions, budgets, and concurrency controls require runtime ownership? |
| 14 | Daemon hosting | When should execution ownership outlive a single client process? |
| 15 | Multi-language clients | Can other languages control Aren without duplicating runtime behaviour? |
| 16 | Workflows | Can reusable, inspectable, and recoverable processes be built from proven execution primitives? |
| 17 | Broader execution types | Which non-agent forms of work genuinely belong under Aren? |
| 18 | Distributed execution | Is there sufficient evidence for remote workers or clustering? |

---

# Foundation

## Phase 0 — Project Foundation and Scope Control

### Primary question

> Can Aren establish a development environment and decision process that make incorrect behaviour visible without building product infrastructure prematurely?

Phase 0 is not an architecture phase. It creates only the minimum repository foundation needed to develop Phase 1 safely.

### Goals

- Select one implementation language for the long-term runtime, then establish a small repository for it.
- Make builds, tests, static checks, and concurrency-safety checks easy to run.
- Define project terminology before APIs proliferate.
- Record explicit scope boundaries.
- Create lightweight architectural decision records.
- Prevent roadmap work from turning into implementation speculation.

### In scope

#### Repository foundation

A minimal repository should contain runtime source, a development CLI entry
point, tests, examples, decision records, a glossary, the selected language's
manifest, and a README.

The exact structure should remain small and may change during Phase 1.

#### Engineering feedback

Documented commands for:

- formatting;
- unit tests;
- concurrency-safety tests;
- static analysis;
- coverage inspection;
- building the CLI.

#### Initial glossary

Define only the terms already required:

- execution;
- run;
- state;
- event;
- result;
- failure;
- cancellation;
- executor.

Definitions may be marked provisional.

#### Scope ledger

Maintain a lightweight record containing:

- current phase scope;
- explicit exclusions;
- deferred ideas;
- evidence required to reconsider each deferred idea.

#### Decision records

Use short decision records only for choices that would otherwise become ambiguous.

Likely initial decisions:

- one implementation language for the long-term runtime;
- library-first runtime;
- no daemon in the initial phases;
- local and in-memory operation initially;
- no stable public API promise before the foundational semantics settle.

### Out of scope

- provider integrations;
- execution persistence;
- databases;
- servers;
- plugin systems;
- workflow definitions;
- SDK generation;
- broad configuration frameworks;
- observability backends;
- production deployment;
- semantic-versioning guarantees.

### Deliverables

- buildable implementation module;
- minimal CLI entry point;
- automated test and concurrency-check commands;
- glossary;
- phase scope ledger;
- decision-record format;
- initial repository README.

### Validation

- A clean checkout can build and test with one documented command.
- The selected toolchain's strongest practical concurrency checks run locally and in CI.
- No production runtime abstractions are introduced solely for repository setup.
- Every nontrivial source module has an immediate purpose.

### Exit gate

Phase 0 is complete when Aren has selected one implementation language for the
long-term runtime and the repository reliably supports Phase 1 development
without carrying speculative runtime architecture.

---

## Phase 1 — Execution Lifecycle

### Primary question

> Can Aren define and enforce the lifecycle of one execution without depending on an LLM, subprocess, network call, or persistent store?

This is the most foundational phase.

The objective is not to create a universal executor interface. It is to establish the minimum semantics that Aren itself owns.

### Conceptual model

A run represents one supervised occurrence of work.

A preliminary lifecycle is:

```text
created
   ↓
running
   ├──→ succeeded
   ├──→ failed
   └──→ cancelled
```

A cancellation request may occur while running:

```text
running
   ↓
cancellation requested
   ↓
cancelled
```

Whether `cancellation_requested` is a state, an event, or both must be resolved through implementation. The lifecycle should remain much smaller than a workflow-engine state machine.

### Goals

- Assign every run a unique identity.
- Define who owns state transitions.
- Define all legal and illegal transitions.
- Represent terminal outcomes explicitly.
- Separate lifecycle state from output data.
- Establish cancellation semantics.
- Establish the first canonical event vocabulary.
- Define event observation behaviour.
- Define a minimum failure representation.
- Make invariant violations visible.

### Run identity

Decide:

- identifier type;
- generation ownership;
- whether callers may provide IDs;
- whether IDs carry semantic information.

Initial direction:

- opaque unique IDs;
- generated by Aren;
- no execution type or timestamp encoded into the ID.

### Lifecycle state

Candidate states:

- `created`;
- `running`;
- `succeeded`;
- `failed`;
- `cancelled`.

Questions to resolve:

- Is `created` externally observable?
- Is cancellation request a state or only an event?
- Can a run fail before entering `running`?
- Is rejection a failed run or a failure to create one?
- Can cleanup failure alter an otherwise successful outcome?

### Terminal outcome

State alone is insufficient. A terminal outcome may include:

- terminal state;
- result value;
- classified failure;
- start time;
- finish time;
- cancellation metadata.

Result and failure must remain coherent.

Valid examples:

```text
succeeded + result
failed + failure
cancelled + cancellation reason
```

Invalid combinations include:

```text
succeeded + failure
failed + successful result
running + terminal timestamp
```

### State ownership

Aren, not executor code, owns legal lifecycle transitions.

Executed work may return an outcome or report progress, but it must not directly mutate the run state machine.

### Cancellation

Cancellation semantics must answer:

- Who may request cancellation?
- Is cancellation idempotent?
- What happens when cancellation races with completion?
- Is cancellation cooperative?
- When is a run considered cancelled?
- What happens if work ignores cancellation?
- Does a cancellation reason belong in the outcome?
- Are repeated cancellation requests observable?

Initial direction:

- cancellation requests are idempotent;
- the first valid terminal outcome wins;
- Aren propagates an implementation-appropriate cancellation signal;
- Aren does not claim the work has stopped until controlled work returns;
- cancellation request and cancellation completion are distinct;
- a successful outcome may remain successful if it wins the race against cancellation.

These rules must be proved by tests, not accepted only in prose.

### Event ordering

The initial event vocabulary should describe lifecycle behaviour only.

Candidate events:

- `run.created`;
- `run.started`;
- `run.cancellation_requested`;
- `run.succeeded`;
- `run.failed`;
- `run.cancelled`.

Likely event metadata:

- run ID;
- per-run sequence number;
- event type;
- occurrence time;
- immutable payload.

Initial direction:

- ordering is established by per-run sequence number, not timestamp;
- state transition and event creation occur together inside the runtime;
- delivery may lag but must preserve order;
- consumer failure must not corrupt run state;
- terminal events occur exactly once.

### Event observation

Avoid building a general event bus.

Possible APIs include:

```text
run.events()
```

or:

```text
subscription = run.subscribe()
```

A raw channel creates semantic questions that must be answered explicitly:

- buffer size;
- slow consumers;
- closure ownership;
- missed events;
- multiple subscribers;
- replay behaviour;
- observer abandonment.

The implementation should choose the smallest observation model whose behaviour can be tested honestly.

### Failure model

Begin with a minimal structured failure representation.

Likely categories:

- execution failure;
- cancellation;
- internal invariant violation.

The representation should preserve:

- machine-readable category;
- human-readable message;
- wrapped cause where applicable;
- whether the failure arose from executed work or Aren itself.

Do not design the future provider and tool error taxonomies yet.

### Time

The runtime needs observable timing information:

- start time;
- terminal time;
- possibly duration.

A public clock abstraction should only be introduced if deterministic testing proves that it is necessary.

### First implementation

Use an in-process work function as a semantic instrument:

```text
work(cancellation) -> result | failure
```

A run controller should:

1. create the run identity;
2. emit creation and start events;
3. invoke the work function;
4. propagate cancellation;
5. classify the returned outcome;
6. perform one legal terminal transition;
7. expose the immutable outcome and events.

This is not yet a promise of the permanent executor API.

### Illustrative public surface

```text
run = aren.start(parent, work)

for event in run.events():
    observe(event)

outcome = run.wait()
run.cancel("user requested cancellation")
```

The PRD should specify behaviour before exact names and signatures.

### Required tests

#### Normal lifecycle

- `created → running → succeeded`;
- result preserved exactly;
- timestamps recorded coherently;
- events ordered correctly;
- waiting after completion returns immediately.

#### Failure lifecycle

- `created → running → failed`;
- original cause retained;
- failure event emitted exactly once;
- no successful result exposed.

#### Cancellation

- cancellation before work begins, if supported;
- cancellation during execution;
- repeated cancellation;
- work that cooperates immediately;
- work that delays acknowledgement;
- work that ignores cancellation until it returns.

#### Completion–cancellation race

Run many iterations with the selected toolchain's strongest practical
concurrency checks enabled.

Required invariants:

- exactly one terminal state;
- exactly one terminal event;
- no data race;
- result agrees with terminal state.

#### Multiple waiters

- all waiters receive the same immutable outcome;
- no waiter consumes the outcome exclusively;
- no deadlock.

#### Event consumers

- slow consumer;
- absent consumer;
- abandoned consumer;
- multiple consumers if supported;
- documented overflow behaviour;
- no leak caused solely by nobody reading events.

#### Panic and invariants

- work panic;
- illegal state transition;
- no silent invariant corruption.

A panic may become an explicitly classified internal execution failure, but it must remain recognisable as a panic rather than being disguised as an ordinary task failure.

### Diagnostic CLI

Provide a small development CLI:

```text
aren dev run success
aren dev run fail
aren dev run cancel
aren dev run race
```

Example output:

```text
000 run.created
001 run.started
002 run.cancellation_requested
003 run.cancelled
```

The CLI is a semantic inspection tool, not the final Aren product interface.

### Deliverables

- lifecycle model;
- run controller;
- terminal outcome model;
- cancellation contract;
- initial event vocabulary;
- deterministic event sequencing;
- minimum failure representation;
- comprehensive unit and concurrency tests;
- diagnostic CLI;
- written lifecycle contract;
- phase review.

### Out of scope

- attempts and retries;
- providers;
- model messages;
- streaming output;
- tools;
- subprocess management;
- persistence;
- replay after process restart;
- multiple execution types;
- daemon transport;
- global event streams;
- workflow composition;
- production telemetry export.

### Exit gate

Phase 1 is complete only when:

1. all lifecycle transitions are defined;
2. cancellation races are proven under repeated concurrency checking;
3. event order is deterministic;
4. exactly one terminal outcome is guaranteed;
5. slow or absent consumers cannot accidentally deadlock or leak a run;
6. a runnable demonstration exists;
7. no foundational semantic ambiguity remains unresolved.

---

## Phase 2 — Controlled Executor and Failure Laboratory

### Primary question

> Does the Phase 1 lifecycle remain coherent when execution exhibits realistic progress, delay, failure, cleanup, partial output, and cancellation behaviour?

Phase 1 proves the lifecycle around a minimal work function. Phase 2 introduces a deliberately configurable executor used to attack those semantics.

It must not grow into a general simulation framework.

### Goals

- Exercise richer controlled behaviour.
- Introduce nonterminal progress.
- Test delayed and partially cooperative cancellation.
- Define cleanup behaviour.
- Test partial output followed by failure.
- Determine whether an executor interface is now genuinely earned.
- Produce a reusable conformance suite for future executors.

### Controlled behaviours

The test executor should be able to:

- succeed immediately;
- succeed after a delay;
- fail immediately;
- fail after progress;
- emit a fixed progress sequence;
- block until externally released;
- cooperate with cancellation;
- delay cancellation response;
- ignore cancellation until a safe point;
- panic;
- fail during cleanup;
- race completion against cancellation;
- produce partial output before failure.

Configuration should be typed and explicit. Do not create a scenario language.

### Progress

Introduce only the smallest useful progress representation.

Initial candidate:

```go
type Progress struct {
    Message string
}
```

A generic metadata bag or universal progress schema should not be introduced without concrete need.

### Cleanup semantics

Decide:

- whether cleanup completes before the terminal event;
- whether cleanup failure can convert success into failure;
- how cleanup behaves after cancellation;
- whether cleanup receives a separate timeout;
- how execution and cleanup errors are preserved together.

Initial direction:

- required cleanup completes before the terminal state is published;
- cleanup errors are never silently discarded;
- mandatory cleanup failure prevents an automatic success outcome;
- cleanup must not run forever under an already cancelled signal.

### Executor boundary

A candidate abstraction may now be tested:

```text
executor.execute(cancellation, reporter) -> result | failure
```

It should only be retained if it provides real value through progress reporting, conformance testing, or ownership clarity. If a function remains simpler and sufficient, keep the function.

### Conformance suite

Future executors should be able to reuse behavioural tests for applicable invariants:

- returns one coherent outcome;
- respects the cancellation contract;
- emits no progress after terminal completion;
- preserves progress ordering;
- releases required resources;
- does not mutate published events;
- retains classified failures.

### Real vertical slice

Add one small non-LLM example. Prefer a platform-independent in-process task initially. Introduce subprocess semantics only if the foundation cannot be validated without them.

### Required tests

- progress before success;
- progress before failure;
- cancellation after partial progress;
- no progress after terminal event;
- cleanup on every exit path;
- cleanup failure;
- delayed cancellation acknowledgement;
- observer detachment;
- panic during execution;
- panic during cleanup;
- a large progress burst;
- documented backpressure behaviour;
- run start failure, if supported.

### Deliverables

- controlled semantic executor;
- progress mechanism;
- cleanup contract;
- executor conformance suite;
- evidence-based decision on an executor interface;
- one runnable non-LLM example;
- revised lifecycle contract;
- phase review.

### Out of scope

- retries;
- attempts unless demanded by evidence;
- model providers;
- arbitrary subprocess supervision;
- remote execution;
- durable progress;
- global event routing;
- plugin registration.

### Exit gate

Phase 2 is complete when the lifecycle contract has survived realistic controlled behaviour and future executor integrations have a clear, tested path.

---

# Model Execution

## Phase 3 — First Real Model Invocation

### Primary question

> Can Aren own one real model request while preserving execution, cancellation, event, and failure semantics?

This is the first external integration. The goal is one excellent concrete model invocation, not provider neutrality.

### Initial scope

Support one current OpenAI API path with:

- one model request type;
- text input;
- text output;
- no tools;
- no structured schema;
- no automatic retry;
- no session abstraction;
- non-streaming operation initially.

### Goals

- Send one real request.
- Map provider activity into Aren lifecycle semantics.
- Propagate cancellation into the network request.
- Capture response metadata honestly.
- Classify the smallest actionable provider failures.
- Preserve provider-specific details behind a bounded escape hatch.
- Validate behaviour with deterministic transport tests and real smoke tests.

### Input and result

Begin with a concrete request rather than a universal message graph.

```go
type ModelRequest struct {
    Model string
    Input string
}
```

Capture, where available:

- generated text;
- provider request ID;
- actual model used;
- stop or completion reason;
- token usage;
- provider metadata required for diagnosis.

Unknown values must remain unknown rather than becoming zero or empty success values.

### Initial provider failures

Classify only failures callers may handle differently:

- authentication;
- invalid request;
- rate limited;
- provider unavailable;
- timeout;
- network or transport failure;
- malformed provider response;
- cancelled;
- unknown provider failure.

Retain causes and relevant metadata such as retry-after information.

### Candidate events

- `model.request.started`;
- `model.request.completed`;
- `model.request.failed`.

Example sequence:

```text
run.started
model.request.started
model.request.completed
run.succeeded
```

Avoid duplicating low-level HTTP details.

### Testing

Use a controlled local HTTP server or injectable transport to simulate:

- success;
- authentication failure;
- rate limiting;
- retry-after metadata;
- delayed response;
- cancellation;
- malformed or partial bodies;
- connection closure;
- provider 5xx;
- unknown status;
- missing usage;
- unknown response fields.

Real smoke tests should be opt-in and use environment-provided credentials.

### Out of scope

- streaming;
- multiple providers;
- provider abstraction;
- fallback;
- retry;
- structured output;
- tool calls;
- prompt templating;
- context reduction;
- persistent sessions.

### Exit gate

One real model invocation is reliable, cancellable, observable, and honestly represented under normal and adverse transport conditions.

---

## Phase 4 — Streaming Model Execution

### Primary question

> Can Aren expose incremental model output while preserving event order, cancellation, terminal consistency, and consumer safety?

Streaming changes runtime concurrency and observation semantics. It is not only a UI feature.

### Goals

- Receive provider streaming responses.
- Emit output deltas in deterministic order.
- Reconstruct the final output exactly.
- Cancel midstream.
- Handle partial and interrupted output honestly.
- Define bounded buffering and backpressure.

### Candidate events

- `model.output.delta`;
- possibly `model.stream.interrupted` if it adds information not already represented by request failure.

A likely sequence is:

```text
run.started
model.request.started
model.output.delta
model.output.delta
model.request.completed
run.succeeded
```

### Delta semantics

Define:

- what a delta contains;
- whether empty deltas are emitted;
- how Unicode boundaries are handled;
- whether final text must equal the ordered concatenation of text deltas;
- how provider annotations are represented;
- whether raw provider chunks are exposed.

Initial direction:

- emit normalised text deltas for the supported request type;
- retain final reconstructed text in the result;
- keep raw chunks internal or behind an explicit low-level path;
- order all events through the run sequence.

### Partial output

A cancelled or interrupted stream may have useful partial output. It should be retained as diagnostic execution data, but not presented as a complete successful model result.

The result must distinguish:

- complete output;
- partial output caused by cancellation;
- partial output caused by provider or transport failure.

### Backpressure

The implementation must choose and document one bounded policy. It must not accidentally create:

- unbounded memory growth;
- a runtime task per event;
- provider deadlock caused by an unread public channel;
- loss of the terminal outcome.

### Required tests

- one delta and many deltas;
- Unicode boundaries;
- empty stream;
- cancellation before first delta;
- cancellation after several deltas;
- interrupted stream;
- malformed event;
- duplicate provider completion signal;
- provider EOF without an explicit final marker;
- slow or absent observer;
- exact final reconstruction;
- partial-output classification;
- no events after the run terminal event.

### Exit gate

Streaming is useful without weakening lifecycle correctness or introducing hidden concurrency behaviour.

---

## Phase 5 — Structured Model Results and Validation

### Primary question

> Can Aren distinguish provider completion from successful task completion through explicit result validation?

```text
Provider execution succeeded
            ≠
Requested result is valid
```

### Goals

- Request structured output through one supported mechanism.
- Parse the returned representation.
- Validate it against one explicit contract.
- Distinguish parsing and validation failures from provider failures.
- Preserve raw output for diagnosis.
- Avoid coupling every execution to JSON or schemas.

### Result layers

Keep these distinct:

1. provider response;
2. raw generated content;
3. parsed representation;
4. validated result.

A malformed result should retain enough evidence to explain what was returned and why it failed.

### Initial validation scope

- one schema mechanism;
- deterministic parsing and validation;
- no LLM evaluator;
- no automatic repair;
- no generic validator registry;
- no arbitrary network-backed validation.

### Candidate events

- `validation.started`;
- `validation.succeeded`;
- `validation.failed`.

Parsing may be represented as one kind of validation failure unless callers need a different control path.

### Outcome semantics

A provider call that returns malformed structured output may have succeeded at the provider layer while the requested structured execution fails overall.

The run should:

- terminate with a validation failure;
- preserve the raw provider response;
- expose no validated result;
- keep provider and task-level evidence distinct.

### Exit gate

Aren can honestly report that a provider responded successfully while the requested result failed validation.

---

## Phase 6 — Retry and Bounded Repair

### Primary question

> Can Aren repeat or repair failed work selectively without hiding the original failure, duplicating effects, or creating uncontrolled loops?

### Goals

- Introduce attempts explicitly.
- Retry only classified transient failures.
- Support cancellable backoff.
- Preserve each attempt’s evidence.
- Define exhaustion clearly.
- Add bounded structured-output repair only after ordinary retry semantics are stable.

### Attempt model

Initial direction:

- an attempt is a subordinate record inside one run;
- attempts have an ordinal and timestamps;
- events identify attempt number;
- one run still reaches exactly one terminal outcome;
- attempts do not initially expose the complete public run API.

Do not recursively define every attempt as a full execution unless later evidence justifies it.

### Retry policy

A concrete policy may contain:

- maximum attempts;
- initial delay;
- maximum delay;
- multiplier;
- optional jitter;
- retryable failure categories.

Default behaviour is no retry unless explicitly configured.

### Candidate events

- `attempt.started`;
- `attempt.failed`;
- `retry.scheduled`;
- `attempt.succeeded`;
- `retry.exhausted`.

The final event set should be reduced after implementation if some events prove redundant.

### Cancellation

Test cancellation:

- during an active request;
- after attempt failure;
- during backoff;
- immediately before another attempt;
- racing with final success.

No wait may ignore cancellation.

### Bounded repair

After retry is stable, support one narrow repair path for invalid structured output.

```text
attempt 1: provider succeeds
validation fails
repair scheduled
attempt 2: original output + validation feedback
validation succeeds
run succeeds
```

Repair must preserve:

- original output;
- validation errors;
- repair instruction;
- repaired output;
- final validation result.

Repair is visible changed-input execution, not a hidden retry.

### Out of scope

- provider fallback;
- hedged requests;
- global retry budgets;
- durable retries after crash;
- automatic tool retries;
- arbitrary repair strategies;
- open-ended self-reflection loops.

### Exit gate

Retries are explicit, bounded, observable, cancellable, and unable to disguise the failures that preceded the final outcome.

---

# Tools and Agents

## Phase 7 — Tool-Call Representation

### Primary question

> Can Aren represent a model’s request for external action independently from how that action is implemented?

This phase separates tool definition from tool execution. It does not yet build the agent loop.

### Goals

- Define a tool for model consumption.
- Represent a requested tool call.
- Parse and validate tool arguments.
- Represent unknown tools without losing the provider response.
- Represent tool results and failures.
- Handle streamed arguments if the provider requires it.
- Avoid coupling tool definitions to implementation-language function signatures.

### Minimum tool definition

- stable name;
- description;
- input schema.

Output schemas should be added only if a real use case requires them.

### Minimum tool call

- provider call ID;
- tool name;
- raw arguments;
- parsed arguments;
- validation status.

An unknown tool call should remain inspectable rather than failing during response parsing.

### Out of scope

- tool registry;
- tool execution;
- parallel tools;
- agent loop;
- remote tools;
- MCP;
- permissions;
- retries;
- client-hosted callbacks.

### Exit gate

Aren can receive, inspect, and validate a requested action without knowing how it will execute.

---

## Phase 8 — Local Tool Execution

### Primary question

> Can Aren supervise a concrete external action with clear lifecycle, result, failure, permission, and cancellation semantics?

Begin with in-process native tools only.

### Goals

- Register a small set of local implementations.
- Resolve validated calls to implementations.
- Execute tools with cancellation.
- Emit tool lifecycle events.
- Validate inputs before side effects.
- Capture typed output or failure.
- Establish an explicit allow or deny permission posture.
- Avoid assuming tools are safe to retry.

### Candidate events

- `tool.started`;
- `tool.completed`;
- `tool.failed`;
- `tool.cancelled`.

### Initial failure categories

- unknown tool;
- invalid arguments;
- permission denied;
- execution failure;
- timeout;
- cancelled;
- invalid result;
- tool panic.

### Permissions

Begin with a direct allow or deny decision. Do not build a policy language.

The key invariant is:

> Required permission is resolved before the tool performs side effects.

### Retry safety

Tool execution is not automatically retried. Idempotency and retry declarations may be added later only when real tools require them.

### Real tools

Start with harmless, useful tools such as:

- reading a text file within an allowed directory;
- listing a directory;
- deterministic calculations.

Do not begin with arbitrary shell execution.

### Exit gate

A model-requested action can be validated and executed locally without weakening Aren’s lifecycle and failure guarantees.

---

## Phase 9 — Minimal Bounded Agent Loop

### Primary question

> Can Aren own a model–tool loop that remains bounded, observable, and understandable?

This is the first conventional agent phase.

### Loop

```text
initial context
    ↓
model invocation
    ↓
assistant output
    ├── final response → stop
    └── tool calls
            ↓
       execute tools
            ↓
       append results
            ↓
       next model invocation
```

### Goals

- Own the loop inside Aren rather than an SDK.
- Support sequential tool calls.
- Maintain explicit conversation state.
- Enforce hard termination conditions.
- Preserve child model and tool evidence.
- Produce one final agent-run outcome.
- Keep control flow linear and direct.

### Required limits

- maximum model turns;
- maximum tool calls;
- maximum elapsed time;
- cancellation;
- optionally token or cost limits when the underlying measurements are trustworthy.

Unknown metrics must not be silently treated as zero.

### Completion semantics

Distinguish:

- model final response;
- exhausted turn limit;
- exhausted tool-call limit;
- unrecoverable model failure;
- tool failure;
- validation failure;
- cancellation.

A model stopping without a final answer is not automatically a successful agent result.

### Initial constraints

- sequential tools only;
- no subagents;
- no planning framework;
- no dynamic tool registration during a run;
- no human approval steps;
- no persistent session;
- no workflow graph.

### Candidate events

```text
run.started
agent.turn.started
model.request.started
model.output.delta
model.request.completed
tool.started
tool.completed
agent.turn.completed
run.succeeded
```

Only retain turn-level events if they add useful boundaries rather than duplicating child events.

### Real workflow

Use Aren to complete one narrow local task, such as inspecting a small directory, reading selected files, and answering a grounded question.

### Exit gate

One bounded agent performs useful work repeatedly without hidden control flow or ambiguous termination.

---

# Longer-Lived Execution

## Phase 10 — Context Engineering

### Primary question

> Can Aren manage growing context deliberately while preserving traceability and avoiding hidden mutation?

Context work should respond to pressure observed in real Phase 9 runs.

### Likely progression

1. Measure context growth.
2. Add deterministic size checks.
3. Separate instructions, history, tool results, and working data.
4. Add simple deterministic reduction.
5. Distinguish durable and reducible entries.
6. Add model-generated summarisation only when deterministic reduction is inadequate.
7. Retain provenance from derived context to original messages.

### Principles

- Original history is not silently rewritten.
- Summaries are explicit derived artefacts.
- Reduction decisions are observable.
- Policies are deterministic where possible.
- Provider token limits do not leak unpredictably into loop logic.
- “Memory” is not introduced as one broad undifferentiated feature.

### Deferred capabilities

- semantic retrieval;
- vector stores;
- long-term memory;
- graph context;
- workspace indexing;
- cross-run memory.

### Exit gate

Aren can support meaningfully longer tasks while every context transformation remains visible and attributable.

---

## Phase 11 — Persistence and Recovery

### Primary question

> Which Aren state must survive process termination, and what recovery guarantees are actually required?

Persistence should follow real lost-work pain, not anticipation.

### Likely progression

1. Persist completed run records.
2. Persist event history.
3. Persist active-run metadata.
4. Detect interrupted runs after restart.
5. Mark interruption honestly.
6. Consider resumability only after interruption semantics are understood.

### Critical distinctions

- recording is not recovery;
- recovery is not resumption;
- resuming an agent loop is not resuming an in-flight network request;
- durable events do not produce exactly-once side effects;
- replay does not imply re-execution safety.

### Storage

Start with one concrete local implementation, selected from demonstrated access patterns. Do not create a general storage abstraction before replacement pressure exists.

### Exit gate

Aren preserves the state users genuinely need and handles interrupted work honestly without claiming unsupported recovery guarantees.

---

## Phase 12 — Execution Composition

### Primary question

> How should several executions be coordinated while keeping control flow explicit and avoiding premature workflow machinery?

Composition should begin as ordinary program control flow.

### Likely progression

- execute A then B;
- branch based on a result;
- run a bounded parallel group;
- propagate cancellation from parent to children;
- retain parent–child run references;
- explain child failures in the parent outcome.

A dependency graph should only appear if repeated direct composition becomes genuinely difficult.

### Principles

- composition is not yet a workflow product;
- child execution ownership is explicit;
- cancellation flows predictably;
- parent completion waits for owned children;
- workflow state is not conflated with run state;
- no visual graph or declarative DSL yet.

### Exit gate

Several executions can be coordinated reliably, and their common composition pressures are understood well enough to inform a later workflow design.

---

## Phase 13 — Policies and Resource Governance

### Primary question

> Which cross-cutting controls require central runtime ownership?

### Likely capabilities

- model and tool allowlists;
- permission decisions;
- token and cost budgets;
- elapsed-time budgets;
- concurrency limits;
- shared rate-limit coordination;
- workspace boundaries;
- approval requirements;
- retry budgets.

### Principles

- policies act on explicit runtime facts;
- unavailable metrics remain unknown;
- denial occurs before side effects;
- policy decisions are observable;
- configuration remains smaller than the controlled behaviour;
- no general policy language without repeated rule complexity.

### Exit gate

Aren safely governs several real forms of execution without spreading duplicated policy logic through applications and executors.

---

# Hosting and External Clients

## Phase 14 — Daemon Hosting

### Primary question

> Has execution ownership outgrown the lifetime and language of a single client process?

A daemon becomes justified by demonstrated needs such as:

- execution surviving client exit;
- multiple clients;
- reconnectable observation;
- central cancellation;
- shared resource controls;
- scheduled execution;
- cross-language access;
- remote management.

### Architecture rule

The daemon hosts the runtime. It does not become the place where lifecycle, retry, tool, or agent-loop behaviour lives.

```text
Aren runtime
    ↓
daemon host
    ↓
transport adapter
```

### Initial daemon scope

- start a run;
- inspect a run;
- subscribe to events;
- cancel a run;
- retrieve its terminal outcome.

Choose one transport based on actual client needs. Do not implement HTTP, WebSockets, SSE, and gRPC simultaneously.

### Command and event distinction

Commands:

- start;
- cancel;
- inspect.

Events:

- lifecycle;
- progress;
- output;
- tool activity;
- terminal outcome.

### Exit gate

A client can disconnect and reconnect without Aren losing execution ownership, while all core semantics remain testable without the daemon.

---

## Phase 15 — Multi-Language Clients and Client-Hosted Tools

### Primary question

> Can Python and JavaScript use Aren naturally without implementing their own runtime behaviour?

### Likely sequence

1. Build one thin client.
2. Support start, observe, cancel, and inspect.
3. Add client-hosted tool callbacks.
4. Define callback timeout and disconnection semantics.
5. Add a second language only after the protocol stabilises.

### Client-hosted tool model

```text
Aren requests tool execution
        ↓
connected client executes callback
        ↓
client returns result
        ↓
Aren continues the loop
```

Questions that must be resolved:

- What happens if the client disconnects?
- Who owns the timeout?
- May another client satisfy the call?
- Can a callback request be replayed?
- How are permissions established?
- Can callbacks be retried safely?

### Exit gate

An application written in another language can use Aren as the execution owner while remaining a thin client rather than becoming a second harness implementation.

---

# Workflows and Broader Execution

## Phase 16 — Workflows

### Primary question

> Can Aren turn proven execution and composition primitives into reusable, inspectable, and recoverable processes without becoming a speculative workflow platform?

This phase is intentionally later than basic execution composition. Aren should first learn how real executions are started, observed, cancelled, retried, persisted, and composed in code. Only then should those patterns be represented as reusable workflow definitions.

### Goals

- Define a reusable process from proven execution primitives.
- Separate workflow definition from workflow instance state.
- Make step transitions explicit and observable.
- Support pause, cancellation, and failure propagation.
- Preserve parent–child execution relationships.
- Recover workflow state after daemon or process restart where supported.
- Introduce human or scheduled steps only when their distinct semantics are understood.

### Initial workflow scope

Begin with a small, linear workflow model:

- named steps;
- sequential execution;
- step inputs from prior outputs;
- explicit terminal states;
- per-step run references;
- cancellation of the workflow and owned child runs;
- persisted workflow-instance state;
- restart detection and honest interruption handling.

A minimal conceptual model:

```text
workflow created
      ↓
step A execution
      ↓
step B execution
      ↓
step C execution
      ↓
workflow completed
```

### Later workflow capabilities, if earned

- conditional branches;
- bounded parallel branches;
- retries at step or workflow level;
- compensation;
- human approval steps;
- scheduled starts;
- waiting for external events;
- reusable subworkflows;
- resumable long-running processes;
- declarative definitions;
- visualisation.

### Critical distinctions

#### Workflow definition vs workflow run

A definition describes reusable structure. A workflow run records one execution of that structure.

#### Workflow state vs execution state

A workflow may be waiting, paused, or blocked even when no child execution is currently running. Those semantics must not be forced into the core run lifecycle.

#### Retry vs resume

Retry starts work again. Resume continues from preserved workflow state. They must not be treated as synonyms.

#### Human waiting vs active execution

A human approval may remain pending for days. It should not pretend to be an active runtime task or network request.

#### Durable orchestration vs exactly-once effects

Persisted state does not guarantee that an external side effect occurred exactly once. Recovery must account for uncertain completion and idempotency.

### Workflow representation

Do not begin with a general graph DSL.

Likely progression:

1. direct typed builder or definition;
2. linear reusable workflow;
3. conditional branching after real use;
4. persisted definitions only if required;
5. external declarative formats only when authorship outside the implementation language is genuinely needed.

### Events

Potential workflow events:

- `workflow.created`;
- `workflow.started`;
- `workflow.step.started`;
- `workflow.step.completed`;
- `workflow.step.failed`;
- `workflow.paused`;
- `workflow.resumed`;
- `workflow.completed`;
- `workflow.failed`;
- `workflow.cancelled`.

The final vocabulary should avoid duplicating child run events. Workflow events should describe orchestration transitions, while run events remain the source of truth for the work itself.

### Failure and recovery tests

- failure in the first, middle, and final step;
- cancellation between steps;
- cancellation during a child execution;
- daemon or process loss after a step finishes but before state is recorded;
- uncertain external side effect;
- restart with an interrupted active step;
- invalid workflow definition;
- missing step implementation;
- incompatible workflow version;
- duplicated resume request;
- human approval timeout, when introduced.

### Out of scope initially

- unrestricted DAGs;
- arbitrary cyclic workflows;
- visual workflow builder;
- distributed workflow scheduling;
- exactly-once claims;
- BPMN compatibility;
- a plugin marketplace;
- dynamic mutation of running workflow structure;
- multi-agent organisations.

### Exit gate

Workflows are justified only when Aren can express repeated real processes more clearly and reliably than direct composition code, while keeping execution semantics, state transitions, and recovery behaviour inspectable.

---

## Phase 17 — Broader Execution Types

### Primary question

> Which additional forms of work genuinely benefit from Aren’s execution semantics?

Candidates should be introduced one at a time:

- shell processes;
- coding-agent subprocesses;
- HTTP operations;
- MCP tools;
- persistent language workers;
- deterministic jobs;
- human approvals;
- scheduled tasks.

Each candidate must answer:

- What does cancellation mean?
- What constitutes success?
- What output is produced?
- Is retry safe?
- What progress can be observed?
- What state is durable?
- What permissions are required?
- Can it survive runtime or client restart?

This is where the hypothesis that “execution is the primitive” is genuinely tested.

Some execution types may not fit cleanly. Aren should revise the abstraction rather than force false uniformity.

### Exit gate

At least two materially different execution types share useful lifecycle and observation semantics without accumulating special cases that invalidate the core model.

---

## Phase 18 — Distributed and Remote Execution

### Primary question

> Is there enough demonstrated need to execute work across machines or workers?

Possible capabilities include:

- worker registration;
- capability discovery;
- leases and heartbeats;
- remote cancellation;
- artefact transfer;
- placement decisions;
- worker-loss handling;
- duplicate execution protection;
- queueing;
- horizontal scaling.

This phase is speculative and should not be treated as inevitable.

### Exit gate

Remote execution solves a measured operational problem that cannot be addressed adequately by local execution or a single daemon.

---

## 4. Cross-Phase Capability Map

| Capability | Introduced | Deepened later |
|---|---:|---|
| Run identity | Phase 1 | Persistence, daemon, workflows |
| Lifecycle states | Phase 1 | Durable interruption and workflows |
| Cancellation | Phase 1 | Providers, tools, daemon, workflows, remote execution |
| Ordered events | Phase 1 | Streaming, persistence, reconnectable clients |
| Progress | Phase 2 | Models, tools, workflows |
| Failure classification | Phase 1 | Provider, validation, tool, workflow taxonomies |
| Model invocation | Phase 3 | Streaming, tools, context |
| Streaming | Phase 4 | Daemon reconnection and clients |
| Structured output | Phase 5 | Tool schemas and workflow inputs |
| Attempts | Phase 6 | Durable retry and policy budgets |
| Tools | Phases 7–8 | Client-hosted, MCP, remote tools |
| Agent loop | Phase 9 | Context and workflow steps |
| Context management | Phase 10 | Retrieval and memory, if earned |
| Persistence | Phase 11 | Recovery, daemon, workflows |
| Composition | Phase 12 | Workflow foundations |
| Policy | Phase 13 | Organisation-wide governance |
| Daemon | Phase 14 | Remote management and scheduling |
| Multi-language support | Phase 15 | Wider client ecosystem |
| Workflows | Phase 16 | Approvals, schedules, durable orchestration |
| Broader execution types | Phase 17 | Workflow step diversity |
| Distributed operation | Phase 18 | Only if justified |

---

## 5. Explicit Non-Commitments

This roadmap does not commit Aren to implementing every listed phase.

Aren is not yet committed to:

- a generic provider abstraction;
- a workflow graph engine;
- distributed execution;
- a plugin ecosystem;
- a general policy language;
- multi-agent coordination;
- vector-based memory;
- a database-backed architecture;
- durable mid-request resumption;
- exactly-once execution;
- every language SDK;
- every proposed execution type.

Each later phase must be approved from evidence produced by previous phases.

---

## 6. Immediate Planning Boundary

The Phase 1 PRD should cover only:

- run identity;
- lifecycle states and transitions;
- terminal outcomes;
- cancellation;
- lifecycle events;
- event observation semantics;
- minimum failure representation;
- deterministic and race testing;
- diagnostic CLI;
- explicit exclusions.

Phase 2 may be considered only enough to ensure Phase 1 is testable. It must not cause Phase 1 to include:

- provider types;
- model abstractions;
- tool interfaces;
- persistence;
- daemon protocols;
- workflow composition;
- generic execution metadata.

The immediate sequence is:

```text
Approve roadmap direction
        ↓
Write Phase 1 PRD
        ↓
Write Phase 1 lifecycle design
        ↓
Write Phase 1 testing and failure design
        ↓
Implement
        ↓
Use in a small real scenario
        ↓
Review and revise Phase 2
```

---

## 7. Final Principle

Aren should not try to prove at inception that execution is a universal primitive.

It should first prove that Aren can own one execution honestly.

Then it should test whether those semantics survive:

- controlled failure;
- real model requests;
- streaming;
- validation;
- retries;
- tools;
- agent loops;
- context growth;
- persistence;
- composition;
- daemon ownership;
- workflows;
- broader forms of work.

If the abstraction survives that progression, it will have earned its generality.

If it does not, Aren should change the abstraction rather than protect the original vision.
<<< END ULTRAPLAN DIRECT PROJECT INPUT >>>

<<< BEGIN ULTRAPLAN DIRECT PROJECT INPUT >>>
ID: project-source-aren-final-language-decision
Kind: project-source-document
Path: projects/aren-phase-01-execution-lifecycle/docs/final-language-decision.md
Assignment: Catalogued project source document: Aren Final Language Decision
Mode: full
Original-Bytes: 45558
Injected-Bytes: 45558

# Final implementation-language decision: Go

> **Status:** Accepted
>
> **Decision date:** 26 August 2026
>
> **Decision:** Aren's production runtime will be implemented in Go.
>
> **Scope:** The complete long-term runtime, including lifecycle control, model execution, tools, agent loops, persistence, daemon hosting, workflows, and any later distributed execution that is justified.
>
> **Evidence state:** Repository revision `ffa347e3dba6` plus the uncommitted prototype and documentation work described below.

## 1. Decision

Aren will use Go as its implementation language.

This is not a temporary choice for Phase 1 and it is not a plan to prototype in Go before rewriting in Rust. Go is the default language for the runtime across the roadmap. The production implementation will not maintain a shadow Rust runtime or require feature parity between two cores.

The Rust prototype remains valuable evidence. It demonstrated that Aren's semantics are compatible with Rust, exposed weaknesses in the Go design, and measured a real memory advantage. It now becomes a frozen comparison instrument rather than a second product implementation.

This decision applies to the runtime, not every program around it. Thin client SDKs may use their host language. External tools may remain separate processes. A future component may use another language behind a narrow protocol boundary if that component has an independent requirement. None of those cases reopens the runtime language by default.

## 2. Basis of the decision

The deciding question is:

> Which language best supports Aren's documented requirements and its method of developing through small, observable, failure-tested vertical slices?

The answer is Go because:

1. Both languages implemented the tested lifecycle and later roadmap slices without changing Aren's core contract. Rust is not required for capability.
2. Aren's defining correctness properties are runtime semantic properties. Neither compiler decided cancellation order, terminal arbitration, event order, cleanup completion, retry safety, or external effect truth.
3. Go kept the state machine and operating-system interactions more direct in the prototypes. That matters because Aren must be inspectable while its semantics are still evolving.
4. Go required less implementation and dependency machinery in the measured prototypes.
5. Go's standard library covers the current local library, CLI, HTTP, JSON, cancellation, process, profiling, and daemon needs.
6. Rust's memory advantage is real, but no current Aren workload or budget makes Go unacceptable. The strongest daemon measurements are preliminary and do not include the real agent process groups likely to dominate memory.
7. Choosing Rust now would pay certain complexity for possible future constraints. That conflicts with Aren's rule that complexity must be earned by demonstrated pressure.

The decision is therefore not “Go is always better than Rust.” It is that Go has the better cost structure for Aren as specified.

## 3. Aren requirements that govern the choice

Aren is not a generic service and the language was not selected from a general feature comparison. The following requirements define the relevant cost function.

| Aren requirement | Source | How Go supports it | Obligation created by choosing Go |
|---|---|---|---|
| One coherent terminal outcome | [Phase 1 lifecycle requirements](../phase-1-prd/02-lifecycle-requirements.md) | A mutex-protected transition gate can commit state, timing, event, and outcome in ordinary control flow. | All terminal paths must pass through one private commit boundary. No caller or executor may publish terminal state directly. |
| Truthful cancellation | [Phase 1 cancellation contract](../phase-1-prd/02-lifecycle-requirements.md#18-cancellation-requirements) | `context.Context` expresses a cooperative stop signal and is accepted by Go's HTTP and process APIs. | A request to cancel must remain distinct from confirmed termination. Every owned goroutine and external operation needs an explicit cancellation and join policy. |
| Ordered, observer-independent events | [Phase 1 event requirements](../phase-1-prd/02-lifecycle-requirements.md) | Private histories, sequence allocation under the transition gate, and pull-based cursors are direct to implement. | Observers must never execute inside the commit path. High-volume deltas, logs, and progress require explicit bounds rather than unbounded channels or histories. |
| Honest failure classification | [Phase 1 failure requirements](../phase-1-prd/02-lifecycle-requirements.md#17-failure-requirements) | Go error wrapping preserves inspectable causes while typed Aren failures distinguish origin and kind. | Public failure types, wrapping rules, and serialization must be defined centrally. Errors must not collapse into strings at durable or protocol boundaries. |
| Small, revisable vertical slices | [Roadmap development rules](../phased-roadmap.md#2-development-rules) | Go supports concrete implementations with little framework setup and allows abstractions to appear after repeated use. | Interfaces must be justified by actual variation. Standard-library use is preferred where it keeps behaviour clear, but dependencies are allowed when they remove proven risk. |
| Failure testing and real use | [Roadmap failure discipline](../phased-roadmap.md#24-failure-testing-is-mandatory) | The race detector, fuzzing, table-driven tests, subprocess tests, and black-box HTTP tests fit normal Go workflows. | Race-enabled CI is mandatory. Green unit tests are insufficient without adversarial, differential, leak, crash, and real-workload evidence. |
| Low, bounded, explainable overhead | [Performance engineering](../performance-engineering.md) | Go exposes runtime, heap, goroutine, mutex, block, trace, and profile data with standard tooling. | GC, allocation, runtime-task growth, retention, tail latency, saturation, and recovery must be measured throughout the roadmap. |
| Permission before effect | [Roadmap Phase 13](../phased-roadmap.md#phase-13--policies-and-resource-governance) | A direct call path can validate, authorize, reserve, and then execute without a separate policy framework. | Old or lower-level entry points must not bypass policy. Denial and admission decisions must be observable before the governed operation begins. |
| Local-first execution with an earned daemon | [Roadmap Phases 14 and 15](../phased-roadmap.md#hosting-and-external-clients) | Go works naturally as an in-process library and CLI, then as a small loopback HTTP daemon using `net/http`. | The daemon must be introduced only when caller-independent lifetime, shared inventory, persistence, or resource ownership requires it. |
| Thin multi-language clients | [Roadmap Phase 15](../phased-roadmap.md#phase-15--multi-language-clients-and-client-hosted-tools) | A stable protocol keeps lifecycle semantics in Go and clients small. | Client SDKs may project the wire contract but must not reimplement lifecycle resolution, retries, or policy. |
| Persistence without false exactly-once claims | [Roadmap Phase 11](../phased-roadmap.md#phase-11--persistence-and-recovery) | Go has mature file, database, and serialization options; the language does not dictate the recovery model. | State and event commits, schema migration, interruption, uncertainty, idempotency, and effect ledgers need explicit designs and crash tests. |
| Broader and distributed execution only when justified | [Roadmap Phases 17 and 18](../phased-roadmap.md#phase-17--broader-execution-types) | Go can supervise processes, HTTP work, daemons, and remote workers through concrete adapters and protocols. | Type-specific cancellation, output, success, cleanup, and effect semantics must remain visible. Distribution cannot be added for speculative scale. |

## 4. Evidence standard

Not all prototype evidence carries the same weight.

| Evidence | Weight in this decision | Reason |
|---|---|---|
| Mirrored lifecycle, failure, cancellation, and race cases | High for the named behaviours | They repeatedly found real defects and exercised the same semantic contract in both languages. |
| Current Go race-enabled suite and Rust safe-code compilation | High for what each tool checks | They confirm the present tested paths, while answering different concurrency questions. |
| Shared 14-case daemon contract | High for those contract cases | The same black-box suite runs unchanged against both daemons. |
| Defects discovered during extension and audit | Medium to high | They reveal concrete failure classes, but author sequence and familiarity prevent a language defect-rate comparison. |
| Current implementation size and dependency shape | Medium | They are useful maintenance indicators, not direct measures of comprehension or quality. |
| Repeated blocked-run footprint | Medium | It demonstrates a substantial memory difference in those two designs, but not full agent capacity. |
| Saved shared-daemon memory results | Preliminary | They used one process per cell on a loaded machine and predate scheduler and accounting repairs. |
| Initial microbenchmarks | Low for the current runtime | They describe an early Phase 1 shape and do not predict model, tool, or agent latency. |
| JavaScript daemon versus embedded Rust benchmark | Medium for deployment architecture only | It changed process boundary, lifetime, data movement, and language at once. |
| Future scale, FFI, WASM, or clustering arguments | None without a requirement | The roadmap explicitly treats later capabilities as hypotheses. |

This weighting prevents one favourable number or one type-system property from overriding the actual product requirements.

## 5. What the prototype programme established

The Go and Rust prototypes were disposable decision instruments. They were deliberately extended far beyond Phase 1 to test whether later capabilities would force a language-specific redesign.

### 5.1 Coverage and remaining limits

| Experiment | Tested capability | Evidence produced | Important limit |
|---|---|---|---|
| [Initial comparison](report.html) | First lifecycle implementation and microbenchmarks | Both languages expressed the basic run lifecycle; Go had lower full-run latency in the early shape and Rust cloned event history faster. | Early benchmarks used different frameworks and predate most runtime behaviour. |
| [Experiments 2 and 3](phase-1-experiment-2-cancellation-and-retention.md) | Cancellation, retention, observers, invariant injection, and stress | Both reached one-terminal, ordered-history, independent-observer, multi-waiter, and collision behaviour. | Hand-written schedules do not exhaust possible interleavings. |
| [Experiment 4](phase-2-experiment-4-progress-cleanup-and-conformance.md) | Progress, cleanup, partial output, and executor conformance | Cleanup occurs before terminal publication and cleanup failure remains visible. | Progress retention is unbounded and cleanup has no separately enforced deadline. |
| [Experiment 5](phase-3-experiment-5-first-model-invocation.md) | Cancellable model call and provider failures | Both separated provider transport, protocol, malformed response, rate limit, timeout, and cancellation. | Only Go has historical live-endpoint evidence; the current smoke was not rerun for this decision. |
| [Experiment 6](phase-4-experiment-6-streaming-model-execution.md) | Ordered model deltas and interruption | Both serialized callbacks, reconstructed output exactly, bounded retained deltas, and kept partial interrupted text diagnostic. | Go's default SSE adapter buffers the response; Rust has no provider network-stream adapter. Neither proves upstream backpressure. |
| [Experiment 7](phase-5-experiment-7-structured-results.md) | Structured result parsing and validation | Both distinguished provider success from task validity and preserved raw invalid output. | No provider-enforced schema, nested production schema, or repair loop. |
| [Experiment 8](phase-6-experiment-8-retry-and-backoff.md) | Attempts, selective retry, backoff, and exhaustion | Both implemented bounded retry policy, cancellable backoff, `Retry-After`, and explicit exhaustion. | No effect idempotency or complete durable attempt record. |
| [Experiment 9](phase-7-experiment-9-tool-call-representation.md) | Tool definitions and requested calls | Both preserved precise raw arguments, strict structure, order, IDs, and inspectable unknown tools. | No streamed tool arguments or provider interoperability matrix. |
| [Experiment 10](phase-8-experiment-10-local-tool-execution.md) | Tool resolution, validation, permission, and local files | Both revalidated arguments and denied before entering implementations. | No OS sandbox, forced stop, output bound, or TOCTOU-resistant filesystem boundary. |
| [Experiment 11](phase-9-experiment-11-bounded-agent-loop.md) | Minimal model and tool loop | Both enforced turn, tool-call, identity, and elapsed limits around a sequential transcript. | No real multi-turn task, token or cost budget, durable transcript, or parallel tool execution. |
| [Experiment 12](phase-10-experiment-12-context-management.md) | Context measurement and compaction | Both measured before model calls and compacted complete non-durable turns deterministically. | Byte limits stand in for token budgets; no evidence store, summarization, or provider-cache measurement. |
| [Experiment 13](phase-11-experiment-13-persistence-and-recovery.md) | Run snapshots, histories, and interruption | Both retained monotonic in-memory records and classified simulated loss honestly as interrupted. | No process-restart durability, fsync, schema evolution, or effect recovery. |
| [Experiment 14](phase-12-experiment-14-execution-composition.md) | Sequence and bounded parallel groups | Both reused child runs with parent cancellation and first-failure semantics. | No durable child evidence, compensation, fairness, or recovery. |
| [Experiment 15](phase-13-experiment-15-policies-and-resources.md) | Model and tool allowlists | Both denied observably before the tested effect. | Policy can be bypassed through older paths; core resource accounting was not part of this slice. |
| [Experiment 16](phase-14-experiment-16-daemon-and-clients.md) | Go daemon, thin clients, and embedded Rust comparison | The daemon established caller-independent ownership; the embedded path demonstrated a different lifetime and lower in-process call overhead. | The performance comparison was architecture plus language, not language alone. Neither path proved production recovery or security. |
| [Shared-daemon gauntlet](phase-16-experiment-17-shared-daemon-gauntlet.md) | One scheduler, admission, resource, retention, and HTTP contract in both languages | Both pass the same 14 black-box cases. Rust used less RSS in every saved matched cell. | Saved memory files predate three correctness repairs and need repeated runs on a quiet reference machine. |
| [Workflow experiment](phase-16-experiment-17-workflows.md) | Definition versus instance, sequential steps, children, and interruption | Both carried lifecycle semantics into a small workflow without a new core. | In-memory store, no daemon hosting, versioning, resume, human wait, or uncertain-effect recovery. |
| [Experiment 18](phase-17-experiment-18-broader-execution-types.md) | Shell and HTTP executions | Both reused the lifecycle and retained type-specific results and failures. | Neither guarantees descendant cleanup, sandboxing, consistent output bounds, streaming, or restart handling. |
| [Experiment 19](phase-18-experiment-19-distributed-and-remote.md) | Remote execution seam | Both crossed the existing daemon contract with failure classification and post-acknowledgement cancellation intact. | No workers, leases, heartbeats, placement, artifact transfer, fencing, or duplicate-effect protection. The need remains unproven. |

The breadth establishes feasibility and architectural continuity. It does not make the later prototype slices production-ready. Persistence is in memory, workflows are linear, policy is narrow, and distribution is only a seam.

### 5.2 Core semantic result

Both implementations preserved the same central lifecycle model through all tested extensions:

- one Aren-owned transition gate;
- exactly one terminal state, event, outcome, and timing record;
- cancellation acceptance distinct from confirmed stop;
- ordered observation that does not control execution;
- classified work, Aren, provider, validation, tool, workflow, shell, HTTP, and remote failures;
- explicit attempts and bounded loops;
- child execution using the same run semantics;
- policy decisions before the tested effects.

No roadmap slice required Rust's ownership model, and no slice exposed a capability Go could not implement directly.

That does not prove the lifecycle abstraction is complete. New executors continued to require type-specific cancellation, output, success, and cleanup rules. This supports the roadmap's plan to broaden through evidence rather than force a universal executor interface.

## 6. What the defects teach

The defect history is more valuable than a count of green tests.

### 6.1 Defects found in both designs

Both languages admitted semantic failures such as:

- cancellation bypassing Aren's acceptance path;
- terminal or progress facts observed in the wrong order;
- completed runs retained by watchers or elapsed-budget tasks;
- stream event order diverging from reconstructed output order;
- permits or reservations not released on early cancellation;
- HTTP cancellation stopping at headers rather than body reads;
- subprocess pipes capable of deadlocking when output was not drained;
- false-positive tests that did not exercise the path they claimed to prove.

These are central Aren risks. They were controlled by the lifecycle contract, a single mutation boundary, explicit ownership, and adversarial tests, not by the language alone.

### 6.2 Rust's protection was real but partial

Rust improved several properties:

- closed state, event, and failure vocabularies use enums and exhaustive matching;
- immutable shared outcomes are natural through ownership and `Arc`;
- `Send` and `Sync` constrain values crossing task boundaries;
- many unsafe sharing errors fail to compile;
- there is no tracing garbage collector.

Rust did not prevent logical deadlocks, task retention, incorrect event order, premature cancellation claims, cancellation-unsafe operations, semaphore leaks, or external-process cleanup defects. Aren's hardest guarantees remained runtime proofs.

### 6.3 Go's weaknesses are accepted, not ignored

The Go prototype exposed defects that Rust's types made less likely, including a mutable outcome shared across waiters and mutable schemas or maps retained across asynchronous execution. Go also admitted a data race in a cleanup window, which the race detector found once the test reached it.

Choosing Go accepts that immutable publication, closed vocabularies, and safe sharing rely more heavily on API design and verification. The engineering rules in section 11 are therefore part of the decision, not optional follow-up work.

The experiment history cannot support a defect-rate comparison. The same author built and repaired the implementations sequentially, and later work benefited from earlier findings.

## 7. Specific weaknesses uncovered and required responses

The table separates the prototype in which a weakness was observed from the lesson Aren must carry into the production Go runtime. Rust-specific findings remain relevant because they identify semantic failures that Go's compiler will not prevent either.

| Weakness uncovered | Observed in | Consequence demonstrated or exposed | Required response in the Go runtime |
|---|---|---|---|
| Parent cancellation bypassed Aren's cancellation-acceptance path. | Go and Rust | A run could receive a stop signal without retaining the first reason or recording the required request event through the authoritative path. | Route explicit, parent, deadline, daemon, workflow, and remote cancellation through one idempotent acceptance function before signalling work. Test every source against the same disposition and event contract. |
| A terminal cancellation event could be recorded before `run.cancellation_requested`. | Rust | History could contradict the causal order promised by Aren even though the final state was plausible. | Commit cancellation acceptance and its event under the transition gate before exposing the signal. Assert causal ordering under repeated completion and cancellation races. |
| Waiters shared a mutable terminal outcome. | Go | One consumer could mutate data observed by another, violating immutable publication and identical-outcome guarantees. | Keep outcome fields private, deep-copy nested collections at commit time, and return values or defensive snapshots. Add mutation attempts to waiter and observer tests. |
| The invariant-failure path recursed, double-closed a channel, and temporarily published incomplete terminal data. | Go | The error path intended to protect coherence could panic or expose state without its matching event and outcome. | Make invariant handling non-recursive. Use one terminal commit routine, one close owner, and a minimal fallback that cannot re-enter ordinary transition logic. Exercise it through explicit fault injection. |
| Terminal outcomes were published without terminal timing. | Go | Consumers could observe a completed run whose outcome violated the Phase 1 terminal record contract. | Construct terminal state, timing, event, and outcome together inside the commit boundary. Reject any terminal constructor input without complete timing. |
| Cleanup introduced an unlock window and a data race before final publication. | Go | Concurrent inspection could observe torn lifecycle facts while cleanup and terminal resolution competed. | Give one owner the classified work result, cleanup result, and final commit. Do not expose an intermediate terminal state. Keep the full suite under the race detector and stress inspection during cleanup. |
| Schemas, argument maps, and work slices were retained after callers could mutate them. | Go | Asynchronous execution could observe a definition different from the one validated at submission. | Validate and deep-copy all mutable caller-owned inputs before acceptance. Store typed immutable forms internally and test caller mutation immediately after submission. |
| Default JSON number decoding lost integer precision. | Go | Tool and structured-result arguments could change value while still looking syntactically valid. | Decode with `json.Decoder.UseNumber`, validate numeric ranges explicitly, reject trailing values, and preserve raw arguments alongside the validated form. |
| Classified Aren failures were flattened into ordinary errors on some paths. | Go | Failure origin and kind could be lost, weakening retry, diagnosis, and protocol behaviour. | Use one typed failure representation with wrapped local causes. Centralize conversion at provider, tool, workflow, persistence, and wire boundaries. Test origin, kind, message, and cause separately. |
| Cancellation during an HTTP body read was classified as an HTTP failure. | Go | The terminal record could lie about why work ended even though the context had been cancelled. | Check context cause around request dispatch, streaming, and body reads before assigning a transport classification. Test cancellation before headers, during headers, and during body consumption. |
| Completed runs were retained by a watch receiver or elapsed-budget task. | Rust | Terminal runs remained reachable after clients and work had finished. | Give observers, timers, and budget goroutines explicit release paths. Cancel and join them on terminal commit, then verify reclamation with leak and retention tests. |
| Watch-channel liveness and re-locking a non-reentrant mutex caused deadlocks. | Rust | Safe memory access did not guarantee forward progress or a valid lock protocol. | Define lock ordering, keep blocking calls and callbacks outside locks, avoid nested acquisition of the transition lock, and add timeout-backed deadlock tests plus model-based schedule exploration. |
| Stream callback order and reconstructed output order could diverge. | Rust | Event evidence and the final result could describe different executions. | Serialize every accepted delta through one ordering owner, assign sequence numbers there, and assert that retained deltas reconstruct the exact final output under concurrent callbacks. |
| A semaphore permit was not released on an early parallel-cancellation path. | Rust | Bounded composition could lose capacity permanently and eventually stall. | Represent every admission and reservation with a scoped release token. Test success, failure, panic, cancellation, and partial-start exits for exact resource reconciliation. |
| Provider streaming was simulated but not live and incrementally backpressured. | Go and Rust | Passing callback tests did not prove network streaming, slow-consumer behaviour, bounded buffering, or cancellation of an active provider body. | Build a real incremental SSE or provider adapter in Go with bounded buffers. Test fragmented frames, slow consumers, malformed late frames, cancellation, upstream closure, and retained-byte limits. |
| Subprocess output could deadlock, remain unbounded, or leave descendants running. | Go and Rust | A direct child could fill a pipe, consume unbounded memory, or terminate while grandchildren survived. | Drain stdout and stderr concurrently, enforce explicit byte limits, and supervise process groups on Unix and Job Objects on Windows. Test infinite output, ignored signals, forked descendants, and cleanup deadlines. |
| Progress histories, output, and daemon tombstones could grow without a lifetime bound. | Go and Rust | Long-lived processes accumulated memory even after heavyweight run data was evicted. | Separate durable lifecycle facts from high-volume data, define per-run and process-wide limits, expire or externalize tombstones, and prove recovery toward baseline in churn and soak tests. |
| Policy could be bypassed through older entry points. | Go and Rust | A correct allowlist on one path did not create an authoritative security or resource boundary. | Make validation, policy, admission, reservation, and execution one unbypassable orchestration path. Keep lower-level executors private or require an already-authorized capability token. |
| Several passing tests contained impossible assertions, empty helpers, or counters unrelated to the behaviour under test. | Go and Rust | Green mirrored suites overstated evidence and allowed defects to survive. | Add negative controls, seeded historical defects, mutation testing where practical, and a language-neutral reference model. Prefer requirement-to-test mapping over cumulative test counts. |

These responses are incorporated into the mandatory Go rules in section 11. They should also seed the production test plan so the prototype's failures remain regression cases rather than historical notes.

## 8. Performance and resource evidence

### 8.1 Latency

The original full-run benchmarks favoured Go for the early lifecycle shape, while Rust favoured event-history cloning. Corrected completed-wait measurements were effectively tied at about 40 nanoseconds.

Those operations are not representative of an agent runtime dominated by provider, network, process, storage, and human waits. They remain regression baselines, not a language decision.

### 8.2 Active-run memory

The repeated 10,000 blocked-run probe is the strongest result in Rust's favour:

- Go mean RSS increase per active run: about 9,604 bytes;
- Rust mean RSS increase per active run: about 2,155 bytes;
- observed ratio: about 4.5 to 1 in favour of the tested Rust design.

The designs were not carrier-identical. Go used three goroutines per run and Rust used two Tokio tasks. RSS also includes stack policies, allocators, schedulers, and runtime state. The result proves that this Rust design was materially lighter, not that every Rust Aren design will use 4.5 times less memory.

No documented Aren workload or budget makes the Go result unacceptable. At 10,000 blocked lifecycle runs the difference is roughly 71 MiB. A real coding-agent run may retain transcript, tool, provider, artifact, and child-process data that outweighs the lifecycle carrier.

### 8.3 Shared daemon

The saved shared-daemon runs reported Rust using 33 to 54 percent less RSS across matched cells. At two simulated 20 MiB runs Rust saved about 43 MiB. At eight, the recorded peaks were about 318 MiB for Go and 206 MiB for Rust.

Resource policy had a larger absolute effect in the same saved data. Reducing the active limit from eight to two saved about 230 MiB in Go and 161 MiB in Rust. This supports Aren's requirement for bounded admission regardless of language.

The shared results are not decision-grade performance measurements. They contain one daemon process per cell, came from a machine with active swap, sampled short workloads coarsely, lack binary digests, and predate repairs to scheduler weighting, reservation accounting, and overflow handling. They establish the direction of Rust's advantage, not its production value.

### 8.4 Performance conclusion

The decision does not deny Rust's memory result. It judges that the known benefit does not justify Rust's continuing complexity without a breached Aren capacity or latency requirement.

Go must now prove its performance against real workload budgets. If Go later misses a hard budget, the first response is measurement, retention control, admission, profiling, and focused design. The answer is not an automatic rewrite.

## 9. Implementation and operational cost

The current prototype shapes provide useful maintenance indicators:

| Component | Go | Rust | Observation |
|---|---:|---:|---|
| Core implementation | 3,567 lines | 4,425 lines | Rust is about 24% larger in this implementation. |
| Shared daemon implementation | 1,532 lines | 1,626 lines | Source sizes are close; the dependency models differ. |
| Diagnostic CLI | 424 lines | 694 lines | Rust is about 64% larger in this implementation. |

The counts exclude tests and generated build output. They are not a measure of quality or developer intelligence. Rust syntax, enums, conversions, and explicit ownership naturally use more code. Go can hide obligations in conventions. The result that matters is that Rust's extra machinery remained present as the prototype expanded.

The Go prototype declares no external module. The Rust prototype declares Tokio, Tokio Utilities, Futures, Serde, Serde JSON, Reqwest, and Axum, resolving a much larger transitive graph. Those are mature libraries that provide real value. They also introduce upgrade, advisory, feature, compile-time, and API-lifecycle work before Aren has demonstrated a need for that stack.

Go's standard library made several relevant paths direct:

- HTTP server and client behaviour;
- JSON request and response handling;
- context propagation into network calls;
- direct-child process cancellation;
- profiling and runtime inspection;
- race-enabled test execution.

Go does not remove the need for careful design. Its advantage is that more of Aren's state machine remains ordinary, reviewable control flow while the product is still learning what belongs in that state machine.

## 10. Alternatives considered

### 10.1 Rust for the permanent runtime

Rust was the only serious alternative. It offers stronger static representation and lower memory in the measured designs. It would be the better choice if Aren already had a hard no-GC requirement, a mandatory native-embedding requirement across several hosts, strict memory density that Go failed, or a team and codebase whose Rust async model made the lifecycle clearer rather than less clear.

Those are not current requirements. Selecting Rust would therefore front-load complexity for predicted pressure while leaving Aren's defining runtime proofs in tests and design.

### 10.2 Go now, Rust later

Rejected. This would make early product knowledge disposable, create two histories, encourage architecture for a port rather than for users, and defer the actual decision.

### 10.3 Two production cores

Rejected. Mirrored prototypes were useful for experiments. Mirrored production runtimes would double implementation, review, conformance, security, and operational work while making neither implementation authoritative.

### 10.4 A split Rust core with a Go daemon

Rejected as the default. It adds an FFI or internal protocol boundary inside the most correctness-sensitive path without a demonstrated requirement. A separately useful Rust component may be proposed later behind a narrow contract, but lifecycle authority remains in Go.

### 10.5 Delay the decision for more benchmarks

Rejected. The missing tests are important engineering work, but no current result identifies a capability Go lacks or a requirement it fails. Keeping the language open would itself carry cost: duplicated prototypes, delayed repository foundation, and continued re-litigation based on hypothetical futures.

## 11. Mandatory engineering rules for the Go runtime

The following rules are consequences of the decision.

### 11.1 One owner for lifecycle mutation

- Run fields are private.
- Every lifecycle transition passes through one commit function or one clearly defined ownership loop.
- State, event sequence, timing, outcome, and waiter release become visible coherently.
- Executors return facts to the owner. They do not mutate run state.

### 11.2 Published data is immutable by contract

- Outcomes, events, schemas, tool arguments, work definitions, and nested collections are copied or frozen before asynchronous retention or publication.
- Accessors return values or defensive snapshots, not mutable internal maps, slices, pointers, or `any` containers.
- Durable formats use explicit types and versions.

### 11.3 Goroutines have owners

- Every Aren-owned goroutine has a named owner, cancellation source, termination condition, and join or release path.
- `context.Background()` is permitted only at an explicit lifetime boundary such as daemon-owned accepted work.
- Fire-and-forget work is prohibited unless its bounded lifetime and irrelevance to correctness are documented.
- Leak tests cover completed runs, abandoned observers, timers, blocked providers, tools, and subprocesses.

### 11.4 Cancellation remains a protocol

- Context cancellation is a signal, not proof of stop.
- The first accepted reason is retained and repeated requests are idempotent.
- Terminal cancellation is published only after the executor-specific stop condition is satisfied.
- Process-tree, HTTP-body, tool, provider, workflow, and remote cancellation receive separate conformance cases.

### 11.5 Go's open type system is narrowed deliberately

- Lifecycle states, event types, failure origins, and failure kinds use named private or closed packages with validated constructors.
- Conversion and serialization switches are centralized.
- CI uses an exhaustive-switch analyzer for closed vocabularies where practical.
- Unknown protocol values remain explicit rather than becoming zero values.

### 11.6 Concurrency testing is part of correctness

- The complete Go suite runs under `go test -race` in CI.
- Stress tests repeat completion and cancellation collisions, concurrent inspection, observers, cleanup, and resource release.
- Fuzzed or model-based command traces compare actual state and history with a small reference model.
- Tests prove their instruments can fail through seeded historical defects, mutation testing, or deliberately broken fixtures.

### 11.7 Work and retention are bounded

- Every queue, event class, stream, tool output, transcript, artifact, cache, retry loop, and tombstone policy declares a bound or a measured reason not to have one.
- Admission reserves resources before effects and releases them on every path.
- Slow observers cannot block execution.
- Resource estimates and measured use remain distinct and overflow-safe.

### 11.8 Performance is measured at product scale

- Each phase records incremental allocations, heap or RSS, goroutines, contention, p50/p95/p99 latency where relevant, saturation, and recovery.
- Benchmarks state revision, diff status, binary digest, toolchain, hardware, workload, repetitions, and spread.
- Profiles precede optimization.
- Real process-group memory is measured when child agents exist; daemon RSS alone is not treated as capacity.

### 11.9 Errors and effects stay honest

- Aren failure origin and kind remain machine-readable.
- Wrapped causes remain inspectable locally; durable and wire representations state what is lost.
- Retry eligibility is explicit and bounded.
- External effects are never described as exactly once without an idempotency or effect mechanism.
- Recovery distinguishes completed, interrupted, and uncertain work.

### 11.10 Dependencies and abstractions must earn their cost

- Standard-library code is preferred when it is clear and sufficient.
- A dependency is accepted when it removes a demonstrated risk or substantial maintained code, not to satisfy a speculative architecture.
- Interfaces follow real variation or a required test boundary.
- Phase reviews may collapse packages, remove options, or delete abstractions.

## 12. Known gaps after the decision

The language decision does not close the following engineering questions:

1. Neither prototype has executed a real, repeated, multi-turn agent workload combining a live provider stream, tools, context growth, child processes, persistence, cancellation, and reconnect.
2. Neither has production persistence or crash consistency.
3. Neither has a complete tool sandbox, process-tree guarantee, consistent output bound, or cross-platform containment proof.
4. Neither has demonstrated live incremental provider streaming with backpressure and provider-enforced structured output.
5. The shared-daemon memory experiment needs corrected, repeated measurements on a quiet reference machine with real process groups.
6. Scheduler weights are contract-tested as configuration but not yet proved under long saturated admission with starvation bounds.
7. Tombstones and several high-volume histories need durable, bounded retention policies.
8. Security work remains incomplete: daemon authentication, identity-scoped policy, filesystem races, secret handling, and hostile protocol input.
9. Test totals and some historical summaries have drifted. Requirement-to-test mapping should replace cumulative counts.
10. Linux evidence dominates. macOS, Windows, ARM64, packaging, upgrade, and rollback need explicit matrices.

These gaps become roadmap work in Go. They do not justify maintaining two language implementations.

## 13. Reversal threshold

This is a definitive decision, not an irreversible claim about the universe. There is no scheduled language checkpoint and no standing plan to reconsider Rust after a particular phase.

A rewrite or replacement may be proposed only if all of the following are true:

1. Aren acquires a hard production requirement, expressed as a measurable threshold or mandatory integration contract.
2. The current Go runtime fails that requirement on a representative workload after profiling and proportionate design work.
3. The failure is attributable to a constraint that a different language can materially change, rather than to retention, admission, algorithms, protocol, storage, or unbounded work.
4. A competing implementation demonstrates the requirement under comparable architecture, workload, hardware, and operational conditions.
5. The migration cost and new operational complexity are lower than the cost of addressing the problem within Go.

Possible triggers include a proved no-GC requirement, mandatory in-process native embedding that a Go boundary cannot satisfy, or a hard laptop-capacity target that corrected process-group measurements show Go cannot meet. Speculative scale, language fashion, isolated microbenchmarks, or a favourable percentage without a product budget are not triggers.

## 14. Consequences

### Positive

- Aren can start one production codebase without a planned port.
- The lifecycle state machine can remain visible in direct control flow.
- The local library, CLI, HTTP transport, subprocess paths, and diagnostics can begin with the standard library.
- Race-enabled CI provides a strong dynamic check over executed Go paths.
- Multi-language consumers can remain thin clients of an earned daemon rather than forcing native bindings into the core.
- The team can spend complexity on lifecycle truth, failure evidence, containment, recovery, and product use instead of async framework integration.

### Costs accepted

- Go does not make published data immutable or lifecycle matches exhaustive by default.
- The garbage collector and goroutine stacks increase memory use and may add latency variance.
- Data races are possible and the race detector observes only executed paths.
- Goroutine leaks, channel deadlocks, and logical ordering errors remain possible.
- Some invalid combinations require constructors, private state, analyzers, and tests rather than compiler ownership rules.
- Native embedding into non-Go hosts is less direct than Rust's N-API path.

The mandatory engineering rules are the mitigation for these costs. If they are not adopted, the rationale for choosing Go no longer holds.

## 15. Current verification

The following checks passed on 26 August 2026 against the current working tree:

```text
Go:
  go vet ./...
  go test -race -count=1 ./...

Rust:
  cargo test --all-targets
  cargo clippy --all-targets -- -D warnings

Shared daemon contract:
  Go:   14 passed, 0 failed
  Rust: 14 passed, 0 failed
```

These checks verify the present test and contract state. They do not reproduce the historical performance measurements or close the production gaps listed above. Published cumulative test totals are intentionally omitted because the experiment reports use inconsistent counting methods.

## 16. Decision outcome

Go is the implementation language for Aren's production runtime.

The choice is based on Aren's requirements, not on a claim that Go dominates Rust in general. Rust demonstrated better static ownership and lower memory in the prototype. Go demonstrated sufficient capability across every tested roadmap slice with a smaller, more direct implementation and lower continuing machinery.

That is the better match for Aren's governing rule: do not take on complexity until evidence shows it is required.

## 17. Evidence sources

Primary requirements and planning:

- [Project lineage and development doctrine](../project-lineage.md)
- [Phased roadmap](../phased-roadmap.md)
- [Phase 1 PRD](../phase-1-prd.md)
- [Phase 1 product definition](../phase-1-prd/01-product-definition-and-scope.md)
- [Phase 1 lifecycle requirements](../phase-1-prd/02-lifecycle-requirements.md)
- [Phase 1 validation and acceptance](../phase-1-prd/03-validation-and-acceptance.md)
- [Phase 1 risks, decisions, and handoff](../phase-1-prd/04-risks-decisions-and-handoff.md)
- [Performance engineering](../performance-engineering.md)
- [Coding-agent context engineering](../coding-agent-context-engineering.md)

Decision analysis:

- [Consolidated Rust versus Go analysis](consolidated-rust-vs-go-analysis.md)
- [Long-term Go recommendation](long-term-go-recommendation.md)
- [Critique of the Rust argument](critique-of-the-rust-argument.md)
- [Language-decision evidence index](README.md)
- [Prototype overview](../../prototypes/README.md)

Experiment reports:

- [Initial prototype report](report.html)
- [Experiment 2: cancellation and retention](phase-1-experiment-2-cancellation-and-retention.md)
- [Experiment 3: observation, invariants, and stress](phase-1-experiment-3-observation-invariants-and-stress.md)
- [Experiment 4: progress, cleanup, and conformance](phase-2-experiment-4-progress-cleanup-and-conformance.md)
- [Experiment 5: first model invocation](phase-3-experiment-5-first-model-invocation.md)
- [Experiment 6: streaming](phase-4-experiment-6-streaming-model-execution.md)
- [Experiment 7: structured results](phase-5-experiment-7-structured-results.md)
- [Experiment 8: retry and backoff](phase-6-experiment-8-retry-and-backoff.md)
- [Experiment 9: tool-call representation](phase-7-experiment-9-tool-call-representation.md)
- [Experiment 10: local tool execution](phase-8-experiment-10-local-tool-execution.md)
- [Experiment 11: bounded agent loop](phase-9-experiment-11-bounded-agent-loop.md)
- [Experiment 12: context management](phase-10-experiment-12-context-management.md)
- [Experiment 13: persistence and recovery](phase-11-experiment-13-persistence-and-recovery.md)
- [Experiment 14: execution composition](phase-12-experiment-14-execution-composition.md)
- [Experiment 15: policies and resources](phase-13-experiment-15-policies-and-resources.md)
- [Experiment 16: daemon and clients](phase-14-experiment-16-daemon-and-clients.md)
- [Shared-daemon gauntlet](phase-16-experiment-17-shared-daemon-gauntlet.md)
- [Workflow experiment](phase-16-experiment-17-workflows.md)
- [Experiment 18: broader execution types](phase-17-experiment-18-broader-execution-types.md)
- [Experiment 19: distributed and remote execution](phase-18-experiment-19-distributed-and-remote.md)
<<< END ULTRAPLAN DIRECT PROJECT INPUT >>>
