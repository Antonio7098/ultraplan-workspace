# Sprint Plan: Distill Stage

> Project: `ultraplan-go`
> Sprint: `19-distill-stage`
> Source: `reasoning.md`
> **Inputs Used:** `projects/ultraplan-go/sprints/19-distill-stage/requirements.md`, `projects/ultraplan-go/sprints/19-distill-stage/reasoning.md`, `projects/ultraplan-go/sprints/19-distill-stage/sprint-index.md`, `projects/ultraplan-go/sprints/19-distill-stage/technical-handbook.md`, `projects/ultraplan-go/docs/PRD.md`, `projects/ultraplan-go/docs/TRD.md`, `projects/ultraplan-go/docs/ARCHITECTURE.md`, `projects/ultraplan-go/project-index.md`, `templates/sprint-reasoning.md`, `templates/sprint-plan.md`; no area reasoning files were present under `projects/ultraplan-go/sprints/19-distill-stage/reasoning/`.

This plan executes `reasoning.md`. It does not invent architecture, scope, or decisions beyond the sprint reasoning document.

## Reasoning Source

- **Sprint Reasoning:** `projects/ultraplan-go/sprints/19-distill-stage/reasoning.md`
- **Sprint Index:** `projects/ultraplan-go/sprints/19-distill-stage/sprint-index.md`
- **Technical Handbook:** `projects/ultraplan-go/sprints/19-distill-stage/technical-handbook.md`
- **Area Reasoning:** none present; the sprint index selected Architecture reasoning, but no area-specific reasoning artifact exists for this sprint.
- **Omitted Evidence:** linked final study reports were not opened directly for this plan because `technical-handbook.md` and `reasoning.md` already distill the selected report evidence sufficiently for implementation planning.

## Sprint Status

- **Status:** complete
- **Owner:** implementation agent
- **Start Date:** 2026-06-14
- **Completion Date:** 2026-06-14

## Decisions To Execute

| Decision | Source Section | Execution Implication |
| --- | --- | --- |
| Distill behavior lives in `internal/sprint` | `reasoning.md#decision-1-distill-behavior-lives-in-internalsprint` | Implement selected-evidence loading, handbook validation, prompt rendering, service methods, flow behavior, store behavior, and state transitions in `internal/sprint`; keep `internal/app` as command wiring only. |
| Selected evidence is loaded through a strict catalog-and-selection manifest | `reasoning.md#decision-2-selected-evidence-is-loaded-through-a-strict-catalog-and-selection-manifest` | Build deterministic manifests from valid `requirements.md`, `project-index.md`, and `sprint-index.md`; reject missing, unreadable, uncataloged, unselected, duplicate, and escaping evidence paths before runtime. |
| Handbook validation is semantic enough to gate the distill stage | `reasoning.md#decision-3-handbook-validation-is-semantic-enough-to-gate-the-distill-stage` | Validate existence, non-empty content, placeholders, required sections, selected evidence traces, unselected evidence references, and implementation-decision language before marking the stage complete. |
| Prompt rendering is deterministic, previewable, selected-only, and non-mutating | `reasoning.md#decision-4-prompt-rendering-is-deterministic-previewable-selected-only-and-non-mutating` | Render `technical-handbook` prompts from validated selected context with workspace-relative paths, selected evidence manifest, required sections, no-decision rules, selected-only rules, output path, and no-mutation instructions. |
| Flow through `technical-handbook` uses runtime only after prerequisite validation and completes only after post-generation validation and atomic state update | `reasoning.md#decision-5-flow-through-technical-handbook-uses-runtime-only-after-prerequisite-validation-and-completes-only-after-post-generation-validation-and-atomic-state-update` | Implement dry-run, runtime-backed generation, post-generation handbook validation, failure handling, safe diagnostics, and atomic `flow-state.json` updates through `technical-handbook` only. |
| CLI wiring is thin, scriptable, and strict about supported stages | `reasoning.md#decision-6-cli-wiring-is-thin-scriptable-and-strict-about-supported-stages` | Add `validate technical-handbook`, `prompt technical-handbook`, and `flow --to technical-handbook` wiring and help while preserving stdout/stderr separation, exit codes, no ANSI output, and unsupported-stage rejection. |
| Verification is offline, fake-first, and reviewable | `reasoning.md#decision-7-verification-is-offline-fake-first-and-reviewable` | Add targeted unit, fixture, fake-runtime, command, race, build, import-boundary, and review evidence without requiring real OpenCode, provider credentials, network, Git mutation, or smoke harness execution. |

