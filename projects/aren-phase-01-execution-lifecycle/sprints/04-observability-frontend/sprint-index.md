# Sprint index: Observability frontend

> Project: `aren-phase-01-execution-lifecycle`
> Sprint: `04-observability-frontend`
> Purpose: select the contracts, evidence, prior lifecycle and performance decisions, reasoning templates, and review protocols needed to ship Aren's first browser observation experience.
> Inputs Used: `projects/aren-phase-01-execution-lifecycle/project-index.md`, `projects/aren-phase-01-execution-lifecycle/roadmap.md`, `projects/aren-phase-01-execution-lifecycle/docs/PRD.md`, `projects/aren-phase-01-execution-lifecycle/docs/observability-mandate.md`, `projects/aren-phase-01-execution-lifecycle/docs/performance-engineering.md`, accepted project reasoning, and completed Sprint 1 through Sprint 3 artifacts.

This index is prepared in advance. Before Sprint 4 starts, confirm every named prior artifact exists and update selections affected by the realised runtime, observation, benchmark, and repository design.

## Sprint Scope

- **Sprint Goal:** Ship a dedicated browser frontend that makes live and completed Phase 1 lifecycle and performance evidence understandable without duplicating runtime semantics or affecting execution.
- **Planned Output:** A versioned read-only observation DTO, process-scoped loopback service, run explorer, ordered timeline, structured outcome and failure presentation, cancellation and diagnostic views, performance evidence views, explicit connection and schema states, accessibility, responsive behaviour, bounded delivery and rendering, and final Phase 1 review.
- **Depends On:** Completed and reviewed Sprint 1 and Sprint 2 runtime work plus the accepted Sprint 3 measurement method, result schemas, and baseline.
- **Non-Goals:** Browser control of lifecycle transitions, persistent daemon hosting, remote or public APIs, authentication, multi-user operation, durable observation storage, production telemetry export, global event streams, model or tool views, and speculative execution hierarchies.

## Source Project Index

- `projects/aren-phase-01-execution-lifecycle/project-index.md` is the authoritative catalog. Every selected contract, report, reasoning template, and protocol below appears there.

## Selected Contracts

| Contract | Why Selected |
|---|---|
| Aren Development Doctrine | Keeps the frontend a bounded Phase 1 observation slice rather than a dashboard platform, control plane, daemon, or speculative universal observability system. |
| Aren Runtime Architecture | Governs runtime-versus-projection authority, process ownership, dependency direction, thin local-service boundaries, and the rule that browser code cannot become lifecycle authority. |
| Aren Execution Lifecycle | Fixes the lifecycle, outcome, failure, timing, and cancellation meanings the browser is allowed to present; the UI may project but not reinterpret them. |
| Aren Cancellation And Lifetimes | Governs process-scoped service/client support lifetimes, cancellation, shutdown, abandonment, and proof that disconnected observers cannot strand runtime work. |
| Aren Events Observation And Waiting | Governs canonical history, sequence order, replay/live handoff, cursor semantics, terminal drain, late observation, and delivery limits exposed through the read-only projection. |
| Aren Verification | Requires interface agreement, real runtime/browser paths, controlled connection schedules, failure explainability, mutation/isolation checks, bounded-performance evidence, and review traceability. |
| Aren Observability | Primary Sprint 4 evidence contract. Governs canonical-versus-diagnostic distinction, passive observation, DTO/UI projection authority, correlation, retention, failure explainability, sensitive data, and the Phase 1 observability gate. |
| Aren Performance Engineering | Governs payload/update/rendering bounds, client/resource growth, observation overhead, saturation behaviour, and preservation of the accepted Sprint 3 measurement method and baseline semantics. |
| Documentation | Governs observation schema, local service, user workflow, accessibility, limitations, and final contract promotion. |
| CLI Surface | Governs the development command that starts bounded observable scenarios and reports service failures. |
| Frontend | Governs feature ownership, dependency direction, state placement, typed transport access, component tests, and rendering cost. |
| Accessibility | Governs semantics, keyboard operation, focus, non-colour status, reduced motion, and accessibility checks. |
| API Contracts | Governs the explicit versioned read-only DTO, compatibility, collection bounds, and stable errors. |
| Privacy And Data | Governs diagnostic and failure payload minimization, redaction, and safe browser presentation. |
| Security | Governs loopback binding, safe defaults, input bounds, browser transport, and prevention of accidental remote exposure. |

## Selected Evidence Reports

