# Area reasoning template: Scale, resource, and profile evidence

## Purpose

Decide how Sprint 3 measures concurrency scaling, latency, allocation, memory, Aren-owned runtime tasks, contention, saturation, and cleanup. This area owns measurement validity and bounds, not performance targets invented before evidence exists.

## Required inputs

- accepted Sprint 3 workload method;
- completed Sprint 2 implementation and concurrency tests;
- performance engineering mandate;
- selected Go performance and concurrency evidence.

## Questions to resolve

1. Which canonical concurrency levels are safe and meaningful on the measurement host?
2. How are throughput and p50, p95, p99, and maximum latency calculated?
3. How are peak and retained heap separated, and when is garbage collection controlled or reported?
4. How are Aren-owned runtime tasks counted without confusing them with test or Go runtime tasks?
5. Which mutex, block, CPU, heap, goroutine, and trace captures answer specific questions?
6. What constitutes recovery after load falls?
7. Which result invalidates the run or blocks Phase 1 acceptance?

## Required analysis

Map every metric to its collection mechanism, measurement error, overhead, and interpretation limit. Define bounded setup, steady-state, release, and quiescence windows. Identify instrumentation that could perturb the measured lifecycle and state how results disclose that cost.

## Required output

- scenario and concurrency matrix;
- metric definitions and collection methods;
- resource baseline and quiescence method;
- profiling command and artifact matrix;
- saturation stop conditions;
- measurement overhead statement;
- blocker versus recorded-baseline classification.

