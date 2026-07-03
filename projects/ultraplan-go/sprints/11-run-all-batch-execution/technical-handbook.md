# Sprint Technical Handbook: 11 Run All Batch Execution

> Project: `ultraplan-go`
> Sprint: `11-run-all-batch-execution`
> Source: `sprint-index.md`
> **Inputs Used:** `projects/ultraplan-go/sprints/11-run-all-batch-execution/sprint-index.md`, `templates/technical-handbook.md`, `projects/ultraplan-go/project-index.md`, `projects/ultraplan-go/sprints/11-run-all-batch-execution/requirements.md`, `projects/ultraplan-go/docs/PRD.md`, `projects/ultraplan-go/docs/TRD.md`, `projects/ultraplan-go/docs/ARCHITECTURE.md`, `studies/go-cli-study/reports/final/01-project-structure.md`, `studies/go-cli-study/reports/final/02-command-architecture.md`, `studies/go-cli-study/reports/final/03-dependency-injection.md`, `studies/go-cli-study/reports/final/04-configuration-management.md`, `studies/go-cli-study/reports/final/05-error-handling.md`, `studies/go-cli-study/reports/final/06-io-abstraction.md`, `studies/go-cli-study/reports/final/07-state-context.md`, `studies/go-cli-study/reports/final/08-concurrency.md`, `studies/go-cli-study/reports/final/09-terminal-ux.md`, `studies/go-cli-study/reports/final/10-logging-observability.md`, `studies/go-cli-study/reports/final/11-testing-strategy.md`, `studies/go-cli-study/reports/final/13-security.md`, `studies/go-cli-study/reports/final/14-performance.md`, `studies/go-cli-study/reports/final/15-philosophy.md`

This handbook distills the studies and reports selected by `sprint-index.md` for sprint reasoning. It does not decide architecture or implementation.

## Selected Studies And Reports

