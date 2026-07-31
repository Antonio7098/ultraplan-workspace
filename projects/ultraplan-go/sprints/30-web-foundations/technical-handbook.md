# Sprint Technical Handbook: Local Web Foundation and Read-Only Dashboard

> Project: `ultraplan-go`
> Sprint: `30-web-foundations`
> Source: `projects/ultraplan-go/sprints/30-web-foundations/sprint-index.md`
> **Inputs Used:** `projects/ultraplan-go/project-index.md`, `projects/ultraplan-go/sprints/30-web-foundations/sprint-index.md`, `projects/ultraplan-go/sprints/30-web-foundations/requirements.md`, `projects/ultraplan-go/docs/ARCHITECTURE.md`, `projects/ultraplan-go/docs/PRD.md`, `projects/ultraplan-go/docs/TRD.md`, `studies/go-cli-study/reports/final/01-project-structure.md`, `studies/go-cli-study/reports/final/02-command-architecture.md`, `studies/go-cli-study/reports/final/03-dependency-injection.md`, `studies/go-cli-study/reports/final/04-configuration-management.md`, `studies/go-cli-study/reports/final/05-error-handling.md`, `studies/go-cli-study/reports/final/06-io-abstraction.md`, `studies/go-cli-study/reports/final/07-state-context.md`, `studies/go-cli-study/reports/final/08-concurrency.md`, `studies/go-cli-study/reports/final/10-logging-observability.md`, `studies/go-cli-study/reports/final/11-testing-strategy.md`, `studies/go-cli-study/reports/final/12-extensibility.md`, `studies/go-cli-study/reports/final/13-security.md`, `studies/go-cli-study/reports/final/14-performance.md`

This handbook distills the studies and reports selected by `sprint-index.md` for sprint reasoning. It does not decide architecture or implementation. The project documents above provide scope context only; the selected reports and the study source locations they cite are the evidence base for the findings below.

## Selected Studies And Reports

