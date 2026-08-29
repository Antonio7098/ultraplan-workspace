# Source Analysis: pydantic-ai

## Plugin and Extension Points

### Source Info

| Field | Value |
|-------|-------|
| Name | pydantic-ai |
| Path | `studies/agent-harness-study/sources/pydantic-ai` |
| Language / Stack | Python (uv workspace: `pydantic-ai-slim`, `pydantic-graph`, `pydantic-evals`) |
| Analyzed | 2026-08-28 |

## Summary

Pydantic AI has no formal "plugin" registry or dynamic loader. Instead it exposes a small set of deeply documented, typed abstract interfaces that together form the extension system: **Capabilities** (`AbstractCapability`) for agent-scoped behavior (instructions, tools, hooks, model selection), **Toolsets** (`AbstractToolset` + `WrapperToolset` family) for tool collections, **Models** (`AbstractModel`/`Model`) + **Providers** (`Provider`) for LLM backends, **Native Tools** (`AbstractNativeTool`) for provider-native features, **MCP** via `MCPToolset`, **UI Adapters** (`UIAdapter`) for frontend protocols, and **Embedding Models** (`EmbeddingModel`). Composition is via `Agent(capabilities=[...], toolsets=[...], model=...)` and per-run `capabilities`/`toolsets`. The middleware model (`CombinedCapability` chaining, `CapabilityOrdering`, two-phase `for_agent` binding) is explicit, tested, and serializable via `AgentSpec`/`CAPABILITY_TYPES`. Isolation is logical (per-run `for_run` copies, `WrapperToolset` boundaries), not sandboxed; runtime discovery is limited to spec-driven declarative loading and on-demand `defer_loading`/`load_capability` plus `ToolSearch`, not OS-level plugin discovery.

## Rating

**7 / 10** — Clear, stable extension model with explicit interfaces, lifecycle hooks, and operational safeguards (ordering, two-phase binding, per-run isolation, spec serialization). Downgraded from 9 because there is no dynamic plugin loading (no entrypoints/importlib discovery), no inter-plugin sandboxing/isolation, and UI/evals/provider additions still require code changes and dependency opt-ins rather than hot-pluggable manifests.

## Evidence Collected

Every entry MUST include a file path with line numbers. Format: `path/to/file.ts:NN`.

