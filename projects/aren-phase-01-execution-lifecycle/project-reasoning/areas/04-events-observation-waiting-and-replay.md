I’ll define the cursor and terminal-drain contract against the prior publication and resource-ownership decisions, while keeping live observation separate from Sprint 1’s retained-history minimum.
# Events Observation Waiting And Replay

> Project: `aren-phase-01-execution-lifecycle`
> Area: Events Observation Waiting And Replay
> Scope: Cross-sprint canonical history, coherent inspection, shared waiting, independent observation, and in-memory replay
> Evidence boundary: Supplied governed documents and report excerpts; no direct repository inspection or test execution

## Purpose

Define how multiple readers observe one run without becoming participants in its execution. The selected semantic model is **one retained canonical lifecycle history, independently advancing readers, and terminal readiness derived from the complete committed terminal record**. Reading history does not execute work, reconstruct runtime state, acknowledge lifecycle transitions, or control retention. [IN-03 §§14–19; IN-06, Candidate authority models and Downstream obligations]

This area consumes the prior authority, outcome, and resource-ownership conclusions. It settles cursor boundaries and completion behavior without selecting Go packages, exported signatures, a specific notification primitive, or a general event-stream abstraction. Sprint 1 needs coherent history snapshots and shared waiting; Sprint 2 adds live observation, sequence cursors, and abandonment verification. [IN-04, both implementation waves; IN-06–IN-08]

## Inputs Used

Repeated injections of the same path are one input. Report-cited repository locations below are second-hand evidence, not separately inspected source files. No omitted report middle, underlying repository file, standalone prototype report, current Aren implementation, or test execution is claimed as inspected.

| Input ID | Input | Kind | Exact workspace-relative path | Exact portion inspected | How it was used | Limits or applicability warning |
| --- | --- | --- | --- | --- | --- | --- |
| IN-01 | Project index | Governance and requirement | `projects/aren-phase-01-execution-lifecycle/project-index.md` | Project Reasoning Policy; Project Scope; evidence catalog; unavailable-dimension table; Prior Decisions; carry-forward policy | Scope, review gate, evidence restrictions | Catalog availability is not runtime proof |
| IN-02 | Project-reasoning index | Governance | `projects/aren-phase-01-execution-lifecycle/project-reasoning/index.md` | Reasoning Areas; this area's evidence and source assignments; Excluded Evidence | Dependencies and decision ownership | Routing text is not comparative evidence |
| IN-03 | Phase 1 PRD | Requirement | `projects/aren-phase-01-execution-lifecycle/docs/PRD.md` | Full supplied copy, especially §§7–19 and §§20–25 | Authoritative lifecycle, event, waiting, observation, timing, and acceptance contract | Normative, not implementation evidence |
| IN-04 | Project roadmap | Requirement and delivery policy | `projects/aren-phase-01-execution-lifecycle/roadmap.md` | Scope Principle; Canonical Project Flow; Cross-Sprint Carry-Forward Rule; both implementation waves; Phase Exit Gate | Sprint minimums, extension, and promotion boundary | Planned commands have not been executed |
| IN-05 | Evidence Assessment And Routing | Dependency reasoning | `projects/aren-phase-01-execution-lifecycle/project-reasoning/areas/00-evidence-map.md` | Supplied Purpose, Inputs Used, Scope, assessment tables, and ending containing Trade-Offs, Evidence, Risks, and Self-critique | Apply evidence restrictions and negative-transfer findings | Excerpt; report-only assessment |
| IN-06 | Lifecycle Authority And Atomic Publication | Dependency reasoning | `projects/aren-phase-01-execution-lifecycle/project-reasoning/areas/01-lifecycle-authority-and-atomic-publication.md` | Supplied Purpose, Governing questions, Normative constraints, Evidence, Candidate authority models, beginning of Atomicity model, and ending containing Trade-Offs, Risks, Downstream obligations, and Self-critique | One authority, coherent reads, publication before waiter return, private capability boundary | Semantic decisions, not executed proof |
| IN-07 | Outcomes Failures And Terminal Resolution | Dependency reasoning | `projects/aren-phase-01-execution-lifecycle/project-reasoning/areas/02-outcomes-failures-and-terminal-resolution.md` | Supplied Purpose, Governing questions, Normative constraints, Evidence, Candidate outcome models, Terminal fact model, visible Resolution matrix, and ending containing Trade-Offs, Risks, Verification obligations, and Self-critique | Valid terminal branches, stable diagnostics, result/cause ownership, panic-event obligations | Excerpt; exact representation remains a sprint choice |
| IN-08 | Cancellation Goroutine Ownership And Cleanup | Dependency reasoning | `projects/aren-phase-01-execution-lifecycle/project-reasoning/areas/03-cancellation-goroutine-ownership-and-cleanup.md` | Supplied Purpose, Governing questions, Normative constraints, Evidence, Cancellation state model through visible acceptance-to-signal introduction, and ending containing Risks, Verification obligations, and Self-critique | Acceptance versus terminality, local abandonment, explicit release, terminal publication versus quiescence | Excerpt; watcher and notification mechanisms are not selected |
| IN-09 | Accepted Go decision | Requirement with historical evidence | `projects/aren-phase-01-execution-lifecycle/docs/final-language-decision.md` | Full initial supplied copy, especially §§4, 6–7, 11, 14–15 | Defensive publication, owned goroutines, declared bounds, historical defects, fresh verification | Historical prototype checks do not verify the new implementation |
| IN-10 | Selected area template | Reasoning requirement | `projects/aren-phase-01-execution-lifecycle/reasoning/project-wide/04-events-observation-waiting-and-replay.md` | Full supplied copy | Required comparisons, schedules, ownership analysis, and critique | Referenced reasoning README unavailable in supplied context |
| IN-11 | Delivery guarantees and idempotency | Study report | `studies/agent-harness-study/reports/final/01.09-delivery-guarantees-and-idempotency.md` | Executive Summary; Approach Models A–C; Pattern Catalog P1–P8; supplied Open Questions and Evidence Index | Separate recording, repeat observation, deduplication, and effects | Report-only; distributed retry and persistence dominate |
| IN-12 | Replay and determinism | Study report | `studies/agent-harness-study/reports/final/01.10-replay-and-determinism.md` | Executive Summary; Approach Models A–B; Pattern Catalog P1–P5 and visible beginning of P6; supplied ending and Evidence Index | Distinguish retained reading from durable reconstruction and re-execution | Report-only; broad model-reexecution claims are restricted by IN-05 |
| IN-13 | Snapshot and checkpoint architecture | Study report | `studies/agent-harness-study/reports/final/02.02-snapshot-and-checkpoint-architecture.md` | Executive Summary; Approach Models A–D; Pattern Catalog 1–12; supplied Evidence Index tail | Committed-only snapshots and negative transfer from checkpoint machinery | Report-only; no persistence or recovery adoption |
| IN-14 | Public API surface | Study report | `studies/agent-harness-study/reports/final/24.01-public-api-surface.md` | Availability disclosure; Executive Summary; Approach Models; Pattern Catalog P1–P6; supplied Per-Source Notes and Evidence Index | Enforced capability boundaries and negative surface tests | Only two material TypeScript sources; no direct Go API proof |
| IN-15 | Go lifecycle ownership and arbitration | Study report | `studies/aren-go-runtime-study/reports/final/01.01-lifecycle-transition-ownership-and-terminal-arbitration.md` | Executive Summary; Approach Models; P1–P8; Key Differences; Tradeoffs; supplied failure-mode tail, Notable Absences, and Evidence Index | Outcome-before-notification, completion funnel, waiter and mock limitations | Report-only; source analyses lacked Aren PRD |
| IN-16 | Go cancellation, ownership, and cleanup | Study report | `studies/aren-go-runtime-study/reports/final/01.02-cancellation-goroutine-ownership-and-cleanup.md` | Executive Summary; Core Thesis; Approach Models; P1–P8; supplied Open Questions and Evidence Index | Distinguish join, close, publication, repeated waiting, and producer abandonment | Report-only; process escalation and worker shutdown excluded |
| IN-17 | Go ordered observation and backpressure | Study report | `studies/aren-go-runtime-study/reports/final/01.04-ordered-observation-live-streaming-and-backpressure.md` | Executive Summary; Core Thesis; Approach Models; Pattern Catalog 1–6; supplied Evidence Index | Sequence identity, atomic handoff, observer isolation, and upstream-loss counterexamples | Report-only; transport/high-volume mismatch and attribution caveat retained |

