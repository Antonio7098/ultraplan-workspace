# Aren Observability Mandate

> **Status:** Project-wide engineering mandate  
> **Applies to:** Every roadmap phase and every Aren-owned runtime capability  
> **Relationship to roadmap:** Strengthens the roadmap requirement that foundations be complete, observable vertical slices. It does not authorize speculative observability infrastructure in early phases.

## 1. Purpose

Aren is intended to own and supervise autonomous execution. A runtime with that responsibility must be able to explain its own behaviour precisely.

> **Aren should never require guesswork to explain its own behaviour.**

Every meaningful action, decision, transition, wait, effect, failure, and materially relevant resource cost owned by Aren must be inspectable at the smallest useful semantic level, without observation changing execution behaviour.

Observability is not an optional telemetry layer added after the runtime works. It is a foundational product property and an engineering requirement for every phase.

This does **not** mean logging every implementation detail, retaining every byte forever, building a durable telemetry backend in Phase 1, or designing a universal observability abstraction before real variation exists. Aren's normal rules still apply: behaviour before abstraction, one primary uncertainty per phase, runnable vertical slices, failure testing, real use, and complexity earned through evidence.

## 2. Core Requirements

### 2.1 Aren-owned behaviour must be explainable

For a meaningful Aren-owned operation, it should be possible, where applicable, to determine:

- that it happened and when;
- which run or enclosing execution it belonged to;
- what caused or requested it;
- which decision Aren made and why, when the reason is part of Aren's semantics;
- what state existed before and after it;
- how long it took and how it terminated;
- what failure occurred;
- which parent and child operations were involved;
- what materially relevant resources it consumed;
- which facts are canonical and which are diagnostic.

The required granularity is the **smallest useful semantic level**, not every internal instruction or synchronization operation.

For example, lifecycle diagnostics may distinguish:

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

A mutex acquisition normally does not deserve a canonical runtime event. Mutex contention belongs in profiling or diagnostic telemetry when it matters.

### 2.2 Observation must not control execution

Observability is passive with respect to runtime correctness.

A slow, absent, abandoned, overloaded, or faulty observer must not determine whether execution can progress, accept cancellation, commit a transition, terminate, release waiters, or clean up owned resources.

### 2.3 The UI is a projection, not a source of truth

The dedicated Aren Observability UI must project evidence already owned by Aren. It must not become a second lifecycle model, infer canonical outcomes from log strings, or require runtime behaviour merely to satisfy a screen.

```text
Aren runtime truth
      │
      ├── canonical facts
      ├── diagnostics
      ├── timings
      ├── metrics
      └── profiles
              │
              ▼
       observation surfaces
              │
       ┌──────┼────────┐
       │      │        │
      CLI     UI    exporters
```

The runtime evidence model comes first. Observation surfaces consume it.

## 3. Observation Forms

Aren must distinguish different kinds of observable information rather than forcing all evidence into one universal structure.

### Canonical semantic evidence

Canonical evidence records execution truth Aren itself owns: lifecycle transitions, accepted cancellation, terminal outcomes, attempt boundaries, permission decisions, tool outcomes, model invocation outcomes, workflow decisions, and later recovery decisions.

Canonical evidence must agree with runtime state. It is not merely a log message.

### Diagnostic telemetry

Diagnostic telemetry explains internal runtime behaviour useful for development and investigation but is not itself canonical truth. Examples include queue depth, observer lag, buffer occupancy, runtime-task count, heap use, garbage collection, scheduler behaviour, and open process or connection counts.

It may be sampled, aggregated, bounded, or environment-specific.

### Structured diagnostic logs

Logs provide contextual diagnostic information that does not warrant canonical semantic representation. They should be structured and correlated with stable Aren identities where applicable.

Semantic truth must not exist only as prose logs. Determining whether a run succeeded must never require parsing `INFO run completed successfully`; run state, canonical history, and outcome own that truth.

### Performance evidence

Performance evidence explains Aren's cost and scaling behaviour. Where meaningful, distinguish Aren overhead from supervised-work time, local processing from external waiting, allocation and retained memory, runtime-task growth, contention, queueing, throughput, tail latency, saturation, and recovery after load falls.

This evidence accumulates phase by phase as defined by `performance-engineering.md`.

### Profiles

Go CPU, heap, goroutine, mutex, block, and runtime traces provide deeper implementation-level evidence when semantic observation and ordinary metrics are insufficient. Profiles complement rather than replace semantic Aren evidence.

### State snapshots

Some questions concern what is true now rather than what happened historically: active runs, current state, accepted cancellation metadata, active child executions, budgets, reservations, or observer positions. Snapshots must remain consistent with canonical state rather than becoming another mutable truth model.

## 4. Retention Is Separate From Observability

"Everything observable" does not mean "everything retained forever."

Aren should distinguish three conceptual retention classes.

### Canonical

