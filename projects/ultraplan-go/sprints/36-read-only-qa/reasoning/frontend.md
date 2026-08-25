> **Inputs Used:** `projects/ultraplan-go/sprints/36-read-only-qa/requirements.md`, `projects/ultraplan-go/sprints/36-read-only-qa/sprint-index.md`, `projects/ultraplan-go/sprints/36-read-only-qa/technical-handbook.md`, `projects/ultraplan-go/sprints/36-read-only-qa/reasoning/architecture.md`, `projects/ultraplan-go/sprints/36-read-only-qa/reasoning/api-design.md`, `projects/ultraplan-go/docs/ARCHITECTURE.md`, `projects/ultraplan-go/docs/PRD.md`, `projects/ultraplan-go/docs/TRD.md`, `studies/go-cli-study/reports/final/05-error-handling.md`, `studies/go-cli-study/reports/final/06-io-abstraction.md`, `studies/go-cli-study/reports/final/07-state-context.md`, `studies/go-cli-study/reports/final/08-concurrency.md`, `studies/go-cli-study/reports/final/09-terminal-ux.md`, `studies/go-cli-study/reports/final/10-logging-observability.md`, `studies/go-cli-study/reports/final/11-testing-strategy.md`, `studies/go-cli-study/reports/final/13-security.md`, `studies/go-cli-study/reports/final/14-performance.md`, `system/reasoning/frontend-reasoning-template.md`, `../ultraplan-go/internal/tui/app.go`, `../ultraplan-go/internal/tui/keys.go`, `../ultraplan-go/internal/tui/model.go`, `../ultraplan-go/internal/tui/views.go`, `../ultraplan-go/internal/tui/run_view_test.go`, `../ultraplan-go/internal/web/handlers.go`, `../ultraplan-go/internal/web/operation_handlers.go`, `../ultraplan-go/internal/web/templates/shell.html`, `../ultraplan-go/internal/web/templates/sprint.html`, `../ultraplan-go/internal/web/templates/operation.html`, `../ultraplan-go/internal/web/templates/components/components.html`, `../ultraplan-go/internal/web/static/app.js`, `../ultraplan-go/internal/web/static/js/app.js`, `../ultraplan-go/internal/web/static/js/operations.js`, `../ultraplan-go/internal/web/static/js/sse.js`, `../ultraplan-go/internal/web/static/css/base.css`, `../ultraplan-go/internal/web/static/css/primitives.css`, `../ultraplan-go/internal/web/static/css/components.css`, `../ultraplan-go/internal/web/static/css/layouts.css`

# Frontend: read-only QA

This area decides how the TUI and local browser present Sprint 36 QA. Both interfaces show the same app-owned facts and provide the same guarded operations. Neither interface interprets verification files, decides outcomes, or treats live delivery as product state.

## Area Decisions

### Keep QA inside sprint navigation

QA is an assurance view under one sprint, not a new top-level product or TUI tab.

The browser adds these HTML routes and renders them through the existing `sprint.html` template hierarchy:

```text
/projects/{project}/sprints/{sprint}/qa
/projects/{project}/sprints/{sprint}/qa/shards/{shard-id}
/projects/{project}/sprints/{sprint}/qa/theories/{theory-id}
```

The sprint sidebar gains a `QA` link. The overview keeps a compact QA summary and links to the QA page. The QA page orders information as follows:

1. Sprint identity and QA phase status.
2. Conformance Review verdict and freshness in a separate panel.
3. QA map fingerprint, implementation fingerprint, freshness, and changed-path coverage.
4. Current blocker and app-provided next action.
5. Bounded shard summary with progress.
6. Theory outcome totals.
7. Synthesis status, contradictions, interactions, and follow-up use.
8. Attempt and durable-run references for recovery and diagnostics.

Focused shard and theory routes keep detail out of the summary response. Ordinary links make every view usable without JavaScript and give refresh, history, and assistive technology a stable location. Unknown or stale IDs return a safe error with a link to current QA status; the browser never guesses a replacement.

The TUI adds `RouteSprintQA`, `RouteSprintQAShard`, and `RouteSprintQATheory` beneath the current sprint route. The sprint menu labels the existing review item `Conformance Review` and adds `QA`. It keeps the current one-column viewport rather than adding a split-pane subsystem. The QA route renders the summary first and then selectable shards. Enter opens focused detail, Escape returns, arrows or `j` and `k` navigate, `r` refreshes authoritative state, and `c` requests cancellation only while an active durable QA run is in view. The footer continues to derive help from the existing key registry.

