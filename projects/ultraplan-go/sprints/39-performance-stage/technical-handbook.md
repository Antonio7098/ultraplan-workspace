# Sprint Technical Handbook: Requirements-Driven Performance Stage

> Project: `ultraplan-go`
> Sprint: `39-performance-stage`
> Source: `projects/ultraplan-go/sprints/39-performance-stage/sprint-index.md`
> **Inputs Used:** `projects/ultraplan-go/project-index.md`, `projects/ultraplan-go/sprints/39-performance-stage/sprint-index.md`, `projects/ultraplan-go/sprints/39-performance-stage/requirements.md`, `studies/go-cli-study/reports/final/01-project-structure.md`, `studies/go-cli-study/reports/final/02-command-architecture.md`, `studies/go-cli-study/reports/final/03-dependency-injection.md`, `studies/go-cli-study/reports/final/04-configuration-management.md`, `studies/go-cli-study/reports/final/05-error-handling.md`, `studies/go-cli-study/reports/final/06-io-abstraction.md`, `studies/go-cli-study/reports/final/07-state-context.md`, `studies/go-cli-study/reports/final/08-concurrency.md`, `studies/go-cli-study/reports/final/09-terminal-ux.md`, `studies/go-cli-study/reports/final/10-logging-observability.md`, `studies/go-cli-study/reports/final/11-testing-strategy.md`, `studies/go-cli-study/reports/final/12-extensibility.md`, `studies/go-cli-study/reports/final/13-security.md`, `studies/go-cli-study/reports/final/14-performance.md`

This handbook distills the studies and reports selected by `sprint-index.md` for sprint reasoning. It does not decide architecture or implementation.

## Selected Studies And Reports

