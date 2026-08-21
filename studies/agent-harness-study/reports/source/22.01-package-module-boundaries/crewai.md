# Source Analysis: crewai

## Package and Module Boundaries

### Source Info

| Field | Value |
|-------|-------|
| Name | crewai |
| Path | `studies/agent-harness-study/sources/crewai` |
| Language / Stack | Python (>=3.10, <3.14); uv-workspace of 6 packages |
| Analyzed | 2026-08-21 |

## Summary

`crewai` is a uv-workspace monorepo that splits the project into six independently versioned, hatchling-built packages — `crewai` (runtime), `crewai-core` (shared utilities), `crewai-cli` (CLI), `crewai-files` (multimodal file processing), `crewai-tools` (tool catalog), and `crewai-devtools` (release tooling) — defined in `pyproject.toml:227-244`. Dependency direction is mostly one-way: `crewai-core` and `crewai-files` are pure leaves with zero internal Python imports (`lib/crewai-core/src/crewai_core/`, `lib/crewai-files/src/crewai_files/` contain only `crewai_core.X` / local references), `crewai-cli` consumes both `crewai-core` and the runtime `crewai`, `crewai` depends on `crewai-core` (58 import sites, see evidence table) plus `crewai-cli`, and `crewai-tools` depends on the runtime `crewai`. Internal circular imports are not avoided but actively broken with documented lazy `__getattr__` resolvers (`crewai/__init__.py:53-66` for `Memory → lancedb`, `crewai/experimental/__init__.py:7-73` for Flow/TouchTask deadlock, `crewai/flow/runtime/__init__.py:162-168` for `ExecutionContext` cycle), dedicated event types files (`crewai/events/types/agent_events.py:1` reads "Agent-related events moved to break circular dependencies"), and an explicit `CircularDependencyError` (`crewai/events/handler_graph.py:16-26`). The answer to the dimension question — "Can you use the tool system without pulling in the entire runtime?" — is **no**: `crewai-tools` declares `crewai==1.14.8a2` as a hard dependency (`lib/crewai-tools/pyproject.toml:13`), and most tool classes import `from crewai.tools import BaseTool` directly. Public/internal separation is reasonably explicit — `crewai/__init__.py:187-205` declares an 18-name `__all__`, PEP 561 `py.typed` markers ship in every published package (`lib/crewai/src/crewai/py.typed`, `lib/crewai-core/src/crewai_core/py.typed`, etc.), and `crewai/experimental/__init__.py:1` flags its surface as unstable — but internal helpers like `_serialize_llm_ref` and `_validate_llm_ref` leak across modules (29 import sites, see evidence table), and four core files (`agent/core.py`, `crew.py`, `task.py`, `llm.py`) each exceed 1400 lines, undermining intra-package module boundaries.

## Rating

**7 / 10** — Clear model at the workspace level (six well-separated packages, unidirectional internal flow, PEP 561 typed, lazy-import cycle breakers, `CircularDependencyError`), with explicit `__all__` and a deprecation shim (`crewai/cli/__init__.py`) that proves the team enforces boundaries across releases. The model breaks down on the dimension's central question: the tool system cannot be loaded independently of the entire runtime (`crewai-tools` pins `crewai==1.14.8a2`), and the giant core files (agent/core.py 1977 LOC, crew.py 2374 LOC, task.py 1463 LOC, llm.py 2674 LOC, flow/runtime/__init__.py 3867 LOC) suggest the team has not extended its clean package boundaries down into intra-package module boundaries.

## Evidence Collected

Every entry includes a file path with line numbers. Format: `path/to/file.py:NN`.

