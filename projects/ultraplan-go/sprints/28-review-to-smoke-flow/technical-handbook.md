# Sprint Technical Handbook: Integrated Review-to-Smoke Verification Flow

> Project: `ultraplan-go`
> Sprint: `28-review-to-smoke-flow`
> Source: `projects/ultraplan-go/sprints/28-review-to-smoke-flow/sprint-index.md`
> **Inputs Used:** `projects/ultraplan-go/project-index.md`, `projects/ultraplan-go/sprints/28-review-to-smoke-flow/sprint-index.md`, `projects/ultraplan-go/sprints/28-review-to-smoke-flow/requirements.md`, `projects/ultraplan-go/docs/ARCHITECTURE.md`, `projects/ultraplan-go/docs/PRD.md`, `projects/ultraplan-go/docs/TRD.md`, `studies/go-cli-study/reports/final/01-project-structure.md`, `studies/go-cli-study/reports/final/02-command-architecture.md`, `studies/go-cli-study/reports/final/03-dependency-injection.md`, `studies/go-cli-study/reports/final/04-configuration-management.md`, `studies/go-cli-study/reports/final/05-error-handling.md`, `studies/go-cli-study/reports/final/06-io-abstraction.md`, `studies/go-cli-study/reports/final/07-state-context.md`, `studies/go-cli-study/reports/final/08-concurrency.md`, `studies/go-cli-study/reports/final/09-terminal-ux.md`, `studies/go-cli-study/reports/final/10-logging-observability.md`, `studies/go-cli-study/reports/final/11-testing-strategy.md`, `studies/go-cli-study/reports/final/13-security.md`, `studies/go-cli-study/reports/final/15-philosophy.md`

This handbook distills the studies and reports selected by `sprint-index.md` for sprint reasoning. It does not decide architecture or implementation. Project documents and sprint requirements were used only to identify relevant pressures; the patterns and trade-offs below are grounded in the selected evidence reports.

## Selected Studies And Reports

