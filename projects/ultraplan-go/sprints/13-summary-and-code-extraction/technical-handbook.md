# Sprint Technical Handbook: 13 Summary And Code Extraction

> Project: `ultraplan-go`
> Sprint: `13-summary-and-code-extraction`
> Source: `sprint-index.md`
> **Inputs Used:** `projects/ultraplan-go/sprints/13-summary-and-code-extraction/sprint-index.md`, `templates/technical-handbook.md`, `projects/ultraplan-go/project-index.md`, `projects/ultraplan-go/sprints/13-summary-and-code-extraction/requirements.md`, `projects/ultraplan-go/docs/PRD.md`, `projects/ultraplan-go/docs/TRD.md`, `projects/ultraplan-go/docs/ARCHITECTURE.md`, `studies/go-cli-study/reports/final/01-project-structure.md`, `studies/go-cli-study/reports/final/02-command-architecture.md`, `studies/go-cli-study/reports/final/03-dependency-injection.md`, `studies/go-cli-study/reports/final/05-error-handling.md`, `studies/go-cli-study/reports/final/06-io-abstraction.md`, `studies/go-cli-study/reports/final/09-terminal-ux.md`, `studies/go-cli-study/reports/final/10-logging-observability.md`, `studies/go-cli-study/reports/final/11-testing-strategy.md`, `studies/go-cli-study/reports/final/13-security.md`, `studies/go-cli-study/reports/final/14-performance.md`, `studies/go-cli-study/reports/final/15-philosophy.md`

This handbook distills the studies and reports selected by `sprint-index.md` for sprint reasoning. It does not decide architecture or implementation.

## Selected Studies And Reports

| Study / Report | Path | Relevant Finding | Confidence |
| --- | --- | --- | --- |
| `01-project-structure` | `studies/go-cli-study/reports/final/01-project-structure.md` | Mature Go CLIs keep entrypoints thin and put product behavior behind `cmd/` or `internal/`; examples include chezmoi `main.go:16` delegating into `internal/cmd`, k9s `cmd/root.go:112-127`, and gh-cli `cmd/gh/main.go:6`. | high |
| `02-command-architecture` | `studies/go-cli-study/reports/final/02-command-architecture.md` | Command factories and thin `RunE` handlers keep CLI wiring testable; gh-cli `pkg/cmdutil/factory.go:16-43`, helm `pkg/cmd/install.go:132-145`, and restic `cmd/restic/main.go:37-114` are repeated evidence. | high |
| `03-dependency-injection` | `studies/go-cli-study/reports/final/03-dependency-injection.md` | Manual DI with a visible composition root and narrow test seams dominates; gh-cli `pkg/cmd/factory/default.go:26-46`, go-task `executor.go:22-24`, and restic `internal/restic/repository.go:18-66` show factory, options, and interface seams. | high |
| `05-error-handling` | `studies/go-cli-study/reports/final/05-error-handling.md` | Mature CLIs preserve error chains and map user-facing outcomes; rclone `fs/fserrors/error.go:22-29`, go-task `errors/errors.go:47-50`, gh-cli `internal/ghcmd/cmd.go:44-49`, and age `cmd/age/tui.go:37-54` show classification, exit mapping, and hints. | high |
| `06-io-abstraction` | `studies/go-cli-study/reports/final/06-io-abstraction.md` | Output and filesystem seams enable deterministic command tests; gh-cli `pkg/iostreams/iostreams.go:551-568`, go-task `executor.go:553-564`, and restic `internal/ui/mock.go:10-53` show in-memory IO and terminal fakes. | high |
| `09-terminal-ux` | `studies/go-cli-study/reports/final/09-terminal-ux.md` | CLI-first commands need concise deterministic output and non-TTY behavior more than rich TUI; chezmoi `internal/cmd/prompt.go:20-256`, gh-cli `internal/prompter/`, and yq `pkg/yqlib/color_print.go:7-9` illustrate output-mode fit. | medium |
| `10-logging-observability` | `studies/go-cli-study/reports/final/10-logging-observability.md` | Diagnostics should be separate from data output and structured when useful; gh-cli `pkg/iostreams/iostreams.go:52-54`, helm `internal/logging/logging.go:31-66`, and k9s `internal/slogs/keys.go:6-231` show stderr/data separation and consistent keys. | medium |
| `11-testing-strategy` | `studies/go-cli-study/reports/final/11-testing-strategy.md` | CLI confidence comes from fixture-backed command tests, golden output, and mocks; chezmoi `internal/cmd/main_test.go:64-174`, helm `internal/test/test.go:43`, go-task `task_test.go:166-169`, and gh-cli `acceptance/acceptance_test.go:26-29` are strong models. | high |
| `13-security` | `studies/go-cli-study/reports/final/13-security.md` | Path, command, and secret boundaries must be explicit; go-task `internal/execext/exec.go:59-66`, opencode `internal/permission/permission.go:44-108`, restic `internal/options/secret_string.go:15-20`, and helm `pkg/registry/transport.go:37-41` show sandboxing and redaction approaches. | high |
| `14-performance` | `studies/go-cli-study/reports/final/14-performance.md` | Large local processing stays safe through lazy init, streaming, bounded search, and cache reuse; gh-cli `pkg/cmdutil/factory.go:27-42`, yq `pkg/yqlib/stream_evaluator.go:78-113`, gdu `pkg/analyze/parallel.go:13`, and rclone `lib/pool/pool.go:17-24` show the relevant pressures. | high |
| `15-philosophy` | `studies/go-cli-study/reports/final/15-philosophy.md` | Coherent projects accept complexity deliberately and keep scope boundaries visible; examples include gh-cli factory discipline `pkg/cmd/factory/default.go:26-46`, restic backend interfaces `internal/backend/backend.go:19-90`, and opencode pub/sub `internal/pubsub/broker.go:10-19`. | medium |

