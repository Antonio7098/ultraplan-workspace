# Sprint Reasoning: Project Domain and Project Index

> Project: `ultraplan-go`
> Sprint: `16-project-domain-and-index`
> Output: `projects/ultraplan-go/sprints/16-project-domain-and-index/reasoning.md`
> **Inputs Used:** `projects/ultraplan-go/project-index.md`, `projects/ultraplan-go/roadmap.md`, `projects/ultraplan-go/sprints/16-project-domain-and-index/requirements.md`, `projects/ultraplan-go/docs/PRD.md`, `projects/ultraplan-go/docs/TRD.md`, `projects/ultraplan-go/docs/ARCHITECTURE.md`, `projects/ultraplan-go/sprints/16-project-domain-and-index/sprint-index.md`, `projects/ultraplan-go/sprints/16-project-domain-and-index/technical-handbook.md`, `templates/sprint-reasoning.md`

This document decides. It synthesizes selected context, handbook evidence, available area-specific reasoning, and contracts into final sprint decisions.

It does not replace `sprint-index.md`, `technical-handbook.md`, or any future `reasoning/*.md` area document.

## Sprint Purpose

- **Goal:** Model `projects/<project>` as a first-class planning root by adding project discovery, project status, and actionable `project-index.md` catalog validation for `ultraplan project list`, `ultraplan project <project> status`, and `ultraplan project <project> validate`.
- **Non-Goals:** Sprint planning artifacts, planning stages, `flow-state.json`, prompt rendering, runtime-backed planning generation, study service consumption, study source/dimension/report models, report validators, rating parsing, summary generation, run-loop scheduling, sprint execution, smoke investigation, automated review generation, issue tracking, Git mutation, browser UI, hosted service, multi-user collaboration, TUI, metrics exporter, local API server, and mutation of existing project docs or planning artifacts.
- **Depends On:** Sprint 1-2 workspace and CLI foundation, Sprint 9 runtime boundary decisions, Sprint 14 validation and diagnostics discipline, Sprint 15 study-side baseline, the roadmap Sprint 16 sequencing, the project index Phase 2 planning context, PRD/TRD/Architecture project planning requirements, selected contracts from `sprint-index.md`, and selected go-cli-study evidence distilled in `technical-handbook.md`.

## Selected Context And Pre-Reasoning Artifacts

| Artifact | Path | How It Was Used |
| --- | --- | --- |
| Sprint Requirements | `projects/ultraplan-go/sprints/16-project-domain-and-index/requirements.md` | Treated as the authoritative sprint contract for outputs, acceptance criteria, non-goals, constraints, dependencies, and review expectations. |
| Project Index | `projects/ultraplan-go/project-index.md` | Established Phase 2 planning scope, available contracts, available evidence reports, source docs, review protocols, and the real catalog table shapes that Sprint 16 must parse. |
| Roadmap | `projects/ultraplan-go/roadmap.md` | Confirmed Sprint 16 sequencing before Sprint 17 flow state and later planning stages; limited the sprint to project discovery, docs, roadmap/index validation, catalog parsing, project list, and project status. |
| Product Requirements | `projects/ultraplan-go/docs/PRD.md` | Confirmed project discovery, inspection, and project-index validation as Phase 2 product surface while deferring sprint execution, smoke, review automation, issues, and Git mutation. |
| Technical Requirements | `projects/ultraplan-go/docs/TRD.md` | Grounded `internal/project` ownership, project model, catalog entries, path rules, project commands, exit codes, validation expectations, and Phase 2 dependency boundaries. |
| Architecture | `projects/ultraplan-go/docs/ARCHITECTURE.md` | Grounded module-owned behavior, `internal/project` responsibility, thin app wiring, platform/product separation, and rejection of global technical-layer packages. |
| Sprint Index | `projects/ultraplan-go/sprints/16-project-domain-and-index/sprint-index.md` | Selected applicable contracts, selected evidence reports, selected Architecture reasoning template, required review protocols, and excluded contexts. |
| Technical Handbook | `projects/ultraplan-go/sprints/16-project-domain-and-index/technical-handbook.md` | Supplied distilled evidence, concrete source references, relevant patterns, trade-offs, anti-patterns, examples, and open questions used to make final decisions. |
| Sprint Reasoning Template | `templates/sprint-reasoning.md` | Provided this document structure and required sections. |

## Area-Specific Reasoning Inputs

No area-specific reasoning files were present under `projects/ultraplan-go/sprints/16-project-domain-and-index/reasoning/` when this document was created. The sprint index selected an Architecture reasoning output path, but no area artifact existed to summarize. Final architecture decisions are therefore resolved directly in this sprint reasoning using `requirements.md`, `sprint-index.md`, `technical-handbook.md`, PRD, TRD, and Architecture inputs.

| Area | Reasoning Document | Key Conclusion | Evidence Basis | Impact On Final Decision |
| --- | --- | --- | --- | --- |
| Architecture | Not present | No separate area-specific conclusion was available. The final reasoning must decide module ownership, dependency direction, CLI thinness, and abstraction boundaries directly. | `sprint-index.md` selected Architecture as relevant; `technical-handbook.md` selected structure, command, DI, IO, error, security, testing, performance, and philosophy evidence; TRD and Architecture prescribe `internal/project` ownership. | Decisions 1, 2, 3, 4, and 5 make architecture final for `plan.md` without requiring a later architecture reopening. |

