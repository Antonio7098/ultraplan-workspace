# Sprint Index: TUI Foundation and Read-Only Dashboard

> Project: `ultraplan-go`
> Sprint: `24-tui-foundations`
> Purpose: selected context for this sprint. Must be a subset of `projects/ultraplan-go/project-index.md`.
> **Inputs Used:** `projects/ultraplan-go/sprints/24-tui-foundations/requirements.md`, `projects/ultraplan-go/project-index.md`, `projects/ultraplan-go/docs/ARCHITECTURE.md`, `projects/ultraplan-go/docs/PRD.md`, `projects/ultraplan-go/docs/TRD.md`

This document selects what must be read, distilled, reasoned through, or checked for this sprint. It does not make implementation decisions. All selections must come from the project index — no items may be included that are not listed in the project index.

## Sprint Scope

- **Sprint Goal:** Add a read-only `ultraplan tui` dashboard over shared application use cases so users can inspect workspace, project, study, and sprint state without invoking runtimes, mutating artifacts, shelling out to UltraPlan, or parsing CLI output.
- **Planned Output:** TUI command wiring, shared read-only app use cases, read-only project/study/sprint use cases, `internal/tui` documentation, program setup, model/update logic, rendering, key bindings, fakes, and deterministic tests for command startup, use cases, model behavior, rendering, error states, previews, and read-only guarantees.
- **Depends On:** Sprint 16 project domain and index; Sprint 17 sprint artifact domain and flow state; Sprints 18-21 planning stages through `plan`; Sprint 22 planning documentation and release; Sprint 23 execute stage, especially typed execute status and redacted execute diagnostics; existing workspace/config/app composition.
- **Non-Goals:** Running validation, planning flow, prompt preview, execute, study run-loop, or other mutating/runtime-backed workflows from the TUI; guarded action dialogs; live runtime progress panes; cancellation controls; operational workflow controls; browser UI; hosted service; API server; multi-user collaboration; project-management features; smoke artifacts; automated review artifacts; issue artifacts; Git mutation; CLI behavior replacement; separate TUI persistence; workflow scheduler; global validation/workflow/scheduler/report/prompt/stage package refactors; broad Sprint 23 hardening beyond safe read-only TUI display needs.

## Source Project Index

- `projects/ultraplan-go/project-index.md` — authoritative source. Any file or item referenced below must appear there.

## Selected Contracts

Each contract applies as a flat whole to this sprint. All paths must appear in the project index's "Active Contract Pool" table.

| Contract | Why Selected |
| -------- | ------------ |
| Architecture | Governs `internal/app` as shared local-interface boundary, `internal/tui` ownership, package dependency direction, product/platform separation, and avoidance of global workflow packages. |
| CLI Surface | `ultraplan tui`, `--workspace`, help visibility, startup errors, exit behavior, and preservation of existing CLI text/JSON automation surfaces are in scope. |
| Configuration | TUI startup must reuse workspace discovery and normal config validation while staying runtime-free and redacting config/runtime values shown in the dashboard. |
| Documentation | `internal/tui/doc.go`, command help expectations, and clear documentation of TUI dependency, ownership boundary, read-only scope, and non-goals are required outputs. |
| Errors | Startup failures, preview failures, invalid workspace/config errors, use-case errors, and actionable TUI error panes must preserve clear classification and diagnostics. |
| Observability | The dashboard displays project, study, sprint, execute status, validation findings, diagnostics, and safe run-state summaries without launching workflows. |
| Performance | TUI startup, refresh, deterministic ordering, narrow-terminal rendering, and bounded artifact previews must avoid unbounded reads or expensive repository scans. |
| Persistence And Migrations | The TUI reads durable workspace artifacts, `flow-state.json`, `execute.md`, and `.run-state.json` as source-of-truth state through product use cases without creating a separate persistence model. |
| Security | Path-safe previews, workspace-root containment, no subprocess self-invocation, no runtime credential requirement, and secret/redaction guarantees are central constraints. |
| Testing | Non-interactive model, rendering, command, use-case, fake-use-case, error-state, and read-only guarantee tests are required acceptance criteria. |
| Workflows | The TUI must display workflow state and execute status safely while excluding workflow execution, preserving product ownership of sprint/study state machines. |

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
| `12-extensibility` | `studies/go-cli-study/reports/final/12-extensibility.md` | Extension points, plugin architecture, package boundaries, API design |
| `13-security` | `studies/go-cli-study/reports/final/13-security.md` | Path safety, secrets, command injection risks, sandboxing |
| `14-performance` | `studies/go-cli-study/reports/final/14-performance.md` | Startup latency, large repos, memory management, performance |
| `15-philosophy` | `studies/go-cli-study/reports/final/15-philosophy.md` | Cross-cutting design philosophy, tradeoffs, maintainability |

## Selected Reasoning Templates

