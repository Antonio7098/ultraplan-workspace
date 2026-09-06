# Source Analysis: crush

## 01.08 Bounded Agent Loop, Context, and Evidence

### Source Info

| Field | Value |
|-------|-------|
| Name | crush |
| Path | `studies/aren-go-runtime-study/sources/crush` |
| Language / Stack | Go (module `github.com/charmbracelet/crush`, LLM via `charm.land/fantasy`, SQLite via `sqlc`, Bubble Tea TUI, pubsub broker) |
| Analyzed | 2026-08-29 |

## Summary

Crush implements a bounded Go agent loop in `internal/agent/agent.go:566` around `fantasy.NewAgent(...).Stream` (`internal/agent/agent.go:685`, `internal/agent/agent.go:796`). The loop is multi-turn with provider-driven `OnToolCall`/`OnToolResult` callbacks persisting a `message.Tool`/`message.Assistant` transcript to SQLite (`internal/message/message.go:163`, `internal/message/message.go:218`). Bounding is heuristic: context-window `StopWhen` (`internal/agent/agent.go:1037-1062`) plus repeated tool-call loop detection (`internal/agent/loop_detection.go:19`) triggers summarization instead of a hard iteration cap. Context compaction is automatic summarization via `Summarize` (`internal/agent/agent.go:1329`) that creates an `IsSummaryMessage` (`internal/message/message.go:174`) and truncates history by `SummaryMessageID` in `getSessionMessages` (`internal/agent/agent.go:1689-1709`). Transcript invariants are actively repaired in `preparePrompt` (`internal/agent/agent.go:1527`) via `filterOrphanedToolResults` (`internal/agent/agent.go:1626`) and `syntheticToolResultsForOrphanedCalls` (`internal/agent/agent.go:1662`). Budgets are post-hoc: token/cost accumulation in `updateSessionUsage` (`internal/agent/agent.go:1906`) and `updateSessionTokenCounters` (`internal/agent/agent.go:1939`) with `fallbackStepUsage` estimation (`internal/agent/usage_fallback.go:18`), no pre-flight enforcement except `MaxOutputTokens` passthrough (`internal/agent/agent.go:791-803`). Cancellation is the strongest area: per-session `activeCancel` map (`internal/agent/agent.go:185`), `dispatchMu` serialization (`internal/agent/agent.go:191`), and monotonic `acceptSeqGen`/`cancelMark` high-water mark (`internal/agent/agent.go:200-206`, `internal/agent/agent.go:493`) covering accepted-but-not-yet-active runs, with `genCtx = context.WithCancel` (`internal/agent/agent.go:644`) cancelled via `Cancel` (`internal/agent/agent.go:1955`) and bounded-blocking terminal delivery `PublishMustDeliver` (`internal/pubsub/broker.go:201`, `internal/agent/agent.go:539`). Evidence preservation is DB-backed messages plus `filetracker.Service` (`internal/filetracker/service.go:16`) but no first-class provenance/artifact manifest.

## Rating

**5/10** — Loop termination relies on model coop + two `StopWhen` conditions, not a hard `MaxSteps`/`MaxToolCalls` guard; fantasy is delegated for iteration bounding with no Crush-side constant. Context compaction preserves invariants (orphan filtering/synthetic results) but compaction is summarize-only (no token-budget truncation) and `disableAutoSummarize` can disable it entirely. Budgets are reported not enforced. Cancellation is well-engineered (accepted-run sequence mark, `CompareAndDelete` cleanup, `WithoutCancel` writes) and well-tested.

## Evidence Collected

Every entry MUST include a file path with line numbers. Format: `path/to/file.ts:NN`.

