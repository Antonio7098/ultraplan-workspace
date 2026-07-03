# Sprint Technical Handbook: Plan Stage

> Project: `ultraplan-go`
> Sprint: `21-plan-stage`
> Source: `sprint-index.md`
> **Inputs Used:** `projects/ultraplan-go/sprints/21-plan-stage/sprint-index.md`, `templates/technical-handbook.md`, `projects/ultraplan-go/project-index.md`, `projects/ultraplan-go/sprints/21-plan-stage/requirements.md`, `projects/ultraplan-go/docs/PRD.md`, `projects/ultraplan-go/docs/TRD.md`, `projects/ultraplan-go/docs/ARCHITECTURE.md`, `studies/go-cli-study/reports/final/01-project-structure.md`, `studies/go-cli-study/reports/final/02-command-architecture.md`, `studies/go-cli-study/reports/final/03-dependency-injection.md`, `studies/go-cli-study/reports/final/04-configuration-management.md`, `studies/go-cli-study/reports/final/05-error-handling.md`, `studies/go-cli-study/reports/final/06-io-abstraction.md`, `studies/go-cli-study/reports/final/07-state-context.md`, `studies/go-cli-study/reports/final/09-terminal-ux.md`, `studies/go-cli-study/reports/final/10-logging-observability.md`, `studies/go-cli-study/reports/final/11-testing-strategy.md`, `studies/go-cli-study/reports/final/13-security.md`, `studies/go-cli-study/reports/final/15-philosophy.md`

This handbook distills the studies and reports selected by `sprint-index.md` for sprint reasoning. It does not decide architecture or implementation. Project docs and requirements were used for sprint scope only; the evidence base below comes from the selected Go CLI study reports.

## Selected Studies And Reports

| Study / Report | Path | Relevant Finding | Confidence |
| --- | --- | --- | --- |
| `01-project-structure` | `studies/go-cli-study/reports/final/01-project-structure.md` | Mature Go CLIs keep entrypoints thin and enforce one-way dependencies from CLI to business logic; examples include k9s `cmd/root.go:112-127`, restic `cmd/restic/main.go:37-114`, and yq `cmd/root.go:9`. | high |
| `02-command-architecture` | `studies/go-cli-study/reports/final/02-command-architecture.md` | Command factories and thin `RunE` wrappers scale better than command handlers containing business logic; examples include helm `pkg/cmd/install.go:132-145`, gh-cli `pkg/cmdutil/factory.go:16-43`, and chezmoi `internal/cmd/config.go:1833-1845`. | high |
| `03-dependency-injection` | `studies/go-cli-study/reports/final/03-dependency-injection.md` | Explicit composition roots, constructor injection, and narrow service interfaces improve testability; evidence includes gh-cli `pkg/cmd/factory/default.go:26-46`, go-task `executor.go:22-24`, and restic `internal/backend/backend.go:19-90`. | high |
| `04-configuration-management` | `studies/go-cli-study/reports/final/04-configuration-management.md` | Clear config precedence and post-merge validation prevent surprising behavior; examples include chezmoi `internal/cmd/config.go:2253-2287`, go-task `internal/flags/flags.go:314-327`, and restic `internal/global/global.go:139,147`. | medium |
| `05-error-handling` | `studies/go-cli-study/reports/final/05-error-handling.md` | First-class error values, `%w` wrapping, sentinels, and exit-code mapping support scriptable diagnostics; examples include go-task `errors/errors.go:47-50`, gh-cli `internal/ghcmd/cmd.go:44-49`, and restic `internal/errors/fatal.go:10`. | high |
| `06-io-abstraction` | `studies/go-cli-study/reports/final/06-io-abstraction.md` | Injectable stdout/stderr/filesystem/runtime seams make CLI behavior testable; examples include gh-cli `pkg/iostreams/iostreams.go:551-568`, go-task `executor.go:553-564`, and restic `internal/ui/terminal.go:10-36`. | high |
| `07-state-context` | `studies/go-cli-study/reports/final/07-state-context.md` | Context propagation and explicit centralized state matter for cancellable operations; examples include gh-cli `internal/ghcmd/cmd.go:142`, helm `pkg/cmd/install.go:333-347`, and restic `cmd/restic/cleanup.go:24-38`. | medium |
| `09-terminal-ux` | `studies/go-cli-study/reports/final/09-terminal-ux.md` | CLI UX should match task type; scriptable commands need calm output, non-TTY awareness, and no unnecessary TUI behavior. Evidence includes chezmoi `internal/cmd/prompt.go:124-137`, gh-cli `internal/prompter/`, and rclone `lib/terminal/terminal.go:82-86`. | medium |
| `10-logging-observability` | `studies/go-cli-study/reports/final/10-logging-observability.md` | Structured logs and stdout/stderr separation are key operational patterns; examples include helm `internal/logging/logging.go:31-66`, k9s `internal/slogs/keys.go:6-231`, and gh-cli `pkg/iostreams/iostreams.go:52-54`. | high |
| `11-testing-strategy` | `studies/go-cli-study/reports/final/11-testing-strategy.md` | Elite CLIs combine table tests, command-level integration, golden fixtures, and fakes; examples include chezmoi `internal/cmd/main_test.go:64-174`, helm `internal/test/test.go:43`, and rclone `cmd/bisync/bisync_test.go:1435-1479`. | high |
| `13-security` | `studies/go-cli-study/reports/final/13-security.md` | Trust boundaries should be explicit, secrets redacted, paths validated, and shell execution constrained; examples include restic `internal/options/secret_string.go:15-20`, helm `pkg/registry/transport.go:37-41`, and opencode `internal/llm/tools/bash.go:41-55`. | high |
| `15-philosophy` | `studies/go-cli-study/reports/final/15-philosophy.md` | Coherent projects accept complexity deliberately and keep non-goals visible; evidence includes age `age.go:18`, lazygit `VISION.md:97`, and go-task `website/src/docs/experiments/index.md:17-21`. | medium |

