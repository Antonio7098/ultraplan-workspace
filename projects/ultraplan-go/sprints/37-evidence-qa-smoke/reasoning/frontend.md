> **Inputs Used:** `projects/ultraplan-go/sprints/37-evidence-qa-smoke/requirements.md`, `projects/ultraplan-go/sprints/37-evidence-qa-smoke/code-context.md`, `projects/ultraplan-go/sprints/37-evidence-qa-smoke/sprint-index.md`, `projects/ultraplan-go/sprints/37-evidence-qa-smoke/technical-handbook.md`, `projects/ultraplan-go/docs/ARCHITECTURE.md`, `projects/ultraplan-go/docs/PRD.md`, `projects/ultraplan-go/docs/TRD.md`, `studies/go-cli-study/reports/final/01-project-structure.md`, `studies/go-cli-study/reports/final/02-command-architecture.md`, `studies/go-cli-study/reports/final/05-error-handling.md`, `studies/go-cli-study/reports/final/06-io-abstraction.md`, `studies/go-cli-study/reports/final/07-state-context.md`, `studies/go-cli-study/reports/final/08-concurrency.md`, `studies/go-cli-study/reports/final/09-terminal-ux.md`, `studies/go-cli-study/reports/final/10-logging-observability.md`, `studies/go-cli-study/reports/final/11-testing-strategy.md`, `studies/go-cli-study/reports/final/13-security.md`, `studies/go-cli-study/reports/final/14-performance.md`

# Frontend reasoning: evidence-producing QA and smoke integration

This area covers the TUI and loopback browser presentation of QA. It decides how operators inspect bounded evidence, adjudication, promoted issues, canonical assessment, smoke compatibility, blockers, cancellation, and recovery. It does not move evidence, verdict, run, or mutation authority into either interface.

## Area Decisions

### One semantic projection, two renderers

`internal/app` will expose one typed, bounded QA projection for CLI, TUI, and browser consumers. `internal/tui` and `internal/web` will map that projection into terminal state and browser view models without reading `verification/`, `qa.md`, `smoke.md`, `flow-state.json`, or external harness files directly. They will not import `internal/sprint`, call CLI handlers, parse terminal output, classify evidence, promote issues, or derive the canonical assessment.

The projection must carry the values on which the interfaces have to agree:

- project, sprint, current semantic attempt, durable run correlation, and update time;
- QA phase and operational lifecycle as separate fields;
- freshness plus governed-input, implementation, review, policy, and evidence identities in bounded form;
- map and shard coverage, evidence counts and statuses, rejected-evidence counts, and truncation facts;
- adjudication status, promoted issue summaries, regression-candidate classification, and exact contained references;
- canonical QA assessment, Conformance Review input, smoke-suite status, blockers, cleanup result, cancellation, recovery, and next action;
- paths and digests for `qa.md`, `smoke.md`, verification records, generated patches, and validated external smoke evidence, subject to the app's path and disclosure policy.

The app projection, not a shared cross-renderer widget abstraction, is the reuse boundary. Browser and TUI composition differ enough that forcing them through one presentation model would either leak HTML concerns into terminal code or reduce both views to the least expressive format.

### Information order and status language

Both interfaces will present the same hierarchy:

1. Identity and freshness.
2. Canonical QA assessment and its current blockers or next action.
3. Current operation, cancellation, cleanup, and recovery state.
4. Coverage and shard progress.
5. Evidence and adjudication.
6. Promoted issues and regression candidates.
7. Smoke-suite compatibility and containing-suite evidence.
8. Canonical artifact and contained evidence references.

The phase, operation result, evidence outcome, adjudication decision, issue severity, review verdict, smoke verdict, and canonical assessment are distinct labels. A completed investigator run is never rendered as a QA pass. A failing command is never rendered as a promoted issue. A narrow or diagnostic smoke pass is never rendered as canonical containing-suite success.

