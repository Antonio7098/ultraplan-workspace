# Source Analysis: docker-agent

## 01.08 Bounded Agent Loop, Context, and Evidence

### Source Info

| Field | Value |
|-------|-------|
| Name | docker-agent |
| Path | `studies/aren-go-runtime-study/sources/docker-agent` |
| Language / Stack | Go 1.24 / `pkg/runtime` + `pkg/session` + `pkg/compaction` + `pkg/model/provider` + `pkg/tools` |
| Analyzed | 2026-08-29 |

## Summary

docker-agent implements a sequential `RunStream` -> `runTurn` loop in `pkg/runtime/loop.go` that mediates user input -> model stream (`pkg/runtime/streaming.go:handleStream`) -> `toolexec.Dispatcher.Process` -> tool results -> compaction/budget/iteration checks -> next turn or `Stopped`. Bounding is layered rather than single-guard: iteration caps (`pkg/runtime/loop_steps.go:enforceMaxIterations`), degenerate-loop detector (`pkg/runtime/toolexec/loop_detector.go`), cost/token/wall-clock budgets (`pkg/runtime/budget.go`), context-window compaction (`pkg/compaction/compaction.go:ShouldCompact` + `pkg/runtime/compactor/compactor.go`), per-result and old-content token truncation (`pkg/session/session.go:capToolResultContent`/`truncateOldToolContent`), and a 5-minute stream idle timeout. Context compaction is the most engineered piece: calibrated heuristic estimator, user/assistant-boundary snapping, synthetic `Session Summary:` messages, and `sanitizeToolCalls` to preserve provider contracts. Evidence preservation is weak: the transcript is durable (`pkg/runtime/persistence_observer.go`), but compaction replaces verbatim tool outputs with an LLM summary and the item schema carries no provenance, citation, or artifact pointer beyond `ToolCallID` and `Cost/Usage`.

## Rating

**6 / 10**

Rationale: Loop termination is multi-redundant (iteration + loop-detector + budget + overflow bounding) and transcript invariants survive compaction thanks to `sanitizeToolCalls` and boundary snapping. Deductions for (1) budgets are advisory and allow one-turn overshoot (`recordBudget` post-turn, `enforceBudget` pre-next-turn), (2) compaction deliberately discards verbatim evidence supporting a later answer with lossy summary, (3) no first-class provenance/source-reference/artifact graph, and (4) cancellation reaches the *currently active* phase (stream *or* tool batch) but stream and tools never run concurrently so there is no simultaneous fan-out guarantee.

## Evidence Collected

Every entry MUST include a file path with line numbers. Format: `path/to/file.ts:NN`.

