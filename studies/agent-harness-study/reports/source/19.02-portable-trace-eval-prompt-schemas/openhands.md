# Source Analysis: openhands

## Portable Trace, Eval, and Prompt Schemas

### Source Info

| Field | Value |
|-------|-------|
| Name | openhands |
| Path | `studies/agent-harness-study/sources/openhands` |
| Language / Stack | Python (Pydantic, litellm, Jinja2, MCP, OTel) + React frontend |
| Analyzed | 2026-08-28 |

## Summary

OpenHands is strongly portable at the **tool-schema layer** — every tool cleanly converts between MCP JSON Schema, OpenAI Chat Completions, and OpenAI Responses via a single Pydantic base (`Schema`) with `_process_schema_node` circular-ref handling. Prompt templates use standard Jinja2 `.j2` files with a `FlexibleFileSystemLoader` and are cache-friendly, making them copy-portable but not registry-portable. Traces are the weak link: although OTEL env-var names are read (`OTEL_EXPORTER_OTLP_*`), the implementation is hard-wired to Laminar (`lmnr` SDK, `LaminarLiteLLMCallback`, `Laminar.start_span`/`use_span`), with no abstraction for Jaeger/Honeycomb/Datadog or a vendor-neutral span serialization format. Internal trajectory persistence is Pydantic JSON per-event (`event-*.json`) and a zip export, not an OTel trace or eval dataset. Critic-based evaluation is internal-only (no portable dataset/task harness like SWE-Bench, no JSONL/HF dataset import/export). There are no migration tools for moving traces, evals, or prompts between platforms.

## Rating

**5/10 — Present but inconsistent, weakly documented, fragile**

Tool schemas are mature (explicit MCP/OpenAI/Responses adapters with tests and discriminated unions) which would merit 7-8 alone. Traces, eval datasets, prompt registries, and cross-platform migration pull the aggregate down: tracing is vendor-locked to Laminar despite OTEL-named env vars, eval datasets are not portable, prompt management is filesystem-only Jinja with no versioning/registry API, and no export/import tooling exists for moving data between tracing or eval platforms.

## Evidence Collected

Every entry MUST include a file path with line numbers. Format: `path/to/file.ts:NN`.

