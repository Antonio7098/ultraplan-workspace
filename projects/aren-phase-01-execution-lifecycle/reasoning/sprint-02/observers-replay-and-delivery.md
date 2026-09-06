# Area reasoning template: Observers, replay, and delivery

> Selected template: `projects/aren-phase-01-execution-lifecycle/reasoning/sprint-02/observers-replay-and-delivery.md`

## Purpose

Decide Sprint 2 multiple-observer behavior, sequence cursors, replay and live handoff, terminal completion, abandonment, and delivery guarantees. The canonical in-memory history remains the source of truth. This document must not introduce persistence, a global event bus, or model-output streaming.

## Inputs Used

Complete this table before analysis. Record only material actually read.

| Input ID | Input | Kind | Exact workspace-relative path | Exact portion inspected | How it was used | Limits or applicability warning |
| --- | --- | --- | --- | --- | --- | --- |
| `[IN-NN]` | `[name]` | `[requirement, project reasoning, handbook, study report, study repository, prototype report, current code, Sprint 1 reasoning or review]` | `[path]` | `[report section; source file and lines or symbol; document heading]` | `[decision, comparison, constraint, or challenge]` | `[report-only, product mismatch, scale mismatch, unverified interpretation, or none]` |

Name each study report by exact path and section. When repository behavior matters, add a separate row for the study repository file and line range or symbol. A study directory, report directory, or repository name alone is not evidence.

## Prior Sprint Decisions Applied

| Sprint 1 history decision | Realized behavior | Status | Sprint 2 extension | Compatibility effect |
| --- | --- | --- | --- | --- |
| `[decision]` | `[behavior]` | `[preserved, extended, superseded, unaffected]` | `[extension]` | `[effect]` |

## Observer use cases

| Use case | Start point | Needs replay? | Needs live follow? | Completion condition | Abandonment behavior |
| --- | --- | --- | --- | --- | --- |
| Observer before run starts | `[sequence]` | `[yes or no]` | `[yes or no]` | `[condition]` | `[behavior]` |
| Observer during run | `[sequence]` | `[yes or no]` | `[yes or no]` | `[condition]` | `[behavior]` |
| Observer after terminal | `[sequence]` | `[yes or no]` | `[yes or no]` | `[condition]` | `[behavior]` |
| Replay from sequence | `[sequence]` | `[yes]` | `[policy]` | `[condition]` | `[behavior]` |

## Candidate observation models

| Option | Replay source | Live wait mechanism | Per-observer state | Producer backpressure | Abandonment cleanup | API complexity |
| --- | --- | --- | --- | --- | --- | --- |
| Pull cursor over retained history | `[source]` | `[mechanism]` | `[state]` | `[behavior]` | `[cleanup]` | `[cost]` |
| Per-observer channel subscription | `[source]` | `[mechanism]` | `[state]` | `[behavior]` | `[cleanup]` | `[cost]` |
| Shared notification plus pull reads | `[source]` | `[mechanism]` | `[state]` | `[behavior]` | `[cleanup]` | `[cost]` |
| Other credible model | `[source]` | `[mechanism]` | `[state]` | `[behavior]` | `[cleanup]` | `[cost]` |

Evaluate exact Phase 1 needs. Do not choose a channel because streaming examples use channels.

## Cursor contract

Define whether the cursor means next requested sequence, last observed sequence, or an opaque position. Then complete:

| Requested cursor | Retained history | Active or terminal | Events returned | Blocks after replay? | Next cursor | Error |
| --- | --- | --- | --- | --- | --- | --- |
| Zero/default | `[history]` | `[state]` | `[events]` | `[yes or no]` | `[cursor]` | `[error]` |
| Existing sequence | `[history]` | `[state]` | `[events]` | `[yes or no]` | `[cursor]` | `[error]` |
| Next sequence | `[history]` | `[state]` | `[events]` | `[yes or no]` | `[cursor]` | `[error]` |
| Beyond next sequence | `[history]` | `[state]` | `[events]` | `[yes or no]` | `[cursor]` | `[error]` |
| Terminal plus exhausted cursor | `[history]` | terminal | `[events]` | no | `[cursor]` | `[error]` |

## Replay and live linearization

Describe the exact mechanism that prevents a gap between reading retained history and waiting for future events. Mark the linearization point for observer registration or notification generation.

Analyze these schedules:

1. an event commits between replay read and live wait;
2. terminal commits while an observer registers;
3. an observer wakes but another observer consumes nothing;
4. replay returns an event previously seen live;
5. observer context ends while a notification is pending;
6. a slow observer never reads again.

## Delivery claims

| Concern | Canonical guarantee | Per-observer guarantee | Explicit non-guarantee | Deduplication method |
| --- | --- | --- | --- | --- |
| Recording | `[guarantee]` | `[guarantee]` | `[non-guarantee]` | `(run_id, sequence)` |
| Replay | `[guarantee]` | `[guarantee]` | `[non-guarantee]` | `(run_id, sequence)` |
| Live notification | `[guarantee]` | `[guarantee]` | `[non-guarantee]` | `[method]` |
| Terminal completion | `[guarantee]` | `[guarantee]` | `[non-guarantee]` | `[method]` |

## Resource and abandonment analysis

| Resource per observer | Bound | Released when | Owner | Slow-observer policy | Leak test |
| --- | --- | --- | --- | --- | --- |
| `[cursor, channel, waiter, goroutine, buffer]` | `[bound]` | `[condition]` | `[owner]` | `[policy]` | `[test]` |

## Area Decisions

| Decision | Exact behavior | Evidence and rationale | Rejected alternative | Failure-mode proof |
| --- | --- | --- | --- | --- |
| `[decision]` | `[behavior]` | `[basis]` | `[alternative]` | `[test]` |

## Trade-Offs

Analyze pull simplicity versus immediate delivery, per-observer isolation versus allocation, duplicate delivery versus loss, cursor clarity versus API size, and terminal observer release versus retained post-terminal replay.

## Evidence

Use `Inputs Used` IDs in every material claim. Cite selected study and prototype reports by exact path and section. Cite relevant NATS, durable-workflow, or other study repository mechanisms by file and line range or symbol. Distinguish report interpretation from direct repository observation. Trace decisions to project-wide observation reasoning, Sprint 1 history behavior, and current code. Explain why durable delivery machinery is not copied into an in-memory per-run model.

## Risks

Include replay/live gaps, observer-controlled backpressure, notification loss mistaken for event loss, duplicate replay mistaken for duplicate recording, a terminal observer that never completes, unbounded observer resources, and global-stream concepts entering Phase 1.

## Verification obligations

| Observer schedule | Controlled checkpoints | Expected events and cursor | Completion behavior | Leak assertion |
| --- | --- | --- | --- | --- |
| `[schedule]` | `[control]` | `[result]` | `[behavior]` | `[assertion]` |

## Self-critique

- Identify the exact replay-to-live linearization point.
- Can an observer miss a committed event without detecting a gap?
- Can an abandoned observer delay work or terminal publication?
- Which claimed guarantee actually belongs to durable delivery and must be removed?
