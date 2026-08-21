# Sprint Technical Handbook: Local Web Hardening and Observable-Product Release

> Project: `ultraplan-go`
> Sprint: `32-hardening-and-release`
> Source: `projects/ultraplan-go/sprints/32-hardening-and-release/sprint-index.md`
> **Inputs Used:** `projects/ultraplan-go/project-index.md`, `projects/ultraplan-go/sprints/32-hardening-and-release/sprint-index.md`, `projects/ultraplan-go/sprints/32-hardening-and-release/requirements.md`, `studies/go-cli-study/reports/final/01-project-structure.md`, `studies/go-cli-study/reports/final/02-command-architecture.md`, `studies/go-cli-study/reports/final/03-dependency-injection.md`, `studies/go-cli-study/reports/final/04-configuration-management.md`, `studies/go-cli-study/reports/final/05-error-handling.md`, `studies/go-cli-study/reports/final/06-io-abstraction.md`, `studies/go-cli-study/reports/final/07-state-context.md`, `studies/go-cli-study/reports/final/08-concurrency.md`, `studies/go-cli-study/reports/final/10-logging-observability.md`, `studies/go-cli-study/reports/final/11-testing-strategy.md`, `studies/go-cli-study/reports/final/12-extensibility.md`, `studies/go-cli-study/reports/final/13-security.md`, `studies/go-cli-study/reports/final/14-performance.md`

This handbook distills the studies and reports selected by `sprint-index.md` for sprint reasoning. It does not decide architecture or implementation.

## Contract Coverage Scope

The selected Go CLI study is supporting evidence, not a substitute for the
selected contracts. Where the study is silent, reasoning must apply the named
contract directly and record that fact.

| Selected contract | Study support | Direct-contract concerns still required |
| --- | --- | --- |
| Architecture | structure, command architecture, dependency injection, extensibility | product ownership, web/app boundary, presentation hierarchy |
| CLI Surface | command architecture, configuration, errors, I/O | stable `serve` lifecycle, help, exit and non-interactive behavior |
| Configuration | configuration management, security | field-level source reporting, schema rejection and redacted `config show` |
| Documentation | testing and command examples only | public-surface ownership, operational recovery and release reconciliation |
| Errors | error handling, I/O, observability | stable cross-transport codes and terminal/durable error agreement |
| LLM Evaluation / Cost / Safety | security and observability only | runtime/model/prompt identity, usage, retry, duration, tool/fallback metadata and safety gates |
| LLM Runtime | cancellation, concurrency and observability only | prompt versioning, structured-output validation, inspectable lifecycle and bounded retry classification |
| Observability | logging, state/context and errors | correlation fields, blocked-as-not-pass and projection-level redaction |
| Performance | performance and concurrency | explicit local latency/scale expectations and rejection at capacity |
| Persistence And Migrations | I/O and state/context only | atomic-write ownership, snapshots, schema ownership and conservative reconciliation/migration |
| Security | security, configuration and I/O | browser Host/Origin/CSRF/session policy and per-projection allowlists |
| Testing | testing, I/O and dependency injection | API fixtures, browser semantics, race/leak, packaging and gated real dependencies |
| Workflows | state/context and concurrency | definition compatibility, exact-once cancellation, durable terminal arbitration and recovery |

## Selected Studies And Reports

