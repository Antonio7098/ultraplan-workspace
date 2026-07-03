# Sprint Technical Handbook: Reason Stage

> Project: `ultraplan-go`
> Sprint: `20-reason-stage`
> Source: `projects/ultraplan-go/sprints/20-reason-stage/sprint-index.md`
> **Inputs Used:** `projects/ultraplan-go/sprints/20-reason-stage/sprint-index.md`, `templates/technical-handbook.md`, `projects/ultraplan-go/project-index.md`, `projects/ultraplan-go/sprints/20-reason-stage/requirements.md`, `projects/ultraplan-go/docs/PRD.md`, `projects/ultraplan-go/docs/TRD.md`, `projects/ultraplan-go/docs/ARCHITECTURE.md`, `studies/go-cli-study/reports/final/01-project-structure.md`, `studies/go-cli-study/reports/final/02-command-architecture.md`, `studies/go-cli-study/reports/final/03-dependency-injection.md`, `studies/go-cli-study/reports/final/04-configuration-management.md`, `studies/go-cli-study/reports/final/05-error-handling.md`, `studies/go-cli-study/reports/final/06-io-abstraction.md`, `studies/go-cli-study/reports/final/07-state-context.md`, `studies/go-cli-study/reports/final/09-terminal-ux.md`, `studies/go-cli-study/reports/final/10-logging-observability.md`, `studies/go-cli-study/reports/final/11-testing-strategy.md`, `studies/go-cli-study/reports/final/13-security.md`, `studies/go-cli-study/reports/final/15-philosophy.md`

This handbook distills the studies and reports selected by `sprint-index.md` for sprint reasoning. It does not decide architecture or implementation.

Project docs were used only for sprint scope context. The evidence base below is the selected evidence reports and their cited source references.

## Selected Studies And Reports