All paths must appear in the project index's "Available Reasoning Templates" table.

| Template | Output Path | Why Selected |
| -------- | ----------- | ------------ |
| Architecture | `projects/ultraplan-go/sprints/24-tui-foundations/reasoning/architecture.md` | Required to reason through `internal/app` use-case extraction, `internal/tui` boundaries, dependency direction, terminal-library containment, read-only scope, and Sprint 23 status/redaction carry-forward constraints. |

## Prior Decisions To Carry Forward

All decision paths must appear in the project index's "Prior Decisions" table.

| Decision | Path | Constraint For This Sprint |
| -------- | ---- | -------------------------- |
| No prior decisions listed | `projects/ultraplan-go/project-index.md` | The project index's Prior Decisions section states that there are no prior decision artifacts; carry-forward constraints for this sprint come from the sprint requirements dependencies and selected source documents. |

## Required Review Protocols

All paths must appear in the project index's "Review Protocols" table.

| Protocol | Path | Required Evidence |
| -------- | ---- | ----------------- |
| Architecture Review | `system/protocols/architecture-review-protocol.md` | Evidence that `internal/tui` depends only on app use cases and simple result types, terminal-library types do not leak into product modules, CLI/TUI share typed use cases, product modules do not import `internal/tui`, and no global workflow package refactor was introduced. |
| Sprint Review | `system/protocols/sprint-review-protocol.md` | Completed-sprint evidence against requirements, acceptance criteria, selected contracts, deterministic tests, `go test ./...`, and `go build ./cmd/ultraplan`. |

## Excluded Context

| Context | Reason Excluded | Revisit If |
| ------- | --------------- | ---------- |
| Implementation execution | The sprint is a read-only TUI foundation; executing `plan.md` tasks, launching runtime work, and mutating execute state from the TUI are non-goals. | Sprint 25 or later explicitly scopes guarded operational workflow controls. |
| Smoke investigation | Project scope and sprint requirements defer smoke investigation, and the TUI must not generate smoke artifacts or run deep smoke behavior. | A later sprint explicitly adds smoke workflows or external smoke evidence collection to TUI scope. |
| Review automation | Automated review artifact generation is deferred and the TUI must not produce review artifacts. | A later sprint explicitly brings conformance review automation into product scope. |
| Issue tracking | Issue artifacts and issue-management behavior are deferred project non-goals. | A later product requirement adds issue integration or local issue artifacts. |
| Git mutation | Automatic Git add, commit, push, checkout, reset, branch, merge, or other Git mutation is explicitly excluded. | A later product requirement adds explicit, guarded Git integration. |
| LLM Runtime | Runtime execution, provider calls, prompt execution, repair, and workflow operation are outside this read-only TUI baseline. | TUI operational workflow controls enter scope with explicit runtime-backed requirements. |
| LLM Evaluation / Cost / Safety | Cost/evaluation behavior is not implemented in this sprint; only safe display/redaction of existing runtime metadata is relevant through Security and Observability. | Runtime-backed TUI operation or cost-aware dashboard behavior enters scope. |
| Deep Smoke Sprint protocol | This sprint does not require real-runtime smoke evidence and must start without runtime credentials, OpenCode availability, provider configuration beyond normal non-runtime config validation, or network access. | Acceptance criteria change to require real-runtime smoke evidence. |
| Browser UI and hosted service | The sprint only adds a local terminal UI in the existing CLI binary. | Product scope changes to local API, browser UI, hosted service, or multi-user operation. |
| Separate TUI persistence model | Durable workspace artifacts remain the source of truth, and TUI state is derived and refreshable. | A future requirement introduces an explicit cache or persistence model with migration rules. |
| Global workflow or validation refactor | Requirements prohibit refactoring product semantics into global `validation`, `workflow`, `scheduler`, `reports`, `prompts`, or `stages` packages. | Repeated concrete implementations prove a shared abstraction is stable and a later sprint scopes it. |
| Study run-loop operation | The TUI may inspect study status but must not start run-loop, validation, synthesis, or other mutating/runtime-backed study actions. | A later sprint scopes guarded TUI workflow operation. |
| Planning flow and prompt preview operation | The TUI may inspect sprint planning state and artifacts but must not run planning flow or prompt preview actions in this baseline. | A later sprint scopes guarded dry-run or prompt-preview controls. |

## Next Artifacts

- `technical-handbook.md` reads from the evidence reports listed above.
- `reasoning/*.md` captures area-specific reasoning.
- `reasoning.md` makes final sprint decisions.
- `plan.md` executes `reasoning.md`.
- `review.md` runs the selected review protocols against implementation.

Hard constraints:
- Select only entries listed in project-index.md.
- Use workspace-relative paths.
- Do not mutate project-index.md, roadmap.md, docs, source repositories, config, Git state, or any artifact other than sprint-index.md.