| Study / Report | Path | Relevant Finding | Confidence |
| --- | --- | --- | --- |
| `01-project-structure` | `studies/go-cli-study/reports/final/01-project-structure.md` | Application CLIs keep entrypoints thin and place product behavior behind one-way package boundaries. Examples include the 12-line gh entrypoint and the restic domain/implementation split at `studies/go-cli-study/repos/gh-cli/cmd/gh/main.go:6` and `studies/go-cli-study/repos/restic/internal/restic/repository.go:18`. | high |
| `02-command-architecture` | `studies/go-cli-study/reports/final/02-command-architecture.md` | Factory-created commands delegate to reusable actions, while long `RunE` functions accumulate policy and side effects. Compare `studies/go-cli-study/repos/helm/pkg/cmd/install.go:132-145` with the 130-line orchestrator at `studies/go-cli-study/repos/opencode/cmd/root.go:49-183`. | high |
| `03-dependency-injection` | `studies/go-cli-study/reports/final/03-dependency-injection.md` | Manual composition roots, constructor injection, and narrow interfaces make volatile dependencies replaceable. The gh factory at `studies/go-cli-study/repos/gh-cli/pkg/cmdutil/factory.go:16-43` contrasts with global configuration and cache state at `studies/go-cli-study/repos/rclone/fs/config.go:14-51` and `studies/go-cli-study/repos/rclone/fs/cache/cache.go:16-21`. | high |
| `04-configuration-management` | `studies/go-cli-study/reports/final/04-configuration-management.md` | Configuration needs explicit precedence and validation after sources are combined. Chezmoi restores changed flags after config loading at `studies/go-cli-study/repos/chezmoi/internal/cmd/config.go:2253-2287`; k9s centralizes cross-field validation at `studies/go-cli-study/repos/k9s/internal/config/k9s.go:423-451`. | high |
| `05-error-handling` | `studies/go-cli-study/reports/final/05-error-handling.md` | Typed or sentinel errors preserve machine decisions while a separate rendering path supplies actionable user guidance. Evidence includes gh exit categories at `studies/go-cli-study/repos/gh-cli/internal/ghcmd/cmd.go:44-49`, go-task's structured missing-task error at `studies/go-cli-study/repos/go-task/errors/errors_task.go:13-32`, and age hints at `studies/go-cli-study/repos/age/cmd/age/tui.go:37-54`. | high |
| `06-io-abstraction` | `studies/go-cli-study/reports/final/06-io-abstraction.md` | Injected streams and boundary interfaces make command behavior testable without a real terminal or external system. Gh provides in-memory streams at `studies/go-cli-study/repos/gh-cli/pkg/iostreams/iostreams.go:551-568`; restic separates terminal, filesystem, and backend interfaces at `studies/go-cli-study/repos/restic/internal/ui/terminal.go:10-36`, `studies/go-cli-study/repos/restic/internal/fs/interface.go:10-31`, and `studies/go-cli-study/repos/restic/internal/backend/backend.go:19-90`. | high |
| `07-state-context` | `studies/go-cli-study/reports/final/07-state-context.md` | Long work needs one propagated cancellation lineage, while cleanup may need a distinct bounded context. Go-task uses `errgroup.WithContext` at `studies/go-cli-study/repos/go-task/task.go:89`; restic separates delayed cleanup cancellation at `studies/go-cli-study/repos/restic/internal/restic/lock.go:290-305`. | high |
| `08-concurrency` | `studies/go-cli-study/reports/final/08-concurrency.md` | Localized launch sites, bounded fan-out, explicit waits, and cancellation make goroutine lifecycles auditable. Relevant examples are restic's errgroup fan-out at `studies/go-cli-study/repos/restic/internal/repository/repository.go:567`, k9s's bounded worker pool at `studies/go-cli-study/repos/k9s/internal/pool.go:21-37`, and opencode's shutdown wait at `studies/go-cli-study/repos/opencode/cmd/root.go:252-279`. | high |
| `09-terminal-ux` | `studies/go-cli-study/reports/final/09-terminal-ux.md` | Long operations need truthful progress and cancellation, but scripts need non-TTY behavior and stable plain output. Chezmoi's fallback appears at `studies/go-cli-study/repos/chezmoi/internal/cmd/prompt.go:20-256`; lazygit's interruptible streaming appears at `studies/go-cli-study/repos/lazygit/pkg/tasks/tasks.go:31-435`. | medium |
| `10-logging-observability` | `studies/go-cli-study/reports/final/10-logging-observability.md` | Stable structured fields and strict stdout/stderr separation improve diagnosis without corrupting machine output. K9s centralizes field names at `studies/go-cli-study/repos/k9s/internal/slogs/keys.go:6-231`; helm routes structured diagnostics to stderr at `studies/go-cli-study/repos/helm/internal/logging/logging.go:31-71`. | high |
| `11-testing-strategy` | `studies/go-cli-study/reports/final/11-testing-strategy.md` | Reliable CLIs combine table-driven unit tests, command-level integration, realistic fixtures, fake process boundaries, and output fixtures. Examples include chezmoi testscript coverage at `studies/go-cli-study/repos/chezmoi/internal/cmd/main_test.go:64-174`, lazygit's fake runner at `studies/go-cli-study/repos/lazygit/pkg/commands/oscommands/fake_cmd_obj_runner.go:17-26`, and helm's golden helper at `studies/go-cli-study/repos/helm/internal/test/test.go:43`. | high |
| `12-extensibility` | `studies/go-cli-study/reports/final/12-extensibility.md` | Versioned metadata and explicit registries make extension contracts inspectable, but subprocess extensions commonly lack resource limits. Helm validates versioned plugin metadata at `studies/go-cli-study/repos/helm/internal/plugin/metadata_v1.go:24-48`; its subprocess runtime at `studies/go-cli-study/repos/helm/internal/plugin/runtime_subprocess.go:65-79` illustrates the timeout concern recorded by the report. | medium |
| `13-security` | `studies/go-cli-study/reports/final/13-security.md` | External execution requires explicit argument arrays, path and schema validation, restrictive temporary storage, and redaction. Evidence includes go-task's controlled interpreter at `studies/go-cli-study/repos/go-task/internal/execext/exec.go:59-66`, k9s schema validation at `studies/go-cli-study/repos/k9s/internal/config/json/validator.go:146`, and chezmoi's private temporary directory at `studies/go-cli-study/repos/chezmoi/gpgencryption.go:151-165`. | high |
| `14-performance` | `studies/go-cli-study/reports/final/14-performance.md` | The recurring performance techniques are lazy initialization, streaming, bounded concurrency, and allocation work justified by profiles. Gh defers expensive dependencies through function fields at `studies/go-cli-study/repos/gh-cli/pkg/cmdutil/factory.go:27-42`; age streams fixed-size chunks at `studies/go-cli-study/repos/age/internal/stream/stream.go:20,195-219`; gdu bounds traversal concurrency at `studies/go-cli-study/repos/gdu/pkg/analyze/parallel.go:13`. | high |

