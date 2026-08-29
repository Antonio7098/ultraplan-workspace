# Source Analysis: pydantic-ai

## Portable Trace, Eval, and Prompt Schemas

### Source Info

| Field | Value |
|-------|-------|
| Name | pydantic-ai |
| Path | `studies/agent-harness-study/sources/pydantic-ai` |
| Language / Stack | Python (Pydantic AI Slim + Pydantic Evals + Pydantic Graph) |
| Analyzed | 2026-08-28 |

## Summary

Pydantic AI's portability story is split: traces and tool/message schemas are built on open standards (OpenTelemetry GenAI semantic conventions and JSON Schema) and are explicitly provider-agnostic, while prompts and eval datasets are portable *within* the Pydantic ecosystem (YAML/JSON + Pydantic validation, Handlebars templates) but have no cross-platform registry or migration tooling. `ModelMessagesTypeAdapter` (`pydantic_ai_slim/pydantic_ai/messages.py:2769`) and `Dataset.to_file/from_file` (`pydantic_evals/pydantic_evals/dataset.py:747`, `pydantic_evals/pydantic_evals/dataset.py:556`) provide round-trippable serialization for messages and eval cases, and `_otel_messages.py` (`pydantic_ai_slim/pydantic_ai/_otel_messages.py:1`) is constrained to spec-only types per repo policy (`pydantic_ai_slim/pydantic_ai/AGENTS.md:16`). Tool definitions emit a single canonical JSON Schema via `_function_schema.function_schema` (`pydantic_ai_slim/pydantic_ai/_function_schema.py:110`) and per-provider `JsonSchemaTransformer` rewrites (`pydantic_ai_slim/pydantic_ai/_json_schema.py:59`) handle strict-mode divergence at the edge. Prompt templating is Handlebars via optional `pydantic-handlebars` (`pydantic_ai_slim/pydantic_ai/template.py:61`) and system prompts/instructions are plain strings or Python callables (`pydantic_ai_slim/pydantic_ai/_system_prompt.py:14`), so migration requires the same engine and Python context. No dedicated trace/eval migration CLI exists; portability relies on OTLP-forwardable OTel spans (`pydantic_ai_slim/pydantic_ai/models/instrumented.py:253`) and generic YAML/JSON files.

## Rating

**6 / 10** — Present but inconsistent: traces, message history, and tool schemas have clear, tested, spec-aligned abstractions with explicit versioning and backward-compat handling; eval datasets and prompts are serializable to standard file formats but are framework-coupled (evaluators remain Python classes loaded via a registry, templates require Handlebars + `deps_type`, instructions are code-adjacent). Missing migration tooling and prompt registry limit cross-platform moves without rewriting.

## Evidence Collected

Every entry MUST include a file path with line numbers. Format: `path/to/file.ts:NN`.