| Area | Evidence | File:Line |
|------|----------|-----------|
| Loop entry | `RunStream` creates buffered `events` chan, `ensureBudget()`, registers `liveSessionEntry`, spawns `runStreamLoop` goroutine, returns `observe` wrapper | `studies/aren-go-runtime-study/sources/docker-agent/pkg/runtime/loop.go:230-260` |
| Iteration bound | `enforceMaxIterations` checks `iteration < max`; non-interactive auto-stops with assistant message; interactive blocks on `resumeChan`, approving raises cap by +10 | `studies/aren-go-runtime-study/sources/docker-agent/pkg/runtime/loop_steps.go:48-127` |
| Loop check site | `newMax, decision := enforceMaxIterations(...); if decision==iterationStop {return}` at top of each turn, before budget and model call | `studies/aren-go-runtime-study/sources/docker-agent/pkg/runtime/loop.go:498-503` |
| Degenerate loop detector | `NewLoopDetector(threshold, exempt...)` canonical JSON signature, exempt polling tools never count, `Record` returns true at `consecutive >= threshold`; reset on agent switch; default threshold 5 | `studies/aren-go-runtime-study/sources/docker-agent/pkg/runtime/toolexec/loop_detector.go:24-91` |
| Detector wiring | `NewLoopDetector(loopThreshold, viewBackground..., viewBackgroundJob)` and `Record(res.Calls)` -> emit `ErrorCodeLoopDetected` + `SetStatus(codes.Error)` + `turnEndReasonLoopDetected` | `studies/aren-go-runtime-study/sources/docker-agent/pkg/runtime/loop.go:441-449` and `928-956` |
| Budget enforcement | `enforceBudget` calls `currentBudget().exceededFor(agentName)`; on breach emits `BudgetExceeded` + `notifyBudgetExceeded` + appends assistant stop message | `studies/aren-go-runtime-study/sources/docker-agent/pkg/runtime/budget.go:338-373` |
| Budget recording | `recordBudget` adds `InputTokens+OutputTokens` and `cost` and `active` to every tracker for the agent's budgets; warns once on `unpricedSpend`; emits `BudgetUsage` | `studies/aren-go-runtime-study/sources/docker-agent/pkg/runtime/budget.go:385-407` |
| Budget call sites | `if r.enforceBudget(...)==iterationStop {streamReason=budgetExceeded; return}` before `ls.iteration++`; `r.recordBudget(sess,a,res.Usage,msgCost,...)` immediately after successful stream | `studies/aren-go-runtime-study/sources/docker-agent/pkg/runtime/loop.go:507-512` and `841-842` |
| Context window resolution | `resolveContextLimit` prefers `provider_opts.context_size` then `modelsdev` catalogue; `effectiveContextLimit` caps primary window by dedicated compaction model's smaller window | `studies/aren-go-runtime-study/sources/docker-agent/pkg/runtime/session_compaction.go:237-249` and `275-295` |
| Proactive compaction trigger | `if contextLimit>0 && sessionCompactionEnabled(a) && compaction.ShouldCompact(input+output+added, contextLimit, threshold) {compactWithReason(threshold)}` before model call | `studies/aren-go-runtime-study/sources/docker-agent/pkg/runtime/loop.go:569-572` |
| Compaction estimator | `ShouldCompact(total > threshold*limit)` with `DefaultThreshold=0.9`; `NewEstimator` calibrates `chars/3.5` heuristic against provider `prompt` deltas, clamped `0.75-2.0`, `calibrationMinTokens=512` | `studies/aren-go-runtime-study/sources/docker-agent/pkg/compaction/compaction.go:19-77` and `99-172` |
| Keep boundary | `SplitIndexForKeep(messages, keepTokenBudget)` walks backward, accumulates calibrated estimates, snaps to `MessageRoleUser/Assistant` boundary, returns `len(messages)` when everything fits | `studies/aren-go-runtime-study/sources/docker-agent/pkg/compaction/compaction.go:279-313` |
| Compactor orchestration | `compactor.RunLLM` clones compaction model (`WithMaxTokens(summaryTokenBudget)`, `WithCompacting`), calls `extractMessages`, creates throwaway `compactionSession`, `RunAgent`, extracts `lastAssistantContentAfter(seedLen)`; `ComputeFirstKeptEntry` shares kept-tail policy for hook-supplied summaries | `studies/aren-go-runtime-study/sources/docker-agent/pkg/runtime/compactor/compactor.go:148-288` |
| Session compaction apply | `doCompact` dispatches `BeforeCompaction` (veto or hook summary), emits `SessionCompaction:started`, runs `compactor.RunLLM` or hook result, snapshots `preInput/preOutput`, `sess.ApplyCompaction(inputTokens,0,Item{Summary,FirstKeptEntry,Cost,Model,Usage})` + `sessionStore.UpdateSession` + `SessionSummary` + `AfterCompaction` | `studies/aren-go-runtime-study/sources/docker-agent/pkg/runtime/session_compaction.go:65-152` |
| Transcript synthesis | `buildSessionSummaryMessages` finds `lastSummaryIndex`, emits synthetic user message `SummaryMessageContent(summary)` (`"Session Summary: "+summary`), `startIndex = kept>0 ? kept : lastSummary+1` | `studies/aren-go-runtime-study/sources/docker-agent/pkg/session/session.go:2003-2037` and `1969-1978` |
| Transcript invariants | `sanitizeToolCalls` injects synthetic `No result provided` for pending tool calls before next user/assistant or EOF, drops orphaned `tool` messages lacking pending `tool_use`, drops duplicates; comments explicitly cite compaction boundary producing Bedrock rejection | `studies/aren-go-runtime-study/sources/docker-agent/pkg/session/session.go:2434-2495` |
| Trimming invariant | `trimMessages` protects system+user messages, marks assistant `ToolCalls` to remove, skips orphaned `tool` results whose assistant was trimmed | `studies/aren-go-runtime-study/sources/docker-agent/pkg/session/session.go:2290-2362` |
| Tool output caps | `AddMessage` calls `capToolResultContent` (middle-out truncation, head+`[...truncated...]`+tail) per result; `GetMessages` reapplies cap as read-time backstop; `truncateOldToolContent` replaces older tool results with `[content truncated]` when `MaxOldToolCallTokens>0` | `studies/aren-go-runtime-study/sources/docker-agent/pkg/session/session.go:829-842` and `2553-2630` and `2506-2535` |
| Agent/system caps | Agent fields `maxIterations`, `maxConsecutiveToolCalls`, `maxOldToolCallTokens`, `maxToolResultTokens` via `pkg/agent/agent.go:48-54` and `pkg/agent/opts.go:193-232`; applied at session (`pkg/session/session.go:326-352`) and at prompt assembly |
| Streaming & context cancel | `handleStream` spawns recv goroutine multiplexed with `select {case recvCh, case ctx.Done(), case idleTimer}`; `ctx.Done()` returns `ctx.Err()` promptly; idle timeout `5m` calls `cancelStream(errStreamIdle)` to close TCP | `studies/aren-go-runtime-study/sources/docker-agent/pkg/runtime/streaming.go:73-104` and `337-354` |
| Tool batch cancellation | Batch `batchCtx, cancelBatch := context.WithCancelCause(ctx)`; tool `Cancel` on `outcome.Canceled` or `StopRun`; per-call early check `if ctx.Err()!=nil {errorResponse(cancellationMessage); return cancellationOutcome}`; `askUser` checks `ctx.Err()` before and while blocked on `Resume` | `studies/aren-go-runtime-study/sources/docker-agent/pkg/runtime/toolexec/dispatcher.go:254-266` and `377-382` and `811-845` |
| Tool handler cancellation | `translateError` maps `context.Canceled` to user-visible `ResultError(cancellationMessage)` with `codes.Ok`; otherwise `RecordError` + `codes.Error` | `studies/aren-go-runtime-study/sources/docker-agent/pkg/runtime/toolexec/dispatcher.go:1111-1132` |
| Message cost ledger | `computeMessageCost(usage, model)` folds `Input/CachedInput/CacheWrite/Output * Cost/1e6`, returns `*float64` (nil = unpriced); single source fed to `after_llm_call` hook payload and persisted assistant message | `studies/aren-go-runtime-study/sources/docker-agent/pkg/runtime/loop.go:1121-1130` |
| Persistence / evidence | `PersistenceObserver.OnEvent` handles `AgentChoice` (streaming row `UpdateMessage`), `MessageAddedEvent`/`UserMessageEvent` -> `store.AddMessage`, `SessionSummaryEvent` -> `AddSummary`, `TokenUsageEvent` -> `UpdateSessionTokens`, `SubSessionCompletedEvent` -> `AddSubSession`, `ErrorEvent` -> `AddError` | `studies/aren-go-runtime-study/sources/docker-agent/pkg/runtime/persistence_observer.go:74-164` |