| Report | Path | Covers |
|---|---|---|
| Streaming Execution Semantics | `studies/agent-harness-study/reports/final/01.08-streaming-execution-semantics.md` | Producer and consumer ownership, terminal delivery, cancellation, and stream completion. Apply only to lifecycle observation. |
| Delivery Guarantees And Idempotency | `studies/agent-harness-study/reports/final/01.09-delivery-guarantees-and-idempotency.md` | Canonical recording versus repeated delivery, stable identity, and delivery claim limits. |
| Replay And Determinism | `studies/agent-harness-study/reports/final/01.10-replay-and-determinism.md` | Ordered reconstruction, replay boundaries, and deterministic interpretation. |
| Snapshot And Checkpoint Architecture | `studies/agent-harness-study/reports/final/02.02-snapshot-and-checkpoint-architecture.md` | Coherent read snapshots. Durable checkpoint and recovery findings remain excluded. |
| Failure Visibility | `studies/agent-harness-study/reports/final/13.03-failure-visibility.md` | Structured cause and diagnostic presentation across result, event, log, and caller views. |
| Export Interoperability And Observability | `studies/agent-harness-study/reports/final/10.04-export-interoperability-observability.md` | Boundary between canonical runtime evidence and external observation representations. |
| Go Ordered Observation Live Streaming And Backpressure | `studies/aren-go-runtime-study/reports/final/01.04-ordered-observation-live-streaming-and-backpressure.md` | Snapshot and live handoff, cursor semantics, abandonment, sequence order, terminal delivery, and backpressure. |
| Error Handling | `studies/go-cli-study/reports/final/05-error-handling.md` | Stable, actionable failures and safe diagnostic presentation. |
| IO Abstraction | `studies/go-cli-study/reports/final/06-io-abstraction.md` | Injectable process boundaries and deterministic command and service tests. |
| Logging And Observability | `studies/go-cli-study/reports/final/10-logging-observability.md` | Structured diagnostics, correlation, verbosity, and separation of canonical output from logs. |
| Testing Strategy | `studies/go-cli-study/reports/final/11-testing-strategy.md` | Unit, integration, fixture, command, and end-to-end test boundaries. |
| Go CLI Performance | `studies/go-cli-study/reports/final/14-performance.md` | Bounded local service, update, payload, and profiling considerations. |

The available report catalog does not contain comparative browser information-architecture research. Sprint 4 UX decisions therefore rest on the PRD, observability mandate, frontend and accessibility contracts, actual Phase 1 evidence, and explicit testing rather than invented external support.

## Selected Reasoning Templates

| Template | Output Path | Why Selected |
|---|---|---|
| Aren Observation Contract And Local Service | `projects/aren-phase-01-execution-lifecycle/sprints/04-observability-frontend/reasoning/observation-contract-and-local-service.md` | Define DTO ownership, versioning, snapshot and live delivery, loopback lifecycle, bounds, privacy, compatibility, and daemon exclusions. |
| Aren Run Explorer Experience | `projects/aren-phase-01-execution-lifecycle/sprints/04-observability-frontend/reasoning/run-explorer-experience.md` | Define information hierarchy, timeline semantics, state and failure language, evidence distinction, accessibility, responsiveness, and realistic density. |
| Aren Frontend Verification And Observer Isolation | `projects/aren-phase-01-execution-lifecycle/sprints/04-observability-frontend/reasoning/frontend-verification-and-observer-isolation.md` | Prove cross-interface agreement, connection schedules, slow and failed client isolation, schema handling, accessibility, rendering bounds, and real-browser behaviour. |

## Selected Project Reasoning

| Document | Path | Why Selected |
|---|---|---|
| Project synthesis | `projects/aren-phase-01-execution-lifecycle/project-reasoning/reasoning.md` | Preserves lifecycle truth and authority while the frontend projects the realised system. |
| Lifecycle authority | `projects/aren-phase-01-execution-lifecycle/project-reasoning/areas/01-lifecycle-authority-and-atomic-publication.md` | Prevents the service or browser from committing or reconstructing lifecycle truth independently. |
| Outcomes and terminal resolution | `projects/aren-phase-01-execution-lifecycle/project-reasoning/areas/02-outcomes-failures-and-terminal-resolution.md` | Governs exact outcome, panic, failure, and cancellation meanings shown in the UI. |
| Cancellation and cleanup | `projects/aren-phase-01-execution-lifecycle/project-reasoning/areas/03-cancellation-goroutine-ownership-and-cleanup.md` | Governs client abandonment, service shutdown, and proof that observation cannot strand producers. |
| Events and observation | `projects/aren-phase-01-execution-lifecycle/project-reasoning/areas/04-events-observation-waiting-and-replay.md` | Governs sequence order, replay, late observation, terminal completion, and delivery limits. |
| Verification and Go correctness | `projects/aren-phase-01-execution-lifecycle/project-reasoning/areas/05-verification-and-go-correctness.md` | Keeps runtime invariants authoritative across DTO, frontend, and browser tests. |

The amended PRD, roadmap, observability mandate, performance document, and active Aren contracts govern Sprint 4 where the earlier project synthesis is silent about the browser frontend.

## Prior Decisions To Carry Forward

