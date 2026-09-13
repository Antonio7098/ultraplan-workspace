# Aren Development Doctrine Contract

> Project: `aren-phase-01-execution-lifecycle`  
> Applies beyond Phase 1: yes  
> Status: active project contract

## Purpose

This contract governs how Aren is allowed to grow.

Aren exists to become a deeply understood execution runtime, not to pre-build the final platform. The contract protects the development discipline established by the Aren lineage, phased roadmap, and accepted Go decision: solve one concrete problem at a time, prove behaviour through runnable vertical slices, and require complexity to earn its cost.

This contract governs development and architectural decision-making. More specific runtime semantics are owned by the other Aren contracts.

## Scope

Apply this contract when a sprint or change:

- introduces a new capability or phase;
- introduces an abstraction, interface, package boundary, daemon, store, event mechanism, workflow mechanism, provider boundary, or extension point;
- broadens the public API;
- pulls deferred roadmap work into the current phase;
- changes the evidence or verification needed before Aren broadens further.

## Requirement Index

| ID | Title | Severity If Violated |
| --- | --- | --- |
| AREN-DEV-001 | One primary uncertainty per phase | High |
| AREN-DEV-002 | Behaviour before abstraction | High |
| AREN-DEV-003 | Foundations must be vertical | High |
| AREN-DEV-004 | Complexity must be earned | High |
| AREN-DEV-005 | Scope exclusions are binding | Blocker |
| AREN-DEV-006 | Failure is part of the feature | High |
| AREN-DEV-007 | Reviews may simplify | Medium |
| AREN-DEV-008 | Real use before broadening | High |
| AREN-DEV-009 | Planning must terminate in implementation | High |
| AREN-DEV-010 | Future architecture remains revisable | Medium |

## Requirements

### AREN-DEV-001 — One Primary Uncertainty Per Phase

Each Aren phase must answer one primary technical question.

A phase may contain several implementation tasks, but those tasks must converge on the same uncertainty. Adjacent capabilities must not be pulled forward merely because they are likely to be useful later.

**Evidence**

- the phase or sprint states its primary uncertainty;
- planned work maps directly to that uncertainty;
- unrelated future capabilities remain deferred.

### AREN-DEV-002 — Behaviour Before Abstraction

Start with concrete behaviour.

An abstraction is earned only by demonstrated pressure such as:

- a second implementation with materially different behaviour;
- repeated conditional logic or duplicated behaviour;
- a real ownership boundary;
- a failure-testing seam that cannot reasonably be exercised otherwise;
- a concept that has remained stable through repeated use;
- a concrete requirement to replace or extend behaviour.

The possibility that an abstraction may be useful later is not sufficient evidence.

**Forbidden by default**

- interfaces whose only justification is future extensibility;
- generic factories around one stable implementation;
- plugin systems without plugin pressure;
- provider-neutral frameworks before provider variation exists;
- framework layers created to make the architecture look complete.

### AREN-DEV-003 — Foundations Must Be Vertical

A foundational capability must be exercised through a runnable end-to-end path.

Types, interfaces, package structure, and design documents alone do not constitute a completed foundation.

A vertical slice may be small, but it must expose enough real behaviour to reveal failure, ownership, lifecycle, and API mistakes.

### AREN-DEV-004 — Complexity Must Be Earned

Do not introduce operational or architectural machinery until a current requirement demonstrates why it is necessary.

This includes, but is not limited to:

- persistence;
- databases;
- daemons;
- event buses;
- workflow engines;
- plugin systems;
- generic executor abstractions;
- broad configuration frameworks;
- remote protocols;
- distributed execution;
- generalized middleware systems.

A dependency or abstraction is acceptable when it removes a demonstrated risk or substantial maintained code. It is not acceptable merely because mature systems commonly have one.

### AREN-DEV-005 — Scope Exclusions Are Binding

Explicit phase exclusions are constraints, not backlog suggestions.

Deferred functionality must remain absent rather than appearing as partially implemented infrastructure, speculative interfaces, placeholder packages, or unused extension points.

A deferred capability may enter scope only through an explicit decision that identifies the new evidence or requirement that justifies it.

### AREN-DEV-006 — Failure Is Part Of The Feature

Every new blocking, stateful, concurrent, external, or long-running capability must define and verify its relevant:

- failure behaviour;
- cancellation behaviour;
- ownership and release behaviour;
- concurrency behaviour;
- partial or uncertain outcome behaviour where applicable.

Happy-path implementation is not feature completion.

### AREN-DEV-007 — Reviews May Simplify

Reduction is valid progress.

A phase or sprint review may conclude that Aren should:

- delete an interface;
- collapse a package;
- remove a configuration option;
- merge or delete an event;
- narrow the public API;
- defer the next planned capability.

Review is not only a gate for adding more functionality.

### AREN-DEV-008 — Real Use Before Broadening

A capability must be runnable and exercised before the next conceptual layer is added.

The review should use real execution to ask whether:

- the API is awkward;
- ownership is unclear;
- events are missing or redundant;
- failure semantics are ambiguous;
- hidden state affects behaviour;
- an abstraction appeared before it was necessary.

### AREN-DEV-009 — Planning Must Terminate In Implementation

Research and planning exist to improve decisions and implementation.

Every substantial Aren planning artefact should eventually connect to at least one of:

- a decision;
- a contract;
- an implementation task;
- a test;
- a rejected alternative;
- a phase boundary.

Documents that no longer affect action should not be expanded indefinitely.

### AREN-DEV-010 — Future Architecture Remains Revisable

Later roadmap phases are hypotheses, not commitments.

Foundational semantics that have been implemented and promoted through review may be treated as contracts. Future package layouts, execution types, protocols, stores, deployment models, and extension mechanisms remain revisable until evidence stabilises them.

## Relationship To Other Contracts

This contract says **how Aren grows**.

Use:

- `02-runtime-architecture.md` for ownership and package/public-boundary rules;
- `03-execution-lifecycle.md` for lifecycle truth and terminal semantics;
- `04-cancellation-and-lifetimes.md` for cancellation and resource ownership;
- `05-events-observation-and-waiting.md` for canonical history and readers;
- `06-verification.md` for proof obligations.

## Primary Governing Sources

- `projects/aren-phase-01-execution-lifecycle/docs/project-lineage.md`
- `projects/aren-phase-01-execution-lifecycle/docs/phased-roadmap.md`
- `projects/aren-phase-01-execution-lifecycle/docs/final-language-decision.md`
