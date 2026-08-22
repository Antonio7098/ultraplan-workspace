# Source Analysis: crewai

## Dimension 24.01: Public API Surface

### Source Info

| Field | Value |
|-------|-------|
| Name | crewai |
| Path | `studies/agent-harness-study/sources/crewai` |
| Language / Stack | Python 3.10–3.13, Pydantic v2, Click (CLI), uv workspace (multi-package monorepo) |
| Analyzed | 2026-08-22 |

All citations below are relative to the source root `studies/agent-harness-study/sources/crewai/`.

## Summary

CrewAI exposes its public surface through five published PyPI packages plus one private dev package, organized as a uv workspace (`pyproject.toml:227-235`): `crewai` (the user-facing framework), `crewai-core` (shared utilities), `crewai-tools` (tool library, optional extra), `crewai-files` (multimodal file handling), and `crewai-cli`. The primary Python API is a deliberately small root namespace — 17 curated symbols in `__all__` (`lib/crewai/src/crewai/__init__.py:187-205`) covering the core abstractions (`Agent`, `Task`, `Crew`, `Flow`, `LLM`, `BaseLLM`, `Process`, `Memory`, `Knowledge`), with deeper extension surfaces split into purposeful subpackages: `crewai.project` (decorator-based crew definition and JSON crew loading), `crewai.tools` (`BaseTool`, `@tool`), `crewai.events` (typed event bus), `crewai.hooks` (LLM/tool call hooks), `crewai.a2a`, and `crewai.state`. A single CLI entry point `crewai` is registered by both the `crewai` and `crewai-cli` packages (`lib/crewai/pyproject.toml:146-147`, `lib/cli/pyproject.toml:34-35`) with ~27 commands/groups (`lib/cli/src/crewai_cli/cli.py:115-1105`). Documentation is extensive, versioned per release, and example-driven; the Python API itself has no generated reference docs — the `api-reference` docs section covers only the enterprise HTTP API (`docs/edge/en/api-reference/kickoff.mdx:6`).

The surface is generally well-layered (stable root exports, documented extension ABCs, explicit `crewai.experimental` namespace with runtime gating), but boundary discipline is inconsistent: internal machinery (`CacheHandler`, `ToolsHandler`, `parse`) is re-exported from `crewai.agents` (`lib/crewai/src/crewai/agents/__init__.py:3-19`), import-time hygiene relies on a global monkey-patch of `warnings.warn` (`lib/crewai/src/crewai/__init__.py:28-49`), and a Pydantic forward-reference workaround injects names into eight modules' `__dict__` via `sys.modules` (`lib/crewai/src/crewai/__init__.py:128-145`). Public-surface tests are thin — two symbols (`lib/crewai/tests/test_imports.py:4-13`).

## Rating

**7 / 10** — Clear model with explicit interfaces, curated root exports, documented extension points (custom LLMs, custom tools, event listeners, hooks), and operational safeguards (experimental gating, deprecation warnings with removal targets). Falls short of 8–9 because the Python API lacks generated reference documentation, the stable/internal boundary is enforced only by naming convention, surface tests are minimal, and import-time global monkey-patching creates fragile coupling that a new integration could trip over.

## Evidence Collected

