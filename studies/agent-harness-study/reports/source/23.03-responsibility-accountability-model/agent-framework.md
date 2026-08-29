# Source Analysis: agent-framework

## Responsibility and Accountability Model

### Source Info

| Field | Value |
|-------|-------|
| Name | agent-framework |
| Path | `studies/agent-harness-study/sources/agent-framework` |
| Language / Stack | Python + .NET (Microsoft Agent Framework) |
| Analyzed | 2026-08-29 |

## Summary

Microsoft Agent Framework (MAF) is a framework/library for building agents, not a deployed agent product. Its responsibility model is therefore **delegation-to-the-builder/operator**: the framework provides attribution primitives (message `author_name`, tool invocation telemetry, `function_approval_request`/`function_approval_response` handshake, session-persisted approval rules) and explicit developer-facing warnings that the builder is responsible for Responsible-AI mitigations, data-boundary control, and human-in-the-loop safeguards. Model outputs are not stamped with automatic disclaimers; tool executions are attributed by tool name / call_id / server_label in typed `Content` objects and OpenTelemetry spans; human approvals are recorded as first-class `Content` types persisted in session state and optionally in history providers, but there is no immutable, centralized audit ledger or formal RACI documented in-code. Accountability guidance lives in `TRANSPARENCY_FAQ.md`, `README.md` Important Notes, and ADR `0006-userapproval.md`, not as enforceable runtime policy.

## Rating

**Rating: 5 / 10**

Rationale: Attribution primitives for agents/tools/humans exist as typed interfaces with tests and telemetry (author_name on Message/AgentResponseUpdate, FunctionTool approval_mode, ToolApprovalMiddleware state, OTel spans). However the holistic “who is responsible” model is informal: delegated to user/operator via documentation disclaimers, without a codified responsibility matrix, without automatic output disclaimers, and without durable, tamper-evident accountability logs. Operational recording depends on opt-in session + history providers + external OTel backends rather than a built-in accountability subsystem. This matches rubric band “present but inconsistent, weakly documented, or fragile”.

## Evidence Collected

Every entry MUST include a file path with line numbers. Format: `path/to/file.ts:NN`.

