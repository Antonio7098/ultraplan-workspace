# Source Analysis: agent-framework

## 21.01 Plugin and Extension Points

### Source Info

| Field | Value |
|-------|-------|
| Name | agent-framework |
| Path | `studies/agent-harness-study/sources/agent-framework` |
| Language / Stack | Python (agent-framework-core, provider packages via uv workspace) + .NET |
| Analyzed | 2026-08-28 |

## Summary

Agent Framework is a core-plus-providers framework where third-party extensibility is achieved through abstract base classes and protocols rather than a runtime plugin registry. The design encourages adding new capabilities by subclassing or implementing protocols in user code and passing instances to `Agent` or `BaseChatClient`. No dynamic filesystem plugin loader, no plugin marketplace, and no isolation: extensions run in-process with full trust. The most mature surfaces are tools/middleware/providers; the newest surfaces (Skills, evals, compaction, FIDES security) are explicitly marked `@experimental` and therefore unstable. Adding a new tool type does **not** require core modification — implement `FunctionTool` or `@tool` or a new `MCPTool` or register a `FunctionTool` via progressive disclosure — but adding a new *tool kind* (e.g., shell) still requires a supporting chat-client protocol such as `SupportsShellTool`.

## Rating

**6 / 10 — Present but inconsistent, weakly documented and fragile under evolution**

**Rationale:** Extension interfaces are explicit, typed (`Protocol`, `ABC`), and tested for the critical paths (tools, middleware, chat clients, executors, context providers). Providers are decoupled via lazy `__getattr__` imports, so a third party can publish `agent-framework-<provider>` without forking core. However: (1) there is no unified plugin loader/discovery (except file-based SKILL.md scanning via `FileSkillsSource`); (2) lifecycle is per-invocation hooks only (`before_run`/`after_run`, `on_checkpoint_save/restore`) with no install/unload/update/version lifecycle; (3) isolation is absent — all extensions share the Python process, `session.state`, and `additional_properties`; (4) many extension points carry `ExperimentalWarning` (`SKILLS`, `HARNESS`, `EVALS`, `PROGRESSIVE_TOOLS`, `MCP_LONG_RUNNING_TASKS`, `FIDES`), signaling breaking-change risk; (5) extension documentation is docstring-only, with no stability contract or versioning for extensions.

## Evidence Collected

Every entry MUST include a file path with line numbers. Format: `path/to/file.ts:NN`.

