> **Inputs Used:** `projects/ultraplan-go/sprints/30-web-foundations/sprint-index.md`, `projects/ultraplan-go/sprints/30-web-foundations/technical-handbook.md`, `projects/ultraplan-go/sprints/30-web-foundations/requirements.md`, `projects/ultraplan-go/docs/ARCHITECTURE.md`, `projects/ultraplan-go/docs/PRD.md`, `projects/ultraplan-go/docs/TRD.md`, `system/reasoning/architecture_reasoning_template.md`

# Architecture: Local Web Foundation

This area covers the package boundaries, composition, state ownership, dependency direction, and lifecycle needed to add `ultraplan serve` as a read-only local interface. The browser is a new adapter over the existing product core, not a new workflow, persistence model, or alternate application layer.

## Area Decisions

### Architectural fit and chosen shape

Proceed after a small composition refactor. The module-driven architecture remains correct, but adding a side-effectful web runner directly inside `internal/app` would create the prohibited cycle `internal/app -> internal/web -> internal/app`. The refactor must make all local-interface runners explicit at `cmd/ultraplan` while preserving `internal/app` as the shared use-case and command boundary.

The required dependency graph is:

```text
cmd/ultraplan
  -> internal/app
  -> internal/tui
  -> internal/web

internal/tui -> internal/app
internal/web -> internal/app
internal/app -> product modules + platform modules
```

No package-global runner, service locator, mutable registry, callback registration in `init`, or context-carried dependency is permitted. `cmd/ultraplan` explicitly constructs the app services and typed read surfaces, constructs TUI and web runners from those surfaces, and supplies the runners to app command dispatch. Serve-only dependencies are constructed lazily only when the selected command is `serve`, so help, version, existing CLI commands, and TUI startup do not parse templates, initialize HTTP state, or acquire a listener.

The primary units are:

| Unit | Ownership |
| --- | --- |
| `cmd/ultraplan/main.go` | Process-level composition, signal-root context, construction order, and final exit. |
| `internal/app/surfaces.go` | Narrow runner function/interface contracts and explicit CLI/TUI/web surface composition types; no adapter implementation. |
| `internal/app/serve_commands.go` | Serve argument parsing, existing workspace/config preflight, explicit-set precedence, loopback option validation, invocation of the injected web runner, and CLI exit/error projection. |
| `internal/app/web_usecases.go` | Cohesive typed read facade for dashboard, project, sprint, study, validation, health, and artifact-preview queries; aggregation and product error identity. |
| `internal/web/server.go` | `net/http` listener/server ownership, browser-launch hook, bounded timeouts, request serving, graceful shutdown, and cleanup. |
| `internal/web/routes.go`, `handlers.go`, `security.go` | Route/method registration, HTTP DTO/template mapping, transport validation, safe error projection, and browser security policy. |
| `internal/web/artifacts.go` | Preview transport/rendering policy over the typed app artifact result; it does not gain arbitrary filesystem access. |
| `internal/web/templates` and `internal/web/static` | Embedded server-rendered presentation assets only. |

Use one cohesive `internal/app` web read facade rather than injecting product services or creating one interface per route. Its methods are use-case-oriented and return plain app result types. This is narrow in capability even though it covers several related reads: it contains no execute, validate-now, review, smoke, runtime, process, cancellation, persistence, or Git method. If later operations enter scope, add a separately reviewed command surface rather than widening the Sprint 30 query facade implicitly.

`internal/web` may depend on the exported app facade contract and result types only. It must not import `internal/project`, `internal/sprint`, `internal/study`, `internal/workspace`, runtime/process adapters, or CLI command handlers. App use cases delegate product interpretation to the owning modules; handlers perform only request validation, use-case invocation, DTO mapping, and rendering.

### Request and lifecycle flow

Startup and shutdown follow this explicit flow:

```text
process signal/root context
  -> cmd composition
  -> app serve command parsing
  -> workspace discovery and merged config validation
  -> loopback listen option validation
  -> injected web runner
  -> listener acquisition and HTTP serving
  -> root cancellation or serve failure
  -> timeout-bounded HTTP shutdown
  -> listener/handler cleanup
  -> app error classification and process exit mapping
```

The listen address must be resolved and rejected if non-loopback before a socket accepts requests. Immutable effective serve configuration is passed to the runner; Sprint 30 has no hot reload. Listener creation remains in `internal/web`, but it occurs only after app preflight succeeds.

The runner blocks until serving fails or its context is cancelled. Request contexts derive from the server context and propagate through app queries and bounded artifact reads. Graceful shutdown stops accepting requests, gives in-flight read handlers a bounded completion period, then closes remaining connections. Every goroutine is owned by the server lifecycle and joined or cancelled before `Run` returns.

The optional browser opener is an injected edge capability, not a direct package-global OS call. It runs only after successful listener acquisition. Launch failure is warning-only because the server remains usable; the warning goes to the diagnostic channel and does not contaminate normal machine-readable output.

### State ownership

- Workspace artifacts and product-owned run/flow/review/smoke state remain the sole durable source of truth. Sprint 30 creates, updates, and deletes no durable product state through HTTP.
- App queries compute fresh typed projections from authoritative state for each request. The web package does not retain a product snapshot or duplicate validation/workflow logic.
- Permitted ephemeral state is limited to immutable effective server configuration, parsed embedded templates/assets, listener/server objects, request contexts, generated request IDs, response view models, and bounded in-flight preview buffers.
- There is no operation hub, confirmation store, SSE subscriber registry, browser session database, cookie state, cache, watcher, or background refresh loop in this sprint.
- Health is derived on demand from server readiness and a lightweight app workspace query. It is not persisted and does not perform runtime/provider or whole-workspace health checks.

Artifact resolution preserves the package boundary: app queries issue opaque safe references and own resolution to an allowlisted, contained Markdown/JSON artifact result. `internal/web/artifacts.go` enforces transport size/rendering behavior over that result. Handlers never accept or open arbitrary absolute/workspace-relative paths and never receive a broad filesystem dependency.

### Error, diagnostics, and compatibility

Lower layers retain wrapped typed error identity. The serve command maps startup/config/workspace/cancellation failures to existing CLI exit classes, while HTTP handlers independently project the safe HTML or versioned JSON representation. Neither boundary compares error strings or exposes raw causes.

Structured server diagnostics use a channel separate from HTML, JSON, and ordinary CLI stdout. Stable fields include request ID, normalized route, method, status, duration, response size, safe error class, server lifecycle event, and shutdown outcome. Workspace paths, artifact contents, environment values, Host/Origin input, secrets, and raw internal errors are omitted or redacted.

The `/api/v1` DTOs and error envelope are transport-owned compatibility contracts; app result structs are not serialized directly. Existing CLI and TUI result/output contracts remain unchanged. The only shared refactor is explicit runner composition, with command-path tests proving existing command dispatch and TUI invocation still use their prior behavior.

### Verification decision

The architecture is accepted only with all of the following evidence:

- Import inspection proves `internal/web -> internal/app` and no direct web import of product, runtime, process, or CLI packages.
- Composition tests prove `cmd/ultraplan` supplies independent CLI/TUI/web dependencies without globals and constructs serve-only resources only for `serve`.
- App tests use deterministic product fakes to cover typed dashboard/detail/validation/health/artifact queries and typed error preservation.
- Web tests use fake app queries and `httptest` for routes, DTOs, templates, security, request cancellation, listener failures, graceful shutdown, and artifact bounds/containment/hostile content.
- Command-path tests cover help, workspace/config precedence, non-loopback rejection before runner invocation, startup failure, browser-launch warning, cancellation, exit mapping, and unchanged existing CLI/TUI dispatch.
- `go test ./...`, `go test -race ./...`, and `go build ./cmd/ultraplan` pass without requiring a real runtime, smoke harness, browser, Node.js, asset server, or network service.

