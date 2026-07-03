# Sprint Index: 06 Study Run Execution

> Project: `ultraplan-go`
> Sprint: `06-study-run-execution`
> Purpose: selected context for this sprint. Must be a subset of `projects/ultraplan-go/project-index.md`.
> **Inputs Used:** `projects/ultraplan-go/project-index.md`, `projects/ultraplan-go/sprints/06-study-run-execution/requirements.md`, `projects/ultraplan-go/docs/PRD.md`, `projects/ultraplan-go/docs/TRD.md`, `projects/ultraplan-go/docs/ARCHITECTURE.md`, `templates/sprint-index.md`

This document selects what must be read, distilled, reasoned through, or checked for this sprint. It does not make implementation decisions. All selections must come from the project index — no items may be included that are not listed in the project index.

## Sprint Scope

- **Sprint Goal:** Make study reports first-class validated artifacts by implementing per-source report validation, final report validation, rating parsing, and actionable validation diagnostics without introducing runtime execution.
- **Planned Output:** Study-owned report validation domain updates, deterministic report path helpers, per-source and final report validators, rating parser, validation diagnostics, unit tests, and sprint artifacts for technical guidance, reasoning, plan, and review.
- **Depends On:** Prior study-side sprints for CLI/package structure, workspace path resolution, study/source/dimension domain, study initialization layout, and Markdown document source applicability behavior.
- **Non-Goals:** Runtime execution, prompt composition, agentwrap/OpenCode wiring, runtime health checks, policy runner wiring, retries, fallback, event mapping, run-all, run-loop, status, cancellation, durable task state, bounded worker pools, summary generation, `summary.csv` writing, code reference extraction, public stable validate command, target workflows, sprint planning, and sprint execution.

## Source Project Index

- `projects/ultraplan-go/project-index.md` — authoritative source. Any file or item referenced below must appear there.

## Selected Contracts

Each contract applies as a flat whole to this sprint. All paths must appear in the project index's "Active Contract Pool" table.

| Contract | Why Selected |
| -------- | ------------ |
| Architecture | Governs keeping study-owned validation, report paths, rating parsing, and diagnostics in `internal/study`; dependency direction and product/platform separation are especially important. |
| Errors | Governs actionable diagnostics, error wrapping, cause preservation, filesystem/parsing failure reporting, and validation failure classification. |
| Security | Governs workspace path safety, source isolation assumptions, safe diagnostics, and avoiding secret leakage in validation output. |
| Testing | Governs deterministic unit and fixture tests for validators, rating parsing, Markdown-source citation exemptions, missing files, missing sections, and diagnostics. |
| Performance | Governs report processing discipline, avoiding unnecessary repository scans, and keeping validation deterministic and bounded for later batch workflows. |

## Selected Evidence Reports

Copied from the project index's "Available Evidence Reports" table. These tell the technical handbook which reports to read – the project index is the authoritative source.

| Report | Path | Covers |
| ------ | ---- | ------ |
| `01-project-structure` | `studies/go-cli-study/reports/final/01-project-structure.md` | Project layout, cmd/internal/pkg, dependency direction, thin entrypoints |
| `05-error-handling` | `studies/go-cli-study/reports/final/05-error-handling.md` | Error wrapping, classification, user-facing diagnostics, sentinel errors |
| `06-io-abstraction` | `studies/go-cli-study/reports/final/06-io-abstraction.md` | Filesystem/stdin/stdout abstraction, test seams, interface design |
| `10-logging-observability` | `studies/go-cli-study/reports/final/10-logging-observability.md` | Logs, diagnostics, structured events, observability patterns |
| `11-testing-strategy` | `studies/go-cli-study/reports/final/11-testing-strategy.md` | Unit, integration, fixture, command-level tests, coverage strategy |
| `13-security` | `studies/go-cli-study/reports/final/13-security.md` | Path safety, secrets, command injection risks, sandboxing |
| `14-performance` | `studies/go-cli-study/reports/final/14-performance.md` | Startup latency, large repos, memory management, performance |
| `15-philosophy` | `studies/go-cli-study/reports/final/15-philosophy.md` | Cross-cutting design philosophy, tradeoffs, maintainability |

## Selected Reasoning Templates

All paths must appear in the project index's "Available Reasoning Templates" table.

| Template | Output Path | Why Selected |
| -------- | ----------- | ------------ |
| Architecture | `projects/ultraplan-go/sprints/06-study-run-execution/reasoning/report-validation.md` | Needed to reason through module ownership, validator boundaries, dependency direction, report/rating API shape, and diagnostics without making implementation decisions in the sprint index. |

## Prior Decisions To Carry Forward

