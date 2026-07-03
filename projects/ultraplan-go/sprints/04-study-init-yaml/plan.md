# Sprint Plan: Study Initialization From YAML

> Project: `ultraplan-go`
> Sprint: `04-study-init-yaml`
> Source: `projects/ultraplan-go/sprints/04-study-init-yaml/reasoning.md`
> **Inputs Used:** `projects/ultraplan-go/project-index.md`, `projects/ultraplan-go/sprints/04-study-init-yaml/requirements.md`, `projects/ultraplan-go/docs/PRD.md`, `projects/ultraplan-go/docs/TRD.md`, `projects/ultraplan-go/docs/ARCHITECTURE.md`, `projects/ultraplan-go/sprints/04-study-init-yaml/sprint-index.md`, `projects/ultraplan-go/sprints/04-study-init-yaml/technical-handbook.md`, `projects/ultraplan-go/sprints/04-study-init-yaml/reasoning/architecture.md`, `projects/ultraplan-go/sprints/04-study-init-yaml/reasoning.md`, `templates/sprint-plan.md`

This plan executes `reasoning.md`. It must not invent architecture, scope, or decisions.

## Reasoning Source

- **Sprint Reasoning:** `projects/ultraplan-go/sprints/04-study-init-yaml/reasoning.md`
- **Sprint Index:** `projects/ultraplan-go/sprints/04-study-init-yaml/sprint-index.md`
- **Technical Handbook:** `projects/ultraplan-go/sprints/04-study-init-yaml/technical-handbook.md`
- **Area Reasoning:** `projects/ultraplan-go/sprints/04-study-init-yaml/reasoning/architecture.md`

## Sprint Status

- **Status:** complete
- **Owner:** implementation agent
- **Start Date:** 2026-05-30
- **Completion Date:** 2026-05-30

## Decisions To Execute

| Decision | Source Section | Execution Implication |
| --- | --- | --- |
| Study-Owned Init Use Case | `reasoning.md#decision-1-study-owned-init-use-case` | Implement parsing, validation, normalization, planning, rendering, filesystem mutation, clone policy, and force/dry-run behavior in `internal/study`; keep `internal/app` as command wiring and output only. |
| Thin CLI Command With Existing Workspace Behavior | `reasoning.md#decision-2-thin-cli-command-with-existing-workspace-behavior` | Add `ultraplan study init <study-init.yml>` with `--dry-run`, `--force`, `--no-clone`, and `--output <dir>` in `internal/app`; reuse existing workspace resolution and exit-code conventions. |
| Single Validated Plan For Dry Run, Writes, Force, And Path Safety | `reasoning.md#decision-3-single-validated-plan-for-dry-run-writes-force-and-path-safety` | Build one deterministic plan listing directories, files, and clone actions; use it for dry-run output and real execution; validate all paths before mutation; scope force to known generated files inside the selected study directory. |
| Explicit YAML-Only Initialization With Deterministic Normalization | `reasoning.md#decision-4-explicit-yaml-only-initialization-with-deterministic-normalization` | Parse only explicit YAML items; fail unsupported count shortages with assisted-completion guidance; normalize dimension numbers and slugs once and render from normalized values. |
| Optional Shallow Clone Through Narrow Non-Shell Runner | `reasoning.md#decision-5-optional-shallow-clone-through-narrow-non-shell-runner` | Execute URL-backed clone actions after artifact writes through a fakeable non-shell runner equivalent to `git clone --depth 1 <url> <dest>`; skip under `--dry-run` and `--no-clone`; aggregate partial failures. |
| Deterministic Generated Artifacts And Reviewable Tests | `reasoning.md#decision-6-deterministic-generated-artifacts-and-reviewable-tests` | Render stable normalized YAML, README, and dimension Markdown; test generated artifacts with fixture/golden-style assertions where stability matters. |
| Deferred Scope Remains Explicitly Deferred | `reasoning.md#decision-7-deferred-scope-remains-explicitly-deferred` | Do not add runtime-assisted completion, OpenCode/agentwrap execution, Markdown frontmatter/applicability, prompt composition, run state, target/sprint workflows, or new stable JSON schemas. |

