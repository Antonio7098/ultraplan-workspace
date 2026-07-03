# Sprint Technical Handbook: Study Initialization From YAML

> Project: `ultraplan-go`
> Sprint: `04-study-init-yaml`
> Source: `sprint-index.md`
> **Inputs Used:** `projects/ultraplan-go/project-index.md`, `projects/ultraplan-go/sprints/04-study-init-yaml/sprint-index.md`, `projects/ultraplan-go/sprints/04-study-init-yaml/requirements.md`, `projects/ultraplan-go/docs/PRD.md`, `projects/ultraplan-go/docs/TRD.md`, `projects/ultraplan-go/docs/ARCHITECTURE.md`, `templates/technical-handbook.md`, `studies/go-cli-study/reports/final/01-project-structure.md`, `studies/go-cli-study/reports/final/02-command-architecture.md`, `studies/go-cli-study/reports/final/03-dependency-injection.md`, `studies/go-cli-study/reports/final/05-error-handling.md`, `studies/go-cli-study/reports/final/06-io-abstraction.md`, `studies/go-cli-study/reports/final/09-terminal-ux.md`, `studies/go-cli-study/reports/final/11-testing-strategy.md`, `studies/go-cli-study/reports/final/13-security.md`, `studies/go-cli-study/reports/final/15-philosophy.md`

This handbook distills the studies and reports selected by `sprint-index.md` for sprint reasoning. It does not decide architecture or implementation.

PRD, TRD, and Architecture documents were used only for sprint scope context. Evidence claims below are grounded in the selected evidence reports and their cited source references.

## Selected Studies And Reports

| Study / Report | Path | Relevant Finding | Confidence |
| --- | --- | --- | --- |
| `01-project-structure` | `studies/go-cli-study/reports/final/01-project-structure.md` | Mature Go CLIs keep entrypoints thin and put application behavior behind `cmd/` or `internal/` boundaries; examples include `cmd/restic/main.go:37-114`, `cmd/root.go:112-127`, and `cmd/gh/main.go:6`. | high |
| `02-command-architecture` | `studies/go-cli-study/reports/final/02-command-architecture.md` | Factory-created commands with thin `RunE` handlers are common in maintainable CLIs; examples include `pkg/cmd/install.go:132-145`, `pkg/cmd/issue/list/list.go:47-118`, and `cmd/restic/main.go:37-114`. | high |
| `03-dependency-injection` | `studies/go-cli-study/reports/final/03-dependency-injection.md` | Manual DI through factories, options, and explicit service boundaries is preferred over global state; examples include `pkg/cmdutil/factory.go:16-43`, `executor.go:22-24`, and `internal/app/app.go:42-81`. | high |
| `05-error-handling` | `studies/go-cli-study/reports/final/05-error-handling.md` | Actionable CLI errors depend on wrapping, sentinels, typed errors, exit-code mapping, and user-facing hints; examples include `cmd/age/tui.go:37-54`, `go-task/errors/errors_task.go:13-32`, and `internal/ghcmd/cmd.go:44-49`. | high |
| `06-io-abstraction` | `studies/go-cli-study/reports/final/06-io-abstraction.md` | Injectable writers, filesystem seams, and fake command boundaries improve command and file-operation tests; examples include `pkg/iostreams/iostreams.go:551-568`, `internal/fs/interface.go:10-31`, and `executor.go:553-564`. | high |
| `09-terminal-ux` | `studies/go-cli-study/reports/final/09-terminal-ux.md` | CLI-first tools should prefer calm text output, non-TTY-safe behavior, and progress only when operations are long; relevant examples include `internal/cmd/prompt.go:20-256`, `lib/terminal/terminal.go:82-86`, and `internal/ui/terminal.go:10-34`. | medium |
| `11-testing-strategy` | `studies/go-cli-study/reports/final/11-testing-strategy.md` | Strong CLI projects combine table-driven tests, integration-style command tests, mock/fake infrastructure, and golden or fixture outputs; examples include `internal/cmd/main_test.go:64-174`, `pkg/httpmock/stub.go:35-199`, and `internal/backend/mock/backend.go:14-26`. | high |
| `13-security` | `studies/go-cli-study/reports/final/13-security.md` | External command execution and path-sensitive operations need explicit trust boundaries: argument arrays, safe executable lookup, validation, and secret redaction; examples include `internal/execext/exec.go:59-66`, `pkg/registry/transport.go:37-41`, and `pkg/surveyext/editor_manual.go:23`. | high |
| `15-philosophy` | `studies/go-cli-study/reports/final/15-philosophy.md` | Durable CLIs intentionally reject scope, keep complexity near its benefit, and use interface/decorator patterns where dry-run, test, or backend flexibility requires them; examples include `internal/chezmoi/system.go:25-45`, `fs/types.go:16-59`, and `internal/backend/backend.go:19-90`. | high |

