# Sprint Technical Handbook: Execute Stage

> Project: `ultraplan-go`
> Sprint: `23-execute-stage`
> Source: `projects/ultraplan-go/sprints/23-execute-stage/sprint-index.md`
> **Inputs Used:** `projects/ultraplan-go/project-index.md`, `projects/ultraplan-go/sprints/23-execute-stage/sprint-index.md`, `projects/ultraplan-go/sprints/23-execute-stage/requirements.md`, `projects/ultraplan-go/docs/ARCHITECTURE.md`, `projects/ultraplan-go/docs/PRD.md`, `projects/ultraplan-go/docs/TRD.md`, `studies/go-cli-study/reports/final/01-project-structure.md`, `studies/go-cli-study/reports/final/02-command-architecture.md`, `studies/go-cli-study/reports/final/03-dependency-injection.md`, `studies/go-cli-study/reports/final/04-configuration-management.md`, `studies/go-cli-study/reports/final/05-error-handling.md`, `studies/go-cli-study/reports/final/06-io-abstraction.md`, `studies/go-cli-study/reports/final/07-state-context.md`, `studies/go-cli-study/reports/final/08-concurrency.md`, `studies/go-cli-study/reports/final/09-terminal-ux.md`, `studies/go-cli-study/reports/final/10-logging-observability.md`, `studies/go-cli-study/reports/final/11-testing-strategy.md`, `studies/go-cli-study/reports/final/13-security.md`, `studies/go-cli-study/reports/final/14-performance.md`, `studies/go-cli-study/reports/final/15-philosophy.md`

This handbook distills the studies and reports selected by `sprint-index.md` for sprint reasoning. It does not decide architecture or implementation.

## Selected Studies And Reports

