# Sprint Index: Workspace, Config, Logging, and Health Skeleton

> Project: `ultraplan-go`
> Sprint: `02-workspace-config-health`
> Purpose: selected context for this sprint. Must be a subset of `projects/ultraplan-go/project-index.md`.
> **Inputs Used:** `projects/ultraplan-go/project-index.md`, `projects/ultraplan-go/sprints/02-workspace-config-health/requirements.md`, `projects/ultraplan-go/docs/ARCHITECTURE.md`, `projects/ultraplan-go/docs/PRD.md`, `projects/ultraplan-go/docs/TRD.md`, `templates/sprint-index.md`

This document selects what must be read, distilled, reasoned through, or checked for this sprint. It does not make implementation decisions. All selections must come from the project index — no items may be included that are not listed in the project index.

## Sprint Scope

- **Sprint Goal:** Make UltraPlan able to find, initialize, inspect, and validate a workspace — establishing workspace discovery, config loading, logging foundation, and health command skeleton.
- **Planned Output:** Sprint 2 implementation guidance and implementation artifacts for workspace discovery, workspace initialization, workspace path/validation helpers, effective config loading and redaction, minimal logging/output helpers, `init-workspace`, `config show`, and `health` command skeletons with deterministic tests.
- **Depends On:** Sprint 1: CLI Foundation for Go module layout, `cmd/ultraplan`, `internal/app`, command routing baseline, and initial package layout.
- **Non-Goals:** Workspace validation beyond structural workspace/config/filesystem checks; config hot-reload; persistent workspace state beyond the initial scaffold; study listing, study initialization, source discovery, and analysis runs; runtime execution, agentwrap integration, OpenCode, and agent health checks; durable file logs, full structured log level policy, run event records, and observability event persistence; target scaffolding, sprint planning, and sprint execution workflows.

## Source Project Index

- `projects/ultraplan-go/project-index.md` — authoritative source. Any file or item referenced below must appear there.

## Selected Contracts

Each contract applies as a flat whole to this sprint. All paths must appear in the project index's "Active Contract Pool" table.

| Contract | Why Selected |
| -------- | ------------ |
| Architecture | Applies to module ownership, dependency direction, thin entrypoints, and product/platform separation for `workspace`, `platform/config`, `platform/logging`, and `app` command wiring. |
| Errors | Applies to meaningful exit codes, error wrapping, actionable diagnostics, and workspace/config/validation error classes. |
| Configuration | Applies to config precedence, validation, redaction, and runtime-facing config boundaries even though runtime execution is out of scope. |
| Observability | Applies to the minimal logging and diagnostics foundation selected for this sprint; durable logs and event persistence remain excluded. |
| Security | Applies to workspace path safety, path escape rejection, secret redaction, and safe diagnostic output. |
| Testing | Applies to deterministic offline unit and command tests for workspace/config/health behavior. |
| Documentation | Applies to user-facing CLI help and maintainable sprint artifacts for the new commands. |
| CLI Surface | Applies to command shape, flags, help output, JSON/text output modes, and deterministic exit codes for `init-workspace`, `config show`, and `health`. |

## Selected Evidence Reports

Copied from the project index's "Available Evidence Reports" table. These tell the technical handbook which reports to read – the project index is the authoritative source.

| Report | Path | Covers |
| ------ | ---- | ------ |
| `01-project-structure` | `studies/go-cli-study/reports/final/01-project-structure.md` | Project layout, cmd/internal/pkg, dependency direction, thin entrypoints |
| `02-command-architecture` | `studies/go-cli-study/reports/final/02-command-architecture.md` | Command routing, flags, help text, command organization, shell completion |
| `03-dependency-injection` | `studies/go-cli-study/reports/final/03-dependency-injection.md` | Dependency construction, seams, testability, constructor patterns |
| `04-configuration-management` | `studies/go-cli-study/reports/final/04-configuration-management.md` | Config loading, precedence, environment variables, paths |
| `05-error-handling` | `studies/go-cli-study/reports/final/05-error-handling.md` | Error wrapping, classification, user-facing diagnostics, sentinel errors |
| `06-io-abstraction` | `studies/go-cli-study/reports/final/06-io-abstraction.md` | Filesystem/stdin/stdout abstraction, test seams, interface design |
| `09-terminal-ux` | `studies/go-cli-study/reports/final/09-terminal-ux.md` | Terminal output, progress indicators, human UX, color and formatting |
| `10-logging-observability` | `studies/go-cli-study/reports/final/10-logging-observability.md` | Logs, diagnostics, structured events, observability patterns |
| `11-testing-strategy` | `studies/go-cli-study/reports/final/11-testing-strategy.md` | Unit, integration, fixture, command-level tests, coverage strategy |
| `13-security` | `studies/go-cli-study/reports/final/13-security.md` | Path safety, secrets, command injection risks, sandboxing |

