# Sprint Reasoning: TUI Foundation and Read-Only Dashboard

> Project: `ultraplan-go`
> Sprint: `24-tui-foundations`
> Output: `projects/ultraplan-go/sprints/24-tui-foundations/reasoning.md`
> **Inputs Used:** `projects/ultraplan-go/project-index.md`, `projects/ultraplan-go/sprints/24-tui-foundations/requirements.md`, `projects/ultraplan-go/docs/ARCHITECTURE.md`, `projects/ultraplan-go/docs/PRD.md`, `projects/ultraplan-go/docs/TRD.md`, `projects/ultraplan-go/sprints/24-tui-foundations/sprint-index.md`, `projects/ultraplan-go/sprints/24-tui-foundations/technical-handbook.md`, `projects/ultraplan-go/sprints/24-tui-foundations/reasoning/architecture.md`, `studies/go-cli-study/reports/final/01-project-structure.md`, `studies/go-cli-study/reports/final/02-command-architecture.md`, `studies/go-cli-study/reports/final/03-dependency-injection.md`, `studies/go-cli-study/reports/final/04-configuration-management.md`, `studies/go-cli-study/reports/final/05-error-handling.md`, `studies/go-cli-study/reports/final/06-io-abstraction.md`, `studies/go-cli-study/reports/final/07-state-context.md`, `studies/go-cli-study/reports/final/09-terminal-ux.md`, `studies/go-cli-study/reports/final/10-logging-observability.md`, `studies/go-cli-study/reports/final/11-testing-strategy.md`, `studies/go-cli-study/reports/final/12-extensibility.md`, `studies/go-cli-study/reports/final/13-security.md`, `studies/go-cli-study/reports/final/14-performance.md`, `studies/go-cli-study/reports/final/15-philosophy.md`

This document decides. It synthesizes selected context, handbook evidence, area-specific reasoning, and contracts into final sprint decisions.

It does not replace `sprint-index.md`, `technical-handbook.md`, or `reasoning/*.md`.

## Sprint Purpose

- **Goal:** Add a read-only `ultraplan tui` dashboard over shared application use cases so users can inspect workspace, project, study, and sprint state without invoking runtimes, mutating artifacts, shelling out to UltraPlan, or parsing CLI output.
- **Non-Goals:** Running validation, planning flow, prompt preview, execute, study run-loop, guarded action dialogs, live runtime progress panes, cancellation controls, browser UI, hosted service, API server, multi-user collaboration, project-management features, smoke artifacts, automated review artifacts, issue artifacts, Git mutation, CLI replacement, separate TUI persistence, workflow scheduler, broad Sprint 23 hardening, or global `validation`, `workflow`, `scheduler`, `reports`, `prompts`, or `stages` package refactors.
- **Depends On:** Sprint 16 project domain and index, Sprint 17 sprint artifact domain and flow state, Sprints 18-21 planning stages through `plan`, Sprint 22 planning documentation and release, Sprint 23 execute stage, Sprint 23 follow-ups for typed execute status and redacted execute diagnostics, and existing workspace/config/app composition.

## Selected Context And Pre-Reasoning Artifacts

| Artifact | Path | How It Was Used |
| --- | --- | --- |
| Project Index | `projects/ultraplan-go/project-index.md` | Established project scope, repository location, selected contract pool, available study evidence, module ownership rules, and review protocols. |
| Requirements | `projects/ultraplan-go/sprints/24-tui-foundations/requirements.md` | Supplied the authoritative sprint contract, required outputs, acceptance criteria, non-goals, constraints, dependencies, and review expectations. |
| Project Docs | `projects/ultraplan-go/docs/ARCHITECTURE.md`, `projects/ultraplan-go/docs/PRD.md`, `projects/ultraplan-go/docs/TRD.md` | Grounded the TUI as one local surface over one product core, confirmed module-driven architecture, app/TUI boundary, config/workspace rules, and read-only TUI baseline requirements. |
| Sprint Index | `projects/ultraplan-go/sprints/24-tui-foundations/sprint-index.md` | Selected Architecture, CLI Surface, Configuration, Documentation, Errors, Observability, Performance, Persistence And Migrations, Security, Testing, and Workflows contracts; selected evidence reports; required Architecture area reasoning and review protocols. |
| Technical Handbook | `projects/ultraplan-go/sprints/24-tui-foundations/technical-handbook.md` | Distilled selected study evidence into patterns, trade-offs, anti-patterns, examples, design pressures, and open questions for sprint decisions. |
| Prior Decision | `projects/ultraplan-go/project-index.md` | The project index states no prior decision artifacts exist; carry-forward constraints come from requirements dependencies and source docs. |
| Architecture Reasoning | `projects/ultraplan-go/sprints/24-tui-foundations/reasoning/architecture.md` | Decided architecture boundaries for `internal/app`, `internal/tui`, dependency direction, terminal-library containment, read-only handling, explicit refresh, previews, and execute redaction/status ownership. |

## Area-Specific Reasoning Inputs

| Area | Reasoning Document | Key Conclusion | Evidence Basis | Impact On Final Decision |
| --- | --- | --- | --- | --- |
| Architecture | `projects/ultraplan-go/sprints/24-tui-foundations/reasoning/architecture.md` | Implement the TUI as a local interface over typed `internal/app` read-only use cases. `internal/tui` owns terminal setup, navigation, key handling, model updates, rendering, preview state, and UI-local error panes only. Product modules keep validation, state machines, durable artifact interpretation, redaction, and persistence. | Project architecture and TRD app/TUI dependency direction; requirements read-only constraints; handbook reports `01`, `02`, `03`, `06`, `07`, `09`, `13`, `14`, and `15`; concrete sources including gh-cli factories, helm cmd/action split, Bubble Tea examples, restic/gh-cli IO seams, opencode permission/redaction patterns, and yq/age streaming evidence. | Final decisions require shared app use cases before TUI rendering, terminal-library containment in `internal/tui`, explicit refresh, bounded path-safe previews, typed execute status/redaction, no CLI handler calls, no self-subprocess, no runtime workflow execution, and no global workflow/refactor packages. |

## Sprint Technical Handbook Summary

