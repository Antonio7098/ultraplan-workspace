# Verification And Go Correctness

> Project: `aren-phase-01-execution-lifecycle`
> Area: Verification And Go Correctness
> Scope: Cross-sprint verification of the Phase 1 lifecycle, publication, cancellation, observation, and Aren-owned resource obligations
> Evidence boundary: Supplied governed documents and report excerpts; no direct repository inspection or test execution
> Status: Verification decisions and proof obligations, not implementation acceptance or a project-review verdict

## Purpose

Define how Phase 1 correctness claims will be challenged by independent oracles, controlled execution schedules, deliberately broken cases, race detection, and resource-release evidence. The verification target is the realized Aren runtime, not a test double that already implements the intended lifecycle. [IN-01 §§21–24; IN-03 §§7, 11.6; IN-15, Pattern Catalog P1]

The selected approach is a small reference model plus exact lifecycle-trajectory assertions, supplemented by controlled concurrency tests and ownership checks. Pure tests establish interpretation of facts; concurrent tests establish publication and progress; leak tests establish release of actual support resources. None substitutes for the others. [IN-03 §§11.1–11.8; IN-17, Core Thesis and Pattern Catalog 4–7; IN-07–IN-10]

This area preserves the dependency decisions. It does not select Go packages, public signatures, synchronization primitives, an observer implementation, or a generic test framework. Sprint reasoning must turn the obligations below into exact test names, production-path seams, commands, and evidence records. [IN-02, Canonical Project Flow; IN-06, Reasoning Areas]

## Inputs Used

Repeated supplied copies of a path are one input. For excerpted documents, only the supplied portions were inspected. Repository locations quoted in reports remain second-hand evidence; they are not additional inspected inputs.

| Input ID | Input | Kind | Exact workspace-relative path | Exact portion inspected | How it was used | Limits or applicability warning |
| --- | --- | --- | --- | --- | --- | --- |
| IN-01 | Phase 1 PRD | Requirement | `projects/aren-phase-01-execution-lifecycle/docs/PRD.md` | Full supplied copy, especially §§7–24 | Authoritative behavior and acceptance coverage | Requirements are not implementation proof |
| IN-02 | Project roadmap | Requirement and delivery policy | `projects/aren-phase-01-execution-lifecycle/roadmap.md` | Full supplied copy; Canonical Project Flow, Cross-Sprint Carry-Forward Rule, both implementation waves, Phase Exit Gate | Sprint-specific gates and carry-forward | Commands are planned, not executed |
| IN-03 | Accepted Go decision | Requirement with historical evidence | `projects/aren-phase-01-execution-lifecycle/docs/final-language-decision.md` | Full initial supplied copy, especially §§4, 6–8, 11–15 | Mandatory verification, ownership, measurement, and prototype regression seeds | Historical dirty-tree evidence; no reproduction here |
| IN-04 | Project lineage | Historical context and doctrine | `projects/aren-phase-01-execution-lifecycle/docs/project-lineage.md` | Full initial supplied copy; Evidence-Grounded Planning, §§3.2–3.4, 3.10, 5.2–5.6 | Proportional verification and simplification | Historical language-open clause does not reopen Go |
| IN-05 | Project index | Governance | `projects/aren-phase-01-execution-lifecycle/project-index.md` | Project Reasoning Policy, Project Scope, unavailable-dimension table, Prior Decisions, carry-forward rules | Required passing review and evidence restrictions | Catalog entries are not evidence of execution |
| IN-06 | Project-reasoning index | Governance | `projects/aren-phase-01-execution-lifecycle/project-reasoning/index.md` | Reasoning Areas; verification assignments; Source Document Assignments; Excluded Evidence | Dependencies and decision routing | Assignment text is not a report finding |
| IN-07 | Evidence Assessment And Routing | Dependency reasoning | `projects/aren-phase-01-execution-lifecycle/project-reasoning/areas/00-evidence-map.md` | Supplied Purpose, Inputs Used, Scope, assessment tables, evidence tail, Risks, Self-critique | Source precedence, report restrictions, negative transfer | Excerpt; omitted assessment material not inspected |
| IN-08 | Lifecycle Authority And Atomic Publication | Dependency reasoning | `projects/aren-phase-01-execution-lifecycle/project-reasoning/areas/01-lifecycle-authority-and-atomic-publication.md` | Supplied Purpose, Governing questions, Normative constraints, Evidence, ending Trade-Offs, Risks, Downstream obligations, Self-critique | Private authority, coherent reads, finality, capability and mutation attacks | Semantic dependency, not executed proof |
| IN-09 | Outcomes Failures And Terminal Resolution | Dependency reasoning | `projects/aren-phase-01-execution-lifecycle/project-reasoning/areas/02-outcomes-failures-and-terminal-resolution.md` | Supplied Purpose, Governing questions, Normative constraints, Evidence, Candidate outcome models, beginning of Terminal fact model, Risks, Verification obligations, Self-critique | Outcome branches, strict cause matching, diagnostics, abnormal exits, containment | Excerpt; detailed implementation remains open |
| IN-10 | Cancellation Goroutine Ownership And Cleanup | Dependency reasoning | `projects/aren-phase-01-execution-lifecycle/project-reasoning/areas/03-cancellation-goroutine-ownership-and-cleanup.md` | Supplied Purpose, Governing questions, Normative constraints, Evidence, beginning of Cancellation state model, Verification obligations, Self-critique | Acceptance-to-signal ordering, terminal eligibility, startup races, resource release | Excerpt; concrete parent integration not selected |
| IN-11 | Events Observation Waiting And Replay | Dependency reasoning | `projects/aren-phase-01-execution-lifecycle/project-reasoning/areas/04-events-observation-waiting-and-replay.md` | Supplied Purpose, Governing questions, Normative constraints, Evidence, beginning of Canonical record versus delivery, Verification obligations, Self-critique | Inclusive cursors, tail behavior, terminal drain, local abort, observer independence | Excerpt; reader implementation not selected |
| IN-12 | Selected verification template | Reasoning requirement | `projects/aren-phase-01-execution-lifecycle/reasoning/project-wide/05-verification-and-go-correctness.md` | Full supplied copy | Required specialist sections and falsification questions | Referenced reasoning README not supplied |
| IN-13 | Mutation discipline and transitions | Study report | `studies/agent-harness-study/reports/final/02.04-mutation-discipline-and-state-transitions.md` | Executive Summary; Approach Models; Pattern Catalog; visible Key Differences; Open Questions; Evidence Index | Independent legal/illegal matrix; distinguish a funnel from enforced guards | Report-only; durability, ORM, and distributed locking excluded |
| IN-14 | Trajectory evaluation | Study report | `studies/agent-harness-study/reports/final/18.02-trajectory-evaluation.md` | Executive Summary; Core Thesis; Approach Models; Pattern Catalog P1–P4 and visible beginning of P5; supplied citation tail | Structured whole-run records as independent test input | Report-only; no LLM scoring, partial-credit scoring, or timestamp ordering adopted |
| IN-15 | Go lifecycle ownership and arbitration | Study report | `studies/aren-go-runtime-study/reports/final/01.01-lifecycle-transition-ownership-and-terminal-arbitration.md` | Executive Summary; Approach Models; Pattern Catalog P1–P8; Key Differences; supplied failure-mode tail; Notable Absences; Open Questions; Evidence Index | Duplicate finalization, early publication, panic boundaries, production/mock divergence | Report-only; authors lacked Aren PRD; no complete Aren timing proof |
| IN-16 | Go cancellation and ownership | Study report | `studies/aren-go-runtime-study/reports/final/01.02-cancellation-goroutine-ownership-and-cleanup.md` | Executive Summary; Core Thesis; Approach Models; Pattern Catalog P1–P5 and visible beginning of P6; supplied Evidence Index | Join versus signal, cause preservation, repeated-close hazards, leak instruments | Report-only; process escalation and worker shutdown excluded |
| IN-17 | Go adversarial verification | Study report | `studies/aren-go-runtime-study/reports/final/01.03-adversarial-concurrency-and-failure-verification.md` | Executive Summary; Core Thesis; Approach Models; Pattern Catalog 1–8 and visible portion of 9; Notable Absences; Per-Source Notes; Open Questions; Evidence Index | Matrices, controlled faults, model traces, seeds, race placement, ownership accounting | Report-only; storage/process/cluster machinery excluded |
| IN-18 | Go ordered observation | Study report | `studies/aren-go-runtime-study/reports/final/01.04-ordered-observation-live-streaming-and-backpressure.md` | Executive Summary; Core Thesis; Approach Models; Pattern Catalog 1–4 and visible beginning of 5; supplied Evidence Index | Handoff, terminal delivery, consumer isolation, upstream-loss counterexamples | Report-only; high-volume transport mismatch and attribution caveat retained |
| IN-19 | Regression gating and CI | Study report | `studies/agent-harness-study/reports/final/18.03-regression-gating-ci-integration.md` | Executive Summary; Core Thesis; Approach Models; Pattern Catalog 1–8; Key Differences; Tradeoffs; supplied Decision Guide portion, Open Questions, Evidence Index | Effective gates versus nominal test presence; fail-closed execution | Only two material sources; four empty trees establish no Go CI behavior |

