# Sprint Requirements: Study Initialization From YAML

> Project: `ultraplan-go`
> Sprint: `04-study-init-yaml`
> Purpose: the authoritative, human-readable sprint contract. All other sprint artifacts must satisfy these requirements.

## Sprint Goal

Implement `ultraplan study init <study-init.yml>` so explicit YAML study definitions generate deterministic study directories, normalized YAML, editable dimension files, a study README, optional shallow source clones, and safe dry-run/force behavior.

## Required Outputs

| Output | Path | Description |
| ---------- | -------- | --------------- |
| Requirements | `projects/ultraplan-go/sprints/04-study-init-yaml/requirements.md` | Sprint contract defining scope, outputs, acceptance criteria, constraints, dependencies, and review expectations. |
| Sprint Index | `projects/ultraplan-go/sprints/04-study-init-yaml/sprint-index.md` | Selected contracts, project docs, and evidence reports for study initialization, filesystem writes, CLI behavior, tests, and safety constraints. |
| Technical Handbook | `projects/ultraplan-go/sprints/04-study-init-yaml/technical-handbook.md` | Implementation guidance for YAML parsing, normalization, generated artifacts, clone behavior, dry-run/force safety, and package boundaries. |
| Architecture Reasoning | `projects/ultraplan-go/sprints/04-study-init-yaml/reasoning/architecture.md` | Reasoning for module ownership, dependency direction, and study initialization composition. |
| Consolidated Reasoning | `projects/ultraplan-go/sprints/04-study-init-yaml/reasoning.md` | Final sprint reasoning that maps implementation decisions to requirements and selected evidence. |
| Implementation Plan | `projects/ultraplan-go/sprints/04-study-init-yaml/plan.md` | Ordered implementation tasks, tests, and verification commands. |
| Study initialization service | `/home/antonioborgerees/coding/ultraplan-go/internal/study/init.go` | Owns study initialization requests, validation, planning, artifact creation, normalization, and force/dry-run behavior. |
| Study initialization YAML parser | `/home/antonioborgerees/coding/ultraplan-go/internal/study/init_yaml.go` | Parses and validates `study-init.yml` input into explicit source and dimension definitions. |
| Study artifact rendering | `/home/antonioborgerees/coding/ultraplan-go/internal/study/init_render.go` | Renders generated dimension Markdown, study README content, and normalized `study-init.yml`. |
| Study source cloning | `/home/antonioborgerees/coding/ultraplan-go/internal/study/init_clone.go` | Provides shallow clone planning/execution for URL-backed source items through a testable command boundary. |
| Study init CLI command | `/home/antonioborgerees/coding/ultraplan-go/internal/app/study_init_commands.go` | Wires `ultraplan study init <study-init.yml>` flags, output, and exit behavior into the existing CLI app. |
| Study initialization tests | `/home/antonioborgerees/coding/ultraplan-go/internal/study/init_test.go` | Covers YAML validation, normalization, generated files, dry-run, force behavior, clone planning, and failure cases. |
| Study init command tests | `/home/antonioborgerees/coding/ultraplan-go/internal/app/study_init_commands_test.go` | Covers command flags, text output, exit codes, dry-run output, overwrite protection, and workspace/output path behavior. |

## Acceptance Criteria

- [ ] `ultraplan study init <study-init.yml>` is available in help output and returns deterministic, actionable errors for invalid usage.
- [ ] The command parses YAML with `name`, `description`, `repos.count`, `repos.items`, `dimensions.count`, and `dimensions.items` fields.
- [ ] Source items support explicit `name`, `url`, `path`, and `description` fields, with unique filesystem-safe names.
- [ ] Dimension items support explicit `number`, `name`, `title`, `description`, `purpose`, `steps`, `citations`, and `questions` fields.
- [ ] Missing required fields, invalid YAML, duplicate source names, duplicate dimension numbers/slugs, unsafe names, and unsafe output paths fail with non-zero actionable errors.
- [ ] `repos.count` and `dimensions.count` cannot be less than the number of explicit items.
- [ ] Because runtime-assisted completion is not in this sprint, `repos.count` or `dimensions.count` greater than explicit items fails with guidance to provide explicit items or wait for assisted completion support.
- [ ] Dimension numbers are normalized to two digits in memory, generated filenames, generated Markdown headings, and normalized `study-init.yml`.
- [ ] Dimension slugs are normalized to kebab-case and used in generated filenames shaped `NN-slug.md`.
- [ ] Generated dimension Markdown preserves supplied purpose, steps, citations, and questions in human-editable form.
- [ ] Successful initialization creates `studies/<study>/study-init.yml`, `studies/<study>/README.md`, `studies/<study>/dimensions/`, `studies/<study>/sources/`, `studies/<study>/reports/source/`, and `studies/<study>/reports/final/` under the resolved workspace unless `--output <dir>` is supplied.
- [ ] Generated `README.md` describes the study, sources, dimensions, generated paths, and currently supported next commands without advertising deferred runtime, target, or sprint workflows as implemented.
- [ ] Generated normalized `study-init.yml` is deterministic, contains accepted explicit source and dimension items, and avoids absolute paths unless the user supplied an explicitly allowed custom output path that must be preserved for diagnostics.
- [ ] `--dry-run` prints the directories, files, and clone actions that would be performed without creating, overwriting, deleting, or cloning anything.
- [ ] Existing study directories are protected by default; re-running initialization for an existing study fails unless `--force` is set.
- [ ] `--force` may overwrite generated files inside the selected study directory, but must not remove or modify paths outside that directory.
- [ ] Source cloning is skipped when `--no-clone` is set.
- [ ] When cloning is enabled for URL-backed source items, clone execution uses `git clone --depth 1 <url> <sources/name>` through a non-shell command invocation.
- [ ] Clone failures are reported without hiding successful generated artifact creation; partial clone completion returns a non-zero partial/failure status according to existing CLI exit-code conventions.
- [ ] Tests cover clone behavior through fakes and do not require network access, GitHub access, OpenCode, provider credentials, or real cloned repositories.
- [ ] Study initialization uses existing workspace discovery and global `--workspace` behavior from prior sprints.
- [ ] `go test ./...` passes from `/home/antonioborgerees/coding/ultraplan-go`.
- [ ] `go build ./cmd/ultraplan` passes from `/home/antonioborgerees/coding/ultraplan-go`.

