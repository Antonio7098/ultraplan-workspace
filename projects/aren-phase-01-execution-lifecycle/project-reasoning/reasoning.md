# Phase 1 Project Reasoning

> Project: `aren-phase-01-execution-lifecycle`
> Target repository: `/home/antonioborgerees/coding/ultraplan/Aren`
> Implementation language: Go
> Scope: One supervised in-process execution lifecycle
> Status: Project decision synthesis awaiting governed review
> Evidence boundary: Supplied project documents and area outputs; no current implementation inspection or test execution

## Purpose

Phase 1 must establish that Aren can own one execution honestly: assign its identity, enforce its lifecycle, interpret completion, propagate accepted cancellation, retain committed events, and support independent readers without inventing guarantees about arbitrary work effects.

This synthesis makes the cross-sprint semantic decisions. It does not select exact Go packages, exported signatures, synchronization primitives, or test names. Those choices belong to sprint reasoning and must satisfy the contract established here.

The evidence is sufficient to proceed to focused sprint design **after the required project-reasoning review returns `pass`**. It is not evidence that Aren already implements these semantics. An inherited governance check remains open: `projects/aren-phase-01-execution-lifecycle/reasoning/README.md` has not been verified against this synthesis or the area outputs. That check must close before a passing review.

## Authority And Scope

### Source Precedence

Source authority is role-specific, rather than a single ranking of all documents:

1. The project index governs scope, catalog status, and the required review gate. The project-reasoning index governs selected areas, assignments, and dependencies.
2. The Phase 1 PRD governs settled runtime behavior. The project roadmap governs delivery boundaries and Sprint 1 carry-forward.
3. The accepted language decision fixes Go and its applicable engineering obligations. Its prototype observations remain historical evidence.
4. The phased roadmap and lineage provide direction and scope-control doctrine. Their provisional or historical statements do not reopen settled Phase 1 requirements.
5. Area reasoning supplies the assessed evidence and proposed semantic resolutions. This synthesis adopts or qualifies those resolutions explicitly.
6. Study reports provide mechanisms, counterexamples, and verification techniques. They cannot override the PRD or substitute for current implementation evidence.

References used below resolve to exact paths in the Evidence section.

### Accepted Constraints

The following are governed constraints, not preferences contingent on implementation convenience:

- Aren supervises one in-process work function with capability equivalent to `func(context.Context) (Result, error)`.
- Aren alone owns lifecycle transitions, sequence allocation, canonical recording, terminal publication, and waiter release.
- The states are `created`, `running`, `succeeded`, `failed`, and `cancelled`.
- The only legal edges are `created -> running` and `running -> succeeded/failed/cancelled`.
- `created` is retained history, not an externally actionable scheduling state.
- Cancellation acceptance is an occurrence while running, not another lifecycle state or a terminal candidate.
- A run may commit at most one terminal outcome. A completed run has exactly one terminal outcome and terminal event.
- Work that never returns need not reach terminality. Aren must not manufacture completion to satisfy an eventual-completion claim.
- State, event, sequence, timing, cancellation metadata, and outcome become visible coherently.
- Waiting observes completed publication; it neither performs finalization nor consumes the only outcome.
- Canonical lifecycle history remains available while the run remains reachable.
- Slow, absent, or abandoned observers cannot control execution progress.
- Lifecycle bookkeeping does not make work transactional, exactly once, reversible, or free of effects on failure.
- Go race-enabled verification, independent semantic checks, defensive-publication tests, and owned-resource evidence are mandatory.

Providers, tools, subprocess supervision, retries, partial output, cleanup hooks, persistence, restart recovery, pause/resume, workflows, daemon hosting, remote APIs, budgets, and a permanent universal executor remain excluded. Test barriers, fault injection, resource probes, and execution of the diagnostic CLI serve current verification requirements; they do not introduce those product capabilities. [PRD §§4–19, 21–25; ROADMAP; GO §11]

## Project conclusions

### Decision Register

“Accepted interpretation” means this synthesis selects a precise meaning where the PRD leaves technical detail open. Such decisions constrain sprint reasoning but remain revisable through explicit evidence and recorded supersession. They are not upstream findings or implemented behavior.

