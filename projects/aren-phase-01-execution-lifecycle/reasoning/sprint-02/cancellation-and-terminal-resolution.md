# Area reasoning template: Cancellation and terminal resolution

> Selected template: `projects/aren-phase-01-execution-lifecycle/reasoning/sprint-02/cancellation-and-terminal-resolution.md`

## Purpose

Decide the complete Sprint 2 cancellation protocol and the deterministic interpretation of completion facts. This document owns cancellation request acceptance, cause retention, acknowledgement, disposition, parent cancellation, and the terminal-resolution truth table. It does not own synchronization implementation, observer delivery, or package restructuring.

## Inputs Used

Complete this table before analysis. Record only material actually read.

| Input ID | Input | Kind | Exact workspace-relative path | Exact portion inspected | How it was used | Limits or applicability warning |
| --- | --- | --- | --- | --- | --- | --- |
| `[IN-NN]` | `[name]` | `[requirement, project reasoning, handbook, study report, study repository, prototype report, current code, Sprint 1 reasoning or review]` | `[path]` | `[report section; source file and lines or symbol; document heading]` | `[decision, comparison, constraint, or challenge]` | `[report-only, product mismatch, scale mismatch, unverified interpretation, or none]` |

Name each study report by exact path and section. When repository behavior matters, add a separate row for the study repository file and line range or symbol. A study directory, report directory, or repository name alone is not evidence.

## Prior Sprint Decisions Applied

Classify every relevant Sprint 1 decision.

| Sprint 1 decision | Source | Status | Why | Sprint 2 consequence |
| --- | --- | --- | --- | --- |
| `[decision]` | `[reasoning path or review evidence]` | `[preserved, extended, superseded, unaffected]` | `[evidence]` | `[consequence]` |

Silence does not supersede a Sprint 1 decision. A superseding decision must identify new evidence and implementation impact.

## Decision scope

- explicit caller cancellation;
- supplied parent-context cancellation;
- one acceptance path;
- `accepted`, `already_requested`, and `already_terminal` disposition;
- first accepted cause and reason;
- exactly one cancellation-request event;
- success after accepted cancellation;
- cancellation-related error after acceptance;
- unrelated error after acceptance;
- panic after acceptance;
- ignored or delayed cancellation.

## Cancellation fact model

| Fact | Produced by | Committed by | Observable when | Can exist without next fact? | Mutability rule |
| --- | --- | --- | --- | --- | --- |
| Request attempted | `[actor]` | `[owner or none]` | `[time]` | `[yes or no]` | `[rule]` |
| Request accepted | `[actor]` | `[owner]` | `[time]` | `[yes or no]` | `[rule]` |
| First cause retained | `[actor]` | `[owner]` | `[time]` | `[yes or no]` | `[rule]` |
| Work context cancelled | `[actor]` | `[owner]` | `[time]` | `[yes or no]` | `[rule]` |
| Work returns | `[actor]` | `[owner captures]` | `[time]` | `[yes or no]` | `[rule]` |
| Terminal cancellation committed | `[actor]` | `[owner]` | `[time]` | `[yes or no]` | `[rule]` |

## Parent-context integration

Decide:

- how an already-cancelled parent is handled;
- whether work still enters `running`;
- how parent cause becomes an Aren cancellation reason;
- whether a watcher goroutine is required;
- how watcher startup and teardown avoid leaks;
- how simultaneous explicit and parent requests select the first cause;
- what happens after terminal commitment.

## Disposition contract

| Run facts at request | Returned disposition | State change | Event appended | Cause changed | Work context action |
| --- | --- | --- | --- | --- | --- |
| Active, no request | `accepted` | `[state rule]` | `[event]` | `[rule]` | `[action]` |
| Active, request already accepted | `already_requested` | none | none | no | `[action or none]` |
| Terminal | `already_terminal` | none | none | no | none |

State whether disposition is a value, error, or both, and why callers need that distinction.

## Terminal-resolution inputs

List raw captured facts. Do not branch on which goroutine reached the mutation boundary first.

| Input | Source | Trust condition | Ambiguous form | Normalization required |
| --- | --- | --- | --- | --- |
| Result | `[source]` | `[condition]` | `[form]` | `[normalization]` |
| Error | `[source]` | `[condition]` | `[form]` | `[normalization]` |
| Panic | `[source]` | `[condition]` | `[form]` | `[normalization]` |
| Accepted cause | `[source]` | `[condition]` | `[form]` | `[normalization]` |
| Context error/cause match | `[source]` | `[condition]` | `[form]` | `[normalization]` |
| Existing terminal | `[source]` | `[condition]` | `[form]` | `[normalization]` |

## Terminal-resolution truth table

Cover all material combinations.

| Accepted cancellation | Work result | Work error | Panic | Error matches run cause | Existing terminal | Outcome | Failure origin/kind | Reason |
| --- | --- | --- | --- | --- | --- | --- | --- | --- |
| `[fact]` | `[fact]` | `[fact]` | `[fact]` | `[fact]` | `[fact]` | `[state]` | `[classification]` | `[policy]` |

Include nil success after acceptance, wrapped cause match, generic `context.Canceled` from an unrelated context, unrelated error, result plus error, panic, and late terminal proposals.

## Error matching analysis

Define how Aren recognizes acknowledgement of its own cancellation cause. Analyze `errors.Is`, `context.Cause`, wrapped errors, `DeadlineExceeded`, custom causes, and false matches. If exact truth is unavailable, choose the conservative classification and say what is lost.

## Area Decisions

| Decision | Exact rule | Evidence and rationale | Rejected alternative | Test oracle |
| --- | --- | --- | --- | --- |
| `[decision]` | `[rule]` | `[basis]` | `[alternative]` | `[oracle]` |

## Trade-Offs

Analyze cause precision versus caller ergonomics, success-after-request truthfulness versus user expectation, stable dispositions versus a smaller API, parent watcher cost versus semantic unity, and conservative failure versus optimistic cancellation classification.

## Evidence

Use `Inputs Used` IDs in every material claim. Cite selected study and prototype reports by exact path and section. Cite relevant study repository files and line ranges or symbols for Go context and cancellation mechanisms. Distinguish report interpretation from direct repository observation. Trace decisions to project-wide cancellation and outcome reasoning and Sprint 1 realized behavior. Separate process cancellation evidence from in-process applicability.

## Risks

Include reporting cancellation before work stops, classifying any `context.Canceled` as Aren cancellation, replacing the first cause, duplicate cancellation events, parent cancellation bypassing acceptance, terminal outcome depending on scheduler order, and a controller API that implies forceful interruption.

## Verification obligations

| Rule | Controlled schedule | Expected disposition | Expected terminal outcome | Failure detected |
| --- | --- | --- | --- | --- |
| `[rule]` | `[schedule]` | `[disposition]` | `[outcome]` | `[defect]` |

## Self-critique

- Can two identical committed fact sets resolve differently because of schedule?
- Can Aren publish `cancelled` while work is still running?
- Which wrapped error produces the greatest classification ambiguity?
- Does the API make cancellation look stronger than cooperative context propagation?
