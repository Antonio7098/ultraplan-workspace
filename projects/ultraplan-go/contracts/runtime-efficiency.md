# Runtime Efficiency Contract

## Purpose

This contract governs the execution shape of UltraPlan model work. It requires each workload to use the smallest executor that can meet its quality, safety, and recovery requirements, with bounded runtime behavior and measurable workflow cost.

## Scope

Select this contract when a sprint changes:

- model or provider routing
- agent, one-off API call, or deterministic worker boundaries
- session creation, continuation, repair, or compaction
- tool exposure, fan-out, arbitration, retries, fallbacks, or queueing
- timeouts, cancellation, concurrency, token budgets, or cost controls
- runtime events, usage accounting, or efficiency evaluation

## Requirement Index

| ID | Title | Severity If Violated |
|---|---|---|
| EFF-INVENTORY-001 | Every runtime operation must be inventoried | High |
| EFF-EXECUTOR-001 | Each workload must use the smallest sufficient executor | High |
| EFF-MODEL-001 | Model and reasoning choices must be evaluation-backed | High |
| EFF-SESSION-001 | Compatible work must reuse session state deliberately | Medium |
| EFF-TOOLS-001 | Tool access and loops must be minimal and bounded | High |
| EFF-FANOUT-001 | Parallelism must have a measured benefit and bounded width | High |
| EFF-FAILURE-001 | Retry, repair, fallback, and cancellation must be bounded | High |
| EFF-METRIC-001 | Efficiency must be measured per successful workflow | Medium |

## Requirements

### EFF-INVENTORY-001: Every Runtime Operation Must Be Inventoried

**Rule**

All work that can invoke a model, tool loop, subprocess, or paid provider must appear in the runtime map.

**Required**

- include primary calls, fan-out workers, arbiters, challengers, repairs, retries, fallbacks, summarizers, and smoke-authoring calls
- record trigger, owner, executor kind, model settings, context contract, tools, output contract, attempts, and expected consumer
- calculate worst-case invocation count for one user operation
- identify deterministic stages and validation paths that invoke no model

**Forbidden**

- hidden provider calls in fallback or error handling
- treating a multi-call stage as one operation in cost and latency estimates

**Evidence**

- tests or diagnostics expose ordered call sequence, role, operation, and attempt

### EFF-EXECUTOR-001: Each Workload Must Use The Smallest Sufficient Executor

**Rule**

Choose the executor from the actual work required.

**Required**

- use deterministic code for parsing, sorting, validation, normalization, fingerprinting, state transitions, queueing, selection by declared references, and repeatable aggregation
- use a one-off model API call when all evidence can be supplied up front, tools are unnecessary, output is bounded, and one validation or repair cycle is sufficient
- use an agent when the model must explore an unknown environment, choose tools from observations, modify several files, verify iteratively, or recover from open-ended state
- keep permissions, retries, state transitions, validation, and publication in deterministic orchestration
- document why an agent remains necessary for any no-tool or single-turn task

**Forbidden**

- using model judgment for deterministic state changes
- starting a general agent runtime for a bounded transformation without a reason
- allowing the model to own artifact promotion or authorization

**Evidence**

- the reasoning record contains an executor decision for every runtime workload

### EFF-MODEL-001: Model And Reasoning Choices Must Be Evaluation-Backed

**Rule**

Model, provider, variant, reasoning effort, and output verbosity must match task difficulty and be changed through representative evaluation.

**Required**

- establish quality, structured-output validity, latency, token, and cost baselines before changing runtime selection
- test lower-cost models or reasoning settings on representative successful and failure cases
- pin settings that materially affect behavior or cache compatibility
- change one major variable at a time when comparing model, prompt, tool, or reasoning behavior
- retain a rollback path for production defaults

**Forbidden**

- routing only by per-token price
- lowering model capability or reasoning effort without acceptance evidence
- relying on provider defaults when stable latency, quality, or cache cohorts require explicit settings

**Evidence**

- eval or smoke results state the tested configuration and acceptance thresholds

### EFF-SESSION-001: Compatible Work Must Reuse Session State Deliberately

**Rule**

Related ordered work should reuse a session when earlier context and observations remain valid.

**Required**

