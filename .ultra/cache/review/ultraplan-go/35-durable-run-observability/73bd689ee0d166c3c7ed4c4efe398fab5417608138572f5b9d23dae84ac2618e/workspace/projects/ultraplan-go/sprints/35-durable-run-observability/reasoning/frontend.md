# Frontend: Workspace Run Observation and Control

> **Inputs Used:** `projects/ultraplan-go/sprints/35-durable-run-observability/technical-handbook.md`, `projects/ultraplan-go/sprints/35-durable-run-observability/requirements.md`, `projects/ultraplan-go/sprints/35-durable-run-observability/sprint-index.md`, `projects/ultraplan-go/sprints/35-durable-run-observability/code-context.md`, `projects/ultraplan-go/docs/ARCHITECTURE.md`, `projects/ultraplan-go/docs/PRD.md`, `projects/ultraplan-go/docs/TRD.md`, `projects/ultraplan-go/sprints/35-durable-run-observability/reasoning/api-design.md`, `studies/go-cli-study/reports/final/05-error-handling.md`, `studies/go-cli-study/reports/final/07-state-context.md`, `studies/go-cli-study/reports/final/08-concurrency.md`, `studies/go-cli-study/reports/final/09-terminal-ux.md`, `studies/go-cli-study/reports/final/10-logging-observability.md`, `studies/go-cli-study/reports/final/11-testing-strategy.md`, `studies/go-cli-study/reports/final/13-security.md`, `studies/go-cli-study/reports/final/14-performance.md`, `system/reasoning/frontend-reasoning-template.md`, `../ultraplan-go/internal/web/static/app.js`, `../ultraplan-go/internal/web/operation_handlers.go`, `../ultraplan-go/internal/web/operations.go`, `../ultraplan-go/internal/web/operations_contract_test.go`

This area defines browser and TUI presentation of durable runs. It preserves the existing server-rendered, progressively enhanced web architecture and treats the durable run API as server state. The browser is never the run store, terminal arbiter, or cancellation owner.

## Area Decisions

### 1. Surface and ownership

Add canonical HTML routes:

- `/runs`: workspace run list
- `/runs/{run_id}`: durable run detail, retained history, control, and recovery

Existing `/operations/{id}` routes resolve or redirect to the durable run page, tombstone, or a precise legacy recovery page.

Keep the current dependency-free JavaScript and `html/template` hierarchy. Handlers map app run DTOs into explicit page/component view models. Templates render complete no-JavaScript snapshots and never read the workspace, interpret product state, or subscribe directly to storage.

Use existing presentation layers:

- primitives: lifecycle/liveness badges, icon+text status, copy button, timestamp, progress meter, action button
- components: active-run summary, run filters, run card/row, identity panel, progress panel, attempt list, event timeline, diagnostics/recovery panel, cancellation panel
- layouts: workspace list and detail layouts
- pages: run list and run detail composition

Do not add a frontend framework, client router, global state store, service worker, WebSocket layer, or JavaScript build pipeline.

### 2. Workspace-wide top bar and list

The top bar reads the canonical workspace active projection, not browser-session operations or planning-stage status. It displays:

- text `No active runs`, `1 active run`, or `<n> active runs`
- a link to `/runs?lifecycle=active`
- a visible degraded indicator when run-control health cannot be read

The full run list defaults to newest retained runs and provides URL-backed filters for active, needs-attention, terminal, project, sprint, study, stage, and operation kind. URL state makes views shareable and functional without JavaScript. `needs-attention` includes stalled, cancellation-uncertain, interrupted, cleanup-uncertain, and storage/reconciliation diagnostics.

Desktop uses a compact table with stable columns: lifecycle, target, operation, progress, liveness, accepted/start time, and action. Below 720px, rows become stacked cards without horizontal scrolling. Active runs sort before equally recent terminal runs only when the active filter is selected; otherwise canonical API order is preserved.