All decision paths must appear in the project index's "Prior Decisions" table.

| Decision | Path | Constraint For This Sprint |
| -------- | ---- | -------------------------- |
| None listed in project index | N/A | Carry forward the requirements-defined prior behavior: Sprint 5 Markdown source kind and applicability behavior remains in force, so inapplicable Markdown source/dimension pairs stay skipped or excluded by callers and must not become validation failures. |

## Required Review Protocols

All paths must appear in the project index's "Review Protocols" table.

| Protocol | Path | Required Evidence |
| -------- | ---- | ----------------- |
| Architecture Review | `system/protocols/architecture-review-protocol.md` | Evidence that validation, report paths, rating parsing, and diagnostics remain in `internal/study`, platform packages do not import study, and no global validation/report/parsing package was introduced. |
| Sprint Review | `system/protocols/sprint-review-protocol.md` | Evidence that sprint artifacts and implementation satisfy requirements, acceptance criteria, non-goals, selected contracts, selected evidence, and `go test ./...` verification. |

## Excluded Context

| Context | Reason Excluded | Revisit If |
| ------- | --------------- | ---------- |
| LLM Runtime contract | Runtime execution, agentwrap/OpenCode wiring, policy runner wiring, retries, fallback, health checks, and event mapping are explicit non-goals for this sprint. | A later sprint implements single-run execution, runtime integration, repair, retry, fallback, or health behavior. |
| LLM Evaluation / Cost / Safety contract | Cost metadata, runtime safety, and LLM evaluation are not needed for local report validation and rating parsing. | Runtime execution or LLM-backed repair becomes sprint scope. |
| Workflows contract | Run-loop, orchestration, retries, cancellation, durable workflow execution, and resumability are explicit non-goals. | A later batch, run-loop, or durable workflow sprint begins. |
| Configuration contract | This sprint does not change config loading, precedence, redaction, or runtime config mapping. | Validation behavior becomes configurable or introduces user-facing configuration fields. |
| Documentation contract | User-facing command/help/generated docs are not direct outputs for this sprint. | A public validate command, stable JSON validation output, or user documentation enters scope. |
| CLI Surface contract | A public stable `ultraplan study <study> validate` command is excluded unless already present from prior work; current scope is internal validation APIs and tests. | Public validation command behavior or stable machine-readable CLI output enters scope. |
| Persistence And Migrations contract | Durable run state, migrations, and persisted workflow schemas are not in scope; any schema-like artifacts should remain internal/test-only or clearly versioned. | Validation results become persisted public artifacts or durable state schemas. |
| `02-command-architecture` evidence report | Command routing, flags, help text, shell completion, and command organization are not the main sprint focus. | The sprint adds or changes public CLI validation commands. |
| `03-dependency-injection` evidence report | The sprint can use existing package/test seams and does not require new runtime or composition-boundary design. | Validator dependencies become external, volatile, or difficult to test without explicit seams. |
| `04-configuration-management` evidence report | Config loading and precedence are not changed by report validation or rating parsing. | Validation behavior becomes configurable. |
| `07-state-context` evidence report | Context propagation, app state, and cancellation are tied to runtime/run-loop behavior that is out of scope. | Runtime execution or cancellation becomes sprint scope. |
| `08-concurrency` evidence report | Worker pools and parallel execution are explicit non-goals. | Batch validation or concurrent report processing enters scope. |
| `09-terminal-ux` evidence report | Terminal progress and human output formatting are not central because this sprint focuses on internal validation APIs and deterministic diagnostics. | A user-facing validate/status command enters scope. |
| `12-extensibility` evidence report | Plugin architecture and broad extension points are not needed for minimal study-owned validators. | Validation APIs become external extension points. |
| Target and sprint workflows | PRD, TRD, project index maintenance notes, and sprint requirements defer target scaffolding, sprint planning, and sprint execution. | The project scope is revised to include target or sprint workflows. |
| Prompt composition | Prompt building for directory analysis, Markdown analysis, or synthesis is explicitly excluded from this sprint. | A later prompt composition or runtime execution sprint begins. |
| Code reference extraction | Directory citation validation checks shape only; resolving cited files and extracting snippets is deferred. | The code extraction sprint begins. |
| Summary generation | Summary generation and `summary.csv` writing are explicit non-goals, though rating parsing should remain usable by later summary code. | The summary generation sprint begins. |

## Next Artifacts

- `technical-handbook.md` reads from the evidence reports listed above.
- `reasoning/*.md` captures area-specific reasoning.
- `reasoning.md` makes final sprint decisions.
- `plan.md` executes `reasoning.md`.
- `review.md` runs the selected review protocols against implementation.
