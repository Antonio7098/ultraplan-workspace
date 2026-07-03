# Sprint Technical Handbook: 09 Runtime Integration

> Project: `ultraplan-go`
> Sprint: `09-runtime-integration`
> Source: `sprint-index.md`
> **Inputs Used:** `projects/ultraplan-go/sprints/09-runtime-integration/sprint-index.md`, `templates/technical-handbook.md`, `projects/ultraplan-go/project-index.md`, `projects/ultraplan-go/sprints/09-runtime-integration/requirements.md`, `projects/ultraplan-go/docs/PRD.md`, `projects/ultraplan-go/docs/TRD.md`, `projects/ultraplan-go/docs/ARCHITECTURE.md`, `studies/go-cli-study/reports/final/01-project-structure.md`, `studies/go-cli-study/reports/final/02-command-architecture.md`, `studies/go-cli-study/reports/final/03-dependency-injection.md`, `studies/go-cli-study/reports/final/04-configuration-management.md`, `studies/go-cli-study/reports/final/05-error-handling.md`, `studies/go-cli-study/reports/final/06-io-abstraction.md`, `studies/go-cli-study/reports/final/07-state-context.md`, `studies/go-cli-study/reports/final/08-concurrency.md`, `studies/go-cli-study/reports/final/10-logging-observability.md`, `studies/go-cli-study/reports/final/11-testing-strategy.md`, `studies/go-cli-study/reports/final/12-extensibility.md`, `studies/go-cli-study/reports/final/13-security.md`, `studies/go-cli-study/reports/final/15-philosophy.md`

This handbook distills the studies and reports selected by `sprint-index.md` for sprint reasoning. It does not decide architecture or implementation.

## Selected Studies And Reports

