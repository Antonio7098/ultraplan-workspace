# Source Analysis: langgraph

## Package and Module Boundaries

### Source Info

| Field | Value |
|-------|-------|
| Name | langgraph |
| Path | `studies/agent-harness-study/sources/langgraph` |
| Language / Stack | Python 3.10+ monorepo (hatchling-built wheels), LangChain ecosystem |
| Analyzed | 2026-08-21 |

## Summary

The repository is a Python monorepo (`libs/*`) that ships eight separately versioned PyPI packages coordinated by a uv workspace: `langgraph-checkpoint` (`libs/checkpoint`), `langgraph-checkpoint-conformance` (`libs/checkpoint-conformance`), `langgraph-checkpoint-postgres` (`libs/checkpoint-postgres`), `langgraph-checkpoint-sqlite` (`libs/checkpoint-sqlite`), `langgraph-cli` (`libs/cli`), `langgraph` (`libs/langgraph`), `langgraph-prebuilt` (`libs/prebuilt`), and `langgraph-sdk` (`libs/sdk-py`); `libs/sdk-js/` is a stub README pointing at a separate repository. The "core" package (`libs/langgraph/langgraph/`) is itself laid out as a layered graph: a stable leaf layer (`channels/`, `errors.py`, `types.py`, `constants.py`, `warnings.py`, `config.py`, `runtime.py`, `managed/`), a private utility layer (`langgraph/_internal/` with leading-underscore names), the Pregel execution layer (`langgraph/pregel/` with both `__init__.py` re-exports and many underscore-prefixed modules), the high-level builder API (`langgraph/graph/`, `langgraph/func/`, `langgraph/stream/`). The top-level `AGENTS.md` declares the downstream dependency map and instructs contributors to run `make format/lint/test` per-library (`AGENTS.md:1-50`). The `langgraph-checkpoint` package reuses the same Python `langgraph` namespace (its code lives under `libs/checkpoint/langgraph/checkpoint/`, `langgraph/store/`, `langgraph/cache/`), so the on-disk boundary between packages is the sub-package name, not the import root. Internal vs public surface is marked three ways: (1) a dedicated `langgraph/_internal/` package whose `__init__.py` declares "This module is not part of the public API, and thus stability is not guaranteed" (`libs/langgraph/langgraph/_internal/__init__.py:1-4`); (2) every public module ships a curated `__all__` (e.g. `libs/langgraph/langgraph/types.py:52-85`, `libs/langgraph/langgraph/errors.py:16-31`); (3) `langgraph.constants` exposes a module-level `__getattr__` that warns when external code reaches for re-exported private constants (`libs/langgraph/langgraph/constants.py:34-62`). All packages publish `py.typed` markers (verified with `find … -name py.typed`), which gives the boundary real teeth at type-check time. The model is mostly clean: independent packages that can be installed or tested in isolation (`libs/checkpoint-sqlite` and `libs/checkpoint-postgres` both depend only on `langgraph-checkpoint` plus their respective drivers; `libs/sdk-py/langgraph_sdk/` has zero production imports from `langgraph.*`); one core-only cycle is handled with a function-local import and an inline comment that calls it out ("ya we have a cyclic import here ¯\\_(ツ)_/¯", `libs/checkpoint/langgraph/checkpoint/serde/jsonplus.py:538-539`); `langgraph-prebuilt` reaches into `langgraph._internal._runnable/_constants/_typing` and `langgraph.pregel._tools`, which is a documented cross-package leak of private symbols.

## Rating

**8 / 10** — Clear, durable model with explicit public/internal separation, per-library wheels, `py.typed` markers, deprecation-aware `__getattr__` shims, and most independent packages actually being independent. Two soft spots keep it from a 9: (a) the core `langgraph` package depends on `langgraph-sdk` purely for one auth `Protocol` (`libs/langgraph/pyproject.toml:29`, `libs/langgraph/langgraph/runtime.py:8`), which is a reverse-direction coupling between two ostensibly independent products; (b) `langgraph-prebuilt` regularly reaches into `langgraph._internal.*` and `langgraph.pregel._*` private modules, so the underscore-prefix boundary is enforced socially, not by packaging.

