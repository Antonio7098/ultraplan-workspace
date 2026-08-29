# Source Analysis: letta

## 21.01 Plugin and Extension Points

### Source Info

| Field | Value |
|-------|-------|
| Name | letta |
| Path | `studies/agent-harness-study/sources/letta` |
| Language / Stack | Python 3.11+ / FastAPI, SQLAlchemy, Pydantic, Temporal, MCP |
| Analyzed | 2026-08-28 |

## Summary

Letta exposes a minimal, string-target plugin registry (`letta/plugins/plugins.py:28`) with exactly two slots — `experimental_check` and `summarizer` — resolved via `importlib.import_module` from `LETTA_PLUGIN_REGISTER` (`letta/settings.py:320`). The only typed contract is `SummarizerProtocol` (`letta/plugins/plugins.py:7`); the experimental checker is an untyped function. There is no entry-point discovery, no sandbox isolation, no versioned manifest, and no install/enable lifecycle beyond lazy `get_plugin` + cached globals with `reset_*` helpers. Broader extensibility (custom tools, MCP tools, file-type registry, LLM provider factory) exists but is implemented as hard-coded enums/factories and database-persisted code objects — not as pluggable interfaces. Adding a new *custom tool* requires no core change, but adding a new *tool type*, provider, memory or eval plugin does require modifying `ToolType` (`letta/schemas/enums.py:212`), `ToolExecutorFactory._executor_map` (`letta/services/tool_executor/tool_execution_manager.py:35`), and `LLMClient.create` (`letta/llm_api/llm_client.py:14`). Documentation is limited to a 22-line README with a stale schema example. Overall posture matches the 1-3 band: present but ad-hoc, unisolated, and unsafe for untrusted third-party code.

## Rating

**3 / 10 — Absent / ad-hoc**

Rationale: Two hard-coded plugin slots behind a semicolon-delimited env var are the entire plugin surface. One has a `Protocol`, one has none. Loading is dynamic (`importlib`) but discovery is manual strings, there is no lifecycle manager, no isolation, no capability/compatibility negotiation, and the README is inconsistent with the implementation. Tool extensibility via `ToolCreate`/`MCPTool` is richer but is persistence+factory-driven, not a plugin contract. Search for `entry_points`, marketplace manifests, plugin versioning, and sandboxing returns only incidental hits (`.gitignore` plugin lists, OTEL `extensions`).

## Evidence Collected

Every entry MUST include a file path with line numbers. Format: `path/to/file.ts:NN`.

