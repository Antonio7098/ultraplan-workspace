# Frontend reasoning: evidence-producing QA and smoke integration

> **Inputs Used:** `projects/ultraplan-go/sprints/37-evidence-qa-smoke/requirements.md`, `projects/ultraplan-go/sprints/37-evidence-qa-smoke/sprint-index.md`, `projects/ultraplan-go/sprints/37-evidence-qa-smoke/technical-handbook.md`, `projects/ultraplan-go/docs/ARCHITECTURE.md`, `projects/ultraplan-go/docs/PRD.md`, `projects/ultraplan-go/docs/TRD.md`, `system/reasoning/frontend-reasoning-template.md`, `studies/go-cli-study/reports/final/01-project-structure.md`, `studies/go-cli-study/reports/final/02-command-architecture.md`, `studies/go-cli-study/reports/final/03-dependency-injection.md`, `studies/go-cli-study/reports/final/05-error-handling.md`, `studies/go-cli-study/reports/final/06-io-abstraction.md`, `studies/go-cli-study/reports/final/07-state-context.md`, `studies/go-cli-study/reports/final/08-concurrency.md`, `studies/go-cli-study/reports/final/09-terminal-ux.md`, `studies/go-cli-study/reports/final/10-logging-observability.md`, `studies/go-cli-study/reports/final/11-testing-strategy.md`, `studies/go-cli-study/reports/final/13-security.md`, `studies/go-cli-study/reports/final/14-performance.md`

> **Selected Area Template:** Frontend, from `system/reasoning/frontend-reasoning-template.md`.

This area covers the TUI and local browser presentation of QA evidence, adjudication, issues, canonical assessment, smoke compatibility, cancellation, and recovery. It does not let either interface derive evidence validity, issue promotion, freshness, or assessment.

## Area Decisions

### Key conclusion

Proceed with a server-rendered QA section on the existing sprint page and a feature-owned TUI QA view. Both consume bounded typed results from `internal/app`. The browser remains useful without JavaScript. Dependency-free JavaScript adds live run updates, reconnect, pagination, disclosure controls, and cancellation, but it never becomes the router, record store, or assessment authority.

The UI must make two distinctions impossible to miss:

- Current attempt status is separate from the canonical QA assessment.
- A failed or cancelled rerun may coexist with an older complete `qa.md`.

Every view labels `current attempt`, `canonical attempt`, freshness, and durable run status separately. It must not show an old passing report as the result of a newer failed attempt.

### Shared presentation contract

The app layer supplies presentation-ready summary and detail records. Templates and TUI code may format them, but cannot inspect verification files or recompute product decisions.

| Presentation fact | Required display |
| --- | --- |
| Assessment | Text label for `incomplete`, `blocked`, `fail`, `pass_with_findings`, or `pass` |
| Attempt | Current attempt ID, status, start or completion time, and durable run ID |
| Canonical report | Canonical attempt ID, `qa.md` reference, and explicit retained-report label when it differs from current |
| Freshness | `current` or `stale`, with the product-provided reason and next action |
| Coverage | Map, shard, theory, generated-check, and required-suite counts |
| Evidence | Bounded accepted, rejected, invalid, inconclusive, flaky, and truncated counts |
| Adjudication | Root-cause groups, rejection reasons, requested follow-up, and sufficiency statement |
| Issues | Severity, repair eligibility, regression-candidate status, and exact evidence references |
| Smoke | Suite, containing-suite identity, canonical or diagnostic status, external run identity, verdict, and `smoke.md` reference |
| Cleanup | Workspace removal and descendant-cleanup certainty, never inferred from cancellation or process exit |
| Recovery | Product-provided blockers and one next action |

Status labels use text and shape as well as color. `blocked`, `cleanup_uncertain`, `stale`, and `fail` remain visually distinct and do not collapse into a generic error badge.

### Browser composition

`internal/web/templates/sprint.html` remains the route-level page. Add one QA section composed from existing primitives and components rather than a separate client application or dashboard route.

The section order is:

1. Canonical assessment header with freshness, current attempt, canonical attempt, and next action.
2. Blocker and recovery panel, shown before evidence when action is required.
3. Guarded action bar for start, focused shard, resume, restart, `qa --suite smoke`, and cancellation.
4. Coverage and evidence-quality summary.
5. Attempt selector and bounded evidence list.
6. Adjudication and rejected-evidence list.
7. Promoted issues and regression candidates.
8. Smoke compatibility panel with canonical-versus-narrow status.
9. Artifact links for `qa.md`, `smoke.md`, and governed bounded patch previews.

No-attempt, no-evidence, no-issues, and not-applicable smoke are separate empty states. "No promoted issues" is not shown when adjudication is incomplete or blocked.

The no-JavaScript path renders the current summary, paged evidence and issues, guarded preparation, a server-rendered confirmation page, start, cancellation, and recovery links. Forms use the existing same-origin, CSRF, confirmation, fingerprint, and authorization checks. The browser never submits argv, writable paths, environment values, evidence decisions, or assessment fields.

`internal/web/qa_handlers.go` maps HTTP requests to app use cases and prepares explicit view models. It does not import `internal/sprint`, read verification JSON, follow external evidence paths directly, or call the smoke harness. Templates receive already bounded records and contained artifact references.

`internal/web/static/js/operations.js` enhances the same controls. It may subscribe to durable run events, update progress text, request cancellation, and refresh the authoritative summary after terminal events or replay gaps. An SSE disconnect changes only the connection indicator. It never marks a run cancelled, failed, or complete.

### TUI composition

`internal/tui/qa_view.go` is a feature-owned view over the same summary and attempt-detail use cases. It uses the existing TUI model, operation dispatch, durable-run progress, and cancellation commands. It does not add a QA-specific operation registry or read verification files.

The TUI has three stable regions when width permits: assessment and blockers, evidence or issues list, and selected-record detail. Narrow terminals stack those regions and keep the assessment and next action first. Detail text wraps to the viewport and never forces horizontal scrolling for ordinary evidence.

Keyboard behavior follows the existing navigation conventions and adds only feature-local actions:

| Key | Action |
| --- | --- |
| `Tab` and `Shift+Tab` | Move among QA regions and action controls |
| Arrow keys or `j` and `k` | Move within the focused bounded list |
| `Enter` | Open or close selected evidence, issue, or artifact detail |
| `PageUp` and `PageDown` | Page through bounded evidence and issues |
| `r` | Refresh authoritative QA and run state |
| `c` | Open cancellation confirmation for an active authorized run |
| `Esc` | Close detail or confirmation and restore prior focus |

Starting, resuming, restarting, and smoke execution use the existing guarded operation confirmation. Cancellation requires a second confirmation and reports only that the request was accepted, already requested, or terminal. Focus returns to the initiating control after cancellation or a failed request.

### State and data

| State category | Owner and use |
| --- | --- |
| Product state | QA summary, attempt details, evidence, adjudication, issues, assessment, and next action from app use cases |
| Durable run state | Lifecycle, progress cursor, cancellation, replay gap, and terminal result from run-control use cases |
| URL state | Browser attempt ID, selected panel, cursor, and bounded page size so refresh is recoverable |
| Local browser state | Open disclosures, current focus, pending form submission, and SSE connection status only |
| Local TUI state | Focused region, selected row, current page, confirmation visibility, and viewport size only |
| Derived presentation | Labels, count text, and layout choice from already classified app fields |

There is no global client store. Browser refresh, TUI restart, session rotation, and observer restart rebuild the view from app and durable-run reads. Dropped delivery triggers a full summary refresh rather than replaying guesses from local state.

### Loading, empty, error, and success states

The initial browser response always contains complete server-rendered content. JavaScript adds a small `Connecting to run updates` status without hiding the existing page. The TUI shows the last authoritative snapshot while a refresh is pending.

Typed states render as follows:

| State | User presentation |
| --- | --- |
| Missing | Explain that QA has not run and show the current admission action or blocker |
| Running | Show durable run ID, current phase, bounded progress, elapsed time, and cancellation when authorized |
| Incomplete | Show preserved evidence counts and the exact missing requirement or next action |
| Blocked | Put the blocker before all success-colored content and show recovery guidance |
| Failed | Show current failure and whether an older canonical report remains available |
| Cancelled | Show cancellation, preserved completed evidence, and cleanup status separately |
| Cleanup uncertain | Use a dedicated high-prominence warning and prohibit passing language |
| Stale | Keep old details inspectable, mark them stale on every panel, and show the current fingerprint action |
| Pass with findings | Show that required evidence passed and list the remaining promoted findings |
| Pass | Show complete current coverage and canonical containing-suite status |
| Replay gap | State that live history is incomplete, then refresh from the durable snapshot and QA summary |

Retry controls appear only when the app result says the action is allowed. The UI does not convert a generic error into a resume or restart recommendation.

### Hostile content and artifact previews

Evidence, issue titles, command excerpts, model-derived summaries, paths, and external identifiers are untrusted text.

Browser templates rely on `html/template` escaping and render evidence as plain text. They do not render investigator or adjudicator Markdown as HTML. External evidence is displayed only through app-validated references. Generated patches use the existing bounded artifact preview with escaped code blocks and explicit truncation.

The TUI removes terminal control sequences and non-printing C0 controls except normalized line breaks and tabs. Invalid UTF-8 is replaced before rendering. Long lines wrap or truncate with a visible marker. OSC links, ANSI sequences, raw URLs, and evidence-supplied color codes are never emitted to the terminal.

### Accessibility

The browser uses one `h1`, ordered section headings, landmarks, table captions, row headers, native buttons, and native `details` elements where disclosure fits. Every form control has a visible label. Errors are associated with their control and summarized at the top of the form.

Live progress uses one restrained `aria-live="polite"` region. Cancellation and terminal outcomes move focus to a status heading, but routine progress events do not steal focus. Color never carries assessment, freshness, severity, or cleanup meaning by itself.

The page remains usable at 320 CSS pixels without horizontal page scrolling. Evidence tables become labeled record cards at narrow widths. `prefers-reduced-motion` disables nonessential transitions, and no progress state depends on animation.

The TUI exposes the same information in plain text, preserves a predictable focus order, and includes the active region and item position in its status line. Resize keeps the selected record when possible and moves focus to the nearest visible region when not.

### Bounds and performance

Summary rendering is constant-size. Evidence and issue views default to 25 records and reject page sizes above 100. One rendered untrusted text field is capped at 8 KiB. Views always disclose omitted record counts and text truncation.

The browser does not embed raw smoke JSON, stdout, stderr, full generated patches, or all attempt history in the initial HTML. Detail is requested by immutable attempt and bounded cursor. JavaScript remains dependency-free and adds no build step. Event handling coalesces progress updates and performs one authoritative refresh after terminal state, reconnect, or replay gap.

The TUI keeps only the current summary, current page, selected detail, and bounded progress lines. It does not accumulate the full durable event stream. Writable attempts remain sequential, so the interface does not offer a concurrency control.

### Testing decision

Use shared parity fixtures for one current pass, one failed rerun with retained canonical report, one blocked admission, one cleanup-uncertain attempt, one stale attempt, and one smoke diagnostic-only result.

Required frontend tests are:

- TUI model tests for keyboard navigation, focus restoration, resize, pagination, cancellation confirmation, dropped delivery, replay gaps, and hostile terminal text.
- Handler tests for strict JSON and form decoding, methods, same-origin, CSRF, confirmation binding, fingerprint changes, session rotation, reconnect, cancellation, and restart recovery.
- Template tests for complete no-JavaScript output, semantic headings and labels, narrow-layout classes, bounded records, hostile HTML and Markdown, stale labeling, and current-versus-canonical attempt distinction.
- JavaScript tests or deterministic browser fixtures for enhancement failure, reconnect, duplicate events, replay gaps, cancellation, and authoritative refresh.
- Golden tests only for documented HTML and TUI contracts. State transitions, bounds, and security rules use field-level assertions.
- Cross-surface fixtures that compare app, CLI, TUI, browser, `qa.md`, and verification-state identities, counts, assessment, blocker, smoke status, and next action.

