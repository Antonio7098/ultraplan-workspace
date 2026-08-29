# Technical Requirements Document: UltraPlan Go

**Version:** 1.8.0
**Status:** Draft
**Owner:** Engineering
**Last Updated:** 2026-08-20

## 1. Purpose

This TRD defines the technical requirements for UltraPlan Go, a production-grade local CLI, TUI, and browser surface that implement the proven UltraPlan workflow. Phase 1 covers study initialization, source analysis, report synthesis, code-reference extraction, resumable orchestration, validation, and operational diagnostics. Phase 2 adds governed project and sprint planning through `plan.md`, then controlled implementation execution through `execute`. Sprints 24 and 25 deliver the local TUI foundation and guarded controls. Phase 3 adds automated conformance review through `review.md`, then sprint-targeted deep smoke through `smoke.md` backed by the external harness cataloged in `project-index.md`. Phase 4 begins with Sprint 30 and adds a loopback-only Go HTTP server plus a simple Go-rendered browser UI with SSE progress. Sprints 33–34 form the grounded-planning track, inserting a validated `code-context` stage and reusing its exact Markdown source pack downstream. Product Phase 5 comprises Sprints 35–39: durable run identity and cross-surface observability, read-only QA decomposition and synthesis, isolated evidence-producing QA with smoke integration, bounded repair, and QA/repair dogfooding and hardening.

This document is implementation-oriented. It defines boundaries, modules, data models, state machines, validators, runtime contracts, error handling, and testing requirements. It does not prescribe every package name or third-party library, but it should be specific enough to guide implementation.

UltraPlan Go must use `github.com/Antonio7098/agentwrap` and its `opencode` adapter for agentic runtime supervision. UltraPlan should not reimplement agentwrap's runtime contract, OpenCode process handling, canonical events, policy execution, output validation, repair, permission translation, health checks, observability records, or run inspection primitives.

## 2. System Scope

UltraPlan Go is responsible for:

- Managing a local UltraPlan workspace.
- Initializing and validating study definitions.
- Running source/dimension analysis tasks through runtime adapters.
- Synthesizing final reports from per-source reports.
- Maintaining durable run state for long-running batches.
- Validating generated artifacts.
- Extracting code snippets from report citations.
- Managing project planning roots under `projects/<project>`.
- Validating project indexes that catalog contracts, evidence, reasoning templates, and review protocols.
- Creating and validating sprint planning artifacts through `plan.md`.
- Executing validated sprint plan tasks through the generic runtime boundary with durable task state.
- Reviewing implemented sprint scope against selected contracts, selected review protocols, technical-handbook guidance, sprint decisions, plan tasks, and verification evidence.
- Running sprint-targeted deep smoke through a cataloged external harness after review.
- Writing the current human-readable sprint `review.md` and `smoke.md` summaries while linking detailed smoke evidence in the harness.
- Providing human-readable and structured operational output.
- Providing a local TUI over the same workspace and workflow services after the CLI workflows are stable.
- Providing a loopback-only local HTTP server and embedded browser UI over the same typed application use cases beginning in Sprint 30.
- Generating a sprint-owned `code-context.md` from read-only inspection of the resolved implementation repository and reusing it as the common downstream prompt foundation beginning in Sprint 33.
- Persisting a stable operational run identity and safe ordered history before execution is exposed, then projecting it consistently through CLI, JSON, TUI, and supported local server instances beginning in Sprint 35.

UltraPlan Go is not responsible for:

- Hosting a multi-user service.
- Providing a hosted, remotely exposed, or multi-user browser service.
- Managing team permissions.
- Replacing source control, issue trackers, or project management systems.
- Owning AI provider billing.
- Guaranteeing semantic correctness of generated prose beyond validation rules.
- Reimplementing agentwrap SDK features that already exist as runtime-neutral primitives.
- Providing a general-purpose issue tracker, remote issue synchronization, or project-management workflow.
- Automatically fixing product code or tests during review or smoke.
- Mutating Git state during planning, execute, review, or smoke.
- Owning detailed smoke run/issue persistence that belongs to the external harness.
- Replacing the CLI or JSON surfaces with a TUI-only workflow.
- Replacing workspace artifacts with browser-owned or server-only durable state.
- Treating operational run persistence as authority for governed Markdown, stage outcomes, Git/source state, or external smoke evidence.
- Treating a code-context pack as a repository index, exclusive source boundary, cache database, or second machine-readable manifest.
- Introducing content identity, retrieval, SQLite product persistence, knowledge-graph persistence, cloud, or Aren capabilities before their explicit post-Sprint-39 evidence gates.

## 3. Architecture Overview

UltraPlan Go should be module-driven, not global-layer-driven. Product behavior should live with the module that owns the relevant state and workflow. Shared platform packages should be limited to genuinely cross-cutting infrastructure such as configuration, logging, filesystem helpers, and runtime execution.

The detailed architecture is defined in [ARCHITECTURE.md](ARCHITECTURE.md). The initial implementation should follow this high-level shape:

```text
cmd/
  ultraplan/
    main.go

internal/
  app/                  # composition root, shared use cases, and interface wiring
  tui/                  # local terminal UI over app use cases, added after execute
  web/                  # loopback HTTP, Go templates/static assets, JSON, and SSE over app use cases
  platform/             # generic infrastructure only
    config/
    logging/
    filesystem/
    runtime/
    process/            # safe external executable/argv execution for the smoke harness
  workspace/            # workspace discovery, paths, and validation
  study/                # study lifecycle, prompts, scheduling, validation, reports, state
  project/              # project docs, project-index cataloging and validation
  sprint/               # planning, execute, review, smoke, prompts, validators, state, and flow through smoke
  codeextract/          # citation parsing, file resolution, snippet extraction
```

Sprint 35 adds the focused `internal/runcontrol` package, exposed through `internal/app`. It owns the adapter-neutral operational model and same-host SQLite repository; product workflow semantics remain in `study` and `sprint`, `internal/web` remains a transport adapter, and no interface process owns the only authoritative run registry.

The core rule is:

```text
Platform owns generic capabilities.
Modules own product behavior.
Logic stays near the state it transforms.
Interfaces appear only at external or volatile boundaries.
```

Do not create global technical-layer packages such as `internal/validation`, `internal/scheduler`, `internal/reports`, or `internal/prompts` unless the behavior is genuinely reusable across multiple product modules. Prefer keeping behavior inside the module that owns the state and workflow, for example `internal/study/validation.go`, `internal/study/scheduler.go`, and `internal/study/reports.go`.

Runtime supervision is delegated to `agentwrap`. The platform runtime package must stay generic: it may know about prompts, working directories, models, timeouts, permissions, events, and execution results, but it must not know about studies, dimensions, sources, synthesis gating, report semantics, project catalogs, sprint stages, or product state machines.

The CLI must not become the only place where use cases are assembled. Shared application operations should be extracted from command-specific glue into testable structs/functions that CLI commands, TUI actions, and local HTTP handlers can call. Neither TUI nor web code should invoke `ultraplan` subprocesses for normal product behavior.

## 4. Design Principles

- Keep the domain model independent of CLI parsing.
- Keep product use cases independent of CLI parsing, TUI widget state, HTTP transport DTOs, and HTML rendering.
- Keep runtime adapters independent of study semantics.
- Treat filesystem artifacts as durable product state.
- Validate outputs before marking tasks successful.
- Make state transitions explicit and testable.
- Prefer deterministic fixtures for tests.
- Do not require real OpenCode for normal unit tests.
- Use agentwrap for runtime execution, retries, validation wrapping, permission policy translation, health/preflight checks, observability, metadata, and OpenCode structured-output adaptation.
- Keep UltraPlan-specific study, source, report, project, and sprint behavior outside agentwrap adapters and wrappers.
- Avoid global mutable state.
- Avoid package cycles.
- Make dry-run behavior available before expensive operations.
- Keep generated Markdown and YAML readable by humans.
- Keep the CLI as the stable automation contract even when a TUI is available.
- Keep the filesystem and product-owned run state authoritative when the local server is running.
- Run review before smoke by default so deterministic conformance failures block unnecessary live-runtime work.
- Keep only `review.md` and `smoke.md` in the sprint root; link detailed smoke evidence from the external harness.
- Compute review/smoke verdicts from validated evidence and explicit severity rules, not from runtime exit success or unstructured model prose alone.
- Preserve exact requirements and code-context bytes in one shared prompt prefix before stage-specific instructions; downstream agents may still inspect additional live source.
- Keep Markdown authoritative initially. Any later search record or knowledge graph is derived and rebuildable; any later SQLite authority is an explicit product-mode decision rather than a web-side shadow store.

## 4.1 Application Surface Requirements

The `internal/app` package is the composition boundary for local interfaces. It should provide:

- shared dependency construction for workspace discovery, config loading, runtime setup, and service creation
- typed use-case functions for project, sprint, study, code extraction, validation, status, flow, execute, review, smoke, and verify operations
- CLI adapters that parse arguments, call use cases, and render text/JSON
- TUI adapters that translate key actions into the same use cases and render terminal models
- Web adapters that map HTTP requests to the same use cases and render HTML, JSON, or SSE without owning product behavior

The intended direction is:

```text
cmd/ultraplan -> internal/app -> product modules/platform
internal/tui  -> internal/app -> product modules/platform
internal/web  -> internal/app -> product modules/platform
```

Avoid this direction:

```text
internal/tui -> CLI command handlers -> stdout parsing
internal/tui -> os/exec("ultraplan", ...)
internal/web -> CLI command handlers or os/exec("ultraplan", ...)
```

Subprocess execution is acceptable only for explicit external runtime behavior already owned by `platform/runtime` and `agentwrap`, not for UltraPlan calling itself.

## 5. Workspace Requirements

### 5.1 Workspace Discovery

The CLI must resolve a workspace root using this precedence:

1. Explicit `--workspace <path>`.
2. `ULTRAPLAN_WORKSPACE` environment variable.
3. Current directory if it contains a workspace marker or required structure.
4. Nearest parent directory containing a workspace marker or required structure.

Workspace discovery must not traverse indefinitely. It must stop at filesystem root.

### 5.2 Workspace Structure

Required top-level structure:

```text
.
  README.md
  ultraplan.yml
  studies/
```

Optional structure:

```text
.
  projects/
  prompts/      # intentional overrides of embedded defaults only
  templates/    # intentional overrides of embedded defaults only
  .ultraplan/
    cache/
    logs/
    locks/
    tmp/
```

### 5.3 Path Handling

- All workspace-managed paths must be normalized.
- Commands must reject paths that escape the workspace when the command is expected to operate only inside the workspace.
- Source directories may be local paths inside a study.
- Absolute paths may be accepted in config only when explicitly allowed.
- Generated artifacts should prefer workspace-relative paths.

## 6. Configuration Requirements

### 6.1 Config File

Primary config file: `ultraplan.yml`.

Required fields:

```yaml
version: 1
runtime:
  default: opencode
models:
  default: provider/model
  primary: provider/model
  backup: provider/model
  stages:
    sprint-index: provider/model
    code-context: provider/model
    technical-handbook: provider/model
    area-reasoning: provider/model
    reasoning: provider/model
    plan: provider/model
    execute: provider/model
execution:
  default_variant: high
  default_parallel: 3
  default_timeout: 30m
  default_retries: 3
logging:
  format: text
  level: info
agentwrap:
  executable: opencode
  required_health:
    - runtime_available
    - structured_output
    - workdir
```

### 6.2 Config Precedence

Effective config must be resolved in this order:

1. Built-in defaults.
2. Workspace config.
3. Environment variables.
4. Command-specific flags.

The effective config object must record which values were defaulted and which were user-provided when relevant for diagnostics.

### 6.3 Secret Handling

- Secrets must not be required in `ultraplan.yml`.
- Runtime/provider secrets should be inherited from the environment or runtime-native config.
- `config show` must redact fields marked sensitive.
- Logs must not print full environment variables by default.

### 6.4 Config Validation

Config validation must check:

- Schema version.
- Required fields.
- Valid duration syntax.
- Positive parallelism, retries, and timeouts.
- Known runtime names.
- Known logging formats.
- Model values are non-empty when required.
- Agentwrap required health check names map to `agentwrap.HealthCheckID` values.
- Agentwrap executable path is non-empty when configured.

Validation failures must include field path and corrective guidance.

### 6.5 Agentwrap Configuration Mapping

UltraPlan config must map into agentwrap concepts instead of inventing parallel runtime concepts:

- Runtime executable maps to `opencode.WithExecutable`.
- Extra runtime args map to `opencode.WithExtraArgs`.
- Runtime environment additions map to `opencode.WithEnv`, with secrets redacted from diagnostics.
- Provider and model map to `agentwrap.RunRequest.Provider` and `agentwrap.RunRequest.Model`.
- Stage-specific model overrides map to `agentwrap.RunRequest.Provider` and `agentwrap.RunRequest.Model` for the selected sprint stage.
- Timeout maps to `agentwrap.RunRequest.Timeout`.
- Permission mode maps to `agentwrap.RunRequest.Permissions`.
- Structured permission policy maps to `agentwrap.RunRequest.PermissionPolicy`.
- Sandbox mode maps to `agentwrap.RunRequest.Sandbox`.
- Required health checks map to `agentwrap.RunRequest.RequireHealth`.
- Required capabilities map to `agentwrap.RunRequest.RequireCaps`.
- Output expectations and repair policy map to `agentwrap.RunRequest.Validation`.

UltraPlan should use agentwrap's `EffectiveConfig`, source-aware config layers, and redaction helpers where practical for runtime-facing configuration summaries.

## 7. CLI Requirements

### 7.1 CLI Behavior

- Every command must return a meaningful exit code.
- `--help` must be available globally and per command.
- `--json` must be supported for status and validation commands.
- `--dry-run` must be supported for expensive or mutating operations where useful.
- `--workspace` must be available globally.
- Commands must be script-friendly and avoid interactive prompts unless explicitly requested.

### 7.2 Exit Codes

