# Sprint Requirements: Select Stage

> Project: `ultraplan-go`
> Sprint: `18-select-stage`
> Purpose: the authoritative, human-readable sprint contract. All other sprint artifacts must satisfy these requirements.

## Sprint Goal

Implement the planning select stage so `sprint-index.md` can be generated, previewed, and validated as the authoritative context selection for a sprint, with subset checks against `project-index.md` and flow-state updates through `sprint-index` only.

## Required Outputs

| Output | Path | Description |
| ------ | ---- | ----------- |
| Sprint index artifact | `projects/ultraplan-go/sprints/18-select-stage/sprint-index.md` | Selects the contracts, evidence reports, reasoning templates, and review protocols used by this sprint and explicitly records excluded context. |
| Sprint index domain and parser | `/home/antonioborgerees/coding/ultraplan-go/internal/sprint/index.go` | Defines sprint-index section parsing, selected-context structures, excluded-context structures, and deterministic diagnostics. |
| Sprint index validation | `/home/antonioborgerees/coding/ultraplan-go/internal/sprint/validation.go` | Validates `sprint-index.md` required sections, placeholder absence, selected catalog entries, excluded context, and subset membership against `project-index.md`. |
| Sprint prompt rendering | `/home/antonioborgerees/coding/ultraplan-go/internal/sprint/prompts.go` | Renders deterministic prompt previews for the `sprint-index` stage from `requirements.md`, project docs, roadmap, and project catalog data. |
| Sprint flow execution | `/home/antonioborgerees/coding/ultraplan-go/internal/sprint/flow.go` | Runs or dry-runs the flow to `sprint-index`, invokes runtime only when requested, validates the expected artifact, and updates `flow-state.json` atomically. |
| Sprint service updates | `/home/antonioborgerees/coding/ultraplan-go/internal/sprint/service.go` | Exposes use cases for validating, prompting, and flowing the `sprint-index` stage while preserving existing status behavior. |
| Sprint filesystem store updates | `/home/antonioborgerees/coding/ultraplan-go/internal/sprint/store_fs.go` | Reads required stage inputs and writes generated planning artifacts through workspace-safe sprint-root paths. |
| Flow-state updates | `/home/antonioborgerees/coding/ultraplan-go/internal/sprint/state.go` | Records `sprint-index` stage completion, failure, and dry-run-safe status changes without introducing deferred stages. |
| Sprint CLI wiring | `/home/antonioborgerees/coding/ultraplan-go/internal/app/sprint_commands.go` | Adds thin command wiring for `ultraplan sprint <project> <sprint> validate sprint-index`, `prompt sprint-index`, and `flow --to sprint-index`. |
| Top-level CLI registration and help | `/home/antonioborgerees/coding/ultraplan-go/internal/app/app.go` | Updates command help and dispatch so the sprint command family exposes status, validate, prompt, and flow surfaces consistently. |
| Sprint index tests | `/home/antonioborgerees/coding/ultraplan-go/internal/sprint/sprint_index_test.go` | Covers parsing, validation, subset checks, excluded context, invalid references, missing sections, and deterministic diagnostics. |
| Sprint flow tests | `/home/antonioborgerees/coding/ultraplan-go/internal/sprint/flow_test.go` | Covers prompt rendering, dry-run behavior, fake-runtime generation, artifact validation, failure handling, and flow-state updates. |
| Sprint command tests | `/home/antonioborgerees/coding/ultraplan-go/internal/app/sprint_commands_test.go` | Covers validate, prompt, flow, help output, exit codes, stdout/stderr separation, no ANSI output, and runtime-free validate/prompt paths. |
| Sprint review evidence | `projects/ultraplan-go/sprints/18-select-stage/review.md` | Records implementation evidence, deviations, verification results, and review conclusions after sprint execution. |

## Acceptance Criteria

