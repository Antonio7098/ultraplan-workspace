# API Design: Guarded Web Operations and SSE Progress

> **Inputs Used:** `projects/ultraplan-go/sprints/31-web-operations/technical-handbook.md`, `projects/ultraplan-go/sprints/31-web-operations/requirements.md`, `projects/ultraplan-go/sprints/31-web-operations/sprint-index.md`, `projects/ultraplan-go/docs/ARCHITECTURE.md`, `projects/ultraplan-go/docs/PRD.md`, `projects/ultraplan-go/docs/TRD.md`, `studies/go-cli-study/reports/final/04-configuration-management.md`, `studies/go-cli-study/reports/final/05-error-handling.md`, `studies/go-cli-study/reports/final/07-state-context.md`, `studies/go-cli-study/reports/final/08-concurrency.md`, `studies/go-cli-study/reports/final/10-logging-observability.md`, `studies/go-cli-study/reports/final/11-testing-strategy.md`, `studies/go-cli-study/reports/final/13-security.md`, `studies/go-cli-study/reports/final/14-performance.md`

This document defines the frontend-owned, loopback-only `/api/v1` contract for preparing, starting, observing, and cancelling existing UltraPlan operations. It does not define new workflow semantics. The app and product modules remain authoritative for normalized scope, fingerprints, conflicts, results, and durable recovery; `internal/web` owns only HTTP validation, safe projection, confirmation/session policy, and bounded ephemeral delivery.

## Area Decisions

### Audience and contract boundary

- The audience is the embedded browser UI and local same-origin clients. This is not a public, partner, remote-worker, or multi-tenant API.
- `/api/v1` route shapes, JSON field names, error codes, operation states, and SSE event names become compatibility-controlled when implemented and documented. Additive optional fields are allowed; changing meanings, removing fields, reusing codes, or changing methods requires a later version.
- HTTP DTOs are explicit web types. Raw app, runtime, provider, process, lock, or persistence models are never serialized. `internal/web` maps only typed `internal/app` requests/results and safe diagnostics.
- The operation API accepts an allowlisted discriminated operation specification, never a shell command, argv, arbitrary executable, raw prompt, unrestricted filesystem path, or product-module object. Supported `kind` values map one-to-one to the existing validation, prompt-preview, dry-run, sprint-flow, execute, review, smoke, verify, study-run-loop, and cancellation use cases.

### Resource shape and methods

The API uses one synchronous preparation command followed by an asynchronous operation resource:

| Method and path | Purpose | Success |
| --- | --- | --- |
| `POST /api/v1/operations/prepare` | Normalize and validate an allowlisted operation specification, summarize impact, and issue confirmation. It starts no work. | `200 OK` with a preparation document. |
| `POST /api/v1/operations` | Consume a current confirmation and create server-owned work. | `202 Accepted` with an operation document and `Location: /api/v1/operations/{id}`. |
| `GET /api/v1/operations/{id}` | Return ephemeral lifecycle state and terminal result when retained; include durable refresh guidance. | `200 OK`. |
| `GET /api/v1/operations/{id}/events` | Replay retained safe events and subscribe to future progress. It cannot cause a state transition. | `200 OK`, `text/event-stream`. |
| `DELETE /api/v1/operations/{id}` | Request canonical cancellation. | `202 Accepted` while cancellation is pending; `200 OK` if cancellation was already requested or the operation is terminal. |

Status and result deliberately share one operation resource. A separate result endpoint would create two retention and not-found contracts without adding authority: terminal results are a bounded convenience projection, while durable product status is the recovery source.

Unknown `/api/` routes return the canonical JSON error envelope and never HTML. Unsupported methods return `405 Method Not Allowed` with `Allow`. New work is not queued by the web hub: start either reserves bounded capacity and returns `202`, or fails before token consumption with `operation_capacity` or `server_draining`. This prevents accepted-but-not-started work from surviving into shutdown ambiguity.

### Operation request

Prepare accepts:

```json
{
  "operation": {
    "kind": "sprint_flow",
    "scope": {
      "project": "ultraplan-go",
      "sprint": "31-web-operations"
    },
    "options": {
      "to_stage": "plan",
      "dry_run": false
    }
  }
}
```

`kind` selects a typed request variant. Each variant has an explicit field allowlist, enum validation, and app-owned normalization. Unknown fields are rejected rather than ignored so a browser cannot believe an option was honored when it was not. Client-supplied operation IDs, fingerprints, mutation classes, affected paths, runtime/provider metadata, result states, executable data, and confirmation claims are rejected.

Start repeats the operation specification and adds the issued token:

```json
{
  "operation": {
    "kind": "sprint_flow",
    "scope": {
      "project": "ultraplan-go",
      "sprint": "31-web-operations"
    },
    "options": {
      "to_stage": "plan",
      "dry_run": false
    }
  },
  "confirmation_token": "opaque-server-value"
}
```

Repeating the specification is intentional. The server normalizes it again and compares the normalized canonical representation and current governed-input fingerprint with the preparation record; it never trusts a client echo of the preparation summary.

### Preparation and confirmation

A successful preparation response contains:

```json
{
  "data": {
    "preparation_id": "prep_example01",
    "operation": {"kind": "sprint_flow", "scope": {}, "options": {}},
    "affected_paths": ["projects/ultraplan-go/sprints/31-web-operations"],
    "mutation_class": "sprint_flow",
    "runtime": {"kind": "opencode", "provider": "openai", "model": "gpt-5"},
    "harness": null,
    "input_fingerprint": "sha256:0123456789abcdef0123456789abcdef0123456789abcdef0123456789abcdef",
    "expires_at": "2026-08-16T12:02:00Z",
    "confirmation_token": "opaque-server-value"
  },
  "request_id": "req_example01"
}
```

- Preparation performs resolution, normalization, containment checks, current-state inspection, effective configuration resolution, and side-effect-free prerequisite checks. It does not acquire a mutation lock, launch a runtime/harness, write governed artifacts, or reserve operation capacity.
- The server keeps a bounded in-memory preparation record and returns an opaque token. The record is bound to the per-process secret, browser session, canonical normalized operation, governed-input fingerprint, and expiry. The token lifetime is two minutes.
- A token is single-use for successful operation creation. Failed transport policy or capacity checks before creation do not consume it; a normalized-request mismatch, stale fingerprint, expiry, or successful start does. Replay returns `confirmation_replayed`.
- Start repeats all app-owned normalization and checks the current fingerprint immediately before operation creation. Expiry returns `confirmation_expired`; session/request mismatch returns `confirmation_mismatch`; changed governed inputs return `stale_confirmation`. These are separate stable codes even though each maps to `409 Conflict`.
- Effective runtime/model/harness details are computed only after complete configuration resolution and validation. This is an application of the post-merge validation finding in `studies/go-cli-study/reports/final/04-configuration-management.md`, not a claim that configuration itself belongs to HTTP.
- CSRF proof and operation confirmation are independent. CSRF proves a same-origin session submitted the request; confirmation proves the user reviewed a particular current operation. Possessing one never substitutes for the other.

### Lifecycle and terminal results

The operation document has stable top-level fields:

```json
{
  "data": {
    "id": "op_example01",
    "kind": "sprint_flow",
    "state": "running",
    "reason": null,
    "created_at": "2026-08-16T12:00:10Z",
    "started_at": "2026-08-16T12:00:10Z",
    "finished_at": null,
    "last_event_id": "17",
    "durable_status": {"available": true, "refresh_path": "/api/v1/dashboard"},
    "result": null
  },
  "request_id": "req_example02"
}
```

Stable states are `accepted`, `running`, `cancelling`, `succeeded`, `failed`, `cancelled`, `interrupted`, and `cleanup_uncertain`. `accepted` is only the short transition between reservation and invocation, not a queue. `reason` is present for cancellation/interruption and is one of `user_request`, `server_shutdown`, `timeout`, or `recovery`; unknown internal strings are not projected.

