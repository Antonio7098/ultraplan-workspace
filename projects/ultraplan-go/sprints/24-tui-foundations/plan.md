# Sprint Plan: TUI Foundation and Read-Only Dashboard

> Project: `ultraplan-go`
> Sprint: `24-tui-foundations`
> Source: `reasoning.md`
> **Inputs Used:** `projects/ultraplan-go/project-index.md`, `projects/ultraplan-go/sprints/24-tui-foundations/requirements.md`, `projects/ultraplan-go/docs/ARCHITECTURE.md`, `projects/ultraplan-go/docs/PRD.md`, `projects/ultraplan-go/docs/TRD.md`, `projects/ultraplan-go/sprints/24-tui-foundations/sprint-index.md`, `projects/ultraplan-go/sprints/24-tui-foundations/technical-handbook.md`, `projects/ultraplan-go/sprints/24-tui-foundations/reasoning/architecture.md`, `projects/ultraplan-go/sprints/24-tui-foundations/reasoning.md`

This plan executes `reasoning.md`. It must not invent architecture, scope, or decisions.

## Reasoning Source

- **Sprint Reasoning:** `reasoning.md`
- **Sprint Index:** `sprint-index.md`
- **Technical Handbook:** `technical-handbook.md`
- **Area Reasoning:** `reasoning/architecture.md`

## Sprint Status

- **Status:** `complete`
- **Owner:** `implementation agent`
- **Start Date:** `2026-07-09`
- **Completion Date:** `2026-07-09`

## Decisions To Execute

| Decision | Source Section | Execution Implication |
| --- | --- | --- |
| Add Thin `ultraplan tui` Command Startup Through `internal/app` | `reasoning.md#decision-1-add-thin-ultraplan-tui-command-startup-through-internalapp` | Add `internal/app/tui_commands.go` command wiring that parses `--workspace`, reuses workspace discovery and normal non-runtime config validation, constructs read-only app use cases, launches the TUI, and does not initialize runtime health checks, OpenCode, provider credentials, network access, CLI handlers, stdout parsing, or self-subprocess integration. |
| Create Shared Read-Only App Use Cases For CLI And TUI | `reasoning.md#decision-2-create-shared-read-only-app-use-cases-for-cli-and-tui` | Add typed read-only result/use-case files in `internal/app` for project, study, sprint, validation findings, artifact paths, and display-safe status summaries; refactor CLI adapters only as needed so CLI and TUI call shared use cases rather than duplicate product workflows. |
| Keep TUI State Derived, Explicitly Refreshed, Context-Aware, And Disposable | `reasoning.md#decision-3-keep-tui-state-derived-explicitly-refreshed-context-aware-and-disposable` | Model only derived UI/session state, use immutable startup workspace/config snapshots, pass `context.Context` through startup and refresh, avoid background periodic refresh, avoid hot config reload, and introduce no TUI persistence, cache, or database. |
| Implement A Contained Testable TUI Model, Rendering, And Key-Binding Layer | `reasoning.md#decision-4-implement-a-contained-testable-tui-model-rendering-and-key-binding-layer` | Implement `internal/tui/app.go`, `model.go`, `views.go`, and `keys.go` with deterministic model/update/render behavior, terminal-library containment inside `internal/tui`, navigable panes, findings, preview/error panes, key help, and narrow-terminal fallback. |
| Provide Path-Safe, Bounded, Read-Only Artifact Previews | `reasoning.md#decision-5-provide-path-safe-bounded-read-only-artifact-previews` | Add preview support for known Markdown and JSON artifacts using existing workspace-safe path resolution, product-root containment, bounded reads, explicit truncation, and actionable missing/path/read/format errors without arbitrary filesystem browsing. |
| Display Execute And Runtime Diagnostics Only Through Typed Redacted Status Results | `reasoning.md#decision-6-display-execute-and-runtime-diagnostics-only-through-typed-redacted-status-results` | Ensure execute status in the TUI comes from typed app/sprint status results backed by durable `flow-state.json`, `execute.md`, and `.run-state.json` interpretation outside `internal/tui`; display only redacted or display-safe runtime/config/model/provider/diagnostic fields. |
| Enforce A Strict Read-Only Trust Boundary And Reject Operational/Extension Surfaces | `reasoning.md#decision-7-enforce-a-strict-read-only-trust-boundary-and-reject-operationalextension-surfaces` | Limit key bindings and model actions to navigation, pane switching, selection, explicit refresh, artifact preview, and quit; do not expose validation runs, flow, prompt preview, execute, study run-loop, runtime-backed work, Git mutation, subprocess/plugin launchers, smoke/review/issue artifact generation, or automatic workspace mutation. |
| Require Deterministic Non-Interactive Test And Documentation Evidence | `reasoning.md#decision-8-require-deterministic-non-interactive-test-and-documentation-evidence` | Add deterministic tests for app use cases, TUI fakes, model updates, rendering, command startup, previews, errors, no-runtime startup, narrow terminals, redaction, and read-only guarantees; document the TUI dependency, ownership boundary, read-only scope, non-goals, and dependency containment in `internal/tui/doc.go`. |

