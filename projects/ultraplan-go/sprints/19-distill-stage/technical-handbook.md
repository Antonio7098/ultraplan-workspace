# Sprint Technical Handbook: Distill Stage

> Project: `ultraplan-go`
> Sprint: `19-distill-stage`
> Source: `sprint-index.md`
> **Inputs Used:** `projects/ultraplan-go/sprints/19-distill-stage/sprint-index.md`, `templates/technical-handbook.md`, `projects/ultraplan-go/project-index.md`, `projects/ultraplan-go/sprints/19-distill-stage/requirements.md`, `projects/ultraplan-go/docs/PRD.md`, `projects/ultraplan-go/docs/TRD.md`, `projects/ultraplan-go/docs/ARCHITECTURE.md`, `studies/go-cli-study/reports/final/01-project-structure.md`, `studies/go-cli-study/reports/final/02-command-architecture.md`, `studies/go-cli-study/reports/final/03-dependency-injection.md`, `studies/go-cli-study/reports/final/04-configuration-management.md`, `studies/go-cli-study/reports/final/05-error-handling.md`, `studies/go-cli-study/reports/final/06-io-abstraction.md`, `studies/go-cli-study/reports/final/07-state-context.md`, `studies/go-cli-study/reports/final/09-terminal-ux.md`, `studies/go-cli-study/reports/final/10-logging-observability.md`, `studies/go-cli-study/reports/final/11-testing-strategy.md`, `studies/go-cli-study/reports/final/13-security.md`, `studies/go-cli-study/reports/final/15-philosophy.md`

This handbook distills the reports selected by `sprint-index.md` for sprint reasoning. It does not decide architecture or implementation.

## Selected Studies And Reports

