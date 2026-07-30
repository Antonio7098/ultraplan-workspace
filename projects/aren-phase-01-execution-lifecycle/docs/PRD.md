# Product Requirements: Aren Phase 1 — Execution Lifecycle

> Project: `aren-phase-01-execution-lifecycle`  
> Target repository: `../Aren/`  
> Implementation language: Go  
> Deployment model: local, in-process library with a diagnostic CLI

## 1. Summary

Phase 1 establishes the smallest execution lifecycle that Aren owns.

Aren must supervise one in-process occurrence of work from creation through exactly one terminal outcome. It must assign the run an identity, control every lifecycle transition, propagate cancellation, retain ordered lifecycle events, publish one immutable outcome, support concurrent waiting and observation, and remain correct under completion-cancellation races.

The atomic unit owned by Aren in this phase is the **lifecycle transition**. Aren does not make arbitrary work transactional. It guarantees that its own state, event, timing, cancellation metadata, and terminal-outcome bookkeeping become visible coherently.

Phase 1 deliberately excludes model providers, subprocesses, tools, persistence, retries, streaming output, workflows, pause/resume, remote execution, and daemon hosting.

## 2. Primary Question

> Can Aren define and enforce the lifecycle of one execution without depending on an LLM, subprocess, network call, or persistent store?

The phase is complete only when this can be answered with runnable implementation evidence, automated tests, race-detector results, and a stable written lifecycle contract.

## 3. Problem

Later execution types cannot be supervised coherently until Aren can answer:

- What is a run?
- Who owns run identity and lifecycle state?
- Which transitions are legal?
- What constitutes completion?
- How are result, failure, cancellation, and outcome separated?
- What happens when cancellation races with completion?
- When may Aren truthfully report a run as cancelled?
- How can multiple callers wait for one outcome?
- How can multiple observers read events without influencing execution?
- Which guarantees apply to Aren's bookkeeping, and which do not apply to work side effects?

Without explicit answers, later providers, tools, subprocesses, agents, and workflows would invent incompatible lifecycle behaviour.

## 4. Scope

### In scope

- one supervised run per invocation;
- Aren-generated opaque run identity;
- lifecycle states and transition guards;
- one in-process work function receiving `context.Context`;
- success, returned failure, and panic handling;
- immutable terminal outcomes;
- lifecycle timing;
- explicit and parent-context cancellation;
- cancellation request acceptance and disposition;
- centralized deterministic terminal resolution;
- completion-cancellation race handling;
- retained per-run lifecycle-event history;
- stable per-run event sequence numbers;
- multiple observers and replay from sequence;
- multiple concurrent waiters;
- race, stress, negative, panic, and leak tests;
- a diagnostic development CLI;
- a promoted lifecycle contract after validation.

### Out of scope

- model providers or messages;
- token or partial-output streaming;
- structured model output;
- retries, attempts, or repair loops;
- tools or tool calls;
- subprocesses;
- cleanup hooks;
- persistence or restart recovery;
- durable replay;
- pause/resume or human approval;
- global event streams;
- workflows or multiple execution types;
- remote execution or daemon hosting;
- network APIs or language SDKs;
- budgets, iteration caps, or stuck-loop detection;
- production telemetry exporters;
- exactly-once arbitrary work execution;
- rollback or compensation of work side effects;
- idempotent run creation;
- a permanent universal executor abstraction.

## 5. Definitions

### Execution

Aren-supervised work and the lifecycle surrounding that work.

### Run

One concrete occurrence of an execution, with its own identity, state, event history, timing, cancellation facts, and outcome.

### Work

The in-process operation supervised by Aren. Work receives a context and returns a result and error according to Go semantics. It may panic. It cannot directly mutate run state.

### State

The current lifecycle position of a run.

Phase 1 states are:

- `created`;
- `running`;
- `succeeded`;
- `failed`;
- `cancelled`.

### Atomic lifecycle transition

One Aren-owned operation that validates the source state, determines the destination, records transition metadata and timing, constructs any terminal outcome, allocates the next event sequence, appends the event, makes the complete transition visible, and releases terminal waiters where applicable.

### Outcome

The immutable final representation of a completed run, including identity, terminal state, timing, and the corresponding result, failure, or cancellation information.