| ID | Conclusion | Status | Primary Basis |
| --- | --- | --- | --- |
| PC-01 | One private per-run authority serializes lifecycle validation, applicable-fact resolution, and complete commitment. | Accepted constraint | PRD §§9, 13–17; A01 LA-01–LA-05 |
| PC-02 | Supported coherent reads return one committed fact set; independent accessor calls are not a transaction. | Accepted interpretation of atomic publication | A01 LA-06–LA-07 |
| PC-03 | Terminal interpretation follows the PRD truth table, separately from exclusion of later terminal writes. | Accepted constraint | PRD §13; A02 OF-02–OF-05 |
| PC-04 | Cancellation acknowledgement uses directional `errors.Is(returnedError, acceptedEffectiveCause)`, without a generic sentinel fallback. | Accepted interpretation | A02 OF-06 |
| PC-05 | Acceptance precedes work-visible signaling; accepted propagation settles before `accepted` returns or terminal publication occurs. | Accepted interpretation and ordering obligation | GO §7; A03 CG-01–CG-05 |
| PC-06 | Terminal publication proves completed invocation and complete lifecycle facts, not automatic joining of every support carrier. Support release remains mandatory and independently verified. | Accepted interpretation | A03 CG-07–CG-08 |
| PC-07 | Outcomes form a validated discriminated contract; Aren-owned lifecycle data is immutable, with explicit result and diagnostic-object ownership limits. | Accepted constraint with ownership interpretation | PRD §§11, 15, 17; A02 OF-01, OF-07; A04 EO-07 |
| PC-08 | Illegal mutation is rejected visibly. Invariant containment cannot fabricate an edge, overwrite terminal truth, or claim active work stopped. | Accepted constraint with containment interpretation | A01 LA-08; A02 OF-08–OF-09 |
| PC-09 | Canonical events are retained directly by lifecycle authority, with contiguous sequences beginning at zero and no observer-dependent recording. | Accepted constraint and derived ordering rule | PRD §§14–16; A04 EO-01–EO-02 |
| PC-10 | Observation uses inclusive next-sequence cursors; the current tail is valid, beyond-tail positions are rejected, and normal exhaustion requires terminal publication plus suffix exhaustion. | Accepted interpretation | A04 EO-03, EO-05 |
| PC-11 | Pull-oriented observation is the preferred implementation direction, not a mandatory Go interface or synchronization mechanism. | Provisional mechanism preference | A04 EO-04 |
| PC-12 | Verification combines an independent small model, exact trajectories, forced schedules, negative controls, race detection, release evidence, measurements, and real CLI execution. | Accepted verification strategy | A05 VG-01–VG-10 |
| PC-13 | Sprint 1 delivers a runnable, race-tested core; Sprint 2 extends its realized decisions explicitly rather than redesigning by implication. | Accepted delivery constraint | ROADMAP; all area handoffs |

### Lifecycle Authority And Publication

Aren generates an opaque identifier before recording `run.created`. The identifier remains unchanged and appears in every event and outcome. Sprint 1 selects its representation, generation mechanism, and construction-failure boundary without adding caller IDs or idempotency keys.

The per-run authority owns current state, retained history, sequence, lifecycle timing, accepted cancellation facts, terminal outcome, and completion readiness. Work supplies a return or captured panic. Callers supply cancellation requests. Observers maintain only their own observation progress.

Logical single-writer behavior does not require a dedicated goroutine. A direct critical section, ownership loop, or complete immutable-record publication can satisfy the contract. No mechanism is selected merely because a studied runtime uses it. Every supported reader must participate in the chosen publication discipline; serialized writes with unsynchronized or lagging read mirrors are insufficient.

The publication boundaries are:

| Boundary | Authoritative Facts | Required Ordering |
| --- | --- | --- |
| Creation | Identity, `created`, creation event at sequence `0`, occurrence time | Identity exists before recording |
| Start | `running`, start time, start event at sequence `1` | Commitment precedes invocation |
| First cancellation acceptance | First reason/effective cause, request event, next sequence; state remains running | Commitment precedes work-visible signal |
| Terminal commitment | Terminal state, finish time, validated outcome, matching terminal event, next sequence, completion readiness | Invocation finished; any accepted propagation settled; complete facts installed before waiter return |

No supported operation may observe a sequence increment without its event, terminal state without terminal history and outcome, or an outcome while the run is nonterminal.

A coherent snapshot must capture its related fields together. Separate calls may observe different valid prefixes. Reading `running` and later reading terminal history is valid. Observing terminality and then starting a fresh outcome or full-history read that lacks terminal facts is not.

Occurrence times describe Aren-owned lifecycle boundaries, not observer receipt or the exact instruction at which user work returned. Start and finish are captured once and reused consistently; sequence, not timestamps, orders events. Equal times are valid. Start-after-finish and active finish timestamps are not. A clock abstraction remains conditional on demonstrated testing need. [PRD §§8–10, 14–18; A01; A02]

### Outcomes And Terminal Resolution

The terminal contract is discriminated by state:

| State | Successful Result | Failure | Terminal Cancellation |
| --- | --- | --- | --- |
| `succeeded` | Exact returned value, including nil or zero | Absent | Absent |
| `failed` | Unavailable | Structured origin, kind, and diagnostics | Absent |
| `cancelled` | Unavailable | No failure classification branch | First accepted cancellation metadata and acknowledgement diagnostics |

Every branch includes identity, start time, and finish time. Accepted cancellation remains historical fact even when the eventual outcome succeeds or fails; it does not create a terminal-cancellation payload on those outcomes.

