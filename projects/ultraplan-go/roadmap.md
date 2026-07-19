# UltraPlan Go Roadmap

> Project: `ultraplan-go`  
> Scope: production-grade Go CLI for UltraPlan study workflows, governed sprint planning and execution, post-execute review and deep smoke, and a local terminal UI over the same workflows.
> Product Phase 1 completed the study-side release scope. Product Phase 2 added governed project and sprint planning through `plan`, then controlled sprint implementation execution through `execute`. Sprints 24 and 25 delivered the TUI foundation and guarded operational controls as an enabling track. Product Phase 3 adds automated sprint review followed by deep smoke, with the full Phase 3 workflow exposed through both CLI and TUI. Hosted SaaS, browser UI, multi-user collaboration, general-purpose issue tracking, and automatic Git mutation remain deferred.

## Scope Principle

Phase 1 of the implementation roadmap is **study-side only**:

- workspace/config/health
- study initialization
- source and dimension discovery
- prompt composition
- runtime execution through `agentwrap` + OpenCode
- per-source analysis
- final synthesis
- report validation
- summary generation
- code-reference extraction
- resumable run-loop
- status, logs, diagnostics, cancellation, retry/fallback

Phase 2 introduces the planning side of UltraPlan:

```text
study -> select -> distill -> reason -> plan -> execute
```

Phase 2 now includes controlled implementation execution from validated sprint plans. The following are outside Phase 2; review and smoke enter scope in Product Phase 3 while the other items remain deferred:

- smoke investigation runs
- conformance review automation
- issue tracking
- automatic Git mutation
- hosted SaaS
- browser UI
- multi-user collaboration

The TUI enabling track introduced a local terminal UI. It uses shared application use cases and product services instead of scraping CLI text or invoking `ultraplan` as a subprocess for normal UltraPlan operations.

Sprints 24 and 25 delivered that TUI foundation before the next product phase. Product Phase 3 is the post-execute verification phase:

```text
study -> select -> distill -> reason -> plan -> execute -> review -> smoke
```

Review runs first because it is cheaper and more deterministic than live smoke. A blocking review verdict stops the default flow before smoke. Deep smoke then proves runtime-facing claims through the external harness cataloged by the project index. The sprint root keeps only the human-readable `review.md` and `smoke.md`; detailed smoke runs and issue evidence remain in the external smoke harness.

---

# Current Implementation Roadmap

## Implementation Wave 0 — Governance and Scope Alignment

### Sprint 0: Roadmap, Scope, and Template Alignment

**Goal:** remove ambiguity before implementation starts.

**Deliverables:**

- Confirm study-side-only release scope.
- Clean up TRD references that still imply current target/sprint implementation.
- Standardize sprint artifact names:

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

**Acceptance:**

- Roadmap is accepted.
- TRD target/sprint references are removed or clearly marked deferred.
- Embedded prompt/template defaults and generated artifact names match the canonical chain; workspace copies remain optional overrides.
- Contract reviews use the phase gates below rather than applying production-only requirements to skeleton sprints.

## Contract Phase Gates

The system contracts remain authoritative, but each sprint must apply only the portions selected by its sprint index and made relevant by implemented behavior.

### Skeleton / Local CLI Gate

Applies to CLI shell, workspace/config/health skeletons, and study discovery/listing.

Required now:

- explicit package ownership and allowed dependency direction from `docs/ARCHITECTURE.md`
- thin CLI adapters with deterministic stdout/stderr and stable exit statuses
- no swallowed failures
- preserved error cause chains when translating errors
- local filesystem path safety
- config precedence, validation, and redaction where config behavior exists
- deterministic unit/command tests for public behavior and negative paths

Allowed now:

- concrete local filesystem collaborators when the side effect is narrow and testable
- app-level command wiring without a full registrar/container
- numeric CLI exit statuses as the public process contract, paired with internal stable error codes where errors are classified

Not required until later gates:

- full canonical structured error payloads on every CLI failure
- request/trace correlation IDs
- diagnostics and alert surfaces
- runtime ports for provider/process execution
- e2e run-loop tests, migration tests, bounded worker pools, and provider-cost metadata

### Runtime / Provider Gate

Starts when prompt execution, `agentwrap`, OpenCode, retries, or external process supervision enter scope.

Required then:

- `context.Context` through runtime-capable services
- runtime and process-execution ports
- bounded retries with visible ownership and exhaustion state
- correlation metadata in runtime logs/events
- provider/model cost drivers in operational signals
- fake-runtime tests plus gated real-provider smoke coverage

### Batch / Durable Workflow Gate

Starts when `run-all`, `run-loop`, cancellation, resumability, status, or durable state enters scope.

Required then:

- bounded concurrency and cancellation propagation
- durable task state with explicit terminal success/failure/cancelled states
- atomic writes and compatibility or migration behavior for durable formats
- diagnostics for current/recent workflow state
- scenario tests for critical user journeys

### Stable Release Gate

Applies before the first production-ready study-side release.

Required then:

- stable JSON output schemas for documented automation surfaces
- canonical structured error payloads or documented safe CLI subsets
- stable machine-readable error codes for all public operational failures
- release documentation for config compatibility and recovery workflows
- smoke release checklist including offline tests, build, and gated OpenCode smoke path

---

## Study Implementation Wave 1 — CLI Foundation

### Sprint 1: Go Module, CLI Shell, and App Composition

**Goal:** establish a buildable Go CLI with module-driven structure.

**Build:**

- `cmd/ultraplan/main.go`
- `internal/app`
- initial package layout:

```text
internal/platform/config
internal/platform/logging
internal/platform/filesystem
internal/platform/runtime
internal/workspace
internal/study
internal/codeextract
```

- `ultraplan --help`
- `ultraplan version`

**Evidence:**

```bash
go test ./...
go build ./cmd/ultraplan
ultraplan --help
ultraplan version
```

### Sprint 2: Workspace, Config, Logging, and Health Skeleton

**Goal:** make UltraPlan able to find, initialize, inspect, and validate a workspace.

**Build:**

- workspace discovery
- `init-workspace`
- `config show`
- basic `health`
- config precedence skeleton
- secret redaction
- text/JSON output foundation

**Commands:**

```bash
ultraplan init-workspace
ultraplan config show
ultraplan health
```

---

## Study Implementation Wave 2 — Study Model and Initialization

### Sprint 3: Study Domain, Listing, and Resolution

**Goal:** model studies, sources, dimensions, and deterministic listing.

**Build:**

- `Study`
- `Source`
- `Dimension`
- source/dimension lookup
- study discovery
- `study list`
- `study <study> list`

**Commands:**

```bash
ultraplan study list
ultraplan study <study> list
```

### Sprint 4: Study Initialization From YAML

**Goal:** generate a study structure from `study-init.yml`.

**Build:**

- YAML parser
- dimension file generation
- study README generation
- normalized `study-init.yml`
- dry-run
- force behavior
- shallow clone support, if included in this sprint

**Command:**

```bash
ultraplan study init <study-init.yml>
ultraplan study init <study-init.yml> --dry-run
```

### Sprint 5: Markdown Document Sources and Applicability

**Goal:** support top-level Markdown document sources with optional `applicable_dimensions` frontmatter.

**Build:**

- Markdown source discovery
- frontmatter parsing
- frontmatter stripping
- dimension normalization
- applicability filtering
- skipped/inapplicable task behavior

**Core helper:**

```go
GetApplicableSources(sources []Source, dimension Dimension) []Source
```

**Acceptance:**

- Directory sources are always applicable.
- Markdown sources without filters apply to all dimensions.
- Markdown sources with filters only apply to matching normalized dimensions.
- Inapplicable pairs are skipped, not failed.

---

## Study Implementation Wave 3 — Reports, Prompt Composition, and Run State

### Sprint 6: Report Validation and Rating Parsing

**Goal:** reports become validated product artifacts, not just runtime output files.

**Build:**

- per-source report validation
- final report validation
- rating parser
- validation diagnostics
- Markdown-source validation rules that do not require code citations by default

### Sprint 7: Prompt Composition

**Goal:** deterministic prompts for directory analysis, Markdown analysis, and synthesis.

**Build:**

- base prompt loading
- report template loading
- directory source prompt builder
- Markdown document prompt builder
- synthesis prompt builder
- dry-run prompt preview

**Acceptance:**

- Directory prompts require source isolation and file-line citations.
- Markdown prompts embed stripped document content and forbid external code/filesystem exploration.
- Synthesis prompts include selected per-source report manifest.

**Command:**

```bash
ultraplan study <study> prompt analysis <dimension-ref> <source-ref>
ultraplan study <study> prompt synthesis <dimension-ref>
```

### Sprint 8: Run State Persistence and Status

**Goal:** create deterministic durable run-state primitives before executing runtime work.

**Build:**

- versioned `run-state.json`
- deterministic analysis and synthesis task construction
- source/dimension applicability in task construction
- atomic state writes
- strict state loading and validation
- completed-output revalidation on resume
- runtime-free status summaries

**Command:**

```bash
ultraplan study <study> status
```

---

## Study Implementation Wave 4 — Runtime and Prompt Execution

### Sprint 9: Agentwrap/OpenCode Runtime Integration

**Goal:** integrate runtime execution correctly without reimplementing OpenCode supervision.

**Build:**