| Study / Report | Path | Relevant Finding | Confidence |
| --- | --- | --- | --- |
| Project structure | `studies/go-cli-study/reports/final/01-project-structure.md` | Successful applications keep interface entrypoints thin and dependencies unidirectional; chezmoi, yq, helm, and restic all delegate inward without reverse imports (`chezmoi/main.go:26-34`, `yq/cmd/root.go:9`, `helm/pkg/action/install.go:73-140`, `restic/internal/restic/repository.go:18`). | high |
| Command architecture | `studies/go-cli-study/reports/final/02-command-architecture.md` | Shared behavior is most reusable when transport wiring delegates to actions or factories; helm's command calls an action and gh-cli centralizes dependencies in a factory (`helm/pkg/cmd/install.go:132-145`, `gh-cli/pkg/cmdutil/factory.go:16-43`). | high |
| Dependency injection | `studies/go-cli-study/reports/final/03-dependency-injection.md` | Explicit manual composition, narrow service interfaces, and constructor injection correlate with testability; package globals and context service locators obscure dependencies (`gh-cli/pkg/cmd/factory/default.go:26-46`, `restic/internal/backend/backend.go:19-90`, `k9s/internal/view/xray.go:412`). | high |
| Configuration management | `studies/go-cli-study/reports/final/04-configuration-management.md` | Configuration is reliable when source precedence is explicit and the merged result is validated once (`chezmoi/internal/cmd/config.go:2253-2287`, `k9s/internal/config/k9s.go:423-451`). | high |
| Error handling | `studies/go-cli-study/reports/final/05-error-handling.md` | Wrapped, typed, or sentinel errors preserve machine classification while a separate rendering path presents safe, actionable messages (`helm/pkg/storage/driver/driver.go:27-48`, `restic/internal/errors/fatal.go:10-53`, `gh-cli/internal/ghcmd/cmd.go:281-301`). | high |
| I/O abstraction | `studies/go-cli-study/reports/final/06-io-abstraction.md` | Injected streams and boundary interfaces permit deterministic tests without real terminals, filesystems, or networks (`gh-cli/pkg/iostreams/iostreams.go:551-568`, `restic/internal/fs/interface.go:10-31`, `chezmoi/internal/cmd/applycmd_test.go:220-241`). | high |
| State and context | `studies/go-cli-study/reports/final/07-state-context.md` | Long-running work needs one propagated cancellation lineage; cleanup may need a deliberately separate bounded context (`helm/pkg/cmd/install.go:333-347`, `restic/internal/restic/lock.go:290-305`). | high |
| Concurrency | `studies/go-cli-study/reports/final/08-concurrency.md` | Auditable launch sites, bounded fan-out, explicit waits, and one-time cleanup reduce leaks and shutdown ambiguity (`go-task/task.go:87`, `k9s/internal/pool.go:21-37`, `opencode/cmd/root.go:252-279`, `rclone/lib/batcher/batcher.go:50`). | high |
| Logging and observability | `studies/go-cli-study/reports/final/10-logging-observability.md` | Structured fields, runtime debug control, and separation of product output from diagnostics make failures inspectable without corrupting machine output (`helm/internal/logging/logging.go:31-71`, `k9s/internal/slogs/keys.go:6-231`, `gh-cli/pkg/iostreams/iostreams.go:52-54`). | high |
| Testing strategy | `studies/go-cli-study/reports/final/11-testing-strategy.md` | Layered behavior tests, compatibility fixtures, centralized fakes, and real command pipelines protect public behavior better than implementation-detail assertions (`gh-cli/acceptance/acceptance_test.go:26-29`, `helm/internal/test/test.go:43`, `restic/internal/backend/mock/backend.go:14-26`). | high |
| Extensibility | `studies/go-cli-study/reports/final/12-extensibility.md` | Capability interfaces, factories, adapters, and optional-interface detection allow controlled growth without exposing internal orchestration (`gh-cli/pkg/cmdutil/factory.go:16-43`, `rclone/fs/features.go:294-370`, `dive/cmd/dive/cli/internal/command/adapter/analyzer.go:13-15`). | high |
| Security | `studies/go-cli-study/reports/final/13-security.md` | Explicit trust boundaries, structured validation, redacting types, safe argv construction, and bounded permission flows distinguish secure tools (`restic/internal/options/secret_string.go:15-20`, `k9s/internal/config/json/validator.go:146`, `opencode/internal/permission/permission.go:44-108`). | high |
| Performance | `studies/go-cli-study/reports/final/14-performance.md` | Long-lived interfaces remain responsive through lazy initialization, streaming, bounded queues, and explicit resource caps rather than speculative micro-optimization (`gh-cli/pkg/cmdutil/factory.go:27-42`, `lazygit/pkg/tasks/tasks.go:189-217`, `k9s/internal/pool.go:26-48`). | high |

