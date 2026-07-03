# Sprint Technical Handbook: Study Domain, Listing, and Resolution

> Project: `ultraplan-go`
> Sprint: `03-study-listing-resolution`
> Source: `sprint-index.md`
> **Inputs Used:** `projects/ultraplan-go/sprints/03-study-listing-resolution/sprint-index.md`, `templates/technical-handbook.md`, `projects/ultraplan-go/project-index.md`, `projects/ultraplan-go/sprints/03-study-listing-resolution/requirements.md`, `projects/ultraplan-go/docs/PRD.md`, `projects/ultraplan-go/docs/TRD.md`, `projects/ultraplan-go/docs/ARCHITECTURE.md`, `studies/go-cli-study/reports/final/01-project-structure.md`, `studies/go-cli-study/reports/final/02-command-architecture.md`, `studies/go-cli-study/reports/final/03-dependency-injection.md`, `studies/go-cli-study/reports/final/04-configuration-management.md`, `studies/go-cli-study/reports/final/05-error-handling.md`, `studies/go-cli-study/reports/final/06-io-abstraction.md`, `studies/go-cli-study/reports/final/09-terminal-ux.md`, `studies/go-cli-study/reports/final/11-testing-strategy.md`, `studies/go-cli-study/reports/final/13-security.md`, `studies/go-cli-study/reports/final/14-performance.md`, `studies/go-cli-study/reports/final/15-philosophy.md`

This handbook distills the studies and reports selected by `sprint-index.md` for sprint reasoning. It does not decide architecture or implementation.

## Selected Studies And Reports

| Study / Report | Path | Relevant Finding | Confidence |
| --- | --- | --- | --- |
| `01-project-structure` | `studies/go-cli-study/reports/final/01-project-structure.md` | Mature Go CLIs keep CLI entrypoints thin and move substantive behavior into `internal/`, `pkg/`, or domain packages; examples include chezmoi `main.go:26-34`, gh-cli `cmd/gh/main.go:6`, and yq `cmd/root.go:9`. | high |
| `02-command-architecture` | `studies/go-cli-study/reports/final/02-command-architecture.md` | Effective command trees use command factory functions and thin `RunE` delegates, such as gh-cli `pkg/cmd/issue/list/list.go:47-118`, helm `pkg/cmd/install.go:132-145`, and restic `cmd/restic/cmd_backup.go:84-115`. | high |
| `03-dependency-injection` | `studies/go-cli-study/reports/final/03-dependency-injection.md` | Explicit composition roots and constructor/factory injection improve traceability and tests; examples include gh-cli `internal/ghcmd/cmd.go:52-132`, `pkg/cmdutil/factory.go:16-43`, and go-task `executor.go:93-114`. | high |
| `04-configuration-management` | `studies/go-cli-study/reports/final/04-configuration-management.md` | Config precedence must distinguish explicit flags from defaults; restic uses `Flag.Changed` tracking at `internal/global/global.go:139,147`, and go-task centralizes precedence in `internal/flags/flags.go:314-327`. | medium |
| `05-error-handling` | `studies/go-cli-study/reports/final/05-error-handling.md` | User-facing errors benefit from wrapping, typed/sentinel errors, and actionable hints; examples include go-task `errors/errors_task.go:13-32`, age `cmd/age/tui.go:37-54`, and gh-cli `internal/ghcmd/cmd.go:281-301`. | high |
| `06-io-abstraction` | `studies/go-cli-study/reports/final/06-io-abstraction.md` | Commands are easier to test when stdout/stderr/filesystem dependencies are injectable; examples include gh-cli `pkg/iostreams/iostreams.go:551-568`, go-task `executor.go:553-564`, and restic `internal/ui/mock.go:10-53`. | high |
| `09-terminal-ux` | `studies/go-cli-study/reports/final/09-terminal-ux.md` | CLI-first listing output should remain calm, non-blocking, and script-friendly; non-TTY fallback and color/TTY checks appear in chezmoi `internal/cmd/prompt.go:20-256` and gh-cli `internal/prompter/`. | medium |
| `11-testing-strategy` | `studies/go-cli-study/reports/final/11-testing-strategy.md` | Strong CLI projects test both domain behavior and command output with table tests, golden/testscript, or helper fakes; examples include chezmoi `internal/cmd/main_test.go:64-174`, go-task `task_test.go:166-169`, and gh-cli `acceptance/acceptance_test.go:26-29`. | high |
| `13-security` | `studies/go-cli-study/reports/final/13-security.md` | Trust boundaries require explicit path/input validation and safe handling of command/file operations; examples include k9s schema validation `internal/config/json/validator.go:146`, gh-cli `safeexec` at `pkg/surveyext/editor_manual.go:23`, and chezmoi private temp handling `gpgencryption.go:151-165`. | medium |
| `14-performance` | `studies/go-cli-study/reports/final/14-performance.md` | Fast CLIs defer expensive work and avoid unbounded scans; examples include gh-cli lazy factory fields `pkg/cmdutil/factory.go:27-42`, yq streaming `pkg/yqlib/stream_evaluator.go:78-113`, and gdu bounded traversal `pkg/analyze/parallel.go:13`. | high |
| `15-philosophy` | `studies/go-cli-study/reports/final/15-philosophy.md` | Coherent tools reject unnecessary complexity and keep patterns tied to product purpose; examples include age's no-keyring constraint `age.go:18`, gh-cli factory pattern `pkg/cmd/factory/default.go:26-46`, and restic backend interface `internal/backend/backend.go:19-90`. | medium |

