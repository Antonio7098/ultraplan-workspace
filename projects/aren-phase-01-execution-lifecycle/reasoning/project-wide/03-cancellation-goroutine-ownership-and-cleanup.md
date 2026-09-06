# Project reasoning template: Cancellation, goroutine ownership, and cleanup

## Purpose

Synthesize Phase 1 evidence about cooperative cancellation, context causes, Aren-owned goroutines, bounded shutdown, and truthful terminal publication. Keep process escalation and durable recovery out of Phase 1.

Read `projects/aren-phase-01-execution-lifecycle/reasoning/README.md` and follow its required reasoning standard.

## Inputs Used

Complete this table before analysis. Record only material actually read.

| Input ID | Input | Kind | Exact workspace-relative path | Exact portion inspected | How it was used | Limits or applicability warning |
| --- | --- | --- | --- | --- | --- | --- |
| `[IN-NN]` | `[name]` | `[requirement, project reasoning, study report, study repository, prototype report, current code, prior review]` | `[path]` | `[report section; source file and lines or symbol; document heading]` | `[question, comparison, constraint, or challenge]` | `[report-only, product mismatch, scale mismatch, unverified interpretation, or none]` |

Name each study report by exact path and section. When repository behavior matters, add a separate row for the study repository file and line range or symbol. A study directory, report directory, or repository name alone is not evidence.

## Governing questions

- Who may request cancellation, and who accepts it?
- What does acceptance prove and what does it not prove?
- How does the first accepted cause reach work?
- Who owns every goroutine created for a run?
- What is each goroutine's stop condition and join path?
- When may Aren publish `cancelled`?
- What happens when work delays or ignores cancellation?

## Normative constraints

| Constraint | Source | Required meaning | Design freedom left |
| --- | --- | --- | --- |
| `[constraint]` | `[path and section]` | `[meaning]` | `[freedom]` |

## Evidence

Use `Inputs Used` IDs in every material claim. Cite the exact report section and, where relevant, the study repository file and line range or symbol. Distinguish report interpretation from direct repository observation.

### Cancellation protocol comparison

| Evidence | Request representation | Cause retention | Acknowledgement | Confirmed stop condition | Repeated request behavior | Aren applicability |
| --- | --- | --- | --- | --- | --- | --- |
| `[source]` | `[mechanism]` | `[mechanism]` | `[mechanism]` | `[condition]` | `[behavior]` | `[adopt, adapt, reject]` |

### Ownership evidence

| Evidence | Spawn owner | Child tracking | Cancellation source | Join mechanism | Abandonment behavior | Failure or leak tests |
| --- | --- | --- | --- | --- | --- | --- |
| `[source]` | `[owner]` | `[mechanism]` | `[source]` | `[mechanism]` | `[behavior]` | `[tests]` |

Distinguish evidence for in-process goroutines from processes, durable workflow tasks, server connections, and distributed workers.

## Cancellation state model

Model cancellation as facts rather than forcing it into the lifecycle state machine.

| Fact | Owner | First-write rule | Observable to | Effect |
| --- | --- | --- | --- | --- |
| Request received | `[owner]` | `[rule]` | `[audience]` | `[effect]` |
| Request accepted | `[owner]` | `[rule]` | `[audience]` | `[effect]` |
| Context cancelled | `[owner]` | `[rule]` | `[audience]` | `[effect]` |
| Work acknowledged cause | `[owner]` | `[rule]` | `[audience]` | `[effect]` |
| Terminal cancellation committed | `[owner]` | `[rule]` | `[audience]` | `[effect]` |

State which facts may occur without the later facts.

## Ownership tree

Produce a conceptual tree for Phase 1 and name every lifetime boundary.

```text
caller parent context
└── Aren run lifetime
    ├── work invocation
    ├── parent-cancellation watcher, if required
    ├── waiter notification mechanism
    └── observer support, if it creates goroutines
```

For every node, record:

| Node | Created by | Cancelled by | Completion observed by | Joined or released by | Can outlive run? |
| --- | --- | --- | --- | --- | --- |
| `[node]` | `[owner]` | `[source]` | `[observer]` | `[owner]` | `[yes or no with reason]` |

## Race analysis

Analyze at least these schedules:

1. completion immediately before explicit cancellation;
2. accepted cancellation immediately before successful return;
3. parent cancellation and explicit cancellation together;
4. accepted cancellation followed by unrelated error;
5. work ignores cancellation indefinitely;
6. cancellation arrives after terminal commitment;
7. waiter abandons its own wait context while the run continues;
8. run finishes while a parent watcher is starting or stopping.

For each, identify observable facts, terminal interpretation, and goroutine cleanup.

## Cleanup boundary

Phase 1 excludes general executor cleanup hooks. Explain which lifecycle cleanup still exists, such as releasing Aren-owned watchers, timers, or channels. Reject any design that silently imports Phase 2 cleanup semantics.

## Project conclusions

| Conclusion | Evidence basis | Sprint 1 obligation | Sprint 2 decision required | Reopen trigger |
| --- | --- | --- | --- | --- |
| `[conclusion]` | `[evidence]` | `[obligation]` | `[decision]` | `[trigger]` |

## Trade-Offs

Address watcher goroutine versus polling or callback integration, cause-aware context versus simpler cancel functions, blocking joins versus truthful return latency, and leak prevention versus unnecessary lifetime machinery.

## Risks

Include cancellation treated as completion, goroutines without owners, cleanup using an already-cancelled context, double close or send-after-close, abandoned consumers stranding producers, context cause mismatch, and process-oriented escalation leaking into Phase 1.

## Verification obligations

| Claim | Controlled schedule | Observable assertion | Leak or quiescence assertion | Negative control |
| --- | --- | --- | --- | --- |
| `[claim]` | `[schedule]` | `[assertion]` | `[assertion]` | `[deliberate defect]` |

## Self-critique

- Which goroutine lacks a named owner or join path?
- Could Aren ever report `cancelled` while work is still executing?
- What happens if cancellation propagation itself races with terminal publication?
- Did evidence from subprocess supervision cause Phase 1 to promise forceful stop?
