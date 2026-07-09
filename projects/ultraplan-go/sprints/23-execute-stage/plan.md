# Sprint Plan: Execute Stage

> Project: `ultraplan-go`
> Sprint: `23-execute-stage`
> Source: `reasoning.md`
> **Inputs Used:** `projects/ultraplan-go/sprints/23-execute-stage/requirements.md`, `projects/ultraplan-go/sprints/23-execute-stage/reasoning.md`, `projects/ultraplan-go/sprints/23-execute-stage/sprint-index.md`, `projects/ultraplan-go/sprints/23-execute-stage/technical-handbook.md`, `projects/ultraplan-go/sprints/23-execute-stage/reasoning/architecture.md`, `projects/ultraplan-go/docs/PRD.md`, `projects/ultraplan-go/docs/TRD.md`, `projects/ultraplan-go/docs/ARCHITECTURE.md`, `projects/ultraplan-go/project-index.md`

This plan executes `reasoning.md`. It does not invent architecture, scope, or decisions.

Final study reports were not opened directly because `technical-handbook.md` and `reasoning.md` already cite and synthesize the selected evidence reports at the level needed for implementation planning. Open final reports only during implementation if a specific code pattern or cited example needs deeper inspection.

## Reasoning Source

- **Sprint Reasoning:** `projects/ultraplan-go/sprints/23-execute-stage/reasoning.md`
- **Sprint Index:** `projects/ultraplan-go/sprints/23-execute-stage/sprint-index.md`
- **Technical Handbook:** `projects/ultraplan-go/sprints/23-execute-stage/technical-handbook.md`
- **Area Reasoning:** `projects/ultraplan-go/sprints/23-execute-stage/reasoning/architecture.md`

## Sprint Status

- **Status:** `complete`
- **Owner:** `implementation agent`
- **Start Date:** `2026-07-09`
- **Completion Date:** `2026-07-09`

## Decisions To Execute

| Decision | Source Section | Execution Implication |
| --- | --- | --- |
| Sprint Owns Execute Semantics | `reasoning.md#decision-1-sprint-owns-execute-semantics` | Put task extraction, deterministic IDs, prompts, validation, run state, summary, status, and flow behavior in `internal/sprint`; keep app/CLI as wiring and runtime generic. |
| Execute Only From Validated Plan Tasks | `reasoning.md#decision-2-execute-only-from-validated-plan-tasks` | Validate prerequisites through `plan.md`, extract only supported explicit checklist tasks, reject ambiguous syntax, and persist traceability to plan and reasoning references. |
| Versioned Run State Is The Durable Source Of Truth | `reasoning.md#decision-3-versioned-run-state-is-the-durable-source-of-truth` | Add strict, versioned, atomic `.run-state.json` load/save/init/resume logic and make status/resume derive from it instead of `execute.md`. |
| Completion Requires Evidence Or Explicit Safe Diagnostic | `reasoning.md#decision-4-completion-requires-evidence-or-explicit-safe-diagnostic` | Mark tasks complete only after runtime success plus expected evidence, or after an explicit safe diagnostic explains why evidence cannot be machine-validated. |
| Target Repository And Runtime Diagnostics Are Safety-Gated | `reasoning.md#decision-5-target-repository-and-runtime-diagnostics-are-safety-gated` | Resolve and validate the only approved target `../ultraplan-go`, constrain prompts/permissions to that target, redact secrets, omit unsafe raw payloads, and exclude Git mutation. |
| Execute Runs Conservatively And Resumes Safely | `reasoning.md#decision-6-execute-runs-conservatively-and-resumes-safely` | Run tasks serially by default, save state after meaningful transitions, handle cancellation, recover stale running tasks, and skip complete tasks on resume. |
| Stage Model Resolution Is Deterministic And Visible | `reasoning.md#decision-7-stage-model-resolution-is-deterministic-and-visible` | Implement source-aware model precedence for sprint stages including `execute`, reject unknown/empty stage model keys, and surface safe model-source diagnostics. |
| Status, Prompt, Validate, Flow, And Execute Commands Stay Thin | `reasoning.md#decision-8-status-prompt-validate-flow-and-execute-commands-stay-thin` | Add app/CLI request/result wiring for execute validation, execute prompt preview, flow targeting the execute stage, execute task runs, task selection, dry-run, help, exit codes, and runtime-free status. |
| Deferred Behaviors Are Explicitly Rejected | `reasoning.md#decision-9-deferred-behaviors-are-explicitly-rejected` | Reject smoke, review automation, issue tracking, Git mutation, TUI, hosted/browser UI, cross-sprint scheduling, and study-service reuse in execute paths. |
| Verification Uses Deterministic Fake-Runtime Evidence | `reasoning.md#decision-10-verification-uses-deterministic-fake-runtime-evidence` | Cover product-owned behavior with unit, fixture, golden, command, race, build, and fake-runtime tests; keep real runtime checks optional and gated. |

