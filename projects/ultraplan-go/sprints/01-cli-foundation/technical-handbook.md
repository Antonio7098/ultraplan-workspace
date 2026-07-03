# Sprint Technical Handbook: CLI Foundation

> Project: `ultraplan-go`
> Sprint: `01-cli-foundation`
> Source: `sprint-index.md`
> **Inputs Used:** `projects/ultraplan-go/sprints/01-cli-foundation/sprint-index.md`, `templates/technical-handbook.md`, `projects/ultraplan-go/project-index.md`, `projects/ultraplan-go/sprints/01-cli-foundation/requirements.md`, `projects/ultraplan-go/docs/PRD.md`, `projects/ultraplan-go/docs/TRD.md`, `projects/ultraplan-go/docs/ARCHITECTURE.md`, `studies/go-cli-study/reports/final/01-project-structure.md`, `studies/go-cli-study/reports/final/02-command-architecture.md`, `studies/go-cli-study/reports/final/03-dependency-injection.md`, `studies/go-cli-study/reports/final/05-error-handling.md`, `studies/go-cli-study/reports/final/06-io-abstraction.md`, `studies/go-cli-study/reports/final/09-terminal-ux.md`, `studies/go-cli-study/reports/final/11-testing-strategy.md`, `studies/go-cli-study/reports/final/12-extensibility.md`, `studies/go-cli-study/reports/final/15-philosophy.md`

This handbook distills the studies and reports selected by `sprint-index.md` for sprint reasoning. It does not decide architecture or implementation.

## Selected Studies And Reports

| Study / Report | Path | Relevant Finding | Confidence |
| --- | --- | --- | --- |
| `01-project-structure` | `studies/go-cli-study/reports/final/01-project-structure.md` | Application CLIs consistently keep entrypoints thin and delegate to protected interior packages; examples include `cmd/gh/main.go:6`, `cmd/restic/main.go:162`, and `cmd/root.go:112-127`. | high |
| `02-command-architecture` | `studies/go-cli-study/reports/final/02-command-architecture.md` | Maintainable commands separate command wiring from business logic; examples include helm's `pkg/cmd/install.go:132-145`, gh-cli's `pkg/cmdutil/factory.go:16-43`, and rclone's `cmd/cmd.go:240-340`. | high |
| `03-dependency-injection` | `studies/go-cli-study/reports/final/03-dependency-injection.md` | Manual composition roots dominate; examples include `internal/ghcmd/cmd.go:52-132`, `internal/app/app.go:42-81`, and `cmd/restic/main.go:181-183`. | high |
| `05-error-handling` | `studies/go-cli-study/reports/final/05-error-handling.md` | Mature CLIs use wrapped errors, user-facing hints, sentinel or typed errors, and exit-code mapping; examples include `cmd/age/tui.go:37-54`, `internal/ghcmd/cmd.go:44-49`, and `go-task/errors/errors_task.go:13-32`. | high |
| `06-io-abstraction` | `studies/go-cli-study/reports/final/06-io-abstraction.md` | Deterministic command tests are enabled by injectable IO streams and test constructors; examples include `pkg/iostreams/iostreams.go:551-568`, `executor.go:553-564`, and `ui.go:19-43`. | high |
| `09-terminal-ux` | `studies/go-cli-study/reports/final/09-terminal-ux.md` | CLI-first tools should keep output simple, script-friendly, and TTY-aware; heavier TUI/progress patterns are justified only by interaction density or long operations. Evidence includes `internal/cmd/prompt.go:124-137`, `lib/terminal/terminal.go:82-86`, and `help.go:371-388`. | medium |
| `11-testing-strategy` | `studies/go-cli-study/reports/final/11-testing-strategy.md` | Strong CLI projects test behavior through command-level or integration-style tests with captured output, golden files, or testscript; examples include `internal/cmd/main_test.go:64-174`, `go-task/task_test.go:166-169`, and `gdu/cmd/gdu/app/app_test.go:682`. | high |
| `12-extensibility` | `studies/go-cli-study/reports/final/12-extensibility.md` | Extensibility should match scope: interface boundaries and functional options can support future growth without introducing plugin/runtime complexity; examples include `go-task/executor.go:20-24,91-122` and `gh-cli/pkg/cmdutil/factory.go:16-43`. | medium |
| `15-philosophy` | `studies/go-cli-study/reports/final/15-philosophy.md` | Coherent Go CLIs make non-goals explicit and accept complexity deliberately; evidence includes age's no-keyring stance at `age.go:18`, fzf's focused filter purpose at `README.md:28-29`, and lazygit's stated complexity rejection at `VISION.md:97`. | high |

