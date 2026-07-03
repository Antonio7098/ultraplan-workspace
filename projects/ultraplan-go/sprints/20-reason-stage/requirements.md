# Sprint Requirements: Reason Stage

> Project: `ultraplan-go`
> Sprint: `20-reason-stage`
> Purpose: the authoritative, human-readable sprint contract. All other sprint artifacts must satisfy these requirements.

## Sprint Goal

Implement the planning reason stage so selected area reasoning and final `reasoning.md` can be generated, previewed, flowed, and validated from valid requirements, `sprint-index.md`, and `technical-handbook.md`, with correct `flow-state.json` updates through `reasoning` and no plan or execution behavior.

## Required Outputs

| Output | Path | Description |
| ------ | ---- | ----------- |
| Sprint index artifact | `projects/ultraplan-go/sprints/20-reason-stage/sprint-index.md` | Selects the contracts, evidence reports, reasoning templates, and review protocols for the reason-stage implementation sprint. |
| Technical handbook artifact | `projects/ultraplan-go/sprints/20-reason-stage/technical-handbook.md` | Distills selected evidence for this sprint without making final implementation decisions. |
| Architecture area reasoning artifact | `projects/ultraplan-go/sprints/20-reason-stage/reasoning/architecture.md` | Records area decisions and trade-offs when the Architecture reasoning template is selected. |
| Final reasoning artifact | `projects/ultraplan-go/sprints/20-reason-stage/reasoning.md` | Consolidates final decisions, expected evidence, assumptions, and risks for implementing the reason stage. |
| Reasoning domain and validation | `/home/antonioborgerees/coding/ultraplan-go/internal/sprint/reasoning.go` | Defines selected-template detection, area reasoning manifests, final reasoning inputs, validation rules, and diagnostics. |
| Sprint prompt rendering updates | `/home/antonioborgerees/coding/ultraplan-go/internal/sprint/prompts.go` | Adds deterministic prompt rendering for `area-reasoning` and `reasoning` stages with selected context, handbook content, output paths, and no-mutation rules. |
| Sprint flow updates | `/home/antonioborgerees/coding/ultraplan-go/internal/sprint/flow.go` | Extends dry-run and runtime-backed flow execution through `area-reasoning` and final `reasoning`. |
| Sprint service updates | `/home/antonioborgerees/coding/ultraplan-go/internal/sprint/service.go` | Exposes use cases for validating, prompting, and flowing area reasoning and final reasoning stages. |
| Sprint filesystem store updates | `/home/antonioborgerees/coding/ultraplan-go/internal/sprint/store_fs.go` | Reads selected templates, requirements, sprint index, handbook, and reasoning artifacts through workspace-safe paths. |
| Flow-state updates | `/home/antonioborgerees/coding/ultraplan-go/internal/sprint/state.go` | Records area-reasoning skipped/complete/failed states and final reasoning completion without introducing deferred stages. |
| Sprint CLI wiring | `/home/antonioborgerees/coding/ultraplan-go/internal/app/sprint_commands.go` | Adds thin command wiring for `validate`, `prompt`, and `flow` support for `area-reasoning` and `reasoning`. |
| Reasoning tests | `/home/antonioborgerees/coding/ultraplan-go/internal/sprint/reasoning_test.go` | Covers selected-template detection, area/final validation, selected-context constraints, and deterministic diagnostics. |
| Sprint flow tests | `/home/antonioborgerees/coding/ultraplan-go/internal/sprint/flow_test.go` | Covers dry-run, fake-runtime generation, validation gates, skipped area reasoning, failure handling, and flow-state transitions through `reasoning`. |
| Sprint command tests | `/home/antonioborgerees/coding/ultraplan-go/internal/app/sprint_commands_test.go` | Covers validate, prompt, flow, help output, exit codes, stdout/stderr separation, no ANSI output, and runtime-free validate/prompt paths. |
| Sprint review evidence | `projects/ultraplan-go/sprints/20-reason-stage/review.md` | Records implementation evidence, deviations, verification results, and review conclusions after sprint execution. |

## Acceptance Criteria

