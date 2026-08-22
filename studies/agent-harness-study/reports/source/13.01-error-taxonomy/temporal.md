# Source Analysis: temporal

## Dimension 13.01: Error Taxonomy

### Source Info

| Field | Value |
|-------|-------|
| Name | temporal |
| Path | `studies/agent-harness-study/sources/temporal` |
| Language / Stack | Go (Temporal server: frontend/history/matching/worker services, gRPC, Cassandra/SQL persistence) |
| Analyzed | 2026-08-21 |

## Summary

The Temporal server implements a layered, typed error taxonomy organized primarily by **source layer** (database driver → persistence → service boundary → cross-service wire → client), rather than by LLM-agent categories like model/provider/tool. There is no single `ErrorType` enum; instead the taxonomy is expressed through four coordinated mechanisms:

1. **Typed Go error structs** per layer: persistence errors (`common/persistence/data_interfaces.go:102-158`), server-specific service errors (`common/serviceerror/*.go`), and the canonical API `serviceerror` package (external module `go.temporal.io/api`, documented in-repo at `docs/architecture/retry.md:26-30`).
2. **Wire-format classification**: gRPC status codes plus typed proto error details (`proto/internal/temporal/server/api/errordetails/v1/message.proto:12-61`) so typed errors survive round-trips between services.
3. **Cause/scope enums for resource exhaustion**: `serviceerror.ResourceExhausted{Cause, Scope}` with ~11 distinct causes observed in code (RPS_LIMIT, BUSY_WORKFLOW, APS_LIMIT, SYSTEM_OVERLOADED, PERSISTENCE_STORAGE_LIMIT, ...).
4. **Workflow-failure taxonomy**: a proto `Failure` union discriminated by failure-info kind — Application / Timeout / Server / Terminated / Canceled / Reset / NexusHandler — constructed via helpers in `common/failure/failure.go:14-46`.

The taxonomy is heavily load-bearing: task executors route retry/drop/DLQ decisions off error types (`service/history/queues/executable.go:520-598`), workflow retries are decided by failure-info kind and non-retryable flags (`service/history/workflow/retry.go:115-152`), and RPC retry policies branch on specific types (`common/util.go:329-347`). The direct dimension question — *can you tell from the error type whether to retry, escalate, or stop?* — is answered yes at every layer.

## Rating

**8 / 10**

Rationale against the rubric: this is a clear, explicit model with operational safeguards and observability (metrics tagged by cause, e.g. `metrics.ResourceExhaustedCauseTag` at `service/history/queues/executable.go:465-466`; internal-error masking at `common/rpc/interceptor/mask_internal_error.go:67-87`), proven under production failure modes (DLQ escalation, shard failover redirect). It falls short of 9-10 because of duplicated conversion logic requiring manual sync (`service/history/chasm_engine.go:1265-1266` "Keep in sync" note vs `service/history/handler.go:2233`), some string/sentinel-equality based classification (`service/history/queues/executable.go:488` compares `err.Error()` strings), persistence errors that discard wrapped causes (`common/persistence/data_interfaces.go:1358-1392` return only `Msg`), and no single authoritative taxonomy document in this repo (the canonical list lives in the external `go.temporal.io/api` module).

## Evidence Collected

