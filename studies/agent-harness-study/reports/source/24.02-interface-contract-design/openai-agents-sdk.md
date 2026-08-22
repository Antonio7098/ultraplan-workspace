# Source Analysis: openai-agents-sdk

## 24.02 Interface Contract Design

### Source Info

| Field | Value |
|-------|-------|
| Name | openai-agents-sdk |
| Path | `studies/agent-harness-study/sources/openai-agents-sdk` |
| Language / Stack | Python 3.10+ (typing Protocols, ABCs, dataclasses, Pydantic, py.typed strict typing) |
| Analyzed | 2026-08-22 |

## Summary

The OpenAI Agents SDK defines its integration seams as a small set of named abstract types, one per extension point: `Model`/`ModelProvider` for LLM backends (`src/agents/models/interface.py:37`, `src/agents/models/interface.py:127`), `Session` for conversation memory (`src/agents/memory/session.py:14`), `Tool` as a closed union of tool dataclasses (`src/agents/tool.py:1355-1369`), `MCPServer` for tool servers (`src/agents/mcp/server.py:224`), `AgentOutputSchemaBase` for structured output (`src/agents/agent_output.py:16`), guardrails and handoffs as callable/dataclass contracts (`src/agents/guardrail.py:72`, `src/agents/handoffs/__init__.py:94`), plus parallel interface families for tracing (`src/agents/tracing/processor_interface.py:9`), realtime (`src/agents/realtime/model.py:151`), voice (`src/agents/voice/model.py:64`) and computer use (`src/agents/computer.py:8`). Contracts are enforced at three layers: type checking (the package ships `src/agents/py.typed`), schema-time validation (strict-schema normalization and `UserError`s raised in `__post_init__` of tools, e.g. `src/agents/tool.py:507-511`), and runtime fail-fast conversion errors at model boundaries (`src/agents/models/chatcmpl_converter.py:881-884`). Substitutability is proven rather than asserted: the test suite is built on independent implementations of the core contracts (`tests/fake_model.py:51`, `tests/mcp/helpers.py:70`, `tests/testing_processor.py:12`), and third-party adapters (`LitellmProvider` in `src/agents/extensions/models/litellm_provider.py:9-23`) satisfy the same `ModelProvider` interface. The most distinctive property is that compatibility is treated as an explicit engineering discipline: positional constructor order of public dataclasses is documented as a contract with in-code comments (`src/agents/tool.py:417-418`) and a repo-wide policy (AGENTS.md, "Public API Positional Compatibility"). Weaknesses are concentrated in a few places: the `Model` signature leaks OpenAI Responses-API concepts into the generic interface (`src/agents/models/interface.py:56-70`), the `Tool` union is closed to third parties, and a couple of runtime checks bypass the declared `Session` protocol by isinstance-testing a concrete class (`src/agents/run_internal/session_persistence.py:89`).

## Rating

**8 / 10** — Clear, documented contracts with dual Protocol+ABC patterns, schema-time validation, semantic guarantees (retry advice, approval semantics), and independent test doubles proving substitutability. Falls short of 9–10 because of OpenAI-specific parameter leakage in the generic `Model` interface, a closed `Tool` union that cannot be extended externally, concrete-class isinstance checks that undermine the `Session` protocol's substitutability promise, and thread-safety/lifecycle obligations stated only in docstrings.

## Evidence Collected