| Study / Report | Path | Relevant Finding | Confidence |
| --- | --- | --- | --- |
| `01-project-structure` | `studies/go-cli-study/reports/final/01-project-structure.md` | Mature Go CLIs keep entrypoints thin and route work inward; examples include tiny entrypoints and unidirectional dependencies such as `cmd/gh/main.go:6`, `main.go:16`, `cmd/root.go:112-127`, and `internal/chezmoi/chezmoi.go:1-2`. | high |
| `02-command-architecture` | `studies/go-cli-study/reports/final/02-command-architecture.md` | Command factories with injected dependencies keep command wiring separate from behavior, as seen in `pkg/cmdutil/factory.go:16-43`, `pkg/cmd/install.go:132-145`, and `cmd/restic/cmd_backup.go:84-115`. | high |
| `03-dependency-injection` | `studies/go-cli-study/reports/final/03-dependency-injection.md` | Central composition roots and explicit collaborators improve testability; evidence includes `internal/ghcmd/cmd.go:52-132`, `pkg/cmd/factory/default.go:26-46`, `executor.go:93-114`, and `internal/app/app.go:42-81`. | high |
| `04-configuration-management` | `studies/go-cli-study/reports/final/04-configuration-management.md` | Config handling works best when precedence is explicit and validation runs after merging; evidence includes `internal/cmd/config.go:2253-2287`, `internal/flags/flags.go:314-327`, and `internal/config/config.go:609-641`. | high |
| `05-error-handling` | `studies/go-cli-study/reports/final/05-error-handling.md` | User-facing diagnostics and operational errors are separated in stronger CLIs; evidence includes `cmd/age/tui.go:37-54`, `internal/model/flash.go:100-103`, `cmd/restic/main.go:199-209`, and exit-code mapping at `internal/ghcmd/cmd.go:44-49`. | high |
| `06-io-abstraction` | `studies/go-cli-study/reports/final/06-io-abstraction.md` | Injectable IO and filesystem seams enable command tests without real terminals or external resources; examples include `pkg/iostreams/iostreams.go:551-568`, `executor.go:553-564`, `internal/chezmoi/system.go:25`, and `internal/ui/mock.go:10-53`. | high |
| `07-state-context` | `studies/go-cli-study/reports/final/07-state-context.md` | Long-running or runtime-backed CLI work benefits from root context propagation and explicit state; evidence includes `internal/ghcmd/cmd.go:142`, `task.go:89`, `pkg/cmd/install.go:333-347`, and `internal/restic/lock.go:290-305`. | high |
| `09-terminal-ux` | `studies/go-cli-study/reports/final/09-terminal-ux.md` | Scriptable CLIs need non-TTY fallback, calm progress, and cancellation-aware output rather than TUI complexity by default; evidence includes `internal/cmd/prompt.go:124-137`, `internal/prompter/prompter.go:16-81`, and `signals.go:11-31`. | medium |
| `10-logging-observability` | `studies/go-cli-study/reports/final/10-logging-observability.md` | Observability is stronger when structured diagnostics are separated from user output; evidence includes `fs/log/slog.go:21`, `internal/logging/logging.go:31-66`, `internal/slogs/keys.go:6-231`, and `pkg/iostreams/iostreams.go:52-54`. | high |
| `11-testing-strategy` | `studies/go-cli-study/reports/final/11-testing-strategy.md` | CLI confidence comes from table-driven tests, fakes/mocks, golden outputs, and integration-like command tests; examples include `internal/cmd/main_test.go:64-174`, `acceptance/acceptance_test.go:26-29`, `internal/test/test.go:43`, and `pkg/httpmock/stub.go:35-199`. | high |
| `13-security` | `studies/go-cli-study/reports/final/13-security.md` | Trust boundaries should be explicit around paths, shell execution, permissions, and secrets; evidence includes `internal/execext/exec.go:59-66`, `internal/permission/permission.go:44-108`, `internal/options/secret_string.go:15-20`, and `pkg/registry/transport.go:37-41`. | high |
| `15-philosophy` | `studies/go-cli-study/reports/final/15-philosophy.md` | Sustainable Go CLI design accepts complexity deliberately and keeps abstractions tied to real product pressure; evidence includes `pkg/cmd/factory/default.go:26-46`, `internal/chezmoi/system.go:25-45`, `internal/pubsub/broker.go:10-19`, and `website/src/docs/experiments/index.md:17-21`. | high |

## Relevant Patterns

