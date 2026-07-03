# Sprint Technical Handbook: Project Domain and Project Index

> Project: `ultraplan-go`
> Sprint: `16-project-domain-and-index`
> Source: `sprint-index.md`
> **Inputs Used:** `projects/ultraplan-go/sprints/16-project-domain-and-index/sprint-index.md`, `projects/ultraplan-go/project-index.md`, `templates/technical-handbook.md`, `projects/ultraplan-go/sprints/16-project-domain-and-index/requirements.md`, `projects/ultraplan-go/docs/PRD.md`, `projects/ultraplan-go/docs/TRD.md`, `projects/ultraplan-go/docs/ARCHITECTURE.md`, `studies/go-cli-study/reports/final/01-project-structure.md`, `studies/go-cli-study/reports/final/02-command-architecture.md`, `studies/go-cli-study/reports/final/03-dependency-injection.md`, `studies/go-cli-study/reports/final/05-error-handling.md`, `studies/go-cli-study/reports/final/06-io-abstraction.md`, `studies/go-cli-study/reports/final/09-terminal-ux.md`, `studies/go-cli-study/reports/final/11-testing-strategy.md`, `studies/go-cli-study/reports/final/13-security.md`, `studies/go-cli-study/reports/final/14-performance.md`, `studies/go-cli-study/reports/final/15-philosophy.md`

This handbook distills the studies and reports selected by `sprint-index.md` for sprint reasoning. It does not decide architecture or implementation. Project requirements and docs were used only for sprint scope context; the technical evidence below comes from the selected evidence reports.

## Selected Studies And Reports

| Study / Report | Path | Relevant Finding | Confidence |
| --- | --- | --- | --- |
| `01-project-structure` | `studies/go-cli-study/reports/final/01-project-structure.md` | Non-trivial Go CLIs consistently keep entrypoints thin and business logic inside protected or domain-owned packages; evidence includes `cmd/gh/main.go:6`, `cmd/root.go:112-127`, `internal/restic/repository.go:18`, and `cmd/evaluate_sequence_command.go:7`. | high |
| `02-command-architecture` | `studies/go-cli-study/reports/final/02-command-architecture.md` | Maintainable command systems use factory-created commands that parse flags, call services, and format output; evidence includes `pkg/cmd/install.go:132-145`, `pkg/cmdutil/factory.go:16-43`, and `internal/cmd/config.go:1833-1845`. | high |
| `03-dependency-injection` | `studies/go-cli-study/reports/final/03-dependency-injection.md` | Manual dependency construction with explicit factories and narrow interfaces dominates; no studied repo used a DI framework, and global state was the main quality divider. Evidence includes `pkg/cmd/factory/default.go:26-46`, `executor.go:22-24`, and `fs/config.go:793`. | high |
| `05-error-handling` | `studies/go-cli-study/reports/final/05-error-handling.md` | Strong CLIs wrap errors with `%w`, preserve programmatic classification, separate user rendering from operational detail, and map exit codes deliberately. Evidence includes `cmd/age/tui.go:47-54`, `go-task/errors/errors_task.go:13-32`, and `internal/ghcmd/cmd.go:44-49`. | high |
| `06-io-abstraction` | `studies/go-cli-study/reports/final/06-io-abstraction.md` | Testable CLI output and filesystem behavior come from injectable `io.Reader`/`io.Writer` streams, filesystem interfaces, and test constructors. Evidence includes `pkg/iostreams/iostreams.go:551-568`, `internal/chezmoi/system.go:25`, and `internal/ui/mock.go:10-53`. | high |
| `09-terminal-ux` | `studies/go-cli-study/reports/final/09-terminal-ux.md` | CLI-first tools should favor concise, non-TTY-safe, scriptable output over rich terminal rendering; prompt/progress machinery is justified mainly for interactive or long operations. Evidence includes `internal/cmd/prompt.go:124-137`, `pkg/iostreams/iostreams.go:116-130`, and `cmd/root.go:213-215`. | medium |
| `11-testing-strategy` | `studies/go-cli-study/reports/final/11-testing-strategy.md` | Durable CLI behavior is protected by table-driven tests, command-level or script-level integration tests, fixture projects, golden/output assertions, and mock/fake infrastructure. Evidence includes `internal/cmd/main_test.go:64-174`, `pkg/httpmock/stub.go:35-199`, and `cmd/bisync/bisync_test.go:1435-1479`. | high |
| `13-security` | `studies/go-cli-study/reports/final/13-security.md` | Trust boundaries must be explicit for paths, subprocesses, credentials, and config; strong examples use safe path/command construction, schema validation, secret redaction, and permission gates. Evidence includes `internal/config/json/validator.go:146`, `pkg/registry/transport.go:37-41`, and `internal/llm/tools/bash.go:41-55`. | high |
| `14-performance` | `studies/go-cli-study/reports/final/14-performance.md` | Fast CLIs defer expensive initialization, stream or bound large work, and avoid unbounded recursive processing. Evidence includes `pkg/cmdutil/factory.go:27-42`, `gdu/pkg/analyze/parallel.go:13`, and `yq/pkg/yqlib/stream_evaluator.go:78-113`. | high |
| `15-philosophy` | `studies/go-cli-study/reports/final/15-philosophy.md` | Coherent tools accept complexity deliberately and reject features outside their purpose; evidence includes age's no-keyring stance at `age.go:18`, lazygit's explicit VISION principles at `VISION.md:5` and `VISION.md:97`, and Helm's plugin-first guidance at `AGENTS.md:88`. | medium |

