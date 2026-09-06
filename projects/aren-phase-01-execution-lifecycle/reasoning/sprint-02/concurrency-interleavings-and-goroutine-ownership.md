# Area reasoning template: Concurrency interleavings and goroutine ownership

> Selected template: `projects/aren-phase-01-execution-lifecycle/reasoning/sprint-02/concurrency-interleavings-and-goroutine-ownership.md`

## Purpose

Attack the realized Sprint 1 design with cancellation, parent watchers, multiple waiters, observers, and concurrent inspection. Decide the exact synchronization and ownership changes required for Sprint 2 without reopening unaffected architecture.

## Inputs Used

Complete this table before analysis. Record only material actually read.

| Input ID | Input | Kind | Exact workspace-relative path | Exact portion inspected | How it was used | Limits or applicability warning |
| --- | --- | --- | --- | --- | --- | --- |
| `[IN-NN]` | `[name]` | `[requirement, project reasoning, handbook, study report, study repository, prototype report, current code, Sprint 1 reasoning or review]` | `[path]` | `[report section; source file and lines or symbol; document heading]` | `[decision, interleaving, constraint, or challenge]` | `[report-only, product mismatch, scale mismatch, unverified interpretation, or none]` |

Name each study report by exact path and section. When repository behavior matters, add a separate row for the study repository file and line range or symbol. A study directory, report directory, or repository name alone is not evidence.

## Prior Sprint Decisions Applied

| Sprint 1 decision | Realized implementation reference | Status | New pressure | Change required |
| --- | --- | --- | --- | --- |
| `[decision]` | `[code or review path]` | `[preserved, extended, superseded, unaffected]` | `[pressure]` | `[change]` |

## Concurrency inventory

List every goroutine and concurrent caller after Sprint 2.

| Actor | Starts when | Reads | Writes or proposes | Blocks on | Stops when | Owned and joined by |
| --- | --- | --- | --- | --- | --- | --- |
| Work goroutine | `[time]` | `[data]` | `[facts]` | `[wait]` | `[condition]` | `[owner]` |
| Parent watcher | `[time]` | `[data]` | `[facts]` | `[wait]` | `[condition]` | `[owner]` |
| Explicit canceller | `[time]` | `[data]` | `[facts]` | `[wait]` | `[condition]` | `[owner]` |
| Waiters | `[time]` | `[data]` | `[facts]` | `[wait]` | `[condition]` | `[owner]` |
| Observers | `[time]` | `[data]` | `[facts]` | `[wait]` | `[condition]` | `[owner]` |

If an actor has no owner, the design is incomplete.

## Synchronization delta

| Shared fact | Sprint 1 mechanism | New concurrent access | Existing mechanism sufficient? | Required change | New risk |
| --- | --- | --- | --- | --- | --- |
| `[fact]` | `[mechanism]` | `[access]` | `[yes or no]` | `[change]` | `[risk]` |

Explain why each change is required. Avoid replacing a working Sprint 1 mechanism merely because Sprint 2 adds complexity elsewhere.

## Interleaving catalog

For each schedule, identify the linearization points, legal observations, terminal winner, and quiescence condition.

| Interleaving | Controlled order | Linearization points | Legal intermediate observations | Final outcome | Goroutines remaining |
| --- | --- | --- | --- | --- | --- |
| Completion versus explicit cancel | `[order]` | `[points]` | `[views]` | `[outcome]` | `[count or none]` |
| Completion versus parent cancel | `[order]` | `[points]` | `[views]` | `[outcome]` | `[count or none]` |
| Explicit versus parent cancel | `[order]` | `[points]` | `[views]` | `[outcome]` | `[count or none]` |
| Terminal commit versus observer registration | `[order]` | `[points]` | `[views]` | `[outcome]` | `[count or none]` |
| Observer abandonment versus append | `[order]` | `[points]` | `[views]` | `[outcome]` | `[count or none]` |
| Waiter cancellation versus run completion | `[order]` | `[points]` | `[views]` | `[outcome]` | `[count or none]` |

Include both sides of every two-way race and at least one three-actor schedule.

## Lock and channel discipline

| Operation | Synchronization held | May block | Consumer-controlled | Panic path | Safe? | Required change |
| --- | --- | --- | --- | --- | --- | --- |
| Accept cancellation | `[mechanism]` | `[yes or no]` | `[yes or no]` | `[path]` | `[assessment]` | `[change]` |
| Append event | `[mechanism]` | `[yes or no]` | `[yes or no]` | `[path]` | `[assessment]` | `[change]` |
| Notify observer | `[mechanism]` | `[yes or no]` | `[yes or no]` | `[path]` | `[assessment]` | `[change]` |
| Release waiters | `[mechanism]` | `[yes or no]` | `[yes or no]` | `[path]` | `[assessment]` | `[change]` |

No observer-controlled send may occur inside lifecycle mutation.

## Shutdown and quiescence protocol

Write the normal and cancellation shutdown sequence. State who closes or cancels each signal, who waits, and what happens if work ignores cancellation. Distinguish run terminal publication from release of observer resources.

## Area Decisions

| Decision | Mechanism | Interleaving handled | Cost accepted | Sprint 1 decision status |
| --- | --- | --- | --- | --- |
| `[decision]` | `[mechanism]` | `[interleaving]` | `[cost]` | `[status]` |

## Trade-Offs

Analyze additional goroutines versus direct synchronization, watcher simplicity versus lifecycle ownership, bounded observer queues versus pull cursors, immediate resource release versus late replay, and global quiescence assertions versus per-run ownership proof.

## Evidence

Use `Inputs Used` IDs in every material claim. Cite selected reports by exact path and section and relevant study repository mechanisms by file and line range or symbol. Distinguish report interpretation from direct repository observation. Use realized Sprint 1 code and review evidence first. Then apply project-wide Go ownership reasoning and prototype failures. A generic concurrency pattern does not override known run semantics.

## Risks

Include startup and shutdown goroutine leaks, close-send races, observer delivery under lock, cancellation watchers retained after success, wait cancellation affecting the run, lock inversion between history and terminal commit, and a race-free implementation with incorrect terminal interpretation.

## Verification obligations

| Interleaving or ownership claim | Deterministic control | Race assertion | Quiescence assertion | Negative control |
| --- | --- | --- | --- | --- |
| `[claim]` | `[control]` | `[assertion]` | `[assertion]` | `[defect]` |

## Self-critique

- Which actor can still run after the terminal outcome is visible, and is that truthful?
- Which channel can be closed by more than one path?
- What happens if parent cancellation arrives while the watcher is being released?
- Which logical race remains even if every memory access is synchronized?