- **Module-owned behavior with thin CLI wiring:** The selected structure and command reports repeatedly show command layers as wiring and delegation only. This matters for a distill stage because validation, selected-evidence loading, prompt rendering, and flow behavior are product behavior rather than argument parsing. Evidence includes thin entrypoints at `cmd/gh/main.go:6`, `main.go:16`, and `cmd/root.go:112-127`, command factories at `pkg/cmdutil/factory.go:16-43`, and command/action separation at `pkg/cmd/install.go:132-145`.
- **Factory and composition-root seams for runtime, store, IO, and clock-like collaborators:** The dependency-injection reports show manual Go composition rather than DI frameworks. For this sprint, reasoning should inspect where collaborators are assembled and where fake runtime and filesystem seams enter. Evidence includes `internal/ghcmd/cmd.go:52-132`, `pkg/cmd/factory/default.go:26-46`, `executor.go:93-114`, `internal/app/app.go:42-81`, and functional options at `executor.go:22-24`.
- **Selected manifest plus validation gate before runtime success counts:** The error, state, and testing reports align around explicit state and output validation. Runtime success should not be treated as artifact success; stronger CLIs use typed or classified failures and validate produced artifacts before completion. Evidence includes `TaskError` at `errors/errors.go:47-50`, fatal/user handling at `cmd/restic/main.go:199-209`, `errgroup.WithContext` at `task.go:89`, and testscript/golden verification at `internal/cmd/main_test.go:64-174` and `internal/test/test.go:43`.
- **Injectable IO and output discipline for scriptable commands:** The IO and observability reports show that command tests and scriptable behavior depend on redirectable stdout/stderr and structured diagnostics. Evidence includes `pkg/iostreams/iostreams.go:551-568`, `ui.go:19-43`, `internal/ui/mock.go:10-53`, and logging/output separation at `pkg/iostreams/iostreams.go:52-54`.
- **Context propagation for cancellable runtime-backed flow:** The state/context report favors a root context passed through long-running or external work rather than `context.Background()` in operation internals. For this sprint, the pressure is runtime-backed flow and failure handling, not worker-pool design. Evidence includes `internal/ghcmd/cmd.go:142`, `internal/ghcmd/cmd.go:194`, `pkg/cmd/install.go:333-347`, `pkg/action/install.go:284`, and `cmd/restic/cleanup.go:24-38`.
- **Explicit trust boundaries around selected evidence paths and runtime commands:** Security evidence favors argument arrays, permission systems, secret redaction, and schema validation rather than implicit trust. The selected-evidence loader and prompt renderer have similar path and redaction pressure. Evidence includes `internal/execext/exec.go:59-66`, `internal/permission/permission.go:44-108`, `internal/options/secret_string.go:15-20`, `pkg/registry/transport.go:37-41`, and `internal/config/json/validator.go:146`.
- **Golden, fake-first, and command-level testing:** The testing report supports table-driven unit tests plus command-level verification of stdout/stderr, exit codes, fixtures, fakes, and golden text where output stability matters. Evidence includes testscript at `acceptance/acceptance_test.go:26-29`, virtual filesystems at `internal/chezmoitest/chezmoitest.go:86-92`, HTTP mocks at `pkg/httpmock/stub.go:35-199`, and goldens at `task_test.go:166-169`.

## Trade-Offs

| Trade-Off | Benefit | Cost | When It Matters |
| --- | --- | --- | --- |
| Keep distill behavior inside the owning sprint module vs. extract global planning/validation/prompt packages | Preserves module-owned behavior and unidirectional dependency pressure seen in `cmd/` to business-logic designs (`internal/chezmoi/chezmoi.go:1-2`, `cmd/root.go:112-127`) | May duplicate some mechanical Markdown or validation helpers temporarily | Matters when reasoning decides where selected-evidence loading, handbook validation, prompt rendering, and flow transitions live |
| Explicit selected-evidence manifest validation vs. permissive file discovery | Gives deterministic inputs and path safety, consistent with schema/path validation patterns (`internal/config/json/validator.go:146`, `internal/execext/exec.go:59-66`) | More diagnostics and fixture coverage are required for missing, unreadable, or mismatched paths | Matters before prompt rendering and before runtime-backed `flow --to technical-handbook` |
| Injectable collaborators and fakes vs. direct filesystem/runtime calls | Enables tests like `pkg/iostreams/iostreams.go:551-568`, fake runners like `fake_cmd_obj_runner.go:17-26`, and command-level validation without real runtime | Adds constructor/service surface area and requires careful dependency grouping | Matters for fake-runtime flow tests, selected-evidence store tests, and command tests proving validate/prompt are runtime-free |
| Structured diagnostics and exit-code mapping vs. plain error strings | Makes validation/runtime/usage failures scriptable and actionable, matching patterns like `internal/ghcmd/cmd.go:44-49` and `errors/errors.go:47-50` | Requires error classification discipline and can overgrow if every condition becomes a custom type | Matters for unsupported stages, invalid evidence, handbook validation failures, and runtime failures |
| Non-TTY scriptable output vs. richer terminal presentation | Keeps commands predictable in CI and tests; evidence favors non-TTY fallback at `internal/cmd/prompt.go:124-137` and IO streams at `pkg/iostreams/iostreams.go:52-54` | Less polished progress UX than BubbleTea/tcell systems | Matters because this sprint's commands are validate, prompt preview, dry-run, and flow rather than interactive TUI work |
| Golden prompt/output tests vs. flexible substring assertions | Golden fixtures catch unintended prompt/help regressions (`internal/test/test.go:43`, `task_test.go:166-169`) | Golden updates can be noisy and brittle if prose is over-specified | Matters for deterministic prompt preview and help text without freezing every sentence unnecessarily |