The central policy is:

| Captured Facts | Resolution |
| --- | --- |
| Existing terminal commitment | Retain unchanged; reject later internal mutation |
| Confirmed normal return with nil error | `succeeded`, even after accepted cancellation |
| Confirmed normal return with non-nil error, accepted cancellation, and matching effective cause | `cancelled` |
| Any other confirmed non-nil error return | `failed`, `work/returned_error` |
| Captured escaping work panic | `failed`, `work/panic`, regardless of cancellation |
| Work still active | No terminal candidate |

A non-nil result plus error is valid Go input to this policy, but never partial success. A non-nil error interface containing a typed nil remains an error. An error-shaped value returned as the result does not make a nil-error return fail.

Work cannot spoof Aren-origin failure by returning an Aren-shaped error. This invocation still classifies it as a work returned error, retaining the original structure beneath that classification unless cancellation matching applies. Only Aren's internal fault path originates this run's `aren/invariant_violation`.

**Cause matching.** The selected predicate is:

```text
cancellation was accepted
AND returnedError is non-nil
AND errors.Is(returnedError, acceptedEffectiveCause)
```

The effective cause is non-nil, retained at first acceptance, and used to signal the work context. Matching is directional and follows Go error equivalence, not message equality or inferred physical causation.

Consequences are explicit:

- A direct or wrapped accepted cause matches.
- An unrelated error with the same text does not match.
- A coarse `context.Canceled` return does not automatically match a distinct custom accepted cause.
- Work acknowledging a custom cause should return or wrap `context.Cause(ctx)`.
- An accepted parent deadline cause can resolve to `cancelled`; there is no timeout terminal state.
- A joined error containing the accepted cause matches even if another branch is unrelated. Preserve the complete acknowledgement error locally rather than hiding that branch.
- Panicking with the accepted cause remains panic failure.

This is a deliberate usability trade-off, not comparative consensus. Shared sentinels and custom `Is` methods cannot authenticate why work stopped.

Error matching and diagnostic preparation can invoke user code. They must occur outside lifecycle mutation ownership. Before committing, Aren must validate that classification used the applicable accepted facts. First acceptance can occur during preparation; the implementation must re-evaluate that newly accepted cause rather than publish a stale candidate. Since acceptance moves only from absent to immutable-present, this is a bounded fact-revalidation obligation, not a retry feature. [PRD §§11–14; A02; A03]

### Cancellation And Owned Lifetimes

Explicit cancellation, parent cancellation, and supplied-parent deadline expiry use one acceptance path. There is no source priority for genuinely concurrent requests. “First” means first accepted under lifecycle authority, not earliest timestamp or earliest unsynchronized call.

Disposition precedence is:

1. Terminal commitment exists: `already_terminal`.
2. Otherwise cancellation already accepted: `already_requested`.
3. Otherwise running: retain the first reason/effective cause and request event, propagate that cause, then return `accepted`.

Rejected and repeated requests do not allocate sequence numbers, replace metadata, restart support, or close readiness again. A response describes its operation's decision point; `accepted` is not a guarantee that the run remains running when the caller receives the response.

Acceptance and propagation may have a visible interval, but that interval has an owner. The accepting operation must settle propagation before returning `accepted`; terminal publication cannot overtake outstanding accepted propagation. This prevents terminal resource release from installing a different context cause before the accepted one reaches work.

A direct automatically cancelling child context plus an independent watcher is not sufficient if parent cancellation can signal work before Aren records acceptance. Sprint 2 must prove the whole context path, including parent values, deadline behavior, `Done`, `Err`, effective cause, startup, and deregistration.

An already-cancelled parent still produces creation and start. Acceptance and signaling complete before work entry, and work is invoked once. Nil-error success, matching cancellation, unrelated failure, and panic then follow normal resolution. Cancellation racing startup must not be lost in a check/registration gap.

Aren distinguishes:

```text
request
    -> acceptance
    -> signal propagation
    -> possible work observation
    -> completed invocation
    -> terminal publication
    -> support release
```

Observation of the signal is not generally knowable by Aren. A matching return is contractual acknowledgement, not authenticated causation. Ignored cancellation leaves the run running indefinitely, including while a work defer remains blocked.

The run owner accounts for its invocation carrier, work-context resources, parent registration or watcher, and any ownership-loop or observation support actually introduced. Waiter and observer caller goroutines remain caller-owned. User-spawned work goroutines are not automatically adopted or joined.

Outcome readiness is not automatically total support quiescence. Any post-terminal epilogue must have a finite release path under ordinary scheduling, with no dependency on observer drainage, parent cancellation, new user cooperation, or another future lifecycle event. Its actual completion must be independently observable in tests. No wall-clock shutdown guarantee is inferred.