| Area | Evidence | File:Line |
|------|----------|-----------|
| Extension interface — SummarizerProtocol | `@runtime_checkable class SummarizerProtocol` with `async def summarize(text: str) -> str` and `def get_name() -> str` | `letta/plugins/plugins.py:7-12` |
| Extension interface — experimental checker | No protocol (`"protocol": None`); impl is bare function `is_experimental_enabled(feature_name: str, **kwargs) -> bool` returning `False` | `letta/plugins/plugins.py:17-20`, `letta/plugins/defaults.py:1-8` |
| Extension registry definition | `DEFAULT_PLUGINS = {"experimental_check": {"target": "letta.plugins.defaults:is_experimental_enabled"}, "summarizer": {"target": "letta.services.summarizer.summarizer:Summarizer"}}` — comment says “Currently this supports one of each plugin type. This can be expanded in the future.” | `letta/plugins/plugins.py:15-25` |
| Plugin loader | `def get_plugin(plugin_type: str)` merges `DEFAULT_PLUGINS` + `settings.plugin_register_dict`, splits target on `":"`, does `importlib.import_module(module_path)` + `getattr`, returns function or instantiated class after optional `Protocol` check | `letta/plugins/plugins.py:28-42` |
| Runtime config surface | `plugin_register: Optional[str] = None` on `Settings`; `plugin_register_dict` splits `";"` then `"="` into `{"target": target}` — no validation, no protocol propagation for user-supplied plugins | `letta/settings.py:320`, `letta/settings.py:496-502` |
| Plugin lifecycle hooks | No install/uninstall/update hooks; only lazy singletons `_experimental_checker`/`_summarizer` with `get_experimental_checker()`, `get_summarizer()`, `reset_experimental_checker()`, `reset_summarizer()` plus TODO `handle coroutines` | `letta/plugins/plugins.py:45-72` |
| Plugin isolation mechanism | None — plugin is `importlib.import_module` into server process, runs with full trust. No subprocess, container, capability, or permission boundary; sandboxing exists only for *tool execution* (E2B/Modal/LOCAL), not plugins | `letta/plugins/plugins.py:34`, `letta/helpers/decorators.py:13-54` |
| Tool extensibility (data-driven) | `ToolType` enum (8 values: `CUSTOM`, `LETTA_CORE`, `LETTA_MEMORY_CORE`, `LETTA_MULTI_AGENT_CORE`, `LETTA_SLEEPTIME_CORE`, `LETTA_BUILTIN`, `LETTA_FILES_CORE`, `EXTERNAL_MCP`); `Tool` schema with `source_code` + `json_schema` generation via `ToolManager`/`derive_openai_json_schema` | `letta/schemas/enums.py:212-224`, `letta/schemas/tool.py:31-107`, `letta/services/tool_manager.py:206-287` |
| Tool execution extension | `ToolExecutor(ABC)` with `abstract async def execute`; `ToolExecutorFactory._executor_map: Dict[ToolType, Type[ToolExecutor]]` hard-coded, fallback `SandboxToolExecutor` for unknown types | `letta/services/tool_executor/tool_executor_base.py:16-46`, `letta/services/tool_executor/tool_execution_manager.py:35-43` |
| External tool protocol (MCP) | `ToolType.EXTERNAL_MCP`, `ToolCreate.from_mcp(cls, mcp_server_name, mcp_tool: MCPTool)`, `generate_tool_schema_for_mcp`, tag `mcp:SCHEMA_STATUS` | `letta/schemas/tool.py:145-169`, `letta/schemas/enums.py:224`, `letta/functions/schema_generator.py:694-902` |
| File-type registry as local registry pattern | `FileTypeRegistry.register(extension, mime_type, ...)` with 30+ default registrations, `get_chunking_strategy_by_extension` — demonstrates how Letta does registry without plugin system | `letta/services/file_processor/file_types.py:27-148` |
| Provider extension (hard-coded factory) | `LLMClient.create(provider_type: ProviderType)` is a `match` over 15+ `ProviderType` values, instantiating `AnthropicClient`, `OpenAIClient`, etc. directly — no registry, no entry-point | `letta/llm_api/llm_client.py:14-145` |
| Decorator extension point for plugins | `@experimental(feature_name, fallback_function, **kwargs)` captures `get_experimental_checker()` at decoration time (`experimental_checker = get_experimental_checker()` inside `decorator`) and dispatches to fallback vs primary | `letta/helpers/decorators.py:19-61` |
| Tests for extensibility | `tests/test_plugins.py:9-96` exercises `settings.plugin_register` override + `experimental` decorator for sync/async cases including Redis-backed inclusion checks; `tests/helpers/plugins_helper.py:5-19` provides `is_experimental_okay` test double | `tests/test_plugins.py:10-65`, `tests/helpers/plugins_helper.py:5-19` |
| Extension documentation | `letta/plugins/README.md:1-22` describes delimited `plugin_name.config_name=module:path` format and `DEFAULT_PLUGINS` shape; example key is `"default"` (stale) while code uses `"target"`; no versioning, stability, or security notes | `letta/plugins/README.md:7-21` |
| Absence of dynamic discovery | `pyproject.toml:1-230` has no `[project.entry-points]` / plugin groups; grep for `entry_points`, `importlib.metadata` shows only `typing-extensions`/`mypy-extensions` deps, no loader. Search for `plugin` across `letta/` returns only settings/plugins/decorators + noise in `.gitignore`/`otel` configs | `pyproject.toml:11-82`, `letta/settings.py:320` |

