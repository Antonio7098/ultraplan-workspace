# Sprint Plan: 08 Run State Persistence

> Project: `ultraplan-go`
> Sprint: `08-run-state-persistence`
> Source: `projects/ultraplan-go/sprints/08-run-state-persistence/reasoning.md`
> **Inputs Used:** `projects/ultraplan-go/sprints/08-run-state-persistence/requirements.md`, `projects/ultraplan-go/sprints/08-run-state-persistence/reasoning.md`, `projects/ultraplan-go/sprints/08-run-state-persistence/sprint-index.md`, `projects/ultraplan-go/sprints/08-run-state-persistence/technical-handbook.md`, `projects/ultraplan-go/docs/ARCHITECTURE.md`, `projects/ultraplan-go/docs/PRD.md`, `projects/ultraplan-go/docs/TRD.md`, `projects/ultraplan-go/project-index.md`, `templates/sprint-plan.md`, `/home/antonioborgerees/coding/ultraplan-go/internal/study/domain.go`, `/home/antonioborgerees/coding/ultraplan-go/internal/study/reports.go`, `/home/antonioborgerees/coding/ultraplan-go/internal/study/validation.go`, `/home/antonioborgerees/coding/ultraplan-go/internal/study/applicability.go`, `/home/antonioborgerees/coding/ultraplan-go/internal/study/prompts.go`, `/home/antonioborgerees/coding/ultraplan-go/internal/app/study_commands.go`

This plan executes `reasoning.md`. It must not invent architecture, scope, or decisions.

## Reasoning Source

- **Sprint Reasoning:** `projects/ultraplan-go/sprints/08-run-state-persistence/reasoning.md`
- **Sprint Index:** `projects/ultraplan-go/sprints/08-run-state-persistence/sprint-index.md`
- **Technical Handbook:** `projects/ultraplan-go/sprints/08-run-state-persistence/technical-handbook.md`
- **Area Reasoning:** none present; `projects/ultraplan-go/sprints/08-run-state-persistence/reasoning/*.md` did not exist when this plan was created.
- **Code Context Resolved:** existing study domain, applicability, prompt/report path, validation, and study command files listed in `Inputs Used`.

## Sprint Status

- **Status:** complete
- **Owner:** implementation agent
- **Start Date:** 2026-05-31
- **Completion Date:** 2026-05-31

## Decisions To Execute

| Decision | Source Section | Execution Implication |
| --- | --- | --- |
| Keep run state study-owned with thin status wiring | `reasoning.md#decision-1-keep-run-state-study-owned-with-thin-status-wiring` | Add run-state domain, construction, persistence, validation, resume validation, and summary logic under `internal/study`; keep `internal/app` limited to command dispatch and human status rendering. |
| Use a versioned JSON run-state schema with explicit task metadata | `reasoning.md#decision-2-use-a-versioned-json-run-state-schema-with-explicit-task-metadata` | Extend `internal/study/domain.go` with run-state, task kind/status, filters, config summary, validation summary, dependency, and agentwrap placeholder types; persist JSON at `studies/<study>/.ultraplan/run-state.json`. |
| Save state atomically with same-directory temp files | `reasoning.md#decision-3-save-state-atomically-with-same-directory-temp-files-and-preserve-prior-valid-state-on-failure` | Implement save/load in `internal/study/state.go` using same-directory temp file, flush/close, rename, wrapped diagnostics, and targeted test seams. |
| Build deterministic planned tasks from applicable inputs only | `reasoning.md#decision-4-build-deterministic-planned-analysis-and-synthesis-tasks-from-applicable-inputs-only` | Add `internal/study/run_state.go` helpers that reuse `GetApplicableSources`, existing report path helpers, and deterministic sorting to create analysis and synthesis tasks. |
| Validate loaded state strictly and revalidate completed outputs | `reasoning.md#decision-5-validate-loaded-state-strictly-and-revalidate-completed-outputs-before-trusting-completion` | Distinguish missing, malformed, and unsupported state; revalidate completed source/final reports with existing validation helpers; reset stale active statuses safely while preserving attempts and safe error details. |
| Persist only narrow, secret-free operational configuration summary | `reasoning.md#decision-6-persist-only-narrow-secret-free-operational-configuration-summary` | Keep `ConfigSummary` fields safe and minimal; do not persist credentials, environment dumps, raw runtime payloads, or unsafe provider diagnostics. |
| Render deterministic runtime-free human status from persisted state | `reasoning.md#decision-7-render-deterministic-runtime-free-human-status-from-persisted-state` | Add `ultraplan study <study> status` handling in `internal/app/study_commands.go` or a focused companion command file; read persisted state, summarize, and render concise deterministic text only. |
| Verify with focused state, persistence, resume, status, and architecture tests | `reasoning.md#decision-8-verify-with-focused-state-persistence-resume-status-and-architecture-tests` | Add tests in `internal/study/state_test.go` and `internal/app/study_status_commands_test.go`; run full test and build commands from `/home/antonioborgerees/coding/ultraplan-go`. |

