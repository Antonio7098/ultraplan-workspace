# Source Analysis: temporal-sdk-go

## Lifecycle Transition Ownership and Terminal Arbitration

### Source Info

| Field | Value |
|-------|-------|
| Name | temporal-sdk-go |
| Path | `studies/aren-go-runtime-study/sources/temporal-sdk-go` |
| Language / Stack | Go (Temporal Go SDK, gRPC against a Temporal server) |
| Analyzed | 2026-08-26 |

All citations below are workspace-relative to `studies/aren-go-runtime-study/sources/temporal-sdk-go`.

## Summary

Temporal's Go SDK gives one component authority over terminal transitions, but it does it in two layers that deserve to be described separately.

The first layer is a per-operation command state machine (`internal/internal_command_state_machine.go`). Every schedulable thing a workflow does (activity, timer, child workflow, signal, nexus operation, cancel requests) gets a small machine with an enumerated state set (12 states, lines 183-196). History events drive transitions in, workflow code drives cancellations in, and any impossible transition panics immediately (`panicIllegalState`, lines 500-507, error code TMPRL1100). There is no way for two completions to race: `handleCompletionEvent` accepts only specific predecessor states, and a second completion hits `failStateTransition`.

The second layer is the workflow-level completion slot. User workflow functions do not publish anything themselves. The root coroutine writes `(result, error)` into a pointer slot (`internal/internal_workflow.go:555-557`), and `executeDispatcher` is the only place that turns that slot into a terminal call: `env.Complete(rp.workflowResult, rp.error)` on normal return (`internal/internal_workflow.go:733`) or `env.Complete(nil, panicErr)` when a coroutine panicked (`internal/internal_workflow.go:696`). `Complete` delegates to one handler wired at construction (`internal/internal_event_handlers.go:150`, `internal/internal_task_handlers.go:701`), and that handler, `workflowExecutionContextImpl.completeWorkflow`, just records `isWorkflowCompleted`, `result`, and `err` under the context mutex (`internal/internal_task_handlers.go:653-657`).

Publication to the outside world is deferred even further. Nothing goes on the wire until the whole workflow task finishes. `CompleteWorkflowTask` collects pending commands, protocol messages, and the outcome into a single `RespondWorkflowTaskCompletedRequest` (`internal/internal_task_handlers.go:1427-1472`). At most one close command exists in that request, chosen by checking the single `err` field in a fixed priority order: canceled, continue-as-new, failed, completed (`internal/internal_task_handlers.go:1884-1958`). Because there is exactly one `err` field and one arbitration function, conflicting terminal outcomes cannot coexist. A cancel that arrives after a failure was recorded simply loses; the failure publishes.

Timing, state, outcome, and notification become visible together because they travel in the same RPC. The server applies the task completion and its commands atomically and appends the resulting history events. On replay, the SDK treats terminal history events (`WORKFLOW_EXECUTION_COMPLETED`, `_FAILED`, `_TIMED_OUT`) as no-ops (`internal/internal_event_handlers.go:1322-1327`), meaning the local runtime never "observes" its own termination. Termination is proposed locally and confirmed by history. That inversion is the most interesting idea in the codebase for Aren's purposes.

The cost is size. Owning lifecycle this precisely requires ~1700 lines of state machine, a cooperative coroutine dispatcher with deadlock detection, a sticky-cache context lock, replay matching, and five separate panic-recovery sites. The mechanism is excellent; the surface area is framework-sized.

## Rating

**9/10** for the property under study.

Rationale:

- Single owner, enforced structurally: one completion handler field, one result slot, one close-command arbitration function. Two goroutines cannot publish conflicting terminals because workflow code never touches the wire and real concurrency is serialized per run by `workflowExecutionContextImpl.mutex` (`internal/internal_task_handlers.go:106`, locked at 817/865/868).
- Request versus confirmation is kept distinct three times over: `WORKFLOW_EXECUTION_CANCEL_REQUESTED` only cancels a context (`internal/internal_event_handlers.go:1692-1694`, `internal/internal_workflow.go:570-574`); the close command waits for user code to return `CanceledError`; and the command machine keeps `commandStateCancellationCommandSent` separate from `commandStateCompletedAfterCancellationCommandSent` (`internal/internal_command_state_machine.go:191-192`) so a completion arriving mid-cancellation still resolves correctly.
- Panics, ordinary errors, and deadlock are normalized into one error channel at each goroutine boundary, then re-specialized by type at publication (`*workflowPanicError`, `*CanceledError`, `*ContinueAsNewError`; `internal/error.go:187-190`).
- Deducting one point for admitted warts: `removeCancelOfResolvedCommand` exists only because cancel commands can be queued for already-resolved timers ("This really should not exist", `internal/internal_command_state_machine.go:1093-1097`), production allows repeated `Complete` calls where the test environment guards against them (`internal/internal_workflow_testsuite.go:1143-1147` versus `internal/internal_task_handlers.go:653-657`), and the `BlockWorkflow` panic policy deliberately wedges a workflow into retry-forever (`internal/internal_task_handlers.go:1338-1342`).

## Evidence Collected

