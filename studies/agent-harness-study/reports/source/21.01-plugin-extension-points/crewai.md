# Source Analysis: crewai

## 21.01 Plugin and Extension Points

### Source Info

| Field | Value |
|-------|-------|
| Name | crewai |
| Path | `studies/agent-harness-study/sources/crewai` |
| Language / Stack | Python / Pydantic, LiteLLM, uv workspace (crewai, crewai-tools, crewai-core, crewai-files, cli) |
| Analyzed | 2026-08-28 |

## Summary

CrewAI exposes a broad but **code-registered, not marketplace-discovered** extension surface. The framework is intentionally composable: every major seam is an explicit Python interface that a third party subclasses or implements without patching core. The dominant answer to "can a third party add a new tool type without modifying core code?" is **yes** — `BaseTool`/`CrewStructuredTool`/`@tool` (`lib/crewai/src/crewai/tools/base_tool.py:103`, `lib/crewai/src/crewai/tools/structured_tool.py:189`, `lib/crewai/src/crewai/tools/base_tool.py:701`) and `BaseLLM` (`lib/crewai/src/crewai/llms/base_llm.py:150`) make that trivial, and the same pattern repeats for memory, knowledge, embeddings, guardrails, event listeners, A2A extensions, and interception hooks. There is **no dynamic plugin loader** (no `entry_points`, no filesystem scan, no hot-reload); registration is programmatic via `register_*`, constructor injection, or process-wide factories. Lifecycle is limited to register/unregister/clear and a `scoped_hooks` context; there is no versioned plugin manifest, no dependency resolution, and no isolation — hooks and tools run in-process with shared mutable state and fail-open error handling.

## Rating

**5/10 — Present but inconsistent, weakly documented, and fragile on lifecycle/isolation.**

Rationale: explicit, well-typed interfaces exist for tools, LLMs, knowledge, memory, embeddings, hooks, events, and A2A extensions, with tests covering the tool/hook seams. However the model is not a *plugin system*: no dynamic discovery/loading, no sandboxing, no enable/disable lifecycle beyond global list mutation, and extension points are split between legacy registries and the newer generic `InterceptionPoint` dispatcher (`lib/crewai/src/crewai/hooks/dispatch.py:40`) without a single稳定, documented contract. This matches the 4-6 band: usable for developers who control the process, but not durable/observable/safe for untrusted third-party plugins.

## Evidence Collected

Every entry MUST include a file path with line numbers. Format: `path/to/file.ts:NN`.