Releasing context resources after work completion is teardown, not a new accepted request. It must not add cancellation history or change the outcome. General work cleanup hooks and cleanup-failure policy remain Phase 2 concerns. [PRD §§10–13, 16–19; A03]

### Failure Visibility And Immutability

Aren protects lifecycle data across every outcome, history, snapshot, and cursor access path. Top-level copies alone are insufficient if nested slices, maps, or stack storage remain writable aliases.

Two ownership limits are explicit:

- Arbitrary successful result graphs remain user-owned. Aren preserves the exact result, not a universal deep-frozen clone.
- Original errors and custom causes are retained as borrowed read-only diagnostic objects where practical. Their identity and wrapping remain inspectable locally; Aren owns stable classifications and diagnostic projections, not every object reachable through those errors.

Supported diagnostic methods must be stable, terminating, and safe for read-only inspection. Moving them outside the lock prevents lifecycle reentrancy deadlocks but does not sandbox hostile or nonterminating methods.

Work panic recovery is narrow: panics escaping the supplied invocation, including its defers, become recognizable `work/panic` with captured stack diagnostics. A broad recovery around the whole runtime must not relabel Aren defects as work panic.

Invariant handling rejects illegal mutation before changing canonical truth. A truthful `running -> failed` containment is permitted when the invocation has finished and the remaining authority, history, and timing are trustworthy. It is not permitted to invent `created -> failed`, replace a completed outcome, or publish terminality while work remains active. Faults beyond trustworthy containment require explicit Aren-origin escalation, not recursive repair.

Sprint 1 must resolve and test abnormal invocation exits such as `runtime.Goexit`, diagnostic-method faults, and the concrete fail-fast boundary. None may become zero-value success, disappear as a clean non-match, or be covered by an unsupported process-isolation claim. These are named sprint-entry decisions, not an invitation to add new public failure categories or a general supervisor. [PRD §§11, 17, 19; GO §§7, 11.2, 11.9; A01; A02]

### Events, Waiting, And Replay

There is one canonical event class and one recording authority. The vocabulary is exactly:

```text
run.created
run.started
run.cancellation_requested
run.succeeded
run.failed
run.cancelled
```

Creation and start occur once. Cancellation request occurs at most once, while running. Exactly one terminal event ends a completed history. Thus ordinary completed histories contain three events without accepted cancellation and four with it. This is a count bound, not a total-byte or process-capacity bound.

Replay preserves original identities, payloads, times, and sequences. It neither appends events nor reconstructs authoritative state nor re-executes work.

**Cursor rules.** Let `n` be committed history length and `c` the next requested sequence:

| Condition | Behavior |
| --- | --- |
| Default | `c = 0` |
| `0 <= c < n` | Read inclusively from `c` |
| `c = n`, active run | Wait for availability or local read cancellation |
| `c = n`, terminal run | Normal exhaustion |
| `c > n`, or invalid representation | Explicit invalid-cursor result; no clamping or speculative future wait |
| Event `s` returned | Advance this traversal to `s+1` |
| Local abort without returned data | No cursor advance and no run mutation |

Validation reflects one coherent decision point. Independent observers have independent cursors. Concurrent advancement of one mutable reader object is not an additional Phase 1 promise unless sprint reasoning explicitly selects and tests it.

Normal observation exhaustion requires terminal publication and no remaining requested suffix. An observer starting at a completed tail requested an empty suffix and need not receive the terminal event again; that event must nevertheless already be retained and retrievable.

Blocking live reads require observation-local abandonment. At a coherent read decision, available data takes precedence, then terminal exhaustion, then local cancellation. A cancellation wake must recheck those predicates. Local abort is not terminal exhaustion or run cancellation.

The preferred realization is pull-oriented reading of the same retained history. It avoids mandatory per-observer producer queues, but still requires a proved tail-check/wait protocol. An append must either already be visible during inspection or be covered by the established availability condition. Wakeups may coalesce; canonical events may not.

Outcome waiting observes reusable terminal publication. It does not wait for observer processing. If a context-cancellable outcome wait is offered, its cancellation remains caller-local and its ready-versus-abort rule must be explicit. Sprint 1 need not add that API merely for symmetry.

No exactly-once delivery or unconditional at-least-once consumer-processing promise is made. One continuing traversal reads its suffix without gaps or unsolicited repeats; deliberate replay can expose the same `(run_id, sequence)` again. [PRD §§14–19; A04]

### Verification Strategy

A05's invariant catalog, `INV-01` through `INV-18`, is the project verification reference. Sprint artifacts must map those invariants to actual tests and evidence rather than copying the catalog into an unmaintained parallel checklist.

