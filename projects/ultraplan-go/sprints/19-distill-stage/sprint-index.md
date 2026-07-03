# Sprint Index: Distill Stage

> Project: `ultraplan-go`
> Sprint: `19-distill-stage`
> Purpose: selected context for this sprint. Must be a subset of `projects/ultraplan-go/project-index.md`.
> **Inputs Used:** `projects/ultraplan-go/sprints/19-distill-stage/requirements.md`, `projects/ultraplan-go/project-index.md`, `projects/ultraplan-go/docs/PRD.md`, `projects/ultraplan-go/docs/TRD.md`, `projects/ultraplan-go/docs/ARCHITECTURE.md`, `templates/sprint-index.md`

This document selects what must be read, distilled, reasoned through, or checked for this sprint. It does not make implementation decisions. All selections must come from the project index -- no items may be included that are not listed in the project index.

## Sprint Scope

- **Sprint Goal:** Implement the planning distill stage so `technical-handbook.md` can be generated, previewed, flowed, and validated from valid `sprint-index.md` selected evidence only, with flow-state updates through `technical-handbook` and without making implementation decisions.
- **Planned Output:** Distill-stage support for `technical-handbook.md`, including selected-evidence loading, prompt preview, validation, flow-state transitions, CLI wiring, tests, and sprint review evidence.
- **Depends On:** Sprint 18 select stage, Sprint 17 sprint artifact domain and flow state, Sprint 16 project catalog boundary, Sprint 14 diagnostics and exit-code conventions, Sprint 9 generic runtime boundary, and the cataloged `go-cli-study` evidence reports selected below.
- **Non-Goals:** Generating or validating `reasoning/*.md`, final `reasoning.md`, or `plan.md`; making implementation or architecture decisions in `technical-handbook.md`; reading or citing unselected evidence reports; sprint execution, smoke investigation, review automation, issue tracking, Git mutation, study-service reuse, global workflow/prompt/validation packages, stable new JSON output, and real-runtime smoke requirements.

## Source Project Index

- `projects/ultraplan-go/project-index.md` -- authoritative source. Any file or item referenced below must appear there.

## Selected Contracts

Each contract applies as a flat whole to this sprint. All paths must appear in the project index's "Active Contract Pool" table.

| Contract | Why Selected |
| -------- | ------------ |
| Architecture | Governs `internal/sprint` ownership, `sprint -> project/workspace/platform` dependencies, thin CLI wiring, and product/platform separation. |
| Errors | Needed for deterministic validation diagnostics, cause preservation, exit-code mapping, and safe runtime failure reporting. |
| Configuration | Needed for runtime configuration boundaries, prompt rendering inputs, config precedence, and redaction requirements. |
| Observability | Needed for calm diagnostics, flow status, runtime metadata boundaries, and safe reporting of validation/runtime failures. |
| Security | Needed for workspace-safe selected evidence paths, secret redaction, no direct shell/OpenCode/Git invocation, and generated prompt safety. |
| Testing | Needed for offline unit/fixture/fake-runtime coverage, command tests, race verification, and deterministic diagnostics. |
| Documentation | Needed for CLI help updates and generated Markdown artifact expectations. |
| CLI Surface | Needed for `validate technical-handbook`, `prompt technical-handbook`, `flow --to technical-handbook`, stdout/stderr behavior, and unsupported-stage rejection. |
| LLM Runtime | Needed because non-dry-run flow must use only the generic runtime boundary backed by agentwrap/OpenCode integration. |
| Workflows | Needed for stage flow behavior, dry-run behavior, failure handling, and next-stage readiness through Planning Phase 2. |
| Persistence And Migrations | Needed for strict versioned `flow-state.json`, atomic writes, durable artifact paths, and rejection of unsupported legacy stages. |

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
| Architecture | `projects/ultraplan-go/sprints/19-distill-stage/reasoning/architecture.md` | Applies to module boundaries, dependency direction, package layout, runtime/product separation, and the next area-reasoning stage after distillation. |

## Prior Decisions To Carry Forward

