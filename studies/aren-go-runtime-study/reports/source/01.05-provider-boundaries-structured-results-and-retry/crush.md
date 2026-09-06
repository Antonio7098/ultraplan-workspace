# Source Analysis: crush

## 01.05 Provider Boundaries, Structured Results, and Retry

### Source Info

| Field | Value |
|-------|-------|
| Name | crush |
| Path | `studies/aren-go-runtime-study/sources/crush` |
| Language / Stack | Go 1.26.6 / charm.land/fantasy v0.41.3, charm.land/catwalk, SQLite/sqlc, Bubble Tea v2 |
| Analyzed | 2026-08-29 |

## Summary

Crush does not own its provider boundary. All live-model semantics are delegated to the external library `charm.land/fantasy` (`go.mod:10`). `Coordinator` (`internal/agent/coordinator.go:114`) resolves a `catwalk.Model`+`ProviderConfig` into a `fantasy.Provider`/`fantasy.LanguageModel` via `buildProvider`/`buildAgentModels`, then `sessionAgent` (`internal/agent/agent.go:685`) drives the run as a single `fantasy.NewAgent(...).Stream(AgentStreamCall{...})` with normalized callbacks (`OnTextDelta`, `OnReasoningDelta`, `OnToolCall`, `OnStepFinish`, `OnRetry`, `OnAuthRefresh`). Request translation (temperature/top-p/thinking/reasoning_effort, cache-control, extra_body) is merged in `getProviderOptions`/`mergeCallOptions` with per-provider `switch` branches. Transport/SSE decoding, `Usage` extraction, and the `ProviderError`/`RetryError` taxonomy live entirely in fantasy. Retry is split: fantasy owns the attempt budget (3 retries, 5s initial, 2x backoff, `Retry-After` aware), crush only contributes `OnRetry` content-reset and `OnAuthRefresh` credential repair (OAuth/AWS SSO/APIKey template) which restarts the budget once. Structured-output (`fantasy.ObjectCall`/`GenerateObject`) is implemented in fantasy (`object.go`) but never used on crush's production path; crush has no schema conversion, validation, or repair prompts — its only schema use is `crush schema` for config JSON Schema (`internal/cmd/schema.go:14`). Late-stream failure preserves partial text under `FinishReasonError`; usage falls back to `approxTokenCount` estimation when providers return zero.

## Rating

**5/10** — Provider-neutral streaming and a layered retry taxonomy exist, but they are outsourced to `fantasy` and the `coordinator` still leaks vendor detail through a large provider-type switch and SDK-specific option parsing. Backoff has no jitter/budget caps and no per-request deadline; rate-limit handling is header-respecting but delegated. Structured-output hardening is absent in crush itself. Usage accounting is the strongest area (fallback estimation + openrouter/hyper cost branches with `EstimatedUsage` flag).

Rationale: Aren Phases 3-6 want a Go boundary that isolates lifecycle from provider SDKs, enforces typed failures on 200, and owns explicit retry/validation budgets. Crush achieves isolation via `fantasy` but reintroduces coupling in `getProviderOptions`; it correctly distinguishes transport vs. auth retry vs. cancellation but leaves structured validation unimplemented and relies on fantasy defaults for timing policy.

## Evidence Collected

Every entry MUST include a file path with line numbers. Format: `path/to/file.ts:NN`.

