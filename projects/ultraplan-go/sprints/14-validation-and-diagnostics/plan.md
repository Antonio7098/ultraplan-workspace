# Sprint Plan: Validation Command, Diagnostics, and JSON Stability

> Project: `ultraplan-go`
> Sprint: `14-validation-and-diagnostics`
> Source: `projects/ultraplan-go/sprints/14-validation-and-diagnostics/reasoning.md`
> **Inputs Used:** `projects/ultraplan-go/sprints/14-validation-and-diagnostics/requirements.md`, `projects/ultraplan-go/sprints/14-validation-and-diagnostics/reasoning.md`, `projects/ultraplan-go/sprints/14-validation-and-diagnostics/sprint-index.md`, `projects/ultraplan-go/sprints/14-validation-and-diagnostics/technical-handbook.md`, `projects/ultraplan-go/docs/PRD.md`, `projects/ultraplan-go/docs/TRD.md`, `projects/ultraplan-go/docs/ARCHITECTURE.md`, `projects/ultraplan-go/project-index.md`, `templates/sprint-plan.md`, `templates/sprint-reasoning.md`

This plan executes `projects/ultraplan-go/sprints/14-validation-and-diagnostics/reasoning.md`. It must not invent architecture, scope, or decisions.

## Reasoning Source

- **Sprint Reasoning:** `projects/ultraplan-go/sprints/14-validation-and-diagnostics/reasoning.md`
- **Sprint Index:** `projects/ultraplan-go/sprints/14-validation-and-diagnostics/sprint-index.md`
- **Technical Handbook:** `projects/ultraplan-go/sprints/14-validation-and-diagnostics/technical-handbook.md`
- **Area Reasoning:** none present; reasoning document records that `projects/ultraplan-go/sprints/14-validation-and-diagnostics/reasoning/*.md` is absent

## Sprint Status

- **Status:** `complete`
- **Owner:** `implementation agent`
- **Start Date:** `2026-06-13`
- **Completion Date:** `2026-06-13`

## Decisions To Execute

| Decision | Source Section | Execution Implication |
| --- | --- | --- |
| Study-owned validation result model | `reasoning.md#decision-1-study-owned-validation-result-model` | Add or extend typed validation result structures in `/home/antonioborgerees/coding/ultraplan-go/internal/study/domain.go`; include schema/result version, overall status, summary counts, check name/ID, status, severity, path, expected condition, observed safe detail, guidance, source/dimension/run-state context where applicable, and redacted diagnostics. |
| Runtime-free study validation command | `reasoning.md#decision-2-runtime-free-study-validation-command` | Add thin CLI wiring in `/home/antonioborgerees/coding/ultraplan-go/internal/app/study_commands.go` for `ultraplan study <study> validate` and `--json`; implement orchestration in `/home/antonioborgerees/coding/ultraplan-go/internal/study/validation_command.go`; do not invoke runtime, network, source cloning, prompts, synthesis, repair, retry, fallback, or subprocess paths. |
| Stable JSON envelope in `internal/app` | `reasoning.md#decision-3-stable-json-envelope-in-internalapp` | Add `/home/antonioborgerees/coding/ultraplan-go/internal/app/json_output.go` as a payload-agnostic renderer for `schema_version`, `command`, `workspace`, `status`, `generated_at`, and `result`; keep command-specific schema fields in owning payloads. |
| Status JSON from durable local state | `reasoning.md#decision-4-status-json-from-durable-local-state` | Extend `ultraplan study <study> status --json` in `/home/antonioborgerees/coding/ultraplan-go/internal/app/study_commands.go`; summarize durable/local state, task counts, locks, run metadata, validation summaries, retry/fallback/cancellation metadata, artifacts, and unknown usage/cost without runtime invocation. |
| Health JSON through existing runtime health boundary | `reasoning.md#decision-5-health-json-through-existing-runtime-health-boundary` | Extend `/home/antonioborgerees/coding/ultraplan-go/internal/app/health_commands.go` with `--json`; use existing platform runtime/agentwrap health boundary and product-agnostic runtime summaries; do not import study semantics into platform runtime. |
| Redaction before text or JSON rendering | `reasoning.md#decision-6-redaction-before-text-or-json-rendering` | Update `/home/antonioborgerees/coding/ultraplan-go/internal/platform/config/redaction.go` or adjacent existing redaction policy so validation, status, and health diagnostics are safe before rendering; avoid secrets, prompt bodies, Markdown bodies, unsafe raw payload bytes, and full stderr. |
| Fake-first schema and behavior tests | `reasoning.md#decision-7-fake-first-schema-and-behavior-tests` | Add deterministic offline tests for validation service, validate command, status JSON, health JSON, redaction, exit codes, stable JSON fields, and no-runtime invocation; run required verification commands from `/home/antonioborgerees/coding/ultraplan-go`. |