| Area | Evidence | File:Line |
|------|----------|-----------|
| Workspace declares 6 members | `[tool.uv.workspace].members = ["lib/crewai", "lib/crewai-tools", "lib/devtools", "lib/crewai-files", "lib/cli", "lib/crewai-core"]` | `pyproject.toml:227-235` |
| Workspace sources mapped explicitly | `[tool.uv.sources]` covers all six packages | `pyproject.toml:238-244` |
| Each package has its own pyproject + hatchling backend | `[build-system] requires = ["hatchling"]` repeated 5x | `lib/crewai/pyproject.toml:149-151`, `lib/crewai-core/pyproject.toml:30-32`, `lib/crewai-tools/pyproject.toml:155-157`, `lib/crewai-files/pyproject.toml:25-27`, `lib/cli/pyproject.toml:37-39` |
| Runtime package version 1.14.8a2 in five places | `crewai-core==1.14.8a2`, `crewai-cli==1.14.8a2`, `crewai-tools==1.14.8a2` | `lib/crewai/pyproject.toml:11-12, 58` |
| `crewai-core` is a leaf with no internal deps | Only `crewai_core.X` imports across 23 sites | `lib/crewai-core/src/crewai_core/` (grep) |
| `crewai-files` is a leaf with no internal deps | Zero matches for `^from crewai\.\|^import crewai` | `lib/crewai-files/src/crewai_files/` (grep) |
| `crewai-devtools` is fully isolated | Zero matches for any `crewai` import | `lib/devtools/` (grep) |
| `crewai-core` deps are infrastructure-only | appdirs, cryptography, httpx, packaging, portalocker, pyjwt, pydantic, rich, opentelemetry-*, tomli | `lib/crewai-core/pyproject.toml:10-23` |
| `crewai-files` deps are infrastructure-only | Pillow, pypdf, python-magic, aiocache, aiofiles, tinytag, av | `lib/crewai-files/pyproject.toml:10-18` |
| `crewai-cli` depends on `crewai-core` and runtime | `crewai-core==1.14.8a2` plus standard infra | `lib/cli/pyproject.toml:11-27` |
| `crewai` depends on `crewai-core` (used 58 times in code) | `crewai_core.printer`, `.lock_store`, `.paths`, `.settings`, `.auth`, `.token_manager`, `.plus_api`, `.project`, `.tool_credentials`, `.user_data`, `.constants`, `.version` | `lib/crewai/src/crewai/` (grep, 58 matches) |
| `crewai` declares `crewai-cli` as runtime dep | `crewai-cli==1.14.8a2` is required at install | `lib/crewai/pyproject.toml:12` |
| `crewai-tools` depends on the runtime `crewai` | `crewai==1.14.8a2` is a hard requirement | `lib/crewai-tools/pyproject.toml:13` |
| `crewai-tools` widely imports `from crewai.tools import BaseTool` | 100+ import sites across tool classes | `lib/crewai-tools/src/crewai_tools/` (grep) |
| `crewai-cli` imports runtime only for project loader + deploy + templates | 20 import sites, mostly in `deploy/`, `templates/`, and `tests/` | `lib/cli/src/crewai_cli/` (grep) |
| Public API explicitly listed in `crewai/__init__.py` | `__all__` lists 18 names: LLM, Agent, BaseLLM, Crew, CrewOutput, Entity, ExecutionContext, Flow, Knowledge, LLMGuardrail, Memory, PlanningConfig, Process, RuntimeState, Task, TaskOutput, __version__ | `lib/crewai/src/crewai/__init__.py:187-205` |
| Public types are eagerly imported at runtime | `from crewai.agent.core import Agent`, `from crewai.crew import Crew`, `from crewai.flow.flow import Flow`, `from crewai.task import Task`, etc. | `lib/crewai/src/crewai/__init__.py:8-21` |
| `Memory` is lazy-loaded to defer lancedb import | `_LAZY_IMPORTS["Memory"] = ("crewai.memory.unified_memory", "Memory")` + `__getattr__` | `lib/crewai/src/crewai/__init__.py:53-66` |
| `crewai/experimental` flags instability | docstring "Experimental CrewAI surface — APIs here may change without major-version bumps" | `lib/crewai/src/crewai/experimental/__init__.py:1` |
| `crewai/cli` deprecation shim | Module-level `DeprecationWarning`, `_ShimFinder` inserts into `sys.meta_path` to map `crewai.cli[.X]` onto `crewai_cli[.X]` | `lib/crewai/src/crewai/cli/__init__.py:1-74` |
| `crewai/cli` exports no symbols — only forwards imports | Empty `__init__.py` apart from shim logic | `lib/crewai/src/crewai/cli/__init__.py:74` |
| Internal helpers leaked across modules | `_serialize_llm_ref` / `_validate_llm_ref` (29 import sites), `_INITIAL_STATE_CLASS_MARKER` (3 sites) | `lib/crewai/src/crewai/agents/agent_builder/base_agent.py:82, 164`; `lib/crewai/src/crewai/agents/crew_agent_executor.py:30`; `lib/crewai/src/crewai/agent/core.py:52-53`; `lib/crewai/src/crewai/agent/planning_config.py:7`; `lib/crewai/src/crewai/crew.py:63-64`; `lib/crewai/src/crewai/flow/flow.py:23, 41`; `lib/crewai/src/crewai/flow/runtime/__init__.py:322-349` |
| `crewai/events/types/agent_events.py` exists to break cycles | docstring "Agent-related events moved to break circular dependencies" | `lib/crewai/src/crewai/events/types/agent_events.py:1` |
| `crewai/events/types/logging_events.py` exists to break cycles | docstring "Agent logging events that don't reference BaseAgent to avoid circular imports" | `lib/crewai/src/crewai/events/types/logging_events.py:1` |
| `crewai/flow/runtime/__init__.py` aliases `ExecutionContext = Any` to avoid init cycle | comment "we can't import it at runtime because `crewai.context` is loaded mid-initialization when this module is imported through `crewai.__init__` (circular)" | `lib/crewai/src/crewai/flow/runtime/__init__.py:162-168` |
| `crewai/flow/persistence/sqlite.py` uses late import | inline comment "Import here to avoid circular imports" | `lib/crewai/src/crewai/flow/persistence/sqlite.py:259` |
| `crewai/experimental/agent_executor.py` uses lazy properties | "Lazily create the StepExecutor (avoids circular imports)", "Lazily create the PlannerObserver (avoids circular imports)" | `lib/crewai/src/crewai/experimental/agent_executor.py:453, 472` |
| `crewai/a2a/config.py` lives apart from experimental.a2a to break cycle | docstring "This module is separate from experimental.a2a to avoid circular imports" | `lib/crewai/src/crewai/a2a/config.py:3` |
| `crewai/events/__init__.py` lazy-resolves every event type | `_LAZY_EVENT_MAPPING` (97 entries) + `__getattr__` defers 12 Pydantic model modules | `lib/crewai/src/crewai/events/__init__.py:154-277` |
| `crewai/experimental/__init__.py` lazy-resolves evaluation + executor | `_LAZY_FROM_AGENT_EXECUTOR`, `_LAZY_FROM_EVALUATION` and `__getattr__` | `lib/crewai/src/crewai/experimental/__init__.py:24-73` |
| `crewai/events/handler_graph.py` defines `CircularDependencyError` | `class CircularDependencyError(Exception)` | `lib/crewai/src/crewai/events/handler_graph.py:16-26` |
| `crewai/events/event_bus.py` detects handler cycles | `crewai_event_bus` raises `CircularDependencyError` | `lib/crewai/src/crewai/events/event_bus.py:814-818` |
| Test for circular-dependency detection | `test_circular_dependency_detection` in depends | `lib/crewai/tests/events/test_depends.py:195-196` |
| LLM providers live in a dedicated sub-package | `crewai/llms/providers/{anthropic,openai,azure,bedrock,gemini,snowflake,openai_compatible}/completion.py` | `lib/crewai/src/crewai/llms/providers/` (ls) |
| LLM provider modules fail loudly if SDK missing | `try: from anthropic import ...; except ImportError: raise ImportError("Anthropic native provider not available, to install: uv add 'crewai[anthropic]'")` | `lib/crewai/src/crewai/llms/providers/anthropic/completion.py:25-38` |
| Flow split into 3 concerns | `flow/dsl` (authoring), `flow/flow_definition` (serializable contract), `flow/runtime` (engine) | `lib/crewai/src/crewai/flow/flow.py:1-14` |
| Flow public class is a 53-line facade over split modules | `class Flow(_ConversationalMixin, RuntimeFlow[T])` plus re-exports | `lib/crewai/src/crewai/flow/flow.py:36-52` |
| `flow/dsl/` further split | `_conditions.py`, `_human_feedback.py`, `_listen.py`, `_router.py`, `_start.py`, `_types.py`, `_utils.py` | `lib/crewai/src/crewai/flow/dsl/` (ls) |
| Agent adapters split by external framework | `agents/agent_adapters/{langgraph,openai_agents}/` | `lib/crewai/src/crewai/agents/agent_adapters/` (ls) |
| BaseAgentAdapter defines the adapter contract | abstract base class with `configure_tools` and `configure_structured_output` | `lib/crewai/src/crewai/agents/agent_adapters/base_agent_adapter.py:15-46` |
| CrewAI public types eagerly model_rebuild | `_BaseAgentExecutor.model_rebuild(force=True, _types_namespace=_full_namespace)` × 8 | `lib/crewai/src/crewai/__init__.py:152-159` |
| PydanticUserError fallback on rebuild failure | `except (ImportError, PydanticUserError): ... logger.warning("model_rebuild() failed; forward refs may be unresolved.")` | `lib/crewai/src/crewai/__init__.py:175-185` |
| PEP 561 typed markers | `py.typed` empty files in 4 packages | `lib/crewai/src/crewai/py.typed`, `lib/crewai-core/src/crewai_core/py.typed`, `lib/crewai-tools/src/crewai_tools/py.typed`, `lib/crewai-files/src/crewai_files/py.typed` |
| Ruff bans relative imports | `[tool.ruff.lint.flake8-tidy-imports] ban-relative-imports = "all"` | `pyproject.toml:90-91` |
| Strict mypy with pydantic plugin | `strict = true`, `disallow_untyped_defs = true`, `plugins = ["pydantic.mypy"]` | `pyproject.toml:120-131` |
| Hatch wheel packages config | `crewai-core`: `packages = ["src/crewai_core"]`; `crewai-cli`: `packages = ["src/crewai_cli"]` | `lib/crewai-core/pyproject.toml:37-38`; `lib/cli/pyproject.toml:44-45` |
| Smoke test for leaf modules | `crewai-core/tests/test_smoke.py` covers version, paths, constants, printer, lock, user_data, telemetry | `lib/crewai-core/tests/test_smoke.py:20-130` |
| Minimal import smoke test for runtime | `tests/test_imports.py` checks `from crewai import TaskOutput` and `from crewai import CrewOutput` | `lib/crewai/tests/test_imports.py:1-15` |
| Optional-dependencies separation test (skipped in CI) | `test_no_optional_dependencies_in_init` builds an isolated uv env and runs `import crewai_tools` | `lib/crewai-tools/tests/test_optional_dependencies.py:34-46` |
| `crewai-files` has its own test suite | `tests/test_file_url.py`, `test_resolved.py`, `test_resolver.py`, `test_upload_cache.py` | `lib/crewai-files/tests/` (ls) |
| Tool specs generated as data | `tool.specs.json` (932 KB) alongside `crewai_tools` package | `lib/crewai-tools/tool.specs.json` |
| CLI tests live in `crewai-cli` package | `tests/test_create_crew.py`, `tests/test_install_crew.py`, `tests/test_crew_run_tui.py`, etc. | `lib/cli/tests/` (ls) |
| Workspace pytest root points at 5 packages | `[tool.pytest.ini_options].testpaths = ["lib/crewai/tests", "lib/crewai-tools/tests", "lib/crewai-files/tests", "lib/cli/tests", "lib/crewai-core/tests"]` | `pyproject.toml:142-148` |