## Requirements / Contracts To Satisfy

| Contract / Requirement ID | Required Behavior | Evidence Planned |
| --- | --- | --- |
| R1 / Required Outputs | Required files in `internal/app` and `internal/tui` are created or updated. | Source review of required paths from `requirements.md` lines 13-30. |
| AC1 | `ultraplan tui` is discoverable in help and starts inside a workspace using CLI workspace discovery/config validation rules. | `internal/app/tui_commands_test.go`; command help review. |
| AC2 | `ultraplan tui --workspace <path>` targets the explicit workspace and reports actionable invalid workspace/config errors. | Startup tests for explicit workspace, invalid workspace, invalid config, and error classification. |
| AC3 | TUI uses typed app/product results and never invokes `ultraplan` as subprocess, calls CLI command handlers, parses CLI stdout/stderr, or scrapes rendered text. | App use-case tests, TUI fake-use-case tests, and source review/search for CLI handler, stdout scraping, and `os/exec` self-invocation absence. |
| AC4 | Dashboard shows projects, studies, and sprints in navigable panes with deterministic ordering. | `internal/tui/model_test.go`, `internal/tui/views_test.go`, and use-case ordering tests. |
| AC5 | Project views show list, status, required docs, project-index health, catalog findings, and key artifact paths. | Project use-case tests and TUI rendering assertions. |
| AC6 | Study views show list, detail, source/dimension summary, status, validation findings where supported, and key artifact paths. | Study use-case tests and TUI model/view assertions. |
| AC7 | Sprint views show sprint list, planning artifact status, flow-state status through `execute`, execute run-state summary where available, validation findings where supported, and key artifact paths. | Sprint use-case tests, execute status/redaction tests, and TUI model/view assertions. |
| AC8 | Artifact preview supports required Markdown and JSON artifacts with bounded reads and clear missing-file errors. | Preview model/view tests for supported artifacts, missing files, truncation, path rejection, read errors, and display states. |
| AC9 | Dashboard starts and remains usable without runtime credentials, OpenCode availability, provider configuration beyond normal non-runtime config validation, or network access. | No-runtime startup tests and dependency review. |
| AC10 | Read-only TUI actions do not intentionally mutate workspace artifacts; any existing status operation that may refresh deterministic status files is explicitly labeled before invocation. | Key-binding/model tests, fake mutation counters, and source review for mutating action absence or explicit labels. |
| AC11 | Terminal rendering keeps critical status/error text visible in narrow terminals and does not rely on ANSI-only meaning. | `internal/tui/views_test.go` narrow-width cases with plain text assertions. |
| AC12 | Unit tests cover TUI navigation/model updates, views, startup command behavior, fakes, errors, no-runtime startup, and read-only guarantees. | Targeted tests in `internal/tui` and `internal/app`. |
| AC13 | Verification passes from `../ultraplan-go`: `go test ./...` and `go build ./cmd/ultraplan`. | Final verification command output. |
| Architecture | `internal/app` is the local interface boundary; `internal/tui` owns terminal UI only; product/platform packages do not import TUI; terminal-library types do not leak outside TUI. | Architecture review import/type checks and source review. |
| CLI Surface | Existing CLI help, exit codes, text output, and JSON surfaces remain stable except adding `tui`. | Existing command tests plus TUI command tests. |
| Configuration | Startup reuses workspace discovery/config precedence and redacts displayed config/runtime values. | `--workspace`, invalid config, no-runtime, and redaction tests. |
| Documentation | TUI dependency, ownership boundary, read-only scope, and non-goals are documented. | `internal/tui/doc.go`, command help, and key help review. |
| Errors | Startup, preview, workspace/config, and use-case failures are classified and actionable. | Command startup tests, model error tests, and view error-pane assertions. |
| Observability | Status, findings, execute summaries, and diagnostics are structured and safe to display without launching workflows. | Use-case tests, TUI rendering tests, and redaction review. |
| Performance | Startup, refresh, ordering, rendering, and previews avoid unbounded reads or expensive scans. | Bounded preview tests, deterministic ordering tests, and source review for no recursive scans. |
| Persistence And Migrations | Durable workspace artifacts remain source of truth; no TUI-specific persistence model is introduced. | Review showing no TUI cache/database/state files; read-only tests. |
| Security | Preview paths are safe; no self-subprocess, runtime credential dependency, network access, raw secrets, or unredacted diagnostics. | Path escape tests, no-runtime tests, redaction tests, and source review/search. |
| Testing | Normal tests are deterministic and non-interactive with fake app use cases and no real terminal/runtime/network. | `go test ./...` and targeted fake-use-case tests. |
| Workflows | TUI displays workflow state but does not execute planning, execute, study run-loop, validation, or runtime workflows. | Key-binding/model tests, use-case surface review, and source search. |

