# UltraPlan Go Roadmap

> Project: `ultraplan-go`  
> Scope: production-grade Go CLI for UltraPlan study workflows, governed sprint planning and execution, post-execute review and deep smoke, a local terminal UI, and a loopback-only Go-served browser UI over the same workflows.
> Product Phase 1 completed the study-side release scope. Product Phase 2 added governed project and sprint planning through `plan`, then controlled sprint implementation execution through `execute`. Sprints 24 and 25 delivered the TUI foundation and guarded operational controls as an enabling track. Product Phase 3 adds automated sprint review followed by deep smoke, with the full Phase 3 workflow exposed through both CLI and TUI. Product Phase 4 begins at Sprint 30 and adds a local Go HTTP server, Go-rendered browser UI, guarded HTTP operations, and SSE progress. Product Phase 5 begins at Sprint 33 and adds a durable `code-context` stage whose exact output is reused across downstream planning, execution, and verification prompts. Hosted SaaS, remote exposure, multi-user collaboration, general-purpose issue tracking, automatic Git mutation, content identity, provenance, QA expansion, retrieval, and automatic repair remain deferred beyond this roadmap chunk.

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

Phase 2 now includes controlled implementation execution from validated sprint plans. The following are outside Phase 2; review and smoke enter scope in Product Phase 3, the local browser UI enters Product Phase 4, and the other items remain deferred:

- smoke investigation runs
- conformance review automation
- issue tracking
- automatic Git mutation
- hosted SaaS
- browser UI until Product Phase 4
- multi-user collaboration

The TUI enabling track introduced a local terminal UI. It uses shared application use cases and product services instead of scraping CLI text or invoking `ultraplan` as a subprocess for normal UltraPlan operations.

Sprints 24 and 25 delivered that TUI foundation before the next product phase. Product Phase 3 is the post-execute verification phase:

```text
study -> select -> distill -> reason -> plan -> execute -> review -> smoke
```

Review runs first because it is cheaper and more deterministic than live smoke. A blocking review verdict stops the default flow before smoke. Deep smoke then proves runtime-facing claims through the external harness cataloged by the project index. The sprint root keeps only the human-readable `review.md` and `smoke.md`; detailed smoke runs and issue evidence remain in the external smoke harness.

Product Phase 4 adds no new product workflow stage:

```text
browser -> local HTTP/SSE adapter -> shared app use cases -> existing product modules
```

The server and browser UI are local interfaces over the existing filesystem-backed product. They do not add hosted service behavior, a remote worker protocol, accounts, tenancy, or a database-backed alternate source of truth.

Product Phase 5 adds the first new workflow through those shared interface boundaries:

```text
requirements
  -> code-context
  -> sprint-index
  -> technical-handbook
  -> area-reasoning
  -> reasoning
  -> plan
  -> execute
  -> review
  -> smoke
```

The new stage gathers implementation evidence once into `code-context.md` and reuses that exact foundation downstream. It does not introduce repository indexing, retrieval, cache ownership, automatic staleness, or a second machine-readable context format.

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

# Product Phase 4 — Local Go Server and Browser UI

Product Phase 4 begins with Sprint 30. It adds a simple local browser surface without changing the governed workflow, workspace layout, or source of truth.

```text
CLI command ─┐
TUI action  ─┼─> shared app use cases -> product modules -> workspace artifacts
HTTP action ─┘

app progress -> bounded server event hub -> SSE -> browser
```

The initial UI is rendered with Go `html/template` and uses embedded CSS plus minimal JavaScript for form submission, progressive refresh, cancellation, and `EventSource`. It does not introduce React, Vite, Node.js, a frontend build step, or a separate frontend process. A richer client framework remains a later option only if real interaction complexity earns it.

## Phase 4 Contract Gate

Required for Phase 4:

- `internal/web` owns only loopback HTTP lifecycle, HTML/templates/static assets, transport DTO mapping, browser security, confirmation tokens, and SSE subscriptions.
- HTTP handlers call typed `internal/app` use cases and do not import product modules, invoke CLI handlers, parse CLI output, or execute `ultraplan` as a subprocess.
- `cmd/ultraplan` explicitly constructs CLI/TUI/web interface dependencies; Phase 4 does not add another package-global mutable runner registration.
- Workspace artifacts and product-owned run state remain authoritative. The server operation hub is ephemeral, bounded, and recoverable by rereading durable state.
- Commands use explicit HTTP methods and structured responses. SSE is a one-way progress channel, not a command or durable-state protocol.
- Browser disconnect cancels only the subscription. Explicit cancellation propagates through the shared app operation context and product/runtime cleanup path.
- Mutating or runtime-backed actions require a short-lived server-issued confirmation bound to normalized scope and current governed-input fingerprints.
- Study and sprint mutation exclusion remains product-owned. HTTP middleware must not become the only concurrency guard.
- The server binds to loopback only, uses same-origin requests, validates Host and Origin, protects mutations from CSRF, bounds bodies/timeouts/streams, redacts secrets, contains paths, and safely renders untrusted Markdown.
- Normal tests use `httptest`, fake app use cases, fake runtimes, and fake harnesses. Gated tests cover real runtime/harness behavior only where needed.
- CLI and TUI behavior remain supported and unchanged except for the explicit composition refactor needed to add the new surface cleanly.

Still deferred:

- hosted SaaS or LAN/public binding
- multi-user authentication, authorization, accounts, teams, tenancy, or collaboration
- remote workers and remote workspace synchronization
- browser editing of arbitrary workspace files
- a database-backed alternate product state
- interactive terminal/session transport, bidirectional agent chat, or WebSockets
- general-purpose issue tracking and automatic Git mutation

### Sprint 30: Local Web Foundation and Read-Only Dashboard

**Goal:** add `ultraplan serve` and a loopback-only, read-only browser dashboard using Go HTTP/templates over existing app use cases.

**Build:**

- refactor TUI/interface runner construction from package-global registration to explicit composition in `cmd/ultraplan`
- `ultraplan serve` with explicit loopback listen address, optional browser opening, signal-aware graceful shutdown, and bounded HTTP timeouts
- `internal/web` package with `net/http` routes, middleware, `html/template` rendering, embedded CSS/minimal JavaScript, and structured error mapping
- browser dashboard for workspace, projects, project sprints, studies, validation summaries, and current run/flow state
- project, sprint, and study detail pages backed by typed app queries
- bounded, allowlisted Markdown/JSON artifact previews with hostile-content-safe rendering
- initial `/api/v1` read-only JSON endpoints and a server health endpoint
- same-origin, Host/Origin, path-containment, redaction, request-limit, and security-header foundation
- CLI help, local-web user documentation, architecture reasoning, and API design documentation

**Acceptance:**

- `ultraplan serve` starts from a workspace or with `--workspace`, listens only on loopback, and shuts down cleanly on context cancellation or signal.
- A browser can inspect the same workspace, project, sprint, study, validation, and artifact state exposed by shared app use cases.
- The web package does not call CLI handlers, parse stdout, import product modules directly, or persist web-specific product state.
- The server and UI run from the Go binary without Node.js, Vite, a separate asset server, or a database.
- Unknown `/api/` routes return structured JSON errors and never fall through to an HTML page.
- Artifact previews reject unsupported and escaping paths, remain bounded, and do not execute workspace HTML/scripts.
- `httptest`, template, security, route, API-shape, shutdown, `go test ./...`, `go test -race ./...`, and `go build ./cmd/ultraplan` checks pass.

### Sprint 31: Guarded Web Operations and SSE Progress

**Goal:** expose the existing guarded local operations through HTTP and stream truthful live progress to the browser without introducing another workflow engine.

**Build:**

