# Source Analysis: openai-agents-sdk

## 21.01 Plugin and Extension Points

### Source Info

| Field | Value |
|-------|-------|
| Name | openai-agents-sdk |
| Path | `studies/agent-harness-study/sources/openai-agents-sdk` |
| Language / Stack | Python 3.10+ / Pydantic, OpenAI Responses API, MCP |
| Analyzed | 2026-08-28 |

## Summary

openai-agents-sdk exposes no generic `Plugin` loader or registry. Extensibility is provided as a set of explicit, typed interface points that third parties implement via subclassing or by passing callables/dataclasses into `Agent`/`RunConfig`/`Runner`. The dominant pattern is "bring your own implementation of an ABC/Protocol" rather than discovery/loading. Providers (`Model`/`ModelProvider`), memory (`Session`/`SessionABC`), MCP servers (`MCPServer`), tracing (`TracingProcessor`/`TraceProvider`), lifecycle hooks (`AgentHooks`/`RunHooks`), tools (`FunctionTool` via `@function_tool`), guardrails, handoffs, and `RunConfig` filters (`call_model_input_filter`, `tool_error_formatter`) are all first-class extension surfaces. Optional capabilities live under `src/agents/extensions/` and are gated behind `pyproject.toml:37-55` extras with lazy import via `__getattr__`. There is no runtime plugin discovery (no entrypoints, no importlib.metadata scanning, no manifest), no lifecycle manager beyond per-object `connect/cleanup` or `shutdown/close`, and no isolation: all extensions share the event loop, memory, and process state, with only exception-swallowing safeguards.

## Rating

**5 / 10 — Present but inconsistent, weakly documented as a unified model, and fragile on lifecycle/isolation.**

Rationale: The SDK offers many stable, typed extension interfaces with tests and docs for individual surfaces (tools, MCP, memory, tracing). However it lacks a unified plugin abstraction: no dynamic loading at runtime, no versioned registry, no isolation, and extensions under `agents.extensions.*` are scattered and classified inconsistently (some are true plugins like memory backends, others are utilities like `ToolOutputTrimmer`). Lifecycle is ad-hoc per-subsystem and requires manual orchestration (e.g., caller must `connect/cleanup` MCP servers). Isolation is absent by design — a faulty tracing processor or MCP filter can only be tolerated via try/except, not sandboxing.

## Evidence Collected

Every entry MUST include a file path with line numbers. Format: `path/to/file.ts:NN`.

