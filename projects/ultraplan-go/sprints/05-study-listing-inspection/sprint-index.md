# Sprint Index: 05 Study Listing Inspection

> Project: `ultraplan-go`
> Sprint: `05-study-listing-inspection`
> Purpose: selected context for this sprint. Must be a subset of `projects/ultraplan-go/project-index.md`.
> **Inputs Used:** `projects/ultraplan-go/sprints/05-study-listing-inspection/requirements.md`, `projects/ultraplan-go/project-index.md`, `projects/ultraplan-go/docs/PRD.md`, `projects/ultraplan-go/docs/TRD.md`, `projects/ultraplan-go/docs/ARCHITECTURE.md`, `templates/sprint-index.md`

This document selects what must be read, distilled, reasoned through, or checked for this sprint. It does not make implementation decisions. All selections must come from the project index — no items may be included that are not listed in the project index.

## Sprint Scope

- **Sprint Goal:** Add top-level Markdown document source discovery, frontmatter-based dimension applicability, and visible study listing/inspection output for source kind and applicability without introducing runtime execution behavior.
- **Planned Output:** `sprint-index.md`, `technical-handbook.md`, `reasoning/source-applicability.md`, `reasoning.md`, `plan.md`, `review.md`, and study-module implementation/test changes for source discovery, Markdown frontmatter parsing, applicability filtering, and deterministic `ultraplan study <study> list` output.
- **Depends On:** Project source documents listed in `projects/ultraplan-go/project-index.md`: Product Requirements, Technical Requirements, and Architecture. Sprint requirements also identify prior build dependencies for the CLI shell, workspace discovery, study listing/resolution, and initialized study layout.
- **Non-Goals:** Runtime execution for Markdown document sources; prompt composition for Markdown document analysis; report validation changes; summary generation changes; `run`, `run-all`, `run-loop`, `status`, and synthesis scheduling changes beyond compile-time compatibility; source cloning; assisted study initialization; YAML schema expansion; JSON output schema stabilization; target workflows; sprint planning; sprint execution.

## Source Project Index

- `projects/ultraplan-go/project-index.md` — authoritative source. Any file or item referenced below must appear there.

## Selected Contracts

Each contract applies as a flat whole to this sprint. All paths must appear in the project index's "Active Contract Pool" table.

| Contract | Why Selected |
| ------------ | -------------------------------------------- |
| Architecture | Governs `internal/study` ownership of source discovery, applicability, CLI behavior, and dependency direction; especially important because behavior must not move into global validation, scheduler, reports, or prompts packages. |
| Errors | Governs wrapped filesystem and malformed frontmatter errors, source path context, actionable CLI diagnostics, and cause preservation. |
| Security | Governs workspace path safety, source isolation, secret-safe diagnostics, and avoiding unsafe traversal or runtime behavior. |
| Testing | Governs unit, fixture, and command-level coverage for deterministic discovery, frontmatter handling, applicability filtering, and listing output. |
| CLI Surface | Governs script-friendly `ultraplan study <study> list` output, command behavior, diagnostics, and exit behavior. |

## Selected Evidence Reports

Copied from the project index's "Available Evidence Reports" table. These tell the technical handbook which reports to read – the project index is the authoritative source.

| Report | Path | Covers |
|---|---|---|
| `01-project-structure` | `studies/go-cli-study/reports/final/01-project-structure.md` | Project layout, cmd/internal/pkg, dependency direction, thin entrypoints |
| `02-command-architecture` | `studies/go-cli-study/reports/final/02-command-architecture.md` | Command routing, flags, help text, command organization, shell completion |
| `05-error-handling` | `studies/go-cli-study/reports/final/05-error-handling.md` | Error wrapping, classification, user-facing diagnostics, sentinel errors |
| `06-io-abstraction` | `studies/go-cli-study/reports/final/06-io-abstraction.md` | Filesystem/stdin/stdout abstraction, test seams, interface design |
| `09-terminal-ux` | `studies/go-cli-study/reports/final/09-terminal-ux.md` | Terminal output, progress indicators, human UX, color and formatting |
| `11-testing-strategy` | `studies/go-cli-study/reports/final/11-testing-strategy.md` | Unit, integration, fixture, command-level tests, coverage strategy |
| `13-security` | `studies/go-cli-study/reports/final/13-security.md` | Path safety, secrets, command injection risks, sandboxing |
| `14-performance` | `studies/go-cli-study/reports/final/14-performance.md` | Startup latency, large repos, memory management, performance |
| `15-philosophy` | `studies/go-cli-study/reports/final/15-philosophy.md` | Cross-cutting design philosophy, tradeoffs, maintainability |

## Selected Reasoning Templates

All paths must appear in the project index's "Available Reasoning Templates" table.

| Template | Output Path | Why Selected |
| ------------ | -------------------------------------------------------------- | ------------ |
| Architecture | `projects/ultraplan-go/sprints/05-study-listing-inspection/reasoning/source-applicability.md` | Reason through study-owned source kinds, frontmatter normalization, applicability filtering, malformed frontmatter behavior, listing semantics, and scope boundaries without making runtime implementation decisions. |