| Study / Report | Path | Relevant Finding | Confidence |
| --- | --- | --- | --- |
| `01-project-structure` | `studies/go-cli-study/reports/final/01-project-structure.md` | Non-trivial Go CLIs keep entrypoints thin and route work inward: chezmoi `main.go:26-34`, gh-cli `cmd/gh/main.go:6`, k9s `cmd/root.go:112-127`, restic `cmd/restic/main.go:37-114`; `internal/` boundaries are used for application-owned behavior such as chezmoi `internal/chezmoi/chezmoi.go:1-2` and restic `internal/restic/repository.go:18`. | high |
| `02-command-architecture` | `studies/go-cli-study/reports/final/02-command-architecture.md` | Mature command trees use factory functions and delegate from command handlers: gh-cli `pkg/cmdutil/factory.go:16-43`, helm `pkg/cmd/install.go:132-145`, restic `cmd/restic/main.go:37-114`; large `RunE` handlers are a warning, such as opencode `cmd/root.go:49-183`. | high |
| `03-dependency-injection` | `studies/go-cli-study/reports/final/03-dependency-injection.md` | Manual composition roots, constructor injection, functional options, and small interfaces dominate over DI frameworks: gh-cli `pkg/cmd/factory/default.go:26-46`, go-task `executor.go:22-24`, helm `pkg/action/action.go:118-146`, opencode `internal/app/app.go:42-81`. | high |
| `04-configuration-management` | `studies/go-cli-study/reports/final/04-configuration-management.md` | Config quality depends on explicit precedence and post-merge validation: chezmoi restores changed flags after config load at `internal/cmd/config.go:2253-2287`, go-task centralizes generic precedence at `internal/flags/flags.go:314-327`, opencode validates centrally at `internal/config/config.go:609-641`. | high |
| `05-error-handling` | `studies/go-cli-study/reports/final/05-error-handling.md` | Strong CLIs wrap errors with `%w`, use sentinels or typed errors for classification, and separate user-facing messages from operational details: rclone `fs/fserrors/error.go:22-29`, go-task `errors/errors.go:47-50`, restic `internal/errors/fatal.go:10`, gh-cli `internal/ghcmd/cmd.go:44-49`. | high |
| `06-io-abstraction` | `studies/go-cli-study/reports/final/06-io-abstraction.md` | Testable CLIs inject stdout, stderr, filesystem, and runtime-like boundaries rather than hardcoding `os.*`: gh-cli `pkg/iostreams/iostreams.go:551-568`, go-task `executor.go:541-577`, restic `internal/ui/terminal.go:10-36`, chezmoi `internal/chezmoi/system.go:25`. | high |
| `07-state-context` | `studies/go-cli-study/reports/final/07-state-context.md` | Long or cancellable workflows create a root context and pass it through I/O and work layers; state is commonly centralized in an app/config/executor object: gh-cli `internal/ghcmd/cmd.go:142`, go-task `task.go:89`, helm `pkg/cmd/install.go:333-347`, opencode `internal/session/session.go:12-23`. | medium |
| `09-terminal-ux` | `studies/go-cli-study/reports/final/09-terminal-ux.md` | Scriptable CLI UX favors non-TTY fallbacks, clean cancellation, and restrained progress. Relevant evidence includes chezmoi non-TTY prompt handling at `internal/cmd/prompt.go:20-256`, gh-cli prompter checks in `internal/prompter/`, and go-task signal handling at `signals.go:11-31`. | medium |
| `10-logging-observability` | `studies/go-cli-study/reports/final/10-logging-observability.md` | Operationally useful CLIs keep stdout for user/data output and stderr or files for diagnostics, with structured logging where useful: helm `internal/logging/logging.go:31-66`, k9s `internal/slogs/keys.go:6-231`, gh-cli `pkg/iostreams/iostreams.go:52-54`, rclone `fs/log/slog.go:21`. | high |
| `11-testing-strategy` | `studies/go-cli-study/reports/final/11-testing-strategy.md` | High-confidence CLI changes use table tests, fake seams, command-level or testscript integration, golden output where stable, and fixture-driven validation: chezmoi `internal/cmd/main_test.go:64-174`, gh-cli `acceptance/acceptance_test.go:26-29`, helm `internal/test/test.go:43`, go-task `task_test.go:166-169`. | high |
| `13-security` | `studies/go-cli-study/reports/final/13-security.md` | Security-relevant CLI work marks trust boundaries explicitly, avoids shell/string command construction, redacts secrets, and validates untrusted config: opencode `internal/permission/permission.go:44-108`, opencode `internal/llm/tools/bash.go:41-55`, restic `internal/options/secret_string.go:15-20`, helm `pkg/registry/transport.go:37-41`, k9s `internal/config/json/validator.go:146`. | high |
| `15-philosophy` | `studies/go-cli-study/reports/final/15-philosophy.md` | Coherent projects accept complexity deliberately and reject scope that does not serve the core workflow. Relevant patterns include gh-cli factory wiring `pkg/cmd/factory/default.go:26-46`, restic backend interfaces `internal/backend/backend.go:19-90`, opencode pub/sub `internal/pubsub/broker.go:10-19`, and go-task experiments `website/src/docs/experiments/index.md:17-21`. | medium |

## Relevant Patterns

- **Thin CLI, sprint-owned behavior:** The CLI layer should remain command parsing, output rendering, exit-code mapping, and service delegation. This is supported by project-structure evidence from chezmoi `main.go:26-34`, gh-cli `cmd/gh/main.go:6`, and k9s `cmd/root.go:112-127`, plus command-architecture evidence from helm `pkg/cmd/install.go:132-145` and gh-cli `pkg/cmdutil/factory.go:16-43`. For this sprint, the pressure is to keep reason-stage rules near sprint planning state rather than in command handlers.

- **Module boundary via `internal/` and unidirectional imports:** Application projects use `internal/` to protect private behavior and keep import flow one-way. Evidence includes chezmoi `internal/chezmoi/chezmoi.go:1-2`, restic `internal/restic/repository.go:18`, yq `cmd/root.go:9`, and yq `cmd/evaluate_sequence_command.go:7`. For this sprint, the pattern argues for preserving product module ownership and avoiding reverse imports from platform code into sprint behavior.