## Relevant Patterns

- **Thin CLI With Module-Owned Behavior:** The studies repeatedly separate command wiring from business logic. `01-project-structure` cites k9s delegating from `cmd/root.go:112-127` into `internal/view`, restic using `cmd/restic/main.go:37-114` as a command root over internal packages, and yq keeping `cmd/` imports one-way into `pkg/yqlib` via `cmd/root.go:9`. `02-command-architecture` reinforces thin command factories through helm `pkg/cmd/install.go:132-145` and gh-cli `pkg/cmdutil/factory.go:16-43`. For this sprint, the pressure is to keep `internal/app` command handlers as routing/output code while plan validation, prompt rendering, prerequisite loading, task/evidence tracing, and flow-state semantics remain in `internal/sprint` for reasoning to evaluate.

- **Explicit Composition Root And Constructor Injection:** `03-dependency-injection` finds that manual composition roots dominate mature Go CLIs, with gh-cli's factory wiring at `pkg/cmd/factory/default.go:26-46`, go-task's functional options at `executor.go:22-24`, and restic's backend abstractions at `internal/backend/backend.go:19-90`. The sprint has a similar pressure: command tests need fake runtime/store seams, but reasoning should decide how narrow those seams are so `internal/sprint` does not become a generic DI framework.

- **IOStreams Or Writer Injection For Scriptable Commands:** `06-io-abstraction` highlights `IOStreams.Test()` in gh-cli `pkg/iostreams/iostreams.go:551-568`, go-task's `WithStdout` at `executor.go:553-564`, and restic's `ui.Terminal` at `internal/ui/terminal.go:10-36`. Plan-stage commands need deterministic stdout/stderr separation and no ANSI output by default, so reasoning should account for output injection and capture without creating a large UI abstraction.

- **Validation As The Completion Gate:** The reports' error, testing, and runtime themes converge on not equating command/runtime success with product success. `05-error-handling` shows error classification and exit-code mapping in gh-cli `internal/ghcmd/cmd.go:44-49`, go-task `errors/errors.go:47-50`, and restic `internal/errors/fatal.go:10`. `11-testing-strategy` shows behavior-focused command testing and golden regression checks through chezmoi `internal/cmd/main_test.go:64-174` and helm `internal/test/test.go:43`. For this sprint, the key pressure is that flow through `plan` should be reasoned as complete only after generated `plan.md` exists, validates, and flow state writes safely.

- **Context-Aware Flow And Explicit State:** `07-state-context` finds high-scoring CLIs create root contexts and propagate cancellation into work, such as gh-cli `internal/ghcmd/cmd.go:142`, helm `pkg/cmd/install.go:333-347`, and restic `cmd/restic/cleanup.go:24-38`. Plan flow is sequential rather than a worker pool, but it still has runtime-backed execution, cancellation, and state update pressures. Reasoning should decide how much context propagation is needed for validate/prompt/dry-run versus non-dry-run flow.