Suggested exit code classes:

- `0`: Success.
- `1`: General failure.
- `2`: Usage or argument error.
- `3`: Config error.
- `4`: Workspace or filesystem error.
- `5`: Validation failure.
- `6`: Runtime failure.
- `7`: Cancellation.
- `8`: Partial completion.

### 7.3 Output Modes

Text mode:

- Default for humans.
- Shows progress and actionable summaries.

JSON mode:

- Stable structure for automation.
- No ANSI formatting.
- Includes command name, workspace, status, timestamps, and result payload.

## 7.4 TUI Requirements

The TUI is a local terminal surface in the same `ultraplan` binary, exposed as `ultraplan tui` unless a later CLI design review chooses a different spelling.

Required baseline behavior:

- discover the active workspace using the same workspace rules as the CLI
- show projects, studies, and sprints in a navigable dashboard
- show project status, sprint flow status, study status, validation findings, and key artifact paths
- open read-only artifact previews where practical
- preserve clear error messages and exit behavior when startup fails
- use shared app use cases and typed domain results, not CLI text scraping

Operational behavior, added after the read-only baseline:

- run validation commands from the TUI
- run dry-run planning flows and prompt previews
- start guarded planning/execute/study workflows only after showing the operation and expected mutation scope
- display live progress from existing progress/event callbacks
- support cancellation through `context.Context`
- leave durable state as the source of truth after cancellation or terminal resize/exit
- show review readiness, scope, selected contracts, reviewer progress, findings, verdict, reruns, and `review.md`
- show smoke readiness, selected harness scope, prerequisites, expected cost/duration class, suite/test progress, external evidence links, open issues, verdict, reruns, and `smoke.md`
- expose the complete `execute -> review -> smoke` workflow through the same typed app use cases as the CLI

Quality requirements:

- normal unit tests must not require an interactive terminal
- TUI model/update logic should be testable with deterministic messages
- runtime-backed TUI behavior must use fake runtimes in normal tests and gated real-runtime smoke where available
- smoke-backed TUI behavior must use a fake harness in normal tests and the cataloged external harness only in gated tests
- TUI rendering must degrade gracefully in narrow terminals

## 7.5 Local HTTP Server And Browser UI Requirements

Product Phase 4 begins with Sprint 30. The same `ultraplan` binary must expose the local web surface through:

```bash
ultraplan serve
```

Server requirements:

- use Go `net/http`; use `html/template` for pages and embedded CSS/minimal JavaScript for browser behavior
- bind to `127.0.0.1` or `::1` by default and reject non-loopback binding in Phase 4
- support an explicit loopback listen address, optional browser opening, signal-aware graceful shutdown, and bounded HTTP timeouts
- expose versioned JSON endpoints for dashboard/detail queries, bounded artifact previews, validation, confirmation, operation start/status, and cancellation
- expose operation progress as `text/event-stream` using SSE; ordinary commands remain explicit HTTP requests
- call typed app use cases rather than CLI command handlers or product modules directly
- keep HTTP request/response types and template models inside `internal/web`
- serve no arbitrary workspace paths and accept no executable command assembled from browser input
- require no Node.js, Vite, separate frontend server, database, or asset build step at runtime

Browser UI requirements:

- show projects, studies, sprints, validation findings, workflow state, and bounded Markdown/JSON artifact previews
- begin with server-rendered pages and progressive enhancement; a client-side application framework is not required for Phase 4
- organize templates into `primitives`, `components`, `layouts`, and `pages`; allow dependencies only in that downward order
- give template definitions stable namespaced names and render them from explicit typed view models prepared by handlers
- keep filesystem reads, application calls, HTTP validation, product-state interpretation, and durable state mutation out of templates
- layer CSS as tokens, base, primitives, components, layouts, and utilities; split JavaScript only by narrow progressive-enhancement capability
- render untrusted Markdown without executing embedded HTML or scripts
- show operation scope, affected paths, runtime/model information, and mutation class before starting guarded work
- use a server-issued, short-lived confirmation bound to the normalized request and current governed-input fingerprint for mutating or runtime-backed work
- stream bounded progress events, support cancellation, and refresh durable state after completion, failure, cancellation, or reconnect
- remain truthful when the browser disconnects: subscriber presence is not task state and loss of an SSE connection does not mark work complete or failed

Local server security requirements:

- same-origin pages and API, no permissive CORS
- strict Host and Origin checks, CSRF protection for mutations, request-body limits, security response headers, and bounded concurrent streams
- safe cookies or an equivalent per-process session mechanism where browser state is required
- path containment, artifact allowlisting, secret redaction, and safe error projection identical in strength to CLI/TUI behavior
- no hosted authentication, team permissions, tenant isolation, or remote-worker protocol in Phase 4

Normal tests must use `httptest`, fake app use cases, deterministic templates, fake runtimes, and fake smoke harnesses. Required coverage includes route/method errors, JSON compatibility, confirmation expiry/staleness, path escape, CSRF/origin rejection, hostile Markdown, SSE ordering/reconnect/slow subscribers, cancellation, shutdown, redaction, and CLI/TUI/web agreement.

## 7.6 Sprint Code-Context Requirements

The delivered grounded-planning track adds `code-context` immediately after `requirements`:

```text
requirements
-> code-context
-> sprint-index
-> technical-handbook
-> area-reasoning
-> reasoning
-> plan
```

Required behavior:

- expose `prompt code-context`, `validate code-context`, and `flow --to code-context` through the existing sprint command/application surfaces
- store exactly one authoritative artifact at `projects/<project>/sprints/<sprint>/code-context.md`; do not add a parallel JSON context manifest
- resolve the target implementation repository through the project index and existing execute/worktree mechanisms
- allow repository reads while limiting stage writes to `code-context.md` in the sprint root
- require validated requirements before execution and require structural code-context validation before sprint-index becomes ready
- include scope interpretation, inspected repository areas, selected exact source excerpts, repository-relative paths, optional well-formed ranges and symbols, relevance explanations, important relationships, constraints, and open questions
- reject unsafe absolute/escaping source paths, placeholders, missing excerpts, missing rationale, and malformed ranges
- atomically replace only `code-context.md` on an explicit rerun
- preserve compatibility for workspaces and flow state created before the stage existed
- provide stage-specific model/variant configuration using existing fallback behavior
- expose readiness, progress, findings, artifact preview, rerun, cancellation, and recovery through shared CLI/TUI/web use cases

Downstream prompt composition must use one shared renderer with this order:

```text
stable shared planning instructions
sprint identity
exact requirements.md
exact code-context.md
other shared context
stage-specific instructions and output contract
```

The requirements/code-context block must remain byte-for-byte identical across compatible downstream calls and must not contain stage names, timestamps, run IDs, output paths, or other dynamic run data. Sprint index, technical handbook, area/final reasoning, plan, execute, Conformance Review, and smoke must receive it whenever they invoke an agent. The pack is a prepared foundation, not an access restriction; agents may inspect additional repository files.

The first implementation must not add a repository index, RAG or embedding system, UltraPlan cache/key subsystem, provider-specific cache-control dependency, automatic staleness detector, context amendment protocol, or hard maximum excerpt count.

## 8. Domain Model

### 8.1 Study

```go
type Study struct {
    Name        string
    Description string
    Root        string
    Sources     []Source
    Dimensions  []Dimension
}
```

Requirements:

- `Name` must be filesystem-safe.
- `Root` must be under `studies/` unless a custom output root is supplied.
- Sources and dimensions must be sorted deterministically for listing and task creation.

### 8.2 Source

```go
type SourceKind string

const (
    SourceKindDirectory SourceKind = "directory"
    SourceKindMarkdown  SourceKind = "markdown"
)

type Source struct {
    Name                 string
    Path                 string
    URL                  string
    Description          string
    Kind                 SourceKind
    ApplicableDimensions []string
    Frontmatter          map[string]any
}
```

Requirements:

- `Name` must be unique inside a study.
- Directory source `Path` must resolve to a directory for execution.
- Markdown source `Path` must resolve to a `.md` file directly under the study's `sources/` directory.
- Markdown source `ApplicableDimensions` must contain normalized two-digit dimension numbers.
- Empty `ApplicableDimensions` means the source applies to all dimensions.
- Source lookup by prefix is allowed only when unambiguous.
- Source discovery must distinguish directory sources from Markdown document sources.

### 8.3 Dimension

```go
type Dimension struct {
    Number      string
    Slug        string
    Title       string
    File        string
    Purpose     string
    Steps       []string
    Citations   []string
    Questions   []string
}
```

Requirements:

- `Number` must be zero-padded.
- `Slug` must be kebab-case.
- Dimension lookup by number, slug, full filename, or unambiguous prefix must be supported.

### 8.4 Report

```go
type Report struct {
    Kind       ReportKind
    Study      string
    Dimension  string
    Source     string
    Path       string
    Status     ValidationStatus
    Score      *float64
}
```

Report kinds:

- Per-source.
- Final synthesis.
- Code extraction bundle.
- Planning handbook.
- Planning reasoning.
- Planning plan.

### 8.5 Project

```go
type Project struct {
    Name       string
    Root       string
    PRDPath    string
    TRDPath    string
    RoadmapPath string
    Sprints    []Sprint
}
```

### 8.6 Planning Sprint

```go
type Sprint struct {
    Project      string
    Slug         string
    Root         string
    ReasoningPath string
    PlanPath     string
    FlowStatePath string
}
```

## 9. Study Initialization

### 9.1 Input Schema

The study initialization YAML must support:

- Study name.
- Description.
- Repository/source target count.
- Repository/source items.
- Dimension target count.
- Dimension items.
- Optional source clone behavior.
- Optional output directory.

The parser must reject unknown required structures only when they prevent safe execution. It may ignore unknown extension fields with a warning if version-compatible.

### 9.2 Assisted Completion

If `repos.count` is greater than the number of explicit repo items, study initialization may request additional source suggestions from the configured runtime.

If `dimensions.count` is greater than the number of explicit dimension items, study initialization may request additional dimension suggestions from the configured runtime.

Runtime suggestions must use structured JSON outputs.

Suggested source schema:

```go
type SuggestedSource struct {
    Name        string `json:"name"`
    URL         string `json:"url"`
    Description string `json:"description"`
}
```

Suggested dimension schema:

```go
type SuggestedDimension struct {
    Number    string   `json:"number"`
    Slug      string   `json:"name"`
    Title     string   `json:"title"`
    Purpose   string   `json:"purpose"`
    Steps     []string `json:"steps"`
    Citations []string `json:"citations"`
    Questions []string `json:"questions"`
}
```

Requirements:

- Assisted completion must be disabled by `--no-assist`.
- Dry-run mode must report the shortage without invoking the runtime.
- Suggestions must be cached under `.ultraplan/cache/study-init/<study>/`.
- Suggestions must be validated before merge.
- Suggestions must be deduplicated by name and URL for sources.
- Suggestions must be deduplicated by number and slug for dimensions.
- Invalid suggestions must be ignored with warnings unless no valid suggestions remain.
- The normalized `study-init.yml` must include accepted suggestions.
- The cache must record runtime, model, prompt hash, created time, and raw response path.

### 9.3 Generated Files

Study initialization must create:

```text
studies/<study>/
  study-init.yml
  README.md
  dimensions/
    NN-slug.md
  sources/
  reports/
    source/
    final/
```

### 9.4 Source Cloning

When cloning is enabled:

- Use `git clone --depth 1`.
- Clone each source into `sources/<source-name>`.
- Skip existing source directories unless `--force` or a specific clean option is used.
- Record clone failures without hiding successful clones.
- Return partial-completion status when some clones fail.
- Optionally verify source URLs before cloning when `--verify-sources` is enabled.
- If URL verification fails and assisted completion is enabled, the CLI may request one replacement source for the same study purpose.
- Replacement sources must pass the same validation and deduplication as other suggestions.

### 9.5 Force Behavior

`--force` may overwrite generated study structure. It must not delete unrelated paths outside the selected study directory.

If destructive replacement is needed, implementation must clearly scope the deletion and should prefer moving old content to a backup or requiring a more explicit flag for deletion.

## 9A. Source Discovery and Applicability

### 9A.1 Source Discovery

Source discovery must inspect `studies/<study>/sources/` and return:

- One directory source for each non-hidden directory.
- One Markdown document source for each top-level `.md` file.

Source discovery must ignore:

- Hidden entries.
- Non-directory, non-`.md` files.
- Nested Markdown files inside directory sources, because those belong to the directory source's repository content.

### 9A.2 Markdown Frontmatter

Markdown document sources may begin with YAML frontmatter:

```markdown
---
applicable_dimensions:
  - 1
  - 3
  - 5
---
# My Guide
```

Required helper behavior:

- `parseFrontmatter(content string) (map[string]any, error)` parses leading `---` frontmatter only.
- `stripFrontmatter(content string) string` removes leading frontmatter and returns the document body.
- Files without leading frontmatter return an empty frontmatter map and the original content.
- Malformed frontmatter should produce a source validation warning or error according to command context.

`applicable_dimensions` rules:

- Accepted values may be numbers or strings.
- Values must normalize to two-digit dimension numbers.
- `1`, `"1"`, and `"01"` all match dimension `01`.
- Empty or absent `applicable_dimensions` means all dimensions are applicable.
- Invalid values must be reported with the source path and offending value.

### 9A.3 Applicability Filtering

The implementation must provide a helper equivalent to:

```go
func GetApplicableSources(sources []Source, dimension Dimension) []Source
```

Requirements:

- Directory sources are always applicable.
- Markdown sources with no applicability filter are applicable to all dimensions.
- Markdown sources with `ApplicableDimensions` are applicable only when the selected dimension number matches.
- Filtering must be used everywhere source/dimension task pairs are created, validated, summarized, or synthesized.

The following flows must respect applicability:

- Single analysis command.
- Full study run.
- Stateful run-loop initial state creation.
- Stateful run-loop resume validation.
- Completed source detection.
- Synthesis gating.
- Summary generation.
- Status calculations.
- Artifact validation for expected reports.

