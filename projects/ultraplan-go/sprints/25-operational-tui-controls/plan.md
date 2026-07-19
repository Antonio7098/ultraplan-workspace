# Sprint Plan: Operational TUI Controls

> Project: `ultraplan-go`
> Sprint: `25-operational-tui-controls`
> Source: roadmap Sprint 25 scope, repository PRD and architecture, applicable system contracts, and the implemented Sprint 24 TUI boundary
> **Inputs Used:** `projects/ultraplan-go/roadmap.md`, `projects/ultraplan-go/project-index.md`, `system/protocols/plan-sprint-protocol.md`, `system/contracts/core/architecture.md`, `system/contracts/core/errors.md`, `system/contracts/core/configuration.md`, `system/contracts/core/observability.md`, `system/contracts/core/security.md`, `system/contracts/core/testing.md`, `system/contracts/surfaces/cli.md`, `system/contracts/runtime/llm.md`, `system/contracts/runtime/workflows.md`, `system/contracts/runtime/performance.md`, `system/contracts/runtime/persistence-and-migrations.md`, `../ultraplan-go/ARCHITECTURE.md`, `../ultraplan-go/PRD.md`, `projects/ultraplan-go/sprints/24-tui-foundations/plan.md`, `../ultraplan-go/internal/app/usecases.go`, `../ultraplan-go/internal/app/tui_commands.go`, `../ultraplan-go/internal/app/sprint_commands.go`, `../ultraplan-go/internal/app/study_commands.go`, `../ultraplan-go/internal/app/study_runs_commands.go`, `../ultraplan-go/internal/app/sprint_usecases.go`, `../ultraplan-go/internal/sprint/flow.go`, `../ultraplan-go/internal/sprint/execute.go`, `../ultraplan-go/internal/study/run_loop.go`, `../ultraplan-go/internal/tui/app.go`, `../ultraplan-go/internal/tui/model.go`, `../ultraplan-go/internal/tui/keys.go`, `../ultraplan-go/internal/tui/views.go`

This is the only Sprint 25 planning artifact requested. No `requirements.md`, `sprint-index.md`, `technical-handbook.md`, or reasoning documents are created by this plan. The plan derives its scope directly from the roadmap and preserves unresolved implementation details as explicit design gates.

## Sprint Status

- **Status:** `complete`
- **Owner:** `implementation agent`
- **Dependencies:** Sprint 24 read-only TUI and shared dashboard use cases; existing sprint validation, prompt, flow, and execute services; existing study validation and run-loop services; generic runtime adapter and durable state formats.

## Goal

Allow users to safely validate local artifacts, preview prompts and dry runs, run selected sprint planning stages, enter execute workflows, and start or resume study run loops from the TUI. Every mutating or runtime-backed action must show its exact scope and require confirmation; progress and cancellation must use typed app-level signals and existing product services rather than CLI handlers or rendered terminal text.

## Scope

### In Scope

- Typed operational use cases in `internal/app` shared by CLI and TUI adapters.
- Project, study, and sprint validation actions.
- Sprint flow dry-run and prompt preview.
- Guarded planning flow execution through a selected stage.
- Execute status plus guarded dry-run, start, and resume entry points.
- Study run-loop start and resume with source/dimension filters and bounded parallelism.
- Typed progress/event panes, loading state, cancellation, terminal outcome, and post-operation refresh.
- Classified, redacted error details with stable code, operation, component, cause/guidance, and retryability where available.
- Deterministic non-interactive model/update/view tests using fake operational use cases and fake runtimes.

### Non-Goals

- Full workflow cockpit polish, related-action shortcuts, run history/cost dashboards, execute task-detail inspection, external-edit recovery, and release documentation reserved for Sprint 26.
- Smoke investigation, automated review, issue tracking, Git add/commit/push, plugins, hosted services, browser UI, or multi-user collaboration.
- A new workflow engine, TUI-owned scheduler, TUI-owned durable state, alternate runtime abstraction, or duplication of product workflow semantics in `internal/tui`.
- Changing documented CLI flags, text/JSON schemas, exit behavior, or making the TUI the automation surface.
- Concurrent execution of multiple foreground operational TUI actions.

