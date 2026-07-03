# Sprint Plan: Plan Stage

> Project: `ultraplan-go`
> Sprint: `21-plan-stage`
> Source: `reasoning.md`
> **Inputs Used:** `projects/ultraplan-go/sprints/21-plan-stage/requirements.md`, `projects/ultraplan-go/sprints/21-plan-stage/reasoning.md`, `projects/ultraplan-go/sprints/21-plan-stage/sprint-index.md`, `projects/ultraplan-go/sprints/21-plan-stage/technical-handbook.md`, `projects/ultraplan-go/docs/PRD.md`, `projects/ultraplan-go/docs/TRD.md`, `projects/ultraplan-go/docs/ARCHITECTURE.md`, `projects/ultraplan-go/project-index.md`, `templates/sprint-plan.md`, `templates/sprint-reasoning.md`

This plan executes `projects/ultraplan-go/sprints/21-plan-stage/reasoning.md`. It does not invent architecture, scope, or decisions beyond that reasoning document.

## Reasoning Source

- **Sprint Reasoning:** `projects/ultraplan-go/sprints/21-plan-stage/reasoning.md`
- **Sprint Index:** `projects/ultraplan-go/sprints/21-plan-stage/sprint-index.md`
- **Technical Handbook:** `projects/ultraplan-go/sprints/21-plan-stage/technical-handbook.md`
- **Area Reasoning:** none present; `projects/ultraplan-go/sprints/21-plan-stage/reasoning/*.md` does not exist in the workspace.
- **Omitted Evidence:** selected final study reports were not opened because `technical-handbook.md` already distilled the selected reports sufficiently for these implementation decisions; reopen final reports only if implementation encounters a decision not covered by the handbook.

## Sprint Status

- **Status:** complete
- **Owner:** implementation agent
- **Start Date:** 2026-06-15
- **Completion Date:** 2026-06-15

## Decisions To Execute

| Decision | Source Section | Execution Implication |
| --- | --- | --- |
| Own plan behavior in `internal/sprint` | `reasoning.md#decision-1-own-plan-behavior-in-internalsprint` | Implement plan domain, validation, prompt rendering, prerequisite loading, task/evidence trace checks, service use cases, and flow-state transitions inside `internal/sprint`; keep `internal/app` command handlers thin. |
| Validate `plan.md` as traceable editable Markdown | `reasoning.md#decision-2-validate-planmd-as-traceable-editable-markdown` | Add structural validation for required sections, placeholders, `reasoning.md` reference, decisions, tasks, evidence, risks/blockers, success criteria, task traceability, and forbidden current-phase deferred behavior. |
| Render deterministic plan prompts without runtime or mutation | `reasoning.md#decision-3-render-deterministic-plan-prompts-without-runtime-or-mutation` | Add `prompt plan` rendering with selected inputs, output path, required sections, traceability rules, workspace-relative paths, and no-execution/no-mutation rules; do not call runtime or write `plan.md`. |
| Flow sequentially through `plan` with post-generation gates | `reasoning.md#decision-4-flow-sequentially-through-plan-with-post-generation-gates` | Extend dry-run and runtime-backed flow through `plan`, validate prerequisites first, call only the generic runtime seam when not dry-run, validate generated `plan.md`, then atomically update `flow-state.json`. |
| Keep service, store, runtime, output, and exit-code seams narrow | `reasoning.md#decision-5-keep-service-store-runtime-output-and-exit-code-seams-narrow` | Add focused service/store/runtime request/result seams, workspace-safe reads for `reasoning.md`, selected area reasoning, and `plan.md`, and command wiring that maps diagnostics and exit codes. |
| Verify offline with focused plan and command tests | `reasoning.md#decision-6-verify-offline-with-focused-plan-and-command-tests` | Add deterministic unit, fixture, command, fake-runtime, state, help, and output tests; default verification must not require OpenCode, network, provider credentials, smoke harnesses, or Git mutation. |