## Relevant Patterns

- **Module-owned domain behavior with thin entrypoints:** The project command layer should remain a thin adapter that parses arguments, delegates, renders, and maps exit codes, because the studied CLIs repeatedly separate command wiring from business logic. Evidence includes the thin CLI thesis in `01-project-structure`, `cmd/gh/main.go:6`, `cmd/root.go:112-127`, `pkg/cmd/install.go:132-145`, and `pkg/cmdutil/factory.go:16-43`.
- **Concrete project package before shared abstractions:** The evidence favors `internal/` or domain packages for application-specific behavior and cautions against global technical-layer packages. `01-project-structure` highlights Go-enforced `internal/` boundaries using `internal/chezmoi/chezmoi.go:1-2`, `internal/restic/repository.go:18`, and `cmd/evaluate_sequence_command.go:7`; `15-philosophy` frames deliberate rejection of unnecessary extensibility as a strength with `age.go:18` and `VISION.md:97`.
- **Factory-created commands and use-case service calls:** Command constructors with injected dependencies make commands testable without embedding business rules. `02-command-architecture` cites `NewCmdList(f *cmdutil.Factory)` at `pkg/cmd/issue/list/list.go:47`, `newInstallCmd(cfg *action.Configuration)` at `pkg/cmd/install.go:132`, and `newBackupCommand(gopts *global.Options)` at `cmd/restic/cmd_backup.go:35`.
- **Explicit manual dependency construction:** Manual DI through constructors, factories, and small interfaces is the dominant idiom. `03-dependency-injection` cites `pkg/cmd/factory/default.go:26-46`, `cmdutil.Factory` at `pkg/cmdutil/factory.go:16-43`, `ExecutorOption` at `executor.go:22-24`, and notes no studied repos used Wire/Fx/Dig.
- **Injectable IO and filesystem seams for deterministic tests:** CLI stdout/stderr and filesystem access should be capturable in tests instead of hardcoded to `os.Stdout`, `os.Stderr`, or raw filesystem calls. `06-io-abstraction` cites `pkg/iostreams/iostreams.go:551-568`, `cmd/gdu/app/app.go:30-49`, `internal/chezmoi/system.go:25`, and `internal/fs/interface.go:10-31`.
- **Actionable, classified validation errors:** Validation failures should preserve causes internally and expose safe, user-repairable context. `05-error-handling` cites `%w` wrapping in `pkg/ssh/ssh_keys.go:64`, `TaskNotFoundError` suggestions at `errors/errors_task.go:13-32`, user hints in `cmd/age/tui.go:47-54`, and exit code mapping in `internal/ghcmd/cmd.go:44-49`.
- **Fail-closed trust boundaries around paths and catalog references:** Catalog validation has a trust-boundary shape similar to config/schema validation and secure command/path handling. `13-security` cites schema validation in `internal/config/json/validator.go:146`, explicit security flags in `pkg/yqlib/security_prefs.go:3-7`, credential scrubbing in `pkg/registry/transport.go:37-41`, and caution around missing path canonicalization in `internal/common/ignore.go:16-37`.
- **Concise CLI-first terminal UX:** This sprint's commands are status/list/validate operations, not interactive TUIs, so the relevant evidence is CLI-first output with non-TTY behavior rather than spinners or prompts. `09-terminal-ux` cites non-TTY fallback in `internal/cmd/prompt.go:124-137`, color/theme adaptation in `pkg/iostreams/iostreams.go:116-130`, and warnings against blocking prompts in `pkg/action/package.go:200-208`.
- **Behavior-level and fixture-driven tests:** The selected reports favor tests that validate command behavior, outputs, exit codes, and fixture filesystem states. `11-testing-strategy` cites testscript integration at `internal/cmd/main_test.go:64-174`, HTTP mocking at `pkg/httpmock/stub.go:35-199`, vfst filesystem isolation at `internal/chezmoitest/chezmoitest.go:86-92`, and golden comparison at `helm/internal/test/test.go:43`.
- **Bounded discovery and lazy work:** Project discovery/status should avoid recursive evidence/source scans and defer expensive work until requested. `14-performance` cites lazy factory fields at `pkg/cmdutil/factory.go:27-42`, concurrency bounding at `gdu/pkg/analyze/parallel.go:13`, and streaming or incremental processing at `yq/pkg/yqlib/stream_evaluator.go:78-113`.