## Relevant Patterns

- **Thin CLI, product-owned modules:** The selected evidence repeatedly favors command layers that parse flags, construct requests, call a product module, render output, and map errors. This matters because this sprint needs `ultraplan study <study> summary` to stay in `internal/study` and `ultraplan code <report>...` behavior to stay in `internal/codeextract`, with `internal/app` acting as a thin adapter. Evidence includes the general thin-entrypoint pattern in `studies/go-cli-study/reports/final/01-project-structure.md` citing `cmd/age/age.go:105`, `cmd/gh/main.go:6`, and `cmd/root.go:112-127`, plus `studies/go-cli-study/reports/final/02-command-architecture.md` citing helm `pkg/cmd/install.go:132-145` and opencode's cautionary `cmd/root.go:49-183` monolithic `RunE`.

- **Factory and request objects for command seams:** Command factories and explicit option/request structs keep dependencies visible without a DI framework. This is relevant for summary and extraction because command tests need to inject workspace paths, output writers, and filesystem-failure conditions without launching runtime behavior. Evidence includes gh-cli `pkg/cmdutil/factory.go:16-43`, gh-cli `pkg/cmd/factory/default.go:26-46`, restic `cmd/restic/cmd_backup.go:84-115`, and go-task `executor.go:22-24` from `studies/go-cli-study/reports/final/02-command-architecture.md` and `studies/go-cli-study/reports/final/03-dependency-injection.md`.

- **Injectable IO and deterministic renderers:** Commands that render human text, JSON, or CSV should be testable with `bytes.Buffer` and in-memory output rather than direct `fmt.Println`/`os.Stdout`. This applies directly to deterministic human output, `summary.csv`, `--json`, and `--output`. Evidence includes gh-cli `pkg/iostreams/iostreams.go:551-568`, go-task `executor.go:553-564`, restic `internal/ui/mock.go:10-53`, and gdu `stdout/stdout_test.go:29-39` in `studies/go-cli-study/reports/final/06-io-abstraction.md`.

- **Error classification with user-action context:** Summary warnings, unresolved extraction references, missing source tables, malformed citations, and write failures should remain inspectable outcomes rather than undifferentiated strings. Evidence includes go-task's `TaskError` interface at `errors/errors.go:47-50`, gh-cli exit-code constants at `internal/ghcmd/cmd.go:44-49`, age's `errorWithHint()` at `cmd/age/tui.go:47-54`, and rclone's behavioral `Retrier` interface at `fs/fserrors/error.go:22-29` in `studies/go-cli-study/reports/final/05-error-handling.md`.

- **Fixture and golden testing for output contracts:** This sprint's CSV, text snippets, JSON shape, warnings, and command help are user-facing artifacts where exact-output or golden-style tests are appropriate. Evidence includes chezmoi's txtar testscript runner `internal/cmd/main_test.go:64-174`, helm `internal/test/test.go:43`, go-task `task_test.go:166-169`, and rclone `cmd/bisync/bisync_test.go:1435-1479` in `studies/go-cli-study/reports/final/11-testing-strategy.md`.

