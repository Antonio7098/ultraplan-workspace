# UltraPlan Go Architecture

## Architecture Philosophy

UltraPlan Go is module-driven, not global-layer-driven.

The guiding question is not:

> What technical category does this file belong to?

It is:

> Which module owns this behavior and state?

A product module should encapsulate a complete slice of behavior:

```text
module = state + logic + workflows + validation + persistence adapters + local interface surface
```

This means study behavior should stay with the study module. Project catalog behavior should stay with the project module. Sprint planning, execute, review, and smoke behavior should stay with the sprint module. Code extraction behavior should stay with the code extraction module. Workspace behavior should stay with the workspace module. Shared platform packages should exist only for genuinely cross-cutting infrastructure.

CLI and TUI are local interfaces over the same product core. They should share application use cases and dependency construction instead of duplicating workflow logic or using CLI output as an integration protocol.

## Core Rule

Prefer this:

```text
internal/study owns its validation, scheduling, reports, prompts, state, and persistence
internal/project owns project docs, project-index cataloging, and project validation
internal/sprint owns planning artifacts, flow state, prompt rendering, stage validation, controlled implementation through execute, automated review, and external-harness-backed smoke
internal/codeextract owns its parsing, resolution, extraction, and validation behavior
internal/workspace owns workspace discovery, path rules, and workspace validation
internal/platform/runtime owns generic execution only
internal/platform/process owns safe generic external executable/argv execution only
internal/app owns local composition and shared use-case wiring for CLI and TUI
internal/tui owns terminal widgets, navigation, key handling, and rendering only
```

Avoid this as the default shape:

```text
internal/validation
internal/scheduler
internal/reports
internal/prompts
internal/study
```

Global technical-layer packages look clean at first, but they tend to fracture product context and create cross-module coupling.

## Target Layout

Use a pragmatic Go layout: one package per product module, split by focused files first, and introduce subpackages only when the module grows enough that the boundary improves comprehension.

```text
cmd/
  ultraplan/
    main.go

internal/
  app/
    app.go                  # composition root, dependency wiring, CLI adapters
    usecases.go             # shared operations used by CLI and TUI, introduced as needed

  tui/
    app.go                  # Bubble Tea or equivalent program setup
    model.go                # TUI state/update model
    views.go                # terminal rendering
    keys.go                 # key bindings and commands

  platform/
    config/
    logging/
    filesystem/
    runtime/
      runtime.go            # generic execution interface
      agentwrap.go          # agentwrap integration
      opencode.go           # opencode-specific adapter, if needed
    process/
      runner.go             # executable/argv/cwd/env/timeout/cancellation boundary

  workspace/
    domain.go
    discovery.go
    validation.go
    paths.go
    store.go

  study/
    domain.go               # Study, Source, Dimension, Report, RunState
    service.go              # use-case coordination
    init.go
    run.go
    run_all.go
    synthesize.go
    scheduler.go
    state.go
    prompts.go
    validation.go
    reports.go
    summary.go
    store_fs.go
    cli.go

  project/
    domain.go               # Project, ProjectIndex, catalog entries
    discovery.go
    validation.go
    service.go
    store_fs.go
    cli.go

  sprint/
    domain.go               # Sprint, PlanningStage, FlowState
    flow.go
    prompts.go
    validation.go
    state.go
    service.go
    store_fs.go
    cli.go
    review.go               # review scope, reviewer orchestration, verdict, review.md
    review_validation.go
    smoke.go                # harness discovery, scope, invocation, verdict, smoke.md
    smoke_validation.go

  codeextract/
    domain.go
    service.go
    parser.go
    resolver.go
    validation.go
    cli.go
```

Do not immediately create a large clean-architecture tree such as:

```text
internal/study/domain/
internal/study/app/
internal/study/ports/
internal/study/adapters/
internal/study/validation/
internal/study/reports/
internal/study/prompts/
```