| Area | Evidence | File:Line |
|------|----------|-----------|
| Tool extension interface | `class FunctionTool(SerializationMixin):` wraps any Python callable with JSON-schema generation, `approval_mode`, `invocation_count` limits, `result_parser` | `python/packages/core/agent_framework/_tools.py:243` |
| Tool decorator | `@tool` overloads convert `Callable` → `FunctionTool`; supports `schema: BaseModel | Mapping` and `max_invocations` | `python/packages/core/agent_framework/_tools.py:1179` |
| Tool normalization/flattening | `normalize_tools()` flattens callables, `FunctionTool`, `MCPTool`, dict specs, iterables and collection wrappers (e.g., Foundry `ToolboxVersionObject`) | `python/packages/core/agent_framework/_tools.py:941` |
| Progressive tool exposure | `FunctionInvocationContext.add_tools()` / `remove_tools()` mutate live `tools` list for next LLM iteration; `@experimental(PROGRESSIVE_TOOLS)` | `python/packages/core/agent_framework/_middleware.py:287` and `python/packages/core/agent_framework/_middleware.py:324` |
| Agent extension protocol | `@runtime_checkable class SupportsAgentRun(Protocol):` defines `run(stream)`, `create_session`, `get_session` — duck-typing, no inheritance required | `python/packages/core/agent_framework/_agents.py:188` |
| Chat client extension protocol | `@runtime_checkable class SupportsChatGetResponse(Protocol):` + `class BaseChatClient(ABC):` with `_inner_get_response` to override | `python/packages/core/agent_framework/_clients.py:82` / `python/packages/core/agent_framework/_clients.py:214` |
| Provider lazy loading (dynamic extension) | Each provider package uses `__getattr__` + `importlib.import_module(IMPORT_PATH)` to lazy-load `agent-framework-*` without hard import | `python/packages/core/agent_framework/azure/__init__.py:27` ; `python/packages/core/agent_framework/openai/__init__.py:14` ; `python/packages/core/agent_framework/redis/__init__.py:13` |
| Middleware extension points | `class AgentMiddleware(ABC): process(context:AgentContext, call_next)` / `FunctionMiddleware` / `ChatMiddleware` + `@agent_middleware` functional form | `python/packages/core/agent_framework/_middleware.py:465` / `_middleware.py:524` / `_middleware.py:588` / `_middleware.py:678` |
| Middleware pipelines | `AgentMiddlewarePipeline.execute()`, `FunctionMiddlewarePipeline`, `ChatMiddlewarePipeline` chain via `MiddlewareWrapper` | `python/packages/core/agent_framework/_middleware.py:843` / `python/packages/core/agent_framework/_middleware.py:930` / `python/packages/core/agent_framework/_middleware.py:1003` |
| Context/memory extension | `class ContextProvider:` with `before_run`/`after_run(source_id, state)` over `SessionContext` and `AgentSession.state` | `python/packages/core/agent_framework/_sessions.py:351` / `_sessions.py:370` / `_sessions.py:391` |
| History/memory built-ins | `class HistoryProvider(ContextProvider):` + `class InMemoryHistoryProvider(HistoryProvider):` + `@experimental(FILE_HISTORY) class FileHistoryProvider` | `python/packages/core/agent_framework/_sessions.py:413` / `python/packages/core/agent_framework/_sessions.py:814` / `python/packages/core/agent_framework/_sessions.py:893` |
| Skills system (experimental) | `class Skill(ABC): frontmatter,get_content,get_resource,get_script` + `SkillFrontmatter` + `SkillResource/SkillScript` | `python/packages/core/agent_framework/_skills.py:491` / `_skills.py:556` / `_skills.py:76` / `_skills.py:260` |
| Class-based skill extensibility | `class ClassSkill(Skill, ABC):` with `@ClassSkill.resource` / `@ClassSkill.script` auto-discovery via `_discover_marked_members` | `python/packages/core/agent_framework/_skills.py:1031` / `_skills.py:1127` / `_skills.py:985` |
| Skills loading / discovery | `class SkillsProvider(ContextProvider):` (1719) + `class SkillsSource(ABC)` + `class FileSkillsSource(SkillsSource):` scanning `SKILL.md` directories | `python/packages/core/agent_framework/_skills.py:1719` / `_skills.py:2360` / `_skills.py:2381` |
| MCP external tool extension | `class MCPTool:` base + `MCPStdioTool` / `MCPStreamableHTTPTool` / `MCPWebsocketTool` with `connect()`, `load_tools()`, allowlist `_MCP_FRAMEWORK_DENYLIST` | `python/packages/core/agent_framework/_mcp.py:333` / `_mcp.py:88` |
| Workflow executor extension | `class Executor(RequestInfoMixin,DictConvertible):` with `@handler` discovery, `input_types`/`output_types`/`workflow_output_types` | `python/packages/core/agent_framework/_workflows/_executor.py:30` / `_workflows/_executor.py:529` |
| Evaluator extension protocol | `@runtime_checkable class Evaluator(Protocol): evaluate(items, eval_name)->EvalResults` + `LocalEvaluator` + `@evaluator` function wrapper | `python/packages/core/agent_framework/_evaluation.py:683` / `_evaluation.py:1519` / `_evaluation.py:1408` |
| Compaction/tokenizer extension | `class CompactionStrategy(Protocol)` + `class TokenizerProtocol(Protocol)` pluggable per-agent/per-client | `python/packages/core/agent_framework/_compaction.py:51` / `_compaction.py:42` |
| Security/FIDES extension (experimental) | `SecureAgentConfig(ContextProvider)` injecting tools/instructions/middleware; `Security` middleware with label tracking | `python/packages/core/agent_framework/security.py:1` |
| Feature-stage marking | All newer extension points gated by `@experimental(feature_id=ExperimentalFeature.SKILLS/HARNESS/EVALS/PROGRESSIVE_TOOLS/FIDES...)` with `ExperimentalWarning` | `python/packages/core/agent_framework/_feature_stage.py:43` / `_feature_stage.py:383` / `_skills.py:76` |
| Public export surface | `agent_framework/__init__.py` re-exports every extension type (`AgentMiddleware`, `FunctionTool`, `ContextProvider`, `Skill`, `Executor`, `Evaluator`, etc.) | `python/packages/core/agent_framework/__init__.py:149` |
| Provider package isolation via workspace | `python/pyproject.toml` declares 18 workspace members (`agent-framework-core`, `agent-framework-openai`, `agent-framework-foundry`, etc.) — third party can add package following cookiecutter template | `python/pyproject.toml:71` |
| Lifecycle hooks | `InMemoryHistoryProvider.before_run` loads `state["messages"]`; `after_run` persists; `Executor.on_checkpoint_save/on_checkpoint_restore` for stateful executors | `python/packages/core/agent_framework/_sessions.py:506` / `_sessions.py:518` / `_workflows/_executor.py:493` / `_workflows/_executor.py:508` |
| Test evidence (tools) | `python/packages/tools/tests/test_shell_truncate_and_quote.py` exists; middleware/skill/tool tests under `python/packages/core/tests` (not enumerated here) | `python/packages/tools/tests/test_shell_truncate_and_quote.py:1` |