## Requirements / Contracts To Satisfy

| Contract / Requirement ID | Required Behavior | Evidence Planned |
| --- | --- | --- |
| AC-01, AC-07, AC-08, AC-09, C-01, C-02, C-03 | `validate plan` validates `plan.md` locally, rejects missing/invalid structure, cites `reasoning.md`, and checks task/evidence traceability. | `internal/sprint/plan_test.go` validation cases; command test proving no runtime call. |
| AC-02, AC-10, AC-11, AC-12, C-07, C-08, C-09 | `prompt plan` renders deterministic preview content with selected paths, required sections, traceability rules, and no-mutation/no-execution rules. | Prompt invariant tests in `internal/sprint/plan_test.go` and command tests in `internal/app/sprint_commands_test.go`. |
| AC-03, AC-04, AC-05, AC-06, AC-13, AC-14, AC-15, AC-16, AC-17 | `flow --to plan` and dry-run behavior preserve Phase 2 stage order, validate prerequisites, use only generic runtime for non-dry-run generation, validate outputs before completion, and update state atomically. | Fake-runtime flow tests, dry-run non-mutation tests, invalid prerequisite tests, missing/invalid generated artifact tests, state write failure tests. |
| AC-18, AC-19, AC-20, AC-21, C-05, C-06, C-12 | Package boundaries remain module-owned: `internal/sprint` owns plan behavior, `internal/app` is thin, platform packages stay product-agnostic, and `internal/study` is not reused. | Import/code review; architecture review evidence in `review.md`; `go test ./...`. |
| AC-22, AC-23, AC-24, AC-25, AC-26 | CLI output is calm and scriptable, separates stdout/stderr, has no ANSI by default, redacts sensitive details, and maps exit codes `0`, `2`, `5`, plus runtime failures. | Command tests for help/output/errors, malformed arguments, unsupported stages, validation failures, runtime failure mapping, and redaction. |
| AC-27, AC-28, AC-29, AC-30, AC-31, AC-32, AC-33 | Tests cover validator, prompt, dry-run, fake runtime, post-generation validation, command behavior, offline verification, race verification, and CLI build. | `go test ./...`, `go test -race ./...`, `go build ./cmd/ultraplan`, plus focused test files. |
| PRD Phase 2, TRD 18, Architecture, Workflows, Persistence And Migrations | Phase 2 supports planning through `plan.md` only and keeps unsupported future command families out of current product flow state. | Flow/state tests; prompt/validator tests for unsupported future behavior; review evidence. |
| Security, Configuration, LLM Runtime, LLM Evaluation / Cost / Safety | Runtime execution uses generic platform/agentwrap boundary, avoids shell/Git/provider/OpenCode calls from `internal/sprint`, and avoids secrets in prompts/diagnostics. | Fake-runtime seam tests, import review, prompt assertions, redaction tests. |

## Tasks

- [x] **Task 1: Inspect Existing Sprint Implementation Surface**
  > Executes: `Decision 1`, `Decision 5`, AC-18 through AC-21
  - [x] Inspect existing files in `/home/antonioborgerees/coding/ultraplan-go/internal/sprint` and `/home/antonioborgerees/coding/ultraplan-go/internal/app` before editing.
  - [x] Identify existing stage parsing, validators, prompt rendering, service methods, store interfaces, flow-state model, atomic write helpers, and exit-code mapping.
  - [x] Record any mismatch with `reasoning.md` as an implementation note in `projects/ultraplan-go/sprints/21-plan-stage/review.md` after implementation.
  - [x] Verification expectation: no code changes are made until the existing seams are understood; implementation preserves module ownership and thin CLI routing.