| Area | Evidence | File:Line |
|------|----------|-----------|
| Command state enumeration | 12 command states incl. cancel-in-flight and completed-after-cancel-sent variants | `internal/internal_command_state_machine.go:183-196` |
| Illegal transition = panic | `failStateTransition` raises TMPRL1100; used by every handler default case | `internal/internal_command_state_machine.go:504-507` |
| Completion acceptance rules | Base machine accepts completion only from initiated/cancel-after-initiated states | `internal/internal_command_state_machine.go:557-566` |
| Cancel idempotence | Cancel on completed or already-canceled machine is a documented no-op | `internal/internal_command_state_machine.go:516-531` |
| Cancel-vs-resolve cleanup | Pending cancel command removed when timer/activity resolves first | `internal/internal_command_state_machine.go:1136-1145`, `1623-1631`, `1097-1107` |
| Child-workflow cancel accepted tracking | New state `commandStateCancellationCommandAccepted` records server acceptance separately | `internal/internal_command_state_machine.go:786-817` |
| Single completion slot | `completeHandler` is one field; `Complete` forwards to it | `internal/internal_event_handlers.go:150`, `413-415` |
| Result written only by root coroutine | Root func stores `workflowResult` via context pointer; dispatcher reads it | `internal/internal_workflow.go:555-557`, `700-704` |
| Dispatcher is sole terminal caller | Panic path and normal-return path both end in `env.Complete` | `internal/internal_workflow.go:692-734` |
| Re-entrancy guard | Second `ExecuteUntilAllBlocked` while running panics | `internal/internal_workflow.go:1287-1292` |
| Per-run serialization | `workflowExecutionContextImpl.mutex` held across task processing; queries included | `internal/internal_task_handlers.go:104-125`, `615-644`, `794-878` |
| Outcome arbitration order | Canceled > ContinueAsNew > Failed > Completed, exactly one close command appended | `internal/internal_task_handlers.go:1884-1958` |
| Coherent publication | Commands + messages + close command in one `RespondWorkflowTaskCompletedRequest` | `internal/internal_task_handlers.go:1427-1472`, `1991-2010` |
| Stop-on-terminal in event loop | Loop breaks as soon as `isWorkflowCompleted` set | `internal/internal_task_handlers.go:1183-1186`, `1223-1234`, `1241-1243` |
| Terminal events are no-ops on replay | COMPLETED/FAILED/TIMED_OUT/CANCELED events perform no operation | `internal/internal_event_handlers.go:1322-1327`, `1391-1392` |
| Cancel request vs confirmation | `CANCEL_REQUESTED` fires context cancel; close command needs `CanceledError` return | `internal/internal_event_handlers.go:1688-1694`, `internal/internal_workflow.go:566-574` |
| Callback single-delivery guard | `scheduledActivity.handle` panics on second delivery via `handled` flag | `internal/internal_event_handlers.go:287-343`, checked at `1595-1597` |
| Coroutine panic capture | Recover in `coroutineState.run` stores `workflowPanicError` with stack | `internal/internal_workflow.go:1234-1245` |
| Deadlock becomes panic error | Ticker expiry synthesizes TMPRL1101 panic error instead of crashing the process | `internal/internal_workflow.go:1156-1186`, detector in `internal/workflow_deadlock.go:12-146` |
| Blocking during panic unwinding forbidden | `yield` detects `runtime.gopanic` on stack and re-panics | `internal/internal_workflow.go:1057-1087` |
| Event-application panic recovery | `ProcessEvent`/`ProcessMessage` recover into `Complete(nil, panicErr)` | `internal/internal_event_handlers.go:1302-1309`, `1518-1525` |
| FailWorkflow panic policy | Policy layer overrides recorded error with application error before publishing | `internal/internal_task_handlers.go:1305-1349` |
| Activity panic boundary | Recover in `activityTaskHandlerImpl.Execute` produces failure response | `internal/internal_task_handlers.go:2439-2458` |
| Local activity panic boundary | Recover in LA goroutine converts panic to result error; result flows via channel | `internal/internal_task_pollers.go:1016-1041`, delivery at `956-967` |
| Activity result arbitration | One `err` decides completed/canceled/failed; cancel reporting gated by `cancelAllowed` | `internal/internal_task_pollers.go:1605-1685` |
| Unexpected-cancel handling | `CanceledError` without allowance wrapped as failure application error | `internal/internal_task_pollers.go:1668-1672` |
| Update protocol mini-machine | accept/reject/complete with `requireState` guards | `internal/internal_update.go:128-215` |
| Replay determinism check | Generated commands matched against history events; mismatch fails task | `internal/internal_task_handlers.go:1285-1295`, `1581-1632` |
| Test env first-complete-wins | Test environment ignores later `Complete` calls, unlike production | `internal/internal_workflow_testsuite.go:1143-1147` |
| SM tests: completion after cancel | `Test_ActivityStateMachine_CompletedAfterCancel` | `internal/internal_command_state_machine_test.go:219-253` |
| SM tests: pending cancel dropped | `Test_ActivityStateMachine_CompletionRemovesPendingCancelWithCustomActivityID` | `internal/internal_command_state_machine_test.go:255-276` |
| SM tests: invalid transition panics | `Test_ActivityStateMachine_PanicInvalidStateTransition` | `internal/internal_command_state_machine_test.go:315-344` |
| SM tests: unusual child ordering | `Test_ChildWorkflow_UnusualCancelationOrdering` asserts no panic on late cancel-accepted | `internal/internal_command_state_machine_test.go:503-535` |
| Integration: panic policies | `TestPanicFailWorkflow`, `TestPanicActivityWorkflow` | `test/integration_test.go:600-618` |
| Integration: cancel/command races | `TestTimerCancellationConcurrentWithOtherCommandDoesNotCausePanic`, `TestCancelChildAndExecuteActivityRace` | `test/integration_test.go:2408-2424`, `2957-2960` |
| Integration: ordering regression doc | Workflow comment explains command-order bug the race test pins | `test/workflow_test.go:2360-2391` |
| Integration: deferred-yield panic | `TestPanicWithDeferredYield` pins the freeze-during-unwinding fix | `test/integration_test.go:9713-9767` |

## Answers to Dimension Questions

**Can two goroutines publish conflicting terminal outcomes, and what prevents it?**

No, for four stacked reasons. First, workflow code runs on simulated coroutines scheduled by one dispatcher loop, and re-entering that loop while it executes panics (`internal/internal_workflow.go:1287-1292`). Second, all five `Complete` producers (normal return, coroutine panic, event-processing panic, message-processing panic, panic-policy override) write through one handler into one `(result, err)` pair guarded by the per-run mutex (`internal/internal_event_handlers.go:150`, `internal/internal_task_handlers.go:106`, `653-657`). Third, publication happens once per task in `workflowTaskHandlerImpl.completeWorkflow`, which inspects that single `err` and emits at most one close command (`internal/internal_task_handlers.go:1884-1958`). Fourth, the server rejects a second close anyway; the SDK even assumes close success and evicts the cached context on completion (`internal/internal_task_handlers.go:630-632`). What "prevents" conflict is less a lock than a funnel: everything converges on one slot, and the last writer within a task wins deterministically because writers run on the same thread in a fixed order.

