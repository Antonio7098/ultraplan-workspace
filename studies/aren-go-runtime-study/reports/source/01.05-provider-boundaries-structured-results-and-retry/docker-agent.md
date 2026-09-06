# Source Analysis: docker-agent

## 01.05 Provider Boundaries, Structured Results, and Retry

### Source Info

| Field | Value |
|-------|-------|
| Name | docker-agent |
| Path | `studies/aren-go-runtime-study/sources/docker-agent` |
| Language / Stack | Go 1.26 / openai-go v3, anthropic-sdk-go, genai (Gemini), Docker Model Runner |
| Analyzed | 2026-08-29 |

## Summary

docker-agent implements a provider-neutral `Provider` boundary (`pkg/model/provider/provider.go:38`) with `CreateChatCompletionStream(ctx, []chat.Message, []tools.Tool) (chat.MessageStream, error)` and a thin `chat.MessageStreamResponse` normalization (`pkg/chat/chat.go:147`). A `Registry` dispatches on `provider`/`api_type`/`base_url` (`pkg/model/provider/factory.go:104`) to leaf clients for OpenAI (chat completions + responses API), Anthropic (Messages + Beta), Gemini (GenerateContent), DMR and Bedrock. Each leaf owns request translation, stream decoding and error wrapping, while the runtime's `fallbackExecutor` (`pkg/runtime/fallback.go:226`) owns retry budget, backoff with jitter (`pkg/backoff/backoff.go:25`), rate-limit handling (429 vs 5xx) and cooldown pinning. Structured output is dual-mode: native JSON schema forwarded via `options.WithStructuredOutput` (`pkg/model/provider/options/options.go:118`) for genuine structured-output providers, and a tool-mode internal tool `__structured_output__` (`pkg/tools/builtin/structuredoutput/structuredoutput.go:27`) with validation + bounded repair via runtime reminder loop (`pkg/runtime/structured_output.go:23`). Retry classification is centralized in `pkg/modelerrors/modelerrors.go:581` via `StatusError` (`pkg/modelerrors/modelerrors.go:25`) wrapped by each SDK adapter (`pkg/model/provider/oaistream/wrap.go:15`, `pkg/model/provider/anthropic/wrap.go:27`, `pkg/model/provider/gemini/wrap.go:16`). Stream adapters normalize to `chat.MessageStreamResponse` but discard partial output on late failure; usage accounting is preserved only via final `response.completed` / last chunk and is stored via `handleStream.recordUsage` (`pkg/runtime/streaming.go:154`). Cancellation propagates via `context.Canceled/DeadlineExceeded` as non-retryable at every layer (`pkg/modelerrors/modelerrors.go:587`, `pkg/runtime/fallback.go:202`).

## Rating

**6 / 10** — Provider interface is neutral and retry/overflow vocabularies are centralized, but request translation still embeds extensive provider-specific escape hatches (`provider_opts` top_k, billing, api_type, transport=websocket), structured-output validation repair is bounded but purely runtime-side (no provider-native repair prompt), and partial output + usage on late stream truncation are not preserved (failed stream returns `Stopped:true` with no content).

## Evidence Collected

Every entry MUST include a file path with line numbers. Format: `path/to/file.ts:NN`.

