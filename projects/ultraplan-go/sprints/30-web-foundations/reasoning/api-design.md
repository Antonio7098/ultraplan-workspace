> **Inputs Used:** `projects/ultraplan-go/sprints/30-web-foundations/sprint-index.md`, `projects/ultraplan-go/sprints/30-web-foundations/technical-handbook.md`, `projects/ultraplan-go/sprints/30-web-foundations/requirements.md`, `projects/ultraplan-go/docs/ARCHITECTURE.md`, `projects/ultraplan-go/docs/PRD.md`, `projects/ultraplan-go/docs/TRD.md`, `system/reasoning/api-design-reasoning-template.md`

# API Design: Read-Only Local Web Foundation

This area covers Sprint 30's server-rendered HTML routes and its frontend-owned, loopback-only `/api/v1` JSON API. The API lets the bundled browser inspect authoritative workspace, project, sprint, study, validation, run/flow-state, health, and allowlisted artifact-preview data. It is not a public integration API and grants no mutation, runtime, review, smoke, cancellation, or Git capability.

## Area Decisions

### Audience and operation model

- Treat `/api/v1` as a compatibility-controlled internal API for the browser shipped in the same `ultraplan` binary. Documentation may show its behavior for troubleshooting, but Sprint 30 does not promise a general public, partner, webhook, or SDK surface.
- Use synchronous request/response for every Sprint 30 route. All API operations are naturally read-only and idempotent; no jobs, operation handles, confirmations, SSE, WebSockets, or cancellation routes are created.
- Read authoritative state on each request through typed `internal/app` query use cases. Do not add a web cache or browser-owned snapshot. A response represents a point-in-time projection and makes no cross-request snapshot-consistency guarantee.
- Keep HTML and JSON as separate representations over the same typed app results. HTML routes render with `html/template`; JSON handlers map app results into transport-owned DTOs. Neither representation calls the other or parses CLI output.

### Route contract

The accepted methods and routes are:

| Method | Route | Contract |
| --- | --- | --- |
| `GET` | `/` | Server-rendered workspace dashboard. |
| `GET` | `/projects` | Server-rendered project list. |
| `GET` | `/projects/{project}` | Server-rendered project detail and sprint summary. |
| `GET` | `/projects/{project}/sprints/{sprint}` | Server-rendered sprint, flow, validation, review, and smoke status. |
| `GET` | `/studies` | Server-rendered study list. |
| `GET` | `/studies/{study}` | Server-rendered study, run-state, and validation detail. |
| `GET` | `/artifacts/{ref}` | Server-rendered bounded Markdown or JSON preview for an app-issued artifact reference. |
| `GET` | `/api/v1/dashboard` | Workspace overview plus bounded project, sprint, study, validation, and current-state summaries. |
| `GET` | `/api/v1/projects` | Bounded project summaries. |
| `GET` | `/api/v1/projects/{project}` | One project and its bounded sprint summaries. |
| `GET` | `/api/v1/projects/{project}/sprints/{sprint}` | One sprint's planning, execute, review, smoke, flow, and validation projection. |
| `GET` | `/api/v1/studies` | Bounded study summaries. |
| `GET` | `/api/v1/studies/{study}` | One study's status, run-state, and validation projection. |
| `GET` | `/api/v1/validations?scope={workspace|project|sprint|study}&ref={opaque-ref}` | Existing validation findings only; it never starts validation. |
| `GET` | `/api/v1/artifacts/{ref}` | Bounded allowlisted artifact metadata and source text. |
| `GET` | `/api/v1/health` | Cheap server readiness and workspace-query availability. |

`{project}`, `{sprint}`, and `{study}` are canonical identifiers, not filesystem paths. `{ref}` is an opaque, URL-safe reference issued by an app query result; clients cannot submit absolute or workspace-relative paths. Detail links and artifact references therefore originate from a successful listing/detail response instead of turning HTTP into a filesystem API.

Only `GET` and implicit `HEAD` for successful `GET` resources are accepted. `HEAD` returns the same status and headers without a body. Other methods return `405` with `Allow: GET, HEAD`. Unknown paths below `/api/`, including unsupported versions, always return the JSON error envelope; unknown browser paths render the HTML error page.

The requirements' “validation” scope is resolved as inspection of existing validation summaries. `POST /api/v1/validations`, operation preparation/start/status, SSE, and `DELETE` cancellation from the broader Phase 4 TRD are deliberately deferred because Sprint 30 prohibits browser-triggered validation and workflow operations.

### Request and resource limits