## Relevant Patterns

- **Thin CLI Wiring, Product Logic Elsewhere:** The evidence repeatedly separates command construction from business logic. `01-project-structure` cites `cmd/restic/main.go:37-114`, `cmd/root.go:112-127`, and `cmd/gh/main.go:6` as thin entrypoints, while `02-command-architecture` cites `pkg/cmd/install.go:132-145` delegating to action logic. For this sprint, this pattern pressures `internal/app` to parse flags and print results while YAML validation, planning, filesystem writes, rendering, and clone behavior stay in `internal/study`.
- **Factory Functions And Explicit Dependencies:** Command factories and DI factories keep command setup testable. `02-command-architecture` cites `pkg/cmd/issue/list/list.go:47-118`, `pkg/cmd/install.go:132`, and `cmd/restic/cmd_backup.go:35`; `03-dependency-injection` cites `pkg/cmdutil/factory.go:16-43` and `pkg/cmd/factory/default.go:26-46`. Sprint reasoning should decide how much construction belongs in the existing app composition root without creating a god object.
- **Narrow Interfaces At Volatile Boundaries:** Evidence favors interfaces for filesystem, external process, terminal, and backend boundaries rather than every helper. `06-io-abstraction` cites `internal/fs/interface.go:10-31`, `internal/backend/backend.go:19-90`, and `pkg/iostreams/iostreams.go:551-568`; `15-philosophy` cites `internal/chezmoi/system.go:25-45`. For this sprint, the volatile boundaries are writes and `git clone`, not pure normalization helpers.
- **Plan-Then-Execute For Dry Run And Force Safety:** Chezmoi-style System wrappers and decorator patterns support dry-run/debug behavior (`internal/chezmoi/system.go:25-45` in `15-philosophy`). Study init has similar pressure: one validated plan should drive both `--dry-run` output and real execution so the dry-run path does not drift from the mutation path.
- **Actionable, Classified Errors:** `05-error-handling` cites `%w` wrapping examples (`vfs/zip.go:21`, `pkg/ssh/ssh_keys.go:64`) and user hints (`cmd/age/tui.go:37-54`, `internal/ghcmd/cmd.go:281-301`). For this sprint, errors should name the YAML field, source item, dimension item, path, clone target, or overwrite condition and preserve cause chains for command-level exit mapping.
- **Partial Failure Visibility:** `05-error-handling` cites multi-error aggregation in `pkg/action/uninstall.go:232-254`, `internal/chezmoierrors/chezmoi errors.go:24`, and `internal/checker/checker.go:70-91`. Clone failures can happen after artifact creation, so sprint reasoning should account for returning a partial/failure status while still telling the user which generated files exist.
- **Non-Shell External Command Boundary:** `13-security` identifies argument arrays as the baseline for subprocess safety, citing `cmd_obj_builder.go:38`, `pkg/surveyext/editor_manual.go:23`, and `internal/execext/exec.go:59-66`. For `git clone --depth 1 <url> <dest>`, the sprint should keep command construction explicit and testable, avoiding `sh -c` or string concatenation.
- **Fixture And Fake-First Tests:** `11-testing-strategy` cites CLI integration via `internal/cmd/main_test.go:64-174`, command fakes like `lazygit/pkg/commands/oscommands/fake_cmd_obj_runner.go:17-26`, and HTTP/IO mocks such as `pkg/httpmock/stub.go:35-199`. Sprint tests should cover YAML parsing, planned artifacts, command output, force/dry-run, and clone calls without network or real GitHub access.
- **Human-Editable Generated Artifacts:** The philosophy report highlights simplicity and readable artifacts as long-lived maintenance choices, with filename-encoded metadata in `internal/chezmoi/chezmoi.go:36-58` and interface-backed dry-run wrappers in `internal/chezmoi/system.go:25-45`. For this sprint, normalized YAML and dimension Markdown should be deterministic and reviewable rather than opaque generated state.

