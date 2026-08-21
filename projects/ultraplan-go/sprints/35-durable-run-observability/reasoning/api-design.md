# API Design: Durable Workspace Runs

> **Inputs Used:** `projects/ultraplan-go/sprints/35-durable-run-observability/technical-handbook.md`, `projects/ultraplan-go/sprints/35-durable-run-observability/requirements.md`, `projects/ultraplan-go/sprints/35-durable-run-observability/sprint-index.md`, `projects/ultraplan-go/sprints/35-durable-run-observability/code-context.md`, `projects/ultraplan-go/docs/ARCHITECTURE.md`, `projects/ultraplan-go/docs/PRD.md`, `projects/ultraplan-go/docs/TRD.md`, `studies/go-cli-study/reports/final/05-error-handling.md`, `studies/go-cli-study/reports/final/07-state-context.md`, `studies/go-cli-study/reports/final/08-concurrency.md`, `studies/go-cli-study/reports/final/10-logging-observability.md`, `studies/go-cli-study/reports/final/11-testing-strategy.md`, `studies/go-cli-study/reports/final/12-extensibility.md`, `studies/go-cli-study/reports/final/13-security.md`, `studies/go-cli-study/reports/final/14-performance.md`, `system/reasoning/api-design-reasoning-template.md`, `../ultraplan-go/internal/web/api_compatibility_test.go`, `../ultraplan-go/internal/web/operations_contract_test.go`, `../ultraplan-go/internal/web/operation_handlers.go`, `../ultraplan-go/internal/web/operations.go`, `../ultraplan-go/internal/web/routes.go`, `../ultraplan-go/internal/web/static/app.js`, `../ultraplan-go/internal/web/handlers.go`

This area defines the stable local API and cross-surface contract for durable operational runs. It does not select the repository implementation, transfer workflow authority from sprint or study modules, or replace agentwrap supervision. The supported public topology is multiple UltraPlan processes on one host attached to the same canonical workspace. Shared-filesystem multi-host liveness, cancellation, and live-follow guarantees are not part of this API contract.

## Area Decisions

### 1. Canonical resource and audience

Proceed with a new canonical `run` resource under `/api/v1/runs`. Keep `/api/v1/operations` as a compatibility projection rather than evolving its frozen transport structs into the durable domain model.

The audience is:

- the dependency-free embedded browser client
- local automation using the documented loopback `/api/v1` API
- CLI and TUI adapters consuming the same application capability directly

The API is not a hosted, tenant-aware, partner, webhook, or remote-worker API. The storage format, lease protocol, process signaling mechanism, and agentwrap records remain internal.

| Method and path | Decision |
| --- | --- |
| `GET /api/v1/runs` | List retained workspace runs with bounded pagination and filters. |
| `GET /api/v1/runs/{run_id}` | Return the current durable snapshot, including compacted or tombstone state. |
| `GET /api/v1/runs/{run_id}/events` | Replay committed events after a cursor, then follow newly committed events with SSE. |
| `DELETE /api/v1/runs/{run_id}` | Persist an authorized, idempotent cancellation request and return the resulting durable snapshot. |
| `POST /api/v1/operations` | Remain the confirmed start command and return the existing operation document shape. Its accepted ID is the durable run ID. |
| `GET /api/v1/operations` | Remain the existing active-operation array shape, but project workspace-wide active runs without originating-session filtering. |
| `GET/DELETE /api/v1/operations/{id}` | Resolve an operation ID to a durable run and project the existing v1 operation document/cancellation behavior. |
| `GET /api/v1/operations/{id}/events` | Project durable events into the existing stable operation event names and frame shape. |

New run routes support `HEAD` wherever the existing router adds `HEAD` to a `GET`. All JSON responses retain the existing `{data, meta}` and `{error, meta}` envelopes, `api_version: "v1"`, `request_id`, `generated_at`, and `Cache-Control: no-store` policy.

The equivalent stable CLI surface is:

