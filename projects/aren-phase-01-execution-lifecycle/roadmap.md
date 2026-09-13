# Aren Phase 1 Roadmap — Execution Lifecycle

> Project: `aren-phase-01-execution-lifecycle`  
> Target repository: `../Aren/`  
> Scope: define, prove, measure, and visually explain the lifecycle of one supervised in-process execution.

## 1. Scope Principle

This project covers only Aren Phase 1.

Its purpose is to prove that Aren can own one execution lifecycle without the supervised work relying on an LLM, subprocess, network call, persistent store, workflow engine, or persistent daemon. Phase 1 also establishes a repeatable performance baseline and a process-scoped browser view of the evidence that lifecycle produces.

The project will not design later Aren capabilities in advance. Provider integration, tools, retries, persistence, workflows, pause/resume, remote execution, persistent daemon hosting, and a universal executor abstraction remain deferred.

The implementation is divided into four sprints because Phase 1 has four distinct uncertainties:

1. Can Aren establish one coherent lifecycle and terminal outcome at all?
2. Does that lifecycle remain truthful and race-free under cancellation, observation, and adversarial concurrency?
3. Can Aren measure the realised lifecycle with a repeatable method that exposes overhead, scaling, memory, runtime-task growth, and contention?
4. Can a browser frontend make that evidence understandable without becoming another source of lifecycle truth or affecting execution?

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
Sprint 3: Performance Methodology And Baseline
    ↓
Sprint 4: Observability Frontend
    ↓