## Answers to Dimension Questions

**1. What can be extended via plugins?**

Via the formal plugin registry: almost nothing — only two slots. `summarizer` (must satisfy `SummarizerProtocol` at `letta/plugins/plugins.py:7-12`, validated at `letta/plugins/plugins.py:39-40`) and `experimental_check` (untyped function `letta/plugins/defaults.py:1`, slot with `protocol: None` at `letta/plugins/plugins.py:18`). Implementation strings are `module:path` pairs like `letta.plugins.defaults:is_experimental_enabled` and `letta.services.summarizer.summarizer:Summarizer` (`letta/plugins/plugins.py:19,23`).

Outside the plugin registry, Letta is extensible in a *data-driven/factory* sense — not a plugin sense:

- **Tools:** Third parties create arbitrary `ToolType.CUSTOM` tools by POSTing `ToolCreate(source_code: str, source_type, json_schema)` (`letta/schemas/tool.py:110-143`); schema is derived via `derive_openai_json_schema` (`letta/functions/functions.py:287`). MCP servers contribute `EXTERNAL_MCP` tools via `ToolCreate.from_mcp` (`letta/schemas/tool.py:145`). Built-in/letta-core tool sets are fixed enums (`letta/schemas/enums.py:212-224`) — adding a new `ToolType` requires core edits to `Tool.refresh_source_code_and_json_schema` (`letta/schemas/tool.py:74-107`) and `ToolExecutorFactory` (`letta/services/tool_executor/tool_execution_manager.py:35-43`).
- **File types:** `FileTypeRegistry.register` (`letta/services/file_processor/file_types.py:102-123`) allows adding MIME/extension mappings at runtime, but only within the process.
- **Providers/LLMs:** `ProviderType` lists 22 providers (`letta/schemas/enums.py:53-78`) and `LLMClient.create` (`letta/llm_api/llm_client.py:32-145`) hard-codes each — adding a provider means editing both.
- **Agents, memory, evals, prompts, policies, UI:** No plugin interface found. Agent variants are `AgentType` enum (`letta/schemas/enums.py:81-94`); memory is `Block`/`Archival` managers; prompts are `letta/prompts/` files; policies/approvals are `ToolRuleType`/`RequiresApprovalToolRule` — all core types, not extension points. UI is not pluggable from server core.

**2. Can plugins be loaded at runtime?**

Partially — **runtime string-target loading without discovery**.

- Mechanism: `LETTA_PLUGIN_REGISTER` env var / `plugin_register` setting (`letta/settings.py:320`) holds `name=module:attr;name2=module2:attr2`. It is parsed at call time in `get_plugin` (`letta/plugins/plugins.py:30-35`) via `importlib.import_module` + `getattr`. The helper `settings.plugin_register_dict` (`letta/settings.py:496-502`) does `split(";")` then `split("=")`.
- Lazy & cached: `get_experimental_checker()`/`get_summarizer()` memoize in `_experimental_checker`/`_summarizer` globals (`letta/plugins/plugins.py:45-62`); `reset_*` clears them (`letta/plugins/plugins.py:65-72`). The `@experimental` decorator captures the checker **at decoration time** (`letta/helpers/decorators.py:27`), so late plugin swaps require re-import or `reset_*` before next decoration.
- No filesystem scan, no hot-reload watcher, no entry-point discovery. `pyproject.toml:1-230` declares no `entry_points`. Search for `importlib.metadata`, `stevedore`, marketplace, or YAML plugin manifest across `letta/` shows no loader. `conf.yaml:1-413` and `config_file.py:1-232` have no plugin section.
- Limitation: only one plugin per type (`DEFAULT_PLUGINS` is flat dict, `letta/plugins/plugins.py:16 comment`). User target overrides default via `dict(DEFAULT_PLUGINS, **settings.plugin_register_dict)` (`letta/plugins/plugins.py:30`) — no merging of multiple plugins per slot.
- Bug/limit: user-supplied `plugin_register_dict` entries have only `{"target": ...}` (`letta/settings.py:501`) but `get_plugin` reads `plugin_register["protocol"]` as if it were `plugin_register[plugin_type]["protocol"]` (`letta/plugins/plugins.py:39`) — should be per-type. This silently skips Protocol validation for custom plugins and raises `KeyError` if `plugin_register["protocol"]` lookup is attempted differently.