## Answers to Dimension Questions

- **What guarantees that a loop terminates even when the model keeps requesting tools?**

  Multiple independent ceilings, checked in strict order at the iteration boundary and after each batch:

  1. **Iteration cap** (`studies/aren-go-runtime-study/sources/docker-agent/pkg/runtime/loop_steps.go:48-58`, `studies/aren-go-runtime-study/sources/docker-agent/pkg/runtime/loop.go:498-503`). `enforceMaxIterations` is evaluated before every model call. `session.MaxIterations==0` means “no limit” at the session layer, but the runtime always has a value (default configured at team load). Non-interactive sessions auto-stop with an assistant `Execution stopped after reaching max_iterations`; interactive ones block on `resumeChan` awaiting approve/reject. Passing `MaxIterations` as `session` config rather than budget makes it visible to the `before_llm_call:max_iterations` hook as well.

  2. **Degenerate identical-call detector** (`studies/aren-go-runtime-study/sources/docker-agent/pkg/runtime/toolexec/loop_detector.go:54-77`, `studies/aren-go-runtime-study/sources/docker-agent/pkg/runtime/loop.go:928-956`). `LoopDetector` records a canonical JSON signature per batch; after `N` consecutive identical batches (`N= sess.MaxConsecutiveToolCalls` or default 5 via `pkg/agent/opts.go:202-208` and `pkg/runtime/loop.go:441-449`) the runtime emits `ErrorCodeLoopDetected`, marks the session span `codes.Error`, calls `notifyError`, and exits with `turnEndReasonLoopDetected`. Purely polling batches (`view_background_agent`, `list_background_agents`, `view_background_job`) are exempt and invisible to the counter.

  3. **Budget ceilings** (`studies/aren-go-runtime-study/sources/docker-agent/pkg/runtime/budget.go:338-383`). At the same boundary, `enforceBudget` checks `cost >= maxCost || tokens >= maxTokens || active >= maxTime` for every budget that applies to the current agent (run budget + named budgets). On breach it emits `BudgetExceeded` and a stop assistant message. Overflow recovery is separately bounded: `handleStreamError` (`studies/aren-go-runtime-study/sources/docker-agent/pkg/runtime/loop_steps.go:179`) allows at most `maxOverflowCompactions` (default 1, `pkg/runtime/runtime.go:386`) consecutive context-overflow auto-compactions.

  4. **Context-overflow retry bound** (`studies/aren-go-runtime-study/sources/docker-agent/pkg/runtime/loop_steps.go:179-196`) prevents infinite compact-and-retry loops when compaction cannot shrink below the window.

  5. **Stream liveness** (`studies/aren-go-runtime-study/sources/docker-agent/pkg/runtime/streaming.go:19-39`, `342-354`) bounds the upstream call itself to 5 minutes idle before cancellation.

  No single guard is formally proven, but together they make an infinite tool-request loop require defeating iteration, detector, and budget simultaneously. Gap: if all numeric limits are deliberately set to 0/unbounded and the budget is omitted, only the loop detector remains, which is triggerable only on exact-signature repetition, not on continually varying args.

