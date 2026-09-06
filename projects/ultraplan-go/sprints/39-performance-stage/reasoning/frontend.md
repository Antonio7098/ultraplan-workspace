> **Inputs Used:** `projects/ultraplan-go/sprints/39-performance-stage/requirements.md`, `projects/ultraplan-go/sprints/39-performance-stage/sprint-index.md`, `projects/ultraplan-go/sprints/39-performance-stage/technical-handbook.md`, `projects/ultraplan-go/sprints/39-performance-stage/reasoning/architecture.md`, `projects/ultraplan-go/sprints/39-performance-stage/reasoning/api-design.md`, `projects/ultraplan-go/docs/ARCHITECTURE.md`, `projects/ultraplan-go/docs/PRD.md`, `projects/ultraplan-go/docs/TRD.md`, `PRODUCT.md`, `DESIGN.md`, `system/reasoning/frontend-reasoning-template.md`, `system/contracts/surfaces/frontend.md`, `system/contracts/surfaces/accessibility.md`, `system/contracts/surfaces/api-contracts.md`, `system/contracts/core/testing.md`, `system/contracts/core/privacy-and-data.md`, `system/contracts/runtime/performance.md`, `system/contracts/core/observability.md`, `studies/go-cli-study/reports/final/01-project-structure.md`, `studies/go-cli-study/reports/final/02-command-architecture.md`, `studies/go-cli-study/reports/final/03-dependency-injection.md`, `studies/go-cli-study/reports/final/05-error-handling.md`, `studies/go-cli-study/reports/final/06-io-abstraction.md`, `studies/go-cli-study/reports/final/07-state-context.md`, `studies/go-cli-study/reports/final/08-concurrency.md`, `studies/go-cli-study/reports/final/09-terminal-ux.md`, `studies/go-cli-study/reports/final/10-logging-observability.md`, `studies/go-cli-study/reports/final/11-testing-strategy.md`, `studies/go-cli-study/reports/final/13-security.md`, `studies/go-cli-study/reports/final/14-performance.md`, `internal/web/templates/sprint.html`, `internal/web/templates/operation-confirm.html`, `internal/web/templates/operation.html`, `internal/web/templates/run_qa.html`, `internal/tui/qa_view.go`, `internal/app/sprint_usecases.go`

# Frontend: requirements-driven performance stage

This area decides how operators inspect and control the optional performance phase in the browser and TUI. It covers placement, composition, state presentation, accessibility, progressive enhancement, bounded evidence, and interface tests. It does not define target parsing, measurement qualification, product verdicts, private record schemas, or durable run arbitration.

## Area Decisions

### 1. Put performance beside the sprint it gates

The browser gets a focused sprint route at `/projects/{project}/sprints/{sprint}/performance`. The TUI gets the equivalent sprint-owned Performance destination. Performance does not belong on a generic operations page because its targets, freshness, blockers, and next action all belong to one sprint.

The sprint overview shows a compact performance gate only when policy is enabled or performance state already exists. It reports policy, freshness, required-target counts, run outcome, blocker, and next action, then links to the focused workbench. The workflow view shows performance between execute and Conformance Review in a separate verification path. It must not add performance to the planning-stage list or imply that a completed operation passed the target gate.

When policy is missing or disabled, the overview preserves the current `execute -> Conformance Review` presentation. A direct visit to the performance route may show a quiet disabled explanation and the project-index path, but it exposes no start, resume, or evidence controls and triggers no target, benchmark, runtime, command, or state work.

This placement follows the product rule that commands sit beside the entity and state they affect. It also keeps Active Runs useful: an accepted performance run remains globally discoverable and links back to its sprint workbench.

### 2. Extend the existing server-rendered interface instead of adding a client application

The browser remains Go-rendered HTML with dependency-free JavaScript for optional live refresh. Handlers call the typed performance application use cases and map their bounded DTOs into explicit page and component view models. Templates render those view models only. They do not read workspace files, parse private JSON, interpret runtime messages, or decide whether a target passed.

Use the established namespaced hierarchy:

```text
page/performance
  -> layout/detail
  -> component/performance-summary
  -> component/performance-path
  -> component/performance-targets
  -> component/performance-evidence
  -> existing confirmation, blocker, artifact-reference, status, and button primitives
```

The performance components remain feature-specific until another verification phase proves the same facts and behavior. Shared status badges, progress, metadata lists, notices, buttons, and artifact references can be reused. A generic verification-workbench framework is not earned by visual similarity alone.