Inapplicable pairs must be skipped, not failed.

## 10. Prompt Composition

### 10.1 Prompt Inputs

Directory source analysis prompt inputs:

- Shared base prompt.
- Dimension Markdown.
- Report template.
- Study name.
- Source name and path.
- Output path.
- Hard rules for source isolation and citation format.

Markdown document source analysis prompt inputs:

- Shared base prompt, with code-exploration rules replaced or scoped for document analysis.
- Dimension Markdown.
- Report template.
- Study name.
- Source name and path.
- Stripped Markdown document content.
- Output path.
- Explicit instruction that all source material is embedded in the prompt.
- Explicit instruction not to access external files, repositories, or code.
- Explicit instruction that code citation requirements do not apply unless the dimension states otherwise.

Synthesis prompt inputs:

- Synthesis prompt.
- Dimension Markdown.
- Final report template.
- List of per-source report paths.
- Output path.

Sprint planning prompt inputs are defined in section 18. Study prompt composition covers analysis and synthesis only.

### 8.1 Study

```go
type Study struct {
    Name        string
    Description string
    Root        string
    Sources     []Source
    Dimensions  []Dimension
}
```

Requirements:

- `Name` must be filesystem-safe.
- `Root` must be under `studies/` unless a custom output root is supplied.
- Sources and dimensions must be sorted deterministically for listing and task creation.

### 8.2 Source

```go
type SourceKind string

const (
    SourceKindDirectory SourceKind = "directory"
    SourceKindMarkdown  SourceKind = "markdown"
)

type Source struct {
    Name                 string
    Path                 string
    URL                  string
    Description          string
    Kind                 SourceKind
    ApplicableDimensions []string
    Frontmatter          map[string]any
}
```

Requirements:

- `Name` must be unique inside a study.
- Directory source `Path` must resolve to a directory for execution.
- Markdown source `Path` must resolve to a `.md` file directly under the study's `sources/` directory.
- Markdown source `ApplicableDimensions` must contain normalized two-digit dimension numbers.
- Empty `ApplicableDimensions` means the source applies to all dimensions.
- Source lookup by prefix is allowed only when unambiguous.
- Source discovery must distinguish directory sources from Markdown document sources.

### 8.3 Dimension

```go
type Dimension struct {
    Number      string
    Slug        string
    Title       string
    File        string
    Purpose     string
    Steps       []string
    Citations   []string
    Questions   []string
}
```

Requirements:

- `Number` must be zero-padded.
- `Slug` must be kebab-case.
- Dimension lookup by number, slug, full filename, or unambiguous prefix must be supported.

### 8.4 Report

```go
type Report struct {
    Kind       ReportKind
    Study      string
    Dimension  string
    Source     string
    Path       string
    Status     ValidationStatus
    Score      *float64
}
```

Report kinds:

- Per-source.
- Final synthesis.
- Code extraction bundle.
- Planning handbook.
- Planning reasoning.
- Planning plan.

### 8.5 Project

```go
type Project struct {
    Name       string
    Root       string
    PRDPath    string
    TRDPath    string
    RoadmapPath string
    Sprints    []Sprint
}
```

### 8.6 Planning Sprint

```go
type Sprint struct {
    Project      string
    Slug         string
    Root         string
    ReasoningPath string
    PlanPath     string
    FlowStatePath string
}
```

## 9. Study Initialization

### 9.1 Input Schema

The study initialization YAML must support:

- Study name.
- Description.
- Repository/source target count.
- Repository/source items.
- Dimension target count.
- Dimension items.
- Optional source clone behavior.
- Optional output directory.

The parser must reject unknown required structures only when they prevent safe execution. It may ignore unknown extension fields with a warning if version-compatible.

### 9.2 Assisted Completion

If `repos.count` is greater than the number of explicit repo items, study initialization may request additional source suggestions from the configured runtime.

If `dimensions.count` is greater than the number of explicit dimension items, study initialization may request additional dimension suggestions from the configured runtime.

Runtime suggestions must use structured JSON outputs.

Suggested source schema:

```go
type SuggestedSource struct {
    Name        string `json:"name"`
    URL         string `json:"url"`
    Description string `json:"description"`
}
```

Suggested dimension schema:

```go
type SuggestedDimension struct {
    Number    string   `json:"number"`
    Slug      string   `json:"name"`
    Title     string   `json:"title"`
    Purpose   string   `json:"purpose"`
    Steps     []string `json:"steps"`
    Citations []string `json:"citations"`
    Questions []string `json:"questions"`
}
```

Requirements:

- Assisted completion must be disabled by `--no-assist`.
- Dry-run mode must report the shortage without invoking the runtime.
- Suggestions must be cached under `.ultraplan/cache/study-init/<study>/`.
- Suggestions must be validated before merge.
- Suggestions must be deduplicated by name and URL for sources.
- Suggestions must be deduplicated by number and slug for dimensions.
- Invalid suggestions must be ignored with warnings unless no valid suggestions remain.
- The normalized `study-init.yml` must include accepted suggestions.
- The cache must record runtime, model, prompt hash, created time, and raw response path.

### 9.3 Generated Files

Study initialization must create:

```text
studies/<study>/
  study-init.yml
  README.md
  dimensions/
    NN-slug.md
  sources/
  reports/
    source/
    final/
```

### 9.4 Source Cloning

When cloning is enabled:

- Use `git clone --depth 1`.
- Clone each source into `sources/<source-name>`.
- Skip existing source directories unless `--force` or a specific clean option is used.
- Record clone failures without hiding successful clones.
- Return partial-completion status when some clones fail.
- Optionally verify source URLs before cloning when `--verify-sources` is enabled.
- If URL verification fails and assisted completion is enabled, the CLI may request one replacement source for the same study purpose.
- Replacement sources must pass the same validation and deduplication as other suggestions.

### 9.5 Force Behavior

`--force` may overwrite generated study structure. It must not delete unrelated paths outside the selected study directory.

If destructive replacement is needed, implementation must clearly scope the deletion and should prefer moving old content to a backup or requiring a more explicit flag for deletion.

## 9A. Source Discovery and Applicability

### 9A.1 Source Discovery

Source discovery must inspect `studies/<study>/sources/` and return:

- One directory source for each non-hidden directory.
- One Markdown document source for each top-level `.md` file.

Source discovery must ignore:

- Hidden entries.
- Non-directory, non-`.md` files.
- Nested Markdown files inside directory sources, because those belong to the directory source's repository content.

### 9A.2 Markdown Frontmatter

Markdown document sources may begin with YAML frontmatter:

```markdown
---
applicable_dimensions:
  - 1
  - 3
  - 5
---
# My Guide
```

Required helper behavior:

- `parseFrontmatter(content string) (map[string]any, error)` parses leading `---` frontmatter only.
- `stripFrontmatter(content string) string` removes leading frontmatter and returns the document body.
- Files without leading frontmatter return an empty frontmatter map and the original content.
- Malformed frontmatter should produce a source validation warning or error according to command context.

`applicable_dimensions` rules:

- Accepted values may be numbers or strings.
- Values must normalize to two-digit dimension numbers.
- `1`, `"1"`, and `"01"` all match dimension `01`.
- Empty or absent `applicable_dimensions` means all dimensions are applicable.
- Invalid values must be reported with the source path and offending value.

### 9A.3 Applicability Filtering

The implementation must provide a helper equivalent to:

```go
func GetApplicableSources(sources []Source, dimension Dimension) []Source
```

Requirements:

- Directory sources are always applicable.
- Markdown sources with no applicability filter are applicable to all dimensions.
- Markdown sources with `ApplicableDimensions` are applicable only when the selected dimension number matches.
- Filtering must be used everywhere source/dimension task pairs are created, validated, summarized, or synthesized.

The following flows must respect applicability:

- Single analysis command.
- Full study run.
- Stateful run-loop initial state creation.
- Stateful run-loop resume validation.
- Completed source detection.
- Synthesis gating.
- Summary generation.
- Status calculations.
- Artifact validation for expected reports.

Inapplicable pairs must be skipped, not failed.

## 10. Prompt Composition

### 10.1 Prompt Inputs

Directory source analysis prompt inputs:

- Shared base prompt.
- Dimension Markdown.
- Report template.
- Study name.
- Source name and path.
- Output path.
- Hard rules for source isolation and citation format.

Markdown document source analysis prompt inputs:

- Shared base prompt, with code-exploration rules replaced or scoped for document analysis.
- Dimension Markdown.
- Report template.
- Study name.
- Source name and path.
- Stripped Markdown document content.
- Output path.
- Explicit instruction that all source material is embedded in the prompt.
- Explicit instruction not to access external files, repositories, or code.
- Explicit instruction that code citation requirements do not apply unless the dimension states otherwise.

Synthesis prompt inputs:

- Synthesis prompt.
- Dimension Markdown.
- Final report template.
- List of per-source report paths.
- Output path.

Sprint planning prompt inputs are defined in section 18. Study prompt composition covers analysis and synthesis only.

### 10.2 Prompt Builder Requirements

- Prompt builders must be deterministic.
- Prompt builders must be unit-tested with golden fixtures.
- Prompt builders must return both prompt text and input manifest.
- Dry-run mode must expose the input manifest.
- Prompt/template lookup must prefer an intentional readable workspace override and otherwise use the embedded default. It must fail before runtime execution only when neither source is available or when an existing override is unreadable or invalid.
- Directory source prompts must preserve source-isolation and file-line citation rules.
- Markdown document source prompts must embed stripped document content and use document-analysis instructions.
- Prompt builders must not embed YAML frontmatter from Markdown document sources.

## 11. Runtime Adapter Requirements

### 11.1 Agentwrap Dependency

UltraPlan Go must use `github.com/Antonio7098/agentwrap` as the runtime SDK and `github.com/Antonio7098/agentwrap/opencode` as the first concrete runtime adapter.

UltraPlan must not define a competing public runtime contract. Product code may define thin internal interfaces for testing and dependency injection, but the implementation boundary to agentic runtimes must be agentwrap's root package API:

```go
type Runtime interface {
    StartRun(context.Context, RunRequest) (Run, error)
    Capabilities(context.Context) (Capabilities, error)
}
```

UltraPlan runtime integration must use:

- `agentwrap.Runtime`
- `agentwrap.Run`
- `agentwrap.RunRequest`
- `agentwrap.RunResult`
- `agentwrap.Capabilities`
- `agentwrap.SDKError`
- `agentwrap.HealthChecker` where the concrete runtime supports it
- `agentwrap.PolicyRunner`
- `agentwrap.ValidatingRuntime`
- `agentwrap.ObservingRuntime`
- `agentwrap.RunStore` or an UltraPlan implementation of that interface

### 11.2 Required Runtime Composition

The default production composition must be:

```text
UltraPlan task runner
  -> agentwrap.ObservingRuntime
  -> agentwrap.ValidatingRuntime
  -> agentwrap.PolicyRunner
  -> opencode.Runtime
```

Rationale and requirements:

- `ObservingRuntime` should wrap the logical run so UltraPlan can inspect active/completed runs, ordered events, final merged metadata, sink failures, and store-backed records.
- `ValidatingRuntime` must be used so runtime exit success is insufficient without UltraPlan artifact validation.
- `PolicyRunner` must be used for bounded retry, wait, rate-limit handling, and fallback decisions.
- `opencode.Runtime` must remain the adapter that owns OpenCode CLI invocation and native event projection.
- UltraPlan may adjust wrapper order only for an explicit reason, such as allowing validation failures to participate in policy fallback for a specific task type.

### 11.3 Run Request Mapping

For each UltraPlan task, build an `agentwrap.RunRequest`:

- `Prompt`: composed UltraPlan prompt.
- `WorkDir`: source directory, study root, or workspace root according to task type.
- `Provider` and `Model`: resolved model configuration.
- `Timeout`: task timeout.
- `Metadata`: study, dimension, source, source kind, output path, task kind, and UltraPlan task ID.
- `RequireHealth`: configured health checks, at least runtime availability, structured output, and working directory for OpenCode-backed tasks.
- `RequireCaps`: structured events and cancellation for normal tasks; permissions, artifacts, usage, or validation events when task policy requires them.
- `PermissionPolicy`: task-specific agent permissions.
- `Sandbox`: configured sandbox mode.
- `WantSession`, `SessionID`, and `SessionAction`: only when UltraPlan intentionally continues or repairs a related run.
- `Validation`: agentwrap validation expectations and custom validators for required UltraPlan artifacts.

### 11.4 OpenCode Adapter Usage

UltraPlan must use `opencode.NewRuntime` and options documented by agentwrap:

- `opencode.WithExecutable`
- `opencode.WithExtraArgs`
- `opencode.WithEnv`
- `opencode.WithStderrLimit`

UltraPlan must rely on the adapter to:

- Launch `opencode run --format json`.
- Pass `--dir`, `--model`, `--session`, extra args, and prompt.
- Decode newline-delimited JSON records.
- Project native OpenCode records into `agentwrap.Event` values.
- Attach native raw payloads as unsafe diagnostics by default.
- Classify malformed output, runtime exit, timeout, cancellation, rate limits, permission issues, and health failures into `agentwrap.SDKError`.
- Translate supported `agentwrap.PermissionPolicy` fields into `OPENCODE_CONFIG_CONTENT`.
- Enforce required health checks before process launch.
- Report best-effort session continuation metadata for supported session actions.

UltraPlan must not parse OpenCode stdout/stderr directly except through agentwrap diagnostics already exposed in `RunResult`, `SDKError`, events, or metadata.

## 12. Canonical Events

### 12.1 Event Envelope

UltraPlan must consume agentwrap canonical events. The source event shape is:

```go
type Event struct {
    ID        agentwrap.EventID
    RunID     agentwrap.RunID
    SessionID agentwrap.SessionID
    Time      time.Time
    Type      string
    Payload   agentwrap.EventPayload
    Raw       *agentwrap.RawPayload
}
```