Normal tests use temporary workspaces, fake app use cases, fake durable events, and fake smoke results. Real browser and harness work remains gated.

## Trade-Offs

| Decision | Benefit | Cost | Rejected alternative |
| --- | --- | --- | --- |
| Add QA to the existing sprint page | Keeps project, sprint, artifacts, review, QA, and smoke in one operator workflow. | The sprint page becomes denser and needs disciplined section hierarchy. | A QA-only browser application was rejected because it would duplicate navigation and state loading. |
| Server rendering first | Gives secure escaping, no-JavaScript operation, refresh recovery, and one rendering authority. | Live updates require small enhancement code and full refresh fallback. | A client-side framework and global store were rejected because the interaction does not justify a second state model or build pipeline. |
| Summary first, paged detail | Keeps initial HTML, TUI memory, and reconnect bounded. | Inspecting many evidence records takes more navigation. | Rendering every evidence record and command output was rejected because hostile or large evidence could exhaust the interface. |
| Plain-text hostile evidence | Prevents script, HTML, ANSI, and terminal-control execution. | Model-authored formatting is lost. | Rich Markdown rendering was rejected because evidence is data, not trusted presentation. |
| Current and canonical attempts shown together | Preserves the truthful failed-rerun and retained-report model. | The header has more status fields than a single-verdict design. | Showing only the latest attempt or only `qa.md` was rejected because either hides important state. |
| Durable events for progress, product reads for truth | Supports live feedback without coupling success to a connection. | Reconnect may cause a visible full refresh. | Treating SSE or TUI messages as the operation record was rejected because observers can disconnect or lose events. |
| One shared action flow | Preserves confirmation, authorization, cancellation, and smoke compatibility across interfaces. | TUI and browser must adapt to operation-oriented results. | QA-specific browser workers and TUI command execution were rejected because they would bypass app and run-control authority. |
| Native browser and terminal controls | Keeps keyboard access, focus, and no-JavaScript behavior testable. | The interface is less visually elaborate than a custom widget system. | Custom div-based controls and animated status widgets were rejected because they add accessibility and state complexity without product value. |

## Evidence

The selected reports support these frontend decisions. The exact composition, bounds, and state labels are sprint-specific conclusions drawn from those findings and the governed requirements.

- `studies/go-cli-study/reports/final/01-project-structure.md` finds that mature CLIs keep entry points thin and dependencies one-way. Its UI-interface examples support TUI and web adapters that render app results without importing sprint internals.
- `studies/go-cli-study/reports/final/02-command-architecture.md` documents thin delegation and factory-built commands. The same rule applies to browser forms and TUI actions: parse, authorize, call the shared operation, and render the result.
- `studies/go-cli-study/reports/final/03-dependency-injection.md` favors explicit composition and warns against globals and context service lookup. This supports injected app use cases and local view state rather than a package-global QA store.
- `studies/go-cli-study/reports/final/05-error-handling.md` finds that typed machine decisions and separate user hints produce stable exit and recovery behavior. The UI therefore renders app-provided blocker codes and next actions instead of classifying error strings.
- `studies/go-cli-study/reports/final/06-io-abstraction.md` shows that injectable output and command boundaries make interface behavior testable. Gh CLI's test IO constructor and Restic's mock terminal support deterministic TUI, handler, and rendering fixtures.
- `studies/go-cli-study/reports/final/07-state-context.md` finds that one root context should reach long work while cleanup has its own bounded lifecycle. This supports separate labels for cancellation request, work termination, cleanup certainty, and final assessment.
- `studies/go-cli-study/reports/final/08-concurrency.md` supports localized goroutine ownership, bounded fan-out, and timeout-backed waits. Browser subscribers and TUI observers therefore remain bounded and cannot block the product operation.
- `studies/go-cli-study/reports/final/09-terminal-ux.md` documents interruptible progress, non-TTY fallbacks, resize behavior, and the limits of presentation as task truth. This directly supports server-rendered fallback, durable refresh, keyboard-first TUI operation, and observation-only progress.
- `studies/go-cli-study/reports/final/10-logging-observability.md` supports structured correlation and strict data-versus-diagnostic output separation. The UI shows safe run and attempt identifiers while raw diagnostics remain outside rendered summaries.
- `studies/go-cli-study/reports/final/11-testing-strategy.md` supports command scenarios, fault-capable fakes, and reviewed golden output while warning against assertions on incidental UI details. The test plan freezes semantic HTML, keyboard behavior, security, and parity rather than arbitrary component counts.
- `studies/go-cli-study/reports/final/13-security.md` supports explicit trust boundaries, default-deny permissions, canonical path checks, and redaction. For the frontend, that means escaped plain text, app-issued artifact references, guarded starts, CSRF checks, and no browser-supplied execution mechanics.
- `studies/go-cli-study/reports/final/14-performance.md` supports streaming, bounded buffers, finite concurrency, and measurement before low-level optimization. Paged records, capped text, coalesced progress, and no full event accumulation apply those findings to TUI and browser rendering.

