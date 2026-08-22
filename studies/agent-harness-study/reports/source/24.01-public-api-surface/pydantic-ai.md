# Source Analysis: pydantic-ai

## 24.01 Public API Surface

### Source Info

| Field | Value |
|-------|-------|
| Name | pydantic-ai |
| Path | `studies/agent-harness-study/sources/pydantic-ai` |
| Language / Stack | Python 3.10+ (uv workspace; pydantic, httpx, anyio, opentelemetry; mkdocs/mkdocstrings docs) |
| Analyzed | 2026-08-22 |

## Summary

Pydantic AI exposes a deliberately layered public API across five published packages in a uv workspace: the umbrella `pydantic-ai` package (`pyproject.toml:13-51`), the core `pydantic-ai-slim` (`pydantic_ai_slim/pyproject.toml:13-16`), the `pydantic-graph` engine (`pydantic_graph/`), the `pydantic-evals` framework (`pydantic_evals/pydantic_evals/__init__.py:1-20`), and the `clai` CLI (`clai/pyproject.toml:13-16`). The primary entry point is the `Agent` class (`pydantic_ai_slim/pydantic_ai/agent/__init__.py:188-207`), backed by a curated top-level export list of ~180 symbols grouped by category (`pydantic_ai_slim/pydantic_ai/__init__.py:166-341`).

The public surface is distinguished from internals through four reinforcing mechanisms: (1) underscore-prefix naming for internal modules (`_agent_graph.py`, `_output.py`, `_run_context.py`, `_utils.py`, etc. in `pydantic_ai_slim/pydantic_ai/`); (2) explicit `__all__` in nearly every public module; (3) a hand-curated API reference that whitelists members per docs page instead of dumping modules (`docs/api/agent.md:3-16`, `docs/api/models/base.md:3-14`); and (4) a written version policy (`docs/version-policy.md:3-17`) backed by a systematic deprecation apparatus — a runtime-visible `PydanticAIDeprecationWarning` (`pydantic_ai_slim/pydantic_ai/_warnings.py:9-15`), 50+ `@deprecated(...)` annotations, and module-level `__getattr__` forwarding shims for renamed symbols (`pydantic_ai_slim/pydantic_ai/__init__.py:348-370`, `pydantic_graph/pydantic_graph/beta/__init__.py:78-114`).

Documentation is unusually rigorous: every Python code example in docs, docstrings, and all three packages is executed and linted in CI (`tests/test_examples.py:81-100,141-297`), so public API examples cannot silently rot. Extension authors get first-class abstract interfaces (`AbstractAgent`, `WrapperAgent`, `AbstractToolset` + wrapper toolsets, `Model`/`StreamedResponse`, `AbstractCapability`), and a lower-abstraction imperative API (`pydantic_ai.direct`, `pydantic_ai_slim/pydantic_ai/direct.py:27-33`) for users who want the model layer without the agent loop.

## Rating

**9/10** — Clear, documented, tested public model with explicit interfaces, stability policy, deprecation machinery, and operational safeguards (CI-executed examples, strict typing). The point is withheld for transitional rough edges: dual CLI entry points during a rename (`pai` vs `clai`), a large top-level export list that leans on comment grouping for navigability, and a few deprecated/semi-public names still reachable (`DeferredToolset`, `HistoryProcessor`, `CAPABILITY_TYPES`). A new integration can identify and use the stable API without reading implementation details.

## Evidence Collected