The template requires reading `projects/aren-phase-01-execution-lifecycle/reasoning/README.md`. Its content was not supplied, and no read-only workspace tool is available in this turn. It is not an inspected input. Compliance with any additional standard there remains unverified and must be checked before a passing project-reasoning review. No current Aren source, standalone prototype report, test execution, or CI configuration was inspected. [IN-12; IN-05, Project Reasoning Policy]

## Governing questions

| Question | Verification decision | Basis |
| --- | --- | --- |
| Which claims are pure? | Legal edges, outcome validity, cancellation dispositions, interpretation of fixed completion facts, and cursor boundaries can be checked with tables and a pure model. | IN-01 §§7, 11–17; IN-09–IN-11 |
| Which require concurrency? | Coherent publication, waiter release, acceptance before signal, matching revalidation, parent handoff, tail waiting, observer independence, and support release. | IN-08, Downstream obligations; IN-10–IN-11, Verification obligations |
| What does `-race` establish? | It detects data races on executed paths. It does not establish legal histories, correct arbitration, absence of deadlock, complete test selection, or quiescence. | IN-03 §11.6; IN-17, Pattern Catalog 1 and Notable Absences |
| How is order controlled? | Work barriers and operation acknowledgements first; narrowly scoped internal checkpoints only where a production-path window cannot otherwise be opened. | IN-10–IN-11, Verification obligations |
| How is release established? | Observe actual exits and registration release, then supplement with stack and retention checks in a still-running test process. | IN-03 §11.3; IN-10, Evidence execution requirements |
| How are instruments challenged? | Each high-value oracle rejects a named broken fixture or seeded production-path defect for the intended reason. | IN-03 §§7, 11.6 |
| What counts as completion evidence? | A requirement-linked record of executed tests, negative controls, race/stress/resource results, real CLI scenarios, and review of the realized contract. | IN-01 §§21–24; IN-02, Phase Exit Gate |

## Normative constraints

Source IDs resolve to exact paths in Inputs Used.

| Required proof | Source | Why ordinary unit tests are insufficient | Expected evidence class |
| --- | --- | --- | --- |
| All four legal edges and rejection of every other state pair | IN-01 §7; IN-08, Downstream obligations | Happy paths never exercise forbidden transitions or unknown vocabulary | Independent matrix, package fault tests, review |
| One immutable terminal commitment | IN-01 §§13–14 | A final-state assertion can miss replacement, duplicate events, or repeated close | Model, production-path duplicate attempts, controlled races |
| Complete state/event/outcome/timing publication | IN-01 §§14, 17–18 | Separately synchronized fields may still expose logically torn facts | Coherent-view oracle, controlled publication windows, race |
| Correct success/error/panic interpretation | IN-01 §§10–13; IN-09, Verification obligations | A mock resolver can hide capture defects and production divergence | Pure table and real invocation tests |
| Truthful cancellation and stable first cause | IN-01 §§12–13; IN-10, Normative constraints | Final cancellation alone says nothing about signal order or stopped work | Controlled work, cause table, explicit/parent races |
| Shared waiting and independent observation | IN-01 §§16–17; IN-11, Verification obligations | One receiver cannot expose exclusive consumption or lost wakeups | Multi-reader contract, handoff schedules, race |
| Defensive lifecycle publication | IN-01 §§15, 17; IN-03 §11.2 | Equal values before mutation do not prove alias isolation | Mutation attacks and independent rereads |
| Aren-owned release | IN-03 §11.3; IN-10, Evidence execution requirements | Successful process exit reclaims leaked goroutines and masks retention | Exit witnesses, registration accounting, leak and retention checks |
| Test instruments can fail | IN-03 §§7, 11.6 | Empty helpers and incorrect counters can produce green suites | Broken fixtures and seeded historical defects |
| Complete race-enabled CI | IN-03 §11.6; IN-01 §22 | A command in documentation does not prove execution or gate enforcement | CI configuration review, executed logs, gate-failure check |
| Real diagnostic execution | IN-01 §20 | Mocked display tests do not prove runtime integration | CLI contract tests and built-artifact smoke |
| Proportionate resource measurements | IN-03 §§11.7–11.8 | Event count and goroutine count do not explain retained bytes or contention | Benchmarks, churn measurements, profiles where indicated |
| Correct sprint sequencing and contract promotion | IN-02, carry-forward rule and Phase Exit Gate | Combined final success can conceal an incomplete Sprint 1 or undocumented supersession | Sprint reviews, traceability, phase review |

## Evidence

### Verification pattern assessment

The repository paths below are anchors quoted by inspected reports, not separately inspected source files or reproduced test results.

| Evidence | Technique | Failure class exposed | Determinism | Cost | Blind spot | Aren use |
| --- | --- | --- | --- | --- | --- | --- |
| IN-13, Explicit transition guard and Formality of the FSM; IN-17, Pattern 5; report anchor `studies/aren-go-runtime-study/sources/runc/libcontainer/state_linux_test.go:26` | Exhaustive legal/illegal matrix | Missing guards and permissive transitions | High for enumerated inputs | Low | Says nothing about publication or concurrent use | Adopt independently authored expected edges |
| IN-14, Core Thesis and P1; report anchor `studies/agent-harness-study/sources/pydantic-ai/pydantic_evals/pydantic_evals/evaluators/agentic.py:250-366` | Evaluation of an explicit trajectory record | Correct-looking final result reached through an invalid path | High for fixed records | Low at three or four lifecycle events | A captured record can omit facts | Adapt exact lifecycle checks; reject partial credit and LLM judges |
| IN-15, P1; report anchors `studies/aren-go-runtime-study/sources/temporal-sdk-go/internal/internal_task_handlers.go:653-657` and `studies/aren-go-runtime-study/sources/temporal-sdk-go/internal/internal_workflow_testsuite.go:1143-1147` | Production/test-double comparison | Repeated completion behaves differently in tests | High for controlled repeated calls | Low | A fake-only test still cannot prove the runtime | Require real commit-path fault tests |
| IN-15, P6; IN-16, P4; report anchor `studies/aren-go-runtime-study/sources/conc/pool/context_pool.go:39-47` | Cause-before-notification ordering | Consequence becomes visible before its cause | High with a held window | Low | Pool error ordering is not Aren acceptance | Adapt signal-time history assertions |
| IN-16, P1–P2; report anchors `studies/aren-go-runtime-study/sources/conc/waitgroup.go:28-43` and `studies/aren-go-runtime-study/sources/temporal-sdk-go/internal/internal_worker_base.go:467-478` | Explicit ownership and join order | Unjoined support and premature closure | High for exit witnesses | Low to moderate | Uninventoried carriers escape accounting | Adopt actual-exit witnesses plus independent leak checks |
| IN-17, Patterns 1–2; report anchor `studies/aren-go-runtime-study/sources/conc/pool/context_pool_test.go:176` | Race detector and bounded observation | Data races and progress failures; sleep-biased tests are a counterexample | Execution-dependent | Moderate under `-race` | Logical races can be race-free; polling does not force order | Adopt detector; use barriers for order and polling only for eventual release |
| IN-17, Pattern 7; report anchor `studies/aren-go-runtime-study/sources/nomad/scheduler/reconciler/reconcile_cluster_prop_test.go:131` | Generated pure-model comparisons | Combinations omitted by examples | Reproducible input trace | Moderate | Shared bugs between model and implementation | Adapt to a small lifecycle model; no scheduler framework |
| IN-17, Pattern 9 and Pebble Approach Model; report anchor `studies/aren-go-runtime-study/sources/pebble/metamorphic/meta.go:258` | Seed, trace, replay, reduction | Irreproducible random failures | Input reproducible; goroutine schedule is not | Moderate | Seed alone cannot replay Go scheduling | Adopt diagnostic records and small-trace reduction |
| IN-18, Pattern 4; report anchors `studies/aren-go-runtime-study/sources/docker-agent/pkg/server/eventlog.go:146-160` and `studies/aren-go-runtime-study/sources/docker-agent/pkg/server/eventlog_test.go:152-180` | Replay/live handoff test | Lost or duplicated event at registration | Controlled boundary plus stress | Low | Local gaplessness does not prove upstream recording | Adapt to canonical history and tail-wait establishment |
| IN-18, Docker Agent and Ollama Approach Models | Slow/absent-consumer tests | Producer coupling and upstream loss | High when consumer inactivity is held | Low | A short test may never fill an incidental buffer | Test the actual dependency, not an arbitrary sleep or oversized buffer |
| IN-19, Executive Summary and Patterns 3, 6, 8; IN-17, Rating Summary | Gate execution and failure-path review | Skipped, advisory, stale, or non-race suites | High for configuration and deliberate gate failure | Low | External branch protection is not proved by repository YAML | Adapt fail-closed checks; reject marker-file and aggregator machinery without need |
| IN-03 §§7, 11 | Historical Aren defect seeding | Mutable outcomes, missing timing, parent bypass, recursive fallback | High when the defect window is held | Moderate | Historical checks themselves were imperfect | Adopt named regressions, not historical pass totals |

