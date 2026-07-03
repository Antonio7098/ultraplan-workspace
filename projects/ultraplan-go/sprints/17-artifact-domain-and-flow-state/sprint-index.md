# Sprint Index: Sprint Artifact Domain and Flow State

> Project: `ultraplan-go`
> Sprint: `17-artifact-domain-and-flow-state`
> Purpose: selected context for this sprint. Must be a subset of `projects/ultraplan-go/project-index.md`.
> **Inputs Used:** `projects/ultraplan-go/project-index.md`, `projects/ultraplan-go/sprints/17-artifact-domain-and-flow-state/requirements.md`, `projects/ultraplan-go/docs/PRD.md`, `projects/ultraplan-go/docs/TRD.md`, `projects/ultraplan-go/docs/ARCHITECTURE.md`, `templates/sprint-index.md`

This document selects what must be read, distilled, reasoned through, or checked for this sprint. It does not make implementation decisions. All selections must come from the project index -- no items may be included that are not listed in the project index.

## Sprint Scope

- **Sprint Goal:** Model planning-stage sprint artifacts and durable `flow-state.json` state through `plan.md`, with runtime-free sprint status inspection and no implementation, smoke, review, or issue workflow stages.
- **Planned Output:** `internal/sprint` package documentation, domain model, discovery, artifact paths, flow-state persistence, filesystem store, status service, sprint CLI wiring, top-level CLI registration/help, sprint and command tests, and sprint review evidence.
- **Depends On:** Sprint 16 project domain and project index for project resolution and planning-root boundaries; Sprint 14 validation, diagnostics, and JSON stability for exit-code and diagnostic discipline; Sprint 12 durable run-loop state for atomic-write and strict-schema precedent; Sprint 9 runtime integration for preserving a runtime-free sprint status boundary; PRD, TRD, and Architecture docs for Phase 2 planning scope.
- **Non-Goals:** Generating, repairing, or content-validating planning artifacts beyond existence and flow-state status inspection; `sprint validate`, `sprint prompt`, or `sprint flow`; technical-handbook, area-reasoning, reasoning, or plan generation; implementation execution; smoke investigation; automated review; issue tracking; Git mutation; study service/model reuse; JSON status output unless already supported without scope expansion; hosted service, browser UI, TUI, API server, metrics exporter, or project-management workflow.

## Source Project Index

- `projects/ultraplan-go/project-index.md` -- authoritative source. Any file or item referenced below must appear there.

## Selected Contracts

Each contract applies as a flat whole to this sprint. All paths must appear in the project index's "Active Contract Pool" table.

| Contract | Why Selected |
| -------- | ------------ |
| Architecture | Governs `internal/sprint` ownership, allowed `sprint -> project/workspace` dependencies, product/platform separation, and the prohibition on global planning/workflow abstractions. |
| Errors | Applies to actionable sprint resolution diagnostics, invalid flow-state failures, cause preservation, and exit-code mapping. |
| Observability | Applies to truthful status output, validation diagnostics, state path visibility, and safe reporting of failed stages without invoking runtime behavior. |
| Security | Applies to workspace path safety, sprint slug validation, artifact paths that must not escape the sprint root, and redaction of unsafe diagnostics. |
| Testing | Applies to deterministic unit and command tests for discovery, state loading, atomic writes, status derivation, CLI output, and runtime-free behavior. |
| Documentation | Applies to `internal/sprint/doc.go`, command help, and reviewable sprint artifacts. |
| CLI Surface | Applies to `ultraplan sprint <project> <sprint> status`, help text, argument parsing, stdout/stderr separation, and exit codes. |
| Persistence And Migrations | Applies to strict, versioned `flow-state.json` loading and atomic writes that preserve the last valid state on failure. |

## Selected Evidence Reports

Copied from the project index's "Available Evidence Reports" table. These tell the technical handbook which reports to read - the project index is the authoritative source.

| Report | Path | Covers |
| ------ | ---- | ------ |
| `01-project-structure` | `studies/go-cli-study/reports/final/01-project-structure.md` | Project layout, cmd/internal/pkg, dependency direction, thin entrypoints |
| `02-command-architecture` | `studies/go-cli-study/reports/final/02-command-architecture.md` | Command routing, flags, help text, command organization, shell completion |
| `03-dependency-injection` | `studies/go-cli-study/reports/final/03-dependency-injection.md` | Dependency construction, seams, testability, constructor patterns |
| `05-error-handling` | `studies/go-cli-study/reports/final/05-error-handling.md` | Error wrapping, classification, user-facing diagnostics, sentinel errors |
| `06-io-abstraction` | `studies/go-cli-study/reports/final/06-io-abstraction.md` | Filesystem/stdin/stdout abstraction, test seams, interface design |
| `07-state-context` | `studies/go-cli-study/reports/final/07-state-context.md` | Context propagation, app state, cancellation |
| `09-terminal-ux` | `studies/go-cli-study/reports/final/09-terminal-ux.md` | Terminal output, progress indicators, human UX, color and formatting |
| `10-logging-observability` | `studies/go-cli-study/reports/final/10-logging-observability.md` | Logs, diagnostics, structured events, observability patterns |
| `11-testing-strategy` | `studies/go-cli-study/reports/final/11-testing-strategy.md` | Unit, integration, fixture, command-level tests, coverage strategy |
| `12-extensibility` | `studies/go-cli-study/reports/final/12-extensibility.md` | Extension points, plugin architecture, package boundaries, API design |
| `13-security` | `studies/go-cli-study/reports/final/13-security.md` | Path safety, secrets, command injection risks, sandboxing |
| `15-philosophy` | `studies/go-cli-study/reports/final/15-philosophy.md` | Cross-cutting design philosophy, tradeoffs, maintainability |