| Area | Evidence | File:Line |
|------|----------|-----------|
| Workspace packages | `pydantic-ai` umbrella, `pydantic-ai-slim`, `pydantic-graph`, `pydantic-evals`, `clai`, `examples` declared as uv workspace members | `pyproject.toml:83-90` |
| Umbrella package = slim + all provider extras | `pydantic-ai` depends on `pydantic-ai-slim[openai,...,spec]=={{ version }}` | `pyproject.toml:49-51` |
| Top-level export list, categorized | `__all__` with `# agent`, `# exceptions`, `# messages`, `# toolsets`, etc. section comments | `pydantic_ai_slim/pydantic_ai/__init__.py:166-341` |
| Primary entry point | `Agent` dataclass with docstring containing minimal runnable example | `pydantic_ai_slim/pydantic_ai/agent/__init__.py:188-207` |
| Run-method API (5 ways to run) | `run`, `run_sync`, `run_stream`, `run_stream_events`, `iter` documented with linked API refs | `docs/agent.md:66-73` |
| CLI entry points | `pai = "pydantic_ai._cli:cli_exit"` with TODO to remove; `clai = "clai:cli"` | `pyproject.toml:72-74`; `clai/pyproject.toml:57-58` |
| `python -m` CLI support | `python -m pydantic_ai` and `python -m clai` run the CLI | `pydantic_ai_slim/pydantic_ai/__main__.py:1-6`; `clai/clai/__main__.py:1-6` |
| CLI agent loading contract | `load_agent` accepts `module:variable` import strings and YAML/JSON `AgentSpec` files | `pydantic_ai_slim/pydantic_ai/_cli/__init__.py:94-120` |
| Declarative API | `Agent.from_spec` / `Agent.from_file` build agents from YAML/JSON/dict specs | `pydantic_ai_slim/pydantic_ai/agent/__init__.py:627-776,833-922` |
| Low-abstraction imperative API | `model_request`, `model_request_sync`, `model_request_stream`, `model_request_stream_sync` in `direct.py` `__all__` | `pydantic_ai_slim/pydantic_ai/direct.py:27-33` |
| Internal-module convention | Underscore-prefixed modules (`_agent_graph.py`, `_output.py`, `_run_context.py`, `_utils.py`, `_warnings.py`, …) | `pydantic_ai_slim/pydantic_ai/` (directory listing) |
| Curated API reference | `docs/api/agent.md` whitelists 12 members of `pydantic_ai.agent` rather than exposing the module wholesale | `docs/api/agent.md:3-16` |
| Curated model interface | `docs/api/models/base.md` exposes only `Model`, `StreamedResponse`, `KnownModelName`, request-gating helpers | `docs/api/models/base.md:3-14` |
| Version policy | No intentional breaking changes in V1 minors; deprecated functionality removed only in V2; enumerated non-breaking categories | `docs/version-policy.md:3-17` |
| Beta stability marker | Beta features "indicated by a `beta` module"; API may change without compat | `docs/version-policy.md:20-22` |
| Deprecation warning class | `PydanticAIDeprecationWarning(UserWarning)` so deprecations are visible by default | `pydantic_ai_slim/pydantic_ai/_warnings.py:9-15` |
| Deprecated-symbol shims | `__getattr__` remaps `BuiltinToolCallPart`→`NativeToolCallPart` etc. with warning | `pydantic_ai_slim/pydantic_ai/__init__.py:345-370` |
| Legacy namespace deprecation | `pydantic_graph.beta` forwards to permanent homes with `PydanticGraphDeprecationWarning` | `pydantic_graph/pydantic_graph/beta/__init__.py:1-11,78-114` |
| Deprecated legacy runner | Importing `Graph`/`GraphRun` from `pydantic_graph` top level warns; `BaseNode`/`End`/`GraphRunContext`/`Edge` survive to v2 | `pydantic_graph/pydantic_graph/__init__.py:5-20,63-94` |
| `@deprecated` usage breadth | 52 occurrences across models, providers, toolsets, messages, result, mcp | e.g. `pydantic_ai_slim/pydantic_ai/result.py:534,594,726`; `pydantic_ai_slim/pydantic_ai/mcp.py:1829` |
| Typed public API | `py.typed` markers in all packages; pyright strict across all source; mypy strict on typed-API tests | `pydantic_ai_slim/pydantic_ai/py.typed`; `pyproject.toml:243-261,281-283` |
| Tested doc examples | `find_examples('docs', 'pydantic_ai_slim', 'pydantic_graph', 'pydantic_evals')` — every example executed and ruff-linted with mocked models | `tests/test_examples.py:86,141-297` |
| Runnable example in docstring | `Agent` class docstring shows 4-line usage with expected output | `pydantic_ai_slim/pydantic_ai/agent/__init__.py:197-206` |
| Standalone examples package | `pydantic-ai-examples` (15 example modules) published as separate workspace package | `pyproject.toml:53-54`; `examples/pydantic_ai_examples/` |
| Toolset extension point | `AbstractToolset` + 12 wrapper toolsets (`Filtered`, `Prefixed`, `Prepared`, `Renamed`, `Combined`, `Wrapper`, …) exported | `pydantic_ai_slim/pydantic_ai/toolsets/__init__.py:19-40` |
| Capability extension point | `AbstractCapability`, wrap-handler types, `CAPABILITY_TYPES` serialization registry exported | `pydantic_ai_slim/pydantic_ai/capabilities/__init__.py:16-34,70-92` |
| Model extension point | `Model`, `StreamedResponse`, `KnownModelName`, `known_model_names()` documented as stable enumeration | `docs/api/models/base.md:3-14`; `pydantic_ai_slim/pydantic_ai/models/__init__.py:76-84` |
| Agent extension point | `AbstractAgent` ABC + `WrapperAgent` re-exported from `pydantic_ai.agent` | `pydantic_ai_slim/pydantic_ai/agent/__init__.py:91-101,111-131` |
| Optional-dependency API grouping | Per-provider extras (`openai`, `anthropic`, `google`, …) and integration extras (`mcp`, `temporal`, `evals`, `ui`, `spec`, …) | `pydantic_ai_slim/pyproject.toml:67-145` |
| Semi-public registry exported | `CAPABILITY_TYPES` in capabilities `__all__` (used for spec serialization) | `pydantic_ai_slim/pydantic_ai/capabilities/__init__.py:120` |
| Deprecated alias still exported | `DeferredToolset` imported with `reportDeprecated` and listed in toolsets `__all__` | `pydantic_ai_slim/pydantic_ai/toolsets/__init__.py:9,30` |
| Graph nodes in public surface | `CallToolsNode`, `ModelRequestNode`, `UserPromptNode` exported top-level for `iter()` users | `pydantic_ai_slim/pydantic_ai/__init__.py:9-15,174-176` |