| Command | Contract |
| --- | --- |
| `ultraplan run list [--active] [--project ...] [--sprint ...] [--study ...] [--json]` | Same filters, lifecycle facts, and ordering as `GET /api/v1/runs`. |
| `ultraplan run show <run-id> [--json]` | Same durable snapshot and recovery guidance as run detail. |
| `ultraplan run follow <run-id> [--after <sequence>] [--json]` | Replay then follow the same sanitized committed event envelope. |
| `ultraplan run cancel <run-id> [--json]` | Invoke the same durable cancellation command and idempotency rules. |
| `ultraplan run diagnostics [<run-id>] [--json]` | Report safe run-control health and local metrics without exposing raw payloads. |

TUI actions use the same application capability. They do not call these HTTP routes or parse CLI output.

### 2. Stable identity and resource representation

The server assigns an opaque, workspace-unique ID with the lexical form `run_<base32-random-128-bit>`. Attempt IDs use `att_<base32-random-128-bit>`. Clients may validate the safe opaque alphabet and maximum length but must not parse, sort, or derive time, workspace, owner, or authority from either ID.

For newly accepted asynchronous operations, the compatibility `operationDocument.id` equals the durable `run_id`. No second browser operation identity is created. Existing `op_*` identifiers remain valid lookup keys only through compatibility resolution.

The canonical detail resource has this conceptual shape:

```json
{
  "schema_version": 1,
  "run_id": "run_...",
  "record_state": "full",
  "operation_kind": "sprint-stage",
  "target": {
    "project": "ultraplan-go",
    "sprint": "35-durable-run-observability",
    "study": null,
    "stage": "area-reasoning",
    "task": null
  },
  "lifecycle": {
    "state": "running",
    "terminal": false,
    "reason": null,
    "accepted_at": "...",
    "queued_at": "...",
    "started_at": "...",
    "finished_at": null
  },
  "liveness": {
    "classification": "live",
    "heartbeat_at": "...",
    "lease_expires_at": "...",
    "owner_id": "owner_...",
    "fencing_generation": 4,
    "process_identity_verified": true
  },
  "current_attempt_id": "att_...",
  "attempts": [],
  "progress": {},
  "events": {
    "oldest_retained_sequence": 12,
    "last_sequence": 48,
    "history_complete": false,
    "retained_until": "..."
  },
  "cancellation": {
    "state": "none",
    "requested_at": null,
    "acknowledged_at": null,
    "reason": null
  },
  "terminal": null,
  "diagnostics": [],
  "correlations": {},
  "links": {}
}
```

The JSON DTO is explicit rather than an encoded internal struct. These invariants apply:

- `record_state` is `full`, `compacted`, or `tombstone`.
- lifecycle states are `accepted`, `queued`, `running`, `cancelling`, `succeeded`, `failed`, `cancelled`, `timed_out`, `interrupted`, or `cleanup_uncertain`.
- `lifecycle.terminal` is the compatibility-safe terminal discriminator. Clients do not infer terminality from a closed enum.
- `stalled` is a liveness classification while lifecycle remains active; reconciliation may later win the immutable `interrupted` or `cleanup_uncertain` terminal outcome.
- `terminal` is null until terminal arbitration commits. Once present, its outcome, winning transition, timestamp, and safe result cannot be rewritten.
- `attempts` contains bounded summaries and stable attempt IDs, not raw provider records.
- `correlations` may contain safe stage/task execution IDs, agentwrap run ID, provider session ID, external harness run ID, and process identity digest. Missing correlation is null or omitted, never fabricated.
- target paths are absent or workspace-relative display paths. Absolute home paths are not returned.
- usage fields retain known/unknown flags; unknown cost or token values are not converted to zero.
- diagnostics contain stable code, component, safe message, retryability, and guidance. They exclude wrapped internal errors, prompts, provider payloads, headers, credentials, and unrestricted stderr.

Product status is linked or summarized as a projection. The run resource never claims that an artifact, flow stage, execute task, review, or smoke verdict succeeded merely because the operational run terminated successfully.

### 3. Listing, ordering, filtering, and bounds