| Area | Evidence | File:Line |
|------|----------|-----------|
| Extension interfaces - Tool | `BaseTool(BaseModel, ABC)` with abstract `_run`, `__init_subclass__` auto-registry `_TOOL_TYPE_REGISTRY`, and `@tool` decorator factory | `lib/crewai/src/crewai/tools/base_tool.py:103`, `lib/crewai/src/crewai/tools/base_tool.py:109`, `lib/crewai/src/crewai/tools/base_tool.py:701` |
| Tool alt API | `CrewStructuredTool` wrapping arbitrary callables, `from_function` with schema inference; `to_structured_tool()` bridge on BaseTool | `lib/crewai/src/crewai/tools/structured_tool.py:189`, `lib/crewai/src/crewai/tools/structured_tool.py:234`, `lib/crewai/src/crewai/tools/base_tool.py:405` |
| Tool registry for checkpoint deserialization | `_TOOL_TYPE_REGISTRY` populated per subclass; `_resolve_tool_dict` resolves `tool_type` dotted path via `importlib` | `lib/crewai/src/crewai/tools/base_tool.py:51`, `lib/crewai/src/crewai/tools/base_tool.py:59` |
| Bulk tool supply | `crewai-tools` package ships 80+ concrete `BaseTool` subclasses (file, web, DB, cloud, scraping) plus adapters | `lib/crewai-tools/src/crewai_tools/tools/__init__.py:1`, `lib/crewai-tools/src/crewai_tools/tools/file_read_tool/file_read_tool.py:11`, `lib/crewai-tools/src/crewai_tools/adapters/mcp_adapter.py:1` |
| LLM extension | `BaseLLM` abstract with `call`/`acall`/`supports_*`/`_apply_stop_words`; concrete providers (OpenAI, Anthropic, Azure, Bedrock, Gemini, OpenAI-compatible) selected in `LLM.__new__` with `SUPPORTED_NATIVE_PROVIDERS` fallback to LiteLLM | `lib/crewai/src/crewai/llms/base_llm.py:150`, `lib/crewai/src/crewai/llm.py:394`, `lib/crewai/src/crewai/llm.py:666`, `lib/crewai/src/crewai/llm.py:494` |
| LLM provider validation & extensibility | `_validate_model_in_constants` + `_matches_provider_pattern` + `_get_native_provider`; unknown provider falls through to LiteLLM lazy load | `lib/crewai/src/crewai/llm.py:586`, `lib/crewai/src/crewai/llm.py:516`, `lib/crewai/src/crewai/llm.py:666` |
| Knowledge source extension | `BaseKnowledgeSource(BaseModel, ABC)` with `validate_content`/`add`/`aadd`; `_KNOWN_SOURCES` dict (string, csv, json, pdf, …) and `_resolve_knowledge_sources` validator | `lib/crewai/src/crewai/knowledge/source/base_knowledge_source.py:16`, `lib/crewai/src/crewai/knowledge/knowledge.py:23`, `lib/crewai/src/crewai/knowledge/knowledge.py:34` |
| Knowledge storage extension | `BaseKnowledgeStorage(BaseModel, ABC)` with `search`/`asave`/`reset`; `Knowledge` accepts custom `storage` or resolves via `resolve_knowledge_storage` | `lib/crewai/src/crewai/knowledge/storage/base_knowledge_storage.py:13`, `lib/crewai/src/crewai/knowledge/knowledge.py:120` |
| Memory storage extension | `StorageBackend` Protocol (`@runtime_checkable`) with `save`/`search`/`delete`/`update`/`get_record`/`list_scopes`; 2 built-ins (LanceDB, QdrantEdge) + `set_memory_storage_factory` process-wide factory | `lib/crewai/src/crewai/memory/storage/backend.py:44`, `lib/crewai/src/crewai/memory/unified_memory.py:232`, `lib/crewai/src/crewai/memory/storage/factory.py:33` |
| Embedder extension | `PROVIDER_PATHS` dict mapping 19 provider names to dotted class paths; `build_embedder`/`build_embedder_from_dict` via `import_and_validate_definition` | `lib/crewai/src/crewai/rag/embeddings/factory.py:90`, `lib/crewai/src/crewai/rag/embeddings/factory.py:352` |
| Interception hooks (core plugin SEAM) | `InterceptionPoint(str, Enum)` listing 10 points (EXECUTION_START/END, INPUT/OUTPUT, PRE/POST_MODEL_CALL, PRE/POST_TOOL_CALL, PRE/POST_STEP); `on(point, agents, tools)` decorator + generic `register`/`dispatch`/`scoped_hooks` | `lib/crewai/src/crewai/hooks/dispatch.py:40`, `lib/crewai/src/crewai/hooks/dispatch.py:400`, `lib/crewai/src/crewai/hooks/dispatch.py:159`, `lib/crewai/src/crewai/hooks/dispatch.py:119` |
| Legacy hook aliases | `_before_llm_call_hooks = get_global_hook_list(PRE_MODEL_CALL)` etc. sharing queue with new `@on` hooks; `before_llm_call`/`after_llm_call`/`before_tool_call`/`after_tool_call` decorators with `tools`/`agents` filters | `lib/crewai/src/crewai/hooks/llm_hooks.py:161`, `lib/crewai/src/crewai/hooks/tool_hooks.py:134`, `lib/crewai/src/crewai/hooks/decorators.py:88`, `lib/crewai/src/crewai/hooks/decorators.py:192` |
| Hook contexts | `LLMCallHookContext` (executor/messages/llm/agent/task/crew/response, mutable `messages` list) and `ToolCallHookContext` (tool_name/tool_input mutable dict, tool_result) both exposing `request_human_input()` | `lib/crewai/src/crewai/hooks/llm_hooks.py:29`, `lib/crewai/src/crewai/hooks/tool_hooks.py:31` |
| Hook lifecycle management | `register`/`unregister`/`clear`/`clear_all` + `scoped_hooks` contextmanager (global first, scoped second); `HookAborted` propagates, other exceptions swallowed fail-open with optional verbose warning | `lib/crewai/src/crewai/hooks/dispatch.py:119`, `lib/crewai/src/crewai/hooks/dispatch.py:124`, `lib/crewai/src/crewai/hooks/dispatch.py:153`, `lib/crewai/src/crewai/hooks/dispatch.py:265` |
| Event bus as extension point | Singleton `crewai_event_bus: CrewAIEventsBus` (`events/event_bus.py:952`) with `on(event_type)` decorator, `off`, `emit`, `aemit`, `replay`, `scoped_handlers`, `flush`, dependency-aware `build_execution_plan`, `is_replaying()` guard | `lib/crewai/src/crewai/events/event_bus.py:952`, `lib/crewai/src/crewai/events/event_bus.py:245`, `lib/crewai/src/crewai/events/event_bus.py:368`, `lib/crewai/src/crewai/events/event_bus.py:572`, `lib/crewai/src/crewai/events/event_bus.py:832` |
| Guardrail extension | `GuardrailCallable = Callable[[TaskOutput|LiteAgentOutput], tuple[bool, Any]]` plus `GuardrailResult` + `process_guardrail`; `Task.guardrail`/`guardrails` fields with string-to-LLMGuardrail coercion | `lib/crewai/src/crewai/utilities/guardrail_types.py:12`, `lib/crewai/src/crewai/utilities/guardrail.py:60`, `lib/crewai/src/crewai/task.py:252` |
| Flow DSL extension | `@start`/`@listen`/`@router` + `or_`/`and_` combinators authoring Flows; `Flow` runtime composes `ConversationalMixin` | `lib/crewai/src/crewai/flow/dsl/__init__.py:10`, `lib/crewai/src/crewai/flow/flow.py:33` |
| A2A extension (protocol) | `ServerExtension` base with `on_request`/`on_response` + `ServerExtensionRegistry`; `ExtensionRegistry` for client-side `A2AExtension` protocol (inject_tools/augment_prompt/process_response/prepare_message_metadata + `ConversationState`) | `lib/crewai/src/crewai/a2a/extensions/server.py:77`, `lib/crewai/src/crewai/a2a/extensions/server.py:190`, `lib/crewai/src/crewai/a2a/extensions/base.py:57`, `lib/crewai/src/crewai/a2a/extensions/base.py:167` |
| Tool MCP/Platform extension | `BaseAgent.mcps: list[str\|MCPServerConfig]` + `get_mcp_tools`; `BaseAgent.apps: list[PlatformAppOrAction]` + `get_platform_tools` + `get_delegation_tools` abstractions | `lib/crewai/src/crewai/agents/agent_builder/base_agent.py:397`, `lib/crewai/src/crewai/agents/agent_builder/base_agent.py:393`, `lib/crewai/src/crewai/agents/agent_builder/base_agent.py:709` |
| Agent as plugin target | `BaseAgent` ABC defines `execute_task`/`aexecute_task`/`create_agent_executor`/`get_delegation_tools`/`get_platform_tools`/`get_mcp_tools` as extension points; `Agent` (`lib/crewai/src/crewai/agent/core.py:1`) concretizes | `lib/crewai/src/crewai/agents/agent_builder/base_agent.py:686`, `lib/crewai/src/crewai/agents/agent_builder/base_agent.py:205` |
| No dynamic loader evidence | Search for `entry_points`/`importlib.metadata`/`stevedore` shows only `stevedore` transitive dep in `uv.lock` and `pydantic.mypy` plugin; no `pyproject.toml` `[project.entry-points]` nor filesystem plugin scanner found | `pyproject.toml:131`, `uv.lock:581` (reviewed), `grep:plugin` negative |
| Isolation evidence (lack) | Hook dispatcher catches all non-HookAborted exceptions and returns `False` (fail-open); verbose-only warnings; `dispatch` telemetry swallowed in `try/except` | `lib/crewai/src/crewai/hooks/dispatch.py:284`, `lib/crewai/src/crewai/hooks/dispatch.py:248` |
| Docs gap | `docs/edge/en/...` not auto-checked; `docs/docs.json` version gating suggests hooks/tools documented as guides, but no `docs/docs.json` entries found for "Extension API" stability contract in code search | No evidence found (searched `docs/docs.json` not present in source path; `grep:extension` only in `a2a/extensions` code) |

