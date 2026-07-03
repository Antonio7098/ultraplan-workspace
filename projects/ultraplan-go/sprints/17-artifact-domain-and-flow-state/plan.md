# Sprint Plan: Sprint Artifact Domain and Flow State

> Project: `ultraplan-go`
> Sprint: `17-artifact-domain-and-flow-state`
> Source: `reasoning.md`
> **Inputs Used:** `projects/ultraplan-go/sprints/17-artifact-domain-and-flow-state/requirements.md`, `projects/ultraplan-go/sprints/17-artifact-domain-and-flow-state/reasoning.md`, `projects/ultraplan-go/sprints/17-artifact-domain-and-flow-state/sprint-index.md`, `projects/ultraplan-go/sprints/17-artifact-domain-and-flow-state/technical-handbook.md`, `projects/ultraplan-go/docs/ARCHITECTURE.md`, `projects/ultraplan-go/docs/PRD.md`, `projects/ultraplan-go/docs/TRD.md`, `projects/ultraplan-go/project-index.md`, `templates/sprint-plan.md`

This plan executes `reasoning.md`. It must not invent architecture, scope, or decisions.

No `reasoning/*.md` area files were present for this sprint. Deeper final study reports were not reopened for this plan because `technical-handbook.md` and `reasoning.md` already distilled the selected reports and no plan decision required additional source-report evidence.

## Reasoning Source

- **Sprint Reasoning:** `projects/ultraplan-go/sprints/17-artifact-domain-and-flow-state/reasoning.md`
- **Sprint Index:** `projects/ultraplan-go/sprints/17-artifact-domain-and-flow-state/sprint-index.md`
- **Technical Handbook:** `projects/ultraplan-go/sprints/17-artifact-domain-and-flow-state/technical-handbook.md`
- **Area Reasoning:** none present; the selected architecture area reasoning artifact was missing and final sprint reasoning resolved the architecture decisions directly.

## Sprint Status

- **Status:** `complete`
- **Owner:** `implementation agent`
- **Start Date:** `2026-06-14`
- **Completion Date:** `2026-06-14`

## Decisions To Execute

| Decision | Source Section | Execution Implication |
| --- | --- | --- |
| `internal/sprint` owns sprint planning state | `reasoning.md#decision-1-internalsprint-owns-sprint-planning-state` | Add a cohesive `internal/sprint` package for sprint discovery, domain, artifact paths, flow-state persistence, status derivation, and service behavior; keep `internal/app` as wiring/rendering only. |
| Model exactly six planning stages and seven planning artifacts | `reasoning.md#decision-2-model-exactly-six-planning-stages-and-seven-planning-artifacts` | Define only `requirements`, `sprint-index`, `technical-handbook`, `area-reasoning`, `reasoning`, and `plan`; define only the listed planning artifacts plus `flow-state.json`; reject execution/review/smoke/issue stages. |
| Store strict versioned flow state with workspace-relative safe paths | `reasoning.md#decision-3-store-strict-versioned-flow-state-with-workspace-relative-safe-paths` | Implement schema version `1`, required project/sprint/timestamp/stage fields, exact statuses, strict stage/status/path/project/sprint validation, and workspace-relative artifact paths contained by sprint root. |
| Refresh status from artifacts only after state is missing or valid | `reasoning.md#decision-4-refresh-status-from-artifacts-only-after-state-is-missing-or-valid` | Invalid persisted state must fail before refresh; missing or valid state may be refreshed from artifact inspection; recorded failed stages remain failed. |
| Derive stage status deterministically with narrow artifact rules | `reasoning.md#decision-5-derive-stage-status-deterministically-with-narrow-artifact-rules` | Inspect only expected sprint artifacts; non-empty files become complete; first unblocked missing stage becomes ready; `area-reasoning` skips only on explicit no-template evidence. |
| Add only thin CLI wiring for `sprint status` | `reasoning.md#decision-6-add-only-thin-cli-wiring-for-sprint-status` | Add `ultraplan sprint <project> <sprint> status` help, argument parsing, service call, deterministic text rendering, and exit-code mapping without app-layer business rules. |
| Keep persistence atomic and failure-loud | `reasoning.md#decision-7-keep-persistence-atomic-and-failure-loud` | Write `flow-state.json` through same-directory temp file, flush/close, rename, and best-effort parent sync; preserve prior valid state on write failures and return non-zero status. |
| Verify with offline behavior tests and review checks | `reasoning.md#decision-8-verify-with-offline-behavior-tests-and-review-checks` | Add sprint package and command tests, run required Go test/build/race commands, and record review evidence after implementation. |