| Area | Evidence | File:Line |
|------|----------|-----------|
| Root public API | Curated `__all__` with 17 symbols (Agent, Crew, Flow, Task, LLM, BaseLLM, Process, Memory, Knowledge, RuntimeState, …) | `lib/crewai/src/crewai/__init__.py:187-205` |
| Version identity | `__version__ = "1.14.8a2"` (pre-release channel) | `lib/crewai/src/crewai/__init__.py:51` |
| Lazy heavy import | `Memory` resolved via module `__getattr__` to avoid importing lancedb eagerly | `lib/crewai/src/crewai/__init__.py:53-66` |
| Import-time monkey patch | Global `warnings.warn` filtered to suppress Pydantic deprecations at package import | `lib/crewai/src/crewai/__init__.py:28-49` |
| Namespace injection hack | `_resolve_namespace` written into `sys.modules[...].__dict__` of 8 modules to resolve Pydantic forward refs | `lib/crewai/src/crewai/__init__.py:128-145` |
| Workspace layout | uv workspace members: crewai, crewai-tools, devtools, crewai-files, cli, crewai-core | `pyproject.toml:227-235` |
| Private package | `crewai-devtools` marked `private = true` + `Private :: Do Not Upload` | `lib/devtools/pyproject.toml:10-11` |
| CLI entry point | `[project.scripts] crewai = "crewai_cli.cli:crewai"` registered by both `crewai` and `crewai-cli` | `lib/crewai/pyproject.toml:146-147`, `lib/cli/pyproject.toml:34-35` |
| CLI command surface | Click group with commands `create`, `version`, `train`, `replay`, `log-tasks-outputs`, `reset-memories`, groups for deploy/auth/org/etc. (~27 registrations) | `lib/cli/src/crewai_cli/cli.py:115-1105` |
| CLI lazy loading | Command classes shimmed behind `__new__` to defer heavy imports until invocation | `lib/cli/src/crewai_cli/cli.py:63-105` |
| Project/decorator API | `crewai.project` exports `@agent`, `@task`, `@crew`, `@before_kickoff`, `@after_kickoff`, `CrewBase`, `load_crew`, JSON/JSONC loaders | `lib/crewai/src/crewai/project/__init__.py:1-40` |
| Tool extension API | `crewai.tools` exports exactly `BaseTool`, `EnvVar`, `tool` | `lib/crewai/src/crewai/tools/__init__.py:1-7` |
| Tool ABC | `class BaseTool(BaseModel, ABC)` with `__init_subclass__` type registry for checkpoint deserialization | `lib/crewai/src/crewai/tools/base_tool.py:103`, registry at `lib/crewai/src/crewai/tools/base_tool.py:47-48` |
| `@tool` decorator | Overloaded `tool()` accepting name/args_schema/caching | `lib/crewai/src/crewai/tools/base_tool.py:653-676` |
| Custom LLM ABC | `BaseLLM(BaseModel, ABC)` with abstract `call()`; docstring invites user subclasses | `lib/crewai/src/crewai/llms/base_llm.py:129-146`, abstract `call` at `lib/crewai/src/crewai/llms/base_llm.py:283-284` |
| Events extension API | `crewai.events` module docstring: monitor/extend via `BaseEventListener`, `crewai_event_bus`, typed events; lazy event-type loading | `lib/crewai/src/crewai/events/__init__.py:1-11` |
| Event bus API | `crewai_event_bus.on(...)` and `emit()` public methods on singleton bus | `lib/crewai/src/crewai/events/event_bus.py:244`, `lib/crewai/src/crewai/events/event_bus.py:570` |
| Hooks API | `crewai.hooks` re-exports `register_before/after_llm_call_hook`, `@before_llm_call`, `@after_llm_call`, tool-hook equivalents | `lib/crewai/src/crewai/hooks/__init__.py:1-30` |
| Experimental namespace | Module docstring: "APIs here may change without major-version bumps" | `lib/crewai/src/crewai/experimental/__init__.py:1` |
| Experimental runtime gate | `CREWAI_EXPERIMENTAL=1` env flag required for Skills Repository; raises `ExperimentalFeatureDisabledError` | `lib/crewai/src/crewai/experimental/skills/_flag.py:9-24` |
| Deprecation policy | `max_retries` warns "removed in CrewAI v1.0.0" via `DeprecationWarning` | `lib/crewai/src/crewai/task.py:561-566` |
| Deprecation policy (fields) | Pydantic `deprecated=True` on Agent fields; `allow_code_execution` warns "removed in v2.0" | `lib/crewai/src/crewai/agent/core.py:227-285`, `lib/crewai/src/crewai/agent/core.py:368-371` |
| Stable-path aliasing | `crewai.plus_api` re-export kept "as a stable import path"; new code directed to `crewai_core.plus_api` | `lib/crewai/src/crewai/plus_api.py:1-11` |
| A2A public config | `crewai.a2a` exports only `A2AConfig`, `A2AClientConfig`, `A2AServerConfig` | `lib/crewai/src/crewai/a2a/__init__.py` |
| State provider API | `crewai.state` exports `CheckpointConfig`, `JsonProvider`, `SqliteProvider` | `lib/crewai/src/crewai/state/__init__.py:1-10` |
| Accidental/internal exports | `crewai.agents` `__all__` includes internals `CacheHandler`, `ToolsHandler`, `parse`, `AgentAction`, `AgentFinish` | `lib/crewai/src/crewai/agents/__init__.py:3-19` |
| Tools package surface | `crewai_tools/__init__.py`: 333 lines, ~209 import statements eagerly exporting 100+ tool classes (SDK usage deferred inside methods) | `lib/crewai-tools/src/crewai_tools/__init__.py:1-333` (e.g. deferred selenium import at `lib/crewai-tools/src/crewai_tools/tools/selenium_scraping_tool/selenium_scraping_tool.py:75`) |
| Surface test coverage | Only `TaskOutput` and `CrewOutput` import tests for the root API | `lib/crewai/tests/test_imports.py:4-13` |
| Docs: quickstart | Runnable Flow+crew tutorial using `crewai create flow`, `load_crew`, `crew.kickoff` | `docs/edge/en/quickstart.mdx:31-100` |
| Docs: canonical imports | `from crewai import Agent, Crew, Task, Process` as the documented entry | `docs/edge/en/concepts/crews.mdx:125` |
| Docs: event listener guide | Full runnable `BaseEventListener` example with `crewai_event_bus.on(...)` | `docs/edge/en/concepts/event-listener.mdx:31-60` |
| Docs: custom LLM guide | Required methods (`__init__`, abstract `call()`), error-handling and function-calling patterns for `BaseLLM` subclasses | `docs/edge/en/learn/custom-llm.mdx:109-148` |
| Docs: custom tools guide | Subclassing `BaseTool` and `@tool` decorator examples | `docs/edge/en/concepts/tools.mdx:168-315` |
| Docs: CLI reference | All CLI commands documented with options | `docs/edge/en/concepts/cli.mdx:35-370` |
| Docs: API reference scope | `api-reference` pages bind to enterprise OpenAPI spec (`/enterprise-api.en.yaml POST /kickoff`), not the Python API | `docs/edge/en/api-reference/kickoff.mdx:6` |
| Docs versioning model | Edge channel + frozen `docs/vX.Y.Z/` snapshots per release; CI guards against snapshot edits | `AGENTS.md:1-40` |
| Enterprise HTTP surface | OpenAPI specs per language under docs (e.g. `enterprise-api.en.yaml`) | `docs/edge/en/enterprise-api.en.yaml` (referenced by `docs/edge/en/api-reference/*.mdx`) |