| Area | Evidence | File:Line |
|------|----------|-----------|
| Persistence error type definitions | Typed structs: `InvalidPersistenceRequestError`, `AppendHistoryTimeoutError`, `CurrentWorkflowConditionFailedError`, `WorkflowConditionFailedError`, `ConditionFailedError`, `ShardAlreadyExistError`, `ShardOwnershipLostError`, `TimeoutError`, `TransactionSizeLimitError` | common/persistence/data_interfaces.go:102-158 |
| Conflict-class predicate used for retry routing | `IsConflictErr` switch over three condition-failed types | common/persistence/data_interfaces.go:1394-1402 |
| Server-specific typed service errors | `ShardOwnershipLost` struct carrying owner/current host + gRPC `Status()` with proto details | common/serviceerror/shard_ownership_lost.go:11-48 |
| Server-specific typed service errors (replication) | `RetryReplication` with event-range fields; encoded as `codes.Aborted` + `RetryReplicationFailure` detail | common/serviceerror/retry_replication.go:9-70 |
| Wire-level error-detail taxonomy | Proto messages: `TaskAlreadyStartedFailure`, `CurrentBranchChangedFailure`, `ShardOwnershipLostFailure`, `RetryReplicationFailure`, `SyncStateFailure`, `StickyWorkerUnavailableFailure`, `ObsoleteDispatchBuildIdFailure`, `ObsoleteMatchingTaskFailure`, `ActivityStartDuringTransitionFailure`, `StalePartitionCountsFailure` | proto/internal/temporal/server/api/errordetails/v1/message.proto:12-61 |
| Cross-service decode dispatch | `FromStatus` switches on gRPC code AND detail type to reconstruct typed server errors, falling back to generic `serviceerror.FromStatus` | common/serviceerror/convert.go:11-57 |
| Driver→persistence classification (infrastructure source) | Cassandra driver errors mapped to `persistence.TimeoutError`, `NotFound`, `ResourceExhausted{SYSTEM_OVERLOADED}`, `ResourceExhausted{PERSISTENCE_STORAGE_LIMIT}` (disk-threshold string match), else `Unavailable` | common/persistence/nosql/nosqlplugin/cassandra/gocql/errors.go:15-64 |
| SQL store classification | Transaction errors pass through known types; unknowns become `Unavailable`; `sql.ErrNoRows`→`NotFound`; serialization failures→`Internal` | common/persistence/sql/common.go:52-81,128-137 |
| Persistence→service conversion (history handler) | `convertError`: ShardOwnershipLost→redirect info, AppendHistoryTimeout/ConditionFailed→`Unavailable`, TransactionSizeLimit→`InvalidArgument`, Timeout→`DeadlineExceeded` | service/history/handler.go:2233-2258 |
| CHASM engine superset conversion | Same mapping plus `ConditionFailedError`; enriches NotFound with archetype display name; requires sync with handler version (comment) | service/history/chasm_engine.go:1261-1329 |
| Serialization→DataLoss normalization | Interceptor converts `serialization.DeserializationError`/`SerializationError` to `serviceerror.DataLoss` and truncates messages | common/rpc/interceptor/service_error_interceptor.go:37-51 |
| Internal-error masking (security safeguard) | Unknown/Internal codes masked as "something went wrong, please retry" + error hash when namespace flag enabled | common/rpc/interceptor/mask_internal_error.go:21,67-87 |
| Resource-exhaustion cause/scope classification | Pre-built `RateLimitServerBusy = ResourceExhausted{CAUSE_RPS_LIMIT, SCOPE_SYSTEM}` | common/rpc/interceptor/rate_limit.go:20-24 |
| Cause enum usage breadth | Cause constants counted across repo: RPS_LIMIT(25), BUSY_WORKFLOW(19), CONCURRENT_LIMIT(11), WORKER_DEPLOYMENT_LIMITS(9), SYSTEM_OVERLOADED(9), PERSISTENCE_LIMIT(6), APS_LIMIT(4), PERSISTENCE_STORAGE_LIMIT(3), OPS_LIMIT, CIRCUIT_BREAKER_OPEN | grep across repo (e.g. common/rpc/interceptor/rate_limit.go:20-24, service/history/workflow/update/update.go:409) |
| Busy-workflow policy error | Update rejection returns `NewResourceExhausted(CAUSE_BUSY_WORKFLOW)` telling caller to retry | service/history/workflow/update/update.go:409 |
| Task-level routing taxonomy | `HandleErr` pipeline: DLQ pattern match → `isInvalidTaskError` (drop) → `isSafeToDropError` (drop) → `isExpectedRetryableError` (retry w/ rewritten sentinel) → unexpected: terminal-check → DLQ-or-drop | service/history/queues/executable.go:517-598 |
| User-vs-system classification at task level | `isUserError`: ResourceExhausted scoped to NAMESPACE counts as user error (excluded from latency accounting); system scope does not | service/history/queues/executable.go:398-405 |
| Retryability by ResourceExhausted cause | Switch on `Cause`: BUSY_WORKFLOW→sentinel, APS_LIMIT→count+reschedule policy, default throttled; metrics tagged per cause | service/history/queues/executable.go:453-467 |
| Terminal/corruption escalation | `isUnexpectedNonRetryableError`: terminal task errors, `DataLoss`, `serviceerror.Internal` → DLQ if enabled else drop, with corruption metric | service/history/queues/executable.go:496-515,571-585 |
| Workflow retry decision from failure kind | `isRetryable(failure)`: Terminated/Canceled→stop; Timeout START_TO_CLOSE/HEARTBEAT→honors `TemporalTimeout:*` non-retryable types; other timeouts→stop; ServerFailure→NonRetryable flag; ApplicationFailure→type-string + NonRetryable flag | service/history/workflow/retry.go:115-152 |
| Server-failure construction with retry flag | `failure.NewServerFailure(msg, nonRetryable)` used for size-limit failures (e.g. heartbeat/result exceeds limit, NonRetryable=true) | common/failure/failure.go:14-23; service/frontend/workflow_handler.go:1290,1468,1667 |
| Timeout failure construction | `failure.NewTimeoutFailure(msg, timeoutType)` sets `Source="Server"`; activity timer path converts exhausted-retry to SCHEDULE_TO_CLOSE | common/failure/failure.go:36-46; service/history/timer_queue_active_task_executor.go:303-319 |
| Workflow-task-failed cause classification | `newWorkflowTaskFailedCause(WORKFLOW_TASK_FAILED_CAUSE_BAD_BINARY, InvalidArgument)` pairs cause enum with service error | service/history/api/respondworkflowtaskcompleted/api.go:364-371 |
| RPC retry routing on types | History handler retry: stale-state sentinels, conflict errors, `IsServiceHandlerRetryableError` (Internal/Unavailabile/MultiOperation recursion) | service/history/api/retry_util.go:12-18 |
| Client transient-error classification | `IsServiceClientTransientError`: namespace-scoped ResourceExhausted NOT transient; ShardOwnershipLost/StalePartitionCounts are | common/util.go:311-327 |
| Context-source classification | `IsContextDeadlineExceededErr`/`IsContextCanceledErr` unify context errors with `serviceerror.DeadlineExceeded`/`Canceled` | common/util.go:300-309 |
| Special throttle policy for ResourceExhausted | `ThrottleRetryContext` takes max of normal backoff and dedicated 1s-10s throttle retrier when err is `*serviceerror.ResourceExhausted` | common/backoff/retry.go:38-48,80-82 |
| Sentinel taxonomy for task engine | `ErrTaskDiscarded`, `ErrTaskVersionMismatch`, `ErrTaskRetry`, `ErrDependencyTaskNotCompleted`, `ErrDuplicate`, `ErrStaleReference`, `ErrStaleState`, quota sentinels etc. | service/history/consts/const.go:20-93 |
| Retry-policy validation ties taxonomy to user config | `NonRetryableErrorTypes` validated; `TimeoutFailureTypePrefix = "TemporalTimeout:"` convention | common/retrypolicy/retry_policy.go:14-20,86-94 |
| Taxonomy documentation | Service Error interface, api-go serviceerror package (general-purpose + specialized), server-specific additions, retry wiring for handler/client | docs/architecture/retry.md:14-45 |
| Conversion round-trip test | `TestFromToStatus` encodes/decodes `ShardOwnershipLost` via status details | common/serviceerror/convert_test.go:11-26 |
| Routing behavior tests | `executableSuite.TestExecuteHandleErr_*` family exercises HandleErr drop/retry/DLQ branches | service/history/queues/executable_test.go:442-594 |

