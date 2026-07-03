# Sprint Requirements: Plan Stage

> Project: `ultraplan-go`
> Sprint: `21-plan-stage`
> Purpose: the authoritative, human-readable sprint contract. All other sprint artifacts must satisfy these requirements.

## Sprint Goal

Implement the planning plan stage so validated `reasoning.md` can be previewed, flowed, generated, and validated into `plan.md`, with `flow-state.json` updated through `plan` and no implementation, smoke, review, issue, or Git behavior.

## Required Outputs

| Output | Path | Description |
| ------ | ---- | ----------- |
| Sprint index artifact | `projects/ultraplan-go/sprints/21-plan-stage/sprint-index.md` | Selects the contracts, evidence reports, reasoning templates, and review protocols for the plan-stage implementation sprint. |
| Technical handbook artifact | `projects/ultraplan-go/sprints/21-plan-stage/technical-handbook.md` | Distills selected evidence for this sprint without making final implementation decisions. |
| Architecture area reasoning artifact | `projects/ultraplan-go/sprints/21-plan-stage/reasoning/architecture.md` | Records area decisions and trade-offs when the Architecture reasoning template is selected. |
| Final reasoning artifact | `projects/ultraplan-go/sprints/21-plan-stage/reasoning.md` | Consolidates final decisions, expected evidence, assumptions, and risks for implementing the plan stage. |
| Plan artifact | `projects/ultraplan-go/sprints/21-plan-stage/plan.md` | Defines decisions to execute, tasks, evidence checklist, risks, blockers, and success criteria for implementing this sprint. |
| Plan domain and validation | `/home/antonioborgerees/coding/ultraplan-go/internal/sprint/plan.go` | Defines plan inputs, plan validation rules, task/evidence trace checks, and diagnostics. |
| Sprint prompt rendering updates | `/home/antonioborgerees/coding/ultraplan-go/internal/sprint/prompts.go` | Adds deterministic `plan` prompt rendering with selected context, reasoning inputs, required sections, output path, and no-execution rules. |
| Sprint flow updates | `/home/antonioborgerees/coding/ultraplan-go/internal/sprint/flow.go` | Extends dry-run and runtime-backed flow execution through `plan`. |
| Sprint service updates | `/home/antonioborgerees/coding/ultraplan-go/internal/sprint/service.go` | Exposes use cases for validating, prompting, and flowing the plan stage. |
| Sprint filesystem store updates | `/home/antonioborgerees/coding/ultraplan-go/internal/sprint/store_fs.go` | Reads `reasoning.md`, selected area reasoning artifacts, and `plan.md` through workspace-safe paths. |
| Flow-state updates | `/home/antonioborgerees/coding/ultraplan-go/internal/sprint/state.go` | Records `plan` complete or failed without introducing deferred stages beyond `plan`. |
| Sprint CLI wiring | `/home/antonioborgerees/coding/ultraplan-go/internal/app/sprint_commands.go` | Adds thin command wiring for `validate plan`, `prompt plan`, and `flow --to plan`. |
| Plan tests | `/home/antonioborgerees/coding/ultraplan-go/internal/sprint/plan_test.go` | Covers plan validation, prompt rendering invariants, reasoning trace checks, and flow behavior through `plan`. |
| Sprint command tests | `/home/antonioborgerees/coding/ultraplan-go/internal/app/sprint_commands_test.go` | Covers plan-stage validate, prompt, flow, help output, exit codes, stdout/stderr separation, no ANSI output, and runtime-free validate/prompt paths. |
| Sprint review evidence | `projects/ultraplan-go/sprints/21-plan-stage/review.md` | Records implementation evidence, deviations, verification results, and review conclusions after sprint execution. |

## Acceptance Criteria