## Answers to Dimension Questions

**1. What can be extended via plugins?**

- **Tools:** Any `Callable` → `FunctionTool` via `@tool` (`python/packages/core/agent_framework/_tools.py:1179`) or `FunctionTool(name, func, input_model)` (`_tools.py:243`); host-provided tools via `MCPTool` family (`_mcp.py:333`). Third party adds new tool type without core change, but hosted shell/code-interpreter tool kinds still require a client implementing `SupportsShellTool`/`SupportsCodeInterpreterTool` (`_clients.py:816`).
- **Agents/Memory/Context:** Subclass `ContextProvider` and override `before_run`/`after_run` to inject messages/tools/instructions/middleware per-run via `SessionContext.extend_messages/extend_tools/extend_middleware` (`_sessions.py:351`, `_sessions.py:220`, `_sessions.py:281`). Built-ins: `InMemoryHistoryProvider`, `FileHistoryProvider`, `MemoryContextProvider`, `FileMemoryProvider`, `TodoProvider`, `SkillsProvider`.
- **Providers/Chat Clients:** Implement `SupportsChatGetResponse` protocol or subclass `BaseChatClient` and override `_inner_get_response(messages, stream, options)` (`_clients.py:82`/`_clients.py:214`). Each provider lives in its own workspace package (`python/pyproject.toml:71`) and is lazy-loaded via `__getattr__` (`azure/__init__.py:27`). Same for embedding clients (`BaseEmbeddingClient`, `_clients.py:919`) and `SupportsGetEmbeddings`.
- **Middleware:** Implement `AgentMiddleware`/`ChatMiddleware`/`FunctionMiddleware` and pass via `Agent(middleware=[...])` or `client_kwargs["middleware"]` (`_middleware.py:465`/`524`/`588`). Functional form via `@agent_middleware` etc. (`_middleware.py:678`). `FunctionInvocationContext` also enables progressive tool disclosure (`add_tools`/`remove_tools`, `_middleware.py:287`).
- **Workflows:** Subclass `Executor` + `@handler` (`_workflows/_executor.py:30`/`529`) or `@executor` function wrapper; builder `WorkflowBuilder` wires edges (`_workflows/_workflow_builder.py:53`).
- **Skills:** Implement `Skill` (file-based `FileSkill`, code-based `InlineSkill`, class-based `ClassSkill`) with `SkillResource`/`SkillScript` and provide via custom `SkillsSource` or `SkillsProvider` (`_skills.py:491`/`1719`/`2360`).
- **Evals:** Implement `Evaluator` protocol (`_evaluation.py:683`) or wrap a function with `@evaluator` (`_evaluation.py:1408`); built-ins `LocalEvaluator`, `keyword_check`, `tool_called_check`.
- **Compaction/Tokenizer/Security:** `CompactionStrategy`, `TokenizerProtocol` (`_compaction.py:42`/`51`); `SecureAgentConfig` (`security.py:1`).
- **UI/Infra:** Not directly plugin-extensible in core — `devui`, `ag_ui`, `a2a`, `foundry_hosting` are separate workspace packages exposing their own extensions.