- **Manual composition root with explicit seams:** Mature CLIs avoid reflection-based DI and instead wire dependencies through factories, constructors, functional options, and app/config structs. Evidence includes gh-cli `pkg/cmd/factory/default.go:26-46`, go-task `executor.go:22-24`, helm `pkg/action/action.go:118-146`, and opencode `internal/app/app.go:42-81`. This matters for fake runtime, fake filesystem, clock, and IO seams in reason-stage tests.

- **Injectable IO for command testability:** CLI output should flow through injected writers or UI abstractions so tests can verify stdout, stderr, and no-ANSI behavior. Evidence includes gh-cli `pkg/iostreams/iostreams.go:551-568`, go-task `executor.go:541-577`, mitchellh-cli `ui.go:19-43`, restic `internal/ui/terminal.go:10-36`, and gdu `stdout/stdout.go:45-58`.

- **Explicit validation gates after generation:** Error-handling and testing reports both favor behavior where successful execution is not enough; outputs and state must be checked. Evidence is indirect but consistent: typed/exit-code errors in go-task `errors/errors.go:47-50`, fatal classification in restic `internal/errors/fatal.go:10`, command lifecycle tests in chezmoi `internal/cmd/main_test.go:64-174`, and golden/fixture checks in helm `internal/test/test.go:43`. This pattern applies to runtime-backed reasoning generation only as a caution: sprint reasoning must decide exactly how validation gates and state updates compose.

- **Context propagation for long or runtime-backed work:** Long operations create an entry context and pass it through execution layers, with signal cancellation when appropriate. Evidence includes gh-cli `internal/ghcmd/cmd.go:142`, go-task `task.go:89`, helm `pkg/cmd/install.go:333-347`, and restic `cmd/restic/cleanup.go:24-38`. For this sprint, this supports treating runtime-backed flow differently from dry-run validate/prompt paths while keeping cancellation explicit.

- **Structured, safe diagnostics with stdout/stderr separation:** Logging evidence favors structured diagnostic paths and strict separation from user output. Evidence includes helm `internal/logging/logging.go:31-66`, gh-cli `pkg/iostreams/iostreams.go:52-54`, yq `pkg/yqlib/logger.go:5`, and rclone `fs/log/slog.go:21`. This pattern matters for validation findings, prompt previews, dry-run output, and runtime errors.

- **Security boundaries around commands, paths, and secrets:** Security evidence favors explicit trust boundaries, redaction types, and no stringly shell execution. Evidence includes opencode permission checks in `internal/permission/permission.go:44-108`, opencode command filtering in `internal/llm/tools/bash.go:41-55`, restic `SecretString` in `internal/options/secret_string.go:15-20`, helm auth-header scrubbing in `pkg/registry/transport.go:37-41`, and k9s config schema validation in `internal/config/json/validator.go:146`.

## Trade-Offs

| Trade-Off | Benefit | Cost | When It Matters |
| --- | --- | --- | --- |
| Keep reason-stage behavior inside `internal/sprint` vs extract generic planning/workflow packages | Preserves module-owned behavior, mirrors `internal/` application boundaries from chezmoi `internal/chezmoi/chezmoi.go:1-2` and restic `internal/restic/repository.go:18`, and reduces premature abstraction | Some duplicated mechanics may remain until repeated needs are proven | When deciding where selected-template parsing, prompt rendering, validation, and flow-state transitions live |
| Command factories and service methods vs larger command handlers | Improves testability and keeps CLI wiring thin, following gh-cli `pkg/cmdutil/factory.go:16-43`, helm `pkg/cmd/install.go:132-145`, and restic `cmd/restic/main.go:37-114` | More small request/response types or constructor arguments may be needed | When adding `validate`, `prompt`, and `flow --to reasoning` support |
| Interface/fake seams vs concrete direct calls | Enables fake runtime, buffer-based output checks, and filesystem/path tests, matching gh-cli `pkg/iostreams/iostreams.go:551-568`, restic `internal/backend/mock/backend.go:14-26`, and go-task `executor.go:541-577` | Too many interfaces can obscure simple code and create maintenance overhead | When choosing seams for runtime, store, clock, filesystem, and IO |
| Strict validation before state completion vs permissive artifact existence checks | Reduces false success and supports actionable repair, aligning with typed error and test evidence such as go-task `errors/errors.go:47-50`, restic `internal/errors/fatal.go:10`, and testscript/golden practices in chezmoi `internal/cmd/main_test.go:64-174` | Validators can become brittle if they overfit prose or enforce implementation plans rather than artifact contract | When defining area reasoning and final reasoning validators |
| Non-interactive scriptable UX vs rich terminal feedback | Keeps prompt previews, dry-runs, and validations stable for automation, supported by IOStreams and non-TTY fallback evidence from gh-cli `pkg/iostreams/iostreams.go:52-54` and chezmoi `internal/cmd/prompt.go:20-256` | Less visually rich during runtime-backed generation | When designing output for validation, prompt preview, dry-run, and failure summaries |
| Context propagation everywhere vs only where work can block | Context on runtime-backed flow enables cancellation and timeouts, as in helm `pkg/cmd/install.go:333-347` and go-task `task.go:89` | Passing context through purely deterministic parsing/rendering helpers can add noise | When separating validate/prompt code paths from runtime-backed flow |