One column is deliberate. The current viewport already follows selection and handles small heights. A second pane would require independent focus, scrolling, width allocation, and help semantics for one feature. At wide widths the row summary can use available space; at narrow widths metadata wraps and focused detail remains a separate route.

### Render one canonical app projection

`internal/tui` and `internal/web` consume typed `internal/app` QA DTOs. Templates and TUI render functions receive bounded view models. They do not open `verification/state.json`, parse map or shard files, derive freshness, join run records, or inspect the target checkout.

The browser handler owns transport mapping and safe view labels. The TUI model owns only selection, route history, viewport position, loading state, the last bounded snapshot, and observation state. There is no global browser store and no persisted adapter state.

| State category | Owner | Frontend rule |
| --- | --- | --- |
| Map, shards, theories, synthesis, fingerprints, blockers, next action | Sprint through app DTOs | Render as supplied; never recompute. |
| Run lifecycle, liveness, cancellation, terminal result, event sequence | Run control through app DTOs | Show beside QA state without merging authorities. |
| Selected shard or theory | HTML route or TUI route | Local navigation only. |
| Expanded sections, current focus, viewport offset | Browser or TUI | Disposable presentation state. |
| Live, reconnecting, snapshot-only, history-gap indicator | Browser or TUI observer | Presentation state only; cannot alter QA status. |

### Preserve exact state axes and vocabulary

The adapters expose the API decision instead of inventing one all-purpose badge.

| Axis | Canonical values and presentation |
| --- | --- |
| QA phase | `missing`, `mapped`, `queued`, `running`, `synthesizing`, `completed`, `blocked`, `cancelled`, `interrupted`, `stale`, `invalid`. Human labels may replace underscores and capitalize words, but the meaning cannot change. |
| Freshness | `Current` or `Stale`, with the app-provided reason and affected fingerprint. Freshness remains visible even when another phase field is terminal. |
| Theory outcome | `confirmed`, `refuted`, `invalid`, `inconclusive`, `blocked`, `cross_shard`, `not_applicable`. Every value has text, not color alone. |
| Conformance Review | Existing review status and verdict, labeled `Conformance Review` for people. `review`, `review.md`, and machine values remain compatible. |
| Cancellation | The run-control cancellation state and reason. A request is shown as requested or cancelling until the canonical terminal result becomes `cancelled`. |
| Observation | `Live`, `Reconnecting`, `Snapshot only`, or `History gap`. This is never displayed as a QA outcome. |

`completed` is rendered as `Read-only QA completed`, never `QA passed`. A confirmed theory is diagnostic evidence, not an issue. The page has no issue count, repair action, promoted finding, QA verdict, or control that changes the independent Conformance Review verdict.

Unknown counts render as `Unknown`, not zero. Blocked, failed, and inconclusive remain distinct: policy or a prerequisite prevented blocked work, execution failed for failed work, and an investigator finished without resolving an inconclusive theory.

### Make server-rendered HTML the complete browser baseline

The first response includes the current QA snapshot, freshness, coverage, shard summary, theory totals, synthesis status, blocker, next action, and run links. Normal HTML forms use the existing prepare and confirmation flow for `qa-dry-run`, `qa-start`, `qa-resume`, and `qa-recover`. Post/Redirect/Get remains the mutation pattern. `qa-recover` copy states that it mutates verification state but never the target repository.

Active runs link to the existing durable run page. That page already has a normal cancellation form, current status, a durable refresh link, and copy stating that page closure does not cancel work. The QA page does not add a second cancellation endpoint or browser operation registry.

No-JavaScript users can:

- Inspect the current QA summary and focused details.
- Prepare and confirm a dry-run map, start, resume, or recovery operation.
- Follow the durable run link and refresh it for committed status.
- Request cancellation through the existing form.
- Return to the QA page after completion to load authoritative product state.

JavaScript may intercept those same forms, render the existing inline confirmation, follow committed operation events, update bounded progress, and reload the canonical page after a terminal event. It cannot add a JavaScript-only action, select a model or budget that the form cannot express, own run identity, infer success, or retain a QA workflow after navigation.

QA adds the map-owned `shard` option to the existing operation-form serializer. The browser cannot submit paths, theory content, commands, prompts, policy, fingerprints, attempt IDs, environment, or limits. `operations.js` remains a small authenticated command helper. Existing SSE and operation code remain the delivery path.

### Keep observation separate from cancellation and recovery

