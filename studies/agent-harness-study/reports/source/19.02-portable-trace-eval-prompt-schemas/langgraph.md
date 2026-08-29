# Source Analysis: langgraph

## Portable Trace, Eval, and Prompt Schemas

### Source Info

| Field | Value |
|-------|-------|
| Name | langgraph |
| Path | `studies/agent-harness-study/sources/langgraph` |
| Language / Stack | Python (langgraph, langchain-core, langsmith) / TypeScript SDKs |
| Analyzed | 2026-08-28 |

## Summary

LangGraph (v1.2.6) delegates all observability to the `langchain-core` callback/tracer system rather than defining its own provider-agnostic trace schema. Execution emits `LangChainTracer` runs (`trace_id`, `dotted_order`, `parent_run_id`) that are LangSmith-specific; optional OpenTelemetry support is declared as dependencies but only referenced as monkey-patched `opentelemetry-instrumentation-langchain` and not wired into `Pregel`. Prompts are a LangChain union type (`str | SystemMessage | Callable | Runnable`) wrapped as `RunnableCallable` with no standard template format or import/export. Tool schemas reuse `langchain_core.tools.BaseTool` / Pydantic JSON Schema (OpenAI `type:function` vs Anthropic `name`) with internal injection filtering, giving partial cross-provider translation but no explicit provider-independent spec. Eval datasets are absent in-core; datasets/evals live in external LangSmith Studio (redirect only). No trace/eval/prompt migration or export/import tooling exists; only internal checkpoint version/DB migrations and `JsonPlusSerializer` for persistence.

## Rating

**Score: 3 / 10 — Absent / ad-hoc portability**

**Rationale:** Tracing is fully coupled to `LangChainTracer`/LangSmith run schema with no abstraction, OTEL is optional/unintegrated. Prompt templates have no standard or portable serialization (code-carried callables). Tool schemas depend on langchain-core with ad-hoc OpenAI/Anthropic branching. No eval dataset schema and no cross-platform export/import. This matches rubric 1–3: implicit, ad-hoc, unsafe for migration without rewriting data.

## Evidence Collected

Every entry MUST include a file path with line numbers. Format: `path/to/file.ts:NN`.

