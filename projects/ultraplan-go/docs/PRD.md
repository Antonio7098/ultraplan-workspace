# Product Requirements Document: UltraPlan Go

**Version:** 1.9.0
**Status:** Draft
**Owner:** Product and Engineering
**Last Updated:** 2026-09-04

## Stage 1: Product Brief

### 1.1 Executive Summary

UltraPlan Go is a production-grade local planning, implementation, review, QA, and research system. Phase 1 implements the study side: architecture studies across source repositories/documents, structured report synthesis, summary generation, validation, and cited code-reference extraction. Phase 2 adds governed project and sprint planning through `plan`, then controlled implementation execution through `execute`. Sprints 24 and 25 add the local TUI foundation and guarded operational controls. Phase 3 completes the initial sprint workflow with automated conformance review followed by sprint-targeted deep smoke. Phase 4, beginning with Sprint 30, adds a loopback-only Go HTTP server and a simple Go-rendered browser UI over the same typed application use cases. Sprints 33–34 form a delivered grounded-planning track that adds `code-context`. Product Phase 5 comprises Sprints 35–40: durable run identity and cross-surface observability, read-only QA decomposition and synthesis, evidence-producing QA with smoke integration, bounded repair, an optional requirements-driven performance stage, then QA/repair/performance dogfooding and hardening.

Sprints 35–40 are the ordered Product Phase 5 delivery group, with each sprint gated on the preceding sprint's evidence. Content identity follows Sprint 40 so its schema is shaped by real theories, evidence, issues, repairs, performance attempts, and findings. Retrieval, alternate product-artifact persistence, knowledge graph, cloud, and Aren remain later gated options. Sprint 35 operational persistence does not select a new authority for authored product artifacts.

- **Problem Statement:** The current prototype proves the workflow, but it is script-like, tightly coupled to a local Bun/TypeScript environment, and not hardened for durable runs, reproducible outputs, reliable retries, clear configuration, or long-term extensibility.
- **Proposed Solution:** Rebuild UltraPlan as a Go CLI with a well-defined domain model, deterministic filesystem layout, resumable orchestration, structured runtime adapters, robust validation, and production-quality testing.
- **Target Users:** Engineers, technical leads, AI workflow builders, and product teams using agentic coding runtimes to study codebases, compare architectural patterns, and plan future implementation work from separately reviewed study outputs.
- **Expected Outcome:** Users can initialize studies, run large batches of source analyses, synthesize final reports, extract code citations, inspect study status, use selected study evidence to generate governed sprint planning artifacts through `plan.md`, execute validated sprint tasks, review implementation against selected contracts and evidence, and deep-smoke runtime-facing claims with durable state and diagnostics.
  Every supported workflow remains available through stable CLI/JSON surfaces. The TUI and browser reuse the same typed use cases, canonical workspace artifacts, and durable workspace-wide run projection. A run remains discoverable and explainable when the observing page, browser session, or local server process changes.

### 1.2 Product Context

The prototype in `ultraplan/cli` demonstrates the core product shape:

- A global `study` CLI.
- Study directories with sources, dimensions, per-source reports, final reports, and summaries.
- YAML-based study initialization.
- Agent runtime execution through OpenCode.
- Parallel and resumable study runs.
- Structured prompts and output templates.
- Code-reference extraction from generated reports.
- Project documents for PRDs, TRDs, roadmaps, project indexes, sprint requirements, sprint indexes, technical handbooks, reasoning, and plans in the prototype.
- Sprint planning and sprint execution workflows in the prototype.

For UltraPlan Go, Phase 2 ports the planning artifact chain from the prototype through `plan.md` and adds controlled implementation execution from validated plans. The delivered TUI enabling track exposes those services locally. Phase 3 replaces manual sprint review and smoke coordination with product-owned `review.md` and `smoke.md` stages. Detailed smoke runs and issue evidence remain in the external harness referenced by the project index. Phase 4 adds a local browser surface, and the grounded-planning track adds `code-context` through the same boundary. Product Phase 5 uses Sprints 35–40 to make run observation durable, add read-only and evidence-producing QA, introduce adjudicated bounded repair, add requirements-driven performance work, and harden the complete verification loop without changing product-artifact authority boundaries.

The Go product keeps those core capabilities but makes them reliable, portable, testable, and extensible.

### 1.3 The Opportunity

- **Why Now:** Teams are increasingly using autonomous coding agents for research and implementation, but durable planning workflows still depend on ad hoc scripts, one-off prompts, manual citation chasing, and fragile long-running terminal sessions.
- **Problem Archetype:** Hair on Fire and Hard Fact. Long-running AI research/planning workflows are valuable, but failures, missing outputs, weak citations, and non-resumable execution create immediate operational pain.
- **Strategic Fit:** UltraPlan Go first becomes the durable study backbone for architecture research, then applies selected findings to governed sprint planning artifacts and controlled implementation execution.
- **North Star Metric:** Percentage of planned study runs that complete with valid artifacts, valid citations, and no manual recovery.
- **Strategic Priority:** High. The product converts exploratory AI planning into repeatable engineering infrastructure.

### 1.4 Product Principles

- **Cited claims over unsupported summaries:** Every architectural claim in generated reports should trace to a concrete source reference where applicable.
- **Runtime success is not product success:** A run is not complete until required outputs exist and pass validation.
- **Resumability by default:** Long-running work must survive interruption, process exit, and partial failures.
- **Explicit state and explicit errors:** Users must be able to inspect what is pending, running, failed, retrying, completed, and why.
- **Small composable primitives:** Studies, dimensions, sources, runs, reports, projects, and planning sprints must remain clear concepts.
- **Product workflows stay product-owned:** Runtime adapters execute prompts; UltraPlan owns study behavior, project catalog behavior, sprint planning behavior, and sprint execute behavior.
- **Portable local-first operation:** The primary product is a CLI that works from a repository checkout without requiring a server.
- **One product core, multiple local surfaces:** CLI, TUI, and local browser surfaces must share product services and use-case wiring instead of duplicating workflow logic or shelling out to each other for normal operation.
- **A run outlives its observer:** Execution identity, lifecycle, and retained safe history must not depend on one browser session, HTTP server, SSE subscriber, or top-level page.
- **Review before expensive proof:** Evidence-grounded review runs before live smoke so deterministic conformance failures block unnecessary runtime work.
- **Summaries local, raw evidence linked:** Sprint roots keep readable `review.md` and `smoke.md`; detailed smoke evidence remains in the cataloged harness and is linked by stable run IDs.

### 1.5 Target Personas

**Architecture Researcher**

- Role: Senior engineer or technical lead.
- Goals: Compare mature open-source implementations across dimensions and extract usable design guidance.
- Pain Points: Manual research does not scale, citations are hard to verify, and synthesis loses source detail.
- Key Behaviors: Creates study definitions, reviews generated reports, inspects code citations, and uses findings as input to later design decisions.

**Agent Workflow Operator**

- Role: Developer running long batches of autonomous agent work.
- Goals: Start, monitor, pause, resume, and recover many agent tasks.
- Pain Points: Runtime failures, rate limits, missing artifacts, and process interruptions require manual babysitting.
- Key Behaviors: Uses `run-loop`, status views, retry policies, validation reports, and logs.

**Local Workflow Navigator**

- Role: Engineer or technical lead supervising multiple UltraPlan artifacts from a terminal.
- Goals: Browse projects, studies, sprints, validation findings, run state, and generated artifacts without remembering every command.
- Pain Points: Long command sequences and dense status output slow down local investigation and recovery.
- Key Behaviors: Opens a TUI dashboard, moves between projects/studies/sprints, runs safe validation/status actions, inspects failures, and launches controlled workflows only after reviewing impact.

**Verification Operator**

- Role: Engineer or technical lead deciding whether an implemented sprint is ready to accept.
- Goals: Review implementation against the sprint's selected contracts and evidence, run the narrowest sufficient real-system smoke, inspect failures, and rerun only the affected checks.
- Pain Points: Manual review is inconsistent, smoke commands are environment-specific, raw run evidence is hard to correlate, and a runtime exit code alone does not prove product behavior.
- Key Behaviors: Runs review before smoke, inspects `review.md` and `smoke.md`, follows linked harness evidence, resolves or explicitly accepts findings, and uses the TUI for progress, cancellation, reruns, and recovery.

**Tool Extender**

- Role: Engineer adding new runtime adapters, validators, output formats, or product commands.
- Goals: Extend UltraPlan without rewriting the study model or breaking existing workflows.
- Pain Points: Prototype code combines CLI parsing, prompts, state, runtime calls, filesystem writes, and validation in one layer.
- Key Behaviors: Works in Go packages with stable interfaces, tests, fixtures, and clear extension points.

### 1.6 User Scenarios

1. **Initialize a comparative architecture study**
   - As an architecture researcher, I want to create a new study from a YAML definition so that sources, dimensions, report folders, and usage docs are generated consistently.

2. **Run one source against one dimension**
   - As a researcher, I want to run a single dimension/source analysis so that I can test prompts, validate a dimension, and inspect output before launching a batch.

3. **Run a complete study**
   - As an operator, I want to run all dimensions against all sources with controlled parallelism so that large research batches complete without hand scheduling.

4. **Resume interrupted work**
   - As an operator, I want a stateful loop to resume incomplete analyses and syntheses so that terminal interruption or runtime failure does not restart completed work.

5. **Synthesize per-source reports**
   - As a researcher, I want all per-source reports for a dimension synthesized into a final comparative report so that I can understand cross-source patterns and tradeoffs.

6. **Extract code references**
   - As a reviewer, I want a command that resolves cited file and line references to source snippets so that claims in a report can be audited quickly.

7. **Inspect status**
    - As an operator, I want clear status for studies, active tasks, failed tasks, retry times, outputs, and summaries so that I can decide whether to wait, retry, inspect, or intervene.