| Area | Evidence | File:Line |
|------|----------|-----------|
| Policy attribution - user responsibility for third-party systems | `README.md:207-212` states Important Note: You are responsible for any usage and associated costs of Third-Party Systems; responsible for reviewing/testing and implementing Responsible AI mitigations (metaprompt, content filters) | `README.md:207` |
| Policy attribution - developer accountability guidance | `TRANSPARENCY_FAQ.md:49-50` Framework-Specific Limitations call out “Accountability and Transparency: establish clear accountability mechanisms…trace decision-making” and “Security & unintended consequences: keep human in the loop” | `TRANSPARENCY_FAQ.md:49` |
| Policy attribution - responsible development practices | `TRANSPARENCY_FAQ.md:75-77` lists Human Oversight, Agent Modularity safeguard agent, LLM Selection, Security Measures, Testing, Monitoring/OTel as developer responsibilities | `TRANSPARENCY_FAQ.md:75` |
| Policy attribution - disclaimer that framework evaluation != model evaluation | `TRANSPARENCY_FAQ.md:26-28` notes AI performance metrics depend on underlying LLM provider; developers should conduct application-specific evaluation | `TRANSPARENCY_FAQ.md:26` |
| Output disclaimers - absence | No code injecting automatic disclaimer text into `ChatResponse`/`AgentResponse` found; grep for `disclaimer` yields no runtime insertion. Output is `Content.type=text` without provenance stamp | `python/packages/core/agent_framework/_types.py:338` |
| Model output attribution - Message author_name | `Message` carries `author_name: str | None` (`_types.py:1731,1754-1755`) documented as name of author; `Content` does not auto-stamp model vs human beyond role | `python/packages/core/agent_framework/_types.py:1731` |
| Model output attribution - propagation | `RawAgent._parse_non_streaming_response` / streaming hooks propagate `author_name` from executor/agent name (`_agents.py:1065-1066`, `_agents.py:1396-1397`) and workflow executors set `author_name=executor_id` (`_workflows/_agent.py:505,555,602,617`) | `python/packages/core/agent_framework/_agents.py:1065` |
| Tool attribution records - FunctionTool metadata | `FunctionTool` stores `name`, `kind`, `approval_mode`, `invocation_count`, `tool_call_id` telemetry (`_tools.py:301-315,398,403`) and serializes via `to_json_schema_spec`/`to_dict` | `python/packages/core/agent_framework/_tools.py:301` |
| Tool attribution records - Content function_call/result | `Content` type union includes `function_call` (call_id,name,arguments) and `function_result` (call_id,result,items) with `from_function_call`/`from_function_result` factories (`_types.py:338-361,793-875`) | `python/packages/core/agent_framework/_types.py:338` |
| Tool attribution records - hosted tool boundary | `FunctionTool` hosted case stores `server_label` in `additional_properties` and ToolApproval differentiates rules by `server_label` (`_harness/_tool_approval.py:61-64,290-294`, `_types.py:1179`) | `python/packages/core/agent_framework/_harness/_tool_approval.py:61` |
| Tool attribution records - OTel spans | `FunctionTool.invoke` emits `gen_ai.tool.name`, `gen_ai.tool.call.id`, `gen_ai.tool.call.arguments`, `gen_ai.tool.call.result`, `agent_framework.function.name` spans (`_tools.py:717-778`, `observability.py:207-214`) plus error.type capture | `python/packages/core/agent_framework/observability.py:207` |
| Tool attribution - MCP allowlist attribution | MCP `_mcp.py:22-23` injects OTel trace context via propagator and filters tool arguments allowlist, preserving provenance | `python/packages/core/agent_framework/_mcp.py:22` |
| Human decision logs - Content handshake | `Content.from_function_approval_request` (id,function_call,user_input_request) and `from_function_approval_response` (approved,bool) plus `to_function_approval_response` helper (`_types.py:1212-1302`) model durable human decision payload | `python/packages/core/agent_framework/_types.py:1212` |
| Human decision logs - session-persisted approval state | `ToolApprovalRule` (tool_name, arguments, server_label) and `ToolApprovalState` (rules, queued_approval_requests, collected_approval_responses) with `to_dict`/`from_dict` persisted to `AgentSession.state[source_id]` (`_harness/_tool_approval.py:86-217,250-277`) | `python/packages/core/agent_framework/_harness/_tool_approval.py:86` |
| Human decision logs - middleware handling creates audit trail in thread | `ToolApprovalMiddleware.process` drains auto-approvable queue, pops queued requests, injects collected responses as `Message(role=user, contents=collected_approval_responses)`, saves state (`_harness/_tool_approval.py:371-407,541-544`) forming thread-level record | `python/packages/core/agent_framework/_harness/_tool_approval.py:371` |
| Human decision logs - checkpoint pending requests | `WorkflowCheckpoint.pending_request_info_events` serialized via `WorkflowEvent.request_info` round-trips through `Memory/File` checkpoint storage (`python/packages/core/tests/workflow/test_checkpoint.py:569-1306`, `python/packages/core/agent_framework/_workflows/_...`) | `python/packages/core/tests/workflow/test_checkpoint.py:569` |
| Human decision logs - standing approvals via metadata | `create_always_approve_tool_response` / `create_always_approve_tool_with_arguments_response` store `tool_approval.always_approve` metadata in `additional_properties` and cloned via `ToolApprovalMiddleware._handle_inbound_approval_response` (`_harness/_tool_approval.py:220-247,515-539`) | `python/packages/core/agent_framework/_harness/_tool_approval.py:220` |
| Accountability documentation - ADR for approvals | ADR `0006-userapproval.md:381-407` defines chosen option 5: base `UserInputRequestContent` + `FunctionApprovalRequestContent`/`FunctionApprovalResponseContent`, logs approvals for debugging/auditing but notes service-managed threads won't support long-term record | `docs/decisions/0006-userapproval.md:381` |
| Accountability documentation - security disclosure process | `SECURITY.md:7-18` defines coordinated vulnerability disclosure via MSRC, not operational accountability for agent actions | `SECURITY.md:7` |
| Accountability documentation - workflow attribution | ADR `foundry-toolbox` decision notes `agent_framework.foundry.toolbox.sources` OTel attribute for local attribution of run to toolbox, no server-side attribution yet (`docs/decisions/0025-foundry-toolbox-support.md:381-421`) | `docs/decisions/0025-foundry-toolbox-support.md:381` |
| Context attribution (secondary) | `ContextProvider` requires `source_id` and `additional_properties.attribution` marker for runtime filtering, stripped before storage (`docs/decisions/0016-python-context-middleware.md:419-463`, `python/packages/core/agent_framework/_sessions.py:source_id`) | `docs/decisions/0016-python-context-middleware.md:419` |
| Approvals enforcement in tools | `LocalShellTool` enforces `approval_mode=always_require` by default, requiring `acknowledge_unsafe=True` to disable (`python/packages/tools/agent_framework_tools/shell/_tool.py:125-159`, tests at `tests/test_security.py:99-113`) | `python/packages/tools/agent_framework_tools/shell/_tool.py:125` |