- **Relevant Patterns:** Thin CLI/local interface over shared use cases; composition root plus explicit dependency seams; injectable IO and terminal abstractions; TUI message/model separation; context propagation for startup and refresh; classified user-facing error panes; bounded preview/streaming discipline; behavior-focused tests with fakes; explicit read-only trust boundary.
- **Important Trade-Offs:** Bubble Tea-style model/update vs widget-heavy tcell/tview; shared typed app use cases vs direct product service calls from TUI; constructor injection vs context service lookup; immutable startup config vs hot reload; golden/rendering tests vs targeted substring assertions; bounded preview reads vs complete artifact display; read-only dashboard vs early guarded actions.
- **Warnings / Anti-Patterns:** Do not let TUI call CLI handlers or scrape CLI output; do not let `internal/tui` own product state machines or persistence; avoid global mutable app/config state; avoid direct `os.Stdout`, `os.Stderr`, or real terminal requirements in normal tests; avoid `context.Background()` inside refresh/status operations; do not expose raw secrets, model/provider values, or unredacted diagnostics; avoid unbounded artifact reads and eager runtime setup; do not implement plugins/subprocess extensions for this dashboard.
- **Evidence Confidence:** High for architecture, command wiring, DI, IO/testability, context, terminal UX, testing, security, and performance because the selected reports studied 16 Go CLI repositories and include concrete source references. Medium for observability, extensibility, and philosophy where evidence is broad and useful but less directly tied to a read-only dashboard implementation.

## Contracts Applied

| Contract / Requirement ID | Constraint | Decision Impact | Expected Evidence |
| --- | --- | --- | --- |
| Architecture | Module ownership and dependency direction must stay clear; `internal/app` is local composition/use-case boundary and `internal/tui` owns terminal UI only. | Decisions 1, 2, 3, 4, 5, and 7 keep product state in product modules and terminal types inside `internal/tui`. | Architecture review import checks; source review showing no product/platform package imports `internal/tui`; no terminal-library types in app/product result structs. |
| CLI Surface | `ultraplan tui` and `--workspace` must be discoverable, startup errors actionable, and existing CLI text/JSON surfaces preserved. | Decision 1 adds TUI command wiring through app startup without changing existing command behavior. | `internal/app/tui_commands_test.go`; command help test; startup failure tests; existing command tests still pass. |
| Configuration | TUI startup must reuse workspace discovery and normal config validation while remaining runtime-free and redacting displayed values. | Decisions 1, 3, and 6 use immutable startup config and display-safe result fields. | Tests for `--workspace`, invalid workspace/config, no runtime credential/OpenCode requirement, and redaction checks. |
| Documentation | TUI dependency, ownership boundary, read-only scope, and non-goals must be documented. | Decision 8 requires `internal/tui/doc.go` and command/key-binding help to state scope and non-goals. | Review of `internal/tui/doc.go`, command help, key help, and absence of misleading action labels. |
| Errors | Startup, preview, invalid workspace/config, and use-case failures must preserve classification and actionable diagnostics. | Decisions 1, 5, and 6 route classified errors into error panes without raw unsafe detail. | Error-state model/view tests; command startup error tests; review of wrapping/classification. |
| Observability | Dashboard displays project, study, sprint, execute status, validation findings, diagnostics, and safe run-state summaries without launching workflows. | Decisions 2, 6, and 7 require typed status summaries and redacted diagnostics. | Use-case tests for project/study/sprint status results; TUI rendering tests for status/finding panes; redaction review. |
| Performance | Startup, refresh, deterministic ordering, narrow rendering, and artifact previews must avoid unbounded reads or expensive scans. | Decisions 3, 5, and 6 require explicit refresh, deterministic results, bounded previews, and no eager runtime setup. | Tests for bounded preview truncation and deterministic ordering; review for no recursive repository scan in TUI startup/preview. |
| Persistence And Migrations | Durable workspace artifacts remain source of truth; TUI must not create a separate persistence model. | Decisions 2, 5, and 6 keep TUI state derived and disposable, with product/use-case interpretation of durable artifacts. | Review showing no TUI cache/database/state files; tests verifying read-only actions do not write artifacts. |
| Security | Path-safe previews, no self-subprocess, no runtime credential requirement, no network access, and redaction are mandatory. | Decisions 1, 5, 6, and 7 reject subprocess integration, runtime actions, unsafe paths, and raw diagnostics. | Search/review for absence of `os/exec` self-invocation in TUI path; path escape tests; redaction tests; no-runtime startup tests. |
| Testing | Normal tests must be deterministic and non-interactive with fake use cases and no real terminal/runtime/network. | Decisions 2, 3, 4, 5, 6, and 8 require fake use-case tests, model tests, rendering tests, command tests, and app use-case tests. | `go test ./...`; `internal/tui/*_test.go`; `internal/app/*usecases*_test.go`; fake-use-case fixture review. |
| Workflows | TUI may display workflow state but must not execute planning, execute, study run-loop, validation, or runtime workflows. | Decisions 2, 3, 6, and 7 make workflow state read-only and product-owned. | Model/key tests proving no workflow execution action exists; review of key bindings and use-case interfaces. |
| R1 / Required Outputs | Required files in `internal/app` and `internal/tui` must be created or updated. | All final decisions define architecture and evidence for those outputs without implementing them here. | Review required output paths listed in requirements lines 13-30. |
| AC1-AC2 | `ultraplan tui` and `--workspace` startup behavior. | Decision 1. | Command help and startup tests. |
| AC3 | TUI uses typed app/product results, not subprocess, CLI handlers, stdout/stderr scraping. | Decisions 1 and 2. | Source review and tests using fake use cases, not CLI output fixtures. |
| AC4-AC8 | Dashboard panes, project/study/sprint views, artifact previews. | Decisions 2, 4, 5, and 6. | Model/view/use-case tests for panes, statuses, findings, previews, and deterministic ordering. |
| AC9 | Runtime-free startup. | Decisions 1 and 3. | No-runtime startup command tests and dependency review. |
| AC10 | Read-only actions must not intentionally write or mutate artifacts, with explicit labels if existing status use cases refresh deterministic status files. | Decisions 2, 3, and 7. | Read-only guarantee tests, fake mutation counters, file snapshot review where appropriate. |
| AC11 | Narrow-terminal rendering keeps critical status and errors visible without ANSI-only meaning. | Decision 4. | `views_test.go` narrow-width cases and plain text assertions. |
| AC12-AC13 | Unit tests and verification pass from `../ultraplan-go`. | Decision 8. | `go test ./...` and `go build ./cmd/ultraplan`. |

