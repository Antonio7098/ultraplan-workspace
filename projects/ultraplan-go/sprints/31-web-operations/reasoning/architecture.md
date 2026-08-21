# Architecture: Guarded Web Operations and SSE Progress

> **Inputs Used:** `projects/ultraplan-go/sprints/31-web-operations/technical-handbook.md`, `projects/ultraplan-go/sprints/31-web-operations/requirements.md`, `projects/ultraplan-go/sprints/31-web-operations/sprint-index.md`, `projects/ultraplan-go/docs/ARCHITECTURE.md`, `projects/ultraplan-go/docs/PRD.md`, `projects/ultraplan-go/docs/TRD.md`, `system/reasoning/architecture_reasoning_template.md`, `studies/go-cli-study/reports/final/01-project-structure.md`, `studies/go-cli-study/reports/final/03-dependency-injection.md`, `studies/go-cli-study/reports/final/04-configuration-management.md`, `studies/go-cli-study/reports/final/05-error-handling.md`, `studies/go-cli-study/reports/final/06-io-abstraction.md`, `studies/go-cli-study/reports/final/07-state-context.md`, `studies/go-cli-study/reports/final/08-concurrency.md`, `studies/go-cli-study/reports/final/10-logging-observability.md`, `studies/go-cli-study/reports/final/11-testing-strategy.md`, `studies/go-cli-study/reports/final/12-extensibility.md`, `studies/go-cli-study/reports/final/13-security.md`, `studies/go-cli-study/reports/final/14-performance.md`

This area covers the package, state, concurrency, and lifecycle architecture for exposing existing UltraPlan operations through the local web server. The key conclusion is to add a small typed operation capability in `internal/app`, product-owned mutation exclusion in the existing product modules, and a bounded transport-lifecycle hub in `internal/web`. The browser remains a third adapter over the same product workflows; it does not gain a workflow engine, scheduler, durable job store, runtime adapter, or product-state authority.

## Area Decisions

### Existing architecture fit

The current module-driven architecture fits with a small refactor before feature work. Existing CLI/TUI command glue that still owns operation assembly must move behind shared `internal/app` use cases before web handlers call it. Product workflow rules remain where their state lives:

```text
cmd/ultraplan
  -> constructs app operation dependencies, web server, and operation hub

internal/web
  -> HTTP/session/confirmation mapping
  -> bounded operation handles and SSE subscriptions
  -> narrow internal/app operation capability only

internal/app
  -> normalization, preparation, dispatch, safe progress/result projection
  -> existing project/sprint/study/runtime/harness use cases

internal/sprint and internal/study
  -> workflow state, durable state, mutation locks, recovery, verdicts

internal/platform/runtime and internal/platform/process
  -> generic cancellable external execution and process-tree cleanup
```

This preserves the documented `internal/web -> internal/app -> product/platform` dependency direction. `internal/web` must not import `internal/sprint`, `internal/study`, `internal/project`, runtime/process adapters, or CLI handlers. Product/platform packages must not import `internal/web` or its transport DTOs.

### Typed app operation capability

Add one cohesive app operation service, separate from the Sprint 30 read-only query facade. Its public concepts are transport-neutral:

- `OperationKind`: the closed set of existing validation, prompt-preview, dry-run, sprint-flow, execute, review, smoke, verify, and study-run-loop capabilities.
- `OperationRequest`: a discriminated request containing typed scope and operation-specific options, with no HTTP, cookie, token, SSE, or template fields.
- `OperationPreparation`: normalized request, affected workspace-relative paths, mutation class, current governed-input fingerprint, effective runtime/model or harness summary, prerequisites, and expiry input for a transport confirmation.
- `OperationEvent`: a safe canonical progress projection suitable for CLI, TUI, and web adapters; it is not a raw provider/runtime event.
- `OperationResult`: a safe typed terminal summary with durable refresh information.
- Typed errors for invalid scope, stale input, mutation conflict, unavailable prerequisite, validation failure, runtime failure, timeout, cancellation, interruption, and cleanup uncertainty.

The service exposes preparation and execution as separate operations. Preparation resolves and normalizes the complete effective request without mutation. Execution accepts the normalized request plus the fingerprint expected by the caller, revalidates current inputs, acquires product locks where required, and runs synchronously with respect to its supplied context while publishing events through a narrow sink. `internal/web` makes execution asynchronous by running this blocking app call under a server-owned context; app and product packages do not know about web operation IDs or subscribers.

