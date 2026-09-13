# Aren Performance Engineering Contract

> Project: `aren-phase-01-execution-lifecycle`  
> Status: active project-wide engineering contract; intended to grow with Aren phase by phase

## Purpose

This contract defines the performance properties Aren must preserve and the evidence required before performance claims or optimizations are accepted.

Aren's primary requirement is correct, explicit execution semantics. Performance engineering exists to ensure those semantics add low, bounded, explainable overhead and continue to scale predictably as concurrency, event volume, execution duration, and runtime capability grow.

This contract is not a collection of universal numeric targets. It governs measurement, attribution, scaling, resource ownership, saturation, regression discipline, and the boundary between justified optimization and speculative performance architecture.

`06-verification.md` governs how claims are proved. This contract governs what performance properties and engineering discipline the runtime must preserve.

## Scope

Apply this contract whenever a change materially affects or introduces:

- runtime hot paths;
- execution creation, transition, terminal resolution, waiting, cancellation, or observation costs;
- concurrency or synchronization;
- retained memory or allocations;
- Aren-owned goroutines/runtime tasks;
- high-volume events, streams, buffers, queues, or telemetry;
- resource limits or admission;
- benchmark infrastructure or result schemas;
- profiling or optimization;
- local services, browser observation, transport, persistence, tools, model calls, workflows, or later execution types whose Aren-owned overhead becomes measurable.

## Requirement Index

| ID | Title | Severity If Violated |
| --- | --- | --- |
| AREN-PERF-001 | Correctness outranks performance | Blocker |
| AREN-PERF-002 | Measure before optimization | High |
| AREN-PERF-003 | Profile before complex optimization | High |
| AREN-PERF-004 | Aren overhead is separated from supervised work | Blocker |
| AREN-PERF-005 | Performance claims identify workload and environment | High |
| AREN-PERF-006 | Independent runs must not accidentally serialize globally | Blocker |
| AREN-PERF-007 | Runtime-task growth is intentional and owned | High |
| AREN-PERF-008 | Retained-memory growth must be explainable | High |
| AREN-PERF-009 | High-volume work and retention are bounded or justified | Blocker |
| AREN-PERF-010 | Tail behaviour matters under contention | High |
| AREN-PERF-011 | Saturation and recovery are measured | High |
| AREN-PERF-012 | Overload must not be hidden by unbounded buffering | Blocker |
| AREN-PERF-013 | Baselines accumulate phase by phase | High |
| AREN-PERF-014 | Numeric gates require stable evidence before enforcement | High |
| AREN-PERF-015 | Material new runtime costs extend the benchmark suite | High |
| AREN-PERF-016 | Optimizations preserve all governing semantic contracts | Blocker |

## Requirements

### AREN-PERF-001 — Correctness Outranks Performance

A performance improvement is invalid if it weakens Aren's lifecycle, cancellation, failure, ordering, observation, ownership, privacy, or verification guarantees.

Do not accept a faster implementation that introduces:

- torn lifecycle publication;
- misleading cancellation or terminality;
- observer-dependent execution progress;
- hidden failure;
- unowned background work;
- weakened retention or evidence without an explicit semantic decision;
- data races, leak risks, or unverifiable synchronization;
- benchmark-only execution paths that differ semantically from the real runtime.

Where a performance mechanism changes semantics, treat it as a product/runtime design change first and a performance change second.

### AREN-PERF-002 — Measure Before Optimization

Do not introduce pools, schedulers, custom allocators, lock-free structures, batching, event coalescing, caches, worker layers, or other performance architecture merely because they may be useful later.

Before a nontrivial optimization:

- reproduce the observed regression or scaling boundary;
- identify the workload that exposes it;
- measure the current behaviour;
- state which metric is materially problematic;
- preserve a baseline suitable for comparison.

Simple clarity is preferred until evidence demonstrates that it is insufficient.

### AREN-PERF-003 — Profile Before Complex Optimization

When the cause of a material performance problem is not already obvious from a narrow benchmark, use appropriate profiling before making a complex optimization.

Relevant Go evidence may include:

- CPU profiles;
- heap/allocation profiles;
- goroutine profiles;
- mutex profiles;
- block profiles;
- runtime traces.

The optimization rationale must identify the measured cause it addresses rather than infer causality from a single aggregate number.

### AREN-PERF-004 — Aren Overhead Is Separated From Supervised Work

Performance evidence must distinguish Aren's own cost from the cost of work Aren supervises wherever meaningful.

Examples include separating:

- lifecycle setup from user work duration;
- event recording from provider/network latency;
- cancellation acceptance from work cooperation latency;
- Aren allocation/retention from subprocess or model memory;
- browser/local-service overhead from runtime execution;
- serialization or persistence overhead from external waiting.

Do not present provider response time, user work duration, or child-process resource consumption as if it were solely Aren runtime overhead.

### AREN-PERF-005 — Performance Claims Identify Workload And Environment

A performance result is meaningful only with enough context to reproduce and interpret it.

Record, as applicable:

- scenario/workload identity;
- source revision and diff state;
- compiler/runtime/toolchain version;
- OS and architecture;
- relevant CPU/hardware details;
- concurrency;
- workload size and duration;
- repetitions and spread/statistical method;
- binary or artifact identity where material;
- instrumentation enabled during the measurement;
- important environment limits or noise caveats.

Do not compare unlike workload shapes as if the difference measured only code quality.

### AREN-PERF-006 — Independent Runs Must Not Accidentally Serialize Globally

Independent executions should not contend on runtime-wide mutable state unless a global semantic requirement explicitly requires coordination.

When shared structures are introduced, measure their contention and justify their ownership.

Examples that require explicit scrutiny include later:

- schedulers;
- registries;
- policy managers;
- resource/accounting managers;
- persistence coordinators;
- global observation structures;
- shared caches.

Per-run synchronization is preferable when the semantics are per-run. Do not create global serialization merely for implementation convenience.

### AREN-PERF-007 — Runtime-Task Growth Is Intentional And Owned

Every long-lived Aren-owned goroutine/runtime task must have a semantic owner as required by `04-cancellation-and-lifetimes.md` and a measured scaling cost when its population can grow materially.

Measure where relevant:

- tasks created per active run;
- tasks retained per blocked run;
- tasks added by observers, cancellation support, timers, streams, transports, or telemetry;
- tasks remaining after completion;
- whether task count returns toward baseline after load is removed.

One task per run may be acceptable. Accidental task multiplication is not.

### AREN-PERF-008 — Retained-Memory Growth Must Be Explainable

Measure transient allocation and retained memory separately where meaningful.

Memory growth should be approximately explainable from active Aren-owned state, retained evidence, buffers, and explicitly owned support structures.

Investigate when equivalent workload growth produces materially superlinear retained-memory growth without a clear semantic reason.

For high-concurrency Aren workloads, retained bytes per active run may matter more than microsecond-level lifecycle latency.

Do not infer memory efficiency from allocation microbenchmarks alone.

### AREN-PERF-009 — High-Volume Work And Retention Are Bounded Or Justified

Any potentially high-volume Aren-owned queue, buffer, history, stream, output, cache, retry record, telemetry series, or client delivery path must declare a practical bound or a measured reason not to have one.

Phase 1's tiny complete lifecycle history is not precedent for retaining all later token deltas, tool output, progress, logs, or telemetry indefinitely.

Where a bound can discard data:

- identify what may be discarded;
- preserve canonical semantic truth;
- expose gap, truncation, sampling, or coalescing semantics honestly where callers need to know;
- ensure the bound cannot silently turn overload into unbounded latency elsewhere.

### AREN-PERF-010 — Tail Behaviour Matters Under Contention

Average latency is insufficient for paths affected by contention or fan-out.

Where useful, record distributions such as:

- p50;
- p95;
- p99;
- bounded-run maximum.

Tail behaviour is especially relevant for:

- lifecycle transitions;
- cancellation acceptance and propagation;
- event append/read paths;
- observer delivery;
- local service/API paths;
- later stream, tool, scheduler, persistence, and daemon control paths.

Do not optimize a mean while allowing unacceptable tail amplification to remain invisible.

### AREN-PERF-011 — Saturation And Recovery Are Measured

For a capability whose scale can grow materially, performance work must eventually determine more than its comfortable-case throughput.

Where relevant, identify:

- the scaling curve;
- the knee where degradation becomes nonlinear;
- the first resource to saturate;
- behaviour at and beyond saturation;
- whether invariants remain correct under load;
- whether memory, runtime tasks, queues, and latency recover after load falls.

Graceful degradation and recoverability are more important than a headline maximum throughput number.

### AREN-PERF-012 — Overload Must Not Be Hidden By Unbounded Buffering

Queues, retries, worker pools, caches, and buffers may move where work waits; they do not remove overload.

A performance mechanism must not silently convert saturation into:

- unbounded memory growth;
- unbounded queueing latency;
- runaway runtime-task creation;
- hidden repeated external work;
- starvation;
- indefinite cleanup delay.

