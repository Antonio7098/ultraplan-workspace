# Source Analysis: temporal-sdk-go

## 01.14 Durable Workflows, Retry, Idempotency, and Replay

### Source Info

| Field | Value |
|-------|-------|
| Name | temporal-sdk-go |
| Path | `studies/aren-go-runtime-study/sources/temporal-sdk-go` |
| Language / Stack | Go (SDK for Temporal server, protobuf history, worker, workflow coroutines) |
| Analyzed | 2026-08-30 |

## Summary

temporal-sdk-go is the reference durable-workflow SDK: workflow decisions are deterministic coroutines replayed from server-persisted history, while all side-effects are isolated to activities, timers, child workflows, Nexus operations, and marker-recorded side effects. Replay validation (`internal/internal_task_handlers.go:1581` `matchReplayWithHistory`) enforces command/history equivalence, `GetVersion`/`SideEffect`/`MutableSideEffect` provide versioned deterministic evolution, and a full command state-machine (`internal/internal_command_state_machine.go:151`) drives retries, timeouts, cancellation, and child ownership. At-most-once workflow decisions are guaranteed by history; at-least-once execution applies to activities/local-activities/heartbeats. Idempotency is enforced at effect boundaries via `WorkflowIdReusePolicy`, `RequestId` deduplication, deterministic activity/child IDs, and `RetryPolicy.MaximumAttempts=1` to disable retries. The machinery is production complete but heavy: history iterators, sticky cache, deadlock detector, update/nexus protocols, and extensive replay fixtures, implying high carry cost for Aren.

## Rating

**8/10** — Most complete open-source Go reference for durable workflows. Separation of deterministic vs effectful code is explicit and enforced (deadlock detector, `IsReplaying`, `IsReadOnly`, replay checker, replayer tests). History persistence, replay, versioning, timers, child workflows, Nexus, cancellation, and local-activity retries are implemented with state machines and validated by `test/replaytests`. Deductions: client-side idempotency is delegated to server `RequestId`/`WorkflowIdReusePolicy`; SDK does not expose storage-level deduplication or idempotency keys beyond those; poison-work handling relies on server retry + `WorkflowPanicPolicy` rather than SDK DLQ.

## Evidence Collected

Every entry MUST include a file path with line numbers. Format: `path/to/file.ts:NN`.

