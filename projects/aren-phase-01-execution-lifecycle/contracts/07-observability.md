# Aren Observability Contract

> Project: `aren-phase-01-execution-lifecycle`  
> Status: active project-wide engineering contract; intended to survive Phase 1 where the semantics remain valid

## Purpose

This contract defines the observability properties every Aren-owned capability must preserve.

Aren owns and supervises execution. It must therefore be possible to explain Aren-owned behaviour from structured evidence without guessing, parsing prose as semantic truth, or allowing observation machinery to control execution.

This contract is broader than `05-events-observation-and-waiting.md`. That contract governs canonical Phase 1 lifecycle history, replay, cursors, and waiting. This contract governs the wider evidence model: canonical semantic facts, diagnostic telemetry, structured logs, metrics, profiles, state snapshots, correlation, retention, decision evidence, observation surfaces, and the phase-level explainability gate.

It does not require a universal telemetry schema, durable observability backend, global event bus, exporter framework, or production monitoring stack before those capabilities are earned.

## Scope

Apply this contract whenever a change introduces or materially changes:

- Aren-owned execution behaviour;
- a decision that materially affects execution semantics;
- canonical events or state;
- diagnostics, logs, metrics, traces, profiles, or snapshots;
- observation transports, CLIs, browser surfaces, or exporters;
- nested executions or correlation identities;
- high-volume observable data;
- failure or resource evidence;
- observability retention or redaction.

## Requirement Index

| ID | Title | Severity If Violated |
| --- | --- | --- |
| AREN-EVID-001 | Aren-owned behaviour must be explainable | Blocker |
| AREN-EVID-002 | Canonical truth and diagnostics remain distinct | Blocker |
| AREN-EVID-003 | Observation must remain passive | Blocker |
| AREN-EVID-004 | Observation surfaces are projections, not authorities | Blocker |
| AREN-EVID-005 | Evidence uses the smallest useful semantic granularity | High |
| AREN-EVID-006 | Correlation preserves meaningful execution relationships | High |
| AREN-EVID-007 | Aren-owned decisions are inspectable where meaningful | High |
| AREN-EVID-008 | Retention is proportional to semantic importance and volume | High |
| AREN-EVID-009 | High-volume observation is explicitly bounded | Blocker |
| AREN-EVID-010 | Semantic truth must not exist only in prose logs | Blocker |
| AREN-EVID-011 | State snapshots agree with canonical state | High |
| AREN-EVID-012 | Observation cost is itself observable and attributable | High |
| AREN-EVID-013 | Sensitive payload capture is explicit and minimal | High |
| AREN-EVID-014 | Every phase closes an observability gate | High |

## Requirements

### AREN-EVID-001 — Aren-Owned Behaviour Must Be Explainable

Every materially meaningful Aren-owned operation must leave enough structured evidence, where applicable, to determine:

- that it occurred;
- when it occurred or the interval it occupied;
- which run or enclosing operation it belonged to;
- what requested or caused it when Aren owns that fact;
- the relevant before/after state;
- how it terminated;
- its failure or disposition;
- the materially relevant resources it consumed;
- whether the evidence is canonical truth or diagnostic support.

Opacity in behaviour Aren owns is an engineering defect unless the missing evidence is explicitly shown to be impossible, misleading, unsafe, or outside Aren's ownership.

### AREN-EVID-002 — Canonical Truth And Diagnostics Remain Distinct

Aren must distinguish evidence classes rather than forcing them into one universal representation.

At minimum, reason separately about:

- **canonical semantic evidence** — facts that constitute or directly record Aren-owned execution truth;
- **diagnostic telemetry** — measurements and internal runtime evidence useful for investigation but not semantic authority;
- **structured diagnostic logs** — contextual diagnostic records that must not be the only representation of semantic truth;
- **performance evidence** — measurements of cost, scaling, contention, saturation, and recovery;
- **profiles** — implementation-level CPU, heap, goroutine, mutex, block, or runtime evidence;
- **state snapshots** — coherent answers to what is true now.

These forms may share identities and timing but need not share one schema, retention rule, delivery mechanism, or storage model.