The terminal `result` is a safe typed summary with operation-specific data or a canonical error object. Once a start has returned `202`, product failure, validation failure, timeout, cancellation, interruption, and cleanup uncertainty are operation outcomes returned by `GET` with `200`, not later HTTP transport failures. This prevents polling clients from confusing transport reachability with work success. An evicted operation returns `404 operation_not_found` plus safe durable refresh guidance; absence from the ephemeral hub does not imply that work never existed or failed.

During graceful shutdown, new starts return `503 server_draining` with `Retry-After`. Existing resources transition through `cancelling` with `reason: server_shutdown`, then to the app-reconciled terminal state. SSE closes only after a terminal event or a bounded `cleanup_uncertain`/`interrupted` outcome has been published. A forced stop has no opportunity to complete this API sequence; the next server instance reports recovery from durable product state rather than recreating old operation IDs.

### SSE contract

SSE is one-way observation. Each frame uses a decimal event ID monotonically increasing within one operation, a stable event name, and a JSON data object:

```text
id: 18
event: progress
data: {"operation_id":"op_example01","time":"2026-08-16T12:00:15Z","sequence":18,"payload":{"message":"Validating plan","current":3,"total":8}}
```

Stable event names are `snapshot`, `progress`, `warning`, `finding`, `artifact`, `cancel_requested`, `recovery_required`, and `terminal`. Event payloads are event-specific safe projections; raw provider events, raw prompts, unrestricted stdout/stderr, authorization data, environment values, and absolute paths are forbidden. Unknown future app event kinds are either mapped to a safe `progress` summary or omitted; native names are not promoted automatically into the API contract.

- The first connection emits a `snapshot`, then replays retained events newer than `Last-Event-ID`, then follows live events. IDs are assigned under the operation hub's serialization point, not copied from concurrently arriving producer IDs.
- If `Last-Event-ID` is older than retained history, the stream emits `recovery_required` with the oldest/newest retained IDs and a durable refresh path, then a current `snapshot`, and continues live. The server never fabricates missing events.
- A heartbeat comment `: heartbeat` is sent every 15 seconds and has no event ID or lifecycle meaning.
- Browser disconnect cancels only the subscriber context. It never invokes operation cancellation or changes the result.
- Each subscriber has a 32-event outbound queue. A full queue disconnects that subscriber; producers never block and the subscriber must reconnect. Buffer rollover is recovered through `recovery_required`, not silent event loss.
- Streams close after the terminal frame is flushed, subscriber eviction, server shutdown completion, or the configured maximum stream lifetime. Reconnect always rechecks session ownership and current hub state.

### Bounds and load shedding

The following are Sprint 31 defaults and hard upper bounds unless a lower validated server configuration is selected:

| Resource | Bound |
| --- | --- |
| JSON request body | 64 KiB |
| Active server-owned operations | 8 |
| Preparations retained | 128 for at most 2 minutes |
| Recent events per operation | 256 events and 256 KiB total |
| Encoded event payload | 16 KiB |
| Terminal result projection | 256 KiB |
| Subscribers per operation | 8 |
| Concurrent SSE streams server-wide | 32 |
| Per-subscriber outbound queue | 32 events |
| Terminal operation retention | 10 minutes |
| SSE heartbeat | 15 seconds |
| One SSE connection lifetime | 30 minutes, followed by reconnect |

Capacity rejection is explicit and never represented as an operation failure: `429 operation_capacity` for operation/preparation/subscriber limits and `503 server_draining` during shutdown. `Retry-After` is included when useful. Product/runtime concurrency may be lower and remains app-owned. These limits apply the bounded queue, slow-consumer, and incremental-retention findings in `studies/go-cli-study/reports/final/08-concurrency.md` and `studies/go-cli-study/reports/final/14-performance.md`; buffer pooling is rejected until profiling shows a hot path.

### Error envelope

All pre-acceptance and transport failures use:

```json
{
  "error": {
    "code": "stale_confirmation",
    "message": "The operation inputs changed after confirmation. Prepare it again.",
    "retryable": true,
    "details": {"action": "prepare_again"}
  },
  "request_id": "req_example03"
}
```

`details` is code-specific and allowlisted. It may contain a safe field name, mutation class, operation ID, conflicting operation ID/kind, expiry, retry delay, durable refresh path, or corrective action. It never contains wrapped error text, provider payloads, unrestricted stderr, absolute paths, tokens, environment values, stack traces, or lock-file internals.

| Code | HTTP status | Retry meaning |
| --- | --- | --- |
| `invalid_request` | `400` | Correct the request. |
| `csrf_failed`, `origin_rejected`, `session_required` | `403` | Establish a valid same-origin session; do not retry unchanged. |
| `operation_not_found` | `404` | Refresh durable state; the ephemeral record may have expired. |
| `confirmation_expired`, `confirmation_mismatch`, `confirmation_replayed`, `stale_confirmation` | `409` | Prepare again, except a mismatch also requires correcting the request/session. |
| `operation_conflict` | `409` | Wait for or cancel the conflicting product-owned mutation when allowed. |
| `validation_failed` | `422` | Correct governed inputs using safe findings. |
| `prerequisite_unavailable` | `424` | Satisfy the named safe prerequisite. |
| `operation_capacity`, `subscriber_capacity` | `429` | Retry after the indicated delay. |
| `server_draining` | `503` | Retry against a running server instance. |
| `internal_failure` | `500` | Inspect correlated logs; no internal detail is exposed. |

App errors are classified with `errors.Is`/`errors.As`, never message matching. Typed conflict, stale-input, validation, prerequisite, cancellation, timeout, runtime, cleanup, and internal errors retain rich data inside the app while the web mapper emits only the compatibility-controlled safe subset. This follows the structured error and user/operational separation evidence in `studies/go-cli-study/reports/final/05-error-handling.md`.

### Session, authorization, retry, and idempotency

- The API has no account authentication or tenant model. Authority derives from loopback binding plus a same-origin per-process browser session with strict Host/Origin checks, secure cookie attributes appropriate to local HTTP, CSRF proof on `POST`/`DELETE`, and security headers.
- Preparations and operation resources are scoped to the initiating server session. Other sessions receive `404`, not an existence-revealing authorization response. Server shutdown still cancels all operations regardless of session.
- `GET` and prepare are safe to retry. `DELETE` is idempotent: repeated requests report the existing cancellation/terminal state and canonical cancellation is invoked at most once.
- Start is not automatically retryable. Its single-use confirmation is the deduplication boundary: replay after accepted creation returns `confirmation_replayed`. A general `Idempotency-Key` store is rejected for this sprint because it would add another retention/identity contract and could be mistaken for durable job state. Clients that lose the `202` response refresh visible operations/durable state and prepare again only when no operation was created.

### Observability and tests

Every request receives a `request_id`; every accepted operation receives an `operation_id`. Safe app/runtime/harness run IDs may be correlated but do not replace either identifier. Structured logs use stable fields for request ID, operation ID, session hash, route, method, operation kind, project/study/sprint reference, state transition, cancellation reason, duration, event sequence, subscriber disconnect reason, buffer rollover, and error code. Confirmation tokens, CSRF values, cookie values, raw payloads, and unsafe paths are never logged. SSE events and logs are separate projections: an event is user progress, while a log may contain safe operational detail that is not part of the API.

Local counters track starts, active operations, terminal outcomes by safe category, start rejection reason, cancellation reason, active streams, slow-subscriber disconnects, replay gaps, event drops before projection, and shutdown cleanup outcomes. Sprint 31 does not add a network metrics endpoint or tracing system.

Required behavioral tests are:

- Route/method, unknown-API-route, content-type, strict JSON/unknown-field, body-limit, and envelope tests.
- Per-operation-variant schema and normalization tests proving caller-controlled IDs, paths, executable data, mutation classes, fingerprints, and runtime claims are rejected.
- Preparation tests for no side effects, effective configuration summary, exact normalized binding, expiry, session mismatch, stale fingerprint, replay, and failed-capacity non-consumption.
- Host, Origin, session, CSRF, cookie, security-header, object-ownership, and redaction tests.
- Start/status/result tests for `202`/`Location`, every lifecycle and terminal outcome, no hidden queue, eviction, durable refresh guidance, conflict mapping, draining rejection, and exact-once cancellation.
- SSE tests for event names and schema, per-operation monotonic IDs under concurrent publication, initial snapshot, replay, `Last-Event-ID`, rollover recovery, heartbeat comments, payload/event bounds, terminal flush, subscriber capacity, slow-reader disconnection, browser disconnect isolation, and shutdown ordering.
- Compatibility fixtures for JSON envelopes and SSE frames, with semantic assertions for lifecycle behavior. Golden tests are limited to representative complete projections where a full diff adds value.
- Race and leak tests with fake clocks, deterministic IDs, fake app operations, blocked/slow writers, cancellation barriers, and temporary workspaces, followed by `go test ./...` and `go test -race ./...`.

The testing split applies the behavior-first, centralized fake, HTTP test-server, and selective golden guidance in `studies/go-cli-study/reports/final/11-testing-strategy.md`. Timing assertions use fake clocks or synchronization barriers rather than sleeps.

## Trade-Offs

| Decision | Benefit | Cost / rejected alternative |
| --- | --- | --- |
| Prepare then start with a single-use server record | Gives the user a current normalized scope and binds approval to exact inputs without client-authored claims. | A signed stateless token would simplify lookup but cannot enforce replay or bounded per-session issuance without another store; direct start has no meaningful confirmation boundary. |
| Async operation resource for all accepted work | Unifies short and long operations, cancellation, shutdown, and recovery semantics. | Synchronous execution is simpler for validation but couples request timeouts/disconnects to work and creates inconsistent operation behavior. |
| One status/result resource | One lifecycle and retention contract; terminal outcome remains explicit. | Separate status and result endpoints add not-found and freshness ambiguity without new authority. |
| No web queue | Shutdown and capacity behavior remain truthful; every `202` owns cancellable work immediately. | A bounded queue could smooth bursts but adds queued cancellation, fairness, and persistence expectations resembling a scheduler. |
| Product outcome inside a `200` operation document | Separates accepted-work failure from HTTP transport failure and makes polling deterministic. | Returning `5xx` from status for a failed run conflates a reachable API with a failed product operation. |
| Drop/disconnect slow subscribers, never producers | Browser behavior cannot alter product completion and memory remains bounded. | Blocking delivery offers gap-free observation but can deadlock work; silently dropping events lies about continuity. |
| Recovery marker plus snapshot after rollover | Truthfully identifies missing observations while allowing live monitoring to continue. | Unlimited replay violates ephemerality; closing with no explanation makes the browser infer failure; durable event persistence creates a second state model. |
| Per-operation event IDs | Monotonic assignment is local, cheap, and deterministic under one hub lock/owner. | Server-global IDs add contention and suggest replay across operations; producer-native IDs may be concurrent, missing, or unsafe. |
| Typed stable error codes with allowlisted details | Supports browser branching and safe guidance while preserving rich internal errors. | Error strings are easy initially but fragile for compatibility and can leak internals; exposing raw typed app objects couples layers. |
| Session-scoped local resources | Reduces cross-tab/session disclosure and binds confirmations to the reviewing browser context. | Unscoped loopback access is simpler but treats every local browser context as equally trusted; full account auth is outside scope. |
| Fixed conservative defaults with lower configuration allowed | Makes acceptance tests and resource protection concrete. | Entirely configurable/unbounded values make behavior non-portable; adaptive limits are premature without measurements. |
| Selective golden fixtures plus semantic lifecycle tests | Protects complete wire projections without making concurrency tests formatting-sensitive. | All-golden testing obscures behavioral failures; only focused assertions can miss accidental contract drift. |