| Study / Report | Path | Relevant Finding | Confidence |
| --- | --- | --- | --- |
| `01-project-structure` | `studies/go-cli-study/reports/final/01-project-structure.md` | Mature Go CLIs keep entrypoints thin and route business behavior into owned internal/domain packages; evidence includes `cmd/gh/main.go:6`, `cmd/root.go:112-127`, and `internal/restic/repository.go:18`. | high |
| `02-command-architecture` | `studies/go-cli-study/reports/final/02-command-architecture.md` | Command factories, shared config/context objects, lifecycle hooks, and wrappers keep command logic thin; evidence includes `pkg/cmd/install.go:132-145`, `pkg/cmdutil/factory.go:16-43`, and `cmd/cmd.go:240-340`. | high |
| `03-dependency-injection` | `studies/go-cli-study/reports/final/03-dependency-injection.md` | Explicit composition roots, constructor injection, functional options, and narrow interfaces improve testability; evidence includes `internal/cmd/config.go:362`, `executor.go:22-24`, and `internal/backend/backend.go:19-90`. | high |
| `04-configuration-management` | `studies/go-cli-study/reports/final/04-configuration-management.md` | Reliable CLIs make precedence explicit, validate after merging sources, and track changed flags; evidence includes `internal/cmd/config.go:2253-2287`, `internal/flags/flags.go:314-327`, and `internal/global/global.go:139,147`. | high |
| `05-error-handling` | `studies/go-cli-study/reports/final/05-error-handling.md` | Mature CLIs use wrapped errors, typed/sentinel classifications, user-operational separation, and exit-code mapping; evidence includes `fs/fserrors/error.go:22-29`, `internal/errors/fatal.go:10`, and `internal/ghcmd/cmd.go:44-49`. | high |
| `06-io-abstraction` | `studies/go-cli-study/reports/final/06-io-abstraction.md` | IO, filesystem, terminal, and backend seams enable deterministic command and state tests; evidence includes `pkg/iostreams/iostreams.go:551-568`, `internal/fs/interface.go:10-31`, and `internal/backend/backend.go:19-90`. | high |
| `07-state-context` | `studies/go-cli-study/reports/final/07-state-context.md` | Long-running CLIs propagate `context.Context`, wire signals to cancellation, and model persistent task/session state explicitly; evidence includes `internal/ghcmd/cmd.go:142`, `task.go:89`, `internal/session/session.go:12-23`, and `internal/restic/lock.go:105`. | high |
| `08-concurrency` | `studies/go-cli-study/reports/final/08-concurrency.md` | Bounded structured concurrency avoids leaks and unbounded goroutine growth; evidence includes `task.go:87`, `gdu/pkg/analyze/parallel.go:13`, and `opencode/cmd/root.go:130,252`. | high |
| `09-terminal-ux` | `studies/go-cli-study/reports/final/09-terminal-ux.md` | Status UX must fit operation length: progress/cancellation for long operations, non-TTY fallback for scripts, and concise output for CLI-first flows; evidence includes `internal/cmd/prompt.go:20-256`, `signals.go:11-31`, and `internal/ui/terminal.go:10-34`. | medium |
| `10-logging-observability` | `studies/go-cli-study/reports/final/10-logging-observability.md` | Structured logs, stdout/stderr separation, debug controls, and consistent fields support recovery; evidence includes `internal/logging/logging.go:31-66`, `internal/slogs/keys.go:6-231`, and `fs/accounting/prometheus.go:78-108`. | high |
| `11-testing-strategy` | `studies/go-cli-study/reports/final/11-testing-strategy.md` | Strong CLI projects combine table tests, CLI/command integration tests, mocks/fakes, and golden outputs; evidence includes `internal/cmd/main_test.go:64-174`, `pkg/httpmock/stub.go:35-199`, and `task_test.go:166-169`. | high |
| `13-security` | `studies/go-cli-study/reports/final/13-security.md` | Trust boundaries should be explicit for shell execution, credentials, paths, permissions, and untrusted input; evidence includes `internal/execext/exec.go:59-66`, `internal/permission/permission.go:44-108`, and `pkg/registry/transport.go:37-41`. | high |
| `14-performance` | `studies/go-cli-study/reports/final/14-performance.md` | Fast, reliable CLIs defer expensive work, stream bounded data, and cap concurrency; evidence includes `pkg/cmdutil/factory.go:27-42`, `age/internal/stream/stream.go:20`, and `gdu/pkg/analyze/parallel.go:13`. | high |
| `15-philosophy` | `studies/go-cli-study/reports/final/15-philosophy.md` | Coherent tools accept complexity deliberately, keep non-goals explicit, and use factory/interface/decorator patterns where they serve product constraints; evidence includes `VISION.md:97`, `pkg/cmd/factory/default.go:26-46`, and `internal/chezmoi/system.go:25-45`. | medium |

## Relevant Patterns

- **Thin command boundary over product-owned execute behavior:** The evidence repeatedly favors command layers that parse flags, construct dependencies, and delegate. Project structure evidence cites thin entrypoints such as `cmd/gh/main.go:6` and `main.go:44-45`, while command architecture evidence shows factories such as `pkg/cmd/install.go:132-145` and `pkg/cmd/issue/list/list.go:47-118`. For this sprint, the pressure is to keep `execute` CLI/status/flow wiring small while sprint-owned behavior handles task extraction, prompt rendering, run state, and validation.

- **Module-local ownership with enforced dependency direction:** Application-style CLIs commonly use `cmd/ + internal/` and unidirectional imports; the project-structure report cites `cmd/restic/main.go:37-114`, `internal/restic/repository.go:18`, and `cmd/root.go:112-127`. This supports sprint reasoning around keeping execute semantics in `internal/sprint`, runtime execution behind generic `platform/runtime`, and avoiding study-service reuse for planning-side task state.

- **Explicit composition root and narrow volatile interfaces:** Dependency-injection evidence favors centralized composition with constructor/factory injection over globals, citing `internal/cmd/config.go:362`, `pkg/cmdutil/factory.go:16-43`, `executor.go:22-24`, and `internal/backend/backend.go:19-90`. Execute reasoning should inspect whether fake runtime, filesystem/state-store, clock/time, and output boundaries are explicit enough for tests without adding broad interfaces for stable internal helpers.

