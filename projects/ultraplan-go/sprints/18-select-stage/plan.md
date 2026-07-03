# Sprint Plan: Select Stage

> Project: `ultraplan-go`
> Sprint: `18-select-stage`
> Source: `projects/ultraplan-go/sprints/18-select-stage/reasoning.md`
> **Inputs Used:** `projects/ultraplan-go/sprints/18-select-stage/requirements.md`, `projects/ultraplan-go/sprints/18-select-stage/reasoning.md`, `projects/ultraplan-go/sprints/18-select-stage/sprint-index.md`, `projects/ultraplan-go/sprints/18-select-stage/technical-handbook.md`, `projects/ultraplan-go/docs/PRD.md`, `projects/ultraplan-go/docs/TRD.md`, `projects/ultraplan-go/docs/ARCHITECTURE.md`, `projects/ultraplan-go/roadmap.md`, `projects/ultraplan-go/project-index.md`, `templates/sprint-plan.md`, `templates/sprint-reasoning.md`; no files were present under `projects/ultraplan-go/sprints/18-select-stage/reasoning/*.md` when checked.

This plan executes `reasoning.md`. It does not invent architecture, scope, or decisions.

## Reasoning Source

| Source | Path | Use In This Plan |
| --- | --- | --- |
| Sprint Reasoning | `projects/ultraplan-go/sprints/18-select-stage/reasoning.md` | Controlling source for decisions, trade-offs, evidence, assumptions, risks, constraints, and handoff requirements. |
| Sprint Requirements | `projects/ultraplan-go/sprints/18-select-stage/requirements.md` | Authoritative sprint contract for required outputs, acceptance criteria, constraints, non-goals, dependencies, and verification. |
| Sprint Index | `projects/ultraplan-go/sprints/18-select-stage/sprint-index.md` | Selected contracts, evidence reports, Architecture reasoning template, required review protocols, and excluded context. |
| Technical Handbook | `projects/ultraplan-go/sprints/18-select-stage/technical-handbook.md` | Study-backed implementation patterns, trade-offs, anti-patterns, and open questions resolved by `reasoning.md`. |
| Area Reasoning | none present | No separate area-specific reasoning files were available; `reasoning.md` closes the Architecture decisions directly. |

## Context Handling

| Evidence | Treatment | Rationale |
| --- | --- | --- |
| PRD, TRD, Architecture, Roadmap, Project Index | Loaded and used. | Needed to verify Phase 2 select-stage scope, package ownership, dependency direction, prompt/flow requirements, artifacts, and non-goals. |
| Selected final study reports | Not reopened for this plan. | `technical-handbook.md` already distilled the selected reports into sprint-specific patterns and cited concrete references; no plan decision required deeper report evidence. |
| Code references in `/home/antonioborgerees/coding/ultraplan-go` | Not resolved for this plan. | This is a planning artifact, not implementation; code lookup should occur only for implementation questions. |

## Sprint Status

| Field | Value |
| --- | --- |
| Status | complete |
| Owner | Codex implementation agent |
| Start Date | 2026-06-14 or when implementation begins |
| Completion Date | 2026-06-14 |

## Decisions To Execute