## Trade-Offs

| Trade-Off | Benefit | Cost | When It Matters |
| --- | --- | --- | --- |
| Project-owned parser/validator in `internal/project` vs generalized catalog/validation packages | Keeps behavior near project state and avoids premature abstraction, matching `internal/` and module-owned evidence from `internal/chezmoi/chezmoi.go:1-2` and `internal/restic/repository.go:18` | Some parsing mechanics may be duplicated later if sprint artifacts add similar tables | Sprint reasoning must decide how narrow the project-index parser stays without blocking future sprint validators |
| Thin CLI handlers vs self-contained command logic | Easier tests, clearer ownership, and lower risk of command-layer business rules, as shown by `pkg/cmd/install.go:132-145` and `pkg/cmd/issue/list/list.go:47-118` | Requires service/domain types before commands feel complete | Matters for `project list`, `status`, and `validate`, where output and exit mapping should not own catalog rules |
| Explicit constructor/factory wiring vs package-level globals | Traceable dependencies and isolated command tests, supported by `pkg/cmdutil/factory.go:16-43` and `pkg/cmd/factory/default.go:26-46` | More boilerplate than global config or singleton access | Matters when commands need workspace, output, and service dependencies but must not call runtime or study services |
| Small local interfaces/fakes vs real filesystem tests only | Improves deterministic tests for missing files, path escapes, malformed rows, and output capture, supported by `pkg/iostreams/iostreams.go:551-568` and `internal/chezmoi/system.go:25` | Interfaces can become inaccurate if they abstract too much too early | Matters for path-safety, project fixture, and command-output tests |
| Structured validation findings vs plain error strings | Enables exit code `5`, sorted diagnostics, section/name/path/action fields, and safe user messages, aligning with `errors/errors.go:47-50`, `internal/ghcmd/cmd.go:44-49`, and `cmd/age/tui.go:47-54` | More domain modeling than returning a single error | Matters for catalog validation where multiple findings should be reported together |
| Fail-closed path resolution vs permissive missing/escaping references | Prevents workspace escape and silent invalid catalogs, consistent with security evidence around schema validation and canonicalization cautions (`internal/config/json/validator.go:146`, `internal/common/ignore.go:16-37`) | Users may see more validation failures during setup | Matters for every catalog path in `project-index.md`, especially selected evidence and contract paths |
| Full command integration tests vs only package unit tests | Catches help, output ordering, stdout/stderr, and exit-code regressions, as shown by testscript and command integration evidence (`internal/cmd/main_test.go:64-174`, `acceptance/acceptance_test.go:26-29`) | More fixture and test helper maintenance | Matters because acceptance criteria explicitly cover CLI help, deterministic text output, and validation exit codes |
| Lazy/bounded scanning vs eager full workspace scan | Keeps project commands fast and avoids accidentally walking study source repositories, matching lazy/bounded performance patterns (`pkg/cmdutil/factory.go:27-42`, `gdu/pkg/analyze/parallel.go:13`) | Some health information remains unavailable until specific validation runs | Matters for `project list` and `project status`, which should inspect direct project/sprint children and selected catalog paths only |

## Anti-Patterns And Warnings