8. **Execute a validated sprint plan**
   - As an operator, I want to run implementation tasks from a validated `plan.md` with durable task state so that agentic implementation can be resumed, inspected, and recovered without manually reconstructing progress.

9. **Inspect work through a TUI**
   - As a local workflow navigator, I want a terminal UI that shows projects, studies, sprints, validation state, and run progress so that I can understand current state without repeatedly composing CLI commands.

10. **Operate workflows through a TUI**
   - As an operator, I want the TUI to run guarded validation, dry-run, flow, and run-loop actions so that I can monitor long-running work and react to failures from one local terminal interface.

11. **Review an implemented sprint**
   - As a verification operator, I want UltraPlan to review the implemented scope against every selected contract, the technical handbook, sprint decisions, plan tasks, and verification evidence so that acceptance is consistent and actionable.

12. **Deep-smoke runtime-facing behavior**
   - As a verification operator, I want UltraPlan to select and run the narrowest sufficient external smoke suite after review so that runtime, CLI, persistence, cancellation, security, and environment claims have executable evidence.

13. **Operate review and smoke through the TUI**
   - As a local workflow navigator, I want review and smoke readiness, scope, confirmation, progress, cancellation, findings, evidence links, reruns, and recovery in the TUI so that the full post-execute workflow is available without losing CLI parity.

14. **Inspect work through a browser**
   - As a local workflow navigator, I want a browser dashboard served by UltraPlan that shows projects, studies, sprints, validation findings, and bounded artifact previews without requiring a separate frontend service.

15. **Operate workflows through a browser**
   - As an operator, I want guarded local workflow actions, live progress, cancellation, and recovery in the browser while durable workspace files remain the source of truth.

16. **Observe current work from any local surface**
   - As an operator, I want every current run in my workspace to appear with the same identity, liveness, progress, history, and recovery action in CLI, TUI, and any supported local browser server, regardless of where it started.

### 1.7 Business and Product Goals

- Provide a durable Go implementation of the proven prototype workflow.
- Reduce manual recovery in large study runs.
- Improve confidence in generated reports through citation validation and report validation.
- Keep study-side workflows compatible with planning-side project, sprint planning, and sprint execute artifacts.
- Provide a governed planning-side workflow that turns selected study evidence into `requirements.md`, `sprint-index.md`, `technical-handbook.md`, `reasoning.md`, and `plan.md`.
- Provide controlled implementation execution from validated `plan.md` tasks with durable state, resumability, and diagnostics.
- Provide a local TUI that improves discovery, monitoring, and guarded operation while preserving the CLI as the stable automation surface.
- Provide deterministic, scope-aware sprint review that replaces the manually coordinated review process and writes the current `review.md`.
- Provide sprint-targeted deep smoke that writes `smoke.md` and links detailed run/issue evidence in the external harness.
- Provide the complete `execute -> review -> smoke` workflow through both CLI and TUI over shared app use cases.
- Provide a loopback-only local HTTP server and simple browser UI over the same app use cases beginning in Sprint 30.
- Provide a workspace-wide durable run list, stable run inspection, replayable progress, conservative liveness, cross-surface cancellation, and correlated diagnostics beginning in Sprint 35.
- Allow future runtime support without changing the study workflow model.
- Keep the CLI understandable enough that users can inspect and edit generated artifacts directly.

### 1.8 Non-Goals

- Do not turn the Phase 4 local server into a hosted SaaS or remotely exposed service.
- Do not add multi-user accounts, team permissions, tenant isolation, or remote workers in Phase 4.
- Do not require users to adopt one particular AI provider.
- Do not hide generated files behind an opaque database-only workflow.
- Do not make study dimensions immutable; users must be able to edit Markdown/YAML artifacts.
- Do not turn UltraPlan into a general-purpose workflow engine.
- Do not implement project management features such as assignment, scheduling, burndown charts, or issue tracker replacement.
- Do not implement smoke investigation or review automation in Phase 2; they enter scope only in Phase 3 after execute is stable.
- Do not turn Phase 3 findings into a general-purpose issue tracker with assignment, scheduling, remote synchronization, or project-management behavior.
- Do not let review or smoke modify product source, product tests, governed planning inputs, or Git state.
- Do not depend on free-form terminal text when a runtime offers structured output.
- Do not make the TUI the only supported interface; scripts and documentation must continue to use stable CLI and JSON surfaces.
- Do not let the TUI call the UltraPlan CLI as its normal integration mechanism. It should share product services and application use cases.
- Do not let the browser UI call CLI subprocesses, own workflow state machines, or persist an alternate product state.
- Do not make an in-memory server registry, browser cookie, open page, or SSE connection the authority for whether a run exists.
- Do not use Sprint 35 operational records as a shadow authority for governed Markdown, flow outcomes, Git/source state, or external smoke evidence.

## Stage 2: Solution Specification

### 2.1 Product Surface

UltraPlan Go ships as a single CLI binary, tentatively named `ultraplan`.

After the planning and execute workflows are stable, the same binary also exposes a local terminal UI, tentatively invoked as:

```bash
ultraplan tui
```

The TUI is a terminal-native surface over the same workspace, project, sprint, study, validation, runtime, and state services. It is not a browser UI and does not introduce a server.

Beginning in Sprint 30, the same binary also exposes a local browser surface:

```bash
ultraplan serve
```

The server binds to a loopback address by default, serves Go-rendered HTML templates plus embedded CSS/JavaScript, exposes a versioned local HTTP API, and streams operation progress with Server-Sent Events (SSE). It does not require Node.js, Vite, a separate frontend process, or a database at runtime.

Beginning in Sprint 33, sprint planning also includes a first-class `code-context` stage immediately after requirements. One requirements-driven agent inspects the target implementation repository and writes a curated `code-context.md` containing exact source excerpts, repository-relative locators, rationale, relationships, and open questions. The stored Markdown is reused unchanged as a common prompt foundation; it is neither a repository index nor a restriction on further source inspection.

The CLI must support these top-level areas:

- Study discovery and inspection.
- Study initialization from YAML.
- Single analysis runs.
- Full study runs.
- Stateful resumable run loops.
- Run status reporting.
- Final report synthesis.
- Summary generation.
- Code-reference extraction.
- Project discovery, inspection, and project-index validation.
- Sprint planning artifact generation and validation through `plan.md`.
- Sprint implementation-context generation and validation through `code-context.md`.
- Sprint implementation execution from validated `plan.md` tasks through `execute`.
- Configuration inspection and validation.
- Health checks for runtime dependencies.
- Local TUI dashboard and guarded workflow operation.
- Local Go HTTP server, browser dashboard, guarded operations, and SSE progress.

The command surface should be stable, scriptable, and documented through help output.

### 2.2 Workspace Model

UltraPlan operates inside a workspace root. The workspace contains shared configuration, studies, projects, runtime state, and generated artifacts. Stable prompt and output-template defaults are embedded in the UltraPlan binary. Workspace `prompts/` and `templates/` files are optional overrides and should exist only when a user intentionally customizes a built-in default, normally after running `ultraplan defaults install`.

Required workspace concepts:

- **Workspace:** The root directory that contains UltraPlan-managed artifacts.
- **Study:** A named comparative research project.
- **Source:** A local source input studied against dimensions. A source may be a directory, usually a code repository, or a Markdown document file placed directly under a study's `sources/` directory.
- **Dimension:** A Markdown-defined architectural concern or analysis prompt.
- **Per-source report:** Output for one dimension/source pair.
- **Final report:** Synthesis for one dimension across all sources.
- **Summary:** CSV or structured summary of source scores across dimensions.
- **Project:** A governed planning root under `projects/<project>` containing product/technical documents, roadmap, project index, and sprint directories.
- **Project index:** The catalog of available contracts, evidence reports, reasoning templates, review protocols, and project-specific source documents.
- **Planning sprint:** A sprint directory under `projects/<project>/sprints/<slug>` that contains governed planning artifacts through `plan.md` and execute artifacts for controlled implementation runs.
- **Code-context pack:** A sprint-owned Markdown snapshot selected from the target implementation repository after requirements. It is authoritative as the prepared common context for that sprint, while the live repository remains authoritative for source code.
- **TUI session:** A local interactive terminal session that reads and mutates the same workspace artifacts through shared application use cases. It does not create an alternate persistence model.
- **Web session:** A local browser session connected to the loopback UltraPlan server. Connection, confirmation, and subscription state are ephemeral; durable run identity and read visibility are not coupled to the session.
- **Operational run:** A workspace-scoped record of one accepted execution, including stable identity, lifecycle, owner/liveness facts, attempts, safe ordered events, terminal outcome, and correlation to product and runtime state. It is authoritative only for those operational concerns.

### 2.3 Functional Requirements

#### Must Have

1. **CLI bootstrap and help**
   - Users can run `ultraplan --help`.
   - Users can discover available commands and command-specific options.
   - Invalid commands return non-zero exit codes and actionable errors.

2. **Workspace discovery**
   - CLI can locate the workspace root from the current directory or explicit `--workspace`.
   - CLI can validate required workspace folders and config files.
   - CLI errors clearly when run outside a workspace.

3. **Configuration loading**
   - CLI reads workspace configuration from a structured config file.
   - CLI supports command-line overrides for model, runtime, variant, timeout, parallelism, batch size, and output path.
   - CLI supports stage-specific model selection for sprint planning and execute stages, with a global/default model fallback.
   - CLI exposes effective configuration for debugging.
   - Sensitive values must not be printed by default.

4. **Study listing**
   - Users can list all studies.
   - Users can list a study's sources and dimensions.
   - Study source listing must show whether each source is a directory source or Markdown document source.
   - Markdown document sources with `applicable_dimensions` frontmatter must show or expose their dimension filter.
   - Listing output is stable enough for humans and scripts.

5. **Study initialization**
   - Users can initialize a study from a YAML file.
   - YAML supports study name, description, repository/source items, dimension items, desired counts, and optional dimension content.
   - CLI creates study folders, dimension Markdown files, source folder, report folders, study README, and normalized `study-init.yml`.
   - CLI supports dry run, force overwrite, no-clone, custom output directory, and field overrides.
   - Source cloning, when enabled, uses shallow clones and reports per-source success/failure.

