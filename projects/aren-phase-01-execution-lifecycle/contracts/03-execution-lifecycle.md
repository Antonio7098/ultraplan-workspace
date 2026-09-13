# Aren Execution Lifecycle Contract

> Project: `aren-phase-01-execution-lifecycle`  
> Status: governing Phase 1 contract; candidate for long-lived promotion after Phase 1 review

## Purpose

This contract defines the authoritative semantics of one Aren-supervised in-process execution.

It governs identity, lifecycle state, transition authority, terminal outcomes, failure interpretation, timing, and coherent publication. It deliberately does not promise transactional work effects, retries, persistence, workflow behaviour, provider semantics, or a permanent executor abstraction.

## Requirement Index

| ID | Title | Severity If Violated |
| --- | --- | --- |
| AREN-LIFE-001 | Run identity is Aren-owned and stable | High |
| AREN-LIFE-002 | Lifecycle vocabulary and legal transitions are closed | Blocker |
| AREN-LIFE-003 | Lifecycle transition is the atomic unit | Blocker |
| AREN-LIFE-004 | One terminal commitment | Blocker |
| AREN-LIFE-005 | Outcome is discriminated by terminal state | Blocker |
| AREN-LIFE-006 | Result plus error is not partial success | High |
| AREN-LIFE-007 | Terminal resolution is centralized | Blocker |
| AREN-LIFE-008 | Panic and runtime faults stay distinct | High |
| AREN-LIFE-009 | Publication is coherent | Blocker |
| AREN-LIFE-010 | Time does not define event order | High |
| AREN-LIFE-011 | Lifecycle immutability has an explicit ownership boundary | High |
| AREN-LIFE-012 | Lifecycle truth is not task truth | High |
| AREN-LIFE-013 | Illegal mutation remains visible | Blocker |

## Requirements

### AREN-LIFE-001 — Run Identity Is Aren-Owned And Stable

Every run has one opaque Aren-generated identifier.

It must:

- exist before `run.created` is recorded;
- remain unchanged for the run lifetime;
- appear in every canonical lifecycle event and terminal outcome;
- carry no required semantic interpretation;
- require no caller-supplied idempotency key or caller-generated identity in Phase 1.

### AREN-LIFE-002 — Lifecycle Vocabulary And Legal Transitions Are Closed

The Phase 1 states are exactly:

```text
created
running
succeeded
failed
cancelled
```

The only legal state transitions are:

```text
created -> running
running -> succeeded
running -> failed
running -> cancelled
```

Terminal states are final.

Cancellation acceptance is an occurrence while `running`; it is not an additional lifecycle state.

### AREN-LIFE-003 — Lifecycle Transition Is The Atomic Unit

Aren's atomic guarantee applies to its own lifecycle bookkeeping.

A transition coherently determines and publishes the Aren-owned facts that belong to that commitment, including the applicable state, event, sequence, timing, cancellation metadata, and terminal outcome.

Aren does **not** guarantee:

- exactly-once work execution;
- transactional user effects;
- rollback;
- compensation;
- absence of external effects after failure or cancellation.

### AREN-LIFE-004 — One Terminal Commitment

A completed run has exactly one terminal outcome and exactly one matching terminal lifecycle event.

Once terminal commitment is authoritative:

- no later path may replace the terminal state;
- no later path may replace the terminal outcome;
- no later terminal event may be appended;
- duplicate internal finalization attempts must be rejected or exposed as a defect.

A run whose work never returns is not required to become terminal.

### AREN-LIFE-005 — Outcome Is Discriminated By Terminal State

A successful outcome contains:

- state `succeeded`;
- the exact successful result;
- run identity;
- start and finish timing;
- no failure;
- no terminal cancellation payload.

A failed outcome contains:

- state `failed`;
- structured failure information;
- run identity;
- start and finish timing;
- no successful result;
- no terminal cancellation branch.

A cancelled outcome contains:

- state `cancelled`;
- the first accepted cancellation metadata and acknowledgement diagnostics as defined by the cancellation contract;
- run identity;
- start and finish timing;
- no successful result;
- no ordinary failure branch.

### AREN-LIFE-006 — Result Plus Error Is Not Partial Success

A completed work return with a non-nil Go error interface is not successful, even if the result value is non-nil.

Phase 1 defines no partial-result semantics.