`Event.Type` preserves the native or adapter-defined event name. UltraPlan should use `Event.Kind()` or `Payload["event_kind"]` for canonical classification.

### 12.2 Required Event Kinds

UltraPlan must handle agentwrap's current event kinds:

- `lifecycle`
- `session`
- `message`
- `progress`
- `tool`
- `artifact`
- `permission`
- `blocking`
- `usage`
- `warning`
- `fatal_error`
- `rate_limit`
- `validation`
- `retry`
- `fallback`
- `final_result`
- `native_extension`

UltraPlan task state and logs should map these event kinds into product-facing status, diagnostics, and run history without losing the original event kind.

### 12.3 Event Compatibility

- Unknown native events projected as `native_extension` must not fail tasks by themselves.
- Unsafe raw native payload bytes must not be persisted by UltraPlan unless an explicit debug retention option is enabled.
- Persisted event records should preserve raw payload presence, source, encoding, safety, and omission reason when bytes are omitted.
- UltraPlan must implement or configure an agentwrap `EventSink` and `RunStore` that correlate agentwrap records with the UltraPlan operational run and support durable local inspection.
- Sprint 35 product-facing event records have a monotonically increasing sequence within each UltraPlan run. They are sanitized and durably appended before live fan-out.
- Replay uses an explicit cursor. Duplicate delivery is safe; an unavailable cursor produces a typed retention-gap result plus the current durable run snapshot rather than silently beginning a partial stream.
- Persistence, compaction, sampling, and subscriber backpressure failures are observable and cannot block the underlying runtime indefinitely.

### 12.4 Agentwrap Metadata

UltraPlan must consume and persist relevant `agentwrap.RunMetadata` fields:

- Runtime context.
- Parent run ID.
- Attempts.
- Policy decisions.
- Status and timing.
- Session metadata.
- Permission metadata.
- Cleanup metadata.
- Validation and repair metadata.
- Artifact references.
- Warnings and errors.
- Usage, estimated cost, and throughput where available.
- Native metadata safe for persistence.

Unknown usage token values must remain unknown. UltraPlan must not convert unknown usage to zero.

## 13. Scheduler and Run State

### 13.1 Task State

```go
type TaskState struct {
    ID              string
    Kind            TaskKind
    Study           string
    DimensionNumber string
    DimensionSlug   string
    Source          string
    SourceKind      SourceKind
    OutputPath      string
    Status          TaskStatus
    Attempts        int
    LastError       *UltraPlanError
    LastAttemptAt   *time.Time
    NextRetryAt     *time.Time
    StartedAt       *time.Time
    CompletedAt     *time.Time
    Validation      *ValidationResult
    AgentRunID      agentwrap.RunID
    SessionID       agentwrap.SessionID
    TurnID          agentwrap.TurnID
}
```

Task kinds:

- Analysis.
- Synthesis.

### 13.2 Run State File

Each stateful study run must persist to:

```text
studies/<study>/.ultraplan/run-state.json
```

State file fields:

- Schema version.
- Run ID.
- Created time.
- Updated time.
- Batch size.
- Filters.
- Config summary.
- Task list.
- Agentwrap run metadata summary for each started task.
- Agentwrap validation, policy, permission, session, cleanup, artifact, and usage summaries where available.
- Completion flag.

### 13.3 Atomic Writes

State writes must be atomic:

1. Write to temporary file in same directory.
2. Flush and close.
3. Rename over previous state file.

If atomic rename is unavailable on a platform, the implementation must use the safest available strategy and document residual risk.

### 13.4 Resume Logic

On resume:

- Load state file.
- Validate schema version.
- Reset stale running tasks to pending or failed according to policy.
- Revalidate completed task outputs.
- Queue missing synthesis tasks when all inputs are valid.
- Exclude inapplicable Markdown source/dimension pairs from required input checks.
- Reconcile task state with agentwrap completed run records when an `agentwrap.RunStore` is available.
- Preserve attempt history.
- Continue respecting filters from the original run unless user explicitly overrides them.

### 13.5 Scheduling

- Use bounded worker pools.
- Never spawn more runtime processes than configured.
- Build the task matrix from applicable source/dimension pairs only.
- Prefer synthesis tasks when they become unblocked, but do not starve analysis tasks.
- Stop scheduling new work on cancellation.
- Save state after each meaningful transition.
- Drain `agentwrap.Run.Events()` while waiting so runtime event channels cannot block.
- Always call `Run.Wait` to obtain final `agentwrap.RunResult` and classified terminal errors.

## 14. Retry, Fallback, and Backoff

### 14.1 Agentwrap Policy Requirements

UltraPlan must use `agentwrap.PolicyRunner` and `agentwrap.BasicPolicy` for default retry, wait, rate-limit, and fallback behavior.

Requirements:

- Configure `MaxAttemptsPerTarget` from UltraPlan execution config.
- Configure `RetryRateLimits` according to workspace policy.
- Use agentwrap backoff implementations such as `agentwrap.ExponentialBackoff` or a compatible configured policy.
- Configure fallback alternatives as `agentwrap.FallbackAlternative` values.
- Preserve original prompt, workdir, timeout, permission policy, sandbox, metadata, and validation expectations when constructing fallback requests unless a specific fallback intentionally overrides them.
- Record policy decisions from `RunMetadata.Policy.Decisions` in UltraPlan task state.
- Surface agentwrap `rate_limit`, `retry`, and `fallback` events in status/log output.
- Do not implement retry/fallback logic inside the OpenCode adapter.

### 14.2 Error Classifications

UltraPlan product errors may include workspace, filesystem, source, dimension, and ambiguity categories. Runtime-facing errors must use or wrap `agentwrap.SDKError` categories.

Agentwrap categories that UltraPlan must handle:

- `configuration`
- `health`
- `runtime_unavailable`
- `provider_unavailable`
- `model_unavailable`
- `authentication`
- `permission`
- `rate_limit`
- `timeout`
- `cancellation`
- `malformed_event`
- `runtime_exit`
- `validation`
- `repair_exhausted`
- `cleanup`
- `unknown`

UltraPlan-specific categories:

- Usage.
- Workspace.
- Filesystem.
- SourceNotFound.
- DimensionNotFound.
- AmbiguousReference.
- MissingOutput.

### 14.3 Error Attributes

UltraPlan must inspect runtime errors with `errors.As` or `agentwrap.ErrorAs`, not by parsing error strings.

For runtime errors, persist or display safe `agentwrap.SDKError` fields:

- Category.
- Operation.
- UserDetail.
- DebugDetail when debug output is enabled.
- Provider.
- Model.
- RuntimeKind.
- ExitCode.
- Signal.
- NativeType.
- RetryAfter.
- Redacted metadata.

### 14.4 Backoff State

Backoff state must be derived from agentwrap policy attempts and persisted in UltraPlan task state:

- Attempt count.
- Target index.
- Attempt number on current task.
- Last classified error.
- Next retry time when known.
- Fallback model/runtime when selected.

### 14.5 Fallback Policy

Fallback may support:

- Backup model through an agentwrap fallback request.
- Backup provider through an agentwrap fallback request.
- Backup runtime in future versions through a different `agentwrap.Runtime`.

Fallback decisions must be visible in task state, agentwrap metadata, and events.
## 15. Validation Requirements

### 15.1 Agentwrap Validation

UltraPlan must use `agentwrap.ValidatingRuntime` and `agentwrap.ValidationSpec` for runtime output validation.

Agentwrap expectation kinds available to UltraPlan:

- `file`
- `directory`
- `artifact`
- `markdown_template`
- `json`
- `metadata`
- `custom`

UltraPlan-specific validation rules must be implemented as `agentwrap.Validator` or `agentwrap.ValidatorFunc` checks where they validate runtime-produced artifacts.

### 15.2 Validation Result Mapping

UltraPlan task state must map from `agentwrap.ValidationResult`:

- `Passed`
- `Skipped`
- `PassedCount`
- `FailedCount`
- `SkippedCount`
- `Checks`
- `Failures`
- `Errors`
- safe native details

Required validation failures must make the logical run fail through agentwrap. Optional validation failures may remain warnings.

### 15.3 Report Validators

Per-source report validator checks:

- File exists.
- File is non-empty.
- Heading exists.
- Source info section exists.
- Summary section exists.
- Rating section exists.
- Rating can be parsed.
- Required question/answer section exists.
- Citation shape appears where required.

Directory source report validation:

- Must enforce code citation shape unless the dimension explicitly disables that requirement.
- Code citations should use file paths and line numbers.

Markdown document source report validation:

- Must not require code citations by default.
- Must still require a rating, summary, and answers to dimension questions.
- May validate document-specific citations or section references if a dimension explicitly requires them.
- Must record the source kind in validation diagnostics.

Final report validator checks:

- File exists.
- File is non-empty.
- Study parameters exist.
- Sources studied table exists.
- Executive summary exists.
- Rating summary exists.
- Pattern or synthesis sections exist.
- Open questions or notable absences section exists.

Planning reasoning validator checks:

- File exists.
- Requirement mapping exists.
- Decisions and tradeoffs exist.
- Risks and assumptions exist.
- Exit criteria exist.

Planning plan validator checks:

- File exists.
- Cites `reasoning.md`.
- Decisions to execute exist.
- Task checklist exists.
- Evidence checklist exists.
- Risks and blockers exist.
- Success criteria exist.
- Phase 2 plan validation does not invoke smoke, review, issues, or Git mutation.

### 15.4 Validation Failure Behavior

Validation failure must:

- Mark the task failed or repairing according to agentwrap validation and repair policy.
- Include failed check names.
- Include expected vs observed details where safe.
- Include path to invalid artifact.
- Be visible in status output.
- Be recorded in `RunMetadata.Validation`.
- Emit agentwrap validation events.

### 15.5 Repair Requirements

UltraPlan may enable bounded repair through `agentwrap.RepairConfig`.

Requirements:

- Repair attempts must be bounded by `MaxAttempts`.
- Repair prompts must be built from safe expected/observed facts and artifact references, not large raw report content by default.
- Repair should use `SessionActionContinue` when same-session repair is useful and supported.
- Repair must inherit original workdir, provider/model, sandbox, permission mode, and permission policy unless UltraPlan explicitly overrides them.
- Exhausted repair must surface `agentwrap.ErrorRepairExhausted`.
- Permission denial during repair remains `agentwrap.ErrorPermission`.

## 16. Code Reference Extraction

### 16.1 Citation Syntax

Supported citation forms:

- `` `path/to/file.go:42` ``
- `` `path/to/file.go:42-58` ``
- `` `path/to/file.go:42,47,53` ``

The extractor should also tolerate an en dash in old reports but should emit normalized output.

### 16.2 Source Table Parsing

Reports must contain a sources table with source names and paths. The extractor must parse rows shaped like:

```markdown
| 1 | source-name | `sources/source-name` |
```

The parser must support workspace-relative and report-relative source paths.

### 16.3 Resolution Algorithm

For each citation:

1. Try path relative to each source root.
2. If citation includes a source-prefixed path, try stripping first path segment.
3. If unresolved, search by basename within source roots, excluding `.git`, `node_modules`, and known ignored directories.
4. Record unresolved reference if no match exists.

### 16.4 Output

Output must include:

- Source name.
- File path.
- Requested line spec.
- Resolved absolute or workspace-relative path.
- Rendered code with line numbers.
- Unresolved reference summary.

JSON output must expose structured refs and resolution status.

## 17. Summary Generation

Summary generator must:

- Discover dimensions in deterministic order.
- Discover sources in deterministic order.
- Parse ratings from per-source reports.
- Write `summary.csv`.
- Include columns for source, each dimension, and total.
- Sort by total descending.
- Represent missing ratings as empty cells.
- Represent inapplicable Markdown source/dimension pairs distinctly from missing expected reports, either as an empty non-error cell or a documented sentinel value.
- Exclude inapplicable pairs from missing-report warnings.

Rating parser must handle:

- `**8 / 10**`
- `8/10`
- `Rating: 8`

Ambiguous ratings should create warnings, not invented values.

## 18. Project, Sprint Planning, Execute, Review, and Smoke Technical Requirements

Phase 2 implements the governed planning and execute side of UltraPlan through `execute`. Phase 3 extends the same product-owned sprint flow through review and smoke:

```text
study -> select -> distill -> reason -> plan -> execute -> review -> smoke
```

The TypeScript prototype includes planning, execution, smoke, review, and issue-tracking behavior. UltraPlan Go Phase 2 ports the planning artifact chain and adds controlled implementation execution. Phase 3 replaces the prototype's manually coordinated review and smoke with product-owned stages while keeping detailed smoke runs and issue evidence in the external harness. General-purpose issue tracking and Git mutation remain deferred.

### 18.1 Package Ownership

Planning behavior must be owned by project and sprint modules:

```text
internal/project
internal/sprint
```

`internal/project` owns:

- project root discovery under `projects/<project>`
- project docs discovery
- roadmap discovery
- `project-index.md` parsing and validation
- catalog entries for contracts, evidence reports, reasoning templates, review protocols, and project source documents
- project status output

`internal/sprint` owns:

- sprint root discovery under `projects/<project>/sprints/<slug>`
- planning-stage domain model
- planning artifact path rules
- `flow-state.json`
- planning-stage prompt rendering
- planning-stage validation
- execute-stage prompt rendering
- execute task extraction, validation, and run-state persistence
- global and per-stage runtime model resolution for sprint planning and execute stages
- sprint status output
- flow execution through `execute`
- review scope, prompt rendering, independent reviewer orchestration, structured result validation, deterministic verdict synthesis, and `review.md`
- smoke harness discovery, review gating, scope selection, safe invocation, evidence-link validation, verdict synthesis, and `smoke.md`
- flow execution through `review` and `smoke`

The dependency direction is:

```text
sprint -> project
sprint -> workspace
sprint -> platform/runtime
sprint -> platform/process
project -> workspace
platform/* -> no product modules
study -> no project or sprint modules
```