## Relevant Patterns

- **Transport adapter over shared capabilities:** The dominant structure is a thin boundary that parses and maps transport concerns, then delegates to an inner action or service. Helm separates command construction from `pkg/action` execution (`helm/pkg/cmd/install.go:132-145`, `helm/pkg/action/install.go:73-140`), while restic defines the domain interface inward of its implementation (`restic/internal/restic/repository.go:18`, `restic/internal/repository/repository.go:1`). For this sprint, that evidence supports examining whether HTTP, HTML, SSE, CLI, and TUI remain projections over shared capabilities rather than independent workflow owners.
- **Explicit composition root with narrow injectable seams:** Gh-cli assembles lazy dependencies through factory functions (`gh-cli/pkg/cmd/factory/default.go:26-46`, `gh-cli/pkg/cmdutil/factory.go:16-43`), while restic supplies replaceable backend and repository interfaces (`restic/internal/backend/backend.go:19-90`, `restic/internal/restic/repository.go:18-66`). This pattern addresses capability testing and controlled substitution, but the reports caution against a central god object such as chezmoi's large config (`chezmoi/internal/cmd/config.go:193-291`).
- **Typed error classification separated from presentation:** Helm combines sentinels and structured wrappers (`helm/pkg/storage/driver/driver.go:27-48`), and restic identifies fatal conditions through an error type before rendering (`restic/internal/errors/fatal.go:10-53`, `restic/cmd/restic/main.go:199-209`). Gh-cli adds user guidance only at the outer renderer (`gh-cli/internal/ghcmd/cmd.go:281-301`). This is relevant to stable API error codes, HTTP status mapping, safe HTML messages, and SSE terminal outcomes without making presentation strings the source of truth.
- **One lifecycle context for work, a distinct bounded context for cleanup:** Helm wires signals into the same context consumed by long-running actions (`helm/pkg/cmd/install.go:333-347`, `helm/pkg/action/install.go:284`), while restic deliberately delays cancellation for lock cleanup (`restic/internal/restic/lock.go:290-305`). The evidence creates a useful distinction between cancellation of product work and bounded final reconciliation; reasoning must determine ownership and ordering for each.
- **Bounded structured concurrency and exact-once teardown:** Go-task uses `errgroup.WithContext` for sibling cancellation (`go-task/task.go:87`), k9s bounds workers with a semaphore (`k9s/internal/pool.go:21-37`), opencode waits for cleanup with a timeout (`opencode/cmd/root.go:252-279`), and rclone uses `sync.Once` for one-time completion (`rclone/lib/batcher/batcher.go:50`). Together these patterns address operation limits, slow subscribers, shutdown deadlines, and exact-once cancellation without prescribing a specific primitive.
- **Streaming with bounded retention and explicit slow-consumer behavior:** Lazygit emits command output incrementally through scanner and channel boundaries (`lazygit/pkg/tasks/tasks.go:189-217`), and opencode streams provider events while accepting that slow consumers may lose events (`opencode/internal/llm/provider/provider.go:56`). Rclone caps pool resources (`rclone/lib/pool/pool.go:17-24,52-53`). This evidence is pertinent to SSE replay windows, gaps, subscriber limits, and terminal delivery, where unbounded fidelity conflicts with bounded server resources.
- **Explicit merged configuration followed by validation:** Chezmoi preserves explicit CLI overrides across config loading (`chezmoi/internal/cmd/config.go:2253-2287`), restic tracks whether a flag was actually changed (`restic/internal/global/global.go:139,147`), and k9s validates the resulting configuration centrally (`k9s/internal/config/k9s.go:423-451`). This pattern matters for documented bind, timeout, stream, retention, and body limits and for ensuring defaults do not silently shadow operator configuration.
- **Redaction at the type and boundary levels:** Restic's `SecretString` makes accidental formatting safe (`restic/internal/options/secret_string.go:15-20`), while helm strips authorization material at its transport logging boundary (`helm/pkg/registry/transport.go:37-41`). The combined pattern is stronger than ad hoc string cleanup and applies across JSON, HTML, SSE, retained events, terminal results, and diagnostics.
- **Layered tests around stable behavior:** Gh-cli uses full acceptance pipelines (`gh-cli/acceptance/acceptance_test.go:26-29`), helm centralizes golden comparison (`helm/internal/test/test.go:43`), and restic provides functional-field fakes (`restic/internal/backend/mock/backend.go:14-26`). Gh-cli and chezmoi also use `httptest` (`gh-cli/pkg/cmd/issue/list/list_test.go:29-31`, `chezmoi/internal/cmd/applycmd_test.go:220-241`). The sprint can draw on this combination for API compatibility, transport behavior, fake-runtime integration, and gated real-system evidence.
- **Capability discovery rather than route-specific branching:** Rclone detects optional backend abilities through interfaces (`rclone/fs/features.go:294-370`), while gh-cli and dive keep consumers behind factories or adapters (`gh-cli/pkg/cmdutil/factory.go:16-43`, `dive/cmd/dive/cli/internal/command/adapter/analyzer.go:13-15`). This pattern directly pressures the proposed future-stage proof: expose a bounded capability vocabulary without moving stage orchestration into web routes.

