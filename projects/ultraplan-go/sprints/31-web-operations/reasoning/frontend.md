# Frontend: Guarded Web Operations and SSE Progress

> **Inputs Used:** `projects/ultraplan-go/sprints/31-web-operations/technical-handbook.md`, `projects/ultraplan-go/sprints/31-web-operations/requirements.md`, `projects/ultraplan-go/sprints/31-web-operations/sprint-index.md`, `projects/ultraplan-go/docs/ARCHITECTURE.md`, `projects/ultraplan-go/docs/PRD.md`, `projects/ultraplan-go/docs/TRD.md`, `system/reasoning/frontend-reasoning-template.md`, `studies/go-cli-study/reports/final/05-error-handling.md`, `studies/go-cli-study/reports/final/06-io-abstraction.md`, `studies/go-cli-study/reports/final/07-state-context.md`, `studies/go-cli-study/reports/final/08-concurrency.md`, `studies/go-cli-study/reports/final/10-logging-observability.md`, `studies/go-cli-study/reports/final/11-testing-strategy.md`, `studies/go-cli-study/reports/final/13-security.md`, `studies/go-cli-study/reports/final/14-performance.md`

This area covers the browser experience for preparing, confirming, observing, cancelling, and recovering existing UltraPlan operations. The key conclusion is to keep server-rendered project, sprint, and study pages useful without JavaScript, then add small dependency-free enhancements for asynchronous start, `EventSource` progress, cancellation, and in-place durable refresh. The page never becomes a workflow engine or authority: server-issued preparation data defines what can be confirmed, operation resources describe ephemeral execution, and refreshed app/product state defines durable truth.

## Area Decisions

### Page and interaction ownership

- Existing project, sprint, and study detail pages own operation entry points because users need current readiness, validation, artifact, and recovery context before acting. A generic command console is rejected.
- Each eligible action submits an allowlisted operation specification to `POST /api/v1/operations/prepare`. The browser does not calculate affected paths, mutation class, fingerprints, runtime/model or harness identity, prerequisites, or confirmation expiry.
- Handlers map app results into explicit page and component view models before rendering. Templates do not read files, call use cases, classify errors, interpret workflow state, or construct operation specifications from arbitrary values.
- Presentation remains embedded Go `html/template`, layered CSS, and minimal dependency-free JavaScript. Sprint 31 extends the established templates and assets without introducing Node.js, a frontend framework, a client router, a service worker, or a build step.
- Reusable presentation is limited to operation-specific status, confirmation, progress, finding, result, error, and recovery components. Product-specific fields remain in route-level view models rather than forcing a generic dynamic form or schema renderer.

### Progressive confirmation flow

The guarded flow has four explicit steps:

1. The user selects an existing operation from a project, sprint, or study page and submits its server-rendered form.
2. The server prepares the normalized request and renders a review panel showing operation kind, target scope, workspace-relative affected paths, mutation class, runtime/model or harness summary when relevant, prerequisites, and expiry.
3. The user explicitly confirms the reviewed preparation. The start request repeats the original operation specification and sends the opaque confirmation token; hidden fields are transport values, not browser authority.
4. An accepted start redirects or progressively updates to an operation view containing its status URL, progress stream, cancellation control, and durable refresh target.

The confirmation control is disabled after submission to reduce accidental double activation, but the server's single-use token remains the actual replay defense. Expired, mismatched, replayed, or stale confirmations replace the review panel with a specific explanation and a primary `Prepare again` action. They never silently re-prepare and start, because changed scope or inputs require renewed user review.

Without JavaScript, prepare and start use ordinary form posts and server-rendered responses; the operation page uses status refresh links and a normal cancellation form using the server's method-preserving HTTP mechanism. With JavaScript, the same forms use the versioned JSON API for in-place prepare/start/cancel behavior. Enhancement must not create a second request schema or bypass normal Host, Origin, session, CSRF, body, and confirmation checks.

### Browser state model

Browser state is deliberately small and derived from server resources:

| State | Browser presentation | Primary action |
| --- | --- | --- |
| `idle` | Current readiness and available allowlisted actions. | Prepare an operation. |
| `preparing` | Busy status attached to the initiating control. | Wait; submission is disabled. |
| `confirmation_required` | Normalized scope and impact summary with expiry. | Confirm or cancel locally. |
| `accepted` / `running` | Operation identity, current phase, latest safe progress, elapsed time, and connection status. | Request cancellation or refresh status. |
| `cancelling` | Cancellation reason and a clear statement that cleanup/reconciliation is still running. | Wait or refresh; cancellation remains idempotent. |
| `succeeded` | Safe result summary, affected artifacts/findings, and durable-state link. | Refresh the owning page or inspect output. |
| `failed` | Typed safe failure, actionable findings, and retry/recovery guidance. | Correct inputs, inspect durable state, or prepare again. |
| `cancelled` | Explicitly cancelled outcome and cleanup result. | Refresh durable state before rerunning. |
| `interrupted` | Work did not reach a clean terminal result, including restart recovery. | Inspect durable state and follow recovery guidance. |
| `cleanup_uncertain` | Prominent warning that owned work or locks may need reconciliation. | Follow the provided recovery action; do not imply rerun safety. |

The operation's server state wins over local JavaScript flags. The client may retain only transient DOM state such as pending submission, the last rendered event ID, and reconnect status. It must not persist operation truth, confirmation tokens, results, event history, or workflow state in local storage, IndexedDB, cookies, or a service worker. URLs may carry only safe route identifiers such as an opaque operation ID; secrets and confirmation values never enter query strings or fragment state.

### Progress and findings

- The operation view subscribes to `GET /api/v1/operations/{id}/events` only after start acceptance. A server-rendered current snapshot is visible before the stream connects, avoiding an empty JavaScript-only shell.
- `snapshot` replaces the current ephemeral status projection. `progress` updates a bounded recent activity list and current phase. `warning`, `finding`, and `artifact` render distinct safe cards or links. `cancel_requested` moves the view to cancelling without claiming completion. `terminal` replaces progress controls with the final result and durable refresh action.
- The UI displays a bounded recent timeline, not an infinite console. Messages are rendered as text; event payloads are never inserted as raw HTML. Raw prompts, provider events, stdout/stderr, executable strings, environment values, credentials, and absolute paths have no browser rendering path.
- Findings preserve severity, safe location/reference, summary, and remediation where supplied by the app DTO. Severity uses text and iconography as well as color. Artifact links use existing allowlisted bounded previews rather than direct filesystem URLs.
- Progress is truthful but non-authoritative. The view labels it as recent activity and always offers a durable-state refresh after completion, cancellation, interruption, eviction, rollover, or reconnect uncertainty.

### Reconnect, gaps, and browser disconnect

`EventSource` automatic reconnect is allowed because SSE remains observation-only. The server and browser use the last event ID to resume retained events. Connection presentation distinguishes `Live`, `Reconnecting`, `History incomplete`, and `Stream ended`; it never maps a network state to operation success or failure.

On `recovery_required`, the browser:

1. marks the recent timeline incomplete rather than erasing or fabricating events;
2. applies the accompanying current snapshot;
3. exposes the durable refresh action prominently; and
4. continues live observation when the server permits it.

When a stream reaches its maximum lifetime or disconnects, JavaScript reconnects with bounded backoff while the operation is non-terminal and falls back to the status endpoint after repeated failures. Page refresh, navigation, tab closure, offline state, and subscriber eviction stop only observation. The UI never calls cancellation from unload handlers and never tells the user that closing the page stops work.

If an operation record is evicted or a server restarts, `operation_not_found` is presented as unavailable ephemeral history, not proof that the operation never ran or failed. The page follows the response's safe durable refresh guidance to reload the owning project, sprint, or study state.

### Cancellation and graceful shutdown

- Cancellation is an explicit, labelled control separate from browser navigation. It requires same-origin mutation protection and a deliberate confirmation naming the operation and target scope.
- After `DELETE`, controls transition to `Cancellation requested`; the result remains pending until app/product cleanup and reconciliation produce a terminal outcome. Repeated activation is disabled locally but remains safe because the endpoint is idempotent.
- A `server_shutdown` reason is rendered as `Server is shutting down` rather than as user cancellation or generic failure. The operation view explains that UltraPlan is stopping active server-owned work and reconciling durable state.
- During draining, preparation/start controls are disabled when the page knows the state, and `server_draining` responses provide a restart/retry message. Existing operation views remain observable until their terminal, interrupted, or cleanup-uncertain event is rendered and the stream closes.
- A stream ending during forced termination is only `Connection lost`; the browser must not invent `server_shutdown`. After the server returns, durable refresh can report recovery-required or interrupted state.