That can recreate the same abstraction problem as global layers. Start with one package per module and multiple focused files. Split into subpackages later when there is a concrete readability or dependency benefit.

## Module Ownership

### `study`

`study` is the main product module for the current build. It owns the full study lifecycle:

```text
Study definitions
Sources
Dimensions
Source applicability
Prompt creation
Single analysis runs
Full study runs
Run-loop state
Scheduling
Per-source report paths
Final synthesis report paths
Study validation
Report parsing
Summary generation
Filesystem persistence for study artifacts
Local interface commands for study workflows
```

Prefer:

```text
internal/study/validation.go
internal/study/scheduler.go
internal/study/reports.go
internal/study/prompts.go
```

Instead of:

```text
internal/validation/study.go
internal/scheduler/study.go
internal/reports/study.go
internal/prompts/study.go
```

The study module may call platform runtime services, workspace path services, and config/logging infrastructure. Runtime and platform packages must not import `study`.

### `project`

`project` owns the planning root under `projects/<project>`:

```text
Project discovery and resolution
Project docs discovery
Roadmap discovery
project-index.md parsing
Catalog entries for contracts, evidence, reasoning templates, and review protocols
Project validation
Project status output
Filesystem persistence for project catalog artifacts
Local interface commands for project workflows
```

`project` may call workspace path services and config/logging infrastructure. It should not import `study`; study outputs are referenced as catalog paths, not consumed through study services.

### `sprint`

`sprint` owns governed planning, execute, review, and smoke artifacts under `projects/<project>/sprints/<slug>` through `smoke`:

```text
requirements.md
sprint-index.md
technical-handbook.md
reasoning/*.md
reasoning.md
plan.md
execute.md
flow-state.json
.run-state.json
review.md
smoke.md
```

It owns:

```text
Planning, execute, review, and smoke stage order through smoke
Stage validation
Sprint-index subset checks against project-index.md
Technical handbook prompt/rendering behavior
Area reasoning prompt/rendering behavior
Final reasoning prompt/rendering behavior
Plan prompt/rendering behavior
Execute prompt/rendering behavior
Plan task extraction and deterministic task IDs
Sprint stage model-resolution rules
Flow state persistence
Execute run-state persistence
Review scope, selected-contract/protocol resolution, structured reviewer orchestration, deterministic verdicts, and review.md
Smoke review-gating, external harness selection/invocation, evidence-link validation, deterministic verdicts, and smoke.md
Sprint status output
Local interface commands for planning-stage sprint workflows
```

`sprint` may depend on `project` for project catalogs, `workspace` for path rules, configuration for stage-specific runtime model selection, `platform/runtime` for generic prompt/reviewer/implementation execution, and `platform/process` for safe external harness invocation. It must not depend on `study` services, source/dimension models, study report validation, rating parsing, summary generation, or study run-loop scheduling.

Phase 2 includes controlled implementation execution through `execute`. Phase 3 extends the same sprint module through `review` and `smoke`. It keeps only the current `review.md` and `smoke.md` in the sprint root; detailed smoke runs and issue evidence remain in the external harness. `sprint` must not model general-purpose issue tracking, automatic product fixes, or Git mutation as workflow stages.

### `app`

`app` is the local composition and use-case boundary. It owns:

```text
CLI argument parsing and command dispatch
shared dependency construction
workspace/config/runtime preflight wiring
typed use-case functions for local interfaces
text and JSON output adapters for CLI commands
stable error classification and process exit mapping
review/smoke progress, cancellation, result, evidence-link, and recovery use cases shared by CLI and TUI
```

As the TUI enters scope, command handlers should be thinned so the real operation lives in shared app use cases or product services. The desired shape is:

```text
CLI command -> parse flags -> call app use case -> render text/JSON
TUI action  -> build request -> call app use case -> update model/view
```

Avoid letting the TUI call CLI command functions that write to stdout, or shelling out to the `ultraplan` binary to parse its own output.