### Cancellation request

A request asking work to stop cooperatively. A request is not proof that work stopped.

### Accepted cancellation

The first cancellation request received while the run is active. Aren records its reason, cancels the work context, and records at most one cancellation-request event.

### Canonical event history

The ordered in-memory sequence of lifecycle events recorded for one run. Event identity is `(run_id, sequence)`.

## 6. Product Principles

1. **Aren owns control flow.** Work may produce facts, but only Aren may transition lifecycle state.
2. **Lifecycle transition is the atomic unit.** Work side effects remain outside Aren's transactional guarantee.
3. **Exactly one terminal outcome.** A started run may commit no more than one terminal outcome and one terminal event.
4. **Terminal resolution is centralized.** Competing paths do not independently write terminal states.
5. **Cancellation is truthful.** Request, acceptance, context propagation, work acknowledgement, and terminal cancellation remain distinct.
6. **Events describe committed truth.** No event may announce a state or outcome the run did not enter.
7. **Recording is not delivery.** Each event is recorded once in canonical history; replay may expose it more than once to a consumer.
8. **Observation cannot control execution.** Slow, absent, or abandoned observers must not block the lifecycle.
9. **Authority remains separated.** Observation, caller control, and internal transition authority are different capabilities.
10. **The smallest honest semantics win.** Phase 1 must not create speculative provider, persistence, workflow, plugin, or executor infrastructure.

## 7. Lifecycle

```text
created
   ↓
running
   ├──→ succeeded
   ├──→ failed
   └──→ cancelled
```

The only legal state transitions are:

- `created -> running`;
- `running -> succeeded`;
- `running -> failed`;
- `running -> cancelled`.

Every terminal state is final. No terminal-to-terminal or terminal-to-running transition is legal.

`created` is canonical history but not an externally actionable scheduling state. A run may already be running or terminal when its handle is returned.

## 8. Run Identity

Every run must have one opaque Aren-generated identifier.

It must:

- exist before `run.created` is recorded;
- appear in every event and outcome;
- remain unchanged for the lifetime of the run;
- require no parsing or semantic interpretation;
- be safe for logs and diagnostic output.

Caller-supplied IDs and idempotency keys are out of scope.

## 9. Authority Boundaries

### Observation authority

May inspect identity and state, wait for completion, retrieve the outcome, and read or subscribe to lifecycle events.

### Caller control authority

May request cancellation. Phase 1 exposes no other caller command.

### Internal transition authority

Only Aren may validate and commit transitions, allocate sequence numbers, record events, publish outcomes, and release waiters.

The exact Go interface layout is a technical decision. The capability separation is a product requirement.

## 10. Work Invocation

The Phase 1 work shape must be equivalent in capability to:

```go
func(context.Context) (Result, error)
```

Aren must:

1. create the run and identity;
2. record `run.created`;
3. transition to `running` and record `run.started`;
4. derive a child context from the caller's parent context;
5. invoke work outside the lifecycle mutation critical section;
6. capture result, error, or panic;
7. resolve the terminal candidate through one central policy;
8. commit exactly one terminal transition and terminal event;
9. publish one immutable outcome;
10. release all waiters only after the complete terminal transition is visible.

The work-function shape is a semantic instrument, not a promised permanent executor API.

## 11. Outcome And Failure Requirements

### Successful outcome

Contains:

- state `succeeded`;
- the exact successful result;
- start and finish time.

It contains no failure or terminal cancellation information.

### Failed outcome

Contains:

- state `failed`;
- structured failure information;
- start and finish time.

It exposes no successful result.

### Cancelled outcome

Contains:

- state `cancelled`;
- cancellation metadata;
- start and finish time.

The first accepted cancellation reason is retained.

### Failure model

Failures must distinguish at least:

```text
origin: work | aren
kind: returned_error | panic | invariant_violation
```

Returned work errors preserve their original cause where practical. Work panic remains recognisable as panic. Internal invariant violation remains recognisable as Aren-origin failure.

If work returns both a result and an error, the invocation is not successful. Phase 1 defines no partial-result semantics.

## 12. Cancellation

Cancellation may originate from:

- explicit run-controller cancellation;
- cancellation of the supplied parent context.

Both sources use one Aren-owned acceptance path.

