# Sprint Requirements: Distill Stage

> Project: `ultraplan-go`
> Sprint: `19-distill-stage`
> Purpose: the authoritative, human-readable sprint contract. All other sprint artifacts must satisfy these requirements.

## Sprint Goal

Implement the planning distill stage so `technical-handbook.md` can be generated, previewed, flowed, and validated from valid `sprint-index.md` selected evidence only, with flow-state updates through `technical-handbook` and without making implementation decisions.

## Required Outputs

| Output | Path | Description |
| ------ | ---- | ----------- |
| Sprint index artifact | `projects/ultraplan-go/sprints/19-distill-stage/sprint-index.md` | Selects the contracts, evidence reports, reasoning templates, and review protocols for this sprint, including the evidence that the distill stage must load. |
| Technical handbook artifact | `projects/ultraplan-go/sprints/19-distill-stage/technical-handbook.md` | Distills selected evidence into relevant patterns, trade-offs, anti-patterns, open questions, and evidence pointers without implementation decisions. |
| Technical handbook domain and validation | `/home/antonioborgerees/coding/ultraplan-go/internal/sprint/handbook.go` | Defines technical-handbook input manifests, selected-evidence loading, required-section validation, evidence-trace validation, and no-decision checks. |
| Sprint prompt rendering updates | `/home/antonioborgerees/coding/ultraplan-go/internal/sprint/prompts.go` | Adds deterministic prompt rendering for the `technical-handbook` stage from `requirements.md`, validated `sprint-index.md`, selected evidence reports, and output instructions. |
| Sprint flow updates | `/home/antonioborgerees/coding/ultraplan-go/internal/sprint/flow.go` | Extends flow execution through `technical-handbook`, including dry-run behavior, runtime-backed generation, post-generation validation, and failure handling. |
| Sprint service updates | `/home/antonioborgerees/coding/ultraplan-go/internal/sprint/service.go` | Exposes use cases for validating, prompting, and flowing the `technical-handbook` stage while preserving existing status and select-stage behavior. |
| Sprint filesystem store updates | `/home/antonioborgerees/coding/ultraplan-go/internal/sprint/store_fs.go` | Reads selected evidence reports through workspace-safe paths and writes generated planning artifacts only under the selected sprint root. |
| Flow-state updates | `/home/antonioborgerees/coding/ultraplan-go/internal/sprint/state.go` | Records `technical-handbook` completion, failure, and next-stage readiness without introducing deferred implementation, smoke, review, issue, or Git stages. |
| Sprint CLI wiring | `/home/antonioborgerees/coding/ultraplan-go/internal/app/sprint_commands.go` | Adds thin command wiring for `ultraplan sprint <project> <sprint> validate technical-handbook`, `prompt technical-handbook`, and `flow --to technical-handbook`. |
| Top-level CLI registration and help | `/home/antonioborgerees/coding/ultraplan-go/internal/app/app.go` | Updates command help and dispatch so the sprint command family exposes the distill stage consistently. |
| Technical handbook tests | `/home/antonioborgerees/coding/ultraplan-go/internal/sprint/handbook_test.go` | Covers selected-evidence loading, missing/unreadable evidence, required sections, evidence traces, placeholder rejection, no-decision checks, and deterministic diagnostics. |
| Sprint flow tests | `/home/antonioborgerees/coding/ultraplan-go/internal/sprint/flow_test.go` | Covers dry-run behavior, fake-runtime generation, artifact validation, evidence-only constraints, flow-state transitions, failure states, and write failures through `technical-handbook`. |
| Sprint command tests | `/home/antonioborgerees/coding/ultraplan-go/internal/app/sprint_commands_test.go` | Covers validate, prompt, flow, help output, exit codes, stdout/stderr separation, no ANSI output, and runtime-free validate/prompt paths for `technical-handbook`. |
| Sprint review evidence | `projects/ultraplan-go/sprints/19-distill-stage/review.md` | Records implementation evidence, deviations, verification results, and review conclusions after sprint execution. |