The top bar refreshes immediately on page load/focus and every 5 seconds while the document is visible. It pauses while hidden. Run-list enhancement uses the same interval only when no detail SSE connection already provides a relevant terminal transition. Polling failure retains the last rendered count with a `Status unavailable` label; it never displays zero as a fallback.

### 3. Stable run detail

The server-rendered detail page contains:

1. identity and target: run ID, attempt, operation, workspace-relative project/sprint/study/stage/task links
2. lifecycle: state, accepted/queued/started/finished timestamps, immutable terminal outcome
3. liveness: live, stalled, owner unreachable, interrupted, cleanup uncertain, or terminal with heartbeat age and guidance
4. progress: stage/task label, completed/total when known, runtime attempts, duration, usage/cost only when known
5. attempts and safe correlations: bounded owner/runtime/agentwrap/product/harness references
6. retained timeline: sequence, time, type, safe summary, omission/compaction markers
7. diagnostics and recovery: typed code, explanation, next action, retention bounds, storage/reconciliation health
8. cancellation: current request/acknowledgement/uncertainty plus the authorized action

Lifecycle and liveness are separate visual concepts. `running + stalled` is not rendered as failed. `cancelling` means a request exists, not that the worker stopped. `cleanup_uncertain` is terminal but does not claim cleanup or product failure. Product artifact/stage status appears as a separately labeled linked authority.

Every state uses text and an icon or shape in addition to color. Unknown future lifecycle values render as `Unknown state` with the explicit API `terminal` flag deciding whether controls remain available.

### 4. Initial render and live-follow model

The server renders the current durable snapshot and retained timeline first. The page includes `data-run-id`, `data-last-sequence`, and `data-oldest-retained-sequence` attributes. JavaScript then opens:

```text
/api/v1/runs/{run_id}/events?after={last_rendered_sequence}
```

This closes the server-render/SSE race because every event committed after the rendered cursor is replayed before live follow. Refresh, a new browser session, or another server repeats the same sequence from durable state; no local browser cache is required for correctness.

JavaScript deduplicates timeline entries by numeric sequence and rejects an event whose run ID differs from the page. It appends only through DOM text APIs. It never uses event HTML or raw provider fields.

Connection states are explicit:

| State | Presentation |
| --- | --- |
| Connecting | `Connecting to live updates` beside the durable snapshot. |
| Live | `Live` with last committed sequence/time. |
| Reconnecting | `Connection interrupted; retained state is safe. Reconnecting.` |
| Replay gap | Persistent gap panel with unavailable range, reason, current snapshot, and retained-history/follow choices. |
| Subscriber capacity | Retry countdown from `Retry-After`; durable page remains usable. |
| Terminal | Close stream, fetch final snapshot, update result/recovery, and stop reconnecting. |
| Store unavailable | Preserve rendered state, label it stale with generation time, and provide diagnostics/retry. |

On EventSource error, the client fetches run detail. If its requested cursor is below the durable lower boundary, it renders the typed replay gap and lets the user choose `Show retained history` or `Follow from current state`. It does not announce transient transport loss as run failure.

### 5. Timeline and bounded detail

Render at most 200 timeline rows initially. Older retained rows load through an explicit `Load earlier retained events` action, 200 at a time. Live rows append at the end and the DOM retains at most 500 rows; removing old DOM rows adds a presentation-only marker and does not alter the durable cursor or imply storage loss.

Equivalent high-frequency progress updates update the current progress panel and are summarized in the timeline rather than causing unbounded row growth. Warnings, findings, artifacts, cancellation, recovery, omission, and terminal events always receive visible entries.

The viewport does not auto-scroll if the user has moved away from the newest event. A `New updates` button appears and moves focus/scroll only when activated. At the bottom, live updates remain visible without animation. Respect `prefers-reduced-motion`; no pulsing status or smooth-scroll requirement is added.

### 6. Cancellation interaction

Show `Request cancellation` only for an active run and a current authorized browser session. Activation opens a concise confirmation dialog containing run ID, operation, target, current state, and the statement that completion may win the race.