Status and result inspection split by authority:

- The web hub reports ephemeral server-lifecycle state for a retained operation handle.
- App query use cases report authoritative durable project/sprint/study state before, during, and after hub retention.
- The app operation result links the two using safe product/run identifiers and refresh guidance, but neither layer copies durable workflow state into the hub.

The capability is intentionally closed and explicit for Sprint 31. A generic command registry, reflection-based dispatcher, plugin API, browser-authored command, or universal `map[string]any` operation is rejected. A later existing use case can extend the discriminated union and app dispatcher without introducing dynamic loading.

### Composition and dependency injection

`cmd/ultraplan` remains the composition root. It constructs concrete product services, runtime/process adapters, clock/ID sources, the app operation service, the web operation hub, and the HTTP server explicitly. Production and test construction differ only in injected implementations.

Use constructor injection for:

- app operation dispatch dependencies;
- clock and operation/preparation ID generation;
- product lock stores/services;
- progress sink and safe diagnostic projection;
- runtime, smoke harness, and process boundaries already used by app/product services;
- web hub limits and retention policy.

Use `context.Context` only for lifecycle, deadlines, cancellation, and request-scoped correlation values. Do not store app services, locks, the hub, or configuration in context values. Do not add package-level mutable registration, a service locator, or a DI framework. Interfaces stay consumer-sized and appear only where multiple current implementations or deterministic test fakes justify them.

### Product-owned mutation exclusion

Mutation exclusion is enforced below HTTP so CLI, TUI, and web cannot bypass it.

- `internal/sprint` owns an exclusive per-sprint mutation lease for sprint flow, execute, review, smoke, and verify. These operations share one exclusion domain because each can read and change overlapping governed artifacts, flow state, run state, review/smoke freshness, or external evidence links.
- Existing study run-loop locking remains owned by `internal/study`. The app operation service invokes the same study capability and lock path used by CLI/TUI.
- Read-only validation, prompt preview, dry-run preparation, status, and result inspection acquire no mutation lease unless an existing product use case already performs mutation. Preparation never reserves a lease.
- A lease contains safe owner metadata sufficient for typed conflict diagnostics and internal recovery, but HTTP receives only an allowlisted scope, operation kind, start time, and recovery action.
- Lock acquisition occurs after request/fingerprint revalidation and before the first mutation or external runtime/harness call. Release occurs only after durable terminal reconciliation or explicit `cleanup_uncertain` recording.
- Cancellation request, SSE disconnect, HTTP response completion, or shutdown deadline expiry does not itself release the lock. The product operation owns release in its terminal cleanup path.
- Stale lock handling is conservative. Startup reconciles durable running state, owned-process evidence, and lock metadata; process absence or partial artifacts never imply success. Ambiguity becomes `interrupted` or `recovery_required`, not an automatic unlock-and-success path.

No cross-product global mutex is added. Per-sprint and per-study scopes preserve unrelated concurrency, while the web hub separately limits total active server-owned operations for resource protection.

### Server-owned operation lifecycle

Every accepted web operation has one lifecycle owner: the server operation hub. The hub owns only:

- an opaque operation ID and initiating session identity;
- immutable normalized request summary and safe correlation fields;
- a root context derived from the server lifecycle and its canonical cancellation function;
- exact-once cancellation arbitration;
- current ephemeral state;
- bounded safe recent events and bounded terminal result;
- bounded subscriber registrations;
- creation, transition, terminal, and eviction timestamps.

The hub does not own workflow steps, retries, runtime policy, lock policy, verdicts, artifact validation, durable state, or recovery decisions. Those remain app/product/runtime responsibilities.

The normal flow is:

1. The handler validates session, CSRF, body, and confirmation policy.
2. The app repeats normalization and fingerprint validation.
3. The hub atomically rejects draining/capacity or creates an `accepted` record and server-derived context.
4. One hub-owned goroutine calls the synchronous app execution method with that context and a non-blocking safe event sink.
5. The app acquires the product lock, starts existing use cases, propagates context through nested runtime/harness/process work, and reconciles durable state.
6. The hub serializes progress publication, assigns monotonic per-operation sequence IDs, and fans out without allowing subscribers to block the producer.
7. Exactly one terminal result wins among completion, failure, timeout, explicit cancellation, or shutdown cancellation.
8. The app completes cleanup/reconciliation and releases product locks; the hub publishes the terminal projection and retains it briefly before eviction.