## Answers to Dimension Questions

### 1. Are errors classified by source?

Yes, but the primary axis is **infrastructure layer**, not agent-style sources (model/provider/tool). Concretely: driver-level faults are classified at the DB boundary into `persistence.TimeoutError` / `ResourceExhausted` / `Unavailable` (`common/persistence/nosql/nosqlplugin/cassandra/gocql/errors.go:15-64`); persistence-layer conflicts have their own structs (`common/persistence/data_interfaces.go:113-148`); validation failures are uniformly `InvalidArgument` (e.g. `service/history/handler.go:161-169`); timeouts form their own category both as persistence errors and as workflow `TimeoutFailureInfo` keyed by `TimeoutType` (`common/failure/failure.go:36-46`). Within ResourceExhausted there is a second axis — **scope** (system vs namespace) — which is explicitly treated as user-vs-system at the task executor (`service/history/queues/executable.go:398-405`). User application failures (closest analog to "tool/model output" failures) are carried as `ApplicationFailureInfo` with a free-form `Type` string and `NonRetryable` flag (`service/history/workflow/retry.go:141-150`). Nexus failures carry an explicit source header distinguishing worker-originated vs server-originated failures (`common/nexus/failure.go:22-27`).

### 2. Is the taxonomy used for handling?

Extensively. The clearest example is the task-executor decision pipeline `HandleErr` (`service/history/queues/executable.go:517-598`) which maps error identity to one of: complete-and-drop (invalid/safe-to-drop), retry with rewritten sentinel (expected: ResourceExhausted by cause, NamespaceNotActive, handover), DLQ escalation or silent drop (terminal/DataLoss/Internal), or bounded-unexpected retry then DLQ (`executable.go:555-557,588-595`). At the workflow level, `isRetryable` decides continue-as-new retry vs terminal failure purely from the serialized failure-info kind and flags (`service/history/workflow/retry.go:115-152`). RPC layers pick retry policies by type (`common/backoff/retry.go:80-82` gives ResourceExhausted its own throttle policy; `common/util.go:329-347` marks Internal/Unavailable handler-retryable). Shard ownership loss triggers host redirection rather than blind retry (`service/history/handler.go:2239-2244`).