## Anti-Patterns And Warnings

- **Do not let command handlers own reasoning business rules:** Command-architecture evidence flags large handlers as a maintainability smell, including opencode `cmd/root.go:49-183`, yq `cmd/evaluate_sequence_command.go:152`, and restic `cmd/restic/cmd_backup.go:84-115` plus the report warning about `runBackup()` length. CLI wiring should not become the reason-stage workflow.

- **Do not introduce global mutable config, runtime, prompt, or validation state:** Dependency and configuration reports warn about globals such as rclone `fs/config.go:793`, yq `pkg/yqlib/lib.go:13`, yq `pkg/yqlib/yaml.go:40`, and opencode `internal/config/config.go:123`. Reason-stage behavior should remain constructable per test and per command invocation.

- **Do not hardcode stdout, stderr, filesystem, runtime, or external command access inside domain logic:** IO evidence flags direct `os.Stdout` and `os.Stderr` usage as testability leaks: rclone `cmd/ls/ls.go:42`, age `cmd/age/tui.go:31`, urfave-cli `cli.go:47`, and yq `pkg/yqlib/logger.go:18`. This sprint has prompt preview and diagnostic requirements that need injectable IO.

- **Do not bypass error wrapping or collapse all failures into unstructured strings:** Error evidence flags non-wrapping and magic exit codes as caution signs, including mitchellh-cli `cli.go:205-206` and `command.go:10-11`. Reason-stage validation and flow errors need cause preservation plus scriptable exit-code mapping.

- **Do not treat runtime success as artifact success:** Error and testing evidence supports explicit classification and output verification. A successful external process without expected files is a false positive; the handbook defers exact rules to reasoning, but warns against state completion without validated artifact existence and content.

- **Do not invoke shell, Git, OpenCode, provider APIs, or agent adapters from sprint domain code directly:** Security evidence favors explicit permission and safe execution boundaries, such as opencode `internal/llm/tools/bash.go:41-55` and gh-cli `pkg/surveyext/editor_manual.go:23`. For this sprint, direct side-effectful calls would cross the selected trust boundary.

- **Do not validate by brittle prose or internal counts:** Testing evidence warns against implementation-detail assertions such as k9s `internal/view/pod_test.go:23` and fragile regex output matching in gh-cli `pkg/cmd/issue/list/list_test.go:88-91`. Reasoning validators should check required artifact contract, selected-context constraints, and placeholders without requiring exact generated wording.

- **Do not allow unselected context to leak into decisions:** The selected-evidence process exists to constrain evidence. Security and configuration reports reinforce explicit input boundaries through k9s `internal/config/json/validator.go:146` and yq `pkg/yqlib/security_prefs.go:3-7`; the sprint-specific equivalent is rejecting unselected contracts, evidence, templates, protocols, and paths as decision sources.

## Examples Worth Inspecting