| Area | Evidence | File:Line |
|------|----------|-----------|
| Trace format abstraction | `_otel_messages.py` TypedDicts (`TextPart`, `ToolCallPart`, `MediaUrlPart`, `UriPart`/`BlobPart` v4+, `ThinkingPart`) header cites the OTel GenAI non-normative examples and states spec alignment | `pydantic_ai_slim/pydantic_ai/_otel_messages.py:1-4`, `pydantic_ai_slim/pydantic_ai/_otel_messages.py:53-97` |
| Trace format abstraction | Policy: `In _otel_*.py modules, implement only spec-defined features — no custom additions` prevents spec drift | `pydantic_ai_slim/pydantic_ai/AGENTS.md:16`, `pydantic_ai_slim/pydantic_ai/CLAUDE.md:16` |
| Trace format abstraction | `InstrumentationSettings.messages_to_otel_messages` normalizes `ModelRequest`/`ModelResponse` into `list[ChatMessage]` (`role` via `_otel_message_role`) for OTel attributes `gen_ai.input.messages`/`gen_ai.output.messages` | `pydantic_ai_slim/pydantic_ai/models/instrumented.py:205-223`, `pydantic_ai_slim/pydantic_ai/models/instrumented.py:410-436` |
| Trace format abstraction | `handle_messages` serializes via `safe_to_json` to `gen_ai.input.messages`, `gen_ai.output.messages`, `gen_ai.system_instructions` and `logfire.json_schema` | `pydantic_ai_slim/pydantic_ai/models/instrumented.py:253-289` |
| Trace format abstraction | Versioned instrumentation (`version: Literal[2,3,4,5,6]`, default 5) with deprecation warning for 2/3/4, and v4 `type='uri'/'blob'` + `modality` vs legacy `image-url`/`binary` | `pydantic_ai_slim/pydantic_ai/models/instrumented.py:78-79`, `pydantic_ai_slim/pydantic_ai/models/instrumented.py:154-164`, `pydantic_ai_slim/pydantic_ai/_otel_messages.py:53-97`, `pydantic_ai_slim/pydantic_ai/messages.py:1113-1134` |
| Trace format abstraction | `provider_attributes`/`model_attributes` emit both `gen_ai.provider.name` (new) and `gen_ai.system` (deprecated) for backward compat | `pydantic_ai_slim/pydantic_ai/_instrumentation.py:284-303` |
| Trace format abstraction | `build_tool_definitions` emits OTel GenAI `gen_ai.tool.definitions` (`type='function'`, `name`, `description`, `parameters`) always, independent of `include_model_request_parameters` | `pydantic_ai_slim/pydantic_ai/_instrumentation.py:370-395` |
| Trace format abstraction | `redact_binary_content`/`_convert_binary_to_otel_part` respect `include_content`/`include_binary_content` flags; `MessageJsonCache` keeps per-request history serialization O(new) not O(history) | `pydantic_ai_slim/pydantic_ai/_instrumentation.py:141-168`, `pydantic_ai_slim/pydantic_ai/messages.py:1064-1080`, `pydantic_ai_slim/pydantic_ai/models/instrumented.py:225-251` |
| Dataset portability | `Dataset` generic `Case[InputsT,OutputT,MetadataT]` with `from_file`/`from_text`/`from_dict` and `to_file` inferring `yaml`/`json` by extension; writes `$schema` reference and sidecar schema JSON | `pydantic_evals/pydantic_evals/dataset.py:177-246`, `pydantic_evals/pydantic_evals/dataset.py:556-634`, `pydantic_evals/pydantic_evals/dataset.py:747-794` |
| Dataset portability | Internal `_DatasetModel`/`_CaseModel` with `extra='forbid'` and `NamedSpec` short-form (`'MyEvaluator'`, `{'MyEvaluator': arg}`, `{'MyEvaluator': {k:v}}`) | `pydantic_evals/pydantic_evals/dataset.py:89-108`, `pydantic_ai_slim/pydantic_ai/_spec.py:50-111` |
| Dataset portability | JSON schema generation `model_json_schema_with_evaluators` builds literal/typed-dict unions from registry for validation | `pydantic_evals/pydantic_evals/dataset.py:797-844`, `pydantic_ai_slim/pydantic_ai/_spec.py:272-336` |
| Dataset portability | Loader requires `custom_evaluator_types` / `custom_report_evaluator_types` for non-default evaluators; errors include valid choices | `pydantic_evals/pydantic_evals/dataset.py:686-740`, `pydantic_ai_slim/pydantic_ai/_spec.py:200-236` |
| Prompt template portability | `TemplateStr` wraps Handlebars via `pydantic-handlebars.compile`, validates `{{` presence, renders against `RunContext.deps` or dict dump | `pydantic_ai_slim/pydantic_ai/template.py:16-49`, `pydantic_ai_slim/pydantic_ai/template.py:50-84`, `pydantic_ai_slim/pydantic_ai/template.py:90-116` |
| Prompt template portability | `validate_from_spec_args` auto-compiles `TemplateStr` fields via `TypeAdapter` using `deps_type`/`deps_schema` from validation context | `pydantic_ai_slim/pydantic_ai/_template.py:14-54` |
| Prompt template portability | `format_as_xml` is an LLM-oriented XML serializer (str/bytes/Mapping/Iterable/dataclass/BaseModel) not a prompt registry | `pydantic_ai_slim/pydantic_ai/format_prompt.py:20-77` |
| Prompt template portability | `SystemPromptRunner`/`resolve_system_prompts` keep static strings and async/sync callables; `dynamic_ref` tracked for re-evaluation | `pydantic_ai_slim/pydantic_ai/_system_prompt.py:14-58` |
| Prompt template portability | `InstructionPart.join`/`sorted` separate static vs dynamic instructions for cache-aware rendering | `pydantic_ai_slim/pydantic_ai/messages.py:1778-1787` |
| Cross-provider schema compatibility | `FunctionSchema` built via Pydantic internal `GenerateSchema` + `TypeAdapter` JSON schema; `single_arg_name` unwrapping for model-like args | `pydantic_ai_slim/pydantic_ai/_function_schema.py:110-265`, `pydantic_ai_slim/pydantic_ai/_function_schema.py:319-367` |
| Cross-provider schema compatibility | `ToolDefinition` holds canonical `parameters_json_schema: ObjectJsonSchema` + `return_schema`, `strict` as tri-state (`None` inferred per provider) | `pydantic_ai_slim/pydantic_ai/tools.py:544-554`, `pydantic_ai_slim/pydantic_ai/tools.py:566-582` |
| Cross-provider schema compatibility | `GenerateToolJsonSchema` strips titles; `JsonSchemaTransformer`/`InlineDefsJsonSchemaTransformer` rewrites schemas (inline `$defs`, nullable simplification, `strict` compatibility) per model | `pydantic_ai_slim/pydantic_ai/tools.py:268-274`, `pydantic_ai_slim/pydantic_ai/_json_schema.py:15-59`, `pydantic_ai_slim/pydantic_ai/_json_schema.py:226-233` |
| Cross-provider schema compatibility | Tools prepare hooks allow per-model/per-run schema mutation without changing canonical definition (`Tool.prepare_tool_def`, `ToolsPrepareFunc`) | `pydantic_ai_slim/pydantic_ai/tools.py:104-163`, `pydantic_ai_slim/pydantic_ai/tools.py:512-529` |
| Export/import tools | `ModelMessagesTypeAdapter = TypeAdapter(list[ModelMessage], ser_json_bytes='base64')` for message history dump/load (Python and JSON modes, preserves `ToolReturnContent` discriminated union) | `pydantic_ai_slim/pydantic_ai/messages.py:2769-2772`, `pydantic_ai_slim/pydantic_ai/messages.py:1239-1257` |
| Export/import tools | `AgentRun.all_messages`, `Result.all_messages`, `exception.all_messages` all serialize via `ModelMessagesTypeAdapter.dump_json` | `pydantic_ai_slim/pydantic_ai/run.py:176`, `pydantic_ai_slim/pydantic_ai/result.py:556`, `pydantic_ai_slim/pydantic_ai/exceptions.py:381-383` |
| Export/import tools | Backward compat for pre-usage-refactor histories (`vendor_details`→`provider_details`, `vendor_id`→`provider_response_id`, missing `conversation_id`) | `pydantic_ai_slim/pydantic_ai/messages.py:2694-2701` (usage rebuild) and `tests/test_messages.py:576-632` showing alias handling, `tests/test_messages.py:892-907` |
| Export/import tools | `BinaryContent` byte handling via `ser_json_bytes='base64'` config and custom identifier; `FileUrl` vendor_metadata forwarding per provider documented in docstrings | `pydantic_ai_slim/pydantic_ai/messages.py:541-547`, `pydantic_ai_slim/pydantic_ai/messages.py:218-244` |
| Export/import tools | Eval `SpanTree` context exporter attaches to global `TracerProvider` via `SimpleSpanProcessor`; degrades gracefully for custom providers without `add_span_processor` | `pydantic_evals/pydantic_evals/otel/_context_in_memory_span_exporter.py:85-169`, `pydantic_evals/pydantic_evals/otel/_context_subtree.py:13-33` |

