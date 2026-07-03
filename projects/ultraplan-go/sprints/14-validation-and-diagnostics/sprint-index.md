# Sprint Index: Validation Command, Diagnostics, and JSON Stability

> Project: `ultraplan-go`
> Sprint: `14-validation-and-diagnostics`
> Purpose: selected context for this sprint. Must be a subset of `projects/ultraplan-go/project-index.md`.
> **Inputs Used:** `projects/ultraplan-go/sprints/14-validation-and-diagnostics/requirements.md`, `projects/ultraplan-go/project-index.md`, `projects/ultraplan-go/docs/PRD.md`, `projects/ultraplan-go/docs/TRD.md`, `projects/ultraplan-go/docs/ARCHITECTURE.md`, `templates/sprint-index.md`

This document selects what must be read, distilled, reasoned through, or checked for this sprint. It does not make implementation decisions. All selections must come from the project index -- no items may be included that are not listed in the project index.

## Sprint Scope

- **Sprint Goal:** Make UltraPlan's study validation, health, and run status surfaces inspectable, automatable, and safe through a study validation command, stable JSON output shapes, and redacted diagnostics.
- **Planned Output:** `ultraplan study <study> validate` with text and JSON output, `ultraplan study <study> status --json`, `ultraplan health --json`, centralized JSON rendering, redacted diagnostics, validation/status/health tests, updated help text, and sprint review evidence.
- **Depends On:** Sprint 05 Markdown source discovery and applicability; Sprint 06 report validation and rating parsing; Sprint 08 run state persistence and status; Sprint 09 runtime integration and health; Sprint 10 single analysis and synthesis; Sprint 11 run-all batch execution; Sprint 12 durable run-loop; Sprint 13 summary generation and code-reference extraction; PRD/TRD JSON and diagnostics requirements.
- **Non-Goals:** Target, sprint planning, sprint execution, or sprint validation commands; browser UI, TUI dashboard, hosted service, metrics exporter, alerting integration, or local API server; new runtime adapters or OpenCode supervision changes outside the agentwrap-backed platform runtime boundary; analyses, synthesis, repairs, retries, fallback execution, source cloning, or prompt generation from validation; schema migrations beyond clear rejection or diagnostics for unsupported durable schema versions; mandatory real OpenCode/provider smoke tests in the default suite; report template, prompt semantic, summary scoring, or code-reference extraction behavior changes except where required for validation diagnostics; unsafe raw runtime payload or full native stderr persistence; a general diagnostics framework unrelated to `validate`, `health --json`, and `status --json`.

## Source Project Index

- `projects/ultraplan-go/project-index.md` -- authoritative source. Any file or item referenced below must appear there.

## Selected Contracts

Each contract applies as a flat whole to this sprint. All paths must appear in the project index's "Active Contract Pool" table.

| Contract | Why Selected |
| ------------ | -------------------------------------------- |
| Architecture | Required for module ownership: study validation remains in `internal/study`, command parsing/rendering remains in `internal/app`, and platform runtime remains product-agnostic. |
| Errors | Required for stable exit-code mapping, actionable validation failures, safe observed details, and command diagnostics. |
| Configuration | Required for health/config diagnostics, config redaction, command preflight behavior, and runtime health summaries. |
| Observability | Required for status, health, run/task metadata, validation summaries, retry/fallback metadata, and truthful diagnostics. |
| Security | Required for path safety, redaction of secrets and unsafe runtime payloads, and safe handling of user-visible diagnostics. |
| Testing | Required for deterministic offline tests, fake-first command tests, JSON schema assertions, race verification, and build verification. |
| Documentation | Required for updated command help and maintainable user-facing command surfaces. |
| CLI Surface | Required for command shape, `--json` flags, help text, text/JSON output behavior, and exit codes. |
| LLM Runtime | Required because health JSON and status diagnostics must use the existing agentwrap/OpenCode runtime boundary without adding product semantics to runtime adapters. |
| LLM Evaluation / Cost / Safety | Required for preserving unknown usage/cost as unknown rather than zero and for safe runtime metadata handling. |
| Workflows | Required for run-loop status, task state summaries, retry/fallback metadata, cancellation state, and runtime-free status behavior. |
| Performance | Required to keep validation bounded to known workspace/study artifacts and avoid recursive source repository scans. |
| Persistence And Migrations | Required for strict durable `run-state.json` schema handling, lock diagnostics, state compatibility diagnostics, and no silent migrations. |

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
| Architecture | `projects/ultraplan-go/sprints/14-validation-and-diagnostics/reasoning/architecture.md` | Needed to reason through study/app/platform ownership, runtime product separation, JSON rendering placement, durable state handling, and avoidance of global validation/diagnostics packages. |

