# Sprint Requirements: TUI Foundation and Read-Only Dashboard

> Project: `ultraplan-go`
> Sprint: `24-tui-foundations`
> Purpose: the authoritative, human-readable sprint contract. All other sprint artifacts must satisfy these requirements.

## Sprint Goal

Add a read-only `ultraplan tui` dashboard over shared application use cases so users can inspect workspace, project, study, and sprint state without invoking runtimes, mutating artifacts, shelling out to UltraPlan, or parsing CLI output.

## Required Outputs

| Output | Path | Description |
| ------ | ---- | ----------- |
| Sprint requirements | `projects/ultraplan-go/sprints/24-tui-foundations/requirements.md` | This sprint contract for the TUI foundation. |
| TUI command wiring | `../ultraplan-go/internal/app/tui_commands.go` | Adds `ultraplan tui` command parsing, workspace/config startup, and TUI program launch through app use cases. |
| Shared app use cases | `../ultraplan-go/internal/app/usecases.go` | Defines typed read-only use-case interfaces/results used by both CLI adapters and the TUI. |
| Read-only project use cases | `../ultraplan-go/internal/app/project_usecases.go` | Exposes project list, project status, and project validation findings without CLI rendering. |
| Read-only study use cases | `../ultraplan-go/internal/app/study_usecases.go` | Exposes study list, study detail, study status, and study validation findings without CLI rendering. |
| Read-only sprint use cases | `../ultraplan-go/internal/app/sprint_usecases.go` | Exposes sprint list, sprint flow status, execute status summary where available, and sprint validation findings without CLI rendering. |
| TUI package documentation | `../ultraplan-go/internal/tui/doc.go` | Documents the selected terminal UI dependency, ownership boundary, read-only scope, and non-goals. |
| TUI program setup | `../ultraplan-go/internal/tui/app.go` | Constructs and runs the terminal program from typed app use cases and startup options. |
| TUI model/update | `../ultraplan-go/internal/tui/model.go` | Implements deterministic dashboard state, navigation messages, refresh messages, and read-only update behavior. |
| TUI rendering | `../ultraplan-go/internal/tui/views.go` | Renders workspace selector/startup state, lists, status panes, findings, artifact previews, error panes, and narrow-terminal fallbacks. |
| TUI key bindings | `../ultraplan-go/internal/tui/keys.go` | Defines keyboard navigation, pane switching, refresh, preview, and quit behavior. |
| TUI fixtures/fakes | `../ultraplan-go/internal/tui/test_fakes_test.go` | Provides fake read-only app use cases for deterministic TUI tests. |
| TUI model tests | `../ultraplan-go/internal/tui/model_test.go` | Covers navigation, selection, refresh, error handling, validation panes, artifact preview state, and read-only guarantees. |
| TUI rendering tests | `../ultraplan-go/internal/tui/views_test.go` | Covers core dashboard rendering and narrow-terminal behavior without requiring an interactive terminal. |
| TUI command tests | `../ultraplan-go/internal/app/tui_commands_test.go` | Covers `ultraplan tui`, `--workspace`, startup failures, no-runtime startup, and no CLI stdout parsing. |
| App use-case tests | `../ultraplan-go/internal/app/usecases_test.go` | Covers read-only use-case extraction and verifies results are reusable by CLI and TUI adapters. |

## Acceptance Criteria

- [ ] `ultraplan tui` is available in command help and starts the TUI from inside a workspace using the same workspace discovery and config validation rules as CLI commands.
- [ ] `ultraplan tui --workspace <path>` starts against the explicitly selected workspace and reports actionable startup errors when the workspace or config is invalid.
- [ ] The TUI uses typed app use cases and product service results; it does not invoke `ultraplan` as a subprocess, call CLI command handlers, parse CLI stdout/stderr, or scrape rendered text.
- [ ] The dashboard shows projects, studies, and sprints in navigable panes with deterministic ordering.
- [ ] Project views show project list, project status, required docs, project-index health, catalog findings, and key artifact paths.
- [ ] Study views show study list, study detail, source/dimension summary, study status, validation findings where existing validators support them, and key artifact paths.
- [ ] Sprint views show sprint list, planning artifact status, flow-state status through `execute`, execute run-state summary when `.run-state.json` exists, validation findings where existing validators support them, and key artifact paths.
- [ ] Artifact preview supports key Markdown and JSON files, including `project-index.md`, `roadmap.md`, docs, `requirements.md`, `sprint-index.md`, `technical-handbook.md`, `reasoning.md`, `plan.md`, `execute.md`, `flow-state.json`, and `.run-state.json`, with bounded reads and clear errors for missing files.
- [ ] The dashboard starts and remains usable without runtime credentials, OpenCode availability, provider configuration beyond normal non-runtime config validation, or network access.
- [ ] Read-only TUI actions do not intentionally write or mutate workspace artifacts; if an existing status use case refreshes deterministic status files, the TUI action label explicitly says so before invocation.
- [ ] Terminal rendering keeps critical status and error text visible in narrow terminals and does not rely on ANSI-only meaning.
- [ ] Unit tests cover TUI navigation/model updates, view rendering, startup command behavior, fake app use cases, error states, no-runtime startup, and read-only action guarantees.
- [ ] Verification passes from `../ultraplan-go`: `go test ./...` and `go build ./cmd/ultraplan`.

## Non-Goals