Every status uses text in addition to color. Stable labels include `current`, `stale`, `running`, `completed`, `blocked`, `cancelled`, `timed out`, `cleanup uncertain`, `rejected evidence`, `promoted issue`, `diagnostic only`, and `containing suite required`. Renderers use the exact typed value for machine-significant terminology and may add a short plain-language explanation. They must not replace a precise blocked or cleanup-uncertain result with a generic failure banner.

When a new attempt fails before canonical publication, the view shows two clearly separated facts: the current attempt failed, and the last complete `qa.md` remains available but is not current evidence for the failed attempt. Freshness appears beside the assessment rather than only in a detail pane.

### Bounded evidence and adjudication detail

The default evidence list shows one row or card per bounded app item. It includes evidence ID, shard, evidence kind, check or plan identity, outcome, current/stale state, repeatability or deterministic-sufficiency result, containment, cleanup, truncation, and adjudication disposition. Detail views may add expectation references, confirmation/refutation conditions, command summary, workspace and target identity digests, generated-patch path and digest, rejection reason, and exact issue links.

Raw provider payloads, unrestricted stdout/stderr, full environment values, arbitrary workspace files, and external harness payloads do not enter these views. A generated patch is represented by bounded metadata and a guarded artifact preview only if the app already authorizes that path and content type. The frontend does not create a general file browser.

Rejected evidence remains visible in a separate subsection because it explains why a failure did not become an issue. Promoted issues appear only from the app's adjudicated issue records. Each issue shows severity, root-cause group, evidence references, promotion reason, repair eligibility, regression-candidate state, and next action. There is no assignment, editing, closing, scheduling, or remote synchronization UI.

### Smoke presentation without authority collapse

Smoke appears within QA as a named suite while retaining its independent compatibility facts. The view shows the invoked entry point, review gate, selected scope, diagnostic or canonical status, required containing suite, enumerated coverage, run ID, counts, verdict, open issue references, cleanup result, `smoke.md`, and validated external evidence location.

`smoke` and `qa --suite smoke` may have different invocation labels but must render the same authoritative execution and evidence facts when they select the same work. The UI does not present two smoke implementations or imply that QA copied raw harness evidence into product state.

### TUI composition and interaction

`internal/tui/qa_view.go` remains a bounded renderer over app DTOs. Its model owns only view state: active section, focused item, scroll position, selected bounded ID, viewport size, and whether a guarded action prompt is open. It does not cache a second QA result or mutate product status in response to progress messages.

The wide layout uses a compact summary header and a master-detail body. The left pane lists sections and bounded records; the right pane shows the selected detail and actions. At narrow widths it becomes a single-column stack. Long paths, IDs, and hostile text wrap or receive an explicit visual truncation marker; the selected detail retains the full bounded value. No semantic status disappears solely because the terminal is narrow.

Keyboard behavior follows the existing TUI navigation model:

- `Tab` and `Shift+Tab` move between section navigation, records, details, and actions;
- arrow keys or the existing list-navigation aliases move within the focused region;
- `Enter` opens the selected detail or the existing guarded prompt for an action;
- `Esc` closes detail or confirmation state before leaving the QA view;
- the displayed action bar exposes refresh, cancel, resume, and recovery only when the app marks each action available.

Focus is always visible without relying on color. The footer shows the active keys. Cancellation requires confirmation, displays the affected durable run and scope, and changes only after the app accepts the request. Completed evidence remains visible while cancellation and bounded cleanup proceed. Exiting or resizing the TUI stops observation only; it does not infer cancellation or success.

Terminal rendering strips or visibly escapes control bytes and embedded ANSI sequences from all evidence-derived strings before width calculation. Newlines in record titles become spaces; multiline content is confined to detail blocks. This is a renderer defense after app-level bounding and redaction, not a substitute for those controls.

### Browser composition and no-JavaScript behavior

