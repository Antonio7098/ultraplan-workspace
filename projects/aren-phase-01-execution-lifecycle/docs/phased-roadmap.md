# Aren Phased Roadmap

## 1. Purpose

Aren will be developed from a small, local execution runtime into a broader autonomous execution system.

The long-term vision is:

> A general-purpose runtime that owns and supervises autonomous execution independently of any particular agent, provider, language, transport, or deployment model.

That vision is directional rather than an initial specification.

Aren will not begin by trying to represent coding agents, model calls, shell commands, workflows, approvals, and scheduled tasks through one complete universal abstraction. It will begin by proving a small execution model against controlled implementations, then broaden only when real variation provides evidence for new abstractions.

The intended progression is:

```text
Execution lifecycle
    ↓
Controlled execution
    ↓
Real model invocation
    ↓
Streaming and structured results
    ↓
Retries and repair
    ↓
Tool execution
    ↓
Bounded agent loop
    ↓
Context management
    ↓
Persistence and recovery
    ↓
Execution composition
    ↓
Policies and resource governance
    ↓
Daemon hosting and remote clients
    ↓
Workflows
    ↓
Additional execution types
    ↓
Distributed execution, if justified
```

Later phases are hypotheses, not commitments. Their order and contents should change in response to evidence from real Aren usage.

---

## 2. Development Rules

### 2.1 One primary uncertainty per phase

Each phase should answer one major technical question. A phase may contain several implementation tasks, but they must converge on that question.

### 2.2 Every phase must produce something runnable

Types, interfaces, and design documents are not sufficient by themselves. Every phase must expose a working vertical slice through one or more of:

- an importable library;
- a development CLI;
- an executable example;
- a real smoke workflow.

### 2.3 Behaviour before abstraction

Start with one concrete implementation.

Introduce a general interface only when there is real evidence, such as:

- a second implementation with meaningful behavioural differences;
- repeated conditional logic;
- duplicated lifecycle handling;
- a testing boundary needed to reproduce failures;
- ownership that genuinely needs to cross a boundary;
- a concept that has remained stable across repeated use.

The possibility that an abstraction may be useful later is not enough.

### 2.4 Failure testing is mandatory

Every phase must validate, where relevant:

- normal operation;
- expected failures;
- cancellation;
- timing and race conditions;
- cleanup;
- partial results;
- recovery or honest termination.

### 2.5 Real use before progression

After implementation and automated testing, the phase capability must be exercised in a small real workflow.

The phase review should ask:

- Is the API awkward?
- Are important events missing?
- Are any events redundant?
- Are failures sufficiently classified?
- Is ownership clear?
- Is hidden state affecting behaviour?
- Did an abstraction appear before it was necessary?
- Should anything be removed before continuing?

### 2.6 Reviews may remove capabilities

A phase review may conclude that:

- an interface should become concrete;
- a package should be collapsed;
- a configuration option should be removed;
- an event should be merged or deleted;
- a planned capability is not yet justified.

Reduction is valid progress.

### 2.7 Foundations must be vertical slices

“Foundations first” does not mean building a large invisible substrate. Each foundation must be exercised through a complete, observable path that reveals real behaviour.

### 2.8 Later phases remain revisable

The first phases establish Aren’s core semantics. Later phases are deliberately less detailed and must be reconsidered after every phase review.

---

## 3. Phase Overview

| Phase | Capability | Primary uncertainty |
|---|---|---|
| 0 | Project foundation | Can Aren maintain strict scope and trustworthy engineering feedback? |
| 1 | Execution lifecycle | Can Aren define and enforce a small, coherent execution lifecycle? |
| 2 | Controlled executor | Can that lifecycle survive progress, failure, cancellation, cleanup, and races? |
| 3 | First real model invocation | Can Aren own one real LLM request without losing lifecycle control? |
| 4 | Streaming model execution | Can streamed output remain ordered, cancellable, bounded, and observable? |
| 5 | Structured model results | Can Aren distinguish valid task results from successful provider calls? |
| 6 | Retry and bounded repair | Can Aren recover selectively without hiding failures or creating uncontrolled loops? |
| 7 | Tool-call representation | Can requested actions be represented independently from their implementation? |
| 8 | Local tool execution | Can Aren supervise external actions through the same reliable principles? |
| 9 | Bounded agent loop | Can Aren own a minimal model–tool loop with explicit termination conditions? |
| 10 | Context management | Can context be managed deliberately without hiding mutation or provenance? |
| 11 | Persistence and recovery | Which state genuinely needs to survive process termination? |
| 12 | Execution composition | How should executions be coordinated before a workflow system is justified? |
| 13 | Policies and resources | Which permissions, budgets, and concurrency controls require runtime ownership? |
| 14 | Daemon hosting | When should execution ownership outlive a single client process? |
| 15 | Multi-language clients | Can other languages control Aren without duplicating runtime behaviour? |
| 16 | Workflows | Can reusable, inspectable, and recoverable processes be built from proven execution primitives? |
| 17 | Broader execution types | Which non-agent forms of work genuinely belong under Aren? |
| 18 | Distributed execution | Is there sufficient evidence for remote workers or clustering? |

---

# Foundation

## Phase 0 — Project Foundation and Scope Control

### Primary question

> Can Aren establish a development environment and decision process that make incorrect behaviour visible without building product infrastructure prematurely?

Phase 0 is not an architecture phase. It creates only the minimum repository foundation needed to develop Phase 1 safely.

### Goals

