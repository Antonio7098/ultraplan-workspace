# Aren Performance Engineering

## 1. Purpose

Performance is a first-class engineering property of Aren because the runtime is intended to supervise many executions concurrently while adding as little avoidable overhead as possible.

This document defines how Aren will assess, benchmark, improve, and preserve performance throughout development.

It is not a promise that every operation will be minimized for raw speed. Aren's primary requirement remains correct and explicit execution semantics. Performance engineering exists to ensure that those semantics continue to behave predictably as concurrency, event volume, execution duration, and runtime capability grow.

The core performance objective is:

> Aren should add low, bounded, explainable overhead to supervised work, and that overhead should scale predictably as the number and complexity of concurrent executions increase.

The most important properties are therefore not a single requests-per-second number. They are:

- incremental runtime overhead;
- concurrency scaling;
- memory efficiency;
- allocation behaviour;
- Aren-owned runtime-task growth;
- synchronization contention;
- tail latency;
- saturation behaviour;
- graceful degradation under overload;
- performance regression visibility over time.

---

## 2. Relationship to Aren's Development Philosophy

Performance follows the same rules as the rest of Aren.

### Measure before abstraction

Do not introduce pools, schedulers, custom allocators, lock-free structures, batching layers, event coalescing, or other performance architecture because they might be useful later.

First demonstrate a real bottleneck or scaling boundary.

### Measure before optimization

A benchmark regression or observed scaling knee should be reproduced and profiled before implementation changes are made.

### Correctness before speed

An optimization must preserve lifecycle, cancellation, event, failure, ordering, and ownership guarantees. A faster implementation that weakens semantics is a regression.

### Real workloads before arbitrary targets

Initial benchmarks establish baselines. Numeric service-level targets should be introduced only after real Aren workloads provide evidence for meaningful limits.

### Performance evidence accumulates phase by phase

Performance remains part of every Aren phase. Each phase inherits the benchmark suite of previous phases and adds focused scenarios for materially new runtime costs.

Phase 1 contains one dedicated sprint to establish the measurement method, runner, evidence format, comparison procedure, and first baseline after the lifecycle is realised. This is not permission to postpone performance work in other sprints or phases. Earlier sprints add the benchmarks their behaviour makes possible, and later phases extend the same method alongside their runtime work.

---

## 3. Performance Dimensions

Aren should evaluate performance across the following dimensions.

### 3.1 Runtime overhead

Measure the cost Aren itself adds around controlled work.

Examples include:

- run creation;
- lifecycle transition;
- terminal resolution;
- state inspection;
- waiter release;
- cancellation acceptance;
- event recording;
- event replay;
- observer registration;
- later, progress recording, stream handling, tool supervision, persistence, policy checks, and daemon transport.

Aren overhead must be distinguished from the latency of model providers, tools, network calls, subprocesses, and user work.

### 3.2 Throughput

Measure how much lifecycle or execution-management work Aren can complete over time.

Throughput is useful only when interpreted with workload shape and concurrency. A single maximum number must not become the primary performance claim.

### 3.3 Concurrency scaling

Measure how performance changes as independent active runs increase.

Canonical early levels should include:

```text
1
10
100
1,000
5,000
10,000
```

Larger levels may be added when justified by the environment and workload.

The purpose is to identify the scaling curve and the point at which degradation becomes nonlinear.

### 3.4 Memory efficiency

Measure both transient allocation and retained memory.

Important metrics include:

- bytes allocated per operation;
- allocations per operation;
- heap growth per active run;
- retained event-history cost;
- memory after run completion;
- memory-management pressure at increasing concurrency.

For high concurrency, memory per active run may matter more than microsecond-level lifecycle latency.

### 3.5 Runtime-task efficiency

Track Aren-owned runtime tasks explicitly.

Aren should know:

- runtime tasks created per active run;
- runtime tasks retained per blocked run;
- runtime tasks remaining after completion;
- whether observers, cancellation, waiting, streaming, or telemetry multiply runtime-task count;
- whether runtime-task count returns to baseline after load is removed.

A one-task-per-run architecture may be entirely acceptable. Unintentional
runtime-task multiplication is not.

### 3.6 Synchronization contention

Measure blocking and mutex contention as concurrency grows.

Independent runs should not contend on global mutable state unless the capability intrinsically requires global coordination.

Per-run synchronization is generally preferable to accidental runtime-wide serialization.

Any shared structure introduced later, such as a scheduler, policy manager, persistence coordinator, registry, or global event stream, must be evaluated for contention explicitly.

### 3.7 Tail latency