## Non-Goals

- Do not implement runtime-assisted source or dimension completion, runtime suggestion prompts, suggestion cache artifacts, or `--no-assist` behavior in this sprint.
- Do not implement source URL verification, replacement source suggestions, repair prompts, or any OpenCode/agentwrap runtime execution.
- Do not implement Markdown document source discovery, frontmatter parsing, `applicable_dimensions`, or applicability filtering; those remain Sprint 5 scope.
- Do not implement prompt composition, analysis runs, synthesis, report validation, summary generation, run state, retries, cancellation, status, diagnostics persistence, or code-reference extraction.
- Do not implement target or sprint CLI workflows; those remain outside the current study-side release.
- Do not add stable public JSON output schemas for study initialization unless already available through existing app output helpers without expanding the sprint scope.
- Do not recursively scan source repositories or inspect cloned source contents during initialization.

## Constraints

- Study initialization behavior must live in `internal/study`; do not create global `internal/validation`, `internal/reports`, `internal/scheduler`, `internal/prompts`, or similar technical-layer packages for this sprint.
- CLI parsing and command output wiring must stay in `internal/app`; YAML validation, artifact planning, rendering, and filesystem behavior must stay outside CLI command handlers.
- `internal/study` may depend on `internal/workspace` and platform packages, but platform packages must not import `internal/study`.
- All workspace-managed paths must be normalized and checked so initialization cannot write outside the resolved workspace or explicitly selected output root.
- Dry-run must share the same planning and validation path as real initialization, with writes and clone execution disabled at the final execution boundary.
- Destructive overwrite behavior must require `--force` and be scoped to the selected study directory only.
- Git clone execution must avoid shell interpolation and must be behind a testable boundary so normal tests remain deterministic and offline.
- Generated artifacts must be deterministic, reviewable in Git, human-editable, and free of secrets.
- Normal unit and command tests must not require OpenCode, provider credentials, network access, existing workspaces, or real source repositories.
- Error messages must preserve cause context and tell the user which field, path, or source item needs correction.

## Dependencies

| Prior Sprint / Output | Required For | Notes |
| ------------------------ | ------------------- | --------- |
| Sprint 1: Go Module, CLI Shell, and App Composition | Command registration, buildable CLI entrypoint, and baseline tests | Provides `cmd/ultraplan`, `internal/app`, exit-code conventions, and package skeletons. |
| Sprint 2: Workspace, Config, Logging, and Health Skeleton | Workspace discovery, path safety, and command context | Study initialization must use existing workspace resolution and global `--workspace` behavior. |
| Sprint 3: Study Domain, Listing, and Resolution | Existing study domain concepts and listing compatibility | Generated studies must be discoverable and listable by the Sprint 3 commands. |
| `projects/ultraplan-go/docs/PRD.md` | Product study initialization behavior | Defines YAML study initialization, generated artifacts, dimension file requirements, dry-run, force, no-clone, and deferred target/sprint scope. |
| `projects/ultraplan-go/docs/TRD.md` | Technical study initialization behavior | Defines input schema, generated structure, force behavior, shallow clone requirements, path handling, and test expectations. |
| `projects/ultraplan-go/docs/ARCHITECTURE.md` | Package boundaries | Defines module ownership and dependency direction for `study`, `workspace`, app, and platform packages. |
| Go toolchain and local `git` executable | Build/test and optional clone execution | Tests must fake clone execution; real clone behavior may be manually reviewed only when Git/network are available. |
| None | Prior sprint review decisions | No prior sprint `review.md` files exist to carry forward. |

## Review Expectations

| What | How Verified |
| ------------- | ----------------------- |
| Sprint scope matches roadmap Sprint 4 | Review `roadmap.md` Sprint 4 against implemented files, commands, flags, and non-goals. |
| Required outputs exist | Check each path listed in Required Outputs. |
| Study init command works | Run `ultraplan study init <study-init.yml>` against a fixture workspace and inspect generated structure. |
| Dry-run is non-mutating | Run `ultraplan study init <study-init.yml> --dry-run` and verify no files, directories, or clone targets are created. |
| Generated artifacts are deterministic | Run initialization from the same fixture twice with clean output and compare generated `study-init.yml`, `README.md`, and dimension files. |
| Force behavior is safe | Verify existing study directories are protected by default and `--force` only affects the selected study directory. |
| YAML validation is actionable | Run invalid fixture cases and verify non-zero exits name the failing field, source, dimension, or path. |
| Clone behavior is bounded | Review tests/fakes and implementation to confirm URL-backed clones use `git clone --depth 1` without shell interpolation and normal tests do not require network. |
| Generated studies remain listable | Run Sprint 3 listing commands against an initialized fixture study and verify sources/dimensions are discovered deterministically. |
| Architecture boundaries are preserved | Review imports to confirm platform packages do not import `internal/study` and CLI handlers do not own initialization logic. |
| Deferred scope remains deferred | Confirm no assisted runtime completion, Markdown frontmatter/applicability, runtime execution, target, or sprint workflow implementation is included. |
| Build and test gate passes | Run `go test ./...` and `go build ./cmd/ultraplan` from `/home/antonioborgerees/coding/ultraplan-go`. |
