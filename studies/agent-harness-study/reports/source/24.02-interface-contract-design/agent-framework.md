# Source Analysis: agent-framework

## Interface Contract Design

### Source Info

| Field | Value |
|-------|-------|
| Name | agent-framework |
| Path | `studies/agent-harness-study/sources/agent-framework` |
| Language / Stack | Python 3 (Protocol + ABC + runtime_checkable; Pydantic, async/await) and .NET 9+ (abstract classes; Microsoft.Extensions.AI abstractions; nullable reference types) |
| Analyzed | 2026-08-22 |

## Summary

Agent Framework maintains parallel implementations (Python and .NET) of the same agent runtime abstraction. The Python side combines `@runtime_checkable` `Protocol` classes for consumer-defined duck-typed surfaces (`_clients.py:82`, `_clients.py:864`, `_compaction.py:42,51`, `_skills.py:1495`, `_workflows/_runner_context.py:98`, `_serialization.py:24`) with `ABC` classes for owned extension points (`_clients.py:214`, `_middleware.py:465,524,588`, `_sessions.py:351`, `_mcp.py:1559`); the .NET side uniformly uses `abstract` classes (`AIAgent.cs:38`, `AIContextProvider.cs:42`, `AgentSession.cs:59`, `ChatHistoryProvider.cs:51`, `DelegatingAIAgent.cs:28`, `Workflows/Executor.cs:164`). Protocols/ABCs are small, named, and consumer-centric in ownership: a chat client, an embedding client, an agent, a script runner, and an internal `RunnerContext` are each a single-method or two-method surface (`_clients.py:82-197`, `_clients.py:864-911`, `_compaction.py:42-63`, `_skills.py:1495-1529`, `_workflows/_runner_context.py:98-280`). Provider substitution is verified by independent test implementations written as plain duck-typed classes rather than framework subclasses (`conftest.py:82-133 MockChatClient`, `conftest.py:273-320 MockAgent`, `test_embedding_client.py:84-104 PlainEmbeddingClient`), and the .NET repository ships a full `AgentConformance.IntegrationTests` project (`dotnet/tests/AgentConformance.IntegrationTests/RunTests.cs:17`) that drives the same `RunAsync` / `RunStreamingAsync` / `CreateSessionAsync` lifecycle assertions across heterogeneous agent fixtures.

Error contracts are hierarchical: a single root `AgentFrameworkException` (`exceptions.py:15-38`) fans into domain-specific bases (`AgentException`, `ChatClientException`, `IntegrationException`, `ToolException`, `MiddlewareException`, `WorkflowException`) and leaves each layer free to attach an `inner_exception` that is automatically debug-logged. Tool execution surfaces a separate `ToolExecutionException` (`exceptions.py:178`) and a `UserInputRequiredException` carrying re-emitted approval-request content (`exceptions.py:184-209`). Compile-time / static checks are layered: typed `OptionsContraT` and `EmbeddingProtocolOptionsT` `TypeVar`s push provider-specific `TypedDict` options through `@overload`-decorated `get_response` signatures (`_clients.py:70-178`), and `pyrightconfig.tests.json` + `pyrefly.toml` are committed to the repo. Runtime validation is partial — `validate_chat_options` enforces numeric ranges on `temperature`/`top_p`/`frequency_penalty`/`presence_penalty`/`max_tokens` (`_types.py:3428-3482`) and `validate_tool_mode` rejects malformed `tool_choice` dicts (`_types.py:3569-3603`), but the only "conformance" test for substitutability is a smoke assertion (`test_embedding_client.py:76-104`); the Python suite does not yet run every protocol through a generic cross-provider conformance harness.

Semantic guarantees are encoded mostly in the .NET side as XML docs that read like a contract specification: `AIContextProvider` documents the two-phase `InvokingAsync` / `InvokedAsync` lifecycle, the default filters (`ExternalOnly` for inputs, no-op for responses), `StateKeys` for safe multi-instance session reuse, and explicit security considerations for system-message injection and indirect prompt injection (`AIContextProvider.cs:42-419`); `ChatHistoryProvider` mirrors this with chronology and filter defaults (`ChatHistoryProvider.cs:51-477`); `AIAgent` documents that sessions are not portable across agents and that serialized sessions must be treated as untrusted input (`AIAgent.cs:24-36, 178-220`, `AgentSession.cs:11-58`). Python is more concise: docstrings cover lifecycle but rely on `@runtime_checkable` plus example-based docs (e.g. `_clients.py:81-198` shows a `CustomChatClient` class with no inheritance and an `isinstance` check).

Substitutability across providers (Anthropic, Bedrock, Azure AI Foundry, Ollama, Claude Agent SDK, Gemini, GitHub Copilot, Mem0, Redis, OpenAI) is implemented by **decorator stacking of layer classes** around `BaseChatClient` rather than interface conformance. The Python stack wires `FunctionInvocationLayer` + `ChatMiddlewareLayer` + `ChatTelemetryLayer` + `BaseChatClient` (`conftest.py:136-177`, `_middleware.py:1104, 794`) so any provider `BaseChatClient` subclass gets middleware, telemetry, and function-invocation for free. The .NET equivalent is `DelegatingAIAgent` (`AIAgent.cs` references / `DelegatingAIAgent.cs:28-101`), `FunctionInvocationDelegatingAgent`, `OpenTelemetryAgent`, `LoggingAgent` (`dotnet/src/Microsoft.Agents.AI/`), and concrete chat-client decorators (`ChatClient/MessageInjectingChatClient.cs`, `PerServiceCallChatHistoryPersistingChatClient.cs`, `NonApprovalRequiredFunctionBypassingChatClient.cs`). This means the architecture deliberately separates "what an agent does" from "what an agent wraps": a provider only implements `_inner_get_response` / `RunCoreAsync`, and the decorator stack supplies observable behavior uniformly.