## Requirements / Contracts To Satisfy

| Contract / Requirement ID | Required Behavior | Evidence Planned |
| --- | --- | --- |
| Architecture contract; SR-CON-01 to SR-CON-03 | `internal/study` owns study init behavior; `internal/app` owns command wiring; platform packages do not import `study`. | Import review, command delegation tests, `go test ./...`, `go build ./cmd/ultraplan`. |
| Errors contract; SR-AC-01, SR-AC-05 to SR-AC-07, SR-AC-19 | Invalid usage, YAML, duplicate names, unsafe paths, count mismatches, overwrite protection, and clone failures return actionable non-zero errors. | Parser/service validation tests, command exit-code tests, clone partial-failure tests. |
| Security contract; SR-AC-03, SR-AC-05, SR-AC-14 to SR-AC-18, SR-CON-04 to SR-CON-09 | Validate names and paths; keep writes under allowed roots; avoid shell interpolation; scope force; avoid secret-generating behavior. | Path escape tests, unsafe name tests, force sentinel tests, fake clone arg assertions, source review for no shell execution strings. |
| Testing contract; SR-AC-20, SR-AC-52, SR-AC-53 | Normal tests are deterministic and offline; build/test gates pass. | `internal/study/init_test.go`, `internal/app/study_init_commands_test.go`, `go test ./...`, `go build ./cmd/ultraplan`. |
| Documentation contract; SR-AC-10 to SR-AC-13 | Generated README, normalized YAML, and dimension Markdown are deterministic, human-editable, and scoped to currently supported commands. | Fixture/golden comparisons and README review for deferred workflow claims. |
| CLI Surface contract; SR-AC-01, SR-AC-14, SR-AC-15, SR-AC-17, SR-AC-19, SR-AC-51 | Command help, flags, workspace behavior, dry-run output, overwrite protection, clone failure output, and exit mapping work. | Command tests for help, usage, flags, workspace/output path behavior, dry-run, force, and partial failure. |
| SR-AC-02 to SR-AC-04 | YAML schema supports required study, source, and dimension fields. | Parser tests for full valid schema and missing-field failures. |
| SR-AC-08 to SR-AC-10 | Dimension numbers and slugs normalize deterministically; Markdown preserves purpose, steps, citations, and questions. | Normalization tests and generated dimension Markdown assertions. |
| SR-AC-11 to SR-AC-13 | Successful init creates required structure and deterministic artifacts under workspace/default or `--output`. | Workflow tests inspecting directory tree and generated file contents. |
| SR-AC-16 to SR-AC-20 | Clone behavior is optional, shallow, non-shell, fake-tested, and partial failures remain visible. | Fake clone runner tests for skip, args, order, and failure aggregation. |
| SR-NG-01 to SR-NG-07 | Deferred runtime, Markdown discovery/applicability, prompt, run, target, sprint, diagnostics, and JSON-schema scope remains absent. | Non-goal review checklist and import/help/README review. |

## Tasks

- [x] **Task 1: Inspect Existing Implementation Boundaries**
  > Executes: `Decision 1`, `Decision 2`, Architecture contract, SR-AC-51
  - [x] Inspect `/home/antonioborgerees/coding/ultraplan-go/internal/app`, `/home/antonioborgerees/coding/ultraplan-go/internal/study`, and `/home/antonioborgerees/coding/ultraplan-go/internal/workspace` for existing command construction, exit-code conventions, workspace resolution APIs, study listing expectations, and test patterns.
  - [x] Identify existing YAML dependency or add the smallest project-consistent YAML parser dependency only if none exists.
  - [x] Record any implementation-time conflict with `reasoning.md` before proceeding instead of inventing new architecture.

- [x] **Task 2: Define Study Init Request, Result, Plan, And Error Shapes**
  > Executes: `Decision 1`, `Decision 3`, `Decision 5`, Errors contract
  - [x] Add `internal/study/init.go` with focused `InitRequest`, `InitResult`, internal plan structures for directories/files/clone actions, and status fields needed by command output.
  - [x] Add typed or sentinel errors only where command mapping needs to distinguish usage, validation, workspace/filesystem, overwrite protection, and clone partial failure.
  - [x] Keep helper types private unless `internal/app` requires them for user output.

