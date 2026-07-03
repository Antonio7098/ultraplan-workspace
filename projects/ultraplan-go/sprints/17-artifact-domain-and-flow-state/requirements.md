# Sprint Requirements: Sprint Artifact Domain and Flow State

> Project: `ultraplan-go`
> Sprint: `17-artifact-domain-and-flow-state`
> Purpose: the authoritative, human-readable sprint contract. All other sprint artifacts must satisfy these requirements.

## Sprint Goal

Model planning-stage sprint artifacts and durable `flow-state.json` state through `plan.md`, with runtime-free sprint status inspection and no implementation, smoke, review, or issue workflow stages.

## Required Outputs

| Output | Path | Description |
| ------ | ---- | ----------- |
| Sprint package documentation | `/home/antonioborgerees/coding/ultraplan-go/internal/sprint/doc.go` | Documents `internal/sprint` ownership, planning-stage scope through `plan.md`, dependency boundaries, and deferred execution/review/smoke/issue behavior. |
| Sprint domain model | `/home/antonioborgerees/coding/ultraplan-go/internal/sprint/domain.go` | Defines `Sprint`, `PlanningStage`, `StageStatus`, `StageState`, `FlowState`, status summaries, and validation findings for planning artifacts. |
| Sprint discovery and resolution | `/home/antonioborgerees/coding/ultraplan-go/internal/sprint/discovery.go` | Discovers sprint directories under `projects/<project>/sprints/`, resolves sprint references safely, validates filesystem-safe sprint slugs, and sorts results deterministically. |
| Planning artifact path rules | `/home/antonioborgerees/coding/ultraplan-go/internal/sprint/artifacts.go` | Centralizes expected artifact paths for `requirements.md`, `sprint-index.md`, `technical-handbook.md`, `reasoning/`, `reasoning.md`, `plan.md`, and `flow-state.json`. |
| Flow-state persistence | `/home/antonioborgerees/coding/ultraplan-go/internal/sprint/state.go` | Loads, validates, refreshes, and atomically writes versioned `flow-state.json` with strict schema and stage/status checks. |
| Sprint filesystem store | `/home/antonioborgerees/coding/ultraplan-go/internal/sprint/store_fs.go` | Reads sprint artifact existence and flow-state files through workspace-safe paths without invoking runtime or study behavior. |
| Sprint status service | `/home/antonioborgerees/coding/ultraplan-go/internal/sprint/service.go` | Provides runtime-free status inspection for one planning sprint and updates flow state when status refresh is requested by the command path. |
| Sprint CLI wiring | `/home/antonioborgerees/coding/ultraplan-go/internal/app/sprint_commands.go` | Adds thin command wiring for `ultraplan sprint <project> <sprint> status`, including help, argument parsing, rendering, and exit-code mapping. |
| Top-level CLI registration and help | `/home/antonioborgerees/coding/ultraplan-go/internal/app/app.go` | Registers the `sprint` command family and documents it in top-level help without disturbing existing commands. |
| Sprint package tests | `/home/antonioborgerees/coding/ultraplan-go/internal/sprint/sprint_test.go` | Covers discovery, reference resolution, stage ordering, artifact status derivation, flow-state strict loading, atomic writes, and unsupported stages/statuses. |
| Sprint command tests | `/home/antonioborgerees/coding/ultraplan-go/internal/app/sprint_commands_test.go` | Covers command help, status output, exit codes, deterministic ordering, runtime-free behavior, invalid sprint refs, and invalid flow-state diagnostics. |
| Sprint review evidence | `projects/ultraplan-go/sprints/17-artifact-domain-and-flow-state/review.md` | Records implementation evidence, deviations, verification results, and review conclusions after sprint execution. |

## Acceptance Criteria