Three restrictions govern transfer. First, IN-19 has two material sources, not six examples of implemented CI. Second, the Go reports did not have the authoritative Aren PRD; their first-wins, shutdown, and streaming recommendations cannot override it. Third, the reported absence of a complete duplicate-terminal harness is an evidence gap, not permission to omit Aren's direct assertion. [IN-07, assessment and Evidence; IN-15–IN-18, availability and absence disclosures]

No persistence crash harness, process-kill ladder, durable replay verifier, broker backpressure system, agent evaluator, or shadow Rust runtime is selected. Those instruments target excluded substrates. [IN-01 §§4, 19, 25; IN-03 §1; IN-17, Patterns 6–8]

### Historical defect register

Names in this table are proposed regression identifiers, not existing tests. “Why earlier tests missed it” distinguishes documented weakness from inference; the supplied summaries do not establish the exact history of every escaped defect. [IN-03 §§6–7]

| Defect | Original evidence | Broken invariant | Why earlier tests missed it | Permanent regression design | Negative control |
| --- | --- | --- | --- | --- | --- |
| Parent cancellation bypassed acceptance | IN-03 §7, parent-bypass row | INV-09, INV-10 | Final outcomes alone would miss missing acceptance; exact original gap not established | `ParentSignalRequiresAcceptedHistory`: inspect canonical facts from work after signal | Directly propagate parent cancellation before acceptance |
| Terminal cancellation preceded request event | IN-03 §7, causal-inversion row | INV-04, INV-10 | Plausible terminal state does not establish causal history | `RequestPrecedesSignalAndTerminal`: hold propagation and completion separately | Signal first or append request after terminal |
| Waiters shared mutable outcome | IN-03 §7, shared-outcome row | INV-07, INV-08 | Read-only equality checks do not attack aliases; exact old assertions unknown | `WaiterSnapshotsDoNotAlias`: mutate lifecycle snapshots, reread independently | Return internal mutable lifecycle storage |
| Recursive invariant handling and double close | IN-03 §7, invariant-failure row | INV-02, INV-03, INV-15 | Error protection was itself an untested failure path | `InvariantFailureIsNonRecursive`: inject before start, during work, after work, and after terminal | Re-enter failed transition or close readiness twice |
| Missing terminal timing | IN-03 §7, terminal-timing row | INV-03, INV-06 | Terminal-state checks can omit timing | `TerminalPublicationIncludesTiming`: every branch and every waiter | Remove finish time or publish it later |
| Unlock window exposed partial publication | IN-03 §7, cleanup-window row | INV-03 | Original defect involved later cleanup; exact old schedule unknown | `PublicationHasNoObservableIntermediateRecord`: adapt only the split-publication failure class | Expose state before outcome/history under otherwise race-free locks |
| Classified failures flattened | IN-03 §7, flattened-failure row | INV-05 | Message equality can pass after cause/type loss | `FailureRetainsOriginKindAndCause`: assert classification, `errors.Is`, and `errors.As` separately | Replace wrapped cause with a new string error |
| Completed run retained by watcher | IN-03 §7, retention row | INV-13 | Process exit or parent cancellation can mask a lifetime leak | `CompletedRunReleasesLiveParentRegistration`: retain parent through repeated completion batches | Leave stop registration attached until parent ends |
| Reentrant lock/watch liveness deadlock | IN-03 §7, deadlock row | INV-12, INV-13 | Memory safety did not establish forward progress | `SupportJoinDoesNotHoldLifecycleOwnership`: hold watcher before its lifecycle access | Join watcher while holding the required lock |
| Empty helpers and unrelated counters passed | IN-03 §7, weak-test row | INV-17 | Documented instrument defects | `VerificationOraclesRejectBrokenCases`: assert the expected violation identity | Feed duplicate terminal, missing event, retained carrier, and changed payload |
| Production and test completion behavior diverged | IN-15 P1, cited production/test paths | INV-02, INV-17 | Test environment implemented a stronger guard than production | Exercise duplicate completion against the real private owner | Guard only the fake completion slot |

Later-phase schema, JSON, provider-stream, subprocess, admission, and policy defects remain excluded. Only their directly applicable mutation, ordering, ownership, or instrument lessons transfer. [IN-01 §4; IN-03 §7; IN-07, negative-transfer assessment]

## Reference model

The model is a test-only, pure description of observable lifecycle facts. It uses its own expected-edge table and outcome rules rather than production guards, constructors, sequence helpers, or resolver code. Shared public vocabulary may be translated at the comparison boundary; shared decision logic is prohibited. [IN-03 §11.6; IN-13, Explicit transition guard; IN-17, Patterns 5 and 7]

Minimum model data consists of an abstract run identity, lifecycle state, logical history, optional first accepted cancellation, propagation status, optional completed work fact, and optional terminal outcome. Time is represented by constraints and stable tokens rather than predicted wall-clock values. Observer cursors are test-side positions over the model history, not runtime subscription machinery. [IN-01 §§8, 12–18; IN-09–IN-11]

`Complete` below means supplying captured work facts to terminal resolution, not allowing work or a caller to choose a terminal state. Recording a completed work fact while propagation is pending is internal model bookkeeping, not a public event or new lifecycle state. [IN-09, Governing questions; IN-10, Verification obligations]

| Command or fact | Model precondition | Model update | Expected output | Illegal case |
| --- | --- | --- | --- | --- |
| Create | No run exists in this model instance | Bind nonempty abstract identity; state `created`; append creation at `0` | Created prefix | Re-create the same run instance or record before identity |
| Start | State `created` | State `running`; capture start token; append start at `1` | Running prefix, no terminal outcome/finish | Any other source state |
| Complete: normal return | Running; invocation including its defers finished | Record result/error fact; resolve when accepted propagation, if any, has settled | Nil error succeeds; matching accepted cause cancels; otherwise work returned-error failure | Completion before invocation finished; malformed completion mode |
| Complete: panic unwind | Running; supervised unwind finished | Record panic fact; resolve failure when eligible | `failed`, `work/panic`, stable diagnostic projection, no result | Treating zero result/error storage as a normal return |
| Request cancellation | Running and no accepted request | Retain first reason/cause; append request; mark propagation owed | `accepted`; state remains running | Request in internal pre-start state is not a supported scheduling command |
| Request cancellation again | Running with acceptance already committed | No change | `already_requested` | Replacing reason/cause or appending another occurrence |
| Request after terminal | Terminal, whether or not cancellation was previously accepted | No change | `already_terminal` | Returning `already_requested` instead, mutating history, or restarting support |
| Propagation settled | Acceptance exists and propagation is pending | Mark accepted signal established; allow pending completion to resolve | Work signal must have corresponding committed acceptance | Signal without acceptance or signal a different effective cause |
| Duplicate internal completion | Terminal already committed | No lifecycle change | Recognizable rejected mutation; existing outcome retained | Replacement, extra terminal event, or repeated close |
| Coherent inspection | Any created run | No change | One valid committed prefix and matching current facts | Joined view assembled from incompatible prefixes |
| Read from cursor `c` | `0 <= c <= n`, where `n` is current history length | Advance only on delivered events | Inclusive suffix; active tail waits; terminal tail exhausts | `c > n` errors without mutation; no future-cursor clamping |
| Abort a pending local read | No ready suffix or terminal exhaustion at its decision point | No cursor or run change | Distinct local abort | Run cancellation, event loss, or advancement without delivery |