- **Can compaction orphan a tool call or remove the evidence supporting the final result?**

  **Orphaning at the wire level — explicitly prevented.** `pkg/compaction/compaction.go:294-313` (`SplitIndexForKeep`) and `pkg/session/session.go:2003-2037` deliberately snap the kept-tail to `user`/`assistant` boundaries so the tail never starts mid-tool-exchange. `pkg/session/session.go:2434-2495` (`sanitizeToolCalls`) is the invariant backstop applied on every `GetMessages` before the provider call: it synthesizes `No result provided` for missing results, drops orphaned `tool` results (whose `ToolCallID` never appeared in the preceding assistant `ToolCalls`), and drops duplicates. The comment on line 2429 explicitly calls out compaction's kept-tail boundary landing between an assistant `tool_use` and its result as the motivating case, with Bedrock's `toolResult/toolUse` count strictness as the failure mode.

  **Evidence supporting the final result — yes, compaction discards it verbatim.** `pkg/session/session.go:2003-2037` and `pkg/runtime/session_compaction.go:128-139` replace every item before `FirstKeptEntry` with a single `Item{Summary}`. `pkg/runtime/compactor/compactor.go:285-311` keeps at most `keepTokenBudget` (default `20k`, scaled to `contextLimit/5`) verbatim. Everything earlier — including the tool outputs that the final assistant message may be summarizing — is gone except for whatever the LLM summary paraphrases. The original tool payloads, document attachments, and per-message `Usage`/`Cost` are not retained alongside the summary except as an aggregate `Usage` on the summary item (`pkg/session/session.go:176-177`, `pkg/runtime/session_compaction.go:157-163`). There is no citation graph; `chat.Message` only links `ToolCallID` (`pkg/session/session.go:517-527`) and session items have no `sourceRef`/`provenance` field.