There is no detached mode and no web-owned work queue. Capacity is reserved before returning acceptance; work begins immediately. If capacity or draining prevents ownership, start fails before an operation is accepted.

### Cancellation ownership

Cancellation is a request, not a terminal-state assignment.

- Explicit browser cancellation and graceful shutdown call the same exact-once cancellation function stored in the hub.
- The cancel function cancels the app operation root context. Existing app/product/runtime/harness/process code remains responsible for stopping scheduling, cancelling nested calls, terminating owned process trees, persisting truthful state, and reporting cleanup metadata.
- Repeated cancellation is idempotent. It may update neither reason nor outcome after the first cancellation request wins.
- Browser refresh, tab closure, SSE timeout, slow-subscriber eviction, and network disconnect cancel only the subscriber/request context.
- Timeouts are represented by operation contexts with deadlines and are distinguished from explicit `user_request` and `server_shutdown` reasons.
- Terminal arbitration does not let a late cancellation overwrite an already-authoritative success/failure, and does not let a late success overwrite cancellation after product state has committed a cancelled/interrupted outcome.
- Cleanup that must continue after work cancellation may use a bounded cleanup context derived from the server shutdown deadline, not `context.Background()`. This preserves bounded ownership without abandoning cleanup when the work context is cancelled.

### Graceful shutdown and forced-stop recovery

Graceful shutdown follows this order:

1. Atomically enter `draining` before HTTP shutdown begins.
2. Reject new preparations that cannot remain valid and all new starts; prevent not-yet-started accepted work from existing by design.
3. Snapshot active operations under the hub lock, then invoke each canonical cancel function exactly once with `reason: server_shutdown` outside the lock.
4. Keep status/SSE available while operations run their canonical cleanup and durable reconciliation paths.
5. Wait for all active operation owners up to the configured cleanup deadline; do not wait while holding hub or product locks.
6. If the deadline expires, escalate owned process-tree cleanup through existing runtime/process adapters and record `interrupted` or `cleanup_uncertain` through product-owned state.
7. Release product locks only through completed reconciliation, publish final safe events/results, close subscribers, shut down HTTP, and return from the server runner.

Shutdown uses a `WaitGroup`-style owner set because operations have independent outcomes; one failure must not cancel sibling cleanup that is already being cancelled deliberately. An `errgroup` remains appropriate inside product workflows whose sibling tasks share fate, but not as the hub's global shutdown result model.

After a crash or forced termination, in-memory hub state is gone. Server startup invokes product recovery use cases before accepting mutations. Those use cases inspect authoritative flow/run state, external harness evidence where applicable, and stale lock metadata. Unresolved running work becomes interrupted or recovery-required. Old web operation IDs, event history, sessions, and confirmations are not reconstructed.

### Event and observability boundary

App/product code emits canonical safe operation events to a sink. The web adapter maps those events to SSE, the TUI maps them to messages, and CLI adapters map them to human/JSON output. Raw runtime/provider payloads remain within runtime diagnostics and are not automatically promoted.

Redaction occurs before an event or terminal result enters the hub. The retained representation cannot contain confirmation/CSRF/session values, credentials, full environment data, raw prompts/provider payloads, unrestricted stderr, executable strings, or absolute local paths. Logs and events share correlation vocabulary but remain separate projections.

Stable structured fields include request ID, web operation ID, safe product run/task ID, operation kind, project/study/sprint reference, mutation class, lifecycle transition, cancellation reason, duration, event sequence, subscriber outcome, cleanup outcome, and typed error category. The architecture adds local counters for capacity rejection, active operations, terminal outcomes, cancellation reasons, subscriber eviction, buffer rollover, and shutdown uncertainty, but no remote telemetry service.

### State classification

| State class | Owner | Examples | Restart behavior |
| --- | --- | --- | --- |
| Durable authoritative | Product modules and existing stores | workspace artifacts, `flow-state.json`, `.run-state.json`, study run state, `review.md`, `smoke.md`, external harness evidence | Re-read and reconcile. |
| Ephemeral lifecycle | Web operation hub | operation handles, cancellation functions, subscriber sets, bounded recent events/results, draining flag | Discarded; never reconstructed as authority. |
| Derived projection | App/web adapters | normalized summaries, browser result/status views, SSE payloads | Recompute from current app/product state where available. |