## Trade-Offs

| Trade-Off | Benefit | Cost | When It Matters |
| --- | --- | --- | --- |
| Explicit plan object before execution vs direct writes | Enables one validation path for dry-run and real init; makes force and path-safety review easier; mirrors decorator/dry-run evidence from `internal/chezmoi/system.go:25-45` | Adds an intermediate data structure and possible duplication between plan and renderer | Matters because `--dry-run`, `--force`, clone planning, and partial failures all need consistent output |
| Minimal clone runner interface vs broad process abstraction | A focused seam is easy to fake and keeps `git clone --depth 1` explicit; supported by command-runner fake evidence `fake_cmd_obj_runner.go:17-26` and argument-array evidence `cmd_obj_builder.go:38` | May need refactoring if later sprints add more external commands | Matters because this sprint only needs local `git`, not a general runtime/process framework |
| In-package `internal/study` files vs new subpackages | Keeps product behavior near study state and aligns with `internal/` application patterns (`cmd/restic/main.go:37-114`, `internal/restic/repository.go:18`) | A single package can grow large if future runtime, reports, and init logic are all mixed | Matters because requirements explicitly constrain study init behavior to `internal/study` and discourage global technical-layer packages |
| Golden/fixture output assertions vs flexible substring assertions | Golden-style comparisons catch accidental changes in README, YAML, and dimension Markdown; evidence includes `helm/internal/test/test.go:43` and `go-task/task_test.go:166-169` | Golden maintenance can slow iteration and produce noisy diffs | Matters for deterministic generated artifacts and command help/dry-run output |
| Typed/sentinel validation errors vs plain formatted errors | Supports exit-code mapping and user-specific messages; evidence includes `go-task/errors/errors.go:47-50`, `pkg/storage/driver/driver.go:27-36`, and `internal/ghcmd/cmd.go:44-49` | More error types to maintain for a small sprint | Matters for duplicate names, unsafe paths, count mismatch, existing output, and clone partial failure |
| Progress output vs quiet deterministic output | Progress can reassure users during long clone operations, as terminal UX evidence supports progress for operations over two seconds (`tui/progress.go:12-68`, `internal/ui/progress/`) | Progress complicates tests and non-TTY behavior | Matters only if cloning multiple repos can take long; dry-run and normal artifact writes should remain calm and script-friendly |

## Anti-Patterns And Warnings

- **Business Rules In `RunE`:** `02-command-architecture` flags long command bodies such as `cmd/root.go:49-183`, `cmd/age/age.go:105-321`, and `cmd/evaluate_sequence_command.go:152`. For this sprint, CLI handlers should not own YAML validation, path planning, artifact rendering, or clone execution policy.
- **Global Mutable Configuration Or IO:** `03-dependency-injection` warns about package-level config and singleton caches (`fs/config.go:793`, `fs/cache/cache.go:16-21`), while `06-io-abstraction` warns about direct `os.Stdout`/`os.Stderr` leaks (`cmd/ls/ls.go:42`, `cmd/age/tui.go:31`). Study init tests will be harder if output streams, workspace roots, or clone runners are hidden globals.
- **Shell String Construction For Clone:** `13-security` treats argument arrays as the baseline and warns against shell command strings. `git clone` should be invoked as executable plus args, not through `sh -c`, to avoid metacharacter interpretation in URLs or paths.
- **Panic For User YAML Or Parse Errors:** `05-error-handling` flags parse-time panics such as `dive/image/docker/manifest.go:18`. Invalid YAML, missing fields, duplicate dimensions, and unsafe names are user errors and should be returned as actionable validation failures.
- **Dry Run That Bypasses Validation:** If dry-run only prints intended paths without using the real planning/validation flow, it can miss path escape, duplicate, force, or count errors. This conflicts with the evidence-backed dry-run/decorator pressure from `internal/chezmoi/system.go:25-45`.
- **Unreviewable Generated Output:** Golden-output evidence exists because user-facing text and generated artifacts regress easily. Non-deterministic YAML ordering, absolute paths in normalized YAML, or generated Markdown that loses supplied steps/citations/questions would undermine the sprint goal.
- **Tests That Require Network Or Real Repositories:** `11-testing-strategy` warns against real external dependencies in unit tests. Clone behavior should use fakes, not GitHub, network, or actual cloned repositories.
- **Implementation-Detail Test Assertions:** The testing report warns about brittle UI-internal assertions like `internal/view/pod_test.go:23`. For this sprint, tests should assert generated files, content, clone args, exit behavior, and user output rather than private map sizes or helper call counts.