| Area | Evidence | File:Line |
|------|----------|-----------|
| Provider interface (neutral) | `type Provider interface { ID() modelsdev.ID; CreateChatCompletionStream(...)(chat.MessageStream,error); BaseConfig() base.Config }` | `pkg/model/provider/provider.go:38` |
| Provider registry & dispatch | `type Factory func(...) (Provider,error)` + `Registry.NewWithModels` routing to `createDirectProvider` via `resolveProviderType` | `pkg/model/provider/factory.go:17` and `pkg/model/provider/factory.go:104` |
| Adapter wrapping for observability | leaf wrapped with `instrumentProvider(p)` so `chat {model}` CLIENT span lives on leaf only | `pkg/model/provider/factory.go:122` |
| Base config shared by all leaves | `type Config { ModelConfig; ModelOptions; Env; ... }` + `CapsOverride()` | `pkg/model/provider/base/base.go:13` and `pkg/model/provider/base/base.go:63` |
| Model options (gateway, structured output, openAIVendor) | `type ModelOptions { gateway, structuredOutput, openAIVendor, transportWrapper }` + `WithOpenAIVendor` trusted bit | `pkg/model/provider/options/options.go:12` and `pkg/model/provider/options/options.go:164` |
| Full provider set | `DefaultFactories()` maps 7 provider types `openai, anthropic, google, dmr, amazon-bedrock, openai_*` | `pkg/model/provider/providers/providers.go:22` |
| Gateway abstraction | `VerifyDockerGatewayAuth`, `GatewayAuthToken`, `GatewayHTTPOptions` shared across leaves | `pkg/model/provider/base/gateway.go:21` |
| OpenAI request translation (chat completions) | `CreateChatCompletionStream` builds `openai.ChatCompletionNewParams`, applies `ConvertParametersToSchema`, sampling opts, reasoning | `pkg/model/provider/openai/client.go:351` |
| OpenAI responses API path | `CreateResponseStream` builds `responses.ResponseNewParams`, deferred tools, reasoning | `pkg/model/provider/openai/client.go:613` |
| Anthropic request translation | `CreateChatCompletionStream` builds `anthropic.MessageNewParams`, counts tokens for clamping | `pkg/model/provider/anthropic/client.go:211` |
| Anthropic single-retry on overflow | `retryableStream.next()` retries once on `isContextLengthError` | `pkg/model/provider/anthropic/retry.go:8` |
| Gemini request translation | `CreateChatCompletionStream` builds `genai.GenerateContentConfig`, converts messages | `pkg/model/provider/gemini/client.go:717` |
| Schema conversion (OpenAI strict) | `ConvertParametersToSchema` makes all required, forces `additionalProperties:false`, ensures types, detects strict-compatibility | `pkg/model/provider/openai/schema.go:20` |
| Schema conversion (Gemini) | `normalizeBooleanSchemas`, `normalizeTypeFields`, `normalizeEnumValues` to fit `genai.Schema` | `pkg/model/provider/gemini/client.go:561` |
| Schema conversion (Anthropic/DMR) | `ConvertParametersToSchema` thin `tools.ConvertSchema` wrapper | `pkg/model/provider/anthropic/client.go:757` and `pkg/model/provider/dmr/schema.go:7` |
| Neutral chat vocabulary | `MessageRole`, `Message`, `MessageStreamResponse`, `MessageStream interface {Recv, Close}` | `pkg/chat/chat.go:20` and `pkg/chat/chat.go:197` |
| Usage vocab | `type Usage { InputTokens, OutputTokens, CachedInputTokens, CacheWriteTokens, ReasoningTokens }` + `Add()` | `pkg/chat/chat.go:175` |
| OpenAI chat stream accumulator | `StreamAdapter.Recv` parses `openai.ChatCompletionChunk`, aggregates toolCalls map, tracks finishReason | `pkg/model/provider/oaistream/adapter.go:35` |
| OpenAI responses stream accumulator | `ResponseStreamAdapter.Recv` handles `response.output_text.delta`, dedup via `output_index`, buffers args | `pkg/model/provider/openai/response_stream.go:83` |
| Anthropic stream adapter (incl. thinking) | `streamAdapter` (standard) + `betaStreamAdapter` share `retryableStream` | `pkg/model/provider/anthropic/client.go:328` and `pkg/model/provider/anthropic/retry.go:24` |
| Gemini stream adapter | `StreamAdapter` iterating `genai.GenerateContentStream` | `pkg/model/provider/gemini/wrap.go:1` (adapter logic in `pkg/model/provider/gemini/client.go:778`) |
| Error vocabulary (StatusError) | `type StatusError { StatusCode, RetryAfter, Err }` + `WrapHTTPError` | `pkg/modelerrors/modelerrors.go:25` and `pkg/modelerrors/modelerrors.go:52` |
| Provider error wrapping | `WrapOpenAIError`, `wrapAnthropicError` (maps 200-in-band SSE errors to HTTP codes), `wrapGeminiError` | `pkg/model/provider/oaistream/wrap.go:15`, `pkg/model/provider/anthropic/wrap.go:27`, `pkg/model/provider/gemini/wrap.go:16` |
| Retry classification | `ClassifyModelError` distinguishes 429 (rateLimited), 5xx/529/408 retryable, context cancel, overflow, transient pattern | `pkg/modelerrors/modelerrors.go:581` |
| Retryable status codes | `isRetryableStatusCode` returns true for 500/502/503/504/529/408, false for 429 | `pkg/modelerrors/modelerrors.go:332` |
| Backoff + jitter | `Calculate(attempt)` exponential 200ms*2^attempt capped 2s ±10%, `SleepWithContext` cancels on ctx | `pkg/backoff/backoff.go:25` and `pkg/backoff/backoff.go:50` |
| Rate-limit Retry-After | `parseRetryAfterHeader` parses seconds or HTTP-date, capped at `MaxRetryAfterWait=60s` in fallback | `pkg/modelerrors/modelerrors.go:549` and `pkg/runtime/fallback.go:438` |
| Attempt budget ownership | `fallbackExecutor.execute` loops `maxAttempts=1+fallbackRetries` per model in chain, `fallbackRetries` from agent or DefaultRetries=2 | `pkg/runtime/fallback.go:259` and `pkg/modelerrors/modelerrors.go:72` |
| Fallback chain + cooldown | `buildModelChain`, `chainStartIndex`, `recordSuccess` pin fallback on non-retryable primary failure | `pkg/runtime/fallback.go:74`, `pkg/runtime/fallback.go:153`, `pkg/runtime/fallback.go:172` |
| Runtime stream handling | `handleStream` reads via goroutine + `recvCh`, enforces `defaultStreamIdleTimeout=5m`, cancels child context on idle, emits `AgentChoice` events | `pkg/runtime/streaming.go:73` and `pkg/runtime/streaming.go:34` |
| Stream error loses partial | `handleStream` on `Recv` error: `return Stopped:true, fmt.Errorf("error receiving from stream: %w", err)` discards `fullContent`/`Usage` | `pkg/runtime/streaming.go:186` |
| Usage preservation (success) | `handleStream.recordUsage` stores `Input+Cached+Write -> Total`, `tel.RecordTokenUsage` once per stream | `pkg/runtime/streaming.go:154` |
| Structured output native forward | `options.WithStructuredOutput` drops tool-mode, else forwards schema to `responses.ResponseNewParams.Text.Format` / Gemini `ResponseJsonSchema` | `pkg/model/provider/options/options.go:118` and `pkg/model/provider/openai/client.go:746` |
| Structured output tool-mode validation | `structuredoutput.New` compiles schema via `gojsonschema`, `Validate` compacts + checks `result.Valid()` | `pkg/tools/builtin/structuredoutput/structuredoutput.go:40` and `pkg/tools/builtin/structuredoutput/structuredoutput.go:128` |
| Structured output retry/bounded repair | `handleStructuredOutputCalls` validates exclusive call, returns `IsError` tool_result for model to retry; `maxStructuredOutputReminders=2` transient system reminder then `ErrorCodeStructuredOutputFailed` | `pkg/runtime/structured_output.go:102` and `pkg/runtime/structured_output.go:27` |
| Provider-specific escape hatches | `top_k` forwarded for Anthropic/Gemini, `operationName` vs `modelName` for Claude, `thinkingBudget` mapping diverges per provider | `pkg/model/provider/anthropic/client.go:294`, `pkg/model/provider/gemini/client.go:425`, `pkg/model/provider/openai/sampling_opts.go:1` |
| OpenAI-compatible shims | `shouldMergeConsecutiveMessages` for Qwen/vLLM, `claudeSchemaInstruction` prompt-injects schema when `response_format` dropped | `pkg/model/provider/openai/client.go:293` and `pkg/model/provider/openai/claude_structured_output.go:55` |
| Transport escape hatch | `getTransport` reads `provider_opts.transport=websocket`, WebSocket disabled when gateway/transportWrapper present | `pkg/model/provider/openai/client.go:865` and `pkg/model/provider/openai/client.go:202` |
| HTTP middleware for error bodies | `ErrorBodyMiddleware` rewrites non-`{"error":{}}` bodies so SDK preserves details | `pkg/model/provider/oaistream/middleware.go:20` |
| Overflow classification | `OverflowKind` tokens/wire/media, `classifyOverflow` tier1 structured fields + tier2 substring, `ContextOverflowError` | `pkg/modelerrors/modelerrors.go:90` and `pkg/modelerrors/modelerrors.go:238` |
| Cancellation distinction | `ClassifyModelError` returns non-retryable for `context.Canceled/DeadlineExceeded`; `fallback.classifyAttemptError` bubbles immediately; `IsStreamTruncationError` excludes cancellation | `pkg/modelerrors/modelerrors.go:587`, `pkg/runtime/fallback.go:202`, `pkg/modelerrors/modelerrors.go:428` |