The template requires reading `projects/aren-phase-01-execution-lifecycle/reasoning/README.md`. Its content was not supplied, and no read-only workspace tool is available in this turn. It is not an inspected input. Compliance with any additional standard there remains unverified and must be checked before a passing project-reasoning review. This document claims neither that compliance nor a review verdict. [IN-10; IN-05, Purpose]

## Governing questions

| Question | Phase 1 answer | Basis |
| --- | --- | --- |
| What is canonical? | The retained sequence of lifecycle occurrences committed by the private authority for one run | IN-03 §§5, 14–16; IN-06 |
| Who allocates sequence? | The same authority that commits the event and associated lifecycle facts | IN-03 §§9, 14–15 |
| What orders events? | Contiguous per-run sequence, reflecting serialized commitment; timestamps are descriptive | IN-03 §15; IN-06, Normative constraints |
| What does delivery guarantee? | An uninterrupted reader advancing from a valid cursor receives the requested suffix in order without gaps or unsolicited repeats; intentional replay can expose the same identities again | Project interpretation of IN-03 §16 |
| How does retained reading become live observation? | A reader reaches the current tail and waits through a race-free availability protocol against the same canonical record; there is no second event source | IN-03 §§14–16; IN-17, Pattern 4 |
| When does observation finish normally? | After terminal publication and exhaustion of the requested suffix, including its terminal event when that event is within the suffix | IN-03 §16; cursor decision below |
| Can readers control execution? | No. Reader cancellation, cursor advancement, slowness, and abandonment cannot cancel work, postpone terminal commitment, or gate retention | IN-03 §§9, 16–17; IN-08 |

## Normative constraints

The source column resolves to exact paths in Inputs Used. Open implementation questions do not permit weaker behavior.

| Constraint | Source | Fixed rule | Open question |
| --- | --- | --- | --- |
| Per-run canonical history | IN-03 §§5, 15–16 | Retain committed lifecycle events in memory; no global bus | Private representation |
| Creation identity and sequence | IN-03 §§8, 15 | Identity exists before creation; creation sequence is `0` | Identity representation belongs to Sprint 1 |
| Atomic publication | IN-03 §14; IN-06 | Event, sequence, state, timing, cancellation facts, and outcome cannot be partially exposed | Synchronization and coherent-read surface |
| Event vocabulary | IN-03 §15 | Six lifecycle event types; cancellation request is not a state | Concrete named Go types and validation |
| Shared waiting | IN-03 §17 | Every completed wait observes the same logical outcome; waiting does not finalize or consume it | Blocking API shape; optional local wait cancellation |
| Independent observers | IN-03 §16 | Each reader has its own cursor; no reader gates another | Stateless cursor operation versus reader object |
| Replay boundary | IN-03 §16 | Requested sequence and default `0` are supported | Inclusive and out-of-range behavior are settled below |
| Terminal drain | IN-03 §§14–17 | Normal observation completion cannot precede terminal-event visibility | Concrete end-of-observation representation |
| Immutable lifecycle data | IN-03 §§15, 17; IN-09 §11.2 | No supported mutation can alter canonical lifecycle containers or payloads | Defensive copies versus immutable internal sharing |
| Arbitrary result exception | IN-03 §17; IN-07, Terminal fact model | Exact successful result is not universally deep-frozen | Result accessor shape and documentation |
| Cause and panic diagnostics | IN-07, Terminal fact model and Self-critique | Stable lifecycle diagnostics; original error/cause uses explicit read-only borrowing | Exact diagnostic types, without serialization infrastructure |
| Owned resource release | IN-09 §11.3; IN-08 | No observer-dependent producer lifetime; support has explicit release paths | Concrete notification/registration lifetime |
| Retention lifetime | IN-03 §16 | History remains available while the run is reachable | Measurements of retained bytes, not an eviction policy |
| Delivery shape | IN-04 | Sprint 1 snapshots and waiters; Sprint 2 live cursors and abandonment | Realized APIs and test names |
| Phase boundary | IN-03 §§4, 19, 25 | No persistence, durable replay, output streaming, acknowledgements for effects, or universal executor | None within this area |

## Evidence

### Observation model comparison

All source behavior in this table is reported behavior. None was independently reproduced here.

| Evidence | Canonical record | Ordering key | Replay model | Live delivery | Slow consumer policy | Terminal completion rule | Phase fit |
| --- | --- | --- | --- | --- | --- | --- | --- |
| IN-17, Crush Approach Model and Patterns 1–3 | Durable message state; notifications are not a complete lifecycle log | No shared sequence cursor in the cited event envelope | Refetch/resynchronize rather than replay lost notifications | Lossy publish plus bounded-blocking “must deliver” class | Per-subscriber drop or timeout | Correlated `RunComplete`, with final text as reconciliation data | Useful negative contrast; neither lossy canonical history nor terminal delivery waits fit Aren |
| IN-17, Docker Agent Approach Model and Patterns 4–5 | Sequenced per-session ring at the server boundary | Monotonic `seq` | Retained backlog plus registered listener; explicit gap outside retention window | Buffered per-listener delivery | Full listener closes and is removed | Sequenced `session_exited` and stream completion | Adapt atomic handoff; reject eviction, gap protocol, session scope, and upstream-loss assumptions |
| IN-17, NATS Approach Model | JetStream store; core subscriptions have a different guarantee | Store and per-consumer sequence/cursor state | Durable consumer resume and redelivery | Isolated outbound staging and flow control | Disconnect or stall; some producer throttling | No single run-terminal primitive | Adapt separation of canonical sequence from delivery progress; reject broker, acknowledgement, and flow-control machinery |
| IN-17, Ollama Approach Model; IN-16 P7 | No retained replay log in the cited request pipe | Channel and transport order | None | Coupled unbuffered response channel | Client pressure propagates into generation; abandonment hazards reported | Final `done:true` response can be lost on disconnect | Counterexample for tying execution to consumption |
| IN-11, Agent Framework Approach Model C and Evidence Index | Reliable-stream sample backed by Redis Streams | Chunk sequence/cursor | Repeated exposure with client deduplication | At-least-once sample delivery | Not sufficient in the excerpt to establish Aren-style isolation | Sample-specific, not an Aren terminal record | Transfer stable identity and recording/delivery distinction only |
| IN-12, Temporal Approach Model B; IN-13, Approach Model A | Persisted workflow history and derived mutable state | Recorded event order within history lineage | State reconstruction and recovery | Durable-engine delivery, not a small local reader | Storage/worker-specific | Recorded workflow termination | Negative scope boundary; no rebuilder or branch architecture is earned |
| IN-03 §§14–17; IN-06; project decision below | One retained in-memory run history | `(run_id, sequence)` | Inclusive reading of an unchanged suffix | Independent pull-oriented availability | Slow reader advances later; no eviction or execution backpressure | Terminal committed and requested suffix exhausted | Selected Phase 1 semantics; implementation still requires proof |

