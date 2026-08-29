# Source Analysis: langgraph

## Plugin and Extension Points

### Source Info

| Field | Value |
|-------|-------|
| Name | langgraph |
| Path | `studies/agent-harness-study/sources/langgraph` |
| Language / Stack | Python (monorepo: libs/langgraph, libs/checkpoint, libs/prebuilt, libs/cli, libs/sdk-py) |
| Analyzed | 2026-08-28 |

## Summary

LangGraph has no unified “plugin” registry or dynamic loader. Extensibility is achieved through **inheritance-based interfaces** and **composition at compile-time**: users subclass/implement `BaseChannel`, `BaseCheckpointSaver`, `SerializerProtocol`, `BaseStore`, `BaseCache`, `ManagedValue`, `StreamTransformer`, and `BaseTool`/`Runnable` nodes, then pass instances to `StateGraph.compile()`, `Pregel()`, or `entrypoint`. Tools are first-class extension via `langchain_core.tools.BaseTool` and `ToolNode` injection (`InjectedState`, `InjectedStore`, `ToolRuntime`). The only runtime dynamic loading is via the CLI’s `langgraph.json` config, which resolves dotted-path strings (`graphs: "./module.py:attr"`, `auth.path`, `store`, `checkpointer`) with `importlib.import_module` and via `JsonPlusSerializer` deserialization (`libs/checkpoint/langgraph/checkpoint/serde/jsonplus.py:193`). There is no plugin lifecycle manager, no versioned registry, no isolation/sandboxing, and no entry-point discovery mechanism.

## Rating

**5 / 10 — Present but inconsistent, weakly documented, fragile**

Rationale: Core extension interfaces are explicit, typed, and tested (channels, checkpoint, store, transformer), but they are not unified under a plugin lifecycle. Loading is static (pass object to builder) except for the CLI deployment config; there is no discovery, enable/disable, or dependency model. Isolation is absent — all extensions share the same process/memory and channel namespace. Documentation is per-interface docstrings, not a stable plugin contract; many interfaces lack guarantees about backward compatibility or capability negotiation.

## Evidence Collected

Every entry MUST include a file path with line numbers. Format: `path/to/file.ts:NN`.

