# Project reasoning template: Outcomes, failures, and terminal resolution

## Purpose

Synthesize the evidence that governs successful results, returned failures, panics, Aren invariant failures, cancellation facts, and exactly one terminal outcome. Establish the semantic vocabulary that both sprints must preserve.

Read `projects/aren-phase-01-execution-lifecycle/reasoning/README.md` and follow its required reasoning standard.

## Inputs Used

Complete this table before analysis. Record only material actually read.

| Input ID | Input | Kind | Exact workspace-relative path | Exact portion inspected | How it was used | Limits or applicability warning |
| --- | --- | --- | --- | --- | --- | --- |
| `[IN-NN]` | `[name]` | `[requirement, project reasoning, study report, study repository, prototype report, current code, prior review]` | `[path]` | `[report section; source file and lines or symbol; document heading]` | `[question, comparison, constraint, or challenge]` | `[report-only, product mismatch, scale mismatch, unverified interpretation, or none]` |

Name each study report by exact path and section. When repository behavior matters, add a separate row for the study repository file and line range or symbol. A study directory, report directory, or repository name alone is not evidence.

## Governing questions

- Which facts describe work completion before Aren interprets them?
- What makes a run successful, failed, or cancelled?
- How do result, failure, and cancellation remain mutually coherent?
- How does panic remain recognizable without escaping the run boundary?
- Which failures originate in work and which originate in Aren?
- Does scheduler arrival order decide the outcome, or does one policy interpret committed facts?

## Normative constraints

| Constraint | Source | Fixed semantic rule | Open implementation question |
| --- | --- | --- | --- |
| `[constraint]` | `[path and section]` | `[rule]` | `[question]` |

## Evidence

Use `Inputs Used` IDs in every material claim. Cite the exact report section and, where relevant, the study repository file and line range or symbol. Distinguish report interpretation from direct repository observation.

### Completion model comparison

| Evidence | Completion vocabulary | Finalization owner | Panic treatment | Cancellation treatment | Result validation | Aren relevance |
| --- | --- | --- | --- | --- | --- | --- |
| `[source]` | `[states or signals]` | `[owner]` | `[mechanism]` | `[mechanism]` | `[mechanism]` | `[finding]` |

Separate framework completion, provider completion, agent task completion, and lifecycle completion. Phase 1 should adopt only the last of these.

### Failure classification evidence

| Failure case | Source treatment | Information preserved | Information lost | Caller behavior enabled | Phase 1 conclusion pressure |
| --- | --- | --- | --- | --- | --- |
| Returned error | `[evidence]` | `[facts]` | `[facts]` | `[behavior]` | `[pressure]` |
| Panic | `[evidence]` | `[facts]` | `[facts]` | `[behavior]` | `[pressure]` |
| Invariant violation | `[evidence]` | `[facts]` | `[facts]` | `[behavior]` | `[pressure]` |
| Cancellation-related error | `[evidence]` | `[facts]` | `[facts]` | `[behavior]` | `[pressure]` |
| Result plus error | `[evidence]` | `[facts]` | `[facts]` | `[behavior]` | `[pressure]` |

## Candidate outcome models

Compare credible representations, including a discriminated terminal outcome, state plus optional fields with validated construction, and separate success/failure types behind a common accessor.

| Model | Invalid combinations prevented by | Cause preservation | Immutability mechanism | Caller cost | Go-specific risk | Phase fit |
| --- | --- | --- | --- | --- | --- | --- |
| `[model]` | `[mechanism]` | `[mechanism]` | `[mechanism]` | `[cost]` | `[risk]` | `[fit]` |

Do not treat a Rust enum as directly available in Go. Explain the constructors, privacy, copying, and validation needed to provide an equivalent runtime contract.

## Terminal fact model

List raw facts before interpreting them.

| Fact | Producer | Recorded when | May race with | Mutable after capture? |
| --- | --- | --- | --- | --- |
| Work result | `[producer]` | `[time]` | `[facts]` | `[yes or no]` |
| Work error | `[producer]` | `[time]` | `[facts]` | `[yes or no]` |
| Panic | `[producer]` | `[time]` | `[facts]` | `[yes or no]` |
| Accepted cancellation and cause | `[producer]` | `[time]` | `[facts]` | `[yes or no]` |
| Existing terminal commitment | `[producer]` | `[time]` | `[facts]` | `[yes or no]` |

## Resolution matrix

Build a conceptual truth table. Do not leave precedence in prose.

| Result | Error | Panic | Cancellation accepted | Error matches run cancellation cause | Existing terminal | Candidate outcome | Reason |
| --- | --- | --- | --- | --- | --- | --- | --- |
| `[fact]` | `[fact]` | `[fact]` | `[fact]` | `[fact]` | `[fact]` | `[state and failure kind]` | `[policy]` |

Include ambiguous and invalid fact combinations. State whether they produce a work failure, Aren invariant failure, or are impossible by construction.

## Information preservation

Determine which information must survive in local errors, immutable outcomes, events, and later diagnostic rendering. Explain where Go error wrapping ends and where explicit structured fields begin.

## Project conclusions

| Conclusion | Evidence basis | Applies in Sprint 1 | Extended in Sprint 2 | Reopen trigger |
| --- | --- | --- | --- | --- |
| `[conclusion]` | `[evidence]` | `[application]` | `[extension]` | `[trigger]` |

## Trade-Offs

Cover closed vocabulary versus forward evolution, defensive copies versus generic result ergonomics, panic recovery versus programmer-fault visibility, preserving causes versus stable serialization, and narrow failure kinds versus premature taxonomy growth.

## Risks

Include result/error incoherence, mutable nested results, panic disguised as an ordinary error, cancellation inferred from `context.Canceled` without accepted cause, error-chain dependence on unstable strings, and a rich future failure taxonomy entering Phase 1.

## Verification obligations

| Semantic claim | Minimal test | Adversarial or negative test | Failure the test must detect |
| --- | --- | --- | --- |
| `[claim]` | `[test]` | `[test]` | `[defect]` |

## Self-critique

- Which invalid outcome can still be constructed inside the package?
- Which failure fact would be lost if the process only retained the terminal event?
- Could unrelated work errors be misclassified as cancellation?
- Which conclusion assumes a future serialization boundary that Phase 1 does not have?