## Trade-Offs

| Decision | Benefit | Cost / Rejected Alternative |
| --- | --- | --- |
| Small explicit composition refactor before adding web | Prevents package cycles and makes runner ownership and test substitution visible | Directly constructing `internal/web` from `internal/app` was rejected because web must import app; package-global registration was rejected because it hides mutable ownership. |
| `cmd/ultraplan` as process composition root | One traceable place owns signal context and side-effectful interface construction | Distributed self-construction was rejected because lifecycle and cleanup become hidden; the root must stay thin to avoid becoming a behavior-heavy god object. |
| Lazy serve-only construction | Existing commands avoid template parsing, HTTP initialization, and web-specific failure/cost | Eagerly constructing every surface can fail earlier, but was rejected because unrelated CLI/TUI paths must not pay web startup or dependency costs. |
| Cohesive read facade | Encodes read-only capability and keeps handlers simple without excessive interface ceremony | Injecting all app/product services was rejected as broad stamp coupling; one interface per route was rejected as premature fragmentation. |
| Transport DTOs separate from app results | Allows HTML/JSON evolution without exposing product internals or binding app refactors to `/api/v1` | Direct serialization is shorter but creates disclosure and compatibility coupling. |
| Live reads without a web cache | Durable workspace state remains authoritative and restart behavior is trivial | Repeated refresh can amplify reads; cached snapshots were rejected because invalidation, staleness, locking, and memory ownership are not justified yet. |
| Sequential app aggregation initially | Deterministic order, simple cancellation, and clear whole-query failure behavior | Parallel dashboard fan-out may reduce latency but was rejected until measurement justifies partial-result and goroutine semantics. |
| Web lifecycle object with explicit dependencies | Cohesively owns listener, server, templates, security middleware, launcher, and cleanup | Package functions plus globals are shorter but make multiple instances and parallel tests interfere. The object owns lifecycle state, not product state. |
| App-owned artifact reference resolution with web-owned rendering | Keeps containment and artifact identity near workspace/product knowledge while preserving web-specific bounds and escaping | Direct filesystem reads in handlers were rejected because they bypass capability control and violate `internal/web -> internal/app`. |
| Immutable startup configuration | Predictable concurrency and validation before listening | Hot reload was rejected because watchers, source precedence changes, and lifecycle races add no Sprint 30 value. |
| Warning-only browser-launch failure | A launcher problem does not invalidate an already healthy server | Treating launch failure as fatal was rejected because browser opening is optional convenience; silently ignoring it was rejected because diagnostics must remain actionable. |

## Evidence

- The handbook's project-structure and command-architecture evidence (`01-project-structure`, `02-command-architecture`) favors thin entrypoints, inward dependencies, interchangeable presentation adapters, and command-to-action delegation. This supports a thin `cmd` composition root and web as another adapter over app behavior.
- Constructor injection and visible application assembly in `03-dependency-injection`, together with the narrow factory/adapter evidence in `12-extensibility`, support explicit runner construction and reject mutable registries or service discovery.
- The handbook's read-only capability pattern (`06-io-abstraction`, `13-security`) supports a query-only app facade and app-issued artifact references. Read-only is enforced by absent mutation/runtime capabilities, not merely by rendering no action controls.
- Configuration evidence (`04-configuration-management`) shows that defaults, workspace config, environment, and explicitly set flags need source-aware merge and post-merge validation. This grounds app-owned serve preflight and rejection before listener acquisition.
- Typed error classification and boundary rendering in `05-error-handling` support preserving causes in app/product layers while CLI and HTTP choose independent safe output and exit/status mappings.
- Signal/context propagation evidence in `07-state-context` supports one process root context, request-derived contexts, explicit dependencies rather than context values, and cancellation reaching blocking reads and server shutdown.
- The concurrency evidence (`08-concurrency`) warns against ownerless goroutines and favors timeout-bounded cleanup. The architecture therefore gives every server/background action one lifecycle owner and starts with sequential query aggregation.
- Structured logging evidence (`10-logging-observability`) separates diagnostics from normal output and standardizes safe fields. This grounds the dedicated server diagnostic channel and redacted request/lifecycle vocabulary.
- Deterministic seam plus real-path testing in `11-testing-strategy` supports fake app use cases, `httptest`, command composition tests, listener/shutdown tests, and full build/race verification rather than handler tests alone.
- Lazy and bounded resource evidence (`14-performance`) supports constructing web facilities only for `serve`, bounded previews, no whole-workspace loading, and measurement before caching or concurrent fan-out.
- Project architecture and sprint requirements explicitly assign shared use cases to `internal/app`, HTTP concerns to `internal/web`, side-effectful construction to `cmd/ultraplan`, and durable state to product modules/workspace artifacts. The selected evidence fits those project constraints without requiring a new global layer or plugin architecture.