## Evidence Collected

Every entry includes a file path with line numbers. Format: `path/to/file.py:NN`.

| Area | Evidence | File:Line |
|------|----------|-----------|
| Monorepo package layout | Eight sibling Python packages under `libs/`, each with its own `pyproject.toml` | `libs/` (ls) |
| Dependency map documented in repo | Textual map of downstream libraries per package | `AGENTS.md:19-43` |
| Core `langgraph` package | Wheel builds the `langgraph` Python package | `libs/langgraph/pyproject.toml:130` |
| Core dependencies | Declares `langchain-core`, `langgraph-checkpoint`, `langgraph-sdk`, `langgraph-prebuilt`, `xxhash`, `pydantic` | `libs/langgraph/pyproject.toml:26-33` |
| `langgraph-checkpoint` minimum deps | Only `langchain-core>=0.2.38` and `ormsgpack>=1.12.0` — no other LangGraph package | `libs/checkpoint/pyproject.toml:14-17` |
| `langgraph-checkpoint-postgres` deps | Only `langgraph-checkpoint>=4.1.0,<5.0.0` + `orjson` + `psycopg` + `psycopg-pool` | `libs/checkpoint-postgres/pyproject.toml:14-19` |
| `langgraph-checkpoint-sqlite` deps | Only `langgraph-checkpoint>=4.1.0,<5.0.0` + `aiosqlite` + `sqlite-vec` | `libs/checkpoint-sqlite/pyproject.toml:14-18` |
| `langgraph-checkpoint-conformance` deps | Only `langgraph-checkpoint>=2.0.0` (test-suite only) | `libs/checkpoint-conformance/pyproject.toml:13-15` |
| `langgraph-prebuilt` deps | `langgraph-checkpoint>=2.1.0,<5.0.0` and `langchain-core>=1.3.1`; `langgraph` is added via `[tool.uv.sources]` path source, not a runtime dep | `libs/prebuilt/pyproject.toml:26-29`, `libs/prebuilt/pyproject.toml:64-68` |
| `langgraph-cli` deps | `click`, `httpx`, `langgraph-sdk` (>=3.11 only), `pathspec`, `python-dotenv` — no `langgraph` import | `libs/cli/pyproject.toml:14-21` |
| `langgraph-cli` does not actually use langgraph-sdk | No `import langgraph_sdk` anywhere in `langgraph_cli/` | `libs/cli/langgraph_cli/` (grep `from langgraph_sdk` returns 0 hits) |
| `langgraph-sdk` deps | `httpx`, `orjson`, `langchain-protocol`, `langchain-core`, `websockets` — no `langgraph` runtime dep | `libs/sdk-py/pyproject.toml:14-20` |
| `langgraph-sdk` zero coupling to langgraph core | Only `langgraph_sdk/runtime.py:20` references `langgraph.store.base.BaseStore`, inside `TYPE_CHECKING` | `libs/sdk-py/langgraph_sdk/runtime.py:20` |
| Reverse dependency: core → SDK | `Runtime` imports `BaseUser` from the SDK | `libs/langgraph/langgraph/runtime.py:8` |
| Reverse dependency declared in pyproject | Core declares `langgraph-sdk>=0.4.2,<0.5.0` as a runtime dep | `libs/langgraph/pyproject.toml:29` |
| Explicit `_internal` package marker | `_internal/__init__.py` declares the package is not public | `libs/langgraph/langgraph/_internal/__init__.py:1-4` |
| Public `__all__` discipline (types) | Curated public surface for `langgraph.types` | `libs/langgraph/langgraph/types.py:52-85` |
| Public `__all__` discipline (errors) | Curated public surface for `langgraph.errors` | `libs/langgraph/langgraph/errors.py:16-31` |
| Public `__all__` discipline (channels) | Curated public surface for `langgraph.channels` | `libs/langgraph/langgraph/channels/__init__.py:14-29` |
| Public `__all__` discipline (graph) | Curated public surface for `langgraph.graph` | `libs/langgraph/langgraph/graph/__init__.py:5-12` |
| Public `__all__` discipline (pregel) | Curated public surface for `langgraph.pregel` (only `Pregel`, `NodeBuilder`) | `libs/langgraph/langgraph/pregel/__init__.py:1-3` |
| Public `__all__` discipline (stream) | Curated public surface for `langgraph.stream` | `libs/langgraph/langgraph/stream/__init__.py:28-45` |
| Public `__all__` discipline (managed) | Curated public surface for `langgraph.managed` | `libs/langgraph/langgraph/managed/__init__.py:1-3` |
| Public `__all__` discipline (func) | Curated public surface for `langgraph.func` | `libs/langgraph/langgraph/func/__init__.py:56` |
| Public `__all__` discipline (prebuilt) | Curated public surface for `langgraph.prebuilt` | `libs/prebuilt/langgraph/prebuilt/__init__.py:14-23` |
| Public/private constant boundary | `__getattr__` warns when external code touches private constants | `libs/langgraph/langgraph/constants.py:34-62` |
| Public/private constant — redundancy by design | `_TAG_HIDDEN` is duplicated in `_internal/_constants.py` to avoid a `langgraph.constants` cycle | `libs/langgraph/langgraph/_internal/_constants.py:107-108` |
| Backwards-compat shim | `langgraph.utils.config` re-exports private helpers and is marked "to be removed in v1" | `libs/langgraph/langgraph/utils/config.py:1-4`, `libs/langgraph/langgraph/utils/__init__.py:1` |
| Backwards-compat shim (runnable) | `langgraph.utils.runnable` re-exports `RunnableCallable`/`RunnableLike` with the same removal notice | `libs/langgraph/langgraph/utils/runnable.py:1-3` |
| Acknowledged cyclic import | `langgraph.checkpoint.serde.jsonplus._send_from_args` lazy-imports `langgraph.types.Send` inside the function with a `¯\_(ツ)_/¯` comment | `libs/checkpoint/langgraph/checkpoint/serde/jsonplus.py:537-543` |
| `CompiledStateGraph` extends `Pregel` | Cross-module class hierarchy, one direction | `libs/langgraph/langgraph/graph/state.py:1391-1395` |
| `StateGraph` imports `Pregel` | `langgraph.graph.state` reaches up into `langgraph.pregel` | `libs/langgraph/langgraph/graph/state.py:73-75` |
| `langgraph.pregel` does not import `langgraph.graph` | Verified by grep: zero `from langgraph.graph` in pregel package | `libs/langgraph/langgraph/pregel/` (grep) |
| `langgraph.func` uses private pregel modules | Imports `langgraph.pregel._call`, `langgraph.pregel._read`, `langgraph.pregel._write` | `libs/langgraph/langgraph/func/__init__.py:36-45` |
| `langgraph.pregel._*` modules used across package | `_write.ChannelWrite/ChannelWriteEntry` reachable from `graph._branch`, `graph.state`, `graph.message`, `func` | `libs/langgraph/langgraph/graph/_branch.py:32`, `libs/langgraph/langgraph/graph/state.py:74-75`, `libs/langgraph/langgraph/graph/message.py:408` (local), `libs/langgraph/langgraph/func/__init__.py:45` |
| `langgraph-prebuilt` reaches into `_internal` | Imports `_runnable`, `_typing`, `_constants` | `libs/prebuilt/langgraph/prebuilt/chat_agent_executor.py:32-33`, `libs/prebuilt/langgraph/prebuilt/tool_node.py:85-86`, `libs/prebuilt/langgraph/prebuilt/tool_validator.py:26` |
| `langgraph-prebuilt` reaches into `pregel._*` | Imports `langgraph.pregel._tools._tool_call_writer` | `libs/prebuilt/langgraph/prebuilt/tool_node.py:89` |
| `langgraph-checkpoint-postgres` stays inside its `langgraph.*` sub-trees | Only imports `langgraph.checkpoint.*` and `langgraph.store.*` | `libs/checkpoint-postgres/langgraph/checkpoint/postgres/__init__.py:10-35` (grep summary) |
| `langgraph-checkpoint-conformance` stays inside `langgraph.checkpoint.*` | Only imports `langgraph.checkpoint.base` and `langgraph.checkpoint.conformance.*` | `libs/checkpoint-conformance/langgraph/checkpoint/conformance/__init__.py:3-4`, `libs/checkpoint-conformance/langgraph/checkpoint/conformance/capabilities.py:9` |
| No relative imports anywhere in core libs | Verified: `^from \.` returns no production hits in `libs/langgraph/langgraph/`, `libs/checkpoint/langgraph/`, `libs/prebuilt/langgraph/` | (grep) |
| All packages ship `py.typed` | 13 `py.typed` markers across 7 libs | `find libs -name py.typed` |
| Lint rule against relative imports | `langgraph-sdk` enables ruff `TID252` (banned relative imports) | `libs/sdk-py/pyproject.toml:76` |
| Lint rule against `typing.TypedDict` direct use | `langgraph` ruff config bans `typing.TypedDict` in favour of `typing_extensions.TypedDict` | `libs/langgraph/pyproject.toml:99-100` |
| Test module reaches into `_internal` (acceptable) | `tests/test_pregel.py` imports `CONFIG_KEY_NODE_FINISHED`, `ERROR`, `PULL` from `_internal._constants` | `libs/langgraph/tests/test_pregel.py:43` |
| Test module reaches into `_internal._queue` (acceptable) | `tests/test_pregel_async.py` imports `AsyncQueue` | `libs/langgraph/tests/test_pregel_async.py:46` |
| `langgraph-channel` exports concrete channels | `AnyValue`, `LastValue`, `LastValueAfterFinish`, `BinaryOperatorAggregate`, `Topic`, `NamedBarrierValue`, `EphemeralValue`, `DeltaChannel`, `UntrackedValue` | `libs/langgraph/langgraph/channels/__init__.py:14-29` |
| `BaseChannel` is the abstract base | Defines the channel contract with typevars and abstract methods | `libs/langgraph/langgraph/channels/base.py:19-120` |
| `ManagedValue` is a small, isolated extension point | Abstract class with one method, used by `IsLastStep`/`RemainingSteps` | `libs/langgraph/langgraph/managed/base.py:18-31`, `libs/langgraph/langgraph/managed/is_last_step.py:9-24` |
| `langgraph.stream` is internally self-contained | Mostly imports from `langgraph.stream._*` and `langgraph.errors`/`langgraph.types` | `libs/langgraph/langgraph/stream/__init__.py:8-26`, `libs/langgraph/langgraph/stream/transformers.py:15-18` |
| `langgraph.pregel.protocol` is a leaf protocol module | Only imports from `langgraph.types` and `langgraph.typing` | `libs/langgraph/langgraph/pregel/protocol.py:11-20` |
| Tooling policy per library | `AGENTS.md` tells contributors to run `make format/lint/test` per library before pushing | `AGENTS.md:7-17` |