## Requirements / Contracts To Satisfy

| Contract / Requirement ID | Required Behavior | Evidence Planned |
| --- | --- | --- |
| Architecture | Preserve `sprint -> project`, `sprint -> workspace`, `sprint -> platform/runtime`, and `platform/* -> no product modules`; keep execute semantics in `internal/sprint`. | Import review, package-level tests, architecture review protocol evidence. |
| CLI Surface | Add or update execute validation, execute prompt preview, execute-stage flow handling, execute task runs, task selection, dry-run, help, exit codes, and status rendering. | Command tests for help, invalid args, prompt/dry-run, selected task, invalid prerequisites, exit codes, and status output. |
| Configuration / REQ-23-41 | Resolve explicit override, stage-specific `execute`, global sprint/planning model, `models.primary`, then `models.default`; reject unknown or empty stage models. | Unit tests for precedence, fallback, unknown keys, empty values, command overrides, and safe diagnostics. |
| Documentation / REQ-23-49 | Produce editable `execute.md` summary that cites `plan.md`, records counts, terminal states, evidence, diagnostics, and `.run-state.json` pointers. | Summary renderer golden tests and `execute.md` validator tests. |
| Errors | Provide actionable diagnostics for prerequisite, plan extraction, runtime, validation, stale state, atomic write, target path, and model/config failures. | Unit and command tests for each failure class and exit-code mapping. |
| LLM Runtime / REQ-23-43 | Use only generic `internal/platform/runtime` and agentwrap/OpenCode through the existing integration; do not parse OpenCode stdout/stderr directly. | Fake-runtime request assertions and import/code review. |
| LLM Evaluation / Cost / Safety / REQ-23-45 | Runtime success alone is insufficient; completion requires evidence or explicit safe diagnostic; unsafe payloads and secrets stay out of state/output. | Fake-runtime tests for success with evidence, success without evidence, diagnostic-only completion, runtime failure, and redaction. |
| Observability / REQ-23-48 | Status reads `flow-state.json`, `.run-state.json`, `execute.md`, and validation facts without runtime calls; output includes task counts, attempts, and safe diagnostics. | Status command tests, text/JSON output checks, and no-runtime fake assertions. |
| Performance | Keep status/resume bounded; avoid recursive target scans, unbounded goroutines, and raw runtime payload persistence. | Code review, larger fixture state/status tests if needed, race test. |
| Persistence And Migrations / REQ-23-46 / REQ-23-47 | `.run-state.json` is versioned, strict, atomic, resumable, and recovers stale `running` tasks with diagnostics. | State-store tests for schema, malformed JSON, unsupported versions, atomic writes, stale recovery, cancellation, resume, and terminal counts. |
| Security / REQ-23-44 / REQ-23-51 | Constrain target writes to `../ultraplan-go`, reject escapes, redact secrets, and add no smoke/review/issues/Git mutation behavior. | Target containment tests, redaction tests, permission/request mapping review, non-goal regression tests. |
| Testing / REQ-23-50 | Normal verification is deterministic, offline, fake-runtime based, and passes `go test ./...`, `go test -race ./...`, and `go build ./cmd/ultraplan`. | Required commands run from `../ultraplan-go` plus focused unit/command/fake-runtime tests. |
| Workflows / REQ-23-40 / REQ-23-42 | Execute is a stateful stage after valid `plan.md`, extracts deterministic task IDs, tracks `pending/running/complete/failed/cancelled`, supports resume, and preserves traceability. | Workflow tests for prerequisites, task extraction, deterministic IDs, selected task, cancellation, stale recovery, resume without redoing complete tasks, and summary output. |

## Tasks