- **Which budgets are enforced before work begins and which are only reported afterward?**

  | Budget / limit | Value & resolution | Enforcement point | Pre-work or post-hoc | File:Line |
  |---|---|---|---|---|
  | `maxIterations` | `session.MaxIterations` or agent `MaxIterations`, routed via `loopState.maxIterations` | Checked before each model call; stop or wait for user | **Pre-work** (hard gate) | `studies/aren-go-runtime-study/sources/docker-agent/pkg/runtime/loop_steps.go:56-58` |
  | `maxConsecutiveToolCalls` | `agent.MaxConsecutiveToolCalls` (0→5), `LoopDetector.threshold` | Checked after tool batch; stop | Post-batch (detected after spend) | `studies/aren-go-runtime-study/sources/docker-agent/pkg/runtime/toolexec/loop_detector.go:54-77` |
  | `maxCost`/`maxTokens`/`maxTime` (run + named budgets) | `latest.BudgetConfig` → `budgetTracker`, shared `budgetSet` via `WithBudget`/`WithNamedBudgets` | `enforceBudget` before each iteration; `recordBudget` after each priced turn adds to trackers | **Both**: breach stops *next* turn; spend of *just-completed* turn is recorded post-hoc (`cost/tokens/active` addition) allowing one-turn overshoot | `studies/aren-go-runtime-study/sources/docker-agent/pkg/runtime/budget.go:338-373` (gate) and `385-407` (record) |
  | Context window / `compactionThreshold` | `resolveContextLimit` vs `effectiveContextLimit`, `compaction.ShouldCompact` `0.9` default | `ShouldCompact` before model call triggers proactive `doCompact`; overflow path retries bounded to 1 | Pre-work (compaction decision) | `studies/aren-go-runtime-study/sources/docker-agent/pkg/compaction/compaction.go:69-77` |
  | `maxToolResultTokens` | `session.MaxToolResultTokens` | Middle-out truncation at `AddMessage` (`capToolResultContent`) and re-applied at `GetMessages` read time, before provider | Pre-send truncation, not a stop budget | `studies/aren-go-runtime-study/sources/docker-agent/pkg/session/session.go:829` and `2553-2604` |
  | `maxOldToolCallTokens` | `session.MaxOldToolCallTokens` | `truncateOldToolContent` during prompt assembly before provider | Pre-send truncation | `studies/aren-go-runtime-study/sources/docker-agent/pkg/session/session.go:2506-2535` |
  | Output length (`FinishReasonLength`) | Provider `finish_reason` | Surfaced via `streamResult.FinishReason`, triggers warning path but no iteration gate | Reported only | `studies/aren-go-runtime-study/sources/docker-agent/pkg/runtime/streaming.go:274-302` |
  | Per-turn cost pricing | `computeMessageCost` `m.Cost * tokens/1e6` | Computed after stream, written to message `Cost` and `recordBudget` | Reported (nil = unpriced) | `studies/aren-go-runtime-study/sources/docker-agent/pkg/runtime/loop.go:1121-1130` |

  The cost/token/time budgets are the only ones where the unit of enforcement (next iteration) lags the unit of accounting (just-finished turn). The pre-work `enforceBudget` gate therefore always evaluates on *prior* spend, so the run can exceed the declared cap by up to one priced turn.

