# Sprint Index: Reason Stage

> Project: `ultraplan-go`
> Sprint: `20-reason-stage`
> Purpose: selected context for this sprint. Must be a subset of `projects/ultraplan-go/project-index.md`.
> **Inputs Used:** `projects/ultraplan-go/sprints/20-reason-stage/requirements.md`, `projects/ultraplan-go/project-index.md`, `projects/ultraplan-go/docs/PRD.md`, `projects/ultraplan-go/docs/TRD.md`, `projects/ultraplan-go/docs/ARCHITECTURE.md`, `templates/sprint-index.md`

This document selects what must be read, distilled, reasoned through, or checked for this sprint. It does not make implementation decisions. All selections must come from the project index -- no items may be included that are not listed in the project index.

## Sprint Scope

- **Sprint Goal:** Implement the planning reason stage so selected area reasoning and final `reasoning.md` can be generated, previewed, flowed, and validated from valid requirements, `sprint-index.md`, and `technical-handbook.md`, with correct `flow-state.json` updates through `reasoning` and no plan or execution behavior.
- **Planned Output:** `internal/sprint` reason-stage domain, validators, prompt rendering, flow support, filesystem loading, flow-state updates, service methods, CLI wiring, tests, selected Architecture area reasoning artifact, final `reasoning.md`, and post-implementation sprint review evidence.
- **Depends On:** Valid `requirements.md`; valid `sprint-index.md`; valid `technical-handbook.md`; cataloged Phase 2 planning context in `project-index.md`; PRD, TRD, and Architecture docs for product and technical boundaries.
- **Non-Goals:** Generating or validating `plan.md`; task breakdowns or implementation checklists in reasoning artifacts; sprint implementation execution; smoke investigation execution; review automation; issue tracking; Git mutation; `.run-state.json`, `smoke.md`, `smoke.json`, generated `review.md`, `issues.md`, or `issues.json`; changing select or distill behavior beyond consuming valid outputs; using unselected contracts, evidence, templates, protocols, docs, or prior reviews as authoritative reasoning inputs; reusing study services or study-owned models; adding generic workflow, prompt, validation, plugin, hosted, browser, TUI, API, collaboration, or stable new JSON-output scope.

## Source Project Index

- `projects/ultraplan-go/project-index.md` -- authoritative source. Any file or item referenced below must appear there.

## Selected Contracts

Each contract applies as a flat whole to this sprint. All paths must appear in the project index's "Active Contract Pool" table.

| Contract | Why Selected |
| -------- | ------------ |
| Architecture | Required for `internal/sprint` ownership, dependency direction, product/platform separation, thin CLI adapters, and avoiding global planning/workflow/validation/prompt packages. |
| Errors | Required for actionable validation diagnostics, safe error messages, cause chains, deterministic ordering, and required exit-code mapping for validate, prompt, and flow commands. |
| Configuration | Required for runtime-facing configuration, command preflight behavior, model/variant/timeout overrides, and redaction of sensitive values in prompt, flow, and diagnostics output. |
| Observability | Required for calm scriptable diagnostics, runtime failure reporting through the generic runtime boundary, validation findings, and flow-state-visible failures. |
| Security | Required for workspace-safe path handling, selected-template path validation, no secrets in generated or previewed prompts, and prohibitions on direct shell, Git, provider, OpenCode, or agentwrap adapter invocation from `internal/sprint`. |
| Testing | Required for deterministic unit, fixture, command, fake-runtime, dry-run, validation, flow-state, race, and build verification. |
| Documentation | Required for user-facing help output and maintainable editable Markdown planning artifacts. |
| CLI Surface | Required for `validate`, `prompt`, and `flow` command behavior, stage argument handling, stdout/stderr separation, no ANSI output, help output, and exit codes. |
| LLM Runtime | Required because non-dry-run reason-stage flow must use only the generic platform runtime boundary and must not directly invoke OpenCode, provider APIs, shell commands, agentwrap adapters, or Git. |
| LLM Evaluation / Cost / Safety | Required for safe runtime execution discipline, validation-after-runtime success, safety boundaries, and avoiding unsafe prompt/runtime diagnostics. |
| Workflows | Required for Planning Phase 2 flow through `reasoning`, stage ordering, skipped area reasoning when no templates are selected, and rejection of unsupported future stages. |
| Persistence And Migrations | Required for strict versioned `flow-state.json`, atomic writes, unsupported state rejection, and preserving last valid state on write failure. |

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
| Architecture | `projects/ultraplan-go/sprints/20-reason-stage/reasoning/architecture.md` | Reason-stage implementation changes module ownership, dependency direction, package boundaries, prompt/rendering ownership, flow-state ownership, and runtime/product separation. |