| Area | Evidence | File:Line |
|------|----------|-----------|
| Loop entry | `func (a *sessionAgent) Run(ctx context.Context, call SessionAgentCall) (result *fantasy.AgentResult, retErr error)` — validates, handles cancel-on-entry, busy-queue, creates `genCtx` | `internal/agent/agent.go:566` |
| Loop core | `fantasy.NewAgent(largeModel.Model, fantasy.WithSystemPrompt(...), fantasy.WithTools(...))` then `agent.Stream(genCtx, fantasy.AgentStreamCall{...})` | `internal/agent/agent.go:685-796` |
| Iteration/turn plumbing | `PrepareStep` drains queue, creates assistant message, sets `currentAssistant`; `OnToolCall`/`OnToolResult` persist tool turns; `OnStepFinish` updates `FinishReason` and `updateSessionUsage` | `internal/agent/agent.go:808-1036` |
| Output/token limit | `MaxOutputTokens *int64` passthrough only if `>0`; no hard cap else delegated to provider/fantasy | `internal/agent/agent.go:791-803` |
| Tool-call limit | No `MaxToolCalls` constant found; only loop-detection heuristic | `internal/agent/loop_detection.go:11-14`, `internal/agent/agent.go:1059` |
| Loop detection | `hasRepeatedToolCalls(steps, loopDetectionWindowSize=10, loopDetectionMaxRepeats=5)` via SHA256 of `toolName+input+output` | `internal/agent/loop_detection.go:11-39` |
| Loop detection StopWhen | `StopWhen: []fantasy.StopCondition{ context-window Check, hasRepeatedToolCalls }` | `internal/agent/agent.go:1037-1062` |
| Context-window threshold | `largeContextWindowThreshold=200_000`, `largeContextWindowBuffer=20_000`, `smallContextWindowRatio=0.2`; computes `remaining=cw-tokens` and sets `shouldSummarize=true` | `internal/agent/agent.go:53-60`, `internal/agent/agent.go:1038-1056` |
| Summarization trigger | `if shouldSummarize { a.activeRequests.Del; a.Summarize(...); if tool calls present re-queue prompt }` | `internal/agent/agent.go:1192-1207` |
| Compaction impl | `Summarize` creates `IsSummaryMessage=true` assistant, streams summary agent, saves `currentSession.SummaryMessageID=summaryMessage.ID`, resets `CompletionTokens, PromptTokens, Cost` | `internal/agent/agent.go:1369-1459` |
| Compaction read | `getSessionMessages` slices `msgs[summaryMsgIndex:]` and mutates `msgs[0].Role=User` so summary becomes context prefix | `internal/agent/agent.go:1689-1709` |
| Transcript repair | `preparePrompt` builds `knownToolCallIDs`/`knownToolResultIDs`, calls `filterOrphanedToolResults` and `syntheticToolResultsForOrphanedCalls` | `internal/agent/agent.go:1539-1587` |
| Orphan filter | Drops `ToolResultPart` where `ToolCallID` not in `knownToolCallIDs`, logs warn | `internal/agent/agent.go:1626-1653` |
| Synthetic results | Injects `ToolResultPart{Output: ToolResultOutputContentError{Error: "tool call was interrupted..."}}` per orphaned `ToolCall` | `internal/agent/agent.go:1662-1687` |
| Context limit enforcement point | Enforced *during* streaming via `StopWhen`, not before `Stream`; zero-context models skip `if cw==0 return false` | `internal/agent/agent.go:1042-1044` |
| Token/cost accounting | `updateSessionUsage` computes `cost = CostPer1M*Tokens/1e6`, override via `openrouterCost`, `extractHyperCredits`; `EstimatedUsage` flag | `internal/agent/agent.go:1906-1946` |
| Token fallback | `fallbackStepUsage` estimates `inputTokens=estimateMessageTokens(messages)`, `outputTokens=estimateStepCompletionTokens(step)` when `step.Usage` zero, `approxTokenCount=(len+3)/4` | `internal/agent/usage_fallback.go:18-176` |
| Session token storage | `Session{PromptTokens, CompletionTokens, Cost, SummaryMessageID, EstimatedUsage}` persisted via `UpdateSession` | `internal/session/session.go:50-63`, `internal/session/session.go:191-221` |
| Evidence: message store | `message.Service.Create/Update/Get/List` backed by `db.Queries`; `Update` debounced (`33ms`) but terminal updates `shouldFlushNow` flush sync + `PublishMustDeliver` | `internal/message/message.go:31-49`, `internal/message/message.go:218-274`, `internal/message/message.go:379-383` |
| Evidence: file tracking | `filetracker.Service{RecordRead, LastReadTime, ListReadFiles}` records via `RecordFileRead`; called from `edit.go:182`, `lsp_rename.go:98` etc. | `internal/filetracker/service.go:16-27`, `internal/agent/tools/edit.go:182` |
| Evidence: tool result pairing | `ToolResult{ToolCallID, Name, Content, Data, MIMEType, IsError}` stored as `ContentPart`; `IsFinished`/`FinishReason` on `Finish` part | `internal/message/content.go:117-127`, `internal/message/content.go:129-135` |
| Cancellation registration | `genCtx, cancel = context.WithCancel(runCtx); ac:=&activeCancel{cancel}; a.activeRequests.Set(call.SessionID, ac); defer CompareAndDelete` | `internal/agent/agent.go:643-657` |
| Cancel dispatch | `Cancel(sessionID)` takes `sessionMu` lock, `ac.cancel()` on `activeRequests`, cancels `"-summarize"` too, sets `cancelMark = max(existing, acceptSeqGen)` only if `acceptedRuns>0`, `clearQueueAndNotify` | `internal/agent/agent.go:1955-2009` |
| Cancel-on-entry | `if call.Accepted != nil && a.canceledBySeq(...) { accept.Close(); Unlock; persistCanceledTurn; publishRunComplete(Cancelled:true); return }` | `internal/agent/agent.go:591-618` |
| Pending cancel semantics | `canceledBySeq(seq)` → `mark==0→false, seq==0→covered if mark>0, else seq <= mark`; `clearPendingCancel` holds dispatch lock | `internal/agent/agent.go:493-500`, `internal/agent/agent.go:476-481` |
| Interruption persistence | `persistCanceledTurn` uses `context.WithoutCancel(ctx)` + `5s` timeout, writes `FinishReasonCanceled` | `internal/agent/agent.go:508-528` |
| Error-path cancel handling | On `Stream` error with `isCancelErr` and `currentAssistant==nil`, calls `persistCanceledTurn`; else appends synthetic `ToolResult` for unfinished `ToolCalls`, `AddFinish(Canceled)` via `cleanupCtx=WithoutCancel` | `internal/agent/agent.go:1067-1189` |
| Pubsub terminal guarantee | `Broker.Publish` lossy (buffer 4096, `DropCount`), `PublishMustDeliver` bounded-blocking `50ms` per subscriber then `MustDeliverDropCount` | `internal/pubsub/broker.go:6-46`, `internal/pubsub/broker.go:165-236` |
| RunComplete terminal event | Deferred `publishRunComplete` echoes `MessageID/Text`, `Cancelled = errors.Is(retErr,Canceled) || ctx.Err()!=nil`, honors `OnComplete` hook else broker, uses `FlushAll` then ordering | `internal/agent/agent.go:749-781` |
| StopTurn short-circuit | `if finishReason==ToolUse { for tr := range stepResult.Content.ToolResults() if tr.StopTurn { finishReason=EndTurn } }` and `ToolResponse.StopTurn` set by permission/hook | `internal/agent/agent.go:1007-1018`, `internal/agent/tools/tools.go:66-70`, `internal/agent/hooked_tool.go:62-71` |
| Hook interruption | `hookedTool.Run` -> `runner.Run(EventPreToolUse)`; `DecisionDeny` or `Halt` returns `NewTextErrorResponse` with `StopTurn=Halt` | `internal/agent/hooked_tool.go:54-73` |
| Config escape hatch | `DisableAutoSummarize bool` in `config.Config.Options`; also `CW==0` disables summarization | `internal/agent/agent.go:179-231`, `internal/config/config.go:321`, `internal/agent/agent.go:1042` |