| Area | Evidence | File:Line |
|------|----------|-----------|
| Extension interfaces — Capability | `AbstractCapability` ABC with `get_toolset`, `get_native_tools`, `get_instructions`, `get_model`, `resolve_model_id`, `get_wrapper_toolset`, `prepare_tools`, 30+ lifecycle hooks (`before/after/wrap/on_*` for run, node, model request, tool validate/execute, output validate/process, event stream) | `pydantic_ai_slim/pydantic_ai/capabilities/abstract.py:162` |
| Extension interfaces — Toolset | `AbstractToolset` ABC with `id`, `get_tools`, `call_tool`, `for_run`, `for_run_step`, `get_instructions` + factory methods (`filtered`, `prefixed`, `prepared`, `approval_required`, `defer_loading`, `with_metadata`) | `pydantic_ai_slim/pydantic_ai/toolsets/abstract.py:76` |
| Extension interfaces — Toolset wrappers | `WrapperToolset` base + concrete wrappers (`FilteredToolset`, `PrefixedToolset`, `PreparedToolset`, `ApprovalRequiredToolset`, `DeferredLoadingToolset`, `CapabilityOwnedToolset`, `ToolSearchToolset`) | `pydantic_ai_slim/pydantic_ai/toolsets/wrapper.py:15`, `pydantic_ai_slim/pydantic_ai/toolsets/__init__.py:1` |
| Extension interfaces — Model/Provider | `AbstractModel` (system/model_name/label) → `Model` (prepare_request, request, request_stream, customize_request_parameters, profile) and `Provider` (name/base_url/client/model_profile) with `infer_provider_class`/`infer_provider` factory | `pydantic_ai_slim/pydantic_ai/models/_abstract.py:19`, `pydantic_ai_slim/pydantic_ai/models/__init__.py:366`, `pydantic_ai_slim/pydantic_ai/providers/__init__.py:42` |
| Extension interfaces — Native tools | `AbstractNativeTool` with `__init_subclass__` registry `NATIVE_TOOL_TYPES` + discriminator `SUPPORTED_NATIVE_TOOLS` and typed subclasses (`WebSearchTool`, `XSearchTool`, `CodeExecutionTool`, etc.) | `pydantic_ai_slim/pydantic_ai/native_tools/__init__.py:35` |
| Extension interfaces — UI adapter | `UIAdapter[RunInputT,MessageT,EventT]` abstract with `build_run_input`, `load_messages`, `build_event_stream`, `toolset`, `state`, `deferred_tool_results` + concrete `AGUIAdapter`/`VercelAIAdapter` pattern | `pydantic_ai_slim/pydantic_ai/ui/_adapter.py:208` |
| Extension interfaces — Embeddings | `EmbeddingModel` ABC with `embed`, `model_name`, `system` | `pydantic_ai_slim/pydantic_ai/embeddings/base.py:8` |
| Plugin loader — declarative spec | `AgentSpec.from_file`/`Agent.from_spec` + `CAPABILITY_TYPES` registry and `custom_capability_types` parameter; `AgentSpec` deserializes capabilities via `get_serialization_name`/`from_spec` | `pydantic_ai_slim/pydantic_ai/agent/spec.py:33`, `pydantic_ai_slim/pydantic_ai/capabilities/__init__.py:72` |
| Plugin loader — absence of dynamic discovery | No `entry_points`, `importlib.metadata`, or plugin manifest loader found; grep for `plugin|entrypoint|entry.point` returns only docs/mkdocs noise | `grep plugin|entrypoint` (no hits in `pydantic_ai_slim/`) |
| On-demand loading | `AbstractCapability.defer_loading` + `id` + `load_capability` typed tool (`LoadCapabilityCallPart`/`LoadCapabilityReturnPart`) + `ToolSearch` capability + `DeferredLoadingToolset` and `ToolDefinition.defer_loading`/`ToolSearchToolset` gating | `pydantic_ai_slim/pydantic_ai/capabilities/abstract.py:216`, `pydantic_ai_slim/pydantic_ai/_deferred_capabilities.py:28`, `pydantic_ai_slim/pydantic_ai/capabilities/_tool_search.py:1`, `pydantic_ai_slim/pydantic_ai/tools.py:404` |
| Tool extension without core modification | `Tool`/`FunctionToolset` (`@tool`, `@tool_plain`, `Tool.prepare`, `ToolsPrepareFunc`) + third-party capability example `pydantic-ai-chdb` documented as external pip package | `pydantic_ai_slim/pydantic_ai/tools.py:292`, `pydantic_ai_slim/pydantic_ai/toolsets/function.py:113`, `docs/capabilities/third-party.md:47` |
| Lifecycle management | Two-phase binding `bind_capabilities_tier(agent, innermost=False/True)` + `for_agent` (construction) and `for_run` (per-run copy) + `CombinedCapability.for_run` rebind via `replace_no_init` + `get_wrapper_toolset` per-run + `_ctx_for_active_cap` gating by `loaded_capability_ids` | `pydantic_ai_slim/pydantic_ai/capabilities/combined.py:857`, `pydantic_ai_slim/pydantic_ai/agent/__init__.py:745` |
| Isolation — per-run state | `for_run` returns fresh instance; `CombinedCapability.for_run` gathers new children atomically; docs note per-run isolation via `for_run` | `pydantic_ai_slim/pydantic_ai/capabilities/abstract.py:309`, `pydantic_ai_slim/pydantic_ai/capabilities/combined.py:101` |
| Isolation — (non-)sandboxing | No sandbox/permission boundary; `WrapperToolset.call_tool` delegates directly; MCP transport is the only network isolation layer; `anyio.Lock` used for lifecycle, not isolation | `pydantic_ai_slim/pydantic_ai/toolsets/wrapper.py:15`, `pydantic_ai_slim/pydantic_ai/mcp.py:863` |
| Ordering/composition | `CapabilityOrdering` (`position`, `wraps`, `wrapped_by`, `requires`) + `sort_capabilities` topological sort in `CombinedCapability.__normalize_capabilities` | `pydantic_ai_slim/pydantic_ai/capabilities/abstract.py:116`, `pydantic_ai_slim/pydantic_ai/capabilities/combined.py:61` |
| Documentation — extension guides | `docs/capabilities/custom.md` (subclass guide), `docs/capabilities/overview.md` (catalog), `docs/toolsets.md`, `docs/native-tools.md`, `docs/extensibility.md` + `.agents/skills` | `docs/capabilities/custom.md:1`, `docs/capabilities/overview.md:1` |
| MCP extensibility | `MCPToolset` as `AbstractToolset` with `client: FastMCPClient`, transports (HTTP/SSE/stdio/in-process), `sampling_model`/`elicitation_handler`, `process_tool_call` hook, `cache_tools` | `pydantic_ai_slim/pydantic_ai/mcp.py:716` |
| Provider extensibility | 35+ providers in `pydantic_ai_slim/pydantic_ai/providers/` each exposing `Provider` + `ModelProfile` + `<Provider>Model` subclassing `OpenAIChatModel`/`Model` | `pydantic_ai_slim/pydantic_ai/providers/__init__.py:142` (glob `providers/*.py`) |
| Agent entry points | `Agent.__init__` accepts `tools`, `toolsets`, `capabilities`, `model` (string/Model/callable selector), `instructions`, `model_settings` with layered merging | `pydantic_ai_slim/pydantic_ai/agent/__init__.py:534` |