Closing a page, hiding a tab, losing SSE, leaving the TUI route, or quitting the TUI stops observation only. The accepted run continues. Cancellation always calls the canonical durable run cancellation use case.

Browser reconnect follows this sequence:

1. Mark the observer `Reconnecting` without changing QA state.
2. Fetch the current durable run and QA snapshots.
3. Request committed events after the last accepted sequence.
4. Deduplicate by run ID and sequence.
5. Replace displayed product facts with the app snapshot.
6. Resume live observation when the cursor remains retained.

A replay gap shows `History gap`, the retained boundary, and a link to refresh durable and QA snapshots. It does not replay browser memory as truth. A terminal event triggers a server refresh before the page claims the QA phase completed.

The TUI uses the same recovery rule. Its operation event channel can drop delivery by design. On the next refresh tick, route entry, terminal message, or explicit `r`, it reads QA status and durable events through app use cases. It never rebuilds outcomes from the local message slice. `q` leaves active work unchanged. The user must use the visible `c` action to request cancellation.

Fingerprint changes keep prior results visible as historical evidence but mark affected current work stale. The UI explains the changed input category, disables unsafe resume by omitting that action from the app-provided action set, and presents the supplied remap or recovery action. It never silently turns stale work current.

### Bound every rendered collection before presentation

The app or handler applies limits before constructing a template or TUI model. A template must not receive an unbounded collection and then hide overflow with CSS.

| Collection | Presentation bound |
| --- | ---: |
| Shards in a summary | 40 |
| Theories on a focused shard view | 24 |
| Recent QA progress entries | 100 by default, 200 maximum |
| Changed paths in one map | 512, shown as coverage totals on summary views |
| Changed paths in one primary shard detail | 64 maximum from the frozen product limit |
| Contextual paths in one shard detail | 128 maximum from the frozen product limit |
| Durable operation events requested by the TUI | 200 |
| Live operation rows in the browser QA panel | 100 |

Stable ID ordering and cursor navigation handle any list that can outgrow its focused bound. Every partial list states `Showing X of Y` and offers a normal link or action for the next page. Truncation, omitted event detail, and replay gaps use separate labels. None changes QA status or implies complete evidence.

The TUI holds one bounded app page and renders only the viewport slice. The browser initial HTML does not embed all theory bodies, all attempt history, raw command output, or all events. Focused routes fetch those records by stable ID. Runtime-free status and map views must not initialize investigator dependencies.

The existing 250 ms progress coalescing and 16 KiB durable event limit remain in force. Visual progress can update at that cadence, but the browser live region announces only phase transitions, blockers, cancellation changes, terminal results, and an aggregate progress summary no more than once every two seconds. This avoids a spoken event flood.

### Use explicit empty, pending, success, and failure states

The interfaces render these states from the same fixture:

| State | Required presentation |
| --- | --- |
| Loading | Keep the page shell and heading stable; show `Loading current QA status`. Do not show a success-colored placeholder. |
| Missing | Explain that QA has not started, show prerequisite state, and offer dry-run mapping or the app-provided start action. |
| Mapped | Show fingerprint, coverage, limits, shard totals, and the guarded start action. |
| Running or synthesizing | Show durable run identity, completed and total shards, current bounded activity, cancellation state, and a durable status link. |
| Completed | Say `Read-only QA completed`, show all outcome totals and synthesis status, and preserve the Conformance Review panel unchanged. |
| Blocked | Name the stable blocker category, affected scope, safe reason, durable state location, and next action. |
| Cancelled or interrupted | Preserve completed shard outcomes, distinguish cancelled from interrupted, and show safe resume or recovery only when supplied by app. |
| Stale or invalid | State which fingerprint or schema check failed. Do not offer ordinary resume when reuse is unsafe. |
| Empty focused view | Explain whether the shard has no theories yet, theories were not applicable, or detail is unavailable. Do not use a generic `Nothing here`. |

An error panel includes the stable code, safe summary, affected map, shard, theory, or run, current/stale/partial state, retryability, correlation ID, and app-provided next action. TUI and browser code do not render raw wrapped errors, provider payloads, unrestricted stderr, environment values, absolute paths, or full command output.

### Meet browser and terminal accessibility requirements

Browser QA uses the existing page `h1` and adds ordered `h2` sections. Shard and outcome tables have captions, column headers, row headers, and text status. Coverage uses a native `progress` element only when total and completed values are known; nearby text repeats the values. Native links, buttons, forms, `details`, and tables take priority over custom ARIA widgets.