## Answers to Dimension Questions

### Is the provider interface shaped by Aren semantics or by one vendor SDK?

Primarily Aren semantics, but with deliberate vendor SDK leakage at translation edges.

The public contract `Provider` (`pkg/model/provider/provider.go:38`) is neutral: it takes `context.Context`, `[]chat.Message` (`pkg/chat/chat.go:52` — roles `system/user/assistant/tool`), `[]tools.Tool` and returns `chat.MessageStream` (`pkg/chat/chat.go:197`) whose chunks are `MessageStreamResponse` (`pkg/chat/chat.go:166`) with `Choices[].Delta` (`MessageDelta`, `pkg/chat/chat.go:148`) and `Usage` (`pkg/chat/chat.go:175`). No OpenAI `ChatCompletionChunk` or Anthropic `MessageStreamEvent` types leak into the interface. `base.Config` (`pkg/model/provider/base/base.go:13`) exposes only `ModelConfig` (YAML-derived), `ModelOptions` (gateway, caps), and `modelsdev.ID` rather than SDK clients.

However leaf implementations re-export SDK decisions through `ProviderOpts` and `ModelOptions.OpenAIVendor` trusted bit: `top_k` forwarding differs by provider (`pkg/model/provider/anthropic/client.go:294` vs `pkg/model/provider/gemini/client.go:425` vs `pkg/model/provider/openai/sampling_opts.go:1`), reasoning-effort strings branch on `isOpenAIVendor` (`pkg/model/provider/factory.go:91` → `pkg/model/provider/options/options.go:164`), Asking for `transport: websocket` (`pkg/model/provider/openai/client.go:866`) disables wrapping, and Claude structured output falls back to prompt injection (`pkg/model/provider/openai/claude_structured_output.go:55`) because OpenAI-compatible gateways drop `response_format`. The injection site `isOpenAIVendor` is intentionally isolated as a factory-set internal bit not reachable from YAML (`pkg/model/provider/options/options.go:164` comment), but `provider_opts` itself remains an untyped `map[string]any` escape hatch.