- [ ] `ultraplan sprint <project> <sprint> validate sprint-index` validates only the `sprint-index` stage and does not invoke runtime, OpenCode, network calls, Git commands, study execution, or artifact generation.
- [ ] `ultraplan sprint <project> <sprint> prompt sprint-index` renders a deterministic prompt preview without invoking runtime or writing `sprint-index.md`; if an output-preview flag is implemented, it may write only to the explicit preview path.
- [ ] `ultraplan sprint <project> <sprint> flow --to sprint-index --dry-run` reports the stages and inputs that would be used without invoking runtime, writing `sprint-index.md`, or marking the stage complete.
- [ ] `ultraplan sprint <project> <sprint> flow --to sprint-index` requires a present and valid `requirements.md` before runtime execution or artifact generation.
- [ ] Runtime-backed flow invokes only the generic platform runtime boundary; `internal/sprint` must not invoke OpenCode, agentwrap adapters, shell commands, or provider APIs directly.
- [ ] Runtime success alone is insufficient: the `sprint-index` stage is complete only when the runtime succeeds, `sprint-index.md` exists, validation passes, and `flow-state.json` is atomically updated.
- [ ] `sprint-index.md` validation fails when the file is missing, empty, contains template placeholders, or omits any required section: sprint scope, selected contracts, selected evidence reports, selected reasoning templates, required review protocols, and excluded context.
- [ ] Selected contracts must match entries present in the active contract pool of `project-index.md`; missing, misspelled, or path-mismatched contracts fail validation with actionable diagnostics.
- [ ] Selected evidence reports must match available evidence reports in `project-index.md`; missing, misspelled, or path-mismatched evidence fails validation with actionable diagnostics.
- [ ] Selected reasoning templates must match available reasoning templates in `project-index.md`; selecting a non-cataloged template fails validation.
- [ ] Required review protocols must match available review protocols in `project-index.md`; selecting an absent protocol fails validation.
- [ ] Empty selected reasoning templates are allowed only when the sprint explicitly records that area reasoning is not needed; this causes `area-reasoning` to remain skipped until a later stage changes the selection.
- [ ] Excluded context is explicit: the sprint index must name at least the known out-of-scope Phase 2 behavior for this sprint, including implementation execution, smoke investigation, review automation, issue tracking, Git mutation, and unrelated study-side execution behavior.
- [ ] Validation diagnostics identify the section, selected entry, referenced path or name, observed problem, and suggested repair when a subset check fails.
- [ ] Prompt rendering includes the project slug, sprint slug, sprint path, requirements path, project docs list, roadmap path, project-index path, available catalog entries, prior sprint carry-forward decisions, and instructions to select only cataloged context.
- [ ] Prompt rendering uses workspace-relative paths in prompt content unless an absolute path is required for local diagnostics.
- [ ] Prompt rendering instructs the runtime to produce editable Markdown at `projects/ultraplan-go/sprints/18-select-stage/sprint-index.md` and not to modify `project-index.md`, roadmap, project docs, prior reviews, source repositories, Git state, or implementation files.
- [ ] Flow updates only supported Planning Phase 2 stages: `requirements`, `sprint-index`, `technical-handbook`, `area-reasoning`, `reasoning`, and `plan`.
- [ ] Flow must not add, preserve as supported, or validate current-stage success for `implementation`, `execute`, `smoke`, `review`, `issues`, `.run-state.json`, `smoke.md`, `smoke.json`, `issues.md`, or `issues.json` stages.
- [ ] A successful flow to `sprint-index` marks `requirements` complete and `sprint-index` complete when both artifacts validate, marks `technical-handbook` ready, leaves later required stages missing, and keeps `area-reasoning` skipped only when no reasoning templates are selected.
- [ ] A failed flow records the `sprint-index` stage as failed with a safe error message and does not mark later stages ready or complete.
- [ ] `flow-state.json` writes remain atomic and strict as established in Sprint 17; failure to write preserves the last valid state and returns a non-zero exit.
- [ ] Sprint-index parsing and validation live in `internal/sprint`; no global `internal/catalog`, `internal/planning`, `internal/workflow`, `internal/stages`, `internal/validation`, `internal/reports`, or `internal/prompts` package is introduced.
- [ ] `internal/sprint` may use `internal/project` catalog data, `internal/workspace` path safety, generic platform configuration, and generic platform runtime, but it must not import `internal/study` or reuse study source, dimension, report, rating, summary, scheduler, or run-loop behavior.
- [ ] `internal/platform/*` packages remain product-agnostic and do not import `internal/sprint`, `internal/project`, `internal/study`, or `internal/app`.
- [ ] CLI handlers in `internal/app` only parse arguments, call sprint services, render output, and map exit codes; selection, validation, prompt, and flow business rules remain in `internal/sprint`.
- [ ] Text output is calm and scriptable, separates stdout from stderr, uses deterministic ordering, includes no ANSI escape sequences by default, and redacts sensitive values before rendering.
- [ ] Exit code `0` is returned for valid sprint-index validation, successful prompt rendering, successful dry-run, and successful flow completion.
- [ ] Exit code `2` is returned for malformed command arguments or unsupported stages such as `execute`, `smoke`, `review`, or `issues`.
- [ ] Exit code `5` is returned for sprint-index validation failures, invalid selected context, missing required sections, placeholder content, invalid flow-state content, or failed post-generation validation.
- [ ] Runtime failures during non-dry-run flow map to the established runtime exit code and preserve safe diagnostic details.
- [ ] Tests cover valid sprint-index content, missing required sections, placeholder content, absent selected contract, absent evidence report, absent reasoning template, absent review protocol, explicit no-template selection, explicit excluded context, and deterministic diagnostic ordering.
- [ ] Tests cover prompt preview content without asserting brittle prose beyond required variables, catalog entries, output path, non-goals, and no-mutation instructions.
- [ ] Tests cover dry-run flow, fake-runtime successful generation, runtime success with missing artifact, runtime success with invalid artifact, runtime failure, validation failure after generation, and flow-state write failure.
- [ ] Command tests prove `validate sprint-index` and `prompt sprint-index` do not call runtime, while non-dry-run `flow --to sprint-index` uses a fake runtime seam.
- [ ] Offline verification passes with `go test ./...` from `/home/antonioborgerees/coding/ultraplan-go`.
- [ ] Race verification passes with `go test -race ./...` from `/home/antonioborgerees/coding/ultraplan-go`.
- [ ] CLI build passes with `go build ./cmd/ultraplan` from `/home/antonioborgerees/coding/ultraplan-go`.

