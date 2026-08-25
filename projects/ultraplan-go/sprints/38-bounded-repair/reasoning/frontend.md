> **Inputs Used:** `projects/ultraplan-go/sprints/38-bounded-repair/requirements.md`, `projects/ultraplan-go/sprints/38-bounded-repair/sprint-index.md`, `projects/ultraplan-go/sprints/38-bounded-repair/technical-handbook.md`, `projects/ultraplan-go/docs/ARCHITECTURE.md`, `projects/ultraplan-go/docs/PRD.md`, `projects/ultraplan-go/docs/TRD.md`, `system/reasoning/frontend-reasoning-template.md`, `studies/go-cli-study/reports/final/01-project-structure.md`, `studies/go-cli-study/reports/final/02-command-architecture.md`, `studies/go-cli-study/reports/final/03-dependency-injection.md`, `studies/go-cli-study/reports/final/04-configuration-management.md`, `studies/go-cli-study/reports/final/05-error-handling.md`, `studies/go-cli-study/reports/final/06-io-abstraction.md`, `studies/go-cli-study/reports/final/07-state-context.md`, `studies/go-cli-study/reports/final/08-concurrency.md`, `studies/go-cli-study/reports/final/09-terminal-ux.md`, `studies/go-cli-study/reports/final/10-logging-observability.md`, `studies/go-cli-study/reports/final/11-testing-strategy.md`, `studies/go-cli-study/reports/final/12-extensibility.md`, `studies/go-cli-study/reports/final/13-security.md`, `studies/go-cli-study/reports/final/14-performance.md`

# Frontend reasoning

Sprint 38 adds a high-consequence workflow to the existing sprint QA views. The operator must understand exactly which issue and production scope a confirmation authorizes, then follow an ordered repair without mistaking progress events for proof. The frontend decision is to extend the existing server-rendered sprint page and TUI QA view with app-owned repair DTOs. It does not add a client-side application, browser-owned state machine, private-file reader, or adapter-specific repair path.

This area covers browser and TUI interaction, rendering, accessibility, reconnect behavior, bounded public data, and interface tests. Product admission, confirmation authority, mutation, reverification, cleanup, and outcome derivation remain outside presentation code.

## Area Decisions

### One product workflow, two presentations

The browser and TUI call the same repair use cases used by CLI and JSON. `internal/web` maps app DTOs to typed template view models. `internal/tui` maps those DTOs to routes and messages. Neither adapter reads `verification/`, interprets patch contents, computes freshness, derives remaining budgets, or maps gate results to a semantic outcome.

The browser remains server-rendered with embedded CSS and small progressive-enhancement scripts. The no-JavaScript path supports packet selection, packet review, operation preparation, confirmation, status refresh, cancellation, resume, recovery, cycle inspection, and terminal result inspection through ordinary links and forms. JavaScript may submit an already defined form asynchronously, attach to durable run events, and refresh bounded status fragments. It cannot confirm on page load, infer completion, or cancel work when the page disconnects.

The selected structure and command reports favor thin interface code and delegated behavior: `studies/go-cli-study/reports/final/01-project-structure.md` and `studies/go-cli-study/reports/final/02-command-architecture.md`. The project architecture fixes the dependency direction as `internal/web -> internal/app` and `internal/tui -> internal/app`. The frontend will preserve that boundary rather than adding repair rules to handlers or models.

### Packet review and confirmation are separate screens and actions

Preparing a repair shows a read-only packet review before any durable mutation confirmation. The review must display the packet ID and digest, issue ID, root-cause group, severity, violated expectation, exact reproducer, expected failing condition, target identity, mode, allowed-path count and bounded preview, forbidden-path count and bounded preview, containing checks, effective limits and their sources, cleanup limit, and stop conditions.

The confirmation page is a second guarded step. It repeats the facts that can invalidate consent: packet digest, target identity, mode, cycle limit, changed-file and changed-byte limits, containing checks, cleanup limit, and stop conditions. The primary action says `Confirm manual repair` or `Confirm automatic repair`; a generic `Continue` label is too weak. The automatic form includes a distinct explicit opt-in control and cannot reuse a manual confirmation token.

GET requests, opening the page, refreshing, reconnecting SSE, restoring browser history, and loading a TUI route are read-only. Confirmation uses the existing CSRF-protected prepared-operation flow. A stale token or changed packet returns a visible stale-request panel with a link to reload the current packet. It never silently re-prepares and submits a changed request.

This design implements AC-2 and AC-3 directly. `studies/go-cli-study/reports/final/13-security.md` supports a visible permission boundary for consequential work, while `studies/go-cli-study/reports/final/04-configuration-management.md` supports showing the fully merged effective limits rather than separate values that have not yet been validated.

