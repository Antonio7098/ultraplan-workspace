# Project Index: ultraplan-go

> Project: `ultraplan-go`
> Purpose: available governance, evidence, and reasoning pool for UltraPlan Go implementation.

## Project Scope

- **Project Slug:** `ultraplan-go`
- **Repository:** `../ultraplan-go/`
- **Target Implementation Directory:** `/home/antonioborgerees/coding/ultraplan/ultraplan-go`
- **Primary Goal:** Build a production-grade Go CLI, local TUI, and loopback-only Go-served browser UI for UltraPlan study workflows, governed project/sprint planning and execution, automated conformance review, deep smoke through `smoke`, and controlled Git publication after valid stage completion.
- **Phase 1 Goal:** Study initialization, source analysis, synthesis, code-reference extraction, resumable orchestration, validation, and diagnostics.
- **Phase 2 Goal:** Project cataloging plus sprint planning and execute artifacts: `requirements.md`, `sprint-index.md`, `technical-handbook.md`, `reasoning/*.md`, `reasoning.md`, `plan.md`, `execute.md`, `flow-state.json`, `.run-state.json`, and configurable global/per-stage models for sprint stages.
- **Phase 3 Goal:** Automated post-execute `review.md`, sprint-targeted `smoke.md`, external smoke-harness evidence, review-before-smoke flow integration, and complete CLI/TUI operation through smoke.
- **Phase 4 Goal:** Beginning with Sprint 30, add `ultraplan serve`, a loopback-only Go HTTP server, Go-rendered browser dashboard, guarded HTTP operations, SSE progress, and browser recovery over the same app use cases and workspace state.
- **Non-Goals:** General-purpose issue tracking, automatic product fixes, hosted SaaS, remote exposure, multi-user collaboration, remote workers, cross-sprint verification scheduling, and general Git automation beyond stage-owned commit and push are explicitly deferred per PRD.
- **Documentation Source Of Truth:** `projects/ultraplan-go/docs/` in this planning workspace is the sole authoritative location for the PRD, TRD, and Architecture documents. The implementation repository does not carry duplicate mirrors.

## Source Documents

| Document | Path | Summary |
|---|---|---|
| Product Requirements | `projects/ultraplan-go/docs/PRD.md` | Product goals, user scenarios, CLI/TUI/local-web surfaces, study workflows, planning and execute workflows, review, smoke, validation, and launch criteria. |
| Technical Requirements | `projects/ultraplan-go/docs/TRD.md` | Go architecture, workspace/config/runtime/state/validation requirements, planning/execute/review/smoke requirements, agentwrap integration, local HTTP/SSE boundary, external smoke-harness boundary, and testing requirements. |
| Architecture | `projects/ultraplan-go/docs/ARCHITECTURE.md` | Module-driven package layout and dependency rules for `study`, `project`, `sprint`, `workspace`, `codeextract`, TUI, web, and platform runtime/process packages. |

## Active Contract Pool

Contracts are selected per sprint and applied through the phase gates in `roadmap.md`. A contract being in the active pool does not mean every production requirement blocks every skeleton sprint.