## Selected Reasoning Templates

All paths must appear in the project index's "Available Reasoning Templates" table.

| Template | Output Path | Why Selected |
| -------- | ----------- | ------------ |
| Architecture | `projects/ultraplan-go/sprints/02-workspace-config-health/reasoning/architecture.md` | Reason through package ownership, dependency direction, and composition boundaries for workspace/config/logging/output/health without making implementation decisions in this index. |

## Prior Decisions To Carry Forward

All decision paths must appear in the project index's "Prior Decisions" table.

| Decision | Path | Constraint For This Sprint |
| -------- | ---- | -------------------------- |
| None yet | None | The project index records no prior decisions; carry forward the project-level deferral that target scaffolding, sprint planning, sprint execution, and hosted/browser/multi-user surfaces are out of scope. |

## Required Review Protocols

All paths must appear in the project index's "Review Protocols" table.

| Protocol | Path | Required Evidence |
| -------- | ---- | ----------------- |
| Architecture Review | `system/protocols/architecture-review-protocol.md` | Evidence that workspace/config/logging/output/health package ownership, dependency direction, platform/product separation, and thin `cmd/ultraplan` entrypoint are respected. |
| Sprint Review | `system/protocols/sprint-review-protocol.md` | Evidence that Sprint 2 acceptance criteria, tests, build commands, output modes, redaction, exit codes, and explicit non-goals are satisfied or called out. |

## Excluded Context

| Context | Reason Excluded | Revisit If |
| ------- | --------------- | ---------- |
| `07-state-context` | Durable run state, context propagation for runtime tasks, and cancellation-heavy run-loop behavior are outside the workspace/config/health skeleton. | Sprint scope includes stateful run-loop, resumability, or cancellation behavior. |
| `08-concurrency` | Worker pools and parallel execution are not required for workspace discovery, config inspection, logging foundation, or health skeleton. | Sprint scope includes batch study execution or concurrent runtime orchestration. |
| `12-extensibility` | Plugin and extension-point design is not part of this sprint's workspace/config/health outputs. | Sprint scope includes public extension seams or additional runtime adapters. |
| `14-performance` | Performance optimization is not a selected driver beyond deterministic local workspace/config operations. | Sprint scope includes large workspace scanning, startup performance budgets, or large artifact handling. |
| `15-philosophy` | Cross-cutting philosophy is not specific evidence for the technical handbook compared with the selected implementation-focused reports. | Later sprint reasoning needs broad tradeoff framing beyond the selected contracts and reports. |
| LLM Runtime contract | Runtime execution, agentwrap integration, OpenCode, provider/model availability checks, and agent health checks are sprint non-goals. | A sprint implements runtime adapter behavior or OpenCode-backed health checks. |
| LLM Evaluation / Cost / Safety contract | Cost metadata, runtime safety evaluation, and provider/model behavior are outside this sprint's skeleton. | A sprint implements runtime execution, cost reporting, or evaluation discipline. |
| Workflows contract | Stateful workflow execution, retries, resumability, and orchestration are not part of this sprint. | A sprint implements run-loop, scheduling, retries, or resumable study execution. |
| Performance contract | Bounded workers, large artifact handling, and startup latency optimization are not central to this sprint's required outputs. | A sprint implements scheduler/concurrency/report processing or explicit performance requirements. |
| Persistence And Migrations contract | Persistent workspace state, migrations, and run state are explicitly excluded beyond initial scaffold creation. | A sprint introduces durable state, schema versions, migrations, or atomic run-state writes. |
| Study, source, analysis, synthesis, summary, and code extraction behavior | Sprint requirements explicitly exclude study listing, study initialization, source discovery, analysis runs, runtime execution, and code extraction. | Sprint scope moves from workspace/config/health foundation into study-side product workflows. |
| Target scaffolding, sprint planning, and sprint execution workflows | Project PRD/TRD and sprint requirements explicitly defer these product areas. | The PRD/TRD are revised and a later sprint selects target or sprint workflow scope. |

## Next Artifacts

- `technical-handbook.md` reads from the evidence reports listed above.
- `reasoning/*.md` captures area-specific reasoning.
- `reasoning.md` makes final sprint decisions.
- `plan.md` executes `reasoning.md`.
- `review.md` runs the selected review protocols against implementation.