## Requirements / Contracts To Satisfy

| Contract / Requirement ID | Required Behavior | Evidence Planned |
| --- | --- | --- |
| Architecture, `REQ-AC-25`, `REQ-AC-26`, `REQ-AC-27`, `REQ-AC-28` | `internal/sprint` owns distill rules; `internal/app` is thin; `internal/platform/*` remains product-agnostic; no global planning/workflow/validation/prompt package is introduced. | Import review; command tests proving handlers delegate to sprint services; code review of package additions. |
| Security, `REQ-AC-05`, `REQ-AC-11`, `REQ-AC-12`, `REQ-AC-15`, `REQ-AC-16`, `REQ-AC-17` | Selected evidence reads are selected-only, cataloged, readable, workspace-safe, and secret-safe; sprint code does not invoke OpenCode, agentwrap adapters, shell, provider APIs, or Git directly. | `handbook_test.go` selected-evidence fixtures; flow tests proving pre-runtime failure; import review for no direct shell/OpenCode/Git calls. |
| Errors, Observability, CLI Surface, `REQ-AC-29`, `REQ-AC-30`, `REQ-AC-31` | Usage, validation, selected-evidence, runtime, and state failures map to deterministic exit codes and safe diagnostics; stdout/stderr remain separated and no ANSI is emitted by default. | `sprint_commands_test.go` for exit codes `0`, `2`, `5`, runtime failure mapping, stderr/stdout, deterministic output, no ANSI, and redaction. |
| Documentation, Configuration, Security, `REQ-AC-02`, `REQ-AC-08`, `REQ-AC-09`, `REQ-AC-10`, `REQ-AC-35` | Prompt preview is deterministic, runtime-free, selected-only, workspace-relative where possible, and forbids mutation outside the handbook output path. | Prompt rendering tests for required variables, selected manifest, output path, required sections, no-decision instructions, selected-only instructions, no-mutation instructions, and no writes. |
| Workflows, LLM Runtime, Persistence And Migrations, `REQ-AC-03`, `REQ-AC-04`, `REQ-AC-06`, `REQ-AC-07`, `REQ-AC-18`, `REQ-AC-19`, `REQ-AC-20`, `REQ-AC-21`, `REQ-AC-22`, `REQ-AC-23`, `REQ-AC-24`, `REQ-AC-36` | Flow validates prerequisites before runtime, uses only the generic runtime boundary, validates generated handbook after runtime, writes state atomically, records failures safely, and supports only Planning Phase 2 stages. | `flow_test.go` for dry-run, fake-runtime success/failure, missing artifact, invalid artifact, validation failure, selected-evidence failure, state write failure, next-stage readiness/skipping, and unsupported-stage exclusion. |
| Testing, `REQ-AC-32`, `REQ-AC-33`, `REQ-AC-34`, `REQ-AC-35`, `REQ-AC-36`, `REQ-AC-37`, `REQ-AC-38`, `REQ-AC-39`, `REQ-AC-40` | Default verification is deterministic, offline, fake-first, and comprehensive enough to review the sprint. | Targeted test files, `go test ./...`, `go test -race ./...`, `go build ./cmd/ultraplan`, and review evidence in `review.md`. |

## Tasks

- [x] **Task 1: Establish Distill-Stage Sprint Domain**
  > Executes: Decision 1, Decision 2, Architecture, Security, `REQ-AC-11`, `REQ-AC-12`, `REQ-AC-25`
  - [x] Add or extend `internal/sprint` stage constants and domain structs so `technical-handbook` is a supported Planning Phase 2 stage and unsupported future stages remain rejected.
  - [x] Add a focused `technical-handbook` input manifest model that records project slug, sprint slug, sprint root, requirements path, project-index path, sprint-index path, handbook output path, and selected evidence entries.
  - [x] Reuse existing project catalog behavior where available to compare selected evidence against `project-index.md` without importing `internal/study`.
  - [x] Keep helpers local to `internal/sprint`; do not introduce global catalog, planning, workflow, validation, reports, stages, or prompts packages.