The TUI follows the same boundary. Its model stores the latest bounded app result plus local selection and viewport state. Rendering code must not read `verification/performance-state.json`, parse `performance.md`, inspect raw samples, or infer outcomes from durable events.

### 3. Lead with current authority, not activity

The workbench uses the existing dark operating-dashboard language from `DESIGN.md`. One dominant current-state panel answers these questions in order:

1. Is performance enabled and current for this sprint?
2. What product phase and durable operation lifecycle are active?
3. How many required and report targets are met, missed, inconclusive, blocked, or baseline-only?
4. Is later verification allowed to proceed?
5. What is the next valid action?

The panel displays operational lifecycle and product outcome as separate labeled facts. An operationally `succeeded` run with product outcome `target_missed` is not styled or worded as a pass. `passed` and `passed_with_reports` are the only non-blocking terminal outcomes. `cancelled`, `cleanup_uncertain`, `stalled`, and `blocked` remain distinct.

The ordered path shows product-owned phase summaries supplied by the application DTO. Its conceptual order is admission, benchmark coverage and freeze, baseline, target-linked optimization when needed, final measurement, correctness, and publication. Presenters render supplied state such as complete, active, skipped, failed, and pending. They do not reconstruct phase completion from event text or artifact presence.

Live events are activity hints below the canonical snapshot. They may show bounded phase, target ID, cycle, completed and total counts, reason code, and evidence reference. They cannot replace the current status, target verdicts, or terminal result.

### 4. Show the complete target contract and a concise result in one record

Each target record displays every governed row field: ID, scenario, metric, comparator, normalized value, unit, gate, samples, and basis. The same record then displays qualification, bounded baseline and candidate aggregates, the comparison fact, outcome, and evidence reference. Missing values read as "not measured" or "not qualified" rather than zero.

Use a scan-friendly target list instead of one very wide table. Each item leads with target ID, scenario, exact target expression, gate, and textual outcome. A native `details` disclosure contains the remaining contract fields, measurement summary, comparison, and evidence reference. This keeps the full contract available on narrow screens without duplicating markup into separate desktop and mobile renderings.

Target outcome always includes text. Semantic color follows the existing state rule but never carries meaning alone:

| Outcome | Visible treatment |
| --- | --- |
| `met` | "Met" with completion state |
| `missed` | "Missed" with required or report gate shown beside it |
| `baseline_recorded` | "Baseline recorded" with the qualified value |
| `report_only` | "Report only" with the observed comparison |
| `inconclusive` | "Inconclusive" with the qualification reason and rerun action |
| `blocked` | "Blocked" with reason code and next action |

A required miss, inconclusive result, or blocker gets a compact attention region near the current-state panel. Report-only findings stay visible but do not receive the same blocking weight. Baseline-only rows never show optimization controls.

### 5. Make preparation and every state transition explicit

Preparation and dry-run are read-only views of deterministic admission facts. They show policy, all target rows, expected benchmark coverage, correctness commands, protected roots, effective limits and sources, missing prerequisites, possible benchmark or implementation promotion, affected paths, runtime and model source, and the confirmation expiry. They never show a progress state because they start no runtime or command.

Start and resume use the existing server-issued confirmation flow. The final confirmation button states the action plainly, for example "Confirm and start performance". A stale confirmation returns to a visible conflict summary with a link or form to prepare again. The browser does not silently refresh the request and proceed against changed inputs.

Available controls come from product state rather than route conventions:

| State | Primary controls |
| --- | --- |
| Enabled, not prepared | Dry run and prepare start |
| Prepared and current | Confirm start or return to sprint |
| Active | Open durable run and request cancellation |
| Interrupted with resumable frozen state | Prepare resume and recover state |
| Stale | Inspect changed authority and prepare a new current run |
| Terminal non-blocking | Continue to Conformance Review and inspect result |
| Terminal blocking | Inspect blocker, bounded evidence, and permitted retry or recovery |
| Cleanup uncertain | Recover state; do not offer an ordinary start that could overlap work |

Cancellation is an explicit idempotent form or button routed through the canonical run cancellation use case. The UI changes its label to "Cancellation requested" when that fact is durable and keeps cleanup progress visible. Navigating away, closing a TUI view, refreshing a page, losing a browser session, or losing SSE never sends cancellation and never marks the run complete.

Recovery is presented as conservative reconciliation, not retry. Its preparation explains what state it may inspect or publish and that it cannot infer success, reset consumed limits, launch runtime work, or widen frozen authority.

