# Source Analysis: openai-agents-sdk

## Portable Trace, Eval, and Prompt Schemas

### Source Info

| Field | Value |
|-------|-------|
| Name | openai-agents-sdk |
| Path | `studies/agent-harness-study/sources/openai-agents-sdk` |
| Language / Stack | Python 3.10+ / openai>=2.36.0, pydantic>=2.12.2, httpx2 |
| Analyzed | 2026-08-28 |

## Summary

`openai-agents-sdk` implements a **proprietary but pluggable** tracing subsystem. Traces and spans use a custom SDK dict format (`object: trace` / `object: trace.span` with nested `span_data`) and are dispatched through an abstract `TracingProcessor`/`TraceProvider` interface that supports full replacement (`set_trace_processors`/`set_trace_provider`) and fan-out (`SynchronousMultiTracingProcessor`). The default `BackendSpanExporter` is hardcoded to OpenAI's `https://api.openai.com/v1/traces/ingest` and applies OpenAI-specific sanitization/truncation; no OpenTelemetry, OTLP, or generic export format is shipped. Prompt handling is a thin wrapper around OpenAI's hosted `ResponsePromptParam` (`id`/`version`/`variables`) with no template language (no Jinja/Mustache/Handlebars). Tool schemas are standard JSON Schema via Pydantic but are aggressively coerced to OpenAI's `strict` mode (`additionalProperties: false`, all props required, `$ref` inlining) via `ensure_strict_json_schema`. No eval dataset, benchmark, or cross-platform migration/import tooling exists in the source. Overall, portability is achieved only by writing custom processors/schemas; out-of-box formats are OpenAI-coupled.

## Rating

**3 / 10 — Absent to ad-hoc portability (low end of 1-3 band).**

Rationale: Tracing exposes a clean provider abstraction and `export()`/`to_json()` serialization (`src/agents/tracing/traces.py:152`, `src/agents/tracing/spans.py:396`) and a replaceable exporter endpoint (`src/agents/tracing/processors.py:62`), demonstrating intentional extensibility. However the wire format is SDK-proprietary, the default exporter is OpenAI-locked, no OTel/OTLP adapter or standard export is included, prompts are bound to `openai.types.responses.response_prompt_param` (`src/agents/prompts.py:8`), unsupported prompts are rejected on non-Responses models (`src/agents/models/openai_chatcompletions.py:85-102`), tool schemas are forced through `ensure_strict_json_schema` (`src/agents/strict_schema.py:115`, `src/agents/function_schema.py:477`), and no eval-dataset or migration tooling is present. Migration requires rewriting data or writing custom code.

## Evidence Collected

Every entry MUST include a file path with line numbers. Format: `path/to/file.ts:NN`.