**2. Can plugins be loaded at runtime?**

Partially, but not via a unified dynamic plugin loader:

- **No generic `PluginLoader`.** No `entry_points`, YAML plugin manifest, or filesystem plug-in directory with hot reload was found. `grep entry_points|plugin.*load` returned only provider lazy imports. Configuration is code-first: `Agent(tools=[...], middleware=[...], context_providers=[...])`.
- **Dynamic surfaces that exist:** (a) `FileSkillsSource` discovers `SKILL.md` files at runtime by scanning configured directories (`_skills.py:2381`, `_skills.py:3080`); (b) `MCPTool.connect()` dynamically loads remote tool specs and synthesizes `FunctionTool`s (`_mcp.py:919`); (c) provider packages are lazy-imported via `__getattr__` so `import agent_framework.openai` triggers `importlib.import_module` only on first attribute access (`openai/__init__.py:14`). (d) `normalize_tools` flattens arbitrary iterables/mappings at run time (`_tools.py:941`).
- **No hot-swap/versioning.** Middleware and providers are fixed at `Agent` construction; replacing them requires constructing a new `Agent`. `SkillsProvider` loads skills lazily per `before_run` but caches `SkillFrontmatter`; updating a `SKILL.md` on disk requires re-instantiating the source.
- **Verdict:** Runtime extensibility is “import-time” dynamic, not “drop-in plugin file” dynamic.

**3. Are plugins isolated from each other?**

No. All extensions share the same OS process, event loop, and mutable framework objects:

- Middleware runs sequentially in the same async task via `Pipeline.execute` — one failing middleware can abort the chain or corrupt `context.result` (`_middleware.py:883`).
- `FunctionInvocationContext.tools` is a live mutable list shared across all function middleware and tool invocations; `add_tools` mutates in-place with only duplicate-name guard (`_middleware.py:321`).
- `AgentSession.state` is a single `dict[str,Any]` partitioned by `source_id`, but any provider can read/write `session.state` or `SessionContext.metadata` (`_sessions.py:772`). No capability restriction, no sandbox, no process isolation.
- `MCPTool` spawns subprocesses (stdio) or HTTP connections, but the `FunctionTool`s it synthesizes still execute in the agent process; no VM/sandbox for skill scripts except the optional `HyperlightCodeActProvider` (experimental, separate package).
- Security `FIDES` labeling is logical (integrity/confidentiality labels on `Content`) not isolation (`security.py:1`).
- Failure modes: a misbehaving `HistoryProvider.save_messages` can drop history; a tool with `max_invocations=None` can loop `DEFAULT_MAX_ITERATIONS=40` (`_tools.py:90`); untrusted skill resources could inject prompt instructions (mitigated only by XML-escaping, `_skills.py:693`).

**4. Are extension points documented and stable?**

Documented via docstrings, not stable:

- Every public extension type has a Google-style docstring and is re-exported in `python/packages/core/agent_framework/__init__.py:149` with `__all__`.
- However, most newer surfaces are `@experimental` and emit `ExperimentalWarning` on first use (`_feature_stage.py:383`): `SKILLS` (`_skills.py:76`), `HARNESS`/`FILE_HISTORY`/`PROGRESSIVE_TOOLS`/`MCP_LONG_RUNNING_TASKS`/`FIDES`/`EVALS`/`FUNCTIONAL_WORKFLOWS`. The `__feature_id__` / `__feature_stage__` attributes are explicitly “optional stage metadata and may disappear” (`_feature_stage.py:43` note).
- No `STABLE`/`ReleaseCandidate` promotion was observed for Skills/Harness/Evals; `ReleaseCandidateFeature` enum is empty (`_feature_stage.py:68`).
- Versioned via `agent-framework-core==1.9.0` pinned in `python/pyproject.toml:26`; provider packages version-lockstep rather than semver range, so extension breakage requires coordinated bump.
- Host-level docs (`docs/specs/`, `docs/decisions/*.md`) describe intent (e.g., `docs/decisions/0007-agent-filtering-middleware.md:13`, `docs/decisions/0021-provider-leading-clients.md:17`) but no `EXTENDING.md` or plugin-author guide enumerates the extension contract or compatibility policy.

## Architectural Decisions

- **Protocol-first, inheritance-optional (`_clients.py:82`, `_agents.py:188`):** `SupportsChatGetResponse` and `SupportsAgentRun` are `@runtime_checkable Protocol`s, enabling duck-typed third-party clients/agents without inheriting framework base classes. Tradeoff: static type checking is weaker for duck-typed clients; `isinstance` checks rely on method presence. Evidence: docstring “Any class implementing the required methods is compatible. No need to inherit.” (`_agents.py:189`).
- **Workspace per-provider, lazy `__getattr__` (`python/pyproject.toml:71`, `azure/__init__.py:27`):** Keeps `agent-framework-core` dependency-free (no `litellm`, `mcp`, etc.) and lets third parties publish `agent-framework-myprovider`. Tradeoff: provider discovery is not automatic — users must `pip install` and `import`; no entry-point auto-registration.
- **Middleware as pipeline wrapper, not aspect annotation (`_middleware.py:843`):** Three distinct pipelines (Agent/Chat/Function) adapt via `MiddlewareWrapper` rather than decorators on methods. Enables composable cross-cutting concerns (logging, retry, compaction) without modifying agents/clients. Tradeoff: three pipeline types + `categorize_middleware` increases API surface; mis-tagged `AgentMiddleware` added via `ContextProvider.extend_middleware` raises `MiddlewareException` (`_sessions.py:302`).
- **ContextProvider as the harness extension point (`_sessions.py:351`):** All “memory, skills, harness, eval” contributions funnel through the same `before_run`/`after_run` contract keyed by `source_id`. This unifies attribution and enables multi-provider composition (e.g., `create_harness_agent` composes `TodoProvider` + `MemoryProvider`). Tradeoff: `ContextProvider` is async-only and state is untyped `dict[str,Any]`, so providers can collide on `source_id` or pollute `session.state`.
- **Skills via filesystem + code + custom sources (`_skills.py:2360`):** `SkillsSource` abstraction allows plugging custom backends (REST, DB) while default `FileSkillsSource` scans disk. Tradeoff: discovery is synchronous walk plus async read, with no built-in watching/reload; path-traversal guards (`_skills.py:2765`) mitigate but not sandbox content.
- **MCP as external tool extension with allowlist (`_mcp.py:88`):** Server-declared `inputSchema.properties` form the allowlist; `additional_tool_argument_names` is the only escape hatch, explicitly opt-in (`_mcp.py:180`). `"*"` global key applies to every tool. Tradeoff: correctly defends against framework-kwargs leakage (`chat_options`, `session`, `_meta`, etc.) but adds per-tool configuration complexity.

## Notable Patterns

