# Aren Phase 1 Roadmap — Execution Lifecycle

> Project: `aren-phase-01-execution-lifecycle`  
> Target repository: `../Aren/`  
> Scope: define and prove the lifecycle of one supervised in-process execution.

## 1. Scope Principle

This project covers only Aren Phase 1.

Its purpose is to prove that Aren can own one execution lifecycle without relying on an LLM, subprocess, network call, persistent store, workflow engine, or daemon.

The project will not design later Aren capabilities in advance. Provider integration, tools, retries, persistence, workflows, pause/resume, remote execution, daemon hosting, and a universal executor abstraction remain deferred.

The implementation is divided into two sprints because the lifecycle has two distinct uncertainties:

1. Can Aren establish one coherent lifecycle and terminal outcome at all?
2. Does that lifecycle remain truthful and race-free under cancellation, observation, and adversarial concurrency?

The split is not by topic. Each sprint must end with a coherent, runnable, tested state.

---

## 2. Canonical Project Flow

```text
Phase 1 PRD
    ↓
Sprint 1: Core Lifecycle
    ↓
Sprint 1 reasoning + implementation + review
    ↓
Sprint 2 consumes Sprint 1 as prior project context
    ↓
Sprint 2: Cancellation And Concurrency Hardening
    ↓
Phase review + promoted lifecycle contract
```

Each sprint uses the normal UltraPlan chain:

```text
requirements -> sprint index -> technical handbook -> area reasoning -> sprint reasoning -> plan -> execute -> review
```

The technical handbook distills selected studies. It does not decide architecture or implementation. Area reasoning investigates the sprint from selected perspectives. Top-level `reasoning.md` makes final sprint decisions. `plan.md` executes those decisions and must not invent new architecture.

---

## 3. Cross-Sprint Carry-Forward Rule

Sprint 2 must not rediscover Sprint 1 from scratch.

Its planning inputs must explicitly include the completed Sprint 1 artifacts:

```text
sprints/01-core-lifecycle/requirements.md
sprints/01-core-lifecycle/sprint-index.md
sprints/01-core-lifecycle/technical-handbook.md
sprints/01-core-lifecycle/reasoning/*.md
sprints/01-core-lifecycle/reasoning.md
sprints/01-core-lifecycle/plan.md
sprints/01-core-lifecycle/execute.md          # when present
sprints/01-core-lifecycle/review.md           # when present
```

The Sprint 2 `sprint-index.md` must classify them as prior project decisions and realised implementation context.

Rules:

- Sprint 1 `reasoning.md` is the authoritative synthesis of Sprint 1 decisions.
- Sprint 1 area reasoning remains available for detailed rationale and rejected alternatives.
- Sprint 1 technical-handbook evidence may be reused, but Sprint 2 should select additional reports targeted to cancellation, event delivery, and concurrency rather than regenerate the same handbook blindly.
- Sprint 2 may supersede a Sprint 1 decision only when new evidence or realised implementation behaviour proves it insufficient or incorrect.
- A superseding decision must name the prior decision, explain the new evidence, describe the impact, and update the implementation plan accordingly.
- Silence does not supersede a prior decision.

This carry-forward rule should be copied into the Sprint 2 requirements and sprint index when that sprint is initialized.

---

## Implementation Wave 1 — Coherent Core Lifecycle

### Sprint 1: Core Lifecycle

> Slug: 01-core-lifecycle
> Status: planned
> Depends On:

#### Goal

Establish a complete, understandable, non-cancellation execution lifecycle for one in-process work function.

At sprint completion, Aren must be able to create a run, invoke work once, resolve success, returned failure, or panic, publish one coherent terminal outcome, retain basic lifecycle history, and allow multiple callers to wait safely.

#### Uncertainty

> Can Aren represent and enforce one lifecycle correctly before cancellation and adversarial concurrency are added?

#### Build

- opaque Aren-generated run identity;
- lifecycle states:
  - `created`;
  - `running`;
  - `succeeded`;
  - `failed`;
  - `cancelled` as vocabulary, without implementing caller cancellation yet unless needed internally;
- legal and illegal transition definitions;
- Aren-only transition authority;
- one in-process work function receiving `context.Context`;
- normal success;
- returned work failure;
- panic recovery and classification;
- immutable terminal outcome;
- lifecycle start and finish timing;
- multiple waiters;
- one central transition-commit boundary;
- basic retained per-run lifecycle history;
- event sequence identity `(run_id, sequence)`;
- contract, negative, panic, waiter, and atomicity tests;
- initial package architecture and dependency direction.

#### Deferred

- caller cancellation API;
- cancellation disposition;
- parent-context cancellation integration beyond the minimum needed to invoke work safely;
- deterministic cancellation-versus-completion resolution;
- cancellation-request events;
- multiple live observers and replay cursors if basic history snapshots are sufficient for Sprint 1;
- slow and abandoned observer behaviour;
- high-iteration cancellation races;
- diagnostic race CLI;
- phase-wide lifecycle-contract promotion.

Deferred items must not be predesigned through speculative interfaces. Sprint 1 should leave clear implementation seams only where its own requirements earn them.

#### Notes