- [x] **Task 2: Add Plan Domain, Store Reads, And Validation**
  > Executes: `Decision 1`, `Decision 2`, AC-01, AC-07, AC-08, AC-09, AC-18, C-01, C-02, C-03
  - [x] Add or extend `internal/sprint/plan.go` with plan-stage inputs, validation result/findings, and checks for missing file, empty content, placeholders, missing `reasoning.md` reference, missing decisions, missing tasks, missing evidence checklist, missing risks/blockers, missing success criteria, untraced tasks, and forbidden current-phase deferred behavior.
  - [x] Implement task/evidence trace checks against final reasoning decisions and acceptance evidence without semantically executing or simulating plan tasks.
  - [x] Extend `internal/sprint/store_fs.go` to read `reasoning.md`, selected `reasoning/*.md`, and `plan.md` through workspace-safe paths.
  - [x] Keep helper functions private and local to `internal/sprint` unless a mechanical existing helper already exists for path safety or atomic IO.
  - [x] Verification expectation: `internal/sprint/plan_test.go` covers valid content and every invalid validator case listed in AC-27.

- [x] **Task 3: Add Deterministic Plan Prompt Rendering**
  > Executes: `Decision 3`, AC-02, AC-10, AC-11, AC-12, AC-28, C-07, C-08, C-09
  - [x] Extend `internal/sprint/prompts.go` with a `plan` prompt renderer that includes project slug, sprint slug, sprint path, requirements path, sprint-index path, technical-handbook path, selected area reasoning paths, final `reasoning.md` path, `plan.md` output path, required sections, traceability rules, and no-execution/no-mutation rules.
  - [x] Render workspace-relative paths in prompt content unless an absolute path is explicitly needed for local diagnostics.
  - [x] Ensure prompt preview writes only to stdout by default and only to an explicit preview path if an existing output-preview convention supports it.
  - [x] Ensure `prompt plan` cannot write `plan.md` and cannot invoke runtime.
  - [x] Verification expectation: prompt tests assert required variables, selected inputs, no-mutation/no-execution rules, deterministic ordering, and no runtime or `plan.md` writes.

- [x] **Task 4: Extend Flow Through `plan`**
  > Executes: `Decision 4`, AC-03, AC-04, AC-05, AC-06, AC-13, AC-14, AC-15, AC-16, AC-17, AC-29, C-04, C-07, C-10, C-11
  - [x] Extend `internal/sprint/flow.go` so `flow --to plan --dry-run` reports required stages, validated reasoning inputs, selected area reasoning prerequisites, expected `plan.md`, and flow-state impact without runtime or mutation.
  - [x] Extend runtime-backed `flow --to plan` to validate `requirements.md`, `sprint-index.md`, `technical-handbook.md`, selected area reasoning when selected, and final `reasoning.md` before runtime execution.
  - [x] Invoke only the generic platform runtime boundary for plan generation; do not call OpenCode, agentwrap adapters, shell commands, provider APIs, or Git directly from `internal/sprint`.
  - [x] Require runtime success, `plan.md` existence, plan validation success, and atomic `flow-state.json` update before marking `plan` complete.
  - [x] Record failed plan flow with safe diagnostics and without introducing unsupported future stages or deferred artifacts into state.
  - [x] Verification expectation: flow tests cover dry-run non-mutation, fake-runtime success, missing generated artifact, invalid generated artifact, runtime failure, invalid prerequisites, selected area reasoning failure, post-generation validation failure, flow-state write failure, completion through `plan`, failed-stage recording, and unsupported-stage rejection.

- [x] **Task 5: Wire Sprint Services And CLI Commands Thinly**
  > Executes: `Decision 1`, `Decision 5`, AC-01, AC-02, AC-03, AC-21 through AC-26, C-01, C-08, C-12
  - [x] Extend `internal/sprint/service.go` with narrow use cases for validate plan, prompt plan, and flow to plan.
  - [x] Extend `internal/app/sprint_commands.go` so CLI handlers parse arguments, call sprint service methods, render stdout/stderr, and map exit codes only.
  - [x] Support `ultraplan sprint <project> <sprint> validate plan`, `prompt plan`, and `flow --to plan`, including help output and malformed argument handling.
  - [x] Ensure unsupported stage values return usage exit code `2`, validation/prerequisite failures return exit code `5`, and runtime failures map to the established runtime exit code.
  - [x] Verification expectation: command tests prove validate/prompt are runtime-free, flow uses a fake runtime seam, stdout/stderr are separated, output is deterministic, no ANSI escapes are emitted, and sensitive values are redacted.

