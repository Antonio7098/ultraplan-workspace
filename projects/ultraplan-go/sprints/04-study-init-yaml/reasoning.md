# Sprint Reasoning: Study Initialization From YAML

> Project: `ultraplan-go`
> Sprint: `04-study-init-yaml`
> Output: `projects/ultraplan-go/sprints/04-study-init-yaml/reasoning.md`
> **Inputs Used:** `projects/ultraplan-go/project-index.md`, `projects/ultraplan-go/sprints/04-study-init-yaml/requirements.md`, `projects/ultraplan-go/docs/PRD.md`, `projects/ultraplan-go/docs/TRD.md`, `projects/ultraplan-go/docs/ARCHITECTURE.md`, `projects/ultraplan-go/sprints/04-study-init-yaml/sprint-index.md`, `projects/ultraplan-go/sprints/04-study-init-yaml/technical-handbook.md`, `projects/ultraplan-go/sprints/04-study-init-yaml/reasoning/architecture.md`, `templates/sprint-reasoning.md`

This document decides. It synthesizes selected context, handbook evidence, area-specific reasoning, and contracts into final sprint decisions.

It does not replace `sprint-index.md`, `technical-handbook.md`, or `reasoning/*.md`.

## Sprint Purpose

- **Goal:** Implement `ultraplan study init <study-init.yml>` so explicit YAML study definitions generate deterministic study directories, normalized YAML, editable dimension files, a study README, optional shallow source clones, and safe dry-run/force behavior.
- **Non-Goals:** Do not implement runtime-assisted completion, runtime suggestion prompts, suggestion caches, `--no-assist`, source URL verification, replacement source suggestions, OpenCode/agentwrap runtime execution, Markdown document source discovery, frontmatter parsing, applicability filtering, prompt composition, analysis runs, synthesis, report validation, run state, retries, cancellation, diagnostics persistence, code-reference extraction, target workflows, sprint workflows, or new stable public JSON schemas.
- **Depends On:** Sprint 1 CLI shell and exit conventions, Sprint 2 workspace discovery/path safety/global `--workspace`, Sprint 3 study domain/listing compatibility, and the PRD/TRD/Architecture docs.

## Selected Context And Pre-Reasoning Artifacts

| Artifact | Path | How It Was Used |
| --- | --- | --- |
| Project Index | `projects/ultraplan-go/project-index.md` | Confirmed project scope, target implementation directory, selected source documents, available contracts, evidence report pool, and lack of prior decision artifacts. |
| Requirements | `projects/ultraplan-go/sprints/04-study-init-yaml/requirements.md` | Treated as the sprint contract for outputs, acceptance criteria, non-goals, constraints, dependencies, and review expectations. |
| PRD | `projects/ultraplan-go/docs/PRD.md` | Used for product context: YAML study initialization, generated folders, dry-run/force/no-clone expectations, shallow cloning, user-editable artifacts, and deferred target/sprint/runtime workflows. |
| TRD | `projects/ultraplan-go/docs/TRD.md` | Used for technical context: module-driven architecture, study init schema, generated structure, clone behavior, force behavior, path handling, exit-code classes, testing, and security requirements. |
| Architecture | `projects/ultraplan-go/docs/ARCHITECTURE.md` | Used as the package-boundary source of truth: product behavior belongs in `internal/study`; CLI/composition belongs in `internal/app`; `study -> workspace` is allowed; platform packages must not import product modules. |
| Sprint Index | `projects/ultraplan-go/sprints/04-study-init-yaml/sprint-index.md` | Selected contracts, evidence reports, reasoning artifacts, review protocols, excluded context, and non-goals for this sprint. |
| Technical Handbook | `projects/ultraplan-go/sprints/04-study-init-yaml/technical-handbook.md` | Supplied evidence patterns for thin CLI handlers, explicit dependencies, narrow process seams, plan/dry-run consistency, actionable errors, non-shell subprocess calls, fixture/fake tests, deterministic artifacts, and scope discipline. |
| Architecture Reasoning | `projects/ultraplan-go/sprints/04-study-init-yaml/reasoning/architecture.md` | Accepted as the area-specific decision input for module ownership, plan-then-execute flow, clone runner seam, force/dry-run mutation boundaries, and rejected alternatives. |

## Area-Specific Reasoning Inputs

| Area | Reasoning Document | Key Conclusion | Evidence Basis | Impact On Final Decision |
| --- | --- | --- | --- | --- |
| Architecture | `projects/ultraplan-go/sprints/04-study-init-yaml/reasoning/architecture.md` | Implement a focused `internal/study` initialization use case with plan-then-execute, deterministic renderers, scoped filesystem mutation, and a narrow fakeable clone runner. Keep `internal/app` limited to command wiring/output and do not add global technical-layer packages. | `technical-handbook.md` findings from `01-project-structure`, `02-command-architecture`, `03-dependency-injection`, `05-error-handling`, `06-io-abstraction`, `11-testing-strategy`, `13-security`, and `15-philosophy`; Architecture doc module ownership and dependency rules. | Final decisions adopt the `internal/study` ownership model, one shared plan for dry-run/execution, non-shell clone runner seam, typed/actionable errors, known-file-only force behavior, sequential clone execution, and offline fixture/fake test strategy. |

