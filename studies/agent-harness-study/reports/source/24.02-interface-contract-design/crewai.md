# Source Analysis: crewai

## Interface Contract Design

### Source Info

| Field | Value |
|-------|-------|
| Name | crewai |
| Path | `studies/agent-harness-study/sources/crewai` |
| Language / Stack | Python 3 (Pydantic v2, ABC + `typing.Protocol`, mypy plugin, pytest) |
| Analyzed | 2026-08-22 |

## Summary

CrewAI defines its cross-boundary contracts through two complementary mechanisms: **Pydantic-backed abstract base classes** for extension points it owns (tools `lib/crewai/src/crewai/tools/base_tool.py:103`, LLMs `lib/crewai/src/crewai/llms/base_llm.py:129`, flow persistence `lib/crewai/src/crewai/flow/persistence/base.py:18`, knowledge storage `lib/crewai/src/crewai/knowledge/storage/base_knowledge_storage.py:13`, state providers `lib/crewai/src/crewai/state/provider/core.py:10`, MCP transports `lib/crewai/src/crewai/mcp/transports/base.py:25`, agent adapters `lib/crewai/src/crewai/agents/agent_adapters/base_agent_adapter.py:15`) and **`runtime_checkable` Protocols** for consumer-side structural contracts it does not own (memory storage `lib/crewai/src/crewai/memory/storage/backend.py:44-45`, third-party agent frameworks `lib/crewai/src/crewai/agents/agent_adapters/openai_agents/protocols.py:20-74`, human-input providers `lib/crewai/src/crewai/core/providers/human_input.py:59-66`, event handlers `lib/crewai/src/crewai/events/depends.py:17-37`). Contracts go beyond method signatures in several places: Pydantic schemas are auto-derived from `_run` signatures and validated at every call (`base_tool.py:200-247`, `272-293`); error semantics are documented in docstrings and enforced by dedicated exception types (`base_llm.py:315-318`; `memory/storage/backend.py:11-41`; `events/handler_graph.py:15-16`); lifecycle is explicit (`connect`/`disconnect`/async-context-manager on transports, `flush`/`shutdown`/atexit on the event bus, `event_bus.py:897-954`); and registries built via `__init_subclass__` make serialized tool/persistence classes round-trippable (`tools/base_tool.py:49-56,109-112`; `flow/persistence/base.py:32-35`). Validation is layered — schema-time (Pydantic model validation), signature-time (`structured_tool.py:248-270` checks that a function's parameters match its declared `args_schema`), registration-time (event-bus dependency cycle detection `event_bus.py:810-829`), and runtime per-call validation with actionable error hints. The main weaknesses are stylistic inconsistency (ABCs vs Protocols vs duck-typed adapters for the same kind of boundary), a few weakly-typed seams (`callbacks: list[Any]` in the LLM contract, `func: Any` on structured tools), an inconsistent failure contract between `BaseTool.run` (returns limit errors as strings, `base_tool.py:322-324`) and `CrewStructuredTool.invoke` (raises `ToolUsageLimitExceededError`, `structured_tool.py:349-352`), and no formal conformance suite proving independent `StorageBackend` implementations satisfy the protocol.

## Rating

**7 / 10** — Clear contract model with tests, explicit interfaces, and operational safeguards. The protocol+ABC mix, documented error taxonomy (`EmbeddingDimensionMismatchError` even documents *why* it is not a `RuntimeError`, `backend.py:11-41`), dependency-cycle validation at registration time, thread-safe usage accounting, and ~39 dedicated tool-contract tests (`tests/tools/test_base_tool.py`) earn the 7–8 band. It stops short of 9–10 because substitutability of storage backends rests on structural typing alone (no conformance suite or runtime check was found), several contract seams are typed `Any`, and the sync tool path violates LSP-style expectations by returning error strings instead of raising.

## Evidence Collected

Every entry includes a file path with line numbers relative to the source root.

| Area | Evidence | File:Line |
|------|----------|-----------|
| Tool interface (ABC + Pydantic) | `class BaseTool(BaseModel, ABC)` with required `name`/`description` fields and abstract `_run` | lib/crewai/src/crewai/tools/base_tool.py:103,375-391 |
| Async tool contract | Default `_arun` raises `NotImplementedError` with guidance to override or use sync `run()` | lib/crewai/src/crewai/tools/base_tool.py:356-365 |
| Schema-time argument contracts | `args_schema` auto-generated from `_run` signature via `field_validator`; falls back to `_arun` params | lib/crewai/src/crewai/tools/base_tool.py:200-247 |
| Runtime argument validation | `_validate_kwargs` validates kwargs against `args_schema` and appends a JSON-schema hint on failure | lib/crewai/src/crewai/tools/base_tool.py:272-293 |
| Result contracts | Optional `result_schema` inferred from return annotation; output validated before reaching the agent, degrading to `str()` with `RuntimeWarning` | lib/crewai/src/crewai/tools/base_tool.py:154-158,249-258; lib/crewai/src/crewai/tools/structured_tool.py:54-83 |
| Signature/schema consistency check | `_validate_function_signature` rejects required function parameters missing from `args_schema` at construction | lib/crewai/src/crewai/tools/structured_tool.py:248-270 |
| Usage-limit lifecycle | Thread-safe `_claim_usage` atomic counter; `max_usage_count <= 0` rejected by validator | lib/crewai/src/crewai/tools/base_tool.py:260-265,295-312 |
| Limit-failure inconsistency | `BaseTool.run` returns limit error as string; `CrewStructuredTool.invoke` raises `ToolUsageLimitExceededError` | lib/crewai/src/crewai/tools/base_tool.py:322-324; lib/crewai/src/crewai/tools/structured_tool.py:111-113,349-352 |
| Serialization registry | `_TOOL_TYPE_REGISTRY` populated in `__init_subclass__` so `list[BaseTool]` checkpoint fields resolve concrete classes; polymorphic core-schema validator dispatches on `tool_type` | lib/crewai/src/crewai/tools/base_tool.py:49-56,109-112,114-137 |
| Typed decorator entry point | `@tool` overloaded decorator requires docstring and annotations; `Tool(BaseTool, Generic[P, R])` uses ParamSpec | lib/crewai/src/crewai/tools/base_tool.py:496,676-761 |
| Storage consumer-side protocol | `@runtime_checkable class StorageBackend(Protocol)` — 13 fully documented methods incl. async mirrors (`asave`/`asearch`/`adelete`) | lib/crewai/src/crewai/memory/storage/backend.py:44-212 |
| Semantic error contract | `EmbeddingDimensionMismatchError(ValueError)` documents migration cause and remediation; explicitly not `RuntimeError` so background saves don't swallow it | lib/crewai/src/crewai/memory/storage/backend.py:11-41 |
| Error-contract test | Tests assert mismatch raises on save/search/update, that the error is not a `RuntimeError`, and that background saves propagate it while shutdown `RuntimeError`s stay swallowed | lib/crewai/tests/memory/test_dimension_mismatch.py:28-166,138 |
| Data contract | `MemoryRecord` Pydantic model with constrained `importance` (0–1) and embedding excluded from serialization | lib/crewai/src/crewai/memory/types.py:20-73 |
| Backend substitution plumbing | `Memory.storage` accepts `StorageBackend \| str` via `PlainValidator` pass-through; process-wide factory `set_memory_storage_factory` with explicit-instance-wins rule | lib/crewai/src/crewai/memory/unified_memory.py:92-95; lib/crewai/src/crewai/memory/storage/factory.py:33-55 |
| Structural implementations | `LanceDBStorage` and `QdrantEdgeStorage` implement the protocol without inheriting it (duck-typed conformance) | lib/crewai/src/crewai/memory/storage/lancedb_storage.py:42,289-660; lib/crewai/src/crewai/memory/storage/qdrant_edge_storage.py:81 |
| Factory selection tests | Tests cover factory registration, raw-spec delivery, explicit instance bypass | lib/crewai/tests/memory/test_storage_factory.py:25-58 |
| LLM interface | `BaseLLM(BaseModel, ABC)` with abstract `call`; documented raises: `ValueError`, `TimeoutError`, `RuntimeError`; `acall` optional | lib/crewai/src/crewai/llms/base_llm.py:129,283-319,321-357 |
| Config round-trip contract | `to_config_dict()` serializes so `LLM(**config)` reconstructs; subclasses told to extend via `super()` | lib/crewai/src/crewai/llms/base_llm.py:261-281 |
| Context propagation | ContextVar-based `llm_call_context` assigns per-call `call_id` used by all events; fallback UUID + warning when emitted outside context | lib/crewai/src/crewai/llms/base_llm.py:74-90,118-126 |
| Concurrency-safe override | `call_stop_override` scopes per-instance stop lists via ContextVar keyed by `id(llm)`, never mutating instance state | lib/crewai/src/crewai/llms/base_llm.py:93-115,200-214 |
| Capability predicates | `supports_stop_words`, `supports_multimodal`, `get_context_window_size`, `get_file_uploader`, `_effective_max_tokens` let callers branch without downcasting | lib/crewai/src/crewai/llms/base_llm.py:372-389,433-482 |
| Structured-output validation | `_validate_structured_output` parses JSON (with regex extraction fallback) into the requested Pydantic model, raising `ValueError` on failure | lib/crewai/src/crewai/llms/base_llm.py:835-871 |
| Message-format guardrails | `_format_messages` raises `ValueError` for non-dict messages or missing `role`/`content`; multimodal files on incapable models rejected with actionable message | lib/crewai/src/crewai/llms/base_llm.py:745-782,796-806 |
| Init-time normalization | Model validator requires non-empty `model`, coerces `stop` shapes, folds unknown kwargs into `additional_params` | lib/crewai/src/crewai/llms/base_llm.py:228-259 |
| Event envelope contract | `BaseEvent` carries `event_id`, `parent_event_id`, `previous_event_id`, `triggered_by_event_id`, `started_event_id`, `emission_sequence` | lib/crewai/src/crewai/events/base_events.py:66-88 |
| Listener lifecycle | `BaseEventListener.__init__` auto-registers listeners then runs `validate_dependencies()` eagerly | lib/crewai/src/crewai/events/base_event_listener.py:13-25 |
| Bus registration contract | Handler arity contract (2 or 3 args incl. optional `RuntimeState`); typed `on`/`off`; execution-plan cache invalidated on mutation | lib/crewai/src/crewai/events/event_bus.py:244-279,367-398 |
| Registration-time validation | `validate_dependencies()` builds plans for all dependent handlers to surface cycles/cross-event deps before emission; `CircularDependencyError` raised from cycle detection | lib/crewai/src/crewai/events/event_bus.py:810-829; lib/crewai/src/crewai/events/handler_graph.py:15-16,92 |
| Handler isolation | Sync handler errors captured per-handler and printed, not propagated; async handlers run under `gather(return_exceptions=True)` | lib/crewai/src/crewai/events/event_bus.py:415-427,444-455 |
| Replay semantics | `replay()` preserves stored event ids/sequence and sets an `is_replaying()` ContextVar so side-effectful listeners opt out | lib/crewai/src/crewai/events/event_bus.py:67-80,671-730 |
| Lifecycle/shutdown | Lazy executor init; `flush(timeout)` waits on tracked futures; `shutdown(wait)` drains loop+executor and clears handlers; registered via `atexit` | lib/crewai/src/crewai/events/event_bus.py:165-190,732-767,897-949,954 |
| Handler typing | Contravariant generic `EventHandler(Protocol[EventT_co])` plus FastAPI-style `Depends` with identity-based equality/hash | lib/crewai/src/crewai/events/depends.py:14-40,43-105 |
| Event-bus tests | Specific/multiple handler dispatch, error isolation, thread-safety, shutdown, RW-lock, async bus suites | lib/crewai/tests/utilities/events/test_crewai_event_bus.py:12-71; lib/crewai/tests/utilities/events/test_thread_safety.py; lib/crewai/tests/utilities/events/test_shutdown.py; lib/crewai/tests/utilities/events/test_rw_lock.py; lib/crewai/tests/utilities/events/test_async_event_bus.py |
| Transport interface | `BaseTransport(ABC)`: abstract `transport_type`, `connect()` documented to raise `ConnectionError`, `disconnect()`, mandatory async context-manager pair; stream accessors raise `RuntimeError` until connected | lib/crewai/src/crewai/mcp/transports/base.py:25-115 |
| Persistence interface | `FlowPersistence(BaseModel, ABC)`: abstract `init_db`/`save_state`/`load_state`; optional pending-feedback hooks default to safe no-ops; concrete subclasses auto-registered | lib/crewai/src/crewai/flow/persistence/base.py:18,32-35,70-116 |
| Knowledge storage interface | Abstract sync/async pairs `search`/`asearch`, `save`/`asave`, `reset`/`areset` returning typed `SearchResult` | lib/crewai/src/crewai/knowledge/storage/base_knowledge_storage.py:13-51 |
| Checkpoint provider interface | `BaseProvider`: abstract `checkpoint`/`acheckpoint` (returning location id), `prune` (returns removed count), `extract_id`, `from_checkpoint`/`afrom_checkpoint` | lib/crewai/src/crewai/state/provider/core.py:10-111 |
| Agent adapter interfaces | `BaseAgentAdapter(BaseAgent, ABC)` with `configure_tools`/`configure_structured_output`; `BaseToolAdapter(ABC)` with `configure_tools` | lib/crewai/src/crewai/agents/agent_adapters/base_agent_adapter.py:15-46; lib/crewai/src/crewai/agents/agent_adapters/base_tool_adapter.py:13-40 |
| Third-party isolation protocols | `runtime_checkable` protocols (`OpenAIAgent`, `OpenAIRunner`, `OpenAIAgentsModule`, `OpenAIFunctionTool`) describe external SDK shapes instead of importing them | lib/crewai/src/crewai/agents/agent_adapters/openai_agents/protocols.py:9-74; lib/crewai/src/crewai/agents/agent_adapters/langgraph/protocols.py:6-28 |
| HITL provider protocol | `HumanInputProvider(Protocol)` + `ExecutorContext` Protocol describing exactly what providers may touch on the executor | lib/crewai/src/crewai/core/providers/human_input.py:20-66 |
| Guardrail result contract | `GuardrailResult` standardizes guardrail returns; validator enforces `error`/`result` exclusivity; `from_tuple` adapts legacy `(bool, payload)` returns | lib/crewai/src/crewai/utilities/guardrail.py:60-106 |
| Compile-time support | Custom mypy plugin models attributes injected by `@CrewBase` so user code type-checks against the framework's implicit contract | lib/crewai/src/crewai/mypy.py:13-46 |
| Contract tests (tool) | ~39 tests: kwargs validation matrix (missing/wrong/extra keys, positional bypass), usage-count increment ordering vs validation errors, result-schema degradation warnings, async mirror parity | lib/crewai/tests/tools/test_base_tool.py:199-689 |

## Answers to Dimension Questions

### 1. Are interfaces small, coherent, and owned by the consumer side?

Mixed but trending toward consumer ownership. Consumer-side protocols exist where CrewAI consumes foreign code: `StorageBackend` (lib/crewai/src/crewai/memory/storage/backend.py:44), `EventHandler` (lib/crewai/src/crewai/events/depends.py:17), `ExecutorContext`/`HumanInputProvider` (lib/crewai/src/crewai/core/providers/human_input.py:20,59), and the OpenAI/LangGraph SDK shapes (lib/crewai/src/crewai/agents/agent_adapters/openai_agents/protocols.py:20-74). Provider-side ABCs are larger: `BaseLLM` spans ~1,045 lines including event-emission helpers, hooks invocation, message formatting, and token accounting (lib/crewai/src/crewai/llms/base_llm.py:129-1045), so a custom LLM inherits substantial machinery rather than implementing a narrow port. `StorageBackend` is also wide — 13 methods (backend.py:48-185) — though async mirrors are grouped separately (187-212). `BaseTool` mixes the execution contract (`_run`, `run`/`arun`) with serialization infrastructure (`__get_pydantic_core_schema__`, registry bookkeeping, base_tool.py:49-137).

### 2. Do contracts specify behavior, not just method signatures?

Yes, substantially, though unevenly. Docstring-documented error contracts appear on `BaseLLM.call` (raises `ValueError`/`TimeoutError`/`RuntimeError`, lib/crewai/src/crewai/llms/base_llm.py:315-318) and `BaseTransport.connect` (raises `ConnectionError`, lib/crewai/src/crewai/mcp/transports/base.py:75-77). Ordering/lifecycle behavior is specified: event scope pairing validates that ending events match their starting counterparts (lib/crewai/src/crewai/events/event_bus.py:547-564), replayed events skip mutation and expose `is_replaying()` (671-730 with 67-80), and `flush` guarantees handler completion within a timeout (732-767). Semantic guarantees are encoded in data contracts: `MemoryRecord.importance` bounded 0–1 (lib/crewai/src/crewai/memory/types.py:40-45), embeddings excluded from serialization (54-59), search results ordered by relevance (backend.py:76-78) and records listed newest-first (131-132). The dimension-mismatch error explains causality and three remediations inline (backend.py:23-41). Counter-examples: `BaseLLM.call`'s return type is `str | Any` (base_llm.py:293) leaving tool-call results untyped, and `callbacks: list[Any]` (288) specifies nothing.

### 3. Can providers, tools, stores, and runtimes be replaced safely?

Largely yes, with caveats. Memory stores are safely replaceable: the protocol is `runtime_checkable` (backend.py:44), the factory hook is process-scoped with an explicit "explicit instance always wins" rule (factory.py:33-43), and both built-ins conform structurally without inheritance (lancedb_storage.py:42). Tools substitute through one abstraction (`to_structured_tool`, base_tool.py:393-408) with schema carried over and tested (tests/tools/test_base_tool.py:576). LLM providers swap behind `BaseLLM` with capability predicates (`supports_stop_words`, `supports_multimodal`, base_llm.py:372-447) preventing callers from downcasting, and `to_config_dict` gives each provider a documented extension point (261-267). Caveats: hidden assumptions exist inside implementations rather than the interface — e.g., LanceDB adds out-of-protocol methods like `touch_records`/`optimize` (lancedb_storage.py:339,618), implying the unified memory layer may rely on capabilities no custom backend is obliged to provide; and nothing found verifies a third-party backend implements all 13 protocol members (see Gaps).

### 4. Are compatibility failures caught early by tests or validation?

Mostly. Registration-time: `BaseEventListener.__init__` calls `validate_dependencies()` so circular handler dependencies fail at listener construction, not mid-run (base_event_listener.py:17; handler_graph.py:92). Construction-time: tool function signatures must match declared schemas (structured_tool.py:265-270), `@tool` requires docstrings/annotations (base_tool.py:709-712), `max_usage_count <= 0` is rejected (260-265), and `BaseLLM` requires a non-empty `model` (base_llm.py:234-235). Call-time: kwargs validation with schema hints (base_tool.py:284-293), structured-output parsing (base_llm.py:852-871). Test-time: dedicated suites for tool validation matrices (test_base_tool.py:253-689), event-bus error isolation/thread-safety/shutdown (tests/utilities/events/*), storage factory selection (test_storage_factory.py:25-58), dimension-mismatch propagation across background saves (test_dimension_mismatch.py:138-166), and transport environment handling (test_stdio_transport.py:12-86). Not caught early: protocol-level conformance of custom storage backends (structural typing is only checked by static analysis, if users run it), and the `str`-returning limit failure on `BaseTool.run` which silently reads like a successful result.

## Architectural Decisions

1. **Pydantic-model ABCs as the canonical extension seam.** Every owned extension point composes `BaseModel, ABC` (tools base_tool.py:103, LLMs base_llm.py:129, persistence flow/persistence/base.py:18, knowledge knowledge/storage/base_knowledge_storage.py:13, providers state/provider/core.py:10), gaining declarative field validation, serialization, and schema generation "for free" at the cost of coupling extensions to Pydantic.
2. **Protocols for boundaries CrewAI does not own.** Foreign SDK shapes and consumer callbacks are `runtime_checkable` Protocols (openai_agents/protocols.py:20; depends.py:17; human_input.py:59), avoiding hard imports and keeping dependency direction inward.
3. **Schemas derived once, enforced everywhere.** A tool's `args_schema` is generated from the implementation signature (base_tool.py:200-247), validated per-call (272-293), embedded into the model-facing description (482-490), serialized/deserialized for checkpoints (160-170, structured_tool.py:28-37), and used to build failure hints (structured_tool.py:90-108).
4. **Registries over reflection for evolution.** `__init_subclass__`-populated registries let serialized payloads name concrete classes for later reconstruction — tools (base_tool.py:51,111-112) and flow persistence (flow/persistence/base.py:32-35) — supporting checkpoint compatibility across processes.
5. **ContextVars for cross-cutting context.** Call IDs (base_llm.py:74-90), stop overrides (93-115), replay flags (event_bus.py:67-80), and runtime state (83-91) propagate through async/thread-pool dispatch without polluting public signatures; sync dispatch copies contexts onto executor threads (event_bus.py:512-514,630-632).
6. **Errors-as-values only where agents consume them.** Tool-limit failures become agent-readable strings on the sync path (base_tool.py:307-312) because they are fed back into prompts, while programmatic paths raise typed exceptions (structured_tool.py:111-113) — deliberate but undocumented asymmetry.
7. **Compile-time accommodation of magic.** A shipped mypy plugin teaches type-checkers about `@CrewBase`-injected attributes (mypy.py:13-46), acknowledging that decorator-injected APIs break naive static analysis.

## Notable Patterns

- **Dual sync/async contracts with default fallback**: `acall` defaults to `NotImplementedError` (base_llm.py:357), `_arun` likewise (base_tool.py:356-365), letting providers adopt async incrementally; storage and knowledge interfaces declare async members as first-class protocol parts (backend.py:187-212).
- **Capability predicate methods** instead of interface splitting: `supports_stop_words()`, `supports_multimodal()`, `get_context_window_size()` (base_llm.py:372-447) — single interface, runtime capability discovery.
- **Template-method validation**: `run()` wraps abstract `_run()` with validation, usage claiming, and coroutine promotion (base_tool.py:314-331); subclasses never see those concerns.
- **Graceful-degradation warnings**: invalid `result_schema` outputs fall back to `str()` with a `RuntimeWarning` naming the schema (structured_tool.py:72-83), tested at tests/tools/test_base_tool.py:539-574.
- **Adapter triad for framework interop**: `BaseAgentAdapter`/`BaseToolAdapter`/converter adapters normalize LangGraph/OpenAI-Agents objects into CrewAI types (agents/agent_adapters/*.py), backed by structural protocols describing the foreign side (openai_agents/protocols.py).
- **Scope-based handler testing**: `scoped_handlers()` snapshots and restores bus registrations around a block (event_bus.py:831-895), making the singleton testable.
- **Explanatory exceptions**: exceptions carry constructor-computed diagnostics and remediation steps (backend.py:23-41), and validators append machine-derived schema hints to messages (base_tool.py:288-293).

## Tradeoffs

- **Inheritance tax vs. narrow ports**: `BaseLLM` implementers inherit event emission, hook plumbing, and formatting helpers (~900 lines around a single abstract method, base_llm.py:484-1045). Lower barrier to correct integrations, higher risk that overriding internals wrongly breaks invariants; a leaner protocol would be more substitutable.
- **Structural typing without enforcement**: Protocols keep imports clean, but nothing at runtime or in CI verifies third-party backends implement all members (no `isinstance(..., StorageBackend)` gate found anywhere in `src/` or `tests/`); failures surface as late `AttributeError`s.
- **String-typed failure channel on tools**: returning the limit error keeps the agent loop alive (base_tool.py:322-324) but means callers cannot distinguish failure from output without string matching — the structured-tool path made the opposite choice (raises, structured_tool.py:349-352).
- **Registry-driven deserialization vs. security/fragility**: resolving `tool_type` dotted paths imports arbitrary modules during `model_validate` (base_tool.py:59-65), trading checkpoint flexibility for an import-execution surface tied to persisted data.
- **Singleton bus with global state**: simplifies wiring (event_bus.py:952-954) and enables atexit safety, but forces snapshot/restore gymnastics (`scoped_handlers`) for test isolation and makes multi-tenant isolation impossible.
- **`Any`-typed convenience seams**: `embedder: Any` (unified_memory.py:96-99), `callbacks: list[Any]` (base_llm.py:288), `CrewStructuredTool.func: Any` (structured_tool.py:136) ease adoption but forfeit compile-time checking exactly where dynamic values flow.

## Failure Modes / Edge Cases

- **Handler crashes never propagate**: both sync and async dispatch swallow handler exceptions after printing (event_bus.py:415-427,450-455); a broken telemetry listener cannot kill a kickoff, but neither can callers learn about it programmatically except via `flush`'s collected futures (759-766).
- **Emitting outside a call context**: `get_current_call_id()` fabricates a fresh UUID with only a log warning (base_llm.py:118-126), so events emitted off-context correlate with nothing — trace assembly can silently fragment.
- **Positional arguments bypass validation entirely**: `run(*args)` skips `_validate_kwargs` when positionals are present (base_tool.py:319-320, tested at test_base_tool.py:285), an intentional escape hatch that weakens the schema guarantee for that call shape.
- **Cross-store dimension drift**: reopening a store with a changed embedder raises the actionable mismatch error everywhere data flows, including background-save paths, while genuine shutdown errors remain suppressed (test_dimension_mismatch.py:87-166; backend.py:11-41).
- **Unconnected transport access**: reading streams before `connect()` raises `RuntimeError` with remediation text (transports/base.py:54-66) rather than hanging on `None`.
- **Persistence resume with missing feedback markers**: `load_pending_feedback` defaults to `None` (flow/persistence/base.py:90-106), so backends that don't implement async-feedback support degrade to plain resume rather than corrupting state.
- **Replayed events double-firing side effects**: mitigated by the `is_replaying()` flag contract (event_bus.py:67-80,671-685), but compliance is voluntary per-listener — nothing enforces that checkpoint-writing listeners check it.

## Future Considerations

- Introduce a shared conformance suite (or a `verify_backend()` helper using `isinstance` against the `runtime_checkable` protocol plus behavioral probes) so third-party `StorageBackend` implementations are validated at startup instead of failing lazily.
- Unify the tool limit-failure contract: either raise `ToolUsageLimitExceededError` consistently and convert to agent-facing strings at the executor layer, or document the string sentinel as part of `run`'s contract.
- Narrow `BaseLLM` into composable ports (generation, streaming, capability discovery, observability) or formally freeze which protected methods subclasses may rely on; today the effective surface includes ~15 underscore helpers.
- Type the remaining `Any` seams (`callbacks`, `embedder`, `available_functions`) with protocols analogous to `EventHandler`.
- Document and version the `tool_type` dotted-path format used by the deserialization registry (base_tool.py:59-78) since it is now a wire-format commitment.

## Questions / Gaps

- **No cross-implementation storage conformance evidence.** Searched `tests/` for `isinstance(..., StorageBackend)`, `issubclass`, and "conformance" — no matches. Both built-ins are exercised through `Memory` integration tests (e.g., tests/memory/test_unified_memory.py, test_qdrant_edge_storage.py), but whether a minimal custom backend satisfies real call sites is unproven by any suite found.
- **Whether `QdrantEdgeStorage` implements every protocol member** was verified only by method-name inspection (qdrant_edge_storage.py:81); exhaustive signature matching was not mechanically confirmed within this study's scope.
- **Dependency direction of `BaseLLM` toward events**: the base class directly imports the concrete event bus and event types (base_llm.py:31-44) — the "interface" is entangled with the observability mechanism; whether custom LLMs can suppress this was not investigated.
- **MCP tool wrapper contract**: `tools/mcp_native_tool.py` and `mcp_tool_wrapper.py` were located but not analyzed in depth; how MCP tool schemas map into the `args_schema` contract deserves a follow-up pass.
- **A2A extension protocols** (a2a/extensions/base.py:39,56) and `hooks/types.py:16` were identified as additional runtime-checkable contracts but not individually audited.

---

Generated by `24.02-interface-contract-design` against `crewai`.