- `agentwrap` runtime wiring
- `agentwrap/opencode` adapter construction
- runtime health checks
- `ObservingRuntime`
- `ValidatingRuntime`
- `PolicyRunner`
- permission policy mapping
- event/log mapping

**Architecture constraint:**

```text
study -> platform/runtime -> agentwrap/opencode
platform/runtime must not know study semantics
```

### Sprint 10: Single Analysis and Synthesis

**Goal:** run one source/dimension pair and synthesize one dimension.

**Build:**

```bash
ultraplan study <study> run <dimension-ref> <source-ref>
ultraplan study <study> synthesize <dimension-ref>
```

**Acceptance:**

- Runtime success alone is insufficient.
- Expected output must exist.
- Output validation must pass.
- Inapplicable Markdown pairs skip cleanly.

---

## Study Implementation Wave 5 — Batch Execution and Resumability

### Sprint 11: `run-all` Batch Execution

**Goal:** run all applicable analysis tasks with bounded parallelism.

**Build:**

```bash
ultraplan study <study> run-all
```

**Features:**

- filter by source and dimension
- bounded worker pool
- synthesis after applicable reports pass
- summary generation after completion

### Sprint 12: Durable `run-loop`, Retry, and Cancellation

**Goal:** production-grade resumable orchestration.

**Build:**

```bash
ultraplan study <study> run-loop
ultraplan study <study> status
```

**Features:**

- atomic `run-state.json`
- task state machine
- stale running task recovery
- retry/fallback metadata from `agentwrap`
- cancellation handling
- per-study lock file

### Sprint 13: Summary Generation and Code Reference Extraction

**Goal:** make completed study outputs easy to review, compare, and audit.

**Build:**

- deterministic `summary.csv`
- source/dimension score matrix
- missing vs inapplicable distinction
- total score sorting
- missing/ambiguous rating warnings
- source table parser
- citation parser
- reference resolver
- line/range/list extraction
- unresolved reference reporting
- text and JSON output

**Commands:**

```bash
ultraplan study <study> summary
ultraplan code <report>...
```

---

## Study Implementation Wave 6 — Hardening and Release

### Sprint 14: Validation Command, Diagnostics, and JSON Stability

**Goal:** make the CLI inspectable and automatable.

**Build:**

```bash
ultraplan study <study> validate
ultraplan health --json
ultraplan study <study> status --json
```

**Features:**

- stable JSON shapes
- actionable validation failures
- redacted diagnostics
- run metadata summaries

### Sprint 15: Docs, Packaging, and Smoke Release

**Goal:** first production-ready study-side release.

**Build:**

- user docs
- recovery docs
- config docs
- OpenCode smoke instructions
- Linux/macOS builds
- checksums

**Release gate:**

```bash
go test ./...
go build ./cmd/ultraplan
```

plus gated OpenCode smoke tests when environment is available.

---

# Planning Phase 2 — Governed Project, Sprint Planning, and Execute

Planning Phase 2 implements the planning and execution side of UltraPlan in Go. It ports the proven `cli` planning flow into production-grade modules, then adds controlled implementation execution from validated sprint plans.

The phase workflow is:

```text
study -> select -> distill -> reason -> plan -> execute
```

## Planning And Execute Scope Principle

Planning behavior belongs to project and sprint modules, not to the study module:

- `study` produces validated research artifacts.
- `project` catalogs governance documents, contracts, evidence, reasoning templates, and review protocols.
- `sprint` creates and validates planning artifacts selected from the project catalog, then executes validated plan tasks through the runtime boundary.
- `platform/runtime` remains generic prompt execution infrastructure.

Reusable study-side infrastructure is allowed when it is genuinely generic:

- workspace discovery and path safety
- config loading, precedence, and redaction
- generic runtime prompt execution
- command exit codes and output discipline
- atomic file/JSON writes when used by more than one module

Study semantics are not reusable abstractions for planning. Do not abstract study services, source/dimension models, report validation, rating parsing, summary generation, or run-loop scheduling into shared packages for Planning Phase 2.

## Planning And Execute Artifact Chain

