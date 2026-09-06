# Project reasoning template: Events, observation, waiting, and replay

## Purpose

Synthesize evidence about canonical lifecycle history, ordered publication, multiple readers, waiting, replay, and observer independence. Separate in-memory replay from durable replay and high-volume streaming.

Read `projects/aren-phase-01-execution-lifecycle/reasoning/README.md` and follow its required reasoning standard.

## Inputs Used

Complete this table before analysis. Record only material actually read.

| Input ID | Input | Kind | Exact workspace-relative path | Exact portion inspected | How it was used | Limits or applicability warning |
| --- | --- | --- | --- | --- | --- | --- |
| `[IN-NN]` | `[name]` | `[requirement, project reasoning, study report, study repository, prototype report, current code, prior review]` | `[path]` | `[report section; source file and lines or symbol; document heading]` | `[question, comparison, constraint, or challenge]` | `[report-only, product mismatch, scale mismatch, unverified interpretation, or none]` |

Name each study report by exact path and section. When repository behavior matters, add a separate row for the study repository file and line range or symbol. A study directory, report directory, or repository name alone is not evidence.

## Governing questions

- What is the canonical event record for one run?
- Who allocates sequence and at what point?
- Is event order defined by sequence, commit order, or time?
- What does observer delivery guarantee?
- How does a reader move from retained history to live observation without a gap?
- When does a live observer finish?
- Can a waiter or observer influence execution?

## Normative constraints

| Constraint | Source | Fixed rule | Open question |
| --- | --- | --- | --- |
| `[constraint]` | `[path and section]` | `[rule]` | `[question]` |

## Evidence

Use `Inputs Used` IDs in every material claim. Cite the exact report section and, where relevant, the study repository file and line range or symbol. Distinguish report interpretation from direct repository observation.

### Observation model comparison

| Evidence | Canonical record | Ordering key | Replay model | Live delivery | Slow consumer policy | Terminal completion rule | Phase fit |
| --- | --- | --- | --- | --- | --- | --- | --- |
| `[source]` | `[record]` | `[key]` | `[model]` | `[model]` | `[policy]` | `[rule]` | `[fit]` |

### Wait model comparison

| Evidence | Completion signal | Multiple waiters | Late wait | Wait cancellation | Outcome sharing | Leak risk |
| --- | --- | --- | --- | --- | --- | --- |
| `[source]` | `[signal]` | `[behavior]` | `[behavior]` | `[behavior]` | `[mechanism]` | `[risk]` |

Explain which sources expose message streams, provider deltas, durable logs, or distributed subscriptions that exceed Phase 1 lifecycle history.

## Canonical record versus delivery

Define these separately:

| Concern | Required Phase 1 truth | Claim deliberately not made |
| --- | --- | --- |
| Canonical recording | `[guarantee]` | `[non-guarantee]` |
| Replay | `[guarantee]` | `[non-guarantee]` |
| Live delivery | `[guarantee]` | `[non-guarantee]` |
| Deduplication | `[guarantee]` | `[non-guarantee]` |
| Retention | `[guarantee]` | `[non-guarantee]` |
| Process restart | `[guarantee]` | `[non-guarantee]` |

## Event vocabulary analysis

| Candidate event | Fact represented | State transition? | Required payload | Duplicate canonical recording allowed? | Sprint introduced |
| --- | --- | --- | --- | --- | --- |
| `run.created` | `[fact]` | `[yes or no]` | `[payload]` | `[rule]` | `[sprint]` |
| `run.started` | `[fact]` | `[yes or no]` | `[payload]` | `[rule]` | `[sprint]` |
| `run.cancellation_requested` | `[fact]` | `[yes or no]` | `[payload]` | `[rule]` | `[sprint]` |
| Terminal events | `[fact]` | `[yes]` | `[payload]` | `[rule]` | `[sprint]` |

Challenge whether every event adds information. Avoid events that merely echo internal function calls.

## Replay and live-handoff schedules

Analyze at least:

1. observer starts before work;
2. observer starts between two nonterminal events;
3. observer starts during terminal commitment;
4. observer starts after completion;
5. observer requests a sequence equal to the next event;
6. observer abandons without draining;
7. one observer is slow while another is current;
8. replay exposes an event previously delivered live.

For each schedule, state the cursor, returned events, blocking behavior, completion condition, and deduplication identity.

## Immutability and data ownership

Examine event payload copying, slices and maps, timestamps, outcome references, and data returned by history accessors. Explain how Go callers are prevented from mutating canonical history.

## Project conclusions

| Conclusion | Evidence basis | Sprint 1 minimum | Sprint 2 extension | Reopen trigger |
| --- | --- | --- | --- | --- |
| `[conclusion]` | `[evidence]` | `[minimum]` | `[extension]` | `[trigger]` |

## Trade-Offs

Address pull cursors versus channels, replayable history versus memory bounds, snapshot access versus live subscription, independent observer state versus centralized fan-out, and exact delivery claims versus simple at-least-once observation.

## Risks

Include missed replay/live boundary events, terminal close before event visibility, observer backpressure controlling execution, unbounded history, mutable payloads, duplicate delivery mistaken for duplicate recording, and durable-replay concepts entering an in-memory phase.

## Verification obligations

| Claim | Schedule or fixture | Assertion | Failure it catches |
| --- | --- | --- | --- |
| `[claim]` | `[schedule]` | `[assertion]` | `[defect]` |

## Self-critique

- Which observer schedule remains underspecified?
- Can a late observer reconstruct the same canonical history as an early observer?
- Does any public channel require the producer to wait for a consumer?
- Did streaming-system evidence cause the design to solve output deltas that Phase 1 excludes?