| Area | Evidence | File:Line |
|------|----------|-----------|
| Deterministic workflow definition | `workflow` type interface doc: workflow code must be deterministic, use `workflow.Channel/Selector/Go`, no map range, use `GetTime(ctx)` (`internal/internal_workflow.go:99`) | `studies/aren-go-runtime-study/sources/temporal-sdk-go/internal/internal_workflow.go:99-101` |
| Deterministic dispatcher | Coroutine dispatcher `ExecuteUntilAllBlocked` deterministic order, deadlock detection (`internal/internal_workflow.go:82`, `1281`) | `studies/aren-go-runtime-study/sources/temporal-sdk-go/internal/internal_workflow.go:82-96,1281-1339` |
| Deterministic wrappers | `WorkflowRandomStream`, `DeterministicKeys`, `NewChannel/NewSelector/Await` must be used instead of native (`workflow/deterministic_wrappers.go:9`, `workflow/workflow.go:906`) | `studies/aren-go-runtime-study/sources/temporal-sdk-go/workflow/deterministic_wrappers.go:9-30`, `studies/aren-go-runtime-study/sources/temporal-sdk-go/workflow/workflow.go:906-918` |
| SideEffect / MutableSideEffect | `SideEffect` records result into history marker, returns cached on replay; `MutableSideEffect` deduplicates via ID+call counter (`internal/internal_event_handlers.go:1054`, `1155`) | `studies/aren-go-runtime-study/sources/temporal-sdk-go/internal/internal_event_handlers.go:1054-1091`, `studies/aren-go-runtime-study/sources/temporal-sdk-go/internal/internal_event_handlers.go:1155-1246` |
| GetVersion versioning | `GetVersion` records `Version` marker + search attribute, returns `DefaultVersion` on replay if first call (`internal/internal_event_handlers.go:958`, `workflow/workflow.go:595`) | `studies/aren-go-runtime-study/sources/temporal-sdk-go/internal/internal_event_handlers.go:958-1002`, `studies/aren-go-runtime-study/sources/temporal-sdk-go/workflow/workflow.go:529-597` |
| IsReplaying / IsReadOnly guards | `IsReplaying` warning not to branch on it; `IsReadOnly` for query/validator/sideEffect; `assertNotInReadOnlyState` panics on blocking in read-only (`workflow/workflow.go:723`, `internal/internal_workflow.go:760`) | `studies/aren-go-runtime-study/sources/temporal-sdk-go/workflow/workflow.go:723-748`, `studies/aren-go-runtime-study/sources/temporal-sdk-go/internal/internal_workflow.go:760-785` |
| History recording & replay | `history` struct, `nextTask`, `prepareTask` reorders events, handles stale cache via `resetHistory` (`internal/internal_task_handlers.go:186`, `242`, `371`) | `studies/aren-go-runtime-study/sources/temporal-sdk-go/internal/internal_task_handlers.go:186-200,242-559` |
| Replay determinism check | `matchReplayWithHistory` compares replay commands vs history events, skips flags (`internal/internal_task_handlers.go:1581`) + `isCommandMatchEvent` per command type (`1642`) | `studies/aren-go-runtime-study/sources/temporal-sdk-go/internal/internal_task_handlers.go:1581-1632,1642-1817` |
| Replay fixtures / replayer | `worker.NewWorkflowReplayer().ReplayWorkflowHistoryFromJSONFile` with 30+ JSON histories covering timers, LA, SideEffect, updates, child, Nexus (`test/replaytests/replay_test.go:120`) | `studies/aren-go-runtime-study/sources/temporal-sdk-go/test/replaytests/replay_test.go:120-132,150-583` |
| Command state machines | `commandsHelper` + 15 command types + state machine transitions `CREATED->SENT->INITIATED->COMPLETED` (`internal/internal_command_state_machine.go:151`, `183`) | `studies/aren-go-runtime-study/sources/temporal-sdk-go/internal/internal_command_state_machine.go:151-170,183-229` |
| Activity retry policy | `RetryPolicy` fields `InitialInterval/BackoffCoefficient/MaximumInterval/MaximumAttempts/NonRetryableErrorTypes` (`internal/internal_workflow_testsuite.go` test `84`), server default unlimited unless `MaximumAttempts=1` (`internal/activity.go:151`) | `studies/aren-go-runtime-study/sources/temporal-sdk-go/internal/activity.go:151-161`, `studies/aren-go-runtime-study/sources/temporal-sdk-go/internal/workflow_test.go:84-118` |
| Activity timeouts | `ScheduleToClose/StartToClose/ScheduleToStart/HeartbeatTimeout` set on `ScheduleActivityTaskCommandAttributes` (`internal/internal_event_handlers.go:803`) | `studies/aren-go-runtime-study/sources/temporal-sdk-go/internal/internal_event_handlers.go:803-811` |
| Activity heartbeat & redelivery | `temporalInvoker.Heartbeat` throttled 80% of `HeartbeatTimeout`, batched, handles `CancelRequested/Paused/Reset`, activity cancellation via `context.CancelCause` (`internal/internal_task_handlers.go:2186`, `2360`) | `studies/aren-go-runtime-study/sources/temporal-sdk-go/internal/internal_task_handlers.go:2186-2311,2360-2400` |
| Local activity retry with backoff | `getRetryBackoff`/`getRetryBackoffWithNowTime` exponential, respects `MaximumAttempts`, `NonRetryableErrorTypes`, `expireTime`; short backoffs via `time.AfterFunc` local retry, large backoffs via server timer marker (`internal/internal_task_handlers.go:1381`, `1351`) | `studies/aren-go-runtime-study/sources/temporal-sdk-go/internal/internal_task_handlers.go:1351-1426` |
| Timer / cancellation | `NewTimer`/`RequestCancelTimer` creates `timerCommandStateMachine` with cancel state machine (`internal/internal_command_state_machine.go:652`, `931`) processed in `workflowExecutionEventHandlerImpl.ProcessEvent` `TIMER_STARTED/FIRED/CANCELED` (`internal/internal_event_handlers.go:1379`) | `studies/aren-go-runtime-study/sources/temporal-sdk-go/internal/internal_command_state_machine.go:652-676,1602-1641`, `studies/aren-go-runtime-study/sources/temporal-sdk-go/internal/internal_event_handlers.go:1379-1388` |
| Child workflow ownership | `ExecuteChildWorkflow` with `ParentClosePolicy`, `WorkflowIdReusePolicy`, deterministic `WorkflowID = currentRunID + GenerateSequenceID` if empty (`internal/internal_event_handlers.go:571`), `childWorkflowCommandStateMachine` handles `Initiated/Started/Completed/Canceled` (`internal/internal_command_state_machine.go:730`) | `studies/aren-go-runtime-study/sources/temporal-sdk-go/internal/internal_event_handlers.go:565-648`, `studies/aren-go-runtime-study/sources/temporal-sdk-go/internal/internal_command_state_machine.go:730-816` |
| Nexus operations | `ExecuteNexusOperation` `ScheduleNexusOperation` + `RequestCancelNexusOperation` state machines, `scheduledEventIDToNexusSeq` mapping, cancellation types (`internal/internal_event_handlers.go:650`, `internal/internal_command_state_machine.go:123`) | `studies/aren-go-runtime-study/sources/temporal-sdk-go/internal/internal_event_handlers.go:650-686`, `studies/aren-go-runtime-study/sources/temporal-sdk-go/internal/internal_command_state_machine.go:123-144,1196-1281` |
| Checkpoint / sticky cache & worker loss | `WorkflowCache` LRU, `previousStartedEventID`/`lastHandledEventID` staleness check, `GetOrCreateWorkflowContext` `StickyCacheHit/Miss`, `resetHistory` on eviction, `clearState` on completion/error (`internal/internal_task_handlers.go:603`, `794`) | `studies/aren-go-runtime-study/sources/temporal-sdk-go/internal/internal_task_handlers.go:603-694,794-878` |
| Idempotency / deduplication | `WorkflowIdReusePolicy` on `ChildWorkflowOptions` (`internal/workflow.go:489`), `RequestId` deduplication on `StartWorkflow` (`internal/client.go:467`), duplicate child check `childWorkflowExistsWithId` (`internal/internal_command_state_machine.go:1399`), duplicate `SignalChannel` buffer check (`internal/workflow.go:932`) | `studies/aren-go-runtime-study/sources/temporal-sdk-go/internal/workflow.go:489-493`, `studies/aren-go-runtime-study/sources/temporal-sdk-go/internal/client.go:467`, `studies/aren-go-runtime-study/sources/temporal-sdk-go/internal/internal_command_state_machine.go:1399-1414`, `studies/aren-go-runtime-study/sources/temporal-sdk-go/internal/workflow.go:928-938` |
| Poison work / panic handling | `WorkflowPanicPolicy` `FailWorkflow` vs `BlockWorkflow` (infinite WorkflowTaskTimeout retry) (`internal/internal_task_handlers.go:1332`), `newWorkflowPanicError` vs `PanicError` (`internal/error.go:836`), activity panic captured and returned as `PanicError` (`internal/internal_task_handlers.go:2440`) | `studies/aren-go-runtime-study/sources/temporal-sdk-go/internal/internal_task_handlers.go:1332-1346`, `studies/aren-go-runtime-study/sources/temporal-sdk-go/internal/error.go:836-862`, `studies/aren-go-runtime-study/sources/temporal-sdk-go/internal/internal_task_handlers.go:2440-2458` |
| Updates / idempotent re-delivery | `bufferedUpdateRequests` queued until handler registered, `DrainUnhandledUpdates`, dedup duplicate update IDs in testsuite (`internal/internal_event_handlers.go:174`, `internal/workflow_testsuite.go:3461`) | `studies/aren-go-runtime-study/sources/temporal-sdk-go/internal/internal_event_handlers.go:174-180,1097-1123`, `studies/aren-go-runtime-study/sources/temporal-sdk-go/internal/internal_workflow_testsuite.go:3461-3675` |
| Deadlock detection | `deadlockDetector` per-workflow `time.Ticker` default 1s, `Pause/Resume` for data converters doing remote calls (`internal/workflow_deadlock.go:12`, `29`) | `studies/aren-go-runtime-study/sources/temporal-sdk-go/internal/workflow_deadlock.go:12-83` |
| ContinueAsNew checkpoint truncation | `ContinueAsNewError` carries new workflow type/input/header/timeouts, `completeWorkflow` appends `CONTINUE_AS_NEW` command (`internal/error.go:192`, `internal/internal_task_handlers.go:1897`) | `studies/aren-go-runtime-study/sources/temporal-sdk-go/internal/error.go:192-237`, `studies/aren-go-runtime-study/sources/temporal-sdk-go/internal/internal_task_handlers.go:1897-1929` |