| Area | Evidence | File:Line |
|------|----------|-----------|
| Provider interface ownership | `fantasy.LanguageModel`/`fantasy.Provider`/`fantasy.Agent` consumed; no crush-defined interface | `internal/agent/agent.go:154-158` |
| Producer-supplied | `Model.Model fantasy.LanguageModel` wraps fantasy type | `internal/agent/agent.go:154` |
| Request entrypoint | `fantasy.NewAgent(largeModel.Model, WithSystemPrompt, WithTools(...)).Stream(AgentStreamCall{...})` | `internal/agent/agent.go:685-796` |
| Provider selection | `buildProvider` switches on `providerCfg.Type` (openai/anthropic/openrouter/vercel/azure/bedrock/google/openaicompat/hyper/custom) | `internal/agent/coordinator.go:1099-1162` |
| Per-provider constructors | `buildAnthropicProvider`, `buildOpenaiProvider`, `buildOpenrouterProvider`, `buildOpenaiCompatProvider`, `buildBedrockProvider`, `buildGoogleProvider` | `internal/agent/coordinator.go:891-1090` |
| Request merging | `mergeCallOptions` + `getProviderOptions` merges catwalk/provider/model JSON via `jsons.Merge` | `internal/agent/coordinator.go:360-622` |
| Reasoning effort policy | `effectiveReasoningEffort` prefers modelCfg → catwalk default → first level | `internal/agent/coordinator.go:342-358` |
| Provider-specific extra_body | `openai`/`anthropic`/`google`/`openaicompat` branches inject `reasoning_effort`/`thinking`/`extra_body` per vendor ID | `internal/agent/coordinator.go:414-608` |
| Stream normalization (consumer) | `OnTextDelta`/`OnReasoningDelta`/`OnReasoningEnd`/`OnToolCall`/`OnToolResult`/`OnStepFinish`/`StopWhen` wired on `AgentStreamCall` | `internal/agent/agent.go:880-1062` |
| Neutral error vocabulary | Handles `*fantasy.ProviderError`, `*fantasy.Error`, `fantasy.IsTransportError`, `fantasy.RetryError`, plus `context.Canceled` | `internal/agent/agent.go:1155-1182` |
| Cancel transport error title | `NewTransportError`/`IsTransportError` with `http2TransportErrorFragments` and `TransientStreamErrorTypes` | `fantasy/errors.go:189-270` (vendored `charm.land/fantasy@v0.41.3`) |
| 200-OK typed failure | `TransientStreamErrorTypes = {server_error, overloaded_error, api_error...}` → `TransientError=true` → retryable despite 200 | `fantasy/errors.go:257-268` |
| Content-filter typed failure | `FinishReasonContentFilter` mapped to `message.FinishReasonContentFilter` with UI banner | `internal/agent/agent.go:995-1005` |
| Retry budget (fantasy) | `DefaultRetryOptions{MaxRetries:3, InitialDelay:5s, BackoffFactor:2.0}` + `getRetryDelayInMs` parses `retry-after-ms`/`retry-after` | `fantasy/retry.go:38-105` (vendored) |
| Retry execution | Recursive `retryWithExponentialBackoff` with `MaxRetries`, `BackoffFactor`, `OnRetry` hook, `isRetryableError`, `isAuthError` | `fantasy/retry.go:106-205` |
| Crush OnRetry hook | Resets streamed content so retry doesn't concatenate: `currentAssistant.ResetStreamedContent()` | `internal/agent/agent.go:931-941` |
| Auth retry (OAuth/template/AWS) | `makeAuthRefreshCallback` → `retryAfterUnauthorized` handles `OAuthToken`/`APIKeyTemplate`/`AWSAuthRefresh` | `internal/agent/coordinator.go:1349-1358` |
| OAuth refresh + 1-shot wait | `refreshOAuth2Token` then `waitForInteractiveReauth` (5m detached, `WaitForTokenChange`, rebuild models) | `internal/agent/coordinator.go:1274-1338` |
| AWS SSO refresh | `refreshAWSCredentials` publishes `TypeAWSSSOAuth`, runs `sh -c $AWSAuthRefresh` detached 5m, rebuilds models | `internal/agent/aws_sso_refresh.go:45-90` |
| Structured output (unused) | `ObjectCall`/`ObjectResponse`/`ObjectMode{Auto,JSON,Tool,Text}` defined in fantasy but not called in crush prod; only test fakes implement `GenerateObject` | `fantasy/object.go:14-70` + `internal/agent/dispatch_cancel_test.go:5-9` |
| Schema use (config only) | `invopop/jsonschema` reflects `config.Config` for `crush schema` | `internal/cmd/schema.go:8-22` |
| Tool-input validation repair | `sanitizeToolInput` replaces invalid JSON with `{}` and marks `sanitizedToolCalls` → synthetic error tool result | `internal/agent/agent.go:950-973`, `internal/agent/agent.go:2272-2287` |
| Orphan repair | `filterOrphanedToolResults` drops tool_result without call; `syntheticToolResultsForOrphanedCalls` injects synthetic error | `internal/agent/agent.go:1626-1687` |
| Rate-limit delay | `getRetryDelayInMs` prefers `retry-after-ms` then `retry-after` (seconds or RFC1123), clamped `<60s` or `<exponentialDelay` | `fantasy/retry.go:18-55` |
| Usage fallback | `fallbackStepUsage` → `estimateMessageTokens`/`estimateStepCompletionTokens` via `approxTokenCount(s)=(len+3)/4` when `usageIsZero` | `internal/agent/usage_fallback.go:18-176` |
| Wired usage per step | `stepMessages` cloned in `PrepareStep`, `fallbackStepUsage(stepMessages, stepResult)` then `updateSessionUsage` | `internal/agent/agent.go:860-862`, `internal/agent/agent.go:1027-1028` |
| Cost branches | `openrouterCost` reads `ProviderMetadata["openrouter"]`; `extractHyperCredits` reads `ProviderMetadata["openai"].extraField("remaining")` | `internal/agent/agent.go:1874-1904` |
| Cancel taxonomy | `isAbortError` (context.Canceled/DeadlineExceeded) → not retried; mapped to `FinishReasonCanceled`; detached `cleanupCtx` preserves writes | `fantasy/retry.go:198-205` + `internal/agent/agent.go:1159-1161` |
| Partial output on failure | On exhaust, partial text remains under `FinishReasonError` in DB via `cleanupCtx`; on retry, reset keeps first-attempt fragment from concatenating | `internal/agent/agent.go:1082-1141` |
| No jitter/budget caps | No `jitter` symbol; no crush-owned rate limiter or token bucket | `internal/agent/coordinator.go:340-622` (absence) |

