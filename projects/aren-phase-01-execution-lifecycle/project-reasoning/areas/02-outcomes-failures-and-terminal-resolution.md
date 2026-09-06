I’ll settle the outcome validity and cause-matching rules against the supplied publication contract, separating project decisions from report evidence and sprint-level implementation choices.
# Outcomes Failures And Terminal Resolution

> Project: `aren-phase-01-execution-lifecycle`
> Area: Outcomes Failures And Terminal Resolution
> Scope: Cross-sprint outcome validity, failure visibility, cancellation-cause matching, and terminal interpretation
> Evidence boundary: Supplied governed documents and report excerpts; no direct repository inspection or test execution

## Purpose

Define how Aren interprets one work invocation without confusing returned work facts, cancellation requests, and committed lifecycle outcomes. The central decision is to preserve the PRD's terminal policy through one validated terminal record: nil-error return succeeds, an accepted-cause-matching error cancels, another returned error fails, and a captured work panic fails recognizably as panic. An existing terminal commitment is never replaced. [IN-03 §§11–14]

This area consumes the private publication authority established by the lifecycle area. It does not choose Go packages, exported signatures, synchronization primitives, observer delivery mechanisms, or cleanup hooks. Its conclusions supply semantic constraints and falsifiable cases for those later choices. [IN-02, Reasoning Areas; IN-04, Canonical Project Flow; IN-07, Project conclusions LA-01–LA-10]

## Inputs Used

Repeated supplied copies are counted once. Source locations quoted inside a study report remain report-only evidence; no repository file or test body was separately inspected. The exact report sections below identify the portions used, not omitted portions of excerpted reports.

| Input ID | Input | Kind | Exact workspace-relative path | Exact portion inspected | How it was used | Limits or applicability warning |
| --- | --- | --- | --- | --- | --- | --- |
| IN-01 | Project index | Requirement and governance | `projects/aren-phase-01-execution-lifecycle/project-index.md` | Project Reasoning Policy; Project Scope; evidence catalog; unavailable-dimension table; Prior Decisions; Cross-Sprint Decision Carry-Forward | Scope, review gate, evidence availability | Catalog status is not runtime proof |
| IN-02 | Project-reasoning index | Governance | `projects/aren-phase-01-execution-lifecycle/project-reasoning/index.md` | Reasoning Areas; selected area's Evidence Assignments and Source Document Assignments; Excluded Evidence | Dependencies and decision routing | Assignment text is not an empirical finding |
| IN-03 | Phase 1 PRD | Requirement | `projects/aren-phase-01-execution-lifecycle/docs/PRD.md` | Full supplied copy, especially §§7–19 and §§21–25 | Authoritative outcome, cancellation, timing, publication, and acceptance rules | Normative, not implementation evidence |
| IN-04 | Project roadmap | Requirement and delivery policy | `projects/aren-phase-01-execution-lifecycle/roadmap.md` | Scope Principle; Canonical Project Flow; Cross-Sprint Carry-Forward Rule; both implementation waves; Phase Exit Gate | Separate complete Sprint 1 behavior from Sprint 2 extension | Planned commands have not been executed |
| IN-05 | Accepted Go decision | Requirement with historical evidence | `projects/aren-phase-01-execution-lifecycle/docs/final-language-decision.md` | §§1, 4, 6–7, 11, 14–15 | Private construction, cause retention, defensive publication, panic/invariant distinctions, regression seeds | Historical prototype observations and checks are not new-runtime verification |
| IN-06 | Evidence Assessment And Routing | Dependency reasoning | `projects/aren-phase-01-execution-lifecycle/project-reasoning/areas/00-evidence-map.md` | Supplied Purpose, Inputs Used, Scope, assessment tables, Phase 1 decision map, visible Requirement coverage; supplied ending containing conclusions, Trade-Offs, Evidence, Risks, and Self-critique | Apply source precedence, coverage restrictions, and routed uncertainties | Excerpt; report-only assessment; reasoning README remains unverified |
| IN-07 | Lifecycle Authority And Atomic Publication | Dependency reasoning | `projects/aren-phase-01-execution-lifecycle/project-reasoning/areas/01-lifecycle-authority-and-atomic-publication.md` | Full supplied copy, especially Atomicity model, Counterexample analysis, Project conclusions, and Downstream obligations | Consume coherent publication, immutable terminality, callback-free ownership, and containment constraints | Semantic dependency, not executed proof |
| IN-08 | Selected area template | Reasoning requirement | `projects/aren-phase-01-execution-lifecycle/reasoning/project-wide/02-outcomes-failures-and-terminal-resolution.md` | Full supplied copy | Required comparisons, fact table, resolution matrix, preservation analysis, and critique | Referenced reasoning README unavailable in supplied context |
| IN-09 | Completion and finalization | Study report | `studies/agent-harness-study/reports/final/03.09-completion-and-finalization-semantics.md` | Executive Summary; Core Thesis; Approach Models A–E; Pattern Catalog 1, 5, 7–9; Key Differences 1–5; supplied Open Questions and Evidence Index | Separate lifecycle completion from provider, framework, and task completion; compare narrow terminals | Report-only; no Aren result/error or timing proof |
| IN-10 | Error taxonomy | Study report | `studies/agent-harness-study/reports/final/13.01-error-taxonomy.md` | Executive Summary; Approach Models; Pattern Catalog P3, P5–P7, P9; Anti-Patterns; Notable Absences; Evidence Index | Structured classification, cause chains, closed vocabulary, string-matching hazards | Only Langfuse and OpenHands contain material evidence; eight empty payloads establish no upstream behavior |
| IN-11 | Failure visibility | Study report | `studies/agent-harness-study/reports/final/13.03-failure-visibility.md` | Executive Summary; Approach Models; Pattern Catalog P1, P3, P5, P7, P8, P10; Key Differences; supplied Per-Source Notes and Evidence Index | Preserve diagnostic facts across outcome/event/rendering boundaries; identify traceback loss and missing failure surfaces | Report-only; model feedback, durable records, redaction platforms, and telemetry exporters excluded |
| IN-12 | Recovery versus escalation | Study report | `studies/agent-harness-study/reports/final/13.04-recovery-vs-escalation.md` | Executive Summary; Approach Models, especially OPA, CrewAI, and Temporal; Pattern Catalog P2, P5, P7; Key Differences; supplied Open Questions and Evidence Index | Distinguish work failure, control flow, and runtime defects; reject recovery that hides bugs | Report-only; no retries, human escalation, or durable recovery in Phase 1 |
| IN-13 | Go lifecycle ownership and arbitration | Study report | `studies/aren-go-runtime-study/reports/final/01.01-lifecycle-transition-ownership-and-terminal-arbitration.md` | Full supplied copy, especially P1–P3, P5–P7, Anti-Patterns, Notable Absences, and Evidence Index | Completion funnel, first-write limitations, panic capture, invariant detection, and production/mock divergence | Report-only; source analyses lacked Aren PRD; no complete local timing/publication implementation demonstrated |
| IN-14 | Go cancellation, ownership, and cleanup | Study report | `studies/aren-go-runtime-study/reports/final/01.02-cancellation-goroutine-ownership-and-cleanup.md` | Executive Summary; Core Thesis; Approach Models; P1–P4, P6, P10–P12; Key Differences; supplied ending and Evidence Index | Cause-typed cancellation, primary-error preservation, stop versus join, and concurrent-close hazards | Report-only; exact Aren cause matcher is not established; process escalation and service shutdown excluded |
| IN-15 | Aren phased roadmap | Directional project document | `projects/aren-phase-01-execution-lifecycle/docs/phased-roadmap.md` | Development Rules; Phase 1 Terminal outcome, Cancellation, Failure model, and Time; Phase 2 Cleanup semantics; Immediate Planning Boundary | Reconcile provisional first-winner language and exclude cleanup-based outcome changes | Subordinate to the settled Phase 1 PRD |