| Area | Evidence | File:Line |
|------|----------|-----------|
| Extension interfaces — Model | `class Model(abc.ABC)` with abstract `get_response` / `stream_response`, `ModelTracing` enum | `src/agents/models/interface.py:37-135` |
| Extension interfaces — ModelProvider | `class ModelProvider(abc.ABC)` abstract `get_model`, `aclose` hook | `src/agents/models/interface.py:138-161` |
| Extension interfaces — MultiProvider routing | Prefix-based router `litellm/` / `any-llm/` fallback with `MultiProviderMap` add/remove | `src/agents/models/multi_provider.py:18-61`, `src/agents/models/multi_provider.py:164-174` |
| Extension interfaces — Tool | `FunctionTool` dataclass, `ComputerTool`, `ShellTool`, `HostedMCPTool`, `function_tool` decorator alias | `src/agents/tool.py:440-600`, `src/agents/decorators.py:1-19`, `src/agents/__init__.py:140-205` |
| Extension interfaces — Tool creation helper | `function_tool` wraps Python func via `FuncSchema` + `function_schema()` | `src/agents/function_schema.py:290-490`, `src/agents/tool.py:698-757` |
| Extension interfaces — MCP Server | `class MCPServer(abc.ABC)` with `connect/cleanup/list_tools/call_tool` + approval/policy normalization | `src/agents/mcp/server.py:543-603`, `src/agents/mcp/server.py:867-1139` |
| Extension interfaces — MCP ToolFilter | `ToolFilterContext`, `ToolFilterStatic`, `ToolFilterCallable`, `MCPToolMetaResolver`, `MCPToolCustomDataExtractor` | `src/agents/mcp/util.py:92-209` |
| Extension interfaces — Memory/Session | `Session` Protocol and `SessionABC` ABC, `OpenAIResponsesCompactionAwareSession` | `src/agents/memory/session.py:15-107` |
| Extension interfaces — Memory backends (extensions) | Lazy export map for `RedisSession`, `SQLAlchemySession`, `EncryptedSession`, etc. | `src/agents/extensions/memory/__init__.py:41-51`, `src/agents/extensions/memory/__init__.py:54-73` |
| Extension interfaces — Tracing | `TracingProcessor` ABC (6 abstract methods), `TracingExporter`, `TraceProvider` ABC | `src/agents/tracing/processor_interface.py:9-129`, `src/agents/tracing/provider.py:222-297` |
| Extension interfaces — Tracing concrete provider | `DefaultTraceProvider`, `SynchronousMultiTracingProcessor` with `add_tracing_processor/set_processors/shutdown/force_flush` | `src/agents/tracing/provider.py:93-220`, `src/agents/tracing/provider.py:300-511` |
| Extension interfaces — Tracing registration | Public `add_trace_processor`, `set_trace_processors`, `set_trace_provider`, `set_tracing_disabled` | `src/agents/tracing/__init__.py:94-113`, `src/agents/tracing/setup.py:27-65` |
| Extension interfaces — Agent hooks | `RunHooksBase`, `AgentHooksBase` with `on_llm_start/on_llm_end/on_agent_start/on_tool_start/on_tool_end/on_handoff` | `src/agents/lifecycle.py:13-107`, `src/agents/lifecycle.py:203-207` |
| Extension interfaces — Guardrails | `InputGuardrail`/`OutputGuardrail` dataclasses, `ToolInputGuardrail`/`ToolOutputGuardrail` + `tool_input_guardrail` decorator | `src/agents/guardrail.py:71-185`, `src/agents/tool_guardrails.py:152-206`, `src/agents/decorators.py:6-8` |
| Extension interfaces — Handoff customization | `Handoff` dataclass with `input_filter`, `HandoffInputFilter` type alias, registry mutators `set_conversation_history_wrappers` | `src/agents/handoffs/__init__.py:125-173`, `src/agents/handoffs/history.py` (re-exported) |
| Extension interfaces — Prompt extension | `Prompt` TypedDict, `DynamicPromptFunction`, `PromptUtil.to_model_input` | `src/agents/prompts.py:23-48`, `src/agents/prompts.py:56-82` |
| Extension interfaces — RunConfig hooks | `CallModelInputFilter`, `ToolErrorFormatter`, `OutputGuardrailBlockedMessageFormatter`, `ToolExecutionConfig` | `src/agents/run_config.py:76-132`, `src/agents/__init__.py:113-124` |
| Extension interfaces — Sandbox | `SandboxRunConfig`, `SandboxArchiveLimits`, `SandboxConcurrencyLimits` | `src/agents/run_config.py:218-347` |
| Extension interfaces — Sandbox hosted extensions | 7 hosted sandbox clients under `extensions/sandbox/{blaxel,cloudflare,daytona,e2b,modal,runloop,vercel}` | `src/agents/extensions/sandbox/blaxel/sandbox.py:1`, `src/agents/extensions/sandbox/e2b/sandbox.py:1` etc. (glob `src/agents/extensions/sandbox/**/*`) |
| Extension interfaces — Hosted tool extensions | `ApplyPatchTool`, `ShellTool`, `ComputerTool` with lifecycle callbacks | `src/agents/tool.py:841-865`, `src/agents/tool.py:1362-1454` |
| Plugin loaders — no dynamic discovery | Only `import importlib.metadata.version` for package version; no entrypoint scanning | `src/agents/version.py:1-5` |
| Plugin loaders — lazy optional loading | `extensions/memory/__init__.py:54-73` uses `import_module` + `_LAZY_EXPORTS` with optional-dep error via `raise_optional_dependency_error`; similar lazy fallback in `MultiProvider._create_fallback_provider` uses deferred import | `src/agents/extensions/memory/__init__.py:54-73`, `src/agents/models/multi_provider.py:164-173` |
| Plugin loaders — agent aggregates | `AgentBase.get_all_tools()` merges `self.tools` + `mcp_tools` via `MCPUtil.get_all_function_tools`; no registry file | `src/agents/agent.py:250-293` |
| Plugin loaders — MCP tool acquisition | `MCPUtil.get_all_function_tools` / `get_function_tools` enumerates servers, checks duplicate names, builds prefixed overrides | `src/agents/mcp/util.py:266-339` |
| Plugin lifecycle hooks — MCP | `connect()`, `cleanup()`, `invalidate_tools_cache()`, `_cache_dirty`, `cache_tools_list` flag | `src/agents/mcp/server.py:592-624`, `src/agents/mcp/server.py:867-970`, `src/agents/mcp/server.py:1102-1104` |
| Plugin lifecycle hooks — Model | `Model.close()`, `Model._cleanup_on_run_end(owner)`, `ModelProvider.aclose()`, `MultiProvider.aclose()` fan-out | `src/agents/models/interface.py:47-57`, `src/agents/models/interface.py:155-161`, `src/agents/models/multi_provider.py:254-278` |
| Plugin lifecycle hooks — Tracing | `TracingProcessor.shutdown()`, `force_flush()`, `TraceProvider.shutdown()` / `force_flush()`, `DefaultTraceProvider.shutdown(timeout)` with deadline slicing | `src/agents/tracing/processor_interface.py:109-129`, `src/agents/tracing/provider.py:283-297`, `src/agents/tracing/provider.py:177-220` |
| Plugin lifecycle hooks — Computer | `ComputerProvider(create/dispose)`, `resolve_computer` / `dispose_resolved_computers` via WeakKeyDictionary | `src/agents/tool.py:346-371`, `src/agents/tool.py:891-970` |
| Plugin lifecycle hooks — Run hooks | No `setup/teardown` — only per-event callbacks on `RunHooksBase`/`AgentHooksBase` | `src/agents/lifecycle.py:13-104` |
| Plugin isolation — exception swallowing only | `SynchronousMultiTracingProcessor` wraps each processor call in try/except + diagnostic log; `MCPUtil._apply_dynamic_tool_filter` catches and excludes tool on error | `src/agents/tracing/provider.py:117-175`, `src/agents/mcp/server.py:1060-1084` |
| Plugin isolation — absence of sandboxing | All extensions run in same asyncio loop; `MCPServer` docs state caller must manage lifecycle manually | `src/agents/agent.py:198-205`, `src/agents/mcp/server.py:592-615` |
| Extension documentation — surface | No single "plugin author guide"; each interface has its own docs: tracing, mcp, tools, realtime, sessions, sandbox clients, visualization | `docs/tracing.md`, `docs/mcp.md`, `docs/tools.md`, `docs/sessions/*`, `docs/sandbox/clients.md`, `docs/visualization.md`, `docs/ref/extensions/*` |
| Extension documentation — extensions namespace | Docstrings note extensions are optional backends gated by extras; sessions index lists 7 backends with install hints | `src/agents/extensions/memory/__init__.py:1-8`, `docs/sessions/index.md` equivalents referenced via `pyproject.toml:37-55` |
| Extension documentation — API stability | `AGENTS.md:79-108` public API positional-compatibility contract, ExecPlan compatibility risk callout | `AGENTS.md:79-108` |