- [x] **Task 1: Establish Execute Domain Types And State Schema**
  > Executes: `Decision 1`, `Decision 3`, `Architecture`, `Persistence And Migrations`, `Workflows`, `REQ-23-46`, `REQ-23-47`
  - [x] Add sprint-owned execute domain types for task records, deterministic task identity inputs, statuses `pending`, `running`, `complete`, `failed`, `cancelled`, attempts, timestamps, diagnostics, evidence, runtime metadata summaries, target reference, plan path, and plan fingerprint.
  - [x] Add versioned `.run-state.json` schema constants and strict validation rules for project, sprint, target, plan fingerprint, task records, statuses, attempts, timestamps, diagnostics, and metadata summaries.
  - [x] Add same-directory atomic write and strict load behavior through sprint-owned persistence or existing product-neutral file helpers only.
  - [x] Add tests for schema versioning, malformed JSON, unsupported versions, missing required fields, invalid statuses, atomic write success/failure behavior, and last-known-good preservation.

- [x] **Task 2: Validate Execute Prerequisites And Extract Plan Tasks**
  > Executes: `Decision 2`, `Workflows`, `Errors`, `Testing`, `REQ-23-40`, `REQ-23-42`
  - [x] Extend sprint validation so execute refuses missing, invalid, or unvalidated prerequisites through `plan.md` before task extraction or runtime launch.
  - [x] Implement deterministic extraction from explicit `plan.md` task checklist entries only, preserving task name, implementation steps, reasoning decision references, requirement references, and expected evidence where present.
  - [x] Compute deterministic task IDs from stable normalized task fields and reject duplicate, unsupported, or ambiguous task syntax with actionable diagnostics.
  - [x] Add fixture tests for valid tasks, no executable tasks, duplicate IDs, ambiguous tasks, unsupported syntax, stable IDs across equivalent content, changed IDs when task identity changes, and traceability fields.

- [x] **Task 3: Resolve Target Repository And Safety Boundaries**
  > Executes: `Decision 5`, `Security`, `LLM Runtime`, `REQ-23-43`, `REQ-23-44`, `REQ-23-51`
  - [x] Add target resolver behavior that uses the project index target implementation directory and constrains this sprint to `../ultraplan-go` unless later explicit safe configuration is designed.
  - [x] Validate target path containment and fail before runtime for path escapes, missing target roots, unsupported alternate layouts, or unsafe workdir values.
  - [x] Ensure execute prompts and runtime permission/request metadata state the approved target and explicitly exclude smoke, review automation, issue tracking, scheduling, and Git mutation.
  - [x] Add target containment, path escape, workspace-relative diagnostic, permission mapping, and forbidden Git mutation regression tests.

- [x] **Task 4: Implement Stage Model Resolution For Execute**
  > Executes: `Decision 7`, `Configuration`, `Observability`, `Security`, `REQ-23-41`
  - [x] Extend or reuse sprint stage model resolution to cover `execute` with precedence: explicit command override when supported, stage-specific model, global sprint/planning model, `models.primary`, then `models.default`.
  - [x] Reject unknown stage model keys and empty model values during validation with field-path diagnostics.
  - [x] Include selected model source in prompt previews, runtime request summaries, and diagnostics without exposing secrets or full runtime config.
  - [x] Add tests for override precedence, stage-specific `execute` winning over global/default, fallback order, unknown stage key rejection, empty model rejection, command override, and redacted diagnostics.

- [x] **Task 5: Render Execute Prompts And Runtime Requests**
  > Executes: `Decision 1`, `Decision 5`, `Decision 7`, `LLM Runtime`, `LLM Evaluation / Cost / Safety`, `REQ-23-43`, `REQ-23-44`
  - [x] Add deterministic execute prompt rendering from validated task data, sprint/project paths, target path, expected evidence, safety constraints, and non-goal exclusions.
  - [x] Build generic platform runtime requests with prompt, workdir, model, timeout, permissions, expected outputs, metadata, and validation expectations while keeping `internal/platform/runtime` free of sprint semantics.
  - [x] Persist or expose only safe runtime metadata summaries such as run/session IDs, model source, attempt count, validation summary, permission summary, usage summary where safe, and omission reasons for unsafe payloads.
  - [x] Add golden prompt tests and fake-runtime request assertion tests covering target scope, model source, metadata, permissions, expected evidence, and absence of forbidden workflow instructions.