## Tasks

- [x] **Task 1: Add Shared Read-Only App Use-Case Boundary**
  > Executes: `Decision 2`, `Architecture`, `AC3`, `AC4-AC7`, `AC10`
  - [x] Add narrow typed read-only use-case interfaces/results in `internal/app/usecases.go` for local-interface dashboard needs without terminal-library types.
  - [x] Add project read-only use cases in `internal/app/project_usecases.go` for project list, status, required docs, project-index health, catalog findings, deterministic ordering, and key artifact paths.
  - [x] Add study read-only use cases in `internal/app/study_usecases.go` for study list, detail, source/dimension summary, status, supported validation findings, deterministic ordering, and key artifact paths.
  - [x] Add sprint read-only use cases in `internal/app/sprint_usecases.go` for sprint list, planning artifact status, flow status through `execute`, typed execute status summary, supported validation findings, deterministic ordering, and key artifact paths.
  - [x] Refactor existing CLI adapters only where needed so shared read-only use cases are a real boundary used by CLI and TUI without changing documented CLI text/JSON behavior.
  - [x] Keep product interpretation in `internal/project`, `internal/study`, and `internal/sprint`; do not add global `validation`, `workflow`, `scheduler`, `reports`, `prompts`, or `stages` packages.

- [x] **Task 2: Add Thin TUI Command Startup Through App**
  > Executes: `Decision 1`, `AC1`, `AC2`, `AC3`, `AC9`, `CLI Surface`, `Configuration`, `Errors`, `Security`
  - [x] Add `internal/app/tui_commands.go` command wiring for `ultraplan tui` and `--workspace` using the same workspace discovery and normal non-runtime config validation as CLI commands.
  - [x] Build read-only app use cases and startup options through the existing app composition path.
  - [x] Launch the TUI program without requiring OpenCode, runtime credentials, provider configuration beyond normal non-runtime config validation, runtime health checks, network access, or eager runtime setup.
  - [x] Preserve existing command help, exit codes, text output, and JSON surfaces except for adding the new `tui` command.
  - [x] Return actionable classified startup errors for invalid workspace and config failures.