## Trade-Offs

| Trade-Off | Benefit | Cost | When It Matters |
| --- | --- | --- | --- |
| Narrow interfaces versus concrete app types | Interfaces enable fake stages, transport-independent tests, and decorators, as in restic's backend contracts (`restic/internal/backend/backend.go:19-90`). | Too many or overly broad interfaces obscure the concrete graph and can become god interfaces; gdu's UI abstraction warning illustrates this pressure (`gdu/cmd/gdu/app/app.go:30-49`). | Defining the shared web and operation capability surface. |
| One central composition object versus grouped dependencies | A composition root makes construction traceable (`gh-cli/pkg/cmd/factory/default.go:26-46`). | Large containers couple unrelated consumers and become difficult to evolve (`chezmoi/internal/cmd/config.go:193-291`). | Wiring the HTTP server, operation hub, runtime/harness fakes, clocks, IDs, and diagnostics. |
| Structured error types versus a small fixed error vocabulary | Typed errors retain fields and causes for mapping and recovery (`go-task/errors/errors_task.go:13-32`). | A large hierarchy increases compatibility and maintenance obligations; sentinels may suffice for boolean conditions. | Freezing API error codes and projecting errors consistently into HTML, JSON, and SSE. |
| Parent cancellation versus independent cleanup context | Parent context propagation makes cancellation prompt and predictable (`helm/pkg/cmd/install.go:333-347`). | Cleanup tied to the parent can be cancelled before durable reconciliation; a detached cleanup context can itself outlive intended ownership (`restic/internal/restic/lock.go:290-305`). | Shutdown, cancellation, process cleanup, lock release, and `cleanup_uncertain` persistence. |
| `errgroup` fail-fast versus independent `WaitGroup` completion | `errgroup` collects failures and cancels siblings (`go-task/task.go:87`). | Fail-fast cancellation is wrong when independent subscribers or cleanup tasks must all reach their own terminal handling; `WaitGroup` requires manual error arbitration. | Concurrent operations, terminal arbitration, subscriber delivery, and shutdown fan-out. |
| Complete replay versus bounded event retention | Replay improves reconnect continuity; streaming avoids loading full histories (`lazygit/pkg/tasks/tasks.go:189-217`). | Complete retention grows without bound, while bounded retention creates explicit gaps and fallback work; opencode accepts dropped events for slow consumers (`opencode/internal/llm/provider/provider.go:56`). | SSE event IDs, reconnect, rollover, slow clients, and durable refresh. |
| Eager validation versus lazy initialization | Early validation fails before serving malformed limits; lazy factories avoid unnecessary startup work (`gh-cli/pkg/cmdutil/factory.go:27-42`). | Eager loading slows simple lifecycle paths; lazy loading defers errors into requests or operations. | Server startup validation, embedded asset parsing, runtime/harness clients, and cheap status routes. |
| Golden compatibility fixtures versus semantic assertions | Golden fixtures expose public shape changes clearly (`helm/internal/test/test.go:43`). | Large snapshots can be mechanically updated and may obscure whether a change is intentional; implementation-detail assertions are brittle (`k9s/internal/view/pod_test.go:23`). | Stable `/api/v1` envelopes, route matrices, rendered HTML, and documentation examples. |
| Rich observability versus secret and payload minimization | Structured component fields improve diagnosis (`k9s/internal/slogs/keys.go:6-231`). | More fields increase leakage and retention risk; debug output can corrupt user output if routing is weak (`fzf/src/core.go:325`). | Operation events, shutdown diagnostics, runtime metadata, and release evidence. |
| Streaming versus full buffering | Streaming bounds memory and improves responsiveness (`age/internal/stream/stream.go:195-219`). | It complicates ordering, terminal flush, replay, and gap handling; buffering is simpler but scales with payload size. | SSE, operation results, artifact rendering, and large diagnostics. |