## Examples Worth Inspecting

| Example | Path / Source | Why It Is Useful |
| --- | --- | --- |
| Thin command-to-action delegation | `helm/pkg/cmd/install.go:132-145`, `helm/pkg/action/install.go:73-140` | Shows how command code can collect flags and delegate actual operation to a non-CLI action layer. |
| Factory with testable dependencies | `gh-cli/pkg/cmdutil/factory.go:16-43`, `gh-cli/pkg/cmd/factory/default.go:26-46` | Useful when deciding how the app composition root supplies IO, workspace, and study services to `study init`. |
| IO test constructor | `gh-cli/pkg/iostreams/iostreams.go:551-568` | Shows a simple shape for command tests that capture stdout/stderr without real terminals. |
| Filesystem interface and mock terminal | `restic/internal/fs/interface.go:10-31`, `restic/internal/ui/mock.go:10-53` | Useful for reasoning about how much filesystem abstraction is warranted for artifact generation tests. |
| Dry-run/debug filesystem wrapper | `chezmoi/internal/chezmoi/system.go:25-45` | Shows a pattern where write-like behavior can be wrapped for dry-run/read-only/debug modes. |
| Actionable hint errors | `age/cmd/age/tui.go:37-54`, `gh-cli/internal/ghcmd/cmd.go:281-301` | Useful for validation errors that should tell users exactly which YAML field or path to fix. |
| Multi-error or partial failure aggregation | `helm/pkg/action/uninstall.go:232-254`, `restic/internal/checker/checker.go:70-91` | Relevant to clone failures after successful artifact writes. |
| Command runner fake | `lazygit/pkg/commands/oscommands/fake_cmd_obj_runner.go:17-26` | Useful for testing clone planning and execution without network or Git. |
| Argument-array command execution | `go-task/internal/execext/exec.go:59-66`, `lazygit/cmd_obj_builder.go:38` | Useful for keeping `git clone` non-shell and injection-resistant. |
| Golden/fixture test infrastructure | `helm/internal/test/test.go:43`, `go-task/task_test.go:166-169` | Useful for deterministic README, normalized YAML, dimension Markdown, help, and dry-run output assertions. |

## Design Pressures

- Keep study initialization owned by `internal/study` while `internal/app` remains a command and output wiring layer.
- Support a single deterministic planning path that feeds validation, dry-run, force checks, artifact writes, and clone planning.
- Preserve human editability: generated dimension Markdown, `README.md`, and normalized `study-init.yml` should remain reviewable and stable.
- Treat YAML files, source names, dimension names, URLs, and output paths as untrusted input that must be normalized and validated before write or clone execution.
- Keep `git clone` behind a small fakeable boundary with explicit arg arrays and no shell interpolation.
- Make clone partial failure visible without hiding successful artifact creation.
- Avoid runtime-assisted completion and source URL repair paths in this sprint, even though broader docs mention them as future scope.
- Prefer behavior-focused unit and command tests that do not need network, GitHub, OpenCode, provider credentials, or real cloned repositories.
- Decide how much output polish is useful without adding interactive UX, spinners, or TUI dependencies to a batch-style init command.
- Keep generated paths relative and workspace-scoped unless an explicitly selected output root requires different diagnostics.

## Open Questions For Reasoning