## Relevant Patterns

- **Thin CLI, Domain-Owned Behavior:** The selected structure and command reports repeatedly show CLI layers as construction, flag, and output boundaries while domain packages own business behavior. Chezmoi's `main.go:26-34`, gh-cli's `cmd/gh/main.go:6`, helm's `pkg/cmd/install.go:132-145`, and yq's `cmd/root.go:9` support keeping study discovery and reference resolution outside command handlers.
- **Command Factory With Injected Dependencies:** Command reports and DI reports point to `newXxxCmd`/factory construction as the dominant scalable pattern. Gh-cli uses `pkg/cmdutil/factory.go:16-43` and `pkg/cmd/issue/list/list.go:47-118`; restic uses `cmd/restic/main.go:37-114` and `cmd/restic/cmd_backup.go:84-115`; helm wires root commands through `pkg/cmd/root.go:105` and delegates to actions at `pkg/action/install.go:73-140`.
- **Explicit Service Surface Over Global State:** DI evidence favors a composition root plus explicit service/factory fields over package-level configuration or mutable singletons. Examples include gh-cli `internal/ghcmd/cmd.go:52-132`, go-task `executor.go:93-114`, and opencode `internal/app/app.go:42-81`; the caution side is rclone global config `fs/config.go:14-51` and yq parser/log singletons `pkg/yqlib/lib.go:13-21`.
- **Injectable Output For Command Tests:** IO and testing reports support passing writers or IO stream structs so listing output can be captured in tests. Gh-cli's `iostreams.Test()` at `pkg/iostreams/iostreams.go:551-568`, go-task `WithStdout` at `executor.go:553-564`, and mitchellh-cli `MockUi` at `ui_mock.go:27-33` are direct evidence.
- **Actionable Reference Errors:** Error evidence supports structured missing/ambiguous reference errors with context and hints. Go-task's `TaskNotFoundError` carries `TaskName` and `DidYouMean` at `errors/errors_task.go:13-32`; age's `errorWithHint()` at `cmd/age/tui.go:47-54` and gh-cli's `printError()` at `internal/ghcmd/cmd.go:281-301` show user-facing diagnostic patterns.
- **Deterministic, Bounded Discovery:** Performance and security reports both point toward bounded traversal and deliberate input scope. Gdu bounds parallel traversal with `pkg/analyze/parallel.go:13`, yq provides bounded streaming mode at `pkg/yqlib/stream_evaluator.go:78-113`, and security cautions call out missing path canonicalization in gdu `internal/common/ignore.go:16-37`; this sprint's shallow study/source/dimension listing should preserve bounded scope.
- **Behavior-Focused Unit And Command Tests:** Testing evidence favors table-driven tests, fixture workspaces, output assertions, and CLI-level integration where output is user-facing. Chezmoi's testscript execution at `internal/cmd/main_test.go:64-174`, go-task golden normalization at `task_test.go:166-169`, and gdu `runApp` helper at `cmd/gdu/app/app_test.go:682` are useful analogues.

## Trade-Offs