> Can a third party add a new tool type without modifying core code? **No for a new `ToolType`; yes for a new custom tool instance.** Creating a `ToolCreate(source_code="def my_tool(...): ...")` (`letta/schemas/tool.py:110`) needs no core change — it is DB-persisted and executed via `SandboxToolExecutor`/`Modal`. But introducing a new category like `ToolType.MY_CUSTOM_KIND` requires editing `ToolType` (`letta/schemas/enums.py:212`), `Tool` validator (`letta/schemas/tool.py:82-107`), and `ToolExecutorFactory._executor_map` (`letta/services/tool_executor/tool_execution_manager.py:35`). No plugin hook exists to register a new `ToolExecutor` subclass externally.

**3. Are plugins isolated from each other?**

**No.** Plugins are `importlib.import_module` imports that run in the Letta server process with full Python permissions (`letta/plugins/plugins.py:34-35`). There is no subprocess, container, seccomp, capability, or per-plugin try/catch boundary.

- Both plugins share global interpreter state. `get_experimental_checker()` is evaluated at decorator definition time (`letta/helpers/decorators.py:27`); a faulty checker raises inside the wrapped function's `call_function` path (`letta/helpers/decorators.py:42-46`) and propagates to caller — no bulkheading.
- Failure modes propagate directly: `get_plugin` raises `TypeError("Unknown plugin type")` (`letta/plugins/plugins.py:42`) or `TypeError(f"{plugin} does not implement {Protocol}")` (`letta/plugins/plugins.py:40`) synchronously; there is no health-check or circuit breaker.
- Contrast with tool sandboxing, which *does* have isolation — `SandboxType` (`letta/schemas/enums.py:262-265`: `E2B`/`MODAL`/`LOCAL`) and `tool_settings.e2b_api_key`/`modal_token_id` (`letta/settings.py:24-28,57-71`) gate `SandboxToolExecutor` vs `E2B` vs `Modal` execution (`letta/services/tool_executor/tool_execution_manager.py:57`). This isolation does not apply to plugins; `ToolExecutor` ABC itself (`letta/services/tool_executor/tool_executor_base.py:16`) is not a plugin boundary for experimental/summarizer.
- No resource quotas, timeout, or per-plugin metrics. Metrics exist per-tool (`MetricRegistry().tool_execution_time_ms_histogram` at `letta/services/tool_executor/tool_execution_manager.py:113-115`) but not per-plugin.
- Redis flag check in test helper (`tests/helpers/plugins_helper.py:15-18` using `get_redis_client().check_inclusion_and_exclusion`) hints at external state sharing without isolation.

**4. Are extension points documented and stable?**

**Weakly documented, unstable/unstabilized.**