### Manual-first gating is explicit

Manual repair is the only enabled start action until the app projection reports a current qualifying manual proof. The frontend does not reduce automatic availability to a configuration toggle.

The automatic section always renders so its absence cannot be confused with a loading bug. Before proof, it shows `Unavailable`, the bounded proof-freshness reasons, and the required next action. After proof, it shows the proof run reference, qualifying outcome, policy/schema match, cleanup verification, and a separate `Prepare automatic repair` action. If proof becomes stale between preparation and confirmation, the server rejects the operation and the UI reloads the current reasons.

Manual runs show a single mutation cycle and never show a `next automatic cycle` control. Automatic runs show consumed and remaining limits, stagnation state, reopenings, repeated-patch detection, deadline state, and the reason another cycle may or may not start. Presentation code does not calculate progress or decide that a limit is exhausted.

### The repair view follows the operator's decision sequence

The sprint repair page uses one route-level page composed from existing template layers. It adds feature-specific components only where the information recurs in packet, status, and result views.

| Page section | Content and behavior |
| --- | --- |
| Repair summary | Current phase, freshness, run lifecycle, semantic outcome, issue and severity, mode, confirmation state, blocker, and next action. |
| Authority | Packet digest, target fingerprint, governed-input and policy fingerprints, confirmer label, confirmation time, and manual-proof status. |
| Scope | Allowed-path and forbidden-path totals with bounded previews, actual changed paths after apply, changed-file and changed-byte counts, and scope-enforcement result. |
| Limits | Effective, consumed, and remaining counters plus absolute deadline and cleanup limit. Sources are shown without unsafe environment values. |
| Reverification ladder | Fixed ordered gates with `pending`, `running`, `passed`, `failed`, `blocked`, or `skipped`, including the reason and next permitted action. |
| Cleanup | Process-tree termination, workspace removal, lock reconciliation, completion, and any uncertainty. |
| Cycle history | Paginated summaries with cycle number, proposal status, actual scope, issue-set delta, severity delta, progress fact, cleanup, and evidence links. |
| Result | One terminal outcome, reason, unresolved issues, completed limits, evidence pointers, and recovery or escalation action. |

The browser uses disclosure elements for long path and evidence previews so the page remains readable without hiding totals or truncation facts. The TUI uses dedicated packet, cycle, and result routes with paging actions. It does not cram every cycle into the main QA screen. Both presentations label the durable operation run ID and repair run ID separately because cancellation and evidence queries use different identities.

### Canonical status wins over events

Progress events show only correlation, phase, cycle, gate, bounded counters, and sanitized text. The browser and TUI reload canonical repair status after terminal events, cancellation acknowledgement, reconnect, replay gap, server restart, route revisit, and periodic refresh while a run is active.

SSE loss changes the connection indicator to `Updates paused`; it does not change run status. The page retains the last canonical snapshot, marks its observation time, and offers an ordinary refresh. Browser disconnect removes only the subscriber. TUI navigation removes only the local observer. Neither action cancels or completes repair.

The state and concurrency evidence supports one root cancellation path and explicit ownership of background work: `studies/go-cli-study/reports/final/07-state-context.md` and `studies/go-cli-study/reports/final/08-concurrency.md`. The terminal evidence supports interruptible progress and non-interactive fallback: `studies/go-cli-study/reports/final/09-terminal-ux.md`. The Sprint 38 inference is stricter: progress is useful feedback, but only the refreshed app projection can state outcome, cleanup, freshness, or remaining authority.

### Every state has a concrete operator action

| State | Browser and TUI behavior |
| --- | --- |
| Loading | Keep identity and last canonical facts visible, mark the requested resource as loading, and disable only actions whose current authority cannot be proven. |
| No eligible issue | Explain that repair requires one current adjudicated repair-eligible issue and link back to QA assessment and issue views. |
| Prepared, unconfirmed | Show the packet review and a separate confirmation action. No progress UI appears. |
| Accepted, not dispatched | Show durable acceptance and confirmation publication state. Prevent duplicate starts while allowing canonical refresh and cancellation where supported. |
| Running | Show cycle, gate, limits, scope facts, cleanup state, and a guarded cancel action. Disable prepare/start actions for conflicting mutation. |
| Cancellation requested | Keep status active until the canonical lifecycle records cancellation, interruption, cleanup uncertainty, or an already authoritative terminal result. |
| Interrupted or resumable | Show the latest proven boundary, consumed limits, deadline, freshness checks, and a fresh guarded resume action. |
| Blocked | Show the missing prerequisite or capability and exact recovery action without implying product failure. |
| Escalated | Show the scope, severity, drift, design, evidence, or cleanup reason requiring human adjudication. Do not offer an automatic retry. |
| Terminal | Show operation lifecycle and repair outcome separately, all completed limits, cleanup proof, unresolved findings, evidence links, and the next action. |
| Retention gap | Show that older detail is unavailable, preserve the canonical current summary, and identify the earliest retained cycle. |

