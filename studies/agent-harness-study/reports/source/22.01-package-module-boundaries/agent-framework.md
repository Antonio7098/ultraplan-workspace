# Source Analysis: agent-framework

## Package and Module Boundaries

### Source Info

| Field | Value |
|-------|-------|
| Name | agent-framework |
| Path | `studies/agent-harness-study/sources/agent-framework` |
| Language / Stack | Python (3.10+) and .NET (multi-TFM), Microsoft Agent Framework |
| Analyzed | 2026-08-21 |

## Summary

The repository is a dual-language implementation of the Microsoft Agent Framework: a Python monorepo (`python/packages/*`) and a .NET solution (`dotnet/src/Microsoft.Agents.AI.*`). Both languages separate concerns into distinct packages/projects at the distribution level, but Python keeps everything inside a single `agent_framework` package namespace plus optional provider namespaces, while .NET splits each integration into its own NuGet assembly. The .NET side has a dedicated `Abstractions` assembly that referenced by every other project and contains only interfaces/abstract types — the cleanest layering in the codebase. The Python side uses private-by-convention modules (leading underscore names like `_agents.py`, `_clients.py`, `_workflows/`, `_harness/`) and a carefully curated `__all__` re-export in `__init__.py`, with lazy `__getattr__` namespaces for optional providers (e.g. `azure/`, `openai/`, `foundry/`) so a user can install `agent-framework-core` alone and pull providers only when needed. Dependency direction is mostly one-way: `core` has no deps on providers, providers depend on `core`, and cross-cutting features (e.g. evaluation, workflow inside evaluation) use late imports to avoid cycles. A few "most things import everything" Coupling hotspots exist inside Python's core package because all the call-graph heavyweight modules live in one wheel.

## Rating

**8 / 10** — Clear, durable model with explicit public/internal separation, lazy loading, and tests that enforce optional-dependency boundaries. Two soft spots keep it from a 9: (a) the Python core package ships a single wheel that re-exports everything from one `__init__.py`, so you cannot install just the runtime sub-tree without the harness/workflow/evaluation modules; (b) one runtime circular import is broken by a late import in `_tools.py` (`_tools.py:1947`), which is a smell even though it is wrapped in a function-local import.

## Evidence Collected

Every entry includes a file path with line numbers. Format: `path/to/file.py:NN`.