The template requires reading `projects/aren-phase-01-execution-lifecycle/reasoning/README.md`. Its content was not supplied, and no read-only workspace tool is available in this turn. It is not an inspected input. Any additional requirements there must be checked before a passing project-reasoning review. This document does not claim that compliance, a review verdict, or implementation acceptance. [IN-08; IN-06, Purpose; IN-07, Inputs Used]

## Governing questions

| Question | Phase 1 answer | Basis |
| --- | --- | --- |
| What does work produce before interpretation? | Either a completed return containing result and error, or a panic captured at the supervised invocation boundary. Neither is a public terminal outcome. | IN-03 §§10–13 |
| What constitutes success? | A completed return with a nil Go error interface. A nil or zero result is still a successful result; Aren performs no task-quality validation. | IN-03 §§11, 13; IN-09, Key Differences 1–2 |
| What constitutes cancellation? | Work has finished, cancellation was accepted, and the returned error matches the retained effective cancellation cause under the rule below. Context cancellation alone is insufficient. | IN-03 §§12–13; IN-14 P3 |
| What constitutes failure? | An unmatched returned error, a captured work panic, or a containable Aren invariant fault. Origin and kind remain explicit. | IN-03 §11; IN-05 §7 |
| Who chooses the terminal state? | Aren's central resolver, using the accepted facts applicable at terminal commitment. Neither cancellation callers nor work choose a terminal state. | IN-03 §§9, 13–14; IN-07 LA-03–LA-04 |
| Does first arrival decide meaning? | No. Serialization determines which requests precede commitment; the resolution policy determines their meaning. A first-write guard only protects an already committed record. | IN-03 §13; IN-13 P1–P2 |
| What may panic containment promise? | Recoverable panics unwinding through the supervised work call become `work/panic`. This is not process isolation or recovery of panics in arbitrary work-spawned goroutines. | IN-03 §§10–11, 19; IN-13 P5 |
| Does terminality prove all resource release? | It proves the supervised invocation has finished and terminal facts are complete. Additional Aren-owned release and join obligations belong to the cancellation/resource area. | IN-07, Downstream obligations; IN-14 P1–P2 |

## Normative constraints

| Constraint | Source | Fixed semantic rule | Open implementation question |
| --- | --- | --- | --- |
| Legal terminal edges | IN-03 §7; IN-07 LA-08 | Only `running -> succeeded/failed/cancelled`; terminal states cannot change | Concrete transition validation and invariant exposure |
| Success payload | IN-03 §§11, 13 | Exact returned result, complete timing, no failure or terminal cancellation payload | Result type and accessor shape |
| Non-success payload | IN-03 §11 | Failed and cancelled outcomes expose no successful result; no partial-result semantics | Private representation of absent payloads |
| Failure classification | IN-03 §11; IN-05 §11.9 | Machine-readable origin and kind distinguish returned error, panic, and Aren invariant fault | Concrete Go vocabulary and accessors |
| Nil-error success after cancellation | IN-03 §13 | Accepted cancellation does not override a completed nil-error return | Integration into realized Sprint 1 resolver |
| Cancellation acknowledgement | IN-03 §§12–13 | Acceptance, work completion, and cause match are all required | Cause capture, signal ownership, and implementation of matching outside lifecycle ownership |
| Stable first reason | IN-03 §12 | Later requests do not replace the accepted reason or effective cause | Explicit-cancel reason parameter and normalization |
| Immutable terminal publication | IN-03 §§14, 17; IN-07 LA-02, LA-07 | State, event, timing, and outcome are complete before waiting returns | Concrete publication and defensive-copy mechanics |
| Result ownership exception | IN-03 §§11, 17; IN-07 LA-09 | Preserve exact successful value without promising deep immutability of arbitrary result objects | Documentation and tests for reference-bearing results |
| Non-recursive invariant handling | IN-05 §7; IN-07 LA-08 | Reject illegal mutation; no recursive finalization, double close, fabricated edge, or replacement outcome | Minimal trustworthy containment path and fault-injection seam |
| Timing | IN-03 §§15, 18; IN-07 LA-10 | Start and finish captured once at lifecycle boundaries; start not after finish; sequence orders events | Time representation and tests, without a speculative clock interface |
| Delivery boundary | IN-04, both implementation waves | Sprint 1 proves success/error/panic and ordinary publication; Sprint 2 extends cancellation and adversarial observation | Exact test and command names |
| Exclusions | IN-03 §§4, 19, 25 | No cleanup hooks, retries, partial results, task validation, persistence, or exactly-once work effects | None within Phase 1 |

## Evidence

### Completion model comparison

All repository behavior in this comparison is reported behavior, not direct observation in this turn.