| Proof Class | Required Evidence | What It Cannot Establish Alone |
| --- | --- | --- |
| Independent model and exact trajectories | Separate legal-edge table, outcome policy, acceptance facts, cursor rules, and committed-prefix assertions | Go memory ordering, actual resource ownership, or all schedules |
| Real invocation and commit tests | Success, error, panic, invalid construction, duplicate finalization, timing, cause preservation | Every concurrent interleaving |
| Controlled schedules | Acceptance/signal, stale matching, startup, publication, waiter readiness, tail handoff, observer abandonment, registration release | Exhaustive scheduler coverage |
| Negative controls | Broken fixtures and named production-path defects rejected for the intended reason | Correctness outside the challenged instrument |
| Complete race-enabled suite | Current execution of all required semantic tests under `go test -race` | Logical deadlocks, correct arbitration, leak freedom, or gate completeness |
| Release and retention checks | Actual exit/stop witnesses, registration accounting, live-parent tests, supplementary stacks and heap evidence | Universal reclamation from one goroutine count or finalizer |
| Measurements | Allocations, active/retained footprint, carrier growth, relevant contention and latency, churn recovery | A production capacity guarantee without a workload budget |
| Real CLI and review | Built-runtime scenarios, useful statuses, implementation-linked contract, effective CI gates | Library leak freedom merely because the process exits |

The five-state transition matrix has four legal and 21 forbidden ordered pairs, plus invalid vocabulary cases. Cancellation-request recording is not a legal `running -> running` transition.

Tests must not use production guards or resolution helpers as the expected oracle. Real work facts, request responses, and externally established order constrain concurrent-history checks; the final runtime state cannot be the sole source of expectation.

Barriers and acknowledgements establish important orderings. Sleeps do not. Test watchdogs bound failure and gather diagnostics; they do not create runtime timeout semantics. Intentionally blocked work must be released on every test exit before final leak assertions.

Historical negative controls include parent bypass, request/terminal inversion, mutable outcomes, missing finish timing, flattened causes, recursive invariant handling, duplicate finalization, lost wakeups, and retained parent watchers. A hook that masks its designated defect is not proof.

A05's proposed 20 repeated controlled schedules and 10,000 mixed collisions are starting budgets only. Sprint reasoning must finalize counts, commands, nonempty test selectors, timeout budgets, and ongoing cadence from measured cost. Reducing duplicated work is acceptable; silently dropping a required ordering or race-enabled case is not.

Evidence records identify requirement, invariant, actual test/subtest, command, revision, diff status, toolchain, result, negative control, and limitation. Stress adds seed, iteration, operation trace, and reached checkpoints. Measurements and built CLI execution add artifact digest and relevant hardware/workload provenance. A random seed reproduces inputs, not Go scheduling. Retry-to-green is not acceptance. [A05; GO §§7, 11.3, 11.6–11.8]

## Contradictions And Decisions

| Tension | Resolution | Remaining Obligation |
| --- | --- | --- |
| Phased roadmap leaves cancellation-request state open | PRD settles it as a running-state event only | No extra state in API or model |
| First-winner language suggests cancellation may defeat success | First acceptance determines retained cancellation facts; PRD interpretation still makes nil-error return succeed | Test interpretation separately from one-commit exclusion |
| Language remains open in lineage | Accepted Go decision supersedes that historical clause | No shadow runtime or planned port |
| Broad immutable-publication language conflicts with exact arbitrary results | Protect Aren-owned lifecycle data; preserve PRD's explicit result exception | Document and test alias boundaries |
| Cause identity conflicts with universal deep freezing | Borrow original causes read-only and publish stable diagnostic projections | Sprint 1 diagnostic-fault mechanics; no blanket immutability claim |
| Strict cause matching conflicts with familiar `ctx.Err()` usage | Adopt strict effective-cause matching without sentinel fallback | Sprint 2 documentation and custom-cause matrix |
| Joined cause plus unrelated error could reasonably imply failure | Adopt standard `errors.Is` match and retain the full acknowledgement error | Any failure-precedence change requires explicit supersession |
| Parent-derived context may cancel before acceptance | Automatic propagation cannot bypass Aren authority | Sprint 2 must prove context values/deadline/cause behavior and startup handoff |
| Atomic commitment might be read as requiring signal delivery inside the lock | Commit acceptance first; signaling is an ordered consequence; terminality cannot overtake it | No user callbacks, blocking delivery, or joins under mutation ownership |
| Cleanup language might imply all carriers joined before outcome wait returns | Terminal readiness and support release are separate; no indefinitely parked epilogue is allowed | Prove actual release independently |
| “Terminal event before observation completion” conflicts with opening at terminal tail | Exhaustion is relative to requested suffix; terminal record must already exist | Test empty suffix separately from premature EOF |
| Illegal mutation visibility could be mistaken for a requirement to append a fault event after terminality | Reject and expose the defect without changing terminal history | Concrete out-of-band internal fault boundary in Sprint 1 |
| Broad study ratings suggest more evidence than populated sources support | Use only assessed, populated-source findings; repeated reports are not independent reproductions | Preserve A00 restrictions in sprint distillation |
| Mandatory measurements suggest fixed capacity targets | Measure applicable local behavior without importing prototype percentages as budgets | Record unexplained growth and relevant metrics; justify inapplicable metrics explicitly |