- **Does cancellation reach an active model stream and active tool at the same time?**

  No — by construction the two phases are sequential, so they are never both active simultaneously, but cancellation does propagate promptly to whichever is active.

  - **Model stream**: `handleStream` (`studies/aren-go-runtime-study/sources/docker-agent/pkg/runtime/streaming.go:89-104`, `337-340`) runs `stream.Recv()` in a dedicated goroutine multiplexed with `select { case ctx.Done(): return ctx.Err() }`. The parent `RunStream` context is derived from `httpclient.ContextWithSessionID` and `genai.WithConversationID`; a `cancel()` delivers `context.Canceled` to the stream path within one `select` iteration. The `idleTimer` path adds a second route: on stall it calls `cancelStream(errStreamIdle)` which cancels the HTTP transport context and unblocks the `Recv` reader.

  - **Tool batch**: `Dispatcher.Process` (`studies/aren-go-runtime-study/sources/docker-agent/pkg/runtime/toolexec/dispatcher.go:254-266`) derives `batchCtx, cancelBatch := context.WithCancelCause(ctx)`. Every parallel tool call runs against `batchCtx`; `cancelBatch` is invoked on `outcome.Canceled` or `StopRun`, aborting siblings. Per-call, `call.run` (`378-382`) checks `ctx.Err()` before dispatch, and `askUser` (`811-845`) checks `ctx.Err()` both before emitting the confirmation event and while blocked on `Resume`, emitting `canceled` with `ApprovalSourceContextCanceled`. `invoke`/`translateError` (`1111-1132`) map handler `context.Canceled` to a user-visible cancellation rather than a hard error, and on entry each tool's span is annotated accordingly.

  - **Loop-level polling**: `runTurn` defers `turn_end` with `context.WithoutCancel(ctx)` so it still fires (`pkg/runtime/loop.go:758`), and top-of-loop `if err:=ctx.Err()` (`pkg/runtime/loop.go:515`) bails before starting the next model call.

  Simultaneous fan-out would require a model stream and a tool dispatch sharing the same instant; the loop alternates `fallback.execute` (which owns the stream) and `processToolCalls` (which owns the tool batch), so the answer is “each phase is cancelable, but not jointly.” A mid-stream cancel prevents the subsequent tool phase from starting; a mid-tool cancel aborts siblings. There is a narrow case where filesystem/background-job polling tools running as part of a long batch could survive `ctx` cancellation long enough to outlive `observe` channel shutdown, but the parent `finalizeEventChannel`'s non-blocking `StreamStopped` + `channel close` ensures the UI does not deadlock.

## Architectural Decisions

- **Loop as `LocalRuntime` with per-RunStream `loopState`** (`studies/aren-go-runtime-study/sources/docker-agent/pkg/runtime/loop.go:628-653`). Mutable per-run state (`iteration`, `toolModelOverride`, `loopDetector`, `overflowCompactions`, `prevTurnMadeToolCalls`) lives in a struct threaded by pointer, keeping `runTurn` signatures stable while the run loop owns persistent policy. Tradeoff: concurrency is limited to tool-batch parallelism; the LLM stream is strictly sequential.

- **Transcript middleware chain** (`studies/aren-go-runtime-study/sources/docker-agent/pkg/session/session.go:2434-2495`, `2290-2362`, `2553-2604`). Assembly is `buildInvariantSystemMessages` + instruction context + summary + conversation interleaving + `NumHistoryItems` trim + `truncateOldToolContent` + `normalizeMessageContent` + `sanitizeToolCalls` + `capToolResultContent`. Each stage is a pure function on `[]chat.Message` slices, enabling unit tests like `TestSanitizeToolCalls` and `TestTrimMessagesWithToolCalls`.

- **Compaction as separate `pkg/runtime/compactor` package** (`studies/aren-go-runtime-study/sources/docker-agent/pkg/runtime/compactor/compactor.go:1-18`). Import direction is strictly `runtime -> compactor`, avoiding cycles. The runtime supplies a `RunAgent` callback so the compactor can spin a throwaway `session.New` sub-runtime. Keeps the glossary of token budgets (`maxKeepTokens=20k`, `MaxSummaryTokens=16k`, scaled to `contextLimit/5` and `/4`) out of the main loop.

- **Calibrated estimator over raw `len/4`** (`studies/aren-go-runtime-study/sources/docker-agent/pkg/compaction/compaction.go:99-165`). Uses provider `prompt` deltas between consecutive assistant usage anchors, clamped to `512` heuristic mass and `0.75-2.0` ratio, to reconcile `charsPerToken=3.5` heuristics with actual tokenizers. Bias is toward compacting slightly early.

- **Budget as shared `budgetSet`** (`studies/aren-go-runtime-study/sources/docker-agent/pkg/runtime/budget.go:229-320`). A single `budgetSet` is lazily installed on first root stream (`ensureBudget`) and shared across sub-sessions via `message_queue` fan-out, so delegated agents spend against the root's wallet. Guarded by `budgetMu`.

## Notable Patterns