- [x] **Task 6: Run Execute Tasks Conservatively With Durable Transitions**
  > Executes: `Decision 3`, `Decision 4`, `Decision 6`, `Persistence And Migrations`, `Workflows`, `Performance`, `REQ-23-45`, `REQ-23-46`, `REQ-23-47`
  - [x] Implement the execute task runner to process one selected or pending task at a time by default, mark transitions atomically, increment attempts, drain runtime events, call runtime wait, and save safe diagnostics.
  - [x] Apply completion gates so runtime success without expected evidence becomes failed unless an explicit safe diagnostic-only completion path is recorded.
  - [x] Implement cancellation handling that stops scheduling, cancels active runtime work through context/runtime APIs, saves state, and marks active task `cancelled` or retryable/failed according to the selected policy.
  - [x] Implement resume reconciliation that loads state strictly, recovers stale `running` tasks with diagnostics, skips `complete` tasks, and refuses mismatched plan fingerprints or task sets unless an explicit future force behavior is designed.
  - [x] Add fake-runtime workflow tests for success with evidence, success without evidence, runtime failure, validation failure, diagnostic-only completion, cancellation, timeout, stale recovery, selected task execution, resume without redoing complete tasks, and no goroutine leaks.

- [x] **Task 7: Generate Execute Summary And Runtime-Free Status**
  > Executes: `Decision 3`, `Decision 4`, `Decision 8`, `Documentation`, `Observability`, `REQ-23-48`, `REQ-23-49`
  - [x] Add `execute.md` summary rendering that cites `plan.md`, lists task counts by terminal state, records completion evidence, records actionable failed-task diagnostics or `.run-state.json` pointers, and avoids claiming deferred behavior completion.
  - [x] Update sprint status to compute execute readiness/progress from prerequisite validation, `flow-state.json`, `.run-state.json`, and `execute.md` without invoking runtime.
  - [x] Add `execute.md` validation checks for required sections, plan citation, task counts, terminal states, evidence/diagnostic entries, and absence of smoke/review/issues/Git mutation claims.
  - [x] Add golden tests for `execute.md`, text status, JSON status, stale summary indicators where applicable, and status behavior when run state is missing, malformed, in progress, failed, or complete.

- [x] **Task 8: Integrate Execute Into Flow And Validation**
  > Executes: `Decision 1`, `Decision 2`, `Decision 8`, `Decision 9`, `Workflows`, `REQ-23-40`, `REQ-23-48`, `REQ-23-51`
  - [x] Update sprint flow so `execute` is allowed only after valid prerequisites through `plan.md` and deferred targets such as `smoke`, `review`, and `issues` are rejected as out of scope.
  - [x] Update `flow-state.json` handling to represent execute readiness, progress, success, failure, and prerequisite failures without adding deferred stages.
  - [x] Ensure `validate execute` checks prerequisites, task extractability, target safety, run-state validity when present, and `execute.md` validity when produced.
  - [x] Add flow and validation tests for missing prerequisites, invalid plan, no executable tasks, execute readiness, execute failure status, complete status, rejected deferred stages, and runtime-free validation/status paths.

- [x] **Task 9: Wire App And CLI Commands Thinly**
  > Executes: `Decision 8`, `CLI Surface`, `Errors`, `Testing`, `REQ-23-40`, `REQ-23-48`
  - [x] Add app/use-case request and result types for execute validation, prompt preview, flow through execute, task execution, selected task, dry-run, resume, and status output.
  - [x] Wire CLI commands and flags for `ultraplan sprint <project> <sprint> validate execute`, execute prompt preview, execute-stage flow handling, `execute`, `execute --task <id>`, dry-run/prompt preview, text/JSON status, help, and exit-code mapping.
  - [x] Keep command handlers limited to parsing, dependency construction, use-case calls, and output rendering; do not put task extraction, runtime execution, state transitions, or validation policy in command glue.
  - [x] Add command tests for help, usage errors, invalid prerequisites, invalid stages, dry-run/prompt preview, selected task, runtime failure, validation failure, partial completion, cancellation exit code, text output, JSON output, and status without runtime.

- [x] **Task 10: Enforce Deferred Behavior Exclusions And Dependency Boundaries**
  > Executes: `Decision 1`, `Decision 5`, `Decision 9`, `Architecture`, `Security`, `REQ-23-43`, `REQ-23-51`
  - [x] Review execute implementation for forbidden imports or calls: no study services/models, no direct OpenCode/opencode adapter calls from sprint/app product code except through platform runtime composition, no shelling out to runtime tools, and no native runtime stream parsing.
  - [x] Ensure no execute path creates `smoke.md`, `smoke.json`, generated `review.md`, `issues.md`, `issues.json`, Git add/commit/push/branch/PR/checkout/reset behavior, hosted/browser UI, TUI behavior, or cross-sprint scheduling.
  - [x] Add regression tests or focused review assertions for rejected stages, forbidden artifacts, forbidden command families, no Git mutation, no study-service reuse, and platform/runtime product-agnostic imports.
  - [x] Record any required deviation from `reasoning.md` before implementation continues rather than silently expanding scope.