| Area | Evidence | File:Line |
|------|----------|-----------|
| Extension interface — channels | `BaseChannel` ABC with `get`, `update`, `checkpoint`, `from_checkpoint`, `ValueType`/`UpdateType` | `libs/langgraph/langgraph/channels/base.py:19` |
| Channel implementations | 8 built-ins: `LastValue`, `Topic`, `BinaryOperatorAggregate`, `DeltaChannel`, `EphemeralValue`, `AnyValue`, `NamedBarrierValue`, `UntrackedValue` | `libs/langgraph/langgraph/channels/__init__.py:1-29` |
| Channel — DeltaChannel | Beta delta-channel backed by `BaseCheckpointSaver.get_delta_channel_history` | `libs/langgraph/langgraph/channels/delta.py:25` |
| Extension interface — checkpoint | `BaseCheckpointSaver` ABC with `get_tuple`, `put`, `put_writes`, `list`, async variants, `get_delta_channel_history` | `libs/checkpoint/langgraph/checkpoint/base/__init__.py:176` |
| Checkpoint — InMemory impl | Reference impl `InMemorySaver(BaseCheckpointSaver[str])` with serde, blob storage, parent-chain traversal | `libs/checkpoint/langgraph/checkpoint/memory/__init__.py:33` |
| Checkpoint — serde | `SerializerProtocol` (`dumps_typed`/`loads_typed`) + `UntypedSerializerProtocol` | `libs/checkpoint/langgraph/checkpoint/serde/base.py:15` |
| Checkpoint — serde implementation | `JsonPlusSerializer` with `importlib.import_module` allowlist deserialization | `libs/checkpoint/langgraph/checkpoint/serde/jsonplus.py:193,651` |
| Extension interface — store | `BaseStore` ABC with `batch`/`abatch`, `get`/`put`/`search`/`list_namespaces`, TTL/index config | `libs/checkpoint/langgraph/store/base/__init__.py:700` |
| Store — InMemory impl | `InMemoryStore` dictionary-backed with optional vector search | `libs/checkpoint/langgraph/store/memory/__init__.py:1` |
| Extension interface — cache | `BaseCache` ABC with `get`/`aget`/`set`/`aset`/`clear`/`aclear` | `libs/checkpoint/langgraph/cache/base/__init__.py:15` |
| Extension interface — managed values | `ManagedValue` ABC + `ManagedValueSpec = type[ManagedValue]`, `IsLastStep`, `RemainingSteps` | `libs/langgraph/langgraph/managed/base.py:18` |
| Extension interface — stream | `StreamTransformer` ABC with `init`, `process`/`aprocess`, `finalize`/`afinalize`, `fail`/`afail`, `schedule`, `scope`, `requires_async`, `required_stream_modes` | `libs/langgraph/langgraph/stream/_types.py:44` |
| Stream — transformers | Built-ins: `LifecycleTransformer`, `MessagesTransformer`, `SubgraphTransformer`, `ValuesTransformer` used by Pregel | `libs/langgraph/langgraph/stream/transformers.py:1` |
| Stream — channel | `StreamChannel` generic typed queue for transformer projections | `libs/langgraph/langgraph/stream/stream_channel.py:14` |
| Graph compilation — transformer injection | `StateGraph.compile(..., transformers: Sequence[Callable[[tuple[str,...]],Any]])` normalized to scope-aware factories | `libs/langgraph/langgraph/graph/state.py:1174` |
| Pregel — node definition | `PregelNode` + `NodeBuilder` fluent builder (`subscribe_only`, `subscribe_to`, `do`, `write_to`, `meta`, `add_retry_policies`) | `libs/langgraph/langgraph/pregel/main.py:205` |
| Pregel — protocol | `PregelProtocol(Runnable)` defines `invoke`/`stream`/`get_state`/`get_graph` contracts | `libs/langgraph/langgraph/pregel/protocol.py:25` |
| Functional API | `@entrypoint` + `@task` decorators produce `Pregel` nodes via `get_runnable_for_entrypoint` | `libs/langgraph/langgraph/func/__init__.py:262` |
| Tools — ToolNode | `ToolNode(RunnableCallable)` with `handle_tool_errors`, `messages_key`, `InjectedState`/`InjectedStore` injection, parallel execution | `libs/prebuilt/langgraph/prebuilt/tool_node.py:622` |
| Tools — injection types | `InjectedState`, `InjectedStore`, `ToolRuntime(State, context, config, tool_call_id, store, writer)` | `libs/prebuilt/langgraph/prebuilt/tool_node.py:1643,1817` |
| Tools — prebuilt agent | `create_react_agent(model, tools, prompt, checkpointer, store, interrupt_before)` composes `StateGraph` + `ToolNode` | `libs/prebuilt/langgraph/prebuilt/chat_agent_executor.py:1` |
| Tool condition helper | `tools_condition` routing helper for `ToolNode` | `libs/prebuilt/langgraph/prebuilt/__init__.py:8` |
| Validation extension | `ValidationNode` for pydantic schema validation of tool calls | `libs/prebuilt/langgraph/prebuilt/tool_validator.py:1` |
| CLI config — graph loading | `python_config_to_docker` + `_assemble_local_deps` + `_update_graph_paths` rewrites `graphs` strings to container paths | `libs/cli/langgraph_cli/config.py:827,709,1263` |
| CLI config — deployment extension points | `graphs`, `store`, `auth.path`, `encryption.path`, `checkpointer.path`, `http.app`, `ui` resolved from `langgraph.json` | `libs/cli/langgraph_cli/config.py:368,564,927,968,1011,1056` |
| CLI config — validation | `validate_config` enforces `python_version`, `dependencies`, `image_distro`, `source.kind` | `libs/cli/langgraph_cli/config.py:323` |
| Pregel runtime objects | `Runtime[ContextT]` injects `context`, `store`, `stream_writer`, `previous` into nodes | `libs/langgraph/langgraph/runtime.py:1` |
| Constants — START/END | Reserved channel keys `START="__start__"`, `END="__end__"` used as entry/exit sentinels | `libs/langgraph/langgraph/constants.py:28` |
| Dynamic import — deprecation shim | `langgraph.constants.__getattr__` uses `importlib.import_module` to lazily import private constants | `libs/langgraph/langgraph/constants.py:43` |
| Tests — StreamTransformer lifecycle | `_TwoChannelTransformer(StreamTransformer)`, `ChannelPusher`, `TestPregelStreamEventsV3` verify `init/process/finalize` | `libs/langgraph/tests/test_interleave_arrival_order.py:24`, `libs/langgraph/tests/test_pregel_stream_events_v3.py:1160` |
| Docs — prebuilt README | Describes `create_react_agent`, `ToolNode`, `ValidationNode` as public extension surface | `libs/prebuilt/README.md:19` |