## Answers to Dimension Questions

**What guarantees that a loop terminates even when the model keeps requesting tools?**

Two `StopWhen` guards in `internal/agent/agent.go:1037-1062` are the only Crush-owned termination signals:

1. Context-window guard (`internal/agent/agent.go:1038-1056`) returns `true` when `remaining <= threshold` ( `20k` for `cw>200k` else `20%`), setting `shouldSummarize=true`. This does not directly stop tool recursion; it ends the fantasy `Stream` and then `Summarize` is invoked (`internal/agent/agent.go:1192-1207`). If the assistant still has tool calls after summarization, the original user prompt is re-queued (`internal/agent/agent.go:1198-1205`), so termination is not guaranteed in a single turn — it restarts.

2. Loop-detection guard (`internal/agent/agent.go:1059`, `internal/agent/loop_detection.go:19`) returns true when `>5` identical `toolName+input+output` signatures appear in the last `10` steps (`internal/agent/loop_detection.go:11-14`). Signature computed via `getToolInteractionSignature` (`internal/agent/loop_detection.go:45`).

Neither is a hard `MaxSteps` or `MaxToolCalls`/`MaxIterations` bound. `MaxOutputTokens` (`internal/agent/agent.go:791`) is passed to the provider only and not enforced by Crush. Actual iteration bounding is delegated to `charm.land/fantasy` defaults (not inspected within isolation boundary; no `fantasy.WithMaxSteps` call found in `internal/agent/agent.go:685`). A model that keeps producing *distinct* tool calls with distinct outputs will not trigger loop detection and will only stop when the context window is exhausted or the provider itself enforces `finish_reason=length` (`internal/agent/agent.go:989` maps to `FinishReasonMaxTokens`). This is a best-effort, not a hard guarantee.

