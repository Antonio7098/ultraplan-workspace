# Sprint Index: 08 Run State Persistence

> Project: `ultraplan-go`
> Sprint: `08-run-state-persistence`
> Purpose: selected context for this sprint. Must be a subset of `projects/ultraplan-go/project-index.md`.
> **Inputs Used:** `projects/ultraplan-go/project-index.md`, `projects/ultraplan-go/sprints/08-run-state-persistence/requirements.md`, `projects/ultraplan-go/docs/ARCHITECTURE.md`, `projects/ultraplan-go/docs/PRD.md`, `projects/ultraplan-go/docs/TRD.md`, `templates/sprint-index.md`

This document selects what must be read, distilled, reasoned through, or checked for this sprint. It does not make implementation decisions. All selections must come from the project index — no items may be included that are not listed in the project index.

## Sprint Scope

- **Sprint Goal:** Implement study-owned durable run-state persistence for planned analysis and synthesis tasks, including atomic load/save behavior and deterministic status summaries, without executing runtime work.
- **Planned Output:** A sprint handbook, run-state persistence reasoning, implementation plan, review, and Go implementation updates for study-owned run-state schema, construction, persistence, resume validation, status summaries, and `ultraplan study <study> status` human output.
- **Depends On:** Prior study-side foundations for CLI composition, workspace/config paths, study/source/dimension resolution, study initialization layout, Markdown source applicability, report validation, and prompt/output path concepts; product and technical context from Product Requirements, Technical Requirements, and Architecture source documents.
- **Non-Goals:** Agentwrap/OpenCode runtime integration, actual analysis or synthesis execution, worker-pool orchestration, retry/fallback/backoff execution, runtime health gates, event sinks, durable agentwrap run stores, user interrupt handling, active cancellation, per-study lock files, force-unlock behavior, stable public JSON status output, summary CSV writing, code reference extraction, target workflows, sprint planning, and sprint execution.

## Source Project Index

- `projects/ultraplan-go/project-index.md` — authoritative source. Any file or item referenced below must appear there.

## Selected Contracts

Each contract applies as a flat whole to this sprint. All paths must appear in the project index's "Active Contract Pool" table.

| Contract | Why Selected |
| ------------ | -------------------------------------------- |
| Architecture | Governs keeping durable run-state behavior in `internal/study`, CLI rendering in `internal/app`, and platform packages free of product imports. |
| Errors | Governs actionable diagnostics for missing, malformed, unsupported, or invalid run-state files and validation failures. |
| Configuration | Governs safe config summaries persisted in state and redaction of sensitive runtime/provider values. |
| Observability | Governs deterministic status diagnostics, run/task metadata, and truthful reporting of current or recent workflow state. |
| Security | Governs workspace-relative paths, secret-free persisted state, source isolation boundaries, and safe local filesystem behavior. |
| Testing | Governs unit, fixture, command-level, persistence, and scenario coverage required for durable state behavior. |
| CLI Surface | Governs `ultraplan study <study> status` command shape, human output, diagnostics, and exit behavior. |
| Workflows | Governs durable task state, explicit terminal states, resumability boundaries, diagnostics, and workflow status semantics while runtime execution remains deferred. |
| Performance | Governs deterministic status loading and summary calculation without unnecessary repository scans or runtime inspection. |
| Persistence And Migrations | Governs schema versioning, unsupported-version rejection, atomic writes, durable file formats, and migration-aware diagnostics. |

## Selected Evidence Reports

Copied from the project index's "Available Evidence Reports" table. These tell the technical handbook which reports to read – the project index is the authoritative source.

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
| Architecture | `projects/ultraplan-go/sprints/08-run-state-persistence/reasoning/run-state-persistence.md` | Reason through study-owned state ownership, package boundaries, persistence placement, status command separation, runtime/product separation, and dependency direction. |

## Prior Decisions To Carry Forward

All decision paths must appear in the project index's "Prior Decisions" table.

| Decision | Path | Constraint For This Sprint |
| ------------ | -------- | -------------------------- |
| None listed in project index | N/A | No selectable prior-decision rows are available in the project index; sprint requirements identify carry-forward dependencies, but they are not listed in the project index's Prior Decisions table. |

## Required Review Protocols

All paths must appear in the project index's "Review Protocols" table.

| Protocol | Path | Required Evidence |
| ------------ | --------------------------------------- | ----------------- |
| Architecture Review | `system/protocols/architecture-review-protocol.md` | Evidence that run-state behavior remains study-owned, CLI status rendering stays in `internal/app`, platform packages do not import product modules, and no global state/scheduler/report package is introduced. |
| Sprint Review | `system/protocols/sprint-review-protocol.md` | Evidence that required outputs exist, acceptance criteria are checked, implementation scope excludes runtime execution, and `go test ./...` plus `go build ./cmd/ultraplan` results are recorded. |

## Excluded Context

| Context | Reason Excluded | Revisit If |
| ----------- | --------------- | ------------- |
| `LLM Runtime` contract | Agentwrap/OpenCode runtime integration and runtime execution are explicit sprint non-goals; this sprint may persist placeholders only where required by state schema. | A later sprint implements analysis/synthesis runtime execution, runtime health gates, or agentwrap run-store integration. |
| `LLM Evaluation / Cost / Safety` contract | Cost, safety, evaluation, and provider metadata behavior depend on runtime execution, which is deferred. | Runtime execution, fallback, cost accounting, or evaluation gates enter sprint scope. |
| `Documentation` contract | This sprint creates internal planning artifacts and deterministic human status output, not user-facing generated docs or release documentation. | User-facing help, generated docs, recovery documentation, or release docs enter sprint scope. |
| `08-concurrency` evidence report | Worker pools, cancellation, and goroutine management are not implemented in this sprint beyond persisted status values and safe stale-state handling. | Run-loop orchestration, bounded workers, cancellation, or concurrent runtime execution enters sprint scope. |
| `12-extensibility` evidence report | Plugin and extension boundaries are not needed for the minimal study-owned persistence implementation. | Public extension points, plugins, or reusable persistence APIs enter scope. |
| Runtime execution workflows | Actual analysis/synthesis task execution, OpenCode invocation, provider credentials, subprocesses, worker pools, retries, fallback, repair, and backoff are non-goals. | A run-all or run-loop execution sprint begins. |
| Summary generation and `summary.csv` | Sprint status summaries are derived from persisted state; CSV summary generation remains out of scope. | A summary-generation sprint begins. |
| Code reference extraction | Citation parsing and snippet extraction are unrelated to run-state persistence and are explicitly excluded. | Code extraction commands or validation enter scope. |
| Target and sprint workflows | Product scope defers target scaffolding, sprint planning, and sprint execution. | The PRD/TRD are revised to include target or sprint-side implementation. |
| Per-study locks and force unlock | Lock files and force-unlock behavior are explicit sprint non-goals. | Concurrent mutation prevention is added to a later run-loop sprint. |
| Prior sprint review files as selectable decisions | Sprint requirements mention prior sprint review carry-forwards, but the project index's Prior Decisions table contains no selectable decision paths. | The project index adds those review files or decisions to its Prior Decisions table. |

## Next Artifacts

- `technical-handbook.md` reads from the evidence reports listed above.
- `reasoning/*.md` captures area-specific reasoning.
- `reasoning.md` makes final sprint decisions.
- `plan.md` executes `reasoning.md`.
- `review.md` runs the selected review protocols against implementation.