| Area | Evidence | File:Line |
|------|----------|-----------|
| Trace format abstraction | `OTEL_EXPORTER_*` keys declared as `_OBSERVABILITY_ENV_KEYS` but `maybe_init_laminar()` immediately branches to `lmnr.Laminar` imports and `LaminarLiteLLMCallback`; no generic OTel span exporter abstraction | `_sdk_inspect/sdk/observability/laminar.py:25-30`, `_sdk_inspect/sdk/observability/laminar.py:57-112` |
| Trace format abstraction | `RootSpan` wraps `Laminar.start_span` / `Laminar.use_span` with `set_trace_session_id`; no OpenTelemetry `Span`/`Tracer` interface, no `opentelemetry-sdk` exporter configuration | `_sdk_inspect/sdk/observability/laminar.py:231-266` |
| Trace format abstraction | `@observe` decorator is a lazy wrapper around `lmnr.observe` (`from lmnr import observe as laminar_observe`); pass-through when `should_enable_observability()==False` — entirely Laminar-typed (`span_type: DEFAULT/LLM/TOOL`) | `_sdk_inspect/sdk/observability/laminar.py:115-196` |
| Trace format abstraction | Conversations own a `RootSpan` via `_observability_root_span` attribute looked up by name at every `observe` entry (`_ROOT_SPAN_ATTR = "_observability_root_span"`) — couples trace hierarchy to Laminar context propagation hack | `_sdk_inspect/sdk/observability/laminar.py:228`, `_sdk_inspect/sdk/conversation/base.py:120-151` |
| Trace format abstraction | `Telemetry` writes LLM call logs as Pydantic-serialized JSON via `_safe_json` (`ModelResponse.model_dump`) to local files or a callback — not OTLP, not Honeycomb/Jaeger format | `_sdk_inspect/sdk/llm/utils/telemetry.py:288-402` |
| Trace format abstraction | Dependencies pin `opentelemetry-api==1.39.1` and `opentelemetry-exporter-otlp-proto-grpc==1.39.1` in lockfile but no code imports `opentelemetry.trace` or configures an OTel `TracerProvider`; import is only via transitive `lmnr` dependency | `pyproject.toml:64-65`, `pyproject.toml:196-197` |
| Dataset portability | `CriticBase.evaluate(events, ...) -> CriticResult` and concrete critics (`AgentFinishedCritic`, `APIBasedCritic`, `EmptyPatchCritic`, `PassCritic`) operate on in-memory `Event` lists — no dataset loader, no JSONL/CSV/Parquet/HF `datasets` integration, no evaluation harness config | `_sdk_inspect/sdk/critic/base.py:65-83`, `_sdk_inspect/sdk/critic/impl/agent_finished.py:26-33`, `_sdk_inspect/sdk/critic/impl/api/critic.py:58-76` |
| Dataset portability | Agent `CriticMixin._should_evaluate_with_critic` and `ResponseDispatch._evaluate_with_critic` gate evaluation per-action/mode (`off`/`finish`/`finish_and_message`) — runtime hook, not a reusable eval dataset | `_sdk_inspect/sdk/agent/critic_mixin.py:34-73`, `_sdk_inspect/sdk/agent/response_dispatch.py:123-275` |
| Dataset portability | No `evaluation/` or `benchmark/` directory, no SWE-Bench harness, no dataset versioning — grep for `eval`/`benchmark` hits only critic/mixing code, not portable datasets | `_sdk_inspect/sdk/agent/agent.py:893-895` (critic gate only) |
| Prompt template portability | Central Jinja2 renderer `render_template(prompt_dir, template_name, **ctx) -> str` via `FlexibleFileSystemLoader` (supports absolute + relative paths), `FileSystemBytecodeCache` under `~/.openhands/cache/jinja`, `refine` filter for win32 terminal rewriting | `_sdk_inspect/sdk/context/prompts/prompt.py:16-45`, `_sdk_inspect/sdk/context/prompts/prompt.py:57-114` |
| Prompt template portability | System prompts are versioned `.j2` files: `system_prompt.j2`, `system_prompt_long_horizon.j2`, `system_prompt_planning.j2`, `security_policy.j2`, plus model-specific overrides `model_specific/anthropic_claude.j2`, `model_specific/google_gemini.j2`, `model_specific/openai_gpt/*.j2` — standard Jinja, copy-portable, no external prompt registry (Langfuse, PromptLayer, etc.) | `_sdk_inspect/sdk/agent/prompts/system_prompt.j2:1-149`, `_sdk_inspect/sdk/agent/prompts/model_specific/anthropic_claude.j2:1`, `_sdk_inspect/sdk/agent/prompts/model_specific/google_gemini.j2:1` |
| Prompt template portability | Condenser prompt `summarizing_prompt.j2` and `ask_agent_template.j2` / `skill_knowledge_info.j2` / `system_message_suffix.j2` are likewise Jinja2 statics with no variable schema or version pinning | `_sdk_inspect/sdk/context/condenser/prompts/summarizing_prompt.j2:1`, `_sdk_inspect/sdk/context/prompts/templates/ask_agent_template.j2:1` |
| Cross-provider schema compatibility | `Schema.to_mcp_schema()` resolves `$ref` via `_process_schema_node` (strips `anyOf`, resolves `$defs`, handles circular refs with shallow `{"type":"object"}` fallback) and strips discriminated-union `kind` fields — produces MCP-compatible JSON Schema | `_sdk_inspect/sdk/tool/schema.py:70-198` |
| Cross-provider schema compatibility | `Schema.from_mcp_schema(model_name, schema)` dynamically constructs `create_model(..., __base__=cls)` mapping JSON Schema `type` → Python type via `py_type()`, preserving `required` vs `T|None` | `_sdk_inspect/sdk/tool/schema.py:201-239` |
| Cross-provider schema compatibility | `ToolDefinition.to_mcp_tool()` / `to_openai_tool()` / `to_responses_tool()` — three explicit adapters: MCP `{inputSchema, outputSchema}`, Chat Completions `ChatCompletionToolParam`, Responses `FunctionToolParam` with `strict: False` | `_sdk_inspect/sdk/tool/tool.py:379-498` |
| Cross-provider schema compatibility | `MessageToolCall` is transport-agnostic with `origin: completion|responses`, `id`, `responses_item_id` for prefix-cache preservation, plus `to_chat_dict()` and `to_responses_dict()` serializers (Responses requires `call_id` + `id` duality) | `_sdk_inspect/sdk/llm/message.py:24-119` |
| Cross-provider schema compatibility | `Message` provides dual serialization: `to_chat_dict(cache_enabled, vision_enabled, function_calling_enabled, ...)` and `to_responses_dict(vision_enabled)` / `from_llm_chat_message` / `from_llm_responses_output` — handles `thinking_blocks`, `reasoning_content`, `responses_reasoning_item` normalization | `_sdk_inspect/sdk/llm/message.py:278-722` |
| Cross-provider schema compatibility | `MCPToolDefinition._create_mcp_action_type()` caches `Schema.from_mcp_schema` per `mcp.types.Tool.name`; `action_from_arguments()` validates against dynamic schema before execution, sanitizes via `exclude_none` and strips `kind` | `_sdk_inspect/sdk/mcp/tool.py:117-223` |
| Cross-provider schema compatibility | `NonNativeToolCallingMixin` enables prompt-mocked function calling for models without native FC: `convert_fncall_messages_to_non_fncall_messages` / `convert_non_fncall_messages_to_fncall_messages` with `STOP_WORDS` injection | `_sdk_inspect/sdk/llm/mixins/non_native_fc.py:34-105` |
| Export/import tools | `EventLog.append()` persists each `Event` as `event.model_dump_json(exclude_none=True)` to `event-{idx:05d}-{event_id}.json` via `FileStore.write()` with `flock` + null `write_guard` — Pydantic JSON, not a standard trace format | `_sdk_inspect/sdk/conversation/event_store.py:132-145` |
| Export/import tools | Persistence constants `BASE_STATE="base_state.json"`, `EVENTS_DIR="events"`, `EVENT_NAME_RE = ^event-(?P<idx>\d{5})-(?P<event_id>...)\.json$` define the on-disk layout; no OTLP/Parquet/Arrow export | `_sdk_inspect/sdk/conversation/persistence_const.py:4-9` |
| Export/import tools | `AppConversationService.export_conversation(conversation_id) -> bytes` / `open_conversation_export -> AsyncGenerator[bytes]` defines the trajectory export contract as a zip of JSON event files + meta — download-only, no import, no cross-platform converter | `openhands/app_server/app_conversation/app_conversation_service.py:160-189` |
| Export/import tools | Live implementation adds Redis distributed lock (`_EXPORT_LOCK_KEY_PREFIX='app_conversation_export'`), `export_max_events=10000`, `export_lock_ttl_seconds=3600` — safeguards but still outputs only internal zip | `openhands/app_server/app_conversation/live_status_app_conversation_service.py:153`, `openhands/app_server/app_conversation/live_status_app_conversation_service.py:283-286`, `openhands/app_server/app_conversation/live_status_app_conversation_service.py:2845-2915` |
| Export/import tools | HTTP route `GET /api/.../export_conversation` streams `open_conversation_export` zip to client — no format negotiation, no `Accept: application/vnd.open-telemetry` etc. | `openhands/app_server/app_conversation/app_conversation_router.py:1619-1639` |