- The only dedicated doc is `letta/plugins/README.md:1-22` — 22 lines, describes a `plugin_name.config_name=class_or_function` delimited format and shows `DEFAULT_PLUGINS` with key `"default"` (`letta/plugins/README.md:18-21`), while actual code uses `"target"` (`letta/plugins/plugins.py:19`). No install guide, versioning, compatibility matrix, deprecation policy, or security model.
- Code comments signal instability: `letta/plugins/plugins.py:15` “Currently this supports one of each plugin type. This can be expanded in the future.” and `letta/plugins/plugins.py:49` “TODO handle coroutines”. `letta/functions/functions.py:12` “THIS FILE WILL BE DEPRECATED”. Summarizer settings TODO at `letta/settings.py:88-93` “should be deprecated or moved”.
- No stability contract. `SummarizerProtocol` (`letta/plugins/plugins.py:7-12`) is `@runtime_checkable` but not semver-guaranteed; experimental checker has no protocol at all. There is no `CHANGELOG` section for plugins, no `EXTENSIONS.md`, and external docs (`README.md:1-122`) point to `docs.letta.com` rather than local plugin author guide. Search for `plugin` in `fern/openapi.json` shows no plugin REST surface.
- Tests are the de facto spec: `tests/test_plugins.py:9-96` pins `plugin_register = "experimental_check=tests.helpers.plugins_helper:is_experimental_okay"` and exercises sync/async fallback logic — but test was not updated when README drifted, confirming doc/code skew.

## Architectural Decisions

- **String-target registry over entry-points** (`letta/plugins/plugins.py:28-42`, `letta/settings.py:496-502`): Choosing `LETTA_PLUGIN_REGISTER=name=module:attr;...` keeps the core dependency-free (no `importlib.metadata` scanning) and works with Docker env injection. Decision trades discoverability, validation, and multi-plugin-per-type for minimal plumbing. The `ToolExecutorFactory` (`letta/services/tool_executor/tool_execution_manager.py:35-43`) makes the same trade for tools — a dict of `ToolType -> ExecutorClass` rather than a pluggable registry.
- **Two-slot allowlist, not open plugin daemon** (`letta/plugins/plugins.py:15-25`): Limiting to `experimental_check` + `summarizer` avoids the lifecycle complexity of enable/disable/update/rollback and the security review of arbitrary code loading. This matches Letta’s threat model where *tool code* is sandboxed (Modal/E2B/LOCAL) but *server plugins* are trusted operator-supplied modules — documented implicitly by requiring a fully-qualified import path rather than an untrusted upload.
- **Protocol for summarizer, function for experimental check** (`letta/plugins/plugins.py:7-12` vs `letta/plugins/defaults.py:1`): A typed `Protocol` gives structural checking for the complex summarizer contract, while the experimental gate stays a loose function to allow any predicate signature (`**kwargs` forwarding at `letta/helpers/decorators.py:42,54`). This explains why `DEFAULT_PLUGINS["experimental_check"]["protocol"] is None` — flexibility over type safety.
- **Decorator-time checker capture** (`letta/helpers/decorators.py:27`): Evaluating `get_experimental_checker()` when `decorator(f)` runs (not when `f` is called) simplifies async/sync branching (`letta/helpers/decorators.py:28-46`) but freezes the plugin instance for the lifetime of the import. Decision favors per-process plugin stability over hot-swap.
- **Persistence + factory for tools vs. registry for plugins** (`letta/services/tool_manager.py:283-370`, `letta/schemas/tool.py:31-107`): Tools are DB rows (`ToolModel`) with `source_code`/`json_schema` and `ToolExecutor` routing; plugins are env-var strings. Decision separates user-content extension (tools, unlimited cardinality, per-org scoping via `organization_id` at `letta/services/tool_manager.py:318`) from operator extension (plugins, global, 1-per-type). This is why `create_or_update_tool_async` has atomic PostgreSQL ON CONFLICT handling (`letta/services/tool_manager.py:285-343`) but `get_plugin` has no concurrency handling beyond globals.
- **Stateless per-request tool sandbox vs. singleton plugin** (`letta/services/tool_executor/sandbox_tool_executor.py` vs `letta/plugins/plugins.py:45-62`): Tools get fresh sandbox config per execution (`ToolExecutionManager` at `letta/services/tool_executor/tool_execution_manager.py:68-92`); plugins are singletons reused across requests. Decision keeps hot path (tool loop) isolated while keeping cold path (experimental gating/summarization) fast via caching.