| Decision | Source Section | Execution Implication |
| --- | --- | --- |
| Keep select-stage ownership in `internal/sprint`. | `reasoning.md#decision-1-keep-select-stage-ownership-in-internalsprint` | Put sprint-index domain structures, parsing, validation, prompt rendering, flow rules, state updates, and service use cases in `internal/sprint`; do not introduce forbidden global packages or reuse `internal/study`. |
| Validate `sprint-index.md` as structured editable Markdown with catalog subset checks. | `reasoning.md#decision-2-validate-sprint-indexmd-as-structured-editable-markdown-with-catalog-subset-checks` | Parse required sections and tables, reject placeholders and missing sections, validate selected contracts/evidence/templates/protocols against `project-index.md`, and produce deterministic actionable diagnostics. |
| Render deterministic runtime-free prompt previews for `sprint-index`. | `reasoning.md#decision-3-render-deterministic-runtime-free-prompt-previews-for-sprint-index` | Implement prompt preview without runtime invocation or artifact writes; include required paths, catalog entries, carry-forward constraints, output path, and no-mutation instructions. |
| Gate non-dry-run flow on requirements validation, generic runtime, artifact validation, and atomic state. | `reasoning.md#decision-4-gate-non-dry-run-flow-on-requirements-validation-generic-runtime-artifact-validation-and-atomic-state` | Dry-run is non-mutating; non-dry-run validates `requirements.md`, uses only the generic platform runtime, validates generated `sprint-index.md`, and updates `flow-state.json` atomically only after all gates pass. |
| Use thin CLI wiring and a sprint service with explicit store, runtime, clock, and output boundaries. | `reasoning.md#decision-5-use-thin-cli-wiring-and-a-sprint-service-with-explicit-store-runtime-clock-and-output-boundaries` | Keep `internal/app` focused on argument parsing, delegation, rendering, and exit-code mapping; keep workspace-safe reads/writes in `internal/sprint/store_fs.go`. |
| Keep CLI output calm, redacted, deterministic, and exit-code stable. | `reasoning.md#decision-6-keep-cli-output-calm-redacted-deterministic-and-exit-code-stable` | Return `0`, `2`, `5`, and established runtime failure codes as specified; keep stdout/stderr separated, no ANSI by default, and avoid leaking secrets or unsafe runtime payloads. |
| Verify with fixture-first unit, flow, and command tests plus offline build gates. | `reasoning.md#decision-7-verify-with-fixture-first-unit-flow-and-command-tests-plus-offline-build-gates` | Implement parser, validation, prompt, flow, state, command, architecture, security, race, and build evidence before review. |

## Requirements / Contracts To Satisfy

| Contract / Requirement ID | Required Behavior | Evidence Planned |
| --- | --- | --- |
| Architecture; REQ-AC-23 through REQ-AC-26 | `internal/sprint` owns select-stage business rules; `internal/app` stays thin; platform packages stay product-agnostic; no forbidden global packages are added. | Import review, code review, `internal/sprint` tests, and `internal/app` command tests. |
| Errors; REQ-AC-14, REQ-AC-29 through REQ-AC-31 | Diagnostics identify section, entry, name/path, observed problem, and repair suggestion; exit codes map usage, validation, and runtime failures correctly. | Validation tests and command tests for deterministic diagnostics and exit codes. |
| Configuration and Security; REQ-AC-15 through REQ-AC-17, REQ-AC-27 | Prompt and flow output uses workspace-relative paths by default, redacts sensitive values, and never prints unsafe raw runtime/config payloads. | Prompt tests, command redaction tests, and review checks. |
| Observability and CLI Surface; REQ-AC-01 through REQ-AC-03, REQ-AC-27 through REQ-AC-31 | Validate, prompt, and dry-run paths are runtime-free where required, scriptable, deterministic, separated by stream, and no ANSI by default. | `sprint_commands_test.go` covering stdout/stderr, help, unsupported stages, no runtime calls, and exit codes. |
| LLM Runtime; REQ-AC-04 through REQ-AC-06, REQ-AC-34 through REQ-AC-36 | Non-dry-run flow uses only generic platform runtime, requires generated artifact existence and validation, and uses fake runtime in default tests. | `flow_test.go` fake-runtime cases and import review against direct OpenCode, agentwrap adapter, shell, provider, or Git calls from `internal/sprint`. |
| Workflows and Persistence And Migrations; REQ-AC-18 through REQ-AC-22 | Supported Planning Phase 2 stages remain `requirements`, `sprint-index`, `technical-handbook`, `area-reasoning`, `reasoning`, and `plan`; Sprint 18 executes only through `sprint-index`; state writes stay atomic. | Flow-state tests for success, failure, dry-run, unsupported stages, area-reasoning skip semantics, and write failure. |
| Testing; REQ-AC-32 through REQ-AC-38 | Parser, validator, prompt, flow, and command behavior are covered with deterministic offline tests and fake runtime; full test, race, and build gates pass. | `go test ./...`, `go test -race ./...`, and `go build ./cmd/ultraplan` from `/home/antonioborgerees/coding/ultraplan-go`. |
| Documentation; Required Outputs | CLI help and review evidence record supported commands, unsupported stages, deviations, verification results, and review conclusions. | Help output tests and `projects/ultraplan-go/sprints/18-select-stage/review.md` after implementation. |

## Tasks