The Sprint 1 index should select only the reasoning documents needed to resolve current uncertainty. Likely candidates are:

- `reasoning/architecture.md` — package ownership, public/internal boundaries, run aggregate placement, dependency direction;
- `reasoning/lifecycle-and-state.md` — states, legal transitions, atomic unit, outcome coherence, event facts;
- `reasoning/outcomes-and-failures.md` — result/error rules, panic, Aren invariant failures, timing;
- `reasoning/concurrency.md` — single-writer mutation, waiters, transition publication, lock scope;
- `reasoning/events.md` — minimum retained history and sequence semantics needed in Sprint 1;
- `reasoning/testing.md` — contract tables, illegal transitions, panic, waiters, atomicity, race detector.

The sprint index may merge or omit areas when the distinction does not justify a separate reasoning document.

#### Evidence

Strong candidates from the agent-harness study:

- `01.01-execution-model-taxonomy.md`;
- `01.02-control-flow-ownership.md`;
- `01.03-step-turn-task-atomicity.md`;
- `02.01-state-taxonomy-and-ownership.md`;
- `02.04-mutation-discipline-and-state-transitions.md`;
- `03.09-completion-and-finalization-semantics.md`.

Supporting Go reports may include project structure, state/context, concurrency, testing, and philosophy.

#### Deliverables

Normal UltraPlan sprint artifacts:

```text
requirements.md
sprint-index.md
technical-handbook.md
reasoning/*.md
reasoning.md
plan.md
flow-state.json
```

Implementation and evidence should produce, at minimum:

- a buildable Aren Go module or the required Phase 1 package foundation;
- the core run lifecycle implementation;
- state and transition vocabulary;
- outcome and failure types;
- panic boundary;
- retained basic history;
- waiting API;
- comprehensive Sprint 1 tests;
- review evidence.

#### Acceptance

Sprint 1 is complete only when:

- `created -> running -> succeeded` is proven;
- `created -> running -> failed` is proven for returned error and panic;
- exactly one terminal outcome and terminal event are recorded;
- all illegal transitions are rejected or exposed as Aren invariant violations;
- state, event, timing, outcome, and waiter release cannot be observed partially;
- multiple waiters observe the same outcome;
- tests pass under `go test -race`;
- no observer or cancellation behaviour is falsely claimed as implemented;
- the implementation remains small enough to inspect directly;
- Sprint 1 review records all realised API and architecture decisions needed by Sprint 2.

#### Evidence

Expected commands will be finalized by Sprint 1 reasoning, but should include equivalents of:

```bash
go test ./...
go test -race ./...
go build ./...
```

Any diagnostic command introduced in this sprint must exercise the real lifecycle rather than a mock display path.

---

## Implementation Wave 2 — Cancellation And Concurrency Hardening

### Sprint 2: Cancellation And Concurrency

> Slug: 02-cancellation-and-concurrency
> Status: planned
> Depends On: 1

#### Goal

Extend the realised Sprint 1 lifecycle with truthful cooperative cancellation, deterministic terminal resolution, multiple event observers, replay, and adversarial concurrency proof.

At sprint completion, Aren must remain coherent under cancellation-completion races, concurrent waiting, event subscription, parent-context cancellation, slow or abandoned observers, and repeated stress under the Go race detector.

#### Uncertainty

> Does the Sprint 1 lifecycle remain truthful, deterministic, and leak-resistant when cancellation and concurrent observation are introduced?

#### Notes

Sprint 2 must include all completed Sprint 1 planning and review artifacts listed in the cross-sprint carry-forward rule.

Its top-level reasoning must contain a section named or equivalent to:

```text
Prior Sprint Decisions Applied
```

That section must state for each relevant Sprint 1 decision whether Sprint 2:

- preserves it;
- extends it;
- supersedes it;
- leaves it unaffected.

#### Build

- explicit run-controller cancellation;
- parent-context cancellation through the same acceptance path;
- cancellation disposition:
  - `accepted`;
  - `already_requested`;
  - `already_terminal`;
- first accepted cancellation reason retention;
- exactly one `run.cancellation_requested` event;
- separation of request, acceptance, context propagation, work acknowledgement, and terminal cancellation;
- centralized deterministic terminal-resolution policy;
- success after accepted cancellation;
- cancellation-related error after accepted cancellation;
- unrelated error after accepted cancellation;
- work that delays cancellation acknowledgement;
- work that ignores cancellation;
- completion-cancellation race proof;
- multiple concurrent waiters;
- multiple event observers;
- replay from event sequence;
- observer registration before, during, and after terminal commitment;
- slow and abandoned observers;
- goroutine ownership and leak resistance;
- transition atomicity under concurrent reads;
- stress and race testing;
- diagnostic CLI scenarios;
- final lifecycle contract;
- phase review and simplification review.

#### Notes

Likely candidates are:

- `reasoning/cancellation-and-terminal-resolution.md` — acceptance, causes, outcome truth table, parent context, ignored cancellation;
- `reasoning/concurrency.md` — cancellation/completion interleavings, single-writer commitment, waiter and observer races, goroutine ownership;
- `reasoning/events.md` — retained history, cursor semantics, replay, stream completion, abandoned observers;
- `reasoning/api-and-authority.md` — view/control separation, cancellation result surface, observer API;
- `reasoning/testing.md` — race harness, stress strategy, leak checks, deterministic orchestration;
- `reasoning/architecture.md` only if Sprint 2 reveals a material package or ownership change.

The sprint index should avoid repeating Sprint 1 area reasoning unchanged. It should create new area documents where Sprint 2 introduces genuine uncertainty and point directly to Sprint 1 reasoning elsewhere.

#### Evidence

Strong candidates from the agent-harness study:

- `01.02-control-flow-ownership.md`;
- `01.04-termination-and-loop-bounds.md` selected only for stop-reason and runtime-ownership evidence, not loop budgets;
- `01.05-pause-resume-interrupt-semantics.md`;
- `01.07-concurrency-and-parallel-advancement.md`;
- `01.09-delivery-guarantees-and-idempotency.md`;
- `01.10-replay-and-determinism.md`;
- `02.04-mutation-discipline-and-state-transitions.md`;
- `03.09-completion-and-finalization-semantics.md`.

Supporting Go reports may include state/context, concurrency, IO abstraction, testing strategy, error handling, and philosophy.

#### Deliverables

Normal UltraPlan sprint artifacts plus:

- cancellation implementation;
- deterministic terminal-resolution function;
- replayable multi-observer event history;
- cancellation and race tests;
- leak-resistance evidence;
- diagnostic CLI;
- promoted execution-lifecycle contract;
- Phase 1 review.

#### Commands

```text
aren dev run success
aren dev run fail
aren dev run cancel
aren dev run race
```

Recommended:

```text
aren dev run panic
aren dev run parent-cancel
aren dev run ignore-cancel
```

#### Acceptance

Sprint 2 is complete only when:

- explicit and parent cancellation use one acceptance path;
- repeated cancellation is idempotent;
- cancellation disposition is observable and stable;
- exactly one cancellation-request event is recorded;
- terminal cancellation is never claimed before work returns;
- success, cancellation-related failure, unrelated failure, and panic resolve according to one documented policy;
- scheduler timing cannot reinterpret an identical committed fact set;
- multiple waiters and observers remain coherent;
- replay from sequence reconstructs canonical history;
- slow or abandoned observers cannot block execution or leak Aren-owned producer goroutines;
- high-iteration race tests produce no duplicate terminal outcome, inconsistent history, deadlock, or data race;
- all tests pass under `go test -race`;
- the diagnostic CLI exercises the real runtime;
- the final lifecycle contract matches the realised implementation;
- phase review finds no unresolved foundational ambiguity.

#### Evidence

Expected commands should include equivalents of:

```bash
go test ./...
go test -race ./...
# repeated targeted race/stress tests
aren dev run success
aren dev run fail
aren dev run cancel
aren dev run race
```

The plan should name exact targeted test commands after reasoning selects package and test names.

---

## 4. Phase Exit Gate

Aren Phase 1 is complete only when both sprints are accepted and the combined evidence proves:

1. Aren owns one coherent run lifecycle.
2. Lifecycle transition is the atomic bookkeeping unit.
3. Exactly one terminal outcome and terminal event are committed.
4. Terminal resolution is centralized and deterministic.
5. Cancellation request, acceptance, observation, and terminal cancellation remain distinct.
6. State, event, timing, cancellation metadata, outcome, and waiter release become visible coherently.
7. Event history is ordered and replayable in memory by `(run_id, sequence)`.
8. Multiple waiters and observers are safe.
9. Work error, panic, and Aren invariant failure remain distinguishable.
10. Arbitrary work effects are explicitly outside exactly-once and rollback guarantees.
11. Race, stress, negative, panic, observer, waiter, and leak tests pass.
12. Diagnostic execution provides runnable evidence outside isolated unit assertions.
13. The final lifecycle contract is promoted and reflects tested reality.
14. No later-phase feature has entered through speculative infrastructure.

Open questions affecting these conditions block the phase. Naming, package-layout details already isolated behind tested behaviour, and questions belonging solely to later execution types do not.

---

## 5. Optional Third Sprint Rule

Do not schedule a third sprint upfront.

A third Phase Closure sprint may be added only when Sprint 2 review demonstrates a distinct, bounded body of required work that cannot be completed honestly inside Sprint 2, such as:

- substantial model-based or property-based verification;
- a separate real-runtime smoke harness;
- contract promotion requiring broad documentation reconciliation;
- simplification work that materially changes public semantics.

A third sprint must not be created merely to move unfinished tests, documentation, or review out of Sprint 2.

---

## 6. Deferred Beyond This Project

The following belong to later Aren UltraPlan projects:

- controlled progress and partial output;
- executor abstraction pressure from multiple execution types;
- model providers and model-call semantics;
- structured output and validation;
- retries and attempts;
- tools and tool execution;
- subprocess supervision;
- persistent state and restart recovery;
- pause/resume and approval;
- agent loops;
- workflows and routing;
- daemon hosting;
- remote APIs and multi-language clients.

Phase 1 may record questions about these topics but must not design or implement them.

