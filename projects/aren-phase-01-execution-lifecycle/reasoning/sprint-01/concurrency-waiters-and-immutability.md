# Area reasoning template: Concurrency, waiters, and immutability

> Selected template: `projects/aren-phase-01-execution-lifecycle/reasoning/sprint-01/concurrency-waiters-and-immutability.md`

## Purpose

Decide Sprint 1 synchronization, goroutine ownership, multiple-waiter behavior, snapshot reads, and publication immutability. This document owns concurrent access mechanics. It does not own the state model itself, Sprint 2 cancellation races, or multiple live observer delivery.

## Inputs Used

Complete this table before analysis. Record only material actually read.

| Input ID | Input | Kind | Exact workspace-relative path | Exact portion inspected | How it was used | Limits or applicability warning |
| --- | --- | --- | --- | --- | --- | --- |
| `[IN-NN]` | `[name]` | `[requirement, project reasoning, handbook, study report, study repository, prototype report, current code, prior review]` | `[path]` | `[report section; source file and lines or symbol; document heading]` | `[decision, comparison, constraint, or challenge]` | `[report-only, product mismatch, scale mismatch, unverified interpretation, or none]` |

Name each study report by exact path and section. When repository behavior matters, add a separate row for the study repository file and line range or symbol. A study directory, report directory, or repository name alone is not evidence.

## Decision scope

- one writer for lifecycle mutation;
- concurrent state, outcome, and history reads;
- work invocation outside lifecycle mutation;
- multiple waiters before and after completion;
- no waiter consumes the outcome;
- Aren-owned goroutine accounting;
- defensive publication under Go.

## Current concurrency facts

| Fact | Code or prototype reference | Proven behavior | Unproven assumption |
| --- | --- | --- | --- |
| `[fact]` | `[path and lines]` | `[behavior]` | `[assumption]` |

## Shared-state inventory

| Data | Writers | Readers | Lifetime | Synchronization candidate | Aliasing risk |
| --- | --- | --- | --- | --- | --- |
| `[data]` | `[writers]` | `[readers]` | `[lifetime]` | `[mechanism]` | `[risk]` |

If a field has more than one writer, justify why the single-writer lifecycle claim still holds.

## Candidate synchronization designs

| Option | Mutation serialization | Wait mechanism | Read mechanism | Goroutines added per run | Benefits | Logical race risks |
| --- | --- | --- | --- | --- | --- | --- |
| Mutex plus completion signal | `[mechanism]` | `[mechanism]` | `[mechanism]` | `[count]` | `[benefits]` | `[risks]` |
| Ownership goroutine plus commands | `[mechanism]` | `[mechanism]` | `[mechanism]` | `[count]` | `[benefits]` | `[risks]` |
| Other credible design | `[mechanism]` | `[mechanism]` | `[mechanism]` | `[count]` | `[benefits]` | `[risks]` |

Compare shutdown, inspectability, blocking behavior, test control, and risk of user code running inside the mutation path.

## Waiter contract

| Schedule | Wait call behavior | Outcome returned | Caller wait-context behavior | Run affected? | Resource released by |
| --- | --- | --- | --- | --- | --- |
| Wait before completion | `[behavior]` | `[outcome]` | `[behavior]` | `[yes or no]` | `[owner]` |
| Several concurrent waiters | `[behavior]` | `[outcome]` | `[behavior]` | `[yes or no]` | `[owner]` |
| Wait after completion | `[behavior]` | `[outcome]` | `[behavior]` | `[yes or no]` | `[owner]` |
| One waiter abandons | `[behavior]` | `[outcome]` | `[behavior]` | `[yes or no]` | `[owner]` |

Decide whether `Wait` accepts a context. If it does, distinguish cancellation of waiting from cancellation of the run.

## Goroutine ownership

| Goroutine or async activity | Created by | Necessary? | Stop condition | Join or release path | Leak test |
| --- | --- | --- | --- | --- | --- |
| Work invocation | `[owner]` | `[reason]` | `[condition]` | `[path]` | `[test]` |
| Wait support | `[owner or none]` | `[reason]` | `[condition]` | `[path]` | `[test]` |
| Event support | `[owner or none]` | `[reason]` | `[condition]` | `[path]` | `[test]` |

Prefer designs that do not create a goroutine per waiter or observer unless that cost buys a required behavior.

## Lock and callback audit

List every operation that might occur while synchronization is held.

| Operation | Under lock or owner turn? | Can block? | Can call user code? | Can panic? | Decision |
| --- | --- | --- | --- | --- | --- |
| `[operation]` | `[yes or no]` | `[yes or no]` | `[yes or no]` | `[yes or no]` | `[move, retain, precompute]` |

No work function, observer callback, channel send to an unbounded consumer, or caller-controlled method may run inside lifecycle mutation.

## Publication and copying

| Published value | Internal representation | Snapshot or alias | Copy point | Mutation test | Guarantee stated to caller |
| --- | --- | --- | --- | --- | --- |
| State | `[representation]` | `[form]` | `[point]` | `[test]` | `[guarantee]` |
| Outcome | `[representation]` | `[form]` | `[point]` | `[test]` | `[guarantee]` |
| Event history | `[representation]` | `[form]` | `[point]` | `[test]` | `[guarantee]` |

## Area Decisions

| Decision | Mechanism | Race or leak prevented | Cost accepted | Sprint 2 impact |
| --- | --- | --- | --- | --- |
| `[decision]` | `[mechanism]` | `[failure]` | `[cost]` | `[impact]` |

## Trade-Offs

Analyze direct locking versus command serialization, context-aware waits versus a smaller API, copies versus allocation, broadcast closure versus condition variables, and goroutine count versus clearer ownership.

## Evidence

Use `Inputs Used` IDs in every material claim. Name the exact Go-runtime study report and section. When a repository mechanism matters, cite its study repository file and line range or symbol. Distinguish report interpretation from direct repository observation. Use prototype races, language-decision obligations, and exact current code. Explain why race-free memory access does not itself prove one terminal outcome or correct publication order.

## Risks

Include lock held around work, deadlock through nested inspection, outcome aliasing, waiter goroutine leaks, channel close races, a wait context accidentally cancelling the run, copied slice headers sharing mutable elements, and race tests that never enter the intended collision.

## Verification obligations

| Claim | Concurrent schedule | Race-detector role | Logical assertion | Quiescence assertion |
| --- | --- | --- | --- | --- |
| `[claim]` | `[schedule]` | `[what it checks]` | `[assertion]` | `[assertion]` |

## Self-critique

- Name every goroutine created for one successful run.
- Can one abandoned waiter retain the run or block completion?
- Which mutable reference escapes despite copying the container?
- What logical race would remain invisible under `go test -race`?
