# Aren Runtime Architecture Contract

> Project: `aren-phase-01-execution-lifecycle`  
> Status: active project contract

## Purpose

This contract defines the architectural constraints for Aren's runtime.

Aren's architecture is primarily organised around **authority, ownership, and observable behaviour**, not a predetermined application-layer template. Package structure and synchronization mechanisms may change, but the runtime must retain one clear owner for lifecycle truth and a deliberately small public surface.

## Scope

Apply this contract when a change:

- adds or changes runtime packages;
- changes lifecycle ownership or synchronization;
- changes exported runtime APIs;
- introduces an adapter, CLI entrypoint, transport, provider, store, or process boundary;
- introduces interfaces or replaceable collaborators;
- creates another representation of lifecycle state, history, timing, cancellation facts, or outcomes.

## Requirement Index

| ID | Title | Severity If Violated |
| --- | --- | --- |
| AREN-ARCH-001 | Authority determines boundaries | High |
| AREN-ARCH-002 | Lifecycle mutation has one owner | Blocker |
| AREN-ARCH-003 | Entrypoints stay thin | High |
| AREN-ARCH-004 | Concrete before interface | High |
| AREN-ARCH-005 | Public surface is deliberate | High |
| AREN-ARCH-006 | No hidden competing state | Blocker |
| AREN-ARCH-007 | Work executes outside mutation ownership | Blocker |
| AREN-ARCH-008 | Future architecture is not current architecture | High |
| AREN-ARCH-009 | Shared infrastructure must not own product semantics | High |
| AREN-ARCH-010 | Stable behaviour outranks package shape | Medium |

## Requirements

### AREN-ARCH-001 — Authority Determines Boundaries

Package and type boundaries should follow real ownership of behaviour and state.

Do not force Aren into a domain/use-case/adapter template merely for architectural symmetry. A boundary must protect a cohesive responsibility, ownership boundary, volatile dependency, or genuine test seam.

### AREN-ARCH-002 — Lifecycle Mutation Has One Owner

For each run, one private Aren authority owns mutation of:

- lifecycle state;
- legal transition enforcement;
- event sequence allocation;
- canonical lifecycle history;
- lifecycle timing;
- accepted cancellation facts;
- terminal outcome;
- completion readiness.

Several goroutines may invoke operations on that authority. They must not become independent lifecycle writers.

Work, callers, observers, transports, and adapters may provide facts or requests. They do not publish lifecycle truth.

The contract requires logical single-writer behaviour; it does not require a dedicated owner goroutine. A mutex-protected commit boundary, ownership loop, or another proved mechanism may satisfy it.

### AREN-ARCH-003 — Entrypoints Stay Thin

CLI handlers and later transport entrypoints are adapters.

They may:

- parse and validate boundary input;
- call the runtime through its supported public surface;
- map returned facts to output and process status.

They must not own lifecycle transitions, terminal resolution, cancellation interpretation, event recording, or runtime state mutation.

### AREN-ARCH-004 — Concrete Before Interface

Prefer concrete packages, functions, and types until current behaviour earns an abstraction.

An interface with one implementation is not automatically wrong, but it requires a specific demonstrated reason such as a real volatile boundary, independent implementation, or necessary failure-testing seam.

Do not introduce a permanent universal executor interface during Phase 1.

### AREN-ARCH-005 — Public Surface Is Deliberate

Export only capabilities callers genuinely require.

Public observation or cancellation capabilities must not expose internal transition authority.

Keep mutation helpers, transition constructors, sequence allocation, and internal outcome construction private unless a later requirement proves otherwise.

Go exports are compatibility cost. New exported identifiers require a concrete caller need.

### AREN-ARCH-006 — No Hidden Competing State

Do not maintain separately authoritative mirrors of lifecycle state, history, timing, accepted cancellation facts, outcome, or completion that later require reconciliation.

Derived views are allowed only when their relationship to canonical run facts is explicit and cannot silently become another authority.

A cache, transport projection, CLI model, or later persistence representation must not independently decide lifecycle truth.

### AREN-ARCH-007 — Work Executes Outside Mutation Ownership

The following must not execute while lifecycle mutation ownership is held:

- user work;
- user callbacks;
- user error or formatting methods that may execute arbitrary code;
- observer processing;
- blocking event delivery;
- resource joins;
- slow external operations.

The commit boundary should validate and install Aren-owned facts, not execute user-controlled work.

### AREN-ARCH-008 — Future Architecture Is Not Current Architecture

Do not shape current packages or APIs around deferred capabilities such as:

- providers;
- subprocesses;
- tools;
- retries;
- persistence;
- workflows;
- daemon hosting;
- remote execution;
- distributed workers.

Leave clear internals, not speculative extension points.

### AREN-ARCH-009 — Shared Infrastructure Must Not Own Product Semantics

Generic helpers may provide mechanical capabilities such as clocks, IDs, synchronization, serialization, or diagnostics when earned.

They must not become the authority for Aren lifecycle semantics merely because the implementation is reusable.

Product meaning remains in the runtime boundary that owns the run.

### AREN-ARCH-010 — Stable Behaviour Outranks Package Shape

Package structure is revisable.

A refactor may collapse, split, or rename packages if it preserves the governing behavioural contracts and leaves ownership clearer. Do not preserve a weak package boundary solely because it already exists.

## Architectural Consequences For Phase 1

The Phase 1 architecture should remain directly inspectable. It should be possible to identify, without reconstructing a framework:

1. where a run is created;
2. where legal transitions are checked;
3. where terminal facts are interpreted;
4. where complete lifecycle facts are committed;
5. where waiters/readers observe publication;
6. where Aren-owned support resources are released.

## Relationship To Generic Workspace Contracts

Where a generic workspace architecture rule conflicts with this project-specific contract, this contract is the Aren-specific narrowing layer for the project.

In particular, generic layering, registrar, facade, or port rules must not force speculative abstractions that violate `AREN-DEV-002` or `AREN-ARCH-004`.

## Primary Governing Sources

- `projects/aren-phase-01-execution-lifecycle/docs/PRD.md`
- `projects/aren-phase-01-execution-lifecycle/docs/final-language-decision.md`
- `projects/aren-phase-01-execution-lifecycle/project-reasoning/reasoning.md`
- `projects/aren-phase-01-execution-lifecycle/project-reasoning/areas/01-lifecycle-authority-and-atomic-publication.md`
