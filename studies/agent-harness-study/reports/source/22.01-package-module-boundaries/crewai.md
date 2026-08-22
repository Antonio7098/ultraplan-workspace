# Source Analysis: crewai

## 22.01 Package and Module Boundaries

### Source Info

| Field | Value |
|-------|-------|
| Name | crewai |
| Path | `studies/agent-harness-study/sources/crewai` |
| Language / Stack | Python 3.10–3.13, uv workspace monorepo, hatchling builds, Pydantic v2 |
| Analyzed | 2026-08-21 |

## Summary

crewAI is organized as a six-package uv workspace (`pyproject.toml:227-235`): `crewai` (framework meta-package), `crewai-core` (leaf utilities), `crewai-cli` (CLI), `crewai-tools` (tool catalog), `crewai-files` (multimodal file handling), and `crewai-devtools` (private release tooling). The intended dependency direction is `core → (nothing)`, `cli → core`, `crewai → core + cli`, `crewai-tools → crewai`, with `crewai-files` fully standalone.

In practice the layering is real but asymmetric and enforced only by convention. `crewai-core` and `crewai-files` are genuinely independent bottom layers (verified by import scans: zero workspace imports). The top of the graph, however, is entangled: `crewai-tools` hard-depends on the entire `crewai` distribution just to subclass `BaseTool` (`lib/crewai-tools/pyproject.toml:13`), while `crewai` reaches back into `crewai_tools` at five call sites via lazy function-level imports (`lib/crewai/src/crewai/task.py:1182`, `lib/crewai/src/crewai/agent/core.py:1179`, `lib/crewai/src/crewai/project/crew_base.py:324`, `lib/crewai/src/crewai/mcp/tool_resolver.py:194`, `lib/crewai/src/crewai/project/json_loader.py:1888`). `crewai` even declares `crewai-tools` as an optional extra pointing back at its dependent (`lib/crewai/pyproject.toml:57-58`), a package-level cycle broken only by optionality. Inside `crewai` itself, circularity is handled by a hand-rolled namespace-injection mechanism in `__init__.py` (`lib/crewai/src/crewai/__init__.py:128-145`) plus 104 files of `TYPE_CHECKING` guards and lazy `__getattr__` loaders — workable but evidence that module boundaries are managed manually rather than by clean layering.

Public vs internal API distinction is weak-to-moderate: `crewai` and `crewai-files` curate `__all__` exports; `crewai-core` exports nothing at its package root (`lib/crewai-core/src/crewai_core/__init__.py:1` is a bare version string), forcing consumers to reach into submodules. There is no import-linter, no layering test, and no `_internal/` convention beyond one directory (`lib/crewai/src/crewai/agent/internal/`). CI enforces strict mypy and ruff with all relative imports banned (`pyproject.toml:90-91`, `pyproject.toml:120-131`, `.github/workflows/type-checker.yml:61`), which guarantees import hygiene but not architectural direction. All six packages are version-pinned to each other with exact `==` pins (`1.14.8a2`), so they ship in lockstep and are not independently installable in practice.

## Rating

**5 / 10** — Present but inconsistent. The workspace split, the pure `crewai-core`/`crewai-files` bottom layers, lazy-loading discipline, and the deprecation shim show deliberate boundary engineering. But the `crewai ↔ crewai-tools` cycle, the undeclared `crewai-cli → crewai` runtime import (`lib/cli/src/crewai_cli/deploy/validate.py:41`), the duplicate console entrypoint (`lib/crewai/pyproject.toml:146-147` duplicating `lib/cli/pyproject.toml`), the absence of any automated layering enforcement, and the lockstep version pins keep this from the "clear model with tests and explicit interfaces" band (7-8). The dimension's litmus question — *can you use the tool system without pulling in the entire runtime?* — answers **no at the distribution level**: `crewai.tools.base_tool` is a leaf module (stdlib + pydantic only, `lib/crewai/src/crewai/tools/base_tool.py:1-15`), but it lives inside the `crewai` wheel, so installing `crewai-tools` transitively installs chromadb, lancedb, openai, and the whole framework (`lib/crewai/pyproject.toml:14-56`).

## Evidence Collected

Every entry includes a file path with line numbers.

