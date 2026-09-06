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
