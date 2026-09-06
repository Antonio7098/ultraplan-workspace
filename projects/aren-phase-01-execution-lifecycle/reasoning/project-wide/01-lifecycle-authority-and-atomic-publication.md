# Project reasoning template: Lifecycle authority and atomic publication

## Purpose

Synthesize Phase 1 evidence about who owns lifecycle truth, what Aren commits atomically, and when state becomes observable. Produce cross-sprint conclusions without choosing exact Go packages or synchronization code.

Read `projects/aren-phase-01-execution-lifecycle/reasoning/README.md` and follow its required reasoning standard.

## Inputs Used

Complete this table before analysis. Record only material actually read.

| Input ID | Input | Kind | Exact workspace-relative path | Exact portion inspected | How it was used | Limits or applicability warning |
| --- | --- | --- | --- | --- | --- | --- |
| `[IN-NN]` | `[name]` | `[requirement, project reasoning, study report, study repository, prototype report, current code, prior review]` | `[path]` | `[report section; source file and lines or symbol; document heading]` | `[question, comparison, constraint, or challenge]` | `[report-only, product mismatch, scale mismatch, unverified interpretation, or none]` |

Name each study report by exact path and section. When repository behavior matters, add a separate row for the study repository file and line range or symbol. A study directory, report directory, or repository name alone is not evidence.

## Governing questions

- What fact is Aren's atomic unit?
- Which component alone may validate and commit a transition?
- What does work report, and what must it never mutate?
- Which state, event, timing, outcome, and notification facts must appear together?
- What prevents competing terminal producers from publishing conflicting truth?
- Which behaviour must fail loudly as an invariant violation?

## Normative constraints

Extract the exact PRD and language-decision constraints. Classify each as fixed, open to implementation choice, or deferred.

| Constraint | Source | Classification | Consequence for reasoning |
| --- | --- | --- | --- |
| `[constraint]` | `[path and section]` | `[fixed, implementation choice, deferred]` | `[what it rules in or out]` |

## Evidence

Use `Inputs Used` IDs in every material claim. Cite the exact report section and, where relevant, the study repository file and line range or symbol. Distinguish report interpretation from direct repository observation.

### Evidence ledger

| Evidence | Observed mechanism | Guarantee actually demonstrated | Source assumption | Aren applicability | Confidence |
| --- | --- | --- | --- | --- | --- |
| `[report or prototype]` | `[single owner, CAS, channel owner, lock, transition table, checkpoint]` | `[specific guarantee]` | `[durable engine, in-process library, distributed system, agent app]` | `[adopt, adapt, reject, later]` | `[rating]` |

### Evidence conflicts

Compare at least:

- first-completion-wins with interpretation of a committed fact set;
- mutex-protected direct commit with a dedicated ownership loop;
- state as source of truth with event log or checkpoint as source of truth;
- panic on invariant violation with classified Aren-origin failure;
- notification inside the commit path with publication after commit.

For each conflict, explain whether the sources disagree or solve different problems.

## Candidate authority models

Analyze credible models without selecting exact code.

| Model | Transition owner | Commit boundary | Competing producer handling | Publication ordering | Benefits | Failure modes | Phase fit |
| --- | --- | --- | --- | --- | --- | --- | --- |
| Direct single-owner critical section | `[owner]` | `[facts committed]` | `[mechanism]` | `[ordering]` | `[benefits]` | `[failures]` | `[fit]` |
| Serialized ownership loop | `[owner]` | `[facts committed]` | `[mechanism]` | `[ordering]` | `[benefits]` | `[failures]` | `[fit]` |
| Other credible model | `[owner]` | `[facts committed]` | `[mechanism]` | `[ordering]` | `[benefits]` | `[failures]` | `[fit]` |

Do not include durable databases or distributed fencing unless the comparison explains why Phase 1 rejects them.

## Atomicity model

Define the minimum coherent fact set conceptually.

| Transition kind | Preconditions | Facts produced | Observers released | Forbidden partial views |
| --- | --- | --- | --- | --- |
| Creation | `[facts]` | `[facts]` | `[none or observers]` | `[views]` |
| Start | `[facts]` | `[facts]` | `[views]` | `[views]` |
| Terminal | `[facts]` | `[facts]` | `[waiters and observers]` | `[views]` |

State whether event occurrence time belongs inside the commitment and how sequence order relates to clock time.

## Counterexample analysis

Walk through schedules that could expose false truth:

1. terminal state visible before outcome;
2. terminal event visible before state;
3. waiter released before terminal history append;
4. two goroutines both construct and publish outcomes;
5. work mutates shared result after publication;
6. illegal transition is ignored and execution continues.

For each schedule, identify the violated invariant and the kind of design that prevents it.

## Project conclusions

| Conclusion | Evidence basis | Fixed across both sprints? | Sprint decision still required | Falsification or reopen trigger |
| --- | --- | --- | --- | --- |
| `[conclusion]` | `[evidence]` | `[yes or no]` | `[decision]` | `[trigger]` |

Conclusions should describe semantics and authority. Leave exact type names, package names, locks, channels, and API signatures to sprint reasoning.

## Trade-Offs

Address inspectability versus serialized indirection, fail-fast invariants versus returned failures, direct state access versus snapshots, and narrow Phase 1 truth versus future executor generality.

## Risks

Include logical races that the Go race detector cannot find, accidental multiple writers, mutation after publication, hidden notification order, false generalization from durable engines, and premature universal executor design.

## Downstream obligations

| Obligation | Sprint 1 owner | Sprint 2 extension | Required observable proof |
| --- | --- | --- | --- |
| `[semantic obligation]` | `[area document]` | `[area document or unaffected]` | `[test or runtime evidence]` |

## Self-critique

- What is the strongest credible argument for a different authority model?
- Which atomicity claim is currently prose rather than executable proof?
- Could the proposed conclusions permit an observer to see a combination the PRD forbids?
- Did future durability concerns distort the in-memory Phase 1 model?
