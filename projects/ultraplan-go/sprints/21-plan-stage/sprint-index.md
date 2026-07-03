# Sprint Index: Plan Stage

> Project: `ultraplan-go`
> Sprint: `21-plan-stage`
> Purpose: selected context for this sprint. Must be a subset of `projects/ultraplan-go/project-index.md`.
> **Inputs Used:** `projects/ultraplan-go/sprints/21-plan-stage/requirements.md`, `projects/ultraplan-go/project-index.md`, `projects/ultraplan-go/docs/PRD.md`, `projects/ultraplan-go/docs/TRD.md`, `projects/ultraplan-go/docs/ARCHITECTURE.md`, `templates/sprint-index.md`

This document selects what must be read, distilled, reasoned through, or checked for this sprint. It does not make implementation decisions. All selections must come from the project index -- no items may be included that are not listed in the project index.

## Sprint Scope

- **Sprint Goal:** Implement the planning plan stage so validated `reasoning.md` can be previewed, flowed, generated, and validated into `plan.md`, with `flow-state.json` updated through `plan` and no implementation, smoke, review, issue, or Git behavior.
- **Planned Output:** `internal/sprint` plan-stage domain and validation, deterministic plan prompt rendering, flow execution through `plan`, service and filesystem-store updates, flow-state updates, thin CLI wiring, plan and command tests, selected Architecture area reasoning artifact, final `reasoning.md`, `plan.md`, and post-implementation sprint review evidence.
- **Depends On:** Sprint 20 reason stage, Sprint 19 distill stage, Sprint 18 select stage, Sprint 17 sprint artifact domain and flow state, Sprint 16 project catalog boundary, Sprint 14 diagnostics and exit-code conventions, Sprint 9 generic runtime boundary, the cataloged `ultraplan-go` project index, PRD, TRD, and Architecture docs.
- **Non-Goals:** Executing implementation tasks from `plan.md`; running smoke investigations; producing, validating, or automating product-generated conformance reviews; creating or managing local issues; mutating Git state; generating or requiring `.run-state.json`, `smoke.md`, `smoke.json`, product-generated `review.md`, `issues.md`, or `issues.json`; changing select, distill, or reason-stage contracts beyond consuming valid prerequisites; reading unselected context as authoritative plan input; mutating project catalog, docs, roadmap, prior reviews, selected evidence, source repositories, workspace config, or Git state; reusing study services or study-owned models; adding generic workflow, prompt, validation, plugin, hosted, browser, TUI, API, collaboration, or multi-stage implementation run-loop behavior.

## Source Project Index

- `projects/ultraplan-go/project-index.md` -- authoritative source. Any file or item referenced below must appear there.

## Selected Contracts

Each contract applies as a flat whole to this sprint. All paths must appear in the project index's "Active Contract Pool" table.

| Contract | Why Selected |
| -------- | ------------ |
| Architecture | Required for `internal/sprint` ownership of plan validation, prompt rendering, flow-state transitions, and prerequisite reasoning checks; also governs thin CLI adapters, product/platform separation, and avoiding global planning/workflow/validation/prompt packages. |
| Errors | Required for actionable plan validation diagnostics, cause preservation, deterministic findings, safe runtime failure messages, and required exit-code mapping for validate, prompt, dry-run, and flow commands. |
| Configuration | Required for runtime-facing model, variant, timeout, preflight, and redaction behavior in prompt previews, flow execution, and diagnostics. |
| Observability | Required for calm scriptable output, stdout/stderr separation, validation findings, runtime-failure reporting through the generic runtime boundary, and flow-state-visible plan failures. |
| Security | Required for workspace-safe artifact and reasoning paths, secret redaction, no unsafe prompt diagnostics, and prohibitions on direct shell, Git, provider, OpenCode, or agentwrap adapter invocation from `internal/sprint`. |
| Testing | Required for deterministic unit, fixture, command, fake-runtime, dry-run, plan-validation, flow-state, race, and build verification. |
| Documentation | Required for user-facing help output, prompt-preview clarity, editable Markdown planning artifacts, and maintainable review evidence. |
| CLI Surface | Required for `validate plan`, `prompt plan`, `flow --to plan`, unsupported-stage rejection, help output, stdout/stderr separation, no ANSI output, scriptable text, and exit codes. |
| LLM Runtime | Required because non-dry-run plan flow must use only the generic platform runtime boundary and must not directly invoke OpenCode, provider APIs, shell commands, agentwrap adapters, or Git. |
| LLM Evaluation / Cost / Safety | Required for safe runtime execution discipline, validation-after-runtime success, safe metadata handling, and avoiding unsafe prompt/runtime diagnostics. |
| Workflows | Required for Planning Phase 2 flow through `plan`, dry-run behavior, prerequisite validation, stage ordering, skipped area reasoning when no templates are selected, failure handling, and rejection of unsupported future stages. |
| Persistence And Migrations | Required for strict versioned `flow-state.json`, atomic writes, unsupported state rejection, and preserving the last valid state on write failure. |