| Study / Report | Path | Relevant Finding | Confidence |
| --- | --- | --- | --- |
| `01-project-structure` | `studies/go-cli-study/reports/final/01-project-structure.md` | Mature application-style CLIs keep entrypoints thin and dependencies inward; gdu also demonstrates interchangeable presentation adapters through a shared UI boundary (`gdu/cmd/gdu/app/app.go:30-49`). | high |
| `02-command-architecture` | `studies/go-cli-study/reports/final/02-command-architecture.md` | Commands work best as explicit wiring and delegation points, as in Helm's command-to-action handoff (`helm/pkg/cmd/install.go:132-145`); large bootstrap handlers such as `opencode/cmd/root.go:49-183` obscure ownership. | high |
| `03-dependency-injection` | `studies/go-cli-study/reports/final/03-dependency-injection.md` | Constructor injection and visible composition roots make dependency ownership testable; examples include `opencode/internal/app/app.go:42-81` and gh-cli's lazy factory (`gh-cli/pkg/cmd/factory/default.go:26-46`). | high |
| `04-configuration-management` | `studies/go-cli-study/reports/final/04-configuration-management.md` | Configuration sources need explicit precedence, explicit-set tracking, and post-merge validation (`chezmoi/internal/cmd/config.go:2253-2287`; `restic/internal/global/global.go:139,147`; `opencode/internal/config/config.go:609-641`). | high |
| `05-error-handling` | `studies/go-cli-study/reports/final/05-error-handling.md` | Lower layers should preserve and classify errors while the outer boundary renders safe, contextual diagnostics (`go-task/errors/errors.go:47-50`; `restic/cmd/restic/main.go:199-209`; `gh-cli/internal/ghcmd/cmd.go:281-301`). | high |
| `06-io-abstraction` | `studies/go-cli-study/reports/final/06-io-abstraction.md` | Narrow I/O and storage boundaries permit deterministic tests and capability control, with production and memory implementations demonstrated by `restic/internal/backend/backend.go:19-90` and `restic/internal/backend/mem/mem_backend.go:52-255`. | high |
| `07-state-context` | `studies/go-cli-study/reports/final/07-state-context.md` | Root cancellation must propagate into blocking work, while dependencies remain explicit rather than hidden in context values (`helm/pkg/cmd/install.go:333-347`; `helm/pkg/action/install.go:284`; `opencode/internal/app/app.go:25-40`). | high |
| `08-concurrency` | `studies/go-cli-study/reports/final/08-concurrency.md` | Sequential work is a sound baseline; independent work can use structured, bounded concurrency and timeout-bounded cleanup (`go-task/task.go:87,197`; `gdu/pkg/analyze/parallel.go:13,36`; `opencode/cmd/root.go:261-279`). | high |
| `10-logging-observability` | `studies/go-cli-study/reports/final/10-logging-observability.md` | Structured diagnostics should use stable fields and a channel separate from normal output (`k9s/internal/slogs/keys.go:6-231`; `helm/internal/logging/logging.go:31-71`). | high |
| `11-testing-strategy` | `studies/go-cli-study/reports/final/11-testing-strategy.md` | Focused units should be complemented by real command-path tests and normalized golden fixtures (`chezmoi/internal/cmd/main_test.go:64-174`; `go-task/task_test.go:166-169`; `helm/internal/test/test.go:43`). | high |
| `12-extensibility` | `studies/go-cli-study/reports/final/12-extensibility.md` | Narrow factories and adapters preserve extension options without committing to mutable registries or public plugin contracts (`gh-cli/pkg/cmdutil/factory.go:16-43`; `dive/cmd/dive/cli/internal/command/adapter/analyzer.go:13-15`). | high |
| `13-security` | `studies/go-cli-study/reports/final/13-security.md` | Local and read-only behavior must be capability and validation boundaries, with deliberate redaction (`k9s/internal/config/json/validator.go:146`; `opencode/internal/options/secret_string.go:15-20`; `helm/pkg/registry/transport.go:37-41`). | high |
| `14-performance` | `studies/go-cli-study/reports/final/14-performance.md` | Expensive facilities should initialize lazily and data work should remain bounded or incremental (`helm/pkg/action/lazyclient.go:35-53`; `k9s/internal/pool.go:26-48`; `yq/pkg/yqlib/stream_evaluator.go:78-113`). | high |

## Relevant Patterns

