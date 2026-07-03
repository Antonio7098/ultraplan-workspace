# Sprint Index: 11 Run All Batch Execution

> Project: `ultraplan-go`
> Sprint: `11-run-all-batch-execution`
> Purpose: selected context for this sprint. Must be a subset of `projects/ultraplan-go/project-index.md`.
> **Inputs Used:** `projects/ultraplan-go/sprints/11-run-all-batch-execution/requirements.md`, `projects/ultraplan-go/project-index.md`, `projects/ultraplan-go/docs/PRD.md`, `projects/ultraplan-go/docs/TRD.md`, `projects/ultraplan-go/docs/ARCHITECTURE.md`, `templates/sprint-index.md`

This document selects what must be read, distilled, reasoned through, or checked for this sprint. It does not make implementation decisions. All selections must come from the project index — no items may be included that are not listed in the project index.

## Sprint Scope

- **Sprint Goal:** Implement `ultraplan study <study> run-all` so a study can execute all selected applicable analysis tasks with bounded parallelism, validate outputs, synthesize completed dimensions, and generate a deterministic summary without introducing durable `run-loop` orchestration.
- **Planned Output:** Sprint artifacts plus study-owned batch execution, summary generation, service wiring, thin CLI command wiring, and fake-runtime/unit tests for batch execution, synthesis gating, summary semantics, cancellation, and command behavior.
- **Depends On:** Sprint 5 Markdown source discovery and applicability; Sprint 6 report validation and rating parsing; Sprint 7 prompt composition; Sprint 8 run-state/status primitives as optional truthful-status support only; Sprint 9 agentwrap/OpenCode runtime integration; Sprint 10 single analysis and synthesis execution semantics; project docs PRD, TRD, and ARCHITECTURE.
- **Non-Goals:** No `run-loop` command, production-grade resumable orchestration, durable retry/backoff scheduling outside agentwrap policy use, stale-running recovery, per-study lock files, force unlock, stable public JSON output unless already required by existing conventions, code-reference extraction, target/sprint CLI commands, new runtime adapters, direct OpenCode parsing, default real OpenCode smoke tests, generic workflow engine, DAG authoring system, plugin system, hosted service, browser UI, or multi-user collaboration features.

## Source Project Index

- `projects/ultraplan-go/project-index.md` — authoritative source. Any file or item referenced below must appear there.

## Selected Contracts

Each contract applies as a flat whole to this sprint. All paths must appear in the project index's "Active Contract Pool" table.

| Contract | Why Selected |
| ------------ | -------------------------------------------- |
| Architecture | Batch behavior must stay study-owned with thin CLI wiring, no product imports from `internal/platform/runtime`, and no global scheduler, reports, summary, prompts, or validation packages. |
| Errors | The command needs usage, validation, runtime, cancellation, and partial-completion behavior with safe diagnostics and preserved cause chains. |
| Configuration | Parallelism, runtime/model settings, timeout, permission policy, health requirements, and command overrides must resolve through existing configuration rules. |
| Observability | Runtime events must be consumed, task/run diagnostics must remain truthful, and human output must report completed, failed, skipped, pending, warnings, and artifact paths. |
| Security | Runtime requests, workspace paths, source isolation, Markdown document handling, redaction, and safe user output are in scope. |
| Testing | Normal verification must use fake runtimes and local fixtures; `go test ./...` must not require OpenCode, provider credentials, network access, or real subprocess execution. |
| Documentation | CLI help output and sprint artifacts must accurately document the new `run-all` behavior and exclusions. |
| CLI Surface | `ultraplan study <study> run-all` needs command shape, flags, help, argument validation, exit codes, and deterministic text output. |
| LLM Runtime | Analysis and synthesis tasks must use the existing agentwrap/OpenCode runtime boundary and runtime request mapping. |
| LLM Evaluation / Cost / Safety | Runtime execution must preserve safe metadata and avoid leaking prompts, embedded Markdown bodies, secrets, sensitive environment values, or unsafe native payloads. |
| Workflows | Bounded batch scheduling, cancellation, task status, synthesis gating, and partial completion are in scope, while durable `run-loop` resumability remains excluded. |
| Performance | Worker concurrency must be bounded, task matrix creation deterministic, event consumption non-blocking, and summary/report processing stable for larger studies. |
| Persistence And Migrations | `summary.csv` and any optional state updates must use workspace-contained explicit paths and safe file-write behavior without introducing Sprint 12 durable orchestration. |