- send the governed shared context once at session start
- continue with the current task, changed facts, and validation findings
- keep repair turns in the originating compatible session
- preserve ordered task execution where later tasks depend on earlier repository changes
- fall back to a fresh session with the complete required context when continuation is unavailable or incompatible
- start fresh when permissions, work directory, provider, model compatibility, source revision, policy, or objective changes materially

**Forbidden**

- resending the full prompt for every task in an ordered queue by default
- relying on session memory without persisted task and outcome records
- continuing a session after its trust or freshness assumptions become invalid

**Evidence**

- tests cover compact continuation, compatible repair, and fresh-session fallback

### EFF-TOOLS-001: Tool Access And Loops Must Be Minimal And Bounded

**Rule**

Expose only the tools needed for the workload and bound every model-controlled loop.

**Required**

- define allowed capabilities, maximum turns, tool calls, context growth, duration, and output size
- keep stable tool definitions and ordering within a reusable prompt cohort
- prefer deterministic orchestration for known reads, writes, validation, and command sequences
- record tool calls by kind and outcome
- cancel tool work with its owning operation

**Forbidden**

- broad shell, network, or write access for a read-only judgment
- unbounded tool loops or recursive discovery
- repeated reads of evidence already delivered and still current without a stated verification reason

**Evidence**

- policy tests reject unsupported capabilities and runtime metrics expose tool use

### EFF-FANOUT-001: Parallelism Must Have A Measured Benefit And Bounded Width

**Rule**

Parallel model work must reduce elapsed time or improve decision quality enough to justify added calls, context, reconciliation, and provider pressure.

**Required**

- define fan-out width, queueing, ownership, cancellation, partial-failure behavior, and arbiter inputs
- give each worker only its shared foundation and assigned evidence
- deduplicate shared context and schedule compatible calls within useful provider cache windows when practical
- compare parallel and sequential behavior under representative rate and cache conditions
- preserve deterministic result ordering and conflict handling

**Forbidden**

- unbounded agent spawning
- sending all worker-specific evidence to every worker
- adding parallel calls where serial session reuse has lower end-to-end time or cost at equivalent quality

**Evidence**

- measurements include fan-out width, queue delay, cache behavior, provider latency, and reconciliation cost

### EFF-FAILURE-001: Retry, Repair, Fallback, And Cancellation Must Be Bounded

**Rule**

Failure handling must minimize repeated work and terminate predictably.

**Required**

- set maximum attempts, timeout, backoff, and retryable error categories
- distinguish transport retry, model repair, session continuation failure, provider fallback, and user cancellation
- send focused validation findings for repair
- prevent a retry or fallback from duplicating already committed external effects
- preserve the last valid governed artifact on failure
- propagate cancellation to model, tools, subprocesses, and queued work

**Forbidden**

- retrying deterministic validation failures with an unchanged request
- silent fallback that changes model behavior, permissions, cost, or output contract
- abandoning owned work after cancellation or timeout

**Evidence**

- failure-path tests prove attempt limits, idempotency, cancellation, and artifact preservation

### EFF-METRIC-001: Efficiency Must Be Measured Per Successful Workflow

**Rule**

Efficiency claims must compare complete successful workflows at equivalent quality.

**Required**

- record model, provider, settings, session mode, call role, queue position, attempts, tool calls, prompt bytes, prefix bytes and digest, input tokens, cached-input tokens, cache-write tokens, output tokens, reasoning tokens, duration, and status where available
- report time to first useful output and end-to-end duration for user-visible operations
- calculate realized cost per successful workflow, including failed attempts, repairs, arbiters, and fallbacks
- track output validation, acceptance coverage, defect escape, timeout, and malformed-output rates beside efficiency measures
- compare representative before and after cohorts

**Forbidden**

- claiming an improvement from fewer calls while success or evidence quality falls
- reporting cache keys as proof of cache hits
- excluding failed attempts or repair cost from workflow totals

**Evidence**

- a performance or evaluation record contains baseline, changed result, cohort definition, and acceptance decision

## Review Rejection Criteria

Reject an affected change if it:

- leaves runtime calls, repairs, or fallbacks outside the call inventory
- uses an agent for deterministic work without justification
- changes model or reasoning settings without representative evaluation
- reconstructs compatible context for each task instead of using continuation
- exposes excessive tools or unbounded loops
- adds fan-out without bounded ownership and measured benefit
- permits retries or fallback to multiply cost invisibly
- reports request-level savings without successful-work quality, latency, and cost