## Answers to Dimension Questions

### 1. What can be extended via plugins?

**No formal plugin type, but 9 inherited interfaces allow extension without core edits:**

- **Channels (state)**: subclass `BaseChannel` (`libs/langgraph/langgraph/channels/base.py:19`). Built-ins in `libs/langgraph/langgraph/channels/__init__.py:1`; custom channels passed via `StateGraph._add_schema` (`libs/langgraph/langgraph/graph/state.py:342`) and `Pregel(channels=...)` (`libs/langgraph/langgraph/pregel/main.py:757`). `DeltaChannel` is beta and requires saver support (`libs/langgraph/langgraph/channels/delta.py:25`).
- **Checkpoint persistence**: subclass `BaseCheckpointSaver` (`libs/checkpoint/langgraph/checkpoint/base/__init__.py:176`). `InMemorySaver`, `PostgresSaver`, `SqliteSaver` are reference impls. Serializer pluggable via `SerializerProtocol` (`libs/checkpoint/langgraph/checkpoint/serde/base.py:15`) — default `JsonPlusSerializer` swappable at `BaseCheckpointSaver.serde` and `BaseCache.serde`.
- **Cross-thread memory**: subclass `BaseStore` (`libs/checkpoint/langgraph/store/base/__init__.py:700`) — methods `batch`/`abatch`, `get`/`put`/`search`. Includes TTL and vector-index (`IndexConfig`/`TTLConfig`).
- **Caching**: subclass `BaseCache` (`libs/checkpoint/langgraph/cache/base/__init__.py:15`) — `get`/`set`/`clear` sync+async.
- **Managed values**: subclass `ManagedValue` (`libs/langgraph/langgraph/managed/base.py:18`) — `IsLastStep`/`RemainingSteps` are the two shipped; others injected via `channels` dict as `ManagedValueSpec`.
- **Streaming projections**: subclass `StreamTransformer` (`libs/langgraph/langgraph/stream/_types.py:44`) — declare `required_stream_modes`, `before_builtins`, `requires_async`/`supports_sync`; instantiated per-run via `StateGraph.compile(transformers=...)` (`libs/langgraph/langgraph/graph/state.py:1174`) and `Pregel(stream_transformers=...)` (`libs/langgraph/langgraph/pregel/main.py:832`).
- **Agent logic**: any `Runnable`/`callable` added via `StateGraph.add_node` (`libs/langgraph/langgraph/graph/state.py:662`) or `NodeBuilder.do` (`libs/langgraph/langgraph/pregel/main.py:303`) or `@task`/`@entrypoint` (`libs/langgraph/langgraph/func/__init__.py:59,262`). Tools are `langchain_core.tools.BaseTool` executed by `ToolNode` (`libs/prebuilt/langgraph/prebuilt/tool_node.py:622`) with injection annotations `InjectedState`/`InjectedStore`/`ToolRuntime`.
- **Prebuilt behaviors**: `create_react_agent` (`libs/prebuilt/langgraph/prebuilt/chat_agent_executor.py:1`) is itself a composition point — model, tools, prompt, `checkpointer`/`store`/`interrupt_before` are all injectable. `ValidationNode` (`libs/prebuilt/langgraph/prebuilt/tool_validator.py:1`) is an alternative tool-validation extension.
- **Deployment-level**: `langgraph.json` `auth`, `encryption`, `http.app`, `store`, `checkpointer`, `ui` each point to an importable Python object (`"<module>:<attr>"`) validated in `libs/cli/langgraph_cli/config.py:521,530,537,545`.