Docker Agent's single-lock handoff is the closest concrete mechanism, but its ring is downstream of a lossy application fan-out and its runtime observer boundary can block synchronously. Consequently, a gapless server-log subscription does not prove that the log contains every lifecycle fact produced upstream. Aren must record directly at its lifecycle authority, before any delivery mechanism. [IN-17, Rating Summary, Docker Agent Approach Model, Patterns 4–6; IN-05, Selection and canonicality assessment]

### Wait model comparison

| Evidence | Completion signal | Multiple waiters | Late wait | Wait cancellation | Outcome sharing | Leak risk |
| --- | --- | --- | --- | --- | --- | --- |
| IN-15, Conc Approach Model; IN-16, Conc Approach Model | Owner's scoped join | `WaitGroup` and pool APIs must be distinguished; the pool wait closes/reset resources and concurrent calls have a reported double-close hazard | Scope reuse is not immutable run identity | Cooperative work cancellation; no general local wait-abort contract established | Private collection exposed at owner wait | Uncooperative work prevents return; owner-side finalization is not Aren observation |
| IN-15, Temporal SDK P1–P2 and Approach Model | Future/completion slot and later command publication | Workflow-runtime consumers, not proof of arbitrary concurrent Go callers | Resolved future semantics within its runtime | SDK-specific | One slot, with reported production/test repeat-completion divergence | Coroutine and worker ownership cannot be inferred from future readiness |
| IN-16, Docker Agent P2 | Guaranteed stream close on the cited path | A consumed channel is not a broadcast outcome store | Closed channel does not replay prior events | Request/consumer context behavior | Events plus close, not retained immutable outcome sharing | Producer/forwarder lifetime and blocked consumer paths still need ownership |
| IN-03 §17; IN-06; IN-08 | Readiness after complete terminal publication | Required, without a unique consuming receiver | Immediate logical outcome availability | Optional local abort must be distinct from run cancellation | Same immutable logical container; arbitrary successful result remains user-owned | No Aren per-wait producer is earned; actual support needs release proof |

Waiting is therefore **observation of independently completed publication**, not a shutdown command, queue drain, or exclusive receive of the only result. A capacity-one result channel can protect one abandoned receiver from stranding a producer, but cannot by itself satisfy shared waiting and late retrieval. [IN-03 §17; IN-16 P6; IN-06, Candidate authority models]

### Load-bearing anchors

These are locations quoted by inspected reports, not additional inspected repository inputs.

| Bounded lesson | Inspected report section | Report-cited repository location | Transfer limit |
| --- | --- | --- | --- |
| Backlog capture and live registration need one ordering boundary | IN-17, Pattern 4 | `studies/aren-go-runtime-study/sources/docker-agent/pkg/server/eventlog.go:146-160`; `studies/aren-go-runtime-study/sources/docker-agent/pkg/server/eventlog_test.go:152-180` | Local handoff evidence, not upstream completeness or an Aren pull-cursor implementation |
| Observer isolation cannot repair upstream dropped facts | IN-17, Docker Agent Approach Model | `studies/aren-go-runtime-study/sources/docker-agent/pkg/app/app.go:956-988`; `studies/aren-go-runtime-study/sources/docker-agent/pkg/runtime/observer.go:14-18` | Reject these loss/blocking boundaries for canonical lifecycle events |
| “Must deliver” naming does not establish lossless delivery | IN-17, Crush Approach Model and Pattern 5 | `studies/aren-go-runtime-study/sources/crush/internal/pubsub/broker.go:201-236` | Bounded blocking and timeout drop still violate an observer-independent commit path |
| Consumption can strand a producer | IN-16 P7; IN-17, Ollama Approach Model | `studies/aren-go-runtime-study/sources/ollama/server/routes.go:649-792`; `studies/aren-go-runtime-study/sources/ollama/llm/llama_server.go:1703-1708` | Disconnect severity is partly inferred in the ownership report; no need to assume it to reject coupling |
| Outcome precedes derived notification | IN-15 P6 | `studies/aren-go-runtime-study/sources/conc/pool/context_pool.go:40-47` | Adapt ordering, not sibling orchestration |
| Waiting and closure ownership must be separate | IN-16, Executive Summary and Conc Approach Model | `studies/aren-go-runtime-study/sources/conc/pool/pool.go:72-82` | Reported close/reset hazard is not a claim about every Conc wait API |
| Snapshots capture committed state | IN-13, Pattern 1 | `studies/agent-harness-study/sources/agent-framework/dotnet/src/Microsoft.Agents.AI.Workflows/Execution/StateManager.cs:245-248`; `studies/agent-harness-study/sources/temporal/service/history/workflow/mutable_state_impl.go:7539-7546` | No checkpoint storage or recovery transfer |
| Repeat delivery needs stable identity | IN-11, Approach Model C and Evidence Index | `studies/agent-harness-study/sources/agent-framework/dotnet/samples/04-hosting/DurableAgents/AzureFunctions/08_ReliableStreaming/RedisStreamResponseHandler.cs:88-93` | No Redis or reliable-stream subsystem adoption |
| Public boundaries benefit from negative assertions | IN-14, Pattern P3 | `studies/agent-harness-study/sources/openhands/__tests__/library-entrypoints.test.ts:10-49` | TypeScript package evidence, not an exact Go capability layout |

### Evidence restrictions

The source reports expose provider deltas, message streams, persisted workflow histories, session rings, and distributed subscriptions. These systems face data volumes, disconnections, acknowledgement protocols, and restart requirements absent from Phase 1. Their two-path publisher, debounce, ring eviction, heartbeat, durable cursor, and resynchronization machinery is not selected. Phase 1 has one event class: required lifecycle facts. [IN-03 §§4, 15–16; IN-11–IN-13; IN-17, Patterns 1–6]

The public-API report has only two material sources. The replay report's broad statements about model re-execution are not adopted. The Go observation report's ambiguous blocking-sink attribution is avoided by using its explicitly attributed Docker Agent model and source paths. Neither report ratings nor repeated citations to the same implementation are independent proof. [IN-05, Selection and canonicality assessment; IN-12, Executive Summary; IN-14, Availability disclosure; IN-17, Core Thesis]

## Canonical record versus delivery

| Concern | Required Phase 1 truth | Claim deliberately not made |
| --- | --- | --- |
| Canonical recording | Each accepted lifecycle occurrence is appended once with its committed facts; one completed run has one terminal event | Exactly-once arbitrary work execution or transactional effects |
| Replay | Read the retained requested suffix with its original identities, payloads, and times | Re-execute work, replay commands, or reconstruct authoritative runtime state |
| Live delivery | A continuing reader receives each requested canonical event in order and can reach terminal exhaustion without another reader's cooperation | Eventual processing by an abandoned, unscheduled, or failed consumer |
| Deduplication | `(run_id, sequence)` identifies an event consistently across history, live reading, and replay | Consumer side effects occur once merely because identities are available |
| Retention | Full canonical history remains available while the run is reachable; no observer acknowledgement is required | Global memory bound across arbitrary caller-retained runs or result graphs |
| Process restart | No retained-history guarantee survives process termination | Durable replay, reconnect across restart, recovery, or idempotent run creation |