These rules carry forward strict directional `errors.Is(returnedError, acceptedEffectiveCause)`, including the selected treatment of wrapped and joined causes. A coarse `ctx.Err()` sentinel does not match a distinct custom cause merely because both indicate cancellation. The model uses explicit match facts for pure arbitration tests; real Go errors separately exercise cause inspection and capture through the runtime. [IN-09, Normative constraints, Verification obligations, Self-critique; IN-10, Normative constraints]

The matrix includes normal, nil, zero, and reference-bearing results; non-nil error interfaces containing typed nils; same-text distinct errors; direct, wrapped, joined, custom, and unmatched causes; and panic despite accepted cancellation. Work-returned Aren-shaped errors remain work-origin returned errors. [IN-09, Candidate outcome models and Verification obligations]

The model deliberately omits synchronization code, goroutine scheduling, context implementation, callback execution, heap reachability, arbitrary result-object internals, and work effects. It therefore cannot prove those properties. An abnormal exit such as `runtime.Goexit` must not become success; Sprint 1 must pin its concrete exposure and containment behavior before acceptance rather than inventing that behavior inside this model. [IN-09, Verification obligations and Self-critique]

For controlled concurrent histories, the harness establishes the relevant order externally. For genuinely colliding operations, compare the observed results against legal serializations consistent with invocation/response order and recorded canonical acceptance. Do not derive the entire expected result from the runtime's final state: independent work facts and request responses must constrain the comparison. No general linearizability framework is required for this small operation set. [IN-01 §13; IN-17, Pattern 7; project verification decision]

## Invariant catalog

“Both” means Sprint 1 proves the implemented core and Sprint 2 preserves it while adding cancellation and live observation. Required suites are proof classes, not existing package names. [IN-02, both implementation waves]

| Invariant ID | Statement | Scope | Observable oracle | Concurrent schedule needed? | Required suites |
| --- | --- | --- | --- | --- | --- |
| INV-01 | One opaque identity exists before creation, remains constant, and appears in every event/outcome; separate invocations do not share run-local state | Both | Work-entry count, identity equality, concurrent distinct-run fixtures | Yes for isolation | Contract, race, review |
| INV-02 | A completed run has exactly one terminal commitment and event; later mutation cannot change either | Both | Full history, stable outcome/timing before and after duplicate attempts | Yes | Model, fault, race |
| INV-03 | Every supported coherent view is a valid committed prefix; returned waiters can immediately retrieve complete terminal facts | Both | Independent prefix checker over state, history, outcome, cancellation, timing | Yes | Contract, publication schedule, race |
| INV-04 | Creation sequence is `0`; sequences are contiguous; terminal is last; cancellation request occurs at most once while running | Both | Exact three-event or four-event trajectory, not counts alone | Yes | Model, contract, race |
| INV-05 | Outcome follows captured return/panic facts and accepted cause; non-success exposes no successful result; origin/kind/cause remain distinguishable | Both | Independent truth table, exact result token, classification and cause assertions | Yes for acceptance/revalidation | Pure, invocation, fault, race |
| INV-06 | Start and finish are captured once; finish is absent while active; completed timing is coherent and identical across readers | Both | Before/after observations and cross-view equality | Yes | Contract, publication, race |
| INV-07 | Waiting observes independently completed publication; all early/late waiters receive the same logical outcome without consuming or finalizing it | Both | Multiple waiters plus completion without any waiter | Yes | Contract, race, release |
| INV-08 | Published Aren-owned containers, event data, and diagnostic collections cannot be mutated through supported access | Both | Mutation then independent fresh read and another observer's retained view | Yes | Mutation, race, API review |
| INV-09 | Exactly one active request is accepted; first reason and effective cause remain stable; repeat/terminal dispositions are correct | Sprint 2 | Request responses, history, signal cause, outcome metadata | Yes | Model, cancellation, race |
| INV-10 | Acceptance is committed before work-visible signal; accepted propagation settles before terminal publication | Sprint 2 | Signal-time coherent read; held propagation with completed work | Yes | Controlled schedule, race |
| INV-11 | No terminality, finish time, waiter outcome, or normal observer exhaustion occurs while supervised work remains active | Sprint 2; core completion in Sprint 1 | Held work and defer barriers with acknowledged inspection opportunities | Yes | Contract, schedule, race |
| INV-12 | Readers are independent; tail handoff is gapless; ready suffix drains before normal completion; local abort has no control effect | Sprint 2 | Exact per-reader suffixes, cursor state, request history, progress witnesses | Yes | Cursor matrix, handoff, race |
| INV-13 | After controlled work completion, all Aren-owned support reaches its declared release boundary without consumer drainage or parent termination | Both | Actual exit/stop witnesses, registrations, stack and retention evidence | Yes | Ownership, leak, churn |
| INV-14 | History survives while the run remains reachable; replay preserves identity, payload, and time without adding canonical events | Both for retention; Sprint 2 for replay | Before/after history equality and independently repeated suffix | Yes | Contract, replay, mutation |
| INV-15 | Illegal transitions and unknown internal vocabulary are rejected visibly; containment never fabricates an edge, stops active work falsely, or recursively finalizes | Both | Independent negative matrix and state-specific fault injection | Yes for fault interactions | Matrix, fault, race, review |
| INV-16 | Observation/control/internal mutation capabilities remain separate through the actual exposed values | Both; controller in Sprint 2 | External-consumer surface tests and dynamic capability checks | Not primarily | API tests, architecture review |
| INV-17 | Every high-value instrument rejects its designated broken case, and required CI lanes execute rather than silently skip | Both | Named violation, mutant result, test-selection and gate evidence | For schedule/leak instruments | Negative controls, CI review |
| INV-18 | Resource claims distinguish bounded lifecycle count from user-sized bytes and caller-created concurrency; CLI claims match runtime guarantees | Both; CLI in Sprint 2 | Measurements, ownership ledger, real CLI output and exit behavior | Yes for churn | Measurement, smoke, review |

The state matrix contains 25 ordered pairs over the five states: four legal edges and 21 forbidden edges, plus invalid vocabulary cases. Cancellation-request recording is not an extra legal `running -> running` state transition. In Sprint 1, the cancellation edge can be validated as vocabulary/guard policy without claiming implemented cancellation integration. [IN-01 §§7, 15; IN-02, Sprint 1 Build and Deferred]

Exact equality applies to Aren-owned logical facts, not a universal deep comparison of arbitrary result/error graphs. Tests use stable fixtures and the dependency's explicit borrowed-cause/result boundary. Concurrently mutating a caller-owned error or successful result is not an Aren lifecycle-immutability test. [IN-01 §17; IN-09, Risks and Self-critique; IN-11, Verification obligations]

## Schedule-control design

Work functions should expose `entered`, `signalObserved`, and test-owned release barriers as needed. Tests must release intentionally blocked work in cleanup even after an assertion fails. Internal checkpoints are earned only for windows such as acceptance-to-signal, precommit revalidation, tail-wait establishment, and registration teardown. They must not become public configuration, replace the production owner, or add ordering that conceals the target defect. [IN-10–IN-11, Verification obligations; IN-03 §§11.3, 11.6]

