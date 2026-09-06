# Source Analysis: temporal

## Event Envelope, Ordering, Causation, and Replay

### Source Info

| Field | Value |
|-------|-------|
| Name | temporal |
| Path | `studies/ultraplan-daemon-events-study/sources/temporal` |
| Language / Stack | Go (server) + protobuf-defined HistoryEvent schema |
| Analyzed | 2026-09-02 |

## Summary

Temporal models durable events as a per-workflow append-only history that is the
event log, the source of truth for state-machine reconstruction, and the primary
replication payload across clusters. The server distinguishes:

- **Event identity** = monotonic `EventId` (`int64`) plus `Version` (`int64`),
  stamped at `EventFactory.createHistoryEvent` (`service/history/historybuilder/event_factory.go:1085-1096`).
- **Operation/transition identity** = `VersionedTransition{NamespaceFailoverVersion, TransitionCount}`
  on `WorkflowExecutionInfo.TransitionHistory` (`api/persistence/v1/executions.pb.go:275`,
  `api/persistence/v1/hsm.pb.go:457-465`), used as a cluster-wide logical clock
  for staleness checks (`common/persistence/transitionhistory/transition_history.go:51-69, 81-125`).
- **Attempt identity** = `attempt` field on individual command attributes
  (`service/history/historybuilder/event_factory.go:112-119`, etc.); not a
  generic envelope field.

Sequence assignment is in-memory by `EventStore.AllocateEventID`
(`service/history/historybuilder/event_store.go:58-62`), starting from the
shard's `nextEventID` for mutable state, then committed in batches through
`AppendHistoryNodesRequest{PrevTransactionID, TransactionID}` — the DB-side
concurrency control overwrites (`common/persistence/data_interfaces.go:776-792`,
`common/persistence/data_interfaces.go:788-791`). On read, the public surface
returns an opaque continuation token
(`api/token/v1/message.pb.go:30-42` — `HistoryContinuation{FirstEventId, NextEventId, PersistenceToken, BranchToken, VersionHistoryItem, VersionedTransition}`)
that lets a caller resume without gaps or duplicates.

The event_type enum lives in the `go.temporal.io/api` dependency (the
HistoryEvent oneof is keyed on past-tense names like
`EVENT_TYPE_WORKFLOW_EXECUTION_COMPLETED`). Causation is not a per-event
attribute; it is encoded through `scheduledEventId`/`initiatedEventId` pointer
fields plus the request-id fields wired by `EventStore.wireEventIDs`
(`service/history/historybuilder/event_store.go:381-454`).

Replay: the API is `GetWorkflowExecutionHistory(WaitNewEvent=true)` which, after
returning the current page, long-polls via
`events.Notifier.WatchHistoryEvent`/`NotifyNewHistoryEvent`
(`service/history/events/notifier.go:128-153, 246-248`) gated by per-workflow
channels and then re-checks mutable state before returning
(`service/history/api/get_workflow_util.go:128-260`).

## Rating

**8 / 10** — Clear model with explicit interfaces, transactional invariants,
strong staleness and branch-token checks, and operational safeguards for
in-flight resumptions. Below mature because (a) no first-class causation field
on the envelope (causality lives in scattered pointer fields), (b) notification
fanout has a non-blocking channel with silent drop semantics under load, and
(c) the public continuation token couples the consumer to a server-versioned
opaque blob with both per-event and storage-internal cursors.

## Evidence Collected