| Evidence | Completion vocabulary | Finalization owner | Panic treatment | Cancellation treatment | Result validation | Aren relevance |
| --- | --- | --- | --- | --- | --- | --- |
| IN-09, Approach Model A and Pattern 1: Pydantic AI and OpenAI Agents SDK | Narrow typed agent success terminals, with other control outcomes | Agent loop | Cited exception models do not establish a Go panic boundary | Separate stream/interruption mechanisms | Some paths validate structured output or guardrails | Adapt explicit success/failure distinctions; reject model-quality criteria and interruption variants |
| IN-09, Approach Models B–C: Letta and LangGraph | Stop reason/status mapping; graph exhaustion | Application run manager or graph runtime | Framework exception handling, not an Aren panic contract | Product-specific status and interruption mappings | Graph completion need not validate output | Shows why framework completion cannot define Phase 1 success |
| IN-09, Approach Models D–E: CrewAI and OpenHands | Parser finish or externally supplied terminal status | Parser-driven loop or external runtime | Not established for the Aren boundary | External or product-specific paths | Heuristic or limited validation | Negative transfer: do not infer lifecycle truth from text, parser success, or a status string |
| IN-13, Approach Models and P5: Conc | Scoped work finished at owner `Wait` | Scope owner after join | Captures value and stack, then rethrows or exposes an error bridge | Context signal distinct from worker exit | No task validation | Adapt panic diagnostics and completion-before-read; reject waiter-driven finalization and pool reuse |
| IN-13, P1, P3, P5: Temporal SDK | Completion slot interpreted into one close command | Root execution path plus terminal arbiter | Typed recovered panic channel | Request signal does not itself close the workflow | No equivalent arbitrary-result validation requirement | Adapt one interpreter and separate request/return; reject durable commands and its particular priority policy |
| IN-13, Approach Models and P5: Crush | Implicit active state plus terminal publication | Several producers cooperating through a guarded emit path | No recovery around the studied `Run` boundary; dependency behavior remains uncertain | Active entry remains until execution returns | Agent-specific | Negative evidence for scattered terminal guards and an undeclared panic boundary |

Provider completion means a provider response ended. Agent task completion may include validation, tool completion, or loop policy. Framework completion may mean no runnable graph tasks remain. Aren Phase 1 adopts only **lifecycle completion of the supervised invocation**, interpreted according to the PRD. It does not inspect the successful result to determine whether a business task was accomplished. [IN-03 §§4, 10–13; IN-09, Core Thesis and Key Differences 1–2]

### Failure classification evidence

| Failure case | Source treatment | Information preserved | Information lost | Caller behavior enabled | Phase 1 conclusion pressure |
| --- | --- | --- | --- | --- | --- |
| Returned error | IN-10 P3 describes cause-chain classification; IN-14 P4 describes preserving the primary error before cancellation | Wrapped causes and selected classification | String-only conversion loses identity and type information | Inspect cause without parsing presentation text | Retain original local error beneath an Aren-owned classification |
| Panic | IN-13 P5 cites Conc's recovered value/stack and Temporal's typed panic channel | Panic recognition and stack at capture | A first-panic latch may drop secondary panics; generic message conversion loses panic identity | Distinguish programmer failure from ordinary returned error | Capture at the invocation boundary and publish explicit `work/panic` diagnostics |
| Invariant violation | IN-13 P7 describes immediate illegal-transition detection; IN-12, OPA Approach Model, reserves fail-stop behavior for consistency faults | Fault location and runtime origin | Log-and-continue paths can conceal corruption | Stop trusting the operation rather than treating it as routine work failure | Reject before mutation; contain only when a truthful failed transition remains possible |
| Cancellation-related error | IN-14 P3 uses typed context causes; P4 protects causal errors from cancellation artifacts | Why the signal was issued and the original failure | Generic context sentinels alone may erase provenance | Separate cancellation acknowledgement from unrelated failure | Match against accepted cause, not any cancellation-looking error |
| Result plus error | Comparative systems have different partial-output policies; IN-03 §§11, 13 directly settle Aren's case | Returned error and its classification | The result is intentionally not published as success | Reliably distinguish successful from unsuccessful invocation | Suppress result for every non-success outcome; no partial-result channel |

### Evidence limits and conflicts

The error-taxonomy report supports narrow lessons from two TypeScript systems, not a ten-system Go taxonomy consensus. Its recommendations about retryability, wire compatibility, and federated subsystem classifications are not adopted: Phase 1 has neither retries nor those boundaries. [IN-10, Executive Summary, P1–P6, Notable Absences; IN-06, Selection and canonicality assessment]

The Go arbitration report supplies a valuable funnel pattern but also records a production/test divergence: repeated `Complete` overwrites in the cited production path while the test environment ignores repeats. Its descriptions of structural impossibility therefore cannot substitute for testing Aren's real commit guard. The report cites `studies/aren-go-runtime-study/sources/temporal-sdk-go/internal/internal_task_handlers.go:653-657` and `studies/aren-go-runtime-study/sources/temporal-sdk-go/internal/internal_workflow_testsuite.go:1143-1147`. Those files were not independently read here. [IN-13 P1, Anti-Patterns]

Likewise, neither a CAS latch nor cancellation-priority arbitration establishes the PRD's meaning of success. The older phased roadmap's first-winner language is superseded for Phase 1 by nil-error success after accepted cancellation. Docker Agent's reported context recheck before choosing an end reason is not copied as an unconditional cancellation override. [IN-15, Phase 1 Cancellation; IN-03 §13; IN-13 P2; IN-14 P4]

Failure visibility supports one classified fact serving multiple surfaces, but does not require identical payload size everywhere. CrewAI's reported traceback loss and OpenAI Agents SDK's missing stream error variant illustrate why outcome and event assertions must be separate. Temporal's protocol failure and truncation machinery is evidence about information loss, not a reason to introduce serialization or retention infrastructure. [IN-11, Executive Summary, P8, P10, Evidence Index]

No assigned report establishes the exact cause matcher selected below, nor the complete Aren timing and invariant-containment contract. Those are explicit project interpretations requiring local tests, not findings attributed to upstream repositories. [IN-06, Phase 1 decision map; IN-13, Notable Absences; IN-14 P3]

## Candidate outcome models

| Model | Invalid combinations prevented by | Cause preservation | Immutability mechanism | Caller cost | Go-specific risk | Phase fit |
| --- | --- | --- | --- | --- | --- | --- |
| Discriminated terminal outcome | Private discriminant and branch payload; centralized validation | Local error reference within failure branch | Private fields, immutable diagnostic values, defensive collection access | Switch once on terminal state, then inspect the matching branch | Go has no native exhaustive sum type; zero values and internal invalid construction remain possible | Strong semantic model |
| State plus optional fields with validated construction | Private construction enforces exactly the payload allowed by state | Same local cause reference | Same privacy and copy obligations | Common metadata is easy to inspect; accessors must expose availability explicitly | Public fields or unchecked internal literals can create result/failure/cancellation mixtures | Credible smallest concrete realization |
| Separate success, failure, and cancellation types behind a common accessor | Each concrete type contains only its valid branch | Failure and cancellation variants retain appropriate causes | Private concrete data and defensive accessors | Type switches or visitor-style access; shared metadata may repeat | Typed-nil interfaces, incomplete switches, and interface machinery without current value | Possible, but not earned solely to imitate a Rust enum |

