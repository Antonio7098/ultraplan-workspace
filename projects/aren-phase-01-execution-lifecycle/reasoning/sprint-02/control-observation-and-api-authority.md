# Area reasoning template: Control, observation, and API authority

> Selected template: `projects/aren-phase-01-execution-lifecycle/reasoning/sprint-02/control-observation-and-api-authority.md`

## Purpose

Decide the Go API capability split for inspection, waiting, observation, and cancellation. The API must expose required behavior without giving callers transition authority or promising permanent executor abstractions.

## Inputs Used

Complete this table before analysis. Record only material actually read.

| Input ID | Input | Kind | Exact workspace-relative path | Exact portion inspected | How it was used | Limits or applicability warning |
| --- | --- | --- | --- | --- | --- | --- |
| `[IN-NN]` | `[name]` | `[requirement, project reasoning, handbook, study report, study repository, prototype report, current code, Sprint 1 reasoning or review]` | `[path]` | `[report section; source file and lines or symbol; document heading]` | `[decision, comparison, constraint, or challenge]` | `[report-only, product mismatch, scale mismatch, unverified interpretation, or none]` |

Name each study report by exact path and section. When repository behavior matters, add a separate row for the study repository file and line range or symbol. A study directory, report directory, or repository name alone is not evidence.

## Prior Sprint Decisions Applied

| Sprint 1 API or authority decision | Realized caller use | Status | Sprint 2 pressure | Decision needed |
| --- | --- | --- | --- | --- |
| `[decision]` | `[use]` | `[preserved, extended, superseded, unaffected]` | `[pressure]` | `[decision]` |

## Actors and authority

| Actor | May inspect | May wait | May observe | May request cancellation | May commit transition | May publish event/outcome |
| --- | --- | --- | --- | --- | --- | --- |
| Run caller | `[yes or no]` | `[yes or no]` | `[yes or no]` | `[yes or no]` | no | no |
| Read-only observer | `[yes or no]` | `[yes or no]` | `[yes or no]` | no | no | no |
| Work function | `[yes or no]` | `[yes or no]` | `[yes or no]` | `[yes or no]` | no | no |
| Aren runtime owner | yes | yes | yes | accepts | yes | yes |

If the chosen Go values do not enforce this distinction structurally, explain the package privacy and construction rules that do.

## Caller workflows

Write representative caller code for:

1. start, inspect, wait;
2. start, follow from sequence zero, receive terminal event;
3. start, request cancellation, inspect disposition, wait for truthful outcome;
4. observe a run without cancellation authority;
5. abandon waiting or observation without affecting the run.

Use these workflows to expose awkward method placement, ambiguous contexts, or authority leakage.

## Candidate API models

| Option | Returned start value | Observation capability | Cancellation capability | Construction control | Benefits | Misuse risks |
| --- | --- | --- | --- | --- | --- | --- |
| One run handle with methods | `[shape]` | `[shape]` | `[shape]` | `[mechanism]` | `[benefits]` | `[risks]` |
| Separate view and controller values | `[shape]` | `[shape]` | `[shape]` | `[mechanism]` | `[benefits]` | `[risks]` |
| Narrow interfaces over private implementation | `[shape]` | `[shape]` | `[shape]` | `[mechanism]` | `[benefits]` | `[risks]` |
| Other credible model | `[shape]` | `[shape]` | `[shape]` | `[mechanism]` | `[benefits]` | `[risks]` |

Assess whether interfaces are earned by distinct capabilities or merely mirror one implementation.

## Method contract table

| Operation | Receiver capability | Inputs | Blocking | Result | Expected errors or dispositions | State mutation |
| --- | --- | --- | --- | --- | --- | --- |
| Start | `[receiver]` | `[inputs]` | `[behavior]` | `[result]` | `[failures]` | `[mutation]` |
| State/snapshot | `[receiver]` | `[inputs]` | no | `[result]` | `[failures]` | none |
| Wait | `[receiver]` | `[inputs]` | `[behavior]` | `[result]` | `[failures]` | none |
| Observe/replay | `[receiver]` | `[inputs]` | `[behavior]` | `[result]` | `[failures]` | observer-local only |
| Cancel | `[receiver]` | `[inputs]` | `[behavior]` | `[disposition]` | `[failures]` | cancellation facts only |

## Context semantics

For each context accepted by the API, state what its cancellation means. Distinguish parent run lifetime, cancellation request cause, waiting lifetime, and observer lifetime. Do not reuse one context for unrelated authority.

## Construction and invariant protection

Analyze zero values, caller construction of handles, forged run IDs, copied handles, concurrent method calls, exported fields, and returned pointers. State which types are safe to copy and which are intentionally opaque.

## Compatibility posture

Phase 1 does not promise a stable public API, but careless exposure still creates migration cost. Identify which semantics should remain stable and which names or shapes remain provisional. Explain how examples and diagnostic CLI depend on the API without freezing future executor design.

## Area Decisions

| Decision | Caller-visible contract | Authority protected | Evidence and rationale | Rejected alternative |
| --- | --- | --- | --- | --- |
| `[decision]` | `[contract]` | `[authority]` | `[basis]` | `[alternative]` |

## Trade-Offs

Analyze capability clarity versus API volume, interfaces versus concrete handles, context-aware blocking versus method complexity, zero-value usability versus invariant construction, and near-term convenience versus future compatibility.

## Evidence

Use `Inputs Used` IDs in every material claim. Cite selected package or API study reports by exact path and section and relevant study repository APIs by file and line range or symbol. Distinguish report interpretation from direct repository observation. Use the PRD authority boundary, project-wide reasoning, Sprint 1 API experience, and code-context callers. Do not copy HTTP resource or distributed handle designs into an in-process library.

## Risks

Include caller mutation of state, a cancellation method that implies synchronous stop, read-only users receiving control methods, wait-context cancellation changing run state, interfaces introduced for imagined SDKs, and public fields that permit invalid outcomes or events.

## Verification obligations

| API claim | Compile-time check, test, or example | Misuse attempted | Expected result |
| --- | --- | --- | --- |
| `[claim]` | `[evidence]` | `[misuse]` | `[result]` |

## Self-critique

- Can a caller obtain internal transition authority through any exported value?
- Does `Cancel` read as a request or a confirmed stop?
- Which interface has only one consumer and one implementation?
- What public shape is most likely to be regretted in Phase 2 or Phase 3?
