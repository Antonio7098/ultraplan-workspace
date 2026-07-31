> **Inputs Used:** `projects/ultraplan-go/sprints/30-web-foundations/sprint-index.md`, `projects/ultraplan-go/sprints/30-web-foundations/technical-handbook.md`, `projects/ultraplan-go/sprints/30-web-foundations/requirements.md`, `projects/ultraplan-go/docs/ARCHITECTURE.md`, `projects/ultraplan-go/docs/PRD.md`, `projects/ultraplan-go/docs/TRD.md`, `system/reasoning/frontend-reasoning-template.md`

# Frontend: Read-Only Local Dashboard

This area covers Sprint 30's browser navigation, server-rendered pages, status presentation, validation findings, and bounded Markdown/JSON artifact previews. The frontend is a local inspection surface over typed app queries. It does not start or cancel work, edit artifacts, own durable state, or require a separate frontend runtime.

## Area Decisions

### Ownership and composition

- Keep all browser presentation in `internal/web`: embedded `html/template` pages and partials, transport-specific view models, embedded CSS, and minimal JavaScript. Product interpretation remains in typed `internal/app` results; templates do not read the filesystem or call product modules.
- Build pages from a small shared shell containing a skip link, header, primary navigation, workspace context, main content landmark, and footer/server status. Page-specific templates own dashboard, project, sprint, study, validation, artifact, and error content. Reuse partials only for repeated presentation such as status badges, finding lists, metadata rows, and empty/error panels; do not introduce a general component framework.
- Render useful, complete HTML on the initial response. JavaScript may progressively refresh explicitly marked status regions and preserve ordinary link navigation, but it must not be required to discover content, inspect details, read errors, or preview artifacts.
- Use the route hierarchy as the information architecture: dashboard first, then projects or studies, then project/sprint or study detail, then allowlisted artifact preview. Breadcrumbs and page titles expose the current location. Artifact and detail links come from app-issued identifiers/references rather than user-entered paths.

### Page behavior and states

| Surface | Primary content | Empty and failure behavior |
| --- | --- | --- |
| Dashboard | Workspace summary, bounded recent/current project, sprint, study, validation, and run/flow status | Named empty panels explain that no matching artifacts exist; a query failure renders a safe page-level error with a normal refresh link. |
| Project list/detail | Project summaries, documentation availability, sprint summaries, validation state | Empty sprint lists remain valid; missing/stale identifiers render `404` without exposing paths. |
| Sprint detail | Planning stages, execute/review/smoke summaries, current flow/run state, validation findings, artifact links | Unknown values display “Not available,” not zero or success; failures retain unaffected explanatory metadata where supplied by the app result. |
| Study list/detail | Study summaries, dimensions/sources, run state, validation, report links | Inapplicable, pending, failed, and missing states remain visually and textually distinct. |
| Artifact preview | Safe metadata and bounded Markdown/JSON source text | Rejected, stale, unsupported, or escaping references render the same safe not-found state; truncation is announced next to the preview. |
| Health/error | Server/workspace readiness or a safe actionable error | Health does not imply artifact validity; internal causes, absolute paths, secrets, and raw diagnostics are omitted. |

The normal loading state is browser document navigation; no client-side skeleton is needed. For progressive refresh, retain the current rendered region while a small textual “Refreshing” status is exposed, then replace only that region on success. Refresh failure leaves the prior content visible, marks it as potentially stale, and offers an ordinary page reload. The UI has no mutation-pending or disabled action state because Sprint 30 registers no browser operations.

### State and data placement

- Server state consists only of the typed app result used for the current response. Each navigation or explicit refresh reads authoritative workspace/product state again; rendered pages use `Cache-Control: no-store` and do not establish a client cache as a second source of truth.
- URL state is limited to route identifiers and the documented validation scope/reference query. It must be bookmarkable and must not contain absolute paths, secrets, serialized product models, or arbitrary filesystem references.
- Local UI state is limited to ephemeral presentation preferences such as an open disclosure element or current focus. It need not survive refresh and is not persisted in cookies, local storage, or a database.
- There is no global client state, service worker, browser session model, optimistic update, polling loop, operation store, or SSE connection in Sprint 30.
- HTML handlers and JSON handlers independently map the same plain app results into template view models or API DTOs. Browser JavaScript may consume only the documented `/api/v1` success/error envelopes and must treat omitted values, empty collections, and truncation metadata explicitly.

