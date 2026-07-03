# Sprint Index: Documentation and Packaging

> Project: `ultraplan-go`
> Sprint: `15-documentation-and-packaging`
> Purpose: selected context for this sprint. Must be a subset of `projects/ultraplan-go/project-index.md`.
> **Inputs Used:** `projects/ultraplan-go/sprints/15-documentation-and-packaging/requirements.md`, `projects/ultraplan-go/project-index.md`, `projects/ultraplan-go/docs/PRD.md`, `projects/ultraplan-go/docs/TRD.md`, `projects/ultraplan-go/docs/ARCHITECTURE.md`, `templates/sprint-index.md`

This document selects what must be read, distilled, reasoned through, or checked for this sprint. It does not make implementation decisions. All selections must come from the project index — no items may be included that are not listed in the project index.

## Sprint Scope

- **Sprint Goal:** Prepare the first production-ready study-side release by completing user/operator documentation, packaging Linux and macOS CLI builds with checksums, and recording offline plus gated OpenCode smoke-release evidence.
- **Planned Output:** User guide, configuration guide, recovery runbook, OpenCode smoke instructions, release checklist, CLI command reference, README release update, four release binaries, checksums, and smoke evidence.
- **Depends On:** Sprint 1 CLI shell and app composition; Sprint 2 workspace/config/health; Sprint 3 study domain/listing; Sprint 4 study initialization; Sprint 5 Markdown sources/applicability; Sprint 6 report validation/rating parsing; Sprint 7 prompt composition; Sprint 8 run-state/status; Sprint 9 agentwrap/OpenCode runtime integration; Sprint 10 single analysis/synthesis; Sprint 11 run-all batch execution; Sprint 12 durable run-loop; Sprint 13 summary/code extraction; Sprint 14 validation/diagnostics/JSON stability; External OpenCode environment when available for gated smoke.
- **Non-Goals:** New product behavior beyond documentation, packaging scripts or commands if needed, and release verification artifacts; target scaffolding; sprint planning; sprint execution; target/sprint validators; target/sprint runtime task kinds; hosted service; browser UI; multi-user auth; organization permissions; collaboration features; new runtime adapters beyond agentwrap/OpenCode; mandatory real OpenCode/provider smoke tests in default test suites; schema migrations unless already required by Sprint 14; GitHub release publication, tags, signing, notarization, or artifact upload unless separately requested.

## Source Project Index

- `projects/ultraplan-go/project-index.md` — authoritative source. Any file or item referenced below must appear there.

## Selected Contracts

Each contract applies as a flat whole to this sprint. All paths must appear in the project index's "Active Contract Pool" table.

| Contract | Why Selected |
| ------------ | -------------------------------------------- |
| Architecture | Release docs and packaging must preserve module ownership, dependency direction, thin entrypoints, and product/platform separation from `docs/ARCHITECTURE.md`. |
| Errors | Recovery docs, CLI reference, smoke evidence, and release checks must describe actionable diagnostics, exit codes, task failures, and validation failures accurately. |
| Configuration | Configuration guide, OpenCode smoke instructions, health/config docs, redaction behavior, precedence, runtime mapping, and schema-version rejection are core sprint outputs. |
| Observability | Smoke evidence, status docs, diagnostics, retry/fallback metadata, run/task metadata, and health/preflight truthfulness are acceptance-critical. |
| Security | Documentation and release artifacts must avoid secrets, preserve path/source isolation, describe permission posture, and avoid unsafe runtime payload disclosure. |
| Testing | Offline release gates, race tests, fake-first verification, gated integration expectations, and smoke evidence all depend on the testing contract. |
| Documentation | This sprint is primarily documentation and release-readiness; docs must be user-facing, maintainable, accurate, and aligned with help/output surfaces. |
| CLI Surface | CLI command reference, README quickstart, output modes, help behavior, stable JSON surfaces, and exit/error expectations must reflect the public CLI surface. |
| LLM Runtime | OpenCode smoke instructions and runtime documentation must keep integration through agentwrap/OpenCode and avoid product-owned direct OpenCode supervision claims. |
| LLM Evaluation / Cost / Safety | Runtime metadata, cost/usage handling where available, safe evidence, and real-runtime smoke disclosure must remain safety-conscious and audit-ready. |
| Workflows | Run-all, run-loop, status, retries, cancellation, resumability, stale locks, and recovery guidance depend on workflow semantics. |
| Performance | Release docs and checks should avoid promising unsupported behavior while preserving bounded workers, startup expectations, and large artifact handling guidance. |
| Persistence And Migrations | Recovery docs, run-state behavior, atomic writes, schema rejection, lock handling, durable artifacts, and release checklist review are persistence-sensitive. |