- Select one implementation language for the long-term runtime, then establish a small repository for it.
- Make builds, tests, static checks, and concurrency-safety checks easy to run.
- Define project terminology before APIs proliferate.
- Record explicit scope boundaries.
- Create lightweight architectural decision records.
- Prevent roadmap work from turning into implementation speculation.

### In scope

#### Repository foundation

A minimal repository should contain runtime source, a development CLI entry
point, tests, examples, decision records, a glossary, the selected language's
manifest, and a README.

The exact structure should remain small and may change during Phase 1.

#### Engineering feedback

Documented commands for:

- formatting;
- unit tests;
- concurrency-safety tests;
- static analysis;
- coverage inspection;
- building the CLI.

#### Initial glossary

Define only the terms already required:

- execution;
- run;
- state;
- event;
- result;
- failure;
- cancellation;
- executor.

Definitions may be marked provisional.

#### Scope ledger

Maintain a lightweight record containing:

- current phase scope;
- explicit exclusions;
- deferred ideas;
- evidence required to reconsider each deferred idea.

#### Decision records

Use short decision records only for choices that would otherwise become ambiguous.

Likely initial decisions:

- one implementation language for the long-term runtime;
- library-first runtime;
- no daemon in the initial phases;
- local and in-memory operation initially;
- no stable public API promise before the foundational semantics settle.

### Out of scope

- provider integrations;
- execution persistence;
- databases;
- servers;
- plugin systems;
- workflow definitions;
- SDK generation;
- broad configuration frameworks;
- observability backends;
- production deployment;
- semantic-versioning guarantees.

### Deliverables

- buildable implementation module;
- minimal CLI entry point;
- automated test and concurrency-check commands;
- glossary;
- phase scope ledger;
- decision-record format;
- initial repository README.

### Validation

- A clean checkout can build and test with one documented command.
- The selected toolchain's strongest practical concurrency checks run locally and in CI.
- No production runtime abstractions are introduced solely for repository setup.
- Every nontrivial source module has an immediate purpose.

### Exit gate

Phase 0 is complete when Aren has selected one implementation language for the
long-term runtime and the repository reliably supports Phase 1 development
without carrying speculative runtime architecture.

---

## Phase 1 — Execution Lifecycle

### Primary question

> Can Aren define and enforce the lifecycle of one execution without depending on an LLM, subprocess, network call, or persistent store?

This is the most foundational phase.

The objective is not to create a universal executor interface. It is to establish the minimum semantics that Aren itself owns.

### Conceptual model

A run represents one supervised occurrence of work.

A preliminary lifecycle is:

```text
created
   ↓
running
   ├──→ succeeded
   ├──→ failed
   └──→ cancelled
```

A cancellation request may occur while running:

```text
running
   ↓
cancellation requested
   ↓
cancelled
```

Whether `cancellation_requested` is a state, an event, or both must be resolved through implementation. The lifecycle should remain much smaller than a workflow-engine state machine.

### Goals

- Assign every run a unique identity.
- Define who owns state transitions.
- Define all legal and illegal transitions.
- Represent terminal outcomes explicitly.
- Separate lifecycle state from output data.
- Establish cancellation semantics.
- Establish the first canonical event vocabulary.
- Define event observation behaviour.
- Define a minimum failure representation.
- Make invariant violations visible.

### Run identity

Decide:

- identifier type;
- generation ownership;
- whether callers may provide IDs;
- whether IDs carry semantic information.

Initial direction:

- opaque unique IDs;
- generated by Aren;
- no execution type or timestamp encoded into the ID.

### Lifecycle state

Candidate states:

- `created`;
- `running`;
- `succeeded`;
- `failed`;
- `cancelled`.

Questions to resolve:

- Is `created` externally observable?
- Is cancellation request a state or only an event?
- Can a run fail before entering `running`?
- Is rejection a failed run or a failure to create one?
- Can cleanup failure alter an otherwise successful outcome?

### Terminal outcome

State alone is insufficient. A terminal outcome may include:

- terminal state;
- result value;
- classified failure;
- start time;
- finish time;
- cancellation metadata.

Result and failure must remain coherent.

Valid examples:

```text
succeeded + result
failed + failure
cancelled + cancellation reason
```

Invalid combinations include:

```text
succeeded + failure
failed + successful result
running + terminal timestamp
```

### State ownership

Aren, not executor code, owns legal lifecycle transitions.

Executed work may return an outcome or report progress, but it must not directly mutate the run state machine.

### Cancellation

Cancellation semantics must answer:

- Who may request cancellation?
- Is cancellation idempotent?
- What happens when cancellation races with completion?
- Is cancellation cooperative?
- When is a run considered cancelled?
- What happens if work ignores cancellation?
- Does a cancellation reason belong in the outcome?
- Are repeated cancellation requests observable?

Initial direction:

- cancellation requests are idempotent;
- the first valid terminal outcome wins;
- Aren propagates an implementation-appropriate cancellation signal;
- Aren does not claim the work has stopped until controlled work returns;
- cancellation request and cancellation completion are distinct;
- a successful outcome may remain successful if it wins the race against cancellation.

These rules must be proved by tests, not accepted only in prose.

### Event ordering

The initial event vocabulary should describe lifecycle behaviour only.

Candidate events:

- `run.created`;
- `run.started`;
- `run.cancellation_requested`;
- `run.succeeded`;
- `run.failed`;
- `run.cancelled`.

Likely event metadata:

- run ID;
- per-run sequence number;
- event type;
- occurrence time;
- immutable payload.

Initial direction:

- ordering is established by per-run sequence number, not timestamp;
- state transition and event creation occur together inside the runtime;
- delivery may lag but must preserve order;
- consumer failure must not corrupt run state;
- terminal events occur exactly once.

