# Sprint Plan: Workspace, Config, Logging, and Health Skeleton

> Project: `ultraplan-go`
> Sprint: `02-workspace-config-health`
> Source: `projects/ultraplan-go/sprints/02-workspace-config-health/reasoning.md`
> **Inputs Used:** `projects/ultraplan-go/sprints/02-workspace-config-health/requirements.md`, `projects/ultraplan-go/sprints/02-workspace-config-health/reasoning.md`, `projects/ultraplan-go/sprints/02-workspace-config-health/sprint-index.md`, `projects/ultraplan-go/sprints/02-workspace-config-health/technical-handbook.md`, `projects/ultraplan-go/docs/PRD.md`, `projects/ultraplan-go/docs/TRD.md`, `projects/ultraplan-go/docs/ARCHITECTURE.md`, `projects/ultraplan-go/project-index.md`, `templates/sprint-plan.md`; area reasoning omitted because `projects/ultraplan-go/sprints/02-workspace-config-health/reasoning/` was not present.

This plan executes `reasoning.md`. It does not introduce architecture, scope, command behavior, config precedence, output contracts, or runtime checks beyond the consolidated sprint reasoning.

## Reasoning Source

- **Sprint Reasoning:** `projects/ultraplan-go/sprints/02-workspace-config-health/reasoning.md`
- **Sprint Index:** `projects/ultraplan-go/sprints/02-workspace-config-health/sprint-index.md`
- **Technical Handbook:** `projects/ultraplan-go/sprints/02-workspace-config-health/technical-handbook.md`
- **Area Reasoning:** none present; `reasoning.md` records the absent `reasoning/architecture.md` artifact and makes final package-boundary decisions directly from project docs, selected contracts, and handbook evidence.
- **Plan Template:** `templates/sprint-plan.md`

## Sprint Status

- **Status:** complete
- **Owner:** implementation agent
- **Start Date:** 2026-05-30
- **Completion Date:** 2026-05-30

## Decisions To Execute

| Decision | Source Section | Execution Implication |
| --- | --- | --- |
| Package ownership and dependency direction | `reasoning.md#decision-1-package-ownership-and-dependency-direction` | Keep `cmd/ultraplan/main.go` thin; put command registration/composition/output/error mapping in `internal/app`; put workspace behavior in `internal/workspace`; put config/logging infrastructure in `internal/platform/config` and `internal/platform/logging`; prevent platform packages from importing product modules. |
| Workspace discovery, markers, and initialization | `reasoning.md#decision-2-workspace-discovery-markers-and-initialization` | Implement discovery precedence `--workspace`, `ULTRAPLAN_WORKSPACE`, current directory, nearest parent; use `ultraplan.yml` as the primary marker for Sprint 2; implement idempotent `init-workspace` with dry-run planning and required scaffold creation only inside the selected path. |
| Workspace path safety and structural validation | `reasoning.md#decision-3-workspace-path-safety-and-structural-validation` | Normalize workspace-managed paths, reject workspace escapes, and validate only required Sprint 2 structural files and directories. |
| Effective config, precedence, validation, and redaction | `reasoning.md#decision-4-effective-config-precedence-validation-and-redaction` | Build effective config in fixed order: defaults, workspace config, environment, explicit CLI flags; validate after merge; expose only redacted projections through `config show`, health diagnostics, and logs. |
| Minimal secret-safe logging foundation | `reasoning.md#decision-5-minimal-secret-safe-logging-foundation` | Implement text/JSON, level-aware, structured, stderr-based logging with redaction; do not implement durable file logs, event records, or runtime observability persistence. |
| Deterministic text and JSON output helpers | `reasoning.md#decision-6-deterministic-text-and-json-output-helpers` | Render command payloads through deterministic helpers in `internal/app`; support text and `--json` for `config show` and `health`; keep diagnostics/logs off stdout. |
| Error classification and exit codes | `reasoning.md#decision-7-error-classification-and-exit-codes` | Implement classified errors or equivalent mapping for usage, config, workspace/filesystem, validation, general, and reserved future runtime/cancellation/partial classes using TRD exit codes. |
| Command scope and behavior | `reasoning.md#decision-8-command-scope-and-behavior` | Register only global `--workspace`, `init-workspace`, `config show`, and `health`; keep health limited to workspace, config, filesystem, and environment basics; exclude OpenCode/provider/runtime/study/target/sprint behavior. |
| Testing and verification strategy | `reasoning.md#decision-9-testing-and-verification-strategy` | Use table-driven unit tests, command tests with injected IO and temp directories, JSON decode assertions, output assertions, and final `go test ./...` plus `go build ./cmd/ultraplan`. |

