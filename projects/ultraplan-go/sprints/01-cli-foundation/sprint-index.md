# Sprint Index: CLI Foundation

> Project: `ultraplan-go`
> Sprint: `01-cli-foundation`
> Purpose: selected context for this sprint. Must be a subset of `projects/ultraplan-go/project-index.md`.
> **Inputs Used:** `projects/ultraplan-go/project-index.md`, `projects/ultraplan-go/sprints/01-cli-foundation/requirements.md`, `projects/ultraplan-go/docs/PRD.md`, `projects/ultraplan-go/docs/TRD.md`, `projects/ultraplan-go/docs/ARCHITECTURE.md`, `templates/sprint-index.md`

This document selects what must be read, distilled, reasoned through, or checked for this sprint. It does not make implementation decisions. All selections must come from the project index -- no items may be included that are not listed in the project index.

## Sprint Scope

- **Sprint Goal:** Establish a buildable Go module with a thin `ultraplan` CLI shell, app composition root, version/help behavior, and the initial module-owned package layout for future study-side features.
- **Planned Output:** Go module definition, thin `cmd/ultraplan` entrypoint, `internal/app` composition root with help/version/invalid-command behavior, deterministic app command tests, and documented skeleton packages for `platform/config`, `platform/logging`, `platform/filesystem`, `platform/runtime`, `workspace`, `study`, and `codeextract`.
- **Depends On:** Phase 0 roadmap and scope alignment; `projects/ultraplan-go/docs/PRD.md`; `projects/ultraplan-go/docs/TRD.md`; `projects/ultraplan-go/docs/ARCHITECTURE.md`; Go toolchain. No prior sprint carry-forward decisions exist for this first implementation sprint.
- **Non-Goals:** Workspace discovery, config behavior, health checks, study modeling or workflows, runtime execution, agentwrap/OpenCode integration, code-reference extraction behavior, target/sprint workflows, packaging/release artifacts, generated user documentation, and real-runtime smoke tests.

## Source Project Index

- `projects/ultraplan-go/project-index.md` -- authoritative source. Any file or item referenced below must appear there.

## Selected Contracts

Each contract applies as a flat whole to this sprint. All paths must appear in the project index's "Active Contract Pool" table.

| Contract | Why Selected |
| --- | --- |
| Architecture | Applies to module ownership, dependency direction, thin entrypoint, app composition root, and product/platform separation. |
| Errors | Applies to actionable invalid-command diagnostics and exit code mapping, especially usage errors returning `2`. |
| Testing | Applies to deterministic offline command tests and normal `go test ./...` verification. |
| Documentation | Applies to user-facing help/version text and maintainable package boundary documentation. |
| CLI Surface | Applies to command shape, help behavior, script-friendly output, and meaningful exit codes. |

## Selected Evidence Reports

Copied from the project index's "Available Evidence Reports" table. These tell the technical handbook which reports to read -- the project index is the authoritative source.

| Report | Path | Covers |
| --- | --- | --- |
| `01-project-structure` | `studies/go-cli-study/reports/final/01-project-structure.md` | Project layout, cmd/internal/pkg, dependency direction, thin entrypoints |
| `02-command-architecture` | `studies/go-cli-study/reports/final/02-command-architecture.md` | Command routing, flags, help text, command organization, shell completion |
| `03-dependency-injection` | `studies/go-cli-study/reports/final/03-dependency-injection.md` | Dependency construction, seams, testability, constructor patterns |
| `05-error-handling` | `studies/go-cli-study/reports/final/05-error-handling.md` | Error wrapping, classification, user-facing diagnostics, sentinel errors |
| `06-io-abstraction` | `studies/go-cli-study/reports/final/06-io-abstraction.md` | Filesystem/stdin/stdout abstraction, test seams, interface design |
| `09-terminal-ux` | `studies/go-cli-study/reports/final/09-terminal-ux.md` | Terminal output, progress indicators, human UX, color and formatting |
| `11-testing-strategy` | `studies/go-cli-study/reports/final/11-testing-strategy.md` | Unit, integration, fixture, command-level tests, coverage strategy |
| `12-extensibility` | `studies/go-cli-study/reports/final/12-extensibility.md` | Extension points, plugin architecture, package boundaries, API design |
| `15-philosophy` | `studies/go-cli-study/reports/final/15-philosophy.md` | Cross-cutting design philosophy, tradeoffs, maintainability |