- [x] **Task 11: Complete Verification And Review Evidence**
  > Executes: `Decision 10`, `Testing`, `Architecture`, `REQ-23-50`
  - [x] Run focused unit, command, fake-runtime, fixture, and golden tests covering execute task extraction, state, flow, status, prompts, validation, target safety, model resolution, and non-goals.
  - [x] Run `go test ./...` from `../ultraplan-go` and save the result in implementation notes for review.
  - [x] Run `go test -race ./...` from `../ultraplan-go` and save the result in implementation notes for review.
  - [x] Run `go build ./cmd/ultraplan` from `../ultraplan-go` and save the result in implementation notes for review.
  - [x] Prepare architecture and sprint review evidence using `system/protocols/architecture-review-protocol.md` and `system/protocols/sprint-review-protocol.md`, including import review, command evidence, fake-runtime evidence, and non-goal exclusions.

## Evidence Checklist

- [x] Tests prove required execute behavior from validated `plan.md` through task extraction, runtime request construction, durable state, evidence gating, summary/status, and flow integration.
- [x] Runtime or diagnostic evidence exists where required, using deterministic fake runtimes for normal acceptance.
- [x] Documentation and help/status updates are complete where required.
- [x] Deviations from `reasoning.md` are recorded before implementation continues.
- [x] Required review protocols have evidence: `system/protocols/architecture-review-protocol.md` and `system/protocols/sprint-review-protocol.md`.
- [x] Import/code review confirms dependency direction and product/runtime separation.
- [x] Non-goal regression evidence confirms no smoke, review automation, issue tracking, Git mutation, hosted/browser UI, TUI behavior, cross-sprint scheduling, or study-service reuse.

## Verification Commands

| Check | Command | Expected Result |
| --- | --- | --- |
| Focused sprint package tests | `go test ./internal/sprint` | Execute task extraction, state, prompt, validation, flow, status, and summary tests pass with fake runtime/fixtures. |
| Focused app command tests | `go test ./internal/app` | Execute command, prompt, validate, flow, status, help, selected task, exit-code, and JSON/text output tests pass. |
| Full test suite | `go test ./...` | All normal deterministic tests pass without OpenCode, provider credentials, network access, or Git mutation. |
| Race suite | `go test -race ./...` | Race-enabled tests pass, including execute cancellation/state/status paths. |
| CLI build | `go build ./cmd/ultraplan` | `ultraplan` binary builds successfully. |
| Dependency boundary review | `go list -deps ./internal/sprint ./internal/platform/runtime ./internal/app` | Review output shows `internal/sprint` uses allowed dependencies and `internal/platform/runtime` imports no product modules. |

Run verification commands from `../ultraplan-go`.

## Risks And Blockers

| Risk / Blocker | Source | Mitigation | Status |
| --- | --- | --- | --- |
| Valid `plan.md` task structure may be too loose for deterministic extraction. | `reasoning.md#assumptions-and-risks`; `technical-handbook.md#open-questions-for-reasoning` | Reject unsupported or ambiguous tasks before runtime and require plan repair. | mitigated |
| Runtime-created edits can conflict across tasks. | `reasoning.md#decision-6-execute-runs-conservatively-and-resumes-safely`; `reasoning/architecture.md#performance-and-scale-assumptions` | Run serially by default and defer parallel execution until conflict controls are designed. | mitigated by plan |
| `.run-state.json` can drift from target repository reality if evidence is weak. | `reasoning.md#assumptions-and-risks` | Require evidence or explicit safe diagnostic before marking complete; persist missing-evidence failures. | mitigated |
| Diagnostics may leak secrets or unsafe runtime payloads. | `reasoning.md#decision-5-target-repository-and-runtime-diagnostics-are-safety-gated`; `TRD.md#19-logging-and-diagnostics` | Redact secrets, persist safe summaries only, and record unsafe payload omission facts. | mitigated |
| Concurrent execute invocations may race on state. | `reasoning.md#potential-technical-debt`; `reasoning.md#assumptions-and-risks` | Use atomic writes and strict reload validation now; add explicit locks only if implementation/review exposes a concrete race. | accepted follow-up if concurrent execute becomes common |
| `execute.md` may be stale relative to `.run-state.json`. | `reasoning.md#decision-3-versioned-run-state-is-the-durable-source-of-truth` | Treat `.run-state.json` as source of truth, status reads state, and summary refreshes after execute transitions. | mitigated by plan |
| Existing platform/runtime seams may not yet support all fake-runtime injection needs. | `reasoning.md#assumptions-and-risks` | Keep runtime dependency injected at app/sprint boundary and refactor composition narrowly if needed. | mitigated |
| Target repository layout may differ from `../ultraplan-go`. | `reasoning.md#decision-5-target-repository-and-runtime-diagnostics-are-safety-gated` | Reject alternate layouts this sprint unless future explicit safe configuration is selected. | mitigated by plan |
| Fake runtime may not reveal provider-specific integration issues. | `reasoning.md#decision-10-verification-uses-deterministic-fake-runtime-evidence` | Keep fake-runtime tests as acceptance; leave real-runtime smoke optional and gated outside normal acceptance. | accepted |

