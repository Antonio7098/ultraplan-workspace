# Source Analysis: langgraph

## Dimension 24.01: Public API Surface

### Source Info

| Field | Value |
|-------|-------|
| Name | langgraph |
| Path | `studies/agent-harness-study/sources/langgraph` |
| Language / Stack | Python (monorepo of 8 packages) + JS/TS SDK stub; hatchling build, uv, ruff/ty linting, pytest |
| Analyzed | 2026-08-22 |

## Summary

LangGraph's public API surface is organized as a monorepo of independently versioned packages (`libs/langgraph` core at v1.2.6 per `libs/langgraph/pyproject.toml:7`, `libs/prebuilt` at v1.1.0 per `libs/prebuilt/pyproject.toml:7`, plus `checkpoint*`, `sdk-py`, and `cli`). The most distinctive design choice is that the top-level `langgraph` package deliberately ships **no** `__init__.py` (confirmed via `git ls-files libs/langgraph/langgraph`); the public API is exposed exclusively through curated submodule entry points — `langgraph.graph` (`libs/langgraph/langgraph/graph/__init__.py:1-12`), `langgraph.func`, `langgraph.types`, `langgraph.config`, `langgraph.errors`, `langgraph.pregel`, `langgraph.channels`, `langgraph.stream`, `langgraph.runtime`, and `langgraph.warnings`. Each uses an explicit `__all__`, so the stable import paths are enumerable from source.

The surface is layered by abstraction level: low-level graph construction (`StateGraph.compile()` at `libs/langgraph/langgraph/graph/state.py:1164-1259`), a functional API (`task`/`entrypoint` in `libs/langgraph/langgraph/func/__init__.py:56,132,262`), high-level agent factories (`create_react_agent` in `libs/prebuilt/langgraph/prebuilt/chat_agent_executor.py:272-308`, now formally deprecated in favor of `langchain.agents.create_agent`), a server SDK (`get_client` in `libs/sdk-py/langgraph_sdk/_async/client.py:29`), an operator CLI (`langgraph = "langgraph_cli.cli:cli"` in `libs/cli/pyproject.toml:36-37`), and checkpointer/store extension base classes (`BaseCheckpointSaver` at `libs/checkpoint/langgraph/checkpoint/base/__init__.py:176`, `BaseStore` exported in `libs/checkpoint/langgraph/store/base/__init__.py:1299-1313`). Internals are fenced into a dedicated `langgraph._internal` package whose docstring states "This module is not part of the public API" (`libs/langgraph/langgraph/_internal/__init__.py:1-4`), and lifecycle is actively managed through versioned deprecation warning classes (`libs/langgraph/langgraph/warnings.py:13-69`) backed by tests (`libs/langgraph/tests/test_deprecation.py:28-379`).

Weaknesses: documentation has been moved off-repo to docs.langchain.com (`docs/llms.txt:3`), the in-repo examples directory is explicitly archival ("no longer updated", `examples/README.md:3-5`), the JS SDK directory is an empty pointer stub (`libs/sdk-js/README.md:7`), the checkpoint base module lacks an explicit `__all__` (grep across `libs/checkpoint/langgraph/` found none for checkpoint/base), and one legacy shim module leaks internal classes under a public path while claiming removal "in v1" although the package is already at v1.2.6 (`libs/langgraph/langgraph/utils/runnable.py:1-2`).

## Rating

**8 / 10.** The public API model is clear, explicitly bounded (`__all__` everywhere in the core lib, `_internal` segregation, versioned deprecation taxonomy with test enforcement), and runnable examples are embedded directly into docstrings (`libs/langgraph/langgraph/config.py:59-121`, `libs/langgraph/langgraph/func/__init__.py:174-214,319-434`). It falls short of 9–10 because discoverability depends on external hosted docs, the checkpoint package does not meet the same export-discipline standard as core, compat shims outlive their stated removal targets, and experimental surfaces (`stream_events(version="v3")`) are reachable from public methods with only runtime warnings.

## Evidence Collected