| Area | Evidence | File:Line |
|------|----------|-----------|
| Envelope schema (event) | HistoryEvent oneof attributes; EventId/Version/EventTime/TaskId set in factory | `service/history/historybuilder/event_factory.go:1085-1096` |
| Envelope fields per event type | Past-tense names; EventType enum in `go.temporal.io/api/enums/v1`; server factory in `event_factory.go:39,112,136,166,190,213,234,264,285,304,323,340,357,373,389,413,438,454,466,481,512,524,538,555,574,587,604,619,649,677,700,727,741,761,787,807,829,851,893,917,942,967,992,1015,1038,1058,1106` | `service/history/historybuilder/event_factory.go:39-1106` |
| Sequence assignment | `AllocateEventID` increments `nextEventID`; batches rolled at max bytes | `service/history/historybuilder/event_store.go:58-66, 120-132` |
| Buffer + flush ordering | Cassandra reorder fix-up before commit; `reorderBuffer` puts "reorderBuffer" events last | `service/history/historybuilder/event_store.go:161-200, 461-491` |
| Causation wiring | `wireEventIDs` maps scheduled/started/initiated IDs; request IDs to event IDs | `service/history/historybuilder/event_store.go:381-454` |
| Request-id causation | `EVENT_TYPE_WORKFLOW_EXECUTION_OPTIONS_UPDATED.AttachedRequestId` and `EVENT_TYPE_WORKFLOW_EXECUTION_SIGNALED.RequestId` recorded | `service/history/historybuilder/event_store.go:441-452` |
| Operation transition identity | `VersionedTransition` proto | `api/persistence/v1/hsm.pb.go:456-465` |
| Staleness / causation-equivalence check | Compares NFV + TransitionCount; rejects stale references | `common/persistence/transitionhistory/transition_history.go:51-69, 81-125` |
| Branch/version history | `VersionHistoryItem`, `VersionHistory`, `VersionHistories`, `StrippedHistoryEvent` protos | `api/history/v1/message.pb.go:73-79, 126-132, 179-185, 336-342, 388-393` |
| Cross-branch ordering / LCA | `FindLCAVersionHistoryItem`, `SplitVersionHistoryByLastLocalGeneratedItem`, `ContainsVersionHistoryItem` | `common/persistence/versionhistory/version_history.go:95-137, 185-196` |
| Append durability | `AppendHistoryNodesRequest` with `PrevTransactionID`/`TransactionID` chain + "larger wins" tie-break | `common/persistence/data_interfaces.go:776-792` |
| Read pagination | `ReadHistoryBranchRequest{MinEventID, MaxEventID, PageSize, NextPageToken}` | `common/persistence/data_interfaces.go:820-847` |
| Page-boundary validation | `VerifyHistoryIsComplete` rejects incomplete results, emits `ServiceErrIncompleteHistoryCounter` | `service/history/api/get_history_util.go:216-229` |
| Gap-fill on pagination race | Re-queries mutable state to detect events committed between paginated DB fetches | `service/history/api/getworkflowexecutionhistory/api.go:456-480` |
| Public continuation token | `HistoryContinuation` carries FirstEventId, NextEventId, PersistenceToken, BranchToken, VersionHistoryItem, VersionedTransition | `api/token/v1/message.pb.go:30-42` |
| Public API request fields | `GetWorkflowExecutionHistoryRequest{MaximumPageSize, NextPageToken, WaitNewEvent, HistoryEventFilterType}` | external `go.temporal.io/api/workflowservice/v1/request_response.pb.go:1157-1252` |
| Resume / long-poll API | `GetWorkflowExecutionHistory` handler delegates to history service | `service/history/handler.go:1941-1959`, `service/history/history_engine.go:1105-1110` |
| Long-poll subscribe | `GetWorkflowExecutionHistoryReverse` and Watch in `get_workflow_util.go` | `service/history/api/get_workflow_util.go:128-260` |
| Notification envelope | `Notification{ID, LastFirstEventID, LastFirstEventTxnID, NextEventID, PreviousStartedEventID, Timestamp, WorkflowState, WorkflowStatus, VersionHistories, TransitionHistory}` | `service/history/events/notifier.go:37-48` |
| Notification flow | Single-process shard notifier, sub-channel size 1, fan-out drops on overflow | `service/history/events/notifier.go:120, 180-200, 202-212` |
| Engine emit | `NotifyNewHistoryEvent` called after transaction commit | `service/history/history_engine.go:925-930`; `service/history/workflow/transaction_impl.go:678, 718` |
| Replay contract on resume | `GetMutableStateWithConsistencyCheck` compares branch token + VersionedTransition against caller ref | `service/history/api/get_workflow_util.go:297-348` |
| Out-of-order delivery on subscribe | Notifier may send out-of-date event; caller compares and continues | `service/history/api/get_workflow_util.go:218-221` |
| Cross-cluster replication ordering | Local-vs-remote split by `SplitVersionHistoryByLastLocalGeneratedItem` | `service/history/replication/eventhandler/history_events_handler.go:114-148` |
| Branch-token enforcement | `ValidateBranchTokenForExecution` used on resume; foreign branch handled specially | `service/history/api/getworkflowexecutionhistory/api.go:260-272` |
| Backwards-compat at replay | `FixFollowEvents` rewrites last event to ContinuedAsNew for older SDKs | `service/history/api/get_history_util.go:554-638` |
| Schema evolution affordances | `reserved 1; reserved 2;` on `TransientWorkflowTaskInfo`; deprecated `EventType` enum field reordering (`event_id=1, version=4`) | `proto/internal/temporal/server/api/history/v1/message.proto:11-12`; `api/history/v1/message.pb.go:338-339` |
| Test coverage | `TestGetWorkflowExecutionHistory`, `TestGetWorkflowExecutionHistoryReverse_BranchTokenNotOwnedByExecution`, `TestGetWorkflowExecutionHistory_LongPollDiscardsRequestBranchToken` | `service/history/history_engine_test.go:5773, 7129, 7158, 7194, 7220` |
| Notifier tests | Drop-on-overflow semantics and long-poll semantics are tested | `service/history/events/notifier_test.go:85-150` |
| Ordering metrics | `OutOfOrderBufferedEventsCounter` emitted when a Completed/Started pair arrives out of order | `service/history/historybuilder/event_store.go:493-549` |