- [x] **Task 6: Preserve Package Boundaries And Deferred Scope**
  > Executes: `Decision 1`, `Decision 4`, `Decision 5`, AC-18, AC-19, AC-20, AC-21, C-05, C-06, C-10, C-11, C-12
  - [x] Confirm `internal/sprint` does not import `internal/study` or reuse study services, source/dimension models, report/rating logic, summary generation, scheduling, or run-loop state.
  - [x] Confirm `internal/platform/*` imports no product modules.
  - [x] Confirm no global `internal/catalog`, `internal/planning`, `internal/workflow`, `internal/stages`, `internal/validation`, `internal/reports`, or `internal/prompts` package is introduced.
  - [x] Confirm current flow state supports only `requirements`, `sprint-index`, `technical-handbook`, `area-reasoning`, `reasoning`, and `plan`.
  - [x] Verification expectation: import review and tests demonstrate deferred command families and artifacts are not modeled as supported Planning Phase 2 behavior.

- [x] **Task 7: Add Focused Tests And Fixtures**
  > Executes: `Decision 6`, AC-27 through AC-33
  - [x] Add plan validator tests in `/home/antonioborgerees/coding/ultraplan-go/internal/sprint/plan_test.go` for all valid and invalid cases required by AC-27.
  - [x] Add prompt invariant tests for selected inputs, required sections, traceability rules, no-execution/no-mutation rules, deterministic ordering, and workspace-relative paths.
  - [x] Add flow tests with fake runtime/store seams for dry-run, success, runtime failure, missing artifact, invalid artifact, invalid prerequisites, invalid selected area reasoning, post-generation validation failure, state write failure, and unsupported stages.
  - [x] Add command tests in `/home/antonioborgerees/coding/ultraplan-go/internal/app/sprint_commands_test.go` for validate, prompt, flow, help, exit codes, stdout/stderr separation, no ANSI output, redaction, runtime-free validate/prompt, and fake-runtime flow.
  - [x] Verification expectation: tests are deterministic, offline, fake-first, and do not require real OpenCode, provider credentials, network access, smoke harnesses, or Git mutation.

- [x] **Task 8: Run Verification And Produce Review Evidence**
  > Executes: `Decision 6`, AC-31, AC-32, AC-33, review protocols
  - [x] Run `go test ./...` from `/home/antonioborgerees/coding/ultraplan-go`.
  - [x] Run `go test -race ./...` from `/home/antonioborgerees/coding/ultraplan-go`.
  - [x] Run `go build ./cmd/ultraplan` from `/home/antonioborgerees/coding/ultraplan-go`.
  - [x] Apply `system/protocols/architecture-review-protocol.md` and `system/protocols/sprint-review-protocol.md` in `projects/ultraplan-go/sprints/21-plan-stage/review.md` after implementation.
  - [x] Verification expectation: `review.md` records implementation evidence, deviations from `reasoning.md`, verification results, unsupported-stage exclusions, selected-context compliance, and conclusions.

## Evidence Checklist