## Answers to Dimension Questions

**Which code must remain deterministic and which code may perform external effects?**

*Must remain deterministic:* All code executed inside `workflow.Context` — the workflow function itself, coroutines created via `workflow.Go`/`GoNamed`, `Await` conditions, signal/query/update handlers, and update validators. Implementation enforces this via `syncWorkflowDefinition` coroutine dispatcher (`studies/aren-go-runtime-study/sources/temporal-sdk-go/internal/internal_workflow.go:99-104` doc + `internal/internal_workflow.go:676-688` `newDispatcher`), which requires use of `workflow.Channel/Selector/WaitGroup/Mutex/Semaphore` (`studies/aren-go-runtime-study/sources/temporal-sdk-go/workflow/deterministic_wrappers.go:9-52`), `workflow.Now` (`workflow/deterministic_wrappers.go:174`), `workflow.NewTimer/Sleep` (`workflow/deterministic_wrappers.go:180-208`), `DeterministicKeys` iteration, and bans `time.Now/rand/map iteration/go chan/select`. `assertNotInReadOnlyState` panics on blocking calls inside queries/validators/side-effects (`internal/internal_workflow.go:760-768`), and `deadlockDetector` (1s default, `internal/workflow_deadlock.go:72`) fails workflow tasks that do blocking I/O. `SideEffect`/`MutableSideEffect` are the *only* sanctioned escape hatches for non-deterministic snippets; they execute once, record a `MarkerRecorded` event, and return the cached payload on replay (`internal/internal_event_handlers.go:1054-1091`).

