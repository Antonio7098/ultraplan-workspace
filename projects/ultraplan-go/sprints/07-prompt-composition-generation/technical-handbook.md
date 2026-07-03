# Sprint Technical Handbook: 07 Prompt Composition Generation

> Project: `ultraplan-go`
> Sprint: `07-prompt-composition-generation`
> Source: `sprint-index.md`
> **Inputs Used:** `projects/ultraplan-go/project-index.md`, `projects/ultraplan-go/sprints/07-prompt-composition-generation/sprint-index.md`, `templates/technical-handbook.md`, `projects/ultraplan-go/sprints/07-prompt-composition-generation/requirements.md`, `projects/ultraplan-go/docs/ARCHITECTURE.md`, `projects/ultraplan-go/docs/PRD.md`, `projects/ultraplan-go/docs/TRD.md`, `studies/go-cli-study/reports/final/01-project-structure.md`, `studies/go-cli-study/reports/final/02-command-architecture.md`, `studies/go-cli-study/reports/final/03-dependency-injection.md`, `studies/go-cli-study/reports/final/05-error-handling.md`, `studies/go-cli-study/reports/final/06-io-abstraction.md`, `studies/go-cli-study/reports/final/11-testing-strategy.md`, `studies/go-cli-study/reports/final/13-security.md`, `studies/go-cli-study/reports/final/14-performance.md`, `studies/go-cli-study/reports/final/15-philosophy.md`

This handbook distills the studies and reports selected by `sprint-index.md` for sprint reasoning. It does not decide architecture or implementation.

## Selected Studies And Reports

| Study / Report | Path | Relevant Finding | Confidence |
| --- | --- | --- | --- |
| `01-project-structure` | `studies/go-cli-study/reports/final/01-project-structure.md` | Mature Go CLIs keep CLI entrypoints thin and push business logic into protected or domain-owned packages; evidence includes `cmd/age/age.go:105`, `cmd/gh/main.go:6`, `cmd/root.go:112-127`, and unidirectional imports such as `cmd/root.go:9` and `cmd/evaluate_sequence_command.go:7`. | high |
| `02-command-architecture` | `studies/go-cli-study/reports/final/02-command-architecture.md` | Command factories and thin `RunE` handlers improve testability and prevent command-layer business logic; examples include `pkg/cmd/install.go:132-145`, `pkg/cmdutil/factory.go:16-43`, `pkg/cmd/issue/list/list.go:47-118`, and `cmd/root.go:49-183` as a monolithic caution. | high |
| `03-dependency-injection` | `studies/go-cli-study/reports/final/03-dependency-injection.md` | Manual composition roots, constructor injection, functional options, and minimal globals dominate high-scoring CLIs; evidence includes `internal/ghcmd/cmd.go:52-132`, `pkg/cmd/factory/default.go:26-46`, `executor.go:22-24`, and global-state cautions such as `fs/cache/cache.go:16-21`. | high |
| `05-error-handling` | `studies/go-cli-study/reports/final/05-error-handling.md` | Mature CLIs wrap errors with `%w`, use sentinels or typed errors for actionable branching, and separate user-facing diagnostics from operational details; examples include `pkg/storage/driver/driver.go:27-36`, `go-task/errors/errors_task.go:13-32`, and `gh-cli/internal/ghcmd/cmd.go:281-301`. | high |
| `06-io-abstraction` | `studies/go-cli-study/reports/final/06-io-abstraction.md` | Injectable `io.Reader`/`io.Writer`, filesystem abstractions, and test constructors make command output and file behavior testable; evidence includes `pkg/iostreams/iostreams.go:551-568`, `executor.go:553-564`, `internal/chezmoi/system.go:25`, and `internal/fs/interface.go:10-31`. | high |
| `11-testing-strategy` | `studies/go-cli-study/reports/final/11-testing-strategy.md` | Table-driven tests, CLI integration tests, golden output checks, and fake command runners are the common reliable-test patterns; examples include `chezmoi/internal/cmd/main_test.go:64-174`, `helm/internal/test/test.go:43`, `go-task/task_test.go:166-169`, and `lazygit/pkg/commands/oscommands/fake_cmd_obj_runner.go:17-26`. | high |
| `13-security` | `studies/go-cli-study/reports/final/13-security.md` | Explicit trust boundaries, safe subprocess construction, secret redaction, schema/input validation, and permission flows matter most when tools handle untrusted inputs; evidence includes `internal/llm/tools/bash.go:41-55`, `pkg/registry/transport.go:37-41`, and `internal/config/json/validator.go:146`. | medium |
| `14-performance` | `studies/go-cli-study/reports/final/14-performance.md` | Lazy initialization, streaming, bounded memory, and avoiding unnecessary startup work are the dominant performance lessons; evidence includes `gh-cli/pkg/cmdutil/factory.go:27-42`, `age/internal/stream/stream.go:20`, and `yq/pkg/yqlib/stream_evaluator.go:78-113`. | medium |
| `15-philosophy` | `studies/go-cli-study/reports/final/15-philosophy.md` | Coherent tools accept complexity deliberately, document non-goals, and use factory/interface/decorator patterns where they serve the product; evidence includes `pkg/cmd/factory/default.go:26-46`, `internal/chezmoi/system.go:25-45`, and `VISION.md:97`. | medium |