## Answers to Dimension Questions

### 1. What can be extended via plugins?

**Tools (without touching core): yes.** Any function becomes a tool via `Tool`/`FunctionToolset` and `@agent.tool` / `@capability.tool` (`pydantic_ai_slim/pydantic_ai/tools.py:292`, `pydantic_ai_slim/pydantic_ai/toolsets/function.py:113`). `ToolPrepareFunc` (`pydantic_ai_slim/pydantic_ai/tools.py:104`) and `ToolsPrepareFunc` allow per-step filtering. Third-party confirmation: `pydantic-ai-chdb` ships `ChDBCapability` + `ChDBToolset` as a pip package (`docs/capabilities/third-party.md:47`).

**Toolsets (collections): yes.** Subclass `AbstractToolset` (`pydantic_ai_slim/pydantic_ai/toolsets/abstract.py:76`) or use the fluent wrappers (`FilteredToolset`, `PrefixedToolset`, `PreparedToolset`, `ApprovalRequiredToolset`, `DeferredLoadingToolset` in `pydantic_ai_slim/pydantic_ai/toolsets/__init__.py:1`). `WrapperToolset` (`pydantic_ai_slim/pydantic_ai/toolsets/wrapper.py:15`) is the designated composition point.

**Capabilities (agent behavior): primary extension point.** Subclass `AbstractCapability` (`pydantic_ai_slim/pydantic_ai/capabilities/abstract.py:162`) to contribute toolsets, native tools, instructions, model settings, model selection/resolution, and 30+ lifecycle hooks. The convenience `Capability` dataclass (`pydantic_ai_slim/pydantic_ai/capabilities/capability.py:29`) bundles instructions+tools without subclassing. The harness (`pydantic-ai-harness`) and third-party packages deliver whole agents (e.g., `Coder`, `FileSystem`, `Memory`) as capabilities (`docs/capabilities/overview.md:23`).