### `tui`

`tui` owns the interactive terminal experience:

```text
dashboard navigation
lists, panes, filters, and detail views
key bindings
terminal rendering and resize behavior
read-only artifact previews
guarded action prompts for mutating/runtime operations
progress/event display
review scope, reviewer progress, findings, verdict, reruns, and review.md preview
smoke scope, prerequisites, suite/test progress, linked harness evidence, issues, verdict, reruns, and smoke.md preview
```

It may depend on `app` use cases and domain result types. It must not own product state machines, validation rules, prompt construction, runtime execution, smoke-harness invocation, verdict synthesis, or artifact persistence. Durable workspace files and linked harness evidence remain the source of truth.

If a third-party terminal UI library is used, it should stay behind `internal/tui` and not leak into product modules.

### `platform/runtime`

Runtime is platform-level because it is generic execution infrastructure, not study behavior.

It may understand:

```text
Prompt
Working directory
Model
Timeout
Environment
Permissions
Expected output path
Execution events
Execution result
```

It must not understand:

```text
Study
Project
Sprint
Dimension
Source
Synthesis gating
Report semantics
Study state machines
Project catalog semantics
Sprint stage semantics
Summary generation
```

The dependency direction is:

```text
study -> platform/runtime
sprint -> platform/runtime
platform/runtime -> no product modules
```

### `platform/process`

The process package is a generic volatile-boundary adapter used by Phase 3 to run the cataloged smoke harness. It may understand:

```text
Executable
Argument vector
Working directory
Allowlisted environment
Timeout and context cancellation
Stdout/stderr capture limits
Exit code and process-tree cleanup
```

It must not understand projects, sprints, reviews, smoke levels, harness issue semantics, verdicts, or `review.md`/`smoke.md`. Those remain in `internal/sprint`.

The dependency direction is:

```text
sprint -> platform/process
platform/process -> no product modules
```

Runtime supervision is delegated to `agentwrap`. UltraPlan's runtime integration should translate generic execution requests into agentwrap requests and translate generic results/events back to UltraPlan's platform-facing runtime model.

### `workspace`

`workspace` owns where UltraPlan is operating and how workspace paths are resolved.

It owns:

```text
Root discovery
Workspace marker/config lookup
Canonical workspace paths
Workspace-level validation
Path safety rules
Embedded prompt/template default registry and override resolution
Opt-in default materialization
```

Product modules own the semantics of their prompts and generated artifacts. The workspace package owns the mechanical asset policy:

```text
readable workspace file at prompts/<name> or templates/<name>
  -> intentional local override
otherwise
  -> embedded binary default
```

`init-workspace` must not export prompt or template files. `ultraplan defaults install` is the explicit operation for materializing editable copies. A workspace without `prompts/` or `templates/` remains complete and operational, and embedded Phase 3 defaults are updated in the sprint that implements the corresponding stage.

It should not know:

```text
Which dimensions exist
How study synthesis works
How report summaries are generated
How code extraction parses citations
What review or smoke content means
```

Boundary:

```text
workspace = where am I and where are things?
study = what study work happens here?
```

### `codeextract`

`codeextract` owns code-reference extraction as a distinct product capability.

It owns:

```text
Parsing report citations
Resolving referenced source files
Extracting source snippets
Producing extraction output
Validating extraction requests
CLI commands for extraction workflows
```

It may consume report paths or metadata produced by `study`, but it should not become a generic `reports` package unless the behavior is genuinely shared by multiple modules.

## Dependency Rules

Use these rules to keep module boundaries clear:

```text
Product modules may depend on platform packages.
Product modules may depend on workspace when they need workspace paths.
Product modules should not depend on other product modules unless there is a clear product relationship.
Platform packages must not depend on product modules.
Shared helpers must not become a dumping ground for product behavior.
Runtime must not import study, project, or sprint.
```

Expected dependency direction for the current build:

```text
cmd/ultraplan -> internal/app
internal/app -> product modules + platform modules
study -> workspace
study -> platform/runtime
study -> platform/config/logging/filesystem as needed
project -> workspace
project -> platform/config/logging/filesystem as needed
sprint -> project
sprint -> workspace
sprint -> platform/runtime
sprint -> platform/config/logging/filesystem as needed
codeextract -> workspace
codeextract -> platform/filesystem/logging as needed
workspace -> platform/filesystem as needed
platform/* -> no product modules
```

## Encapsulation in Practice

A module should expose a small use-case-oriented surface and hide internal helpers.

Example shape:

```go
type Service struct {
    store   Store
    runtime Runtime
    clock   Clock
}

func (s *Service) InitStudy(ctx context.Context, req InitStudyRequest) error
func (s *Service) RunDimension(ctx context.Context, req RunDimensionRequest) error
func (s *Service) RunAll(ctx context.Context, req RunAllRequest) error
func (s *Service) Synthesize(ctx context.Context, req SynthesizeRequest) error
```

Internally, the module can call helpers such as:

```go
validateStudy(...)
buildAnalysisPrompt(...)
resolveSources(...)
scheduleRuns(...)
parseReports(...)
writeRunState(...)
```

Callers should not need to know those helpers exist.

## Interface Guidance

Interfaces should appear at external or volatile boundaries, not everywhere by default.

Good interface boundaries:

```text
Runtime execution
Filesystem persistence where tests need fakes
Clock/time
External process execution
```

Avoid introducing interfaces for every internal helper. If a function is purely internal to a module and not volatile, keep it concrete.

## Contract Interpretation For Go Sprints

The shared contracts are production standards. Apply them through this Go module architecture rather than literal Python-style package shapes such as `use_cases/ports.py` or `bootstrap.py`.

Current CLI foundation and study-discovery sprints may use concrete local filesystem collaborators when:

```text
The side effect is local and narrow.
The package boundary is explicit.
The behavior is tested through the public CLI or module surface.
The design does not block later introduction of runtime, persistence, or process-execution ports.
```

Do not reject current-sprint code only because it lacks a port or registrar. Reject it when a concrete dependency crosses a volatile boundary, hides side effects, makes tests require private mutation, or couples product modules in a way this document does not allow.

The following become mandatory when their capability enters scope:

```text
Runtime/provider execution: context propagation, cancellation, runtime ports, correlation IDs, retry ownership, bounded provider calls, and cost metadata.
Batch/run-loop execution: bounded concurrency, durable task state, diagnostics, terminal failure state, and resumability.
Stable public JSON/release: canonical structured error payloads, stable machine-readable error codes, scenario tests, documented compatibility, and migration/rejection behavior for durable formats.
```

`study -> workspace` is an allowed dependency for workspace path resolution and safety helpers. It is not a cross-module violation unless `study` starts depending on workspace-owned behavior that knows study semantics or reaches around the exported workspace package API.

## Reuse Boundary For Phase 2

Reuse infrastructure, not study semantics.

Good candidates for reuse:

```text
workspace discovery and path safety
config loading, precedence, and redaction
generic runtime prompt execution
command exit codes and output discipline
atomic file and JSON writes
durable task-state helpers when they remain product-neutral
small result/diagnostic structs without product semantics
```

Do not extract these from `study` for Phase 2:

```text
study Service
Source and Dimension models
study prompt builders
report validation
rating parsing
summary generation
run-loop scheduling
task state machines
```

If two modules need the same mechanical filesystem behavior, extract the mechanical helper. If two modules merely have similar product workflows, keep the behavior in the owning modules until repeated concrete implementations prove a shared abstraction is stable. Execute task semantics belong to `internal/sprint`; do not move them into a global scheduler or workflow package.

## Final Principle

```text
Platform owns generic capabilities.
Modules own product behavior.
Logic stays near the state it transforms.
Interfaces appear only at external or volatile boundaries.
```
