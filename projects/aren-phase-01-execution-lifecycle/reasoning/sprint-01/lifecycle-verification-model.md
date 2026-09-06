# Area reasoning template: Core lifecycle verification model

> Selected template: `projects/aren-phase-01-execution-lifecycle/reasoning/sprint-01/lifecycle-verification-model.md`

## Purpose

Turn Sprint 1 lifecycle decisions into an independent state model, invariant catalog, controlled tests, and failure-producing negative controls. This document owns proof design. It must challenge the architecture rather than restate its implementation.

## Inputs Used

Complete this table before analysis. Record only material actually read.

| Input ID | Input | Kind | Exact workspace-relative path | Exact portion inspected | How it was used | Limits or applicability warning |
| --- | --- | --- | --- | --- | --- | --- |
| `[IN-NN]` | `[name]` | `[requirement, project reasoning, handbook, study report, study repository, prototype report, current code, prior review]` | `[path]` | `[report section; source file and lines or symbol; document heading]` | `[claim, test design, constraint, or challenge]` | `[report-only, product mismatch, scale mismatch, unverified interpretation, or none]` |

Name each study report by exact path and section. When repository behavior matters, add a separate row for the study repository file and line range or symbol. A study directory, report directory, or repository name alone is not evidence.

## Decision scope

- success, returned failure, and panic;
- legal and illegal transitions;
- exactly one terminal commitment;
- timing and state/outcome agreement;
- multiple waiters;
- history ordering and immutability;
- concurrent inspection during publication;
- required `go test -race` evidence.

## Claims under test

| Decision source | Claim | Testable observation | Why the claim could fail |
| --- | --- | --- | --- |
| `[area or requirement]` | `[claim]` | `[observation]` | `[failure mechanism]` |

## Independent reference model

Define model state without copying the implementation fields or lock structure.

| Model command | Preconditions | Next model state | Output | Invalid-command result |
| --- | --- | --- | --- | --- |
| Create | `[conditions]` | `[state]` | `[output]` | `[result]` |
| Start | `[conditions]` | `[state]` | `[output]` | `[result]` |
| Complete success | `[conditions]` | `[state]` | `[output]` | `[result]` |
| Complete error | `[conditions]` | `[state]` | `[output]` | `[result]` |
| Complete panic | `[conditions]` | `[state]` | `[output]` | `[result]` |

Explain how generated or enumerated command traces compare implementation snapshots with model snapshots.

## Invariant catalog

| Invariant ID | Precise statement | Oracle | Checked continuously or finally | Required suite |
| --- | --- | --- | --- | --- |
| `[INV-S1-NN]` | `[statement]` | `[oracle]` | `[timing]` | `[suite]` |

Include transition legality, one terminal state, one terminal event, one outcome, state/outcome agreement, monotonic sequence, timing coherence, waiter equality, history immutability, and no partial terminal view.

## Behavior matrix

| Scenario | Work behavior | Expected state | Expected result | Expected failure | Expected events | Wait behavior |
| --- | --- | --- | --- | --- | --- | --- |
| Success | `[behavior]` | `succeeded` | `[result]` | none | `[sequence]` | `[behavior]` |
| Returned error | `[behavior]` | `failed` | none | `[origin and kind]` | `[sequence]` | `[behavior]` |
| Result plus error | `[behavior]` | `failed` | `[rule]` | `[origin and kind]` | `[sequence]` | `[behavior]` |
| Panic | `[behavior]` | `failed` | none | `[origin and kind]` | `[sequence]` | `[behavior]` |

## Controlled concurrency schedules

| Schedule | Checkpoints and release order | Readers involved | Assertions inside window | Final assertions |
| --- | --- | --- | --- | --- |
| Two terminal proposals | `[control]` | `[readers]` | `[assertions]` | `[assertions]` |
| Many waiters before completion | `[control]` | `[readers]` | `[assertions]` | `[assertions]` |
| State/history reads during terminal commit | `[control]` | `[readers]` | `[assertions]` | `[assertions]` |
| Late wait and history read | `[control]` | `[readers]` | `[assertions]` | `[assertions]` |

No required test may depend only on `time.Sleep` to enter a race window.

## Negative controls

| Test instrument | Deliberate defect or broken fixture | Expected failure | What a pass would reveal about the test |
| --- | --- | --- | --- |
| `[instrument]` | `[defect]` | `[failure]` | `[test weakness]` |

Include at least duplicate terminal publication, waiter release before outcome, mutable returned history, and an illegal transition silently ignored.

## Race, repetition, and time bounds

Define exact targeted commands, repetition counts, deterministic seeds, per-test timeouts, and evidence retained on failure. Explain which tests need `-race`, which need `-count`, and which are already exhaustive.

## Area Decisions

| Decision | Proof mechanism | Why this level is needed | Cheaper alternative rejected | Residual blind spot |
| --- | --- | --- | --- | --- |
| `[decision]` | `[mechanism]` | `[reason]` | `[alternative]` | `[blind spot]` |

## Trade-Offs

Analyze independent model cost versus confidence, deterministic hooks versus runtime intrusion, exhaustive matrices versus test maintenance, and exact goroutine assertions versus behavioral quiescence.

## Evidence

Use `Inputs Used` IDs in every material claim. Name the exact Go-runtime study or prototype report and section. Cite a study repository file and line range or symbol when a technique or defect depends on source behavior. Distinguish report interpretation from direct repository observation. Use the project verification synthesis, language-decision rules, and final area decisions. Link every permanent regression to the defect it would detect.

## Risks

Include tests mirroring implementation, impossible assertions, race tests that miss the collision, brittle exact timing, a fake that bypasses the real transition path, and passing `-race` treated as proof of semantic correctness.

## Review and command evidence

| Evidence | Exact command or artifact | Pass condition | Blocked condition |
| --- | --- | --- | --- |
| `[evidence]` | `[command]` | `[condition]` | `[condition]` |

## Self-critique

- Which invariant has the weakest oracle?
- What implementation defect could satisfy the reference model because the model copied it?
- Which concurrency test does not control the schedule tightly enough?
- If all tests pass, what important Sprint 1 claim remains unproved?