### 3. Are error categories documented?

Partially. `docs/architecture/retry.md:14-33` documents the ServiceError interface concept, points to the api-go `serviceerror` package as the canonical general-purpose set, and names the server-specific additions (`ShardOwnershipLost`, `TaskAlreadyStarted`). The errordetails proto carries inline comments explaining each wire error (`proto/internal/temporal/server/api/errordetails/v1/message.proto:3,48-61`), and `AGENTS.md` documents severity conventions for logging (`logger.Fatal` for invariant violations, `DPanic` otherwise). However, there is no single document enumerating the full taxonomy, the retry/drop semantics per category, or the persistence↔service mapping table; that knowledge is distributed between `docs/architecture/retry.md`, proto comments, and the routing code itself.

### 4. Can new error types be added without breaking existing handling?

Largely yes, by design. The extension recipe is additive: define a proto detail message, implement a Go struct with `Status()` attaching it (`common/serviceerror/shard_ownership_lost.go:35-48`), add a reconstruction case in `FromStatus` (`common/serviceerror/convert.go:18-54`), and rely on the fallback `serviceerror.FromStatus(st)` (`convert.go:56`) so older receivers still get a valid generic service error. Type switches in routing code use default branches rather than exhaustive enums (e.g. `switch resourceExhaustedErr.Cause //nolint:exhaustive` with a default at `service/history/queues/executable.go:455-463`), so new causes fall into sane throttled-retry behavior. Two caveats: (a) the two `convertError` implementations must be manually kept in sync (`service/history/chasm_engine.go:1265-1266`), so a new persistence error added to only one path will silently degrade to pass-through on the other; (b) several routing checks compare sentinels or message strings (`err == consts.ErrTaskVersionMismatch` at `executable.go:426`, `err.Error() == consts.ErrNamespaceHandover.Error()` at `executable.go:488`, regex-on-message DLQ patterns at `executable.go:600-621`), which new errors cannot participate in without knowing those conventions.