### Event observation

Avoid building a general event bus.

Possible APIs include:

```text
run.events()
```

or:

```text
subscription = run.subscribe()
```

A raw channel creates semantic questions that must be answered explicitly:

- buffer size;
- slow consumers;
- closure ownership;
- missed events;
- multiple subscribers;
- replay behaviour;
- observer abandonment.

The implementation should choose the smallest observation model whose behaviour can be tested honestly.

### Failure model

Begin with a minimal structured failure representation.

Likely categories:

- execution failure;
- cancellation;
- internal invariant violation.

The representation should preserve:

- machine-readable category;
- human-readable message;
- wrapped cause where applicable;
- whether the failure arose from executed work or Aren itself.

Do not design the future provider and tool error taxonomies yet.

### Time

The runtime needs observable timing information:

- start time;
- terminal time;
- possibly duration.

A public clock abstraction should only be introduced if deterministic testing proves that it is necessary.

### First implementation

Use an in-process work function as a semantic instrument:

```text
work(cancellation) -> result | failure
```

A run controller should:

1. create the run identity;
2. emit creation and start events;
3. invoke the work function;
4. propagate cancellation;
5. classify the returned outcome;
6. perform one legal terminal transition;
7. expose the immutable outcome and events.

This is not yet a promise of the permanent executor API.

### Illustrative public surface

```text
run = aren.start(parent, work)

for event in run.events():
    observe(event)

outcome = run.wait()
run.cancel("user requested cancellation")
```

The PRD should specify behaviour before exact names and signatures.

### Required tests

#### Normal lifecycle

- `created → running → succeeded`;
- result preserved exactly;
- timestamps recorded coherently;
- events ordered correctly;
- waiting after completion returns immediately.

#### Failure lifecycle

- `created → running → failed`;
- original cause retained;
- failure event emitted exactly once;
- no successful result exposed.

#### Cancellation

- cancellation before work begins, if supported;
- cancellation during execution;
- repeated cancellation;
- work that cooperates immediately;
- work that delays acknowledgement;
- work that ignores cancellation until it returns.

#### Completion–cancellation race

Run many iterations with the selected toolchain's strongest practical
concurrency checks enabled.

Required invariants:

- exactly one terminal state;
- exactly one terminal event;
- no data race;
- result agrees with terminal state.

#### Multiple waiters

- all waiters receive the same immutable outcome;
- no waiter consumes the outcome exclusively;
- no deadlock.

#### Event consumers

- slow consumer;
- absent consumer;
- abandoned consumer;
- multiple consumers if supported;
- documented overflow behaviour;
- no leak caused solely by nobody reading events.

#### Panic and invariants

- work panic;
- illegal state transition;
- no silent invariant corruption.

A panic may become an explicitly classified internal execution failure, but it must remain recognisable as a panic rather than being disguised as an ordinary task failure.

### Diagnostic CLI

Provide a small development CLI:

```text
aren dev run success
aren dev run fail
aren dev run cancel
aren dev run race
```

Example output:

```text
000 run.created
001 run.started
002 run.cancellation_requested
003 run.cancelled
```

The CLI is a semantic inspection tool, not the final Aren product interface.

### Deliverables

- lifecycle model;
- run controller;
- terminal outcome model;
- cancellation contract;
- initial event vocabulary;
- deterministic event sequencing;
- minimum failure representation;
- comprehensive unit and concurrency tests;
- diagnostic CLI;
- written lifecycle contract;
- phase review.

### Out of scope

- attempts and retries;
- providers;
- model messages;
- streaming output;
- tools;
- subprocess management;
- persistence;
- replay after process restart;
- multiple execution types;
- daemon transport;
- global event streams;
- workflow composition;
- production telemetry export.

### Exit gate

Phase 1 is complete only when:

1. all lifecycle transitions are defined;
2. cancellation races are proven under repeated concurrency checking;
3. event order is deterministic;
4. exactly one terminal outcome is guaranteed;
5. slow or absent consumers cannot accidentally deadlock or leak a run;
6. a runnable demonstration exists;
7. no foundational semantic ambiguity remains unresolved.

---

## Phase 2 — Controlled Executor and Failure Laboratory

### Primary question

> Does the Phase 1 lifecycle remain coherent when execution exhibits realistic progress, delay, failure, cleanup, partial output, and cancellation behaviour?

Phase 1 proves the lifecycle around a minimal work function. Phase 2 introduces a deliberately configurable executor used to attack those semantics.

It must not grow into a general simulation framework.

### Goals

- Exercise richer controlled behaviour.
- Introduce nonterminal progress.
- Test delayed and partially cooperative cancellation.
- Define cleanup behaviour.
- Test partial output followed by failure.
- Determine whether an executor interface is now genuinely earned.
- Produce a reusable conformance suite for future executors.

### Controlled behaviours

The test executor should be able to:

- succeed immediately;
- succeed after a delay;
- fail immediately;
- fail after progress;
- emit a fixed progress sequence;
- block until externally released;
- cooperate with cancellation;
- delay cancellation response;
- ignore cancellation until a safe point;
- panic;
- fail during cleanup;
- race completion against cancellation;
- produce partial output before failure.

Configuration should be typed and explicit. Do not create a scenario language.

### Progress

Introduce only the smallest useful progress representation.

Initial candidate:

```go
type Progress struct {
    Message string
}
```

A generic metadata bag or universal progress schema should not be introduced without concrete need.

### Cleanup semantics

Decide:

- whether cleanup completes before the terminal event;
- whether cleanup failure can convert success into failure;
- how cleanup behaves after cancellation;
- whether cleanup receives a separate timeout;
- how execution and cleanup errors are preserved together.

