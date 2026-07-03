# Sprint Index: 07 Prompt Composition Generation

> Project: `ultraplan-go`
> Sprint: `07-prompt-composition-generation`
> Purpose: selected context for this sprint. Must be a subset of `projects/ultraplan-go/project-index.md`.
> **Inputs Used:** `projects/ultraplan-go/project-index.md`, `projects/ultraplan-go/sprints/07-prompt-composition-generation/requirements.md`, `projects/ultraplan-go/docs/ARCHITECTURE.md`, `projects/ultraplan-go/docs/PRD.md`, `projects/ultraplan-go/docs/TRD.md`, `templates/sprint-index.md`

This document selects what must be read, distilled, reasoned through, or checked for this sprint. It does not make implementation decisions. All selections must come from the project index — no items may be included that are not listed in the project index.

## Sprint Scope

- **Sprint Goal:** Implement deterministic study prompt composition for directory analysis, Markdown document analysis, and synthesis, including manifest-backed prompt previews, without invoking agentwrap or OpenCode.
- **Planned Output:** Sprint planning artifacts plus study prompt domain updates, prompt composition implementation, prompt preview CLI wiring, and unit/command tests for deterministic prompt text, manifests, template failures, Markdown frontmatter stripping, applicability handling, synthesis inputs, and no-runtime preview behavior.
- **Depends On:** Prior sprint capabilities listed in `requirements.md`: buildable CLI shell and app composition; workspace discovery, template path resolution, and path safety; study/source/dimension resolution; generated study layout; Markdown document source kind, frontmatter stripping, and applicability filtering; report validation/path/rating concepts for synthesis manifests.
- **Non-Goals:** Agentwrap/OpenCode runtime integration, actual analysis or synthesis runtime execution, runtime health checks, permission policy mapping, event/log mapping, retry/fallback/repair/cancellation/provider cost metadata, `run-all`, `run-loop`, durable task state, status calculations, worker pools, resumability, report validation changes beyond consuming prior path concepts, summary generation, code reference extraction, assisted study initialization, source cloning, Markdown source discovery changes, target workflows, sprint planning, and sprint execution.

## Source Project Index

- `projects/ultraplan-go/project-index.md` — authoritative source. Any file or item referenced below must appear there.

## Selected Contracts

Each contract applies as a flat whole to this sprint. All paths must appear in the project index's "Active Contract Pool" table.

| Contract | Why Selected |
| ------------ | -------------------------------------------- |
| Architecture | Governs study-owned prompt behavior in `internal/study`, CLI wiring in `internal/app`, dependency direction, and avoiding global prompt/template packages. |
| Errors | Governs actionable diagnostics for missing templates, missing Markdown documents, unknown or ambiguous study/dimension/source references, inapplicable pairs, and missing synthesis inputs. |
| Security | Governs workspace path safety, source isolation, secret-safe generated artifacts, and the sprint constraint that prompt preview must not invoke subprocesses, OpenCode, providers, network access, or runtime execution. |
| Testing | Governs deterministic unit, fixture, and command tests for prompt composition, manifests, missing input failures, and no-runtime preview behavior. |
| Documentation | Governs maintainable generated sprint artifacts and user-facing command/help behavior for prompt preview. |
| CLI Surface | Governs prompt preview command shape, flags, help behavior, output paths/stdout behavior, exit codes, and script-friendly diagnostics. |

## Selected Evidence Reports

Copied from the project index's "Available Evidence Reports" table. These tell the technical handbook which reports to read – the project index is the authoritative source.

| Report | Path | Covers |
| ----------- | ------------------------------------------------------ | ------------------------------- |
| `01-project-structure` | `studies/go-cli-study/reports/final/01-project-structure.md` | Project layout, cmd/internal/pkg, dependency direction, thin entrypoints |
| `02-command-architecture` | `studies/go-cli-study/reports/final/02-command-architecture.md` | Command routing, flags, help text, command organization, shell completion |
| `03-dependency-injection` | `studies/go-cli-study/reports/final/03-dependency-injection.md` | Dependency construction, seams, testability, constructor patterns |
| `05-error-handling` | `studies/go-cli-study/reports/final/05-error-handling.md` | Error wrapping, classification, user-facing diagnostics, sentinel errors |
| `06-io-abstraction` | `studies/go-cli-study/reports/final/06-io-abstraction.md` | Filesystem/stdin/stdout abstraction, test seams, interface design |
| `11-testing-strategy` | `studies/go-cli-study/reports/final/11-testing-strategy.md` | Unit, integration, fixture, command-level tests, coverage strategy |
| `13-security` | `studies/go-cli-study/reports/final/13-security.md` | Path safety, secrets, command injection risks, sandboxing |
| `14-performance` | `studies/go-cli-study/reports/final/14-performance.md` | Startup latency, large repos, memory management, performance |
| `15-philosophy` | `studies/go-cli-study/reports/final/15-philosophy.md` | Cross-cutting design philosophy, tradeoffs, maintainability |