- **Synthetic summary message**: `Session Summary: …` (`studies/aren-go-runtime-study/sources/docker-agent/pkg/session/session.go:1969-1978`) is a *user*-role message so the conversation stays conversational rather than introducing a synthetic system turn that providers might cache differently.
- **Orphan sanitization as read-time invariant**: `sanitizeToolCalls` and `trimMessages` both guarantee `toolResult` ↔ `toolUse` duality, satisfying Anthropic/Bedrock strict validation without persisting synthetic rows.
- **Exempt polling loop-detector**: `NewLoopDetector`'s `isExemptBatch` (`studies/aren-go-runtime-study/sources/docker-agent/pkg/runtime/toolexec/loop_detector.go:79-91`) prevents background-job polling (`view_background_agent`) from masking a genuine tight loop.
- **Before-compaction hook injection**: `summaryFromHook` (`studies/aren-go-runtime-study/sources/docker-agent/pkg/runtime/session_compaction.go:175-191`) lets a hook supply `Summary` text verbatim, reusing `ComputeFirstKeptEntry` so the kept-tail policy stays identical across LLM and hook strategies.
- **Spend vs gate separation**: `recordBudget` after the model call, `enforceBudget` before the next call, mirrors classic “account then police” quota design.
- **Context-limit capping to compaction model** (`studies/aren-go-runtime-study/sources/docker-agent/pkg/runtime/session_compaction.go:221-249`) ensures the UI gauge and proactive threshold never describe a window that the summary call cannot ingest.

## Tradeoffs

- **Strong bounding, lossy evidence.** The vertically integrated caps make an unbounded run unlikely, but verbatim tool evidence needed to audit the final answer is sacrificed for prompt fitting. For Aren Phase 9/10 (where tool results *are* the evidence), this is a regression risk.
- **Post-hoc cost budgets.** Guarantees are “best-effort stop after overshoot” rather than “never exceed.” A single expensive turn (e.g., a large `MaxTokens` setting) can double the intended cap before the next `enforceBudget` fires. `unpricedSpend` detection (`studies/aren-go-runtime-study/sources/docker-agent/pkg/runtime/budget.go:205-212`) mitigates silent under-counting but only as a warning.
- **Dual compaction trigger (proactive threshold + reactive overflow).** Proactive `ShouldCompact` fires at ~90% but is optional; reactive overflow retries once. The conservative estimator (compact early) trades extra summarization cost for fewer hard overflows. Small windows (local GGUF with `provider_opts.context_size: 8k`) are handled via scaled budgets, but scaling down `keepBudget` to `contextLimit/5` leaves less recent context verbatim.
- **Sequential phases, not concurrent tool+stream.** Simplicity of cancellation and budget accounting over throughput: a long tool batch still blocks the next `cancel` delivery until the batch context is observed, unlike an executor that multiplexes stream deltas and tool I/O concurrently.
- **No provenance chain.** The schema stores flat `Item{Message, SubSession, Summary, Cost, Model, Usage}` (`studies/aren-go-runtime-study/sources/docker-agent/pkg/session/session.go:136-184`) but no edge from an assistant answer to the tool-result slices that justify it. Evidence consumers must grep the transcript.

## Failure Modes / Edge Cases