Keyboard order follows visual order. Inline confirmation receives focus after preparation. On validation failure, focus moves to the error summary. Closing a focused detail or confirmation restores focus to the action that opened it. Routine progress never steals focus. The polite live region receives coalesced summaries; blocking request failures use an alert. Existing `:focus-visible`, skip link, minimum control size, and `prefers-reduced-motion` behavior apply to QA.

At 320 CSS pixels, panels stack, labels remain next to their values, and wide tables scroll inside a named region instead of widening the page. No control requires hover. Status always includes text and supports high contrast and monochrome display.

The TUI uses ASCII text for every QA state and does not rely on color, animation, or symbols. It sanitizes ANSI escapes and non-printing controls in theory claims, evidence summaries, errors, paths, and runtime messages before width calculation and rendering. Views remain navigable at `80x24` and `40x12`; narrow output wraps metadata and keeps the selected item and help action visible without horizontal scrolling.

Full-screen TUI startup requires interactive input and output. In a non-TTY environment, `ultraplan tui` returns a concise usage error directing automation to the ordinary QA text command or `--json`. It does not emit an alternate-screen sequence or silently switch to another renderer. CLI text and JSON remain the stable non-interactive and screen-reader fallback.

### Test semantics before layout details

One canonical QA fixture is projected through app results, CLI text, CLI JSON, TUI, browser HTML, browser JSON, and durable run detail. Tests compare map fingerprint, freshness, changed-path coverage, shard status and progress, theory outcomes, synthesis status, blocker, cancellation, terminal result, replay gap, and next action.

TUI tests cover:

- Summary, shard, and theory routes at wide and narrow terminal sizes.
- Existing key bindings, visible help, selection retention, refresh, back, and active-run cancellation.
- Monochrome output and absence of color-only meaning.
- All phase and theory outcome labels.
- Bounded rows, explicit `Showing X of Y` copy, and viewport behavior.
- Dropped operation messages followed by authoritative status refresh.
- Requested cancellation versus terminal cancellation.
- ANSI, control-character, oversized, and multiline hostile content.
- Non-TTY startup without ANSI or alternate-screen bytes.

Browser tests cover:

- Full server-rendered snapshots with JavaScript disabled for every major state.
- Summary and focused routes, stable IDs, unknown IDs, and bounded collection metadata.
- Dry-run, start, resume, recover, and cancellation through normal forms.
- Host, Origin, session, CSRF, body-limit, method, authorization, and strict input rejection.
- Escaped hostile HTML, Markdown, URLs, terminal controls, paths, and oversized evidence.
- Keyboard order, headings, landmarks, table headers, labels, focus placement and restoration, live regions, visible focus, and reduced motion.
- Mobile and desktop layouts, including scroll-contained tables at 320 CSS pixels.
- SSE ordering, duplicate sequence suppression, bounded rows, reconnect, replay gaps, session rotation, tab visibility, server restart, and terminal snapshot refresh.

Use semantic assertions for fields, actions, states, accessibility, and limits. Normalized golden fixtures may cover representative no-JavaScript HTML and TUI output, but omit timestamps, run IDs, cursor values, animation frames, and incidental spacing. Normal tests remain offline with fake app/runtime collaborators and temporary workspaces. Race tests cover delivery, refresh, cancellation, and observer replacement. Real-repository dogfood stays gated and must prove unchanged target identity.

## Trade-Offs

### Focused routes instead of one expanding sprint page

Separate QA, shard, and theory HTML routes add handler and route fixtures. They keep the sprint overview readable, make no-JavaScript detail usable, and prevent one response from embedding every theory and event. A query-driven client panel was rejected because the current router rejects unknown query parameters and because URLs should identify focused content without JavaScript state.

### One-column TUI instead of a master-detail split

A split view could show more context on large terminals. It would add another focus model, two scroll positions, width breakpoints, and narrow-terminal behavior that the current TUI does not have. Route-based focused detail reuses the current viewport and key registry and has fewer ways to drift from browser navigation.

### Server refresh instead of optimistic product state

Reloading after terminal events and recovery can feel less immediate than applying local reducers. It guarantees that the page reflects sprint-owned QA state and run-control facts after persistence, cancellation races, reconnect, or server restart. Optimistic completion was rejected because event delivery is not product authority.

### Several visible state axes instead of one status badge

Showing QA phase, freshness, Conformance Review verdict, run lifecycle, cancellation, and observer connection takes more space. Combining them would produce false claims such as a completed run meaning QA passed, a dropped connection meaning cancellation, or current QA upgrading a failed review.