## Relevant Patterns

- **Keep policy and orchestration inside the application domain.** The structural study found consistent one-way flow from CLI packages into protected product packages, including yq's `cmd` to library direction at `studies/go-cli-study/repos/yq/cmd/root.go:9` and `studies/go-cli-study/repos/yq/cmd/evaluate_sequence_command.go:7`. For this sprint, that evidence puts pressure on reasoning to keep target parsing, qualification, comparison, and terminal outcomes out of CLI, TUI, browser, and runtime adapters.

- **Use thin lifecycle commands over one shared operation.** Helm's install command delegates through a small command wrapper at `studies/go-cli-study/repos/helm/pkg/cmd/install.go:132-145`, and rclone centralizes cross-cutting execution behavior at `studies/go-cli-study/repos/rclone/cmd/cmd.go:240-340`. The performance lifecycle has enough entry points that duplicated start, resume, cancel, and recover policy would drift across interfaces.

- **Make volatile dependencies explicit.** Gh's factory carries lazy clients and IO as function fields at `studies/go-cli-study/repos/gh-cli/pkg/cmdutil/factory.go:16-43`, while go-task uses functional options at `studies/go-cli-study/repos/go-task/executor.go:22-24`. Benchmark runners, isolation, clocks, environment capture, persistence, and runtime proposal generation are likely test seams. The evidence favors explicit construction over package globals.

- **Separate governed values from operational settings.** Configuration studies show how defaults or direct environment reads can bypass intended precedence. Gdu's flag defaults shadow file values at `studies/go-cli-study/repos/gdu/cmd/gdu/main.go:46-112`, and opencode reads some environment values outside its config layer at `studies/go-cli-study/repos/opencode/internal/config/config.go:163`. Sprint reasoning must account for the stronger rule here: operational inputs may lower limits but may not become an alternate target source.

- **Validate after complete parsing, then freeze identity.** Dive's post-load hook validates combined options at `studies/go-cli-study/repos/dive/cmd/dive/cli/internal/options/analysis.go:48-53`, while helm's versioned metadata at `studies/go-cli-study/repos/helm/internal/plugin/metadata_v1.go:24-48` makes contract versions explicit. The target table, benchmark mapping, parser version, and effective limits need similarly inspectable validation boundaries before measurement starts. The exact schema and digest model remain reasoning questions.

- **Represent outcomes as product facts, not prose.** Go-task's `TaskError` carries an exit code at `studies/go-cli-study/repos/go-task/errors/errors.go:47-50`, gh maps stable exit categories at `studies/go-cli-study/repos/gh-cli/internal/ghcmd/cmd.go:44-49`, and helm exposes sentinels at `studies/go-cli-study/repos/helm/pkg/storage/driver/driver.go:27-48`. Performance target and run outcomes need a typed path that all presenters can render without interpreting runtime text.

- **Run external work through a narrow, bounded execution boundary.** The security report treats explicit argv as the baseline and points to controlled execution in `studies/go-cli-study/repos/go-task/internal/execext/exec.go:59-66`. The extensibility report warns that `studies/go-cli-study/repos/helm/internal/plugin/runtime_subprocess.go:65-79` lacks a built-in timeout. Benchmark, profile, correctness, and cleanup commands therefore create pressure for explicit argv, fixed directories, allowlisted environments, finite output, deadlines, and process-tree cleanup.