## Answers to Dimension Questions

1. **What ordering is actually guaranteed?**
   - Per workflow run, total order by `EventId`. Event IDs are allocated
     monotonically in the shard's history builder
     (`service/history/historybuilder/event_store.go:58-62, 85-106`) and persisted
     atomically in a transactional write to the events table (batched).
   - Per branch, the storage layer enforces `TransactionID` ordering with "larger
     wins" semantics (`common/persistence/data_interfaces.go:788-791`).
   - Across branches in the same run, `VersionHistory` provides a sequence of
     `VersionHistoryItem{event_id, version}` breakpoints that identify which
     `EventId` range belongs to which cluster/generation
     (`api/history/v1/message.pb.go:73-132`,
     `common/persistence/versionhistory/version_history.go:96-108`).
   - Cross-cluster ordering uses `VersionedTransition{NamespaceFailoverVersion, TransitionCount}`
     for staleness checks (`common/persistence/transitionhistory/transition_history.go:51-125`).
   - Buffered events that arrive out of order (Cassandra reorder) are detected
     and reordered before commit; out-of-order completion-before-started pairs
     are explicitly counted via
     `metrics.OutOfOrderBufferedEventsCounter`
     (`service/history/historybuilder/event_store.go:183-185, 493-549`).
2. **Can a reconnecting client resume without gaps or duplicates?**
   - Yes for the public `GetWorkflowExecutionHistory` reader: the server returns
     a continuation token bundling `FirstEventId`, `NextEventId`, and the
     storage cursor (`api/token/v1/message.pb.go:30-42`).
   - The server explicitly defends against races between DB fetches in
     pagination by re-reading mutable state and bridging gaps with
     `fetchGapEvents` (`service/history/api/getworkflowexecutionhistory/api.go:316-338, 456-480`).
   - Stale branch tokens or versioned transitions are detected on resume and
     turned into `CurrentBranchChanged` errors with the new branch token
     (`service/history/api/get_workflow_util.go:155-179, 218-244`).
   - The replication fetcher exposes explicit inclusive/exclusive event-id
     windows so cross-cluster replay does not double-deliver
     (`service/history/replication/eventhandler/remote_history_paginated_fetcher.go:30-53`).
   - Caveat: `Notifier` is **best-effort, non-blocking**. Subscribers can miss
     notifications if their channel buffer (size 1) is full or if the
     queue/dispatcher overflow (`service/history/events/notifier.go:120, 191-197, 205-211`).
     Clients using `WaitNewEvent` after a long-poll always re-validate by
     re-reading mutable state, which closes the gap; a pure notification
     consumer must add its own at-least-once handling.