### Safe content rendering

- Render all workspace-provided names, descriptions, findings, paths, and artifact content through contextual `html/template` escaping. Do not convert any workspace string to `template.HTML`, inject it with `innerHTML`, or permit inline event handlers.
- Present Markdown preview as escaped source text in a labelled `<pre><code>` region. This deliberately favors auditability and hostile-content safety over rich Markdown rendering in the foundation sprint; embedded HTML and scripts remain inert text.
- Present JSON as escaped, formatted text only after bounded parsing succeeds. If parsing fails, show bounded escaped source plus a safe invalid-JSON notice. Never turn JSON values into markup or executable links automatically.
- Show preview size and truncation before the content. Long lines must wrap or scroll within the preview without widening the page, and copying source text must not require JavaScript.
- Static assets use fixed embedded paths and a restrictive Content Security Policy. JavaScript uses external embedded files with no inline script requirement and no third-party assets, fonts, analytics, or network calls.

### Visual and responsive system

- Use a restrained documentation-console visual language suited to local engineering inspection: high-contrast neutral surfaces, a readable system sans-serif for navigation and prose, and a system monospace stack for identifiers, state values, and artifacts. Status meaning always combines text, shape/border, and color.
- Use one responsive content grid. Wide screens may place navigation/context beside the main view and arrange summary panels in columns; narrow screens collapse to one reading column without horizontal page scrolling. Tables that cannot collapse receive a labelled horizontal scroll region, while key-value/status data should use definition lists or stacked rows.
- Keep content density appropriate for operational data: headings and summaries precede details, findings are grouped by severity/status, and long sections use native `<details>` only when the summary remains descriptive. Do not hide critical errors or the current state inside collapsed regions.
- Avoid animation by default. Any refresh indicator uses a text change rather than continuous motion; `prefers-reduced-motion` disables nonessential transitions.

### Accessibility contract

- Pages use valid heading order, semantic landmarks, lists, tables with captions/headers where tabular data is necessary, and a first-focusable skip link. Every page has a unique document title and visible `h1`.
- All navigation and disclosure behavior works with keyboard alone and native controls. Focus indicators remain visible. Progressive region replacement restores focus only when the focused element was replaced; otherwise it does not steal focus.
- Current navigation uses text plus `aria-current`; validation and health statuses have visible labels rather than color-only dots. Dates, counts, truncation, unknown values, and error severity are available as text.
- Refresh status uses a restrained `role="status"`/polite live region. Page-level and field/request errors use headings and direct explanatory text; assertive announcements are reserved for a refresh that invalidates the current view.
- Touch targets and links remain usable at narrow widths and zoom to 200 percent. The layout supports reflow without loss of information, except intentionally scrollable code/table regions.

### Performance and verification

- Embed one small stylesheet and at most one small dependency-free script. There is no Node.js dependency, framework bundle, hydration, client router, external font, or asset build pipeline.
- Bound initial content to the app/API collection and artifact limits. Do not render recursive trees or every artifact into hidden DOM. Link to detail pages instead of preloading them, and refresh no faster than an explicit user action in this sprint.
- Parse templates once during serve startup and fail startup on parse errors. Execute templates into a buffer before writing headers so rendering failures can produce a safe complete error response rather than partial HTML.
- Test every embedded template for parse success and representative dashboard, project, sprint, study, artifact, empty, truncation, and error rendering. Assertions must verify escaping of hostile names, Markdown, JSON, URLs, and error text.
- Add semantic tests for landmarks, title/heading structure, labels, table headers, `aria-current`, live-region behavior, keyboard-operable native controls, no-JavaScript navigation, responsive overflow classes, security headers, and absence of inline/third-party assets.
- Use fake app results for focused template tests and `httptest` with temporary workspaces for route-to-page tests. Normalize request IDs and timestamps only where a whole-page golden adds value; prefer semantic assertions for volatile status content.

