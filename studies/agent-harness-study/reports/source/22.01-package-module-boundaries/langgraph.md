# Source Analysis: langgraph

## Dimension 22.01: Package and Module Boundaries

### Source Info

| Field | Value |
|-------|-------|
| Name | langgraph |
| Path | `studies/agent-harness-study/sources/langgraph` |
| Language / Stack | Python (hatchling + uv monorepo); JS/TS SDK moved out to a separate repository |
| Analyzed | 2026-08-21 |

## Summary

LangGraph is a multi-package Python monorepo: eight libraries live under `libs/` (checkpoint, checkpoint-postgres, checkpoint-sqlite, checkpoint-conformance, cli, langgraph core, prebuilt, sdk-py), each an independently versioned and independently published distribution (`Makefile:1-58`, `.github/workflows/release.yml`). All Python packages install into a single shared `langgraph` implicit namespace (PEP 420 — no top-level `langgraph/__init__.py` exists in any lib), so the runtime, persistence, tools, SDK, and CLI are physically separate distributions that compose at import time.

The intended layering is documented in `AGENTS.md:33-55`: `langgraph-checkpoint` sits at the bottom; `langgraph-prebuilt` and the core `langgraph` package depend on it; core additionally depends on `langgraph-sdk`; the CLI depends on the SDK; sdk-js is standalone (now only a pointer stub, `libs/sdk-js/README.md:7`). For five of the packages this layering holds and is genuinely enforced by per-library CI environments. The major exception is the prebuilt/core relationship: **the declared dependency direction is inverted relative to actual code flow**. Core declares `langgraph-prebuilt` as a production dependency (`libs/langgraph/pyproject.toml:30`) but never imports it in its own code, while prebuilt imports core modules extensively — including private ones — without declaring core as a dependency at all (`libs/prebuilt/pyproject.toml:26-29`). This creates an undeclared circular dependency between distributions that is masked by the fact that installing `langgraph` always pulls in `langgraph-prebuilt`.

Public vs internal API separation is handled through explicit mechanisms: an underscore-prefixed `_internal` subpackage with a stability disclaimer (`libs/langgraph/langgraph/_internal/__init__.py:1-4`), `__all__` lists on every public module, versioned deprecation warning classes with machine-readable `since`/`expected_removal` fields (`libs/langgraph/langgraph/warnings.py:13-48`), and `@beta` markers for experimental surfaces such as the v3 streaming protocol (`libs/langgraph/langgraph/stream/run_stream.py:30`). A dedicated conformance library programmatically validates that any third-party checkpointer satisfies the checkpoint storage contract (`libs/checkpoint-conformance/langgraph/checkpoint/conformance/validate.py:45`), which is unusually strong boundary testing for the persistence interface.

## Rating

**6 / 10** — Present but inconsistent; strong in most places, structurally unsound at one seam.

Rationale against the rubric:

- The bottom of the stack (checkpoint/store base, sdk-py, conformance) is cleanly separated, dependency-light, and provably independent (verified imports, CI-tested in isolation). That alone would score 8.
- The prebuilt↔core pair fails two of the dimension's core tests simultaneously: dependencies do not flow in one direction (undeclared cycle), and prebuilt cannot be used independently despite being packaged as if it could be. Prebuilt also reaches into core's private `_internal` and `pregel._tools` namespaces.
- There is no automated guard (import-linter contract, ruff banned-api entry for cross-package imports) that would catch these violations; enforcement relies on human discipline plus the AGENTS.md diagram.

The model is clear and mostly operationalized, but the flagship "tools/prebuilt" layer violates its own declared boundaries — squarely a 6 ("present but inconsistent ... fragile") rather than 7–8 ("explicit interfaces and operational safeguards").

## Evidence Collected