A checkpoint must say what has actually happened. “Waiter goroutine started” is not proof that it entered the blocking wait. A nonblocking receive that finds no result is not proof that an early-return defect had a chance to execute. Use a test-visible wait decision, an appropriate toolchain-supported scheduling instrument, or a held production boundary with an acknowledged observation opportunity. Timeouts are failure bounds, not ordering mechanisms. [IN-03 §7; IN-17, Pattern 2; project verification decision]

The repetition policy below is a proposed starting budget: **D** is one forced execution of each named ordering on every required run; **R** is a proposed 20 repetitions of those schedules under `-race`; **S** is a proposed 10,000 mixed collisions for Sprint 2 closure. These numbers are not measurements or coverage probabilities. Sprint reasoning must finalize them after measuring suite cost and record any adjustment without deleting a required ordering. [IN-02, Evidence command policy; IN-10, Evidence execution requirements]

| Race or window | Controlled checkpoints | Forbidden use of sleep | Repetitions | Assertions during race | Assertions after quiescence |
| --- | --- | --- | --- | --- | --- |
| Work entry/return versus inspection | Hold work and its final defer; read before release | Do not delay work “long enough” for readers | D; R | Running, no finish/outcome; creation/start already recorded | INV-01–INV-07 |
| Terminal publication versus readers/waiters | Hold attempted publication at relevant internal boundaries; start fresh coherent reads | Do not sample only after a guessed publication duration | D per boundary; R | Reads either await ownership or return complete valid facts, never a torn record | All returned waiters retrieve matching terminal history |
| Conflicting duplicate finalization | Invoke real private commit boundary with controlled conflicting candidates | No probabilistic reliance on simultaneous launch | D in both orders; R | One valid candidate commits; later attempt rejected | Outcome/history/timing unchanged; no double close |
| Explicit cancellation versus completion | Force acceptance before completion; separately complete before request; then collide | No sleep-biased winner | D for each completion branch; R; S | Correct disposition and no false stop | Independent truth table, exact history, release |
| Explicit versus parent cancellation | Force each source first with distinct causes; then collide | No assumption watcher scheduling follows parent cancel call immediately | D each order; R; S | One accepted fact, one signal cause | First reason retained across history/context/outcome |
| Acceptance versus signal/teardown | Pause after accepted event but before propagation; release work to return nil | No delay used to imply signal delivery | D; R | Nonterminal until accepted propagation settles | Success, settled propagation, no helper cycle |
| Startup parent handoff | Cancel before start; between initial check and integration; during installation; complete during installation | No reliance on an asynchronous watcher eventually running before work | D each window; R | Already-cancelled parent signaled before invocation; other races honor serialized acceptance | No missed signal or late retained stop handle |
| Preliminary cause inspection versus new acceptance | Hold classification after “not accepted”; accept; resume | No statistical attempt to hit the stale-candidate gap | D; R | Readers/control still progress | Final interpretation includes newly committed acceptance |
| User error inspection versus lifecycle access | Custom `Is` or formatter performs safe reentrant observation, or waits on a controlled barrier | No timeout treated as successful callback isolation | D | Other lifecycle operations can progress; no user method under ownership | Correct stable diagnostics and release |
| Delayed/ignored cooperation | Hold work after signal or make it ignore context until release | No “cancel then sleep and assume stopped” | D for nil/error/panic exits; R | Running, no terminal event or normal EOF, repeat requests inert | Chosen return fact resolves correctly; all support releases |
| Tail check versus append | Append immediately before and after tail-wait establishment | No random registration delay | D each ordering for start/request/terminal availability as applicable; R | No active-tail EOF or lost next event | Exact suffix, no gap or unsolicited repeat |
| Backlog versus terminal readiness | Pause reader with unread events; complete work | No drain assumption based on elapsed time | D for each valid cursor; R | Ready backlog wins over terminal exhaustion | Terminal consumed when in suffix, then normal exhaustion |
| Local abort versus availability | Force availability before read decision; separately abort while no data is ready | No raw random `select` used as expected policy | D each ordering; R | Ready suffix/exhaustion wins when visible; otherwise local abort without cursor movement | Other reader/work unaffected; replay complete |
| Idle/abandoned reader versus progress | Hold one reader idle; abort another pending read; advance a fast reader and complete work | No assumption a brief idle period exercises capacity | D; R; S churn | Work, cancellation, fast reader, and waiters progress without drainage | No retained observer producer/registration |
| Watcher lifecycle access versus join | Hold watcher before acquisition; finish work and trigger support release | No sleep loop presented as a join | D; R | No lock-held join cycle; terminal view remains coherent | Watcher exit witnessed while parent remains alive |
| Completion versus late registration release | Finish before installation stores its stop handle; resume installation | No hope that installation normally wins | D; R | No post-terminal mutation | Late-installed registration is released |
| Mutation versus concurrent readers | Mutate only returned lifecycle copies while other readers inspect | No unsynchronized caller-owned result mutation | D per mutable field; R | Other observations unchanged and race-free | Canonical and retained snapshots unchanged |
| Cross-run interference | Hold one run active with idle observers; complete/cancel another | No global sleep-based progress assumption | D; R; S | One run cannot control another's lifecycle | Independent histories, identities, causes, release |

A correctly blocked coherent read need not expose a snapshot while the transition owner is paused. That is not a failure. Conversely, a test cannot declare atomicity proved merely because its hook forced every reader to wait. The publication instrument must be challenged with a deliberately split, race-free publication implementation that exposes the intermediate state and is rejected. [IN-08, Risks and Self-critique; IN-03 §11.6]

## Test portfolio

All commands below are proposed execution forms for the Aren implementation directory, `/home/antonioborgerees/coding/ultraplan/Aren`. No command was run here. Selector and benchmark names are placeholders until sprint reasoning supplies actual names; plans must verify selected tests exist and execute. [IN-02, both sprint Evidence sections; IN-17, stale-target warnings]

| Test level | Behaviors proved | Real dependencies | Faults injected | Command | Release gate |
| --- | --- | --- | --- | --- | --- |
| Pure/model | Legal matrix, validity, dispositions, fixed-fact arbitration, cursor table | Test-only independent model | Invalid commands, malformed records, generated traces | `go test -count=1 ./...` including model tests | Both sprints for implemented behavior |
| Package contract | Real invocation/capture, history, timing, causes, waiting, defensive access | Actual runtime package | Returned errors, panic/defer panic, typed nil, abnormal exit, invalid candidate | `go test -count=1 ./...` | Both |
| Static/API | Private authority, vocabulary handling, dependency direction, usable external surface | Actual package exports and build | Invalid consumer examples where useful; unknown internal values | `go vet ./...`; `go build ./...`; selected analyzer if earned | Both |
| Race | All required contract, model-driver, fault, and schedule cases | Actual goroutines and synchronization | Controlled collisions and mutation of returned copies | `go test -race -count=1 ./...` | Mandatory CI in both |
| Repeated schedules | Repeated publication, completion, cancellation, observer and teardown boundaries | Real runtime with narrow test checkpoints | Forced orderings | `go test -race -count=20 -run '<actual targeted selector>' ./...` | Proposed required targeted lane; sprint plan fixes budget |
| Collision stress | Mixed work returns, explicit/parent requests, waiters, reads, observer churn | Real runtime without schedule-forcing hooks | Seeded workloads and concurrent launch | Actual stress selector and configured count, with `-race` | Sprint 2 closure; bounded presubmit subset |
| Leak/quiescence | Carrier exits, registration release, abandonment, live-parent retention | Real contexts, support carriers, still-running process | Held/abandoned readers, late installation, missing-release controls | Actual ownership/leak selector with `-race` | Both for existing resources |
| Measurement | Allocations, active/retained footprint, contention, churn recovery, relevant latency spread | Real runtime; non-race measurement artifact | Representative success/failure/panic and Sprint 2 cancellation/observer workloads | `go test -run '^$' -bench '<actual benchmark selector>' -benchmem -count=5 ./...` plus named churn/profile commands | Measurement evidence at sprint/phase review |
| Diagnostic CLI | Actual runtime integration, ordered presentation, request/terminal distinction, status mapping | Built Aren CLI using real runtime | Required success/fail/cancel/race scenarios | `aren dev run success`, `fail`, `cancel`, `race` | Sprint 2 and phase exit |
| Review | Whole-program authority and memory-order argument; gate execution; scope and contract accuracy | Real source, tests, CI, recorded evidence | Compare guarantees against actual bypass/release paths | Required architecture/sprint/phase review process | Governed acceptance |