| Area | Evidence | File:Line |
|------|----------|-----------|
| No top-level package `__init__.py`; API via submodules only | `git ls-files` shows no `libs/langgraph/langgraph/__init__.py`; only subpackage inits exist | `libs/langgraph/langgraph/` (tree) |
| Primary graph-building entry point | `from langgraph.graph.state import StateGraph`; `__all__ = ("END","START","StateGraph","add_messages","MessagesState","MessageGraph")` | `libs/langgraph/langgraph/graph/__init__.py:1-12` |
| Graph compile lifecycle | `def compile(...) -> CompiledStateGraph` documents checkpointer/thread_id contract; implements Runnable interface | `libs/langgraph/langgraph/graph/state.py:1164-1223` |
| Functional API exports exactly two symbols | `__all__ = ("task", "entrypoint")` with docstring examples for sync/async tasks | `libs/langgraph/langgraph/func/__init__.py:56,174-214` |
| Public types module | `__all__` lists Command/interrupt/Send/StreamMode/RetryPolicy/etc. | `libs/langgraph/langgraph/types.py:52-85` |
| Runtime injection objects | `__all__ = ("BaseUser","ExecutionInfo","RunControl","Runtime","ServerInfo","get_runtime")` | `libs/langgraph/langgraph/runtime.py:16-23` |
| Curated channel abstractions | Grouped `__all__`: BaseChannel + value/topic channels | `libs/langgraph/langgraph/channels/__init__.py:14-29` |
| Internal boundary declaration | "This module is not part of the public API, and thus stability is not guaranteed." | `libs/langgraph/langgraph/_internal/__init__.py:1-4` |
| High-level prebuilt API | `langgraph.prebuilt.__all__` = create_react_agent, ToolNode, ToolCallTransformer, tools_condition, ValidationNode, InjectedState/Store, ToolRuntime | `libs/prebuilt/langgraph/prebuilt/__init__.py:14-23` |
| create_react_agent formally deprecated | `@deprecated("...moved to \`langchain.agents\`...", category=LangGraphDeprecatedSinceV10)` above signature | `libs/prebuilt/langgraph/prebuilt/chat_agent_executor.py:272-276` |
| Versioned deprecation warnings as API policy | Classes encode since/expected_removal versions (v0.5→v2.0, v1.0→v2.0, v1.1→v3.0) | `libs/langgraph/langgraph/warnings.py:51-69` |
| Deprecated import redirection with warning | Module-level `__getattr__` warns then re-exports Send/Interrupt from `langgraph.types`; private constants warn "let the LangGraph team know if you need this value" | `libs/langgraph/langgraph/constants.py:34-62` |
| Deprecation behavior pinned by tests | 15+ tests assert exact warning categories/messages (retry→retry_policy, input→input_schema, MessageGraph removal, config_schema→context_schema) | `libs/langgraph/tests/test_deprecation.py:28-348` |
| Experimental streaming labeled @beta | `@beta(message="The v3 streaming protocol on Pregel is experimental.")` on `_pregel_stream_v3` / `_apregel_stream_v3`; public `stream_events(version="v3")` docstrings repeat "experimental and may change" | `libs/langgraph/langgraph/pregel/main.py:3519,3575,3630-3674` |
| New stream module exported but beta-flagged internals | `GraphRunStream`, transformers etc. re-exported; module decorated `@beta` | `libs/langgraph/langgraph/stream/__init__.py:16-45`; `libs/langgraph/langgraph/stream/run_stream.py:30,303` |
| SDK client factory | Top-level `__all__ = ["Auth","Encryption","EncryptionContext","get_client","get_sync_client"]`, version 0.4.2 | `libs/sdk-py/langgraph_sdk/__init__.py:1-8` |
| SDK resource-scoped clients | `LangGraphClient` composes AssistantsClient/CronClient/RunsClient/StoreClient/ThreadsClient | `libs/sdk-py/langgraph_sdk/_async/client.py:8-12,143` |
| SDK quick-start runnable example | get_client → assistants.search → threads.create → runs.stream snippet | `libs/sdk-py/README.md:31-46` |
| Server auth extension point | `class Auth` documented with `langgraph.json` wiring example for custom auth handlers | `libs/sdk-py/langgraph_sdk/auth/__init__.py:13-40` |
| CLI command group | `[project.scripts] langgraph = "langgraph_cli.cli:cli"`; commands up/build/dockerfile/dev/validate/new/prepare defined | `libs/cli/pyproject.toml:36-37`; `libs/cli/langgraph_cli/cli.py:232,278,419,550,763,873,920,982` |
| Checkpointer extension base class | `class BaseCheckpointSaver(Generic[V])`; InMemorySaver with debug-only guidance; MemorySaver alias kept for backcompat | `libs/checkpoint/langgraph/checkpoint/base/__init__.py:176`; `libs/checkpoint/langgraph/checkpoint/memory/__init__.py:33,41,631` |
| Store extension base class | Explicit `__all__` for BaseStore, Item, Op types, embedding helpers | `libs/checkpoint/langgraph/store/base/__init__.py:1299-1313` |
| Beta metadata field marked in schema | `counters_since_delta_snapshot` carries `!!! warning "Beta"` note tied to DeltaChannel stabilization | `libs/checkpoint/langgraph/checkpoint/base/__init__.py:60-77` |
| Docstrings carry runnable examples | get_store/get_stream_writer show complete StateGraph + functional-API snippets with expected outputs | `libs/langgraph/langgraph/config.py:32-196` |
| Docs relocated off-repo | "LangGraph documentation has moved to docs.langchain.com"; reference link to reference.langchain.com | `docs/llms.txt:3,31` |
| Examples archive-only | "retained purely for archival purposes and is no longer updated" | `examples/README.md:3-5` |
| JS SDK is a redirect stub | "This package has moved to langchain-ai/langgraphjs" | `libs/sdk-js/README.md:7` |
| Compat shims leak internals | `utils/runnable.py` re-exports `RunnableCallable`, `RunnableLike` from `langgraph._internal._runnable`; header says "to be removed in v1" though package version is 1.2.6 | `libs/langgraph/langgraph/utils/runnable.py:1-2`; `libs/langgraph/pyproject.toml:7` |
| Package maturity classifier | `'Development Status :: 5 - Production/Stable'` on both core and prebuilt | `libs/langgraph/pyproject.toml:15`; `libs/prebuilt/pyproject.toml:17` |

