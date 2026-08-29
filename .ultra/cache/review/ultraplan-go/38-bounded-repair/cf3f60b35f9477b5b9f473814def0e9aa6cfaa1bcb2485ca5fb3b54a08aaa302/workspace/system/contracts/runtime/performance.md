# Performance Contract

## Purpose

This contract defines how performance-sensitive changes must be designed, measured, and reviewed.

It governs:

- latency and throughput assumptions
- CLI/service/UI responsiveness where applicable
- query, scan, runtime, and provider cost
- caching
- concurrency
- resource limits
- profiling-before-optimization discipline

## Scope

This contract applies when a change:

- adds user-facing runtime paths
- changes hot paths, queries, filesystem scans, loops, rendering, caching, runtime/provider/model calls, batch jobs, or bundle size
- introduces concurrency, streaming, polling, retries, or expensive computation

## Requirement Index

| ID | Title | Applies To | Severity If Violated |
| --- | --- | --- | --- |
| PERF-BUDGET-001 | Important paths must declare performance expectations | hot/user-facing paths | Medium |
| PERF-MEASURE-001 | Optimization must be evidence-driven | performance work | Medium |
| PERF-BOUND-001 | Expensive work must be bounded | APIs, CLIs, scans, jobs, UI | High |
| PERF-UI-001 | UI-facing changes must consider loading, interaction, and layout stability | frontend/TUI/status UI | High |
| PERF-CACHE-001 | Caching must define correctness and invalidation | caches | High |
| PERF-CONC-001 | Concurrency must be safe, bounded, and observable | async/parallel paths | High |
| PERF-COST-001 | Runtime/provider/cloud-cost paths must expose cost drivers | AI/runtime/provider paths | Medium |

## Core Principles

1. Correct and clear first; optimize with evidence.
2. Known scale constraints are design inputs, not afterthoughts.
3. Every unbounded operation is a future incident.
4. Caches trade freshness for speed; define the trade.
5. Concurrency improves throughput only when bounded and safe.
6. UI and status-output performance are user experience.
7. Runtime/provider/model cost is operational performance.

## Requirements

### PERF-BUDGET-001: Important Paths Must Declare Performance Expectations

**Rule**
Important user-facing, runtime-backed, batch, provider-backed, or high-volume paths must have explicit performance expectations.

**Required**
- identify latency/throughput expectations where meaningful
- define expected data size, repository size, task count, or request volume assumptions
- document when performance is intentionally not optimized yet

**Forbidden**
- adding high-cost paths with no volume or latency assumptions
- pretending small-test performance proves production behaviour

**Evidence**
- reasoning docs identify performance assumptions for non-trivial paths

### PERF-MEASURE-001: Optimization Must Be Evidence-Driven

**Rule**
Performance optimizations should be justified by profiling, benchmark, telemetry, or known constraints.

**Required**
- measure before complex optimization where practical
- keep optimized code isolated and documented if it reduces clarity
- add regression benchmarks for critical hot paths where useful

**Forbidden**
- introducing obscure or fragile optimizations without evidence
- rejecting simple clear code because of hypothetical bottlenecks

**Evidence**
- benchmark/profiling/telemetry notes support meaningful optimization work

### PERF-BOUND-001: Expensive Work Must Be Bounded

**Rule**
Loops, queries, filesystem scans, file processing, runtime/provider calls, UI rendering, jobs, and batch work must have practical bounds.

**Required**
- define max page size, batch size, file size, directory traversal scope, retry count, timeout, concurrency, and memory limits where relevant
- stream or chunk large work when needed
- expose progress for long-running user-visible work

**Forbidden**
- unbounded list endpoints, recursive filesystem scans, or runtime/provider loops
- loading large files/datasets fully into memory without reason
- infinite retry/polling loops

**Evidence**
- tests/config/review notes show relevant bounds

### PERF-UI-001: UI-Facing Changes Must Consider Loading, Interaction, And Layout Stability

**Rule**
Frontend, TUI, and status-output changes must avoid unnecessary harm to loading speed, responsiveness, readability, and visual stability.

**Required**
- avoid unnecessary large dependencies and bundles
- split code where route/feature boundaries justify it
- prevent layout shift with stable dimensions/skeletons where practical
- avoid repeated unnecessary renders, status redraws, filesystem scans, and duplicate network requests
- handle loading states without blank hangs

**Forbidden**
- adding heavy packages for minor UI helpers
- blocking initial render with avoidable work
- causing layout jumps from late-loaded content without mitigation

**Evidence**
- bundle, redraw, status-output, or UI performance impact is considered for meaningful UI changes

### PERF-CACHE-001: Caching Must Define Correctness And Invalidation

**Rule**
Caches must declare what they cache, why, for how long, and when they invalidate.

**Required**
- define cache key, TTL, invalidation, source of truth, and stale behaviour
- avoid caching sensitive data unless justified and protected
- distinguish local memoization, file cache, HTTP cache, app cache, DB cache, and derived stores

**Forbidden**
- adding caches to hide slow code without correctness rules
- caching authorization-sensitive responses under insufficient keys
- treating stale cached data as authoritative accidentally

**Evidence**
- cache tests or docs cover invalidation/staleness for important paths

### PERF-CONC-001: Concurrency Must Be Safe, Bounded, And Observable

**Rule**
Parallelism, async tasks, workers, goroutines, promises, and batch execution must have ownership, limits, cancellation, and failure visibility.

**Required**
- set concurrency limits where fan-out can grow
- use structured ownership for spawned work so cancellation and failures flow back to the caller or owning supervisor
- propagate cancellation/timeouts where possible
- connect process, command, request, or workflow cancellation signals to long-running work where applicable
- collect and surface partial failures deliberately
- avoid shared mutable state races
- recover or surface panics/fatal failures from worker execution rather than letting them disappear silently

**Forbidden**
- unbounded goroutines/promises/tasks
- fire-and-forget work without ownership when the result matters
- swallowing parallel task failures
- launching workers that cannot be stopped, joined, reconciled, or observed

**Evidence**
- concurrency tests or runtime safeguards cover important paths
- review can identify worker ownership, cancellation propagation, failure collection, and concurrency limits

### PERF-COST-001: Runtime/Provider/Cloud-Cost Paths Must Expose Cost Drivers

**Rule**
Runtime paths that consume paid providers, models, external services, storage, or compute must expose meaningful cost drivers.

**Required**
- track token count, model/provider, request count, duration, payload size, storage growth, or job count where relevant
- set cost-related limits for abuse-prone paths
- avoid hidden repeated provider calls

**Forbidden**
- unbounded LLM/runtime/provider calls from one command, request, or task
- retry/fallback behaviour that multiplies cost invisibly

**Evidence**
- logs/events/metrics include cost-relevant metadata where meaningful

## Review Rejection Criteria

Reject a performance-sensitive change if it:

- adds unbounded expensive work
- introduces heavy frontend/TUI/status dependencies without justification
- adds caching with no invalidation/correctness model
- creates unbounded concurrency or hidden owned work
- introduces concurrent work with no cancellation, ownership, or failure propagation model
- optimizes complexly with no evidence
- exposes high-cost provider/model paths with no limits or cost visibility

---