## Answers to Dimension Questions

**1. What is the intended public API surface?**
Three tiers. (a) Application developers: the `Agent` class and its run methods (`pydantic_ai_slim/pydantic_ai/agent/__init__.py:188-207`, `docs/agent.md:66-73`), the ~180-symbol top-level namespace (`pydantic_ai_slim/pydantic_ai/__init__.py:166-341`), declarative `AgentSpec` loading (`pydantic_ai_slim/pydantic_ai/agent/__init__.py:627-776`), and the `clai`/`pai` CLI (`clai/pyproject.toml:57-58`, `pyproject.toml:72-74`). (b) Power users: the imperative `pydantic_ai.direct` model API (`pydantic_ai_slim/pydantic_ai/direct.py:27-33`) and the `agent.iter()` graph-node API with `ModelRequestNode`/`CallToolsNode` exported (`pydantic_ai_slim/pydantic_ai/__init__.py:174-176`). (c) Extension authors: `AbstractAgent`/`WrapperAgent`, `Model`/`StreamedResponse`, `AbstractToolset` + wrapper toolsets, and `AbstractCapability` (`pydantic_ai_slim/pydantic_ai/agent/__init__.py:111-131`; `docs/api/models/base.md:3-14`; `pydantic_ai_slim/pydantic_ai/toolsets/__init__.py:22-40`; `pydantic_ai_slim/pydantic_ai/capabilities/__init__.py:97-130`).

**2. Is the stable API easy to distinguish from internal implementation details?**
Yes, via layered conventions. Internal modules are underscore-prefixed (`_agent_graph.py`, `_output.py`, `_run_context.py` in `pydantic_ai_slim/pydantic_ai/`); public modules declare `__all__`; the API reference hand-picks documented members (`docs/api/agent.md:3-16` lists 12 of dozens of module members). The distinction is enforced culturally too: contributor guidance in `AGENTS.md` requires being "thoughtful and deliberate about new abstractions, public APIs" and states tests must go "through public APIs, not private methods (prefixed with `_`)". The one soft spot: `pydantic_ai.agent.__all__` re-exports `PydanticAIDeprecationWarning` and `ToolsPrepareFunc` (`pydantic_ai_slim/pydantic_ai/agent/__init__.py:129-130`) that are not in the top-level `__all__`, so the two export lists are not perfectly aligned.

**3. Does the API expose the right level of abstraction for agent harness users?**
Yes, with explicit opt-in to lower levels. The default path (`Agent(...)` + `run/run_sync/run_stream`) hides the graph engine entirely; the docstring example is 4 lines (`pydantic_ai_slim/pydantic_ai/agent/__init__.py:197-206`). Users needing loop control get `iter()` yielding graph nodes without losing the agent facade (`docs/agent.md:73`). Users needing no agent at all get `direct.model_request*` ("the only abstraction is input and output schema translation", `pydantic_ai_slim/pydantic_ai/direct.py:1-7`). Runtime internals (parts manager, OTel encoding, SSRF guards) stay private (`_parts_manager.py`, `_otel_messages.py`, `_ssrf.py`), while cross-cutting behavior is exposed through the composable toolset/capability wrappers rather than `Agent` kwargs (`pydantic_ai_slim/pydantic_ai/toolsets/__init__.py:22-40`; capabilities guidance in `pydantic_ai_slim/pydantic_ai/capabilities/AGENTS.md`).