The sprint page keeps its existing server-rendered template hierarchy. The current `component/run-qa` in `internal/web/templates/run_qa.html` should be extended or composed into `internal/web/templates/sprint.html`; Sprint 37 does not add a single-page application, frontend framework, client router, client store, or asset build step.

The complete current QA result, evidence summaries, adjudication, issues, smoke status, blockers, artifact references, and next actions render in the initial HTML response. Read navigation uses ordinary links and anchors. Guarded start, cancel, resume, and recovery paths remain usable through ordinary forms and server-rendered confirmation/result pages. The forms retain the existing same-origin, CSRF, session, normalized-request, governed-fingerprint, authorization, and one-time confirmation checks. JavaScript is never required to discover a blocker, inspect current evidence, cancel an authorized run, or recover after interruption.

`internal/web/static/js/operations.js` may enhance the baseline with bounded SSE progress, inline cancellation, and targeted refresh. It treats events as hints. On terminal events, replay gaps, disconnect, session rotation, or observer restart it fetches or reloads the authoritative app query. It never computes an assessment, promotes an issue, changes freshness, or treats subscriber loss as operation failure.

Templates receive explicit typed view models. They do not read files, call use cases, validate requests, or interpret product state. Evidence strings use `html/template` escaping. No evidence-derived value enters `template.HTML`, an inline script, a style attribute, or an unvalidated URL. Bounded Markdown and JSON previews use the existing safe preview path with embedded HTML and active content disabled.

Browser layout uses semantic landmarks, a single page heading, ordered section headings, labelled status groups, table captions, row headers where tables are used, and descriptive link and button names. Status never relies on color or icon shape alone. Dynamic progress uses a polite live region for coarse state changes; high-frequency event text does not flood assistive technology. After an enhanced action, focus moves to the result heading or error summary. Validation errors link back to the affected control. Reduced-motion preferences disable nonessential animation, and small screens stack summary fields and cards without requiring page-level horizontal scrolling.

### Progress, cancellation, and recovery

Both renderers distinguish durable state from transient delivery. Progress may be sampled, coalesced, truncated, or absent. The UI shows replay gaps and sampled history explicitly, then offers an authoritative refresh. It does not reconstruct missing state from the events that remain.

Cancellation has three visible steps when supplied by the app: requested, stopping work, and cleanup result. `cancelled` and `cleanup uncertain` are separate terminal facts. A cancellation request stops new scheduling but does not erase completed valid evidence. Browser disconnect, page navigation, TUI exit, and SSE subscriber loss do not request cancellation.

Recovery actions are driven by typed availability and next-action fields. A stale input prompts rerun from current inputs; interrupted ownership prompts the existing recovery use case; invalid state presents the exact safe recovery guidance; cleanup uncertainty remains blocking. The frontend does not turn a generic retry button into an unconditional rerun.

### Performance bounds

App DTO limits are the primary rendering bound. Each view receives capped collections, safe summaries, explicit total and omitted counts, and bounded preview text. The browser does not serialize all detailed verification records into HTML or JavaScript. The TUI renders only the current viewport and selected bounded detail. Neither interface retains an unbounded event list.

The browser uses the initial HTML as the baseline and one bounded SSE subscription only while observing active work. The TUI coalesces progress updates before redraw. Both refresh canonical state after a terminal event instead of replaying every event into a local model. No buffer pool, virtualized browser framework, or frontend cache is added without measurement; the existing bounded projections and viewport rendering are sufficient for this sprint.

### Verification strategy

One adapter-independent fixture will drive app, CLI/JSON, TUI, browser HTML/JSON, durable-run, `qa.md`, and verification-state parity checks. Parity assertions compare semantic fields, not identical prose or layout.

TUI tests use deterministic messages and injected app use cases. They cover every phase and assessment combination, wide and narrow widths, focus order, keyboard-only actions, cancellation confirmation, completed-evidence preservation, dropped progress, replay gaps, recovery availability, hostile ANSI/control text, long IDs and paths, explicit truncation, and absence of an invented pass verdict.