## Answers to Dimension Questions

**1. Who is responsible for agent actions?**
Framework delegates responsibility to **builder / operator / organization**, not to model/runtime/tool author per se. `README.md:207-212` Important Notes and `TRANSPARENCY_FAQ.md:75-87,106-124` place duty for Responsible-AI mitigations (content filters, metaprompt), data-boundary stewardship, tool selection, testing, and maintaining human oversight on the developer deploying the agent. Runtime does not assume liability; `TRANSPARENCY_FAQ.md:26-28` explicitly says evaluation and safety depend on chosen LLM provider and application implementation. Tool authors are responsible via `approval_mode` and invocation limits on `FunctionTool` (`_tools.py:302-402`), but framework does not enforce a centralized policy owner. No formal RACI enum or config key was found; 10+ greps for `responsib*` only hit docs, not code interfaces.

**2. Is model output attributed?**
Partially. Model-generated vs human vs tool messages are distinguished by `Message.role` and `Message.author_name` (`_types.py:1731`). Agent name is propagated to `author_name` when missing (`_agents.py:1065-1066,1396-1397`), workflows set `author_name=executor_id` (`_workflows/_agent.py:505`). However there is **no automatic disclaimer** or provenance watermark stamped onto text content; `Content.type=text` carries only `text` and `annotations` (`_types.py:594-609`). OTel spans record `gen_ai.agent.name`, `gen_ai.response.id`, usage tokens (`observability.py:231-245`), but that is telemetry-side attribution, not user-visible disclosure. Searches for `disclaimer` found no runtime injection.

**3. Are tool decisions attributed?**
Yes — tool invocations are strongly attributed at API and telemetry layers. `FunctionTool` carries `name`, `kind`, `approval_mode`, `invocation_count` (`_tools.py:369-407`) and emits `FunctionInvocationContext` with live mutable `tools` list. Execution results are normalized to `Content.from_function_result` with `call_id` linkage (`_types.py:817-874`). Tool calls themselves are `Content.from_function_call` with `call_id`/`name`/`arguments` (`_types.py:793-815`) plus `server_label` for hosted tools (`_harness/_tool_approval.py:61`). OTel captures `gen_ai.tool.name`, `gen_ai.tool.call.id`, arguments, result, duration (`_tools.py:717-779`, `observability.py:207-214`). Methodologically this answers “who/what executed which tool with which args” without guessing, but there is no higher-level policy-to-tool ownership registry beyond these per-invocation artifacts.