**4. Are examples sufficient to use the API correctly without reading internals?**
Yes — this is the strongest aspect of the surface. Every code example in docs, docstrings, and the three packages is discovered by `find_examples('docs', 'pydantic_ai_slim', 'pydantic_graph', 'pydantic_evals')` and executed against a mocked `infer_model` with ruff linting in CI (`tests/test_examples.py:86,152-153,280-297`), so examples are executable truth, not prose. Docs follow a context→code→caveats structure with expected output inline (`docs/agent.md:24-60`), and there are 39 top-level doc pages plus a full API reference (`mkdocs.yml:20-222`). A separate published examples package (`pydantic-ai-examples`, `pyproject.toml:53-54`) provides larger end-to-end apps.

## Architectural Decisions

1. **Slim core + extras packaging.** `pydantic-ai-slim` has minimal required deps (`pydantic_ai_slim/pyproject.toml:55-65`) and ~30 optional groups split per provider and per integration (`pydantic_ai_slim/pyproject.toml:67-145`); the umbrella `pydantic-ai` package exists only to pin slim with all extras (`pyproject.toml:49-51`). This keeps the public import surface installable without pulling every provider SDK.
2. **Multi-package split with versioned contracts.** Graph engine (`pydantic-graph`) and evals (`pydantic-evals`) are independently published packages with their own `__init__` contracts and deprecation warnings (`pydantic_graph/pydantic_graph/__init__.py:97-137`; `pydantic_evals/pydantic_evals/__init__.py:13-20`), so the agent harness depends on them as libraries, not folders.
3. **Curated-docs-as-API-contract.** Rather than auto-documenting whole modules, each `docs/api/*.md` whitelists members (`docs/api/agent.md:3-16`), making the documented API an explicit, reviewable artifact.
4. **Deprecation via forwarding shims, not removal.** Renamed symbols keep working through module `__getattr__` with warnings (`pydantic_ai_slim/pydantic_ai/__init__.py:355-370`; `pydantic_graph/pydantic_graph/beta/__init__.py:102-114`), and deprecated classes remain exported with `@deprecated` markers (`pydantic_ai_slim/pydantic_ai/toolsets/__init__.py:9`).
5. **Runtime-visible deprecations.** `PydanticAIDeprecationWarning` deliberately inherits `UserWarning`, not `DeprecationWarning`, so users actually see deprecation notices (`pydantic_ai_slim/pydantic_ai/_warnings.py:9-15`).

## Notable Patterns

- **Category-commented `__all__`**: the top-level export list is organized with `# agent`, `# messages`, `# toolsets` section comments, doubling as an API map (`pydantic_ai_slim/pydantic_ai/__init__.py:166-341`).
- **Stable enumeration helpers over type introspection**: `known_model_names()` is provided and documented as the public way to list models, explicitly discouraging `get_args()` on the `KnownModelName` alias (`pydantic_ai_slim/pydantic_ai/models/__init__.py:76-84`).
- **Wrapper-based extension axes**: cross-cutting behavior (filter, prefix, rename, prepare, approve, defer) is delivered by wrapping `AbstractToolset`/`AbstractCapability` instead of growing `Agent` kwargs — the toolsets and capabilities AGENTS.md files codify this rule.
- **Examples as tests**: `tests/test_examples.py:141-297` runs every docs/docstring example with mocked model inference, env seeding, and HTTP mocking; `pyproject.toml:136` (`pytest-examples`) and `mkdocs.yml:341-372` (mkdocstrings) tie docs generation directly to source docstrings.
- **Graceful CLI degradation**: the CLI extra raises a targeted ImportError telling users which extra group to install (`pydantic_ai_slim/pydantic_ai/_cli/__init__.py:39-43`).

## Tradeoffs