## Sprint Technical Handbook Summary

- **Relevant Patterns:** Thin CLI wiring with product logic elsewhere; factory/manual dependency injection; narrow interfaces at volatile boundaries; plan-then-execute for dry-run and force safety; actionable classified errors; non-shell command execution; fake-first command/filesystem tests; human-editable generated artifacts.
- **Important Trade-Offs:** An explicit plan/result model adds internal structure but keeps dry-run, force checks, clone actions, and command output aligned. A minimal clone runner seam avoids broad process abstractions. Keeping init files in one `internal/study` package avoids premature subpackages but requires naming discipline. Golden/fixture assertions improve determinism but can create maintenance work.
- **Warnings / Anti-Patterns:** Do not put validation/rendering/filesystem policy in `RunE`; do not introduce global validation/reports/workflow packages; do not build shell command strings for clone; do not panic on user YAML; do not let dry-run bypass validation; do not require network/GitHub/OpenCode in normal tests; do not generate non-deterministic or unreviewable artifacts.
- **Evidence Confidence:** High for architecture, command, DI, IO, error handling, security, testing, and philosophy evidence because the selected reports cite mature Go CLI examples. Medium for terminal UX because this sprint needs calm text output but not interactive prompts or progress UI.

## Contracts Applied

| Contract / Requirement ID | Constraint | Decision Impact | Expected Evidence |
| --- | --- | --- | --- |
| Architecture contract; SR-CON-01 to SR-CON-03 | Study initialization behavior stays in `internal/study`; CLI parsing/output stays in `internal/app`; dependency direction remains `app -> study -> workspace/platform`; platform must not import `study`. | Use `internal/study/init.go`, `init_yaml.go`, `init_render.go`, and `init_clone.go` for behavior; use `internal/app/study_init_commands.go` only for command registration, flags, output, and exit mapping. | Import review, architecture review protocol, `go test ./...`, and command tests showing CLI delegates to study behavior. |
| Errors contract; SR-AC-01, SR-AC-05 to SR-AC-07, SR-AC-19 | Invalid usage, YAML, duplicate names, unsafe paths, unsupported count shortages, overwrite protection, and clone failures must return actionable non-zero errors while preserving cause context. | Add field/path/source/dimension-specific validation errors and a partial clone failure representation that maps to existing exit-code conventions. | Unit tests for invalid YAML/fields/counts/duplicates/paths; command tests for usage and exit-code mapping; clone failure test showing generated artifacts remain visible. |
| Security contract; SR-AC-03, SR-AC-05, SR-AC-14 to SR-AC-18, SR-CON-04 to SR-CON-09 | User paths and names are untrusted; writes must stay under the resolved workspace or explicitly selected safe output root; `--force` is scoped; clone uses argument arrays and no shell interpolation; artifacts must not contain secrets. | Normalize and validate study/source/dimension names, output roots, generated paths, and clone destinations before mutation; use `git clone --depth 1 <url> <dest>` via a testable non-shell runner; overwrite known generated files only. | Path escape tests, unsafe name tests, force safety tests with outside sentinels, fake clone arg assertions, source review for no `sh -c`, artifact content review. |
| Testing contract; SR-AC-20, SR-AC-52, SR-AC-53 | Normal tests must be deterministic and offline; build/test gates must pass. | Use temp workspaces, fixture YAML, fake clone runners, command IO capture, and generated-content assertions; no network, GitHub, provider credentials, OpenCode, or real cloned repositories in normal tests. | `go test ./...`, `go build ./cmd/ultraplan`, unit tests in `internal/study/init_test.go`, command tests in `internal/app/study_init_commands_test.go`. |
| Documentation contract; SR-AC-10 to SR-AC-13 | Generated README, normalized YAML, and dimension Markdown must be deterministic, human-editable, and not advertise deferred workflows as implemented. | Render README with study description, sources, dimensions, generated paths, and only currently supported next commands; render normalized YAML and Markdown in stable order. | Golden or fixture comparisons for README/YAML/dimension files; review check that deferred runtime/target/sprint workflows are not advertised as available. |
| CLI Surface contract; SR-AC-01, SR-AC-14, SR-AC-15, SR-AC-17, SR-AC-19, SR-AC-51 | `ultraplan study init <study-init.yml>` must have help, flags, calm dry-run/success/failure text, global workspace behavior, and meaningful exit codes. | Register command and flags in `internal/app`; use existing workspace resolution; print directories/files/clone actions for dry-run; map partial clone failure to partial completion if available, otherwise the closest existing failure code with clear text. | Help output tests, dry-run output tests, overwrite protection tests, workspace flag tests, partial failure command test. |
| SR-AC-02 to SR-AC-04 | YAML must parse `name`, `description`, `repos.count`, `repos.items`, `dimensions.count`, and `dimensions.items`; source and dimension item fields must be supported. | Define explicit init YAML structs and normalize into internal source/dimension models. | Parser tests for valid full schema and required item fields. |
| SR-AC-08 to SR-AC-10 | Dimension numbers and slugs must be normalized; filenames are `NN-slug.md`; Markdown preserves purpose, steps, citations, and questions. | Normalize numbers to two digits and names to kebab-case once during planning; render generated YAML, filenames, headings, and Markdown from normalized values. | Unit tests for normalization and generated dimension Markdown. |
| SR-AC-11 to SR-AC-13 | Successful init creates required structure and deterministic artifacts under workspace/default or `--output`; normalized YAML avoids absolute paths except explicit diagnostics. | Plan directories/files first, then write stable generated files; preserve only explicitly selected output diagnostics where necessary. | Workflow tests inspecting directory tree and artifact contents. |
| SR-AC-16 to SR-AC-19 | Clone behavior is optional, shallow, non-shell, and reports partial failures without hiding artifact success. | Execute clone actions after file writes, sequentially, through fakeable runner; skip with `--no-clone`; aggregate clone results. | Fake clone runner tests for args/order/skip/failure; command output review. |