| Area | Evidence | File:Line |
|------|----------|-----------|
| Top-level package layout | Python monorepo with 31 packages under `python/packages/` | `python/packages/` (ls) |
| Top-level package layout | .NET solution with 36 project folders under `dotnet/src/` | `dotnet/src/` (ls) |
| Single Python core package | All core abstractions, workflows, harness, MCP, skills, evaluation live under one `agent_framework` directory | `python/packages/core/agent_framework/` (ls) |
| Two-language split | README and Python AGENTS.md describe "agent-framework-core" + providers | `python/AGENTS.md:50-72` |
| Two-language split | Both languages maintain their own AGENTS.md, skills, test layout | `python/AGENTS.md:1-10`, `dotnet/AGENTS.md:1-30` |
| Core package only depends on foundation libs | `dependencies` list is 4 entries: typing-extensions, pydantic, python-dotenv, opentelemetry-api | `python/packages/core/pyproject.toml:25-30` |
| Optional-deps declared as extras, not runtime deps | All providers listed under `[project.optional-dependencies].all` | `python/packages/core/pyproject.toml:32-59` |
| Core cannot list `tools` as runtime dep (avoids cycle) | `[dependency-groups].dev` carries it for type-check only with explanatory comment | `python/packages/core/pyproject.toml:61-67` |
| Provider packages depend on core only | e.g. `agent-framework-anthropic` declares `agent-framework-core>=1.9.0,<2` | `python/packages/anthropic/pyproject.toml:25-28` |
| Provider packages depend on core only | `agent-framework-orchestrations` depends solely on core | `python/packages/orchestrations/pyproject.toml:25-27` |
| Provider packages depend on core only | `agent-framework-tools` depends on core + psutil (no transitive cycle) | `python/packages/tools/pyproject.toml:24-31` |
| Lazy provider namespace | `azure/__init__.py` uses `_IMPORTS` table + `__getattr__` + `__dir__` | `python/packages/core/agent_framework/azure/__init__.py:11-41` |
| Lazy provider namespace | `openai/__init__.py` uses the same pattern | `python/packages/core/agent_framework/openai/__init__.py:17-47` |
| Lazy provider namespace | `foundry/__init__.py` lazy-loads from 4 packages via single dict | `python/packages/core/agent_framework/foundry/__init__.py:15-59` |
| Type stubs for lazy namespaces | `azure/__init__.pyi` declares re-exports for static type checkers | `python/packages/core/agent_framework/azure/__init__.pyi:1-36` |
| Internal modules use single underscore prefix | `_agents.py`, `_clients.py`, `_types.py`, `_tools.py`, `_middleware.py`, `_sessions.py`, `_mcp.py`, `_skills.py`, `_compaction.py`, `_evaluation.py`, `_harness/`, `_workflows/`, `security.py`, `observability.py`, `_serialization.py`, `_settings.py`, `_telemetry.py`, `_feature_stage.py`, `_docstrings.py`, `exceptions.py` | `python/packages/core/agent_framework/` (ls) |
| Public API defined in `__init__.py` | `from ._agents import Agent, BaseAgent, RawAgent, SupportsAgentRun` | `python/packages/core/agent_framework/__init__.py:20` |
| Public API defined in `__init__.py` | `__all__` lists 277 names | `python/packages/core/agent_framework/__init__.py:342-618` |
| Harness module is part of public API | `from ._harness._agent import DEFAULT_HARNESS_INSTRUCTIONS, create_harness_agent` and 8 other harness imports | `python/packages/core/agent_framework/__init__.py:85-147` |
| Workflows re-exported from public API | 20 imports from `._workflows.*` | `python/packages/core/agent_framework/__init__.py:258-331` |
| `py.typed` marker present | Empty `py.typed` file marks the package as typed | `python/packages/core/agent_framework/py.typed` |
| `_harness` is a sibling of `_workflows`, both hidden | `_harness/__init__.py` is empty (0 bytes) | `python/packages/core/agent_framework/_harness/__init__.py` |
| `_workflows` package namespace | `_workflows/__init__.py` is empty, container only | `python/packages/core/agent_framework/_workflows/__init__.py` |
| Workflows do not import from `_harness` | `_workflows` only references its own submodules and core (`_agents.py`, `_clients.py`, `_types.py`, `_middleware.py`, `_serialization.py`) | `python/packages/core/agent_framework/_workflows/` (grep) |
| Harness does not import from `_workflows` | `_harness/_agent.py` imports `.._agents`, `.._clients`, `.._compaction`, `.._feature_stage`, `.._sessions`, `.._skills` only | `python/packages/core/agent_framework/_harness/_agent.py:18-29` |
| Evaluation uses late imports to avoid cycle | `from ._workflows._agent_executor import AgentExecutorResponse` is inside `TYPE_CHECKING`/function block | `python/packages/core/agent_framework/_evaluation.py:62-63, 950-951, 1900` |
| One runtime late import in `_tools.py` | `from ._harness._tool_approval import ToolApprovalState` inside `_get_tool_approval_state` | `python/packages/core/agent_framework/_tools.py:1947` |
| `_harness` modules import upward to core only | e.g. `_tool_approval.py:12-16` imports `_feature_stage`, `_middleware`, `_serialization`, `_sessions`, `_types` | `python/packages/core/agent_framework/_harness/_tool_approval.py:12-16` |
| `_harness/_memory.py` imports core abstractions only | `.._clients`, `.._compaction`, `.._feature_stage`, `.._sessions`, `.._tools`, `.._types`, `..exceptions` | `python/packages/core/agent_framework/_harness/_memory.py:21-28` |
| Separation test for optional deps | `test_optional_dependencies.py` hides `opentelemetry.sdk` and `mcp` via monkey-patched `__import__` | `python/packages/core/tests/core/test_optional_dependencies.py:14-181` |
| Separation test for lazy loader | `test_azure_namespace.py` asserts `azure.CosmosHistoryProvider is CosmosHistoryProvider` | `python/packages/core/tests/core/test_azure_namespace.py:8-10` |
| Separation test for lazy loader | `test_foundry_namespace.py` (1 KB) sanity-checks `foundry` namespace | `python/packages/core/tests/core/test_foundry_namespace.py` |
| Separation test for lazy loader | `test_hyperlight_namespace.py` sanity-checks `hyperlight` namespace | `python/packages/core/tests/core/test_hyperlight_namespace.py` |
| Public/internal split on .NET side | `Microsoft.Agents.AI.Abstractions` is the leaf assembly: only depends on `Microsoft.Extensions.AI.Abstractions` | `dotnet/src/Microsoft.Agents.AI.Abstractions/Microsoft.Agents.AI.Abstractions.csproj:29-35` |
| Abstractions assembly contains only types, no behavior | `AIAgent.cs`, `AgentSession.cs`, `AIContextProvider.cs`, `ChatHistoryProvider.cs`, `DelegatingAIAgent.cs`, `InMemoryChatHistoryProvider.cs`, `MessageAIContextProvider.cs` | `dotnet/src/Microsoft.Agents.AI.Abstractions/` (ls) |
| `Microsoft.Agents.AI` (the implementation) only depends on Abstractions + MEAI | `ProjectReference` to Abstractions only; several `PackageReference` to MEAI family | `dotnet/src/Microsoft.Agents.AI/Microsoft.Agents.AI.csproj:21-33` |
| Internal access is gated with `InternalsVisibleTo` | `Microsoft.Agents.AI.csproj` lists test/peer assemblies | `dotnet/src/Microsoft.Agents.AI/Microsoft.Agents.AI.csproj:49-55` |
| Internal access is gated with `InternalsVisibleTo` | `Mcp` package exposes internals to `Microsoft.Agents.AI.Mcp.UnitTests` | `dotnet/src/Microsoft.Agents.AI.Mcp/Microsoft.Agents.AI.Mcp.csproj:35-37` |
| Workflows references both Abstractions and Microsoft.Agents.AI | `ProjectReference` to both | `dotnet/src/Microsoft.Agents.AI.Workflows/Microsoft.Agents.AI.Workflows.csproj:23-26` |
| Workflows exposes internals to DurableTask and test assemblies | `InternalsVisibleTo` | `dotnet/src/Microsoft.Agents.AI.Workflows/Microsoft.Agents.AI.Workflows.csproj:28-32` |
| Provider packages depend on Abstractions only (no runtime) | `A2A` references only `Microsoft.Agents.AI.Abstractions` | `dotnet/src/Microsoft.Agents.AI.A2A/Microsoft.Agents.AI.A2A.csproj:27` |
| Provider packages depend on Abstractions only | `CopilotStudio` references only Abstractions | `dotnet/src/Microsoft.Agents.AI.CopilotStudio/Microsoft.Agents.AI.CopilotStudio.csproj:15` |
| Provider packages depend on Abstractions only | `GitHub.Copilot` references only Abstractions | `dotnet/src/Microsoft.Agents.AI.GitHub.Copilot/Microsoft.Agents.AI.GitHub.Copilot.csproj:18` |
| Provider packages depend on Abstractions only | `Mcp` references only `Microsoft.Agents.AI` (implementation) | `dotnet/src/Microsoft.Agents.AI.Mcp/Microsoft.Agents.AI.Mcp.csproj:40` |
| Provider packages depend on Abstractions only | `Hyperlight` references only Abstractions | `dotnet/src/Microsoft.Agents.AI.Hyperlight/Microsoft.Agents.AI.Hyperlight.csproj:19` |
| Provider packages depend on Abstractions only | `Mem0` references only Abstractions | `dotnet/src/Microsoft.Agents.AI.Mem0/Microsoft.Agents.AI.Mem0.csproj:24` |
| Provider packages depend on Abstractions only | `CosmosNoSql` references Abstractions + Workflows | `dotnet/src/Microsoft.Agents.AI.CosmosNoSql/Microsoft.Agents.AI.CosmosNoSql.csproj:27-28` |
| Provider packages depend on Abstractions only | `LocalCodeAct` references only Abstractions | `dotnet/src/Microsoft.Agents.AI.LocalCodeAct/Microsoft.Agents.AI.LocalCodeAct.csproj:19` |
| Provider packages depend on Abstractions only | `Valkey` references only Abstractions | `dotnet/src/Microsoft.Agents.AI.Valkey/Microsoft.Agents.AI.Valkey.csproj:28` |
| Provider packages depend on Abstractions only | `Tools.Shell` references only Abstractions | `dotnet/src/Microsoft.Agents.AI.Tools.Shell/Microsoft.Agents.AI.Tools.Shell.csproj:37` |
| Harness depends on Microsoft.Agents.AI (+ Tools.Shell on net8.0+) | `Microsoft.Agents.AI.Harness.csproj` | `dotnet/src/Microsoft.Agents.AI.Harness/Microsoft.Agents.AI.Harness.csproj:15-19` |
| Hosting gives DI the building blocks | Depends on Abstractions + Workflows + Microsoft.Agents.AI | `dotnet/src/Microsoft.Agents.AI.Hosting/Microsoft.Agents.AI.Hosting.csproj:16-18` |
| Declarative adds on top of core + workflows | `Microsoft.Agents.AI.Declarative.csproj` references Abstractions + Microsoft.Agents.AI | `dotnet/src/Microsoft.Agents.AI.Declarative/Microsoft.Agents.AI.Declarative.csproj:33-34` |
| Workflows.Declarative adds on top of Workflows | `Microsoft.Agents.AI.Workflows.Declarative.csproj` references only Workflows | `dotnet/src/Microsoft.Agents.AI.Workflows.Declarative/Microsoft.Agents.AI.Workflows.Declarative.csproj:41` |
| Workflows.Declarative.Foundry composes two leaf projects | References both Foundry and Workflows.Declarative | `dotnet/src/Microsoft.Agents.AI.Workflows.Declarative.Foundry/Microsoft.Agents.AI.Workflows.Declarative.Foundry.csproj:33-34` |
| Workflows.Declarative.Mcp composes Workflows.Declarative | `Microsoft.Agents.AI.Workflows.Declarative.Mcp.csproj:34` | `dotnet/src/Microsoft.Agents.AI.Workflows.Declarative.Mcp/Microsoft.Agents.AI.Workflows.Declarative.Mcp.csproj:34` |
| Python docs describe the layering | `Python AGENTS.md` "Package Relationships" notes lazy loading via `__getattr__` in provider folders | `python/AGENTS.md:68-72` |
| Python docs describe the layering | `core/AGENTS.md` "Module Structure" lists each module's role | `python/packages/core/AGENTS.md:5-21` |
| Workflows internal-API documented | `output_from` / `intermediate_output_from` allow-list for emissions | `python/AGENTS.md:132-143` (referenced in core/AGENTS.md) |
| Feature-stage gating for experimental APIs | `ExperimentalFeature` / `experimental` decorator used in `_feature_stage.py` | `python/packages/core/agent_framework/_feature_stage.py` |
| Package lifecycle is tracked per package | `PACKAGE_STATUS.md` enumerates lifecycle state for every package | `python/PACKAGE_STATUS.md:13-46` |