- [ ] `ultraplan sprint <project> <sprint> validate area-reasoning` validates selected area reasoning files without invoking runtime, OpenCode, network calls, Git commands, study execution, or artifact generation.
- [ ] `ultraplan sprint <project> <sprint> validate reasoning` validates final `reasoning.md` and required area reasoning prerequisites without invoking runtime or mutating artifacts.
- [ ] `ultraplan sprint <project> <sprint> prompt area-reasoning` renders deterministic prompt previews for selected reasoning templates without invoking runtime or writing `reasoning/*.md`; if an output-preview flag exists, it may write only to the explicit preview path.
- [ ] `ultraplan sprint <project> <sprint> prompt reasoning` renders a deterministic final-reasoning prompt preview without invoking runtime or writing `reasoning.md`; if an output-preview flag exists, it may write only to the explicit preview path.
- [ ] `ultraplan sprint <project> <sprint> flow --to reasoning --dry-run` reports the stages, selected templates, expected area reasoning outputs, and final reasoning output without invoking runtime, creating reasoning artifacts, or marking reasoning complete.
- [ ] `ultraplan sprint <project> <sprint> flow --to reasoning` requires valid `requirements.md`, valid `sprint-index.md`, and valid `technical-handbook.md` before runtime execution for reasoning stages.
- [ ] Area reasoning is skipped only when the validated `sprint-index.md` selects no reasoning templates; when the Architecture template is selected, `projects/ultraplan-go/sprints/20-reason-stage/reasoning/architecture.md` is required and must validate before final reasoning completes.
- [ ] Runtime-backed flow invokes only the generic platform runtime boundary; `internal/sprint` must not invoke OpenCode, agentwrap adapters, shell commands, provider APIs, or Git directly.
- [ ] Runtime success alone is insufficient: a reasoning stage is complete only when runtime succeeds, the expected artifact exists, that artifact passes its validator, and `flow-state.json` is atomically updated.
- [ ] Area-reasoning prompt rendering includes project slug, sprint slug, sprint path, requirements path, sprint-index path, technical-handbook path, selected reasoning template path, output path, selected contracts/evidence summary, required sections, and selected-context-only rules.
- [ ] Final-reasoning prompt rendering includes project slug, sprint slug, sprint path, requirements path, sprint-index path, technical-handbook path, required area reasoning artifact paths when selected, final `reasoning.md` output path, required sections, and no-plan/no-execution rules.
- [ ] Prompt rendering uses workspace-relative paths in prompt content unless an absolute path is required for local diagnostics.
- [ ] Prompt rendering instructs the runtime to produce editable Markdown only at the expected reasoning output path and not to modify `requirements.md`, `sprint-index.md`, `technical-handbook.md`, `plan.md`, project docs, prior reviews, source repositories, implementation files outside explicit runtime output, workspace config, or Git state.
- [ ] Selected reasoning templates must be present in `project-index.md`, selected by `sprint-index.md`, readable, and workspace-safe unless explicitly cataloged as external; missing, unreadable, unselected, or escaping paths fail with actionable diagnostics.
- [ ] Area reasoning validation fails when a selected area file is missing, empty, contains template placeholders, omits area decisions, omits trade-offs, or references unselected contracts, evidence, templates, or protocols as decision sources.
- [ ] Final reasoning validation fails when `reasoning.md` is missing, empty, contains template placeholders, omits final decisions, omits expected evidence, omits assumptions and risks, or omits required selected-area reasoning references when templates are selected.
- [ ] Final reasoning includes decisions, expected evidence, assumptions, and risks, and it traces decisions to selected requirements, selected contracts, selected evidence, `technical-handbook.md`, and selected area reasoning where applicable.
- [ ] Reasoning validation does not require `plan.md`, task checklists, implementation steps, smoke evidence, review evidence, issue artifacts, or Git state.
- [ ] A successful flow through `reasoning` marks `requirements`, `sprint-index`, `technical-handbook`, selected area reasoning, and final `reasoning` complete when their artifacts validate; it may mark `area-reasoning` skipped only when no templates are selected.
- [ ] A successful flow through `reasoning` must not mark `plan` complete and must not create or validate execution, smoke, review, issue, or Git stages.
- [ ] A failed area-reasoning or final-reasoning flow records the failed stage with a safe error message and does not mark later stages complete.
- [ ] `flow-state.json` writes remain atomic and strict as established in Sprint 17; failure to write preserves the last valid state and returns a non-zero exit.
- [ ] Flow updates only supported Planning Phase 2 stages: `requirements`, `sprint-index`, `technical-handbook`, `area-reasoning`, `reasoning`, and `plan`.
- [ ] Flow must reject unsupported stages such as `implementation`, `execute`, `smoke`, `review`, `issues`, `.run-state.json`, `smoke.md`, `smoke.json`, `review.md`, `issues.md`, and `issues.json` as current Planning Phase 2 behavior.
- [ ] Reasoning validation, prompt rendering, selected-template detection, selected-context loading, and flow behavior live in `internal/sprint`; no global `internal/catalog`, `internal/planning`, `internal/workflow`, `internal/stages`, `internal/validation`, `internal/reports`, or `internal/prompts` package is introduced.
- [ ] `internal/sprint` may use `internal/project` catalog data, `internal/workspace` path safety, generic platform configuration, generic platform runtime, and generic atomic file helpers, but it must not import `internal/study` or reuse study source, dimension, report, rating, summary, scheduler, or run-loop behavior.
- [ ] `internal/platform/*` packages remain product-agnostic and do not import `internal/sprint`, `internal/project`, `internal/study`, or `internal/app`.
- [ ] CLI handlers in `internal/app` only parse arguments, call sprint services, render output, and map exit codes; reasoning business rules remain in `internal/sprint`.
- [ ] Text output is calm and scriptable, separates stdout from stderr, uses deterministic ordering, includes no ANSI escape sequences by default, and redacts sensitive values before rendering.
- [ ] Exit code `0` is returned for valid reasoning validation, successful prompt rendering, successful dry-run, and successful flow completion.
- [ ] Exit code `2` is returned for malformed command arguments or unsupported stages such as `execute`, `smoke`, `review`, or `issues`.
- [ ] Exit code `5` is returned for missing or invalid `requirements.md`, invalid `sprint-index.md`, invalid `technical-handbook.md`, invalid selected templates, invalid area reasoning, invalid final reasoning, placeholder content, unsupported flow-state content, or failed post-generation validation.
- [ ] Runtime failures during non-dry-run flow map to the established runtime exit code and preserve safe diagnostic details.
- [ ] Tests cover selected-template detection for no templates, one template, missing template paths, unreadable template paths, unselected template references, and deterministic template ordering.
- [ ] Tests cover area reasoning validation for valid content, missing file, empty file, placeholders, missing area decisions, missing trade-offs, unselected-context references, and workspace path escapes.
- [ ] Tests cover final reasoning validation for valid content, missing file, empty file, placeholders, missing final decisions, missing expected evidence, missing assumptions/risks, missing required area references, and selected-area success.
- [ ] Tests cover prompt preview content without asserting brittle prose beyond required variables, selected template manifest, selected context summary, output path, required sections, selected-only rules, and no-mutation rules.
- [ ] Tests cover dry-run flow, fake-runtime successful area generation, fake-runtime successful final generation, runtime success with missing artifact, runtime success with invalid artifact, runtime failure, validation failure after generation, no-template area skip, selected-template failure before runtime, and flow-state write failure.
- [ ] Command tests prove `validate area-reasoning`, `validate reasoning`, `prompt area-reasoning`, and `prompt reasoning` do not call runtime, while non-dry-run `flow --to reasoning` uses a fake runtime seam.
- [ ] Offline verification passes with `go test ./...` from `/home/antonioborgerees/coding/ultraplan-go`.
- [ ] Race verification passes with `go test -race ./...` from `/home/antonioborgerees/coding/ultraplan-go`.
- [ ] CLI build passes with `go build ./cmd/ultraplan` from `/home/antonioborgerees/coding/ultraplan-go`.