A cancellation request reports one of:

- `accepted`;
- `already_requested`;
- `already_terminal`.

The first accepted request:

- records the reason;
- cancels the work context;
- records one `run.cancellation_requested` event.

Repeated requests do not replace the reason or duplicate events.

Aren must not report terminal cancellation merely because cancellation was requested. The run remains `running` until work returns and terminal resolution commits an outcome.

If work ignores cancellation, Aren does not force-kill the goroutine and must not falsely claim that it stopped.

Starting with an already-cancelled parent context still follows `created -> running`; work receives a cancelled context and normal terminal resolution applies.

## 13. Terminal Resolution

All completion facts pass through one terminal-resolution function.

The Phase 1 policy is:

- result with nil error -> `succeeded`, even if cancellation was accepted shortly before return;
- accepted cancellation plus an error matching the run context cancellation cause -> `cancelled`;
- any other returned error -> `failed` with `origin=work`, `kind=returned_error`;
- panic -> `failed` with `origin=work`, `kind=panic`;
- already committed terminal outcome -> retain it unchanged and reject later mutation.

Scheduler timing may determine whether cancellation is accepted before completion. Given the same committed facts, interpretation must not depend on which goroutine attempts to write first.

## 14. Transition Atomicity

Every transition must occur through one Aren-owned critical section.

No supported concurrent operation may observe:

- terminal state without terminal event;
- terminal event without terminal outcome;
- outcome while state is nonterminal;
- a sequence increment without its event;
- two terminal outcomes;
- two terminal events;
- event transition data inconsistent with committed state.

Work execution itself is outside this atomic boundary.

## 15. Lifecycle Events

Phase 1 supports:

- `run.created`;
- `run.started`;
- `run.cancellation_requested`;
- `run.succeeded`;
- `run.failed`;
- `run.cancelled`.

Every event contains:

- run identity;
- type;
- per-run sequence number;
- occurrence time;
- source and destination state where applicable;
- immutable type-specific data where required.

`run.created` uses sequence `0`. Event order is defined by sequence, not timestamps.

`run.cancellation_requested` is an occurrence while the state remains `running`.

Exactly one terminal event is recorded in canonical history. Phase 1 does not claim exactly-once observer delivery.

## 16. Event Observation

Phase 1 uses retained, in-memory, per-run lifecycle history. It does not introduce a global event bus.

Requirements:

- multiple observers are supported;
- each observer has an independent cursor;
- observation can begin from a requested sequence;
- default observation begins from sequence `0`;
- observers created after completion can read the complete history;
- live observation completes only after the terminal event is visible;
- slow or abandoned observers do not block work, cancellation, transitions, waiters, or other observers;
- lifecycle history remains available while the run object remains reachable;
- no guarantee exists after process termination or collection;
- replay may return the same canonical event again;
- consumers can deduplicate with `(run_id, sequence)`.

## 17. Waiting And Concurrency

A caller can wait for the terminal outcome.

Waiting must:

- block while active;
- return only after complete terminal publication;
- return immediately after completion;
- provide the same logical outcome to every waiter;
- never consume the outcome exclusively.

Lifecycle mutation behaves as a single-writer system. Concurrent readers are supported.

The implementation must remain race-free under concurrent waiting, state inspection, explicit cancellation, parent cancellation, event subscription, event replay, terminal completion, and outcome inspection.

Aren guarantees immutability of the outcome container, not deep immutability of arbitrary result objects supplied by user code.

## 18. Timing

Aren records start and terminal time.

Requirements:

- start is not after terminal time;
- terminal time exists only after completion;
- all waiters see identical timing;
- ordering relies on event sequence, not timestamp comparison.

A clock abstraction is introduced only if deterministic testing proves difficult without one.

## 19. Guarantee Boundaries

Aren guarantees:

- one canonical lifecycle history per run;
- at most one terminal lifecycle commitment;
- one terminal outcome;
- coherent publication of state, event, timing, cancellation metadata, and outcome.

Aren does not guarantee:

- exactly-once work execution;
- transactional work effects;
- rollback or compensation;
- deduplication across separate runs;
- forceful interruption of a goroutine that ignores context;
- that failed or cancelled work produced no external effects.

## 20. Diagnostic CLI

