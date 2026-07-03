# Sprint Index: Project Domain and Project Index

> Project: `ultraplan-go`
> Sprint: `16-project-domain-and-index`
> Purpose: selected context for this sprint. Must be a subset of `projects/ultraplan-go/project-index.md`.
> **Inputs Used:** `projects/ultraplan-go/project-index.md`, `projects/ultraplan-go/sprints/16-project-domain-and-index/requirements.md`, `projects/ultraplan-go/docs/PRD.md`, `projects/ultraplan-go/docs/TRD.md`, `projects/ultraplan-go/docs/ARCHITECTURE.md`, `templates/sprint-index.md`

This document selects what must be read, distilled, reasoned through, or checked for this sprint. It does not make implementation decisions. All selections must come from the project index -- no items may be included that are not listed in the project index.

## Sprint Scope

- **Sprint Goal:** Model `projects/<project>` as a first-class planning root with project discovery, project status, and actionable `project-index.md` catalog validation without depending on study internals.
- **Planned Output:** `internal/project` domain, discovery, filesystem store, project-index parser, validation service, service API, project CLI wiring, top-level help registration, tests, and sprint review evidence for `ultraplan project list`, `ultraplan project <project> status`, and `ultraplan project <project> validate`.
- **Depends On:** Project index Phase 2 planning context, source documents, active contracts, selected go-cli-study evidence reports, available Architecture reasoning template, and required review protocols listed in `projects/ultraplan-go/project-index.md`.
- **Non-Goals:** Sprint planning artifacts, planning stages, `flow-state.json`, planning prompt generation, runtime-backed planning generation, study service consumption, study source/dimension models, study report validators, rating parsing, summary generation, run-loop scheduling, sprint execution, smoke investigation, automated review generation, issue tracking, Git mutation, browser UI, hosted service, multi-user collaboration, TUI, metrics exporter, local API server, and mutation of existing project docs or planning artifacts from project status or validation commands.

## Source Project Index

- `projects/ultraplan-go/project-index.md` -- authoritative source. Any file or item referenced below must appear there.

## Selected Contracts

Each contract applies as a flat whole to this sprint. All paths must appear in the project index's "Active Contract Pool" table.

| Contract | Why Selected |
| -------- | ------------ |
| Architecture | Required for `internal/project` ownership, dependency direction, thin entrypoints, and product/platform separation; especially important for keeping project behavior out of study internals and global shared packages. |
| Errors | Required for validation failure categories, preserved cause chains, actionable CLI diagnostics, and exit code mapping including validation exit code `5`. |
| Security | Required for workspace path safety, fail-closed catalog path resolution, path escape rejection, and secret-safe diagnostics. |
| Testing | Required for deterministic unit, fixture, and command tests covering discovery, parsing, validation, path failures, output ordering, and help behavior. |
| Documentation | Required for package documentation, top-level/project help text, and clear user-facing descriptions of project indexes as catalogs rather than sprint plans. |
| CLI Surface | Required for `ultraplan project list`, `ultraplan project <project> status`, `ultraplan project <project> validate`, help output, script-friendly text output, and exit codes. |
| Performance | Required for bounded project scans, deterministic ordering, and avoiding recursive scans of source repositories or evidence source trees. |

## Selected Evidence Reports

Copied from the project index's "Available Evidence Reports" table. These tell the technical handbook which reports to read -- the project index is the authoritative source.

| Report | Path | Covers |
| ------ | ---- | ------ |
| `01-project-structure` | `studies/go-cli-study/reports/final/01-project-structure.md` | Project layout, cmd/internal/pkg, dependency direction, thin entrypoints |
| `02-command-architecture` | `studies/go-cli-study/reports/final/02-command-architecture.md` | Command routing, flags, help text, command organization, shell completion |
| `03-dependency-injection` | `studies/go-cli-study/reports/final/03-dependency-injection.md` | Dependency construction, seams, testability, constructor patterns |
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
| -------- | ----------- | ------------ |
| Architecture | `projects/ultraplan-go/sprints/16-project-domain-and-index/reasoning/architecture.md` | Required to reason through `internal/project` ownership, CLI adapter thinness, dependency direction, study/runtime separation, and exclusion of speculative shared catalog or validation abstractions. |

## Prior Decisions To Carry Forward

All decision paths must appear in the project index's "Prior Decisions" table.

