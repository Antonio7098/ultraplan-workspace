# Sprint Technical Handbook: Bounded Manual and Automatic Repair

> Project: `ultraplan-go`
> Sprint: `38-bounded-repair`
> Source: `sprint-index.md`
> **Inputs Used:** `projects/ultraplan-go/project-index.md`, `projects/ultraplan-go/sprints/38-bounded-repair/sprint-index.md`, `projects/ultraplan-go/sprints/38-bounded-repair/requirements.md`, `studies/go-cli-study/reports/final/01-project-structure.md`, `studies/go-cli-study/reports/final/02-command-architecture.md`, `studies/go-cli-study/reports/final/03-dependency-injection.md`, `studies/go-cli-study/reports/final/04-configuration-management.md`, `studies/go-cli-study/reports/final/05-error-handling.md`, `studies/go-cli-study/reports/final/06-io-abstraction.md`, `studies/go-cli-study/reports/final/07-state-context.md`, `studies/go-cli-study/reports/final/08-concurrency.md`, `studies/go-cli-study/reports/final/09-terminal-ux.md`, `studies/go-cli-study/reports/final/10-logging-observability.md`, `studies/go-cli-study/reports/final/11-testing-strategy.md`, `studies/go-cli-study/reports/final/12-extensibility.md`, `studies/go-cli-study/reports/final/13-security.md`, `studies/go-cli-study/reports/final/14-performance.md`

This handbook distills the studies and reports selected by `sprint-index.md` for sprint reasoning. It does not decide architecture or implementation.

## Selected Studies And Reports