## Answers to Dimension Questions

1. **Are modules cleanly separated?** Mostly yes. Each PyPI wheel is its own top-level sub-package with an explicit, curated `__all__`; private symbols are concentrated under `langgraph/_internal/` with a docstring-declared "not public API" contract. The single notable blurs are: (a) `langgraph.runtime` reaches into `langgraph_sdk.auth.types` for `BaseUser` (`libs/langgraph/langgraph/runtime.py:8`) and that dep is hard-wired into `libs/langgraph/pyproject.toml:29`; (b) `langgraph-prebuilt` reaches into `langgraph._internal._runnable/_constants/_typing` and `langgraph.pregel._tools` (`libs/prebuilt/langgraph/prebuilt/chat_agent_executor.py:32-33`, `libs/prebuilt/langgraph/prebuilt/tool_node.py:85-89`). Both leaks are one-way across package boundaries, not cyclic.

2. **Do dependencies flow in one direction?** Within the core `langgraph` package the direction is one-way: `channels/`/`errors`/`types`/`constants` are leaves, `_internal` is a private utility layer, `pregel/` sits on top of those, and `graph.state`/`func` sit on top of `pregel`. `libs/langgraph/langgraph/graph/state.py:73` imports `Pregel`, but `libs/langgraph/langgraph/pregel/` has no `from langgraph.graph` import (verified by grep). The only acknowledged cycle is between `langgraph-checkpoint` and `langgraph-core`: `libs/checkpoint/langgraph/checkpoint/serde/jsonplus.py:538-539` does `from langgraph.types import Send` inside the function body with a `¯\_(ツ)_/¯` comment, broken by a lazy local import. Across packages: `checkpoint` → `(nothing else)`, `checkpoint-{postgres,sqlite,conformance}` → `checkpoint`, `prebuilt` → `checkpoint` (and socially → `langgraph` via path sources), `cli` → `langgraph-sdk` (declared but unused), `sdk-py` → standalone, `langgraph` core → `checkpoint + sdk + prebuilt`.