- [x] **Task 3: Implement YAML Parsing And Explicit Validation**
  > Executes: `Decision 4`, SR-AC-02 to SR-AC-07, SR-AC-13
  - [x] Add `internal/study/init_yaml.go` with structs for `name`, `description`, `repos.count`, `repos.items`, `dimensions.count`, and `dimensions.items`.
  - [x] Validate required fields with YAML field paths such as `repos.items[0].name` and `dimensions.items[0].number`.
  - [x] Validate source item `name`, `url`, `path`, and `description`; require unique filesystem-safe source names.
  - [x] Validate dimension item `number`, `name`, `title`, `description`, `purpose`, `steps`, `citations`, and `questions`; reject duplicate numbers and duplicate normalized slugs.
  - [x] Reject `repos.count` or `dimensions.count` lower than explicit item counts.
  - [x] Reject `repos.count` or `dimensions.count` greater than explicit item counts with guidance that assisted completion is deferred and explicit items are required now.

- [x] **Task 4: Implement Normalization And Plan Construction**
  > Executes: `Decision 3`, `Decision 4`, Security contract, SR-AC-08, SR-AC-09, SR-AC-11, SR-AC-13 to SR-AC-16
  - [x] Normalize dimension numbers to two digits for memory, filenames, headings, and normalized YAML.
  - [x] Normalize dimension slugs to kebab-case from dimension names and use filenames shaped `NN-slug.md`.
  - [x] Resolve default output as `studies/<study>` under the resolved workspace and resolve `--output <dir>` through existing workspace/path-safety behavior.
  - [x] Validate planned directories, files, and clone destinations cannot escape the selected safe output root.
  - [x] Build one plan containing `study-init.yml`, `README.md`, `dimensions/`, `sources/`, `reports/source/`, `reports/final/`, all dimension files, and URL-backed clone actions.
  - [x] Implement existing-directory protection by default and force eligibility checks before writes.

- [x] **Task 5: Implement Deterministic Artifact Rendering**
  > Executes: `Decision 6`, Documentation contract, SR-AC-10 to SR-AC-13
  - [x] Add `internal/study/init_render.go` for normalized `study-init.yml`, study `README.md`, and dimension Markdown.
  - [x] Render normalized YAML in stable item order without absolute paths unless an explicitly allowed custom output path must be preserved for diagnostics.
  - [x] Render README with study description, sources, dimensions, generated paths, and only currently supported next commands.
  - [x] Render dimension Markdown preserving purpose, steps, citations, and questions in human-editable sections.

- [x] **Task 6: Implement Scoped Execution, Dry-Run, Force, And Cloning**
  > Executes: `Decision 3`, `Decision 5`, SR-AC-14 to SR-AC-20, SR-CON-04 to SR-CON-09
  - [x] Execute the validated plan by creating directories and writing known generated files with deterministic permissions.
  - [x] Ensure `--dry-run` returns or prints the plan without creating, overwriting, deleting, or cloning anything.
  - [x] Ensure `--force` overwrites only known generated files inside the selected study directory and never deletes unknown paths or modifies paths outside that directory.
  - [x] Add `internal/study/init_clone.go` with a narrow clone runner seam and a production runner that invokes executable plus argument array equivalent to `git clone --depth 1 <url> <dest>`.
  - [x] Execute clone actions after artifact writes unless `--dry-run` or `--no-clone` is set.
  - [x] Aggregate clone failures so generated artifact success remains visible and command status is non-zero partial/failure according to existing conventions.

- [x] **Task 7: Wire CLI Command And Output**
  > Executes: `Decision 2`, CLI Surface contract, SR-AC-01, SR-AC-14, SR-AC-15, SR-AC-17, SR-AC-19, SR-AC-51
  - [x] Add `/home/antonioborgerees/coding/ultraplan-go/internal/app/study_init_commands.go` or extend the existing study command registration file if that is the project convention.
  - [x] Register `ultraplan study init <study-init.yml>` with help and `--dry-run`, `--force`, `--no-clone`, and `--output <dir>` flags.
  - [x] Reuse existing global `--workspace` behavior and command context.
  - [x] Print calm text output for success, dry-run plan, overwrite protection, validation failures, skipped clone actions, clone failures, and generated artifact locations.
  - [x] Map errors to existing exit-code conventions; use partial completion exit code if available, otherwise general non-zero failure with explicit partial wording.