| Area | Evidence | File:Line |
|------|----------|-----------|
| Package structure | uv workspace with 6 members: crewai, crewai-tools, devtools, crewai-files, cli, crewai-core | `pyproject.toml:227-235` |
| Package structure | Workspace source mapping for all 6 packages | `pyproject.toml:238-244` |
| Dependency direction | `crewai` depends on `crewai-core==1.14.8a2` and `crewai-cli==1.14.8a2` | `lib/crewai/pyproject.toml:11-12` |
| Dependency direction | `crewai-tools` depends on `crewai==1.14.8a2` (full framework, exact pin) | `lib/crewai-tools/pyproject.toml:13` |
| Dependency direction | `crewai-cli` depends on `crewai-core` but NOT on `crewai` | `lib/cli/pyproject.toml:11-24` |
| Dependency direction | `crewai-core` dependencies are all third-party (appdirs, cryptography, httpx, pydantic, rich, opentelemetry, etc.) — no workspace deps | `lib/crewai-core/pyproject.toml:10-22` |
| Dependency direction | `crewai-files` dependencies all third-party (Pillow, pypdf, aiocache, aiofiles, av) — no workspace deps | `lib/crewai-files/pyproject.toml:10-18` |
| Circular dependency (package level) | `crewai` optional extra `tools = ["crewai-tools==1.14.8a2"]` points back at its own dependent | `lib/crewai/pyproject.toml:57-58` |
| Circular dependency (package level) | `crewai → crewai_tools` lazy imports at 5 sites; `crewai_tools → crewai` in 88 source files | `lib/crewai/src/crewai/task.py:1182`, `lib/crewai-tools/src/crewai_tools/__init__.py:1-60` |
| Undeclared dependency | `crewai_cli.deploy.validate` imports `crewai.project.json_loader` at module level; `crewai` absent from `crewai-cli` deps | `lib/cli/src/crewai_cli/deploy/validate.py:41` |
| Duplicate entrypoint | `crewai = "crewai_cli.cli:crewai"` script defined in both `crewai` and `crewai-cli` packages | `lib/crewai/pyproject.toml:146-147`, `lib/cli/pyproject.toml` `[project.scripts]` |
| Independence: crewai-core | Only stdlib + third-party imports across all of `crewai_core` (verified by import scan); zero `crewai` references | `lib/crewai-core/src/crewai_core/` (whole tree) |
| Independence: crewai-files | Only stdlib + third-party imports; explicit re-exports of 60+ public symbols in `__init__` | `lib/crewai-files/src/crewai_files/__init__.py:3-80` |
| Independence: event bus | `event_bus.py` runtime imports limited to `crewai.events.*` leaves; `RuntimeState` only under TYPE_CHECKING | `lib/crewai/src/crewai/events/event_bus.py:26-27` |
| Intra-package circularity handling | Namespace injection: `_resolve_namespace` written into 8 modules' `__dict__`, then 9 `model_rebuild()` calls | `lib/crewai/src/crewai/__init__.py:128-159` |
| Lazy loading | `Memory` lazily imported via module `__getattr__` to avoid lancedb at init | `lib/crewai/src/crewai/__init__.py:53-66` |
| Lazy loading | Event types lazy-loaded via `_LAZY_EVENT_MAPPING` to avoid ~12 Pydantic modules at init | `lib/crewai/src/crewai/events/__init__.py:9-11,154,257-260` |
| TYPE_CHECKING discipline | 104 files under `lib/crewai/src` use `if TYPE_CHECKING:` guards for upward imports | e.g. `lib/crewai/src/crewai/utilities/tool_utils.py:24-27`, `lib/crewai/src/crewai/llms/base_llm.py` |
| Public API surface | `crewai.__all__` curates 17 public symbols | `lib/crewai/src/crewai/__init__.py:187-205` |
| Public API surface | `crewai.tools.__all__` exports exactly BaseTool, EnvVar, tool | `lib/crewai/src/crewai/tools/__init__.py:1-8` |
| Public API surface | `crewai-core` package root exports only `__version__` — consumers import submodules directly | `lib/crewai-core/src/crewai_core/__init__.py:1` |
| Public API surface | All 5 distributable packages ship `py.typed` markers | `lib/*/src/*/py.typed` |
| Internal API convention | `crewai/agent/internal/` package for extension metaclass, imported only by `base_agent.py` | `lib/crewai/src/crewai/agent/internal/meta.py:1-13`, `lib/crewai/src/crewai/agents/agent_builder/base_agent.py:26` |
| Deprecation shim | `crewai.cli` meta-path finder remaps legacy imports to `crewai_cli` with DeprecationWarning | `lib/crewai/src/crewai/cli/__init__.py:1-66` |
| Optional-dep guard (good) | Platform tools import wrapped in try/except, returns `[]` on failure | `lib/crewai/src/crewai/agent/core.py:1177-1186` |
| Optional-dep guard (good) | JSON loader resolves tool classes from `crewai_tools` with ImportError fallback returning None | `lib/crewai/src/crewai/project/json_loader.py:1887-1893` |
| Optional-dep gap | `MCPServerAdapter` import unguarded — raises ImportError if `crewai-tools` extra not installed but `mcp_server_params` set | `lib/crewai/src/crewai/project/crew_base.py:324` |
| Import hygiene enforcement | ruff `ban-relative-imports = "all"` | `pyproject.toml:90-91` |
| Import hygiene enforcement | mypy strict mode across `lib/` in CI | `pyproject.toml:120-131`, `.github/workflows/type-checker.yml:61` |
| Import hygiene enforcement | ruff check + format in CI | `.github/workflows/linter.yml:54-57` |
| No layering enforcement | No import-linter config, no dependency-cruiser, no layering test found (searched pyproject, pre-commit, workflows, tests for "import.linter", "circular", "layered") | — |
| Lockstep versioning | All cross-package pins are exact `==1.14.8a2`; versions duplicated in source `__init__` files | `lib/crewai/pyproject.toml:11-12`, `lib/crewai-tools/pyproject.toml:13`, `lib/crewai/src/crewai/__init__.py:51` |
| Separation tests | `crewai-core` has exactly 1 smoke test file covering leaf modules | `lib/crewai-core/tests/test_smoke.py:1-40` |
| Separation tests | `crewai-tools` import test asserts no Pydantic deprecation warnings at import | `lib/crewai-tools/tests/tools/test_import_without_warnings.py:1-9` |
| Separation tests | Public-API importability test (TaskOutput, CrewOutput from `crewai`) — minimal, 15 lines | `lib/crewai/tests/test_imports.py:1-15` |
| Runtime circular-dependency guard | `CircularDependencyError` exists for event-handler dependency graph (not module imports) | `lib/crewai/src/crewai/events/handler_graph.py:15` |
| Test isolation | Root `conftest.py` shared by all packages; per-package test dirs (207/26/8/32/2/1 test files) | `conftest.py:1-10`, `pyproject.toml:142-148` |

