> **Inputs Used:** `projects/ultraplan-go/sprints/33-code-context-stage/technical-handbook.md`, `projects/ultraplan-go/sprints/33-code-context-stage/requirements.md`, `projects/ultraplan-go/sprints/33-code-context-stage/sprint-index.md`, `projects/ultraplan-go/docs/ARCHITECTURE.md`, `projects/ultraplan-go/docs/PRD.md`, `projects/ultraplan-go/docs/TRD.md`, `system/reasoning/api-design-reasoning-template.md`, `studies/go-cli-study/reports/final/02-command-architecture.md`, `studies/go-cli-study/reports/final/04-configuration-management.md`, `studies/go-cli-study/reports/final/05-error-handling.md`, `studies/go-cli-study/reports/final/07-state-context.md`, `studies/go-cli-study/reports/final/10-logging-observability.md`, `studies/go-cli-study/reports/final/11-testing-strategy.md`, `studies/go-cli-study/reports/final/12-extensibility.md`, `studies/go-cli-study/reports/final/13-security.md`

# API Design: Code-Context Stage

This area covers the internal application API, CLI and stable JSON projections, and existing local HTTP operation API needed to expose `code-context`. It does not define a public network service or a new workflow protocol. The audience is UltraPlan's CLI, TUI, and loopback browser adapters, all calling typed shared application capabilities.

## Area Decisions

### Contract shape

Add `code-context` as a value in the existing sprint-stage contract, immediately after `requirements`, rather than as a top-level command, route family, operation kind, or plugin. The same stage value must be accepted by:

- `ultraplan sprint <project> <sprint> prompt code-context`
- `ultraplan sprint <project> <sprint> validate code-context`
- `ultraplan sprint <project> <sprint> flow --to code-context`
- sprint status and stable JSON stage projections
- stage-specific effective model and variant configuration projections
- typed app requests used by CLI and web adapters

This is an additive internal/CLI/local-HTTP contract change. `internal/sprint` remains the semantic owner; `internal/app` maps typed requests and results; CLI and `internal/web` only parse, validate, invoke, and render. This follows the requirements' shared-boundary constraints and avoids making CLI text an integration protocol.

### Synchronous and asynchronous operations

Keep prompt preview, structural validation, status, readiness, findings retrieval, effective configuration display, and bounded artifact preview synchronous. They must not construct or invoke the runtime, write `code-context.md`, or mutate flow state merely because they were queried.

Run generation and explicit rerun through the existing asynchronous operation capability. Reuse the established sequence:

1. Prepare a normalized stage-operation request and return a short-lived confirmation bound to the request and current governed-input fingerprint.
2. Start the confirmed operation and return the existing ephemeral operation identity.
3. Observe bounded safe progress through existing operation status and SSE event capabilities.
4. Request cancellation through the existing cancellation capability.
5. Refresh authoritative sprint status, findings, and artifact state after a terminal outcome or reconnect.

No new `code-context` HTTP route or durable web session model is permitted. The existing generic local endpoints remain the transport surface, including operation prepare/start/status/events/cancel and allowlisted artifact preview. The browser may retain only bounded ephemeral operation state; filesystem artifacts and sprint-owned state remain authoritative.

### Request contract

The typed stage request must carry only the data needed to identify and control the operation:

- project and sprint references, resolved and contained by the application boundary
- stage fixed to `code-context`
- mode distinguishing preview/validate/run and explicit rerun where the existing API requires it
- optional explicit model and variant overrides, with whether each override was supplied preserved
- dry-run where supported
- current confirmation token and governed-input fingerprint for browser-started runtime work

Adapters must reject unknown stages, missing or ambiguous project/sprint references, stale or mismatched confirmations, invalid model/variant values, attempts to inject an output path or implementation target, and incompatible option combinations. The output path is always the canonical sprint-root `code-context.md`; target resolution is product-owned and cannot be caller-selected through a raw path field.

Prompt preview and dry-run return a deterministic safe projection of inputs and effective runtime selection, not the full prompt in status, logs, operation events, or browser operation metadata. Existing HTTP body, stream, event-buffer, diagnostic, and artifact-preview bounds apply; this sprint must not create a separate payload-size or concurrency regime for `code-context`.

### Response and error contract

Shared app results are typed and transport-neutral. Surface projections should expose the existing safe subset needed by operators:

- project, sprint, stage, readiness, execution status, and authoritative artifact path
- validation state and actionable findings
- operation identity and bounded progress for an active web operation
- terminal outcome and required next action
- effective runtime/model/variant and their sources where already supported
- stable attempt, prompt identity/version/checksum, and correlation metadata where available