## Selected Evidence Reports

Copied from the project index's "Available Evidence Reports" table. These tell the technical handbook which reports to read – the project index is the authoritative source.

| Report | Path | Covers |
| ----------- | ------------------------------------------------------ | ------------------------------- |
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
| `12-extensibility` | `studies/go-cli-study/reports/final/12-extensibility.md` | Extension points, plugin architecture, package boundaries, API design |
| `13-security` | `studies/go-cli-study/reports/final/13-security.md` | Path safety, secrets, command injection risks, sandboxing |
| `14-performance` | `studies/go-cli-study/reports/final/14-performance.md` | Startup latency, large repos, memory management, performance |
| `15-philosophy` | `studies/go-cli-study/reports/final/15-philosophy.md` | Cross-cutting design philosophy, tradeoffs, maintainability |

## Selected Reasoning Templates

All paths must appear in the project index's "Available Reasoning Templates" table.

| Template | Output Path | Why Selected |
| ------------ | -------------------------------------------------------------- | ------------ |
| Architecture | `projects/ultraplan-go/sprints/15-documentation-and-packaging/reasoning/architecture.md` | Confirm documentation, packaging, release artifacts, and any packaging helper code preserve module boundaries, dependency direction, runtime/product separation, and the single CLI entrypoint. |

## Prior Decisions To Carry Forward

All decision paths must appear in the project index's "Prior Decisions" table.

| Decision | Path | Constraint For This Sprint |
| ------------ | -------- | -------------------------- |
| None yet — this is the first project. | N/A | No prior decision records are available in the project index; carry forward only the project-index maintenance notes, PRD/TRD deferred scope, and sprint dependencies listed in `requirements.md`. |

## Required Review Protocols

All paths must appear in the project index's "Review Protocols" table.

| Protocol | Path | Required Evidence |
| ------------ | --------------------------------------- | ----------------- |
| Architecture Review | `system/protocols/architecture-review-protocol.md` | Evidence that documentation, packaging, and release verification do not violate module boundaries, runtime/product separation, the single CLI entrypoint, or study-side scope. |
| Sprint Review | `system/protocols/sprint-review-protocol.md` | Final review of required outputs, acceptance criteria, release artifacts, checksums, smoke evidence, offline gates, gated OpenCode smoke disposition, and excluded non-goals. |

## Excluded Context

| Context | Reason Excluded | Revisit If |
| ----------- | --------------- | ------------- |
| Target scaffolding and target commands | Explicitly deferred by PRD/TRD, project index non-goals, and sprint requirements; README/docs must not claim support. | PRD/TRD are revised to include target workflows in current release scope. |
| Sprint planning and sprint execution workflows | Explicitly deferred by PRD/TRD and sprint non-goals; this sprint is about study-side release documentation and packaging only. | A future release adds sprint planning/execution to the active product scope. |
| Hosted SaaS, browser UI, multi-user auth, organization permissions, and collaboration features | Explicit project and sprint non-goals; release artifacts are local CLI outputs only. | Product scope changes to include hosted or collaborative operation. |
| Additional runtime adapters beyond agentwrap/OpenCode | Sprint constraints require the existing agentwrap/OpenCode path and forbid new runtime adapters. | A later sprint adds a new runtime adapter contract and implementation scope. |
| Mandatory real OpenCode/provider smoke tests in normal test suites | Normal verification must stay fake-first and offline; real OpenCode smoke is gated and may be skipped with evidence. | CI or release policy explicitly provisions real runtime credentials and opts into gated smoke. |
| GitHub release publication, tags, binary signing, macOS notarization, and artifact upload | Sprint non-goals exclude publication and distribution-platform actions beyond local packaging artifacts. | A separate release-publishing request is made. |
| Direct OpenCode process supervision by UltraPlan | TRD and project index require using `github.com/Antonio7098/agentwrap` and `agentwrap/opencode`; UltraPlan must not document product-owned direct OpenCode supervision. | Agentwrap integration requirements are replaced by a revised TRD. |
| Raw prompts, full report bodies, provider tokens, full sensitive environment variables, and unsafe runtime payloads in release evidence | Sprint constraints require release artifacts and docs to avoid secrets and unsafe payload disclosure. | A secure debug-retention policy is defined and explicitly requested for non-release diagnostics. |

## Next Artifacts

- `technical-handbook.md` reads from the evidence reports listed above.
- `reasoning/*.md` captures area-specific reasoning.
- `reasoning.md` makes final sprint decisions.
- `plan.md` executes `reasoning.md`.
- `review.md` runs the selected review protocols against implementation.