Study outputs may be referenced by project indexes and sprint artifacts as evidence paths, but `project` and `sprint` must not depend on study services, source/dimension models, report validators, run-loop scheduling, summary generation, or rating parsing.

### 18.2 Reuse Boundary

Phase 2 may reuse generic infrastructure from Phase 1:

- workspace discovery and path safety
- config precedence and redaction
- command exit code and output conventions
- generic runtime prompt execution
- filesystem helpers that are genuinely cross-module, such as atomic writes
- small validation/result structs when they do not carry study semantics

Phase 2 must not prematurely abstract study behavior into shared packages. Keep these local to `study`:

- study `Service`
- source and dimension models
- study prompt builders
- report validation
- rating parsing
- summary generation
- task scheduling and run-loop state

If sprint and study both need the same mechanical file operation, extract the file operation. If they both happen to validate Markdown headings, implement locally first and extract only after the second concrete use proves the shared behavior is stable.

### 18.3 Project Model

```go
type Project struct {
    Name         string
    Root         string
    DocsDir      string
    RoadmapPath  string
    IndexPath    string
    SprintsDir   string
}

type ProjectIndex struct {
    SourceDocuments    []CatalogEntry
    Contracts          []CatalogEntry
    EvidenceReports    []CatalogEntry
    ReasoningTemplates []CatalogEntry
    ReviewProtocols    []CatalogEntry
}

type CatalogEntry struct {
    Name        string
    Path        string
    Summary     string
    AppliesTo   string
}
```

Requirements:

- Project names must be filesystem-safe.
- Project roots must resolve under `projects/`.
- Project docs must be Markdown files under `projects/<project>/docs/`.
- `project-index.md` is a catalog, not a sprint plan.
- Catalog paths must resolve within the workspace unless explicitly marked external.
- Missing catalog entries must produce actionable diagnostics.

### 18.4 Sprint Planning Model

```go
type PlanningStage string

const (
    StageRequirements      PlanningStage = "requirements"
    StageSprintIndex       PlanningStage = "sprint-index"
    StageTechnicalHandbook PlanningStage = "technical-handbook"
    StageAreaReasoning     PlanningStage = "area-reasoning"
    StageReasoning         PlanningStage = "reasoning"
    StagePlan              PlanningStage = "plan"
    StageExecute           PlanningStage = "execute"
)

type Sprint struct {
    Project string
    Slug    string
    Root    string
    Stages  []StageState
}

type StageState struct {
    Stage     PlanningStage
    Status    StageStatus
    Path      string
    LastRunAt string
    Error     string
}
```

Supported sprint artifacts through Phase 3:

```text
requirements.md
sprint-index.md
technical-handbook.md
reasoning/
reasoning.md
plan.md
execute.md
flow-state.json
.run-state.json
review.md
smoke.md
```

Deferred artifacts:

```text
issues.md
issues.json
```

Detailed smoke JSON and issue records belong to the external harness `runs/` and `issues/` directories rather than the sprint root. Historical manual `review.md` and `deep-smoke.md` files may exist until migration; Phase 3 review and smoke atomically replace the canonical root `review.md` and `smoke.md` when explicitly run.

### 18.5 Flow State

`flow-state.json` must be versioned and stored in the sprint root.

Required fields:

- schema version
- project name
- sprint slug
- updated timestamp
- per-stage status
- per-stage artifact path
- per-stage last run timestamp
- per-stage error, if any

Allowed statuses:

- missing
- ready
- complete
- failed
- skipped

Requirements:

- Writes must be atomic.
- Loading must reject malformed JSON and unsupported schema versions with clear diagnostics.
- Existing artifacts must be inspected to initialize or refresh flow state.
- `area-reasoning` is skipped only when `sprint-index.md` selects no reasoning templates.
- Execute may appear in Phase 2 flow state only after valid prerequisites through `plan`.
- No smoke, review, or issue stages may appear in Phase 2 flow state.
- Phase 3 flow state adds `review` after `execute` and `smoke` after `review`.
- Review and smoke stage state must distinguish execution status from verdict so a successfully completed investigation may truthfully report a failing verdict.
- Review and smoke state must record the governed input fingerprint used by the current artifact. A fingerprint mismatch makes the artifact stale.

### 18.5.1 Execute Run State

`.run-state.json` must be versioned and stored in the sprint root when execute begins.

Required fields:

- schema version
- project name
- sprint slug
- target repository path or workspace-relative target reference
- source `plan.md` path and content fingerprint where practical
- updated timestamp
- deterministic task records
- per-task status
- per-task attempts
- per-task timestamps
- per-task safe diagnostics
- per-task runtime metadata summaries where available

Allowed task statuses:

- pending
- running
- complete
- failed
- cancelled

Requirements:

- Writes must be atomic.
- Loading must reject malformed JSON and unsupported schema versions with clear diagnostics.
- Interrupted runs must be resumable without redoing complete tasks unless explicitly forced.
- Stale running tasks must be recovered to a retryable or failed state with a diagnostic.
- Task IDs must be stable for the same validated `plan.md`.
- Execute must not mutate Git state automatically.
- Execute must not run smoke, review, or issue workflows.
- Target repository writes must be constrained to the configured target implementation directory unless explicitly configured otherwise.

### 18.6 Stage Validators

Planning validators are product behavior owned by `internal/sprint`.

`requirements.md` validator checks:

- file exists
- no placeholders
- required sections exist: sprint goal, required outputs, acceptance criteria, non-goals

`sprint-index.md` validator checks:

- file exists
- required sections exist: sprint scope, selected contracts, selected evidence reports, selected reasoning templates, required review protocols
- all selected catalog paths are present in `project-index.md`
- excluded context is explicit when relevant

`technical-handbook.md` validator checks:

- file exists
- selected evidence reports are present and readable
- required sections exist: selected studies/reports, relevant patterns, trade-offs
- document does not make implementation decisions

`reasoning/*.md` validator checks:

- files exist for selected reasoning templates
- required sections exist: area decisions and trade-offs
- no placeholders

`reasoning.md` validator checks:

- file exists
- required sections exist: final decisions, expected evidence, assumptions and risks
- required area reasoning files exist when selected
- no placeholders

`plan.md` validator checks:

- file exists
- cites `reasoning.md`
- required sections exist: decisions to execute, tasks, evidence checklist
- tasks trace to reasoning decisions
- no placeholders
- does not invoke smoke, review, issue tracking, or Git mutation as Phase 2 CLI behavior
- executable tasks are explicit enough for the execute stage to assign deterministic task IDs

`execute.md` validator checks:

- file exists after execute produces a summary
- cites `plan.md`
- summarizes task counts and terminal states
- records failed tasks with actionable diagnostics or pointers to `.run-state.json`
- does not claim smoke, review, issue, or Git mutation completion

`review.md` validator checks:

- file exists after review completes and contains no placeholders
- identifies project, sprint, target implementation, input fingerprint, review date, selected contracts, and selected protocols
- includes decision conformance, plan execution, verification evidence, contract conformance, handbook conformance, applicability/deferred scope, findings, deviations, and final assessment
- includes one valid result for every selected contract plus the technical handbook
- uses repository/workspace-relative contained evidence paths and valid line references
- records missing/failed reviewer tasks rather than treating them as pass
- uses only `pass`, `pass_with_findings`, `fail`, or `blocked` as the review verdict
- applies deterministic severity rules: applicable blocker/high findings cannot produce pass or pass-with-findings
- does not claim smoke execution or modify other sprint artifacts

`smoke.md` validator checks:

- file exists after smoke completes and contains no placeholders
- identifies project, sprint, review verdict/fingerprint, selected harness, selected scope, selection rationale, environment/runtime/model, and smoke date
- includes run ID, safe argv display, result counts, verdict, external summary/detail paths, relevant issue references, and required next action
- verifies referenced harness evidence exists under the cataloged harness and matches the reported run ID
- uses only `pass`, `pass_with_open_issues`, `fail`, `blocked`, or `not_applicable` as the smoke verdict
- never treats missing prerequisites or missing required coverage as pass
- does not copy raw event streams, secrets, or unrestricted stdout/stderr into the sprint artifact

### 18.7 Prompt Rendering and Runtime Execution

Sprint planning prompts and output templates are product assets shipped as embedded defaults. The sprint module owns their semantics; the workspace package owns safe default lookup and opt-in materialization through `ultraplan defaults install`. Files at matching workspace `prompts/` or `templates/` paths are optional overrides, not prerequisites. Prompt rendering must support:

- project name substitution
- sprint slug substitution
- sprint path substitution
- workspace-relative paths
- dry-run preview
- optional output file for prompt preview

Runtime execution for planning stages uses the generic `platform/runtime` request model. The runtime request may include prompt, working directory, model, variant, timeout, permissions, and expected output path, but it must not include project or sprint semantics in the platform package.

Sprint runtime model resolution must be product-owned and deterministic:

1. explicit command override for the requested stage, when supported
2. configured stage-specific model for the requested stage
3. configured global sprint/planning model, when present
4. `models.primary`
5. `models.default`

Supported stage-specific keys are `sprint-index`, `technical-handbook`, `area-reasoning`, `reasoning`, `plan`, `execute`, and `review`. Smoke model/runtime selection belongs to the external harness request and its explicit configuration rather than an agentwrap review model key. Validation must reject unknown stage keys and empty model values. Diagnostics and prompt previews must show the selected model source without leaking secrets.

Runtime success is insufficient. A planning stage is complete only when:

- runtime execution succeeds, when runtime is invoked
- expected artifact exists
- expected artifact passes its stage validator
- `flow-state.json` is updated atomically

Runtime success is insufficient for execute tasks. An execute task is complete only when:

- the task was extracted from a valid `plan.md`
- runtime execution succeeds, when runtime is invoked
- expected task evidence is present or the task records an explicit diagnostic explaining why evidence cannot be machine-validated
- `.run-state.json` is updated atomically
- task completion is reflected in sprint status

### 18.8 Commands

Required commands through Phase 3:

```bash
ultraplan project list
ultraplan project <project> status
ultraplan project <project> validate
ultraplan sprint <project> <sprint> status
ultraplan sprint <project> <sprint> validate [stage]
ultraplan sprint <project> <sprint> prompt <stage>
ultraplan sprint <project> <sprint> flow --to <stage>
ultraplan sprint <project> <sprint> execute [--task <id>]
ultraplan sprint <project> <sprint> review
ultraplan sprint <project> <sprint> smoke
ultraplan sprint <project> <sprint> verify [--to review|smoke]
```

Flow options:

- `--from <stage>`
- `--to <stage>`
- `--force`
- `--no-skip`
- `--dry-run`
- model, variant, and timeout overrides where runtime execution is available
- stage-specific model overrides for sprint planning, execute, and review runtime requests
- bounded review parallelism and review resume controls
- explicit smoke level/suite/test selection, smoke timeout, and review-failure diagnostic override
- stable `--json` output for review, smoke, verify, and status

Phase 2 valid `--to` stages end at `execute`. Phase 3 adds `review` and then `smoke`. `issues` remains invalid.

Sprints 36–38 extend the verification surface without adding QA or repair to the planning-stage enum:

```bash
ultraplan sprint <project> <sprint> conformance-review
ultraplan sprint <project> <sprint> qa [--dry-run|--restart|--shard <id>|--suite smoke] [--json]
ultraplan sprint <project> <sprint> repair --issue <id>
ultraplan sprint <project> <sprint> verify --to conformance-review|qa [--repair --max-cycles <n>]
```

`review` remains a compatibility alias/projection for Conformance Review. `smoke` remains compatible and may later map to `qa --suite smoke` only after parity. General-purpose `issues` remains invalid; QA issue records are bounded evidence artifacts, not a project-management surface.

### 18.9 Phase 3 Review And Deep Smoke Requirements

#### 18.9.1 Review Scope And Execution

Review preflight must:

- require valid governed inputs through `execute`
- resolve every selected contract and review protocol from `project-index.md`
- reject missing, duplicate, unknown, unreadable, or escaping catalog paths
- compute a deterministic fingerprint over requirements, sprint index, technical handbook, area reasoning, final reasoning, plan, execute state/summary, selected contracts/protocols, target implementation identity, and explicit changed-path scope
- show the selected model source, reviewer count, concurrency, target scope, and permitted writes in dry-run/TUI confirmation

UltraPlan, not the model, owns reviewer fan-out. It must issue one independent structured agentwrap request per selected contract plus one handbook request, with bounded concurrency, context cancellation, read-only permissions, and validated JSON results. A reviewer must classify each requirement/guidance item as direct, partial, not triggered, or explicitly deferred before judging conformance.

Deterministic product checks must cover decision conformance, plan-task execution evidence, approved verification-command results, citation containment/line validity, complete reviewer coverage, deviations, and missing evidence. Runtime exit success alone cannot produce a passing review.

The final review verdict is computed by product code. The final Markdown is written atomically to the sprint-root `review.md`, replacing the prior manual or generated file only when the command has permission to run. A failed/incomplete new review must not corrupt the last complete review artifact; current failed state and diagnostics remain visible through flow state and status.

#### 18.9.2 Smoke Harness Contract And Execution

`project-index.md` must catalog the smoke harness root and versioned manifest. The manifest must describe an executable plus argument prefix, protocol version, machine-readable discovery/run commands, bounded authoring paths, evidence directories, and supported capability flags. UltraPlan must never execute a shell command parsed from Markdown or README prose.

Before discovery on every non-dry smoke run, UltraPlan must invoke the configured smoke author model in the external harness. The author receives the governed sprint inputs, execute/review evidence, target implementation, deterministic test inventory, existing harness, and manifest-declared writable paths. It must create or update a durable sprint-specific suite for real boundaries that deterministic unit/integration tests cannot prove. Those boundaries include real provider/model behavior where the product path uses a provider, process and signal lifecycle, real filesystem/configuration, network, credentials, browser engines, timing, cancellation, and platform behavior. It must not duplicate ordinary deterministic verification or introduce a provider call into an otherwise offline product path.