## Selected Reasoning Templates

All paths must appear in the project index's "Available Reasoning Templates" table.

| Template | Output Path | Why Selected |
| -------- | ----------- | ------------ |
| Architecture | `projects/ultraplan-go/sprints/17-artifact-domain-and-flow-state/reasoning/architecture.md` | Sprint requires explicit reasoning about `internal/sprint` ownership, dependency direction, runtime-free status behavior, flow-state persistence boundaries, and deferred workflow stages. |

## Prior Decisions To Carry Forward

All decision paths must appear in the project index's "Prior Decisions" table.

| Decision | Path | Constraint For This Sprint |
| -------- | ---- | -------------------------- |
| None yet -- this is the first project. | N/A | No prior-decision artifacts are cataloged in the project index; carry forward the sprint requirements to keep product behavior module-owned, keep CLI adapters thin, preserve cause chains, redact unsafe diagnostics, avoid speculative shared abstractions, and keep generated/user-editable artifacts as filesystem state. |

## Required Review Protocols

All paths must appear in the project index's "Review Protocols" table.

| Protocol | Path | Required Evidence |
| -------- | ---- | ----------------- |
| Architecture Review | `system/protocols/architecture-review-protocol.md` | Evidence that package ownership, dependency direction, `internal/sprint` boundaries, platform/product separation, and absence of global workflow abstractions match the sprint scope. |
| Sprint Review | `system/protocols/sprint-review-protocol.md` | Evidence that all required sprint outputs, acceptance criteria, verification commands, deviations, and non-goals are reviewed after implementation. |

## Excluded Context

| Context | Reason Excluded | Revisit If |
| ------- | --------------- | ---------- |
| Configuration contract | Sprint status does not introduce new config precedence, runtime configuration, provider/model selection, or config display behavior. | Sprint work changes command preflight configuration, effective config rendering, or runtime config mapping. |
| LLM Runtime contract | Sprint status inspection and flow-state refresh must not invoke agentwrap, OpenCode, prompt generation, network calls, or runtime health checks. | Later planning sprints implement `sprint prompt` or runtime-backed `sprint flow`. |
| LLM Evaluation / Cost / Safety contract | No runtime execution, usage metadata, cost tracking, or LLM evaluation behavior is in scope. | Planning artifact generation begins using a runtime. |
| Workflows contract | This sprint models planning stage state and status only; it does not implement workflow execution, retries, scheduling, or runtime-backed orchestration. | `sprint flow --to <stage>` or another workflow executor enters scope. |
| Performance contract | No scheduler, concurrency, large report processing, or performance tuning work is selected beyond deterministic direct-child sprint discovery. | Sprint status must meet a new explicit latency or scalability target. |
| `04-configuration-management` evidence report | Configuration loading, environment precedence, and runtime mapping are not changed by this sprint. | The sprint command starts accepting or resolving new config inputs. |
| `08-concurrency` evidence report | Sprint status and flow-state refresh are single-command inspection behavior with no worker pools or goroutine scheduling scope. | Future sprint flow execution introduces concurrent planning stages or runtime tasks. |
| `14-performance` evidence report | The sprint requires deterministic filesystem inspection, not broader startup, memory, or large-repository performance design. | Status inspection is expanded to large recursive scans or performance budgets. |
| Deep Smoke Sprint protocol | Requirements explicitly exclude smoke investigation and runtime invocation from this sprint. | A completed sprint later requires real-runtime smoke evidence. |
| `ultraplan-go-smoke` smoke harness | Runtime-backed smoke evidence is not needed for a runtime-free status and flow-state sprint. | Review requests deep smoke or a later runtime-backed planning sprint needs external evidence. |
| Deferred sprint planning commands | `sprint validate`, `sprint prompt`, and `sprint flow` are non-goals for this sprint. | Sprint 18 or later brings validation, prompt preview, or runtime-backed planning flow into scope. |
| Deferred execution, smoke, review, issues, and Git mutation behavior | PRD, TRD, project index, and sprint requirements all defer these beyond Phase 2 planning through `plan.md`. | A future requirements revision explicitly adds implementation execution, smoke investigation, review automation, issue tracking, or Git mutation. |
| Study service, source/dimension model, report validator, rating parser, summary generator, and run-loop scheduling reuse | Phase 2 planning context requires reuse of infrastructure only, not study-owned product semantics. | A later sprint proves a shared mechanical helper is needed without carrying study semantics. |

## Next Artifacts

- `technical-handbook.md` reads from the evidence reports listed above.
- `reasoning/*.md` captures area-specific reasoning.
- `reasoning.md` makes final sprint decisions.
- `plan.md` executes `reasoning.md`.
- `review.md` runs the selected review protocols against implementation.