**4. Are human approvals recorded?**
Yes — approvals are first-class, recorded in three places: (a) as `function_approval_request` / `function_approval_response` `Content` objects threaded as `Message` contents (`_types.py:1212-1253`) with `approved: bool` and `function_call` echo; (b) as durable `ToolApprovalState` in `AgentSession.state[source_id]` with `rules`, `queued_approval_requests`, `collected_approval_responses` serialized to dict (`_harness/_tool_approval.py:159-277`); (c) as persisted workflow `pending_request_info_events` in checkpoints (`tests/workflow/test_checkpoint.py:569-706`). The middleware flow queuing/injecting/collecting ensures a thread-replayable record. However recording is **opt-in** (requires `AgentSession` + `ToolApprovalMiddleware` or workflow checkpoint storage); without a session the middleware throws `RuntimeError` (`_harness/_tool_approval.py:373`). “Always approve” standing rules are also persisted (`_harness/_tool_approval.py:220-247`). Tests cover approval_mode enforcement (`tests/test_security.py:99-113`) and workflow request_info (`tests/workflow/test_workflow_status.py:52-125`).

**5. Is accountability documented?**
Weakly/documentarily, not operationally. `TRANSPARENCY_FAQ.md:49-87` and `README.md:207-212` document expectations (traceability, human-in-the-loop, testing, monitoring) but do not define a codified accountability model or audit-log retention policy. ADR `0006-userapproval.md:381-407` documents approval content design and notes that while local threads retain `FunctionApprovalRequestContent`, **service-managed threads won’t** support that content type persistently (“log approvals so there is a trace for debugging/auditing” `0006-userapproval.md:95`). `SECURITY.md` covers vulnerability disclosure, not action audit. `docs/decisions/0025-foundry-toolbox-support.md:385-430` acknowledges incomplete server-side toolbox attribution. No evidence found of a durable, tamper-evident, queryable accountability ledger; accountability relies on compositable opt-in layers: session state, `HistoryProvider` (InMemory/File/Cosmos/etc.), and external OTel exporters (`observability.py:365-533`).

## Architectural Decisions

- **Framework vs. runtime responsibility split** (`README.md:207`, `TRANSPARENCY_FAQ.md:26`, `docs/decisions/0024-codeact-integration.md:24`): Framework owns approval capabilities and isolation configuration; backend/host owns isolation boundary. Consequence: portability across OpenAI/Foundry/Anthropic but accountability must be rebuilt per-deployment.

- **Approval as typed Content handshake** (`_types.py:1212-1253` + `docs/decisions/0006-userapproval.md:266-381`): Chose option 5 — base `UserInputRequestContent` with derived `FunctionApprovalRequestContent` / `FunctionApprovalResponseContent` (Python side as `function_approval_request`/`function_approval_response` string types). Enables suspend/resume across remote A2A and service-managed threads, unlike callback-only option 2.

- **Session-state–backed standing approvals** (`_harness/_tool_approval.py:86-162,250-277`): `ToolApprovalRule` with optional exact-argument matching and `server_label` scoping, persisted as dict in `AgentSession.state`. Allows “always approve this tool / this tool with these args on this server” without re-prompting.

- **Attribution via author_name + OTel semantic conventions** (`_types.py:1731`, `_agents.py:1065`, `observability.py:176-311`): Follow OTel GenAI conventions (`gen_ai.agent.name`, `gen_ai.tool.*`, `gen_ai.usage.*`) rather than custom headers; keeps model vs tool vs human distinguishable in traces even though not watermarked in chat text.

- **No server-side content persistence assumption for approvals** (`docs/decisions/0006-userapproval.md:92-95`): Acknowledges `FunctionApprovalRequestContent` may not be storable in service-managed thread types; therefore recommends logging rather than guaranteeing long-term history — trades durability for provider compatibility.

## Notable Patterns

- **Suspend/resume approval loop**: Agent returns `function_approval_request` content, caller collects user decision into `function_approval_response`, re-invokes agent with that response on same thread/session — implements human-in-the-loop as a workflow continuation, not a blocking callback (`docs/decisions/0006-userapproval.md:384-417`, `_harness/_tool_approval.py:371-406`).