## Answers to Dimension Questions

**Is the provider interface shaped by Aren semantics or by one vendor SDK?**

By `fantasy`, not Aren and not a single vendor, but vendor detail still leaks. `Coordinator` (`internal/agent/coordinator.go:1099-1162`) creates a concrete SDK per `providerCfg.Type` (8+ branches plus custom `litellm/ollama/omlx`), and `getProviderOptions` (`internal/agent/coordinator.go:360-612`) contains a vendor-keyed `switch` that parses per-provider options (`openai.ParseOptions`, `anthropic.ParseOptions`, `google.ParseOptions`, etc.), injects `reasoning_effort`/`thinking`/`extra_body` per provider ID (e.g. `catwalk.InferenceProviderIoNet`, `AlibabaSingapore`, `MiniMax`), and sets `anthropic-beta: interleaved-thinking-2025-05-14` for Anthropic thinking (`internal/agent/coordinator.go:1106-1112`). The neutral surface crush programs to is `fantasy.Provider`/`fantasy.LanguageModel`/`fantasy.AgentStreamCall`/`fantasy.ProviderOptions` (`internal/agent/agent.go:154`, `internal/agent/agent.go:685`), with transport details normalized to `fantasy.Usage`/`FinishReason`/`ProviderMetadata`/`CallWarning`. So the boundary is neutral at the `Stream` call site, but the translation layer is vendor-parameterized in crush itself — Aren would want that behind a provider registry rather than an expanding switch.

**Can a successful HTTP status still produce a typed execution failure?**

