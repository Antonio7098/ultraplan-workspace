# Sprint Index: Study Domain, Listing, and Resolution

> Project: `ultraplan-go`
> Sprint: `03-study-listing-resolution`
> Purpose: selected context for this sprint. Must be a subset of `projects/ultraplan-go/project-index.md`.
> **Inputs Used:** `projects/ultraplan-go/sprints/03-study-listing-resolution/requirements.md`, `projects/ultraplan-go/project-index.md`, `projects/ultraplan-go/docs/PRD.md`, `projects/ultraplan-go/docs/TRD.md`, `projects/ultraplan-go/docs/ARCHITECTURE.md`, `templates/sprint-index.md`

This document selects what must be read, distilled, reasoned through, or checked for this sprint. It does not make implementation decisions. All selections must come from the project index - no items may be included that are not listed in the project index.

## Sprint Scope

- **Sprint Goal:** Implement the study domain model, deterministic study/source/dimension discovery, reference resolution, and listing commands needed to inspect existing studies without running analysis or initialization workflows.
- **Planned Output:** `internal/study` domain, discovery, resolution, and listing service; `internal/app` study listing commands; unit and command tests for deterministic discovery, resolution, output, and errors.
- **Depends On:** Sprint 1 command registration and buildable CLI entrypoint; Sprint 2 workspace discovery and global `--workspace` behavior; PRD study listing behavior; TRD study/source/dimension model and lookup rules; ARCHITECTURE package boundaries.
- **Non-Goals:** Study initialization and YAML parsing; Markdown document source discovery, frontmatter parsing, and applicability filtering; prompt composition; runtime execution; synthesis; report validation; summary generation; run state; retry; cancellation; code-reference extraction; target or sprint CLI workflows; real OpenCode or `agentwrap` runtime wiring; non-trivial JSON listing output.

## Source Project Index

- `projects/ultraplan-go/project-index.md` - authoritative source. Any file or item referenced below must appear there.

## Selected Contracts

Contracts apply through the Skeleton / Local CLI gate in `roadmap.md`, not as a flat production-release gate. All paths must appear in the project index's "Active Contract Pool" table.

| Contract | Why Selected |
| ------------ | -------------------------------------------- |
| Architecture | Governs `internal/study`, `internal/app`, `internal/workspace`, and platform dependency direction; especially important for keeping domain discovery and resolution out of CLI handlers. |
| Errors | Required for actionable missing and ambiguous study/source/dimension reference errors and meaningful non-zero command behavior. |
| Security | Required for workspace path safety, shallow discovery, and preventing listing commands from escaping the resolved workspace or scanning source repository contents recursively. |
| Testing | Required for deterministic unit and command tests covering ordering, ignored entries, resolution, output, and errors. |
| CLI Surface | Required for `ultraplan study list`, `ultraplan study <study> list`, global `--workspace`, help/exit behavior, and stable human-facing output. |

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
| `09-terminal-ux` | `studies/go-cli-study/reports/final/09-terminal-ux.md` | Terminal output, progress indicators, human UX, color and formatting |
| `11-testing-strategy` | `studies/go-cli-study/reports/final/11-testing-strategy.md` | Unit, integration, fixture, command-level tests, coverage strategy |
| `13-security` | `studies/go-cli-study/reports/final/13-security.md` | Path safety, secrets, command injection risks, sandboxing |
| `14-performance` | `studies/go-cli-study/reports/final/14-performance.md` | Startup latency, large repos, memory management, performance |
| `15-philosophy` | `studies/go-cli-study/reports/final/15-philosophy.md` | Cross-cutting design philosophy, tradeoffs, maintainability |

## Selected Reasoning Templates

All paths must appear in the project index's "Available Reasoning Templates" table.

| Template | Output Path | Why Selected |
| ------------ | -------------------------------------------------------------- | ------------ |
| Architecture | `projects/ultraplan-go/sprints/03-study-listing-resolution/reasoning/architecture.md` | Needed to reason through `study`, `workspace`, `app`, and platform boundaries before implementing domain discovery, resolution, and command wiring. |