| Contract | Path | Applies To | Selection Notes |
|---|---|---|---|
| Architecture | `system/contracts/core/architecture.md` | All sprints | Module ownership, dependency direction, thin entrypoints, product/platform separation. |
| Errors | `system/contracts/core/errors.md` | All implementation sprints | Broad error categories, wrapping, actionable CLI diagnostics, exit codes, task failures. |
| Configuration | `system/contracts/core/configuration.md` | Workspace/config/runtime config sprints | Config precedence, command preflight validation, redaction, runtime mapping. |
| Observability | `system/contracts/core/observability.md` | Runtime, run-loop, status, logs | Events, diagnostics, structured logs, run/task metadata, health/preflight truthfulness. |
| Security | `system/contracts/core/security.md` | Runtime execution, paths, permissions | Workspace path safety, secret redaction, subprocess/runtime policy, source isolation. |
| Testing | `system/contracts/core/testing.md` | All implementation sprints | Unit/fixture/fake-runtime/gated integration expectations. |
| Documentation | `system/contracts/core/documentation.md` | CLI/help/artifact docs | Generated docs, user-facing help, maintainable artifacts. |
| CLI Surface | `system/contracts/surfaces/cli.md` | CLI command sprints | Command shape, flags, help, exit codes, JSON/text output. |
| LLM Runtime | `system/contracts/runtime/llm.md` | Agentwrap/OpenCode integration | Runtime behavior, provider/model boundaries, adapter use. |
| LLM Evaluation / Cost / Safety | `system/contracts/runtime/llm-evaluation-cost-safety.md` | Runtime execution/validation | Cost metadata, safety, evaluation discipline. |
| Workflows | `system/contracts/runtime/workflows.md` | Run-loop, orchestration, retries | Stateful batch/workflow execution, retries, cancellation, resumability. |
| Performance | `system/contracts/runtime/performance.md` | Scheduler/concurrency/report processing | Bounded workers, startup latency, repository scans, large artifact handling. |
| Persistence And Migrations | `system/contracts/runtime/persistence-and-migrations.md` | Run state/workspace artifacts | Atomic file writes, durable format versions, compatibility and migrations. |

## Phase 2 Planning And Execute Context

Planning-side sprints should select only the contracts and evidence needed for the artifact or execute stage being implemented. The phase workflow is:

```text
study -> select -> distill -> reason -> plan -> execute
```

Planning-side modules must follow these ownership rules:

- `internal/project` owns project docs and `project-index.md` catalog behavior.
- `internal/sprint` owns planning artifacts, execute artifacts, `flow-state.json`, and `.run-state.json` through `execute`.
- `internal/study` remains independent and does not become a shared planning abstraction.
- `internal/platform/runtime` remains generic prompt execution infrastructure.

Phase 2 may reuse workspace discovery, config/redaction, command output conventions, generic runtime execution, and atomic file/JSON writes. It must not abstract study services, source/dimension models, report validation, rating parsing, summary generation, or run-loop scheduling for planning or execute behavior.

## Phase 3 Review And Smoke Context

Phase 3 extends the sprint workflow without adding a second artifact hierarchy:

```text
study -> select -> distill -> reason -> plan -> execute -> review -> smoke
```

- `internal/sprint` owns review and smoke stage semantics, validation, verdicts, `review.md`, `smoke.md`, and flow-state integration.
- `internal/platform/runtime` remains generic and executes independent structured review requests through agentwrap.
- The configured smoke author model first builds or updates a durable sprint-specific suite inside manifest-declared harness authoring paths; a generic external-process boundary then discovers and invokes it with explicit argv, bounded environment forwarding, timeout, and cancellation.
- `review.md` is the current automated sprint conformance review and may replace the older manually produced file.
- `smoke.md` is the current sprint smoke summary and may replace older `deep-smoke.md` files.
- Raw smoke run JSON, stdout/stderr captures, test artifacts, and open/resolved issue records stay in the smoke harness directories cataloged below.
- Review runs before smoke. Blocking/high review findings stop default smoke execution.
- Review and smoke are available through the CLI and the TUI using the same typed app use cases.
- Review and smoke agents must not modify product source, product tests, governed planning inputs, or Git state. After valid canonical output is persisted, UltraPlan may commit and push only the paths owned by that completed stage.

## Phase 4 Local Web Context

Phase 4 starts with Sprint 30 and adds a local interface rather than another workflow stage:

```text
browser -> internal/web HTTP/SSE -> internal/app use cases -> existing product modules
```

- `internal/web` owns loopback HTTP lifecycle, Go templates/static assets, transport mapping, browser security, confirmation tokens, and bounded SSE subscriptions.
- `internal/web` must not import product modules directly, call CLI handlers, parse CLI output, or persist an alternate workflow state.
- Go `html/template`, embedded CSS, and minimal JavaScript are the initial UI stack; Node.js, Vite, and a separate frontend process are not required.
- HTTP POST/DELETE carries commands and cancellation; SSE carries one-way progress only.
- Workspace artifacts and product-owned run state remain authoritative across browser refresh, disconnect, and server restart.
- Hosted service, remote binding, accounts, permissions, tenants, collaboration, remote workers, and WebSockets remain deferred.

