# Sprint Technical Handbook: 05 Study Listing Inspection

> Project: `ultraplan-go`
> Sprint: `05-study-listing-inspection`
> Source: `sprint-index.md`
> **Inputs Used:** `projects/ultraplan-go/sprints/05-study-listing-inspection/sprint-index.md`, `templates/technical-handbook.md`, `projects/ultraplan-go/project-index.md`, `projects/ultraplan-go/sprints/05-study-listing-inspection/requirements.md`, `projects/ultraplan-go/docs/PRD.md`, `projects/ultraplan-go/docs/TRD.md`, `projects/ultraplan-go/docs/ARCHITECTURE.md`, `studies/go-cli-study/reports/final/01-project-structure.md`, `studies/go-cli-study/reports/final/02-command-architecture.md`, `studies/go-cli-study/reports/final/05-error-handling.md`, `studies/go-cli-study/reports/final/06-io-abstraction.md`, `studies/go-cli-study/reports/final/09-terminal-ux.md`, `studies/go-cli-study/reports/final/11-testing-strategy.md`, `studies/go-cli-study/reports/final/13-security.md`, `studies/go-cli-study/reports/final/14-performance.md`, `studies/go-cli-study/reports/final/15-philosophy.md`

This handbook distills the studies and reports selected by `sprint-index.md` for sprint reasoning. It does not decide architecture or implementation.

## Selected Studies And Reports

| Study / Report | Path | Relevant Finding | Confidence |
| --- | --- | --- | --- |
| `go-cli-study / 01-project-structure` | `studies/go-cli-study/reports/final/01-project-structure.md` | Non-trivial Go CLIs keep entrypoints thin and put business logic behind `internal/` or a domain package; `cmd/` imports business logic and not the reverse (`cmd/restic/main.go:37-114`, `internal/restic/repository.go:18`, `cmd/root.go:112-127`). | high |
| `go-cli-study / 02-command-architecture` | `studies/go-cli-study/reports/final/02-command-architecture.md` | Strong command designs use factory functions, option structs, and thin command wrappers; `RunE` functions that grow past orchestration become a maintenance warning (`pkg/cmd/install.go:132-145`, `pkg/cmdutil/factory.go:16-43`, `opencode/cmd/root.go:49-183`). | high |
| `go-cli-study / 05-error-handling` | `studies/go-cli-study/reports/final/05-error-handling.md` | Mature CLIs preserve error chains with `%w`, use typed or sentinel errors for programmatic handling, and separate user-facing diagnostics from operational detail (`rclone/vfs/zip.go:21`, `helm/pkg/storage/driver/driver.go:27-36`, `age/cmd/age/tui.go:37-54`). | high |
| `go-cli-study / 06-io-abstraction` | `studies/go-cli-study/reports/final/06-io-abstraction.md` | Testable CLIs inject stdout/stderr and filesystem seams at boundaries; direct `os.Stdout` or `os.Stderr` paths create output-test gaps (`pkg/iostreams/iostreams.go:551-568`, `internal/ui/terminal.go:10-36`, `cmd/ls/ls.go:42`). | high |
| `go-cli-study / 09-terminal-ux` | `studies/go-cli-study/reports/final/09-terminal-ux.md` | CLI-first tools should keep output calm, stable, and scriptable; progress/TUI machinery is reserved for long-running or interactive workflows (`pkg/yqlib/color_print.go:7-9`, `internal/cmd/prompt.go:20-256`, `internal/cmd/prompt.go:351-360`). | medium |
| `go-cli-study / 11-testing-strategy` | `studies/go-cli-study/reports/final/11-testing-strategy.md` | High-confidence CLI changes are covered by behavior-focused unit tests, command-level tests, fixtures, and golden/output regression checks (`chezmoi/internal/cmd/main_test.go:64-174`, `helm/internal/test/test.go:43`, `go-task/task_test.go:166-169`). | high |
| `go-cli-study / 13-security` | `studies/go-cli-study/reports/final/13-security.md` | Trust boundaries should be explicit for paths, shell execution, secrets, and untrusted inputs; argument arrays, path validation, and redaction are recurring safety patterns (`pkg/registry/transport.go:37-41`, `pkg/surveyext/editor_manual.go:23`, `internal/config/json/validator.go:146`). | high |
| `go-cli-study / 14-performance` | `studies/go-cli-study/reports/final/14-performance.md` | Fast CLIs defer expensive work, avoid recursive or unbounded scans, stream large inputs, and bound concurrency when parallelism enters scope (`pkg/cmdutil/factory.go:27-42`, `yq/pkg/yqlib/stream_evaluator.go:78-113`, `gdu/pkg/analyze/parallel.go:13`). | high |
| `go-cli-study / 15-philosophy` | `studies/go-cli-study/reports/final/15-philosophy.md` | Coherent tools state non-goals, accept complexity deliberately, and use factory/interface patterns only where they serve a clear product constraint (`VISION.md:97`, `pkg/cmd/factory/default.go:26-46`, `internal/chezmoi/system.go:25-45`). | high |