## Answers to Dimension Questions

1. **Are modules cleanly separated?**
   Yes, at the package/namespace level. .NET: 36 distinct projects under `dotnet/src/` with a single `Abstractions` leaf (`Microsoft.Agents.AI.Abstractions.csproj:29-35`). Python: 31 packages plus 23 sub-namespaces under `agent_framework/`. Internal modules are prefixed with `_` (`_agents.py`, `_clients.py`, `_mcp.py`, `_harness/_loop.py`, `_workflows/_edge.py`). The `__init__.py` is the only public surface and it does not export the underscore-prefixed modules directly — they are reached only through the curated 277-name `__all__` (`python/packages/core/agent_framework/__init__.py:342-618`). The clear weak point is runtime: `agent-framework-core` ships all of these in one wheel, so a "smallest possible install" is a wheel boundary, not a module boundary.

2. **Do dependencies flow in one direction?**
   Yes, with two minor exceptions handled with discipline.
   - Core has no dependency on providers (`python/packages/core/pyproject.toml:25-30`).
   - Providers depend on core (`python/packages/anthropic/pyproject.toml:25-28`, `python/packages/orchestrations/pyproject.toml:25-27`, `python/packages/tools/pyproject.toml:24-31`).
   - `_harness` depends on root top-level modules (`_clients.py`, `_compaction.py`, `_sessions.py`, etc.) and never on `_workflows` (`_harness/_agent.py:18-29`).
   - `_workflows` does not import `_harness` (no grep hit in either direction).
   - The only cross-cutting late import is `_tools.py:1947` importing `_harness._tool_approval` inside a function; this is a deliberate cycle break, flagged here as a smell.