These guarantees follow from the PRD; the precise advancing-cursor rule below is this area's project interpretation. “At-least-once observation” is acceptable only as shorthand for possible repeat exposure through replay, not an unconditional delivery promise. The more useful contract is gapless ordered reading within one traversal and stable identity across traversals. [IN-03 §§15–19; IN-11, Executive Summary and P4; IN-12, Approach Models]

The history and current run facts share one authority. A history snapshot is a committed prefix, not an eventually consistent mirror. Separate calls can observe different valid prefixes while a run advances; they do not become a multi-call transaction. Where one operation returns history with state, timing, or outcome, those fields must describe the same committed version. [IN-06, Evidence conflicts and Trade-Offs; IN-13, Pattern 1]

## Event vocabulary analysis

Every event carries run identity, event type, sequence, and committed occurrence time. State-transition endpoints appear where applicable. Replay preserves these values rather than assigning delivery times or new sequence numbers. [IN-03 §15; IN-06, Normative constraints]

| Candidate event | Fact represented | State transition? | Required payload | Duplicate canonical recording allowed? | Sprint introduced |
| --- | --- | --- | --- | --- | --- |
| `run.created` | Aren created this run and identity | Initialization, not an edge from a fabricated prior state | Common envelope; destination/initial state `created`; no invented source state | No; exactly once at sequence `0` | Sprint 1 |
| `run.started` | Aren entered running before invoking work | `created -> running` | Common envelope, endpoints, committed start time | No; exactly once at sequence `1` | Sprint 1 |
| `run.cancellation_requested` | First cancellation accepted; Aren owes propagation | No state change; occurrence while `running` | Common envelope; running context; immutable accepted-source/reason projection consistent with retained cause metadata | At most one; rejected/repeated requests allocate no sequence | Sprint 2 |
| `run.succeeded` | Nil-error work return resolved successfully | `running -> succeeded` | Common envelope, endpoints, terminal timing and success classification; exact result remains accessible through the outcome | No; mutually exclusive with other terminal events | Sprint 1 |
| `run.failed` | Work error, work panic, or containable Aren invariant fault resolved as failure | `running -> failed` | Common envelope, endpoints, terminal timing, origin/kind, stable diagnostics; panic classification and captured stack where applicable | No; mutually exclusive with other terminal events | Sprint 1 |
| `run.cancelled` | Finished work matched the accepted cancellation cause | `running -> cancelled` | Common envelope, endpoints, terminal timing, first accepted cancellation metadata consistent with the outcome | No; mutually exclusive with other terminal events | Sprint 2 |

The event projection must preserve the outcome area's distinctions without duplicating every outcome object into history. A success event need not contain an arbitrary mutable result graph. A failure event must remain useful for identifying origin/kind and panic diagnostics; it need not serialize the original Go error graph. A cancellation event must not replace the first accepted reason with the returned acknowledgement error. The exact fields and accessors are sprint choices within these requirements. [IN-07, Terminal fact model, Verification obligations, and Self-critique]

Each event adds a distinct fact. Creation establishes identity; start establishes the invocation boundary and timing; acceptance records a request that may precede success, failure, or cancellation; the terminal event records the resolved outcome. No `observer.attached`, `wait.started`, `event.delivered`, `context.signalled`, repeated-request, or internal-finalizer event is added. Such events would describe implementation activity, expand history with observer behavior, and undermine the small retention argument. [IN-03 §§6, 12–16; IN-08, Cancellation state model; project conclusion]

For a correctly represented completed run, history contains three events without accepted cancellation and four with it. Active runs contain the corresponding prefix. This is an event-count bound derived from the lifecycle and first-acceptance rule, not a byte bound or permission to silently drop a fifth append caused by a defect. An impossible append is an invariant violation handled under the prior containment contract. [IN-03 §§7, 12–15; IN-07, Risks and Verification obligations; IN-09 §§11.7–11.8]

## Cursor contract

The project selects an **inclusive next-sequence cursor**. Let `n` be the number of events in a coherent committed history; present sequences are `0` through `n-1`. A cursor `c` names the next event the reader requests. These are semantic names, not exported Go fields or signatures. [Project interpretation of IN-03 §§15–16]

| Boundary | Selected behavior |
| --- | --- |
| Default cursor | `c = 0` |
| Opening with `0 <= c < n` | Valid; return retained events starting with exactly `c` |
| Opening with `c = n`, run active | Valid; wait for the next append or local read cancellation |
| Opening with `c = n`, run terminal | Valid; immediately report normal exhaustion |
| Opening with `c > n` | Explicit invalid-cursor error; do not clamp, skip, or wait for a hypothetical future sequence |
| Negative/unrepresentable cursor, where the API permits such input | Explicit invalid input; never wrap into another sequence |
| Returning event sequence `s` | Reader's next cursor becomes `s+1`; a batch advances to one past its last returned event |
| No event returned because the read is locally cancelled | Cursor does not advance |
| Intentional replay | Open another traversal at an earlier valid sequence; no canonical mutation |
| Two observers | Independent cursors; advancing one has no effect on the other |
| Concurrent calls on one mutable reader object | Not an additional Phase 1 requirement; callers serialize that object's advancement unless sprint reasoning explicitly earns and tests stronger behavior |

Opening validation is serialized against the current committed prefix. Thus a concurrent append may determine whether an opening cursor is at the tail or beyond it. The operation's answer is valid at that boundary; it is not guaranteed to remain an up-to-date description when the caller later receives it. Requiring a captured valid cursor avoids speculative future positions and indefinite waiting for an impossible sequence. [IN-06, Governing questions and coherent-read constraints; project decision]

Normal exhaustion means **the run is terminal and no event in the requested suffix remains**. It does not mean the current active history happens to be empty at the cursor. A reader opened at a completed tail has intentionally requested an empty suffix; it can finish without redelivering a terminal event that lies before its cursor. This does not weaken terminal visibility: the terminal record must already be committed and retrievable. [IN-03 §§14–17; project interpretation]

A blocking live read needs an observation-local cancellation path so a caller can leave an active tail wait without requesting run cancellation. Exact context or reader-close API shape belongs to Sprint 2. A local abort is not normal terminal exhaustion, does not append an event, and does not discard canonical history. [IN-03 §§9, 16; IN-08, Verification obligations; project decision]

For a read whose cancellation races with availability, prefer already available canonical data, then terminal exhaustion, then local cancellation at the read's coherent decision point. A cancellation wake must recheck availability before reporting abort. This avoids a `select` scheduling accident hiding a ready terminal suffix, while preserving caller-local cancellation when no event or terminal exhaustion is available. No cursor advances on the abort path. [Project decision derived from IN-03 §§14–17; IN-06, publication discipline]

## Replay and live-handoff schedules

Notation: `E[s]` is canonical event sequence `s`, `t` is the terminal sequence, and `c` is the next requested sequence. All returned events retain deduplication identity `(run_id, s)`. Blocking below means waiting for future availability, not routine synchronization needed to inspect a committed prefix. These schedules are proof obligations, not executed results. [IN-03 §§15–17; cursor contract above]

