# Sprint Index: 09 Runtime Integration

> Project: `ultraplan-go`
> Sprint: `09-runtime-integration`
> Purpose: selected context for this sprint. Must be a subset of `projects/ultraplan-go/project-index.md`.
> **Inputs Used:** `projects/ultraplan-go/sprints/09-runtime-integration/requirements.md`, `projects/ultraplan-go/project-index.md`, `projects/ultraplan-go/docs/PRD.md`, `projects/ultraplan-go/docs/TRD.md`, `projects/ultraplan-go/docs/ARCHITECTURE.md`, `templates/sprint-index.md`

This document selects what must be read, distilled, reasoned through, or checked for this sprint. It does not make implementation decisions. All selections must come from the project index — no items may be included that are not listed in the project index.

## Sprint Scope

- **Sprint Goal:** Integrate UltraPlan's generic runtime boundary with `github.com/Antonio7098/agentwrap` and `agentwrap/opencode` so later study execution can use OpenCode through agentwrap wrappers, health checks, policies, validation, permissions, and events without UltraPlan reimplementing runtime supervision.
- **Planned Output:** Runtime integration sprint artifacts plus the platform runtime API, agentwrap/opencode composition, runtime health mapping, policy and permission mapping, event and diagnostics mapping, runtime config extensions, health command runtime wiring, and fake-runtime/config/health tests listed in the sprint requirements.
- **Depends On:** Sprints 1-8 carry forward CLI/app wiring, workspace/config/logging/health foundations, study metadata context without platform-runtime product imports, workspace and prompt locations, Markdown source applicability context, report validation hooks, prompt composition outputs, and run-state ID/task metadata placeholders. Runtime execution remained deferred in those prior outputs and enters scope only through the generic platform runtime boundary in this sprint.
- **Non-Goals:** Public study run, synthesize, run-all, run-loop, batch scheduling, durable workflow orchestration, per-study locks, stale-running recovery, active task cancellation, summary generation, code-reference extraction, target workflows, sprint planning, sprint execution, stable new public runtime JSON schemas, and any UltraPlan-owned OpenCode subprocess supervisor, native event decoder, retry engine, permission translator, or validation wrapper that duplicates agentwrap behavior.

## Source Project Index

- `projects/ultraplan-go/project-index.md` — authoritative source. Any file or item referenced below must appear there.

## Selected Contracts

Each contract applies as a flat whole to this sprint. All paths must appear in the project index's "Active Contract Pool" table.

| Contract | Why Selected |
| ------------ | ------------ |
| Architecture | Governs module ownership, dependency direction, thin entrypoints, and the required `study -> platform/runtime -> agentwrap/opencode` separation. |
| Errors | Governs wrapped errors, actionable diagnostics, exit behavior, and preserving agentwrap `SDKError` classifications without string matching. |
| Configuration | Governs runtime config precedence, validation, mapping into agentwrap/opencode options, command preflight behavior, and redaction. |
| Observability | Governs canonical events, safe diagnostics, structured logs, runtime metadata, health truthfulness, usage/cost signals, and run/task correlation fields. |
| Security | Governs secret redaction, source/workspace path safety, sandbox and permission posture, and avoiding direct unsafe runtime process handling. |
| Testing | Governs fake-runtime coverage, fixture/unit tests, gated real-runtime smoke tests, and default test runs that do not require OpenCode or credentials. |
| CLI Surface | Governs `ultraplan health` runtime reporting, help/exit-code behavior, and stable human/script-facing output for health failures. |
| LLM Runtime | Governs agentwrap/OpenCode runtime boundary, provider/model handling, health checks, adapter use, wrapper composition, validation, policy, permissions, and events. |
| LLM Evaluation / Cost / Safety | Governs runtime validation discipline, safe usage/cost metadata, unknown usage handling, and safety around native payloads and diagnostics. |
| Workflows | Governs context propagation, cancellation semantics, retry/fallback ownership through agentwrap, and operational signals while full batch/run-loop workflow scope remains excluded. |

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
| `08-concurrency` | `studies/go-cli-study/reports/final/08-concurrency.md` | Worker pools, cancellation, parallel execution, goroutine management |
| `10-logging-observability` | `studies/go-cli-study/reports/final/10-logging-observability.md` | Logs, diagnostics, structured events, observability patterns |
| `11-testing-strategy` | `studies/go-cli-study/reports/final/11-testing-strategy.md` | Unit, integration, fixture, command-level tests, coverage strategy |
| `12-extensibility` | `studies/go-cli-study/reports/final/12-extensibility.md` | Extension points, plugin architecture, package boundaries, API design |
| `13-security` | `studies/go-cli-study/reports/final/13-security.md` | Path safety, secrets, command injection risks, sandboxing |
| `15-philosophy` | `studies/go-cli-study/reports/final/15-philosophy.md` | Cross-cutting design philosophy, tradeoffs, maintainability |