| Trade-Off | Benefit | Cost | When It Matters |
| --- | --- | --- | --- |
| Keep discovery/resolution in `internal/study` rather than command handlers | Matches thin CLI evidence from chezmoi `main.go:26-34`, helm `pkg/cmd/install.go:132-145`, and yq `cmd/root.go:9`; improves unit testing and reuse | Requires an explicit service/API boundary between `internal/app` and `internal/study` | Matters for `ultraplan study list` and `ultraplan study <study> list`, where commands should format output while study owns filesystem and reference rules |
| Use explicit command factories and injected writers/services | Enables command tests with buffer output, matching gh-cli `pkg/iostreams/iostreams.go:551-568` and go-task `executor.go:553-564` | More setup than direct `fmt.Println` or package globals | Matters because acceptance requires command output and error behavior tests without real OpenCode, network, or credentials |
| Support prefix resolution for usability | Improves UX like go-task's typo-aware `TaskNotFoundError` at `errors/errors_task.go:13-32`; reduces typing for study/source/dimension refs | Requires precise ambiguity detection and actionable errors to avoid surprising matches | Matters for dimension references by number, slug, full filename, exact reference, or unambiguous prefix |
| Keep discovery shallow instead of recursively scanning source trees | Fast startup and safer trust boundary; aligns with performance evidence around bounded traversal `pkg/analyze/parallel.go:13` and streaming/bounded modes `pkg/yqlib/stream_evaluator.go:78-113` | Listing cannot infer nested repository metadata or Markdown document sources deferred to later scope | Matters because source directories may be large repositories, and sprint scope explicitly forbids recursive source scans |
| Add only minimal output formatting now | Calm, script-friendly listing aligns with CLI-first UX evidence and avoids TUI/progress overhead | Less polished than rich tables, colors, or JSON output modes | Matters because sprint scope is read-only listing and non-trivial JSON listing is a non-goal |

## Anti-Patterns And Warnings

- **Do Not Put Domain Rules In `RunE`:** The command architecture report flags long command handlers such as opencode `cmd/root.go:49-183`, yq `cmd/evaluate_sequence_command.go:152`, and restic `runBackup()` as logic leakage. Study discovery, sorting, normalization, and prefix resolution should not become command-handler internals.
- **Avoid Global Mutable Config Or IO State:** DI/config/IO reports warn about yq globals `pkg/yqlib/lib.go:13-21`, rclone fallback globals `fs/config.go:793`, and hardcoded output paths like rclone `cmd/ls/ls.go:42`. Listing commands should consume the existing workspace context and injected writers rather than reading global state directly.
- **Do Not Recursively Inspect Source Repositories For Listing:** Performance and security evidence warns against unbounded in-memory accumulation and unbounded traversal; examples include table cloning in k9s `internal/model/table.go:200`, path canonicalization gaps in gdu `internal/common/ignore.go:16-37`, and the general bounded traversal pattern at `gdu/pkg/analyze/parallel.go:13`.
- **Do Not Collapse Missing And Ambiguous References Into Generic Errors:** Error reports show mature CLIs distinguish recoverable/domain conditions with typed or sentinel errors, such as helm `pkg/storage/driver/driver.go:27-36`, go-task `errors/errors.go:47-50`, and gh-cli `internal/ghcmd/cmd.go:44-49`. Generic `not found` messages would undercut the sprint acceptance criteria.
- **Avoid Over-Broad Abstractions For A Small Listing Sprint:** Philosophy and structure evidence warn about god structs and accumulated complexity, including chezmoi's large `internal/cmd/config.go:193-291` and lazygit's acknowledged God Struct migration. This sprint needs study discovery/listing seams, not a generalized workflow, report, scheduler, or validation framework.

## Examples Worth Inspecting

| Example | Path / Source | Why It Is Useful |
| --- | --- | --- |
| Thin entrypoint delegation | `chezmoi main.go:26-34`, `gh-cli cmd/gh/main.go:6`, `yq cmd/root.go:9` | Shows how little business logic needs to live at the binary/command boundary. |
| Command factory pattern | `gh-cli pkg/cmd/issue/list/list.go:47-118`, `restic cmd/restic/main.go:37-114`, `helm pkg/cmd/root.go:105` | Shows consistent command construction with injected dependencies. |
| Injectable command output | `gh-cli pkg/iostreams/iostreams.go:551-568`, `go-task executor.go:553-564`, `mitchellh-cli ui_mock.go:27-33` | Shows practical test seams for command output assertions. |
| Domain-level interface boundary | `restic internal/restic/repository.go:18-66`, `restic internal/backend/backend.go:19-90`, `chezmoi internal/chezmoi/system.go:25` | Shows interfaces at external/volatile boundaries rather than every helper. |
| Actionable reference-style errors | `go-task errors/errors_task.go:13-32`, `age cmd/age/tui.go:47-54`, `gh-cli internal/ghcmd/cmd.go:281-301` | Shows how user-facing errors can include identifiers and correction hints. |
| Shallow/bounded traversal discipline | `gdu pkg/analyze/parallel.go:13`, `yq pkg/yqlib/stream_evaluator.go:78-113`, `age internal/stream/stream.go:20` | Shows bounded work as a design pressure even when exact mechanisms differ. |
| CLI behavior tests | `chezmoi internal/cmd/main_test.go:64-174`, `go-task task_test.go:166-169`, `gdu cmd/gdu/app/app_test.go:682` | Shows command-level tests for output and workflows. |

## Design Pressures