## Answers to Dimension Questions

**1. Can traces be moved between platforms?**
Partially yes via standards, with caveats. The wire format is OTel GenAI semantic conventions, not a proprietary schema: `gen_ai.input.messages`/`gen_ai.output.messages` carry normalized `ChatMessage` parts (`pydantic_ai_slim/pydantic_ai/models/instrumented.py:253`), tool definitions use `gen_ai.tool.definitions` (`pydantic_ai_slim/pydantic_ai/_instrumentation.py:370`), and provider identity is carried in `gen_ai.provider.name`/`gen_ai.system` (`pydantic_ai_slim/pydantic_ai/_instrumentation.py:284`). Any OTLP-compatible collector (Logfire, Jaeger, Honeycomb, Datadog) can ingest spans by configuring a different `TracerProvider` (`pydantic_ai_slim/pydantic_ai/models/instrumented.py:145`). However, payload richness depends on `version` (v6 changes `tool` role vs legacy `user` for tool responses `pydantic_ai_slim/pydantic_ai/models/instrumented.py:410`), on `include_content`/`include_binary_content` flags, and on provider-specific blobs in `provider_details`/`provider_name` (`pydantic_ai_slim/pydantic_ai/messages.py:1898-1916`) that are opaque and only round-trippable to the same provider. Moving traces loses nothing structural but loses provider-specific rendering hints (`code_arg_name`/`code_arg_language` in `pydantic_ai_slim/pydantic_ai/_instrumentation.py:346`). No built-in exporter *converter* beyond the OTel pipeline.