3. **Are events named as immutable past-tense facts?**
   - Yes. Server-side factory names are `WorkflowExecutionStarted`,
     `WorkflowTaskScheduled`, `WorkflowTaskStarted`, `WorkflowTaskCompleted`,
     `WorkflowTaskTimedOut`, `ActivityTaskScheduled`, `ActivityTaskStarted`,
     `ActivityTaskCompleted/Failed/TimedOut/Canceled`,
     `WorkflowExecutionCompleted/Failed/TimedOut/Terminated/Canceled/ContinuedAsNew/OptionsUpdated/UpdateAccepted/UpdateCompleted/Signaled/CancelRequested/Paused/Unpaused/TimeSkippingTransitioned`,
     `TimerStarted/Fired/Canceled`,
     `StartChildWorkflowExecutionInitiated`,
     `ChildWorkflowExecutionStarted/Completed/Failed/TimedOut/Canceled/Terminated`,
     `RequestCancelExternalWorkflowExecutionInitiated`,
     `ExternalWorkflowExecutionSignaled/Failed`,
     `SignalExternalWorkflowExecutionInitiated`,
     `MarkerRecorded`, `UpsertWorkflowSearchAttributes`,
     `WorkflowPropertiesModified`,
     `NexusOperationScheduled/Started/Completed/Failed/Canceled/CancelRequested/TimedOut`
     (e.g., `service/history/historybuilder/event_factory.go:39, 112, 136, 166, 190, 213, 234, 264, 285, 304, 323, 340, 357, 373, 389, 413, 438, 454, 466, 481, 512, 524, 538, 555, 574, 587, 604, 619, 649, 677, 700, 727, 741, 761, 787, 807, 829, 851, 893, 917, 942, 967, 992, 1015, 1038, 1058, 1106`).
   - All in the past tense; no event encodes a "request to do" as a fact — only
     "X started", "X completed", "X failed".
4. **Can the consumer tell why an event happened?**
   - Partially. Each event carries its scheduled/initiated/started chain
     (`scheduledEventId`, `initiatedEventId`, `startedEventId`,
     `lastCompletedEventId`) wired by `EventStore.wireEventIDs`
     (`service/history/historybuilder/event_store.go:381-454`), plus request-id
     causation on a subset (`WorkflowExecutionOptionsUpdated.AttachedRequestId`,
     `WorkflowExecutionSignaled.RequestId` — same file lines 441-452).
   - The full `Request` that triggered a command is embedded in the originating
     event (e.g., `WorkflowExecutionStartedEventAttributes` includes `Input`,
     `Identity`, `RetryPolicy`, `Header`, etc. —
     `service/history/historybuilder/event_factory.go:50-89`).
   - Nexus/completion callbacks carry `Links` and `RequestId`
     (`api/historyservice/v1/request_response.pb.go:8764, 9022`; events see
     `Links` propagated into history events such as the start event via
     `AddWorkflowExecutionStartedEvent` —
     `service/history/historybuilder/history_builder.go:179-187`).
   - **What is missing**: there is no generic envelope-level
     `causationId`/`parentEventId` field. Causal traversal must be done by
     walking the typed pointer fields.
5. **How does a schema evolve without breaking old clients?**
   - Protobuf `reserved` field numbers protect against accidental reuse:
     `TransientWorkflowTaskInfo` reserves 1 and 2
     (`proto/internal/temporal/server/api/history/v1/message.proto:11-12`).
   - `StrippedHistoryEvent` skips to field number 4 for `version`
     (`api/history/v1/message.pb.go:338-339`), leaving gaps for future use.
   - Public API surface uses `// Deprecated:` doc comments and helper wrappers
     for old proto field paths
     (`api/history/v1/message.pb.go:60-69, 106-117, 159-176, 212-228, 265-274, 316-326, 369-386, 420-432`).
   - The frontend has an explicit backwards-compat layer
     (`FixFollowEvents` — `service/history/api/get_history_util.go:554-638`)
     that synthesizes a synthetic ContinuedAsNew event for older clients that
     do not understand `NewExecutionRunId`. This pattern (rewrite at read time)
     is the operational evidence of a migration strategy.
   - Server-internal tokens are versioned through opaque blobs
     (`api/token/v1/message.pb.go:30-42`); the history service can detect old
     tokens via `consts.ErrInvalidNextPageToken`
     (`service/history/api/getworkflowexecutionhistory/api.go:232-239`).

## Architectural Decisions

- **Two-tier identity (event vs. transition)**: The server uses `EventId`
  for intra-shard ordering and `VersionedTransition{NFV, TC}` for cross-cluster
  ordering. The `VersionedTransition` lives next to workflow execution info
  and on every persisting task object
  (`api/persistence/v1/executions.pb.go:223, 275, 315-336, 1178-1181, 1444, 1837, 2206-2218, 2844, 3314, 3406, 3547`,
  `api/persistence/v1/hsm.pb.go:457-465`). Rationale: failover safety requires
  a transition count that survives namespace-version resets, while intra-shard
  ordering just needs a monotonic integer.