## Evidence

- **Requirements and project contracts:** `projects/ultraplan-go/sprints/31-web-operations/requirements.md` requires methods, confirmation binding, bounded ephemerality, SSE semantics, explicit/server-shutdown cancellation, typed errors, product-owned locks, recovery, and transport-only web dependencies. `projects/ultraplan-go/docs/TRD.md` section 18A supplies the initial route family and establishes `/api/v1` compatibility control. `projects/ultraplan-go/docs/ARCHITECTURE.md` makes browser projections non-authoritative and requires `internal/web -> internal/app`; `projects/ultraplan-go/docs/PRD.md` defines the browser as a guarded local interface over the same use cases.
- **Handbook orientation:** `projects/ultraplan-go/sprints/31-web-operations/technical-handbook.md` identifies typed safe projection, context-owned operations, explicit lock ownership, bounded pub/sub, correlated events, and durable refresh as the relevant cross-cutting pressures. The decisions above resolve those pressures for the HTTP contract rather than copying the handbook's candidate patterns.
- **Normalized preparation:** `studies/go-cli-study/reports/final/04-configuration-management.md` reports that high-quality tools merge all sources before centralized validation (`k9s/internal/config/k9s.go:423-451`) and distinguish explicit overrides (`restic/internal/global/global.go:139,147`). Report finding: effective inputs should be validated after resolution. Sprint inference: confirmation binds the complete canonical operation and effective governed-input fingerprint, not partially decoded browser fields.
- **Error contract:** `studies/go-cli-study/reports/final/05-error-handling.md` finds that typed errors carry actionable data (`restic/internal/restic/lock.go:47`; `go-task/errors/errors_task.go:13-32`), `%w` preserves classification, and user rendering should be separate from operational detail (`k9s/internal/model/flash.go:100-103`). Report finding: strings are a weak programmatic boundary. Sprint inference: stable API codes and allowlisted details are projections of typed app errors, while accepted-work failures live in the operation result.
- **Lifecycle ownership:** `studies/go-cli-study/reports/final/07-state-context.md` finds predictable cancellation where a root context reaches long-running I/O (`helm/pkg/cmd/install.go:333-347`) and identifies lock/deduplication as explicit operation coordination (`restic/internal/restic/lock.go:105`; `go-task/task.go:438-469`). It also warns that `context.Background()` severs lifecycle propagation. Sprint inference: start creates server-owned context, DELETE and shutdown cancel that same lineage, and SSE subscriber contexts remain separate.
- **Bounded concurrency:** `studies/go-cli-study/reports/final/08-concurrency.md` supports localized goroutine ownership, bounded semaphores (`k9s/internal/pool.go:21,30,37`), explicit subscription cleanup (`opencode/internal/pubsub/broker.go:67-82`), `sync.Once` for exact-once cleanup, and timeout-bounded waits (`opencode/cmd/root.go:261-279`). Report warning: unbuffered/abandoned consumers can deadlock producers. Sprint inference: no web queue, bounded subscriber buffers, exact-once cancellation, and slow-subscriber disconnection.
- **Correlation without projection leakage:** `studies/go-cli-study/reports/final/10-logging-observability.md` finds stable structured keys (`k9s/internal/slogs/keys.go:6-231`), component tagging, and strict user-output/diagnostic separation. Sprint inference: request/operation IDs and stable lifecycle fields correlate HTTP, SSE, app, and logs, but logs and browser events remain different safe views.
- **Behavioral verification:** `studies/go-cli-study/reports/final/11-testing-strategy.md` supports centralized fakes and HTTP mocks (`gh-cli/pkg/httpmock/stub.go:35-199`), behavior assertions, and selective golden output (`helm/internal/test/test.go:43`), while warning against implementation-detail assertions (`k9s/internal/view/pod_test.go:23`). Sprint inference: use `httptest`, fake clocks/operations, semantic lifecycle assertions, and a small wire-format compatibility fixture set.
- **Trust boundary:** `studies/go-cli-study/reports/final/13-security.md` finds that explicit trust boundaries and permission/confirmation points improve safety (`opencode/internal/permission/permission.go:44-108`), typed redaction prevents accidental disclosure (`restic/internal/options/secret_string.go:15-20`), and explicit argument arrays avoid shell interpretation (`lazygit/cmd_obj_builder.go:38`). Sprint inference: accept typed operation capabilities only, separate CSRF from confirmation, session-scope resources, and redact before data enters retained buffers.
- **Streaming and resource bounds:** `studies/go-cli-study/reports/final/14-performance.md` finds that streaming, incremental bounded structures, and semaphore limits preserve responsiveness (`opencode/internal/llm/provider/provider.go:56`; `k9s/internal/pool.go:26-48`), and reports pragmatic slow-consumer event drops. It cautions against speculative pooling. Sprint inference: bounded event/result retention, explicit rollover recovery, producer-independent subscribers, and no buffer pool until measured.