**Can compaction orphan a tool call or remove the evidence supporting the final result?**

Compaction is summarization, not truncation-by-token. `Summarize` (`internal/agent/agent.go:1329`) creates a persistent `IsSummaryMessage` (`internal/agent/agent.go:1373`) and sets `session.SummaryMessageID` (`internal/agent/agent.go:1455`), after which `getSessionMessages` (`internal/agent/agent.go:1695-1707`) slices history to `msgs[summaryMsgIndex:]` and rewrites `msgs[0].Role = User`. Evidence *before* the summary pointer is dropped from future `history` (`internal/agent/agent.go:783`, `internal/agent/agent.go:1351`) — tool calls/results, file reads, and costs from earlier turns are no longer in context.

Crush explicitly guards against orphaning *within* the retained window: `preparePrompt` (`internal/agent/agent.go:1527-1587`) builds `knownToolCallIDs`/`knownToolResultIDs` across all retained `msgs` and then (a) drops orphaned tool results (`internal/agent/agent.go:1626`) that would brick the provider API, and (b) injects synthetic error results for orphaned calls (`internal/agent/agent.go:1662`) so every `tool_use` is paired. `preparePrompt` also skips empty assistant messages (`internal/agent/agent.go:1562`) and respects `supportsImages` filtering (`internal/agent/agent.go:1573`). Tests `TestPreparePrompt_OrphanedToolUse` (`internal/agent/agent_test.go:757`) and `TestPreparePrompt_OrphanedToolUseMixed` (`internal/agent/agent_test.go:824`) verify this.

However, *across* summarization boundary, original tool outputs are lost: only the natural-language summary remains. `filetracker` state (`internal/filetracker/service.go:16`) is not compacted or referenced in prompts. There is no artifact manifest linking `SummaryMessageID` to the discarded `ToolResult` IDs, and `session.Cost/Tokens` are reset (`internal/agent/agent.go:1456-1458`) without retaining per-turn cost ledger beyond the overwritten counters. If `DisableAutoSummarize=true` (`internal/config/config.go:321`), no compaction occurs and context will eventually exceed `cw`, likely causing provider `FinishReasonLength` errors rather than clean compaction.

**Which budgets are enforced before work begins and which are only reported afterward?**