The evidence is high-confidence for dependency direction, explicit composition, capability restriction, cancellation, diagnostics, bounded work, and testing. Exact Go type names and file-level decomposition remain implementation details so long as the decided ownership and import rules are preserved.

## Risks

- **Composition-root growth:** Adding a third surface can turn `main.go` into a god function. Mitigation: keep it to construction and handoff, with app-owned use cases and adapter-owned behavior; use small constructors rather than a generic dependency container.
- **Cycle pressure:** App command dispatch needs a web runner while web needs app queries. Mitigation: define the runner contract in app or as a plain function consumed by app, implement it in web, and wire both only in `cmd`; enforce with import review and `go test`.
- **Facade creep:** Future operation work may append mutation methods to the read facade. Mitigation: keep Sprint 30's facade query-only and introduce a separately named command/operation capability only when a later sprint governs confirmation, locking, and recovery.
- **Duplicate mapping logic:** CLI, TUI, HTML, and JSON can repeat product interpretation. Mitigation: app results contain shared semantics; adapters duplicate only representation-specific mapping. Do not extract a generic renderer that couples unrelated surfaces.
- **Hidden filesystem access:** Convenience reads in web handlers could bypass app containment and test fakes. Mitigation: prohibit direct web imports/readers for workspace data and test artifact behavior through opaque app references.
- **Shutdown leaks:** Browser launch hooks, handlers, or listener goroutines can outlive cancellation. Mitigation: inject blocking/async edges, give each an owner, propagate context, bound shutdown, and run race/leak-oriented lifecycle tests.
- **Freshness and latency:** Live reads are truthful but may become expensive for large workspaces. Mitigation: bounded results, no recursive repository scans, request cancellation, and measurement before introducing cache or concurrency semantics.
- **Error disclosure:** A shared internal error may be safe for CLI diagnostics but unsafe for HTTP. Mitigation: retain typed identity and apply transport-specific projection/redaction; never serialize internal errors or compare their text.
- **Broader Phase 4 scope leakage:** Project docs describe guarded operations and SSE, but Sprint 30 explicitly excludes them. Mitigation: do not construct operation/runtime/process capabilities in web and test that no mutation, confirmation, operation, SSE, or cancellation routes are registered.
- **Existing-surface regression:** Refactoring runner composition can alter CLI/TUI startup, help, output, or exit behavior. Mitigation: preserve existing adapters unchanged behind explicit dependencies and add command-path non-regression tests before accepting the architecture.
- **Final reasoning handoff:** `projects/ultraplan-go/sprints/30-web-foundations/reasoning.md` must reference this document and carry forward its composition, dependency, state-ownership, lifecycle, and capability decisions. That artifact is outside this repair and is not modified here.

No architecture question blocks implementation. Public/remote serving, authentication, operation state, confirmations, SSE, mutation locks, web caching, and concurrent partial dashboard aggregation remain deferred until a later sprint explicitly selects and governs them.
