I’ll trace cancellation from acceptance through work completion and resource release, using the prior authority and outcome decisions to constrain the race analysis.
# Cancellation Goroutine Ownership And Cleanup

> Project: `aren-phase-01-execution-lifecycle`
> Area: Cancellation Goroutine Ownership And Cleanup
> Scope: Cross-sprint cancellation facts, work-context propagation, Aren-owned lifetimes, and resource release
> Evidence boundary: Supplied governed documents and report excerpts; no direct repository inspection or test execution

## Purpose

Define what cancellation acceptance means, how its first cause reaches work, when terminal cancellation becomes truthful, and who releases the resources supporting that protocol. This area consumes the established private lifecycle authority and outcome-resolution policy rather than creating another terminal writer. [IN-03 §§9–17; IN-07, Atomicity model; IN-08, Resolution matrix]

The central conclusion is that **acceptance, propagation, work completion, terminal publication, and resource quiescence are distinct boundaries**. Aren must connect them through explicit ownership and ordering. It must not collapse them into a call to a cancel function, a timeout, an observer-channel close, or a claim that a goroutine has stopped. [IN-03 §§12–19; IN-14, Core Thesis and P1–P4]

This is a semantic contract for sprint reasoning, not a selection of Go packages, public signatures, a watcher implementation, or an executor abstraction. General cleanup hooks, forced goroutine termination, subprocess escalation, persistence, and daemon shutdown remain outside Phase 1. [IN-02, Reasoning Areas; IN-03 §4; IN-04, both implementation waves]

## Inputs Used

Repeated supplied copies are counted once. Repository locations quoted by reports are second-hand evidence, not separately inspected source files. No current Aren implementation, standalone prototype report, test body, or command result was inspected.

| Input ID | Input | Kind | Exact workspace-relative path | Exact portion inspected | How it was used | Limits or applicability warning |
| --- | --- | --- | --- | --- | --- | --- |
| IN-01 | Project index | Governance and requirement | `projects/aren-phase-01-execution-lifecycle/project-index.md` | Project Reasoning Policy; Project Scope; unavailable-dimension table; Prior Decisions; Cross-Sprint Decision Carry-Forward | Scope, review gate, evidence restrictions | Catalog entries are not runtime proof |
| IN-02 | Project-reasoning index | Governance | `projects/aren-phase-01-execution-lifecycle/project-reasoning/index.md` | Reasoning Areas; this area's evidence and source assignments; Excluded Evidence | Dependency order and decision ownership | Routing text is not empirical evidence |
| IN-03 | Phase 1 PRD | Requirement | `projects/aren-phase-01-execution-lifecycle/docs/PRD.md` | Full supplied copy, especially §§4–19 and §§21–25 | Cancellation, lifecycle, observation, and acceptance contract | Normative, not implementation evidence |
| IN-04 | Project roadmap | Requirement and delivery policy | `projects/aren-phase-01-execution-lifecycle/roadmap.md` | Scope Principle; Canonical Project Flow; Cross-Sprint Carry-Forward Rule; both implementation waves; Phase Exit Gate | Sprint boundaries and explicit supersession | Planned commands were not executed |
| IN-05 | Accepted Go decision | Requirement with historical evidence | `projects/aren-phase-01-execution-lifecycle/docs/final-language-decision.md` | §§4, 6–7, 8.2, 11, 14–15 | Owned goroutines, immutable causes and metadata, release obligations, historical defects | Prototype observations and measurements do not verify the new runtime |
| IN-06 | Evidence Assessment And Routing | Dependency reasoning | `projects/aren-phase-01-execution-lifecycle/project-reasoning/areas/00-evidence-map.md` | Supplied Purpose, Inputs Used, Scope, assessment tables, and ending containing Trade-Offs, Evidence, Risks, Self-critique | Source precedence, report restrictions, negative transfer | Excerpt; omitted material was not inspected |
| IN-07 | Lifecycle Authority And Atomic Publication | Dependency reasoning | `projects/aren-phase-01-execution-lifecycle/project-reasoning/areas/01-lifecycle-authority-and-atomic-publication.md` | Supplied Governing questions, Normative constraints, Evidence, Candidate authority models, Atomicity model, and ending containing conclusions and downstream obligations | One authority, committed acceptance before signaling, coherent publication, no callbacks or joins under ownership | Semantic dependency, not executed proof |
| IN-08 | Outcomes Failures And Terminal Resolution | Dependency reasoning | `projects/aren-phase-01-execution-lifecycle/project-reasoning/areas/02-outcomes-failures-and-terminal-resolution.md` | Supplied Governing questions, Terminal fact model, Resolution matrix, Cause-matching rule, and ending containing conclusions, risks, verification, and critique | Strict cause matching, nil-error success, panic distinction, borrowed cause ownership, revalidation | Excerpt; exact implementation remains unproved |
| IN-09 | Selected area template | Reasoning requirement | `projects/aren-phase-01-execution-lifecycle/reasoning/project-wide/03-cancellation-goroutine-ownership-and-cleanup.md` | Full supplied copy | Required comparisons, state model, ownership tree, schedules, and critique | Referenced reasoning README not supplied |
| IN-10 | Aren phased roadmap | Directional project document | `projects/aren-phase-01-execution-lifecycle/docs/phased-roadmap.md` | Development Rules; Phase 1; Phase 2, especially Cleanup semantics; Immediate Planning Boundary | Separate Phase 1 resource release from later executor cleanup policy | Provisional Phase 1 language is subordinate to IN-03 |
| IN-11 | Timeouts and cancellation | Study report | `studies/agent-harness-study/reports/final/07.04-timeouts-and-cancellation.md` | Coverage caveat; Executive Summary; Approach Models; P2, P4–P8, including only the supplied beginning of P8; Open Questions; Evidence Index | Cooperative ceiling, provenance, cleanup ownership, negative process transfer | Report-only; three material sources, not nine; Python and frontend mechanisms |
| IN-12 | Resource locking and isolation | Study report | `studies/agent-harness-study/reports/final/07.05-resource-locking-and-isolation.md` | Executive Summary; Approach Models; Pattern Catalog 4–6 and 9–10; Tradeoffs; Anti-Patterns; Notable Absences; Evidence Index | Owner-local synchronization, lock ordering, idempotent release, work outside locks | Report-only; empty LangGraph checkout; durable and frontend boundaries differ |
| IN-13 | Recovery versus escalation | Study report | `studies/agent-harness-study/reports/final/13.04-recovery-vs-escalation.md` | Executive Summary; Core Thesis; Approach Models, especially OPA and CrewAI; Pattern Catalog P2, P5, P7, P9; supplied Evidence Index | Keep control facts and runtime faults distinct; reject retries and false watchdog completion | Report-only; primarily recovery, agent, and service behavior |
| IN-14 | Go cancellation, ownership, and cleanup | Study report | `studies/aren-go-runtime-study/reports/final/01.02-cancellation-goroutine-ownership-and-cleanup.md` | Executive Summary; Core Thesis; Approach Models; Pattern Catalog P1–P10; supplied Per-Source Notes, Open Questions, Evidence Index | Primary Go ownership, cause, join, repeated-stop, abandonment, and leak evidence | Report-only; source analyses lacked Aren PRD; process shutdown recommendations restricted |
| IN-15 | Go lifecycle ownership and arbitration | Study report | `studies/aren-go-runtime-study/reports/final/01.01-lifecycle-transition-ownership-and-terminal-arbitration.md` | Executive Summary; Approach Models; P1–P8; Tradeoffs; Anti-Patterns; Notable Absences; Evidence Index | Request versus termination, cause-before-signal, terminal funnel, repeated-completion hazards | Report-only; no complete Aren publication or timing solution |
| IN-16 | Go adversarial verification | Study report | `studies/aren-go-runtime-study/reports/final/01.03-adversarial-concurrency-and-failure-verification.md` | Executive Summary; Approach Models; Pattern Catalog 1–5 and 7–11; Anti-Patterns; Notable Absences; Evidence Index | Controlled schedules, independent assertions, counters, negative controls, race and leak gates | Report-only; no complete Aren cancellation/quiescence proof |