- **Thin transport over typed application behavior:** gdu separates TUI, stdout, and report presentation behind one UI boundary (`gdu/cmd/gdu/app/app.go:30-49`), Helm delegates command work into an action (`helm/pkg/cmd/install.go:132-145`), and restic separates a domain repository interface from its implementation (`restic/internal/restic/repository.go:18-66`; `restic/internal/repository/repository.go:1`). For this sprint, these examples support reasoning about HTML and JSON as delivery adapters over typed read use cases rather than as owners of workflow logic. Evidence: `01-project-structure`, `02-command-architecture`, and `03-dependency-injection`.
- **Explicit composition and narrow injection:** opencode assembles application dependencies visibly (`opencode/internal/app/app.go:42-81`), while gh-cli exposes dependency-producing functions through a factory (`gh-cli/pkg/cmdutil/factory.go:16-43`). A visible serve composition path could clarify who owns configuration, app queries, rendering, logging, listener lifecycle, and cleanup without requiring handlers to discover globals. Evidence: `03-dependency-injection` and `12-extensibility`.
- **Read-only as a capability boundary:** Mature codebases distinguish dangerous capabilities and gate operations (`k9s/internal/config/plugin.go:64`; `k9s/internal/view/actions.go:142`), while narrow filesystem/backend interfaces make permitted I/O explicit (`restic/internal/fs/interface.go:10-31`; `restic/internal/backend/backend.go:19-90`). Sprint reasoning should examine whether web-facing dependencies can omit mutation and execution capabilities entirely, rather than relying only on the absence of UI controls. Evidence: `06-io-abstraction` and `13-security`.
- **Resolve and validate configuration before lifecycle start:** chezmoi preserves source precedence while applying flags (`chezmoi/internal/cmd/config.go:2253-2287`), restic tracks explicitly set flags (`restic/internal/global/global.go:139,147`), and opencode centralizes validation (`opencode/internal/config/config.go:609-641`). These patterns are relevant to listen address, workspace, browser-opening, timeout, and limit settings before a socket is opened. Evidence: `04-configuration-management`.
- **Boundary-specific error projection:** go-task uses typed errors for classification (`go-task/errors/errors.go:47-50`), restic separates fatal presentation from ordinary operational errors (`restic/cmd/restic/main.go:199-209`), and gh-cli renders errors according to command context (`gh-cli/internal/ghcmd/cmd.go:281-301`). App/query errors can therefore retain causes and identity while HTTP and CLI boundaries independently choose safe representations. Evidence: `05-error-handling`.
- **End-to-end cancellation with explicit state ownership:** Helm turns signals into cancellation and passes context into operations (`helm/pkg/cmd/install.go:333-347`; `helm/pkg/action/install.go:284`); opencode keeps dependencies in explicit application state (`opencode/internal/app/app.go:25-40`). This applies to server shutdown and request cancellation reaching artifact reads or workspace queries, without using context as a service locator. Evidence: `07-state-context`.
- **Sequential baseline, localized bounded fan-out:** Sequential execution remains effective in `urfave-cli/command_run.go:92` and `yq/pkg/yqlib/stream_evaluator.go:52`; when work is independent, go-task uses structured fan-out (`go-task/task.go:87,197`) and gdu or k9s add bounds (`gdu/pkg/analyze/parallel.go:13,36`; `k9s/internal/pool.go:21-37`). Dashboard aggregation need not become concurrent until latency or cardinality warrants it. Evidence: `08-concurrency` and `14-performance`.
- **Separate structured diagnostics from user representations:** Helm dynamically controls debug logging and sends diagnostics to stderr (`helm/internal/logging/logging.go:31-71`), while k9s centralizes field names (`k9s/internal/slogs/keys.go:6-231`). Server lifecycle and request diagnostics can remain searchable without contaminating HTML, JSON, or ordinary CLI output. Evidence: `10-logging-observability`.
- **Deterministic boundaries plus real-stack verification:** restic's functional-field mock (`restic/internal/backend/mock/backend.go:14-26`) supports focused behavior tests, while chezmoi exercises the real command path (`chezmoi/internal/cmd/main_test.go:64-174`). Normalized golden testing (`go-task/task_test.go:166-169`; `helm/internal/test/test.go:43`) may fit stable HTML, JSON, help, and error representations. Evidence: `06-io-abstraction` and `11-testing-strategy`.
- **Lazy and bounded resource use:** gh-cli and Helm delay expensive dependency construction (`gh-cli/pkg/cmdutil/factory.go:27-42`; `helm/pkg/action/lazyclient.go:35-53`), and yq processes documents incrementally (`yq/pkg/yqlib/stream_evaluator.go:78-113`). Web initialization should not burden unrelated commands, and dashboard or artifact reads should not imply whole-workspace materialization. Evidence: `03-dependency-injection` and `14-performance`.

## Trade-Offs