- **Propagate cancellation through work and reserve bounded time for cleanup.** Helm wires signals into action context at `studies/go-cli-study/repos/helm/pkg/cmd/install.go:333-347`; restic's delayed cleanup context appears at `studies/go-cli-study/repos/restic/internal/restic/lock.go:290-305`; opencode waits only for a fixed shutdown interval at `studies/go-cli-study/repos/opencode/cmd/root.go:252-279`. This pattern matches a durable operation that must stop runtimes and descendants without treating uncertain cleanup as success.

- **Bound fan-out and retained data.** K9s limits workers at `studies/go-cli-study/repos/k9s/internal/pool.go:21-37`, age holds memory flat through chunked streaming at `studies/go-cli-study/repos/age/internal/stream/stream.go:20,195-219`, and yq processes one document at a time at `studies/go-cli-study/repos/yq/pkg/yqlib/stream_evaluator.go:78-113`. Performance evidence includes samples, command output, profiles, patches, and attempts, so bounds must apply to storage as well as execution.

- **Profile before adding specialized optimizations.** Fzf reuses slabs in a measured hot path at `studies/go-cli-study/repos/fzf/src/matcher.go:183-185`; restic wraps buffer pools with size limits at `studies/go-cli-study/repos/restic/internal/archiver/buffer.go:24-46`; several simpler tools do neither. The study explicitly cautions against speculative pooling. That evidence fits the sprint's profile-led cycle and argues against generic optimization machinery without a target-linked miss.

- **Project one bounded fact model into every interface.** Gdu defines a UI interface at `studies/go-cli-study/repos/gdu/cmd/gdu/app/app.go:30-49`, and gh's test IO constructor at `studies/go-cli-study/repos/gh-cli/pkg/iostreams/iostreams.go:551-568` lets callers share behavior while changing presentation. The CLI, JSON, TUI, and browser should consume the same bounded performance result rather than recompute phase or verdict policy.

- **Test the trust boundary, not only the arithmetic.** Chezmoi's testscript suite exercises command behavior at `studies/go-cli-study/repos/chezmoi/internal/cmd/main_test.go:64-174`; restic uses functional backend fakes at `studies/go-cli-study/repos/restic/internal/backend/mock/backend.go:14-26`; lazygit records process invocations through `studies/go-cli-study/repos/lazygit/pkg/commands/oscommands/fake_cmd_obj_runner.go:17-26`. The sprint needs parser and comparator tables, but its highest-risk failures sit at isolation, process cleanup, writer fencing, promotion, recovery, and interface agreement.

## Trade-Offs