## Requirements / Contracts To Satisfy

| Contract / Requirement ID | Required Behavior | Evidence Planned |
| --- | --- | --- |
| AC31, TRD 9A, TRD 15, Persistence And Migrations | Validation checks study structure, sources, dimensions, applicable per-source reports, final reports, `summary.csv`, and optional `studies/<study>/.ultraplan/run-state.json`; malformed or unsupported state is diagnosed, not migrated. | `internal/study/validation_command_test.go` fixtures for valid, missing, invalid, malformed, inapplicable, malformed state, and unsupported state cases. |
| AC32, TRD 9A.3, PRD 2.5.1 | Inapplicable Markdown source/dimension pairs are skipped or reported as inapplicable, not missing failures. | Validation tests for Markdown frontmatter applicability and status/summary behavior. |
| AC33, AC43, Errors, Observability | Failures include check name, severity, artifact path where applicable, expected condition, observed safe detail, and corrective guidance. | Validation model tests and command output tests assert structured fields and guidance. |
| AC34, TRD 7.2, CLI Surface | Validation returns exit code `0` when checks pass and exit code `5` for validation failures while preserving other error categories. | `internal/app/study_validate_commands_test.go` asserts exit codes and distinct non-validation behavior. |
| AC35, C66, LLM Runtime, Performance | Validation and status do not invoke agentwrap, OpenCode, network, source cloning, runtime execution, synthesis, repair, prompt generation, retry, fallback, or subprocess paths. | Command tests with nil/fake runtime seams assert no runtime invocation; import/package review. |
| AC36-40, AC42, TRD 7.3, Documentation, CLI Surface | `validate --json`, `status --json`, and `health --json` emit parseable, ANSI-free JSON with stable envelope fields and deterministic text behavior. | Command tests parse JSON and assert envelope/result fields, no ANSI, no mixed human text, and updated help text. |
| AC37, AC41, AC43, Configuration, Security | Health JSON includes stable health-check IDs/statuses, runtime/config/workspace summaries, redacted diagnostics, and guidance. | `internal/app/health_commands_test.go` covers success/failure, runtime-health diagnostics, schema fields, guidance, and redaction. |
| AC38, C68, Workflows, LLM Evaluation / Cost / Safety | Status JSON exposes run counts, task summaries, locks, validation summaries, retry/fallback metadata, run metadata, and unknown usage/cost as unknown/null/known flags, never invented zero. | `internal/app/study_status_commands_test.go` parses JSON and asserts counts, metadata, locks, validation summaries, and unknown usage/cost representation. |
| AC41, C69-70, Security | All visible diagnostics are redacted before text or JSON rendering; paths prefer workspace-relative output. | Redaction tests seed API-key-like, env/config secret, runtime stderr, prompt, and Markdown-like body values and assert absence. |
| AC45-46, Architecture | Study behavior remains in `internal/study`, command parsing/rendering in `internal/app`, runtime product-agnostic in `internal/platform/runtime`; no global `internal/validation`, `internal/reports`, `internal/diagnostics`, `internal/status`, or workflow-engine package. | Architecture Review protocol, import/package review, and implementation diff review. |
| AC47-49, Testing | Offline tests, race tests, and CLI build pass. | `go test ./...`, `go test -race ./...`, and `go build ./cmd/ultraplan` from `/home/antonioborgerees/coding/ultraplan-go`. |

## Tasks

- [x] **Task 1: Baseline Current Implementation**
  > Executes: `Decision 1-7`, AC45-46, C63-75
  - [x] Inspect existing `/home/antonioborgerees/coding/ultraplan-go/internal/study/domain.go`, validation/report/status/run-state helpers, and tests to identify reusable validators and state summaries.
  - [x] Inspect existing `/home/antonioborgerees/coding/ultraplan-go/internal/app/study_commands.go`, `/home/antonioborgerees/coding/ultraplan-go/internal/app/health_commands.go`, and command test seams.
  - [x] Inspect `/home/antonioborgerees/coding/ultraplan-go/internal/platform/config/redaction.go` to understand the existing redaction API before extending it.
  - [x] Record any implementation-only deviations or missing seams in `projects/ultraplan-go/sprints/14-validation-and-diagnostics/review.md` during execution; do not change sprint scope without updating reasoning/review.