**2. Can eval datasets be reused across systems?**
As plain data, yes; as executable evaluations, no without Pydantic AI. Serialized form is standard YAML or JSON with an optional JSON Schema sidecar (`pydantic_evals/pydantic_evals/dataset.py:769`). Example files validate via `Dataset.model_json_schema_with_evaluators` (`pydantic_evals/pydantic_evals/dataset.py:797`). Cases hold generic `inputs`/`expected_output`/`metadata` (`pydantic_evals/pydantic_evals/dataset.py:111`), so the *data* can be read by any system. But evaluator bindings are `NamedSpec` registry entries (`pydantic_ai_slim/pydantic_ai/_spec.py:50`) that require the Python class to be supplied via `custom_evaluator_types` on load (`pydantic_evals/pydantic_evals/dataset.py:562`); the error path lists valid choices (`pydantic_ai_slim/pydantic_ai/_spec.py:223`). Report evaluators, lifecycle hooks, and `SpanTree`-based evaluators (`pydantic_evals/pydantic_evals/otel/_context_in_memory_span_exporter.py:39`) are code, not declarative data. No export to Hugging Face datasets, LangSmith, or CSV beyond hand-rolling `model_dump(mode='json')` per case.

**3. Can prompts be migrated?**
Bare prompt strings migrate trivially (they are `str` in `SystemPromptPart.content` `pydantic_ai_slim/pydantic_ai/messages.py:187` or `InstructionPart.content` `pydantic_ai_slim/pydantic_ai/messages.py:1762`). Templated prompts using `TemplateStr` (`pydantic_ai_slim/pydantic_ai/template.py:16`) require `pydantic-handlebars` and a `deps_type`/`deps_schema` context (`pydantic_ai_slim/pydantic_ai/template.py:104`), so the *source* string `v._source` (`pydantic_ai_slim/pydantic_ai/template.py:111`) is portable but rendering semantics are not without the same engine and Python type. Dynamic system prompts are Python callables (`pydantic_ai_slim/pydantic_ai/_system_prompt.py:14`) and toolset-derived instructions (`pydantic_ai_slim/pydantic_ai/messages.py:1778`) — not serializable to a file without code. `format_as_xml` (`pydantic_ai_slim/pydantic_ai/format_prompt.py:20`) is a helper, not a standard template language (no Jinja, Mustache, or LangChain hub compatibility). There is no prompt registry, versioning, or CRUD API.

**4. Are tool schemas provider-independent?**
Yes at the canonical layer, no at the wire layer — by design. `Tool`/`FunctionSchema` generates one `ObjectJsonSchema` via Pydantic (`pydantic_ai_slim/pydantic_ai/_function_schema.py:244`) stored as `ToolDefinition.parameters_json_schema` (`pydantic_ai_slim/pydantic_ai/tools.py:554`). Provider adapters then apply `JsonSchemaTransformer` rewrites (`pydantic_ai_slim/pydantic_ai/_json_schema.py:59`) — inlining `$defs` (`pydantic_ai_slim/pydantic_ai/_json_schema.py:226`), `strict` compatibility (`pydantic_ai_slim/pydantic_ai/tools.py:566`), and per-provider `prepare` hooks (`pydantic_ai_slim/pydantic_ai/tools.py:512`) — to meet provider quirks (e.g., OpenAI strict requires `additionalProperties: false`, Google `VALIDATED` mode). `ToolReturn` parameterization (`pydantic_ai_slim/pydantic_ai/messages.py:970`) and `return_schema` propagation (`pydantic_ai_slim/pydantic_ai/_function_schema.py:268`) are likewise canonical; injection as description fallback for non-supporting models happens at adapter time. So a `ToolDefinition` can move between providers, but the exact JSON sent over the wire is intentionally provider-specific.

