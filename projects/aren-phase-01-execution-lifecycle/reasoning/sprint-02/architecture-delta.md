# Area reasoning template: Sprint 2 architecture delta

> Selected template: `projects/aren-phase-01-execution-lifecycle/reasoning/sprint-02/architecture-delta.md`

## Selection gate

Select this template only when Sprint 2 evidence requires a material change to Sprint 1 package ownership, dependency direction, construction, or public/internal boundaries. Do not select it merely because Sprint 2 adds cancellation methods or observer behavior inside the established runtime owner.

## Purpose

Determine whether the realized Sprint 1 architecture still fits cancellation, multi-observer replay, and adversarial concurrency. Preserve working decisions unless new evidence demonstrates a concrete defect.

## Inputs Used

Complete this table before analysis. Record only material actually read.

| Input ID | Input | Kind | Exact workspace-relative path | Exact portion inspected | How it was used | Limits or applicability warning |
| --- | --- | --- | --- | --- | --- | --- |
| `[IN-NN]` | `[name]` | `[requirement, project reasoning, handbook, study report, study repository, prototype report, current code, Sprint 1 reasoning or review]` | `[path]` | `[report section; source file and lines or symbol; document heading]` | `[trigger, comparison, constraint, or challenge]` | `[report-only, product mismatch, scale mismatch, unverified interpretation, or none]` |

Name each study report by exact path and section. When repository behavior matters, add a separate row for the study repository file and line range or symbol. A study directory, report directory, or repository name alone is not evidence.

## Trigger evidence

| Trigger | Sprint 1 decision affected | Realized code reference | Failure or unacceptable cost | Why local extension is insufficient |
| --- | --- | --- | --- | --- |
| `[trigger]` | `[decision]` | `[path and lines]` | `[failure]` | `[reason]` |

If this table has no concrete entry, stop and mark this area not applicable.

## Prior Sprint Decisions Applied

| Sprint 1 architecture decision | Status | Evidence | Exact delta | Compatibility impact |
| --- | --- | --- | --- | --- |
| `[decision]` | `[preserved, extended, superseded, unaffected]` | `[evidence]` | `[delta]` | `[impact]` |

## Current dependency and authority map

Show the realized Sprint 1 graph, not the planned graph.

```text
[real current package and call relationships]
```

Mark:

- lifecycle mutation owner;
- work execution boundary;
- observer and waiter access;
- cancellation acceptance path;
- any dependency that now points in the wrong direction;
- any caller that can bypass authority.

## Change options

| Option | Packages or ownership changed | Defect fixed | New abstraction | Migration cost | New failure mode | Scope effect |
| --- | --- | --- | --- | --- | --- | --- |
| Extend current owner in place | `[change]` | `[defect]` | `[none or small]` | `[cost]` | `[failure]` | `[effect]` |
| Refactor before adding behavior | `[change]` | `[defect]` | `[abstraction]` | `[cost]` | `[failure]` | `[effect]` |
| Other credible option | `[change]` | `[defect]` | `[abstraction]` | `[cost]` | `[failure]` | `[effect]` |

## Abstraction recheck

| Abstraction | Sprint 1 evidence | New Sprint 2 variation | Does it remove present complexity? | Cost | Decision |
| --- | --- | --- | --- | --- | --- |
| `[interface, component, coordinator, registry]` | `[evidence]` | `[variation]` | `[yes or no]` | `[cost]` | `[decision]` |

Explicitly reject a daemon, durable repository, global event bus, executor registry, and workflow engine unless the PRD changed. Sprint 2 requirements do not earn them.

## Migration and intermediate correctness

Explain the order for moving existing behavior without creating a temporary second lifecycle owner. Identify tests that remain green before, during, and after the refactor. No intermediate execution path may bypass terminal or cancellation invariants.

## Area Decisions

| Decision | Existing decision changed | New mechanism | Evidence requiring change | Rejected smaller change |
| --- | --- | --- | --- | --- |
| `[decision]` | `[decision]` | `[mechanism]` | `[evidence]` | `[alternative]` |

## Trade-Offs

Analyze refactor risk versus preserving a strained design, package clarity versus fragmentation, compatibility versus correcting authority, and generalization versus the smallest Sprint 2 delta.

## Evidence

Use `Inputs Used` IDs in every material claim. Cite each study report by exact path and section and each relevant study repository mechanism by file and line range or symbol. Distinguish report interpretation from direct repository observation. Prioritize realized Sprint 1 code, tests, and review findings. Use project-wide reasoning and studies to interpret the pressure, not to replace direct implementation evidence.

## Risks

Include two temporary lifecycle owners, event or cancellation bypass during migration, package cycles, test seams becoming public abstractions, lost Sprint 1 behavior, and refactoring used to import later-phase infrastructure.

## Verification obligations

| Architecture claim | Before/after proof | Integration path | Regression prevented |
| --- | --- | --- | --- |
| `[claim]` | `[proof]` | `[path]` | `[regression]` |

## Self-critique

- Is this change required by evidence or merely cleaner on paper?
- Could the same defect be fixed inside the existing owner with less risk?
- Does any migration step permit duplicate authority?
- What later-phase abstraction is trying to enter under the label "refactor"?