## Notable Patterns

- **Gateway import pattern:** `get_plugin` (`letta/plugins/plugins.py:33-34`) is the single gateway for both plugin types, splitting on `":"` — a micro pattern repeated in `letta/functions/functions.py:327-350` with `get_function_from_module`/`get_json_schema_from_module`. Consistent delimited-target convention across codebase.
- **Experimental gating via decorator** (`letta/helpers/decorators.py:19-62`): Higher-order `experimental(feature_name, fallback_function, **kwargs)` that merges decorator kwargs with call kwargs (`dict(_kwargs, **kwargs)` at `letta/helpers/decorators.py:42,54`) and branches to `fallback_function` when checker returns falsy. Used for Redis-backed group gating in tests (`tests/helpers/plugins_helper.py:15-18`).
- **Lazy singleton with explicit reset for tests** (`letta/plugins/plugins.py:45-72`): Module globals `_experimental_checker`, `_summarizer` + `reset_*` functions allow `tests/test_plugins.py:17,30,60,65` to mutate `settings.plugin_register` between cases without process restart — testability pattern, not an operational lifecycle.
- **Registry class template reused elsewhere:** `FileTypeRegistry` (`letta/services/file_processor/file_types.py:40-148`) shows the intended shape for a mature registry (dedicated class, `register` method, `get_supported_extensions`, `get_chunking_strategy_by_extension`) — but plugin registry did not adopt it, staying as a plain dict.
- **Hash-based redeploy detection for sandbox tools** (`letta/services/tool_manager.py:305-310,1007-1014`): `compute_tool_hash` stored in `metadata_["tool_hash"]` drives Modal redeploy; analogous version/hash mechanism is absent for plugins.

## Tradeoffs

| Tradeoff | Pro | Con | Evidence |
|----------|-----|-----|----------|
| Two-slot string env var vs open plugin marketplace | Zero-dep, Docker-friendly, no supply-chain risk | Expressiveness ceiling: one summarizer, one checker, no namespacing, no multi-tenant per-plugin enable | `letta/plugins/plugins.py:15-25`, `letta/settings.py:320,496` |
| `importlib.import_module` in-process vs sandboxed process/container | Fast, simple import, debuggable | No isolation: a broken plugin crashes server; no resource limits/timeouts | `letta/plugins/plugins.py:34-41`, no sandbox logic in plugin folder |
| Decoration-time checker capture vs call-time lookup | Fast inner loop (no import per call), handles sync/async split cleanly | Stale after `plugin_register` change unless caller resets globals or reimports module | `letta/helpers/decorators.py:27-54` |
| Protocol validation only for summarizer vs none for checker | Strong contract where it matters (async summarization) | Checker can be any function — typo in path or signature fails late at runtime, not at startup validation | `letta/plugins/plugins.py:39-41` vs `letta/plugins/defaults.py:1` |
| DB-persisted tools + factory vs plugin registry for tools | Auditable, per-org scoping (`organization_id`), atomic upsert (`ON CONFLICT`) — supports thousands of tools | New ToolType still requires core change; factory is closed to external executors without forking | `letta/services/tool_manager.py:285-343`, `letta/services/tool_executor/tool_execution_manager.py:35-43` |
| Hard-coded LLM provider `match` vs pluggable provider | Exhaustive `match` gives compile-time coverage (`ty` check at `pyproject.toml:200-216`) | Third-party model provider cannot be added without editing `ProviderType` + `LLMClient` + schemas | `letta/llm_api/llm_client.py:32-145`, `letta/schemas/enums.py:53-78` |

## Failure Modes / Edge Cases