- **One-turn budget overrun**: Because `recordBudget` (`studies/aren-go-runtime-study/sources/docker-agent/pkg/runtime/budget.go:385`) happens after the turn, a single turn can spend arbitrarily beyond `maxCost`/`maxTokens`. Tests assert only that the *next* turn stops.
- **Persistent `InputTokens` after compaction**: `ApplyCompaction(inputTokens, 0, item)` resets `InputTokens` to the summary estimate (`studies/aren-go-runtime-study/sources/docker-agent/pkg/session/session.go:921-927`, `pkg/runtime/session_compaction.go:133-139`) but the summary token estimate is heuristic (`EstimateMessageTokens`). Under-estimation defers the next proactive compaction; over-estimation triggers extra compaction.
- **Orphan tool-result drop as silent data loss**: `sanitizeToolCalls` silently drops duplicated/orphaned tool results after a compaction boundary that split an exchange. No warning event is emitted; downstream tool logic cannot distinguish “result was dropped at assembly” from “tool never ran.”
- **Unpriced spend invisible to cost caps**: Models absent from `modelsdev` and without `cost:` config return `nil` cost (`studies/aren-go-runtime-study/sources/docker-agent/pkg/runtime/loop.go:1122-1123`), so `budgetTracker.unpriced` stays true and `enforceBudget` may never trip `maxCost` even as real provider spend accrues. Only a warning is emitted.
- **Context limit 0 disables compaction entirely**: `ShouldCompact` returns `false` when `contextLimit==0` (`studies/aren-go-runtime-study/sources/docker-agent/pkg/compaction/compaction.go:70`). Runs against an uncatalogued local model with no `provider_opts.context_size` will never auto-compact and will fail with `ContextOverflowError` up to `maxOverflowCompactions` times before surfacing.
- **Hook veto stranding**: `BeforeCompaction` returning `Decision: block` results in a `compaction:started` without a terminal `completed` (`skipped`/`failed`/`applied`) on the same code path when no summary was attempted (`studies/aren-go-runtime-study/sources/docker-agent/pkg/runtime/session_compaction.go:71-76`). Consumers tolerating unpaired terminals must handle this.
- **Cancellation vs `turn_end` ordering**: `runTurn`'s deferred `executeTurnEndHooks` runs even on `ctx` cancel via `WithoutCancel`, but `streamSpan.End()` happens twice-guarded inline; a panic-recovered branch that skips the deferred read of `ctx.Err` could report `turnEndReasonNormal` rather than `canceled`.
- **Token truncation middle-out may hide error details**: `capToolResultContent` keeps head+tail but the removed middle is what often contains the stack trace; the stop message for a budget breach is unaffected, but tool-error diagnostics become less actionable at tight caps.

## Future Considerations

- **Pre-flight budget admission control**: Price the prompt (`EstimateMessageTokens` + `summaryTokenBudget`) before calling `fallback.execute` and refuse to start the turn when the projection already exceeds `maxCost`/`maxTokens`, closing the one-turn overrun gap.
- **Evidence-preserving compaction**: Store compacted tool-result messages as a compressed attachment or hash-chained log alongside the summary item instead of dropping them, giving final-answer consumers a cryptographically auditable trail without repopulating the prompt.
- **Provenance edges**: Add `SourceRefs []{ItemIndex, ToolCallID, hash}` to `session.Item` and propagate through `SessionSummaryEvent` so downstream evaluators can distinguish “summarized from” vs “hallucinated.”
- **Structured cancellation contract**: Document the sequential stream→tool invariant and expose a single `RunStreamContext` that fan-out cancels both phases (e.g., a shared `toolCtx` tied to the stream goroutine) to justify “simultaneous” claims under parallel tool batches.
- **Compaction policy configurability for Aren**: Expose `keepBudget`/`summaryBudget` per agent so low-window local models used in Aren Phase 9 can retain more recent tool evidence verbatim.

## Questions / Gaps

- No evidence of a per-tool-call cost breakdown beyond the aggregate `Cost` on the assistant `chat.Message` (`studies/aren-go-runtime-study/sources/docker-agent/pkg/session/session.go:517-527`) — unclear how multi-tool batches attribute individual tool latency/cost to budgets.
- No formal “evidence-preservation” or “artifact” mechanism was found; `pkg/content/store.go` (`studies/aren-go-runtime-study/sources/docker-agent/pkg/content/store.go:46-97`) manages OCI registry artifacts for agent distribution, not conversation evidence. `Is there an equivalent to Aren's source-reference/artifact ledger?` — **No clear evidence found** after grepping `provenance`, `source_reference`, `artifact` (outside OCI content) across `pkg/` and reading `pkg/session/session.go` and `pkg/runtime/event.go`.
- Cancellation tests (`studies/aren-go-runtime-study/sources/docker-agent/pkg/runtime/turn_end_test.go:244-320`, `studies/aren-go-runtime-study/sources/docker-agent/pkg/runtime/toolexec/dispatcher_test.go:322-363`) drive single path cancels; no end-to-end test asserts that an in-flight `handleStream` and an in-flight tool batch are cancelled by the same `RunStream` `context.WithCancel` within the same iteration.
- Whether `sanitizeToolCalls`'s synthesized `No result provided` error messages are excluded from future `LoopDetector` signature hashing was not verified; current `callSignature` ignores synthetic inputs, but an attacker crafting matching synthetic failures could influence detector state.

---

Generated by `01.08-bounded-agent-loop-context-and-evidence` against `docker-agent`.