## Requirements / Contracts To Satisfy

| Contract / Requirement ID | Required Behavior | Evidence Planned |
| --- | --- | --- |
| AC 31, TRD 13.2, Persistence And Migrations | State file exists at `studies/<study>/.ultraplan/run-state.json` and includes schema version, run ID, timestamps, filters, config summary, tasks, and completion flag. | `internal/study/state_test.go` construction and JSON round-trip tests. |
| AC 32-33, PRD 2.6, TRD 13.1 | Analysis and synthesis task records include stable IDs, kind, study/dimension/source metadata, output paths, status, attempts, timestamps, errors, retry time, validation summary, and future agentwrap metadata placeholders. | Domain tests and fixture JSON assertions in `internal/study/state_test.go`. |
| AC 34-36, AC 45, TRD 9A.3 | Initial state includes only applicable source/dimension pairs, creates synthesis tasks only for dimensions with applicable inputs, and orders tasks deterministically. | Task construction tests covering directory sources, applicable Markdown sources, and inapplicable Markdown exclusions. |
| AC 37-38, TRD 13.3, PRD 3.5 | State writes use same-directory temp file, flush/close, and rename; failed writes preserve the prior valid state where feasible. | Atomic save/load and injected failure tests in `internal/study/state_test.go`. |
| AC 39-43, Errors, Workflows | Loading rejects unsupported schema versions, distinguishes missing from malformed state, revalidates completed outputs, resets stale active statuses safely, and preserves attempts/errors. | Missing/malformed/unsupported/resume validation tests in `internal/study/state_test.go`. |
| AC 44-46, Observability, CLI Surface | Status summary deterministically counts all required statuses and human output includes state path, run ID, completion flag, totals, failed count, active count, retry count, and next retry time when present. | Summary tests and command output tests. |
| AC 47, Performance, Security | Status command does not invoke runtime work, agentwrap/OpenCode, provider credentials, network calls, subprocess execution, worker pools, or process inspection. | Command tests using existing dependency seams and review of imports/call graph. |
| AC 48-49, Testing | Unit and command tests cover state construction, stable IDs, filtering, synthesis dependencies, persistence, load diagnostics, resume validation, status aggregation, and no-runtime status behavior. | `internal/study/state_test.go`; `internal/app/study_status_commands_test.go`. |
| AC 50-51 | Full Go tests and CLI build pass. | `go test ./...`; `go build ./cmd/ultraplan`. |
| Architecture Constraints 70-72 | Study-owned behavior stays in `internal/study`; CLI rendering stays in `internal/app`; platform packages do not import `study`; no global state/scheduler/report package. | Architecture review and import/package inspection. |
| Security Constraints 75, 80-81 | Persisted paths are workspace-relative where practical; persisted state avoids secrets; tests use local fixtures/temp dirs; errors preserve cause chains and name relevant paths. | State JSON assertions, command diagnostics tests, security review. |

## Tasks

- [x] **Task 1: Define Run-State Domain Types**
  > Executes: Decisions 1, 2, 6; AC 31-33, 44; Architecture, Security
  - [x] Extend `/home/antonioborgerees/coding/ultraplan-go/internal/study/domain.go` with `RunState`, `TaskState`, `TaskKind`, `TaskStatus`, `RunFilters`, `ConfigSummary`, `TaskError`, `ValidationSummary`, `AgentMetadata`, `SynthesisDependency`, and `StatusSummary` concepts.
  - [x] Define task statuses for `pending`, `running`, `validating`, `completed`, `failed`, `cancelled`, `skipped`, `waiting`, and `retrying`; include any existing or needed planning-only status only if mapped into deterministic summaries.
  - [x] Keep all fields JSON-tagged with stable snake_case names and no secret-bearing runtime payload fields.
  - [x] Add constants for current schema version and the study-local state file name/path components.