| Schedule | Cursor | Returned events | Blocking behavior | Completion condition | Deduplication identity |
| --- | --- | --- | --- | --- | --- |
| 1. Observer starts before work | `0`, after handle availability with work held at a controlled invocation boundary | Creation and start from history, then later request/terminal events | May wait only at the active tail; execution never waits for registration | After returning the requested suffix through `E[t]` | Original `(run_id, s)` |
| 2. Observer starts between two nonterminal events | `0`, or `2` after creation/start and before accepted cancellation | From `0`: retained creation/start followed by request; from `2`: request when committed | Tail read waits until request or terminal; no synthetic request is required | Terminal suffix drained | Original identities |
| 3. Observer starts during terminal commitment | Any currently valid `c` | Either a preterminal prefix followed by the terminal event, or the already complete suffix | May serialize behind the commit; cannot see terminal state with missing terminal data | Terminal suffix drained, never readiness alone | Original identities |
| 4. Observer starts after completion | `0` | Complete canonical history through `E[t]` | No wait for future events | Exhaustion after final event | Same identities as an early observer |
| 5. Observer requests the next event | `c = n` | Next append at sequence `n`, if active | Waits while active; if already terminal, immediate empty-suffix exhaustion | Once terminal and cursor is beyond its terminal event | Returned events keep assigned identities; exhaustion has none |
| 6. Observer abandons without draining | Any valid `c` | Only events already returned to that observer | No producer may remain waiting for its consumption; a pending read exits through its local cancellation path | Local abandonment, not successful lifecycle exhaustion | Later replay uses unchanged identities |
| 7. One observer is slow while another is current | Slow `c_s < c_f`; fast near current tail | Slow reader receives its retained suffix later; fast receives newly committed events | Slow reader's inactivity does not block the fast reader, work, cancellation, or waiters | Each drains its own suffix independently | Shared canonical IDs, independent progress |
| 8. Replay exposes a previously live-delivered event | New traversal starts at `c <= s` after earlier delivery of `E[s]` | Includes the same `E[s]` again | Retained events require no future-event wait | According to the new traversal's tail and terminal state | Exactly the same `(run_id, s)`, not a new event |
| 9. Opening cursor exceeds the current tail | `c > n` | No event | Immediate invalid-cursor result at validation boundary | Not terminal exhaustion | None |
| 10. Local cancellation races with terminal availability | Current active-tail cursor | Terminal event if available at the decision point; otherwise local abort without advancement | No work cancellation or observer-dependent terminal delay | Normal exhaustion only after the requested terminal suffix; abort remains distinguishable | Terminal identity unchanged |
| 11. Work ignores accepted cancellation | Reader at tail after request event | No invented terminal event; later actual completion is delivered | Reader may remain pending or cancel its own read | Only actual terminal publication and suffix drain | Request keeps its original identity |
| 12. Repeated cancellation after completion | Any observer cursor | No new event from repeated requests | No new wait introduced | Existing terminal exhaustion unchanged | Existing identities unchanged |

“Before work” must not become an externally actionable `created` scheduling state. The PRD allows work to start or finish before a handle returns. Tests can hold the work body at an entry barrier or use a narrow internal invocation seam; production execution must never await observer attachment. Physical subscription between creation and start is not a required public operation. Both events remain available through replay. [IN-03 §§7, 10, 16; IN-04, Sprint 1 Deferred]

### Handoff argument

A pull-oriented reader has no separate backlog channel and live channel to reconcile. It repeatedly consults the same record. Nevertheless, reaching the tail and arranging to wait still creates a lost-wakeup risk. Sprint 2 must prove this protocol: [IN-17, Pattern 4; IN-06, Downstream obligations]

1. Inspect the cursor, committed history, terminal status, and the applicable availability notification through one coherent synchronization protocol.
2. If requested data exists, return it and advance only this reader.
3. If the run is terminal and the suffix is exhausted, return normal completion.
4. Otherwise establish a wait that covers every append after the inspected prefix, then release lifecycle ownership before blocking.
5. After notification or local cancellation, recheck canonical availability and terminality. A wakeup is not an event or evidence of termination.

An append before inspection is found in history. An append after the waiting boundary is established must wake the reader or leave a persistent readiness condition. The forbidden interval is an append after observing an empty tail but before acquiring the notification that should represent that append. Docker Agent's locked backlog/registration operation demonstrates one way to eliminate this class of interval; Aren's concrete pull mechanism must supply its own Go memory-order argument. [IN-17, Pattern 4; IN-09 §11.6]

Notifications may coalesce because history retains the facts. Canonical events may not coalesce or drop. The terminal predicate cannot be checked ahead of unread history in a way that discards the final suffix. No callback, blocking event send, or join executes under lifecycle ownership. [IN-03 §§14–17; IN-06, Evidence conflicts; IN-08, Normative constraints]

## Waiting and completion

A completed outcome wait returns only after the supervised invocation has finished and state, terminal timing, outcome, and terminal event are fully published. If cancellation was accepted, the prior area's propagation-settlement requirement also applies. Waiting does not need to wait for observers to receive or process the terminal event. [IN-03 §§10–17; IN-07, Terminal fact model; IN-08, Governing questions and Self-critique]

Normal outcome waiting and event-stream exhaustion answer different questions:

| Operation | What completion establishes | What it does not establish |
| --- | --- | --- |
| Outcome wait | Complete authoritative terminal outcome is available | Every observer drained history |
| History snapshot | One committed prefix was captured | Run is terminal unless the snapshot's committed facts say so |
| Normal live-observation exhaustion | Run is terminal and this traversal's requested suffix is exhausted | Consumer side effects completed or any other observer caught up |
| Locally cancelled read/wait | This caller stopped waiting | Run cancellation was requested or work stopped |
| Internal resource-release observation | The named Aren-owned support reached its release boundary | Arbitrary user-spawned goroutines or external effects ceased |

The PRD does not require a context-cancellable outcome-wait API. If Sprint 1 or Sprint 2 exposes one, its abort must be caller-local and distinct from a terminal outcome. It should use the same coherent ready-before-abort principle rather than initiate cancellation or finalization. An already published outcome remains available to later and concurrent waiters. [IN-03 §§9, 17; IN-08, Verification obligations; project extension rule]

Outcome readiness is not automatically a join of every Aren support goroutine. The prior resource area permits only a finite internal epilogue independent of consumer behavior and requires separate release proof. Observation must neither strengthen that claim accidentally nor weaken the release obligation by leaving producers or registrations parked indefinitely. [IN-08, Self-critique; IN-09 §11.3]

## Immutability and data ownership

The ownership boundary is stronger than copying a top-level struct and narrower than freezing every object reachable from user code. Aren-owned lifecycle data must remain immutable across every supported accessor, snapshot, waiter, and cursor. Arbitrary successful results and borrowed original diagnostic causes retain the explicit exceptions established upstream. [IN-03 §§15, 17; IN-07, Terminal fact model; IN-09 §11.2]