## Answers to Dimension Questions

### 1. Can traces be moved between platforms?

**Partially, but not as traces.** The SDK's observability layer is vendor-locked to Laminar despite reading generic `OTEL_*` env vars. `maybe_init_laminar()` (`_sdk_inspect/sdk/observability/laminar.py:57-112`) checks `OTEL_EXPORTER_OTLP_TRACES_ENDPOINT` etc. via `should_enable_observability()` but then unconditionally imports `lmnr` and registers `LaminarLiteLLMCallback` with `litellm.callbacks`. `RootSpan` (`_sdk_inspect/sdk/observability/laminar.py:231-266`) delegates to `Laminar.start_span`/`Laminar.use_span` and `Laminar.set_trace_session_id`; there is no generic `opentelemetry.trace.Tracer` abstraction or pluggable exporter.

What *is* portable is the **conversation event log**, not the trace: `EventLog` (`_sdk_inspect/sdk/conversation/event_store.py:132-145`) persists each `Event` as Pydantic JSON (`event.model_dump_json`) to a file store, and `export_conversation` (`openhands/app_server/app_conversation/app_conversation_service.py:160-189`) zips those files + `base_state.json`. Rehydration is via `Event.model_validate_json` (`_sdk_inspect/sdk/conversation/event_store.py:105-112`). This is an internal OpenHands format — it can be copied between OpenHands deployments (same Pydantic models, `DiscriminatedUnionMixin` with `kind`), but it is not an OTLP trace, not Honeycomb/Jaeger/Datadog-native, and requires custom code to translate into any external tracing platform. Moving to another tracing backend requires rewriting the `lmnr` integration; setting `OTEL_EXPORTER_OTLP_ENDPOINT` to a non-Laminar collector still goes through Laminar's exporter path.