Initial direction:

- required cleanup completes before the terminal state is published;
- cleanup errors are never silently discarded;
- mandatory cleanup failure prevents an automatic success outcome;
- cleanup must not run forever under an already cancelled signal.

### Executor boundary

A candidate abstraction may now be tested:

```text
executor.execute(cancellation, reporter) -> result | failure
```

It should only be retained if it provides real value through progress reporting, conformance testing, or ownership clarity. If a function remains simpler and sufficient, keep the function.

### Conformance suite

Future executors should be able to reuse behavioural tests for applicable invariants:

- returns one coherent outcome;
- respects the cancellation contract;
- emits no progress after terminal completion;
- preserves progress ordering;
- releases required resources;
- does not mutate published events;
- retains classified failures.

### Real vertical slice

Add one small non-LLM example. Prefer a platform-independent in-process task initially. Introduce subprocess semantics only if the foundation cannot be validated without them.

### Required tests

- progress before success;
- progress before failure;
- cancellation after partial progress;
- no progress after terminal event;
- cleanup on every exit path;
- cleanup failure;
- delayed cancellation acknowledgement;
- observer detachment;
- panic during execution;
- panic during cleanup;
- a large progress burst;
- documented backpressure behaviour;
- run start failure, if supported.

### Deliverables

- controlled semantic executor;
- progress mechanism;
- cleanup contract;
- executor conformance suite;
- evidence-based decision on an executor interface;
- one runnable non-LLM example;
- revised lifecycle contract;
- phase review.

### Out of scope

- retries;
- attempts unless demanded by evidence;
- model providers;
- arbitrary subprocess supervision;
- remote execution;
- durable progress;
- global event routing;
- plugin registration.

### Exit gate

Phase 2 is complete when the lifecycle contract has survived realistic controlled behaviour and future executor integrations have a clear, tested path.

---

# Model Execution

## Phase 3 — First Real Model Invocation

### Primary question

> Can Aren own one real model request while preserving execution, cancellation, event, and failure semantics?

This is the first external integration. The goal is one excellent concrete model invocation, not provider neutrality.

### Initial scope

Support one current OpenAI API path with:

- one model request type;
- text input;
- text output;
- no tools;
- no structured schema;
- no automatic retry;
- no session abstraction;
- non-streaming operation initially.

### Goals

- Send one real request.
- Map provider activity into Aren lifecycle semantics.
- Propagate cancellation into the network request.
- Capture response metadata honestly.
- Classify the smallest actionable provider failures.
- Preserve provider-specific details behind a bounded escape hatch.
- Validate behaviour with deterministic transport tests and real smoke tests.

### Input and result

Begin with a concrete request rather than a universal message graph.

```go
type ModelRequest struct {
    Model string
    Input string
}
```

Capture, where available:

- generated text;
- provider request ID;
- actual model used;
- stop or completion reason;
- token usage;
- provider metadata required for diagnosis.

Unknown values must remain unknown rather than becoming zero or empty success values.

### Initial provider failures

Classify only failures callers may handle differently:

- authentication;
- invalid request;
- rate limited;
- provider unavailable;
- timeout;
- network or transport failure;
- malformed provider response;
- cancelled;
- unknown provider failure.

Retain causes and relevant metadata such as retry-after information.

### Candidate events

- `model.request.started`;
- `model.request.completed`;
- `model.request.failed`.

Example sequence:

```text
run.started
model.request.started
model.request.completed
run.succeeded
```

Avoid duplicating low-level HTTP details.

### Testing

Use a controlled local HTTP server or injectable transport to simulate:

- success;
- authentication failure;
- rate limiting;
- retry-after metadata;
- delayed response;
- cancellation;
- malformed or partial bodies;
- connection closure;
- provider 5xx;
- unknown status;
- missing usage;
- unknown response fields.

Real smoke tests should be opt-in and use environment-provided credentials.

### Out of scope

- streaming;
- multiple providers;
- provider abstraction;
- fallback;
- retry;
- structured output;
- tool calls;
- prompt templating;
- context reduction;
- persistent sessions.

### Exit gate

One real model invocation is reliable, cancellable, observable, and honestly represented under normal and adverse transport conditions.

---

## Phase 4 — Streaming Model Execution

### Primary question

> Can Aren expose incremental model output while preserving event order, cancellation, terminal consistency, and consumer safety?

Streaming changes runtime concurrency and observation semantics. It is not only a UI feature.

### Goals

- Receive provider streaming responses.
- Emit output deltas in deterministic order.
- Reconstruct the final output exactly.
- Cancel midstream.
- Handle partial and interrupted output honestly.
- Define bounded buffering and backpressure.

### Candidate events

- `model.output.delta`;
- possibly `model.stream.interrupted` if it adds information not already represented by request failure.

A likely sequence is:

```text
run.started
model.request.started
model.output.delta
model.output.delta
model.request.completed
run.succeeded
```

### Delta semantics

Define:

- what a delta contains;
- whether empty deltas are emitted;
- how Unicode boundaries are handled;
- whether final text must equal the ordered concatenation of text deltas;
- how provider annotations are represented;
- whether raw provider chunks are exposed.

Initial direction:

- emit normalised text deltas for the supported request type;
- retain final reconstructed text in the result;
- keep raw chunks internal or behind an explicit low-level path;
- order all events through the run sequence.

### Partial output

A cancelled or interrupted stream may have useful partial output. It should be retained as diagnostic execution data, but not presented as a complete successful model result.

The result must distinguish:

- complete output;
- partial output caused by cancellation;
- partial output caused by provider or transport failure.