- **Config precedence with stage-specific override visibility:** Configuration evidence converges on explicit ordering and changed-flag tracking, including `internal/cmd/config.go:2253-2287`, `internal/flags/flags.go:314-327`, `internal/global/global.go:139,147`, and `pkg/cli/environment.go:97-149`. This matters because execute has a stage-specific model override requirement; the handbook evidence supports reasoning about a single deterministic resolution path and diagnostics that show source of value without exposing secrets.

- **Runtime success must be paired with artifact validation:** Error-handling and testing evidence support treating outputs and validation as first-class completion gates. Mature error handling uses typed categories and `%w` wrapping (`fs/fserrors/error.go:22-29`, `go-task/errors/errors.go:47-50`, `internal/errors/fatal.go:10`), while testing evidence recommends fake runtime and golden/fixture checks (`task_test.go:166-169`, `pkg/httpmock/stub.go:35-199`). For execute, completion should be reasoned as a product state transition based on runtime result plus evidence or explicit diagnostic.

- **Atomic durable state and resumable task lifecycle:** State/context evidence favors explicit state plus context cancellation (`task.go:89`, `cmd/restic/cleanup.go:24-38`, `internal/session/session.go:12-23`), while concurrency evidence favors bounded coordination (`gdu/pkg/analyze/parallel.go:13`, `restic/internal/repository/repository.go:567`). This sprint needs reasoning about `.run-state.json` schema, stale running-task recovery, terminal states, and atomic writes.

- **Inspectable diagnostics with safe observability:** Logging evidence supports structured fields and separation of user output from logs (`internal/logging/logging.go:31-66`, `internal/slogs/keys.go:6-231`, `pkg/iostreams/iostreams.go:52-54`). Error evidence adds user-operational separation (`internal/model/flash.go:100-103`, `cmd/age/tui.go:37-54`). For execute, status and `execute.md` should point users to what failed, where state is, and what evidence is missing without leaking runtime secrets.

- **Workspace-safe and permission-aware execution boundary:** Security evidence highlights argument-array execution, permission dialogs, banned/allowlist models, credential scrubbing, private temp dirs, and schema validation (`internal/permission/permission.go:44-108`, `internal/llm/tools/bash.go:41-55`, `pkg/registry/transport.go:37-41`, `internal/config/json/validator.go:146`). For this sprint, the relevant pattern is not to copy those systems wholesale, but to make target repository paths explicit, workspace-safe, permission policy translation visible, and Git mutation excluded.

- **Behavior-focused test harnesses with fakes and golden fixtures:** Testing evidence favors table-driven cases, fake/mock infrastructure, command integration, and golden files (`internal/cmd/main_test.go:64-174`, `internal/chezmoi/mockpersistentstate.go:4-89`, `internal/backend/mock/backend.go:14-26`, `helm/internal/test/test.go:43`). Execute reasoning should plan fake-runtime tests for success-with-evidence, success-without-evidence, runtime failure, cancellation, malformed state, stale running tasks, and deterministic task IDs.

## Trade-Offs