*Enforced before work*: effectively none. `ValidateCall` (`internal/agent/agent.go:556`) checks only `Prompt != ""` or text attachment and `SessionID != ""`. `MaxOutputTokens` is set per-call (`internal/agent/agent.go:791-794`, `internal/agent/coordinator.go:254-257` from `DefaultMaxTokens` or `ModelCfg.MaxTokens`) and sent as `MaxOutputTokens: &int64` to `AgentStreamCall` (`internal/agent/agent.go:802`), but enforcement is provider-side. Context-window check is evaluated *during* `Stream` in `StopWhen` (`internal/agent/agent.go:1045`), not as a pre-flight reject.

*Reported afterward*: all. `fallbackStepUsage` (`internal/agent/usage_fallback.go:18`) derives `usage` from `step.Usage` or estimated via `approxTokenCount` (`internal/agent/usage_fallback.go:171`). `updateSessionUsage` (`internal/agent/agent.go:1906`) computes cost via `CatwalkCfg.CostPer1M*` (`internal/agent/agent.go:1912-1915`), overridden by `openrouterCost` (`internal/agent/agent.go:1874`) or zeroed if `estimated` or `FlatRate`. `updateSessionTokenCounters` (`internal/agent/agent.go:1939`) overwrites `session.PromptTokens/CompletionTokens` (not cumulative except `Cost+=`). Summary cost uses `resp.TotalUsage` (`internal/agent/agent.go:1451`). Title generation cost is accumulated via `UpdateTitleAndUsage` (`internal/agent/agent.go:1866`). No token/time/cost/tool-call budget rejects a turn before `fantasy.Stream` starts, and no per-turn `Cost` ceiling or wall-clock timeout exists in `sessionAgent.Run`; the only time budget is the `5s` bounded `FlushAll`/`publishRunComplete` (`internal/agent/agent.go:753`, `internal/agent/agent.go:446`).

**Does cancellation reach an active model stream and active tool at the same time?**

Yes via a single `genCtx`. `Run` creates `genCtx, cancel = context.WithCancel(runCtx)` where `runCtx` carries `SessionIDContextKey` (`internal/agent/agent.go:643-646`) and registers `ac` in `activeRequests` (`internal/agent/agent.go:646`) under `dispatchMu`. `Cancel(sessionID)` (`internal/agent/agent.go:1955`) holds `sessionMu` and calls `ac.cancel()` on the stored func (`internal/agent/agent.go:1971`). Because `fantasy.Stream(genCtx, ...)` (`internal/agent/agent.go:796`) is the callee, cancellation propagates simultaneously to (a) the provider stream (fantasy monitors `ctx.Done()`), and (b) tool execution — tools are executed by fantasy with `callContext` derived from `genCtx` (the `PrepareStep` callback receives `callContext` (`internal/agent/agent.go:808`) that inherits `genCtx`; tool `Run` receives that `ctx` (`internal/agent/tools/tools.go:14` defines context keys)). The `OnRetry` path explicitly logs provider retry (`internal/agent/agent.go:931`) and resets streamed content under `genCtx`. 

Persistence after cancel uses `context.WithoutCancel(ctx)` + `5s` timeout so writes are not lost when `genCtx` is canceled: `persistCanceledTurn` (`internal/agent/agent.go:509`), error-path synthetic tool results (`internal/agent/agent.go:1097`), `FlushAll` (`internal/agent/agent.go:753`), and `publishRunComplete` (`internal/agent/agent.go:779`) all detach. The per-session `dispatchMu` + `acceptedRuns`/`cancelMark` sequence ensures a cancel racing `Run` is not lost: `cancelMark` is set under `acceptedMu` only when `acceptedRuns>0` (`internal/agent/agent.go:1999`), and `Run` checks `canceledBySeq` under `sessionMu` before `IsSessionBusy` (`internal/agent/agent.go:591`), draining queued prompts atomically (`internal/agent/agent.go:399-424`). Active tool goroutines that ignore `ctx` will run to completion but their results are still written via `ctx` (parent context, `internal/agent/agent.go:967-981` comment "Use parent ctx instead of genCtx").