| Task | Executes | Implementation Steps | Required Evidence |
| --- | --- | --- | --- |
| [x] Task 1: Add sprint-index domain and parser in `internal/sprint/index.go`. | Decisions 1 and 2; REQ-AC-07 through REQ-AC-13, REQ-AC-32 | Define selected-context and excluded-context structures; parse required sections and Markdown tables; detect empty content, placeholders, missing sections, malformed rows, duplicates, and explicit no-template selection; keep helpers sprint-local. | `sprint_index_test.go` valid/invalid parsing, duplicate row, placeholder, and deterministic/actionable diagnostic assertions. |
| [x] Task 2: Implement semantic validation in `internal/sprint/validation.go`. | Decisions 1 and 2; REQ-AC-08 through REQ-AC-14, REQ-AC-30, REQ-AC-32 | Load project catalog entries through `internal/project`; validate selected contracts, evidence reports, reasoning templates, and review protocols as catalog subsets; validate excluded context includes known out-of-scope Phase 2 behavior; return diagnostics with section, selected entry, referenced name/path, observed problem, and suggested repair. | `sprint_index_test.go` covers invalid selected entries, path subset checks, excluded-context requirements, and actionable diagnostics. |
| [x] Task 3: Implement deterministic prompt rendering in `internal/sprint/prompts.go`. | Decisions 1 and 3; REQ-AC-02, REQ-AC-15 through REQ-AC-17, REQ-AC-27, REQ-AC-28, REQ-AC-33 | Build a runtime-free prompt preview from project slug, sprint slug, sprint path, requirements path, docs list, roadmap path, project-index path, available catalog entries, carry-forward constraints, output path, workspace-relative paths, selected-context instructions, and no-mutation rules. | Prompt tests assert required variables, catalog entries, output path, non-goals, no-mutation instructions, workspace-relative paths, no runtime calls, no artifact writes, no ANSI, and no absolute root leakage. |
| [x] Task 4: Extend store and service seams in `internal/sprint/store_fs.go` and `internal/sprint/service.go`. | Decisions 1 and 5; REQ-AC-23 through REQ-AC-26 | Add workspace-safe reads for requirements, project docs, roadmap, project index, and sprint-index artifacts; expose validate, prompt, and flow use cases through service methods; inject store, runtime, clock, and safe output dependencies without global registries. | Service and command tests prove validate/prompt/flow delegation and runtime-free preview paths. |
| [x] Task 5: Implement select-stage flow in `internal/sprint/flow.go` and `internal/sprint/state.go`. | Decisions 1, 4, and 6; REQ-AC-03 through REQ-AC-06, REQ-AC-18 through REQ-AC-22, REQ-AC-31, REQ-AC-34 through REQ-AC-35 | Implement `flow --to sprint-index --dry-run` as non-mutating; implement non-dry-run preflight on valid `requirements.md`; invoke only generic platform runtime; require generated `sprint-index.md`; run validation; atomically update `flow-state.json`; record safe failure state; reject unsupported future/execution stages. | `sprint_index_test.go` covers dry-run, fake-runtime success, invalid artifact validation failure, state transitions, non-mutation, and unsupported-stage rejection through command tests. |
| [x] Task 6: Wire CLI commands in `internal/app/sprint_commands.go` and `internal/app/app.go`. | Decisions 5 and 6; REQ-AC-01 through REQ-AC-03, REQ-AC-27 through REQ-AC-31 | Add thin command handling for `validate sprint-index`, `prompt sprint-index`, and `flow --to sprint-index`; update help and dispatch for sprint status/validate/prompt/flow; map usage, validation, runtime, and success exit codes; keep output plain, deterministic, and stream-separated. | `sprint_commands_test.go` covers help output, malformed args, unsupported stages, validate/prompt no-runtime behavior, dry-run flow, stdout/stderr separation, no ANSI, redaction, and exit codes. |
| [x] Task 7: Complete verification and review evidence. | Decision 7; REQ-AC-36 through REQ-AC-38; Review Expectations | Run offline tests, race tests, and CLI build; inspect imports and package boundaries; record deviations and verification in sprint review evidence; apply Architecture Review and Sprint Review protocols selected by `sprint-index.md`. | `go test ./...`, `go test -race ./...`, and `go build ./cmd/ultraplan` passed from `/home/antonioborgerees/coding/ultraplan-go`; review evidence recorded in `review.md`. |

## Evidence Checklist