- browser validation, prompt-preview, dry-run, flow, execute, review, smoke, verify, study run-loop, and cancellation actions as supported by current app use cases
- `POST /api/v1/operations/prepare` returning scope, paths, runtime/model information, mutation class, input fingerprint, expiry, and a bound confirmation token
- confirmed operation start, status/result, and explicit cancellation endpoints
- bounded ephemeral operation hub with operation IDs, cancellation functions, safe recent events, SSE subscribers, terminal results, and short retention
- SSE event IDs, typed event names, heartbeat comments, reconnect support, slow-subscriber handling, and durable-status refresh guidance
- browser progress, result, finding, failure, recovery, and cancellation views
- product-owned per-sprint mutation locking plus conflict diagnostics; existing study locking remains authoritative
- structured conflict, stale-confirmation, cancellation, validation, runtime, and internal error responses

**Acceptance:**

- Mutating or runtime-backed work cannot start without a valid current server-issued confirmation matching the normalized request.
- Commands use HTTP POST/DELETE; SSE carries progress only.
- A disconnected or slow browser cannot block, complete, fail, or silently cancel product work.
- SSE reconnect resumes from bounded recent events when available and otherwise directs the browser to current durable status.
- Explicit cancellation reaches the shared operation context, preserves cleanup metadata, and leaves durable state recoverable from CLI, TUI, and browser.
- Conflicting sprint or study mutations fail with actionable scope/lock information rather than running concurrently.
- Normal tests cover success, validation failure, runtime failure, cancellation, stale confirmation, reconnect, buffer rollover, slow subscriber, shutdown, redaction, and CLI/TUI/web agreement with fake dependencies.

## Post-Sprint-31 Delivery Chunk 1 — Observable Grounded Planning

This chunk finishes the browser as a supported product surface, then proves its extensibility by adding `code-context` as the first new workflow operated through that surface.

The boundary is deliberately limited to Sprints 32–34:

```text
Sprint 32: Web hardening and release
Sprint 33: Code-context vertical slice
Sprint 34: Shared-context integration and release
```

The integrated product roadmap and this workspace roadmap use different phase numbering. Their mapping for this chunk is:

| Integrated product roadmap | Workspace product roadmap |
|---|---|
| Phases 1–3: web foundation, operations, and hardening | Product Phase 4, Sprints 30–32 |
| Phase 4: code-context | Product Phase 5, Sprints 33–34 |
| Phase 5 onward: content identity, provenance, and QA | Later roadmap chunks |

**Immediate prerequisite:** Complete Sprint 31's planning chain before materializing Sprint 32. Its current `flow-state.json` records `technical-handbook` as failed and `area-reasoning`, `reasoning`, and `plan` as missing.

### Sprint 32: Local Web Hardening and Observable-Product Release

**Goal:** turn the Sprint 30–31 browser implementation into a supported local interface and prove that the application boundary is ready to accept new stages.

**Build:**

- stabilize `/api/v1` JSON responses and typed error envelopes, with documentation and compatibility fixtures
- complete local-web user, configuration, packaging, security, recovery, and troubleshooting documentation
- accessibility and keyboard-navigation pass for dashboard, details, confirmations, progress, findings, and errors
- cache policy for HTML/static/API responses; bounded polling/refetch behavior for changes made by CLI or TUI
- CSRF, Host/Origin, session, CSP/security-header, request-smuggling/body-limit, hostile-Markdown, path-containment, and redaction audit
- operation/SSE concurrency, leak, race, slow-client, reconnect, cancellation, reconciliation, graceful-shutdown, and recovery tests
- explicit shutdown coverage proving that server-owned active operations are cancelled, bounded cleanup produces a truthful durable outcome, browser disconnection does not cancel work, and restart reconciles interrupted or cleanup-uncertain work
- representative browser integration tests over temporary workspaces and fake runtime/harness dependencies
- gated real-runtime and real smoke-harness browser operation evidence
- release packaging and checks confirming templates/static assets are embedded in the single Go binary
- an interface capability test proving that stage status, artifacts, commands, progress, cancellation, and recovery are exposed through shared application abstractions rather than route-specific workflow logic

**Release gate:**

```bash
go test ./...
go test -race ./...
go build ./cmd/ultraplan
```

plus:

