# Sprint Reasoning: Workspace, Config, Logging, and Health Skeleton

> Project: `ultraplan-go`
> Sprint: `02-workspace-config-health`
> Output: `projects/ultraplan-go/sprints/02-workspace-config-health/reasoning.md`
> **Inputs Used:** `projects/ultraplan-go/sprints/02-workspace-config-health/requirements.md`, `projects/ultraplan-go/sprints/02-workspace-config-health/reasoning.md`, `projects/ultraplan-go/sprints/02-workspace-config-health/sprint-index.md`, `projects/ultraplan-go/sprints/02-workspace-config-health/technical-handbook.md`, `projects/ultraplan-go/docs/PRD.md`, `projects/ultraplan-go/docs/TRD.md`, `projects/ultraplan-go/docs/ARCHITECTURE.md`, `projects/ultraplan-go/project-index.md`, `templates/sprint-reasoning.md`; area reasoning omitted because `projects/ultraplan-go/sprints/02-workspace-config-health/reasoning/` was not present.

This document decides. It synthesizes selected context, handbook evidence, available reasoning inputs, and selected contracts into final sprint decisions.

It does not execute implementation and does not replace `sprint-index.md`, `technical-handbook.md`, or any area-specific reasoning documents.

## Sprint Purpose

- **Goal:** Make UltraPlan able to find, initialize, inspect, and validate a workspace by implementing workspace discovery, workspace initialization, workspace path safety, structural validation, effective config loading, secret redaction, deterministic command output, minimal logging, and `init-workspace`, `config show`, and `health` command skeletons.
- **Non-Goals:** No study listing, study initialization, source discovery, analysis runs, synthesis, summary generation, code extraction, runtime execution, agentwrap/OpenCode execution or health checks, provider/model availability checks, durable run state, durable file logs, event persistence, target scaffolding, sprint planning, or sprint execution behavior.
- **Depends On:** Sprint 1 CLI foundation for Go module layout, `cmd/ultraplan`, `internal/app`, command routing baseline, and initial package layout.

## Requirement IDs Used By This Reasoning

The sprint requirements file does not define stable IDs, so this reasoning uses the following local IDs for traceability:

| ID | Requirement |
| --- | --- |
| AC-01 | `go test ./...` succeeds from `../ultraplan-go` without OpenCode, provider credentials, network access, or a configured workspace. |
| AC-02 | `go build ./cmd/ultraplan` succeeds and produces a runnable CLI binary. |
| AC-03 | `ultraplan init-workspace [--path <dir>] [--dry-run]` creates the required workspace structure. |
| AC-04 | `ultraplan config show [--json]` prints effective config with precedence and redacted sensitive values. |
| AC-05 | `ultraplan health [--json]` reports basic workspace, config, filesystem, and environment status without OpenCode or provider credentials. |
| AC-06 | Config precedence is built-in defaults < workspace config < environment variables < CLI flags. |
| AC-07 | Secret redaction is implemented for sensitive config values. |
| AC-08 | Text and JSON output modes are deterministic for config and health commands. |
| AC-09 | CLI help output includes `init-workspace`, `config`, and `health` commands. |
| AC-10 | Exit codes are deterministic and meaningful per error class. |
| C-01 | Use Go and the standard Go toolchain. |
| C-02 | Workspace discovery precedence is explicit `--workspace`, `ULTRAPLAN_WORKSPACE`, current directory, nearest parent. |
| C-03 | Config loading precedence is defaults, workspace config, environment variables, CLI flags. |
| C-04 | Secrets must be redacted and never logged or exposed. |
| C-05 | Keep `cmd/ultraplan/main.go` thin; logic lives under `internal/`. |
| C-06 | Platform packages must not import product modules. |
| C-07 | Tests must be deterministic and offline. |
| C-08 | CLI output must be script-friendly with deterministic text, JSON mode, and meaningful exit codes. |

## Selected Context And Pre-Reasoning Artifacts

| Artifact | Path | How It Was Used |
| --- | --- | --- |
| Project Index | `projects/ultraplan-go/project-index.md` | Established project scope, active contract pool, available evidence reports, and the fact that no prior project decisions exist. |
| Product Requirements | `projects/ultraplan-go/docs/PRD.md` | Confirmed local-first CLI product goals, workspace/config/health command expectations, script-friendly UX, redaction, and target/sprint workflow deferral. |
| Technical Requirements | `projects/ultraplan-go/docs/TRD.md` | Provided workspace discovery precedence, required workspace structure, config schema, config precedence, validation needs, exit-code classes, output mode requirements, logging/diagnostic constraints, testing constraints, and current-scope clarification. |
| Architecture | `projects/ultraplan-go/docs/ARCHITECTURE.md` | Grounded package ownership and dependency direction: `workspace` owns workspace behavior, `platform/config` and `platform/logging` own cross-cutting infrastructure, `internal/app` composes CLI behavior, and `cmd/ultraplan` stays thin. |
| Sprint Index | `projects/ultraplan-go/sprints/02-workspace-config-health/sprint-index.md` | Selected contracts, evidence reports, review protocols, scope, non-goals, and excluded context for this sprint. |
| Technical Handbook | `projects/ultraplan-go/sprints/02-workspace-config-health/technical-handbook.md` | Supplied source-backed implementation patterns, warnings, design pressures, and open questions for workspace/config/output/logging/health decisions. |
| Sprint Reasoning Template | `templates/sprint-reasoning.md` | Supplied required reasoning structure and quality bar. |