## Non-Goals

- Generating or validating `plan.md`; that belongs to Sprint 21.
- Creating task breakdowns, implementation checklists, evidence checklists for execution, or plan-stage success criteria in `reasoning.md`.
- Implementing sprint implementation execution, smoke investigation execution, conformance review automation, issue tracking, project-management behavior, or Git mutation.
- Generating, validating, or requiring `.run-state.json`, `smoke.md`, `smoke.json`, product-generated `review.md`, `issues.md`, or `issues.json` as Planning Phase 2 artifacts.
- Changing the select-stage or distill-stage contract beyond what is needed to consume valid `sprint-index.md` and `technical-handbook.md` as prerequisites.
- Reading, citing, or deciding from unselected contracts, evidence reports, reasoning templates, protocols, project docs, or prior reviews as authoritative reasoning inputs.
- Mutating `project-index.md`, roadmap, project docs, prior sprint reviews, source repositories, selected evidence reports, workspace config, or Git state.
- Reusing study services, study prompt builders, study validators, source/dimension models, report/rating logic, summary generation, scheduler logic, or run-loop state for planning reasoning.
- Adding a generic workflow engine, global prompt package, global validation package, plugin system, hosted service, browser UI, TUI, local API server, or multi-user collaboration features.
- Adding stable JSON output for sprint validate, prompt, or flow unless it is already supported by existing app infrastructure without expanding scope.
- Requiring real OpenCode, provider credentials, network access, or real-runtime smoke tests for default verification.