- Running validation, planning flow, prompt preview, execute, study run-loop, or any other mutating/runtime-backed workflow from the TUI.
- Adding guarded action dialogs, live runtime progress panes, cancellation controls, or operational workflow controls; those belong to Sprint 25.
- Creating a browser UI, hosted service, API server, multi-user collaboration surface, or project-management features.
- Generating smoke artifacts, automated review artifacts, issue artifacts, or any Sprint Phase 2 deferred artifacts from the TUI.
- Performing Git add, commit, push, checkout, reset, branch, merge, or other automatic Git mutation.
- Replacing, weakening, or changing documented CLI command behavior or JSON/text automation surfaces.
- Building a separate TUI persistence model, database, cache of product truth, or workflow scheduler.
- Refactoring study, project, sprint, or runtime product semantics into global `validation`, `workflow`, `scheduler`, `reports`, `prompts`, or `stages` packages.
- Solving all Sprint 23 execute hardening follow-ups unless they are required for the read-only TUI surface to display accurate status safely.

## Constraints

- `internal/tui` owns terminal widgets, navigation, key handling, model updates, and rendering only; it must not own product state machines, validation rules, prompt rendering, runtime execution, artifact persistence, or workspace mutation logic.
- `internal/app` is the shared local-interface boundary; CLI and TUI must call common typed use cases rather than duplicating workflow logic.
- The TUI must not depend on CLI command handlers, subprocess calls to `ultraplan`, or text parsing of CLI output for normal UltraPlan operations.
- CLI remains the stable automation surface; TUI work must not break existing command help, exit codes, text output, or documented JSON surfaces.
- Durable workspace artifacts remain the source of truth; TUI state is derived and refreshable, not authoritative.
- TUI startup and refresh must use `context.Context` so later cancellation support can be added without redesign.
- TUI dependency selection must stay contained behind `internal/tui` and must not leak terminal-library types into product modules or platform runtime packages.
- Normal tests must be deterministic and non-interactive; they must use fake app use cases and must not require OpenCode, provider credentials, network access, or a real terminal.
- Path handling for artifact previews must use existing workspace-safe path resolution and must reject paths outside the workspace or selected product roots.
- Read-only preview reads must be bounded so large artifacts or repositories do not create unbounded memory use.
- Execute status shown in the TUI must come from durable `.run-state.json`, `execute.md`, and sprint status use cases without launching runtime work.
- Carry-forward Sprint 23 safety gaps that affect display must be handled before exposing them in the TUI: execute model/config values and runtime diagnostics displayed by the TUI must be redacted, and execute status must be available through typed status/use-case results rather than direct JSON parsing in `internal/tui`.

## Dependencies

| Prior Sprint / Output | Required For | Notes |
| --------------------- | ------------ | ----- |
| Sprint 16 project domain and index | Project list/status panes | Provides project discovery, project-index parsing, catalog validation, and project status services. |
| Sprint 17 sprint artifact domain and flow state | Sprint list/flow status panes | Provides sprint discovery, artifact paths, `flow-state.json`, and status derivation. |
| Sprints 18-21 planning stages through `plan` | Planning artifact status and validation panes | Provides validators and artifact semantics for `sprint-index.md`, `technical-handbook.md`, `reasoning.md`, and `plan.md`. |
| Sprint 22 planning documentation and release | Baseline CLI/help/docs expectations | Roadmap marks planning release through `plan.md` complete; no `review.md` was available in the sprint directory. |
| Sprint 23 execute stage | Execute status pane and `.run-state.json` preview | Execute is accepted with follow-ups; Sprint 24 may consume durable execute artifacts but must not run execute from the TUI. |
| Sprint 23 required fix R-3 | Accurate execute status in TUI | Execute status/progress must be surfaced by typed sprint status/use-case results before or during this sprint; `internal/tui` must not parse `.run-state.json` directly as a substitute for product status logic. |
| Sprint 23 required fixes R-5 and R-8 | Safe TUI display of config/diagnostics | Any model, variant, diagnostic, evidence, or runtime summary text displayed by the TUI must go through redaction before rendering. |
| Existing workspace/config/app composition | TUI startup and shared use cases | The TUI must reuse workspace discovery, config validation, and dependency construction from `internal/app`. |

## Review Expectations

| What | How Verified |
| ---- | ------------ |
| TUI command exists and starts correctly | Run `go test ./internal/app` and inspect command help tests for `ultraplan tui` and `--workspace`. |
| Shared use-case boundary is real | Review `internal/app/*usecases*.go` and confirm CLI/TUI adapters call typed use cases instead of duplicating product workflows. |
| No CLI text scraping or self-subprocess integration | Search/review `internal/tui` and TUI app wiring for absence of CLI handler calls, `os/exec` self-invocation, and stdout/stderr parsing. |
| Read-only scope is preserved | Review TUI actions and tests; confirm no runtime execution, planning flow, execute, study run-loop, Git mutation, or artifact-generation action is reachable. |
| Dashboard coverage | Run deterministic TUI model/view tests covering project, study, sprint, validation, artifact preview, and error panes. |
| Runtime-free startup | Tests must show dashboard startup with fake use cases and no runtime credentials/OpenCode dependency. |
| Execute artifact visibility is safe | Review sprint status/use-case results and TUI rendering for execute counts, `.run-state.json` pointer, diagnostics, and redaction. |
| Terminal robustness | Run TUI view tests for narrow widths and confirm critical errors/status are still visible. |
| Architecture boundaries | Review imports: `internal/tui` may depend on `internal/app` and simple result types, but product/platform packages must not import `internal/tui`, and terminal-library types must not leak into product modules. |
| Verification gate | Run from `../ultraplan-go`: `go test ./...` and `go build ./cmd/ultraplan`. |