## Selected Evidence Reports

Copied from the project index's "Available Evidence Reports" table. These tell the technical handbook which reports to read - the project index is the authoritative source.

| Report | Path | Covers |
| ------ | ---- | ------ |
| `01-project-structure` | `studies/go-cli-study/reports/final/01-project-structure.md` | Project layout, cmd/internal/pkg, dependency direction, thin entrypoints |
| `02-command-architecture` | `studies/go-cli-study/reports/final/02-command-architecture.md` | Command routing, flags, help text, command organization, shell completion |
| `03-dependency-injection` | `studies/go-cli-study/reports/final/03-dependency-injection.md` | Dependency construction, seams, testability, constructor patterns |
| `04-configuration-management` | `studies/go-cli-study/reports/final/04-configuration-management.md` | Config loading, precedence, environment variables, paths |
| `05-error-handling` | `studies/go-cli-study/reports/final/05-error-handling.md` | Error wrapping, classification, user-facing diagnostics, sentinel errors |
| `06-io-abstraction` | `studies/go-cli-study/reports/final/06-io-abstraction.md` | Filesystem/stdin/stdout abstraction, test seams, interface design |
| `07-state-context` | `studies/go-cli-study/reports/final/07-state-context.md` | Context propagation, app state, cancellation |
| `09-terminal-ux` | `studies/go-cli-study/reports/final/09-terminal-ux.md` | Terminal output, progress indicators, human UX, color and formatting |
| `10-logging-observability` | `studies/go-cli-study/reports/final/10-logging-observability.md` | Logs, diagnostics, structured events, observability patterns |
| `11-testing-strategy` | `studies/go-cli-study/reports/final/11-testing-strategy.md` | Unit, integration, fixture, command-level tests, coverage strategy |
| `13-security` | `studies/go-cli-study/reports/final/13-security.md` | Path safety, secrets, command injection risks, sandboxing |
| `15-philosophy` | `studies/go-cli-study/reports/final/15-philosophy.md` | Cross-cutting design philosophy, tradeoffs, maintainability |

## Selected Reasoning Templates

All paths must appear in the project index's "Available Reasoning Templates" table.

| Template | Output Path | Why Selected |
| -------- | ----------- | ------------ |
| Architecture | `projects/ultraplan-go/sprints/21-plan-stage/reasoning/architecture.md` | Plan-stage implementation changes module ownership, dependency direction, package boundaries, prompt/rendering ownership, flow-state ownership, prerequisite reasoning handling, and runtime/product separation. |

## Prior Decisions To Carry Forward

All decision paths must appear in the project index's "Prior Decisions" table.

| Decision | Path | Constraint For This Sprint |
| -------- | ---- | -------------------------- |
| No cataloged prior decisions | None cataloged in project index | The project index records no selectable prior-decision artifacts. Carry-forward constraints stated in `requirements.md` still constrain the sprint: keep product behavior module-owned, keep CLI adapters thin, keep platform packages product-agnostic, preserve cause chains, use workspace-safe paths, keep state writes atomic, avoid speculative shared planning abstractions, consume valid Sprint 20 reasoning and Sprint 19 handbook artifacts without regenerating them, preserve Sprint 18 selected-context validation, and preserve Sprint 17 exclusion of implementation, smoke, review, issue, and Git stages. |