## Architectural Decisions

- **Layer-owned error definitions with upward conversion.** Each layer defines its own error vocabulary and each boundary has an explicit converter (driver: `gocql/errors.go:15-64`; SQL: `sql/common.go:52-81`; history handler: `handler.go:2237-2258`; CHASM: `chasm_engine.go:1267-1314`). Errors never leak raw driver types to clients.
- **Wire-preserving typed errors.** Instead of flattening to a code+message, typed server errors serialize rich fields into gRPC status details (`common/serviceerror/retry_replication.go:52-70`) and are reconstructed client-side (`common/serviceerror/convert.go:11-57`), keeping the taxonomy intact across process boundaries.
- **Dual-axis resource-exhaustion model.** `ResourceExhausted{Cause, Scope}` separates *what limit* was hit from *who hit it*, enabling distinct behaviors: namespace-scoped = user fault (no latency accounting, not transient for clients: `common/util.go:317-319`), system-scoped = server fault (throttle-retried: `common/backoff/retry.go:80-82`).
- **Persistence of failure taxonomy, not just errors.** Workflow failures are stored as a discriminated proto union (`ApplicationFailureInfo`/`TimeoutFailureInfo`/`ServerFailureInfo`...) so retry decisions survive restarts and are evaluated later from history (`service/history/workflow/retry.go:115-152`), including truncation that deliberately preserves the NonRetryable flags (`common/failure/failure.go:69-81`).
- **Fail-open reconstruction.** `FromStatus` falls back to the generic api-go parser for unrecognized codes/details (`common/serviceerror/convert.go:56`), so forward-compatible evolution of the taxonomy is possible.

## Notable Patterns

- **Classifier-predicate routing**: small boolean predicates over error types (`isUserError`, `isInvalidTaskError`, `isSafeToDropError`, `isExpectedRetryableError`, `isUnexpectedNonRetryableError`) compose into a single decision function (`service/history/queues/executable.go:398-515`), making routing policy inspectable and testable.
- **Sentinel + wrapper escalation**: terminal failures wrap the original in `fmt.Errorf("%w: %v", ErrTerminalTaskFailure, err)` so downstream `Execute` can detect terminality while retaining the cause (`service/history/queues/executable.go:579-581`).
- **Observability keyed by taxonomy**: every routing branch emits a distinct metric (`TaskSkipped`, `TaskThrottledCounter` with cause tag, `TaskCorruptionCounter`, `TaskTerminalFailures`), turning the taxonomy into dashboards (`service/history/queues/executable.go:412,465-466,508,575,580`).
- **Defense-in-depth normalization at the edge**: serialization failures become `DataLoss` before leaving the process (`common/rpc/interceptor/service_error_interceptor.go:37-42`), and Internal/Unknown responses can be masked with a stable hash for support correlation (`common/rpc/interceptor/mask_internal_error.go:80-86`).
- **Convention-encoded sub-taxonomy**: timeout failure types embed in the user-facing `NonRetryableErrorTypes` list under a reserved prefix `"TemporalTimeout:"`, validated at policy admission (`common/retrypolicy/retry_policy.go:14-20,86-94`).

## Tradeoffs