| Area | Evidence | File:Line |
|------|----------|-----------|
| Monorepo layout | Eight libs under `libs/`, each with own `pyproject.toml`, `uv.lock`, Makefile | `Makefile:3-5`, `libs/*/pyproject.toml` |
| Documented dependency map | Maintainer-maintained graph of downstream dependents | `AGENTS.md:33-55` |
| Shared namespace | No top-level `langgraph/__init__.py` in any lib (implicit PEP 420 namespace) | `libs/checkpoint/langgraph/` (no `__init__.py` at namespace root), verified across all 7 Python libs |
| Wheel packaging scope | Each wheel ships only its `langgraph/*` subtree | `libs/langgraph/pyproject.toml:129-130`, `libs/prebuilt/pyproject.toml:70-71` |
| Core deps | Core depends on langchain-core, langgraph-checkpoint, langgraph-sdk, langgraph-prebuilt | `libs/langgraph/pyproject.toml:26-33` |
| Core→SDK usage | `RemoteGraph` and remote run streaming built on `langgraph_sdk.client` / streams | `libs/langgraph/langgraph/pregel/remote.py:26-41`, `libs/langgraph/langgraph/pregel/_remote_run_stream.py:11-14` |
| Core→checkpoint usage | Pregel loop imports `langgraph.checkpoint.base`, `langgraph.store.base` | `libs/langgraph/langgraph/pregel/main.py:47-52` |
| Checkpoint base is dep-minimal | Only langchain-core + ormsgpack | `libs/checkpoint/pyproject.toml:14-17` |
| Duck-typed optional backend | `RedisCache` accepts any client object; no redis runtime dep | `libs/checkpoint/langgraph/cache/redis/__init__.py:11-16` |
| Prebuilt undeclared dep | Production deps omit `langgraph`; core appears only in test group via editable path | `libs/prebuilt/pyproject.toml:26-29`, `libs/prebuilt/pyproject.toml:44`, `libs/prebuilt/pyproject.toml:65` |
| Prebuilt→core code flow | `create_react_agent` imports `langgraph.graph`, `langgraph.types`, `langgraph.runtime`, etc. | `libs/prebuilt/langgraph/prebuilt/chat_agent_executor.py:32-43` |
| Prebuilt→core private APIs | ToolNode imports `langgraph._internal._constants`, `_runnable`, and `langgraph.pregel._tools._tool_call_writer` | `libs/prebuilt/langgraph/prebuilt/tool_node.py:85-92` |
| Coupling acknowledged | Docstring states core hook file is "Read by `ToolRuntime.emit_output_delta` (in `langgraph.prebuilt`)" | `libs/langgraph/langgraph/pregel/_tools.py:31` |
| Core bundles prebuilt without importing | `langgraph-prebuilt` declared as prod dep; zero `langgraph.prebuilt` imports in core source (only docstring mentions) | `libs/langgraph/pyproject.toml:30`, grep over `libs/langgraph/langgraph/**` found matches only at `libs/langgraph/langgraph/pregel/_tools.py:31`, `libs/langgraph/langgraph/runtime.py:138` |
| CLI declares unused dep | `langgraph-sdk>=0.1.0` declared; no `langgraph_sdk` import anywhere in `libs/cli` source | `libs/cli/pyproject.toml:17` (grep over `libs/cli/**/*.py` excluding lockfile returned nothing) |
| CLI optional extras point outside repo | `inmem` extra requires external `langgraph-api` / `langgraph-runtime-inmem` distributions | `libs/cli/pyproject.toml:24-28` |
| sdk-py standalone | Depends only on httpx, orjson, langchain-protocol, langchain-core, websockets; zero sibling-lib imports | `libs/sdk-py/pyproject.toml:14-20`, grep over `libs/sdk-py/langgraph_sdk/**` returned none |
| Internal-API marker | `_internal` docstring: "not part of the public API, and thus stability is not guaranteed" | `libs/langgraph/langgraph/_internal/__init__.py:1-4` |
| Public surface pinning | `__all__` on public modules: graph, pregel, prebuilt, stream, sdk, func, ui | `libs/langgraph/langgraph/graph/__init__.py:5-12`, `libs/langgraph/langgraph/pregel/__init__.py:2`, `libs/prebuilt/langgraph/prebuilt/__init__.py:14-23`, `libs/langgraph/langgraph/stream/__init__.py:28-44`, `libs/sdk-py/langgraph_sdk/__init__.py:1-7`, `libs/langgraph/langgraph/graph/ui.py:10-18` |
| Versioned deprecations | Warning classes carry `since` and `expected_removal` versions | `libs/langgraph/langgraph/warnings.py:25-48` |
| Experimental markers | `@beta(message="The v3 streaming protocol on Pregel is experimental.")` on new stream API | `libs/langgraph/langgraph/stream/run_stream.py:30,303`, `libs/langgraph/langgraph/pregel/main.py:3519,3575` |
| Export regression test | Tests assert sdk client symbols remain importable "to ensure backwards compatibility during refactoring" | `libs/sdk-py/tests/test_client_exports.py:1-78` |
| Conformance suite for checkpointers | `validate()` checks blob round-trips, metadata, namespace isolation for any `BaseCheckpointSaver` subclass | `libs/checkpoint-conformance/langgraph/checkpoint/conformance/validate.py:45`, capability table in `libs/checkpoint-conformance/README.md` |
| Per-lib CI isolation | Matrix runs lint+tests per lib dir on Python 3.10–3.14 with `uv sync --frozen --group test --no-dev` (only declared deps) and a clean-worktree assertion afterwards | `.github/workflows/ci.yml:58-107`, `.github/workflows/_test.yml:19-22,45,50,53-56` |
| Dev-time sibling wiring | Sibling libs wired via editable path sources under `[tool.uv.sources]` | `libs/langgraph/pyproject.toml:83-89`, `libs/prebuilt/pyproject.toml:64-68` |