- [x] **Task 3: Document TUI Package Boundary And Dependency Choice**
  > Executes: `Decision 4`, `Decision 8`, `Documentation`, `Architecture`, `AC12`
  - [x] Add `internal/tui/doc.go` documenting the selected terminal UI dependency, why it fits deterministic model/update/render tests, and that dependency types are contained inside `internal/tui`.
  - [x] Document that `internal/tui` owns terminal program setup, navigation, key handling, model updates, rendering, preview state, and UI-local error panes only.
  - [x] Document read-only scope and non-goals: no validation runs, planning flow, prompt preview, execute, study run-loop, runtime-backed operations, Git mutation, subprocess/plugin launchers, smoke/review/issue artifact generation, or TUI persistence.
  - [x] Include read-only scope in command or key help where user-visible labels could otherwise imply mutation.

- [x] **Task 4: Implement TUI Program Setup And Derived Model**
  > Executes: `Decision 3`, `Decision 4`, `AC4`, `AC9`, `AC10`, `Architecture`, `Performance`, `Persistence And Migrations`
  - [x] Add `internal/tui/app.go` to construct and run the terminal program from typed read-only app use cases and immutable startup options.
  - [x] Add `internal/tui/model.go` with deterministic UI/session state derived from app results, not product-owned state or persisted cache.
  - [x] Implement startup/load messages, explicit refresh messages, selection/navigation state, use-case error state, preview state, and quit behavior.
  - [x] Propagate caller `context.Context` through startup and refresh operations; do not create isolated `context.Background()` inside status or refresh operations.
  - [x] Avoid background periodic refresh, hot config reload, workflow execution, runtime health checks, and separate TUI persistence.

- [x] **Task 5: Implement Read-Only Key Bindings And Action Boundary**
  > Executes: `Decision 7`, `AC3`, `AC10`, `Workflows`, `Security`, `Documentation`
  - [x] Add `internal/tui/keys.go` for navigation, pane switching, selection, explicit refresh, artifact preview, and quit.
  - [x] Ensure no key binding can run validation, sprint flow, prompt preview, execute, study run-loop, runtime work, external commands, Git commands, plugins, or artifact generation.
  - [x] Make refresh labels explicit and non-mutating; if an existing status use case may refresh deterministic status files, label that action before invocation and cover the mutation scope in tests.
  - [x] Keep future operational action placeholders out of Sprint 24 unless they are inert documentation with no executable behavior.

- [x] **Task 6: Implement Rendering And Narrow-Terminal Fallbacks**
  > Executes: `Decision 4`, `AC4-AC7`, `AC11`, `Errors`, `Observability`, `Testing`
  - [x] Add `internal/tui/views.go` rendering for workspace/startup state, project/study/sprint panes, status panes, validation findings, key artifact paths, preview panes, actionable error panes, and key help.
  - [x] Render project status including required docs, project-index health, catalog findings, and key artifact paths.
  - [x] Render study status including source/dimension summary, study detail, supported validation findings, and key artifact paths.
  - [x] Render sprint status including planning artifact status, flow-state status through `execute`, execute run-state summary where available, supported validation findings, and key artifact paths.
  - [x] Add narrow-terminal fallback rendering where critical status and error text remains visible in plain text and meaning does not depend only on ANSI styling.

- [x] **Task 7: Implement Path-Safe Bounded Artifact Previews**
  > Executes: `Decision 5`, `AC8`, `AC10`, `Security`, `Performance`, `Errors`
  - [x] Support previews for known Markdown and JSON artifacts: `project-index.md`, `roadmap.md`, project docs, `requirements.md`, `sprint-index.md`, `technical-handbook.md`, `reasoning.md`, `plan.md`, `execute.md`, `flow-state.json`, and `.run-state.json`.
  - [x] Resolve preview paths through existing workspace-safe path resolution and selected product-root containment; reject escapes and unsupported arbitrary filesystem paths.
  - [x] Bound preview reads by byte and/or line cap, display truncation clearly, and avoid unbounded memory use or recursive repository scans.
  - [x] Represent missing files, path rejection, read errors, truncation, invalid JSON, and Markdown display states as explicit read-only model/view states.
  - [x] Do not write, repair, normalize, or mutate previewed artifacts from TUI preview actions.