Do not expose internal filesystem models, absolute implementation paths, full prompts, source excerpts in operation events, environment values, raw provider payloads, unrestricted stderr, or unbounded diagnostics. Artifact content is available only through the existing contained and allowlisted bounded-preview capability.

Errors must preserve wrapped causes internally and map to existing stable classifications at adapters. The material distinctions are usage/reference, prerequisite/validation, stale confirmation, operation conflict, configuration, runtime, missing output, invalid output, persistence, cancellation, interruption, and cleanup uncertainty. Add a new error class only if no existing stable class can represent one of these outcomes. Human text includes a safe next action; JSON and web code branch on stable codes or typed identity, never error-string parsing. Runtime exit success with absent or invalid output maps to a validation or missing-output failure, not success.

### Authorization and trust boundary

The CLI uses the local process identity and workspace access model. The browser remains loopback-only and uses the existing same-origin, Host, Origin, CSRF, body-limit, session, and confirmation controls; Sprint 33 adds no hosted authentication, users, tenants, or role model.

Authorization is capability- and path-based: the operation may read the resolved implementation repository and governed sprint inputs, but its only stage output is the canonical sprint `code-context.md`. API callers cannot broaden that scope. Artifact preview is allowlisted to this sprint-owned artifact, paths shown in generated content must be repository-relative and contained, and every public projection applies existing redaction.

### Idempotency, conflict, and cancellation

Prompt, validate, status, readiness, and preview operations are naturally retryable reads. Cancellation requests are idempotent through the existing canonical cancellation function. A browser disconnect cancels only its subscription and never the operation.

Generation is not automatically retried by the transport because it invokes a costly runtime and writes authoritative state. The existing operation conflict/mutation-lock behavior prevents concurrent stage mutation. An explicit rerun is the only duplicate-write behavior: it generates and validates an isolated candidate and atomically replaces only `code-context.md`; failure or cancellation preserves the last valid artifact and records the truthful operation outcome. A repeated start with an expired, stale, consumed, or mismatched confirmation is rejected rather than deduplicated heuristically.

### Compatibility

The CLI change is additive: existing sprint command shapes remain intact and gain one accepted stage value. Stable JSON retains existing envelopes and meanings; `code-context` appears through existing stage, artifact, finding, operation, and next-action fields rather than a parallel response schema. New optional fields may be added only under the project's existing JSON compatibility rules.

Persisted pre-code-context flow state must load deterministically without losing prior outcomes. Read-only status must not perform a hidden migration write. Compatibility reconciliation must project truthful current readiness and stage state, while any durable rewrite occurs only through an explicit mutating path. This sprint does not add API version negotiation or a second state format.

### Observability and cost safety

Reuse existing command/operation events and stable structured fields. Safe correlation includes command/operation ID, project, sprint, stage, attempt, status, duration, runtime, model, variant, configuration source, prompt identity/version/checksum, validation result, cancellation reason, and recovery outcome where available. Unknown usage or cost remains unknown rather than zero.

Logs and events are diagnostics, not the result payload. Structured JSON output must remain parseable, so diagnostics do not go to stdout. Full prompts, source content, unsafe paths, secrets, raw provider events, and unrestricted stderr are excluded. Event buffers, SSE subscribers, active operations, and cleanup waits remain bounded. There is no stage-internal parallel fan-out or automatic transport retry in this sprint; cost and cancellation remain attributable to one canonical operation.

### Testing contract

Use layered tests according to the boundary under proof:

- table-driven app and command tests for accepted stage values, help, prompt preview, validation, dry-run non-mutation, explicit override tracking, exit/error mapping, and stable JSON
- compatibility fixtures for old flow state, including proof that status does not write migration state
- fake-runtime tests for success, runtime failure, missing/invalid output, cancellation, metadata projection, and preserved last-valid artifact on failed rerun
- temporary repository/filesystem tests for source-read/output-write isolation, path containment, atomic replacement, and conflict behavior
- `httptest` and operation-contract tests for generic prepare/start/status/events/cancel reuse, stale confirmation, bounded safe events, disconnect behavior, shutdown cancellation, recovery, and absence of a route-specific operation kind
- golden tests only for stable help, embedded defaults, and stable JSON/Markdown shapes; use structural assertions for timestamps, operation IDs, paths, provider details, and progress

Compatibility tests must compare the same typed state across CLI/app/web projections. Normal tests remain offline and deterministic; real-runtime proof is explicitly deferred to Sprint 34.

## Trade-Offs