**Cannot be extended without core change:** the Pregel execution loop (`libs/langgraph/langgraph/pregel/_loop.py:1`, `_runner.py`, `_algo.py`) is not pluggable; stream event types (`libs/langgraph/langgraph/types.py:120`) and protocol envelopes (`libs/langgraph/langgraph/stream/_types.py:14`) are closed unions.

### 2. Can plugins be loaded at runtime?

**Partially, but not as a plugin system:**

- **Compile-time wiring (preferred):** Most extensions are instantiated in user code and passed to `StateGraph.compile(checkpointer=..., store=..., cache=..., transformers=...)` (`libs/langgraph/langgraph/graph/state.py:1164`) or `Pregel(...)` (`libs/langgraph/langgraph/pregel/main.py:757`). This is static for the process lifetime — new channels/checkpointers require recompilation, but graphs are lightweight to rebuild.
- **CLI deployment dynamic import:** For deployment, `langgraph.json` strings like `{"graphs": {"agent": "./src/agent.py:graph"}}` and `{"auth": {"path": "./src/auth.py:handle_auth"}}` are validated (`libs/cli/langgraph_cli/config.py:323`) and rewritten to container paths (`libs/cli/langgraph_cli/config.py:827,927`). At runtime the LangGraph API server dynamically imports them via `importlib.import_module` (indirectly through the generated Dockerfile env `LANGSERVE_GRAPHS`, `LANGGRAPH_AUTH`, etc. — `libs/cli/langgraph_cli/config.py:1132`). This is the closest to runtime plugin loading.
- **Serde dynamic import:** `JsonPlusSerializer` reconstructs arbitrary types by `importlib.import_module(tup[0])` (`libs/checkpoint/langgraph/checkpoint/serde/jsonplus.py:651,664`), but that is deserialization, not a plugin loader, and is gated by a msgpack allowlist built in `StateGraph.compile` (`libs/langgraph/langgraph/graph/state.py:1221,1235`).
- **No entry-point discovery:** No `importlib.metadata.entry_points`, no `setuptools` plugin groups, no file-system watch or hot reload; `grep` across `libs/langgraph` and `libs/checkpoint` shows only incidental `importlib` uses (`libs/langgraph/langgraph/constants.py:43`, `libs/checkpoint/langgraph/checkpoint/serde/jsonplus.py:6`), not a registry.
- **API surface gap:** `langgraph_api` (not in this source snapshot) is what hosts the dynamically imported graphs; the core library itself has no `load_plugin(path)` API.

### 3. Are plugins isolated from each other?

**No.**

- **Same process, shared namespace:** All nodes, channels, and transformers share the in-process `Pregel` graph (`libs/langgraph/langgraph/pregel/main.py:449`) and a flat channel dict (`Pregel.channels: dict[str, BaseChannel|ManagedValueSpec]` at `libs/langgraph/langgraph/pregel/main.py:706`). A custom `BaseChannel` can read/write any channel name it knows; there is no capability scoping.
- **No sandboxing or privilege model:** `ToolNode` executes tools directly (`libs/prebuilt/langgraph/prebuilt/tool_node.py:758` forward) with full access to `Runtime` (`store`, `config`, `stream_writer`, `context`). Store namespaces reserve `"langgraph"` (`libs/checkpoint/langgraph/store/base/__init__.py:1272`) but that is just a label check, not isolation.
- **Resource sharing, not bulkheading:** Checkpointer `blobs`, `writes`, `storage` are shared structures (`libs/checkpoint/langgraph/checkpoint/memory/__init__.py:68`); a misbehaving channel or checkpointer saver can poison checkpoints for all plugins on that thread. `DeltaChannel` caveats explicitly warn that deleting intermediate checkpoints breaks reconstruction (`libs/checkpoint/langgraph/checkpoint/base/__init__.py:333,380`).
- **Transformer ordering hazards:** `StreamTransformer.before_builtins` flag (`libs/langgraph/langgraph/stream/_types.py:98`) and `_normalize_stream_transformer_factories` (`libs/langgraph/langgraph/pregel/main.py:416`) reveal that transformers share an ordered pipeline and can desync built-ins if they mutate events — documented as a foot-gun, not prevented.
- **Failure propagation:** Errors bubble via `GraphBubbleUp` (`libs/prebuilt/langgraph/prebuilt/tool_node.py:87`) and cancel the super-step; an isolated failure domain (per-plugin try/catch) is not provided — it is per-node retry policy (`RetryPolicy` in `libs/langgraph/langgraph/types.py:415`).

