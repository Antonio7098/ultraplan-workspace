# Aren Verification Contract

> Project: `aren-phase-01-execution-lifecycle`  
> Status: active project contract; Phase 1 proof obligations are candidates for long-lived promotion where the semantics survive review

## Purpose

This contract defines how Aren correctness claims must be challenged.

Aren's hardest properties are semantic and concurrent. No single test style, race detector, test double, or green CI run can establish them. Verification must combine independent semantic expectations, controlled schedules, negative controls, race detection, release evidence, and real runtime execution.

## Scope

Apply this contract when a change affects:

- lifecycle state or transitions;
- outcome or failure interpretation;
- cancellation;
- goroutine/resource ownership;
- canonical events or observation;
- waiting or publication;
- public runtime or CLI behaviour;
- concurrency or synchronization.

## Requirement Index

| ID | Title | Severity If Violated |
| --- | --- | --- |
| AREN-VERIFY-001 | Correctness requires multiple proof classes | High |
| AREN-VERIFY-002 | Expected behaviour must be independent | Blocker |
| AREN-VERIFY-003 | Verify every transition pair | High |
| AREN-VERIFY-004 | Barriers beat sleeps | High |
| AREN-VERIFY-005 | Race detection is mandatory but insufficient | Blocker |
| AREN-VERIFY-006 | Test the test | High |
| AREN-VERIFY-007 | Release requires direct evidence | High |
| AREN-VERIFY-008 | Mutation attacks verify immutable publication | High |
| AREN-VERIFY-009 | Stress failures leave reproducible evidence | High |
| AREN-VERIFY-010 | Retry-to-green is not acceptance | High |
| AREN-VERIFY-011 | Real runtime paths matter | High |
| AREN-VERIFY-012 | Evidence is requirement-linked | High |
| AREN-VERIFY-013 | Measurements remain claims with boundaries | Medium |

## Requirements

### AREN-VERIFY-001 — Correctness Requires Multiple Proof Classes

Use the proof classes relevant to the changed behaviour.

For Phase 1 these include:

- an independently defined small semantic model or truth table;
- exact lifecycle trajectory assertions;
- real invocation and commit tests;
- controlled concurrency schedules;
- deliberately broken negative controls;
- complete applicable race-detector execution;
- Aren-owned resource release evidence;
- stress or generated-operation tests where useful;
- real CLI/runtime execution;
- review of the realized contract.

No one class substitutes for all the others.

### AREN-VERIFY-002 — Expected Behaviour Must Be Independent

Tests must not use the production transition implementation, terminal resolver, cursor implementation, or equivalent runtime decision code as the expected oracle for the same behaviour.

Expected legal edges, outcome rules, cancellation dispositions, and cursor rules should be stated independently enough to catch production defects rather than repeat them.

A test double that already implements the intended lifecycle is not proof of the real runtime path.

### AREN-VERIFY-003 — Verify Every Transition Pair

For the five Phase 1 lifecycle states, verification must cover:

- the four legal ordered state transitions;
- every other ordered state pair as forbidden;
- invalid/unknown vocabulary where the implementation boundary can receive it.

Cancellation-request recording is not a legal `running -> running` state transition.

The expected transition matrix must be independently authored.

### AREN-VERIFY-004 — Barriers Beat Sleeps

Use synchronization barriers, channels, latches, acknowledgements, or narrowly scoped production-path checkpoints to establish important concurrent orderings.

Do not use arbitrary sleeps as the primary proof that one operation happened before another.

Sleeps, polling, and watchdog timeouts may be used to:

- bound a failed test;
- observe eventual release;
- gather diagnostics;

but they do not create runtime semantics and do not prove the intended interleaving by themselves.

### AREN-VERIFY-005 — Race Detection Is Mandatory But Insufficient

The complete applicable Go test suite must run under `go test -race` in CI for Phase 1 acceptance.

A green race detector establishes only that no data race was detected on executed paths.

It does not prove:

- legal lifecycle histories;
- correct terminal arbitration;
- deadlock freedom;
- leak freedom;
- cancellation truthfulness;
- complete test selection;
- observer isolation;
- absence of logical races.

### AREN-VERIFY-006 — Test The Test

High-value verification instruments must demonstrate that they can fail for the intended reason.

Use one or more of:

- a deliberately broken fixture;
- a seeded historical defect;
- a mutation of the production path;
- a negative implementation variant.

Important Phase 1 regression seeds include:

- duplicate finalization;
- parent cancellation bypassing Aren acceptance;
- request/terminal publication inversion;
- mutable published outcomes or events;
- missing finish timing;
- flattened error causes;
- recursive invariant handling;
- lost observer wakeups;
- retained parent watchers.

### AREN-VERIFY-007 — Release Requires Direct Evidence

Do not infer absence of leaks merely because:

- the test process exits;
- a final global goroutine count is near its starting value;
- terminal outcome became ready.

Where Aren creates support resources, tests should observe their actual exit, stop, deregistration, or release when possible.

Supplement direct witnesses with stack, heap, retention, or churn evidence when useful.

Intentionally blocked work must be released on every test exit before final leak assertions.

### AREN-VERIFY-008 — Mutation Attacks Verify Immutable Publication

Published Aren-owned lifecycle data must be attacked through any mutable aliases exposed by the API.

Tests should attempt to mutate returned maps, slices, stack data, payload structures, snapshots, or equivalent reference-bearing data and then independently reread canonical facts.

Equality before mutation is not proof of isolation.

This requirement does not attempt to deep-freeze arbitrary user result graphs.

### AREN-VERIFY-009 — Stress Failures Leave Reproducible Evidence

Stress, generated, or randomized tests must record enough evidence to investigate failures.

Where applicable record:

- seed;
- iteration;
- operation trace;
- reached checkpoints/barriers;
- relevant input facts;
- observed history/outcome;
- toolchain and revision.

A random seed reproduces generated inputs. It does not claim to reproduce the Go scheduler.

Prefer reducing a failing trace to a smaller deterministic schedule when possible.

### AREN-VERIFY-010 — Retry-To-Green Is Not Acceptance

A flaky test that is rerun until green remains an unresolved correctness problem.

Acceptance evidence must not silently discard failed attempts or use automatic retry as a substitute for understanding nondeterminism.

### AREN-VERIFY-011 — Real Runtime Paths Matter

Critical lifecycle and CLI demonstrations must exercise the real runtime path.

For Phase 1 the final diagnostic surface must prove the supported scenarios through the built implementation rather than a display mock or fake lifecycle.

Required final Phase 1 scenarios include equivalents of:

```text
aren dev run success
aren dev run fail
aren dev run cancel
aren dev run race
```

The CLI output must not imply persistence, exactly-once delivery, or stronger cancellation guarantees than the runtime provides.

### AREN-VERIFY-012 — Evidence Is Requirement-Linked

Verification evidence should identify:

- the requirement or invariant being tested;
- the actual test/subtest or scenario;
- the command executed;
- revision/diff state;
- result;
- negative control where applicable;
- relevant limitation.

Do not substitute test counts for requirement coverage.

### AREN-VERIFY-013 — Measurements Remain Claims With Boundaries

When performance or resource behaviour is measured, record enough context to keep the claim honest, including relevant:

- revision and diff state;
- toolchain;
- workload;
- repetitions/spread;
- hardware where material;
- artifact/binary identity where material.

Measure before optimizing.

A benchmark or goroutine count must not be promoted into a production capacity guarantee without a representative workload and explicit budget.

## Minimum Phase 1 Gate

Before Phase 1 contract promotion, verification must include fresh equivalents of:

```bash
go build ./...
go vet ./...
go test ./...
go test -race ./...
```

plus the controlled semantic, negative-control, release/leak, stress, and real CLI evidence required by the realized implementation.

## Relationship To Generic Workspace Testing Rules

This contract is the Aren-specific verification layer. Generic workspace testing guidance remains useful where compatible, but it must not force speculative interfaces or fake-port architecture solely to satisfy a testing pattern.

Test seams should be created when required to prove real behaviour, not as ceremony.

## Primary Governing Sources

- `projects/aren-phase-01-execution-lifecycle/docs/PRD.md`
- `projects/aren-phase-01-execution-lifecycle/docs/final-language-decision.md`
- `projects/aren-phase-01-execution-lifecycle/project-reasoning/reasoning.md`
- `projects/aren-phase-01-execution-lifecycle/project-reasoning/areas/05-verification-and-go-correctness.md`