## Architecture And Safety Constraints

1. `cmd/ultraplan` remains thin; composition remains in `internal/app`.
2. `internal/project`, `internal/sprint`, and `internal/study` continue to own validation and workflow semantics. `internal/platform/runtime` remains generic and must not import product or TUI packages.
3. `internal/tui` owns dialogs, local navigation/form state, rendering, and translation of typed progress into view state only. It must not invoke CLI handlers, parse stdout/stderr, shell out to `ultraplan`, interpret provider-native payloads, or write product artifacts directly.
4. Bubble Tea and Glamour types remain contained inside `internal/tui`; app and product APIs exchange plain Go request/result/event types.
5. Durable workspace artifacts remain the source of truth. Existing atomic/versioned persistence rules remain intact; Sprint 25 introduces no TUI state file or cache.
6. Exactly one foreground operation is owned by the TUI at a time. It has an operation-scoped `context.Context`, cancel function, bounded event channel/buffer, and explicit terminal result. Cancellation is requested once, waits for the owned operation to terminate, then refreshes durable status.
7. Validation and prompt preview are non-runtime actions but may read or deterministically refresh documented artifacts. Dry-run must not construct a real runtime or perform runtime-backed writes.
8. A mutating or runtime-backed operation cannot start until a confirmation view displays operation kind, project/study/sprint, selected stage or task scope, affected workspace paths, filters, parallelism/model source when relevant, and whether durable files may change.
9. Errors and events shown in the TUI must be classified and display-safe. Secrets, environment values, raw provider payloads, unrestricted prompts, and unredacted runtime diagnostics must not be exposed.
10. All fan-out, event collection, runtime calls, refreshes, preview sizes, and rendered event history are bounded. No background polling or detached goroutine may outlive its owning operation.

## Applicable Contract Map

| Contract IDs | Sprint Obligation | Planned Evidence |
| --- | --- | --- |
| `ARCH-CORE-001`, `ARCH-CORE-002`, `ARCH-LAYER-002`, `ARCH-ENTRY-001`, `ARCH-COMP-001`, `ARCH-SHARED-001` | Keep workflow ownership in product modules, expose narrow app use cases, centralize runtime composition, and keep terminal/vendor types out of app/product packages. | Import review, constructor/use-case tests, and source searches for forbidden dependencies. |
| `ERR-CORE-001`, `ERR-CODE-001`, `ERR-TRANS-001`, `ERR-STARTUP-001`, `ERR-TASK-001`, `ERR-DEP-001`, `ERR-DATA-001`, `ERR-REDACT-001`, `ERR-USER-001` | Preserve causes and stable classification; make every action end complete, failed, cancelled, or partial; separate concise guidance from inspectable redacted detail. | Fake failure matrices, cancellation tests, error-pane tests, and redaction assertions. |
| Configuration contract | Reuse CLI config precedence and validation; expose effective model/runtime/parallelism safely; never print secrets. | App composition tests and CLI/TUI parity fixtures. |
| Observability contract | Convert existing canonical runtime/product progress into typed operation events with run/task/stage identity and terminal outcomes. | Ordered event tests and progress-pane model/view tests. |
| Security contract | Confirm paths and scope before execution; enforce workspace/target containment, runtime permission policy, and redaction. | Path-escape tests, confirmation snapshots, permission-scope tests, and source review. |
| Testing contract | Default to deterministic fakes and temporary workspaces; keep real-runtime smoke gated. | Unit, model, view, integration-fixture, race, and full-suite results. |
| CLI Surface contract | Keep CLI behavior stable and reuse the same app operations without routing TUI work through CLI rendering. | Existing CLI regression tests plus shared-use-case parity tests. |
| LLM Runtime contract | Use the configured generic runtime boundary, propagate context, retain validation and policy wrappers, and expose safe progress/cost drivers where available. | Fake-runtime request/event assertions and gated smoke instructions. |
| Workflow contract | Give long-running work one owner, explicit state transitions, bounded retry/resume behavior, idempotent existing workflow semantics, and inspectable cancellation. | Lifecycle/state-transition tests and durable resume/cancel fixtures. |
| `PERF-BOUND-001`, `PERF-UI-001`, `PERF-CONC-001`, `PERF-COST-001` | Bound parallelism and retained events; keep rendering responsive; surface model/provider/request/task/duration drivers safely. | Bound tests, race tests, large-event fixture, and render responsiveness review. |
| `PERSIST-SCHEMA-001`, `PERSIST-ATOMIC-001`, `PERSIST-INTEGRITY-001`, `PERSIST-DERIVED-001`, `PERSIST-RECOVERY-001` | Product modules alone mutate their owned versioned artifacts; TUI state remains derived; cancellation and failure leave recoverable state. | Before/after artifact fixtures, atomic-write tests already owned by product modules, and recovery tests. |

