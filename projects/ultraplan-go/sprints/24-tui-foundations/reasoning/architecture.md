# Architecture Reasoning: TUI Foundation and Read-Only Dashboard

> **Inputs Used:** `projects/ultraplan-go/sprints/24-tui-foundations/sprint-index.md`, `projects/ultraplan-go/sprints/24-tui-foundations/technical-handbook.md`, `projects/ultraplan-go/sprints/24-tui-foundations/requirements.md`, `projects/ultraplan-go/docs/ARCHITECTURE.md`, `projects/ultraplan-go/docs/PRD.md`, `projects/ultraplan-go/docs/TRD.md`, `system/reasoning/architecture_reasoning_template.md`

Scope: this document decides the architecture shape for Sprint 24's read-only `ultraplan tui` foundation. It covers `internal/app` use-case extraction, `internal/tui` ownership, dependency direction, terminal-library containment, read-only state handling, and carry-forward execute status/redaction constraints. It does not decide non-architecture contract details or implement any source changes.

## Area Decisions

The TUI should be implemented as a local interface over shared application use cases, not as a new product owner. The required flow is `cmd/ultraplan -> internal/app -> product modules/platform` for command startup and `internal/tui -> internal/app -> product modules/platform` for dashboard behavior. `internal/tui` owns terminal program setup, navigation, key handling, model updates, rendering, preview state, and UI-local error panes only. Product modules continue to own state interpretation, validators, workflow state machines, durable artifacts, and persistence.

The sprint should introduce a narrow typed read-only use-case boundary in `internal/app`. `internal/app/usecases.go`, `project_usecases.go`, `study_usecases.go`, and `sprint_usecases.go` should expose result types that both CLI adapters and TUI code can consume without parsing rendered CLI output. These use cases should call existing `internal/project`, `internal/study`, `internal/sprint`, `internal/workspace`, and platform config/path services. They should not introduce global `validation`, `workflow`, `scheduler`, `reports`, `prompts`, or `stages` packages.

The TUI must not call CLI command handlers, invoke `ultraplan` as a subprocess, parse stdout/stderr, launch runtime work, run planning or execute workflows, mutate Git state, create review/smoke/issue artifacts, or create a separate persistence model. Durable workspace files remain the source of truth. TUI state is derived, refreshable, and disposable.

The current architecture fits with a small refactor-before-feature step: extract or add shared app use cases before wiring the TUI. A larger package restructuring is rejected because the project architecture already defines the desired module-driven shape and the sprint explicitly excludes broad workflow, validation, scheduler, report, prompt, or stage refactors.

The terminal UI dependency should be contained entirely inside `internal/tui`. Bubble Tea-style model/update/rendering is the preferred architectural shape because the handbook evidence favors deterministic message/model tests and library containment. The final library choice can happen during implementation, but terminal-library types must not appear in `internal/app`, product modules, platform packages, or typed use-case result structs.

Refresh should be explicit user-triggered refresh for Sprint 24, not background periodic refresh. This keeps lifecycle, cancellation, test determinism, and read-only guarantees simple while preserving `context.Context` propagation for later Sprint 25 operational controls. Startup should use an immutable workspace/config snapshot after normal discovery and validation; refresh should reload product status artifacts through use cases, not hot-reload runtime credentials or start runtime health checks.

Artifact previews should be implemented as read-only, bounded, workspace-contained reads exposed through app/product path-safe use cases or helpers. Preview state belongs in `internal/tui`, but path authorization, product artifact discovery, execute status derivation, and redaction must happen before rendering. Missing files, invalid paths, truncation, invalid JSON/Markdown, and redaction should be displayed as actionable read-only states, not repaired or silently ignored.

Execute display must come from typed sprint status/use-case results that interpret `.run-state.json`, `execute.md`, and `flow-state.json` safely. `internal/tui` must not parse `.run-state.json` directly as product logic. Any model, provider, runtime diagnostic, evidence, or config value shown in the TUI must be redacted before reaching rendering or must be represented as display-safe fields.