**Models/Providers/evals/memory/prompts/policies/UI:**
- **Models:** subclass `Model`/`AbstractModel` (`pydantic_ai_slim/pydantic_ai/models/_abstract.py:19`, `pydantic_ai_slim/pydantic_ai/models/__init__.py:366`) and register via string `model="myprovider:my-model"` resolved through `Provider.resolve_model_id` or `AbstractCapability.resolve_model_id`.
- **Providers:** subclass `Provider` (`pydantic_ai_slim/pydantic_ai/providers/__init__.py:42`) and extend `infer_provider_class`.
- **Native tools / provider policy:** subclass `AbstractNativeTool` (auto-registered in `NATIVE_TOOL_TYPES` at `pydantic_ai_slim/pydantic_ai/native_tools/__init__.py:108`) and advertise via `ModelProfile.supported_native_tools`.
- **Embeddings:** subclass `EmbeddingModel` (`pydantic_ai_slim/pydantic_ai/embeddings/base.py:8`).
- **MCP / external systems:** add an `MCPToolset` (`pydantic_ai_slim/pydantic_ai/mcp.py:716`) pointing at any URL/stdio/in-process server; no core change.
- **UI:** subclass `UIAdapter` (`pydantic_ai_slim/pydantic_ai/ui/_adapter.py:208`) — only two concrete adapters ship (AG-UI, Vercel AI), so the surface is proven but narrow.
- **Evals:** separate `pydantic_evals` package; not a capability but a distinct extension dimension (`pydantic_evals/`).
- **Prompts/policies:** via capabilities (`get_instructions`, `ProcessHistory`, `Guardrails` in harness, `PrepareTools`).

Limit: adding a new *provider-native* feature that has no `AbstractNativeTool` abstraction still requires a core `native_tools/*.py` edit; thin provider-specific tools without cross-provider semantics are documented to stay as `NativeTool` wrapped natives until a cross-provider abstraction emerges (`pydantic_ai_slim/pydantic_ai/native_tools/AGENTS.md:1`).

### 2. Can plugins be loaded at runtime?

**No dynamic discovery loader.** The codebase contains no `importlib.metadata.entry_points`, plugin manifest, or filesystem plug-in directory scan (grep for `plugin|entrypoint` is clean inside `pydantic_ai_slim/`). Extensions are loaded by Python import and passed to `Agent(capabilities=[MyCap()], toolsets=[my_toolset], model=...)` or to `agent.run(capabilities=[...])`.

**Declarative runtime loading exists via specs:** `AgentSpec.from_file(path)` / `Agent.from_spec(spec, custom_capability_types=[...])` (`pydantic_ai_slim/pydantic_ai/agent/spec.py:33`, `pydantic_ai_slim/pydantic_ai/capabilities/__init__.py:72`) deserializes a YAML/JSON file into an `Agent` with capabilities resolved through `CAPABILITY_TYPES` (`pydantic_ai_slim/pydantic_ai/capabilities/__init__.py:72`) by `get_serialization_name`/`from_spec`. This is data-driven, not hot-pluggable — you must restart the process with a new spec.

**On-demand reveal at runtime (not hot-load):** `AbstractCapability(defer_loading=True, id="...")` stays hidden until the model calls the built-in `load_capability` tool (`pydantic_ai_slim/pydantic_ai/_deferred_capabilities.py:48`). `ToolSearch` (`pydantic_ai_slim/pydantic_ai/capabilities/_tool_search.py`) and `DeferredLoadingToolset` (`pydantic_ai_slim/pydantic_ai/toolsets/deferred_loading.py:15`) implement search-gated tool reveal (`ToolDefinition.defer_loading` at `pydantic_ai_slim/pydantic_ai/tools.py:404`, `ModelRequestParameters.tool_visibility` at `pydantic_ai_slim/pydantic_ai/models/__init__.py:161`). History-derived `loaded_capability_ids`/`discovered_tool_names` (`pydantic_ai_slim/pydantic_ai/_run_context.py:228`) gate subsequent steps. These reduce prompt cost, not enable loading new code from disk mid-run.