| Study / Report | Path | Relevant Finding | Confidence |
| --- | --- | --- | --- |
| `01-project-structure` | `studies/go-cli-study/reports/final/01-project-structure.md` | Mature Go CLIs keep entrypoints thin and business behavior inside protected/domain packages; command code imports business logic, never the reverse. Evidence includes chezmoi `main.go:16`, k9s `cmd/root.go:112-127`, restic `cmd/restic/main.go:37-114`, and yq `cmd/root.go:9`. | high |
| `02-command-architecture` | `studies/go-cli-study/reports/final/02-command-architecture.md` | Command factories and thin `RunE` delegates are the dominant scalable pattern; large `RunE` functions are flagged as logic leakage. Evidence includes gh-cli `pkg/cmdutil/factory.go:16-43`, helm `pkg/cmd/install.go:132-145`, rclone `cmd/cmd.go:240-340`, and opencode `cmd/root.go:49-183` as a caution. | high |
| `03-dependency-injection` | `studies/go-cli-study/reports/final/03-dependency-injection.md` | Manual composition roots, constructor injection, factory functions, and narrow interfaces enable testable CLIs without DI frameworks. Evidence includes gh-cli `pkg/cmd/factory/default.go:26-46`, go-task `executor.go:22-24`, restic `internal/restic/repository.go:18-66`, and opencode `internal/app/app.go:42-81`. | high |
| `04-configuration-management` | `studies/go-cli-study/reports/final/04-configuration-management.md` | Effective config needs explicit precedence and post-merge validation; CLI flags should override env/config/defaults without defaults accidentally shadowing config. Evidence includes chezmoi `internal/cmd/config.go:2253-2287`, go-task `internal/flags/flags.go:314-327`, k9s `internal/config/k9s.go:423-451`, and gdu `cmd/gdu/main.go:46-112` as a caution. | high |
| `05-error-handling` | `studies/go-cli-study/reports/final/05-error-handling.md` | Production CLI errors preserve chains with `%w`, classify user-visible outcomes, and map typed/sentinel errors to actionable diagnostics and exit behavior. Evidence includes rclone `fs/fserrors/error.go:22-29`, gh-cli `internal/ghcmd/cmd.go:44-49`, go-task `errors/errors_task.go:13-32`, and age `cmd/age/tui.go:37-54`. | high |
| `06-io-abstraction` | `studies/go-cli-study/reports/final/06-io-abstraction.md` | Testable CLIs inject stdout/stderr/filesystem/runtime-like boundaries instead of hardcoding `os.*`; test constructors and buffers are common. Evidence includes gh-cli `pkg/iostreams/iostreams.go:551-568`, restic `internal/ui/terminal.go:10-36`, go-task `executor.go:541-577`, and opencode `internal/app/app.go:156` as a caution. | high |
| `07-state-context` | `studies/go-cli-study/reports/final/07-state-context.md` | Long-running work should receive a root context, signal-triggered cancellation, and context propagation through I/O/parallel work. Evidence includes gh-cli `internal/ghcmd/cmd.go:142`, helm `pkg/cmd/install.go:333-347`, go-task `task.go:89`, and restic `cmd/restic/cleanup.go:24-38`. | high |
| `08-concurrency` | `studies/go-cli-study/reports/final/08-concurrency.md` | Bounded parallelism appears as `errgroup.WithContext`, worker pools, or semaphore channels; unbounded goroutine-per-item and unread channels are explicit cautions. Evidence includes go-task `task.go:87`, gdu `pkg/analyze/parallel.go:13`, k9s `internal/pool.go:21,30,37`, and gh-cli `pkg/cmd/extension/manager.go:196-206` as a caution. | high |
| `09-terminal-ux` | `studies/go-cli-study/reports/final/09-terminal-ux.md` | Human output should match interaction density: batch commands need concise progress/status and non-TTY-safe deterministic output rather than full TUI complexity. Evidence includes restic `internal/ui/terminal.go:10-34`, rclone `cmd/progress.go:24-71`, chezmoi `internal/cmd/prompt.go:124-137`, and helm `pkg/action/package.go:200-208` as a CI-prompt caution. | medium |
| `10-logging-observability` | `studies/go-cli-study/reports/final/10-logging-observability.md` | Structured logs, debug control, and stdout/stderr separation are the operational baseline; diagnostic output should not pollute user data output. Evidence includes k9s `internal/slogs/keys.go:6-231`, helm `internal/logging/logging.go:31-66`, gh-cli `pkg/iostreams/iostreams.go:52-54`, and fzf `src/core.go:325` as a caution. | high |
| `11-testing-strategy` | `studies/go-cli-study/reports/final/11-testing-strategy.md` | Behavior-focused tests, fake dependencies, golden/fixture output, and command-level integration catch CLI regressions; sparse runtime tests are a major risk. Evidence includes gh-cli `acceptance/acceptance_test.go:26-29`, helm `internal/test/test.go:43`, lazygit `pkg/commands/oscommands/fake_cmd_obj_runner.go:17-26`, and opencode `internal/tui/theme/theme_test.go:1-89` as a sparse-coverage caution. | high |
| `13-security` | `studies/go-cli-study/reports/final/13-security.md` | Trust boundaries should be explicit: argument arrays, secret redaction types, permission gates, config validation, and safe temp/path handling. Evidence includes restic `internal/options/secret_string.go:15-20`, helm `pkg/registry/transport.go:37-41`, opencode `internal/permission/permission.go:44-108`, and dive `cmd/dive/cli/internal/command/build.go:25` as a caution. | high |
| `14-performance` | `studies/go-cli-study/reports/final/14-performance.md` | Fast and scalable CLIs defer expensive initialization, stream large data, bound goroutines, and avoid unbounded in-memory accumulation. Evidence includes gh-cli `pkg/cmdutil/factory.go:27-42`, age `internal/stream/stream.go:20`, yq `pkg/yqlib/stream_evaluator.go:78-113`, and k9s `internal/model/table.go:200` as a memory caution. | high |
| `15-philosophy` | `studies/go-cli-study/reports/final/15-philosophy.md` | Coherent tools explicitly choose what they do not build; complexity should be accepted deliberately and kept near the owning domain. Evidence includes lazygit `VISION.md:97`, go-task `website/src/docs/experiments/index.md:17-21`, opencode `internal/pubsub/broker.go:10-19`, and gdu `pkg/analyze/parallel.go:13` as a global-state caution. | medium |

## Relevant Patterns

- **Study-owned thin-command delegation:** The selected structural and command evidence favors keeping command wiring thin and putting behavior in the owning domain package. Chezmoi, k9s, restic, and yq show one-way CLI-to-domain dependencies (`main.go:16`, `cmd/root.go:112-127`, `cmd/restic/main.go:37-114`, `cmd/root.go:9`), while helm and gh-cli show command factories delegating to action/domain logic (`pkg/cmd/install.go:132-145`, `pkg/cmdutil/factory.go:16-43`). For this sprint, the pressure is to keep `run-all` orchestration in `internal/study` and leave `internal/app` as argument validation, service invocation, and deterministic rendering only.