## Anti-Patterns And Warnings

- **Route or template ownership of workflow rules:** Large command handlers become untestable and non-reusable (`opencode/cmd/root.go:49-183`, `yq/cmd/evaluate_sequence_command.go:152`). Equivalent workflow branching in HTTP handlers or templates would undermine the shared capability goal.
- **Hidden globals or context as a service locator:** Rclone's nil-context fallback masks missing configuration (`rclone/fs/config.go:793`), and k9s retrieves services with unchecked context assertions (`k9s/internal/dao/ds.go:72`). Either pattern can create cross-test contamination and runtime-only failures in a concurrent server.
- **Fresh background contexts inside cancellable work:** Chezmoi template functions detach work with `context.Background()` (`chezmoi/internal/cmd/templatefuncs.go:215`), and go-task deferred work does likewise (`go-task/task.go:341`). Detachment without an explicit cleanup ownership contract risks orphaned product work or premature lock release.
- **Fire-and-forget goroutines and unbounded fan-out:** Dive launches an untracked notification goroutine (`dive/cmd/dive/cli/internal/command/adapter/resolver.go:70`), while gh-cli launches one goroutine per extension and waits without a timeout (`gh-cli/pkg/cmd/extension/manager.go:196-206`). Both conflict with leak-free shutdown and bounded operation requirements.
- **Direct output or logging bypasses:** Chezmoi template functions write directly to stderr (`chezmoi/internal/cmd/templatefuncs.go:296`), and urfave-cli tracing bypasses its injectable error writer (`urfave-cli/cli.go:46-47`). Any bypass can evade capture, redaction, response-shape control, and compatibility tests.
- **String-only error and state interpretation:** Mitchellh-cli loses error-chain inspection (`mitchellh-cli/cli.go:205-206`), while k9s uses stringly typed errors (`k9s/internal/client/errors.go:9-14`). HTTP status, stable error code, terminal state, and recovery guidance should not depend on parsing display strings.
- **Silent or partial configuration failure:** Gdu logs config errors and continues (`gdu/cmd/gdu/main.go:242-244`), and direct environment access can bypass a configured precedence model (`opencode/internal/config/config.go:163`). For security and resource bounds, silent fallback can turn an invalid operator intent into unsafe server behavior.
- **Warnings-only schema or registry failures:** K9s continues after plugin-schema errors (`k9s/internal/config/plugin.go:158-164`), and rclone registries can silently overwrite names (`rclone/fs/rc/registry.go:41-48`). Analogous duplicate routes, template definitions, operation kinds, or capability registrations should not become load-order-dependent behavior.
- **Unredacted or default-allow trust boundaries:** Gh-cli can fall back to plaintext credentials (`gh-cli/internal/config/config.go:368-372`), while yq enables file and environment operations by default (`yq/pkg/yqlib/security_prefs.go:3-7`). The local-only server still needs explicit fail-closed decisions because browser input, paths, Markdown, and runtime output are untrusted.
- **Unbounded memory and wait paths:** Full accumulation risks memory growth (`yq/pkg/yqlib/stream_evaluator.go:78-113` contrasts stream mode), and bare `WaitGroup.Wait()` can hang forever (`gh-cli/pkg/cmd/extension/manager.go:205`). Event retention, payloads, streams, subscribers, results, cleanup, and shutdown waits all need inspectable bounds.
- **Tests coupled to internals rather than contracts:** K9s's hint-count assertion (`k9s/internal/view/pod_test.go:23`) and fragile regex output checks (`gh-cli/pkg/cmd/issue/list/list_test.go:88-91`) show how tests can fail for irrelevant changes or miss shape regressions. Compatibility and accessibility tests need behavior-level assertions with intentional fixtures.
- **Speculative pooling or concurrency:** The performance report recommends pooling only after profiling (`restic/internal/archiver/buffer.go:24-46`) and shows sequential designs can avoid races entirely (`urfave-cli/command_run.go:92`). More concurrency or caching is not automatically safer or faster for a loopback interface.