| Study / Report | Path | Relevant Finding | Confidence |
| --- | --- | --- | --- |
| `01-project-structure` | `studies/go-cli-study/reports/final/01-project-structure.md` | Mature applications keep entrypoints thin and imports unidirectional; chezmoi and yq provide concrete CLI-to-core examples (`chezmoi/main.go:26-34`, `yq/cmd/root.go:9`, `yq/cmd/evaluate_sequence_command.go:7`). | high |
| `02-command-architecture` | `studies/go-cli-study/reports/final/02-command-architecture.md` | Factory-built commands delegate to reusable operations, while large `RunE` functions and TUI-only command systems reduce testability and scriptability (`gh-cli/pkg/cmdutil/factory.go:16-43`, `helm/pkg/cmd/install.go:132-145`, `opencode/cmd/root.go:49-183`, `k9s/internal/view/command.go:176`). | high |
| `03-dependency-injection` | `studies/go-cli-study/reports/final/03-dependency-injection.md` | Explicit composition roots and boundary interfaces support substitutes without a DI framework; global state is the principal source of hidden coupling (`gh-cli/pkg/cmd/factory/default.go:26-46`, `restic/internal/backend/backend.go:19-90`, `rclone/fs/config.go:793`). | high |
| `04-configuration-management` | `studies/go-cli-study/reports/final/04-configuration-management.md` | Effective configuration needs explicit precedence, explicit flag-change tracking, and validation after all sources merge (`chezmoi/internal/cmd/config.go:2253-2287`, `restic/internal/global/global.go:139,147`, `k9s/internal/config/k9s.go:423-451`). | high |
| `05-error-handling` | `studies/go-cli-study/reports/final/05-error-handling.md` | Wrapped, typed or sentinel errors preserve machine-readable behavior while a separate rendering path can add user recovery hints and stable exit codes (`go-task/errors/errors_task.go:13-32`, `age/cmd/age/tui.go:37-54`, `gh-cli/internal/ghcmd/cmd.go:44-49`). | high |
| `06-io-abstraction` | `studies/go-cli-study/reports/final/06-io-abstraction.md` | Injected streams and volatile-boundary interfaces make output and process behavior testable without real terminals or external systems (`gh-cli/pkg/iostreams/iostreams.go:551-568`, `restic/internal/ui/terminal.go:10-36`, `lazygit/pkg/commands/oscommands/cmd_obj_runner.go:18-23`). | high |
| `07-state-context` | `studies/go-cli-study/reports/final/07-state-context.md` | Long-running work benefits from root-context propagation, signal cancellation, explicit state ownership, and a distinct cleanup policy (`helm/pkg/cmd/install.go:333-347`, `go-task/task.go:89`, `restic/internal/restic/lock.go:290-305`). | high |
| `08-concurrency` | `studies/go-cli-study/reports/final/08-concurrency.md` | Structured concurrency, bounded fan-out, localized goroutine launch sites, and timeout-bounded shutdown are safer than scattered or fire-and-forget work (`go-task/task.go:87`, `k9s/internal/pool.go:21,30,37`, `opencode/cmd/root.go:261-279`). | high |
| `09-terminal-ux` | `studies/go-cli-study/reports/final/09-terminal-ux.md` | Interactive surfaces need truthful progress and interruptibility, but automation still needs a non-TTY path and independently scriptable operations (`chezmoi/internal/cmd/prompt.go:20-256`, `lazygit/pkg/tasks/tasks.go:144`, `gh-cli/pkg/iostreams/iostreams.go:116-130`). | high |
| `10-logging-observability` | `studies/go-cli-study/reports/final/10-logging-observability.md` | Structured diagnostics, runtime debug controls, and separation of result output from logs improve operator recovery (`helm/internal/logging/logging.go:31-71`, `k9s/internal/slogs/keys.go:6-231`, `gh-cli/pkg/iostreams/iostreams.go:52-54`). | high |
| `11-testing-strategy` | `studies/go-cli-study/reports/final/11-testing-strategy.md` | Strong CLI suites combine behavior-focused command tests, reusable fakes, and golden fixtures for stable user-visible output (`gh-cli/acceptance/acceptance_test.go:26-29`, `lazygit/pkg/commands/oscommands/fake_cmd_obj_runner.go:17-26`, `helm/internal/test/test.go:43`). | high |
| `13-security` | `studies/go-cli-study/reports/final/13-security.md` | External execution should use explicit executable/argument arrays, explicit trust boundaries, path validation, and redaction (`lazygit/cmd_obj_builder.go:38`, `gh-cli/pkg/iostreams/iostreams.go:228`, `restic/internal/options/secret_string.go:15-20`). | high |
| `15-philosophy` | `studies/go-cli-study/reports/final/15-philosophy.md` | Coherent tools accept complexity only for a named pressure and preserve narrow interfaces and explicit non-goals (`gh-cli/pkg/cmd/factory/default.go:26-46`, `restic/internal/backend/backend.go:19-90`, `lazygit/VISION.md:97`). | medium |

## Relevant Patterns

- **Thin surfaces over a shared operation:** The consistent application pattern is for command code to parse inputs, call a reusable operation, and render a result. Helm's command delegates to an action (`helm/pkg/cmd/install.go:132-145`), gh-cli passes a dependency factory into command constructors (`gh-cli/pkg/cmdutil/factory.go:16-43`), and gdu exposes alternate UI implementations through one interface (`gdu/cmd/gdu/app/app.go:30-49`). For this sprint, the pattern is relevant to keeping CLI text, JSON, and TUI behavior aligned without making any one surface the workflow owner.