- Reject request bodies on Sprint 30 API routes. A non-empty body returns `400 invalid_request`; a declared body larger than 64 KiB returns `413 request_too_large` before reading it fully.
- Limit the request target to 8 KiB, each decoded identifier or opaque reference to 128 bytes, and accepted query parameters to the documented names. Reject duplicate, malformed, control-character-bearing, slash-bearing identifier values and unknown query parameters with `400 invalid_request`.
- Return at most 200 entries per collection or findings list. Responses report `returned_count`, `total_count`, and `truncated`; Sprint 30 does not add caller-controlled pagination.
- Limit an artifact preview to the first 256 KiB after containment and type checks. Return `size_bytes`, `returned_bytes`, and `truncated`. Read only enough bytes to establish the bound; do not load or recursively scan the workspace to produce a preview.
- Use the server's bounded request concurrency and HTTP read/write/idle timeouts for all routes. App calls receive the request context so client cancellation and server shutdown stop pending reads.

These fixed foundation limits are intentionally conservative and can change only with a documented `/api/v1` compatibility assessment. They are not environment- or query-controlled in this sprint.

### Success contract

Every JSON success response has this top-level shape:

```json
{
  "data": {},
  "meta": {
    "api_version": "v1",
    "request_id": "server-generated-id",
    "generated_at": "RFC3339 timestamp"
  }
}
```

Collection metadata also contains `returned_count`, `total_count`, and `truncated`. Domain absence is represented deliberately: optional timestamps and optional terminal results are omitted when unknown; an empty known collection is `[]`, not `null`; missing numeric measurements remain absent rather than becoming zero. Transport DTOs use stable JSON field names and workspace-relative display references. Raw product models, absolute paths, raw provider events, unrestricted stderr, internal error causes, secrets, and environment data are never exposed.

Artifact success data contains `ref`, a safe workspace-relative display path, `media_type` (`text/markdown` or `application/json`), `content` as plain source text, `size_bytes`, `returned_bytes`, and `truncated`. JSON content may additionally include a parsed value only when it fits the same byte bound and parses successfully. Markdown embedded HTML is never marked trusted: the HTML page escapes or sanitizes it, and the JSON endpoint serves it as data with `X-Content-Type-Options: nosniff`.

Health data contains `status` (`ok` or `unavailable`), a cheap `server` check, and a lightweight `workspace` check. It returns `200` only when the server can answer queries for the configured workspace and `503` otherwise. It does not scan every project/study, check runtime/provider health, invoke smoke, or disclose paths and configuration. Health is a truthful readiness projection, not proof that every artifact is valid.

All successful state and preview responses use `Cache-Control: no-store` so browser refresh reads current authoritative workspace state.

### Error contract

Every `/api/` failure uses `application/json` and this envelope:

```json
{
  "error": {
    "code": "not_found",
    "message": "The requested resource was not found."
  },
  "meta": {
    "api_version": "v1",
    "request_id": "server-generated-id"
  }
}
```

`code` is machine-readable and compatibility-controlled; `message` is safe and actionable but must not include internal causes or absolute paths. A bounded `details` object is allowed only for safe field-level validation facts. The HTTP mapping is:

| Status | Code | Meaning |
| --- | --- | --- |
| `400` | `invalid_request` | Malformed identifier, reference, query, or unexpected body. |
| `403` | `request_rejected` | Host or Origin policy rejected the request without disclosing policy internals. |
| `404` | `not_found` | Unknown route/resource or invalid, stale, unsupported, escaping, or non-allowlisted artifact reference. |
| `405` | `method_not_allowed` | Route exists but method is unsupported. |
| `413` | `request_too_large` | Request limit exceeded. |
| `500` | `internal_error` | Unclassified server/app failure with details retained only in diagnostics. |
| `503` | `unavailable` | Configured workspace or required read dependency is temporarily unavailable. |

Handlers classify typed app errors with `errors.Is`/`errors.As`-style identity and then project this safe set. They never parse error strings or serialize internal errors. Artifact rejections deliberately collapse to `404` to avoid revealing whether an arbitrary local path exists.

### Security boundary

- There is no account authentication or tenant authorization in Sprint 30. Access is constrained by loopback-only binding, strict Host validation, same-origin pages/API, no permissive CORS, and Origin validation whenever the header is present. The server still treats every URL, query, header, and workspace artifact as hostile input.
- An absent Origin is allowed for normal top-level navigation and local non-browser `GET` clients after Host validation. A present Origin must exactly match the effective loopback server origin. `Origin: null`, malformed origins, non-loopback origins, and cross-origin requests are rejected.
- All routes receive only query capabilities. HTTP composition does not inject mutating repositories, runtime/process adapters, CLI handlers, Git operations, or arbitrary filesystem readers. This enforces read-only behavior structurally rather than through hidden buttons.
- Security headers apply to HTML and JSON, including a restrictive Content Security Policy for embedded assets, `X-Content-Type-Options: nosniff`, `Referrer-Policy: no-referrer`, and frame denial. No cookies or browser session state are needed for this read-only sprint.