When a queue or buffer is introduced, define capacity, ownership, saturation behaviour, and recovery semantics.

Where rejection or backpressure is eventually required, expose it truthfully rather than pretending work was accepted without bound.

### AREN-PERF-013 — Baselines Accumulate Phase By Phase

Performance evidence grows with Aren's capabilities.

Each phase should:

1. retain the relevant prior benchmark scenarios;
2. add focused scenarios for materially new Aren-owned costs;
3. compare incremental cost against an appropriate prior baseline;
4. record changes in workload meaning or measurement method;
5. avoid benchmarking future capabilities before they exist.

Phase 1 Sprint 3 establishes the first formal measurement method and attributable lifecycle baseline. It does not mark the end of performance work.

### AREN-PERF-014 — Numeric Gates Require Stable Evidence Before Enforcement

Initial baselines are measurements, not service-level promises.

Do not create hard percentage or capacity gates until:

- the benchmark is stable enough to distinguish signal from noise;
- the workload represents a meaningful supported use case;
- variance and environment sensitivity are understood;
- the comparison method is repeatable;
- the threshold has a product or operational rationale.

Before that point, performance regressions should produce evidence and investigation rather than arbitrary pass/fail policy.

Correctness regressions, leaks, unbounded growth, global serialization without semantic cause, or observer-controlled execution remain blockers even without numeric budgets.

### AREN-PERF-015 — Material New Runtime Costs Extend The Benchmark Suite

When a phase introduces a materially new Aren-owned cost, add the smallest useful benchmark/scenario needed to observe it.

Examples include later:

- model request bookkeeping;
- stream delta handling;
- structured validation;
- retries/backoff;
- tool supervision;
- agent-loop iteration;
- context construction/compaction;
- persistence and recovery;
- execution composition;
- admission and policy checks;
- daemon transport;
- workflow coordination.

Do not add benchmark matrices merely for completeness when the new capability does not create a material performance question.

### AREN-PERF-016 — Optimizations Preserve All Governing Semantic Contracts

Any optimization must continue to satisfy the active Aren contracts.

In particular:

- lifecycle facts remain coherent;
- cancellation remains truthful;
- observation remains passive;
- failures remain inspectable;
- evidence remains attributable;
- bounds remain explicit;
- owned resources still release;
- verification remains capable of falsifying the optimized path.

When performance pressure suggests changing a semantic guarantee, explicitly supersede the guarantee through governed reasoning rather than weakening it implicitly inside optimized code.

## Canonical Performance Dimensions

Aren's benchmark and review vocabulary should distinguish, where applicable:

- incremental runtime latency/overhead;
- throughput;
- concurrency scaling;
- bytes and allocations per operation;
- peak and retained heap;
- memory per active run;
- Aren-owned runtime-task growth;
- mutex/block contention;
- p50/p95/p99 tail latency;
- saturation point and saturated behaviour;
- recovery after load falls.

No single requests-per-second number is the primary definition of Aren performance.

## Canonical Synthetic Workload Classes

Use controlled work to isolate Aren costs where useful:

- **zero-work** — immediate return to expose lifecycle overhead;
- **short blocked work** — controlled short waits to expose active-run bookkeeping and scheduling;
- **long blocked work** — many blocked runs to measure idle active-run footprint and runtime-task scaling;
- **event-heavy work** — controlled evidence volume once such behaviour exists;
- **cancellation storm** — concurrent cancellation pressure once cancellation exists.

Representative real workloads should supplement these synthetic classes when the corresponding Aren capabilities exist.

## Relationship To Verification And Observability

`06-verification.md` governs whether the benchmark/scenario actually exercises the claimed path, whether evidence is reproducible, and whether the claim's limitations are recorded.

`07-observability.md` governs how performance/resource evidence is correlated, distinguished from canonical semantic truth, exposed through observation surfaces, and bounded so measurement itself does not control execution.

A performance dashboard, benchmark report, or browser view is a projection of measured evidence. It must not manufacture unsupported precision or causal attribution.

## Primary Governing Sources

- `projects/aren-phase-01-execution-lifecycle/docs/performance-engineering.md`
- `projects/aren-phase-01-execution-lifecycle/docs/observability-mandate.md`
- `projects/aren-phase-01-execution-lifecycle/docs/final-language-decision.md`
- `projects/aren-phase-01-execution-lifecycle/contracts/04-cancellation-and-lifetimes.md`
- `projects/aren-phase-01-execution-lifecycle/contracts/06-verification.md`