No database, durable web job table, event log, operation-history file, shadow lock file, or dual-write path is introduced. Existing durable schemas change only where product-owned interruption/cleanup truth requires an explicit supported state; such changes belong to the owning module and its migration/compatibility rules.

### Test architecture

The core behavior must be testable without a browser, real runtime, smoke harness, or external process. Required seams are narrow app operation fakes, fake clock/ID sources, fake runtime/harness/process cleanup, deterministic product stores/temporary workspaces, and controllable event subscribers.

Architectural tests cover:

- import direction proving `internal/web` has no direct product/runtime/process/CLI dependency;
- one shared app operation dispatch path used by CLI, TUI, and web projections;
- preparation has no mutation, lock acquisition, runtime invocation, or hub reservation;
- product locks reject cross-surface conflicts and remain held through cleanup/reconciliation;
- root context reaches nested product/runtime/harness/process work;
- browser disconnect affects only subscription while DELETE and shutdown affect operation context;
- exact-once cancellation and one terminal outcome under completion/cancel/timeout/shutdown races;
- bounded operations, buffers, subscribers, retention, and shutdown waits with slow or abandoned consumers;
- graceful shutdown ordering and forced-stop startup reconciliation;
- redaction before retention and durable-state authority after hub eviction/restart;
- race-enabled and goroutine-leak coverage for hub, locks, cancellation, and subscriber cleanup.

Behavior assertions are preferred over private map-size or goroutine-count implementation assertions, except where a public bound or leak invariant is the behavior under test.

## Trade-Offs

| Decision | Benefit | Cost / rejected alternative |
| --- | --- | --- |
| Small refactor to shared app operations before web handlers | Preserves one product core and CLI/TUI/web agreement. | Directly calling existing CLI handlers would be faster initially but couples HTTP to flags, output parsing, and process semantics. |
| Closed typed operation union | Makes authority, confirmation, validation, and safe projection auditable. | A generic command registry is more extensible but permits browser-driven behavior not governed by current use cases. |
| Synchronous app execution with asynchronous web ownership | App remains transport-neutral and easy to test; web owns only its lifecycle. | Returning an app-level job handle would spread web-like asynchronous state into CLI/TUI and risk a second scheduler. |
| Constructor injection from one composition root | Dependencies and test substitutions remain explicit. | Context-carried services or globals reduce plumbing but hide authority and contaminate tests. |
| Product-owned locks plus hub capacity bounds | Correctness holds across all surfaces while web resource limits remain local. | HTTP middleware locking cannot prevent CLI/TUI conflicts; one global process mutex destroys unrelated concurrency. |
| Exclusive per-sprint mutation domain | Safest treatment of overlapping flow, execute, review, smoke, verify artifacts and freshness. | Finer lock classes offer more concurrency but require a proven non-overlap matrix and increase deadlock/recovery risk. |
| Ephemeral hub with durable refresh | Avoids a second authority and migration burden. | Durable web jobs would survive restart but require leases, heartbeats, idempotency, ownership transfer, and reconciliation beyond this sprint. |
| No queue | Every accepted operation has immediate server ownership and clear shutdown semantics. | A queue smooths bursts but becomes scheduler state and creates cancellation/fairness/persistence questions. |
| Independent subscriber buffers with disconnect on overflow | Slow browsers cannot block product work and memory remains bounded. | Blocking delivery preserves every event but lets transport backpressure alter workflow completion; silent drops misrepresent continuity. |
| WaitGroup-style global shutdown ownership | All independent operation cleanups run and produce their own truthful outcomes. | Fail-fast `errgroup` would incorrectly make one cleanup failure cancel or obscure siblings. |
| Bounded cleanup context after work cancellation | Product cleanup can reconcile without becoming detached. | Reusing the cancelled work context prevents cleanup; `context.Background()` can outlive the server indefinitely. |
| Explicit adapters, no plugin/registry framework | Meets current three-surface need with traceable code. | Dynamic registration and generic plugin machinery are unearned and broaden the local browser's authority. |
| Real temporary workspaces plus focused fakes | Exercises filesystem locking/recovery while keeping runtime/harness deterministic. | A universal virtual filesystem would hide platform behavior and create a broad abstraction not required by the product architecture. |

## Evidence