- [x] **Task 2: Add Study Validation Result Model**
  > Executes: `Decision 1`, AC31-34, AC36, AC41-45
  - [x] Define minimal typed validation result and finding structures in `internal/study/domain.go` or the existing study domain location.
  - [x] Include stable result schema/version fields, overall status, summary counts, check ID/name, status, severity, path, expected condition, safe observed detail, guidance, and optional source/dimension/run-state context.
  - [x] Use explicit statuses for pass, fail, warning, skipped, and inapplicable where needed by validation/status JSON.
  - [x] Keep the model study-owned and avoid app-local duplicates of study semantics.

- [x] **Task 3: Implement Study Validation Service**
  > Executes: `Decision 2`, AC31-35, AC40-45, C71-72
  - [x] Add `/home/antonioborgerees/coding/ultraplan-go/internal/study/validation_command.go` with bounded validation orchestration over known study artifacts.
  - [x] Validate study structure, source discovery, dimension files, expected applicable per-source reports, final reports, and `summary.csv` using existing study/report helpers where possible.
  - [x] Respect Markdown source applicability and report inapplicable pairs as inapplicable or skipped rather than missing failures.
  - [x] Validate optional `studies/<study>/.ultraplan/run-state.json` strictly; diagnose missing, malformed, and unsupported versions without migration.
  - [x] Ensure validation never calls runtime, network, source cloning, prompt generation, synthesis, repair, retry, fallback, or subprocess paths.
  - [x] Construct all observed details as safe summaries, not raw report bodies, prompt bodies, Markdown document bodies, unsafe bytes, or full stderr.

- [x] **Task 4: Wire Validate Command And Help**
  > Executes: `Decision 2`, AC34-36, AC40, AC44, Documentation, CLI Surface
  - [x] Extend `/home/antonioborgerees/coding/ultraplan-go/internal/app/study_commands.go` with `ultraplan study <study> validate` and `--json`.
  - [x] Keep command handler thin: parse flags, call study validation service, render text/JSON, and map validation failures to exit code `5`.
  - [x] Add deterministic concise text output with actionable guidance and no unrelated behavior changes.
  - [x] Update help text to document `study <study> validate`, `study <study> validate --json`, and `study <study> status --json`.

- [x] **Task 5: Add Shared JSON Envelope Renderer**
  > Executes: `Decision 3`, AC36-40, AC42-44, C67-70
  - [x] Add `/home/antonioborgerees/coding/ultraplan-go/internal/app/json_output.go` with a payload-agnostic envelope renderer.
  - [x] Emit `schema_version`, `command`, `workspace`, `status`, `generated_at`, and `result` for JSON mode.
  - [x] Ensure JSON mode writes only valid JSON to stdout and no ANSI, progress text, or human-only diagnostics.
  - [x] Keep helper free of imports from `internal/study` and `internal/platform/runtime` unless existing app package structure already imports payload types elsewhere.

- [x] **Task 6: Extend Status JSON**
  > Executes: `Decision 4`, AC38-45, C68, C74
  - [x] Extend `ultraplan study <study> status --json` in `internal/app/study_commands.go` using the shared envelope.
  - [x] Build or reuse study-owned/local summaries for run counts, task states, artifact status, validation summaries, locks, run metadata, cancellation state, retry/fallback metadata, and safe diagnostics.
  - [x] Represent unknown usage/cost with explicit unknown/null/known flags and never numeric zero unless the value is known zero.
  - [x] Keep status runtime-free and bounded to durable/local study artifacts.

- [x] **Task 7: Extend Health JSON**
  > Executes: `Decision 5`, AC37, AC39-45, C66, C69, C73
  - [x] Extend `/home/antonioborgerees/coding/ultraplan-go/internal/app/health_commands.go` with `--json` using the shared envelope.
  - [x] Include stable health-check IDs, statuses, runtime/config/workspace summaries, safe diagnostics, and guidance.
  - [x] Use existing platform runtime/agentwrap health seams and fakes; do not duplicate OpenCode checks in app code.
  - [x] Preserve platform runtime product-agnostic boundaries and prevent study imports into runtime packages.