Browser tests use `httptest`, fake app use cases, deterministic templates, and fake run/process dependencies. They cover route and method handling, strict JSON, no-JavaScript full-page snapshots, ordinary-form actions, confirmation expiry and staleness, same-origin and CSRF rejection, authorization, hostile HTML/Markdown/URL values, redaction, bounded output, reconnect, replay gaps, session rotation, observer restart, cancellation, and last-complete versus current-failure presentation.

Focused template tests assert heading order, landmarks, labels, table semantics, status text independent of color, focus targets, and no inline or third-party executable content. CSS/browser tests cover reduced motion, keyboard focus, mobile stacking, wrapping, and zoom. A gated real-browser smoke check covers navigation, guarded confirmation, live progress, cancellation, refresh/reconnect recovery, and state agreement; normal tests remain offline.

Golden or snapshot tests protect stable no-JavaScript HTML and TUI rendering, but focused assertions protect security and authority rules. Golden updates require review of dynamic-field normalization and every semantic change. They cannot replace tests that assert hostile content stays inert or that a diagnostic smoke pass is not canonical.

## Trade-Offs

### Shared app DTO versus shared renderer

The decision is to share typed semantics and keep renderer composition separate. This adds mapping code in `internal/tui` and `internal/web`, but preserves their different accessibility, layout, and interaction needs. A common UI interface or generic component tree was rejected because the selected project-structure evidence warns that broad UI interfaces can become god interfaces, and the project already defines app use cases as the stable cross-surface boundary.

### Complete server-rendered baseline versus JavaScript-first operation

The decision is to render a complete HTML baseline and use JavaScript only for progress and convenience. A client-side application would make live updates easier, but it would add a second router and state model, weaken no-JavaScript recovery, and tempt the browser to infer authority from events. That cost is not earned by the bounded QA view.

### Bounded summaries versus raw evidence browsing

The decision is to show bounded summaries and guarded previews. Rendering raw output would help ad hoc debugging, but it would expand the secret, hostile-content, memory, and path-containment boundary and would duplicate the external harness and verification stores. Operators still receive exact IDs, digests, contained paths, rejection reasons, and links needed to inspect authoritative records through approved paths.

### Authoritative refresh versus optimistic local state

The decision is to refresh app state after actions, terminal events, and delivery gaps. Optimistically marking a run cancelled or an assessment passed would feel faster but can disagree with fencing, cleanup, publication, and terminal arbitration. A brief pending state is preferable to a false terminal claim.

### Dense dashboard versus staged disclosure

The decision is to put assessment, freshness, blockers, and next action first, then disclose bounded evidence detail. A single dense table would expose more fields at once but fails on narrow terminals, mobile screens, and assistive navigation. Detail views add interaction, yet they preserve a readable default without hiding authority-bearing facts.

### Rich motion and color versus durable semantics

The decision is to use restrained progress feedback and textual status. Spinners and color may aid scanning, but they are transient and inaccessible as the sole signal. The no-JavaScript page, non-color terminal rendering, reduced-motion mode, and durable artifact paths remain complete.

### One smoke presentation versus a second QA-specific smoke view

The decision is to render smoke as a QA suite with explicit compatibility metadata. A separate QA smoke dashboard was rejected because it would obscure whether `smoke` and `qa --suite smoke` reuse the same discovery, containing-suite, invocation, evidence, and verdict path.

### Goldens plus focused assertions versus either alone

The decision is to use both. Goldens catch broad presentation drift, while focused tests keep authority, escaping, cancellation, and accessibility rules visible. Goldens alone can bless a security defect; substring assertions alone miss accidental omissions and layout regressions.

## Evidence

### Report findings

