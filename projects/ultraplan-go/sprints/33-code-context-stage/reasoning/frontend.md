> **Inputs Used:** `projects/ultraplan-go/sprints/33-code-context-stage/technical-handbook.md`, `projects/ultraplan-go/sprints/33-code-context-stage/requirements.md`, `projects/ultraplan-go/sprints/33-code-context-stage/sprint-index.md`, `projects/ultraplan-go/docs/ARCHITECTURE.md`, `projects/ultraplan-go/docs/PRD.md`, `projects/ultraplan-go/docs/TRD.md`, `system/reasoning/frontend-reasoning-template.md`, `studies/go-cli-study/reports/final/01-project-structure.md`, `studies/go-cli-study/reports/final/02-command-architecture.md`, `studies/go-cli-study/reports/final/04-configuration-management.md`, `studies/go-cli-study/reports/final/05-error-handling.md`, `studies/go-cli-study/reports/final/06-io-abstraction.md`, `studies/go-cli-study/reports/final/07-state-context.md`, `studies/go-cli-study/reports/final/08-concurrency.md`, `studies/go-cli-study/reports/final/10-logging-observability.md`, `studies/go-cli-study/reports/final/11-testing-strategy.md`, `studies/go-cli-study/reports/final/12-extensibility.md`, `studies/go-cli-study/reports/final/13-security.md`

# Frontend: Code-Context Sprint Presentation

This area covers the existing server-rendered sprint page and minimal progressive enhancement needed to present code-context readiness, findings, artifact access, execution progress, explicit rerun, cancellation, and recovery. It does not add a route, client-side application, browser-owned workflow state, or code-context-specific operation protocol.

## Area Decisions

### Ownership and placement

Extend the existing sprint page in `internal/web/templates/sprint.html` and its established typed view model. Handlers continue to map shared app results into presentation data before rendering. Templates do not read workspace files, call application services, validate requests, interpret product state, or start operations.

Reuse the existing page-to-layout-to-component-to-primitive hierarchy. Code-context is a sprint stage row/detail section composed from current domain-neutral presentation patterns: stage/status badge, readiness or validation summary, artifact-preview link/panel, operation status/progress, confirmation form, cancellation control, error panel, and next-action guidance. Introduce a new primitive or generic component only if no existing component can express the state without code-context-specific branching; keep the stage-specific composition at the sprint-page level.

No public export, frontend framework, JavaScript package, route family, or cross-feature store is needed. Existing embedded CSS tokens and layers remain authoritative; any styling change belongs to the existing component/layout layer and must work without JavaScript.

### User flow

The sprint page presents code-context in canonical order immediately after requirements and before sprint-index.

1. The operator sees its status, readiness, prerequisite explanation, latest operation outcome, artifact validity, and required next action.
2. When ready and not running, the operator can inspect the normalized operation scope, resolved runtime/model/variant source, affected path, and mutation class before requesting confirmation.
3. Starting or rerunning uses the existing generic prepare/start operation flow; rerun is explicit and clearly states that only `code-context.md` may be replaced after successful validation.
4. While active, the page shows bounded progress and an explicit cancel action. Refresh or SSE reconnect reloads operation and authoritative sprint state without changing ownership.
5. On completion, failure, cancellation, interruption, or cleanup uncertainty, the page refreshes durable status and presents findings, artifact availability, and a safe next action.

The no-JavaScript experience remains complete: server-rendered status, findings, artifact access, prepare/start forms, cancellation endpoint behavior, and post-action redirects are usable without SSE. JavaScript only enhances progress streaming and refresh behavior.

### State and data model

- **Server state:** shared app projections for sprint stage status, readiness, findings, artifact metadata, effective runtime selection, latest operation outcome, and next action.
- **URL state:** existing project/sprint route identity and existing safe artifact/operation references only.
- **Local UI state:** transient disclosure state, confirmation form state, and SSE connection status; none is authoritative.
- **Ephemeral server operation state:** operation ID, bounded recent safe events, subscribers, cancellation function, and terminal result under existing retention limits.
- **Durable state:** sprint `code-context.md` and sprint-owned flow state, never browser storage.
- **Derived presentation state:** badge labels, enabled/disabled actions, explanatory copy, and whether the artifact or rerun controls are shown, all computed by handlers from typed app results.

The view model must distinguish the latest operation outcome from the validity and availability of a preserved prior artifact. A failed rerun may therefore show a valid existing artifact and a failed latest attempt simultaneously; the UI must not collapse these into a success badge.

### State presentation

- **Missing/not ready:** show the blocking requirements finding and no start action.
- **Ready:** show the allowed read/write scope and an enabled prepare action.
- **Preparing/confirmation required:** show exact normalized scope, model/variant source, canonical output path, and confirmation expiry/staleness behavior.
- **Running:** show operation identity, current safe progress, start time where available, and explicit cancellation. Disable duplicate start/rerun actions.
- **Validation failure:** show actionable structural findings without raw provider output and retain a prior valid artifact link if one exists.
- **Runtime/persistence failure:** show the stable safe error summary and required next action; do not infer failure details from strings in the template.
- **Cancelled/interrupted/cleanup uncertain:** use distinct text and status semantics, with recovery guidance and no success styling.
- **Complete:** show validated artifact access, completion metadata, and an explicit rerun action rather than an implicit refresh.
- **Empty findings:** say that no structural findings are present; do not imply semantic completeness.