## Answers to Dimension Questions

1. **Are modules cleanly separated?** Mostly yes at the package level, no at the intra-package level. Six hatchling-built wheels exist with distinct names, `py.typed` markers, and explicit `__all__` lists. Within `crewai`, the four heaviest modules — `agent/core.py` (1977 LOC), `crew.py` (2374 LOC), `task.py` (1463 LOC), `llm.py` (2674 LOC) — are monoliths that import widely from `agents/`, `tools/`, `events/`, `memory/`, `knowledge/`, `flow/`, and `security/` (`lib/crewai/src/crewai/agent/core.py:1-100`). Internal helpers like `_serialize_llm_ref` / `_validate_llm_ref` are exported from `agents/agent_builder/base_agent.py:82, 164` and re-imported by `agent/core.py:52-53`, `agent/planning_config.py:7`, `agents/crew_agent_executor.py:30`, `crew.py:63-64`. Flow is the only subsystem that has actually been refactored into separate concerns (`flow/dsl`, `flow/runtime`, `flow/flow_definition`).

2. **Do dependencies flow in one direction?** Yes at the package level. `crewai-core` and `crewai-files` are leaves (zero internal imports). `crewai-cli` → `crewai-core` + `crewai` (for project loader / templates / deploy). `crewai` → `crewai-core` + `crewai-cli`. `crewai-tools` → `crewai`. `crewai-devtools` has no internal deps (`lib/devtools/`). Confirmed by `grep '^from crewai\.\|^import crewai'` returning zero matches in `lib/crewai-core/`, `lib/crewai-files/`, and `lib/devtools/`. Within the runtime, direction is mostly one-way — `agents/agent_builder` is a leaf for `BaseAgent` and `BaseAgentExecutor`, `tools/base_tool.py` is the leaf for `BaseTool`, `llms/base_llm.py` is the parent for all providers, `events/event_bus.py` is the central bus — but the runtime core files re-import from many places to wire everything together.