## Answers to Dimension Questions

**1. What can be extended via plugins?**

Yes, via explicit interfaces — not via a single plugin registry:

- **Models/providers:** Implement `Model` + `ModelProvider` (`src/agents/models/interface.py:37-161`). `MultiProvider` adds prefix routing (`litellm/`, `any-llm/`) and custom `MultiProviderMap` (`src/agents/models/multi_provider.py:18-61`). Third-party adapters `LitellmProvider` / `AnyLLMProvider` already ship under `src/agents/extensions/models/` (gated by `pyproject.toml:40-41` extras).
- **Tools:** Register any `FunctionTool` via `@function_tool` / `function_tool()` (`src/agents/tool.py:440-600`, `src/agents/decorators.py:10`), plus `ShellTool`, `ComputerTool`, `ApplyPatchTool`, `CustomTool`, `FileSearchTool`, `WebSearchTool`, `HostedMCPTool`. Agent-as-tool via `Agent.as_tool()` (`src/agents/agent.py:583-1040`) is another extensible surface.
- **MCP servers (tool providers):** Subclass `MCPServer` (`src/agents/mcp/server.py:543`). Filtering/customization via `ToolFilter`, `MCPToolMetaResolver`, `MCPToolCustomDataExtractor` (`src/agents/mcp/util.py:92-209`).
- **Memory/sessions:** Implement `Session` Protocol or subclass `SessionABC` (`src/agents/memory/session.py:15-107`). 7 production backends ship as extensions: `AsyncSQLiteSession`, `RedisSession`, `SQLAlchemySession`, `MongoDBSession`, `DaprSession`, `AdvancedSQLiteSession`, `EncryptedSession` (`src/agents/extensions/memory/__init__.py:41-51`).
- **Tracing:** Implement `TracingProcessor` / `TracingExporter` and optionally `TraceProvider` (`src/agents/tracing/processor_interface.py:9-142`, `src/agents/tracing/provider.py:222-297`). Registered globally via `add_trace_processor`/`set_trace_provider` (`src/agents/tracing/__init__.py:94-113`, `src/agents/tracing/setup.py:27-50`).
- **Guardrails:** `InputGuardrail`/`OutputGuardrail` (`src/agents/guardrail.py:71-185`) and `ToolInputGuardrail`/`ToolOutputGuardrail` (`src/agents/tool_guardrails.py:152-206`) with `@input_guardrail` etc. decorators.
- **Handoffs:** Fully customizable via `handoff()` helper and `HandoffInputFilter`/`HandoffHistoryMapper` (`src/agents/handoffs/__init__.py:125-173`), plus built-in filters like `remove_all_tools` (`src/agents/extensions/handoff_filters.py:33-56`).
- **Lifecycle/hooks:** `AgentHooks` / `RunHooks` (`src/agents/lifecycle.py:13-207`) and `RunConfig` callbacks (`call_model_input_filter`, `tool_error_formatter`, etc. — `src/agents/run_config.py:76-132`). Prompts via `Prompt`/`DynamicPromptFunction` (`src/agents/prompts.py:23-48`). Tool output trimming as a `CallModelInputFilter` (`src/agents/extensions/tool_output_trimmer.py:88-136`).
- **Sandbox/hosted execution:** 7 sandbox clients (`blaxel`, `cloudflare`, `daytona`, `e2b`, `modal`, `runloop`, `vercel`) under `src/agents/extensions/sandbox/**/*`, configured via `SandboxRunConfig` (`src/agents/run_config.py:218-290`).