## Examples Worth Inspecting

| Example | Path / Source | Why It Is Useful |
| --- | --- | --- |
| Thin transport-to-action delegation | `studies/go-cli-study/reports/final/02-command-architecture.md` -> `helm/pkg/cmd/install.go:132-145`, `helm/pkg/action/install.go:73-140` | Shows a boundary mapping inputs and outputs while reusable logic remains inward. |
| Lazy explicit factory | `studies/go-cli-study/reports/final/03-dependency-injection.md` -> `gh-cli/pkg/cmd/factory/default.go:26-46`, `gh-cli/pkg/cmdutil/factory.go:16-43` | Shows test seams and deferred dependencies without a reflection-based DI framework. |
| Typed error plus outer rendering | `studies/go-cli-study/reports/final/05-error-handling.md` -> `helm/pkg/storage/driver/driver.go:27-48`, `gh-cli/internal/ghcmd/cmd.go:281-301` | Separates stable machine classification from actionable user text. |
| Cancellation and cleanup separation | `studies/go-cli-study/reports/final/07-state-context.md` -> `helm/pkg/cmd/install.go:333-347`, `restic/internal/restic/lock.go:290-305` | Contrasts work cancellation with cleanup that must complete under its own bound. |
| Bounded workers and timed shutdown | `studies/go-cli-study/reports/final/08-concurrency.md` -> `k9s/internal/pool.go:21-37`, `opencode/cmd/root.go:252-279` | Useful for operation, subscriber, and server-shutdown lifecycle reasoning. |
| Injectable test streams and HTTP | `studies/go-cli-study/reports/final/06-io-abstraction.md` -> `gh-cli/pkg/iostreams/iostreams.go:551-568`, `chezmoi/internal/cmd/applycmd_test.go:220-241` | Demonstrates deterministic boundary tests with buffers and `httptest`. |
| Structured diagnostics | `studies/go-cli-study/reports/final/10-logging-observability.md` -> `helm/internal/logging/logging.go:31-71`, `k9s/internal/slogs/keys.go:6-231` | Shows runtime-controlled levels, stderr routing, and consistent fields. |
| Secret-safe value and transport logging | `studies/go-cli-study/reports/final/13-security.md` -> `restic/internal/options/secret_string.go:15-20`, `helm/pkg/registry/transport.go:37-41` | Demonstrates defense against accidental formatting and boundary-specific credential leakage. |
| Acceptance plus golden plus fake layers | `studies/go-cli-study/reports/final/11-testing-strategy.md` -> `gh-cli/acceptance/acceptance_test.go:26-29`, `helm/internal/test/test.go:43`, `restic/internal/backend/mock/backend.go:14-26` | Provides complementary evidence for integration, compatibility, and deterministic failure-path tests. |
| Optional capability detection | `studies/go-cli-study/reports/final/12-extensibility.md` -> `rclone/fs/features.go:294-370` | Offers a reference for exposing new stage abilities without hard-coded route orchestration. |
| Streaming and bounded pools | `studies/go-cli-study/reports/final/14-performance.md` -> `lazygit/pkg/tasks/tasks.go:189-217`, `rclone/lib/pool/pool.go:17-24,52-53` | Useful for SSE flow, slow-consumer policy, and resource-cap reasoning. |