### AREN-EVID-003 — Observation Must Remain Passive

A slow, absent, abandoned, overloaded, disconnected, or faulty observer must not determine whether Aren can:

- execute work;
- accept or propagate cancellation;
- commit lifecycle transitions;
- resolve terminal state;
- release waiters;
- release Aren-owned resources.

Observation may consume resources, and those costs must be bounded and measured where material, but observer progress must not become a correctness dependency.

### AREN-EVID-004 — Observation Surfaces Are Projections, Not Authorities

CLIs, browser frontends, local services, exporters, dashboards, and future operator surfaces consume Aren-owned evidence.

They must not:

- become a second lifecycle model;
- infer canonical state from log text when canonical state exists;
- independently decide terminal outcomes;
- mutate canonical history to make a presentation convenient;
- require runtime semantics that exist only for the UI;
- silently reconcile disagreeing runtime facts into a fabricated answer.

If projections disagree, investigate the evidence boundary rather than choosing a presentation as the source of truth.

### AREN-EVID-005 — Evidence Uses The Smallest Useful Semantic Granularity

Instrument meaningful semantic boundaries, not every implementation instruction.

Good evidence granularity may include, where applicable:

```text
run created
transition committed
cancellation accepted
cancellation signal propagated
work returned
terminal candidate resolved
terminal transition committed
waiters released
```

An internal mutex acquisition normally does not deserve a canonical semantic event. Contention may instead belong in diagnostic telemetry or profiling when it materially affects behaviour.

Do not create an event merely because something can be observed. Choose event, interval/span, metric, log, profile, or state snapshot according to the question being answered.

### AREN-EVID-006 — Correlation Preserves Meaningful Execution Relationships

Observable evidence must carry stable identities sufficient to relate it to the Aren-owned operation it describes.

Phase 1 primarily uses `run_id`. Later phases may earn identities such as attempt, model invocation, tool call, child run, workflow run, or workflow step.

Requirements:

- introduce identities only when the capability exists;
- preserve parent/child or causal relationships when they are part of Aren's semantics;
- avoid inventing unrelated identifiers for one semantic operation;
- do not flatten meaningful governed child operations into an opaque parent success record.

### AREN-EVID-007 — Aren-Owned Decisions Are Inspectable Where Meaningful

When Aren makes a decision that materially changes execution semantics, the evidence model should expose enough structured information to understand the decision.

Examples include later decisions about:

- terminal resolution;
- retry eligibility and delay;
- permission or policy denial;
- admission or resource reservation;
- context compaction;
- workflow routing;
- recovery or uncertainty handling.

Prefer structured decision inputs, selected rule/policy, and resulting disposition over retrospective prose explanation.

Do not manufacture explanations for decisions Aren did not actually make or inputs it did not actually record.

### AREN-EVID-008 — Retention Is Proportional To Semantic Importance And Volume

Observability and retention are separate decisions.

Classify retained evidence according to its role:

- **canonical** — required to preserve Aren's promised semantic truth for the relevant lifetime;
- **diagnostic** — useful for explanation or debugging but not required to establish canonical truth;
- **ephemeral** — high-volume or continuously sampled evidence that may be aggregated, coalesced, sampled, or discarded.

Do not generalize Phase 1's complete small lifecycle history to token deltas, tool output, progress streams, profiling samples, queue-depth samples, or other high-volume telemetry.

Retention policy must follow semantic importance, volume, safety, and demonstrated use.

### AREN-EVID-009 — High-Volume Observation Is Explicitly Bounded

Any observation mechanism whose volume can grow materially must declare its bounds or a measured reason for not having one.

Relevant bounds may include:

- retained item count or bytes;
- buffer capacity;
- sampling/coalescing policy;
- payload size;
- subscriber or client count;
- update rate;
- history window;
- serialization or rendering limits.

Unbounded observability is not justified by the requirement for complete explainability.

When a bound can discard detail, the system must remain explicit about which evidence is canonical and whether any semantic guarantee is affected.

### AREN-EVID-010 — Semantic Truth Must Not Exist Only In Prose Logs