- **Bounded source lookup and local-only resolution:** Code extraction should resolve local paths under parsed source roots and avoid unbounded scans. The closest evidence comes from security and performance studies: explicit argument/path discipline in lazygit `cmd_obj_builder.go:38`, safe command lookup in gh-cli `pkg/surveyext/editor_manual.go:23`, path/trust boundary warnings in `studies/go-cli-study/reports/final/13-security.md`, and bounded concurrency/search pressure in gdu `pkg/analyze/parallel.go:13` plus streaming/bounded approaches in yq `pkg/yqlib/stream_evaluator.go:78-113` from `studies/go-cli-study/reports/final/14-performance.md`.

## Trade-Offs

| Trade-Off | Benefit | Cost | When It Matters |
| --- | --- | --- | --- |
| Thin CLI adapter vs command-owned product logic | Keeps summary and extraction logic testable without CLI setup; mirrors evidence from helm `pkg/cmd/install.go:132-145` and gh-cli `pkg/cmdutil/factory.go:16-43`. | Requires request/result types and rendering boundary decisions. | Every `internal/app` command path for `study summary` and `code`. |
| Exact/golden output tests vs flexible substring assertions | Catches CSV, text, JSON, and help regressions; supported by helm `internal/test/test.go:43` and go-task `task_test.go:166-169`. | Can make intentional formatting changes noisier and requires disciplined updates. | Deterministic `summary.csv`, extraction text output, unresolved summaries, JSON output, and help text. |
| Basename fallback search vs strict source-relative resolution | Recovers citations that omit source prefixes or use legacy report styles. | Risks ambiguity, expensive scans, and accidental traversal unless bounded and ignored dirs are explicit. | Code extraction resolution after direct source-relative and source-prefixed attempts fail. |
| Structured JSON output vs sprint-local schema only | Enables deterministic automation and tests for `ultraplan code --json`. | Broader schema stability can constrain future release changes if over-promised. | This sprint requires deterministic tested JSON but explicitly does not decide release-wide stable schema. |
| Warnings as non-fatal results vs hard failures | Allows one report's unresolved references or one summary cell warning to be reviewed without hiding other useful results. | Requires clear exit-code mapping so partial validation outcomes are visible. | Missing ratings, ambiguous ratings, unresolved code refs, and multiple report inputs. |
| Filesystem abstraction vs concrete local filesystem helpers | Improves fixture tests and simulated write failures; supported by restic `internal/fs/interface.go:10-31` and chezmoi `internal/chezmoi/system.go:25`. | Too broad an interface can add unnecessary machinery. | Atomic summary writes, optional extraction output writes, source file reads, and out-of-range line extraction tests. |

## Anti-Patterns And Warnings

- **Do not let `RunE` become the product implementation:** `studies/go-cli-study/reports/final/02-command-architecture.md` flags opencode `cmd/root.go:49-183` and yq long `RunE` functions as caution signs. For this sprint, parsing citation syntax, resolving paths, calculating summary totals, and writing CSV should not live in `internal/app`.

- **Do not duplicate parsers that already exist in the owning module:** The sprint scope requires reusing study discovery, report path, applicability, and rating parsing behavior. The project-structure evidence favors logic staying near state it transforms, and `studies/go-cli-study/reports/final/01-project-structure.md` warns about global technical-layer packages and bidirectional imports with examples like `internal/cmd`/`internal/chezmoi` import direction.

- **Do not silently downgrade write failures or unresolved references:** `studies/go-cli-study/reports/final/05-error-handling.md` identifies silent failure logging as an anti-pattern using lazygit `pkg/commands/git_commands/file_loader.go:52`. Summary write failure and extraction output write failure need visible non-zero outcomes, not only logged warnings.

- **Do not hardcode output to `os.Stdout`/`os.Stderr`:** `studies/go-cli-study/reports/final/06-io-abstraction.md` calls out rclone `cmd/ls/ls.go:42`, dive `dive/image/docker/cli.go:27`, and opencode direct output as testability gaps. This sprint needs in-memory command tests for text, JSON, warnings, and help.

- **Do not treat path strings as trusted:** `studies/go-cli-study/reports/final/13-security.md` warns against path and shell trust gaps, including gdu raw ignore path matching at `internal/common/ignore.go:16-37` and unrestricted shell-style pass-through in dive `cmd/dive/cli/internal/command/build.go:25`. Extraction should normalize, contain, and report path escape attempts.

- **Do not allow unbounded fallback scans:** `studies/go-cli-study/reports/final/14-performance.md` warns against unbounded accumulation and unbounded goroutine/search patterns, with gdu's bounded `pkg/analyze/parallel.go:13` as a contrast. Basename lookup should be scoped to parsed source roots and ignore `.git`, `node_modules`, and documented ignored dirs.

