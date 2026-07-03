# Sprint Index: Study Initialization From YAML

> Project: `ultraplan-go`
> Sprint: `04-study-init-yaml`
> Purpose: selected context for this sprint. Must be a subset of `projects/ultraplan-go/project-index.md`.
> **Inputs Used:** `projects/ultraplan-go/project-index.md`, `projects/ultraplan-go/sprints/04-study-init-yaml/requirements.md`, `projects/ultraplan-go/docs/PRD.md`, `projects/ultraplan-go/docs/TRD.md`, `projects/ultraplan-go/docs/ARCHITECTURE.md`, `templates/sprint-index.md`

This document selects what must be read, distilled, reasoned through, or checked for this sprint. It does not make implementation decisions. All selections must come from the project index — no items may be included that are not listed in the project index.

## Sprint Scope

- **Sprint Goal:** Implement `ultraplan study init <study-init.yml>` so explicit YAML study definitions generate deterministic study directories, normalized YAML, editable dimension files, a study README, optional shallow source clones, and safe dry-run/force behavior.
- **Planned Output:** Study initialization service, YAML parser, artifact renderer, clone executor boundary, CLI command, unit tests, command tests, technical handbook, reasoning, and implementation plan for Sprint 4.
- **Depends On:** Sprint 1 command registration and CLI shell; Sprint 2 workspace discovery, path safety, and global `--workspace`; Sprint 3 study domain, listing, and resolution compatibility; PRD, TRD, and Architecture docs listed in the project index.
- **Non-Goals:** Runtime-assisted source or dimension completion, runtime suggestion prompts, suggestion cache artifacts, `--no-assist`, source URL verification, replacement source suggestions, repair prompts, OpenCode/agentwrap runtime execution, Markdown document source discovery, frontmatter parsing, `applicable_dimensions`, applicability filtering, prompt composition, analysis runs, synthesis, report validation, summary generation, run state, retries, cancellation, diagnostics persistence, code-reference extraction, target workflows, sprint workflows, and new stable public JSON schemas.

## Source Project Index

- `projects/ultraplan-go/project-index.md` — authoritative source. Any file or item referenced below must appear there.

## Selected Contracts

Each contract applies as a flat whole to this sprint. All paths must appear in the project index's "Active Contract Pool" table.

| Contract | Why Selected |
| ------------ | -------------------------------------------- |
| Architecture | Required for module ownership, dependency direction, thin CLI entrypoints, and keeping study initialization behavior in `internal/study` while CLI wiring stays in `internal/app`. |
| Errors | Required for actionable validation errors, path/source/dimension diagnostics, exit-code conventions, clone failure reporting, and partial completion behavior. |
| Security | Required for workspace path safety, preventing writes outside the selected root, avoiding shell interpolation for `git clone`, and ensuring generated artifacts do not contain secrets. |
| Testing | Required for unit and command tests covering YAML parsing, filesystem planning, dry-run, force behavior, clone fakes, and offline deterministic verification. |
| Documentation | Required for generated `README.md`, help output, normalized human-editable artifacts, and user-facing command guidance. |
| CLI Surface | Required for `ultraplan study init <study-init.yml>`, flags, help, output, usage errors, dry-run output, and exit behavior. |

## Selected Evidence Reports

Copied from the project index's "Available Evidence Reports" table. These tell the technical handbook which reports to read – the project index is the authoritative source.

| Report | Path | Covers |
|---|---|---|
| `01-project-structure` | `studies/go-cli-study/reports/final/01-project-structure.md` | Project layout, cmd/internal/pkg, dependency direction, thin entrypoints |
| `02-command-architecture` | `studies/go-cli-study/reports/final/02-command-architecture.md` | Command routing, flags, help text, command organization, shell completion |
| `03-dependency-injection` | `studies/go-cli-study/reports/final/03-dependency-injection.md` | Dependency construction, seams, testability, constructor patterns |
| `05-error-handling` | `studies/go-cli-study/reports/final/05-error-handling.md` | Error wrapping, classification, user-facing diagnostics, sentinel errors |
| `06-io-abstraction` | `studies/go-cli-study/reports/final/06-io-abstraction.md` | Filesystem/stdin/stdout abstraction, test seams, interface design |
| `09-terminal-ux` | `studies/go-cli-study/reports/final/09-terminal-ux.md` | Terminal output, progress indicators, human UX, color and formatting |
| `11-testing-strategy` | `studies/go-cli-study/reports/final/11-testing-strategy.md` | Unit, integration, fixture, command-level tests, coverage strategy |
| `13-security` | `studies/go-cli-study/reports/final/13-security.md` | Path safety, secrets, command injection risks, sandboxing |
| `15-philosophy` | `studies/go-cli-study/reports/final/15-philosophy.md` | Cross-cutting design philosophy, tradeoffs, maintainability |

## Selected Reasoning Templates

All paths must appear in the project index's "Available Reasoning Templates" table.