- **One explicit composition root with narrow volatile-boundary seams:** The study found manual composition rather than DI frameworks across all examined projects. Gh-cli's lazy factory (`gh-cli/pkg/cmd/factory/default.go:26-46`) and restic's backend boundary (`restic/internal/backend/backend.go:19-90`) permit deterministic substitutes while keeping construction traceable. Review runtime and smoke-process dependencies create similar volatility, but evidence also warns that a central object can become a god object (`chezmoi/internal/cmd/config.go:193-291`).

- **Explicit, inspectable lifecycle state:** Central app/config/executor structs give long operations clear state ownership (`helm/pkg/action/action.go:119-146`, `opencode/internal/app/app.go:25-40`, `go-task/executor.go:27-84`). Restic's lock lifecycle additionally separates cleanup from cancelled work (`restic/internal/restic/lock.go:290-305`). The sprint's resumability, stale propagation, preservation of last complete artifacts, and recovery guidance all pressure reasoning toward explicit states and transitions rather than deriving truth from transient UI activity.

- **Cancellation from entrypoint through every blocking boundary:** Signal-to-context wiring is visible in Helm (`helm/pkg/cmd/install.go:333-347`), while go-task couples parallel work through `errgroup.WithContext` (`go-task/task.go:89`). The reports specifically warn that isolated `context.Background()` calls make nominally cancellable operations ignore their caller (`chezmoi/internal/cmd/templatefuncs.go:215`, `go-task/task.go:341`). Review, harness execution, progress delivery, waits, and cleanup each need a deliberate context policy.

- **Bounded and localized concurrency:** Restic and go-task localize fan-out behind `errgroup` (`restic/internal/repository/repository.go:567`, `go-task/task.go:87`); k9s bounds work with a semaphore-backed pool (`k9s/internal/pool.go:21,30,37`); opencode bounds shutdown waits (`opencode/cmd/root.go:261-279`). This is relevant if review requests or progress streams run concurrently, but the evidence does not imply that sequential review-before-smoke orchestration itself should be parallelized.

- **Typed failures plus surface-specific recovery rendering:** Go-task carries corrective data in `TaskNotFoundError` (`go-task/errors/errors_task.go:13-32`), gh-cli maps known conditions to stable exit codes (`gh-cli/internal/ghcmd/cmd.go:44-49`), and age adds user hints without discarding the underlying error (`age/cmd/age/tui.go:37-54`). The sprint has several distinct non-success outcomes, including stale, blocked, malformed evidence, timeout, cancellation, and failed verdict; evidence supports preserving those distinctions for programmatic consumers while rendering a useful next action for people.

- **Output and external effects injected at boundaries:** Gh-cli's test IO constructor provides in-memory stdin/stdout/stderr (`gh-cli/pkg/iostreams/iostreams.go:551-568`), restic abstracts terminal behavior (`restic/internal/ui/terminal.go:10-36`), and lazygit records external command calls behind `ICmdObjRunner` (`lazygit/pkg/commands/oscommands/cmd_obj_runner.go:18-23`). These seams are directly relevant to parity tests and fake smoke-harness execution without requiring the real terminal, OpenCode, or harness.

- **Explicit configuration precedence and post-merge validation:** Chezmoi restores explicitly changed flags after config loading (`chezmoi/internal/cmd/config.go:2253-2287`), restic records whether a flag was changed (`restic/internal/global/global.go:139,147`), and k9s validates the merged configuration (`k9s/internal/config/k9s.go:423-451`). Harness command discovery, timeout, runtime selection, and diagnostic overrides need similarly visible source and validation behavior so defaults cannot silently override operator intent.

- **Trust-boundary execution using executable plus argv:** The security report identifies explicit argument arrays as the baseline for external commands (`lazygit/cmd_obj_builder.go:38`) and points to safe executable resolution in gh-cli (`gh-cli/pkg/surveyext/editor_manual.go:23`, `gh-cli/pkg/iostreams/iostreams.go:228`). This matters because the sprint invokes a cataloged harness: command prose or Markdown should not become an executable shell program.