Harness discovery must return structured levels, suites, enumerated tests, per-test coverage IDs, sprint mappings with required coverage IDs, runtime/network/credential prerequisites, expected duration/cost class, and evidence schema version. A complete mapping with empty suites, empty tests, missing required coverage, or unassigned coverage must not pass. Smoke selection prefers the authored sprint-specific suite, then directly mapped suites, then explicit tests for investigation. A narrow passing rerun does not replace required evidence from its containing suite.

Smoke execution must:

- require a current passing or non-blocking review by default
- record the smoke author run/model and every changed harness path
- show scope, prerequisites, model/runtime, duration/cost class, allowed mutation roots, and external evidence destination before execution
- use explicit argv, contained cwd, bounded environment forwarding, timeout, context cancellation, and descendant-process cleanup
- validate the harness run ID, exit/result counts, evidence paths, evidence identity, and relevant open/resolved issue references
- keep raw run JSON, stdout/stderr, per-test artifacts, and issue files in the harness
- atomically write the linked human-readable sprint-root `smoke.md` after validating the authoring scope and enumerated executed test identities
- classify unavailable required environment as blocked and irrelevant scope as not applicable

The product review/smoke workflows must not edit product source, product tests, governed sprint inputs, or Git state. Smoke authoring may edit only manifest-declared harness authoring paths; run and issue evidence remains confined to the manifest-declared evidence roots. Any other harness change or any product/governed-input identity change fails smoke.

#### 18.9.3 Flow, Freshness, And TUI Parity

The default order is `execute -> review -> smoke`. Blocking/high review findings stop smoke unless an explicit diagnostic override is confirmed. Flow state records execution status, verdict, artifact path, timestamp, input fingerprint, run/evidence ID where applicable, and safe diagnostics for both stages.

Any governed input, selected contract/protocol, execute evidence, or implementation fingerprint change makes the prior review stale. A stale review also makes smoke stale. Status, JSON, and TUI must agree on readiness, running/completed/failed/cancelled state, verdict, staleness, artifact, evidence link, and required next action.

Every Phase 3 delivery sprint must update the TUI in the same sprint. TUI operations use the shared app use cases and expose the same review/smoke scope, dry-run, confirmation, progress, cancellation, results, focused reruns, evidence links, and recovery behavior as the CLI. TUI code must not invoke CLI handlers, parse terminal output, interpret native provider payloads, or persist an alternate verification state.

Normal tests use fake review runtimes and a fake smoke harness. Required tests cover success, findings, blocking failure, missing reviewer result, malformed structured output, stale inputs, missing harness, missing coverage, blocked environment, not-applicable scope, timeout, cancellation, malformed harness output, missing evidence, path escape, redaction, CLI/JSON/TUI agreement, and recovery.

### 18.10 Deferred Technical Requirements

The following remain explicitly deferred:

- general-purpose issue tracking and remote issue synchronization
- automatic product/test fixes from findings
- Git add/commit/push
- cross-sprint review/smoke scheduler
- hosted review or smoke service

## 18A. Phase 4 Local Web Surface Technical Requirements

Phase 4 is an interface expansion, not a new product workflow. `internal/web` owns HTTP routing, transport validation, HTML/template rendering, embedded static assets, browser-session protection, and SSE connections. It must not own study/project/sprint state machines, prompt construction, runtime or smoke invocation, verdict computation, artifact persistence, durable run identity, or durable recovery.

The required dependency direction is:

```text
cmd/ultraplan -> internal/app + internal/tui + internal/web
internal/web  -> internal/app
internal/web  -> no product modules directly
internal/app  -> product modules + platform modules
```

Interface runners and other side-effectful surface dependencies should be constructed explicitly at the composition root and passed inward. Phase 4 must not add another package-global mutable runner registration.

Required initial route capabilities:

```text
GET    /                         browser dashboard
GET    /projects/...             browser project and sprint pages
GET    /studies/...              browser study pages
GET    /api/v1/dashboard         structured dashboard state
GET    /api/v1/artifacts/{ref}   bounded allowlisted preview
POST   /api/v1/validations       read-only validation
POST   /api/v1/operations/prepare guarded-operation confirmation
POST   /api/v1/operations        start confirmed operation
GET    /api/v1/operations/{id}   compatibility projection to stable run status/result
GET    /api/v1/operations/{id}/events compatibility projection to replayable progress
DELETE /api/v1/operations/{id}   compatibility cancellation request
GET    /api/v1/runs              workspace-wide durable run projection
GET    /api/v1/runs/{id}         durable run status/result/recovery
GET    /api/v1/runs/{id}/events  cursor-based replay and live progress
DELETE /api/v1/runs/{id}         authorized idempotent cancellation request
```

Exact resource detail routes may evolve, but `/api/v1` JSON and error envelopes are compatibility-controlled once documented. Unknown `/api/` routes return structured JSON errors and must never fall through to an HTML page.

The server delivery hub is ephemeral and bounded. It may hold normalized requests, current connections, and delivery buffers, but not the only run IDs, cancellation routes, history, or terminal results. Server restart recovery reads the durable run projection and package-owned product state. Slow or disconnected SSE subscribers must not block product execution. Compatibility operation URLs must resolve to the durable run, a retained tombstone, or a typed recovery response; session expiry or in-memory reaping alone cannot produce an unexplained 404 for a valid run.

Graceful server shutdown owns cancellation of every worker the server owns under the selected Sprint 35 topology. Shutdown must enter draining, reject new mutations and queued starts, request canonical cancellation exactly once with reason `server_shutdown`, propagate through retries/runtime/process trees, wait for bounded cleanup, and persist one truthful terminal outcome. Valid outcomes are clean cancellation, interruption, cleanup uncertainty, or a failure/completion that was already authoritative before cancellation won. HTTP/SSE closes only after durable outcome or bounded uncertainty recording. Browser navigation, tab close, refresh, or SSE disconnect never cancels a run. Startup and periodic reconciliation must handle stale active runs and locks after forced termination without inferring success from process absence or artifact presence.

Commands and streams remain separate:

```text
browser -> HTTP POST/DELETE -> app command/cancellation
browser <- SSE GET          <- safe bounded progress events
browser -> HTTP GET         -> refreshed durable state
```

Product-owned mutation locks remain mandatory. Existing study locks stay in `internal/study`; sprint mutation exclusion belongs in `internal/sprint`, not in HTTP middleware. Conflicting work must return an actionable conflict rather than execute concurrently by accident.

The first Phase 4 UI uses embedded `html/template` files, CSS, and minimal JavaScript. No frontend framework or JavaScript build system is required. A later client-side framework is allowed only after demonstrated interaction complexity and must still consume the same versioned HTTP/app boundary.

Sprint 32 must converge on this presentation structure:

```text
internal/web/
  templates/
    primitives/
    components/
    layouts/
    pages/
  static/
    css/
      tokens.css
      base.css
      primitives.css
      components.css
      layouts.css
      utilities.css
    js/
      app.js
      operations.js
      sse.js
```

File names may evolve with the actual component set, but the ownership and dependency layers are normative. Template parsing must fail at server startup on duplicate/invalid definitions. Focused tests must render representative primitives, components, layouts, and pages; validate complete no-JavaScript output; assert escaping of hostile values; and verify embedded asset paths, accessibility semantics, and absence of inline or third-party executable content.

## 18B. Delivered Grounded-Planning Technical Requirements

`internal/sprint` owns the `code-context` stage, including ordered-stage membership, artifact paths, prerequisites, validation, runtime execution, flow transitions, state compatibility, and downstream prompt integration. It reuses project-owned target implementation resolution and the generic runtime boundary; neither `internal/web` nor a new repository-index package owns source selection.

The implementation sequence is split deliberately:

1. Sprint 33 adds the vertical slice: stage/state/artifact model, template and prompt, repository read boundary, output restriction, runtime config, validation, CLI/app/web surfaces, compatibility, and deterministic tests.
2. Sprint 34 adds one shared prompt-prefix renderer, injects exact requirements and code-context content into every downstream agent request, materializes the manual skill, completes documentation, and dogfoods a representative real repository.

Required tests include ordered-stage transitions, old flow-state compatibility, path/range validation, source-read/output-write isolation, fake runtime output, missing/invalid output, atomic rerun, model/variant fallback, CLI/JSON/TUI/web parity, exact prefix equality across representative downstream stages, absence of dynamic prefix data, flow execution exactly once, cancellation, and recovery.

The web extensibility gate requires the new stage to expose status, artifact preview, operation start, progress, cancellation, findings, and recovery through the existing application capability model. Adding route-specific product logic for `code-context` fails the gate.

## 18C. Sprint 35 Durable Run Identity And Cross-Surface Observability

The current in-memory operation hub is a compatibility/delivery mechanism, not an adequate execution control plane. Sprint 35 must provide one workspace-wide model shared by CLI, TUI, browser pages, and every local server instance within the explicitly supported topology.

The durable run envelope must include, directly or through stable references:

```text
workspace-scoped run ID
operation kind and project/study/sprint/stage target
accepted / queued / running / terminal lifecycle timestamps
attempt, stage, task, agentwrap run/session, process, and external-harness correlations
owner identity, renewable lease, heartbeat, fencing value, and safe process-birth identity
latest durable event sequence and retention boundary
cancellation request/acknowledgement facts
one immutable arbitrated terminal result or explicit cleanup uncertainty
safe diagnostics, schema version, and migration metadata
```

Acceptance of a runtime-backed execution is atomic with durable run creation and happens before child execution. If required acceptance persistence fails, the child does not start. Event records are redacted before persistence, assigned monotonically increasing per-run sequences, and committed before live fan-out. The store may compact detail, but it must retain an explicit snapshot and lower replay boundary so clients can distinguish complete replay from a gap.

Active-run queries are workspace-wide. Product-page status may filter the common projection, but the global running count cannot be derived from the current page, browser session, or one server's memory. Planning stage readiness and operational execution lifecycle are separate types; UI code must not infer `running` from a stage-status enum that cannot represent it.

Read visibility is independent from originating-session continuity. Mutation still requires fresh authorization and a resolvable current owner. Cancellation is idempotent, durably recorded, routed through the canonical product/runtime cancellation path, and arbitrated against completion, failure, timeout, interruption, shutdown, and reconciliation so only one terminal outcome wins.

Liveness is conservative. A PID alone is insufficient because of PID reuse. Lease expiry, missing process, missing owner, artifact presence, or lost server memory cannot independently prove completion. Reconciliation uses fencing or equivalent stale-writer protection, records its evidence and decisions, and exposes `stalled`, `interrupted`, or `cleanup_uncertain` where truth cannot be proven.

The selected implementation contract is:

- same-host local-filesystem support only; shared-filesystem multi-host access is unsupported
- direct SQLite writers using `modernc.org/sqlite`, WAL, `synchronous=FULL`, foreign keys, a five-second busy timeout, and private `0700`/`0600` storage
- one random owner per process, exact boot/PID/birth identity, five-second heartbeats, 15-second leases, repository fencing generations, and a 45-second reconciliation grace
- no adoption after owner loss; product-owned resume creates deliberate new work
- immutable terminal compare-and-set across completion, failure, timeout, cancellation, interruption, persistence degradation, and reconciliation
- commit-before-delivery events capped at 16 KiB, bounded replay with `(run_id, sequence)` deduplication, explicit gaps, seven-day full history, and thirty-day tombstones
- a 512 MiB hard workspace quota, 496 MiB soft admission threshold, and 16 MiB reserved headroom for active control/terminal evidence
- workspace-visible sanitized reads across sessions, with fresh same-origin session and CSRF authority for browser cancellation
- bounded exporter-neutral metrics, JSON diagnostics, reconciliation evidence, and private redacted support export; OpenTelemetry remains deferred
- fail-closed acceptance/claim and conservative cancellation/reconciliation when later persistence cannot prove success

Every chosen mechanism must pass the same cross-process failure matrix: two CLI runs visible globally; a second supported server replaying and then following a current run; browser-session expiry; observer restart; abrupt owner death; stale lease and PID reuse; duplicate cancellation; completion/cancellation races; slow subscribers; cursor expiry; corrupt/torn records; disk-full or permission failure; bounded retention; and redaction. Tests must prove that operational records never override canonical product artifacts or stage verdicts.

## 18D. Sprint 36 Read-Only QA Decomposition And Synthesis

Sprint 36 must preserve `review` and `review.md` compatibility while presenting the existing analytical capability as Conformance Review. It introduces a separate `VerificationPhase` model and schema-versioned QA state outside detailed `flow-state.json` ownership.

The QA mapper must deterministically derive bounded behavioral verification surfaces from execute evidence, changed paths, requirements, code-context, selected contracts and protocols, Conformance Review findings, adjacent tests, and risk tags. Every changed path has one primary shard; explicit boundary shards and bounded overlap are allowed where behavior crosses packages, interfaces, state transitions, or public APIs.

Read-only investigators may inspect approved context and run safe non-mutating checks. They may formulate falsifiable theories and record confirmation, refutation, inconclusive, invalid, blocked, and cross-shard outcomes, but they may not create tests, mutate production or verification code, promote issues, or repair code. Global synthesis may deduplicate, combine, challenge, and request bounded follow-up without issue promotion.

QA maps, shards, theories, attempts, and synthesis may use deterministic schema-versioned identifiers scoped to the verification model. Those identifiers must have explicit compatibility or migration behavior and must not claim to be the later workspace-wide content identity contract.

Required tests cover deterministic mapping, changed-path coverage, budgets, overlap, cancellation, resume, fingerprint invalidation, atomic state, bounded cross-shard follow-up, and CLI/JSON/TUI/browser agreement.

## 18E. Sprint 37 Evidence-Producing QA And Smoke Integration

Writable investigation requires one validated isolated workspace per shard attempt. Inability to prove isolation, target identity, containment, process cleanup, or path safety is `blocked`. Investigators may create targeted tests, fixtures, probes, smoke scenarios, and bounded experiments only inside that workspace; generated verification patches remain evidence or regression candidates rather than production repair.