## Available Studies

| Study | Path | Useful For | Status |
|---|---|---|---|
| `go-cli-study` | `studies/go-cli-study/` | Go CLI structure, command architecture, config, testing, error handling, concurrency, observability, security, performance, extensibility, philosophy | Current |

## Available Evidence Reports

| Report | Path | Covers |
|---|---|---|
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

## Available Reasoning Templates

| Template | Path | Useful For | Status |
|---|---|---|---|
| Architecture | `system/reasoning/architecture_reasoning_template.md` | Module boundaries, dependency direction, package layout, runtime/product separation | Current |
| API Design | `system/reasoning/api-design-reasoning-template.md` | Phase 4 HTTP resources, versioning, errors, confirmation, cancellation, and SSE semantics | Current |
| Frontend | `system/reasoning/frontend-reasoning-template.md` | Phase 4 server-rendered browser flows, state placement, progressive enhancement, and accessibility | Current |

## Smoke Harnesses

| Harness | Path | Manifest | Evidence | Useful For | Status |
|---|---|---|---|---|---|
| `ultraplan-go-smoke` | `/home/antonioborgerees/coding/ultraplan/ultraplan-go-smoke/` | `ultraplan-smoke.json` | `runs/` and `issues/` under the harness root | Existing real-runtime evidence for `ultraplan-go`, including OpenCode execution, CLI/TUI behavior, diagnostics, persisted state, cancellation, security, redaction, and sprint-specific deep-smoke suites. `smoke.md` links the relevant run IDs and evidence paths. | Harness current; versioned Phase 3 manifest planned in Sprint 27 |

## Prior Decisions

None yet — this is the first project.

## Review Protocols

| Protocol | Path | Required When |
|---|---|---|
| Architecture Review | `system/protocols/architecture-review-protocol.md` | Package layout, dependency boundaries, runtime/module separation |
| Sprint Review | `system/protocols/review-sprint-protocol.md` | Every completed implementation sprint in Phase 3; writes the current `review.md` |
| Deep Smoke Sprint | `system/protocols/deep-smoke-sprint-protocol.md` | Completed sprints that need real-runtime smoke evidence |

## Maintenance Notes

- Keep this index as a catalog, not a sprint plan.
- All 15 go-cli-study reports are available — select only the ones relevant to the sprint's scope via `sprint-index.md`.
- Reasoning templates are added on-demand, not upfront. Start with Architecture only.
- Planning-side sprints may implement project and sprint artifact workflows through `execute`; Phase 3 extends the same sprint module through `review` and `smoke`.
- Keep Phase 3 sprint artifacts simple: current `review.md` and `smoke.md` in the sprint root. Do not add a parallel verification directory or copy raw harness evidence into the project workspace.
- The external smoke harness owns durable authored suites and detailed `runs/` and `issues/` evidence. UltraPlan owns agent-driven sprint-suite authoring, authoring-scope validation, enumerated coverage validation, smoke selection/invocation, summary validation, flow state, and the sprint-root `smoke.md` link summary.
- Verification findings are evidence inside `review.md`, `smoke.md`, and the harness issue records; do not turn Phase 3 into a general-purpose issue tracker.
- Every Phase 3 sprint must update the TUI for the CLI/use-case functionality introduced by that sprint.
- Product Phase 4 starts at Sprint 30 and keeps the browser UI local, Go-rendered, progressively enhanced, and backed by shared app use cases.
- Built-in prompts and output templates ship embedded in UltraPlan. Workspace `prompts/` and `templates/` paths are optional intentional overrides, never required project inputs.
- Agents may not mutate Git. UltraPlan may commit and push stage-owned paths after valid completion when `git.stage_completion` enables publication.
- TRD requires `github.com/Antonio7098/agentwrap` and `agentwrap/opencode` as the runtime SDK. Do not invent competing runtime contracts.