## Requirements / Contracts To Satisfy

| Contract / Requirement ID | Required Behavior | Evidence Planned |
| --- | --- | --- |
| AC-01, C-07, Testing | `go test ./...` passes offline without OpenCode, credentials, network, or configured workspace. | Full test suite from `../ultraplan-go`. |
| AC-02, C-01 | CLI builds with standard Go toolchain. | `go build ./cmd/ultraplan` from `../ultraplan-go`. |
| AC-03, C-02, Security | `init-workspace [--path <dir>] [--dry-run]` plans and creates only required workspace structure inside selected path. | `internal/workspace/init_test.go`, `internal/app/workspace_commands_test.go`, manual dry-run/create review. |
| AC-04, AC-06, AC-07, C-03, C-04, Configuration | `config show [--json]` prints effective config with correct precedence and redacted sensitive values. | `internal/platform/config/config_test.go`, `redaction_test.go`, `internal/app/config_commands_test.go`. |
| AC-05 | `health [--json]` reports basic workspace, config, filesystem, and environment status without runtime/provider dependency. | `internal/app/health_commands_test.go`, offline `go test ./...`. |
| AC-08, C-08, CLI Surface | Text and JSON outputs are deterministic, script-friendly, and ANSI-free in JSON. | Output helper tests and command tests that decode JSON and assert text. |
| AC-09, Documentation, CLI Surface | Help output includes new commands and relevant flags. | Help-output command tests or sprint review evidence. |
| AC-10, Errors | Error classes map to deterministic exit codes. | Error mapping tests and command tests for usage, config, workspace/filesystem, validation, and general failures. |
| Architecture, C-05, C-06 | Entrypoint remains thin; platform packages do not import product modules. | Architecture review, import inspection, build/test evidence. |
| Observability, Security | Logging and diagnostics redact secrets and keep data stdout separate from diagnostics stderr. | Logging tests and command stderr/stdout assertions. |

## Tasks

- [x] **Task 1: Inspect Sprint 1 Baseline And Preserve Existing Shape**
  > Executes: `Decision 1`, `Decision 9`, `AC-01`, `AC-02`, `C-05`, `C-06`
  - [x] Inspect `../ultraplan-go` module layout, existing `cmd/ultraplan`, `internal/app`, command router, test style, and any existing error/output helpers.
  - [x] Record any mismatch with `reasoning.md` as an implementation note before editing, without changing sprint scope.
  - [x] Identify the smallest integration points for app dependencies, injected IO, command construction, and exit-code mapping.
  - [x] Confirm no existing unrelated worktree changes are reverted or overwritten.

- [x] **Task 2: Establish App Composition, IO, Output, And Error Skeleton**
  > Executes: `Decision 1`, `Decision 6`, `Decision 7`, `Decision 9`, `AC-08`, `AC-10`, `C-05`, `C-08`
  - [x] Keep `cmd/ultraplan/main.go` delegating to `internal/app` without workspace/config/health business logic.
  - [x] Add or extend app-level command construction with a small dependency container for stdout, stderr, env lookup, filesystem access where needed, logger, config loader, and workspace services.
  - [x] Implement deterministic app output helpers for text and JSON command payloads, including stable field order through struct-based JSON output.
  - [x] Implement classified app errors or an equivalent classifier for exit-code mapping: `0`, `1`, `2`, `3`, `4`, `5`, and reserved `6`, `7`, `8`.
  - [x] Add tests for output rendering, stdout/stderr separation where feasible, and error-to-exit-code mapping.

- [x] **Task 3: Implement Workspace Discovery And Path Safety**
  > Executes: `Decision 2`, `Decision 3`, `AC-03`, `AC-05`, `C-02`, `C-07`
  - [x] Create `internal/workspace/discovery.go` to resolve explicit path, `ULTRAPLAN_WORKSPACE`, current directory, and nearest parent in that order.
  - [x] Treat `ultraplan.yml` as the primary marker for this sprint and avoid accepting `.ultraplan/` alone as a complete workspace.
  - [x] Create `internal/workspace/paths.go` with workspace root normalization, workspace-relative display helpers where useful, and workspace escape rejection for workspace-scoped paths.
  - [x] Use `filepath` APIs and temp-directory tests; avoid Unix-only assumptions unless unavoidable.
  - [x] Add `discovery_test.go` and `paths_test.go` covering precedence, parent traversal, invalid candidates, normalization, and escape rejection.