Evidence required to preserve or reconstruct Aren's semantic truth for the lifetime promised by the capability.

### Diagnostic

Detailed evidence useful for debugging or explanation but not required to establish canonical execution truth. Examples may include bounded stream deltas, subprocess output, adapter diagnostics, or detailed timing samples.

### Ephemeral

High-volume or continuously sampled information that may be aggregated, coalesced, sampled, or discarded, such as queue-depth samples, goroutine-count samples, observer lag, or profiling samples.

Retention must be chosen according to semantic importance, volume, recovery requirements, and demonstrated use. High-volume telemetry must not inherit Phase 1's complete lifecycle-history retention merely because both are observable.

## 5. Do Not Create a Universal Observation Type Prematurely

Different evidence forms answer different questions:

| Form | Primary question |
|---|---|
| Event | What happened? |
| Span or interval | Where was time spent? |
| Metric | How much or how many? |
| Log | What diagnostic context was reported? |
| Profile | Where are runtime resources being consumed? |
| State snapshot | What is true now? |

They should share correlation identities where useful, but they do not need one schema, storage policy, delivery mechanism, or retention rule.

A general abstraction should appear only after real implementations demonstrate stable common behaviour.

## 6. Correlation and Causality

As Aren gains nested execution, observable evidence must preserve enough relationship information to reconstruct execution structure.

Likely later identities include:

```text
run_id
attempt_id
model_invocation_id
tool_call_id
child_run_id
workflow_run_id
workflow_step_id
```

The exact vocabulary must be earned by the relevant phases.

Nested execution must not collapse into opaque parent operations. If a later programmatic execution invokes 200 Aren tools, the UI must not expose only `program succeeded` while hiding the governed nested calls. Those calls retain their normal validation, permission, lifecycle, timing, cancellation, failure, and evidence semantics.

Likewise, an agent run should eventually be explorable through its turns, model invocations, tool calls, retries, context changes, and child executions without each component inventing an incompatible tracing model.

## 7. Decision Observability

Aren should expose important Aren-owned decisions when they materially affect execution semantics.

Examples may eventually include:

- why a terminal candidate resolved to success, failure, or cancellation;
- why an error was considered retryable;
- why a retry delay was selected;
- why a tool request was denied;
- why admission was rejected;
- why a resource reservation was granted or refused;
- why context was compacted;
- why a workflow route was selected;
- why recovery resumed, replayed, interrupted, or refused uncertain work.

Prefer structured decision inputs, selected policy, and resulting disposition over arbitrary prose explanations generated after the fact.

## 8. Dedicated Aren Observability UI

Aren will have a dedicated observability UI.

Its long-term purpose is broader than production monitoring. It is a **runtime microscope** for development, testing, demonstration, performance investigation, debugging, and eventually operation.

The UI must evolve from proven runtime evidence rather than define that evidence prematurely.

### Run explorer

A run should eventually be inspectable as one coherent object showing, where available:

- identity and state;
- result or failure;
- timing and cancellation facts;
- canonical event history;
- logs and diagnostics;
- child operations and attempts;
- model invocations and tool calls;
- context operations;
- resource use and performance breakdowns.

### Unified timeline

Different observation forms should be alignable on a common time axis where useful:

```text
0.000 ms   run.created
0.018 ms   run.started
48.204 ms  cancellation accepted
48.219 ms  cancellation signal propagated
142.991 ms work returned success
143.014 ms terminal resolution -> succeeded
143.027 ms run.succeeded
143.041 ms waiters released
```

Canonical facts and diagnostic instrumentation should be visually distinguishable.

### Execution hierarchy

As composition appears, the UI should support drilling through execution structure rather than presenting only a flat log stream.

```text
Run
└── agent execution
    ├── turn
    │   └── model invocation
    ├── tool call
    │   └── subprocess
    └── turn
        └── model invocation
```

The hierarchy must follow Aren's eventual semantics rather than being fixed in advance.

### Performance inspection

The UI should make Aren overhead distinguishable from external work, while displaying only measurements Aren can support honestly. It must not imply precision or causal attribution that was not actually measured.

### Failure laboratory

The UI should eventually help inspect deliberately adversarial scenarios including cancellation/completion races, ignored cancellation, retry exhaustion, stream interruption, cleanup failure, observer backpressure, policy denial, resource saturation, crash and recovery, and workflow partial failure.

Observability is part of Aren's validation machinery, not merely an operator convenience.

## 9. Incremental Delivery

The mandate applies from the first runtime phase. The full observability product does not.

### Phase 0

Establish only the engineering foundation needed to collect trustworthy evidence. Do not build an observability backend, server, database, exporter framework, or dashboard platform.

### Phase 1 — Execution lifecycle

The first observable vertical slice consists of evidence already required by the lifecycle contract:

- run identity;
- canonical lifecycle history;
- deterministic sequence order;
- lifecycle timing;
- cancellation acceptance and disposition;
- structured failures and invariant diagnostics;
- immutable terminal outcome;
- diagnostic CLI scenarios;
- Go race and profiling evidence where useful.

The diagnostic CLI is the first Aren observation surface. Phase 1 then adds a dedicated browser frontend over the realised lifecycle and performance evidence. It uses a versioned read-only contract and a process-scoped loopback service. The service owns no lifecycle state and ends with its development process.

Phase 1 does not require production telemetry export, a global event stream, durable telemetry, authentication, multi-user operation, remote access, or persistent daemon hosting.

### Phase 2 — Controlled execution

New behaviour adds corresponding observation only for capabilities actually introduced, such as progress, partial output, cleanup, or richer timing.

### Phase 3 — Model invocation

Make provider interaction inspectable without confusing provider time with Aren overhead. Observe request lifecycle, provider failure classification, cancellation, timings, and provider usage metadata actually supplied by the provider.

### Phase 4 — Streaming

Streaming introduces high-volume observation pressure. Explicitly establish bounded behaviour for deltas, allocation, buffering, throughput, cancellation latency, and observer impact. This phase tests the distinction between canonical truth and bounded high-volume telemetry.

### Later model, tool, and agent phases

Each capability adds its own semantic evidence as it is introduced: structured validation, attempts and retry decisions, tool requests and permission decisions, tool execution, agent turns, termination and budget decisions, and context measurements and transformations.

The observability model grows from these concrete cases rather than from a speculative universal schema.

### Persistence and daemon phases

Persistence forces a deliberate decision about which observation state survives process termination.

Daemon hosting is the boundary for making the Phase 1 browser frontend shared, reconnectable across process lifetimes, remotely accessible, or operationally persistent. The Phase 1 process-scoped service must not grow those responsibilities solely to host observability.

### Workflows and broader execution

Workflow observation should compose existing evidence rather than replace it. A workflow should be inspectable as both a run and a hierarchy of governed child executions, decisions, data relationships, and failures.

## 10. Phase-Level Observability Gate

Every roadmap phase must make the capability it introduces inspectable.

During design and review, ask where applicable:

1. Can we determine that the operation happened?
2. Can we determine when it happened?
3. Can we correlate it with the enclosing run or parent operation?
4. Can we determine its outcome?
5. Can we inspect its failure without relying on an unstructured message?
6. Can we reconstruct important Aren-owned decisions?
7. Can we distinguish Aren overhead from supervised-work time where meaningful?
8. Can we determine materially relevant resource cost where meaningful?
9. Are high-volume observations explicitly bounded?
10. Is canonical truth distinguishable from diagnostics?
11. Can a slow or failed observer affect execution correctness?
12. Can the capability be demonstrated and investigated through the current Aren observation surface?

A capability that is functionally correct but materially opaque is incomplete unless the phase explicitly documents why the missing visibility is impossible, misleading, or outside Aren's ownership.

## 11. Observability and Failure Testing

Every failure test should ask not only whether Aren produced the correct outcome but whether the available evidence sufficiently explains the path that produced it.

A failure that can only be understood by adding temporary print statements after the test fails indicates an observability gap worth evaluating.

That does not mean every internal detail becomes permanent telemetry. The right response may be a canonical event, structured diagnostic, profile, state snapshot, or improved failure context depending on the missing evidence's semantic importance.

## 12. Observability Has a Cost

Observation mechanisms themselves must obey Aren's performance principles.

Risks include:

- allocation from fine-grained events;
- memory retention from histories;
- runtime-task multiplication from observers;
- contention in shared telemetry structures;
- copying large payloads;
- high-cardinality metric growth;
- serialization overhead;
- UI or exporter backpressure leaking into runtime execution.

Measure these costs when they become material. Do not weaken execution semantics for cheaper telemetry, but do not allow the mandate to become justification for unbounded instrumentation either.

## 13. Security and Sensitive Data

Complete observability does not mean indiscriminate payload capture.

Later phases involving prompts, model output, tool arguments, files, environment data, credentials, HTTP requests, or user content must explicitly decide:

- what may be recorded;
- what must be redacted;
- what should be represented by metadata or hashes instead of content;
- who may inspect it;
- how long it is retained.

Aren should preserve structural execution evidence even when payload contents cannot safely be retained.

## 14. Project-Wide Rule

The enduring rule is:

> **Every Aren capability must carry enough structured evidence to explain what Aren did, why it did it when the reason is part of Aren's semantics, when it happened, how it terminated, and what it cost where materially relevant. Observation must remain passive, canonical truth must remain distinct from diagnostics, and retention must be proportional to semantic importance and volume.**

The dedicated observability UI will make that evidence navigable and visual. It will not manufacture the evidence or become the runtime's source of truth.

Opacity in an Aren-owned capability should be treated as an engineering defect, not as something to solve after the runtime is complete.