## Answers to Dimension Questions

**1. What can be extended via plugins?**

Broadly: **tools, LLMs, knowledge, memory, embeddings, prompts/history, output validation, control flow, event observation, and A2A protocol behavior.**

- **Tools**: Subclass `BaseTool` (`lib/crewai/src/crewai/tools/base_tool.py:103`) implementing `_run`/`_arun`, or wrap any callable via `CrewStructuredTool.from_function` (`lib/crewai/src/crewai/tools/structured_tool.py:234`) or `@tool` (`lib/crewai/src/crewai/tools/base_tool.py:701`). 80+ examples ship in `lib/crewai-tools/src/crewai_tools/tools/**`. MCP wrappers (`lib/crewai/src/crewai/tools/mcp_tool_wrapper.py:16`, `lib/crewai/src/crewai/tools/mcp_native_tool.py:17`) let external MCP servers inject tools at runtime via `BaseAgent.mcps` (`lib/crewai/src/crewai/agents/agent_builder/base_agent.py:397`). Platform apps (`apps` field, `lib/crewai/src/crewai/agents/agent_builder/base_agent.py:393`) similarly inject enterprise tools.
- **LLMs**: Extend `BaseLLM` (`lib/crewai/src/crewai/llms/base_llm.py:150`) with `call`/`acall`. `LLM.__new__` (`lib/crewai/src/crewai/llm.py:394`) routes `openai/`, `anthropic/`, `azure`, `bedrock`, `gemini`, etc. via `SUPPORTED_NATIVE_PROVIDERS` (`lib/crewai/src/crewai/llm.py:328`) and falls back to LiteLLM lazy loading (`lib/crewai/src/crewai/llm.py:494`). Custom OpenAI-compatible endpoints are supported via `custom_openai=True` + `base_url`.
- **Memory**: Pluggable `StorageBackend` Protocol (`lib/crewai/src/crewai/memory/storage/backend.py:44`) — implement `save`/`search`/`delete`/`update`/`get_record`/`list_scopes` etc. Built-ins: `LanceDBStorage` and `QdrantEdgeStorage` (`lib/crewai/src/crewai/memory/unified_memory.py:232`). Process-wide override via `set_memory_storage_factory` (`lib/crewai/src/crewai/memory/storage/factory.py:33`). `Memory` also exposes `scope()`/`slice()` views (`lib/crewai/src/crewai/memory/unified_memory.py:898`).
- **Knowledge**: `BaseKnowledgeSource` (`lib/crewai/src/crewai/knowledge/source/base_knowledge_source.py:16`) for new source types; `BaseKnowledgeStorage` (`lib/crewai/src/crewai/knowledge/storage/base_knowledge_storage.py:13`) for new vector stores. Registration is the `_KNOWN_SOURCES` map (`lib/crewai/src/crewai/knowledge/knowledge.py:23`) plus `_resolve_knowledge_sources` coercion (`lib/crewai/src/crewai/knowledge/knowledge.py:34`) — extending the map requires code change, but passing an instance is open.
- **Embeddings**: 19 providers enumerated in `PROVIDER_PATHS` (`lib/crewai/src/crewai/rag/embeddings/factory.py:90`). New provider added by adding an entry and a `BaseEmbeddingsProvider` subclass; resolved via `import_and_validate_definition`.
- **Hooks / Prompt & Policy injection**: The most general extension: 10 `InterceptionPoint` values (`lib/crewai/src/crewai/hooks/dispatch.py:40`) covering execution boundaries (`EXECUTION_START`, `INPUT`, `OUTPUT`, `EXECUTION_END`), model/tool calls, and steps. Hooks receive mutable `LLMCallHookContext`/`ToolCallHookContext` and can mutate `messages` in-place, filter inputs, block/redirect via `HookAborted` (`lib/crewai/src/crewai/hooks/dispatch.py:75`), or replace results. Registration via `@on(point, agents=[...], tools=[...])` (`lib/crewai/src/crewai/hooks/dispatch.py:400`) or legacy `@before_llm_call` etc. (`lib/crewai/src/crewai/hooks/decorators.py:88`).
- **Events**: Any event type can be observed via `crewai_event_bus.on(EventClass)` (`lib/crewai/src/crewai/events/event_bus.py:245`). Emitters are all core classes (`lib/crewai/src/crewai/llms/base_llm.py:592` for `LLMCallStartedEvent`, etc.). Listeners can be sync or async with dependency ordering (`handler_graph`).
- **Guardrails / evals**: `guardrail: GuardrailCallable | str` on `Task` (`lib/crewai/src/crewai/task.py:252`) and `guardrails` list; string values are auto-wrapped into `LLMGuardrail`. `process_guardrail` (`lib/crewai/src/crewai/utilities/guardrail.py:123`) standardizes `tuple[bool, Any]` → `GuardrailResult`.
- **Flows / UI-agnostic orchestration**: Decorators `@start`/`@listen`/`@router` plus `or_`/`and_` (`lib/crewai/src/crewai/flow/dsl/__init__.py:10`) let third parties define arbitrary DAGs and state machines over the runtime (`lib/crewai/src/crewai/flow/flow.py:33`).
- **A2A protocol**: Two extension families — protocol-level `ServerExtension` with `is_active`/`on_request`/`on_response` (`lib/crewai/src/crewai/a2a/extensions/server.py:77`) managed by `ServerExtensionRegistry` (`lib/crewai/src/crewai/a2a/extensions/server.py:190`), and CrewAI wrapper `A2AExtension` Protocol (`lib/crewai/src/crewai/a2a/extensions/base.py:57`) for `inject_tools`/`augment_prompt`/`process_response`.