- **Do not test only internal state:** `studies/go-cli-study/reports/final/11-testing-strategy.md` flags k9s `internal/view/pod_test.go:23` hint-count assertions as brittle. Sprint tests should prefer observable command output, files, exit categories, parsed results, and unresolved diagnostics.

## Examples Worth Inspecting

| Example | Path / Source | Why It Is Useful |
| --- | --- | --- |
| gh-cli factory and IOStreams | `pkg/cmdutil/factory.go:16-43`, `pkg/iostreams/iostreams.go:551-568` | Good model for command dependencies and buffer-backed command output tests. |
| helm thin command/action boundary | `pkg/cmd/install.go:132-145`, `pkg/action/install.go:73-140` | Shows command-level construction delegating to product logic rather than embedding behavior. |
| go-task functional options | `executor.go:22-24`, `executor.go:553-564` | Useful for optional test seams around IO or filesystem behavior without broad interfaces. |
| restic terminal and filesystem mocks | `internal/ui/mock.go:10-53`, `internal/fs/interface.go:10-31` | Useful when reasoning about minimal abstractions for output and file-backed tests. |
| chezmoi txtar/testscript infrastructure | `internal/cmd/main_test.go:64-174` | Good reference for command-level workflows with generated files and exact outputs. |
| rclone error behavior interfaces | `fs/fserrors/error.go:22-29`, `cmd/cmd.go:497-516` | Useful for partial/validation outcomes and exit-code reasoning. |
| age hint-based CLI errors | `cmd/age/tui.go:37-54` | Useful for actionable user diagnostics without exposing unsafe internals. |
| yq streaming evaluator | `pkg/yqlib/stream_evaluator.go:78-113` | Useful for bounded processing mindset when reading reports and extracting lines. |
| opencode permission boundary | `internal/permission/permission.go:44-108`, `internal/llm/tools/bash.go:41-55` | Useful as evidence for explicit trust-boundary design, even though this sprint must not run runtime operations. |
| gdu bounded analysis | `pkg/analyze/parallel.go:13` | Useful as a cautionary model for avoiding unbounded filesystem traversal in fallback search. |

## Design Pressures

- Summary generation must reuse existing study-owned helpers while preserving deterministic CSV semantics, warning collection, inapplicable `N/A` cells, total-descending ordering, and atomic writes.

- Code extraction must parse enough citation syntax to audit existing reports while avoiding becoming a general report-processing or global parser package.

- Human output must be useful for review, but JSON/text surfaces also need deterministic ordering for tests and scripts.

- Unresolved extraction references are validation-style outcomes, not ordinary success and not fatal to processing other report arguments.

- Source table parsing is a trust boundary because report Markdown is generated content and may contain malformed paths, missing tables, duplicated source names, or paths outside intended roots.

- Basename fallback is useful but dangerous: it can recover weak citations, but it creates ambiguity, scan cost, and path containment risks.

- Summary and extraction writes are local filesystem mutations in an otherwise non-runtime sprint, so write paths, permissions, atomicity, and failure visibility must be reasoned explicitly.

- The sprint must preserve runtime/product separation: no agentwrap/OpenCode calls, no prompt composition changes, and no import pressure from `internal/platform/runtime` toward product modules.

- Tests need to exercise missing ratings, ambiguous ratings, missing reports, inapplicable Markdown pairs, source table parsing, malformed citations, path escape attempts, ambiguous basename matches, output files, and repeated-run determinism.

## Open Questions For Reasoning

- What exact result and warning model lets summary generation distinguish missing report, missing rating, ambiguous rating, and inapplicable pair without duplicating validation internals?

- Should `summary.csv` use only `N/A` for inapplicable pairs while all other invalid/missing rating conditions remain empty cells, and how should human output summarize those conditions?

- What is the smallest `internal/codeextract` API that supports CLI text/JSON rendering without leaking CLI-specific formatting into the module?

- How should the extractor decide that a report requires a source table, especially if a report has no citations or contains only non-code document citations?

- How should source root resolution handle workspace-relative versus report-relative paths when both could exist?

- What exact ignored-directory set should basename fallback use beyond `.git` and `node_modules`, and should it be fixed or locally configurable for this sprint?

- What exit-code distinction should represent unresolved references across multiple reports: validation failure, partial completion, or another existing app outcome?

- How much JSON structure is necessary for deterministic tests without implying release-wide stable schema compatibility?

- Which write helper, if any, should be reused for atomic `summary.csv` writes and optional extraction output writes without creating a broad filesystem abstraction prematurely?

- How should output redaction avoid prompt bodies and sensitive environment details while still showing report paths, source names, resolved paths, and snippets requested by the user?