- [x] **Task 8: Add Typed Execute Status And Redacted Diagnostic Display**
  > Executes: `Decision 6`, `AC7`, `AC8`, `AC9`, `AC10`, `Observability`, `Security`, `Persistence And Migrations`
  - [x] Ensure sprint read-only use cases expose execute status summaries from product-owned interpretation of durable `flow-state.json`, `execute.md`, and `.run-state.json`.
  - [x] Keep `.run-state.json` product semantics outside `internal/tui`; TUI may render typed display-safe fields and may preview the file as a bounded artifact, but must not parse it as status logic.
  - [x] Redact or omit raw model, provider, runtime diagnostic, environment, native payload, config, evidence, and execution-summary values before rendering.
  - [x] Prefer product/app redaction helpers and display-safe fields; if redaction logic must be touched, keep it outside terminal rendering.
  - [x] Show unavailable/actionable execute status instead of unsafe raw details when typed status cannot safely expose a field.

- [x] **Task 9: Add Deterministic App And TUI Tests**
  > Executes: `Decision 8`, `AC1-AC13`, `Testing`, `Security`, `Performance`, `Workflows`
  - [x] Add `internal/app/usecases_test.go` covering read-only use-case extraction, deterministic ordering, display-safe results, and CLI/TUI reuse.
  - [x] Add `internal/app/tui_commands_test.go` covering help visibility, `ultraplan tui`, `--workspace`, invalid workspace/config startup failures, no-runtime startup, and no CLI stdout parsing.
  - [x] Add `internal/tui/test_fakes_test.go` with fake read-only app use cases, deterministic fixtures, context recording, and mutation counters.
  - [x] Add `internal/tui/model_test.go` covering navigation, selection, pane switching, explicit refresh, context propagation, error handling, validation panes, artifact preview state, execute status state, and read-only guarantees.
  - [x] Add `internal/tui/views_test.go` covering core dashboard rendering, status/finding panes, preview/error panes, unsafe diagnostic redaction, key help, and narrow-terminal behavior.
  - [x] Add preview tests for success, missing file, path escape rejection, truncation, invalid JSON/Markdown display state, read errors, and no mutation.
  - [x] Add key/action tests proving no runtime-backed workflow, validation, execute, study run-loop, Git mutation, subprocess/plugin launcher, smoke/review/issue artifact generation, or hidden mutation action is reachable.

- [x] **Task 10: Run Verification And Prepare Review Evidence**
  > Executes: `Decision 8`, `AC13`, `Architecture`, `Testing`, required review protocols
  - [x] Run `go test ./...` from `../ultraplan-go` and capture pass/fail evidence.
  - [x] Run `go build ./cmd/ultraplan` from `../ultraplan-go` and capture pass/fail evidence.
  - [x] Review import boundaries: product/platform packages do not import `internal/tui`, terminal-library types do not leak into `internal/app` or product result structs, and `internal/tui` depends on app use cases/simple results only.
  - [x] Review read-only boundary: no self-subprocess, no CLI handler/stdout scraping, no runtime execution, no Git mutation, no smoke/review/issue artifact generation, no TUI persistence, and no global workflow/refactor packages.
  - [x] Record any deviation from `reasoning.md` before implementation continues.

## Evidence Checklist