*May perform external effects:* Activity functions (`activity` package, `internal/activity.go:64`), local activities (`internal/internal_event_handlers.go:865-888`), Nexus operation handlers (outside workflow), and client calls. These are invoked via `ExecuteActivity`/`ExecuteLocalActivity`/`ExecuteChildWorkflow` which emit commands that become history events; the effect itself runs outside the replay domain and is retried independently per `RetryPolicy`. Effects must be idempotent or guarded by dedup because they may execute more than once (see next question).

**What can execute more than once after worker or daemon failure?**

*Workflow decisions* do not — they are re-executed deterministically from history after stick-cache eviction or worker crash (`internal/internal_task_handlers.go:828-862` staleness + `resetHistory`). The `replayCommands == historyEvents` check (`internal/internal_task_handlers.go:1287-1291`) ensures decisions are at-most-once logically.

*Effects* are at-least-once:

- **Activities:** Server retries per `RetryPolicy` independent of worker; `maximumAttempts=0` default is unlimited (`internal/activity.go:156-158`). After worker loss, task redelivered to another poller; activity may have already written side effects. Activity heartbeat timeout determines detection (`internal/internal_task_handlers.go:2408-2412` throttle calc + `internal/activity.go:364` deadline).
- **Local activities:** Retried locally via `getRetryBackoff` (`internal/internal_task_handlers.go:1381-1424`). Short backoffs (< `WorkflowTaskTimeout`) via local `time.AfterFunc` tunneled through `laRetryCh`; long backoffs create a server timer and replay via `LocalActivity` marker (`internal/internal_task_handlers.go:1356-1376`). Can re-execute on same worker after WFT failure or after cache eviction on different worker (marker replay + `pendingLaTasks` rebuild).
- **Timers/child workflows/Nexus:** Commands are idempotently replayed, but underlying external actions (child workflow start, Nexus call) are exactly-once at server level only if `WorkflowIdReusePolicyRejectDuplicate`/`RequestId` holds; otherwise duplicate initiation fails via `childWorkflowExistsWithId` (`internal/internal_command_state_machine.go:1409`). Cancellation requests may be retried and result in `CanceledError`.
- **SideEffect functions:** Doc warns no at-most-once guarantee (`workflow/doc.go:377-381`): under workflow task failure before marker committed, function can execute twice; hence must be short and side-effect-free beyond return value.
- **Updates/signals:** `bufferedUpdateRequests` and testsuite dedup of duplicate update IDs (`internal/internal_workflow_testsuite.go:3675`) show server may redeliver updates; handler validation must be idempotent.

**How are workflow definition changes made compatible with existing histories?**