| Area | Evidence | File:Line |
|------|----------|-----------|
| Trace format abstraction | `Trace` abstract interface defining `trace_id`, `name`, `export()`, `to_json()` with per-trace `tracing_api_key` | `src/agents/tracing/traces.py:51-185` |
| Trace format abstraction | `Span[TSpanData]` abstract interface with `trace_id`, `span_id`, `parent_id`, `span_data`, `export()`, `trace_metadata` | `src/agents/tracing/spans.py:59-214` |
| Trace format abstraction | `TracingProcessor` abstract with `on_trace_start/end`, `on_span_start/end`, `shutdown`, `force_flush`; `TracingExporter.export(items: list[Trace|Span])` | `src/agents/tracing/processor_interface.py:9-142` |
| Trace format abstraction | `TraceProvider` abstract: `create_trace`, `create_span`, `register_processor`, `set_processors`, `gen_trace_id/span_id/group_id` | `src/agents/tracing/provider.py:222-298` |
| Trace format abstraction | `DefaultTraceProvider` generates `trace_<hex>` / `span_<hex>` / `group_<hex>` IDs, honors `OPENAI_AGENTS_DISABLE_TRACING`, delegates to `SynchronousMultiTracingProcessor` | `src/agents/tracing/provider.py:300-511` |
| Trace format abstraction | Global provider indirection: `get_trace_provider()`, `set_trace_provider()`, lazy `default_processor()`/`default_exporter()` with `atexit` shutdown | `src/agents/tracing/setup.py:27-66` |
| Trace format abstraction | Public surface re-exports `Trace`, `TracingProcessor`, `TraceProvider`, `add_trace_processor`, `set_trace_processors`, `set_trace_provider`, `flush_traces` | `src/agents/tracing/__init__.py:1-130` |
| Trace serialization | `TraceImpl.export()` returns `{"object":"trace","id":trace_id,"workflow_name":name,"group_id":..., "metadata":...}` | `src/agents/tracing/traces.py:568-575` |
| Trace serialization | `SpanImpl.export()` returns `{"object":"trace.span","id":span_id,"trace_id":...,"parent_id":..., "span_data": span_data.export(), "error":...}` with whitelisted `metadata` routing (`agent_harness_id`) | `src/agents/tracing/spans.py:396-423` |
| Trace state portability | `TraceState` dataclass with `to_json`/`from_json` mapping `id`/`trace_id`, `workflow_name`, `group_id`, `metadata`, `tracing_api_key_hash` for `RunState` persistence and `reattach_trace` without re-emitting start | `src/agents/tracing/traces.py:195-404` |
| Span data types | 13 SDK-specific `SpanData` subtypes (`AgentSpanData`, `GenerationSpanData`, `FunctionSpanData`, `ResponseSpanData`, etc.) each with custom `type` and `export()` dict — no OTel `SpanKind`/`SpanContext` mapping | `src/agents/tracing/span_data.py:28-452` |
| Default ingest coupling | `BackendSpanExporter._OPENAI_TRACING_INGEST_ENDPOINT = "https://api.openai.com/v1/traces/ingest"`; headers include `OpenAI-Beta: traces=v1`, `Authorization: Bearer <api_key>` | `src/agents/tracing/processors.py:44-50,138-157` |
| OpenAI sanitization | Endpoint-aware sanitization: truncates `input`/`output` to 100k bytes, strips `usage` unless `span_data.type=="generation"`, normalizes allowed usage keys | `src/agents/tracing/processors.py:258-319,322-531` |
| Endpoint configurability | `BackendSpanExporter.__init__(..., endpoint=..., api_key, organization, project, max_retries, base_delay, max_delay)` allows non-OpenAI endpoint, but sanitization is still conditional on matching default endpoint | `src/agents/tracing/processors.py:57-90` |
| Batch processor | `BatchTraceProcessor` queues via `queue.Queue(maxsize=8192)` with background thread, `schedule_delay=5.0`, `export_trigger_ratio=0.7`, `force_flush`/`shutdown(timeout)` semantics | `src/agents/tracing/processors.py:541-718` |
| Ecosystem note | Docs list 27+ third-party processors (W&B, Phoenix, MLflow, Braintrust, LangSmith, Langfuse, Datadog etc.) as external adapters — confirms portability via custom processors, not via standard format | `docs/tracing.md:196-226` |
| Prompt template portability | `Prompt` TypedDict is `{"id": str, "version"?: str, "variables"?: dict[str, ResponsesPromptVariables]}` directly wrapping `openai.types.responses.response_prompt_param.ResponsePromptParam` | `src/agents/prompts.py:23-34` |
| Prompt template portability | `PromptUtil.to_model_input` resolves static `Prompt` dict or `DynamicPromptFunction` and returns `ResponsePromptParam` — no Jinja/handlebars/mustache engine | `src/agents/prompts.py:56-82` |
| Prompt template portability | `Agent.get_prompt()` delegates to `PromptUtil.to_model_input` returning `ResponsePromptParam | None` | `src/agents/agent.py:1073-1083` |
| Prompt model coupling | `OpenAIChatCompletionsModel._handle_unsupported_prompt` raises `UserError` when `strict_feature_validation` else warns: "Reusable prompts are only supported by the Responses API" — Chat Completions path rejects hosted prompts | `src/agents/models/openai_chatcompletions.py:85-102` |
| Prompt cache key | `prompt_cache_key.py` references OpenAI `prompt` in cache key computation, reinforcing Responses-only semantics | `src/agents/run_internal/prompt_cache_key.py:1-...` (checked) |
| Eval dataset portability | No `eval`/`dataset`/`benchmark` module in `src/agents`; `grep` for `eval|dataset|Dataset` returns unrelated hits (`unevaluatedItems`, `evaluate_needs_approval` etc.) | `src/agents/*` (no evidence found) |
| Tool schema portability | `FuncSchema` holds `params_pydantic_model`, `params_json_schema: dict`, `strict_json_schema: bool=True`, `to_call_args` | `src/agents/function_schema.py:22-45` |
| Tool schema portability | `function_schema(..., strict_json_schema=True)` calls `ensure_strict_json_schema(json_schema)` after Pydantic `json_schema()` | `src/agents/function_schema.py:296-489` |
| Tool schema strictness | `FunctionTool` stores `params_json_schema`, `strict_json_schema`; constructor copies and calls `ensure_strict_json_schema(_copy_json_schema(schema))`; `_normalize_function_tool_output_json_schema` enforces same | `src/agents/tool.py:452-599,2751-2762` |
| Tool schema strictness | `ensure_strict_json_schema` mutates schema to enforce `additionalProperties: false`, `required: all props`, `$ref` expansion, depth/node budgets, `allOf`/`oneOf` normalization, rejecting open objects | `src/agents/strict_schema.py:115-405` |
| Cross-provider tool handling | MCP tool conversion copies then `ensure_strict_json_schema` unless `strict=False`, preserving OpenAI strictness for remote tools | `src/agents/mcp/util.py:538-577` |
| Tool schema providers | `LiteLLMModel`/`AnyLLMModel` under `src/agents/extensions/models/` bridge to third-party providers via `Converter` but inherit strict schemas; no generic schema-unwrapping | `src/agents/extensions/models/litellm_model.py:1-...`, `src/agents/extensions/models/any_llm_model.py:1-...` |
| Export/import tools | No CLI or library for trace/dataset/prompt migration; only programmatic `Trace.export()`, `Span.export()`, `TraceState.to_json()`, and processor-level `export()` batching | `src/agents/tracing/traces.py:152,171`, `src/agents/tracing/spans.py:181,396` |
| Tests confirm coupling | `test_agent_prompt.py` exercises `Prompt = {"id": "my_prompt", "version":"1", "variables":{...}}` against `FakeModel` expecting `ResponsePromptParam` passthrough | `tests/test_agent_prompt.py:52-60` |
| Tests confirm tracing format | `testing_processor.py` normalizes exported traces asserting `object == "trace"` and `object == "trace.span"`, stripping `id` prefixes `trace_*`/`span_*` | `tests/testing_processor.py:96-145` |