### 2. Can eval datasets be reused across systems?

**No.** There is no portable eval dataset abstraction. Evaluation is via an in-process `CriticBase` plugin system (`_sdk_inspect/sdk/critic/base.py:65-83`) with four built-in implementations: `AgentFinishedCritic` (checks finish), `APIBasedCritic` (LLM-as-judge via `classify_trace`), `EmptyPatchCritic`, `PassCritic`. They consume a live `list[Event]` trace, not a serialized dataset file. Gating is per-agent `CriticMixin` (`_sdk_inspect/sdk/agent/critic_mixin.py:34-73` — modes `off`/`finish`/`finish_and_message`) and `ResponseDispatch` (`_sdk_inspect/sdk/agent/response_dispatch.py:123-275`). No JSONL/CSV/Parquet, no HF `datasets`, no SWE-Bench harness, no shared schema for tasks/expected patches/metrics that another system could ingest. `Telemetry`/`Metrics` (`_sdk_inspect/sdk/llm/utils/telemetry.py:81-246`) records latency/tokens/cost per `ModelResponse` but does not expose a portable eval result schema.

### 3. Can prompts be migrated?

**As Jinja2 source files, yes — with caveats.** Prompts are plain `.j2` templates rendered by a thin Jinja2 wrapper (`_sdk_inspect/sdk/context/prompts/prompt.py:88-114` — `render_template(prompt_dir, template_name, **ctx)`). Key templates: `system_prompt.j2`, `system_prompt_interactive.j2`, `system_prompt_long_horizon.j2`, `system_prompt_planning.j2`, `self_documentation.j2`, `security_policy.j2`, model-specific overrides `model_specific/anthropic_claude.j2`, `model_specific/google_gemini.j2`, `model_specific/openai_gpt/gpt-5.j2` etc., plus `context/prompts/templates/*` and `context/condenser/prompts/summarizing_prompt.j2`. Because Jinja2 is de-facto standard in Python, copying the `.j2` tree to another system and calling `render_template` with the right `prompt_dir` and context (`security_policy_filename`, `model_family`, `enable_browser`, `llm_security_analyzer`, etc.) reproduces the prompts.