## Relevant Patterns

- **Study-owned behavior behind a thin CLI:** The project-structure report repeatedly shows CLI entrypoints delegating to business logic while preserving one-way imports, with examples such as k9s `cmd/root.go:112-127`, restic `cmd/restic/main.go:37-114`, and yq `cmd/root.go:9` plus `cmd/evaluate_sequence_command.go:7`. For this sprint, the pressure is to keep source kind discovery, frontmatter handling, applicability filtering, and list output in the study module rather than spreading behavior into generic global packages.

- **Factory and option-shaped command construction:** The command architecture report identifies `newXxxCmd` factories and command-specific options as the common shape in mature CLIs, citing gh-cli `pkg/cmd/issue/list/list.go:47-118`, helm `pkg/cmd/install.go:132-145`, and restic `cmd/restic/cmd_backup.go:84-115`. This matters for preserving command testability when `study <study> list` output changes without letting command handlers absorb parsing, discovery, and filtering rules.

- **Injectable output for command tests:** The IO report highlights `IOStreams.Test()` in gh-cli (`pkg/iostreams/iostreams.go:551-568`), restic's `ui.Terminal` (`internal/ui/terminal.go:10-36`), and mitchellh-cli's `Ui` (`ui.go:19-43`) as patterns that make CLI output capturable. This sprint's list output is acceptance-critical, so output construction should remain observable through buffers or command-level seams.

- **Error wrapping with source-path context:** The error-handling report treats `%w` wrapping as a baseline pattern, with examples in rclone `vfs/zip.go:21`, gh-cli `pkg/ssh/ssh_keys.go:64`, helm `pkg/storage/driver/secrets.go:76`, and opencode `internal/llm/agent/agent.go:238`. Malformed frontmatter and invalid `applicable_dimensions` values need enough context to identify the Markdown source path and offending value while preserving lower-level YAML or filesystem causes.

- **Behavior-focused fixture tests and output regression checks:** The testing report shows CLI integration via testscript (`chezmoi/internal/cmd/main_test.go:64-174`, `gh-cli/acceptance/acceptance_test.go:26-29`), golden comparisons (`helm/internal/test/test.go:43`, `go-task/task_test.go:166-169`), and fixture-driven tests (`gdu/internal/testdir/test_dir.go:9`). This sprint has deterministic discovery and listing requirements, making local filesystem fixtures and explicit output assertions more relevant than private implementation assertions.

- **Direct-child discovery instead of broad scans:** The performance report warns against unbounded accumulation and unnecessary traversal, while highlighting lazy init and streaming (`pkg/cmdutil/factory.go:27-42`, `yq/pkg/yqlib/stream_evaluator.go:78-113`). This sprint's evidence pressure is to inspect only direct children of `studies/<study>/sources/` for listing, not recurse through source repositories.

- **Explicit trust boundaries for user-authored files:** The security report shows schema/config validation and safe path handling as recurring trust-boundary tools, including k9s schema validation (`internal/config/json/validator.go:146`) and shell-safe argument construction (`cmd_obj_builder.go:38`, `pkg/surveyext/editor_manual.go:23`). Markdown source files and frontmatter are user-authored inputs, so parsing and diagnostics need to treat them as untrusted data rather than trusted internal metadata.

- **CLI-first output instead of interactive UX:** The terminal UX report distinguishes CLI-first tools from full TUI tools and points to non-TTY fallback and stable output patterns (`internal/cmd/prompt.go:20-256`, `cmd/evaluate_sequence_command.go:67`). Because this sprint only changes listing/inspection, the evidence supports stable plain text over progress bars, spinners, prompts, or TUI machinery.

## Trade-Offs