`GET /api/v1/runs` defaults to the 50 newest retained runs ordered by `accepted_at DESC, run_id DESC`. It accepts:

- `lifecycle=active`, `terminal`, or one exact lifecycle state
- one or more of `project`, `sprint`, `study`, `stage`, `task`, and `operation_kind`
- `limit`, default 50 and maximum 200
- `page_after`, an opaque repository-issued pagination cursor distinct from an event sequence

Invalid combinations or limits return `400 invalid_request`. The response metadata reports returned count and total count when the store can compute it within the bounded query. A missing total is represented by omission, not a guessed value.

The browser top bar and compatibility `GET /api/v1/operations` explicitly request/project `lifecycle=active`. Active means `accepted`, `queued`, `running`, or `cancelling`, regardless of which CLI, TUI, or web process accepted the run. Planning-stage readiness is never used as active-run state.

### 4. Event and replay contract

The canonical run event envelope is:

```json
{
  "schema_version": 1,
  "run_id": "run_...",
  "sequence": 49,
  "committed_at": "...",
  "type": "progress",
  "attempt_id": "att_...",
  "stage": "execute",
  "task": "T-03",
  "payload": {},
  "omissions": []
}
```

The event's decimal `sequence` is also the SSE `id`. Events are sanitized, size-bounded, assigned a per-run sequence, and durably committed before the SSE adapter may emit them. The delivery guarantee is at-least-once after a committed sequence. Clients deduplicate by `(run_id, sequence)` and must tolerate reconnect duplicates. The API does not promise exactly-once network delivery.

`GET /api/v1/runs/{id}/events` accepts either `Last-Event-ID` or `after=<decimal sequence>`. If both are supplied and differ, it returns `400 invalid_cursor`. An omitted cursor starts at the oldest retained event and does not claim that earlier compacted history exists. A requested cursor greater than the current last sequence returns `409 cursor_ahead` with the current snapshot.

If a requested cursor is below `oldest_retained_sequence - 1`, the server returns `409 replay_gap` as JSON before starting SSE. The error details contain:

- `requested_after`
- `oldest_retained_sequence`
- `last_sequence`
- `gap_reason`, such as `retention`, `compaction`, or `corruption_recovery`
- the current durable run snapshot
- actions to replay from the retained boundary or follow only future events

This typed pre-stream response avoids pretending a partial stream is complete. The browser performs a detail preflight when resuming a saved cursor and, on an EventSource connection error, fetches detail to distinguish a gap from transport loss.

Each connection replays at most 512 historical events before the server closes it. EventSource reconnect sends the last delivered ID and continues catch-up; live mode begins only after the retained backlog is exhausted. A single encoded event remains capped at 16 KiB. Oversized producer detail is replaced before append with a safe warning/omission event. SSE heartbeat comments remain transport-only and carry no run sequence.

Slow subscribers are disconnected when their bounded delivery queue fills. They recover from their last committed cursor; subscriber lag cannot block execution or journal commits. Journal retention and compaction may remove detail, but the snapshot always exposes the lower replay boundary and whether history is complete.

Canonical event types are `snapshot`, `progress`, `warning`, `finding`, `artifact`, `cancel_requested`, `cancel_acknowledged`, `recovery_required`, and `terminal`. Unknown future types do not terminate the client; clients display a generic progress item and refresh the snapshot. The compatibility operation stream continues to expose only its frozen event names and maps new canonical events to the nearest safe existing name.

### 5. Acceptance and retry semantics

`POST /api/v1/operations` remains asynchronous and returns `202 Accepted` only after the durable run acceptance record exists. Child execution cannot begin before that commit. Persistence failure returns `503 run_persistence_unavailable`, sets `Retry-After` when useful, and guarantees that no child was started.

The existing confirmation token remains bound to the normalized request, current governed-input fingerprint, browser session, and expiry. The run store records a non-reversible digest of the consumed token as the acceptance deduplication key. Repeating the same consumed token and identical canonical request during retained dedupe history returns the same run and `Location`; the raw token is never stored. The same token with a different request returns `409 idempotency_conflict`. Unknown, expired, stale, or never-accepted tokens retain the existing typed confirmation errors.