## Answers to Dimension Questions

### 1. What is the intended public API surface?

Five tiers, each its own distribution:

1. **Core framework** (`pip install langgraph`): `StateGraph`/`START`/`END`/`MessagesState` (`libs/langgraph/langgraph/graph/__init__.py:5-12`), functional `task`/`entrypoint` (`libs/langgraph/langgraph/func/__init__.py:56`), control/data types like `Command`, `interrupt`, `Send`, `StreamMode`, `RetryPolicy` (`libs/langgraph/langgraph/types.py:52-85`), runtime accessors `get_config`/`get_store`/`get_stream_writer` (`libs/langgraph/langgraph/config.py:17,32,126`), errors with machine-readable codes (`ErrorCode` at `libs/langgraph/langgraph/errors.py:34-39`), and the low-level `Pregel` engine itself (`libs/langgraph/langgraph/pregel/__init__.py:1-3`).
2. **Prebuilt agents** (`langgraph-prebuilt`): eight symbols around `create_react_agent` and `ToolNode` (`libs/prebuilt/langgraph/prebuilt/__init__.py:14-23`).
3. **Persistence/storage contracts** (`langgraph-checkpoint[-sqlite/-postgres]`): `BaseCheckpointSaver`, serializers, `BaseStore`, cache bases — the seams third parties implement.
4. **Server interaction**: `langgraph-sdk` Python client (`libs/sdk-py/langgraph_sdk/__init__.py:8`) plus the `Auth`/`Encryption` extension hooks, and the `langgraph` CLI for operators (`libs/cli/pyproject.toml:36-37`).
5. **HTTP/RPC routes**: no route table lives in this repo; the REST surface is consumed via sdk-py resource clients (assistants/threads/runs/crons/store sub-clients composed at `libs/sdk-py/langgraph_sdk/_async/client.py:143+`). No OpenAPI spec was found inside the studied source (searched `libs/sdk-py` for schema/route definitions; only typed client methods exist).