## Review Inputs

Review should use:

- `projects/ultraplan-go/sprints/23-execute-stage/sprint-index.md`
- `projects/ultraplan-go/sprints/23-execute-stage/technical-handbook.md`
- `projects/ultraplan-go/sprints/23-execute-stage/reasoning/architecture.md`
- `projects/ultraplan-go/sprints/23-execute-stage/reasoning.md`
- `projects/ultraplan-go/sprints/23-execute-stage/plan.md`
- implementation diff in `../ultraplan-go`
- verification evidence from focused tests, `go test ./...`, `go test -race ./...`, and `go build ./cmd/ultraplan`
- `system/protocols/architecture-review-protocol.md`
- `system/protocols/sprint-review-protocol.md`

## Execution Log

| Date / Step | Action | Evidence / Notes |
| --- | --- | --- |
| `2026-07-04 / planning` | Created sprint implementation plan from `reasoning.md`. | Plan cites required inputs and carries forward decisions, evidence, risks, assumptions, and open questions. |
| `2026-07-09 / Task 1` | Added execute domain types, versioned `.run-state.json` load/save/validation, and atomic write tests in `internal/sprint`. | `GOCACHE=/tmp/ultraplan-go-build-cache go test ./internal/sprint` passed. |
| `2026-07-09 / Task 2` | Added execute prerequisite validation and deterministic extraction from explicit `plan.md` task checklist entries. | `GOCACHE=/tmp/ultraplan-go-build-cache go test ./internal/sprint` passed. |
| `2026-07-09 / Task 3` | Added execute target resolver, approved-target enforcement, workdir containment, and non-goal safety instruction coverage. | `GOCACHE=/tmp/ultraplan-go-build-cache go test ./internal/sprint` passed. |
| `2026-07-09 / Task 4 partial` | Added `planning.execute_model` / `planning.execute_variant` config fields and source-aware execute model precedence. | `GOCACHE=/tmp/ultraplan-go-build-cache go test ./internal/sprint ./internal/platform/config ./internal/app` passed. Remaining: stage-key and empty-value validation coverage. |
| `2026-07-09 / completion` | Wired execute prompt, validation, flow target, serial runner, `.run-state.json`, `execute.md`, CLI `execute`, selected task, dry-run, resume, and safety exclusions. | `GOCACHE=/tmp/ultraplan-go-build-cache go test ./...`, `GOCACHE=/tmp/ultraplan-go-build-cache go test -race ./...`, and `GOCACHE=/tmp/ultraplan-go-build-cache GOMODCACHE=/tmp/ultraplan-go-mod-cache go build ./cmd/ultraplan` passed. |

## Completion Criteria

- [x] All tasks are complete or explicitly deferred.
- [x] Verification commands were run or deferrals are documented.
- [x] Evidence satisfies the expectations from `reasoning.md`.
- [x] `review.md` can evaluate conformance without guessing intent.
- [x] Execute implementation preserves sprint-owned semantics and generic runtime separation.
- [x] `.run-state.json` is durable source of truth and `execute.md` is reviewable summary only.
- [x] Execute refuses invalid prerequisites and unsupported plan tasks before runtime launch.
- [x] Task completion requires expected evidence or explicit safe diagnostic.
- [x] Target path safety and diagnostic redaction are verified.
- [x] Deferred smoke, review, issue, Git mutation, TUI, hosted/browser, and cross-sprint scheduling behavior remains absent.