| Area | Evidence | File:Line |
|------|----------|-----------|
| Model interface | `Model` ABC: two abstract methods (`get_response`, `stream_response`) sharing identical 10-parameter signatures; default no-op `close()`; optional `get_retry_advice()` hook | `src/agents/models/interface.py:37-124`, `src/agents/models/interface.py:40-46`, `src/agents/models/interface.py:48-54` |
| Provider lookup contract | `ModelProvider.get_model(model_name) -> Model` plus overridable async `aclose()` lifecycle | `src/agents/models/interface.py:127-150` |
| Consumer-side resolution | Runner resolves models through `run_config.model_provider.get_model(...)` or an injected `Model` instance — the runner owns the abstraction, providers just implement it | `src/agents/run_internal/turn_preparation.py:126-135` |
| Session protocol + ABC duality | `Session` is a `@runtime_checkable Protocol`; a separate internal-only `SessionABC` exists for in-repo implementations, explicitly telling third parties to implement the protocol instead | `src/agents/memory/session.py:13-54`, `src/agents/memory/session.py:57-65` |
| Capability detection | Compaction support detected via a dedicated `OpenAIResponsesCompactionAwareSession` protocol and defensive `TypeGuard` helper instead of inheritance | `src/agents/memory/session.py:131-150` |
| Tool contract shape | `FunctionTool` dataclass centers the contract on `on_invoke_tool(ctx, json_str)` with documented return obligations (structured output types or str()-able values) | `src/agents/tool.py:380-406` |
| Closed tool union | `Tool = FunctionTool \| FileSearchTool \| WebSearchTool \| ...` — enumeration, not an interface | `src/agents/tool.py:1355-1369` |
| Output schema contract | `AgentOutputSchemaBase` couples JSON schema production with `validate_json` error semantics ("must raise `ModelBehaviorError`") | `src/agents/agent_output.py:16-51`, reference impl `src/agents/agent_output.py:55-164` |
| MCP server contract | `MCPServer` ABC: `connect/name/cleanup/list_tools/call_tool/list_prompts/get_prompt` abstract; resources are optional via default `NotImplementedError` methods naming the missing override | `src/agents/mcp/server.py:264-302`, `src/agents/mcp/server.py:328-383` |
| Guardrail contract | `InputGuardrail`/`OutputGuardrail` wrap a callable `(context, agent, input) -> MaybeAwaitable[GuardrailFunctionOutput]` with tripwire semantics documented on the dataclass | `src/agents/guardrail.py:71-130`, `src/agents/guardrail.py:133-185` |
| Handoff contract | `Handoff` dataclass exposes `on_invoke_handoff(ctx, args_json) -> TAgent` plus input filter hook typed as `HandoffInputFilter` alias; `HandoffInputData.run_context` optional "for backwards compatibility" | `src/agents/handoffs/__init__.py:93-140`, `src/agents/handoffs/__init__.py:42-87` |
| Tracing processor contract | `TracingProcessor` ABC with lifecycle (`shutdown`, `force_flush`) and behavioral notes: "All methods should be thread-safe", "Should not block" — doc-level, not enforced | `src/agents/tracing/processor_interface.py:9-129` |
| Realtime / voice contracts | `RealtimeModel` connect/add_listener/send_event/close; voice `TTSModel`, `STTModel`, `StreamedTranscriptionSession` with explicit streaming semantics ("expected to return only after `close()` is called") | `src/agents/realtime/model.py:151-177`, `src/agents/voice/model.py:64-100` |
| Computer contract | Sync `Computer` and async `AsyncComputer` ABCs define the full action surface (screenshot/click/type/drag...) consumed by `ComputerTool` | `src/agents/computer.py:8-69`, `src/agents/tool.py:714-738` |
| Schema-time validation | `ensure_strict_json_schema` rewrites schemas and raises `UserError` on unsupported `additionalProperties`; applied automatically in `FunctionTool.__post_init__` and output-schema construction | `src/agents/strict_schema.py:18-64`, `src/agents/tool.py:507-510`, `src/agents/agent_output.py:112-120` |
| Constructor-time validation | `ShellTool.__post_init__` rejects invalid executor/environment combinations; `ToolExecutionConfig` validates concurrency limits | `src/agents/tool.py:1238-1256`, `src/agents/run_config.py:112-118` |
| Runtime fail-fast conversion | Chat Completions adapter rejects hosted tools and Responses-only function-tool features with explicit `UserError` naming the backend | `src/agents/models/chatcmpl_converter.py:865-884`, `src/agents/tool.py:1408-1423` |
| Error hierarchy | Single `AgentsException` base with typed subclasses per failure mode (MaxTurnsExceeded, ModelBehaviorError, tripwire exceptions, ToolTimeoutError, MCPToolCancellationError) | `src/agents/exceptions.py:46-162` |
| Tool error contract | Errors become model-visible strings via `default_tool_error_function`; timeout behavior configurable between `error_as_result` and `raise_exception` | `src/agents/tool.py:1609-1618`, `src/agents/tool.py:182`, `src/agents/tool.py:436-444` |
| Semantic retry contract | `ModelRetryAdvice` carries replay-safety/retry-after/provider guidance; `Model` adapters can override `get_retry_advice(request)` — behavior, not just signatures | `src/agents/retry.py:93-137`, `src/agents/models/interface.py:48-54` |
| Approval/HITL semantics | `needs_approval` accepts bool or callable receiving `(run_context, params, call_id)`; interruption resolved via `RunState.approve()/reject()` — contract documented on every tool variant | `src/agents/tool.py:426-433`, `src/agents/tool.py:1219-1231` |
| Independent Model double | `FakeModel(Model)` implements the full interface including stream event sequencing; used across the entire test suite | `tests/fake_model.py:51-365` |
| Independent MCP double | `FakeMCPServer(MCPServer)` implements all abstract members without touching real MCP transport; `_TestFilterServer` subclasses the internal client-session base | `tests/mcp/helpers.py:70-164`, `tests/mcp/helpers.py:47-67` |
| Tracing test double | Thread-safe `SpanProcessorForTests(TracingProcessor)` honors the documented thread-safety obligation | `tests/testing_processor.py:12-59` |
| Third-party provider adapter | `LitellmProvider(ModelProvider)` is a ~15-line implementation proving the provider seam is small enough to satisfy independently | `src/agents/extensions/models/litellm_provider.py:9-23` |
| Prefix-routed provider | `MultiProvider` composes providers by model-name prefix with explicit alias/model-id mode options to avoid breaking historical semantics | `src/agents/models/multi_provider.py:61-156` |
| Positional compat policy | Comment pins field order: "Keep guardrail fields before needs_approval to preserve v0.7.0 positional constructor compatibility"; new fields appended kw-only | `src/agents/tool.py:417-435` |
| Persisted-name compat | `ComputerTool.name` deliberately keeps preview-era name `"computer_use_preview"` for hooks and persisted RunState compatibility while tracing uses the GA name | `src/agents/tool.py:733-743` |
| Documented design goal | Docs enumerate where each customization seam belongs (`ModelProvider` at run level vs `Model` instance) tying interface choice to usage guidance | `docs/models/index.md:219`, `docs/models/index.md:226`, `docs/models/index.md:251` |
| Substitutability gap | Runtime code isinstance-checks concrete `OpenAIConversationsSession` rather than a protocol capability, so third-party sessions cannot opt into that behavior path | `src/agents/run_internal/session_persistence.py:89`, `src/agents/run_internal/session_persistence.py:292`, `src/agents/run_internal/session_persistence.py:565-569` |