- **Factory and composition-root wiring with fakeable runtime seams:** Manual dependency wiring is the prevailing pattern; no selected report found DI frameworks in mature Go CLIs. Gh-cli uses factory function fields for lazy dependencies (`pkg/cmd/factory/default.go:26-46`), go-task uses functional options (`executor.go:22-24`), restic exposes repository/backend interfaces (`internal/restic/repository.go:18-66`), and opencode composes services centrally in `app.New` (`internal/app/app.go:42-81`). For `run-all`, reasoning should inspect how to pass the existing runtime dependency through the study service without creating product imports in `internal/platform/runtime`.

- **Explicit configuration precedence before launch:** The configuration report shows CLI/config merging must be resolved before work starts. Chezmoi saves and restores changed flags so CLI values win (`internal/cmd/config.go:2253-2287`), go-task centralizes precedence in generic accessors (`internal/flags/flags.go:314-327`), and restic tracks `Flag.Changed` to avoid confusing defaults with explicit input (`internal/global/global.go:139,147`). This matters for parallelism, runtime/model, timeout, health checks, permission policy, and command overrides; invalid parallelism should fail before any runtime task is queued.

- **Context-propagated bounded concurrency:** The sprint needs batch execution but not durable orchestration. The evidence supports a bounded coordination primitive rather than open-ended goroutine spawning: `errgroup.WithContext` for fail-fast fan-out (`go-task/task.go:87`, `restic/internal/repository/repository.go:567`), semaphore channels for explicit limits (`gdu/pkg/analyze/parallel.go:13`, `k9s/internal/pool.go:21,30,37`), and signal-to-context wiring (`helm/pkg/cmd/install.go:333-347`, `restic/cmd/restic/cleanup.go:24-38`). Reasoning must decide whether task failures cancel siblings or permit partial completion, but whichever choice is made should keep launch sites localized and concurrency bounded.

- **Runtime success plus artifact validation as product gate:** The error and testing evidence aligns with the PRD principle that runtime success is insufficient. Helm, gh-cli, rclone, and restic wrap errors and preserve classification (`pkg/storage/driver/driver.go:27-48`, `internal/ghcmd/cmd.go:44-49`, `fs/fserrors/error.go:22-29`, `internal/errors/fatal.go:10`), while testing evidence emphasizes behavior assertions and golden output (`helm/internal/test/test.go:43`, `go-task/task_test.go:166-169`). For `run-all`, each analysis and synthesis task should be considered successful only after the expected report exists and validates.

- **Injected IO and deterministic output capture:** Gh-cli's `IOStreams.Test()` (`pkg/iostreams/iostreams.go:551-568`), restic's terminal interface (`internal/ui/terminal.go:10-36`), go-task's `WithStdout`/`WithStderr` options (`executor.go:541-577`), and gdu's `CreateStdoutUI(output io.Writer)` (`stdout/stdout.go:45-58`) show a testable output seam. For this sprint, deterministic human output and command tests are easier if `run-all` rendering can be captured without real terminal state.

- **Event consumption and observability without leaking unsafe payloads:** Opencode's pub/sub and event patterns (`internal/pubsub/broker.go:10-19`) and logging report findings around structured logs (`internal/logging/logger.go:25-62`, `internal/slogs/keys.go:6-231`) support consuming runtime events as diagnostics rather than printing raw payloads. The security report cautions that raw command/runtime details can expose secrets and unsafe native payloads; restic `SecretString` redacts by type (`internal/options/secret_string.go:15-20`) and helm scrubs Authorization headers (`pkg/registry/transport.go:37-41`). `run-all` should preserve safe diagnostics while avoiding prompt/source-body/runtime payload leakage.

- **Fake-runtime and fixture-first testing:** The testing report favors centralized mocks, fake runners, golden/fixture comparisons, and command-level integration: lazygit records command-runner calls without executing Git (`pkg/commands/oscommands/fake_cmd_obj_runner.go:17-26`), gh-cli uses testscript acceptance (`acceptance/acceptance_test.go:26-29`) and HTTP mocks (`pkg/httpmock/stub.go:35-199`), and helm uses shared golden assertions (`internal/test/test.go:43`). For this sprint, fake runtime tests should cover bounded parallelism, filtering, validation gates, synthesis gating, cancellation, summary determinism, and CLI output.