- [x] **Task 4: Implement Workspace Initialization And Structural Validation**
  > Executes: `Decision 2`, `Decision 3`, `AC-03`, `AC-05`, `C-07`, `C-08`
  - [x] Create `internal/workspace/init.go` with a dry-run plan describing files/directories that would be created and a creation path that is idempotent for expected existing files.
  - [x] Create starter `ultraplan.yml`, `prompts/base.md`, `prompts/synthesize.md`, `templates/repo-analysis.md`, `templates/report.md`, and `studies/` only under the selected target directory.
  - [x] Create `internal/workspace/validation.go` that checks only required Sprint 2 structural files/directories.
  - [x] Add unit tests for scaffold planning, no writes in dry-run, real creation, idempotence, invalid path handling, missing required files, and missing required directories.

- [x] **Task 5: Wire `init-workspace` Command**
  > Executes: `Decision 2`, `Decision 6`, `Decision 7`, `Decision 8`, `AC-03`, `AC-09`, `AC-10`
  - [x] Register `init-workspace [--path <dir>] [--dry-run]` in `internal/app` with help text.
  - [x] Use workspace init planning for dry-run and creation for normal execution.
  - [x] Render deterministic text output for planned and created operations.
  - [x] Return usage exit code for invalid flags/arguments and workspace/filesystem exit code for path or write failures.
  - [x] Add command tests for help, dry-run, creation, invalid path handling, idempotence, output, and exit codes.

- [x] **Task 6: Implement Effective Config Loading And Validation**
  > Executes: `Decision 4`, `AC-04`, `AC-06`, `C-03`, `C-07`
  - [x] Create `internal/platform/config/config.go` with explicit built-in defaults aligned to TRD starter fields for schema version, runtime, models, execution, logging, and agentwrap-related config fields that are syntactic only in this sprint.
  - [x] Load workspace `ultraplan.yml` after workspace discovery when a workspace is available.
  - [x] Apply environment overrides from centralized `ULTRAPLAN_` variables only through the config package.
  - [x] Apply CLI overrides only when flags are explicitly set so flag defaults do not shadow workspace config.
  - [x] Validate schema version, required fields, duration syntax, positive numeric values, known logging formats, and syntactic runtime names without checking runtime availability.
  - [x] Record source metadata where practical for diagnostics, especially JSON output, but do not block precedence behavior on rich metadata.
  - [x] Add unit tests for defaults, workspace file loading, environment overrides, explicit CLI overrides, precedence, validation failures with field paths, and no provider/OpenCode dependency.

- [x] **Task 7: Implement Config Redaction And `config show` Command**
  > Executes: `Decision 4`, `Decision 6`, `Decision 7`, `Decision 8`, `AC-04`, `AC-07`, `AC-08`, `AC-09`, `C-04`
  - [x] Create `internal/platform/config/redaction.go` with centralized sensitive-field detection and a redacted config projection.
  - [x] Ensure sensitive values cannot appear in text output, JSON output, validation diagnostics, or logging fields used by this sprint.
  - [x] Register `config show [--json]` and relevant help text.
  - [x] Render deterministic text and JSON effective-config output from the redacted projection.
  - [x] Return config exit code for invalid config and workspace/filesystem exit code for missing/unreadable workspace config when the command requires it.
  - [x] Add redaction tests with known secret values and command tests for text, JSON, precedence, redaction, and exit codes.

- [x] **Task 8: Implement Minimal Secret-Safe Logging**
  > Executes: `Decision 5`, `AC-07`, `C-04`, `C-08`, `Observability`, `Security`
  - [x] Create `internal/platform/logging/logging.go` with text and JSON formats, levels `debug`, `info`, `warn`, `error`, structured fields, and stderr emission.
  - [x] Apply redaction before emission for fields whose keys or values are classified as sensitive by this sprint's redaction policy.
  - [x] Keep the logger instance injectable; avoid package-level mutable global loggers.
  - [x] Do not create durable file logs, event records, run stores, agentwrap observability sinks, or runtime diagnostics in this sprint.
  - [x] Add tests for text/JSON log formatting, level fields, stderr-oriented usage where feasible, and secret absence.