After authorized `DELETE`:

- `202` changes the panel to `Cancellation requested` and keeps following events
- duplicate `200` displays the current cancellation state without another optimistic update
- terminal `200` renders the immutable winner
- stale/unreachable owner displays `Cancellation delivery uncertain` and reconciliation guidance
- authorization or persistence failure restores the button and focuses an actionable error

The button never changes lifecycle to `cancelled` locally. Browser close, navigation, refresh, and SSE disconnect never invoke cancellation. Retry/resume remain product-owned guarded operations rather than run-detail buttons in this sprint.

### 7. Loading, empty, error, success, and recovery states

| Condition | Required UI |
| --- | --- |
| Initial list loading | Server-rendered list or a small status line; no full-page spinner. |
| No retained runs | Explain that accepted runtime work will appear here and link to relevant commands/docs. |
| No active runs | Calm `No active runs` state, not an error. |
| Filter has no matches | Preserve filters and offer `Clear filters`. |
| Run not found | Distinguish malformed/current unknown ID and link to workspace runs/diagnostics. |
| Tombstone | Show known identity, target, terminal outcome if known, expiry, and recovery; do not show an empty timeline as complete history. |
| Legacy operation not retained | Explain the pre-durable limitation and link to workspace/product status, never generic `Operation not retained`. |
| Replay gap | Show exact available boundary and current durable snapshot. |
| Stalled owner | Warning state with heartbeat age and conservative guidance; no success/failure inference. |
| Cleanup uncertainty | Terminal attention state with product-status and diagnostics links. |
| Storage degraded | Mark snapshot generation time and disable mutations whose persistence cannot be guaranteed. |
| Terminal success/failure/cancel | Show immutable operational outcome and then separately refresh product state. |

### 8. Accessibility

The server-rendered page provides a logical heading hierarchy, skip link, landmarks, table headers, list semantics, and ordinary links/forms before JavaScript runs.

Keyboard and focus rules are:

- all filters, rows, copy controls, timeline paging, recovery choices, and cancellation are reachable in document order
- opening cancellation moves focus into the dialog and traps it until cancel/submit
- closing the dialog returns focus to the invoking button
- after start redirect, focus lands on the run-detail heading
- newly appended events do not steal focus
- `New updates` moves to the first unseen event only on activation

Use one throttled `aria-live="polite"` summary for connection, progress milestones, cancellation acknowledgement, and terminal state. Do not announce every event. Blocking mutation errors use `role="alert"`. Labels include text; color is supplementary. Timestamps use `<time datetime>` and display both absolute local time and an accessible relative description where useful.

### 9. CLI, TUI, and browser agreement

Use one presentation vocabulary across surfaces:

- lifecycle labels come from the canonical lifecycle projection
- liveness is shown separately
- unknown totals remain `unknown`, not zero
- `cancellation requested` and `cancellation acknowledged` are not terminal labels
- replay gaps and omitted detail are explicit
- product outcome and operational outcome remain separately named
- recovery guidance uses the same typed diagnostic code and next action

TUI receives the same app snapshots/events and keeps its own bounded message channel only for rendering. Dropped TUI delivery triggers a durable refresh by sequence. The TUI model does not retain a competing run identity or terminal state.

### 10. Security and privacy

Render only allowlisted safe fields from explicit view models. Use `textContent` for event detail, keep existing hostile-Markdown protections, and do not construct links from arbitrary event paths.

Do not show raw prompts, native provider payloads, headers, cookies, environment, unrestricted stdout/stderr, absolute home paths, raw process command lines, or unverified artifact paths. Copy buttons copy only safe opaque IDs or workspace-relative references.

Read visibility is not session ownership. Mutation still requires current same-origin session and CSRF proof. The UI must not hide a readable run merely because the session rotated, and must not retain mutation authority in local storage.

### 11. Testing

Required frontend verification includes:

- server rendering for full, compacted, tombstone, active, stalled, interrupted, cleanup-uncertain, and terminal snapshots
- no-JavaScript navigation, filters, detail, recovery, and cancellation form behavior
- active count agreement for two CLI-started runs across two local servers
- session rotation and observer restart without disappearing detail
- server-render/SSE race, duplicate sequences, reconnect, 512-event catch-up reconnect, replay gap, slow subscriber, and terminal refresh
- cancellation confirmation, keyboard focus, duplicate request, terminal race, uncertain delivery, CSRF rejection, and persistence failure
- unknown lifecycle/event values and explicit terminal fallback
- hostile event/diagnostic/Markdown content, path redaction, and no unsafe DOM insertion
- mobile card layout, narrow detail layout, zoom, keyboard-only use, focus visibility, live-region throttling, reduced motion, and non-color state recognition
- bounded timeline DOM and paused hidden-document polling
- API fixture, TUI model, browser integration, and gated real-runtime dogfood agreement

Tests use `httptest`, fake app capabilities, deterministic clocks/events, and a browser engine for behavior that DOM unit tests cannot prove. Assertions target visible labels, focus, actions, and durable cursor behavior rather than internal JavaScript variables.

## Trade-Offs

| Decision | Benefit | Cost and rejected alternative |
| --- | --- | --- |
| Server-rendered pages plus narrow JS | Durable snapshots work on refresh, session rotation, errors, and no-JS paths. | Some interaction code remains manual. Rejected: SPA/global client store because it would duplicate authority and add a build system. |
| Canonical workspace run page | One discoverable place for CLI/TUI/web work and attention states. | Adds routes/components and filter complexity. Rejected: only embedding runs on owning project/study pages. |
| Separate lifecycle and liveness | Truthfully represents `running but stalled` and uncertainty. | More labels for users to learn. Rejected: collapsing all uncertainty into failed/running badges. |
| Server snapshot before SSE | Fast meaningful first paint and no dependence on browser cache. | Requires cursor handoff and final fetch. Rejected: empty page waiting for SSE. |
| Typed gap panel | Users know history is partial and retain a current snapshot. | More recovery UI. Rejected: silently replaying from the oldest available event. |
| Bounded timeline DOM | Keeps long-run pages responsive. | Older visible rows require paging and may leave the DOM. Rejected: unbounded append. |
| Poll top bar every 5 seconds | Cross-process discovery without a workspace-wide browser socket. | Counts are briefly stale. Rejected: per-page in-memory counts and a permanent global SSE connection before evidence warrants it. |
| Confirmation before cancellation | Makes the mutation and race semantics explicit. | One extra interaction. Rejected: immediate destructive click or browser disconnect cancellation. |
| No automatic scroll away from user position | Preserves reading and keyboard context. | Requires a `New updates` action. Rejected: forced scroll on every event. |
| Shared vocabulary, surface-specific rendering | CLI, TUI, and browser agree without sharing framework types. | Mapping tests are required. Rejected: parsing CLI output or putting browser DTOs in product modules. |

## Evidence

### Repository evidence

- `projects/ultraplan-go/sprints/35-durable-run-observability/code-context.md` identifies the current top bar as a session-local `/api/v1/operations` projection and the operation page/SSE client as process-memory dependent. The workspace list and durable detail directly repair those observed seams.
- `../ultraplan-go/internal/web/static/app.js` currently starts EventSource from an operation ID, announces reconnect, reloads on terminal, and uses a direct cancellation request. The selected design preserves progressive enhancement while changing the source, cursor handoff, gap recovery, and cancellation truthfulness.
- `../ultraplan-go/internal/web/operation_handlers.go` contains the generic `Operation not retained` page and session-filtered detail/cancellation. Stable run/tombstone/recovery views and separate read/mutation policy replace that behavior.
- `../ultraplan-go/internal/web/operations_contract_test.go` freezes existing operation lifecycle values and event names. The frontend keeps compatibility projection while using richer canonical run/liveness fields on new pages.
- `projects/ultraplan-go/docs/ARCHITECTURE.md` requires page-to-layout-to-component-to-primitive dependencies, explicit view models, dependency-free JavaScript, and no browser authority. The component and state placement follows that rule.
- `projects/ultraplan-go/sprints/35-durable-run-observability/reasoning/api-design.md` defines canonical run resources, explicit terminal flags, typed replay gaps, tombstones, workspace-readable history, and authorized idempotent cancellation. The browser behavior maps those contracts without inventing alternate state.