The final frontend decision is to proceed with server-rendered, progressively enhanced, read-only pages. The main trade-off is less client-side interactivity in exchange for a smaller security boundary, complete no-JavaScript behavior, single-binary distribution, and direct alignment with authoritative workspace state.

## Trade-Offs

| Decision | Benefit | Cost / Rejected Alternative |
| --- | --- | --- |
| Server-rendered HTML as the primary experience | Delivers useful content without hydration or a frontend runtime and keeps request/state ownership explicit | A React/Vue-style client application was rejected because Sprint 30 does not justify a build pipeline, client router, duplicated state model, or larger security/dependency surface. |
| Progressive enhancement only | Preserves links, errors, and previews when JavaScript is absent or fails | A JavaScript-owned navigation/data flow could feel more app-like, but was rejected because it would make the foundation dependent on client execution. |
| Escaped Markdown source instead of rich rendering | Guarantees workspace HTML/scripts remain inert and keeps citations/content auditable | Rich Markdown rendering was rejected for now because adding a parser and sanitizer creates a second content-security contract without being required for inspection. |
| Live navigation/explicit refresh without client cache | Keeps workspace files and product run state authoritative and avoids stale cache semantics | Automatic polling and normalized client stores were rejected because they add duplicate fetches, synchronization, hidden background load, and browser-owned state. |
| Shared shell plus a few earned partials | Gives consistent navigation, status, empty, and error treatment without over-abstraction | Page-local duplication is acceptable until repetition is concrete; a generic component/design-system layer was rejected as premature for embedded Go templates. |
| Semantic HTML and native controls | Provides keyboard, focus, and assistive-technology behavior with minimal code | Custom interactive widgets were rejected because they require more JavaScript and accessibility state without adding Sprint 30 capability. |
| Textual operational density with responsive reflow | Preserves precise statuses and artifact metadata on desktop and mobile | A card-only dashboard was rejected because it obscures comparison and wastes space; unbounded wide tables were rejected because they fail narrow layouts. |
| Explicit user refresh only | Predictable load, no duplicate requests, and no surprise focus/live-region churn | Timed polling was rejected because refresh cadence, stale-state messaging, cancellation, and load bounds have not been justified. |
| Buffered template rendering | Prevents partial success responses and allows safe error projection | Direct streaming can reduce first-byte latency, but was rejected for these bounded pages because partial template failures are harder to recover and test. |

## Evidence

- The handbook's thin-transport evidence (`01-project-structure`, `02-command-architecture`, `03-dependency-injection`) supports templates and JSON as presentation adapters over typed app results. This grounds frontend-owned view models without product logic, CLI invocation, or raw product model exposure.
- The read-only capability and I/O-boundary evidence (`06-io-abstraction`, `13-security`) supports omitting all mutation/runtime capabilities from the browser surface. It also grounds app-issued artifact references and escaped content instead of direct path reads or UI-only enforcement.
- The error evidence (`05-error-handling`) supports preserving typed causes below the transport while rendering safe page/API messages. This is the basis for explicit empty, unavailable, not-found, and error states without exposing absolute paths or internal diagnostics.
- The state and cancellation evidence (`07-state-context`) keeps dependencies explicit and requests cancellable. The frontend therefore adds no global client service locator or detached refresh work; navigation and any enhanced refresh remain owned by the current request/view.
- The handbook favors sequential, bounded, incremental work (`08-concurrency`, `14-performance`) and warns against whole-workspace materialization. This supports bounded lists/previews, detail links rather than hidden preloads, no polling, and no speculative parallel browser fetches.
- Structured diagnostic separation (`10-logging-observability`) supports keeping lifecycle/request details out of visible pages while showing safe user-facing state. Browser output contains only presentation data and a request ID suitable for correlating a safe diagnostic.
- Deterministic fakes, real command paths, and normalized representations (`06-io-abstraction`, `11-testing-strategy`) support template tests with fake app results plus `httptest` route/security tests and selective normalized goldens.
- The extensibility evidence (`12-extensibility`) cautions against accidental public contracts and premature registries. A small embedded page/partial set is therefore preferred over a plugin UI, mutable component registry, or framework abstraction.
- The project documents require `html/template`, embedded CSS/minimal JavaScript, progressive enhancement, safe hostile-content handling, and no separate runtime. Sprint requirements narrow the broader Phase 4 vision to read-only inspection, so confirmations, forms that trigger work, operation handles, SSE, and cancellation UI are intentionally absent.