### 4. Are extension points documented and stable?

**Mixed — some mature, others beta/inconsistent:**

- **Well-documented:** `BaseChannel` (`libs/langgraph/langgraph/channels/base.py:19` docstring + methods), `BaseCheckpointSaver` (module docstring at `libs/checkpoint/langgraph/checkpoint/base/__init__.py:176` plus per-method signatures), `BaseStore` (rich docstrings at `libs/checkpoint/langgraph/store/base/__init__.py:700` with examples), `BaseCache` (`libs/checkpoint/langgraph/cache/base/__init__.py:15`), `StreamTransformer` (100+ lines of class docstring at `libs/langgraph/langgraph/stream/_types.py:44` covering `process`/`aprocess`, `schedule`, `requires_async`, `before_builtins`). `ToolNode` has module-level design patterns (`libs/prebuilt/langgraph/prebuilt/tool_node.py:1`). `create_react_agent` args are typed (`libs/prebuilt/langgraph/prebuilt/chat_agent_executor.py:1`).
- **Weakly documented:** `ManagedValue` is 31 lines total (`libs/langgraph/langgraph/managed/base.py:18`) with no usage guide beyond `IsLastStep`. `SerializerProtocol` is 64 lines (`libs/checkpoint/langgraph/checkpoint/serde/base.py:6`) without stability guarantees.
- **Explicitly unstable:** `DeltaChannel` and `BaseCheckpointSaver.get_delta_channel_history` are marked `!!! warning "Beta"` (`libs/checkpoint/langgraph/checkpoint/base/__init__.py:588`) with notes that field names and semantics may change. The serde allowlist plumbing (`_serde.STRICT_MSGPACK_ENABLED`, `build_serde_allowlist` at `libs/langgraph/langgraph/graph/state.py:1221`) is internal and not versioned for external plugins.
- **Stability mechanisms:** Semver-tracked `Checkpoint.v` and `_migrate_checkpoint` (`libs/langgraph/langgraph/pregel/main.py:1135`) handle checkpoint format evolution, but there is no deprecation policy for plugin interfaces. Tests serve as the de facto contract: `test_channels.py`, `test_stream_events_v3.py`, `test_checkpoint_migration.py`, `test_tool_node_validation_error_filtering.py`, `test_interleave_arrival_order.py` pin transformer and channel behavior.
- **Missing centralized contract:** No `EXTENSIONS.md` or plugin SDK docs in this repo’s `docs/` (only `docs/.gitignore|redirects.json` at `docs`). Public reference docs live outside the repo (pointed to in `libs/prebuilt/README.md:19`).

## Architectural Decisions