| Trade-Off | Benefit | Cost | When It Matters |
| --- | --- | --- | --- |
| Thin CLI wrappers vs self-contained command handlers | Better testability and clearer ownership; supported by `pkg/cmd/install.go:132-145` and `pkg/cmdutil/factory.go:16-43` | More internal service/request types and delegation | Execute commands must integrate flow, status, runtime, and state without putting workflow semantics in CLI glue |
| Central composition root vs decentralized wiring | Dependencies and fake seams are traceable; supported by `internal/cmd/config.go:362` and `pkg/cmdutil/factory.go:16-43` | Risk of a large god object, as cautioned by `internal/cmd/config.go:193-291` | Execute needs runtime, config, workspace, state store, and output dependencies wired consistently |
| Interface injection vs concrete local helpers | Fakes are easier for runtime/filesystem/clock boundaries; evidence includes `internal/backend/backend.go:19-90` and `internal/chezmoi/system.go:25` | Too many interfaces obscure simple module-local behavior | Use at volatile/external boundaries, but keep deterministic plan parsing and task ID helpers concrete unless reasoning finds a testing need |
| Explicit config precedence vs ad-hoc override checks | Users can predict stage model selection; evidence includes `internal/flags/flags.go:314-327` and `restic/internal/global/global.go:139,147` | Requires source tracking and validation for unknown/empty stage keys | Execute supports global/default model plus stage-specific overrides, where stage-specific values win |
| Validation-gated completion vs runtime-result-only completion | Prevents false success and aligns with product requirements; testing evidence supports output validation through golden/fixture tests | Requires more task evidence rules and failure diagnostics | Execute must not mark tasks complete just because runtime exited successfully |
| Bounded concurrency vs maximum throughput | Prevents process/resource exhaustion; evidence includes `gdu/pkg/analyze/parallel.go:13`, `k9s/internal/pool.go:26-48`, and `go-task/setup.go:277-279` | Slower completion for many tasks; requires tuning | Running multiple plan tasks through runtime could otherwise spawn too much work or overwhelm target repo state |
| Persist every meaningful transition vs lower write frequency | Strong resumability and recovery after interruption | More file writes and schema/version care | `.run-state.json` must survive cancellation and stale running-task recovery |
| Structured diagnostics vs concise human output | Better automation and debugging; evidence includes `internal/slogs/keys.go:6-231` and `internal/ghcmd/cmd.go:281-301` | More output paths to maintain | Status and `execute.md` must be actionable for failed tasks while remaining readable |
| Permission/path safety checks vs frictionless execution | Reduces accidental writes and policy violations; evidence includes `internal/permission/permission.go:44-108` and `internal/config/json/validator.go:146` | Additional validation can block legitimate local workflows if too strict | Target repository writes must be explicit and workspace-safe; automatic Git mutation remains out of scope |

## Anti-Patterns And Warnings

- **Do not let `execute` become a monolithic `RunE`:** The command-architecture report flags long command handlers such as `cmd/root.go:49-183`, `runBackup()` at `cmd/restic/cmd_backup.go`, and `evaluateSequence()` at `cmd/evaluate_sequence_command.go:152` as logic leakage. Keep complex execute behavior out of command glue.

- **Do not create a global workflow/scheduler package for sprint task semantics:** The project-structure report supports module ownership and Go-enforced internal boundaries; the philosophy report warns against accumulated undocumented complexity. Execute task state and semantics should stay with sprint unless reasoning identifies a stable product-neutral helper.

- **Do not rely on package-level mutable config or runtime state:** Dependency and config reports warn about globals such as `pkg/yqlib/lib.go:13-21`, `fs/config.go:793`, and `cmd/root.go:121`. For execute, globals would make stage-specific models, fake runtime tests, and concurrent status/execute interactions harder to reason about.

- **Do not bypass IO abstractions in user-facing output or tests:** IO evidence flags direct `os.Stdout`/`os.Stderr` paths such as `cmd/ls/ls.go:42`, `dive/image/docker/cli.go:27`, and `pkg/yqlib/logger.go:18`. Execute status and summaries should be testable through injected writers or app-level output adapters.

- **Do not mark success without artifact/evidence validation:** Product requirements align with report evidence that runtime success is insufficient. Treat success-without-evidence as failed or explicitly diagnostic, not complete.

- **Do not spawn unbounded goroutines or runtime runs:** Concurrency evidence warns about unbounded goroutine spawning in `gh-cli/manager.go:198-203` and no timeout on waits in `gh-cli/manager.go:205`. Execute task loops need bounded parallelism and cleanup paths.

- **Do not persist unsafe raw runtime payloads or secrets by default:** Logging and security reports emphasize redaction and safe fields, citing `pkg/registry/transport.go:37-41`, `internal/options/secret_string.go:15-20`, and structured logging patterns. Execute diagnostics should preserve useful metadata without leaking environment or secret payloads.

- **Do not implement smoke, review, issue tracking, or Git mutation through execute:** The sprint requirements and scope documents explicitly exclude these behaviors. Evidence from security and philosophy reports shows that plugins/shell/Git mutation expand trust boundaries and should be deliberately scoped, not added incidentally.