3. **Can modules be used independently?**
   - .NET: yes. Each provider is its own NuGet package with a focused `<ProjectReference>` set; `A2A` (`Microsoft.Agents.AI.A2A.csproj:27`), `Mem0` (`Microsoft.Agents.AI.Mem0.csproj:24`), `Hyperlight` (`Microsoft.Agents.AI.Hyperlight.csproj:19`), `Valkey` (`Microsoft.Agents.AI.Valkey.csproj:28`) and others only depend on `Abstractions`, so a consumer can ship a single provider without dragging in the core implementation.
   - Python: providers can be installed independently (`python/packages/anthropic/pyproject.toml:25-28`) and the core `azure`, `openai`, `foundry` namespaces lazy-load providers on first attribute access (`python/packages/core/agent_framework/azure/__init__.py:27-37`). However, the **contents** of `agent-framework-core` are inseparable at runtime — `_harness`, `_workflows`, `_evaluation`, `_mcp`, `_skills`, `_compaction` are all in the same wheel and eagerly imported via `__init__.py` lines 20-340.

4. **Are public APIs distinguished from internal ones?**
   - Python: enforced by convention. The leading underscore on every internal module (`_agents.py`, `_clients.py`, `_tools.py`, `_harness/`, `_workflows/`, etc.) is consistent across the core. The public surface is the `__all__` list in `__init__.py` (`python/packages/core/agent_framework/__init__.py:342-618`). The `azure/`, `openai/`, `foundry/` namespaces also use `__dir__` to advertise only the names they can actually resolve (`python/packages/core/agent_framework/openai/__init__.py:46`).
   - .NET: enforced by language and tooling. `Microsoft.Agents.AI.Abstractions` exists precisely to host interface-only types (`AIAgent.cs`, `AgentSession.cs`, `AIContextProvider.cs`, `ChatHistoryProvider.cs`, `DelegatingAIAgent.cs`). Implementation files live in `Microsoft.Agents.AI/ChatClient/`, `Microsoft.Agents.AI/Harness/`, etc. with `InternalsVisibleTo` exposing internals only to specific test assemblies (`Microsoft.Agents.AI/Microsoft.Agents.AI.csproj:49-55`).