The compatibility response preserves the existing operation document fields and `Location: /api/v1/operations/{run_id}`. It adds an HTTP `Link: </api/v1/runs/{run_id}>; rel="canonical"` header rather than adding required JSON fields to the frozen operation shape.

Runtime-free synchronous queries are not runs. Any operation accepted for asynchronous execution through the shared operation runner is a durable run, whether or not that particular path eventually starts an agentwrap child. This keeps cancellation, inspection, and compatibility semantics uniform while making runtime-backed acceptance fail closed as required.

### 6. Read visibility and mutation authorization

Run list, detail, retained events, operation compatibility reads, and HTML run detail are workspace-readable from any fresh browser session served by a loopback UltraPlan server attached to that workspace. They do not compare the current session ID with the session that accepted the run.

Read responses expose only the sanitized resource described above. They remain protected by loopback binding, strict Host/Origin policy where applicable, same-origin browser policy, no permissive CORS, path containment, and response redaction. The contract does not add accounts, tenants, bearer tokens, or network sharing.

Cancellation requires current authority, not the originating session:

- browser/API cancellation requires a valid current same-origin session and CSRF proof
- CLI/TUI cancellation requires successful workspace resolution and local OS access through the shared application capability
- no caller may provide owner ID, fencing generation, PID, process birth, terminal state, event sequence, or arbitrary cancellation reason
- retry is not added as a run mutation in this sprint; product-owned resume/retry commands remain separate guarded operations

This is an intentional widening of read visibility within one local workspace and an intentional preservation of fresh authorization for mutation.

### 7. Idempotent cancellation and terminal races

`DELETE /api/v1/runs/{id}` has no request body. The server assigns the canonical reason `user_request`; server shutdown and timeout use internal trusted reasons.

The response rules are:

| Condition | Status | Result |
| --- | --- | --- |
| First cancellation request durably committed | `202` | Snapshot with `cancellation.state` set to `requested`, `routing`, `acknowledged`, or `uncertain`. |
| Duplicate request while cancellation is pending | `200` | Current snapshot; no second signal or terminal proposal is created. |
| Run already terminal | `200` | Immutable winning terminal snapshot; cancellation does not rewrite it. |
| Request persistence fails | `503` | `cancellation_persistence_failed`; the API does not claim cancellation was requested. |
| Current owner is stale or unreachable after persistence | `202` | Durable request remains visible as `uncertain`; reconciliation guidance is returned. |
| Caller lacks current mutation authority | `403` | No cancellation record or owner signal. |

Owner acknowledgement is a separate durable fact from terminal cancellation. Completion, failure, timeout, interruption, cleanup uncertainty, and cancellation race through the same terminal compare-and-set. A stale owner or repeated observer can only read the winner.

### 8. Errors, tombstones, and recovery

Keep the current v1 error envelope and add stable codes with typed safe details. Do not classify storage, retention, ownership, or cancellation failures by substring matching.

| HTTP status | Code | Meaning and required details |
| --- | --- | --- |
| `400` | `invalid_run_id`, `invalid_cursor`, `invalid_request` | Malformed safe input; include the invalid field, not raw body content. |
| `403` | `mutation_forbidden` | Current session/CSRF/local authority is insufficient. |
| `404` | `run_not_found` | A syntactically current run ID has no record or tombstone; point to workspace run list and diagnostics. |
| `409` | `cursor_ahead`, `replay_gap`, `idempotency_conflict`, `terminal_conflict` | Include cursor boundaries or the immutable current snapshot as applicable. |
| `410` | `events_expired`, `legacy_operation_not_retained`, `run_tombstone_expired` | Identity/history was known or is recognizably legacy but detail is no longer retained; include precise recovery actions. |
| `429` | `subscriber_capacity`, `replay_capacity` | Bounded delivery capacity is full; include `Retry-After`. |
| `503` | `run_store_unavailable`, `run_persistence_unavailable`, `cancellation_persistence_failed`, `reconciliation_unavailable` | Required durable truth is unavailable; no stale in-memory success projection. |