The model-based requirement can be satisfied without adding a property-testing dependency: fixed exhaustive cases plus a small generated command-trace driver compared with the independent model are sufficient in shape. If native Go fuzzing is added, the plan must name the actual package and target, run seed cases in normal CI, schedule exploration, and preserve failing inputs. Merely checking in a fuzz function is not evidence of exploration. [IN-03 §11.6; IN-17, Pattern 7 and Notable Absences]

### Required scenario coverage

The implementation evidence ledger must map every PRD §21 item, not only these families, to actual tests and commands. One test may satisfy several rows when its assertions genuinely cover them. [IN-01 §21; IN-03 §12.9]

| Scenario family | Required cases | Invariants | Sprint |
| --- | --- | --- | --- |
| Ordinary completion | Normal/nil/zero result; result plus error; wrapped failure; panic and work-defer panic | INV-01–INV-08 | 1 |
| Negative lifecycle | All forbidden pairs, unknown vocabulary, malformed terminal records, duplicate finalization, abnormal non-return | INV-02, INV-03, INV-05, INV-15 | 1 |
| Concurrent core | Early/late/multiple waits; no waiters; concurrent inspection; mutation attacks; timing coherence; cross-run isolation | INV-01–INV-08, INV-13–INV-16 | 1 |
| Cancellation sources | Explicit, repeated, parent, already-cancelled parent, supplied deadline expiry, explicit/parent collisions | INV-09–INV-11 | 2 |
| Post-acceptance outcomes | Nil-error success; direct/wrapped/joined matching cause; custom-cause mismatch; unrelated error; panic | INV-05, INV-09–INV-11 | 2 |
| Cooperation | Immediate acknowledgement, delayed acknowledgement, ignored signal with later controlled release | INV-10, INV-11, INV-13 | 2 |
| Observation | Registration before/during/after completion; all valid cursors; tail/invalid cursor; gapless handoff; terminal drain; independent cursors; replay | INV-04, INV-12, INV-14 | 2 |
| Abandonment | Absent/slow/idle observers; pending local abort; abort/availability races; optional locally cancellable outcome wait | INV-07, INV-12, INV-13 | 2 |
| Broad hardening | High-iteration mixed collisions, watcher races, retained-parent churn, leak checks, full suite under `-race` | All applicable invariants | 2, preserving 1 |
| Runnable evidence | Required real-runtime CLI scenarios and useful statuses; recommended panic/parent/ignore scenarios when implemented | INV-18 | 2 |

Deadline tests must distinguish actual supplied-parent expiry from explicit cancellation carrying a deadline-shaped error. Use an already-expired real parent for one case and a controlled real expiry or suitable test-time facility for the live path. Waiting for `parent.Done()` establishes expiry; sleeping for the nominal duration does not establish acceptance or propagation. [IN-01 §12; IN-09–IN-10, Verification obligations]

### Quiescence and retention

Terminal publication and support quiescence are separate assertions. Ordinary outcome waiting is not assumed to join every support goroutine; the dependency permits a finite internal epilogue that does not depend on observers or parent termination. Tests must reach that epilogue's actual completion rather than treating readiness closure as proof. [IN-10, Governing questions and Self-critique]

Each realized spawn, context callback, timer, registration, and support reference needs a row in a small ownership ledger: creator, owner, stop trigger, release path, and observable witness. A decrement in a counter before a goroutine actually exits is not an exit witness. Test-owned waiter goroutines must also be joined so the harness does not confuse its own leaks with Aren's. [IN-03 §11.3; IN-16, P1; IN-10, Evidence execution requirements]

Use three complementary fixtures:

1. Retain completed run handles while the parent remains active. History/outcome retention is expected; watchers and delivery producers must still release.
2. Drop unneeded run handles and temporary snapshots while retaining the parent. Registration accounting and targeted heap evidence must not show completed runs pinned by parent integration.
3. Abort/drop readers during an active run, then release work and verify support exits. Canonical history remains available through any deliberately retained run handle.

Supplement explicit witnesses with a package/process leak sweep and stack diagnostics. Heap/retention measurements are needed for callbacks or references that leak without a goroutine. Finalizers, one forced GC, raw goroutine-count equality, or RSS returning to its original value are not decisive alone. A process must remain alive through the checks; CLI exit cannot certify library leak freedom. [IN-03 §§11.3, 11.7–11.8; IN-10–IN-11, Verification obligations]

### Measurement scope

A completed normal lifecycle has three events, or four with accepted cancellation. That is a bound on event count, not on arbitrary result, error, panic-diagnostic, or caller-retained snapshot bytes. No global run-admission system or result-size budget is introduced to make a benchmark convenient. [IN-01 §§4, 15–17; IN-09, Risks; IN-11, Self-critique]

Measure completed-run allocation, blocked active runs, repeated completion while a parent remains alive, retained versus dropped handles, concurrent inspection, and Sprint 2 reader churn. Record latency percentiles where the workload yields meaningful samples; report allocations and heap or RSS, goroutines, contention, workload scaling, and recovery after releases. Use non-race builds for performance numbers and separately retain race-enabled correctness results. No historical prototype percentage becomes a production threshold. [IN-03 §§4, 8, 11.8]

## Falsification plan

An oracle can be self-tested with broken records, but a broken-record test does not prove a production-path window was exercised. Every critical schedule or release instrument therefore needs both a positive real-runtime case and a designated negative control. Mutations must be isolated from the normal implementation and removed from the accepted artifact. [IN-03 §§7, 11.6; IN-15 P1]

| Instrument | Deliberately broken case | Required failure | What would invalidate the instrument |
| --- | --- | --- | --- |
| Transition matrix | Permit `created -> succeeded` or terminal restart | Exact forbidden edge identified | Expected edges obtained from production guard |
| Trajectory checker | Duplicate terminal; omit start/request; reuse sequence; mismatch terminal type | Specific identity/order/state violation | Only counts events or accepts an arbitrary terminal suffix |
| Outcome checker | Cancellation overrides nil-error success; failed branch retains result | Exact branch mismatch | Expected outcome calculated by production resolver |
| Cause/diagnostic checker | Flatten error; drop stack; match by text; accept coarse sentinel after distinct custom cause | Missing identity/type/diagnostic or wrong state | Message equality is the only assertion |
| Publication schedule | Publish terminal state and release readiness before outcome/history using race-free split writes | Returned coherent view or waiter exposes a forbidden combination | Hook itself serializes all readers past the defect |
| Duplicate-finalization harness | Remove production terminal guard but leave fake guard | Real owner produces replacement/duplicate and test fails | Only fake completion is invoked |
| Cancellation-order harness | Signal before recording acceptance | Work's signal-time observation lacks required acceptance | Work is prevented from inspecting until after recording |
| Propagation/revalidation harness | Allow terminal teardown to overtake signal; commit stale “no acceptance” classification | Wrong eligibility or outcome detected | Test waits for propagation before opening completion |
| Cursor/handoff harness | Separate tail check and wait registration without reconciliation | Missing next event or diagnosed lost-wakeup timeout | No checkpoint proves append occurred inside the gap |
| Terminal-drain checker | Check terminal flag before unread backlog | Missing suffix/terminal event | Test starts only at terminal tail |
| Observer-independence harness | Require idle consumer acknowledgement or a mandatory unread send | Fast reader/work/waiter progress fails at named boundary | Test accidentally drains the supposedly abandoned reader |
| Immutability checker | Return shared event slice, nested diagnostic, or stack storage | Another view changes or race detector reports alias access | Test mutates only a detached top-level scalar |
| Exit/registration ledger | Omit watcher stop or close witness prematurely | Outstanding actual resource identified | Counter is unrelated to exit or decremented early |
| Independent leak check | Seed a blocked support goroutine or retained parent callback | Stack/registration/retention evidence detects it | Broad ignore hides it or process exits before checking |
| CI gate | Fail a known required test; use an empty targeted selector in a gate self-check | Required lane blocks acceptance; empty selection is rejected | Success marker, skip, retry, or advisory job turns it green |
| CLI integration | Stub fixed display or disconnect CLI from real runtime | Scenario/runtime integration assertion fails | Golden text alone is treated as runtime evidence |