- `ultraplan serve` is loopback-only and needs no external frontend/runtime dependency
- browser, CLI, and TUI agree on durable state, readiness, verdicts, artifacts, and next actions
- guarded browser operations, progress, reconnect, cancellation, and recovery pass representative scenarios
- browser refresh or server restart never promotes ephemeral state over durable workspace truth
- one substantial study workflow and one substantial sprint workflow can be observed and operated through the browser
- no high-severity local-web security, path, redaction, concurrency, cancellation, shutdown, or accessibility finding remains
- adding another planning stage does not require a parallel web workflow implementation
- hosted/multi-user/remote-worker behavior has not entered the implementation accidentally

---

# Product Phase 5 — Grounded Planning Foundation

Product Phase 5 gathers implementation evidence once per sprint, preserves it as a durable artifact, and reuses the exact same foundation across planning, execution, and verification.

```text
requirements
  -> code-context
  -> sprint-index
  -> technical-handbook
  -> area-reasoning
  -> reasoning
  -> plan
  -> execute
  -> review
  -> smoke
```

The authoritative artifact is:

```text
projects/<project>/sprints/<sprint>/code-context.md
```

## Phase 5 Contract Gate

Required for Phase 5:

- `internal/sprint` owns `code-context` stage order, prerequisites, status, validation, artifact paths, flow transitions, and downstream prompt integration.
- The stage reuses the existing project implementation-repository/worktree resolution and may read the source repository while writing only the sprint's `code-context.md`.
- `code-context.md` remains the single authoritative context-pack format; no parallel JSON manifest is added.
- Existing workspaces and persisted sprint state that predate `code-context` remain usable through explicit compatibility handling.
- Downstream prompts use one shared sprint-context renderer. Exact `requirements.md` and `code-context.md` content appear before stage-specific instructions without rewriting or dynamic run data.
- Downstream agents may inspect any additional repository files needed to verify assumptions or complete their work.
- CLI, TUI, and browser use the same application operations and durable state. The browser gains no route-specific `code-context` workflow semantics.
- Runtime success alone is insufficient: the artifact must exist and pass structural validation before the stage completes.
- No repository index, RAG system, cache subsystem, cache key, provider-specific cache dependency, automatic staleness system, or amendment workflow enters this phase.

### Sprint 33: Code-Context Stage Vertical Slice

**Goal:** make `code-context` a fully operational stage capable of inspecting the implementation repository and producing a validated source-context pack.

**Build:**

- add `StageCodeContext` immediately after `StageRequirements` in every canonical ordered stage list
- add artifact-path, readiness, prerequisite, status, cumulative-flow, and flow-state support
- preserve compatibility with existing sprint state that predates the stage
- add the embedded `create-code-context` prompt and `code-context.md` template, including defaults-install registration
- add stage-specific runtime model and variant configuration with existing fallback behavior
- resolve the project repository/worktree through existing implementation-repository mechanisms
- permit repository reads while restricting stage output to `code-context.md`
- add prompt preview, structural validation, and runtime-backed application use cases
- validate required sections, selected excerpts, repository-relative paths, rationale, language-tagged fenced code, safe paths, and well-formed line ranges
- add CLI support for `prompt code-context`, `validate code-context`, `flow --to code-context`, help, status, and stable JSON projections
- surface readiness, progress, validation findings, artifact preview, explicit rerun, cancellation, and recovery through the existing web operation model
- add focused stage-order, runtime, state, compatibility, CLI, application, web, defaults, and validation tests

**Acceptance:**

- completed valid requirements make `code-context` ready, and a valid context pack makes `sprint-index` ready
- a fake runtime or gated real runtime can inspect the source repository and produce a valid `code-context.md`
- missing or invalid output fails truthfully, records actionable findings, and does not mark the stage complete
- the browser observes and controls the stage through shared application and operation abstractions without new route-specific workflow semantics
- existing workspaces and pre-stage flow state remain usable
- no repository index, RAG system, cache subsystem, JSON context manifest, or automatic staleness system is introduced