### Interaction and accessibility

Use semantic headings, lists, forms, buttons, links, tables, and status regions already established by the sprint page. The stage order in the DOM matches the canonical workflow order. Every action has a visible text label; status is conveyed by text and semantics, not color alone.

Confirmation and cancellation are keyboard-operable standard forms. Validation errors associate with their stage and use a focusable summary where the current pattern supports it. After a synchronous action or page reload, focus moves to the resulting status/error heading or is restored predictably. Enhanced progress uses a polite live region for meaningful state changes, not every raw event. Cancellation and terminal errors use assertive announcement only when the existing component semantics warrant it. Respect reduced-motion preferences and avoid progress animation as the sole indicator of activity.

### Security and bounded rendering

All output is escaped by `html/template`. Markdown artifact previews use the existing safe renderer and allowlist; source excerpts are never inserted into operation events or unbounded page data. Do not render raw HTML, full prompts, absolute implementation paths, environment values, raw provider payloads, unrestricted stderr, secrets, or caller-controlled artifact paths.

Mutating/runtime actions retain existing same-origin, Host, Origin, CSRF, request-body, short-lived confirmation, normalized-request fingerprint, and draining-server controls. Browser disconnect remains subscription loss only. The page must not offer a start/rerun action when the app reports a conflict, stale confirmation, unmet prerequisite, or draining server.

### Performance and testing

The bundle impact is effectively none: no framework or build pipeline is added, and existing minimal JavaScript handles generic operation/SSE behavior. The page performs no repository scan and does not embed the full artifact by default. Use bounded preview loading and existing operation-state refresh paths; do not poll in parallel with a healthy SSE stream unless the current generic implementation already does so.

Tests cover typed view-model mapping, canonical stage placement, every state listed above, no-JavaScript rendering, action enablement, preserved-artifact/latest-attempt distinction, hostile Markdown and escaped values, safe path allowlisting, redaction, confirmation staleness, CSRF/origin rejection, SSE ordering/reconnect/slow subscribers, browser-disconnect behavior, explicit cancellation, shutdown/recovery, and generic-operation contract reuse. Prefer semantic and structural assertions over brittle full-page snapshots or exact counts.

Final decision: proceed by extending the existing sprint page and generic components/operation enhancement. No frontend refactor is required; if the current page hardcodes a stage subset, make the smallest typed projection change needed to render the canonical stage list rather than adding a code-context-only branch throughout the template.

## Trade-Offs

| Decision | Benefit | Cost / Limitation |
| --- | --- | --- |
| Extend the existing sprint page | Preserves navigation, visual language, and shared state mapping | The page view model gains one more stage projection |
| Reuse generic status, finding, preview, and operation components | Keeps behavior consistent across stages and reduces duplicated interaction logic | Generic components must express code-context states without becoming stage-aware |
| Server rendering first, SSE enhancement second | Complete accessible fallback and durable-state truth | Live feedback is less app-like than a client-side state store |
| Separate artifact validity from latest operation outcome | Truthfully represents failed reruns that preserve prior output | Requires a slightly richer view model and clearer copy |
| Explicit rerun confirmation | Makes runtime cost and mutation scope visible | Adds an extra operator step |
| Bounded artifact preview | Protects responsiveness and disclosure boundaries | Operators may need to open the artifact directly for full content |
| Structural assertions over whole-page goldens | Tests durable behavior and accessibility semantics | Small visual regressions need focused component checks rather than one snapshot |

Rejected alternatives:

- **A dedicated code-context page or route:** rejected because the capability is one sprint stage and the existing sprint detail hierarchy already owns stage presentation.
- **A code-context-specific JavaScript controller or client store:** rejected because generic operation/SSE behavior already exists and browser state is not authoritative.
- **Automatic rerun when requirements or source change:** rejected because automatic staleness and amendment are out of scope and runtime work requires explicit intent.
- **Cancel on tab close or SSE disconnect:** rejected because subscriber lifecycle is not operation lifecycle.
- **Inline full artifact rendering in initial HTML:** rejected because previews must remain bounded and safe.
- **A new design-system primitive for every code-context state:** rejected because status, findings, progress, action, and error patterns are domain-neutral and already established.
- **Raw runtime logs as progress UI:** rejected because they are unstable, noisy, potentially sensitive, and unsuitable for accessible announcements.

## Evidence

### Report findings