Average latency is insufficient under contention.

Where applicable, record:

- p50;
- p95;
- p99;
- maximum observed latency during bounded benchmark runs.

Tail behaviour is especially important for lifecycle transitions, cancellation, stream delivery, tool supervision, and daemon control paths.

### 3.8 Saturation and overload behaviour

Aren must eventually be tested beyond comfortable capacity.

The objective is to identify:

- the scaling knee;
- maximum sustainable concurrency for a scenario;
- the resource that saturates first;
- behaviour at and beyond saturation;
- recovery once load is removed.

Graceful degradation is more important than maximizing the highest possible benchmark number.

---

## 4. Benchmark Layers

Aren's performance suite should contain several layers. They answer different questions and should not be collapsed into one benchmark category.

## 4.1 Layer 1 — Microbenchmarks

Microbenchmarks measure isolated Aren primitives with minimal work attached.

Early Phase 1 candidates include:

```text
BenchmarkRunSuccess
BenchmarkRunFailure
BenchmarkRunCancellation
BenchmarkStateRead
BenchmarkWaitCompleted
BenchmarkEventReplay
BenchmarkObserverRegistration
BenchmarkLifecycleTransition
BenchmarkTerminalResolution
```

Record runtime-appropriate benchmark metrics such as:

```text
ns/op
B/op
allocs/op
```

Microbenchmarks are useful for identifying local regressions and allocation changes, but they do not establish high-concurrency behaviour.

## 4.2 Layer 2 — Concurrent runtime benchmarks

Run canonical scenarios at increasing levels of concurrency.

Example dimensions:

```text
success/concurrency=1
success/concurrency=10
success/concurrency=100
success/concurrency=1000
success/concurrency=5000
success/concurrency=10000
```

Equivalent matrices should exist for blocked work, event-heavy work, and cancellation once those behaviours exist.

Capture, where practical:

- completed runs per second;
- latency distribution;
- bytes and allocations per run;
- peak and retained heap;
- peak runtime tasks;
- CPU use;
- block and mutex contention.

The main output of these tests is the shape of the scaling curve and the location of its knee.

## 4.3 Layer 3 — Stress and saturation tests

Stress tests deliberately increase load until performance degrades materially or a resource boundary is reached.

Example progression:

```text
1,000 active runs
2,000
5,000
10,000
20,000
50,000
...
```

The test should record not only where the process stops scaling well, but how it fails.

Desirable overload behaviour includes:

- increasing latency without invariant corruption;
- bounded queueing where queues exist;
- explicit rejection where limits exist;
- existing work continuing correctly;
- memory and runtime-task counts returning toward baseline after load falls.

Undesirable behaviour includes:

- uncontrolled memory growth;
- runaway runtime-task creation;
- deadlock;
- starvation caused by global contention;
- memory-management collapse;
- process termination from avoidable resource exhaustion;
- correctness failures caused by load.

## 4.4 Layer 4 — Representative workload benchmarks

As Aren gains capabilities, synthetic benchmarks must be supplemented by realistic mixtures.

Potential later profiles include:

### Agent-light

Many long-lived runs mostly waiting on external model or tool I/O with relatively little local processing.

### Stream-heavy

Many active model streams producing frequent deltas or progress events.

### Tool-heavy

A high proportion of runs supervising subprocess, shell, or external tool execution.

### Cancellation-heavy

Large numbers of active runs receiving concurrent cancellation requests.

### Mixed autonomous workload

A realistic mixture of waiting model calls, active streams, tool executions, agent-loop transitions, observers, and cancellations.

These profiles should be introduced only when the corresponding runtime capabilities exist.

## 4.5 Layer 5 — Long-duration soak tests

Once Aren owns persistent, daemon-hosted, or long-running execution, add long-duration tests that look for slow degradation rather than immediate throughput limits.

Measure:

- heap growth over time;
- runtime-task drift;
- event or state accumulation;
- descriptor or subprocess leaks;
- scheduler degradation;
- latency drift;
- cleanup after repeated execution cycles.

---

## 5. Canonical Synthetic Workloads

To isolate Aren overhead from work behaviour, use a small set of controlled work classes.

### 5.1 Zero-work

Work returns immediately.

Purpose:

- expose lifecycle overhead;
- maximize sensitivity to allocations and synchronization;
- compare implementation changes.

### 5.2 Short blocked work

Work blocks for controlled durations such as:

```text
1 ms
10 ms
100 ms
```

Purpose:

- expose scheduler effects;
- represent short I/O-style waits;
- test concurrent active-run bookkeeping.

### 5.3 Long blocked work