- [x] **Task 2: Implement Selected-Evidence Loading And Store Safety**
  > Executes: Decision 2, Decision 4, Security, Configuration, `REQ-AC-11`, `REQ-AC-12`, `REQ-AC-34`
  - [x] Extend `internal/sprint/store_fs.go` or focused sprint store helpers to read selected evidence through workspace-safe paths.
  - [x] Validate that each selected evidence path is selected in `sprint-index.md`, present in `project-index.md`, readable, not duplicated ambiguously, and not escaping the workspace unless explicitly cataloged external.
  - [x] Sort selected evidence manifest entries deterministically for prompt rendering, dry-run output, diagnostics, and tests.
  - [x] Return actionable diagnostics for missing, unreadable, uncataloged, unselected, mismatched, duplicate, and escaping paths before any runtime call.

- [x] **Task 3: Implement Technical Handbook Validation**
  > Executes: Decision 3, Errors, Security, Testing, `REQ-AC-13`, `REQ-AC-14`, `REQ-AC-15`, `REQ-AC-16`, `REQ-AC-17`, `REQ-AC-33`
  - [x] Add `internal/sprint/handbook.go` with validation for missing file, empty file, template placeholders, and required sections: selected studies/reports, relevant patterns, trade-offs, anti-patterns, open questions, and evidence pointers.
  - [x] Validate that handbook evidence traces cite selected report names or selected report paths in selected studies/reports or evidence pointers sections.
  - [x] Reject uncataloged or unselected evidence report names and paths when used as evidence sources.
  - [x] Reject implementation decisions, final architecture decisions, task-plan sections, and decision-language sections that belong in `reasoning.md` or `plan.md`.
  - [x] Allow observations, warnings, anti-patterns, trade-off framing, and open questions when framed as evidence-backed guidance rather than final decisions.

- [x] **Task 4: Implement Deterministic Technical-Handbook Prompt Rendering**
  > Executes: Decision 4, Documentation, CLI Surface, Security, `REQ-AC-02`, `REQ-AC-08`, `REQ-AC-09`, `REQ-AC-10`, `REQ-AC-35`
  - [x] Extend `internal/sprint/prompts.go` to render a deterministic `technical-handbook` prompt from validated requirements, validated sprint index, selected evidence manifest, and output instructions.
  - [x] Include project slug, sprint slug, sprint path, requirements path, sprint-index path, technical-handbook output path, selected evidence manifest, required handbook sections, no-decision instructions, selected-evidence-only constraints, and no-mutation instructions.
  - [x] Use workspace-relative paths in prompt content unless an absolute path is required for local diagnostics.
  - [x] Ensure prompt preview invokes no runtime and writes no handbook artifact; if an existing explicit preview-output path feature exists, write only to that path.
  - [x] Keep prompt assertions focused on required variables and constraints rather than brittle full-prose matching unless existing test style requires goldens.

- [x] **Task 5: Implement Flow Through `technical-handbook`**
  > Executes: Decision 5, Workflows, LLM Runtime, Persistence And Migrations, `REQ-AC-03`, `REQ-AC-04`, `REQ-AC-05`, `REQ-AC-06`, `REQ-AC-07`, `REQ-AC-18`, `REQ-AC-19`, `REQ-AC-20`, `REQ-AC-21`, `REQ-AC-22`, `REQ-AC-23`, `REQ-AC-24`, `REQ-AC-36`
  - [x] Extend `internal/sprint/flow.go` so `flow --to technical-handbook --dry-run` reports stages and selected evidence without invoking runtime, writing `technical-handbook.md`, or marking the stage complete.
  - [x] For non-dry-run flow, validate `requirements.md`, validate `sprint-index.md`, and validate selected evidence before invoking runtime.
  - [x] Invoke only the generic platform runtime boundary for generation; do not call OpenCode, agentwrap adapters, shell commands, provider APIs, or Git directly from `internal/sprint`.
  - [x] After runtime success, require `technical-handbook.md` to exist and pass handbook validation before completing the stage.
  - [x] Atomically update `flow-state.json` only after post-generation validation succeeds, marking `requirements`, `sprint-index`, and `technical-handbook` complete and readying or skipping `area-reasoning` according to selected reasoning templates.
  - [x] On runtime failure, selected-evidence failure, missing artifact, invalid artifact, validation failure, or write failure, return non-zero, record safe failure state when appropriate, and never mark later reasoning or plan stages ready or complete.