## Answers to Dimension Questions

**1. Can traces be moved between platforms?**
Partially — via code, not via data. The SDK defines a provider-agnostic processor abstraction (`TracingProcessor` at `src/agents/tracing/processor_interface.py:9`, `TraceProvider` at `src/agents/tracing/provider.py:222`, `set_trace_provider`/`set_trace_processors`/`add_trace_processor` at `src/agents/tracing/__init__.py:94-105` and `src/agents/tracing/setup.py:27`). You can replace the default OpenAI exporter with any custom processor. The docs explicitly list ~27 community processors (`docs/tracing.md:198-226`) as the intended integration path. However the underlying serialization is proprietary (`TraceImpl.export()` at `src/agents/tracing/traces.py:568` → `{"object":"trace",...}`, `SpanImpl.export()` at `src/agents/tracing/spans.py:396` → `{"object":"trace.span",...}`) with custom `span_data.type` values (`src/agents/tracing/span_data.py:28-452`). There is no OTel/OTLP adapter, no `opentelemetry` dependency in `pyproject.toml:9-18`, and the default `BackendSpanExporter` is hardcoded to `https://api.openai.com/v1/traces/ingest` (`src/agents/tracing/processors.py:45`) with OpenAI-specific header `OpenAI-Beta: traces=v1`. Batch export (`BatchTraceProcessor` at `src/agents/tracing/processors.py:541`) handles fan-out but format translation must be written by the consumer.