### Loading, empty, failure, and recovery views

- **Loading:** Keep the current page stable, attach `aria-busy` to the affected region, show text describing preparation, start, cancellation, or refresh, and avoid replacing the whole page with a spinner.
- **Empty progress:** Show the operation state and `Waiting for progress` rather than an empty box; lack of events does not mean the operation is stalled.
- **Validation or prerequisite failure:** Render safe field/finding summaries next to the operation panel and link to the owning artifact/readiness view. Do not expose wrapped internal errors.
- **Conflict:** Name the safe conflicting operation kind/scope and offer status or wait guidance. Do not imply that hiding/disabling the browser control is concurrency enforcement.
- **Runtime/internal failure:** Show the stable safe message and request/operation correlation ID, with logs or CLI diagnostics as the recovery path when supplied; never show unrestricted stderr or stack traces.
- **Success:** Show the final typed summary, findings/artifacts where applicable, completion time, and a durable-state refresh link. A successful transport response alone never triggers success styling.
- **Interruption/uncertainty:** Use a persistent warning panel with explicit recovery steps. Do not offer an enabled rerun until refreshed product state says mutation ownership and prerequisites are safe.

### Accessibility and interaction behavior

- All operation paths use native forms, buttons, headings, links, lists, tables, and `<dialog>` only if its no-JavaScript fallback remains usable. Div-based custom controls are rejected.
- Preparation errors move focus to an error summary linked to affected fields. Opening the confirmation view moves focus to its heading; cancelling it returns focus to the initiating action. Terminal transitions move focus only when user initiated or when necessary, avoiding disruptive focus theft during routine progress.
- A polite status live region announces phase changes, reconnect state, cancellation requested, and terminal outcome. Individual high-frequency progress events are not all announced. Failures and cleanup uncertainty use an assertive alert once.
- Buttons retain visible focus, have operation-specific labels, and expose disabled/pending state in text. Color is never the sole signal for severity, stream state, or outcome.
- Progress animation is nonessential and respects `prefers-reduced-motion`; no auto-scrolling steals reading position. The bounded activity region can be reached and read by keyboard.
- Confirmation expiry is shown as an absolute local time plus concise remaining-time text, but server rejection remains authoritative if client clocks differ.

### Performance and verification

No framework or asset pipeline is added, so bundle impact remains limited to operation and SSE enhancement modules. JavaScript attaches behavior through data attributes and event delegation, opens at most one stream for the visible operation view, bounds rendered activity to the server's recent-event window, and replaces or removes old nodes to prevent an unbounded DOM. Reconnect must not create duplicate streams or duplicate event rendering.

Verification uses `httptest`, deterministic template view models, fake app operations, fake clocks/IDs, and controlled SSE writers. Required coverage includes:

- no-JavaScript prepare, confirm, start, status refresh, and cancellation paths;
- enhanced prepare/start/cancel requests using the same DTOs and security policy;
- exact rendering for every lifecycle state, cancellation reason, error category, empty state, finding severity, and durable recovery action;
- monotonic event application, duplicate suppression, reconnect, rollover recovery, terminal close, slow subscriber eviction, and browser disconnect isolation;
- graceful-shutdown draining, `server_shutdown`, interrupted, and cleanup-uncertain presentations without false success;
- keyboard order, focus restoration, labels, live-region behavior, color-independent meaning, and reduced motion;
- hostile event/artifact/error strings proving escaped text, safe URLs, path containment, and absence of secrets or raw diagnostics;
- bounded timeline DOM, one active subscription, stable layout during updates, and CLI/TUI/browser agreement on readiness and terminal outcomes.

Use semantic assertions for interaction and state behavior. Keep representative template/golden fixtures for complete confirmation, running, findings, failure, and recovery panels where full-output review catches accidental presentation drift.

## Trade-Offs