- **Queued vs auto-approved triage**: `_harness/_tool_approval.py:560-606` batches approval requests; auto-approves those matching `rules` or `auto_approval_rules` callbacks (receiving `function_call` content), queues second+ unresolved requests, and strips auto-approved/queued ids from outbound messages — mirrors .NET bypass for mixed tool-call batches.

- **Server-label scoping for hosted tools**: Standing approvals keyed by `(tool_name, server_label, arguments)` prevents same-named tools on different MCP/hosted servers from sharing approvals (`_harness/_tool_approval.py:61-64,290-323`).

- **Opt-in accountability layers**: `InMemoryHistoryProvider` auto-injection only when session exists and no service-side storage (`_agents.py:1201-1213`); `ToolApprovalMiddleware` requires `AgentSession`; OTel export guarded by `OBSERVABILITY_SETTINGS.ENABLED` with sticky disable (`observability.py:798-829`). Pattern composes but makes “always accountable” non-default.

- **Context attribution via source_id**: `ContextProvider.source_id` plus `additional_properties.attribution` marker enables filtering ephemeral vs durable context (`docs/decisions/0016-python-context-middleware.md:446-463`), orthogonally to tool attribution.

## Tradeoffs

- **Builder-owns-accountability vs framework-owns-accountability**: Delegating to builder via docs/OTel keeps framework lightweight and provider-agnostic, but means two deployments with identical code can have radically different audit guarantees depending on which history provider / OTel exporter the builder wires. No enforcement of minimum audit retention.

- **Typed Content handshake vs stringly-typed tool args**: Strong typing and `server_label` scoping improves precise attribution and prevents confused-deputy approvals; cost is added complexity in session state serialization (`_harness/_tool_approval.py:67-84` `Content.from_dict` round-trips) and need for `Acknowledge_unsafe` opt-out for shell tools.

- **Session-state persistence vs service-managed threads**: Storing approvals in `AgentSession.state` (local) ensures round-trip fidelity for long-running workflows, but breaks when delegating to Foundry/hosted services that don’t understand `function_approval_request` content type — hence ADR’s logging-only fallback. Tradeoff sacrifices end-to-end audit for interoperability.

- **Telemetry as accountability surrogate**: Rich OTel spans (`observability.py:207-245`, `_tools.py:717-779`, `_mcp.py:22`) give queryable attribution in APM backends, but are sampling-configurable and not tamper-evident; disabling via `disable_instrumentation()` is sticky (`observability.py:1106-1129`), so accountability can be silently lost if an integration disables telemetry.

- **Auto-approval callbacks vs explicit human gate**: `auto_approval_rules: Sequence[ToolApprovalRuleCallback]` (`_harness/_tool_approval.py:358`) lets policies auto-approve low-risk calls (reducing friction), but heuristic callbacks receive only `function_call` content — no rich risk context — and run without requiring session-backed justification record beyond the resulting `function_approval_response`.

## Failure Modes / Edge Cases

- **No session → no accountability**: `ToolApprovalMiddleware.process` raises `RuntimeError("ToolApprovalMiddleware requires an AgentSession.")` (`_harness/_tool_approval.py:373`) and `_sessions.py` logic shows approvals not recorded; headless/autonomous runs skip human decision logging entirely.

- **Service-managed thread amnesia**: Per `docs/decisions/0006-userapproval.md:92-95`, approvals placed in a service-managed conversation may be dropped on persistence (service doesn’t support `FunctionApprovalRequestContent` type). Thread replay on that service loses the “why was this tool denied” provenance.

- **Ephemeral storage without HistoryProvider**: Stock `AgentSession.state` lives in-memory; process crash without `FileHistoryProvider`/`Cosmos`/`Redis` loses `ToolApprovalState` including standing “always approve” rules. Checkpoint tests show correct serialization (`test_checkpoint.py:569`) but that only helps workflow checkpoints, not ad-hoc agent runs.