| Trade-Off | Benefit | Cost | When It Matters |
| --- | --- | --- | --- |
| Central composition root vs decentralized wiring | One place exposes runners, stores, clocks, limits, IO, and runtime selection; gh shows the traceability benefit at `studies/go-cli-study/repos/gh-cli/pkg/cmdutil/factory.go:16-43`. | A broad root can become a large coupled object, as chezmoi's config at `studies/go-cli-study/repos/chezmoi/internal/cmd/config.go:193-291` demonstrates. | When performance adds many replaceable boundaries but should not enlarge the existing service into an opaque dependency bag. |
| Lazy admission dependencies vs eager validation | Lazy function fields avoid paying for runtime, process, or repository work during status and dry-run, as in `studies/go-cli-study/repos/gh-cli/pkg/cmdutil/factory.go:27-42`. | First-use failures surface later and may be cached by `sync.Once`. | When dry-run must remain runtime-free while start must still fail before durable work if a prerequisite is missing. |
| Streaming evidence vs complete in-memory snapshots | Streaming keeps memory bounded, shown by age at `studies/go-cli-study/repos/age/internal/stream/stream.go:195-219` and yq at `studies/go-cli-study/repos/yq/pkg/yqlib/stream_evaluator.go:78-113`. | Incremental parsing and atomic publication require careful framing and recovery after partial writes. | When command output, samples, profiles, or progress can exceed public and private record limits. |
| Structured concurrency vs sequential execution | Errgroups and bounded workers improve throughput and propagate failure, as in `studies/go-cli-study/repos/restic/internal/repository/repository.go:567` and `studies/go-cli-study/repos/k9s/internal/pool.go:21-37`. | Concurrency complicates stable benchmark conditions, cancellation, attribution, and cleanup. | When deciding whether independent target preparation may run in parallel while measurements may need serialization for environmental stability. |
| Immutable startup policy vs live configuration | Immutable settings are easier to fingerprint and reproduce. Most one-shot CLIs choose this model, according to `04-configuration-management`. | Operators cannot adjust a running operation, and a restart may be required to lower a limit. | When resume must restore the exact effective limits and reject policy drift. |
| Subprocess isolation vs integration depth | Separate processes contain crashes and allow explicit IO protocols; helm's subprocess runtime is at `studies/go-cli-study/repos/helm/internal/plugin/runtime_subprocess.go:65-79`. | Serialization, startup cost, environment control, timeout handling, and cleanup proof all move to the host. | When a runtime proposes benchmark mappings or patches but cannot receive direct product authority. |
| Rich private evidence vs small public state | Detailed records support recovery and audit, while a bounded projection keeps status fast and safe. Opencode's disk-backed session data at `studies/go-cli-study/repos/opencode/internal/message/message.go:37-42` shows the resource benefit of moving long-lived detail out of memory. | Two representations can disagree unless publication and fingerprints share one owner. | When attempt records retain raw facts but `flow-state.json` and interface DTOs expose only bounded summaries. |
| Golden output fixtures vs semantic assertions | Golden files catch complete output changes, as in `studies/go-cli-study/repos/helm/internal/test/test.go:43`. | Large fixture diffs can hide accidental changes and require deliberate regeneration. | When canonical Markdown, JSON, CLI, TUI, and HTML projections must stay aligned without making every test brittle to formatting. |
| General extension points vs closed descriptors | Registries and versioned metadata ease future additions, as helm shows at `studies/go-cli-study/repos/helm/internal/plugin/metadata.go:114-130`. | General registries can admit collisions, unbounded commands, or authority not required by the sprint. Rclone's silent registry overwrite at `studies/go-cli-study/repos/rclone/fs/rc/registry.go:41-48` is the warning. | When defining v1 parser and benchmark descriptors without pulling a general plugin system into scope. |

## Anti-Patterns And Warnings

- **Shell strings or runtime-selected commands.** Explicit argument arrays are the security baseline. Free-form command text would give a model, requirements cell, or benchmark file execution authority that the product cannot inspect. The report contrasts safe argument construction with unrestricted forwarding such as `studies/go-cli-study/repos/dive/cmd/dive/cli/internal/command/build.go:25`.

- **External work without a timeout or cleanup owner.** The extensibility report flags the missing timeout in `studies/go-cli-study/repos/helm/internal/plugin/runtime_subprocess.go:65-79`. The concurrency report also warns about waits that can hang forever, including `studies/go-cli-study/repos/gh-cli/pkg/cmd/extension/manager.go:196-206`.

- **Package globals for policy, parser, cache, or current run state.** Yq mutates global preferences at `studies/go-cli-study/repos/yq/pkg/yqlib/yaml.go:40`; rclone falls back to global config at `studies/go-cli-study/repos/rclone/fs/config.go:793-802`. Globals would make concurrent attempts, tests, resume, and writer fencing ambiguous.

- **Using context as a service locator.** K9s retrieves application services from context keys at `studies/go-cli-study/repos/k9s/internal/keys.go:10-38`, which the state report calls type-unsafe and prone to missing-value panics. Context should carry cancellation, deadlines, and small request facts, not hidden performance authority.

- **Eager work in status, help, or dry-run.** Gh's eager configuration load penalizes unrelated commands, while its lazy factory at `studies/go-cli-study/repos/gh-cli/pkg/cmdutil/factory.go:27-42` shows the alternative. Performance admission must not accidentally launch a runtime, discover by executing benchmarks, or mutate state during dry-run.

- **Unbounded goroutine, output, or evidence growth.** The concurrency report calls out goroutine-per-item fan-out at `studies/go-cli-study/repos/gh-cli/pkg/cmd/extension/manager.go:198-203`. The performance report warns against full in-memory accumulation and unbounded caches. Sample counts do not justify unbounded retained command output.