A compacted run returns `200` with `record_state: "compacted"`, the current snapshot, and explicit event boundaries. A retained tombstone returns `200` with `record_state: "tombstone"`, stable identity, target summary, terminal outcome if known, retention facts, and recovery guidance. Event requests against a tombstone return `410 events_expired` with the tombstone snapshot.

An unresolvable `op_*` URL from the pre-durable implementation cannot have history fabricated. It returns `410 legacy_operation_not_retained`, explains that the old process-local record never entered durable storage, and links to the workspace run list plus product status pages. This replaces the generic `Operation not retained` dead end while remaining truthful.

### 9. Compatibility and migration

The change is additive at `/api/v1/runs` and intentionally behavior-changing only where Sprint 35 requires workspace visibility.

- Existing operation routes, methods, JSON field order, success/error envelopes, lifecycle names, and stable SSE event names remain compatibility fixtures.
- The existing operation document `id` carries the durable run ID for new accepted work; clients that treat it as opaque continue to work.
- `GET /api/v1/operations` keeps the existing array item schema but changes from session-local to workspace-wide active results.
- Operation status and event reads stop requiring the originating session. Mutation still requires a current authorized session.
- Canonical lifecycle values that do not exist in the operation contract are projected: `queued` to `accepted`, `timed_out` to `failed`, and canonical cancellation routing states to `cancelling`.
- Existing operation SSE clients continue to receive `snapshot`, `progress`, `warning`, `finding`, `artifact`, `cancel_requested`, `recovery_required`, and `terminal` only.
- HTML `/operations/{id}` resolves to or redirects to stable run detail when mapped, renders tombstone/recovery state when retained, and never uses an unexplained generic 404 for a recognized legacy ID.
- Existing sprint/study locks, flow state, execute state, stage session checkpoints, cleanup markers, and external smoke evidence are correlation sources only. No migration imports them as invented run history or lets the run API rewrite them.
- Run JSON and persisted schema versions evolve independently. API clients rely on `meta.api_version`; repository migrations remain internal and fail visibly through health/errors.

Future API evolution is additive within v1 where old clients can safely ignore fields. A required field removal, semantic reuse, route/method change, or incompatible event envelope requires `/api/v2`. New lifecycle or event values require clients to use explicit `terminal` and generic unknown-event behavior rather than fail closed on enum parsing.

### 10. Operational observability surfaces

Every API command, event, and error log carries safe `request_id`, `run_id`, `attempt_id`, operation kind, project/sprint/study/stage/task identifiers, owner/fencing generation when relevant, runtime/agentwrap correlation, durable sequence, lifecycle transition, and terminal winner. High-cardinality IDs are structured log fields, not metric labels.

Extend `GET /api/v1/health` additively with a safe `run_control` component summarizing store availability, append/terminal persistence health, stale-owner count, reconciliation backlog, oldest backlog age, and compaction health. Detailed local metrics are exposed by `ultraplan run diagnostics --json`, not by a new unauthenticated Prometheus endpoint. A support bundle is an explicit CLI export with bounded size and the same redaction policy.

OpenTelemetry export is deferred. The selected observability report found no implemented OpenTelemetry precedent across the studied CLIs, while structured logs and one local metrics precedent were well supported. Deferral does not change the stable correlation vocabulary, so a later optional exporter can consume it without changing run API resources.

### 11. Verification contract

The API is accepted only with all of the following behavior-level tests:

- route/method and exact compatibility fixtures for old operation APIs plus additive run routes
- JSON schema tests for full, compacted, tombstone, active, stalled, interrupted, cleanup-uncertain, and terminal resources
- strict input, identifier, pagination, cursor, unknown-field, event-size, and output-redaction tests
- read visibility across browser-session rotation and mutation rejection without current CSRF/session authority
- start deduplication before and after response loss, including a retry through another local server
- fail-closed acceptance when durable creation fails and proof that no child starts
- committed-event ordering, duplicate delivery, bounded replay reconnect, cursor-ahead, typed replay-gap with snapshot, slow subscriber, and observer restart tests
- idempotent duplicate cancellation, stale/unreachable owner, cancellation acknowledgement, and completion/cancellation/timeout/reconciliation terminal races
- legacy `op_*` resolution, operation-to-run projection, tombstones, event expiry, unsupported persisted schema, and rollback/migration fixtures
- two concurrent CLI runs visible through `GET /api/v1/operations`, `GET /api/v1/runs`, browser top bar, TUI, and CLI JSON
- two local servers observing the same workspace, including retained replay followed by newly committed events
- `go test ./...`, `go test -race ./...`, browser tests, CLI build, and gated real-runtime dogfood

Use injected clocks, process identities, repository failures, subscriber speed, and cancellation coordinators. Assertions target public snapshots, events, statuses, error codes, and absence of child starts, not internal map shape or goroutine counts.

## Trade-Offs

| Decision | Benefit | Cost and rejected alternative |
| --- | --- | --- |
| Add `/runs` and retain `/operations` | Gives the durable model a coherent resource without breaking the frozen browser contract. | Two projections must be maintained. Rejected: expanding `operationDocument` into the domain model, because its field order and meanings are compatibility fixtures. |
| Workspace-readable history with separately authorized mutation | Survives session expiry and supports every local observer. | Any local process able to reach the loopback server can see sanitized run metadata. Rejected: originating-session object ownership, because it recreates the current defect. |
| Same-host multi-process guarantee | Matches loopback security, process identity, and existing local coordination evidence. | Does not promise shared-filesystem multi-host cancellation or liveness. Rejected: implying distributed guarantees without a remote identity, clock, or worker protocol. |
| Opaque run and attempt IDs | Stable correlation without leaking path, PID, timestamp, or provider identity. | IDs cannot be sorted or interpreted by clients. Rejected: reusing agentwrap, product-run, PID, or timestamp-derived IDs because each has narrower authority or reuse hazards. |
| Typed snapshot plus separate event journal | Polling, replay, tombstones, and current truth remain useful after compaction. | Snapshot/event consistency and migration need explicit repository transactions. Rejected: SSE or in-memory event buffers as source of truth. |
| At-least-once sequence delivery | Reconnect is implementable and duplicate-safe across observers. | Clients must deduplicate. Rejected: exactly-once network delivery, which would add acknowledgement state without product value. |
| JSON `409 replay_gap` before SSE | Makes unavailable history machine-readable and includes a current snapshot. | Browser EventSource needs preflight/fallback logic. Rejected: silently starting from oldest retained history or injecting non-durable synthetic sequence events. |
| Bounded 512-event replay connections | Prevents one subscriber from monopolizing memory or a server goroutine. | Large catch-up requires transparent reconnects. Rejected: unbounded replay on one connection. |
| Durable token-digest deduplication | A lost `202` can be retried through another observer without duplicate work or stored secrets. | Dedup history consumes bounded index space. Rejected: process-local dedup only and caller-selected mutable run IDs. |
| `DELETE` as idempotent cancel request | Preserves the current route, makes repeated calls safe, and separates command from observation. | `202` means requested, not necessarily delivered or terminal. Rejected: treating SSE disconnect as cancel or returning success before request persistence. |
| Tombstone as a `200` run representation | Stable links retain useful identity and recovery after detail compaction. | Clients must inspect `record_state`. Rejected: using generic `404` for every missing detail case. |
| Stable error envelope with typed codes/details | Supports browser recovery and automation while preserving safe causes. | More codes require compatibility maintenance. Rejected: substring classification and string-only failures. |
| CLI diagnostics plus additive health summary | Keeps detailed telemetry local and explicit while allowing browser health display. | No scrape-ready Prometheus or traces in this sprint. Rejected: exposing high-cardinality or sensitive run data on a broad metrics endpoint. |

## Evidence

### Repository and contract evidence