## Trade-Offs

| Trade-Off | Benefit | Cost | When It Matters |
| --- | --- | --- | --- |
| `errgroup.WithContext` fail-fast vs worker pool that records partial completion | `errgroup` gives automatic cancellation and clear error collection (`go-task/task.go:87`, `restic/internal/repository/repository.go:567`). A worker pool can preserve completed successes and report failed/skipped/pending counts. | Fail-fast may conflict with partial-completion semantics; worker pools require more explicit error aggregation and cancellation handling. | Decide how `run-all` behaves when one task fails after other tasks are active or pending. |
| Static task matrix upfront vs incremental scheduling | A full deterministic matrix supports stable ordering, summaries, and preflight filter validation. Incremental scheduling can reduce memory and adapt as synthesis unlocks. | Full matrices can over-model skipped/inapplicable work; incremental scheduling is harder to reason about and test. | Applies to building selected dimension/source pairs, Markdown applicability skipping, and synthesis gating. |
| Thin CLI rendering vs richer progress UX | Thin deterministic output is script-friendly and testable, matching gh-cli/restic injected IO patterns (`pkg/iostreams/iostreams.go:551-568`, `internal/ui/terminal.go:10-36`). Rich progress improves perceived responsiveness for long runs. | Rich progress introduces TTY/non-TTY branching and nondeterministic output. | Sprint requires concise deterministic human output; reasoning should determine whether progress is just final counts or includes live task updates. |
| Persisting optional run state vs keeping `run-all` ephemeral | Optional state can make status truthful and reviewable, aligning with state/context patterns. Ephemeral operation avoids drifting into durable `run-loop`. | Any state update risks implying resumability, migrations, locks, stale recovery, and future Sprint 12 scope. | Sprint allows run-state primitive updates only if needed for truthful status; durable orchestration is explicitly excluded. |
| Centralized batch result types vs reusing existing single-run result types | Centralized result types make summary rendering and partial completion straightforward. Reuse minimizes new surface area and preserves existing semantics. | New types can become a premature workflow engine; over-reuse can obscure per-task states like skipped, pending, failed validation, and synthesis blocked. | Applies to `internal/study/domain.go`, `run_all.go`, and command output. |
| Configurable parallelism from flags vs only workspace defaults | Flag overrides are consistent with CLI precedence evidence and sprint requirements. Defaults reduce command complexity. | Flag handling increases validation and precedence risk, especially if defaults shadow config values (`cmd/gdu/main.go:46-112`). | Applies to `--parallel` or equivalent command override behavior. |

## Anti-Patterns And Warnings

- **Do not let CLI code become the scheduler:** Large `RunE` functions and command-layer logic leakage are recurring cautions: opencode's `cmd/root.go:49-183`, yq's long command logic, and age's command-contained logic (`cmd/age/age.go:105-321`) are cited as concerns. Keep batch behavior in the study module, not `internal/app`.

- **Do not spawn unbounded goroutines or leave events unread:** The concurrency report flags goroutine-per-item launch sites (`gh-cli/pkg/cmd/extension/manager.go:196-206`), unbuffered channels causing deadlock (`rclone` batch client `client.go:981`), and missing `WaitGroup` timeouts (`gh-cli/manager.go:205`). `run-all` must cap runtime calls and drain or consume every started run's event stream.

- **Do not treat runtime exit as success:** The PRD principle is reinforced by error and testing patterns. A task that returns from runtime without a valid report should not be marked completed; validation failures need cause-preserving internal errors and safe user diagnostics.

- **Do not bypass config precedence or validate after launch:** Gdu's flag-default shadowing (`cmd/gdu/main.go:46-112`) and silent config error handling (`cmd/gdu/main.go:242-244`) are explicit cautions. Invalid source/dimension filters and parallelism less than 1 must fail before runtime work starts.

- **Do not leak prompts, Markdown bodies, secrets, or unsafe runtime payloads:** The security report's `SecretString` and credential-scrubbing patterns (`internal/options/secret_string.go:15-20`, `pkg/registry/transport.go:37-41`) highlight why diagnostics need redaction. Runtime native payloads and embedded source documents should stay out of human output.