- **Progressive disclosure for Skills (`_skills.py:3`, `_skills.py:1694`):** Advertise (`system` prompt injection of `frontmatter` names/descriptions) → Load (`load_skill` tool returns full `SKILL.md`) → Read resources/scripts on demand. Mirrors the Agent Skills spec (`agentskills.io`). Evidence: `SkillsProvider` docstring at `_skills.py:15`.
- **Progressive tool exposure during invocation (`_middleware.py:287`):** A tool that receives `FunctionInvocationContext` can call `ctx.add_tools([factorial])`; change takes effect on the *next* model iteration via `context.tools[:] = merged` (`_middleware.py:322`). Enables tool-agents that unlock capabilities mid-run.
- **Decorator auto-registration for ClassSkill (`_skills.py:1127`, `_skills.py:985`):** `@ClassSkill.resource` / `@ClassSkill.script` stamps `_skill_resource_marker`; `_discover_marked_members` walks the MRO including `property` descriptors, caching `InlineSkillResource` on first access (`_skills.py:1252`). Keeps subclass constructors stable — new `SkillFrontmatter` fields require no subclass changes.
- **Executor compile-time handler validation (`_workflows/_executor.py:329`, `_workflows/_executor.py:694`):** `__init__` calls `_discover_handlers` + `validate_workflow_context_annotation`; duplicate message types or missing `@handler` raise at construction, not at run.
- **Per-tool invocation budget (`_tools.py:301`):** `max_invocations` and `max_invocation_exceptions` per `FunctionTool` instance, counted locally, plus global `FunctionInvocationConfiguration["max_function_calls"]` for per-request cap. No reset across requests — intentional to surface long-running singleton misuse in docstring note (`_tools.py:329`).

## Tradeoffs

- **Explicit interfaces vs. runtime plugin registry:** Giving every extension an `ABC`/`Protocol` yields IDE autocomplete and mypy/pyright coverage, but there is no `register_plugin(name, factory)` or `load_plugins_from_entry_points()` — every extension must be wired in code.
- **Flexible flattening via `normalize_tools` (`_tools.py:941`) vs. strict typing:** Accepts `ToolTypes | Callable | Sequence | Mapping.tools | Iterable[Tool]` and spreads toolbox wrappers, which lets `tools=[toolbox, my_func]` work naturally but also silently accepts unexpected iterables; wrong input produces a `log.warning("Can't parse tool.")` rather than a typed error (`_tools.py:1036`).
- **Shared mutable state vs. isolation:** `SessionContext.tools` / `session.state` mutability enables progressive disclosure and session resumption (`AgentSession.to_dict/from_dict`, `_sessions.py:779`/`793`), but any extension can corrupt cross-provider data; no schema validation on `state`.
- **Experimental velocity vs. stability:** Gating new areas behind `@experimental` lets the team iterate without semver, but third-party authors cannot rely on `Skills`, `Evals`, `Harness`, `FIDES`, or `ProgressiveTools` remaining compatible; no deprecation policy for experimental removal.
- **Lazy provider import vs. discoverability:** `__getattr__` lazy loading keeps import cheap (`openai/__init__.py:14`), but `pip list` does not reveal available providers, and `dir(agent_framework)` hides them until first access.

## Failure Modes / Edge Cases