- **Structured status separate from diagnostics:** Helm routes structured logs to stderr (`helm/internal/logging/logging.go:31-71`), k9s standardizes diagnostic keys (`k9s/internal/slogs/keys.go:6-231`), and gh-cli centralizes streams (`gh-cli/pkg/iostreams/iostreams.go:52-54`). A typed status result can therefore remain the source for text, JSON, and TUI while debug diagnostics stay out of stable machine output.

- **Behavior-focused tests at multiple boundaries:** Testscript scenarios exercise full command behavior (`chezmoi/internal/cmd/main_test.go:64-174`, `gh-cli/acceptance/acceptance_test.go:26-29`), fake runners isolate effects (`lazygit/pkg/commands/oscommands/fake_cmd_obj_runner.go:17-26`), and golden helpers detect output drift (`helm/internal/test/test.go:43`). This combination maps well to domain transition tests, fake runtime/harness tests, and parity fixtures, while gated tests can retain a smaller role for real integrations.

## Trade-Offs

| Trade-Off | Benefit | Cost | When It Matters |
| --- | --- | --- | --- |
| Shared typed operation vs surface-owned behavior | One semantic result can drive CLI text, JSON, and TUI, following Helm/gh-cli delegation (`helm/pkg/cmd/install.go:132-145`, `gh-cli/pkg/cmdutil/factory.go:16-43`). | Adds an application boundary and result types; an oversized shared factory can become a god object (`chezmoi/internal/cmd/config.go:193-291`). | When parity and future reuse outweigh the convenience of implementing directly in each handler. |
| Explicit interfaces vs concrete collaborators | Runtime/process/IO fakes become easy and boundary contracts are visible (`restic/internal/backend/backend.go:19-90`, `gh-cli/pkg/iostreams/iostreams.go:551-568`). | Interface maintenance and indirection; large interfaces burden every fake (`restic/internal/fs/interface.go:10-31`). | At volatile external boundaries or where normal tests must avoid real services, not for every internal helper. |
| Centralized state vs recomputation from artifacts | Gives one inspectable place for lifecycle, attempts, verdicts, and next actions (`go-task/executor.go:27-84`, `opencode/internal/app/app.go:25-40`). | State can drift from editable artifacts or grow into an over-coupled aggregate. | When reasoning defines which values are durable facts, fingerprints, or deterministic derivations and how reconciliation works. |
| Preserve last complete artifact vs expose partial progress | Atomic preservation supports recovery and stable evidence; progressive feedback prevents long operations from looking frozen (`lazygit/pkg/tasks/tasks.go:31-435`, `restic/internal/ui/termstatus/status.go:197-205`). | Requires separate transient progress/state from the canonical completed artifact and careful interruption cleanup. | During review or smoke cancellation, timeout, malformed output, and reruns. |
| Fail-fast structured concurrency vs collect independent findings | `errgroup.WithContext` promptly cancels dependent siblings and simplifies error propagation (`go-task/task.go:87`, `restic/internal/repository/repository.go:567`). | Early cancellation can hide useful independent findings; `WaitGroup` collection requires more manual classification. | If multiple review checks are independent, reasoning must decide whether one failure invalidates siblings or whether complete evidence is more valuable. |
| Narrow rerun vs containing-scope confidence | Focused execution lowers time and cost, consistent with deliberate bounded work. | A narrow pass may not establish the broader invariant; the study's real-fixture evidence also shows realism costs more time and setup (`dive/dive/image/docker/testing.go:12-34`). | When deciding how a test-level or suite-level rerun contributes to the canonical smoke verdict. |
| Rich typed errors vs a small error surface | Structured fields and `errors.Is`/`errors.As` enable stable recovery and exit mapping (`rclone/fs/fserrors/error.go:22-29`, `gh-cli/internal/ghcmd/cmd.go:44-49`). | Too many types or behavioral interfaces create taxonomy overhead and classification ambiguity. | When blocked, stale, cancelled, timed out, malformed, or failed outcomes need distinct machine behavior. |
| Human progress vs stable automation output | Progress improves confidence for operations longer than a few seconds (`chezmoi/internal/cmd/readhttpresponse.go:25-108`, `rclone/cmd/progress.go:24-71`). | ANSI updates and concurrent messages can corrupt JSON or piped output. | When CLI text and TUI show live work while JSON must remain deterministic and non-interactive. |
| Configuration convenience vs explicit trust | Layered config reduces repetitive flags and supports local harness setup (`go-task/internal/flags/flags.go:314-327`). | Precedence, migration, environment leakage, and hidden defaults increase complexity; direct `os.Getenv` bypasses are inconsistent (`opencode/internal/config/config.go:163`). | When harness executable, argv, timeout, environment, and diagnostic overrides can come from more than one source. |
| Golden parity tests vs semantic assertions | Golden files expose complete output drift in reviewable diffs (`helm/internal/test/test.go:43`, `go-task/task_test.go:166-169`). | Fixtures require deliberate updates and can create noisy diffs or bless defects. | For stable text/JSON/TUI projections; domain verdict rules still need focused semantic tests. |
| Sequential flow vs concurrency | Sequential ordering is simpler and naturally preserves review-before-smoke; bounded concurrency can speed independent review work. | Parallelism adds cancellation, ordering, race, and partial-state complexity (`gh-cli/pkg/cmd/extension/manager.go:196-206`). | Only where independent work is demonstrably expensive; smoke gating itself remains an ordering constraint from sprint scope. |