**Is a cancellation request represented separately from confirmed termination?**

Yes, consistently at three levels. At the workflow level, `WORKFLOW_EXECUTION_CANCEL_REQUESTED` triggers only `d.cancel()` on a `workflow.Context` (`internal/internal_workflow.go:570-574`); the terminal `COMMAND_TYPE_CANCEL_WORKFLOW_EXECUTION` is emitted only if user code actually returns a `CanceledError` (`internal/internal_task_handlers.go:1890-1896`). At the command level, request and confirmation get separate states (`commandStateCancellationCommandSent`, `commandStateCancellationCommandAccepted`, `commandStateCompletedAfterCancellationCommandSent`; `internal/internal_command_state_machine.go:190-195`). For child workflows there is even a dedicated handler distinguishing "cancel delivered" from "child finished" (`handleExternalWorkflowExecutionCancelRequested`, `internal/internal_command_state_machine.go:809-817`). Inbound confirmation events are no-ops because the SDK treats its own proposal as non-authoritative until history confirms it.

**Can observers see a terminal state before its outcome or terminal event is available?**

Not through the published record. The close command, its payload/failure details, and all pending commands leave in one request (`internal/internal_task_handlers.go:1427-1472`), so external observers see state and outcome atomically in history. Queries answered from cached state hold the same mutex as the completing task (`internal/internal_task_handlers.go:817`), so a query either sees pre-terminal state or is destroyed with the context on error (`Unlock(err)`, `internal/internal_task_handlers.go:626-644`). One caveat: `isWorkflowCompleted` flips true the moment user code returns, before the RPC lands, so in-process code sharing the context could observe "completed" milliseconds before the server does. The SDK contains this by stopping event processing immediately (`internal/internal_task_handlers.go:1232-1234`) and by clearing current-task state at publication.

**Which parts can Aren adopt without importing a framework-sized lifecycle model?**

Adoptable without the rest of the framework:

- The funnel pattern: one completion handler field, one `(state, outcome)` slot, one function that maps outcome to terminal action. This is maybe 100 lines (`internal/internal_task_handlers.go:653-657`, `1884-1958`).
- Enumerated states with fail-fast illegal transitions, including the cancel-in-flight states. The core of `commandStateMachineBase` is ~120 lines (`internal/internal_command_state_machine.go:483-601`) and needs no Temporal types to imitate.
- The `handled` flag guard on result callbacks (`internal/internal_event_handlers.go:287-343`), a cheap guarantee against double delivery.
- Panic normalization at each goroutine boundary into one typed error carrying a stack trace (`internal/error.go:185-190`, recovery sites listed above), plus the refusal to let deferred code block during unwinding (`internal/internal_workflow.go:1057-1087`).
- Proposal/confirmation asymmetry: treat a local terminal as a proposal that only externally confirmed events make authoritative.

What would drag the framework along: the full command taxonomy (16 command types), version markers and their event-ID arithmetic (`getNextID` skipping version-marker slots, `internal/internal_command_state_machine.go:1044-1067`), the sticky cache with stale-detection heuristics (`ResetIfStale`, `internal/internal_task_handlers.go:1511-1526`), protocol-message outboxes, and SDK feature flags. Aren's contract, as stated in the dimension, does not need replay-equivalent determinism checking; that requirement alone justifies half of this code.

Abstractions smaller than Aren's contract also exist here and are worth stealing individually: `futureImpl.Set` panics on second set (`internal/internal_workflow.go:423-434`), and the naive state machines show how little machinery a fire-and-forget operation needs (`internal/internal_command_state_machine.go:85-113`).

## Architectural Decisions