## Requirements / Contracts To Satisfy

| Contract / Requirement ID | Required Behavior | Evidence Planned |
| --- | --- | --- |
| Architecture contract; `requirements.md:30,49-52` | `internal/sprint` owns sprint behavior; no study dependency; platform remains product-agnostic; no global workflow/planning packages. | Import review, `internal/sprint/doc.go`, package layout review, `go test ./...`. |
| CLI Surface contract; `requirements.md:21-23,45-48,54` | Expose only `ultraplan sprint <project> <sprint> status` for this sprint, with help, deterministic text, stdout/stderr separation, stable exit codes, and no ANSI. | `internal/app/sprint_commands_test.go`, help-output tests, command-output assertions. |
| Persistence And Migrations contract; `requirements.md:37-42` | Versioned strict `flow-state.json`, one entry per supported stage, exact statuses, strict loading, atomic writes, write-failure preservation. | `internal/sprint/sprint_test.go` strict load/write tests and write-failure tests. |
| Security contract; `requirements.md:31-33,40,80-83` | Resolve projects/sprints safely, validate filesystem-safe slugs, reject path escapes, keep diagnostics safe, avoid runtime/shell/Git/network invocations. | Slug/reference/path tests, unsafe flow-state path tests, runtime-free command tests, review checklist. |
| Observability contract; `requirements.md:43-49` | Status reports project, sprint, sprint root, flow-state path, all stages, statuses, paths, and safe failures truthfully. | Stable status-output tests for empty, partial, complete, failed, skipped, and invalid-state fixtures. |
| Errors contract; `requirements.md:47-48` | Invalid `flow-state.json` returns exit code `5`; usage/workspace/reference failures use established app exit codes; causes are preserved internally. | Error classification tests, command tests for stderr diagnostics and exit codes. |
| Testing contract; `requirements.md:53-58,82` | Tests cover discovery, resolution, stage ordering, artifacts, strict flow-state loading, atomic writes, command behavior, offline verification, race, and build. | `go test ./...`, `go test -race ./...`, `go build ./cmd/ultraplan`. |
| Documentation contract; `requirements.md:15,23,100-114` | Package docs and help explain planning-stage-only scope, dependencies, status behavior, deferred stages, and review expectations. | `internal/sprint/doc.go`, command help tests, final `review.md`. |
| PRD Phase 2 planning scope | Planning side stops at `plan.md`; no implementation execution, smoke, review automation, issue tracking, or Git mutation. | Diff review for absent deferred commands/stages; command tests reject or omit deferred behavior. |
| TRD 18.1-18.5, 18.8-18.9 | `internal/project`/`internal/sprint` ownership, planning stages/artifacts, flow-state fields/statuses, status command, deferred technical requirements. | Implementation review, domain/state/service tests, command tests. |

## Tasks

- [x] **Task 1: Establish Sprint Package Scope And Domain**
  > Executes: `Decision 1`, `Decision 2`, Architecture contract, `requirements.md:15-16,30,34-36,49-52`
  - [x] Add `/home/antonioborgerees/coding/ultraplan-go/internal/sprint/doc.go` documenting package ownership, planning scope through `plan.md`, allowed dependencies, and deferred implementation/smoke/review/issue/Git behavior.
  - [x] Add `/home/antonioborgerees/coding/ultraplan-go/internal/sprint/domain.go` with `Sprint`, `PlanningStage`, `StageStatus`, `StageState`, `FlowState`, status summary types, and validation findings.
  - [x] Define stage order as exactly `requirements`, `sprint-index`, `technical-handbook`, `area-reasoning`, `reasoning`, `plan`.
  - [x] Define allowed statuses as exactly `missing`, `ready`, `complete`, `failed`, and `skipped`.
  - [x] Keep unsupported stages such as `implementation`, `execute`, `smoke`, `review`, and `issues` out of domain constants and strict accepted values.

- [x] **Task 2: Implement Discovery, Reference Resolution, And Artifact Paths**
  > Executes: `Decision 1`, `Decision 2`, `Decision 5`, Security contract, `requirements.md:17-18,31-36,80-81`
  - [x] Add `/home/antonioborgerees/coding/ultraplan-go/internal/sprint/discovery.go` to inspect only direct child directories under `projects/<project>/sprints/`, ignore hidden and non-directory entries, validate filesystem-safe slugs, and sort deterministically.
  - [x] Implement exact and unambiguous-prefix sprint reference resolution with actionable missing, invalid, and ambiguous diagnostics.
  - [x] Add `/home/antonioborgerees/coding/ultraplan-go/internal/sprint/artifacts.go` with centralized expected paths for `requirements.md`, `sprint-index.md`, `technical-handbook.md`, `reasoning/`, `reasoning.md`, `plan.md`, and `flow-state.json`.
  - [x] Normalize artifact paths and ensure each expected path remains inside the selected sprint root before filesystem access or persistence.
  - [x] Store and render artifact paths as workspace-relative paths, not absolute persisted paths.