- [x] **Task 8: Add Study Init Unit And Workflow Tests**
  > Executes: `Decision 3` through `Decision 6`, Testing contract, SR-AC-02 to SR-AC-20, SR-AC-52
  - [x] Add `/home/antonioborgerees/coding/ultraplan-go/internal/study/init_test.go` covering valid schema parsing and invalid YAML/missing field errors.
  - [x] Cover duplicate source names, duplicate dimension numbers/slugs, unsafe names, unsafe output paths, and count mismatch failures.
  - [x] Cover two-digit dimension number normalization, kebab-case slug normalization, filenames, headings, and normalized YAML.
  - [x] Cover generated README and dimension Markdown content with fixture/golden-style assertions where stable output matters.
  - [x] Cover successful directory/file creation in a temp workspace.
  - [x] Cover dry-run non-mutation and no clone execution.
  - [x] Cover existing study protection, force overwrite behavior, outside sentinel preservation, and unknown file preservation.
  - [x] Cover fake clone runner args, `--no-clone` skip behavior, dry-run no-exec behavior, and partial clone failure results.

- [x] **Task 9: Add CLI Command Tests**
  > Executes: `Decision 2`, Testing contract, CLI Surface contract, SR-AC-01, SR-AC-14, SR-AC-15, SR-AC-17, SR-AC-19, SR-AC-51
  - [x] Add `/home/antonioborgerees/coding/ultraplan-go/internal/app/study_init_commands_test.go` or extend existing command tests if that is the project convention.
  - [x] Cover help output and flag presence.
  - [x] Cover usage errors for missing or extra arguments.
  - [x] Cover dry-run output with directories, files, and clone actions.
  - [x] Cover overwrite protection and `--force` command behavior.
  - [x] Cover `--workspace` and `--output` path behavior.
  - [x] Cover clone partial failure output and exit mapping using a fake or injectable runner.

- [x] **Task 10: Review Scope, Boundaries, And Verification Gates**
  > Executes: `Decision 7`, SR-NG-01 to SR-NG-07, Architecture Review, Sprint Review
  - [x] Review imports to confirm `internal/app -> internal/study`, `internal/study -> internal/workspace` is allowed, and platform packages do not import `internal/study`.
  - [x] Review code and generated README/help text to confirm no runtime-assisted completion, OpenCode/agentwrap execution, Markdown frontmatter/applicability, prompt composition, run state, target/sprint workflows, diagnostics persistence, or new stable JSON schemas were added.
  - [x] Run `go test ./...` from `/home/antonioborgerees/coding/ultraplan-go`.
  - [x] Run `go build ./cmd/ultraplan` from `/home/antonioborgerees/coding/ultraplan-go`.
  - [x] Capture deviations, omitted evidence, or unresolved open questions for sprint review.

## Evidence Checklist

- [x] Tests prove YAML parsing, validation, normalization, generated artifacts, dry-run, force safety, clone behavior, CLI output, and exit-code behavior.
- [x] Runtime or diagnostic evidence exists where required: no agent runtime evidence is required; clone behavior is evidenced by fake runner tests and non-shell argument review.
- [x] Documentation updates are complete where required: generated README fixture/help output accurately describe current capabilities only.
- [x] Deviations from `reasoning.md` are recorded before implementation continues.
- [x] Required review protocols have evidence: Architecture Review and Sprint Review checks can be run from implementation diff, tests, and verification commands.

## Verification Commands