## Answers to Dimension Questions

### 1. Are interfaces small, coherent, and owned by the consumer side?

Mostly yes. The consumer (runner) owns the abstractions it calls: `get_model()` in the turn-preparation layer consumes only `Model`/`ModelProvider` (`src/agents/run_internal/turn_preparation.py:126-135`), and the agent consumes only `MCPServer.list_tools/call_tool` via `MCPUtil` (`src/agents/agent.py:224-244`). Provider implementations stay thin — `LitellmProvider` is a single method (`src/agents/extensions/models/litellm_provider.py:22-23`). Two caveats: `FunctionTool` has grown to ~20 fields mixing public config, approval policy, timeouts, and private runtime metadata (`src/agents/tool.py:380-494`), and `Model.get_response` takes ten parameters including three OpenAI Responses-API-specific ones (`previous_response_id`, `conversation_id`, `prompt`, `src/agents/models/interface.py:66-69`) whose docstrings admit they are "generally not used by the model, except for the OpenAI Responses API" (`src/agents/models/interface.py:81-84`) — leakage of one provider's vocabulary into the generic seam.

### 2. Do contracts specify behavior, not just method signatures?

Substantially yes, with varying rigor. Strong cases: `AgentOutputSchemaBase.validate_json` specifies the exception type callers must raise (`src/agents/agent_output.py:47-51`); retry advice encodes replay safety and delay semantics (`src/agents/retry.py:93-100`); compaction modes are enumerated with meaning ("auto"/"previous_response_id"/"input", `src/agents/memory/session.py:113-119`); approval semantics specify the exact resume mechanism (`RunState.approve()/reject()`, `src/agents/tool.py:429-433`); the streamed transcription session specifies when the iterator must terminate relative to `close()` (`src/agents/voice/model.py:90-95`). Weaker case: `TracingProcessor` obligations ("thread-safe", "should not block", `src/agents/tracing/processor_interface.py:48-50`) exist only as docstring guidance with no enforcement, and `FunctionTool.on_invoke_tool`'s return contract is partially structural — "or something we can call `str()` on" leaves output shape loosely typed as `Any` (`src/agents/tool.py:400-406`).