## Acceptance Criteria

- [ ] `ultraplan sprint <project> <sprint> validate technical-handbook` validates only the technical-handbook stage and does not invoke runtime, OpenCode, network calls, Git commands, study execution, or artifact generation.
- [ ] `ultraplan sprint <project> <sprint> prompt technical-handbook` renders a deterministic prompt preview without invoking runtime or writing `technical-handbook.md`; if an output-preview flag exists, it may write only to the explicit preview path.
- [ ] `ultraplan sprint <project> <sprint> flow --to technical-handbook --dry-run` reports the stages and selected evidence that would be used without invoking runtime, writing `technical-handbook.md`, or marking the stage complete.
- [ ] `ultraplan sprint <project> <sprint> flow --to technical-handbook` requires valid `requirements.md` and valid `sprint-index.md` before runtime execution or artifact generation.
- [ ] Runtime-backed flow invokes only the generic platform runtime boundary; `internal/sprint` must not invoke OpenCode, agentwrap adapters, shell commands, provider APIs, or Git directly.
- [ ] Runtime success alone is insufficient: the `technical-handbook` stage is complete only when runtime succeeds, `technical-handbook.md` exists, handbook validation passes, and `flow-state.json` is atomically updated.
- [ ] Technical-handbook prompt rendering includes project slug, sprint slug, sprint path, requirements path, sprint-index path, technical-handbook output path, selected evidence manifest, required handbook sections, no-decision instructions, and selected-evidence-only constraints.
- [ ] Prompt rendering uses workspace-relative paths in prompt content unless an absolute path is required for local diagnostics.
- [ ] Prompt rendering instructs the runtime to produce editable Markdown at `projects/ultraplan-go/sprints/19-distill-stage/technical-handbook.md` and not to modify `project-index.md`, roadmap, project docs, prior reviews, source repositories, Git state, implementation files, `sprint-index.md`, reasoning artifacts, or `plan.md`.
- [ ] Selected evidence loading reads only evidence reports selected by `sprint-index.md`, plus the minimum project/sprint artifacts required to identify and validate those selections.
- [ ] Selected evidence paths must be present in `project-index.md`, selected by `sprint-index.md`, readable, and workspace-safe unless explicitly cataloged as external; missing, unreadable, unselected, or escaping paths fail with actionable diagnostics.
- [ ] `technical-handbook.md` validation fails when the file is missing, empty, contains template placeholders, or omits required sections: selected studies/reports, relevant patterns, trade-offs, anti-patterns, open questions, and evidence pointers.
- [ ] Handbook content must trace to selected evidence reports by report name or selected report path in the selected studies/reports or evidence pointers section.
- [ ] Handbook validation fails when it cites uncataloged or unselected evidence report names or paths as evidence sources.
- [ ] Handbook validation fails when it makes implementation decisions, prescribes final architecture choices, creates task plans, or uses decision-language sections that belong to `reasoning.md` or `plan.md`.
- [ ] Handbook validation may allow observations, patterns, warnings, anti-patterns, and open questions when they are framed as evidence-backed guidance rather than final sprint decisions.
- [ ] A successful flow through `technical-handbook` marks `requirements`, `sprint-index`, and `technical-handbook` complete when their artifacts validate, readies the next applicable reasoning stage, and does not mark `reasoning.md` or `plan.md` complete.
- [ ] `area-reasoning` may be marked skipped only when the validated `sprint-index.md` selects no reasoning templates; otherwise the next applicable area-reasoning stage must remain missing or ready, not complete.
- [ ] A failed flow records the `technical-handbook` stage as failed with a safe error message and does not mark later reasoning or plan stages ready or complete.
- [ ] `flow-state.json` writes remain atomic and strict as established in Sprint 17; failure to write preserves the last valid state and returns a non-zero exit.
- [ ] Flow updates only supported Planning Phase 2 stages: `requirements`, `sprint-index`, `technical-handbook`, `area-reasoning`, `reasoning`, and `plan`.
- [ ] Flow must not add, preserve as supported, or validate current-stage success for `implementation`, `execute`, `smoke`, `review`, `issues`, `.run-state.json`, `smoke.md`, `smoke.json`, `review.md`, `issues.md`, or `issues.json` stages.
- [ ] Technical-handbook validation, prompt rendering, selected-evidence loading, and flow behavior live in `internal/sprint`; no global `internal/catalog`, `internal/planning`, `internal/workflow`, `internal/stages`, `internal/validation`, `internal/reports`, or `internal/prompts` package is introduced.
- [ ] `internal/sprint` may use `internal/project` catalog data, `internal/workspace` path safety, generic platform configuration, generic platform runtime, and generic atomic file helpers, but it must not import `internal/study` or reuse study source, dimension, report, rating, summary, scheduler, or run-loop behavior.
- [ ] `internal/platform/*` packages remain product-agnostic and do not import `internal/sprint`, `internal/project`, `internal/study`, or `internal/app`.
- [ ] CLI handlers in `internal/app` only parse arguments, call sprint services, render output, and map exit codes; handbook, prompt, evidence, validation, and flow business rules remain in `internal/sprint`.
- [ ] Text output is calm and scriptable, separates stdout from stderr, uses deterministic ordering, includes no ANSI escape sequences by default, and redacts sensitive values before rendering.
- [ ] Exit code `0` is returned for valid handbook validation, successful prompt rendering, successful dry-run, and successful flow completion.
- [ ] Exit code `2` is returned for malformed command arguments or unsupported stages such as `execute`, `smoke`, `review`, or `issues`.
- [ ] Exit code `5` is returned for missing or invalid `requirements.md`, invalid `sprint-index.md`, technical-handbook validation failures, invalid selected evidence, missing required sections, placeholder content, invalid flow-state content, or failed post-generation validation.
- [ ] Runtime failures during non-dry-run flow map to the established runtime exit code and preserve safe diagnostic details.
- [ ] Tests cover valid handbook content, missing file, empty file, placeholders, missing required sections, missing selected studies/reports, missing patterns, missing trade-offs, missing anti-patterns, missing open questions, missing evidence pointers, implementation-decision wording, and unselected evidence references.
- [ ] Tests cover selected-evidence loading for multiple selected reports, missing selected report files, unreadable selected report files, project-index path mismatches, path escape attempts, and deterministic evidence manifest ordering.
- [ ] Tests cover prompt preview content without asserting brittle prose beyond required variables, selected evidence manifest, output path, required sections, no-decision instructions, selected-evidence-only instructions, and no-mutation instructions.
- [ ] Tests cover dry-run flow, fake-runtime successful generation, runtime success with missing artifact, runtime success with invalid artifact, runtime failure, validation failure after generation, selected-evidence failure before runtime, and flow-state write failure.
- [ ] Command tests prove `validate technical-handbook` and `prompt technical-handbook` do not call runtime, while non-dry-run `flow --to technical-handbook` uses a fake runtime seam.
- [ ] Offline verification passes with `go test ./...` from `/home/antonioborgerees/coding/ultraplan-go`.
- [ ] Race verification passes with `go test -race ./...` from `/home/antonioborgerees/coding/ultraplan-go`.
- [ ] CLI build passes with `go build ./cmd/ultraplan` from `/home/antonioborgerees/coding/ultraplan-go`.