**2. Can plugins be loaded at runtime?**

**No — registration is import-time/programmatic, not dynamic discovery.**

Evidence: no `[project.entry-points]` or `importlib.metadata` plugin loader, no filesystem scan, no `plugin` config key. All search hits for `plugin` refer to the mypy plugin (`lib/crewai/src/crewai/mypy.py:1`) or Claude Code marketplace docs (`README.md:111`), not a CrewAI runtime loader. The actual mechanisms are:

- Direct instantiation: `Agent(tools=[MyTool()])`, `Crew(agents=[...])`, `Knowledge(sources=[MySource()])`, `Memory(storage=MyBackend())`.
- Global registration: `register_before_llm_call_hook(fn)` (`lib/crewai/src/crewai/hooks/llm_hooks.py:191`), `register(point, fn)` (`lib/crewai/src/crewai/hooks/dispatch.py:119`), `crewai_event_bus.on(Event)(fn)` (`lib/crewai/src/crewai/events/event_bus.py:245`).
- Process-wide factories for storage: `set_memory_storage_factory(factory)` (`lib/crewai/src/crewai/memory/storage/factory.py:33`) and `resolve_knowledge_storage` inside `Knowledge.__init__` (`lib/crewai/src/crewai/knowledge/knowledge.py:124`).
- Scoped (execution-bound) hooks: `scoped_hooks()` context manager (`lib/crewai/src/crewai/hooks/dispatch.py:159`) and `event_bus.scoped_handlers()` (`lib/crewai/src/crewai/events/event_bus.py:832`) — still code-driven.
- Closest to dynamic: MCP servers declared as strings/URLs in `BaseAgent.mcps` (`lib/crewai/src/crewai/agents/agent_builder/base_agent.py:397`) — resolved at agent creation, not hot-loaded mid-run; and `EmbeddingProvider` spec dicts that are `import_and_validate_definition`-loaded on demand (`lib/crewai/src/crewai/rag/embeddings/factory.py:251`).