## Sprint Technical Handbook Summary

- **Relevant Patterns:** Module-owned domain behavior with thin entrypoints; concrete project package before shared abstractions; factory-created commands and use-case services; explicit manual dependency construction; injectable IO and filesystem seams; actionable classified validation errors; fail-closed path boundaries; concise CLI-first terminal UX; behavior-level fixture tests; bounded discovery and lazy work.
- **Important Trade-Offs:** Keep parsing and validation in `internal/project` instead of global catalog/validation packages; keep CLI handlers thin even though this requires service/domain types first; use structured validation findings instead of plain error strings; fail closed on catalog path issues; run command-level tests in addition to package unit tests; avoid eager recursive scans even if status output is less exhaustive.
- **Warnings / Anti-Patterns:** Do not let `internal/app` become the project domain; do not introduce broad shared catalog, validation, planning, reports, or prompts packages; do not consume study services or models; do not hardcode output streams or real paths into test-sensitive code; do not flatten validation into unclassified strings; do not silently ignore malformed catalog rows or missing paths; do not add prompts, spinners, color, TUI, runtime checks, recursive source scans, or workflow execution.
- **Evidence Confidence:** High for package boundaries, command architecture, DI, IO seams, error handling, security, testing, and performance because the handbook cites multiple mature Go CLI sources and final reports. Medium for terminal UX and philosophy because these are more contextual but still aligned with Sprint 16's script-friendly, deliberately narrow command scope.

## Requirement ID Key

Requirement IDs below refer to the acceptance criteria in `requirements.md` in order.

| ID | Acceptance Criterion Summary |
| --- | --- |
| AC-01 | `ultraplan project list` lists direct project roots in deterministic name order. |
| AC-02 | Discovery ignores hidden entries, non-direct children, and non-directory entries under `projects/`. |
| AC-03 | Project names are filesystem-safe and invalid names produce actionable diagnostics. |
| AC-04 | Project status reports presence and health of docs, Markdown docs, roadmap, project index, sprints, and sprint directories. |
| AC-05 | Project validate exits `0` when required files and catalog references are valid. |
| AC-06 | Project validate exits `5` for validation failures. |
| AC-07 | Project-index parsing recognizes source documents, active contracts, evidence reports, reasoning templates, and review protocols. |
| AC-08 | Catalog paths resolve relative to the workspace and fail closed on workspace escape unless explicitly external. |
| AC-09 | Missing catalog paths include section, entry name, referenced path, and suggested action. |
| AC-10 | `project-index.md` is treated as a catalog only, not a sprint plan. |
| AC-11 | Project commands do not import or call study internals. |
| AC-12 | `internal/project` may use workspace and generic platform helpers, while platform packages do not import product modules. |
| AC-13 | Project command output is deterministic, concise, and script-friendly in text mode. |
| AC-14 | Command failures preserve cause chains internally and render safe, actionable user-facing errors. |
| AC-15 | Project validation does not invoke agentwrap, OpenCode, runtime health checks, network calls, source cloning, study execution, or sprint flow execution. |
| AC-16 | Tests cover valid fixture, missing roadmap/index, malformed rows, unresolved references, path escapes, empty docs, and ambiguous or missing references. |
| AC-17 | Top-level and project-specific help document the project command family. |
| AC-18 | No global catalog, validation, reports, prompts, planning, or workflow-engine package is introduced. |
| AC-19 | `go test ./...` passes from `/home/antonioborgerees/coding/ultraplan-go`. |
| AC-20 | `go test -race ./...` passes from `/home/antonioborgerees/coding/ultraplan-go`. |
| AC-21 | `go build ./cmd/ultraplan` passes from `/home/antonioborgerees/coding/ultraplan-go`. |

## Contracts Applied

