# Sprint Index: 23-execute-stage

> Project: `ultraplan-go`
> Sprint: `23-execute-stage`
> Purpose: selected context for this sprint. Must be a subset of `projects/ultraplan-go/project-index.md`.
> **Inputs Used:** `projects/ultraplan-go/project-index.md`, `projects/ultraplan-go/roadmap.md`, `projects/ultraplan-go/sprints/23-execute-stage/requirements.md`, `projects/ultraplan-go/docs/ARCHITECTURE.md`, `projects/ultraplan-go/docs/PRD.md`, `projects/ultraplan-go/docs/TRD.md`

This document selects what must be read, distilled, reasoned through, or checked for this sprint. It does not make implementation decisions. All selections must come from the project index - no items may be included that are not listed in the project index.

## Sprint Scope

- **Sprint Goal:** Execute validated `plan.md` implementation tasks through the generic runtime boundary with durable task state, safe target-repository boundaries, resumability, and clear diagnostics.
- **Planned Output:** Context selection for execute-stage implementation planning, followed by `technical-handbook.md`, `reasoning/architecture.md`, `reasoning.md`, `plan.md`, `execute.md`, `.run-state.json`, `flow-state.json`, execute domain/persistence/flow/prompt/service/validation updates, CLI command wiring, and execute-focused tests.
- **Depends On:** Valid `requirements.md`, project catalog validation, completed planning prerequisites through `plan.md`, runtime model configuration, flow-state persistence, and the explicit target implementation directory `../ultraplan-go`.
- **Non-Goals:** Smoke investigation; automated review artifact generation; issue tracking; Git state mutation; hosted SaaS; browser UI; project-management workflows; generic workflow engine; cross-sprint scheduling; reuse of study services or study-specific models for execute behavior.

## Source Project Index

- `projects/ultraplan-go/project-index.md` - authoritative source. Any file or item referenced below must appear there.

## Selected Contracts

Each contract applies as a flat whole to this sprint. All paths must appear in the project index's "Active Contract Pool" table.

| Contract | Why Selected |
| ------------ | -------------------------------------------- |
| Architecture | Execute behavior must preserve `internal/sprint` ownership, allowed dependency direction, thin app/CLI wiring, and generic `platform/runtime` separation. |
| CLI Surface | Sprint 23 adds or updates `validate execute`, `prompt execute`, `flow --to execute`, `execute`, task selection, dry-run, help, exit codes, text output, and status behavior. |
| Configuration | Execute requires global/default sprint model resolution plus stage-specific overrides, command preflight validation, config diagnostics, and redaction. |
| Documentation | Execute artifacts and help/status output must remain readable, maintainable, and aligned with editable Markdown artifact expectations. |
| Errors | Execute prerequisites, task extraction, runtime failures, validation failures, stale state recovery, and target path violations need actionable diagnostics and exit behavior. |
| LLM Evaluation / Cost / Safety | Runtime-backed execute work must treat runtime success as insufficient, record safe metadata where available, and avoid unsafe payload or secret leakage. |
| LLM Runtime | Execute must use the existing generic runtime boundary backed by agentwrap/OpenCode and must not parse native runtime streams or invent a competing runtime contract. |
| Observability | Execute status, logs, run-state metadata, diagnostics, attempts, task counts, and runtime summaries must be inspectable without requiring runtime calls. |
| Performance | Execute task extraction, status, resume, and state loading must remain bounded and deterministic for larger plans and repositories. |
| Persistence And Migrations | `.run-state.json` and `flow-state.json` require versioned schemas, strict loading, same-directory atomic writes, and resumable recovery behavior. |
| Security | Execute writes target another repository, so path safety, workspace-safe diagnostics, target containment, permissions, and secret redaction are central constraints. |
| Testing | Acceptance requires deterministic offline fake-runtime tests, command tests, race tests, build verification, and non-goal regression coverage. |
| Workflows | Execute is a stateful workflow stage with prerequisite gates, task states, attempts, resumability, cancellation/stale recovery, and terminal-state summaries. |

## Selected Evidence Reports

Copied from the project index's "Available Evidence Reports" table. These tell the technical handbook which reports to read - the project index is the authoritative source.

| Report | Path | Covers |
| ----------- | ------------------------------------------------------ | ------------------------------- |
| `01-project-structure` | `studies/go-cli-study/reports/final/01-project-structure.md` | Project layout, cmd/internal/pkg, dependency direction, thin entrypoints |
| `02-command-architecture` | `studies/go-cli-study/reports/final/02-command-architecture.md` | Command routing, flags, help text, command organization, shell completion |
| `03-dependency-injection` | `studies/go-cli-study/reports/final/03-dependency-injection.md` | Dependency construction, seams, testability, constructor patterns |
| `04-configuration-management` | `studies/go-cli-study/reports/final/04-configuration-management.md` | Config loading, precedence, environment variables, paths |
| `05-error-handling` | `studies/go-cli-study/reports/final/05-error-handling.md` | Error wrapping, classification, user-facing diagnostics, sentinel errors |
| `06-io-abstraction` | `studies/go-cli-study/reports/final/06-io-abstraction.md` | Filesystem/stdin/stdout abstraction, test seams, interface design |
| `07-state-context` | `studies/go-cli-study/reports/final/07-state-context.md` | Context propagation, app state, cancellation |
| `08-concurrency` | `studies/go-cli-study/reports/final/08-concurrency.md` | Worker pools, cancellation, parallel execution, goroutine management |
| `09-terminal-ux` | `studies/go-cli-study/reports/final/09-terminal-ux.md` | Terminal output, progress indicators, human UX, color and formatting |
| `10-logging-observability` | `studies/go-cli-study/reports/final/10-logging-observability.md` | Logs, diagnostics, structured events, observability patterns |
| `11-testing-strategy` | `studies/go-cli-study/reports/final/11-testing-strategy.md` | Unit, integration, fixture, command-level tests, coverage strategy |
| `13-security` | `studies/go-cli-study/reports/final/13-security.md` | Path safety, secrets, command injection risks, sandboxing |
| `14-performance` | `studies/go-cli-study/reports/final/14-performance.md` | Startup latency, large repos, memory management, performance |
| `15-philosophy` | `studies/go-cli-study/reports/final/15-philosophy.md` | Cross-cutting design philosophy, tradeoffs, maintainability |