## Repos Studied / Source Evidence Used

| Source / Repo / Report | Concrete Reference | Relevant Finding | Why It Matters For This Sprint | Used In Decision(s) |
| --- | --- | --- | --- | --- |
| `01-project-structure` | `studies/go-cli-study/reports/final/01-project-structure.md`, Pattern 1, Pattern 4, Pattern 6; `cmd/gh/main.go:6`, `cmd/root.go:9`, `cmd/gdu/app/app.go:30-49` | Mature Go CLIs keep entrypoints thin, preserve unidirectional imports, and use UI abstractions when multiple surfaces exist. | Supports `cmd/ultraplan` and TUI as local interfaces over app/product logic, not product owners. | Decisions 1, 2, 7 |
| `02-command-architecture` | `pkg/cmd/install.go:132-145`, `pkg/cmdutil/factory.go:16-43`, `internal/view/command.go:176` | Command wiring should parse/delegate through reusable logic; TUI-hidden command interpreters create CLI/TUI bifurcation risk. | Supports adding `ultraplan tui` as thin command wiring and rejecting CLI handler/text scraping. | Decisions 1, 2 |
| `03-dependency-injection` | `internal/ghcmd/cmd.go:52-132`, `pkg/cmdutil/factory.go:16-43`, `internal/view/xray.go:412`, `internal/dao/ds.go:72` | Manual composition roots, factories, interfaces, and fakes dominate; context service lookup is ergonomic but less type-safe. | Supports explicit app use-case construction and fake TUI use cases rather than global state or service locators. | Decisions 1, 2, 3, 8 |
| `04-configuration-management` | `internal/cmd/config.go:2253-2287`, `opencode/internal/config/config.go:609-641`, `internal/config/k9s.go:317-337` | Config precedence should be explicit and validated after merge; TUI sessions can use immutable startup config unless hot reload is justified. | Supports reusing CLI workspace/config startup and deferring hot reload. | Decisions 1, 3 |
| `05-error-handling` | `internal/errors/fatal.go:10`, `cmd/restic/main.go:205`, `cmd/age/tui.go:37-54`, `internal/model/flash.go:100-103` | Mature CLIs preserve error chains and separate user-facing rendering from operational details. | Supports classified startup/preview/use-case errors and actionable TUI error panes. | Decisions 1, 5, 6 |
| `06-io-abstraction` | `pkg/iostreams/iostreams.go:551-568`, `internal/ui/terminal.go:10-36`, `internal/ui/mock.go:10-53`, `executor.go:541-577` | Injectable streams, terminal interfaces, and mocks enable non-interactive tests. | Supports fake app use cases, string-rendered views, and tests that do not require a real terminal. | Decisions 4, 8 |
| `07-state-context` | `internal/ghcmd/cmd.go:142`, `task.go:89`, `pkg/cmd/install.go:333-347`, `internal/keys.go:10-38` | Context propagation enables cancellation and refresh discipline; context-as-service-locator is risky. | Supports passing context through startup/refresh while rejecting hidden service lookup. | Decisions 1, 3 |
| `09-terminal-ux` | Bubble Tea examples `internal/chezmoibubbles/passwordinputmodel.go:8-53`, opencode overlays `internal/tui/tui.go:717-723`, fzf renderers `src/tui/tcell.go:1-19`, `src/tui/light.go:110-143` | Full-screen TUIs should match interaction density, remain testable/interruptible, and degrade in narrow/non-TTY contexts. | Supports model/update/render separation, narrow fallback, and contained terminal-library dependency. | Decisions 4, 5 |
| `10-logging-observability` | `pkg/iostreams/iostreams.go:52-54`, `internal/logging/logger.go:25-62`, `internal/slogs/keys.go:6-231`, `internal/ui/progress/printer.go:15-34` | Status/diagnostics/logs should be structured and separated from user output; diagnostics displayed in UI must be safe. | Supports typed status results and redacted execute diagnostics, not raw stdout/log scraping. | Decisions 2, 6 |
| `11-testing-strategy` | `internal/cmd/main_test.go:64-174`, `pkg/httpmock/stub.go:35-199`, `internal/backend/mock/backend.go:14-26`, `task_test.go:166-169` | Behavior tests, fakes, compile-time interface checks, and golden/targeted output tests prevent regressions. | Supports command tests, app use-case tests, fake use-case model tests, view tests, and full verification. | Decision 8 |
| `12-extensibility` | `helm/internal/plugin/runtime_subprocess.go:65-79`, `rclone/fs/registry.go:407`, `k9s/internal/view/actions.go:206-265`, `go-task/executor.go:20-24` | Plugin/subprocess extension systems carry security and lifecycle costs and should be explicit. | Supports rejecting plugins, subprocess launchers, and external command extension surfaces from this read-only TUI. | Decisions 2, 7 |
| `13-security` | `internal/options/secret_string.go:15-20`, `pkg/registry/transport.go:37-41`, `internal/permission/permission.go:44-108`, `internal/config/json/validator.go:146`, `pkg/yqlib/security_prefs.go:3-7` | Trust boundaries require safe path handling, secret redaction, permission gates, schema validation, and default-deny dangerous behavior. | Supports path-contained previews, no runtime actions, no self-subprocess, and redacted display fields. | Decisions 5, 6, 7 |
| `14-performance` | `pkg/cmdutil/factory.go:27-42`, `age/internal/stream/stream.go:20`, `yq/pkg/yqlib/stream_evaluator.go:78-113`, `gdu/pkg/analyze/parallel.go:13` | Startup should defer expensive work; large reads and concurrency must be bounded. | Supports runtime-free lazy startup, bounded previews, explicit refresh, and no recursive scans. | Decisions 3, 5, 6 |
| `15-philosophy` | `age.go:18`, `VISION.md:97`, `pkg/cmd/factory/default.go:26-46`, `fs/types.go:16-59` | Good CLIs accept complexity deliberately and reject out-of-scope behavior. | Supports strict read-only scope and rejecting early guarded actions, scheduler abstractions, and plugin surfaces. | Decisions 2, 7, 8 |
| PRD | `projects/ultraplan-go/docs/PRD.md` lines 46-56, 113-118, 131-144, 149-157, 195 | CLI and TUI share product services; TUI is local terminal, not browser/server; CLI remains stable; no CLI self-invocation. | Product-level source for the read-only local TUI baseline and future operational separation. | Decisions 1, 2, 7 |
| TRD | `projects/ultraplan-go/docs/TRD.md` sections 4.1 and 7.4 | `internal/app` is the local-interface boundary; TUI must use shared app use cases and typed results; normal TUI tests must be non-interactive and narrow rendering must degrade gracefully. | Technical source for app/TUI direction and baseline TUI requirements. | Decisions 1, 2, 4, 8 |
| Architecture doc | `projects/ultraplan-go/docs/ARCHITECTURE.md` lines 21-38, 247-285, 378-409 | Module-driven ownership; `internal/app` owns composition/use cases; `internal/tui` owns terminal widgets/navigation/rendering only. | Governs final architecture and import boundaries. | Decisions 1, 2, 7 |