## Answers to Dimension Questions

**1. Are modules cleanly separated?**
Partially. The workspace gives six physically separate distributions with their own pyproject, tests, and `py.typed`. The bottom layers are clean: `crewai-core` (verified zero workspace imports across its 27 modules, e.g. `lib/crewai-core/src/crewai_core/telemetry.py` uses only opentelemetry/stdlib) and `crewai-files` (zero `crewai` references; explicit 80-line export surface at `lib/crewai-files/src/crewai_files/__init__.py`). The top is not clean: `crewai` ↔ `crewai-tools` form a package-level cycle broken only by lazy imports and an optional extra, and `crewai-cli` has an undeclared runtime import of `crewai` (`lib/cli/src/crewai_cli/deploy/validate.py:41`). Within `crewai`, `utilities/` and `events/listeners/` import upward into `task`/`crew` at runtime (`lib/crewai/src/crewai/utilities/task_output_storage_handler.py:12`, `lib/crewai/src/crewai/utilities/planning_handler.py:9`), so "utilities" is not a lower layer despite its name.

**2. Do dependencies flow in one direction?**
Core does; the periphery does not. Verified one-directional edges: `crewai-files → (nothing)`, `crewai-core → (nothing)`, `cli → crewai-core`, `crewai → crewai-core`, `crewai → crewai-cli`. Non-one-directional edges: `crewai-tools → crewai` (declared, `lib/crewai-tools/pyproject.toml:13`) combined with five reverse lazy imports from `crewai` into `crewai_tools` (`lib/crewai/src/crewai/agent/core.py:1179`, `lib/crewai/src/crewai/task.py:1182`, `lib/crewai/src/crewai/project/crew_base.py:324`, `lib/crewai/src/crewai/mcp/tool_resolver.py:194`, `lib/crewai/src/crewai/project/json_loader.py:1888`); plus the undeclared `crewai-cli → crewai` import. The `crewai` optional extra `tools` (`lib/crewai/pyproject.toml:57-58`) formally documents the cycle.