The template requires reading `projects/aren-phase-01-execution-lifecycle/reasoning/README.md`. Its content was not supplied, and no read-only workspace tool is available in this turn. It is not an inspected input. Compliance with additional requirements there remains unverified and must be checked before a passing project-reasoning review. [IN-09; IN-06, Purpose]

## Governing questions

| Question | Phase 1 answer | Basis |
| --- | --- | --- |
| Who may request cancellation? | A holder of caller-control authority, or Aren's integration with the supplied parent context. An observer or waiter does not gain control merely by observing. | IN-03 §§9, 12, 16–17 |
| Who accepts it? | The same private per-run authority that serializes lifecycle commitment. Neither the parent watcher nor the work context is an independent acceptance authority. | IN-07, Atomicity model |
| What does acceptance prove? | One request was accepted while running; its first reason and effective cause are fixed; its canonical request occurrence exists; Aren owes propagation of that cause. It does not prove that work observed the signal or stopped. | IN-03 §§12–15; IN-08, Terminal fact model |
| How does the first cause reach work? | Through an Aren-controlled, cause-aware work context after canonical acceptance. Direct parent cancellation must not race around that path. | IN-05 §7; IN-14 P3–P4 |
| Who owns run goroutines? | The run's execution owner owns its invocation carrier and any parent integration or other support carrier. Callers own goroutines that call wait or observation methods. | IN-05 §11.3; IN-14 P1 |
| When may Aren publish `cancelled`? | After the supervised invocation has finished, accepted propagation is settled, and the returned error matches the retained effective cause under the outcome area's rule. | IN-03 §§12–14; IN-08, Cause-matching rule; project conclusion below |
| What if work delays or ignores cancellation? | State remains `running`, with no outcome or finish time. Aren offers no maximum work-stop latency and performs no forced stop. | IN-03 §§12, 18–19 |
| Does a returned outcome prove total quiescence? | It proves work completion and complete terminal facts. Run-owned support must also have an explicit terminal-triggered release path, but ordinary outcome waiting is not automatically a join of every support goroutine. | IN-07, Risks; IN-08, Governing questions; project conclusion below |

## Normative constraints

Sources below resolve to the exact paths in Inputs Used.

| Constraint | Source | Required meaning | Design freedom left |
| --- | --- | --- | --- |
| Five-state lifecycle | IN-03 §§7, 15 | Cancellation acceptance is an occurrence while `running`, not a new state | Private protocol facts may exist without becoming public lifecycle states |
| One cancellation acceptance path | IN-03 §12; IN-05 §7 | Explicit, parent cancellation, and parent deadline expiry share authority and first-write rules | Watcher, callback, or another proved integration |
| First reason retained | IN-03 §12; IN-08, Cause-matching rule | Later requests cannot replace reason, effective cause, or canonical event | Public reason representation and concrete normalization API |
| Three dispositions | IN-03 §12 | `accepted`, `already_requested`, `already_terminal` describe the operation's serialized decision | Exact return type |
| Acceptance before signal | IN-05 §7; IN-07, Minimum coherent fact set | Work cannot receive Aren's cancellation signal before the acceptance event and metadata are committed | Signal placement and internal coordination |
| Work completion before terminality | IN-03 §§10–13 | A stop request, elapsed deadline, or abandoned reader cannot finalize active work | Invocation carrier and completion handoff |
| Fixed terminal interpretation | IN-03 §13; IN-08, Resolution matrix | Nil-error success, accepted-cause-match cancellation, unrelated-error failure, panic failure | Implementation of the existing resolver |
| Cause matching is directional | IN-08, Cause-matching rule | Use `errors.Is(returnedError, acceptedEffectiveCause)`; no blanket `ctx.Err()` fallback | Safe preparation and accepted-fact revalidation outside mutation ownership |
| Already-cancelled parent | IN-03 §12 | Still record creation and start; accept cancellation and signal before invoking work | Startup integration mechanics |
| No blocking work under lifecycle ownership | IN-03 §10; IN-07, Evidence conflicts | No work, user error methods, observer execution, blocking delivery, or joins under the mutation critical section | Internal nonblocking readiness and synchronization primitives |
| Every Aren carrier has an owner | IN-05 §11.3 | Name its creator, stop condition, and join or release path, including error paths | Prefer eliminating carriers over building generic supervision |
| Observation is independent | IN-03 §§16–17 | Absent or abandoned consumers cannot block completion or retain a producer waiting for consumption | Pull cursors or another proved mechanism |
| Resource cleanup is not executor cleanup | IN-03 §4; IN-10, Phase 2 Cleanup semantics | Release Aren-owned support without cleanup hooks, cleanup-result aggregation, or a new terminal policy | Concrete teardown ordering |
| Two independently correct sprints | IN-04, both implementation waves | Sprint 1 proves its owned resources; Sprint 2 adds full cancellation and observer hardening | Exact tests and implementation seams |
| No forced stop or effect rollback | IN-03 §§12, 19 | Cancellation does not guarantee goroutine termination or absence of effects | None within Phase 1 |

## Evidence

### Cancellation protocol comparison

These are report interpretations. They support individual mechanisms and counterexamples, not wholesale adoption of a studied runtime.

| Evidence | Request representation | Cause retention | Acknowledgement | Confirmed stop condition | Repeated request behavior | Aren applicability |
| --- | --- | --- | --- | --- | --- | --- |
| IN-14, Conc Approach Model and P3–P4; IN-15 P3 | Derived context cancellation | Generic context cancellation; causal error separately recorded before cancel-on-error | Task returns; no Aren-style acceptance ledger | Owner's scoped join | Cancel function is repeat-safe; pool `Wait` has a reported double-close hazard | Adapt scoped ownership and cause-before-signal; reject pool waiting as Aren's multi-waiter implementation |
| IN-15 P3, Crush | Cancel function plus active-request retention and dispatch marks | Product-specific cancellation bookkeeping | `Run` returns and records its finish reason | Execution return, not deletion at request time | Marks and guards coordinate several dispatch paths | Adapt active-until-return; reject session queues and high-water marks for one immediately started invocation |
| IN-15 P3; IN-14, Temporal SDK Approach Model and P3 | Workflow cancellation signal; worker/task cause-aware contexts | Typed causes distinguish worker shutdown and other sources | Workflow returns a cancellation error; worker paths have separate stop protocols | Relevant coroutine or worker joins and close ordering | Top-level worker stop has a reported concurrent-close window | Adapt cause retention and separate stop facts; reject durable commands and worker shutdown policy |
| IN-14, Docker Agent Approach Model and P2–P4 | Request cancellation and cause-aware batch contexts | Custom causes distinguish product stop reasons | Tool/stream completion according to its runtime | Guaranteed stream close on the cited path | Uses several once/handshake mechanisms, with residual ownership caveats | Adapt explicit lifetime boundaries; do not adopt its cancellation-based end-reason override |
| IN-11, Pydantic AI Approach Model and P2, P4–P5 | Token and task cancellation | Issued-cancellation accounting, with a documented attribution limitation | Resolution at run boundary | Async task draining; synchronous worker thread may continue after abandonment | Token idempotency and controller finish behavior | Adapt provenance discipline and cooperative ceiling; reject Python accounting and abandonment as Aren completion |
| IN-11 P6; IN-14 P5, process supervisors | OS signals and kill escalation | Process-specific reasons and errors | Signal delivery is not reap | Process wait/reap, sometimes with additional pipe waits | Varies; repeated destruction can rerun hooks | Negative transfer only: a goroutine has no equivalent kill-and-reap protocol |