Planning Phase 2 supports these sprint artifacts:

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
```

The following artifacts are added by Product Phase 3, not Planning Phase 2:

```text
review.md
smoke.md
```

Detailed smoke run data and issue evidence stay in the external smoke harness. General-purpose `issues.md` / `issues.json` artifacts remain deferred. Existing historical sprint `review.md` and `deep-smoke.md` files may remain until Phase 3 migration replaces them with the current generated `review.md` and `smoke.md`; Planning Phase 2 itself does not generate or automate either stage.

## Planning Phase 2 Contract Gate

Required for planning-stage sprints:

- project and sprint packages own planning behavior and validation
- sprint-index selections are validated as a subset of `project-index.md`
- technical handbooks distill only selected evidence
- reasoning artifacts decide from selected contracts, evidence, and handbook content
- plans trace to `reasoning.md`
- execute runs trace to validated `plan.md` tasks
- generated artifacts are editable Markdown
- flow state is versioned JSON with atomic writes and strict loading
- execute run state is versioned JSON with atomic writes, resumability, and explicit task terminal states
- prompt previews and dry runs are available before runtime execution
- runtime execution, when used for artifact generation, validates expected output before marking a stage complete
- runtime execution, when used for implementation tasks, records task status, diagnostics, attempts, and evidence before marking a task complete
- text output is calm and scriptable; JSON output is stable where documented

Not required in Planning Phase 2:

- smoke investigation execution
- conformance review automation
- local issue tracking
- Git mutation
- cross-sprint implementation run loops

### Sprint 16: Project Domain and Project Index

**Goal:** model `projects/<project>` as a first-class planning root.

**Build:**

- `internal/project`
- project discovery and resolution
- `docs/` discovery
- `roadmap.md` and `project-index.md` validation
- catalog parsing for contracts, evidence reports, reasoning templates, and review protocols
- `ultraplan project list`
- `ultraplan project <project> status`

**Acceptance:**

- Project commands do not depend on study internals.
- Project index remains a catalog, not a sprint plan.
- Catalog references resolve or fail with actionable diagnostics.

### Sprint 17: Sprint Artifact Domain and Flow State

**Goal:** model planning-stage sprint artifacts and durable stage state.

**Build:**

- `internal/sprint`
- sprint discovery and resolution
- stage model through `execute`:

```text
requirements
sprint-index
technical-handbook
area-reasoning
reasoning
plan
execute
```

- versioned `flow-state.json`
- atomic flow-state writes
- strict flow-state loading
- `ultraplan sprint <project> <sprint> status`

**Acceptance:**

- Existing artifacts can be inspected and reflected in flow state.
- Missing, ready, complete, failed, and skipped states are explicit.
- No smoke, review, or issue stages are modeled.

### Sprint 18: Select Stage

**Goal:** create and validate `sprint-index.md` as the authoritative context selection.

**Build:**

- `ultraplan sprint <project> <sprint> validate sprint-index`
- `ultraplan sprint <project> <sprint> prompt sprint-index`
- `ultraplan sprint <project> <sprint> flow --to sprint-index`
- subset validation against `project-index.md`
- selected contracts, evidence reports, reasoning templates, and review protocols

**Acceptance:**

- `sprint-index.md` cannot reference contracts, evidence, templates, or protocols absent from `project-index.md`.
- Excluded context is explicit.
- The stage supports dry run and prompt preview.

### Sprint 19: Distill Stage

**Goal:** produce `technical-handbook.md` from selected evidence only.

**Build:**

- technical handbook prompt rendering
- selected evidence loading
- handbook validation
- flow execution through `technical-handbook`

**Acceptance:**

- Handbook content traces to selected studies/reports.
- It captures relevant patterns, trade-offs, anti-patterns, open questions, and evidence pointers.
- It does not make implementation decisions.

### Sprint 20: Reason Stage

**Goal:** produce optional area reasoning and final `reasoning.md`.

**Build:**

- selected reasoning-template detection
- `reasoning/*.md` generation and validation when selected
- `reasoning.md` generation and validation
- flow execution through `reasoning`

**Acceptance:**

- Area reasoning is skipped only when no templates are selected.
- Final reasoning includes decisions, expected evidence, assumptions, and risks.
- Reasoning does not invent unselected contracts or evidence.

### Sprint 21: Plan Stage

**Goal:** produce `plan.md` from validated reasoning.

**Build:**

- plan prompt rendering
- plan validation
- task/evidence checklist validation
- flow execution through `plan`

**Acceptance:**

- `plan.md` cites `reasoning.md`.
- Tasks map to decisions and acceptance evidence.
- The plan stage stops after plan generation and validation.
- Implementation execution is invoked only by the later `execute` stage.

### Sprint 22: Planning Documentation and Release Gate

**Goal:** document and verify the planning-side release through `plan.md`.

**Build:**

- CLI reference for `project` and planning-stage `sprint` commands
- recovery docs for failed planning stages
- migration notes from `cli`
- offline fixture tests
- gated real-runtime planning smoke, if environment is available

**Release gate:**

```bash
go test ./...
go build ./cmd/ultraplan
```

plus gated planning-runtime smoke when environment is available.

**Status:** completed in `sprints/22-documentaiton-and-release/plan.md`.

### Sprint 23: Execute Stage

**Goal:** execute validated `plan.md` implementation tasks through the generic runtime boundary with durable task state, safe target-repository boundaries, resumability, and clear diagnostics.

**Build:**

- global and per-stage runtime model selection for sprint flow and execute
- execute prompt rendering
- executable task extraction from `plan.md`
- `.run-state.json` task state
- `execute.md` execution summary
- flow execution through `execute`
- sprint status updates that show execute progress

**Acceptance:**

- `execute` requires valid prerequisites through `plan`.
- Runtime model selection supports a global/default model and stage-specific overrides for planning stages and execute; stage-specific values win over the global/default value.
- Tasks trace to validated plan entries and preserve deterministic task IDs.
- Runtime-backed implementation uses only the generic platform runtime boundary.
- Target repository paths are explicit and workspace-safe.
- Runtime success alone is insufficient; task completion requires expected evidence or diagnostics.
- `.run-state.json` is versioned, atomic, resumable, and records pending, running, complete, failed, and cancelled task states.
- No smoke, automated review, issue tracking, Git add/commit/push, or cross-sprint scheduler behavior is invoked.

---

# Delivered TUI Enabling Track — Local Terminal UI

The post-execute TUI phase adds a TUI without changing the workspace model or making the CLI secondary. The TUI is a local terminal surface over existing product services and durable artifacts.

Before TUI feature work grows, reorganize app-layer glue so both CLI and TUI can call typed use cases:

```text
CLI command -> parse flags -> app use case -> text/JSON renderer
TUI action  -> app use case -> terminal model/view update
```

Do not integrate the TUI by shelling out to `ultraplan` or parsing CLI stdout. Shared use cases should cover workspace discovery, config loading, runtime setup, project status, sprint status/validation/flow, study listing/status/validation/run-loop, code extraction, and execute status/actions as needed.

External terminal UI dependencies belong at the first sprint that exposes `ultraplan tui` as an actual interactive terminal experience. Because Sprint 24 ships that command as a user-facing TUI, it introduces Bubble Tea inside `internal/tui` for the full-screen event loop and Glamour for Markdown preview rendering while keeping terminal-library and renderer types out of `internal/app` and product packages. Later TUI sprints may add richer widgets, split-pane resizing, async refresh commands, loading overlays, and guarded workflow dialogs on top of those contained dependencies.

## TUI Contract Gate

Required for TUI sprints:

- CLI remains the documented automation surface.
- TUI startup uses the same workspace discovery and config validation as CLI commands.
- Mutating and runtime-backed actions show the operation scope before execution.
- TUI cancellation propagates through `context.Context`.
- Durable workspace state remains the source of truth.
- TUI tests use deterministic model/update tests and fake runtimes by default.
- Terminal rendering handles narrow terminals without hiding critical status or errors.

Not required for TUI sprints:

- hosted service
- browser UI
- multi-user collaboration
- issue tracking
- Git mutation
- smoke/review automation

### Sprint 24: TUI Foundation and Read-Only Dashboard

**Goal:** add a read-only TUI that helps users inspect the workspace without changing artifacts or invoking runtimes.

**Build:**

- `ultraplan tui`
- `internal/tui` package, deterministic model/update/render foundation, Bubble Tea event loop, and contained Markdown preview rendering
- app-layer use-case extraction for read-only operations
- workspace startup state
- top-level Projects and Studies tabs with focus movement between tab bar and content
- project list, project status summary, and project detail navigation
- project-scoped sprint list plus sprint artifact navigation
- study list, study detail navigation, dimensions, sources, and run-state preview
- selection-follow viewport behavior for lists and preview-scroll behavior for artifact previews
- validation finding panes for project, study, and sprint artifacts where existing validators support it
- name-only navigation with bounded Markdown/JSON previews for key artifacts

**Dependency note:** Sprint 24 introduces Bubble Tea because `ultraplan tui` is a real interactive terminal command, not a text dump. It also introduces Glamour for terminal Markdown rendering. Both dependencies remain contained inside `internal/tui`; `internal/app` and product packages continue to exchange plain Go result types.

**Acceptance:**

- TUI does not call CLI command handlers or parse stdout.
- TUI can start from inside a workspace or with `--workspace`.
- Dashboard works without runtime credentials.
- Projects and Studies are top-level tabs; sprints are navigated through the selected project rather than a top-level tab.
- Artifact navigation shows names rather than paths, with paths retained as internal preview metadata.
- Markdown previews render headings, tables, inline code, and lists as terminal Markdown rather than raw Markdown text.
- Read-only actions do not mutate flow state except where existing status operations intentionally refresh deterministic status files; any such mutation is documented in the TUI action label.
- Unit tests cover tab focus, route navigation, selection-follow scrolling, preview scrolling, Markdown rendering, model updates, and fake app use cases.

### Sprint 25: Operational TUI Controls

**Goal:** allow guarded local operation of validation, dry-run, prompt preview, planning flow, execute status, and study run-loop monitoring from the TUI.

**Build:**

- guarded action dialogs for mutating or runtime-backed operations
- project/study/sprint validation actions
- sprint flow dry-run and prompt preview actions
- sprint planning flow actions through selected stages
- execute status and execute action entry points
- study run-loop start/resume controls with filters and parallelism
- live progress/event panes backed by existing progress callbacks
- cancellation handling and post-cancel state refresh
- error detail panes with classified error codes and guidance

**Acceptance:**

- Mutating/runtime actions require explicit confirmation after showing affected workspace paths and stage/task scope.
- Running workflows update from typed progress/events, not terminal text scraping.
- Cancellation leaves durable state inspectable from both TUI and CLI.
- Normal tests use fake runtimes; gated smoke may cover real runtime operation when available.
- CLI command behavior remains unchanged.
- External terminal-library types remain contained inside `internal/tui`; app and product result types stay plain Go data.

---

# Product Phase 3 — Review and Deep Smoke

Product Phase 3 completes the governed sprint workflow after `execute`:

```text
requirements
  -> sprint-index
  -> technical-handbook
  -> area-reasoning
  -> reasoning
  -> plan
  -> execute
  -> review
  -> smoke
```

The phase deliberately keeps its sprint artifacts simple:

```text
review.md
smoke.md
```

`review.md` is the current automated conformance review and replaces the older manually produced review when the review stage is run. `smoke.md` is the current sprint smoke summary. Raw smoke JSON, command output, per-test artifacts, and open/resolved issue files remain under the external smoke harness referenced by `project-index.md`; `smoke.md` links the relevant run IDs and evidence paths.

## Phase 3 Contract Gate

Required for Phase 3:

- `internal/sprint` owns review and smoke stage semantics, validation, prompt construction, artifact paths, flow ordering, and verdict rules.
- `internal/platform/runtime` remains generic and executes structured review requests through agentwrap without learning sprint or contract semantics.
- External smoke execution uses an explicit executable plus argument list, bounded environment forwarding, `context.Context`, timeout, and cancellation; it must not evaluate a shell command assembled from Markdown.
- Review resolves selected contracts and review protocols from `project-index.md`; it does not hardcode a fixed contract map.
- Reviewers are read-only and return structured results. UltraPlan synthesizes the final verdict and atomically writes `review.md`.
- Review runs before smoke. Blocking or high-severity applicable review findings stop the default flow before smoke.
- Smoke selects the narrowest harness suite that covers the sprint, records why it is sufficient, and atomically writes `smoke.md` with links to the external evidence.
- Missing required runtime, credentials, harness coverage, or environment produces a blocked result, never a false pass.
- Runtime-free or otherwise irrelevant smoke scope is reported as `not_applicable`, not as a passing run.
- A code, contract, handbook, reasoning, plan, or execute-state change makes prior review and smoke results stale until rerun.
- Review and smoke do not modify product source, product tests, governed planning artifacts, or Git state. The smoke harness may update its own tests, runs, and issue evidence only through a separate harness-maintenance action.
- CLI and TUI call the same typed app use cases. Every Phase 3 sprint updates the TUI for the functionality introduced by that sprint.
- Phase 3 prompt and output-template defaults ship embedded in UltraPlan. Workspace `prompts/` and `templates/` files remain optional intentional overrides and are never stage prerequisites.
- Default tests use fake runtimes and a fake smoke harness. Real OpenCode and live smoke remain gated.

Still deferred:

- general-purpose issue tracking, assignment, scheduling, and remote issue synchronization
- automatic product fixes during review or smoke
- Git add, commit, push, branch, merge, or reset
- hosted review services, browser UI, and multi-user collaboration
- cross-project or cross-sprint verification scheduling

### Sprint 26: Automated Sprint Review

**Goal:** replace the manual sprint review with a product-owned, evidence-grounded review stage that writes the current `review.md` and is fully operable from the TUI.

**Build:**

- Phase 3 domain additions in `internal/sprint` for review scope, reviewer tasks, findings, verdicts, validation, and flow integration
- dynamic selected-contract and selected-review-protocol resolution through `project-index.md`
- a frozen review manifest held in run state while the command executes, including governed input hashes, target implementation identity, changed-path scope, plan tasks, and execute evidence
- one independent, structured agentwrap review request per selected contract plus one technical-handbook review request
- bounded review concurrency, cancellation, retry, and in-process resume of completed reviewer results
- read-only reviewer permissions and validated structured reviewer output
- deterministic checks for decision conformance, plan-task execution evidence, verification-command results, citation containment, and missing reviewer coverage
- deterministic severity and verdict synthesis; an LLM may summarize findings but cannot choose the machine verdict
- updated embedded review prompt and `review.md` output-template defaults, available for optional customization through `ultraplan defaults install`
- atomic replacement of sprint-root `review.md`
- `ultraplan sprint <project> <sprint> review`
- review support in `status`, `validate review`, `prompt review`, and `flow --to review`
- text and stable JSON output
- TUI review readiness/status, dry-run and prompt preview, confirmation, live reviewer progress, cancellation, result display, finding navigation, and `review.md` preview

**Acceptance:**

- Review requires valid prerequisites through `execute`, unless an explicit review-only diagnostic mode is later approved.
- Every selected contract resolves from the project catalog; missing, duplicate, unknown, or escaping paths fail preflight.
- Review does not assume the runtime can spawn subagents; UltraPlan owns bounded fan-out over independent runtime calls.
- A failed reviewer task is recorded and cannot be silently converted into a passing review.
- Applicable blocker/high findings produce `fail`; only medium/low/info findings may produce `pass_with_findings`.
- All citations are workspace/target-contained and checked before the final artifact is accepted.
- Runtime success is insufficient: `review.md` must exist, contain the required sections, and pass review validation before the stage is complete.
- Re-running review atomically replaces the old manual or generated `review.md`; unrelated sprint artifacts are unchanged.
- TUI and CLI expose the same options, progress, cancellation, and final verdict through shared typed use cases.
- Offline fake-runtime tests, `go test ./...`, `go test -race ./...`, and `go build ./cmd/ultraplan` pass.

### Sprint 27: Deep Smoke Harness Integration

**Goal:** run the narrowest sufficient real-system smoke after review, keep raw evidence in the external harness, write the current sprint `smoke.md`, and expose the full operation through the TUI.

**Build:**

- a versioned smoke-harness manifest in `ultraplan-go-smoke` describing its entrypoint, supported discovery/run commands, evidence directories, and protocol version
- project-index discovery of the harness and its manifest
- structured harness discovery for levels, suites, tests, prerequisites, expected duration/cost class, and sprint mappings
- safe external-process execution without shell interpolation
- review-verdict gate, smoke preflight, narrowest-sufficient suite selection, explicit level/suite/test overrides, and dry-run
- timeout, cancellation, descendant-process cleanup, and post-cancel evidence refresh
- structured import of run ID, command arguments, result counts, model/runtime, duration, evidence paths, hashes, and relevant open/resolved issue IDs
- an embedded `smoke.md` output-template default, available for optional customization through `ultraplan defaults install`; no workspace template is required
- atomic replacement of sprint-root `smoke.md`
- `ultraplan sprint <project> <sprint> smoke`
- smoke support in `status`, `validate smoke`, and `flow --to smoke`
- text and stable JSON output
- TUI smoke readiness/status, scope selection, prerequisite display, cost/duration class, confirmation, live suite/test progress, cancellation, result display, issue summary, and `smoke.md` preview

**Acceptance:**

- Smoke runs after a passing or non-blocking review by default; an explicit force flag is required to investigate after a failed review.
- Harness commands and models are discovered or explicitly configured; no developer-specific absolute command is built into the protocol.
- The product never treats README prose as an executable command.
- Raw evidence remains in the smoke harness `runs/` and `issues/` directories. `smoke.md` contains stable run IDs and links rather than copying the full evidence set.
- A failing selected smoke path produces `fail`; unavailable required prerequisites produce `blocked`; irrelevant smoke produces `not_applicable`.
- Product source and governed sprint inputs remain unchanged.
- Cancellation terminates the owned process tree and leaves a clear, inspectable result.
- TUI and CLI expose the same scope, confirmation, progress, cancellation, and verdict through shared typed use cases.
- Fake-harness tests and gated real-harness tests cover success, failure, blocked, not-applicable, timeout, cancellation, malformed JSON, missing evidence, and path escape.

### Sprint 28: Integrated Review-to-Smoke Verification Flow

**Goal:** make `execute -> review -> smoke` a coherent resumable flow with stale-result detection, focused reruns, clear recovery, and a complete TUI workflow.

**Build:**

- flow ordering and prerequisite enforcement through smoke
- input fingerprints in flow state for review and smoke freshness
- combined sprint status showing review and smoke execution state, verdict, artifact, run ID, and required next action
- `ultraplan sprint <project> <sprint> verify` as a convenience command over the same review/smoke use cases
- focused review reruns and smoke test/suite reruns without bypassing final containing-suite evidence
- deterministic overall assessment rendered in status from the current `review.md`, `smoke.md`, flow state, and referenced harness issues; no third assessment artifact is added
- recovery for interruption, stale inputs, malformed artifacts, missing harness evidence, and externally edited review/smoke summaries
- TUI end-to-end verification action from execute status through review and smoke, with gate explanations, focused rerun controls, linked raw evidence, current overall assessment, and recovery guidance

**Acceptance:**

- `flow --to smoke` and `verify` always apply review before smoke unless the user explicitly selects a diagnostic override.
- A governed-input or implementation change marks prior review and smoke evidence stale.
- A passing narrow smoke test does not hide failure of the containing required suite.
- Relevant open harness issues prevent a clean smoke pass and are listed in `smoke.md` and status.
- The overall assessment is deterministic and cannot contradict the stage verdicts.
- Interrupted review or smoke can be rerun without corrupting flow state or the prior complete Markdown artifact.
- CLI, JSON, and TUI agree on current state, verdict, staleness, and next action.
- TUI supports the complete normal Phase 3 workflow without invoking CLI handlers or parsing terminal output.

### Sprint 29: Phase 3 Documentation, Hardening, and Release

**Goal:** stabilize review and smoke as supported CLI and TUI workflows, migrate away from manual artifacts, dogfood the phase, and pass the release gate.

**Build:**

- CLI reference, user guide, TUI guide, recovery guide, configuration reference, smoke-harness guide, migration guide, and release checklist updates
- stable JSON schema documentation for review, smoke, verify, and status
- documentation source-of-truth declaration and validation of the single authoritative project planning-doc set
- migration instructions that allow old manual `review.md` and `deep-smoke.md` files to be removed and replaced by current generated `review.md` and `smoke.md`
- security/redaction audit, path-containment audit, bounded-concurrency audit, race tests, cancellation tests, and durable-state compatibility fixtures
- dogfood review and smoke against representative planning, execute, and TUI sprints
- gated real-runtime test of the review stage and gated real smoke-harness test
- final TUI Phase 3 polish: review/smoke dashboard summaries, evidence links, responsive progress, error and recovery panes, narrow-terminal behavior, keyboard help, and documentation

**Release gate:**

```bash
go test ./...
go test -race ./...
go build ./cmd/ultraplan
```

plus:

- fake-runtime review suite passes
- fake smoke-harness suite passes
- deliberate contract violation produces a failing `review.md`
- deliberate runtime behavior failure produces a failing `smoke.md` linked to harness evidence
- required-environment absence produces blocked, not pass
- cancellation leaves CLI and TUI state recoverable
- authoritative project planning docs are internally consistent and current
- no relevant blocker/high review finding or open smoke issue remains
- the entire normal `execute -> review -> smoke` workflow is operable from both CLI and TUI

---

# Roadmap Review Gates

Before Sprint 1 starts:

- [ ] Roadmap accepted.
- [ ] Phase 1 study-side scope confirmed.
- [ ] TRD target/sprint contradictions resolved or replaced by Planning Phase 2 scope.
- [ ] Canonical sprint artifact names confirmed.

Before runtime integration starts:

- [ ] Agentwrap dependency and version strategy confirmed.
- [ ] OpenCode smoke environment documented.
- [ ] Runtime validation strategy reviewed.

Before first release:

- [ ] `go test ./...` passes.
- [ ] CLI builds as a single binary.
- [ ] Fake-runtime test suite passes.
- [ ] Gated OpenCode smoke path documented.
- [ ] User-facing docs explain normal and recovery workflows.

Before TUI foundation starts:

- [ ] App-layer read-only use cases are separated from CLI renderers.
- [ ] TUI dependency choice is reviewed and contained inside `internal/tui`.
- [ ] Terminal test strategy is documented.

Before operational TUI starts:

- [ ] Runtime-backed app use cases support progress callbacks and cancellation.
- [ ] Confirmation policy for mutating TUI actions is documented.
- [ ] Fake-runtime TUI test harness exists.

Before automated review starts:

- [x] Phase 3 PRD, TRD, architecture, roadmap, project-index, and protocol updates are accepted; embedded review defaults are updated in Sprint 26.
- [x] `review.md` is confirmed as the canonical current review artifact and may replace manual review output.
- [x] Review permissions, selected-contract resolution, verdict policy, and fake-runtime strategy are documented.

Before deep smoke starts:

- [ ] `ultraplan-go-smoke` exposes a versioned machine-readable harness contract.
- [x] `smoke.md` is confirmed as the canonical sprint summary while detailed evidence stays in the external harness.
- [x] Review-before-smoke gating, environment classification, timeout, cancellation, and mutation boundaries are documented.

Before Phase 3 release:

- [ ] Review and smoke JSON schemas are stable and documented.
- [ ] CLI/TUI parity tests cover every Phase 3 operation.
- [ ] Legacy manual review/deep-smoke migration is documented.
- [ ] Authoritative planning-doc validation passes.
- [ ] Phase 3 dogfood evidence exists for review, smoke, cancellation, failure, and recovery paths.