| Trade-Off | Benefit | Cost | When It Matters |
| --- | --- | --- | --- |
| Narrow web query interfaces vs broad app/repository injection | Encodes read-only capability, limits coupling, and simplifies fakes (`restic/internal/restic/repository.go:18-66`) | Adds mapping types and can become ceremony if split too finely | When deciding what `internal/web` may receive from `internal/app` |
| Central composition root vs distributed construction | Makes ownership and shutdown traceable (`opencode/internal/app/app.go:42-81`) | Can become a god object like the broad shape at `chezmoi/internal/cmd/config.go:193-291` | When wiring CLI, TUI, and web without globals or cycles |
| Eager construction vs lazy factories | Eager construction fails early; lazy factories avoid costs for unrelated commands (`gh-cli/pkg/cmd/factory/default.go:26-46`) | Lazy failures occur later and lifecycle can be less obvious | When serve-only resources would otherwise affect help, version, CLI, or TUI startup |
| Immutable startup configuration vs hot reload | Immutable state is predictable and concurrency-friendly | Restart is needed for changes; hot reload adds watcher and precedence complexity (`k9s/internal/ui/config.go:134-190`) | When deciding whether a foundation sprint needs runtime configuration changes |
| Live per-request reads vs cached snapshots | Live reads maximize freshness and minimize cache ownership | Repeated work can amplify refresh load; snapshots add invalidation, staleness, and memory costs | When defining dashboard freshness and authoritative-state projection |
| Sequential aggregation vs structured fan-out | Sequential work is deterministic and easy to cancel | Independent slow sources may increase latency; fan-out adds cancellation, ordering, and partial-result semantics (`go-task/task.go:87,197`) | When one view combines several independent queries |
| Fail-fast fan-out vs partial rendering | Fail-fast simplifies consistency and sibling cancellation | One failed panel can suppress useful results; partial rendering needs per-result errors | When a dashboard section fails but other read-only data remains available |
| Golden representation tests vs semantic assertions | Goldens reveal complete HTML/JSON changes (`helm/internal/test/test.go:43`) | Volatile values cause brittleness and updates require review | When deciding which parts of server-rendered output are compatibility-sensitive |
| Internal evolving API vs documented versioned API | Internal representations can change quickly | A documented `/api/v1` creates schema and error-envelope compatibility obligations; explicit version conversion has maintenance cost (`helm/internal/plugin/metadata_v1.go:24-48`; `helm/internal/plugin/metadata.go:114-130`) | Before clients beyond bundled browser code depend on JSON |
| Rich debug diagnostics vs strict disclosure minimization | Detailed logs accelerate local diagnosis | Paths, artifact content, environment data, or secrets may leak unless fields are classified and redacted (`opencode/internal/options/secret_string.go:15-20`) | When choosing request fields, health detail, and debug behavior |

## Anti-Patterns And Warnings

- **Transport-owned business logic:** Avoid reproducing app validation, workspace interpretation, or workflow logic in handlers. Large orchestration functions such as `opencode/cmd/root.go:49-183` show how wiring and behavior become difficult to separate. Evidence: `02-command-architecture`.
- **Global services, configuration, or registries:** Package globals hide ownership and make concurrent tests or multiple server instances interfere; examples include `dive/internal/bus/bus.go:5`, `dive/internal/log/log.go:9`, and `rclone/fs/cache/cache.go:16-21`. Evidence: `01-project-structure` and `03-dependency-injection`.
- **Context as a dependency container or detached context:** Context-carried services create runtime-only dependencies (`k9s/internal/keys.go:10-38`), while `context.Background()` inside work disconnects cancellation (`chezmoi/internal/cmd/templatefuncs.go:215`). Accepting but ignoring context, as at `dive/image/analysis.go:20`, is equally misleading. Evidence: `07-state-context`.
- **Error strings as contracts or raw internal errors as responses:** String comparison loses wrapping and identity (`age/cmd/age/age.go:77,353,381`), while verbatim internal errors can disclose implementation details. Preserve causes for diagnostics and classify only distinctions that alter boundary behavior. Evidence: `05-error-handling`.
- **Mixed injected and direct I/O:** A single direct `os.*` or default-client path can bypass a testable boundary (`chezmoi/internal/cmd/templatefuncs.go:296`; `gh-cli/pkg/updates/updates.go:268`). Handlers should not gain arbitrary filesystem or network access through convenience calls. Evidence: `06-io-abstraction`.
- **UI-only read-only enforcement:** The lack of buttons does not prevent accidental mutation when handlers hold command executors or mutable repositories. The explicit dangerous-operation distinctions at `k9s/internal/config/plugin.go:64` and `k9s/internal/view/actions.go:142` are a warning to constrain capabilities structurally. Evidence: `13-security`.
- **Treating local input as trusted:** URL paths, query values, headers, artifact references, and bodies still cross a parser boundary. Central validation (`k9s/internal/config/json/validator.go:146`), bounded inputs, containment, and secret-safe rendering remain necessary. Evidence: `13-security`.
- **Unbounded or ownerless goroutines:** Goroutine-per-item fan-out (`gh-cli/pkg/cmd/extension/manager.go:196-206`) and fire-and-forget work (`dive/cmd/dive/cli/internal/command/adapter/resolver.go:70`) complicate request cancellation and shutdown. Every goroutine needs an owner, bound, cancellation path, and completion policy. Evidence: `08-concurrency` and `14-performance`.
- **Whole-workspace loading for a first response:** Materializing all plans, artifacts, logs, or history increases latency and memory; incremental processing at `yq/pkg/yqlib/stream_evaluator.go:78-113` and bounded pools at `k9s/internal/pool.go:26-48` provide counter-patterns. Evidence: `14-performance`.
- **Handler-only tests:** Unit tests cannot reveal command wiring, bind failures, listener cleanup, shutdown, or output-channel regressions. Real command-path tests such as `chezmoi/internal/cmd/main_test.go:64-174` address that gap. Evidence: `11-testing-strategy`.
- **Accidental public extension contract:** A mutable HTTP registry can silently overwrite duplicate names (`rclone/fs/rc/registry.go:41-48`), and an externally consumed API acquires versioning costs. Do not infer that a read-only foundation needs plugin registration or an ungoverned public API. Evidence: `12-extensibility`.
- **Diagnostics mixed with responses or left unredacted:** Helm's stderr separation (`helm/internal/logging/logging.go:71`) and credential scrubbing (`helm/pkg/registry/transport.go:37-41`) show why logs, HTML, JSON, and CLI stdout need distinct policies. Evidence: `10-logging-observability` and `13-security`.