### Bounded summaries instead of exhaustive rendering

Focused links and paging add navigation. Rendering everything would make large attempts slow, produce noisy live regions, and turn polling into an unbounded read. The fixed limits also make TUI and browser parity testable.

### Existing components instead of a QA design subsystem

QA reuses the current status, metadata, notice, error, progress, table, operation form, and durable run patterns. Feature-specific namespaced QA components may live with `sprint.html` during the current template migration. No new frontend framework, build pipeline, global store, generic dashboard component library, or custom widget is justified.

### Rejected alternatives

- A single-page QA application, client router, or authoritative browser store.
- A new top-level QA tab detached from sprint context.
- Direct TUI or web reads of detailed verification files.
- A browser-only operation registry, event history, or cancellation endpoint.
- Cancelling work when the TUI exits, the browser closes, or SSE disconnects.
- Showing `QA passed`, treating confirmed theories as issues, or merging QA with Conformance Review.
- Rendering all shards, theories, evidence, attempts, and events on the initial page.
- Announcing every progress event to assistive technology.
- A JavaScript-only action, modal-first workflow, or color-only state indicator.
- A second non-interactive TUI renderer instead of CLI text and JSON.
- Raw Markdown HTML, provider output, terminal control sequences, or shell command text in presentation.
- Exact whole-screen snapshots as the primary test oracle.

## Evidence

### Governed project evidence

`projects/ultraplan-go/sprints/36-read-only-qa/requirements.md` requires TUI and browser QA views, a useful no-JavaScript snapshot, escaped hostile content, reconnect and restart recovery, keyboard and focus tests, bounded rendering, and agreement across CLI, JSON, TUI, HTML, run detail, and durable progress. It also forbids issue promotion, repair, alternate state authority, and cancellation on observer loss.

`projects/ultraplan-go/docs/ARCHITECTURE.md` assigns navigation and rendering to `internal/tui`, transport and templates to `internal/web`, and shared use cases to `internal/app`. It fixes server rendering and progressive enhancement as the browser baseline and says SSE delivery cannot own durable state. `projects/ultraplan-go/docs/PRD.md` and `projects/ultraplan-go/docs/TRD.md` add the local browser security, accessibility, cross-surface recovery, responsive layout, and Phase 5 QA requirements.

The architecture and API area decisions for this sprint freeze the authority split, app DTO boundary, phase vocabulary, operation kinds, HTTP resources, cancellation path, reconnect behavior, and presentation bounds. This document applies those decisions to interaction and layout rather than creating a competing contract.

The live implementation supports a conservative extension. `../ultraplan-go/internal/tui/model.go:27-80` models route and disposable view state with app DTOs, `../ultraplan-go/internal/tui/keys.go:24-59` defines one key registry, and `../ultraplan-go/internal/tui/views.go:610-653` already renders a height-bounded viewport. `../ultraplan-go/internal/tui/app.go:232-310` durably accepts operations and allows local event delivery to drop, while `:357-388` recovers run snapshots and bounded events through app use cases.

For the browser, `../ultraplan-go/internal/web/templates/sprint.html:24-45` establishes summary-first sprint presentation and normal operation forms, and `../ultraplan-go/internal/web/templates/operation.html:11-24` supplies a no-JavaScript cancellation and refresh path. `../ultraplan-go/internal/web/operation_handlers.go:292-408` implements prepare, confirm, redirect, status, and cancellation for ordinary forms. `../ultraplan-go/internal/web/static/app.js:718-732` bounds operation rows, while `:876-909` reconnects SSE and reloads after terminal delivery. `:1375-1476` shows the existing durable cursor, replay-gap, tab-visibility, and reconnect pattern. QA should reuse these paths, not fork them.

### Report evidence and project inference

The selected reports compare other Go tools. They support the constraints below, but the specific QA routes, labels, limits, and state composition are UltraPlan decisions.

