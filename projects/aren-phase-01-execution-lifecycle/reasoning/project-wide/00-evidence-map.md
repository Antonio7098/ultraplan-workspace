# Project reasoning template: Phase 1 evidence assessment and routing

## Purpose

Assess the evidence selected for Aren Phase 1 before thematic reasoning begins. `project-index.md` catalogs what is available, and the project-reasoning index selects what enters this analysis. This document evaluates the selected evidence, records canonical or superseded relationships, tests applicability, exposes conflicts and gaps, and routes findings to decision areas. It does not duplicate either index and does not choose the runtime architecture.

Read `projects/aren-phase-01-execution-lifecycle/reasoning/README.md` and follow its required reasoning standard.

## Inputs Used

Complete this table before analysis. Record only material that was actually read. Do not copy every entry from `project-index.md`.

| Input ID | Input | Kind | Exact workspace-relative path | Exact portion inspected | How it was used | Limits or applicability warning |
| --- | --- | --- | --- | --- | --- | --- |
| `[IN-NN]` | `[name]` | `[requirement, project document, study report, study repository, prototype report, current code, prior review]` | `[path]` | `[report section; source file and lines or symbol; document heading]` | `[question, comparison, constraint, or challenge]` | `[report-only, product mismatch, scale mismatch, unverified interpretation, or none]` |

For a study report, name its exact report path and section. When repository behavior matters, add a separate study-repository row with the concrete file and line range or symbol. A study directory, dimension directory, report directory, or repository name alone does not count as an input reference.

## Scope

- Phase: Aren Phase 1, execution lifecycle.
- Sprints served: Core Lifecycle; Cancellation and Concurrency.
- Required inputs: Phase 1 PRD, project roadmap, Aren phased roadmap, project lineage, final language decision, selected study reports, prototype reports, and any current Aren code named by the project-reasoning manifest.
- Excluded subjects: provider integration, tools, subprocesses, persistence, workflows, daemon hosting, remote execution, and later-phase designs.

## Selection and canonicality assessment

Assess every selected evidence input. Use the IDs from `Inputs Used`. Do not repeat paths or other catalog metadata here, and do not equate a dimension brief with its final report.

| Input ID | Status | Canonical replacement | Phase 1 relevance | Applicability boundary | Reason |
| --- | --- | --- | --- | --- | --- |
| `[IN-NN]` | `[canonical, partial, duplicate, superseded, unusable]` | `[input ID or none]` | `[high, medium, low, none]` | `[where its assumptions stop matching Aren]` | `[concrete evidence for the assessment]` |

Status decisions must cite the concrete defect or relationship. Examples include duplicate filenames, a synthesis built from no applicable sources, a report superseded by a corrected version, or conclusions that apply only to durable workflows.

## Evidence quality assessment

For each canonical input, assess what it can prove.

| Evidence ID | Direct observation | Author interpretation | Missing proof | Source assumptions | Confidence |
| --- | --- | --- | --- | --- | --- |
| `[ID]` | `[code, test, experiment, or contract observed]` | `[interpretation made by report]` | `[what was not tested or inspected]` | `[scale, durability, deployment, language, product]` | `[high, medium, low]` |

Do not rate confidence from prose quality or repository popularity. Prefer direct code paths, adversarial tests, reproduced prototype failures, and matching product constraints.

## Phase 1 decision map

Route evidence to Aren decision families. The mapping is many-to-many.

| Decision family | Required questions | Primary evidence | Contradicting evidence | Missing evidence | Intended downstream document |
| --- | --- | --- | --- | --- | --- |
| Lifecycle authority and atomic publication | `[questions]` | `[IDs]` | `[IDs and contradiction]` | `[gap]` | `01-lifecycle-authority-and-atomic-publication.md` |
| Outcomes, failures, and terminal resolution | `[questions]` | `[IDs]` | `[IDs and contradiction]` | `[gap]` | `02-outcomes-failures-and-terminal-resolution.md` |
| Cancellation, ownership, and cleanup | `[questions]` | `[IDs]` | `[IDs and contradiction]` | `[gap]` | `03-cancellation-goroutine-ownership-and-cleanup.md` |
| Events, observation, waiting, and replay | `[questions]` | `[IDs]` | `[IDs and contradiction]` | `[gap]` | `04-events-observation-waiting-and-replay.md` |
| Verification and Go correctness | `[questions]` | `[IDs]` | `[IDs and contradiction]` | `[gap]` | `05-verification-and-go-correctness.md` |

## Requirement coverage

Map every material Phase 1 requirement to evidence or identify that it rests only on a product decision.

| Requirement or invariant | Evidence supporting it | Evidence challenging it | Coverage | Required follow-up |
| --- | --- | --- | --- | --- |
| `[PRD section or invariant]` | `[IDs]` | `[IDs or none]` | `[strong, partial, product decision only, missing]` | `[reasoning, experiment, test, or none]` |

Missing external precedent does not invalidate a product requirement. It changes the verification burden.

## Contradiction register

Record real disagreements. Do not flatten them into a vague combined recommendation.

| Conflict | Evidence on each side | Why sources differ | Aren consequence | Resolution owner |
| --- | --- | --- | --- | --- |
| `[question]` | `[IDs and positions]` | `[product or scale assumptions]` | `[what Phase 1 must decide or test]` | `[project document or sprint area]` |

## Negative transfer register

Identify attractive patterns that Phase 1 must not copy.

| Pattern | Source | Why it works there | Why it is premature or wrong for Phase 1 | Revisit trigger |
| --- | --- | --- | --- | --- |
| `[pattern]` | `[ID]` | `[mechanism and assumption]` | `[specific scope or semantic conflict]` | `[measured future pressure]` |

At minimum, inspect durable workflow engines, distributed schedulers, brokers, databases, general executor interfaces, and plugin mechanisms for phase leakage.

## Evidence gaps and focused work

| Gap | Why existing evidence is insufficient | Smallest way to answer it | Owner | Blocking? |
| --- | --- | --- | --- | --- |
| `[gap]` | `[reason]` | `[prototype, code inspection, state model, test, or explicit deferral]` | `[project or sprint]` | `[yes or no]` |

Do not answer a runtime semantic question with more comparative reading when a small executable collision test would provide stronger evidence.

## Project conclusions

List only conclusions about evidence use, not implementation choices.

| Conclusion | Basis | Authority | Applies to | Reopen when |
| --- | --- | --- | --- | --- |
| `[which evidence is canonical or how it must be interpreted]` | `[IDs]` | `[evidence selection or project constraint]` | `[document or sprint]` | `[condition]` |

## Trade-Offs

Analyze the cost of the selected evidence boundary. Include reports excluded for weak applicability and knowledge that may be lost by selecting only final syntheses.

## Evidence

Use `Inputs Used` IDs in every material claim. Provide the exact report section or repository file and line range or symbol for every status, contradiction, and exclusion decision. Distinguish what a report concludes from what its study repository demonstrates. Record `report-only` when no underlying repository source was inspected.

## Risks

Cover at least stale report revisions, duplicate reports, synthesis drift, source licensing limits, over-weighting mature infrastructure, under-weighting Aren prototypes, and accidental later-phase design.

## Self-critique

- Which exclusion could remove evidence that changes a Phase 1 decision?
- Which high-confidence report rests on assumptions least like Aren?
- Where did the map rely on a report summary without checking its evidence index?
- What is the clearest sign that this evidence process has become larger than the decision it supports?