## Non-Goals

- Generating or validating `reasoning/*.md` or final `reasoning.md`; that belongs to Sprint 20.
- Generating or validating `plan.md`; that belongs to Sprint 21.
- Making implementation decisions, architecture decisions, task breakdowns, evidence checklists, or final trade-off decisions in `technical-handbook.md`.
- Changing the select-stage contract beyond what is needed to consume a valid `sprint-index.md` and selected evidence for the distill stage.
- Reading, summarizing, or citing unselected evidence reports as handbook evidence.
- Implementing sprint implementation execution, smoke investigation execution, conformance review automation, issue tracking, project-management behavior, or Git mutation.
- Generating, validating, or requiring `.run-state.json`, `smoke.md`, `smoke.json`, product-generated `review.md`, `issues.md`, or `issues.json` as Planning Phase 2 artifacts.
- Mutating `project-index.md`, roadmap, project docs, prior sprint reviews, source repositories, workspace config, selected evidence reports, or Git state.
- Reusing study services, study prompt builders, study validators, source/dimension models, report/rating logic, summary generation, scheduler logic, or run-loop state for planning distillation.
- Adding a generic workflow engine, global prompt package, global validation package, plugin system, hosted service, browser UI, TUI, local API server, or multi-user collaboration features.
- Adding stable JSON output for sprint validate, prompt, or flow unless it is already supported by existing app infrastructure without expanding scope.
- Requiring real OpenCode, provider credentials, network access, or real-runtime smoke tests for default verification.