- **Duplicate tool name collision:** Adding two different tools with the same `name` raises `ValueError: Duplicate tool name` via `_raise_duplicate_tool_name` (`_tools.py:896`, `_tools.py:927`). Same-object duplicate is deduped silently (`_tools.py:924`). Easy to hit when composing `FileSkillsSource` + `InlineSkill` with overlapping names or when `tool_name_prefix` is not used for MCP (`_mcp.py:388`).
- **No isolation failure — cascading crash:** An exception in any `FunctionMiddleware.process` or `FunctionTool.invoke` propagates through `FunctionMiddlewarePipeline.execute` (`_middleware.py:998`) and aborts the entire function-calling loop; retries must be implemented as custom `AgentMiddleware` (e.g., `RetryMiddleware` example in `_middleware.py:482`).
- **Silent `FileHistoryProvider` path safety edge:** `session_id` with `../`, control characters, or Windows-reserved stems (`CON`, `PRN`) is base64-encoded via `_ENCODED_SESSION_PREFIX` and then `is_relative_to(self._storage_root)` checked (`_sessions.py:1089`, `_sessions.py:1122`); encoded names can collide if not `urlsafe_b64encode`-roundtripped, and cross-process locking uses 64-stripe `threading.Lock` + per-loop `asyncio.Lock` (`_sessions.py:916`, `_sessions.py:991`) without cross-host coordination.
- **MCP long-running task abandonment:** If `call_tool_as_task` exceeds `max_task_wait` (default `None` = unbounded) or the server drops connection before `task_id` is known, the framework raises `ToolExecutionException` without retrying the `tools/call` to avoid double-execution (`_mcp.py:231` comment, `_mcp.py:243` `_MCPTaskAbandoned`).
- **Progressive-tools misuse:** Calling `ctx.add_tools` outside a live function-calling loop raises `RuntimeError: not bound to a live agent run` (`_middleware.py:312`); recommendations grounded in “batch applied on next iteration” mean a tool must not assume its added tools are available in the same turn.
- **Experimental removal risk:** Any class decorated `@experimental(SKILLS|HARNESS|EVALS|FIDES|PROGRESSIVE_TOOLS|MCP_LONG_RUNNING_TASKS)` (`_skills.py:76`, `security.py:77`, `_evaluation.py:68`, `_middleware.py:287`) can be renamed or removed without deprecation; `getattr(obj, "__feature_stage__")` may disappear (`_feature_stage.py:372` comment).
- **Middleware classification error:** `ContextProvider.extend_middleware` calls `categorize_middleware` and raises `MiddlewareException("Context providers may only add chat or function middleware.")` if `AgentMiddleware` is included (`_sessions.py:302`); easy to hit when a provider tries to add agent-loop control.
- **Provider not installed:** Accessing `agent_framework.anthropic.AnthropicChatClient` without `agent-framework-anthropic` raises `ModuleNotFoundError` with `pip install agent-framework-anthropic` hint (`openai/__init__.py:50` pattern mirrored in all provider `__getattr__`s).

## Future Considerations

- Fill the “plugin manager” gap with an optional `PluginRegistry` that discovers entry points (`importlib.metadata.entry_points(group="agent_framework.providers")`) and loads/unloads providers without code changes — keep it optional to preserve lightweight-core principle.
- Stabilize the experimental surfaces under a documented lifecycle: promote `Skills` and `HistoryProvider` from `@experimental` to `@release_candidate` with a 2-version deprecation window rather than silent removal.
- Add isolation primitives: run `InlineSkillScript` and file-based `SkillScript` via `SkillScriptRunner` in an opt-in sandbox (e.g., `HyperlightCodeActProvider`-style microVM) and gate resource reads with `allowed_tools` / filesystem `allowlist` beyond XML-escaping.
- Version `source_id` namespaces and add `session.state` JSON-schema validation so one misbehaving provider cannot corrupt another’s state.
- Publish an `EXTENDING.md` with per-extension-point examples, stability tier table, and a `third-party-tool-type` cookbook showing how to add a new tool kind without modifying core (e.g., implement `SupportsCustomTool` protocol and register alongside `FunctionTool`).

## Questions / Gaps

- No evidence found for a unified plugin lifecycle (install, enable/disable, update, uninstall, version-negotiation, health-check, or crash-isolation) — searched `plugin`, `extension`, `registry`, `hook`, `entry_points` across `python/packages/core/agent_framework/` and found only provider `__getattr__` lazy imports and `FileSkillsSource._discover_skill_directories`.
- No evidence found for runtime hot-reload of skills/tools/middleware — `FileSkillsSource` scans once per provider instantiation and caches; file watching (`watchdog`/`inotify`) mechanism absent.
- No evidence found for plugin isolation or permission model — `Security` module implements information-flow labels but no sandbox/process boundary for extensions; verified by absence of `subprocess`, `hyperlight`, or container logic in core (only `packages/hyperlight` optional package references it).
- Cross-language plugin story unclear — `.dotnet` side has its own middleware/skill abstractions (not inspected in this task, isolated by hard rule); whether a plugin written for Python is portable was not assessable.
- `lab` (`python/packages/lab`) and `declarative` (`python/packages/declarative`) extensions are experimental and lack stable plugin-author docs; whether they will become supported extension points is unknown.

---

Generated by `21.01-plugin-and-extension-points` against `agent-framework`.