| Contract / Requirement ID | Constraint | Decision Impact | Expected Evidence |
| --- | --- | --- | --- |
| Architecture | `internal/project` owns project docs, project-index parsing, catalog validation, and status; app stays thin; platform remains product-agnostic. | Drives one concrete `internal/project` package and service API, with command wiring in `internal/app` only. | Import review; package placement review; tests proving commands delegate to project service; no platform import of product modules. |
| Errors | Validation and command failures preserve cause chains, classify validation failures, and render safe actionable diagnostics. | Drives structured `ValidationFinding` and app exit-code mapping to `5` for validation failures. | Unit tests for wrapped errors and findings; command tests for stderr, exit code `5`, and safe messages. |
| Security | Workspace-managed paths must be normalized; catalog path escapes fail closed; diagnostics must avoid unsafe leakage. | Drives workspace-safe path resolution in store/validation and explicit external-entry modeling. | Tests for `..` path escape, absolute path escape, unresolved paths, redaction/safe diagnostic rendering, and no runtime/network calls. |
| Testing | Deterministic unit, fixture, and command tests are required. | Drives package tests for discovery/parser/validation and command tests for list/status/validate/help/exit codes. | `internal/project/project_test.go`, `internal/app/project_commands_test.go`, `go test ./...`, `go test -race ./...`. |
| Documentation | Package docs and CLI help must explain ownership and catalog-only semantics. | Drives `internal/project/doc.go` and help text for `project`, `list`, `status`, and `validate`. | Doc review; help command tests; top-level help includes project command family. |
| CLI Surface | Commands must be stable, script-friendly, deterministic, and documented with meaningful exit codes. | Drives text-mode output policy, no interactive UX, and thin command handlers. | Command tests for stdout/stderr, help text, deterministic ordering, and exit mapping. |
| Performance | Project list/status/validate must avoid unbounded recursive scans. | Drives direct-child project/sprint scans, `docs/*.md` only, and catalog-reference checks only. | Tests with hidden/nested/irrelevant entries; review confirms no recursive source/evidence tree traversal. |
| AC-01, AC-02, AC-03 | Discovery and reference resolution must be deterministic, scoped, and filesystem-safe. | Drives `discovery.go` decisions. | Discovery tests for direct children, hidden entries, files, nested dirs, invalid names, sorted names, missing/ambiguous references. |
| AC-04, AC-05, AC-06, AC-09, AC-14 | Status and validation must produce actionable results and correct exit codes. | Drives status summaries, validation findings, and CLI rendering. | Command tests for status summaries, validate success, validate failure, stderr guidance, and exit code `5`. |
| AC-07, AC-08, AC-10 | Project index is a catalog with specific recognized tables and safe path rules. | Drives narrow parser and catalog-only validation. | Parser tests based on real project-index format; validation tests for missing/escaping paths and external entries. |
| AC-11, AC-12, AC-15, AC-18 | Project commands must not depend on study internals, runtime, agentwrap, workflow execution, or global shared abstractions. | Drives dependency boundaries and non-goal exclusions. | Import review, absence checks in review, tests with no runtime fake needed. |
| AC-16, AC-17, AC-19, AC-20, AC-21 | Test and verification coverage must be explicit and complete. | Drives required test files and final verification commands. | `go test ./...`, `go test -race ./...`, `go build ./cmd/ultraplan`, and `review.md` evidence after implementation. |

## Repos Studied / Source Evidence Used

| Source / Repo / Report | Concrete Reference | Relevant Finding | Why It Matters For This Sprint | Used In Decision(s) |
| --- | --- | --- | --- | --- |
| `01-project-structure` report | `studies/go-cli-study/reports/final/01-project-structure.md`; cited `cmd/gh/main.go:6`, `cmd/root.go:112-127`, `internal/restic/repository.go:18`, `cmd/evaluate_sequence_command.go:7` | Non-trivial Go CLIs keep entrypoints thin and domain behavior in protected/domain-owned packages. | Supports `internal/project` ownership and rejection of global catalog/validation layers. | Decisions 1, 5 |
| `02-command-architecture` report | `studies/go-cli-study/reports/final/02-command-architecture.md`; cited `pkg/cmd/install.go:132-145`, `pkg/cmdutil/factory.go:16-43`, `internal/cmd/config.go:1833-1845` | Command constructors parse, delegate, and render instead of owning business rules. | Supports app-level command wiring that delegates to a project service. | Decisions 1, 5 |
| `03-dependency-injection` report | `studies/go-cli-study/reports/final/03-dependency-injection.md`; cited `pkg/cmd/factory/default.go:26-46`, `executor.go:22-24`, `fs/config.go:793` | Manual constructor/factory wiring and narrow interfaces dominate; DI frameworks were not used. | Supports simple service/store construction and test seams without global state. | Decisions 1, 5, 6 |
| `05-error-handling` report | `studies/go-cli-study/reports/final/05-error-handling.md`; cited `cmd/age/tui.go:47-54`, `go-task/errors/errors_task.go:13-32`, `internal/ghcmd/cmd.go:44-49` | Strong CLIs wrap errors, preserve classification, render hints, and map exit codes deliberately. | Supports `ValidationFinding`, preserved cause chains, safe user errors, and exit code `5`. | Decisions 4, 5 |
| `06-io-abstraction` report | `studies/go-cli-study/reports/final/06-io-abstraction.md`; cited `pkg/iostreams/iostreams.go:551-568`, `internal/chezmoi/system.go:25`, `internal/ui/mock.go:10-53` | Injectable streams and filesystem seams make CLI and filesystem behavior testable. | Supports command tests and project store seams without hardcoded stdout/stderr or uncontrolled filesystem calls. | Decisions 3, 5, 6 |
| `09-terminal-ux` report | `studies/go-cli-study/reports/final/09-terminal-ux.md`; cited `internal/cmd/prompt.go:124-137`, `pkg/iostreams/iostreams.go:116-130`, `cmd/root.go:213-215` | CLI-first tools should favor concise, non-TTY-safe, scriptable output for non-interactive commands. | Supports deterministic plain text output and no prompts/spinners/color by default. | Decision 5 |
| `11-testing-strategy` report | `studies/go-cli-study/reports/final/11-testing-strategy.md`; cited `internal/cmd/main_test.go:64-174`, `pkg/httpmock/stub.go:35-199`, `cmd/bisync/bisync_test.go:1435-1479` | Durable CLI behavior is protected by table-driven tests, command-level tests, fixtures, fakes, and golden assertions. | Supports package and command tests for discovery, parser, validation, help, output, exit codes, and no runtime behavior. | Decision 6 |
| `13-security` report | `studies/go-cli-study/reports/final/13-security.md`; cited `internal/config/json/validator.go:146`, `pkg/registry/transport.go:37-41`, `internal/llm/tools/bash.go:41-55`, `internal/common/ignore.go:16-37` | Paths, credentials, subprocesses, and config-like inputs require explicit trust boundaries, canonicalization, validation, and redaction. | Supports workspace-safe catalog resolution, fail-closed path escape behavior, and safe diagnostics. | Decisions 3, 4, 5 |
| `14-performance` report | `studies/go-cli-study/reports/final/14-performance.md`; cited `pkg/cmdutil/factory.go:27-42`, `gdu/pkg/analyze/parallel.go:13`, `yq/pkg/yqlib/stream_evaluator.go:78-113` | Fast CLIs defer expensive initialization and avoid unbounded recursive processing. | Supports direct-child scans and catalog-reference validation only, not recursive source/evidence scans. | Decisions 2, 3 |
| `15-philosophy` report | `studies/go-cli-study/reports/final/15-philosophy.md`; cited `age.go:18`, `VISION.md:5`, `VISION.md:97`, `AGENTS.md:88` | Coherent tools reject features outside purpose and avoid accidental extensibility. | Supports deferring sprint flow, runtime planning, shared catalog frameworks, plugins, TUI, and Git mutation. | Decisions 1, 4, 5 |