## Relevant Patterns

- **Study-owned prompt logic behind thin CLI wiring:** The selected structure and command reports repeatedly favor thin entrypoints that delegate to domain logic rather than embedding behavior in command handlers. `01-project-structure` cites thin entrypoints such as `cmd/age/age.go:105`, `cmd/gh/main.go:6`, and `cmd/root.go:112-127`; `02-command-architecture` cites thin command constructors like `pkg/cmd/install.go:132-145` and cautions against `cmd/root.go:49-183` as a 130-line `RunE`. For this sprint, this supports reasoning about CLI preview commands as routing/output code while prompt assembly, manifests, applicability, and template validation remain in study-owned code.
- **Factory-style command construction with explicit dependencies:** Command architecture and DI reports show command factories receiving shared config/dependencies as a stable pattern. Evidence includes `pkg/cmdutil/factory.go:16-43`, `pkg/cmd/issue/list/list.go:47-118`, `cmd/restic/cmd_backup.go:35`, and `pkg/cmd/factory/default.go:26-46`. For prompt preview, reasoning should evaluate how to wire workspace resolution, study services, and IO streams without adding global prompt/template registries.
- **Lazy, deterministic input loading:** Performance and DI evidence favors deferring expensive or unnecessary work until commands need it, using `func()` fields, `sync.Once`, or narrow constructors. Evidence includes `gh-cli/pkg/cmdutil/factory.go:27-42`, `restic/internal/repository/repository.go:52-53`, `helm/pkg/action/lazyclient.go:35-53`, and `age/cmd/age/age.go:476-484`. Prompt preview should reason about deterministic reads of required prompt/template/report inputs without triggering runtime, network, provider, or repository-wide scan paths.
- **Injectable IO for preview output and tests:** IO abstraction evidence highlights test constructors and writer injection as the simplest way to test CLI output. `pkg/iostreams/iostreams.go:551-568` returns test streams, `executor.go:553-564` injects stdout, `ui.go:19-43` abstracts UI, and `stdout/stdout_test.go:29-39` captures stdout behavior. Prompt preview should keep stdout/file-writing behavior capturable and command-testable.
- **Structured, actionable error paths:** Error-handling reports favor `%w` wrapping, sentinel/typed errors, and hinting for user-facing diagnostics. Evidence includes `%w` examples at `pkg/ssh/ssh_keys.go:64`, sentinel errors at `pkg/storage/driver/driver.go:27-36`, task suggestions at `go-task/errors/errors_task.go:13-32`, and contextual rendering at `gh-cli/internal/ghcmd/cmd.go:281-301`. Prompt composition reasoning should preserve causes for filesystem/template failures and identify the missing template, source, dimension, report, or output path.
- **Golden and fixture tests for deterministic text artifacts:** Testing reports show output regression protection through golden/snapshot fixtures and script-level CLI coverage. Evidence includes `helm/internal/test/test.go:43`, `go-task/task_test.go:166-169`, `dive/cmd/dive/cli/testdata/snapshots/cli_build_test.snap:1`, and `chezmoi/internal/cmd/main_test.go:64-174`. This sprint's deterministic prompt text and manifests are good candidates for table-driven fixture or golden-style tests.
- **Explicit trust boundaries for no-runtime preview:** Security evidence emphasizes command filtering, permission boundaries, and safe diagnostics in tools that may interact with untrusted inputs or subprocesses. Evidence includes opencode's allowlist/banned command checks at `internal/llm/tools/bash.go:41-55`, permission flow at `internal/permission/permission.go:44-108`, and credential scrubbing at `pkg/registry/transport.go:37-41`. For this sprint, preview must remain a prompt-rendering operation, not a runtime or subprocess operation.