- **Do not let `internal/app` become the project domain:** `02-command-architecture` warns that long `RunE` functions and command-contained business logic are weak architecture, citing `opencode/cmd/root.go:49-183`, `cmd/age/age.go:105-321`, and `yq/cmd/evaluate_sequence_command.go:152`.
- **Do not introduce broad shared packages just because catalog parsing resembles validation:** `01-project-structure` and `15-philosophy` warn against generic layers and accidental complexity; caution examples include monolithic or over-broad packages such as `internal/cmd/config.go:1`, `src/terminal.go`, and global package patterns.
- **Do not consume study services or models for project validation:** The selected structure/philosophy evidence supports domain boundaries and deliberate feature rejection; study reports may be catalog paths, but project behavior should not reuse study source/dimension/report semantics just because those artifacts are in the workspace.
- **Do not hardcode `os.Stdout`, `os.Stderr`, or real filesystem paths into test-sensitive code:** `06-io-abstraction` identifies hardcoded output as a testability smell in `cmd/age/tui.go:31`, `cmd/ls/ls.go:42`, `pkg/yqlib/logger.go:18`, and `internal/cmd/templatefuncs.go:296`.
- **Do not flatten all validation failures into one unclassified error:** `05-error-handling` warns against missing `%w`, no sentinels, magic exit numbers, and string-only error handling, citing `mitchellh-cli/cli.go:205-206`, `dive/image/docker/manifest.go:18`, and `mitchellh-cli/command.go:10-11`.
- **Do not silently ignore bad catalog rows or unresolved paths:** Security evidence treats untrusted inputs as explicit boundaries; `13-security` cautions that no path canonicalization before comparison (`internal/common/ignore.go:16-37`) and default-allow dangerous operations (`pkg/yqlib/security_prefs.go:3-7`) create brittle safety models.
- **Do not add interactive prompts, spinners, color, or TUI mechanics to script-first project commands by default:** `09-terminal-ux` shows rich UX pays off for interactive tools, but also warns that prompts can hang CI (`pkg/action/package.go:200-208`) and color/progress must respect TTY state.
- **Do not recursively scan evidence source trees or repository contents for status:** `14-performance` warns against unbounded in-memory accumulation and goroutine-per-subdir patterns, citing `gdu/pkg/analyze/parallel.go:13` as the bounding pattern and `k9s/internal/model/table.go:200` as a memory-pressure caution.
- **Do not make tests assert incidental internals instead of user-visible behavior:** `11-testing-strategy` flags hint-count tests (`internal/view/pod_test.go:23`) and fragile regex output (`pkg/cmd/issue/list/list_test.go:88-91`) as caution signs.

## Examples Worth Inspecting

| Example | Path / Source | Why It Is Useful |
| --- | --- | --- |
| Thin command delegation | `helm/pkg/cmd/install.go:132-145`, cited by `02-command-architecture` | Shows command creation and delegation without making command code own the action's business rules. |
| Command dependency factory | `gh-cli/pkg/cmdutil/factory.go:16-43`, cited by `02-command-architecture` and `03-dependency-injection` | Shows a dependency carrier that enables test injection and lazy resolution. |
| Command-level output capture | `gh-cli/pkg/iostreams/iostreams.go:551-568`, cited by `06-io-abstraction` | Shows a production/test IO split suitable for deterministic command tests. |
| Filesystem interface seam | `chezmoi/internal/chezmoi/system.go:25`, cited by `06-io-abstraction` and `15-philosophy` | Shows filesystem operations as an interface with debug/dry-run/test implementations. |
| User-facing hints | `age/cmd/age/tui.go:47-54`, cited by `05-error-handling` | Shows actionable error rendering without losing lower-level error semantics. |
| Exit code classification | `gh-cli/internal/ghcmd/cmd.go:44-49`, cited by `05-error-handling` | Shows script-relevant exit categories as named constants rather than ad hoc numbers. |
| Schema validation boundary | `k9s/internal/config/json/validator.go:146`, cited by `13-security` | Shows validation as a trust boundary before config-like data is trusted. |
| Bounded project-like scanning | `gdu/pkg/analyze/parallel.go:13`, cited by `14-performance` | Shows explicit resource bounds for filesystem traversal patterns. |
| CLI integration testing | `chezmoi/internal/cmd/main_test.go:64-174`, cited by `11-testing-strategy` | Shows script-style command behavior tests with controlled filesystem state. |
| Golden output comparison | `helm/internal/test/test.go:43`, cited by `11-testing-strategy` | Shows a pattern for guarding user-facing text output changes when stable output matters. |