## Architectural Decisions

- **Single Python core package, curated public surface.** All agent abstractions, workflows, harness, MCP, skills, evaluation, and telemetry live in `agent-framework-core`. The 23 sub-namespaces inside `agent_framework/` are reached via the `__init__.py` re-export tuple, so users get a flat `from agent_framework import Agent, Workflow, SkillsProvider, FileMemoryProvider` API. Trade-off: install size and import cost; gain: ergonomic single import surface.
- **Dedicated Abstractions assembly on .NET.** `Microsoft.Agents.AI.Abstractions` contains interfaces and abstract classes only (`AIAgent`, `AgentSession`, `AIContextProvider`, `ChatHistoryProvider`, `DelegatingAIAgent`, `InMemoryChatHistoryProvider`, `MessageAIContextProvider`). This is the single graph root that virtually every other package depends on, enabling pluggable providers without depending on the implementation.
- **Lazy provider namespaces.** `azure/__init__.py`, `openai/__init__.py`, `foundry/__init__.py` use `__getattr__` + an `_IMPORTS` table (`python/packages/core/agent_framework/azure/__init__.py:11-24`) so the optional provider packages do not need to be installed to import `agent_framework`. Errors are surfaced with actionable messages naming the missing package and `pip install` command.
- **Type stubs alongside lazy loaders.** `azure/__init__.pyi` (`python/packages/core/agent_framework/azure/__init__.pyi:1-36`) re-declares the names for static type checkers, separating runtime behavior from type information. Other namespaces follow the same pattern.
- **Optional-dependency discipline.** `agent-framework-core` keeps `mcp` and `opentelemetry-sdk` out of runtime deps; tests verify the core remains importable when they are missing (`python/packages/core/tests/core/test_optional_dependencies.py:14-181`).
- **Cycle avoidance via late imports.** `_evaluation.py` and `_tools.py` use function-local imports to break the would-be cycles (`_evaluation.py:62-63, 950-951, 1900`; `_tools.py:1947`). The `_harness` ↔ `_tools` cycle is the only one that needed a runtime break.
- **Feature-stage gating.** `_feature_stage.py` defines `ExperimentalFeature` and `ReleaseCandidateFeature` (`python/packages/core/agent_framework/__init__.py:84`) so experimental/rc APIs can be marked without removing them from the public surface. `PACKAGE_STATUS.md` enumerates feature-level staged APIs (`python/PACKAGE_STATUS.md:58-83`).
- **`InternalsVisibleTo` for friend assemblies.** Both `Microsoft.Agents.AI` and `Microsoft.Agents.AI.Mcp` opt into peer test visibility (`Microsoft.Agents.AI/Microsoft.Agents.AI.csproj:49-55`) while restricting it to specific named assemblies.