### 2. Is the stable API easy to distinguish from internal implementation details?

Yes, unusually well. Three mechanisms work together: (a) the underscore convention consolidated into `langgraph._internal` with an explicit stability disclaimer (`libs/langgraph/langgraph/_internal/__init__.py:1-4`); (b) exhaustive `__all__` declarations in every public module (e.g., `libs/langgraph/langgraph/types.py:52-85`, `libs/prebuilt/langgraph/prebuilt/__init__.py:14-23`); (c) deprecated legacy import paths that fail loud rather than silently working — importing `Send` from `langgraph.constants` emits `LangGraphDeprecatedSinceV10` (`libs/langgraph/constants.py:34-46`, tested at `libs/langgraph/tests/test_deprecation.py:86-97`), and even private constants now behind `__getattr__` tell users to contact the team if they need them (`libs/langgraph/constants.py:48-60`). Gaps: the checkpoint base package predates this discipline and has no `__all__` (grep over `libs/checkpoint/langgraph/` returned matches only for store/base), and the `langgraph.utils` shims keep internal names (`RunnableCallable`) importable from a non-underscored path (`libs/langgraph/langgraph/utils/runnable.py:2`).

### 3. Does the API expose the right level of abstraction for agent harness users?

The layering supports harness integrators well. A minimal harness can target `StateGraph` + checkpointer without touching Pregel internals; the compiled graph implements LangChain's `Runnable` interface so invoke/stream/batch come for free (`libs/langgraph/langgraph/graph/state.py:1176-1180`). Extension authors implement well-defined base classes (`BaseCheckpointSaver`, `BaseStore`) rather than protocols inferred from usage. Escape hatches are graded: `Pregel` and channel primitives are public but positioned below `StateGraph` in the docs hierarchy, and injectable runtime parameters (`previous`, `runtime`, `config` documented as a table in `libs/langgraph/langgraph/func/__init__.py:276-281`) avoid global lookups where typing can express them. Two friction points: `NodeBuilder` is exported alongside `Pregel` (`libs/langgraph/langgraph/pregel/__init__.py:1-3`) despite being a low-level assembly tool used mainly in tests (e.g., `libs/langgraph/tests/test_deprecation.py:154`), and `interrupt_before`/`interrupt_after`/`pre_model_hook`/`post_model_hook` on the prebuilt agent expose graph topology decisions through the high-level factory's signature (`libs/prebuilt/langgraph/prebuilt/chat_agent_executor.py:296-303`).

### 4. Are examples sufficient to use the API correctly without reading internals?

For the Python core and SDK, yes — examples live in the primary reference (docstrings) with expected outputs, e.g. the full interrupt/resume workflow including `Command(resume=...)` in `libs/langgraph/langgraph/func/__init__.py:319-380`, `entrypoint.final` round-trip arithmetic (`libs/langgraph/langgraph/func/__init__.py:403-434`), and store/stream-writer recipes with `pycon` output blocks (`libs/langgraph/langgraph/config.py:59-121,137-193`). The SDK README shows an end-to-end client session (`libs/sdk-py/README.md:31-46`). However, repo-level example coverage regressed intentionally: `examples/` is frozen (`examples/README.md:3-5`) and all narrative docs moved to an external site (`docs/llms.txt:3`), so correctness of examples is no longer exercised against the codebase in-repo (no doctest runner found in `libs/*/Makefile` targets). Error-path examples are strong: exceptions embed troubleshooting URLs generated from `ErrorCode` values (`libs/langgraph/langgraph/errors.py:42-47`).

## Architectural Decisions