## Architectural Decisions

- **Fantasy as loop driver.** Crush delegates iteration to `charm.land/fantasy` (`internal/agent/agent.go:685`, `internal/agent/agent.go:796`), adding only two `StopWhen` conditions. This keeps provider differences (Anthropic/OpenAI/Google/Bedrock etc. via `internal/agent/coordinator.go:1099-1162`) behind one interface but leaves hard iteration bounding to the external library.

- **Summarize-as-compaction.** Instead of token-budget sliding window, Crush uses a secondary LLM summary (`internal/agent/agent.go:1364-1416` with `templates/summary.md:0`) stored as `IsSummaryMessage` (`internal/message/message.go:174`) and referenced by `SummaryMessageID` (`internal/session/session.go:58`). `getSessionMessages` (`internal/agent/agent.go:1689`) treats the summary as new `User` context prefix. This preserves natural language evidence but discards structured tool outputs.

- **Transcript repair not strict validation.** `preparePrompt` (`internal/agent/agent.go:1527`) repairs orphaned pairs rather than rejecting the turn. This avoids permanently locked sessions after interruption (`internal/agent/agent.go:1107-1149` also injects synthetic error results on cancel), at cost of injecting synthetic data into history.

- **Debounced persistence with terminal must-deliver.** `message.Service.Update` (`internal/message/message.go:218`) coalesces streaming deltas (`33ms`, `internal/message/message.go:19`) but `shouldFlushNow` (`internal/message/message.go:419`) flushes sync for `IsFinished`/tool-call changes, and `flushOne` (`internal/message/message.go:352-383`) uses `PublishMustDeliver` (`internal/pubsub/broker.go:201`) for terminal events while `Publish` (`internal/pubsub/broker.go:165`) remains lossy for intermediate deltas (buffer `4096`, `internal/pubsub/broker.go:40`).

- **Monotonic accept sequence for cancellation.** `BeginAccepted` increments `acceptSeqGen` under `acceptedMu` (`internal/agent/agent.go:305-311`), `Cancel` sets `cancelMark = max(existing, mark)` (`internal/agent/agent.go:2002`), and `canceledBySeq` compares `seq <= mark` (`internal/agent/agent.go:493`). This prevents idle `Escape` poisoning (`internal/agent/agent.go:204-206`) and ensures prompts accepted after a cancel are not poisoned (`internal/agent/dispatch_cancel_test.go:370-447`).

## Notable Patterns

- **Compare-and-delete active guard.** `activeRequests` entries are removed via `CompareAndDelete` (`internal/agent/agent.go:657`) so a deferred cleanup does not delete a newer run's entry registered in the window between explicit `Del` and return.

- **Quiet-period queue drain.** `drainQueueForStep` (`internal/agent/agent.go:399`) under `dispatchMu` partitions queued prompts into `fold` (no RunID) vs `keep` (RunID) vs `canceledWithRunID`, with `publishCanceledQueueDrops` (`internal/agent/agent.go:435`) emitting `Cancelled:true` `RunComplete` for RunID-bearing drops so `crush run` (`internal/cmd/run.go:0`) does not hang.

- **StopTurn protocol.** Permission denial (`internal/agent/tools/tools.go:66`) and hook halt (`internal/agent/hooked_tool.go:70`) set `ToolResponse.StopTurn=true`; `OnStepFinish` (`internal/agent/agent.go:1011-1017`) converts `FinishReasonToolUse` → `EndTurn`, cleanly terminating the turn without another model call.

- **Media workaround for provider limits.** `workaroundProviderMediaLimitations` (`internal/agent/agent.go:2154`) splits media `ToolResultPart` into text placeholder + synthetic `User` `FilePart` for non-Anthropic/Bedrock providers, or replaces with placeholder for text-only models (`internal/agent/agent.go:2184-2196`), verified by `TestWorkaroundProviderMediaLimitations_*` (`internal/agent/agent_test.go:896-1034`).