Limitations: (a) no prompt registry/versioning API (no Langfuse/PromptLayer integration, no `prompt_id@version`); templates are versioned only by git; (b) Jinja-specific features (`{% include %}`, `| refine` filter (`_sdk_inspect/sdk/context/prompts/prompt.py:48-73`), `FlexibleFileSystemLoader` absolute-path support) tie migration to Jinja2; Mustache/Handlebars or provider prompt stores would need rewriting; (c) `FileSystemBytecodeCache` under `~/.openhands/cache/jinja` (`_sdk_inspect/sdk/context/prompts/prompt.py:62-66`) is an optimization detail that does not affect portability but signals filesystem coupling; (d) variable contracts are implicit — `agent/prompts/system_prompt.j2:79-88` expects `security_policy_filename`, `llm_security_analyzer`, `enable_browser`, `model_family` etc. without a declared schema, so cross-system reuse risks missing variables.

### 4. Are tool schemas provider-independent?

**Yes — this is the strongest portability story.** The `Schema` base (`_sdk_inspect/sdk/tool/schema.py:173-239`) normalizes JSON Schema `anyOf`/`$ref`/`$defs` via `_process_schema_node` (`_sdk_inspect/sdk/tool/schema.py:70-170`) and produces MCP-compatible schemas, then dynamically reconstructs Pydantic models from MCP schemas. `ToolDefinition` (`_sdk_inspect/sdk/tool/tool.py:379-498`) exposes three orthogonal adapters from a single declaration:

- `to_mcp_tool()` → `{name, description, inputSchema, outputSchema, annotations}`
- `to_openai_tool()` → `ChatCompletionToolParam` (Chat Completions)
- `to_responses_tool()` → `FunctionToolParam` (Responses API, `strict: False`)

Internally, `MessageToolCall` (`_sdk_inspect/sdk/llm/message.py:24-119`) abstracts the duality (`origin: completion|responses`, `id`/`responses_item_id` preservation for prefix-cache, `to_chat_dict()`/`to_responses_dict()`). `Message` (`_sdk_inspect/sdk/llm/message.py:278-722`) further handles `thinking_blocks` (Anthropic), `reasoning_content` (DeepSeek/o1), and `responses_reasoning_item` (OpenAI) with dual `to_chat_dict`/`to_responses_dict` serializers. MCP tools are first-class: `_create_mcp_action_type` (`_sdk_inspect/sdk/mcp/tool.py:120-143`) caches `Schema.from_mcp_schema` per tool name, and `MCPToolDefinition.action_from_arguments` (`_sdk_inspect/sdk/mcp/tool.py:192-222`) validates before dispatch, preventing vendor-specific schema drift from reaching the executor. `NonNativeToolCallingMixin` (`_sdk_inspect/sdk/llm/mixins/non_native_fc.py:34-105`) additionally mocks function calling via prompt engineering for providers without native FC (e.g., `openhands-lm`, `devstral`), with `STOP_WORDS` stop-sequence injection.

Result: adding a new LLM provider requires only that `litellm` supports it; tool definitions themselves need no per-provider fork. The limitation is that extensibility beyond the three fixed targets (MCP, Chat Completions, Responses) requires new methods — there is no generic JSON Schema registry or provider-agnostic function registry beyond these.

## Architectural Decisions