A global adjudicator alone validates expectation grounding, implementation freshness, setup validity, containment, repeatability, flakiness, evidence sufficiency, root-cause grouping, issue promotion, repair eligibility, and regression-candidate classification. A failing check is not automatically an issue.

Canonical user-facing QA uses `qa.md`; versioned detailed maps, shards, evidence, synthesis, and issues live under dedicated verification state, while `flow-state.json` retains only canonical summaries, freshness, verdicts, and pointers.

The current smoke protocol may be wrapped as a QA suite/executor only while preserving harness discovery, selection and containing-suite semantics, environment allowlisting, bounded processes, timeout, cancellation, cleanup, evidence validation, diagnostic-only behavior, and canonical-versus-narrow evidence. `smoke` and `smoke.md` compatibility remain until parity is proven.

## 18F. Sprint 38 Manual Repair And Bounded Automatic Repair

Production repair is permitted only from a frozen adjudicated issue packet containing the confirmed claim, current supporting evidence, violated expectations, allowed paths, acceptance criteria, exact reproducer, affected shards, and containing checks. Repair requires explicit confirmation and may not weaken or rewrite evidence, tests, requirements, or acceptance criteria.

Reverification widens progressively from exact reproducer to affected shard, linked theories, neighboring/boundary shards, containing QA suites, and Conformance Review delta. A repair agent cannot declare its own global success.

Automatic repair follows only after manual repair works end to end. It must bound cycles, reopenings, issue-set stagnation, scope growth, severity growth, uncertainty, target drift, cleanup failure, and design decisions. Terminal outcomes distinguish `verified`, `verified_with_findings`, `failed`, `blocked`, `escalated`, and `stalled`.

## 18G. Sprint 39 QA And Repair Dogfooding And Hardening

Sprint 39 must dogfood broad multi-package, boundary, concurrency, cancellation, persistence, recovery, invalid-setup, cross-shard, manual-repair, bounded-automatic-repair, restart, and cleanup-uncertain cases through CLI, JSON, TUI, and browser.

Measure deterministic shard quality, false-positive and inconclusive rates, evidence validity, isolation reliability, investigation cost, issue quality, repair convergence, cancellation/recovery behavior, and operator usability. Automatic repair remains disabled if issue sets do not shrink reliably or severity, scope, or uncertainty grows.

Sprint 39 exits only when read-only and writable boundaries are reliable, adjudication rejects invalid evidence consistently, manual repair succeeds on real work, automatic repair either demonstrates bounded convergence or remains disabled, and representative QA artifacts exist to inform the later content contract.

## 18H. Gated Post-Sprint-39 Technical Direction

The following sequence remains directional until a later sprint passes the preceding gate:

1. **Content identity/provenance pilot:** design optional metadata, semantic block IDs, revision-aware evidence, contradiction, and supersession from real QA theories, evidence, issues, repairs, and findings; preserve legacy Markdown and verification-state compatibility.
2. **Joint content/QA schema dogfood:** validate identifier stability, traceability, authoring burden, and migration before retrieval.
3. **Derived retrieval:** build deterministic semantic records and a measured lexical baseline first. Indexes remain disposable and governed context selection remains authoritative.
4. **Product persistence/SQLite:** classify product state and extract focused package-owned repositories only for proven needs. Sprint 35 operational SQLite does not select authored-artifact authority.
5. **Knowledge graph:** begin, if justified, with a deterministic in-memory read-only projection over explicit relationships and provenance.
6. **Authority/cloud/Aren:** compare filesystem, SQLite/server, and hybrid modes locally before remote authority or typed Aren artifact tools.

Every branch must stop when its evidence gate fails: no content schema based only on speculation, no retrieval before schema dogfood, no product SQLite merely for architectural cleanliness, no graph persistence before useful bounded traversal, and no cloud authority before local authority and recovery semantics are proven.

## 19. Logging and Diagnostics

UltraPlan logging and diagnostics must consume agentwrap observability records rather than directly inspecting native runtime process streams.

### 19.1 Logger Requirements

- Context-aware logging.
- Text and JSON output.
- Log levels: debug, info, warn, error.
- Redaction support.
- Stable fields for structured logs.

Structured log fields:

- Timestamp.
- Level.
- Command.
- Workspace.
- Study.
- Project.
- Sprint.
- Run ID.
- Run schema version.
- Owner and fencing identity when applicable.
- Task ID.
- Stage/execution ID.
- Attempt.
- Runtime.
- Model.
- Event type.
- Message.
- Agentwrap run ID.
- Agentwrap session ID.
- Agentwrap event kind.
- Agentwrap policy decision when applicable.
- Durable event sequence and replay boundary when applicable.
- Lifecycle transition, lease state, reconciliation decision, and terminal-outcome winner when applicable.

### 19.2 Diagnostics

Diagnostics must include:

- Runtime executable path.
- Runtime version when available.
- Working directory.
- Timeout.
- Exit code.
- Safe stderr excerpt.
- Validation failures.
- State file path.
- Agentwrap `RunMetadata` summaries for attempts, policy, validation, repair, sessions, permissions, cleanup, artifacts, usage, warnings, and errors.
- Agentwrap raw payload omission/safety facts when relevant.
- Durable run record path/store identity and schema version.
- Owner, lease, heartbeat age, fencing, and safe process-birth facts.
- Current lifecycle snapshot, last event sequence, oldest retained sequence, and any replay gap.
- Acceptance/append/terminal persistence health, reconciliation backlog, and last decision evidence.
- Cancellation request, routing, acknowledgement, and cleanup certainty.

Diagnostics must not include:

- API keys.
- Full sensitive environment.
- Secret config values.

### 19.3 Agentwrap Stores and Sinks

UltraPlan must use agentwrap observability hooks:

- `agentwrap.ObservingRuntime` for active/completed run projections.
- `agentwrap.EventSink` for event fan-out to logs and local event records.
- `agentwrap.RunStore` for run inspection.

`agentwrap.MemoryRunStore` is suitable for isolated tests only. Sprint 35 production paths use the UltraPlan `internal/runcontrol` SQLite store and correlate safe agentwrap identifiers when observed. Agentwrap remains runtime supervision authority; its store is not substituted for the workspace run-control journal.

Required sink failures must be treated according to agentwrap semantics: returned from `Wait` when the primary runtime outcome succeeded. Best-effort sink failures must be recorded as warnings without replacing the primary outcome.

## 20. Concurrency and Cancellation

### 20.1 Concurrency

- Use `context.Context` for cancellation.
- Use bounded goroutine pools for tasks.
- Avoid unbounded channels.
- Drain `agentwrap.Run.Events()` or deliberately discard through a draining goroutine until the channel closes.
- Ensure `agentwrap.Run.Wait()` is called for every started run.
- Let the agentwrap adapter own subprocess waiting and cleanup.
- Ensure no goroutine leaks in tests.
- Bound HTTP request concurrency, active operations, SSE subscribers, and per-operation event buffers.
- Never let an SSE write or slow browser subscriber block the underlying app operation.
- Serialize or fence competing lifecycle writers so stale owners cannot renew leases, append authoritative transitions, or overwrite the winning terminal outcome.
- Bound durable journal queues, append latency, retained bytes/events, replay work, reconciliation concurrency, and support-bundle size.

### 20.2 Cancellation

Cancellation sources:

- User interrupt.
- Context cancellation.
- Timeout.
- Runtime failure.

Cancellation behavior:

- Stop scheduling new tasks.
- Cancel active task contexts and call `agentwrap.Run.Cancel`.
- Let the agentwrap adapter attempt runtime process termination and cleanup.
- Save state.
- Mark active tasks cancelled or retryable based on policy.
- Return cancellation exit code when user initiated.
- Preserve agentwrap cleanup metadata separately from the primary run result.
- Disconnecting an SSE subscriber cancels only that subscription; explicit operation cancellation uses the operation context and shared app cancellation path.
- Graceful server shutdown stops accepting new work, closes subscribers, cancels workers owned by that server under the selected topology, and leaves durable product and operational state recoverable.
- Shutdown cancellation records `reason: server_shutdown`, is idempotent, and participates in the same single-terminal-outcome arbitration as user cancellation, failure, timeout, and completion.
- The server waits only for a configured bounded cleanup period, does not release locks before ownership is reconciled, and records `interrupted` or `cleanup_uncertain` when complete cleanup cannot be proven.
- Forced termination is reconciled at startup; stale `running` state, missing processes, or partial artifacts never imply successful completion.
- Cancellation from another supported process is first persisted against the run, then routed to the current fenced owner. Repetition is safe and an unreachable owner remains visibly pending or uncertain until reconciliation.

## 21. Persistence and File Writes

### 21.1 File Write Policy

- Generated files must be written with explicit paths.
- State files must be atomic.
- Report files written by runtimes are validated after runtime exit.
- UltraPlan-created helper files should use `0644` permissions by default.
- Directories should use `0755` by default.
- Durable operational acceptance, event append, lease, cancellation, and terminal writes use the atomicity and durability contract selected for Sprint 35; partial success must be detectable and recoverable.
- Operational retention/compaction never deletes the only current lifecycle snapshot or hides the lower replay boundary.

### 21.2 Locks

UltraPlan must prevent accidental concurrent mutation of package-owned workflow state and stale concurrent mutation of shared operational run state.

Lock requirements:

- Per-study run lock for `run-loop`.
- Lock file includes PID, command, and timestamp.
- Stale lock detection should be conservative.
- Users can force unlock with explicit command or flag.
- Sprint 35 run ownership uses leases plus fencing or an equivalently safe mechanism; a lock file containing only PID and timestamp is not sufficient authority.

## 22. Security Requirements

- Redact secrets from logs and config output.
- Use agentwrap redaction helpers and `SDKError` safe fields for runtime diagnostics where practical.
- Validate workspace paths.
- Avoid shell interpolation for runtime commands where possible.
- UltraPlan product code must not invoke OpenCode directly; `agentwrap/opencode` owns `exec.CommandContext` process execution.
- Treat source repository content as untrusted input.
- Do not execute commands from source repositories unless explicitly configured.
- Disable post-run hooks by default.
- Require explicit opt-in for commands that mutate Git state.
- Runtime permission posture must be expressed through `agentwrap.PermissionPolicy` where possible.
- Unsupported required permission policy features must fail before runtime launch unless explicitly configured as best-effort.
- The Phase 4 server must remain loopback-only, same-origin, CSRF-protected, body-limited, timeout-bounded, and strict about Host/Origin validation.
- Browser routes must not expose arbitrary workspace paths, raw provider payloads, unrestricted stderr, secrets, or executable HTML from workspace Markdown.

### 22.1 Permission Policy Requirements

UltraPlan must support task-level permission policies using agentwrap:

- Read-only study analysis should default to allowing read/search/list-like tools and denying or asking for edit/shell behavior unless a command explicitly needs mutation.
- Study-side commands must not edit source repositories unless explicitly configured for generated artifacts.
- Shell behavior should default to ask or deny unless explicitly enabled.
- Permission policy metadata must be persisted from `RunMetadata.Permissions`.
- Permission events must appear in logs/status where relevant.
- Path-level permission rules must respect agentwrap OpenCode adapter limitations: unsupported required path rules fail before launch unless best-effort behavior is configured.

## 23. Testing Requirements

### 23.1 Unit Tests

Required unit coverage:

- Workspace discovery.
- Path normalization.
- Config loading and precedence.
- Study YAML parsing.
- Dimension generation.
- Source and dimension resolution.
- Markdown source discovery.
- Markdown frontmatter parsing.
- Markdown frontmatter stripping.
- Applicable dimension normalization and filtering.
- Prompt composition.
- State transitions.
- Backoff calculation.
- Error classification and mapping from `agentwrap.SDKError`.
- Report validation.
- Rating parsing.
- Code citation parsing.
- Code reference resolution.
- Summary generation.
- Project catalog validation.
- Sprint planning artifact validation.
- Agentwrap `RunRequest` construction for analysis, synthesis, and planning artifact generation tasks.
- Agentwrap wrapper composition.
- Permission policy construction.
- HTTP route and method mapping, request validation, safe error projection, confirmation binding/expiry, and template model construction.
- SSE event encoding, bounded buffering, slow-subscriber behavior, cursor replay/gaps, reconnect behavior, durable-run lookup, and cancellation routing.
- Run envelope validation, lifecycle transitions, monotonic sequence allocation, lease renewal/expiry, fencing, terminal arbitration, retention/compaction, redaction-before-append, tombstones, and migration.

### 23.2 Fixture Tests

Fixtures required:

- Minimal valid study.
- Study with missing source.
- Study with ambiguous dimension reference.
- Study with directory sources and Markdown document sources.
- Markdown document source with `applicable_dimensions`.
- Markdown document source with malformed frontmatter.
- Valid per-source report.
- Invalid per-source report.
- Valid final report.
- Invalid final report.
- Code citations with single lines, ranges, lists, missing files, and source-prefixed paths.
- OpenCode structured event stream.
- Malformed runtime event stream.
- Rate limit event or diagnostic.
- Timeout fixture through fake runtime.
- Agentwrap event records for lifecycle, validation, retry, fallback, permission, and cleanup.
- Durable run records for active, stalled, interrupted, cleanup-uncertain, terminal, migrated, tombstoned, compacted, corrupt/torn, and unsupported-version cases.
- Event journals with duplicate delivery, missing cursors, explicit retention gaps, redacted payloads, and terminal races.

### 23.3 Fake Runtime Tests

Fake runtime tests should use agentwrap-compatible fakes. If the public agentwrap package does not expose a reusable fake, UltraPlan should define local test fakes that implement `agentwrap.Runtime` and `agentwrap.Run`.

Fake runtime must support:

- Successful run that writes expected artifact.
- Successful runtime exit without artifact.
- Runtime failure.
- Timeout.
- Cancellation.
- Malformed events.
- Rate limit.
- Delayed events for concurrency tests.
- Validation failure and repair paths through `agentwrap.ValidatingRuntime`.
- Retry and fallback paths through `agentwrap.PolicyRunner`.
- Observability records through `agentwrap.ObservingRuntime`.