## Examples Worth Inspecting

| Example | Path / Source | Why It Is Useful |
| --- | --- | --- |
| Multiple presentation adapters | `gdu/cmd/gdu/app/app.go:30-49` via `studies/go-cli-study/reports/final/01-project-structure.md` | Shows presentation variation behind a shared boundary rather than duplicated application behavior. |
| Thin command delegation | `helm/pkg/cmd/install.go:132-145` via `studies/go-cli-study/reports/final/02-command-architecture.md` | Shows a command configuring and delegating to reusable behavior. |
| Explicit application composition | `opencode/internal/app/app.go:42-81` via `studies/go-cli-study/reports/final/03-dependency-injection.md` | Shows centralized dependency assembly and staged initialization pressures. |
| Flag/config precedence | `chezmoi/internal/cmd/config.go:2253-2287` and `restic/internal/global/global.go:139,147` via `studies/go-cli-study/reports/final/04-configuration-management.md` | Shows how explicit CLI values can override config without defaults doing so accidentally. |
| Error classification and presentation | `go-task/errors/errors.go:47-50`, `restic/cmd/restic/main.go:199-209`, and `gh-cli/internal/ghcmd/cmd.go:281-301` via `studies/go-cli-study/reports/final/05-error-handling.md` | Separates domain identity, process behavior, and user-facing rendering. |
| Production and memory I/O boundaries | `restic/internal/backend/backend.go:19-90` and `restic/internal/backend/mem/mem_backend.go:52-255` via `studies/go-cli-study/reports/final/06-io-abstraction.md` | Demonstrates deterministic substitution at a meaningful boundary. |
| Signal and operation cancellation | `helm/pkg/cmd/install.go:333-347` and `helm/pkg/action/install.go:284` via `studies/go-cli-study/reports/final/07-state-context.md` | Shows lifecycle cancellation reaching blocking work. |
| Bounded work and shutdown | `k9s/internal/pool.go:21-37` and `opencode/cmd/root.go:261-279` via `studies/go-cli-study/reports/final/08-concurrency.md` | Shows concurrency limits and a bounded completion wait. |
| Structured logging vocabulary | `k9s/internal/slogs/keys.go:6-231` and `helm/internal/logging/logging.go:31-71` via `studies/go-cli-study/reports/final/10-logging-observability.md` | Shows consistent fields, runtime levels, and diagnostic-channel separation. |
| Real command-path testing | `chezmoi/internal/cmd/main_test.go:64-174` via `studies/go-cli-study/reports/final/11-testing-strategy.md` | Provides a model for testing composition and startup rather than handlers alone. |
| Existing CLI-hosted HTTP server | `fzf/src/server.go:40-145` via `studies/go-cli-study/reports/final/12-extensibility.md` | Offers a concrete server lifecycle to inspect critically; it is evidence, not a template to copy. |
| Validation and redaction | `k9s/internal/config/json/validator.go:146`, `opencode/internal/options/secret_string.go:15-20`, and `helm/pkg/registry/transport.go:37-41` via `studies/go-cli-study/reports/final/13-security.md` | Shows centralized validation and secret-safe output boundaries. |
| Lazy and incremental resource use | `helm/pkg/action/lazyclient.go:35-53` and `yq/pkg/yqlib/stream_evaluator.go:78-113` via `studies/go-cli-study/reports/final/14-performance.md` | Helps evaluate serve-only initialization and bounded artifact/data processing. |