**Verdict:** interface core is provider-neutral; translation layer is intentionally vendor-aware to keep the core clean, at the cost of a loosely-typed `provider_opts` escape hatch.

### Can a successful HTTP status still produce a typed execution failure?

Yes. Three distinct mechanisms produce typed failures from HTTP 200 streams or validated-but-rejected successes:

1. **In-band SSE error with 200**: Anthropic mid-stream `type: error` events arrive as SSE with HTTP 200. `wrapAnthropicError` (`pkg/model/provider/anthropic/wrap.go:36`) maps `shared.ErrorType.*` → HTTP status (529→529, api_error→500, overloaded→529) so `modelerrors.ClassifyModelError` (`pkg/modelerrors/modelerrors.go:599`) treats it identically to a transport 5xx and triggers retry/fallback.

2. **Context-overflow as typed error**: Substring/structured detection (`pkg/modelerrors/modelerrors.go:238`) wraps the underlying error in `ContextOverflowError` (`pkg/modelerrors/modelerrors.go:118`). `handleStream` has a single-retry path for Anthropic (`pkg/model/provider/anthropic/client.go:332` → `retry.go:24`) that calls the `countTokens` API and re-issues with clamped `maxTokens`; the runtime `fallback.execute` path also wraps as `NewContextOverflowError` (`pkg/runtime/fallback.go:370`) to trigger auto-compaction rather than surfacing HTTP text.

3. **Structured-output validation failure**: Tool-mode structured output (`pkg/tools/builtin/structuredoutput/structuredoutput.go:128`) always succeeds at HTTP level but `Validate` against `gojsonschema` can fail. `handleStructuredOutputCalls` (`pkg/runtime/structured_output.go:155`) does **not** return an execution error to the caller; it records an `IsError=true` tool_result (`pkg/tools/builtin/structuredoutput/structuredoutput.go:157`) and lets the model retry on next iteration. Native mode forwards schema (`pkg/model/provider/openai/client.go:746`) and a successful HTTP response with invalid JSON would surface only as a downstream parse error (no typed failure type is synthesized).

A successful HTTP status therefore very much can produce `StatusError`, `ContextOverflowError`, or a structured-output `ToolCallResult.IsError` execution state, each with distinct handling.