## Prior Decisions To Carry Forward

No prior decision records are listed in the project index. Carry-forward constraints from the sprint requirements are: Sprint 1 provides command registration and the buildable CLI entrypoint; Sprint 2 provides workspace discovery and global `--workspace` behavior that listing commands must use.

## Required Review Protocols

All paths must appear in the project index's "Review Protocols" table.

| Protocol | Path | Required Evidence |
| ------------ | --------------------------------------- | ----------------- |
| Architecture Review | `system/protocols/architecture-review-protocol.md` | Evidence that `internal/study` owns domain discovery and resolution, `internal/app` owns CLI wiring, workspace behavior is reused, and platform packages do not import `internal/study`. |
| Sprint Review | `system/protocols/sprint-review-protocol.md` | Evidence that required outputs exist, acceptance criteria are satisfied, deferred scope remains excluded, and `go test ./...` plus `go build ./cmd/ultraplan` pass. |

## Excluded Context

| Context | Reason Excluded | Revisit If |
| ----------- | --------------- | ------------- |
| Configuration contract | Sprint 2 already owns the config foundation; this sprint only consumes existing workspace discovery and global `--workspace` behavior. | Listing needs new configuration precedence, validation, redaction, or runtime mapping behavior. |
| Observability contract | Sprint 3 does not implement runtime logs, run metadata, event records, or status workflows. | Study run-loop, status, logs, or runtime diagnostics enter scope. |
| Documentation contract | No standalone documentation artifacts are required beyond CLI help/output behavior covered by CLI Surface. | User-facing docs, generated docs, or artifact documentation become a required output. |
| LLM Runtime contract | Real OpenCode and `agentwrap` runtime wiring are explicitly out of scope. | Analysis, synthesis, prompt execution, health checks, or runtime adapters enter scope. |
| LLM Evaluation / Cost / Safety contract | No model execution, cost metadata, runtime safety evaluation, or generated report evaluation occurs in this sprint. | Runtime execution or model-output validation enters scope. |
| Workflows contract | No run-loop, orchestration, retries, resumability, or stateful workflow execution is implemented. | Batch study execution or resumable workflows enter scope. |
| Performance contract | Bounded workers and report processing are out of scope; only shallow deterministic listing performance is relevant and covered through selected performance evidence. | Scheduler, concurrency, large artifact processing, or run status performance enters scope. |
| Persistence And Migrations contract | No run state, persisted schemas, migrations, or atomic state writes are required for listing. | Study state, run state, schema versions, or migrations enter scope. |
| `07-state-context` evidence report | Context propagation and cancellation are not central to read-only listing commands. | Long-running commands, cancellation, or run-loop state enter scope. |
| `08-concurrency` evidence report | Listing commands must be shallow and deterministic, not concurrent worker orchestration. | Scheduler, run-all, or parallel runtime execution enters scope. |
| `10-logging-observability` evidence report | Runtime observability and structured event logging are outside this sprint's required outputs. | Status, logs, event records, or runtime diagnostics enter scope. |
| `12-extensibility` evidence report | Plugin and extension-point design is not needed for the minimal study listing service. | New public extension seams or plugin architecture enter scope. |
| Markdown document source discovery | Sprint requirements explicitly defer top-level `.md` source support, frontmatter parsing, and applicability filtering to a later sprint. | Sprint 5 or revised requirements include Markdown document sources. |
| Target and sprint workflows | PRD, TRD, and sprint requirements explicitly defer target and sprint CLI workflows. | A future requirements revision brings target or sprint workflows into the study-side build. |

## Next Artifacts

- `technical-handbook.md` reads from the evidence reports listed above.
- `reasoning/*.md` captures area-specific reasoning.
- `reasoning.md` makes final sprint decisions.
- `plan.md` executes `reasoning.md`.
- `review.md` runs the selected review protocols against implementation.