## Notable Patterns

- **`_IMPORTS` dictionary + `__getattr__` + `__dir__` triplet.** Each lazy namespace (`azure/`, `openai/`, `foundry/`, `chatkit/`, `a2a/`, etc.) repeats the same shape: `dict[str, tuple[import_path, package_name]]` + module-level `__getattr__` that raises an install hint on `ModuleNotFoundError` + `__dir__` returning the dict keys. This makes the implementation discoverable but also predictable.
- **Empty `__init__.py` for hidden subpackages.** Both `_harness/__init__.py` (0 bytes) and `_workflows/__init__.py` (1 line) are placeholders, so the subpackages are reachable through explicit dotted paths (`_harness._agent`) and not through wildcard `from _harness import *`.
- **Eager top-level re-export in `__init__.py`.** Lines 20-340 of `python/packages/core/agent_framework/__init__.py` make every public name available at the top of the package. There is no deprecation gate; `__all__` is the curated allow-list.
- **Per-package `py.typed` marker.** `python/packages/core/agent_framework/py.typed` is present (verified empty file), signaling to mypy/pyright that the package ships type annotations.
- **Per-module `AGENTS.md`.** Each major package/core has its own `AGENTS.md` (`python/packages/core/AGENTS.md`, `python/packages/anthropic/AGENTS.md`, `dotnet/AGENTS.md`) that documents the module's role, public classes, and conventions. This is a documentation-as-boundary pattern.
- **`shared_tasks.toml` for cross-package tooling.** `python/packages/core/pyproject.toml:122-124` references `../../shared_tasks.toml` so the same `poe` tasks (mypy, test, etc.) are reused across every Python package.

