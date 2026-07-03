# Sprint Plan: Reason Stage

> Project: `ultraplan-go`
> Sprint: `20-reason-stage`
> Source: `projects/ultraplan-go/sprints/20-reason-stage/reasoning.md`
> **Inputs Used:** `projects/ultraplan-go/project-index.md`, `projects/ultraplan-go/sprints/20-reason-stage/requirements.md`, `projects/ultraplan-go/docs/PRD.md`, `projects/ultraplan-go/docs/TRD.md`, `projects/ultraplan-go/docs/ARCHITECTURE.md`, `projects/ultraplan-go/sprints/20-reason-stage/sprint-index.md`, `projects/ultraplan-go/sprints/20-reason-stage/technical-handbook.md`, `projects/ultraplan-go/sprints/20-reason-stage/reasoning.md`, `templates/sprint-plan.md`

This plan executes `reasoning.md`. It does not invent architecture, scope, or decisions beyond the final sprint reasoning.

Linked final study reports were not reopened for this plan because `technical-handbook.md` already distilled the selected reports with concrete source references, and no plan-level decision required deeper evidence.

## Reasoning Source

- **Sprint Reasoning:** `projects/ultraplan-go/sprints/20-reason-stage/reasoning.md`
- **Sprint Index:** `projects/ultraplan-go/sprints/20-reason-stage/sprint-index.md`
- **Technical Handbook:** `projects/ultraplan-go/sprints/20-reason-stage/technical-handbook.md`
- **Area Reasoning:** none present under `projects/ultraplan-go/sprints/20-reason-stage/reasoning/*.md` when this plan was written; the selected Architecture area output remains a required generated/validated artifact during implementation flow.
- **Key Boundary:** Sprint 20 stops at `reasoning`; it must not generate, validate, or complete `plan.md` as product behavior.

## Sprint Status

- **Status:** complete
- **Owner:** implementation agent
- **Start Date:** 2026-06-15
- **Completion Date:** 2026-06-15

## Decisions To Execute

| Decision | Source Section | Execution Implication |
| --- | --- | --- |
| Reason-stage ownership stays in `internal/sprint` | `reasoning.md#decision-1-reason-stage-ownership-stays-in-internalsprint` | Add selected-template detection, area/final reasoning validation, prompt rendering, flow behavior, service use cases, and filesystem loading inside `internal/sprint`; do not add global planning, workflow, validation, prompts, reports, or catalog packages. |
| Selected reasoning templates become a deterministic manifest | `reasoning.md#decision-2-selected-reasoning-templates-become-a-deterministic-manifest` | Derive selected template entries from validated `sprint-index.md` and `project-index.md`; include template name/path, output path, selected-context summary, ordering, and diagnostics. |
| Reasoning validators enforce artifact contracts without brittle prose matching | `reasoning.md#decision-3-reasoning-validators-enforce-artifact-contracts-without-brittle-prose-matching` | Implement area and final validators that check existence, non-empty content, placeholders, required sections, selected-context discipline, and selected-area prerequisites without exact prose matching. |
| Prompt rendering is deterministic, selected-context-only, and non-mutating except expected output | `reasoning.md#decision-4-prompt-rendering-is-deterministic-selected-context-only-and-non-mutating-except-expected-output` | Render `area-reasoning` and `reasoning` prompt previews with required paths, selected context, required sections, output paths, workspace-relative paths, and no-mutation rules. |
| Runtime-backed flow uses only the generic runtime boundary and validates after generation | `reasoning.md#decision-5-runtime-backed-flow-uses-only-the-generic-runtime-boundary-and-validates-after-generation` | Extend non-dry-run flow to invoke only the generic platform runtime seam, then require expected artifact existence, stage validation, and atomic state update before completion. |
| Flow-state remains strict, atomic, and limited to Planning Phase 2 stages | `reasoning.md#decision-6-flow-state-remains-strict-atomic-and-limited-to-planning-phase-2-stages` | Support only `requirements`, `sprint-index`, `technical-handbook`, `area-reasoning`, `reasoning`, and `plan` awareness; reject execution/smoke/review/issues stages and preserve atomic writes. |
| CLI wiring is thin, scriptable, and exit-code disciplined | `reasoning.md#decision-7-cli-wiring-is-thin-scriptable-and-exit-code-disciplined` | Keep `internal/app/sprint_commands.go` limited to parsing, delegation, output rendering, help text, and exit-code mapping. |