## Acceptance Criteria

- **AC1:** TUI validation actions exist for project, study, and every supported sprint artifact stage, and results render typed findings without calling CLI handlers.
- **AC2:** Sprint prompt preview and flow dry-run are available without constructing a runtime or mutating runtime-backed artifacts; preview content is bounded and redacted according to its sensitivity.
- **AC3:** A user can select a sprint planning target through `plan`, review its ordered stage/path scope, confirm it, observe typed progress, and inspect its terminal result.
- **AC4:** Execute status is refreshable, and execute dry-run/start/resume entry points use the existing sprint execute service, validation gates, approved target containment, permission policy, and durable run state.
- **AC5:** A user can start or resume a study run loop with normalized source/dimension filters and validated bounded parallelism.
- **AC6:** Every mutating or runtime-backed action requires an explicit confirmation step showing affected paths and stage/task scope. Cancel/back from confirmation performs no operation and no product-state mutation.
- **AC7:** Running operations update from typed progress/events or product results, never CLI text, stdout/stderr scraping, or provider-native terminal output.
- **AC8:** Pressing cancel propagates through the exact operation context, produces an explicit cancelled or already-terminal outcome, waits for owned work to settle, and refreshes status from durable artifacts.
- **AC9:** Error panes expose a stable display-safe code/category, operation, component, concise message, retryability, cause summary, and actionable guidance when available; secret values and raw runtime payloads never render.
- **AC10:** Only one foreground operation can run; conflicting actions are disabled with a clear explanation, event retention is bounded, and no owned goroutine or channel is leaked.
- **AC11:** CLI commands and documented text/JSON behavior remain unchanged, while shared operational app use cases are exercised by both adapters where practical.
- **AC12:** Narrow-terminal layouts preserve confirmation scope, destructive/runtime warning, progress status, cancellation state, and critical errors in plain text.
- **AC13:** Normal tests are non-interactive, offline, and use fake runtimes; `go test -race ./internal/app ./internal/tui ./internal/sprint ./internal/study`, `go test ./...`, and `go build ./cmd/ultraplan` pass.

## Tasks

- [x] **Task 1: Define The Operational Use-Case And Event Contracts**
  > Satisfies: `AC1-AC11`, architecture, errors, observability, workflows
  - [x] Add the first narrow plain-Go operational interface and request/result types in `internal/app` for project, study, and sprint-stage validation; compose `ReadOnlyUseCases` into the broader TUI-facing facade. Preview/dry-run/flow/execute/run-loop contracts remain pending.
  - [x] Define an operation identity, lifecycle states (`preparing`, `running`, `cancelling`, `complete`, `failed`, `cancelled`, `partial`), typed progress records, terminal result, and display-safe classified error detail.
  - [x] Define bounded event delivery and ownership semantics: ordering, buffer/retention limit, terminal event guarantee, slow-consumer behavior, and channel closure.
  - [x] Define confirmation metadata returned before execution: operation, subject refs, stages/tasks, affected paths, runtime/mutation flags, filters, parallelism, model source, and permission summary.
  - [x] Keep runtime/provider, Bubble Tea, CLI renderer, and durable schema types out of the public app contract.