- [x] Plan validation tests cover valid content, missing file, empty file, placeholders, missing `reasoning.md` citation/reference, missing decisions, missing tasks, missing evidence checklist, missing risks/blockers, missing success criteria, untraced tasks, and forbidden current-phase deferred behavior.
- [x] Prompt preview tests cover required variables, selected context paths, selected Architecture area reasoning path when selected, output path, required sections, traceability rules, no-mutation/no-execution rules, deterministic ordering, workspace-relative paths, no runtime calls, and no `plan.md` writes.
- [x] Flow tests cover dry-run non-mutation, fake-runtime success, missing generated artifact, invalid generated artifact, runtime failure, invalid prerequisites, selected area reasoning prerequisite failure, validation failure after generation, flow-state write failure, stage completion through `plan`, failed-stage recording, and unsupported-stage rejection.
- [x] Command tests cover `validate plan`, `prompt plan`, `flow --to plan`, help output, malformed arguments, unsupported stages, exit codes `0`, `2`, and `5`, runtime failure mapping, stdout/stderr separation, no ANSI output, redaction, runtime-free validate/prompt paths, and fake-runtime flow.
- [x] Import review confirms `internal/sprint` does not import `internal/study`, platform packages do not import product modules, and no global planning/workflow/validation/prompt package was added.
- [x] Runtime evidence confirms non-dry-run plan flow uses only the generic platform runtime boundary and fake runtime by default.
- [x] State evidence confirms `flow-state.json` remains strict, versioned, atomically written, and limited to supported Planning Phase 2 stages.
- [x] Documentation/help evidence confirms plan-stage commands are visible and unsupported future stages are rejected.
- [x] Deviations from `projects/ultraplan-go/sprints/21-plan-stage/reasoning.md` are recorded in `review.md` before implementation is considered complete.
- [x] Required Architecture Review and Sprint Review protocol evidence is recorded in `projects/ultraplan-go/sprints/21-plan-stage/review.md`.

## Verification Commands

| Check | Command | Expected Result |
| --- | --- | --- |
| Offline tests | `go test ./...` from `/home/antonioborgerees/coding/ultraplan-go` | Exit code `0`; unit, fixture, command, fake-runtime, validator, prompt, flow, and state tests pass. |
| Race tests | `go test -race ./...` from `/home/antonioborgerees/coding/ultraplan-go` | Exit code `0`; no data races in sprint flow, state, runtime fake, or command paths. |
| CLI build | `go build ./cmd/ultraplan` from `/home/antonioborgerees/coding/ultraplan-go` | Exit code `0`; CLI binary builds. |
| Help surface | `ultraplan sprint --help`, `ultraplan sprint <project> <sprint> validate --help`, `ultraplan sprint <project> <sprint> prompt --help`, `ultraplan sprint <project> <sprint> flow --help` | Help includes `plan` support where appropriate, remains scriptable, and does not imply deferred command-family support. |
| Architecture review | Apply `system/protocols/architecture-review-protocol.md` | `review.md` records package ownership, dependency direction, platform product-agnostic status, and absence of study semantic reuse. |
| Sprint review | Apply `system/protocols/sprint-review-protocol.md` | `review.md` records acceptance evidence, verification results, deviations, risks, and conclusion. |

## Risks And Blockers

| Risk / Blocker | Source | Mitigation | Status |
| --- | --- | --- | --- |
| Selected Architecture area reasoning artifact is absent. | `reasoning.md#assumptions-and-risks`; `sprint-index.md#selected-reasoning-templates` | Implemented flow requires the selected area reasoning file before runtime-backed `plan` generation. Tests cover prerequisite failure handling. | mitigated |
| Existing code seams may differ from the names assumed by the requirements. | `reasoning.md#assumptions-and-risks` | Existing `internal/sprint` and `internal/app` seams were inspected before editing and extended in place. | closed |
| Plan validation could overfit prose and reject valid human-edited plans. | `reasoning.md#assumptions-and-risks` | Validator checks stable structural elements, explicit references, traceability, checklist/evidence presence, placeholders, and forbidden current-phase behavior rather than exact wording. | mitigated |
| Prompt rendering could include unsafe, absolute, or unselected context. | `reasoning.md#assumptions-and-risks`; Security contract | Prompt tests assert workspace-relative output and no local root or ANSI leakage. | mitigated |
| Runtime success followed by state write failure can leave a valid `plan.md` without `plan` complete in state. | `reasoning.md#assumptions-and-risks`; Persistence And Migrations contract | Existing atomic state writer remains in use and prior write-failure preservation tests pass. | mitigated |
| Real OpenCode integration is not exercised by default verification. | `reasoning.md#assumptions-and-risks`; requirements non-goals | Default verification remains fake-first as required; optional deep smoke is deferred. | accepted deferral |