1. **Version markers (`GetVersion`).** First call generates a `Version` marker + `TemporalChangeVersion` search attribute (`internal/internal_event_handlers.go:983-995`). Marker name/details (`change-id`, `version`, `version-search-attribute-updated`) stored via `recordVersionMarker` (`internal/internal_command_state_machine.go:1283-1314`) and `handleVersionMarker` population of `versionMarkerLookup` (`1317-1331`). On replay, `GetVersion` returns recorded version (`internal/internal_event_handlers.go:959-961`), `validateVersion` (`945-956`) panics if removed support violates `minSupported/maxSupported`. Pattern: `v := GetVersion(ctx,"fooChange",DefaultVersion,1); if v==DefaultVersion { old } else { new }` (`workflow/workflow.go:529-597` doc). SDK flags (`SDKFlagLimitChangeVersionSASize` etc., `internal/internal_event_handlers.go:985`) ensure old histories without flags still replay (`test/replaytests/replay_test.go:327-338` version loop tests).

2. **`WorkflowReplayer` + fixtures.** `worker.NewWorkflowReplayer().ReplayWorkflowHistoryFromJSONFile` validates old JSON histories against current code (`test/replaytests/replay_test.go:120-132`). 30+ histories cover LA, SideEffect, mutable side effect, search attributes, continue-as-new, Nexus cancels. CI is expected to run `go run ./internal/cmd/build check` and integration replayer before merge (`AGENTS.md:9-22`).

3. **Build ID / worker deployment versioning.** `VersioningBehavior` (`Pinned` vs `AutoUpgrade`) (`internal/workflow.go:45-68`) on workflow type and `WorkerDeploymentVersion` on `WorkerDeploymentOptions` route tasks; `preferredVersionProvider` resolves which `GetVersion` branch to pick during replay (`internal/internal_event_handlers.go:1004-1033`).

4. **MutableSideEffect for config.** Instead of new code path, `MutableSideEffect(id,f,equals)` records a new marker only when `!equals(old,new)` (`internal/internal_event_handlers.go:1187-1206`), allowing dynamic config without breaking replay while still deterministic (`workflow/workflow.go:501-523` doc).

If change is not wrapped in `GetVersion`/`MutableSideEffect`, replay fails with `nondeterministic workflow: missing/extra replay command` (`internal/internal_task_handlers.go:1616-1625`) or `lookup failed for scheduledEventID` (`test/replaytests/replay_test.go:152-162`), and `WorkflowPanicPolicy` determines whether workflow blocks forever or fails (`internal/internal_task_handlers.go:1332-1346`).

**What concrete Aren use would justify adopting any of this machinery?**

- **Current Aren is single-node CLI daemon with synchronous jobs** (roadmap Phase 16 deferred). None of the heavy machinery is justified today: history persistence is just log files, retries are simple `RetryPolicy` on jobs, there is no need for sticky cache, worker loss, or replay.

- **Justified increment (minimal subset) when Aren gains long-running, resumable agent loops spanning daemon restarts:** 
  - *Deterministic decision log + replay of LLM tool calls as activities.* Treat planner/executor decision steps as workflow decisions, LLM API calls and tool executions as activities with retry policy (`RetryPolicy{MaximumAttempts:1, NonRetryableTypes:["PayloadValidationError"]}` for non-idempotent tools). This gives crash resume without adopting timers/child workflows/ Nexus.
  - *`SideEffect` for UUID/random + `GetVersion` for prompt evolution.* Small, self-contained; cost is marker table.
  - *Timers via `NewTimer/Sleep` for polling/scheduling.* Only if Aren needs durable sleep across restarts; otherwise simple `time.Sleep` suffices.
  - *WorkflowIdReusePolicy + RequestId for at-least-once launch.* If Aren exposes an HTTP launch API, use idempotency key to reject duplicate job creation (mirrors `internal/client.go:467`).

- **Defer:** child workflows (fan-out), Nexus/remote operations, sessions, update protocol, worker deployment versioning, deadlock detection, payload codecs with remote calls. These require server, history store, versioned workers, and extensive operational runbooks. Aren should only adopt after demonstrating pain from manual checkpointing (e.g., multi-hour research workflows that must survive `aren daemon` upgrades without losing progress). Even then, prefer a lightweight embedded history (sqlite + event replay) over full Temporal server unless cross-host scale is needed.

## Architectural Decisions