- **Structured Diagnostics With Safe Redaction:** `10-logging-observability` emphasizes structured logs, debug controls, and stdout/stderr separation, with helm `internal/logging/logging.go:31-66`, k9s field keys in `internal/slogs/keys.go:6-231`, and gh-cli `pkg/iostreams/iostreams.go:52-54`. `13-security` adds secret redaction evidence from restic `internal/options/secret_string.go:15-20` and helm `pkg/registry/transport.go:37-41`. Plan-stage reasoning should treat runtime errors, invalid artifacts, and flow-state write failures as safe diagnostics rather than raw provider/OpenCode output dumps.

- **Fake-First And Golden/Fixture-Oriented Tests:** `11-testing-strategy` shows command and artifact tests built around testscript, golden files, fixtures, and mocks: chezmoi `internal/cmd/main_test.go:64-174`, go-task `task_test.go:166-169`, helm `internal/test/test.go:43`, and rclone `fstest/mockfs/mockfs.go:31-62`. This fits sprint needs for deterministic prompt rendering, fake-runtime plan flow, plan validation fixtures, and command-level exit/output assertions without real OpenCode or network use.

- **Deliberate Non-Goals As Architecture Control:** `15-philosophy` shows that strong projects keep non-goals visible: age's no-keyring design at `age.go:18`, lazygit's rejection of features not worth complexity at `VISION.md:97`, and go-task's experiments framework for controlled evolution at `website/src/docs/experiments/index.md:17-21`. For this sprint, the pressure is to preserve Phase 2's stop at `plan` and avoid pulling implementation execution, smoke, review, issues, Git mutation, plugin, TUI, or workflow-engine behavior into the plan stage.

## Trade-Offs

| Trade-Off | Benefit | Cost | When It Matters |
| --- | --- | --- | --- |
| Keep plan behavior in `internal/sprint` instead of global `validation`, `prompts`, or `workflow` packages | Mirrors module-owned behavior in `01-project-structure` and `15-philosophy`; keeps ownership clear and avoids cross-module product coupling | Some validation/prompt code may duplicate study-shaped concepts until a stable reusable mechanical helper is proven | When deciding where plan validation, reasoning prerequisite checks, prompt rendering, and flow-state transitions live |
| Explicit service/store/runtime injection vs direct filesystem/runtime calls | Testability follows `03-dependency-injection` and `06-io-abstraction`; fake runtime and fake store tests become straightforward | More constructor parameters and potential central object growth, like chezmoi's large `Config` concern in `internal/cmd/config.go:193-291` | When wiring `internal/sprint.Service`, filesystem store, prompt preview, non-dry-run flow, and command tests |
| Strict plan validation vs permissive editable Markdown | Prevents invalid success and supports scriptable diagnostics, matching validation and testing findings from `05-error-handling` and `11-testing-strategy` | Risk of false negatives if validation asserts brittle prose rather than structural evidence | When defining checks for `reasoning.md` citation, decisions, task checklist, evidence checklist, risks/blockers, success criteria, and forbidden deferred-stage content |
| Golden/fixture tests vs flexible substring assertions | Golden/fixture tests catch prompt/output regressions, as shown by helm `internal/test/test.go:43` and go-task `task_test.go:166-169` | Golden files can become maintenance overhead and can overfit prose | When testing deterministic plan prompt rendering, help output, and validation diagnostics |
| Centralized flow-state model vs minimal artifact-only checks | Explicit state makes failures and recovery visible, matching `07-state-context` state lessons and the project's flow-state pressure | More schema validation and atomic-write error paths to test | When deciding how `plan` complete/failed/skipped states are represented and how unsupported future stages are rejected |
| Runtime-backed generation through generic runtime boundary vs direct OpenCode invocation | Preserves product/platform separation and agentwrap ownership; aligns with security evidence against unconstrained command execution | Requires adapter-facing mapping and fake-runtime tests rather than simple command spawning | When implementing non-dry-run `flow --to plan` and runtime failure diagnostics |

## Anti-Patterns And Warnings

- **Command Handlers Accumulating Plan Rules:** `02-command-architecture` flags long `RunE` functions such as opencode `cmd/root.go:49-183` and yq's long evaluation command as command-layer logic leakage. For this sprint, `internal/app/sprint_commands.go` should not become the owner of plan validation, prompt contents, stage ordering, or flow-state semantics.

- **Global Mutable State And Hidden Initialization:** `03-dependency-injection` identifies rclone globals like `fs/config.go:14-51` and yq's singleton parser at `pkg/yqlib/lib.go:13` as testability risks. Plan-stage state, runtime config, and selected context should not depend on package globals that make tests order-dependent.