3. **Can modules be used independently?** Yes, with caveats. `langgraph-checkpoint` (`libs/checkpoint`) depends only on `langchain-core` and `ormsgpack` and ships its own `BaseCheckpointSaver` interface (`libs/checkpoint/langgraph/checkpoint/base/__init__.py:38-86`); the `_DeltaSnapshot` and serde machinery in `libs/checkpoint/langgraph/checkpoint/serde/types.py` is self-contained and has its own Protocol mirrors of `langgraph.channels.base.BaseChannel` and `Send` (`libs/checkpoint/langgraph/checkpoint/serde/types.py:39-66`). `langgraph-checkpoint-postgres` and `langgraph-checkpoint-sqlite` are each usable with just `langgraph-checkpoint` plus their respective driver (`libs/checkpoint-postgres/pyproject.toml:14-19`, `libs/checkpoint-sqlite/pyproject.toml:14-18`). `langgraph-sdk` is the most independent — it never imports `langgraph.*` at runtime (only `TYPE_CHECKING`). The caveat: `langgraph-prebuilt` is published as a separate wheel, but its tests import `langgraph.pregel._checkpoint`, so consumers who install `langgraph-prebuilt` against a released `langgraph` wheel inherit that social contract.

4. **Are public APIs distinguished from internal ones?** Three layered mechanisms:
   - A physically separate `langgraph/_internal/` package whose `__init__.py:1-4` explicitly states it is not part of the public API.
   - Per-module curated `__all__` lists in every public module (`types.py:52-85`, `errors.py:16-31`, `channels/__init__.py:14-29`, `graph/__init__.py:5-12`, `pregel/__init__.py:1-3`, `stream/__init__.py:28-45`, `managed/__init__.py:1-3`, `prebuilt/__init__.py:14-23`).
   - `langgraph.constants.__getattr__` emits a `LangGraphDeprecatedSinceV10` warning when external code touches underscored constants re-exported from `_internal._constants` (`libs/langgraph/langgraph/constants.py:34-62`), making misuse loud rather than silent.
   - `py.typed` markers ship with every package, so type-checkers will surface leaks statically.