## Architectural Decisions

- **OTel GenAI as the trace boundary.** Chose `pydantic_ai_slim/pydantic_ai/_otel_messages.py:1` and `pydantic_ai_slim/pydantic_ai/models/instrumented.py:253` to be spec-only with versioned evolution (`pydantic_ai_slim/pydantic_ai/models/instrumented.py:154`) rather than a proprietary event bus. Tradeoff: stability against spec churn vs ability to add bespoke Pydantic AI events (deferred tool deferrals now demoted to not set `ERROR` in v5 `pydantic_ai_slim/pydantic_ai/models/instrumented.py:129` shows the tension). Evidence: policy in `pydantic_ai_slim/pydantic_ai/AGENTS.md:16`.

- **Pydantic TypeAdapter as the serialization contract for messages and eval artifacts.** `ModelMessagesTypeAdapter` (`pydantic_ai_slim/pydantic_ai/messages.py:2769`) with `ser_json_bytes='base64'` and discriminated unions (`pydantic_ai_slim/pydantic_ai/messages.py:2765`) plus `Dataset._DatasetModel` (`pydantic_evals/pydantic_evals/dataset.py:99`) and `NamedSpec` (`pydantic_ai_slim/pydantic_ai/_spec.py:50`) give schema-validated, forward-compatible storage. Tradeoff: strong Python/Pydantic coupling vs dependency on Pydantic version for schema validity; alias handling for `vendor_details` shows migration cost is borne in-code.

- **Registry-based eval extensibility.** `build_registry`/`load_from_registry` (`pydantic_ai_slim/pydantic_ai/_spec.py:159`, `pydantic_ai_slim/pydantic_ai/_spec.py:200`) keep dataset files small (short-form `{'MyEvaluator': arg}` `pydantic_ai_slim/pydantic_ai/_spec.py:82`) but require code registration on load, deliberately sacrificing zero-code portability for type safety and arbitrary Python evaluator logic.

- **Canonical JSON Schema + deferred provider transform.** `GenerateToolJsonSchema` (`pydantic_ai_slim/pydantic_ai/tools.py:268`) plus `JsonSchemaTransformer.walk` (`pydantic_ai_slim/pydantic_ai/_json_schema.py:59`) centralizes provider variance at request-build time, preserving a single source of truth for tool authoring.

- **Handlebars via `TemplateStr` rather than Jinja.** `pydantic_ai_slim/pydantic_ai/template.py:61` leans on the Pydantic ecosystem (`pydantic-handlebars`) with `deps_type` validation (`pydantic_ai_slim/pydantic_ai/template.py:104`), aligning prompt validation with model validation. Tradeoff: type-checked templates vs industry-standard Jinja/Mustache portability.

## Notable Patterns

- **Spec-guarded telemetry module.** `_otel_messages.py` (`pydantic_ai_slim/pydantic_ai/_otel_messages.py:1`) is deliberately thin and type-only, with no business logic, enforcing that only spec-defined `MessagePart` variants are emitted. Strength: prevents accidental coupling to internal `Message` types.

- **Versioned instrumentation with graceful degradation.** `InstrumentationSettings(version=...)` (`pydantic_ai_slim/pydantic_ai/models/instrumented.py:78`) supports 2-6 with warnings, and `version >= 4` branches in message serialization (`pydantic_ai_slim/pydantic_ai/messages.py:1114`) demonstrate additive evolution without breaking old collectors.

- **Discriminated union for tool returns with fallback discriminator.** `_tool_return_content_discriminator` (`pydantic_ai_slim/pydantic_ai/messages.py:1173`) plus `WrapValidator`/`WrapSerializer` (`pydantic_ai_slim/pydantic_ai/messages.py:1246`) distinguishes real `MultiModalContent` from user dicts that happen to share a `kind` value — a portability safeguard for round-tripping mixed content.