## Area-Specific Reasoning Inputs

No area-specific reasoning files were present at `projects/ultraplan-go/sprints/02-workspace-config-health/reasoning/*.md` when this document was created. The sprint index selected an Architecture reasoning artifact at `projects/ultraplan-go/sprints/02-workspace-config-health/reasoning/architecture.md`, but the reasoning directory did not exist. Because there are no actual area reasoning conclusions to summarize, this consolidated reasoning makes the final architecture, package ownership, dependency direction, and composition decisions directly from the project docs, selected contracts, and technical handbook.

| Area | Reasoning Document | Key Conclusion | Evidence Basis | Impact On Final Decision |
| --- | --- | --- | --- | --- |
| Architecture | Not present | No area-specific conclusion was available. Final package-boundary decisions are made in this document. | `ARCHITECTURE.md`, `TRD.md`, `sprint-index.md`, `technical-handbook.md`, and selected go-cli-study reports. | `plan.md` must execute the package-boundary decisions in this document and must not reopen architecture because the area artifact is absent. |

## Sprint Technical Handbook Summary

- **Relevant Patterns:** Thin `cmd/` entrypoint with `internal/` business packages; command factories and shared dependency construction; manual composition root; fixed layered config precedence with explicit changed-flag handling; post-merge config validation; injected IO for deterministic command tests; classified errors mapped to exit codes; redaction by type or projection; structured logging with stdout/stderr separation; behavior-focused command and unit tests.
- **Important Trade-Offs:** Prefer `internal/` packages over public `pkg/`; use a small shared app dependency container without letting it become a god object; use fixed config precedence rather than dynamic config-source priority; use projection-based redaction now while preserving a sensitive-field boundary; use command-level output assertions where user-visible output matters; build minimal logging now and defer durable observability.
- **Warnings / Anti-Patterns:** Do not put workspace/config logic in `cmd/ultraplan`; do not let flag defaults shadow config; do not bypass centralized config with scattered env reads; do not create global mutable config or logging state; do not print secrets through generic formatting; do not mix stdout data and stderr diagnostics; do not hardcode `os.Stdout`, `os.Stderr`, or filesystem behavior where tests need isolation; do not panic for user-fixable workspace/config failures.
- **Evidence Confidence:** High for project structure, command architecture, dependency injection, configuration management, error handling, IO abstraction, testing, and security because the handbook cites multiple mature Go CLIs and concrete source references. Medium for terminal UX and logging because those reports are still applicable but this sprint intentionally implements only skeleton behavior.

## Contracts Applied

| Contract / Requirement ID | Constraint | Decision Impact | Expected Evidence |
| --- | --- | --- | --- |
| Architecture, C-05, C-06 | Keep entrypoint thin; product modules own product behavior; platform packages do not import product modules. | `cmd/ultraplan` delegates to `internal/app`; `internal/workspace` owns workspace behavior; `internal/platform/config` and `internal/platform/logging` stay product-agnostic; output helpers live in `internal/app`. | Architecture review, import direction inspection, `go test ./...`, `go build ./cmd/ultraplan`. |
| Errors, AC-10, C-08 | User-facing failures must be classified with deterministic exit codes and actionable messages. | Define error classes for usage, config, workspace/filesystem, validation, and general failures; command wiring maps classes to TRD suggested exit codes. | Command tests for invalid args, invalid config, missing workspace, path escape, validation failures, and filesystem errors. |
| Configuration, AC-04, AC-06, AC-07, C-03, C-04 | Resolve effective config in required precedence and redact sensitive values. | `platform/config` owns defaults, workspace config loading, env override parsing, CLI override application, validation, source metadata where practical, and redacted projections. | Unit tests for each precedence layer, changed-flag handling, validation, and redaction; command tests for `config show` text/JSON. |
| Observability, C-04, C-08 | Establish safe logging/diagnostic foundation without durable event persistence. | Implement minimal text/JSON logger that writes diagnostics to stderr, supports levels, accepts structured fields, and redacts secret fields. | Unit or command tests proving data output stays stdout, diagnostics stay stderr, JSON logs are valid when enabled, and secrets are absent. |
| Security, AC-03, AC-07, C-02, C-04 | Normalize workspace-managed paths, reject workspace escapes, and never expose secrets. | `workspace` owns canonical path helpers and escape checks; config/logging/output use redaction before display; init writes only inside selected path. | Path tests for traversal/symlink-normalized escape attempts where feasible, init command tests, redaction tests. |
| Testing, AC-01, C-07 | Tests must be deterministic, offline, and independent of real workspace/provider/OpenCode. | Use temp directories, injected IO, table-driven tests, and command-level tests; do not require network or provider credentials. | `go test ./...` from `../ultraplan-go`. |
| Documentation, AC-09 | Help output must expose new commands and flags. | Command registration must include `init-workspace`, `config show`, `health`, `--workspace`, `--json`, `--dry-run`, and command help text. | Help-output command tests or review check. |
| CLI Surface, AC-03, AC-04, AC-05, AC-08, AC-09, AC-10 | Commands must be script-friendly, deterministic, support JSON where required, and return meaningful exit codes. | Shared output helpers render deterministic text/JSON envelopes; command tests assert stable outputs and exit codes. | Command-level tests and manual review commands. |