| Decision | Benefit | Cost / Limitation |
| --- | --- | --- |
| Extend the existing stage-valued contracts | Preserves command discoverability, typed app reuse, and one workflow vocabulary | Every canonical stage projection must be updated consistently and tested for drift |
| Reuse generic operation APIs | Preserves cancellation, confirmation, progress, conflict, and recovery semantics already proven by the web surface | The generic request/result types must remain expressive enough without leaking stage semantics into `internal/web` |
| Keep reads synchronous and generation asynchronous | Makes preview/status cheap and deterministic while retaining observable long-running work | Clients must refresh authoritative state after operation completion or reconnect |
| Use stable error classes plus wrapped causes | Gives CLI/JSON/web reliable branching and safe recovery guidance | Over-classification would expand the compatibility surface, so one-off details remain wrapped causes |
| Preserve explicit override-source metadata | Makes model/variant selection truthful and prevents defaults from masquerading as user intent | Parsing and request types must retain whether a flag was supplied, not only its value |
| Reject caller-controlled target/output paths | Enforces the single-artifact write boundary and avoids path injection | The API is intentionally less flexible than a generic runtime endpoint |
| Atomic explicit rerun without transport retry | Preserves the last valid artifact and avoids duplicate cost | A caller must deliberately start another run after a terminal failure |
| Additive legacy-state interpretation without read-time writes | Keeps old workspaces usable and status side-effect free | Compatibility logic remains necessary until an explicit migration policy removes it |
| Safe bounded projections instead of raw runtime detail | Protects secrets, source content, JSON stability, and browser safety | Deep provider diagnosis remains in explicitly safe runtime diagnostics rather than normal API responses |

Rejected alternatives:

- **A dedicated `/api/v1/code-context` route family:** rejected because it would duplicate generic operation, confirmation, cancellation, status, and recovery behavior and move product semantics into the web adapter.
- **A new durable code-context operation/session record:** rejected because sprint artifacts and flow state already own durable truth, while the web operation hub is intentionally ephemeral.
- **A new top-level CLI command:** rejected because `code-context` is one canonical sprint stage and must compose with existing prompt, validate, flow, status, and configuration commands.
- **A dynamic stage or plugin registry:** rejected because the stage set is compile-time product workflow, and registry/plugin collision, lifecycle, and versioning costs do not solve a Sprint 33 requirement.
- **Caller-supplied implementation or output paths:** rejected because they weaken containment, fingerprinting, and the single-output guarantee.
- **Automatic request retries or idempotency keys for generation:** rejected because runtime cost, explicit rerun intent, atomic candidate handling, and existing mutation exclusion provide clearer semantics.
- **Raw runtime/provider payloads in JSON or SSE:** rejected because they are unstable and may contain prompts, source content, paths, environment data, or secrets.

## Evidence

### Report findings

- `studies/go-cli-study/reports/final/02-command-architecture.md` finds that high-scoring CLIs use thin command factories and delegate behavior to reusable action/application layers (`gh-cli/pkg/cmdutil/factory.go:16-43`, `helm/pkg/cmd/install.go:132-145`). Its `rclone/cmd/cmd.go:240-340` evidence also supports one shared lifecycle wrapper. This supports extending existing sprint commands and shared operation orchestration rather than creating stage-specific handlers.
- `studies/go-cli-study/reports/final/04-configuration-management.md` finds that reliable precedence requires explicit ordering, tracking whether a flag changed, and validating after merge (`go-task/internal/flags/flags.go:314-327`, `restic/internal/global/global.go:139,147`, `k9s/internal/config/k9s.go:423-451`). This directly supports preserving supplied/not-supplied state for code-context model and variant overrides and projecting the effective source.
- `studies/go-cli-study/reports/final/05-error-handling.md` finds that wrapped causes, sentinel or typed identity, adapter-specific rendering, and exit mapping enable reliable recovery (`gh-cli/pkg/ssh/ssh_keys.go:64`, `helm/pkg/storage/driver/driver.go:27-48`, `gh-cli/internal/ghcmd/cmd.go:44-49,281-301`). This supports stable safe error classes shared across CLI and web without exposing raw causes.
- `studies/go-cli-study/reports/final/07-state-context.md` finds that cancellation is credible only when a caller-owned context reaches actual work, and that most CLIs do not need an extra session abstraction (`helm/pkg/cmd/install.go:333-347`, `opencode/internal/session/session.go:12-23`). Its cleanup evidence (`restic/internal/restic/lock.go:290-305`) supports separately representing bounded cleanup and uncertainty. This supports canonical operation cancellation, disconnect-as-subscription-loss, and no new durable session model.
- `studies/go-cli-study/reports/final/10-logging-observability.md` finds that structured fields plus strict result/diagnostic stream separation preserve scriptability (`helm/internal/logging/logging.go:31-71`, `k9s/internal/slogs/keys.go:6-231`, `go-task/internal/output/output.go:12-14`). This supports bounded safe operation metadata and keeping diagnostics out of JSON result stdout.
- `studies/go-cli-study/reports/final/11-testing-strategy.md` finds that confidence comes from layered command, fake, fixture, integration, and selective golden tests rather than one style (`chezmoi/internal/cmd/main_test.go:64-174`, `helm/internal/test/test.go:43`, `restic/cmd/restic/integration_helpers_test.go:188-235`). This supports separate contract tests for CLI/JSON, app behavior, real filesystem isolation, and generic web operations.
- `studies/go-cli-study/reports/final/12-extensibility.md` finds that factories and additive internal options are lower-cost extension seams for fixed workflows, while registries/plugins add collision, lifecycle, and versioning concerns (`go-task/executor.go:20-24,91-122`, `gh-cli/pkg/cmdutil/factory.go:16-43`, `rclone/fs/rc/registry.go:41-48`). This supports explicit canonical-stage registration and rejects a dynamic stage registry.
- `studies/go-cli-study/reports/final/13-security.md` finds that explicit trust boundaries, centralized validation, argument-safe execution, path controls, and redaction distinguish robust CLIs (`k9s/internal/config/json/validator.go:146`, `lazygit/cmd_obj_builder.go:38`, `restic/internal/options/secret_string.go:15-20`, `helm/pkg/registry/transport.go:37-41`). This supports fixed target/output resolution, contained artifact preview, central request validation, and omission of unsafe diagnostics.

