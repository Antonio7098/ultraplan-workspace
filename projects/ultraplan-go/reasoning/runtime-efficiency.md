# Runtime Efficiency Reasoning Guide

Use this guide when a sprint changes model calls, agents, tools, session behavior, fan-out, retries, model routing, latency, or cost. Analyze the complete user operation, including failure paths, before choosing an optimization.

## Required Analysis

### 1. Runtime Call Graph

List deterministic work and every possible runtime call in execution order.

| Sequence | Trigger | Operation and role | Executor | Model settings | Session mode | Tools | Max calls | Output consumer |
|---|---|---|---|---|---|---|---:|---|
| `1` | `<...>` | `<stage/worker/arbiter/repair>` | `<code/API/agent/subprocess>` | `<provider/model/reasoning/verbosity>` | `<fresh/continue/none>` | `<...>` | `<n>` | `<...>` |

Include:

- primary stage calls
- fan-out workers and challengers
- arbitration and consolidation
- schema or semantic repairs
- transport retries and provider fallback
- smoke authoring and external subprocesses
- calls suppressed by deterministic validation

Calculate best-case, expected, and worst-case call counts for one user operation.

### 2. Smallest-Executor Decision

Classify each workload independently.

| Workload | Rules deterministic? | Evidence complete up front? | Adaptive tools needed? | Side effects? | Chosen executor | Reason |
|---|---|---|---|---|---|---|
| `<...>` | `<yes/no>` | `<yes/no>` | `<yes/no>` | `<...>` | `<code/API/agent>` | `<...>` |

Use deterministic Go for parsing, validation, canonical ordering, selection by declared references, fingerprints, state transitions, queueing, retry policy, and artifact promotion.

Use one bounded model API call when the evidence fits in the request, the task needs judgment, tools add no value, output is schema-bound, and at most one focused repair is expected.

Keep an agent for unknown-environment exploration, adaptive tool selection, multi-file implementation, iterative verification, or recovery from open-ended state. Identify the exact behavior that requires agency.

### 3. Session And Repair Design

For related ordered work, answer:

- Can one session receive the governed context once and process an ordered queue?
- Which facts or repository changes from earlier tasks help later tasks?
- What compact continuation payload identifies the current task and changed facts?
- Can validation repair continue in the same session?
- Which provider, model, policy, permission, worktree, or revision changes force a fresh session?
- What complete prompt is used when continuation fails?
- How are task and outcome records persisted independently of session memory?

For execute work, assess shared-session execution before independent agents. Keep full details for the current task; give the queue enough summary to preserve order and cross-task awareness.

### 4. Tool And Loop Budget

Record for every agent-backed call:

| Call | Allowed capabilities | Max turns | Max tool calls | Timeout | Max context/output | Cancellation owner |
|---|---|---:|---:|---|---|---|
| `<...>` | `<...>` | `<n>` | `<n>` | `<duration>` | `<...>` | `<...>` |

Move known reads, writes, command sequences, and validation into deterministic orchestration when that reduces turns without weakening evidence or permission controls.

### 5. Fan-Out And Scheduling

For each parallel group, define:

- worker roles and disjoint evidence ownership
- maximum width and queue behavior
- shared prefix and worker-specific suffix
- cancellation and partial-failure handling
- deterministic result order
- reconciliation or arbiter cost
- provider rate-limit and prompt-cache effects
- measured sequential alternative

Consider a brief leader-first cache schedule only when real requests share a cacheable prefix and measurements show net latency or cost benefit. Do not add a synthetic warm-up call by default.

### 6. Model And Reasoning Route

Define the minimum quality bar before selecting a cheaper model or lower reasoning setting.

| Workload | Baseline configuration | Candidate | Quality/eval result | P50/P95 duration | Tokens and cache | Cost per success | Decision |
|---|---|---|---|---|---|---|---|
| `<...>` | `<...>` | `<...>` | `<...>` | `<...>` | `<...>` | `<...>` | `<...>` |

Pin settings that affect behavior and cache compatibility. Compare one major change at a time and retain rollback for production defaults.

### 7. Failure Economics

Map transport retries, validation repairs, session-continuation failure, provider fallback, timeout, cancellation, and malformed output separately.

For each path, state:

- retryable classification and maximum attempts
- repeated input and tool work
- focused repair payload
- idempotency and external-effect protection
- last-valid-artifact behavior
- worst-case added latency and cost
- operator-visible outcome

### 8. Measurement Plan

Measure complete successful workflows at equivalent quality. Capture:

- output validation and acceptance-criteria coverage
- defect escape, retry, timeout, and malformed-output rates
- model, provider, reasoning, verbosity, and output format
- prompt bytes, stable-prefix bytes, prefix digest, and cache cohort
- fresh-input, cached-input, cache-write, output, and reasoning tokens
- session mode, turns, tool calls by kind, queue delay, and attempts
- time to first useful output and end-to-end duration
- realized cost including repairs, fallbacks, workers, and arbiters

Use representative workflows. Separate measured values, deterministic byte comparisons, and estimates.

## Decisions To Produce

The area reasoning must include these decisions under the standard `## Area Decisions` heading:

1. runtime call graph and worst-case invocation count
2. executor for every workload
3. model and reasoning route with evaluation threshold
4. session, continuation, and repair policy
5. tool and loop budgets
6. fan-out width and scheduling policy
7. retry, fallback, timeout, and cancellation behavior
8. workflow-level measurement and rollout gate

Under `## Trade-Offs`, compare quality, latency, cost, operational complexity, provider limits, and recovery. Under `## Evidence`, include runtime paths, request traces, evals, and before/after measurements. Under `## Risks`, cover quality loss, cache misses, context growth, rate limits, runaway tools, retry multiplication, and stale session state.

## Rejection Questions

- Is a model doing work that deterministic code can produce exactly?
- Is a no-tool bounded judgment paying for a general agent loop without evidence?
- Does each ordered task reconstruct a compatible shared context?
- Can a repair continue from valid prior context?
- Does fan-out improve elapsed time or quality after reconciliation and provider limits?
- Are tools, turns, attempts, time, concurrency, and context growth bounded?
- Are failed attempts and fallback included in cost per success?
- Is a cheaper model or reasoning setting supported by representative evaluation?