## Repos Studied / Source Evidence Used

| Source / Repo / Report | Concrete Reference | Relevant Finding | Why It Matters For This Sprint | Used In Decision(s) |
| --- | --- | --- | --- | --- |
| `go-cli-study` report `01-project-structure` | `studies/go-cli-study/reports/final/01-project-structure.md`; cited `cmd/restic/main.go:37-114`, `internal/restic/repository.go:18`, `internal/chezmoi/chezmoi.go:1-2` | Mature Go CLIs keep entrypoints thin and route business logic into `internal/` packages. | Directly supports keeping workspace/config/health behavior out of `cmd/ultraplan`. | Decisions 1, 8 |
| `go-cli-study` report `02-command-architecture` | `studies/go-cli-study/reports/final/02-command-architecture.md`; cited `pkg/cmd/install.go:132-145`, `pkg/cmdutil/factory.go:16-43`, `internal/cmd/config.go:2236-2454` | Command factories and shared dependency/context objects keep commands thin while sharing setup. | Supports `internal/app` command wiring with shared workspace/config/output/logging dependencies. | Decisions 1, 8 |
| `go-cli-study` report `03-dependency-injection` | `studies/go-cli-study/reports/final/03-dependency-injection.md`; cited `internal/cmd/config.go:362`, `pkg/cmdutil/factory.go:16-43`, `executor.go:22-24` | Manual composition roots and constructor injection dominate; no DI framework is needed. | Supports explicit app dependencies and test seams without global mutable state. | Decisions 1, 8, 9 |
| `go-cli-study` report `04-configuration-management` | `studies/go-cli-study/reports/final/04-configuration-management.md`; cited `internal/cmd/config.go:2253-2287`, `internal/flags/flags.go:314-327`, `restic/internal/global/global.go:139,147` | Strong CLIs implement explicit layered precedence, flag-change tracking, env prefixes, and post-merge validation. | Directly grounds effective config behavior and CLI override handling. | Decision 4 |
| `go-cli-study` report `05-error-handling` | `studies/go-cli-study/reports/final/05-error-handling.md`; cited `internal/ghcmd/cmd.go:44-49`, `cmd/cmd.go:497-516`, `errors/errors.go:47-50` | Errors are wrapped and classified so scripts get deterministic exit codes. | Supports typed or classified sprint errors for usage/config/workspace/validation failures. | Decision 7 |
| `go-cli-study` report `06-io-abstraction` | `studies/go-cli-study/reports/final/06-io-abstraction.md`; cited `pkg/iostreams/iostreams.go:551-568`, `internal/ui/mock.go:10-53`, `ui_mock.go:27-33` | Injected stdout/stderr/stdin and filesystem seams make command tests deterministic. | Supports output helper design and command-level tests without real terminal dependencies. | Decisions 6, 9 |
| `go-cli-study` report `09-terminal-ux` | `studies/go-cli-study/reports/final/09-terminal-ux.md`; cited `internal/cmd/prompt.go:20-256`, `lib/terminal/terminal.go:82-86` | CLI-first tools preserve scriptability and avoid interactive assumptions unless needed. | Supports deterministic text output and no prompts in workspace/config/health commands. | Decisions 6, 8 |
| `go-cli-study` report `10-logging-observability` | `studies/go-cli-study/reports/final/10-logging-observability.md`; cited `internal/logging/logging.go:31-71`, `internal/slogs/keys.go:6-231`, `cmd/age/tui.go:37-40` | Strong CLIs separate stdout data from stderr diagnostics and use structured level-aware logging. | Supports minimal logging foundation, stderr diagnostics, structured fields, and deferring durable logs. | Decision 5 |
| `go-cli-study` report `11-testing-strategy` | `studies/go-cli-study/reports/final/11-testing-strategy.md`; cited `internal/cmd/main_test.go:64-174`, `helm/internal/test/test.go:43`, `go-task/task_test.go:166-169` | Table-driven unit tests, command tests, fixtures, and golden-like output checks catch CLI regressions. | Supports required test matrix for workspace/config/health behavior and deterministic outputs. | Decision 9 |
| `go-cli-study` report `13-security` | `studies/go-cli-study/reports/final/13-security.md`; cited `internal/options/secret_string.go:15-20`, `pkg/registry/transport.go:37-41`, `internal/config/json/validator.go:146`, `status.go:332-338` | Secret wrappers, redaction projections, credential scrubbing, schema validation, and path canonicalization define CLI trust boundaries. | Supports path escape rejection and redacted config/log output. | Decisions 3, 4, 5 |

## Trade-Off And Debt Analysis

### Accepted Trade-Offs