## Design Pressures

- The browser is an additional presentation surface, so shared typed app queries must remain independent of CLI, TUI, HTTP, templates, and JSON while still producing enough information for useful views.
- Read-only behavior must survive future growth; exposing only query capabilities now reduces the chance that later handlers accidentally gain mutation, runtime, process, or Git powers.
- A server is longer-lived than an ordinary command. Listener ownership, root cancellation, request cancellation, handler completion, and timeout-bounded shutdown must be coherent without creating browser-owned durable state.
- Loopback operation reduces exposure but does not remove browser-origin, hostile-input, path-containment, redaction, or local multi-user trust concerns. The selected security evidence establishes the pressure for explicit boundaries but does not settle HTTP-specific controls.
- HTML and JSON have different representation needs but should not drive duplicate app logic. Safe error projection must also preserve enough internal detail for diagnostics without leaking it to either representation.
- Artifact previews combine potentially hostile content, path selection, type checks, size limits, parsing, and escaping. Broad filesystem injection would work against the required capability and containment boundaries.
- Dashboard freshness competes with bounded work. Per-request reads are simple and current but can repeat scans; snapshots reduce request cost but introduce stale-state and synchronization questions.
- Concurrency can improve multi-panel latency but also determines whether errors fail the whole view or permit partial results. The evidence favors adding only measured, localized, bounded concurrency.
- Existing CLI and TUI behavior must not pay web startup, dependency, logging, or configuration costs. Serve-only resources and explicit surface composition are therefore under pressure to remain isolated.
- Once `/api/v1` is documented, browser and other consumers may treat its envelopes as stable. The sprint needs enough representation discipline for compatibility without prematurely creating a general extension platform.
- Tests must cover narrow behavior and actual composition. Deterministic fakes, normalized representations, real listener/server paths, cancellation, and cleanup each reveal different failure classes.
- Embedded templates and assets avoid a second build/runtime toolchain, but template parsing, rendering, escaping, and static delivery become binary-startup and test concerns.

## Open Questions For Reasoning