## Repos Studied / Source Evidence Used

| Source / Repo / Report | Concrete Reference | Relevant Finding | Why It Matters For This Sprint | Used In Decision(s) |
| --- | --- | --- | --- | --- |
| `01-project-structure` | `studies/go-cli-study/reports/final/01-project-structure.md`; cited `cmd/restic/main.go:37-114`, `cmd/root.go:112-127`, `cmd/gh/main.go:6` | Mature Go CLIs keep entrypoints thin and put behavior behind `cmd/` or `internal/` boundaries. | Supports keeping `cmd/ultraplan` and `internal/app` thin while `internal/study` owns init logic. | Decision 1, Decision 2 |
| `02-command-architecture` | `studies/go-cli-study/reports/final/02-command-architecture.md`; cited `pkg/cmd/install.go:132-145`, `pkg/cmd/issue/list/list.go:47-118`, `cmd/restic/cmd_backup.go:84-115` | Factory-created commands with thin handlers delegate to action logic. | Shapes command registration, flag collection, and delegation from app to study service. | Decision 2 |
| `03-dependency-injection` | `studies/go-cli-study/reports/final/03-dependency-injection.md`; cited `pkg/cmdutil/factory.go:16-43`, `executor.go:22-24`, `internal/app/app.go:42-81` | Manual DI through factories and explicit service boundaries is preferable to globals. | Supports passing workspace root, output options, IO, and clone runner explicitly. | Decision 1, Decision 5 |
| `05-error-handling` | `studies/go-cli-study/reports/final/05-error-handling.md`; cited `cmd/age/tui.go:37-54`, `go-task/errors/errors_task.go:13-32`, `internal/ghcmd/cmd.go:44-49`, `pkg/action/uninstall.go:232-254` | CLI errors should wrap causes, classify failures, include hints, and aggregate partial failures. | Drives field-specific validation errors and clone partial-failure reporting. | Decision 4, Decision 5 |
| `06-io-abstraction` | `studies/go-cli-study/reports/final/06-io-abstraction.md`; cited `pkg/iostreams/iostreams.go:551-568`, `internal/fs/interface.go:10-31`, `executor.go:553-564` | Injectable IO/filesystem/process seams improve command and file-operation tests. | Supports fake clone runner and command IO capture without over-abstracting pure helpers. | Decision 5, Decision 6 |
| `09-terminal-ux` | `studies/go-cli-study/reports/final/09-terminal-ux.md`; cited `internal/cmd/prompt.go:20-256`, `lib/terminal/terminal.go:82-86`, `internal/ui/terminal.go:10-34` | CLI output should be calm, non-TTY-safe, and avoid unnecessary interactive UX. | Shapes success, dry-run, and partial-failure text; rejects spinners/TUI/progress in this sprint. | Decision 2, Decision 6 |
| `11-testing-strategy` | `studies/go-cli-study/reports/final/11-testing-strategy.md`; cited `internal/cmd/main_test.go:64-174`, `pkg/httpmock/stub.go:35-199`, `fake_cmd_obj_runner.go:17-26`, `helm/internal/test/test.go:43`, `go-task/task_test.go:166-169` | Strong CLI projects use table-driven tests, command tests, fakes, fixtures, and golden-style assertions. | Drives parser/service/command tests, fake clone tests, and deterministic artifact comparisons. | Decision 6 |
| `13-security` | `studies/go-cli-study/reports/final/13-security.md`; cited `internal/execext/exec.go:59-66`, `cmd_obj_builder.go:38`, `pkg/surveyext/editor_manual.go:23`, `pkg/registry/transport.go:37-41` | Path-sensitive operations and subprocess execution need trust boundaries and argument arrays. | Drives path safety, no shell interpolation, explicit clone args, and secret-free output. | Decision 3, Decision 5 |
| `15-philosophy` | `studies/go-cli-study/reports/final/15-philosophy.md`; cited `internal/chezmoi/system.go:25-45`, `fs/types.go:16-59`, `internal/backend/backend.go:19-90`, `VISION.md:97` | Durable CLIs reject scope, keep complexity near benefits, and use wrappers/interfaces where dry-run/test/backend flexibility requires them. | Supports plan-then-execute, no assisted completion, no broad abstractions, and human-editable artifacts. | Decision 1, Decision 3, Decision 7 |