## Answers to Dimension Questions

**1. Are modules cleanly separated?**
Mostly yes, with one structural exception. Runtime (core `langgraph`), tools/agents (prebuilt), persistence (checkpoint family), server SDK (sdk-py), and CLI occupy separate distributions under separate directories with separate locks and version numbers. The seams that are clean: checkpoint base ↔ implementations (postgres/sqlite/conformance depend only on the base, `libs/checkpoint-postgres/pyproject.toml:14-19`), sdk-py ↔ everything (no internal imports), CLI ↔ core (imports only its own `langgraph_cli.*` modules). The unclean seam: prebuilt ⇄ core. Prebuilt's entire agent/tool API is implemented *on top of* core classes (`libs/prebuilt/langgraph/prebuilt/chat_agent_executor.py:32-43`) yet ships as if it were a lower layer than core.

**2. Do dependencies flow in one direction?**
At the level of *declared* dependencies, yes — the graph in `AGENTS.md:38-53` is acyclic and matches every `pyproject.toml`. At the level of *actual imports*, no: prebuilt → core imports form an undeclared cycle because core declares prebuilt (`libs/langgraph/pyproject.toml:30`) while prebuilt imports core (`tool_node.py:85-92`) without declaring it. Two smaller declaration/usage mismatches compound the picture: the CLI declares `langgraph-sdk` but never imports it (`libs/cli/pyproject.toml:17`), and core declares prebuilt but never imports it. Net effect: declarations describe intent better than reality in both directions.

**3. Can modules be used independently?**
Yes for checkpoint base (+ postgres/sqlite/conformance implementations), sdk-py, and the CLI — each installs and imports with only third-party deps. No for `langgraph-prebuilt`: although published as a standalone distribution, importing anything from `langgraph.prebuilt` triggers `from langgraph._internal...` / `from langgraph.graph...` at module load time, so a bare `pip install langgraph-prebuilt` yields a package whose `__init__.py` raises `ModuleNotFoundError` unless something else installed core `langgraph` first. In practice users always install `langgraph`, which pulls prebuilt, masking the defect. Answering the dimension's guiding question — "can you use the tool system without pulling in the entire runtime?" — the evidence says **no**: the tool layer (`ToolNode`, `ToolRuntime`, `create_react_agent`) is hard-wired to the full core runtime.