- [x] **Task 9: Implement `health` Command Skeleton**
  > Executes: `Decision 3`, `Decision 4`, `Decision 6`, `Decision 7`, `Decision 8`, `AC-05`, `AC-08`, `AC-09`, `AC-10`
  - [x] Register `health [--json]` and help text.
  - [x] Report workspace discovery result, structural validation result, config parse/validation result, filesystem readability/writability checks for required paths, and relevant environment override presence/shape.
  - [x] Mark OpenCode, provider credentials, model availability, agentwrap runtime health, study state, source repositories, target files, and sprint files as out of scope rather than checking them.
  - [x] Render deterministic text and JSON health output with a stable top-level envelope and append-only check list shape.
  - [x] Return deterministic exit codes for invalid config, missing/invalid workspace, validation failure, and success according to the app error classifier.
  - [x] Add command tests for valid workspace, missing workspace, invalid workspace structure, invalid config, text output, JSON output, no credential/OpenCode dependency, and exit codes.

- [x] **Task 10: Complete Help, Import, Scope, And Verification Review**
  > Executes: `Decision 1`, `Decision 8`, `Decision 9`, `AC-01`, `AC-02`, `AC-09`, `C-05`, `C-06`
  - [x] Add or update help-output tests for `init-workspace`, `config`, `config show`, `health`, `--workspace`, `--json`, and `--dry-run` where applicable.
  - [x] Inspect imports to confirm platform packages do not import product modules and `cmd/ultraplan` remains thin.
  - [x] Search implementation for excluded behavior: study listing/init/runs, source discovery, runtime execution, agentwrap/OpenCode health checks, provider checks, durable logs, event persistence, target/sprint workflows.
  - [x] Run `go test ./...` from `../ultraplan-go`.
  - [x] Run `go build ./cmd/ultraplan` from `../ultraplan-go`.
  - [x] Record verification evidence and any deviations or deferrals before sprint review.

## Evidence Checklist

- [x] Tests prove workspace discovery precedence, parent traversal, init planning, init creation, idempotence, path normalization, escape rejection, and structural validation.
- [x] Tests prove effective config defaults, workspace config loading, environment overrides, explicit CLI overrides, validation, source precedence where implemented, and redaction.
- [x] Tests prove `config show` and `health` text/JSON output is deterministic and secret-safe.
- [x] Tests prove logging text/JSON formatting, levels, stderr-oriented diagnostics, and redaction.
- [x] Tests prove error classes map to deterministic exit codes.
- [x] Tests prove help output includes Sprint 2 commands and flags.
- [x] `go test ./...` passes offline from `../ultraplan-go`.
- [x] `go build ./cmd/ultraplan` succeeds from `../ultraplan-go`.
- [x] Architecture review evidence shows package ownership and dependency direction match `reasoning.md`.
- [x] Sprint review evidence shows explicit non-goals were not implemented.
- [x] Deviations from `reasoning.md` are recorded before implementation continues.

## Verification Commands

| Check | Command | Expected Result |
| --- | --- | --- |
| Unit and command tests | `go test ./...` from `../ultraplan-go` | Passes without OpenCode, provider credentials, network access, or configured UltraPlan workspace. |
| CLI build | `go build ./cmd/ultraplan` from `../ultraplan-go` | Builds a runnable `ultraplan` binary. |
| Workspace dry-run review | `./ultraplan init-workspace --path <temp-dir> --dry-run` from built binary location | Prints deterministic planned operations and writes no files. |
| Workspace init review | `./ultraplan init-workspace --path <temp-dir>` from built binary location | Creates only `ultraplan.yml`, required prompt/template files, and `studies/` under the selected directory. |
| Config text review | `./ultraplan --workspace <temp-dir> config show` | Prints effective config with sensitive values redacted. |
| Config JSON review | `./ultraplan --workspace <temp-dir> config show --json` | Emits valid deterministic JSON with no secrets and no ANSI formatting. |
| Health text review | `./ultraplan --workspace <temp-dir> health` | Reports workspace/config/filesystem/environment basics without runtime/provider checks. |
| Health JSON review | `./ultraplan --workspace <temp-dir> health --json` | Emits valid deterministic JSON with basic checks and no secrets. |
| Help review | `./ultraplan --help`, `./ultraplan init-workspace --help`, `./ultraplan config show --help`, `./ultraplan health --help` | Help includes Sprint 2 commands and relevant flags. |

## Risks And Blockers