- **Direct OS IO Bypassing Test Seams:** `06-io-abstraction` warns about direct `os.Stdout`/`os.Stderr` writes such as rclone `cmd/ls/ls.go:42`, urfave-cli `cli.go:47`, and chezmoi template funcs `internal/cmd/templatefuncs.go:296`. Plan commands need output capture in tests, so direct printing from domain logic is a caution sign.

- **Context Accepted But Ignored:** `07-state-context` calls out misleading context signatures, such as dive accepting context while not propagating it through analysis at `dive/image/analysis.go:20`. Runtime-backed plan flow should not accept a context and then use `context.Background()` for long or cancellable work.

- **Unsafe Or Overbroad Runtime/Shell Behavior:** `13-security` highlights opencode command filtering at `internal/llm/tools/bash.go:41-55`, go-task's pure interpreter at `internal/execext/exec.go:59-66`, and dive's unrestricted `DisableFlagParsing` at `cmd/dive/cli/internal/command/build.go:25`. This sprint should reason carefully about keeping runtime execution behind the generic runtime boundary and not adding shell, Git, provider, or OpenCode process calls to `internal/sprint`.

- **Config Precedence Or Diagnostics That Surprise Users:** `04-configuration-management` warns that flag defaults can shadow config files, as in gdu `cmd/gdu/main.go:46-112`, and that direct environment access can bypass config layers, as in opencode `internal/config/config.go:163`. Plan prompt previews and runtime-backed flow should use resolved config consistently and redact sensitive values.

- **Brittle Tests Of Incidental Details:** `11-testing-strategy` warns against implementation-detail tests like k9s hint count assertions at `internal/view/pod_test.go:23` and fragile regex output matching in gh-cli `pkg/cmd/issue/list/list_test.go:88-91`. Plan-stage tests should assert required variables, paths, sections, exit codes, and traceability rules rather than exact prose unrelated to behavior.

## Examples Worth Inspecting

| Example | Path / Source | Why It Is Useful |
| --- | --- | --- |
| Thin command-to-domain delegation | `helm/pkg/cmd/install.go:132-145`; `gh-cli/pkg/cmdutil/factory.go:16-43` | Shows command factories accepting dependencies and delegating work without embedding domain logic. |
| One-way package boundary | `k9s/cmd/root.go:112-127`; `yq/cmd/root.go:9`; `restic/cmd/restic/main.go:37-114` | Useful when checking that `internal/app` routes into `internal/sprint`, not the reverse. |
| Central factory/composition root | `gh-cli/pkg/cmd/factory/default.go:26-46`; `chezmoi/internal/cmd/config.go:362`; `go-task/executor.go:93-114` | Shows ways to assemble command dependencies while preserving test seams. |
| IO test constructor | `gh-cli/pkg/iostreams/iostreams.go:551-568`; `mitchellh-cli/ui_mock.go:27-33` | Useful for command tests that assert stdout/stderr separation and no ANSI output. |
| User/operational error separation | `age/cmd/age/tui.go:37-54`; `restic/cmd/restic/main.go:199-209`; `gh-cli/internal/ghcmd/cmd.go:281-301` | Useful for safe validation diagnostics and actionable CLI error output. |
| Context cancellation wiring | `helm/pkg/cmd/install.go:333-347`; `go-task/signals.go:18`; `restic/cmd/restic/cleanup.go:24-38` | Relevant to non-dry-run flow cancellation and state preservation. |
| Golden and CLI integration tests | `chezmoi/internal/cmd/main_test.go:64-174`; `helm/internal/test/test.go:43`; `go-task/task_test.go:166-169` | Useful for prompt preview and command-output regression testing strategy. |
| Secret redaction patterns | `restic/internal/options/secret_string.go:15-20`; `helm/pkg/registry/transport.go:37-41`; `gh-cli/status.go:332-338` | Relevant to runtime diagnostics, prompt previews, and config summaries that must not leak credentials. |
| Explicit non-goal discipline | `age/age.go:18`; `lazygit/VISION.md:97`; `go-task/website/src/docs/experiments/index.md:17-21` | Helps reasoning keep future sprint execution, smoke, review, issues, and Git behavior out of plan-stage scope. |

## Design Pressures