### Which failures are safe to retry, and who owns the attempt budget?

**Safe to retry (same model, backoff + jitter):**
- HTTP 500, 502, 503, 504, 529 (Anthropic overloaded), 408 (`pkg/modelerrors/modelerrors.go:332`)
- Transient pattern override: `number of function response parts` on otherwise 400 (`pkg/modelerrors/modelerrors.go:359`)
- Network: `timeout`, `connection reset/refused`, `no such host`, `overloaded` (`pkg/modelerrors/modelerrors.go:384`)
- Stream truncation: `unexpected eof`, `unexpected end of json input` (`pkg/modelerrors/modelerrors.go:415`) — cleared for DMR/llama.cpp idle prefill resets, checked before non-retryable `invalid` pattern

**Not safe to retry same model (skip to next in chain):**
- 429 rate-limit — `rateLimited=true` (`pkg/modelerrors/modelerrors.go:600`) breaks the per-model loop (`pkg/runtime/fallback.go:413`)
- 4xx (400, 401, 403, 404), `rate limit`, `quota`, `invalid`, `unauthorized`, `api key` (`pkg/modelerrors/modelerrors.go:447`)
- `ContextOverflowError` of any kind (`pkg/modelerrors/modelerrors.go:593`) — documented as "retrying same oversized payload will always fail"

**Budget ownership:**
- `fallbackExecutor` owns the entire budget (`pkg/runtime/fallback.go:226`). Per-model `maxAttempts = 1 + fallbackRetries` (`pkg/runtime/fallback.go:259`), where `fallbackRetries` comes from `agent.FallbackRetries()` or `modelerrors.DefaultRetries=2` (`pkg/modelerrors/modelerrors.go:72`, `pkg/runtime/fallback.go:136`). `backoff.Calculate(attempt)` (`pkg/backoff/backoff.go:25`) gives 200ms·2^attempt capped 2s ±10%. 429 with `retryOnRateLimit` enabled honors `Retry-After` header (`pkg/modelerrors/modelerrors.go:549`) capped at `MaxRetryAfterWait=60s` (`pkg/backoff/backoff.go:20`), otherwise 429 skips immediately. The cooldown manager (`pkg/runtime/cooldown_manager.go`, wired in `pkg/runtime/fallback.go:53`) owns cross-invocation pinning: after a fallback rescues a non-retryable primary failure, `recordSuccess` (`pkg/runtime/fallback.go:172`) pins for `DefaultCooldown=1m` (`pkg/modelerrors/modelerrors.go:76`).
- Anthropic leaf owns a separate **single** retry for context-length specifically (`pkg/model/provider/anthropic/retry.go:24`), independent of `fallbackExecutor`'s budget and limited to one token-counting retry.
- No caller-owned attempt budget beyond that: user cannot configure per-error backoff curves; only `retries` and `cooldown` on the agent.

### Are token usage and partial output preserved when a stream fails late?

**Partial output: No.** `handleStream` (`pkg/runtime/streaming.go:73`) accumulates `fullContent`, `fullReasoningContent`, `toolCalls` and `messageUsage` on success, but on any `Recv` error (including mid-stream truncation, adapter `WrapOpenAIError`/`wrapAnthropicError`, or `idleTimer` firing) it returns `streamResult{Stopped:true}, fmt.Errorf("error receiving from stream: %w", err)` with zero content (`pkg/runtime/streaming.go:186`). The `fullContent` builder and `toolCalls` slice are discarded. The idle path (`pkg/runtime/streaming.go:342`) similarly returns empty result plus `errStreamIdle` (`pkg/runtime/streaming.go:39`). `fallbackExecutor.execute` (`pkg/runtime/fallback.go:338`) does not persist partial either; it classifies the error and either retries the same model or falls through to next model, so late-stream text already emitted via `events.Emit(AgentChoice(...))` (`pkg/runtime/streaming.go:326`) is visible in the UI but not returned as a result nor attached to `streamResult.Content`.