A third party can add a new tool type without modifying core by publishing a Python package that exports a `BaseTool` subclass and having the user `from my_package import MyTool; Agent(tools=[MyTool()])` — but the user must change *their* code/imports. There is no "drop into `plugins/` and it auto-registers."

**3. Are plugins isolated from each other?**

**No isolation: single-process, shared mutable state, cooperative execution.**

- **Same interpreter & memory**: All hooks/tools/handlers share the process. `LLMCallHookContext.messages` is a direct reference to `executor.messages` (`lib/crewai/src/crewai/hooks/llm_hooks.py:88`); `ToolCallHookContext.tool_input` is the live dict mutation target (`lib/crewai/src/crewai/hooks/tool_hooks.py:39`). One hook mutating these affects all subsequent hooks and the core loop — documented with warnings ("Modify in-place… Do NOT replace the list") but no copy-on-write.
- **Execution order, not boundaries**: `dispatch` concatenates global + scoped hooks (`lib/crewai/src/crewai/hooks/dispatch.py:201`) and runs sequentially in registration order. Failure semantics are fail-open: a hook exception (non-`HookAborted`) is swallowed (`lib/crewai/src/crewai/hooks/dispatch.py:284`) and optionally printed only if `verbose` (`lib/crewai/src/crewai/hooks/dispatch.py:285`). `HookAborted` (`lib/crewai/src/crewai/hooks/dispatch.py:75`) aborts the intercepted operation globally, but still runs earlier hooks' side effects.
- **Event bus concurrency is cooperative**: Sync handlers run in a shared `ThreadPoolExecutor(max_workers=10)` (`lib/crewai/src/crewai/events/event_bus.py:179`), async handlers on a single daemon `asyncio` loop (`lib/crewai/src/crewai/events/event_bus.py:184`). No per-plugin thread/process, no resource quotas, no timeouts. `flush(timeout=30)` (`lib/crewai/src/crewai/events/event_bus.py:734`) waits globally.
- **No permission model / dependency isolation**: No capability restriction, no versioned API boundary, no dependency graph for hooks/tools (unlike `handler_graph` which is only for event-listener ordering). A buggy tool can starve the event loop, leak secrets via hook mutation, or abort arbitrary operations.
- **One safeguard**: `HookAborted` carries `source`/`reason` for telemetry (`lib/crewai/src/crewai/hooks/dispatch.py:84`, `lib/crewai/src/crewai/hooks/dispatch.py:343`), and hook dispatch emits `HookDispatchedEvent` (`lib/crewai/src/crewai/hooks/dispatch.py:224`). This is observability, not isolation.