Yes — two paths. (1) Mid-stream SSE error events inside a 200 are classified via `fantasy/errors.go:257` `TransientStreamErrorTypes = {server_error, overloaded_error, api_error, rate_limit_error, internal_error}` into `ProviderError{TransientError:true}`, causing `IsRetryable=true` (`fantasy/errors.go:63`) and a retry via `RetryWithExponentialBackoffRespectingRetryHeaders` despite `StatusCode==200`. (2) A 200 that delivers `FinishReasonContentFilter` (Anthropic `stop_reason=refusal`, OpenAI `content_filter`) is typed as `message.FinishReasonContentFilter` (`internal/agent/agent.go:995-1005`) and persisted for the TUI's REFUSED banner. An incomplete stream without terminal (`NewIncompleteStreamError` with `io.ErrUnexpectedEOF` cause) is also retryable (`fantasy/errors.go:169-177`). Crush's `OnStepFinish` and post-stream error mapping then convert any remaining `*fantasy.ProviderError`/`*fantasy.Error`/`IsTransportError` into `FinishReasonError` with `Title`/`Message` (`internal/agent/agent.go:1163-1182`), so a 200 can land as a non-retryable typed `FinishReasonError` or `ContentFilter`.

**Which failures are safe to retry, and who owns the attempt budget?**

Taxonomy is owned by `fantasy/errors.go` and `fantasy/retry.go`, with crush contributing a one-shot credential repair:

*Retryable (transport/server):* `ProviderError.IsRetryable` true when `TransientError` set, `Cause==io.ErrUnexpectedEOF`, `IsTransportError` (HTTP/2 stream/connection/GOAWAY/reset or header `x-should-retry: true`), or status `408/409/429/5xx` (`fantasy/errors.go:71-97`). Covers `NewIncompleteStreamError`, `WrapTransportError`, network `net.Error`, and mid-stream `server_error` etc. These trigger `retryWithExponentialBackoff` with `DefaultRetryOptions MaxRetries=3, InitialDelay=5s, BackoffFactor=2` (`fantasy/retry.go:88-94`), delay overridden by `Retry-After`/`retry-after-ms` (`fantasy/retry.go:18-55`). Crush's `OnRetry` hook (`internal/agent/agent.go:931`) runs after delay is chosen but before it elapses; crush resets content there but does not change the budget.

*Auth retry (single shot):* `fantasy/retry.go:60-78` `RetryWithExponentialBackoffRespectingRetryHeaders` checks `isAuthError` (`ProviderError.StatusCode==401 || AuthError`) and, if `OnAuthRefresh` is set, calls it once; on `nil` it restarts the entire retry pass with a fresh budget. Crush supplies that hook via `makeAuthRefreshCallback` (`internal/agent/coordinator.go:1349`), which delegates to `retryAfterUnauthorized` (`internal/agent/coordinator.go:1274`) — OAuth refresh (`refreshOAuth2Token`), APIKey re-resolve, or `refreshAWSCredentials` (spawns `sh -c` 5m detached, streams URL via `notify.TypeAWSSSOAuth` — `internal/agent/aws_sso_refresh.go:45-90`). Interactive revoke waits 5m via `WaitForTokenChange` (`internal/agent/coordinator.go:1313-1318`). At most one refresh is attempted; a second 401 surfaces the original error.

*Non-retryable:* `400`/`403`/`404`/`422` etc. without retry marker, plus `isAbortError` (`context.Canceled`/`DeadlineExceeded`) which bypasses retry (`fantasy/retry.go:130`). Validation failures (invalid tool JSON → `sanitizeToolInput`, orphan tool calls) are repaired locally, not by transport retry.

*Ownership:* fantasy owns counting/backoff/jitter(untuned)-free logic; coordinator owns credential refresh policy; sessionAgent owns content-reset side-effect.

**Are token usage and partial output preserved when a stream fails late?**

Partly preserved, with last-error-wins semantics:

*Partial output:* Text/reasoning/tool_call deltas are appended to `currentAssistant` and flushed via `message.Service.Update` debounced + terminal flush (`internal/message/message.go:218-274`). On retry, `OnRetry` does `currentAssistant.ResetStreamedContent()` and `messages.Update(genCtx, *currentAssistant)` (`internal/agent/agent.go:937-940`), discarding the failed-attempt fragment so the next attempt doesn't concatenate. The comment explicitly notes *on final attempt (no more retries), any partial content stays under the error* (`internal/agent/agent.go:935-936`). On terminal stream error, crush maps it to `FinishReasonError` and writes with a detached `cleanupCtx` (5s, `WithoutCancel`) so shutdown cancel can't drop it, leaving the partial text visible beneath the error (`internal/agent/agent.go:1097-1188`). Dropped tool calls that were started but never finished are autoclosed with `Input="{}"` and synthetic `ToolResult{IsError:true, Content:"There was an error while executing the tool"}` (`internal/agent/agent.go:1107-1154`).

*Token usage:* Per-step, `OnStepFinish` captures `stepResult.Usage` and, if zero, falls back via `fallbackStepUsage(stepMessages, stepResult)` (`internal/agent/agent.go:1027`) which estimates with `approxTokenCount` (`internal/agent/usage_fallback.go:171-176`) and flags `estimated=true`. `updateSessionUsage` (`internal/agent/agent.go:1906-1937`) accumulates `session.Cost` (via catwalk `CostPer1M*` or `openrouterCost`), updates `PromptTokens`/`CompletionTokens` only when non-zero, sets `EstimatedUsage`, and emits `eventTokensUsed` only when not estimated. On the summary path, `usageIsZero` also gates `EstimatedUsage` (`internal/agent/agent.go:1458`). So a late failure that delivered a final `Usage` frame contributes exact counts; a failure that truncated before usage falls back to an `estimated` count and zero cost, but the partial text and any tool results are durably recorded.

## Architectural Decisions

* **Delegate provider breadth to `fantasy` and `catwalk`.** Builds on `charm.land/fantasy` for streaming agent loop and on `charm.land/catwalk` for catalog (models/providers/pricing/context windows). Keeps `crush` free of per-provider SSE parsers but ties roadmap upgrades to `fantasy` releases — `go.mod:9-10`.
* **Centralized `Config` → `SelectedModel` → `ProviderConfig` → `fantasy.Provider` chain.** `Coordinator.buildAgentModels` resolves large/small models from `config.SelectedModel{Provider, Model}` against `Providers.Get(provider)` and `catwalk.Model` (`internal/agent/coordinator.go:806-889`), then `buildProvider` instantiates the SDK with `BaseURL/APIKey/Headers/ExtraBody` and debug HTTP client. Decision centralizes credential resolution (`cfg.Resolve`) but scatters vendor quirks across `getProviderOptions`.
* **Single `fantasy.Agent.Stream` per turn with rich callbacks.** Rather than polling, crush hands `PrepareStep`/`OnTextDelta`/`OnToolCall` etc. to fantasy. This normalizes streaming across providers but means crush cannot observe raw frames or apply custom backpressure.
* **Vendor-aware but "neutral" request merging.** `jsons.Merge(catwalkOpts, providerCfgOpts, cfgOpts)` + per-type `ParseOptions` (`internal/agent/coordinator.go:388-609`) lets a user-level `provider_options` override everything, while still applying vendor-specific `reasoning_effort`/`thinking` defaults per model capability (`CanReason`, `ReasoningLevels`). Explicit handling for edge providers (`IoNet`, `MiniMax`, `Baseten`, `Fireworks`, `Copilot Responses` — `internal/agent/coordinator.go:522-591`).
* **Media workaround as pre-request translation.** `workaroundProviderMediaLimitations` rewrites tool_result media into text placeholder + follow-up user `FilePart` for providers lacking tool-result media support (`internal/agent/agent.go:2133-2238`). Encapsulates provider capability divergence at one seam.
* **Two-layer retry.** Layer 1: fantasy's exponential backoff over `isRetryableError` with rate-limit header respect. Layer 2: crush's `OnAuthRefresh` that refreshes credentials and restarts Layer 1 once. Separation prevents infinite auth loops while allowing autos-retry of transient blips.
* **Detached contexts for durability.** Terminal writes (`FlushAll`, `publishRunComplete`, canceled-turn persistence, AWS SSO refresh) use `context.WithoutCancel(ctx)` + 5s timeout so workspace shutdown cancel doesn't drop DB / terminal pubsub state (`internal/agent/agent.go:752-756`, `internal/agent/agent.go:1097-1100`, `internal/agent/aws_sso_refresh.go:62`).
* **Queue/cancel marks for prompt-at-a-time execution.** Late-stream failure paths race queue drain and cancellation via `cancelMark` high-water mark and `dispatchMu` — not directly provider-boundary but affects whether a retried stream replays queued prompts.