- [x] **Task 2: Extract Shared Validation Operations From CLI Adapters**
  > Satisfies: `AC1`, `AC9`, `AC11`
  - [x] Add app operations for project validation, study validation, and sprint stage validation using existing product services.
  - [x] Return deterministic typed, redacted findings and mapped errors rather than pre-rendered strings.
  - [x] Refactor CLI validation handlers only as needed to call the same app operations while preserving existing stdout, stderr, JSON, and exit statuses.
  - [x] Add parity tests proving equivalent product results for CLI and TUI callers and negative tests for missing/invalid workspace artifacts.

- [x] **Task 3: Add Sprint Prompt Preview And Dry-Run Operations**
  > Satisfies: `AC2`, `AC6`, `AC7`, security, performance
  - [x] Expose supported planning-stage prompt previews and flow/execute dry runs through typed app requests.
  - [x] Reuse product preflight, stage ordering, artifact validation, model selection, target-path checks, and config precedence.
  - [x] Prove dry-run does not construct or call a real runtime and document any deterministic artifact refresh it may perform.
  - [x] Bound prompt display and preserve an explicit truncation indicator; redact configuration/runtime metadata while retaining useful scope.

- [x] **Task 4: Add Runtime-Backed Sprint Flow And Execute Operations**
  > Satisfies: `AC3`, `AC4`, `AC7-AC10`, runtime and workflow contracts
  - [x] Move reusable sprint flow orchestration out of CLI-only control flow into an app use case or product-owned service without changing product module ownership.
  - [x] Support selected flow targets through `plan`, skipping already-valid stages according to existing behavior and emitting typed stage start/progress/terminal events.
  - [x] Expose execute status plus dry-run/start/resume requests using the existing execute validation, plan-task extraction, target containment, model resolution, permission policy, durable state, and evidence checks.
  - [x] Centralize runtime construction in app composition and apply the same observing, validating, and policy wrappers used by CLI runtime operations.
  - [x] Ensure partial failures and retry exhaustion remain visible and durable, and never translate runtime success alone into task completion.

- [x] **Task 5: Add Study Run-Loop Operational Use Cases**
  > Satisfies: `AC5`, `AC7-AC10`, workflows, performance, persistence
  - [x] Expose study validation and run-loop start/resume through app request/result types backed by `internal/study` services.
  - [x] Accept source and dimension filters using existing resolution/normalization rules and reject unknown or inapplicable selections before runtime work.
  - [x] Validate parallelism against configured practical bounds; retain existing lock, retry, fallback, synthesis gating, run-history, and atomic run-state behavior.
  - [x] Adapt existing run-loop/runtime events or progress callbacks into typed app events without duplicating the study scheduler in `internal/app` or `internal/tui`.
  - [x] Treat an existing active lock, stale lock, resume mismatch, partial failure, and cancellation as distinct inspectable outcomes with guidance.

- [x] **Task 6: Implement Confirmation Policy And Operation-Scoped TUI State**
  > Satisfies: `AC6`, `AC8`, `AC10`, `AC12`
  - [x] Add TUI model states for action selection, parameter entry, scope review, confirmation, running, cancelling, terminal result, and refreshed detail.
  - [x] Require an affirmative confirmation key distinct from ordinary navigation; back/escape from parameter or confirmation views must be side-effect free.
  - [x] Display exact affected workspace paths, stages/tasks, runtime/mutation status, filters, parallelism, and safe model/permission source before execution.
  - [x] Permit only one foreground operation and disable navigation/actions that could conflict while it runs, while keeping cancel and progress inspection available.
  - [x] Keep all form/dialog state disposable and derived; introduce no TUI persistence.