## Examples Worth Inspecting

| Example | Path / Source | Why It Is Useful |
| --- | --- | --- |
| gh-cli factory dependency wiring | `studies/go-cli-study/reports/final/03-dependency-injection.md`; evidence `pkg/cmdutil/factory.go:16-43`, `pkg/cmd/factory/default.go:26-46` | Shows lazy dependency functions and command test seams that may inform execute wiring without global state. |
| Helm command/action split | `studies/go-cli-study/reports/final/02-command-architecture.md`; evidence `pkg/cmd/install.go:132-145`, `pkg/action/install.go:73-140` | Demonstrates thin command construction and action-owned business logic. |
| Go-task executor/options and structured concurrency | `studies/go-cli-study/reports/final/03-dependency-injection.md`, `08-concurrency.md`; evidence `executor.go:22-24`, `task.go:87`, `task.go:89` | Useful for reasoning about runtime task options, cancellation, and bounded task execution. |
| Restic durable state/error/IO boundaries | `studies/go-cli-study/reports/final/05-error-handling.md`, `06-io-abstraction.md`, `07-state-context.md`; evidence `internal/errors/fatal.go:10`, `internal/ui/terminal.go:10-36`, `internal/restic/lock.go:105` | Shows explicit fatal handling, terminal abstraction, and persistent lock/session style patterns. |
| Agent-style permission boundary in opencode | `studies/go-cli-study/reports/final/13-security.md`; evidence `internal/permission/permission.go:44-108`, `internal/llm/tools/bash.go:41-55` | Relevant for reasoning about runtime permissions, blocked operations, and user-facing permission diagnostics. |
| IOStreams test constructor | `studies/go-cli-study/reports/final/06-io-abstraction.md`; evidence `pkg/iostreams/iostreams.go:551-568` | Useful for designing output capture in status/execute tests. |
| Golden and fake test infrastructure | `studies/go-cli-study/reports/final/11-testing-strategy.md`; evidence `internal/cmd/main_test.go:64-174`, `helm/internal/test/test.go:43`, `internal/backend/mock/backend.go:14-26` | Supports planning fake-runtime, fixture, and summary-output tests. |
| Explicit config precedence patterns | `studies/go-cli-study/reports/final/04-configuration-management.md`; evidence `internal/cmd/config.go:2253-2287`, `internal/flags/flags.go:314-327` | Useful when reasoning through global/default model and execute-stage override behavior. |
| Bounded concurrency/resource management | `studies/go-cli-study/reports/final/08-concurrency.md`, `14-performance.md`; evidence `gdu/pkg/analyze/parallel.go:13`, `restic/internal/archiver/file_saver.go:56-58` | Useful for execute task loop limits and cancellation behavior. |

## Design Pressures

- Execute must preserve product/module ownership: `internal/sprint` owns task extraction, prompt rendering, `.run-state.json`, `execute.md`, status, and flow through execute, while runtime execution remains generic.

- Runtime execution must use the existing generic runtime/agentwrap boundary rather than product-specific process management.

- Task IDs must be deterministic and trace back to validated `plan.md` entries, which pushes toward stable parsing, normalized task identity inputs, and tests that exercise equivalent plan content.

- Completion semantics must distinguish runtime success, validation success, explicit diagnostic completion, failure, cancellation, and stale running recovery.

- State durability must be strong enough for interruption and resume, implying schema versioning, atomic writes, malformed-state rejection, and conservative stale-task recovery.

- Target repository boundaries must be explicit and safe. Execute needs to write implementation changes only inside an approved target directory and must not add automatic Git mutation.

- Status and `execute.md` need to be concise enough for humans while carrying enough structured diagnostics for recovery.

- Config/model selection must remain predictable across planning and execute stages, with stage-specific overrides winning over global defaults.

- Testing pressure is high because real runtime behavior is nondeterministic. Fake runtime, fixture state files, and output/golden checks are central to confidence.