- **Direct `os.Stdout`, `os.Stderr`, or filesystem access behind shared use cases.** Such calls bypass test and adapter boundaries. Examples include `studies/go-cli-study/repos/rclone/cmd/ls/ls.go:42` and `studies/go-cli-study/repos/chezmoi/internal/cmd/templatefuncs.go:296`.

- **Panics or prose strings for expected validation failures.** Dive panics on malformed manifests at `studies/go-cli-study/repos/dive/dive/image/docker/manifest.go:18`; string-only errors cannot support stable exit behavior or recovery. Malformed targets, parser uncertainty, noisy measurements, and drift are expected product states.

- **Silent duplicate registration or warning-only schema rejection.** Rclone can overwrite registry entries at `studies/go-cli-study/repos/rclone/fs/rc/registry.go:41-48`, and k9s continues after plugin schema warnings at `studies/go-cli-study/repos/k9s/internal/config/plugin.go:158-164`. A target, benchmark mapping, or terminal record collision must fail closed.

- **Optimizing before profiling.** The performance study supports pools and slabs only in demonstrated hot paths such as `studies/go-cli-study/repos/fzf/src/matcher.go:183-185`. Generic caches, pooling, parallel measurement, or GC tuning could add variance and correctness risk without improving a governed target.

- **Tests tied to incidental internals.** K9s's exact hint-count assertion at `studies/go-cli-study/repos/k9s/internal/view/pod_test.go:23` breaks when unrelated UI details change. Performance tests should assert authority, containment, identities, numeric boundaries, and published outcomes rather than helper call counts unless those counts are governed limits.

## Examples Worth Investigating

| Example | Path / Source | Why It Is Useful |
| --- | --- | --- |
| Gh command factory | `studies/go-cli-study/repos/gh-cli/pkg/cmdutil/factory.go:16-43` | Shows explicit, lazy, test-replaceable dependencies shared by many commands. It is useful when comparing a dedicated performance dependency set with additions to an existing application service. |
| Helm command/action split | `studies/go-cli-study/repos/helm/pkg/cmd/install.go:132-145`; `studies/go-cli-study/repos/helm/pkg/action/install.go:73-140` | Shows how adapters can collect flags and render output while an application/domain object owns behavior. |
| Chezmoi precedence restoration | `studies/go-cli-study/repos/chezmoi/internal/cmd/config.go:2253-2287` | Demonstrates how an apparently simple default can override a stronger source unless explicit-set state is tracked. This is directly relevant to lower-only operational limits. |
| Gh stable exit mapping | `studies/go-cli-study/repos/gh-cli/internal/ghcmd/cmd.go:44-49,281-301` | Combines machine-usable exit categories with user-specific guidance. It can inform reasoning about disabled, blocked, missed, cancelled, and cleanup-uncertain CLI results. |
| Restic delayed cleanup context | `studies/go-cli-study/repos/restic/internal/restic/lock.go:290-305` | Shows a cleanup path that can continue after main work cancellation without becoming unbounded. |
| Opencode bounded shutdown wait | `studies/go-cli-study/repos/opencode/cmd/root.go:252-279` | Provides a concrete local-server cleanup pattern with a finite wait, useful for cancellation and restart reasoning. |
| Restic functional backend mock | `studies/go-cli-study/repos/restic/internal/backend/mock/backend.go:14-26` | Shows concise failure injection at a volatile boundary without a large mocking framework. |
| Chezmoi testscript command suite | `studies/go-cli-study/repos/chezmoi/internal/cmd/main_test.go:64-174` | Demonstrates full command workflows with controlled files, environment, output, and timeouts. |
| Helm versioned metadata conversion | `studies/go-cli-study/repos/helm/internal/plugin/metadata_v1.go:24-48`; `studies/go-cli-study/repos/helm/internal/plugin/metadata.go:114-130` | Useful for examining strict schema versions and conversion boundaries before designing private performance record migration. |
| K9s schema validation | `studies/go-cli-study/repos/k9s/internal/config/json/validator.go:1-187` | Provides an example of central structural validation for user-authored configuration, relevant to exact Markdown table and JSON record validation. |
| Age fixed-chunk streaming | `studies/go-cli-study/repos/age/internal/stream/stream.go:20,195-219` | Shows a simple bounded-memory design for unbounded input. The same pressure applies to command output capture and retained evidence. |
| Fzf slab reuse | `studies/go-cli-study/repos/fzf/src/constants.go:44-45`; `studies/go-cli-study/repos/fzf/src/matcher.go:183-185` | A good counterexample to speculative tuning: specialized allocation appears only in a measured, repeated hot path. |
| Restic bounded buffer pool | `studies/go-cli-study/repos/restic/internal/archiver/buffer.go:24-46` | Shows that even pooling needs size checks and oversize exclusion. |
| Gh in-memory IO constructor | `studies/go-cli-study/repos/gh-cli/pkg/iostreams/iostreams.go:551-568` | Useful for testing the same product DTO and operation through text and JSON adapters without terminal coupling. |