**Usage: Partially preserved on success, lost on failure.** On success, `handleStream` keeps the latest `messageUsage` observed across chunks (`pkg/runtime/streaming.go:192`) and `recordUsage` (`pkg/runtime/streaming.go:154`) persists `Input = InputTokens+Cached+CacheWrite` and `OutputTokens` to the session and emits `RecordTokenUsage`. OpenAI responses path correctly populates only `response.completed` (`pkg/model/provider/openai/response_stream.go:349`); chat completions and Gemini paths populate incremental chunks but `handleStream` overwrites with latest, so last value wins. On late failure, `messageUsage` is still discarded (same error return as above), though intermediate chunks may have updated `sess` via prior `recordUsage` calls only if a success path had been hit — a truncated failure before final chunk leaves no usage persisted. No `Usage` is attached to the error itself, so callers cannot attribute tokens to a failed attempt.

**Deadline/rate-limit preservation:** `Retry-After` is extracted into `StatusError.RetryAfter` (`pkg/modelerrors/modelerrors.go:56`) via `parseRetryAfterHeader` (`pkg/modelerrors/modelerrors.go:549`) and honored only on 429 retry (`pkg/runtime/fallback.go:430`); otherwise jittered backoff (`pkg/backoff/backoff.go:42`) applies. Deadlines rely on `context.DeadlineExceeded` transitively via `ctx.Done()` (`pkg/runtime/streaming.go:337`, `pkg/runtime/fallback.go:263`) — no per-request deadline is synthesized; the operation simply aborts and returns `ctx.Err()`.

## Architectural Decisions

| Decision | Location | Rationale | Consequence |
|----------|----------|-----------|-------------|
| Minimal `Provider` interface over SDK clients | `pkg/model/provider/provider.go:38` | Keeps runtime neutral; callers cannot import `openai-go` without opting into `providers.NewDefaultRegistry` | Translation burden pushed into each leaf; `provider_opts` stays untyped to keep interface small |
| `Registry` + `ProviderFactory` dispatch | `pkg/model/provider/factory.go:17` | Allows YAML `provider: anthropic` → constructor without compile-time SDK dependency | Factories registered globally; `DefaultRegistry` comment warns to use explicit registry for full provider set |
| Trusted `openAIVendor` bit not from YAML | `pkg/model/provider/options/options.go:164` and `pkg/model/provider/factory.go:91` | Prevents user from spoofing `reasoning_effort: none` onto non-OpenAI alias | Gating is opaque; audit requires tracing factory path |
| `ctx` + `MessageStream` streaming, not `chan` | `pkg/chat/chat.go:197` | SSE fits SDK `Next()/Current()/Err()` idiom; allows adapter reuse | Caller must enforce idle timeout via goroutine + `select` on `recvCh` (`pkg/runtime/streaming.go:89`) |
| `StatusError` as neutral wrapper | `pkg/modelerrors/modelerrors.go:25` | All SDK errors converge to `(StatusCode, RetryAfter)` for `ClassifyModelError` without provider imports | Providers must remember to wrap; unwrapped SDK errors fall to regex extraction (`pkg/modelerrors/modelerrors.go:309`) |
| `fallbackExecutor` owns retry+cooldown | `pkg/runtime/fallback.go:43` | Isolates per-agent sticky fallback state; primary+fallbacks tried in one span `runtime.fallback` | Retry policy not per-provider configurable; single global `DefaultRetries/DefaultCooldown` |
| Anthropic single overflow retry inside leaf | `pkg/model/provider/anthropic/retry.go:24` and `pkg/model/provider/anthropic/client.go:332` | Clamp `maxTokens` after calling `CountTokens` API without surfacing to runtime first | Double retry budget: leaf retry + `fallbackExecutor` retry; token-count API adds latency on overflow path |
| Tool-mode structured output as internal tool | `pkg/tools/builtin/structuredoutput/structuredoutput.go:27` and `pkg/runtime/structured_output.go:62` | Unified with normal tool loop: validation error becomes tool_result the model can repair | Extra turn latency vs native; bounded by `maxStructuredOutputReminders=2` (`pkg/runtime/structured_output.go:27`) before typed failure |

## Notable Patterns