- **Mixed approval batches and bypass fragility**: `_harness/_tool_approval.py:580-606` hides auto-approved and queued requests via `id(content)` set removal. If a caller clones `Content` (deepcopy changes identity), remove logic misses and duplicate prompts are emitted. Also, `len(unresolved) <=1` early-return path (`_harness/_tool_approval.py:579`) means single unresolved request is still surfaced — could starve multi-tool atomicity.

- **Standing rule poisoning**: `create_always_approve_tool_response` writes `additional_properties.tool_approval.always_approve` (`_harness/_tool_approval.py:241-246`); if a malicious tool echoes back crafted `function_approval_response` with that metadata, `_handle_inbound_approval_response` will create a permanent `ToolApprovalRule` for that tool/server, bypassing future human gates.

- **Invocation-budget overflow hides attribution**: `FunctionTool.max_invocations` / `max_invocation_exceptions` (`_tools.py:325-405`) throws `ToolException` after limit, but does not emit distinct `accountability` event; downstream logs show `error.type` but not that budget caused the denial vs policy.

- **Missing disclaimer enables misattribution**: Since text outputs carry no model watermark, downstream consumers (chat UI, A2A peer) cannot distinguish model hallucination from human-authored text without trusting `role`/`author_name`, which are caller-controllable (framework sets `author_name` only when `None`).

## Future Considerations

- **Formal responsibility matrix artifact**: Codify RACI (model vs runtime vs tool author vs operator vs organization) as a shipped decision record or `responsibility.yaml` rather than scattered `README` Important Notes + `TRANSPARENCY_FAQ` prose — would enable audit questionnaires to reference a versioned policy file.

- **Automatic provenance stamp**: Add optional `SystemInstructions` or `Message` header stamping generator (e.g., `model: gpt-4o via FoundryChatClient`) controlled by policy flag, plus text-level disclaimer hook in `AgentMiddlewareLayer` so builders can satisfy regulatory disclosure without hand-rolling middleware.

- **Durable, immutable approval ledger**: Promote `ToolApprovalState` from ephemeral session dict to pluggable `ApprovalLedger` interface with `File`/`OTel LogRecord`/`Cosmos` implementations that append-only and support query (“who approved tool X with args Y at time T”), addressing service-managed thread amnesia.

- **Server-side toolbox attribution contract**: Resolve open question in `docs/decisions/0025-foundry-toolbox-support.md:421-430` — negotiate custom request header/metadata for server-side request-log attribution back to toolbox version, so hosted runs are attributable even when client OTel is off.

- **Approval audit completeness**: Log rejected approvals as `function_result` with “Function invocation denied” already (`docs/decisions/0006-userapproval.md:400-403`) but promote to structured `error.type=approval_denied` OTel event and include `approved_by` (user id / tenant id from `ChannelIdentity`) when available — requires threading identity through `ApprovalResponse`.

## Questions / Gaps

- **Who signs accountability?** No `CODEOWNERS`-style mapping of tool `approval_mode` decisions to owners; `additional_properties` holds no owner contact field. Searched `additional_properties` usages — only `server_label` and `tool_approval` metadata found. Gap: cannot answer “who is accountable for allowing tool X to run without approval”.

- **Are human identities captured?** `ToolApprovalState` records `approved: bool` + `function_call`, but not `approved_by` user principal, timestamp, or justification. OTel spans likewise lack `user.id`. Cannot attribute a specific human to a specific approval beyond session id.

- **Cross-source compliance checks?** Per Hard Rule 1, sibling sources not inspected. Whether .NET implementation’s `Message`/`Content` mapping diverges (e.g., `ApprovalRequestContent` class shape) was not verified — could affect cross-language audit parity.

- **Regulatory disclaimer requirements?** No evidence of EU AI Act / AI transparency disclosure templating. `TRANSPARENCY_FAQ.md:38-41,48` mentions “Lack of Transparency” and “Content Harms” but points to external content moderation services, not in-framework disclosure.

---

Generated by `Dimension 23.03: Responsibility and Accountability Model` against `agent-framework`.