- [x] **Task 7: Add Operational Navigation, Forms, And Guarded Actions**
  > Satisfies: `AC1-AC6`, `AC12`
  - [x] Extend project, study, and sprint routes with validation action entries without replacing existing artifact navigation.
  - [x] Add project/study/sprint validation actions and a result pane with finding guidance.
  - [x] Add sprint stage selection for prompt preview, dry-run, planning flow, execute status, execute dry-run, execute start, and execute resume.
  - [x] Add study run-loop parameter controls for start/resume, source filters, dimension filters, and parallelism with inline validation.
  - [x] Make runtime-backed and mutating actions visually and textually distinct without relying on color or ANSI styling alone.

- [x] **Task 8: Add Typed Live Progress, Cancellation, And Refresh**
  > Satisfies: `AC7`, `AC8`, `AC10`, observability, workflows, performance
  - [x] Bridge typed app events into Bubble Tea commands/messages with operation IDs so stale events cannot update a newer action.
  - [x] Render bounded event history plus current stage/task, counts, attempt, duration, and safe runtime/model/cost drivers when available.
  - [x] Wire cancel to the operation-scoped cancel function, transition immediately to `cancelling`, reject repeated cancellation safely, and await the terminal event.
  - [x] Refresh dashboard/status from durable artifacts after complete, failed, partial, or cancelled outcomes without hiding the terminal summary.
  - [x] Ensure terminal exit or TUI quit during an operation requests cancellation and does not abandon owned work silently.

- [x] **Task 9: Add Classified Error And Guidance Panes**
  > Satisfies: `AC9`, `AC12`, error and security contracts
  - [x] Render a concise user message separately from optional operator detail.
  - [x] Preserve stable code/category, component, operation, retryability, safe cause chain, relevant subject IDs, and actionable guidance through app-to-TUI translation.
  - [x] Add explicit views for configuration/preflight, validation, path/security, dependency/runtime, timeout, concurrency/lock, persistence, cancellation, partial, and unexpected internal failures.
  - [x] Apply centralized redaction before values reach terminal rendering; add adversarial fixtures containing tokens, environment values, provider payloads, prompts, paths outside scope, and multiline diagnostics.

- [x] **Task 10: Add Deterministic Tests And Contract Regression Coverage**
  > Satisfies: `AC1-AC13`, testing, security, workflows, performance
  - [x] Extend TUI fakes with controllable event streams, blocking operations, terminal results, mutation counters, context recording, and cancellation acknowledgements.
  - [x] Add model/update tests for every dialog transition, invalid parameter, confirmation accept/reject, stale event, progress sequence, cancellation race, terminal result, refresh, and single-operation guard.
  - [x] Add view tests for confirmation scope, prompt truncation, progress, errors/guidance, filters, execute/run-loop states, plain-text warnings, and narrow terminals.
  - [x] Add app integration fixtures for validation, flow dry-run, selected planning flow, execute dry-run/resume, study start/resume, durable state after cancellation, and CLI behavior parity.
  - [x] Add race/leak-oriented tests for event delivery, cancellation, slow consumers, full buffers, runtime failures, and TUI quit during owned work.
  - [x] Prove no TUI path invokes CLI handlers, parses terminal output, shells out to `ultraplan`, mutates Git, or reaches provider/native types directly.

- [x] **Task 11: Update User-Facing Help And Package Documentation**
  > Satisfies: documentation, `AC6`, `AC11`
  - [x] Update `ultraplan tui --help`, in-TUI key help, and `internal/tui/doc.go` for operational scope, confirmation semantics, cancellation behavior, and durable state ownership.
  - [x] Document which validation/dry-run actions are non-runtime, which operations mutate artifacts, which paths can change, and how to inspect/recover with equivalent CLI commands.
  - [x] Keep full cockpit/release documentation deferred to Sprint 26, but provide enough operational guidance to prevent unsafe or ambiguous use in Sprint 25.

- [x] **Task 12: Verify Architecture, Safety, And Release Readiness**
  > Satisfies: `AC13` and all selected contracts
  - [x] Run targeted app/TUI/product tests, the race-enabled concurrency suite, the full Go suite, and binary build for the implemented validation slice.
  - [x] Review import direction, runtime construction, event ownership, cancellation propagation, path containment, durable writes, redaction, and terminal dependency containment.
  - [x] Exercise a deterministic temporary-workspace scenario covering validation, prompt preview, dry-run, confirmed planning operation, cancellation, refresh, and CLI inspection of resulting state.
  - [x] If a real runtime is configured, run only the separately gated smoke path after normal tests pass; record it as optional evidence rather than a normal-test dependency.
  - [x] Record deviations from this plan before implementation continues and defer Sprint 26 cockpit features rather than expanding Sprint 25.