## Answers to Dimension Questions

**1. What is the intended public API surface?**
Four layers, each explicit in code:
- **Python root API**: the 17-symbol `__all__` (`lib/crewai/src/crewai/__init__.py:187-205`) — `Agent`, `Task`, `Crew`, `Flow`, `LLM`, `BaseLLM`, `Process`, `Memory`, `Knowledge`, `CrewOutput`, `TaskOutput`, `RuntimeState`, `ExecutionContext`, `PlanningConfig`, `LLMGuardrail`, `Entity`, `__version__`.
- **Extension subpackages**: `crewai.project` (`lib/crewai/src/crewai/project/__init__.py:22-40`), `crewai.tools` (`lib/crewai/src/crewai/tools/__init__.py:3-7`), `crewai.events` (`lib/crewai/src/crewai/events/__init__.py:14-16`), `crewai.hooks` (`lib/crewai/src/crewai/hooks/__init__.py:3-30`), `crewai.a2a`, `crewai.state`.
- **CLI**: the `crewai` console script (`lib/cli/pyproject.toml:34-35`) with scaffold/run/train/replay/deploy/auth commands (`lib/cli/src/crewai_cli/cli.py:115-1105`).
- **Enterprise HTTP API**: OpenAPI-documented endpoints (`docs/edge/en/api-reference/kickoff.mdx:6`).
Separate distributions (`crewai-tools`, `crewai-files`, `crewai-core`) are published packages with their own versioned identity; `crewai-devtools` is explicitly private (`lib/devtools/pyproject.toml:10-11`).