6. **Assisted study completion**
   - When a YAML file declares a desired source count greater than the explicit source items, the CLI can ask the configured runtime to suggest additional sources.
   - When a YAML file declares a desired dimension count greater than the explicit dimension items, the CLI can ask the configured runtime to suggest additional dimensions.
   - Suggested additions must be written to a cache artifact before they are merged into the normalized study definition.
   - Runtime-generated source and dimension suggestions must be validated before use.
   - Dry-run mode must show how many sources or dimensions would be requested without invoking the runtime.
   - Users must be able to disable assisted completion and proceed only with explicit YAML entries.

7. **Dimension file generation**
   - Generated dimension files include purpose, steps, expected citations, questions, rating rubric, and output guidance.
   - Generated files are human-editable Markdown.
   - Dimension numbering is stable and zero-padded.

8. **Prompt composition**
   - CLI composes analysis prompts from shared base instructions, selected dimension file, selected source path, and report template.
   - CLI composes synthesis prompts from synthesis instructions, selected dimension file, per-source report manifest, and final report template.
   - Directory source prompts instruct the runtime to explore only the source directory and cite code paths with line numbers.
   - Markdown document source prompts embed the stripped document body directly and instruct the runtime to analyze only the embedded document content without external code or filesystem exploration.
   - Dry-run mode prints or writes prompt previews without executing runtime work.

9. **Single analysis run**
   - Users can run one dimension against one source.
   - If a Markdown document source declares that the selected dimension is not applicable, the command must skip the analysis with a clear message instead of invoking the runtime.
   - CLI creates the output folder if needed.
   - Runtime is invoked with configured working directory, model, variant, timeout, and permission mode.
   - CLI validates that the expected per-source report was written.
   - CLI returns non-zero status if runtime execution or output validation fails.

10. **Full study run**
   - Users can run all matching dimension/source pairs with controlled parallelism.
   - Users can filter by dimensions and sources.
   - The task matrix must exclude Markdown document sources for dimensions not listed in their `applicable_dimensions` frontmatter.
   - CLI synthesizes final reports after per-source analyses complete.
   - CLI generates or updates a summary after analyses and final reports.

11. **Stateful run loop**
   - Users can start or resume a long-running batch.
   - CLI stores durable run state in the study directory.
   - Run state includes task status, attempts, errors, timestamps, retry times, and completion markers.
   - CLI detects completed output files and does not redo valid completed tasks unless requested.
   - CLI validates that state marked complete still has the required files.
   - CLI handles interrupted runs by saving state before exit when possible.
   - CLI state creation, resume validation, completed-source detection, and synthesis gating must ignore inapplicable Markdown document source/dimension pairs.
   - CLI schedules synthesis only when all applicable source reports for a dimension exist and pass validation.

12. **Retry and backoff**
    - CLI supports retry with bounded backoff.
    - CLI classifies failures enough to decide whether retry is allowed.
    - Rate limits are classified separately from generic runtime failures.
    - Users can configure retry count, backoff schedule, timeout, and fallback model/runtime options.

13. **Runtime health checks**
    - CLI can check whether required runtime executables are available.
    - CLI can check configured provider/model availability where runtime support allows it.
    - CLI can fail fast before expensive runs when required setup is missing.

14. **Runtime adapter for OpenCode**
    - First production runtime adapter supports OpenCode.
    - Adapter uses structured output when available.
    - Adapter preserves native event payloads for diagnostics.
    - Adapter maps native events into canonical UltraPlan events.
    - Adapter captures stderr diagnostics, exit status, timeout, cancellation, usage metadata when available, and final status.

15. **Canonical event stream**
    - CLI and internal packages share a canonical event model for lifecycle, messages, tool activity, artifacts, warnings, errors, rate limits, usage, validation, retry, fallback, and final result.
    - Unknown future event types must not crash the system by default.

16. **Report validation**
    - Per-source reports must be checked for existence, non-empty content, required headings, rating, and citation shape where applicable.
    - Final reports must be checked for existence, non-empty content, required headings, rating summary, source table, and citation discipline where applicable.
    - Validation failures must include clear repair context.

17. **Code-reference extraction**
    - Users can run a command against one or more reports to resolve inline code references.
    - The extractor parses a sources table from a report.
    - The extractor resolves references such as `path/to/file.go:42`, ranges, and selected line lists.
    - Output includes source name, path, line numbers, and code snippets.
    - Unresolved references are reported with enough detail to fix the report.
    - Users can write extraction output to a file.

18. **Summary generation**
    - CLI generates a summary of source scores across dimensions.
    - Summary is deterministic.
    - Missing scores are represented distinctly from zero.
    - Inapplicable Markdown document source/dimension pairs are represented distinctly from missing expected reports.
    - Sources are sorted by total score by default.

19. **Markdown document sources**
    - Users can place Markdown files directly in `studies/<study>/sources/` alongside source directories.
    - Markdown document sources may include YAML frontmatter.
    - Supported frontmatter must include `applicable_dimensions`, a list of dimension numbers that the document should be analyzed against.
    - Dimension numbers in `applicable_dimensions` must be normalized so `1` and `01` match the same dimension.
    - Markdown document sources without `applicable_dimensions` are applicable to every dimension by default.
    - Frontmatter must be stripped before document content is embedded in analysis prompts.
    - Markdown document analysis must not require code citations or code exploration.

20. **Logging and diagnostics**
    - CLI exposes human-readable progress logs by default.
    - CLI supports structured JSON logs for automation.
    - Logs include run IDs, task IDs, dimension/source identifiers, attempt numbers, runtime, model, duration, and status.
    - Debug logs include safe diagnostics without leaking secrets.

21. **Cancellation**
    - Users can interrupt active runs.
    - CLI attempts to cancel owned runtime processes.
    - State is saved before exit when possible.
    - Cancellation is visible in run state and logs.

22. **Tests and fixtures**
    - The product includes unit tests for domain models, config, path resolution, state transitions, report validation, code extraction, and prompt composition.
    - The product includes fake runtime fixtures for deterministic runtime behavior.
    - OpenCode integration tests are gated so normal test runs do not require OpenCode.

23. **Local HTTP server**
    - Users can start the local server with `ultraplan serve`.
    - The server binds to loopback by default and rejects non-loopback binding in Phase 4.
    - HTTP handlers call typed app use cases rather than CLI handlers or product modules directly.
    - The server exposes versioned JSON endpoints for dashboard, detail, artifact-preview, validation, confirmation, operation, and cancellation behavior.
    - The server owns only ephemeral connection and subscription state; durable product and operational state are reached through shared application boundaries.

24. **Local browser UI and progress streaming**
    - The Go server renders browser pages from embedded `html/template` assets and serves embedded CSS and minimal JavaScript.
    - Browser presentation is organized as primitives, components, layouts, and route-level pages; composition flows from pages down to primitives.
    - Handlers provide explicit typed view models. Templates do not read files, call application use cases, validate HTTP requests, or decide product workflow state.
    - Users can inspect projects, studies, sprints, findings, state, and bounded Markdown/JSON artifact previews.
    - Runtime-backed or mutating actions require a server-validated confirmation bound to the normalized request and current input state.
    - Commands use ordinary HTTP requests; live operation progress may use SSE and cancellation uses an explicit HTTP request.
    - Browser disconnect or refresh does not determine product task success and does not replace durable run-state recovery.

25. **Sprint 35: Durable run identity and cross-surface observability**
    - Every runtime-backed execution is durably accepted with a stable run ID before child work begins.
    - Active-run counts and lists are workspace-wide and agree across CLI, JSON, TUI, dashboard, and detail pages.
    - Stable run inspection and retained safe events remain available across browser-session expiry, observer restart, and supported local server changes.
    - Event delivery resumes from a durable ordered cursor. Retention gaps, sampled detail, and degraded persistence are explicit rather than silently presented as a complete stream.
    - Owner liveness, cancellation routing, terminal arbitration, and reconciliation are conservative, idempotent, and diagnosable.
    - Read visibility and mutation authorization are separate policies: losing the originating browser session does not erase the run, while cancellation and retry still require current authorization.
    - The selected boundary is same-host `internal/runcontrol` with direct fenced SQLite writers, exact process-birth checks, no worker adoption, bounded full/compacted/tombstone retention, and explicit persistence degradation. It does not replace product-owned workflow or artifact authority.

26. **Grounded sprint code context**
    - The canonical planning order is `requirements -> code-context -> sprint-index -> technical-handbook -> area-reasoning -> reasoning -> plan`.
    - `code-context` can read the resolved implementation repository but may write only the sprint's `code-context.md`.
    - The artifact contains selected exact source excerpts with repository-relative paths, useful ranges/symbols, relevance explanations, important relationships, and explicit uncertainties.
    - Exact `requirements.md` and `code-context.md` content forms a stable common prefix before downstream stage-specific instructions.
    - Sprint index, handbook, reasoning, plan, execute, Conformance Review, and smoke receive the stored pack whenever they invoke an agent.
    - Downstream agents remain free to inspect additional live repository files.
    - Legacy sprints without the stage remain usable through explicit compatibility behavior.
    - No repository index, RAG system, cache subsystem, parallel JSON manifest, automatic staleness system, or provider-specific caching dependency is introduced.

27. **Sprint 36: Read-only QA decomposition and synthesis**
    - Users can generate a deterministic QA map from current execute evidence and governed sprint context.
    - Every changed path belongs to a bounded primary verification surface, with explicit boundary shards where behavior crosses interfaces or state transitions.
    - Read-only investigators can record falsifiable theories, refutations, inconclusive results, static evidence, and bounded cross-shard requests without mutating production or verification code.
    - QA map, shard, theory, synthesis, cancellation, resume, and recovery state is inspectable through CLI, JSON, TUI, and browser.