## Prior Decisions To Carry Forward

All decision paths must appear in the project index's "Prior Decisions" table.

No prior decision artifacts are listed in the project index. Carry-forward constraints for this sprint therefore come from the selected source documents and sprint requirements: preserve the existing CLI shell and package structure, preserve workspace discovery and path safety, extend rather than replace deterministic study listing/resolution, and respect the initialized `studies/<study>/sources/` and `dimensions/` layout.

| Decision | Path | Constraint For This Sprint |
| ------------ | -------- | -------------------------- |
| No prior decision artifact listed in project index | None | Do not invent prior decision paths; use the selected source documents and sprint requirements for carry-forward constraints. |

## Required Review Protocols

All paths must appear in the project index's "Review Protocols" table.

| Protocol | Path | Required Evidence |
| ------------ | --------------------------------------- | ----------------- |
| Architecture Review | `system/protocols/architecture-review-protocol.md` | Evidence that source discovery, frontmatter parsing, applicability filtering, and listing behavior remain owned by `internal/study`; platform packages do not import product modules; no global validation, scheduler, reports, or prompts package is introduced for this sprint. |
| Sprint Review | `system/protocols/sprint-review-protocol.md` | Evidence that implementation, tests, and sprint artifacts satisfy requirements; `go test ./...` result is recorded; non-goals and exclusions were respected. |

## Excluded Context

| Context | Reason Excluded | Revisit If |
| ----------- | --------------- | ------------- |
| `Configuration` contract | This sprint does not change workspace config, config precedence, runtime mapping, or config validation. | A later sprint changes config files, environment overrides, runtime model selection, or config diagnostics. |
| `Observability` contract | Runtime logs, structured events, run metadata, and health/preflight reporting are outside the current listing/inspection scope. | Runtime execution, run-loop, status, logs, or diagnostics enter scope. |
| `Documentation` contract | User-facing generated documentation and help documentation are not required beyond preserving script-friendly listing behavior. | The sprint adds or changes generated docs, command help docs, or artifact documentation. |
| `LLM Runtime` contract | Markdown document runtime execution and agentwrap/OpenCode integration are explicit non-goals. | Prompt execution, OpenCode adapter use, or runtime request mapping enters scope. |
| `LLM Evaluation / Cost / Safety` contract | Cost metadata, runtime safety, and evaluation discipline are tied to runtime execution, not source listing. | Runtime execution or validation of generated model outputs enters scope. |
| `Workflows` contract | Run-loop orchestration, retries, cancellation, resumability, and stateful scheduling are explicit non-goals except compile-time compatibility. | Batch execution, stateful run-loop, or synthesis scheduling behavior is implemented. |
| `Persistence And Migrations` contract | Durable run state, schema versions, migrations, and atomic persistence formats are not changed by listing inspection. | The sprint writes durable state or changes persisted artifact schemas. |
| `03-dependency-injection` evidence report | This sprint can be covered through existing package seams and command tests without introducing new construction patterns. | New volatile collaborators or runtime/filesystem ports are introduced. |
| `04-configuration-management` evidence report | Config loading and precedence are unchanged. | Workspace or runtime configuration behavior changes. |
| `07-state-context` evidence report | App state and cancellation are runtime/run-loop concerns outside the source listing scope. | Context propagation, cancellation, or task state behavior changes. |
| `08-concurrency` evidence report | Bounded workers and goroutine management are not relevant to direct child source discovery and listing. | Batch scheduling or parallel execution enters scope. |
| `10-logging-observability` evidence report | Logging and structured observability are not changed by this sprint. | Status, logs, event records, or diagnostics are expanded. |
| `12-extensibility` evidence report | Plugin and extension-point design is not needed for this narrow study-module extension. | Public extension seams or package APIs are introduced. |
| Runtime execution for Markdown document sources | Explicit sprint non-goal; listing can expose kind and applicability without invoking runtime behavior. | A future sprint implements `run`, `run-all`, or prompt execution for Markdown sources. |
| Prompt composition for Markdown document analysis | Explicit sprint non-goal; frontmatter stripping may be implemented only as source-model support required by this sprint. | A future sprint builds prompts for Markdown document sources. |
| Report validation and summary generation changes | Explicit sprint non-goals; inapplicable pairs only need listing-derived applicability checks in this sprint. | Validation, synthesis gating, or summary behavior enters scope. |
| Target and sprint workflows | Project and sprint non-goals; UltraPlan Go current build is study-side only. | The PRD/TRD are revised to include target or sprint execution workflows. |

## Next Artifacts

- `technical-handbook.md` reads from the evidence reports listed above.
- `reasoning/*.md` captures area-specific reasoning.
- `reasoning.md` makes final sprint decisions.
- `plan.md` executes `reasoning.md`.
- `review.md` runs the selected review protocols against implementation.