**2. Is the stable API easy to distinguish from internal implementation details?**
Partially. The root `__all__` and the small, single-purpose subpackage `__init__` files (`crewai.tools` exports exactly three names) draw a clear line, and `crewai.experimental` is labeled "may change without major-version bumps" (`lib/crewai/src/crewai/experimental/__init__.py:1`) with a runtime env-var gate for its riskiest feature (`lib/crewai/src/crewai/experimental/skills/_flag.py:9-24`). However, several boundaries leak: `crewai.agents` re-exports executor/cache/parser internals in `__all__` (`lib/crewai/src/crewai/agents/__init__.py:12-19`); `crewai.utilities.*`, `crewai.telemetry.*`, and `crewai.plus_api` are plain-importable with no underscore convention or `__getattr__` guard; and the root `__init__` itself reaches into `sys.modules` to mutate eight modules' namespaces (`lib/crewai/src/crewai/__init__.py:128-145`), which means importing `crewai` visibly changes other modules' global namespaces. There is no formal stability policy (semver promise, deprecation policy doc) in the repo — only per-field `DeprecationWarning`s with removal targets (`lib/crewai/src/crewai/task.py:561-566`, `lib/crewai/src/crewai/agent/core.py:368-371`) and the current pre-release version `1.14.8a2` (`lib/crewai/src/crewai/__init__.py:51`).

**3. Does the API expose the right level of abstraction for agent harness users?**
Yes for the primary personas. Application developers get declarative composition (`Crew`/`Flow`/`Task` plus `crewai.project` decorators and `load_crew` for JSON-configured crews, `lib/crewai/src/crewai/project/__init__.py:22-36`); extension authors get three well-scoped ABCs/decorators — `BaseLLM.call()` (`lib/crewai/src/crewai/llms/base_llm.py:283-284`), `BaseTool`/`@tool` (`lib/crewai/src/crewai/tools/base_tool.py:103`, `:653-676`) — plus a typed event bus (`lib/crewai/src/crewai/events/event_bus.py:244,570`) and hook registration (`lib/crewai/src/crewai/hooks/__init__.py:96`). Operators get a CLI covering scaffold-through-deploy (`docs/edge/en/concepts/cli.mdx:35-370`). One abstraction concern: `RuntimeState`/`Entity` and the `model_rebuild` machinery in the root `__init__` (`lib/crewai/src/crewai/__init__.py:152-176`) are exported publicly but only make sense given Pydantic serialization internals, and `CrewAgentExecutor` is deprecated yet still surfaced through `crewai.agents` (`lib/crewai/src/crewai/agents/__init__.py:14-16`, deprecation at `lib/crewai/src/crewai/agent/core.py:158-161`).

**4. Are examples sufficient to use the API correctly without reading internals?**
Mostly yes. Docs are versioned per release with frozen snapshots and CI-enforced link integrity (`AGENTS.md:1-40`), and the key extension paths each have runnable examples: quickstart Flow+crew (`docs/edge/en/quickstart.mdx:31-100`), event listeners (`docs/edge/en/concepts/event-listener.mdx:31-60`), custom LLMs including required-method contracts (`docs/edge/en/learn/custom-llm.mdx:109-148`), and custom tools (`docs/edge/en/concepts/tools.mdx:168-315`). Gaps: there is no generated Python API reference at all — the `api-reference` section documents only the enterprise HTTP API (`docs/edge/en/api-reference/kickoff.mdx:6`) — so signature-level truth lives in source/docstrings; and the docs do not state which parts of `crewai.agents`/`crewai.utilities` are safe to import, leaving users to infer boundaries from `__all__`.

## Architectural Decisions

- **Monorepo with published sub-packages**: `crewai` depends on `crewai-core` and `crewai-cli` as pinned workspace members (`lib/crewai/pyproject.toml:14-15`, `pyproject.toml:227-244`), while heavier capabilities (`crewai-tools`) are optional extras (`lib/crewai/pyproject.toml` `[project.optional-dependencies] tools`). This keeps the core install small and makes the CLI a first-class, separable surface.
- **Curated root namespace with lazy heavy imports**: only 17 symbols at top level; `Memory` (which pulls lancedb) is resolved through module `__getattr__` (`lib/crewai/src/crewai/__init__.py:53-66`), and event-type classes lazy-load to avoid importing ~12 Pydantic modules (`lib/crewai/src/crewai/events/__init__.py:8-11`).
- **Explicit experimental tier**: a dedicated `crewai.experimental` package with a docstring stability disclaimer (`lib/crewai/src/crewai/experimental/__init__.py:1`) and an opt-in env gate (`CREWAI_EXPERIMENTAL=1`) that raises a typed error when unset (`lib/crewai/src/crewai/experimental/skills/_flag.py:12-24`).
- **Docs-as-versioned-artifact**: every release freezes `docs/edge/` into `docs/vX.Y.Z/` with CI guards (`AGENTS.md:1-40`), so documented API matches shipped versions rather than main.
- **CLI as separate package with lazy command dispatch**: `crewai_cli` shims command classes behind `__new__` factories to keep CLI startup fast (`lib/cli/src/crewai_cli/cli.py:63-105`).
- **Stable-path aliasing over breakage**: `crewai.plus_api` is kept as a re-export "stable import path" while pointing new code at `crewai_core.plus_api` (`lib/crewai/src/crewai/plus_api.py:1-11`) — a deliberate compatibility shim during the core-package extraction.