## Requirements / Contracts To Satisfy

| Contract / Requirement ID | Required Behavior | Evidence Planned |
| --- | --- | --- |
| Sprint goal; `requirements.md` lines 7-10 | Implement reason-stage validation, prompt preview, flow generation, and `flow-state.json` updates through `reasoning` only. | Sprint flow tests, command tests, and review evidence. |
| Area/final validate; `requirements.md` lines 31-36, 47-50 | `validate area-reasoning` and `validate reasoning` run without runtime or mutation and reject missing, empty, placeholder, structurally incomplete, or unselected-context artifacts. | `internal/sprint/reasoning_test.go` and command no-runtime tests. |
| Prompt previews; `requirements.md` lines 35-36, 42-45, 69 | `prompt area-reasoning` and `prompt reasoning` render deterministic selected-context prompts without writing normal output artifacts or invoking runtime. | Prompt tests asserting required variables, paths, selected context, no-mutation rules, no ANSI, and no runtime invocation. |
| Dry-run flow; `requirements.md` lines 37, 70 | `flow --to reasoning --dry-run` reports stages, selected templates, expected outputs, and does not create artifacts or mark reasoning complete. | Flow dry-run tests against artifacts and `flow-state.json`. |
| Runtime-backed flow gates; `requirements.md` lines 38-41, 65, 70-72 | Non-dry-run flow validates prerequisites, uses only generic runtime, validates generated artifacts, handles runtime/missing/invalid outputs, and records safe failures. | Fake-runtime flow tests for success, missing artifact, invalid artifact, runtime failure, and post-generation validation failure. |
| Selected templates; `requirements.md` lines 39, 46-49, 66-68 | Area reasoning is skipped only when no templates are selected; selected Architecture requires `reasoning/architecture.md`; selected template paths must be cataloged, readable, and workspace-safe. | Manifest tests for no templates, Architecture selection, missing/unreadable/unselected/escaping paths, deterministic ordering, and selected area requirement. |
| Flow-state scope; `requirements.md` lines 51-56, 64 | Successful flow through `reasoning` marks prerequisites and reasoning complete, does not mark `plan` complete, rejects unsupported future stages, and preserves last valid state on write failure. | State and flow tests for strict stage set, no-plan completion, unsupported stages, malformed state, atomic write success, and write failure. |
| Architecture contract; `requirements.md` lines 57-60, 92-99, 102 | Keep behavior in `internal/sprint`, CLI thin, platform product-agnostic, no `internal/study` import, no global planning abstractions. | Import/package review, architecture review protocol, and `go test ./...`. |
| CLI surface; `requirements.md` lines 60-65, 71, 125, 129, 135 | Output is calm, scriptable, deterministic, no ANSI by default, stdout/stderr separated, with exit codes `0`, `2`, `5`, and runtime failure mapping. | `internal/app/sprint_commands_test.go` help, output, no-runtime, fake-runtime, and exit-code cases. |
| Security/configuration; `requirements.md` lines 40, 45, 56-58, 94-101 | Use workspace-safe paths, redacted diagnostics, no direct shell/Git/OpenCode/provider/agent adapter calls from `internal/sprint`, and no secrets in prompts. | Path escape tests, prompt diagnostics tests, import review, and code review. |
| Verification; `requirements.md` lines 72-74, 136 | Offline tests, race tests, and CLI build pass from `/home/antonioborgerees/coding/ultraplan-go`. | `go test ./...`, `go test -race ./...`, `go build ./cmd/ultraplan`. |

## Tasks

- [x] **Task 1: Baseline Existing Sprint Surface**
  > Executes: Decisions 1, 6, 7; Architecture and CLI Surface contracts
  - [x] Inspect current `/home/antonioborgerees/coding/ultraplan-go/internal/sprint` files for existing stage constants, validators, flow services, prompt rendering, filesystem store, and flow-state write helpers.
  - [x] Inspect current `/home/antonioborgerees/coding/ultraplan-go/internal/app/sprint_commands.go` and related tests for command parsing, help, exit-code mapping, and output conventions.
  - [x] Identify reusable generic infrastructure already present for workspace-safe paths, atomic writes, runtime seams, diagnostics, and exit codes without moving product rules out of `internal/sprint`.
  - [x] Record any implementation shape mismatch as a narrow integration note; do not reopen reasoning decisions unless a requirement conflict is found.