- [x] **Task 6: Preserve Strict Flow-State Stage Support**
  > Executes: Decision 5, Persistence And Migrations, `REQ-AC-18`, `REQ-AC-19`, `REQ-AC-20`, `REQ-AC-21`, `REQ-AC-22`, `REQ-AC-23`, `REQ-AC-24`
  - [x] Extend `internal/sprint/state.go` to record `technical-handbook` completion, failure, and next-stage readiness while supporting only `requirements`, `sprint-index`, `technical-handbook`, `area-reasoning`, `reasoning`, and `plan`.
  - [x] Reject or fail validation for unsupported legacy/current stages such as `implementation`, `execute`, `smoke`, `review`, and `issues` instead of preserving them as supported Planning Phase 2 state.
  - [x] Preserve Sprint 17 atomic write behavior and last-valid-state preservation on write failure.
  - [x] Ensure `.run-state.json`, `smoke.md`, `smoke.json`, product-generated `review.md`, `issues.md`, and `issues.json` are not generated, validated, or marked supported by this sprint.

- [x] **Task 7: Add Service Methods And CLI Wiring**
  > Executes: Decision 1, Decision 6, CLI Surface, Errors, Observability, `REQ-AC-01`, `REQ-AC-02`, `REQ-AC-03`, `REQ-AC-29`, `REQ-AC-30`, `REQ-AC-31`, `REQ-AC-37`
  - [x] Extend `internal/sprint/service.go` with use cases for validating, prompting, dry-running, and flowing the `technical-handbook` stage.
  - [x] Update `internal/app/sprint_commands.go` so `validate technical-handbook`, `prompt technical-handbook`, and `flow --to technical-handbook` parse arguments, call sprint service methods, render output, and map exit codes only.
  - [x] Update `internal/app/app.go` help and dispatch so the sprint command family exposes `technical-handbook` consistently.
  - [x] Reject malformed command arguments and unsupported stages such as `execute`, `smoke`, `review`, and `issues` with exit code `2`.
  - [x] Return exit code `0` for successful validation, prompt rendering, dry-run, and completed flow; exit code `5` for validation and selected-evidence failures; and the established runtime exit code for runtime failures.
  - [x] Keep text output calm, deterministic, stdout/stderr-separated, no ANSI by default, and redacted.

- [x] **Task 8: Add Focused Test Coverage**
  > Executes: Decision 7, Testing, all acceptance criteria with test requirements
  - [x] Add `internal/sprint/handbook_test.go` coverage for valid handbook content, missing file, empty file, placeholders, every required section, selected-evidence traces, unselected evidence references, decision-language rejection, allowed evidence language, selected evidence loading, path escapes, missing/unreadable files, project-index mismatches, duplicate ordering, and deterministic diagnostics.
  - [x] Add or extend prompt rendering tests for required variables, selected evidence manifest, output path, required sections, no-decision instructions, selected-only instructions, no-mutation instructions, workspace-relative paths, no runtime invocation, no handbook write, and deterministic ordering.
  - [x] Add `internal/sprint/flow_test.go` coverage for dry-run, fake-runtime successful generation, runtime success with missing artifact, runtime success with invalid artifact, runtime failure, selected-evidence failure before runtime, validation failure after runtime, flow-state write failure, safe failure messages, no later-stage completion, and area-reasoning ready/skipped behavior.
  - [x] Add `internal/app/sprint_commands_test.go` coverage for validate, prompt, flow, help output, unsupported stages, exit codes, stdout/stderr separation, no ANSI output, safe redaction, runtime-free validate/prompt paths, and fake-runtime flow seam.
  - [x] Include import-boundary checks in tests or review evidence proving `internal/sprint` does not import `internal/study`, `internal/platform/*` does not import product modules, and no prohibited global packages were introduced.