- [x] **Task 3: Implement Strict Flow-State Persistence**
  > Executes: `Decision 3`, `Decision 4`, `Decision 7`, Persistence And Migrations contract, Security contract, Errors contract, `requirements.md:19,37-42,47,79-80`
  - [x] Add `/home/antonioborgerees/coding/ultraplan-go/internal/sprint/state.go` to load, validate, refresh, and write schema version `1` `flow-state.json`.
  - [x] Validate required top-level fields: schema version, project name, sprint slug, updated timestamp, and one state entry per supported planning stage.
  - [x] Validate each stage state field: supported stage, supported status, workspace-relative artifact path inside the sprint root, optional last-run timestamp, and optional safe error detail.
  - [x] Reject malformed JSON, unsupported schema versions, missing required fields, unknown stages, duplicate stages, unknown statuses, project mismatches, sprint mismatches, and path escapes.
  - [x] Ensure invalid persisted state fails before artifact refresh and is not repaired or overwritten by `status`.
  - [x] Write state atomically through same-directory temp file, flush and close, rename over previous state, and best-effort sync the parent directory.
  - [x] Preserve the last valid `flow-state.json` when temp write, flush, close, or rename fails, and preserve the state path plus underlying cause internally.

- [x] **Task 4: Implement Filesystem Store And Status Service**
  > Executes: `Decision 1`, `Decision 4`, `Decision 5`, Observability contract, Security contract, `requirements.md:20-21,43-49,61-68`
  - [x] Add `/home/antonioborgerees/coding/ultraplan-go/internal/sprint/store_fs.go` to read planning artifact existence, file non-emptiness, reasoning directory entries, and flow-state files through workspace-safe paths.
  - [x] Add `/home/antonioborgerees/coding/ultraplan-go/internal/sprint/service.go` for runtime-free status inspection and flow-state refresh.
  - [x] Derive file-backed stages as `complete` only for present non-empty files.
  - [x] Mark blocked later missing stages as `missing` and the first missing stage after completed predecessors as `ready`.
  - [x] Mark `area-reasoning` as `complete` when `reasoning/` contains at least one readable non-empty Markdown file.
  - [x] Mark `area-reasoning` as `skipped` only when readable `sprint-index.md` explicitly selects no reasoning templates; otherwise do not silently skip missing reasoning files.
  - [x] Preserve recorded `failed` stages and safe error detail during artifact-derived refresh.
  - [x] Ensure status inspection does not invoke agentwrap, OpenCode, runtime health checks, prompt generation, network calls, source cloning, shell execution, Git, project-index mutation, planning-stage generation, or study run-loop behavior.

- [x] **Task 5: Wire Sprint CLI Command Thinly**
  > Executes: `Decision 6`, CLI Surface contract, Errors contract, Documentation contract, `requirements.md:22-23,45-48,54,75`
  - [x] Add `/home/antonioborgerees/coding/ultraplan-go/internal/app/sprint_commands.go` for `ultraplan sprint <project> <sprint> status` help, argument parsing, service invocation, rendering, and exit-code mapping.
  - [x] Update `/home/antonioborgerees/coding/ultraplan-go/internal/app/app.go` to register the `sprint` command family and document it in top-level help without disturbing existing commands.
  - [x] Keep stage ordering, artifact paths, flow-state schema validation, status derivation, and atomic-write behavior out of `internal/app`.
  - [x] Render deterministic plain text including project, sprint, sprint root, flow-state path, every supported stage, each status, each artifact path, and safe error details where present.
  - [x] Map invalid flow-state content, unsupported stages/statuses, project/sprint mismatch, and unsafe artifact paths to exit code `5`.
  - [x] Map missing projects, missing sprint directories, invalid references, and malformed arguments to established workspace or usage exit codes.
  - [x] Do not add `sprint validate`, `sprint prompt`, `sprint flow`, JSON status output, execution, smoke, review, issue, or Git commands in this sprint.

