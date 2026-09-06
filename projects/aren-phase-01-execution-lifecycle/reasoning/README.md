# Aren Phase 1 reasoning templates

These files are source templates, not completed reasoning. UltraPlan injects a selected sprint template and writes the completed document under the sprint's `reasoning/` directory.

The templates are grouped by decision scope:

```text
reasoning/
├── project-wide/   # Phase 1 research synthesis, for the planned project-reasoning flow
├── sprint-01/      # Core lifecycle area reasoning
└── sprint-02/      # Cancellation and concurrency area reasoning
```

The project-wide templates cannot yet run through the current sprint area-reasoning operation. They are ready for the planned UltraPlan project-reasoning feature. Do not select them as sprint area templates merely to work around that missing feature.

## Required reasoning standard

Every completed Aren reasoning document must do all of the following:

1. Separate four kinds of statement:
   - observed source or repository fact;
   - interpretation of that fact;
   - Aren-specific decision;
   - unresolved assumption.
2. Complete the `Inputs Used` section before analysis. Record only inputs actually read, using exact workspace-relative paths and the specific sections, files, symbols, or lines inspected.
3. Trace each material conclusion to an exact requirement, prior decision, study report, study repository reference, prototype result, or current-code reference.
4. Cite a study report by its exact report path and section. When repository behavior supports or challenges the report, also cite the concrete study repository file and line range or symbol. A study root, dimension directory, report directory, or repository name alone is not evidence.
5. Distinguish report interpretation from repository observation. If only the report was inspected, label the evidence `report-only`. If direct repository inspection changes or qualifies the report, record both sources and explain the difference.
6. Compare at least two credible designs when a genuine choice exists. Do not invent a weak alternative merely to reject it.
7. State the source system's product and scale assumptions before adopting its pattern.
8. Record negative evidence, contradictions, and missing evidence. Insufficient evidence must produce an experiment, test, or explicit deferral.
9. Explain what would falsify or reopen each material conclusion.
10. Map decisions to invariants and observable tests. A test name without the failure it detects is not enough.
11. Audit Phase 1 scope. Provider, tool, subprocess, persistence, workflow, daemon, and remote-execution designs must not enter through speculative seams.
12. Name accepted costs. "Simple", "robust", "flexible", and "idiomatic" are not rationales without a mechanism.
13. End with a self-critique that identifies the strongest objection, the most fragile assumption, and the easiest way the proposed design could lie about execution.

## Authority rules

- The PRD owns required behaviour and exclusions.
- Project-wide reasoning synthesizes evidence and records Phase 1 conclusions shared by both sprints.
- Sprint area reasoning makes detailed decisions for one sprint concern.
- Sprint `reasoning.md` is the authoritative sprint decision set.
- `plan.md` implements sprint `reasoning.md` and must not promote a provisional research conclusion into architecture.

## Evidence catalog and assessment

- `project-index.md` catalogs the studies, reports, and other evidence available to the project. It does not judge evidence quality or make conclusions.
- The project-reasoning index selects the evidence needed for project-wide reasoning and assigns it to decision areas.
- `project-wide/00-evidence-map.md` assesses only that selected evidence. It records canonicality, source quality, contradictions, applicability limits, negative transfer, gaps, and routing. It must not copy the project catalog.
- Each thematic project or sprint reasoning document records what it actually read in `Inputs Used` and cites exact report sections and study repository locations in its analysis.

## Project-wide templates

| Template | Purpose |
| --- | --- |
| `project-wide/00-evidence-map.md` | Assess selected evidence for canonicality, quality, conflicts, applicability, gaps, and routing. It does not duplicate the project catalog. |
| `project-wide/01-lifecycle-authority-and-atomic-publication.md` | Establish Phase 1 conclusions about transition ownership and coherent publication. |
| `project-wide/02-outcomes-failures-and-terminal-resolution.md` | Synthesize outcome, failure, panic, and terminal-arbitration evidence. |
| `project-wide/03-cancellation-goroutine-ownership-and-cleanup.md` | Synthesize truthful cancellation and Go ownership evidence. |
| `project-wide/04-events-observation-waiting-and-replay.md` | Synthesize canonical history, delivery, cursor, observer, and waiter evidence. |
| `project-wide/05-verification-and-go-correctness.md` | Define the Phase 1 proof model and permanent regression obligations. |

## Sprint 1 templates

| Template | Decision boundary |
| --- | --- |
| `sprint-01/runtime-architecture-and-authority.md` | Package ownership, dependency direction, work boundary, and capabilities. |
| `sprint-01/lifecycle-transition-and-publication.md` | State machine, transition gate, linearization, and terminal publication. |
| `sprint-01/outcome-and-failure-contract.md` | Result, failure, panic, timing, and immutable outcome rules. |
| `sprint-01/concurrency-waiters-and-immutability.md` | Locking, goroutine ownership, waiters, copying, and race boundaries. |
| `sprint-01/event-history-and-observation.md` | Minimum retained history and observation semantics earned in Sprint 1. |
| `sprint-01/lifecycle-verification-model.md` | State-model, negative, race, panic, waiter, and publication proof. |

## Sprint 2 templates

| Template | Decision boundary |
| --- | --- |
| `sprint-02/cancellation-and-terminal-resolution.md` | Acceptance, causes, acknowledgement, and outcome resolution. |
| `sprint-02/concurrency-interleavings-and-goroutine-ownership.md` | Cancellation races, watcher ownership, shutdown, and leak prevention. |
| `sprint-02/observers-replay-and-delivery.md` | Multiple observers, cursor replay, live handoff, and abandonment. |
| `sprint-02/control-observation-and-api-authority.md` | View, caller-control, and internal mutation capabilities. |
| `sprint-02/race-stress-and-leak-verification.md` | Adversarial schedules, repeated stress, negative controls, and leak proof. |
| `sprint-02/architecture-delta.md` | Conditional review of material changes to Sprint 1 architecture. |

## Selection rule

Select a template only when its decision boundary contains real sprint uncertainty. Sprint 2 must not select `architecture-delta.md` unless implementation evidence creates a material package, ownership, or dependency change.
