# Source Analysis: docker-agent

## 01.11 Execution Composition and Child Ownership

### Source Info

| Field | Value |
|-------|-------|
| Name | docker-agent |
| Path | `studies/aren-go-runtime-study/sources/docker-agent` |
| Language / Stack | Go 1.27 (`go.mod:3`), charm/bubbletea TUI, MCP, expr, provider SDKs |
| Analyzed | 2026-08-30 |

## Summary

docker-agent implements child ownership via `session.Session` parent linkage and two distinct execution paths rather than a generic `errgroup`/`pool` composition library. `SubSessionConfig` + `newSubSession` (`pkg/runtime/agent_delegation.go:115-275`) constructs identity-coupled child sessions (system prompt, `ParentID`, `DelegationLineage`, `AgentName` pinning, permissions inheritance). `runForwarding` (`pkg/runtime/agent_delegation.go:340-412`) provides sequential, blocking, event-forwarding execution for `transfer_task` and fork-skills; `runCollecting` (`pkg/runtime/agent_delegation.go:422-545`) provides detached, non-interactive execution for `run_background_agent` backed by `pkg/tools/builtin/agent.Handler` (`pkg/tools/builtin/agent/agent.go:236-440`) with explicit concurrency caps (`maxConcurrentTasks=20`, `maxTotalTasks=100`). Lifecycle guards (`validateDelegation` depth 10 + cycle detection), session-scoped event filtering (`pkg/runtime/event.go:18-22`, `pkg/runtime/persistence_observer.go:74-79`, `pkg/runtime/event_sink.go:14-36`), and context detachment (`context.WithoutCancel` at `pkg/tools/builtin/agent/agent.go:337`) define ownership. There is no generic fail-fast parallel, collect-all ordered, or bounded-pool helper for arbitrary child sets; parallelism is limited to independent background tasks queried individually by `taskID`. Panic propagation is absent; parent terminal semantics allow background children to outlive the parent stream.

## Rating

**Rating: 5 / 10**

