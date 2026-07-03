# Sprint Index: Select Stage

> Project: `ultraplan-go`
> Sprint: `18-select-stage`
> Purpose: selected context for this sprint. Must be a subset of `projects/ultraplan-go/project-index.md`.
> **Inputs Used:** `projects/ultraplan-go/project-index.md`, `projects/ultraplan-go/sprints/18-select-stage/requirements.md`, `projects/ultraplan-go/docs/PRD.md`, `projects/ultraplan-go/docs/TRD.md`, `projects/ultraplan-go/docs/ARCHITECTURE.md`, `projects/ultraplan-go/roadmap.md`, `projects/ultraplan-go/sprints/17-artifact-domain-and-flow-state/review.md`, `projects/ultraplan-go/sprints/16-project-domain-and-index/review.md`, `projects/ultraplan-go/sprints/14-validation-and-diagnostics/review.md`, `projects/ultraplan-go/sprints/09-runtime-integration/review.md`, `templates/sprint-index.md`

This document selects what must be read, distilled, reasoned through, or checked for this sprint. It does not make implementation decisions. All selections must come from the project index — no items may be included that are not listed in the project index.

## Sprint Scope

- **Sprint Goal:** Implement the planning select stage so `sprint-index.md` can be generated, previewed, and validated as the authoritative context selection for a sprint, with subset checks against `project-index.md` and flow-state updates through `sprint-index` only.
- **Planned Output:** A valid `projects/ultraplan-go/sprints/18-select-stage/sprint-index.md` plus select-stage parsing, validation, prompt rendering, flow execution, flow-state updates, CLI wiring, and tests in the Go implementation.
- **Depends On:** Sprint 16 project catalog parsing and validation, Sprint 17 sprint artifact domain and atomic flow-state handling, Sprint 14 validation diagnostics and safe output conventions, Sprint 9 generic platform runtime integration, the `ultraplan-go` project catalog, roadmap, PRD, TRD, and architecture document.
- **Non-Goals:** Do not generate or validate `technical-handbook.md`, `reasoning/*.md`, `reasoning.md`, or `plan.md`; do not implement `flow --to technical-handbook`, `flow --to area-reasoning`, `flow --to reasoning`, or `flow --to plan` except unsupported-stage rejection where required; do not implement sprint implementation execution, smoke investigation execution, review automation, issue tracking, project-management behavior, Git mutation, study-side execution behavior, `.run-state.json`, `smoke.md`, `smoke.json`, `issues.md`, or `issues.json`; do not mutate `project-index.md`, `roadmap.md`, project docs, prior sprint reviews, source repositories, workspace config, or Git state.

## Source Project Index

- `projects/ultraplan-go/project-index.md` — authoritative source. Any file or item referenced below must appear there.

## Selected Contracts

Each contract applies as a flat whole to this sprint. All paths must appear in the project index's "Active Contract Pool" table.

| Contract | Why Selected |
| -------- | ------------ |
| Architecture | Governs `internal/sprint` ownership, dependency direction, thin CLI adapters, product/platform separation, and the prohibition on global planning, workflow, validation, reports, or prompts packages. |
| Errors | Governs deterministic validation diagnostics, preserved cause chains, actionable repair guidance, and exit-code mapping for validate, prompt, and flow failures. |
| Configuration | Governs runtime config override handling, command preflight behavior, and redaction of config or provider-sensitive values in prompt previews and diagnostics. |
| Observability | Governs calm status output, safe diagnostics, runtime metadata summaries, validation findings, and truthful flow-state reporting. |
| Security | Governs workspace-safe sprint paths, generated artifact paths, secret redaction, runtime permission posture, and the prohibition on direct shell/OpenCode/Git/provider calls from sprint logic. |
| Testing | Governs fixture-first tests, fake-runtime tests, deterministic diagnostics, command tests, offline verification, race verification, and CLI build checks. |
| Documentation | Governs user-facing help, editable Markdown artifacts, prompt-preview clarity, and review evidence for the sprint. |
| CLI Surface | Governs command shape, help text, stdout/stderr separation, no ANSI output, scriptable text, and exit codes for `validate sprint-index`, `prompt sprint-index`, and `flow --to sprint-index`. |
| LLM Runtime | Governs routing runtime-backed flow through the generic platform runtime boundary and forbids direct OpenCode, agentwrap adapter, provider API, or shell invocation from `internal/sprint`. |
| Workflows | Governs select-stage flow behavior, dry-run behavior, post-runtime validation gates, stage transitions, failure handling, and exclusion of unsupported implementation/smoke/review/issues stages. |
| Persistence And Migrations | Governs strict versioned `flow-state.json`, atomic writes, preservation of last valid state on write failure, and deterministic loading/rejection of invalid durable state. |