| Study / Report | Path | Relevant Finding | Confidence |
| --- | --- | --- | --- |
| `01-project-structure` | `studies/go-cli-study/reports/final/01-project-structure.md` | Mature applications keep entrypoints thin, protect application policy in `internal/`, and preserve one-way imports. Helm separates command wiring from action behavior at `helm/pkg/cmd/install.go:347` and `helm/pkg/action/install.go:73-140`; gdu places multiple renderers behind `gdu/cmd/gdu/app/app.go:30-49`. | high |
| `02-command-architecture` | `studies/go-cli-study/reports/final/02-command-architecture.md` | Factory-built commands and shared execution wrappers keep command parsing separate from lifecycle and business behavior. Evidence includes gh-cli `pkg/cmdutil/factory.go:16-43`, Helm `pkg/cmd/install.go:132-145`, and rclone `cmd/cmd.go:240-340`. | high |
| `03-dependency-injection` | `studies/go-cli-study/reports/final/03-dependency-injection.md` | Explicit composition roots, constructor injection, and interfaces make runtime and storage behavior replaceable in tests. Gh-cli uses lazy dependency functions in `pkg/cmd/factory/default.go:26-46`; restic defines repository and backend boundaries at `internal/restic/repository.go:18-66` and `internal/backend/backend.go:19-90`. | high |
| `04-configuration-management` | `studies/go-cli-study/reports/final/04-configuration-management.md` | Layered configuration needs explicit precedence and validation after all sources merge. Go-task centralizes precedence at `internal/flags/flags.go:314-327`; restic tracks explicit flags at `internal/global/global.go:139,147`; k9s centralizes validation at `internal/config/k9s.go:423-451`. | high |
| `05-error-handling` | `studies/go-cli-study/reports/final/05-error-handling.md` | Typed or behavioral errors preserve facts needed for recovery and exit-code mapping. Rclone defines retry and fatal behavior at `fs/fserrors/error.go:22-99`; restic separates fatal outcomes at `internal/errors/fatal.go:10-53`; gh-cli maps distinct exit classes at `internal/ghcmd/cmd.go:44-49`. | high |
| `06-io-abstraction` | `studies/go-cli-study/reports/final/06-io-abstraction.md` | Injectable terminal, filesystem, process, and backend boundaries permit failure tests without real side effects. Gh-cli provides test streams at `pkg/iostreams/iostreams.go:551-568`; restic defines terminal and filesystem interfaces at `internal/ui/terminal.go:10-36` and `internal/fs/interface.go:10-31`. | high |
| `07-state-context` | `studies/go-cli-study/reports/final/07-state-context.md` | Long work needs root-context propagation, while cleanup may need a deliberately bounded context after work cancellation. Helm wires signals to action cancellation at `pkg/cmd/install.go:333-347`; restic separates cleanup lifetime at `internal/restic/lock.go:290-305`. | high |
| `08-concurrency` | `studies/go-cli-study/reports/final/08-concurrency.md` | Auditable concurrency has localized launch sites, bounded fan-out, cancellation, and explicit waits. Restic uses `errgroup` at `internal/repository/repository.go:567`; k9s bounds workers at `internal/pool.go:21-37`; opencode waits for shutdown with a timeout at `cmd/root.go:261-279`. | high |
| `09-terminal-ux` | `studies/go-cli-study/reports/final/09-terminal-ux.md` | Long operations need visible progress, non-TTY behavior, and interruptible UI updates. Chezmoi implements non-TTY prompt paths at `internal/cmd/prompt.go:20-256`; lazygit couples streaming views to stop channels at `pkg/tasks/tasks.go:31-435`; restic routes progress through a terminal worker at `internal/ui/termstatus/status.go:197-205`. | medium |
| `10-logging-observability` | `studies/go-cli-study/reports/final/10-logging-observability.md` | Structured, correlated diagnostics should remain separate from user output. K9s centralizes structured keys at `internal/slogs/keys.go:6-231`; restic wraps backends for logging at `internal/backend/logger/log.go:22-77`; Helm sends structured diagnostics to stderr at `internal/logging/logging.go:31-71`. | high |
| `11-testing-strategy` | `studies/go-cli-study/reports/final/11-testing-strategy.md` | Strong CLI suites combine table tests, fakes, real command-path integration, and selective golden fixtures. Chezmoi runs txtar workflows at `internal/cmd/main_test.go:64-174`; gh-cli uses testscript at `acceptance/acceptance_test.go:26-29`; restic exposes functional backend mocks at `internal/backend/mock/backend.go:14-26`. | high |
| `12-extensibility` | `studies/go-cli-study/reports/final/12-extensibility.md` | Versioned metadata and subprocess isolation can stabilize external execution, but registries and plugins create collision, timeout, and trust problems. Helm's subprocess runtime is at `internal/plugin/runtime_subprocess.go:65-79` and version conversion at `internal/plugin/metadata.go:114-130`; rclone's registry silently overwrites at `fs/rc/registry.go:41-48`. | medium |
| `13-security` | `studies/go-cli-study/reports/final/13-security.md` | Consequential operations benefit from explicit permission checks, argument-array execution, private temporary storage, and redaction. Opencode gates commands at `internal/permission/permission.go:44-108`; chezmoi creates private temporary directories at `gpgencryption.go:151-165`; restic redacts secrets through `internal/options/secret_string.go:15-20`. | high |
| `14-performance` | `studies/go-cli-study/reports/final/14-performance.md` | Bounded concurrency, streaming, incremental state, and lazy initialization prevent work from growing without limit. Gh-cli defers dependencies at `pkg/cmdutil/factory.go:27-42`; lazygit streams lines at `pkg/tasks/tasks.go:189-217`; restic bounds file work at `internal/archiver/file_saver.go:56-58`. | high |

## Relevant Patterns

- **One product path behind thin adapters.** The structure and command reports consistently place behavior behind command wiring. Helm's command/action split at `helm/pkg/cmd/install.go:132-145` and `helm/pkg/action/install.go:73-140`, plus gdu's renderer interface at `gdu/cmd/gdu/app/app.go:30-49`, give reasoning a concrete model for one repair protocol consumed by CLI, TUI, and browser adapters. The reports do not settle the exact UltraPlan package split.