The request/stop distinction has strong comparative support, but the PRD determines Aren's interpretation. In particular, Docker Agent's reported context recheck before selecting an end reason cannot override a nil-error return in Aren. Conc's first-error behavior and Temporal's close-command priority likewise do not replace the accepted outcome matrix. [IN-03 §13; IN-08, Evidence limits and conflicts; IN-14 P4; IN-15 P1–P3]

### Ownership evidence

| Evidence | Spawn owner | Child tracking | Cancellation source | Join mechanism | Abandonment behavior | Failure or leak tests |
| --- | --- | --- | --- | --- | --- | --- |
| IN-14 P1, Conc | Scoped group | Registration before spawn | Derived context when configured | Group wait after children exit | Forgotten wait can retain work; uncooperative work can prevent return | Report describes cancellation/panic tests but no dedicated leak harness |
| IN-14 P1–P2, P10, Temporal SDK | Worker scope | Separate worker and poller groups | Stop channels and contexts | Ordered joins and single closer | Some bounded/detached paths remain; timeout is not universal quiescence | Report cites `synctest` shutdown tests and an integration leak sweep |
| IN-14 P2, P6–P8, Docker Agent | Runtime and concrete connection supervisors | Scope-specific tracking and adoption/reaping | Request or longer connection lifetime | Close contracts and stop handshakes | Buffered one-result delivery avoids a stranded send but cannot end an uncooperative producer | Issue-linked concurrency tests; no repo-wide leak harness reported |
| IN-14, Ollama Approach Model and P5–P7 | Partial ownership across server, scheduler, and runner | Mixed channels, done signals, and unjoined loops | Request/root cancellation and process signals | Some runner waits; no complete server join tree | Report identifies producer sends that may strand resources when consumers disappear | One goroutine-count regression; disconnect severity partly inferred |
| IN-12, Pattern 10; IN-14, Docker Agent Evidence Index | Resource owner | Collect or identify work while locked | Owner-specific shutdown | Cleanup performed outside lock | Avoids waiting for a child that needs the same lock | Mechanism evidence, not an Aren deadlock proof |
| IN-16, Pattern 4 and Leak philosophy | Test-visible resource owner | Atomic counters, close-time accounting, leak tools | Controlled test completion | Explicit release accounting plus supplementary stack checks | Makes forgotten resources independently observable | Coverage varies; absolute counts and broad ignores can mislead |

The smallest transferable lesson is named ownership plus an observable release path. It is not “use a supervisor type for every goroutine.” For a lifecycle with bounded-count history, removing observer producers may satisfy more requirements with less machinery than adopting a stream supervisor. Exact observation design belongs to the next project area and Sprint 2. [IN-03 §§15–17; IN-05 §11.10; IN-14 P1, P6–P8]

### Load-bearing report anchors

The following paths are **locations cited by inspected reports**, not additional inspected inputs.

| Bounded claim | Inspected report section | Report-cited repository location | Limit |
| --- | --- | --- | --- |
| Register work before spawn and join before scoped completion | IN-14 P1 | `studies/aren-go-runtime-study/sources/conc/waitgroup.go:28-43` | Does not establish Aren's active inspection or multiple-waiter semantics |
| Record cause before cancellation side effects | IN-14 P4; IN-15 P6 | `studies/aren-go-runtime-study/sources/conc/pool/context_pool.go:39-47` | Pool error ordering, not a complete cancellation-acceptance protocol |
| Close ownership differs from repeated waiting | IN-14 Executive Summary and P1 | `studies/aren-go-runtime-study/sources/conc/pool/pool.go:72-82` | Concurrent-close failure is reported from source reasoning, not reproduced here |
| Typed context causes support provenance | IN-14 P3 | `studies/aren-go-runtime-study/sources/temporal-sdk-go/internal/internal_task_handlers.go:2389-2399,2480-2485` | SDK cause categories are not Aren's taxonomy |
| Join producers before closing their shared queue | IN-14 P2 | `studies/aren-go-runtime-study/sources/temporal-sdk-go/internal/internal_worker_base.go:467-478` | Worker queue shutdown, not retained event-history delivery |
| Droppable stop event differs from reliable closure | IN-14 P2 | `studies/aren-go-runtime-study/sources/docker-agent/pkg/runtime/loop.go:187-211,327-337` | Aren instead requires canonical terminal-event retention |
| Abandoned consumers can strand producers | IN-14 P7 | `studies/aren-go-runtime-study/sources/ollama/server/routes.go:649-792` | Report explicitly leaves downstream transport behavior partly unverified |
| Cancellation can be separated from confirmed return | IN-15 P3 | `studies/aren-go-runtime-study/sources/crush/internal/agent/agent.go:1966-1973` | Session dispatch machinery is not needed in Phase 1 |
| Shutdown tests need deterministic schedules and leak evidence | IN-14 P10 | `studies/aren-go-runtime-study/sources/temporal-sdk-go/internal/internal_worker_base_test.go:665-1127`; `studies/aren-go-runtime-study/sources/temporal-sdk-go/test/integration_test.go:109-134` | Reported test presence is not execution evidence in this turn |
| Fixed sleep bias is weak ordering evidence | IN-16, Pattern 2 and Anti-Patterns | `studies/aren-go-runtime-study/sources/conc/pool/context_pool_test.go:176` | Transfer the testing warning, not a conclusion that all time-based tests are invalid |

### Restrictions and negative transfer

IN-11 contains material cancellation evidence from only Agent Framework, OpenHands, and Pydantic AI. Its six empty or near-empty sources establish no upstream behavior. IN-12's Temporal material is from the workflow server, whereas IN-14 and IN-15 discuss the Go SDK; they are not interchangeable evidence for one lifetime boundary. [IN-11, Coverage caveat; IN-12, Sources Studied; IN-14, Sources Studied; IN-06, Evidence quality assessment]

A reliable channel close can be a notification or join mechanism, but a bare close cannot establish Aren's terminal state, cause, timing, or history. Nor are events inherently lossy: Phase 1 explicitly retains canonical lifecycle events. The useful lesson from stream-oriented reports is to separate canonical facts from delivery, not to substitute EOF for a terminal record. [IN-03 §§14–17; IN-14 P2; IN-07 LA-05–LA-07]

Process escalation, service grace periods, polling shutdown loops, retry budgets, and human escalation solve different problems. They do not justify an Aren work timeout, an `abandoned` state, a forced-stop claim, or recovery that re-executes work. [IN-03 §§4, 7, 12, 19; IN-11 P5–P8; IN-13 P5, P7, P9; IN-14 P5]

## Cancellation state model

### Facts and visibility

Cancellation is modeled as monotonic facts attached to a run, not as another lifecycle state machine exposed to callers. Internal propagation and release tracking are permitted only where the realized mechanism needs them. [IN-03 §§7, 12, 15; IN-07, Atomicity model]