- `projects/ultraplan-go/sprints/35-durable-run-observability/code-context.md` reports that `internal/web/operations.go` currently owns IDs, session filtering, events, subscribers, cancellation, counts, and terminal state in one process-local map. The API decision to make `/runs` canonical and `/operations` a projection is an inference from that defect plus the requirement that `internal/web` remain an adapter.
- `../ultraplan-go/internal/web/api_compatibility_test.go` freezes the `/api/v1/operations` route/method matrix and exact operation DTO field order. This directly supports an additive route instead of incompatible in-place replacement.
- `../ultraplan-go/internal/web/operations_contract_test.go` freezes lifecycle values, terminal classification, stable SSE names, frame shape, and error codes consumed by the embedded browser. This supports explicit compatibility mapping for new states and events.
- `../ultraplan-go/internal/web/operation_handlers.go` currently returns `202` with an operation `Location`, filters reads/cancel/events by browser session, parses `Last-Event-ID`, uses strict JSON decoding, and collapses missing records into `operation_not_found`. The proposed start, visibility, cursor, and typed recovery contracts address those specific seams.
- `../ultraplan-go/internal/web/operations.go` already bounds event count, bytes, event size, subscribers, queues, and stream lifetime, and disconnects slow subscribers. It also reveals that sequencing and deduplication disappear with the process. The API preserves bounded delivery while moving sequence authority and deduplication to durable state.
- `../ultraplan-go/internal/web/static/app.js` creates EventSource from the operation ID, understands the frozen event names, reconnects automatically, and reloads on terminal. This supports compatibility projection plus a canonical-client preflight for typed gaps.
- `projects/ultraplan-go/docs/ARCHITECTURE.md`, `projects/ultraplan-go/docs/PRD.md`, and `projects/ultraplan-go/docs/TRD.md` consistently separate durable operational truth from product artifacts, place web behind app use cases, keep loopback as the security boundary, and already describe additive `/api/v1/runs` capabilities. The API shape follows those project constraints rather than selecting a storage technology.

### Selected report findings

- `studies/go-cli-study/reports/final/05-error-handling.md` finds that typed errors carrying structured data, preserved causes, programmatic classification, and separate user/operator rendering outperform flat strings. The run error code table and safe details are the sprint-specific application of that report; the exact codes are this document's inference.
- `studies/go-cli-study/reports/final/07-state-context.md` finds that context cancellation and explicit persistent session/job identity solve different concerns, and that lock refresh plus a cleanup context can outlive work cancellation. This supports durable run identity independent from request/session context and cancellation acknowledgement independent from terminal outcome. The report explicitly finds no distributed session coordination, supporting the same-host boundary rather than a multi-host claim.
- `studies/go-cli-study/reports/final/08-concurrency.md` finds that localized lifecycle ownership, bounded queues/semaphores, explicit waiting, timeout cleanup, and `sync.Once` prevent leaks and duplicate cleanup. This supports bounded replay, slow-subscriber disconnection, idempotent cancellation, and finite stream lifetimes.
- `studies/go-cli-study/reports/final/10-logging-observability.md` finds structured fields, stable keys, output separation, and runtime debug controls to be the strongest CLI observability pattern. It finds OpenTelemetry absent and first-class metrics rare outside rclone. This supports the shared correlation vocabulary, local diagnostics command, additive health summary, and trace-export deferral.
- `studies/go-cli-study/reports/final/11-testing-strategy.md` finds that real command-path integration, golden/fixture contracts, centralized fakes, and behavior-focused assertions provide the most trustworthy compatibility coverage. This supports route/schema fixtures, process-level multi-server tests, injected failure seams, and avoiding implementation-detail assertions.
- `studies/go-cli-study/reports/final/12-extensibility.md` finds explicit versioned metadata and conversion functions safer than unversioned schema drift, while warning that formal extension systems add substantial lifecycle and compatibility burden. This supports separate API and persisted schema versions, additive v1 evolution, and rejecting a plugin/remote-worker protocol for run observation.
- `studies/go-cli-study/reports/final/13-security.md` finds that explicit trust boundaries, schema validation, secret-aware types, credential scrubbing, bounded diagnostics, and permission separation materially improve safety. This supports strict input validation, read/mutation policy separation, token-digest deduplication, redaction before persistence/fan-out, and omission of owner PID/raw provider data.
- `studies/go-cli-study/reports/final/14-performance.md` finds that streaming, bounded structures, bounded concurrency, and disk-backed long-session state keep long-lived tools viable, while noting the latency and lock costs of persistence. This supports cursor streaming, bounded replay and event size, paginated lists, and snapshots that remain useful after compaction. It does not by itself select SQLite or a filesystem journal.