- [ ] `ultraplan sprint <project> <sprint> validate plan` validates `plan.md` without invoking runtime, OpenCode, network calls, Git commands, study execution, or artifact generation.
- [ ] `ultraplan sprint <project> <sprint> prompt plan` renders a deterministic prompt preview without invoking runtime or writing `plan.md`; if an output-preview flag exists, it may write only to the explicit preview path.
- [ ] `ultraplan sprint <project> <sprint> flow --to plan --dry-run` reports required stages, validated reasoning inputs, expected `plan.md` output, and flow-state impact without invoking runtime, creating `plan.md`, or marking `plan` complete.
- [ ] `ultraplan sprint <project> <sprint> flow --to plan` requires valid `requirements.md`, valid `sprint-index.md`, valid `technical-handbook.md`, valid selected area reasoning when templates are selected, and valid final `reasoning.md` before runtime execution for the plan stage.
- [ ] Runtime-backed plan flow invokes only the generic platform runtime boundary; `internal/sprint` must not invoke OpenCode, agentwrap adapters, shell commands, provider APIs, or Git directly.
- [ ] Runtime success alone is insufficient: the plan stage is complete only when runtime succeeds, `projects/ultraplan-go/sprints/21-plan-stage/plan.md` exists, the file passes plan validation, and `flow-state.json` is atomically updated.
- [ ] `plan.md` validation fails when the file is missing, empty, contains template placeholders, omits a citation or explicit reference to `reasoning.md`, omits decisions to execute, omits a task checklist, omits an evidence checklist, omits risks/blockers, or omits success criteria.
- [ ] `plan.md` tasks trace to final reasoning decisions and acceptance evidence; each implementation task must have a clear verification or evidence expectation.
- [ ] `plan.md` may describe implementation tasks for the current sprint, but must not invoke, automate, or model sprint implementation execution, smoke investigation, conformance review automation, issue tracking, Git mutation, or multi-stage implementation run loops as Phase 2 CLI behavior.
- [ ] Plan prompt rendering includes project slug, sprint slug, sprint path, requirements path, sprint-index path, technical-handbook path, required area reasoning paths when selected, final `reasoning.md` path, `plan.md` output path, required sections, traceability rules, and no-execution/no-mutation rules.
- [ ] Prompt rendering uses workspace-relative paths in prompt content unless an absolute path is required for local diagnostics.
- [ ] Prompt rendering instructs the runtime to produce editable Markdown only at the expected `plan.md` output path and not to modify `requirements.md`, `sprint-index.md`, `technical-handbook.md`, `reasoning/*.md`, `reasoning.md`, project docs, prior reviews, source repositories, implementation files outside the explicit runtime output, workspace config, or Git state.
- [ ] A successful flow through `plan` marks `requirements`, `sprint-index`, `technical-handbook`, selected area reasoning or skipped area reasoning, final `reasoning`, and `plan` complete when their artifacts validate.
- [ ] A failed plan flow records the `plan` stage as failed with a safe error message and does not introduce or mark later implementation, smoke, review, issue, or Git stages.
- [ ] `flow-state.json` writes remain atomic and strict as established in Sprint 17; failure to write preserves the last valid state and returns a non-zero exit.
- [ ] Flow updates only supported Planning Phase 2 stages: `requirements`, `sprint-index`, `technical-handbook`, `area-reasoning`, `reasoning`, and `plan`.
- [ ] Flow rejects unsupported stages such as `implementation`, `execute`, `smoke`, `review`, `issues`, `.run-state.json`, `smoke.md`, `smoke.json`, product-generated `review.md`, `issues.md`, and `issues.json` as current Planning Phase 2 behavior.
- [ ] Plan validation, prompt rendering, reasoning input loading, task/evidence trace checks, and flow behavior live in `internal/sprint`; no global `internal/catalog`, `internal/planning`, `internal/workflow`, `internal/stages`, `internal/validation`, `internal/reports`, or `internal/prompts` package is introduced.
- [ ] `internal/sprint` may use `internal/project` catalog data, `internal/workspace` path safety, generic platform configuration, generic platform runtime, and generic atomic file helpers, but it must not import `internal/study` or reuse study source, dimension, report, rating, summary, scheduler, or run-loop behavior.
- [ ] `internal/platform/*` packages remain product-agnostic and do not import `internal/sprint`, `internal/project`, `internal/study`, or `internal/app`.
- [ ] CLI handlers in `internal/app` only parse arguments, call sprint services, render output, and map exit codes; plan-stage business rules remain in `internal/sprint`.
- [ ] Text output is calm and scriptable, separates stdout from stderr, uses deterministic ordering, includes no ANSI escape sequences by default, and redacts sensitive values before rendering.
- [ ] Exit code `0` is returned for valid plan validation, successful prompt rendering, successful dry-run, and successful flow completion.
- [ ] Exit code `2` is returned for malformed command arguments or unsupported stages such as `execute`, `smoke`, `review`, or `issues`.
- [ ] Exit code `5` is returned for missing or invalid prerequisite artifacts, invalid `plan.md`, placeholder content, failed task/evidence trace checks, unsupported flow-state content, or failed post-generation validation.
- [ ] Runtime failures during non-dry-run flow map to the established runtime exit code and preserve safe diagnostic details.
- [ ] Tests cover plan validation for valid content, missing file, empty file, placeholders, missing `reasoning.md` citation, missing decisions, missing tasks, missing evidence checklist, missing risks/blockers, missing success criteria, untraced tasks, and forbidden deferred-stage content.
- [ ] Tests cover prompt preview content without asserting brittle prose beyond required variables, reasoning inputs, output path, required sections, traceability rules, and no-mutation/no-execution rules.
- [ ] Tests cover dry-run flow, fake-runtime successful plan generation, runtime success with missing artifact, runtime success with invalid artifact, runtime failure, validation failure after generation, invalid prerequisite reasoning, and flow-state write failure.
- [ ] Command tests prove `validate plan` and `prompt plan` do not call runtime, while non-dry-run `flow --to plan` uses a fake runtime seam.
- [ ] Offline verification passes with `go test ./...` from `/home/antonioborgerees/coding/ultraplan-go`.
- [ ] Race verification passes with `go test -race ./...` from `/home/antonioborgerees/coding/ultraplan-go`.
- [ ] CLI build passes with `go build ./cmd/ultraplan` from `/home/antonioborgerees/coding/ultraplan-go`.