## Prior Decisions To Carry Forward

All decision paths must appear in the project index's "Prior Decisions" table.

| Decision | Path | Constraint For This Sprint |
| ------------ | -------- | -------------------------- |
| No prior decisions listed in project index | N/A | Project index lists no formal prior decision artifacts; carry-forward constraints are taken from this sprint's requirements instead. |
| Preserve cause chains | `projects/ultraplan-go/sprints/14-validation-and-diagnostics/requirements.md` | Diagnostics and errors must retain causality while rendering safe, redacted user-facing details. |
| Runtime success is insufficient without artifact validation | `projects/ultraplan-go/sprints/14-validation-and-diagnostics/requirements.md` | Validation and status must not treat successful runtime execution as sufficient when required artifacts are missing or invalid. |
| Keep status runtime-free | `projects/ultraplan-go/sprints/14-validation-and-diagnostics/requirements.md` | `study <study> status --json` must inspect durable/local state and artifacts without invoking agentwrap, OpenCode, network, or subprocess paths. |
| Keep diagnostics redacted | `projects/ultraplan-go/sprints/14-validation-and-diagnostics/requirements.md` | Validation, health, and status text/JSON outputs must use the established redaction policy before rendering. |
| Address deferred stable JSON status surface | `projects/ultraplan-go/sprints/14-validation-and-diagnostics/requirements.md` | `study <study> status --json` must provide a stable, versioned machine-readable shape with counts, task summaries, lock diagnostics, run metadata summaries, and validation summaries. |

## Required Review Protocols

All paths must appear in the project index's "Review Protocols" table.

| Protocol | Path | Required Evidence |
| ------------ | --------------------------------------- | ----------------- |
| Architecture Review | `system/protocols/architecture-review-protocol.md` | Evidence that study validation is owned by `internal/study`, command wiring/rendering stays in `internal/app`, platform runtime remains product-agnostic, no forbidden global packages are introduced, and dependency direction is preserved. |
| Sprint Review | `system/protocols/sprint-review-protocol.md` | Evidence for implemented outputs, deviations, validation/status/health JSON schemas, redaction tests, runtime-free validate/status behavior, health runtime-boundary behavior, deterministic offline tests, race tests, build status, and review conclusions in `review.md`. |

## Excluded Context

| Context | Reason Excluded | Revisit If |
| ----------- | --------------- | ------------- |
| `08-concurrency` evidence report | This sprint validates and renders status/health/diagnostics; it does not introduce new scheduling, worker-pool, or runtime concurrency behavior beyond existing run-loop metadata inspection. | Implementation changes task scheduling, goroutine ownership, cancellation mechanics, or introduces new concurrent validation execution. |
| `12-extensibility` evidence report | Plugin architecture and broad extension points are not needed for the selected validation, status JSON, health JSON, and redaction surfaces. | The sprint introduces new public extension seams, plugin APIs, or reusable cross-module validation frameworks. |
| Target and sprint workflow requirements | The PRD/TRD and sprint requirements explicitly defer target scaffolding, sprint planning, sprint execution, and sprint validation commands. | The project scope is revised to include target or sprint-side product workflows. |
| Runtime adapter implementation details beyond health/status diagnostics | The sprint must not add new runtime adapters or change OpenCode process supervision outside the existing agentwrap-backed boundary. | Health JSON requires a capability that cannot be exposed through the existing platform runtime/agentwrap boundary. |
| Prompt generation, analysis execution, synthesis execution, repair, retry, fallback execution, and source cloning flows | The validation command must not invoke runtime work or mutate/generate study outputs, and status must remain runtime-free. | A later sprint explicitly adds repair or execution behavior triggered by validation results. |
| Browser UI, TUI dashboard, hosted service, metrics exporter, alerting integration, and local API server | These are explicit sprint non-goals and outside the current CLI-only automation surfaces. | Product requirements change to add non-CLI operational surfaces. |
| Schema migrations | This sprint only diagnoses or rejects unsupported durable schema versions and must not silently migrate state. | A later migration sprint is planned and added to the project index. |
| Full raw runtime payloads, full native stderr, prompt bodies, and embedded Markdown document bodies in output | User-visible diagnostics must be safe, redacted, and bounded; unsafe raw bytes and sensitive or large bodies are explicitly excluded. | An explicit debug-retention feature with safe controls is added to requirements. |

## Next Artifacts

- `technical-handbook.md` reads from the evidence reports listed above.
- `reasoning/*.md` captures area-specific reasoning.
- `reasoning.md` makes final sprint decisions.
- `plan.md` executes `reasoning.md`.
- `review.md` runs the selected review protocols against implementation.