- **Rich typed errors vs boilerplate**: every cross-service typed error needs a proto detail message, Status() plumbing, and a FromStatus case (compare `shard_ownership_lost.go` and `retry_replication.go` structures). The payoff is lossless routing; the cost is that most server-specific errors duplicate ~60 lines of encode/decode scaffolding.
- **Two conversion implementations**: legacy `Handler.convertError` and the CHASM superset coexist with a manual-sync contract (`chasm_engine.go:1265-1266`); the legacy path already lacks `ConditionFailedError` handling that CHASM has (`handler.go:2247-2250` vs `chasm_engine.go:1296-1299`).
- **String-based escape hatches**: regex DLQ matching on error text is operationally flexible (configurable via dynamicconfig key `HistoryTaskDLQErrorPattern`, `executable.go:600-621`) but couples behavior to message wording; similarly the Cassandra disk-limit check matches on lowercased message text (`gocql/errors.go:47`).
- **Information loss at the persistence layer**: persistence error structs carry only `Msg` strings without wrapping the underlying driver error (`data_interfaces.go:1358-1392`), trading debuggability inside persistence for clean abstraction above it.

## Failure Modes / Edge Cases

- **Drift between duplicate converters**: a new persistence error handled in CHASM but not in the legacy handler passes through unconverted (`service/history/handler.go:2257` returns `err` unchanged), potentially surfacing a non-gRPC-safe error to the wire.
- **First-attempt-only invalid-task detection**: `isInvalidTaskError` only marks a task invalid on attempt 1, because later attempts cannot distinguish pre-existing invalidity from a failed write in a prior attempt (`service/history/queues/executable.go:537-542`) — a subtle correctness guard against misclassifying self-inflicted errors.
- **Guard against classifier bugs**: `isExpectedRetryableError` uses a deferred check preventing `(true, nil)` from silently completing a task (`service/history/queues/executable.go:443-451`).
- **Corruption containment**: DataLoss/Internal errors are treated as likely data corruption and either DLQ'd or dropped-with-metric rather than endlessly retried (`executable.go:496-515,571-584`).
- **Retry multiplication**: the docs call out that handler-side and client-side retries multiply; handler retries are capped at one extra attempt (`docs/architecture/retry.md:52-55`).
- **Namespace handover window**: handover errors are explicitly non-retryable for handlers (`common/util.go:330-332`) but surfaced as expected retryables for standby tasks (`executable.go:488-491`), showing the same category gets opposite treatment per consumer.

## Future Considerations

- Consolidate the two `convertError` implementations behind one shared mapper parameterized by extras (request-ID enrichment, NotFound archetypes) to eliminate the keep-in-sync risk (`service/history/chasm_engine.go:1265`, `service/history/handler.go:2233`).
- Replace string/sentinel-equality checks with typed or `errors.Is`-wrapped errors (e.g. `executable.go:426,488`) so new errors integrate predictably.
- Wrap causes inside persistence errors while preserving the public surface, restoring stack traces for infra debugging (`common/persistence/data_interfaces.go:1358-1392`).
- Produce a single taxonomy reference table (category → representative types → retry/drop/DLQ semantics → metric) extending `docs/architecture/retry.md`.

## Questions / Gaps

- The canonical general-purpose service error definitions (`InvalidArgument`, `NotFound`, `ResourceExhausted`, ...) live in the external `go.temporal.io/api` module (declared at `go.mod`, pinned `v1.62.15...`) and are not inspectable within this repository; analysis of that half of the taxonomy is limited to usage sites and `docs/architecture/retry.md:26-31`. No evidence found in-repo of the full enumerated list.
- No evidence found of automated tests asserting that *every* typed server error round-trips through `Status()`/`FromStatus`; coverage appears limited to `ShardOwnershipLost` (`common/serviceerror/convert_test.go:11-26`) plus scattered interceptor tests (e.g. `common/rpc/interceptor/service_error_interceptor_test.go` exists but was not exhaustively reviewed).
- Search boundary: the dimension's "model/provider" source categories have no direct analog; nearest mappings identified were worker/application failures (`ApplicationFailureInfo`) and Nexus worker-originated failures (`common/nexus/failure.go:25-26`). If finer-grained worker-side error taxonomies matter to this study, they reside in the separate SDK repositories, outside this source's isolation boundary and therefore not inspected.

---

Generated by `13.01-error-taxonomy` against `temporal`.
