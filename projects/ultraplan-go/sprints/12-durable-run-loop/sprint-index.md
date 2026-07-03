# Sprint Index: 12 Durable Run Loop

> Project: `ultraplan-go`
> Sprint: `12-durable-run-loop`
> Purpose: selected context for this sprint. Must be a subset of `projects/ultraplan-go/project-index.md`.
> **Inputs Used:** `projects/ultraplan-go/sprints/12-durable-run-loop/requirements.md`, `projects/ultraplan-go/project-index.md`, `projects/ultraplan-go/docs/PRD.md`, `projects/ultraplan-go/docs/TRD.md`, `projects/ultraplan-go/docs/ARCHITECTURE.md`, `templates/sprint-index.md`

This document selects what must be read, distilled, reasoned through, or checked for this sprint. It does not make implementation decisions. All selections must come from the project index - no items may be included that are not listed in the project index.

## Sprint Scope

- **Sprint Goal:** Implement `ultraplan study <study> run-loop` as durable, resumable study orchestration with atomic state updates, bounded execution, retry/fallback metadata, cancellation handling, stale-task recovery, and per-study locking.
- **Planned Output:** Sprint artifacts plus durable run-loop implementation, lock handling, run-state domain/persistence/resume updates, study service and CLI wiring, status updates, and fake-runtime/state/lock/CLI test coverage.
- **Depends On:** Sprint 5 source discovery/applicability, Sprint 6 validation/rating parsing, Sprint 7 prompt composition, Sprint 8 run-state/status primitives, Sprint 9 agentwrap/OpenCode runtime integration, Sprint 10 single analysis/synthesis, Sprint 11 bounded `run-all` patterns, and Sprint 11 review follow-ups named by the sprint requirements.
- **Non-Goals:** Code-reference extraction, target scaffolding, sprint planning, sprint execution, generic workflow engines, new runtime adapters, real OpenCode smoke tests in default tests, stable public JSON schemas, broad migrations, hosted/browser/multi-user features, and destructive filesystem or Git mutations outside the selected study state/lock and report artifact paths.

## Source Project Index

- `projects/ultraplan-go/project-index.md` - authoritative source. Any file or item referenced below must appear there.

## Selected Contracts

Each contract applies as a flat whole to this sprint. All paths must appear in the project index's "Active Contract Pool" table.

| Contract | Why Selected |
| ------------ | -------------------------------------------- |
| Architecture | Enforces study-owned orchestration, thin CLI wiring, product/platform separation, and no global scheduler/workflow/prompts/reports packages. |
| Errors | Required for actionable resume/load/write/runtime diagnostics, cancellation exit behavior, classified task failures, and safe cause preservation. |
| Configuration | Required for resolved parallelism, retries, timeout, runtime/model overrides, permission/runtime mapping, and redacted config diagnostics. |
| Observability | Required for durable task metadata, runtime-free status, retry/fallback/cost/usage diagnostics, and truthful progress reporting. |
| Security | Required for workspace-contained state/lock paths, secret redaction, safe runtime diagnostics, permission posture, and source/document isolation. |
| Testing | Required for fake-runtime, fixture, command, atomic-write, lock, cancellation, race, and offline verification coverage. |
| Documentation | Required for command help, user-facing flags, recovery guidance, sprint artifacts, and concise deterministic output documentation. |
| CLI Surface | Required for `ultraplan study <study> run-loop`, status rendering, flags, help, exit codes, and human/scriptable output behavior. |
| LLM Runtime | Required because run-loop execution must go through agentwrap/OpenCode runtime boundaries without direct product-owned OpenCode invocation. |
| LLM Evaluation / Cost / Safety | Required for validation discipline, usage/cost preservation, fallback safety, and safe handling of runtime metadata. |
| Workflows | Required for stateful batch execution, resumability, retries, cancellation, synthesis unblocking, and durable terminal states. |
| Performance | Required for bounded workers, no unbounded goroutines/channels, efficient runtime-free status, and large task-matrix handling. |
| Persistence And Migrations | Required for atomic state writes, schema compatibility or explicit rejection, durable format versioning, locks, and prior-state preservation on write failure. |

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
| `07-state-context` | `studies/go-cli-study/reports/final/07-state-context.md` | Context propagation, app state, cancellation |
| `08-concurrency` | `studies/go-cli-study/reports/final/08-concurrency.md` | Worker pools, cancellation, parallel execution, goroutine management |
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
| Architecture | `projects/ultraplan-go/sprints/12-durable-run-loop/reasoning/durable-run-loop.md` | Reason through study ownership of run-loop state machines, lock semantics, runtime/product separation, bounded scheduling, persistence boundaries, and status diagnostics. |