## Tradeoffs

- **Single Python wheel for core.** Installing `agent-framework-core` imports workflows, harness, MCP, skills, evaluation, telemetry, and observability code at startup (`python/packages/core/agent_framework/__init__.py:20-340`). A user who only wants `Agent` pays the import cost (and accepts runtime types like `Workflow`, `FileMemoryProvider`, `MCPTool`) until the package is split. The codebase hides this behind a curated `__all__`, but the import-time cost is real.
- **`_harness` shipped in core.** The harness module is `@experimental` and lives in core (`python/packages/core/agent_framework/_harness/_loop.py` is decorated via `_feature_stage.py:ExperimentalFeature` per `python/AGENTS.md` references). This is convenient but couples core to a still-evolving subsystem; the cyclic dependency with `_tools.py` is the price.
- **No `__all__` in internal modules.** Internal modules (`_agents.py`, `_clients.py`, `_mcp.py`, etc.) do not declare `__all__`, so `from _agents import *` is technically possible. The leading underscore convention is the only enforcement.
- **One runtime late import.** `_tools.py:1947` is a function-local `from ._harness._tool_approval import ToolApprovalState`. The comment in `pyproject.toml:62-66` and the late-import pattern in `_evaluation.py` together suggest the author expected cycles and accepted them as the cost of single-package core.
- **Provider packages reference `Microsoft.Agents.AI` (the implementation) rather than only `Abstractions` for some integrations.** `CosmosNoSql` (`Microsoft.Agents.AI.CosmosNoSql.csproj:27-28`) and `Hosting` (`Microsoft.Agents.AI.Hosting/Microsoft.Agents.AI.Hosting.csproj:16-18`) take dependencies on both. This is the right call for richer integrations but does pull the implementation package into more callers than strictly necessary.
- **Lots of assemblies on .NET.** 36 projects in `dotnet/src/` means every consumer's `dotnet restore` resolves many packages. The `agent-framework-release.slnf` filter (`dotnet/agent-framework-release.slnf`) is the mechanism for keeping release builds reasonable.

## Failure Modes / Edge Cases

- **Late-import fragility.** If someone moves `_get_tool_approval_state` out of `_tools.py` to a top-level module, the cycle break breaks. The `_evaluation.py` late-imports (`python/packages/core/agent_framework/_evaluation.py:62-63`) are similarly fragile.
- **Lazy-loader error messages.** When a name is missing, the error message tells the user which package to install (`python/packages/core/agent_framework/azure/__init__.py:33-36`). If a user calls a name that does not exist in `_IMPORTS`, the error is "Module `azure` has no attribute `name`" — useful but easy to confuse with the install message.
- **`py.typed` is empty.** Verified: `python/packages/core/agent_framework/py.typed` is 0 bytes. The `py.typed` marker is what enables downstream type checking; the lack of explicit pointer just means downstream consumers should rely on the wheel's inline annotations on public names.
- **Test-only `agent-framework-tools` dependency.** `python/packages/core/pyproject.toml:61-67` lists `agent-framework-tools` under `[dependency-groups].dev` only because core cannot list it as a runtime dep. This is a documented design choice but it does mean integration tests under `tests/core/test_harness_*.py` need both packages present.
- **`InternalsVisibleTo` amplification.** Each project that uses `InternalsVisibleTo` (`Microsoft.Agents.AI.csproj:49-55`, `Microsoft.Agents.AI.Mcp.csproj:35-37`, `Microsoft.Agents.AI.Workflows.csproj:28-32`) widens the surface that friend assemblies can see. If a friend assembly is later repurposed, the exposed internals follow it.
- **Multiple `Microsoft.Agents.AI.*` namespace roots on .NET.** `Microsoft.Agents.AI.Abstractions.csproj:4` (`<RootNamespace>Microsoft.Agents.AI</RootNamespace>`) shares the root namespace with `Microsoft.Agents.AI` itself. This is intentional (the abstraction types live in `Microsoft.Agents.AI` namespace) but can confuse static analyzers that key off namespace/file path.

