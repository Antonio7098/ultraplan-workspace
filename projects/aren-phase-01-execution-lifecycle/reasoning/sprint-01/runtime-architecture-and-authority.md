# Area reasoning template: Aren runtime architecture and authority

> Selected template: `projects/aren-phase-01-execution-lifecycle/reasoning/sprint-01/runtime-architecture-and-authority.md`

## Purpose

Decide the smallest Sprint 1 package and authority design that can own one in-process run. This document owns package boundaries, dependency direction, capabilities, and the work invocation boundary. It does not own the detailed transition algorithm, outcome fields, event API, or test matrix.

Apply the reasoning standard in the project template README. Use the Phase 1 PRD, selected project-wide reasoning, technical handbook, code context, and current Aren source. Label observed code facts separately from proposed design.

## Inputs Used

Complete this table before analysis. Record only material actually read.

| Input ID | Input | Kind | Exact workspace-relative path | Exact portion inspected | How it was used | Limits or applicability warning |
| --- | --- | --- | --- | --- | --- | --- |
| `[IN-NN]` | `[name]` | `[requirement, project reasoning, handbook, study report, study repository, prototype report, current code, prior review]` | `[path]` | `[report section; source file and lines or symbol; document heading]` | `[decision, comparison, constraint, or challenge]` | `[report-only, product mismatch, scale mismatch, unverified interpretation, or none]` |

Name each study report by exact path and section. When repository behavior matters, add a separate row for the study repository file and line range or symbol. A study directory, report directory, or repository name alone is not evidence.

## Decision scope

- Goal: place lifecycle authority where every transition and publication path must pass through it.
- Must decide: package responsibilities, public/internal boundary, run handle or capability division, dependency construction, and work ownership.
- Must defer: cancellation API, providers, tools, subprocesses, persistence, daemon hosting, workflows, and a universal executor interface.

## Current architecture facts

| Fact | Exact source reference | Why it constrains this sprint | Observed or inferred |
| --- | --- | --- | --- |
| `[fact]` | `[repository path and lines]` | `[constraint]` | `[observed or inferred]` |

State whether Aren currently contains only prototypes, reusable production code, or both. Do not treat prototype package shape as accepted architecture.

## Required authority map

| Capability | Required owner | Who may request it | Who must never perform it | Observable result |
| --- | --- | --- | --- | --- |
| Generate run identity | `[owner]` | `[caller]` | `[excluded actor]` | `[fact]` |
| Start work | `[owner]` | `[caller]` | `[excluded actor]` | `[fact]` |
| Commit transition | `[owner]` | `[producer]` | `[excluded actor]` | `[fact]` |
| Inspect state | `[owner]` | `[reader]` | `[restriction]` | `[snapshot]` |
| Wait for outcome | `[owner]` | `[reader]` | `[restriction]` | `[outcome]` |
| Read history | `[owner]` | `[reader]` | `[restriction]` | `[events]` |

## Candidate package designs

Compare at least two designs that could plausibly meet Sprint 1.

| Option | Package layout | Lifecycle owner | Public API | Dependency direction | Benefits | Failure modes | Future pressure |
| --- | --- | --- | --- | --- | --- | --- | --- |
| Minimal concrete runtime package | `[layout]` | `[owner]` | `[API]` | `[direction]` | `[benefits]` | `[failures]` | `[pressure]` |
| Split domain and orchestration packages | `[layout]` | `[owner]` | `[API]` | `[direction]` | `[benefits]` | `[failures]` | `[pressure]` |
| Other credible design, if evidence supports it | `[layout]` | `[owner]` | `[API]` | `[direction]` | `[benefits]` | `[failures]` | `[pressure]` |

Reject an option because of a concrete ownership, testing, or dependency defect. File count alone is not a rationale.

## Work boundary

Define what Aren supplies to work and what work may return. Explain:

- why `context.Context` belongs at this boundary;
- whether result typing is generic, concrete, or intentionally provisional;
- how panic crosses the boundary;
- why work cannot mutate run state;
- where work runs relative to the lifecycle critical section;
- which API choices are semantic instruments rather than compatibility promises.

## Dependency direction

Provide the proposed dependency graph.

```text
[entry point]
    -> [runtime composition]
        -> [lifecycle owner]
            -> [small injected dependencies]
```

For every injected dependency, state the concrete variation or deterministic-test need that earns an interface. Prefer direct concrete values where no variation exists.

## Abstraction pressure audit

| Proposed abstraction | Current concrete variation | Testing need | Cost introduced | Earned now? | Decision |
| --- | --- | --- | --- | --- | --- |
| `[interface, factory, option, registry, generic type]` | `[variation]` | `[need]` | `[cost]` | `[yes or no]` | `[keep, narrow, remove]` |

Explicitly assess a general executor interface, event bus, repository/store, clock interface, plugin registry, and configuration object.

## Area Decisions

For each decision include the chosen design, mechanism, rationale, project-wide conclusion applied, rejected credible alternative, and evidence needed to prove it.

| Decision | Mechanism | Evidence and inherited constraint | Rejected alternative | Sprint 2 consequence |
| --- | --- | --- | --- | --- |
| `[decision]` | `[mechanism]` | `[evidence]` | `[alternative and reason]` | `[preserve, extend, or revisit]` |

## Trade-Offs

Cover at least package clarity versus fragmentation, API convenience versus capability separation, concrete code versus test seams, and provisional public shape versus future compatibility.

## Evidence

Use `Inputs Used` IDs in every material claim. Cite each study report by exact path and section and each relevant study repository mechanism by file and line range or symbol. Distinguish report interpretation from direct repository observation. Trace every material decision to selected project reasoning, current code, or a requirement. Include contrary evidence and explain why it does not control this sprint.

## Risks

Analyze bypass paths, lifecycle authority leaking into a CLI, a public API freezing too early, hidden global state, context used as a service locator, cyclic package dependencies, and future features entering through speculative interfaces.

## Verification and review obligations

| Architecture claim | Static or runtime check | Failure it detects | Review owner |
| --- | --- | --- | --- |
| `[claim]` | `[check]` | `[failure]` | `[plan task or review]` |

## Self-critique

- Can any supported caller publish state or outcome without the lifecycle owner?
- Which interface has only one implementation and no genuine test need?
- What is the strongest reason the proposed package split may be wrong after Sprint 1?
- Did the design make a later execution type easier by making Sprint 1 harder?