### Observability and tests

Generate a request ID at the server boundary rather than trusting a caller-supplied value. Structured request diagnostics record request ID, normalized route name, method, status, duration, response byte count, and safe error code. Lifecycle diagnostics record listen host without workspace secrets, startup result, shutdown cause, shutdown duration, and cleanup failures. Do not log raw URLs, query values, artifact content, absolute workspace paths, Host/Origin values, or internal errors at normal level; redacted debug diagnostics may retain classified causes.

Required tests are:

- Table-driven route/method tests for every HTML and JSON route, `HEAD`, `405`/`Allow`, unsupported versions, and unknown `/api/` JSON fallback.
- Semantic schema tests for top-level envelopes, stable field names, null/omission behavior, collection truncation, content types, cache policy, and health `200`/`503` truthfulness.
- Typed app-error mapping tests for every public status/code and assertions that causes, paths, environment values, and secrets are absent.
- Host/Origin/no-CORS, absent-Origin, security-header, request-target/body/query/identifier limit, and request-cancellation tests.
- Artifact tests for valid Markdown/JSON, opaque-reference round trips, stale/forged refs, traversal and encoded traversal, symlink escape, unsupported extension, exact/over-limit size, invalid JSON, hostile Markdown/HTML/script, and `nosniff` behavior.
- Fake-use-case handler tests plus `httptest` tests over temporary workspaces to prove handlers use typed app queries, state agrees with existing CLI/TUI projections, and repeated reads cause no mutation.
- Normalized JSON golden fixtures only for the compatibility-sensitive success/error envelopes; use semantic assertions for timestamps, request IDs, paths, counts, and other volatile fields.

## Trade-Offs

| Decision | Benefit | Cost / Rejected Alternative |
| --- | --- | --- |
| Versioned resource-oriented `GET` routes | Clear read-only semantics, natural retry safety, and a stable browser boundary | A single RPC-style `/api/query` was rejected because operation names and union payloads obscure cache, method, and compatibility behavior. |
| Frontend-owned but compatibility-controlled `/api/v1` | Lets bundled HTML/minimal JavaScript evolve with a documented stable envelope | Calling it a public API was rejected because authentication, remote access, broad client guarantees, pagination, and release policy are outside scope. |
| One envelope with route-specific DTOs | Stable error handling and metadata without exposing internal models | Serializing app/domain structs directly was rejected because product refactors would become transport breaks and could disclose unsafe fields. |
| Opaque app-issued artifact references | Prevents the HTTP layer from becoming an arbitrary path reader and centralizes allowlisting | Caller-provided workspace-relative paths were rejected even with `filepath.Clean`; path strings increase traversal, symlink, encoding, and disclosure risk. |
| Live per-request reads with `no-store` | Preserves workspace files and product run state as the source of truth and avoids invalidation state | Server snapshots/caches were rejected for the foundation because they add staleness, synchronization, memory, and restart semantics. Repeated read cost remains. |
| Sequential app aggregation by default | Deterministic ordering, straightforward cancellation, and simple complete-error semantics | Handler-owned fan-out and partial panel responses were rejected until measurement demonstrates a need; they add goroutine ownership and partial-result contracts. |
| Fixed collection and preview limits without pagination | Prevents accidental large scans and keeps the initial client simple | Unbounded responses were rejected for safety. Cursor pagination was also rejected because no Sprint 30 scenario requires arbitrary-depth browsing. |
| No authentication token or cookie | Matches a single-process loopback read-only tool and avoids inventing session state | Treating loopback as fully trusted was rejected; Host/Origin, containment, disclosure controls, and no-CORS remain mandatory. Hosted authentication is explicitly out of scope. |
| Validation summaries via `GET` only | Meets inspection requirements without acquiring execution capability | `POST /validations` was rejected for Sprint 30 because it starts product behavior and belongs with guarded operations in Sprint 31. |
| Cheap workspace-aware health | Truthfully distinguishes a listening process from a usable configured workspace | Deep health that scans artifacts or checks runtime/provider availability was rejected because it is expensive, can leak detail, and conflates server readiness with product diagnostics. |

## Evidence

