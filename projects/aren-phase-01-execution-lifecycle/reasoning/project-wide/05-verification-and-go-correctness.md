# Project reasoning template: Verification and Go correctness

## Purpose

Define how Phase 1 claims will be falsified and proved in Go. Convert study findings and prototype defects into permanent invariants, controlled schedules, negative controls, and review evidence.

Read `projects/aren-phase-01-execution-lifecycle/reasoning/README.md` and follow its required reasoning standard.

## Inputs Used

Complete this table before analysis. Record only material actually read.

| Input ID | Input | Kind | Exact workspace-relative path | Exact portion inspected | How it was used | Limits or applicability warning |
| --- | --- | --- | --- | --- | --- | --- |
| `[IN-NN]` | `[name]` | `[requirement, project reasoning, study report, study repository, prototype report, current code, prior review]` | `[path]` | `[report section; source file and lines or symbol; document heading]` | `[question, comparison, constraint, or challenge]` | `[report-only, product mismatch, scale mismatch, unverified interpretation, or none]` |

Name each study report by exact path and section. When repository behavior matters, add a separate row for the study repository file and line range or symbol. A study directory, report directory, or repository name alone is not evidence.

## Governing questions

- Which claims are pure decision functions and which require concurrent execution?
- Which failures can the race detector find, and which logical races require assertions?
- How will tests control schedule without depending on sleeps?
- How will leak and quiescence claims be measured?
- How will tests prove that their own instruments can fail?
- Which prototype defects become named regression tests?

## Normative constraints

| Required proof | Source | Why ordinary unit tests are insufficient | Expected evidence class |
| --- | --- | --- | --- |
| `[proof]` | `[path and section]` | `[reason]` | `[model, unit, race, stress, leak, CLI, review]` |

## Evidence

Use `Inputs Used` IDs in every material claim. Cite the exact report section and, where relevant, the study repository file and line range or symbol. Distinguish report interpretation from direct repository observation.

### Verification pattern assessment

| Evidence | Technique | Failure class exposed | Determinism | Cost | Blind spot | Aren use |
| --- | --- | --- | --- | --- | --- | --- |
| `[source]` | `[race detector, property test, model, fault seam, virtual time, census, regression]` | `[failure]` | `[level]` | `[cost]` | `[blind spot]` | `[adopt, adapt, reject]` |

### Historical defect register

| Defect | Original evidence | Broken invariant | Why earlier tests missed it | Permanent regression design | Negative control |
| --- | --- | --- | --- | --- | --- |
| `[defect]` | `[path]` | `[invariant]` | `[reason]` | `[test]` | `[deliberate mutation or fixture]` |

## Reference model

Define the smallest pure model needed to check lifecycle commands.

| Command or fact | Model precondition | Model update | Expected output | Illegal case |
| --- | --- | --- | --- | --- |
| Create | `[precondition]` | `[update]` | `[output]` | `[case]` |
| Start | `[precondition]` | `[update]` | `[output]` | `[case]` |
| Complete | `[precondition]` | `[update]` | `[output]` | `[case]` |
| Request cancellation | `[precondition]` | `[update]` | `[output]` | `[case]` |

State what the model deliberately omits. A reference model that repeats the implementation structure provides weak independent evidence.

## Invariant catalog

| Invariant ID | Statement | Scope | Observable oracle | Concurrent schedule needed? | Required suites |
| --- | --- | --- | --- | --- | --- |
| `[INV-ID]` | `[precise invariant]` | `[Sprint 1, Sprint 2, both]` | `[state, event, outcome, goroutine, CLI]` | `[yes or no]` | `[suites]` |

Include at least one-terminal commitment, state/outcome agreement, terminal event agreement, monotonic sequence, waiter equality, cancellation truthfulness, immutable publication, observer independence, and Aren-owned quiescence.

## Schedule-control design

List every required race and the mechanism that opens each window.

| Race | Controlled checkpoints | Forbidden use of sleep | Number of repetitions | Assertions during race | Assertions after quiescence |
| --- | --- | --- | --- | --- | --- |
| `[race]` | `[barrier, channel, hook, fake work]` | `[how timing is avoided]` | `[count and basis]` | `[assertions]` | `[assertions]` |

Test hooks must not become public runtime configuration or alter production semantics.

## Test portfolio

| Test level | Behaviors proved | Real dependencies | Faults injected | Command | Release gate |
| --- | --- | --- | --- | --- | --- |
| Pure/model | `[behaviors]` | `[none]` | `[commands]` | `[command]` | `[gate]` |
| Package contract | `[behaviors]` | `[runtime package]` | `[faults]` | `[command]` | `[gate]` |
| Race/stress | `[behaviors]` | `[real goroutines]` | `[schedules]` | `[command]` | `[gate]` |
| Leak/quiescence | `[behaviors]` | `[real goroutines]` | `[abandonment]` | `[command]` | `[gate]` |
| Diagnostic CLI | `[behaviors]` | `[real Aren runtime]` | `[scenarios]` | `[command]` | `[gate]` |

## Falsification plan

For every high-value test instrument, state how a deliberately broken implementation, seeded historical defect, or mutation should make it fail. Avoid tests that assert only eventual success or count unrelated activity.

## Flake and suite economics

Separate deterministic required tests from repeated stress and longer soak runs. Define seeds, reproduction output, time bounds, platform expectations, and what a timeout means.

## Project conclusions

| Conclusion | Evidence basis | Required in Sprint 1 | Added in Sprint 2 | Reopen trigger |
| --- | --- | --- | --- | --- |
| `[conclusion]` | `[evidence]` | `[proof]` | `[proof]` | `[trigger]` |

## Trade-Offs

Address deterministic seams versus production-code intrusion, high iteration counts versus controlled schedules, goroutine census versus precise ownership assertions, fake clocks versus time acceleration, and exhaustive matrices versus maintainable suites.

## Risks

Include false-positive tests, race-detector overconfidence, goroutine-count brittleness, sleeps that hide scheduling defects, stress without reproducible seeds, tests coupled to internal locks, and expensive suites quietly removed from CI.

## Self-critique

- Which claimed invariant has no independent oracle?
- Which race test only makes the failure less likely rather than controlling it?
- What defect would pass every proposed test?
- Can a clean test process still hide a goroutine retained until process exit?