| Study / Report | Path | Relevant Finding | Confidence |
| --- | --- | --- | --- |
| `01-project-structure` | `studies/go-cli-study/reports/final/01-project-structure.md` | Mature CLIs keep entrypoints thin and enforce one-way dependencies from CLI to interior logic; opencode uses `cmd/` plus `internal/` with service/pubsub boundaries, while restic shows strict `cmd/restic/` to `internal/` layering (`cmd/root.go:255-259`, `cmd/restic/main.go:37-114`, `internal/restic/repository.go:18`). | high |
| `02-command-architecture` | `studies/go-cli-study/reports/final/02-command-architecture.md` | Factory-created commands with shared lifecycle hooks and thin `RunE` functions keep runtime setup out of command bodies; examples include gh-cli factory wiring and helm action delegation (`pkg/cmdutil/factory.go:16-43`, `pkg/cmd/install.go:132-145`). | high |
| `03-dependency-injection` | `studies/go-cli-study/reports/final/03-dependency-injection.md` | Manual composition roots, constructor injection, functional options, and decorators dominate; no studied repo uses DI frameworks, and global config/service singletons are repeatedly flagged as testability risks (`internal/app/app.go:42-81`, `executor.go:22-24`, `ui_concurrent.go:9-12`). | high |
| `04-configuration-management` | `studies/go-cli-study/reports/final/04-configuration-management.md` | Clear precedence, merged validation, and explicit flag-change tracking are the stable config patterns; opencode and chezmoi demonstrate flag/config restore, go-task generic precedence, and opencode validation (`internal/config/config.go:218-227`, `internal/flags/flags.go:314-327`, `internal/config/config.go:609-641`). | high |
| `05-error-handling` | `studies/go-cli-study/reports/final/05-error-handling.md` | Runtime-facing errors should preserve typed/sentinel classification via wrapping and `errors.Is`/`errors.As`; rclone behavioral error interfaces, restic fatal errors, and gh-cli exit-code mapping are the strongest examples (`fs/fserrors/error.go:22-29`, `internal/errors/fatal.go:10`, `internal/ghcmd/cmd.go:44-49`). | high |
| `06-io-abstraction` | `studies/go-cli-study/reports/final/06-io-abstraction.md` | IO and filesystem seams make commands and health output testable; gh-cli `IOStreams.Test`, restic `Terminal`/`Backend`, and go-task `WithStdout` show production/test symmetry (`pkg/iostreams/iostreams.go:551-568`, `internal/ui/terminal.go:10-36`, `executor.go:553-564`). | high |
| `07-state-context` | `studies/go-cli-study/reports/final/07-state-context.md` | Long-running or external-runtime work benefits from root context propagation, signal cancellation, explicit session state, and avoiding `context.Background()` in work paths (`internal/ghcmd/cmd.go:142`, `pkg/cmd/install.go:333-347`, `internal/session/session.go:12-23`). | high |
| `08-concurrency` | `studies/go-cli-study/reports/final/08-concurrency.md` | Structured goroutine lifecycle, cancellation, and bounded fan-out avoid leaks; errgroup, semaphores, wait-with-timeout, and event brokers appear across go-task, restic, gdu, and opencode (`task.go:87`, `internal/repository/repository.go:567`, `pkg/analyze/parallel.go:13`, `cmd/root.go:261-279`). | high |
| `10-logging-observability` | `studies/go-cli-study/reports/final/10-logging-observability.md` | Structured logs, runtime debug controls, stdout/stderr separation, and event/pubsub models support operator diagnostics without mixing user output and logs (`internal/logging/logger.go:25-62`, `internal/pubsub/broker.go:10-19`, `pkg/iostreams/iostreams.go:52-54`). | high |
| `11-testing-strategy` | `studies/go-cli-study/reports/final/11-testing-strategy.md` | Reliable CLI projects pair table-driven unit tests with fake/mock infrastructure and command-level or script-level integration; sparse coverage around opencode was explicitly identified as a risk (`internal/chezmoitest/chezmoitest.go:86-92`, `pkg/httpmock/stub.go:35-199`, `opencode/internal/tui/theme/theme_test.go:1-89`). | high |
| `12-extensibility` | `studies/go-cli-study/reports/final/12-extensibility.md` | Extension systems usually prefer subprocess isolation, registries, versioned metadata, or protocol boundaries; agent/tool systems like opencode use formal tool/MCP interfaces rather than direct core command mutation (`internal/llm/tools/tools.go:69-72`, `internal/llm/agent/mcp-tools.go:106-129`). | medium |
| `13-security` | `studies/go-cli-study/reports/final/13-security.md` | Agent-style execution needs explicit permissions, safe command filtering, secret redaction, and fail-fast validation at trust boundaries; opencode permission flow and banned/safe command filters are especially relevant (`internal/permission/permission.go:44-108`, `internal/llm/tools/bash.go:41-55`). | high |
| `15-philosophy` | `studies/go-cli-study/reports/final/15-philosophy.md` | Strong Go CLIs accept complexity deliberately and reject scope that does not serve the tool; factory/options, interfaces, wrappers, and pub/sub are recurring patterns, but plugin/security/config complexity must be justified (`pkg/cmd/factory/default.go:26-46`, `internal/pubsub/broker.go:10-19`, `internal/permission/permission.go:74-108`). | medium |

## Relevant Patterns

- **Thin Runtime Boundary Behind Thin Commands:** Keep CLI health wiring as a caller of runtime capability, not the owner of runtime behavior. Evidence across project-structure and command reports favors thin entrypoints and command wrappers (`cmd/gh/main.go:6`, `cmd/root.go:112-127`, `pkg/cmd/install.go:132-145`). For this sprint, this pattern pressures `ultraplan health` to report runtime state without directly launching OpenCode or embedding study execution.

- **Manual Composition Root With Explicit Wrapper Stack:** Mature Go CLIs wire dependencies in a visible constructor or app root rather than using hidden globals or DI frameworks. gh-cli uses `cmdutil.Factory` lazy dependencies (`pkg/cmdutil/factory.go:16-43`), opencode uses `app.New` as a composition root (`internal/app/app.go:42-81`), and go-task exposes functional options (`executor.go:22-24`). For this sprint, reasoning should decide where the agentwrap/opencode stack is built so wrapper order, stores, sinks, validation, policy, and fake runtimes remain inspectable.

- **Decorator/Wrapper For Cross-Cutting Runtime Concerns:** Decorators are common for logging, retry, debug, and UI safety. Chezmoi wraps systems for debug/dry-run, mitchellh-cli wraps UI, and restic wraps backends (`config.go:2314-2316`, `ui_concurrent.go:9-12`, `internal/backend/logger/log.go:22-77`). Agentwrap already supplies runtime wrappers, so the sprint pressure is to compose wrappers rather than duplicate their internals.