## Anti-Patterns And Warnings

- **Do not put verification semantics in CLI or TUI handlers:** Large command functions such as `opencode/cmd/root.go:49-183` and TUI-only dispatch such as `k9s/internal/view/command.go:176` make behavior harder to reuse and script. Surface code should not become a second verdict, freshness, or recovery engine.

- **Do not create a global mutable verification registry:** Rclone's global config fallback hides missing context (`rclone/fs/config.go:793`), and singleton caches pollute tests (`rclone/fs/cache/cache.go:16-21`). Package-level review runtime, process runner, current sprint, or alternate TUI state would create the same coupling.

- **Do not use context as an untyped service locator:** K9s relies on runtime assertions from context values (`k9s/internal/keys.go:10-38`), which can panic when wiring is absent. Context should carry cancellation, deadlines, and request-scoped metadata; core collaborators remain explicit.

- **Do not sever cancellation with background contexts:** Chezmoi's background template calls (`chezmoi/internal/cmd/templatefuncs.go:215`) and go-task's detached deferred work (`go-task/task.go:341`) demonstrate how cleanup or child work can outlive caller intent. Any deliberate detached cleanup needs a bounded, documented lifetime.

- **Do not launch fire-and-forget progress or process goroutines:** Dive's unwaited notification goroutine (`dive/cmd/dive/cli/internal/command/adapter/resolver.go:70`) and unbounded waits in gh-cli (`gh-cli/pkg/cmd/extension/manager.go:196-206`) expose leak and hang risks. Every spawned activity needs ownership, termination, and bounded waiting.

- **Do not infer success from process exit alone:** The selected reports emphasize artifact/behavior validation over implementation signals; testscript verifies observable CLI behavior (`chezmoi/internal/cmd/main_test.go:64-174`), and typed errors preserve classification. A zero exit from review runtime or harness cannot substitute for valid, current evidence.

- **Do not construct shell strings from catalog or Markdown input:** Explicit argv prevents metacharacters from becoming code (`lazygit/cmd_obj_builder.go:38`). The security report flags string concatenation, unrestricted shells, and missing external-command timeouts as brittle trust boundaries.

- **Do not let debug output corrupt stable output:** Fzf writes stats to stdout (`fzf/src/core.go:325`), identified as an observability anti-pattern. JSON output should not receive spinners, logs, ANSI, or unrestricted harness output.