**3. Can modules be used independently?**
`crewai-core`: yes — it is a genuine leaf library (paths, printer, telemetry, auth tokens, user data) with its own smoke test (`lib/crewai-core/tests/test_smoke.py:6-13`). `crewai-files`: yes — fully standalone multimodal file handling with provider formatters and uploaders (`lib/crewai-files/src/crewai_files/`). `crewai-cli`: mostly — it depends only on `crewai-core` per its manifest, but the deploy-validate path breaks without `crewai` installed (`lib/cli/src/crewai_cli/deploy/validate.py:41`). `crewai-tools`: no — it requires the full `crewai` distribution (`lib/crewai-tools/pyproject.toml:13`), which transitively pulls chromadb, lancedb, openai, instructor, and opentelemetry (`lib/crewai/pyproject.toml:14-56`). At module level the tool abstraction is decoupled — `crewai/tools/base_tool.py:1-15` imports only stdlib and pydantic — but it is packaged inside `crewai`, so the answer to the dimension's litmus question is no at install granularity.

**4. Are public APIs distinguished from internal ones?**
Moderately. `crewai` curates `__all__` with 17 symbols and uses lazy `__getattr__` for heavy imports (`lib/crewai/src/crewai/__init__.py:53-66,187-205`); `crewai.tools` exports exactly three names (`lib/crewai/src/crewai/tools/__init__.py:1-8`); `crewai-files` re-exports its full public surface explicitly (`lib/crewai-files/src/crewai_files/__init__.py:3-80`). But `crewai-core`'s root namespace is empty (`lib/crewai-core/src/crewai_core/__init__.py:1`), so every consumer — including `crewai` and `crewai-cli` — reaches into submodules (`lib/cli/src/crewai_cli/config.py` imports `crewai_core.settings`), making every core module de-facto public. There is a single `internal/` package (`lib/crewai/src/crewai/agent/internal/`) but no repo-wide `_internal` or underscore convention; conversely, deeply nested modules like `crewai.utilities.prompts` are injected into the public namespace by the init-time machinery (`lib/crewai/src/crewai/__init__.py:78-81`). No stability annotations (`@public`, `__all__` in most submodules) or API-surface tests beyond the 15-line `lib/crewai/tests/test_imports.py`.

## Architectural Decisions

1. **Workspace monorepo with lockstep exact pins.** All six packages pin each other at `==1.14.8a2` (`lib/crewai/pyproject.toml:11-12`, `lib/crewai-tools/pyproject.toml:13`). Packages are structurally separate but released as one unit — the boundary buys code organization, not independent versioning or installability.
2. **Extract shared leaf utilities into `crewai-core`.** Version, paths, settings, telemetry, printer, token management, and plus-API client live in a dependency-free package consumed by both `crewai` and `crewai-cli` (53 import sites in `crewai`, e.g. `lib/crewai/src/crewai/constants.py`; 5+ in cli, e.g. `lib/cli/src/crewai_cli/plus_api.py`). This deduplicates auth/telemetry that previously existed in both trees.
3. **Keep `BaseTool` in `crewai`, not `crewai-tools`.** The tool contract (`lib/crewai/src/crewai/tools/base_tool.py:103`) is a leaf module inside the framework package; `crewai-tools` is a pure catalog of concrete tools. This avoids a tools-core package but forces tool authors to install the entire runtime.
4. **Manage intra-package cycles at runtime, not by design.** The init-time namespace injection (`lib/crewai/src/crewai/__init__.py:128-159`), 104 TYPE_CHECKING-guarded files, and lazy `__getattr__` loaders (`lib/crewai/src/crewai/__init__.py:58-66`, `lib/crewai/src/crewai/events/__init__.py:257`) let a highly connected object graph (Agent ↔ Task ↔ Crew ↔ Flow) exist without import errors, at the cost of import-order sensitivity.
5. **Compatibility shim instead of a breaking change.** `crewai.cli` was extracted to `crewai-cli`, and a meta-path finder transparently remaps legacy imports (`lib/crewai/src/crewai/cli/__init__.py:38-66`) — a deliberate boundary-migration technique.
6. **Optional extras for heavy integrations.** Provider and integration deps (anthropic, bedrock, qdrant, docling, a2a, tools) are extras (`lib/crewai/pyproject.toml:53-110`), keeping the core install smaller while code references them defensively.

