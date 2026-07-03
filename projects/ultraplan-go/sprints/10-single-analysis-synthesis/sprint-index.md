# Sprint Index: 10 Single Analysis and Synthesis

> Project: `ultraplan-go`
> Sprint: `10-single-analysis-synthesis`
> Purpose: selected context for this sprint. Must be a subset of `projects/ultraplan-go/project-index.md`.
> **Inputs Used:** `projects/ultraplan-go/project-index.md`, `projects/ultraplan-go/sprints/10-single-analysis-synthesis/requirements.md`, `projects/ultraplan-go/docs/ARCHITECTURE.md`, `projects/ultraplan-go/docs/PRD.md`, `projects/ultraplan-go/docs/TRD.md`, `templates/sprint-index.md`

This document selects what must be read, distilled, reasoned through, or checked for this sprint. It does not make implementation decisions. All selections must come from the project index - no items may be included that are not listed in the project index.

## Sprint Scope

- **Sprint Goal:** Implement `ultraplan study <study> run <dimension-ref> <source-ref>` and `ultraplan study <study> synthesize <dimension-ref>` so one analysis task or one synthesis task executes through the agentwrap-backed runtime and succeeds only after the expected report artifact exists and passes study-owned validation.
- **Planned Output:** Sprint artifacts plus study-owned single analysis orchestration, single synthesis orchestration, runtime dependency wiring, execution domain fields, thin CLI command wiring, app composition updates, fake-runtime study tests, and CLI execution tests listed in the sprint requirements.
- **Depends On:** Sprints 1-9 carry forward buildable CLI/app wiring, workspace/config/logging/health foundations, study/dimension/source resolution, study artifact layout, Markdown source discovery and applicability, report validators and rating parsing, deterministic prompt builders, task/state vocabulary compatibility, and agentwrap/OpenCode platform runtime integration. Runtime execution was deferred until this sprint, while batch/run-loop mutation, summary generation, code extraction, and target/sprint workflows remain deferred.
- **Non-Goals:** `study run-all`, task matrix execution, bounded worker pools, multi-task scheduling, `study run-loop`, resumable batch mutation, stale running recovery, retry scheduling across process restarts, cancellation persistence, per-study lock files, automatic synthesis after analysis, summary generation, code reference extraction, report repair workflows beyond configured runtime validation wrappers, durable agentwrap run-store persistence across process exits, stable public execution JSON schemas, target workflows, sprint planning, sprint execution, and any UltraPlan-owned OpenCode subprocess supervisor, retry engine, permission translator, event decoder, or validation wrapper duplicating agentwrap behavior.

## Source Project Index

- `projects/ultraplan-go/project-index.md` - authoritative source. Any file or item referenced below must appear there.

## Selected Contracts

Each contract applies as a flat whole to this sprint. All paths must appear in the project index's "Active Contract Pool" table.

| Contract | Why Selected |
| ------------ | ------------ |
| Architecture | Governs module ownership, dependency direction, thin entrypoints, and keeping study workflow behavior in `internal/study` while `internal/platform/runtime` remains generic. |
| Errors | Governs wrapped failures, actionable CLI diagnostics, exit-code mapping, validation failures, runtime failures, and preserving cause chains without parsing human-readable runtime strings. |
| Configuration | Governs runtime/provider/model/timeout resolution, command preflight validation, config redaction, and mapping execution settings into runtime requests. |
| Observability | Governs runtime events, task metadata, correlation fields, safe diagnostics, usage/cost fields when available, and truthful reporting of validation/runtime outcomes. |
| Security | Governs workspace and source path safety, secret redaction, runtime permission posture, document-source isolation, and avoiding unsafe direct process handling. |
| Testing | Governs fake-runtime tests, fixture/unit coverage, command tests, and default verification without OpenCode, provider credentials, network access, or real subprocess execution. |
| CLI Surface | Governs `study run` and `study synthesize` command shape, help text, argument errors, stdout/stderr separation, exit codes, and optional JSON behavior if added. |
| LLM Runtime | Governs agentwrap/OpenCode runtime boundary, provider/model handling, health/capability checks, request metadata, permission policy, validation hooks, event consumption, and final runtime result handling. |
| LLM Evaluation / Cost / Safety | Governs runtime validation discipline, safe handling of usage/cost metadata, native payload safety, and the principle that runtime success is insufficient without product artifact validation. |

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
| `12-extensibility` | `studies/go-cli-study/reports/final/12-extensibility.md` | Extension points, plugin architecture, package boundaries, API design |
| `13-security` | `studies/go-cli-study/reports/final/13-security.md` | Path safety, secrets, command injection risks, sandboxing |
| `15-philosophy` | `studies/go-cli-study/reports/final/15-philosophy.md` | Cross-cutting design philosophy, tradeoffs, maintainability |