### Backpressure

The implementation must choose and document one bounded policy. It must not accidentally create:

- unbounded memory growth;
- a runtime task per event;
- provider deadlock caused by an unread public channel;
- loss of the terminal outcome.

### Required tests

- one delta and many deltas;
- Unicode boundaries;
- empty stream;
- cancellation before first delta;
- cancellation after several deltas;
- interrupted stream;
- malformed event;
- duplicate provider completion signal;
- provider EOF without an explicit final marker;
- slow or absent observer;
- exact final reconstruction;
- partial-output classification;
- no events after the run terminal event.

### Exit gate

Streaming is useful without weakening lifecycle correctness or introducing hidden concurrency behaviour.

---

## Phase 5 — Structured Model Results and Validation

### Primary question

> Can Aren distinguish provider completion from successful task completion through explicit result validation?

```text
Provider execution succeeded
            ≠
Requested result is valid
```

### Goals

- Request structured output through one supported mechanism.
- Parse the returned representation.
- Validate it against one explicit contract.
- Distinguish parsing and validation failures from provider failures.
- Preserve raw output for diagnosis.
- Avoid coupling every execution to JSON or schemas.

### Result layers

Keep these distinct:

1. provider response;
2. raw generated content;
3. parsed representation;
4. validated result.

A malformed result should retain enough evidence to explain what was returned and why it failed.

### Initial validation scope

- one schema mechanism;
- deterministic parsing and validation;
- no LLM evaluator;
- no automatic repair;
- no generic validator registry;
- no arbitrary network-backed validation.

### Candidate events

- `validation.started`;
- `validation.succeeded`;
- `validation.failed`.

Parsing may be represented as one kind of validation failure unless callers need a different control path.

### Outcome semantics

A provider call that returns malformed structured output may have succeeded at the provider layer while the requested structured execution fails overall.

The run should:

- terminate with a validation failure;
- preserve the raw provider response;
- expose no validated result;
- keep provider and task-level evidence distinct.

### Exit gate

Aren can honestly report that a provider responded successfully while the requested result failed validation.

---

## Phase 6 — Retry and Bounded Repair

### Primary question

> Can Aren repeat or repair failed work selectively without hiding the original failure, duplicating effects, or creating uncontrolled loops?

### Goals

- Introduce attempts explicitly.
- Retry only classified transient failures.
- Support cancellable backoff.
- Preserve each attempt’s evidence.
- Define exhaustion clearly.
- Add bounded structured-output repair only after ordinary retry semantics are stable.

### Attempt model

Initial direction:

- an attempt is a subordinate record inside one run;
- attempts have an ordinal and timestamps;
- events identify attempt number;
- one run still reaches exactly one terminal outcome;
- attempts do not initially expose the complete public run API.

Do not recursively define every attempt as a full execution unless later evidence justifies it.

### Retry policy

A concrete policy may contain:

- maximum attempts;
- initial delay;
- maximum delay;
- multiplier;
- optional jitter;
- retryable failure categories.

Default behaviour is no retry unless explicitly configured.

### Candidate events

- `attempt.started`;
- `attempt.failed`;
- `retry.scheduled`;
- `attempt.succeeded`;
- `retry.exhausted`.

The final event set should be reduced after implementation if some events prove redundant.

### Cancellation

Test cancellation:

- during an active request;
- after attempt failure;
- during backoff;
- immediately before another attempt;
- racing with final success.

No wait may ignore cancellation.

### Bounded repair

After retry is stable, support one narrow repair path for invalid structured output.

```text
attempt 1: provider succeeds
validation fails
repair scheduled
attempt 2: original output + validation feedback
validation succeeds
run succeeds
```

Repair must preserve:

- original output;
- validation errors;
- repair instruction;
- repaired output;
- final validation result.

Repair is visible changed-input execution, not a hidden retry.

### Out of scope

- provider fallback;
- hedged requests;
- global retry budgets;
- durable retries after crash;
- automatic tool retries;
- arbitrary repair strategies;
- open-ended self-reflection loops.

### Exit gate

Retries are explicit, bounded, observable, cancellable, and unable to disguise the failures that preceded the final outcome.

---

# Tools and Agents

## Phase 7 — Tool-Call Representation

### Primary question

> Can Aren represent a model’s request for external action independently from how that action is implemented?

This phase separates tool definition from tool execution. It does not yet build the agent loop.

### Goals

- Define a tool for model consumption.
- Represent a requested tool call.
- Parse and validate tool arguments.
- Represent unknown tools without losing the provider response.
- Represent tool results and failures.
- Handle streamed arguments if the provider requires it.
- Avoid coupling tool definitions to implementation-language function signatures.

### Minimum tool definition

- stable name;
- description;
- input schema.

Output schemas should be added only if a real use case requires them.

### Minimum tool call

- provider call ID;
- tool name;
- raw arguments;
- parsed arguments;
- validation status.

An unknown tool call should remain inspectable rather than failing during response parsing.

### Out of scope

- tool registry;
- tool execution;
- parallel tools;
- agent loop;
- remote tools;
- MCP;
- permissions;
- retries;
- client-hosted callbacks.

### Exit gate

Aren can receive, inspect, and validate a requested action without knowing how it will execute.

---

## Phase 8 — Local Tool Execution

### Primary question

> Can Aren supervise a concrete external action with clear lifecycle, result, failure, permission, and cancellation semantics?

Begin with in-process native tools only.

### Goals

- Register a small set of local implementations.
- Resolve validated calls to implementations.
- Execute tools with cancellation.
- Emit tool lifecycle events.
- Validate inputs before side effects.
- Capture typed output or failure.
- Establish an explicit allow or deny permission posture.
- Avoid assuming tools are safe to retry.

