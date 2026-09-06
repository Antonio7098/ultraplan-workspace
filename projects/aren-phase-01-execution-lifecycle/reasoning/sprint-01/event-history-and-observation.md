# Area reasoning template: Event history and minimum observation

> Selected template: `projects/aren-phase-01-execution-lifecycle/reasoning/sprint-01/event-history-and-observation.md`

## Purpose

Decide the minimum retained lifecycle history and inspection behavior Sprint 1 must implement. Keep multiple live observers, cursor subscriptions, cancellation events, and abandoned-observer hardening in Sprint 2 unless Sprint 1 correctness requires a smaller seam.

## Inputs Used

Complete this table before analysis. Record only material actually read.

| Input ID | Input | Kind | Exact workspace-relative path | Exact portion inspected | How it was used | Limits or applicability warning |
| --- | --- | --- | --- | --- | --- | --- |
| `[IN-NN]` | `[name]` | `[requirement, project reasoning, handbook, study report, study repository, prototype report, current code, prior review]` | `[path]` | `[report section; source file and lines or symbol; document heading]` | `[decision, comparison, constraint, or challenge]` | `[report-only, product mismatch, scale mismatch, unverified interpretation, or none]` |

Name each study report by exact path and section. When repository behavior matters, add a separate row for the study repository file and line range or symbol. A study directory, report directory, or repository name alone is not evidence.

## Decision scope

- `run.created`, `run.started`, `run.succeeded`, and `run.failed`;
- per-run sequence starting at zero;
- event append within lifecycle transition commitment;
- immutable history snapshots;
- history access before and after completion;
- explicit boundary between Sprint 1 inspection and Sprint 2 live observation.

## Inherited constraints

| Input | Constraint or evidence | Applied rule | Question remaining |
| --- | --- | --- | --- |
| `[input]` | `[finding]` | `[rule]` | `[question]` |

## Event purpose test

Every retained event must describe committed lifecycle truth that a caller cannot obtain as reliably from a single current-state read.

| Candidate event | Fact represented | Transition or occurrence | Required by PRD? | Payload needed | Keep, merge, or defer |
| --- | --- | --- | --- | --- | --- |
| `run.created` | `[fact]` | `[kind]` | `[yes or no]` | `[payload]` | `[decision]` |
| `run.started` | `[fact]` | `[kind]` | `[yes or no]` | `[payload]` | `[decision]` |
| `run.succeeded` | `[fact]` | `[kind]` | `[yes or no]` | `[payload]` | `[decision]` |
| `run.failed` | `[fact]` | `[kind]` | `[yes or no]` | `[payload]` | `[decision]` |

## Event schema

| Field | Type and invariants | Allocated by | Copied or normalized | Reason caller needs it | Future compatibility pressure |
| --- | --- | --- | --- | --- | --- |
| Run identity | `[contract]` | `[owner]` | `[rule]` | `[reason]` | `[pressure]` |
| Sequence | `[contract]` | `[owner]` | `[rule]` | `[reason]` | `[pressure]` |
| Type | `[contract]` | `[owner]` | `[rule]` | `[reason]` | `[pressure]` |
| Occurrence time | `[contract]` | `[owner]` | `[rule]` | `[reason]` | `[pressure]` |
| State transition | `[contract]` | `[owner]` | `[rule]` | `[reason]` | `[pressure]` |
| Type-specific data | `[contract]` | `[owner]` | `[rule]` | `[reason]` | `[pressure]` |

Reject generic metadata maps unless a concrete Phase 1 field cannot be represented directly.

## Ordering and commitment

Define sequence allocation, timestamp capture, append timing, and relation to state publication. Explain why timestamps never decide order.

| Schedule | History allowed | State allowed | Outcome allowed | Forbidden mismatch |
| --- | --- | --- | --- | --- |
| Before start commit | `[history]` | `[state]` | `[outcome]` | `[mismatch]` |
| During terminal commit | `[history]` | `[state]` | `[outcome]` | `[mismatch]` |
| After terminal commit | `[history]` | `[state]` | `[outcome]` | `[mismatch]` |

## Candidate Sprint 1 APIs

Compare snapshot access, iterator access over a fixed snapshot, and an early cursor API. Select the smallest API that does not force Sprint 2 to break a valid Sprint 1 contract.

| Option | Multiple callers | Mutation safety | Blocks? | Claims live delivery? | Sprint 2 compatibility | Cost |
| --- | --- | --- | --- | --- | --- | --- |
| `[option]` | `[behavior]` | `[mechanism]` | `[yes or no]` | `[claim]` | `[impact]` | `[cost]` |

## Deferred behavior contract

State exactly what Sprint 1 does not promise:

- live subscription;
- independent cursors;
- replay from arbitrary sequence;
- slow-observer behavior;
- observer completion signaling;
- cancellation-request event;
- retention after process lifetime.

If any seam for later behavior is retained, show the current requirement that earns it.

## Area Decisions

| Decision | Exact behavior | Evidence and rationale | Rejected alternative | Sprint 2 extension point |
| --- | --- | --- | --- | --- |
| `[decision]` | `[behavior]` | `[basis]` | `[alternative]` | `[extension]` |

## Trade-Offs

Cover history snapshots versus allocations, concrete event fields versus extensibility, returning full history versus a bounded API, and an intentionally narrow Sprint 1 API versus likely Sprint 2 change.

## Evidence

Use `Inputs Used` IDs in every material claim. Cite each study report by exact path and section and each relevant study repository mechanism by file and line range or symbol. Distinguish report interpretation from direct repository observation. Trace schema and ordering decisions to PRD requirements, project-wide reasoning, prototype defects, and current code. Separate durable event-log evidence from in-memory applicability.

## Risks

Analyze mutable history, sequence allocated outside transition commitment, duplicated terminal events, timestamp ordering, public generic payloads, speculative cursor abstractions, and a Sprint 1 snapshot API that makes Sprint 2 semantics misleading.

## Verification obligations

| Claim | Test schedule | Mutation or alias check | Failure detected |
| --- | --- | --- | --- |
| `[claim]` | `[schedule]` | `[check]` | `[defect]` |

## Self-critique

- Does each event report a committed fact rather than an intention?
- Can callers mutate canonical history through any nested value?
- Which Sprint 2 requirement would force a breaking change to the chosen API?
- Did this document accidentally design a general event bus?