| Data | Owner and publication rule | Mutation challenge |
| --- | --- | --- |
| Canonical history storage | Private run authority; never expose writable backing storage | Replace returned entries, reorder them, append within spare capacity, then reread history |
| Event envelopes and transition fields | Value data or immutable private representation | Change a returned event's identity, sequence, type, endpoints, or time |
| Nested event collections | Aren-owned immutable storage; defensive copies on any mutable access | Mutate nested slices/maps through one observer and verify another is unchanged |
| Start, finish, and occurrence times | Captured once into committed facts, returned by value | Verify all readers see identical times; replay must not resample |
| Outcome containers referenced from events or snapshots | Private immutable representation or independent defensive snapshots | Mutate any public container copy, then wait/replay again |
| Panic diagnostics | Stable classification and captured stack; no writable canonical stack slice | Mutate returned stack bytes or frames and compare subsequent event/outcome access |
| Failure and cancellation projections | Aren-owned stable metadata and diagnostic text | Check that presentation does not recompute changing fields from borrowed objects |
| Original errors and custom causes | Stable borrowed references under read-only, safe-inspection contract | Assert local cause inspection remains possible; do not claim arbitrary implementations are deep-frozen |
| Successful result | Exact captured result, including reference identity where applicable | Demonstrate documented aliasing without allowing that alias to mutate lifecycle metadata |
| Observer cursor | Observer-local progress, not a mutable field in canonical events | Advance/replay one reader and verify another's next sequence and history remain unchanged |

A history accessor returning a slice copy is insufficient if its elements still contain writable slices, maps, pointers to mutable lifecycle containers, or raw mutable panic payloads. Immutable internal sharing is acceptable only when the full supported API leaves no mutation path; otherwise defensive copies are required. No universal reflection-based copier or JSON round trip is justified. [IN-09 §§7, 11.2, 11.10; IN-07, Trade-Offs]

User-defined `Error`, `Is`, `Unwrap`, formatting, or rendering behavior must not execute inside lifecycle mutation ownership or history-copy critical sections. Observation should use the stable diagnostic projections prepared under the outcome contract. Cause inspection performed by a caller is not a reason to rerun classification or modify canonical history. [IN-06, Risks; IN-07, Risks and Self-critique]

An observation-only surface must not expose a cancellation function, public controller through an obvious dynamic assertion, private transition method, writable history, or closable internal readiness channel. Capability narrowing needs tests against the actual delivered value, not merely a smaller static interface. This is API authority separation, not a sandbox against unsafe in-process code. [IN-03 §9; IN-06, Risks; IN-14, Pattern P3]

## Resource and retention model

The selected direction is pull-oriented observation with independent cursors and no mandatory per-observer producer goroutine. A synchronous blocked read runs in its caller's goroutine. Dropping an idle reader requires no drain; cancelling a pending read must release any operation-local waiting registration. The concrete notification mechanism remains Sprint 2's decision. [IN-03 §16; IN-08, Ownership evidence; IN-09 §§11.3, 11.10]

| Resource | Required ownership and bound | Release or retention rule |
| --- | --- | --- |
| Canonical history | At most four events for a correctly completed Phase 1 run | Retain while the run is reachable; no reader-driven eviction |
| Outcome and lifecycle diagnostics | One terminal container plus required stable diagnostic data | Same run lifetime; bytes are not bounded merely by event count |
| Independent cursor | Constant-size progress state per observer, excluding caller-retained copies | Caller owns its lifetime; it must not require a run-owned producer |
| Pending read support | Only concrete synchronization state earned by the implementation | Release on data return, terminal exhaustion, or local abort |
| Availability notification | Run-owned or operation-owned as explicitly chosen | Must not retain abandoned observers or completed runs through external registrations |
| Returned snapshots/events | Caller-owned copies or immutable views | Caller retention is intentional; must not be mistaken for hidden runtime leakage |
| Parent integration and invocation support | Existing run resource owners | Preserve the prior terminal-triggered release and separate quiescence obligations |

A retained reader may legitimately keep the run reachable, just as a retained run handle does. That is different from a run or parent retaining an abandoned reader through a producer, registry, or callback after its useful lifetime. Tests must distinguish required reachability from unintended retention. [IN-03 §16; IN-08, Risks and Verification obligations]

The event-count bound does not bound error text, stack diagnostics, custom-cause graphs, arbitrary results, the number of retained runs, or caller-created observers and snapshots. Sprint measurements must state these categories and their provenance. Silent truncation, ring eviction, and process-wide admission machinery are not introduced to manufacture a stronger bound than the product supports. [IN-09 §§11.7–11.8; IN-07, Risks; IN-08, retention discussion]

## Project conclusions

These are project semantic decisions and proof obligations, not implementation acceptance.

| Conclusion | Evidence basis | Sprint 1 minimum | Sprint 2 extension | Reopen trigger |
| --- | --- | --- | --- | --- |
| EO-01: Record lifecycle facts directly through the run's private authority, before delivery | IN-03 §§14–16; IN-06; IN-17, Docker Agent upstream-loss counterexample | Complete retained history for success, error, and panic | Add request/cancelled occurrences through the same authority | Any delivery path can omit, reorder, or independently mint canonical facts |
| EO-02: Sequence is contiguous, begins at zero, and identifies commitment order rather than delivery time | IN-03 §15; IN-11; IN-06 | Identity and sequence tests for three-event trajectories | Request insertion and collision tests | Gaps, replay renumbering, timestamp sorting, or sequences allocated for rejected requests |
| EO-03: Use inclusive next-sequence cursors; current tail is valid, beyond-tail input is rejected | IN-03 §16; explicit Cursor contract | No live cursor machinery required | Implement boundary behavior and independent traversal | Realized API cannot express the contract cleanly; record any semantic revision explicitly |
| EO-04: Prefer pull-oriented reading of one history over central event fan-out | IN-08, Ownership evidence; IN-17, Patterns 4–6; IN-09 §11.10 | Defensive history snapshots only | Race-free availability waits with no mandatory per-observer producer | A concrete requirement or measured cost justifies another equally honest mechanism |
| EO-05: Normal observation completion requires terminal publication and requested-suffix exhaustion | IN-03 §§14–17; IN-16 P2, restricted | Terminal history is available when outcome waits return | Drain-before-exhaustion, late-tail, and local-abort tests | A reader closes at an active tail or skips its terminal suffix |
| EO-06: Waiting observes one reusable logical outcome and never owns finalization | IN-03 §17; IN-15–IN-16, Conc limits | Multiple early/late waiters, no exclusive receive | Wait/cancel/observer concurrency and optional local wait abort | Repeated waits close resources, consume the only outcome, or wait on observers |
| EO-07: Lifecycle payloads are defensive; arbitrary result and borrowed-cause exceptions remain explicit | IN-07, Terminal fact model; IN-09 §11.2 | Mutation attacks on history, outcomes, and panic diagnostics | Repeat attacks through cursors and multiple readers | A supported alias changes canonical lifecycle data |
| EO-08: Abandonment affects only the observer; no drain or acknowledgement gates execution | IN-03 §16; IN-08; IN-16 P6–P7 | Snapshot access and waiting create no stranded producers | Pending-read abort, idle-reader abandonment, slow/fast isolation, release evidence | Reader inactivity controls execution or retains Aren-owned support indefinitely |
| EO-09: Retain complete lifecycle history; declare the three/four-event bound separately from bytes | IN-03 §§7, 12, 15–16; IN-09 §§11.7–11.8 | Three-event histories and baseline retention evidence | Four-event histories, repeated-request bound, reader churn | New event classes or measured diagnostic retention pressure; later streaming requires separate reasoning |
| EO-10: No exactly-once delivery, durable replay, or consumer-effect guarantee | IN-03 §§16, 19; IN-11–IN-12 | Accurate library documentation | Replay and diagnostic CLI demonstrate stable IDs and repeat exposure | A later governed phase explicitly introduces durable or effect semantics |