- The handbook's thin-transport pattern cites gdu's multiple presentation adapters, Helm's command-to-action delegation, and restic's repository boundary (`01-project-structure`, `02-command-architecture`, `03-dependency-injection`). This supports HTML and JSON mapping the same typed app read results rather than invoking CLI handlers or owning product logic.
- The handbook's capability-boundary finding (`06-io-abstraction`, `13-security`) shows that dangerous capabilities should be absent from injected interfaces. The API therefore receives query use cases and app-issued artifact references, not mutable services, command runners, or broad filesystem access.
- The boundary-specific error evidence (`05-error-handling`) shows typed lower-layer errors and context-specific rendering. This is the basis for retaining internal causes in diagnostics while exposing a small stable HTTP status/code set.
- The state/context evidence (`07-state-context`) requires cancellation to propagate into blocking work without using context as a service locator. Each handler passes the request context to explicit app dependencies, and server shutdown owns cancellation.
- The concurrency and performance evidence (`08-concurrency`, `14-performance`) favors a sequential baseline and bounded, incremental reads. This supports fixed list/preview limits, no whole-workspace materialization, and no speculative handler fan-out.
- The observability evidence (`10-logging-observability`) separates structured diagnostics from user representations and uses stable fields. Request IDs and normalized route/status fields therefore stay in logs and response metadata without placing raw internal diagnostics in JSON or HTML.
- The testing evidence (`06-io-abstraction`, `11-testing-strategy`) combines deterministic fakes, real command/server paths, and normalized golden fixtures. The test decision combines fake app handlers, `httptest`, temporary-workspace integration, and narrow envelope goldens.
- The extensibility evidence (`12-extensibility`) warns that documented versioned APIs acquire compatibility obligations and mutable registries can become accidental public contracts. `/api/v1` therefore has explicit DTO/envelope ownership but no plugin registry or claim of public extensibility.
- The project architecture requires `internal/web -> internal/app`, transport-owned DTOs, and no direct web imports of product/runtime/process modules. The sprint requirements further require stable read-only `/api/v1` objects, structured unknown-API errors, bounded allowlisted previews, and unchanged CLI/TUI behavior.

The evidence basis is high-confidence for dependency direction, error separation, capability restriction, cancellation, bounded work, diagnostics, and testing. HTTP-specific constants and exact route names are sprint decisions derived from those patterns and the project constraints rather than claims copied from the comparative reports.

## Risks

- **Compatibility drift:** Bundled browser code may tempt unreviewed DTO changes. Mitigation: keep explicit transport DTOs, versioned envelope fixtures, and require additive changes within `v1`; remove/rename/type changes require a new version or coordinated migration decision.
- **Opaque-reference staleness:** A referenced artifact can move or change between listing and preview. Decision: return `404 not_found` and require a state refresh; do not silently redirect or preview a different file.
- **Symlink and encoding bypass:** Lexical path checks alone are insufficient. Mitigation: app-owned reference resolution must canonicalize against the workspace, reject escapes after symlink resolution, allow only Markdown/JSON artifact classes, and test encoded traversal and symlink cases.
- **Read amplification:** Live dashboard refresh can repeat filesystem work. Mitigation: fixed response bounds, no recursive source-repository scans, request cancellation, `no-store`, and measurement before adding cache or bounded fan-out. A future cache must define invalidation and source-of-truth semantics explicitly.
- **Local disclosure:** Another local process or malicious webpage may probe the service. Mitigation: loopback bind validation, exact Host/Origin policy, no CORS, opaque references, no secrets/raw paths, safe errors, and restrictive response headers. This remains a local single-user model, not an isolation boundary against a compromised host account.
- **Health misuse:** External scripts may interpret `200` as full product validity. Mitigation: document health strictly as server/workspace-query readiness and keep validation findings in their own resource.
- **Partial-state expectations:** One failing dashboard query currently fails the aggregate API response rather than returning mixed panel errors. This is intentional for a coherent foundation contract; revisit only with measured failures and an explicit per-section error schema.
- **Arbitrary initial limits:** The 200-entry and 256 KiB limits are safety defaults without production measurements. Tests and documentation must make truncation visible; later tuning must preserve fields and compatibility behavior.
- **Scope collision with broader Phase 4 docs:** The PRD/TRD describe guarded operations and SSE, while Sprint 30 defers them. Route registration tests must prove no mutation/operation endpoints exist in this sprint so future work adds those capabilities deliberately.
- **Final reasoning handoff:** `projects/ultraplan-go/sprints/30-web-foundations/reasoning.md` must reference this document and carry forward its route, envelope, error, health, artifact-reference, and read-only decisions. That file is intentionally not written by this area-reasoning step.

No unresolved question blocks implementation. Questions about public API support, pagination, caching, partial dashboard responses, configurable bounds, guarded validation, confirmations, and SSE are deferred until a later sprint selects those capabilities and their security/compatibility requirements.