No unresolved contradiction is silently delegated as an implementation preference. Fixed product rules remain binding. The concrete context, diagnostic-fault, and release mechanisms remain open and must close in their assigned sprint reasoning before the relevant implementation plan is finalized.

## Sprint Handoff

### Sprint 1: Core Lifecycle

Sprint 1 must deliver a buildable, runnable non-cancellation core with identity, guarded transitions, work invocation, success/error/panic, terminal timing and outcomes, shared waiting, coherent inspection, defensive retained history, and complete ordinary negative and race tests.

Its reasoning must select:

- Minimal package ownership and actual observation capability.
- Identity generation and pre-creation failure behavior.
- Concrete commit/read publication protocol and memory-order argument.
- Outcome, failure, event, and snapshot representations.
- Narrow panic recovery, abnormal-exit behavior, diagnostic-fault handling, and non-recursive invariant exposure.
- Existing carrier ownership and release witnesses.
- Supported toolchain/platforms, independent model, actual test selections, CI gates, and baseline measurements.

Caller cancellation, full parent integration, live cursors, and observer abandonment remain deferred unless a current Sprint 1 requirement genuinely earns part of them. A context parameter does not by itself justify claiming Phase 1 cancellation completion. The realized intermediate behavior and limitations must be documented honestly.

Sprint 1 acceptance requires fresh equivalents of `go build ./...`, `go test ./...`, `go vet ./...`, and complete `go test -race` execution, together with invariant-linked negative and release evidence. Commands here describe required execution forms, not results.

### Sprint 2: Cancellation And Concurrency

Sprint 2 implements the accepted cancellation and cursor semantics against the realized core. It must decide and prove:

- Controller capability, reason/cause normalization, and all three dispositions.
- Parent context construction, deadline/value behavior, startup handoff, and support deregistration.
- Acceptance-to-signal settlement and cause-match revalidation.
- Independent live observation, local abandonment, strict cursor boundaries, and tail notification.
- Completion/cancellation collisions, watcher installation/release races, observer churn, mutation attacks, and leak resistance.
- Required diagnostic CLI scenarios and useful process statuses.
- Final measured behavior, phase review, and lifecycle-contract promotion.

The required real-runtime scenarios are:

```text
aren dev run success
aren dev run fail
aren dev run cancel
aren dev run race
```

Panic, parent cancellation, and ignored cancellation are recommended demonstrations. CLI output must show identity, sequence-ordered events, acceptance distinct from terminal cancellation, and the outcome without implying durable replay or exactly-once effects.

### Carry-Forward And Gates

Sprint 2 must explicitly consume these artifacts under `projects/aren-phase-01-execution-lifecycle/`:

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

Its synthesis includes **Prior Sprint Decisions Applied**, classifying each relevant decision as preserved, extended, superseded, or unaffected. Supersession names the prior decision, new evidence, behavioral and architectural impact, and updated plan/tests. Changed goldens or silence do not supersede semantics.

| Gate | Required Closure |
| --- | --- |
| Before either sprint starts | Verify the named reasoning README; reconcile applicable standards with all outputs; obtain actual project-reasoning review verdict `pass` |
| Before Sprint 1 implementation planning completes | Resolve its concrete authority, publication, failure, abnormal-exit, API, and verification choices |
| Before Sprint 1 acceptance | Runnable core, fresh required checks, effective negative controls, resource evidence, and review recording realized decisions |
| Before Sprint 2 implementation planning completes | Consume Sprint 1 artifacts and close context, cause, signal, cursor, and release mechanics |
| Before Phase 1 acceptance | Both sprint acceptances, full required race/stress/leak evidence, real CLI scenarios, and no unresolved foundational semantic ambiguity |
| Before stable contract promotion | Final contract agrees with realized tests and phase review; promotion does not precede validation |

No third sprint is scheduled. A closure sprint requires the roadmap's evidence-based exception, not routine deferral of unfinished testing or documentation.

## Trade-Offs

**Strict semantics over convenient cancellation labels.** Nil-error success after acceptance and strict cause matching preserve the selected contract but may surprise callers expecting every post-cancel return to be cancelled. Explicit examples and tests are preferable to a hidden fallback.

**One authority over independently optimized fields.** Coherent publication imposes synchronization discipline. The five-state, three/four-event lifecycle does not yet justify complex read mirrors or lock-free machinery. Mechanisms remain revisable after measurement.

**Pull-oriented observation over unsolicited fan-out.** Pull reading reduces per-observer queues and producer lifetimes. It still requires a precise availability protocol and local abort. Channel ergonomics remain possible only if they earn their additional proof cost.