## Design Pressures

- Requirements rows are the only target source. Configuration evidence shows that ordinary precedence systems can accidentally let defaults or environment reads win, so this sprint needs a sharper separation between target authority and lower-only operational policy.
- The performance phase sits between execute and Conformance Review. The architecture must add ordering without turning performance into a planning command or duplicating lifecycle rules in each adapter.
- Dry-run must report meaningful admission facts without launching runtimes or commands. Lazy dependency evidence supports that goal, but delayed errors must not leak past durable acceptance into avoidable partial runs.
- Benchmark identity must remain constant across baseline, candidates, and final measurement. The selected studies cover versioned schemas and validation but not benchmark identity freezing, so reasoning must define that contract rather than infer it from the studies.
- Measurements can be noisy or malformed. Product code must own parsing, qualification, comparison, and inclusive boundaries. None of the selected reports supplies a complete statistical qualification model.
- Optimization proposals combine hostile runtime output with canonical source mutation. Subprocess isolation helps, but the reports show that process isolation alone does not provide time, output, path, environment, or cleanup limits.
- Cancellation crosses runtime work, child processes, disposable copies, private publication, durable ownership, CLI, TUI, and browser observers. A browser disconnect cannot stand in for the operation's cancellation context.
- Resume and recovery need exact durable boundaries. Global state, implicit setup, and warning-only validation would make it impossible to prove what was frozen or consumed before interruption.
- Detailed evidence can grow much faster than `flow-state.json` or public DTOs. The design needs bounded private retention and one canonical projection path so public summaries do not diverge from terminal evidence.
- Correctness must outrank speed. The performance reports discuss resource techniques, but they do not authorize test changes, benchmark changes after freeze, or accepting a faster incorrect candidate.
- Benchmarking itself can perturb the environment. Concurrency and background-refresh patterns that help ordinary application throughput may increase variance during qualification.
- Existing disabled projects must remain byte-for-byte compatible. Activation parsing therefore needs an additive default and must avoid eager target, benchmark, runtime, command, or state work.

## Project Reasoning Applied

| Project Reasoning Document | Accepted Constraint | Sprint Interpretation |
| --- | --- | --- |
| None selected by `projects/ultraplan-go/sprints/39-performance-stage/sprint-index.md` | No accepted current project synthesis is available. | This handbook applies no project-reasoning decision. Sprint reasoning must evaluate the selected study evidence against the governed requirements and must not invent a missing project decision. |

## Open Questions For Reasoning