## Constraints

- `internal/sprint` owns technical-handbook validation, selected-evidence loading, prompt rendering, and flow-state transitions for the distill stage.
- `technical-handbook.md` must be treated as editable Markdown and as evidence distillation only; it is not a decision document or implementation plan.
- The distill stage must use only evidence reports selected by a valid `sprint-index.md`; it must not infer, silently add, or summarize uncataloged or unselected evidence.
- `project-index.md` remains a catalog, not a sprint plan, and must not be modified by the distill stage.
- Planning prompt execution must go through the generic platform runtime boundary only; product code must not directly invoke OpenCode or agentwrap adapters.
- `internal/sprint` must not depend on `internal/study` or extract study-owned behavior into shared packages for planning use.
- Flow-state schema remains strict and versioned; unsupported stages such as `implementation`, `execute`, `smoke`, `review`, and `issues` must not be accepted as supported Planning Phase 2 stages.
- Generated or previewed prompt content must avoid secrets, provider credentials, full environment dumps, unsafe raw runtime payloads, and unnecessary absolute local paths.
- Default tests must be deterministic, offline, fake-first, and free of long sleeps, network calls, provider credentials, real OpenCode execution, and Git mutation.
- Preserve prior review decisions: keep product behavior module-owned, keep CLI adapters thin, keep platform packages product-agnostic, preserve cause chains, use workspace-safe paths, keep state writes atomic, and avoid speculative abstractions.

## Dependencies