**2. Can eval datasets be reused across systems?**
No evidence found. The source tree contains no `eval`, `evaluation`, `dataset`, or `benchmark` abstractions. `grep` for `eval|dataset` in `src/agents` yields only `unevaluatedItems` (`src/agents/extensions/tool_output_trimmer.py:62`) and policy-evaluation helpers (`evaluate_needs_approval_setting` etc.). No dataset schema, loader, versioning, or portability spec was located. The question is unanswerable from implementation — treated as absent.

**3. Can prompts be migrated?**
No — prompts are OpenAI-hosted references, not portable templates. `src/agents/prompts.py:23` defines `Prompt` as `id`/`version`/`variables` wrapping `openai.types.responses.response_prompt_param.ResponsePromptParam` (`src/agents/prompts.py:8`). `PromptUtil.to_model_input` (`src/agents/prompts.py:56`) passes this through to `Model.get_response(prompt= ResponsePromptParam)` (`src/agents/models/interface.py:80`). There is no Jinja, Mustache, Handlebars, or file-based prompt template engine. Dynamic prompts are Python callables returning a `Prompt` dict (`src/agents/prompts.py:47`), still resolved to the same OpenAI type. `OpenAIChatCompletionsModel._handle_unsupported_prompt` (`src/agents/models/openai_chatcompletions.py:85-102`) explicitly rejects prompts for non-Responses paths, confirming lack of cross-provider portability.

**4. Are tool schemas provider-independent?**
Nominally JSON Schema, but in practice OpenAI-coupled. Tool parameter schemas are generated from Python type hints via Pydantic (`src/agents/function_schema.py:22-45`) and then forced through `ensure_strict_json_schema` (`src/agents/strict_schema.py:115`) which enforces `additionalProperties: false`, injects `required` for all properties, inlines `$ref`, and normalizes `oneOf`→`anyOf`/`allOf`. `FunctionTool` (`src/agents/tool.py:452-599`) defaults `strict_json_schema=True` and mutates schemas in-place. MCP tools receive the same treatment (`src/agents/mcp/util.py:538-577`). Cross-provider bridges (`src/agents/extensions/models/litellm_model.py`, `src/agents/extensions/models/any_llm_model.py`) must accept these strict schemas; setting `strict_json_schema=False` opts out (`AgentOutputSchema` at `src/agents/agent_output.py:85-118`, `tool.py:2557-2606`). Hosted tools (`WebSearchTool`, `FileSearchTool`, `HostedMCPTool` etc. in `src/agents/tool.py`) have no portable equivalent. So migration requires either disabling strict mode or rewriting hosted tool definitions.

## Architectural Decisions

- **Abstract processor/provider seam over standard format** — The SDK chooses pluggability (`TracingProcessor`, `TraceProvider`, `SynchronousMultiTracingProcessor` at `src/agents/tracing/provider.py:93`) over conformance to OTel. Trade: trivial to add a vendor exporter but no zero-effort data interoperability.
- **OpenAI ingest as default path** — `BackendSpanExporter` (`src/agents/tracing/processors.py:44`) hardcodes the OpenAI endpoint, API-key grouping, retries/backoff, and endpoint-conditional sanitization (`_should_sanitize_for_openai_tracing_api` at `src/agents/tracing/processors.py:258`). Non-OpenAI backends must supply a custom exporter or override `endpoint`.
- **Exported dict as wire schema** — Spans/traces serialize to plain dicts with `object` discriminator (`src/agents/tracing/traces.py:568-575`, `src/agents/tracing/spans.py:396-423`) rather than protobuf/OTLP or OpenTelemetry `ReadableSpan`. Reduces dependencies but cedes ecosystem tooling.
- **Strict JSON Schema as invariant** — `ensure_strict_json_schema` (`src/agents/strict_schema.py:115-405`) with depth/node budgets (`_MAX_SCHEMA_DEPTH=100`, `_MAX_SCHEMA_NODES=100000`) guarantees Responses API compatibility at the cost of portability and open-object support. A `reject_open_objects` mode exists but defaults off.
- **Hosted prompt reference over client-side templating** — Prompts are stored server-side and referenced by `id` (`src/agents/prompts.py:23-34`), aligning with OpenAI's Prompt management feature rather than a local template file. No local `prompt.md` loader beyond sandbox memory prompts which are internal.