- [ ] `internal/sprint` exists and owns sprint discovery, planning artifact path rules, flow-state persistence, flow-state validation, and sprint status behavior.
- [ ] `ultraplan sprint <project> <sprint> status` resolves the project through the existing project boundary and resolves the sprint under `projects/<project>/sprints/`.
- [ ] Sprint discovery inspects only direct child directories under the selected project's `sprints/` directory, ignores hidden and non-directory entries, validates filesystem-safe slugs, and returns deterministic sorted results.
- [ ] Sprint reference resolution supports exact references and unambiguous prefixes; missing, invalid, and ambiguous references fail with actionable diagnostics.
- [ ] The modeled planning stages are exactly `requirements`, `sprint-index`, `technical-handbook`, `area-reasoning`, `reasoning`, and `plan`, in that order.
- [ ] No `implementation`, `execute`, `smoke`, `review`, `issues`, Git, or project-management stages are modeled, persisted, rendered as supported stages, or accepted during strict flow-state loading.
- [ ] The modeled planning artifacts are exactly `requirements.md`, `sprint-index.md`, `technical-handbook.md`, `reasoning/`, `reasoning.md`, `plan.md`, and `flow-state.json`.
- [ ] `flow-state.json` records a schema version, project name, sprint slug, updated timestamp, and one state entry per supported planning stage.
- [ ] Each stage state records the stage name, status, workspace-relative artifact path, optional last-run timestamp, and optional safe error detail.
- [ ] Allowed stage statuses are exactly `missing`, `ready`, `complete`, `failed`, and `skipped`.
- [ ] Strict flow-state loading rejects malformed JSON, unsupported schema versions, missing required fields, unknown stages, duplicate stages, unknown statuses, paths that escape the sprint root, and project/sprint mismatches.
- [ ] Flow-state writes are atomic: write a temporary file in the same directory, flush and close it, rename it over the prior file, and best-effort sync the parent directory.
- [ ] Atomic write failure preserves the last valid `flow-state.json` and surfaces a non-zero command result with the state path and cause preserved internally.
- [ ] Existing planning artifacts are inspected to derive status without invoking runtime: present non-empty files are reflected as `complete`; missing stages blocked by earlier incomplete stages are `missing`; the first missing stage after completed predecessors is `ready`; recorded stage errors remain `failed`.
- [ ] `area-reasoning` is marked `skipped` only when a readable `sprint-index.md` explicitly selects no reasoning templates; otherwise missing `reasoning/*.md` files are not silently skipped.
- [ ] `ultraplan sprint <project> <sprint> status` prints deterministic plain text that includes project, sprint, sprint root, flow-state path, and every supported planning stage with status and artifact path.
- [ ] `ultraplan sprint <project> <sprint> status` returns exit code `0` when sprint resolution succeeds and flow state is valid or can be refreshed from artifact inspection.
- [ ] `ultraplan sprint <project> <sprint> status` returns exit code `5` for invalid `flow-state.json` content, unsupported stages, unsupported statuses, project/sprint mismatch, or unsafe artifact paths.
- [ ] `ultraplan sprint <project> <sprint> status` returns the established workspace or usage exit codes for missing projects, missing sprint directories, invalid references, and malformed command arguments.
- [ ] Status inspection and flow-state refresh do not call agentwrap, OpenCode, runtime health checks, prompt generation, network calls, source cloning, study execution, project-index mutation, or planning-stage generation.
- [ ] `internal/sprint` may import `internal/project` for project-root and catalog boundaries, `internal/workspace` for path safety, and generic platform helpers, but it must not import `internal/study` or reuse study source/dimension/report/run-loop concepts.
- [ ] `internal/platform/*` packages remain product-agnostic and do not import `internal/sprint`, `internal/project`, `internal/study`, `internal/app`, or other product modules.
- [ ] The implementation introduces no global `internal/planning`, `internal/workflow`, `internal/stages`, `internal/validation`, `internal/reports`, `internal/prompts`, or generic workflow-engine package.
- [ ] Tests cover a valid sprint with no artifacts, a sprint with all planning artifacts present, a sprint with partial artifacts, a missing sprint, an ambiguous sprint reference, malformed flow-state JSON, unsupported schema version, unsupported `implementation` or `review` stages, unknown statuses, and path escape attempts.
- [ ] Command tests assert stdout/stderr separation, deterministic ordering, help text, exit codes, no ANSI escape sequences, and no runtime invocation.
- [ ] Offline verification passes with `go test ./...` from `/home/antonioborgerees/coding/ultraplan-go`.
- [ ] Race verification passes with `go test -race ./...` from `/home/antonioborgerees/coding/ultraplan-go`.
- [ ] CLI build passes with `go build ./cmd/ultraplan` from `/home/antonioborgerees/coding/ultraplan-go`.

## Non-Goals

- Generating, repairing, or content-validating `requirements.md`, `sprint-index.md`, `technical-handbook.md`, `reasoning/*.md`, `reasoning.md`, or `plan.md` beyond artifact existence and flow-state status inspection.
- Implementing `ultraplan sprint <project> <sprint> validate [stage]`; stage validation begins in later planning sprints.
- Implementing `ultraplan sprint <project> <sprint> prompt <stage>` or prompt preview output.
- Implementing `ultraplan sprint <project> <sprint> flow --to <stage>` or runtime-backed planning artifact generation.
- Validating `sprint-index.md` selections against `project-index.md`; that belongs to Sprint 18.
- Distilling `technical-handbook.md`, generating area reasoning, generating final reasoning, or generating `plan.md`; those belong to Sprints 19-21.
- Executing implementation tasks from `plan.md`, running smoke investigations, generating conformance reviews, creating issues, mutating Git state, or modeling those as current flow stages.
- Reusing study services, study source/dimension models, study report validators, rating parsing, summary generation, task scheduling, run-loop state, or study prompt builders.
- Adding JSON output for sprint status unless it is already supported by existing app infrastructure without changing the sprint scope.
- Building a hosted service, browser UI, TUI, multi-user collaboration, metrics exporter, local API server, or project-management workflow.

## Constraints