### 6. Treat stale and historical state as first-class states

Stale performance evidence remains inspectable but loses current-state styling. The workbench shows the stable freshness reason codes, a safe summary of which fingerprint changed, the blocking effect, the historical result reference, and the exact rerun action. It must not show stale `met` targets as satisfying the current verification gate.

If a run page refers to an older performance attempt, it shows a historical notice and links to the sprint's current performance attempt. The event journal remains scoped to the selected durable run. Current target facts must not be mixed into that historical journal.

The browser and TUI use app-supplied references for `performance.md`, current state, and terminal result. Paths are workspace-relative and paired with digests where available. The presenters do not decide freshness from timestamps.

### 7. Expose bounded evidence metadata, never hostile detail by default

The evidence region uses the paged app query defined by the API reasoning. It shows stable metadata such as attempt, target, cycle, kind, bounded summary, path, digest, size or omission facts, and containment or cleanup status where relevant. Filters are limited to the closed attempt, target, cycle, and kind set. Pagination uses server-generated links in HTML and the opaque cursor in enhanced requests.

The default page size is 50 and the maximum is 200. Retention gaps and stale cursors are visible states with a reset action. Opening content goes through the existing allowlisted bounded artifact preview. The workbench itself never receives or renders raw samples, full command output, profiles, patches, prompts, provider payloads, environment values, secrets, private record structs, or arbitrary workspace paths.

This is deliberate. Operators get enough identity, qualification, comparison, correlation, and omission data to judge the result without turning the browser or TUI into a second private-evidence store.

### 8. Keep JavaScript optional and subordinate to committed state

Every read, confirmation, start, resume, cancellation, recovery, filter, pagination, result, and artifact-preview path works with ordinary links and forms. Without JavaScript, an accepted action lands on the operation page. The page states that the run continues independently and provides explicit Refresh status and Cancel run forms.

With JavaScript, the browser subscribes to durable run events from the last committed cursor. A committed event schedules one debounced fetch of the canonical performance HTML snapshot. Only one snapshot request may be in flight. Reconnect uses the durable cursor; a replay gap or fetch failure leaves the last committed snapshot visible and shows a refresh control. There is no independent polling loop while SSE is healthy.

Snapshot replacement preserves open native disclosures by stable semantic IDs and does not replace the focused subtree while focus is inside it. It updates after focus leaves or the user requests refresh. JavaScript holds only connection, cursor, disclosure, and transient request state. It is not a router, product store, outcome calculator, or cancellation owner.

### 9. Apply the same accessibility contract to browser and TUI

The browser uses native headings, landmarks, lists, `details`, links, buttons, labels, forms, and `progress` before ARIA. The performance path includes visible text for every phase state and a visible legend. Target and run outcomes use words and programmatic text in addition to color and shape. Counts use tabular numerals, but no percentage or count stands alone without its label and denominator.

Keyboard order follows the visual reading order: breadcrumb, current state, next action, target list, evidence, then secondary identity and limits. Forms keep explicit labels and associate field errors with `aria-describedby`. On validation failure, focus moves to the error summary; on dialog close, it returns to the invoking control. Focus is never moved merely because an SSE event arrived.

One polite live region announces coalesced phase changes, cancellation acknowledgement, blockers, reconnect gaps, and terminal outcome. It does not announce every sample or replace the visible event journal. Unexpected request failures use an alert. Known stale, blocked, and cancellation states use status unless immediate interruption is required.

Motion is limited to short state transitions and the existing active-state cue. `prefers-reduced-motion: reduce` disables those transitions, animated markers, and smooth scrolling. The layout stacks in source order below 48rem. Controls keep the existing minimum target size and visible focus treatment.

The TUI provides the same textual states, not color-only glyphs. Every action has a named key binding shown in the help model, selection remains visible, narrow terminals fall back to one column, and leaving the view does not cancel work. Screen output keeps the current phase, outcome, blocker, and next action readable without animation.

### 10. Test the state matrix and authority boundary

Use one set of additive canonical performance fixtures to drive app projection, CLI JSON, browser view-model, rendered HTML, and TUI rendering tests. The fixtures cover disabled, ready, preparing, active, cancelling, passed, passed-with-reports, target-missed, blocked, cancelled, cleanup-uncertain, stalled, stale, interrupted/resumable, historical, evidence-gap, and recovery-conflict states.

Focused web tests must prove:

- route and method behavior, typed error mapping, CSRF and stale-confirmation rejection;
- complete server-rendered start, resume, cancel, recover, evidence pagination, and refresh flows with JavaScript absent;
- explicit labels, heading order, landmarks, live-region policy, focus targets, textual status, and keyboard-operable native controls;
- safe escaping and omission of hostile scenarios, summaries, paths, runtime text, raw samples, profiles, patches, prompts, provider payloads, environment values, and secrets;
- bounded target and evidence rendering, cursor gaps, one in-flight snapshot fetch, SSE reconnect, slow subscriber behavior, and preservation of focus and open disclosures;
- browser refresh, tab loss, session expiry, and SSE loss do not cancel or complete a run.

Focused TUI tests use deterministic messages and fake app use cases. They cover narrow and wide rendering, target selection, textual outcome distinctions, guarded actions, cancellation acknowledgement, stale and historical results, recovery, and observer exit without operation cancellation. Browser-level smoke covers the critical keyboard and no-JavaScript journeys against a real local server with fake runtime and process boundaries. Exact cosmetic wrapping is not a contract; state, authority, controls, and safe content are.

## Trade-Offs

| Decision | Benefit | Cost | Rejected alternative |
| --- | --- | --- | --- |
| Focused sprint workbench plus bounded overview summary | Keeps the phase next to its target authority while preserving a scannable overview | Adds one route and TUI destination | Putting all controls on the overview would overload it; a generic operations page would separate commands from the sprint they change. |
| Server-rendered HTML with optional SSE refresh | Preserves local build simplicity, first-load truth, accessibility, and no-JavaScript operation | Snapshot replacement needs careful focus and disclosure preservation | A client-side application would duplicate routing and state authority before the interaction complexity earns it. |
| Canonical snapshot after event hints | Every refresh comes from one product projection | One extra bounded GET follows meaningful committed events | Rendering product state directly from event payloads would be faster but could turn incomplete activity into a verdict. |
| Feature-specific performance components | Names and layouts can match target and measurement facts | Some markup resembles QA and repair views | A generic verification UI framework would force unlike phases into one schema and widen this sprint. |
| Expandable target records instead of one wide table | Keeps every governed field available and works on narrow screens with one semantic representation | Comparing many secondary fields takes an extra disclosure action | A wide table would require horizontal scanning; separate mobile cards would duplicate markup and drift. |
| Bounded evidence metadata and allowlisted preview | Keeps status fast and prevents hostile private data from becoming public interface state | Deep diagnosis requires an explicit local evidence step | Embedding raw records would expose secrets and unbounded content and couple the UI to private schemas. |
| Separate operational lifecycle and product outcome | A completed run can truthfully remain target-missed or blocked | Presenters and fixtures carry two state axes | One combined badge would make operational success look like target success. |
| Coalesced accessibility announcements | Screen-reader users hear meaningful transitions without sample-level noise | Fine-grained activity remains visual or available in the event journal | Announcing every event would make a long benchmark run unusable. |

## Evidence