## Verification Commands

Run from `/home/antonioborgerees/coding/ultraplan-go`.

| Check | Command | Expected Result |
| --- | --- | --- |
| App and TUI tests | `go test ./internal/app ./internal/tui` | Operational use cases, dialogs, typed events, errors, confirmation, and rendering pass offline and non-interactively. |
| Product workflow tests | `go test ./internal/project ./internal/sprint ./internal/study` | Existing validation, planning, execute, run-loop, persistence, and cancellation behavior remains correct. |
| Race-enabled operational suite | `go test -race ./internal/app ./internal/tui ./internal/sprint ./internal/study` | No data race, leaked ownership, unsafe event delivery, or cancellation race is detected. |
| Full test suite | `go test ./...` | All packages pass without real credentials, network access, OpenCode, or a real terminal. |
| Binary build | `go build ./cmd/ultraplan` | CLI builds with operational TUI controls and contained terminal dependencies. |
| CLI regression | `go test ./internal/app -run 'Test.*(Sprint|Study|Project|TUI)'` | Existing command parsing, output, errors, and exit behavior remain stable. |
| Forbidden integration search | `rg -n 'os/exec|exec\.Command|runSprint\(|runStudy\(|fmt\.(Fprint|Fprintf)|tea\.|glamour' internal/tui internal/app internal/project internal/sprint internal/study` | Matches are reviewed: no self-subprocess or CLI-handler calls from TUI; rendering/vendor types remain in `internal/tui`; app/product operation results are plain Go data. |
| Architecture review | `review using ../ultraplan-go-workspace/system/protocols/architecture-review-protocol.md` | Module ownership, dependency direction, composition, progress seams, and TUI dependency containment conform. |
| Sprint review | `review using ../ultraplan-go-workspace/system/protocols/sprint-review-protocol.md` | Acceptance criteria, contract map, verification evidence, risks, and non-goals are evaluable without guessing. |

## Execution Evidence

### 2026-07-10 — validation-control vertical slice

- Added `app.OperationalUseCases`, typed validation requests/results, and a shared dispatcher backed directly by `project.Service`, `study.Service`, and `sprint.Service`.
- Added project, study, and sprint-plan validation entries to existing TUI routes, asynchronous Bubble Tea command delivery, a typed findings/guidance pane, and side-effect-free escape/back behavior.
- Kept validation values display-safe through the existing centralized finding translations and configuration redaction.
- Updated TUI composition/help and deterministic model/view coverage. Existing artifact navigation and CLI validation handlers were not changed.
- Passed: `GOCACHE=/tmp/ultraplan-go-cache go test ./internal/app ./internal/tui`.
- Passed: `GOCACHE=/tmp/ultraplan-go-cache go test ./...`.
- Passed: `GOCACHE=/tmp/ultraplan-go-cache go test -race ./internal/app ./internal/tui ./internal/sprint ./internal/study`.
- Passed: `GOCACHE=/tmp/ultraplan-go-cache go build ./cmd/ultraplan`. Go emitted a non-fatal warning because its module stat cache is outside the writable sandbox; exit status was zero.
- This initial deferral was subsequently closed by the completion pass below.

### 2026-07-10 — completion pass