**Cause fidelity over universal object isolation.** Preserving Go error identity enables useful local inspection but requires a read-only object contract. Stable diagnostic projections protect lifecycle metadata without pretending arbitrary user graphs are immutable.

**Truthful terminality over artificial shutdown bounds.** Work can ignore cancellation indefinitely. Separating publication from support release avoids false stop claims and lock-held joins, but requires explicit resource witnesses rather than one reassuring `Wait` call.

**Independent proof over green-test totals.** Model, schedule, mutation, and release checks cost more than happy-path assertions. Their value comes from rejecting named defects, not from suite size. The proof effort should remain a small set of maintained instruments rather than a generic testing platform.

**Bounded decisions over complete comparative coverage.** Missing final reports leave local proof obligations, not an automatic research program. Further study is justified only when a named unresolved question cannot be settled through focused inspection or a small experiment.

## Evidence

### Inputs Used

The supplied full project documents and supplied portions of all six area outputs were used. A00, A03, A04, and A05 contain omitted middles; this synthesis does not claim to have inspected those omitted sections. Study findings enter through the area assessments, not fresh inspection of the underlying reports or repository files.

| Reference | Exact Path | Role |
| --- | --- | --- |
| INDEX | `projects/aren-phase-01-execution-lifecycle/project-index.md` | Scope, catalog, review gate |
| SELECTION | `projects/aren-phase-01-execution-lifecycle/project-reasoning/index.md` | Selected areas, evidence, dependencies |
| PRD | `projects/aren-phase-01-execution-lifecycle/docs/PRD.md` | Authoritative product behavior |
| ROADMAP | `projects/aren-phase-01-execution-lifecycle/roadmap.md` | Delivery and carry-forward |
| LINEAGE | `projects/aren-phase-01-execution-lifecycle/docs/project-lineage.md` | Scope-control and earned-abstraction doctrine |
| PHASES | `projects/aren-phase-01-execution-lifecycle/docs/phased-roadmap.md` | Directional later-phase boundary |
| GO | `projects/aren-phase-01-execution-lifecycle/docs/final-language-decision.md` | Accepted language, engineering obligations, historical defects |
| A00 | `projects/aren-phase-01-execution-lifecycle/project-reasoning/areas/00-evidence-map.md` | Assessment and claim-level evidence restrictions |
| A01 | `projects/aren-phase-01-execution-lifecycle/project-reasoning/areas/01-lifecycle-authority-and-atomic-publication.md` | Authority and coherent publication |
| A02 | `projects/aren-phase-01-execution-lifecycle/project-reasoning/areas/02-outcomes-failures-and-terminal-resolution.md` | Outcome, matching, diagnostics, containment |
| A03 | `projects/aren-phase-01-execution-lifecycle/project-reasoning/areas/03-cancellation-goroutine-ownership-and-cleanup.md` | Acceptance, propagation, ownership, release |
| A04 | `projects/aren-phase-01-execution-lifecycle/project-reasoning/areas/04-events-observation-waiting-and-replay.md` | Recording, cursor, waiting, delivery |
| A05 | `projects/aren-phase-01-execution-lifecycle/project-reasoning/areas/05-verification-and-go-correctness.md` | Independent verification and sprint gates |

### Evidence Weight

The direct Go studies are the closest mechanism evidence, as assessed by A00:

- `studies/aren-go-runtime-study/reports/final/01.01-lifecycle-transition-ownership-and-terminal-arbitration.md`
- `studies/aren-go-runtime-study/reports/final/01.02-cancellation-goroutine-ownership-and-cleanup.md`
- `studies/aren-go-runtime-study/reports/final/01.03-adversarial-concurrency-and-failure-verification.md`
- `studies/aren-go-runtime-study/reports/final/01.04-ordered-observation-live-streaming-and-backpressure.md`

Their useful contributions are concrete publication, ownership, handoff, and verification mechanisms plus negative examples. Their source analyses lacked the authoritative Aren PRD. A first-write latch is not Aren's terminal policy; worker shutdown is not goroutine force-stop; local replay handoff is not end-to-end canonical completeness.

A00 assesses all 23 selected reports. Its material restrictions remain binding:

- Empty source trees establish evidence unavailability, not absent upstream capabilities or poor implementation quality.
- Several cancellation, taxonomy, API, interface, and CI reports have much narrower material coverage than their nominal source lists.
- Repeated citations to the same implementation are not independent reproductions.
- Durable-engine and high-volume streaming mechanisms require narrow translation.
- Missing and regeneration-required dimension reports remain unavailable comparative evidence. The Go reports supplement questions; they are not silently relabeled replacements.

The historical Aren defects in GO §7 have unusually strong relevance as regression seeds. Its dated green commands, prototype architecture, and memory measurements do not establish current correctness or capacity.