- **Large flat top-level namespace.** ~180 exported names at `pydantic_ai.*` maximize discoverability (`from pydantic_ai import Agent, Tool, ModelResponse` works for everything) but require the comment-grouping and docs to navigate; submodule imports (`pydantic_ai.toolsets.PreparedToolset`) remain the precise path.
- **Deprecated surface retained for a full major cycle.** Keeping renamed/deprecated names importable through v1 (`pydantic_ai_slim/pydantic_ai/toolsets/__init__.py:9,30`; `_a2a.py` shim slated "removed in v2" per `pyproject.toml:362`) preserves upgrade confidence at the cost of a larger apparent surface until v2.
- **`KnownModelName` literal grows unboundedly** with each provider/model addition; mitigated by accepting plain `str` everywhere ("the actual list of allowed models changes frequently", `pydantic_ai_slim/pydantic_ai/agent/__init__.py:349`) rather than making the alias load-bearing.
- **Graph node types in the public surface.** Exporting `ModelRequestNode`/`CallToolsNode`/`UserPromptNode` (`pydantic_ai_slim/pydantic_ai/__init__.py:174-176`) is required for `iter()` users but couples the top-level API to the graph execution model.

## Failure Modes / Edge Cases

- **Two CLI names during migration**: `pai` and `clai` both exist, with a TODO to remove `pai` (`pyproject.toml:73`, `pydantic_ai_slim/pyproject.toml:151`); users may install either and get the same code path (`clai/clai/__init__.py:9-11` delegates to `pydantic_ai._cli.cli_exit`).
- **Misaligned export lists**: `pydantic_ai.agent.__all__` includes `PydanticAIDeprecationWarning` and `ToolsPrepareFunc` (`pydantic_ai_slim/pydantic_ai/agent/__init__.py:129-130`) absent from the top-level list, so `from pydantic_ai import ToolsPrepareFunc` fails while `from pydantic_ai.agent import ...` works — a discoverability trap, not a correctness one.
- **Semi-public registries**: `CAPABILITY_TYPES` (`pydantic_ai_slim/pydantic_ai/capabilities/__init__.py:70-92`) is exported and mutable-by-design for spec serialization; external mutation could affect `AgentSpec` capability resolution.
- **Deprecation warnings visible by default** can surface in user logs for code that still imports legacy names (`pydantic_ai_slim/pydantic_ai/_warnings.py:9-15`) — intentional, but noisy for slow migrators.
- **Typo resilience**: module `__getattr__` fallbacks raise a clear `AttributeError` for unknown names (`pydantic_ai_slim/pydantic_ai/__init__.py:370`; `pydantic_graph/pydantic_graph/__init__.py:94`) — no silent attribute absorption.

## Future Considerations

- Complete the `pai` → `clai` CLI rename and drop the legacy entry point (`pyproject.toml:73`).
- Reconcile the `pydantic_ai.agent` vs top-level `__all__` drift (e.g., decide whether `ToolsPrepareFunc`, `PydanticAIDeprecationWarning` belong top-level).
- Execute the documented v2 removals: `pydantic_graph.beta` namespace (`pydantic_graph/pydantic_graph/beta/__init__.py:10`), legacy `BaseNode` runner (`pydantic_graph/pydantic_graph/__init__.py:13-19`), `_a2a.py` (`pyproject.toml:362`), and deprecated aliases (`DeferredToolset`, `HistoryProcessor`).
- Consider a `pydantic_ai.__all__` lint/CI check that asserts every `__all__` member resolves and docs reference pages stay in sync with export lists.

## Questions / Gaps

- No formal semver/`__version__`-gated feature flags were found beyond the version policy and deprecation warnings; stability is enforced socially (policy doc + review) and via CI, not mechanically. Evidence searched: `docs/version-policy.md`, `pydantic_ai_slim/pydantic_ai/_warnings.py`, `pyproject.toml` (no `api-stability` tooling configured).
- Whether the curated `docs/api/*.md` member lists are verified against `__all__` automatically: no such check was found in `tests/test_examples.py` or the mkdocs hooks (`docs/.hooks/` not audited in depth within this study's boundary); the docs build runs with `strict: true` (`mkdocs.yml:3`), which catches broken links but not missing members.
- HTTP/RPC surface exposure (A2A, AG-UI, web UI) was noted as package extras (`pydantic_ai_slim/pyproject.toml:129-135`) but endpoint-level API stability was not audited; that belongs to a service-endpoint dimension.

---

Generated by `24.01-public-api-surface` against `pydantic-ai`.