The project selects a **validated discriminated terminal contract**, not a mandatory Go sum-type implementation. A private state-plus-fields representation can satisfy it if construction and publication make invalid combinations inaccessible through supported APIs. Separate variant types are acceptable only if sprint reasoning demonstrates that they simplify the actual callers and tests. [IN-03 §11; IN-05 §§11.2, 11.5, 11.10; IN-09 Pattern 1; IN-10 P6]

The valid outcome branches are:

| Terminal state | Successful result | Failure | Terminal cancellation | Common facts |
| --- | --- | --- | --- | --- |
| `succeeded` | Present as the exact returned value, including nil or zero | Absent | Absent | Run identity, start time, finish time |
| `failed` | Unavailable | Present, with valid origin/kind and diagnostics | Absent | Same common facts |
| `cancelled` | Unavailable | Absent as a failure classification | Present, retaining the first accepted reason and effective cause | Same common facts |

Accepted cancellation remains a historical lifecycle fact even when the outcome succeeds or fails. It does not populate a terminal-cancellation branch on those outcomes. A successful result that happens to contain an `error` value is still a result when the separate returned error is nil. Conversely, a non-nil error interface containing a typed nil is still a returned error; Aren must not reinterpret it as success through reflection or formatting. [IN-03 §§11–13; IN-07, Atomicity model; project interpretation of the specified Go return semantics]

Private constructors do not make invalid internal values impossible in Go. The publication boundary must validate its inputs, zero values must never be confused with a completed run, and package-internal tests must challenge invalid construction. A Rust enum analogy does not remove these runtime obligations. [IN-05 §§7, 11.5–11.6; IN-07 LA-08]

## Terminal fact model

A completion record distinguishes **normal return** from **panic unwind**. Merely finding zero result/error storage is not evidence that work returned successfully. Work-local defers are part of the supervised invocation: a panic escaping one of them is a panic completion, not a successful return followed by a second outcome. [IN-03 §§10–13; IN-13 P5; project interpretation of the work boundary]

| Fact | Producer | Recorded when | May race with | Mutable after capture? |
| --- | --- | --- | --- | --- |
| Completion mode: returned or panicked | Aren's invocation boundary | After normal return or during recovery from escaping work panic | Cancellation acceptance | No; mutually exclusive authoritative modes |
| Work result | Work, captured by Aren | On completed normal return | Acceptance and inspection | Captured value is fixed; arbitrary referenced objects remain user-owned |
| Work error | Work, captured by Aren | On completed normal return | Acceptance | Error reference is fixed; its implementation must remain safe for read-only use |
| Panic diagnostic | Aren's work boundary | While the work stack is available during unwind | Acceptance | Aren-owned type/text/stack snapshot is immutable; no unrestricted raw mutable panic payload is required |
| Accepted cancellation and effective cause | Private lifecycle authority | At the first accepted request | Return, panic, another request, terminal commitment | Never replaced after acceptance |
| Cause-match result | Aren's classification preparation | Against one specific captured error and accepted cause | First acceptance if not yet captured | Valid only for those facts; not an independent terminal decision |
| Existing terminal commitment | Private lifecycle authority | At complete terminal publication | Late finalization attempts and cancellation | Immutable |
| Detected invariant fault | Aren's guard or internal validation boundary | At detection, before applying the invalid mutation | Work completion and acceptance if work is still active | Diagnostic remains stable; detection does not itself authorize terminal publication |
| Start time | Lifecycle authority | At `created -> running` commitment | Readers | Immutable |
| Finish time | Lifecycle authority | At the terminal commitment | Readers and waiters | Immutable and absent before commitment |

Errors and custom causes are not universally cloneable values. Preserving their original Go identity while guaranteeing deep immutability of every object reachable through them is not a promise Aren can honestly make. Aren must retain a stable classification and immutable diagnostic snapshot, while documenting the original error/cause as a borrowed, read-only diagnostic object. Mutating such an object after handing it to Aren is outside the supported sharing contract. This is an explicit ownership limitation, not permission to expose Aren-owned slices or metadata. [IN-03 §11; IN-05 §§11.2, 11.9; IN-07 LA-09 and Self-critique]

## Resolution matrix

The matrix applies to a normally functioning authority with valid identity, history, and start timing. `Any` includes nil or zero result values. `Matching` means the exact predicate in the next section. Rows are ordered: an existing terminal commitment is protected before any later candidate is considered. [IN-03 §§7, 11–14; IN-07 LA-04, LA-08]

| Result | Error | Panic | Cancellation accepted | Error matches run cancellation cause | Existing terminal | Candidate outcome | Reason |
| --- | --- | --- | --- | --- | --- | --- | --- |
| Any | Any | Any | Any | Any | Yes | Retain existing outcome; reject later internal finalization | Terminal truth is immutable; no new timestamp, sequence, or event |
| Any, including nil | Nil | No; normal return confirmed | No | Not applicable | No | `succeeded` | Nil-error return defines success |
| Any, including nil | Nil | No; normal return confirmed | Yes | Not applicable | No | `succeeded` | Acceptance is not a cancellation override |
| Any | Non-nil | No; normal return confirmed | No | Any apparent match | No | `failed`, `work/returned_error` | No accepted cancellation exists; suppress result |
| Any | Non-nil | No; normal return confirmed | Yes | Yes | No | `cancelled` | Accepted cancellation plus matching return; suppress result |
| Any | Non-nil | No; normal return confirmed | Yes | No | No | `failed`, `work/returned_error` | Unrelated return error remains failure |
| Unavailable | No completed return | Captured escaping work panic | No | Not applicable | No | `failed`, `work/panic` | Panic is not an ordinary error return |
| Unavailable | No completed return | Captured escaping work panic | Yes | Even if panic value is the cause | No | `failed`, `work/panic` | Panic is never cancellation acknowledgement |
| None yet | None yet | None yet; work active | Yes | Not applicable | No | No candidate; remain `running` | A request does not prove work stopped |
| Any | Any | Both authoritative normal-return and panic modes asserted | Any | Any | No | Aren invariant fault | Capture representation is contradictory; ordinary work cannot author both completion modes |
| Any | Any | Valid completed mode | Yes, but effective cause missing or replaced | Undefined | No | Aren invariant fault | Accepted-cause record violates its own contract |
| Any | Any | No confirmed completion mode | Any | Any | No | Reject terminal proposal as Aren invariant fault | Zero storage or a signal is not completion evidence |
| Any | Any | Valid completed mode | Any | Any | No, but state is not `running` or history/timing is invalid | Reject mutation; apply containment boundary below | Resolver cannot invent a legal source state or repair history by publishing a plausible outcome |

