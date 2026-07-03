# Sprint Requirements: Project Domain and Project Index

> Project: `ultraplan-go`
> Sprint: `16-project-domain-and-index`
> Purpose: the authoritative, human-readable sprint contract. All other sprint artifacts must satisfy these requirements.

## Sprint Goal

Model `projects/<project>` as a first-class planning root with project discovery, project status, and actionable `project-index.md` catalog validation without depending on study internals.

## Required Outputs

| Output | Path | Description |
| ------ | ---- | ----------- |
| Project domain model | `/home/antonioborgerees/coding/ultraplan-go/internal/project/domain.go` | Defines `Project`, `ProjectIndex`, catalog entries, project status, validation findings, and stable status categories for planning roots. |
| Project package documentation | `/home/antonioborgerees/coding/ultraplan-go/internal/project/doc.go` | Documents `internal/project` ownership, boundaries, and the rule that project indexes are catalogs rather than sprint plans. |
| Project discovery and resolution | `/home/antonioborgerees/coding/ultraplan-go/internal/project/discovery.go` | Discovers projects under `projects/`, resolves project references, validates filesystem-safe project names, and sorts results deterministically. |
| Project filesystem store | `/home/antonioborgerees/coding/ultraplan-go/internal/project/store_fs.go` | Reads project docs, `roadmap.md`, `project-index.md`, and sprint directory metadata through workspace-safe paths. |
| Project index parser | `/home/antonioborgerees/coding/ultraplan-go/internal/project/index.go` | Parses the catalog tables from `project-index.md` for source documents, contracts, evidence reports, reasoning templates, and review protocols. |
| Project validation service | `/home/antonioborgerees/coding/ultraplan-go/internal/project/validation.go` | Validates required project files, Markdown docs, catalog entry shape, and catalog path resolution with actionable diagnostics. |
| Project service API | `/home/antonioborgerees/coding/ultraplan-go/internal/project/service.go` | Provides project list, status, and validation use cases for the CLI without exposing parser or filesystem internals. |
| Project CLI wiring | `/home/antonioborgerees/coding/ultraplan-go/internal/app/project_commands.go` | Adds thin command wiring for `ultraplan project list`, `ultraplan project <project> status`, and `ultraplan project <project> validate`. |
| Top-level CLI registration and help | `/home/antonioborgerees/coding/ultraplan-go/internal/app/app.go` | Registers the `project` command family and documents it in top-level help without disturbing existing commands. |
| Project domain tests | `/home/antonioborgerees/coding/ultraplan-go/internal/project/project_test.go` | Covers project discovery, reference resolution, index parsing, path-safety failures, missing catalog references, and status summaries. |
| Project command tests | `/home/antonioborgerees/coding/ultraplan-go/internal/app/project_commands_test.go` | Covers command help, list/status/validate output, exit codes, deterministic ordering, and runtime-free behavior. |
| Sprint review evidence | `projects/ultraplan-go/sprints/16-project-domain-and-index/review.md` | Records implementation evidence, deviations, verification results, and review conclusions after sprint execution. |

## Acceptance Criteria

- [ ] `ultraplan project list` lists project roots discovered directly under `projects/` in deterministic name order.
- [ ] Project discovery ignores hidden entries, non-direct children, and non-directory entries under `projects/`.
- [ ] Project names are validated as filesystem-safe names; invalid names fail with actionable diagnostics.
- [ ] `ultraplan project <project> status` reports the presence and health of `docs/`, Markdown docs, `roadmap.md`, `project-index.md`, `sprints/`, and discovered sprint directories.
- [ ] `ultraplan project <project> validate` returns exit code `0` when required project files and catalog references are valid.
- [ ] `ultraplan project <project> validate` returns exit code `5` for validation failures such as missing required files, malformed catalog rows, unresolved catalog paths, or workspace-escaping paths.
- [ ] `project-index.md` parsing recognizes source documents, active contracts, evidence reports, reasoning templates, and review protocols from the catalog tables used by the existing `ultraplan-go` project index.
- [ ] Catalog paths resolve relative to the workspace and fail closed if they escape the workspace unless the entry is explicitly modeled as external.
- [ ] Missing catalog paths produce diagnostics that include the catalog section, entry name, referenced path, and suggested corrective action.
- [ ] `project-index.md` is treated as a catalog only; validation must not require or infer sprint task plans from it.
- [ ] Project commands do not import or call `internal/study`, study services, source/dimension models, report validators, rating parsing, summary generation, or run-loop scheduling.
- [ ] `internal/project` may use workspace path helpers and generic platform helpers, but platform packages must not import `internal/project` or any other product module.
- [ ] Project command output is deterministic, concise, and script-friendly in text mode.
- [ ] Project command failures preserve cause chains internally and render safe, actionable user-facing errors.
- [ ] Project validation does not invoke agentwrap, OpenCode, runtime health checks, prompt generation, network calls, source cloning, study execution, or sprint flow execution.
- [ ] Tests cover a valid fixture project, missing `roadmap.md`, missing `project-index.md`, malformed catalog table rows, unresolved catalog references, path escape attempts, empty `docs/`, and ambiguous or missing project references.
- [ ] Top-level help lists the `project` command family, and project-specific help documents `list`, `status`, and `validate`.
- [ ] The implementation introduces no global `internal/catalog`, `internal/validation`, `internal/reports`, `internal/prompts`, `internal/planning`, or workflow-engine package.
- [ ] Offline verification passes with `go test ./...` from `/home/antonioborgerees/coding/ultraplan-go`.
- [ ] Race verification passes with `go test -race ./...` from `/home/antonioborgerees/coding/ultraplan-go`.
- [ ] CLI build passes with `go build ./cmd/ultraplan` from `/home/antonioborgerees/coding/ultraplan-go`.