| Area | Evidence | File:Line |
|------|----------|-----------|
| Trace format abstraction | `RunnableCallable.invoke` and `ainvoke` create `LangChainTracer` runs, set context via `langsmith.run_helpers._set_tracing_context` and emit `callback_manager.on_chain_start/on_chain_end` with `trace` flag controlling whether tracing occurs | `libs/langgraph/langgraph/_internal/_runnable.py:400-426`, `libs/langgraph/langgraph/_internal/_runnable.py:473-503`, `libs/langgraph/langgraph/_internal/_runnable.py:69-102` |
| Trace format abstraction | `RunnableSeq.invoke/ainvoke/stream` manually manage `LangChainTracer` run_map and `set_config_context` with `trace_inputs` filtering per step | `libs/langgraph/langgraph/_internal/_runnable.py:651-692`, `libs/langgraph/langgraph/_internal/_runnable.py:695-746`, `libs/langgraph/langgraph/_internal/_runnable.py:748-800` |
| Trace format abstraction | `coerce_to_runnable` always sets `trace=True` for nodes; `Pregel.NodeBuilder.do` coerces with `trace=True`, making tracing opt-out only via explicit `trace=False` on internal helpers | `libs/langgraph/langgraph/_internal/_runnable.py:529-553`, `libs/langgraph/langgraph/pregel/main.py:303-314`, `libs/langgraph/langgraph/pregel/_call.py:182-236` |
| Trace format abstraction | Graph lifecycle exposes only `GraphCallbackHandler` (`on_interrupt`/`on_resume`) wrapping `BaseCallbackHandler`; no OTEL span interface | `libs/langgraph/langgraph/callbacks.py:87-112`, `libs/langgraph/langgraph/callbacks.py:219-345` |
| Trace provider coupling | Interop test asserts LangSmith run fields `dotted_order`, `parent_run_id`, `trace_id` equality across nested graphs via `LangChainTracer` mock | `libs/langgraph/tests/test_tracing_interops.py:11-18`, `libs/langgraph/tests/test_tracing_interops.py:92-118` |
| Trace provider coupling | `FakeTracer` extends `langchain_core.tracers.BaseTracer` and records `Run.trace_id` mapping, confirming dependency on langchain-core tracer types | `libs/langgraph/tests/fake_tracer.py:7-68` |
| Trace OTEL optional | `opentelemetry-api`, `opentelemetry-sdk`, `opentelemetry-exporter-otlp-proto-http` listed as optional `python_full_version >= '3.11'` dependencies but not imported in pregel; only mention is `opentelemetry-instrumentation-langchain monkey-patch` comment | `libs/langgraph/uv.lock:1602-1605`, `libs/cli/uv.lock:1016-1018`, `libs/langgraph/tests/test_graph_callbacks.py:283-284` |
| Trace / checkpoint traces | `RemotePregel` notes `distributed_tracing` sends LangSmith headers, not OTEL headers | `libs/langgraph/langgraph/pregel/remote.py:163-164` |
| Prompt template portability | `Prompt = SystemMessage | str | Callable[[StateSchema], LanguageModelInput] | Runnable[StateSchema, LanguageModelInput]` type definition | `libs/prebuilt/langgraph/prebuilt/chat_agent_executor.py:121-126` |
| Prompt template portability | `_get_prompt_runnable` converts `str`→`SystemMessage`, `SystemMessage`→prepend, `callable`/`async callable`→`RunnableCallable`, `Runnable`→passthrough; no serialization or standard template language | `libs/prebuilt/langgraph/prebuilt/chat_agent_executor.py:137-170` |
| Prompt template portability | `create_react_agent` docs show prompt examples as `str` or callable, dynamic prompt composed as `_get_prompt_runnable(prompt) | model` | `libs/prebuilt/langgraph/prebuilt/chat_agent_executor.py:366-372`, `libs/prebuilt/langgraph/prebuilt/chat_agent_executor.py:590-616` |
| Prompt template portability | Example graph uses `langchain_core.prompts.ChatPromptTemplate.from_messages` / `from_template` externally, not a LangGraph-provided template registry | `libs/cli/examples/graphs/storm.py:16-32`, `libs/cli/examples/graphs/storm.py:85-146` |
| Tool schema portability | `ToolNode` built on `langchain_core.tools.BaseTool`, validates via `Pydantic ValidationError`→`ToolInvocationError`, filters injected args from LLM schema | `libs/prebuilt/langgraph/prebuilt/tool_node.py:76-84`, `libs/prebuilt/langgraph/prebuilt/tool_node.py:957-966`, `libs/prebuilt/langgraph/prebuilt/tool_node.py:510-563` |
| Tool schema portability | `_should_bind_tools` branches on OpenAI `type==function` + `function.name` vs Anthropic `name` field to compare bound tools; indicates provider-specific shape handling rather than unified spec | `libs/prebuilt/langgraph/prebuilt/chat_agent_executor.py:202-210` |
| Tool schema portability | `ToolCallRequest` / `ToolCallWithContext` internal TypedDict with `tool_call`, `__type`, `state`; tool call payload is `ToolCall` dict (`name`, `args`, `id`, `type`) from `langchain_core.messages.ToolCall` | `libs/prebuilt/langgraph/prebuilt/tool_node.py:132-150`, `libs/prebuilt/langgraph/prebuilt/tool_node.py:286-307`, `libs/prebuilt/langgraph/prebuilt/tool_node.py:1014-1045` |
| Tool schema portability | Injection hides `InjectedState`/`InjectedStore`/`ToolRuntime` keys from tool JSON schema via `all_injected_keys` filtering and dict-merge `{**stripped_args, **injected_args}` | `libs/prebuilt/langgraph/prebuilt/tool_node.py:1315-1431`, `libs/prebuilt/langgraph/prebuilt/tool_node.py:1967-2030` |
| Dataset / eval portability | No `Dataset`, `Eval`, or `Experiment` types in `libs/langgraph/langgraph` or `libs/prebuilt`; only external redirect to LangSmith datasets studio | `docs/redirects.json:123` |
| Export/import tools | No `export*`, `import*` for traces/evals/prompts found; search across `libs/langgraph` returns only DB `MIGRATIONS` for checkpoint/store SQLite/Postgres | `libs/checkpoint-sqlite/langgraph/store/sqlite/base.py:718-751`, `libs/checkpoint-postgres/langgraph/store/postgres/base.py:1092-1135` |
| Export/import tools | Checkpoint persistence uses `JsonPlusSerializer` (`SerializerProtocol`) with `allowed_msgpack_modules` allowlist and optional `EncryptedSerializer`, but this is for state snapshots, not traces/evals | `libs/checkpoint/langgraph/checkpoint/serde/jsonplus.py:82-131`, `libs/checkpoint/langgraph/checkpoint/base/__init__.py:21-209`, `libs/checkpoint/langgraph/checkpoint/serde/encrypted.py:12-40` |
| Data migration | Only internal checkpoint version migration `_migrate_checkpoint` (`pending_sends` → `TASKS` channel) and delta-channel migration preservation tests | `libs/langgraph/langgraph/pregel/main.py:1135-1143`, `libs/langgraph/tests/test_delta_channel_migration.py:1-11`, `libs/checkpoint-sqlite/tests/test_delta_channel_migration.py:1-11` |
| Stream / trace portability | Debug/stream payloads `TaskPayload`, `TaskResultPayload`, `CheckpointPayload`, `StreamPart` are LangGraph-internal checkpoint/stream schemas, not OTEL or OpenInference | `libs/langgraph/langgraph/types.py:142-219`, `libs/langgraph/langgraph/types.py:262-351`, `libs/langgraph/langgraph/pregel/debug.py:41-71` |