| Trade-Off | Benefit | Cost | When It Matters |
| --- | --- | --- | --- |
| Keep source discovery concrete in `internal/study` vs introduce filesystem abstraction now | Minimal change and aligns with thin CLI plus module-owned behavior evidence (`cmd/root.go:112-127`, `internal/restic/repository.go:18`) | Harder to unit-test exotic filesystem failures unless tests use real temp dirs | Relevant because sprint scope is local direct-child discovery, not reusable filesystem infrastructure |
| Parse frontmatter during source discovery vs defer parsing until filtering/listing | One source model can expose kind and applicability consistently, reducing repeated reads | Discovery can fail or warn because of Markdown metadata, increasing command error surface | Relevant because list output must show applicability and invalid values must include source path context |
| Text listing with explicit annotations vs adding/stabilizing JSON output | Keeps current command script-friendly and avoids premature output schema commitments | Text parsing by scripts can be brittle if formatting changes | Relevant because requirements call for deterministic text output and JSON schema stabilization is a non-goal |
| Directory sources always applicable vs one unified frontmatter applicability model for every source | Preserves repository source behavior and keeps applicability rules simple | Two source kinds have asymmetric filtering rules | Relevant because Markdown document sources need filters while directory repositories remain baseline sources |
| Golden-style output assertions vs flexible substring assertions | Catches unintended formatting drift in user-facing list output, matching testing evidence from helm and go-task (`helm/internal/test/test.go:43`, `go-task/task_test.go:166-169`) | Requires intentional fixture updates when output legitimately changes | Relevant because the sprint's visible behavior is deterministic command output |
| Leading-only frontmatter parsing vs permissive frontmatter scanning anywhere in a document | Avoids treating content examples as metadata and keeps document body predictable | Users with non-leading metadata blocks get no applicability behavior | Relevant because requirements constrain parsing to leading `---` blocks and stripping must leave body intact |

## Anti-Patterns And Warnings

- **Do not let `study <study> list` become a large command handler:** The command report calls out large command functions as a caution, including opencode `cmd/root.go:49-183`, yq `cmd/evaluate_sequence_command.go:152`, and age `cmd/age/age.go:105-321`. Keep command code as wiring and presentation while discovery and applicability remain study logic.

- **Do not recurse into source repositories for listing:** The performance report warns against unbounded in-memory accumulation and unnecessary traversal, while the sprint requires direct child inspection only. Nested Markdown files inside directory sources should remain repository content, not separate study sources.

- **Do not break error chains or replace causes with string-only diagnostics:** The error report explicitly flags `fmt.Errorf` without `%w` in mitchellh-cli `cli.go:205-206` and inconsistent wrapping in fzf `src/server.go:68`. Frontmatter, filesystem, and path errors should retain underlying causes.

- **Do not hardcode stdout/stderr in paths that command tests need to inspect:** The IO report flags direct output such as rclone `cmd/ls/ls.go:42`, age `cmd/age/tui.go:31`, and chezmoi template stderr bypass `internal/cmd/templatefuncs.go:296`. Listing output should remain captured by existing command test seams.

- **Do not add runtime, prompt, scheduler, report-validation, or synthesis behavior as part of listing:** The philosophy report warns that complexity should be accepted deliberately, not accumulated. Sprint scope is inspection only; runtime and scheduling concerns belong to later reasoning.

- **Do not treat malformed user-authored Markdown metadata as trusted internal state:** The security report emphasizes explicit trust boundaries and validation for untrusted inputs (`internal/config/json/validator.go:146`, `pkg/yqlib/security_prefs.go:3-7`). Applicability metadata needs validation with source context.

- **Do not assert private implementation details in tests when user behavior is the contract:** The testing report identifies hint-count assertions as brittle (`internal/view/pod_test.go:23`). Prefer assertions on discovered sources, normalized applicability, errors, and rendered listing output.

## Examples Worth Inspecting

| Example | Path / Source | Why It Is Useful |
| --- | --- | --- |
| Thin CLI delegating into module logic | `cmd/root.go:112-127`, `cmd/restic/main.go:37-114`, `cmd/root.go:9` | Shows command entrypoints staying small while domain packages own behavior. |
| Command factory and options pattern | `pkg/cmd/issue/list/list.go:47-118`, `pkg/cmd/install.go:132-145`, `cmd/restic/cmd_backup.go:84-115` | Useful when reasoning about how `study <study> list` should wire output without absorbing business logic. |
| Capturable command IO | `pkg/iostreams/iostreams.go:551-568`, `ui.go:19-43`, `internal/ui/terminal.go:10-36` | Shows how tests can capture output without a real terminal. |
| Error wrapping and user-facing hints | `rclone/vfs/zip.go:21`, `age/cmd/age/tui.go:37-54`, `gh-cli/internal/ghcmd/cmd.go:281-301` | Useful for malformed frontmatter diagnostics that need both cause preservation and actionable messages. |
| Fixture and output regression tests | `chezmoi/internal/cmd/main_test.go:64-174`, `helm/internal/test/test.go:43`, `go-task/task_test.go:166-169` | Provides evidence for command-level fixtures and deterministic output checks. |
| Security validation boundary | `internal/config/json/validator.go:146`, `pkg/yqlib/security_prefs.go:3-7`, `pkg/registry/transport.go:37-41` | Shows validation/redaction approaches for user-provided data and sensitive diagnostics. |
| Lazy initialization and bounded scanning | `pkg/cmdutil/factory.go:27-42`, `yq/pkg/yqlib/stream_evaluator.go:78-113`, `gdu/pkg/analyze/parallel.go:13` | Useful when avoiding startup cost and accidental broad repository scans. |

