# Area reasoning template: Outcome and failure contract

> Selected template: `projects/aren-phase-01-execution-lifecycle/reasoning/sprint-01/outcome-and-failure-contract.md`

## Purpose

Decide Sprint 1 result, outcome, timing, returned-error, panic, and Aren-invariant contracts. This document owns the semantic type model and construction rules. It does not own transition synchronization or Sprint 2 cancellation interpretation.

## Inputs Used

Complete this table before analysis. Record only material actually read.

| Input ID | Input | Kind | Exact workspace-relative path | Exact portion inspected | How it was used | Limits or applicability warning |
| --- | --- | --- | --- | --- | --- | --- |
| `[IN-NN]` | `[name]` | `[requirement, project reasoning, handbook, study report, study repository, prototype report, current code, prior review]` | `[path]` | `[report section; source file and lines or symbol; document heading]` | `[decision, comparison, constraint, or challenge]` | `[report-only, product mismatch, scale mismatch, unverified interpretation, or none]` |

Name each study report by exact path and section. When repository behavior matters, add a separate row for the study repository file and line range or symbol. A study directory, report directory, or repository name alone is not evidence.

## Decision scope

- Success with exact result.
- Returned work error, including result plus error.
- Work panic with recognizable panic classification.
- Aren-origin invariant failure.
- Start and finish timing.
- Immutable publication under Go's mutable type system.

## Inherited constraints

| Input | Constraint or finding | Consequence | Uncertainty remaining |
| --- | --- | --- | --- |
| `[input]` | `[finding]` | `[consequence]` | `[uncertainty]` |

## Caller use cases

List how callers must inspect terminal state, retrieve a result, distinguish failures, unwrap causes, identify panic, and compare outcomes returned to multiple waiters. Avoid adding methods without a concrete caller need.

## Candidate type models

| Option | Success representation | Failure representation | Invalid combinations prevented by | Cause access | Generic result handling | Cost |
| --- | --- | --- | --- | --- | --- | --- |
| Validated outcome value | `[shape]` | `[shape]` | `[mechanism]` | `[mechanism]` | `[mechanism]` | `[cost]` |
| Separate terminal variants | `[shape]` | `[shape]` | `[mechanism]` | `[mechanism]` | `[mechanism]` | `[cost]` |
| Other credible Go model | `[shape]` | `[shape]` | `[mechanism]` | `[mechanism]` | `[mechanism]` | `[cost]` |

Analyze compile-time guarantees honestly. Named types and private fields do not make nested caller-owned data deeply immutable.

## Validity matrix

| Terminal state | Result present | Failure present | Cancellation present | Start time | Finish time | Valid? | Construction path |
| --- | --- | --- | --- | --- | --- | --- | --- |
| `[state]` | `[yes or no]` | `[yes or no]` | `[yes or no]` | `[state]` | `[state]` | `[yes or no]` | `[constructor or impossible]` |

Include nonterminal and contradictory combinations.

## Failure taxonomy

| Origin | Kind | Trigger | Cause preserved? | Human message source | Caller action | Sprint introduced |
| --- | --- | --- | --- | --- | --- | --- |
| `work` | `returned_error` | `[trigger]` | `[mechanism]` | `[source]` | `[action]` | Sprint 1 |
| `work` | `panic` | `[trigger]` | `[mechanism]` | `[source]` | `[action]` | Sprint 1 |
| `aren` | `invariant_violation` | `[trigger]` | `[mechanism]` | `[source]` | `[action]` | Sprint 1 |

Reject additional kinds unless a Phase 1 caller changes behavior because of them.

## Panic boundary

Decide:

- where recovery occurs;
- which panic value and stack information are retained;
- whether panic data can contain sensitive values;
- whether runtime panics and Aren invariant panics share a representation;
- how terminal publication proceeds if panic normalization itself fails;
- what diagnostic CLI output may expose.

## Result ownership and immutability

| Value | Created by | Copied before publication? | Copy depth | Returned to callers as | Mutation risk accepted |
| --- | --- | --- | --- | --- | --- |
| Work result | `[owner]` | `[yes or no]` | `[depth]` | `[value, pointer, interface]` | `[risk]` |
| Failure | `[owner]` | `[yes or no]` | `[depth]` | `[form]` | `[risk]` |
| Outcome | `[owner]` | `[yes or no]` | `[depth]` | `[form]` | `[risk]` |

State the precise guarantee if Aren cannot deep-copy arbitrary result values.

## Timing semantics

Decide clock ownership, capture points, zero-value behavior, monotonic duration use, equality across waiters, and whether a clock abstraction is earned. Do not introduce injectable time without identifying a test that cannot be deterministic otherwise.

## Area Decisions

| Decision | Exact contract | Evidence and rationale | Rejected alternative | Verification |
| --- | --- | --- | --- | --- |
| `[decision]` | `[contract]` | `[basis]` | `[alternative]` | `[test]` |

## Trade-Offs

Cover type safety versus API weight, cause preservation versus stable failure vocabulary, stack retention versus safety, generic results versus copying guarantees, and test clock seams versus direct `time.Now`.

## Evidence

Use `Inputs Used` IDs in every material claim. Cite each study or prototype report by exact path and section and each relevant study repository mechanism by file and line range or symbol. Distinguish report interpretation from direct repository observation. Trace each failure kind and validity rule to a requirement, project-wide conclusion, prototype defect, or current caller need. Record absent evidence instead of filling gaps with conventions.

## Risks

Analyze mutable results shared across waiters, error causes collapsed to strings, panic mistaken for returned error, invariant failure overwriting a work failure, zero times interpreted as real values, and a Phase 1 taxonomy shaped by future provider errors.

## Verification obligations

| Contract | Positive case | Negative case | Aliasing or mutation check | Failure detected |
| --- | --- | --- | --- | --- |
| `[contract]` | `[test]` | `[test]` | `[test]` | `[defect]` |

## Self-critique

- Which invalid combination can package-internal code still create?
- Can two waiters mutate data visible to each other?
- Which panic information is useful enough to retain and safe enough to expose?
- Did the design promise more immutability than Go can enforce for arbitrary results?