## Non-Goals

- Generating or validating `technical-handbook.md`; that belongs to Sprint 19.
- Generating or validating `reasoning/*.md` or final `reasoning.md`; that belongs to Sprint 20.
- Generating or validating `plan.md`; that belongs to Sprint 21.
- Implementing `flow --to technical-handbook`, `flow --to area-reasoning`, `flow --to reasoning`, or `flow --to plan` beyond rejecting unsupported or out-of-scope stage requests when not implemented.
- Implementing sprint implementation execution, smoke investigation execution, conformance review automation, issue tracking, project-management behavior, or Git mutation.
- Generating, validating, or requiring `.run-state.json`, `smoke.md`, `smoke.json`, `review.md`, `issues.md`, or `issues.json` as Planning Phase 2 artifacts.
- Mutating `project-index.md`, `roadmap.md`, project docs, prior sprint reviews, source repositories, workspace config, or Git state.
- Reusing study services, study prompt builders, study validators, source/dimension models, report/rating logic, summary generation, scheduler logic, or run-loop state for planning selection.
- Adding a generic workflow engine, global prompt package, global validation package, plugin system, hosted service, browser UI, TUI, local API server, or multi-user collaboration features.
- Adding stable JSON output for sprint validate, prompt, or flow unless it is already supported by existing app infrastructure without expanding scope.
- Requiring real OpenCode, provider credentials, network access, or real-runtime smoke tests for default verification.

## Constraints

- `internal/sprint` owns planning artifact validation, prompt rendering, and flow state for the select stage.
- `sprint-index.md` must be treated as editable Markdown and as the authoritative selected context for later planning stages.
- All selected context must be validated as a subset of entries cataloged in `project-index.md`; do not infer or silently add uncataloged contracts, evidence, templates, or protocols.
- `project-index.md` remains a catalog, not a sprint plan, and must not be modified by the select stage.
- Planning prompt execution must go through the generic platform runtime boundary only; product code must not directly invoke OpenCode or agentwrap adapters.
- `internal/sprint` must not depend on `internal/study` or extract study-owned behavior into shared packages for planning use.
- Flow-state schema remains strict and versioned; unsupported stages such as `implementation` and `review` must not be accepted as supported Planning Phase 2 stages.
- Generated or previewed prompt content must avoid secrets, provider credentials, full environment dumps, and unsafe raw runtime payloads.
- Default tests must be deterministic, offline, fake-first, and free of long sleeps, network calls, provider credentials, real OpenCode execution, and Git mutation.
- Preserve prior review decisions: keep product behavior module-owned, keep CLI adapters thin, keep platform packages product-agnostic, preserve cause chains, use workspace-safe paths, and avoid speculative abstractions.