- **Submodule-as-entry-point instead of root namespace re-export.** With no `libs/langgraph/langgraph/__init__.py`, imports must name a facet (`langgraph.graph`, `langgraph.func`, ...), which keeps each facet's dependency direction explicit and avoids a mega-module that would force circular-import management. Cost: discoverability relies entirely on external docs/reference site.
- **Monorepo with decoupled versioning and declared compatibility ranges.** Core pins `langgraph-checkpoint>=4.1.0,<5` and `langgraph-prebuilt>=1.1.0,<1.2` (`libs/langgraph/pyproject.toml:26-33`), letting storage backends and agent helpers evolve on separate cadences while the dependency map in `AGENTS.md` documents blast radius.
- **Deprecation taxonomy as first-class API.** Versioned warning subclasses carrying `since`/`expected_removal` tuples (`libs/langgraph/langgraph/warnings.py:29-48`) turn breaking-change communication into typed, greppable, testable artifacts rather than changelog prose.
- **Internal consolidation under `langgraph._internal`.** Formerly scattered private modules were gathered behind one disclaimer boundary, and the ruff config additionally bans `typing.TypedDict` in favor of `typing_extensions` (`libs/langgraph/pyproject.toml:99-100`), showing import hygiene enforced by tooling, not just convention.
- **SDK mirrors server resources as client sub-objects.** `LangGraphClient` exposes `.assistants/.threads/.runs/.crons/.store` (`libs/sdk-py/langgraph_sdk/_async/client.py:8-12`), matching operator mental models of the deployment API rather than exposing raw HTTP verbs.

## Notable Patterns

- **Docstring-driven developer experience:** every major public symbol carries a runnable example block with expected output, effectively making docstrings the executable tutorial (`libs/langgraph/langgraph/func/__init__.py:174-214`; `libs/langgraph/langgraph/config.py:90-92`).
- **Graceful migration ladders:** old symbols are not deleted but rerouted with escalating warnings — `MemorySaver = InMemorySaver` alias (`libs/checkpoint/langgraph/checkpoint/memory/__init__.py:631`), `NodeInterrupt` decorated deprecated pointing to `langgraph.types.interrupt` (`libs/langgraph/langgraph/errors.py:110-115`), `create_react_agent` pointing to `langchain.agents.create_agent` (`libs/prebuilt/langgraph/prebuilt/chat_agent_executor.py:272-276`).
- **Beta gating at method granularity:** the v3 stream protocol is reachable through public `stream_events(version="v3")` but implemented in underscored `@beta` methods `_pregel_stream_v3`/`_apregel_stream_v3` (`libs/langgraph/langgraph/pregel/main.py:3519,3575`), keeping the experimental transport out of the nominal method body.
- **Machine-readable failure taxonomy:** `ErrorCode` enum + message factory that appends a deterministic troubleshooting URL per code (`libs/langgraph/langgraph/errors.py:34-47`).
- **Typed extensibility:** generic typevars with variance published for integrators (`StateT_co`, `ContextT_contra`, etc. in `libs/langgraph/langgraph/typing.py:7-30`).

## Tradeoffs

- **Discoverability vs. hygiene:** omitting a root `__init__.py` prevents accidental god-imports but means IDE auto-complete on `import langgraph` yields nothing; users must already know the submodule map that `AGENTS.md` records for contributors, not consumers.
- **Docs centralization vs. reproducibility:** moving all guides off-repo (`docs/llms.txt:3`) eliminates doc drift within the repo but severs the compile-time link between docs and code; nothing in this checkout validates that documented signatures still exist.
- **Explicit `__all__` vs. maintenance cost:** every addition requires editing the export list (visible in the grouped comments of `libs/langgraph/langgraph/channels/__init__.py:14-29`); the checkpoint package evidently opted out, creating inconsistent auditability between packages.
- **Long-lived deprecation shims vs. clean slate:** aliases and `__getattr__` redirects preserve downstream ecosystems but accumulate — `langgraph/utils` claims removal "in v1" (`libs/langgraph/langgraph/utils/__init__.py:1`) yet persists in v1.2.6, signaling removal targets are aspirational.
- **Rich factory signatures vs. abstraction leakage:** `create_react_agent(..., interrupt_before=..., post_model_hook=..., version=...)` (`libs/prebuilt/langgraph/prebuilt/chat_agent_executor.py:296-305`) offers one-call power but couples callers to graph-shape concepts that the deprecation notice says will move to middleware anyway.