## Trade-Off And Debt Analysis

Analyze the consequences of the chosen direction before listing final decisions.

### Accepted Trade-Offs

| Trade-Off | Benefit | Cost / Constraint Accepted | Why Acceptable Now | Revisit Trigger |
| --- | --- | --- | --- | --- |
| Shared typed app use cases instead of direct TUI calls into product services | One local-interface boundary for CLI and TUI; avoids CLI text scraping and duplicate workflow logic. | Adds request/result types and a maintenance layer. | Requirements explicitly require shared use cases and two real consumers exist. | Result types become broad mirrors of product internals or CLI-only behavior cannot reuse them. |
| Bubble Tea-style model/update/rendering over widget-heavy tcell/tview design | Deterministic message/model tests, contained terminal dependency, easier narrow rendering checks. | Some dashboard layouts require manual rendering rather than prebuilt widgets. | Sprint 24 prioritizes correctness, read-only navigation, and tests over rich operational widgets. | Sprint 25+ requires complex operational controls that are materially easier with another library. |
| Explicit constructor/factory injection over context service lookup | Compile-time safety, traceable dependencies, simple fakes. | More setup plumbing through TUI app construction. | Testability and boundary review are core acceptance criteria. | Constructor plumbing becomes demonstrably noisy after multiple independent TUI subcomponents. |
| Immutable startup workspace/config snapshot over hot reload | Predictable startup matching CLI precedence; avoids accidental runtime credential checks. | Users must restart TUI for config edits. | Read-only foundation does not need session config mutation. | Future TUI operational controls require long-lived sessions to adapt to config changes safely. |
| Explicit user-triggered refresh over periodic background refresh | Deterministic tests, simpler context lifecycle, no hidden filesystem churn. | Less live-dashboard feel. | Live runtime progress panes and operational monitoring are non-goals. | Requirements add active monitoring or watch-mode behavior. |
| Bounded previews over complete artifact display | Protects memory/latency and supports large artifacts. | Users may need external editor/CLI for full content. | Preview is inspection aid, not editor or pager replacement. | Users need in-TUI full-file paging with measured safe streaming design. |
| Read-only dashboard over guarded workflow actions | Lowest security/mutation risk; clean architecture handoff. | TUI is less operationally powerful initially. | Sprint 24 acceptance criteria explicitly exclude mutating/runtime-backed workflows. | Sprint 25 scopes guarded validation, dry-run, flow, execute, or study operation. |
| Typed execute status in product/app use cases over direct TUI `.run-state.json` parsing | Product modules retain durable artifact semantics and redaction rules. | Requires sprint status/use-case work before rendering execute details. | Requirements require typed execute status and safe redaction. | Product status cannot expose needed display fields without excessive coupling. |

### Potential Technical Debt

| Debt / Shortcut | Why It Might Accrue | Current Mitigation | Owner / Follow-Up |
| --- | --- | --- | --- |
| App result types become a parallel product model | TUI display needs may tempt broad structs that duplicate product state. | Keep results use-case-oriented and display-safe; product modules interpret artifacts. | Sprint 24 implementation review; future app boundary decision if repeated duplication appears. |
| TUI rendering helpers grow into product logic | Rendering status/finding/preview text may tempt local parsing of artifacts. | Require status/finding/preview metadata from app/product use cases; tests assert fake typed inputs. | Architecture review and Sprint Review. |
| Redaction logic is scattered | Execute diagnostics and config fields may be redacted in multiple places if not centralized. | Display-safe fields should be produced by product/app use cases or platform redaction helpers before `internal/tui`. | Sprint 23 follow-up carry-forward and Sprint 24 review. |
| Preview bounds become arbitrary | Initial byte/line caps may not fit all artifacts. | Make truncation explicit and tests assert bounded behavior, not exact universal cap as a product promise. | Future UX tuning when real use reveals preview pain. |
| Rendering tests become brittle | Layout changes can break overly exact golden tests. | Test model behavior separately; use targeted assertions for critical status/error/narrow text unless a stable golden is valuable. | TUI test maintenance in Sprint 24. |
| No periodic refresh may feel stale | Read-only state requires manual refresh. | Key binding and help make refresh explicit; durable artifacts remain source of truth. | Sprint 25 operational/monitoring requirements. |

### Future Considerations