- **Shim for OpenAI-compatible servers**: `shouldMergeConsecutiveMessages` (`pkg/model/provider/openai/client.go:293`) coalesces same-role system messages for Qwen3/vLLM/OVHcloud that reject >1 system message; `MergeConsecutiveMessages` (`pkg/model/provider/oaistream/messages.go`) keeps DMR parity.
- **Provider-specific thinking budgets mapped to neutral YAML**: single `thinking_budget: {tokens, effort}` in `latest.ModelConfig` fans out to `reasoning_effort` (OpenAI), `thinking budget_tokens` (Anthropic), `ThinkingBudget` vs `ThinkingLevel` (Gemini 2.5 vs 3) (`pkg/model/provider/gemini/client.go:463`).
- **Error body preservation middleware**: `ErrorBodyMiddleware` (`pkg/model/provider/oaistream/middleware.go:20`) rewrites string/plain-text error bodies into `{"error": ...}` so `openai-go` gjson extraction retains details; paired with `parseProviderError` (`pkg/modelerrors/modelerrors.go:687`) lifting structured details for user display.
- **Cache control via provider opts**: prompt-cache breakpoints flow from `options.WithMaxTokens` and `messageCacheBreakpoints` into provider-specific `CacheControl` fields (Anthropic ephemeral), but stored messages strip them on ingestion (`pkg/chat/chat.go:105`).
- **WebSocket/SSE transport abstraction**: `responseEventStream` (`pkg/model/provider/openai/event_stream.go:5`) abstracts SSE vs WebSocket behind `Next/Current/Err/Close`; OpenAI client chooses at runtime and falls back on failure (`pkg/model/provider/openai/client.go:773`).

## Tradeoffs

- **Neutral interface vs expressiveness**: `Provider.CreateChatCompletionStream` cannot express embeddings (`EmbeddingProvider`, `pkg/model/provider/provider.go:57`), reranking (`RerankingProvider`, `pkg/model/provider/provider.go:71`), or batch; those live on sub-interfaces accessed via type assertion, so callers that need them must probe capabilities.
- **Untyped `provider_opts` escape hatch**: Allows per-provider knobs (`top_k`, `api_version`, `transport`, `supports_deferred_tools`) without interface churn, but loses compile-time safety and central audit. Misspelled keys are silently ignored (e.g. `top_k` vs `topK`) — `GetProviderOptBool` (`pkg/model/provider/anthropic/client.go:295`) returns ok=false.
- **Single-retry overflow inside Anthropic leaf vs runtime compaction**: avoids one extra round-trip through `fallbackExecutor` on token overflow, but funds a blocking `CountTokens` API call (`pkg/model/provider/anthropic/client.go:815`) that may itself fail, and the generic overflow path (`NewContextOverflowError` → compaction) is bypassed for that first retry.
- **Discard partial on error**: Simplifies `streamResult` contract (empty on error), but loses already-generated text for UI/usage and makes late truncation unrecoverable without re-generating entire turn.
- **Global `DefaultRetries=2`**: Sensible for transient 5xx without configuration, but interacts multiplicatively with fallback chain length; a 3-model chain can make 9 attempts (3 models × 3 attempts) before ultimate failure.

## Failure Modes / Edge Cases

- **HTTP 200 with error event not classified as overflow**: `wrapAnthropicError` (`pkg/model/provider/anthropic/wrap.go:36`) maps error types to HTTP codes so `classifyOverflow` (`pkg/modelerrors/modelerrors.go:238`) can inspect `StatusError.StatusCode==413` tier-1; unknown in-band types map to 500 and become retryable 5xx instead of typed overflow — compaction won't trigger.
- **Non-JSON error bodies truncated at 1 MB**: `firstJSONObject` (`pkg/modelerrors/modelerrors.go:744`) caps at 1 MB; huge proxy error bodies beyond that limit produce no structured `parseProviderError` detail and surface raw SDK truncation.
- **Mixed `provider/model` inline specs vs named models**: `resolveRoutedModel` (`pkg/model/provider/factory.go:63`) forbids routed targets with their own `Routing`; error is `"model X has routing rules and cannot be used as a routing target"` — no cycle detection beyond that, so circular fallback definitions deadlock via recursion depth until factory error surfaces.
- **Orphaned function_call without output**: `convertMessagesToResponseInput` (`pkg/model/provider/openai/client.go:1080`) injects placeholder `"(no output — tool call was not executed)"` for orphaned calls (e.g. user cancellation mid-tool); without this the Responses API would reject the next request.
- **Streaming truncated mid-SSE frame**: `StreamAdapter.Recv` (`pkg/model/provider/oaistream/adapter.go:38`) and `ResponseStreamAdapter.Recv` (`pkg/model/provider/openai/response_stream.go:84`) call `oaistream.WrapOpenAIError` on `stream.Err()`; truncation `io.ErrUnexpectedEOF` / `unexpected end of json input` is classified separately as `IsStreamTruncationError` (`pkg/modelerrors/modelerrors.go:428`) and becomes retryable, but `handleStream` drop of partial still applies.
- **Cancellation vs idle distinguishability**: Both set `Strapped:true` but idle uses `errStreamIdle` (`pkg/runtime/streaming.go:39`) as context cause and formats `"model stream stalled after %s"` (`pkg/runtime/streaming.go:353`), while cancellation returns `ctx.Err()` (`pkg/runtime/streaming.go:340`). `ClassifyModelError` exclusion (`pkg/modelerrors/modelerrors.go:587`) prevents idle from being mis-retried as transient.
- **Retry-After overflow**: `parseRetryAfterHeader` (`pkg/modelerrors/modelerrors.go:549`) can parse HTTP-date future times; `handleModelError` (`pkg/runtime/fallback.go:436`) caps at `MaxRetryAfterWait=60s` to avoid indefinite hang from misbehaving gateway.
- **Structured-output collision**: tool name `__structured_output__` reserved (`pkg/tools/builtin/structuredoutput/structuredoutput.go:27`); `appendStructuredOutputTool` (`pkg/runtime/structured_output.go:74`) fails explicitly if user tool claims it.