## Prior Decisions To Carry Forward

All decision paths must appear in the project index's "Prior Decisions" table.

| Decision | Path | Constraint For This Sprint |
| -------- | ---- | -------------------------- |
| None cataloged | None cataloged in project index | The project index's "Prior Decisions" section says no prior decisions are cataloged, so no prior-decision artifact is selected here. Carry-forward constraints stated in `requirements.md` still constrain the sprint scope but are not selected as prior-decision artifacts. |

## Required Review Protocols

All paths must appear in the project index's "Review Protocols" table.

| Protocol | Path | Required Evidence |
| -------- | ---- | ----------------- |
| Architecture Review | `system/protocols/architecture-review-protocol.md` | Evidence that `internal/sprint` owns reasoning behavior, CLI adapters stay thin, platform packages remain product-agnostic, `internal/study` is not imported, selected context is workspace-safe, and no global planning/workflow/validation/prompt package is introduced. |
| Sprint Review | `system/protocols/sprint-review-protocol.md` | Evidence that the required outputs, acceptance criteria, verification commands, deviations, and conclusions are recorded after implementation. |

## Excluded Context

| Context | Reason Excluded | Revisit If |
| ------- | --------------- | ---------- |
| `08-concurrency` evidence report | Reason-stage flow is sequential stage orchestration with fake-runtime tests; no new scheduler, worker-pool, or concurrent run-loop behavior is selected for this sprint. | A later sprint changes scheduler, run-loop, cancellation fan-out, or bounded parallel execution behavior. |
| `12-extensibility` evidence report | Plugin systems, extension frameworks, hosted services, and broad extensibility points are non-goals; the sprint should keep behavior in `internal/sprint`. | A later sprint explicitly introduces extension points or stable external APIs. |
| `14-performance` evidence report | No performance-sensitive repository scanning, report processing, or scheduler work is selected beyond ordinary deterministic prompt/validation/flow behavior. | A later sprint changes startup latency, large artifact processing, worker pools, or scanning behavior. |
| Smoke harness `ultraplan-go-smoke` | Real-runtime smoke evidence is not required for default verification; requirements call for offline deterministic verification with fake runtime by default. | The completed sprint is selected for deep smoke evidence. |
| Deep Smoke Sprint protocol | The sprint requirements exclude real-runtime smoke tests for default verification and do not require deep-smoke execution for this index. | A reviewer explicitly requests real-runtime smoke evidence after implementation. |
| Plan-stage implementation context | Generating or validating `plan.md` belongs to Sprint 21 and must not be required by reasoning validation or flow through `reasoning`. | Sprint 21 starts or requirements expand reason-stage scope to include `plan.md`. |
| Implementation, execute, smoke, review automation, issues, and Git mutation stages | Phase 2 current behavior supports planning through `plan`; this sprint must reject unsupported future stages and must not create or validate deferred artifacts. | A later requirements revision brings these product stages into scope. |
| Unselected contracts, evidence reports, reasoning templates, protocols, project docs, and prior reviews | Requirements require selected-context-only reasoning and validation; final reasoning must not cite unselected context as authoritative. | `project-index.md` and a future sprint index explicitly select additional context. |
| Study services, source/dimension models, report/rating logic, summary generation, scheduler logic, and run-loop state | Phase 2 may reuse infrastructure but must not reuse or extract study-owned semantics for planning reasoning. | A concrete future requirement proves a stable shared mechanical helper is needed outside study semantics. |

## Next Artifacts

- `technical-handbook.md` reads from the evidence reports listed above.
- `reasoning/*.md` captures area-specific reasoning.
- `reasoning.md` makes final sprint decisions.
- `plan.md` executes `reasoning.md`.
- `review.md` runs the selected review protocols against implementation.