## Non-Goals

- Executing implementation tasks from `plan.md`.
- Running smoke investigations or generating smoke evidence.
- Producing, validating, or automating product-generated conformance `review.md` files.
- Creating or managing local issues, issue tracker artifacts, `issues.md`, or `issues.json`.
- Mutating Git state, including `git add`, commit, push, branch creation, or automatic repository cleanup.
- Generating, validating, or requiring `.run-state.json`, `smoke.md`, `smoke.json`, product-generated `review.md`, `issues.md`, or `issues.json` as Planning Phase 2 artifacts.
- Changing the select, distill, or reason-stage contracts beyond consuming their valid artifacts as prerequisites for `plan`.
- Reading, citing, or deciding from unselected contracts, evidence reports, reasoning templates, protocols, project docs, or prior reviews as authoritative plan inputs.
- Mutating `project-index.md`, roadmap, project docs, prior sprint reviews, selected evidence reports, source repositories, workspace config, or Git state.
- Reusing study services, study prompt builders, study validators, source/dimension models, report/rating logic, summary generation, scheduler logic, or run-loop state for planning.
- Adding a generic workflow engine, global prompt package, global validation package, plugin system, hosted service, browser UI, TUI, local API server, or multi-user collaboration features.
- Requiring real OpenCode, provider credentials, network access, or real-runtime smoke tests for default verification.

## Constraints

- `internal/sprint` owns plan validation, plan prompt rendering, reasoning prerequisite checks, task/evidence trace checks, and flow-state transitions for the plan stage.
- `plan.md` must be editable Markdown and must trace executable sprint tasks to `reasoning.md`; it is not a runtime state file or an execution log.
- `reasoning.md` is the authoritative decision input for `plan.md`; `plan.md` must not invent unselected contracts, evidence, templates, protocols, project docs, or prior review conclusions.
- `project-index.md` remains a catalog, not a sprint plan, and must not be modified by the plan stage.
- Planning prompt execution must go through the generic platform runtime boundary only; product code must not directly invoke OpenCode or agentwrap adapters.
- `internal/sprint` must not depend on `internal/study` or extract study-owned behavior into shared packages for planning use.
- Flow-state schema remains strict and versioned; unsupported stages such as `implementation`, `execute`, `smoke`, `review`, and `issues` must not be accepted as supported Planning Phase 2 stages.
- Generated or previewed prompt content must avoid secrets, provider credentials, full environment dumps, unsafe raw runtime payloads, and unnecessary absolute local paths.
- Default tests must be deterministic, offline, fake-first, and free of long sleeps, network calls, provider credentials, real OpenCode execution, and Git mutation.
- Preserve prior review decisions: keep product behavior module-owned, keep CLI adapters thin, keep platform packages product-agnostic, preserve cause chains, use workspace-safe paths, keep state writes atomic, and avoid speculative abstractions.

## Dependencies