Large numbers of runs block on a channel or controlled wait until released.

Purpose:

- measure idle active-run footprint;
- measure Aren-owned runtime tasks per active run;
- model agent workloads waiting on network or tools;
- test memory scaling independently from work CPU.

### 5.4 Event-heavy work

Once progress or streaming exists, generate controlled event counts such as:

```text
1 event/run
10 events/run
100 events/run
1,000 events/run
```

Purpose:

- expose retention and delivery cost;
- detect allocation pressure;
- evaluate observer impact;
- identify when coalescing or bounded retention may become justified.

### 5.5 Cancellation storm

Create many active runs and request cancellation concurrently across a large subset.

Purpose:

- test cancellation-path contention;
- measure cancellation tail latency;
- combine correctness and performance stress;
- expose shared-state bottlenecks.

---

## 6. The Concurrency Matrix

A canonical concurrency matrix should evolve alongside Aren.

An early conceptual matrix is:

| Concurrent runs | Immediate success | Blocked work | Cancellation | Event-heavy |
|---:|---|---|---|---|
| 1 | benchmark | benchmark | benchmark | later |
| 10 | benchmark | benchmark | benchmark | later |
| 100 | benchmark | benchmark | benchmark | later |
| 1,000 | benchmark | benchmark | benchmark | later |
| 5,000 | benchmark | benchmark | benchmark | later |
| 10,000 | benchmark | benchmark | benchmark | later |

Each cell is a scenario, not merely a pass/fail result.

A scenario record should eventually contain enough information to reproduce and interpret the measurement:

```text
scenario
commit
compiler and runtime version
OS
architecture
CPU
concurrency
benchmark duration
runs completed
throughput
p50
p95
p99
ns/op where meaningful
B/op
allocs/op
peak heap
retained heap
peak runtime tasks
CPU utilization
notes
```

Not every metric must be collected by every lightweight benchmark. Heavier collection belongs in scheduled or release benchmark runs.

---

## 7. Performance Invariants and Architectural Guardrails

Before stable numeric targets exist, Aren should maintain qualitative performance invariants.

### Independent runs should remain independent

Independent executions should not acquire a runtime-wide lock or serialize through shared mutable state unless global coordination is an explicit semantic requirement.

### Observation should not enter the execution critical path

Slow, absent, abandoned, or faulty observers must not determine whether execution can progress.

This is already a Phase 1 correctness rule and should remain a performance rule as event volume grows.

### Runtime-task growth must be intentional

Every long-lived Aren-owned runtime task should have a clear owner, lifecycle,
termination condition, and measured scaling cost.

### High-volume retention must be bounded or justified

Phase 1 lifecycle histories are deliberately small and retained in memory. That rule must not automatically generalize to high-volume token deltas, progress updates, tool output, or telemetry.

Aren should distinguish between:

- canonical lifecycle facts that may deserve full retention;
- high-volume execution telemetry that may require bounded buffering, coalescing, sampling, streaming, or external persistence.

### Memory growth should be approximately explainable from active state

If doubling a representative number of equivalent active runs more than doubles Aren's retained memory without a clear reason, investigate.

### Performance mechanisms must not hide overload

Queues, retries, buffers, and workers may improve throughput, but they must not silently turn overload into unbounded latency or memory accumulation.

---

## 8. Phase-by-Phase Performance Evidence

Performance evidence should grow with the roadmap rather than anticipate future architecture.

### Phase 0 — Project foundation

Establish repeatable benchmark commands and record the environment required for meaningful comparison.

No performance target is required.

### Phase 1 — Execution lifecycle

Establish the first baseline for Aren's primitive lifecycle cost.

Recommended initial benchmarks:

- immediate successful run;
- returned failure;
- cancellation path;
- lifecycle transition;
- terminal resolution;
- event replay;
- completed wait;
- concurrent runs at increasing levels;
- many blocked active runs for memory and runtime-task footprint.

Phase 1 performance establishes an attributable baseline rather than a universal numeric target. Correctness failures, leaks, unbounded resource growth, global serialization of independent runs, or observation cost that controls execution block acceptance. Other results are recorded and interpreted rather than converted into arbitrary gates.

### Phase 2 — Controlled executor

Add benchmarks for the richer controlled execution behaviour introduced by the phase, such as progress and cleanup if they are implemented.

Compare incremental cost against the Phase 1 baseline.

From this phase onward, every phase review should include explicit performance evidence.

### Phase 3 — First real model invocation

Separate provider latency from Aren overhead around request setup, execution bookkeeping, cancellation, and result handling.