## Anti-Patterns And Warnings

- **Do not let CLI handlers become business logic:** Reports call out long `RunE` functions and all-in-one mains as maintainability risks, such as `cmd/root.go:49-183`, `cmd/age/age.go:105-321`, and `runBackup()` at `cmd/restic/cmd_backup.go`. For this sprint, app-level command code should not own selected-evidence semantics, no-decision checks, or flow-state rules.
- **Do not read unselected or uncataloged evidence opportunistically:** Security and config evidence favors explicit validation over discovery side effects. Schema validation at `internal/config/json/validator.go:146` and shell/path trust-boundary patterns at `internal/execext/exec.go:59-66` support failing on invalid inputs rather than broad scanning.
- **Do not treat runtime success as sufficient:** Error and testing reports show artifact validation and explicit failure classification are central. Evidence includes runtime/error classification patterns at `errors/errors.go:47-50`, `cmd/restic/main.go:199-209`, and regression verification at `internal/cmd/main_test.go:64-174`.
- **Do not bypass IO abstractions in command output:** Direct `os.Stdout`, `os.Stderr`, or `fmt.Println` paths are called out as testability gaps, including `cmd/ls/ls.go:42`, `cli.go:47`, and `internal/app/app.go:156`. This sprint's validate/prompt/flow output should remain captureable.
- **Do not introduce global mutable config or service singletons for convenience:** Reports flag global config and parser state as hidden coupling, including `fs/config.go:793`, `pkg/yqlib/lib.go:13`, and `pkg/yqlib/yaml.go:40`. Sprint reasoning should inspect whether any proposed state belongs in explicit service/store/runtime collaborators instead.
- **Do not hide diagnostics in logs only:** The error report warns about silent failure logging such as `lazygit/pkg/commands/git_commands/file_loader.go:52`; validation failures should reach the user through actionable command diagnostics while preserving safe operational context.
- **Do not use decision language in the handbook artifact:** The selected sprint scope treats `technical-handbook.md` as distillation, not final reasoning. The handbook should describe observed patterns, pressures, trade-offs, cautions, and questions, leaving final choices to reasoning.

## Examples Worth Inspecting

| Example | Path / Source | Why It Is Useful |
| --- | --- | --- |
| Thin command entrypoints | `cmd/gh/main.go:6`, `main.go:16`, `cmd/root.go:112-127` | Useful when evaluating whether CLI wiring remains thin and behavior lives in the sprint module. |
| Command factories with dependencies | `pkg/cmdutil/factory.go:16-43`, `pkg/cmd/issue/list/list.go:47-118`, `pkg/cmd/install.go:132-145` | Useful for shaping command construction without embedding business rules in handlers. |
| Composition root and lazy collaborator construction | `internal/ghcmd/cmd.go:52-132`, `pkg/cmd/factory/default.go:26-46`, `internal/app/app.go:42-81` | Useful for deciding how services receive runtime, store, config, and IO collaborators. |
| IOStreams test constructor | `pkg/iostreams/iostreams.go:551-568` | Useful for command tests that assert stdout/stderr separation and no ANSI output. |
| Virtual filesystem and mock persistent state | `internal/chezmoitest/chezmoitest.go:86-92`, `internal/chezmoi/mockpersistentstate.go:4-89` | Useful for selected-evidence and flow-state tests without mutating real workspace state. |
| Error exit-code mapping | `internal/ghcmd/cmd.go:44-49`, `errors/errors.go:47-50`, `errors/errors_task.go:13-32` | Useful for mapping usage, validation, runtime, and cancellation failures. |
| Context cancellation from command entry | `internal/ghcmd/cmd.go:142`, `pkg/cmd/install.go:333-347`, `cmd/restic/cleanup.go:24-38` | Useful for runtime-backed flow and cancellation behavior. |
| Structured logging and safe fields | `internal/logging/logging.go:31-66`, `internal/slogs/keys.go:6-231`, `fs/log/slog.go:21` | Useful for safe diagnostics and deterministic logs. |
| Security redaction and permission boundaries | `internal/options/secret_string.go:15-20`, `pkg/registry/transport.go:37-41`, `internal/permission/permission.go:44-108` | Useful for redaction and runtime permission posture. |
| Golden and integration test helpers | `internal/test/test.go:43`, `acceptance/acceptance_test.go:26-29`, `task_test.go:166-169` | Useful for prompt/help/output regression tests. |