## Design Pressures

- The same product operation must remain semantically consistent across CLI, TUI, HTML, JSON, and SSE while each surface has different presentation and transport constraints.
- Stable `/api/v1` envelopes and error codes create a compatibility surface, but internal error causes and product state must remain evolvable.
- A browser disconnect is a subscriber-lifecycle event, not product cancellation; explicit cancellation and server shutdown have different authority.
- Shutdown needs both prompt cancellation and enough independent bounded time to persist a truthful terminal or cleanup-uncertain outcome.
- Event replay and terminal visibility compete with strict limits on retained events, payloads, subscribers, stream lifetime, and memory.
- The local boundary reduces exposure but does not remove Host, Origin, CSRF, path, hostile content, forged reference, and secret-leak risks.
- Diagnostics must be rich enough to explain lifecycle and recovery while excluding prompts, tokens, cookies, raw provider data, unsafe paths, stderr, and retained secrets.
- The capability model must demonstrate extensibility without becoming a generic plugin system, alternate scheduler, route registry with hidden side effects, or web-owned durable state.
- No-JavaScript completeness and progressive enhancement require durable truth to remain server/product-owned while enhanced operation views recover after refresh and SSE gaps.
- Compatibility, accessibility, race, leak, lifecycle, packaging, and gated real-system evidence need distinct test layers without making normal tests depend on live providers or harnesses.
- Lazy startup and streaming improve responsiveness, but startup must still fail closed on invalid templates, duplicate definitions, unsafe configuration, and impossible bounds.
- Release observability must distinguish blocked external prerequisites from passed checks and must not infer success from missing processes or the mere presence of artifacts.

## Open Questions For Reasoning

- What is the smallest shared capability vocabulary that covers status, artifacts, commands, progress, cancellation, and recovery without encoding stage-specific workflow rules?
- Which interfaces belong on the application side of the web boundary, and which concrete types are simpler and safer than adding another abstraction?
- Which error classifications and fields are stable API contract, and how are internal causes reduced into safe HTML, JSON, and SSE projections?
- How are HTTP status, API error code, operation terminal state, durable workflow state, and recovery guidance kept consistent without one presentation layer interpreting another?
- What context owns a started operation, what context owns server shutdown cleanup, and what exact bound applies after the operation's normal context is cancelled?
- Where is exact-once cancellation arbitrated when user cancellation, product failure, completion, and server shutdown race?
- What lock ordering permits cancellation, terminal arbitration, durable reconciliation, hub updates, and subscriber closure without waiting while holding hub locks?
- What are the hard limits for operations, preparations, retained events, event and result payloads, subscribers, concurrent streams, heartbeat, lifetime, polling, and cleanup waits?
- What replay-gap contract tells a browser to refresh durable state without treating an event gap as product failure?
- Which terminal information must be flushed to current subscribers, and what happens when a slow client cannot receive it within its stream bound?
- Which configuration sources control web limits and bind behavior, what is their exact precedence, and which invalid combinations fail startup?
- Does `config show` report each effective field's source without exposing secret values, and which source/redaction assertions freeze that behavior?
- Which values need secret-safe types, which values require explicit allowlists, and where must final projection-level redaction occur?
- Which runtime-backed browser projections expose runtime/provider/model and prompt version, token, retry, duration, tool, cost, and fallback metadata; which fields are unavailable or unsafe?
- Where are structured runtime outputs validated locally, and how are retry causes and limits represented in inspectable operation state?
- What API fixture format makes field additions, omissions, method changes, status changes, and stable error-code changes deliberate and reviewable?
- Which accessibility properties can be asserted deterministically from rendered HTML, and which still require manual keyboard, zoom, reflow, color, and reduced-motion evidence?
- How can the future-stage capability fixture prove extensibility without introducing production-only genericity or a second registration mechanism?
- What evidence distinguishes durable completion, durable cancellation, durable failure, and cleanup uncertainty after abrupt interruption or deadline exhaustion?
- Which product layer owns atomic writes, snapshots and schema versions, and how do migrations or restart reconciliation preserve recoverable prior state?
- Which tests use deterministic clocks, IDs, and barriers, and which timing behaviors need race/leak or gated integration coverage instead?
- What diagnostic fields are necessary for release support, and which are forbidden because they could reveal secrets, prompts, provider payloads, raw stderr, or unsafe paths?
- Which resources should initialize lazily, and which templates, routes, configuration, and security invariants must be validated before the listener serves requests?
- What release evidence is deterministic and mandatory everywhere, and what real-runtime or smoke-harness evidence can truthfully be blocked by unavailable prerequisites?