**4. Are public APIs distinguished from internal ones?**
Yes, through several coordinated conventions: an explicitly disclaimed `langgraph._internal` package (`_internal/__init__.py:1-4`); underscore-private modules inside otherwise-public packages (`stream/_types.py`, `pregel/_loop.py`, `checkpoint/serde/_msgpack.py`); curated `__all__` tuples on all public modules; versioned deprecation warnings encoding removal timelines (`warnings.py:29-48`); and `@beta` decorators marking unstable surfaces (`stream/run_stream.py:30`). Enforcement is partial: sdk-py has a dedicated export-compatibility test (`test_client_exports.py:35`), and syrupy snapshot tests exist for graph behavior (`tests/__snapshots__/test_pregel.ambr`), but no test pins the core or prebuilt export lists, and nothing mechanically prevents internal imports across package boundaries.

## Architectural Decisions

1. **Single shared namespace across distributions** (`langgraph.*`, PEP 420 implicit namespaces). All seven Python libs contribute subtrees (`langgraph.checkpoint.*`, `langgraph.prebuilt.*`, `langgraph_sdk.*`, `langgraph_cli.*`) without a root `__init__.py`. This lets persistence, SDK, and runtime evolve and release independently while presenting one coherent import tree. Cost: ownership of a namespace path is purely conventional — nothing stops one lib from writing into another's subtree, which is exactly how prebuilt ends up coupled to core internals.

2. **Persistence as a separate dependency axis.** Checkpointing was extracted into `langgraph-checkpoint` with versioned implementations (`>=4.1.0,<5.0.0` pins in `libs/langgraph/pyproject.toml:28`, `libs/checkpoint-postgres/pyproject.toml:15`). The core runtime consumes only the base interfaces (`pregel/main.py:47-52`), so storage backends are pluggable and testable without the runtime.

3. **Conformance library instead of trust.** Third-party checkpointer authors get an executable contract (`checkpoint-conformance/validate.py:45`) covering required vs auto-detected extended capabilities (README capability table). This converts the base-class interface from documentation into an enforced boundary — the strongest single mechanism in the repo for keeping the persistence boundary honest.

4. **Versioned deprecation ladder.** `LangGraphDeprecatedSinceV05/V10/V11` classes each hard-code `since` and `expected_removal` versions (`warnings.py:51-69`), giving consumers machine-readable migration signals rather than free-form warnings.

5. **Per-lib CI cells as the boundary firewall.** Each lib is linted and tested in its own directory with `uv sync --frozen --group test --no-dev` across Python 3.10–3.14, then checked for a clean worktree (`.github/workflows/_test.yml:45,50,53-56`). Because dev installs resolve siblings via editable path sources (`[tool.uv.sources]`), a change to checkpoint or prebuilt is exercised against every dependent lib in the same PR — the practical mechanism behind the AGENTS.md warning "Changes to a library may impact all of its dependents" (`AGENTS.md:55`).

6. **Dependency-injection for heavy backends.** `RedisCache` takes a duck-typed `redis: Any` client (`cache/redis/__init__.py:12-16`) so the base distribution avoids a redis dependency; similarly the CLI keeps its dev-server behind the optional `inmem` extra pointing at out-of-repo distributions (`libs/cli/pyproject.toml:24-28`).

## Notable Patterns

- **Curated `__all__` everywhere**: every public module enumerates exports (`graph/__init__.py:5-12`, `prebuilt/__init__.py:14-23`, `stream/__init__.py:26-44`, `sdk-py/__init__.py:1-7`), making the public surface greppable and diffable.
- **Private-by-prefix convention at three scales**: whole package (`_internal/`), submodule (`stream/_types.py`), symbol (`_tool_call_writer`).
- **Docstring-negotiated cross-package coupling**: when prebuilt must read a core private hook, the coupling is documented on the core side (`pregel/_tools.py:31` says the file is "Read by `ToolRuntime.emit_output_delta` (in `langgraph.prebuilt`)"). Honest, but manual.
- **Snapshot tests (syrupy)** pin observable behavior of compiled graphs (`tests/__snapshots__/test_pregel.ambr`, `prebuilt/tests/__snapshots__/test_react_agent_graph.ambr`), indirectly protecting boundary-relevant semantics like stream event shapes during refactors.
- **Export-compatibility regression tests** in sdk-py framed explicitly as refactor safety (`test_client_exports.py:1-4`).

## Tradeoffs