28. **Sprint 37: Evidence-producing QA and smoke integration**
    - Writable investigators operate only in validated isolated workspaces and may create targeted tests, fixtures, probes, and smoke scenarios there.
    - A global adjudicator validates expectations and evidence before promoting an issue.
    - Users receive canonical `qa.md` plus versioned detailed verification state.
    - Existing smoke behavior is available as a QA executor without losing protocol, containment, timeout, cancellation, cleanup, evidence, or compatibility guarantees.

29. **Sprint 38: Manual and bounded automatic repair**
    - Users can repair one frozen evidence-backed issue within explicitly allowed scope and progressively reverify the result.
    - Repair cannot weaken evidence, requirements, tests, or acceptance criteria.
    - Automatic repair is optional, bounded by cycles and reopenings, resumable, and exposes failed, blocked, escalated, and stalled outcomes distinctly.

30. **Sprint 39: Requirements-driven performance stage**
    - Projects can opt into a `performance` verification phase through explicit project policy; the phase is disabled by default.
    - An enabled sprint takes its complete target set from a validated `Performance Targets` section in `requirements.md`.
    - The performance runtime can author missing benchmarks, establish a repeatable baseline, profile missed targets, and propose bounded implementation changes in isolated copies.
    - Product code freezes target and benchmark identities, runs configured correctness gates, promotes only validated changes, and derives target verdicts from measurements.
    - Required targets block later verification and merge when they remain missed; report-only and baseline targets remain visible without manufacturing a pass.

31. **Sprint 40: QA, repair, and performance dogfooding and hardening**
    - Representative real workflows measure shard quality, false positives, inconclusive outcomes, invalid evidence, isolation reliability, issue quality, and repair convergence.
    - Automatic repair remains disabled unless bounded convergence is demonstrated.
    - Performance dogfood measures benchmark stability, target coverage, correctness preservation, optimization convergence, and stale-evidence handling.
    - QA, repair, and performance artifacts become the evidence base for the later content identity and provenance contract.

#### Should Have

1. **Machine-readable output modes**
   - Commands that list, status, validate, or inspect state support JSON output.

2. **Repair prompts**
   - Validation failures can trigger bounded repair attempts when configured.
   - Repair prompts include missing output details and previous run context.

3. **Run history inspection**
   - Users can inspect prior run attempts and validation results.

4. **Config initialization**
   - Users can generate a starter workspace config.

5. **Template validation**
   - CLI can validate prompt and report templates for required placeholders.

6. **Git integration hooks**
   - CLI can optionally run configured post-write commands such as formatting, tests, commit, or push.
   - These hooks are disabled by default and must be explicit.

7. **Runtime usage metadata**
   - CLI records tokens, cost estimates, duration, and model/provider metadata where available.

8. **Report path normalization**
   - CLI normalizes report paths in generated prompts and outputs so workspaces remain relocatable.

9. **Schema versioning**
   - Run state, config, and structured artifacts include schema versions.

10. **Workspace migration checks**
    - CLI can detect old schema versions and recommend or run migrations.

#### Could Have

1. Additional runtime adapters beyond OpenCode.
2. Direct provider/model workers for non-coding document generation.
3. Watch mode for status dashboards.
4. Remote artifact storage.
5. Pluggable report output formats beyond Markdown and CSV.
6. Advanced source acquisition beyond Git clone, such as local tarballs or registries.

#### Won't Have In The Phase 3 Release (Historical Gate)

Browser UI was intentionally excluded from Phase 3 and enters scope only in Product Phase 4 beginning with Sprint 30. The other exclusions remain deferred unless a later product phase explicitly changes them.

1. Hosted multi-user service.
2. Browser UI.
3. Organization-level permissions.
4. Built-in issue tracker integration.
5. Full workflow DAG engine.
6. General-purpose issue tracking or remote issue synchronization.
7. Silent auto-commit or auto-push by default.
8. Automatic product fixes during review or smoke.
9. Cross-project or cross-sprint verification scheduling.

### 2.4 Command Requirements

The exact command names may be refined during implementation, but the production CLI must cover this surface.

#### Global Commands

- `ultraplan serve`
  - Starts the loopback-only local HTTP server and embedded browser UI.
  - Supports explicit loopback listen address, optional browser opening, and graceful shutdown.

- `ultraplan list`
  - Lists studies in Phase 1. Project listing is handled by `ultraplan project list` in Phase 2.

- `ultraplan init-workspace`
  - Creates the shared workspace structure and starter config.

- `ultraplan config show`
  - Shows effective config with secrets redacted.

- `ultraplan health`
  - Runs workspace and runtime health checks.

- `ultraplan code <report>...`
  - Extracts cited code snippets from reports.

#### Project Commands

- `ultraplan project list`
  - Lists planning projects.

- `ultraplan project <project> status`
  - Shows project docs, project index health, sprint directories, and planning-stage readiness.

- `ultraplan project <project> validate`
  - Validates project docs and `project-index.md` catalog references.

#### Study Commands

- `ultraplan study list`
  - Lists studies.

- `ultraplan study init <study-init.yml>`
  - Initializes a study from YAML.

- `ultraplan study <study> list`
  - Lists sources and dimensions for one study.

- `ultraplan study <study> run <dimension-ref> <source-ref>`
  - Runs one analysis.

- `ultraplan study <study> run-all`
  - Runs all selected analyses and synthesis.

- `ultraplan study <study> run-loop`
  - Runs or resumes stateful batch execution.

- `ultraplan study <study> status`
  - Prints run state and progress.

- `ultraplan study <study> synthesize <dimension-ref>`
  - Synthesizes one dimension from existing per-source reports.

- `ultraplan study <study> summary`
  - Regenerates score summary.

- `ultraplan study <study> validate`
  - Validates study structure and generated artifacts.

#### Sprint Planning, Code Context, Execute, Performance, Review, QA, Smoke, And Repair Commands

Phase 2 supports sprint planning commands through `plan.md` and controlled implementation execution through `execute`. Phase 3 extends the same sprint surface through review and smoke. The grounded-planning track inserts `code-context` after requirements without creating a separate workflow surface:

- `ultraplan sprint <project> <sprint> status`
  - Shows planning artifacts and `flow-state.json`.

- `ultraplan sprint <project> <sprint> validate [stage]`
  - Validates one stage or all planning stages.

- `ultraplan sprint <project> <sprint> prompt <stage>`
  - Renders the prompt for a planning stage without executing runtime work.

- `ultraplan sprint <project> <sprint> flow --to <stage>`
  - Runs planning and execute stages from missing/invalid state through the requested stage.

- `ultraplan sprint <project> <sprint> execute [--task <id>]`
  - Runs pending implementation tasks from a validated `plan.md` or a selected task when supported by the current sprint implementation.

- `ultraplan sprint <project> <sprint> performance [--dry-run|status|resume|cancel]`
  - For projects with performance enabled, reads the frozen targets from `requirements.md`, authors or locates matching benchmarks, records a repeatable baseline, and runs a bounded profile/change/remeasure loop without changing the targets or benchmark definition after baseline.

- `ultraplan sprint <project> <sprint> review`
  - Runs the selected contract and handbook reviewers, deterministic plan/decision/evidence checks, and atomically replaces the sprint's current `review.md`.

- `ultraplan sprint <project> <sprint> conformance-review`
  - Compatibility alias and clearer human-facing name for the existing analytical review capability.

- `ultraplan sprint <project> <sprint> qa [--dry-run|--shard <id>|--suite smoke]`
  - Maps and investigates bounded verification surfaces, then gathers and adjudicates evidence only as allowed by the current QA phase.

- `ultraplan sprint <project> <sprint> smoke`
  - After review, uses the configured smoke author to build or update a durable sprint-specific suite for non-deterministic real boundaries in the cataloged external harness, validates its enumerated coverage, runs it safely, and atomically replaces the sprint's current `smoke.md` with linked evidence.

- `ultraplan sprint <project> <sprint> repair --issue <id>`
  - Repairs one frozen adjudicated issue within confirmed scope and progressively reverifies it.

- `ultraplan sprint <project> <sprint> verify [--to performance|conformance-review|qa] [--repair --max-cycles <n>]`
  - Orchestrates the verification phases without introducing a third canonical assessment artifact.

Supported stages:

```text
requirements
code-context
sprint-index
technical-handbook
area-reasoning
reasoning
plan
execute
review     # compatibility planning-stage projection for Conformance Review
smoke      # compatibility planning-stage projection for the smoke QA executor
```

`performance`, `conformance-review`, `qa`, and `repair` are verification phases rather than additions to the canonical planning-stage list. Compatibility projections preserve current review/smoke clients during migration. Performance runs after execute and before Conformance Review when project policy enables it.

For `code-context`, the existing `prompt`, `validate`, `flow --to`, status, JSON, TUI, and browser surfaces apply. `flow --to plan` runs it exactly once when required. Explicit reruns replace only `code-context.md` atomically.

General-purpose `ultraplan issue ...` and automatic Git mutation remain deferred. Sprint 38 repair is limited to adjudicated verification issues and bounded confirmed scope. Performance mutation is limited to an enabled project, frozen sprint targets, bounded isolated proposals, and product-owned promotion.

### 2.5 Study Initialization Requirements

Input YAML must support:

```yaml
name: go-cli-study
description: Comparative architecture study for Go CLIs
repos:
  count: 3
  items:
    - name: example
      url: https://github.com/org/example
      description: Why this source matters
    - name: guide
      path: sources/guide.md
      description: Markdown guide source analyzed only for dimensions declared in frontmatter
dimensions:
  count: 2
  items:
    - number: "01"
      name: project-structure
      title: Project Structure
      description: Boundaries and package layout
      purpose: What this dimension analyzes
      steps:
        - Inspect module and package layout
      citations:
        - Source files implementing key boundaries
      questions:
        - How are responsibilities separated?
```

Acceptance criteria:

- Missing required fields produce actionable validation errors.
- Counts cannot be less than item counts unless explicitly allowed.
- Dimension numbers are normalized to two digits.
- Generated dimension files preserve supplied steps, citations, and questions.
- Runtime-suggested missing sources and dimensions are cached, validated, deduplicated, and included in the normalized `study-init.yml`.
- Users can opt out of runtime-suggested additions.
- Existing study directories are protected unless `--force` is set.
- Dry run shows planned directories and files.

### 2.5.1 Markdown Document Source Requirements

A Markdown document source is a `.md` file placed directly under a study's `sources/` directory, for example:

```text
studies/<study>/sources/my-guide.md
```

Markdown document sources may declare dimension applicability with YAML frontmatter:

```markdown
---
applicable_dimensions:
  - 1
  - 3
  - 5
---
# My Guide

Content relevant to dimensions 1, 3, and 5.
```

Requirements:

- Source discovery must discover both source directories and top-level `.md` files under `sources/`.
- Directory sources keep the existing behavior: the runtime explores source files and cites code file paths with line numbers.
- Markdown document sources are read as documents, not repositories.
- Frontmatter is parsed before scheduling.
- Frontmatter is stripped before document content is embedded into prompts.
- `applicable_dimensions` values are normalized to two-digit dimension numbers for matching.
- A Markdown document source with `applicable_dimensions: [1, 3, 5]` is analyzed only against dimensions `01`, `03`, and `05`.
- A Markdown document source without `applicable_dimensions` is analyzed against all dimensions.
- Inapplicable document source/dimension pairs are skipped in single runs, batch runs, run-loop state creation, resume validation, completed-source detection, synthesis gating, and summary generation.
- Skipped inapplicable pairs are not failures.
- Markdown document prompts must say that all material is in the embedded document and that the runtime must not access external files or code.
- Markdown document report validation must not require code citations unless a dimension explicitly requires document section citations or another non-code citation rule.

### 2.6 Study Execution Requirements

Each analysis task must have:

- Stable task ID.
- Study name.
- Dimension number and slug.
- Source name.
- Source kind: directory or Markdown document.
- Output path.
- Runtime config.
- Attempt number.
- Status.
- Start/end timestamps.
- Error classification.
- Validation result.

Lifecycle states:

- Pending.
- Ready.
- Running.
- Waiting.
- Retrying.
- Validating.
- Repairing.
- Completed.
- Failed.
- Cancelled.
- Skipped.

Per-source analysis success requires:

- Runtime exits successfully.
- Expected report file exists.
- Report is non-empty.
- Required sections are present.
- Rating can be parsed.
- Required citation format is present for directory sources unless the dimension explicitly permits otherwise.
- Markdown document source reports satisfy document-analysis validation rules and are not failed for lack of code citations.

Synthesis success requires:

- All required applicable per-source reports are valid.
- Runtime exits successfully.
- Final report file exists.
- Final report is non-empty.
- Required sections are present.
- Source summary table exists.
- Rating summary exists.

### 2.7 Applying Study Findings Through Planning, Execution, Performance, Review, And Smoke

The TypeScript prototype demonstrates a second side of UltraPlan: applying study findings to project requirements, roadmaps, sprint reasoning, plans, execution, smoke runs, review, and issue tracking.

Phase 2 includes governed planning and controlled implementation execution. Phase 3 extends the governed flow. Product Phase 5 adds an optional performance phase after execute:

```text
study -> select -> distill -> reason -> plan -> execute
  -> performance, when enabled
  -> conformance review -> QA -> bounded repair -> verified
```

Required planning and execute behavior:

- Users can inspect a project catalog and sprint planning status.
- Users can validate that `sprint-index.md` selects only items present in `project-index.md`.
- Users can generate or validate `technical-handbook.md` from selected evidence reports.
- Users can generate or validate optional `reasoning/*.md` and final `reasoning.md`.
- Users can generate or validate `plan.md` from validated reasoning.
- Users can execute implementation tasks from a validated `plan.md`.
- Users can inspect execute task state, attempts, failures, and completion evidence.
- Users can resume interrupted execute runs without redoing completed tasks unless explicitly requested.
- Users can configure a default runtime model for all sprint stages and override the model for specific stages such as `sprint-index`, `technical-handbook`, `reasoning`, `plan`, and `execute`.
- Projects can enable performance verification explicitly; projects without that policy keep the existing flow unchanged.
- An enabled sprint must declare its targets in a machine-parseable `Performance Targets` section of `requirements.md`. Project and workspace configuration cannot supply or override target values.
- Performance work writes or locates benchmarks, freezes their identity before optimization, measures a repeatable baseline, profiles missed targets, and attempts bounded changes without weakening correctness checks or target definitions.
- Conformance Review and QA run against the post-performance implementation. Any later repair or implementation mutation makes performance evidence stale and requires a current performance recheck before verification or merge can complete.
- Users can dry-run planning stages and render prompts before runtime execution.
- Flow state records planning-stage and execute-stage progress, failures, skipped optional area reasoning, and timestamps.
- Execute run state records task-level progress, attempts, diagnostics, runtime metadata where available, and terminal task states.
- Review resolves only contracts and protocols selected through the project/sprint indexes, runs scope-aware independent reviewers, checks decisions and plan execution, and writes an actionable `review.md` with a deterministic verdict.
- Smoke runs only after review by default, selects the narrowest sufficient cataloged harness scope, and writes `smoke.md` linking the external run and issue evidence.
- Review and smoke state, staleness, verdicts, artifacts, and next actions are visible from CLI, JSON, and TUI.
- Review and smoke are cancellable and recoverable without letting stale evidence satisfy a changed implementation or sprint input.

Deferred behavior remains general-purpose issue management, automatic product fixes, Git mutation, and cross-sprint scheduling.

### 2.8 User Experience Requirements

- Default output should be calm, concise, and useful during long runs.
- Status output should show totals, completed counts, failed tasks, pending tasks, retry times, and current active work.
- Errors should include what failed, why it failed, where to inspect state, and a suggested next command.
- Destructive operations such as overwrite should require explicit flags.
- Dry-run modes should be available for initialization, prompt execution and batch runs.
- Commands should work from nested directories inside a workspace.
- Generated files should be readable and editable without proprietary tooling.

### 2.9 Data and Artifact Requirements

UltraPlan must create deterministic, reviewable artifacts.

Required generated artifacts:

- Study dimension Markdown files.
- Per-source analysis Markdown files.
- Final synthesis Markdown files.
- Study summary CSV.
- Study README.
- Normalized study YAML.
- Run state JSON.
- Run logs or event records.
- Optional code extraction bundles.
- Project `project-index.md`.
- Sprint `requirements.md`.
- Sprint `code-context.md` for sprints using the grounded-planning stage chain.
- Sprint `sprint-index.md`.
- Sprint `technical-handbook.md`.
- Sprint `reasoning/*.md`, when selected.
- Sprint `reasoning.md`.
- Sprint `plan.md`.
- Sprint `flow-state.json`.
- Sprint `performance.md` when project policy enables the performance phase.
- Versioned detailed performance target, attempt, measurement, profile, proposal, cleanup, and result records under the sprint's `verification/` directory.
- Sprint execute summary `execute.md`.
- Sprint execute `.run-state.json`.
- Sprint automated conformance `review.md`.
- Sprint deep-smoke summary `smoke.md`.

Detailed smoke run JSON, stdout/stderr, per-test artifacts, and issue files are generated by the external harness under its cataloged `runs/` and `issues/` directories rather than copied into the sprint. General local issue-tracker artifacts remain deferred.

Artifact requirements:

- Use stable paths.
- Avoid absolute paths in generated files unless necessary for local diagnostics.
- Include schema/version metadata where the artifact is machine-read.
- Be safe to review in Git.
- Avoid secret values.

## Stage 3: Execution Plan

### 3.1 Technical Implications

- **Language:** Go.
- **Distribution:** Single CLI binary.
- **Storage:** Local filesystem artifacts in the workspace.
- **Runtime Integration:** Adapter-based execution with OpenCode as the first adapter.
- **Concurrency:** Bounded worker pools with durable task state.
- **Validation:** First-class validators for config, study structure, reports, code references, and run state.
- **Testing:** Unit tests, fixture tests, fake runtime tests, and gated real-runtime integration tests.
- **Compatibility:** Linux and macOS are required for first release; Windows support should not be intentionally blocked but may be later-release.
- **Local Web Surface:** Go `net/http`, `html/template`, embedded static assets, and SSE; no separate frontend runtime or database.

### 3.2 System Architecture

UltraPlan Go should be organized around these responsibilities:

- CLI command layer.
- Workspace discovery and path resolution.
- Configuration loading and validation.
- Domain model for studies, runs, reports, and runtime events.
- Study initialization service.
- Prompt composition service.
- Runtime adapter service.
- Run scheduler and state store.
- Report validation service.
- Synthesis service.
- Code-reference extraction service.
- Project catalog service.
- Sprint planning artifact service.
- Logging and diagnostics.
- Local HTTP and browser interface adapter.

Data flow for a per-source analysis:

```text
CLI command
  -> workspace discovery
  -> config resolution
  -> study/dimension/source resolution
  -> prompt composition
  -> runtime adapter execution
  -> event stream and logs
  -> expected artifact validation
  -> run state update
  -> summary update
```

Data flow for a study batch:

```text
CLI command
  -> task graph creation
  -> durable run state
  -> bounded worker pool
  -> per-source analysis tasks
  -> validation
  -> synthesis task scheduling
  -> final report generation
  -> summary generation
  -> terminal and structured status
```

Data flow for sprint planning:

```text
CLI command
  -> workspace discovery
  -> config resolution
  -> project/sprint resolution
  -> project-index validation
  -> selected context validation
  -> prompt composition or dry-run preview
  -> runtime adapter execution, when not dry-run
  -> expected artifact validation
  -> flow-state update
```