### 3. Are plugins isolated from each other?

**Logical isolation, not sandbox isolation.** Each `AbstractCapability.for_run(ctx)` is expected to return a fresh per-run copy for state isolation (`pydantic_ai_slim/pydantic_ai/capabilities/abstract.py:309`), and `CombinedCapability.for_run` rebuilds via `replace_no_init` (`pydantic_ai_slim/pydantic_ai/capabilities/combined.py:101`). `CapabilityOrdering` topologically sorts middleware (`pydantic_ai_slim/pydantic_ai/capabilities/abstract.py:116`) so execution order is deterministic and documented.

However:
- All capabilities/toolsets/models run in the **same Python process and event loop**; a capability's `wrap_run`/`wrap_model_request` can observe or suppress errors from peers, but there is no memory, permission, or network sandbox between them.
- `WrapperToolset.call_tool` (`pydantic_ai_slim/pydantic_ai/toolsets/wrapper.py:15`) delegates directly to the wrapped tool; a slow or leaking `call_tool` affects sibling toolsets.
- `MCPToolset` (`pydantic_ai_slim/pydantic_ai/mcp.py:863`) is the only boundary that executes outside the host process (via FastMCP transport / subprocess / HTTP).
- `durable_exec` (`pydantic_ai_slim/pydantic_ai/durable_exec/`) wraps leaf toolsets by `id` for replay, but still shares the process model.
- When `defer_loading=True`, inactive capabilities are fully skipped via `_ctx_for_active_cap` (`pydantic_ai_slim/pydantic_ai/capabilities/combined.py:880`), which is the closest to isolation (unloaded code never runs), but once loaded it joins the same chain.

No mechanism caps CPU, memory, or I/O per extension; no inter-plugin communication bus beyond the capability chain ordering.

### 4. Are extension points documented and stable?

**Documented: yes, extensively.** `docs/capabilities/custom.md:1` is a full subclassing guide (construction, typing, tools, `for_agent`/`for_run` lifecycle, 30+ hooks). `docs/capabilities/overview.md:1` catalogs every first-party and harness capability with Package column. `docs/toolsets.md`, `docs/native-tools.md`, `docs/extensibility.md:45`, and `docs/durable_execution/overview.md` cover the other surfaces. `.agents/skills/building-pydantic-ai-agents` mirrors the same guidance for AI assistants. `AGENTS.md:56` states that "strong primitives and extension points" are the project's philosophy.

**Stability signals:**
- Abstract bases use `@abstractmethod` and versioned serialization: `get_serialization_name` / `from_spec` (`pydantic_ai_slim/pydantic_ai/capabilities/abstract.py:261`) and `CAPABILITY_TYPES` registry (`pydantic_ai_slim/pydantic_ai/capabilities/__init__.py:72`) pin wire names; tests enforce backward compatibility (`tests/test_agent_spec.py`, `tests/models/test_model_settings_support.py` per `pydantic_ai_slim/pydantic_ai/models/AGENTS.md`).
- `Agent` construction keeps raw `model` strings un-inferred through `for_agent` binding so resolvers can intercept before provider construction (`pydantic_ai_slim/pydantic_ai/agent/__init__.py:638`).
- `AbstractNativeTool.__init_subclass__` auto-registration is a stable contract (`pydantic_ai_slim/pydantic_ai/native_tools/__init__.py:107`); `Provider` contract is frozen via `gen_ai.system` notes (`pydantic_ai_slim/pydantic_ai/providers/__init__.py:70`).

**Weakness:** No semver stability promise or deprecation policy is emitted from the capability/toolset modules themselves; stability is inferred from docs and `PydanticAIDeprecationWarning` usage rather than an explicit guarantee file.

## Architectural Decisions