A non-nil result together with a non-nil error is not an invariant violation. It is a legal Go return combination with no successful-result semantics in Aren. It becomes either `cancelled` or `failed` according to acceptance and cause matching. [IN-03 §§11, 13]

A work function returning an error whose dynamic type resembles an Aren failure does not gain lifecycle authority. For this invocation, it is still a returned work error, wrapped with `origin=work`, `kind=returned_error`, unless it matches the accepted cause. Its original structure remains inspectable underneath. Only Aren's internal fault path may originate this run's `aren/invariant_violation` classification. [IN-03 §§9–11; IN-07 LA-03, LA-08; project conclusion]

### Cause-matching rule

The project interpretation of “matches the run context cancellation cause” is:

```text
accepted cancellation exists
AND work returned a non-nil error
AND errors.Is(returnedError, acceptedEffectiveCause)
```

`acceptedEffectiveCause` is the non-nil cause retained by Aren's first acceptance and used to cancel the work context. Matching is directional, from returned error to accepted cause. It is not string equality, arbitrary comparison of two messages, or a check that either context happens to be cancelled. [IN-03 §§12–13; IN-14 P3; project decision, not a studied Aren implementation]

| Case | Interpretation |
| --- | --- |
| Returned error is the accepted cause | Matches |
| Returned error wraps the accepted cause using Go error wrapping | Matches |
| Returned error has identical text but no matching identity or `Is` relationship | Does not match |
| Work returns `context.Canceled`, and that is the effective accepted cause | Matches |
| Work returns `context.Canceled`, but the effective accepted cause is a distinct custom error | Does not match unless the returned error explicitly matches that custom cause |
| Work returns the accepted parent's `context.DeadlineExceeded` cause, directly or wrapped | Matches; terminal state is `cancelled`, not a new timeout state |
| Work returns `context.DeadlineExceeded` from a work-owned timeout unrelated to accepted run cancellation | Does not match merely because it is a timeout |
| Work panics with the accepted cause as its panic value | `work/panic`, not cancellation |
| Returned joined error contains the accepted cause and an unrelated error | Matches under `errors.Is`; retain the complete returned error as acknowledgement diagnostics |

This deliberately does **not** add a blanket `errors.Is(err, ctx.Err())` fallback. With custom causes, `ctx.Err()` may be a coarse sentinel rather than the effective cause. Treating that sentinel as sufficient would broaden the PRD's cause-based rule and could classify a different cancellation source as the run's cancellation. Work that needs reliable custom-cause acknowledgement should return or wrap `context.Cause(ctx)`. Sprint 2 must document and test this behavior explicitly. [IN-03 §13; IN-14 P3; project decision]

The joined-error case is an explicit trade-off. Under this matcher, an error containing the accepted cause is not “any other returned error,” even if it contains another branch. Aren preserves that branch for inspection rather than changing the terminal state to failure. If product review instead requires every branch to be cancellation-related, that is a different matching policy requiring an explicit superseding decision, not a hidden special case in implementation. [IN-03 §13; IN-05 §11.9; project decision]

The predicate identifies declared error equivalence, not physical causation. Shared standard sentinels and custom `Is` methods cannot prove which operation actually produced an error. Aren supervises cooperative in-process code; it does not authenticate work's explanation for stopping. [IN-03 §§10, 19; IN-07, Authority and ownership]

### Matching and publication

Go error inspection can invoke user-defined `Is`, `Unwrap`, or formatting methods. Such code must not run while lifecycle ownership is held. Nevertheless, a match computed before acceptance cannot be published later as though it described the final committed fact set. [IN-07, Timing and notification; IN-05 §7]

The semantic implementation obligation is:

1. Capture the work completion once.
2. Obtain the current accepted-cause fact without exposing mutable lifecycle state.
3. Perform user-defined error inspection outside lifecycle ownership.
4. Re-enter authority and validate that the inspected cause is still the applicable fact before resolving and committing.
5. If first acceptance occurred in the gap, inspect that cause outside ownership and revalidate again.

Accepted cancellation moves only from absent to present and is then immutable. This supplies a bounded lifecycle revalidation problem, not a new retry policy, callback registry, or independently writable classification cache. The sprint may realize the obligation differently, but cannot publish a stale classification or move error callbacks under the commit lock. [IN-03 §§12–14; IN-07, Linearization and terminal interpretation; project derivation]

Error and cause methods must obey the supported read-only diagnostic contract: terminate, remain stable, and not panic. Aren cannot bound an arbitrary blocking `Is` or `Error` method without supervising additional user work. A method fault must not be converted into nil-error success or silently treated as a clean non-match. The concrete diagnostic-fault boundary must preserve visibility without blaming it on an ordinary work return; no forceful-stop or sandbox guarantee is added. [IN-03 §§4, 19; IN-05 §§7, 11.9; IN-07, Risks]

## Information preservation

### Surface contract

| Information | Local outcome and error access | Canonical lifecycle event | Diagnostic rendering |
| --- | --- | --- | --- |
| Run identity, state, start and finish timing | Complete immutable terminal record | Matching identity, event type, transition, sequence, and occurrence time | Render committed facts, never independently infer a terminal state |
| Successful result | Exact returned value; no deep result-object freeze | Success fact; no requirement to duplicate arbitrary result into the event | Show the result only through success access |
| Returned failure | `work/returned_error`, stable message, original error available for `errors.Is`/`errors.As` | Same origin/kind and stable diagnostic summary | May shorten display, but must not replace structured classification |
| Work panic | `work/panic`, captured dynamic type description, stable value description, and captured stack | Panic-recognizable failure payload with immutable diagnostics, including stack access | Clearly label panic; stack can be shown separately from the compact event line |
| Aren invariant fault | `aren/invariant_violation` with operation, relevant state, and invariant description where safely available | `run.failed` only if truthful containment can commit it | An uncontained fault is reported as an invariant defect, not a fabricated run outcome |
| Accepted cancellation | First reason and effective cause remain in run-local cancellation facts | One request event with stable first-reason metadata | Distinguish request acceptance from terminal cancellation |
| Terminal cancellation | Cancellation metadata plus original returned acknowledgement error available locally | Matching cancellation metadata and stable acknowledgement summary | Do not present cancellation as proof that no effects occurred |
| Original custom error object or arbitrary panic object graph | Error causes retained read-only where practical; no mandatory unrestricted raw panic-object accessor | No requirement to expose mutable object graphs | Rendering is a projection, not a serialized equivalent of the original object |

This is a semantic minimum, not an event field layout. The events area must realize the information contract using immutable lifecycle payloads and explicit access to any larger diagnostic fields. A compact CLI line need not print a full stack, but the stack must not be destroyed merely because the first display is compact. [IN-03 §§11, 15, 20; IN-05 §§7, 11.2, 11.9; IN-11 P1, P8, P10; project conclusion]