- **Explicit composition and replaceable side effects.** Gh-cli's `Factory` at `pkg/cmdutil/factory.go:16-43` and restic's terminal, filesystem, and backend interfaces at `internal/ui/terminal.go:10-36`, `internal/fs/interface.go:10-31`, and `internal/backend/backend.go:19-90` show how to inject runtime, process, storage, clock, and output behavior. This matters for cancellation, stale-writer, partial-publication, and cleanup failure tests. The DI report also warns that a central object can grow into a large dependency container, as in chezmoi `internal/cmd/config.go:193-291`.

- **Admission before expensive or mutable work.** Configuration reports favor merging sources and validating the effective result once, as in k9s `internal/config/k9s.go:423-451`. Security reports add an explicit approval boundary, represented by opencode `internal/permission/permission.go:44-108`. Together they support reasoning about a distinct prepare and confirmation sequence before runtime startup, without deciding how UltraPlan encodes that sequence.

- **Immutable authority plus typed terminal facts.** Helm's versioned metadata conversion at `internal/plugin/metadata.go:114-130` illustrates strict schema evolution. Rclone's behavioral errors at `fs/fserrors/error.go:22-99`, restic's fatal error handling at `internal/errors/fatal.go:10-53`, and gh-cli's exit classes at `internal/ghcmd/cmd.go:44-49` show that operational categories should remain inspectable rather than collapse into message strings. Sprint reasoning still needs to map durable run states and repair outcomes without confusing the two vocabularies.

- **Cancellation reaches work, then bounded cleanup proves shutdown.** Helm propagates a signal-derived context into actions at `pkg/cmd/install.go:333-347`. Restic's `delayedCancelContext` at `internal/restic/lock.go:290-305` demonstrates why cleanup lifetime may differ from work lifetime. Opencode's timed shutdown wait at `cmd/root.go:261-279` shows the need for a finite end to cleanup waiting. This combination closely matches repair's process-tree, workspace, lock, and server-shutdown pressures.

- **Localized and bounded concurrency.** Restic's `errgroup` use at `internal/repository/repository.go:567` provides sibling cancellation and error collection, while k9s `internal/pool.go:21-37` and gdu `pkg/analyze/parallel.go:13,36` show explicit concurrency limits. The repair protocol is mostly ordered, but runtime subprocesses, observers, cancellation, and cleanup still create concurrent ownership risks that reasoning must make auditable.

- **Isolated external execution with explicit trust limits.** Subprocess execution is the dominant extension boundary in the extensibility report, including Helm `internal/plugin/runtime_subprocess.go:65-79`. The security report strengthens that pattern with argument arrays, permission gates, private temporary directories, and redacted values. Relevant examples are opencode `internal/llm/tools/bash.go:41-55`, chezmoi `gpgencryption.go:151-165`, and restic `internal/options/secret_string.go:15-20`. These sources support isolation as a pressure, not a decision to copy a plugin architecture.

- **Canonical status separated from transient progress.** Restic's terminal worker at `internal/ui/termstatus/status.go:197-205` and lazygit's streamed view manager at `pkg/tasks/tasks.go:31-435` treat progress as an update channel with its own lifecycle. Structured logging examples such as k9s `internal/slogs/keys.go:6-231` add stable correlation fields. For repair, this raises a clear question about how observers refresh canonical state after lost or coalesced events.

- **Finite configuration and resource policy.** Go-task's generic precedence accessor at `internal/flags/flags.go:314-327` and restic's explicit flag tracking at `internal/global/global.go:139,147` show how to distinguish defaults from user choices. Performance examples bound active work and avoid complete buffering, including restic `internal/archiver/file_saver.go:56-58` and yq `pkg/yqlib/stream_evaluator.go:78-113`. Reasoning must apply those lessons to lower-only budgets, retained cycles, output, patch bytes, commands, and wall time.

- **Layered behavioral tests.** The testing report pairs table-driven unit coverage with real command-path tests and purpose-built fakes. Chezmoi's txtar runner at `internal/cmd/main_test.go:64-174`, gh-cli's acceptance runner at `acceptance/acceptance_test.go:26-29`, and restic's functional backend mock at `internal/backend/mock/backend.go:14-26` are useful models for testing the same protocol at domain, app, and adapter levels. Golden fixtures can protect public envelopes, while semantic assertions remain necessary for fencing and terminal arbitration.