| Example | Path / Source | Why It Is Useful |
| --- | --- | --- |
| Thin entrypoint and inward delegation | chezmoi `main.go:26-34`; gh-cli `cmd/gh/main.go:6`; k9s `cmd/root.go:112-127` | Shows how command entrypoints avoid owning business rules. |
| Command factory with injected dependencies | gh-cli `pkg/cmdutil/factory.go:16-43`; helm `pkg/cmd/install.go:132-145`; restic `cmd/restic/main.go:37-114` | Useful for reasoning about sprint command wiring without embedding validation or flow logic in CLI handlers. |
| Functional options for testable services | go-task `executor.go:22-24`; go-task `executor.go:541-577` | Shows a Go-native way to configure optional seams, including IO. |
| IOStreams test constructor | gh-cli `pkg/iostreams/iostreams.go:551-568` | Strong reference for command tests that assert stdout/stderr behavior. |
| Filesystem abstraction | chezmoi `internal/chezmoi/system.go:25`; restic `internal/fs/interface.go:10-31` | Useful for workspace-safe store testing and path validation. |
| Error classification and exit mapping | gh-cli `internal/ghcmd/cmd.go:44-49`; go-task `errors/errors.go:47-50`; restic `internal/errors/fatal.go:10` | Relevant to validation failures, unsupported stages, and runtime failures. |
| Context cancellation around long operations | helm `pkg/cmd/install.go:333-347`; go-task `task.go:89`; restic `cmd/restic/cleanup.go:24-38` | Useful when reasoning about runtime-backed flow cancellation and dry-run separation. |
| Prompt/non-TTY fallback discipline | chezmoi `internal/cmd/prompt.go:20-256`; gh-cli `internal/prompter/` | Relevant to prompt preview and no-interactive default behavior. |
| Structured logging and diagnostic routing | helm `internal/logging/logging.go:31-66`; k9s `internal/slogs/keys.go:6-231`; gh-cli `pkg/iostreams/iostreams.go:52-54` | Useful for safe validation diagnostics and runtime error reporting. |
| Security boundary for agent-like commands | opencode `internal/permission/permission.go:44-108`; opencode `internal/llm/tools/bash.go:41-55` | Useful as a contrast: reason-stage code should not recreate agent permission or shell behavior. |
| Secret redaction | restic `internal/options/secret_string.go:15-20`; helm `pkg/registry/transport.go:37-41` | Useful for prompt previews, config summaries, and diagnostic safety. |
| Golden and script-style CLI tests | chezmoi `internal/cmd/main_test.go:64-174`; helm `internal/test/test.go:43`; go-task `task_test.go:166-169` | Useful for command behavior and output regression tests. |

## Design Pressures

- The sprint has two modes with different risk profiles: deterministic validate/prompt/dry-run paths and runtime-backed flow paths.

- Reasoning artifacts are editable Markdown, but validators must still reject missing sections, placeholders, unselected context, and absent prerequisites.

- Selected templates create a conditional stage: area reasoning may be required, or may be skipped only when no templates are selected.

- The command surface must stay scriptable, deterministic, stdout/stderr separated, and no-ANSI by default.

- Runtime-backed generation must be observable enough to diagnose failure without leaking secrets or raw sensitive runtime/config payloads.

- `flow-state.json` is durable state, so completion should represent validated facts, not intention or runtime return code alone.

- Workspace paths and selected template/report paths are trust boundaries; path escape, unreadable selected files, and unselected references need clear diagnostics.

- The sprint must reuse infrastructure carefully: workspace safety, config/redaction, generic runtime, and atomic files are plausible reuse points; study report semantics and run-loop scheduling are not evidence-supported reuse points for planning.

- Tests need to cover behavior through command, service, validator, store, fake runtime, and state seams without invoking real runtime, network, Git, or shell by default.

- Prompt rendering has evidence pressure from configuration, security, and terminal UX reports: include enough context to guide runtime output, but do not include secrets, unnecessary absolute paths, or mutation permission beyond expected output.

## Open Questions For Reasoning

- What exact internal model should represent selected reasoning templates and their deterministic output paths without turning `project-index.md` parsing into a new generic catalog package?