## Non-Goals

- Modeling sprint planning artifacts, planning stages, or `flow-state.json`; that belongs to Sprint 17.
- Generating or validating `requirements.md`, `sprint-index.md`, `technical-handbook.md`, `reasoning/*.md`, `reasoning.md`, or `plan.md`.
- Implementing select, distill, reason, or plan stage prompts, dry-run previews, or runtime-backed planning generation.
- Consuming study services, study source/dimension models, study report validators, rating parsing, summary generation, or run-loop scheduling.
- Implementing sprint execution, smoke investigation, automated review generation, issue tracking, project management features, or Git mutation.
- Adding a browser UI, hosted service, multi-user collaboration, TUI, metrics exporter, or local API server.
- Mutating `project-index.md`, `roadmap.md`, docs, sprint directories, or generated planning artifacts from project status or validation commands.
- Creating a generalized Markdown table parser or catalog framework outside `internal/project` unless required by this sprint's concrete project-index parsing behavior.

## Constraints

- `internal/project` owns project discovery, project docs discovery, project-index parsing, catalog validation, and project status behavior.
- `internal/app` may parse command arguments, call the project service, render text, and map exit codes only; project business rules must not live in command handlers.
- `internal/project` must not import `internal/study` or reuse study-owned concepts such as `Study`, `Source`, `Dimension`, report validation, rating parsing, summary generation, task scheduling, or run state.
- `internal/platform/runtime` must remain product-agnostic and must not know about projects, project indexes, catalog entries, or planning stages.
- All project-managed paths must be normalized through workspace-safe path handling before filesystem access.
- Catalog validation must fail closed on workspace path escapes and missing paths; it must not silently ignore invalid references.
- Project listing and status must avoid recursive scans of source repositories or evidence source trees; scans are limited to project roots, `docs/*.md`, catalog-referenced files, and direct sprint directories.
- Output must be deterministic across filesystems by sorting project names, docs, catalog entries, diagnostics, and sprint directory names.
- Default tests must be deterministic, offline, and fake-first; they must not require OpenCode, provider credentials, network access, Git commands, or long sleeps.
- The sprint must preserve prior review decisions: keep product behavior module-owned, keep CLI adapters thin, preserve cause chains, redact unsafe diagnostics, and avoid speculative shared abstractions.
- The implementation must not modify unrelated dirty files or revert user changes in the target repository.

## Dependencies

| Prior Sprint / Output | Required For | Notes |
| --------------------- | ------------ | ----- |
| Sprint 1-2 CLI, workspace, config, and health foundation | Command wiring and workspace resolution | Project commands must use existing workspace discovery and command output conventions. |
| Sprint 9 runtime integration | Boundary preservation | Project validation must not call runtime paths; platform runtime remains reusable only as generic infrastructure for later planning stages. |
| Sprint 14 validation, diagnostics, and JSON stability | Exit-code and diagnostics discipline | Project validation should reuse established exit-code conventions, safe diagnostics, and redaction practices. |
| Sprint 15 study-side release baseline | Phase transition | Planning Phase 2 starts after the study-side release scope; no Sprint 15 review was present in the required input set. |
| `projects/ultraplan-go/project-index.md` | Catalog shape and fixture expectations | Existing catalog sections and paths define the concrete project-index format Sprint 16 must parse and validate. |
| `projects/ultraplan-go/roadmap.md` | Sprint scope | Defines Sprint 16 as project domain and project-index validation, before sprint artifact flow state. |
| `projects/ultraplan-go/docs/PRD.md` | Product command requirements | Requires project discovery, inspection, and project-index validation in Phase 2. |
| `projects/ultraplan-go/docs/TRD.md` | Technical boundaries and validators | Defines `internal/project`, project model, catalog entries, path rules, and project command requirements. |
| `projects/ultraplan-go/docs/ARCHITECTURE.md` | Package ownership | Requires `internal/project` ownership and prohibits project behavior from depending on study internals. |

## Review Expectations

| What | How Verified |
| ---- | ------------ |
| Project commands exist and are documented | Run `ultraplan project --help`, `ultraplan project list --help`, `ultraplan project <project> status --help`, and `ultraplan project <project> validate --help`; inspect top-level help. |
| Project discovery is deterministic and scoped | Review `internal/project/project_test.go` fixtures for direct child directories, hidden entries, non-direct children, invalid names, and sorted output. |
| Project-index parsing matches the existing catalog format | Tests parse a fixture based on `projects/ultraplan-go/project-index.md` and assert all required catalog categories and entries. |
| Catalog reference validation is actionable | Tests cover missing paths, malformed rows, path escapes, and diagnostics containing section, entry, path, and guidance. |
| Project status reflects required files and sprint directories | Command tests assert status output for docs, roadmap, project index, catalog health, and direct sprint directory counts. |
| Runtime and study boundaries are preserved | Import review confirms `internal/project` does not import `internal/study`; project command tests assert no runtime invocation. |
| Architecture remains module-owned | Inspect file placement and imports; confirm no global catalog/validation/planning package and no product imports from platform packages. |
| Text output is safe and deterministic | Command tests assert stable stdout/stderr ordering, no ANSI escape sequences, no secret-like diagnostics, and correct exit codes. |
| Non-goals stayed excluded | Review diff for absence of sprint flow-state modeling, planning-stage generation, implementation execution, smoke/review/issues, Git mutation, and runtime-backed planning execution. |
| Default verification passes | Run `go test ./...`, `go test -race ./...`, and `go build ./cmd/ultraplan` from `/home/antonioborgerees/coding/ultraplan-go`. |