## Risks

- **Fingerprint incompleteness:** If an app operation omits a governed input, runtime/model selection, harness identity, or normalized option from its fingerprint, a confirmation can remain apparently current after material change. Each operation variant needs a canonicalization/fingerprint contract test.
- **Preparation side effects:** Runtime or harness preflight can accidentally launch work, warm mutable state, or acquire locks. Preparation must use explicitly side-effect-free inspection APIs; expensive or mutating checks move to start and may still fail after confirmation with a safe prerequisite/conflict outcome.
- **Start response loss:** A connection can fail after token consumption and operation creation. With no idempotency-key contract, the browser must query current visible operations/durable state before preparing again. Tests must force this boundary and prove duplicate work is not started by token replay.
- **Terminal arbitration races:** Completion, explicit cancellation, timeout, and server shutdown can race. App/hub integration must select one truthful terminal outcome, invoke cancellation once, retain cleanup metadata, and never overwrite an already-authoritative success/failure after cancellation loses the race.
- **SSE gap misunderstanding:** A browser may treat `recovery_required` as operation failure or treat a snapshot as complete history. Documentation and UI must label event history as partial observation and direct the user to durable status.
- **Redaction after retention:** Redacting only during JSON/SSE encoding leaves unsafe data in memory and reconnect buffers. App-to-web projection must sanitize before publication and retention; tests should inject secrets, provider payloads, stderr, and absolute paths at every event/result/error boundary.
- **Memory multiplication:** Per-operation buffers, per-subscriber queues, terminal results, and concurrent streams multiply even when each bound looks small. The configured limits require aggregate allocation tests and race tests; a limit increase must consider the full product, not one collection.
- **Session lifecycle:** A per-process session rotation or cookie loss can make an active operation invisible to its original page while the server still owns it. This is safer than cross-session disclosure but requires durable refresh guidance and an operator-visible aggregate shutdown path.
- **Local-only threat assumptions:** Loopback and same-origin reduce exposure but do not make browser input trusted; DNS rebinding, malicious local pages, extensions, and other local processes remain concerns. Host/Origin/CSRF/session/body/path checks must apply consistently to every operation and SSE route.
- **Compatibility drift:** Operation-specific result fields and new progress kinds can leak raw app models if added ad hoc. Every new field/event/error code requires a safe DTO, compatibility fixture, and redaction review; unknown app/runtime events remain omitted or generically summarized.
- **Unmeasured default limits:** The concrete Sprint 31 limits are conservative engineering defaults, not benchmark-derived optima. Normal tests must prove enforcement; later tuning should use measured event sizes, operation mix, and browser behavior without weakening the bounded contract.
- **Open question for later versions:** A durable worker or multi-process server would invalidate same-process token, session, operation-ID, and retention assumptions. That architecture must introduce explicit leases, durable idempotency, ownership, authentication, and replay semantics under a new API version rather than extending this ephemeral contract invisibly.