The implementation unit should be a small set of concrete modules and structs rather than broad abstractions. `internal/app` earns read-only use-case interfaces/results because CLI and TUI are two real consumers and tests need fakes. `internal/tui` earns a model/message boundary because deterministic non-interactive tests are required. Additional plugin, registry, scheduler, workflow-engine, or generic dashboard abstractions are rejected as premature and outside sprint scope.

## Trade-Offs

Shared typed app use cases are chosen over direct product service calls from `internal/tui`. The benefit is a stable local-interface boundary shared by CLI and TUI, avoiding duplicated workflow logic and CLI text scraping. The cost is a small layer of request/result types that must be maintained. This cost is justified because the sprint explicitly requires CLI and TUI to share use cases and because review must verify `internal/tui` depends only on app use cases and simple result types.

Bubble Tea-style model/update/rendering is favored over a widget-heavy tcell/tview-style architecture. The benefit is testable message handling, deterministic model transitions, and easier narrow-terminal rendering tests. The cost is that some dashboard widgets may need more manual rendering. This is acceptable because Sprint 24 values reliable read-only inspection and tests over rich operational controls.

Explicit constructor/factory dependency injection is chosen over context-based service lookup. The benefit is traceable dependencies, compile-time safety, and simple fake use cases in tests. The cost is more plumbing through TUI setup. This is acceptable because the handbook identifies context-as-service-locator as ergonomic but risky, while requirements demand deterministic tests and clear boundary review.

Immutable startup config plus explicit product-state refresh is chosen over hot-reloading config during the TUI session. The benefit is predictable startup behavior that matches CLI config precedence and avoids accidental runtime credential requirements. The cost is that users must restart the TUI to pick up config edits. This is acceptable for a read-only foundation because runtime-backed operation and session overrides are Sprint 25 or later concerns.

Explicit refresh is chosen over periodic background refresh. The benefit is simpler context lifecycle, no hidden filesystem churn, and deterministic model tests. The cost is less live-dashboard feel. This is acceptable because live runtime progress panes and operational monitoring are non-goals for Sprint 24.

Bounded preview reads are chosen over complete artifact display. The benefit is predictable memory and latency for large artifacts and repositories. The cost is that users may need to open files outside the TUI for complete content. This is acceptable if the TUI clearly marks truncation and provides artifact paths.

Read-only dashboard actions are chosen over guarded operational actions. The benefit is lower security, mutation, and workflow-coupling risk. The cost is that users cannot run validation, planning flow, execute, or study run-loop from the initial TUI. This is required by the sprint non-goals and preserves a clean handoff to later operational TUI work.

Rejected alternatives include building a separate TUI persistence cache, moving common planning/status behavior into global workflow or validation packages, parsing CLI JSON/text as an integration protocol, embedding terminal types in app/product results, and letting the TUI directly inspect execute run-state JSON as its own product logic. Each alternative would either duplicate ownership, weaken boundaries, increase coupling, or violate explicit sprint constraints.

## Evidence

The project architecture states that UltraPlan Go is module-driven, not global-layer-driven, and that behavior should live with the module that owns the state. It names `internal/app` as local composition and shared use-case wiring, while `internal/tui` owns only terminal widgets, navigation, key handling, and rendering. It also states that CLI and TUI should share application use cases and dependency construction instead of using CLI output as an integration protocol.

The sprint requirements require `ultraplan tui` to inspect workspace, project, study, and sprint state without invoking runtimes, mutating artifacts, shelling out to UltraPlan, or parsing CLI output. They also require typed app use cases, product service results, deterministic ordering, bounded previews, runtime-free startup, path-safe previews, redacted execute diagnostics, and tests for read-only guarantees.

The technical handbook's project-structure and command-architecture evidence supports thin entrypoints and one-way imports from command/UI layers into business logic. It specifically warns against TUI behavior hidden inside an interactive command interpreter and against CLI/TUI bifurcation through rendered text. This supports the decision to route TUI behavior through `internal/app` typed use cases.