| Fact | Owner | First-write rule | Observable to | Effect |
| --- | --- | --- | --- | --- |
| Request received | Explicit caller path or parent integration | Each invocation may present a request; wall-clock arrival does not assign priority | Requesting path; not necessarily retained history | Submit to the common acceptance authority |
| Request accepted | Private lifecycle authority | Only the first eligible request while running | Controller disposition and canonical request event | Fix reason and non-nil effective cause; create a propagation obligation; state remains running |
| Context cancelled | Aren's work-context propagation owner | Active-run cancellation uses only the accepted effective cause | Work through its context | Cooperative signal; not acknowledgement or completion |
| Work observed the signal | Work | Not generally observable by Aren | Work; controlled test instrumentation | May influence work but creates no lifecycle event or terminal permission |
| Work acknowledged cause | Work supplies error; Aren evaluates it | One completed return checked against the stable accepted cause | Resolver and terminal diagnostics | Cause match makes cancellation eligible; does not independently publish it |
| Work invocation finished | Aren's invocation boundary | One normal return or captured escaping panic | Completion owner | Enables terminal resolution after required propagation settlement |
| Terminal cancellation committed | Private lifecycle authority | At most one terminal commitment, only from running | Waiters, state/outcome readers, retained history | Publish complete cancelled outcome and event |
| Run support released | Each named support owner | Release is single-owner or explicitly idempotent | Internal join/accounting tests; public surface only if separately specified | Removes active watchers, registrations, timers, and blocked support paths |

A request may be rejected without acceptance. Acceptance and signaling may occur without acknowledgement or work completion. Work may observe cancellation and then return nil, an unrelated error, or panic. Work can finish without ever observing cancellation. A terminal cancelled outcome establishes the contractual matching rule, not authenticated proof of why user code returned that error. [IN-03 §§12–13; IN-08, Cause-matching rule and Self-critique]

### Acceptance and disposition order

The common acceptance operation decides against the serialized run facts:

1. If a terminal outcome is committed, return `already_terminal`.
2. Otherwise, if cancellation is already accepted, return `already_requested`.
3. Otherwise, while running, retain the first reason and effective cause, append exactly one request event, and return `accepted` after fulfilling the operation's propagation obligation.

This ordering makes `already_terminal` the disposition after completion even if cancellation was accepted earlier. “Stable disposition” means one operation's answer is final and consistent with its decision point, not that every later call must repeat the same answer forever. A repeated request neither appends an event nor allocates a sequence. [IN-03 §12; IN-07, Minimum coherent fact set; project interpretation]

“First” means first accepted by the common authority. Concurrent explicit and parent requests have no source priority and no timestamp-based ordering. A cancellation call can return `accepted` after another goroutine has observed terminal success: the request was accepted earlier, its propagation completed, and work then succeeded before the requesting goroutine resumed. The disposition is not a fresh state snapshot at response receipt. [IN-03 §§12–13; IN-07 LA-06; project interpretation]

### Cause ownership

The effective cause accepted by Aren must be the cause used to signal the work context. A human-readable reason, if represented separately, must remain consistent with that cause and must not replace it as the matching key. When an explicit request supplies no distinct cause, its normalization must yield a non-nil standard cancellation cause. Parent requests retain the supplied parent's cancellation cause, including custom causes and deadline causes. Exact public parameters remain a sprint decision. [IN-03 §12; IN-08, Cause-matching rule; IN-14 P3]

Original custom errors and causes follow the outcome area's read-only borrowing contract. Aren owns stable metadata and diagnostic projections; it cannot universally clone or freeze arbitrary error implementations. Matching or formatting may invoke user methods, so those operations must not occur under lifecycle mutation ownership. [IN-08, Terminal fact model and Risks; IN-07, Risks]

The outcome rule is preserved unchanged:

```text
normal return with nil error:
    succeeded

normal return with non-nil error, accepted cancellation,
and errors.Is(returnedError, acceptedEffectiveCause):
    cancelled

other normal return with non-nil error:
    failed, work/returned_error

escaping work panic:
    failed, work/panic
```

Returning only `ctx.Err()` is not sufficient when it does not match a distinct accepted custom cause. Work can return or wrap `context.Cause(ctx)` to acknowledge that cause. A joined error containing the accepted cause follows the outcome area's matching rule and retains the complete acknowledgement error diagnostically. [IN-08, Resolution matrix and Cause-matching rule]

### Acceptance-to-signal ordering

The project selects the following ordering constraint:

```text
running committed
    -> acceptance metadata and request event committed
    -> accepted cause propagated to work context
    -> terminal commitment, if work has finished
    -> terminal-triggered support release
```

There may be an observable acceptance-to-signal interval. Acceptance records Aren's decision and propagation obligation; it does not certify that work already observed `Done`. However, an accepted request must not be stranded in that interval. The accepting operation must complete propagation before returning `accepted`, and terminal publication must not overtake an outstanding accepted propagation. This closes the race in which terminal teardown would otherwise substitute a generic context-release cause before the accepted cause reaches work. [IN-05 §7; IN-07, Minimum coherent fact set; IN-08, Terminal fact model; project conclusion]

This is not permission to wait under the mutation critical section. Sprint reasoning must provide either a safe ordered signal step or explicit coordination outside that section. The implementation must not introduce a cycle in which propagation waits for terminal publication while terminal publication waits for propagation. An `already_requested` disposition proves the first acceptance exists; it is not a join of another caller's propagation operation. [IN-07, Evidence conflicts; IN-12, Pattern 6; project conclusion]

### Parent integration and startup

A direct child context with uncontrolled automatic parent propagation is insufficient if it can signal work before Aren records acceptance. Adding a watcher beside such a context does not repair the causal inversion. Parent values and supplied deadline information must remain accounted for by the chosen context design, but parent `Done`, `Err`, and cause visibility must not bypass the Aren-owned acceptance path. [IN-03 §§10–12; IN-05 §7; IN-07, Risks]

A parent deadline is a source of cancellation, not an Aren budget or a new terminal state. It is routed through the same acceptance path, and its accepted cause governs matching. No independent run deadline timer may compete to replace that cause. A context-detachment mechanism, if selected, must explicitly preserve the intended value/deadline behavior rather than silently turning parent-derived work into background-owned work. Exact Go construction requires Sprint 2 verification. [IN-03 §§4, 12–13; IN-05 §11.3; IN-08, Cause-matching rule]

Startup must distinguish two cases:

- If the parent is already cancelled before invocation, creation and start are committed first, then acceptance and propagation complete before work is invoked. Skipping work or taking `created -> cancelled` would violate the contract.
- If the parent becomes cancelled concurrently with startup, the integration must not miss it. Registration and the pre-invocation check need a race-safe handoff. Its acceptance can race with later work completion; parent wall-clock cancellation alone does not reserve a terminal outcome.

No separate dispatch queue, pending-run ledger, or session-wide cancellation mark is earned by this requirement. [IN-03 §§7, 10, 12–13; IN-15 P4; project conclusion]

## Ownership tree

The tree is conceptual. It names ownership even where the implementation needs no goroutine.