- **Do not use globals for batch coordination:** Global config/state/semaphore patterns are repeatedly marked risky: rclone's global config fallback (`fs/config.go:793`), gdu's package-level `concurrencyLimit` (`pkg/analyze/parallel.go:13`), and yq's mutable globals (`pkg/yqlib/yaml.go:40`). `run-all` should use per-request/per-service state.

- **Do not introduce plugin/workflow-engine abstractions for this sprint:** The philosophy report warns that plugins and extensibility add security and maintenance surface. The sprint excludes generic workflow engines, DAG systems, plugin systems, durable retry scheduling, and `run-loop`.

- **Do not make tests depend on real OpenCode, network, or subprocesses:** Sparse or external-dependency-heavy test suites are flagged as risks (`opencode/internal/tui/theme/theme_test.go:1-89`, dive Docker fixtures). The sprint requires fake-runtime and local fixture verification.

## Examples Worth Inspecting

| Example | Path / Source | Why It Is Useful |
| --- | --- | --- |
| Thin command to action delegation | `helm/pkg/cmd/install.go:132-145`, `pkg/action/install.go:73-140` | Shows command wiring that constructs options and delegates to domain/action logic. |
| CLI dependency factory | `gh-cli/pkg/cmdutil/factory.go:16-43`, `pkg/cmd/factory/default.go:26-46` | Useful for understanding lazy dependency functions, IOStreams injection, and command test seams. |
| Bounded semaphore pattern | `gdu/pkg/analyze/parallel.go:13`, `k9s/internal/pool.go:21,30,37` | Directly relevant to configured maximum analysis runtime concurrency. |
| Structured fan-out with cancellation | `go-task/task.go:87`, `restic/internal/repository/repository.go:567` | Useful when reasoning about cancellation propagation and error aggregation. |
| Signal-to-context cancellation | `helm/pkg/cmd/install.go:333-347`, `restic/cmd/restic/cleanup.go:24-38` | Useful for command cancellation semantics and graceful shutdown. |
| IO capture for command tests | `gh-cli/pkg/iostreams/iostreams.go:551-568`, `restic/internal/ui/mock.go:10-53`, `go-task/executor_test.go:146-151` | Useful for deterministic CLI output assertions. |
| Fake command/runtime seam | `lazygit/pkg/commands/oscommands/fake_cmd_obj_runner.go:17-26`, `restic/internal/backend/mock/backend.go:14-26` | Useful for fake-runtime tests that record calls and inject failures without real subprocesses. |
| Error typing and exit mapping | `go-task/errors/errors.go:47-50`, `gh-cli/internal/ghcmd/cmd.go:44-49`, `urfave-cli/errors.go:95-99` | Useful for partial completion, runtime failure, validation failure, and cancellation exit behavior. |
| Secret redaction | `restic/internal/options/secret_string.go:15-20`, `helm/pkg/registry/transport.go:37-41` | Useful for safe diagnostics and runtime metadata redaction. |
| Streaming and bounded memory | `age/internal/stream/stream.go:20`, `yq/pkg/yqlib/stream_evaluator.go:78-113`, `lazygit/pkg/tasks/tasks.go:189-217` | Useful for summary/report processing and avoiding unnecessary full-source scans. |

## Design Pressures

- `run-all` is batch execution but not durable orchestration; it needs bounded workers and partial completion without Sprint 12 `run-loop` semantics.

- Analysis tasks must reuse existing single-run prompt, runtime request, workdir, permission, validation, and source-applicability behavior instead of creating a second execution path.

- Synthesis has a stricter gate than analysis: it can start only when all applicable selected per-source reports for that dimension are valid.

- Markdown document sources add asymmetric semantics: inapplicable source/dimension pairs are skipped, not missing, not failed, and not counted in totals.

- Human output must be deterministic and concise while still exposing enough counts and artifact paths to be useful after long-running work.

- Runtime event channels must be consumed for every started task; otherwise, worker deadlocks can occur even if the runtime call itself is correct.

- Partial completion requires preserving successful work and failed diagnostics together; a single returned error value is probably not enough evidence for rendering.

- Summary generation is both output and validation-adjacent: it must distinguish missing ratings, ambiguous ratings, missing reports, and inapplicable cells deterministically.

- Architecture boundaries create tension with convenience: `internal/study` should own batch behavior, but `internal/app` owns CLI flags/output and platform runtime must remain generic.