| Consideration | Deferred Until | Reason Deferred | What Should Be Preserved Now |
| --- | --- | --- | --- |
| Guarded validation, dry-run, flow, execute, and study operations from TUI | Sprint 25 or later | Sprint 24 explicitly excludes mutating/runtime-backed operations. | Context-aware startup/refresh, read-only use-case boundary, and clear action labels. |
| Live runtime progress panes and cancellation controls | Sprint 25 or later | Requires operational runtime/event integration and cancellation UX. | Context propagation and typed status/event seams. |
| Config hot reload or session overrides | Future long-lived TUI session sprint | Not needed for read-only foundation and can complicate credential/redaction behavior. | Immutable startup config semantics and explicit restart expectation. |
| Full-file pager/editor behavior | Future UX refinement | Sprint only requires bounded artifact preview. | Path-safe artifact registry, bounded read abstraction, truncation state. |
| Browser UI/API server | Future product phase only | Explicit non-goal for this project phase and sprint. | Keep product use cases independent of terminal-library types. |
| Separate persistence/cache | Future requirement only | Current durable workspace artifacts remain source of truth. | TUI model stays derived and disposable. |
| Plugin/extension surfaces | Future explicit extensibility sprint | Security/lifecycle cost is not justified for read-only dashboard. | No subprocess/plugin registries introduced by accident. |

## Decisions

The final sprint decisions are recorded in detail below. In summary, Sprint 24 will add a thin `ultraplan tui` startup path through `internal/app`, create shared read-only app use cases for CLI and TUI adapters, keep TUI state derived and explicitly refreshed, contain terminal UI implementation inside `internal/tui`, provide path-safe bounded previews, display execute diagnostics only through typed redacted status results, enforce a strict read-only trust boundary, and require deterministic non-interactive tests plus TUI package documentation.

## Final Decisions

### Decision 1: Add Thin `ultraplan tui` Command Startup Through `internal/app`

- **Decision:** Add `ultraplan tui` command wiring in `internal/app/tui_commands.go` that parses `--workspace`, reuses existing workspace discovery and normal non-runtime config validation, builds read-only app use cases, and launches the TUI program. Startup must not require OpenCode, runtime credentials, provider configuration beyond normal non-runtime config validation, network access, or runtime health checks.
- **Rationale:** The CLI remains the stable automation surface while the TUI becomes another local interface over shared use cases. Startup must behave like existing workspace-aware CLI commands but avoid expensive/runtime-backed dependencies.
- **Study / Source Grounding:** `technical-handbook.md` relevant patterns on thin local interfaces, composition roots, config precedence, and context propagation; `02-command-architecture.md` thin-delegate evidence from helm `pkg/cmd/install.go:132-145` and gh-cli `pkg/cmdutil/factory.go:16-43`; `03-dependency-injection.md` composition-root evidence; `04-configuration-management.md` explicit precedence evidence; `07-state-context.md` `ExecuteContextC`/context evidence; PRD and TRD require `ultraplan tui` over the same services and forbid CLI self-invocation.
- **Trade-Offs Accepted:** Startup adds a TUI command and dependency construction path but must keep existing commands unchanged. Immutable startup config is accepted over hot reload.
- **Technical Debt / Future Impact:** Future operational TUI work can reuse the startup context and app construction, but runtime dependencies must remain opt-in when operational actions enter scope.
- **Alternatives Rejected:** Rejected adding a separate binary because the PRD/TRD define the same `ultraplan` binary. Rejected shelling out to `ultraplan` because it violates requirements and command-architecture evidence. Rejected eager runtime health checks because runtime-free startup is mandatory.
- **Contracts Satisfied:** Architecture, CLI Surface, Configuration, Errors, Security, Testing, Workflows; AC1, AC2, AC3, AC9, AC13; required output `internal/app/tui_commands.go`.
- **Evidence Required:** `internal/app/tui_commands_test.go` covers help visibility, `ultraplan tui`, `--workspace`, invalid workspace/config startup failures, no-runtime startup, and no CLI stdout parsing; `go test ./...`; `go build ./cmd/ultraplan`; review confirms existing CLI behavior and JSON/text surfaces are preserved.

### Decision 2: Create Shared Read-Only App Use Cases For CLI And TUI

- **Decision:** Introduce typed read-only interfaces/results in `internal/app/usecases.go`, `project_usecases.go`, `study_usecases.go`, and `sprint_usecases.go` for project list/status/validation findings, study list/detail/status/validation findings, sprint list/flow/execute status/validation findings, artifact paths, and display-safe summaries. CLI adapters and TUI adapters must call these use cases rather than duplicating product workflows or parsing rendered output.
- **Rationale:** A shared app boundary is required because CLI and TUI are two local surfaces over the same product core. Product modules keep ownership of domain semantics while app use cases expose stable, display-safe results.
- **Study / Source Grounding:** `technical-handbook.md` shared typed app use-case trade-off; `01-project-structure.md` thin CLI and UI abstraction evidence; `02-command-architecture.md` warns against bifurcated CLI/TUI behavior hidden in interactive interpreters; `03-dependency-injection.md` supports explicit interfaces/fakes; `10-logging-observability.md` supports structured status/diagnostic separation; Architecture reasoning lines 9-12 decide this boundary.
- **Trade-Offs Accepted:** Adds a small layer of request/result types. This is acceptable because it prevents TUI/product coupling and provides test seams.
- **Technical Debt / Future Impact:** Result types must not become a complete duplicate domain model; future work should add fields only when a real local interface needs them.
- **Alternatives Rejected:** Rejected direct product service calls from `internal/tui` because that would bypass the shared local-interface boundary. Rejected parsing CLI text/JSON because CLI output is an automation surface, not an internal integration protocol. Rejected global workflow/validation abstractions because architecture requires module-owned product behavior.
- **Contracts Satisfied:** Architecture, Observability, Persistence And Migrations, Security, Testing, Workflows; AC3, AC4, AC5, AC6, AC7, AC10; required outputs `internal/app/usecases.go`, `project_usecases.go`, `study_usecases.go`, `sprint_usecases.go`.
- **Evidence Required:** `internal/app/usecases_test.go` verifies read-only use-case extraction and reuse by CLI/TUI adapters; tests assert deterministic ordering and display-safe results; architecture review confirms `internal/tui` depends on app use cases/simple results and product modules do not import `internal/tui`.

### Decision 3: Keep TUI State Derived, Explicitly Refreshed, Context-Aware, And Disposable