Sprint 2 must consume realized Sprint 1 requirements, handbook, area reasoning, authoritative sprint synthesis, plan, and implementation/review evidence. It must identify decisions as preserved, extended, superseded, or unaffected rather than replacing the snapshot/publication design by implication. A counterexample to fixed PRD semantics requires repair or a governed requirement change, not weaker delivery wording. [IN-04, Cross-Sprint Carry-Forward Rule; IN-01]

## Trade-Offs

**Pull cursors versus channels.** Pull-oriented reading avoids a delivery goroutine, per-observer event queue, and send/close ownership for every reader. It still requires a precise tail-wait protocol and a caller-local abort path. Buffered channels can be made correct, particularly for a tiny event count, but introduce storage duplication and lifetime obligations without a demonstrated Phase 1 benefit. A single shared event channel is not broadcast. Exact Go API choice remains with Sprint 2. [IN-03 §16; IN-16 P6–P7; IN-17, Patterns 4–5]

**Replayable history versus memory bounds.** Retaining every lifecycle event removes retention-gap recovery and makes slow readers harmless to recording. It costs memory for reachable runs and their diagnostics. The count is small; total bytes are not universally bounded. A ring would violate the complete reachable-history requirement and solve future output-volume pressure prematurely. [IN-03 §§15–17; IN-09 §§11.7–11.8; IN-17, Docker Agent Approach Model]

**Snapshots versus live subscription.** Snapshots give Sprint 1 a complete useful inspection surface with fewer notification races. They do not satisfy Phase 1's eventual live-observation requirement by themselves. Sprint 2 should extend the same canonical history rather than create a second event source or claim that repeated polling is already a proven gapless subscription. [IN-04, both implementation waves; IN-13, Pattern 1]

**Independent state versus centralized fan-out.** Independent cursors make slow readers fall behind only themselves. Central fan-out can provide channel ergonomics but needs subscriber registration, buffers, shutdown, and abandonment accounting. Shared availability notification is still possible without a registry of per-observer event producers. It must be measured and verified rather than described as cost-free. [IN-03 §16; IN-08, Ownership evidence; IN-17, Pattern 5]

**Exact delivery claims versus repeatable observation.** A single advancing traversal can be gapless and nonduplicating without promising exactly-once consumer processing. Replay intentionally repeats canonical IDs. An unconditional “at least once” promise is also too strong for an abandoned observer. The selected wording is precise about retained availability, traversal behavior, and consumer responsibility. [IN-03 §§16, 19; IN-11, Executive Summary]

**Strict cursor validation versus future-position convenience.** Rejecting `c > n` detects mistakes early and avoids waiting for a terminal sequence that can never exist. Accepting the current tail still supports live following. Validation racing with an append may produce an error that a retry would not, but that is an explicit linearized result rather than silent cursor clamping. [Cursor contract; IN-06, coherent-read constraints]

**Terminal readiness versus complete quiescence.** Releasing outcome waiters without waiting for consumer processing preserves independence. It must not excuse a permanently retained watcher or producer. The prior resource contract's finite internal epilogue and separate release evidence remain necessary. [IN-08, Self-critique; IN-09 §11.3]

## Risks

| Risk | Consequence | Control |
| --- | --- | --- |
| Missed replay/live boundary event | Reader sleeps forever or skips an event despite correct canonical history | One availability protocol; force append at the empty-tail/wait boundary. [IN-17, Pattern 4] |
| Terminal close before event visibility | EOF or outcome readiness is observed without matching terminal history | Complete publication before readiness; drain requested history before exhaustion. [IN-03 §§14–17] |
| Observer backpressure controls execution | Slow or absent consumer delays cancellation, work completion, or waiters | No callback or blocking delivery in commitment; test an unread observer beside a fast one. [IN-03 §16] |
| Event-count bound mistaken for byte bound | Retention claims conceal large diagnostics or user objects | Separate count, payload bytes, caller retention, and run-owned support measurements. [IN-09 §§11.7–11.8] |
| Mutable nested payloads | A consumer changes another reader's event, timing, failure, or stack | Ownership inventory and mutation attacks across every access path. [IN-07; IN-09 §11.2] |
| Duplicate delivery mistaken for duplicate recording | Replay is reported as a lifecycle violation or canonical duplicates are dismissed as replay | Check canonical history separately from per-traversal delivery; compare stable IDs. [IN-03 §§15–16] |
| Durable-replay concepts enter an in-memory phase | Unnecessary stores, acknowledgements, gap protocols, or rebuilders appear | Keep exclusions explicit and use current lifecycle behavior to justify every mechanism. [IN-03 §4; IN-12–IN-13] |
| Current tail mistaken for terminality | Active observation returns normal EOF before work finishes | EOF requires committed terminal status as well as suffix exhaustion. [IN-03 §16] |
| Local abort becomes run cancellation | Merely leaving a read appends a cancellation request or stops another waiter | Separate observation lifetime from control capability and supplied work parent. [IN-03 §9; IN-08] |
| Reader retained by hidden registration | Completed runs or abandoned readers remain reachable | Audit notification registrations and actual exit paths; do not rely on finalizers. [IN-08, Verification obligations] |
| Read-only interface exposes controller dynamically | Observer can recover caller-control authority | Test delivered concrete capability surface, not just method declarations. [IN-06, Risks; IN-14 P3] |
| Joined snapshot assembled from independent reads | History and outcome describe different versions in one claimed coherent operation | Capture related fields through one publication discipline. [IN-06; IN-13, Pattern 1] |
| Channel closure confused with work stop | Observer exit falsely implies a terminal run | Distinguish normal terminal exhaustion, local abort, and internal support release. [IN-08; IN-16 P2, restricted] |
| Evidence breadth overstated | A local handoff test or historical green suite is treated as Phase 1 proof | Preserve report-only labels and require fresh production-path evidence. [IN-05; IN-09 §§4, 15] |
| Additional reasoning standard unchecked | Passing review overlooks requirements in the missing README | Inspect the named reasoning README before approval. [IN-10] |

## Verification obligations

All rows describe required evidence, not executed results. Assertions must operate on the realized runtime and an independent expected trajectory, not a fake that already enforces the desired contract. Controlled barriers establish ordering; timeouts detect lack of progress rather than create the tested order. [IN-03 §21; IN-09 §11.6; IN-15, production/test divergence]