## Trade-Off And Debt Analysis

### Accepted Trade-Offs

| Trade-Off | Benefit | Cost / Constraint Accepted | Why Acceptable Now | Revisit Trigger |
| --- | --- | --- | --- | --- |
| Explicit init plan/result model before execution | One validated source of truth for dry-run, writes, force checks, clone actions, and command output. | Adds internal structs and possible duplication between plan and render concepts. | Sprint has dry-run, force safety, clone, and partial failure requirements that need consistent behavior. | Multiple product modules need the same planning/execution abstraction or the plan becomes hard to understand. |
| Keep initialization in one `internal/study` package split by focused files | Preserves module ownership and avoids premature global layers. | `internal/study` grows and needs disciplined file naming. | Architecture and requirements explicitly constrain study behavior to `internal/study`. | Study package becomes too large to navigate or dependency boundaries within study become unclear. |
| Minimal clone runner interface only | Keeps subprocess boundary fakeable and secure without broad process framework. | Later external commands may require a more general platform process helper. | This sprint only executes local `git clone`; runtime execution remains deferred. | A future sprint adds multiple non-runtime subprocess use cases. |
| Sequential clone execution | Simpler partial-failure reporting and deterministic tests. | Slower when many URL-backed sources are initialized. | Initialization volume is expected to be small/moderate and performance evidence was excluded. | Requirements add clone concurrency, many-source initialization, or progress/throughput targets. |
| Force overwrites known generated files only and does not delete unknown paths | Minimizes accidental loss of user edits or unrelated files. | Existing clone directories may remain if stale, and users may need manual cleanup. | Requirements require scoped force behavior but do not require destructive cleanup. | A later sprint defines explicit clean/replace source behavior with backup or confirmation semantics. |
| Human text output only for this sprint | Keeps CLI behavior simple and matches no-new-stable-JSON-schemas non-goal. | Automation cannot rely on a structured study-init JSON result yet. | Requirements explicitly avoid adding stable public JSON schemas beyond existing helpers. | A future CLI contract requires JSON output for init or release compatibility. |

### Potential Technical Debt

| Debt / Shortcut | Why It Might Accrue | Current Mitigation | Owner / Follow-Up |
| --- | --- | --- | --- |
| Init plan shape could become a generic workflow engine | Future run/synthesis/state work may try to reuse it outside initialization. | Keep names and types init-specific; do not export more than command output requires. | Future workflow/run-loop sprint should decide its own orchestration model. |
| YAML normalization may become schema compatibility surface accidentally | Users may start relying on exact normalized YAML formatting. | Keep deterministic tests and do not claim a public stable schema beyond sprint requirements. | Release/schema sprint should decide compatibility guarantees. |
| Clone runner may duplicate future platform process helper | Later commands may also need subprocess execution. | Keep interface narrow and easy to replace; use argument arrays now. | Future subprocess-heavy sprint can extract a platform process helper if a second real consumer exists. |
| README supported-next-command list may drift as new commands are implemented | Generated docs can become stale if not updated with command surface. | Test for current sprint wording and avoid advertising deferred workflows. | Each command sprint should update README rendering intentionally. |
| Force behavior may not clean stale generated directories | Known-file-only overwrites preserve safety but may leave obsolete clone directories or files. | Document outcome in command text; do not delete unknown files. | Add explicit `--clean` or backup behavior only if requirements later demand it. |

### Future Considerations

| Consideration | Deferred Until | Reason Deferred | What Should Be Preserved Now |
| --- | --- | --- | --- |
| Runtime-assisted source/dimension completion and suggestion caches | Assisted completion sprint | Sprint 4 explicitly fails count shortages instead of invoking runtime. | Validation errors must clearly explain explicit items are required now and avoid reserving incompatible YAML fields. |
| Source URL verification and replacement suggestions | Source acquisition/repair sprint | Requirements exclude URL verification and repair prompts. | Clone action model should not imply verification was performed. |
| Markdown document source discovery and applicability | Sprint 5 | Requirements explicitly defer Markdown discovery/frontmatter/applicability. | Do not recursively scan sources during init; preserve top-level source layout compatibility. |
| Study run/synthesis/report validation/state | Runtime execution sprints | Prompt composition, run state, validation, retry, and summaries are non-goals. | Generated dimensions/reports directories should match future discovery expectations. |
| Stable JSON output for study init | Public automation/release sprint | No stable public schema is required now. | Internal result should be structured enough to support future JSON without changing domain behavior. |
| Cross-platform filesystem nuances | Release hardening | First target is Linux/macOS; Windows should not be intentionally blocked. | Use standard path/file APIs and avoid shell-specific behavior. |

## Final Decisions

### Decision 1: Study-Owned Init Use Case