- **In-flight auth coalescing.** `coordinator.run` (`internal/agent/coordinator.go:281-337`) buffers `OnComplete` per attempt and publishes exactly one `RunComplete` via `PublishMustDeliver` after unauthorized→re-auth→retry, with `MarkRunCompletePublished` to suppress duplicate from backend.

## Tradeoffs

- **Heuristic vs hard bound:** Loop detection via SHA256 signature is content-aware and avoids false termination on alternating tools (`internal/agent/loop_detection_test.go:128`), but an attacker or errant model producing unique tool inputs/outputs each step evades it indefinitely. A hard `MaxSteps` (e.g., 50) would deterministically bound cost/latency but might cut legitimate long tasks.

- **Summarization fidelity vs determinism:** LLM-generated summaries preserve intent and todos (`internal/agent/agent.go:2241-2253`) better than blind truncation, but are non-deterministic, incur extra cost/latency, and may hallucinate. The summary completion tokens approximation (`internal/agent/agent.go:1948`) falls back to `approxTokenCount` when provider usage is zero.

- **Lossy streaming vs guaranteed terminal:** `Publish` dropping under `bufferSize=4096` contention (`internal/pubsub/broker.go:40`) keeps TUI responsive during fast token streams, with `PublishMustDeliver` (`50ms` timeout, `internal/pubsub/broker.go:45`) for terminal events. Tradeoff: slow subscribers still drop terminal events after `MustDeliverDropCount` (`internal/pubsub/broker.go:156`), requiring `Re-fetch on next session-visible event` recovery (comment `internal/pubsub/broker.go:198-200`).

- **WithoutCancel writes vs leak:** Using `context.WithoutCancel` for final writes (`internal/agent/agent.go:509`, `internal/agent/agent.go:1097`) ensures evidence is not lost on cancellation/workspace shutdown, but a hung DB write can outlive the session `5s` timeout and delay `CancelAll`'s `5s` spin (`internal/agent/agent.go:2026-2034`).

- **Single model for summarize:** Summarization reuses `largeModel` (`internal/agent/agent.go:1335`) even for long contexts where a smaller/cheaper model might suffice, simplifying auth/token handling but increasing cost.

## Failure Modes / Edge Cases

- **Unbounded distinct tool loops:** With `windowSize=10`/`maxRepeats=5`, 100 steps of `read("a.txt")` with varying `output` each time produce distinct signatures and never trigger `StopWhen`. If `cw==0` (custom/local model, `internal/agent/agent.go:1042`) auto-summarize is disabled entirely, so even context exhaustion does not stop — provider may error with `FinishReasonLength` but Crush will persist it as `FinishReasonMaxTokens` (`internal/agent/agent.go:989-990`) and not retry.

- **DisableAutoSummarize deadlock:** `DisableAutoSummarize=true` (`internal/config/config.go:321`) or `smallContextWindowRatio` path for small windows may let context exceed provider limit, causing repeated `FinishReasonLength` without recovery; no pre-flight token check rejects the turn.

- **SummaryMessageID dangling:** If `Summarize` fails after creating `summaryMessage` but before `sessions.Save` (`internal/agent/agent.go:1459`), `SummaryMessageID` is not updated; subsequent `getSessionMessages` retains full history and may OOM. Error path marks `FinishReasonError` (`internal/agent/agent.go:1425`) but does not delete the partial summary (unlike cancel path `internal/agent/agent.go:1420`).

- **Orphan synthesis poisoning:** Synthetic error result `"tool call was interrupted..."` (`internal/agent/agent.go:1676`) becomes part of persistent history; a later model may misinterpret it as authoritative failure and spin. No TTL or marker distinguishes synthetic vs real results beyond content.

- **Filetracker not checkpointed:** `ListReadFiles` (`internal/filetracker/service.go:77`) is not consulted in `preparePrompt` or summarization, so compaction loses which files grounded the answer; reproduction from summary alone is impossible.