Required scenarios:

```text
aren dev run success
aren dev run fail
aren dev run cancel
aren dev run race
```

Recommended scenarios:

```text
aren dev run panic
aren dev run parent-cancel
aren dev run ignore-cancel
```

The CLI must exercise the real runtime, display run identity and ordered events, distinguish cancellation request from terminal cancellation, show the terminal outcome, return useful process status, and avoid implying persistence or exactly-once delivery.

## 21. Required Verification

The phase must include automated evidence for:

- normal success;
- returned failure;
- panic;
- explicit cancellation;
- repeated cancellation;
- parent-context cancellation;
- already-cancelled parent context;
- delayed cancellation acknowledgement;
- ignored cancellation;
- success after accepted cancellation;
- cancellation-related error after accepted cancellation;
- unrelated error after accepted cancellation;
- high-iteration completion-cancellation races;
- multiple waiters;
- observer registration before, during, and after completion;
- replay from sequence;
- slow and abandoned observers;
- transition atomicity under concurrent reads;
- immutable lifecycle-event payloads;
- illegal transition detection;
- timing coherence;
- broad concurrency stress under `go test -race`;
- absence of Aren-owned observer/waiter goroutine leaks.

## 22. Acceptance Criteria

Phase 1 is accepted only when:

- all states, legal transitions, and illegal transitions are defined;
- lifecycle transition is established as Aren's atomic unit;
- Aren alone owns transition enforcement;
- terminal resolution is centralized and deterministic;
- exactly one terminal outcome and terminal event are guaranteed;
- cancellation request, acceptance, observation, and terminal cancellation remain distinct;
- explicit and parent cancellation use one path;
- state, event, timing, outcome, and waiter release are coherent;
- retained per-run history and sequence replay work;
- multiple observers and waiters are supported;
- absent or abandoned observers cannot deadlock execution;
- returned errors, panics, and Aren invariant failures remain distinguishable;
- lifecycle guarantees are explicitly separated from work-side-effect guarantees;
- all required tests pass under the Go race detector;
- diagnostic CLI scenarios provide runnable evidence;
- a written phase review finds no unresolved foundational semantic ambiguity.

## 23. Two-Sprint Delivery Shape

### Sprint 1 — Core Lifecycle

Establish the coherent non-cancellation lifecycle:

- identity;
- states and transition guards;
- work invocation;
- success, returned failure, and panic;
- terminal outcomes;
- waiting;
- transition atomicity;
- basic retained lifecycle history;
- contract and negative tests.

Sprint 1 must end in a complete runnable state. It must not defer ordinary correctness testing to Sprint 2.

### Sprint 2 — Cancellation And Concurrency Hardening

Extend and attack the Sprint 1 lifecycle through:

- cancellation acceptance and disposition;
- parent cancellation;
- deterministic terminal resolution;
- completion-cancellation races;
- multiple waiters and observers;
- replay cursors;
- slow and abandoned observers;
- race, stress, and leak testing;
- diagnostic CLI;
- lifecycle-contract promotion;
- phase review.

Sprint 2 must consume Sprint 1 reasoning and decisions as prior project context. It may supersede them only with explicit evidence and recorded rationale.

## 24. Exit Gate

The project is complete only when:

1. both sprints have completed their own acceptance criteria;
2. all lifecycle invariants are implemented and proven;
3. race and stress testing reveal no duplicate terminal outcome, inconsistent history, deadlock, or Aren-owned leak;
4. the diagnostic CLI demonstrates required scenarios;
5. the final lifecycle contract reflects the realised and tested semantics;
6. the phase review identifies no unresolved ambiguity that would undermine Phase 2;
7. later concerns remain deferred rather than hidden inside speculative abstractions.

## 25. Deferred Questions

- Should Aren eventually expose a general executor interface?
- How should progress and partial output be represented?
- How should cleanup affect terminal outcomes?
- How should attempts relate to runs?
- How are retry side effects made idempotent?
- Which event types are shared across model, tool, and subprocess execution?
- Which state must survive process termination?
- How are durable events, pause, resume, and approval represented?
- When does Aren require a daemon?
- Which Phase 1 semantics prove universal and which remain executor-specific?

These questions may inform later study selection but must not expand this project.