- **Merged Config Then Validate Once:** Configuration reports repeatedly favor a clear precedence chain followed by validation after all sources are merged. Examples include chezmoi flag restore (`internal/cmd/config.go:2253-2287`), go-task generic config precedence (`internal/flags/flags.go:314-327`), and opencode centralized validation (`internal/config/config.go:609-641`). Runtime config should follow this pressure: resolve executable/provider/model/policy/health/permissions before constructing a runtime, then reject unsupported values before launch.

- **Typed Error Classification Over String Matching:** Error studies favor `%w` chains, sentinels, and typed errors that callers inspect through `errors.Is`/`errors.As`. Rclone's `Retrier` interface (`fs/fserrors/error.go:22-29`), restic's `fatalError` (`internal/errors/fatal.go:10`), and gh-cli exit classes (`internal/ghcmd/cmd.go:44-49`) support this pattern. Runtime integration should preserve agentwrap classifications instead of converting them into unstructured messages.

- **Injectable IO And Fake Runtime Seams:** Testable CLIs inject stdout/stderr, filesystems, and external services. Gh-cli's `IOStreams.Test` returns buffers for tests (`pkg/iostreams/iostreams.go:551-568`), restic exposes `Terminal` and backend interfaces (`internal/ui/terminal.go:10-36`, `internal/backend/backend.go:19-90`), and go-task's `WithStdout` configures output (`executor.go:553-564`). This applies to runtime health output, runtime fakes, event sinks, and config redaction tests.

- **Context As Lifecycle, Not Service Locator:** State/context evidence supports root contexts and cancellation propagation, while warning against service-locator usage. gh-cli uses `ExecuteContextC` (`internal/ghcmd/cmd.go:194`), helm wires SIGTERM to cancellation (`pkg/cmd/install.go:333-347`), and k9s context keys are flagged as a service-locator risk (`internal/keys.go:10-38`). Runtime calls should receive context directly and avoid `context.Background()` in launch/wait paths.

- **Structured Events And Logs For Operator Diagnostics:** Observability evidence favors structured logs and pub/sub/event fan-out for runtime inspection. Opencode's pub/sub broker (`internal/pubsub/broker.go:10-19`) and slog wrapper (`internal/logging/logger.go:25-62`) show event/log separation. This sprint should map agentwrap event kinds and safe metadata into UltraPlan diagnostics without persisting unsafe native payload bytes by default.

- **Permission And Trust Boundaries As First-Class Runtime Inputs:** Security evidence from opencode shows interactive permission requests and command filters as explicit trust boundaries (`internal/permission/permission.go:44-108`, `internal/llm/tools/bash.go:41-55`). For UltraPlan, this translates into explicit permission policy mapping and fail-fast behavior when agentwrap/opencode cannot represent required policy features.

- **Behavioral Tests Over Implementation Detail Tests:** Testing evidence favors testscript, golden output, centralized mocks, and behavior assertions. Chezmoi's txtar test harness (`internal/cmd/main_test.go:64-174`), gh-cli HTTP mocks (`pkg/httpmock/stub.go:35-199`), and restic functional backend mocks (`internal/backend/mock/backend.go:14-26`) point to fake-runtime coverage that asserts observable results, events, diagnostics, and exit behavior rather than private wrapper internals.

## Trade-Offs