**4. Are extension points documented and stable?**

**Partially. Core tool/LLM/agent surfaces are documented and semver-stabilized; hook/A2A/memory storage contracts are present in code but lack a published stability guarantee.**

- **What is documented**: `BaseTool`/`CrewStructuredTool`/`@tool` have docstring examples (`lib/crewai/src/crewai/tools/base_tool.py:701`, `lib/crewai/src/crewai/tools/structured_tool.py:234`) and a large body of `docs/edge/en` guides (e.g., `build-with-ai.mdx:35` referencing the Claude Code plugin marketplace — note: this is about distributing *skills*, not a runtime plugin API). `BaseLLM` docstring (`lib/crewai/src/crewai/llms/base_llm.py:150`) explicitly says "Users can extend this class." `BaseKnowledgeSource` and `StorageBackend` are typed protocols/ABCs. The public re-export surface (`lib/crewai/src/crewai/__init__.py:69`) stabilizes under Pydantic rebuild with `_base_namespace`.
- **What is weakly documented / unstable**: 
  - The 10 `InterceptionPoint` values (`lib/crewai/src/crewai/hooks/dispatch.py:40`) and the `on()` decorator (`lib/crewai/src/crewai/hooks/dispatch.py:400`) are the *new* dialect; legacy decorators (`lib/crewai/src/crewai/hooks/decorators.py:88`) alias the same queues (`lib/crewai/src/crewai/hooks/llm_hooks.py:161`). No versioned contract or deprecation schedule is declared; the comment in `dispatch.py:7` admits "Legacy adapters" coexist.
  - `resolve_memory_storage`/`resolve_knowledge_storage` factories are undocumented as a public extension API; their docstrings target "application startup" one-time use.
  - A2A `ServerExtension` vs `A2AExtension` distinction is code-commented (`lib/crewai/src/crewai/a2a/extensions/base.py:6`) but not surfaced in user docs search.
  - Tests exist for tool hooks (`lib/crewai/tests/tools/**`) and event listeners, but no contract tests pin the stability of interception point semantics (e.g., that mutating `messages` in a `POST_MODEL_CALL` persists to next iteration — only documented in `LLMCallHookContext` docstring, `lib/crewai/src/crewai/hooks/llm_hooks.py:44`).

## Architectural Decisions

| Decision | Evidence | Rationale / Tradeoff |
|----------|----------|----------------------|
| Tool extensibility via Pydantic BaseModel subclass + auto-registry | `lib/crewai/src/crewai/tools/base_tool.py:103`, `lib/crewai/src/crewai/tools/base_tool.py:109`, `lib/crewai/src/crewai/tools/base_tool.py:51` | Leverages Pydantic validation/serialization and checkpoint serde (`tool_type` dotted path). Tradeoff: ties tools to Pydantic model lifecycle; `CrewStructuredTool` duplicate path (`structured_tool.py:189`) creates two mental models. |
| Dual-tool model (BaseTool for domain authors, CrewStructuredTool for LangChain compat/infra) | `lib/crewai/src/crewai/tools/base_tool.py:405`, `lib/crewai/src/crewai/tools/base_tool.py:521`, `lib/crewai/src/crewai/tools/structured_tool.py:189` | Allows wrapping arbitrary callables and legacy LangChain tools. Tradeoff: `to_structured_tool`/`from_langchain` conversions plus `_original_tool` back-pointer add subtle formatting drift. |
| LLM routing in `LLM.__new__` with native-provider fast path + LiteLLM fallback | `lib/crewai/src/crewai/llm.py:394`, `lib/crewai/src/crewai/llm.py:666`, `lib/crewai/src/crewai/llm.py:494`, `lib/crewai/src/crewai/llms/base_llm.py:150` | Optimizes latency/deps (native SDKs) while keeping open model support. Tradeoff: `_matches_provider_pattern` heuristics risk mis-routing new models; provider string becomes implicit extension API. |
| Memory/knowledge storage as protocols + process-wide factory | `lib/crewai/src/crewai/memory/storage/backend.py:44`, `lib/crewai/src/crewai/memory/storage/factory.py:33`, `lib/crewai/src/crewai/knowledge/storage/base_knowledge_storage.py:13` | Makes storage swappable (LanceDB ↔ Qdrant ↔ fake for tests) without subclassing `Memory`/`Knowledge`. Tradeoff: factory is global mutable state; not thread-safe for multi-tenant processes. |
| Generic hook dispatcher unifying legacy and new dialects | `lib/crewai/src/crewai/hooks/dispatch.py:7`, `lib/crewai/src/crewai/hooks/dispatch.py:97`, `lib/crewai/src/crewai/hooks/llm_hooks.py:161`, `lib/crewai/src/crewai/hooks/tool_hooks.py:134` | Avoids breaking old `register_before_llm_call_hook` code while introducing typed `InterceptionPoint`. Tradeoff: two decorator families, shared queue, subtle filter semantics (`_wrap_with_filters`). |
| Event bus with dependency-ordered execution plan caching | `lib/crewai/src/crewai/events/event_bus.py:58`, `lib/crewai/src/crewai/events/event_bus.py:245`, `lib/crewai/src/crewai/events/handler_graph.py:1` | Enables observability plugins (checkpoint, trace, console) without core changes, with level-wise parallel execution. Tradeoff: plan cache invalidated only on handler registration, not on dependency mutation. |
| MCP/AMP as external tool extension rather than Python subclass | `lib/crewai/src/crewai/agents/agent_builder/base_agent.py:397`, `lib/crewai/src/crewai/tools/mcp_tool_wrapper.py:16` | Delegates sandboxing/network isolation to external servers. Tradeoff: no streaming/thinking parity guarantees for external tools. |