- **Decision:** The TUI model stores only UI/session state derived from read-only app use-case results. Startup receives an immutable workspace/config snapshot. Refresh is explicit user-triggered and calls app use cases with `context.Context`; no background periodic refresh, workflow execution, runtime health checks, or separate TUI persistence is added.
- **Rationale:** Durable workspace artifacts remain the source of truth. Explicit refresh keeps lifecycle and tests deterministic while preserving a context seam for future cancellation/operation support.
- **Study / Source Grounding:** `technical-handbook.md` context propagation, immutable config, explicit read-only boundary, and performance warnings; `04-configuration-management.md` immutable config evidence; `07-state-context.md` root context and cancellation evidence; `14-performance.md` lazy initialization and background-refresh trade-off; Architecture reasoning lines 19 and 35-37.
- **Trade-Offs Accepted:** Less live-dashboard feel and users restart for config changes. This is acceptable because live monitoring and operational controls are non-goals.
- **Technical Debt / Future Impact:** Future watch/live progress features must add lifecycle management deliberately rather than piggybacking on hidden background goroutines.
- **Alternatives Rejected:** Rejected background periodic refresh because it complicates tests and can create hidden work. Rejected config hot reload because read-only baseline does not need it. Rejected a TUI cache/database because it violates persistence constraints.
- **Contracts Satisfied:** Configuration, Performance, Persistence And Migrations, Testing, Workflows; AC4, AC9, AC10; constraints on context, source of truth, and runtime-free startup.
- **Evidence Required:** `internal/tui/model_test.go` covers refresh messages, context-aware fake use cases, no hidden mutation, and deterministic state transitions; review confirms no TUI persistence files and no runtime setup in startup/refresh path.

### Decision 4: Implement A Contained Testable TUI Model, Rendering, And Key-Binding Layer

- **Decision:** Implement `internal/tui/app.go`, `model.go`, `views.go`, and `keys.go` around a deterministic model/update/render structure. Terminal-library types stay inside `internal/tui`; product modules, platform packages, and app result structs must not expose terminal-library types. Rendering must support navigable project/study/sprint panes, status panes, findings, artifact previews, error panes, key help, and narrow-terminal fallbacks where critical status/error text remains visible without relying on ANSI-only meaning.
- **Rationale:** A full-screen terminal dashboard needs an explicit UI state model, but that model must remain local to the UI and testable without a real terminal.
- **Study / Source Grounding:** `technical-handbook.md` Bubble Tea-style model/update evidence; `09-terminal-ux.md` Bubble Tea vs tcell/tview trade-off and narrow/non-TTY warnings; `06-io-abstraction.md` terminal mocks and IO streams; `11-testing-strategy.md` behavior/golden test evidence; Architecture reasoning lines 17 and 31.
- **Trade-Offs Accepted:** A Bubble Tea-style architecture may require manual multi-pane rendering compared with widget-heavy libraries. The benefit is deterministic model and rendering tests.
- **Technical Debt / Future Impact:** If future operational widgets become complex, the contained `internal/tui` boundary permits changing library internals without product API changes.
- **Alternatives Rejected:** Rejected exposing terminal-library types in app/product results because it would leak UI concerns. Rejected real-terminal-only tests because requirements demand deterministic non-interactive tests. Rejected ANSI-only status meaning because narrow/accessibility requirements require text visibility.
- **Contracts Satisfied:** Architecture, Documentation, Errors, Observability, Performance, Testing; AC4, AC5, AC6, AC7, AC8, AC11, AC12; required outputs `internal/tui/app.go`, `model.go`, `views.go`, `keys.go`.
- **Evidence Required:** `internal/tui/model_test.go` covers navigation, pane switching, selection, refresh, errors, validation panes, artifact preview state, and read-only behavior; `internal/tui/views_test.go` covers core dashboard rendering and narrow terminals; import review confirms terminal dependency containment.

### Decision 5: Provide Path-Safe, Bounded, Read-Only Artifact Previews

- **Decision:** Artifact previews support key Markdown and JSON files required by the sprint, including project docs, `project-index.md`, `roadmap.md`, docs, `requirements.md`, `sprint-index.md`, `technical-handbook.md`, `reasoning.md`, `plan.md`, `execute.md`, `flow-state.json`, and `.run-state.json`, using workspace-safe path resolution and bounded reads. Missing files, path rejection, truncation, and read/format errors become explicit read-only UI states.
- **Rationale:** Previews are high-value for local inspection but create security and performance risk. The TUI should preview known artifacts safely, not become a file browser or editor.
- **Study / Source Grounding:** `technical-handbook.md` bounded preview and path-safety design pressure; `13-security.md` path/trust-boundary and redaction evidence; `14-performance.md` streaming/bounded-read evidence from age and yq; `05-error-handling.md` user-facing error separation; Architecture reasoning line 21.
- **Trade-Offs Accepted:** Users may not see full files in the TUI. Truncation indicators and artifact paths make the limitation explicit.
- **Technical Debt / Future Impact:** Preview caps may need tuning; future full-file paging should build on the same path-safe artifact registry and bounded read abstraction.
- **Alternatives Rejected:** Rejected arbitrary filesystem browsing because it expands security scope. Rejected unbounded reads because large artifacts can harm responsiveness. Rejected silent missing-file handling because errors must be actionable.
- **Contracts Satisfied:** Security, Performance, Persistence And Migrations, Errors, Testing; AC8, AC10, AC11; path handling constraints.
- **Evidence Required:** TUI model/view tests for preview success, missing file, path escape rejection, truncation, invalid JSON/Markdown display state, and read-only behavior; source review confirms existing workspace-safe path resolution is used.

### Decision 6: Display Execute And Runtime Diagnostics Only Through Typed Redacted Status Results