**2. Can plugins be loaded at runtime?**

No generic runtime plugin loading. Evidence: global searches for `importlib.metadata`, `entrypoint`, `plugin` show only `version.py:1` and per-backend lazy imports (`src/agents/extensions/memory/__init__.py:54-73`, `src/agents/models/multi_provider.py:166-172`). Adding a new model provider requires explicit code: `MultiProvider(provider_map=MultiProviderMap().add_provider(...))` or `RunConfig(model_provider=my_provider)`. Adding a memory backend requires `Runner.run(..., session=MySession())`. Adding MCP tools requires `Agent(mcp_servers=[MyMCPServer()])`. No config-driven discovery, no hot-reload, no manifest file — construction is programmatic and happens at agent/run creation time. Optional deps are lazily imported on first attribute access, not discovered.

**3. Are plugins isolated from each other?**

No. All extensions execute in the host process's asyncio loop with shared memory, global tracing state (`GLOBAL_TRACE_PROVIDER` in `src/agents/tracing/setup.py:11-14`), and shared `RunContextWrapper`.

- **MCP isolation:** Only logical: `cache_tools_list`/`invalidate_tools_cache()` (`src/agents/mcp/server.py:1102-1104`) and per-server `list_tools` error redaction; a misbehaving server can still block `get_all_function_tools` or exhaust retries.
- **Tracing isolation:** Best-effort try/except around each processor in `SynchronousMultiTracingProcessor` (`src/agents/tracing/provider.py:117-175`) so one faulty processor does not crash the run, but no thread/process/memory isolation and no resource quotas.
- **Tool isolation:** `FunctionTool.on_invoke_tool` runs as `await` in the same task group; user code exceptions are caught only if `failure_error_function` is set (`src/agents/tool.py:653-667`). No sandboxing.
- **Memory isolation:** None cross-session; sessions are caller-provided objects.