| Claim | Schedule or fixture | Assertion | Failure it catches |
| --- | --- | --- | --- |
| Canonical sequence is contiguous | Success, returned failure, panic; then accepted cancellation variants | Identity constant; creation `0`, start `1`, exactly one terminal; length three or four as applicable | Duplicate recording, skipped sequence, delivery-time allocation |
| Rejected requests do not grow history | Repeat explicit/parent requests before and after completion | At most one request occurrence; unchanged sequence after rejected requests | Recording API calls rather than accepted facts |
| Outcome wait sees complete publication | Release many waiters while terminal commitment is instrumented | Every returned waiter immediately retrieves matching terminal history, state, and timing | Early readiness or partial fact installation |
| Coherent snapshots are committed prefixes | Concurrent snapshot readers at every transition | Each combined view matches one independently valid prefix and outcome state | Separate-field snapshot assembly |
| Replay is inclusive | Open at every valid sequence in completed three- and four-event histories | Exact suffix begins at requested sequence; original times and IDs unchanged | Off-by-one start, replay renumbering, regenerated payloads |
| Tail and invalid cursors differ | Active `c=n`, terminal `c=n`, `c>n`, and invalid representation | Active tail waits; terminal tail exhausts; beyond-tail errors without mutation | Infinite future wait, clamping, active-tail EOF |
| Registration before work does not become scheduling | Hold controlled work entry, then observe; also complete before handle is consumed | Both paths expose creation/start; execution never waits on registration | Observer attachment controls invocation |
| Handoff is gapless | Force append immediately before and after tail-wait establishment | Next requested event is returned without a skipped ID or indefinite wait | Lost wakeup between check and wait |
| Terminal suffix drains | Pause reader with unread events, then complete work | Reader receives every requested event through terminal before normal exhaustion | Terminal flag checked ahead of backlog |
| Late and early readers agree | Early reader from `0`; late reader from `0` after completion | Equal logical canonical histories, including diagnostics and timing | Live-only terminal, history omission, projection drift |
| Cursors are independent | Readers begin at different valid sequences and advance at different rates | Each gets its own expected suffix; no event is consumed globally | Shared receive queue or shared mutable cursor |
| Slow/absent observers cannot block | Leave one reader idle; advance another and wait for outcome | Cancellation, work completion, waiters, and fast reader progress independently | Mandatory consumer delivery in commit path |
| Pending read can be abandoned locally | Cancel a tail read while work and a second reader remain active | No cancellation-request event; first read returns distinct local abort without advancement | Observation context controls work |
| Abort/terminal race is explicit | Force terminal availability before decision; separately force abort before any availability | Ready suffix wins when visible at decision; otherwise abort; later replay remains complete | Random readiness selection or cursor advancement on abort |
| Ignored cancellation remains active | Hold work after accepted request and signal | No terminal event or normal observer exhaustion until actual work release | Stop request or notification treated as termination |
| Replay repeats delivery, not recording | Consume `E[s]` live, then open at `s` again | Same ID and payload delivered again; canonical history unchanged | Duplicate canonical append during replay |
| Containers are defensive | Mutate returned history, event fields, nested diagnostics, and stack bytes | Fresh history/outcome and another reader remain unchanged | Shallow-copy aliasing |
| Result/cause exceptions remain narrow | Reference-bearing successful result and wrapped custom cause | Exact result/local cause semantics retained; lifecycle metadata remains protected | Universal-copy distortion or blanket immutability overclaim |
| Observation capability has no control | Inspect/use the actual observation value from an external consumer test | No supported mutation or recoverable controller capability | Static-only authority narrowing |
| Waiting is shared observation | Multiple early/late waits and repeated calls | Same logical outcome; no double close, exclusive consumption, or finalization by wait | Pool/shutdown semantics copied into outcome waiting |
| Abandoned readers leave no support | Repeated idle-reader drops and pending-read aborts, followed by controlled run completion | Named registrations/carriers release; retained run history remains available | Producer or external registration retention |
| Resource checks distinguish ownership | Retain/drop run handles and snapshots in separate fixtures | Intentional retained data is explained; run-owned support returns to its release baseline | Treating all retained memory as either bounded or leaked |
| Instruments detect seeded defects | Introduce a lost wakeup, early terminal exhaustion, shared payload, or blocking delivery separately | Corresponding assertion fails for the intended reason | Vacuous helpers, unrelated counters, oracle copying implementation |

Sprint 1 must run its actual success/error/panic, snapshot, timing, mutation, negative-transition, shared-waiter, and release tests under `go test -race`. Sprint 2 must add all cursor, handoff, cancellation, observer, abandonment, and broad concurrency cases under the same detector, plus repeated stress and independent resource accounting. A green race run alone cannot prove absence of logical deadlocks, missed notifications, or retained registrations. [IN-04, both implementation waves; IN-03 §§21–24; IN-09 §11.6; IN-08, Evidence execution requirements]

Stress evidence must record the realized revision, worktree status, toolchain, workload, repetitions, and relevant timing/allocation spread. Random seeds are useful but do not reproduce Go scheduling by themselves; important interleavings need deterministic fixtures. Resource verification must release intentionally blocked test work before asserting quiescence and supplement named release observations with targeted leak/retention checks. [IN-09 §§11.6–11.8; IN-08, Evidence execution requirements]

The required diagnostic CLI must exercise the real runtime and show identity, sequence-ordered events, request-versus-terminal cancellation, and the resolved outcome for `success`, `fail`, `cancel`, and `race`. It must not label replayed display as duplicate recording or describe output as durable. Exact commands and status assertions belong to Sprint 2 planning; contract promotion follows phase validation, not this reasoning output. [IN-03 §§20, 24; IN-04, Sprint 2 Commands and Phase Exit Gate]

## Self-critique

**Which observer schedule remains underspecified?** Inclusive start, current tail, beyond-tail input, terminal drain, and local-abort races are settled semantically. Exact Go notification mechanics and concurrent operations on one mutable reader object remain sprint decisions; separate observers are required, while concurrent advancement of one object is not silently promised. The most dangerous implementation gap is still the empty-tail-to-wait interval, which needs a concrete memory-order argument and a forced-schedule test. [Cursor contract; IN-17, Pattern 4; IN-09 §11.6]

**Can a late observer reconstruct the same canonical history as an early observer?** A late observer starting at `0` must read the same logical history, including original identities, times, transition facts, and stable diagnostics. “Reconstruct” here means read and compare retained facts, not rebuild runtime state or reproduce arbitrary result-object contents. Readers starting at another cursor intentionally receive a suffix. [IN-03 §§15–17; IN-07, Terminal fact model]

**Does any public channel require the producer to wait for a consumer?** None is selected. The pull-oriented direction avoids that obligation. If Sprint 2 chooses a channel surface, it must independently prove broadcast, late replay, terminal drain, abandonment, and no execution backpressure; buffering alone is not that proof. [IN-16 P6–P7; IN-17, Patterns 4–5]

**Did streaming-system evidence cause the design to solve output deltas that Phase 1 excludes?** No dual publisher, debounce, throttle, eviction ring, heartbeat, acknowledgement window, or transport reconnection protocol is selected. Their useful contribution is identifying loss and ownership boundaries. The lifecycle remains one retained class of required facts with three or four events in a correctly completed run. [IN-03 §§4, 15; IN-17, Patterns 1–6]

**Is “nonblocking observation” overstated?** It means consumer inactivity cannot create a dependency on execution progress. Coherent reads still use synchronization and copying, and many callers can create contention. This document does not promise lock-free access, real-time latency, infinite observer capacity, or zero allocation. Measurements and the concrete mechanism must make those costs visible. [IN-03 §16; IN-09 §§11.7–11.8]

**Is strict future-cursor rejection the only valid interpretation?** No. Waiting for arbitrary future positions could be made explicit, but introduces impossible-tail and skipping semantics without a current requirement. This area chooses rejection as the smallest clear contract. Changing it is a semantic decision that needs recorded rationale and updated tests, not an accidental difference between snapshot and live implementations. [IN-03 §16; Cursor contract]

**What remains unproven?** All Aren-specific behavior remains unexecuted. The comparative reports establish useful mechanisms and counterexamples, not a complete implementation of this cursor/publication/release contract. The reasoning README remains unchecked, and the realized defensive-copy, notification, capability, and resource-release boundaries must survive sprint reasoning and tests. These limits prevent implementation or review-pass claims; they do not justify adding more platform scope. [IN-01, Project Reasoning Policy; IN-05; IN-09 §§4, 15; IN-10]