- **Decision:** Execute status shown in the TUI must come from `internal/app` sprint use-case results backed by `internal/sprint` status interpretation of durable `flow-state.json`, `execute.md`, and `.run-state.json`. `internal/tui` must not parse `.run-state.json` as product logic. Any model, provider, runtime diagnostic, evidence, config, environment, native payload, or execution summary text displayed by the TUI must be redacted before rendering or represented only as display-safe fields.
- **Rationale:** Sprint 23 carry-forward gaps make typed execute status and redacted diagnostics prerequisite to safe TUI display. Product modules own durable execute semantics; TUI owns only presentation.
- **Study / Source Grounding:** Requirements dependencies for Sprint 23 fixes R-3, R-5, and R-8; `technical-handbook.md` warnings on raw secrets/diagnostics and direct JSON parsing in UI layers; `10-logging-observability.md` structured diagnostics separation; `13-security.md` `SecretString`, credential scrubbing, and permission evidence; Architecture reasoning lines 23 and 75.
- **Trade-Offs Accepted:** TUI execute panes depend on typed sprint status work and may initially show summaries rather than raw detail. Safety is prioritized over completeness.
- **Technical Debt / Future Impact:** If redaction logic is duplicated, future work should consolidate it in platform/app helpers. The TUI must preserve display-safe seams for future live execute monitoring.
- **Alternatives Rejected:** Rejected direct `.run-state.json` parsing in `internal/tui` because it duplicates product status logic. Rejected showing raw runtime diagnostics because it risks secret/model/provider leakage. Rejected launching execute/status refresh workflows from TUI because this sprint is read-only.
- **Contracts Satisfied:** Configuration, Observability, Persistence And Migrations, Security, Workflows, Testing; AC7, AC8, AC9, AC10; Sprint 23 dependency constraints R-3, R-5, R-8.
- **Evidence Required:** Sprint use-case tests for execute status summaries and redaction; TUI view tests proving unsafe sample diagnostics do not render; review confirms TUI uses typed status results and no product JSON parsing.

### Decision 7: Enforce A Strict Read-Only Trust Boundary And Reject Operational/Extension Surfaces

- **Decision:** TUI key bindings and model actions are limited to navigation, pane switching, selection, explicit refresh, artifact preview, and quit. The TUI must not run validation, planning flow, prompt preview, execute, study run-loop, runtime-backed operations, external commands, Git commands, plugins, review/smoke/issue artifact generation, or automatic workspace mutations. If an existing status use case may refresh deterministic status files, the action label must explicitly say so before invocation and tests must cover it.
- **Rationale:** This sprint establishes foundation and inspection only. Mutating operations require guarded action design, permission/impact display, cancellation, and runtime evidence that are intentionally deferred.
- **Study / Source Grounding:** Requirements non-goals and constraints; `12-extensibility.md` plugin/subprocess lifecycle/security costs; `13-security.md` permission/trust-boundary evidence; `15-philosophy.md` explicit non-goals and deliberate complexity control; Architecture reasoning lines 13, 25, and 41-43.
- **Trade-Offs Accepted:** Users cannot operate workflows from the first TUI. This is required to avoid accidental mutation and architecture creep.
- **Technical Debt / Future Impact:** Future operational actions will need guarded dialogs, mutation labels, cancellation, and runtime status integration. The current key/action model should leave room for adding those later without changing product ownership.
- **Alternatives Rejected:** Rejected adding guarded validation/execute early because Sprint 24 non-goals exclude it. Rejected external command/plugin launchers because extensibility/security evidence shows lifecycle costs. Rejected Git mutation because project and sprint requirements explicitly defer automatic Git operations.
- **Contracts Satisfied:** Security, Workflows, Persistence And Migrations, Architecture, Documentation, Testing; AC3, AC9, AC10; non-goals lines 48-58 and constraints lines 60-73.
- **Evidence Required:** Key-binding tests prove no mutating/runtime-backed action is reachable; fake use cases record no mutation calls; source review/search confirms no `os/exec` self-invocation, no Git mutation, no runtime execution, and no review/smoke/issue artifact generation from TUI.

### Decision 8: Require Deterministic Non-Interactive Test And Documentation Evidence

- **Decision:** The implementation must include deterministic tests for app use cases, TUI fake use cases, model updates, rendering, command startup, error states, previews, narrow terminals, no-runtime startup, and read-only guarantees. Documentation in `internal/tui/doc.go` must state the selected terminal UI dependency, ownership boundary, read-only scope, non-goals, and dependency containment. Final verification is `go test ./...` and `go build ./cmd/ultraplan` from `../ultraplan-go`.
- **Rationale:** A TUI can regress silently if tested only manually. The sprint requires non-interactive tests and clear docs so future operational work does not reinterpret the boundary.
- **Study / Source Grounding:** `technical-handbook.md` structured/behavior-focused testing pattern; `06-io-abstraction.md` IO/test constructors; `11-testing-strategy.md` behavior tests, fakes, golden/targeted output evidence; `09-terminal-ux.md` narrow rendering guidance; requirements review expectations and verification gate.
- **Trade-Offs Accepted:** Tests add upfront work and some rendering assertions may need maintenance. The value is preventing boundary, UX, and read-only regressions.
- **Technical Debt / Future Impact:** Avoid brittle full-layout golden tests where targeted assertions suffice. Future operational TUI work should extend fake use cases rather than requiring real runtime tests.
- **Alternatives Rejected:** Rejected manual-only TUI verification because normal tests must be deterministic and non-interactive. Rejected real terminal/runtime/network requirements because they would violate acceptance criteria. Rejected undocumented dependency choice because terminal-library containment must be reviewable.
- **Contracts Satisfied:** Documentation, Testing, Architecture, CLI Surface, Security, Performance; AC1-AC13; required outputs `internal/tui/doc.go`, `test_fakes_test.go`, `model_test.go`, `views_test.go`, `tui_commands_test.go`, `usecases_test.go`.
- **Evidence Required:** `go test ./...`; `go build ./cmd/ultraplan`; specific tests listed above; Architecture Review and Sprint Review protocol checks.

## Expected Evidence