**4. Are extension points documented and stable?**

Partially. Each interface is well-documented individually (`docs/mcp.md`, `docs/tools.md`, `docs/tracing.md`, `docs/sessions/*.md`, `docs/sandbox/clients.md`, etc.) and typed with `py.typed` (`src/agents/py.typed`). The repository advertises a public API positional-compatibility contract (`AGENTS.md:79-108`) and requires ExecPlan compatibility risk assessment. However there is no single "author a plugin" guide, no versioned extension SDK, and the `agents.extensions.*` namespace is treated as "beyond core" without a stability promise — some extensions (e.g., `experimental/codex`) are explicitly marked experimental (`src/agents/extensions/experimental/__init__.py:1`), and optional backends can break on extra-version bumps (e.g., `e2b==2.20.0` pinned in `pyproject.toml:52`). Stability is therefore interface-by-interface, not system-wide.

## Architectural Decisions

- **Interface-per-concern vs. unified plugin SPI.** The SDK chose 8+ discrete ABCs/Protocols (`Model`, `Session`, `MCPServer`, `TracingProcessor`, `AgentHooks`, etc.) over a single `Plugin` type. Citations: `src/agents/models/interface.py:37`, `src/agents/memory/session.py:15`, `src/agents/mcp/server.py:543`, `src/agents/tracing/processor_interface.py:9`, `src/agents/lifecycle.py:13`. This maximizes type safety and keeps core (`src/agents/`) dependency-free from optional extras, at the cost of no discoverable catalog of "what can be extended."

- **Extras-gated extensions namespace.** Third-party integrations live under `src/agents/extensions/` and are optional dependencies (`pyproject.toml:37-55` lists `litellm`, `any-llm`, `sqlalchemy`, `redis`, `mongodb`, `dapr`, etc.). Loading is lazy via `__getattr__` + `import_module` with `raise_optional_dependency_error` (`src/agents/extensions/memory/__init__.py:54-73`). This keeps install size small but means extension availability is import-time, not runtime-pluggable.

- **Caller-owned lifecycle.** `MCPServer` javadoc explicitly requires callers to `connect()` before and `cleanup()` after use (`src/agents/agent.py:198-205`). Similarly `Model.close`/`Provider.aclose` and `TracingProcessor.shutdown` are opt-in. No central orchestrator (e.g., `PluginManager.start_all/stop_all`) exists.

- **Prefix routing for model extensibility.** `MultiProvider` interprets `model_name` prefixes (`openai/`, `litellm/`, `any-llm/`, custom) to select providers (`src/agents/models/multi_provider.py:227-253`). This is a lightweight extension mechanism that avoids registry files but couples naming to routing semantics and requires code-level registration of custom prefixes.

- **Global tracing singleton.** `GLOBAL_TRACE_PROVIDER` with `threading.Lock` + `atexit` shutdown (`src/agents/tracing/setup.py:11-14`, `src/agents/tracing/provider.py:300-330`) gives a process-wide, observable extension point. Tradeoff: easy for third parties (`add_trace_processor`), but not isolated per-run and subject to global state pitfalls.

## Notable Patterns