| Check | Command | Expected Result |
| --- | --- | --- |
| Unit and command tests | `go test ./...` | Passes from `/home/antonioborgerees/coding/ultraplan-go` without network, GitHub, OpenCode, provider credentials, or real cloned repositories. |
| CLI build | `go build ./cmd/ultraplan` | Builds the `ultraplan` command from `/home/antonioborgerees/coding/ultraplan-go`. |
| Help smoke check | `go run ./cmd/ultraplan study init --help` | Shows `study init`, required argument, and `--dry-run`, `--force`, `--no-clone`, and `--output` flags if command tests do not already cover the exact invocation. |
| Manual dry-run fixture check | `go run ./cmd/ultraplan --workspace <fixture-workspace> study init <fixture-study-init.yml> --dry-run` | Prints planned directories, files, and clone actions without creating, overwriting, deleting, or cloning anything, if a manual fixture is used during review. |

## Risks And Blockers

| Risk / Blocker | Source | Mitigation | Status |
| --- | --- | --- | --- |
| Existing workspace APIs may not fully cover default and `--output` path-safety needs. | `reasoning.md#assumptions-and-risks` | Used exported `workspace.ResolveInside`; kept `--output` workspace-scoped. | closed |
| Partial clone failure exit-code mapping may not have an exact existing convention. | `reasoning.md#assumptions-and-risks`; `reasoning/architecture.md#risks-assumptions-and-open-questions` | Mapped clone partial failures to existing `ExitPartial`. | closed |
| `--output <dir>` outside-workspace behavior is subtle. | `reasoning.md#assumptions-and-risks`; `reasoning/architecture.md#risks-assumptions-and-open-questions` | Conservative behavior: outside-workspace output fails validation. | closed |
| Users may expect `--force` to replace stale clone directories. | `reasoning.md#assumptions-and-risks` | Force overwrites known generated files only; clone destination conflicts are surfaced by clone failure. | carried forward |
| User-supplied clone URLs may contain credentials. | `reasoning.md#assumptions-and-risks` | Output prints clone names and destinations, not full command strings. Normalized YAML preserves user-provided metadata. | mitigated |
| Normalized YAML formatting may become a de facto compatibility surface. | `reasoning.md#assumptions-and-risks` | Kept deterministic formatting but added no public schema guarantee. | mitigated |
| README supported-next-command guidance may drift in future sprints. | `reasoning.md#assumptions-and-risks` | README lists only current study list commands and is covered by tests. | mitigated |
| Sequential clones may be slow for many repositories. | `reasoning.md#assumptions-and-risks` | Accepted for Sprint 4; no concurrency requirement. | carried forward |

## Review Inputs

Review should use:

- `projects/ultraplan-go/sprints/04-study-init-yaml/sprint-index.md`
- `projects/ultraplan-go/sprints/04-study-init-yaml/technical-handbook.md`
- `projects/ultraplan-go/sprints/04-study-init-yaml/reasoning/architecture.md`
- `projects/ultraplan-go/sprints/04-study-init-yaml/reasoning.md`
- `projects/ultraplan-go/sprints/04-study-init-yaml/plan.md`
- implementation diff
- verification evidence from `go test ./...` and `go build ./cmd/ultraplan`
- `system/protocols/architecture-review-protocol.md`
- `system/protocols/sprint-review-protocol.md`

## Execution Log

| Date / Step | Action | Evidence / Notes |
| --- | --- | --- |
| 2026-05-30 / planning | Created implementation plan from `reasoning.md` and selected sprint evidence. | No implementation performed. |
| 2026-05-30 / implementation | Implemented YAML study initialization, deterministic rendering, safe dry-run/force behavior, optional non-shell shallow cloning, and CLI wiring. | Changed `internal/study/init*.go`, `internal/app/study_commands.go`; added tests in `internal/study/init_test.go` and `internal/app/study_init_commands_test.go`; added `gopkg.in/yaml.v3`. |
| 2026-05-30 / verification | Ran sprint verification commands. | `go test ./...` passed; `go build ./cmd/ultraplan` passed; `go run ./cmd/ultraplan study init --help` passed. |

## Completion Criteria

- [x] All tasks are complete or explicitly deferred with requirement and decision references.
- [x] Verification commands were run or deferrals are documented with cause.
- [x] Evidence satisfies the expectations from `reasoning.md`.
- [x] `review.md` can evaluate conformance without guessing intent.
- [x] No deferred runtime, target, sprint, Markdown applicability, diagnostics persistence, or stable JSON-schema behavior was introduced.