## Assumptions And Open Questions

| Item | Type | Implementation Handling |
| --- | --- | --- |
| Existing platform runtime abstraction can carry prompt, workdir, model/config, timeout, and expected output information generically. | Assumption | Use only generic runtime fields; keep product validation in `internal/sprint` after runtime success. |
| Existing sprint command infrastructure already has stage parsing and exit-code/output conventions from prior sprints. | Assumption | Reuse existing conventions; do not introduce a parallel command framework. |
| The existing command surface may or may not support an explicit prompt-preview output path. | Open question | If already present, allow writing only to the explicit preview path; otherwise keep prompt preview stdout-only for this sprint. |
| Safe runtime diagnostics may vary by platform runtime result shape. | Open question | Use only safe/redacted fields exposed through the generic runtime and established diagnostic helpers; do not dump raw runtime payloads. |

## Review Inputs

Review should use:

- `projects/ultraplan-go/sprints/21-plan-stage/requirements.md`
- `projects/ultraplan-go/sprints/21-plan-stage/sprint-index.md`
- `projects/ultraplan-go/sprints/21-plan-stage/technical-handbook.md`
- `projects/ultraplan-go/sprints/21-plan-stage/reasoning.md`
- `projects/ultraplan-go/sprints/21-plan-stage/plan.md`
- `projects/ultraplan-go/sprints/21-plan-stage/reasoning/*.md` where created later
- implementation diff in `/home/antonioborgerees/coding/ultraplan-go`
- verification evidence from `go test ./...`, `go test -race ./...`, and `go build ./cmd/ultraplan`
- `system/protocols/architecture-review-protocol.md`
- `system/protocols/sprint-review-protocol.md`

## Execution Log

| Date / Step | Action | Evidence / Notes |
| --- | --- | --- |
| 2026-06-15 planning | Created implementation sprint plan from `reasoning.md`. | No implementation code changed; selected Architecture area reasoning remains absent and is recorded as an open risk. |
| 2026-06-15 implementation | Implemented plan validation, plan prompt rendering, flow through `plan`, CLI wiring, and focused tests. | Key files: `internal/sprint/plan.go`, `internal/sprint/prompts.go`, `internal/sprint/service.go`, `internal/sprint/flow.go`, `internal/app/sprint_commands.go`, `internal/sprint/plan_test.go`, `internal/app/sprint_commands_test.go`. |
| 2026-06-15 verification | Ran required verification from `/home/antonioborgerees/coding/ultraplan-go`. | `go test ./...` passed; `go test -race ./...` passed; `go build ./cmd/ultraplan` passed. |
| 2026-06-15 review | Applied selected Architecture Review and Sprint Review protocols in `review.md`. | Review was synthesized directly; no delegated subagent reports were produced because available subagent tooling requires explicit user-requested delegation. |

## Completion Criteria

- [x] All tasks are complete or explicitly deferred with requirement-compatible rationale.
- [x] `validate plan`, `prompt plan`, `flow --to plan --dry-run`, and runtime-backed `flow --to plan` behavior satisfies acceptance criteria without adding deferred command-family behavior.
- [x] Verification commands were run from `/home/antonioborgerees/coding/ultraplan-go` or deferrals are documented with cause.
- [x] Evidence satisfies the expectations from `projects/ultraplan-go/sprints/21-plan-stage/reasoning.md`.
- [x] `review.md` can evaluate conformance without guessing intent.