- **Project boundary:** `projects/ultraplan-go/docs/ARCHITECTURE.md` defines `internal/web -> internal/app`, product-module ownership, ephemeral hub state, and authoritative filesystem/run state. `projects/ultraplan-go/sprints/31-web-operations/requirements.md` makes shared app use cases, product locks, server-owned cancellation, and no durable web state acceptance requirements. The architecture above specializes those contracts for Sprint 31.
- **Thin adapter and one-way dependencies:** `studies/go-cli-study/reports/final/01-project-structure.md` finds a stable thin-entrypoint/business-logic split across mature applications and cites Helm's command-to-action delegation (`helm/pkg/cmd/install.go:347`, `helm/pkg/action/install.go:73-140`) and gdu's UI abstraction (`gdu/cmd/gdu/app/app.go:30-49`). Report finding: boundary adapters should delegate inward. Sprint inference: web handlers translate transport only, while app/product use cases own operations.
- **Explicit composition and narrow seams:** `studies/go-cli-study/reports/final/03-dependency-injection.md` finds centralized manual composition, constructor injection, and core-service interfaces correlated with traceability/testability (`gh-cli/pkg/cmdutil/factory.go:16-43`, `opencode/internal/app/app.go:42-81`). It identifies global services and context service locators as hidden-coupling risks. Sprint inference: construct hub/app/server at `cmd/ultraplan`, inject clock/IDs/services, and reserve context for lifecycle.
- **Effective preparation:** `studies/go-cli-study/reports/final/04-configuration-management.md` reports that validation after all sources are merged catches invalid combinations (`k9s/internal/config/k9s.go:423-451`) and explicit override tracking avoids ambiguous precedence (`restic/internal/global/global.go:139,147`). Sprint inference: app preparation normalizes complete effective operation inputs before fingerprinting; web does not duplicate configuration logic.
- **Typed failures and safe projection:** `studies/go-cli-study/reports/final/05-error-handling.md` supports structured lock errors (`restic/internal/restic/lock.go:47`), wrapped classification, and separate user/operational rendering (`k9s/internal/model/flash.go:100-103`). Sprint inference: app/product retain typed errors while web maps only safe fields; lock and cleanup outcomes are never inferred from strings.
- **Testable volatile boundaries:** `studies/go-cli-study/reports/final/06-io-abstraction.md` finds production/test constructors and injectable I/O effective (`gh-cli/pkg/iostreams/iostreams.go:551-568`, `restic/internal/ui/mock.go:10-53`, `go-task/executor_test.go:146-151`). Sprint inference: event sinks, clocks, IDs, operations, runtimes, harnesses, and process cleanup receive deterministic fakes without replacing all filesystem behavior.
- **Context and lock lifecycle:** `studies/go-cli-study/reports/final/07-state-context.md` finds reliable cancellation where a root context reaches long-running I/O (`helm/pkg/cmd/install.go:333-347`), explicit lock coordination (`restic/internal/restic/lock.go:105`), and bounded post-cancel cleanup (`restic/internal/restic/lock.go:290-305`). It warns that `context.Background()` severs ownership. Sprint inference: server context owns web-started work, product code owns locks/reconciliation, and cleanup receives a bounded derived context.
- **Structured concurrency:** `studies/go-cli-study/reports/final/08-concurrency.md` supports localized goroutine launch, bounded fan-out (`k9s/internal/pool.go:21,30,37`), exact-once cleanup, explicit subscriber cleanup (`opencode/internal/pubsub/broker.go:67-82`), and timeout-bounded shutdown (`opencode/cmd/root.go:261-279`). Report warning: fire-and-forget and blocked channels leak or deadlock. Sprint inference: one owned goroutine per accepted operation, non-blocking bounded subscribers, and an explicit shutdown owner set.
- **Operational correlation:** `studies/go-cli-study/reports/final/10-logging-observability.md` reports benefits from stable structured fields (`k9s/internal/slogs/keys.go:6-231`), component tagging, and separation of diagnostic logs from user output. Sprint inference: app/web/product share safe correlation fields but events, logs, and durable state remain distinct representations.
- **Behavior-first verification:** `studies/go-cli-study/reports/final/11-testing-strategy.md` supports centralized mocks/fakes, HTTP test servers, behavior assertions, and selective golden fixtures, while identifying exact internal-count assertions as brittle (`k9s/internal/view/pod_test.go:23`). Sprint inference: tests prove ownership, transitions, bounds, recovery, and cross-surface agreement rather than hub implementation shape.
- **Delay unearned extension machinery:** `studies/go-cli-study/reports/final/12-extensibility.md` finds adapters/factories can provide useful evolution without dynamic loading (`dive/cmd/dive/cli/internal/command/adapter/analyzer.go:13-15`, `gh-cli/pkg/cmdutil/factory.go:16-43`) and documents global registry collision/initialization risks. Sprint inference: a closed typed capability is sufficient; no plugin, reflection, or registration framework is introduced.
- **Trust boundary and pre-retention redaction:** `studies/go-cli-study/reports/final/13-security.md` finds explicit permission/confirmation points (`opencode/internal/permission/permission.go:44-108`), type-enforced secret redaction (`restic/internal/options/secret_string.go:15-20`), and explicit argument arrays safer than shell strings. Sprint inference: browser requests select typed capabilities only and unsafe data is removed before hub retention.
- **Bounded ephemeral state:** `studies/go-cli-study/reports/final/14-performance.md` finds bounded queues/concurrency and incremental streaming preserve responsiveness (`k9s/internal/pool.go:26-48`, `opencode/internal/llm/provider/provider.go:56`) and advises against speculative pooling. Sprint inference: hub collections and subscriber queues are bounded, slow consumers are isolated, and optimization waits for measurement.