Do not use provider response time as the runtime performance metric.

### Phase 4 — Streaming model execution

Measure:

- per-delta processing cost;
- allocation rate;
- event or buffer growth;
- observer impact;
- many simultaneous streams;
- cancellation latency during heavy streaming.

This phase should explicitly test whether fine-grained streaming creates unacceptable allocation or retention pressure.

### Phase 5 — Structured model results

Measure validation and decoding cost separately from provider latency.

### Phase 6 — Retry and bounded repair

Measure per-attempt bookkeeping and ensure retry machinery does not create
unbounded retained state, runtime tasks, or hidden queues.

### Phase 7 — Tool-call representation

Representation itself should remain cheap. Benchmark only if serialization, copying, or validation becomes material.

### Phase 8 — Local tool execution

Measure concurrent tool supervision, subprocess/resource overhead where applicable, cancellation, cleanup, and descriptor/process limits.

### Phase 9 — Bounded agent loop

Measure Aren overhead per loop iteration and scaling across many simultaneous loops.

### Phase 10 — Context management

Measure context assembly, copying, pruning, serialization, and memory growth as conversation or execution history expands.

### Phase 11 — Persistence and recovery

Measure durable-write latency, batching effects, recovery throughput, and the incremental cost persistence adds to lifecycle operations.

### Phase 12 — Execution composition

Measure fan-out/fan-in behaviour, dependency coordination, and any shared scheduling structures.

### Phase 13 — Policies and resources

Measure concurrency-limit enforcement, budget checks, admission decisions, fairness, and contention in shared governance state.

### Phase 14 — Daemon hosting

Add transport and process-level metrics:

- request/control-path latency;
- connection count;
- memory per hosted run;
- daemon throughput;
- client fanout;
- serialization cost;
- recovery after client disconnect.

### Phase 15 — Multi-language clients

Measure client transport and serialization overhead without duplicating runtime benchmarks in every language.

### Phase 16 — Workflows

Measure coordination overhead, large workflow state, scheduling, recovery, and high fan-out/fan-in patterns.

### Phase 17 — Broader execution types

Add workload-specific benchmark scenarios only where new execution behaviour creates materially different pressure.

### Phase 18 — Distributed execution

If this phase remains justified, add network scheduling, worker coordination, remote state, cross-node backpressure, and distributed saturation testing.

---

## 9. Regression Strategy

Performance regression detection should become stricter only after benchmarks are stable enough to support it.

### 9.1 Pull-request benchmarks

Run a small deterministic benchmark subset on relevant changes.

Candidate cheap benchmarks include:

- lifecycle transition;
- run creation/completion;
- terminal resolution;
- event append/replay;
- important allocation-sensitive primitives.

Use statistical comparison tooling such as `benchstat` rather than comparing one raw run against another.

Initially, benchmark changes should be reported rather than automatically blocked.

### 9.2 Allocation regressions

Allocation-count changes often provide a cleaner signal than nanosecond changes on noisy CI hardware.

Changes such as:

```text
0 allocations/op -> 1
1 allocation/op  -> 3
```

should be visible and investigated when they occur on high-frequency paths.

They are not automatically defects; the question is whether the additional allocation is justified and material at expected concurrency.

### 9.3 Scheduled or release benchmarks

Heavier suites should run outside the normal PR critical path once they become expensive.

These runs can include:

- the full concurrency matrix;
- saturation tests;
- representative workload mixes;
- memory and runtime-task measurements;
- CPU, block, and mutex profiles;
- long-duration tests where appropriate.

### 9.4 Thresholds

Do not establish strict percentage gates before measuring normal benchmark variance on the actual CI environment.

A future policy may classify regressions into report, warning, and investigation thresholds, but those values must be evidence-based.

A percentage alone must not decide whether an optimization is worthwhile. A 50% increase from 1 microsecond to 1.5 microseconds may be immaterial, while a small increase in memory per active run may be critical at high concurrency.

---

## 10. Benchmark History

Aren should retain enough historical data to answer when and why a scaling characteristic changed.

At minimum, important benchmark runs should be attributable to:

- commit;
- compiler and runtime version;
- operating system;
- architecture;
- CPU or runner class;
- benchmark scenario;
- concurrency level.

The storage mechanism should remain simple initially.

Possible progression:

1. local benchmark output during Phase 1;
2. CI artifacts once automated comparisons are useful;
3. a lightweight machine-readable canonical baseline when historical comparison becomes valuable;
4. dedicated benchmark storage or visualization only if the volume justifies it.

Do not build a performance database merely because one may eventually be useful.