- [x] Tests prove the required behavior.
- [x] Runtime or diagnostic evidence exists where required.
- [x] Documentation updates are complete where required.
- [x] Deviations from `reasoning.md` are recorded before implementation continues.
- [x] Required review protocols have evidence.
- [x] `internal/app` shared read-only use cases are used by CLI and TUI adapters where applicable.
- [x] `internal/tui` terminal dependency remains contained.
- [x] Read-only TUI guarantees are proven by tests and source review.
- [x] Execute status and diagnostics are typed, redacted, and display-safe.
- [x] Artifact previews are path-safe, bounded, and non-mutating.

## Verification Commands

| Check | Command | Expected Result |
| --- | --- | --- |
| Full Go test suite | `go test ./...` | All normal tests pass without requiring OpenCode, provider credentials, network access, or a real terminal. |
| CLI binary build | `go build ./cmd/ultraplan` | `ultraplan` binary builds successfully with the TUI command and contained terminal dependency. |
| App package tests | `go test ./internal/app` | TUI command startup, `--workspace`, help visibility, no-runtime startup, startup errors, and shared use-case behavior pass. |
| TUI package tests | `go test ./internal/tui` | Model, view, key-binding, fake-use-case, preview, redaction, narrow-terminal, and read-only guarantee tests pass non-interactively. |
| Architecture review | `review using system/protocols/architecture-review-protocol.md` | Import boundaries, terminal-library containment, shared use cases, no global workflow packages, and product ownership are confirmed. |
| Sprint review | `review using system/protocols/sprint-review-protocol.md` | Implementation evidence satisfies requirements, acceptance criteria, selected contracts, verification commands, and `reasoning.md`. |

## Risks And Blockers

| Risk / Blocker | Source | Mitigation | Status |
| --- | --- | --- | --- |
| Existing project/study/sprint services may not expose needed read-only status without focused extraction. | `reasoning.md#assumptions-and-risks` | Added narrow `internal/app` use cases over existing product services; no global validation/workflow packages introduced. | `closed` |
| Sprint 23 execute status/redaction gaps may block safe execute display. | `reasoning.md#decision-6-display-execute-and-runtime-diagnostics-only-through-typed-redacted-status-results` | TUI renders typed execute counts and unavailable state only; no raw runtime diagnostics are parsed in `internal/tui`. | `mitigated` |
| App result types may become a parallel product model. | `reasoning.md#potential-technical-debt` | Results are dashboard-specific summaries and artifacts, backed by product services. | `mitigated` |
| `internal/tui` may accidentally interpret product artifacts or own state machines. | `reasoning/architecture.md#risks` | `internal/tui` consumes app result types and preview states only; sprint execute semantics stay in app/sprint use cases. | `closed` |
| Read-only boundary may be weakened by existing status use cases that write deterministic status files. | `reasoning.md#assumptions-and-risks` | Sprint status refresh may recompute `flow-state.json`; command/key help and sprint result labels state this explicitly. | `mitigated` |
| Raw diagnostics or config/model/provider values may leak in TUI views. | `reasoning.md#decision-6-display-execute-and-runtime-diagnostics-only-through-typed-redacted-status-results` | Display strings are redacted in app use cases and tests assert redacted rendering. | `mitigated` |
| Artifact preview path escapes or unbounded reads may create security/performance issues. | `reasoning.md#decision-5-provide-path-safe-bounded-read-only-artifact-previews` | Preview uses workspace containment, supported artifact allowlist, 32 KiB cap, truncation, missing, invalid JSON, and escape tests. | `closed` |
| Rendering tests may become brittle. | `reasoning.md#potential-technical-debt` | Tests assert critical text and narrow behavior rather than full-screen goldens. | `mitigated` |
| Explicit refresh may feel stale compared with live dashboard expectations. | `reasoning.md#decision-3-keep-tui-state-derived-context-aware-and-disposable` | Documented in `internal/tui/doc.go` and key help; live monitoring deferred. | `carried forward` |
| Exact terminal library, narrow-width minimum, and preview truncation cap remain implementation questions. | `reasoning/architecture.md#risks` | Chose Bubble Tea for the interactive event loop, narrow fallback below 72 columns, and 32 KiB preview cap; terminal-library types remain contained inside `internal/tui`. | `closed` |