All decision paths must appear in the project index's "Prior Decisions" table.

| Decision | Path | Constraint For This Sprint |
| -------- | ---- | -------------------------- |
| No cataloged prior decisions | N/A | The project index records no prior-decision artifacts; carry-forward constraints are taken from the sprint requirements dependencies and selected contracts rather than invented prior-decision paths. |

## Required Review Protocols

All paths must appear in the project index's "Review Protocols" table.

| Protocol | Path | Required Evidence |
| -------- | ---- | ----------------- |
| Architecture Review | `system/protocols/architecture-review-protocol.md` | Evidence that `internal/sprint` owns handbook validation, prompt rendering, selected-evidence loading, and flow behavior; `internal/app` remains thin; platform packages stay product-agnostic; and no global planning/workflow/validation/prompt package is introduced. |
| Sprint Review | `system/protocols/sprint-review-protocol.md` | Evidence in `projects/ultraplan-go/sprints/19-distill-stage/review.md` covering implementation results, deviations, verification commands, selected-evidence constraints, and acceptance-criteria conclusions. |

## Excluded Context

| Context | Reason Excluded | Revisit If |
| ------- | --------------- | ---------- |
| `08-concurrency` evidence report | Distill-stage flow is not adding worker pools, parallel scheduling, or concurrency primitives beyond existing runtime/context boundaries. | A later sprint changes planning flow to run multiple stages or runtime jobs concurrently. |
| `12-extensibility` evidence report | Plugin architecture and extension points are explicit non-goals for this sprint. | A later requirements revision adds supported extension/plugin behavior. |
| `14-performance` evidence report | No sprint requirement targets performance tuning beyond deterministic ordering and safe file reads covered by selected contracts. | Evidence loading or prompt rendering introduces measurable large-artifact performance constraints. |
| Deep Smoke Sprint protocol | The sprint requires offline, fake-first verification and explicitly does not require real OpenCode, provider credentials, network access, or real-runtime smoke tests. | A completed sprint is selected for real-runtime smoke evidence. |
| `ultraplan-go-smoke` harness | Real-runtime smoke execution is outside this sprint's required verification and non-goals. | Real-runtime smoke evidence is explicitly requested after implementation. |
| Unselected evidence reports | `technical-handbook.md` must read and cite only evidence reports selected by this sprint index. | The sprint index is deliberately revised before handbook generation. |
| Study services and study-owned semantics | Requirements forbid importing `internal/study` or reusing study source, dimension, report, rating, summary, scheduler, or run-loop behavior for planning distillation. | A future architecture decision explicitly extracts proven generic infrastructure without product semantics. |
| `reasoning/*.md`, `reasoning.md`, and `plan.md` generation or validation | Sprint 19 stops at distillation; reasoning belongs to Sprint 20 and plan belongs to Sprint 21. | The project advances to the reason or plan sprint. |
| Implementation, execute, smoke, review automation, issues, and Git stages | Planning Phase 2 supports only `requirements`, `sprint-index`, `technical-handbook`, `area-reasoning`, `reasoning`, and `plan`; later stages are deferred. | A later product phase explicitly adds execution, smoke, review automation, issues, or Git mutation. |
| `.run-state.json`, `smoke.md`, `smoke.json`, product-generated `review.md`, `issues.md`, and `issues.json` | These are not supported Planning Phase 2 artifacts and must not be generated, validated, or marked supported by this sprint. | A future requirements revision moves those artifacts into scope. |
| Mutating project catalog, docs, roadmap, prior reviews, selected evidence reports, source repositories, workspace config, or Git state | The distill stage consumes selected context and writes sprint planning artifacts only under the selected sprint root. | A separate maintenance sprint authorizes catalog, docs, evidence, source, config, or Git mutation. |

## Next Artifacts

- `technical-handbook.md` reads from the evidence reports listed above.
- `reasoning/*.md` captures area-specific reasoning.
- `reasoning.md` makes final sprint decisions.
- `plan.md` executes `reasoning.md`.
- `review.md` runs the selected review protocols against implementation.