- **Decision:** Implement study initialization as a focused use case in `internal/study`, split across `init.go`, `init_yaml.go`, `init_render.go`, and `init_clone.go`. Expose a small request/result API for `internal/app`; keep parser, validation, normalization, rendering, planning, filesystem mutation, and clone policy outside CLI handlers.
- **Rationale:** Study initialization transforms study definitions, sources, dimensions, and study filesystem artifacts. The project Architecture says product modules own product behavior and `study` owns study definitions, validation, reports, and persistence. Keeping this logic in `internal/study` prevents `internal/app` from becoming a business-rule layer and preserves Sprint 3 listing compatibility.
- **Study / Source Grounding:** `technical-handbook.md` cites thin entrypoints from `01-project-structure` (`cmd/restic/main.go:37-114`, `cmd/root.go:112-127`, `cmd/gh/main.go:6`) and command delegation from `02-command-architecture` (`pkg/cmd/install.go:132-145`, `pkg/cmd/issue/list/list.go:47-118`). `03-dependency-injection` supports explicit factories (`pkg/cmdutil/factory.go:16-43`). Architecture reasoning chose this option and rejected CLI-owned logic and global technical packages.
- **Trade-Offs Accepted:** Accept several focused files in one package and a small init request/result model rather than a generic clean-architecture tree or global validation/report packages.
- **Technical Debt / Future Impact:** The `internal/study` package can grow; future sprints may split subpackages only if a real comprehension or dependency benefit appears.
- **Alternatives Rejected:** Rejected putting YAML validation/rendering/writes in `internal/app` because it violates sprint constraints and command-architecture evidence. Rejected `internal/validation`, `internal/reports`, `internal/workflow`, or filesystem-planner packages because Architecture warns against global technical layers for product behavior.
- **Contracts Satisfied:** Architecture contract; SR-CON-01 to SR-CON-03; SR-AC-01 to SR-AC-13; SR-AC-51.
- **Evidence Required:** Import review showing allowed dependency direction, unit tests through `internal/study`, command tests showing app delegation, architecture review protocol, `go test ./...`, and `go build ./cmd/ultraplan`.

### Decision 2: Thin CLI Command With Existing Workspace Behavior

- **Decision:** Add `ultraplan study init <study-init.yml>` in `internal/app` with `--dry-run`, `--force`, `--no-clone`, and `--output <dir>` flags. The command resolves the workspace using existing global behavior, constructs `study.InitRequest`, delegates to `internal/study`, prints calm text output, and maps errors to existing exit-code conventions.
- **Rationale:** CLI handlers should own command shape, flags, help, output, and exit mapping, not study business rules. This also preserves the global `--workspace` behavior required from prior sprints.
- **Study / Source Grounding:** `technical-handbook.md` cites command factories and thin handlers in `02-command-architecture`, manual DI from `03-dependency-injection`, and calm non-TTY-safe output from `09-terminal-ux` (`internal/cmd/prompt.go:20-256`, `lib/terminal/terminal.go:82-86`, `internal/ui/terminal.go:10-34`).
- **Trade-Offs Accepted:** Accept human-readable text as the only required output mode for this sprint and avoid stable public JSON schemas.
- **Technical Debt / Future Impact:** Future JSON output can be added from structured `InitResult` without changing study behavior if a later contract requires it.
- **Alternatives Rejected:** Rejected interactive prompts/spinners/TUI progress because init is a batch command and terminal UX evidence favors calm script-safe output. Rejected embedding workspace discovery/path safety directly in the command because existing workspace behavior must be reused.
- **Contracts Satisfied:** CLI Surface contract; Documentation contract; Errors contract; SR-AC-01, SR-AC-14, SR-AC-15, SR-AC-17, SR-AC-19, SR-AC-51.
- **Evidence Required:** Help output tests, usage error tests, dry-run output tests, global `--workspace` tests, overwrite protection tests, partial clone exit/output tests, and review that no deferred workflows are advertised.

### Decision 3: Single Validated Plan For Dry Run, Writes, Force, And Path Safety

- **Decision:** Build a deterministic initialization plan after parsing and validation. The plan lists all directories, files, and clone actions. `--dry-run` prints this plan without mutation; real execution applies the same plan. All planned paths are normalized and checked to remain under the resolved workspace study root or explicitly selected safe output root. `--force` overwrites known generated files inside the selected study directory only and never deletes unknown files or touches paths outside that directory.
- **Rationale:** A shared plan prevents dry-run drift, makes path safety reviewable before mutation, and scopes force behavior precisely. It also makes command output and tests concrete.
- **Study / Source Grounding:** `technical-handbook.md` identifies plan-then-execute as a relevant pattern using `15-philosophy` and `internal/chezmoi/system.go:25-45`. `13-security` supports explicit path trust boundaries. Architecture reasoning concluded the plan is the single source of truth for dry-run output and execution.
- **Trade-Offs Accepted:** Accept an internal plan model and known-file-only overwrite behavior. This may leave stale unknown files or old clone directories, but it avoids accidental data loss.
- **Technical Debt / Future Impact:** Future explicit clean/backup behavior may be needed if users want source directory replacement; this sprint must not invent destructive cleanup semantics.
- **Alternatives Rejected:** Rejected separate dry-run logic because it can bypass validation and path checks. Rejected rollback/delete-on-clone-failure because requirements say clone failures must not hide successful artifact creation. Rejected force deletion of whole study directories because it risks user files and violates scoped overwrite constraints.
- **Contracts Satisfied:** Security contract; Architecture contract; SR-AC-11 to SR-AC-16; SR-CON-04 to SR-CON-08.
- **Evidence Required:** Dry-run non-mutation tests, path escape tests, output-root tests, existing-directory failure tests, force tests with outside/unknown-file sentinels, deterministic plan/artifact assertions, and security review.