The exact cause matcher, cursor boundaries, propagation-settlement rule, and publication/quiescence distinction are **project interpretations requiring local proof**. No selected source is claimed to demonstrate the complete combined contract.

## Risks

| Risk | Consequence | Control And Owner |
| --- | --- | --- |
| Reasoning approval assumed | Sprint work begins without the required gate | Verify the inherited README requirement and obtain a real `pass`; project review |
| Parent context bypass or deadline mismatch | Signal precedes history, or work sees inconsistent cause/deadline behavior | Whole-context argument and startup/collision tests; Sprint 2 |
| Accepted propagation coordination cycles | Completion waits for signaling while signaling waits on lifecycle completion | Explicit dependency ordering, no lock-held waits, held-gap negative control; Sprint 2 |
| Reentrant or faulty diagnostic methods | Deadlock, false classification, or hidden abnormal exit | Outside-lock inspection, stable diagnostic contract, explicit fault boundary; Sprint 1 and Sprint 2 extension |
| Properly synchronized but torn publication | Race detector passes despite invalid lifecycle facts | Independent coherent-prefix oracle and observable split-publication mutant; both sprints |
| Fault containment becomes another state machine | Illegal recovery edge, false stop claim, duplicate terminal | Trustworthy-state preconditions and non-recursive fault tests; Sprint 1 |
| Snapshot or stack aliases escape | One reader changes another's lifecycle evidence | Field-level ownership inventory and mutation attacks; both sprints |
| Observer handoff loses a wake | Correct history exists but a reader never progresses | Atomic availability protocol and forced append-in-gap schedule; Sprint 2 |
| Completed run remains parent-retained | Embedded hosts accumulate references without obvious goroutine growth | Live-parent registration and retention tests; Sprint 2 |
| Quiescence wording excuses parked support | Terminal runs retain watchers or producers indefinitely | Actual exit witnesses and independent leak evidence; both sprints |
| Strict matching is undocumented | Custom-cause work unexpectedly fails instead of cancelling | Public contract examples and full matcher matrix; Sprint 2 |
| Count bounds become memory claims | Large diagnostics, results, or caller-retained runs are concealed | Separate lifecycle count, payload bytes, support, and caller retention in measurements |
| Tests cannot fail meaningfully | Green suites repeat prototype false-positive mistakes | Named negative controls and effective CI gate checks; both sprints |
| Project reasoning selects too much machinery | A small run acquires supervisors, brokers, or compatibility layers | Keep representations provisional; simplification review may remove them |
| Finite evidence is described as exhaustive proof | Unmodeled schedules or unsupported platforms are overlooked | State evidence limits, retain focused exploration, review every new mutation/resource boundary |

## Self-critique

**Has this synthesis settled semantics without pretending to settle implementation?** The terminal policy, cause matching, propagation ordering, cursor boundaries, ownership limits, and required evidence are explicit. Exact packages, context construction, notification mechanics, and public signatures remain sprint decisions. Pull-oriented observation is a preference, not a prematurely fixed framework.

**Which decisions are least supported by comparative evidence?** Strict custom-cause matching, joined-error treatment, beyond-tail rejection, and propagation settlement before terminality are project choices. They are precise and falsifiable, but should not be presented as consensus from the studies. Concrete evidence can justify explicit revision; implementation convenience alone cannot silently change them.

**Does the reasoning leave foundational ambiguity unresolved?** It retains bounded implementation questions with named closure gates: context deadline/value behavior, abnormal invocation exits, diagnostic-method faults, and concrete invariant escalation. These must close before the relevant implementation plan and acceptance. The unchecked reasoning README is a separate project-review prerequisite. None is represented as already resolved or tested.

**Could this contract still permit an implementation bug to pass?** Yes. A model and implementation can share a misunderstanding; a hook can conceal a publication defect; an untracked registration can evade resource accounting; an unmodeled platform schedule can expose a new race. Independent fixtures, negative controls, source review, and supplementary stack/heap evidence reduce these risks without claiming exhaustive proof.

**Is support release weaker than users may assume from completion?** Yes, deliberately. Outcome waiting guarantees completed invocation and complete terminal truth, not the simultaneous exit of every internal carrier or user-spawned goroutine. The compensating obligation is a finite, consumer-independent release path with actual evidence. Documentation must not blur either side of that boundary.

**Has planning become larger than the runtime?** The area evidence is broad, but its output should now be a small core and focused tests. Regenerating every missing report, repeating every evidence table in sprint artifacts, or building generic supervision and verification frameworks would violate the project's scope-control doctrine.

**What does this document actually establish?** It establishes the authoritative project synthesis proposed for governed review, distinguishes inherited requirements from selected interpretations and provisional mechanisms, and assigns the remaining decisions and proof obligations. It establishes neither a review verdict nor runtime correctness, sprint completion, or Phase 1 acceptance.