## Notable Patterns

* **Callback accumulation pattern.** `currentAssistant *message.Message` pointer is shared across `OnTextDelta`/`OnToolInputStart`/`OnToolCall` and `OnStepFinish`/`OnRetry` closures (`internal/agent/agent.go:734-950`), letting the agent treat the DB-backed message as streamed state. `OnRetry` is the only callback allowed to mutate accumulated state by resetting it.
* **`OnRetry` content-reset idiom.** Documented in `fantasy/retry.go:22-33` as consumer responsibility; crush implements it verbatim with `ResetStreamedContent()` — preserves correctness of retry concatenation without buffering original wire bytes.
* **`ProviderOptions` as typed bag.** `fantasy.ProviderOptions` (`map[string]any` of typed per-provider structs) flows from `getProviderOptions` through `PrepareStep` (`internal/agent/agent.go:810-811` clears; `internal/agent/agent.go:845-852` re-adds cache control) to `AgentStreamCall.ProviderOptions`. Avoids `any`-typed leaks in crush code while still allowing provider-specific fields.
* **Estimated-usage flag.** `EstimatedUsage bool` on `session.Session` distinguishes billable exact vs. fallback (`internal/agent/usage_fallback.go:9-33`, `internal/agent/agent.go:1907-1922`). Generated `session.go` stores it alongside token counters, letting UI cost views degrade gracefully.
* **Auth refresh as retry-adapter.** Implementing `OnAuthRefresh func(ctx context.Context, *ProviderError) error` both inside fantasy's retry middleware and inside crush's `Coordinator` bridges interactive auth (browser URL publish via `pubsub`) into a non-branching retry loop without duplicating backoff logic.

## Tradeoffs

* **Abstraction reuse vs. leakage.** Using `fantasy` avoids reimplementing 9 SDKs, but `coordinator.getProviderOptions` still grows linearly with vendor quirks (320+ lines of `switch` on provider type/ID). Adding a provider means touching crush's coordinator, not just `catwalk`.
* **Global retry defaults vs. tunability.** `DefaultRetryOptions` is hard-coded (`MaxRetries=3`, `InitialDelay=5s`, `BackoffFactor=2`, no jitter, no max delay) and crush never overrides `MaxRetries` (`internal/agent/agent.go:685-796` passes none). This is simple and matches OpenAI SDK precedent, but long-running agents pay 5s+10s+20s = 35s worst-case on a transient blip, with no deadline to bound it — a noisy Bedrock endpoint can dominate user-perceived latency.
* **Retry-After compliance vs. abuse.** Honoring `retry-after-ms`/`retry-after` up to 60s (or `exponentialDelay`) correctly backs off on 429/5xx, but a malicious `Retry-After: 60` can stall `crush run` (`ctx` may not have a deadline). No circuit-breaker or per-session rate limiter mitigates it.
* **Single auth refresh vs. rotation.** One retry budget after a successful refresh avoids loop, but a flaky OIDC that intermittently returns 401 within the same turn won't be re-refreshed — the turn fails with `Provider Error: unauthorized`. More resilient would be to refresh per-401 up to a cap.
* **Content-reset simplicity vs. UX continuity.** Discarding failed-attempt text keeps output coherent but drops context a user may have seen streaming before the reset; the final error route preserves partial, but a successful retry erases the blip entirely — an observer sampling pubsub mid-retry sees a flicker.
* **Estimated usage heuristic vs. precision.** `approxTokenCount = (len+3)/4` is provider-agnostic and cheap, but diverges from BPE/Claude tokenizers by 10-30%; the `estimated` flag makes the divergence explicit, yet downstream budgets that sum `session.Cost` may silently underbill/overbill for providers that legitimately emit zero usage.
* **No structured output in crush.** Leaving `fantasy.ObjectCall` unused avoids schema/repair complexity, but Aren Phase 5 expects bounded structured results — crush's pattern offers no reusable guidance for that phase.
* **Media workaround correctness vs. token cost.** Injecting a synthetic user message with `FilePart` for tool images doubles the effective turn length for non-Anthropic providers; vision models pay extra input tokens, text-only models fall back to placeholder — tradeoff between compatibility and cost correctness.