### Sprint-specific inference

The technical handbook synthesizes the same evidence into pressures for thin adapters, source-aware configuration, typed error projection, caller-owned cancellation, bounded observability, and explicit fixed-stage registration. Applying those findings to the sprint requirements leads to one conclusion: **proceed by extending the existing stage and generic operation contracts; do not split out a code-context-specific API.**

This inference is also constrained by `projects/ultraplan-go/docs/ARCHITECTURE.md`, which assigns stage semantics to `internal/sprint`, shared use cases and projections to `internal/app`, and only HTTP mapping/SSE state to `internal/web`; by `projects/ultraplan-go/docs/TRD.md`, which compatibility-controls `/api/v1` and defines generic operation endpoints; and by the sprint acceptance criteria requiring shared CLI/app/web behavior with no route-specific operation kind.

## Risks

- **Canonical-stage drift:** stage order, help, parsing, JSON, configuration, app results, artifact allowlists, and browser presentation can diverge if each keeps a separate list. Mitigation: one explicit canonical definition where the current design permits it, plus parity tests across every projection; do not solve drift with a dynamic registry.
- **Stable JSON regression:** adding fields or a stage can accidentally rename statuses, alter nullability, or leak internal models. Mitigation: additive fields only, compatibility fixtures, structural JSON assertions, and safe DTO mapping at adapters.
- **False idempotency:** a client may retry a start after a timeout and create duplicate costly work. Mitigation: confirmation consumption/staleness checks and sprint mutation conflict handling; do not advertise generation as automatically retryable.
- **Cancellation race:** completion, cancellation, runtime failure, and cleanup may race. Mitigation: one canonical operation context, idempotent cancellation, single terminal-outcome arbitration, bounded cleanup, and explicit `interrupted` or `cleanup_uncertain` outcomes rather than inferred success.
- **Last-valid artifact ambiguity:** preserving a previous artifact after a failed rerun can mislead a client into treating the rerun as successful. Mitigation: expose artifact validity separately from latest operation outcome and required next action.
- **Diagnostic disclosure:** source excerpts, prompts, implementation paths, provider payloads, or stderr may enter logs/events/errors. Mitigation: projection-time allowlists, bounded excerpts only in the artifact-preview capability, redaction tests, and no raw payload projection.
- **Compatibility hidden mutation:** eager state migration during status would violate read-only API semantics. Mitigation: deterministic in-memory compatibility interpretation during reads and explicit atomic persistence only on mutating operations.
- **Overstated validation:** structural validation can prove shape, containment, and excerpt presence but not semantic completeness. Mitigation: label findings accurately and never project a structurally valid pack as complete repository coverage.
- **Open verification question:** implementation review must confirm that the existing generic operation request can identify a sprint stage without adding a route-specific operation kind and that its current conflict key excludes concurrent mutation of the same sprint. If either capability is absent, extend the generic typed app contract narrowly; do not add a code-context transport workflow.
- **Open compatibility question:** implementation review must identify the currently documented stable JSON nullability and unknown-stage behavior before adding fixtures. The decision remains to preserve the existing envelope and semantics rather than introduce a new schema.