## Design Pressures

- Preserve `internal/study` ownership while extending the model for source kind, frontmatter metadata, and applicability.

- Keep `ultraplan study <study> list` deterministic, human-readable, and script-friendly without committing to a new JSON schema.

- Balance source discovery simplicity against enough metadata parsing to show Markdown applicability in listing output.

- Treat Markdown files as user-authored, potentially malformed inputs with source-path diagnostics.

- Avoid recursive scans into potentially large source repositories; direct-child discovery is enough for the sprint's inspection scope.

- Preserve directory source behavior while adding Markdown-specific applicability filters.

- Keep tests focused on public behavior: discovered source order, ignored entries, normalized dimensions, applicability filtering, error messages, and listing output.

- Avoid introducing runtime execution semantics, prompt composition, report validation, summary generation, or scheduler changes while still leaving future flows able to reuse applicability helpers.

## Open Questions For Reasoning

- Should malformed leading frontmatter fail `study <study> list`, produce a validation-style warning, or be handled differently depending on command context?

- What exact text format should represent `directory`, unfiltered Markdown applicability, and filtered Markdown applicability while remaining stable for command tests?

- Should Markdown source names include the `.md` extension, strip it, or follow the existing source name derivation rules?

- If a Markdown source declares duplicate applicable dimensions, should the normalized list deduplicate silently, preserve duplicates, or report a validation issue?

- How should ordering compare a directory and Markdown file that normalize to the same display name?

- Should invalid `applicable_dimensions` values reject the whole source list or only the offending source?

- Which helper surface best supports future flows without overexposing internals: `GetApplicableSources`, a method on `Study`, or an unexported helper used through existing public service methods?

- How should inapplicable source/dimension pairs be represented in listing-derived checks without implying a runtime skipped task for this sprint?

## Evidence Pointers

- `studies/go-cli-study/reports/final/01-project-structure.md`: Inspect Pattern 1, Pattern 2, Pattern 4, and the Evidence Index for thin CLI, `internal/` protection, and dependency direction.

- `studies/go-cli-study/reports/final/02-command-architecture.md`: Inspect Pattern 1, Pattern 2, and Anti-Patterns for command factory design and command-handler size warnings.

- `studies/go-cli-study/reports/final/05-error-handling.md`: Inspect Pattern 1, Pattern 4, Pattern 7, and Anti-Patterns for wrapped errors, user diagnostics, hints, and panic/string-only failures.

- `studies/go-cli-study/reports/final/06-io-abstraction.md`: Inspect IOStreams, `Terminal`, and `Ui` examples for output capture and filesystem seam trade-offs.

- `studies/go-cli-study/reports/final/09-terminal-ux.md`: Inspect CLI-first and non-TTY fallback sections to keep listing output stable and non-interactive.

- `studies/go-cli-study/reports/final/11-testing-strategy.md`: Inspect testscript, golden files, fixture extraction, and brittle-test warnings for the source and command test approach.

- `studies/go-cli-study/reports/final/13-security.md`: Inspect shell-safe argument construction, schema validation, secret redaction, and trust-boundary sections for path and frontmatter safety reasoning.

- `studies/go-cli-study/reports/final/14-performance.md`: Inspect lazy initialization, streaming, and anti-patterns around recursive/unbounded work for direct-child discovery constraints.

- `studies/go-cli-study/reports/final/15-philosophy.md`: Inspect explicit non-goals, factory/interface patterns, and complexity warnings to avoid scope creep.

## Handoff To Reasoning

- Use this handbook as evidence input.
- Validate whether the observed patterns fit this project's constraints.
- Do not copy external patterns without sprint-specific reasoning.
- Resolve frontmatter error behavior, exact listing text, source naming, duplicate normalization, and helper export shape in sprint reasoning before implementation planning.