| Trade-Off | Benefit | Cost / Constraint Accepted | Why Acceptable Now | Revisit Trigger |
| --- | --- | --- | --- | --- |
| Keep workspace as one `internal/workspace` package split by files instead of subpackages. | Clear ownership and low ceremony for discovery, paths, validation, and init. | Package can grow as more workspace behaviors appear. | Sprint scope is small and ARCHITECTURE.md recommends one package per module first. | Workspace package becomes hard to navigate or develops independent subdomains. |
| Put deterministic output helpers in `internal/app` instead of `internal/platform/output`. | Avoids premature platform abstraction for a CLI surface detail. | Future product modules cannot import a generic output package directly. | This sprint only needs command rendering for app-level commands; platform packages should not know product payloads. | Multiple modules need shared rendering outside app command wiring. |
| Use fixed config precedence rather than a dynamic source-priority framework. | Simple, reviewable, and directly matches TRD and sprint requirements. | Less flexible for future config sources or plugin-defined layers. | Current requirements define an exact chain; dynamic priority would add complexity without need. | Requirements add user profiles, XDG config, project overrides, or plugin-provided config layers. |
| Use redacted projection as the command-output boundary, with sensitive field detection centralized in `platform/config`. | Fast to implement and easy to test for `config show`, health diagnostics, and logs. | Plain config values may exist internally and must not be formatted directly. | Scope is small and can enforce redaction through output/logging APIs and tests. | Runtime/provider secret model expands enough to justify dedicated secret wrapper types everywhere. |
| Implement logging skeleton only, not durable observability. | Satisfies safe diagnostics foundation while respecting non-goals. | Later run event persistence may require adapting interfaces or field names. | Runtime execution, run records, and agentwrap observability are explicitly excluded. | A later sprint introduces runtime execution, run loops, event stores, or status history. |
| Command-level tests focus on behavior and output rather than full golden fixture files for every command. | Lower fixture maintenance while still asserting script-facing contracts. | Some formatting regressions may require broader expected strings. | Sprint output is small and can be asserted with exact text/JSON snippets. | Output schemas become larger or documented as stable public contracts. |

### Potential Technical Debt

| Debt / Shortcut | Why It Might Accrue | Current Mitigation | Owner / Follow-Up |
| --- | --- | --- | --- |
| Minimal `ultraplan.yml` schema may omit runtime fields that later execution needs. | This sprint must initialize config but runtime execution is out of scope. | Include TRD-aligned starter fields without validating provider availability; keep config structs extensible but explicit. | Runtime/config sprint should extend validation and agentwrap mapping. |
| App dependency container could become too broad. | Workspace, config, logging, output, command registration, and errors are all composed in `internal/app`. | Keep the container small, concrete, and limited to Sprint 2 dependencies; avoid adding study/runtime state. | Reassess when study commands are added. |
| Projection-only redaction relies on discipline. | Future code may accidentally print raw config structs. | Centralize command output through redacted projections and add tests that search for known secret values in output/logs. | Consider sensitive wrapper types when provider/runtime secrets are added. |
| Health JSON schema may become a public contract before runtime checks exist. | Automation may depend on early skeleton fields. | Use stable top-level envelope and explicit check list names; include status and messages without overclaiming runtime support. | Version or document health schema when runtime health is implemented. |
| No area-specific architecture reasoning artifact exists. | Future readers may expect separate architecture reasoning because sprint-index selected it. | This document records final architecture decisions directly and names the absent artifact. | Create or backfill area reasoning only if project governance requires it, not as an implementation blocker. |

### Future Considerations

| Consideration | Deferred Until | Reason Deferred | What Should Be Preserved Now |
| --- | --- | --- | --- |
| Agentwrap/OpenCode health checks | Runtime integration sprint | Sprint 2 health is limited to workspace/config/filesystem/environment basics. | Health check model should allow adding runtime checks without changing command shape. |
| Durable logs and event records | Run-loop/observability sprint | Current sprint excludes durable file logs and event persistence. | Logging fields should be structured and redaction-safe. |
| Study workflow validation | Study initialization/listing/running sprints | Sprint 2 only validates workspace structure and config basics. | Workspace validation should not absorb study-specific rules. |
| Config migrations | Schema evolution sprint | No persisted config migration requirement exists yet beyond starter `version: 1`. | Config validation should report schema version clearly. |
| Public output schemas | Documentation/release hardening | Sprint needs deterministic output but not formal schema docs. | JSON envelopes should be stable, simple, and test-covered. |
| Windows-specific path behavior | Cross-platform hardening | First release prioritizes Linux/macOS and this environment is Linux. | Path helpers should use `filepath` and avoid Unix-only assumptions where practical. |

## Final Decisions

### Decision 1: Package Ownership And Dependency Direction

- **Decision:** Implement Sprint 2 logic under `internal/` with `cmd/ultraplan/main.go` remaining a thin entrypoint. `internal/app` owns command registration, dependency construction, error-to-exit mapping, and output rendering. `internal/workspace` owns discovery, canonical paths, init planning/creation, and structural validation. `internal/platform/config` owns effective config, defaults, workspace config parsing, env/CLI overrides, validation, and redaction. `internal/platform/logging` owns minimal safe logging. Platform packages must not import product modules; `workspace` may import only platform packages where genuinely needed.
- **Rationale:** This matches ARCHITECTURE.md's module-driven rule and keeps behavior near the state it transforms. It also leaves later `study` and `codeextract` modules free to consume workspace/config services without creating global technical layers.
- **Study / Source Grounding:** `technical-handbook.md` relevant patterns and warnings; `01-project-structure` cites thin entrypoints and `internal/` boundaries from restic and chezmoi; `02-command-architecture` and `03-dependency-injection` cite command factories and manual composition roots.
- **Trade-Offs Accepted:** Accepts `internal/` encapsulation over public package reuse and a small app composition layer over per-command ad hoc wiring.
- **Technical Debt / Future Impact:** The app dependency container must stay small; split only when concrete complexity appears. Missing area architecture reasoning is superseded by this final decision.
- **Alternatives Rejected:** Put all logic in `cmd/ultraplan` was rejected because it violates C-05 and study evidence warns against logic-heavy entrypoints. Create global packages like `internal/validation`, `internal/output`, or `internal/paths` was rejected because ARCHITECTURE.md warns that global technical layers fracture product context.
- **Contracts Satisfied:** Architecture, Testing, CLI Surface, AC-01, AC-02, C-01, C-05, C-06, C-07.
- **Evidence Required:** Import-direction review; `go test ./...`; `go build ./cmd/ultraplan`; command help tests showing `cmd/ultraplan` delegates to app command registration.