- [x] **Task 8: Extend Redaction Policy For New Surfaces**
  > Executes: `Decision 6`, AC33, AC37-43, C69-70, C73-75
  - [x] Update `/home/antonioborgerees/coding/ultraplan-go/internal/platform/config/redaction.go` or adjacent existing redaction policy to support validation, status, and health diagnostics.
  - [x] Apply redaction before both text and JSON rendering.
  - [x] Cover nested diagnostic strings and safe excerpts without introducing a broad diagnostics framework.
  - [x] Prefer workspace-relative paths in output and avoid exposing secret-bearing absolute paths.

- [x] **Task 9: Add Validation Service Tests**
  > Executes: `Decision 1`, `Decision 2`, `Decision 6`, `Decision 7`, AC31-35, AC41-42
  - [x] Add `/home/antonioborgerees/coding/ultraplan-go/internal/study/validation_command_test.go`.
  - [x] Cover valid study validation, missing artifacts, invalid reports, malformed frontmatter, inapplicable Markdown pairs, run-state diagnostics, malformed/unsupported state, and secret-free diagnostics.
  - [x] Assert structured findings include stable check names, severity/status, paths, expected/observed safe details, and guidance.

- [x] **Task 10: Add Validate Command Tests**
  > Executes: `Decision 2`, `Decision 3`, `Decision 6`, `Decision 7`, AC34-36, AC39-44
  - [x] Add `/home/antonioborgerees/coding/ultraplan-go/internal/app/study_validate_commands_test.go`.
  - [x] Cover validate text/JSON output, exit code `0`, exit code `5`, checkable failures, redaction, stable envelope/result fields, no ANSI in JSON, and no runtime invocation.
  - [x] Assert JSON diagnostics are fields, not mixed human text.

- [x] **Task 11: Extend Status And Health Tests**
  > Executes: `Decision 3-7`, AC37-44, C68-69, C73-74
  - [x] Extend or add `/home/antonioborgerees/coding/ultraplan-go/internal/app/study_status_commands_test.go` for stable status JSON shape, counts, task states, lock diagnostics, run metadata summaries, validation summaries, retry/fallback metadata, unknown usage/cost, redaction, and no runtime invocation.
  - [x] Extend `/home/antonioborgerees/coding/ultraplan-go/internal/app/health_commands_test.go` for `health --json` success/failure, health-check IDs/statuses, runtime-health diagnostics, schema fields, guidance, redaction, and fake runtime behavior.
  - [x] Keep tests deterministic, offline, fake-first, and independent of OpenCode, provider credentials, network access, source cloning, and long sleeps.

- [x] **Task 12: Run Verification And Prepare Review Evidence**
  > Executes: AC47-49, Architecture Review, Sprint Review
  - [x] Run `go test ./...` from `/home/antonioborgerees/coding/ultraplan-go`.
  - [x] Run `go test -race ./...` from `/home/antonioborgerees/coding/ultraplan-go`.
  - [x] Run `go build ./cmd/ultraplan` from `/home/antonioborgerees/coding/ultraplan-go`.
  - [x] Review imports and package additions for architecture boundary conformance and forbidden package absence.
  - [x] Record implementation evidence, deviations, verification status, and review conclusions in `projects/ultraplan-go/sprints/14-validation-and-diagnostics/review.md`.

## Evidence Checklist

- [x] Tests prove study validation behavior for valid, missing, invalid, malformed, inapplicable, run-state, unsupported-state, and secret-free cases.
- [x] Command tests parse validation/status/health JSON and assert stable envelope/result fields.
- [x] Command tests assert validation exit code `0` for success and `5` for validation failures.
- [x] Tests prove validation and status do not invoke runtime seams.
- [x] Tests prove health JSON uses fakeable runtime-health seams through the existing platform boundary.
- [x] JSON output has no ANSI, progress output, or mixed human diagnostics.
- [x] Text output remains deterministic, concise, actionable, and backward-compatible unless tests document a deliberate sprint change.
- [x] Redaction tests cover validation, status, and health text/JSON.
- [x] Help text documents `study <study> validate`, `study <study> validate --json`, and `study <study> status --json`.
- [x] Architecture Review evidence confirms study/app/platform ownership, dependency direction, and no forbidden global packages.
- [x] Sprint Review evidence is recorded in `projects/ultraplan-go/sprints/14-validation-and-diagnostics/review.md`.