3. **Can modules be used independently?** Partially. `crewai-core` can be installed and used on its own (no `crewai` deps, has its own smoke tests at `lib/crewai-core/tests/test_smoke.py:1-130`). `crewai-files` can be installed on its own (zero `crewai` imports). `crewai-cli` cannot be used without the runtime `crewai` (it imports `from crewai.project import CrewBase, crew` at `lib/cli/src/crewai_cli/deploy/archive.py:249-250`). **The tool system cannot be used without pulling in the entire runtime**: `crewai-tools` pins `crewai==1.14.8a2` (`lib/crewai-tools/pyproject.toml:13`) and 100+ tool classes import `from crewai.tools import BaseTool` (e.g. `lib/crewai-tools/src/crewai_tools/tools/arxiv_paper_tool/arxiv_paper_tool.py:11`). Conversely, the runtime `crewai` imports nothing from `crewai-tools` — tool discovery is one-way.

4. **Are public APIs distinguished from internal ones?** Yes for the top-level package, inconsistently within. `crewai/__init__.py:187-205` declares an 18-name `__all__`; `crewai/events/__init__.py:280-383` lists 96 names; `crewai/experimental/__init__.py:76-102` lists 26 names; `crewai/tools/__init__.py` and others follow the same pattern. `crewai/experimental/__init__.py:1` explicitly flags that surface as unstable. The `crewai/cli/__init__.py:1-74` shim proves the team tracks and enforces cross-package boundaries across releases. PEP 561 `py.typed` markers ship in every published package (`lib/crewai/src/crewai/py.typed`, `lib/crewai-core/src/crewai_core/py.typed`, `lib/crewai-tools/src/crewai_tools/py.typed`, `lib/crewai-files/src/crewai_files/py.typed`). Ruff enforces `ban-relative-imports = "all"` (`pyproject.toml:90-91`) and mypy runs in `strict` mode (`pyproject.toml:120-131`). However, the internal helpers `_serialize_llm_ref`, `_validate_llm_ref`, `_INITIAL_STATE_CLASS_MARKER`, `_coerce_checkpoint`, etc. are leaked across module boundaries with no `__all__` enforcement inside `agents/agent_builder/base_agent.py` — they have leading underscores but are imported by 4+ other modules (`lib/crewai/src/crewai/agent/core.py:52-53`; `lib/crewai/src/crewai/agents/crew_agent_executor.py:30`; `lib/crewai/src/crewai/crew.py:63-64`; `lib/crewai/src/crewai/flow/runtime/__init__.py:322-349`).