### Candidate events

- `tool.started`;
- `tool.completed`;
- `tool.failed`;
- `tool.cancelled`.

### Initial failure categories

- unknown tool;
- invalid arguments;
- permission denied;
- execution failure;
- timeout;
- cancelled;
- invalid result;
- tool panic.

### Permissions

Begin with a direct allow or deny decision. Do not build a policy language.

The key invariant is:

> Required permission is resolved before the tool performs side effects.

### Retry safety

Tool execution is not automatically retried. Idempotency and retry declarations may be added later only when real tools require them.

### Real tools

Start with harmless, useful tools such as:

- reading a text file within an allowed directory;
- listing a directory;
- deterministic calculations.

Do not begin with arbitrary shell execution.

### Exit gate

A model-requested action can be validated and executed locally without weakening Aren’s lifecycle and failure guarantees.

---

## Phase 9 — Minimal Bounded Agent Loop

### Primary question

> Can Aren own a model–tool loop that remains bounded, observable, and understandable?

This is the first conventional agent phase.

### Loop

```text
initial context
    ↓
model invocation
    ↓
assistant output
    ├── final response → stop
    └── tool calls
            ↓
       execute tools
            ↓
       append results
            ↓
       next model invocation
```

### Goals

- Own the loop inside Aren rather than an SDK.
- Support sequential tool calls.
- Maintain explicit conversation state.
- Enforce hard termination conditions.
- Preserve child model and tool evidence.
- Produce one final agent-run outcome.
- Keep control flow linear and direct.

### Required limits

- maximum model turns;
- maximum tool calls;
- maximum elapsed time;
- cancellation;
- optionally token or cost limits when the underlying measurements are trustworthy.

Unknown metrics must not be silently treated as zero.

### Completion semantics

Distinguish:

- model final response;
- exhausted turn limit;
- exhausted tool-call limit;
- unrecoverable model failure;
- tool failure;
- validation failure;
- cancellation.

A model stopping without a final answer is not automatically a successful agent result.

### Initial constraints

- sequential tools only;
- no subagents;
- no planning framework;
- no dynamic tool registration during a run;
- no human approval steps;
- no persistent session;
- no workflow graph.

### Candidate events

```text
run.started
agent.turn.started
model.request.started
model.output.delta
model.request.completed
tool.started
tool.completed
agent.turn.completed
run.succeeded
```

Only retain turn-level events if they add useful boundaries rather than duplicating child events.

### Real workflow

Use Aren to complete one narrow local task, such as inspecting a small directory, reading selected files, and answering a grounded question.

### Exit gate

One bounded agent performs useful work repeatedly without hidden control flow or ambiguous termination.

---

# Longer-Lived Execution

## Phase 10 — Context Engineering

### Primary question

> Can Aren manage growing context deliberately while preserving traceability and avoiding hidden mutation?

Context work should respond to pressure observed in real Phase 9 runs.

### Likely progression

1. Measure context growth.
2. Add deterministic size checks.
3. Separate instructions, history, tool results, and working data.
4. Add simple deterministic reduction.
5. Distinguish durable and reducible entries.
6. Add model-generated summarisation only when deterministic reduction is inadequate.
7. Retain provenance from derived context to original messages.

### Principles

- Original history is not silently rewritten.
- Summaries are explicit derived artefacts.
- Reduction decisions are observable.
- Policies are deterministic where possible.
- Provider token limits do not leak unpredictably into loop logic.
- “Memory” is not introduced as one broad undifferentiated feature.

### Deferred capabilities

- semantic retrieval;
- vector stores;
- long-term memory;
- graph context;
- workspace indexing;
- cross-run memory.

### Exit gate

Aren can support meaningfully longer tasks while every context transformation remains visible and attributable.

---

## Phase 11 — Persistence and Recovery

### Primary question

> Which Aren state must survive process termination, and what recovery guarantees are actually required?

Persistence should follow real lost-work pain, not anticipation.

### Likely progression

1. Persist completed run records.
2. Persist event history.
3. Persist active-run metadata.
4. Detect interrupted runs after restart.
5. Mark interruption honestly.
6. Consider resumability only after interruption semantics are understood.

### Critical distinctions

- recording is not recovery;
- recovery is not resumption;
- resuming an agent loop is not resuming an in-flight network request;
- durable events do not produce exactly-once side effects;
- replay does not imply re-execution safety.

### Storage

Start with one concrete local implementation, selected from demonstrated access patterns. Do not create a general storage abstraction before replacement pressure exists.

### Exit gate

Aren preserves the state users genuinely need and handles interrupted work honestly without claiming unsupported recovery guarantees.

---

## Phase 12 — Execution Composition

### Primary question

> How should several executions be coordinated while keeping control flow explicit and avoiding premature workflow machinery?

Composition should begin as ordinary program control flow.

### Likely progression

- execute A then B;
- branch based on a result;
- run a bounded parallel group;
- propagate cancellation from parent to children;
- retain parent–child run references;
- explain child failures in the parent outcome.

A dependency graph should only appear if repeated direct composition becomes genuinely difficult.

### Principles

- composition is not yet a workflow product;
- child execution ownership is explicit;
- cancellation flows predictably;
- parent completion waits for owned children;
- workflow state is not conflated with run state;
- no visual graph or declarative DSL yet.

### Exit gate

Several executions can be coordinated reliably, and their common composition pressures are understood well enough to inform a later workflow design.

---

## Phase 13 — Policies and Resource Governance