- [x] **Task 2: Add Reasoning Domain And Selected-Template Manifest**
  > Executes: Decisions 1, 2; selected-template and security requirements
  - [x] Add `internal/sprint/reasoning.go` or equivalent sprint-owned file for selected reasoning template detection, manifest entries, final reasoning input assembly, and validation request/result types.
  - [x] Load selected reasoning templates only from validated `sprint-index.md` entries that are present in `project-index.md`.
  - [x] Validate selected template paths are workspace-safe, readable, cataloged, and not path escapes unless explicitly cataloged as external.
  - [x] Map each selected template to a deterministic required output path under the sprint `reasoning/` directory, with Architecture mapping to `reasoning/architecture.md` for this sprint.
  - [x] Sort manifest entries deterministically by template name or output path and reuse the manifest for validate, prompt, dry-run, runtime flow, and final reasoning prerequisites.
  - [x] Include selected contracts, selected evidence report names/paths, required review protocols, and excluded-context summary needed by prompt rendering without reselecting context.

- [x] **Task 3: Extend Filesystem Store For Reasoning Inputs**
  > Executes: Decisions 1, 2, 3, 4; workspace and selected-context requirements
  - [x] Add store methods in `internal/sprint/store_fs.go` or existing store files to read requirements, sprint index, technical handbook, selected template files, area reasoning files, and final `reasoning.md` using workspace-safe paths.
  - [x] Ensure missing, unreadable, empty, or escaping files return classified diagnostics with artifact path and stage context.
  - [x] Keep project catalog loading through `internal/project`; do not duplicate project-index parsing as a new generic catalog layer.
  - [x] Preserve workspace-relative paths in public prompt/diagnostic content unless absolute local diagnostics are required.

- [x] **Task 4: Implement Area And Final Reasoning Validators**
  > Executes: Decision 3; validation requirements
  - [x] Implement `ValidateAreaReasoning` behavior that handles no-template skip, selected-template required files, non-empty content, placeholder rejection, area decisions, trade-offs, and selected-context references.
  - [x] Implement `ValidateReasoning` behavior that requires final `reasoning.md`, final decisions, expected evidence, assumptions and risks, selected requirement/contract/evidence mapping, no placeholders, and required selected-area references.
  - [x] Do not require `plan.md`, task checklists, smoke evidence, review evidence, issue artifacts, Git state, or implementation output for reasoning validation.
  - [x] Keep validation structural and enforceable; avoid exact generated prose or brittle heading-count assertions beyond required contract sections.
  - [x] Return deterministic, actionable diagnostics that identify stage, file, missing section or invalid reference, safe cause, and suggested repair target.

- [x] **Task 5: Render Area And Final Reasoning Prompts**
  > Executes: Decision 4; prompt preview requirements
  - [x] Add deterministic `area-reasoning` prompt rendering with project slug, sprint slug, sprint path, requirements path, sprint-index path, technical-handbook path, selected reasoning template path, output path, selected contracts/evidence summary, required sections, and selected-context-only rules.
  - [x] Add deterministic final `reasoning` prompt rendering with project slug, sprint slug, sprint path, requirements path, sprint-index path, technical-handbook path, required selected area reasoning paths, final output path, required sections, and no-plan/no-execution rules.
  - [x] Include explicit no-mutation instructions for `requirements.md`, `sprint-index.md`, `technical-handbook.md`, `plan.md`, project docs, prior reviews, source repositories, implementation files outside expected runtime output, workspace config, and Git state.
  - [x] Ensure prompt preview commands do not invoke runtime and do not write `reasoning/*.md` or `reasoning.md`; if an output-preview flag already exists, write only to the explicit preview path.
  - [x] Avoid embedding secrets, provider credentials, full environment dumps, unsafe raw runtime payloads, or unnecessary absolute paths.