## Selected Evidence Reports

Copied from the project index's "Available Evidence Reports" table. These tell the technical handbook which reports to read – the project index is the authoritative source.

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
| Architecture | `projects/ultraplan-go/sprints/18-select-stage/reasoning/architecture.md` | Needed to reason through `internal/sprint` ownership, `sprint -> project/workspace/platform` dependencies, platform runtime separation, flow-state boundaries, and exclusion of study semantics. |

## Prior Decisions To Carry Forward

All decision paths must appear in the project index's "Prior Decisions" table.

| Decision | Path | Constraint For This Sprint |
| -------- | ---- | -------------------------- |
| No cataloged prior decisions | Not applicable | The project index records no selectable prior-decision entries. Requirements still require preserving prior sprint outcomes: product behavior remains module-owned, CLI adapters stay thin, platform packages stay product-agnostic, cause chains are preserved, paths remain workspace-safe, state writes remain atomic, and speculative shared planning abstractions are avoided. |

## Required Review Protocols

All paths must appear in the project index's "Review Protocols" table.

| Protocol | Path | Required Evidence |
| -------- | ---- | ----------------- |
| Architecture Review | `system/protocols/architecture-review-protocol.md` | Verify `internal/sprint` owns selection, validation, prompt rendering, and flow rules; `internal/app` remains thin; platform runtime stays product-agnostic; no forbidden global packages or study-semantic reuse are introduced. |
| Sprint Review | `system/protocols/sprint-review-protocol.md` | Verify all Sprint 18 outputs, selected contracts, acceptance criteria, tests, validation behavior, prompt preview behavior, dry-run behavior, flow-state transitions, excluded stages, and verification commands are satisfied. |

## Excluded Context

| Context | Reason Excluded | Revisit If |
| ------- | --------------- | ---------- |
| `technical-handbook.md` generation and validation | Belongs to Sprint 19 distill stage, not the Sprint 18 select stage. | Sprint 19 starts or a later requirement explicitly expands Sprint 18. |
| Area reasoning and final `reasoning.md` generation or validation | Belongs to Sprint 20 reason stage; Sprint 18 only selects reasoning templates and output paths. | Sprint 20 starts or select-stage requirements are revised. |
| `plan.md` generation or validation | Belongs to Sprint 21 plan stage; Sprint 18 stops at validated `sprint-index.md` and readying later stages. | Sprint 21 starts or flow scope is expanded. |
| Sprint implementation execution | Phase 2 explicitly stops at planning artifacts and does not execute implementation tasks. | A future requirements revision adds implementation execution to the product. |
| Smoke investigation execution and deep smoke evidence | Sprint 18 must not run real-runtime smoke or model smoke stages; default verification is offline and fake-first. | A completed sprint explicitly requests external smoke evidence outside Planning Phase 2 artifact flow. |
| Automated conformance review and Planning Phase 2 `review.md` generation | Review evidence is produced by the development process, not by product flow stages. | A future product phase adds review automation. |
| Issue tracking and issue artifacts | Project-management and local issue tracking behavior are deferred non-goals. | A future product phase adds issue-management scope. |
| Git mutation | Automatic Git add, commit, push, or other Git mutation is out of scope and must not be invoked by select-stage commands. | A future product phase adds explicit opt-in Git behavior. |
| Study-side execution semantics | Planning may reference cataloged study reports as evidence paths but must not import or reuse study services, source/dimension models, report validators, ratings, summaries, scheduler, or run-loop behavior. | A repeated mechanical helper is proven generic and contains no study semantics. |
| Unselected evidence reports `08-concurrency`, `12-extensibility`, and `14-performance` | Sprint 18 does not introduce worker-pool scheduling, plugin architecture, public extension points, or performance-sensitive repository scanning beyond bounded catalog and sprint-file reads. | Later stages introduce concurrency, extension, or performance-sensitive behavior. |
| Deep Smoke Sprint protocol | Not required because Sprint 18 default verification must not require real OpenCode, provider credentials, network access, or real-runtime smoke tests. | The sprint is complete and maintainers explicitly request external smoke validation. |

## Next Artifacts

- `technical-handbook.md` reads from the evidence reports listed above.
- `reasoning/*.md` captures area-specific reasoning.
- `reasoning.md` makes final sprint decisions.
- `plan.md` executes `reasoning.md`.
- `review.md` runs the selected review protocols against implementation.