A nil or zero result with a nil error is still a successful return.

### AREN-LIFE-007 — Terminal Resolution Is Centralized

All completed invocation facts pass through one terminal-resolution policy.

The Phase 1 policy is:

```text
existing terminal commitment
    -> retain it unchanged and reject later mutation

normal return with nil error
    -> succeeded

normal return with non-nil error
AND cancellation was accepted
AND errors.Is(returnedError, acceptedEffectiveCause)
    -> cancelled

any other non-nil returned error
    -> failed / origin=work / kind=returned_error

panic escaping the supervised invocation
    -> failed / origin=work / kind=panic
```

Accepted cancellation does not override a completed nil-error return.

A panic remains panic failure even if its value resembles the accepted cancellation cause.

### AREN-LIFE-008 — Panic And Runtime Faults Stay Distinct

A panic escaping the supplied work invocation, including its defers, is recognisable as `work/panic` and should retain useful captured diagnostics.

An Aren invariant violation is Aren-origin failure or escalation. It must not be relabelled as an ordinary work failure.

Broad recovery around the entire runtime is forbidden if it can disguise an Aren defect as user-work panic.

### AREN-LIFE-009 — Publication Is Coherent

No supported concurrent observation may expose:

- terminal state without its terminal outcome;
- terminal state without its terminal event;
- terminal event without its terminal outcome;
- an outcome while state is nonterminal;
- a sequence increment without its event;
- finish timing while the run is active;
- two terminal outcomes;
- two terminal events;
- event transition data inconsistent with committed state.

Waiter readiness must not become observable before complete terminal publication.

Separate accessor calls are not themselves a transaction and may observe different valid committed prefixes. Any API that promises a coherent snapshot must capture related facts together.

### AREN-LIFE-010 — Time Does Not Define Event Order

Aren records start and terminal timing once at the relevant lifecycle boundaries.

Requirements:

- start is not after finish;
- finish is absent while active;
- all waiters see identical committed timing;
- event ordering is defined by per-run sequence, not timestamp comparison;
- equal timestamps are valid.

Introduce a clock abstraction only when deterministic testing or a later runtime requirement demonstrates the need.

### AREN-LIFE-011 — Lifecycle Immutability Has An Explicit Ownership Boundary

Aren-owned lifecycle containers and payloads are immutable by contract after publication.

Supported APIs must not expose mutable aliases that can rewrite canonical state, event history, cancellation facts, failure classification, or timing.

Aren does not promise deep immutability of arbitrary successful result graphs supplied by user code. It preserves the exact result value according to the public API contract.

Original Go errors and custom causes may remain inspectable local diagnostic objects; Aren owns stable classification and diagnostic projections, not every mutable object reachable through user errors.

### AREN-LIFE-012 — Lifecycle Truth Is Not Task Truth

A successful Phase 1 invocation means the supervised work returned with a nil error and Aren committed a successful lifecycle outcome.

It does not prove that a future model, tool, business operation, or user-level task achieved its higher-level objective.

Execution success, model/provider success, and task-quality success must remain distinct concepts as later capabilities are introduced.

### AREN-LIFE-013 — Illegal Mutation Remains Visible

Illegal lifecycle mutation must be rejected before corrupting canonical truth.

Invariant handling must not:

- fabricate an otherwise illegal transition;
- overwrite an already terminal outcome;
- append a second terminal event;
- claim active work has stopped when it has not;
- recursively re-enter the transition path that exposed the invariant failure.

Where a runtime defect cannot be truthfully contained inside the current lifecycle, fail or escalate explicitly rather than manufacture a valid-looking history.

## Phase Boundary

This contract does not introduce:

- retries;
- cleanup hooks that can rewrite outcomes;
- partial output;
- provider or tool failure taxonomies;
- persistence or restart recovery;
- workflow composition;
- exactly-once execution;
- forceful goroutine interruption.

## Primary Governing Sources

- `projects/aren-phase-01-execution-lifecycle/docs/PRD.md`
- `projects/aren-phase-01-execution-lifecycle/project-reasoning/reasoning.md`
- `projects/aren-phase-01-execution-lifecycle/project-reasoning/areas/01-lifecycle-authority-and-atomic-publication.md`
- `projects/aren-phase-01-execution-lifecycle/project-reasoning/areas/02-outcomes-failures-and-terminal-resolution.md`