| Decision | Benefit | Cost / rejected alternative |
| --- | --- | --- |
| Server-rendered baseline with narrow progressive enhancement | Every guarded operation remains understandable and recoverable without a client application or build runtime. | A client SPA could make transitions smoother but would duplicate routing/state ownership and is not justified by current complexity. |
| Operation entry points on owning detail pages | Preserves readiness and scope context before confirmation. | A universal command palette is more compact but obscures product context and invites a generic browser command surface. |
| Server-issued review panel | Users confirm normalized current scope rather than browser guesses. | Client-generated summaries feel faster but can diverge from app normalization, fingerprints, and effective runtime/harness selection. |
| Explicit confirm after stale/expired preparation | Preserves meaningful consent when inputs change. | Automatic re-prepare/retry is convenient but could start work the user did not review. |
| Ephemeral operation view plus durable refresh | Keeps browser progress useful without treating retained events as workflow state. | Persisting client history improves continuity but creates stale competing authority and disclosure risk. |
| Bounded recent timeline | Protects DOM and memory while showing useful progress. | An unbounded console preserves more detail but grows indefinitely and encourages unsafe raw-output projection. |
| Reconnect with visible gap state | Live monitoring can continue without pretending missing events were delivered. | Unlimited replay or silent gap recovery would respectively violate bounds or mislead the user. |
| Explicit cancellation only | Navigation cannot accidentally terminate work and matches server lifecycle ownership. | Cancelling on page close appears intuitive to some users but is unreliable and conflates subscription with execution. |
| Status text plus restrained live announcements | Keeps screen-reader users informed without announcing every event. | Announcing every progress item is more exhaustive but produces unusable noise during active operations. |
| Selective template fixtures and semantic interaction tests | Protects complete critical views while keeping lifecycle tests robust. | All-golden tests are brittle; only small assertions can miss unsafe or misleading full-page combinations. |

## Evidence

- **Product and project contracts:** `projects/ultraplan-go/sprints/31-web-operations/requirements.md` requires server-rendered confirmation, progress, findings, failure, explicit/server-shutdown cancellation, reconnect, and durable recovery views. `projects/ultraplan-go/docs/PRD.md` lines 399-406 and `projects/ultraplan-go/docs/TRD.md` lines 2101-2133 define embedded templates, ordinary HTTP commands, SSE observation, explicit cancellation, browser-disconnect independence, and no frontend runtime. Sprint inference: build a progressively enhanced operation panel over the versioned app/web boundary rather than client-owned workflow state.
- **Presentation ownership:** `projects/ultraplan-go/docs/ARCHITECTURE.md` lines 329-365 assigns HTML/view models, embedded assets, confirmations, and subscribers to `internal/web`, while forbidding templates from reading files or deciding product state. Sprint inference: handlers produce explicit view models and route pages compose bounded operation components; JavaScript only enhances transport interactions.
- **Safe error presentation:** `studies/go-cli-study/reports/final/05-error-handling.md` finds typed errors useful for actionable data and cites separation of operational and user messages (`restic/internal/restic/lock.go:47`, `k9s/internal/model/flash.go:100-103`). Sprint inference: UI branches on stable error codes and renders allowlisted details, not message-string parsing or wrapped diagnostics.
- **Deterministic surface testing:** `studies/go-cli-study/reports/final/06-io-abstraction.md` identifies in-memory/test constructors and synchronized capture as effective (`gh-cli/pkg/iostreams/iostreams.go:551-568`, `go-task/executor_test.go:146-151`). Sprint inference: fake operation sources, clocks, and controlled streams test pages and enhancement without a real runtime or browser service.
- **Lifecycle separation:** `studies/go-cli-study/reports/final/07-state-context.md` finds reliable cancellation where one root context reaches long-running work (`helm/pkg/cmd/install.go:333-347`) and warns against severing ownership. Sprint inference: the browser displays server/app lifecycle but never derives operation cancellation from the SSE request or page lifecycle.
- **Bounded concurrency and subscriptions:** `studies/go-cli-study/reports/final/08-concurrency.md` supports explicit subscriber cleanup and bounded fan-out (`opencode/internal/pubsub/broker.go:67-82`, `k9s/internal/pool.go:21,30,37`) and warns about abandoned consumers. Sprint inference: one visible subscription, bounded rendered events, explicit reconnect states, and no producer backpressure from UI behavior.
- **Progress versus diagnostics:** `studies/go-cli-study/reports/final/10-logging-observability.md` supports stable structured fields and separation of user output from logs (`k9s/internal/slogs/keys.go:6-231`, `opencode/internal/logging/logger.go:25-62`). Sprint inference: show safe correlation IDs and user progress while keeping logs, provider detail, and raw stderr out of the page.
- **Behavior-focused coverage:** `studies/go-cli-study/reports/final/11-testing-strategy.md` supports centralized fakes, HTTP test servers, behavior assertions, and selective goldens. Sprint inference: test state transitions, focus, reconnect, cancellation, bounds, and recovery semantically, with complete fixtures only for critical rendered projections.
- **Explicit trust boundary and redaction:** `studies/go-cli-study/reports/final/13-security.md` cites explicit permission/confirmation points and type-enforced secret protection (`opencode/internal/permission/permission.go:44-108`, `restic/internal/options/secret_string.go:15-20`). Sprint inference: the browser confirms server-normalized scope, uses typed operations, renders text safely, and never stores tokens or unsafe payloads as client state.
- **Incremental bounded rendering:** `studies/go-cli-study/reports/final/14-performance.md` finds bounded streaming and slow-consumer isolation preserve responsiveness (`opencode/internal/llm/provider/provider.go:56`, `k9s/internal/pool.go:26-48`). Sprint inference: cap event retention and DOM growth, isolate reconnect, and reject a framework or speculative optimization without measured need.