## Trade-Offs

| Trade-Off | Benefit | Cost | When It Matters |
| --- | --- | --- | --- |
| One shared operation path versus adapter-specific shortcuts | Shared behavior reduces drift between CLI, TUI, and browser. Helm's command/action split and gdu's UI interface show the benefit. | DTO mapping and dependency wiring grow, and the shared coordinator can become too broad. | Repair facts and mutations must agree across all interfaces. Evidence: `studies/go-cli-study/reports/final/01-project-structure.md` and `02-command-architecture.md`. |
| Eager admission validation versus lazy runtime setup | Early validation rejects stale or unsafe work before side effects. | Eager loading can slow harmless queries or surface dependencies that a query does not need, as the performance report notes for gh-cli config loading. | Packet preparation and confirmation need strong checks, while status and help should remain cheap. Evidence: `k9s/internal/config/k9s.go:423-451`; `gh-cli/pkg/cmdutil/factory.go:27-42`. |
| Central composition root versus narrowly scoped constructors | One composition point makes the authority chain traceable and testable. | A large root can become a god object, like chezmoi's `Config` at `internal/cmd/config.go:193-291`. | Repair adds runtime, process, persistence, writer-fence, clock, and adapter dependencies. Evidence: `studies/go-cli-study/reports/final/03-dependency-injection.md`. |
| Structured concurrency versus a mostly sequential protocol | `errgroup` and bounded workers give reliable cancellation and waiting for independent work. | Parallel work complicates gate ordering, writer ownership, and deterministic evidence. | Reverification is ordered, but process cleanup and observers may overlap. Evidence: `restic/internal/repository/repository.go:567`; `k9s/internal/pool.go:21-37`. |
| Streaming progress versus canonical refresh | Streaming gives immediate feedback and bounds memory. | Events can be coalesced, dropped, reordered, or observed after authority has changed. | CLI progress, TUI refresh, browser SSE, reconnect, and restart must not derive outcomes from events. Evidence: `lazygit/pkg/tasks/tasks.go:189-217`; `opencode/internal/llm/provider/provider.go:56`. |
| Subprocess isolation versus deep runtime integration | A process boundary limits host corruption and permits independent runtime implementation. | Serialization, startup, timeout, process-tree cleanup, and trust handling become explicit work. | Proposal generation must remain isolated from production mutation. Evidence: `helm/internal/plugin/runtime_subprocess.go:65-79`; `studies/go-cli-study/reports/final/13-security.md`. |
| Golden public-output tests versus semantic assertions | Goldens catch accidental CLI, JSON, and rendered-output drift. | Large fixture diffs can hide mistakes and do not prove races, fencing, or authority. | Stable public envelopes need regression coverage, while mutation and recovery need fact-based tests. Evidence: `helm/internal/test/test.go:43`; `rclone/cmd/bisync/bisync_test.go:1435-1479`. |
| Automatic reuse of the manual protocol versus a separate automatic engine | Reuse keeps scope checks, apply, cleanup, and outcomes aligned. | One coordinator must express extra limits, progress checks, and resume state without becoming opaque. | Automatic mode must remain lower-only and manual-first. The reports support shared wrappers and bounded work but do not decide the coordinator shape. |

## Anti-Patterns And Warnings

- **Do not put repair semantics in interface handlers.** Long command handlers such as opencode `cmd/root.go:49-183` and yq `cmd/evaluate_sequence_command.go:152` are cited as logic leakage. Separate adapters should not derive admission, scope, progress, or outcomes.

- **Do not hide dependencies in globals or context values.** Rclone's global fallback at `fs/config.go:793-802`, yq's mutable parser/config state at `pkg/yqlib/lib.go:13-21`, and k9s's context-key service lookup at `internal/keys.go:10-38` make ownership and test isolation harder. Repair writer identity, policy, target identity, and consumed limits need visible inputs.