---

## 11. Profiling Workflow

Optimization begins only after a benchmark or real workload exposes a meaningful problem.

The default loop is:

```text
measure
  ↓
reproduce the regression or scaling knee
  ↓
profile
  ↓
form a hypothesis
  ↓
make one focused change
  ↓
remeasure
  ↓
rerun correctness, race, and stress tests
```

Use the selected language's benchmark runner, statistical comparison tools,
CPU and heap profilers, contention profiler, and scheduler tracing where
available. Record the exact commands with each reproducible benchmark.

Different symptoms require different evidence.

### CPU-heavy degradation

Use CPU profiling.

### Unexpected memory growth

Use allocation and heap profiles and compare active versus completed-run retention.

### Poor scaling without high CPU

Inspect contention profiles, runtime-task states, and scheduler traces.

### Tail-latency spikes

Correlate scheduler, contention, GC, buffering, and workload behaviour rather than assuming the slowest function in a CPU profile is responsible.

---

## 12. Performance Review in Every Phase

From Phase 2 onward, each phase review should answer:

1. What incremental runtime cost did this phase introduce?
2. Which existing canonical benchmarks changed materially?
3. Did allocation, memory, runtime-task, contention, or tail-latency behaviour change?
4. Where does the capability begin to degrade under increasing concurrency?
5. Did a new global coordination point appear?
6. Did the phase introduce any unbounded queue, retention, buffer, or runtime-task behaviour?
7. Does any regression matter for intended workloads, or is it benchmark noise/immaterial cost?
8. If optimization occurred, what profile demonstrated the bottleneck?
9. Did the optimization preserve all correctness and concurrency guarantees?
10. What new scaling risk should become a benchmark in the next phase?

The review may conclude that no optimization is warranted.

That is a valid and often desirable result.

---

## 13. Initial Phase 1 Benchmark Plan

Phase 1 should establish measurement without allowing performance work to expand the phase scope.

### Required baseline scenarios

The initial implementation should have a small benchmark set covering:

```text
BenchmarkRunSuccess
BenchmarkRunFailure
BenchmarkRunCancellation
BenchmarkLifecycleTransition
BenchmarkTerminalResolution
BenchmarkEventReplay
BenchmarkWaitCompleted
```

Names may differ with the final API.

### Concurrency scenarios

At minimum, exercise immediate or blocked runs at representative levels such as:

```text
1
10
100
1,000
5,000
10,000
```

Very high levels may be reduced on constrained CI hardware, but a heavier local or scheduled scenario should still exist when practical.

### Memory and runtime-task scenario

Create a large number of controlled blocked runs and measure:

- heap increase;
- approximate memory per active run;
- Aren-owned runtime-task increase;
- whether runtime tasks and retained memory return toward baseline after release and completion.

### Phase 1 interpretation

These numbers establish the first historical baseline.

They are not yet service-level objectives.

Phase 1 should not be delayed to shave nanoseconds unless measurement reveals one of the following:

- pathological allocation behaviour;
- super-linear resource growth;
- accidental global contention;
- runaway runtime-task creation;
- observer-induced execution slowdown or blockage;
- memory use large enough to undermine the intended high-concurrency direction;
- performance severe enough to distort race or stress evidence.

---

## 14. Future Numeric Budgets

Aren may eventually define explicit budgets such as:

- lifecycle p99 overhead;
- bytes retained per blocked run;
- maximum Aren-owned runtime tasks per active run;
- supported concurrent active runs on a reference machine;
- event-processing throughput;
- daemon control-path latency;
- saturation recovery time.

No values should be invented now.

Budgets become useful only after:

1. the relevant capability exists;
2. benchmark variance is understood;
3. realistic workloads have been exercised;
4. the cost matters to intended Aren usage.

The objective is not to win microbenchmarks. It is to preserve predictable and efficient supervision at the scale Aren actually needs.

---

## 15. Performance Success Model

Aren's performance discipline is working when:

- important runtime costs are measured before they become mysterious;
- every major capability has a representative benchmark;
- high-concurrency behaviour is understood rather than assumed;
- memory and runtime-task cost per active execution remain visible;
- global contention and unbounded retention are detected early;
- regressions can be tied to commits and capabilities;
- profiling evidence guides optimization;
- overload behaviour is explicit and eventually governed rather than accidental;
- correctness guarantees remain intact under benchmark and saturation load;
- performance work does not itself become speculative architecture.

The long-term standard is simple:

> Aren should be able to supervise many executions in parallel without the runtime becoming the dominant source of latency, memory pressure, contention, or instability.