- **Silent validation bypass:** `plugin_register_dict` (`letta/settings.py:500-501`) populates only `{"target": target}`. `get_plugin` (`letta/plugins/plugins.py:39`) then checks `plugin_register["protocol"]` (the *dict of all plugins*) not `plugin_register[plugin_type]["protocol"]`. For user-supplied plugins this key is missing or refers to wrong scope; branch is either skipped or raises `KeyError` depending on path — protocol enforcement is effectively dead for custom plugins. Also comparison `not isinstance(plugin, type(plugin_register["protocol"]))` (`letta/plugins/plugins.py:39`) is semantically wrong (checks class object vs protocol, not instance); runtime checkable protocol intended `isinstance(plugin(), SummarizerProtocol)` but code checks the class.
- **Fragile env parsing:** `plugin_register_dict` assumes `";"` and `"="` delimiters with no escaping, trimming, or error handling (`letta/settings.py:498-501`). Trailing `;` or extra `=` raises `ValueError: not enough values to unpack`; empty string yields `{}` silently. No early validation on startup.
- **Import side-effects and security:** `get_plugin` (`letta/plugins/plugins.py:34-35`) does `importlib.import_module` on operator-supplied string with no allowlist. Any importable module (including `os`, `subprocess`) can be instantiated if its class matches the loose check — attack surface for config injection.
- **Global memoization staleness:** `_experimental_checker` / `_summarizer` cached after first call (`letta/plugins/plugins.py:53-62`). Changing `LETTA_PLUGIN_REGISTER` at runtime requires manual `reset_experimental_checker()`/`reset_summarizer()`; health probes and `/v1/health` (`letta/settings.py:600-637` readiness) have no plugin health gate — a bad plugin that throws keeps returning error until process restart.
- **Doc/code drift:** `letta/plugins/README.md:18-21` documents `DEFAULT_PLUGINS[...]["default"]` while code uses `["target"]` (`letta/plugins/plugins.py:19`). New contributor following README will produce `KeyError` at `impl_path = plugin_register[plugin_type]["target"]` (`letta/plugins/plugins.py:32`).
- **No per-plugin failure bulkhead:** Any plugin exception propagates to the agent loop or request handler. In `@experimental` async wrapper (`letta/helpers/decorators.py:41-46`), a throwing checker aborts both primary and fallback paths; no fallback-to-fallback or metrics. Tool execution similarly lacks per-executor circuit breaker (`letta/services/tool_executor/tool_execution_manager.py:131-155` logs but does not isolate a poisoned executor).
- **Scalability ceiling:** Only one plugin per type means multi-tenant operators cannot enable different summarizers per organization — global singleton (`letta/plugins/plugins.py:45-46`) plus `summarizer_settings` (`letta/settings.py:645`) are process-wide.
- **Missing isolation tests:** No tests assert plugin sandboxing, mutual interference, or concurrent `get_plugin` race. Existing `tests/test_plugins.py:9-96` covers only positive checker behavior and Redis-backed gating.

## Future Considerations