| Evidence Type | Required Evidence | Source / Command / Review Check |
| --- | --- | --- |
| Tests | Full Go test suite passes. | From `../ultraplan-go`: `go test ./...`. |
| Build | CLI binary builds after adding TUI command and dependency. | From `../ultraplan-go`: `go build ./cmd/ultraplan`. |
| Command | `ultraplan tui` appears in help; `--workspace` works; invalid workspace/config failures are actionable; no-runtime startup succeeds. | `internal/app/tui_commands_test.go`; command help review. |
| Use Cases | Project/study/sprint read-only use cases return typed, deterministic, display-safe results reusable by CLI and TUI. | `internal/app/usecases_test.go`; review of `internal/app/*usecases*.go`. |
| TUI Model | Navigation, pane switching, selection, refresh, error states, validation panes, artifact preview state, and read-only guarantees are deterministic. | `internal/tui/model_test.go`; `internal/tui/test_fakes_test.go`. |
| TUI Rendering | Core dashboard panes and narrow-terminal fallback keep critical status and errors visible without ANSI-only meaning. | `internal/tui/views_test.go`. |
| Preview | Markdown/JSON previews are path-safe, bounded, read-only, truncation-aware, and error-aware. | `internal/tui/model_test.go`, `internal/tui/views_test.go`, and path-safety tests where preview use cases live. |
| Redaction | Execute status and diagnostics shown in TUI are display-safe; unsafe model/provider/env/native diagnostic samples do not render. | Sprint use-case tests; TUI rendering tests; security review. |
| Read-Only | No runtime execution, validation, planning flow, execute, study run-loop, Git mutation, self-subprocess, smoke/review/issue artifact generation, or separate TUI persistence is reachable. | Key-binding/model tests; fake mutation counters; source review/search. |
| Architecture | `internal/tui` depends only on `internal/app` and simple result types; terminal-library types do not leak outside `internal/tui`; product/platform packages do not import `internal/tui`; no global workflow/refactor packages introduced. | Architecture Review: `system/protocols/architecture-review-protocol.md`. |
| Documentation | TUI dependency, ownership boundary, read-only scope, non-goals, and future operational boundary documented. | `internal/tui/doc.go`; command/key help review. |
| Review | Completed sprint evidence matches requirements, acceptance criteria, contracts, and verification. | Sprint Review: `system/protocols/sprint-review-protocol.md`. |

## Assumptions And Risks

| Item | Type | Impact | Mitigation / Follow-Up |
| --- | --- | --- | --- |
| Existing project/study/sprint services can expose needed read-only status without broad rewrites. | Assumption | If false, implementation may need focused service extraction before TUI panes can render. | Keep extraction minimal and product-owned; do not create global validation/workflow packages. |
| Sprint 23 execute status/redaction gaps can be addressed enough for safe display. | Assumption | If false, execute panes may be incomplete or unsafe. | Gate execute display through typed redacted use-case results; show unavailable/actionable status rather than raw data. |
| Terminal library choice can remain contained in `internal/tui`. | Assumption | If library types leak, app/product APIs become coupled to UI internals. | Architecture review import/type checks; use simple result structs across app boundary. |
| App result types may grow too broad. | Risk | Duplicated product model and high maintenance cost. | Keep result types use-case-specific, display-safe, and backed by product-owned interpretation. |
| TUI accidentally becomes a workflow executor. | Risk | Violates read-only scope and security/workflow contracts. | Restrict key bindings/use-case interfaces; tests with fake mutation counters; review for runtime/Git/subprocess usage. |
| Raw diagnostics or config/model/provider values leak. | Risk | Security and privacy issue. | Redact before rendering; tests with unsafe fixture values; use display-safe fields. |
| Artifact preview path escapes or unbounded reads. | Risk | Workspace safety and performance issue. | Use existing workspace-safe path resolution; reject escapes; bound reads; test truncation and rejection. |
| Rendering tests become brittle. | Risk | Slows iteration and encourages deleting tests. | Separate model tests from view tests; assert critical text and behavior rather than full screen where appropriate. |
| No periodic refresh disappoints users. | Risk | Dashboard may feel stale. | Document explicit refresh; defer live monitoring to Sprint 25 where lifecycle can be designed. |
| Existing status use cases may write deterministic status files. | Risk | TUI read-only guarantee becomes ambiguous. | Prefer non-writing status APIs; if unavoidable, label action explicitly and test mutation scope. |

## Implementation Constraints

- `internal/tui` owns terminal program setup, navigation, key handling, model updates, rendering, preview state, and UI-local error panes only.
- `internal/app` is the shared local-interface boundary for CLI and TUI; CLI/TUI adapters call typed use cases rather than parsing rendered output.
- Product modules own state interpretation, validators, workflow state machines, prompt construction, runtime execution, durable artifact semantics, and persistence.
- Terminal-library types must not appear in `internal/app`, product modules, platform packages, or app use-case result structs.
- TUI startup and refresh must accept and propagate `context.Context`; do not create isolated `context.Background()` inside status/refresh operations.
- TUI startup must not initialize runtime execution, require OpenCode, require provider credentials, run runtime health checks, or access the network.
- TUI actions must not run validation, planning flow, prompt preview, execute, study run-loop, runtime-backed work, external plugins, Git mutation, smoke, review, or issue artifact generation.
- Durable workspace artifacts remain the source of truth; no separate TUI persistence model, database, or cache is introduced.
- Artifact preview paths must be workspace-contained or selected product-root-contained through existing safe path resolution and must reject escapes.
- Artifact preview reads must be bounded and must display truncation clearly.
- Execute status must come from typed sprint/app use-case results; `internal/tui` must not parse `.run-state.json` directly as product logic.
- Any runtime/config/model/provider/diagnostic/evidence text shown in TUI must be redacted before rendering or represented as display-safe fields.
- Deterministic ordering is required for project, study, sprint, artifact, and finding lists.
- Tests must be deterministic and non-interactive; normal tests must not require OpenCode, provider credentials, network access, or a real terminal.
- Existing CLI help, exit codes, text output, and documented JSON surfaces must not be weakened or changed except for adding the `tui` command.
- Do not introduce global `validation`, `workflow`, `scheduler`, `reports`, `prompts`, or `stages` packages for this sprint.

## Plan Handoff

`plan.md` must execute these decisions. It must not invent architecture, scope, or decisions beyond this document.

The plan must carry forward:

- final decisions
- applicable contracts and requirement IDs
- expected evidence
- risks and assumptions
- required review protocols

## Phase Exit Criteria

- [x] Selected context was read and used.
- [x] Area-specific reasoning documents were completed where applicable or explicitly marked not applicable.
- [x] Area-specific reasoning conclusions are reflected or explicitly overridden in final decisions.
- [x] Contracts are explicitly mapped to decisions and expected evidence.
- [x] Final decisions are clear enough for `plan.md` to execute.
- [x] Expected evidence is specific and reviewable.