| Risk / Blocker | Source | Mitigation | Status |
| --- | --- | --- | --- |
| Sprint 1 CLI foundation may differ from expected command router or package layout. | `reasoning.md#assumptions-and-risks` | Inspect baseline first and preserve established Sprint 1 patterns while enforcing `reasoning.md` dependency rules. | closed |
| `ultraplan.yml` as primary marker excludes directories with only `.ultraplan/`. | `reasoning.md#decision-2-workspace-discovery-markers-and-initialization` | `init-workspace` creates `ultraplan.yml`; future marker compatibility can be added only with a concrete requirement. | closed |
| Projection-only redaction can be bypassed by direct formatting of raw config. | `reasoning.md#assumptions-and-risks` | Centralize command output through redacted projections and add known-secret tests for output/logging. | closed |
| Path safety around symlinks and canonical paths can be subtle. | `reasoning.md#assumptions-and-risks` | Use `filepath.Clean`, absolute/canonical comparisons where feasible, temp-dir tests, and conservative rejection when resolution is ambiguous. | closed |
| Health JSON may become consumed before runtime health exists. | `reasoning.md#assumptions-and-risks` | Keep stable top-level envelope and append checks later rather than renaming existing fields. | closed |
| App dependency container may become a god object. | `reasoning.md#trade-off-and-debt-analysis` | Keep it concrete and limited to Sprint 2 services; do not add study/runtime state. | closed |
| Missing area-specific architecture reasoning artifact may confuse governance review. | `reasoning.md#area-specific-reasoning-inputs` | Treat consolidated `reasoning.md` as final architecture decision source unless governance explicitly requires a backfill artifact. | closed |
| Runtime/cancellation/partial exit codes are reserved but not exercised in Sprint 2. | `reasoning.md#decision-7-error-classification-and-exit-codes` | Define the mapping now and test only in-scope classes; avoid runtime behavior. | closed |

## Review Inputs

Review should use:

- `projects/ultraplan-go/sprints/02-workspace-config-health/sprint-index.md`
- `projects/ultraplan-go/sprints/02-workspace-config-health/technical-handbook.md`
- `projects/ultraplan-go/sprints/02-workspace-config-health/reasoning.md`
- `projects/ultraplan-go/sprints/02-workspace-config-health/plan.md`
- Implementation diff in `../ultraplan-go`
- Verification evidence from `go test ./...` and `go build ./cmd/ultraplan`
- Architecture Review protocol from `system/protocols/architecture-review-protocol.md`
- Sprint Review protocol from `system/protocols/sprint-review-protocol.md`

## Execution Log

| Date / Step | Action | Evidence / Notes |
| --- | --- | --- |
| 2026-05-30 / planning | Created evidence-grounded sprint plan from `reasoning.md`. | No implementation performed; area reasoning directory was absent and is recorded as omitted. |
| 2026-05-30 / implementation | Implemented app command routing, workspace discovery/init/path/validation helpers, effective config loading/redaction, minimal logging, `init-workspace`, `config show`, and `health`. | Changed files are in `../ultraplan-go/internal/app`, `../ultraplan-go/internal/workspace`, `../ultraplan-go/internal/platform/config`, and `../ultraplan-go/internal/platform/logging`; `cmd/ultraplan/main.go` remains thin and unchanged. |
| 2026-05-30 / tests | Ran `go test ./...` from `/home/antonioborgerees/coding/ultraplan-go`. | Passed; no OpenCode, provider credentials, network, or configured workspace required. |
| 2026-05-30 / build | Ran `go build -o /tmp/ultraplan-check ./cmd/ultraplan` from `/home/antonioborgerees/coding/ultraplan-go`. | Passed; generated binary was outside the repo. |
| 2026-05-30 / manual CLI review | Ran temp-workspace checks for `init-workspace --dry-run`, `init-workspace`, `config show --json`, and `health --json`; validated JSON with `python3 -m json.tool`. | Passed; dry-run wrote no files, init created only required scaffold, config/health JSON was valid, and health marked runtime as skipped/out of scope. |
| 2026-05-30 / scope review | Searched Go implementation for excluded behavior using `rg` terms covering study workflows, runtime execution, OpenCode/provider checks, durable logs/events, target, and sprint. | No out-of-scope implementation found; only placeholders/doc comments, starter config values, and explicit health skipped runtime check were present. |

## Completion Criteria

- [x] All tasks are complete or explicitly deferred with evidence.
- [x] `go test ./...` was run from `../ultraplan-go` or a deferral/blocker is documented.
- [x] `go build ./cmd/ultraplan` was run from `../ultraplan-go` or a deferral/blocker is documented.
- [x] Evidence satisfies `reasoning.md#expected-evidence`.
- [x] Architecture Review can evaluate package ownership, dependency direction, platform/product separation, and thin entrypoint without guessing intent.
- [x] Sprint Review can evaluate acceptance criteria, tests, build, output modes, redaction, exit codes, and non-goals without guessing intent.
- [x] `review.md` can be created from this plan, the implementation diff, and verification evidence.