## Failure Modes / Edge Cases

- **Silent divergence of compat shims:** `RunnableLike`/`RunnableCallable` remain importable from `langgraph.utils.runnable` (`libs/langgraph/langgraph/utils/runnable.py:2`); if `_internal` refactors rename them, this public path breaks without any stability guarantee applying to it.
- **Experimental surface reachable by default-typing users:** `stream_events(version="v3")` is a one-string change from the stable `"v2"` default (`libs/langgraph/langgraph/pregel/main.py:3630-3646`); protection is a runtime `@beta` warning, not a type-level barrier, so CI without `-W error` will not catch adoption.
- **Deprecated kwargs accepted indefinitely:** `task(retry=...)` and `entrypoint(config_schema=...)` still function through `**kwargs` sniffing (`libs/langgraph/langgraph/func/__init__.py:216-223,449-456`); code depending on warning-as-error policies breaks loudly, but plain users get no hard signal until v2/v3 removal actually ships.
- **Checkpoint package export drift:** without `__all__`, new names added to `libs/checkpoint/langgraph/checkpoint/base/__init__.py` (it already transitively imports serializer internals like `JsonPlusSerializer` at lines 12-24) become public implicitly, and removing any would be an untracked breaking change.
- **Python-version-dependent availability:** `get_config` raises on async contexts before Python 3.11 (`libs/langgraph/langgraph/config.py:18-25`), and docstrings warn `get_store`/`get_stream_writer` silently won't propagate contextvars on older runtimes (`libs/langgraph/langgraph/config.py:53-56`) — a documented but easy-to-miss environmental constraint baked into the API contract.

## Future Considerations

- Add an explicit `__all__` (and ideally an export snapshot test) to `libs/checkpoint/langgraph/checkpoint/base/__init__.py` to match core-lib discipline.
- Either honor or update the "to be removed in v1" statements in `libs/langgraph/langgraph/utils/*.py` — stale removal promises erode trust in the versioned-deprecation system.
- Consider a lightweight in-repo doctest/smoke runner for docstring examples so the embedded tutorials stay executable as the codebase evolves (Makefile currently covers format/lint/test/bench only, `libs/langgraph/Makefile:1-40`).
- Promote the v3 streaming protocol out of `@beta` once stable, or hide it behind a dedicated opt-in import, so type checkers can flag experimental adoption.
- Publish a single-page "public surface map" (the equivalent of AGENTS.md's contributor view) aimed at consumers, mitigating the no-root-init discoverability cost.

## Questions / Gaps

- **REST route ownership:** the LangGraph Server HTTP API itself (route handlers/OpenAPI spec) is not in this repository — only typed SDK clients (`libs/sdk-py/langgraph_sdk/_async/runs.py`) consume it. Route-level API-stability guarantees could not be assessed here (searched `libs/sdk-py`, `libs/cli` for server route tables; not present).
- **No export-surface regression tests found:** searched `libs/*/tests` for `test_import*`, `__all__` assertions, or API-snapshot fixtures; only deprecation-warning tests pin behavior (`libs/langgraph/tests/test_deprecation.py`). Whether upstream relies on the external `reference.langchain.com` generation to catch surface changes is unknown from this source.
- **sdk-js content:** the directory contains only a redirect README (`libs/sdk-js/README.md:7`); the actual TS SDK lives in another repository and was out of scope under source-isolation rules.
- **CLI programmatic API:** `langgraph_cli` modules (`config.py`, `schemas.py`, `docker.py`) have `py.typed` (`libs/cli/langgraph_cli/py.typed` exists) but no documented contract for importing the CLI as a library; treated as operator-facing binary only.

---

Generated by `Dimension 24.01: Public API Surface` against `langgraph`.