## Trade-Off And Debt Analysis

### Accepted Trade-Offs

| Trade-Off | Benefit | Cost / Constraint Accepted | Why Acceptable Now | Revisit Trigger |
| --- | --- | --- | --- | --- |
| Keep project-index parsing and validation inside `internal/project` instead of shared catalog/validation packages. | Preserves module ownership and avoids speculative abstractions. | Similar parsing may need to be repeated or extracted later when `internal/sprint` validates sprint artifacts. | Sprint 16 only validates project catalogs; the concrete table shapes are known from the current project index. | Revisit after at least two product modules implement stable, repeated catalog mechanics with no product semantics. |
| Implement a narrow Markdown table parser for known project-index sections. | Keeps behavior small and aligned with current catalog requirements. | It will not parse arbitrary Markdown tables or unknown catalog families. | Acceptance criteria name specific sections only, and `project-index.md` is a catalog, not a general Markdown database. | Revisit if project indexes add new catalog sections that are requirements-backed. |
| Use structured validation findings rather than a single error string. | Enables actionable diagnostics, sorting, exit code mapping, and reviewable tests. | Adds domain model and rendering work. | Requirements explicitly require section/name/path/action diagnostics and preserved cause chains. | Revisit only if the finding model becomes too generic or duplicates a later stable diagnostic type. |
| Fail closed on missing or workspace-escaping catalog paths. | Protects trust boundaries and prevents invalid catalogs from silently passing. | Users may need to fix legacy or partially written indexes before validation succeeds. | Sprint 16 validates planning roots, so false success is worse than explicit repair work. | Revisit only for explicit external entry rules or migration support. |
| Keep CLI output plain text, deterministic, and non-interactive. | Makes commands script-friendly and testable. | Output is less rich than TTY-specific UX. | Commands are list/status/validate operations and must work in CI. | Revisit if a future TUI or JSON status mode is explicitly scoped. |
| Limit status/list scans to direct project roots, docs, sprint directories, and catalog references. | Keeps project commands fast and avoids walking source repositories or evidence trees. | Status cannot infer every deep artifact issue. | Detailed sprint flow and study validation are out of scope. | Revisit if future project status requirements need deeper but bounded checks. |

### Potential Technical Debt

| Debt / Shortcut | Why It Might Accrue | Current Mitigation | Owner / Follow-Up |
| --- | --- | --- | --- |
| Local table parsing may be duplicated by future sprint validators. | Sprint validators will likely parse selected-context tables. | Keep parser narrow, product-owned, and easy to extract only after repeated concrete use. | Sprint 17 or later planning validators. |
| Text output could become semi-stable before JSON output exists. | Automation may start depending on text mode. | Keep text concise and deterministic; do not promise JSON for Sprint 16 unless implementation already has a supported output mode. | Future JSON output requirements from TRD `--json` should define schema explicitly. |
| Project status may expose a simplified catalog health summary. | Full validation may be more detailed than status. | Make status derive catalog health from the same validation service where practical and document summary semantics in tests. | Sprint 16 implementation and review. |
| External catalog entry modeling may start minimal. | The current project index appears workspace-relative, but requirements allow explicitly external entries. | Model external entries explicitly enough to skip workspace path checks only when clearly marked external; reject ambiguous path escapes. | Future project-index schema changes. |
| Command tests may need fixture helpers that resemble product builders. | Fixtures for missing files, malformed rows, and path escapes can become verbose. | Keep helpers test-local and behavior-focused, not exported product APIs. | Test maintenance during Sprint 16 implementation. |

### Future Considerations