```text
embedding caller
|-- supplied parent context and any parent timer
|-- caller-owned wait invocations and wait contexts
|-- caller-owned observers and observation contexts
`-- Aren run execution ownership
    |-- one supervised work invocation
    |-- work context and its cancellation/release capability
    |-- parent-cancellation integration, if a carrier is required
    |-- lifecycle mutation carrier, only if an ownership loop is selected
    |-- terminal readiness mechanism
    `-- observer support, only if the selected design creates resources

reachable run object
`-- immutable terminal outcome and retained canonical history
```

Execution lifetime and object lifetime differ. A terminal run retains its outcome and history while reachable, but does not need an active parent watcher or event producer merely to preserve those records. An observer retaining a run intentionally retains its history; a parent registration retaining a completed run after clients release it is a different retention problem. [IN-03 §16; IN-05 §§7, 11.3; IN-14 P1, P7]

| Node | Created by | Cancelled by | Completion observed by | Joined or released by | Can outlive run? |
| --- | --- | --- | --- | --- | --- |
| Supplied parent context and timer | Embedding caller | Caller or parent deadline | Aren's parent integration | Caller retains responsibility for its cancel function/timer | Yes; Aren must not cancel its parent |
| Run execution owner | Aren start path | Cancellation does not terminate the owner; it initiates the protocol | Lifecycle and internal completion accounting | Owner's terminal/release path | Only a finite post-publication release epilogue; no ongoing work loop |
| Work invocation carrier | Run execution owner | Accepted work-context signal, cooperatively | Invocation boundary captures return or panic | Invocation owner; termination accounted for independently of whether anyone waits | User invocation cannot continue after terminal publication; carrier may finish its bookkeeping epilogue |
| Work context | Run execution owner | Accepted cause during activity; resource release after work completion | Work, resolver's captured facts, internal teardown | Run execution owner | May be retained by user code, but Aren must remove unnecessary registrations and references |
| Parent watcher or callback registration | Run execution owner | Terminal release; first acceptance can also make further parent observation unnecessary | Registration/exit handshake or concrete stop result | Named run support owner, outside mutation ownership | Only finite release epilogue; never until the parent eventually cancels |
| Lifecycle ownership loop, if selected | Aren start path | Terminal shutdown protocol | Private loop-exit signal | Run support owner or equivalent explicit release protocol | No continuing command loop after terminal release |
| Terminal readiness mechanism | Lifecycle owner | Not caller-cancellable as a shared fact | All waiters | Single publication owner makes it ready once; retained with run | Yes as inert completed state; no goroutine required |
| Wait invocation | Calling goroutine | Its own wait context, if supported | Calling goroutine | Calling goroutine; Aren must not leave a per-wait helper behind | May overlap completion briefly; abandonment does not cancel run |
| Observer cursor | Observer | Observation context or explicit close, if supported | Observer | Observer-local release; run must not require drainage | Yes; post-completion replay is required |
| Observer producer or registration, if introduced | Named observation owner | Observer abandonment or terminal delivery completion | Concrete exit/stop handshake | Observation owner, without run waiting for consumption | No active producer stranded by unread output |
| Work-spawned goroutines | User work | User work's own protocol | User work or embedding caller | User code | Possibly; Aren does not discover or join arbitrary descendants |
| Outcome and canonical history | Lifecycle owner | Not cancellation-controlled | Readers | Reachability and collection | Yes; retained data is required, not an active-resource leak |

Work-spawned goroutines and user-installed context callbacks are not silently adopted into an Aren execution tree. The terminal guarantee concerns the supplied invocation, including its ordinary defer unwinding, not all work that user code may have launched elsewhere. This boundary must appear in the promoted contract so “work stopped” does not imply process isolation or descendant supervision. [IN-03 §§10, 19; IN-08, Governing questions; IN-05 §11.3]

## Race analysis

Each schedule assumes the central outcome policy and one publication authority. Ordering is established by barriers or internal test seams, not by sleep duration. [IN-07 LA-01–LA-07; IN-08, Resolution matrix; IN-16, Pattern 2]

| Schedule | Observable facts | Terminal interpretation | Goroutine and resource cleanup |
| --- | --- | --- | --- |
| 1. Completion immediately before explicit cancellation | If terminal commitment precedes acceptance, the request returns `already_terminal`, with no request event. Physical work return alone does not close acceptance. | Preserve committed outcome. If only work return occurred first, acceptance may still precede commitment, and the resolver uses those accepted facts. | Terminal release stops parent integration; late request starts no new carrier or timer. |
| 2. Accepted cancellation immediately before successful return | One request event; accepted cause reaches context; work returns nil error. | `succeeded`, with request retained only as history, not terminal-cancellation payload. | Settle propagation, commit success, release support. No wait for an acknowledgement that work did not provide. |
| 3. Parent and explicit cancellation together | Exactly one accepted cause. Other operation sees `already_requested` while active or `already_terminal` after completion. | Apply return/error/panic policy against the first accepted cause only. | Losing path releases without signaling its own cause; parent integration cannot remain registered unnecessarily. |
| 4. Accepted cancellation followed by unrelated error | Request event precedes failure; returned error does not match accepted cause. | `failed`, `work/returned_error`; preserve original error. | Same support release as success or cancellation; no retry or cancellation override. |
| 5. Work ignores cancellation indefinitely | Request recorded and context signaled; state remains running, outcome and finish time absent, shared wait remains pending. | No terminal candidate. | Invocation remains active by the explicit cooperative limit. Redundant parent observation can stop after acceptance; no accumulating observers, timers, or per-request workers are permitted. Tests later release work. |
| 6. Cancellation arrives after terminal commitment | `already_terminal`; history, timing, and outcome unchanged. | Existing terminal remains authoritative. | Do not restart a watcher, signal a newly created context, or close readiness again. |
| 7. Waiter abandons its own wait context while run continues | Local wait returns its own cancellation result, if such a wait API is offered; other waiters remain pending; run has no new request event. | Run interpretation is unaffected. | Caller exits wait; any internal registration is removed. No helper stays blocked delivering an outcome to the departed caller. |
| 8. Run finishes while parent watcher starts or stops | Terminal history is complete regardless of watcher scheduling. A delayed parent callback must recheck terminality. | Either accepted-before-terminal facts apply, or the late request is rejected. No post-terminal request event. | Registration must be prevented, stopped, or tracked through exit even when terminal release races its installation. No join while holding the lifecycle lock. |
| 9. Parent is already cancelled before invocation | Creation, start, request, and signal all precede work entry. Work runs once with cancelled context. | Nil return succeeds; matching error cancels; unrelated error or panic fails. | No dormant watcher is needed solely to rediscover the known cancellation. Remaining support releases on completion. |
| 10. Acceptance committed, propagation paused, work returns | Readers may see running plus request; accepted caller has not completed its propagation obligation. | Terminal publication waits for accepted propagation settlement, then uses the ordinary policy. | Coordination must remain outside mutation ownership and must not wait on work or observer cooperation. |
| 11. Cause matching pauses while first acceptance occurs | Preliminary “not accepted” classification is stale. | Revalidate accepted facts before commitment and, when necessary, match against the newly accepted cause outside ownership. | No classification helper is abandoned; no asynchronous worker is created merely to time out a user error method. |
| 12. Accepted cancellation followed by work panic or blocked defer | Request does not establish completion; defer execution remains part of invocation. | Escaping panic becomes `work/panic` only after the supervised unwind reaches its boundary. A blocked defer prevents terminality. | Release follows completed capture, not the initial panic occurrence or context signal. |

There is no deterministic source winner for genuinely concurrent parent and explicit requests. Determinism applies to interpretation of the committed facts, while controlled tests separately force each acceptance ordering. A stress test that accepts either winner must still assert that history, context cause, and outcome all agree on the same winner. [IN-03 §13; IN-08, Verification obligations; IN-16, Pattern 9]

If a context-aware wait or cursor API offers local cancellation, simultaneous local cancellation and terminal readiness must have a documented return rule selected by its owning sprint area. Regardless of that rule, local abandonment cannot mutate the run, consume its unique outcome, or create a terminal-cancellation event. This does not require adding a new public wait-context API in Sprint 1. [IN-03 §§9, 16–17; IN-04, Sprint 1 Deferred; project boundary]

## Cleanup boundary

### What Phase 1 cleans up

Phase 1 cleanup consists of retiring the machinery Aren created to supervise the invocation:

- Stop or unregister parent-cancellation integration.
- Complete any accepted propagation obligation before terminal publication.
- Release child-context resources without manufacturing a new accepted request.
- Stop Aren-created timers, if the realized implementation genuinely needs any.
- End an ownership loop, if one exists, without stranding queued caller replies.
- Release observer registrations or producers without requiring consumer drainage.
- Remove unnecessary references to work closures, parent contexts, and temporary completion storage while retaining the required outcome and history.

These are ownership obligations, not user hooks or a general cleanup phase. They apply on success, returned error, captured work panic, cancellation, and any containable internal-fault path. [IN-05 §§7, 11.3; IN-03 §§4, 16; IN-07 LA-08; IN-08 OF-09]

Releasing a child context after work ends is resource management, not evidence of an accepted cancellation request. It must not append `run.cancellation_requested`, alter the retained first cause, or reclassify the completed invocation. Classification uses captured work facts and accepted cancellation facts, never a generic cancellation artifact introduced by teardown. [IN-03 §§12–15; IN-08, Terminal fact model; IN-14 P4]

### Publication versus quiescence

The project distinguishes three boundaries:

| Boundary | Required truth |
| --- | --- |
| Invocation finished | The supplied function has returned or unwound through its supervised capture boundary; no more code from that invocation is running |
| Terminal publication | Accepted propagation is settled; state, event, timing, outcome, and readiness are coherent and final |
| Run-support quiescence | Aren-owned support has exited or its registration has been definitively released; no blocked producer or watcher keeps the completed execution active |

Ordinary waiting returns the terminal outcome after the second boundary. It must not be described as an unconditional join of every support goroutine unless the realized implementation supplies that stronger guarantee. Terminal publication must initiate all remaining support release, with only a finite internal epilogue allowed. That epilogue may depend on scheduling and short synchronization, but not on future parent cancellation, user work, observer consumption, or another call to `Wait`. Private joins or equivalent accounting must make the third boundary testable. [IN-03 §17; IN-05 §11.3; IN-07, Risks; IN-08, Governing questions; project conclusion]

This avoids a self-join trap: a terminal owner cannot wait for its own goroutine to exit, and a parent callback cannot be joined while blocked on the lock held by its joiner. It also avoids concealing leaks behind the phrase “eventually cleaned up.” Every remaining carrier needs a concrete exit path and a completion observation used by tests. No additional public quiescence API is required by this decision. [IN-12, Pattern 6; IN-14 P1–P2; project conclusion]

“Finite epilogue” is a structural guarantee under normal scheduling and supported, terminating context/error behavior, not a hard real-time deadline. The PRD provides no maximum cancellation latency and no mechanism to enforce one against arbitrary in-process code. Test watchdogs diagnose violations; they do not convert a blocked invocation into a terminal run. [IN-03 §§4, 12, 19; IN-08, Self-critique; IN-16, Pattern 2]

### What Phase 1 does not clean up

Aren does not invoke arbitrary executor cleanup callbacks, merge cleanup errors into work outcomes, run detached cleanup with a grace budget, compensate effects, kill goroutines, or adopt work-created background activity. Those would introduce the later controlled-executor or process-supervision contract. User defers remain part of ordinary work execution and need no Aren hook API. [IN-03 §§4, 10, 25; IN-10, Phase 2 Cleanup semantics]

An already-cancelled work context must not make Aren skip its own resource release. Conversely, introducing a fresh background cleanup context does not solve local join ordering and could conceal an unowned lifetime. A concrete new blocking resource must earn its lifetime design; Phase 1 must not build a generic cleanup-context factory in anticipation. [IN-05 §§11.3, 11.10; IN-11 P8; IN-14, Approach Models]

## Project conclusions

| Conclusion | Evidence basis | Sprint 1 obligation | Sprint 2 decision required | Reopen trigger |
| --- | --- | --- | --- | --- |
| CG-01: One authority accepts explicit and parent requests; acceptance is not terminality | IN-03 §§12–15; IN-07 LA-01–LA-04; IN-15 P3 | Preserve private lifecycle ownership and truthful completion | Implement common acceptance and disposition path | Any source signals or records cancellation outside authority |
| CG-02: First accepted reason and effective cause remain fixed; terminal disposition has precedence after completion | IN-03 §12; IN-08, Cause-matching rule | Preserve outcome/cause ownership boundaries | Select concrete reason API and normalization; test both source orderings | Cause replacement, duplicate request, or inconsistent disposition |
| CG-03: Commit acceptance before signaling; settle accepted propagation before returning `accepted` or publishing terminal truth | IN-05 §7; IN-07, Atomicity model; IN-14 P4 | Avoid a completion design that requires callbacks or blocking joins under ownership | Prove signal coordination and terminal-race handling | Signal bypass, stranded propagation, teardown cause overwrite, or dependency cycle |
| CG-04: Already-cancelled parents still produce creation/start and one work invocation receiving cancellation | IN-03 §§7, 10, 12 | Document deferred parent behavior honestly | Select race-safe registration/pre-invocation handoff | Missed startup cancellation or skipped invocation |
| CG-05: Preserve the outcome area's strict cause-based policy | IN-03 §13; IN-08, Resolution matrix | Complete success/error/panic resolver | Integrate accepted facts, safe matching, and revalidation | Equal committed facts produce different outcomes |
| CG-06: Work can remain running indefinitely; no force-stop or abandon-as-terminal policy | IN-03 §§12, 19; IN-11 P5; IN-14, Conc Approach Model | Publish only after invocation completion | Delayed/ignored cancellation cases and truthful CLI behavior | Any terminal result while invocation or its defer remains active |
| CG-07: Every Aren-owned carrier and registration has a named stop and release path | IN-05 §11.3; IN-14 P1, P7–P8 | Inventory and prove all Sprint 1 carriers | Extend inventory for parent and observation integration | Untracked spawn, retained completed run, or consumer-dependent exit |
| CG-08: Terminal outcome waiting and support quiescence are distinct; all post-terminal release is finite and independently verifiable | IN-07, Risks; IN-08, Governing questions; IN-14 P1–P2 | Define actual publication and carrier-exit boundaries | Pin joins/stop results and eliminate watcher startup/teardown gaps | `Wait` overclaims joins, or support waits indefinitely after terminality |
| CG-09: Waiter/observer abandonment affects only that observation lifetime | IN-03 §§9, 16–17; IN-14 P6–P7 | Multiple waiters without exclusive consumption or producer dependence | Choose observation abandonment and cursor resources | Abandonment cancels execution or strands a producer |
| CG-10: Cleanup means releasing Aren resources, not adding executor hooks or cleanup outcome policy | IN-03 §4; IN-10, Phase 2 | Release on all implemented exits | Extend local release only | A proposed cleanup abstraction lacks a current Phase 1 requirement |
| CG-11: Race, semantic, and leak evidence are complementary and must prove their instruments can fail | IN-03 §21; IN-05 §11.6; IN-16 | Fresh race-enabled ordinary lifecycle/resource tests | Controlled collisions, repeated stress, retention checks, negative controls | Green tests rely on sleeps, ignored leaks, mocks, or historical passes |

Sprint 2 must classify the realized Sprint 1 decisions as preserved, extended, superseded, or unaffected. A changed mechanism requires named evidence and recorded impact; a fixed PRD guarantee cannot be weakened by an implementation preference. No conclusion here grants either sprint acceptance or a passing project review. [IN-04, Cross-Sprint Carry-Forward Rule; IN-01, Project Reasoning Policy]

## Trade-Offs

**Watcher goroutine versus callback integration.** A dedicated parent watcher makes its select conditions and exit signal explicit, but adds a carrier per active run and a startup/teardown handshake. Callback integration can avoid an idle watcher, but callback execution, deregistration, concurrent start, and completion still have ownership costs. A stop request is not automatically a callback join. Sprint 2 should choose the smaller implementation after examining the supported Go toolchain's concrete semantics, not assume either option is free. [IN-05 §§11.3, 11.10; IN-14 P1, P8]

**Polling versus event-driven cancellation.** Periodic polling adds latency, timers, and a tunable interval without satisfying a present requirement. Polling only at completion fails to propagate cancellation to blocked work. A direct pre-invocation parent check is useful for startup, but is not a substitute for active parent integration. Polling is therefore rejected as the sole parent-cancellation mechanism. [IN-03 §12; IN-15, Anti-Patterns; project conclusion]

**Cause-aware context versus a simple cancel function.** Simple cancellation communicates “stop” but loses the distinction between a custom accepted cause and a generic sentinel. Cause-aware propagation supports the selected matching policy, with a documented usability cost for work returning only `ctx.Err()`. Separate provenance metadata without consistent context propagation would create two conflicting truths. [IN-08, Cause-matching rule and Trade-Offs; IN-14 P3]

**Signal ordering versus additional coordination.** Requiring accepted propagation to settle before terminal publication adds an ordering obligation, but prevents a completed run from retaining an accepted reason that never reached its work context. The extra mechanism should be minimal and private. It must not become another public state, independent terminal arbiter, or per-request goroutine. [IN-05 §7; IN-07, Atomicity model; CG-03]

**Blocking joins versus truthful return latency.** Joining stopped internal support can give a strong quiescence guarantee, but joins under lifecycle ownership can deadlock, and waiting for arbitrary work after a deadline cannot provide a bounded return. The selected contract separates outcome readiness from the finite support-release epilogue while requiring independent release proof. A stronger join-before-wait-return guarantee may be implemented if inexpensive, but is not presumed. [IN-03 §§12, 17, 19; IN-12, Pattern 6; IN-14 P1–P2]

**Leak prevention versus unnecessary lifetime machinery.** A generic supervisor, adoption ledger, resource registry, or cleanup timeout framework would exceed the small run's needs. Conversely, “just a goroutine” is not an ownership argument. Prefer no carrier for waiting and observation where the chosen API permits it; name and test every carrier that remains. [IN-05 §§11.3, 11.10; IN-14 P6–P8]

**History retention versus memory reclamation.** Retaining canonical history and outcome while the run is reachable is required. Retaining an active parent registration after completion is not. Lifecycle event count is bounded by the vocabulary and first-acceptance rule, but arbitrary result graphs, causes, and diagnostics are not thereby byte-bounded. Measurements must separate these categories rather than label all retained memory a leak or all memory bounded. [IN-03 §§15–17; IN-05 §§11.7–11.8; IN-08, Risks]

## Risks

| Risk | Failure mode | Required control |
| --- | --- | --- |
| Cancellation treated as completion | Request or deadline publishes `cancelled` while work continues | Gate every terminal path on completed invocation facts; hold work after signal in tests. [IN-03 §12] |
| Parent bypass | Child `Done` or cause becomes visible before canonical acceptance | Audit complete context chain, not just watcher code; test work's first signal observation against retained history. [IN-05 §7] |
| Stranded accepted propagation | Acceptance returns or terminal publishes while its signal operation is lost | Explicit propagation owner and settlement ordering; test the acceptance-to-signal gap. [CG-03] |
| Context cause mismatch | Parent, explicit request, and teardown install different causes | One accepted effective cause; no independent cancellation writer; strict matching tests. [IN-08, Cause-matching rule] |
| Unknown goroutine owner | Helper survives because nobody observes its exit | Complete carrier inventory, including callbacks and library-created registrations. [IN-05 §11.3] |
| Join under lifecycle lock | Watcher needs the lock held by its joiner | Stop/join outside ownership; test watcher paused before acceptance. [IN-12, Pattern 6] |
| Cleanup uses already-cancelled work context | Release is skipped or returns early without freeing support | Resource release independent of work cooperation; no cleanup callback framework. [IN-05 §11.3; CG-10] |
| Double close or send-after-close | Repeated cancellation/waiting races a notification closer | One close owner; waiting never owns closure; producer shutdown protocol where channels exist. [IN-14, Executive Summary and P2] |
| Abandoned consumer strands producer | Wait or observer disappears while producer sends | Avoid producer carriers where possible; otherwise prove an abandonment exit and terminal delivery contract. [IN-14 P6–P7] |
| Late registration retains completed run | Terminal release runs before watcher stop handle is installed | Registration/terminal handshake and explicit stop-or-exit accounting. [CG-07–CG-08] |
| Cleanup artifact changes outcome | Generic context release makes a success look cancelled | Resolve from captured work and accepted facts, not teardown `Err`. [IN-08, Terminal fact model] |
| Internal fault falsely stops active work | Invariant fallback fabricates a terminal state | Preserve outcome area's containment preconditions; do not use fault handling as force-stop. [IN-08 OF-09] |
| Work descendants are mistaken for Aren children | Terminal wording implies all user goroutines stopped | Document supervised-call boundary and user ownership of spawned activity. [IN-03 §19] |
| Process escalation leaks into scope | Grace timer manufactures terminality or adds kill infrastructure | Explicit rejection of process and abandoned-state transfer. [IN-03 §4; IN-14 P5] |
| Leak tests excuse active support | Broad stack ignores hide watcher/observer leaks | Controlled work release, targeted accounting, narrow justified exclusions. [IN-16, Anti-Patterns] |
| Governance compliance assumed | Missing README standard is overlooked | Verify the named reasoning README before approval. [IN-09] |

## Verification obligations

These are required instruments and assertions, not executed results. The verification area must connect them to the realized runtime and an independent lifecycle oracle. A fake that already enforces terminal order is insufficient. [IN-03 §21; IN-05 §11.6; IN-15, Anti-Patterns]

| Claim | Controlled schedule | Observable assertion | Leak or quiescence assertion | Negative control |
| --- | --- | --- | --- | --- |
| Acceptance precedes work-visible signal | Work waits on context; trigger explicit cancellation, then parent cancellation in a separate case | On signal observation, a fresh coherent history contains the request and matching first reason | Acceptance operation and parent integration reach their recorded release boundaries | Signal context before appending request |
| Acceptance propagation cannot be overtaken | Pause after acceptance commitment but before signal; let work return nil | State remains nonterminal until propagation settles; eventual outcome succeeds | No helper remains waiting on a terminal/propagation cycle | Allow terminal teardown to run first |
| First request is stable | Force explicit-first, parent-first, then true collision | One request event; disposition, accepted metadata, context cause, and final diagnostics agree | Losing paths create no surviving support | Replace cause on repeated request |
| Already-cancelled parent works | Cancel parent before start; instrument work entry | Creation/start/request precede invocation; work runs once with cancelled context | No dormant watcher waits on an already-settled source | Rely only on asynchronously scheduled watcher |
| Startup handoff misses no parent signal | Cancel parent between initial check and integration installation | Event and signal occur through common acceptance unless terminal already won | Registration stops on immediate completion | Remove post-registration reconciliation |
| Nil-error success survives cancellation | Work observes signal, then waits for release and returns nil | Request followed by success; no terminal cancellation payload | All run support releases | Cancellation-priority terminal override |
| Cause matching remains strict | Return direct/wrapped custom cause, coarse sentinel, same-text error, unrelated timeout, and joined error | Outcomes match IN-08's matrix; causes remain inspectable | No matching helper survives; no user method runs under mutation ownership | Add blanket `ctx.Err()` fallback |
| Delayed acknowledgement is not completion | Hold work after observing cancellation | Running state, absent finish/outcome, waiters still pending | Parent integration does not multiply or remain needlessly active after acceptance | Publish terminal on `Done` |
| Ignored cancellation remains truthful | Work ignores context and blocks on a test-owned release channel | Repeated requests add no events; no terminality before release | Release work even on assertion failure, then join/account for all Aren support | Timeout and mark run cancelled |
| Unrelated error and panic remain failures | Accept cancellation, then release to error or panic | Correct work failure kind; exactly one terminal event; no result | Same release obligations as normal success | Classify all post-cancel exits as cancelled |
| Terminal requests are inert | Complete run, then issue concurrent cancellation requests | All report `already_terminal`; history/timing/outcome unchanged | No watcher or timer restarts | Close readiness or append request on each call |
| Multiple waiters only observe | Many callers wait before and after completion | Same logical outcome; no exclusive consumption or duplicate finalization | No Aren per-wait goroutine remains | Reuse a one-consumer outcome channel or close in `Wait` |
| Local wait abandonment is isolated | Cancel one wait context while another wait and work remain active, if API supports it | No run request event; other wait remains pending | Abandoned wait leaves no registration/helper | Derive work lifetime from wait context |
| Observer abandonment cannot block run | Stop reading before terminal; run a second fast observer | Work/cancellation/waiters continue; retained terminal remains available | Producer/registration exits without drainage | Unbuffered mandatory delivery from commit path |
| Watcher startup and terminal release compose | Pause installation, finish run, resume installation | No post-terminal mutation | Stop handle installed late is still released; no parent-retained completed run | Store stop handle after teardown without reconciliation |
| Watcher stop does not deadlock | Pause watcher before lifecycle acquisition; finish run | Terminal facts and concurrent reads remain coherent | Owner releases lock before joining; watcher exits | Join watcher while holding mutation lock |
| Teardown does not create cancellation history | Complete uncancelled success/error/panic; inspect after support release | Still exactly creation/start/terminal; original outcome unchanged | Child-context resources released | Route terminal context release through acceptance |
| Completed runs release parent retention | Keep parent alive across repeated completed runs; drop unneeded run handles | No late parent cancellation changes completed outcomes | Registrations and owned stacks return to baseline; targeted retention evidence distinguishes required snapshots | Leave callback registered until parent ends |
| Instrumentation detects real omissions | Seed one missing exit/stop/close-order defect at a time | Specific invariant fails rather than only global timeout | Leak/accounting instrument identifies the deliberately retained support | Empty helper or unrelated counter that always passes |

### Evidence execution requirements

Sprint 1 must already prove ordinary invocation completion, panic capture, shared waiting, complete publication, and release of its actual carriers under `go test -race`. Full caller/parent cancellation and observer behavior must not be claimed before Sprint 2 implements them. [IN-04, Sprint 1 Acceptance and Deferred]

Sprint 2 must combine deterministic schedules with repeated completion-cancellation collisions, concurrent inspection, observer churn, and parent lifetime tests. Every randomized test must retain its seed and the observed committed trace; a seed alone does not reproduce Go scheduler interleavings. Controlled traces must force both sides of important races independently. [IN-16, Patterns 7 and 9; IN-03 §21]

Quiescence evidence should start with explicit exit/stop observations and resource accounting. Supplementary goroutine-stack and retention checks catch resources omitted from the inventory. Raw `runtime.NumGoroutine` differences, garbage collection timing, or finalizers alone are not decisive. Test-owned blocked work must always be released before asserting the final resource baseline. [IN-14 P10; IN-16, Leak philosophy and Anti-Patterns; IN-06, Risks]

Exact stress iterations, test names, toolchain-dependent scheduling tools, leak tooling, and measurement commands belong to sprint reasoning. Recorded evidence must identify revision, worktree status, toolchain, workload, repetitions, and relevant spread. Historical prototype memory figures and dated green commands cannot satisfy these obligations. [IN-05 §§4, 8.2, 11.6–11.8, 15]

## Self-critique

**Which goroutine lacks a named owner or join path?** None is intentionally unnamed in the conceptual inventory. That is not proof that the implementation will match it: context callbacks, propagation helpers, and an optional ownership loop are the most likely hidden carriers. Sprint review must inspect actual spawn and registration sites and connect each to an observable release path. Creating a helper solely to wait with a timeout is particularly suspect because the helper can outlive the timed-out caller. [IN-05 §11.3; IN-14 P1, P6–P8]

**Could Aren report `cancelled` while work is still executing?** Not under this contract. The invocation must finish, accepted propagation must settle, and the returned error must match the accepted cause. Observing context cancellation, a test watchdog firing, or an invariant fault is insufficient. The important wording limit is that arbitrary work-spawned goroutines are outside the supervised invocation and remain user-owned. [IN-03 §§12–13, 19; IN-08 OF-09; CG-06]

**What happens if cancellation propagation races with terminal publication?** Acceptance and terminal eligibility serialize through the same authority, but signal execution must also have an explicit owner. This document adds the constraint that accepted propagation settles before terminal publication, preventing context-release teardown from overtaking it. The main implementation risk is turning that constraint into a lock-held wait or dependency cycle; the paused-propagation and watcher-stop schedules must falsify those designs. [CG-03; IN-12, Pattern 6]

**Is the quiescence boundary weaker than the cleanup evidence suggests?** It is deliberately narrower than “every goroutine has exited when `Wait` returns.” The PRD requires complete terminal publication, while the language decision requires owned release paths. The selected contract permits only a finite, non-consumer-dependent internal epilogue and requires separate proof of its completion. If implementation or API documentation promises join-before-return, that stronger claim needs its own test rather than being inferred from readiness closure. [IN-03 §17; IN-05 §11.3; CG-08]

**Does the context design remain under-specified?** The semantic obligations are fixed: parent cancellation cannot bypass acceptance, the first accepted cause must reach work, already-cancelled parents are handled before invocation, and parent-derived values/deadline information must be accounted for. The concrete watcher/callback and child-context construction remains a Sprint 2 decision requiring inspection of the selected Go toolchain. This area does not claim that a plain child context plus watcher already satisfies those obligations. [CG-03–CG-04; IN-03 §§10–12]

**Did subprocess evidence cause a forceful-stop promise?** No. Signal escalation, PID ownership, reaping, pipe draining, and process-exit reclamation are explicitly rejected as Phase 1 mechanisms. They explain why substrate-specific stop contracts matter; they cannot supply one for an arbitrary Go function. [IN-11 P6; IN-14 P5; IN-03 §§4, 19]

**What remains unproven?** All Aren-specific behavior and release assertions remain unexecuted. Report evidence is second-hand, several reports have coverage caveats, and no source demonstrates this exact acceptance/publication/quiescence combination. The missing reasoning README also prevents claiming complete process-standard compliance. These limitations require targeted sprint proof and the governed passing review, not broader scope or an unsupported declaration of correctness. [IN-01, Project Reasoning Policy; IN-06, Evidence quality assessment; IN-09; IN-14, Open Questions]