## Selected Reasoning Templates

All paths must appear in the project index's "Available Reasoning Templates" table.

| Template | Output Path | Why Selected |
| ------------ | -------------------------------------------------------------- | ------------ |
| Architecture | `projects/ultraplan-go/sprints/23-execute-stage/reasoning/architecture.md` | Execute-stage planning must reason through module ownership, dependency direction, package layout, runtime/product separation, persistence ownership, and target-repository safety before final sprint decisions. |

## Prior Decisions To Carry Forward

All decision paths must appear in the project index's "Prior Decisions" table.

| Decision | Path | Constraint For This Sprint |
| ------------ | -------- | -------------------------- |
| None yet | None yet - this is the first project. | No prior decisions are cataloged in the project index; carry forward only the roadmap, PRD, TRD, architecture, requirements, and selected catalog entries. |

## Required Review Protocols

All paths must appear in the project index's "Review Protocols" table.

| Protocol | Path | Required Evidence |
| ------------ | --------------------------------------- | ----------------- |
| Architecture Review | `system/protocols/architecture-review-protocol.md` | Evidence that execute implementation preserves package layout, dependency boundaries, sprint-owned product semantics, generic runtime/module separation, and target path safety. |
| Sprint Review | `system/protocols/sprint-review-protocol.md` | Evidence that completed sprint artifacts and implementation satisfy selected contracts, requirements acceptance criteria, validators, tests, and deferred-behavior exclusions. |

## Excluded Context

| Context | Reason Excluded | Revisit If |
| ----------- | --------------- | ------------- |
| Implementation execution outside validated `plan.md` tasks | Sprint 23 brings controlled implementation execution into scope only through validated `plan.md` task extraction and the generic runtime boundary. Ad hoc execution, broad operators, and cross-sprint execution are not selected. | A later sprint explicitly adds broader execution modes or cross-sprint orchestration. |
| Smoke investigation | Requirements, PRD, TRD, roadmap, and project index explicitly defer smoke investigation and `smoke.md` / `smoke.json` artifacts. | A later sprint explicitly selects smoke investigation or the Deep Smoke Sprint protocol as product scope. |
| Review automation | Requirements, PRD, TRD, roadmap, and project index defer automated conformance `review.md` generation and review automation. Review protocols are selected only as later review checks against completed implementation. | A later sprint explicitly adds automated review generation as product behavior. |
| Issue tracking | Requirements, PRD, TRD, roadmap, and project index defer issue files, issue JSON, assignment, scheduling, and project-management workflows. | A later requirements revision adds issue or project-management scope. |
| Git mutation | Requirements, PRD, TRD, roadmap, and project index exclude automatic Git add, commit, push, branch, PR, checkout, reset, or other Git state mutation. | A later sprint adds explicit opt-in Git hook or mutation requirements. |
| Deep Smoke Sprint protocol | Smoke investigation is out of scope for this sprint, so the cataloged smoke protocol is not selected as a required review protocol. | A completed implementation sprint needs real-runtime smoke evidence as selected scope. |
| `12-extensibility` evidence report | Execute-stage work is about controlled execution, state, runtime boundaries, validation, diagnostics, and command behavior rather than plugin architecture or new public extension points. | The sprint scope expands to plugin behavior, extension APIs, or public extension contracts. |
| Study services and study semantics | Architecture, TRD, and project index require Phase 2 to reuse generic infrastructure only, not study services, source/dimension models, report validation, rating parsing, summary generation, or run-loop scheduling. | A concrete shared infrastructure need appears that is product-neutral and stable. |
| Hosted SaaS and browser UI | PRD, TRD, roadmap, and project index exclude hosted service and browser UI from Phase 2 and the first production release. | Product scope changes to include hosted or browser-based operation. |
| TUI behavior | Sprint 23 precedes the post-execute TUI phase and does not implement `internal/tui` or terminal dashboard behavior. | Sprint 24 or later selects TUI foundation or operational controls. |

## Next Artifacts

- `technical-handbook.md` reads from the evidence reports listed above.
- `reasoning/*.md` captures area-specific reasoning.
- `reasoning.md` makes final sprint decisions.
- `plan.md` executes `reasoning.md`.
- `review.md` runs the selected review protocols against implementation.