- [x] **Task 2: Implement Deterministic Run-State Construction**
  > Executes: Decisions 2, 4; AC 34-36, 45; TRD 9A.3, TRD 13.4
  - [x] Add `/home/antonioborgerees/coding/ultraplan-go/internal/study/run_state.go` for initial state construction helpers.
  - [x] Reuse `GetApplicableSources` so inapplicable Markdown pairs are absent from required analysis tasks and missing-output counts.
  - [x] Build stable analysis task IDs from task kind, study name, dimension number/slug, source name, and source kind; avoid timestamps, run IDs, absolute paths, or ordering-dependent values in IDs.
  - [x] Use existing `SourceReportPath` and `FinalReportPath` helpers for output metadata, converting to workspace-relative paths where a workspace root is available.
  - [x] Create synthesis tasks only for dimensions with at least one applicable source and record deterministic dependencies for each required applicable source report.
  - [x] Sort dimensions, sources, dependencies, and tasks deterministically for identical study inputs.

- [x] **Task 3: Implement State Load, Validation, Save, And Atomic Persistence**
  > Executes: Decisions 2, 3, 5; AC 31, 37-40; Persistence And Migrations, Errors
  - [x] Add `/home/antonioborgerees/coding/ultraplan-go/internal/study/state.go` with helpers to resolve `studies/<study>/.ultraplan/run-state.json`.
  - [x] Implement `LoadRunState` with missing-state detection distinct from malformed JSON and unsupported schema versions.
  - [x] Validate required top-level fields and task fields enough to reject unsupported or malformed state with diagnostics naming the state file path.
  - [x] Implement `SaveRunState` using same-directory temporary files, deterministic JSON indentation, file flush/sync where supported, close, and rename over the prior state file.
  - [x] Ensure write, flush, close, and rename failures wrap underlying causes and leave the previous valid state in place where filesystem semantics allow.
  - [x] Add only targeted clock/filesystem seams needed for deterministic timestamps and atomic-write failure tests; do not create a global filesystem abstraction package.

- [x] **Task 4: Implement Resume Validation And Status Summary Logic**
  > Executes: Decisions 4, 5, 7; AC 41-46; Workflows, Observability
  - [x] Revalidate completed analysis tasks through `ValidateSourceReport` and completed synthesis tasks through `ValidateFinalReport` before treating completion as trusted.
  - [x] Reset stale active statuses such as `running`, `validating`, `waiting`, and `retrying` to a safe retryable or pending state without marking them complete.
  - [x] Preserve attempt count, safe `LastError`, retry time, and useful timestamps when resume validation changes task status.
  - [x] Compute deterministic status counts for pending, running, validating, completed, failed, cancelled, skipped, waiting, retrying, and total.
  - [x] Compute active count from active statuses and retry count/next retry time from persisted retry metadata.
  - [x] Ensure status calculations derive from persisted state and known study/report metadata, not live process inspection.

- [x] **Task 5: Wire Runtime-Free Study Status Command**
  > Executes: Decisions 1, 7; AC 46-47, 49; CLI Surface, Performance
  - [x] Add `status` handling for `ultraplan study <study> status` in `/home/antonioborgerees/coding/ultraplan-go/internal/app/study_commands.go` or a focused `/home/antonioborgerees/coding/ultraplan-go/internal/app/study_status_commands.go` companion file.
  - [x] Keep the command thin: discover workspace, resolve study through existing service/listing behavior, load run state, call study-owned validation/summary helpers, and render output.
  - [x] Render deterministic human output with state file path, run ID, completion flag, total tasks, failed count, active count, retry count, and next retry time when present.
  - [x] Distinguish missing state from malformed/unsupported state in user diagnostics and include the state file path.
  - [x] Update `studyHelp()` to list `ultraplan study <study> status`.
  - [x] Do not initialize runtime dependencies, call agentwrap/OpenCode, inspect processes, start worker pools, call network APIs, or read provider credentials.