## Notable Patterns

- **Pydantic-native API design**: public classes are Pydantic models (`BaseLLM(BaseModel, ABC)` at `lib/crewai/src/crewai/llms/base_llm.py:129`; `BaseTool(BaseModel, ABC)` at `lib/crewai/src/crewai/tools/base_tool.py:103`), giving users validation, serialization, and `deprecated=` field markers (`lib/crewai/src/crewai/agent/core.py:227-285`) for free.
- **Registry-based polymorphic deserialization**: `BaseTool.__init_subclass__` populates `_TOOL_TYPE_REGISTRY` so serialized `tool_type` dicts resolve to concrete tool classes (`lib/crewai/src/crewai/tools/base_tool.py:47-48, 57-77`) — an extension contract that third-party tools inherit implicitly.
- **Decorator + JSON dual definition**: the same crew can be defined in Python via `@crew`-style decorators (`lib/crewai/src/crewai/project/__init__.py:8-20`) or loaded from JSONC files (`load_crew`, `load_agent`, `strip_jsonc_comments`, `lib/crewai/src/crewai/project/__init__.py:31-36`), and the quickstart teaches the JSON path (`docs/edge/en/quickstart.mdx:31-100`).
- **Eager class exports, lazy SDK imports in `crewai-tools`**: 100+ tools are importable from one module (`lib/crewai-tools/src/crewai_tools/__init__.py:1-333`), but each tool defers its third-party SDK import into methods (`lib/crewai-tools/src/crewai_tools/tools/selenium_scraping_tool/selenium_scraping_tool.py:75`), so optional extras don't break the base import.
- **Deprecation-with-deadline**: warnings name the removal version ("removed in CrewAI v1.0.0", `lib/crewai/src/crewai/task.py:561-566`; "removed in v2.0", `lib/crewai/src/crewai/agent/core.py:368`), giving users concrete migration horizons.

## Tradeoffs