### 23.4 Integration Tests

Gated integration tests may require:

- OpenCode installed.
- Provider/model configured.
- Network access if runtime requires it.

These tests must be opt-in and skipped by default unless required environment variables are present.

OpenCode integration tests must go through `agentwrap/opencode`, not direct process invocation.

Local web integration tests must use `httptest` with real app use cases over temporary workspaces and fake runtime/process dependencies. Browser-level tests must cover read-only navigation, guarded confirmation, live progress, cancellation, refresh/reconnect recovery, hostile artifact content, and CLI/TUI/web state agreement without requiring a live provider.

Sprint 35 cross-process tests must cover two concurrent CLI-started runs visible from the browser top bar, a run observed through a second supported local server, session expiry, observing-server restart, retained replay followed by new live events, owner crash at every commit boundary, stale lease, PID reuse, fencing of stale writers, duplicate cancellation, cancellation/completion races, slow subscribers, cursor expiry, compaction, disk-full/permission failure, corrupted/torn records, schema migration, and a redacted support bundle. The test harness must distinguish observer restart from worker-owner restart and assert the selected topology contract explicitly.

### 23.5 Golden Tests

Golden tests should cover:

- Generated dimension Markdown.
- Generated study README.
- Directory source prompt composition.
- Markdown document source prompt composition with stripped frontmatter.
- Summary CSV.
- Code extraction output.
- Project catalog output.
- Sprint planning artifacts through `plan.md`.
- Sprint review Markdown with contract, handbook, plan, evidence, finding, and verdict sections.
- Sprint smoke Markdown with selected scope, run ID, result, evidence links, issues, and verdict.
- Review/smoke/status JSON envelopes.

Golden updates must be explicit.

## 24. Build and Release Requirements

### 24.1 Build

- Go module builds with `go build ./...`.
- Tests run with `go test ./...`.
- CLI binary is produced by `go build ./cmd/ultraplan`.
- Go module must depend on `github.com/Antonio7098/agentwrap`.
- Local development may use a `replace` directive or workspace file to point at the sibling `agentwrap` checkout.
- Release builds must pin an explicit agentwrap version or commit.
- Agentwrap documentation under `agentwrap/docs` is the canonical reference for runtime integration behavior.
- Go builds embed the Phase 4 HTML templates, CSS, and minimal JavaScript. A JavaScript package manager, Vite server, or separate asset build is not required to build or run the initial web UI.

### 24.2 Version Command

CLI must support:

```bash
ultraplan version
```

Version output:

- Version.
- Commit SHA when available.
- Build date when available.
- Go version.

### 24.3 Distribution

First release should support:

- Source build.
- Linux binary.
- macOS binary.

Release artifacts should include checksums.

## 25. Migration From Prototype

The Go implementation must preserve the important prototype workflows:

- Study listing.
- Study initialization.
- Single dimension/source run.
- Full study run.
- Stateful run-loop.
- Status.
- Synthesis.
- Summary.
- Code extraction.
- Project catalog validation.
- Sprint planning flow through `plan`.
- Sprint execute flow through `execute`.
- Sprint automated review through root `review.md`.
- Sprint deep-smoke summary through root `smoke.md` with detailed evidence retained in the external harness.
- CLI/TUI operation of review, smoke, focused reruns, cancellation, evidence inspection, and recovery.
- Markdown document sources with `applicable_dimensions`.

The Go implementation may intentionally change:

- Command grouping for clarity.
- Config file format.
- Internal state schema.
- Runtime adapter implementation.
- Validation strictness.
- Logging format.

The Go implementation must intentionally change runtime execution from prototype-owned OpenCode process management to agentwrap-owned runtime supervision.

Migration support should include:

- Ability to read existing study directory layouts where practical.
- Validation command that reports incompatible or missing files.
- Documentation for converting prototype config into `ultraplan.yml`.
- Documentation for removing/replacing manual `review.md` and legacy `deep-smoke.md` with the Phase 3 generated `review.md` and `smoke.md`.
- Validation that recognizes legacy artifacts without treating them as current Phase 3 evidence.

## 26. Acceptance Criteria

UltraPlan Go is technically acceptable when:

- `go test ./...` passes.
- CLI builds as a single binary.
- `ultraplan --help` and command help are complete.
- Workspace init and validation work.
- Study init from YAML creates expected structure.
- Source discovery finds directories and top-level Markdown document sources.
- Markdown source `applicable_dimensions` filters source/dimension task creation and validation.
- Study list commands work.
- Single fake-runtime analysis writes and validates a report.
- Runtime execution uses `github.com/Antonio7098/agentwrap` and the `agentwrap/opencode` adapter.
- Default runtime composition includes `ObservingRuntime`, `ValidatingRuntime`, and `PolicyRunner`.
- Fake-runtime run-loop resumes after interruption.
- Synthesis task runs only after valid per-source reports exist.
- Summary generation is deterministic.
- Code extraction resolves valid citations and reports unresolved ones.
- OpenCode adapter passes gated smoke tests.
- Runtime health gates use agentwrap health checks.
- Runtime validation uses agentwrap validation expectations and custom validators.
- Runtime retry/fallback uses agentwrap policy metadata and events.
- Runtime permission policies use agentwrap permission metadata and OpenCode translation.
- Project catalog commands validate required files.
- Sprint planning loads required files and updates flow-state validation through `plan`.
- Sprint execute extracts validated plan tasks, updates `.run-state.json`, and reports task progress through status.
- Sprint review dynamically covers every selected contract plus the technical handbook, applies deterministic verdict rules, and atomically writes valid `review.md`.
- Sprint smoke runs after review by default, selects or validates the narrowest sufficient external harness scope, and atomically writes valid `smoke.md` linked to matching harness evidence.
- Review/smoke freshness, blocking, cancellation, rerun, and recovery behavior agree across CLI text, JSON, status, and TUI.
- Every Phase 3 sprint exposes its completed functionality in the TUI through shared typed app use cases.
- Logs redact secrets.
- State writes are atomic.
- Cancellation saves state.
- Documentation explains normal workflows and recovery workflows.
- `ultraplan serve` starts the loopback-only local server and serves the embedded browser UI from the same binary.
- HTTP handlers use shared app use cases; browser pages do not shell out to the CLI or duplicate product workflow logic.
- Guarded web operations require a current server-issued confirmation, stream bounded SSE progress, support explicit cancellation, and recover from durable state after refresh or restart.
- Local web security, redaction, path-containment, SSE, browser, race, and graceful-shutdown tests pass.
- Graceful server shutdown cancels every server-owned active operation, performs bounded cleanup, records a truthful outcome, and keeps browser disconnection independent from run cancellation.
- `code-context` is ordered immediately after requirements, writes a structurally valid `code-context.md` from read-only target-repository inspection, and makes sprint-index ready.
- Existing pre-code-context sprint state remains usable through explicit compatibility handling.
- Exact requirements and code-context content appears in a stable shared prefix for every downstream agent-backed stage, while additional repository inspection remains allowed.
- CLI, JSON, TUI, and browser agree on code-context readiness, progress, findings, artifact, rerun, cancellation, and recovery.
- No repository index, RAG/cache subsystem, parallel context manifest, or automatic staleness mechanism is required by the grounded-planning track.
- Every runtime-backed execution is durably accepted with a stable workspace run ID before child start, or fails closed without starting the child.
- Workspace-wide run queries and the browser top bar include active CLI-, TUI-, and web-started work from the shared lifecycle projection.
- A supported second local server can inspect a current run, replay retained ordered events, and receive subsequently committed events without depending on the originating browser session.
- Refresh, session expiry, observing-server restart, bounded delivery reaping, and legacy operation URLs resolve to durable run state, an explicit tombstone/gap, or precise recovery rather than an unexplained 404.
- Lease/fencing, owner loss, PID reuse, cancellation routing, terminal arbitration, reconciliation, retention, persistence failure, and backpressure have deterministic race and fault-injection coverage.
- Logs, metrics, health, diagnostics, and support exports correlate UltraPlan run/attempt/stage/task identities with agentwrap/runtime/process identities safely.
- Sprint 35 operational persistence does not change authority for governed artifacts, flow outcomes, Git/source state, or external smoke evidence.
- Sprint 36 QA mapping is deterministic, covers every changed path with a bounded verification surface, persists resumable theory outcomes, and prevents investigator mutation.
- Sprint 37 writable investigation fails closed without validated isolation; global adjudication alone promotes evidence-backed issues; smoke-as-QA preserves current protocol and evidence guarantees.
- Sprint 38 repair consumes frozen adjudicated issue packets, enforces allowed scope, progressively reverifies changes, and exposes bounded stalled/escalated outcomes.
- Sprint 39 dogfood measures evidence validity, isolation reliability, false-positive/inconclusive rates, repair convergence, cancellation, recovery, and cross-surface agreement before content identity proceeds.

## 27. Open Technical Questions

- Should the workspace marker be `ultraplan.yml`, `.ultraplan/`, or both?
- Should report validators be configurable per study, or fixed for the first release?
- Should retries be per task only, or should there also be run-level retry budgets?
- How should UltraPlan persist agentwrap `RunStore` records durably across process exits?
- Which agentwrap wrapper order should be used when validation failures should be eligible for fallback?
- Should lock files be mandatory for all mutating commands or only long-running run loops?
- Should code extraction support non-local source paths in the first release?
- Should generated report templates be versioned independently from the CLI binary?
- What is the minimum stable JSON schema for status output?
- What long-term compatibility guarantee should the external smoke-harness protocol provide across independently released harness versions?
- Which target identity should be mandatory when Git metadata is unavailable: execute-state fingerprint, explicit changed paths, or a full contained tree manifest?
- What exact topology does Sprint 35 support: multiple local processes, shared-filesystem multi-host observers/owners, or both, and which guarantees are conditional?
- Should execution control use a coordinator/daemon, fenced direct repository writers, or a hybrid, and where does that responsibility live without duplicating product modules or agentwrap?
- Which durable mechanism best satisfies acceptance atomicity, ordered append, concurrent observation, recovery, and bounded retention: filesystem journal/snapshots, SQLite, agentwrap store composition, or another local design?
- Which workers may be adopted after owner loss, if any, and which must be cancelled or reconciled as interrupted?
- What is the stable identity hierarchy among operation, product run, attempt, stage execution, task, agentwrap run/session, OS process, and smoke-harness run?
- What lease, heartbeat, clock, fencing, and process-birth rules prevent stale writers and PID reuse without falsely killing slow healthy work?
- What event delivery guarantee, replay cursor, deduplication, retention, compaction, snapshot, backpressure, and disk-limit policy is supportable?
- Which run history is visible across sessions, and which cancel/retry controls require fresh session authorization?
- How long do legacy operation mappings and terminal tombstones remain, and what typed response replaces a missing/expired mapping?
- Which metrics and health surfaces are always local, and is OpenTelemetry export an optional Sprint 35 adapter or a later concern?
- What is the safe fail-closed or degraded behavior for acceptance, mid-run append, heartbeat, and terminal persistence failures?

## 28. Changelog

| Version | Date | Change | Rationale |
|---------|------|--------|-----------|
| 1.0.0 | 2026-05-25 | Initial TRD | Define production technical requirements for UltraPlan Go. |
| 1.1.0 | 2026-05-25 | Ground runtime supervision in agentwrap | Require UltraPlan Go to use agentwrap runtime contracts, wrappers, OpenCode adapter, validation, policy, health, observability, metadata, and permissions. |
| 1.2.0 | 2026-06-13 | Add project and sprint planning through plan | Scope Phase 2 to governed planning artifacts while deferring execution, smoke, review automation, issue tracking, and Git mutation. |
| 1.3.0 | 2026-07-02 | Add sprint execute scope | Expand Phase 2 to controlled implementation execution from validated `plan.md` tasks while keeping smoke, review automation, issue tracking, and Git mutation deferred. |
| 1.5.0 | 2026-07-17 | Add Phase 3 review and deep smoke | Define root `review.md`/`smoke.md`, dynamic structured review, external harness integration, review-before-smoke gates, freshness, and full CLI/TUI parity. |
| 1.6.0 | 2026-07-22 | Add Phase 4 local web surface | Define the Sprint 30+ loopback Go HTTP server, embedded Go-rendered browser UI, guarded commands, SSE progress, and local-only security boundary. |
| 1.7.0 | 2026-08-15 | Add grounded planning and gated future architecture | Define the `code-context` stage and stable downstream prompt prefix, strengthen server shutdown ownership, formalize the embedded primitives/components/layouts/pages hierarchy, and record evidence gates for content, QA/repair, retrieval, persistence, optional graph, and cloud/Aren. |
| 1.8.0 | 2026-08-20 | Add durable run control and observability | Define stable workspace run identity, durable ordered safe events, cross-surface/server projection, lease/fencing and reconciliation, compatible stable inspection, and correlated operational diagnostics. |
| 1.9.0 | 2026-08-21 | Add Product Phase 5 Sprints 35–39 | Treat Sprints 33–34 as the grounded-planning track; define one sprint each for durable run observability, read-only QA, evidence QA and smoke integration, bounded repair, and QA/repair dogfood before content identity. |


## Current Scope Clarification

UltraPlan has two connected product sides: (1) studying source repositories/documents and producing validated research artifacts, and (2) applying selected study findings to governed project/sprint planning, controlled implementation, Conformance Review, empirical QA, and bounded repair. Phase 4 adds the local browser interface, and the delivered grounded-planning track adds `code-context` through the same application boundaries. Product Phase 5 uses Sprints 35–39 for durable operational run identity, read-only QA, isolated evidence production and smoke integration, adjudicated bounded repair, and QA/repair dogfooding. Content identity, retrieval, alternate product-artifact persistence, knowledge graphs, cloud, and Aren remain gated after Sprint 39. General-purpose issue tracking, hosted/multi-user service, remote workers, automatic Git mutation, and unbounded autonomous repair remain non-goals.