- [x] **Task 6: Extend Flow And Flow-State Through Reasoning**
  > Executes: Decisions 5, 6; workflow, runtime, persistence, and observability requirements
  - [x] Extend stage parsing and flow targets to support `area-reasoning` and `reasoning` while still rejecting `implementation`, `execute`, `smoke`, `review`, `issues`, `.run-state.json`, `smoke.md`, `smoke.json`, generated `review.md`, `issues.md`, and `issues.json` as current stages.
  - [x] Preflight non-dry-run flow in order: valid requirements, valid sprint index, valid technical handbook, valid selected templates, then selected area/final prerequisites as appropriate.
  - [x] Implement dry-run output that reports stages, selected templates, expected area outputs, final output, and planned state transitions without runtime invocation or artifact/state completion.
  - [x] For non-dry-run selected area reasoning, invoke only the generic platform runtime boundary, validate the expected area artifact after runtime success, and atomically update `flow-state.json` only after validation passes.
  - [x] For non-dry-run final reasoning, require selected area artifacts when templates are selected, invoke only the generic platform runtime boundary, validate final `reasoning.md`, and atomically update state only after validation passes.
  - [x] Mark `area-reasoning` as `skipped` only when the validated sprint index selects no reasoning templates.
  - [x] On area/final runtime or validation failure, record the failed stage with a safe error message, preserve prior valid state on write failure, return the correct non-zero exit code, and do not mark later stages complete.
  - [x] Successful flow through `reasoning` must mark prerequisites and reasoning stages complete or skipped as applicable and must not mark `plan` complete.

- [x] **Task 7: Add Sprint Service Use Cases**
  > Executes: Decisions 1, 5, 7; service boundary requirements
  - [x] Add or extend `internal/sprint/service.go` request/result methods for validating `area-reasoning`, validating `reasoning`, rendering `area-reasoning` prompts, rendering `reasoning` prompts, dry-run flow, and runtime-backed flow through `reasoning`.
  - [x] Keep service methods responsible for selected-context loading, validator calls, prompt builders, runtime seam calls, and flow-state transitions.
  - [x] Keep runtime-backed service paths context-aware; keep deterministic validate/prompt paths runtime-free.
  - [x] Use explicit small seams for runtime, store, clock/time, and output summaries only where tests or external volatility require them.

- [x] **Task 8: Wire CLI Commands Thinly**
  > Executes: Decision 7; CLI Surface and Documentation contracts
  - [x] Extend `internal/app/sprint_commands.go` so `ultraplan sprint <project> <sprint> validate area-reasoning` delegates to sprint validation and maps valid output to exit code `0`.
  - [x] Extend `validate reasoning` with final reasoning validation and required selected-area prerequisite checks.
  - [x] Extend `prompt area-reasoning` and `prompt reasoning` to render previews to stdout or existing explicit preview output path without runtime invocation.
  - [x] Extend `flow --to reasoning` and `flow --to reasoning --dry-run` to call sprint flow services and map results to stdout/stderr and exit codes.
  - [x] Update help output to show `area-reasoning` and `reasoning` support and reject unsupported future stages with usage exit code `2`.
  - [x] Preserve calm deterministic text, no ANSI by default, stdout for normal output, stderr for diagnostics, and validation exit code `5` for invalid artifacts or selected context.

- [x] **Task 9: Add Sprint Reasoning And Flow Tests**
  > Executes: Decisions 2, 3, 4, 5, 6; Testing contract
  - [x] Create or extend `/home/antonioborgerees/coding/ultraplan-go/internal/sprint/reasoning_test.go` for no templates, one Architecture template, missing/unreadable/unselected/escaping template paths, deterministic ordering, selected output paths, and selected-context summaries.
  - [x] Test area reasoning validation for valid content, missing file, empty file, placeholders, missing area decisions, missing trade-offs, unselected-context references, and path escapes.
  - [x] Test final reasoning validation for valid content, missing file, empty file, placeholders, missing final decisions, missing expected evidence, missing assumptions/risks, missing required selected-area references, and selected-area success.
  - [x] Test prompt rendering without brittle prose by asserting required variables, selected template manifest, selected context summary, output path, required sections, selected-only rules, no-mutation rules, workspace-relative paths, and no runtime invocation.
  - [x] Extend `/home/antonioborgerees/coding/ultraplan-go/internal/sprint/flow_test.go` for dry-run no mutation, fake-runtime successful area generation, fake-runtime successful final generation, runtime success with missing artifact, runtime success with invalid artifact, runtime failure, validation failure after generation, no-template area skip, selected-template failure before runtime, no plan completion, unsupported stage rejection, and flow-state write failure.