| Prior Sprint / Output | Required For | Notes |
| --------------------- | ------------ | ----- |
| Sprint 18 select stage | Valid `sprint-index.md`, selected-context parsing, prompt/flow command shape, and flow-state transition model through `sprint-index` | Sprint 19 extends the same sprint command family and must consume `sprint-index.md` rather than reselecting context. |
| Sprint 17 sprint artifact domain and flow state | Sprint resolution, artifact paths, supported planning stages, and atomic `flow-state.json` writes | Sprint 19 must preserve the exclusion of implementation, smoke, review, issue, and Git stages; the existing Sprint 19 `flow-state.json` contains unsupported legacy stages that must not be treated as supported. |
| Sprint 16 project domain and project index | Catalog loading and selected evidence path validation | Sprint 19 must validate selected evidence against the accepted project catalog boundary rather than reparsing project roots ad hoc in app code. |
| Sprint 14 validation, diagnostics, and JSON stability | Exit-code discipline, redaction, deterministic rendering, and safe validation findings | Distill-stage diagnostics should follow established validation and safe-output conventions. |
| Sprint 9 runtime integration | Runtime-backed non-dry-run flow | Runtime execution must use the generic platform runtime boundary and fake-runtime tests by default. |
| `projects/ultraplan-go/project-index.md` | Authoritative catalog for selected evidence and planning context | Selected evidence reports must be present here before the distill stage can read them. |
| `projects/ultraplan-go/sprints/19-distill-stage/sprint-index.md` | Authoritative selected evidence for this sprint | Must exist and validate before technical-handbook prompt rendering or non-dry-run flow can complete. |
| `studies/go-cli-study/reports/final/*.md` selected by sprint index | Source material for `technical-handbook.md` | Only selected and readable reports may be distilled or cited as evidence. |
| `projects/ultraplan-go/roadmap.md` | Sprint sequence and Planning Phase 2 scope | Defines Sprint 19 as the distill stage and keeps reason and plan stages separate. |
| `projects/ultraplan-go/docs/PRD.md` | Product planning behavior and non-goals | Requires governed planning through `plan.md`, editable artifacts, selected evidence application, and defers execution, smoke, review automation, issues, and Git mutation. |
| `projects/ultraplan-go/docs/TRD.md` | Technical requirements for planning validators, prompt rendering, flow, and commands | Defines `technical-handbook.md` validator checks, prompt rendering, runtime completion criteria, and required Phase 2 command shape. |
| `projects/ultraplan-go/docs/ARCHITECTURE.md` | Package ownership and dependency direction | Requires `internal/sprint` ownership, allows `sprint -> project/workspace/platform`, and forbids planning reuse of study semantics. |
| Prior sprint reviews | Carry-forward decisions and deferrals | Especially Sprint 17 atomic flow-state behavior, Sprint 16 catalog-only project boundary, Sprint 14 diagnostics/redaction, Sprint 9 runtime boundary, and Sprint 18 select-stage scope. |

## Review Expectations

| What | How Verified |
| ---- | ------------ |
| Required commands exist and are documented | Run `ultraplan sprint --help`, `ultraplan sprint <project> <sprint> validate --help`, `ultraplan sprint <project> <sprint> prompt --help`, and `ultraplan sprint <project> <sprint> flow --help`; inspect help for `technical-handbook` support and unsupported-stage rejection. |
| Technical-handbook validation is semantic, not just existence-based | Tests cover required sections, placeholders, selected evidence traces, unselected evidence references, no-decision checks, and missing/unreadable evidence diagnostics. |
| Selected evidence loading is constrained | Fixture tests prove only `sprint-index.md` selected evidence reports are read and that unselected, uncataloged, missing, unreadable, or escaping paths fail before runtime. |
| Prompt preview is deterministic and runtime-free | Command tests assert stable prompt variables, workspace-relative paths, selected evidence manifest, output path, required sections, no-decision instructions, no-mutation instructions, and no runtime invocation. |
| Flow dry-run is non-mutating | Tests assert dry-run does not create or overwrite `technical-handbook.md` and does not mark `technical-handbook` complete in `flow-state.json`. |
| Runtime-backed flow validates expected output | Fake-runtime tests cover successful generation, missing generated artifact, invalid generated artifact, runtime failure, selected-evidence failure, and validation failure after runtime success. |
| Flow-state transitions are correct | Tests assert successful flow marks `requirements`, `sprint-index`, and `technical-handbook` complete, readies or skips the next applicable reasoning stage correctly, and does not introduce unsupported stages. |
| Unsupported future stages remain excluded | Tests and code review confirm `execute`, `implementation`, `smoke`, `review`, `issues`, `.run-state.json`, `smoke.md`, `review.md`, and issue artifacts are not modeled as current planning stages. |
| Architecture boundaries are preserved | Import review confirms `internal/sprint` does not import `internal/study`, platform packages do not import product modules, and no global planning/workflow/validation/prompt package was added. |
| CLI output and errors are safe | Command tests assert stdout/stderr separation, no ANSI escapes, deterministic diagnostics, safe redaction, and correct exit-code mapping. |
| Default verification passes | Run `go test ./...`, `go test -race ./...`, and `go build ./cmd/ultraplan` from `/home/antonioborgerees/coding/ultraplan-go`. |