### Primary question

> Which cross-cutting controls require central runtime ownership?

### Likely capabilities

- model and tool allowlists;
- permission decisions;
- token and cost budgets;
- elapsed-time budgets;
- concurrency limits;
- shared rate-limit coordination;
- workspace boundaries;
- approval requirements;
- retry budgets.

### Principles

- policies act on explicit runtime facts;
- unavailable metrics remain unknown;
- denial occurs before side effects;
- policy decisions are observable;
- configuration remains smaller than the controlled behaviour;
- no general policy language without repeated rule complexity.

### Exit gate

Aren safely governs several real forms of execution without spreading duplicated policy logic through applications and executors.

---

# Hosting and External Clients

## Phase 14 — Daemon Hosting

### Primary question

> Has execution ownership outgrown the lifetime and language of a single client process?

A daemon becomes justified by demonstrated needs such as:

- execution surviving client exit;
- multiple clients;
- reconnectable observation;
- central cancellation;
- shared resource controls;
- scheduled execution;
- cross-language access;
- remote management.

### Architecture rule

The daemon hosts the runtime. It does not become the place where lifecycle, retry, tool, or agent-loop behaviour lives.

```text
Aren runtime
    ↓
daemon host
    ↓
transport adapter
```

### Initial daemon scope

- start a run;
- inspect a run;
- subscribe to events;
- cancel a run;
- retrieve its terminal outcome.

Choose one transport based on actual client needs. Do not implement HTTP, WebSockets, SSE, and gRPC simultaneously.

### Command and event distinction

Commands:

- start;
- cancel;
- inspect.

Events:

- lifecycle;
- progress;
- output;
- tool activity;
- terminal outcome.

### Exit gate

A client can disconnect and reconnect without Aren losing execution ownership, while all core semantics remain testable without the daemon.

---

## Phase 15 — Multi-Language Clients and Client-Hosted Tools

### Primary question

> Can Python and JavaScript use Aren naturally without implementing their own runtime behaviour?

### Likely sequence

1. Build one thin client.
2. Support start, observe, cancel, and inspect.
3. Add client-hosted tool callbacks.
4. Define callback timeout and disconnection semantics.
5. Add a second language only after the protocol stabilises.

### Client-hosted tool model

```text
Aren requests tool execution
        ↓
connected client executes callback
        ↓
client returns result
        ↓
Aren continues the loop
```

Questions that must be resolved:

- What happens if the client disconnects?
- Who owns the timeout?
- May another client satisfy the call?
- Can a callback request be replayed?
- How are permissions established?
- Can callbacks be retried safely?

### Exit gate

An application written in another language can use Aren as the execution owner while remaining a thin client rather than becoming a second harness implementation.

---

# Workflows and Broader Execution

## Phase 16 — Workflows

### Primary question

> Can Aren turn proven execution and composition primitives into reusable, inspectable, and recoverable processes without becoming a speculative workflow platform?

This phase is intentionally later than basic execution composition. Aren should first learn how real executions are started, observed, cancelled, retried, persisted, and composed in code. Only then should those patterns be represented as reusable workflow definitions.

### Goals

- Define a reusable process from proven execution primitives.
- Separate workflow definition from workflow instance state.
- Make step transitions explicit and observable.
- Support pause, cancellation, and failure propagation.
- Preserve parent–child execution relationships.
- Recover workflow state after daemon or process restart where supported.
- Introduce human or scheduled steps only when their distinct semantics are understood.

### Initial workflow scope

Begin with a small, linear workflow model:

- named steps;
- sequential execution;
- step inputs from prior outputs;
- explicit terminal states;
- per-step run references;
- cancellation of the workflow and owned child runs;
- persisted workflow-instance state;
- restart detection and honest interruption handling.

A minimal conceptual model:

```text
workflow created
      ↓
step A execution
      ↓
step B execution
      ↓
step C execution
      ↓
workflow completed
```

### Later workflow capabilities, if earned

- conditional branches;
- bounded parallel branches;
- retries at step or workflow level;
- compensation;
- human approval steps;
- scheduled starts;
- waiting for external events;
- reusable subworkflows;
- resumable long-running processes;
- declarative definitions;
- visualisation.

### Critical distinctions

#### Workflow definition vs workflow run

A definition describes reusable structure. A workflow run records one execution of that structure.

#### Workflow state vs execution state

A workflow may be waiting, paused, or blocked even when no child execution is currently running. Those semantics must not be forced into the core run lifecycle.

#### Retry vs resume

Retry starts work again. Resume continues from preserved workflow state. They must not be treated as synonyms.

#### Human waiting vs active execution

A human approval may remain pending for days. It should not pretend to be an active runtime task or network request.

#### Durable orchestration vs exactly-once effects

Persisted state does not guarantee that an external side effect occurred exactly once. Recovery must account for uncertain completion and idempotency.

### Workflow representation

Do not begin with a general graph DSL.

Likely progression:

1. direct typed builder or definition;
2. linear reusable workflow;
3. conditional branching after real use;
4. persisted definitions only if required;
5. external declarative formats only when authorship outside the implementation language is genuinely needed.

### Events

Potential workflow events:

- `workflow.created`;
- `workflow.started`;
- `workflow.step.started`;
- `workflow.step.completed`;
- `workflow.step.failed`;
- `workflow.paused`;
- `workflow.resumed`;
- `workflow.completed`;
- `workflow.failed`;
- `workflow.cancelled`.

The final vocabulary should avoid duplicating child run events. Workflow events should describe orchestration transitions, while run events remain the source of truth for the work itself.

### Failure and recovery tests