| Consideration | Deferred Until | Reason Deferred | What Should Be Preserved Now |
| --- | --- | --- | --- |
| Sprint artifact model, `flow-state.json`, and planning-stage validators. | Sprint 17 or later. | Explicit Sprint 16 non-goal. | Keep `internal/project` service usable by future `internal/sprint` without depending on study internals or runtime. |
| Runtime-backed prompt rendering or planning generation. | Sprint planning flow sprint. | Project validation must not invoke agentwrap, OpenCode, prompt generation, or network calls. | Keep platform runtime product-agnostic and unused by `project` commands. |
| Shared catalog/table parser. | After repeated concrete need across `project` and `sprint`. | Premature shared package is prohibited by requirements and Architecture. | Keep parsing functions small, tested, and product-context names explicit. |
| JSON output schema for project status/validate. | Future CLI surface stabilization. | Sprint requirements only require deterministic script-friendly text mode. | Keep service result structs structured enough that JSON can be added without moving business rules to `internal/app`. |
| Migration support for older project indexes. | Later compatibility sprint. | No persisted external compatibility requirement is active for Sprint 16. | Produce clear diagnostics rather than silent compatibility shims. |

## Final Decisions

### Decision 1: Own Project Behavior In `internal/project`

- **Decision:** Implement `internal/project` as the owner of project domain types, package documentation, project discovery, filesystem store, project-index parsing, catalog validation, status behavior, and use-case service API. `internal/app` will only register commands, parse arguments, call the service, render output, and map exit codes.
- **Rationale:** The sprint goal is a project planning root, not a generic catalog framework. TRD section 18.1 and Architecture require `internal/project` to own project docs, `project-index.md` parsing, catalog entries, validation, and status. Keeping behavior local prevents `internal/app` and platform packages from accumulating project rules.
- **Study / Source Grounding:** `technical-handbook.md` patterns for module-owned behavior, concrete project package before shared abstractions, factory-created commands, and manual DI. Evidence reports `01-project-structure`, `02-command-architecture`, `03-dependency-injection`, and `15-philosophy` cite thin entrypoints and protected/domain packages, including `cmd/gh/main.go:6`, `cmd/root.go:112-127`, `internal/restic/repository.go:18`, `pkg/cmd/install.go:132-145`, `pkg/cmdutil/factory.go:16-43`, and `age.go:18`.
- **Trade-Offs Accepted:** More files and domain types are created before the commands feel complete, but this keeps business rules out of CLI wiring and avoids global shared packages.
- **Technical Debt / Future Impact:** Future sprint validators may need related parsing mechanics, but extraction is deferred until repeated stable behavior exists. This decision preserves a clean `sprint -> project` future dependency without letting `project` depend on `sprint`.
- **Alternatives Rejected:** Put project logic directly in `internal/app` command handlers, rejected because it violates thin adapter guidance and makes tests couple to CLI internals. Create `internal/catalog` or `internal/validation`, rejected because Sprint 16 has one concrete owner and AC-18 prohibits speculative global packages. Reuse `internal/study` models/services, rejected because AC-11 and TRD section 18.2 prohibit study-semantic reuse.
- **Contracts Satisfied:** Architecture, Documentation, CLI Surface, AC-10, AC-11, AC-12, AC-15, AC-18; required outputs `domain.go`, `doc.go`, `discovery.go`, `store_fs.go`, `index.go`, `validation.go`, `service.go`.
- **Evidence Required:** Import review confirms `internal/project` does not import `internal/study` or runtime packages; platform packages do not import product modules; no global catalog/validation/planning/reports/prompts package appears; package docs state project indexes are catalogs, not sprint plans.

### Decision 2: Discover And Resolve Projects By Direct Safe Names Only

- **Decision:** Project discovery will inspect only direct non-hidden directories under workspace `projects/`, validate project names as filesystem-safe single path segments, sort discovered projects by name, and resolve project references only when exact or unambiguous according to the project service's reference rules.
- **Rationale:** The requirements ask for deterministic project listing, hidden/non-directory/nested exclusion, actionable invalid-name diagnostics, and bounded scans. Direct safe names are the smallest model that satisfies project roots without modeling sprint or study behavior.
- **Study / Source Grounding:** `technical-handbook.md` performance and security patterns support bounded discovery and safe path boundaries. Evidence reports `13-security` and `14-performance` cite schema/path validation and bounded traversal examples, including `internal/config/json/validator.go:146`, `internal/common/ignore.go:16-37`, `pkg/cmdutil/factory.go:27-42`, and `gdu/pkg/analyze/parallel.go:13`.
- **Trade-Offs Accepted:** Discovery will not recursively search for project-like directories or infer projects from nested files. This avoids expensive and surprising behavior.
- **Technical Debt / Future Impact:** If future workspaces support aliases or external project roots, a new explicit resolution rule will be needed. The current model should not add compatibility shims until such requirements exist.
- **Alternatives Rejected:** Recursively scan the workspace for `project-index.md`, rejected because it is slow, may enter source/evidence trees, and violates AC-02 and performance constraints. Accept arbitrary path references as project names, rejected because it creates path escape risks and weakens workspace-root semantics.
- **Contracts Satisfied:** Architecture, Security, Performance, Testing, AC-01, AC-02, AC-03, AC-11, AC-12, AC-13, AC-15.
- **Evidence Required:** Unit tests cover direct child discovery, sorted output, hidden entries ignored, files ignored, nested directories ignored, invalid names rejected, missing references, and ambiguous references; command tests verify deterministic `project list` output.