## Required Review Protocols

All paths must appear in the project index's "Review Protocols" table.

| Protocol | Path | Required Evidence |
| -------- | ---- | ----------------- |
| Architecture Review | `system/protocols/architecture-review-protocol.md` | Evidence that `internal/sprint` owns plan validation, prompt rendering, reasoning prerequisite checks, task/evidence trace checks, and flow-state transitions; `internal/app` remains thin; platform packages remain product-agnostic; `internal/study` is not imported; selected context is workspace-safe; and no global planning/workflow/validation/prompt package is introduced. |
| Sprint Review | `system/protocols/sprint-review-protocol.md` | Evidence in `projects/ultraplan-go/sprints/21-plan-stage/review.md` covering implementation results, deviations, verification commands, selected-context constraints, plan-stage acceptance criteria, and review conclusions. |

## Excluded Context

| Context | Reason Excluded | Revisit If |
| ------- | --------------- | ---------- |
| `08-concurrency` evidence report | Plan-stage flow is sequential stage orchestration with fake-runtime tests; no new scheduler, worker-pool, or concurrent run-loop behavior is selected for this sprint. | A later sprint changes scheduler, run-loop, cancellation fan-out, or bounded parallel execution behavior. |
| `12-extensibility` evidence report | Plugin systems, extension frameworks, hosted services, and broad extensibility points are non-goals; the sprint should keep behavior in `internal/sprint`. | A later sprint explicitly introduces extension points or stable external APIs. |
| `14-performance` evidence report | No performance-sensitive repository scanning, report processing, or scheduler work is selected beyond ordinary deterministic prompt, validation, selected artifact loading, and flow behavior. | A later sprint changes startup latency, large artifact processing, worker pools, or scanning behavior. |
| Smoke harness `ultraplan-go-smoke` | Real-runtime smoke evidence is not required for default verification; requirements call for deterministic offline verification and fake runtime by default. | The completed sprint is selected for deep smoke evidence. |
| Deep Smoke Sprint protocol | The sprint requirements exclude real-runtime smoke tests for default verification and do not require deep-smoke execution for this index. | A reviewer explicitly requests real-runtime smoke evidence after implementation. |
| Select, distill, and reason-stage implementation changes | Sprint 21 consumes valid prerequisites and must not change those stage contracts beyond loading, validating, and flowing prerequisites needed for `plan`. | A later requirements revision expands plan-stage scope to alter prior planning stages. |
| Implementation, execute, smoke, review automation, issues, and Git mutation stages | Phase 2 current behavior supports planning through `plan`; this sprint must reject unsupported future stages and must not create, validate, invoke, or mark deferred artifacts or commands. | A later requirements revision brings these product stages into scope. |
| `.run-state.json`, `smoke.md`, `smoke.json`, product-generated `review.md`, `issues.md`, and `issues.json` | These are not supported Planning Phase 2 artifacts and must not be generated, validated, required, or modeled in current flow state. | A future requirements revision moves those artifacts into scope. |
| Unselected contracts, evidence reports, reasoning templates, protocols, project docs, and prior reviews | Requirements require selected-context-only planning; `plan.md` and prompt rendering must not cite unselected context as authoritative. | `project-index.md` and a future sprint index explicitly select additional context. |
| Study services, source/dimension models, report/rating logic, summary generation, scheduler logic, and run-loop state | Phase 2 may reuse infrastructure but must not reuse or extract study-owned semantics for planning plan behavior. | A concrete future requirement proves a stable shared mechanical helper is needed outside study semantics. |
| Mutating project catalog, docs, roadmap, prior reviews, selected evidence reports, source repositories, workspace config, or Git state | The plan stage consumes selected context and writes only expected sprint planning artifacts under the selected sprint root. | A separate maintenance sprint authorizes catalog, docs, evidence, source, config, or Git mutation. |

## Next Artifacts

- `technical-handbook.md` reads from the evidence reports listed above.
- `reasoning/*.md` captures area-specific reasoning.
- `reasoning.md` makes final sprint decisions.
- `plan.md` executes `reasoning.md`.
- `review.md` runs the selected review protocols against implementation.