## Architectural Decisions

- **Single namespace, multiple wheels.** All eight PyPI packages register sub-packages under the same `langgraph` Python namespace (e.g. `langgraph.checkpoint`, `langgraph.store`, `langgraph.checkpoint.postgres`, `langgraph.prebuilt`, `langgraph_sdk`). The naming convention gives users one import root even when they only need a subset of features.
- **Layered core.** `langgraph/types.py`, `errors.py`, `constants.py`, `warnings.py`, `config.py`, `runtime.py`, `managed/` are intentionally small leaves. They define data classes, error types, public constants, and the `Runtime`/`ManagedValue` extension points. Everything heavier sits on top.
- **Two-graph duality.** `langgraph.graph.state.StateGraph` is a builder API that returns `CompiledStateGraph`, which **is-a** `Pregel` (`libs/langgraph/langgraph/graph/state.py:1391-1395`). This forces a single inheritance direction (`graph` → `pregel`) and lets users call any `Pregel` method on a compiled graph. The price is that `graph.state` cannot be split out of `langgraph` without also shipping `pregel`.
- **Pregel decomposition.** `langgraph.pregel/__init__.py` exposes only `Pregel` and `NodeBuilder` (`libs/langgraph/langgraph/pregel/__init__.py:1-3`); the runtime is sliced into underscore-prefixed modules (`_algo`, `_call`, `_checkpoint`, `_executor`, `_io`, `_loop`, `_messages`, `_read`, `_remote_run_stream`, `_retry`, `_runner`, `_tools`, `_utils`, `_validate`, `_write`). These private modules are imported by other parts of `langgraph` (e.g. `_write` by `graph.state` and `func`) and by `prebuilt` (`_tools`).
- **Internal namespace package.** `langgraph/_internal/` is reserved for cross-cutting helpers (`_config`, `_constants`, `_fields`, `_future`, `_pydantic`, `_queue`, `_replay`, `_retry`, `_runnable`, `_scratchpad`, `_serde`, `_timeout`, `_typing`) and carries an explicit "not public API" warning at `libs/langgraph/langgraph/_internal/__init__.py:1-4`.
- **Checkpoint protocol mirrors.** `libs/checkpoint/langgraph/checkpoint/serde/types.py:39-66` defines `ChannelProtocol` and `SendProtocol` that "mirror" the langgraph-core classes (`libs/checkpoint/langgraph/checkpoint/serde/types.py:40,60` comments). This lets `langgraph-checkpoint` ship without depending on the core, while still deserializing values produced by the core.
- **Auth lives in the SDK.** `BaseUser` is defined in `langgraph_sdk.auth.types` (`libs/sdk-py/langgraph_sdk/auth/types.py:182`) and re-imported into the core by `langgraph.runtime` (`libs/langgraph/langgraph/runtime.py:8`). The core declares `langgraph-sdk` as a runtime dependency (`libs/langgraph/pyproject.toml:29`).
- **Conformance harness as a separate package.** `langgraph-checkpoint-conformance` (`libs/checkpoint-conformance`) ships only the conformance test suite and depends solely on `langgraph-checkpoint`, so third-party checkpointer authors can run the suite against their own implementation without pulling in the rest of LangGraph.
- **CLI side-steps the SDK.** `langgraph-cli` declares `langgraph-sdk` as a runtime dep (`libs/cli/pyproject.toml:17`) but does not import it; the CLI talks to the LangGraph Server via `httpx` (`libs/cli/langgraph_cli/config.py:12`, `libs/cli/langgraph_cli/host_backend.py:8`). The dependency is currently dormant.
- **Backwards-compat shims.** `langgraph.utils.config` and `langgraph.utils.runnable` are explicitly marked "to be removed in v1" (`libs/langgraph/langgraph/utils/config.py:1`, `libs/langgraph/langgraph/utils/runnable.py:1`) and re-export a subset of `_internal` symbols; `langgraph.constants.__getattr__` is the runtime-side analogue for the constants namespace.
- **Tooling policy per library.** The `AGENTS.md` of the monorepo codifies the per-library workflow (`make format`, `make lint`, `make test`) and gives a textual dependency map (`AGENTS.md:7-50`), making each library behave like its own project inside the monorepo.