- What is the smallest typed read-use-case surface that supports dashboard, detail, validation, health, and artifact views without exposing product mutations?
- Which component owns construction and closure of workspace readers, templates, listener, HTTP server, browser launcher, and root cancellation?
- Should handlers receive one cohesive read facade or several narrower query interfaces, and where would either choice become excessive coupling or ceremony?
- Which serve settings participate in existing configuration precedence, and how are untouched flag defaults distinguished from explicit overrides?
- Is configuration immutable after startup, or is any reload behavior justified in this sprint?
- What exact domain error distinctions lead to different HTTP statuses, API error codes, HTML states, CLI diagnostics, or exit codes?
- What public error details are safe, and which paths, artifact content, environment values, and internal causes require redaction?
- Does each request read authoritative workspace state directly, consume an immutable snapshot, or combine both approaches? What are the freshness and stale-data semantics?
- If one dashboard section fails, should successful sections still render? How would that choice affect error aggregation and any use of structured concurrency?
- Which reads are sufficiently independent and costly to justify concurrency, what bounds apply, and who owns cancellation when a client disconnects or the server shuts down?
- What are the concrete bounds for request bodies, query values, response rows, artifact bytes, concurrent requests, read duration, and shutdown duration?
- How are artifact references represented so that only allowlisted Markdown/JSON within the workspace can be selected without accepting arbitrary filesystem paths?
- Which Host and Origin values are valid for each supported loopback listen form, and how should requests with absent or malformed Origin be treated?
- What structured fields are useful for startup, requests, query failures, rejection, and shutdown, and which fields are safe at debug versus normal levels?
- Is `/api/v1` intended only for the bundled browser or as a documented external local API, and which success/error fields become compatibility-controlled now?
- Which HTML and JSON outputs merit whole-file goldens, which need semantic assertions, and how will ports, paths, times, and IDs be normalized?
- How will tests exercise actual serve command composition, bind rejection, graceful shutdown, listener cleanup, and CLI/TUI non-regression in addition to handler behavior?
- Does optional browser opening need an injected launcher, and is launch failure fatal, warning-only, or merely diagnostic?
- Are health results limited to process readiness, or do they project workspace/query availability; how can either remain truthful without expensive whole-workspace work?
- Which accessibility and progressive-enhancement behavior can be guaranteed without expanding the embedded template/CSS/minimal-JavaScript scope?

## Evidence Pointers

- `studies/go-cli-study/reports/final/01-project-structure.md`: inspect thin entrypoints, `internal` boundaries, dependency direction, and gdu's presentation abstraction.
- `studies/go-cli-study/reports/final/02-command-architecture.md`: inspect command-to-action delegation and the monolithic bootstrap cautions.
- `studies/go-cli-study/reports/final/03-dependency-injection.md`: inspect composition roots, constructor injection, lazy factories, and god-object pressure.
- `studies/go-cli-study/reports/final/04-configuration-management.md`: inspect source precedence, explicit-set tracking, post-merge validation, and reload trade-offs.
- `studies/go-cli-study/reports/final/05-error-handling.md`: inspect wrapping, sentinels and typed errors, exit classification, and context-specific rendering.
- `studies/go-cli-study/reports/final/06-io-abstraction.md`: inspect narrow stream/filesystem/backend seams and in-memory test implementations.
- `studies/go-cli-study/reports/final/07-state-context.md`: inspect root signal contexts, propagation breadth, application/session ownership, and detached-context warnings.
- `studies/go-cli-study/reports/final/08-concurrency.md`: inspect sequential baselines, structured fan-out, bounds, goroutine ownership, and shutdown waits.
- `studies/go-cli-study/reports/final/10-logging-observability.md`: inspect structured keys, dynamic levels, stderr separation, and expensive-debug guards.
- `studies/go-cli-study/reports/final/11-testing-strategy.md`: inspect full command-path tests, functional mocks, normalized goldens, and integration setup.
- `studies/go-cli-study/reports/final/12-extensibility.md`: inspect adapter/factory boundaries, the fzf HTTP server, registry collision risk, and version conversion.
- `studies/go-cli-study/reports/final/13-security.md`: inspect capability distinctions, centralized validation, secret wrappers, credential scrubbing, and safe process invocation.
- `studies/go-cli-study/reports/final/14-performance.md`: inspect lazy initialization, incremental processing, bounded pools, profiling trade-offs, and long-session resource pressure.

## Handoff To Reasoning

- Use this handbook as evidence input.
- Validate whether the observed patterns fit this project's constraints.
- Resolve the open questions in architecture, API-design, frontend, and sprint-level reasoning before implementation choices are made.
- Treat project requirements and docs as scope constraints, not as substitutes for the selected study evidence.
- Do not copy external patterns without sprint-specific reasoning.