| Evidence | Required Before Sprint Review |
| --- | --- |
| [x] Parser tests | Valid content, empty file, missing sections, placeholder content, malformed tables, duplicate entries, explicit no-template selection, excluded context, and deterministic diagnostic ordering. |
| [x] Catalog subset tests | Invalid selected contracts, evidence reports, reasoning templates, review protocols, and path mismatches fail with actionable diagnostics. |
| [x] Prompt tests | Prompt preview includes required variables and no-mutation instructions without runtime invocation or artifact writes. |
| [x] Flow tests | Dry-run, fake-runtime success, runtime success with invalid artifact, post-generation validation failure, and flow-state transitions. |
| [x] State transition tests | Success marks `requirements` and `sprint-index` complete, readies `technical-handbook`, leaves later required stages missing, and supports `area-reasoning` skip semantics. |
| [x] Unsupported-stage tests | `execute`, `implementation`, `smoke`, `review`, `issues`, `.run-state.json`, smoke artifacts, review artifacts, and issue artifacts are not accepted as supported Planning Phase 2 stages. |
| [x] Command tests | Help, malformed args, unsupported stages, exit codes, stdout/stderr separation, no ANSI, safe redaction, validate/prompt no-runtime paths, and dry-run flow. |
| [x] Architecture review | `internal/sprint` does not import `internal/study`; platform packages do not import product modules; no forbidden global planning, workflow, validation, prompt, report, or stage package exists. |
| [x] Security review | Workspace-safe paths, redacted diagnostics, no direct shell/OpenCode/provider/Git calls from sprint logic, and prompt no-mutation rules. |
| [x] Verification commands | `go test ./...`, `go test -race ./...`, and `go build ./cmd/ultraplan` pass from `/home/antonioborgerees/coding/ultraplan-go`. |
| [x] Review evidence | `projects/ultraplan-go/sprints/18-select-stage/review.md` records implementation evidence, deviations, verification results, and review conclusions. |

## Verification Commands

| Check | Command | Expected Result |
| --- | --- | --- |
| Offline test suite | `go test ./...` | Passes without network calls, provider credentials, real OpenCode execution, Git mutation, or long sleeps. |
| Race test suite | `go test -race ./...` | Passes with fake-runtime and deterministic fixture tests. |
| CLI build | `go build ./cmd/ultraplan` | Produces a buildable CLI binary. |
| Sprint help check | `ultraplan sprint --help` | Shows sprint status, validate, prompt, and flow surfaces consistently. |
| Validate help check | `ultraplan sprint <project> <sprint> validate --help` | Documents supported validation shape and rejects unsupported stages through usage errors. |
| Prompt help check | `ultraplan sprint <project> <sprint> prompt --help` | Documents `prompt sprint-index` as preview-only and runtime-free. |
| Flow help check | `ultraplan sprint <project> <sprint> flow --help` | Documents `--to sprint-index`, `--dry-run`, and unsupported-stage behavior. |

## Risks And Blockers

| Risk / Blocker | Source | Mitigation | Status |
| --- | --- | --- | --- |
| `internal/project` catalog APIs may not expose enough structured entries for subset validation. | `reasoning.md#assumptions-and-risks` | Use existing project APIs first; extend `internal/project` catalog parsing only if needed; do not create a global catalog package. | open |
| Sprint 17 flow-state model may need small additions for Sprint 18 transitions. | `reasoning.md#assumptions-and-risks` | Preserve strict schema versioning and add only Sprint 18-required fields with load/write tests. | open |
| Generic platform runtime seam may be too narrow for planning prompt execution. | `reasoning.md#assumptions-and-risks` | Extend platform runtime generically for prompt, workdir, timeout, model, expected output, metadata, and safe result fields only. | open |
| Markdown table parsing may reject reasonable human edits. | `reasoning.md#assumptions-and-risks` | Support reasonable whitespace/table formatting and provide precise diagnostics; revisit tolerance only from concrete rejected edits. | open |
| Multiple validation failures may create noisy diagnostics. | `reasoning.md#assumptions-and-risks` | Sort deterministically, group by section, and keep repair suggestions concise. | open |
| Prompt preview could leak sensitive config or absolute local details. | `reasoning.md#assumptions-and-risks` | Use redacted config summaries, workspace-relative paths by default, and command tests with sensitive strings. | open |
| Runtime writes `sprint-index.md` but atomic state write fails. | `reasoning.md#assumptions-and-risks` | Return non-zero, preserve last valid state, and require rerun or validation; test write-failure path. | open |
| Unsupported future stages could be accidentally modeled as current workflow. | `reasoning.md#assumptions-and-risks` | Keep allowed stage constants explicit and test rejection of execute, implementation, smoke, review, and issues. | open |
| `internal/sprint.Service` may grow too broad. | `reasoning.md#assumptions-and-risks` | Keep methods use-case-oriented and helpers file-scoped; review service shape before Sprint 19 extends it. | open |
| Fake runtime may not reveal real OpenCode integration gaps. | `reasoning.md#assumptions-and-risks` | Keep default verification fake-first; add gated external smoke only if maintainers request it after implementation. | open |