## Selected Reasoning Templates

All paths must appear in the project index's "Available Reasoning Templates" table.

| Template | Output Path | Why Selected |
| ------------ | -------------------------------------------------------------- | ------------ |
| Architecture | `projects/ultraplan-go/sprints/07-prompt-composition-generation/reasoning/prompt-composition.md` | Needed to reason through prompt ownership, template boundaries, source-kind behavior, synthesis manifest shape, preview boundaries, and runtime/product separation. |

## Prior Decisions To Carry Forward

All decision paths must appear in the project index's "Prior Decisions" table.

| Decision | Path | Constraint For This Sprint |
| ------------ | -------- | -------------------------- |
| None yet | None | The project index lists no prior decisions. Carry-forward dependencies from prior sprint reviews are recorded in `requirements.md` and constrain this sprint through the listed completed capabilities: Markdown source applicability/frontmatter behavior and report path/validation concepts. |

## Required Review Protocols

All paths must appear in the project index's "Review Protocols" table.

| Protocol | Path | Required Evidence |
| ------------ | --------------------------------------- | ----------------- |
| Architecture Review | `system/protocols/architecture-review-protocol.md` | Evidence that prompt behavior remains study-owned, CLI preview wiring remains app-owned, platform packages do not import product modules, no global prompt/template package is introduced, and runtime execution remains out of scope. |
| Sprint Review | `system/protocols/sprint-review-protocol.md` | Evidence that required sprint artifacts and implementation outputs exist, acceptance criteria are checked, no non-goal scope was added, and `go test ./...` plus `go build ./cmd/ultraplan` results are recorded. |

## Excluded Context

| Context | Reason Excluded | Revisit If |
| ----------- | --------------- | ------------- |
| LLM Runtime contract | Agentwrap/OpenCode runtime integration and actual runtime execution are explicit non-goals for this prompt composition sprint. | A later sprint implements analysis or synthesis runtime execution. |
| LLM Evaluation / Cost / Safety contract | Provider cost metadata, runtime safety evaluation, and runtime execution behavior are out of scope. | Runtime-backed execution, usage accounting, or provider safety gates enter scope. |
| Observability contract | Event/log mapping, runtime diagnostics, status, and run metadata are non-goals beyond actionable local prompt-preview errors. | Runtime execution, run-loop, status, or structured event persistence enters scope. |
| Workflows contract | `run-all`, `run-loop`, retries, cancellation, resumability, durable task state, and orchestration are explicit non-goals. | Batch or durable workflow execution enters scope. |
| Configuration contract | Runtime configuration mapping and health/preflight behavior are non-goals; this sprint only needs resolved workspace paths for prompt/template inputs. | Config precedence, runtime options, health checks, or provider/model mapping enter scope. |
| Persistence And Migrations contract | Durable run state, schema migration, and compatibility management are out of scope; prompt manifests are deterministic previews, not durable workflow state. | Persisted state formats or migration behavior enter scope. |
| Runtime-backed analysis and synthesis | This sprint composes and previews prompts only and must not invoke agentwrap, OpenCode, providers, network access, or subprocess runtime execution. | The product moves from prompt preview to actual `study run` or `study synthesize` execution. |
| Code reference extraction | Extraction is a later product capability and acceptance criteria only require prompts to request citations for directory sources where applicable. | The code extraction command or citation resolver enters scope. |
| Summary generation | Summary CSV writing and rating aggregation are explicit non-goals. | Summary generation enters scope. |
| Target and sprint workflows | Project index maintenance notes, PRD, and TRD defer target scaffolding, sprint planning, and sprint execution. | A later requirements revision includes target or sprint-side workflows. |

## Next Artifacts

- `technical-handbook.md` reads from the evidence reports listed above.
- `reasoning/*.md` captures area-specific reasoning.
- `reasoning.md` makes final sprint decisions.
- `plan.md` executes `reasoning.md`.
- `review.md` runs the selected review protocols against implementation.