Go wrapping preserves local error identity and inspection behavior. Explicit Aren fields preserve origin, kind, message, and panic/invariant distinctions. Neither replaces the other: `%w` alone does not classify origin, while a message string alone does not preserve `errors.Is` or `errors.As`. [IN-05 §11.9; IN-10 P3; IN-11 P10]

A terminal event is not a lossless substitute for the entire outcome. An event-only consumer need not retain the arbitrary successful result, original Go error identity, or full wrapper graph. In particular, the terminal event alone cannot reconstruct the full earlier cancellation-request history of a successful run. Canonical retained history and local outcome access remain complementary, with no persistence promise. [IN-03 §§15–19; IN-07 LA-05; project conclusion]

### Panic and invariant boundaries

A recoverable panic escaping the supervised work call, including a work defer, becomes `work/panic`. A panic recovered internally by work followed by a normal nil-error return is success from Aren's boundary. Recovery must be narrowly placed so it cannot accidentally relabel Aren's resolver or commit defects as work panic. Panic detection must also avoid equating a nil recovered value with a completed normal return; supported toolchain behavior for `panic(nil)` requires a regression test. [IN-03 §§10–13; IN-05 §7; IN-13 P5; project interpretation]

Aren does not promise to contain process exit, unrecoverable runtime failures, or panics in goroutines created independently by work. `runtime.Goexit` is not a normal return or a recoverable panic and must never become success merely because result/error slots remain zero. Sprint 1 must make that abnormal-exit boundary explicit without introducing general execution isolation. [IN-03 §§10, 19; IN-13 P5; project boundary]

Invariant handling follows the publication dependency rather than inventing an emergency state graph:

| Detection boundary | Required behavior |
| --- | --- |
| Trustworthy `running` run, supervised work finished, candidate invalid before publication | Reject the invalid candidate; construct one minimal `aren/invariant_violation` failure through the same terminal publication owner, without re-entering the failed ordinary transition path |
| Trustworthy `running` run, work still active | Reject the invalid mutation and expose the invariant diagnostic; do not publish terminal failure or cancellation as proof of stopped work. Any containable terminal fault waits for actual work completion |
| Existing terminal commitment | Preserve it byte-for-byte at the lifecycle level; reject and expose the later internal fault without another lifecycle event or close |
| Before a legal start, or after authority/history can no longer be trusted | Fail fast through an explicit internal invariant-fault boundary; do not fabricate `created -> failed`, invent timing, or recursively attempt repair |

The fail-fast path is an explicitly identifiable internal defect, not an ordinary work error and not a promise of host-process recovery. A guarded internal panic is a permissible concrete mechanism; silent logging followed by continuing the illegal mutation is not. Fault injection must demonstrate the distinction between a containable candidate-construction defect and corruption beyond truthful recovery. [IN-03 §§7, 11–14; IN-05 §7; IN-07, Counterexample analysis and LA-08; IN-12, OPA Approach Model]

### Timing

Start time is captured at the start commitment; finish time is captured at the terminal commitment. The outcome and corresponding terminal event use that one committed finish fact. It is not recomputed per waiter, copied from cancellation acceptance, or advertised as the exact instruction at which user work returned. [IN-03 §§15, 18; IN-07 LA-10]

Equal start and finish times are valid. Start-after-finish is not. Sequence remains the event ordering authority, including when clock readings coincide. No public duration field or clock abstraction is selected here; if duration is exposed later in sprint reasoning, it must derive from the same committed timing rather than observer receipt time. [IN-03 §18; IN-07, Timing and notification]

## Project conclusions

| Conclusion | Evidence basis | Applies in Sprint 1 | Extended in Sprint 2 | Reopen trigger |
| --- | --- | --- | --- | --- |
| OF-01: Use one validated discriminated terminal contract, independently of its concrete Go representation | IN-03 §11; IN-05 §§11.2, 11.5; IN-09 Pattern 1 | Success/error/panic branches, common timing, invalid-construction tests | Cancellation branch | Realized representation permits invalid public combinations or imposes avoidable caller complexity |
| OF-02: Distinguish captured return/panic facts from terminal outcomes | IN-03 §§10–14; IN-07 LA-03–LA-04; IN-13 P1 | One invocation capture and one terminal interpreter | Interpret accepted cancellation at commitment | A candidate bypasses authority or a zero capture becomes success |
| OF-03: Nil-error return succeeds, including nil/zero results and success after accepted cancellation | IN-03 §§11, 13 | Full ordinary success contract | Acceptance-before-success case | Only a governed product revision can change this policy |
| OF-04: Every non-success outcome suppresses the work result; result plus error is legal but not partial success | IN-03 §§11, 13 | Returned failure and panic | Cause-matching cancellation with a returned result | A later phase explicitly earns partial-result semantics |
| OF-05: Origin and kind describe this invocation's boundary, not arbitrary labels supplied by work | IN-03 §§9–11; IN-05 §11.9 | `work/returned_error`, `work/panic`, and containable `aren/invariant_violation` | Cancellation remains a separate outcome | A concrete new failure boundary is introduced through governed scope |
| OF-06: Match returned errors directionally with `errors.Is` against the first accepted effective cause; no generic sentinel fallback | IN-03 §§12–13; IN-14 P3; project interpretation | Preserve error identity and keep classifier seams private | Full custom-cause, deadline, wrapper, join, and mismatch matrix | A concrete integration demonstrates that this interpretation is insufficient; change must be explicit |
| OF-07: Retain original local causes plus stable immutable diagnostic projections | IN-05 §§7, 11.2, 11.9; IN-10 P3; IN-11 P10 | Error inspection and panic stack; read-only cause ownership | Accepted cause and acknowledgement error | Cause loss, mutable lifecycle aliases, or a genuinely introduced serialization boundary |
| OF-08: Capture work panics narrowly; do not recover Aren defects as work failures | IN-03 §11; IN-13 P5; IN-05 §7 | Panic boundary and invariant fault injection | Panic after accepted cancellation | A new callback or execution boundary changes actual fault ownership |
| OF-09: Invariant containment never invents an edge, overwrites a terminal, or claims active work has stopped | IN-07 LA-08; IN-05 §7; IN-12, OPA Approach Model | Complete containable/non-containable boundary | Collision tests preserve the same boundary | Realized fault path cannot meet truthful publication preconditions |
| OF-10: Complete state/event/outcome/timing precedes waiter return; one terminal guard is separate from resolution policy | IN-03 §§13–18; IN-07 LA-02, LA-04, LA-07; IN-13 P1–P2 | Multiple waiters and complete retained history | Cancellation and observer collisions | Any torn view, stale classification, duplicate terminal, or lost diagnostic |
| OF-11: Preserve Sprint 1 decisions explicitly rather than redesigning its outcome model in Sprint 2 | IN-04, Cross-Sprint Carry-Forward Rule | Record realized API and tests | Preserve, extend, or explicitly supersede each relevant decision | New evidence proves a named earlier decision incorrect or insufficient |

