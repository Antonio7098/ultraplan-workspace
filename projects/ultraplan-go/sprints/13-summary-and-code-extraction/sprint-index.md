# Sprint Index: 13 Summary And Code Extraction

> Project: `ultraplan-go`
> Sprint: `13-summary-and-code-extraction`
> Purpose: selected context for this sprint. Must be a subset of `projects/ultraplan-go/project-index.md`.
> **Inputs Used:** `projects/ultraplan-go/sprints/13-summary-and-code-extraction/requirements.md`, `projects/ultraplan-go/project-index.md`, `projects/ultraplan-go/docs/PRD.md`, `projects/ultraplan-go/docs/TRD.md`, `projects/ultraplan-go/docs/ARCHITECTURE.md`, `templates/sprint-index.md`

This document selects what must be read, distilled, reasoned through, or checked for this sprint. It does not make implementation decisions. All selections must come from the project index - no items may be included that are not listed in the project index.

## Sprint Scope

- **Sprint Goal:** Implement standalone study summary generation and code-reference extraction so completed study outputs can be reviewed, compared, and audited from deterministic text, CSV, and JSON surfaces.
- **Planned Output:** Sprint artifacts for summary generation and code extraction, plus implementation outputs for `ultraplan study <study> summary`, `ultraplan code <report>...`, deterministic `summary.csv`, text/JSON extraction output, and fixture-backed tests.
- **Depends On:** Prior study-side outputs named in the sprint requirements for Markdown source discovery and applicability, report validation and rating parsing, prompt/report artifact conventions, completed-output inspection, existing summary behavior, and run-loop output compatibility. No prior decision artifacts are listed in the project index.
- **Non-Goals:** Runtime execution, prompt composition changes, report template changes, run-loop orchestration, target scaffolding, sprint planning or execution, remote source fetching, browser/TUI/hosted workflows, plugin systems, workflow engines, default real OpenCode smoke tests, and release-wide stable public JSON schema guarantees.

## Source Project Index

- `projects/ultraplan-go/project-index.md` - authoritative source. Any file or item referenced below must appear there.

## Selected Contracts

Each contract applies as a flat whole to this sprint. All paths must appear in the project index's "Active Contract Pool" table.

| Contract | Why Selected |
| ------------ | -------------------------------------------- |
| Architecture | Summary behavior must stay in `internal/study`, code extraction in `internal/codeextract`, CLI wiring thin, and `internal/platform/runtime` product-free. |
| Errors | Summary warnings, validation-style extraction failures, unresolved references, write failures, and exit-code mapping need actionable diagnostics. |
| Security | Code extraction path resolution, basename fallback, ignored directories, workspace containment, and output redaction are central constraints. |
| Testing | The sprint requires deterministic unit, fixture, command-level, golden/exact-output, build, and offline verification coverage. |
| Documentation | CLI help text, sprint artifacts, deterministic user-facing outputs, and reviewable generated files must remain clear and maintainable. |
| CLI Surface | Both `study <study> summary` and top-level `code <report>...` require command shape, flags, help, text/JSON output, and exit-code handling. |
| Performance | Code extraction must avoid unbounded repository scans, cache lookup work within a command, and skip ignored directories during fallback search. |
| Persistence And Migrations | `summary.csv` and optional extraction output writes require explicit safe paths, stable artifact behavior, and atomic summary writes. |

## Selected Evidence Reports

Copied from the project index's "Available Evidence Reports" table. These tell the technical handbook which reports to read - the project index is the authoritative source.