- **Short-form evaluator spec serialization.** `NamedSpec.serialize` (`pydantic_ai_slim/pydantic_ai/_spec.py:95`) emits the most compact representation that still round-trips, with `serializes_as_string_keyed_dict` (`pydantic_ai_slim/pydantic_ai/_spec.py:37`) guarding against ambiguous dict-arg collisions.

- **Context-propagation for eval span capture.** `_ContextInMemorySpanExporter` with `ContextVar` ` _EXPORTER_CONTEXT_ID` (`pydantic_evals/pydantic_evals/otel/_context_in_memory_span_exporter.py:34`, `pydantic_evals/pydantic_evals/otel/_context_in_memory_span_exporter.py:85`) isolates concurrent eval runs without global state, falling back to `SpanTreeRecordingError` when no provider is configured (`pydantic_evals/pydantic_evals/otel/_context_subtree.py:13`).

## Tradeoffs

- **Strict spec compliance vs feature velocity.** The OTel-only rule (`pydantic_ai_slim/pydantic_ai/AGENTS.md:16`) keeps traces portable but forces framework-specific concepts (deferred tools, approvals) to map onto generic span attributes or be omitted.

- **Canonical schema vs per-provider quirks.** Keeping one `parameters_json_schema` (`pydantic_ai_slim/pydantic_ai/tools.py:554`) maximizes portability; `JsonSchemaTransformer` then pays the cost of provider-specific rewrites (inline defs, nullable handling `pydantic_ai_slim/pydantic_ai/_json_schema.py:204`). Users rarely see the transformed form, making debugging harder.

- **Template type safety vs portability.** `TemplateStr` with `deps_type` gives compile-time checks against `AgentDepsT` but ties prompts to Python types and the optional `pydantic-handlebars` engine (`pydantic_ai_slim/pydantic_ai/template.py:125`), unlike plain string substitution that any platform can render.

- **Eval data/code separation.** Cases are data (`pydantic_evals/pydantic_evals/dataset.py:111`), evaluators are code; this yields strong type safety but prevents shipping a dataset to a non-Python runner without also shipping source.

- **Message history fidelity vs size.** `ser_json_bytes='base64'` (`pydantic_ai_slim/pydantic_ai/messages.py:544`) preserves binary payloads losslessly; `include_binary_content=False` redaction (`pydantic_ai_slim/pydantic_ai/_instrumentation.py:141`) must then walk arbitrarily nested `ToolReturn`/`DeferredToolRequests` structures (`pydantic_ai_slim/pydantic_ai/_instrumentation.py:170`) to stay consistent, adding complexity.

## Failure Modes / Edge Cases

- **Cross-version message history breakage without aliases.** Without `vendor_details`→`provider_details` aliases, loading a v1 history with `ModelMessagesTypeAdapter.validate_python` would `ValidationError`. Mitigated by `TypeAdapter` alias tests (`tests/test_messages.py:576`) but depends on keeping aliases forever.

- **Dataset load failure for unknown evaluator.** `_spec.load_from_registry` raises `ValueError` listing valid names (`pydantic_ai_slim/pydantic_ai/_spec.py:223`); without `custom_evaluator_types` the dataset is unusable, a silent data-only migration failure mode.

- **Eval span capture silently empty for unsupported TracerProviders.** `_add_context_span_exporter` returns `SpanTreeRecordingError` for providers without `add_span_processor` (e.g., `ddtrace`) or unconfigured `ProxyTracerProvider` (`pydantic_evals/pydantic_evals/otel/_context_in_memory_span_exporter.py:147`), degrading eval to metrics-only without crashing — but evaluators expecting `SpanTree` get an error object they must handle.

- **Handlebars injection via `{{` heuristic.** `TemplateStr.__get_pydantic_core_schema__` treats any string containing `{{` as a template (`pydantic_ai_slim/pydantic_ai/template.py:100`); user content with literal `{{` triggers compilation and may raise, or fall through to `str` branch if validation fails, yielding inconsistent typing in `Union[TemplateStr, str]`.