## Design Pressures

- The sprint needs evidence distillation, not implementation decisions, so validation must distinguish evidence-backed guidance from decision language.
- Selected evidence paths must be both cataloged and selected, creating pressure for deterministic manifest parsing and clear diagnostics.
- Runtime-backed flow must remain behind a generic runtime boundary, while sprint-specific prompt content and artifact validation remain product-owned.
- `technical-handbook.md` is editable Markdown, but it is also validated for required sections, placeholder rejection, selected-evidence traces, and no-decision wording.
- Prompt preview and dry-run output must be useful without invoking runtime or writing the handbook artifact.
- The flow-state update is a product state transition and should be treated like durable state: validate prerequisites, validate outputs, then atomically mark completion.
- Tests need to prove absence of side effects as much as presence of behavior: no runtime for validate/prompt/dry-run, no unselected evidence reads, no future-stage mutation, and no Git/OpenCode direct calls from sprint code.
- Output must remain calm and scriptable: no ANSI by default, deterministic ordering, stdout/stderr separation, safe redaction, and actionable validation messages.
- The selected reports show many tempting abstractions, but the sprint scope is narrow; reasoning should resist global packages or study-service reuse unless the evidence and project constraints justify it.

## Open Questions For Reasoning

- Where should the line be drawn between `internal/project` catalog validation and `internal/sprint` selected-evidence loading so the sprint module owns distill behavior without reparsing catalogs ad hoc?
- What evidence-trace rules are sufficient for `technical-handbook.md`: report names, selected report paths, source-line citations from reports, or all of these?
- How strict should no-decision validation be without blocking legitimate trade-off wording such as "prefer," "risk," "when to use," or "best fit" from the evidence reports?
- Should prompt preview tests use golden files for the full prompt, focused golden sections, or structured assertions over required variables and manifests?
- What is the minimum service/store/runtime interface surface needed for fake-runtime and filesystem tests without creating a broad planning framework?
- How should diagnostics report multiple invalid evidence entries: aggregate all findings like multi-error patterns (`pkg/action/uninstall.go:232-254`) or fail fast at the first unsafe path?
- What exact flow-state transition should occur when `requirements.md` and `sprint-index.md` validate but runtime generation succeeds with an invalid `technical-handbook.md`?
- How should dry-run present selected evidence without leaking absolute paths or local-only diagnostics unless needed for repair?
- Which command help changes are necessary for discoverability without expanding stable JSON or unsupported future stages?

## Evidence Pointers