- Sprint-specific authority comes from `projects/ultraplan-go/sprints/39-performance-stage/requirements.md:95-104` and `requirements.md:194-220`. Those lines require adapter-independent performance DTOs, browser and TUI parity, exact target and run outcomes, stale-state visibility, canonical cancellation, recovery, and observer independence. `requirements.md:151-175` and `requirements.md:206-220` make preparation runtime-free, reserve verdicts for product code, and require every interface to show the rerun action after drift.
- `projects/ultraplan-go/sprints/39-performance-stage/reasoning/api-design.md:47-67` defines the bounded public fact model used here, and `api-design.md:96-120` defines shared run resources, sprint product queries, paged evidence, and browser-disconnect semantics. This frontend reasoning presents those decisions; it does not add transport fields or a second operation protocol.
- `projects/ultraplan-go/sprints/39-performance-stage/reasoning/architecture.md:29-45` places performance between execute and Conformance Review while keeping it outside planning-stage semantics. `architecture.md:84-118` makes bounded current state, digest-bound references, product-owned publication, freshness, and finite retained work the source for the display rules above.
- `projects/ultraplan-go/sprints/39-performance-stage/technical-handbook.md:29-53` distills thin adapters, typed outcomes, bounded execution, cancellation, bounded retention, profile-led optimization, and one fact model across interfaces. `technical-handbook.md:112-125` specifically warns that a browser disconnect cannot own cancellation and that public DTOs must remain much smaller than private evidence.
- `projects/ultraplan-go/docs/ARCHITECTURE.md:303-368` assigns typed use cases to `internal/app`, presentation to TUI and web, and durable cancellation to run control. `ARCHITECTURE.md:390-405` requires server rendering, the `page -> layout -> component -> primitive` hierarchy, explicit view models, namespaced templates, and capability-focused JavaScript.
- `projects/ultraplan-go/docs/PRD.md:400-423` requires Go-rendered progressive enhancement, typed handlers, bounded previews, guarded confirmation, durable runs, and truthful disconnect behavior. `PRD.md:782-815` and `PRD.md:1346-1354` establish the performance placement, target authority, stale-state rules, and cross-interface agreement used in this design.
- `projects/ultraplan-go/docs/TRD.md:366-408` requires complete server-rendered browser behavior, accessibility, bounded SSE, and local security. `TRD.md:2129-2205` fixes web ownership, transport direction, disconnect and shutdown semantics, template layering, and no-JavaScript tests. `TRD.md:2291-2328` supplies the exact Sprint 39 target, phase, outcome, evidence, and freshness facts.
- `PRODUCT.md:25-49` places entity-scoped commands, Projects, Studies, and Active Runs in the established navigation model and requires keyboard, live-region, reduced-motion, and no-JavaScript behavior. `DESIGN.md:172-184`, `DESIGN.md:206-268`, and `DESIGN.md:270-292` provide the incumbent responsive dashboard, current-state priority, component vocabulary, status, progress, live-refresh, and evidence-shape rules.
- `system/contracts/surfaces/frontend.md:280-319` supports server, route, and local-state separation plus typed DTO mapping. `frontend.md:357-419` requires accessible UI, bounded rendering and fetching, and complete form states. Because UltraPlan uses Go templates rather than the contract's example React directory tree, the project-specific namespaced template hierarchy is the applicable ownership mechanism.
- `system/contracts/surfaces/accessibility.md:51-175` requires semantic HTML, keyboard operation, deliberate focus, accessible form status, non-color state, reduced motion, and accessibility checks. These are direct requirements for the performance workbench, confirmation, cancellation, disclosures, and live updates.
- `system/contracts/surfaces/api-contracts.md:52-71`, `api-contracts.md:117-170`, and `api-contracts.md:135-152` support explicit DTOs, idempotent write behavior, stable errors, and bounded cursor pagination. `system/contracts/core/privacy-and-data.md:79-96` and `privacy-and-data.md:118-174` support references, hashes, summaries, explicit field allowlists, and omission of raw private content.
- `system/contracts/runtime/performance.md:83-119` requires bounded work, useful progress, stable UI loading, and no duplicate requests or scans. `performance.md:139-178` requires owned cancellation and visible cost drivers. These findings support bounded snapshots and evidence pages, but they do not justify speculative browser caching or a new frontend dependency.
- `system/contracts/core/observability.md:45-96` and `observability.md:168-192` require structured failures, durable correlation, and explicit terminal state. `system/contracts/core/testing.md:122-190` requires negative, compatibility, and end-to-end coverage based on stable behavior rather than cosmetic phrasing.
- The selected reports are comparative evidence, not Sprint 39 authority. `studies/go-cli-study/reports/final/01-project-structure.md:32-40`, `02-command-architecture.md:32-38`, and `03-dependency-injection.md:32-58` support thin interface adapters, one-way dependencies, explicit construction, and avoidance of global state. Applied here, web and TUI render app facts rather than owning performance logic.
- `studies/go-cli-study/reports/final/05-error-handling.md:32-38` supports typed machine decisions with separate user guidance. `06-io-abstraction.md:32-38` supports replaceable boundaries and deterministic view tests. These findings support stable error codes and fake app use cases; the exact performance states still come from sprint requirements.
- `studies/go-cli-study/reports/final/07-state-context.md:32-48` and `08-concurrency.md:32-46` support propagated cancellation, localized asynchronous ownership, explicit waits, and bounded fan-out. They reinforce the decision that page, session, and TUI lifecycles observe rather than own the durable run.
- `studies/go-cli-study/reports/final/09-terminal-ux.md:80-102` and `09-terminal-ux.md:149-157` support interruptible progress and non-interactive fallbacks for long work. `10-logging-observability.md:32-38` supports structured operational facts and separation of user output from diagnostics. The Sprint 39 inference is to render bounded canonical facts while keeping raw runtime and command detail out of the interface.
- `studies/go-cli-study/reports/final/11-testing-strategy.md:32-38` and `11-testing-strategy.md:85-100` support behavior-focused integration tests, shared fakes, and deliberate output fixtures. `13-security.md:115-125` and `13-security.md:179-187` support explicit argument boundaries and private temporary evidence, which argues against exposing private records or executable strings in browser controls.
- `studies/go-cli-study/reports/final/14-performance.md:32-40`, `14-performance.md:89-111`, and `14-performance.md:123-166` support lazy setup, incremental bounded data, streaming, and bounded concurrency. They support paged evidence and event-triggered snapshots. They do not supply the target thresholds, qualification rules, or product verdicts shown by the UI.
- `internal/web/templates/sprint.html:1-40` shows the current sprint dashboard, entity breadcrumbs, adjacent operation forms, textual status, semantic progress, and mobile-compatible source order. `internal/web/templates/operation-confirm.html:1-28` and `internal/web/templates/operation.html:11-27` provide the current confirmation, durable operation, cancellation, refresh, and no-JavaScript patterns to extend.
- `internal/web/templates/run_qa.html:1-43` and `run_qa.html:108-125` demonstrate the closest existing canonical-snapshot pattern: current phase and freshness first, events subordinate to product state, native disclosures, identity and limits in a secondary rail, and no private evidence parsing. `internal/tui/qa_view.go:8-53` explicitly renders only the bounded app projection. `internal/app/sprint_usecases.go:78-94` and `sprint_usecases.go:196-255` demonstrate typed use-case and DTO boundaries for verification interfaces.