## Notable Patterns

- **Decorator + Registry + Context mutation**: Every seam (hooks, events) follows: decorator marks a function (`_create_hook_decorator`, `lib/crewai/src/crewai/hooks/decorators.py:18`), registry stores ordered list (`_global_hooks` dict, `lib/crewai/src/crewai/hooks/dispatch.py:97`), dispatcher resolves global+scoped and runs with reducer (`dispatch`/`run_hooks`, `lib/crewai/src/crewai/hooks/dispatch.py:357`), context mutated in place (`LLMCallHookContext.messages`, `lib/crewai/src/crewai/hooks/llm_hooks.py:57`).
- **Filter-aware wrappers**: Both hook and old decorator paths sanitize names and wrap filters (`sanitize_tool_name`, `lib/crewai/src/crewai/hooks/dispatch.py:375`, `lib/crewai/src/crewai/hooks/decorators.py:43`) so `tools=["delete_file"]` works regardless of `Tool Name:` formatting.
- **Checkpoint-friendly identity**: `tool_type` dotted path (`lib/crewai/src/crewai/tools/base_tool.py:203`) and `llm_type`/`executor_type` registries (`lib/crewai/src/crewai/agents/agent_builder/base_agent.py:76`, `lib/crewai/src/crewai/agents/agent_builder/base_agent.py:138`) allow serde of extension choices.
- **Fail-open hook dispatch**: Non-abort exceptions swallowed; `HookAborted` is the only abort path (`lib/crewai/src/crewai/hooks/dispatch.py:280`, `lib/crewai/src/crewai/hooks/dispatch.py:284`). Mirrors long-standing CrewAI philosophy of "framework protects against buggy user hook".
- **Extension Registry pattern**: `ExtensionRegistry`/`ServerExtensionRegistry`/`A2AExtensionRegistry` all follow `register` + `inject/augment/process` fan-out (`lib/crewai/src/crewai/a2a/extensions/base.py:167`, `lib/crewai/src/crewai/a2a/extensions/server.py:190`).

## Tradeoffs

- **Breadth vs depth**: Many narrow extension points (each with its own Protocol/ABC) rather than a single `Plugin` manifest. Easy to find the seam, hard to govern cross-cutting concerns (e.g., a tool that also needs memory and event hooks requires touching 3 registries).
- **Inheritance vs composition**: Tools/LLMs require subclassing `Base*`; memory/storage use composition. The former couples to Pydantic/model lifecycle; the latter is cleaner but inconsistent style.
- **Global vs scoped state**: Global `register_*` persists for process lifetime; `scoped_hooks`/`scoped_handlers` are contextvar-bound. No per-crew/per-agent plugin manifest; tenant isolation requires manual `clear_all_global_hooks()` (`lib/crewai/src/crewai/hooks/__init__.py:47`).
- **Observability without isolation**: `HookDispatchedEvent` telemetry (`lib/crewai/src/crewai/hooks/dispatch.py:224`) and `MemorySaveFailedEvent` etc. make failures visible, but not contained — a slow handler blocks the pool.
- **Dynamic idioms without a loader**: Embedder and MCP paths use lazy `import_and_validate_definition` (`lib/crewai/src/crewai/rag/embeddings/factory.py:251`), but tools still require Python import. MCP partially fills the "install without code change" gap via external servers.