### Decision 3: Read Project Files Through A Workspace-Safe Filesystem Store

- **Decision:** Implement a project filesystem store that reads `docs/*.md`, `roadmap.md`, `project-index.md`, and direct sprint directory metadata through workspace-safe path helpers. Store methods will normalize paths before filesystem access and return structured status inputs to the service.
- **Rationale:** Project status and validation are read-only operations over workspace-managed artifacts. The filesystem boundary must be explicit enough for tests, path safety, and future sprint reuse, but it must not become a generic persistence framework.
- **Study / Source Grounding:** `technical-handbook.md` IO abstraction, security, and performance patterns support injectable IO/filesystem seams, safe path normalization, and bounded scans. Evidence includes `pkg/iostreams/iostreams.go:551-568`, `internal/chezmoi/system.go:25`, `internal/config/json/validator.go:146`, `internal/common/ignore.go:16-37`, and `yq/pkg/yqlib/stream_evaluator.go:78-113`.
- **Trade-Offs Accepted:** A local store abstraction adds a small amount of structure but keeps tests deterministic and prevents commands from directly walking the filesystem.
- **Technical Debt / Future Impact:** Store methods should stay project-specific. If future project/sprint modules share atomic or path operations, only mechanical helpers should move to platform/workspace.
- **Alternatives Rejected:** Let command handlers call `os.ReadDir` and `os.Stat` directly, rejected because it hides path rules in CLI code and weakens testability. Add a broad platform filesystem repository with project semantics, rejected because platform packages must remain product-agnostic.
- **Contracts Satisfied:** Architecture, Security, Testing, Performance, AC-04, AC-08, AC-11, AC-12, AC-13, AC-15, AC-16.
- **Evidence Required:** Tests cover docs directory present/empty, Markdown doc discovery, missing roadmap, missing project index, direct sprint directory counting, deterministic status inputs, and path normalization before catalog file access; review confirms no recursive source/evidence tree scan.

### Decision 4: Parse Current Project-Index Catalog Tables Narrowly And Validate Findings Structurally

- **Decision:** The project-index parser will recognize the concrete catalog sections used by the current project index: Source Documents, Active Contract Pool, Available Evidence Reports, Available Reasoning Templates, and Review Protocols. Validation will produce sorted structured findings with severity, section, entry name, referenced path, problem/cause, and suggested action. Catalog paths resolve relative to the workspace and fail closed on missing files or workspace escapes unless an entry is explicitly modeled as external.
- **Rationale:** Sprint 16 needs actionable validation for the existing catalog format, not a general Markdown parser. Structured findings directly satisfy missing-path diagnostics and exit-code behavior while preserving internal causes.
- **Study / Source Grounding:** `technical-handbook.md` supports actionable classified validation, fail-closed path boundaries, and concrete package-owned parsing. Evidence reports `05-error-handling` and `13-security` cite `%w` wrapping, hint-based errors, exit classification, schema validation, path canonicalization cautions, and redaction, including `cmd/age/tui.go:47-54`, `go-task/errors/errors_task.go:13-32`, `internal/ghcmd/cmd.go:44-49`, `internal/config/json/validator.go:146`, and `internal/common/ignore.go:16-37`. `15-philosophy` supports rejecting broad abstractions outside current purpose.
- **Trade-Offs Accepted:** Unknown catalog sections will not automatically become validated project behavior. This is acceptable because Sprint 16 acceptance criteria name exact sections and require catalog-only semantics.
- **Technical Debt / Future Impact:** Future project-index schema evolution may require extending recognized sections. The finding model should be stable enough for future JSON output but not over-generalized now.
- **Alternatives Rejected:** Build a generalized Markdown table or catalog framework, rejected because it is out of scope and AC-18 prohibits a global catalog package. Treat malformed rows or missing paths as warnings only, rejected because AC-06 and AC-08 require validation failure. Infer sprint plans from catalog entries, rejected because AC-10 requires catalog-only behavior.
- **Contracts Satisfied:** Architecture, Errors, Security, Testing, Documentation, AC-05, AC-06, AC-07, AC-08, AC-09, AC-10, AC-14, AC-15, AC-16, AC-18.
- **Evidence Required:** Parser tests based on `projects/ultraplan-go/project-index.md`; tests for malformed rows, unresolved references, path escapes, explicit external entries if supported, missing catalog section/name/path/action fields, sorted findings, and validation exit code `5` via command tests.

### Decision 5: Expose Project Commands As Thin, Deterministic, Runtime-Free CLI Wiring