- Where should policy-aware requirements validation join project policy and sprint content without making project policy an alternate target source?
- What Markdown scanning model can prove that only an actual level-two heading and its immediately associated table count, while excluding fences, comments, blockquotes, and unrelated sections?
- Which target fields should retain normalized decimal text in addition to parsed numeric values so digest stability and exact comparison semantics remain inspectable?
- Should performance persistence reuse QA's writer token and atomic record machinery directly, or should a verification-generic type be introduced without broadening Sprint 39?
- Which benchmark descriptor types can cover Go benchmark output and the UltraPlan JSON envelope while preventing runtime-authored argv, parser, environment, or working-directory authority?
- How will benchmark authoring promotion distinguish test and measurement-support paths from production paths without allowing fixtures or expected outputs to weaken correctness?
- What aggregation method, coefficient-of-variation rule, environment compatibility check, and non-finite handling will qualify baseline and candidate samples?
- Which work may run concurrently without contaminating measurement stability? Preparation and parsing may have different constraints from warmups and measured samples.
- How should a baseline-relative percentage behave near zero, and which exact facts should explain an inconclusive result?
- At what durable boundary does an accepted benchmark or implementation mutation invalidate execute, performance, review, and QA fingerprints?
- How should a one-target optimization cycle prove improvement without allowing an already-met required target to regress under noisy measurements?
- Which correctness commands are frozen, who supplies them, and how does admission report a missing or stale correctness policy?
- How should lower-only configured maxima interact with a requirements-owned `Samples` value when the requested sample count exceeds the effective operational limit?
- What minimum public DTO fields let text, JSON, TUI, and browser views agree while withholding raw samples, profiles, patches, prompts, secrets, and provider payloads?
- How should cancellation arbitrate with late command completion, persistence failure, expired leases, and cleanup uncertainty so exactly one terminal result wins?
- Which private records need independent immutable digests, and which current-state pointers must bind those digests to prevent partial publication from appearing current?
- How will old flow-state schemas and projects with missing performance policy retain their current behavior without creating performance files or changing existing artifact bytes?
- Which tests should use pure table cases, fake process runners, real temporary worktrees, failure hooks, race execution, or gated real-runtime dogfood?

## Evidence Pointers

- `studies/go-cli-study/reports/final/01-project-structure.md`: inspect Thin CLI Entry Point, `internal/` Package Protection, Unidirectional Dependency Flow, UI Interface Abstraction, and the monolith warnings.
- `studies/go-cli-study/reports/final/02-command-architecture.md`: inspect Factory Function Command Creation, `cmd.Run()` Wrapper, lifecycle hooks, and long `RunE` cautions.
- `studies/go-cli-study/reports/final/03-dependency-injection.md`: inspect Constructor Injection, Lazy Initialization, Centralized Composition Root, and global-state cautions.
- `studies/go-cli-study/reports/final/04-configuration-management.md`: inspect explicit precedence, post-load validation, immutability, direct environment bypass, and flag-default shadowing.
- `studies/go-cli-study/reports/final/05-error-handling.md`: inspect typed errors, sentinels, user/operational separation, exit-code mapping, and panic cautions.
- `studies/go-cli-study/reports/final/06-io-abstraction.md`: inspect IOStreams, filesystem and backend interfaces, concurrent-safe test buffers, and direct `os.*` leaks.
- `studies/go-cli-study/reports/final/07-state-context.md`: inspect signal-context wiring, centralized state, delayed cleanup contexts, and context service-locator warnings.
- `studies/go-cli-study/reports/final/08-concurrency.md`: inspect errgroup fan-out, semaphore bounds, explicit shutdown waits, `sync.Once` cleanup, and unbounded goroutine warnings.
- `studies/go-cli-study/reports/final/09-terminal-ux.md`: inspect non-TTY fallback, progressive interruptible output, signal-safe exit, and no-loading-state warnings.
- `studies/go-cli-study/reports/final/10-logging-observability.md`: inspect structured keys, output separation, component tags, bounded debug controls, and missing evidence discipline.
- `studies/go-cli-study/reports/final/11-testing-strategy.md`: inspect testscript, golden files, fake command runners, failure isolation, and behavior-vs-internal assertions.
- `studies/go-cli-study/reports/final/12-extensibility.md`: inspect subprocess isolation, versioned metadata, duplicate registry failures, schema warnings, and absent plugin resource limits.
- `studies/go-cli-study/reports/final/13-security.md`: inspect explicit argv, private temporary directories, schema validation, redaction, trust markers, and command-timeout cautions.
- `studies/go-cli-study/reports/final/14-performance.md`: inspect lazy initialization, streaming, bounded concurrency, profiling hooks, measured pooling, and speculative optimization warnings.

## Handoff To Reasoning

- Use this handbook as evidence input.
- Validate whether the observed patterns fit this project's constraints.
- Do not copy external patterns without sprint-specific reasoning.