## Architectural Decisions

- **uv-workspace with hatchling wheels** instead of a single fat wheel. The team chose `[tool.uv.workspace].members` with 6 entries (`pyproject.toml:227-235`) and per-package `[build-system] requires = ["hatchling"]` (`lib/crewai/pyproject.toml:149-151`, `lib/crewai-core/pyproject.toml:30-32`, `lib/crewai-tools/pyproject.toml:155-157`, `lib/crewai-files/pyproject.toml:25-27`, `lib/cli/pyproject.toml:37-39`). Each package pins the same version `1.14.8a2` to keep releases synchronized (`lib/crewai/pyproject.toml:11-12, 58`).
- **Promote shared infra to `crewai-core`** so the runtime, the CLI, and `crewai-tools` can all consume the same `printer`, `lock_store`, `paths`, `settings`, `auth`, `plus_api`, `telemetry`, `user_data`, `version`, and `constants` without re-implementing them. 58 import sites in `crewai/` (see evidence table), 28 in `cli/`. `crewai-core` is a pure leaf that depends only on infrastructure packages (`lib/crewai-core/pyproject.toml:10-23`).
- **Promote multimodal file handling to `crewai-files`** so the heavy Pillow + pypdf + av + tinytag stack is opt-in via `crewai[file-processing]` (`lib/crewai/pyproject.toml:112-114`, `lib/crewai-files/pyproject.toml:10-18`). `crewai-files` is a pure leaf with zero internal imports.
- **Extract the CLI to `crewai-cli`** and keep a deprecation shim in `crewai/cli/__init__.py` that intercepts `from crewai.cli.X import Y` via a `_ShimFinder` inserted into `sys.meta_path` (`lib/crewai/src/crewai/cli/__init__.py:43-74`). The shim emits a `DeprecationWarning` (`lib/crewai/src/crewai/cli/__init__.py:23-27`) and forwards to `crewai_cli`. The runtime `crewai` still declares `crewai-cli==1.14.8a2` as a hard dep (`lib/crewai/pyproject.toml:12`) so the `crewai` console script remains wired up after the split.
- **LLM providers live in `crewai/llms/providers/<vendor>/completion.py`** and each module wraps its vendor SDK import in `try/except ImportError: raise ImportError("Anthropic native provider not available, to install: uv add 'crewai[anthropic]'")` (`lib/crewai/src/crewai/llms/providers/anthropic/completion.py:25-38`). The extras `anthropic`, `aws`, `bedrock`, `gemini`, `voyageai`, `watson`, `litellm`, `azure-ai-inference` are declared in `lib/crewai/pyproject.toml:80-105`.
- **Flow split into `dsl/` (decorators), `flow_definition.py` (serializable contract), `runtime/` (engine)** with a 53-line `flow/flow.py` re-export facade (`lib/crewai/src/crewai/flow/flow.py:1-52`). The public `Flow` class composes a `_ConversationalMixin` over the runtime engine (`lib/crewai/src/crewai/flow/flow.py:36-37`).
- **Event types live in `crewai/events/types/`** and are lazy-resolved by `crewai/events/__init__.py:154-277` via a 97-entry `_LAZY_EVENT_MAPPING` table to avoid loading all Pydantic models at package import. `crewai/events/types/agent_events.py:1` documents that the module exists specifically to break circular dependencies.
- **Circular-dependency handling is explicit and tested**. `crewai/events/handler_graph.py:16-26` defines `CircularDependencyError`. `crewai/events/event_bus.py:814-818` raises it. `crewai/tests/events/test_depends.py:195-196` covers the path. Runtime cycles between Flow and `experimental` are broken by lazy `__getattr__` (`lib/crewai/src/crewai/experimental/__init__.py:7-73`).
- **PEP 561 typed markers** ship in 4 packages so downstream users get type checking. mypy runs in strict mode with the pydantic plugin (`pyproject.toml:120-131`).
- **Ruff enforces explicit imports** via `ban-relative-imports = "all"` (`pyproject.toml:90-91`).