- **Decision:** Add `ultraplan project list`, `ultraplan project <project> status`, and `ultraplan project <project> validate` in `internal/app/project_commands.go`, register the command family in `internal/app/app.go`, and keep handlers limited to argument parsing, service calls, deterministic text rendering, and exit-code mapping. Text output will be concise, stable, non-interactive, and free of ANSI/color/prompt/progress behavior.
- **Rationale:** The project commands are inspection and validation commands. They must be script-friendly, deterministic, and isolated from runtime/study behavior while following the established CLI surface.
- **Study / Source Grounding:** `technical-handbook.md` command, terminal UX, error, IO, and DI patterns. Evidence reports cite command factories and thin delegates in `pkg/cmd/install.go:132-145`, `pkg/cmd/issue/list/list.go:47`, `pkg/cmdutil/factory.go:16-43`, output capture in `pkg/iostreams/iostreams.go:551-568`, non-TTY fallback in `internal/cmd/prompt.go:124-137`, and exit classification in `internal/ghcmd/cmd.go:44-49`.
- **Trade-Offs Accepted:** CLI rendering will be intentionally plain. JSON output is not made a Sprint 16 architecture dependency unless existing app conventions already provide it, because the sprint acceptance criteria emphasize deterministic text mode.
- **Technical Debt / Future Impact:** Automation may rely on text output before JSON is added. Keeping service results structured will allow later JSON support without moving business rules into command handlers.
- **Alternatives Rejected:** Add interactive repair prompts or rich terminal status, rejected because project validation is read-only and AC-13 requires script-friendly text. Call runtime health checks or agentwrap during validation, rejected because AC-15 explicitly forbids runtime/network behavior. Put exit-code classification inside `internal/project`, rejected because command exit mapping belongs in `internal/app` while domain findings belong in `internal/project`.
- **Contracts Satisfied:** CLI Surface, Errors, Documentation, Testing, Performance, AC-01, AC-04, AC-05, AC-06, AC-13, AC-14, AC-15, AC-17.
- **Evidence Required:** Command tests cover top-level help, `project --help`, `project list --help`, `project <project> status --help`, `project <project> validate --help`, list/status/validate stdout and stderr, exit code `0` for valid validation, exit code `5` for validation failures, deterministic output ordering, no ANSI escape sequences, and no runtime invocation.

### Decision 6: Verification Will Be Fixture-First, Command-Level, And Review-Protocol Driven

- **Decision:** Add `internal/project/project_test.go` for domain/discovery/parser/validation/status behavior and `internal/app/project_commands_test.go` for CLI help, output, exit codes, deterministic ordering, and runtime-free behavior. Sprint review evidence will be recorded in `projects/ultraplan-go/sprints/16-project-domain-and-index/review.md` after implementation.
- **Rationale:** Acceptance criteria include both domain behavior and CLI behavior. Fixture-first tests avoid OpenCode, credentials, network, Git, long sleeps, and source repository scans.
- **Study / Source Grounding:** `technical-handbook.md` testing strategy and IO abstraction patterns. Evidence report `11-testing-strategy` cites command-level/testscript integration, mocks, fixture filesystem isolation, and golden assertions including `internal/cmd/main_test.go:64-174`, `pkg/httpmock/stub.go:35-199`, `internal/chezmoitest/chezmoitest.go:86-92`, and `helm/internal/test/test.go:43`. `06-io-abstraction` supports captured streams and filesystem seams.
- **Trade-Offs Accepted:** Command tests and fixtures require maintenance, but they protect user-visible contracts and acceptance criteria more directly than private helper tests alone.
- **Technical Debt / Future Impact:** Fixture helpers should remain test-local to avoid creating a hidden test framework. Future JSON output can add golden/schema tests when scoped.
- **Alternatives Rejected:** Only run package unit tests, rejected because help, stdout/stderr, and exit codes are command-level behavior. Use real workspace or runtime dependencies in tests, rejected because AC-15 and test constraints require offline deterministic tests. Defer race/build verification, rejected because AC-20 and AC-21 are explicit acceptance criteria.
- **Contracts Satisfied:** Testing, CLI Surface, Errors, Documentation, Security, Performance, AC-16, AC-17, AC-19, AC-20, AC-21.
- **Evidence Required:** Passing `go test ./...`, passing `go test -race ./...`, passing `go build ./cmd/ultraplan`, review evidence covering Architecture Review and Sprint Review protocols, and tests for all fixtures named in AC-16.

## Expected Evidence

| Evidence Type | Required Evidence | Source / Command / Review Check |
| --- | --- | --- |
| Tests | Discovery, reference resolution, filesystem-safe name validation, direct-child filtering, deterministic sorting, status summary, parser categories, malformed catalog rows, missing roadmap, missing project index, unresolved references, workspace escapes, empty docs, missing/ambiguous project references. | `/home/antonioborgerees/coding/ultraplan-go/internal/project/project_test.go`; `go test ./...`. |
| Tests | Command help, top-level help registration, list/status/validate output, stdout/stderr separation, exit code `0`, exit code `5`, deterministic ordering, no ANSI output, runtime-free behavior. | `/home/antonioborgerees/coding/ultraplan-go/internal/app/project_commands_test.go`; `go test ./...`. |
| Verification | Offline unit and command test suite passes. | `go test ./...` from `/home/antonioborgerees/coding/ultraplan-go`. |
| Verification | Race verification passes. | `go test -race ./...` from `/home/antonioborgerees/coding/ultraplan-go`. |
| Verification | CLI binary builds. | `go build ./cmd/ultraplan` from `/home/antonioborgerees/coding/ultraplan-go`. |
| Review | Architecture Review confirms package ownership, dependency direction, no study/runtime coupling, no global shared packages, and platform product-agnosticism. | `system/protocols/architecture-review-protocol.md`; record result in sprint `review.md`. |
| Review | Sprint Review confirms required outputs, acceptance criteria, non-goal exclusions, tests, commands, diagnostics, and verification results. | `system/protocols/sprint-review-protocol.md`; `projects/ultraplan-go/sprints/16-project-domain-and-index/review.md`. |
| Documentation | Package documentation states `internal/project` ownership and catalog-only semantics; CLI help documents `project`, `list`, `status`, and `validate`. | `/home/antonioborgerees/coding/ultraplan-go/internal/project/doc.go`; `/home/antonioborgerees/coding/ultraplan-go/internal/app/app.go`; command help tests. |