- **Event-sourced history as source of truth.** Every `ScheduleActivityTask`, `StartTimer`, `StartChildWorkflow`, `RecordMarker`, `UpsertSearchAttributes` is a command appended deterministically; server persists as `HistoryEvent` stream (`HistoryEvent` protobuf). Worker never persists workflow state locally beyond sticky cache; on `Unlock` error or cache size 0, `clearState` forces rebuild from history (`internal/internal_task_handlers.go:636-644`). Tradeoff: strong durability, but requires strict determinism.

- **Coroutine dispatcher instead of goroutines.** `dispatcherImpl` with `coroutineState` channels (`aboutToBlock/unblock`) guarantees sequential, deterministic interleaving (`internal/internal_workflow.go:152-179`, `1044-1088`). Enables `WorkflowPanicPolicy BlockWorkflow` infinite retry on nondeterminism while still allowing `Workflow.Go` concurrency within replay.

- **State-machine per command type.** Activity/timer/child/Nexus each have explicit state (`Created/CommandSent/Initiated/Started/Completed/...`) (`internal/internal_command_state_machine.go:183-196`). `commandsHelper.orderedCommands` list preserves command order (`1068-1091`); `moveCommandToBack` handles child cancellation reordering (`1109-1116`). This isolates replay mismatch to specific command (`failStateTransition` panics with `[TMPRL1100]`).

- **Separate local-activity path.** Local activities bypass server task queue, executed by `localActivityTunnel` poller (`internal/internal_worker.go:393-438`). Retry backoff split local vs server timer avoids holding WFT open for long retries (`internal/internal_task_handlers.go:1356-1376`). Complexity: separate `localActivityCounterID`, `pendingLaTasks/unstartedLaTasks`.

- **Failure conversion layer.** `DataConverter` + `FailureConverter` with workflow/activity serialization context (`converter.With...SerializationContext`) ensures payloads and errors are decoded with correct namespace/workflow/activity keys (`internal/internal_task_handlers.go:234-241`, `2379-2387`). Failure types are preserved across wire as `failurepb.Failure` (`internal/error.go:346-365`).

- **Pluggable retry with `IsRetryable`.** `NonRetryableErrorTypes` exact string match, `ApplicationError.NonRetryable`, timeout type limited to `START_TO_CLOSE/HEARTBEAT` (`internal/error.go:1088-1119`). No jitter config in SDK; server computes backoff.

## Notable Patterns

- **Marker pattern for otherwise non-deterministic work.** `SideEffect` (payload in `MarkerRecorded` details `side-effect-id/data`), `Version` (change-id/version), `LocalActivity` (data/result/backoff), `MutableSideEffect` (id + callCounter + data) all share `markerCommandStateMachine` path (`internal/internal_command_state_machine.go:232-244`, `1414-1396`). Unified `RecordMarker` command with  `completeOnSendStateMachine` that completes once sent.

- **Determinism escape hatch pairing.** Workflow code calls `workflow.IsReplaying`/`IsReadOnly` only for logging/metrics (`workflow/workflow.go:723-748` warning), while SDK provides `GetLogger`/`GetMetricsHandler` that are replay-aware and auto-suppressed (`internal/internal_event_handlers.go:269-282`).

- **Sticky-cache + incremental history.** `previousStartedEventID` watermark + `lastHandledEventID` in `workflowExecutionContextImpl` (`internal/internal_task_handlers.go:107-122`) enables partial history fetch (`isFullHistory` check `880-885`) and `resetHistory` fallback (`717-725`). LRU eviction via `WorkerCache` (not shown) emits `StickyCacheTotalForcedEviction`.

- **Test-time simulation via `testsuite`.** `WorkflowTestSuite`/`TestWorkflowEnvironment` mocks activities/child workflows/nexus, provides `OnActivity/OnWorkflow/OnGetVersion/OnSideEffect` and `RegisterNexusAsyncOperationCompletion` (`internal/workflow_testsuite.go:435-794`), and simulates time advances, enabling unit tests without server.

- **Update protocol as outbox + SDK flags.** `Send` appends to `outbox` with `eventPredicate` (`internal/internal_event_handlers.go:365-377`), `addProtocolMessage` gated by `SDKFlagProtocolMessageCommand`, replay matches via `isCommandMatchEvent` protocol message branch (`internal/internal_task_handlers.go:1643-1651`). Ensures forward compatibility with old histories.