- **Do not bypass injected streams or runners:** Direct terminal output in otherwise abstracted systems creates untestable paths (`chezmoi/internal/cmd/templatefuncs.go:296`, `rclone/cmd/ls/ls.go:42`). CLI, TUI, and fake-harness tests need all relevant effects to pass through their declared boundary.

- **Do not collapse blocked, cancelled, stale, malformed, and failed into one string error:** Plain strings and magic exit codes limit programmatic recovery (`mitchellh-cli/command.go:10-11`, `mitchellh-cli/cli.go:205-206`). Conversely, avoid an elaborate hierarchy where a small typed result or sentinel is sufficient.

- **Do not overfit tests to internal representation:** K9s's exact hint-count test (`k9s/internal/view/pod_test.go:23`) breaks on harmless UI changes. Parity tests should compare required semantics, with goldens reserved for intentionally stable projections.

- **Do not add abstraction without a sprint pressure:** The philosophy report warns about accumulated god structs and configuration sprawl (`lazygit/pkg/gui/gui.go`, `lazygit/Config.md`) and praises deliberate rejection of complexity (`lazygit/VISION.md:97`). This sprint needs external-boundary seams and shared use cases, not a general workflow engine or plugin system.

## Examples Worth Inspecting

| Example | Path / Source | Why It Is Useful |
| --- | --- | --- |
| Thin command delegating to an action | `helm/pkg/cmd/install.go:132-145` via `studies/go-cli-study/reports/final/02-command-architecture.md` | Shows where command parsing can stop and reusable behavior can begin. |
| Lazy dependency factory | `gh-cli/pkg/cmd/factory/default.go:26-46` and `gh-cli/pkg/cmdutil/factory.go:16-43` via reports `02`, `03`, and `15` | Useful for reasoning about one CLI/TUI composition path without eagerly constructing every runtime dependency. |
| Test IO bundle | `gh-cli/pkg/iostreams/iostreams.go:551-568` via `studies/go-cli-study/reports/final/06-io-abstraction.md` | Demonstrates production and in-memory streams with the same command path. |
| Fake external command runner | `lazygit/pkg/commands/oscommands/fake_cmd_obj_runner.go:17-26` via `studies/go-cli-study/reports/final/11-testing-strategy.md` | Relevant to proving executable, argv, cwd, environment, cancellation, and timeout behavior without the real harness. |
| Signal and context cancellation | `helm/pkg/cmd/install.go:333-347` via reports `07` and `09` | Shows signal cancellation reaching a long-running operation. |
| Structured fan-out | `go-task/task.go:87-89` via reports `07` and `08` | Provides a compact example of sibling cancellation and collected errors for independent bounded work. |
| Bounded worker pool | `k9s/internal/pool.go:21,30,37` via `studies/go-cli-study/reports/final/08-concurrency.md` | Useful if review fan-out needs an explicit concurrency cap. |
| Bounded shutdown wait | `opencode/cmd/root.go:261-279` via `studies/go-cli-study/reports/final/08-concurrency.md` | Shows why cancellation alone is insufficient when goroutines or subprocess cleanup can hang. |
| Typed recoverable error with suggestion | `go-task/errors/errors_task.go:13-32` via `studies/go-cli-study/reports/final/05-error-handling.md` | Demonstrates carrying next-action data without parsing an error string. |
| User-facing error hints | `age/cmd/age/tui.go:37-54` via `studies/go-cli-study/reports/final/05-error-handling.md` | Shows separation between an underlying failure and actionable presentation. |
| Explicit config precedence | `chezmoi/internal/cmd/config.go:2253-2287` and `restic/internal/global/global.go:139,147` via `studies/go-cli-study/reports/final/04-configuration-management.md` | Useful for diagnostic overrides and for distinguishing explicit flags from defaults. |
| Explicit argv execution | `lazygit/cmd_obj_builder.go:38` and `gh-cli/pkg/iostreams/iostreams.go:228` via `studies/go-cli-study/reports/final/13-security.md` | Establishes a concrete non-shell process boundary. |
| Secret-safe value | `restic/internal/options/secret_string.go:15-20` via `studies/go-cli-study/reports/final/13-security.md` | Useful for evaluating whether environment/config diagnostics need type-level redaction. |
| CLI acceptance scenarios | `gh-cli/acceptance/acceptance_test.go:26-29` via `studies/go-cli-study/reports/final/11-testing-strategy.md` | Shows multi-step, behavior-level command verification. |
| Golden output helper | `helm/internal/test/test.go:43` via `studies/go-cli-study/reports/final/11-testing-strategy.md` | Useful for stable text and JSON projection tests with intentional update review. |
| Progressive cancellable output | `lazygit/pkg/tasks/tasks.go:31-435` and `lazygit/pkg/tasks/tasks.go:144` via `studies/go-cli-study/reports/final/09-terminal-ux.md` | Demonstrates separating long-running progress from final durable results. |