| Prior Sprint / Output | Required For | Notes |
| --------------------- | ------------ | ----- |
| Sprint 20 reason stage | Valid `reasoning.md`, selected area reasoning handling, prompt/flow command shape, and flow-state transition model through `reasoning` | Sprint 21 must consume validated reasoning and must not regenerate or reinterpret reasoning decisions outside the plan stage. |
| Sprint 19 distill stage | Valid `technical-handbook.md`, selected evidence loading precedent, and validation of evidence-only distillation | Plan generation depends on handbook context only through selected and validated planning inputs. |
| Sprint 18 select stage | Valid `sprint-index.md`, selected-context parsing, and subset validation against `project-index.md` | Sprint 21 must not reselect or silently add context. |
| Sprint 17 sprint artifact domain and flow state | Sprint resolution, artifact paths, supported planning stages, and atomic `flow-state.json` writes | Sprint 21 must preserve the exclusion of implementation, smoke, review, issue, and Git stages. |
| Sprint 16 project domain and project index | Catalog loading and catalog-only project boundary | Project catalog behavior remains read-only and project-owned. |
| Sprint 14 validation, diagnostics, and JSON stability | Exit-code discipline, redaction, deterministic rendering, and safe validation findings | Plan-stage diagnostics should follow established validation and safe-output conventions. |
| Sprint 9 runtime integration | Runtime-backed non-dry-run flow | Runtime execution must use the generic platform runtime boundary and fake-runtime tests by default. |
| `projects/ultraplan-go/project-index.md` | Authoritative catalog for selected contracts, evidence, reasoning templates, and protocols | Selected context must remain a subset of this catalog. |
| `projects/ultraplan-go/roadmap.md` | Sprint sequence and Planning Phase 2 scope | Defines Sprint 21 as the plan stage and keeps execution, smoke, review, issues, and Git mutation out of scope. |
| `projects/ultraplan-go/docs/PRD.md` | Product planning behavior and non-goals | Requires governed planning through `plan.md`, editable artifacts, selected evidence application, and deferred execution/smoke/review/issues/Git behavior. |
| `projects/ultraplan-go/docs/TRD.md` | Technical requirements for planning validators, prompt rendering, flow, and commands | Defines `plan.md` validator checks, prompt rendering, runtime completion criteria, and Phase 2 command shape. |
| `projects/ultraplan-go/docs/ARCHITECTURE.md` | Package ownership and dependency direction | Requires `internal/sprint` ownership, allows `sprint -> project/workspace/platform`, and forbids planning reuse of study semantics. |
| Prior sprint reviews | Carry-forward decisions and deferrals | Especially Sprint 20 reasoning scope, Sprint 19 fake-first verification, Sprint 17 atomic flow-state behavior, Sprint 16 catalog-only project boundary, Sprint 14 diagnostics/redaction, Sprint 9 runtime boundary, and Sprint 18 select-stage scope. |

## Review Expectations

| What | How Verified |
| ---- | ------------ |
| Required commands exist and are documented | Run `ultraplan sprint --help`, `ultraplan sprint <project> <sprint> validate --help`, `ultraplan sprint <project> <sprint> prompt --help`, and `ultraplan sprint <project> <sprint> flow --help`; inspect help for `plan` support and unsupported-stage rejection. |
| Plan validation is structural and traceable | Tests cover required sections, placeholders, `reasoning.md` citation, task-to-decision mapping, evidence checklist, risks/blockers, success criteria, and forbidden deferred-stage content. |
| Prompt preview is deterministic and runtime-free | Command tests assert stable prompt variables, workspace-relative paths, reasoning inputs, output path, required sections, traceability instructions, no-mutation rules, and no runtime invocation. |
| Flow dry-run is non-mutating | Tests assert dry-run does not create or overwrite `plan.md` and does not mark `plan` complete in `flow-state.json`. |
| Runtime-backed flow validates expected output | Fake-runtime tests cover successful generation, missing generated artifact, invalid generated artifact, runtime failure, invalid prerequisites, and validation failure after runtime success. |
| Flow-state transitions are correct | Tests assert successful flow marks stages through `plan` complete, failed flow records only the failed stage, and no unsupported later stages are introduced. |
| Unsupported future stages remain excluded | Tests and code review confirm `execute`, `implementation`, `smoke`, `review`, `issues`, `.run-state.json`, smoke artifacts, review artifacts, and issue artifacts are not modeled as current planning stages. |
| Architecture boundaries are preserved | Import review confirms `internal/sprint` does not import `internal/study`, platform packages do not import product modules, and no global planning/workflow/validation/prompt package was added. |
| CLI output and errors are safe | Command tests assert stdout/stderr separation, no ANSI escapes, deterministic diagnostics, safe redaction, and correct exit-code mapping. |
| Default verification passes | Run `go test ./...`, `go test -race ./...`, and `go build ./cmd/ultraplan` from `/home/antonioborgerees/coding/ultraplan-go`. |
