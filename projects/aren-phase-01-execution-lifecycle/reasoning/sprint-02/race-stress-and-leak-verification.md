# Area reasoning template: Race, stress, and leak verification

> Selected template: `projects/aren-phase-01-execution-lifecycle/reasoning/sprint-02/race-stress-and-leak-verification.md`

## Purpose

Define the adversarial proof required to accept Sprint 2 and close Phase 1. This document owns deterministic race schedules, repeated stress, observer abandonment, leak checks, negative controls, CLI evidence, and exact release commands.

## Inputs Used

Complete this table before analysis. Record only material actually read.

| Input ID | Input | Kind | Exact workspace-relative path | Exact portion inspected | How it was used | Limits or applicability warning |
| --- | --- | --- | --- | --- | --- | --- |
| `[IN-NN]` | `[name]` | `[requirement, project reasoning, handbook, study report, study repository, prototype report, current code, Sprint 1 reasoning or review]` | `[path]` | `[report section; source file and lines or symbol; document heading]` | `[claim, test design, constraint, or challenge]` | `[report-only, product mismatch, scale mismatch, unverified interpretation, or none]` |

Name each study report by exact path and section. When repository behavior matters, add a separate row for the study repository file and line range or symbol. A study directory, report directory, or repository name alone is not evidence.

## Prior Sprint Evidence Applied

| Sprint 1 invariant or test | Existing evidence | Preserve, extend, or replace | New failure pressure | Reason |
| --- | --- | --- | --- | --- |
| `[invariant or test]` | `[path or command]` | `[status]` | `[pressure]` | `[reason]` |

Do not rewrite sound Sprint 1 tests merely to move them into a larger suite.

## Phase exit claims

| Claim | Invariant ID | Required oracle | Evidence level | Why weaker evidence is insufficient |
| --- | --- | --- | --- | --- |
| Exactly one terminal commitment under races | `[ID]` | `[oracle]` | `[model, deterministic race, stress]` | `[reason]` |
| Truthful cancellation | `[ID]` | `[oracle]` | `[level]` | `[reason]` |
| Observer independence | `[ID]` | `[oracle]` | `[level]` | `[reason]` |
| Aren-owned quiescence | `[ID]` | `[oracle]` | `[level]` | `[reason]` |

## Adversarial schedule matrix

| Schedule ID | Actors | Controlled checkpoints | Release orders covered | Assertions inside window | Final assertions |
| --- | --- | --- | --- | --- | --- |
| `[SCH-NN]` | `[actors]` | `[barriers or hooks]` | `[orders]` | `[assertions]` | `[assertions]` |

Cover explicit cancel versus completion, parent cancel versus completion, two cancel sources, observer registration versus terminal commit, replay versus new append, abandoned observer versus terminal, many waiters, and concurrent state/history/outcome reads.

## Cancellation behavior matrix

| Work behavior | Request timing | Work return | Expected disposition | Expected terminal state | Expected cancellation event count |
| --- | --- | --- | --- | --- | --- |
| Immediate cooperation | `[timing]` | `[error]` | `[disposition]` | `[state]` | `[count]` |
| Delayed acknowledgement | `[timing]` | `[error]` | `[disposition]` | `[state]` | `[count]` |
| Ignores cancellation then succeeds | `[timing]` | nil | `[disposition]` | `succeeded` | `[count]` |
| Accepted cancellation then unrelated error | `[timing]` | `[error]` | `[disposition]` | `failed` | `[count]` |
| Panic after request | `[timing]` | panic | `[disposition]` | `failed` | `[count]` |

## Observer matrix

| Observer schedule | Count | Cursor | Consumer speed | Abandonment | Expected events | Expected completion |
| --- | --- | --- | --- | --- | --- | --- |
| `[schedule]` | `[count]` | `[cursor]` | `[speed]` | `[behavior]` | `[events]` | `[completion]` |

## Leak and quiescence model

Define Aren-owned resources precisely. Exclude unrelated runtime goroutines and document the baseline method.

| Resource | Owner | Created by | Expected release point | Detection method | Timeout meaning |
| --- | --- | --- | --- | --- | --- |
| `[goroutine, timer, subscription, channel state]` | `[owner]` | `[path]` | `[point]` | `[method]` | `[meaning]` |

Prefer owner-specific completion signals and behavioral quiescence over brittle global goroutine counts. If a census tool is used, explain filtering and stabilization.

## Stress design

| Stress target | Iterations or duration | Parallelism | Seed handling | Per-iteration oracle | Failure artifact |
| --- | --- | --- | --- | --- | --- |
| `[target]` | `[bound]` | `[bound]` | `[method]` | `[oracle]` | `[artifact]` |

Explain what stress adds beyond deterministic schedules. High iteration count is not a substitute for opening the intended race window.

## Negative controls and mutation proof

| Instrument | Seeded defect | Expected failing assertion | Historical defect represented | Required frequency |
| --- | --- | --- | --- | --- |
| `[instrument]` | `[defect]` | `[assertion]` | `[defect]` | `[per change, CI, periodic]` |

At minimum, prove detection of duplicate terminal publication, parent cancellation bypass, terminal close before event append, blocked observer controlling execution, leaked watcher, and false cancellation classification.

## Suite tiers and commands

| Tier | Exact command | Deterministic? | Expected duration | Required in normal CI | Failure handling |
| --- | --- | --- | --- | --- | --- |
| Fast contract | `[command]` | `[yes or no]` | `[budget]` | `[yes or no]` | `[handling]` |
| Full race | `go test -race ./...` or selected equivalent | `[yes or no]` | `[budget]` | yes | `[handling]` |
| Repeated targeted stress | `[command]` | `[yes or seeded]` | `[budget]` | `[gate]` | `[handling]` |
| Diagnostic CLI | `[commands]` | `[yes or no]` | `[budget]` | `[gate]` | `[handling]` |

## Area Decisions

| Decision | Proof mechanism | Defect class | Cost accepted | Residual blind spot |
| --- | --- | --- | --- | --- |
| `[decision]` | `[mechanism]` | `[defect]` | `[cost]` | `[blind spot]` |

## Trade-Offs

Analyze deterministic hooks versus production intrusion, suite runtime versus release confidence, leak tools versus owner-specific proof, broad `-race` runs versus targeted collision tests, and model-based tests versus readable regression cases.

## Evidence

Use `Inputs Used` IDs in every material claim. Cite selected study and prototype reports by exact path and section. Cite relevant study repository test files, helpers, and symbols when a technique or historical defect depends on their behavior. Distinguish report interpretation from direct repository observation. Trace tests to project-wide verification conclusions, Sprint 1 reviews, and language-decision rules. State which claims remain beyond what `go test -race` can prove.

## Risks

Include false-positive completion, tests that never enter the race, flakes hidden by retries, global goroutine noise, stress without reproduction data, provider or network dependencies in normal tests, and a passing CLI demonstration mistaken for invariant proof.

## Phase closure evidence

| Exit criterion | Required artifacts and commands | Pass condition | Blocking ambiguity |
| --- | --- | --- | --- |
| `[criterion]` | `[evidence]` | `[condition]` | `[ambiguity]` |

## Self-critique

- Which phase-exit claim still rests on sampling rather than an invariant oracle?
- Which leak can survive until process exit without failing the proposed test?
- Which negative control is too coupled to the current implementation?
- What would make a green test run untrustworthy?