The contracts are **durable, observable, and extensible** within their design domain, but several seams show clear fragility: many Protocol/ABC members use `**kwargs` (`_clients.py:177`, `_workflows/_runner_context.py:170`), `SupportsChatGetResponse.additional_properties` and `BaseAgent.additional_properties` are required for protocol conformance (`_clients.py:127`, `_agents.py:409`) but not declared on the runtime check path of `MockAgent` in conftest (it inherits them via `SupportsAgentRun`'s implicit Protocol rules), and the Python `ChatResponse`/`AgentResponse` defaults are anchored on a single `ResponseStream` class (`_types.py:2939`) that mixes an `AsyncIterable` with a finalizer hook — non-trivial for an alternative implementation to reproduce.

## Rating

**8 / 10** — Clear model with tests, explicit interfaces, and operational safeguards. Earned by (a) narrow, named, owner-defined Protocols/ABCs for every cross-boundary extension point; (b) decorator-stack composability that allows a single implementor to satisfy middleware/telemetry/function-invocation without knowing about those layers; (c) a real cross-provider conformance harness on the .NET side (`AgentConformance.IntegrationTests`); (d) a typed exception hierarchy with explicit `inner_exception` plumbing and `log_level` controls; and (e) explicit security-consideration XML docs on every state-touching .NET abstraction (`AIAgent.cs:24-36`, `AIContextProvider.cs:34-40`, `ChatHistoryProvider.cs:43-49`, `AgentSession.cs:46-53`). Falls short of 9–10 because (i) Python has no equivalent conformance harness — substitutability is only smoke-tested with one mock per protocol (`test_embedding_client.py:76-104`, `test_clients.py:33-43`); (ii) many Protocol methods accept `**kwargs` and any-object `Mapping[str, Any]` options, leaving semantic behavior in the implementation; (iii) optional `additional_properties` and decorator-stack configuration create hidden assumptions about which layer is responsible for persistence or history loading; (iv) Python type validation (`validate_chat_options`) only checks numeric ranges, not provider-specific option semantics.

## Evidence Collected

Every entry includes a workspace-relative file path with line numbers.

| Area | Evidence | File:Line |
|------|----------|-----------|
| Agent Protocol (duck-typed contract) | `@runtime_checkable class SupportsAgentRun(Protocol)` with `id`, `name`, `description`, overloaded `run`, `create_session`, `get_session` | `sources/agent-framework/python/packages/core/agent_framework/_agents.py:187-312` |
| Minimal agent base class | `class BaseAgent(SerializationMixin)` with no `run` implementation — explicitly tells the reader subclassing is required | `sources/agent-framework/python/packages/core/agent_framework/_agents.py:318-449` |
| Agent as tool contract check | `if not isinstance(self, SupportsAgentRun): raise TypeError(...)` — uses Protocol `isinstance` for runtime gating | `sources/agent-framework/python/packages/core/agent_framework/_agents.py:530-532` |
| Chat client Protocol | `@runtime_checkable class SupportsChatGetResponse(Protocol[OptionsContraT])` with `additional_properties` field and 4 `@overload`-decorated `get_response` signatures | `sources/agent-framework/python/packages/core/agent_framework/_clients.py:81-197` |
| Chat client ABC base | `class BaseChatClient(SerializationMixin, ABC, Generic[OptionsCoT])` exposing shared logic; subclasses implement `_inner_get_response` | `sources/agent-framework/python/packages/core/agent_framework/_clients.py:214-470` |
| Tool-specific capability Protocols | `SupportsCodeInterpreterTool`, `SupportsWebSearchTool`, `SupportsImageGenerationTool`, `SupportsMCPTool`, `SupportsFileSearchTool`, `SupportsShellTool` — each `@runtime_checkable` with one method | `sources/agent-framework/python/packages/core/agent_framework/_clients.py:665-833` |
| Embedding Protocol | `@runtime_checkable class SupportsGetEmbeddings(Protocol[...])` plus `BaseEmbeddingClient` ABC | `sources/agent-framework/python/packages/core/agent_framework/_clients.py:864-911, 919-957` |
| Provide-as-Agent convenience | `BaseChatClient.as_agent(...)` constructs `Agent` with the client and `default_options` — keeps the contract narrow at runtime | `sources/agent-framework/python/packages/core/agent_framework/_clients.py:568-655` |
| Middleware ABCs | `AgentMiddleware`, `FunctionMiddleware`, `ChatMiddleware` as ABCs with `@abstractmethod async def process(self, context, call_next)` | `sources/agent-framework/python/packages/core/agent_framework/_middleware.py:465-628` |
| Middleware pipeline abstraction | `class BaseMiddlewarePipeline(ABC)` plus per-layer `ChatMiddlewareLayer`, `FunctionInvocationLayer`, `AgentMiddlewareLayer` (decorator-stack pattern) | `sources/agent-framework/python/packages/core/agent_framework/_middleware.py:777-1260` |
| Context provider / history base classes | `class ContextProvider` with default no-op `before_run`/`after_run`; `class HistoryProvider(ContextProvider)` with two `@abstractmethod`s `get_messages`, `save_messages` | `sources/agent-framework/python/packages/core/agent_framework/_sessions.py:351-496` |
| Serialization Protocol + Mixin | `@runtime_checkable SerializationProtocol(Protocol)` with `to_dict` / `from_dict`; `SerializationMixin` provides nested serialization | `sources/agent-framework/python/packages/core/agent_framework/_serialization.py:23-130, 137-419` |
| Tokenizer + Compaction Protocols | `@runtime_checkable class TokenizerProtocol(Protocol)` and `@runtime_checkable class CompactionStrategy(Protocol)` callable contract | `sources/agent-framework/python/packages/core/agent_framework/_compaction.py:41-63` |
| Skill Runner Protocol | `@runtime_checkable class SkillScriptRunner(Protocol)` — any callable matching `__call__(self, skill, script, args)` satisfies it | `sources/agent-framework/python/packages/core/agent_framework/_skills.py:1493-1529` |
| Skill ABCs | `class Skill(ABC)` with abstract `frontmatter` property and `get_content`; `SkillResource`, `SkillScript`, `SkillsSource` also ABC | `sources/agent-framework/python/packages/core/agent_framework/_skills.py:492-528, 77-107, 261-296, 2360-2371` |
| MCP transport ABC | `class MCPTool` with `@abstractmethod def get_mcp_client(self)` — transport seam is its own ABC | `sources/agent-framework/python/packages/core/agent_framework/_mcp.py:1559-1566` |
| Workflow Runner Context Protocol | `@runtime_checkable class RunnerContext(Protocol)` covering messaging, events, checkpointing, streaming | `sources/agent-framework/python/packages/core/agent_framework/_workflows/_runner_context.py:97-280` |
| Exception hierarchy | Single base `AgentFrameworkException` with `inner_exception` and configurable `log_level`; subclasses by domain | `sources/agent-framework/python/packages/core/agent_framework/exceptions.py:15-262` |
| TypedDict option validation | `validate_chat_options` enforces numeric ranges on `temperature`/`top_p`/`frequency_penalty`/`presence_penalty`/`max_tokens` | `sources/agent-framework/python/packages/core/agent_framework/_types.py:3428-3482` |
| Tool-mode validation | `validate_tool_mode` rejects malformed `tool_choice` dicts (raises `ContentError`) | `sources/agent-framework/python/packages/core/agent_framework/_types.py:3569-3603` |
| ChatOptions TypedDict (structural contract) | `_ChatOptionsBase(TypedDict, total=False)` enumerates common provider options; provider-specific TypedDicts extend it | `sources/agent-framework/python/packages/core/agent_framework/_types.py:3352-3422` |
| ResponseStream (streaming wire contract) | `class ResponseStream(AsyncIterable[UpdateT], Generic[UpdateT, FinalT])` with finalizer hook | `sources/agent-framework/python/packages/core/agent_framework/_types.py:2939` |
| Provider lazy-loading seam | `agent_framework/openai/__init__.py` lazily `importlib.import_module`s providers; raises `ModuleNotFoundError` with install hint when missing | `sources/agent-framework/python/packages/core/agent_framework/openai/__init__.py:33-43` |
| Duck-typed test fixture (replaces MockChatClient) | Plain class implementing only `get_response` + `additional_properties` (no inheritance) is asserted against `SupportsChatGetResponse` | `sources/agent-framework/python/packages/core/tests/core/conftest.py:82-133` |
| Duck-typed MockAgent fixture | `class MockAgent(SupportsAgentRun)` exercises `isinstance(client, SupportsChatGetResponse)` and is invoked through generic test cases | `sources/agent-framework/python/packages/core/tests/core/test_clients.py:33-51`, `conftest.py:273-320` |
| Plain duck-typed embedding client | `class PlainEmbeddingClient` with no inheritance satisfies `SupportsGetEmbeddings`; wrong class does not | `sources/agent-framework/python/packages/core/tests/core/test_embedding_client.py:81-104` |
| Conformance-style lifecycle tests | `test_agent_type`, `test_agent_run`, `test_agent_run_streaming`, `test_agent_create_session` all parameterized against `SupportsAgentRun` | `sources/agent-framework/python/packages/core/tests/core/test_agents.py:128-152, 2273-2316` |
| Agent rejection of stray kwargs | `Agent.__init__` raises `TypeError` when given unknown kwarg (e.g., `legacy_key`) — compile-time guard at runtime | `sources/agent-framework/python/packages/core/tests/core/test_agents.py:223-225` |
| Cross-provider .NET conformance suite | `AgentConformance.IntegrationTests` with `AgentTests<TAgentFixture>` base, `RunTests`, `RunStreamingTests`, `StructuredOutputRunTests` | `sources/agent-framework/dotnet/tests/AgentConformance.IntegrationTests/AgentTests.cs:13`, `RunTests.cs:1-60` |
| .NET agent abstract base | `public abstract partial class AIAgent` with `RunCoreAsync` / `RunCoreStreamingAsync` / `CreateSessionCoreAsync` / `SerializeSessionCoreAsync` / `DeserializeSessionCoreAsync` abstract methods | `sources/agent-framework/dotnet/src/Microsoft.Agents.AI.Abstractions/AIAgent.cs:38-235, 366-370, 502-506` |
| .NET decorator-pattern base | `public abstract class DelegatingAIAgent : AIAgent` forwarding all operations to `InnerAgent` | `sources/agent-framework/dotnet/src/Microsoft.Agents.AI.Abstractions/DelegatingAIAgent.cs:28-101` |
| .NET AIContextProvider two-phase lifecycle | Abstract `AIContextProvider` with default `InvokingCoreAsync` (filter + merge + source-stamp) and `InvokedCoreAsync` (skip-on-error + filter + store); `ProvideAIContextAsync` and `StoreAIContextAsync` are the override seams | `sources/agent-framework/dotnet/src/Microsoft.Agents.AI.Abstractions/AIContextProvider.cs:42-419` |
| .NET ChatHistoryProvider | Abstract base with default `InvokingCoreAsync` (chronology + source-stamp) and `InvokedCoreAsync` (skip-on-error + filter + store); `ProvideChatHistoryAsync` and `StoreChatHistoryAsync` override seams | `sources/agent-framework/dotnet/src/Microsoft.Agents.AI.Abstractions/ChatHistoryProvider.cs:51-477` |
| .NET AgentSession contract | `abstract class AgentSession` with `StateBag`, `GetService`; docstring explicitly says sessions are not portable across agents | `sources/agent-framework/dotnet/src/Microsoft.Agents.AI.Abstractions/AgentSession.cs:11-119` |
| .NET workflow Executor ABC | `public abstract class Executor : IIdentified`; subclasses override `ConfigureProtocol(ProtocolBuilder)`; `Executor<TInput>` and `Executor<TInput,TOutput>` with abstract `HandleAsync` | `sources/agent-framework/dotnet/src/Microsoft.Agents.AI.Workflows/Executor.cs:164-422` |
| .NET workflow context interface | `public interface IWorkflowContext` plus `IExternalRequestContext`, `IExternalRequestEnvelope`, `IMessageRouter`, `IResettableExecutor`, `IIdentified`, `IWorkflowExecutionEnvironment` | `sources/agent-framework/dotnet/src/Microsoft.Agents.AI.Workflows/IWorkflowContext.cs:13`, `IExternalRequestEnvelope.cs:30`, `IMessageRouter.cs`, `IResettableExecutor.cs:11`, `IIdentified.cs:8`, `IWorkflowExecutionEnvironment.cs:12` |
| .NET checkpoint seam | `public interface ICheckpointStore<TStoreObject>`, `public interface IWireMarshaller<TWireContainer>` | `sources/agent-framework/dotnet/src/Microsoft.Agents.AI.Workflows/Checkpointing/ICheckpointStore.cs:12`, `IWireMarshaller.cs:11` |
| .NET chat-client decorator stack | `ChatClientAgent` calls `IChatClient.GetResponseAsync`; `FunctionInvocationDelegatingAgent`, `LoggingAgent`, `OpenTelemetryAgent`, `PerServiceCallChatHistoryPersistingChatClient` wrap it | `sources/agent-framework/dotnet/src/Microsoft.Agents.AI/ChatClient/ChatClientAgent.cs:205-263`, `dotnet/src/Microsoft.Agents.AI/FunctionInvocationDelegatingAgent.cs`, `LoggingAgent.cs`, `OpenTelemetryAgent.cs`, `dotnet/src/Microsoft.Agents.AI/ChatClient/PerServiceCallChatHistoryPersistingChatClient.cs` |
| JSON Schema for durable agent entity state | `durable-agent-entity-state.json` Draft 2020-12 schema with `$defs` for `usage`, `dataContent`, `errorContent`, `hostedFileContent` and required `schemaVersion` | `sources/agent-framework/schemas/durable-agent-entity-state.json:1-217` |
| SerializationMixin recursive typing | `to_dict` walks nested `SerializationProtocol`, lists, and dicts; uses `ClassVar` `DEFAULT_EXCLUDE` to omit private state | `sources/agent-framework/python/packages/core/agent_framework/_serialization.py:289-369` |
| Pyright config | `pyrightconfig.tests.json` enables stricter type-checking on tests; `pyrefly.toml`, `ty.samples.toml`, `pyrefly.samples.toml` committed | `sources/agent-framework/python/pyrightconfig.tests.json`, `pyrefly.toml`, `ty.samples.toml`, `pyrefly.samples.toml` |

## Answers to Dimension Questions

**1. Are interfaces small, coherent, and owned by the consumer side?**

Yes for the narrow extension points; mixed for the higher-level abstractions. The truly consumer-facing surfaces — `SupportsAgentRun` (3 properties + 2 methods + 1 optional), `SupportsChatGetResponse` (1 dict + 1 method), `SupportsGetEmbeddings` (1 method), the `Tool*` capability protocols (1 method each), `TokenizerProtocol` (1 method), `CompactionStrategy` (1 async callable), `SkillScriptRunner` (1 `__call__`), and `SerializationProtocol` (1 method + 1 classmethod) — are uniformly one-method-or-two-method Protocols with explicit ownership ("classes don't need to inherit from this protocol" in `_clients.py:93-95`, `_agents.py:195-198`). The wider `RunnerContext` (`_workflows/_runner_context.py:98-280`) is intentionally larger because the workflow runner needs messaging, events, checkpointing, and streaming in one object — but each method is small and documented. The .NET `AIAgent` is much larger (create/serialize/deserialize/run/runStreaming + GetService), justified because `AgentSession` is a separate class and the contract is intentionally clustered.

Ownership direction is consistently consumer-side: `SupportsChatGetResponse` lives in the core, and providers (OpenAI, Anthropic, etc.) implement the protocol via `BaseChatClient`; `SupportsAgentRun` lives in the core, and concrete agents implement it. The .NET `AIAgent` similarly owns session lifecycle, serialization, and service lookup at the base class, with subclasses only overriding the protected `*CoreAsync` hooks.

**2. Do contracts specify behavior, not just method signatures?**

Mostly yes on the .NET side via XML docs, partially yes on the Python side. .NET `AIContextProvider` XML docs spell out the **two-phase lifecycle** (InvokingCoreAsync filter-then-merge-then-source-stamp; InvokedCoreAsync skip-on-exception-then-filter-then-store) and embed **security considerations** as `<para>` blocks (`AIContextProvider.cs:33-40, 99-114, 215-220, 326-330`). `ChatHistoryProvider` documents chronology, source-stamp merging, and the default filters (`ChatHistoryProvider.cs:25-49, 155-187`). `AIAgent` documents that sessions are not portable and serialized sessions must be treated as untrusted (`AIAgent.cs:24-36, 178-220`, `AgentSession.cs:11-58`).

Python carries the same intent but in a more compact form. `BaseChatClient._validate_options` documents "Subclasses should call this at the start of _inner_get_response to validate options" (`_clients.py:326-340`), and `validate_chat_options` enforces numeric bounds (`_types.py:3428-3482`). The MCP layer explicitly specifies `MCPTaskOptions` semantics including max-wait, cancel-on-abandonment, and resubmit-vs-track reconnect policy. Behavioral contract that is **not** encoded anywhere includes: ordering guarantees between `Middleware.process` and `call_next`, what happens if a `ContextProvider.before_run` raises, what state an `AgentSession` must round-trip across processes (only the .NET side has a schema for that — `schemas/durable-agent-entity-state.json`), and lifecycle of in-flight `ResponseStream` when its consumer disconnects.

**3. Can providers, tools, stores, and runtimes be replaced safely?**

Yes for the contract surface; partial for hidden behaviors. Independently written test classes (`MockChatClient`, `PlainEmbeddingClient`, `MockAgent`) prove that a class with only the required methods + `additional_properties` (for chat) is treated identically to a framework subclass. Any provider package that subclasses `BaseChatClient` and implements `_inner_get_response` is automatically wrapped by `FunctionInvocationLayer`, `ChatMiddlewareLayer`, and `ChatTelemetryLayer` when constructed via the public `BaseChatClient.__init_subclass__` + decorator stack (the test `MockBaseChatClient` in `conftest.py:136-177` shows this).

Tools can be replaced safely when they expose `FunctionTool` (the most constrained form) and when an MCP tool is wrapped properly — `MCPTool._prepare_call_kwargs` enforces an **allowlist** built from the schema's `properties` plus user-supplied `additional_tool_argument_names`, and `_MCP_FRAMEWORK_DENYLIST` is a safety net for framework-named params that a server declares in its schema. The .NET side uses `IChatClient` from `Microsoft.Extensions.AI` as the provider seam — every Microsoft.Extensions.AI chat client works with `ChatClientAgent`.

Where replacement is fragile:
- `_inner_get_response` requires the implementor to handle **both** streaming and non-streaming via a single method signature returning either an awaitable or a `ResponseStream` (`_clients.py:413-431`); an alternative runtime that always returns a streaming iterator would not satisfy this without adaptation.
- `require_per_service_call_history_persistence` (`_agents.py:678, 470-475`) shifts history persistence from a per-invocation to a per-service-call middleware, with warnings logged for providers that have `load_messages=True` — this is a behavioral contract that an alternative store would need to honor.
- Streaming consumers must call `await stream.get_final_response()` to read the assembled response; nothing in the protocol prevents a consumer from discarding the stream early, leaving the underlying service call dangling.

**4. Are compatibility failures caught early by tests or validation?**

Yes on multiple layers, with gaps. Compile-time: `pyrightconfig.tests.json`, `pyrefly.toml`, and `ty.samples.toml` are committed; `OptionsContraT` / `EmbeddingProtocolOptionsT` `TypeVar`s push provider options through `@overload` declarations so a provider-specific TypedDict mismatch is caught statically (`_clients.py:70-178`, `_clients.py:856-911`). Runtime validation: `validate_chat_options` (`_types.py:3428-3482`), `validate_tool_mode` (`_types.py:3569-3603`), and `validate_tools` (`_types.py:3518-3566`); the unknown-kwarg rejection in `Agent.__init__` (`test_agents.py:223-225`).

Contract-shape conformance tests: `test_embedding_client.py:76-104` (mock + plain + wrong-class), `test_clients.py:33-51` (chat client + chat client base), `test_agents.py:128-152` (agent type / run / streaming). Cross-provider .NET conformance: `AgentConformance.IntegrationTests/RunTests.cs:17-53` runs `RunWithNoMessageDoesNotFailAsync`, `RunWithStringReturnsExpectedResultAsync`, `RunWithChatMessageReturnsExpectedResultAsync` against any `IAgentFixture` — covering at least Azure AI Foundry, Foundry Local, GitHub Copilot, AzureAI.Persistent, OpenAI, Anthropic, Bedrock, CosmosNoSql, Mem0, OpenAIAssistant, OpenAIResponse per the `dotnet/tests/` directory listing.

Gaps:
- No Python equivalent of `AgentConformance.IntegrationTests` — protocol conformance is verified only by individual mock tests, not by a parameterized harness that runs every provider against the same assertions.
- `validate_chat_options` does not validate provider-specific options (e.g., `OpenAIChatOptions.reasoning_effort`, `tools[].strict`). Each provider subclass must validate its own options; `BaseChatClient._validate_options` only calls the common validator.
- `SerializationMixin.to_dict` ignores fields listed in `DEFAULT_EXCLUDE` and silently drops them — there is no test that confirms every serialization round-trip is lossless for every subclass.
- `SkillsProvider`, `AgentLoopMiddleware`, `ToolApprovalMiddleware` carry no explicit contract test asserting that two independent implementations produce equivalent observable behavior.

## Architectural Decisions

- **Protocol + ABC split (Python).** `@runtime_checkable Protocol` for surfaces that any third-party class can implement (`SupportsAgentRun`, `SupportsChatGetResponse`, `SupportsGetEmbeddings`, `Tool*` capability protocols, `TokenizerProtocol`, `CompactionStrategy`, `SkillScriptRunner`, `SerializationProtocol`, `RunnerContext`); `ABC` for owned extension points where the framework provides default behavior subclasses are expected to inherit (`BaseChatClient`, `BaseEmbeddingClient`, `ContextProvider`, `HistoryProvider`, `Skill`, `SkillResource`, `SkillScript`, `AgentMiddleware`, `FunctionMiddleware`, `ChatMiddleware`, `MCPTool.get_mcp_client`). The split matches Python's structural-vs-nominative typing reality.

- **Decorator-stack composability.** Every layer that can be added to a chat client is a separate base class that takes the inner client in `__init__` and forwards calls — `FunctionInvocationLayer`, `ChatMiddlewareLayer`, `ChatTelemetryLayer`, `BaseChatClient`. The .NET equivalent is `DelegatingAIAgent` (`DelegatingAIAgent.cs:28-101`), `FunctionInvocationDelegatingAgent`, `LoggingAgent`, `OpenTelemetryAgent`. This means observability, function-calling, and middleware are properties of the stack, not the implementor, and adding a new layer never requires touching providers.

- **`*CoreAsync` protected hooks (Python `BaseChatClient._inner_get_response`; .NET `AIAgent.RunCoreAsync`/`RunCoreStreamingAsync`).** Public methods do argument validation, session construction, and lifecycle; `*Core` hooks do the actual work. This is the **template-method pattern** and lets the framework enforce invariants (security considerations, async-context propagation, session-creation ordering) without each provider re-implementing them.

- **Protocol-side `isinstance` for runtime gating.** `BaseAgent.as_tool` rejects objects that don't satisfy `SupportsAgentRun` at runtime via `isinstance(self, SupportsAgentRun)` (`_agents.py:530-532`). This makes the Protocol a real contract, not just a type-checker hint.

- **TypedDict options + `@overload`.** `ChatOptions` is a `TypedDict`, not a class. Provider packages extend it (`OpenAIChatOptions`, `BedrockChatOptions`, etc.); `get_response` carries four `@overload`s so `OptionsContraT` and the response-model type are captured statically (`_clients.py:128-197`). The .NET side uses `ChatOptions` similarly but as a record.

- **Provider lazy-loading.** `agent_framework/openai/__init__.py:33-43` `importlib.import_module`s providers on first access and raises a friendly `ModuleNotFoundError` with the install command when the provider package is missing. This decouples the core from optional provider dependencies.

- **Two-phase context lifecycle with default filtering.** Both `AIContextProvider` (`AIContextProvider.cs:42-419`) and `ChatHistoryProvider` (`ChatHistoryProvider.cs:51-477`) implement `InvokingAsync` (call `InvokingCoreAsync`) / `InvokedAsync` (call `InvokedCoreAsync`) where the public method delegates to a virtual `*Core*` that performs default filtering, merging, and source-stamping. Subclasses override `ProvideAIContextAsync` / `StoreAIContextAsync` for the simple case, or `*Core*` for full control. Python mirrors this in `ContextProvider.before_run` / `after_run` (`_sessions.py:351-411`).

- **`ResponseStream` with finalizer hook.** Streaming responses are first-class objects with both `__aiter__` and a finalizer callback (`_types.py:2939`). This unified shape lets non-streaming and streaming paths share downstream code while keeping the `isinstance(SupportsChatGetResponse)` check meaningful.

- **Storage contracts by configuration, not by type.** `HistoryProvider` has `load_messages`, `store_inputs`, `store_context_messages`, `store_context_from`, `store_outputs` (`_sessions.py:434-459`) so the same base class serves as primary memory, audit-only logging, or evaluation storage — fewer subclasses to keep in sync.

- **Security consideration as a documented contract.** Every state-touching .NET abstraction (`AIAgent`, `AIContextProvider`, `ChatHistoryProvider`, `AgentSession`) embeds a `<para>Security considerations</para>` block describing prompt-injection, PII, and untrusted-output risks. These are part of the contract, not just docs.

## Notable Patterns

- **Two-method Protocol with optional kwargs.** Most Protocols follow the shape `{property_or_field, async_method(self, *, options=None, **kwargs)}` — explicit options, additional-properties dict, and `**kwargs` for level-specific extensions (`_clients.py:81-197`, `_clients.py:864-911`). The kwargs is a pragmatic concession that lets providers add behavior without breaking the protocol.

- **`additional_properties` as the escape hatch.** Every serializable object carries an `additional_properties: dict[str, Any]` field (`_clients.py:127`, `_agents.py:409`, `_types.py:469-535`, `SerializationMixin.__init__`). This is the "vendor extensions" pattern — JSON-Schema's `additionalProperties` carried into the type system so user code can decorate objects without subclassing.

- **Consumer-defined duck typing proven by tests.** `PlainEmbeddingClient` and `MockChatClient` write `additional_properties: dict = {}` and the right `async def get_*`, and pass `isinstance(obj, SupportsGetEmbeddings)` (`test_embedding_client.py:84-104`, `conftest.py:82-133`). This is the canonical demonstration of the Protocol contract.

- **Stacked layer base classes** (`FunctionInvocationLayer[OptionsCoT]`, `ChatMiddlewareLayer[OptionsCoT]`, `ChatTelemetryLayer[OptionsCoT]`, `BaseChatClient[OptionsCoT]`) — the `MockBaseChatClient` in conftest (`conftest.py:136-177`) inherits from all four to show that a full-featured client is a one-line composition.

- **Pre-handler caching + post-handler transformation.** `BaseChatClient._finalize_response_updates` and `_build_response_stream` (`_clients.py:339-360`) are reusable helpers for any subclass that wants consistent streaming/non-streaming response shaping.

- **`get_service` for cross-component introspection.** Both `AIAgent` (`AIAgent.cs:118-136`) and `AIContextProvider` (`AIContextProvider.cs:344-362`) implement a Microsoft.Extensions.DependencyInjection-style `GetService(Type)` / `GetService<T>(...)`. This is the seams-out composition pattern: a decorator can be queried for the underlying agent without exposing it via property.

- **Override-either-Core-or-leaf hook.** `InvokingAsync` (public, delegates to `InvokingCoreAsync`) → `InvokingCoreAsync` (virtual, default behavior) → `ProvideAIContextAsync` (virtual, override seam). Most users override only the leaf, but full control is available.

## Tradeoffs

- **Protocols with `**kwargs` accept any caller behavior.** The chat-client `get_response` accepts `compaction_strategy`, `tokenizer`, `function_invocation_kwargs`, `client_kwargs` (`_clients.py:168-178`), and most other Protocol methods accept `**kwargs`. This keeps the protocol stable across providers but means **semantic compatibility must be tested, not just signature compatibility** — `kwargs` content is unchecked.

- **`additional_properties` everywhere.** Every framework object carries it. This is a feature (escape hatch) but also a leak (callers can store anything, and serialization round-trips don't validate it).

- **Decorator stack vs. single mega-class.** Stacking `FunctionInvocationLayer` + `ChatMiddlewareLayer` + `ChatTelemetryLayer` + `BaseChatClient` keeps responsibilities clean but requires order-of-inheritance discipline; the `MockBaseChatClient` in `conftest.py:142-177` shows this. A subclass that inherits `BaseChatClient` without `FunctionInvocationLayer` does not get function invocation.

- **Runtime `isinstance` for Protocol gating.** `BaseAgent.as_tool` uses `isinstance(self, SupportsAgentRun)` (`_agents.py:531`). `Protocol` `isinstance` only checks **method presence**, not behavioral correctness. A class with the right method names but wrong semantics will pass.

- **Python static vs. runtime validation gap.** `validate_chat_options` only checks numeric ranges and tool normalization (`_types.py:3428-3482`); provider-specific options (e.g., `OpenAIChatOptions.reasoning_effort` value range) are validated by the provider subclass, and any subclass that forgets the call to `_validate_options` skips the common validation entirely.

- **Generic `OptionsContraT` propagation.** Provider options are typed via `TypedDict` extensions and `@overload`s. Powerful, but the TypeVar inference relies on `default=` syntax (Python 3.13+) or `typing_extensions.TypeVar` (3.10+). Projects pinned to older Python may see degraded type-checking.

- **Streaming `ResponseStream` must be iterated to completion.** `await stream.get_final_response()` consumes the underlying generator. A consumer that forgets to await will leave the underlying HTTP connection open. The .NET equivalent uses `IAsyncEnumerable<AgentResponseUpdate>` which has the same risk.

- **Two-phase lifecycle creates a "what if both before_run raise" question.** `ContextProvider.before_run` failure mode is documented only in the docstring (`_sessions.py:370-410`); behavior in error case is not enforced by tests.

## Failure Modes / Edge Cases

- **`MCPTool` tool calling without `properties` schema.** Documented behavior: if the server returns `additionalProperties: true`, the framework forwards only user-configured extras and strips framework runtime kwargs (`_mcp.py` per `AGENTS.md` entry; security-relevant).

- **Long-running MCP tasks.** `MCPTaskOptions.max_task_wait` bounds polling + result fetching, and the abandon-vs-terminal distinction (`_MCPTaskAbandoned` marker) decides whether to fire `tasks/cancel` first. Unparseable success response after the server accepted the augmented call **does not** fall back, raising `ToolExecutionException` to avoid double-execution.

- **HistoryProvider with `require_per_service_call_history_persistence=True` and a service-managed session.** The per-service-call middleware owns persistence; the once-per-run path skips `HistoryProvider.after_run` (`_agents.py:470-475`). If both flags are set and the client does not store history server-side, the middleware also loads providers around each model call and drives the function loop with a local conversation.

- **Session portability across agents.** Documented as not supported (`AIAgent.cs:35-37`, `AgentSession.cs:35-37`). A test that passes a session from one agent to another is not part of the conformance suite.

- **Serialized session from untrusted source.** Documented as equivalent to accepting untrusted input (`AIAgent.cs:216-220`, `AgentSession.cs:46-53`). `Deserialized` JSON could carry elevated-trust roles or adversarial content; the framework does not validate.

- **Mock fixture violating Protocol without `isinstance` failing.** `MockAgent` declares `id`, `name`, `description`, `run`, `create_session` but does not declare `additional_properties: dict` explicitly — `Protocol` `isinstance` checks method presence only, so it passes (`conftest.py:273-320`). Real implementors need to remember the field.

- **`ResponseStream` finalized multiple times.** The class permits `await stream.get_final_response()` to be awaited multiple times if the implementation does not track the consumed state. Tests do not assert single-finalize behavior.

- **`@runtime_checkable` Protocols only check method presence.** A class with `async def run(self) -> None:` (returning `None`) satisfies `SupportsAgentRun.run`'s signature-only check, even though `run` is overloaded to return `Awaitable[AgentResponse[Any]]` for non-streaming and `ResponseStream[AgentResponseUpdate, AgentResponse[Any]]` for streaming. Mypy/pyright catch this, runtime does not.

- **Decorator ordering matters.** If `ChatTelemetryLayer` is applied before `ChatMiddlewareLayer`, telemetry observes middleware behavior; if after, telemetry is outside middleware. Neither the Python nor the .NET layer ordering is enforced by static checks.

- **Sampling-confused-deputy.** MCP `sampling_callback` is deny-by-default; passing `lambda params: True` re-enables auto-approve. Default `_DEFAULT_SAMPLING_MAX_REQUESTS` and `_DEFAULT_SAMPLING_MAX_TOKENS` apply rate / size limits. Requests and denials log at WARNING; content is not logged (`AGENTS.md` entry).

## Future Considerations

- **Python conformance harness.** The .NET `AgentConformance.IntegrationTests` model could be ported to Python (parameterize test cases over `SupportsChatGetResponse`, `SupportsAgentRun`, and tool-protocol fixtures and run identical assertions against every provider package).

- **Provider-specific option validation as part of the protocol.** Currently `validate_chat_options` only checks common numeric ranges. Adding a per-provider option schema (`OpenAIChatOptions`, `BedrockChatOptions`) validated at `BaseChatClient._inner_get_response` time would catch silent typos earlier.

- **JSON-Schema for `ChatOptions`.** The durable-agent-entity-state schema (`schemas/durable-agent-entity-state.json`) is published; a schema for `ChatOptions` and `AgentRunOptions` would let non-Python consumers validate requests before sending them.

- **Lossless round-trip tests for `SerializationMixin`.** Currently `to_dict` / `from_dict` is exercised per-class but no project-wide invariant asserts that `obj == obj.from_dict(obj.to_dict())` for every public type.

- **Single-source-of-truth for security considerations.** .NET XML docs encode security contracts inline; Python encodes them in `AGENTS.md` and inline docstrings. Consolidating them into the type-level docstring would let tools surface them at the call site.

- **Behavioral contracts for streaming lifecycle.** `ResponseStream` lifecycle (when finalize is mandatory, what happens on consumer disconnect) is currently docstring-only. A test fixture that asserts no leaked HTTP connections when the consumer abandons a stream would harden the contract.

- **Async-context propagation contract.** `AIAgent.CurrentRunContext` (`AIAgent.cs:97-106`) propagates via `AsyncLocal<T>`. The Python side has `token=None` / `function_invocation_kwargs` propagation through explicit parameter passing but no equivalent `AsyncLocal` channel. A framework-wide `ContextVar` (or explicit context object) would let cross-cutting middleware avoid long parameter lists.

## Questions / Gaps

- **No evidence found** for a Python conformance suite that runs every chat-client provider (OpenAI, Anthropic, Bedrock, Gemini, GitHub Copilot, Foundry Local, Ollama, Mem0, Azure AI Foundry) through the same `isinstance(client, SupportsChatGetResponse)` plus call-shape assertions. The .NET equivalent is `dotnet/tests/AgentConformance.IntegrationTests/`.

- **No evidence found** for a contract test that asserts `obj.from_dict(obj.to_dict())` round-trips losslessly across all `SerializationMixin` subclasses. Subclass-level round-trip tests exist; cross-class invariant does not.

- **No evidence found** for streaming lifecycle contracts — what happens when `await stream.get_final_response()` is never called; what happens when the consumer disconnects mid-stream. Both Python (`_types.py:2939`) and .NET (`AIAgent.cs:464-479`) `RunStreamingAsync` implementations leave the answer in prose only.

- **No evidence found** for a documented compatibility matrix — which `BaseChatClient` decorator layers compose safely (the test `MockBaseChatClient` in `conftest.py:142` inherits all four, but no docstring enumerates the order constraints or the layers that must come before/after telemetry).

- **Possible gap** in `ChatHistoryProvider.InvokingCoreAsync` (`.NET`): the default implementation merges chat-history messages *before* caller-provided messages (`ChatHistoryProvider.cs:150-152`), but the Python equivalent `ContextProvider.before_run` does not specify an order. Two independent `ChatHistoryProvider` implementations might disagree about ordering.

- **Possible gap** in the Python `_validate_options` enforcement: subclasses of `BaseChatClient` are instructed to call `await self._validate_options(options)` at the start of `_inner_get_response` (`_clients.py:326-340`), but no static check enforces this. A subclass that forgets skips numeric-range validation.

- **Possible gap** in the Protocol `additional_properties` field: `MockAgent` in `conftest.py:273-320` does not declare `additional_properties`, yet `isinstance(MockAgent(), SupportsAgentRun)` passes at runtime. Real implementors must remember to add the field; the protocol docstring documents it but does not enforce it.

- **Possible gap** in `SerializationProtocol` for typed Pydantic models: the protocol's `to_dict(**kwargs)` and `from_dict(value, /, **kwargs)` signatures assume dict-keyed serialization, but pydantic-anchored types in `_types.py` (e.g., `Content`, `Message`) use `_SHALLOW_COPY_FIELDS` and `raw_representation`. An alternative serialization implementation that ignores `raw_representation` will lose fidelity.

---

Generated by `dimensions/24.02-interface-contract-design.md` against `agent-framework`.