## Relevant Patterns

- **Thin process entrypoint with app-owned command handling:** The project-structure report identifies thin CLI entrypoints as the dominant application pattern, with `cmd/gh/main.go:6` and `main.go:23` delegating almost immediately and restic keeping command construction under `cmd/restic/main.go:37-114`. For this sprint, that evidence supports treating `cmd/ultraplan/main.go` as a process boundary only while leaving command behavior and tests in `internal/app`.
- **Internal package boundary for application-only behavior:** The project-structure report shows pure applications using `internal/` to enforce encapsulation, including chezmoi's `internal/chezmoi/chezmoi.go:1-2`, restic's `internal/restic/repository.go:18`, and opencode's `internal/app/` plus `internal/llm/` structure. This is relevant because the sprint creates an application CLI shell, not a public Go SDK.
- **Single composition root with explicit wiring:** The DI report finds that high-scoring projects use visible manual composition roots rather than DI frameworks, including gh-cli's `internal/ghcmd/cmd.go:52-132`, opencode's `internal/app/app.go:42-81`, and restic's `cmd/restic/main.go:181-183`. This supports reasoning about `internal/app` as the place where command parsing, IO streams, version metadata, and exit status mapping meet.
- **Injectable IO for command tests:** The IO report highlights gh-cli's `IOStreams.Test()` at `pkg/iostreams/iostreams.go:551-568`, go-task's `WithStdout()` at `executor.go:553-564`, and mitchellh-cli's `Ui` interface at `ui.go:19-43`. For this sprint's deterministic help/version/invalid-command tests, output should be capturable without real terminal state.
- **User-facing error rendering separated from error classification:** The error-handling report shows age's `errorf()` and `errorWithHint()` at `cmd/age/tui.go:37-54`, gh-cli's exit constants at `internal/ghcmd/cmd.go:44-49`, and go-task's `TaskNotFoundError` with `DidYouMean` at `errors/errors_task.go:13-32`. The sprint's unknown-command behavior needs actionable output and exit code `2`, but reasoning still needs to decide the minimal classification shape.
- **Behavior-first command tests:** The testing report favors output and exit-code assertions over implementation inspection, with testscript examples at `internal/cmd/main_test.go:64-174`, command helper integration at `gdu/cmd/gdu/app/app_test.go:682`, and golden comparison for user output at `go-task/task_test.go:166-169`. This fits the sprint requirement to assert implemented command behavior rather than manual output inspection.
- **Explicit non-goals as design control:** The philosophy report shows focused tools rejecting features deliberately, such as age avoiding a global keyring at `age.go:18`, fzf staying a fuzzy finder at `README.md:28-29`, and lazygit documenting features not worth complexity at `VISION.md:97`. This matters because the sprint intentionally excludes workspace, config, study, runtime, target, and sprint workflows from the initial help surface.

## Trade-Offs

| Trade-Off | Benefit | Cost | When It Matters |
| --- | --- | --- | --- |
| `cmd/ultraplan` thin entrypoint vs richer root command file | Keeps process concerns separate from app behavior; mirrors thin entry evidence such as `cmd/gh/main.go:6` and `cmd/restic/main.go:37-114`; makes `internal/app` testable without `os.Exit`. | Adds one delegation hop and requires a clear app API. | Sprint reasoning must choose how `main.go` calls `internal/app` and receives a status code. |
| `internal/` application packages vs public `pkg/` surface | Go enforces private application boundaries; aligns with application projects using `internal/chezmoi/`, `internal/restic/`, and opencode internals. | Future public API extraction would require refactoring to `pkg/`. | This sprint creates application-owned skeletons and no public SDK. |
| Minimal custom command shell vs adopting a full CLI framework | A small shell can satisfy help/version/unknown-command behavior with low dependency and no framework conventions. | A framework may provide help formatting, subcommands, flag handling, and future growth patterns already solved. | Sprint reasoning must weigh current tiny command surface against future study command expansion. |
| Explicit IO streams vs direct `fmt.Println`/`os.Stdout` | Capturable output enables deterministic tests, following gh-cli `IOStreams.Test()` and go-task stream options. | Slightly more plumbing in `internal/app`. | Help, version, and invalid-command tests must not require a real terminal. |
| Simple error/status values vs typed error taxonomy | Minimal shape is enough for usage errors now; avoids premature categories. | Later config/workspace/runtime errors may need richer mapping. | Unknown command currently needs only actionable diagnostic plus exit code `2`; future commands may expand the taxonomy. |
| Static skeleton packages now vs delayed package creation | Makes intended module boundaries visible early and matches sprint deliverables. | Empty packages can imply unearned architecture if over-documented or too numerous. | Skeleton package docs should describe boundaries without implementing deferred behavior. |