## Prior Decisions To Carry Forward

All decision paths must appear in the project index's "Prior Decisions" table.

| Decision | Path | Constraint For This Sprint |
| ------------ | -------- | -------------------------- |
| None indexed | Not listed | The project index records no prior-decision rows. Sprint 11 follow-ups are carried only through the sprint requirements dependencies and acceptance criteria, not as selectable prior-decision artifacts. |

## Required Review Protocols

All paths must appear in the project index's "Review Protocols" table.

| Protocol | Path | Required Evidence |
| ------------ | --------------------------------------- | ----------------- |
| Architecture Review | `system/protocols/architecture-review-protocol.md` | Confirm run-loop behavior lives in `internal/study`, CLI wiring remains thin, `internal/platform/runtime` has no product imports, and no deferred global workflow or technical-layer packages were introduced. |
| Sprint Review | `system/protocols/sprint-review-protocol.md` | Record acceptance evidence, verification commands, deviations, residual risks, and carry-forward follow-ups for the completed durable run-loop sprint. |

## Excluded Context

| Context | Reason Excluded | Revisit If |
| ----------- | --------------- | ------------- |
| `12-extensibility` evidence report | Plugin architecture and broad extension-point design are not needed for this sprint's bounded study-owned run-loop and are adjacent to explicitly deferred plugin/workflow scope. | A later sprint introduces plugin architecture, external extension APIs, or new runtime adapter extension points. |
| Code-reference extraction | Sprint requirements explicitly exclude `ultraplan code <report>...` and code-reference extraction implementation. | A code extraction sprint selects the `codeextract` scope from the PRD/TRD. |
| Target and sprint command families | PRD, TRD, project index, and sprint requirements defer target scaffolding, sprint planning, and sprint execution. | Project requirements are revised to include target/sprint workflows. |
| Stable public JSON schemas | Sprint requirements assign JSON stability to Sprint 14 and exclude stable public schemas for status, run-loop, validation, and diagnostics. | Sprint 14 or an explicit schema-stability sprint begins. |
| Real OpenCode smoke tests in default verification | Sprint requirements require normal tests to use fake runtimes and local fixtures without OpenCode, credentials, network, or real subprocess execution. | Gated integration testing or release smoke validation is explicitly selected. |
| New runtime adapters or direct provider workers | PRD/TRD and sprint requirements require existing agentwrap/OpenCode integration and exclude new adapters or direct provider workers. | A future runtime adapter sprint is selected from updated requirements. |
| Generic workflow engine, DAG authoring, TUI, hosted service, browser UI, and multi-user collaboration | These are PRD/TRD non-goals or future-scope items and would dilute the durable study run-loop implementation. | Product scope changes to include hosted, UI, collaboration, or general workflow capabilities. |
| Broad migration framework | Sprint requirements allow only explicit compatibility or rejection behavior for the current run-state schema. | A migration-focused sprint is selected or persisted public formats require compatibility guarantees. |
| Destructive filesystem or Git mutations outside selected study state/lock and reports | Sprint requirements constrain mutations to workspace-contained study state/lock areas and expected artifact paths. | A future command explicitly owns broader destructive behavior with safety rules. |

## Next Artifacts

- `technical-handbook.md` reads from the evidence reports listed above.
- `reasoning/*.md` captures area-specific reasoning.
- `reasoning.md` makes final sprint decisions.
- `plan.md` executes `reasoning.md`.
- `review.md` runs the selected review protocols against implementation.