| Decision | Evidence | Consequence |
|----------|----------|-------------|
| **Laminar as sole observability backend, OTEL env vars as feature flag only** | `_sdk_inspect/sdk/observability/laminar.py:25-30` declares `OTEL_*` keys but `maybe_init_laminar()` branches to `lmnr.Laminar` (`_sdk_inspect/sdk/observability/laminar.py:90-112`) | Portable-looking config that is actually single-vendor; switching backends requires code change, not just env-var change. |
| **Per-conversation `RootSpan` owned by `BaseConversation` via named attribute `_observability_root_span`** | `_sdk_inspect/sdk/conversation/base.py:120-151`, `_sdk_inspect/sdk/observability/laminar.py:228`, `_sdk_inspect/sdk/observability/laminar.py:299-330` | Avoids global stack collisions (deprecated `SpanManager`/`start_active_span` in `_sdk_inspect/sdk/observability/laminar.py:352-465`) and survives async/task hops via `Laminar.use_span` re-attachment; but couples span lifecycle to conversation object identity. |
| **Lazy `observe` decorator that bypasses `lmnr` import until enabled** | `_sdk_inspect/sdk/observability/laminar.py:115-196` (`should_enable_observability()` check, `lmnr` import inside `_build_wrapped`) | Zero overhead when tracing disabled, no import-time side effects; downside is no compile-time validation of span attributes. |
| **Pydantic `DiscriminatedUnionMixin` + `kind` for event/tool polymorphism** | `_sdk_inspect/sdk/tool/schema.py:173-198` (strips `kind` on export), `_sdk_inspect/sdk/tool/tool.py:273-302` (serializes `action_type` via `kind_of`), `_sdk_inspect/sdk/event/base.py:20-23` | Events/tools round-trip through JSON by `kind`, enabling zip export/import within OpenHands; fragile across SDK version skew if `kind` registry changes. |
| **Single `Schema` base handles both directions (to/from MCP)** | `_sdk_inspect/sdk/tool/schema.py:25-52` `py_type()`, `_sdk_inspect/sdk/tool/schema.py:70-170` `_process_schema_node`, `_sdk_inspect/sdk/tool/schema.py:201-239` `from_mcp_schema` | Tool schemas stay MCP-canonical while auto-deriving OpenAI parameters; circular-ref fallback to `{"type":"object"}` (`_sdk_inspect/sdk/tool/schema.py:64-67`) loses recursive structure. |
| **Jinja2 filesystem-prompts with `FlexibleFileSystemLoader` + bytecode cache** | `_sdk_inspect/sdk/context/prompts/prompt.py:16-86` | Prompts are ordinary files, easily copied; but loading is filesystem- and Jinja-specific, not a prompt registry. |
| **FileStore abstraction for persistence (local/memory/S3/GCS)** | `_sdk_inspect/sdk/io/base.py:1-30`, `_sdk_inspect/sdk/conversation/event_store.py:132-145` (`Flock` + `write_guard`), `_sdk_inspect/sdk/conversation/persistence_const.py:4-9` | Conversation trajectories are portable as files across storage backends, but format remains internal Pydantic JSON, not an open trace/eval standard. |

## Notable Patterns

- **Triple-target tool adapters from one declaration:** single `ToolDefinition` yields MCP, Chat Completions, and Responses payloads via `_get_tool_schema` + `_prioritize_schema_fields` (`_sdk_inspect/sdk/tool/tool.py:413-498`) — avoids per-provider tool duplication. `_prioritize_schema_fields` moves `security_risk`/`summary` to front to survive token truncation.
- **Transport-agnostic tool-call ID dual-encoding:** `MessageToolCall.id` (canonical) + `responses_item_id` (original `fc_*` prefix) (`_sdk_inspect/sdk/llm/message.py:24-40`) preserved through `to_chat_dict`/`to_responses_dict` so prefix-cache hits survive round-trip and replay stays byte-identical (`_sdk_inspect/sdk/llm/message.py:100-119`).
- **Caching-aware message serialization:** `Message.to_chat_dict` (`_sdk_inspect/sdk/llm/message.py:278-326`) and `_list_serializer`/`_string_serializer` switch on `cache_enabled`/`vision_enabled`/`function_calling_enabled`/`force_string_serializer`, with `cache_control: {type: ephemeral}` promotion for `tool` role (`_sdk_inspect/sdk/llm/message.py:365-382`).
- **Non-native function-calling mock:** `NonNativeToolCallingMixin` (`_sdk_inspect/sdk/llm/mixins/non_native_fc.py:34-105`) converts FC messages ↔ prompt-mocked messages and injects `STOP_WORDS` when `get_features(model).supports_stop_words` — allows tool use on providers without native FC.
- **Event log as append-only JSON with distributed locking:** `EventLog` (`_sdk_inspect/sdk/conversation/event_store.py:132-145`) + Redis export lock (`openhands/app_server/app_conversation/live_status_app_conversation_service.py:2900-2915`) provides process-safe trajectory persistence.

