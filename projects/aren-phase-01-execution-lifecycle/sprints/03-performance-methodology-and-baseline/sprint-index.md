# Sprint index: Performance methodology and baseline

> Project: `aren-phase-01-execution-lifecycle`
> Sprint: `03-performance-methodology-and-baseline`
> Purpose: select the contracts, evidence, prior lifecycle decisions, reasoning templates, and review protocols needed to establish Aren's repeatable performance method and first Phase 1 baseline.
> Inputs Used: `projects/aren-phase-01-execution-lifecycle/project-index.md`, `projects/aren-phase-01-execution-lifecycle/roadmap.md`, `projects/aren-phase-01-execution-lifecycle/docs/PRD.md`, `projects/aren-phase-01-execution-lifecycle/docs/performance-engineering.md`, accepted project reasoning, and completed Sprint 1 and Sprint 2 artifacts.

This index is prepared in advance. Before Sprint 3 starts, confirm every named prior artifact exists and update selections affected by the realised lifecycle implementation.

## Sprint Scope

- **Sprint Goal:** Establish a reproducible performance experiment method over the completed Phase 1 lifecycle and record the first attributable baseline without weakening runtime semantics.
- **Planned Output:** Canonical workloads, benchmark and bounded scenario runners, environment and source metadata, statistical comparison, concurrency and resource evidence, profiling commands, execution tiers, machine-readable results, and the first Phase 1 baseline report.
- **Depends On:** Accepted project reasoning plus completed and reviewed Sprint 1 and Sprint 2 planning, implementation, tests, contracts, and diagnostic commands.
- **Non-Goals:** Invented service-level targets, benchmark-driven semantic changes, a performance database, hosted or distributed runners, an optimization campaign without measured cause, daemon hosting, frontend implementation, providers, tools, persistence, or later execution types.

## Source Project Index

- `projects/aren-phase-01-execution-lifecycle/project-index.md` is the authoritative catalog. Every selected contract, report, reasoning template, and protocol below appears there.

## Selected Contracts

| Contract | Why Selected |
|---|---|
| Architecture | Keeps measurement code outside lifecycle authority, preserves dependency direction, and prevents a duplicate execution path. |
| Errors | Governs invalid, partial, cancelled, timed-out, and failed measurements without false success. |
| Observability | Governs attributable evidence, correlation, diagnostic separation, and measurement overhead visibility. |
| Testing | Governs deterministic workloads, fixtures, repeatability, race checks, failure injection, and evidence traceability. |
| Documentation | Requires exact reproduction commands, environment assumptions, result interpretation, and baseline limits. |
| CLI Surface | Governs benchmark command discoverability, bounded flags, output, failures, and exit status. |
| Performance | Governs workload bounds, measurement before optimization, concurrency, profiling, resource evidence, and future regression policy. |

## Selected Evidence Reports

| Report | Path | Covers |
|---|---|---|
| Go Adversarial Concurrency And Failure Verification | `studies/aren-go-runtime-study/reports/final/01.03-adversarial-concurrency-and-failure-verification.md` | Controlled schedules, stress, negative controls, race execution, and leak verification over the realised lifecycle. |
| Go Ordered Observation Live Streaming And Backpressure | `studies/aren-go-runtime-study/reports/final/01.04-ordered-observation-live-streaming-and-backpressure.md` | Observer costs, event replay, delivery bounds, and backpressure risks that require measurement. |
| Dataset And Golden Task Management | `studies/agent-harness-study/reports/final/18.01-dataset-golden-task-management.md` | Stable scenario identity, reproducible inputs, versioning, and limits on importing LLM evaluation machinery. |
| Regression Gating And CI Integration | `studies/agent-harness-study/reports/final/18.03-regression-gating-ci-integration.md` | Evidence tiers, repeatability, regression reporting, and the distinction between signal and gates. |
| Cost Latency And Quality Evaluation | `studies/agent-harness-study/reports/final/18.04-cost-latency-quality-evaluation.md` | Attribution, latency interpretation, workload shape, sampling, and evaluation caveats. |
| Concurrency | `studies/go-cli-study/reports/final/08-concurrency.md` | Bounded workers, cancellation, task ownership, and concurrent test safety. |
| Testing Strategy | `studies/go-cli-study/reports/final/11-testing-strategy.md` | Unit, integration, fixture, command, race, and environment-sensitive test boundaries. |
| Go CLI Performance | `studies/go-cli-study/reports/final/14-performance.md` | Go benchmark mechanics, profiling, allocation, startup, concurrency, and performance tradeoffs. |
| Philosophy | `studies/go-cli-study/reports/final/15-philosophy.md` | Scope control and resistance to speculative optimization machinery. |

## Selected Reasoning Templates

| Template | Output Path | Why Selected |
|---|---|---|
| Aren Benchmark Method And Workloads | `projects/aren-phase-01-execution-lifecycle/sprints/03-performance-methodology-and-baseline/reasoning/benchmark-method-and-workloads.md` | Define workload identity, measured intervals, tiers, sampling, comparison, metadata, invalid runs, and later extension. |
| Aren Scale Resource And Profile Evidence | `projects/aren-phase-01-execution-lifecycle/sprints/03-performance-methodology-and-baseline/reasoning/scale-resource-and-profile-evidence.md` | Decide how to measure scaling, latency, allocation, memory, runtime tasks, contention, saturation, quiescence, and profiles. |
| Aren Performance Evidence And Regression Policy | `projects/aren-phase-01-execution-lifecycle/sprints/03-performance-methodology-and-baseline/reasoning/performance-evidence-and-regression-policy.md` | Define result schemas, attribution, variance, failure truth, retention, baseline interpretation, and entry conditions for future gates. |

## Selected Project Reasoning