| Decision | Path | Constraint For This Sprint |
| -------- | ---- | -------------------------- |
| None yet -- this is the first project. | N/A | Carry forward the project-index maintenance note that the index is a catalog, not a sprint plan; apply requirements constraints preserving module-owned product behavior, thin CLI adapters, cause chains, safe diagnostics, and no speculative shared abstractions. |

## Required Review Protocols

All paths must appear in the project index's "Review Protocols" table.

| Protocol | Path | Required Evidence |
| -------- | ---- | ----------------- |
| Architecture Review | `system/protocols/architecture-review-protocol.md` | Evidence that `internal/project` owns project discovery, project docs, project-index parsing, catalog validation, and status behavior; `internal/app` stays thin; platform packages do not import product modules; no global catalog, validation, planning, reports, or prompts package is introduced. |
| Sprint Review | `system/protocols/sprint-review-protocol.md` | Evidence for implemented required outputs, acceptance criteria, test coverage, command help, deterministic output, validation diagnostics, non-goal exclusions, and verification with `go test ./...`, `go test -race ./...`, and `go build ./cmd/ultraplan`. |

## Excluded Context

| Context | Reason Excluded | Revisit If |
| ------- | --------------- | ---------- |
| Configuration contract | This sprint does not change config precedence, runtime mapping, or config validation; it only reuses existing workspace and command conventions where needed. | A later sprint changes project command configuration, workspace config schema, or runtime config behavior. |
| Observability contract | This sprint requires concise status and diagnostics but does not implement runtime events, structured logs, run metadata, or health/preflight observability. | A later sprint adds project/sprint event records, structured logging, or runtime-backed planning execution. |
| LLM Runtime contract | Project validation must not invoke agentwrap, OpenCode, runtime health checks, prompt generation, or network calls. | Sprint planning flow or runtime-backed planning generation enters scope. |
| LLM Evaluation / Cost / Safety contract | No runtime execution, token/cost accounting, model evaluation, or provider safety behavior is in scope. | Runtime-backed planning, evaluation, or cost reporting enters scope. |
| Workflows contract | Stateful batch workflow execution, retries, cancellation, resumability, and run-loop scheduling are non-goals for project catalog validation. | Sprint `flow-state.json` or planning-stage flow execution enters scope. |
| Persistence And Migrations contract | The sprint reads project files and validates catalogs; it does not introduce durable state formats, migrations, or atomic writes beyond any existing test fixtures. | A later sprint creates or mutates `flow-state.json`, run state, migrations, or other versioned durable formats. |
| `04-configuration-management` evidence report | Config loading and precedence are not being changed by project discovery, status, or validation. | Project commands add config-specific behavior or new config fields. |
| `07-state-context` evidence report | Context propagation and app state are not the central risk for read-only project catalog commands. | Long-running project or sprint flows with cancellation/state enter scope. |
| `08-concurrency` evidence report | Project discovery, status, and validation should be deterministic and bounded without worker pools or parallel scheduling. | Catalog validation becomes parallelized or long-running. |
| `10-logging-observability` evidence report | Runtime/event observability and structured logs are outside this sprint's read-only project command scope. | Structured project diagnostics or event logs are added. |
| `12-extensibility` evidence report | Plugin/extension-point design is not required; the minimal module-owned project service API is sufficient for this sprint. | External project APIs, plugins, or extension points become a concrete requirement. |
| Study services, source/dimension models, report validators, rating parsing, summary generation, and run-loop scheduling | Phase 2 planning context says project and sprint behavior must not abstract or consume study internals; requirements make this an acceptance criterion. | A future requirements revision explicitly defines a shared product abstraction and updates the project index. |
| Sprint planning artifacts and `flow-state.json` | Requirements make sprint artifacts, planning stages, and `flow-state.json` non-goals for Sprint 16; Sprint 17 owns that scope. | Sprint 17 or later planning-stage requirements are active. |
| Sprint execution, smoke investigation, review automation, issue tracking, and Git mutation | Project index and source documents explicitly defer these behaviors beyond Phase 2 planning through `plan.md`. | A later product phase changes the project index and PRD/TRD scope. |
| Browser UI, hosted service, multi-user collaboration, TUI, metrics exporter, and local API server | These product surfaces are explicitly outside the project catalog command scope and Phase 2 non-goals. | Product requirements add these surfaces to the active project index. |

## Next Artifacts

- `technical-handbook.md` reads from the evidence reports listed above.
- `reasoning/*.md` captures area-specific reasoning.
- `reasoning.md` makes final sprint decisions.
- `plan.md` executes `reasoning.md`.
- `review.md` runs the selected review protocols against implementation.