| Decision | Why it exists | Consequence |
|----------|---------------|-------------|
| **Capabilities over Agent kwargs** (`pydantic_ai_slim/pydantic_ai/capabilities/AGENTS.md:1`) — every cross-cutting concern is a capability | Prevent constructor bloat; keep `Agent.__init__` narrow while allowing composable behavior (instructions + tools + hooks in one unit) | Clear extension story, but discoverability depends on reading `capabilities/__init__.py` or `docs/capabilities/overview.md` rather than IDE autocomplete of `Agent` |
| **Two-phase `for_agent` binding** (`pydantic_ai_slim/pydantic_ai/capabilities/combined.py:857` + `pydantic_ai_slim/pydantic_ai/agent/__init__.py:745`) — non-innermost caps bind first, extract toolsets, then innermost (durability) caps wrap them | Lets durability wrappers see the complete `agent.toolsets` set without capabilities being able to mutate the set mid-binding | Newcomers who put toolsets on an innermost capability hit a silent "no contribution" failure documented only in `AGENTS.md` |
| **CombinedCapability as middleware chain** (`pydantic_ai_slim/pydantic_ai/capabilities/combined.py:46`) with reverse-order `after_*`/`on_*_error` and chain builders (`_make_run_wrap` etc.) | Mirrors ASGI/Express middleware intuition (first = outermost) so ordering is predictable without teaching a new mental model | Bugs arise when a capability's `for_run` returns a new instance and instance-identity `wraps`/`wrapped_by` refs no longer match (`abstract.py:138`) |
| **`WrapperToolset` as universal decorator** (`pydantic_ai_slim/pydantic_ai/toolsets/wrapper.py:15`) with 7 built-in wrappers | Keeps leaf `AbstractToolset` implementations focused on listing/calling; cross-cutting policy (filter, prefix, approval) composes without subclass explosion | Same policy must be re-applied per toolset unless a capability's `get_wrapper_toolset` does it centrally |
| **Spec-driven construction** (`pydantic_ai_slim/pydantic_ai/agent/spec.py:33`, `CAPABILITY_TYPES` at `capabilities/__init__.py:72`) | Enables YAML/JSON-defined agents and harness codegen without code changes; `custom_capability_types` keeps the registry open | Any capability that holds callables/Tool objects cannot be serialized (`Capability.get_serialization_name` returns `None` at `capabilities/capability.py:105`), so its "plugin" is code-only |
| **Provider-native tool deferral to `ModelProfile`** (`pydantic_ai_slim/pydantic_ai/models/__init__.py:484` + `native_tools/__init__.py:35`) | Allows the same logical tool (e.g., `WebSearchTool`) to render via provider-native or local fallback without caller branching | Provider-specific behavior leaks into `ModelProfile` booleans; adding a new native feature requires a profile flag + adapter `supported_tool_deferral_modes` handshake |

## Notable Patterns

- **Per-run copy via `for_run`**: every built-in wrapper/capability that holds mutable state returns `replace(self, ...)` in `for_run`, ensuring concurrent runs do not share cursors. Visible in `Instrumentation.for_run` (`pydantic_ai_slim/pydantic_ai/capabilities/instrumentation.py:142`) and leaf toolsets (`pydantic_ai_slim/pydantic_ai/toolsets/abstract.py:112`).
- **Deferred reveal as capability visibility gate**: `ToolDefinition.capability_id` + `ModelRequestParameters.deferred_capability_ids` (`pydantic_ai_slim/pydantic_ai/models/__init__.py:201`) + history-derived `loaded_capability_ids` (`pydantic_ai_slim/pydantic_ai/_run_context.py:228`) together decide whether a tool appears on the wire (`visible`/`deferred`/`withheld`/`via_history` at `pydantic_ai_slim/pydantic_ai/models/__init__.py:161`). This is a novel, cache-aware form of lazy plugin activation.
- **`CapabilityOwnedToolset` stamping**: `CombinedCapability.get_toolset` wraps each child toolset in `CapabilityOwnedToolset(capability=cap)` (`pydantic_ai_slim/pydantic_ai/capabilities/combined.py:214`), so deferred-loading and durability can attribute tools to their owner without polluting `AbstractToolset` itself.
- **`AbstractNativeTool` auto-registry**: `__init_subclass__` writes `NATIVE_TOOL_TYPES[cls.kind]=cls` (`pydantic_ai_slim/pydantic_ai/native_tools/__init__.py:108`), validated via `supports_tool_return_schema` / `supported_native_tools` intersection in `Model.profile` (`pydantic_ai_slim/pydantic_ai/models/__init__.py:884`).
- **UI as capability-injected middleware**: `UIAdapter.run_stream_native` injects `ReinjectSystemPrompt` per run (`pydantic_ai_slim/pydantic_ai/ui/_adapter.py:548`) and unions `DeferredToolRequests` into `output_type` when a frontend toolset exists, showing how UI plugins are just run-scoped capabilities, not a separate plugin host.