- Performance pressure is less about raw speed than bounded work: avoid expensive scans, unbounded task fan-out, excessive state reads, and loading large artifacts unnecessarily.

## Open Questions For Reasoning

- What exact `plan.md` task syntax is executable for this sprint, and which fields participate in deterministic task ID generation?

- What evidence qualifies as machine-validatable task completion, and when is an explicit diagnostic acceptable instead of direct evidence?

- Should execute run one task at a time by default, or allow bounded parallel execution? If bounded, what safety constraints prevent conflicting edits in the same target repository?

- How should stale `running` tasks in `.run-state.json` be recovered on resume: always pending, failed with diagnostic, or policy-dependent?

- Which status fields are required in text output versus JSON output so humans and scripts can both inspect execute progress?

- How much agentwrap run metadata should be copied into `.run-state.json` versus referenced through run IDs or summarized diagnostics?

- Should `execute.md` be written after every task transition, only at terminal states, or only when explicitly requested after execution?

- How should target repository path configuration be represented so it remains workspace-safe but supports the project index's target implementation directory?

- Which configuration diagnostics should show the selected model source for execute without exposing environment or runtime secrets?

- How should cancellation map into task state and process exit: cancelled, retryable pending, failed, or mixed partial completion?

## Evidence Pointers

- `studies/go-cli-study/reports/final/01-project-structure.md`: Inspect thin CLI entrypoint, `internal/` boundary, unidirectional dependency, and global-state cautions.

- `studies/go-cli-study/reports/final/02-command-architecture.md`: Inspect factory command creation, options structs, lifecycle hooks, `cmd.Run()` wrappers, and long `RunE` anti-patterns.

- `studies/go-cli-study/reports/final/03-dependency-injection.md`: Inspect composition root, constructor injection, functional options, interface boundaries, lazy initialization, and globals-heavy warnings.

- `studies/go-cli-study/reports/final/04-configuration-management.md`: Inspect flag/config/env precedence, changed-flag tracking, post-load validation, and direct environment bypass warnings.

- `studies/go-cli-study/reports/final/05-error-handling.md`: Inspect wrapped errors, sentinel/typed errors, fatal/recoverable distinctions, user-facing hints, and exit-code mapping.

- `studies/go-cli-study/reports/final/06-io-abstraction.md`: Inspect IOStreams/test constructors, filesystem/backend interfaces, mock terminals, and hardcoded `os.Stdout`/`os.Stderr` warnings.

- `studies/go-cli-study/reports/final/07-state-context.md`: Inspect signal-context wiring, errgroup cancellation, session modeling, context-as-DI cautions, and stale context/background anti-patterns.

- `studies/go-cli-study/reports/final/08-concurrency.md`: Inspect errgroup, semaphore, stop-channel, timeout-on-wait, and fire-and-forget goroutine warnings.

- `studies/go-cli-study/reports/final/09-terminal-ux.md`: Inspect non-TTY fallback, progress/cancellation UX, status-line patterns, and prompt blocking warnings.

- `studies/go-cli-study/reports/final/10-logging-observability.md`: Inspect structured logging, field naming, debug controls, stdout/stderr separation, and log redaction practices.

- `studies/go-cli-study/reports/final/11-testing-strategy.md`: Inspect testscript, golden files, fake/mock packages, behavior-focused assertions, and global-state test isolation warnings.

- `studies/go-cli-study/reports/final/13-security.md`: Inspect path/trust boundaries, permission request flow, credential redaction, shell-safe execution, and dangerous default warnings.

- `studies/go-cli-study/reports/final/14-performance.md`: Inspect lazy initialization, streaming, bounded concurrency, profiling hooks, and unbounded memory/goroutine cautions.

- `studies/go-cli-study/reports/final/15-philosophy.md`: Inspect explicit non-goals, factory/interface/decorator rationale, plugin/security trade-offs, and complexity warnings.

## Handoff To Reasoning

- Use this handbook as evidence input.
- Validate whether the observed patterns fit this project's constraints.
- Do not copy external patterns without sprint-specific reasoning.