## Future Considerations

- **Persist partial output + usage on truncation**: Modify `handleStream` to return `streamResult{Content: fullContent.String(), Usage: messageUsage}` alongside the truncation error, allowing `fallbackExecutor` to record usage and (optionally) stream partial to transcript rather than discarding. This aligns with the Aren question on late-stream failure preservation.
- **Centralize overflow budget**: Move Anthropic overflow retry into `fallbackExecutor`'s `ClassifyModelError` path so the `CountTokens` + `clampMaxTokens` (`pkg/model/provider/anthropic/client.go:800`) logic is shared across providers (OpenAI also clamps via `clampMaxTokens` `pkg/model/provider/openai/client.go:328` but with no counting pass). Prevents duplicate budgets.
- **Typed per-provider capability negotiation**: Replace `provider_opts map[string]any` escape hatches with a generated `ProviderCapabilities` struct per provider (codegen from JSON schema) to recover compile-time safety while preserving extensibility.
- **Deadline propagation**: Add per-request `context.WithTimeout` derived from model `contextLimit` / `SamplingConfig` so idle timeout is not the only temporal guard; expose `provider_opts.request_timeout` as a first-class field.
- **Native structured-output repair prompt**: For `structured_output.mode: native`, failed schema validation currently has no retry (error would be surfaced to runtime as tool-result from `structuredoutput.ToolName` only in tool mode). Consider a bounded native repair that re-prompts the model with validation errors inside `handleStructuredOutputCalls` before falling back to tool-mode reminder.

## Questions / Gaps

- No fake-server tests for malformed frames / incomplete success bodies were found under the isolated source scope. `pkg/model/provider/oaistream/middleware_test.go` and `pkg/model/provider/openai/response_stream_test.go` cover streaming adapters but no malformed-SSE-frame injection via `httptest.Server` was located; rate-limit / cancellation coverage relies on unit tests with mocked `StatusError` rather than fake HTTP servers.
- No explicit deadline / context timeout configuration was found for a single model request beyond `defaultStreamIdleTimeout=5m` (`pkg/runtime/streaming.go:32`). Whether `CreateChatCompletionStream` respects a caller `context.WithTimeout` end-to-end (including Anthropic `CountTokens` retry) is exercised only via `ctx.Done()` checks (`pkg/runtime/fallback.go:263`, `pkg/runtime/streaming.go:337`).
- Partial-output handling on late truncation needs empirical verification (e.g., kill downstream `httptest` server mid-chunk and assert whether `handleStream` still emits partial `AgentChoice` events before returning error).
- Schema conversion for Bedrock/Vertex AI not inspected in depth (`pkg/model/provider/bedrock`, `pkg/model/provider/vertexai`) — whether they share the same strict-compatibility enforcement as OpenAI (`pkg/model/provider/openai/schema.go:30`) is not evidenced here.
- Token usage accounting for failed attempts: `fallbackExecutor` records `RecordToolCall` per structured-output attempt (`pkg/runtime/structured_output.go:256`) but does not aggregate `Usage` across retry attempts; cost accounting for retries is not evidenced.

---

Generated by `01.05-provider-boundaries-structured-results-and-retry` against `docker-agent`.