Outcome styling combines text, status badge shape, and icon or prefix. Color never carries meaning alone. `verified_with_findings` is not rendered as ordinary success, and `stalled` appears only for automatic mode. Cancellation and interruption are lifecycle states, not repair outcomes.

### Accessibility and keyboard operation are required behavior

Server-rendered forms use visible labels, fieldsets for mode and automatic opt-in, descriptive button text, and an error summary linked to invalid controls. Confirmation loads with focus on the page heading, not the destructive action. After a failed submission, focus moves to the error summary. After cancellation, resume, or recovery, focus moves to the canonical status heading.

Status updates use a polite live region for phase and gate changes. Cancellation and terminal outcomes use an assertive announcement only once. The full cycle history is not placed in a live region because repeated table announcements would be unusable. TUI routes preserve a stable tab and key order, provide explicit keys for prepare, confirm, cancel, resume, recover, next page, and previous page, and require a second confirmation message before mutation.

Motion is not required to understand progress. Any spinner or reconnect animation stops under reduced-motion preferences. The no-JavaScript page remains complete, and narrow browser widths stack definition groups rather than creating horizontal scrolling for the whole page. Long digests and paths wrap at safe boundaries while retaining their exact displayed value.

### Untrusted retained text is treated as data

Issue claims, paths, blocker details, command diagnostics, provider errors, and evidence labels are hostile input. Browser templates rely on `html/template` escaping and never mark retained repair text as trusted HTML. TUI rendering strips or replaces terminal control sequences. Both use app-level redaction and field bounds before rendering.

Public views exclude proposal patch bodies, production file contents, full prompts, raw provider payloads, unrestricted command output, secrets, unsafe environment values, and private evidence bodies. Evidence links point only to allowlisted bounded query resources. Path previews show total and returned counts so truncation cannot look complete.

The security report gives the relevant external pattern: explicit permission checks, argument arrays, private temporary storage, validation, and secret redaction in `studies/go-cli-study/reports/final/13-security.md`. The logging report supports keeping structured diagnostics separate from user-visible output in `studies/go-cli-study/reports/final/10-logging-observability.md`.

### Rendering and queries remain bounded

The initial HTML response contains the canonical summary, current cycle, fixed reverification ladder, bounded scope previews, latest cleanup facts, and terminal result when present. Cycle history and larger path or evidence collections use paginated app queries. JavaScript does not preload every retained cycle.

Only the active page or selected TUI route refreshes detail. A disconnected or background browser does not keep an unbounded event buffer. Reconnect begins with canonical status and then resumes from a durable cursor when available. A cursor gap is rendered as a gap, not silently replaced with a partial event stream.

`studies/go-cli-study/reports/final/14-performance.md` favors lazy setup, streaming, bounded concurrency, and incremental state. For this UI, that means bounded first render, paged history, fixed-size event projections, and no frontend cache that competes with durable state.

### Tests prove parity and authority, not screenshots alone

Frontend tests use fake app use cases and canonical DTO fixtures. Handler tests cover methods, authorization, CSRF, body limits, stale confirmations, duplicate submissions, escaping, redaction, pagination, reconnect, replay gaps, and no-JavaScript form completion. Template tests render every repair phase, all six semantic outcomes, cancellation, interruption, cleanup uncertainty, missing proof, stale proof, retention gaps, and hostile values.

TUI model tests cover route transitions, separate manual and automatic confirmation, disabled automatic mode, durable start, reconnect, status refresh, cancellation, resume, recovery, paging, narrow terminals, and terminal outcomes. A shared canonical fixture asserts that CLI JSON, TUI view data, browser JSON, and HTML view models carry the same product facts. Browser and TUI tests must assert that event delivery alone never changes the canonical outcome.

Golden fixtures are useful for stable JSON envelopes, help text, and representative complete HTML. They are not evidence for confirmation single use, stale-request rejection, cancellation races, or hostile-text safety. Those require semantic assertions. This follows `studies/go-cli-study/reports/final/11-testing-strategy.md`, which combines command-path tests, fakes, and selective goldens rather than using snapshots as the only proof.