- [x] **Task 6: Add State And Persistence Tests**
  > Executes: Decision 8; AC 48, 50; Testing, Persistence And Migrations
  - [x] Add `/home/antonioborgerees/coding/ultraplan-go/internal/study/state_test.go` covering schema fields, stable IDs, deterministic order, applicability filtering, synthesis dependencies, and secret-free config summary.
  - [x] Test atomic save/load round trip, malformed JSON, unsupported schema version, missing state, and failed atomic write preserving a prior valid state where feasible.
  - [x] Test stale active status reset, completed output revalidation success/failure, attempt preservation, last-error preservation, status aggregation, retry count, and next retry time.
  - [x] Use temporary directories and small fixture studies; do not require OpenCode, provider credentials, network access, real source repositories, or long-running subprocesses.

- [x] **Task 7: Add Study Status Command Tests**
  > Executes: Decisions 7, 8; AC 49; CLI Surface, Performance
  - [x] Add `/home/antonioborgerees/coding/ultraplan-go/internal/app/study_status_commands_test.go` covering missing state, valid state, malformed/unsupported state, failed tasks, retry next time, deterministic output, and no-runtime behavior.
  - [x] Reuse existing app command test patterns from `/home/antonioborgerees/coding/ultraplan-go/internal/app/study_commands_test.go` and prompt command tests.
  - [x] Assert diagnostics include the state file path and do not depend on nondeterministic timestamps except fixed fixture values.

- [x] **Task 8: Run Verification And Architecture Review Checks**
  > Executes: Decision 8; AC 50-51; Architecture, Testing
  - [x] Run `go test ./...` from `/home/antonioborgerees/coding/ultraplan-go`.
  - [x] Run `go build ./cmd/ultraplan` from `/home/antonioborgerees/coding/ultraplan-go`.
  - [x] Review package paths/imports to confirm run-state logic stayed in `internal/study`, CLI status rendering stayed in `internal/app`, platform packages do not import `study`, and no global state/scheduler/report package was introduced.
  - [x] Record verification command results, architecture findings, acceptance status, and any carry-forward risks in `review.md` during sprint review.

## Evidence Checklist

- [x] Tests prove run-state construction includes all required schema fields and task metadata.
- [x] Tests prove stable task IDs and deterministic task ordering for identical study inputs.
- [x] Tests prove inapplicable Markdown pairs are absent from required analysis tasks and excluded from status missing-output counts.
- [x] Tests prove synthesis tasks are created only for dimensions with applicable inputs and record deterministic dependencies.
- [x] Tests prove atomic save/load behavior and prior valid state preservation on failed writes where feasible.
- [x] Tests prove missing, malformed, and unsupported state are distinguishable and path-aware.
- [x] Tests prove resume validation revalidates completed outputs, resets stale active statuses safely, and preserves attempts/errors.
- [x] Tests prove status aggregation counts pending, running, validating, completed, failed, cancelled, skipped, waiting, retrying, and total.
- [x] Command tests prove deterministic human status output and no runtime behavior.
- [x] `go test ./...` passes from `/home/antonioborgerees/coding/ultraplan-go`.
- [x] `go build ./cmd/ultraplan` passes from `/home/antonioborgerees/coding/ultraplan-go`.
- [x] Architecture review evidence confirms package ownership and dependency direction.
- [x] Sprint review records any deviations from `reasoning.md` before implementation continues.

## Verification Commands

| Check | Command | Expected Result |
| --- | --- | --- |
| Full test suite | `go test ./...` | All packages pass without OpenCode, provider credentials, network access, real source repositories, or long-running subprocesses. |
| CLI build | `go build ./cmd/ultraplan` | Command builds successfully. |
| Architecture review | Manual import/path review using `system/protocols/architecture-review-protocol.md` | Run-state behavior is in `internal/study`, status rendering is in `internal/app`, no platform package imports `study`, and no global state/scheduler/report package exists. |
| Sprint review | Manual review using `system/protocols/sprint-review-protocol.md` | Required outputs, acceptance criteria, verification evidence, scope exclusions, risks, and carry-forward decisions are recorded in `review.md`. |

## Risks And Blockers