For an intentionally broken case, a generic package timeout is usually insufficient: record which boundary was reached, what was expected to progress, and the relevant stacks or trace. A lost-wakeup or leak control may legitimately fail by timeout, but only after the harness proves it opened the intended window. [IN-17, Pattern 2 and Pebble Approach Model]

A mutation campaign need not become a permanent mutation platform. Keep a small, named set of high-value controls, record the targeted command and expected failure, and rerun them when their instrument or production boundary changes. Broken-fixture self-tests can remain ordinary fast tests. [IN-03 §11.6; IN-04, Evidence-Grounded Planning]

## Flake and suite economics

### Required and extended lanes

The required per-change lane contains deterministic contract/model cases, negative-oracle fixtures, owned-resource tests for the realized runtime, and the complete Go suite under `-race`. Targeted repeated schedules supplement that lane. Longer collision/churn runs are required for Sprint 2 closure and must also have a named ongoing owner and execution cadence if retained as regression protection. They must not disappear into an optional command nobody runs. [IN-01 §§21–24; IN-03 §11.6; IN-17, Notable Absences]

Long soaks and exploratory fuzzing may be separated for cost, but required semantic assertions cannot be moved into a non-gating lane merely because they expose flakes. Performance measurements may run without `-race`; correctness scenarios still require race-enabled execution. [IN-01 §22; IN-03 §§11.6–11.8]

### Reproduction and failure bounds

Every generated or stressed failure must record the source revision, diff status, Go version, OS/architecture, `GOMAXPROCS`, workload seed, iteration, command, test selection, work facts, request responses, canonical history, and reached checkpoints. Save the smallest failing command trace when practical. A seed reproduces generated inputs, not the Go scheduler; promote discovered races into controlled regressions. [IN-03 §11.8; IN-17, Pattern 9; IN-10, Evidence execution requirements]

Use bounded test watchdogs with diagnostic output and cleanup headroom. Sprint reasoning must select concrete budgets from observed normal and race-enabled runtimes. A timeout means the required progress or test completion was not established within the test budget. It does not prove work was forcefully stopped, create a runtime timeout feature, or justify marking a run cancelled. [IN-01 §§4, 12, 19; IN-10, Self-critique]

Automatic retry-to-green is not acceptance. A rerun can gather evidence, but the original failure remains a defect or unexplained flake until resolved. Stress must not require both winners to appear by chance; deterministic cases force each legal order separately. [IN-03 §11.6; IN-17, Pattern 2 and CI caveats]

### Toolchain and platform

Sprint 1 must record the supported Go version and practical CI platform matrix. `testing/synctest` or another toolchain-dependent facility may be used only after verifying its availability and semantics for the selected version. No clock abstraction is assumed. Channel barriers cover most ordering; actual parent-deadline tests cover real context integration. [IN-01 §18; IN-03 §12.10; IN-16, cited shutdown tests]

Linux-hosted evidence does not automatically establish other platforms. Where a target cannot run the race detector, record that limitation and obtain race-enabled evidence on a supported target without mislabeling a skipped detector run as a pass. Package portability and race coverage are separate claims. [IN-03 §§11.6, 12.10]

### Evidence and gate integrity

Each sprint review must link requirement, invariant, actual test/subtest, command, revision, execution result, negative control, and residual limitation. Targeted `-run` and benchmark commands must be checked for nonempty selection. Required tests must not be hidden behind `-short`, undocumented environment variables, unconditional skips, or `continue-on-error`. [IN-01 §§21–24; IN-17, stale/non-gating suite warnings; IN-19, Patterns 3 and 6]

For measurements and built CLI smoke, add artifact digest, hardware, workload size, repetitions, and spread. Review actual CI dependency/required-check behavior; repository configuration alone may not prove externally configured branch protection. No custom reporter, aggregator, artifact format, or PR-comment system is required when ordinary Go/CI output suffices. [IN-03 §11.8; IN-19, Open Questions; IN-04, scope-control doctrine]

## Project conclusions

| Conclusion | Evidence basis | Required in Sprint 1 | Added in Sprint 2 | Reopen trigger |
| --- | --- | --- | --- | --- |
| VG-01: Use an independent small model and exact trajectories, not final-state-only tests | IN-01 §§7–18; IN-14 P1; IN-17 Patterns 5, 7 | Independent edge/outcome model and core history assertions | Acceptance, cause, cursor, and collision traces | Model copies production decisions or cannot express a realized legal history |
| VG-02: Separate interpretation, one-commit exclusion, and coherent publication proofs | IN-08–IN-10; IN-15 P1–P2 | Real capture/resolution/commit tests and duplicate attempts | Acceptance/revalidation/propagation races | A new path can publish without the tested owner |
| VG-03: Force important orders; use stress only as supplemental exploration | IN-03 §11.6; IN-17 Patterns 2, 9 | Work, publication, waiter, fault boundaries | Every cancellation, observer, and watcher window listed above | A discovered race has no controlled reproducer or hooks mask behavior |
| VG-04: Defend Aren-owned publication while preserving the explicit result/cause ownership boundary | IN-01 §17; IN-09 and IN-11, Verification obligations | Outcome/history/diagnostic mutation attacks | Cancellation payload and live-reader mutation attacks | New mutable lifecycle field or changed borrowed-object contract |
| VG-05: Prove actual release independently from outcome readiness and process exit | IN-03 §11.3; IN-10, Self-critique | Ledger and exit evidence for existing carriers | Live-parent registration, observer abandonment, late-installation and churn checks | New spawn/callback/timer or unexplained retention |
| VG-06: Every high-value verification instrument must reject a named broken case | IN-03 §§7, 11.6 | Matrix, outcome, publication, mutation, invariant, release controls | Signal, stale-match, handoff, observer, watcher controls | Instrument or tested boundary changes; a relevant mutant survives |
| VG-07: Complete race-enabled CI and effective gate execution are mandatory | IN-01 §22; IN-03 §11.6; IN-19 restrictions | Full suite, actual selections, gate-failure evidence | Expanded full suite and required closure stress | Required tests become skipped, advisory, stale, or retry-to-green |
| VG-08: Measure resource behavior proportionately; do not invent a runtime budget or infrastructure | IN-03 §§11.7–11.8; IN-04 | Allocations, active/retained footprint, carrier release, relevant contention | Reader/cancellation churn and release recovery | Measured growth is unexplained or a real budget is introduced |
| VG-09: Carry forward realized Sprint 1 decisions and tests explicitly | IN-02, Cross-Sprint Carry-Forward Rule | Record APIs, ownership, tests, limits, and review findings | Preserve/extend/supersede/unaffected classification with evidence | New evidence proves an earlier decision insufficient |
| VG-10: Phase acceptance requires real runtime CLI evidence and a realized contract, not this document | IN-01 §§20–24; IN-05, review policy | Runnable tested core; no premature cancellation/observer claims | Required CLI, phase review, then contract promotion | Documentation disagrees with tested behavior or foundational ambiguity remains |

Sprint 2 must consume Sprint 1 requirements, index, handbook, every completed area output, authoritative `reasoning.md`, plan, and available execution/review evidence. Existing core regression tests remain requirements, not disposable scaffolding. A changed semantic oracle requires an explicit superseding decision, not merely updated golden output. [IN-02, Cross-Sprint Carry-Forward Rule]

No additional broad study is required to proceed with this verification design. The remaining work is the contained reasoning-standard check, sprint-level choice of actual mechanisms and commands, and fresh implementation evidence. The missing README check and unexecuted obligations prevent this area from claiming a passing project review. [IN-05; IN-07, Self-critique; IN-12]

## Trade-Offs