### Decision 4: Explicit YAML-Only Initialization With Deterministic Normalization

- **Decision:** Parse only explicit YAML inputs for this sprint. Required fields are `name`, `description`, `repos.count`, `repos.items`, `dimensions.count`, and `dimensions.items`. Source items support `name`, `url`, `path`, and `description`. Dimension items support `number`, `name`, `title`, `description`, `purpose`, `steps`, `citations`, and `questions`. Dimension numbers normalize to two digits and dimension slugs normalize to kebab-case for memory, filenames, headings, and normalized YAML. Count values lower than explicit items fail, and count values greater than explicit items fail with guidance that assisted completion is deferred.
- **Rationale:** The sprint goal is deterministic explicit initialization, not runtime-assisted completion. Rejecting shortages now avoids silently producing incomplete studies or invoking runtime behavior outside scope.
- **Study / Source Grounding:** `technical-handbook.md` warnings against panics for user YAML and unreviewable generated output are grounded in `05-error-handling` and `15-philosophy`. The PRD/TRD define broader assisted completion, but `requirements.md` explicitly narrows Sprint 4 to explicit YAML and failure guidance for count shortages.
- **Trade-Offs Accepted:** Users must provide all sources and dimensions explicitly until assisted completion is implemented. This is stricter than the long-term PRD/TRD capability but matches the sprint contract.
- **Technical Debt / Future Impact:** Future assisted completion must add suggestion cache/merge behavior without weakening current validation. Normalized YAML may become user-visible; keep ordering deterministic and reviewable.
- **Alternatives Rejected:** Rejected runtime-assisted suggestions and `--no-assist` because they are explicit non-goals. Rejected auto-filling placeholder dimensions/sources because generated placeholders would undermine deterministic, reviewable artifacts. Rejected accepting unsafe names with best-effort slugging when uniqueness or filesystem safety is ambiguous.
- **Contracts Satisfied:** Errors contract; Documentation contract; SR-AC-02 to SR-AC-10; SR-AC-13; SR-NG-01; SR-CON-09.
- **Evidence Required:** Parser tests for valid schema; validation tests for missing fields, invalid YAML, duplicate source names, duplicate dimension numbers/slugs, unsafe names, unsafe output paths, and count mismatches; generated YAML/Markdown fixture tests.

### Decision 5: Optional Shallow Clone Through Narrow Non-Shell Runner

- **Decision:** For URL-backed sources, generate clone actions targeting `sources/<source-name>` and execute them after file writes unless `--dry-run` or `--no-clone` is set. Use a narrow fakeable runner that invokes `git clone --depth 1 <url> <dest>` through executable plus argument array, never shell interpolation. Run clone actions sequentially. Aggregate clone results so generated artifact success remains visible and any clone failure returns a non-zero partial/failure status according to existing CLI conventions.
- **Rationale:** Clone execution is an external process and possible network operation, so it needs a test seam and strict argument handling. Running after artifact writes satisfies the requirement to report clone failures without hiding successful generated artifacts. Sequential execution keeps results deterministic and easier to test.
- **Study / Source Grounding:** `technical-handbook.md` cites non-shell command execution from `13-security` (`internal/execext/exec.go:59-66`, `cmd_obj_builder.go:38`) and fake command runners from `11-testing-strategy` (`fake_cmd_obj_runner.go:17-26`). `05-error-handling` multi-error evidence supports partial failure aggregation (`pkg/action/uninstall.go:232-254`, `internal/checker/checker.go:70-91`).
- **Trade-Offs Accepted:** Sequential clones may be slower, and the clone runner is intentionally not a broad process abstraction. That is acceptable because no performance or generic subprocess requirement exists in this sprint.
- **Technical Debt / Future Impact:** If later commands need subprocess execution, the runner can be replaced by a platform process helper. If clone concurrency is needed, the result model should preserve deterministic output order.
- **Alternatives Rejected:** Rejected `sh -c` or string-built clone commands due to injection risk. Rejected network/GitHub-backed tests because normal tests must be offline. Rejected URL verification and replacement suggestions because they are non-goals.
- **Contracts Satisfied:** Security contract; Testing contract; Errors contract; SR-AC-16 to SR-AC-20; SR-CON-09.
- **Evidence Required:** Fake clone tests asserting executable/args/destination, `--no-clone` skip tests, dry-run no-exec tests, clone failure partial result tests, source review for no shell interpolation, and offline test execution.

### Decision 6: Deterministic Generated Artifacts And Reviewable Tests