The dependency-injection and IO-abstraction evidence supports composition roots, constructor/factory injection, injectable streams, terminal mocks, and fake dependencies. This supports app-level dependency assembly, fake read-only use cases, and non-interactive TUI tests rather than global mutable state or real terminal dependencies.

The terminal-UX evidence supports model/message separation, actionable state, interruptible full-screen interfaces, and graceful narrow-terminal behavior. This supports keeping TUI state and rendering in `internal/tui` and preferring a deterministic model/update architecture over widget code that leaks terminal-library concepts into product modules.

The state/context evidence supports root context propagation and explicit app/session state while warning against `context.Background()` inside operations and against context-based service locators. This supports passing `context.Context` through startup and refresh use cases while keeping side-effectful dependencies explicit.

The error-handling and observability evidence supports preserving error chains, classifying failures, separating user-facing rendering from operational detail, and keeping diagnostics structured and safe. This supports actionable TUI error panes backed by classified app/product errors, not raw runtime/config dumps.

The security evidence supports path containment, redaction, default-deny dangerous behavior, and explicit trust boundaries. This supports rejecting subprocess self-invocation, runtime-backed actions, path escapes, raw diagnostic display, and unredacted config/model/provider output.

The performance evidence supports lazy dependency setup, bounded reads, streaming discipline, and avoiding unbounded repository scans. This supports bounded artifact previews, known-artifact discovery through product use cases, and no eager runtime initialization during TUI startup.

The PRD and TRD both define the TUI as a local terminal surface over the same product services, not a browser/server surface or a CLI replacement. The TRD states the intended dependency direction as `cmd/ultraplan -> internal/app -> product modules/platform` and `internal/tui -> internal/app -> product modules/platform`, and explicitly forbids `internal/tui -> CLI command handlers -> stdout parsing` and `internal/tui -> os/exec("ultraplan", ...)`.

## Risks

There is a coupling risk if `internal/app` result types become too broad or mirror every product detail. Mitigation: keep result types use-case-oriented and display-safe, pass only fields needed by CLI/TUI, and keep product interpretation in `internal/project`, `internal/study`, and `internal/sprint`.

There is an ownership risk if `internal/tui` starts interpreting `flow-state.json`, `.run-state.json`, validation findings, or project indexes directly. Mitigation: require sprint/project/study status use cases to return typed summaries and make review check that product state machines and durable artifact semantics remain outside `internal/tui`.

There is a read-only boundary risk if existing status use cases perform deterministic refresh writes. Mitigation: avoid such use cases for Sprint 24 where possible; if an existing status operation can write, the TUI action label and tests must make the mutation scope explicit, matching the sprint acceptance criteria.

There is a redaction risk for Sprint 23 execute diagnostics and runtime metadata. Mitigation: expose only display-safe fields from sprint use cases and test that raw model, provider credential, environment, native diagnostic, and unsafe runtime payload values do not appear in TUI views.

There is a performance risk from previews and status aggregation over large workspaces. Mitigation: discover known product artifacts through use cases, sort deterministically, bound preview reads, avoid recursive source scans, and show truncation rather than loading complete files.

There is a UX risk that a read-only dashboard may appear less useful than users expect from a TUI. Mitigation: document the read-only scope in `internal/tui/doc.go`, key-binding help, and command help, and reserve validation/workflow controls for later guarded operational sprints.

There is a test fragility risk for terminal rendering. Mitigation: keep model/update behavior testable independently from rendering, use fake app use cases, prefer stable golden or targeted string assertions for critical status/error content, and include narrow-terminal cases.

There is an open implementation question about the exact terminal library. The architecture decision is that the library must support deterministic model/update/render tests and remain contained in `internal/tui`; Bubble Tea-style architecture is preferred by evidence, but the final dependency can be confirmed during implementation without changing product boundaries.

There is an open implementation question about the minimum narrow-terminal width and exact preview truncation size. These should be decided in implementation tests, but the architectural requirement is fixed: critical status and errors remain visible, previews are bounded, and truncation is explicit.