## Design Pressures

- The workflow must expose one coherent truth across CLI text, stable JSON, and TUI while retaining interface-specific rendering and interaction.
- Review must be current and acceptable before default smoke execution, but diagnostic override behavior must remain explicit, visible, and testable.
- Editable Markdown artifacts and strict versioned JSON state can diverge after external edits, so freshness, reconciliation, and malformed-evidence behavior need precise ownership.
- A stale review propagates to smoke, while focused reruns should avoid unnecessary work without letting narrow evidence erase containing-scope failures.
- Cancellation and timeout must stop runtime/process work and leave the system resumable, while bounded cleanup may need to outlive the cancelled work context.
- The last complete `review.md` and `smoke.md` must remain trustworthy during partial or failed replacement, creating a distinction between transient run state and canonical completed evidence.
- External harness evidence remains outside the workspace, so status must be informative through stable links and identifiers without copying raw output or silently treating missing evidence as success.
- Open harness issues affect smoke assessment, but the product must not expand into issue ownership, assignment, or remediation.
- The overall assessment is derived rather than a third artifact; its inputs, precedence, contradictions, and stale/blocked behavior must therefore be deterministic.
- Runtime and process boundaries need enough metadata for recovery and diagnostics without learning sprint verdict semantics or leaking secrets and unrestricted output.
- Configuration may come from defaults, workspace config, environment, and flags; explicit operator choices must remain distinguishable from defaults.
- Normal tests must prove cancellation, parity, freshness, process invocation, and recovery using fakes, while real runtime/harness evidence remains gated.
- The sprint integrates existing execute, review, and smoke behavior rather than replacing it, so new abstractions must justify migration and regression risk.
- Verification is read-only with respect to product source, tests, governed inputs, and Git state; external execution must preserve that trust boundary.

## Open Questions For Reasoning