## Anti-Patterns And Warnings

- **Do not let `cmd/ultraplan/main.go` accumulate command logic:** The project-structure and command reports flag large entrypoints and `RunE` bodies as caution signs, including age's `cmd/age/age.go:105-321` and opencode's `cmd/root.go:49-183`. The sprint acceptance criteria explicitly require a thin entrypoint.
- **Do not hardcode stdout/stderr in the app logic that tests need to exercise:** The IO report identifies direct `os.Stdout`/`os.Stderr` usage as a testability leak, including `cmd/age/tui.go:31`, `cmd/ls/ls.go:42`, and `pkg/yqlib/logger.go:18`. Direct process IO belongs at the outer boundary, not in command behavior under test.
- **Do not add advertised commands for deferred workflows:** The selected sprint scope excludes workspace, config, health, study, runtime, code extraction, target, and sprint commands. The philosophy report's focused-tool examples show that explicit rejection of features is a strength when scope is narrow, as in `age.go:18` and `README.md:28-29`.
- **Do not introduce global mutable config or service singletons for a shell sprint:** The DI report warns against config and service singletons such as yq's `ConfiguredYamlPreferences` mutation at `cmd/root.go:121` and rclone's cache singleton at `fs/cache/cache.go:16-21`. This sprint does not need shared mutable state.
- **Do not create premature plugin, runtime, or registry mechanics:** The extensibility report shows plugin systems bring versioning, validation, lifecycle, and security concerns, such as helm metadata at `internal/plugin/metadata_v1.go:24-48` and subprocess runtime at `internal/plugin/runtime_subprocess.go:65-79`. Runtime and plugin behavior are non-goals for this sprint.
- **Do not test implementation details instead of behavior:** The testing report flags brittle implementation-detail assertions such as k9s hint counts at `internal/view/pod_test.go:23`. This sprint's tests should assert output and exit codes for help, version, and unknown commands.

## Examples Worth Inspecting

| Example | Path / Source | Why It Is Useful |
| --- | --- | --- |
| Minimal thin entrypoint | `gh-cli/cmd/gh/main.go:6`; report `01-project-structure` | Shows an entrypoint that defers real work elsewhere. |
| Root command factory | `restic/cmd/restic/main.go:37-114`; report `01-project-structure` and `02-command-architecture` | Useful for reasoning about a single app-owned command construction point without putting logic in `main`. |
| Factory dependency object | `gh-cli/pkg/cmdutil/factory.go:16-43`; reports `02-command-architecture`, `03-dependency-injection`, `12-extensibility` | Shows explicit dependency aggregation and test seams. |
| App composition root | `opencode/internal/app/app.go:42-81`; report `03-dependency-injection` | Relevant to `internal/app` as a composition boundary, though opencode's CLI orchestration is also cited as a caution. |
| Injectable output streams | `gh-cli/pkg/iostreams/iostreams.go:551-568`; report `06-io-abstraction` | Shows production/test IO separation with buffers. |
| Exit code mapping | `gh-cli/internal/ghcmd/cmd.go:44-49`; report `05-error-handling` | Useful for script-friendly CLI status mapping. |
| Actionable user hints | `age/cmd/age/tui.go:37-54`; report `05-error-handling` | Shows user-facing error text separated from internal errors. |
| Behavior-oriented command tests | `gdu/cmd/gdu/app/app_test.go:682`; report `11-testing-strategy` | Shows command-level testing through an app helper rather than manual process checks. |
| Golden or fixture output checks | `go-task/task_test.go:166-169`; report `11-testing-strategy` | Relevant if help/version text should be protected from accidental drift. |
| Focused non-goal philosophy | `VISION.md:97`, `age.go:18`, `README.md:28-29`; report `15-philosophy` | Supports strict scope control in the CLI foundation sprint. |