## Notable Patterns

- **Lazy attribute resolution** with `__getattr__` is used in 3 places to defer heavy imports: `crewai/__init__.py:53-66` (defers lancedb via `Memory`), `crewai/events/__init__.py:154-277` (defers 97 event types), `crewai/experimental/__init__.py:24-73` (defers evaluation + executor).
- **Lazy property factories** for circular-import avoidance: `crewai/experimental/agent_executor.py:453, 472` uses property-based lazy creation for `StepExecutor` and `PlannerObserver`.
- **Meta-path import shim** for the CLI deprecation: `crewai/cli/__init__.py:30-69` installs a `_ShimFinder` that maps `crewai.cli[.X]` onto `crewai_cli[.X]`.
- **Pydantic namespace plumbing**: `crewai/__init__.py:128-145` injects a `_resolve_namespace` dict into the `__dict__` of multiple modules so cross-module `model_rebuild` calls can resolve forward refs (`crewai/__init__.py:152-159`).
- **Provider packages with guarded imports** rather than dynamic dispatch: each provider is its own module with a hard `try/except` at the top, so a missing SDK is a loud install-time error rather than a runtime AttributeError.
- **Workspace-pinned versions**: all six `pyproject.toml` files declare `1.14.8a2` to keep releases coordinated (`lib/crewai/pyproject.toml:11-12, 58`).
- **Per-package test suites**: each workspace member has its own `tests/` directory and the root `[tool.pytest.ini_options].testpaths` lists them all (`pyproject.toml:142-148`).
- **Generated data file `tool.specs.json`** (932 KB) lives alongside `crewai_tools` and is regenerated by `lib/crewai-tools/src/crewai_tools/generate_tool_specs.py:9` (`from crewai.tools.base_tool import BaseTool, EnvVar`).