- `studies/go-cli-study/reports/final/01-project-structure.md` finds that interface layers should remain thin over inward-owned behavior and that UI abstractions help support multiple projections (`cmd/gdu/app/app.go:30-49`, Helm's command/action split). This supports a presentation-only web change over shared app results.
- `studies/go-cli-study/reports/final/02-command-architecture.md` finds that shared lifecycle wrappers prevent command/interface handlers from accumulating business logic (`rclone/cmd/cmd.go:240-340`, `helm/pkg/cmd/install.go:132-145`). By inference, the browser should reuse generic operation lifecycle rather than implement stage-specific orchestration.
- `studies/go-cli-study/reports/final/04-configuration-management.md` finds that effective values must retain their source and distinguish explicit overrides from defaults (`restic/internal/global/global.go:139,147`, `go-task/internal/flags/flags.go:314-327`). This supports displaying model/variant source in confirmation without treating defaults as user choices.
- `studies/go-cli-study/reports/final/05-error-handling.md` finds that user-facing rendering should be separate from operational errors and based on typed identity (`age/cmd/age/tui.go:37-54`, `gh-cli/internal/ghcmd/cmd.go:281-301`). This supports safe actionable error panels rather than raw error strings.
- `studies/go-cli-study/reports/final/06-io-abstraction.md` finds that UI and stream abstraction enables deterministic capture and testing, while direct global output bypasses those seams (`gh-cli/pkg/iostreams/iostreams.go:551-568`, `go-task/executor_test.go:146-151`). This supports typed view models and render tests without live infrastructure.
- `studies/go-cli-study/reports/final/07-state-context.md` finds that operation/session state should be explicit only when needed and cancellation must follow the caller-owned context (`helm/pkg/cmd/install.go:333-347`, `opencode/internal/session/session.go:12-23`). This supports ephemeral operation identity, explicit cancellation, and no browser-owned durable session.
- `studies/go-cli-study/reports/final/08-concurrency.md` finds that asynchronous work needs localized ownership, bounded buffers, cancellation, and explicit waiting (`lazygit/pkg/gui/background.go:35,46,123`, `opencode/cmd/root.go:261-279`). This supports non-blocking bounded SSE subscribers and shutdown recovery rather than per-page goroutine ownership.
- `studies/go-cli-study/reports/final/10-logging-observability.md` finds that structured diagnostics must remain separate from result output and use consistent fields (`k9s/internal/slogs/keys.go:6-231`, `go-task/internal/output/output.go:12-14`). This supports concise safe progress events and stable status presentation.
- `studies/go-cli-study/reports/final/11-testing-strategy.md` finds that behavior-focused command/component tests, integration fixtures, and selective goldens outperform brittle implementation-detail assertions (`chezmoi/internal/cmd/main_test.go:64-174`, `helm/internal/test/test.go:43`, `k9s/internal/view/pod_test.go:23` as a caution). This supports semantic render and interaction tests.
- `studies/go-cli-study/reports/final/12-extensibility.md` finds that additive internal seams are preferable to plugin/registry machinery for a fixed capability (`go-task/executor.go:20-24,91-122`, `rclone/fs/rc/registry.go:41-48`). This supports reusing existing components rather than adding a frontend extension system.
- `studies/go-cli-study/reports/final/13-security.md` finds that centralized validation, explicit trust boundaries, path safety, and redaction are necessary around untrusted input (`k9s/internal/config/json/validator.go:146`, `restic/internal/options/secret_string.go:15-20`, `helm/pkg/registry/transport.go:37-41`). This supports escaped rendering, allowlisted previews, bounded diagnostics, and existing guarded-action controls.

### Sprint-specific inference

The project architecture already defines the browser as a server-rendered projection over typed app use cases, with generic guarded operations and SSE. The sprint requirements explicitly require code-context display and operation parity while forbidding route-specific semantics and durable web-owned state. Therefore the smallest honest frontend design is a new stage composition within the existing sprint page, not a new feature shell or client application.

## Risks

- **Truth collapse:** a preserved valid artifact could cause a failed rerun to look successful. Mitigation: separate artifact validity, stage state, and latest operation outcome in the view model and labels.
- **Stage-order drift:** a hardcoded template list could place code-context incorrectly. Mitigation: render or map from canonical typed stage order and assert DOM order.
- **Generic-component leakage:** adding code-context branches inside domain-neutral components could make them product-aware. Mitigation: handlers/page composition provide labels and state; primitives/components remain presentation-only.
- **Unsafe preview disclosure:** source excerpts or paths could bypass bounds or escaping. Mitigation: existing allowlist, contained artifact references, safe Markdown rendering, response limits, and hostile-content tests.
- **SSE authority confusion:** stale or dropped events could be treated as final truth. Mitigation: SSE is advisory progress; terminal/reconnect paths refresh authoritative app state.
- **Cancellation confusion:** disabling a button or closing a stream might be mistaken for cancellation. Mitigation: only explicit confirmed cancellation invokes the canonical app capability, and the UI waits for a truthful outcome.
- **Accessibility noise:** announcing every runtime event can overwhelm assistive technology. Mitigation: coalesce to meaningful status changes and keep full diagnostics outside live regions.
- **No-JavaScript regression:** progress enhancement could become required for operation completion or recovery. Mitigation: test complete server-rendered forms, redirects, status, findings, and artifact access without scripts.
- **Open verification question:** implementation must confirm whether the existing sprint view model can represent a valid artifact alongside a failed latest attempt. If not, add explicit fields at the handler/app projection boundary rather than deriving this distinction in the template.