| Document | Path | Why Selected |
|---|---|---|
| Project synthesis | `projects/aren-phase-01-execution-lifecycle/project-reasoning/reasoning.md` | Preserves lifecycle authority, semantic truth, cancellation, observer isolation, and correctness requirements while Sprint 3 measures their cost. |
| Lifecycle authority | `projects/aren-phase-01-execution-lifecycle/project-reasoning/areas/01-lifecycle-authority-and-atomic-publication.md` | Prevents measurement code from creating a second lifecycle path or moving instrumentation into transition authority without cause. |
| Cancellation and cleanup | `projects/aren-phase-01-execution-lifecycle/project-reasoning/areas/03-cancellation-goroutine-ownership-and-cleanup.md` | Supplies ownership and quiescence facts needed for runtime-task and retained-resource measurement. |
| Events and observation | `projects/aren-phase-01-execution-lifecycle/project-reasoning/areas/04-events-observation-waiting-and-replay.md` | Defines the observer and replay semantics whose costs Sprint 3 measures. |
| Verification and Go correctness | `projects/aren-phase-01-execution-lifecycle/project-reasoning/areas/05-verification-and-go-correctness.md` | Keeps correctness evidence authoritative when benchmark pressure suggests changes. |

The amended PRD, roadmap, and performance mandate govern Sprint 3 where the earlier project synthesis is silent about the dedicated performance sprint.

## Prior Decisions To Carry Forward

| Decision or artifact | Path | Constraint For This Sprint |
|---|---|---|
| Sprint 1 decision synthesis | `projects/aren-phase-01-execution-lifecycle/sprints/01-core-lifecycle/reasoning.md` | Preserve realised lifecycle, outcome, failure, timing, history, and waiter semantics. |
| Sprint 1 review | `projects/aren-phase-01-execution-lifecycle/sprints/01-core-lifecycle/review.md` | Carry confirmed implementation facts and unresolved measurement-relevant findings. |
| Sprint 2 requirements | `projects/aren-phase-01-execution-lifecycle/sprints/02-cancellation-and-concurrency/requirements.md` | Measure the complete cancellation, observer, replay, and concurrency behaviour rather than an earlier subset. |
| Sprint 2 technical handbook | `projects/aren-phase-01-execution-lifecycle/sprints/02-cancellation-and-concurrency/technical-handbook.md` | Reuse applicable Go and concurrency evidence without repeating lifecycle analysis. |
| Sprint 2 area reasoning | `projects/aren-phase-01-execution-lifecycle/sprints/02-cancellation-and-concurrency/reasoning/` | Preserve accepted ownership, delivery, terminal resolution, and verification decisions. |
| Sprint 2 decision synthesis | `projects/aren-phase-01-execution-lifecycle/sprints/02-cancellation-and-concurrency/reasoning.md` | Treat this as the authoritative realised lifecycle decision set. Silence does not supersede it. |
| Sprint 2 plan | `projects/aren-phase-01-execution-lifecycle/sprints/02-cancellation-and-concurrency/plan.md` | Use intended boundaries together with realised code and review evidence. |
| Sprint 2 review | `projects/aren-phase-01-execution-lifecycle/sprints/02-cancellation-and-concurrency/review.md` | Carry correctness results, implementation facts, and unresolved performance questions into Sprint 3. |

Before Sprint 3 starts, add any completed `execute.md`, `smoke.md`, or promoted contract paths that materially affect measurement and record why an expected artifact is absent.

## Required Review Protocols

| Protocol | Path | Required Evidence |
|---|---|---|
| Architecture Review | `system/protocols/architecture-review-protocol.md` | Measurement ownership, dependency direction, real-runtime execution, instrumentation boundaries, bounded runners, and absence of benchmark-only semantic paths. |
| Sprint Review | `system/protocols/review-sprint-protocol.md` | Requirement traceability, workload validity, reproducibility, metadata, statistical comparison, resource measurement, failure truth, tests, and baseline limitations. |
| Deep Smoke Sprint | `system/protocols/deep-smoke-sprint-protocol.md` | Real commands that reproduce a quick suite, an extended bounded scenario, a comparison, and at least one profile capture. |

## Excluded Context

| Context | Reason Excluded | Revisit If |
|---|---|---|
| Implementation execution during planning | Planning must settle workload and evidence rules before source mutation begins. | `plan.md` validates and governed execution begins. |
| Smoke investigation during planning | Smoke must reproduce the realised benchmark commands and evidence, not a predicted runner. | Implementation and review produce a valid target. |
| Review automation during planning | Review must inspect realised measurements, failure paths, and retained artifacts. | Execution produces reviewable code and baseline evidence. |
| Issue tracking mutation | Findings belong in Sprint 3 review before separate work receives ownership. | Review identifies a deferred performance issue that warrants another project. |
| Git mutation | Planning and benchmark execution do not authorize commits, branches, merges, resets, or index changes. | A governed merge stage owns repository history. |
| Frontend implementation | Sprint 3 defines performance evidence that Sprint 4 can consume but does not build the browser UI. | Sprint 4 begins with accepted baseline and schemas. |
| Arbitrary numeric targets | Host-specific baseline data cannot establish a product promise without representative workload evidence and measured variance. | Later operational evidence supports a governed target. |
| Performance database or hosted runner | Local and continuous-integration artifacts are enough for the first baseline. | Result volume or comparison needs prove a separate service is necessary. |

## Next Artifacts

- `technical-handbook.md` distills the selected evidence and applicable prior decisions.
- `reasoning/*.md` resolves the three selected measurement areas.
- `reasoning.md` makes the final Sprint 3 decisions and classifies relevant earlier decisions.
- `plan.md` implements only those decisions.
- `review.md` and `smoke.md` verify reproducibility against the realised benchmark system.