## Trade-Offs

| Decision | Benefit | Cost | Rejected alternative |
| --- | --- | --- | --- |
| Extend the existing sprint QA page and TUI routes | Keeps repair beside the issue and QA evidence that authorize it. | The QA presentation gains more states and needs careful grouping. | A separate repair application would duplicate navigation, status, and durable operation behavior. |
| Server-rendered HTML with progressive enhancement | Confirmation, status, cancellation, resume, and recovery work without JavaScript. | Full-page form transitions are less fluid than a client application. | A frontend framework would add a second state model without a demonstrated need and could make browser state look authoritative. |
| Separate packet review from confirmation | The operator sees the exact authority before any start action. | Manual repair takes one extra step. | One-click prepare-and-start or confirmation on page load cannot prove informed consent. |
| Always render automatic availability | Missing or stale proof has an explicit explanation and next action. | The page shows a disabled section before the feature is usable. | Hiding automatic mode makes proof failure indistinguishable from loading or authorization failure. |
| Canonical reload after event loss or terminal notice | Refresh, reconnect, and server restart remain truthful. | Adds status requests and brief stale-snapshot indicators. | Deriving state from SSE is faster but fails under replay gaps, coalescing, and owner changes. |
| Separate operation lifecycle from repair outcome | Cancellation and semantic verification retain their distinct authorities. | Users see two related status fields. | A single badge would mislabel cancelled, interrupted, or cleanup-uncertain work as a repair verdict. |
| Paginated cycle and evidence detail | Bounds response size and TUI rendering for automatic runs. | Operators may navigate more than one page. | Rendering all retained cycles can exhaust memory and creates unbounded hostile-text exposure. |
| App-level redaction plus adapter escaping | Protects all interfaces and still uses the correct output-specific escaping. | The same text passes through two safety layers. | Adapter-only filtering leaves other interfaces exposed; app-only sanitization cannot replace HTML escaping or terminal-control handling. |
| Purpose-built repair components, not new design-system atoms | Keeps domain-specific labels and ordering clear. | Some markup is local to the repair view. | Premature generic status or workflow components would hide repair-specific authority and stop rules. |

## Evidence

The report findings below are external evidence. The frontend decisions are Sprint 38 inferences constrained by the governed requirements and project documents.

| Decision area | Evidence finding | Sprint 38 conclusion |
| --- | --- | --- |
| Adapter ownership | `studies/go-cli-study/reports/final/01-project-structure.md` finds that mature Go CLIs keep entrypoints thin and preserve one-way imports. `studies/go-cli-study/reports/final/02-command-architecture.md` favors delegated command wrappers. | Browser handlers and TUI actions map typed requests and results. They do not own repair rules. |
| Dependency construction | `studies/go-cli-study/reports/final/03-dependency-injection.md` favors explicit composition roots and warns about global state and context service locators. | Inject repair use cases through existing app dependencies. Do not place confirmation or current repair state in package globals. |
| Effective authority | `studies/go-cli-study/reports/final/04-configuration-management.md` favors explicit precedence and validation after all sources merge. | Confirmation displays the effective lower-only limits and safe source labels returned by the app. |
| Error presentation | `studies/go-cli-study/reports/final/05-error-handling.md` supports typed errors, user and operational separation, and actionable exit mapping. | Render stable blocker categories and next actions. Do not parse messages to select views or outcomes. |
| Testable output | `studies/go-cli-study/reports/final/06-io-abstraction.md` finds that injectable terminal and output boundaries permit behavior tests without real side effects. | Test TUI and web rendering from DTOs and fake use cases. Keep direct filesystem and runtime access out of presentation code. |
| Cancellation lifetime | `studies/go-cli-study/reports/final/07-state-context.md` supports root-context propagation and a separate bounded cleanup lifetime. | UI cancellation requests the canonical operation cancellation path and continues to show cleanup until canonical state settles. |
| Observer ownership | `studies/go-cli-study/reports/final/08-concurrency.md` favors localized launch sites, explicit waits, and no fire-and-forget work. | Browser scripts and TUI observers never launch repair workers. Disconnecting an observer cannot alter work ownership. |
| Interactive fallback | `studies/go-cli-study/reports/final/09-terminal-ux.md` supports non-TTY fallbacks, visible interruptible progress, and signal-safe behavior. | Preserve complete no-JavaScript operation and guarded keyboard flows while keeping progress observational. |
| Correlated diagnostics | `studies/go-cli-study/reports/final/10-logging-observability.md` favors structured correlation and separation of diagnostics from user output. | Display bounded run correlations and safe facts; keep raw provider and command diagnostics out of public views. |
| Test layers | `studies/go-cli-study/reports/final/11-testing-strategy.md` supports table tests, command-path integration, fakes, and selective golden fixtures. | Use semantic tests for authority and races, plus limited fixtures for stable public rendering. |
| Extension restraint | `studies/go-cli-study/reports/final/12-extensibility.md` documents the versioning and trust cost of plugin systems and dynamic registries. | Do not add a frontend framework, plugin view registry, or user-defined repair controls for one workflow. |
| Mutation warning | `studies/go-cli-study/reports/final/13-security.md` supports explicit permission gates, validation, and redaction for consequential operations. | Use a separate confirmation page, explicit mode labels, CSRF protection, stale rejection, and hostile-text defenses. |
| Bounded views | `studies/go-cli-study/reports/final/14-performance.md` favors streaming, lazy setup, incremental state, and bounded collections. | Page cycle history and evidence, keep event payloads finite, and avoid loading all retained repair detail into one screen. |