## Tradeoffs

| Tradeoff | Pro | Con |
|----------|-----|-----|
| Laminar-only tracing | `lmnr` + `LaminarLiteLLMCallback` gives rich LLM+tool spans with minimal code; `RootSpan` trick fixes ~60% orphan-trace bug (`_sdk_inspect/sdk/observability/laminar.py:242-250`) | Cannot switch to Jaeger/Honeycomb/Datadog without code fork; `OTEL_*` env vars are misleadingly generic. No OTLP file export for offline analysis. |
| Pydantic JSON event log vs OTel spans | Strong typing, `DiscriminatedUnionMixin` ensures round-trip correctness, `model_validate_json` is strict; trivial to zip/export | Not queryable in trace UIs, not convertible without custom mapping; version skew breaks deserialization (no schema registry). |
| Jinja2 static files vs prompt registry | No service dependency, works offline, `FileSystemBytecodeCache` avoids reparsing | No versioning, A/B, or no-code editing; variable contracts undocumented; migration to non-Jinja runtimes requires template rewrite. |
| MCP-canonical tool schemas | Single source of truth, auto-handles `anyOf`/`$ref`/circular refs (`_process_schema_node`) | Circular refs collapse to generic `object` (`_sdk_inspect/sdk/tool/schema.py:64-67`) — tree/graph tools lose structure in MCP view. |
| Triple adapters vs generic converter | Explicit, tested paths for top 3 targets (MCP/Chat/Responses) | Adding Anthropic/Google-native or new OpenAI APIs requires new method; no plugin point for providers. |
| Zip trajectory export | Simple, includes `base_state.json` + per-event JSON + Redis lock + size guard (`export_max_events=10000`) | Download-only, no import path, no format negotiation, no incremental export, no cross-platform trace conversion. |

## Failure Modes / Edge Cases

- **Trace context loss across async boundaries:** documented prior use of `Laminar.start_active_span` lost parent for ~60% of conversations (`_sdk_inspect/sdk/observability/laminar.py:244-250`); fix is `RootSpan` + `Laminar.use_span` re-attachment at every `@observe` entry (`_sdk_inspect/sdk/observability/laminar.py:299-330`). Failure returns to pass-through (`should_enable_observability()==False`) rather than crashing, but traces silently disappear if env vars are mis-configured.
- **Stale EventLog index:** `_get_single_item` rebuilds from disk on `KeyError` (`_sdk_inspect/sdk/conversation/event_store.py:80-95`) with warning `Stale EventLog index at %d; rebuilding` — NFS/`LocalFileStore.flock` unreliability noted (`_sdk_inspect/sdk/conversation/event_store.py:18-25`) means concurrent writes can still race despite `LOCK_TIMEOUT_SECONDS=30`.
- **Circular schema degradation:** `_process_schema_node` returns `_shallow_expand_circular_ref` generic `{"type":"object"}` (`_sdk_inspect/sdk/tool/schema.py:64-67`, `_sdk_inspect/sdk/tool/schema.py:114-120`) — recursive tools (e.g., tree editors) silently lose type fidelity in `to_mcp_schema` output.
- **Export lock unavailable:** `open_conversation_export` (`openhands/app_server/app_conversation/live_status_app_conversation_service.py:2903-2912`) raises `ConversationExportAlreadyRunning` if `export_lock_required==True` else logs `lock_unavailable_proceeding_without_lock` and proceeds without lock — risk of torn zip or duplicate concurrent exports.
- **Export size guard:** `export_max_events=10000` (`openhands/app_server/app_conversation/live_status_app_conversation_service.py:284-285`); `ConversationExportTooLarge` if `event_count > export_max_events` (`openhands/app_server/app_conversation/live_status_app_conversation_service.py:2845-2853`) — large trajectories cannot be exported without raising the limit (no pagination/streaming of filtered export).
- **Deprecated span globals:** `SpanManager`, `start_active_span`/`end_active_span` remain as deprecated shims until 1.27.0 (`_sdk_inspect/sdk/observability/laminar.py:352-465`, `warn_deprecated` inside) — risk of mixed usage if external consumers still call globals, leading to orphan spans on the LIFO stack.
- **Image/vision toggle divergence:** `_list_serializer` drops `ImageContent` when `vision_enabled==False` (`_sdk_inspect/sdk/llm/message.py:373-378`) and `to_responses_dict` does similar (`_sdk_inspect/sdk/llm/message.py:470-474`) — traces rendered without vision lose image payloads, making cross-system replay non-equivalent.
- **`force_string_serializer` provider quirks:** `Message.to_chat_dict` (`_sdk_inspect/sdk/llm/message.py:300-307`) forces string content for HuggingFace/Groq-like providers — a trace recorded with list serializer will not replay identically on those providers without re-serialization.