## Trade-Offs

| Trade-Off | Benefit | Cost | When It Matters |
| --- | --- | --- | --- |
| Keep prompt composition inside `internal/study` instead of creating a global prompt/template package | Matches domain-owned boundary evidence: business logic lives behind CLI entrypoints and import direction stays one-way (`cmd/root.go:112-127`, `internal/chezmoi/chezmoi.go:1-2`, `cmd/evaluate_sequence_command.go:7`) | `internal/study` may grow another focused file and must manage study-specific template/report concepts locally | When deciding package placement for prompt builders, manifests, and template loading |
| Thin preview command plus study builder vs self-contained command handler | Easier unit tests and clearer separation; aligns with factory/thin command evidence (`pkg/cmd/install.go:132-145`, `pkg/cmdutil/factory.go:16-43`, `pkg/cmd/issue/list/list.go:47-118`) | Requires passing request/response structs between app and study packages and may feel like indirection for one command | When shaping `internal/app/study_prompt_commands.go` and deciding what command tests assert directly |
| Golden/fixture prompt assertions vs flexible substring assertions | Strong regression protection for deterministic prompt text, similar to golden practices (`helm/internal/test/test.go:43`, `go-task/task_test.go:166-169`, `rclone/cmd/bisync/bisync_test.go:1435-1479`) | Golden diffs can be noisy and require deliberate updates when wording changes | When testing full prompt text, manifests, and stable ordering |
| Workspace-relative manifest paths vs absolute diagnostic paths | Generated artifacts stay relocatable and reviewable, consistent with artifact discipline and path guidance; security evidence favors avoiding leaked local details where not needed (`pkg/registry/transport.go:37-41`, `internal/options/secret_string.go:15-20`) | Some local failures may need absolute paths or resolved paths for immediate debugging | When defining manifest fields and CLI error text |
| Fail-fast missing input checks vs partial preview output | Prevents producing runnable-looking prompts with missing templates/reports and aligns with actionable error patterns (`pkg/storage/driver/driver.go:27-36`, `gh-cli/internal/ghcmd/cmd.go:281-301`) | Users may not see any partial prompt when one input is missing | When resolving templates, Markdown documents, report templates, and synthesis input manifests |
| Strict no-runtime preview vs sharing future runtime command paths | Security and test isolation are simpler; avoids subprocess, provider, OpenCode, and network pathways, matching explicit trust-boundary evidence (`internal/llm/tools/bash.go:41-55`, `internal/permission/permission.go:44-108`) | Later runtime execution may need parallel wiring that reuses prompt builders without reusing preview command code | When deciding whether preview invokes run/synthesize placeholders or calls prompt builders directly |

## Anti-Patterns And Warnings