Logs are diagnostic context, not the sole representation of lifecycle or execution truth.

Determining facts such as successful completion, failure classification, cancellation disposition, permission outcome, or later recovery status must not require parsing human log messages when Aren owns a structured semantic representation.

Structured logs should preserve correlation and useful diagnostics, but log wording is not a lifecycle protocol.

### AREN-EVID-011 — State Snapshots Agree With Canonical State

A state snapshot or diagnostics view answers what is true at one coherent observation point.

It must not become an independently mutable truth store.

Where a snapshot includes related lifecycle facts, they must be captured through a coherence model compatible with the governing lifecycle/publication contracts.

Separate reads at different times may legitimately show different committed prefixes. A single claimed snapshot must not mix impossible combinations.

### AREN-EVID-012 — Observation Cost Is Itself Observable And Attributable

Observability machinery consumes CPU, allocations, memory, runtime tasks, bandwidth, serialization work, and rendering capacity.

When those costs become material, measure them explicitly and distinguish them from supervised work.

In particular, investigate:

- allocation from fine-grained evidence;
- retained-memory growth;
- observer/runtime-task multiplication;
- shared telemetry contention;
- copying or serialization of large payloads;
- high-cardinality metric growth;
- client/rendering backpressure;
- profiling or instrumentation perturbation.

Observability must obey `08-performance-engineering.md`.

### AREN-EVID-013 — Sensitive Payload Capture Is Explicit And Minimal

Complete observability does not authorize indiscriminate content capture.

When a capability can contain prompts, model output, tool arguments, files, environment data, credentials, HTTP data, or user content, explicitly decide:

- whether content is needed at all;
- whether metadata, lengths, hashes, classes, or identifiers are sufficient;
- what must be redacted or omitted;
- who may inspect the evidence;
- how long it may be retained.

Preserve structural execution evidence even when payload contents cannot safely be retained.

Never weaken semantic structure solely because sensitive content must be excluded.

### AREN-EVID-014 — Every Phase Closes An Observability Gate

Every roadmap phase must make the capability it introduces sufficiently inspectable.

During reasoning and review, ask where applicable:

1. Can we determine that the operation happened?
2. Can we determine when it happened?
3. Can we correlate it with its enclosing Aren operation?
4. Can we determine its outcome?
5. Can we inspect its failure without relying on an unstructured message?
6. Can we reconstruct important Aren-owned decisions?
7. Can we distinguish Aren overhead from supervised-work time where meaningful?
8. Can we determine materially relevant resource cost where meaningful?
9. Are high-volume observations explicitly bounded?
10. Is canonical truth distinguishable from diagnostics?
11. Can a slow or failed observer affect execution correctness?
12. Can the capability be demonstrated and investigated through the current supported observation surface?

A phase review must either close the applicable questions or explicitly record why a question is inapplicable or deliberately deferred without making a stronger observability claim.

## Relationship To Other Aren Contracts

- `03-execution-lifecycle.md` defines core lifecycle truth and terminal semantics.
- `05-events-observation-and-waiting.md` defines canonical Phase 1 lifecycle history, replay, cursors, and waiting.
- this contract defines the broader evidence model and explainability obligations around all Aren capabilities.
- `06-verification.md` proves that the supported evidence actually explains the tested path.
- `08-performance-engineering.md` governs observability overhead and performance claims.

Canonical lifecycle events are one class of Aren observable evidence. Logs, metrics, profiles, spans/intervals, state snapshots, and future high-volume execution telemetry must not be forced into the lifecycle-event schema merely for uniformity.

## Primary Governing Sources

- `projects/aren-phase-01-execution-lifecycle/docs/observability-mandate.md`
- `projects/aren-phase-01-execution-lifecycle/docs/PRD.md`
- `projects/aren-phase-01-execution-lifecycle/docs/performance-engineering.md`
- `projects/aren-phase-01-execution-lifecycle/contracts/03-execution-lifecycle.md`
- `projects/aren-phase-01-execution-lifecycle/contracts/05-events-observation-and-waiting.md`