## Future Considerations

- **Extract observability behind an interface:** replace direct `lmnr` imports with a `Tracer` protocol (`start_span`, `use_span`, `set_attribute`) and Pluggable OTel exporter; then `maybe_init_laminar` becomes `maybe_init_tracer(provider="laminar"|"otel"|"datadog")` — makes `OTEL_*` env vars truthful.
- **Add OTLP/Parquet export:** alongside zip, support `export?format=otlp|parquet|jsonl` that maps `Event` → `opentelemetry.proto.trace.v1.Span` or Arrow, enabling move to Jaeger/Honeycomb without custom code.
- **Introduce import path:** `import_trajectory(zip_or_otlp)` that calls `Event.model_validate_json` per file and validates `kind` registry version — needed for migration and cross-deployment replay.
- **Prompt registry:** promote `.j2` tree to versioned store with `prompt_id`, `version`, `variables_json_schema` and compatibility test (render + snapshot) — allows marketplace/discovery and safe cross-system migration; consider adding Mustache/Handlebars adapter or standardizing on Jinja2 explicitly.
- **Portable eval harness:** define `EvalDataset` (tasks + expected outputs + scorers) with `load_jsonl`/`to_hf_dataset` and a `run_eval(dataset, critic)` CLI — makes `CriticBase` results comparable across platforms; align scorer output with `Telemetry.metrics` snapshot.
- **Schema versioning:** embed `schema_version` in `Event`/`Schema` payloads and implement up-migration in `_handle_deprecated_fields` (`_sdk_inspect/sdk/llm/message.py:247-261`) — prevents silent breakage when moving trajectories between SDK versions.

## Questions / Gaps

- No evidence found for cross-platform trace migration tooling — searched `_sdk_inspect/sdk/observability/*:1-503`, `openhands/app_server/**/*:1-195`, and codebase grep for `otel`/`jaeger`/`zipkin`/`honeycomb`; only Laminar + raw OTEL env reading exists.
- No evidence found for eval dataset portability or external benchmark harness — grep for `dataset`/`swe-bench`/`benchmark` in `_sdk_inspect/sdk` yields only critic/mixing code; no `evaluation/` directory.
- No evidence found for prompt registry/versioning beyond git-tracked `.j2` files — `prompt.py:57-114` is the sole loader, with no DB or API for prompt CRUD.
- No evidence found for export/import of traces in vendor-neutral format — `app_conversation_service.py:160-189` and `live_status_app_conversation_service.py:2833-2915` are the only export paths, both producing internal zip.
- Open: does `enterprise/` SaaS layer add additional trace storage/analytics beyond OSS `Telemetry` file logs? `enterprise/server/sharing/*` suggests shared event services but not tracing-platform portability.

---

Generated by `19.02-portable-trace-eval-and-prompt-schemas` against `openhands`.