- **Sprint scope is read-only and inspection-oriented:** Evidence supports thin command surfaces and bounded discovery; the sprint should avoid run-loop, prompt, runtime, and scheduler patterns even though later product requirements need them.
- **Discovery must be deterministic:** Reports emphasize stable output and behavior-focused tests; study/source/dimension ordering must be explicit enough to test and reason about.
- **Resolution must balance convenience and safety:** Prefix matching improves ergonomics, but error evidence pressures the design to distinguish missing from ambiguous references and include the conflicting candidates or useful next steps.
- **Workspace safety matters even for listing:** Security evidence treats filesystem paths as trust boundaries; listing should stay inside the resolved workspace and avoid recursively entering source repositories.
- **Testing needs both domain and CLI coverage:** IO/test evidence pressures the design toward injectable writers and fixture workspaces so command tests can cover output and errors without external dependencies.
- **Package boundaries should stay module-driven:** Structure/philosophy evidence supports `internal/study` owning study behavior and `internal/app` owning command wiring, not global technical packages.

## Open Questions For Reasoning

- What exact public surface should `internal/study` expose so `internal/app` can list studies and one study's details without learning discovery internals?
- Should reference resolution return typed domain errors, sentinel-wrapped errors, or a small structured error type that command rendering can inspect for missing vs ambiguous cases?
- What candidate strings participate in dimension prefix matching: number, slug, full filename, normalized filename, or all of them, and how should ambiguity be reported?
- How much path information should listing output include while keeping user-facing paths workspace-relative and avoiding absolute path leakage?
- Should absent `studies/`, `sources/`, or `dimensions/` directories be treated as empty listings or actionable errors in this sprint's listing commands?
- What minimum command output shape is stable enough for humans and tests without expanding into non-trivial JSON output?
- Which filesystem operations need abstraction now, and which can remain concrete because the sprint only needs shallow local reads and temp-dir tests?

## Evidence Pointers

- `studies/go-cli-study/reports/final/01-project-structure.md`: Inspect Pattern 1, Pattern 2, and Evidence Index for `cmd/` plus `internal/` boundary examples such as `main.go:26-34`, `cmd/gh/main.go:6`, and `internal/chezmoi/chezmoi.go:1-2`.
- `studies/go-cli-study/reports/final/02-command-architecture.md`: Inspect Pattern 1, Thin-Delegate model, and Anti-Patterns for command factory and command-handler size evidence including `pkg/cmd/issue/list/list.go:47-118`, `pkg/cmd/install.go:132-145`, and `cmd/root.go:49-183`.
- `studies/go-cli-study/reports/final/03-dependency-injection.md`: Inspect constructor/factory injection and global-state cautions, especially `pkg/cmdutil/factory.go:16-43`, `executor.go:93-114`, `fs/config.go:793`, and `pkg/yqlib/lib.go:13-21`.
- `studies/go-cli-study/reports/final/05-error-handling.md`: Inspect typed errors, sentinel errors, and hint rendering at `errors/errors_task.go:13-32`, `cmd/age/tui.go:37-54`, `internal/ghcmd/cmd.go:281-301`, and `pkg/storage/driver/driver.go:27-36`.
- `studies/go-cli-study/reports/final/06-io-abstraction.md`: Inspect IOStreams, writer injection, and hardcoded IO cautions at `pkg/iostreams/iostreams.go:551-568`, `executor.go:553-564`, `cmd/ls/ls.go:42`, and `internal/cmd/templatefuncs.go:296`.
- `studies/go-cli-study/reports/final/11-testing-strategy.md`: Inspect testscript, golden tests, and behavior-focused test guidance at `internal/cmd/main_test.go:64-174`, `task_test.go:166-169`, `internal/test/test.go:43`, and `cmd/gdu/app/app_test.go:682`.
- `studies/go-cli-study/reports/final/13-security.md`: Inspect path/input trust boundary guidance and cautions at `internal/config/json/validator.go:146`, `internal/common/ignore.go:16-37`, and `pkg/surveyext/editor_manual.go:23`.
- `studies/go-cli-study/reports/final/14-performance.md`: Inspect lazy initialization and bounded traversal guidance at `pkg/cmdutil/factory.go:27-42`, `pkg/analyze/parallel.go:13`, and `pkg/yqlib/stream_evaluator.go:78-113`.
- `studies/go-cli-study/reports/final/15-philosophy.md`: Inspect simplicity and deliberate complexity guidance at `age.go:18`, `pkg/cmd/factory/default.go:26-46`, `internal/backend/backend.go:19-90`, and the anti-patterns section.

## Handoff To Reasoning

- Use this handbook as evidence input.
- Validate whether the observed patterns fit this project's constraints.
- Do not copy external patterns without sprint-specific reasoning.
- Keep final implementation decisions in sprint reasoning, especially around `internal/study` service shape, reference error types, output format, and absent-directory behavior.