- [x] **Task 10: Add Command Tests**
  > Executes: Decision 7; CLI Surface, Errors, Documentation, Observability contracts
  - [x] Extend `/home/antonioborgerees/coding/ultraplan-go/internal/app/sprint_commands_test.go` for `validate area-reasoning`, `validate reasoning`, `prompt area-reasoning`, `prompt reasoning`, `flow --to reasoning --dry-run`, and `flow --to reasoning` with fake runtime.
  - [x] Assert validate and prompt command paths do not call runtime.
  - [x] Assert non-dry-run flow uses only the generic fake runtime seam.
  - [x] Assert help output includes supported reasoning stages and unsupported future stages fail with exit code `2`.
  - [x] Assert validation failures return exit code `5`, runtime failures map to the established runtime failure code, stdout/stderr separation is preserved, and no ANSI escape sequences appear by default.

- [x] **Task 11: Run Verification And Record Review Evidence**
  > Executes: Expected Evidence and Review Expectations
  - [x] Run `go test ./...` from `/home/antonioborgerees/coding/ultraplan-go`.
  - [x] Run `go test -race ./...` from `/home/antonioborgerees/coding/ultraplan-go`.
  - [x] Run `go build ./cmd/ultraplan` from `/home/antonioborgerees/coding/ultraplan-go`.
  - [x] Review imports and package additions to confirm `internal/sprint` does not import `internal/study`, platform packages do not import product modules, and no global planning/workflow/validation/prompt package was added.
  - [x] Review prompt and diagnostic output for workspace-relative paths, selected-context-only content, redaction, no mutation permission beyond expected output, and no direct shell/Git/OpenCode/provider calls from `internal/sprint`.
  - [x] Record implementation evidence, deviations, verification results, and Architecture Review/Sprint Review conclusions in `projects/ultraplan-go/sprints/20-reason-stage/review.md` after implementation.

## Evidence Checklist

- [x] Selected-template manifest tests prove catalog subset validation, workspace path safety, deterministic ordering, no-template skip, and selected Architecture output requirements.
- [x] Area reasoning validator tests prove required artifact structure, placeholders, selected-context constraints, missing/unreadable paths, and path escape diagnostics.
- [x] Final reasoning validator tests prove final decisions, expected evidence, assumptions and risks, selected-area prerequisites, placeholders, and no plan/task/smoke/review/issue/Git dependency.
- [x] Prompt tests prove deterministic previews with required variables, selected context, workspace-relative paths, output paths, required sections, selected-only rules, no-mutation rules, and no runtime invocation.
- [x] Flow tests prove dry-run no mutation, fake-runtime success, missing/invalid artifact failure, runtime failure, validation failure after runtime, no-template area skip, selected-template preflight failure, strict stage set, no plan completion, and atomic state write failure behavior.
- [x] Command tests prove help, stage parsing, stdout/stderr separation, no ANSI by default, exit codes, validate/prompt no-runtime paths, and fake-runtime flow invocation.
- [x] Architecture review confirms package ownership, import direction, product/platform separation, no study semantic reuse, and no global planning/workflow/validation/prompt package.
- [x] Security review confirms workspace-safe paths, selected template safety, redacted prompts/diagnostics, no secrets, and no direct shell/Git/OpenCode/provider/agent adapter calls from `internal/sprint`.
- [x] Verification commands pass or documented deferrals are reviewed before sprint completion.
- [x] Required review protocols have evidence in `review.md` after implementation.

## Verification Commands

| Check | Command | Expected Result |
| --- | --- | --- |
| Offline tests | `go test ./...` | All packages pass without network, provider credentials, real OpenCode, Git mutation, or long sleeps. |
| Race tests | `go test -race ./...` | Race suite passes with no data races or goroutine leaks attributable to reason-stage flow. |
| CLI build | `go build ./cmd/ultraplan` | `ultraplan` command builds successfully. |
| Help review | `ultraplan sprint --help` and stage-specific help commands available in the implemented test harness | Help includes validate/prompt/flow support for `area-reasoning` and `reasoning` and rejects unsupported future stages. |
| Import review | package/import inspection during code review | `internal/sprint` does not import `internal/study`; `internal/platform/*` imports no product modules; no global planning/workflow/validation/prompt package was added. |

## Risks And Blockers