These are project reasoning conclusions and proof obligations. No tests, race runs, or implementation review were performed in this area stage. Historical prototype passes do not satisfy either sprint's acceptance. [IN-01, Project Reasoning Policy; IN-04; IN-05 §§4, 15]

## Trade-Offs

**Closed vocabulary versus forward evolution.** A small validated origin/kind vocabulary makes Phase 1 switches and tests auditable. It requires deliberate revision when new execution types arrive. No persisted unknown-value compatibility layer is needed now, and internal invalid values must not silently become a generic successful or unknown outcome. [IN-03 §§4, 11; IN-05 §11.5; IN-10 P5–P6]

**Defensive copies versus result ergonomics.** Copying Aren-owned stacks, event payloads, and metadata protects concurrent observers. Copying arbitrary successful result graphs could change identity, fail on unsupported values, or imply impossible universal freezing. The contract keeps result identity and assigns mutable result use to the caller. Original error causes similarly require an explicit read-only sharing obligation. [IN-03 §§11, 17; IN-05 §11.2; IN-07 LA-09]

**Panic containment versus programmer-fault visibility.** A classified work panic protects the run boundary while retaining evidence. A catch-all recovery around the entire runtime could disguise internal corruption. Narrow recovery plus explicit invariant escalation costs more care but preserves meaningful origin. [IN-03 §11; IN-05 §7; IN-13 P5, P7]

**Cause fidelity versus stable serialization.** Local wrapping retains Go identity and type inspection. An immutable textual projection is easier to render but loses those capabilities. Phase 1 retains both where applicable and builds no serializer; any later wire or durable representation must declare its losses separately. [IN-05 §11.9; IN-11 P8, P10; IN-03 §4]

**Strict cause matching versus familiar `ctx.Err()` returns.** Strict matching avoids automatically accepting a coarse sentinel for a custom cause. It makes `context.Cause(ctx)` important for custom-cause-aware work. Adding an implicit sentinel fallback would be more permissive, but weakens attribution and changes the stated rule. [IN-03 §13; IN-14 P3; OF-06]

**Joined-error fidelity versus failure precedence.** Standard `errors.Is` matching treats a joined cause as acknowledgement even alongside another error. Keeping the complete returned error prevents information loss, but the terminal state alone does not advertise every branch. Failure-on-any-unrelated-branch would be a distinct policy and is not silently added. [IN-03 §13; IN-05 §11.9; OF-06–OF-07]

**Narrow failure kinds versus premature taxonomy growth.** Returned error, panic, and invariant violation cover current handling differences. Retryability, provider categories, cleanup failure policy, and transport codes would add names without Phase 1 consumers. Those remain deferred. [IN-03 §§4, 11, 25; IN-10 P1–P2; IN-12, Approach Models]

## Risks

| Risk | Consequence | Control |
| --- | --- | --- |
| Result/error incoherence | A failed or cancelled run exposes a successful result | Validate branch payloads at publication; test non-nil result plus every non-success path. [IN-03 §11] |
| Mutable nested lifecycle data | One observer changes another's outcome, event, or stack | Private storage, defensive access, mutation attacks; do not confuse user-result exceptions with lifecycle payloads. [IN-05 §7; IN-07 LA-09] |
| Mutable custom error or cause | Message or matching behavior changes after capture; concurrent reads may race | Stable diagnostic snapshots and explicit read-only cause contract; no deep-freeze claim. [IN-05 §11.9; OF-07] |
| Panic disguised as ordinary returned error | Caller cannot distinguish work panic from expected failure | Explicit panic mode, kind, stack, and narrow recovery boundary. [IN-03 §11; IN-13 P5] |
| Cancellation inferred from `context.Canceled` alone | Unaccepted or different cancellation is mislabeled | Require acceptance and effective-cause matching; test same-text and custom-cause mismatches. [IN-03 §13; OF-06] |
| Shared sentinel or custom `Is` overstates causation | Matching succeeds without proving the physical cause of the return | Describe matching as cooperative declaration, not authenticated provenance. [IN-03 §19; OF-06] |
| Cause inspection under lifecycle ownership | Reentrancy or blocking deadlocks readers and cancellation | Inspect outside ownership and revalidate accepted facts before commitment. [IN-07, Timing and notification] |
| Stale classification | A cause accepted after preliminary inspection is ignored at commit | Monotonic accepted-fact revalidation and controlled gap tests. [IN-07, Linearization and terminal interpretation] |
| Recursive invariant fallback | Duplicate close, terminal replacement, or incomplete timing | Minimal non-recursive containment with explicit trustworthy-state preconditions. [IN-05 §7] |
| Diagnostic rendering changes semantics | Formatting panic or string comparison changes outcome or erases cause | Stable typed classification; no string-based matching; rendering is not resolution. [IN-10 P9; IN-11 P10] |
| Bounded event count mistaken for bounded diagnostic bytes | Large error text, stacks, or results produce misleading resource claims | Declare retained diagnostic behavior and measure representative failure workloads; any truncation must be explicit. [IN-05 §§11.7–11.8; IN-07, Minimum coherent fact set] |
| Future taxonomy or serialization enters early | Phase 1 grows into a provider/runtime framework | Require a current behavioral consumer for each field and abstraction. [IN-03 §§4, 25] |
| Report or prototype confidence exceeds inspection | Prose is accepted as proof of the new runtime | Fresh real-runtime tests, independent assertions, and required reviews. [IN-06; IN-05 §§4, 15] |
| Additional reasoning standard unchecked | Project review overlooks missing method requirements | Inspect the named reasoning README before approval. [IN-08] |

## Verification obligations

These are required test intentions, not reported execution results. The verification area must map them to the realized API and independent oracle. Tests must challenge the production capture/resolution/publication path, not only a mock that already enforces the desired policy. [IN-03 §21; IN-05 §11.6; IN-13, Anti-Patterns]