- **Cost non-cumulative overwritten tokens:** `updateSessionTokenCounters` (`internal/agent/agent.go:1939`) overwrites `PromptTokens`/`CompletionTokens` with last usage, not sum; `session.Cost` is cumulative but `PromptTokens` is not, so total prompt tokens over session lifetime cannot be reconstructed post-summarization (`internal/agent/agent.go:1456-1457` resets to summary-only tokens).

- **Cancel window between dequeue and register:** Closed by `BeginAccepted` before `mu.Unlock` (`internal/agent/agent.go:1313`), but if `CancelAll` iterates `activeRequests` without holding `dispatchMu` (`internal/agent/agent.go:2022`), a freshly dequeued recursive `Run` that has called `BeginAccepted` but not yet `activeRequests.Set` will not be canceled — it will take `cancel-on-entry` on next `Run` dispatch instead, preserving correctness but delaying cancellation by one turn.

- **Hook double-fire risk:** `PreToolUse` hooks fire via `hookedTool` (`internal/agent/hooked_tool.go:56`) which runs `hooks.Runner` (parallel shell commands, `internal/hooks/runner.go:0`) without timeout propagation to the tool itself; a hung hook blocks the tool `Run` under `genCtx` but `Cancel` will still cancel `genCtx`, leaving hook subprocesses potentially orphaned (process supervision is separate dimension).

## Future Considerations

- Add hard `MaxSteps`/`MaxToolCalls`/`wall-clock` guards in `StopWhen` (e.g., `maxSteps=80`, `maxToolCalls=100`) alongside heuristic loop detection, with configurable `agent.max_steps` in `internal/config/config.go:321` and pre-flight rejection when `ApproxTokenCount(history) > cw - buffer`.

- Persist compaction evidence: store discarded `ToolResult` IDs and `filetracker` snapshot as `Session.SummaryMetadata` alongside `SummaryMessageID`, and include `Cost`/`Tokens` ledger (not just last counters) for audit.

- Enforce budgets before `Stream`: validate `PromptTokens + estimate(prompt+attachments) < cw` and `Cost + projected < budget` before `agent.Stream`, returning `FinishReasonError` rather than starting a doomed turn.

- Introduce time/cost tool-call budgets: `fantasy.AgentStreamCall` could carry `MaxToolCalls`/`Timeout` propagated to `OnToolCall` context with `context.WithTimeout` per tool, so cancellation reaches tools independently of model stream.

- Expose compaction selection policy: allow caller to choose `summarize` vs deterministic `truncate-oldest` vs `keep-tool-results` modes; current `preparePrompt` filtering is fixed and summarization prompt (`internal/agent/templates/summary.md:0`) is embedded, not configurable.

## Questions / Gaps

- What is `charm.land/fantasy`'s default `MaxSteps` and does it enforce a hard iteration cap when Crush supplies none? No `MaxSteps` wiring found in `internal/agent/agent.go:685` or `internal/agent/agent.go:1364`; fantasy vendor code was out of isolation boundary and not inspected.

- Is provider `Usage` reliable for all providers, and does `fallbackStepUsage` (`internal/agent/usage_fallback.go:18`) ever double-count `ReasoningTokens` vs `OutputTokens`? No test asserts fallback against real provider `Usage` shapes.

- How are tool outputs bounded? File reads go through `filetracker` but `view`/`edit` tool result sizes are not capped in `preparePrompt`; large tool results could exceed context without summarization.

- What is the provenance contract for `FilePart` media in `ToolResult`? `workaroundProviderMediaLimitations` decodes base64 (`internal/agent/agent.go:2199`) but does not persist the original `media.Data` hash, so integrity cannot be verified post-compaction.

- Does `CancelAll`'s `5s` spin + `200ms` poll (`internal/agent/agent.go:2026-2034`) guarantee all `activeRequests` goroutines have reached `FlushAll` before DB close during workspace shutdown? No shutdown test was found.

---

Generated by `01.08 Bounded Agent Loop, Context, and Evidence` against `crush`.