## Verification Commands

| Check | Command | Expected Result |
| --- | --- | --- |
| Offline unit and command tests | `go test ./...` | Passes from `/home/antonioborgerees/coding/ultraplan-go` without OpenCode, provider credentials, network access, source cloning, or long sleeps. |
| Race verification | `go test -race ./...` | Passes from `/home/antonioborgerees/coding/ultraplan-go`. |
| CLI build | `go build ./cmd/ultraplan` | Builds successfully from `/home/antonioborgerees/coding/ultraplan-go`. |

## Risks And Blockers

| Risk / Blocker | Source | Mitigation | Status |
| --- | --- | --- | --- |
| Existing report validators may be too task-specific to reuse directly. | `reasoning.md#assumptions-and-risks` | Reused existing per-report validators from `internal/study/validation.go`; added only bounded study orchestration. | mitigated |
| Existing status/run-state code may not expose enough metadata for JSON summaries. | `reasoning.md#assumptions-and-risks` | Reused `SummarizeRunState` and shaped app-owned JSON from durable `RunState` fields. | mitigated |
| Existing redaction policy may be string-oriented and insufficient for nested payloads. | `reasoning.md#assumptions-and-risks` | Added boundary redaction helpers for validation/status/health and expanded token-like marker detection. | mitigated |
| JSON schema vocabulary may become breaking once automation consumes it. | `reasoning.md#assumptions-and-risks` | Added envelope/result schema versions and command tests asserting public fields. | mitigated |
| Validation/status could accidentally invoke runtime through shared command factories. | `reasoning.md#assumptions-and-risks` | Added command tests that replace the runtime factory and assert validate/status do not invoke it. | mitigated |
| Safe diagnostics may omit details useful for debugging. | `reasoning.md#assumptions-and-risks` | Findings include check names, paths, expected/observed details, and guidance; raw payload/debug retention remains future scope. | carried forward |
| Unsupported run-state versions may frustrate users without migration. | `reasoning.md#assumptions-and-risks` | Validation emits unsupported-schema diagnostics and guidance without migration, as required. | mitigated |
| Area-specific architecture reasoning artifact is absent. | `reasoning.md#assumptions-and-risks` | Use final architecture decisions in `reasoning.md`; do not reopen architecture in implementation. | mitigated |

## Review Inputs

Review should use:

- `projects/ultraplan-go/sprints/14-validation-and-diagnostics/sprint-index.md`
- `projects/ultraplan-go/sprints/14-validation-and-diagnostics/technical-handbook.md`
- `projects/ultraplan-go/sprints/14-validation-and-diagnostics/reasoning.md`
- `projects/ultraplan-go/sprints/14-validation-and-diagnostics/plan.md`
- implementation diff in `/home/antonioborgerees/coding/ultraplan-go`
- verification evidence from `go test ./...`, `go test -race ./...`, and `go build ./cmd/ultraplan`
- `system/protocols/architecture-review-protocol.md`
- `system/protocols/sprint-review-protocol.md`

## Execution Log

| Date / Step | Action | Evidence / Notes |
| --- | --- | --- |
| 2026-06-12 / planning | Created sprint plan from requirements, reasoning, sprint index, technical handbook, PRD, TRD, Architecture, project index, and sprint plan template. | Planning only; no code implementation performed. |
| 2026-06-13 / implementation | Added study validation result model, study validation service, validate command, shared JSON envelope, status JSON, health JSON, redaction updates, and tests. | Implementation files changed under `/home/antonioborgerees/coding/ultraplan-go/internal/{study,app,platform/config}`. |
| 2026-06-13 / verification | Ran required verification commands. | `go test ./...` passed; `go test -race ./...` passed; `go build ./cmd/ultraplan` passed. |

## Completion Criteria

- [x] All tasks are complete or explicitly deferred with rationale in `review.md`.
- [x] Verification commands were run or deferrals are documented in `review.md`.
- [x] Evidence satisfies the expectations from `reasoning.md`.
- [x] `review.md` can evaluate conformance without guessing intent.
- [x] Architecture Review confirms boundaries and forbidden package constraints.
- [x] Sprint Review confirms outputs, deviations, tests, build status, and review conclusions.