- **Namespace composition vs enforceable boundaries**: the shared `langgraph` namespace gives users a unified API but removes any structural barrier between packages; only conventions separate public from cross-package-private.
- **Bundling vs layering for prebuilt**: declaring prebuilt as a core dependency maximizes convenience (`pip install langgraph` yields batteries included) but froze the architecture into a cycle once prebuilt started building on core types (`Runtime`, `CompiledStateGraph`). A strict-layer design would have required either duplicating those types or accepting prebuilt-as-upper-layer with core not depending on it.
- **Editable-path dev installs vs hermetic builds**: `[tool.uv.sources]` editable wiring makes monorepo development fast and keeps dependents exercised, but means local dev never reproduces the "fresh consumer installs published wheels" scenario where prebuilt's missing dependency would surface immediately.
- **Duck typing vs declared optional deps**: `redis: Any` keeps the base lib lean at the cost of static checking on the client object (`cache/redis/__init__.py:12`).

## Failure Modes / Edge Cases

- **Standalone `langgraph-prebuilt` install breaks at import**: any `import langgraph.prebuilt` fails with `ModuleNotFoundError` unless core `langgraph` happens to be present transitively. No `ImportError` guard or lazy import mitigates this; the failure occurs inside `prebuilt/__init__.py:3-12`.
- **Silent breakage across releases via private imports**: because prebuilt consumes `langgraph._internal._runnable` and `langgraph.pregel._tools._tool_call_writer` (`tool_node.py:85-89`), a routine core-internal refactor can break the published prebuilt wheel even when both pass their own CI (each tested against same-commit siblings, not against previously released versions). The `_internal` disclaimer puts the burden on prebuilt maintainers.
- **Undeclared-dependency drift is undetectable by current tooling**: the CLI shipping an unused `langgraph-sdk` dep and prebuilt shipping missing core deps shows declarations are neither generated nor validated against imports; ruff configs select only `E,F,I,TID251,UP` (`libs/prebuilt/pyproject.toml:77-79`) — TID251 bans `typing.TypedDict`, not cross-package imports.
- **Namespace collision risk**: two distributions owning overlapping `langgraph.*` paths would silently shadow each other given PEP 420; there is no build-time check that wheel contents don't overlap (all wheels simply `include = ["langgraph"]`).

## Future Considerations

- Declare the real dependency edge: add `langgraph` to `langgraph-prebuilt`'s production deps and drop it from core's (making prebuilt a true upper layer), or move the shared types (`RunnableCallable`, `MISSING`, stream protocols) down into `langgraph-checkpoint`/a leaf package so prebuilt needs nothing from core. Either resolves the cycle; the second preserves today's user-facing bundling.
- Add an import-contract linter (e.g., import-linter "forbidden/layered" contracts) to CI per lib, encoding the AGENTS.md map so violations fail builds instead of relying on review.
- Generate or validate dependency declarations from actual imports (e.g., in a `check-sdk-methods`-style CI job like the existing `.github/workflows/ci.yml:109-120`) to catch both directions of drift (missing prebuilt dep, unused cli dep).
- Pin the public API surface of core and prebuilt with export-snapshot tests analogous to `sdk-py/tests/test_client_exports.py`, since those are the highest-churn boundaries.
- Promote frequently consumed `_internal` members (e.g., `RunnableCallable`) into a documented stable location, shrinking the private-API surface prebuilt depends on.

## Questions / Gaps

- No automated circularity check was found anywhere in the repo (searched `.github/workflows/`, pyproject/ruff configs, Makefiles): "No clear evidence found" for mechanical cycle detection beyond human review.
- Whether `langgraph-prebuilt`'s missing core dependency breaks real users could not be observed directly (no issue tracker in scope); the claim rests on static import analysis of `prebuilt/__init__.py:3-12` and its transitive modules.
- The repo docs directory contains only redirects/llms.txt (`docs/`), so the canonical public-API reference lives outside this source; the study treats `__all__` surfaces as the operative definition of "public".
- sdk-js is a stub pointing to `langchain-ai/langgraphjs` (`libs/sdk-js/README.md:7`); its boundary properties are out of scope for this isolated source.

---

Generated by `dimensions/22.01-package-and-module-boundaries.md` against `langgraph`.