## Failure Modes / Edge Cases

* **Malformed tool JSON bricks follow-on turns if not sanitized.** Without `sanitizeToolInput` (`internal/agent/agent.go:2272`), a provider-emitted truncated JSON tool call would cause the next provider to reject the tool_result with `400 invalid_json` forever. The fix replaces it with `{}` and forces `IsError=true` follow-on (`internal/agent/agent.go:969-972`), but silently swallows the original malformed `input` for observability (only `slog.Warn` + `input_len`).
* **Orphaned tool call/result deadlocks conversation.** `preparePrompt` detects mismatched call/result sets (`internal/agent/agent.go:1543-1687`) and either drops the orphan result or injects a synthetic `ToolResultPart{Error: "tool call was interrupted..."}`. Without this, a canceled turn that emitted `tool_call` but no `tool_result` would make every subsequent fantasy request fail validation.
* **`ResetStreamedContent` races `messages.Update` on cancel.** If `genCtx` is canceled concurrently with `OnRetry`, `a.messages.Update(genCtx, *currentAssistant)` may fail with `context.Canceled`; it is logged but not retried (`internal/agent/agent.go:938-940`). A subsequent successful retry then starts from a DB state that still contains the pre-reset fragment for a brief window until the next `Update` flushes — mitigated by `FlushAll` in the terminal `defer`.
* **`retry-after: <date>` far-future stall.** `getRetryDelayInMs` does `time.Until(t)` on an RFC1123 date with no ceiling beyond the 60s/exponential check (`fantasy/retry.go:48-53`). A date 10 minutes in the future yields a 10-minute sleep, canceled only via `ctx.Done` (`fantasy/retry.go:160-164`). `crush run` without a caller deadline hangs.
* **Context-too-large signal lost on fallback.** `ProviderError.IsContextTooLarge` (`fantasy/errors.go:189`) checks `ContextTooLargeErr || ContextMaxTokens>0`. `fallbackStepUsage` never sets those fields — it only synthesizes `InputTokens/OutputTokens`. A turn that exhausts context but gets zero provider usage will fall back to estimated usage and then trigger `StopWhen` summarization via context-window check (`internal/agent/agent.go:1038-1056`) instead of surfacing a typed `context_too_large` for explicit handling.
* **Rate-limit vs. quota distinguished only by header truthiness.** `IsRetryable` treats any `429` as retryable regardless of whether it's quota exhaustion (non-retryable) vs. burst limit (retryable). `x-should-retry: false` could mark it non-retryable, but most providers don't send that header.
* **AWS SSO refresh blocks on `notify==nil`.** `refreshAWSCredentials` returns `errNoInteractiveAuth` (`internal/agent/aws_sso_refresh.go:47`) — surfaced as original 401 — if run without a notifier (headless/CI). The failure is correctly non-retryable but loses the actionable `run "aws sso login"` hint outside the TUI path.
* **Post-refresh `UpdateModels` can still fail.** If catalog fetch fails after credential refresh, `refreshOAuth2Token`/`refreshAWSCredentials` return that error and fantasy treats it as refresh failure, surfacing original 401 (`fantasy/retry.go:72-74`, `internal/agent/coordinator.go:1365`). The turn fails despite credentials having been refreshed — a transient catwalk outage masks a successful auth.

## Future Considerations