### Selected report findings

- `studies/go-cli-study/reports/final/05-error-handling.md` supports machine-inspectable failures plus safe actionable user guidance. This informs distinct gap, tombstone, storage, owner, authorization, and cancellation states instead of one generic error.
- `studies/go-cli-study/reports/final/07-state-context.md` distinguishes durable identity from cancellation context and notes explicit sessions/locks for long-running work. This supports refresh/session-independent reads and mutation controls that use fresh authority.
- `studies/go-cli-study/reports/final/08-concurrency.md` supports bounded queues, localized lifecycle, timeout cleanup, and slow-consumer isolation. This informs bounded DOM/event queues and reconnect rather than backpressure on execution.
- `studies/go-cli-study/reports/final/09-terminal-ux.md` finds progressive streaming, interruptibility, non-TTY fallback, calm progress, and explicit cancellation central to long operations. This supports server-first content, live milestones, no-JS behavior, and truthful cancellation states.
- `studies/go-cli-study/reports/final/10-logging-observability.md` supports consistent structured fields and separation of user and operator detail. This informs visible safe diagnostic codes/guidance while detailed correlations remain in logs/support diagnostics.
- `studies/go-cli-study/reports/final/11-testing-strategy.md` supports behavior-level command/browser integration, fixtures, and controlled fakes. This informs browser-engine scenarios and avoids assertions about internal timers or arrays.
- `studies/go-cli-study/reports/final/13-security.md` supports explicit trust boundaries, redaction, validation, and permission separation. This informs text-only event insertion, allowlisted fields, same-origin mutation, and no local-storage authority.
- `studies/go-cli-study/reports/final/14-performance.md` supports streaming and bounded data structures for long sessions. This informs timeline paging, DOM limits, visibility-aware polling, and coalesced progress rendering.

## Risks

- The richer lifecycle/liveness vocabulary can overwhelm users. Copy, grouping, and cross-surface labels must be tested with real failure scenarios rather than adding unexplained badges.
- EventSource hides non-200 response bodies. The detail refresh on connection error is required for typed gap UX; omitting it would regress to an ambiguous reconnect loop.
- Five-second active polling creates brief count staleness and aggregate load with many tabs. Visibility pause, request de-duplication per page, and measurements are required.
- A bounded DOM can be mistaken for durable event loss. Presentation markers must clearly distinguish `not currently rendered` from repository compaction or replay gaps.
- Live updates can produce excessive screen-reader announcements. Milestone throttling and manual access to the timeline are mandatory.
- Session-independent reads increase local metadata visibility. Strict redaction and loopback/same-origin policy remain prerequisites for every run field.
- Cancellation confirmation does not guarantee delivery. The UI must preserve `requested`, `acknowledged`, `uncertain`, and terminal winner as separate states.
- Unknown event/lifecycle handling can hide newly important states if it is too generic. Unknown values must trigger a snapshot refresh and visible `Unknown state`, not silent omission.
- Server-rendered and JavaScript mappings can drift. Shared fixture tables must cover CLI JSON, app DTO, HTML model, operation compatibility, canonical run API, and TUI labels.
- Mobile card conversion can obscure sorting and timestamps. Responsive tests must verify the same facts remain visible and ordered.

The frontend decision is to proceed within the existing server-rendered architecture, adding canonical workspace run list/detail pages and narrowly scoped durable replay enhancement. The main trade-off is accepting bounded polling and explicit recovery UI to avoid a client-side authority or new frontend stack.