- **Do not let command handlers become prompt composition engines:** `02-command-architecture` flags long `RunE` functions like `opencode cmd/root.go:49-183`, `yq cmd/evaluate_sequence_command.go:152`, and `restic cmd/restic/cmd_backup.go:84-115` as maintainability risks. Preview command code should not build large prompt strings inline.
- **Do not introduce hidden globals for templates, config, or IO:** DI evidence warns about global config and singleton cache pollution, citing `pkg/yqlib/lib.go:13`, `fs/cache/cache.go:16-21`, and `fs/config.go:793`. Prompt template paths, workspace root, and output streams should be explicit inputs, not package-level state.
- **Do not bypass injectable output streams:** IO evidence flags direct `os.Stdout`/`os.Stderr` access such as `cmd/ls/ls.go:42`, `cmd/age/tui.go:31`, and `pkg/yqlib/logger.go:18`. Preview output must be capturable in command tests.
- **Do not break error chains or rely on string matching:** Error evidence warns against `fmt.Errorf` without `%w` at `mitchellh-cli/cli.go:205-206` and stringly sentinels at `k9s/internal/client/errors.go:9-14`. Missing templates, inapplicable pairs, and absent reports should preserve causes and support structured handling where useful.
- **Do not treat Markdown document sources like repositories:** Security and prompt-scope pressures require trust boundaries. The sprint scope requires embedded Markdown-only analysis, and security evidence warns against unrestricted or implicit shell/filesystem trust boundaries (`cmd/dive/cli/internal/command/build.go:25`, `pkg/yqlib/security_prefs.go:3-7`). Markdown prompts should not invite runtime exploration of external files or code.
- **Do not add speculative performance machinery:** Performance evidence says to profile before pooling and reserve buffer pools for hot paths (`fzf/src/matcher.go:183-185`, `rclone/lib/pool/pool.go:17-24`, `restic/internal/archiver/buffer.go:24-46`). Prompt rendering is deterministic local IO; avoid adding concurrency, pools, or caching unless reasoning identifies a real need.
- **Do not assert brittle implementation details in tests:** Testing evidence warns about hint-count assertions at `k9s/internal/view/pod_test.go:23` and fragile regex output at `pkg/cmd/issue/list/list_test.go:88-91`. Tests should assert externally relevant prompt text, manifests, error categories/messages, and no-runtime behavior.

## Examples Worth Inspecting

| Example | Path / Source | Why It Is Useful |
| --- | --- | --- |
| Thin command delegation | `pkg/cmd/install.go:132-145`, `pkg/action/install.go:73-140` | Demonstrates command construction delegating real work to action/domain logic. Useful when reasoning about preview command vs study prompt builder boundaries. |
| Factory DI | `pkg/cmdutil/factory.go:16-43`, `pkg/cmd/factory/default.go:26-46` | Shows shared dependencies as explicit fields/functions and supports lazy construction for command tests. |
| IO test constructor | `pkg/iostreams/iostreams.go:551-568` | Useful model for making preview stdout/stderr behavior easy to capture in tests. |
| System/filesystem abstraction | `internal/chezmoi/system.go:25`, `internal/fs/interface.go:10-31` | Shows how filesystem boundaries can be abstracted when testability requires it, while reasoning can decide whether existing project filesystem helpers are enough. |
| Actionable CLI errors | `gh-cli/internal/ghcmd/cmd.go:281-301`, `go-task/errors/errors_task.go:13-32` | Useful for missing source/dimension/template diagnostics and typo/ambiguity handling. |
| Golden output helpers | `helm/internal/test/test.go:43`, `go-task/task_test.go:166-169` | Useful for deterministic prompt and manifest regression tests. |
| Script-level CLI tests | `chezmoi/internal/cmd/main_test.go:64-174`, `gh-cli/acceptance/acceptance_test.go:26-29` | Useful if preview command behavior benefits from end-to-end CLI tests rather than only direct command invocation. |
| Permission and no-runtime boundary | `internal/permission/permission.go:44-108`, `internal/llm/tools/bash.go:41-55` | Shows how agent-facing tools make trust boundaries explicit; sprint reasoning should adapt the caution by keeping preview runtime-free. |
| Lazy init | `gh-cli/pkg/cmdutil/factory.go:27-42`, `helm/pkg/action/lazyclient.go:35-53` | Useful for avoiding unnecessary runtime/provider/config work during prompt preview. |

## Design Pressures