* **Expose per-provider retry policy from config.** Allow `crushrc` `provider_options.max_retries/initial_delay/backoff_factor/jitter` to flow into `fantasy.RetryOptions.MaxRetries` via `AgentCall.MaxRetries` (`fantasy/agent.go:149`) instead of relying on global defaults; add `context.WithTimeout` per `AgentStreamCall` to bound tail latency.
* **Add jitter and max-delay cap.** Replace the deterministic `InitialDelay * BackoffFactor^n` with `jittered = rand(0, delay)` and cap at 30s to avoid thundering herd on multi-agent retries.
* **Formalize structured-output path.** Adopt `fantasy.ObjectCall` for Aren Phase 5 with explicit `Schema` (via `fantasy/schema` + `jsonrepair`) and a bounded validation-repair loop distinct from transport retry (e.g., 1 repair prompt reusing the same `ProviderError`-typed taxonomy), using `NoObjectGeneratedError` to distinguish parse vs. validation failure.
* **Make usage fallback tokenizer-aware.** Replace `(len+3)/4` with a tokenizer lookup per model (e.g., `tiktoken` for OpenAI, Anthropic count) when `catwalk.Model` indicates the family, retaining `estimated` flag but narrowing error band for cost-sensitive budgets.
* **Deadline and rate-limit edge handling.** Clamp `Retry-After` to `min(exponentialDelay, 60s)` and enforce overall `ctx` deadline in `Coordinator.run` (`context.WithTimeout` per turn) so interactive cancellations cannot be delayed by an adversarial header.
* **Auth refresh with cap, not single-shot.** Change `RetryWithExponentialBackoffRespectingRetryHeaders` to allow `k=2` refreshes within one turn (e.g., token expiry mid-tool-loop) with idempotent `UpdateModels`, matching how tool loops can re-hit 401 after a successful refresh.
* **Centralize `getProviderOptions` registry.** Extract the vendor `switch` into a `providerRegistry map[catwalk.Type]func(Model, ProviderConfig) (ProviderOptions, error)` so new providers don't widen `coordinator.go` and Aren's lifecycle core stays free of SDK imports.

## Questions / Gaps

* **No fake-server tests for provider boundaries found in crush.** Workspace search found no `fantasy` fake-server tests for malformed frames/incomplete bodies/rate-limit/cancellation in the crush repo (tests at `internal/agent/*_test.go` use mock `SessionAgent`/`LanguageModel`, not HTTP fakes). Verification of streaming frame corruption and `TransientError` classification relies on `fantasy`'s own `providertests`/`retry_test.go` — not evidence in `studies/aren-go-runtime-study/sources/crush`.
* **No per-request deadline observed.** `AgentStreamCall` construction (`internal/agent/agent.go:796`) sets no `context.WithTimeout`/`Deadline`; cancellation is only cooperative via `Cancel`/`context.Canceled`. Whether the underlying `http.Client` per provider has a default timeout is owned by `fantasy` provider constructors — not visible in crush evidence.
* **No bounded repair loop for structured output.** Since crush never calls `GenerateObject`, the questions "schema conversion, parse/validation failures, repair prompts" and "retry budgets for validation vs. transport" cannot be answered from crush implementation; search returned zero crush call sites beyond test fakes.
* **Partial-output preservation under transport share.** The retained-partial logic is proven for error-token paths but its interaction with provider-executed tools (`ProviderExecuted=true` tool results) during a mid-stream transport retry is not exercised in crush tests — `usage_fallback_test.go` covers token math only.
* **Token estimate granularity.** `estimateMessagePartTokens` handles `TextPart/ReasoningPart/FilePart/ToolCallPart/ToolResultPart` (`internal/agent/usage_fallback.go:93-118`) but `estimateStepCompletionTokens` ignores `DocumentContent`/`Other` content types — if a provider returns those, `InputTokens` will be undercounted while `OutputTokens` fallback diverges.

---

Generated by `dimensions/01.05-provider-boundaries-structured-results-and-retry.md` against `crush`.