- **Import-time hygiene vs. user quietness**: `_suppress_pydantic_deprecation_warnings` globally monkey-patches `warnings.warn` when `crewai` is imported (`lib/crewai/src/crewai/__init__.py:28-49`). This keeps user output clean but also silences Pydantic deprecations the user's *own* code may trigger — the library trades user observability for polish.
- **Forward-reference resolution vs. namespace pollution**: to make Pydantic models rebuild with cross-module types, the root `__init__` injects `_resolve_namespace` (including every symbol in `BaseAgent`'s module `__dict__`) into eight modules (`lib/crewai/src/crewai/__init__.py:128-145`). This avoids typing failures at the cost of unpredictable module globals — a new integration that shadows any injected name can hit subtle conflicts.
- **Small curated root vs. discoverability**: the 17-symbol `__all__` is easy to learn but pushes everything else into subpackages; combined with the absence of a generated Python API reference, discoverability depends on the concept guides (`docs/edge/en/concepts/*.mdx`) rather than a canonical symbol index.
- **One CLI name, two packages**: both `crewai` and `crewai-cli` register the `crewai` script (`lib/crewai/pyproject.toml:146-147`; `lib/cli/pyproject.toml:34-35`). Convenient for end users, but it means the framework package and CLI package can race for the same entry point in constrained environments.
- **Versioned docs vs. maintenance cost**: freezing `docs/vX.Y.Z/` per release gives users stable API documentation (`AGENTS.md:1-40`) at the price of an append-only images policy and CI snapshot guards.

## Failure Modes / Edge Cases

- **Silent degradation on rebuild failure**: if the forward-ref rebuild block fails, the root `__init__` logs a warning and sets `RuntimeState = None` (`lib/crewai/src/crewai/__init__.py:178-185`). Users importing `RuntimeState` get `None` rather than an import error, deferring failure to first use.
- **Experimental feature surprise**: using Skills without `CREWAI_EXPERIMENTAL=1` raises `ExperimentalFeatureDisabledError` at call time (`lib/crewai/src/crewai/experimental/skills/_flag.py:16-24`) — a runtime, not import-time, failure that surfaces only when the feature is exercised.
- **Lazy attribute errors**: `from crewai import Memory` works, but direct `crewai.Memory` access before first `__getattr__` resolution depends on the lazy path (`lib/crewai/src/crewai/__init__.py:58-66`); typos raise `AttributeError` with a clear message, but static analyzers cannot see lazily-exported names (they are absent from module globals and `__all__` lists them anyway — `Memory` is in `__all__` at line 198 while only defined under `TYPE_CHECKING` at line 25).
- **Deprecated executor still reachable**: `CrewAgentExecutor` warns at construction ("deprecated and will be removed in a future release", `lib/crewai/src/crewai/agent/core.py:158-161`) but remains in `crewai.agents.__all__` (`lib/crewai/src/crewai/agents/__init__.py:14-16`), so new code can still adopt a removal-bound API.
- **Namespace-injection collisions**: because `sys.modules[_BaseAgent.__module__].__dict__` is bulk-updated with `_resolve_namespace` (`lib/crewai/src/crewai/__init__.py:128-131,135-145`), any module-level name in `base_agent.py` becomes visible in `Agent`, `Crew`, `Flow`, `Task`, and executor modules — a refactor that adds a common name (e.g., `Process`) to one of those modules would silently shadow the public `crewai.Process` inside those namespaces.

## Future Considerations

- **Publish a generated Python API reference**: mkdocstrings-style pages derived from `__all__` + subpackage exports would close the largest documentation gap; today the only signature-level truth is source code (`docs/edge/en/api-reference/` covers only the enterprise HTTP surface, `docs/edge/en/api-reference/kickoff.mdx:6`).
- **Codify a stability policy**: a short SEMVER/deprecation policy document (what `crewai.utilities` guarantees, when `experimental` graduates) would convert the ad-hoc markers (`lib/crewai/src/crewai/experimental/__init__.py:1`, `lib/crewai/src/crewai/task.py:561-566`) into a contract.
- **Shrink the internal-export leak**: move `CacheHandler`, `ToolsHandler`, `parse` out of `crewai.agents.__all__` (`lib/crewai/src/crewai/agents/__init__.py:12-19`) or mark them internal once checkpoint/serialization callers migrate.
- **Harden surface tests**: extend `lib/crewai/tests/test_imports.py:4-13` to assert the full `__all__` (importability, absence of `None` after failed rebuilds, lazy-name resolution) — cheap insurance for the namespace-injection machinery.
- **Retire the warnings monkey-patch**: replace the global `warnings.warn` filter (`lib/crewai/src/crewai/__init__.py:28-49`) with targeted Pydantic config or upstream fixes once the pinned Pydantic range (`pyproject.toml` root dev deps) allows.

## Questions / Gaps

- **No explicit semver/stability statement found.** Searched `README.md`, `lib/crewai/README.md`, and `AGENTS.md` for "semantic", "stability", "breaking" — only marketing-level claims ("production-grade standards", `README.md:800`) and docs-versioning rules. No evidence of a formal public-API stability policy beyond per-call deprecation warnings.
- **No generated Python API reference found.** `docs/edge/en/api-reference/` contains only enterprise HTTP endpoint pages bound to OpenAPI (`docs/edge/en/api-reference/kickoff.mdx:1-8`); no mkdocstrings/sphinx configuration exists in the repo (searched for mkdocs/sphinx/mkdocstrings configs — none found).
- **Whether `crewai.utilities`/`crewai.telemetry` are considered public could not be determined from code.** They lack `__all__`-driven curation and underscore prefixes; no doc page references them as supported entry points. Treated as internal-by-convention only.
- **Enterprise HTTP API surface not fully audited.** The OpenAPI specs (`docs/edge/en/enterprise-api.*.yaml`) define the operator-facing endpoints, but this study focused on the Python/CLI surface; a dedicated dimension pass over the YAML specs would be needed for endpoint-level claims.

---

Generated by `24.01-public-api-surface` against `crewai`.