Rationale: Core parent-child identity, sequential delegation, and bounded detached concurrency are well-engineered and heavily tested (delegation lineage isolation, pinned vs shared `currentAgent` handling #3886, explicit depth/cycle guards). However the composition model lacks the structured Go primitives the dimension expects: no `errgroup.WithContext` fail-fast group, no `SetLimit` semaphore pool for children, no deterministic ordered result collection, no sibling-cancellation policy, and no panic-to-error propagation. Background children can outlive parent terminal via `WithoutCancel`, and `Handler.tasks` iteration order is non-deterministic, failing ordered-results criteria. The implementation preserves concrete executor semantics well for its two paths but does not generalize to Phase-12 composition patterns.

## Evidence Collected

Every entry MUST include a file path with line numbers. Format: `path/to/file.ts:NN`.

| Area | Evidence | File:Line |
|------|----------|-----------|
| Parent identity | `ParentID` field, `WithParentID`, `IsSubSession()` (sub-session predicate) | `pkg/session/session.go:422-425`, `pkg/session/session.go:1527-1531`, `pkg/session/session.go:1621-1623` |
| Delegation lineage | `DelegationLineage []string` not persisted, snapshot accessor, lineage guard | `pkg/session/session.go:426-433`, `pkg/session/session.go:1369-1373`, `pkg/runtime/agent_delegation.go:60-85` |
| Child config | `SubSessionConfig` 15 fields covering system/user message, `PinAgent`, `DelegationLineage`, tool filters, `DisableStructuredOutput` | `pkg/runtime/agent_delegation.go:115-176` |
| Child factory | `newSubSession` snapshots `AttachedFiles`, builds lineage, merges `ExcludedTools`, clones permissions | `pkg/runtime/agent_delegation.go:214-275` |
| Sequential execution | `runForwarding` blocks parent, optionally `swapCurrentAgent`, drives `RunStream`, drains `ErrorEvent`, emits `SubSessionCompleted` | `pkg/runtime/agent_delegation.go:340-412` |
| Detached execution | `runCollecting` pins `AgentName`, uses `context.WithoutCancel` descendant, drops events except `TokenUsage`, persists via `persistBackgroundSubSession` | `pkg/runtime/agent_delegation.go:422-545`, `pkg/tools/builtin/agent/agent.go:337` |
| Bounded parallelism (cap) | `maxConcurrentTasks=20`, `maxTotalTasks=100`, `maxOutputBytes=10MB`, enforced before spawn with `pruneCompleted()` | `pkg/tools/builtin/agent/agent.go:33-40`, `pkg/tools/builtin/agent/agent.go:319-329` |
| Bounded parallelism (errgroup elsewhere) | `errgroup.WithContext` + `SetLimit` used only for RAG embedding/file-indexing, not child agents | `pkg/rag/embed/embed.go:120-122`, `pkg/rag/strategy/vector_store.go:306-307` |
| Sibling cancellation | No errgroup cancellation; each background task gets own `cancel` from `context.WithCancel(WithoutCancel(ctx))`; `StopAll` CAS each `running→stopped` and `wg.Wait()` | `pkg/tools/builtin/agent/agent.go:337`, `pkg/tools/builtin/agent/agent.go:507-515` |
| Current-agent isolation | `swapCurrentAgent` emits `AgentSwitching`/`AgentInfo` and hooks; pinned parents degrade switch to pin (`if parent.AgentName != "" → req.PinAgent=true`) | `pkg/runtime/agent_delegation.go:300-319`, `pkg/runtime/agent_delegation.go:356-366` |
| Context detachment | `taskCtx, cancel := context.WithCancel(context.WithoutCancel(ctx))` — parent cancellation does not kill background task | `pkg/tools/builtin/agent/agent.go:337` |
| Result collection (single) | `tools.ResultSuccess(s.GetLastAssistantMessageContent())` — only last assistant message string | `pkg/runtime/agent_delegation.go:411` |
| Result collection (background) | `RunResult{Result string, ErrMsg string}` per task; retrieved via `HandleView`/`HandleList` by `taskID`, unordered map iteration | `pkg/tools/builtin/agent/agent.go:73-77`, `pkg/tools/builtin/agent/agent.go:443-464`, `pkg/tools/builtin/agent/agent.go:466-482` |
| Panic handling | No child panic recovery; only channel `recover()` for `send on closed channel` in bridge and sink | `pkg/runtime/elicitation.go:108-109`, `pkg/runtime/event_sink.go:36`, `pkg/runtime/event_sink.go:63` |
| Parent terminal | `RunStream` creates buffered channel, `finalizeEventChannel` emits `StreamStopped` best-effort non-blocking then closes channel; background `StopAll` not tied to parent stream end, only `LocalRuntime.Close()` | `pkg/runtime/loop.go:230-259`, `pkg/runtime/loop.go:196-228` |
| Cleanup completion | `AddLiveSubSession` marks `liveAttached=true` to avoid double-counted cost; `SessionUsage` excludes live-attached cost from parent aggregation | `pkg/session/session.go:1101-1107`, `pkg/runtime/event.go:413-428` |
| Event attributability | `SessionScoped` interface, `TokenUsageEvent.GetSessionID`, `ErrorEvent.GetSessionID`, `ElicitationRequestEvent.SessionID`; observer filters `scoped.GetSessionID()!=sess.ID` | `pkg/runtime/event.go:18-22`, `pkg/runtime/event.go:256-258`, `pkg/runtime/event.go:401`, `pkg/runtime/event.go:678-679` |
| Persistence observer filtering | `OnRunStart` skips sub-sessions; `OnEvent` drops `IsSubSession()` and `SessionScoped` mismatches before `AddSubSession` | `pkg/runtime/persistence_observer.go:60-80` |
| Background event bridge | `OnBackgroundEvent` + `backgroundEventMu` for `TokenUsage` from detached tasks; `emitBackgroundEvent` in `runCollecting` | `pkg/runtime/runtime.go:179-184`, `pkg/runtime/runtime.go:356-361`, `pkg/runtime/agent_delegation.go:465` |
| SubagentStop attribution | `executeSubagentStopHooks` always runs against `*parent` agent's executor, deferred to guarantee on error | `pkg/runtime/agent_delegation.go:372-380`, `pkg/runtime/hooks.go:727-733` |
| Depth/cycle guard | `maxDelegationDepth=10` fixed, `validateDelegation` clones lineage slice to avoid sharing backing array (sibling isolation test) | `pkg/runtime/agent_delegation.go:60`, `pkg/runtime/agent_delegation.go:69-85`, `pkg/runtime/agent_delegation_test.go:744-754` |
| Live session registry | `registerLiveSession`/`finishLiveSession` with `liveSessionsMu`, `activeRootStreams` counter; sub-sessions share root budget not fresh allowance | `pkg/runtime/loop.go:250`, `pkg/runtime/loop.go:240-244` |
| Skill sub-session | Fork-skill `DisableStructuredOutput:true`, `ExcludedTools:[run_skill]`, `PinAgent` when parent pinned, no lineage edge | `pkg/runtime/skill_runner.go:118-138` |
| Branch/Fork title isolation | `setParentIDs` recursively re-parents branched sessions; `Clone` preserves `ParentID`/`AgentName` where `BranchSession` resets them | `pkg/session/branch.go:381-391`, `pkg/session/branch.go:210-214` |
| Observer fan-out | `observe` wraps inner channel, calls `OnRunStart` then synchronous per-event `OnEvent` in registration order, back-pressures consumer | `pkg/runtime/observer.go:66-81` |
| Tests for composition | Concurrent pinned nested transfers, `PinAgent` isolation, cycle rejection, depth boundary, `RunAgent_NestedBackgroundSubagentStopFiresOnPinnedCaller` | `pkg/runtime/agent_delegation_test.go:944-1190`, `pkg/runtime/agent_delegation_test.go:1343-1416` |
| Background handler tests | Concurrency cap, total cap auto-prune, concurrent `HandleView`/`HandleList`/`HandleStop` race, legacy vs session-aware resolver | `pkg/tools/builtin/agent/agent_test.go:499-513`, `pkg/tools/builtin/agent/agent_test.go:646-664`, `pkg/tools/builtin/agent/agent_test.go:684-729` |

## Answers to Dimension Questions

**When is a parent allowed to become terminal relative to child cleanup?**
Parent `RunStream` may become terminal (emit `StreamStopped` and `close(events)`) while background children are still running. `loop.go:201-211` `finalizeEventChannel` emits `StreamStopped` non-blocking *before* session-end hooks and then `restoreAndClose` closes the parent channel (`pkg/runtime/elicitation.go:137-142`). Foreground `runForwarding` (`pkg/runtime/agent_delegation.go:382-395`) blocks until `childEvents` drains, then unconditionally `parent.AddLiveSubSession(s)` and emits `SubSessionCompleted` (`pkg/runtime/agent_delegation.go:401-402`) — parent does not close until child loop ends. Background `runCollecting` detaches via `context.WithoutCancel` (`pkg/tools/builtin/agent/agent.go:337`) and outlives parent; it persists via `persistBackgroundSubSession` with `context.WithoutCancel(ctx)` (`pkg/runtime/agent_delegation.go:498-521`) even after parent `StreamStopped`. The documented ordering (`loop.go:187-194`) explicitly makes channel close, not `StreamStopped`, the terminal signal, and cleanup hooks run with `context.WithoutCancel`. No barrier waits for background children on parent termination except global `StopAll` at `LocalRuntime.Close()` (`pkg/tools/builtin/agent/agent.go:507-515`).

**Does one child failure cancel siblings, and is that policy explicit?**
No. Sibling cancellation does not occur. `transfer_task` is strictly sequential (`runForwarding` blocks, `loop.go` has no concurrent sibling dispatch), so there are no siblings to cancel. `run_background_agent` tasks are independent entries in `Handler.tasks *concurrent.Map` (`pkg/tools/builtin/agent/agent.go:240`) each with its own `cancel` (`pkg/tools/builtin/agent/agent.go:337-353`) and `atomic Int32 status` (`pkg/tools/builtin/agent/agent.go:138-161`). `HandleRun` checks only caps (`pkg/tools/builtin/agent/agent.go:319-329`); on failure it `storeStatus(taskFailed)` for that task alone (`pkg/tools/builtin/agent/agent.go:412`). No `errgroup.WithContext` links siblings; `runCollecting` `break` on `ErrorEvent` only exits that child's loop (`pkg/runtime/agent_delegation.go:476-479`). The policy is implicitly "isolated" (detached contexts prevent cascade), but no explicit `FailFast` or `CancelSiblings` flag is exposed; `StopAll` only runs at shutdown, not on single failure.

**Are results ordered by submission, start, or completion?**
None deterministically. `transfer_task` returns a single `ToolCallResult` with `s.GetLastAssistantMessageContent()` (`pkg/runtime/agent_delegation.go:411`) — ordering collapses to call-site sequence (one-at-a-time blocking). Background results are `RunResult` per `taskID` (`pkg/tools/builtin/agent/agent.go:73-77`), stored in a `concurrent.Map` whose `HandleList` iterates `Range` in random map order (`pkg/tools/builtin/agent/agent.go:448-456`), and `view_background_agent` requires caller-supplied `taskID` (`pkg/tools/builtin/agent/agent.go:466-482`). No helper sorts by submission, start `startTime`, or completion; `formatView` prints individual tasks unsorted (`pkg/tools/builtin/agent/agent.go:178-234`). The only ordered collection in the codebase is `errgroup` for RAG embeddings (ordered by batch index, `pkg/rag/embed/embed.go:120-137`), not applicable to agent children.

**Can a child outlive or mutate state after its parent is terminal?**
Yes — by design for background agents. `HandleRun` creates `taskCtx` from `context.WithoutCancel(ctx)` (`pkg/tools/builtin/agent/agent.go:337`), documented as "not killed when the parent message context is cancelled" (`pkg/tools/builtin/agent/agent.go:332-336`). `runCollecting` defers `executeSubagentStopHooks` against the parent agent even if `tracedCtx.Err()!=nil` (`pkg/runtime/agent_delegation.go:446-448`, `pkg/runtime/agent_delegation.go:422-424`). After parent channel close, the child still mutates shared state: `parent.AddLiveSubSession(s)` under `sess.mu` (`pkg/session/session.go:1101-1107`, `pkg/runtime/agent_delegation.go:510`), `emitBackgroundEvent(NewTokenUsageEvent(...))` via `OnBackgroundEvent` sink (`pkg/runtime/agent_delegation.go:502-504`), and `persistBackgroundSubSession` writing to `session.Store.AddSubSession` (`pkg/runtime/agent_delegation.go:595-603`). `loop.go:193` notes consumers must rely on channel close, not `StreamStopped`, because hooks run after. A foreground child via `runForwarding` cannot outlive parent because the parent goroutine blocks on `for event := range childEvents` (`pkg/runtime/agent_delegation.go:383-395`), but once the parent is terminal, no further `AddLiveSubSession` occurs; background is the exception.

## Architectural Decisions

- **Two-path composition instead of generic pool**: `SubSessionConfig`/`delegationRequest` centralize construction (`pkg/runtime/agent_delegation.go:115-209`) but split execution into `runForwarding` (interactive, event-forwarding, `swapCurrentAgent`) vs `runCollecting` (detached, `PinAgent`, `NonInteractive=true`). Preserves concrete executor semantics (sequential delegation keeps `currentAgent` semantics; background preserves pinned agent) but avoids abstract `Executor` interface. Evidence: `pkg/runtime/agent_delegation.go:340-440`, `pkg/runtime/skill_runner.go:118-138`.

- **ParentID + DelegationLineage for ownership**: `Session.ParentID` (`pkg/session/session.go:422-425`) gives identity; `DelegationLineage` (`pkg/session/session.go:426-433`) gives ancestor chain for cycle/depth guards. Not persisted (`json:"-"`), re-derived per process. Verified by `validateDelegation` cloning lineage (`pkg/runtime/agent_delegation.go:70`) and `newSubSession` nil-vs-explicit lineage handling (`pkg/runtime/agent_delegation.go:232-235`).

- **Detached contexts for background ownership**: `context.WithoutCancel` + `context.WithCancel` (`pkg/tools/builtin/agent/agent.go:337`) plus `tools.WithoutInteractivePrompts` (`pkg/runtime/loop.go:294-296`) ensure background children survive TUI message cancellation and fail fast on OAuth elicitation instead of hanging. Tradeoff: leaks beyond parent lifetime until explicit `stop_background_agent` or `StopAll`.

- **Event attributability without global bus**: `AgentContext` + `SessionScoped.GetSessionID()` (`pkg/runtime/event.go:18-22`, `pkg/runtime/event.go:25-33`) plus `PersistenceObserver` filtering (`pkg/runtime/persistence_observer.go:74-80`) and per-agent `newAgentContext` keep child events attributable. Background `TokenUsage` forwarded via dedicated `OnBackgroundEvent` (`pkg/runtime/runtime.go:179-184`) not parent channel, avoiding global bus fan-out. Elicitation has dedicated `OnElicitationRequest` sink (`pkg/runtime/elicitation.go:352-356`) for exactly-once delivery across concurrent bridges.

- **Explicit lineage isolation for fan-out**: `validateDelegation` returns `append(parent.DelegationLineageSnapshot(), caller)` freshly allocated (`pkg/runtime/agent_delegation.go:70`) so concurrent `transfer_task` from one parent never shares backing arrays (`pkg/runtime/agent_delegation_test.go:744-754`). Mirrors `mergeExcludedTools` deduplication (`pkg/runtime/agent_delegation.go:279-298`).

- **Budget sharing, not per-child allowance**: First root stream installs `budgetSet` via `ensureBudget()` (`pkg/runtime/loop.go:240-244`) guarded by `budgetMu` (`pkg/runtime/runtime.go:324-328`); sub-sessions reach same `budget` in `recordBudget` (`pkg/runtime/loop.go:841`), so fan-out spends against one wallet (`pkg/runtime/runtime.go:318-324`).

## Notable Patterns

- **Session-pinned agent resolution**: `resolveSessionAgent(sess)` used in `runForwarding`/`runCollecting`/`handleTaskTransfer` to resolve caller from session not `CurrentAgent()` (`pkg/runtime/agent_delegation.go:347-348`, `pkg/runtime/agent_delegation.go:430`, `pkg/runtime/agent_delegation.go:681`). Fix #3886 verified by `TestTransferTask_PinnedParentDoesNotMutateSharedCurrentAgent` (`pkg/runtime/agent_delegation_test.go:944-1025`) and nested background attribution (`pkg/runtime/agent_delegation_test.go:1343-1416`).

- **Deferred `subagent_stop` attribution**: Both paths `defer r.executeSubagentStopHooks(ctx, parent, s, callerAgent, ...)` (`pkg/runtime/agent_delegation.go:378-380`, `pkg/runtime/agent_delegation.go:446-448`) so parent's executor observes completion even on `ErrorEvent` or `ctx.Done()`, with emptiness signalling failure.

- **Live vs embedded cost accounting**: `AddLiveSubSession` sets `liveAttached=true` (`pkg/session/session.go:1101-1107`); `SessionUsage` sums `OwnCost()+EmbeddedSubSessionCost()` but `PersistenceObserver`/`event.go:413-428` exclude live-attached from embedded to avoid double-counting when aggregating per-session `TokenUsage` events.

- **Non-blocking StreamStopped**: `nonBlocking(&channelSink{ch: events}).Emit(StreamStopped(...))` (`pkg/runtime/loop.go:211`) with `event_sink.go:36` `recover()` swallow — prevents deadlock #3070 (`pkg/runtime/loop.go:196-201`).

- **Elicitation waiters map**: `elicitationWaiters.register` before `emitElicitationRequest` (`pkg/runtime/elicitation.go:512-530`) fixes TOCTOU; `elicitationBridge.send` holds `RLock` across channel send with `ctx.Done` select and `restoreAndClose` holds `Lock` to avoid send-on-closed panic (`pkg/runtime/elicitation.go:104-142`).

## Tradeoffs

- **Simplicity vs generality**: Single-child blocking + map-keyed detached tasks are easy to reason about and test (#3886 lineage tests), but no `Group`/`Pool` abstraction exists to express "run 5 helpers, collect in submission order, fail-fast, limit 3 concurrent" — callers must orchestrate via LLM tool calls (`run_background_agent` → poll `view_background_agent`).
- **Outlive guarantee vs resource leak**: `WithoutCancel` prevents premature cancellation and enables parallel background work during TUI interaction, but children hold `session.Store` writes and `wg` references beyond parent; no per-parent `WaitGroup` or `context.Context` tree cancels them — only global `maxConcurrentTasks`/`StopAll` bound resource use.
- **Strong typing vs erasure**: `SubSessionConfig` preserves `Permissions`, `SafetyPolicy`, `AllowedTools`, `ExtraToolSets` per child (`pkg/runtime/agent_delegation.go:132-168`), so tool semantics survive delegation. Cost: every new delegation feature adds fields to config rather than a generic `Opts` bag.
- **Deterministic lineage vs unordered store**: Lineage ordering is deterministic (slice append), but task storage is `concurrent.Map` — background result observation is non-deterministic, pushing ordering responsibility to higher layers.
- **Isolation vs observability**: `PersistenceObserver` dropping `SessionScoped` mismatches (`pkg/runtime/persistence_observer.go:78-79`) keeps parent transcript clean, but means sub-session streaming `AgentChoice` events never reach parent's store via that observer — only final `SubSessionCompleted` / background `TokenUsage` survive, losing incremental child observability in persistence.

## Failure Modes / Edge Cases

- **Sibling failure silent**: One `run_background_agent` failure `storeStatus(taskFailed)` (`pkg/tools/builtin/agent/agent.go:412`) does not cancel or notify sibling tasks; caller must poll each `taskID` individually (`pkg/tools/builtin/agent/agent.go:466-482`). No fail-fast group exists; simultaneous failures reported only via individual `view_background_agent` `Error` fields.

- **Depth exhaustion masquerading as tool error**: Exceeding `maxDelegationDepth` returns `tools.ResultError` not Go `error` (`pkg/runtime/agent_delegation.go:688-691`, `pkg/tools/builtin/agent/agent.go:650-652`), so `RunStream` tool dispatch treats it as model-visible tool output rather than run error — may be retried by LLM.

- **Concurrent cap error ordering**: `HandleRun` checks `runningTaskCount >= maxConcurrentTasks` before `totalTaskCount` pruning (`pkg/tools/builtin/agent/agent.go:319-329`); when both caps are hit, "maximum concurrent" error fires even if pruning would free slots — tested at `pkg/tools/builtin/agent/agent_test.go:666-680`.

- **Panic in child tool not propagated**: Tool handlers (`probeTool` at `pkg/runtime/agent_delegation_test.go:899-908`) have no `recover`; panic in `RunAgent` goroutine (`pkg/tools/builtin/agent/agent.go:356`) would crash handler goroutine and leak `task` in `running` state with no `CAS` to `failed` — `StopAll` would wait via `wg.Wait()` but status remains `running`.

- **Output truncation silent**: `t.writeOutput` respects `maxOutputBytes` cap by dropping excess (`pkg/tools/builtin/agent/agent.go:165-173`); `formatView` adds `[output truncated at 10MB limit]` (`pkg/tools/builtin/agent/agent.go:221`), but `RunResult.Result` via `GetLastAssistantMessageContent()` (`pkg/runtime/agent_delegation.go:527`) is not truncated — inconsistency between live buffer and final result.

- **Channel full deadlock avoided but lossy**: `StreamStopped` non-blocking loss (`pkg/runtime/loop.go:211`) means consumers relying on that event may miss it and rely on channel close fallback; background bridge `go func() { r.elicitation.send(...) }` (`pkg/runtime/elicitation.go:541-545`) also best-effort, with sync sink as reliable path — dual delivery invariant fragile if second path added.

- **Branch after concurrent append race**: `BranchSession` snapshots via `snapshotItems()` under `mu` (`pkg/session/branch.go:43`), but `recalculateSessionTotals` sums `Usage` from messages (`pkg/session/branch.go:394-412`) which may have been updated concurrently via `SetTokensAndCost` — window between snapshot and recalc could tear totals.

- **Budget double-count risk**: `runForwarding` parent and child share `budget` (`pkg/runtime/runtime.go:324-328`); if child emits its own `BudgetUsageEvent` with sub-session `SessionID` but parent also aggregates via `SessionUsage` (`pkg/runtime/event.go:413`), UI aggregating per-session budgets must deduplicate by session ID — not enforced in `BudgetUsage` construction (`pkg/runtime/event.go:763-797`).

## Future Considerations

- Add a generic `Group`/`Pool` helper (e.g., `pkg/runtime/composition` wrapping `errgroup.WithContext` + `semaphore.Weighted`) that enforces `validateDelegation` lineage, `PinAgent` semantics, deterministic ordered collection by submission index, and explicit `FailFast` vs `CollectAll` policy — would close Phase-12 gap without breaking `runForwarding`/`runCollecting` callers.
- Introduce `ctx` tree for background children tied to parent's `WithoutCancel` descendant with `CancelCause` linking sibling failure to sibling `cancel` when fail-fast enabled; surface as explicit `SubSessionConfig.FailFast bool`.
- Preserve result ordering: store `submissionOrder int` in `task` and sort `HandleList` / future `CollectResults` by it; add `OrderedResults []RunResult` API alongside `taskID` map for `collect-all` pattern.
- Harden panic propagation: `defer recover()` in `Handler.wg.Go` body (`pkg/tools/builtin/agent/agent.go:356`) converting panic to `storeStatus(taskFailed)` + `errMsg`; similarly wrap `toolexec.Dispatcher` parallel tool execution (`pkg/runtime/toolexec/dispatcher.go:95-235`) which currently not shown to recover.
- Tie background lifecycle to parent: optionally register background tasks in parent session's `liveSessions` entry and expose `WaitForChildren(ctx)` barrier before parent `StreamStopped`; or make `StopAll` per-session scoped, not runtime-global.
- Make `StreamStopped` reliable via sequence counter rather than best-effort drop, or document consumer contract to await `SubSessionCompleted` + channel close for full child persistence before considering parent clean.

## Questions / Gaps

- No evidence found for panic handling tests for child tool execution — searched `pkg/runtime` for `panic`/`recover` (`grep` returned only channel guards `event_sink.go:36`, `elicitation.go:108` and loop detector, not child panic recovery). What is expected child panic surface per dimension?
- No evidence of simultaneous child failure test — delegation tests cover sequential depth/cycle (`agent_delegation_test.go:602-1455`) and concurrent pinned transfers staying isolated (`agent_delegation_test.go:1087-1190`) but not "two background children fail concurrently, verify sibling cancellation or collect-all error aggregation."
- Bounded-pool execution for children not found — `maxConcurrentTasks` is admission control returning tool error, not a semaphore-bound `errgroup.SetLimit` pool that queues and runs with concurrency limit; is queueing vs rejection the intended Phase-12 bounded-pool semantics?
- Cannot answer deterministic result ordering — background tasks have no ordered collection API; `concurrent.Map.Range` order is random (Go map semantics). Is ordering by submission/start/completion intended to be provided by a new helper or by `session.Messages` item order (which is completion order via `AddLiveSubSession`)?
- Child outlive mutation window not quantified — `persistBackgroundSubSession` writes with detached context (`pkg/runtime/agent_delegation.go:498-502`) but no test asserts writes after parent `StreamStopped`; does persistence observer handle `AddSubSession` race with `BranchSession` `setParentIDs` (`pkg/session/branch.go:381-391`)?
- No evidence of parent becoming `terminal` before child cleanup completing being *prevented* — `runForwarding` does wait, but no `LocalRuntime` test asserts parent `Run` does not return before `SubSessionCompleted` persistence succeeds (persistence observer is async via `observe` goroutine `pkg/runtime/observer.go:70-80`).

---

Generated by `01.11-execution-composition-and-child-ownership` against `docker-agent`.