## Evidence Pointers

- `studies/go-cli-study/reports/final/01-project-structure.md`: Inspect Pattern 1, Pattern 2, and Anti-Patterns for thin CLI, `internal/` boundaries, and global-layer warnings; key citations include `cmd/gh/main.go:6`, `cmd/root.go:112-127`, `internal/chezmoi/chezmoi.go:1-2`, and `internal/cmd/config.go:1`.

- `studies/go-cli-study/reports/final/02-command-architecture.md`: Inspect Thin-Delegate Pattern and Pattern Catalog for command factories and command grouping; key citations include `pkg/cmdutil/factory.go:16-43`, `pkg/cmd/install.go:132-145`, `cmd/cmd.go:240-340`, and `cmd/root.go:49-183` as a caution.

- `studies/go-cli-study/reports/final/03-dependency-injection.md`: Inspect constructor injection, functional options, and global-state anti-patterns; key citations include `pkg/cmd/factory/default.go:26-46`, `executor.go:22-24`, `internal/restic/repository.go:18-66`, `fs/cache/cache.go:16-21`, and `pkg/yqlib/lib.go:13`.

- `studies/go-cli-study/reports/final/05-error-handling.md`: Inspect typed errors, sentinels, user/operational separation, and exit-code mapping; key citations include `errors/errors.go:47-50`, `internal/ghcmd/cmd.go:44-49`, `cmd/age/tui.go:47-54`, `fs/fserrors/error.go:22-29`, and `pkg/commands/git_commands/file_loader.go:52`.

- `studies/go-cli-study/reports/final/06-io-abstraction.md`: Inspect IOStreams, functional IO options, filesystem interfaces, and hardcoded IO warnings; key citations include `pkg/iostreams/iostreams.go:551-568`, `executor.go:553-564`, `internal/ui/mock.go:10-53`, `internal/fs/interface.go:10-31`, and `cmd/ls/ls.go:42`.

- `studies/go-cli-study/reports/final/09-terminal-ux.md`: Inspect non-TTY fallback, CLI-first output, cancellation, and progress guidance; key citations include `internal/cmd/prompt.go:20-256`, `internal/prompter/`, `cmd/progress.go:24-71`, and `pkg/yqlib/color_print.go:7-9`.

- `studies/go-cli-study/reports/final/10-logging-observability.md`: Inspect stdout/stderr separation, structured keys, and debug control; key citations include `pkg/iostreams/iostreams.go:52-54`, `internal/logging/logging.go:31-66`, `internal/slogs/keys.go:6-231`, and `src/core.go:325` as a stdout-mixing warning.

- `studies/go-cli-study/reports/final/11-testing-strategy.md`: Inspect testscript integration, golden files, centralized mocks, and behavior-focused assertions; key citations include `internal/cmd/main_test.go:64-174`, `task_test.go:166-169`, `internal/test/test.go:43`, `acceptance/acceptance_test.go:26-29`, and `internal/view/pod_test.go:23` as a caution.

- `studies/go-cli-study/reports/final/13-security.md`: Inspect shell-safe argument construction, secret redaction, path validation cautions, and trust-boundary visibility; key citations include `cmd_obj_builder.go:38`, `internal/options/secret_string.go:15-20`, `pkg/registry/transport.go:37-41`, `internal/permission/permission.go:44-108`, and `internal/common/ignore.go:16-37`.

- `studies/go-cli-study/reports/final/14-performance.md`: Inspect lazy initialization, streaming, bounded concurrency, and unbounded accumulation warnings; key citations include `pkg/cmdutil/factory.go:27-42`, `pkg/yqlib/stream_evaluator.go:78-113`, `pkg/analyze/parallel.go:13`, `lib/pool/pool.go:17-24`, and `internal/model/table.go:200` as a memory warning.

- `studies/go-cli-study/reports/final/15-philosophy.md`: Inspect evidence for deliberate scope, factory/options patterns, interface-driven abstraction, and complexity cautions; key citations include `pkg/cmd/factory/default.go:26-46`, `internal/backend/backend.go:19-90`, `internal/pubsub/broker.go:10-19`, `VISION.md:97`, and `pkg/analyze/parallel.go:13`.

## Handoff To Reasoning

- Use this handbook as evidence input.
- Validate whether the observed patterns fit this project's constraints.
- Do not copy external patterns without sprint-specific reasoning.
- Defer final choices about output shapes, warning types, exit-code mapping, ignored-directory sets, atomic write helpers, and JSON fields to the area reasoning and consolidated reasoning artifacts.