## Notable Patterns

- **Lazy deprecation via `__getattr__`.** `langgraph.constants` uses a module-level `__getattr__` (`libs/langgraph/langgraph/constants.py:34-62`) to redirect Send/Interrupt to `langgraph.types` and to issue a `LangGraphDeprecatedSinceV10` warning when external code accesses private constants re-exported from `_internal._constants`. This pattern keeps old import paths working while giving a clear, versioned deprecation signal.
- **Protocol duplication across the checkpoint/core boundary.** `langgraph.checkpoint.serde.types` declares `ChannelProtocol` and `SendProtocol` (`libs/checkpoint/langgraph/checkpoint/serde/types.py:39-66`) that explicitly "mirror" the core classes — a pragmatic way to avoid a hard runtime dependency between the packages.
- **Explicit cycle handling.** The only known cyclic import is documented in source: `libs/checkpoint/langgraph/checkpoint/serde/jsonplus.py:537-543` lazy-imports `langgraph.types.Send` inside the function body and acknowledges the cycle with a comment.
- **`CompiledStateGraph` IS-A `Pregel`.** A single inheritance relationship between two modules lets the builder API return a fully featured runtime object without a separate type hierarchy (`libs/langgraph/langgraph/graph/state.py:1391-1395`).
- **No relative imports.** All internal libraries use absolute `from langgraph.X` imports (no `from .` patterns were found in `libs/langgraph/langgraph/`, `libs/checkpoint/langgraph/`, or `libs/prebuilt/langgraph/`). The `langgraph-sdk` package actively enforces this with ruff's `TID252` rule (`libs/sdk-py/pyproject.toml:76`).
- **`py.typed` everywhere.** Every published package ships a `py.typed` marker (13 files across 7 libraries), so the public/internal split is also enforced at type-check time.