## Selected Evidence Reports

Copied from the project index's "Available Evidence Reports" table. These tell the technical handbook which reports to read – the project index is the authoritative source.

| Report | Path | Covers |
|---|---|---|
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
| Architecture | `projects/ultraplan-go/sprints/11-run-all-batch-execution/reasoning/run-all-batch-execution.md` | Reason through study-owned batch orchestration, runtime/product separation, bounded workers, synthesis gating, summary placement, cancellation, and exclusions without making implementation decisions in this index. |

## Prior Decisions To Carry Forward

All decision paths must appear in the project index's "Prior Decisions" table.

| Decision | Path | Constraint For This Sprint |
| ------------ | -------- | -------------------------- |
| None | N/A | The project index lists no prior decisions; carry forward the PRD/TRD/ARCHITECTURE constraints and prior sprint dependencies named in the sprint requirements instead. |

## Required Review Protocols

All paths must appear in the project index's "Review Protocols" table.

| Protocol | Path | Required Evidence |
| ------------ | --------------------------------------- | ----------------- |
| Architecture Review | `system/protocols/architecture-review-protocol.md` | Verify study-owned batch behavior, thin CLI wiring, allowed dependency direction, generic runtime boundary, and absence of global technical-layer packages for batch behavior. |
| Sprint Review | `system/protocols/sprint-review-protocol.md` | Verify required sprint artifacts, acceptance criteria evidence, command behavior, fake-runtime tests, summary determinism, safe diagnostics, exclusions, `go test ./...`, and `go build ./cmd/ultraplan`. |

## Excluded Context

| Context | Reason Excluded | Revisit If |
| ----------- | --------------- | ------------- |
| `ultraplan study <study> run-loop` and production-grade resumable orchestration | Explicit sprint non-goal; Sprint 11 implements bounded `run-all` batch execution only and must not replace Sprint 12 durable orchestration. | Sprint 12 or a revised requirement brings durable run-loop behavior into scope. |
| Durable retry/backoff scheduling outside agentwrap policy use | Explicit sprint non-goal; batch tasks may use existing agentwrap policy behavior but must not create a separate retry scheduler. | Requirements add product-owned retry orchestration beyond existing agentwrap policy. |
| Stale-running recovery, per-study lock files, force unlock, and multi-process coordination | Explicit sprint non-goal; these are durable orchestration and coordination concerns rather than this sprint's bounded batch command. | Durable run-loop or multi-process safety enters scope. |
| Stable public JSON output for `run-all`, `status`, or summary commands | Explicit sprint non-goal unless existing app conventions already require it; this sprint focuses on deterministic human output. | Existing CLI conventions require JSON for this command or product requirements are revised. |
| Code-reference extraction and `ultraplan code <report>...` | Explicit sprint non-goal and separate product capability listed in project scope. | A code extraction sprint selects the related context. |
| Target scaffolding, sprint planning, sprint execution, and target/sprint CLI commands | Project and sprint non-goals defer target/sprint workflows from the current study-side build. | PRD/TRD are revised to include target/sprint workflows. |
| New runtime adapters, direct OpenCode parsing, or bypassing `agentwrap/opencode` | Explicit sprint non-goal and TRD constraint; runtime execution must use the existing platform runtime abstraction backed by agentwrap. | A runtime-adapter sprint explicitly selects that work. |
| Default real OpenCode smoke tests | Explicit sprint non-goal; normal tests must be deterministic, offline, and fake-runtime based. | A gated integration/smoke-test task is explicitly requested outside default verification. |
| Generic workflow engine, DAG authoring system, plugin system, hosted service, browser UI, and multi-user collaboration features | Project and sprint non-goals; these are broader product horizons unrelated to `run-all` batch execution. | Product scope changes beyond the local study CLI. |
| `12-extensibility` evidence report | Plugin architecture and extension-point design are not central to this sprint's concrete batch execution scope. | Reasoning identifies a concrete extension-point or API boundary decision that cannot be resolved from architecture and dependency-injection evidence. |

## Next Artifacts

- `technical-handbook.md` reads from the evidence reports listed above.
- `reasoning/*.md` captures area-specific reasoning.
- `reasoning.md` makes final sprint decisions.
- `plan.md` executes `reasoning.md`.
- `review.md` runs the selected review protocols against implementation.