## Notable Patterns

- **Pluggable pipeline**: `DefaultTraceProvider` → `SynchronousMultiTracingProcessor` (fan-out with per-processor error isolation at `src/agents/tracing/provider.py:117-175`) → `BatchTraceProcessor` (async queue with `schedule_delay`/`export_trigger_ratio` at `src/agents/tracing/processors.py:548-669`) → `BackendSpanExporter` (grouping by `tracing_api_key` at `src/agents/tracing/processors.py:125-128`). Duplicated in docs as the recommended customization point (`docs/tracing.md:144-152`).
- **ContextVar scoping**: Traces/spans stored in `Scope` via `contextvars` (`src/agents/tracing/traces.py:305-335`, `src/agents/tracing/spans.py:289-407`) enabling concurrent nesting without explicit propagation.
- **Reattach without re-emit**: `ReattachedTrace` (`src/agents/tracing/traces.py:305-388`) and `TraceState` (`src/agents/tracing/traces.py:195-277`) rebuild a live context from persisted JSON without calling `on_trace_start`, preserving `trace_id` across `RunState` resumes while avoiding duplicate trace creation. Fingerprinting `tracing_api_key_hash` avoids persisting secrets.
- **Conditional sanitization**: Payload sanitization gated by endpoint match (`_should_sanitize_for_openai_tracing_api` at `src/agents/tracing/processors.py:258`) so custom endpoints receive full `span_data` including arbitrary `input`/`output`/`usage`.

## Tradeoffs

- **Extensibility vs. interoperability** — The processor interface makes adding a vendor exporter trivial (docs list 27 adapters), but because the payload is SDK-defined, each adapter must translate from the custom dict rather than consuming OTLP directly. Teams migrating from/to LangSmith/Langfuse etc. rewrite the translation layer.
- **Strict mode safety vs. expressiveness** — Strict schemas catch missing `required`/`additionalProperties` violations early for OpenAI, but reject legitimate JSON Schema patterns (e.g., `additionalProperties: {}`, nullable unions without `object` root) and require `strict_json_schema=False` for portable schemas. Error messages are explicit (`src/agents/strict_schema.py:27-50`).
- **Background batching vs. durability** — `BatchTraceProcessor` reduces ingest QPS and tolerates transient failures (exponential backoff at `src/agents/tracing/processors.py:158-221`), but queues are in-memory with hard cap (`max_queue_size=8192`, drops on `queue.Full` at `src/agents/tracing/processors.py:603`). No disk spill; shutdown/flush must be called for short-lived workers (`docs/tracing.md:48-96`).
- **Hosted prompts vs. local portability** — Server-side prompts enable central governance and `variables` substitution without shipping template files, but create hard vendor lock-in; non-Responses models cannot use them at all.
- **No eval/dataset layer** — Keeping evaluation out of scope simplifies the SDK, but leaves teams to invent dataset schemas and serialization for offline eval portability.

## Failure Modes / Edge Cases