## Dependencies

| Prior Sprint / Output | Required For | Notes |
| --------------------- | ------------ | ----- |
| Sprint 16 project domain and project index | Catalog loading and subset validation | Sprint 18 must use the accepted project catalog boundary rather than reparsing project roots ad hoc in app code. |
| Sprint 17 sprint artifact domain and flow state | Sprint resolution, artifact paths, supported planning stages, and atomic `flow-state.json` writes | Sprint 18 extends the existing `internal/sprint` package with select-stage behavior and must preserve the exclusion of implementation, smoke, review, issue, and Git stages. |
| Sprint 14 validation, diagnostics, and JSON stability | Exit-code discipline, redaction, and deterministic rendering | Select-stage diagnostics should follow established validation and safe-output conventions. |
| Sprint 9 runtime integration | Runtime-backed non-dry-run flow | Runtime execution must use the generic platform runtime boundary and fake-runtime tests by default. |
| `projects/ultraplan-go/project-index.md` | Authoritative catalog for subset validation | Selected contracts, evidence reports, reasoning templates, and review protocols must be present here. |
| `projects/ultraplan-go/roadmap.md` | Sprint sequence and Planning Phase 2 scope | Defines Sprint 18 as the select stage and keeps later distill, reason, and plan stages separate. |
| `projects/ultraplan-go/docs/PRD.md` | Product planning behavior and non-goals | Requires governed planning through `plan.md` and defers execution, smoke, review automation, issues, and Git mutation. |
| `projects/ultraplan-go/docs/TRD.md` | Technical requirements for planning validators, prompt rendering, flow, and commands | Defines `sprint-index.md` validation, prompt rendering, runtime completion criteria, and required Phase 2 command shape. |
| `projects/ultraplan-go/docs/ARCHITECTURE.md` | Package ownership and dependency direction | Requires `internal/sprint` ownership, allows `sprint -> project/workspace/platform`, and forbids planning reuse of study semantics. |
| Prior sprint reviews | Carry-forward decisions and deferrals | Especially Sprint 17's deferral of semantic sprint-index validation, Sprint 16's catalog-only project boundary, and prior runtime/durable-state redaction and atomic-write decisions. |

## Review Expectations

| What | How Verified |
| ---- | ------------ |
| Required commands exist and are documented | Run `ultraplan sprint --help`, `ultraplan sprint <project> <sprint> validate --help`, `ultraplan sprint <project> <sprint> prompt --help`, and `ultraplan sprint <project> <sprint> flow --help`; inspect help for supported stages and non-goal stage rejection. |
| Sprint-index validation is semantic, not just existence-based | Tests cover required sections, placeholders, selected contracts, evidence, reasoning templates, review protocols, excluded context, and missing catalog entries. |
| Selected context is a subset of `project-index.md` | Fixture tests include valid catalog selections and invalid names/paths for each catalog category. |
| Prompt preview is deterministic and runtime-free | Command tests assert stable prompt variables, workspace-relative paths, catalog entries, output path, non-goal instructions, and no runtime invocation. |
| Flow dry-run is non-mutating | Tests assert dry-run does not create or overwrite `sprint-index.md` and does not mark `sprint-index` complete in `flow-state.json`. |
| Runtime-backed flow validates expected output | Fake-runtime tests cover successful generation, missing generated artifact, invalid generated artifact, runtime failure, and validation failure after runtime success. |
| Flow-state transitions are correct | Tests assert successful flow marks `requirements` and `sprint-index` complete, readies `technical-handbook`, and does not introduce unsupported stages. |
| Unsupported future stages remain excluded | Tests and code review confirm `execute`, `implementation`, `smoke`, `review`, `issues`, `.run-state.json`, `smoke.md`, `review.md`, and issue artifacts are not modeled as current planning stages. |
| Architecture boundaries are preserved | Import review confirms `internal/sprint` does not import `internal/study`, platform packages do not import product modules, and no global planning/workflow/validation/prompt package was added. |
| CLI output and errors are safe | Command tests assert stdout/stderr separation, no ANSI escapes, deterministic diagnostics, safe redaction, and correct exit-code mapping. |
| Default verification passes | Run `go test ./...`, `go test -race ./...`, and `go build ./cmd/ultraplan` from `/home/antonioborgerees/coding/ultraplan-go`. |
