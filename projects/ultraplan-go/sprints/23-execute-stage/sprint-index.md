# Sprint Index: 23-execute-stage

> Project: `ultraplan-go`
> Sprint: `23-execute-stage`
> Purpose: selected context for this sprint. Must be a subset of `projects/ultraplan-go/project-index.md`.
> **Inputs Used:** `projects/ultraplan-go/project-index.md`, `projects/ultraplan-go/sprints/23-execute-stage/requirements.md`, `projects/ultraplan-go/docs/ARCHITECTURE.md`, `projects/ultraplan-go/docs/PRD.md`, `projects/ultraplan-go/docs/TRD.md`

This document selects what must be read, distilled, reasoned through, or checked for this sprint. It does not make implementation decisions. All selections must come from the project index - no items may be included that are not listed in the project index.

## Sprint Scope

- **Sprint Goal:** Execute validated `plan.md` implementation tasks through the generic runtime boundary with durable task state, safe target-repository boundaries, resumability, and clear diagnostics.
- **Planned Output:** `sprint-index.md`, `technical-handbook.md`, `reasoning.md`, `plan.md`, execute prompt rendering, deterministic task extraction, `.run-state.json`, `execute.md`, flow through `execute`, and status updates showing execute progress.
- **Depends On:** Project index catalog validation, valid planning prerequisites through `plan`, runtime model config, and flow-state persistence.
- **Non-Goals:** Smoke investigation; review automation; issue tracking; automatic Git mutation; hosted SaaS; browser UI; project-management features; cross-sprint scheduling.

## Source Project Index

- `projects/ultraplan-go/project-index.md` - authoritative source. Any file or item referenced below must appear there.

## Selected Contracts

Each contract applies as a flat whole to this sprint. All paths must appear in the project index's "Active Contract Pool" table.

| Contract | Why Selected |
| ------------ | -------------------------------------------- |
| Architecture | Execute-stage planning must preserve module ownership, dependency direction, and product/platform separation for `internal/sprint`, `internal/project`, `workspace`, and `platform/runtime`. |
| CLI Surface | Execute-stage planning concerns sprint commands, stage validation, status output, prompt rendering, and scriptable behavior. |
| Configuration | Execute-stage planning needs global and per-stage model selection, runtime config mapping, preflight validation, and redaction. |
| Documentation | This sprint produces editable planning Markdown artifacts and must keep generated docs readable and maintainable. |
| Errors | Execute-stage planning must account for actionable diagnostics, validation failures, runtime failures, cancellation, and task failure reporting. |
| LLM Evaluation / Cost / Safety | Execute-stage runtime planning needs usage/cost metadata, safety discipline, and bounded validation expectations where runtime metadata is available. |
| LLM Runtime | Execute-stage planning depends on agentwrap/OpenCode integration boundaries, runtime behavior, provider/model mapping, and adapter use. |
| Observability | Execute-stage planning needs task metadata, flow/run state visibility, diagnostics, structured logs, and truthful status output. |
| Performance | Execute-stage planning must avoid unbounded work, respect startup/status performance, and plan for large artifacts and bounded execution. |
| Persistence And Migrations | Execute-stage planning depends on durable `flow-state.json`, `.run-state.json`, atomic writes, schema versions, and resumability. |
| Security | Execute-stage planning must preserve workspace path safety, secret redaction, runtime permission policy, source isolation, and no automatic Git mutation. |
| Testing | Planning must identify unit, fixture, fake-runtime, and gated integration expectations for execute-stage behavior. |
| Workflows | Execute-stage planning is governed by stateful workflow execution, retries, cancellation, resumability, and controlled task execution. |

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
| Architecture | `projects/ultraplan-go/sprints/23-execute-stage/reasoning/architecture.md` | Execute-stage planning must reason through module boundaries, dependency direction, package layout, and runtime/product separation before final sprint decisions. |

## Prior Decisions To Carry Forward

All decision paths must appear in the project index's "Prior Decisions" table.

| Decision | Path | Constraint For This Sprint |
| ------------ | -------- | -------------------------- |
| None yet | N/A | No prior decisions are listed in the project index. |

## Required Review Protocols

All paths must appear in the project index's "Review Protocols" table.

| Protocol | Path | Required Evidence |
| ------------ | --------------------------------------- | ----------------- |
| Architecture Review | `system/protocols/architecture-review-protocol.md` | Evidence that execute-stage planning preserves package layout, dependency boundaries, and runtime/module separation. |
| Sprint Review | `system/protocols/sprint-review-protocol.md` | Evidence that completed planning artifacts satisfy sprint requirements, selected catalog constraints, and acceptance criteria. |

## Excluded Context

| Context | Reason Excluded | Revisit If |
| ----------- | --------------- | ------------- |
| Implementation execution outside validated `plan.md` tasks | The sprint includes controlled execute behavior only for validated `plan.md` implementation tasks through the runtime boundary; ad hoc execution remains excluded. | A later sprint adds broader execution modes or operators. |
| Smoke investigation | Sprint requirements, PRD, TRD, and project scope exclude smoke investigation from planning generation and Phase 2 product behavior. | A later sprint explicitly requires real-runtime smoke evidence outside the planning flow. |
| Deep Smoke Sprint protocol | Smoke investigation is a non-goal for this sprint and deferred by PRD/TRD scope. | A completed implementation sprint needs real-runtime smoke evidence. |
| Review automation | Sprint requirements and project scope defer review automation; this sprint may select review protocols for later human/gated review evidence but does not implement automated review execution. | A later sprint explicitly adds review automation as product behavior. |
| `12-extensibility` evidence report | Execute-stage planning is focused on controlled execution, state, diagnostics, runtime boundaries, and validation rather than plugin or extension architecture. | The sprint scope expands to new extension points, plugin behavior, or public APIs. |
| Hosted SaaS and browser UI context | PRD and TRD exclude hosted service and browser UI from the first production release. | Product scope changes to include hosted or browser-based operation. |
| Issue tracking | PRD and TRD explicitly defer issue tracking, assignment, scheduling, and project-management features. | A later requirements revision adds issue or project-management scope. |
| Git mutation | Requirements, PRD, TRD, and project index explicitly exclude automatic Git mutation from planning and execute behavior. | A later sprint adds explicit opt-in Git hook or mutation requirements. |
| Study implementation internals | Architecture and project index require Phase 2 to reuse infrastructure, not study services, source/dimension models, report validation, rating parsing, summary generation, or run-loop scheduling. | A concrete shared infrastructure need appears that is product-neutral and stable. |

## Next Artifacts

- `technical-handbook.md` reads from the evidence reports listed above.
- `reasoning/*.md` captures area-specific reasoning.
- `reasoning.md` makes final sprint decisions.
- `plan.md` turns `reasoning.md` decisions into ordered implementation tasks and evidence checks.
- Execute behavior runs validated `plan.md` tasks through the generic runtime boundary and records `.run-state.json` plus `execute.md`.