- **Do not use shell strings or unrestricted subprocesses.** The security report treats argument arrays as the baseline and flags missing timeouts as a recurring risk. Opencode's command filter at `internal/llm/tools/bash.go:41-55` and restic's quote-aware splitting at `internal/backend/shell_split.go:45-76` are safer evidence points, though neither alone satisfies repair isolation.

- **Do not launch fire-and-forget work.** Dive's untracked notification goroutine at `cmd/dive/cli/internal/command/adapter/resolver.go:70` and no-timeout waits such as gh-cli `pkg/cmd/extension/manager.go:196-206` show how shutdown can hang or outlive ownership. Every repair goroutine and child process needs cancellation and a bounded join.

- **Do not let fan-out, output, or retained state grow without a hard bound.** The concurrency report warns about goroutine-per-item patterns, while the performance report warns about full buffering and unbounded caches. Examples include gh-cli `pkg/cmd/extension/manager.go:198-203` and lazygit's unevicted pipe cache noted in the performance report.

- **Do not treat cleanup timeout as successful cleanup.** Opencode's timed shutdown at `cmd/root.go:261-279` proves only that waiting is bounded. Restic's cleanup context at `internal/restic/lock.go:290-305` shows cleanup needs separate lifecycle care. Reasoning must preserve uncertainty when process, workspace, or lock cleanup cannot be proved.

- **Do not bypass injected IO, filesystem, or logging paths.** Direct writes such as rclone `cmd/ls/ls.go:42`, urfave-cli `cli.go:47`, and chezmoi `internal/cmd/templatefuncs.go:296` create untestable and potentially unredacted paths. Public repair output and diagnostics need consistent bounds and redaction.

- **Do not collapse user messages, operational diagnostics, and machine outcomes.** K9s `internal/model/flash.go:100-103` deliberately routes UI and logs separately, while gh-cli maps distinct exit codes at `internal/ghcmd/cmd.go:44-49`. A friendly explanation cannot replace the stored repair outcome or blocker category.

- **Do not accept silent schema, registry, or configuration conflicts.** Rclone's registry overwrite at `fs/rc/registry.go:41-48`, k9s's warning-only plugin schema failures at `internal/config/plugin.go:158-164`, and gdu's non-fatal config error at `cmd/gdu/main.go:242-244` are weak precedents for security-sensitive state. Repair should fail closed on unknown schema or conflicting authority.

- **Do not add a general plugin framework to solve isolated proposal execution.** The extensibility report documents versioning, collision, timeout, signature, and resource-limit gaps across plugin systems. Sprint reasoning should distinguish a fixed product-owned runtime boundary from user extensibility.

- **Do not use output snapshots as the only concurrency or security proof.** The testing report warns that snapshots and implementation-detail assertions can be brittle. Fencing, path containment, symlink and hard-link handling, cancellation races, and exactly-one terminal publication need direct behavioral assertions.

## Examples Worth Investigating