- `studies/go-cli-study/reports/final/05-error-handling.md:32-38`, `:63-89`, and `:125-176` support typed, cause-preserving errors and top-level user guidance. This supports stable QA error panels instead of matching or exposing wrapped strings.
- `studies/go-cli-study/reports/final/06-io-abstraction.md:63-85` and `:155-170` support injected terminal and process streams for deterministic renderer tests. They do not justify a broad frontend filesystem abstraction.
- `studies/go-cli-study/reports/final/07-state-context.md:119-165` and `:219-238` distinguish cancellation plumbing from persistent state, locks, and bounded cleanup. This supports treating route or connection loss as observer loss only.
- `studies/go-cli-study/reports/final/08-concurrency.md:71-112` and `:357-377` support bounded admission and shutdown. The frontend inference is to bound retained rows and avoid making a slow observer part of execution backpressure.
- `studies/go-cli-study/reports/final/09-terminal-ux.md:96-102`, `:144-156`, and `:194-209` support single-owner progress, interruptible interaction, and deliberate non-TTY behavior. The report notes little explicit accessibility evidence at `:227-239`, so UltraPlan's browser and terminal accessibility rules come from project requirements rather than imitation.
- `studies/go-cli-study/reports/final/10-logging-observability.md:188-213` and `:277-317` support structured fields and separation of user output from diagnostics. They do not prove redaction or cross-surface parity; UltraPlan requires both before public rendering.
- `studies/go-cli-study/reports/final/11-testing-strategy.md:63-114`, `:204-212`, and `:273-305` support scripted scenarios, normalized goldens, recording fakes, fault injection, and selected real boundaries. This supports semantic parity fixtures plus gated dogfood.
- `studies/go-cli-study/reports/final/13-security.md:89-103`, `:113-167`, and `:259-269` support explicit arguments, permission checks, canonical paths, and redacting values. The frontend applies that evidence by treating all theory, evidence, path, and runtime text as hostile data.
- `studies/go-cli-study/reports/final/14-performance.md:73-111`, `:150-166`, and `:235-262` support streaming, bounded queues, hard caps, and disk-backed long-lived state. The presentation inference is to page focused data and cap DOM and terminal history rather than rendering the complete durable record.

The key conclusion is simple: the user can leave and return without changing the work. The app and durable run records provide current truth; TUI and browser state only decide what part of that truth is visible.

## Risks

- The app DTO must freeze one mapping from canonical values to human labels. If TUI and browser title-case, group, or explain states independently, parity will drift even when raw JSON agrees.
- `stale` appears in the phase vocabulary while freshness is also a separate field. View models must show both without inventing a previous phase or hiding which fingerprint changed.
- The current TUI renders some raw `error.Error()` values. QA paths must use bounded safe app errors and sanitize control characters before rendering, or provider and path details can bypass the public error contract.
- The current TUI runner does not visibly enforce interactive input and output in `internal/tui/app.go`. Non-TTY tests must drive the final startup decision; otherwise alternate-screen bytes may leak into redirected output.
- The operation form serializer and Go transport DTO do not yet carry `shard`. Adding it in only HTML, JavaScript, or Go would break focused-run parity or silently start the full map.
- The browser has both transient operation rows capped at 100 and durable run rows capped at 500. QA pages must use the 100-entry QA summary limit and label links to the separate durable history so users do not confuse the two retention views.
- A replay gap can coincide with a valid current QA snapshot. The UI must show current product state and incomplete observed history together rather than marking the whole QA attempt invalid.
- Partial fingerprint invalidation may leave current and stale shards in one historical attempt. Every row needs explicit freshness, and synthesis must state whether it covers the current map.
- A confirmed theory may look like a failure badge. Outcome copy and neutral phase styling must keep `confirmed` distinct from transport failure, issue promotion, and Conformance Review verdict.
- Live-region updates can still overwhelm users if every 250 ms visual coalescing tick is announced. The two-second aggregate announcement limit and transition-only messages need explicit tests.
- Responsive tables can hide relationships if the first columns scroll out of view. Focused mobile fixtures must prove shard identity, status, blocker, and action remain reachable at 320 CSS pixels.
- Existing templates are midway through a namespaced hierarchy migration. QA components must use stable names and downward-only dependencies even if their definitions temporarily remain in `sprint.html`.
- Browser read visibility and mutation authorization differ after session rotation. A visible stale page must not leave start, resume, recovery, or cancellation controls that bypass current session and CSRF checks.
- Fixed presentation limits may hide a legitimate large change unless totals and paging remain accurate. Hitting a display bound is not a QA blocker and must never be reported as complete collection coverage.
- Terminal accessibility remains constrained by terminal and screen-reader support. Text labels, no color-only meaning, predictable keys, and CLI text/JSON fallback are required, but they do not guarantee every full-screen terminal combination works.
- Atomic publication and summary reconciliation remain product concerns. Frontends must expose `invalid`, `stale`, persistence degradation, or cleanup uncertainty and the supplied recovery action rather than attempting repair during a read.