| Risk / Blocker | Source | Mitigation | Status |
| --- | --- | --- | --- |
| Existing report validation helpers may not expose all detail wanted for resume validation summaries. | `reasoning.md#assumptions-and-risks` | Used existing `ValidationResult` and adapter summaries only. | mitigated |
| Existing prompt/report path helpers may not cover every desired workspace-relative persisted path. | `reasoning.md#assumptions-and-risks`; code context in `reports.go` and `prompts.go` | Reused `SourceReportPath`, `FinalReportPath`, and `workspace.Rel` via study-local helper. | mitigated |
| Config summary may be narrower than final runtime diagnostics. | `reasoning.md#decision-6-persist-only-narrow-secret-free-operational-configuration-summary` | Persisted only safe known operational facts; redacted runtime metadata remains deferred. | carried forward |
| Atomic rename semantics vary by filesystem/platform. | `reasoning.md#decision-3-save-state-atomically-with-same-directory-temp-files-and-preserve-prior-valid-state-on-failure` | Implemented same-directory temp, flush, close, rename, and best-effort directory sync; residual platform semantics carried forward. | mitigated |
| Concurrent commands may race because locks are out of scope. | `reasoning.md#assumptions-and-risks` | No multi-writer safety claimed; locks remain deferred. | carried forward |
| Inapplicable pairs absent from tasks may surprise users expecting a full matrix. | `reasoning.md#decision-4-build-deterministic-planned-analysis-and-synthesis-tasks-from-applicable-inputs-only` | Tests prove inapplicable pairs are absent from task totals. | mitigated |
| Stale active status reset policy may need refinement after real cancellation/retry policies exist. | `reasoning.md#decision-5-validate-loaded-state-strictly-and-revalidate-completed-outputs-before-trusting-completion` | Implemented safe reset to pending while preserving attempts, last error, and retry metadata. | mitigated |
| Manual edits to `run-state.json` may produce unsupported or malformed state. | `reasoning.md#assumptions-and-risks` | Added path-aware missing/malformed/unsupported diagnostics and tests. | mitigated |

## Review Inputs

Review should use:

- `projects/ultraplan-go/sprints/08-run-state-persistence/requirements.md`
- `projects/ultraplan-go/sprints/08-run-state-persistence/sprint-index.md`
- `projects/ultraplan-go/sprints/08-run-state-persistence/technical-handbook.md`
- `projects/ultraplan-go/sprints/08-run-state-persistence/reasoning.md`
- `projects/ultraplan-go/sprints/08-run-state-persistence/plan.md`
- implementation diff
- verification evidence from `go test ./...` and `go build ./cmd/ultraplan`
- `system/protocols/architecture-review-protocol.md`
- `system/protocols/sprint-review-protocol.md`

## Execution Log

| Date / Step | Action | Evidence / Notes |
| --- | --- | --- |
| 2026-05-31 planning | Created implementation plan from requirements, reasoning, sprint index, technical handbook, PRD/TRD/Architecture, project index, plan template, and focused code context. | No code implemented. Area reasoning directory was absent. |
| 2026-05-31 implementation | Added run-state schema, deterministic construction, atomic persistence, resume validation, status summaries, and `ultraplan study <study> status`. | Implemented in `internal/study/domain.go`, `internal/study/run_state.go`, `internal/study/state.go`, and `internal/app/study_commands.go`. |
| 2026-05-31 tests | Added state/persistence tests and study status command tests. | `go test ./internal/study ./internal/app` passed. |
| 2026-05-31 verification | Ran required full verification and review checks. | `go test ./...` passed. `go build ./cmd/ultraplan` passed. Architecture review found run-state behavior in `internal/study`, rendering in `internal/app`, no new global state/scheduler/report package, and no platform package importing `study`. |

## Completion Criteria

- [x] All tasks are complete or explicitly deferred with requirement impact recorded.
- [x] Verification commands were run or deferrals are documented with cause.
- [x] Evidence satisfies the expectations from `reasoning.md` Decisions 1-8.
- [x] `review.md` can evaluate conformance without guessing intent.
- [x] No runtime execution, worker pool orchestration, retry/fallback execution, cancellation handling, locks, summary CSV generation, code extraction, target workflows, sprint planning, or sprint execution scope was added.