| Risk / Blocker | Source | Mitigation | Status |
| --- | --- | --- | --- |
| No Architecture area reasoning artifact existed when `reasoning.md` and this plan were written. | `reasoning.md#assumptions-and-risks` | Implemented selected-area validation so final reasoning completion fails when Architecture is selected and `reasoning/architecture.md` is missing or invalid. | mitigated |
| Existing Sprint 17-19 APIs may differ from the file names or seams assumed in the plan. | `reasoning.md#assumptions-and-risks` | Adapted to existing `Service`, `FlowRequest`, `FlowState`, `FSStore`, prompt, and CLI patterns without new global packages. | closed |
| Selected-context reference validation may not catch every semantic unsupported claim. | `reasoning.md#assumptions-and-risks` | Implemented structural/path/name validation and recorded this as residual review risk. | carried forward |
| Fake runtime does not prove real OpenCode behavior. | `reasoning.md#assumptions-and-risks` | Verified through fake runtime only, per sprint scope; no real-runtime smoke was selected. | carried forward |
| Strict flow-state schema may reject older prototype or manually edited state. | `reasoning.md#assumptions-and-risks` | Existing strict schema preserved; unsupported stages remain rejected. | accepted |
| Prompt builders may become long because they must include many required variables and safety rules. | `reasoning.md#trade-off-and-debt-analysis` | Kept prompt rendering table/manifest-oriented in `internal/sprint/prompts.go` and asserted required variables in tests. | mitigated |
| Markdown validators may duplicate mechanics from other planning stages. | `reasoning.md#trade-off-and-debt-analysis` | Kept validators local; reused existing Markdown helpers only inside `internal/sprint`. | accepted |

## Review Inputs

Review should use:

- `projects/ultraplan-go/sprints/20-reason-stage/sprint-index.md`
- `projects/ultraplan-go/sprints/20-reason-stage/technical-handbook.md`
- `projects/ultraplan-go/sprints/20-reason-stage/reasoning/architecture.md` if created during implementation
- `projects/ultraplan-go/sprints/20-reason-stage/reasoning.md`
- `projects/ultraplan-go/sprints/20-reason-stage/plan.md`
- implementation diff in `/home/antonioborgerees/coding/ultraplan-go`
- verification outputs for `go test ./...`, `go test -race ./...`, and `go build ./cmd/ultraplan`
- `system/protocols/architecture-review-protocol.md`
- `system/protocols/sprint-review-protocol.md`
- `projects/ultraplan-go/sprints/20-reason-stage/review.md` after sprint execution

## Execution Log

| Date / Step | Action | Evidence / Notes |
| --- | --- | --- |
| 2026-06-15 / planning | Created Sprint 20 plan from selected requirements, sprint reasoning, sprint index, technical handbook, PRD, TRD, Architecture doc, project index, and plan template. | Planning only; no implementation code changed. |
| 2026-06-15 / implementation | Added reason-stage manifest, validators, prompt rendering, flow support, service methods, filesystem loading, and CLI wiring inside `internal/sprint` and `internal/app`. | Changed `internal/sprint/reasoning.go`, `internal/sprint/service.go`, `internal/sprint/prompts.go`, `internal/sprint/flow.go`, `internal/sprint/store_fs.go`, `internal/sprint/handbook.go`, `internal/app/sprint_commands.go`, and tests. |
| 2026-06-15 / verification | Ran `go test ./internal/sprint ./internal/app`, `go test ./...`, `go test -race ./...`, and `go build ./cmd/ultraplan`. | All commands passed. Import/security inspection found no `internal/sprint` imports of `internal/study`, no direct shell/Git/OpenCode/provider calls from `internal/sprint`, and no platform imports of product modules. |

## Completion Criteria

- [x] All tasks are complete or explicitly deferred with requirement-grounded rationale.
- [x] `validate area-reasoning`, `validate reasoning`, `prompt area-reasoning`, `prompt reasoning`, `flow --to reasoning --dry-run`, and `flow --to reasoning` satisfy Sprint 20 acceptance criteria.
- [x] Selected Architecture area reasoning is generated and validated before final reasoning completes, or `area-reasoning` is skipped only when the selected template manifest is empty.
- [x] Flow through `reasoning` does not generate, validate, or complete `plan.md` and does not create implementation, smoke, review automation, issue, or Git stages.
- [x] Verification commands were run or deferrals are documented and accepted in review.
- [x] Evidence satisfies the expectations from `reasoning.md` Decisions 1-7.
- [x] Architecture Review and Sprint Review evidence is recorded in `review.md` after implementation.
- [x] `review.md` can evaluate conformance without guessing intent.