### 3. Can providers, tools, stores, and runtimes be replaced safely?

Models, model providers, MCP servers, and tracing processors: demonstrably yes — the suite runs entirely against independent fakes implementing only the public contract (`tests/fake_model.py:51`, `tests/mcp/helpers.py:70`), and LiteLLM/any-LLM adapters plug in through the same ABCs (`src/agents/extensions/models/litellm_provider.py:9-23`, `src/agents/models/multi_provider.py:61`). Sessions: mostly — any object satisfying the `Session` protocol works (`src/agents/memory/session.py:13-14`), but two runtime paths silently special-case the SDK's own `OpenAIConversationsSession` via concrete isinstance checks (`src/agents/run_internal/session_persistence.py:89`, `src/agents/run_internal/session_persistence.py:565-569`), so a structurally identical third-party session gets different ID-matching behavior without any error. Tools: not fully — `Tool` is a closed union (`src/agents/tool.py:1355-1369`); a custom tool type cannot be added without editing the SDK, and converters respond to unknown types with `UserError` rather than an extension point (`src/agents/models/chatcmpl_converter.py:881-884`). Functionality can still be added via `FunctionTool` composition, which is the sanctioned path.

### 4. Are compatibility failures caught early by tests or validation?

Three early-detection layers exist. (a) Schema/construction time: strict-schema normalization raises `UserError` when a tool or output type cannot meet the strict standard (`src/agents/strict_schema.py:59-64`, applied at `src/agents/tool.py:507-510` and `src/agents/agent_output.py:112-120`), and dataclass `__post_init__` validators catch misconfiguration before any run starts (`src/agents/tool.py:1238-1256`). (b) Type time: `py.typed` plus strict typing means interface drift surfaces in `make typecheck`. (c) Run boundaries: backend mismatch is caught at first conversion with actionable messages naming the unsupported feature and backend (`src/agents/tool.py:1408-1423`). Additionally, compatibility itself is regression-tested: positional-constructor compatibility is called out with comments and tests exercising old positional call patterns per AGENTS.md policy, and serialized-state evolution is versioned via `CURRENT_SCHEMA_VERSION` with per-version summaries in `src/agents/run_state.py` (referenced from AGENTS.md "Agents Core Runtime Guidelines"). What is not caught early: docstring-level obligations (thread safety, non-blocking processors) and the concrete-isinstance session paths described above.

## Architectural Decisions

1. **ABCs for infrastructure seams, Protocols for user-facing data-ish contracts.** `Model`, `ModelProvider`, `MCPServer`, `TracingProcessor`, `Computer` are ABCs requiring inheritance (`src/agents/models/interface.py:37`, `src/agents/tracing/processor_interface.py:9`), while `Session` is a `runtime_checkable` Protocol with an internal `SessionABC` reserved for SDK implementations — the docstring explicitly directs third-party libraries to the protocol (`src/agents/memory/session.py:57-65`). This splits "framework extension point" from "duck-typed capability".