- **Terminal publication is separated from terminal determination.** Five different failures can determine the terminal outcome, but only `executeDispatcher` and the panic-policy layer ever call `Complete`, and only `workflowTaskHandlerImpl.completeWorkflow` translates outcome into a wire message. Each stage has one job (`internal/internal_workflow.go:692-734`, `internal/internal_task_handlers.go:1305-1349`, `1846-1958`).
- **One error value encodes the outcome class.** Rather than a status enum, the SDK uses typed errors (`CanceledError`, `ContinueAsNewError`, others) inspected with `errors.As` in priority order (`internal/internal_task_handlers.go:1886-1953`). Simple, but it means error types are load-bearing API, not diagnostics.
- **Illegal transitions crash the workflow task, not the process.** State machine violations panic, the panic is recovered at the event boundary, converted to a `workflowPanicError`, and published as a workflow failure or WFT failure depending on configured policy (`internal/internal_command_state_machine.go:500-507`, `internal/internal_event_handlers.go:1302-1309`, `internal/internal_task_handlers.go:1332-1345`).
- **Deadlock detection reuses the panic path.** A coroutine that does not yield within the deadline gets a synthesized TMPRL1101 panic error with its stack captured from another goroutine (`internal/internal_workflow.go:1170-1185`, stack identification trick at `1118-1148`). Runtime failure and user panic end in the same channel.
- **Authoritative state lives in history, not memory.** Replay re-executes user code against recorded events; the deterministic-match check compares regenerated commands to recorded ones and fails the task on drift (`internal/internal_task_handlers.go:1278-1295`, `1581-1632`). Memory is an optimization (sticky cache), never the source of truth.
- **Cancellation is cooperative and observable.** Context cancellation propagates to timers and activities through the interceptor layer; `waitForCancelRequest` options let workflows await confirmed cancellation instead of assuming it (`internal/internal_event_handlers.go:851-863`, activity option plumbing in `internal/internal_event_handlers.go:56-87`).

## Notable Patterns

- **Handled-flag idempotence guards.** Every scheduled callback wrapper (`scheduledTimer`, `scheduledActivity`, `scheduledChildWorkflow`, `scheduledCancellation`, `scheduledSignal`) panics on second delivery (`internal/internal_event_handlers.go:287-343`). Cheap, loud, and catches producer bugs early.
- **Cancel-of-resolved-command removal.** When a timer or activity resolves, any queued cancel command for it is silently dropped (`handleActivityTaskClosed` at `internal/internal_command_state_machine.go:1136-1145`). Necessary because cancellation callbacks can queue commands after the resolving event was already applied, an ordering wart the code comments own up to (`1093-1097`).
- **First-write-wins in tests, last-write-wins in production.** The test environment's `Complete` ignores repeat calls (`internal/internal_workflow_testsuite.go:1143-1147`); production overwrites. The divergence is harmless in practice because production writers are ordered, but it is a place where the test double and the real thing disagree.
- **Priority-ordered error inspection.** `errors.As` chains checked in a fixed sequence give a total order over mutually exclusive outcomes without an explicit state variable (`internal/internal_task_handlers.go:1890-1953`).
- **Event-ID prediction.** The SDK predicts server-assigned event IDs from the WFT started event ID plus a counter adjusted for unsent commands (`setCurrentWorkflowTaskStartedEventID`, `internal/internal_command_state_machine.go:1023-1042`). This is what makes replay matching possible, and it is fragile enough that version markers need special-case skipping (`1044-1067`).

## Tradeoffs

- **Precision versus size.** Getting terminal coherence right took a per-command-type state machine, a custom coroutine runtime, and replay infrastructure. An Aren-sized runtime cannot afford all of it; the transferable part is the discipline, not the code volume.
- **Fail-fast panics versus resilience.** Panicking on illegal transitions surfaces nondeterminism immediately and turns silent corruption into a failed task. The cost is that benign edge cases (late cancels, duplicate events) need dedicated states or removal hacks to avoid false alarms.
- **Typed-error arbitration versus explicit status.** Using the returned error as the outcome avoids a parallel status enum, but couples error identity to control flow. Anyone wrapping these errors in middleware risks changing terminal semantics.
- **Deferred publication versus latency.** Buffering all commands until task completion buys atomicity for free, but a workflow that schedules work and then blocks for minutes holds its outcome hostage to the next task boundary. The heartbeat/force-complete path exists precisely to mitigate this (`internal/internal_task_handlers.go:962-1051`).
- **Simulated goroutines versus real ones.** Cooperative scheduling eliminates data races on workflow state entirely, at the price of implementing channels, selectors, timers, and mutexes from scratch (`internal/internal_workflow.go:802-1031`, `1397-1575`).