## Design Pressures

- Project behavior must be module-owned and not drift into app-level command handlers or shared technical packages.
- Project index parsing must be concrete enough to parse current catalog tables while avoiding a speculative general Markdown table framework.
- Catalog validation must report enough structured detail for repair without leaking unsafe paths or secrets.
- Project commands are script-friendly status/list/validate commands, so deterministic plain text and exit codes matter more than rich terminal UX.
- Discovery and status must be bounded to direct project roots, direct sprint directories, docs, and catalog references, not recursive source/evidence scans.
- Tests must cover both package behavior and CLI behavior because acceptance criteria include help, output order, exit codes, and runtime-free operation.
- Dependency wiring should preserve clear product/platform boundaries while keeping simple filesystem and output seams testable.
- Validation should aggregate multiple findings where useful, but sprint reasoning must decide the exact result model and rendering policy.

## Open Questions For Reasoning

- What is the smallest project-domain model that can express projects, project indexes, catalog entries, status, and validation findings without becoming a planning or study abstraction?
- Should the project-index parser support only the exact catalog table headings used today, or allow a narrow amount of heading variation for future project indexes?
- What structured fields must a validation finding carry to satisfy actionable diagnostics: section, entry name, referenced path, severity, cause, and suggestion?
- How should path resolution model explicitly external entries, if any, while failing closed for normal workspace-relative paths?
- Where should stdout/stderr and exit code mapping live so `internal/app` remains thin but command tests can assert exact output?
- How much fixture infrastructure is justified for this sprint: package-level fixtures only, command-level fixtures, or golden files for help/output?
- Should status perform catalog parsing/validation or only report file presence plus a summarized health from explicit validation logic?
- What deterministic ordering rules should apply across projects, docs, catalog entries, diagnostics, and sprint directories?

## Evidence Pointers

- `studies/go-cli-study/reports/final/01-project-structure.md`: Inspect Pattern 1, Pattern 2, Pattern 4, tradeoffs around `internal/` vs `pkg/`, and cautions on monolithic packages and bidirectional imports.
- `studies/go-cli-study/reports/final/02-command-architecture.md`: Inspect Thin-Delegate Pattern, factory function command creation, options structs, and cautions on long `RunE` functions.
- `studies/go-cli-study/reports/final/03-dependency-injection.md`: Inspect centralized composition root, constructor injection, functional options, interface abstraction, lazy initialization, and global-state cautions.
- `studies/go-cli-study/reports/final/05-error-handling.md`: Inspect `%w` wrapping, typed errors, sentinel errors, user vs operational rendering, hint-based errors, multi-error aggregation, and exit code mapping.
- `studies/go-cli-study/reports/final/06-io-abstraction.md`: Inspect IOStreams/Test constructors, System/FS interfaces, MockTerminal/MockUi, and hardcoded `os.*` cautions.
- `studies/go-cli-study/reports/final/09-terminal-ux.md`: Inspect CLI-first vs TUI tradeoffs, non-TTY fallback, color/TTY checks, and warnings against blocking prompts in CI.
- `studies/go-cli-study/reports/final/11-testing-strategy.md`: Inspect testscript integration, centralized mocks, golden files, fixture extraction, behavior assertions, and anti-patterns around fragile internal assertions.
- `studies/go-cli-study/reports/final/13-security.md`: Inspect schema validation, trust boundaries, secret redaction, safe command/path handling, and path canonicalization cautions.
- `studies/go-cli-study/reports/final/14-performance.md`: Inspect lazy initialization, bounded traversal/concurrency, streaming, and anti-patterns around eager initialization and unbounded accumulation.
- `studies/go-cli-study/reports/final/15-philosophy.md`: Inspect explicit non-goals, factory/options patterns, interface-driven abstractions, deliberate complexity, and cautions on god structs, plugin/security surfaces, and configuration bloat.

## Handoff To Reasoning

- Use this handbook as evidence input.
- Validate whether the observed patterns fit this project's constraints.
- Do not copy external patterns without sprint-specific reasoning.
- Resolve the open questions in area-specific reasoning or final sprint reasoning before implementation.