The evidence is high-confidence for transport separation, capability restriction, bounded content, diagnostic separation, and deterministic testing. Exact typography, spacing, and breakpoints are implementation details, but semantic reflow, visible focus, textual status, no-JavaScript operation, and hostile-content escaping are acceptance constraints.

## Risks

- **Template data leakage:** Adding a field to an app result may be rendered without reviewing disclosure. Mitigation: use explicit web view models, never pass broad models/maps to templates, and test that paths, secrets, environment values, raw errors, and provider payloads are absent.
- **Escaping bypass:** A convenience use of `template.HTML`, unsafe URL construction, or `innerHTML` could execute workspace content. Mitigation: prohibit trusted-markup conversion for workspace data, use app-issued references, keep scripts static, and test hostile HTML, scripts, URLs, Markdown, and JSON.
- **Status ambiguity:** Color-only badges or default zero values can misstate unknown, pending, failed, skipped, and inapplicable states. Mitigation: require textual labels, explicit optionality, and fixtures for every meaningful state.
- **Large-page degradation:** Dense workspaces could create slow rendering and unusable DOM size. Mitigation: enforce collection limits, expose truncation, link to detail views, avoid recursive trees/preloading, and measure before adding pagination or client virtualization.
- **Mobile overflow:** Operational tables and artifact lines can force viewport scrolling. Mitigation: prefer definition lists/stacked rows, wrap metadata, constrain preview regions, and give unavoidable tables/code labelled local scrolling.
- **Progressive-enhancement regression:** JavaScript may accidentally become required as refresh behavior grows. Mitigation: test navigation and all inspection scenarios with scripts disabled; every enhanced action must retain an ordinary link/reload fallback.
- **Stale refresh messaging:** A failed enhanced refresh may leave old content appearing current. Mitigation: visibly mark retained content stale, announce failure, and provide a full reload path; do not silently swallow fetch or envelope errors.
- **Accessibility drift in partials:** Reused status/finding partials can spread invalid heading, table, or ARIA patterns. Mitigation: keep partial contracts small, render them in representative pages, and combine semantic assertions with keyboard and zoom/reflow review.
- **Over-abstraction:** Treating Go templates as a frontend component framework could create indirection without reusable behavior. Mitigation: begin with page-owned templates and extract only repeated, semantically stable fragments.
- **Scope leakage:** Broader Phase 4 docs describe guarded operations and SSE. Mitigation: render no operational forms/buttons, register no operation client code, and test that Sprint 30 pages expose inspection and refresh only.
- **Source-preview expectations:** Users may expect rendered Markdown rather than escaped source. This is an accepted foundation limitation for security and auditability; documentation and preview labels must state the behavior. Rich rendering requires a later sanitizer/parser decision and hostile-content test contract.
- **Final reasoning handoff:** `projects/ultraplan-go/sprints/30-web-foundations/reasoning.md` must reference this document and carry forward its server-rendering, state, hostile-content, accessibility, responsive, and progressive-enhancement decisions. That artifact is outside this repair and is not modified here.

No frontend question blocks implementation. Rich Markdown rendering, automatic polling, client caching, a frontend framework, browser persistence, operation confirmations, SSE progress, cancellation controls, and other mutating UI remain deferred until a later sprint explicitly selects and governs them.