### Decision 2: Workspace Discovery, Markers, And Initialization

- **Decision:** Workspace discovery must resolve root in this order: explicit `--workspace`, `ULTRAPLAN_WORKSPACE`, current directory, then nearest parent. A valid discovered workspace is identified by `ultraplan.yml` plus enough required structure for structural validation. `init-workspace` creates `ultraplan.yml`, `prompts/base.md`, `prompts/synthesize.md`, `templates/repo-analysis.md`, `templates/report.md`, and `studies/`, with dry-run returning the planned operations without writing. Initialization must be idempotent for existing expected files and must write only inside the selected target directory.
- **Rationale:** TRD 5.1 mandates discovery precedence and TRD 5.2 plus sprint AC-03 mandate the required structure. Treating `ultraplan.yml` as the primary marker resolves the TRD open question for this sprint while structural validation prevents false positives.
- **Study / Source Grounding:** `technical-handbook.md` design pressures for workspace discovery/init and path safety; `13-security` path canonicalization and trust-boundary evidence; `06-io-abstraction` for temp-directory and filesystem test seams. Concrete study sources are relevant because this decision involves path safety and testable filesystem behavior.
- **Trade-Offs Accepted:** Using `ultraplan.yml` as the primary marker is simpler than accepting `.ultraplan/` alone, but requires users to initialize or provide config before discovery succeeds. Dry-run planning adds a small implementation layer but makes mutating behavior reviewable.
- **Technical Debt / Future Impact:** If future migrations introduce alternate markers, discovery may need marker-version handling. Current behavior must not assume optional `.ultraplan/` directories exist.
- **Alternatives Rejected:** Using `.ultraplan/` as the only marker was rejected because the required workspace config is `ultraplan.yml` and optional `.ultraplan/` directories are not required by TRD 5.2. Auto-creating a workspace when discovery fails was rejected because it would make commands unexpectedly mutating and script-hostile.
- **Contracts Satisfied:** Architecture, Security, CLI Surface, Testing, AC-01, AC-03, AC-05, C-02, C-07, C-08.
- **Evidence Required:** `internal/workspace/discovery_test.go` for precedence and parent traversal; `init_test.go` for planning, creation, dry-run, idempotence, and invalid path handling; command tests for `init-workspace --dry-run` and real init in temp dirs.

### Decision 3: Workspace Path Safety And Structural Validation

- **Decision:** `internal/workspace` must expose focused path helpers that normalize workspace-managed paths, return canonical workspace-relative paths for display where possible, and reject paths that escape the workspace for workspace-scoped operations. Structural validation checks required files/directories only: `ultraplan.yml`, `prompts/base.md`, `prompts/synthesize.md`, `templates/repo-analysis.md`, `templates/report.md`, and `studies/`.
- **Rationale:** Sprint scope requires workspace inspection and validation but excludes study validation and runtime behavior. Path escape rejection is required by TRD 5.3 and the Security contract.
- **Study / Source Grounding:** `technical-handbook.md` design pressures and anti-patterns; `13-security` cites safe path/command handling and trust-boundary examples; `06-io-abstraction` supports temp-dir tests for filesystem behavior.
- **Trade-Offs Accepted:** Validation remains structural instead of semantic. This keeps health fast and local but does not catch invalid study definitions or runtime readiness.
- **Technical Debt / Future Impact:** Later study validation must live in `internal/study`, not be added to workspace validation. Symlink handling may need hardening if future commands operate on untrusted source trees.
- **Alternatives Rejected:** Performing deep recursive workspace scans was rejected because non-goals exclude study/source analysis and PRD performance guidance warns against unnecessary scans. Allowing absolute paths freely was rejected because TRD 5.3 allows absolute paths only when explicitly allowed and Security requires path validation.
- **Contracts Satisfied:** Architecture, Security, Testing, AC-03, AC-05, C-02, C-07.
- **Evidence Required:** `paths_test.go` for normalization and escape rejection; `validation_test.go` for missing required files/directories; health command tests for invalid workspace status.

### Decision 4: Effective Config, Precedence, Validation, And Redaction