## Failure Modes / Edge Cases

- **Hook ordering non-determinism**: Global registration order is preserved, but imports that side-effect-register hooks (decorators at import time) depend on import order. Two packages registering `PRE_MODEL_CALL` hooks get interleaved based on `sys.modules` import sequence — not a defined contract.
- **Mutation contract violation**: Reassigning `context.messages = []` in `LLMCallHookContext` breaks the executor (documented at `lib/crewai/src/crewai/hooks/llm_hooks.py:44`), but not enforced. Similarly `context.tool_input = {}` is silently ignored (`lib/crewai/src/crewai/hooks/tool_hooks.py:41`).
- **Swallowed errors hide bugs**: Hook exceptions other than `HookAborted` are silently dropped unless `verbose=True` on the agent (`lib/crewai/src/crewai/hooks/tool_hooks.py:167`, `lib/crewai/src/crewai/hooks/dispatch.py:284`). Users may believe a hook ran when it failed.
- **Tool validation fast-fail on agent construction**: `BaseAgent.validate_tools` (`lib/crewai/src/crewai/agents/agent_builder/base_agent.py:526`) raises `ValueError` for any tool lacking `name`/`func`/`description`; this is eagerly validated on every import, not lazily — a broken tool package blocks all agents.
- **Factory override races**: `set_memory_storage_factory` (`lib/crewai/src/crewai/memory/storage/factory.py:33`) is global mutable; late override doesn't affect already-constructed `Memory` instances, and concurrent use in tests can leak state.
- **Embedding dimension mismatch**: `EmbeddingDimensionMismatchError` (`lib/crewai/src/crewai/memory/storage/backend.py:11`) requires manual `crewai reset-memories` after default embedder upgrade — an extension (custom embedder) that changes dims invalidates prior storage with no migration path.
- **Checkpoint guardrail loss**: `serialize_guardrail_for_json` drops callable guardrails with a warning (`lib/crewai/src/crewai/utilities/guardrail.py:24`); restored checkpoints silently run without that guardrail — a silent security regression.
- **Event bus shutdown races**: `CrewAIEventsBus.shutdown` flushes then clears handlers (`lib/crewai/src/crewai/events/event_bus.py:897`); late `emit` during shutdown is silently dropped (`lib/crewai/src/crewai/events/event_bus.py:602`), losing telemetry.

## Future Considerations

- Introduce a declarative `crewai.plugins` manifest (or `pyproject.toml` [tool.crewai] section) with versioned `Plugin` interface (activate/deactivate, dependencies, isolation level) so dynamic discovery via `importlib.metadata.entry_points(group="crewai.plugins")` can replace import-time side effects.
- Version the interception contract: freeze `InterceptionPoint` values, promise semantic versioning for context fields, and add contract tests (mutating `messages` in POST_MODEL_CALL persists) so new points don't break old hooks.
- Replace fail-open silent swallow with structured `HookErrorEvent` + per-hook timeout + optional `strict` mode that can be enabled in CI.
- Provide per-crew/per-agent plugin scopes instead of only global + contextvar-scoped; add `Crew(plugins=[...])` / `Agent(plugins=[...])` so `clear_all_global_hooks` is not needed for test isolation.
- Unify tool extension story: deprecate duplicate `Tool` vs `CrewStructuredTool` naming and offer a single `BaseTool[TParams, TResult]` Generic with first-class async (`_arun`) examples.
- Document extension points as a reference page (stable/unstable markers) and generate it from the actual Protocol/ABC signatures so drift is CI-caught.

## Questions / Gaps

- No evidence found for a published compatibility matrix mapping CrewAI versions → supported extension APIs; search of `docs/docs.json` and source comments yielded no `stable`/`experimental` annotations beyond `lib/crewai/src/crewai/agent/internal` implicit convention.
- No evidence found for plugin health/readiness probes (e.g., hook returning "unhealthy" to be disabled). `HookAborted` aborts the operation, not the hook.
- No evidence found for policy/UI extension points in OSS: `plus_api.py` (`lib/crewai/src/crewai/plus_api.py:1`) suggests enterprise tier may host policy evaluation, but no OSS `PolicyExtension` was located.
- Unclear fate of `crewai-tools` adapters (zapier, lancedb, enterprise) — each adapter (`adapters/mcp_adapter.py:1`, `adapters/enterprise_adapter.py:1`) has bespoke initialization not governed by a shared plugin lifecycle.

---
Generated by `Dimension 21.01: Plugin and Extension Points` against `crewai`.