## Open Questions

| Question | Source | Handling For Sprint 18 |
| --- | --- | --- |
| Does the existing project catalog API fully support subset validation by name and path for contracts, evidence, reasoning templates, and protocols? | `reasoning.md` assumption | Implementation must inspect and use existing `internal/project` APIs first; if insufficient, extend project-owned parsing without global catalog abstractions. |
| Does the generic runtime boundary already support planning-stage expected-output metadata? | `reasoning.md` assumption | Implementation may extend generic platform runtime only with product-neutral fields; `internal/sprint` must not call OpenCode or agentwrap adapters directly. |
| How tolerant should Markdown table parsing be beyond the required section/table shapes? | `technical-handbook.md#open-questions-for-reasoning` and `reasoning.md` Decision 2 | Start strict on required sections and catalog semantics while allowing ordinary whitespace; record concrete rejected edits as follow-up evidence. |
| Should stable JSON output be added for sprint validate, prompt, or flow? | `requirements.md` non-goals and `reasoning.md` Decision 6 | Do not add it unless existing app infrastructure supports it without expanding scope; text output and exit codes are required evidence. |

## Review Inputs

Review should use:

- `projects/ultraplan-go/sprints/18-select-stage/requirements.md`
- `projects/ultraplan-go/sprints/18-select-stage/sprint-index.md`
- `projects/ultraplan-go/sprints/18-select-stage/technical-handbook.md`
- `projects/ultraplan-go/sprints/18-select-stage/reasoning.md`
- `projects/ultraplan-go/sprints/18-select-stage/plan.md`
- `system/protocols/architecture-review-protocol.md`
- `system/protocols/sprint-review-protocol.md`
- Implementation diff in `/home/antonioborgerees/coding/ultraplan-go`
- Verification evidence from `go test ./...`, `go test -race ./...`, and `go build ./cmd/ultraplan`

## Execution Log

| Date / Step | Action | Evidence / Notes |
| --- | --- | --- |
| 2026-06-14 planning | Created implementation plan from `reasoning.md`. | No code implemented; no final study reports reopened because `technical-handbook.md` supplied sufficient sprint-specific evidence. |
| 2026-06-14 implementation | Implemented sprint-index parser, validator, prompt preview, flow service, state transitions, and CLI commands. | Changed `/home/antonioborgerees/coding/ultraplan-go/internal/sprint/*`, `/home/antonioborgerees/coding/ultraplan-go/internal/app/sprint_commands.go`, `/home/antonioborgerees/coding/ultraplan-go/internal/app/app.go`, and related tests. |
| 2026-06-14 verification | Ran required Go verification from `/home/antonioborgerees/coding/ultraplan-go`. | `go test ./...` passed; `go test -race ./...` passed; `go build ./cmd/ultraplan` passed. A manual `go run ... --workspace /home/antonioborgerees/coding/ultra/.ultra ...` check failed because `.ultra` lacks `ultraplan.yml`, so command behavior is covered by workspace-fixture command tests instead. |

## Completion Criteria

| Criterion | Required State |
| --- | --- |
| [x] Tasks complete or deferred | All tasks above are complete. |
| [x] Verification complete | Required commands were run from `/home/antonioborgerees/coding/ultraplan-go`; workspace-specific CLI smoke limitation is documented above. |
| [x] Evidence satisfies `reasoning.md` | Parser, validation, prompt, flow, state, command, architecture, security, test, race, build, and review evidence exists. |
| [x] Scope preserved | No implementation execution, smoke execution, review automation, issue tracking, Git mutation, global workflow engine, or forbidden package was introduced. |
| [x] Review can proceed | `review.md` evaluates conformance against requirements, `sprint-index.md`, `reasoning.md`, this plan, and selected review protocols. |