## Tradeoffs

- **Tools cannot be used independently of the runtime**. `crewai-tools` requires `crewai==1.14.8a2` (`lib/crewai-tools/pyproject.toml:13`) and 100+ tool classes import `from crewai.tools import BaseTool`. The original dimension question — "Can you use the tool system without pulling in the entire runtime?" — has the answer **no**. Consumers who want just `BaseTool` must install `crewai`.
- **Several core files are 1400+ lines**, weakening intra-package module boundaries: `agent/core.py` (1977 LOC), `crew.py` (2374 LOC), `task.py` (1463 LOC), `llm.py` (2674 LOC), `flow/runtime/__init__.py` (3867 LOC). They import widely across `agents/`, `tools/`, `events/`, `memory/`, `knowledge/`, `flow/`, `security/`, `state/`, `llms/`, and `utilities/`.
- **Internal helpers leak across modules** via leading-underscore imports. `_serialize_llm_ref` and `_validate_llm_ref` are imported from 4+ modules (`lib/crewai/src/crewai/agent/core.py:52-53`; `lib/crewai/src/crewai/agents/crew_agent_executor.py:30`; `lib/crewai/src/crewai/crew.py:63-64`). They are not gated by `__all__` inside `agents/agent_builder/base_agent.py`. This means the "internal helpers" are stable in practice even though they look private.
- **Pydantic model_rebuild is fragile** — `crewai/__init__.py:152-159` calls `model_rebuild` on 8 types and falls back to a warning if it fails (`crewai/__init__.py:175-185`). Any forward-ref addition that misses the namespace dict breaks at runtime.
- **Inconsistent lazy-vs-eager imports** in `__init__.py` files. `crewai/events/__init__.py` lazy-loads all 97 types, `crewai/__init__.py` lazy-loads only `Memory`, `crewai/agents/__init__.py` lazy-loads only `CrewAgentExecutor` (`lib/crewai/src/crewai/agents/__init__.py:23-28`), and most subpackages are fully eager. There is no consistent policy.
- **`crewai-cli` re-declares most of `crewai-core.auth`** under `crewai_cli/authentication/` (`lib/cli/src/crewai_cli/authentication/`). 28 `crewai_core.X` imports in `cli/src/crewai_cli/` but the auth subpackage is re-implemented. Some drift risk.
- **CLI imports runtime `crewai` for templates and deploy** (`lib/cli/src/crewai_cli/deploy/archive.py:249-250`; `lib/cli/src/crewai_cli/templates/crew/crew.py:2-3`). The CLI therefore pulls in the runtime even when not strictly needed for `crewai version`, `crewai config`, etc.
- **Workspace-level safety overrides** (security-driven dependency pins in `[tool.uv] override-dependencies`, `pyproject.toml:201-225`) cross-cut every package and add review friction for routine dep bumps.

## Failure Modes / Edge Cases

- **Cycle regression risk**: any new module added to the runtime that imports `crewai.experimental.*` at top level risks the Flow/TouchTask deadlock that the lazy `__getattr__` in `crewai/experimental/__init__.py:7-73` was added to prevent. The docstring calls it out explicitly.
- **`crewai/cli` shim meta_path insertion**: `lib/crewai/src/crewai/cli/__init__.py:73` does `sys.meta_path.insert(0, _finder)`. Inserting at position 0 changes Python's import resolution order and could shadow other finders (e.g. importlib_hooks, sitecustomize finders) for any module starting with `crewai.cli`.
- **Heavy optional-dep leak in `crewai/rag/embeddings/providers/`**: `chromadb.utils.embedding_functions.cohere_embedding_function` is imported unconditionally at `lib/crewai/src/crewai/rag/embeddings/providers/cohere/cohere_provider.py:3` even though `cohere` is an optional extra. Users who install `crewai[cohere]` get the full ChromaDB dep chain at import time. Same pattern in `lib/crewai/src/crewai/rag/embeddings/providers/{voyageai, onnx, huggingface, ollama, jina, instructor, openai, sentence_transformer, openclip, text2vec, roboflow, google}/...` — `chromadb.utils.embedding_functions.*` is imported unconditionally even when the embedding package is not installed.
- **`crewai-events.types.agent_events.py` directly imports `BaseAgent` and `BaseTool`**. If a downstream user shims `crewai.tools.base_tool.BaseTool` or `crewai.agents.agent_builder.base_agent.BaseAgent`, the event types break at first emit.
- **`crewai/__init__.py:185` fallback**: if `model_rebuild` fails, `RuntimeState = None` is assigned. Any later code that does `if RuntimeState is None` fails closed silently with only a logger warning (`crewai/__init__.py:181-184`).
- **Optional-dependencies test is skipped in CI**: `lib/crewai-tools/tests/test_optional_dependencies.py:34` is marked `@pytest.mark.skip(reason="Test takes too long in GitHub Actions (>30s timeout) due to dependency installation")`. The only automated boundary test for the tool catalog is therefore dormant.
- **`crewai/agent/internal/meta.py:11`** imports `from pydantic._internal._model_construction import ModelMetaclass` — reaching into Pydantic's private API. Pydantic 3 will break this.