### Requirement mapping

| Sprint outcome | API decision |
| --- | --- |
| Durable identity before child start | Server-issued opaque run ID; `202` only after durable acceptance; typed `503` otherwise. |
| Workspace-wide discovery | Canonical run list plus compatibility active projection without session filtering. |
| Stable inspection | Full/compacted/tombstone snapshots with lifecycle, liveness, attempts, progress, safe diagnostics, result, and recovery. |
| Replayable delivery | Committed decimal sequence, `Last-Event-ID`/`after`, at-least-once delivery, bounded replay, typed gap plus snapshot. |
| Ownership and reconciliation | Safe liveness projection and immutable terminal winner; no caller-controlled owner/fencing facts. |
| Cross-surface cancellation | Authorized idempotent `DELETE`, durable request before routing, visible acknowledgement/uncertainty. |
| Compatibility and migration | Additive run routes, frozen operation projection, durable aliases, honest legacy `410`, retained tombstones. |
| Security and data minimization | Loopback-only policy, workspace-relative targets, pre-append redaction, bounded safe DTOs, no raw payloads. |
| Operational telemetry | Stable correlation fields, additive run-control health, CLI diagnostics, no mandatory OTEL. |
| Fault verification | Schema, auth, dedupe, replay, retention, owner, race, storage, browser, process, and real-runtime tests. |

## Risks

- Some existing JSON consumers may reject unknown routes or lifecycle values even though the operation DTO remains unchanged. Compatibility fixtures must include a strict old-client decoder and the old browser bundle.
- Changing `GET /api/v1/operations` from session-local to workspace-wide intentionally reveals sanitized run metadata to other local browser sessions. The loopback/same-origin boundary and redaction policy must be tested as security invariants, not assumed.
- EventSource does not expose a non-200 JSON body to JavaScript. The canonical browser client must preflight saved cursors and fetch run detail after connection errors; otherwise typed `replay_gap` information is present on the wire but not useful in the UI.
- Closing replay after 512 events relies on standards-compliant EventSource reconnection and correct `Last-Event-ID` propagation. Browser tests must prove no gap or live-event overtaking during repeated catch-up connections.
- Durable token-digest deduplication only works while the dedupe/tombstone record is retained. The API reports expiry and does not promise indefinite safe replay of a start request.
- `record_state: "tombstone"` with `200` is useful for stable links but can surprise clients that assume every successful response has event history. Schema docs and clients must branch on `record_state`.
- Broader canonical lifecycle values can drift from the frozen operation vocabulary. Projection tests must prove every canonical state has a conservative old-client mapping and that `terminal` remains authoritative.
- Safe provider session, process identity, target, and diagnostics fields can still be sensitive in some workspaces. The allowlist must be field-based and redaction must occur before storage, logs, JSON, SSE, and support export.
- A store outage during reads must not fall back to a stale web-hub snapshot. This makes the UI unavailable rather than misleading, which is the selected safety trade-off.
- The API intentionally does not expose remote-host ownership or authorization. If final architecture reasoning were to select multi-host support, this area must be revisited for host identity, clock, transport security, and remote cancellation authorization rather than silently broadening the contract.
- Exact retention durations, disk quota, lease intervals, repository transactions, and compaction implementation are architecture/persistence decisions. This API decision is final about their observable semantics: clients receive explicit boundaries, health, tombstones, and typed failures and may never infer a fixed retention duration.

The area decision is to proceed with the split canonical/compatibility API. The main trade-off is maintaining two projections in exchange for durable semantics without breaking the existing local browser and v1 JSON clients.