- **Decorator-factory + dataclass duality.** Tools and guardrails expose both `@function_tool`/`@input_guardrail` decorators and raw dataclass constructors, allowing inline and programmatic registration (`src/agents/decorators.py:1-19`, `src/agents/tool.py:698-757`, `src/agents/guardrail.py:200-270`).
- **Protocol + ABC dual-shape.** Memory offers both a structural `Session` Protocol and a `SessionABC` ABC (`src/agents/memory/session.py:15-107`) — third parties implement the Protocol (avoids inheritance), first-party and extensions subclass the ABC. Same pattern for tracing via duck-typed `TracingProcessor` ABC vs. any object with same methods.
- **Filter/Mapper as first-class extension.** Almost every subsystem exposes a filter callback: `ToolFilter`/`HandoffInputFilter`/`HandoffHistoryMapper`/`CallModelInputFilter`/`MCPToolMetaResolver`, making composition via closures the extension mechanism (`src/agents/mcp/util.py:92-138`, `src/agents/handoffs/__init__.py:118-122`, `src/agents/run_config.py:76-132`).
- **WeakKeyDictionary per-run caches.** `ComputerTool` resolver stores per-`RunContextWrapper` computers (`src/agents/tool.py:879-888`, `src/agents/tool.py:891-939`) and MCP `MultiProvider` caches fallback providers (`src/agents/models/multi_provider.py:193-197`) — per-run resource binding without leaking across runs, but no eviction policy beyond GC.
- **Example: third-party tool type without core modification.** A caller can add a wholly new hosted tool type by wrapping it in a `FunctionTool` (declare `name`, `description`, `params_json_schema`, provide `on_invoke_tool` ctx+args → Awaitable) and passing it in `Agent(tools=[my_tool])` — no core change needed. Hosted tools that require wire-format support (e.g., a new Responses API tool) would require core change because serialization lives in `src/agents/models/openai_responses.py`.

## Tradeoffs

- **Typed interfaces vs. discovery.** Strong typing and explicit wiring make extensions self-documenting and testable, but sacrifice zero-config extensibility — every extension requires code changes at composition time; there is no plugin directory or config that can be dropped in without touching `Agent` construction.
- **Optional extras vs. cohesion.** Gating behind extras keeps core slim and avoids importing heavy clients (SQLAlchemy, redis, dapr), but fragments documentation and makes cross-extension behavior (e.g., combining `EncryptedSession` + `SQLAlchemySession`) rely on wrapper composition rather than a registry.
- **Global tracing vs. per-run scoping.** Global `TraceProvider` simplifies adoption (one `set_trace_provider` call) but makes multi-tenant or test isolation harder; `SynchronousMultiTracingProcessor` fans out synchronously so a slow exporter blocks the traced operation.
- **Manual lifecycle vs. framework orchestration.** Requiring explicit `connect/cleanup` gives control to the host (important for network resources), but places burden on every consumer and yields inconsistent error messages if forgotten (`MCPServer` errors bubble as `UserError` or `McpError`).
- **Prefix routing simplicity vs. namespace collision risk.** String-prefix routing is trivial to understand but can confuse users when model IDs legitimately contain slashes (e.g., `openai/gpt-4o` alias vs. literal model_id — mitigated by `openai_prefix_mode`/`unknown_prefix_mode` flags in `src/agents/models/multi_provider.py:89-90`).

## Failure Modes / Edge Cases