## Tradeoffs

- **Strong replay guarantee vs developer velocity.** Every workflow change must be version-gated; forgotten `GetVersion` breaks all open workflows with `nondeterministic` panic and triggers `BlockWorkflow` infinite task failure until code fixed. High operational burden for small teams.

- **Activity at-least-once vs exactly-once illusion.** SDK docs state activity may execute more than once; users must implement idempotency externally (keyed by `ActivityID` which defaults to `GenerateSequence` `internal/internal_event_handlers.go:798-801` but can be overridden). No built-in dedup store; server retries are transparent.

- **Local activity latency vs durability.** Local activities finish in seconds but block WFT; poller runs in same process, so worker shutdown delays completion (`internal/internal_task_handlers.go:984-987` `laTunnel.stopCh`). Not heartbeatable, not scalable, not routable.

- **History size vs debuggability.** Search attributes/memo/metering metadata appended to every `RespondWorkflowTaskCompleted` (`internal/internal_task_handlers.go:1979-2008`); large payloads trigger `payloadSizeError` checks (`internal/internal_task_handlers.go:2557`). `ContinueAsNewSuggestedReasons` (`HistorySizeTooLarge/TooManyHistoryEvents/TooManyUpdates`) signal need for `ContinueAsNew` truncation (`workflow/workflow.go:66-83`), which fragments workflow execution and complicates queries.

- **Payload visitor/ codec extensibility vs deadlock risk.** Custom `DataConverter` that does HTTP can trigger deadlock detector; SDK provides `DataConverterWithoutDeadlockDetection` that pauses detector (`internal/workflow_deadlock.go:43-51`) at cost of disabling deadlock safety during encode.

- **Build-ID vs deployment versioning evolution.** Deprecated `WorkerBuildID/UseBuildIDForVersioning` plus new `WorkerDeploymentVersion`/`DeploymentOptions` coexist (`internal/internal_worker.go:154-168`), increasing configuration surface.

## Failure Modes / Edge Cases

- **Nondeterministic workflow code.** Adding/removing/reordering activity/timer/child calls without `GetVersion` → `matchReplayWithHistory` returns `historyMismatchError` with `[TMPRL1100]` (`internal/internal_task_handlers.go:1616-1625`), bubbled to `workflowPanicError`, then per `WorkflowPanicPolicy` either failed workflow (app error) or blocked task spinning indefinitely (`internal/internal_task_handlers.go:1332-1346`). Stack trace included.

- **Duplicate activity/child signal send after resolve.** `removeCancelOfResolvedCommand` (`internal/internal_command_state_machine.go:1093-1107`) suppresses late cancel-timer/activity cancels that race with `TimerFired/ActivityCompleted` in same WFT; without this, replay would generate extra clear command and mismatch.

- **Worker loss / sticky cache eviction.** If worker crashes while WFT open, server reassigns to new worker; `GetOrCreateWorkflowContext` detects `history.Events[0].EventId != previousStartedEventID+1` → evicts cache and calls `resetHistory` to fetch full history (`internal/internal_task_handlers.go:828-857`). Local activities in `pendingLaTasks` are reconstructed from markers; in-flight attempts lost, retried via marker `Backoff` time.

- **Poison work.** Activity that panics returns `PanicError` → `ActivityError` → retried per policy until `MaximumAttempts` exhausted (`internal/error.go:840`, `internal/internal_task_handlers.go:2440`). Workflow that panics repeatedly hits `WorkflowPanicPolicy.BlockWorkflow` → WFT timed out infinitely, occupant must deploy fix or terminate via CLI; no automatic DLQ.

- **Heartbeat loss / duplicate detection.** Activity heartbeat records via `RecordActivityTaskHeartbeat` (`internal/internal_task_handlers.go:2635-2655`); if worker dies, server heartbeat timeout fires (`TimeoutType=HEARTBEAT`), activity retried on new worker. Throttled batching may lose last details if worker shuts down gracefully — flushed in `Close(flush=true)` only if not completed (`2422-2427`).

- **Duplicate update / signal redelivery.** Server may redeliver same `Update` protocol message after WFT failure; SDK dedup via `protocol.Registry` and testsuite `dedup duplicate update IDs` (`internal/internal_workflow_testsuite.go:3675` caches result). Workflow validator must be side-effect free; stray map iteration without `DeterministicKeys` can cause heisenbugs.