## Assumptions And Risks

| Item | Type | Impact | Mitigation / Follow-Up |
| --- | --- | --- | --- |
| The current `project-index.md` table shapes are the Sprint 16 parser contract. | Assumption | Parser may reject future table variants not represented in the current index. | Keep parser narrow but well tested; extend only when a future requirement updates the catalog shape. |
| Existing workspace path helpers are sufficient for workspace-safe project and catalog path resolution. | Assumption | If helpers are missing or study-specific, implementation may need small workspace/platform helper additions. | Add only generic path helpers to workspace/platform; do not move product rules out of `internal/project`. |
| Project reference resolution can be exact plus unambiguous prefix or equivalent existing CLI convention. | Assumption | Ambiguity behavior must be clear for users and tests. | Implement explicit missing/ambiguous diagnostics and deterministic candidate ordering. |
| External catalog entries may not exist in current fixtures. | Risk | External modeling could be underspecified and path validation could either over-accept or over-reject. | Accept external only when explicitly marked by modeled fields or URL-like path; otherwise fail closed. Add tests for at least one explicit external case if implementation supports it. |
| Status output could duplicate validation logic. | Risk | Divergent status and validate behavior would confuse users. | Derive catalog health from shared validation service or common internal helper; test status against valid and invalid fixtures. |
| Text output may become automation-sensitive before JSON exists. | Risk | Future output changes may break scripts. | Keep text concise, sorted, and documented by command tests; plan future JSON schema separately. |
| Adding project commands could accidentally touch runtime or study imports through existing app composition. | Risk | Violates non-goals and dependency boundaries. | Add import review in sprint review; keep service constructors project/workspace/store only; no runtime fakes in project command tests. |
| Validation diagnostics may leak absolute local paths. | Risk | Unsafe or noisy output. | Render workspace-relative paths where possible; preserve causes internally; include safe user-facing path and corrective action. |

## Implementation Constraints

- `internal/project` must own project discovery, docs discovery, roadmap/index discovery, project-index parsing, catalog validation, status result construction, and project service use cases.
- `internal/project` must not import `internal/study`, study-owned source/dimension/report validators, rating parsing, summary generation, run-loop scheduling, `internal/sprint`, `internal/platform/runtime`, agentwrap, or OpenCode.
- `internal/app` must stay a thin adapter: parse args, call project service, render deterministic text, and map exit codes.
- Platform packages must not import `internal/project` or any product package.
- Do not introduce `internal/catalog`, `internal/validation`, `internal/reports`, `internal/prompts`, `internal/planning`, or a workflow-engine package.
- Discovery must scan only direct children under `projects/`; status must scan `docs/*.md`, direct sprint directories, required project files, and catalog-referenced paths only.
- Project names must be filesystem-safe path segments; invalid, missing, and ambiguous references must fail with actionable diagnostics.
- Catalog paths must be normalized through workspace-safe helpers before filesystem access and must fail closed on workspace escapes unless explicitly modeled as external.
- Validation findings must include section, entry name, referenced path when applicable, problem/cause, and suggested corrective action.
- Command output must be deterministic, concise, script-friendly, text-only by default, and free of ANSI, prompts, spinners, runtime health checks, network calls, source cloning, study execution, and sprint flow execution.
- Tests must be deterministic, offline, and fake/fixture-first; default verification must not require OpenCode, provider credentials, network access, Git commands, or long sleeps.
- Implementation must not mutate `project-index.md`, `roadmap.md`, docs, sprint directories, generated planning artifacts, unrelated dirty files, or Git state.

## Plan Handoff

`plan.md` must execute these decisions. It must not invent architecture, scope, or decisions beyond this document.

The plan must carry forward:

- final decisions in this document
- selected contracts and AC-01 through AC-21 mappings
- required outputs from `requirements.md`
- expected tests and verification commands
- risks and mitigations
- Architecture Review and Sprint Review protocols
- non-goal exclusions, especially study/runtime/sprint-flow/Git mutation exclusions

## Phase Exit Criteria

- [x] Selected context was read and used.
- [x] Area-specific reasoning documents were completed where applicable or explicitly marked not present.
- [x] Area-specific reasoning conclusions are reflected or explicitly replaced by final sprint decisions.
- [x] Contracts are explicitly mapped to decisions and expected evidence.
- [x] Final decisions are clear enough for `plan.md` to execute.
- [x] Expected evidence is specific and reviewable.
