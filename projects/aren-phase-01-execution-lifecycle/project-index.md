# Project Index: Aren Phase 1 — Execution Lifecycle

> Project: `aren-phase-01-execution-lifecycle`  
> Purpose: governance, evidence, reasoning, and sprint planning for Phase 1 of Aren.

## Project Scope

- **Project Slug:** `aren-phase-01-execution-lifecycle`
- **Target Repository:** `../Aren/`
- **Expected Implementation Directory:** `/home/antonioborgerees/coding/Aren`
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
| Aren Project Lineage | `../Aren/docs/project-lineage.md` or current equivalent | Historical context from Elevate through 24-Hour Testers, AgentWrap, UltraPlan, and Aren. |
| Aren Phased Roadmap | `../Aren/docs/phased-roadmap.md` | Aren-wide phase sequence and Phase 1 boundary. |

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
| Go CLI Study | `studies/go-cli-study/` | Diagnostic CLI structure, IO seams, exit codes, testing, project layout | Supporting only |

## Available Agent-Harness Evidence Reports

Select only reports relevant to the current sprint through its `sprint-index.md`.

| Report | Path | Covers |
|---|---|---|
| Execution Model Taxonomy | `studies/agent-harness-study/reports/final/01.01-execution-model-taxonomy.md` | Units of atomic progress, execution archetypes, durability pressure. |
| Control-Flow Ownership | `studies/agent-harness-study/reports/final/01.02-control-flow-ownership.md` | Runtime authority, typed transitions, control overrides. |
| Step / Turn / Task Atomicity | `studies/agent-harness-study/reports/final/01.03-step-turn-task-atomicity.md` | Atomic units, state/event consistency, crash and retry boundaries. |
| Termination And Loop Bounds | `studies/agent-harness-study/reports/final/01.04-termination-and-loop-bounds.md` | Runtime-driven termination and structured stop reasons. |
| Pause / Resume / Interrupt | `studies/agent-harness-study/reports/final/01.05-pause-resume-interrupt-semantics.md` | Cooperative interruption, cancellation acknowledgement, durable pause distinctions. |
| Concurrency And Parallel Advancement | `studies/agent-harness-study/reports/final/01.07-concurrency-and-parallel-advancement.md` | Single-writer models, deterministic ordering, cancellation and sibling behaviour. |
| Delivery Guarantees And Idempotency | `studies/agent-harness-study/reports/final/01.09-delivery-guarantees-and-idempotency.md` | Canonical recording versus delivery, stable identities, exactly-once warnings. |
| Replay And Determinism | `studies/agent-harness-study/reports/final/01.10-replay-and-determinism.md` | Replay boundaries, deterministic reconstruction, ordering. |
| State Taxonomy And Ownership | `studies/agent-harness-study/reports/final/02.01-state-taxonomy-and-ownership.md` | State classes, owners, durable versus ephemeral boundaries. |
| Mutation Discipline And State Transitions | `studies/agent-harness-study/reports/final/02.04-mutation-discipline-and-state-transitions.md` | Single mutation entry points, transition guards, commit boundaries. |
| Completion And Finalization Semantics | `studies/agent-harness-study/reports/final/03.09-completion-and-finalization-semantics.md` | Completion causes, finalization, terminal outcome semantics. |

## Available Supporting Go CLI Reports

| Report | Path | Covers |
|---|---|---|
| Project Structure | `studies/go-cli-study/reports/final/01-project-structure.md` | Thin entrypoints and internal package boundaries. |
| Error Handling | `studies/go-cli-study/reports/final/05-error-handling.md` | Exit mapping and actionable diagnostics. |
| IO Abstraction | `studies/go-cli-study/reports/final/06-io-abstraction.md` | Injectable streams and deterministic command tests. |
| State And Context | `studies/go-cli-study/reports/final/07-state-context.md` | Go context propagation and app state. |
| Concurrency | `studies/go-cli-study/reports/final/08-concurrency.md` | Goroutine ownership, cancellation, worker safety. |
| Testing Strategy | `studies/go-cli-study/reports/final/11-testing-strategy.md` | Unit, command, integration, race, and fixture strategies. |
| Philosophy | `studies/go-cli-study/reports/final/15-philosophy.md` | Scope control and deliberate complexity. |

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