| Trade-off | Selected position | Benefit | Cost and limit |
| --- | --- | --- | --- |
| Deterministic seams versus production intrusion | Prefer work barriers and public behavior; add narrow internal checkpoints only for inaccessible windows | Reliable coverage of specific races | Hooks can conceal ordering defects; each needs a negative control and real-path review. [IN-03 §11.6; IN-08, Risks] |
| High iteration counts versus controlled schedules | Force both important orders first; repeat and collide afterward | Avoids relying on rare scheduling luck | Stress still samples only some interleavings; counts are not proof probabilities. [IN-17 Patterns 2, 9] |
| Goroutine census versus precise ownership | Actual exit/registration witnesses primary; independent stack/heap evidence supplementary | Separates known resources from process noise and catches inventory omissions | Counters can lie; census can be brittle; neither alone proves release. [IN-10, Evidence execution requirements] |
| Fake clocks versus accelerated real time | No general clock abstraction; barriers for order, real context expiry or verified test-time facilities for deadlines | Small runtime and meaningful context integration | Real expiry needs generous bounds; accelerated sleeps do not prove ordering. [IN-01 §18; IN-16, shutdown-test evidence] |
| Exhaustive matrices versus maintainable suites | Exhaust the tiny state/cursor domains; use focused interaction fixtures and generated traces elsewhere | Strong coverage without a Cartesian product of every representation detail | Generators and adapters need independent review; missed dimensions remain possible. [IN-17 Patterns 5, 7] |
| Strict cause policy versus convenient sentinel matching | Test the dependency's strict accepted-cause rule explicitly | Prevents an accidental cancellation override | Work returning only `ctx.Err()` after a distinct custom cause may fail; tests and docs must show this. [IN-09, Self-critique] |
| Full race coverage versus suite speed | Keep semantic scenarios under `-race`; separate non-race measurements and longer exploration | Maintains the accepted Go mitigation | Runtime cost must be measured; expensive mandatory cases cannot quietly become advisory. [IN-03 §§11.6–11.8] |
| Proof records versus planning overhead | One requirement-linked evidence ledger using ordinary test/CI artifacts | Traceable acceptance with little infrastructure | It must be maintained when tests move; no dataset platform or generic conformance engine is earned. [IN-04, Evidence-Grounded Planning] |

## Risks

| Risk | Consequence | Control |
| --- | --- | --- |
| False-positive helpers | Green tests establish nothing about the claimed boundary | Named broken fixtures and production-path mutants; require specific failures. [IN-03 §7] |
| Race-detector overconfidence | Race-free but illegal histories, deadlocks, or leaks are accepted | Independent semantic, progress, and release oracles. [IN-17, Notable Absences] |
| Tests coupled to internal locks | Refactoring breaks tests or tests accidentally serialize away bugs | Assert committed behavior; use only narrow semantic checkpoints and challenge their effectiveness. [IN-08, Self-critique] |
| Sleeps hide schedules | Slow/flaky tests still miss the target race | Barriers establish order; watchdogs only bound failure. [IN-17 Pattern 2] |
| Goroutine-count brittleness | Background activity causes flakes or broad ignores hide leaks | Precise ownership plus narrowly scoped stack/retention checks. [IN-10, Evidence execution requirements] |
| Premature negative assertion | Test says “still blocked” before the defective path had a chance to run | Acknowledged checkpoints or verified scheduling facilities; do not rely on an immediate empty receive. [IN-03 §7] |
| Model and implementation share the same error | Both agree on incorrect policy | Independent table, hand-authored trajectory fixtures, and mutants violating the PRD. [IN-01 §13; IN-17 Pattern 7] |
| Random seed mistaken for schedule replay | Stress failures cannot be reproduced | Preserve observed checkpoints and operation history; add controlled regression. [IN-10, Evidence execution requirements] |
| Weak test cleanup | Intentionally blocked work or waiters pollute later leak checks | Release test-owned barriers and join harness goroutines on every exit. [IN-03 §11.3] |
| Registration leak without a goroutine | Goroutine sweep passes while completed runs remain parent-retained | Registration accounting and live-parent retention fixtures. [IN-10, Verification obligations] |
| Hostile diagnostic methods | Blocking or panicking user methods interfere with resolution outside the lock | Preserve the dependency's stable, terminating inspection contract; test supported reentrancy and record containment limits. [IN-09, Self-critique] |
| Abnormal-exit handling remains vague | `Goexit` or failed capture becomes zero-value success or strands ownership | Sprint 1 must pin exposure and test cleanup/containment before acceptance. [IN-09, Verification obligations] |
| Expensive suites quietly disappear | Phase-exit evidence ceases to protect future changes | Named owner/cadence, actual selection checks, no silent skip or advisory conversion. [IN-17; IN-19] |
| Measurements imply unsupported bounds | Three/four events are described as bounded memory or unlimited observers | Separate count, retained bytes, caller concurrency, and measured workload. [IN-03 §§11.7–11.8; IN-11, Self-critique] |
| CLI exit masks library leaks | Smoke passes while an embedded host accumulates support | In-process quiescence checks before test-process exit. [IN-16, runc Approach Model; IN-10] |
| Report-only evidence is overstated | Upstream test presence or old prototypes are called current proof | Preserve provenance and require fresh executed evidence. [IN-03 §§4, 15; IN-07] |
| Reasoning-standard compliance is assumed | Project review passes despite an unread required input | Inspect the exact named README before approval. [IN-12; IN-05] |

## Self-critique

**Which claimed invariant has no independent oracle?** No implemented invariant has been verified here. In the proposed design, whole-program private authority and the completeness of the resource inventory are the hardest to establish independently: tests cannot enumerate every future bypass or hidden reference. External-surface checks, actual mutation/spawn/registration review, and independent stack/heap evidence supplement the behavioral model. Opaque-ID collision tests also provide regression evidence, not a mathematical uniqueness proof. [IN-01 §§8–9; IN-03 §§11.1, 11.3; IN-08, Risks]

**Which race test only makes failure less likely rather than controlling it?** The mixed collision and churn lanes remain scheduler-sampled. Their seeds control workloads, not interleavings. Each named critical ordering has a deterministic obligation, but the actual sprint implementation must demonstrate that its checkpoints open the intended window rather than merely delaying a goroutine near it. [IN-17 Pattern 9; IN-10, Evidence execution requirements]

**What defect could pass every proposed test?** A synchronization defect confined to an unmodeled interleaving or unsupported platform could survive. So could a retained reference omitted from both the ownership ledger and sampled retention fixtures, or a new lifecycle field copied incorrectly after the mutation inventory becomes stale. These are reasons to review each new mutation/resource boundary and preserve exploration, not to claim exhaustive proof from a finite suite. [IN-03 §§11.2–11.8, 12.10]

**Can a clean test process hide a goroutine retained until process exit?** Yes. A test may finish successfully while the Go process still owns blocked support, and process termination then hides it. Release checks must run before process exit with parents still alive, actual support-exit witnesses, and an independent leak sweep. Reference-only retention additionally needs registration and heap evidence. [IN-10, Verification obligations; IN-16, ownership comparisons]

**Could a coherent-snapshot test manufacture the property it claims to prove?** Yes. If the hook or adapter takes a test lock around all access, it may hide a broken production publication protocol. The required split-publication negative control must remain observable through the same supported read/wait paths. A test that cannot fail against that control is not atomicity evidence. [IN-03 §11.6; IN-08, Self-critique]

**Are the proposed suite budgets authoritative performance facts?** No. The repetition values are starting proposals to make sprint planning concrete. Exact commands, selections, runtime bounds, and ongoing cadence must be finalized from the realized suite. Reducing cost may change repetitions or eliminate duplicate fixtures; it may not remove a required semantic ordering or race-enabled acceptance case without explicit governance. [IN-02, Evidence policy; IN-03 §11.6]

**Has verification expanded Phase 1?** The reference model, controlled work, test-only fault seams, CLI execution, and resource probes serve current requirements. They do not introduce cleanup hooks, runtime watchdogs, forced stop, storage, workflows, provider calls, a universal executor, or a second runtime. Test-process execution of the diagnostic CLI is not subprocess supervision as an Aren product capability. [IN-01 §§4, 20–24; IN-04 §§3.2, 3.10]

**What remains before a passing claim is possible?** The required reasoning README must be inspected; sprint reasoning must resolve concrete APIs, abnormal-exit exposure, publication and release mechanics, and executable test selections; the actual implementation must pass the linked positive and negative evidence; and the governed review must return `pass`. This document establishes how to challenge those claims. It establishes neither that the implementation is correct nor that project reasoning has passed review. [IN-05, Project Reasoning Policy; IN-02, Phase Exit Gate; IN-12]