- `internal/sprint` owns sprint planning artifact behavior and durable `flow-state.json` state through `plan.md`.
- `internal/app` may parse command arguments, call the sprint service, render output, and map exit codes only; sprint business rules must not live in command handlers.
- `internal/sprint` must not import `internal/study` or reuse study-owned types or behavior.
- `internal/sprint` may depend on `internal/project` only for the project planning root and catalog boundary established in Sprint 16.
- `internal/platform/runtime` must remain generic and must not know about projects, sprints, planning stages, or artifact semantics.
- Flow-state schema must be strict and versioned; do not add backward-compatibility paths for unsupported pre-sprint state shapes unless a concrete migration requirement is added to the sprint plan.
- Artifact and flow-state paths must be normalized and verified to stay inside the selected sprint root before filesystem access or persistence.
- Status derivation must be deterministic across filesystems by sorting sprint directories, diagnostics, and reasoning files.
- Default tests must be deterministic, offline, and fake-first; they must not require OpenCode, provider credentials, network access, Git commands, subprocess execution, or long sleeps.
- The sprint must preserve prior review decisions: keep product behavior module-owned, keep CLI adapters thin, preserve cause chains, redact unsafe diagnostics, avoid speculative shared abstractions, and keep generated/user-editable artifacts as filesystem state.
- The implementation must not modify unrelated dirty files or revert user changes in the target repository.

## Dependencies

| Prior Sprint / Output | Required For | Notes |
| --------------------- | ------------ | ----- |
| Sprint 16 project domain and project index | Project resolution and planning-root boundary | Sprint 17 must build on accepted `internal/project` discovery/status/validation behavior without duplicating project catalog ownership. |
| Sprint 14 validation, diagnostics, and JSON stability | Exit-code and diagnostics discipline | Sprint status should reuse established exit-code conventions, safe diagnostics, redaction, and deterministic command rendering. |
| Sprint 12 durable run-loop state | Persistence precedent | Flow-state persistence should follow the accepted atomic-write and strict-schema discipline without reusing study run-loop semantics. |
| Sprint 9 runtime integration | Boundary preservation | Sprint status must remain runtime-free; later planning generation may use platform runtime, but this sprint must not call it. |
| `projects/ultraplan-go/project-index.md` | Catalog and package-boundary context | Establishes Planning Phase 2 ownership rules and available contracts/evidence. |
| `projects/ultraplan-go/roadmap.md` | Sprint scope and sequence | Defines Sprint 17 as sprint artifact domain plus flow state, before select/distill/reason/plan generation. |
| `projects/ultraplan-go/docs/PRD.md` | Product planning scope | Requires sprint planning artifacts and `flow-state.json` through `plan.md` while deferring execution, smoke, review automation, issues, and Git mutation. |
| `projects/ultraplan-go/docs/TRD.md` | Technical requirements | Defines `internal/sprint`, planning stage model, flow-state fields/statuses, command surface, strict loading, atomic writes, and deferred technical requirements. |
| `projects/ultraplan-go/docs/ARCHITECTURE.md` | Package ownership | Requires `internal/sprint` ownership and prohibits sprint planning behavior from depending on study internals. |

## Review Expectations

| What | How Verified |
| ---- | ------------ |
| Sprint status command exists and is documented | Run `ultraplan sprint --help` and `ultraplan sprint <project> <sprint> status --help`; inspect top-level help for the `sprint` command family. |
| Sprint discovery is deterministic and scoped | Review `internal/sprint/sprint_test.go` fixtures for direct child directories, hidden entries, non-direct children, invalid slugs, ambiguous prefixes, and sorted output. |
| Stage model is limited to planning through `plan` | Inspect `internal/sprint/domain.go` and tests to confirm exactly six supported stages and no implementation/smoke/review/issues stages. |
| Flow-state strict loading is enforced | Tests cover malformed JSON, unsupported schema version, unknown stages, duplicate stages, unknown statuses, unsafe paths, and project/sprint mismatch. |
| Flow-state writes are atomic and loud on failure | Tests inject write/rename failures and assert prior valid state is preserved and errors include the flow-state path. |
| Artifact inspection reflects existing planning files | Fixture tests cover empty, partial, and complete sprint artifact sets and assert derived `missing`, `ready`, `complete`, `failed`, and `skipped` statuses. |
| CLI output is safe and deterministic | Command tests assert stable stdout/stderr, stage ordering, exit codes, no ANSI escape sequences, and actionable diagnostics. |
| Runtime and study boundaries are preserved | Import review confirms `internal/sprint` does not import `internal/study`; command tests assert no runtime invocation. |
| Architecture remains module-owned | Inspect file placement and imports; confirm no global planning/workflow/stage/validation package and no product imports from platform packages. |
| Non-goals stayed excluded | Review diff for absence of `sprint validate`, `sprint prompt`, `sprint flow`, runtime-backed planning generation, implementation execution, smoke/review/issues, Git mutation, and project-index mutation. |
| Default verification passes | Run `go test ./...`, `go test -race ./...`, and `go build ./cmd/ultraplan` from `/home/antonioborgerees/coding/ultraplan-go`. |