## Tradeoffs

| Tradeoff | Pro | Con |
|----------|-----|-----|
| Code-only extension (no entrypoints/manifests) vs auto-discovery | No import-time side effects, no supply-chain manifest parsing, fully type-checked; spec `custom_capability_types` keeps the registry explicit | Cannot drop a file in a plugins folder and be discovered; third-party capabilities must be installed and imported by name |
| One `AbstractCapability` to rule them all (50+ hooks) vs many narrow interfaces | A single unit can touch every phase (model selection, tool validate/execute, output, events) without needing five registrations | Subclasses carry a large API surface; mis-overriding `_has_wrap_node_run` vs `wrap_node_run` or `for_agent` vs `for_run` is a subtle source of bugs |
| Logical isolation only vs sandbox | Zero overhead; capabilities can share `RunContext` cheaply and wrap synchronously in-process | A misbehaving capability (leaked `asyncio.create_task` in `wrap_run` without shield/`CancelScope`) can keep the event loop alive or swallow `CancelledError` — docs warn at `docs/capabilities/custom.md:462` but do not enforce |
| Spec serializability requires `get_serialization_name` | Enables reproducible YAML-defined agents that round-trip through `AgentSpec` | Any callable, closure, or `AgentDepsT`-bound tool is non-serializable (`capabilities/capability.py:105` returns `None`), so not all extensions benefit |
| `innermost` durability tier being last to bind | Lets temporal/dbos/prefect wrappers observe the final toolset list | Capability authors who chose `innermost` to "run last" cannot contribute tools — a category error caught only by reading the tier docs |

## Failure Modes / Edge Cases

| Mode | Evidence | Impact |
|------|----------|--------|
| **Duplicate tool names across toolsets/capabilities** | `AbstractToolset` enforces uniqueness only via `ToolManager`'s call-time lookup (`pydantic_ai_slim/pydantic_ai/tool_manager.py:506`); `label`/`tool_name_conflict_hint` (`pydantic_ai_slim/pydantic_ai/toolsets/abstract.py:107`) suggests `PrefixedToolset` | Without prefixing, later `get_tools` entries silently shadow earlier ones; `CallToolsNode` may invoke the wrong function |
| **Deferred capability never loaded** | `DeferredCapabilityLoader` resolves `load_capability` id against `ctx.active_capability_ids` (`pydantic_ai_slim/pydantic_ai/toolsets/_deferred_capability_loader.py:34`); unknown id returns an error return part, not an exception | Model that hallucinates an id gets an error message in history rather than a crash, but may retry indefinitely |
| **`for_run` returns new instance with stale identity refs** | `CapabilityOrdering.wraps: CapabilityRef` docs warn that instance refs use `is` and break when `for_run` replaces the target (`pydantic_ai_slim/pydantic_ai/capabilities/abstract.py:138`) | Ordering silently falls back to list order; middleware executes in unexpected sequence |
| **Unawaited tasks in `wrap_run`** | `docs/capabilities/custom.md:462` documents that `wrap_run` must shield or gather child tasks; no runtime check enforces it | Cancelled runs leak tasks; durable execution replay may diverge because child side effects are not deterministic |
| **Native tool unsupported on chosen model** | `AbstractNativeTool.optional` (`pydantic_ai_slim/pydantic_ai/native_tools/__init__.py:80`) controls fallback; non-optional native tool on unsupported model raises `UserError` in `Model.prepare_request` (`pydantic_ai_slim/pydantic_ai/models/__init__.py:671`) | Run fails late (at first request) rather than at `Agent(...)` construction, unless caller tests the profile first |
| **Spec references unknown capability** | `AgentSpec` validates unknown names via registry lookup in `agent/spec.py:271` | `UserError` at `from_spec` time; no hot-fix without a new deployment that adds `custom_capability_types` |
| **History-replayed `load_capability` without matching toolset** | `_refresh_loaded_capability_ids` intersects `parse_loaded_capabilities(messages)` with `capabilities.keys()` (`pydantic_ai_slim/pydantic_ai/_agent_graph.py:2407`) | Stale capability ids in persisted history are ignored; deferred tools stay withheld rather than erroring |