## Future Considerations

- **Split `crewai-tools` into per-tool wheels** (or at least group tools by dependency footprint) so consumers can install just the tools they need without paying the install cost of `crewai-tools`'s full transitive closure (`lib/crewai-tools/pyproject.toml:10-19`).
- **Carve `crewai.tools` into a minimal `crewai-tools-core` wheel** with just `BaseTool`, `BaseAgent`, `BaseLLM`, and the event-bus interfaces so external tool authors can depend on the contract without the runtime.
- **Break up the giant core files**: `crewai/agent/core.py` (1977 LOC), `crewai/crew.py` (2374 LOC), `crewai/task.py` (1463 LOC), `crewai/llm.py` (2674 LOC), `crewai/flow/runtime/__init__.py` (3867 LOC) are good candidates for the same `dsl + runtime + definition` split that Flow already uses.
- **Lock down `agents/agent_builder/base_agent.py` internal helpers** with an `_INTERNAL_API = {...}` namespace or move them to a `crewai/_internals/` package so the cross-module underscore imports stop being a contract.
- **Make `crewai/rag/embeddings/providers/*` honor their extras** by wrapping the `chromadb.utils.embedding_functions.*` imports in `try/except ImportError` like the LLM providers already do.
- **Re-enable the optional-dependencies test** in `lib/crewai-tools/tests/test_optional_dependencies.py:34` outside the GitHub Actions timeout window (e.g. nightly, or run it in a separate job with a longer timeout).
- **Replace `pydantic._internal._model_construction.ModelMetaclass`** with the public API in `crewai/agent/internal/meta.py:11` and `crewai/flow/runtime/__init__.py:56` before the Pydantic 3 release.

## Questions / Gaps

- Is there a documented module-boundary contract (e.g. a `ARCHITECTURE.md` or `CONTRACTS.md`) that names what each submodule can import? The deprecation shim (`lib/crewai/src/crewai/cli/__init__.py`) and the version pinning (`lib/crewai/pyproject.toml:11-12, 58`) suggest implicit contracts only.
- Is the `crewai-cli` CLI the intended replacement for `crewai.cli` long-term, or is it a transitional split? The shim emits a `DeprecationWarning` but `lib/cli/pyproject.toml:34-35` registers `crewai` as a console script that overlaps with `lib/crewai/pyproject.toml:146-147`.
- Why is `crewai-core` not exposed via the `crewai` console-script? Its `cli.py`-style usage is hidden inside `crewai-cli`, so users who need `version` checks must install `crewai-cli`.
- Is the lazy `__getattr__` for `Memory` (`lib/crewai/src/crewai/__init__.py:53-66`) intended as a model for the rest of the runtime, or is it just an optimization for one heavy import?
- Does the team have a CI check that enforces the lazy-import contract for `crewai.experimental` and `crewai.events`? `crewai/tests/test_imports.py` only covers 2 symbols.
- No clear evidence found for a top-level boundary-enforcement test that fails if `crewai-core` accidentally gains a `crewai` import. The grep I ran was manual; nothing in `tests/` enforces the leaf property.
- No clear evidence found for a documented policy on what constitutes "stable" vs "experimental" surface beyond the `crewai/experimental/__init__.py:1` docstring. Tools that want a stability signal must infer it.

---

Generated by `22.01-package-module-boundaries` against `crewai`.