For grounded-planning sprints, requirements completion first enables `code-context`. That stage resolves the implementation target, performs read-only repository exploration, writes and validates only `code-context.md`, and then enables sprint-index. One shared prompt renderer places exact requirements and code-context content before stage-specific instructions for all later agent-backed stages.

Data flow for sprint execute:

```text
CLI command
  -> workspace discovery
  -> config resolution
  -> project/sprint resolution
  -> prerequisite validation through plan
  -> plan task extraction
  -> execute prompt composition or dry-run preview
  -> runtime adapter execution, when not dry-run
  -> expected task evidence or diagnostic recording
  -> .run-state.json update
  -> execute.md summary update when supported
  -> flow-state update
```

Data flow for sprint performance:

```text
CLI, TUI, or browser action
  -> project performance-policy resolution
  -> validated requirements target parsing and immutable target packet
  -> execute-complete and target-worktree identity gate
  -> isolated benchmark authoring or discovery
  -> correctness check and frozen benchmark identity
  -> repeated baseline measurement and variance qualification
  -> bounded profile, hypothesis, isolated change, and remeasurement loop
  -> product-owned comparison, scope validation, and accepted-patch promotion
  -> final correctness and target verdict
  -> performance.md plus bounded detailed verification state
  -> flow-state summary and downstream freshness update
```

Data flow for sprint review:

```text
CLI or TUI action
  -> workspace/config/project/sprint resolution
  -> prerequisite validation through execute
  -> selected contract/protocol resolution
  -> frozen review scope and input fingerprints
  -> bounded independent agentwrap review requests
  -> deterministic decision/plan/verification/citation checks
  -> deterministic verdict synthesis
  -> atomic review.md replacement
  -> flow-state update
```

Data flow for sprint smoke:

```text
CLI or TUI action
  -> current review verdict and freshness gate
  -> project-index smoke harness resolution
  -> machine-readable harness discovery and preflight
  -> narrowest sufficient scope or explicit suite/test selection
  -> safe external process execution with cancellation
  -> harness run and issue evidence validation
  -> atomic smoke.md replacement with evidence links
  -> flow-state update
```

Data flow for the local browser surface:

```text
Browser request
  -> loopback HTTP route, origin/CSRF/body-limit checks
  -> HTTP request/response DTO mapping
  -> shared app use case
  -> owning product module and durable workspace state
  -> Go-rendered HTML or versioned JSON response

Browser operation start
  -> server-validated confirmation
  -> shared app operation use case
  -> durable run acceptance before child start
  -> product/runtime owner with lease and heartbeat
  -> sanitized ordered event journal before fan-out
  -> transient SSE delivery from a replay cursor
  -> durable run plus product-state refresh on completion, cancellation, or reconnect
```

Graceful server shutdown enters a draining state, rejects new mutations, requests cancellation for every server-owned active operation, waits for bounded cleanup and reconciliation, persists a truthful `cancelled`, `interrupted`, `cleanup_uncertain`, or already-authoritative failure/completion outcome, closes SSE/HTTP, and exits. Closing or refreshing a browser tab does not cancel work.

General issue tracking, unbounded automatic product fixes, and automatic Git mutation remain deferred. The performance phase cannot edit requirements, targets, benchmark definitions after baseline, correctness commands, verification evidence, or Git state.

### 3.3 Dependencies

Required external dependencies:

- Git executable for cloning sources when enabled.
- OpenCode executable for the first runtime adapter.
- One or more configured AI providers, depending on runtime configuration.

Optional dependencies:

- User-configured post-run commands.
- Additional runtime executables in future releases.

Internal dependencies:

- Stable embedded prompt and output-template defaults with explicit workspace-override precedence.
- Workspace config schema.
- Validator definitions.
- Embedded Go HTML templates, CSS, and minimal JavaScript for the local browser surface.

### 3.4 Performance Requirements

- CLI startup for simple commands should complete in under 250 ms on a normal developer machine, excluding filesystem scans of large source trees.
- Listing studies should avoid recursively scanning source repositories.
- Run-loop status should load state in under 1 second for thousands of tasks.
- Code-reference extraction should cache file lookup work within one command invocation.
- Bounded parallelism must prevent unbounded process creation.
- Large reports should be streamed or read intentionally; avoid accidental loading of entire source repositories into memory.
- Markdown document sources should be read once per prompt build or validation pass and should not trigger repository-style recursive scans.

### 3.5 Reliability Requirements

- Run state must be written atomically.
- Partial writes must not corrupt the last known good state.
- Interrupted runs must be resumable.
- Completed tasks must be revalidated before being trusted.
- Runtime exit code 0 must not bypass output validation.
- Cleanup failures must be visible.
- Failed tasks must retain error details and next action.
- Retry loops must be bounded.

### 3.6 Security Requirements

- Secrets must not be logged by default.
- Config display must redact sensitive fields.
- Runtime environment construction must separate safe and sensitive values.
- Generated artifacts must not include API keys or secret env values.
- Source analysis must respect source isolation requirements.
- Markdown document source analysis must respect document isolation requirements: no external filesystem or code access should be requested by the prompt.
- Commands that overwrite or delete user files must require explicit flags.
- User-provided paths must be normalized and checked against the workspace where appropriate.

### 3.7 Observability Requirements

UltraPlan must expose:

- Human-readable progress.
- Structured JSON logs when requested.
- Per-run event records.
- Attempt history.
- Runtime/model/provider metadata.
- Durations.
- Validation results.
- Retry/fallback decisions.
- Warnings and diagnostics.

Status output must answer:

- What is running?
- What is done?
- What failed?
- What is waiting?
- What will retry and when?
- Which artifacts were produced?
- Which validations failed?

### 3.8 Non-Goals and Icebox

Hosted and collaborative service scope remains excluded:

- Hosted service mode.
- Multi-user auth.
- Team collaboration features.
- Full remote execution service.
- Plugin marketplace.
- Complex workflow DAG authoring.
- Automatic source selection from a vague natural-language goal.

Deferred possibilities:

- Remote worker mode.
- Additional runtime adapters.
- Rich client-side application behavior beyond the simple Go-rendered UI.
- Integration with GitHub issues or pull requests.
- General-purpose local or remote issue tracking.
- Automatic Git add/commit/push from sprint execution.
- Automatic product fixes during review or smoke.
- Cross-project/cross-sprint verification scheduling.
- Workspace-level or project-level performance target values that override sprint requirements.
- A performance database, unrestricted benchmark commands, benchmark rewriting after baseline, or indefinite optimization.

Product Phase 5 proceeds in this sprint order:

1. Complete Sprint 35 durable run identity and cross-surface observability.
2. Complete Sprint 36 read-only QA decomposition, investigation, and synthesis before permitting generated checks.
3. Add Sprint 37 isolated evidence-producing QA, global adjudication, canonical QA, and smoke integration.
4. Add Sprint 38 adjudicated manual repair, then bounded automatic repair only if convergence is demonstrated.
5. Add Sprint 39 optional requirements-driven performance verification after execute and require current performance evidence again after later repair.
6. Complete Sprint 40 QA, repair, and performance dogfooding and hardening; disable automatic repair or optimization paths when evidence quality, isolation, measurement stability, or convergence is weak.

After Sprint 40:

1. Design minimal artifact identity, authority, provenance, semantic blocks, and revision-aware evidence from the resulting real QA artifacts while preserving legacy Markdown.
2. Dogfood the combined content and QA schema before retrieval.
3. Add derived lexical retrieval before embeddings; keep indexes disposable and governed selections authoritative.
4. Extract package-owned persistence contracts and consider authored-artifact SQLite only for proven needs.
5. Consider an in-memory read-only knowledge graph only after stable relationships and a retrieval baseline prove a real multi-hop traceability need.
6. Decide filesystem, SQLite/server, or hybrid authority before cloud or Aren integration.

The post-Sprint-40 directions are not current release commitments. A failed evidence gate stops or simplifies the branch.

### 3.9 Risks and Mitigations

- **Risk:** Runtime structured output changes.
  - **Mitigation:** Preserve raw events, version native decoders, fixture-test unknown and malformed events.

- **Risk:** Generated reports omit required files or sections.
  - **Mitigation:** Treat validators as product gates, not optional checks.

- **Risk:** Long runs create corrupted state after interruption.
  - **Mitigation:** Atomic state writes and startup revalidation.

- **Risk:** CLI surface grows confusing.
  - **Mitigation:** Keep commands grouped by workspace, study, and config.

- **Risk:** Too much workflow logic enters runtime adapters.
  - **Mitigation:** Keep adapters limited to execution and event projection; product services own study orchestration, project catalog behavior, sprint planning, and validation.

- **Risk:** Users rely on non-deterministic generated content without review.
  - **Mitigation:** Generated files remain Markdown/YAML and validation is explicit.

- **Risk:** Source repositories are huge.
  - **Mitigation:** Avoid recursive scans except in targeted code-reference lookup, cache lookups within commands, and allow explicit source paths.

## Stage 4: Validation and Launch Plan

### 4.1 Success Metrics

**Primary KPI: Completed Valid Runs**

- Definition: Percentage of requested study analysis and synthesis tasks that complete with valid required artifacts.
- Target: 95% for deterministic fake-runtime test suites; 85% for real OpenCode smoke batches in controlled environments.
- Measurement: Run state records and validation results.

**Secondary KPI: Manual Recovery Rate**

- Definition: Percentage of batch runs requiring manual file repair, state editing, or task rescheduling.
- Target: Less than 10% for supported runtime configurations.
- Measurement: Failed task classifications, repair attempts, and user-reported interventions.

**Secondary KPI: Citation Resolution Rate**

- Definition: Percentage of code references in reports that resolve to local source files.
- Target: 95% or higher for reports generated from local source repositories.
- Measurement: `ultraplan code` resolution summary.

**Guardrail Metric: Invalid Success Rate**

- Definition: Runs marked completed despite missing or invalid required artifacts.
- Target: 0 known cases.
- Measurement: Validator failures, test fixtures, and manual audit.