- **Decision:** Render deterministic `study-init.yml`, `README.md`, and `dimensions/NN-slug.md` files. Generated dimension Markdown must preserve supplied purpose, steps, citations, and questions in human-editable form. README must describe the study, sources, dimensions, generated paths, and currently supported next commands only. Tests should use fixture/golden-style assertions where generated content stability matters and direct structured assertions where behavior matters.
- **Rationale:** The product requires local, reviewable, editable artifacts. Deterministic rendering keeps generated studies discoverable and version-control friendly, while targeted golden tests catch accidental wording/order regressions without testing private helper implementation details.
- **Study / Source Grounding:** `technical-handbook.md` cites human-editable generated artifacts from `15-philosophy`, golden/fixture tests from `11-testing-strategy` (`helm/internal/test/test.go:43`, `go-task/task_test.go:166-169`), and calm output guidance from `09-terminal-ux`.
- **Trade-Offs Accepted:** Golden-style tests may require updates when intended text changes. This is acceptable for artifact formats users will review and edit.
- **Technical Debt / Future Impact:** README wording must evolve as future study commands are implemented, but this sprint must not advertise deferred runtime, target, or sprint workflows as available.
- **Alternatives Rejected:** Rejected non-deterministic map-order YAML or prose generation because it undermines reviewability. Rejected opaque machine-only artifacts because PRD requires readable/editable files. Rejected broad substring-only tests for generated artifacts because they can miss regressions in required sections.
- **Contracts Satisfied:** Documentation contract; Testing contract; SR-AC-08 to SR-AC-13; SR-AC-20; SR-AC-52; SR-AC-53.
- **Evidence Required:** Fixture/golden comparisons for normalized YAML, README, and dimension Markdown; study listing compatibility checks; `go test ./...`; `go build ./cmd/ultraplan`.

### Decision 7: Deferred Scope Remains Explicitly Deferred

- **Decision:** Implement no runtime-assisted completion, no OpenCode/agentwrap execution, no prompt composition, no Markdown document discovery/frontmatter/applicability filtering, no report validation/synthesis/summary/run state, no target or sprint workflows, and no stable public JSON schema for study init.
- **Rationale:** The sprint contract is deliberately narrower than the full PRD/TRD. Keeping deferred scope out avoids premature architecture, protects CLI clarity, and lets future sprints make decisions with their own evidence.
- **Study / Source Grounding:** `technical-handbook.md` `15-philosophy` evidence supports rejecting scope and keeping complexity near benefit. No runtime evidence is relevant because the selected sprint index explicitly excluded LLM Runtime, Workflows, Observability, Performance, and Persistence contracts for this sprint.
- **Trade-Offs Accepted:** Some PRD/TRD long-term capabilities remain unavailable after Sprint 4. Users with incomplete YAML receive actionable failure guidance instead of assistance.
- **Technical Debt / Future Impact:** Future sprints must add assisted completion and runtime behavior intentionally, likely extending the init result/schema and generated README. Current code should leave seams but not implement dormant runtime hooks.
- **Alternatives Rejected:** Rejected dormant `--no-assist` or hidden runtime stubs because they create unsupported surface area. Rejected Markdown source discovery during init because Sprint 5 owns that behavior. Rejected target/sprint command references in README because they are outside current study-side release.
- **Contracts Satisfied:** Architecture contract; Documentation contract; SR-NG-01 to SR-NG-07; sprint-index excluded context.
- **Evidence Required:** Review check that no deferred flags/commands/runtime packages are added; README/help review; import review showing no agentwrap/OpenCode init usage; non-goal checklist in sprint review.

## Expected Evidence

| Evidence Type | Required Evidence | Source / Command / Review Check |
| --- | --- | --- |
| Tests | YAML parser accepts full explicit schema and rejects invalid YAML, missing fields, duplicate source names, duplicate dimension numbers/slugs, unsafe names, unsafe output paths, and unsupported count shortages. | `internal/study/init_test.go`; `go test ./...` from `/home/antonioborgerees/coding/ultraplan-go`. |
| Tests | Normalization renders two-digit dimension numbers, kebab-case slugs, deterministic filenames, normalized YAML, README, and dimension Markdown preserving purpose/steps/citations/questions. | `internal/study/init_test.go`; fixture/golden assertions. |
| Tests | Successful init creates `study-init.yml`, `README.md`, `dimensions/`, `sources/`, `reports/source/`, and `reports/final/` under workspace/default or supplied output root. | `internal/study/init_test.go`; temp workspace fixtures. |
| Tests | Dry-run prints planned directories/files/clone actions and creates, overwrites, deletes, and clones nothing. | `internal/study/init_test.go`; `internal/app/study_init_commands_test.go`. |
| Tests | Existing study directory fails without `--force`; `--force` overwrites known generated files only and does not touch outside paths or unknown files. | `internal/study/init_test.go`; sentinel-file checks. |
| Tests | `--no-clone` skips clone execution; URL-backed sources produce fake clone calls equivalent to `git clone --depth 1 <url> <sources/name>`; clone failures return partial/failure while artifacts remain. | `internal/study/init_test.go`; fake clone runner. |
| Tests | CLI help includes `study init`, flags, usage errors, dry-run output, overwrite protection, workspace/output behavior, and exit-code mapping. | `internal/app/study_init_commands_test.go`. |
| Build | Full test and build gates pass. | `go test ./...` and `go build ./cmd/ultraplan` from `/home/antonioborgerees/coding/ultraplan-go`. |
| Review | Architecture boundaries are preserved. | Import review: `internal/app -> internal/study`; `internal/study -> internal/workspace` allowed; platform packages do not import `internal/study`; CLI handlers do not own business rules. |
| Review | Security constraints are preserved. | Review path normalization/safety, no shell interpolation, scoped force behavior, secret-free generated artifacts, and no recursive source scanning. |
| Review | Deferred scope remains deferred. | Confirm no runtime-assisted completion, OpenCode/agentwrap execution, Markdown frontmatter/applicability, target/sprint workflows, prompt composition, run state, diagnostics persistence, or new stable JSON schema. |
| Documentation | Generated README and command help are accurate for current capabilities only. | Review generated README fixture and help output. |