## Notable Patterns

- **Lazy `__getattr__` module pattern** for heavy optional deps: `Memory → lancedb` (`lib/crewai/src/crewai/__init__.py:53-66`) and ~40 event types (`lib/crewai/src/crewai/events/__init__.py:154,257-260`), with documented rationale ("avoid importing ~12 Pydantic model modules", `lib/crewai/src/crewai/events/__init__.py:9-11`).
- **Pydantic model_rebuild with explicit types namespace** to close forward-reference cycles across Agent/Task/Crew/Flow (`lib/crewai/src/crewai/__init__.py:152-176`).
- **Meta-path import shim** for cross-package renames with deprecation warnings (`lib/crewai/src/crewai/cli/__init__.py:38-66`).
- **Function-level imports as cycle breakers** for optional sibling packages, sometimes guarded (`lib/crewai/src/crewai/agent/core.py:1178-1186`), sometimes not (`lib/crewai/src/crewai/project/crew_base.py:324`).
- **Extension metaclass in an `internal/` package** (`lib/crewai/src/crewai/agent/internal/meta.py:13`) consumed by the public `BaseAgent` — the only explicit internal-API carve-out.
- **Runtime dependency-graph validation for event handlers** — `CircularDependencyError` in the handler graph (`lib/crewai/src/crewai/events/handler_graph.py:15`) — the circularity concern is enforced for handler wiring even though module-level circularity is not.

## Tradeoffs

- **Tool-catalog convenience vs. boundary integrity.** Putting `BaseTool` in `crewai` makes the tool API co-versioned with the runtime, but every tool author and `crewai-tools` consumer installs chromadb, lancedb, openai, and telemetry deps (`lib/crewai/pyproject.toml:14-56`). A `crewai-tools-core` split would answer the dimension's litmus question affirmatively.
- **Lazy imports vs. failure visibility.** Guarded lazy imports degrade silently (platform tools return `[]` on any Exception, `lib/crewai/src/crewai/agent/core.py:1184-1186`; JSON loader returns None, `lib/crewai/src/crewai/project/json_loader.py:1892-1893`), while unguarded ones fail late at feature use (`lib/crewai/src/crewai/project/crew_base.py:324`). Neither failure mode is surfaced at configuration time.
- **Namespace injection vs. explicit layering.** The `__init__.py` machinery keeps the Agent/Task/Crew graph importable without refactoring into layers, but it makes import order semantically significant and wraps the entire setup in one try/except that downgrades failures to a log warning (`lib/crewai/src/crewai/__init__.py:178-185`).
- **Lockstep pins vs. independent evolution.** Exact `==` pins eliminate version-skew bugs across packages but mean `crewai-files` — which needs nothing from the workspace — still cannot be consumed at a different version by third parties mixing releases.
- **Empty `crewai-core` root namespace vs. API stability.** Exporting nothing avoids heavy imports, but with no `__all__`, every internal module is implicitly public API for two consumer packages, constraining refactoring.

## Failure Modes / Edge Cases