This sprint is useful independently of downstream prompt reuse: users can explicitly generate, inspect, validate, rerun, and recover the implementation evidence pack.

### Sprint 34: Shared Context Integration and Grounded-Planning Release

**Goal:** reuse the stored requirements and code-context pack unchanged across every downstream agent operation.

**Build:**

- add one shared sprint-context renderer
- compose downstream prompts in this order:

```text
stable shared instructions
sprint identity
exact requirements.md
exact code-context.md
other shared context
stage-specific instructions
```

- keep the requirements and code-context block byte-for-byte identical across compatible downstream calls
- inject the shared context into sprint index, technical handbook, area reasoning, final reasoning, plan, execute, Conformance Review, and smoke wherever those operations invoke an agent; retain the same integration point for later QA without adding QA in this chunk
- ensure `flow --to plan` runs `code-context` exactly once in the correct position
- explicitly allow every downstream agent to inspect additional repository files
- add prefix-stability fixtures proving that timestamps, run IDs, stage names, output paths, and other dynamic run data do not enter the shared prefix
- add and materialize the manually invokable `code-context` stage skill using the canonical CLI operation
- update README, CLI reference, user guide, architecture, recovery, planning-smoke, generated-workspace, skill, and local-web documentation
- dogfood the stage in a temporary representative workspace against a real implementation repository and a gated runtime
- verify prompt reuse, cancellation, rerun, atomic artifact replacement, and browser recovery end to end

**Release gate:**

- every relevant downstream agent prompt contains the exact stored requirements and context pack
- the common prefix remains byte-for-byte stable until stage-specific instructions begin
- rerunning `code-context` atomically replaces only `code-context.md`
- downstream agents remain free to inspect additional source
- CLI, TUI, and browser agree on stage readiness, state, artifacts, findings, and recovery
- the complete requirements-to-plan flow succeeds with `code-context` inserted exactly once
- fake-runtime coverage and gated real-runtime dogfood pass
- `go test ./...`, `go test -race ./...`, and `go build ./cmd/ultraplan` pass
- no caching claim or dependency is introduced; the prompt layout merely enables provider prefix caching where available

## Why This Chunk Stops After Sprint 34

Sprint 34 is a clean product milestone:

```text
observable browser
  -> first new stage added through shared boundaries
  -> implementation evidence gathered once
  -> evidence reused across the full sprint lifecycle
```

The resulting real `code-context.md` examples should inform the next chunk's content-identity, provenance, and QA work. That chunk should begin with an evidence-based artifact inventory and retrieval-question corpus rather than designing metadata in the abstract.

---

# Gated Direction After Sprint 34

The detailed implementation plans establish the following dependency order, but this roadmap does not assign sprint numbers beyond Sprint 34 until the preceding evidence gate passes:

The implementation-repository planning sources are:

- `docs/plans/integrated-roadmap.md`
- `docs/plans/ultraplan-local-server-experiment-plan.md`
- `docs/plans/server-shutdown-run-cancellation-contract.md`
- `docs/plans/sprint-code-context-stage.md`
- `docs/plans/retrieval-ready-content-plan.md`
- `docs/plans/post-execution-qa-and-repair-loop.md`
- `docs/plans/retrieval-ready-content-knowledge-graph-addendum.md`

Those plans remain detailed design inputs. This workspace roadmap owns product sequencing and promotion into committed sprint scope.

```text
minimal content identity and revision-aware provenance
-> read-only QA mapping, investigation, and synthesis
-> isolated evidence-producing QA and smoke integration
-> adjudicated manual repair, then bounded automatic repair
-> joint browser/code-context/content/QA dogfood
-> derived lexical retrieval
-> proven product persistence boundary and optional SQLite
-> optional in-memory knowledge graph
-> explicit authority choice, cloud, and Aren integration
```

## Gate A — Content Contract

- Begin with an artifact inventory and at least twenty real retrieval/traceability questions.
- Pilot optional minimal metadata, selective semantic block IDs, and revision-aware evidence on new artifacts.
- Preserve legacy Markdown validity and explicit project/sprint context selection.
- Stop if metadata is inaccurate, burdensome, or does not improve real discovery and traceability.

