# Context Engineering Reasoning Guide

Use this guide when a sprint changes what any model receives, discovers, reuses, trusts, or passes to another stage. The resulting area reasoning must make the context route reviewable from authoritative source to model call to validated artifact.

## Current UltraPlan Baseline

UltraPlan's canonical sprint input contracts are owned by `internal/sprint/input_contract.go`. Confirm this map against the implementation before using it in a decision.

| Stage | Required prepared foundation | Stage-specific context | Expected exploration |
|---|---|---|---|
| Requirements | project index, roadmap, project docs | prior sprint reviews when selected | planning evidence only |
| Code context | accepted requirements | resolved implementation repository | broad, read-only repository exploration owned by this stage |
| Sprint index | requirements, code context, resolved source evidence | project index, roadmap, project docs | targeted verification only |
| Technical handbook | shared foundation | sprint index and selected evidence | evidence interpretation |
| Area reasoning | shared foundation | project docs, sprint index, handbook, selected context files, one selected area guide | area-specific verification |
| Consolidated reasoning | shared foundation | project/roadmap/docs, sprint index, handbook, selected context files, all completed area reasoning | reconcile decisions and conflicts |
| Plan | shared foundation | indexes, docs, handbook, area reasoning, consolidated reasoning, plan template | task decomposition |
| Execute | shared foundation once per compatible session | plan queue plus full current-task detail | live inspection, edits, and verification required by the current task |
| Review | shared foundation | one coverage source, governed review inputs, changed target files | read-only conformance inspection |
| QA | shared foundation | execution evidence, review outcome, changed target files, approved checks | bounded read-only checks and arbitration |
| Smoke | shared foundation | planning decisions, execution evidence and handoff, review outcome, smoke harness | bounded harness authoring and invocation |
| Merge | workspace record and Git identities | changed paths and QA outcome | deterministic Git validation and integration |

The shared foundation currently consists of exact `requirements.md`, exact `code-context.md`, and resolved source selections. It is a derived, content-addressed, disposable rendering cache. The governed artifacts remain authoritative.

## Required Analysis

### 1. Call And Input Ledger

Enumerate every call affected by the sprint, including repairs, retries, workers, arbiters, and fallbacks.

| Call | Purpose and success condition | Required inputs | Optional inputs | Forbidden inputs | Delivery mode | Tools and permissions | Output and validator |
|---|---|---|---|---|---|---|---|
| `<stage/role>` | `<...>` | `<...>` | `<...>` | `<...>` | `<inline/excerpt/retrieval/live/session>` | `<...>` | `<...>` |

For each input, answer:

- Who produces and owns it?
- Does the model receive exact bytes, a selected excerpt, structured data, a path, retrieval access, or prior session state?
- What revision, digest, or governed artifact establishes freshness?
- Is it authoritative, derived, generated, untrusted, or transient?
- Which other calls receive the same bytes?
- What discovery or interpretation remains for the consumer?

### 2. Code-Context Ownership

Start from the accepted requirements and list the questions repository inspection must answer. Then assess the selected evidence.

| Requirement or decision | Path and symbol | Exact range or full file | Relationship exposed | Selection reason | Remaining uncertainty |
|---|---|---|---|---|---|
| `<AC/decision>` | `<path:symbol>` | `<lines/full>` | `<caller, data flow, config, error, test>` | `<why downstream needs it>` | `<...>` |

Check that the code-context stage covers:

- entry points and ownership boundaries
- callers, callees, data flow, and persisted state
- configuration and runtime selection
- success, failure, cancellation, and recovery behavior
- validation and tests
- partial implementations, contradictions, and notable absences

Record repeated search or read work that downstream calls still perform. Distinguish required live verification from exploration the code-context stage should have resolved.

### 3. Delivery And Deduplication

For each required input, choose one primary delivery mode and explain it.

- Inline full content for small governed artifacts whose exact wording matters.
- Validated selected source for decision-relevant code.
- Structured packet for normalized facts consumed by several calls.
- Retrieval or live repository access for optional, volatile, or open-ended evidence.
- Session state for compatible ordered work after the initial governed packet.

Identify content repeated in the shared prefix, stage suffix, tool result, and conversation history. State which copy will be removed and why the remaining source is sufficient.

### 4. Freshness And Trust

Define:

- authoritative inputs and their identity
- derived context cache key and invalidation inputs
- repository or worktree identity
- behavior when selected source changes during a read
- behavior when execution changes line numbers after the shared foundation is frozen
- conditions that require targeted live inspection or a fresh session
- evidence framing and instruction hierarchy
- redaction and cache-sharing limits for sensitive content

### 5. Prefix And Provider Cache Design

Verify current provider documentation before choosing cache controls. Record the provider, model, endpoint, region or retention boundary, tool set, output schema, reasoning settings, and any provider-specific cache parameters.

Lay out the rendered request in byte order:

1. stable platform policy and tool definitions
2. versioned workflow and task-family instructions
3. governed project or sprint foundation
4. frozen selected source evidence
5. shared cache boundary
6. stage-common material and an optional second boundary
7. volatile stage, task, run, attempt, and diagnostic data

Document:

- exact prefix cohort and version
- boundary byte offset and prefix digest
- fields excluded from the prefix
- whether later turns append or rewrite earlier context
- cache key or routing hint semantics
- minimum reusable length and retention assumptions
- expected reuse count and scheduling window
- cached-input, cache-write, latency, and realized-cost measurements

Provider caching is an acceleration layer. Describe the same correct workflow under a cache miss.

## Decisions To Produce

The area reasoning must include these decisions under the standard `## Area Decisions` heading:

1. input contract changes for each affected call
2. owner and lifecycle of shared code context
3. delivery mode for each required input
4. deduplication and downstream discovery removed
5. freshness identity and invalidation behavior
6. trust and permission boundaries
7. prefix cohort, ordering, and cache boundaries
8. telemetry and evaluation needed to prove the design

Under `## Trade-Offs`, compare prompt size, avoided tool work, freshness, cache write cost, latency, and evidence quality. Under `## Evidence`, cite implementation paths, rendered prompt tests, and measured runs. Under `## Risks`, cover stale context, omitted evidence, oversized context, cache fragmentation, prompt injection, and quality regression.

## Rejection Questions

- Does any downstream agent repeat repository exploration already owned by code context?
- Is a required input only mentioned by path?
- Is the same evidence delivered more than once?
- Can a volatile field change bytes before the intended cache boundary?
- Does the cohort ignore tools, output schema, model, or reasoning settings?
- Can stale or untrusted content become authoritative?
- Would a cache miss change correctness or rerun an authoritative stage?
- Can reviewers measure discovery avoided and quality preserved?