## Tradeoffs

- **Reverse-direction coupling.** Putting `BaseUser` in the SDK and importing it from the core (`libs/langgraph/langgraph/runtime.py:8`) means the core cannot be installed without the SDK, and any change to the `BaseUser` Protocol must be coordinated across both repos. The README explicitly lists the SDK as part of the ecosystem, so the practical impact is small, but the naming is confusing.
- **Underscore-prefix is social, not enforced.** `langgraph-prebuilt` reaches into `langgraph._internal._runnable`/`_constants`/`_typing` (`libs/prebuilt/langgraph/prebuilt/chat_agent_executor.py:32-33`, `libs/prebuilt/langgraph/prebuilt/tool_node.py:85-86`) and into `langgraph.pregel._tools` (`libs/prebuilt/langgraph/prebuilt/tool_node.py:89`). This works because the two packages are developed in lock-step in the same monorepo, but it makes `prebuilt` fragile if the underscore modules change shape.
- **`pregel._*` is leaky inside the core.** Modules like `langgraph.pregel._write.ChannelWrite/ChannelWriteEntry` (`libs/langgraph/langgraph/pregel/_write.py:26-46`) are imported by `langgraph.graph._branch`, `langgraph.graph.state`, `langgraph.graph.message`, and `langgraph.func` (`libs/langgraph/langgraph/graph/_branch.py:32`, `libs/langgraph/langgraph/graph/state.py:74-75`, `libs/langgraph/langgraph/graph/message.py:408`, `libs/langgraph/langgraph/func/__init__.py:45`). If `pregel._write` is refactored, several unrelated-looking modules break in lock-step.
- **Dormant dependency.** `langgraph-cli` declares `langgraph-sdk>=0.1.0` (`libs/cli/pyproject.toml:17`) but does not import it (verified with grep). This costs an extra install at user time for no runtime use today.
- **Internal cycle requires discipline.** The `langgraph.types.Send` ↔ `langgraph.checkpoint.serde.jsonplus` cycle (`libs/checkpoint/langgraph/checkpoint/serde/jsonplus.py:538-539`) is broken with a function-local import. It works, but it relies on every call site happening after both modules are fully initialised.
- **`langgraph.utils.*` is a transitional surface.** Re-exporting `_internal` symbols through `langgraph.utils.config` and `langgraph.utils.runnable` (`libs/langgraph/langgraph/utils/config.py:1-4`, `libs/langgraph/langgraph/utils/runnable.py:1-3`) doubles the surface area users could anchor code on and must be removed at v1.

## Failure Modes / Edge Cases

- **Hard coupling if the SDK `BaseUser` Protocol changes.** Adding or removing a field on `langgraph_sdk.auth.types.BaseUser` (`libs/sdk-py/langgraph_sdk/auth/types.py:182`) will silently affect `Runtime.user` typing in `langgraph.runtime` (`libs/langgraph/langgraph/runtime.py:8`), even though the two packages are released independently.
- **`langgraph-prebuilt` is not as decoupled as the package split suggests.** If `_internal._runnable.RunnableCallable` (`libs/langgraph/langgraph/_internal/_runnable.py`) changes its signature, every `langgraph-prebuilt` release that touches it (`libs/prebuilt/langgraph/prebuilt/chat_agent_executor.py:32`, `libs/prebuilt/langgraph/prebuilt/tool_node.py:86`, `libs/prebuilt/langgraph/prebuilt/tool_validator.py:26`) has to be re-released in lock-step.
- **Cycle reintroduction risk.** The lazy `from langgraph.types import Send` inside `_send_from_args` (`libs/checkpoint/langgraph/checkpoint/serde/jsonplus.py:539`) is the only thing standing between `langgraph-checkpoint` and a hard import-time cycle. Moving the call to module scope, or moving `Send` itself into `_internal`, would break both packages.
- **`CompiledStateGraph` ↔ `Pregel` is a contract.** `CompiledStateGraph` subclasses `Pregel` (`libs/langgraph/langgraph/graph/state.py:1391-1395`) and stores the `StateGraph` builder on `self.builder`. Any change to `Pregel.__init__` or its public surface ripples into every compile() call site.
- **Dormant CLI dependency can confuse static analysis.** `langgraph-cli` lists `langgraph-sdk` in `pyproject.toml:17` but never imports it; lockfiles may pull a version that is incompatible with the SDK used by the core, leading to install-time friction without runtime benefit.
- **Tests reach into private modules.** In-tree tests for `langgraph` and `prebuilt` import `langgraph._internal._config/_constants/_typing/_pydantic/_queue` and `langgraph.pregel._checkpoint` (`libs/langgraph/tests/test_pregel.py:43`, `libs/langgraph/tests/test_pregel_async.py:45-46`, `libs/prebuilt/tests/memory_assert.py:13`). These are intentional white-box tests, but a consumer copy-pasting the pattern would be relying on private surface.