- Added typed operation kinds, lifecycle events/results, confirmation metadata, bounded content/event retention, classified display errors, and configured runtime composition.
- Added all supported sprint-stage validation and prompt actions, plan flow dry-run/run, execute status/dry-run/start/resume, and study run-loop start/resume. Runtime/mutating actions show explicit scope/path warnings and require `y` confirmation.
- Added single foreground-operation ownership, operation-scoped cancellation on quit, terminal results, bounded progress history, and durable dashboard refresh after settlement.
- Runtime work delegates to the existing sprint and study services with the same configuration/runtime factories as CLI paths; TUI contains no CLI calls, subprocess invocation, provider types, workflow persistence, or Git mutation.
- Passed full offline suite, race suite, binary build, whitespace check, and forbidden-integration search (no `os/exec`, `exec.Command`, `runSprint`, or `runStudy` matches in `internal/tui`).

## Review And Sign-Off

- **Status:** complete
- **Completed:** 2026-07-10
- **Blockers:** none
- **Deferrals:** Sprint 26 cockpit polish, history/cost dashboards, related-action shortcuts, and release documentation remain outside this sprint.

## Risks And Mitigations

| Risk | Mitigation | Completion Evidence |
| --- | --- | --- |
| Existing runtime-backed app use cases do not expose progress callbacks or cancellation. | Make typed progress/cancellation the first implementation gate; adapt canonical runtime/product events at the app boundary and do not simulate progress from text. | Fake-runtime event ordering and cancellation tests pass before UI action work is accepted. |
| Moving CLI orchestration could change established CLI behavior. | Extract shared operations behind existing renderers incrementally and lock text/JSON/exit behavior with regression fixtures. | Existing and new CLI tests pass unchanged. |
| TUI cancellation may return before durable state settles. | Keep operation ownership until the product service emits a terminal outcome, then reload durable state before enabling another action. | Cancel-and-resume fixtures are inspectable from both TUI results and CLI status. |
| Event floods can freeze rendering or leak memory. | Bound producer delivery and retained history, coalesce high-frequency progress where safe, and preserve terminal/error events. | Large/slow-consumer tests remain responsive and retain terminal truth. |
| Confirmation text may understate the mutation scope. | Generate scope metadata from the same app request/preflight used for execution; test exact affected paths and tasks. | Confirmation fixtures match actual changed artifacts. |
| Prompt/error/event panes may leak secrets or raw provider data. | Redact and normalize at app/product boundaries; render only display-safe typed fields and bounded content. | Adversarial redaction tests pass. |
| TUI code may duplicate sprint/study workflow state machines. | Limit `internal/tui` to interaction state and delegate every product transition to existing services. | Source/import review finds no product artifact writes or scheduler logic in `internal/tui`. |
| Run-loop filters or parallelism could create unbounded work. | Reuse product normalization and scheduling, validate practical bounds before confirmation, and display task count/cost drivers. | Bounds and invalid-filter tests pass. |
| Existing durable formats may not contain every desired display field. | Do not change formats merely for UI convenience; expose unavailable fields explicitly or make a separately reviewed compatible schema change. | Persistence review records no accidental schema drift. |

## Review Inputs

- This `plan.md` and the Sprint 25 roadmap section.
- Repository `ARCHITECTURE.md` and `PRD.md`.
- Applicable contracts listed in the contract map.
- Sprint 24 TUI implementation and tests.
- Implementation diff and verification output.
- Temporary-workspace lifecycle evidence and optional gated runtime smoke evidence.
- Architecture and sprint review protocols.

## Completion Criteria

- [x] All acceptance criteria have direct automated or review evidence.
- [x] All mutating/runtime actions are scoped, explicitly confirmed, cancellable, and terminally inspectable.
- [x] Validation, prompt preview, dry-run, planning flow, execute, and study run-loop controls use typed shared app operations and existing product services.
- [x] Progress and errors use typed display-safe data; no CLI text scraping or raw provider rendering exists.
- [x] Cancellation and TUI quit leave durable state inspectable and resumable from both TUI and CLI.
- [x] CLI behavior and automation surfaces remain unchanged.
- [x] No TUI-owned durable state, workflow engine, Git mutation, smoke/review automation, or Sprint 26 cockpit scope is introduced.
- [x] Targeted, race-enabled, full-suite, and build verification pass, or any blocker is recorded with actionable evidence.
- [x] Architecture and sprint reviews can evaluate the implementation without inventing intent.