The project documents add rules the comparative reports do not cover. `projects/ultraplan-go/docs/ARCHITECTURE.md` assigns presentation to `internal/tui` and `internal/web`, use-case composition to `internal/app`, QA decisions to `internal/sprint`, and durable lifecycle to run control. `projects/ultraplan-go/docs/PRD.md` and `projects/ultraplan-go/docs/TRD.md` require shared CLI, TUI, and browser results; server-rendered no-JavaScript operation; guarded starts; SSE as observation; cancellation; reconnect; hostile-content safety; and smoke compatibility. The sprint requirements require bounded evidence, adjudication, promoted issues, canonical assessment, failed-rerun preservation, cleanup uncertainty, and exact cross-surface agreement.

## Risks

- A retained passing `qa.md` beside a failed current attempt is easy to misread. Every summary, detail heading, action response, and artifact link must carry both attempt identities and the freshness label.
- Evidence text may contain ANSI, OSC, HTML, Markdown, malformed UTF-8, very long lines, or misleading URLs. Sanitization needs shared fixtures for browser and terminal output, but each adapter still escapes for its own medium.
- A count-only summary can imply success while required evidence is rejected, stale, narrow, or cleanup-uncertain. The assessment and next action must come from product code and remain visually primary.
- Pagination can hide the evidence linked by an issue. Issue rows must provide exact evidence references that open the correct immutable attempt and record, not the current page by position.
- Browser confirmation can become stale between preparation and acceptance. The handler must show the current fingerprint and return to preparation without silently changing the request.
- Session rotation may remove mutation authority while observation remains valid. The interface must keep run and QA reads visible while requiring fresh authorization for cancellation or restart.
- Progress bursts can cause excessive TUI redraws or DOM updates. Coalescing may omit intermediate display events, but it must retain sequence and replay-gap facts and refresh from authoritative state.
- Responsive evidence tables can become unreadable if labels disappear. The narrow layout must repeat field labels per record rather than relying on horizontal scrolling.
- The exact existing TUI key map and CSS tokens may constrain the proposed local bindings and responsive classes. Implementation should preserve established bindings and tokens, but not weaken the keyboard, focus, reduced-motion, or narrow-layout decisions above.
- Accessibility cannot be proven by template snapshots alone. Tests need keyboard and focus scenarios, semantic assertions, reduced-motion coverage, and at least one gated browser pass.
- External smoke references may be unavailable or no longer retained. The UI must show the validated identity and an unavailable-evidence recovery state, not follow arbitrary or stale paths.
- A large number of rejected evidence records can dominate the screen. The fixed page and text bounds are mandatory, and omitted counts must remain visible so truncation is not mistaken for completeness.

The frontend decision is to proceed without a new framework or client-side authority. The hard part is truthful state presentation, especially failed reruns, cleanup uncertainty, and smoke evidence status. The interface should spend its complexity there rather than on custom widgets or animation.