- **Append-only history with branching for replication**: Replication produces
  a new `VersionHistoryItem` (one item per cluster generation, identified by
  `Version = NamespaceFailoverVersion`), and forks produce branches reconciled
  through `FindLCAVersionHistoryItem`
  (`common/persistence/versionhistory/version_history.go:111-137`).
- **Buffered-then-flushed ordering**: Activity/timer/child/nexus completions
  are buffered (`HistoryBuilder.bufferEvent`,
  `service/history/historybuilder/event_store.go:297-360`) and re-ordered by
  `reorderBuffer` (`service/history/historybuilder/event_store.go:461-491`)
  before flush — an explicit acknowledgement that the storage layer can
  reorder and the builder is responsible for restoring event-id order.
- **Pull+push hybrid for resume**: Continuation-token pagination is the
  pull-mode; the in-memory notifier is the push-mode. They share the same
  cursor semantics (continuation token). The push-mode is strictly an
  optimization; correctness does not depend on it because every long-poll
  re-checks mutable state
  (`service/history/api/get_workflow_util.go:128-260`).
- **Race detection across pages**: After the last page returns, the API
  re-queries mutable state, detects events committed during pagination
  (`freshNextEventID > continuationToken.NextEventId`), and back-fills them
  (`service/history/api/getworkflowexecutionhistory/api.go:456-480`).
- **In-place backwards compatibility at read**: `FixFollowEvents` rewrites
  failure/time-out close events into a ContinuedAsNew variant for older SDKs
  that do not understand `NewExecutionRunId`
  (`service/history/api/get_history_util.go:554-638`).

## Notable Patterns

- **Oneof event attributes**: The protobuf `HistoryEvent.Attributes` oneof
  guarantees exactly-one attribute per event; consumers must switch on
  `EventType` (e.g., `service/history/historybuilder/event_factory.go:100-103`).
- **Defensive key validation**: `EventKey{EventID, Version}` plus
  `validateKey` (`service/history/events/cache.go:97-108`) makes the cache
  treat malformed keys as a logged warning, not a crash.
- **Explicit causal-pointers wiring in builder**: `wireEventIDs` walks the
  buffered events and patches `StartedEventId` into completion/failure events
  (`service/history/historybuilder/event_store.go:381-454`). This is the
  closest analogue to a generic causation field.
- **Drop-on-overflow fanout**: The notifier uses buffered channels of size 1
  per subscriber and silently drops on full
  (`service/history/events/notifier.go:131, 191-197`). Subscriber-id is a UUID
  to avoid collisions
  (`service/history/events/notifier.go:132`).
- **Versioned transitions on tasks**: Each persisted task object carries
  `LastUpdateVersionedTransition` so staleness can be detected even after
  shard reload (`api/persistence/v1/executions.pb.go:2844, 3177, 3314, 3384,
  3406, 3526, 3547, 3610`).
- **Local-vs-remote history split during replication**: `splitBatchesToLocalAndRemote`
  partitions history batches by `SplitVersionHistoryByLastLocalGeneratedItem`
  so that import and replication paths each get the correct segment
  (`service/history/replication/eventhandler/history_events_handler.go:114-148`).

## Tradeoffs

- **Opaque continuation token**: Couples callers to a server-versioned payload
  that includes both event-cursor (`FirstEventId`/`NextEventId`) and storage
  cursor (`PersistenceToken`). Pro: smaller wire format, hides storage. Con: any
  schema bump in the token requires careful migration handling.
- **Single mutable-state source of truth per shard**: Ordering and gap-fill
  are easy because the shard serializes writes through mutable state. Tradeoff
  is shard contention — a single workflow's hot shard becomes the bottleneck.
- **Best-effort notifications**: The notifier is fast and non-blocking. Tradeoff
  is silent drops under load; clients must use `WaitNewEvent` long-poll (which
  itself validates via a re-read) rather than rely on push alone.
- **Causation via typed pointer fields**: Avoids a generic "parent event id"
  field that would not fit every shape. Tradeoff is that consumers must
  peruse attribute-specific fields to construct a causal graph; there is no
  generic `CausationEventId`/`CorrelationId` envelope.