## Review Inputs

Review should use:

- `sprint-index.md`
- `technical-handbook.md`
- `reasoning/architecture.md`
- `reasoning.md`
- this `plan.md`
- implementation diff
- verification evidence
- `system/protocols/architecture-review-protocol.md`
- `system/protocols/sprint-review-protocol.md`

## Execution Log

| Date / Step | Action | Evidence / Notes |
| --- | --- | --- |
| `2026-07-09 / planning` | Created sprint implementation plan from required inputs and existing reasoning. | Plan writes only `projects/ultraplan-go/sprints/24-tui-foundations/plan.md`; implementation is not started. |
| `2026-07-09 / implementation` | Added shared read-only app use cases, TUI command wiring, contained `internal/tui` model/render/key package, bounded artifact previews, and deterministic tests. | Changed `/home/antonioborgerees/coding/ultraplan-go/cmd/ultraplan/main.go`, `internal/app/app.go`, `internal/app/usecases.go`, `internal/app/project_usecases.go`, `internal/app/study_usecases.go`, `internal/app/sprint_usecases.go`, `internal/app/tui_commands.go`, `internal/tui/*`, and related tests. |
| `2026-07-09 / verification` | Ran targeted and full verification. | `go test ./internal/app ./internal/tui` passed; `go test ./...` passed; `go build ./cmd/ultraplan` passed. |
| `2026-07-09 / command smoke` | Ran `./ultraplan --workspace /home/antonioborgerees/coding/ultraplan-go-workspace tui`. | Command rendered the read-only dashboard header, project pane, project docs, studies, and sprint status without requiring runtime credentials or network access. A malformed legacy sprint flow-state was converted to an in-dashboard finding instead of aborting startup. |
| `2026-07-09 / review` | Reviewed architecture and read-only boundaries. | `internal/app` does not import `internal/tui`; `cmd/ultraplan` registers a runner seam; `internal/tui` contains model/render/key behavior only; no terminal-library type leakage; no `os/exec`, Git mutation, runtime launch, prompt preview, flow, execute, study run-loop, smoke/review/issue artifact generation, or TUI persistence was introduced. |
| `2026-07-09 / interactive TUI correction` | Added Bubble Tea as the contained external terminal dependency and wired `internal/tui.Run` to a full-screen event loop. | `go.mod`/`go.sum` include `github.com/charmbracelet/bubbletea`; `internal/tui/app.go` adapts existing deterministic model/render/key behavior to Bubble Tea; PTY smoke verified alt-screen startup and `q` quit behavior. |
| `2026-07-09 / UI refinement` | Reworked navigation into top-level Projects and Studies tabs, route-stack drilldown, focus movement, selection-follow viewport behavior, and preview scrolling. | Sprints are now project-scoped, artifact navigation is name-only, `tab` moves focus between tabs and content, `left/right` switch tabs while tab focus is active, `enter` opens routes/previews, and `esc` backs out or closes preview. |
| `2026-07-09 / Markdown preview` | Added contained terminal Markdown rendering and expanded supported preview paths. | `internal/tui/markdown.go` uses Glamour inside `internal/tui`; headings render without raw `#` prefixes; study dimensions, Markdown sources, and study run-state are previewable; `go test ./...` and `go build ./cmd/ultraplan` passed. |

## Completion Criteria

- [x] All tasks are complete or explicitly deferred.
- [x] Verification commands were run or deferrals are documented.
- [x] Evidence satisfies the expectations from `reasoning.md`.
- [x] `review.md` can evaluate conformance without guessing intent.
- [x] Required output files from `requirements.md` exist or deferrals are explicitly justified.
- [x] `go test ./...` and `go build ./cmd/ultraplan` pass from `../ultraplan-go` or failures are documented with actionable blockers.
- [x] Architecture and sprint review protocols have enough evidence to evaluate import boundaries, read-only scope, redaction, previews, tests, and acceptance criteria.