`projects/ultraplan-go/docs/ARCHITECTURE.md` assigns browser transport and templates to `internal/web`, terminal presentation to `internal/tui`, and shared DTOs and use cases to `internal/app`. `projects/ultraplan-go/docs/PRD.md` requires a run to outlive its observer and all local interfaces to agree. `projects/ultraplan-go/docs/TRD.md` requires server-rendered no-JavaScript behavior, guarded mutations, strict browser security, durable status refresh, and browser disconnect independence.

The governed acceptance criteria provide the frontend contract. AC-2 requires visible packet facts and separate confirmation. AC-3 and AC-6 require manual-first delivery and current proof before automatic controls. AC-5 requires the fixed gate order. AC-8 and AC-9 require truthful outcomes, cancellation, cleanup, resume, and recovery. AC-10 requires bounded, escaped, redacted parity across CLI, JSON, TUI, and browser.

## Risks

### Confirmation facts can drift between review and submit

The packet or target may change while the user reads the confirmation page. The server must reject the stale token before mutation. The UI must show which identity changed and require a fresh packet review; it cannot preserve checked controls and silently submit a new authority.

### Automatic availability can be overstated

A proof file may exist but be stale, weaker, or tied to another policy. The frontend must use app-provided availability and mismatch reasons. Any local rule based on proof presence or configuration would expose an unsafe start action.

### Dense repair facts can bury the decision

Packet, limits, scope, gates, cleanup, and cycle history can overwhelm one screen. The chosen order keeps current authority and next action first, followed by scope, progress, cleanup, and history. Usability tests should verify that operators can identify the issue, mode, target, mutable paths, current gate, and next action without opening every disclosure.

### Dual identifiers can cause the wrong action

The repair run ID owns packet and cycle evidence, while the durable operation run ID owns cancellation and observation. Labels, links, and forms must name both explicitly and submit the correct one. Tests should reject cross-sprint and swapped identifiers.

### Progress can look authoritative after ownership changes

A late SSE event or TUI message may arrive after cancellation, recovery, or terminal publication. Event handlers must never set semantic outcome or cleanup state. They trigger a canonical reload and may update only an explicitly labeled transient progress line.

### Hostile text can attack both presentation types

HTML injection, terminal control sequences, long unbroken paths, secrets embedded in diagnostics, and huge retained text require different output defenses. App-level bounds and redaction are necessary but not sufficient. Web escaping and TUI control-character handling need separate tests over the same hostile fixture.

### Disabled controls can hide recovery

When mutation is unavailable, disabling every action would strand the operator. Read-only refresh, packet inspection, result inspection, cancellation status, recovery guidance, and applicable recover operations remain available. Each disabled mutation control must state why and identify the next allowed action.

### No-JavaScript and enhanced paths can diverge

Progressive enhancement can accidentally bypass form validation or submit a different request shape. Browser tests must exercise both ordinary forms and enhanced requests against the same operation preparation and confirmation use cases. JavaScript remains optional and cannot introduce repair-only authority.

### Accessibility can regress under rapid status updates

Announcing every event would flood assistive technology. Only canonical phase, gate, cancellation, and terminal changes enter live regions. Cycle tables and diagnostic detail update without announcements unless the user requests them.

### Decision

Proceed by extending the existing sprint QA presentation with app-owned repair DTOs, a read-only packet review, a separate guarded confirmation page, explicit manual-proof gating, bounded canonical status and cycle views, and complete no-JavaScript operation. TUI and browser show the same facts but use presentation-specific navigation and escaping. Events improve responsiveness; canonical status remains authoritative.