- Normal verification must stay deterministic and offline, so test seams must be available before implementation relies on real runtime behavior.

## Open Questions For Reasoning

- Should a failed analysis task stop scheduling new analysis work immediately, or should `run-all` continue launching remaining selected tasks until cancellation or worker queue exhaustion to maximize partial completion?

- What exact task states are needed in the batch result to render completed, failed, skipped, pending/not-run, synthesis skipped, cancellation, and validation warnings without introducing durable run-loop state?

- How should synthesis scheduling interleave with remaining analysis work: only after all analysis workers finish, or as each dimension becomes unblocked?

- What is the precise configured/flag precedence for parallelism, runtime, model, timeout, permission policy, and health requirements, and where should less-than-1 parallelism be rejected?

- How should cancellation be represented when some tasks completed, some active runtime calls return cancellation, and some tasks were never started?

- What safe diagnostic shape is sufficient for user output while preserving original errors internally for tests/review?

- Should `run-all` update existing run-state/status primitives for truthful status, or should it remain entirely ephemeral for this sprint?

- What exact CSV ordering and total semantics should summary generation use when ratings are missing, ambiguous, or inapplicable?

- How should selected filters interact with Markdown applicability for synthesis: should synthesis require all applicable selected sources only, or all applicable study sources even when source filters are used?

- What command output is stable enough for tests and users without creating a public JSON-output contract that the sprint excludes?

## Evidence Pointers

- `studies/go-cli-study/reports/final/01-project-structure.md`: Inspect Pattern 1, Pattern 4, tradeoffs, and evidence around `internal/` boundaries before deciding package placement.

- `studies/go-cli-study/reports/final/02-command-architecture.md`: Inspect Factory Function Command Creation, `cmd.Run()` wrapper, command grouping, and RunE anti-patterns before shaping `internal/app` wiring.

- `studies/go-cli-study/reports/final/03-dependency-injection.md`: Inspect constructor injection, functional options, interface abstraction, and globals anti-patterns before wiring runtime and fake-runtime tests.

- `studies/go-cli-study/reports/final/04-configuration-management.md`: Inspect explicit precedence, flag-restore, generic accessor, and post-load validation before deciding parallelism/config overrides.

- `studies/go-cli-study/reports/final/05-error-handling.md`: Inspect typed errors, sentinel errors, multi-error aggregation, and exit code mapping for partial completion and validation failures.

- `studies/go-cli-study/reports/final/06-io-abstraction.md`: Inspect IOStreams, terminal interfaces, mock UI, and direct `os.*` cautions for command testability.

- `studies/go-cli-study/reports/final/07-state-context.md`: Inspect signal-context wiring, errgroup, centralized state, and `context.Background()` cautions for cancellation semantics.

- `studies/go-cli-study/reports/final/08-concurrency.md`: Inspect errgroup, semaphore channels, wait-with-timeout, and goroutine leak/deadlock cautions for batch worker design.

- `studies/go-cli-study/reports/final/09-terminal-ux.md`: Inspect non-TTY fallback, progress/status patterns, and CI prompt cautions for deterministic human output.

- `studies/go-cli-study/reports/final/10-logging-observability.md`: Inspect structured logging, stdout/stderr separation, and debug control before exposing runtime diagnostics.

- `studies/go-cli-study/reports/final/11-testing-strategy.md`: Inspect fake runners, golden output, testscript, and sparse coverage cautions before defining unit and command tests.

- `studies/go-cli-study/reports/final/13-security.md`: Inspect secret redaction, permission dialogs, argument arrays, and unsafe shell/default-allow cautions before rendering diagnostics or runtime metadata.

- `studies/go-cli-study/reports/final/14-performance.md`: Inspect lazy initialization, streaming, concurrency bounding, and unbounded accumulation cautions before summary/report processing.

- `studies/go-cli-study/reports/final/15-philosophy.md`: Inspect explicit non-goals, deliberate complexity, and plugin/global-state cautions to keep this sprint from becoming durable orchestration or a workflow engine.

## Handoff To Reasoning

- Use this handbook as evidence input.
- Validate whether the observed patterns fit this project's constraints.
- Do not copy external patterns without sprint-specific reasoning.
- Defer implementation decisions to `reasoning/run-all-batch-execution.md` and consolidated `reasoning.md`.
- Keep the selected evidence boundary intact: this handbook used only reports selected by `sprint-index.md`, with project docs and requirements used for scope context.