## Constraints

- `internal/sprint` owns selected-template detection, area reasoning validation, final reasoning validation, reasoning prompt rendering, and flow-state transitions for the reason stage.
- Reasoning artifacts must be editable Markdown and must remain selected-context decisions, not implementation plans or execution records.
- `reasoning/*.md` may be skipped only when `sprint-index.md` selects no reasoning templates; otherwise every selected template must map to a deterministic required area reasoning output.
- Final `reasoning.md` must not invent unselected contracts, evidence, templates, protocols, project docs, or prior review conclusions.
- `project-index.md` remains a catalog, not a sprint plan, and must not be modified by the reason stage.
- Planning prompt execution must go through the generic platform runtime boundary only; product code must not directly invoke OpenCode or agentwrap adapters.
- `internal/sprint` must not depend on `internal/study` or extract study-owned behavior into shared packages for planning use.
- Flow-state schema remains strict and versioned; unsupported stages such as `implementation`, `execute`, `smoke`, `review`, and `issues` must not be accepted as supported Planning Phase 2 stages.
- Generated or previewed prompt content must avoid secrets, provider credentials, full environment dumps, unsafe raw runtime payloads, and unnecessary absolute local paths.
- Default tests must be deterministic, offline, fake-first, and free of long sleeps, network calls, provider credentials, real OpenCode execution, and Git mutation.
- Preserve prior review decisions: keep product behavior module-owned, keep CLI adapters thin, keep platform packages product-agnostic, preserve cause chains, use workspace-safe paths, keep state writes atomic, and avoid speculative abstractions.

## Dependencies