| Report | Path | Covers |
| ----------- | ------------------------------------------------------ | ------------------------------- |
| `01-project-structure` | `studies/go-cli-study/reports/final/01-project-structure.md` | Project layout, cmd/internal/pkg, dependency direction, thin entrypoints |
| `02-command-architecture` | `studies/go-cli-study/reports/final/02-command-architecture.md` | Command routing, flags, help text, command organization, shell completion |
| `03-dependency-injection` | `studies/go-cli-study/reports/final/03-dependency-injection.md` | Dependency construction, seams, testability, constructor patterns |
| `05-error-handling` | `studies/go-cli-study/reports/final/05-error-handling.md` | Error wrapping, classification, user-facing diagnostics, sentinel errors |
| `06-io-abstraction` | `studies/go-cli-study/reports/final/06-io-abstraction.md` | Filesystem/stdin/stdout abstraction, test seams, interface design |
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
| Architecture | `projects/ultraplan-go/sprints/13-summary-and-code-extraction/reasoning/summary-generation.md` | Reason through study-owned summary behavior, existing discovery/applicability/rating helpers, deterministic CSV semantics, warnings, and atomic writes. |
| Architecture | `projects/ultraplan-go/sprints/13-summary-and-code-extraction/reasoning/code-extraction.md` | Reason through `internal/codeextract` ownership, parsing/resolution boundaries, safe source-root path handling, output shape, and thin CLI wiring. |

## Prior Decisions To Carry Forward

All decision paths must appear in the project index's "Prior Decisions" table.

| Decision | Path | Constraint For This Sprint |
| ------------ | -------- | -------------------------- |
| None listed in project index | N/A | No prior decision artifacts are available to select; carry-forward constraints come from the sprint requirements dependencies and selected project docs. |

## Required Review Protocols

All paths must appear in the project index's "Review Protocols" table.

| Protocol | Path | Required Evidence |
| ------------ | --------------------------------------- | ----------------- |
| Architecture Review | `system/protocols/architecture-review-protocol.md` | Verify summary remains in `internal/study`, extraction remains in `internal/codeextract`, CLI adapters stay thin, runtime stays generic, and no global product-parser/resolver packages appear. |
| Sprint Review | `system/protocols/review-sprint-protocol.md` | Verify required sprint artifacts exist, acceptance criteria are evidenced, deviations and residual risks are recorded, and `go test ./...` plus `go build ./cmd/ultraplan` results are captured. |

## Excluded Context

| Context | Reason Excluded | Revisit If |
| ----------- | --------------- | ------------- |
| Configuration contract and `04-configuration-management` report | Sprint scope does not change config loading, precedence, health mapping, or secret config display. | New config fields, config preflight behavior, or environment precedence changes enter scope. |
| LLM Runtime contract | Sprint must not launch, bypass, extend, or recompose runtime execution; it inspects completed artifacts only. | Implementation changes agentwrap/OpenCode runtime integration or runtime request construction. |
| LLM Evaluation / Cost / Safety contract | No provider execution, cost metadata, runtime evaluation, or model policy behavior is in scope. | Runtime execution, cost reporting, or provider safety behavior enters scope. |
| Workflows contract | Run-loop, retries, cancellation, stale task recovery, locks, and orchestration are non-goals except compatibility with completed outputs. | Sprint scope expands to orchestration or durable workflow mutation. |
| `07-state-context` report | This sprint must not require direct run-state mutation; completed-output inspection is covered by requirements and selected docs. | Summary or extraction starts reading or mutating run-state semantics. |
| `08-concurrency` report | The sprint requires deterministic local processing, not new worker pools or parallel runtime scheduling. | Extraction or summary introduces concurrent scanning or bounded workers. |
| `12-extensibility` report | Plugin architecture and public extension points are excluded; module boundaries are covered by Architecture selections. | A reusable public extraction, report, or plugin extension surface is added. |
| Target scaffolding, sprint planning, and sprint execution | Project scope and sprint requirements explicitly defer target/sprint workflow implementation. | PRD/TRD current-scope clarification is revised. |
| Hosted service, browser UI, TUI dashboard, remote artifact store, plugin system, and workflow engine | These are project and sprint non-goals unrelated to deterministic summary/code extraction artifacts. | Product scope expands beyond local study-side CLI workflows. |
| Remote repository resolution and source fetching | Code extraction must resolve local files only under parsed source roots or explicitly accepted local paths. | Requirements add remote source resolution or cloning for extraction. |

## Next Artifacts

- `technical-handbook.md` reads from the evidence reports listed above.
- `reasoning/*.md` captures area-specific reasoning.
- `reasoning.md` makes final sprint decisions.
- `plan.md` executes `reasoning.md`.
- `review.md` runs the selected review protocols against implementation.