## Design Pressures

- The sprint must establish enough command architecture to be testable, but not enough to pre-decide the later study command tree.
- The entrypoint must be operationally real while remaining small enough that `internal/app` owns command behavior and test seams.
- The package layout must express future module ownership without implementing workspace, config, runtime, study, or code extraction behavior.
- Help output must aid command discovery but must not advertise deferred commands, even if the larger PRD lists them as future product capabilities.
- Version output needs build metadata fields, but the handbook evidence only supports treating build-time variables as simple immutable metadata, not a global state system.
- Tests must run offline and deterministically, so runtime, filesystem workspace discovery, network, provider credentials, and terminal state are design constraints to avoid.
- Error handling must be script-friendly from the start, but the sprint only needs the first usage-error classification shape.
- The future product will need runtime/config/workspace integration, yet this sprint should leave those boundaries documented rather than wired.

## Open Questions For Reasoning

- Should the initial command shell use the Go standard library only, a small custom dispatcher, or a CLI framework, given the current command set is tiny but future study commands will be nested?
- What is the minimal `internal/app` API that keeps `cmd/ultraplan/main.go` thin while making output, args, version metadata, and exit status testable?
- Should unknown-command diagnostics go to stdout or stderr, and how should that interact with deterministic tests and script-friendly behavior?
- Should no-argument help and `--help` share one exact rendering path, or may they differ in preamble while preserving command discovery?
- What version metadata defaults are acceptable before release-time linker injection exists?
- How much documentation should skeleton packages include so boundaries are clear without implying implemented behavior?
- Should command tests assert exact complete output or stable substrings/fields to balance regression protection with maintainability?

## Evidence Pointers

- `studies/go-cli-study/reports/final/01-project-structure.md`: Inspect Thin CLI Entry Point, internal package protection, unidirectional dependency flow, and anti-patterns around large entrypoint files.
- `studies/go-cli-study/reports/final/02-command-architecture.md`: Inspect Thin-Delegate Pattern, factory command creation, central config/context object, and RunE-size warnings.
- `studies/go-cli-study/reports/final/03-dependency-injection.md`: Inspect centralized composition root, constructor injection, functional options, lazy initialization, and global state cautions.
- `studies/go-cli-study/reports/final/05-error-handling.md`: Inspect user vs operational separation, hint-based messages, exit code mapping, and `%w` wrapping guidance.
- `studies/go-cli-study/reports/final/06-io-abstraction.md`: Inspect IOStreams test constructors, functional IO options, UI interfaces, and direct `os.*` anti-patterns.
- `studies/go-cli-study/reports/final/09-terminal-ux.md`: Inspect CLI-first output, non-TTY fallback, help/template behavior, and no-TTY/color cautions.
- `studies/go-cli-study/reports/final/11-testing-strategy.md`: Inspect command-level integration, golden output checks, behavior assertions, and brittle implementation-detail tests.
- `studies/go-cli-study/reports/final/12-extensibility.md`: Inspect functional options, registry/plugin tradeoffs, and warnings about delaying plugin systems until scope requires them.
- `studies/go-cli-study/reports/final/15-philosophy.md`: Inspect explicit non-goals, simplicity vs extensibility tradeoffs, and accepted-complexity cautions.

## Handoff To Reasoning

- Use this handbook as evidence input for sprint architecture reasoning.
- Validate whether the observed patterns fit the sprint's constraints before selecting an implementation shape.
- Keep implementation decisions out of this handbook; decide them in the reasoning artifacts.
- Preserve the sprint non-goals unless later reasoning explicitly identifies a requirements conflict.
- Prefer the smallest command shell that satisfies buildability, help/version/unknown-command behavior, deterministic tests, and the required package skeletons.