## Selected Reasoning Templates

All paths must appear in the project index's "Available Reasoning Templates" table.

| Template | Output Path | Why Selected |
| ------------ | -------------------------------------------------------------- | ------------ |
| Architecture | `projects/ultraplan-go/sprints/09-runtime-integration/reasoning/runtime-integration.md` | Required to reason through runtime/product separation, wrapper order, config mapping, permission posture, event/log mapping, test seams, and deferred study orchestration. |

## Prior Decisions To Carry Forward

All decision paths must appear in the project index's "Prior Decisions" table.

The project index lists no formal prior decisions, so no prior-decision artifacts are selected here. Carry-forward constraints from completed prior sprint outputs are recorded in Sprint Scope as dependencies from this sprint's requirements rather than as project-index decision records.

## Required Review Protocols

All paths must appear in the project index's "Review Protocols" table.

| Protocol | Path | Required Evidence |
| ------------ | --------------------------------------- | ----------------- |
| Architecture Review | `system/protocols/architecture-review-protocol.md` | Evidence that runtime behavior remains generic in `internal/platform/runtime`, platform packages do not import product modules, and OpenCode execution is delegated through `agentwrap/opencode`. |
| Sprint Review | `system/protocols/sprint-review-protocol.md` | Evidence that required sprint artifacts and implementation outputs exist, acceptance criteria are checked, verification commands are recorded, and carry-forward decisions are captured. |

## Excluded Context

| Context | Reason Excluded | Revisit If |
| ----------- | --------------- | ------------- |
| Public study execution commands | This sprint builds the runtime boundary and health wiring only; public `study run`, `synthesize`, `run-all`, and `run-loop` execution remains outside scope. | A later execution sprint selects study task orchestration scope. |
| Full Batch / Durable Workflow Gate scope | Worker pools, run-loop mutation, locks, stale-running recovery, and durable run stores are explicitly non-goals except for metadata compatibility with existing run-state placeholders. | Batch/run-loop execution enters scope. |
| Study report semantics inside `internal/platform/runtime` | Platform runtime must stay generic and expose validation hooks without importing study types or owning report rules. | Study execution wires product validators through the generic runtime boundary. |
| Summary generation | Summary writing and score aggregation are not part of runtime integration. | Summary generation sprint or study batch execution requires it. |
| Code reference extraction | Citation parsing and snippet extraction are not part of runtime integration. | Code extraction implementation or validation scope is selected. |
| Target workflows, sprint planning, and sprint execution | Project scope and sprint requirements defer target and sprint workflows. | Project requirements are revised to include target/sprint behavior. |
| Stable new public runtime JSON schemas | Sprint requirements do not require new stable public JSON schemas for runtime status or health unless already present. | A public JSON output contract is selected and documented. |
| Direct OpenCode process supervision or native decoding by UltraPlan | Agentwrap/opencode owns process launch, stdout/stderr handling, native event projection, retry/fallback, permissions, health, validation, and diagnostics. | Agentwrap cannot represent a required runtime behavior and a reviewed exception is recorded. |
| `09-terminal-ux` evidence report | Health output must be actionable, but this sprint is primarily runtime integration rather than terminal UX design. | A sprint focuses on progress output, TUI behavior, or terminal formatting. |
| `14-performance` evidence report | Broad performance tuning, scheduler throughput, and large artifact handling are not central to this sprint's generic runtime boundary. | Batch scheduling, large-run status, or repository-scan performance enters scope. |
| Performance contract | Scheduler/concurrency/report-processing performance is not selected as a governing contract for this runtime integration slice. | Performance-sensitive batch scheduling or large artifact handling enters scope. |
| Persistence And Migrations contract | Durable run record migration and persisted schema compatibility are not selected beyond existing run-state metadata alignment required by the sprint. | Durable runtime records, migrations, or stable persisted runtime schemas enter scope. |

## Next Artifacts

- `technical-handbook.md` reads from the evidence reports listed above.
- `reasoning/*.md` captures area-specific reasoning.
- `reasoning.md` makes final sprint decisions.
- `plan.md` executes `reasoning.md`.
- `review.md` runs the selected review protocols against implementation.