2. **Closed unions over open interfaces for model-facing payloads.** Both `Tool` (`src/agents/tool.py:1355-1369`) and `RunItem` (`src/agents/items.py:639-654`) are fixed unions of dataclasses anchored to OpenAI wire types (type aliases `TResponseInputItem` etc., `src/agents/items.py:73-80`). This gives exhaustive matching inside converters at the cost of external extensibility.

3. **Wire-format anchoring.** Nearly every contract type aliases or wraps an `openai.types.responses.*` type, so the OpenAI Responses API is the canonical interlingua; other backends must translate into it (`src/agents/models/chatcmpl_converter.py` is a full converter module). This makes OpenAI substitution trivial and everyone else's substitution lossy-by-design.

4. **Capability protocols + TypeGuards instead of marker flags.** Compaction-aware sessions and resource-capable MCP servers are handled by optional protocols and default `NotImplementedError` methods that tell implementers exactly what to override (`src/agents/memory/session.py:131-150`, `src/agents/mcp/server.py:342-345`).

5. **Compatibility as an explicit, tested artifact.** Field ordering pinned with comments (`src/agents/tool.py:417-418`), renamed-but-preserved runtime names (`src/agents/tool.py:733-738`), optional fields marked backwards-compatible in their own docstrings (`src/agents/handoffs/__init__.py:60-64`), versioned persisted state (`src/agents/run_state.py` schema versions), and prefix-mode switches that let new semantics coexist with historical behavior without breaking callers (`src/agents/models/multi_provider.py:68-72`, `116-124`).

## Notable Patterns

- **Dual Protocol+ABC declaration**: same members documented twice, once duck-typed for users, once inherited internally (`src/agents/memory/session.py:13-104`).
- **Default-no-op lifecycle hooks**: `Model.close()`, `ModelProvider.aclose()` default to no-op so stateless implementations need no boilerplate (`src/agents/models/interface.py:40-46`, `144-150`).
- **`MaybeAwaitable[T]` everywhere**: callbacks may be sync or async (`src/agents/util/_types.py:7`), with `inspect.isawaitable` normalization at each call site (`src/agents/guardrail.py:120-125`, `src/agents/agent.py:207-214`).
- **Constructor-time normalization**: tools validate and normalize their own configuration in `__post_init__` so downstream code sees a canonical shape (`src/agents/tool.py:503-511`, `1238-1256`).
- **Fail-fast boundary guards**: shared helpers like `ensure_function_tool_supports_responses_only_features(tool, backend_name=...)` produce uniform cross-backend rejection messages (`src/agents/tool.py:1408-1423`, invoked at `src/agents/models/chatcmpl_converter.py:867-870`).
- **Defensive deserialization at trust boundaries**: `ToolOrigin.from_json_dict` returns `None` on malformed persisted data rather than raising, keeping old snapshots loadable (`src/agents/tool.py:298-322`).
- **Factory-bound invoker rebinding**: copied tools rebind failure-handling invokers via a `__agents_bind_function_tool__` protocol so copies resolve error policy against themselves, avoiding identity bugs after `clone()` (`src/agents/tool.py:503-506`, `526-552`).

## Tradeoffs