## Assumptions And Risks

| Item | Type | Impact | Mitigation / Follow-Up |
| --- | --- | --- | --- |
| Existing workspace APIs can resolve roots and enforce path safety for default and `--output` initialization. | Assumption | If insufficient, implementation could duplicate path logic or leave escapes possible. | Use exported workspace helpers where available; add narrow helpers in `workspace` only if necessary and still keep study semantics in `study`. |
| Existing CLI has a partial-completion exit-code convention or can map to a non-zero failure with clear text. | Assumption | Clone failures need deterministic command behavior. | Inspect existing app exit mapping during implementation; map to partial completion if available, otherwise general failure with explicit partial wording. |
| `--output <dir>` may be explicitly selected but must still be normalized and safety-checked. | Risk | Ambiguous handling could allow writes outside intended roots or produce absolute paths in normalized YAML. | Define output root resolution in the plan: accept only normalized safe output roots according to existing workspace rules or explicit allowed output behavior; test inside/outside cases. |
| Users may expect `--force` to replace stale clone directories. | Risk | Known-file-only force may leave existing source dirs unchanged or cause clone destination conflicts. | Fail or report clone destination conflicts actionably; do not delete unknown source directories unless a future explicit clean behavior is specified. |
| Clone URLs can contain credentials or sensitive tokens. | Risk | Normalized YAML or command output could expose secrets if users provide them. | Do not add secrets; treat supplied URL as user metadata but avoid printing environment or expanded command strings; consider redaction if existing helpers support it. |
| YAML formatting choices may become de facto compatibility. | Risk | Future schema changes can be harder if tests/users rely on exact layout. | Keep deterministic but do not document a stable public schema beyond current requirements; future schema/version sprint can formalize compatibility. |
| Generated README may drift from future command surface. | Risk | Users may see outdated next-command guidance. | Make README rendering tests intentional; update in future command sprints. |
| Sequential clones may be slow for many repos. | Risk | Users with large repo lists wait longer. | Accept for Sprint 4; revisit only if requirements introduce clone concurrency or performance targets. |

## Implementation Constraints

- `internal/study` owns YAML parsing, validation, normalization, planning, rendering, filesystem writes, force/dry-run policy, and clone policy.
- `internal/app` owns command registration, flag parsing, existing workspace/global flag wiring, user-facing output, and exit-code mapping only.
- Do not create `internal/validation`, `internal/reports`, `internal/scheduler`, `internal/prompts`, `internal/workflow`, or similar global technical-layer packages for this sprint.
- `internal/study` may depend on `internal/workspace` and platform packages only through allowed exported APIs; platform packages must not import `internal/study`.
- Dry-run must use the same parse/validate/plan path as real initialization and must disable mutation only at the final execution boundary.
- All generated paths must be normalized and verified before writes or clone execution.
- `--force` must be explicit, scoped to the selected study directory, and must not delete unknown files or modify paths outside that directory.
- Clone execution must use executable plus argument array equivalent to `git clone --depth 1 <url> <dest>` and must be behind a fakeable boundary.
- Normal tests must not require network, GitHub, OpenCode, provider credentials, or real cloned repositories.
- Generated artifacts must be deterministic, human-editable, reviewable in Git, and secret-free except for user-supplied metadata that cannot be safely inferred or redacted without changing meaning.
- Count shortages caused by `repos.count` or `dimensions.count` greater than explicit items must fail with guidance that assisted completion is deferred.
- Do not recursively scan source repositories or inspect cloned contents during initialization.
- Do not add runtime, prompt composition, run-loop, diagnostics persistence, target, sprint, or stable public JSON behavior in this sprint.

## Plan Handoff

`plan.md` must execute these decisions. It must not invent architecture, scope, or decisions beyond this document.

The plan must carry forward:

- final decisions
- applicable contracts and requirement IDs
- expected evidence
- risks and assumptions
- required review protocols

Implementation order should start with `internal/study` parser/normalization/plan/render tests, then filesystem and clone execution tests, then `internal/app` command wiring/tests, then full `go test ./...` and `go build ./cmd/ultraplan` verification.

## Phase Exit Criteria

- [x] Selected context was read and used.
- [x] Area-specific reasoning documents were completed where applicable or explicitly marked not applicable.
- [x] Area-specific reasoning conclusions are reflected or explicitly overridden in final decisions.
- [x] Contracts are explicitly mapped to decisions and expected evidence.
- [x] Final decisions are clear enough for `plan.md` to execute.
- [x] Expected evidence is specific and reviewable.