## Answers to Dimension Questions

**1. Can traces be moved between platforms?**
No. LangGraph does not define a provider-agnostic trace format. All tracing goes through `langchain_core.tracers.langchain.LangChainTracer` and `langsmith.run_helpers._set_tracing_context` (`libs/langgraph/langgraph/_internal/_runnable.py:43`, `libs/langgraph/langgraph/_internal/_runnable.py:77-79`, `libs/langgraph/langgraph/_internal/_runnable.py:412-418`). Tests assert LangSmith-specific fields `trace_id`, `dotted_order`, `parent_run_id` (`libs/langgraph/tests/test_tracing_interops.py:113-118`). OTEL packages are optional dependencies gated on Python 3.11 (`libs/langgraph/uv.lock:1602-1605`) and only referenced as external monkey-patching (`libs/langgraph/tests/test_graph_callbacks.py:283`). `RemotePregel` distributed tracing uses `x-ls-trace` headers (`libs/langgraph/tests/test_remote_graph_v3.py:594`). Moving to Datadog/Grafana/OpenInference would require rewriting collector code and run schema; there is no export to OTLP/JSON.

**2. Can eval datasets be reused across systems?**
No evidence found. Grep for `eval`, `dataset`, `experiment` in `libs/langgraph` and `libs/prebuilt` yields no dataset abstraction; the only dataset reference is a docs redirect `"/cloud/how-tos/datasets_studio": "https://docs.langchain.com/langsmith/use-studio"` (`docs/redirects.json:123`). Evaluations are expected to run in LangSmith Studio outside the harness. No portable dataset file format (JSONL/CSV/Parquet) or loader is shipped.

**3. Can prompts be migrated?**
Partially but not portably. `Prompt` is a code-level union `str | SystemMessage | Callable | Runnable` (`libs/prebuilt/langgraph/prebuilt/chat_agent_executor.py:121-126`) materialized via `_get_prompt_runnable` (`libs/prebuilt/langgraph/prebuilt/chat_agent_executor.py:137-170`). Callable/Runnable prompts are arbitrary Python functions and cannot be serialized. String prompts are plain `SystemMessage` prefixes (`libs/prebuilt/langgraph/prebuilt/chat_agent_executor.py:143-147`). No standard template format (Jinja, Mustache, F-String registry) is enforced; examples import `ChatPromptTemplate` from `langchain_core.prompts` directly (`libs/cli/examples/graphs/storm.py:16`). Migration would require hand-rewriting callables to target platform's template spec; no prompt import/export tool exists.