## Selected Reasoning Templates

All paths must appear in the project index's "Available Reasoning Templates" table.

| Template | Output Path | Why Selected |
| --- | --- | --- |
| Architecture | `projects/ultraplan-go/sprints/01-cli-foundation/reasoning/architecture.md` | Needed to reason through package boundaries, dependency direction, thin entrypoint, app composition root, and runtime/product separation without making implementation decisions in this index. |

## Prior Decisions To Carry Forward

No prior decisions are listed in the project index, and the sprint requirements state that no prior sprint review exists for this first implementation sprint.

## Required Review Protocols

All paths must appear in the project index's "Review Protocols" table.

| Protocol | Path | Required Evidence |
| --- | --- | --- |
| Architecture Review | `system/protocols/architecture-review-protocol.md` | Confirm the package layout, thin entrypoint, app composition root, product/platform dependency direction, and deferred runtime/module separation match the selected architecture context. |
| Sprint Review | `system/protocols/sprint-review-protocol.md` | Confirm all sprint acceptance criteria are reviewed, including `go test ./...`, `go build ./cmd/ultraplan`, help/version/invalid-command behavior, package skeletons, and scope exclusions. |

## Excluded Context

| Context | Reason Excluded | Revisit If |
| --- | --- | --- |
| Configuration contract and `04-configuration-management` | This sprint creates only a config package skeleton; workspace discovery, config precedence, `config show`, redaction behavior, and health checks are explicit non-goals. | A workspace/config sprint introduces real configuration behavior. |
| Observability contract and `10-logging-observability` | This sprint creates only a logging package skeleton; structured logs, diagnostics, runtime events, and run metadata are not implemented. | Runtime, run-loop, status, or logging behavior enters scope. |
| Security contract and `13-security` | This sprint does not implement workspace path handling, runtime permissions, source isolation, secret redaction behavior, or mutating workflows. | Workspace paths, config secrets, runtime execution, or code extraction behavior enters scope. |
| LLM Runtime and LLM Evaluation / Cost / Safety contracts | Agentwrap, OpenCode, provider/model handling, runtime health, usage metadata, and runtime safety are explicitly deferred. | Runtime integration is scheduled and uses `github.com/Antonio7098/agentwrap` and `agentwrap/opencode`. |
| Workflows, Persistence And Migrations, Performance contracts and `07-state-context`, `08-concurrency`, `14-performance` | The sprint has no run-loop, scheduler, durable state, migrations, worker pools, or large-artifact processing; only a buildable CLI shell and skeleton layout are required. | Study run-loop, resumability, persistence, concurrency, or performance-sensitive behavior enters scope. |
| Study workflow implementation | Study domain modeling, study initialization, listing, prompt composition, analysis runs, synthesis, summaries, and validation are explicit non-goals even though the `study` package skeleton is required. | A later study-side sprint selects study behavior requirements. |
| Code-reference extraction implementation | Citation parsing, resolution, snippet extraction, and `ultraplan code` are explicit non-goals even though the `codeextract` package skeleton is required. | A code extraction sprint selects code-reference extraction requirements. |
| Target and sprint workflows | The PRD, TRD, project index, and sprint requirements all defer target scaffolding, sprint planning, sprint execution, target commands, and sprint commands. | The product requirements are revised to include target/sprint workflows. |
| Packaging and release artifacts | Checksums, shell completion, release binaries, generated user documentation, and real-runtime smoke tests are outside this sprint's required outputs and acceptance criteria. | Release hardening or distribution work enters scope. |

## Next Artifacts

- `technical-handbook.md` reads from the evidence reports listed above.
- `reasoning/*.md` captures area-specific reasoning.
- `reasoning.md` makes final sprint decisions.
- `plan.md` executes `reasoning.md`.
- `review.md` runs the selected review protocols against implementation.