| Decision or artifact | Path | Constraint For This Sprint |
|---|---|---|
| Sprint 1 decision synthesis | `projects/aren-phase-01-execution-lifecycle/sprints/01-core-lifecycle/reasoning.md` | Preserve identity, transition, outcome, timing, history, failure, and waiting semantics. |
| Sprint 1 review | `projects/aren-phase-01-execution-lifecycle/sprints/01-core-lifecycle/review.md` | Carry confirmed implementation facts into presentation and agreement tests. |
| Sprint 2 decision synthesis | `projects/aren-phase-01-execution-lifecycle/sprints/02-cancellation-and-concurrency/reasoning.md` | Treat realised cancellation, terminal resolution, observer, replay, and delivery semantics as authoritative. |
| Sprint 2 review | `projects/aren-phase-01-execution-lifecycle/sprints/02-cancellation-and-concurrency/review.md` | Carry correctness, race, leak, CLI, and limitation evidence into Sprint 4. |
| Sprint 3 requirements | `projects/aren-phase-01-execution-lifecycle/sprints/03-performance-methodology-and-baseline/requirements.md` | Present only measurements and interpretations that the accepted method supports. |
| Sprint 3 technical handbook | `projects/aren-phase-01-execution-lifecycle/sprints/03-performance-methodology-and-baseline/technical-handbook.md` | Reuse accepted measurement and evidence guidance without reopening the method casually. |
| Sprint 3 area reasoning | `projects/aren-phase-01-execution-lifecycle/sprints/03-performance-methodology-and-baseline/reasoning/` | Preserve workload, resource, schema, variance, and regression-policy decisions. |
| Sprint 3 decision synthesis | `projects/aren-phase-01-execution-lifecycle/sprints/03-performance-methodology-and-baseline/reasoning.md` | Treat this as the authoritative Phase 1 performance decision set. Silence does not supersede it. |
| Sprint 3 plan | `projects/aren-phase-01-execution-lifecycle/sprints/03-performance-methodology-and-baseline/plan.md` | Use intended boundaries with realised code and review evidence. |
| Sprint 3 review | `projects/aren-phase-01-execution-lifecycle/sprints/03-performance-methodology-and-baseline/review.md` | Carry accepted schemas, commands, baseline results, limitations, and unresolved frontend-relevant findings. |

Before Sprint 4 starts, add completed `execute.md`, `smoke.md`, promoted contract, and baseline artifact paths that materially affect the frontend. Record why any expected artifact is absent.

## Required Review Protocols

| Protocol | Path | Required Evidence |
|---|---|---|
| Architecture Review | `system/protocols/architecture-review-protocol.md` | Runtime truth ownership, DTO and frontend boundaries, local service lifecycle, dependency direction, bounded delivery, and daemon exclusion. |
| Sprint Review | `system/protocols/review-sprint-protocol.md` | Requirement traceability, schema compatibility, interface agreement, live and late observation, failure states, accessibility, responsiveness, isolation, performance, tests, and documentation. |
| Deep Smoke Sprint | `system/protocols/deep-smoke-sprint-protocol.md` | Real browser inspection of success, failure, cancellation, and collision scenarios with live updates, late connection, disconnection, malformed evidence, and cross-interface comparison. |

## Excluded Context

| Context | Reason Excluded | Revisit If |
|---|---|---|
| Implementation execution during planning | Planning must settle authority, data, experience, and proof before source mutation starts. | `plan.md` validates and governed execution begins. |
| Smoke investigation during planning | Browser smoke must inspect the realised service and frontend rather than a proposed interaction. | Implementation and review provide a runnable target. |
| Review automation during planning | Review must compare actual runtime, DTO, HTML, live updates, accessibility, and browser behaviour. | Execution produces reviewable code and evidence. |
| Issue tracking mutation | Findings remain in Sprint 4 review until separate work receives explicit ownership. | Final review identifies a deferred issue that belongs to another phase. |
| Git mutation | Planning and browser testing do not authorize commits, branches, merges, resets, or index changes. | A governed merge stage owns repository history. |
| Browser lifecycle control | Observation must remain passive, so the browser cannot cancel, transition, retry, or otherwise operate a run. | A later control-plane phase defines authority and safety. |
| Persistent daemon and remote API | The process-scoped loopback service exists only for local Phase 1 observation. | Daemon hosting defines persistence, remote identity, authorization, and reconnectable ownership. |
| Durable or global telemetry | Phase 1 reads bounded in-memory lifecycle and benchmark evidence. | Persistence or streaming volume creates a concrete retention requirement. |
| Future execution hierarchy | Phase 1 has one run and must not invent agent, model, tool, attempt, or workflow screens. | The owning execution phases provide realised semantics and evidence. |

## Next Artifacts

- `technical-handbook.md` distills the selected reports and applicable prior decisions.
- `reasoning/*.md` resolves the observation contract, run experience, and proof areas.
- `reasoning.md` makes final Sprint 4 decisions and classifies relevant earlier decisions.
- `plan.md` implements only those decisions.
- `review.md` and `smoke.md` verify the realised browser experience and close Phase 1.