## Evidence Pointers

- `studies/go-cli-study/reports/final/01-project-structure.md` -> "Pattern 4: Unidirectional Dependency Flow" and `helm/pkg/action/install.go:73-140`: inspect boundary direction and inner action ownership.
- `studies/go-cli-study/reports/final/02-command-architecture.md` -> "Thin-Delegate Pattern" and `gh-cli/pkg/cmdutil/factory.go:16-43`: inspect delegation and shared construction.
- `studies/go-cli-study/reports/final/03-dependency-injection.md` -> "Centralized Composition Root" and `restic/internal/backend/backend.go:19-90`: inspect manual wiring and test seams, plus the god-object cautions around `chezmoi/internal/cmd/config.go:193-291`.
- `studies/go-cli-study/reports/final/04-configuration-management.md` -> "Explicit Three-Layer Systems" and `k9s/internal/config/k9s.go:423-451`: inspect precedence and merged validation.
- `studies/go-cli-study/reports/final/05-error-handling.md` -> "User/Operational Separation" and `helm/pkg/storage/driver/driver.go:27-48`: inspect safe rendering and machine classification.
- `studies/go-cli-study/reports/final/06-io-abstraction.md` -> "IOStreams with Test Constructor" and `gh-cli/pkg/iostreams/iostreams.go:551-568`: inspect deterministic boundary testing.
- `studies/go-cli-study/reports/final/07-state-context.md` -> "Signal-Context Wiring" and `restic/internal/restic/lock.go:290-305`: inspect work-versus-cleanup cancellation semantics.
- `studies/go-cli-study/reports/final/08-concurrency.md` -> "Deferred Cancel + Explicit Wait with Timeout" and `opencode/cmd/root.go:252-279`: inspect bounded shutdown and leak risks.
- `studies/go-cli-study/reports/final/10-logging-observability.md` -> "Structured Keys Package" and `helm/internal/logging/logging.go:31-71`: inspect consistent fields, levels, and output separation.
- `studies/go-cli-study/reports/final/11-testing-strategy.md` -> "testscript Integration Layer", "Golden File Regression Prevention", and "Centralized Mock Infrastructure": inspect complementary test layers and their maintenance costs.
- `studies/go-cli-study/reports/final/12-extensibility.md` -> "Optional Interface Detection" and `rclone/fs/features.go:294-370`: inspect capability discovery without route-specific orchestration.
- `studies/go-cli-study/reports/final/13-security.md` -> "Secret Redaction Type", "Credential Scrubbing in Logs", and `opencode/internal/permission/permission.go:44-108`: inspect type-level, boundary-level, and confirmation controls.
- `studies/go-cli-study/reports/final/14-performance.md` -> "Streaming via Channels and bufio" and "Concurrency Bounding via Semaphore Channels": inspect resource bounds before considering pooling or other optimization.

## Handoff To Reasoning

- Use this handbook as evidence input.
- Validate whether the observed patterns fit this project's constraints.
- Resolve the open questions in architecture, API-design, frontend, and sprint reasoning rather than treating report patterns as decisions.
- Preserve the distinction between ephemeral web operation state and authoritative product-owned durable state.
- Do not copy external patterns without sprint-specific reasoning.