## Gate B — QA And Repair

- Keep current `review` compatibility while presenting it as Conformance Review.
- Introduce read-only deterministic QA shards and global synthesis before generated tests or production mutation.
- Require isolated writable workspaces and evidence adjudication before issue promotion.
- Permit repair only from frozen evidence-backed issue packets; bound cycles, reopenings, scope growth, and convergence.
- Preserve smoke protocol guarantees while absorbing smoke only after QA parity is proven.

## Gate C — Retrieval

- Derive disposable semantic records from authoritative artifacts after content dogfood.
- Establish metadata filtering and lexical search before embeddings.
- Measure against the retrieval-question corpus; do not silently inject results into governed context.
- Keep authority, status, provenance, exact source text, freshness, and explicit selection visible.

## Gate D — Persistence And SQLite

- Dogfood the filesystem-backed browser first and identify concrete needs such as drafts, immutable revisions, approvals, cross-project queries, provenance, or measured filesystem cost.
- Extract focused product-owned repository contracts one workflow at a time; do not create a universal storage or virtual-filesystem abstraction.
- Keep Git/source and agent execution workspaces filesystem-native.
- Select one authoritative storage mode at composition, use explicit migration/import/export, and prohibit silent dual writes.
- Compare filesystem and SQLite modes on real work before choosing server-canonical, dual first-class, or hybrid publication authority.

## Gate E — Knowledge Graph, Cloud, And Aren

- Consider a graph only after stable explicit relationships and a retrieval baseline prove multi-hop traceability value.
- Start with a deterministic in-memory read-only graph; preserve edge provenance, origin, status, revision, contradiction, and supersession.
- Treat any persisted graph as derived and rebuildable; do not assume a graph database.
- Move toward Postgres, durable workers, remote identity, and typed Aren artifact tools only after local persistence, authority, execution projection, cancellation, and recovery semantics are proven.

The normative local-server shutdown rule applies throughout server-backed phases: graceful server stop cancels every server-owned active run through the canonical cancellation path, waits for bounded cleanup, records a truthful outcome, and never equates browser disconnection with cancellation.

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

Before Sprint 30 starts:

- [ ] Phase 3 release gate is complete or any remaining exceptions are explicitly recorded.
- [ ] Phase 4 PRD, TRD, architecture, roadmap, project-index, and reasoning selections are accepted.
- [ ] `internal/app` exposes the required typed query, confirmation, operation, progress, cancellation, and error surfaces without CLI dependencies.
- [ ] Explicit CLI/TUI/web composition has a reviewed package-cycle-free design.
- [ ] Loopback-only, same-origin, CSRF, path-containment, SSE, shutdown, and fake-runtime browser test strategies are documented.

Before Phase 4 release:

- [ ] `/api/v1` compatibility and error envelopes are documented and fixture-tested.
- [ ] Browser read-only and guarded-operation state agrees with CLI/TUI/shared app results.
- [ ] SSE slow-client, reconnect, cancellation, shutdown, leak, and race tests pass.
- [ ] Local-web security, hostile-content, path-containment, and redaction audits pass.
- [ ] Gated representative real-runtime and smoke-harness browser evidence exists.
- [ ] A capability test proves that a new stage can expose status, artifacts, commands, progress, cancellation, and recovery without route-specific product logic.

Before Sprint 32 starts:

- [ ] Sprint 31's technical handbook is valid and its area reasoning, final reasoning, and plan are complete.

Before Product Phase 5 release:

- [ ] Existing sprint state remains compatible after inserting `code-context` after requirements.
- [ ] CLI, TUI, and browser parity covers the new stage through shared application abstractions.
- [ ] Prefix-stability fixtures prove exact requirements and code-context reuse across downstream agent prompts.
- [ ] A representative real-repository dogfood flow completes from requirements through plan with `code-context` run exactly once.
- [ ] Content identity, retrieval, caching, automatic staleness, and QA expansion remain outside this chunk.