- `studies/go-cli-study/reports/final/01-project-structure.md`, "Thin CLI Entry Point", "Unidirectional Dependency Flow", and "UI Interface Abstraction", supports thin interface adapters and reusable behavior behind them. Its warning that a UI interface can become a god interface supports sharing the app projection rather than a cross-renderer widget API.
- `studies/go-cli-study/reports/final/02-command-architecture.md`, "Bifurcated CLI/TUI Pattern" and "Bifurcated CLI/TUI hiding real commands", identifies parity loss when meaningful behavior exists only in the TUI. This supports app-owned QA actions and identical smoke execution facts across adapters.
- `studies/go-cli-study/reports/final/05-error-handling.md`, "User/Operational Separation", "Multi-Error Aggregation", and "Exit Code Mapping from Error Types", supports distinct typed outcomes and actionable user messages without collapsing partial failures or diagnostics into one generic error.
- `studies/go-cli-study/reports/final/06-io-abstraction.md`, "MockTerminal / MockUi for UI Testing" and "IOStreams with Test Constructor", supports injected interface dependencies and deterministic renderer tests without a real terminal or server.
- `studies/go-cli-study/reports/final/07-state-context.md`, "Signal-Context Wiring" and practical tip 6, shows that cancellation must propagate and that cleanup may need a separate bounded context. The project-specific conclusion is to render cancellation and cleanup as separate facts.
- `studies/go-cli-study/reports/final/08-concurrency.md`, "Deferred Cancel + Explicit Wait with Timeout" and "Cleanup Models", supports bounded shutdown and explicit lifecycle handling. It does not justify concurrent writable QA, nor does it make transient UI delivery authoritative.
- `studies/go-cli-study/reports/final/09-terminal-ux.md`, "Progressive Streaming with Interruptibility", "Non-TTY Fallback", "Interruptibility Depth", and "No accessibility features", supports bounded progress, explicit cancellation outcomes, narrow-terminal fallback, and deliberate keyboard and non-color accessibility work rather than assuming the TUI library provides it.
- `studies/go-cli-study/reports/final/10-logging-observability.md`, "Output Interface Abstraction" and "Output Separation Strategies", supports separating user presentation, operational diagnostics, and machine output. The report supplies no evidence-authority rules, so logs and progress remain diagnostic inputs only.
- `studies/go-cli-study/reports/final/11-testing-strategy.md`, "Golden File Regression Prevention", "Fake/Stub Command Runner", and "CLI Integration Depth", supports layered renderer, handler, parity, and gated browser tests. Its warning about implementation-detail assertions supports semantic parity checks rather than widget-count tests.
- `studies/go-cli-study/reports/final/13-security.md`, "Secret Redaction Type", "Credential Scrubbing in Logs", "Trust Boundary Visibility", and its path-validation cautions, supports explicit redaction, containment, default-deny previews, and inert rendering of untrusted evidence. The report does not prove a complete browser escaping or symlink-safe artifact design; those controls come from this project's requirements.
- `studies/go-cli-study/reports/final/14-performance.md`, "Streaming via Channels and bufio", "Incremental Operations with Bounded Data Structures", and "Unbounded in-memory accumulation", supports bounded progress and viewport rendering. Its warning against speculative pooling supports not adding frontend resource machinery before measurement.

### Project-specific conclusions

