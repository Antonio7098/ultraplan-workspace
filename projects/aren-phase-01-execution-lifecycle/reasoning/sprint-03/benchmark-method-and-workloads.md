# Area reasoning template: Benchmark method and workloads

## Purpose

Define the performance experiment method that Phase 1 and later Aren phases will reuse. This area owns workload definitions, benchmark tiers, sampling, environment capture, comparison, and extension rules. It does not choose optimizations.

## Required inputs

- Phase 1 PRD and roadmap;
- performance engineering mandate;
- completed Sprint 1 and Sprint 2 artifacts and realised code;
- selected performance, testing, and evaluation evidence.

## Questions to resolve

1. Which real lifecycle entry points does each benchmark exercise?
2. How does each workload separate Aren overhead from controlled-work time?
3. Which parameters and concurrency levels are fixed, configurable, or host-limited?
4. What warm-up, duration, repetition, and sample rules make comparison honest?
5. Which commands belong in quick local, pull-request, and extended tiers?
6. Which environment and source facts make a result reproducible?
7. How will later phases extend workloads without changing old benchmark meaning?

## Required analysis

Define immediate success, returned failure, cancellation, completed wait, event replay, concurrent completion, and blocked active-run workloads. For each, record setup, measured interval, work performed inside that interval, parameters, outputs, expected invariants, cleanup, and invalid-result conditions.

Compare Go benchmarks with a bounded scenario runner where latency distributions or process-level resource measurements require one. Keep one authority for workload definitions and avoid two implementations of the lifecycle path.

## Required output

- workload catalog and naming rules;
- tier and command matrix;
- sample and comparison policy;
- environment metadata contract;
- invalid measurement rules;
- later-phase extension rule;
- rejected alternatives and reopen conditions.