- **Undeclared dependency breaks standalone `crewai-cli`.** `pip install crewai-cli` alone (satisfying its declared deps, `lib/cli/pyproject.toml:10-24`) then running deploy validation raises `ModuleNotFoundError` at `lib/cli/src/crewai_cli/deploy/validate.py:41` because `crewai` is not declared. Masked in practice because `crewai` installs `crewai-cli`, not vice versa.
- **MCP usage without the tools extra crashes at call time.** `mcp_server_params` set but `crewai-tools` not installed → unguarded `from crewai_tools import MCPServerAdapter` (`lib/crewai/src/crewai/project/crew_base.py:324`) raises ImportError deep in crew execution rather than at config validation.
- **Import-order sensitivity from namespace injection.** The rebuild sequence in `lib/crewai/src/crewai/__init__.py:128-176` assumes specific modules are already in `sys.modules` (`__init__.py:130`); if the wrapped block fails, `RuntimeState` is set to `None` (`__init__.py:185`) and downstream attribute errors surface far from the cause.
- **Broad exception swallowing hides missing tools.** `get_platform_tools` catches bare `Exception` — including genuine bugs in tool construction — and returns an empty list (`lib/crewai/src/crewai/agent/core.py:1184-1186`), conflating "extra not installed" with "tool broken."
- **Silent pydantic-warning suppression at import.** `crewai/__init__.py:28-49` globally monkey-patches `warnings.warn` for pydantic deprecations, trading user visibility for clean imports; a test pins the same expectation for `crewai_tools` (`lib/crewai-tools/tests/tools/test_import_without_warnings.py:4-9`).
- **Duplicate console entrypoint.** Both `crewai` and `crewai-cli` wheels register `crewai = "crewai_cli.cli:crewai"` (`lib/crewai/pyproject.toml:146-147`); installing both is the normal case, so the conflict is latent rather than active, but it makes the `crewai` package's CLI ownership ambiguous.
- **No automated detection of new boundary violations.** Nothing in CI (`.github/workflows/linter.yml:54-57`, `.github/workflows/type-checker.yml:61`) or pre-commit (`.pre-commit-config.yaml`) would catch a new upward import from `crewai-core`, a new unguarded `crewai_tools` import, or `crewai-cli` growing more undeclared `crewai` usage.

## Future Considerations

- Add import-linter (or equivalent) contracts to CI encoding the observed intended graph: `crewai-core` ← nothing; `crewai-files` ← nothing; `cli → {core}`; `crewai → {core, cli}`; `crewai-tools → {crewai}`. This is cheap and would have caught `lib/cli/src/crewai_cli/deploy/validate.py:41`.
- Extract `crewai.tools.base_tool` + `structured_tool` into a small `crewai-tools-core` (or move into `crewai-core`) so the tool system is installable without the runtime — directly resolving the litmus question.
- Declare the real dependency: either add `crewai` to `crewai-cli` deps (accepting the cycle) or make the deploy-validate import lazy with a clear error message.
- Give `crewai-core` an explicit root `__all__` re-exporting its stable surface (paths, printer, settings, version, token_manager) so consumers stop importing private paths.
- Replace the unguarded `MCPServerAdapter` import with a guarded one plus a config-time validation error naming the missing extra.
- Remove the duplicate `[project.scripts]` entry from `lib/crewai/pyproject.toml:146-147`, leaving CLI ownership in `crewai-cli`.
- Consider independent versioning (compatible-range pins like `crewai-core>=1.14,<2`) for the leaf packages that already have no upward coupling.

## Questions / Gaps

- **No evidence found** of any automated circular-import detection or layering test for Python module imports. Searched: all `pyproject.toml` files, `.pre-commit-config.yaml`, `.github/workflows/*.yml`, and test trees for "import-linter", "circular", "layered", "boundary"; the only `CircularDependencyError` is for event-handler wiring (`lib/crewai/src/crewai/events/handler_graph.py:15`).
- **No evidence found** of a documented, intended dependency diagram (no ARCHITECTURE.md; `lib/crewai-core/README.md` and `lib/crewai-files/README.md` are one-paragraph descriptions). The dependency graph above is reconstructed from manifests and import scans, not from a stated design document.
- **No evidence found** of per-package publish-time independence: whether `crewai-files` or `crewai-core` are published to PyPI separately, or only as transitive deps of `crewai`. The lockstep pins suggest the latter, but publish configuration (`.github/workflows/publish.yml`) was not analyzed in depth for per-package triggers.
- **No evidence found** of stability/deprecation policy for the de-facto-public `crewai-core` submodule paths consumed by `crewai` and `crewai-cli`.
- The runtime behavior of the unguarded `crew_base.py:324` import path (exact exception surfaced to users when the tools extra is missing) was verified by code reading only, not executed.

---

Generated by `dimension 22.01: Package and Module Boundaries` against `crewai`.