**Phase 3 KPI: Review Coverage Completeness**

- Definition: Percentage of selected contracts plus the technical handbook that produce valid review results for the current input fingerprint.
- Target: 100% for a passing or pass-with-findings review.
- Measurement: Review task results and `review.md` conformance tables.

**Phase 3 KPI: Smoke Evidence Link Integrity**

- Definition: Percentage of `smoke.md` run and issue references that resolve to matching external harness evidence.
- Target: 100%.
- Measurement: Smoke validation and evidence-path/hash checks.

**Phase 4 KPI: Interface State Agreement**

- Definition: Percentage of representative workspace/project/sprint/study states and terminal operation results that agree across shared app results, CLI/TUI projections, and the local browser surface.
- Target: 100% for documented Phase 4 scenarios.
- Measurement: API compatibility fixtures and CLI/TUI/web integration tests over the same temporary workspaces.

**Grounded-Planning KPI: Shared Context Integrity**

- Definition: Percentage of representative downstream agent prompts that contain byte-for-byte identical validated requirements and code-context content before stage-specific instructions.
- Target: 100% across sprint index, handbook, reasoning, plan, execute, Conformance Review, and smoke fixtures.
- Measurement: Prompt-prefix stability tests plus one gated real-repository requirements-to-plan dogfood flow.

**Sprint 35 KPI: Run Observation Agreement**

- Definition: Percentage of representative active, terminal, interrupted, and cleanup-uncertain executions whose identity, lifecycle, liveness, progress cursor, result, and recovery action agree across CLI, JSON, TUI, and every supported local server topology.
- Target: 100% for the Sprint 35 failure matrix, including concurrent CLI runs, session expiry, reconnect, observer restart, owner loss, retention gaps, and cancellation races.
- Measurement: Cross-process integration/browser tests plus gated real-runtime dogfood over one workspace.

**Sprint 36 KPI: Deterministic QA Coverage**

- Definition: Percentage of changed paths assigned to a bounded primary verification surface with reproducible mapping and inspectable read-only theory outcomes.
- Target: 100% changed-path coverage and identical maps for unchanged inputs in representative fixtures.
- Measurement: QA map fixtures, fingerprint tests, and CLI/JSON/TUI/browser agreement scenarios.

**Sprint 37 KPI: Valid Evidence Promotion**

- Definition: Percentage of promoted issues backed by current, grounded, contained, non-flaky evidence, with invalid setups and stale evidence rejected.
- Target: 100% of promoted issues satisfy adjudication rules in deterministic fixtures and dogfood audits.
- Measurement: Evidence-quality fixtures, isolation tests, smoke-parity tests, and issue audits.

**Sprint 38 KPI: Bounded Repair Convergence**

- Definition: Percentage of attempted repairs that either converge within configured bounds or terminate with an accurate failed, blocked, escalated, or stalled outcome.
- Target: 100% bounded termination; zero passes created by weakened evidence or acceptance criteria.
- Measurement: Repair-cycle state, issue-set deltas, reverification results, and adversarial repair tests.

**Sprint 39 KPI: Performance Target Integrity**

- Definition: Percentage of enabled performance runs where every measured target traces to the current validated sprint requirements, uses the frozen benchmark definition, and receives a product-derived verdict from qualified measurements.
- Target: 100% target traceability and frozen-benchmark integrity; zero passes after correctness failure, target drift, benchmark drift, or inconclusive measurement.
- Measurement: Target-packet fixtures, benchmark identity checks, repeated-sample results, correctness gates, stale-state tests, and one representative real-repository dogfood run.

**Sprint 40 KPI: Verification Trustworthiness**

- Definition: Real-work rates for false positives, inconclusive or blocked investigations, invalid evidence, isolation failures, issue reopening, repair convergence, unstable benchmarks, and optimization convergence.
- Target: Thresholds are selected from initial dogfood evidence; automatic repair and optimization remain disabled where their evidence is not trustworthy.
- Measurement: Representative multi-package dogfood, performance attempts, browser operation traces, and verification-state summaries.

### 4.2 Definition of Done

The first production release is done when:

- CLI builds as a Go binary.
- Workspace initialization and validation work.
- Study initialization works from YAML.
- Study listing works.
- Single analysis run works with fake runtime and OpenCode adapter.
- Full study run works with controlled parallelism.
- Stateful run-loop persists and resumes state.
- Synthesis works from validated per-source reports.
- Summary generation works.
- Code-reference extraction works.
- Config and health commands work.
- Unit and fixture tests pass.
- Gated OpenCode smoke test path is documented.
- User-facing docs describe the supported workflow.
- Sprint review writes a valid current `review.md`, uses every selected contract, and applies deterministic verdict rules.
- Sprint smoke writes a valid current `smoke.md` linked to a real external harness run or records a truthful blocked/not-applicable result.
- Review-before-smoke gating, staleness detection, cancellation, focused reruns, and recovery work from both CLI and TUI.

Product Phase 4 is done when:

- `ultraplan serve` starts a loopback-only HTTP server and serves the embedded browser UI without Node.js or a separate frontend process.
- Browser dashboard and detail pages report the same workspace, project, sprint, study, validation, and artifact state as shared app use cases.
- Guarded operations require server-validated confirmation and expose bounded SSE progress plus explicit cancellation.
- Browser refresh, SSE reconnect, and server shutdown preserve truthful recovery through durable product state.
- Path containment, origin/CSRF checks, secret redaction, request limits, and hostile-Markdown rendering tests pass.
- The embedded UI has a documented and tested primitives/components/layouts/pages hierarchy with namespaced template definitions and layered CSS.
- Graceful shutdown cancels all server-owned active work through the canonical cancellation path, records truthful durable outcomes after bounded cleanup, and restart reconciles abrupt interruption without inferring success.
- Browser disconnection remains independent from operation cancellation.

The grounded-planning track is done when:

- `code-context` is a first-class validated planning stage immediately after requirements.
- A real implementation repository can be inspected to produce a useful `code-context.md` without source mutation.
- Existing pre-stage workspaces remain compatible.
- CLI, TUI, and browser agree on readiness, progress, findings, artifact preview, rerun, cancellation, and recovery.
- Exact requirements and code-context content is reused in a stable common prefix across every downstream agent-backed stage.
- A representative requirements-to-plan dogfood flow runs the stage exactly once and passes normal test, race, and build gates.
- No repository index, retrieval system, UltraPlan cache, parallel context manifest, or automatic staleness claim has entered the release.

Sprint 35 — Durable Run Identity and Cross-Surface Observability is done when:

- Every accepted runtime-backed execution has a stable durable run ID before child execution starts.
- Workspace-wide active counts include current CLI-, TUI-, and web-started work without requiring navigation to an owning detail page.
- A supported second local server can inspect a run, replay retained history, and receive subsequently committed events.
- Refresh, browser-session expiry, observing-server restart, and bounded delivery retention never erase a valid run or produce an unexplained operation 404.
- Lifecycle, liveness, event cursor/gap, result, cancellation, and recovery agree across CLI, JSON, TUI, and browser.
- Owner death, stale leases, PID reuse, duplicate cancellation, terminal races, persistence failure, backpressure, and retention limits have deterministic fault-injection coverage.
- Logs, metrics, health, diagnostics, and the redacted support bundle share safe run/attempt/stage/task/runtime correlation.
- Operational run persistence remains separate from canonical authored artifacts and passes migration, race, build, browser, and gated real-runtime release checks.

Sprint 36 — Read-Only QA Decomposition and Synthesis is done when:

- Conformance Review terminology and compatibility are unambiguous.
- `VerificationPhase` is separate from planning-stage semantics.
- Every changed path maps deterministically to a bounded primary verification surface.
- Read-only investigators persist resumable confirmed, refuted, inconclusive, invalid, and blocked theories.
- Global synthesis supports bounded cross-shard follow-up without issue promotion.
- CLI, JSON, TUI, and browser agree on QA maps, progress, outcomes, cancellation, and recovery.

Sprint 37 — Evidence-Producing QA and Smoke Integration is done when:

- Writable investigators operate only in validated isolated workspaces and uncertainty blocks mutation.
- Evidence adjudication rejects stale, malformed, flaky, ungrounded, or uncontained checks before issue promotion.
- Canonical `qa.md` and versioned detailed verification state remain distinct from `flow-state.json` summaries.
- Smoke runs as a QA executor without losing protocol, containment, evidence, timeout, cancellation, cleanup, or compatibility guarantees.

Sprint 38 — Manual Repair and Bounded Automatic Repair is done when:

- One frozen evidence-backed issue can be repaired and progressively reverified end to end.
- Repair cannot modify evidence, weaken requirements, or escape its allowed scope.
- Automatic repair is bounded, resumable, and terminates distinctly as verified, failed, blocked, escalated, or stalled.

Sprint 39, Requirements-Driven Performance Stage, is done when:

- Performance is disabled by default and appears in the flow only for projects with explicit enabled policy.
- Enabled sprints take every target from the current validated `Performance Targets` section in `requirements.md`; configuration and runtime output cannot add or weaken targets.
- Benchmark authoring, baseline qualification, profiling, bounded optimization, correctness checks, target comparison, cancellation, resume, and recovery use one durable product-owned workflow.
- Baseline and candidate measurements use the same frozen benchmark identity, and product code derives every target verdict.
- Required misses block downstream verification and merge, while report-only and baseline outcomes remain truthful.
- CLI, JSON, TUI, and browser agree on activation, target coverage, progress, measurements, changes, freshness, verdict, and next action.

Sprint 40, QA, Repair, and Performance Dogfooding and Hardening, is done when:

- Representative real-work dogfood measures shard quality, evidence validity, isolation reliability, investigation cost, repair convergence, benchmark stability, and optimization convergence.
- Manual repair succeeds on real work.
- Automatic repair and performance optimization either demonstrate bounded convergence or remain disabled for the unsafe case.
- The resulting QA, repair, and performance artifacts are representative enough to design the later content identity and provenance contract.