## Failure Modes / Edge Cases

- **BlockWorkflow wedges permanently.** With that policy, a panicking workflow fails its WFT once, then times out and retries forever until someone fixes the code (`internal/internal_task_handlers.go:1338-1342`, comment says so plainly). Deliberate, but operationally nasty.
- **Late cancel-accepted after child completion.** Covered by `Test_ChildWorkflow_UnusualCancelationOrdering`: cancel delivered after the child already resolved must not panic (`internal/internal_command_state_machine_test.go:503-535`). The `CompletedAfterCancellationCommandSent` intermediate state exists specifically to survive this window.
- **Duplicate child workflow ID.** Rejected eagerly in-process rather than relying on server rejection, with a comment admitting the local check is a workaround for command IDs being equal to workflow IDs (`internal/internal_command_state_machine.go:1399-1414`).
- **Blocking inside a deferred function during panic unwinding.** Freezes the panic and lets post-panic code emit commands; now detected and re-panicked (`internal/internal_workflow.go:1077-1083`), with a regression integration test (`test/integration_test.go:9713-9767`).
- **Activity completes after deadline.** The result is discarded and the task returns a timeout; logged as "Activity complete after timeout" (`internal/internal_task_handlers.go:2485-2495`). Late success is indistinguishable from failure downstream, which is correct but lossy for debugging.
- **Query tasks skip local activities entirely** (`hasPendingLocalActivityWork`, `internal/internal_task_handlers.go:1474-1481`), so query answers never observe side effects of pending LAs. Sound, but worth knowing when reasoning about observed state.

## Future Considerations

For Aren Phase 1, the concrete takeaways in priority order:

1. Build the funnel first: a single completion slot plus a single arbitration function that maps outcome to exactly one terminal action, before any provider or persistence work touches the core. This matches the roadmap constraint that lifecycle comes first.
2. Keep cancellation as a request that mutates only a context/flag, and make confirmed termination require either user-code agreement or an explicit arbiter decision. Copy the cancel-in-flight states; they cost little and eliminate a whole class of races.
3. Normalize panics at every goroutine boundary into one typed error with a captured stack, and forbid blocking operations during unwinding.
4. Add the `handled` guard to every result callback from day one.
5. Decide consciously whether Aren wants replay-as-source-of-truth. If not, skip command/event matching and event-ID prediction entirely; that removes most of the bulk without losing terminal arbitration.
6. Guard the completion slot itself (first-write-wins) rather than relying on writer ordering, closing the small divergence this SDK left between its test and production environments.

## Questions / Gaps

- The dimension cites `Aren/docs/phase-1-prd/02-lifecycle-requirements.md` and `01-product-definition-and-scope.md`. These sit outside the isolated source directory and were not read; the Aren-side statements in this report derive from the requirements summarized in the injected dimension text. If the PRD imposes timing guarantees (for example, maximum delay between terminal decision and publication), this study did not verify them against the source. Nothing in the SDK bounds that delay except the workflow task timeout.
- Whether repeated `Complete` calls in production can ever produce a *different* final outcome than the first writer intended: the writers are ordered on one thread, so the panic-policy override intentionally wins, but a formal argument for that ordering lives mostly in code comments rather than tests. No direct test asserts the final outcome when both an event-processing panic and a policy override occur in one task.
- The sticky-cache eviction path (`onEviction`, `internal/internal_task_handlers.go:659-672`) runs on a separate goroutine and takes the same mutex; I did not trace whether a forced eviction racing a completing task can drop a computed close command (the RPC would then be rebuilt from replayed history, so the effect should be benign duplication, but this is inference from structure, not observed behavior).

---

Generated by `01.01-lifecycle-transition-ownership-and-terminal-arbitration` against `temporal-sdk-go`.