## Risks

- **Confirmation becomes ceremonial:** If the panel omits normalized scope, affected paths, mutation class, runtime/harness summary, or expiry, users cannot make an informed decision. Template tests must assert every applicable summary field and stale-confirmation re-review.
- **Enhanced and non-enhanced paths diverge:** Separate request construction or error handling can make JavaScript behavior less safe than normal forms. Both paths must share handlers, DTO mapping, confirmation checks, and representative fixtures.
- **Local UI state overrides server truth:** Optimistic success or cancellation can misreport work still cleaning up. Only server snapshots/terminal results may select terminal presentation; local state is limited to transport pending/reconnecting flags.
- **Reconnect duplicates or reorders progress:** Automatic reconnect plus manual fallback can create two streams or replay duplicate IDs. The client must maintain one subscription, ignore already-applied IDs, and treat gaps as incomplete history rather than failure.
- **Browser disconnect accidentally cancels work:** Unload handlers, shared abort controllers, or request-context coupling could violate ownership. Interaction and integration tests must close streams/pages while proving the fake app operation continues.
- **Shutdown is rendered as ordinary failure:** Generic connection-loss handling can hide graceful cancellation or falsely attribute forced termination. Render `server_shutdown` only from server data and use neutral connection-loss guidance otherwise.
- **Cancellation is mistaken for completion:** Immediately showing `cancelled` after DELETE can allow unsafe reruns while cleanup or locks remain active. Keep `cancelling` until authoritative terminal reconciliation and suppress rerun when cleanup is uncertain.
- **High-frequency events harm accessibility or responsiveness:** Appending and announcing every event can grow the DOM, shift layout, and flood assistive technology. Bound/reuse nodes, summarize phases, and announce only meaningful transitions.
- **Unsafe content reaches HTML:** Event, finding, artifact, and error values may contain hostile markup, secrets, or paths. Use escaped template/text insertion, allowlisted URLs/fields, pre-retention redaction, and hostile fixtures; never use `innerHTML` for server data.
- **Operation controls imply authorization:** Hidden or disabled controls are only affordances. Server Host/Origin/session/CSRF/confirmation, allowlist, fingerprint, lock, and capacity checks remain mandatory for every request.
- **Ephemeral history appears durable:** A polished timeline can be mistaken for the complete execution record. Label it as recent activity and expose durable refresh after gaps, terminal states, eviction, and restart.
- **Framework-like JavaScript accretes:** Ad hoc global stores, client routing, schema-driven forms, or generalized components could become a second application layer. Keep modules capability-focused and revisit architecture only when measured interaction complexity justifies a separately governed change.
- **Accessibility hardening is deferred too broadly:** Sprint 32 owns release-level audit, but Sprint 31 still must provide native controls, keyboard operation, focus/error behavior, status announcements, reduced motion, and non-color meaning for the new critical workflow.