- [x] **Task 6: Add Sprint Package Behavior Tests**
  > Executes: `Decision 8`, Testing contract, Security contract, Persistence And Migrations contract, `requirements.md:24,53,55-58,82`
  - [x] Add `/home/antonioborgerees/coding/ultraplan-go/internal/sprint/sprint_test.go` tests for exact stage order, exact allowed statuses, exact artifact paths, unsupported stages, and unsupported statuses.
  - [x] Cover sprint discovery with direct child directories, hidden entries, non-directory entries, invalid slugs, deterministic sorting, exact refs, prefix refs, ambiguous refs, and missing refs.
  - [x] Cover status derivation for no artifacts, partial artifacts, all planning artifacts present, selected area reasoning files, explicit no-template skip, missing/unreadable sprint-index behavior, and failed-state preservation.
  - [x] Cover strict flow-state loading for malformed JSON, unsupported schema version, missing fields, unknown stages, duplicate stages, unknown statuses, unsupported `implementation` and `review`, project/sprint mismatch, unsafe paths, and path escape attempts.
  - [x] Cover atomic writes, temp-write failure, flush/close failure where feasible, rename failure, parent-sync best-effort behavior, and preservation of prior valid state.

- [x] **Task 7: Add Command Tests And Help Assertions**
  > Executes: `Decision 6`, `Decision 8`, CLI Surface contract, Observability contract, Errors contract, Documentation contract, `requirements.md:25,45-48,54`
  - [x] Add `/home/antonioborgerees/coding/ultraplan-go/internal/app/sprint_commands_test.go` tests for top-level help, `sprint --help`, `status --help`, malformed arguments, missing project, missing sprint, invalid sprint ref, ambiguous sprint ref, valid status, and invalid flow-state diagnostics.
  - [x] Assert deterministic stdout ordering, stderr-only diagnostics, expected exit codes, no ANSI escape sequences, and stable stage order.
  - [x] Assert status refresh succeeds for missing or valid state and exits `0` when sprint resolution succeeds and refreshed state is valid.
  - [x] Assert invalid flow-state content exits `5` and is not repaired or overwritten during the command.
  - [x] Assert the command path remains runtime-free by using fakes or absence checks that would fail if runtime/study/Git/network behavior is invoked.

- [x] **Task 8: Verify Boundaries, Non-Goals, And Offline Quality Gates**
  > Executes: `Decision 8`, all selected contracts, `requirements.md:49-58,100-114`
  - [x] Review imports to confirm `internal/sprint` does not import `internal/study` and platform packages do not import product packages.
  - [x] Review the diff to confirm no global `internal/planning`, `internal/workflow`, `internal/stages`, `internal/validation`, `internal/reports`, or `internal/prompts` package was introduced.
  - [x] Review the diff to confirm no implementation/smoke/review/issues/Git workflow stages, commands, state entries, or accepted strict-load values were added.
  - [x] Run `go test ./...` from `/home/antonioborgerees/coding/ultraplan-go`.
  - [x] Run `go test -race ./...` from `/home/antonioborgerees/coding/ultraplan-go`.
  - [x] Run `go build ./cmd/ultraplan` from `/home/antonioborgerees/coding/ultraplan-go`.
  - [x] Record implementation evidence, deviations, verification results, and review conclusions in `projects/ultraplan-go/sprints/17-artifact-domain-and-flow-state/review.md` after implementation.

## Evidence Checklist

- [x] Tests prove sprint discovery, reference resolution, strict stage/status modeling, artifact path rules, status derivation, strict flow-state loading, atomic writes, and write-failure preservation.
- [x] Command tests prove help, output, diagnostics, exit codes, stdout/stderr separation, deterministic ordering, no ANSI output, invalid flow-state behavior, and runtime-free operation.
- [x] Documentation updates are complete in `internal/sprint/doc.go` and command help.
- [x] Deviations from `reasoning.md` are recorded before implementation continues.
- [x] Architecture Review and Sprint Review evidence can be produced from the diff, test output, and `review.md`.
- [x] Non-goals are confirmed absent: no runtime-backed planning generation, `sprint validate`, `sprint prompt`, `sprint flow`, implementation execution, smoke investigation, automated review, issue tracking, project-index mutation, or Git mutation.

## Verification Commands

| Check | Command | Expected Result |
| --- | --- | --- |
| Offline tests | `go test ./...` from `/home/antonioborgerees/coding/ultraplan-go` | All packages pass without OpenCode, provider credentials, network access, Git commands, subprocess execution, or long sleeps. |
| Race tests | `go test -race ./...` from `/home/antonioborgerees/coding/ultraplan-go` | All packages pass under the race detector. |
| CLI build | `go build ./cmd/ultraplan` from `/home/antonioborgerees/coding/ultraplan-go` | The `ultraplan` command builds successfully. |
| Help smoke | `ultraplan sprint --help` after build | Help documents the `sprint` command family and status command without exposing deferred commands as implemented behavior. |
| Status smoke | `ultraplan sprint <project> <sprint> status` against a fixture workspace | Output includes project, sprint, sprint root, flow-state path, and all six planning stages in order. |