| Decision | Location | Tradeoff |
|----------|----------|----------|
| Inheritance over registry: all extension points are ABCs/Protocols, not string-keyed registries | `libs/langgraph/langgraph/channels/base.py:19`, `libs/checkpoint/langgraph/checkpoint/base/__init__.py:176`, `libs/checkpoint/langgraph/store/base/__init__.py:700`, `libs/checkpoint/langgraph/cache/base/__init__.py:15`, `libs/langgraph/langgraph/stream/_types.py:44` | Maximum flexibility and static typing, but no discovery lifecycle; users must wire objects manually. |
| Channels as versioned, checkpointable state | `BaseChannel.checkpoint/from_checkpoint` at `libs/langgraph/langgraph/channels/base.py:49,61` | Enables time-travel (`get_state_history`) but makes custom channels responsible for serializability and migration. |
| Single flat channel namespace | `Pregel.channels: dict[str, ...]` at `libs/langgraph/langgraph/pregel/main.py:706`, `StateGraph.channels` at `libs/langgraph/langgraph/graph/state.py:204` | Simple routing, collision-prone; no per-plugin isolation. |
| CLI deployment as thin import shim | `langgraph.json` `graphs`/`auth`/`store`/`checkpointer` + `_update_graph_paths` at `libs/cli/langgraph_cli/config.py:827` | Achieves runtime loading without a plugin daemon; couples deployment to file layout and `importlib.import_module`. |
| Serializer allowlist enforced at compile time | `StateGraph.compile` builds `serde_allowlist` from schemas+channels and clones checkpointer via `with_allowlist` at `libs/langgraph/langgraph/graph/state.py:1221,1239` | Hardens `importlib` deserialization against arbitrary type injection, but adds compile-time coupling between graph shape and checkpointer. |
| StreamTransformer scope-aware factories | `_normalize_stream_transformer_factories` rejects pre-built instances at `libs/langgraph/langgraph/pregel/main.py:416` | Ensures per-subgraph scoping, but prevents singleton transformers and increases boilerplate. |

## Notable Patterns

- **Builder + compile:** `StateGraph` builder (add nodes/edges/conditional branches via `BranchSpec` at `libs/langgraph/langgraph/graph/_branch.py:1`) compiles to immutable `Pregel` (`libs/langgraph/langgraph/pregel/main.py:757`); extensions injected at build time.
- **Annotation injection:** `InjectedState`/`InjectedStore` use `Annotated` + `BaseTool` introspection to inject `Runtime` facets into tool signatures (`libs/prebuilt/langgraph/prebuilt/tool_node.py:1643,1817`).
- **Multiplexed streaming:** `StreamMux` collects `required_stream_modes` from all transformers (`libs/langgraph/langgraph/pregel/main.py:397`) and fans out `ProtocolEvent`s to each; `StreamChannel` projections are auto-wired (`libs/langgraph/langgraph/stream/_types.py:44`).
- **Decorators as sugar:** `@task` + `@entrypoint` lower to `PregelNode`/`RunnableCallable` (`libs/langgraph/langgraph/func/__init__.py:59,262`) — functional API without touching `Pregel`.

## Tradeoffs

- **Flexibility vs. safety:** Arbitrary `BaseChannel`/`BaseCheckpointSaver` implementations give power users full control but allow breaking invariants (e.g., not implementing `from_checkpoint` causes history walk failure — `DeltaChannel` warns at `libs/checkpoint/langgraph/checkpoint/base/__init__.py:380`).
- **Explicit wiring vs. discovery:** No plugin discovery means simpler mental model and deterministic startup, but prevents ecosystem of installable plugins (contrast with LangChain’s `entry_points` for tools).
- **Isolation vs. performance:** In-process execution with shared `Checkpoint`/`Store` avoids serialization overhead but means a buggy tool or channel can corrupt shared state; the beta `DeltaChannel` amplifies this (`libs/checkpoint/langgraph/checkpoint/base/__init__.py:333`).
- **Compile-time allowlist vs. openness:** Strict msgpack allowlist (`libs/langgraph/langgraph/graph/state.py:1221`) prevents deserialization exploits but requires recompilation when adding new channel types.
- **Transformer power vs. foot-guns:** `before_builtins` (`libs/langgraph/langgraph/stream/_types.py:98`) and `schedule()` (`libs/langgraph/langgraph/stream/_types.py:233`) enable advanced projections (PII redaction, moderation) but come with documented desync risks.

## Failure Modes / Edge Cases