- What is the smallest canonical verification result that can represent stage execution, verdict, freshness, artifact/run references, relevant issues, and next action without duplicating durable state?
- Which fields are persisted facts, which are fingerprints, and which are derived at read time from current artifacts and harness evidence?
- What exact input manifest feeds review freshness, smoke freshness, and overall assessment, and how are path ordering, timestamps, missing files, and externally edited files normalized?
- At what point is an interrupted run recorded, and how is that transient attempt distinguished from the last complete canonical artifact and verdict?
- Which error/result categories require distinct CLI exit behavior, and which belong only in the typed status payload or human recovery guidance?
- Does explicit diagnostic override permit smoke execution only, or can it influence the resulting assessment; where is the override and its rationale recorded?
- How does a focused review rerun merge with or supersede prior findings while still producing one current review artifact?
- What containing-suite evidence remains mandatory after a focused smoke level, suite, or test rerun, and how is its freshness compared with the focused run?
- How are relevant harness issues selected, resolved, and fingerprinted without importing issue-management semantics into the product?
- Should independent review checks fail fast or complete to maximize findings, and what bounded concurrency is justified by actual runtime cost?
- Which context should cleanup use after cancellation, what is its maximum duration, and what happens when process-tree cleanup itself fails?
- What process-runner result is sufficiently generic to report exit, signal, timeout, cancellation, bounded stdout/stderr, and cleanup without embedding smoke semantics?
- How are configured executable and argv values validated and surfaced safely, including PATH lookup, cwd containment, environment allowlisting, and redaction?
- What typed app use-case boundary lets the TUI receive progress and cancellation controls while CLI JSON remains non-streaming and deterministic?
- Which parity assertions compare semantic fields across CLI/JSON/TUI, and which complete renderings deserve golden coverage?
- How should stale, blocked, malformed, failed, cancelled, and timed-out states map to a deterministic overall assessment and required next action?
- What schema rejection or migration behavior is needed for older flow state once review/smoke freshness metadata is added?
- Which existing Sprint 26 and Sprint 27 seams can be composed directly, and where would a new abstraction merely duplicate established behavior?

## Evidence Pointers

- `studies/go-cli-study/reports/final/01-project-structure.md`: Inspect “Pattern 1: Thin CLI Entry Point,” “Pattern 4: Unidirectional Dependency Flow,” and the global-state cautions.
- `studies/go-cli-study/reports/final/02-command-architecture.md`: Inspect “Thin-Delegate Pattern,” factory command creation, shared lifecycle wrappers, and the bifurcated CLI/TUI warning.
- `studies/go-cli-study/reports/final/03-dependency-injection.md`: Inspect centralized composition roots, constructor injection, volatile-boundary interfaces, and global-state failure modes.
- `studies/go-cli-study/reports/final/04-configuration-management.md`: Inspect explicit precedence, changed-flag tracking, post-merge validation, and direct-environment bypass warnings.
- `studies/go-cli-study/reports/final/05-error-handling.md`: Inspect typed errors, sentinels, user/operational separation, exit mapping, and hint-based recovery.
- `studies/go-cli-study/reports/final/06-io-abstraction.md`: Inspect IOStreams, fake command/terminal seams, concurrent-safe buffers, and direct `os.*` leakage.
- `studies/go-cli-study/reports/final/07-state-context.md`: Inspect signal-context wiring, centralized application state, cleanup contexts, and context-as-service-locator cautions.
- `studies/go-cli-study/reports/final/08-concurrency.md`: Inspect errgroup fan-out, bounded pools, localized launch sites, shutdown timeouts, and fire-and-forget warnings.
- `studies/go-cli-study/reports/final/09-terminal-ux.md`: Inspect progressive interruptible output, non-TTY fallback, signal-safe exit, and CLI-first automation trade-offs.
- `studies/go-cli-study/reports/final/10-logging-observability.md`: Inspect structured fields, stdout/stderr separation, runtime debug control, and output-abstraction bypasses.
- `studies/go-cli-study/reports/final/11-testing-strategy.md`: Inspect behavior-level acceptance tests, fake runners, golden files, test isolation, and implementation-detail anti-patterns.
- `studies/go-cli-study/reports/final/13-security.md`: Inspect explicit argv execution, safe executable lookup, redaction, path/config validation, and command timeout cautions.
- `studies/go-cli-study/reports/final/15-philosophy.md`: Inspect deliberate complexity, explicit non-goals, factory/interface patterns, and warnings about god structs and configuration sprawl.

## Handoff To Reasoning

- Use this handbook as evidence input.
- Validate whether the observed patterns fit this project's constraints.
- Resolve the open questions in area-specific reasoning or final sprint reasoning.
- Treat project documents and sprint requirements as scope constraints, not additional handbook evidence.
- Do not copy external patterns without sprint-specific reasoning.