- **Queue overflow drops silently**: `BatchTraceProcessor.on_trace_start`/`on_span_end` (`src/agents/tracing/processors.py:597-621`) catches `queue.Full` and logs a single warning — excess traces/spans are dropped without backpressure or persistence.
- **OTel absence**: No `opentelemetry` import or OTLP exporter; attempts to forward via standard collector require a custom `TracingExporter` that re-serializes the SDK dict.
- **Endpoint-locked truncation**: Sanitization truncates `input`/`output` to 100k JSON bytes per field and strips non-`generation` `usage` (`src/agents/tracing/processors.py:46-53,281-291`); workloads with large tool payloads or custom `usage` lose data on the default endpoint without indication beyond truncated preview (`_truncated_preview` at `src/agents/tracing/processors.py:432-446`).
- **Tracing disabled by env or ZDR**: `OPENAI_AGENTS_DISABLE_TRACING=1` (`src/agents/tracing/provider.py:346-352`) or org ZDR policy makes all `trace()` calls return `NoOpTrace` (`src/agents/tracing/traces.py:407`) and `export()`→`None`; migration tooling cannot recover traces that were never captured.
- **Prompt rejection on cross-provider**: Moving an agent with `prompt= {"id": ...}` to `OpenAIChatCompletionsModel` raises `UserError` under `strict_feature_validation` (`src/agents/models/openai_chatcompletions.py:94-95`); otherwise it logs and silently ignores the prompt, causing behavior divergence.
- **Strict schema denial-of-service guards**: Malicious or deeply nested MCP schemas trip `_MAX_SCHEMA_DEPTH=100` or `_MAX_SCHEMA_NODES=100000` (`src/agents/strict_schema.py:21-25`) with `UserError`; legitimate large schemas from third-party tools may need `strict_json_schema=False` or schema simplification.
- **API key grouping failure**: `BackendSpanExporter.export` (`src/agents/tracing/processors.py:125-134`) groups batches by `tracing_api_key` and skips the group if `api_key` is `None` with warning `"OPENAI_API_KEY is not set, skipping trace export"` — traces silently never leave the process.
- **GeneratorExit context mismatch**: `ReattachedTrace`/`TraceImpl` finalization (`src/agents/tracing/traces.py:18-48`, `src/agents/tracing/spans.py:31-56`) tolerates `ValueError` on `ContextVar.reset` when finalization happens on a different task, leaving stale trace context current until the owning task's scope ends.

## Future Considerations

- Add an **OpenTelemetry / OTLP exporter** as a first-class `TracingExporter` (or contrib package) that maps `TraceImpl`/`SpanImpl` fields to OTel `SpanContext`, `SpanKind`, and `Resource`/`InstrumentationScope`, enabling zero-translation migration to any OTel backend and reducing the 27 bespoke adapters.
- Publish a **JSON Schema for the export payload** (versioned `object: trace` / `object: trace.span` envelope and each `SpanData.type` variant) and support `export(format="otlp"|"jsonl")` plus a matching `import_traces(path)` for data migration/backfill.
- Introduce **portable prompt templates** (e.g., `PromptTemplate` with Jinja2/Mustache, file-based `prompt.md`, or `Prompt({messages:[...]})`) alongside hosted `id` prompts, with a migration helper that expands `ResponsePromptParam` into message lists for non-Responses providers (removing the `UserError` at `src/agents/models/openai_chatcompletions.py:85`).
- Offer **strict-schema opt-in per tool** with a `json_schema_compat` flag that emits vanilla JSON Schema (`additionalProperties` preserved, `$ref` unexpanded) for cross-provider use, and document the mapping in `docs/tools.md`.
- Provide **eval dataset scaffolding** (`datasets` module with `Dataset`, `Example`, `EvalRun` types and `to_jsonl`/`from_jsonl`) so evals built on this SDK can be re-run on other harnesses without re-authoring.
- Add **disk-backed spill or synchronous fallback** for `BatchTraceProcessor` queue overflow (or expose `on_drop` hook) to avoid silent loss in bursty workloads and to support cold-start migration exports.

## Questions / Gaps

- No source evidence for evaluative dataset portability — search across `src/agents` and `docs` found no dataset/eval loader; whether this is intentionally out-of-scope vs. deferred is undocumented.
- The SDK advertises third-party processors (`docs/tracing.md:196-226`) but provides no **conformance test suite** that validates a custom processor receives the full export dict correctly — portability claim rests on external implementations.
- `TraceState.tracing_api_key_hash` (`src/agents/tracing/traces.py:187-192,270-273`) is persisted as a SHA-256 fingerprint; no migration guide explains how to rotate or re-key traces during import across providers.
- Prompt `variables` type (`ResponsesPromptVariables` from `openai.types.responses.response_prompt_param` at `src/agents/prompts.py:9`) is opaque to the SDK — variable schema evolution and cross-provider substitution semantics are unspecified.
- No **version negotiation** for the export format: `TraceImpl.export()` has no `version` field, so consumers cannot detect breaking changes to the dict shape during platform migration.

---
Generated by `Dimension 19.02: Portable Trace, Eval, and Prompt Schemas` against `openai-agents-sdk`.