| Mode | Location | Impact |
|------|----------|--------|
| Custom `BaseChannel` does not implement `from_checkpoint` correctly | `libs/langgraph/langgraph/channels/base.py:61` contract | Checkpoint load returns wrong values; delta reconstruction silently yields empty channel (`libs/checkpoint/langgraph/checkpoint/base/__init__.py:597` “`seed` omitted → start empty”). |
| `DeltaChannel` without saver history walk | `libs/checkpoint/langgraph/checkpoint/base/__init__.py:583` | `prune`/`copy_thread`/`delete_for_runs` break history chain — no error raised, silent data loss (caveats at `libs/checkpoint/langgraph/checkpoint/base/__init__.py:340,355,387`). |
| `StreamTransformer` installed as instance not factory | `libs/langgraph/langgraph/pregel/main.py:429` `TypeError` | Fails at `stream_events(version="v3")` startup; no recovery without fixing `StateGraph.compile(transformers=...)`. |
| Async transformer under sync `stream()` | `libs/langgraph/langgraph/stream/_types.py:308` `transformer_requires_async` + `libs/langgraph/langgraph/pregel/main.py:416` check | Registration raises `TypeError`; `schedule()` without event loop raises `RuntimeError` (`libs/langgraph/langgraph/stream/_types.py:270`). |
| `ToolNode` tool not registered | `libs/prebuilt/langgraph/prebuilt/tool_node.py:949` `Error: {requested_tool} is not a valid tool` | Falls back to error template (`TOOL_CALL_ERROR_TEMPLATE` at `libs/prebuilt/langgraph/prebuilt/tool_node.py:111`) or raises `ToolException` if `handle_tool_errors=False`. |
| Serializer allowlist mismatch | `libs/langgraph/langgraph/graph/state.py:1239` `apply_checkpointer_allowlist` | Deserialization raises allowlist violation; mitigated by rebuilding graph, not hotspot-patchable. |
| CLI auth/encryption path outside `dependencies` | `libs/cli/langgraph_cli/config.py:961,1004` `ValueError` | Docker build fails; requires adding parent dir to `dependencies` array. |
| Store namespace `langgraph` prefix | `libs/checkpoint/langgraph/store/base/__init__.py:1272` `InvalidNamespaceError` | Rejected at `put`/`aput` time — no fallback. |

## Future Considerations

- **Plugin registry:** Add `importlib.metadata.entry_points(group="langgraph.channels")` discovery with versioned interface checks; would enable `pip install langgraph-ext-foo` without editing graph code.
- **Lifecycle hooks:** Introduce `on_install`/`on_uninstall`/`on_enable` and health checks for checkpointers/stores; current `BaseCheckpointSaver` has no init/dispose contract beyond serializer injection (`libs/checkpoint/langgraph/checkpoint/base/__init__.py:211`).
- **Isolation:** Namespace-prefix scoping for channels + store, plus per-plugin resource quotas/timeouts (`TimeoutPolicy` exists only per-node at `libs/langgraph/langgraph/types.py:449`). Could sandbox `ToolNode` tools via `BaseStore` proxy.
- **Stabilize delta:** Promote `DeltaChannel` and `get_delta_channel_history` out of beta (`libs/checkpoint/langgraph/checkpoint/base/__init__.py:588`) and add migration tooling for `prune`.
- **Docs contract:** Publish a versioned Extension API page listing `BaseChannel`, `BaseCheckpointSaver`, `SerializerProtocol`, `BaseStore`, `BaseCache`, `ManagedValue`, `StreamTransformer` with semver guarantees — today only per-class docstrings exist.

## Questions / Gaps

- **Evals/prompts/policies/UI extension points** (asked in dimension purpose) — no evidence in core libs; `eval`, `prompt`, `policy`, `UI` extension surfaces may live in `langgraph_api` server, not in this source snapshot. Searched `libs/langgraph`, `libs/checkpoint`, `libs/prebuilt`, `libs/cli` — no `BaseEvaluator`/`BasePrompt`/`Policy` ABCs found.
- **Plugin version negotiation** — no capability flags or `apiVersion` field on checkpointers/stores to negotiate protocol versions.
- **Dynamic reload in dev** — `TEST=` make targets (`AGENTS.md:10`) suggest static test invocation; no watched plugin reload observed.
- **LangChain tool interop stability** — `ToolNode` re-exports/extends `langchain_core.tools.BaseTool` (`libs/prebuilt/langgraph/prebuilt/tool_node.py:76`); stability depends on external `langchain-core` version (`libs/langgraph/pyproject.toml:27` pins `>=1.4.7,<2`) — not modeled as a LangGraph extension guarantee.

---

Generated by `21.01-plugin-and-extension-points` against `langgraph`.