- failure in the first, middle, and final step;
- cancellation between steps;
- cancellation during a child execution;
- daemon or process loss after a step finishes but before state is recorded;
- uncertain external side effect;
- restart with an interrupted active step;
- invalid workflow definition;
- missing step implementation;
- incompatible workflow version;
- duplicated resume request;
- human approval timeout, when introduced.

### Out of scope initially

- unrestricted DAGs;
- arbitrary cyclic workflows;
- visual workflow builder;
- distributed workflow scheduling;
- exactly-once claims;
- BPMN compatibility;
- a plugin marketplace;
- dynamic mutation of running workflow structure;
- multi-agent organisations.

### Exit gate

Workflows are justified only when Aren can express repeated real processes more clearly and reliably than direct composition code, while keeping execution semantics, state transitions, and recovery behaviour inspectable.

---

## Phase 17 — Broader Execution Types

### Primary question

> Which additional forms of work genuinely benefit from Aren’s execution semantics?

Candidates should be introduced one at a time:

- shell processes;
- coding-agent subprocesses;
- HTTP operations;
- MCP tools;
- persistent language workers;
- deterministic jobs;
- human approvals;
- scheduled tasks.

Each candidate must answer:

- What does cancellation mean?
- What constitutes success?
- What output is produced?
- Is retry safe?
- What progress can be observed?
- What state is durable?
- What permissions are required?
- Can it survive runtime or client restart?

This is where the hypothesis that “execution is the primitive” is genuinely tested.

Some execution types may not fit cleanly. Aren should revise the abstraction rather than force false uniformity.

### Exit gate

At least two materially different execution types share useful lifecycle and observation semantics without accumulating special cases that invalidate the core model.

---

## Phase 18 — Distributed and Remote Execution

### Primary question

> Is there enough demonstrated need to execute work across machines or workers?

Possible capabilities include:

- worker registration;
- capability discovery;
- leases and heartbeats;
- remote cancellation;
- artefact transfer;
- placement decisions;
- worker-loss handling;
- duplicate execution protection;
- queueing;
- horizontal scaling.

This phase is speculative and should not be treated as inevitable.

### Exit gate

Remote execution solves a measured operational problem that cannot be addressed adequately by local execution or a single daemon.

---

## 4. Cross-Phase Capability Map

| Capability | Introduced | Deepened later |
|---|---:|---|
| Run identity | Phase 1 | Persistence, daemon, workflows |
| Lifecycle states | Phase 1 | Durable interruption and workflows |
| Cancellation | Phase 1 | Providers, tools, daemon, workflows, remote execution |
| Ordered events | Phase 1 | Streaming, persistence, reconnectable clients |
| Progress | Phase 2 | Models, tools, workflows |
| Failure classification | Phase 1 | Provider, validation, tool, workflow taxonomies |
| Model invocation | Phase 3 | Streaming, tools, context |
| Streaming | Phase 4 | Daemon reconnection and clients |
| Structured output | Phase 5 | Tool schemas and workflow inputs |
| Attempts | Phase 6 | Durable retry and policy budgets |
| Tools | Phases 7–8 | Client-hosted, MCP, remote tools |
| Agent loop | Phase 9 | Context and workflow steps |
| Context management | Phase 10 | Retrieval and memory, if earned |
| Persistence | Phase 11 | Recovery, daemon, workflows |
| Composition | Phase 12 | Workflow foundations |
| Policy | Phase 13 | Organisation-wide governance |
| Daemon | Phase 14 | Remote management and scheduling |
| Multi-language support | Phase 15 | Wider client ecosystem |
| Workflows | Phase 16 | Approvals, schedules, durable orchestration |
| Broader execution types | Phase 17 | Workflow step diversity |
| Distributed operation | Phase 18 | Only if justified |

---

## 5. Explicit Non-Commitments

This roadmap does not commit Aren to implementing every listed phase.

Aren is not yet committed to:

- a generic provider abstraction;
- a workflow graph engine;
- distributed execution;
- a plugin ecosystem;
- a general policy language;
- multi-agent coordination;
- vector-based memory;
- a database-backed architecture;
- durable mid-request resumption;
- exactly-once execution;
- every language SDK;
- every proposed execution type.

Each later phase must be approved from evidence produced by previous phases.

---

## 6. Immediate Planning Boundary

The Phase 1 PRD should cover only:

- run identity;
- lifecycle states and transitions;
- terminal outcomes;
- cancellation;
- lifecycle events;
- event observation semantics;
- minimum failure representation;
- deterministic and race testing;
- diagnostic CLI;
- explicit exclusions.

Phase 2 may be considered only enough to ensure Phase 1 is testable. It must not cause Phase 1 to include:

- provider types;
- model abstractions;
- tool interfaces;
- persistence;
- daemon protocols;
- workflow composition;
- generic execution metadata.

The immediate sequence is:

```text
Approve roadmap direction
        ↓
Write Phase 1 PRD
        ↓
Write Phase 1 lifecycle design
        ↓
Write Phase 1 testing and failure design
        ↓
Implement
        ↓
Use in a small real scenario
        ↓
Review and revise Phase 2
```

---

## 7. Final Principle

Aren should not try to prove at inception that execution is a universal primitive.

It should first prove that Aren can own one execution honestly.

Then it should test whether those semantics survive:

- controlled failure;
- real model requests;
- streaming;
- validation;
- retries;
- tools;
- agent loops;
- context growth;
- persistence;
- composition;
- daemon ownership;
- workflows;
- broader forms of work.

If the abstraction survives that progression, it will have earned its generality.

If it does not, Aren should change the abstraction rather than protect the original vision.