- **Provider-details opacity.** `ThinkingPart.provider_details` can carry a callable merge delta (`pydantic_ai_slim/pydantic_ai/messages.py:161`) that serializes as `null` in JSON mode (`pydantic_ai_slim/pydantic_ai/messages.py:163`) but survives in Python mode — a trace or message history exported as JSON loses transient reasoning state and is not losslessly re-importable.

- **Compaction provenance leakage.** `sanitize_messages` strips `STANDING_PROMPT_PLANTED_KEY` (`pydantic_ai_slim/pydantic_ai/messages.py:3376`) from untrusted histories; a naïve export/import that bypasses `sanitize_messages` could carry provider-specific compaction blobs that the next provider misinterprets or rejects.

- **Binary content mimetype inference platform variance.** `UploadedFile.media_type` inference via `mimetypes.guess_type` (`pydantic_ai_slim/pydantic_ai/messages.py:898`) is noted as platform-dependent; a history exported on Linux and imported on macOS could infer a different type, altering `BinaryContent.format` (`pydantic_ai_slim/pydantic_ai/messages.py:686`) and downstream routing.

## Future Considerations

- **Provide a prompt registry or externalized prompt packs.** Today prompts live in `Agent(instructions=...)` and `TemplateStr` sources; a file-based registry with versioning, A/B metadata, and export to common formats (Jinja, OpenAI chat template) would close the portability gap. The existing `use_short_form` serialization context (`pydantic_evals/pydantic_evals/dataset.py:783`) could be reused.

- **Add dataset export adapters.** Offer `Dataset.to_hf_dataset()` / `Dataset.to_csv()` and corresponding `from_*` that map `inputs`/`expected_output`/`metadata` without evaluator bindings, plus a CLI `pydantic-evals convert` to translate between YAML/JSON/Parquet.

- **Stabilize instrumentation version 6 as default.** `TODO(v3): default to instrumentation format version 6` (`pydantic_ai_slim/pydantic_ai/models/instrumented.py:157`) signals the remaining churn; freezing on `role='tool'` for `tool_call_response` will simplify cross-consumer expectations.

- **Expose a trace migration shim.** A helper `pydantic_ai.messages.ModelMessagesTypeAdapter` is already the message migration path; a symmetric `otel_messages_to_model_messages` or OTLP replay tool would let traces be rehydrated into a different backend without custom code.

- **Document provider-details portability contract.** Clarify which `provider_details` keys are safe to persist across providers (none, currently) and consider a sanitized export mode that strips them for cross-vendor moves, mirroring `sanitize_messages` (`pydantic_ai_slim/pydantic_ai/messages.py:2954`).

- **Consider a pluggable template engine hook.** Allow `TemplateStr` to delegate to Jinja when `deps` is a flat dict, keeping Handlebars as default but enabling drop-in migration from LangChain/Langfuse prompt templates.

## Questions / Gaps

- No evidence found of a CLI or SDK for bulk export/import of traces or eval datasets between platforms (searched `pydantic_ai_slim/pydantic_ai/_instrumentation.py`, `pydantic_evals/pydantic_evals/dataset.py`, `pydantic_evals/pydantic_evals/otel/`). If such tooling exists in `clai/` or Logfire integration, it was not discoverable from the slim/evals packages alone — search boundary was `pydantic_ai_slim/`, `pydantic_evals/`, `pydantic_graph/`.

- No evidence of prompt versioning, lineage, or experiment-reproducibility linking prompts to eval runs beyond `gen_ai.system_instructions` span attribute (`pydantic_ai_slim/pydantic_ai/models/instrumented.py:291`). The eval `EvaluationReport` metadata path (`pydantic_evals/pydantic_evals/dataset.py:1044`) was not inspected deeply enough to confirm whether prompt hashes are recorded.

- No evidence of cross-provider tool schema equivalence tests (e.g., asserting `openai` and `anthropic` receive logically identical tool definitions from the same `ToolDefinition`). Such tests would strengthen the portability claim but were not located in `tests/test_messages.py` search.

- Open question: does `TemplateStr` rendering validate `deps` shape at call time or only at compile time (`pydantic_ai_slim/pydantic_ai/template.py:68` `check_template_compatibility`)? Failure mode for schema-drifted deps in production evals is unclear.

---

Generated by `19.02-portable-trace-eval-and-prompt-schemas` against `pydantic-ai`.