| Trade-Off | Benefit | Cost | When It Matters |
| --- | --- | --- | --- |
| Thin UltraPlan runtime model vs exposing raw agentwrap everywhere | Keeps platform/runtime generic and prevents study/product leakage; matches thin-boundary evidence from `cmd/` to interior packages (`cmd/root.go:112-127`, `internal/restic/repository.go:18`) | Requires mapping types, metadata, and errors carefully; may hide agentwrap details needed by later features | When deciding how much of `agentwrap.RunRequest`, `RunResult`, events, health, and metadata to mirror in UltraPlan internal types |
| Fixed wrapper composition vs task-specific wrapper order | Predictable diagnostics and tests; decorator evidence favors explicit wrapper chains (`ui_concurrent.go:9-12`, `internal/backend/logger/log.go:22-77`) | Validation/fallback interactions may require exceptions, and order changes can be subtle | When reasoning whether validation failures should participate in policy retry/fallback or only fail after policy completes |
| Centralized runtime config validation vs distributed adapter checks | Single failure point with actionable field errors; supported by opencode and k9s validation patterns (`internal/config/config.go:609-641`, `internal/config/json/validator.go:146`) | Central validator can grow if it absorbs adapter semantics; risks duplicating agentwrap capability logic | When mapping required health checks, capabilities, permission features, executable, environment, timeout, retry, fallback, provider, and model |
| Agentwrap-owned retry/fallback vs UltraPlan-owned scheduling policy | Avoids duplicating retry engines and keeps runtime retry classification near `SDKError`/policy metadata; aligns with decorator and typed-error evidence (`fs/fserrors/error.go:22-29`) | UltraPlan must translate policy events/metadata cleanly and may need to preserve enough state for later run-loop features | When deciding which retry/fallback facts are stored in runtime result metadata and which remain agentwrap internals |
| Rich event persistence vs safe minimal diagnostics | Rich events support debugging and future run inspection; pub/sub evidence shows value in event fan-out (`internal/pubsub/broker.go:10-19`) | Persisting native payloads or stderr can leak secrets and inflate artifacts | When deciding whether to store raw payload presence, omission reason, event kind, usage, permissions, validation, and policy decisions |
| Config-driven permission posture vs hardcoded safe defaults | Config supports different runtime contexts and future task kinds; security evidence supports explicit trust boundaries (`internal/permission/permission.go:44-108`) | More config fields increase validation and user explanation burden; philosophy evidence warns that config growth becomes maintenance debt | When deciding the minimal runtime config additions needed for executable, sandbox, permissions, env, policy, health, and capabilities |
| Fake-only default tests vs gated real OpenCode smoke tests | Fake tests are deterministic and do not require credentials; testing reports favor mocks and fixtures (`internal/backend/mock/backend.go:14-26`, `pkg/httpmock/stub.go:35-199`) | Fake behavior can diverge from real adapter behavior if not aligned with agentwrap APIs | When designing runtime unit tests and optional smoke coverage through `agentwrap/opencode` |

## Anti-Patterns And Warnings

- **Do Not Reimplement The Runtime Supervisor:** Evidence favors delegating volatile external behavior to clear boundaries and wrappers. Implementing OpenCode process launch, stdout decoding, retry, validation, permission translation, or native event parsing in UltraPlan would recreate the logic agentwrap is meant to own; opencode's own command orchestration was cited as a monolithic risk (`cmd/root.go:49-183`).

- **Avoid String Matching Runtime Errors:** The error report repeatedly flags missing `%w`, no sentinels, and stringly errors as weak patterns (`mitchellh-cli/cli.go:205-206`, `k9s/internal/client/errors.go:9-14`). Runtime error handling should use agentwrap typed/classified errors and preserve chains.

- **Avoid Global Runtime Config Or Service Singletons:** DI and config reports identify global config/service singletons as hidden coupling and test pollution (`fs/config.go:793`, `pkg/yqlib/lib.go:13`, `internal/config/config.go:123`). Runtime construction should receive resolved config and fakes through explicit seams.

- **Do Not Let CLI Health Invoke Study Workflows:** Command and architecture evidence supports command wrappers that delegate narrowly (`pkg/cmd/install.go:132-145`). Health should call platform runtime health/capability checks only, not prompt composition, report validation, run-state mutation, or scheduling.

- **Do Not Drop Unknown Events Or Unsafe Payload Facts:** Observability and runtime evidence pressure event handling to preserve event kind and diagnostic presence while omitting unsafe bytes. Unknown native events should be represented safely rather than treated as fatal only because the adapter projected a future type.

- **Do Not Convert Unknown Usage/Cost To Zero:** Summary and observability evidence repeatedly treats missing values distinctly from zero values; runtime metadata should keep token/cost unknownness explicit rather than defaulting absent values to numeric zero.

- **Avoid Context Background In Work Paths:** Context reports flag `context.Background()` inside long operations as misleading or uncancellable (`internal/cmd/templatefuncs.go:215`, `task.go:341`). Runtime request start, event drain, wait, cancellation, and health checks should propagate caller context.

- **Avoid Test Assertions On Private Wrapper Implementation Details:** Testing evidence warns against brittle internal assertions such as hint counts (`internal/view/pod_test.go:23`). Tests should assert request mapping, wrapper-observable behavior, result metadata, safe diagnostics, event mapping, and health output.

- **Avoid Permission Best-Effort When Requirements Are Strict:** Security reports show that explicit trust boundary failures are preferable to silent weakening (`internal/config/plugin.go:158-164` warning-only schema failures). If a required permission policy cannot map to agentwrap/opencode, reasoning should choose fail-fast behavior unless explicitly classified as best-effort.