- **OpenAI-first interlingua vs neutrality**: anchoring everything on Responses wire types makes the OpenAI path exact and strongly typed, but pushes complexity into converters and leaks Responses concepts (`previous_response_id`, `prompt`) into `Model`'s required signature (`src/agents/models/interface.py:66-70`) — every alternative provider must accept parameters it ignores.
- **Closed unions vs extensibility**: exhaustive `Tool`/`RunItem` unions enable precise internal handling and serialization but mean new tool categories require coordinated edits across items, steps, streaming, and state modules — a coupling the repo itself documents as mandatory checklist behavior (AGENTS.md, list of files to update when adding item types).
- **Protocol flexibility vs silent divergence**: accepting any protocol-satisfying `Session` invites third-party implementations, yet hardcoded concrete-class checks create undocumented behavior forks (`src/agents/run_internal/session_persistence.py:89`) that neither types nor tests will flag for external authors.
- **Docstring contracts vs enforceability**: rich prose contracts scale documentation cheaply but rely on reviewer diligence; nothing prevents a blocking `TracingProcessor.on_span_end`.
- **kw_only appends vs readability**: appending new `FunctionTool` fields as keyword-only preserves positional compatibility (`src/agents/tool.py:452-494`) but yields very long constructor signatures.

## Failure Modes / Edge Cases

- **Unsupported feature/backend combos** fail loudly at conversion with `UserError` listing the offending features (`src/agents/models/chatcmpl_converter.py:881-884`, `src/agents/tool.py:1419-1423`) rather than silently dropping tools.
- **Malformed persisted tool-origin metadata** degrades gracefully to `None` (`src/agents/tool.py:301-315`).
- **Optional capabilities** degrade to explicit `NotImplementedError` messages naming the override needed (`src/agents/mcp/server.py:342-383`), and misuse of `ToolApprovalItem.to_input_item()` raises immediately to surface filtering mistakes (`src/agents/items.py:631-636`).
- **Tool invocation failures** have a two-way contract: convert to a model-visible string or raise to fail the run, chosen per-tool via `failure_error_function`/`timeout_behavior` (`src/agents/tool.py:400-406`, `439-447`, `182`).
- **Cancellation** is a first-class typed outcome: `MCPToolCancellationError` (`src/agents/exceptions.py:99`) and idempotent streaming cancellation covered by dedicated tests (`tests/test_cancel_streaming.py:80`, `96`).
- **Edge case left implicit**: a third-party `Session` that satisfies the protocol but not the concrete class will take different dedupe/fingerprint paths with no warning — substitutability failure that fails silently rather than loudly.

## Future Considerations

- Narrow the `Model` call signature: group Responses-specific parameters (`previous_response_id`, `conversation_id`, `prompt`) into an options object so non-OpenAI providers stop carrying ignored parameters (`src/agents/models/interface.py:56-105`).
- Replace concrete isinstance checks on sessions with a capability protocol/TypeGuard, mirroring the existing compaction pattern (`src/agents/run_internal/session_persistence.py:565-569` vs `src/agents/memory/session.py:140-150`).
- Define an open extension point or registry for custom tool types, or document the closed-union decision prominently next to `Tool` (`src/agents/tool.py:1355-1369`).
- Promote `TracingProcessor` behavioral obligations (thread safety, non-blocking) into a conformance test helper analogous to `tests/testing_processor.py:12` so third-party processors can self-check.
- Consolidate the scattered sync/async normalization (`inspect.isawaitable` repeated at many call sites, e.g. `src/agents/agent.py:207-214`, `src/agents/guardrail.py:120-130`) into a shared awaitable-resolution utility to remove per-site drift risk.

## Questions / Gaps

- No formal conformance suite was found for `Model` implementations beyond ad-hoc fakes: searches across `tests/` for "conformance" returned nothing; verification of provider correctness relies on per-backend test files (e.g. `tests/models/test_litellm_chatcompletions_stream.py`) rather than a shared contract test that all `Model` implementations could run against.
- The `Session` protocol's `session_settings` attribute defaults to `None` on a Protocol (`src/agents/memory/session.py:22`); how strictly attribute presence is checked by consumers is unclear from the code inspected (protocol checks were not observed being used via `isinstance` for `Session` itself outside concrete-class checks).
- Whether the OpenAI-specific parameters on `Model` are planned to be generalized could not be determined from code alone; no TODO/deprecation markers were found in `src/agents/models/interface.py`.

---

Generated by `24.02-interface-contract-design` against `openai-agents-sdk`.