Product Phase 5 is done when Sprints 35–40 have passed their individual gates and the Sprint 40 dogfood evidence supports release of the durable QA, repair, and performance workflow.

### 4.3 Instrumentation Requirements

Track these event categories:

- Command started.
- Command completed.
- Command failed.
- Runtime health check started/completed/failed.
- Task queued.
- Task started.
- Runtime event received.
- Artifact detected.
- Validation started/completed/failed.
- Retry scheduled.
- Fallback selected.
- Task completed.
- Task failed.
- Task cancelled.
- Summary generated.
- Run durably accepted or rejected before start.
- Owner lease acquired, renewed, fenced, expired, or reconciled.
- Event journal append, replay, cursor gap, compaction, or persistence failure.
- Subscriber connected, lagged, dropped detail, reconnected, or closed.
- Cancellation requested, routed, acknowledged, or left uncertain.
- Terminal outcome proposed, won, rejected as stale, or reconciled.
- Performance target packet frozen, invalidated, or rejected.
- Benchmark authored, frozen, started, sampled, qualified, or rejected as inconclusive.
- Performance profile captured, optimization proposed, accepted, rejected, or exhausted.

Metrics to monitor locally:

- Task duration.
- Runtime duration.
- Validation duration.
- Attempts per task.
- Failure category counts.
- Citation resolution counts.
- Generated artifact counts.
- Pending/running/completed/failed task counts.
- Active runs by lifecycle and owner-liveness state.
- Lease-renewal failures and reconciliation backlog/latency.
- Durable append latency/failures, journal bytes, retained events, and compaction results.
- Replay distance, cursor gaps, subscriber lag, and dropped/sampled event detail.
- Cancellation routing latency and cleanup-uncertain outcomes.
- Performance targets covered, met, missed, report-only, or inconclusive.
- Benchmark variance, sample count, baseline/candidate comparison, optimization cycles, and correctness-gate failures.

### 4.4 Rollout Strategy

**Rollout Step 1: Internal CLI parity**

- Implement the core prototype behavior in Go.
- Use fake runtime tests and fixture reports.
- Validate generated directory structures against prototype expectations.

**Rollout Step 2: Real runtime smoke**

- Run small OpenCode-backed studies.
- Verify structured events, report writes, validation, status, and retry behavior.
- Compare outputs to prototype workflows where applicable.

**Rollout Step 3: Production hardening**

- Add migration checks, richer validation, docs, error polish, and release packaging.

**Rollout Step 4: Governed sprint planning**

- Add project discovery, project-index validation, sprint planning flow state, planning-stage validators, prompt previews, and runtime-backed generation through `plan.md`.

**Rollout Step 5: Governed sprint execute**

- Add validated plan task extraction, execute prompt previews, runtime-backed task execution, durable `.run-state.json`, execute status, and recovery guidance.
- Stop before review and smoke until execute is stable; general-purpose issue tracking and automatic Git mutation remain excluded.

**Rollout Step 6: Automated sprint review**

- Add dynamic contract/protocol resolution, bounded structured reviewers, deterministic checks and verdicts, root `review.md`, flow integration, and full TUI operation.

**Rollout Step 7: Deep smoke and integrated verification**

- Add agent-driven authoring of durable sprint-specific harness suites, the versioned external harness contract with bounded authoring paths and enumerated coverage, narrow smoke selection, root `smoke.md`, review-before-smoke flow, evidence links, focused reruns, recovery, full TUI operation, and Phase 3 release hardening.

**Rollout Step 8: Local web foundation**

- Beginning with Sprint 30, add explicit interface-runner composition, `ultraplan serve`, a loopback-only Go HTTP server, embedded templates/static assets, read-only dashboard/detail pages, and bounded artifact previews.

**Rollout Step 9: Guarded web operations and streaming**

- Add server-validated confirmations, validation and workflow actions, bounded SSE progress, cancellation, reconnect behavior, and product-owned mutation locking.

**Rollout Step 10: Local web hardening and release**

- Add security, API compatibility, race, browser, accessibility, recovery, packaging, and gated real-runtime coverage without introducing hosted or multi-user scope.

**Rollout Step 11: Code-context vertical slice**

- Add stage order, artifact, validation, runtime configuration, repository-read/output-write boundaries, CLI/app/web surfaces, compatibility, and fake-runtime coverage.

**Rollout Step 12: Shared-context grounded-planning release**

- Reuse exact requirements and `code-context.md` through a stable downstream prompt prefix, add the manual stage skill, dogfood the real-repository flow, and stop before content identity, QA, or retrieval expansion.

**Rollout Step 13: Durable run control and observation**

- Establish stable run identity, workspace-wide lifecycle projection, durable sanitized event history, replayable delivery, liveness/lease and reconciliation semantics, cross-surface cancellation, compatibility, correlated diagnostics, and failure-matrix dogfood.
- Select storage, ownership/coordinator, topology, retention, authorization, and telemetry mechanisms through sprint reasoning; do not treat operational persistence as approval for alternate product-artifact authority.

**Rollout Step 14: Sprint 36 read-only QA decomposition and synthesis**

- Add Conformance Review compatibility terminology, `VerificationPhase`, deterministic QA mapping, bounded read-only investigators, durable theory outcomes, global synthesis, and full cross-surface observation.

**Rollout Step 15: Sprint 37 evidence-producing QA and smoke integration**

- Add validated isolated investigation workspaces, generated checks and probes, evidence adjudication, canonical QA state, issue promotion, and smoke-as-QA compatibility without permitting production repair.

**Rollout Step 16: Sprint 38 manual and bounded automatic repair**

- Add frozen evidence-backed issue packets, explicit confirmation, allowed-scope enforcement, progressive reverification, then bounded automatic cycles with stalled and escalated outcomes.

**Rollout Step 17: Sprint 39 requirements-driven performance stage**

- Add explicit project activation, requirements target parsing, immutable target packets, benchmark authoring and freezing, repeatable baseline qualification, bounded isolated optimization, product-owned verdicts, stale evidence, and cross-surface operation.

**Rollout Step 18: Sprint 40 QA, repair, and performance dogfooding and hardening**

- Exercise representative real changes and failure modes, measure evidence quality, benchmark stability, and convergence, harden recovery and browser operation, and leave unsafe automatic paths disabled unless they prove trustworthy.

**Later gated evaluation**

- Design and dogfood content identity/provenance from real Sprint 36–40 artifacts, then sequence retrieval, product persistence/SQLite, optional graph, and cloud/Aren only through the evidence gates described above.

### 4.5 Review History

| Date | Reviewer | Status | Notes |
|------|----------|--------|-------|
| 2026-05-25 | Product/Engineering | Draft | Initial production PRD based on the prototype workflow. |
| 2026-06-13 | Product/Engineering | Draft | Add Phase 2 governed project/sprint planning scope through `plan.md`; keep execute/smoke/review/issues deferred. |
| 2026-07-02 | Product/Engineering | Draft | Expand Phase 2 to include controlled sprint implementation execution through `execute`; keep smoke/review/issues/Git mutation deferred. |
| 2026-07-17 | Product/Engineering | Draft | Define Product Phase 3 as automated review followed by deep smoke, keep `review.md`/`smoke.md` as the only sprint summaries, link raw smoke evidence externally, and require full TUI parity in every Phase 3 sprint. |
| 2026-07-22 | Product/Engineering | Draft | Add Product Phase 4 beginning at Sprint 30: a loopback-only Go HTTP server, Go-rendered browser UI, guarded operations, SSE progress, and no hosted/multi-user scope. |
| 2026-08-15 | Product/Engineering | Draft | Add the grounded-planning track through `code-context`, adopt the server-owned shutdown cancellation contract, and record the evidence-gated later sequence. |
| 2026-08-20 | Product/Engineering | Draft | Add durable run identity and cross-surface observability after live use exposed server-memory, browser-session, and page-local visibility gaps. |
| 2026-08-21 | Product/Engineering | Draft | Define Product Phase 5 as Sprints 35–39 covering durable run observability, read-only QA, evidence QA with smoke integration, bounded repair, and QA/repair dogfood; move content identity after QA. |
| 2026-09-04 | Product/Engineering | Draft | Insert Sprint 39 requirements-driven performance verification before Conformance Review and move QA, repair, and performance dogfood to Sprint 40. |

### 4.6 Changelog

| Version | Date | Change | Rationale |
|---------|------|--------|-----------|
| 1.0.0 | 2026-05-25 | Initial full PRD | Establish production requirements for UltraPlan Go. |
| 1.4.0 | 2026-07-17 | Add Phase 3 review and smoke | Replace manual post-execute verification with automated `review.md`, external-harness-backed `smoke.md`, integrated CLI/TUI operation, and deterministic gates. |
| 1.5.0 | 2026-07-22 | Add Phase 4 local web surface | Introduce a loopback Go server and simple embedded browser UI over shared app use cases, with HTTP commands and SSE progress starting in Sprint 30. |
| 1.6.0 | 2026-08-15 | Add grounded planning and gated future direction | Add `code-context` as a shared source foundation, make server shutdown cancellation explicit, organize the embedded UI as primitives/components/layouts/pages, and align later work with measured stop/go gates. |
| 1.7.0 | 2026-08-20 | Add durable run control and observation | Commit Sprint 35 outcomes for workspace-wide run identity, replayable safe events, cross-server observation, liveness/reconciliation, compatible stable inspection, and correlated operational diagnostics while leaving the implementation topology open for reasoning. |
| 1.8.0 | 2026-08-21 | Add Product Phase 5 Sprints 35–39 | Add one sprint each for durable run observability, read-only QA, evidence-producing QA and smoke integration, bounded repair, and QA/repair dogfood before content identity. |
| 1.9.0 | 2026-09-04 | Add requirements-driven performance stage | Add opt-in project activation, sprint-requirements target authority, frozen benchmarks, bounded optimization, product-owned verdicts, post-repair freshness, and Sprint 40 dogfood. |