- **Adopt class-based registry pattern:** Promote `DEFAULT_PLUGINS` dict to a `PluginRegistry` class modeled on `FileTypeRegistry` (`letta/services/file_processor/file_types.py:27-148`) with `register(name, protocol, factory, version, organization_id)` and `unregister`/`list` — keep `LETTA_PLUGIN_REGISTER` as backwards-compatible bootstrap. This fixes the `plugin_register["protocol"]` scoping bug and allows multiple plugins per type.
- **Fix `get_plugin` validation:** Change `letta/plugins/plugins.py:39` to `if plugin_register[plugin_type].get("protocol") and not isinstance(plugin() if issubclass ... )` with proper `issubclass`/`isinstance` + runtime-checkable guard, and validate allowed modules via allowlist (e.g., `letta.plugins.*`, `plugins.*`).
- **Move to call-time resolver or `importlib.metadata.entry_points`:** Replace decoration-time capture (`letta/helpers/decorators.py:27`) with per-call `get_experimental_checker()` lookup (cached via `lru_cache` with mtime on config) or entry-point group `letta.plugins` to enable pip-installable third-party plugins without env var edits — declare in `pyproject.toml` via `[project.entry-points."letta.plugins"]`.
- **Add isolation for untrusted plugins:** Reuse `SandboxType`/`tool_settings.sandbox_type` (`letta/settings.py:62-71`) infrastructure: run summarizer plugin in same sandbox as `SandboxToolExecutor` when flagged, with timeout (`tool_settings.tool_sandbox_timeout:36`) and char-limit truncation (`letta/services/tool_executor/tool_execution_manager.py:127-128`). Gate with `tool_settings.modal_sandbox_enabled` pattern.
- **Version & lifecycle contract:** Add `PluginManifest {name, version, peer_letta_version, protocol_version}` with semver check in `get_plugin`, persisted in `metadata_` similar to `tool_hash` (`letta/services/tool_manager.py:308`), plus `/v1/plugins` admin API and health probes mirroring `ReadinessSettings` (`letta/settings.py:600-637`). Add `__all__` stability docs alongside `letta/plugins/README.md`.
- **First-class tool-type plugin:** Open `ToolExecutorFactory._executor_map` (`letta/services/tool_executor/tool_execution_manager.py:35`) to external registration: `register_executor(tool_type: str, factory: Callable) -> None` guarded by `mcp_disable_stdio`-style flag, so a new MCP-like tool kind (e.g., `ToolType.EXTERNAL_SANDBOX`) does not require editing `letta/schemas/enums.py:212`.
- **Open provider plugin:** Switch `LLMClient.create` (`letta/llm_api/llm_client.py:32`) from closed `match` to `ProviderRegistry.get(provider_type) or fallback OpenAIClient` — allows custom OpenAI-compatible gateways without adding a `ProviderType` variant.
- **Documentation & testing:** Fix `letta/plugins/README.md` stale key (`"default"` -> `"target"`), add author guide with stable `SummarizerProtocol` version, deprecation timeline, and example pip package layout; add negative tests for malformed `LETTA_PLUGIN_REGISTER` (trailing semicolon, missing colon) and fault-injection test where summarizer throws.

## Questions / Gaps

- No evidence found for a plugin marketplace, filesystem plugin directory, or remote plugin fetching — searched `letta/plugins/`, `pyproject.toml` entry-points, `conf.yaml`, `config_file.py`, `server/`, and grep for `entry_points`/`marketplace`/`plugin.*manifest`; only incidental OTEL `.gitignore` hits returned.
- No evidence found for plugin-scoped configuration or secrets management beyond env-var strings; how a summarizer plugin would receive API keys or per-org settings beyond `settings` globals is unspecified. Search boundary: `letta/settings.py:1-648`, `letta/schemas/providers.py`, `letta/services/`.
- No evidence found for plugin observability: no per-plugin OTEL spans, no metric labels for plugin execution, no log channel per plugin; verified by absence of plugin-specific instrumentation in `letta/plugins/plugins.py`, `letta/helpers/decorators.py`, and telemetry settings (`letta/settings.py:532-598`).
- No evidence found for memory/eval/prompt/policy/UI extension points behind a plugin contract; file layout shows `letta/prompts/`, `letta/schemas/memory.py`, `letta/helpers/tool_rule_solver.py` but no `Plugin` or `Extension` interface for them. Confirmed by listing `letta/` directory (48 entries) and searching for `BaseTool` extension beyond schemas.
- Cross-language plugin story unclear — whether TypeScript tools (`ToolSourceType.typescript` at `letta/schemas/enums.py:236`, `source_type: "typescript"` at `letta/schemas/tool.py:114`) share the same `get_plugin` path or a separate JS runtime was not inspected in depth; `typescript_parser.py` handles source parsing but not plugin loading.
- Plugin hot-reload during udev: `Settings.uvicorn_reload` (`letta/settings.py:398`) reloads process, but in-proc `importlib.reload` for plugins without restart was not observed; whether a developer can iterate on a summarizer plugin without full restart is unknown.

---

Generated by `21.01-plugin-and-extension-points` against `letta`.