- `projects/ultraplan-go/docs/ARCHITECTURE.md` makes `internal/tui` and `internal/web` adapters over `internal/app`, defines the browser's page-to-primitive template order, and states that browser views and SSE are projections rather than authority. The one-projection decision applies those rules to Sprint 37 QA.
- `projects/ultraplan-go/docs/PRD.md` requires one product core, a run that outlives its observer, server-rendered local pages, guarded operations, and durable local artifacts. This rules out browser-session authority and a JavaScript-only QA workflow.
- `projects/ultraplan-go/docs/TRD.md` requires narrow-terminal behavior, explicit typed view models, no-JavaScript HTML, same-origin and CSRF protections, hostile-content handling, reduced-motion and browser tests, and CLI/TUI/web agreement. The composition and test decisions above make those requirements concrete for QA.
- `projects/ultraplan-go/sprints/37-evidence-qa-smoke/requirements.md` requires all surfaces to agree on attempt, fingerprints, evidence, adjudication, issues, assessment, smoke suite, blockers, cancellation, recovery, and next action. It also requires current failure to remain visible while preserving the last complete report and forbids adapters from becoming evidence authorities.
- Prepared source evidence in `projects/ultraplan-go/sprints/37-evidence-qa-smoke/code-context.md` identifies `internal/app/sprint_usecases.go` as the adapter-independent contract, `internal/web/qa_handlers.go` as an app-backed no-JavaScript query adapter, `internal/web/templates/run_qa.html` as the existing QA component lacking Sprint 37 adjudication/issues/assessment/smoke fields, and `internal/tui/qa_view.go` as a bounded verdict-neutral renderer. Extending these paths is smaller and safer than adding a parallel frontend workflow.
- No selected study report defines canonical `qa.md`, product-owned adjudication, smoke-as-QA parity, hostile browser rendering, or a cross-surface QA contract. Those decisions are project-specific inferences from the governed requirements and project documents, not claims imported from the reports.

## Risks

- **Projection drift:** Browser and TUI mapping can still omit different fields even when both consume the same DTO. The shared fixture must compare authority-bearing semantics across all surfaces, and every additive app field must have explicit renderer coverage.
- **Over-broad projection:** Adding raw evidence for convenience could leak secrets, control bytes, arbitrary paths, or oversized output. App DTOs must remain bounded and redacted; renderers add escaping but cannot repair an unsafe contract.
- **Last-complete confusion:** Preserved `qa.md` can look current after a new attempt fails. The current-attempt failure, freshness, artifact identity, and last-complete label must remain adjacent in both interfaces.
- **Cancellation race:** Work may complete, fail, or become cleanup-uncertain while a cancellation request is in flight. The UI must show the app's accepted request and eventual terminal result rather than assigning a local outcome.
- **Delivery gaps:** SSE or TUI progress can be sampled, dropped, or replay-gapped. Both renderers need an explicit incomplete-history message and authoritative refresh path; silent reconstruction is prohibited.
- **Hostile terminal content:** HTML escaping does not protect terminal output, and width calculations can be corrupted by control sequences. TUI sanitization must precede wrapping and truncation, with dedicated adversarial fixtures.
- **Unsafe artifact links:** A valid-looking path can still escape an allowed root or point to a stale external record. Only app-authorized contained references may become links or previews; otherwise the renderer shows inert bounded text and recovery guidance.
- **No-JavaScript action decay:** Short-lived confirmations may expire between prepare and submit. Server-rendered error pages must preserve scope, explain staleness, and offer a fresh preparation path without replaying the mutation.
- **Accessibility regressions:** Dense evidence tables, frequent progress updates, focus resets, and color-only severity are likely failure points. Semantic template tests and gated keyboard/browser checks are required, and CLI text remains the non-TUI terminal fallback.
- **Small-screen overflow:** IDs, digests, paths, and evidence conditions are long. Browser cards must wrap safely, tables need contained responsive behavior, and the TUI must stack rather than horizontally clip authority-bearing text.
- **Snapshot complacency:** Golden updates can normalize away stale/current mistakes or hostile-content regressions. Focused assertions for authority labels, escaping, action availability, and canonical-versus-diagnostic smoke status remain mandatory.
- **DTO completeness:** The current app projection predates writable evidence, adjudication, issues, canonical assessment, and smoke-suite status. If it lacks a required fact, implementation must extend `internal/app`; neither renderer may infer the value from sibling fields or read private state directly.

No frontend authority decision remains open. Implementation still has to verify exact additive app field names and existing key-binding conventions, but those details must satisfy the behavior above rather than reopen the ownership or no-JavaScript decisions.