## Risks

- **App capability becomes a god interface:** Combining every operation behind one service can expose broad authority to web callers. Keep the exported request union closed, methods use-case-oriented, dependencies private, and transport composition limited to the exact selected capabilities.
- **CLI/TUI parity is only nominal:** If web dispatch uses new code while existing surfaces retain separate command logic, behavior will drift. Move orchestration behind app use cases first and add cross-surface fixtures over the same temporary workspace.
- **Lock-domain assumptions are too coarse or too weak:** One per-sprint mutation lease may reduce safe concurrency, but prematurely splitting locks can permit overlapping artifact/state writes. Start exclusive; refine only after an explicit operation/read-write matrix and deadlock analysis.
- **Lock release races cleanup:** A defer that releases on context cancellation rather than completed reconciliation can permit a second mutation while owned processes or writes remain active. Lock tests must control cleanup barriers and assert release ordering.
- **Cancellation does not reach a nested boundary:** Any fresh background context, uncancellable wait, or unmanaged subprocess breaks server ownership. Tests must trace cancellation through retries, runtimes, harnesses, child tasks, and process-tree escalation.
- **Terminal-state race:** Completion, failure, timeout, explicit cancellation, and shutdown can arrive concurrently. A single arbitration point and durable-state comparison are required so exactly one truthful terminal outcome is published.
- **Shutdown deadlock:** Waiting while holding the hub mutex, a subscriber send, or a product lock can block cancellation or reconciliation. Snapshot owners under lock, cancel/wait outside it, and make subscriber delivery non-blocking.
- **Unsafe event enters memory before redaction:** Encoding-time redaction is too late for retained replay buffers. Safe projection must occur in app/product event mapping before the hub receives data, with hostile fixtures at every boundary.
- **Ephemeral result mistaken for authority:** Retained terminal results can outlive changed product state briefly or disappear after restart. Every projection must identify durable refresh sources, and UI recovery must re-query app/product state.
- **Startup reconciliation damages valid work:** Aggressive stale-lock cleanup may release ownership while a child process remains alive or mark partial artifacts successful. Recovery must be conservative and represent uncertainty explicitly.
- **Resource bounds interact:** Active operation, event, result, subscriber, and shutdown limits multiply. Aggregate tests and race runs are required; increasing one bound requires reassessing total memory and goroutine cost.
- **Preparation causes side effects:** Runtime/harness health checks or effective-config resolution may mutate caches or launch dependencies. Preparation APIs must document and test side-effect freedom; checks that cannot be side-effect-free run only after accepted start.
- **Future durable workers conflict with this ownership model:** Allowing operations to survive server exit would require leases, heartbeats, durable idempotency, worker identity, ownership transfer, and event replay. It must be a separately versioned architecture, not a relaxation of Sprint 31 invariants.
- **Open implementation question with a fixed decision boundary:** The exact internal Go type names may follow existing app conventions, but they must preserve the decided package ownership, closed operation union, synchronous app execution contract, safe event sink, product-owned locks, and server-owned lifecycle. Naming does not reopen those decisions.