| Example | Path / Source | Why It Is Useful |
| --- | --- | --- |
| Helm command/action split | `helm/pkg/cmd/install.go:132-145`; `helm/pkg/action/install.go:73-140`, cited by `01-project-structure` and `02-command-architecture` | Shows transport-facing command setup delegating to a reusable operation without reverse imports. |
| Gh-cli lazy factory | `gh-cli/pkg/cmdutil/factory.go:16-43`; `gh-cli/pkg/cmd/factory/default.go:26-46`, cited by `03-dependency-injection` | Shows explicit, testable dependencies that defer expensive setup until a mutation needs it. |
| Restic cleanup lifetime | `restic/internal/restic/lock.go:290-305`, cited by `07-state-context` | Gives a concrete example of cleanup continuing under a controlled context after work cancellation. |
| Opencode shutdown wait | `opencode/cmd/root.go:261-279`, cited by `08-concurrency` | Shows cancellation followed by a finite wait, useful for reasoning about truthful shutdown outcomes. |
| Opencode permission request | `opencode/internal/permission/permission.go:44-108`; `opencode/internal/llm/tools/bash.go:41-55`, cited by `13-security` | Shows a separate, visible authority check before dangerous execution, including deny and allow classifications. |
| Chezmoi private temporary workspace | `chezmoi/gpgencryption.go:151-165`, cited by `13-security` | Demonstrates restrictive temporary-directory handling for sensitive isolated work. |
| Restic replaceable IO and storage | `restic/internal/ui/terminal.go:10-36`; `restic/internal/fs/interface.go:10-31`; `restic/internal/backend/mock/backend.go:14-26`, cited by `06-io-abstraction` and `11-testing-strategy` | Useful for failure injection around persistence, output, and filesystem operations. |
| Lazygit streamed, stoppable view | `lazygit/pkg/tasks/tasks.go:31-435`, cited by `08-concurrency`, `09-terminal-ux`, and `14-performance` | Shows bounded incremental presentation tied to an explicit stop lifecycle. |
| Helm versioned metadata conversion | `helm/internal/plugin/metadata_v1.go:24-48`; `helm/internal/plugin/metadata.go:114-130`, cited by `12-extensibility` | Useful for examining strict version dispatch and compatibility without treating unknown schemas as current. |
| Chezmoi and gh-cli command-path tests | `chezmoi/internal/cmd/main_test.go:64-174`; `gh-cli/acceptance/acceptance_test.go:26-29`, cited by `11-testing-strategy` | Shows end-to-end command workflows that can freeze confirmation, output, and exit behavior. |
| Rclone behavioral error wrappers | `rclone/fs/fserrors/error.go:22-192`, cited by `05-error-handling` | Shows error facts that survive wrapping and remain available for retry or fatal classification. |
| K9s schema validation | `k9s/internal/config/json/validator.go:1-187`, cited by `04-configuration-management` and `13-security` | Shows central validation of user-authored structured input before runtime use. |

## Design Pressures

- Repair has two authorities that must stay distinct: durable operation ownership and semantic repair outcome. The error and state reports support typed facts, but sprint reasoning must define their mapping.

- Confirmation must authorize one exact mutable request. Permission-dialog evidence supports a separate approval step, while configuration evidence requires the effective limits and sources to be fixed before validation.

- Proposal generation and production application have different trust levels. Subprocess and private-workspace evidence support isolation, but product code still needs an independently reasoned apply boundary.

- The mutation sequence is narrow and ordered, while cancellation, observation, persistence, process cleanup, and server shutdown are concurrent. Launch sites and ownership must remain reviewable.

- Automatic work increases state and budget complexity without gaining a second authority path. Limits must remain finite across retries and restarts, and progress must be a product fact rather than a runtime claim.

- Full private evidence and public summaries have different consumers. The reports favor bounded streaming and stable DTOs, but sprint reasoning must decide which facts are canonical, immutable, paged, or summarized.

- Cleanup is part of the semantic result. A bounded wait prevents shutdown hangs but cannot turn missing cleanup evidence into success.

- Four interfaces need the same facts under different interaction constraints. CLI and JSON need stable exits and envelopes; TUI and browser need explicit confirmation, reconnect, hostile-text bounds, and canonical refresh.

- The protected-path set is broader than ordinary path containment. It includes tests, evidence, configuration, Git control files, generated artifacts, and side effects from descendants or formatters.

- Real-runtime proof is stronger than fixtures. Unit and integration tests can prove guards and races, but they cannot manufacture the qualifying manual proof needed for automatic admission.

## Open Questions For Reasoning

- What is the smallest domain API that keeps packet freezing, confirmation validation, proposal, apply, reverification, cleanup, and outcome derivation in one authority chain without creating a single oversized coordinator?

- Which repair dependencies warrant interfaces for failure injection, and which should remain concrete to keep the authority path easy to audit?

- At what exact point does durable acceptance become part of the single-use confirmation digest, and how does replay detection survive process restart?

- How should durable operation states such as `cancelled`, `interrupted`, and `cleanup_uncertain` map to repair records without allowing either store to overwrite the other's authority?

- Which checks can run concurrently, if any, without violating the fixed widening order or obscuring the first authoritative failure?