- [x] **Task 9: Verify And Prepare Review Evidence**
  > Executes: Decision 7, Review Expectations, `REQ-AC-38`, `REQ-AC-39`, `REQ-AC-40`
  - [x] Run `go test ./...` from `/home/antonioborgerees/coding/ultraplan-go` and record results.
  - [x] Run `go test -race ./...` from `/home/antonioborgerees/coding/ultraplan-go` and record results.
  - [x] Run `go build ./cmd/ultraplan` from `/home/antonioborgerees/coding/ultraplan-go` and record results.
  - [x] Inspect command help for `technical-handbook` support and unsupported-stage rejection.
  - [x] Prepare `projects/ultraplan-go/sprints/19-distill-stage/review.md` with implementation evidence, deviations, verification results, architecture review findings, and acceptance-criteria conclusions.

## Evidence Checklist

- [x] Tests prove technical-handbook validation rejects missing, empty, placeholder, incomplete, untraced, unselected-evidence, and decision-language content.
- [x] Tests prove selected evidence loading reads only cataloged and selected reports, rejects unsafe paths, and orders diagnostics deterministically.
- [x] Tests prove prompt preview includes all required variables and constraints without runtime invocation or handbook mutation.
- [x] Tests prove dry-run reports stages and selected evidence without runtime, handbook writes, or completion state changes.
- [x] Tests prove runtime-backed flow completes only after runtime success, generated artifact existence, handbook validation success, and atomic flow-state update.
- [x] Tests prove failed flow records or preserves safe state and does not ready or complete later reasoning or plan stages.
- [x] Tests prove `area-reasoning` is skipped only when no reasoning templates are selected; otherwise it remains missing or ready.
- [x] Tests prove unsupported future stages and deferred artifacts remain excluded from Planning Phase 2 support.
- [x] Tests prove CLI handlers are thin, validate/prompt are runtime-free, flow uses a fake runtime seam in tests, output is deterministic, stdout/stderr are separated, ANSI is absent by default, redaction is applied, and exit codes match the sprint requirements.
- [x] Runtime or diagnostic evidence proves non-dry-run flow uses only the generic platform runtime boundary.
- [x] Architecture review evidence proves `internal/sprint` owns distill rules, `internal/app` remains thin, platform packages stay product-agnostic, `internal/study` is not imported by sprint, and prohibited global packages were not introduced.
- [x] Documentation/help evidence proves `technical-handbook` is discoverable for validate, prompt, and flow commands without adding execute, smoke, review, issues, or Git behavior.
- [x] Deviations from `reasoning.md` are recorded before implementation continues.
- [x] Required review protocols have evidence in `review.md`.

## Verification Commands

| Check | Command | Expected Result |
| --- | --- | --- |
| Unit and fixture tests | `go test ./...` | Passes from `/home/antonioborgerees/coding/ultraplan-go` with deterministic offline tests. |
| Race verification | `go test -race ./...` | Passes from `/home/antonioborgerees/coding/ultraplan-go`. |
| CLI build | `go build ./cmd/ultraplan` | Builds successfully from `/home/antonioborgerees/coding/ultraplan-go`. |
| Sprint help discovery | `ultraplan sprint --help` | Shows the sprint command family without unsupported future stages as active behavior. |
| Validate help discovery | `ultraplan sprint <project> <sprint> validate --help` | Shows `technical-handbook` validation support and usage errors for unsupported stages. |
| Prompt help discovery | `ultraplan sprint <project> <sprint> prompt --help` | Shows `technical-handbook` prompt preview support without runtime execution. |
| Flow help discovery | `ultraplan sprint <project> <sprint> flow --help` | Shows `--to technical-handbook`, `--dry-run`, and unsupported-stage rejection expectations. |

## Risks And Blockers