| Semantic claim | Minimal test | Adversarial or negative test | Failure the test must detect |
| --- | --- | --- | --- |
| Nil-error return succeeds | Return ordinary, nil, and zero values | Return an error-shaped object as the successful result; seed a nonempty-result requirement | Success inferred from result shape rather than returned error |
| Result plus error is not partial success | Return non-nil result and sentinel error | Repeat with wrapped matching cancellation cause | Leaked successful result on non-success |
| Go error nilness is respected | Return a non-nil error interface containing a typed nil | Error formatter panics; verify no conversion to success | Reflection or formatting changes error semantics |
| Original cause remains inspectable | Assert `errors.Is` and `errors.As` through failure access | Several wrapping layers and distinct same-message errors | String flattening or mistaken equality |
| Work cannot spoof runtime origin | Work returns an error carrying Aren-like classification | Return a previously obtained structured failure as work error | Foreign labels grant internal fault authority |
| Work panic is contained and recognizable | Panic with string and error values; assert stack | Panic in a work defer; panic with accepted cause; `panic(nil)` | Panic escape, lost stack, or false cancellation/success |
| Abnormal non-return is not success | Exercise the supported abnormal-exit boundary | `runtime.Goexit` in a controlled test that cleans up its harness | Zero completion storage treated as normal return |
| Accepted nil-error completion succeeds | Accept cancellation, then release work to return nil error | Repeat under concurrent readers and waiters | Cancellation-priority override |
| Cause match is explicit | Direct and wrapped standard/custom cause | Same-text error, unrelated timeout, coarse sentinel after custom cause | Cancellation without accepted provenance |
| Joined matching error preserves diagnostics | Join accepted cause and unrelated error | Assert both remain locally inspectable while outcome is cancelled | Silent branch loss or an undocumented alternative precedence |
| Parent deadlines do not add a terminal state | Already-expired and controlled expiring parent | Work returns unrelated error or nil after deadline acceptance | Automatic timeout/failure/cancellation override |
| First accepted reason is stable | Repeat explicit requests with distinct reasons | Race explicit and parent requests with distinct custom causes | Replacement cause or mismatched request/outcome metadata |
| Matching cannot block lifecycle ownership | Custom `Is` inspects run state or is held at a test barrier | Concurrent cancellation and readers proceed while inspection is held | User callback runs under the transition lock |
| Preliminary matching is revalidated | Pause after reading “not accepted,” then accept | Seed publication of the stale failure candidate | Committed acceptance ignored by terminal interpretation |
| Active cancelled work remains running | Hold work after it observes the signal | Assert no terminal event or waiter return until controlled release | Request confused with stop |
| Terminal timing is complete | All terminal branches have identical committed timing across waiters | Seed missing finish time and finish-before-start | Plausible state with invalid timing |
| Publication is coherent | Returned waiters immediately inspect outcome and terminal history | Inject pauses between attempted field writes | Early waiter release or torn fact set |
| Terminal commitment is immutable | Attempt duplicate internal finalization | Race conflicting candidates against the real owner | Replacement outcome, sequence append, or double close |
| Invariant handling is truthful | Inject invalid candidate after work stops | Inject before start, while work is active, and after terminal | Fabricated edge, premature stop claim, or recursive recovery |
| Lifecycle diagnostics are immutable | Mutate returned event/stack snapshots | One waiter mutates a snapshot while another reads | Shared writable Aren-owned storage |
| Instruments can fail | Corrupt a fixture or seed one historical defect at a time | Remove guard, flatten cause, drop stack, or release waiters early | Vacuous tests or an oracle repeating production logic |

Sprint 1 owns ordinary result/error/panic behavior, timing, invariant containment, defensive publication, and concurrent waiting under `go test -race`. Sprint 2 adds the complete acceptance/matching matrix, explicit/parent collisions, repeated stress, observer interactions, and cause-aware diagnostic CLI evidence. The full Phase 1 suite must run under the race detector; stress supplements controlled schedules rather than replacing them. [IN-04, both implementation waves; IN-03 §§20–24; IN-05 §11.6]

Neither race detection nor a passing truth-table unit test proves resource release, absence of deadlock, or the real publication path. The cancellation and verification areas must join these obligations with owned-resource accounting and controlled work release. [IN-07, Downstream obligations; IN-14 P1, P10]

## Self-critique

**Which invalid outcome can still be constructed inside the package?** A zero value, a terminal state with missing finish time, a failed branch retaining a result, or a cancellation branch without a first accepted cause may still be constructible internally in Go. Privacy protects callers, not every internal author. Publication validation and negative tests remain necessary even if constructors are the intended path. [IN-05 §§7, 11.5–11.6; OF-01]

**Which failure fact would be lost if the process only retained the terminal event?** Under the minimum event projection, original Go cause identity, its full wrapper graph, and arbitrary raw panic-object structure may be unavailable. Panic recognition and stack diagnostics must survive in the failure event contract, but an event is not a serialized replacement for every outcome field. Phase 1 retains both history and outcome in memory and promises neither after process loss. [IN-03 §§15–19; Information preservation]

**Could unrelated work errors be misclassified as cancellation?** Yes, if a custom `Is` method claims equivalence, a shared sentinel is returned for another cancellation, or a joined error contains the accepted cause alongside another failure. The chosen rule is precise and testable but establishes declared error matching, not physical causation. The strict rejection of a coarse sentinel after a distinct custom cause reduces one ambiguity without pretending to eliminate all of them. [OF-06; Resolution matrix]

**Does custom-cause matching impose a usability cost?** Yes. Work returning only `ctx.Err()` may fail rather than cancel when the run uses a distinct custom cause. This document chooses the PRD's cause-based interpretation explicitly, not silently. Sprint 2 documentation and tests must make the consequence visible; changing it requires an explicit decision about the weaker provenance guarantee. [IN-03 §13; OF-06]

**Which conclusion assumes a future serialization boundary?** None requires one. Stable diagnostic fields serve current events and CLI rendering. Discussion of lost Go identity explains why future serialization would need its own contract; it does not authorize a schema version, wire envelope, durable failure graph, or compatibility layer now. [IN-03 §4; OF-07]

**Is arbitrary diagnostic code fully contained?** No. Avoiding error callbacks under lifecycle ownership prevents one serious deadlock class, but it does not make hostile or nonterminating error methods safe. The supported contract requires stable, terminating, non-panicking inspection. A broader adversarial-code containment promise would be a different execution boundary and is not smuggled into this phase. [IN-03 §19; IN-07, Timing and notification]

**What remains unproven?** Every Aren-specific test obligation remains unexecuted here. Exact Go representation, diagnostic-fault mechanics, non-recursive containment implementation, and accepted-fact revalidation must survive sprint reasoning and real-runtime tests. The reasoning README also remains unchecked. These limitations block claims of implementation correctness or a passing project review, not the explicit semantic decisions recorded above. [IN-01, Project Reasoning Policy; IN-04; IN-08]