## Future Considerations

- **Resolve the core/SDK auth coupling.** Either move `BaseUser` into a neutral `langgraph.types` or `langgraph._internal` and have the SDK re-export it, or split `Runtime.user` into an optional `BaseUser | None` typed as `Any` so the core no longer hard-imports the SDK.
- **Promote `langgraph._internal` symbols that `prebuilt` actually needs.** Items like `RunnableCallable`, `RunnableLike`, `is_async_callable`, `MISSING`, `CONF`, `CONFIG_KEY_READ` are now part of `prebuilt`'s de-facto public surface. Moving them into `langgraph.runnable` / `langgraph.constants` would let the underscore modules be refactored freely.
- **Replace the `_send_from_args` cycle with a shared `_internal` symbol.** Pulling `Send` into `langgraph._internal._types` (or another leaf shared by both packages) would let `jsonplus.py` import it at module top-level.
- **Clean up `langgraph.utils.*`.** The two shim modules (`libs/langgraph/langgraph/utils/config.py`, `libs/langgraph/langgraph/utils/runnable.py`) are marked "to be removed in v1" and should be retired in lock-step with a deprecation cycle.
- **Decide what to do about `langgraph-cli`'s dormant `langgraph-sdk` dep.** Either use it (replace ad-hoc `httpx` calls with the SDK client) or drop it from `pyproject.toml` to make the install footprint honest.
- **Strengthen type-level enforcement of `langgraph.constants.__getattr__`.** Today the deprecation warning fires at runtime; a `py.typed` re-export list plus a type-checker rule could prevent the warning from firing in the first place.
- **Document the `pregel._*` social contract.** Since `langgraph.pregel._write`, `_read`, `_call`, `_tools`, `_messages` are routinely imported by other in-tree modules, a code comment or `__init__.py` note acknowledging them as "intra-package private but cross-module shared" would set expectations.

## Questions / Gaps

- The exact mechanism used by the `langgraph-cli` ↔ `langgraph-sdk` integration is not observable in the source tree — `langgraph_cli/` has zero `langgraph_sdk` imports despite the declared dep (`libs/cli/pyproject.toml:17`). Possible explanations (intentional phasing, dead-code to be removed, hidden runtime path) cannot be confirmed without release notes.
- It is unclear whether `libs/checkpoint/langgraph/checkpoint/serde/jsonplus.py:539`'s lazy import is covered by an import-time test; the cycle risk is mitigated by code review only, with no recorded regression test in the visible tests folder (`libs/checkpoint/tests/` not enumerated here).
- No CI configuration is visible from inside the repository that would mechanically verify the absence of cross-package `_internal` imports (e.g. an import-linter rule or a ruff custom check); the boundary is enforced socially.
- The behaviour of `langgraph.constants.__getattr__` when a user imports a non-deprecated attribute that shadows a private one (`libs/langgraph/langgraph/constants.py:51-60`) is partially documented in the warning string but not exercised by an in-tree test in the visible scope.

---

Generated by `studies/agent-harness-study/.opencode/skills/ultraplan/dimensions/22.01-package-module-boundaries` against `langgraph`.