- **Payload size / memo limits.** `setErrorLimits` from namespace capabilities (`internal/internal_worker.go:1408-1421`) enforces `BlobSizeLimitError`; exceeding returns sync workflow error before command sent (`internal/internal_event_handlers.go:837` check, `visitProtoPayloads` size error `2548-2576`).

- **Cancellation race.** `childWorkflowCommandStateMachine.cancel` moves command to back (`internal/internal_command_state_machine.go:770-784`) to keep commands in occurrence order; without, replay would place cancel before activity scheduled in same WFT and mismatch.

- **Version marker search attribute size overflow.** `changeVersionSearchAttrSizeLimit = 2048` (`internal/internal_event_handlers.go:32`); exceeding skips upsert and logs warning but still records marker (`985-995`), leaving visibility incomplete.

## Future Considerations

- **Minimum viable for Aren (Phase 16):** Implement (a) event log (append-only JSONL / sqlite with `HistoryEvent` equivalent), (b) deterministic replay of tool-use decisions vs idempotent activity wrappers, (c) `SideEffect` for `uuid/rand` and `GetVersion` for prompt changes, (d) `RetryPolicy` with `MaximumAttempts=1` for non-idempotent tools, heartbeat for long LLM calls. Defer timers, child workflows, sessions, Nexus, updates, deployment versioning.

- **If full durability needed:** Prefer Temporal Cloud / self-hosted server rather than re-implementing history, matching, sticky queues, deadline propagation, and replayer — cost is ~10k LOC (`internal_task_handlers + command_state_machine + event_handlers ≈ 5k+ LOC`) plus operational runbooks for `WorkflowPanicPolicy`, payload limits, history growth.

- **Storage-level idempotency.** For Aren's external effects (file writes, git commits), adopt server-side dedup pattern: make `ActivityID` deterministic (`ActivityID = "job:<jobID>:step:<seq>"` via `WithActivityOptions` `internal/workflow/activity_options.go:86` example) and store idempotency key in activity implementation, checking before effect.

- **Observability.** Before durable workflows, add `WorkflowInfo`–like metadata (`WorkflowID/RunID/Memo/SearchAttributes`) to Aren job records so future migration can replay; design payload visitor pattern for audit logging.

- **Testing.** Adopt `testsuite` style (`internal/workflow_testsuite.go:887` `ExecuteWorkflow` + `OnActivity`) plus `ReplayWorkflowHistoryFromJSONFile` snapshot tests for any Aren workflow, without requiring live server.

- **Eager activities.** SDK's eager dispatch (`defaultMaxEagerActivityReservationsPerWorkflowTask=3`, `internal/internal_worker.go:69`) is optimization for colocated workflow+activity workers; irrelevant for Aren until it colocates.

## Questions / Gaps

- **No evidence of exactly-once storage dedup in SDK.** Grep for `idempotency/dedup` found only `RequestId` at client level (`internal/client.go:467`, `internal/internal_workflow_client.go:1351`) and workflow reuse policies (`internal/workflow.go:489`). No SDK-level dedup table for activity side effects; need server source to confirm transactional task completion.

- **No direct worker-loss integration test in SDK repo.** Searched `test/workflow_test.go` and `testsuite`; no test named `worker-loss`/`duplicate-delivery`. Replay tests cover history mismatch, not redelivery of task queue messages. Assumed coverage via server tests; not verified in this source.

- **ContinueAsNew as checkpoint vs unbounded history:** SDK exposes reason `HistorySizeTooLarge` (`workflow/workflow.go:73`) but no auto-truncation; caller must explicitly return `ContinueAsNewError`. Unclear how Aren would know when to trigger.

- **Operator intervention docs missing from code.** Panic `FailWorkflow` vs `BlockWorkflow` is configured via `workerExecutionParameters.WorkflowPanicPolicy` (`internal/internal_worker.go:186`), but no CLI/tool code for reset/terminate in SDK; requires `temporal` CLI or API outside SDK.

- **Nexus cancellation semantics partial.** `NexusOperationCancellationType` (`workflow/workflow.go:101-120`) `WaitCompleted vs WaitRequested vs TryCancel` distinguishes but tests only cover replayer (`test/replaytests/replay_test.go:565`), not live cancellation race.

---

Generated by `01.14 Durable Workflows, Retry, Idempotency, and Replay` against `temporal-sdk-go`.