**4. Are tool schemas provider-independent?**
Weakly. Tools are `langchain_core.tools.BaseTool` instances validated via Pydantic JSON Schema (`libs/prebuilt/langgraph/prebuilt/tool_node.py:76-84`, `libs/prebuilt/langgraph/prebuilt/tool_node.py:957-971`). `_should_bind_tools` explicitly handles two provider shapes — OpenAI `{"type":"function","function":{"name":...}}` vs Anthropic `{"name":...}` (`libs/prebuilt/langgraph/prebuilt/chat_agent_executor.py:202-207`) — indicating provider-specific branching rather than a single canonical spec. Injected args (`InjectedState`, `InjectedStore`, `ToolRuntime`) are stripped from the JSON schema sent to the LLM (`libs/prebuilt/langgraph/prebuilt/tool_node.py:1424-1430`) so that part is portable, but name/description/parameters JSON still must be translated per provider outside LangGraph (no built-in converter beyond langchain-core's `bind_tools`). No OpenAPI/JSON Schema export utility for tools is provided.

## Architectural Decisions

- **Delegation to langchain-core for all observability** (`libs/langgraph/langgraph/_internal/_runnable.py:43-64` imports `LangChainTracer`, `_StreamingCallbackHandler`; `libs/langgraph/langgraph/_internal/_config.py:236-272` builds `CallbackManager`). Keeps harness lean but cements LangSmith as the trace sink; OTEL remains an external instrumentation concern.
- **Prompt as Runnable composition** (`libs/prebuilt/langgraph/prebuilt/chat_agent_executor.py:137-170` `_get_prompt_runnable` + `590: _get_prompt_runnable(prompt) | model`). Enables dynamic `callable(state, runtime)->model` patterns but makes prompts unserializable and non-portable.
- **Tool injection via TypedDict + Pydantic filtering** (`libs/prebuilt/langgraph/prebuilt/tool_node.py:286-307` `ToolCallWithContext`, `libs/prebuilt/langgraph/prebuilt/tool_node.py:1967-2030` `_get_all_injected_args`). Hides system state from LLM tool schema, improving security but tying schema to langchain-core's `InjectedToolArg` mechanism.
- **Checkpoint-centric persistence vs trace-centric observability** (`libs/checkpoint/langgraph/checkpoint/serde/jsonplus.py:82` `JsonPlusSerializer` + `libs/langgraph/langgraph/pregel/main.py:1135-1143` `_migrate_checkpoint`). Migration investment goes to state durability, not trace portability; `JsonPlusSerializer` allowlisting (`allowed_msgpack_modules`) governs state, not eval datasets.
- **Optional OTEL as add-on, not first-class** (`libs/langgraph/uv.lock:1602`, `libs/cli/uv.lock:1016`). Declares dependency but never imports or emits spans in `PregelLoop`/`RunnableCallable`, leaving portability to third-party instrumentation.

## Notable Patterns

- **Tracer-coupled Runnable wrappers**: Every user node is wrapped via `coerce_to_runnable(..., trace=True)` (`libs/langgraph/langgraph/_internal/_runnable.py:529-553`) and `RunnableCallable`/`RunnableSeq` explicitly extract `LangChainTracer` from handler list to fetch `run` (`libs/langgraph/langgraph/_internal/_runnable.py:412`, `libs/langgraph/langgraph/_internal/_runnable.py:677`).
- **Injection-aware schema filtering**: `all_injected_keys` collected from `get_all_basemodel_annotations` + `get_type_hints` and stripped before `tool.invoke` (`libs/prebuilt/langgraph/prebuilt/tool_node.py:1999-2001`, `libs/prebuilt/langgraph/prebuilt/tool_node.py:1424-1429`).
- **Stream-mode-agnostic debug payloads**: `TaskPayload`/`TaskResultPayload`/`CheckpointPayload` emitted via `map_debug_tasks`/`map_debug_task_results` (`libs/langgraph/langgraph/pregel/debug.py:41-71`, `libs/langgraph/langgraph/pregel/debug.py:106-128`) — rich internal telemetry but not mapped to OTLP.
- **Versioned checkpoint migration**: `_migrate_checkpoint` upgrades `v<4` checkpoints by moving `pending_sends` to `TASKS` channel (`libs/langgraph/langgraph/pregel/main.py:1135-1143`), mirrored in SQLite/Postgres `MIGRATIONS` (`libs/checkpoint-sqlite/langgraph/store/sqlite/base.py:981-1011`).

## Tradeoffs

- **Single-vendor trace ergonomics vs portability**: Using `LangChainTracer` gives immediate LangSmith tracing with propagated `trace_id`/`dotted_order` and `_set_tracing_context` support for nested `@ls.traceable` (`libs/langgraph/tests/test_tracing_interops.py:68-77`) but prevents switching to OTEL-native backends without adapter code.
- **Code-defined prompts vs declarative templates**: Callable prompts enable context-aware generation (`prompt(state, runtime)` in `libs/prebuilt/langgraph/prebuilt/chat_agent_executor.py:279-282`) at the cost of no portable serialization; declarative Jinja prompts would be exportable but less expressive.
- **Rich checkpoint serde vs trace portability**: `JsonPlusSerializer` with pluggable `SerializerProtocol` and encrypted wrapper (`libs/checkpoint/langgraph/checkpoint/serde/encrypted.py:12`) provides flexible state persistence but is not reused for traces/evals, so eval datasets still require LangSmith.
- **Provider branching in tool binding**: Handling OpenAI and Anthropic shapes in one function (`libs/prebuilt/langgraph/prebuilt/chat_agent_executor.py:202-210`) avoids N provider SDKs but leaks provider specifics into harness and omits other providers (Google/ Bedrock tool formats not handled).

## Failure Modes / Edge Cases

- **OTEL collector receives no spans**: Since `PregelLoop` never creates OTEL spans, deploying OTEL collector + `opentelemetry-instrumentation-langchain` monkey-patch may silently drop LangGraph-internal steps (e.g., `ChannelWrite`, `Send` distribution) that are only visible via `stream_mode="debug"` payloads.
- **Trace context loss on `trace=False` nodes**: Internal helpers (`_read`, `_write`, `_branch`) set `trace=False` (`libs/langgraph/langgraph/pregel/_read.py:48`, `libs/langgraph/langgraph/pregel/_write.py:64`, `libs/langgraph/langgraph/graph/_branch.py:134`) — if a user copies this pattern for custom nodes, child runs disappear from trace without warning.
- **Prompt callable serialization failure**: Attempting to pickle/checkpoint a graph with a `lambda state: ...` prompt fails `JsonPlusSerializer` allowlist checks or produces unimportable references; `allowed_msgpack_modules=None` fallback permits arbitrary deserialization (`libs/checkpoint/langgraph/checkpoint/serde/jsonplus.py:238-255` warning).
- **Tool schema mismatch across providers**: A tool bound for OpenAI (`type:function`) sent to Anthropic without translation raises `Missing tools` or silent argument coercion errors (`libs/prebuilt/langgraph/prebuilt/chat_agent_executor.py:214-217`).
- **Dataset lock-in**: No eval dataset export; migrating off LangSmith requires scraping API traces and reconstructing datasets manually, with no guarantee of `tool_calls`/`ToolMessage` ID correlation preservation (IDs are run-local in `libs/prebuilt/langgraph/prebuilt/tool_node.py:1007-1012`).
- **Checkpoint migration fragility**: `_migrate_checkpoint` only handles `v<4`→`TASKS`; future channel renames (e.g., `BinaryOperatorAggregate`→`DeltaChannel`) require explicit migration tests (`libs/langgraph/tests/test_delta_channel_migration.py:1`) — missing migrations corrupt resumed threads.

## Future Considerations

- Introduce a `TraceExporter` protocol (OTLP/JSON) alongside `SerializerProtocol` so `RunnableCallable`/`RunnableSeq` can tee spans to OTEL without monkey-patching; gate via `RunnableConfig` similar to `store`/`checkpointer`.
- Add `PromptTemplate` registry with standard serialization (e.g., `ChatPromptTemplate` JSON) and allow `create_react_agent(prompt=PromptTemplate.from_file(...))` to replace callable prompts for portable cases, with lint rule flagging non-serializable prompts.
- Provide `BaseTool` → provider-agnostic JSON Schema exporter (`tool.get_json_schema(provider="openai"|"anthropic"|"generic")`) and tests covering translation, removing ad-hoc `_should_bind_tools` branching.
- Ship `langgraph eval` CLI or `ExportDatasets` utility that dumps traces to OpenTelemetry/JSONL and rehydrates via `InMemorySaver`, mirroring existing `checkpoint-conformance` suite for cross-backend validation.
- Document `LANGGRAPH_STRICT_MSGPACK` impact on prompt/tool serialization and add CI job that runs tracing interop tests against both LangSmith and OTEL export paths.

## Questions / Gaps

- No evidence found for eval dataset schema or harness — searched `libs/langgraph/langgraph`, `libs/prebuilt`, `libs/checkpoint`, `docs/` for `eval`, `dataset`, `experiment`; only `docs/redirects.json:123` references external LangSmith studio. Confirm whether eval is intentionally out-of-scope for harness.
- No evidence found for prompt versioning or prompt registry — `_get_prompt_runnable` treats prompts as ephemeral runnables; no `PromptVersion` or storage.
- No evidence found for trace export/import CLI — `libs/cli` contains graph deployment templates (`libs/cli/langgraph_cli/templates.py:44-56`) but no `export traces` command.
- OTEL integration depth unclear — `uv.lock` declares deps but no code imports `opentelemetry` in `libs/langgraph/langgraph`; verify if instrumentation is expected via `opentelemetry-instrumentation-langchain` only or planned native support.
- Cross-provider tool schema translation completeness — `_should_bind_tools` covers OpenAI/Anthropic only (`libs/prebuilt/langgraph/prebuilt/chat_agent_executor.py:202-210`); Google `functionDeclarations` not handled — confirm support matrix.

---

Generated by `19.02-portable-trace-eval-and-prompt-schemas` against `langgraph`.