## Future Considerations

- **Split `agent-framework-core` Python wheel.** Consider separating `agent-framework-core-runtime` (agents, clients, sessions, tools, middleware) from `agent-framework-core-orchestration` (workflows, harness, evaluation). This would let users opt out of the workflow and harness modules while preserving the same import path via a meta-package.
- **Formalize module-level `__all__`.** Add `__all__` to every internal module (`_agents.py`, `_clients.py`, `_tools.py`, `_mcp.py`, `_sessions.py`, etc.) so that wildcard imports cannot reach internals even by accident.
- **Move `_tools.py` late import to a top-level accelerator.** The `_get_tool_approval_state` runtime import of `_harness._tool_approval` (`_tools.py:1947`) could be replaced by a `Protocol` in `_types.py` so the dependency arrow reverses.
- **Stamp `py.typed` everywhere.** Most Python wheels have `py.typed`; verify each provider package (`agent-framework-anthropic`, `agent-framework-openai`, etc.) ships its own. The `py.typed` file in `agent_framework/` is good but the provider wheel may not have one.
- **Slim the .NET implementation package.** Several providers that currently reference `Microsoft.Agents.AI` plus `Microsoft.Agents.AI.Abstractions` (e.g. `CosmosNoSql`, `Hosting`) could in principle depend only on the abstractions for read paths and on the implementation for the few types they need. The current state is acceptable but worth re-evaluating.
- **Single source of truth for module dependency graph.** The Python side documents its layering in `python/AGENTS.md:50-72` and `python/packages/core/AGENTS.md:5-21`. A machine-checked dependency review (e.g. `import-linter` or `pydeps`) would catch the one late-import cycle and any future regressions.

## Questions / Gaps

- Is there an automated check that `_harness` and `_workflows` never import each other? The current evidence is grep-based; an enforced CI check would prevent regressions.
- Does the .NET solution have a `dotnet format`/analyzer rule that fails the build when a new `<ProjectReference>` is added that violates the layering (e.g. a provider depending on `Microsoft.Agents.AI.Hosting`)? No evidence found.
- Why does `_harness/_tool_approval.py` exist in the core package when `agent-framework-tools` is the "tools" package and core explicitly *cannot* depend on it at runtime? The split between "tools-in-core" (`_harness/_tool_approval.py`) and "tools-in-tools-package" (`agent_framework_tools/shell/`) is documented only in `python/packages/core/AGENTS.md`. The boundary is policy, not enforced.
- Is the `__init__.py` re-export of every `from ._harness._foo import *` deliberately eager, or is it a place where convertible lazy attributes would shave startup time? No evidence of a benchmark; assumed deliberate.
- Is the `dotnet/Shared/` directory (which contains `Workflows/`, `Throw/`, `DiagnosticIds/`, `Redaction/`, `StructuredOutput/`) treated as a separate package or is it conditional-included? `dotnet/src/Shared/Workflows/Execution/WorkflowFactory.cs` is referenced under `dotnet/eng` style assets but no `<ProjectReference>` to it was found in the surveyed `src/` projects. **No clear evidence found** of how `Shared` is consumed; further investigation needed.

---

Generated by `22.01-package-module-boundaries` against `agent-framework`.