## Selected Reasoning Templates

All paths must appear in the project index's "Available Reasoning Templates" table.

| Template | Output Path | Why Selected |
| ------------ | -------------------------------------------------------------- | ------------ |
| Architecture | `projects/ultraplan-go/sprints/10-single-analysis-synthesis/reasoning/single-analysis-synthesis.md` | Required to reason through study/runtime ownership, request construction, validation boundaries, Markdown skip behavior, CLI output boundaries, test seams, and deferred batch/run-loop scope. |

## Prior Decisions To Carry Forward

All decision paths must appear in the project index's "Prior Decisions" table.

The project index lists no formal prior decisions, so no prior-decision artifacts are selected here. Carry-forward constraints from completed prior sprint outputs are recorded in Sprint Scope as dependencies from this sprint's requirements rather than as project-index decision records.

## Required Review Protocols

All paths must appear in the project index's "Review Protocols" table.

| Protocol | Path | Required Evidence |
| ------------ | --------------------------------------- | ----------------- |
| Architecture Review | `system/protocols/architecture-review-protocol.md` | Evidence that study execution behavior lives in `internal/study`, CLI wiring stays thin in `internal/app`, `internal/platform/runtime` remains generic, and no product semantics move into platform runtime. |
| Sprint Review | `system/protocols/sprint-review-protocol.md` | Evidence that required sprint artifacts and implementation outputs exist, acceptance criteria are checked, fake-runtime and command tests cover required paths, verification commands are recorded, and carry-forward decisions are captured. |

## Excluded Context

| Context | Reason Excluded | Revisit If |
| ----------- | --------------- | ------------- |
| Full study run and task matrix execution | This sprint implements one selected analysis task and one selected synthesis task only; `study run-all`, worker pools, and multi-task scheduling are non-goals. | Batch execution enters scope. |
| Full Batch / Durable Workflow Gate scope | Run-loop mutation, stale running recovery, retry scheduling across process restarts, cancellation persistence, per-study locks, and durable workflow orchestration are explicitly deferred. | `study run-loop` or durable batch recovery enters scope. |
| Automatic synthesis after analysis | `study run` and `study synthesize` remain separate commands for this sprint. | Requirements select automatic post-analysis synthesis. |
| Summary generation | `summary.csv` updates and score aggregation are explicit non-goals. | Summary generation or batch completion enters scope. |
| Code reference extraction | Citation parsing and snippet extraction are not part of single analysis or synthesis execution. | Code extraction implementation or validation scope is selected. |
| Target workflows, sprint planning, and sprint execution | Project scope and sprint requirements defer target and sprint workflows. | Project requirements are revised to include target/sprint behavior. |
| Report repair product workflow | Repair prompts and bounded repair loops are not required unless already supplied by configured runtime validation wrappers without new product workflow logic. | Requirements select explicit report repair behavior. |
| Durable agentwrap run-store persistence across process exits | This sprint executes single tasks and does not require durable runtime run-store persistence. | Durable run inspection or resumable runtime records enter scope. |
| Stable public execution JSON schemas | New stable JSON schemas are not required unless explicitly added and documented by this sprint. | Public automation contracts are selected for execution commands. |
| Direct OpenCode subprocess supervision or native decoding by UltraPlan | Agentwrap/opencode owns process launch, native event projection, health checks, permissions, validation wrappers, retry/fallback, and diagnostics. | Agentwrap cannot represent a required runtime behavior and a reviewed exception is recorded. |
| Study report semantics inside `internal/platform/runtime` | Platform runtime must stay generic and must not import study types, enforce synthesis gating, or own report validation rules. | A generic runtime API limitation requires a reviewed architecture exception. |
| `14-performance` evidence report | Broad performance tuning, scheduler throughput, and large artifact performance are not central to this single-task execution sprint. | Batch scheduling, large-run status, or repository-scan performance enters scope. |
| Performance contract | Scheduler/concurrency/report-processing performance is not selected as a governing contract for this single-task execution slice. | Performance-sensitive batch scheduling or large artifact handling enters scope. |
| Persistence And Migrations contract | Durable schema migrations and persisted run-store compatibility are not selected beyond compatibility with existing task IDs and validation summaries. | Durable runtime records, migrations, or stable persisted execution schemas enter scope. |
| Workflows contract | Full stateful workflow behavior is deferred; this sprint uses runtime/event handling expectations through LLM Runtime and Observability instead of selecting batch workflow governance. | Run-loop, orchestration retries, or resumable workflow execution enters scope. |

## Next Artifacts

- `technical-handbook.md` reads from the evidence reports listed above.
- `reasoning/*.md` captures area-specific reasoning.
- `reasoning.md` makes final sprint decisions.
- `plan.md` executes `reasoning.md`.
- `review.md` runs the selected review protocols against implementation.