## Examples Worth Inspecting

| Example | Path / Source | Why It Is Useful |
| --- | --- | --- |
| gh-cli dependency factory | `pkg/cmdutil/factory.go:16-43` | Shows explicit command dependency injection with lazy functions and test seams. |
| opencode app composition | `internal/app/app.go:42-81` | Shows a visible composition root for services that can inform runtime stack construction, while its global config concern is a caution. |
| opencode pub/sub broker | `internal/pubsub/broker.go:10-19` | Useful for event fan-out and runtime event sink reasoning. |
| opencode permission service | `internal/permission/permission.go:44-108` | Relevant to explicit permission requests and blocking decisions around agent tools. |
| opencode bash allow/deny policy | `internal/llm/tools/bash.go:41-55` | Useful evidence for permission posture mapping and command safety categories. |
| go-task functional options | `executor.go:22-24`, `executor.go:553-564` | Shows configurable construction and test IO injection without expanding constructors indefinitely. |
| restic backend and terminal abstractions | `internal/backend/backend.go:19-90`, `internal/ui/terminal.go:10-36` | Shows strong external-boundary interfaces and testable UI/IO seams. |
| rclone behavioral errors | `fs/fserrors/error.go:22-29`, `fs/fserrors/error.go:92-99` | Useful for retry/fatal semantics without string matching. |
| helm signal-context wiring | `pkg/cmd/install.go:333-347` | Relevant to cancellation propagation into external operations. |
| chezmoi flag restore | `internal/cmd/config.go:2253-2287` | Useful when reasoning config precedence where flags should override config file values. |
| k9s config schema validation | `internal/config/json/validator.go:146` | Useful for fail-fast validation of user-authored runtime config. |
| gh-cli IO test helper | `pkg/iostreams/iostreams.go:551-568` | Good model for testing health/runtime output with buffers. |
| restic functional mock backend | `internal/backend/mock/backend.go:14-26` | Useful for fake-runtime design that records behavior without real subprocesses. |
| helm subprocess plugin runtime | `internal/plugin/runtime_subprocess.go:65-79` | Useful as an example of external process isolation, though this sprint should rely on agentwrap rather than copy subprocess handling. |

## Design Pressures

- Runtime behavior must stay generic and platform-owned while study/report semantics remain outside the platform runtime boundary.

- Agentwrap is the volatile external runtime boundary; UltraPlan should translate requests, results, health, errors, permissions, policy, validation, and events without competing with the SDK.

- Health output must be actionable and testable, but it must not become a hidden study execution path.

- Runtime config needs enough expressive power for OpenCode/agentwrap options while resisting broad plugin/registry/config growth.

- Validation must make runtime success insufficient for product success, but the platform runtime should expose generic validator hooks rather than owning report rules.

- Events and diagnostics must preserve enough detail for later run inspection while redacting secrets and omitting unsafe raw native payload bytes by default.

- Context and cancellation need to cross every runtime operation, including health checks, event drains, waits, timeouts, and cancellation paths.

- Tests must prove behavior through fakes and buffers by default, with real OpenCode smoke coverage optional and gated.

- Permission posture must be explicit, representable through agentwrap/opencode, and fail fast when a required policy cannot be honored.

- Runtime metadata must align with later run-state fields without importing product modules or embedding study types.

## Open Questions For Reasoning

- What exact UltraPlan-internal runtime types should exist so the platform runtime is generic but still maps agentwrap request/result/health/event concepts without losing classifications?

- Where should the production agentwrap stack be constructed so wrapper order is explicit, testable, and not coupled to CLI health wiring?

- Should validation failures run inside policy retry/fallback decisions, or should policy operate only on underlying runtime failures with validation applied afterward?

- Which health and capability names are required in config for this sprint, and which unsupported names should be rejected before launch?

- How should permission requirements be represented when agentwrap/opencode can only support a subset of path/tool/shell policy features?

- What safe event and diagnostic fields should be persisted or logged now, and which raw/native fields should only record presence and omission reason?

- How should unknown usage, token, and cost fields be represented so they remain distinct from zero?

- What is the minimal fake runtime surface needed to test success, validation failure, runtime failure, timeout, cancellation, malformed event, rate limit, retry, fallback, permission failure, and observability metadata?

- Should any runtime config environment values be accepted directly, or should config store only safe env additions and rely on ambient runtime/provider configuration?