| Template | Output Path | Why Selected |
| ------------ | -------------------------------------------------------------- | ------------ |
| Architecture | `projects/ultraplan-go/sprints/04-study-init-yaml/reasoning/architecture.md` | Needed to reason through `internal/study` ownership, CLI/app boundaries, workspace dependency use, clone command seam, and deterministic filesystem artifact generation. |

## Prior Decisions To Carry Forward

All decision paths must appear in the project index's "Prior Decisions" table.

| Decision | Path | Constraint For This Sprint |
| ------------ | -------- | -------------------------- |
| None yet — this is the first project. | None listed in project index. | No prior decision artifacts are available to constrain this sprint; carry forward only the dependency notes listed in the sprint requirements. |

## Required Review Protocols

All paths must appear in the project index's "Review Protocols" table.

| Protocol | Path | Required Evidence |
| ------------ | --------------------------------------- | ----------------- |
| Architecture Review | `system/protocols/architecture-review-protocol.md` | Evidence that study initialization behavior is owned by `internal/study`, CLI handlers remain thin, workspace/path safety dependencies follow allowed direction, and platform packages do not import `internal/study`. |
| Sprint Review | `system/protocols/sprint-review-protocol.md` | Evidence that required outputs exist, acceptance criteria are checked, selected contracts and evidence were used, non-goals stayed excluded, and `go test ./...` plus `go build ./cmd/ultraplan` pass. |

## Excluded Context

| Context | Reason Excluded | Revisit If |
| ----------- | --------------- | ------------- |
| `04-configuration-management` evidence report | Sprint uses existing workspace/global flag behavior but does not introduce new configuration precedence, environment variable, or config schema behavior beyond prior sprints. | A later sprint changes config loading, command overrides, runtime config mapping, or `config show`. |
| `07-state-context` evidence report | Durable run state, cancellation propagation, and active task context are outside study initialization scope. | A run, run-all, run-loop, status, retry, or cancellation sprint begins. |
| `08-concurrency` evidence report | Study initialization performs bounded local planning, writes, and optional clone calls; worker pools and parallel runtime execution are not in scope. | Batch analysis, synthesis scheduling, or concurrent clone execution enters scope. |
| `10-logging-observability` evidence report | Persistent diagnostics, event streams, structured logs, and runtime observability are explicitly non-goals for this sprint. | Runtime execution, diagnostics persistence, status, or observability work enters scope. |
| `12-extensibility` evidence report | Plugin architecture and broader extension points are not needed for explicit YAML initialization. | New runtime adapters, public extension APIs, plugins, or cross-module extension seams are introduced. |
| `14-performance` evidence report | Large-repository scanning and runtime/report processing performance are outside scope; initialization must not recursively inspect cloned source contents. | Commands begin scanning large repositories, processing large reports, or scheduling many runtime tasks. |
| LLM Runtime contract | OpenCode/agentwrap runtime execution and assisted completion are excluded; clone execution is local `git`, not agent runtime work. | Runtime-assisted completion or study analysis execution enters scope. |
| LLM Evaluation / Cost / Safety contract | Token/cost metadata, evaluation discipline, and runtime safety gates are not part of explicit YAML initialization. | Runtime-backed suggestions, analysis, synthesis, repair, or evaluation enters scope. |
| Workflows contract | Stateful orchestration, retries, resumability, and workflow DAG behavior are excluded from initialization. | Run-loop, run-all, retry, resume, or durable task state implementation begins. |
| Performance contract | No selected Sprint 4 requirement depends on scheduler throughput, report processing scale, or startup performance beyond normal CLI expectations. | Large-scale scheduling, report processing, or performance-sensitive startup work enters scope. |
| Persistence And Migrations contract | Sprint writes deterministic generated artifacts but does not introduce durable schema migrations or run-state persistence. | Versioned state/config/report schemas or migrations are added. |
| Assisted study completion | Requirements explicitly defer runtime-assisted source and dimension completion, `--no-assist`, suggestion prompts, and suggestion cache artifacts. | A later sprint implements runtime suggestions for missing sources or dimensions. |
| Markdown document source discovery and applicability | Requirements explicitly defer Markdown document source discovery, frontmatter parsing, `applicable_dimensions`, and applicability filtering to Sprint 5. | Sprint 5 begins or the requirements are revised. |
| Runtime study execution and report generation | Prompt composition, analysis runs, synthesis, validation, summary generation, retries, cancellation, and diagnostics persistence are non-goals. | Study run, run-all, synthesize, validate, summary, status, or code extraction sprints begin. |
| Target and sprint workflows | Project scope and sprint requirements explicitly defer target and sprint CLI workflows. | Product requirements are revised to include target or sprint workflows. |

## Next Artifacts

- `technical-handbook.md` reads from the evidence reports listed above.
- `reasoning/*.md` captures area-specific reasoning.
- `reasoning.md` makes final sprint decisions.
- `plan.md` executes `reasoning.md`.
- `review.md` runs the selected review protocols against implementation.