- Where is the smallest correct seam between `internal/app` command handlers and `internal/sprint` service methods for `validate area-reasoning`, `validate reasoning`, `prompt area-reasoning`, `prompt reasoning`, and `flow --to reasoning`?

- How strict should Markdown reasoning validation be so it catches placeholders, missing decisions, missing trade-offs, missing expected evidence, and unselected-context references without becoming brittle prose matching?

- What diagnostics shape gives users actionable file/path/section information while preserving cause chains and exit-code mapping?

- How should flow-state represent selected area reasoning completion, skipped area reasoning, final reasoning completion, and failure without accepting deferred implementation/smoke/review/issues stages?

- Which dependencies need interfaces or fakes now, and which should remain concrete until a second implementation proves the need?

- How should prompt rendering summarize selected contracts and selected evidence without treating unselected governance docs as evidence or copying too much source material into prompts?

- What validation order best prevents runtime work when prerequisites are invalid: requirements, sprint-index, project-index subset checks, handbook, selected templates, then area/final reasoning?

- How should unsupported stages fail: stage parser rejection, validation failure, flow preflight failure, or a shared stage validation helper?

- What minimum command tests prove prompt and validate paths do not invoke runtime, and flow paths invoke only the generic fake runtime seam?

## Evidence Pointers

- `studies/go-cli-study/reports/final/01-project-structure.md`: Inspect `Pattern Catalog` for thin CLI, `internal/` protection, unidirectional imports, and anti-patterns around monolithic packages.

- `studies/go-cli-study/reports/final/02-command-architecture.md`: Inspect `Pattern Catalog` for command factory functions, options structs, lifecycle hooks, command grouping, and warnings about long `RunE` handlers.

- `studies/go-cli-study/reports/final/03-dependency-injection.md`: Inspect `Pattern Catalog` and `Tradeoffs` for composition roots, constructor injection, functional options, interface seams, and global-state cautions.

- `studies/go-cli-study/reports/final/04-configuration-management.md`: Inspect precedence, post-load validation, env-prefix conventions, and anti-patterns around direct environment reads and flag defaults shadowing config.

- `studies/go-cli-study/reports/final/05-error-handling.md`: Inspect typed/sentinel errors, user vs operational separation, exit-code mapping, multi-error aggregation, and anti-patterns around panics and non-wrapping errors.

- `studies/go-cli-study/reports/final/06-io-abstraction.md`: Inspect IOStreams, filesystem interfaces, mock terminals, functional IO options, and direct `os.*` warnings.

- `studies/go-cli-study/reports/final/07-state-context.md`: Inspect signal-context wiring, errgroup use, centralized app state, and warnings about `context.Background()` in long operations.

- `studies/go-cli-study/reports/final/09-terminal-ux.md`: Inspect non-TTY fallback, cancellation behavior, calm progress indicators, and warnings around prompts in non-interactive contexts.

- `studies/go-cli-study/reports/final/10-logging-observability.md`: Inspect structured logging, stdout/stderr separation, debug controls, field consistency, and trace/output abstraction warnings.

- `studies/go-cli-study/reports/final/11-testing-strategy.md`: Inspect testscript, golden fixtures, centralized mocks, fake command runners, and anti-patterns around brittle assertions.

- `studies/go-cli-study/reports/final/13-security.md`: Inspect shell-safe execution, secret redaction, permission request flows, schema validation, credential scrubbing, and trust-boundary cautions.

- `studies/go-cli-study/reports/final/15-philosophy.md`: Inspect factory/options, interface-driven abstraction, wrapper/decorator patterns, explicit non-goals, and warnings around undocumented complexity and god structs.

## Handoff To Reasoning

- Use this handbook as evidence input.
- Validate whether the observed patterns fit this project's constraints.
- Do not copy external patterns without sprint-specific reasoning.
- Treat project docs and requirements as scope context, not as evidence replacing the selected reports.
- Defer final decisions about package boundaries, validator strictness, runtime seams, and flow-state representation to area-specific and final sprint reasoning.