- The plan stage has to extend the sprint flow without expanding the product into a general workflow engine.
- `internal/sprint` needs enough ownership to validate, prompt, flow, and update state, but not so much abstraction that it duplicates a framework for all future stages.
- CLI commands must remain scriptable and thin while still surfacing detailed validation and runtime failures with correct exit codes.
- Runtime-backed flow has to use the generic runtime boundary and then validate the expected artifact, not trust runtime success alone.
- Prompt previews and generated prompts must be deterministic, workspace-relative where possible, safe to inspect, and explicit about no-mutation/no-execution rules.
- Flow-state writes are durable product state and must preserve the last valid state on write failure.
- Tests must prove validate and prompt paths are runtime-free, while flow uses a fake runtime seam by default.
- Plan validation must catch missing required structure and forbidden deferred-stage behavior without becoming a brittle prose linter.
- Selected-context discipline matters: the sprint index selects only the evidence reports above, so plan-stage prompt and validation behavior should not silently add unselected evidence as authoritative input.
- Source repository and generated artifact paths are trust boundaries; diagnostics should prefer safe paths and redacted metadata.

## Open Questions For Reasoning

- What is the smallest `internal/sprint` plan-domain model that can validate `plan.md`, trace tasks to `reasoning.md`, and reject forbidden deferred-stage behavior without becoming a generic planner model?
- How should plan validation represent traceability failures: as simple deterministic findings, typed errors, or both?
- Which output assertions are stable enough for command tests, and which prompt text should be covered by required-variable checks rather than exact golden prose?
- How should non-dry-run plan flow record a runtime failure versus post-generation validation failure in `flow-state.json` while preserving safe diagnostics?
- How should the service/store/runtime seams be shaped so validate and prompt commands cannot accidentally invoke runtime, while flow can use fake runtime tests?
- Should prompt preview support writing to an explicit preview path in this sprint, or only stdout, if existing command conventions allow either?
- How should unsupported stages such as `execute`, `smoke`, `review`, and `issues` be rejected consistently across validate, prompt, and flow commands?
- What is the minimum atomic-write helper reuse needed for `flow-state.json` without pulling study run-loop semantics into sprint planning?
- How much of the selected area reasoning prerequisite should be structurally validated before invoking runtime for plan generation?
- What diagnostic detail from runtime failures is safe and useful enough for default stderr versus debug output?

## Evidence Pointers

- `studies/go-cli-study/reports/final/01-project-structure.md`: inspect Pattern 1, Pattern 2, Pattern 4, tradeoffs around `internal/` vs `pkg/`, and anti-patterns for large entrypoints and bidirectional imports.
- `studies/go-cli-study/reports/final/02-command-architecture.md`: inspect factory command creation, thin delegate pattern, command grouping, and warnings about long `RunE` handlers.
- `studies/go-cli-study/reports/final/03-dependency-injection.md`: inspect constructor injection, central composition roots, functional options, and warnings about globals and context-based service locators.
- `studies/go-cli-study/reports/final/04-configuration-management.md`: inspect config precedence, flag-change tracking, post-load validation, redaction implications, and precedence anti-patterns.
- `studies/go-cli-study/reports/final/05-error-handling.md`: inspect `%w` wrapping, sentinel errors, typed errors, user/operational separation, hints, and exit-code mapping.
- `studies/go-cli-study/reports/final/06-io-abstraction.md`: inspect `IOStreams.Test()`, writer injection, mock terminal patterns, direct IO anti-patterns, and filesystem abstraction cautions.
- `studies/go-cli-study/reports/final/07-state-context.md`: inspect signal-context wiring, centralized application state, context anti-patterns, and explicit session/state modeling tradeoffs.
- `studies/go-cli-study/reports/final/09-terminal-ux.md`: inspect non-TTY fallback, calm CLI-first output, no-ANSI expectations, and prompt/progress tradeoffs.
- `studies/go-cli-study/reports/final/10-logging-observability.md`: inspect structured logging, stdout/stderr separation, safe debug controls, and untyped log-field anti-patterns.
- `studies/go-cli-study/reports/final/11-testing-strategy.md`: inspect CLI integration tests, golden fixtures, centralized mocks, fake command runners, and brittle assertion warnings.
- `studies/go-cli-study/reports/final/13-security.md`: inspect shell-safe command construction, secret redaction, permission boundaries, safe diagnostics, and default-allow cautionary examples.
- `studies/go-cli-study/reports/final/15-philosophy.md`: inspect deliberate non-goals, factory/options patterns, interface-driven abstraction, and warnings about undocumented complexity and plugin/security scope creep.

## Handoff To Reasoning

- Use this handbook as evidence input for architecture reasoning and final sprint reasoning.
- Validate whether each observed pattern fits this project's constraints before applying it.
- Keep implementation decisions in `reasoning/*.md`, `reasoning.md`, and `plan.md`, not in this handbook.
- Preserve selected-context discipline: cite only selected evidence as authoritative evidence for this sprint.
- Treat project docs and requirements as sprint scope constraints, not as external evidence reports.