- **Decision:** `internal/platform/config` must build an effective config by applying built-in defaults, workspace `ultraplan.yml`, environment variables with an `ULTRAPLAN_` prefix, then CLI overrides only when flags are explicitly set. Validation runs after the merge. `config show` displays a redacted projection, not raw config. Config source metadata should be recorded where practical for diagnostics, especially in JSON output, but source metadata must not delay the core precedence contract.
- **Rationale:** The TRD and sprint requirements define exact precedence and require redaction. Explicit changed-flag handling prevents flag defaults from shadowing workspace config.
- **Study / Source Grounding:** `technical-handbook.md` configuration patterns and warnings; `04-configuration-management` cites layered precedence and flag-change tracking; `13-security` cites redaction wrappers/projections and schema validation; `05-error-handling` supports validation failures with field paths and guidance.
- **Trade-Offs Accepted:** Fixed precedence is less flexible than a priority abstraction but easier to test and explain. Projection redaction is accepted for Sprint 2, with centralized sensitive-field detection to reduce accidental leaks.
- **Technical Debt / Future Impact:** Runtime-facing config mapping to agentwrap is deferred; future sprints may add secret wrapper types and richer source metadata. Unknown field policy should be conservative: fail unknown core fields that could indicate typos, while allowing explicitly documented extension maps only if introduced later.
- **Alternatives Rejected:** Reading environment variables directly in commands was rejected because it makes precedence unpredictable. Letting workspace config override CLI flags was rejected because it violates TRD and sprint AC-06. Printing raw config structs was rejected because it risks leaking secrets.
- **Contracts Satisfied:** Configuration, Security, Errors, Testing, AC-01, AC-04, AC-06, AC-07, C-03, C-04, C-07.
- **Evidence Required:** `config_test.go` for defaults, file loading, env overrides, explicit CLI overrides, validation, and source precedence; `redaction_test.go` proving sensitive values are redacted; `config_commands_test.go` for text/JSON `config show` outputs and no secret leakage.

### Decision 5: Minimal Secret-Safe Logging Foundation

- **Decision:** Implement `internal/platform/logging` as a minimal text/JSON logger with level support, structured fields, stderr output, and redaction before emission. It must be context-friendly enough for later expansion but must not implement durable file logs, event records, agentwrap observability sinks, run stores, or full runtime diagnostics in this sprint.
- **Rationale:** Sprint requirements include a logging foundation but explicitly exclude durable logs and full observability persistence. Keeping diagnostics on stderr protects stdout as command data for scripts.
- **Study / Source Grounding:** `technical-handbook.md` logging pattern and warning; `10-logging-observability` cites structured, level-aware logging and stdout/stderr separation; `13-security` supports credential scrubbing and safe diagnostics.
- **Trade-Offs Accepted:** Minimal interface may need adaptation when runtime events arrive. Deferring durable logs avoids scope creep into runtime/run-loop behavior.
- **Technical Debt / Future Impact:** Field naming should align with TRD 19 where feasible (`command`, `workspace`, `level`, `message`) so later runtime fields can be added without replacing the logger.
- **Alternatives Rejected:** Writing logs to workspace files now was rejected because durable file logs and event persistence are sprint non-goals. Using package-level global mutable logger state was rejected because tests must be isolated and DI evidence warns against globals.
- **Contracts Satisfied:** Observability, Security, Testing, AC-01, AC-07, C-04, C-07, C-08.
- **Evidence Required:** Unit tests for log formatting and redaction; command tests or review checks proving diagnostics are on stderr and command payloads are on stdout.

### Decision 6: Deterministic Text And JSON Output Helpers

- **Decision:** Implement deterministic output helpers in `internal/app` for command payloads. `config show` and `health` must support text and `--json`. JSON output must be stable, contain no ANSI formatting, and include command name, status, workspace when known, and result payload. Text output must be concise, deterministic, and script-friendly. Diagnostics and logs must not be mixed into stdout payloads.
- **Rationale:** CLI Surface and TRD output requirements require stable text and JSON output. Keeping output rendering in app avoids platform packages importing product payloads.
- **Study / Source Grounding:** `technical-handbook.md` output and IO patterns; `06-io-abstraction` for injected streams; `09-terminal-ux` for scriptability and non-interactive behavior; `10-logging-observability` for stdout/stderr separation.
- **Trade-Offs Accepted:** No standalone `platform/output` abstraction now; output helpers may be moved later if multiple modules need direct rendering support.
- **Technical Debt / Future Impact:** JSON schemas should remain minimal to avoid premature public commitments. Later docs may formalize them.
- **Alternatives Rejected:** Human-only text output was rejected because AC-04, AC-05, and AC-08 require JSON modes. Colorized or TTY-dependent output was rejected because commands must be deterministic and script-friendly.
- **Contracts Satisfied:** CLI Surface, Documentation, Observability, Testing, AC-04, AC-05, AC-08, AC-09, C-08.
- **Evidence Required:** Command tests asserting exact or semantically exact text output; JSON decode tests for `config show --json` and `health --json`; help output tests.

### Decision 7: Error Classification And Exit Codes