- Which context owns bounded cleanup after cancellation, and what facts prove process-tree termination, workspace removal, and lock release before that context expires?

- How should the product compute and compare full target identity around isolated proposal, production apply, and reverification without turning unrelated pre-existing changes into repair scope?

- What strict patch representation and parser can validate paths, renames, deletes, links, changed bytes, and actual post-apply differences without invoking Git or a shell?

- Which lower-only configuration fields share one precedence mechanism, and how will status explain default, workspace, and environment sources without exposing unsafe values?

- What is the minimal immutable evidence needed at each cycle boundary so resume never repeats a committed apply and never resets consumed limits?

- How should repeated patches, unchanged issue sets, severity changes, and new facts be normalized so automatic progress and stagnation remain deterministic?

- Which canonical query DTOs let CLI, JSON, TUI, and browser agree while keeping private command output, prompts, provider payloads, and production content hidden?

- How should progress events identify replay gaps and stale ownership while remaining explicitly non-authoritative?

- Which public outputs merit golden fixtures, and which security, race, recovery, and outcome properties require semantic assertions or real filesystem/process tests?

- How will a real manual proof bind schema, policy, isolation capability, code, governed inputs, cleanup, and target identity tightly enough to invalidate automatic admission for the right reasons?

## Evidence Pointers

- `studies/go-cli-study/reports/final/01-project-structure.md`: inspect "Thin CLI Entry Point," "UI Interface Abstraction," unidirectional imports, and global-state cautions.

- `studies/go-cli-study/reports/final/02-command-architecture.md`: inspect factory-built commands, shared execution wrappers, lifecycle hooks, and oversized `RunE` warnings.

- `studies/go-cli-study/reports/final/03-dependency-injection.md`: inspect composition roots, interface boundaries, lazy factories, large config objects, and context-as-service-locator warnings.

- `studies/go-cli-study/reports/final/04-configuration-management.md`: inspect explicit precedence, post-merge validation, `Flag.Changed` tracking, lower-level direct environment access, and silent config failures.

- `studies/go-cli-study/reports/final/05-error-handling.md`: inspect typed errors, behavioral interfaces, user and operational separation, multi-error aggregation, and exit-code mapping.

- `studies/go-cli-study/reports/final/06-io-abstraction.md`: inspect test IO constructors, terminal/filesystem/backend interfaces, concurrent-safe buffers, and direct `os.*` bypasses.

- `studies/go-cli-study/reports/final/07-state-context.md`: inspect signal-context wiring, `errgroup`, delayed cleanup contexts, explicit lock/session models, and uncancellable background contexts.

- `studies/go-cli-study/reports/final/08-concurrency.md`: inspect localized goroutine ownership, bounded workers, timed waits, `sync.Once` cleanup, and fire-and-forget warnings.

- `studies/go-cli-study/reports/final/09-terminal-ux.md`: inspect non-TTY fallback, progressive interruptible updates, signal-safe exit, and cancellation cleanup gaps.

- `studies/go-cli-study/reports/final/10-logging-observability.md`: inspect structured correlation keys, output separation, component loggers, backend decorators, and redaction risks.

- `studies/go-cli-study/reports/final/11-testing-strategy.md`: inspect testscript workflows, functional mocks, golden output tests, real fixtures, global-state isolation, and brittle implementation assertions.

- `studies/go-cli-study/reports/final/12-extensibility.md`: inspect subprocess isolation, versioned metadata, registry collision handling, plugin timeout gaps, and reasons to delay general extension systems.

- `studies/go-cli-study/reports/final/13-security.md`: inspect permission requests, argument-array execution, command filters, secret redaction, private temporary directories, schema validation, and missing process limits.

- `studies/go-cli-study/reports/final/14-performance.md`: inspect lazy initialization, bounded concurrency, streaming, incremental retention, cleanup on errors, and unbounded accumulation warnings.

## Handoff To Reasoning

- Use this handbook as evidence input.
- Validate whether the observed patterns fit this project's constraints.
- Do not copy external patterns without sprint-specific reasoning.