| Risk / Blocker | Source | Mitigation | Status |
| --- | --- | --- | --- |
| Existing sprint-index parsing may not expose selected evidence in a reusable deterministic structure. | `reasoning.md` assumptions | Keep parsing in `internal/sprint`, add focused fixtures for current Sprint 19 index shape, and fail with diagnostics rather than guessing. | open |
| `project-index.md` catalog APIs may not expose enough structure for selected-evidence comparison. | `reasoning.md` assumptions | Use `internal/project` exported catalog behavior where available; add only focused catalog read methods if needed. | open |
| No area reasoning file exists even though the sprint index selected Architecture reasoning. | `reasoning.md` assumptions | Execute final sprint decisions directly from `reasoning.md`; do not require a missing area file to reopen architecture. | mitigated |
| No-decision checks may reject legitimate evidence wording. | `reasoning.md` risks | Target prohibited sections and final-decision phrases; test allowed trade-off, guidance, warning, and open-question wording. | open |
| Markdown table parsing for selected evidence may be brittle. | `reasoning.md` risks | Use deterministic parser tests, actionable diagnostics, and no fallback discovery that expands scope. | open |
| Flow-state write failure can occur after artifact generation succeeds. | `reasoning.md` risks | Return non-zero, preserve last valid state, record safe failure, and require rerun or revalidation before marking complete. | open |
| Existing Sprint 19 `flow-state.json` may contain unsupported legacy stages. | `reasoning.md` risks | Treat unsupported stages as invalid current Planning Phase 2 state; do not preserve them as supported behavior. | open |
| Runtime result shape may not expose all desired diagnostics. | `reasoning.md` risks | Use safe platform/runtime result fields and cause chains; do not parse OpenCode/native streams directly. | open |
| Help-output tests can become brittle. | `reasoning.md` risks | Assert required command names, supported stages, unsupported-stage behavior, and key usage lines instead of freezing all prose unless existing style requires goldens. | open |

## Review Inputs

Review should use:

- `projects/ultraplan-go/sprints/19-distill-stage/sprint-index.md`
- `projects/ultraplan-go/sprints/19-distill-stage/technical-handbook.md`
- `projects/ultraplan-go/sprints/19-distill-stage/reasoning.md`
- `projects/ultraplan-go/sprints/19-distill-stage/plan.md`
- `reasoning/*.md` where created later, if any
- implementation diff under `/home/antonioborgerees/coding/ultraplan-go`
- verification evidence from `go test ./...`, `go test -race ./...`, and `go build ./cmd/ultraplan`
- `system/protocols/architecture-review-protocol.md`
- `system/protocols/sprint-review-protocol.md`
- `projects/ultraplan-go/sprints/19-distill-stage/review.md`

## Execution Log

| Date / Step | Action | Evidence / Notes |
| --- | --- | --- |
| 2026-06-14 / planning | Created sprint implementation plan from `reasoning.md`. | Plan carries forward decisions, expected evidence, assumptions, risks, and review protocols; implementation has not started. |
| 2026-06-14 / implementation | Implemented technical-handbook selected-evidence manifest, validation, prompt rendering, flow transitions, service methods, and CLI wiring. | Changed `internal/sprint/handbook.go`, `internal/sprint/prompts.go`, `internal/sprint/flow.go`, `internal/sprint/service.go`, `internal/sprint/store_fs.go`, `internal/app/sprint_commands.go`, and focused tests. |
| 2026-06-14 / verification | Ran required offline checks from `/home/antonioborgerees/coding/ultraplan-go`. | `go test ./...` passed; `go test -race ./...` passed; `go build ./cmd/ultraplan` passed. |
| 2026-06-14 / help and boundary evidence | Inspected sprint help and import boundaries. | Help shows `validate technical-handbook`, `prompt technical-handbook`, and `flow --to technical-handbook`; `rg 'internal/study|internal/(planning|workflow|validation|prompts|stages|reports)' internal/sprint internal/platform -n` returned no matches. Live unsupported-stage CLI check against `.ultra` was blocked by workspace config mismatch (`ultraplan.yml` absent); unsupported target rejection is covered by command tests. |

## Completion Criteria

- [x] All tasks are complete or explicitly deferred with a reason tied to `requirements.md` or `reasoning.md`.
- [x] `validate technical-handbook`, `prompt technical-handbook`, and `flow --to technical-handbook` satisfy the required behavior without expanding into deferred stages.
- [x] Verification commands were run from `/home/antonioborgerees/coding/ultraplan-go`, or any deferrals are documented in `review.md` with a blocker and follow-up.
- [x] Evidence satisfies the expectations from `reasoning.md` and the acceptance criteria in `requirements.md`.
- [x] Architecture and sprint review protocol evidence is recorded in `review.md`.
- [x] `review.md` can evaluate conformance without guessing implementation intent.