- **Buffered-then-flushed events**: Lets the builder assign stable `EventId`s
  only at flush time. Tradeoff is the comment at
  `service/history/historybuilder/event_store.go:456-460` — this complexity is
  preserved for backward compatibility (`HasActivityFinishEvent`,
  `hasActivityFinishEvent`) and is explicitly marked TODO.

## Failure Modes / Edge Cases

- **Dropped notification**: A subscriber's channel is full (size 1) when an
  event is published; the event is dropped without an error
  (`service/history/events/notifier.go:191-197`). Mitigated by the
  long-poll always re-reading mutable state.
- **Cassandra-style reorder on buffer flush**: Detected and reordered before
  commit (`service/history/historybuilder/event_store.go:183-185, 461-491`)
  with a dedicated metric
  (`metrics.OutOfOrderBufferedEventsCounter`).
- **Stale continuation token**: Detected via
  `serviceerrors.NewCurrentBranchChanged`
  (`service/history/api/get_workflow_util.go:115-118, 155-179`); the client is
  expected to start a fresh read on the new branch token.
- **Invalid next-page token**: Rejected with `consts.ErrInvalidNextPageToken`
  (`service/history/api/getworkflowexecutionhistory/api.go:232-239`).
- **Data loss / incomplete history**: The API explicitly returns
  `serviceerror.DataLoss` on incomplete results and trims the history node
  to recover (`service/history/api/getworkflowexecutionhistory/api.go:293-308`,
  `service/history/api/get_history_util.go:159-183, 216-229`).
- **Branch ownership violation**: When the requested branch token is not owned
  by the execution, the API either rejects the request or serves it under
  relaxed enforcement
  (`service/history/history_engine_test.go:7129-7220`).
- **Out-of-date notifier event**: Long-poll re-checks mutable state and skips
  events whose `VersionHistoryItem` is older than the caller's
  (`service/history/api/get_workflow_util.go:218-221`).
- **Shard-owned channel deadlock**: `LongPollExpirationInterval` and
  `common.DefaultLongPollBuffer` are used to ensure long-poll never blocks
  past `ctx.Deadline()`
  (`service/history/api/get_workflow_util.go:190-192, 248-259`).

## Future Considerations

- A generic `causation_event_id`/`correlation_id` envelope field would make
  causal tracing first-class instead of the current pattern of
  type-specific `scheduledEventId`/`initiatedEventId`/`startedEventId`
  pointers.
- The notification fanout channel size of 1 is a known sharp edge. A bounded
  ring buffer with explicit overflow signaling (rather than silent drop)
  would let callers detect and re-read.
- `HistoryContinuation` mixes storage-internal cursor with event cursor. If
  storage pagination changed (e.g., new shard layout), every old client
  token would need to be migrated. A versioned cursor with a clear separation
  between "client-safe" and "server-internal" halves would reduce coupling.
- `FixFollowEvents` is an old-SDK compatibility hack. Adding a feature
  flag in newer SDKs (`FeatureFollowsNextRunID`) allows it to be removed once
  all clients cross the threshold
  (`service/history/api/get_history_util.go:573-583`).
- `versionhistory` is still partly duplicated with the new
  `transitionhistory` model. Some helpers (e.g., `ContainsVersionHistoryItem`)
  are kept on the older API while new code uses `VersionedTransition` for
  staleness. A future consolidation would clarify which is the "single source
  of truth" for ordering.

## Questions / Gaps

- The exact wire format of `HistoryContinuation` is private to the server; the
  public client only sees an opaque `NextPageToken`. Whether clients should
  ever inspect or persist these tokens across server upgrades is not
  documented in the analyzed source — search bound: `api/token/v1/message.pb.go`
  does not document a version field.
- Searched for "causation" / "correlation" / "causationId" /
  "correlationId" within the server source. No evidence found in
  `service/history/`, `api/history/`, or `api/persistence/`. Causation is
  conveyed only through typed pointer attributes on individual event types.
- `HistoryEvent` envelope (in `go.temporal.io/api`) carries no
  `correlation_id` or `causation_id` field; the server factory does not add
  one (`service/history/historybuilder/event_factory.go:1085-1096`).
- The notifier has no concept of "last delivered event id" to a subscriber;
  the cursor lives on the caller side. This means observers must keep their
  own watermark.