Phase review + promoted lifecycle and observation contracts
```

Each sprint uses the normal UltraPlan chain:

```text
requirements -> sprint index -> technical handbook -> area reasoning -> sprint reasoning -> plan -> execute -> review
```

The technical handbook distills selected studies. It does not decide architecture or implementation. Area reasoning investigates the sprint from selected perspectives. Top-level `reasoning.md` makes final sprint decisions. `plan.md` executes those decisions and must not invent new architecture.

---

## 3. Cross-sprint carry-forward rule

No sprint may rediscover earlier accepted work from scratch.

Each sprint after Sprint 1 must include the completed artifacts of every sprint it depends on:

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

Sprint 3 must add the equivalent completed Sprint 2 paths. Sprint 4 must add the equivalent completed Sprint 2 and Sprint 3 paths. Each `sprint-index.md` must classify them as prior project decisions and realised implementation context.

Rules:

- Sprint 1 `reasoning.md` is the authoritative synthesis of Sprint 1 decisions.
- Sprint 1 area reasoning remains available for detailed rationale and rejected alternatives.
- Sprint 1 technical-handbook evidence may be reused, but Sprint 2 should select additional reports targeted to cancellation, event delivery, and concurrency rather than regenerate the same handbook blindly.
- A later sprint may supersede an earlier decision only when new evidence or realised implementation behaviour proves it insufficient or incorrect.
- A superseding decision must name the prior decision, explain the new evidence, describe the impact, and update the implementation plan accordingly.
- Silence does not supersede a prior decision.

This carry-forward rule should be copied into each later sprint's requirements and sprint index when that sprint is initialized.

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
- lifecycle hardening review and simplification review.

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
- lifecycle hardening review.

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
- the lifecycle hardening review finds no unresolved foundational ambiguity before measurement and UI work begin.

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

## Implementation Wave 3 — Performance Methodology And Baseline

### Sprint 3: Performance Methodology And Baseline

> Slug: 03-performance-methodology-and-baseline
> Status: planned
> Depends On: 1, 2

#### Goal

Establish a repeatable performance experiment method over the realised Phase 1 lifecycle and record the first attributable baseline.

At sprint completion, Aren must be able to measure lifecycle overhead, concurrency scaling, allocation, retained memory, runtime-task growth, tail latency, and contention without confusing controlled-work time with Aren cost.

#### Uncertainty

> Can Aren produce repeatable and interpretable performance evidence that remains comparable across implementation changes and later phases?

#### Build

- canonical synthetic workloads for immediate success, returned failure, cancellation, completed wait, event replay, concurrent completion, and blocked active runs;
- microbenchmarks for lifecycle operations supported by the realised API;
- bounded concurrency scenarios at the canonical Phase 1 levels that the host can run safely;
- latency distribution, throughput, allocation, peak and retained memory, and runtime-task measurement;
- CPU, heap, goroutine, mutex, block, and runtime-trace capture commands;
- source revision, dirty state, Go toolchain, operating system, architecture, CPU, runtime settings, scenario parameters, duration, and sample metadata;
- machine-readable result records;
- human-readable statistical comparison using `benchstat` or an equivalent selected during reasoning;
- quick local, pull-request, and extended execution tiers;
- measured variance on the intended comparison environment;
- documented rules for adding a benchmark when later phases introduce new costs;
- the first Phase 1 performance baseline and interpretation report.

#### Deferred

- Sprint 3 measures the real Sprint 2 implementation. It does not build a parallel lifecycle model.
- It does not invent universal throughput or latency targets before workload evidence exists.
- It does not add a benchmark database, hosted service, fleet runner, or historical dashboard.
- A benchmark result may justify a focused fix, but no optimization may weaken lifecycle, cancellation, ordering, observation, or failure semantics.
- Correctness failures, leaks, unbounded growth, global serialization of independent runs, and observation that controls execution are release blockers. Other performance results establish the baseline.

#### Deliverables

- benchmark packages and scenario runner selected by Sprint 3 reasoning;
- versioned workload and result schemas;
- reproducible local commands;
- pull-request and extended benchmark entry points;
- profile capture and interpretation instructions;
- machine-readable baseline evidence;
- a human-readable Phase 1 baseline report;
- tests for parsers, metadata, bounds, and benchmark determinism where practical;
- Sprint 3 review.

#### Acceptance

Sprint 3 is complete only when:

- every required workload exercises the real lifecycle implementation;
- Aren time is distinguishable from controlled-work time where the workload contains a delay;
- repeated samples can be compared statistically;
- every retained result identifies its source and measurement environment;
- the quick tier is suitable for ordinary development and the extended tier is explicitly bounded;
- concurrency, blocked-run memory, retained memory, runtime-task growth, and contention have recorded evidence;
- profile capture is demonstrated on at least one benchmark run;
- normal variance is recorded before any regression threshold is proposed;
- benchmark failures cannot be reported as successful measurements;
- measurement code does not alter production lifecycle semantics;
- `go test ./...`, `go test -race ./...`, and the selected benchmark commands pass;
- the review records baseline limitations and does not overstate host-specific results.

---

## Implementation Wave 4 — Observability Frontend

### Sprint 4: Observability Frontend

> Slug: 04-observability-frontend
> Status: planned
> Depends On: 1, 2, 3

#### Goal

Ship the first dedicated Aren browser frontend for inspecting live and completed Phase 1 runs, failures, cancellation decisions, diagnostics, and performance evidence.

At sprint completion, the browser must present the same canonical lifecycle truth as the runtime and CLI through a versioned read-only contract. Browser connection, slowness, disconnection, or failure must not affect execution.

#### Uncertainty

> Can Aren make its first lifecycle and performance evidence directly understandable in a browser without duplicating runtime semantics or introducing daemon architecture?

#### Build

- a versioned read-only observation DTO distinct from internal runtime types;
- a process-scoped local observation service that binds to loopback;
- a development command that runs or attaches the UI to a bounded diagnostic scenario without creating a persistent daemon;
- run identity, state, timing, outcome, and cancellation summary;
- an ordered event timeline keyed by sequence rather than timestamp;
- structured returned-error, panic, cancellation, and invariant-diagnostic presentation;
- visible distinction between canonical facts, diagnostics, and performance measurements;
- live updates and complete late-connection reconstruction from retained in-memory history;
- explicit empty, connecting, live, terminal, disconnected, malformed-evidence, and unsupported-version states;
- keyboard navigation, visible focus, semantic landmarks and headings, non-colour state labels, reduced-motion support, and responsive layouts;
- bounded payload, update, reconnection, and rendering behaviour;
- agreement, observer-isolation, transport, frontend, accessibility, and browser tests.

#### Deferred

- The runtime remains the only owner of lifecycle truth. The service and frontend only project it.
- The browser cannot commit transitions, cancel runs, or construct outcomes.
- The local service is not durable, remotely accessible, multi-user, or a general Aren API.
- Sprint 4 does not add authentication, production telemetry export, a global event stream, or persistent observation storage.
- A UI requirement may not force false runtime events or instrumentation into the execution critical path.

#### Deliverables

- browser frontend and its design-system foundation;
- process-scoped loopback observation service;
- explicit observation schema and compatibility tests;
- live and completed-run views;
- lifecycle timeline and structured outcome, failure, cancellation, diagnostic, and performance views;
- automated frontend, transport, integration, accessibility, and observer-isolation tests;
- diagnostic commands for required lifecycle scenarios;
- promoted lifecycle and observation contracts;
- final Phase 1 review.

#### Commands

Exact command names remain a Sprint 4 reasoning decision. The delivered commands must cover equivalents of:

```text
aren dev observe success
aren dev observe fail
aren dev observe cancel
aren dev observe race
```

#### Acceptance

Sprint 4 is complete only when:

- runtime, CLI, observation DTO, and browser agree on identity, state, sequence, timing, cancellation facts, failure, and outcome;
- the UI never derives canonical truth from prose logs;
- live updates preserve committed sequence order;
- a late browser connection reconstructs the complete retained Phase 1 history;
- slow, disconnected, or faulty clients cannot block lifecycle progress or leak Aren-owned producer tasks;
- malformed evidence and unsupported schema versions fail visibly without presenting partial data as canonical;
- loading, empty, live, terminal, and disconnected states are understandable;
- keyboard-only operation works and focus remains visible and predictable;
- status does not rely on colour, reduced-motion preferences are respected, and layouts remain usable at narrow and wide widths;
- payload and rendering costs remain bounded and are measured with the Sprint 3 method where applicable;
- the observation service binds only to loopback by default and stops with its owning process;
- no durable daemon, public control API, or second lifecycle model has appeared;
- all runtime, race, frontend, accessibility, and browser integration checks pass;
- the final contracts and review match the realised implementation.

---

## 4. Phase Exit Gate

Aren Phase 1 is complete only when all four sprints are accepted and the combined evidence proves:

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
13. Performance workloads, environment capture, sampling, comparison, and profiling are reproducible.
14. The Phase 1 baseline records lifecycle overhead, concurrency scaling, memory, runtime-task growth, tail latency, and contention without inventing unsupported universal targets.
15. The browser frontend presents runtime-owned truth through a versioned read-only contract.
16. Browser observation remains passive under slow consumption, disconnection, malformed evidence, and failure.
17. The process-scoped local observation service remains distinct from persistent daemon hosting and remote APIs.
18. Accessibility, responsive behaviour, and bounded frontend performance are proven.
19. The final lifecycle and observation contracts are promoted and reflect tested reality.
20. No later-phase feature has entered through speculative infrastructure.

Open questions affecting these conditions block the phase. Naming, package-layout details already isolated behind tested behaviour, and questions belonging solely to later execution types do not.

---

## 5. Deferred Beyond This Project

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
- persistent daemon hosting;
- remote or general-purpose APIs and multi-language clients;
- durable observability storage and production telemetry export;
- multi-user observability and authentication.

Phase 1 may record questions about these topics but must not design or implement them.