## Future Considerations

| Consideration | Rationale |
|---------------|-----------|
| Add opt-in `importlib.metadata.entry_points(group="pydantic_ai.capabilities")` discovery for `AgentSpec` without requiring `custom_capability_types` | Would let third-party packages self-register (e.g., `pydantic-ai-chdb`) similar to the pattern used by `pytest`/`mkdocs` plugins, without breaking explicit-import semantics |
| Introduce a per-capability resource budget (timeout, task count, memory hint) and enforce it via `anyio.CapacityLimiter` already used for `max_concurrency` (`pydantic_ai_slim/pydantic_ai/agent/__init__.py:694`) | Mitigates the unguarded `wrap_run` task-leak and "noisy plugin" failure modes |
| Promote `for_run` isolation derive-check (derive `*_safe_at_runtime` from overridden hook introspection rather than a manual `ClassVar`) — tracked at `pydantic_ai_slim/pydantic_ai/capabilities/abstract.py:187` (`#5477`) | Lowers the barrier for third-party capabilities to be used as per-run overrides alongside durability |
| Publish a semver policy for `AbstractCapability`/`AbstractToolset` (e.g., `get_serialization_name` stability window) | Currently documented only via deprecation warnings; an explicit contract would let harness/third-party maintainers plan releases |
| Consider a `UIAdapter` plugin registry parallel to `CAPABILITY_TYPES` | UI is the only surface with just two adapters; a registry + `ui:spec` file would let teams add custom wire protocols without forking `pydantic_ai_slim` |

## Questions / Gaps

| Gap | What was searched | Status |
|-----|-------------------|--------|
| Is there an eval-memory-policy extension beyond harness? | Explored `pydantic_evals/`, `durable_exec/`, `docs/durable_execution/`, capability catalog; no core eval-memory plugin separate from harness | No evidence of a core `pydantic-ai` eval plugin beyond `pydantic_evals`; harness ships memory/eval-adjacent capabilities |
| Are there integration tests proving third-party `AbstractToolset`/`AbstractCapability` do not need core modification? | Checked `tests/test_capabilities.py`, `tests/test_toolsets.py`, `examples/`, `docs/capabilities/third-party.md` | Verified via `pydantic-ai-chdb` docs entry and `Tool` constructor usage in `docs/capabilities/custom.md:106`; no dedicated third-party isolation test suite in `tests/` |
| Does the framework enforce plugin dependency ordering (`requires`)? | Read `CapabilityOrdering.requires` (`pydantic_ai_slim/pydantic_ai/capabilities/abstract.py:157`) and `sort_capabilities` | Ordering declares `requires` but the sort raises `UserError` only if the required type is absent; it is not a soft dependency resolver |

---

Generated by `21.01-plugin-and-extension-points` against `pydantic-ai`.