- **Decision:** Implement classified errors or an equivalent error classification function for usage/argument errors, config errors, workspace/filesystem errors, validation failures, cancellation classes reserved for future commands, runtime classes reserved for future commands, and general failures. Map to TRD suggested exit codes: `0` success, `1` general failure, `2` usage, `3` config, `4` workspace/filesystem, `5` validation, `6` runtime, `7` cancellation, `8` partial completion. Sprint 2 commands should only emit classes that are in scope, but the mapping may reserve future runtime/cancellation/partial codes.
- **Rationale:** Deterministic exit codes are required by AC-10 and TRD 7.2. Reserving the full map now prevents incompatible command behavior later without implementing runtime behavior.
- **Study / Source Grounding:** `technical-handbook.md` error pattern; `05-error-handling` cites wrapped errors, typed/sentinel classification, and exit-code constants from mature CLIs.
- **Trade-Offs Accepted:** A simple classifier is enough now; a larger error package may be needed when runtime and run-loop failures are introduced.
- **Technical Debt / Future Impact:** Error messages must include actionable context without leaking secrets. Runtime codes are reserved but not exercised by Sprint 2 health because OpenCode/provider checks are out of scope.
- **Alternatives Rejected:** Returning `1` for all failures was rejected because it violates AC-10 and weakens automation. Parsing error strings to choose exit codes was rejected because TRD error requirements and study evidence favor typed/classified errors.
- **Contracts Satisfied:** Errors, CLI Surface, Security, Testing, AC-10, C-04, C-08.
- **Evidence Required:** Command tests for bad flags/args exit `2`, invalid config exit `3`, missing/invalid workspace or filesystem path exit `4`, structural validation failure exit `5` where the command treats validation as terminal, and success exit `0`.

### Decision 8: Command Scope And Behavior

- **Decision:** Register and implement only the Sprint 2 command surface: global `--workspace`, `init-workspace [--path <dir>] [--dry-run]`, `config show [--json]`, and `health [--json]`, plus help output for these commands. `health` must report workspace discovery/validation, config parse/validation, basic filesystem accessibility, and relevant environment override presence/shape. It must not check OpenCode, provider credentials, model availability, agentwrap runtime health, study state, source repositories, or target/sprint files.
- **Rationale:** This delivers the sprint goal while enforcing explicit non-goals. Health should be useful before runtime integration but must not pretend to validate runtime readiness.
- **Study / Source Grounding:** `technical-handbook.md` design pressures for health skeleton and command factories; `02-command-architecture` command construction evidence; `09-terminal-ux` scriptability evidence. No runtime study sources are relevant because runtime checks are explicitly out of scope.
- **Trade-Offs Accepted:** Health is intentionally shallow and local. Users will not get provider/runtime readiness yet, but tests remain offline and deterministic.
- **Technical Debt / Future Impact:** Health check names and JSON shape should allow later runtime checks to be appended without breaking existing consumers. Runtime unavailable must not be reported as failure in Sprint 2 unless runtime config is syntactically invalid.
- **Alternatives Rejected:** Adding OpenCode/provider health checks now was rejected because requirements explicitly exclude agentwrap/OpenCode/provider dependencies. Adding study validation to health was rejected because workspace validation beyond structure is a non-goal.
- **Contracts Satisfied:** CLI Surface, Documentation, Configuration, Security, Testing, AC-01, AC-03, AC-04, AC-05, AC-08, AC-09, C-07, C-08.
- **Evidence Required:** `workspace_commands_test.go`, `config_commands_test.go`, and `health_commands_test.go`; manual/review checks for help output; `health --json` tests showing no dependency on OpenCode, credentials, or network.

### Decision 9: Testing And Verification Strategy

- **Decision:** Use table-driven unit tests for workspace, paths, validation, init planning, config, redaction, logging, and error mapping; command-level tests with injected IO and temp directories for `init-workspace`, `config show`, and `health`; JSON decode assertions for structured output; exact or stable-substring assertions for deterministic text output. Verification commands are `go test ./...` and `go build ./cmd/ultraplan` from `../ultraplan-go`.
- **Rationale:** Sprint 2 behavior is mostly local filesystem/config/CLI behavior, so deterministic temp-dir and command tests provide stronger evidence than mocks of external systems. Normal tests must not require OpenCode, provider credentials, network, or a preconfigured UltraPlan workspace.
- **Study / Source Grounding:** `technical-handbook.md` testing pattern; `11-testing-strategy` cites table-driven, command-level, and golden/output checks; `06-io-abstraction` supports injected IO; `03-dependency-injection` supports constructor seams.
- **Trade-Offs Accepted:** Not every text output needs full golden files; exact command assertions are sufficient for small outputs.
- **Technical Debt / Future Impact:** As command outputs grow, golden fixtures may become more maintainable than inline expected strings.
- **Alternatives Rejected:** Manual-only verification was rejected because AC-01 requires test pass and the Testing contract requires deterministic coverage. Real OpenCode/provider integration tests were rejected for this sprint because AC-01 requires normal tests to run offline.
- **Contracts Satisfied:** Testing, CLI Surface, Documentation, AC-01, AC-02, AC-03, AC-04, AC-05, AC-07, AC-08, AC-09, AC-10, C-07.
- **Evidence Required:** Passing `go test ./...`; passing `go build ./cmd/ultraplan`; test files listed in sprint requirements with behavior coverage for all acceptance criteria.

## Expected Evidence