- `studies/go-cli-study/reports/final/01-project-structure.md`: Inspect thin CLI, `internal/` boundary, unidirectional dependency, and global-state caution evidence around `cmd/gh/main.go:6`, `internal/chezmoi/chezmoi.go:1-2`, `cmd/root.go:112-127`, and `pkg/yqlib/lib.go:13-21`.
- `studies/go-cli-study/reports/final/02-command-architecture.md`: Inspect command factory, lifecycle hook, command grouping, and long `RunE` cautions around `pkg/cmdutil/factory.go:16-43`, `pkg/cmd/install.go:132-145`, `internal/cmd/config.go:2236-2454`, and `cmd/root.go:49-183`.
- `studies/go-cli-study/reports/final/03-dependency-injection.md`: Inspect composition roots, functional options, interface seams, and global-state risks around `internal/ghcmd/cmd.go:52-132`, `executor.go:22-24`, `internal/app/app.go:42-81`, and `fs/cache/cache.go:16-21`.
- `studies/go-cli-study/reports/final/04-configuration-management.md`: Inspect precedence, post-load validation, env-prefix conventions, and config anti-patterns around `internal/cmd/config.go:2253-2287`, `internal/flags/flags.go:314-327`, `internal/config/config.go:609-641`, and `cmd/gdu/main.go:46-112`.
- `studies/go-cli-study/reports/final/05-error-handling.md`: Inspect `%w` wrapping, sentinels, user/operational separation, multi-error, and exit-code mapping around `pkg/ssh/ssh_keys.go:64`, `cmd/age/tui.go:37-54`, `cmd/restic/main.go:199-209`, `pkg/action/uninstall.go:232-254`, and `internal/ghcmd/cmd.go:44-49`.
- `studies/go-cli-study/reports/final/06-io-abstraction.md`: Inspect IOStreams, filesystem abstractions, mock UI, and direct IO anti-patterns around `pkg/iostreams/iostreams.go:551-568`, `internal/chezmoi/system.go:25`, `internal/ui/mock.go:10-53`, `cmd/ls/ls.go:42`, and `cli.go:47`.
- `studies/go-cli-study/reports/final/07-state-context.md`: Inspect signal-context wiring, errgroup cancellation, centralized state, and context anti-patterns around `internal/ghcmd/cmd.go:142`, `task.go:89`, `pkg/cmd/install.go:333-347`, `fs/config.go:793-802`, and `internal/cmd/templatefuncs.go:215`.
- `studies/go-cli-study/reports/final/09-terminal-ux.md`: Inspect non-TTY fallback, cancellation, progress, and no-TTY cautions around `internal/cmd/prompt.go:124-137`, `internal/prompter/prompter.go:16-81`, `signals.go:11-31`, and `pkg/action/package.go:200-208`.
- `studies/go-cli-study/reports/final/10-logging-observability.md`: Inspect structured logging, output separation, debug controls, and log anti-patterns around `fs/log/slog.go:21`, `internal/logging/logging.go:31-66`, `internal/slogs/keys.go:6-231`, `src/core.go:325`, and `cli.go:46`.
- `studies/go-cli-study/reports/final/11-testing-strategy.md`: Inspect CLI integration, mocks, golden files, and fragile assertions around `internal/cmd/main_test.go:64-174`, `acceptance/acceptance_test.go:26-29`, `pkg/httpmock/stub.go:35-199`, `internal/test/test.go:43`, and `internal/view/pod_test.go:23`.
- `studies/go-cli-study/reports/final/13-security.md`: Inspect path/shell trust boundaries, redaction, permission systems, and validation around `internal/execext/exec.go:59-66`, `internal/permission/permission.go:44-108`, `internal/options/secret_string.go:15-20`, `pkg/registry/transport.go:37-41`, and `internal/config/json/validator.go:146`.
- `studies/go-cli-study/reports/final/15-philosophy.md`: Inspect deliberate complexity, factory/options patterns, interface abstractions, and anti-patterns around `pkg/cmd/factory/default.go:26-46`, `internal/chezmoi/system.go:25-45`, `internal/pubsub/broker.go:10-19`, `pkg/analyze/parallel.go:13`, and `VISION.md:97`.

## Handoff To Reasoning

- Use this handbook as evidence input.
- Validate whether the observed patterns fit this project's constraints.
- Do not copy external patterns without sprint-specific reasoning.
- Resolve the open questions in area-specific reasoning or sprint reasoning before turning the distill-stage implementation into a plan.