## Risks

| Risk | Consequence | Required control |
| --- | --- | --- |
| The overview treats performance as another planning stage | Flow order and readiness become misleading | Render a separate verification gate and path; never insert it into the planning-stage enum or derive readiness in templates. |
| Operational success is styled as a performance pass | A valid `target_missed` result appears safe for Conformance Review | Show lifecycle and product outcome as separate text fields and fixture-test every valid combination. |
| Disabled projects gain eager status work | Existing projects incur scans, files, or changed bytes without opting in | Branch on the app-supplied disabled policy before target, benchmark, runtime, command, or performance-state work; keep the disabled page read-only. |
| A stale result retains green current-state treatment | Old measurements satisfy the visible gate after source or requirement changes | Give freshness priority over stored target outcomes and show the changed fingerprint, blocking effect, and rerun action. |
| Event payloads become a shadow verdict source | Dropped, reordered, or replay-gapped events produce false phase or target state | Use events only to trigger canonical snapshot reads and render all verdicts from bounded app DTOs. |
| Live replacement closes disclosures or removes focus | Operators lose position or keyboard context during long runs | Preserve disclosures by stable IDs, defer focused-subtree replacement, never move focus on an SSE event, and test reconnect updates. |
| Sample-level announcements overwhelm assistive technology | The interface becomes unusable during repeated measurements | Coalesce announcements to phase, cancellation, blocker, gap, and terminal transitions; keep detail in an explicit journal. |
| Target records hide contract fields on narrow screens | Users cannot verify what was measured | Keep every governed field in the same semantic record and use native disclosures that work without JavaScript. |
| Raw private evidence leaks through HTML, JSON, or TUI | Secrets, hostile output, large payloads, or private schemas become public contracts | Build allowlisted view models from app DTOs, paginate metadata, use bounded artifact preview, and add explicit omission and hostile-content tests. |
| Evidence pagination races retention or a new attempt | The browser mixes attempts or silently skips records | Bind cursors to attempt and filters, display stale cursor and retention-gap states, and offer a reset to the current attempt. |
| JavaScript retries duplicate mutation | A reconnect or network retry starts or promotes work twice | Never auto-retry start, resume, cancel, or recover in JavaScript; rely on explicit forms, confirmation fingerprints, and app idempotency. |
| Cleanup uncertainty offers a normal new start | Two attempts may overlap while process state is unknown | Make recovery the primary action and suppress ordinary start until product state proves a safe boundary. |
| Similar QA, repair, and performance screens grow a premature framework | Product-specific authority and outcomes become hidden behind generic UI state | Reuse presentation primitives only; keep performance composition and view models feature-specific until repeated behavior proves a smaller shared component. |

No frontend decision remains open for final sprint reasoning. Concrete Go type and template file names may change during planning, but browser and TUI implementation must preserve the placement, bounded DTO boundary, state distinctions, accessibility, no-JavaScript flows, and observer-independent cancellation rules above.