- What is the minimal `internal/study` service API that keeps CLI handlers thin while avoiding a large generic service container?
- Should planning produce an explicit list of directories, files, and clone actions as a public internal type, or should this stay private to `InitStudy`?
- What exact fields should a study initialization validation error expose so command output can name YAML paths such as `repos.items[1].name` and `dimensions.items[0].number`?
- How should partial clone failure be represented under existing CLI exit-code conventions: a typed partial error, a result status plus error, or a joined error?
- What generated files may `--force` overwrite, and should it ever remove unknown files inside the study directory or only overwrite known generated paths?
- How should `--output <dir>` interact with workspace path safety and normalized YAML paths?
- Should clone actions be executed after all file writes complete, or interleaved with source directory creation, given the requirement to report successful artifact creation when clones fail?
- Is a filesystem abstraction needed for all artifact writes, or are temp directories plus a small writer/planner seam sufficient for this sprint?
- Which generated outputs deserve golden assertions versus direct structured assertions to balance determinism with maintenance cost?
- What user-facing wording should explain that `repos.count` or `dimensions.count` greater than explicit items is unsupported until assisted completion enters scope?

## Evidence Pointers

- `studies/go-cli-study/reports/final/01-project-structure.md`: Inspect thin CLI, `internal/` boundary, and command-per-file evidence around `cmd/restic/main.go:37-114`, `cmd/root.go:112-127`, `internal/restic/repository.go:18`, and `internal/cmd/config.go:1`.
- `studies/go-cli-study/reports/final/02-command-architecture.md`: Inspect command factories, options structs, lifecycle wrappers, and long `RunE` warnings around `pkg/cmd/install.go:132-145`, `pkg/cmd/issue/list/list.go:47-118`, `cmd/restic/cmd_backup.go:84-115`, and `cmd/root.go:49-183`.
- `studies/go-cli-study/reports/final/03-dependency-injection.md`: Inspect factory/manual DI and global-state cautions around `pkg/cmdutil/factory.go:16-43`, `executor.go:22-24`, `internal/app/app.go:42-81`, `fs/cache/cache.go:16-21`, and `fs/config.go:793`.
- `studies/go-cli-study/reports/final/05-error-handling.md`: Inspect actionable hints, typed errors, sentinel errors, exit-code mapping, and multi-error aggregation around `cmd/age/tui.go:37-54`, `go-task/errors/errors_task.go:13-32`, `internal/ghcmd/cmd.go:44-49`, and `pkg/action/uninstall.go:232-254`.
- `studies/go-cli-study/reports/final/06-io-abstraction.md`: Inspect IO stream tests, filesystem interfaces, and hardcoded IO warnings around `pkg/iostreams/iostreams.go:551-568`, `internal/fs/interface.go:10-31`, `executor.go:553-564`, `cmd/ls/ls.go:42`, and `pkg/yqlib/logger.go:18`.
- `studies/go-cli-study/reports/final/09-terminal-ux.md`: Inspect non-TTY fallback and progress-output guidance around `internal/cmd/prompt.go:20-256`, `lib/terminal/terminal.go:82-86`, `internal/ui/terminal.go:10-34`, and `tui/progress.go:12-68`.
- `studies/go-cli-study/reports/final/11-testing-strategy.md`: Inspect command integration, mocks, fakes, and golden tests around `internal/cmd/main_test.go:64-174`, `pkg/httpmock/stub.go:35-199`, `fake_cmd_obj_runner.go:17-26`, `helm/internal/test/test.go:43`, and `go-task/task_test.go:166-169`.
- `studies/go-cli-study/reports/final/13-security.md`: Inspect command execution and trust-boundary patterns around `internal/execext/exec.go:59-66`, `internal/execext/exec.go:152-157`, `pkg/surveyext/editor_manual.go:23`, `pkg/registry/transport.go:37-41`, and `internal/config/json/validator.go:146`.
- `studies/go-cli-study/reports/final/15-philosophy.md`: Inspect scope discipline, interface/decorator patterns, and complexity cautions around `internal/chezmoi/system.go:25-45`, `fs/types.go:16-59`, `internal/backend/backend.go:19-90`, `VISION.md:97`, and `pkg/analyze/parallel.go:13`.

## Handoff To Reasoning

- Use this handbook as evidence input.
- Validate whether the observed patterns fit this project's constraints.
- Do not copy external patterns without sprint-specific reasoning.
- Resolve the open questions in architecture or sprint reasoning before implementation decisions are finalized.
- Keep PRD/TRD/Architecture context separate from evidence: they define sprint scope and constraints, while selected reports provide the pattern evidence above.