| Prior Sprint / Output | Required For | Notes |
| --------------------- | ------------ | ----- |
| Sprint 19 distill stage | Valid `technical-handbook.md`, selected evidence loading precedent, prompt/flow command shape, and flow-state transition model through `technical-handbook` | Sprint 20 must consume a valid handbook and preserve handbook as evidence distillation rather than decision source of record. |
| Sprint 18 select stage | Valid `sprint-index.md`, selected-context parsing, selected reasoning templates, and subset validation against `project-index.md` | Sprint 20 must not reselect or silently add context. |
| Sprint 17 sprint artifact domain and flow state | Sprint resolution, artifact paths, supported planning stages, and atomic `flow-state.json` writes | Sprint 20 must preserve the exclusion of implementation, smoke, review, issue, and Git stages. |
| Sprint 16 project domain and project index | Catalog loading and selected template/path validation | Reasoning templates and selected context must validate against the accepted project catalog boundary. |
| Sprint 14 validation, diagnostics, and JSON stability | Exit-code discipline, redaction, deterministic rendering, and safe validation findings | Reason-stage diagnostics should follow established validation and safe-output conventions. |
| Sprint 9 runtime integration | Runtime-backed non-dry-run flow | Runtime execution must use the generic platform runtime boundary and fake-runtime tests by default. |
| `projects/ultraplan-go/project-index.md` | Authoritative catalog for selected contracts, evidence, reasoning templates, and protocols | Selected reasoning templates must be present here before area reasoning can run. |
| `projects/ultraplan-go/roadmap.md` | Sprint sequence and Planning Phase 2 scope | Defines Sprint 20 as the reason stage and keeps plan, execution, smoke, review, issues, and Git mutation out of scope. |
| `projects/ultraplan-go/docs/PRD.md` | Product planning behavior and non-goals | Requires governed planning through `plan.md`, editable artifacts, selected evidence application, and deferred execution/smoke/review/issues/Git behavior. |
| `projects/ultraplan-go/docs/TRD.md` | Technical requirements for planning validators, prompt rendering, flow, and commands | Defines `reasoning/*.md` and `reasoning.md` validator checks, prompt rendering, runtime completion criteria, and Phase 2 command shape. |
| `projects/ultraplan-go/docs/ARCHITECTURE.md` | Package ownership and dependency direction | Requires `internal/sprint` ownership, allows `sprint -> project/workspace/platform`, and forbids planning reuse of study semantics. |
| Prior sprint reviews | Carry-forward decisions and deferrals | Especially Sprint 19 selected-evidence and fake-first verification, Sprint 17 atomic flow-state behavior, Sprint 16 catalog-only project boundary, Sprint 14 diagnostics/redaction, Sprint 9 runtime boundary, and Sprint 18 select-stage scope. |

## Review Expectations

| What | How Verified |
| ---- | ------------ |
| Required commands exist and are documented | Run `ultraplan sprint --help`, `ultraplan sprint <project> <sprint> validate --help`, `ultraplan sprint <project> <sprint> prompt --help`, and `ultraplan sprint <project> <sprint> flow --help`; inspect help for `area-reasoning` and `reasoning` support and unsupported-stage rejection. |
| Selected-template detection is constrained | Fixture tests prove only `sprint-index.md` selected templates from `project-index.md` are used, with missing, unreadable, unselected, or escaping paths failing before runtime. |
| Area reasoning validation is semantic, not just existence-based | Tests cover required sections, placeholders, selected context, missing/unreadable files, and unselected-context references. |
| Final reasoning validation is semantic, not just existence-based | Tests cover final decisions, expected evidence, assumptions and risks, selected area references, placeholders, and no plan/task checklist requirements. |
| Prompt previews are deterministic and runtime-free | Command tests assert stable prompt variables, workspace-relative paths, selected template manifest, selected context summary, output path, required sections, no-mutation instructions, and no runtime invocation. |
| Flow dry-run is non-mutating | Tests assert dry-run does not create or overwrite reasoning artifacts and does not mark reasoning complete in `flow-state.json`. |
| Runtime-backed flow validates expected outputs | Fake-runtime tests cover successful generation, missing generated artifact, invalid generated artifact, runtime failure, selected-template failure, and validation failure after runtime success. |
| Flow-state transitions are correct | Tests assert successful flow marks prior prerequisites and reasoning stages complete or skipped as appropriate, does not mark plan complete, and does not introduce unsupported stages. |
| Unsupported future stages remain excluded | Tests and code review confirm `execute`, `implementation`, `smoke`, `review`, `issues`, `.run-state.json`, `smoke.md`, `review.md`, and issue artifacts are not modeled as current planning stages. |
| Architecture boundaries are preserved | Import review confirms `internal/sprint` does not import `internal/study`, platform packages do not import product modules, and no global planning/workflow/validation/prompt package was added. |
| CLI output and errors are safe | Command tests assert stdout/stderr separation, no ANSI escapes, deterministic diagnostics, safe redaction, and correct exit-code mapping. |
| Default verification passes | Run `go test ./...`, `go test -race ./...`, and `go build ./cmd/ultraplan` from `/home/antonioborgerees/coding/ultraplan-go`. |