## Risks And Blockers

| Risk / Blocker | Source | Mitigation | Status |
| --- | --- | --- | --- |
| Existing Sprint 16 project APIs may not expose exactly the project-root boundary needed for sprint resolution. | `reasoning.md#assumptions-and-risks` | Use exported project boundaries first; if needed, add the smallest project-facing method without duplicating project catalog behavior. | open |
| Workspace path helpers may be insufficient for sprint-root containment and workspace-relative persisted paths. | `reasoning.md#assumptions-and-risks` | Reuse or add product-neutral workspace path helpers while keeping sprint artifact semantics in `internal/sprint`. | open |
| The selected Architecture area reasoning file is missing. | `reasoning.md#area-specific-reasoning-inputs` | Proceed from final `reasoning.md`, which resolves architecture decisions directly; review may request backfill if process compliance requires it. | open |
| Narrow no-template detection may not understand every future `sprint-index.md` shape. | `reasoning.md#assumptions-and-risks` | Keep detection conservative; skip only on clear explicit no-template selection; leave full selection validation to Sprint 18. | open |
| Persisted failed stage preservation may surprise users after manual artifact creation. | `reasoning.md#assumptions-and-risks` | Document behavior and keep it tested; future flow/validation commands define failure-clearing semantics. | open |
| Parent directory sync behavior can vary by filesystem. | `reasoning.md#assumptions-and-risks` | Treat parent sync as best-effort after successful rename; test core atomic write preservation. | open |
| `status` mutates `flow-state.json` during refresh. | `reasoning.md#assumptions-and-risks` | Keep mutation limited, deterministic, atomic, and documented because the sprint requirements explicitly require status refresh. | open |

## Review Inputs

Review should use:

- `projects/ultraplan-go/sprints/17-artifact-domain-and-flow-state/sprint-index.md`
- `projects/ultraplan-go/sprints/17-artifact-domain-and-flow-state/technical-handbook.md`
- `projects/ultraplan-go/sprints/17-artifact-domain-and-flow-state/reasoning.md`
- this `plan.md`
- implementation diff
- verification evidence from `go test ./...`, `go test -race ./...`, and `go build ./cmd/ultraplan`
- `system/protocols/architecture-review-protocol.md`
- `system/protocols/sprint-review-protocol.md`

## Execution Log

| Date / Step | Action | Evidence / Notes |
| --- | --- | --- |
| `2026-06-14 / planning` | Created sprint implementation plan from requirements, final reasoning, selected context, handbook, PRD, TRD, Architecture, project index, and sprint-plan template. | No implementation code changed; area reasoning files were absent; deeper final study reports were not reopened because the handbook and reasoning already supplied decision evidence. |
| `2026-06-14 / implementation` | Implemented `internal/sprint` package for planning-stage domain, discovery, artifact paths, strict flow-state load/write, filesystem artifact inspection, status derivation, and service refresh. | Added `/home/antonioborgerees/coding/ultraplan-go/internal/sprint/{doc.go,domain.go,discovery.go,artifacts.go,state.go,store_fs.go,service.go,sprint_test.go}`. |
| `2026-06-14 / CLI wiring` | Registered `ultraplan sprint <project> <sprint> status` as thin app wiring with deterministic plain-text rendering and exit-code mapping. | Added `/home/antonioborgerees/coding/ultraplan-go/internal/app/sprint_commands.go`, `/home/antonioborgerees/coding/ultraplan-go/internal/app/sprint_commands_test.go`; updated app help. |
| `2026-06-14 / verification` | Ran required offline verification. | `go test ./...` passed; `go test -race ./...` passed; `go build ./cmd/ultraplan` passed; `./ultraplan sprint --help` printed scoped sprint status help. |
| `2026-06-14 / boundary review` | Reviewed imports and command surface for non-goals. | `internal/sprint` imports project/workspace plus stdlib only; no global planning/workflow/stages/validation/reports/prompts package added; no `sprint validate`, `sprint prompt`, `sprint flow`, implementation, smoke, review, issue, or Git command added. |

## Completion Criteria

- [x] All implementation tasks are complete or explicitly deferred with requirement impact recorded.
- [x] Verification commands were run or deferrals are documented with cause.
- [x] Evidence satisfies the expectations from `reasoning.md`.
- [x] `review.md` can evaluate conformance without guessing intent.
- [x] Architecture and sprint review protocols have enough evidence to verify boundaries, non-goals, tests, and deviations.