- **Uninitialized MCP server:** Passing an `MCPServer` that was never `connect()`ed yields transport errors at `list_tools` time (`src/agents/mcp/server.py:1212-1236` redacts to `UserError`), not at construction — fail-late.
- **Duplicate tool names:** `MCPUtil.get_all_function_tools` raises `UserError` on name collisions across servers; mitigation requires `include_server_in_tool_names=True` (`src/agents/mcp/util.py:304-335`). Missing this flag in multi-server setups is a common footgun.
- **Faulty tracing processor leaks across runs:** Because `GLOBAL_TRACE_PROVIDER` holds a tuple of processors under a lock (`src/agents/tracing/provider.py:99-116`), an exception inside one processor is caught per-span but retained processors are never auto-removed — a poisoned processor stays registered until manual `set_trace_processors` or process exit. `shutdown()` is best-effort (`_safe_debug` guards closed streams — `src/agents/tracing/provider.py:24-46`).
- **Leaked per-run resources:** `ComputerProvider` disposers are stored in `WeakKeyDictionary` keyed by `RunContextWrapper`; if `dispose_resolved_computers` is not called (e.g., runner short-circuit), disposer closures keep references until GC (`src/agents/tool.py:943-970`).
- **Session protocol evolution:** Third parties implementing `Session` structurally must add `wrapper: RunContextWrapper` to all 4 methods to receive context (`src/agents/memory/session.py:155-196`); old implementations silently lose wrapper without error, breaking context-aware backends.
- **Optional dependency ImportError at attribute access time:** Accessing `agents.extensions.memory.RedisSession` without `pip install openai-agents[redis]` raises late, with message from `raise_optional_dependency_error` (`src/agents/extensions/memory/__init__.py:64-70`) — not at import of `agents`, surprising in production.
- **Tool schema conversion fallback:** `MCPUtil.to_function_tool` attempts `ensure_strict_json_schema` and silently falls back to non-strict on failure (`src/agents/mcp/util.py:544-561`), so malformed MCP schemas produce lenient runtime schemas without warning beyond a debug log.
- **No load isolation:** A long-running or hanging `TracingProcessor.on_span_end` blocks `SynchronousMultiTracingProcessor.on_span_end` (`src/agents/tracing/provider.py:162-175`) because dispatch is synchronous over the tuple snapshot — one slow exporter stalls all.

## Future Considerations

- Introduce a lightweight extension manifest (e.g., `[tool.agents.extensions]` in `pyproject.toml` or a `register_extension()` helper) so third parties can ship a package that auto-registers a `ModelProvider` prefix or a `Session` factory without editing host code — currently requires explicit code wiring at `Agent`/`RunConfig` construction.
- Consolidate lifecycle under an optional `ExtensionManager` that tracks `connect/cleanup` for MCP servers and `shutdown/close` for trace/model providers, with deterministic ordering and aggregate error collection instead of ad-hoc caller responsibility.
- Provide per-run (scoped) tracing provider override in `RunConfig` to avoid global state for tests and multi-tenant servers; already partially mitigated by `tracing_disabled` (`src/agents/run_config.py:397-405`) but not by scoped processor sets.
- Document a single, versioned "Plugin Author Guide" that enumerates all extension points, stability tier (stable/experimental), and SemVer guarantees; today stability is only implied via `py.typed` and `AGENTS.md:79-108` positional compatibility rule.
- Add isolation/testing harness: e.g., a `Fake`/`InMemory` reference implementation for each extension interface (like `testing/model.py`) beyond the current `tests/` — would lower the barrier for third-party authors and enable conformance tests.

## Questions / Gaps

- **Are extension points versioned or feature-flagged?** No evidence of a registry version or capability advertisement beyond type checking. `CURRENT_SCHEMA_VERSION` exists for `RunState` but not for extension interfaces — searched `src/agents/**/*` for `extension.*version` / `capability` with no versioned extension contract.
- **Can a third party add a new tool *wire type* (not just a FunctionTool) without forking?** Implemented behavior suggests no for hosted Responses API tools: adding a new `Hosted*Tool` requires changes in `src/agents/models/openai_responses.py` serialization. Function-tool shape is user-extensible; tool wire type is not.
- **Is there a published stability guarantee for `agents.extensions.*`?** Not found. `AGENTS.md:79-108` covers public constructors but does not enumerate extensions; docs mark Codex as experimental but not other extensions' status. Declared as gap.
- **No automated discoverability for custom extensions.** Verified absence of `importlib.metadata.entry_points` scanning; only TODO-level implicit support via `MultiProviderMap.set_mapping` (`src/agents/models/multi_provider.py:32-34`). Whether this is intentional minimalism vs. missing feature is unclear from docs.

---

Generated by `21.01-plugin-and-extension-points` against `openai-agents-sdk`.
