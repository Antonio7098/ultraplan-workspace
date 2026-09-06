# Area reasoning template: Lifecycle transition and publication

> Selected template: `projects/aren-phase-01-execution-lifecycle/reasoning/sprint-01/lifecycle-transition-and-publication.md`

## Purpose

Decide the Sprint 1 state machine, transition gate, linearization point, and publication order for creation, start, success, returned failure, and panic. This document owns lifecycle mutation semantics. It does not own package placement, cancellation resolution, live observer delivery, or the full verification strategy.

## Inputs Used

Complete this table before analysis. Record only material actually read.

| Input ID | Input | Kind | Exact workspace-relative path | Exact portion inspected | How it was used | Limits or applicability warning |
| --- | --- | --- | --- | --- | --- | --- |
| `[IN-NN]` | `[name]` | `[requirement, project reasoning, handbook, study report, study repository, prototype report, current code, prior review]` | `[path]` | `[report section; source file and lines or symbol; document heading]` | `[decision, comparison, constraint, or challenge]` | `[report-only, product mismatch, scale mismatch, unverified interpretation, or none]` |

Name each study report by exact path and section. When repository behavior matters, add a separate row for the study repository file and line range or symbol. A study directory, report directory, or repository name alone is not evidence.

## Decision scope

- Required states: `created`, `running`, `succeeded`, `failed`, with `cancelled` reserved in the vocabulary as required by the PRD.
- Required transitions: creation to running, then exactly one terminal commitment.
- Required coherent facts: state, event, timing, outcome, and waiter release.
- Sprint 2 boundary: cancellation request and cancellation terminal interpretation remain deferred.

## Inherited conclusions

| Input | Exact conclusion or constraint used | Authority | Applied, narrowed, or challenged |
| --- | --- | --- | --- |
| `[PRD, project reasoning, handbook, code]` | `[conclusion]` | `[requirement, evidence synthesis, current fact]` | `[treatment]` |

## State model

Complete the full transition matrix. Include illegal terminal-to-terminal, repeated start, and internal impossible cases.

| Source state | Operation or fact | Destination | Legal? | Facts committed | Error or invariant response |
| --- | --- | --- | --- | --- | --- |
| `[state]` | `[operation]` | `[state or unchanged]` | `[yes or no]` | `[facts]` | `[response]` |

Explain whether `created` is externally observable and whether starting work can fail before `running` is committed.

## Mutation inventory

| Mutable fact | Sole writer | Read paths | Written in which transition | Protected by | Immutable after terminal? |
| --- | --- | --- | --- | --- | --- |
| State | `[owner]` | `[readers]` | `[transition]` | `[mechanism]` | `[yes or no]` |
| Event sequence/history | `[owner]` | `[readers]` | `[transition]` | `[mechanism]` | `[rule]` |
| Timing | `[owner]` | `[readers]` | `[transition]` | `[mechanism]` | `[rule]` |
| Outcome | `[owner]` | `[readers]` | `[transition]` | `[mechanism]` | `[rule]` |
| Completion notification | `[owner]` | `[waiters]` | `[transition]` | `[mechanism]` | `[rule]` |

## Candidate transition mechanisms

| Option | Linearization point | Lock or owner scope | Work placement | Publication ordering | Benefits | Defects or costs |
| --- | --- | --- | --- | --- | --- | --- |
| Direct critical-section commit | `[point]` | `[scope]` | `[inside or outside]` | `[order]` | `[benefits]` | `[costs]` |
| Serialized transition loop | `[point]` | `[scope]` | `[inside or outside]` | `[order]` | `[benefits]` | `[costs]` |
| Other credible mechanism | `[point]` | `[scope]` | `[inside or outside]` | `[order]` | `[benefits]` | `[costs]` |

Choose by the visibility guarantees and failure schedules, not by stylistic preference.

## Atomic publication sequence

Write the exact conceptual order for each transition. Mark the linearization point and explain which steps occur inside the mutation boundary.

```text
1. [validate source and terminal status]
2. [capture transition facts]
3. [construct immutable values]
4. [allocate sequence and append event]
5. [publish state, timing, and outcome]
6. [release waiters after coherent visibility]
```

If the implementation order differs, explain why no supported operation can observe a partial commitment.

## Visibility schedule analysis

| Concurrent reader schedule | Allowed observation | Forbidden observation | Mechanism preventing the forbidden view |
| --- | --- | --- | --- |
| State read during terminal commit | `[view]` | `[view]` | `[mechanism]` |
| History read during terminal commit | `[view]` | `[view]` | `[mechanism]` |
| Wait returns during terminal commit | `[view]` | `[view]` | `[mechanism]` |
| Two terminal producers collide | `[view]` | `[view]` | `[mechanism]` |

## Invariant violation policy

Decide which illegal transitions return a caller-facing error, which indicate an Aren defect, and how an invariant failure remains visible without causing a second terminal mutation.

## Area Decisions

| Decision | Exact rule | Why it satisfies the PRD | Rejected alternative | Required proof |
| --- | --- | --- | --- | --- |
| `[decision]` | `[rule]` | `[reason]` | `[alternative]` | `[test or review]` |

## Trade-Offs

Analyze lock simplicity versus long critical sections, constructing outcome inside versus before the gate, state snapshots versus direct reads, strict invariant failures versus defensive no-ops, and event-first versus state-first internal assignment.

## Evidence

Use `Inputs Used` IDs in every material claim. Cite selected reports by exact path and section and relevant study repository mechanisms by file and line range or symbol. Distinguish report interpretation from direct repository observation. Use exact project-wide conclusions, prototype defects, and current-code references. Explain why durable checkpoint or database transaction mechanisms are relevant only as semantic analogies, if used at all.

## Risks

Cover partial publication, hidden second writers, callbacks while holding a lock, work executed under the mutation gate, outcome construction that can panic, waiter release before event append, and tests that pass only because readers run after completion.

## Verification obligations

| Invariant | Controlled schedule | Assertion | Negative control |
| --- | --- | --- | --- |
| `[invariant]` | `[schedule]` | `[assertion]` | `[deliberate defect]` |

## Self-critique

- Identify the exact instruction that linearizes a terminal transition.
- Could any allocation, callback, or user code run while lifecycle mutation is locked?
- Which partial observation is hardest to rule out?
- Does reserving `cancelled` in the vocabulary accidentally claim Sprint 2 behavior?