- The project architecture says product behavior belongs with the module that owns its state, while the selected reports show this same pressure through thin CLI and unidirectional dependency patterns.
- Prompt composition must be deterministic text generation, not runtime execution; this creates pressure to isolate prompt builders from future agentwrap/OpenCode paths.
- Directory and Markdown document sources have different trust and citation rules, so a single analysis prompt path must not blur source-kind semantics.
- Synthesis depends on deterministic applicable report manifests, so ordering, missing-report behavior, and inapplicable Markdown pairs need explicit reasoning.
- Preview output must support both stdout and explicit output path behavior, making IO injection and file-write testability important.
- Missing templates, missing Markdown documents, ambiguous references, unknown references, inapplicable pairs, and missing synthesis inputs all need actionable diagnostics without losing filesystem causes.
- The manifest is both a user-facing preview artifact and a test oracle, creating pressure to keep fields stable, workspace-relative where possible, and free of secrets or accidental absolute paths.
- The sprint depends on prior Markdown frontmatter stripping and applicability behavior, so prompt composition should consume those capabilities rather than redefining source discovery or applicability rules.
- Tests must prove absence of runtime invocation, so command wiring should avoid indirect calls into run/synthesize paths that may later grow runtime side effects.

## Open Questions For Reasoning

- What exact prompt request/result types belong in `internal/study/domain.go`, and which fields are necessary for deterministic manifests without over-modeling future runtime execution?
- Should directory, Markdown document, and synthesis prompts use one builder entrypoint with prompt kind dispatch, or separate focused builder methods that share template loading?
- How should an inapplicable Markdown source/dimension pair be represented: a typed non-runnable result, a skipped result, or an error returned only for single-source preview?
- How should dimensions express that directory code citation requirements are disabled, and how should Markdown prompts state that code citations do not apply unless the dimension says otherwise?
- Which paths in manifests should be workspace-relative, study-relative, or absolute-only-for-diagnostics?
- What deterministic order should synthesis manifests use when combining applicable directory and Markdown source reports?
- Should prompt preview write only prompt text, or also expose the manifest through adjacent output, stdout framing, or a separate flag?
- What is the minimal CLI shape for preview that satisfies scriptability without committing to future runtime command names or `run --dry-run` semantics?
- How should tests prove no agentwrap/OpenCode/provider/network/subprocess invocation beyond checking command wiring and absence of runtime dependencies?

## Evidence Pointers

- `studies/go-cli-study/reports/final/01-project-structure.md`: inspect Pattern 1, Pattern 4, Tradeoffs, and Evidence Index for thin CLI, `internal/`, and unidirectional dependency evidence.
- `studies/go-cli-study/reports/final/02-command-architecture.md`: inspect Pattern 1, Pattern 2, Pattern 3, Anti-Patterns, and Evidence Index for command factory and long handler cautions.
- `studies/go-cli-study/reports/final/03-dependency-injection.md`: inspect Pattern 1, Pattern 2, Pattern 7, Global State Usage, and Anti-Patterns for explicit wiring and lazy dependency construction.
- `studies/go-cli-study/reports/final/05-error-handling.md`: inspect Patterns 1, 2, 3, 7, 10 and Anti-Patterns for wrapped errors, sentinels, hints, and exit-code mapping.
- `studies/go-cli-study/reports/final/06-io-abstraction.md`: inspect IOStreams, functional stream injection, System Interface, and Anti-Patterns for preview stdout/file-output testability.
- `studies/go-cli-study/reports/final/11-testing-strategy.md`: inspect testscript, golden-file, fake runner, fixture extraction, and anti-pattern sections for prompt/manifest test strategy.
- `studies/go-cli-study/reports/final/13-security.md`: inspect Shell-Safe Argument Construction, Secret Redaction Type, Permission Request Flow, and anti-patterns for no-runtime and safe artifact boundaries.
- `studies/go-cli-study/reports/final/14-performance.md`: inspect Lazy Initialization and streaming sections for avoiding unnecessary runtime/provider/repo-wide work during preview.
- `studies/go-cli-study/reports/final/15-philosophy.md`: inspect Factory + Options, Interface-Driven Abstraction, and Anti-Patterns for deliberate complexity and explicit non-goals.

## Handoff To Reasoning

- Use this handbook as evidence input.
- Validate whether the observed patterns fit this project's constraints.
- Do not copy external patterns without sprint-specific reasoning.
- Resolve package ownership, prompt kind boundaries, manifest shape, preview CLI surface, and test strategy in sprint reasoning before implementation.