| Evidence Type | Required Evidence | Source / Command / Review Check |
| --- | --- | --- |
| Tests | Full unit and command test suite passes offline without configured workspace, network, OpenCode, or credentials. | `go test ./...` from `../ultraplan-go`. |
| Build | CLI builds as a runnable binary. | `go build ./cmd/ultraplan` from `../ultraplan-go`. |
| Workspace | Discovery precedence, parent traversal, init dry-run, init creation, idempotence, path normalization, escape rejection, and structural validation are covered. | `internal/workspace/*_test.go`, `internal/app/workspace_commands_test.go`. |
| Config | Defaults, workspace config, env overrides, explicit CLI overrides, validation, source precedence, text/JSON display, and redaction are covered. | `internal/platform/config/*_test.go`, `internal/app/config_commands_test.go`. |
| Health | Text/JSON health reports basic workspace/config/filesystem/environment checks and does not require runtime/provider setup. | `internal/app/health_commands_test.go`. |
| Output | Text and JSON outputs are deterministic; JSON decodes and contains no ANSI formatting; stdout data is separated from stderr diagnostics. | Command tests and output helper tests. |
| Logging | Logs support text/JSON format, level fields, stderr emission, and redaction. | `internal/platform/logging/logging_test.go` or command-level stderr assertions. |
| Errors | Usage, config, workspace/filesystem, validation, and general failures map to deterministic exit codes. | Error mapping tests and command tests. |
| Help | CLI help includes `init-workspace`, `config`, `config show`, `health`, relevant flags, and global `--workspace`. | Help-output tests or sprint review manual checks. |
| Review | Package ownership and import direction respect selected architecture. | Architecture Review protocol from `sprint-index.md`. |
| Scope | No study, target, sprint, runtime execution, OpenCode, provider, durable logs, or event persistence behavior was implemented. | Sprint Review protocol and code review. |

## Assumptions And Risks

| Item | Type | Impact | Mitigation / Follow-Up |
| --- | --- | --- | --- |
| Sprint 1 CLI foundation exists and can be extended. | Assumption | If the command router differs from expected, command wiring may need adaptation. | Inspect existing `../ultraplan-go` layout during implementation and preserve Sprint 1 patterns while enforcing these decisions. |
| `ultraplan.yml` is the primary workspace marker for Sprint 2. | Assumption | Workspaces with only `.ultraplan/` will not be discovered. | `init-workspace` creates `ultraplan.yml`; future migration work can add marker compatibility if needed. |
| Config schema can include TRD starter fields without runtime validation. | Assumption | Users may see runtime/model fields before runtime execution exists. | Health and config validation must distinguish syntactic config validity from runtime/provider availability. |
| Health output may be consumed by scripts early. | Risk | Later runtime checks could break consumers if the schema changes abruptly. | Use a stable top-level JSON envelope and append checks rather than renaming existing fields. |
| Projection-only redaction can be bypassed by future direct formatting of raw config. | Risk | Secret leakage in output or logs. | Centralize output through redacted projections and add explicit tests with known secret values. |
| Path safety around symlinks can be subtle. | Risk | Workspace escape checks may miss platform-specific cases. | Use `filepath.Clean`, absolute/canonical comparisons where feasible, temp-dir tests, and conservative rejection when resolution is ambiguous. |
| Missing area-specific architecture reasoning artifact. | Risk | Governance may expect a separate architecture artifact. | This document records final architecture decisions; implementation and plan must follow this document unless governance explicitly requires backfill. |
| Exit-code classes reserved for runtime/cancellation are not exercised now. | Risk | Future code may interpret reserved mappings differently. | Document the mapping in tests or app error code constants now. |

## Implementation Constraints

- Use Go and the standard Go toolchain.
- Keep `cmd/ultraplan/main.go` thin; command setup and behavior must live under `internal/app` and owned internal packages.
- `internal/workspace` owns discovery, init, paths, and structural workspace validation only.
- `internal/platform/config` and `internal/platform/logging` must not import product modules.
- Do not add study, source, analysis, synthesis, code extraction, target, sprint, runtime execution, OpenCode, provider, or agentwrap health behavior in this sprint.
- Workspace discovery precedence must be explicit `--workspace`, `ULTRAPLAN_WORKSPACE`, current directory, nearest parent directory.
- Config precedence must be defaults, workspace config, environment variables, CLI flags, with CLI flags applied only when explicitly set.
- Sensitive values must be redacted before output, diagnostics, or logs.
- Command data goes to stdout; diagnostics/logs go to stderr.
- Text and JSON output must be deterministic and script-friendly.
- JSON output must not contain ANSI formatting.
- `init-workspace --dry-run` must not write files.
- Workspace-managed paths must be normalized and must reject escapes from the selected workspace.
- Tests must use temp directories and injected IO where needed; normal tests must be offline and must not require OpenCode, credentials, network, or a configured workspace.
- Exit codes must map deterministically to the selected error classes.

## Plan Handoff

`plan.md` must execute these decisions. It must not invent architecture, scope, package boundaries, config precedence, health scope, output contracts, error classes, or runtime behavior beyond this document.

The plan must carry forward:

- Final decisions in this document.
- Requirement IDs AC-01 through AC-10 and constraints C-01 through C-08.
- Selected contracts from `sprint-index.md`: Architecture, Errors, Configuration, Observability, Security, Testing, Documentation, and CLI Surface.
- Expected evidence from this document.
- Risks and mitigations from this document.
- Required review protocols from `sprint-index.md`: Architecture Review and Sprint Review.

## Phase Exit Criteria

- [x] Selected context was read and used.
- [x] Area-specific reasoning documents were completed where present or explicitly recorded as absent.
- [x] Area-specific reasoning conclusions are reflected or explicitly unavailable.
- [x] Contracts are explicitly mapped to decisions and expected evidence.
- [x] Final decisions are clear enough for `plan.md` to execute.
- [x] Expected evidence is specific and reviewable.