- What exact health command UX is needed for runtime checks without introducing a stable new public runtime JSON schema beyond existing output conventions?

## Evidence Pointers

- `studies/go-cli-study/reports/final/01-project-structure.md`: Inspect thin entrypoint, `internal/` boundary, and one-way dependency flow evidence around `cmd/root.go:112-127`, `internal/restic/repository.go:18`, and `cmd/root.go:255-259`.

- `studies/go-cli-study/reports/final/02-command-architecture.md`: Inspect factory command creation, lifecycle hooks, and monolithic `RunE` warnings around `pkg/cmdutil/factory.go:16-43`, `pkg/cmd/install.go:132-145`, and `cmd/root.go:49-183`.

- `studies/go-cli-study/reports/final/03-dependency-injection.md`: Inspect composition root, functional options, decorators, and global-state warnings around `internal/app/app.go:42-81`, `executor.go:22-24`, `config.go:2314-2316`, and `fs/config.go:793`.

- `studies/go-cli-study/reports/final/04-configuration-management.md`: Inspect precedence and validation evidence around `internal/cmd/config.go:2253-2287`, `internal/flags/flags.go:314-327`, `internal/config/config.go:609-641`, and `restic/internal/global/global.go:139,147`.

- `studies/go-cli-study/reports/final/05-error-handling.md`: Inspect wrapping, sentinels, behavioral error interfaces, and warning examples around `fs/fserrors/error.go:22-29`, `internal/errors/fatal.go:10`, `pkg/storage/driver/driver.go:27-36`, and `mitchellh-cli/cli.go:205-206`.

- `studies/go-cli-study/reports/final/06-io-abstraction.md`: Inspect IO and fake/test seam evidence around `pkg/iostreams/iostreams.go:551-568`, `internal/ui/terminal.go:10-36`, `internal/backend/backend.go:19-90`, and `executor.go:553-564`.

- `studies/go-cli-study/reports/final/07-state-context.md`: Inspect context propagation, cancellation, session modeling, and `Background()` warnings around `internal/ghcmd/cmd.go:142`, `pkg/cmd/install.go:333-347`, `internal/session/session.go:12-23`, and `task.go:341`.

- `studies/go-cli-study/reports/final/08-concurrency.md`: Inspect structured concurrency, bounded fan-out, event coordination, and wait timeout evidence around `task.go:87`, `pkg/analyze/parallel.go:13`, `internal/pubsub/broker.go:10`, and `cmd/root.go:261-279`.

- `studies/go-cli-study/reports/final/10-logging-observability.md`: Inspect structured logging, debug controls, stdout/stderr separation, and pub/sub evidence around `internal/logging/logger.go:25-62`, `internal/pubsub/broker.go:10-19`, `pkg/iostreams/iostreams.go:52-54`, and `src/core.go:325`.

- `studies/go-cli-study/reports/final/11-testing-strategy.md`: Inspect fake/mock, golden, integration, and sparse coverage evidence around `internal/cmd/main_test.go:64-174`, `pkg/httpmock/stub.go:35-199`, `internal/backend/mock/backend.go:14-26`, and `opencode/internal/tui/theme/theme_test.go:1-89`.

- `studies/go-cli-study/reports/final/12-extensibility.md`: Inspect extension boundary and registry evidence around `internal/llm/tools/tools.go:69-72`, `internal/llm/agent/mcp-tools.go:106-129`, `helm/internal/plugin/runtime_subprocess.go:65-79`, and `rclone/fs/registry.go:407`.

- `studies/go-cli-study/reports/final/13-security.md`: Inspect permission, command filtering, secret redaction, and validation evidence around `internal/permission/permission.go:44-108`, `internal/llm/tools/bash.go:41-55`, `internal/options/secret_string.go:15-20`, and `pkg/registry/transport.go:37-41`.

- `studies/go-cli-study/reports/final/15-philosophy.md`: Inspect complexity discipline, explicit non-goals, and factory/interface/pubsub philosophy around `pkg/cmd/factory/default.go:26-46`, `internal/pubsub/broker.go:10-19`, `VISION.md:97`, and `internal/permission/permission.go:74-108`.

## Handoff To Reasoning

- Use this handbook as evidence input.
- Validate whether the observed patterns fit this project's constraints.
- Do not copy external patterns without sprint-specific reasoning.
