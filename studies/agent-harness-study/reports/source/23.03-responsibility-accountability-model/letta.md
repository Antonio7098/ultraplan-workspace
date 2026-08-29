# Source Analysis: letta

## 23.03 Responsibility and Accountability Model

### Source Info

| Field | Value |
|-------|-------|
| Name | letta |
| Path | `studies/agent-harness-study/sources/letta` |
| Language / Stack | Python (Pydantic, SQLAlchemy, FastAPI) / PostgreSQL + pgvector |
| Analyzed | 2026-08-29 |

## Summary

Letta does not implement an explicit responsibility/accountability model that can answer "who is responsible for this action?". Responsibility is fragmented across three layers: (1) legal/organizational policy in `TERMS.md:15` ("You are fully responsible for Your Content") and contributor policy in `AI_POLICY.md:11` ("human-in-the-loop must fully understand all code"); (2) runtime technical attribution via low-level persistence — `Message` rows (`letta/orm/message.py:23`), `Step` rows (`letta/orm/step.py:20`), `Run` rows (`letta/orm/run.py:22`), and `Tool` metadata (`letta/schemas/tool.py:67`) — each carrying `created_by_id`, `organization_id`, `sender_id`, `model`, and `provider_name`; and (3) a mature human-in-the-loop approval subsystem (`ApprovalRequestMessage`/`ApprovalResponseMessage` in `letta/schemas/letta_message.py:306`/`328`) persisted to `messages.approvals` (`letta/orm/message.py:81`). No model-output disclaimer or watermark is emitted; assistant content (`letta/schemas/letta_message.py:350`) and reasoning (`letta/schemas/letta_message.py:156`) are returned without `AI-generated` attribution. Tool execution itself is attributed to the LLM-initiated `ToolCallMessage` (`letta/schemas/letta_message.py:222`) and server-side `ToolReturnMessage` (`letta/schemas/letta_message.py:279`) tied to a `step_id`/`run_id`, not to a runtime/tool-author identity the operator could query as an accountability view. Human approvals are durably recorded with per-tool-call granularity (`ApprovalReturn.tool_call_id` in `letta/schemas/letta_message.py:33`), but the system provides no unified accountability document, dashboard, or `GET /accountability` API explaining the chain model → runtime → tool author → user/operator/organization.

## Rating

**4 / 10 — Present but inconsistent, weakly documented, and fragile**

Rationale: Human approvals are fully persisted and validated with tests, and low-level attribution fields exist across ORM/pydantic layers, which lifts the score above "absent" (1-3). However, Letta fails the dimension's core question: there is no coherent, documented responsibility assignment for agent actions, no output disclaimer, no tool-authorship attribution surfaced at runtime, and no accountability interface that aggregates `actor`, `model`, `tool`, and `human decision` into a single observable answer. The attribution that does exist is ad-hoc (scattered `created_by_id`/`sender_id`/`model` columns) and not enforced as an invariant.

## Evidence Collected

Every entry MUST include a file path with line numbers. Format: `path/to/file.ts:NN`.

| Area | Evidence | File:Line |
|------|----------|-----------|
| Policy attribution | Terms assign responsibility to user: "You are fully responsible for Your Content" + "You are responsible for safeguarding your password... Letta will not be liable for your acts" | `TERMS.md:14-16` |
| Policy attribution | AI usage policy requires disclosure and human understanding: "All AI usage must be disclosed" + "human-in-the-loop must fully understand all code" | `AI_POLICY.md:7-14` |
| Tool attribution records | `Tool` schema carries provenance: `created_by_id`, `last_updated_by_id`, `metadata_`, `source_code`, `tool_type`, `default_requires_approval` | `letta/schemas/tool.py:67-69`, `letta/schemas/tool.py:59-61` |
| Tool attribution records | `Message` ORM stores `sender_id` (identity/agent that sent message), `tool_calls`, `tool_returns`, `model`, `step_id`, `run_id` | `letta/orm/message.py:59-65`, `letta/orm/message.py:44-53` |
| Tool attribution records | `LettaMessage` base carries `sender_id`, `step_id`, `run_id` for every user-facing message type | `letta/schemas/letta_message.py:97-101` |
| Tool attribution records | `Step` ORM captures execution provenance: `organization_id`, `agent_id`, `provider_name`, `provider_category`, `model`, `model_handle`, `trace_id` | `letta/orm/step.py:28-50` |
| Tool attribution records | `Run` ORM links agent interaction to organization: `organization_id`, `agent_id`, `conversation_id`, `stop_reason`, timing metrics | `letta/orm/run.py:36-70` |
| Human decision logs | Approval data model: `ApprovalCreate.approvals: List[LettaMessageReturnUnion]`, legacy `approve`/`approval_request_id` fields | `letta/schemas/message.py:178-197` |
| Human decision logs | `ApprovalReturn` discriminated union: `tool_call_id`, `approve: bool`, `reason` | `letta/schemas/letta_message.py:31-35` |
| Human decision logs | `Message` ORM persists approvals: `approvals: ApprovalsColumn`, `approve`, `denial_reason`, `approval_request_id` | `letta/orm/message.py:75-83` |
| Human decision logs | Custom column serialization for approvals preserves structured history: `ApprovalsColumn.process_bind_param` / `process_result_value` | `letta/orm/custom_columns.py:116-126` |
| Human decision logs | Approval helpers validate attribution: `validate_approval_tool_call_ids`, `validate_persisted_tool_call_ids`, converter `serialize_approvals`/`deserialize_approvals` | `letta/agents/helpers.py:121-134`, `letta/helpers/converters.py:289-353` |
| Human decision logs | Agent loop enforces HITL gate: `RequiresApprovalToolRule` category, `is_requires_approval_tool()`, `get_requires_approval_tools()`, `create_approval_request_message_from_llm_response` | `letta/helpers/tool_rule_solver.py:48-50`, `letta/helpers/tool_rule_solver.py:186-196`, `letta/agents/letta_agent_v3.py:1681-1709` |
| Human decision logs | Message role includes `approval` and approval messages carry explicit `step_id` binding; `to_letta_messages` routes `role=="approval"` to `ApprovalRequestMessage` vs `ApprovalResponseMessage` | `letta/schemas/message.py:322`, `letta/schemas/message.py:519-569` |
| Accountability documentation | No disclaimer/watermark code found — grep for `disclaimer`, `AI-generated`, `generated by` returns only policy/license boilerplate, zero runtime watermark | `TERMS.md:10`, `AI_POLICY.md:23-24` (no disclaimer impl) |
| Accountability documentation | Base schema formally tracks `created_by_id`/`last_updated_by_id` on all ORM entities via `OrmMetadataBase` | `letta/schemas/letta_base.py:100-103` |
| Accountability documentation | Log context propagation carries actor attribution into observability: `update_log_context(agent_id, actor_id)` and `LogContextFilter` injects `actor_id` | `letta/log_context.py:24-27`, `letta/log.py:153-175` |

## Answers to Dimension Questions

**1. Who is responsible for agent actions?**

Partially and implicitly. Legal docs assign end responsibility to the user/organization (`TERMS.md:14` "You are fully responsible for Your Content"; `TERMS.md:16` "You are responsible for safeguarding your password"; `TERMS.md:30` disclaimer of warranties + `$4.20` liability cap in `TERMS.md:33`). Contributor policy assigns responsibility to the human reviewer of AI-assisted code (`AI_POLICY.md:11-12`). At runtime, every persistence write is gated by an `actor: PydanticUser` parameter (e.g., `letta/agent.py:178`, `letta/agents/letta_agent_v3.py:209`, `letta/services/step_manager.py:44-47`), which scopes access via `organization_id` and is stored as `created_by_id`/`organization_id` on `Agent` (`letta/orm/agent.py:52`), `Message`, `Step`, `Run`, etc. However, there is no single "responsibility model" object or doc that maps the six principals (model, runtime, tool author, user, operator, organization) to action types. The system cannot answer "who is responsible for this tool call?" without manual correlation of `step.model` + `message.sender_id` + `tool.created_by_id` + `run.organization_id`.

**2. Is model output attributed?**

Weakly at the persistence layer, not at the presentation layer. Every `Message` row optionally records `model` (`letta/orm/message.py:44`, `letta/schemas/message.py:278`), and every `Step` records `model`, `model_handle`, `provider_name`, `provider_category`, `provider_id` (`letta/orm/step.py:34-48`). Tool-call and assistant messages surfaced as `ToolCallMessage` (`letta/schemas/letta_message.py:222`), `AssistantMessage` (`letta/schemas/letta_message.py:350`), `ReasoningMessage` (`letta/schemas/letta_message.py:156`) carry `id/date/sender_id/step_id/run_id` but no `model` field in the public API shape and no textual disclaimer. No watermark, no `AI-generated` prefix, no per-message `generated_by_model` label is injected into user-visible output. The attribution is therefore forensic (query `steps` table) rather than user-facing.

**3. Are tool decisions attributed?**

Partially. The decision to call a tool is attributed to the LLM via the persisted `tool_calls: List[OpenAIToolCall]` on `Message` (`letta/orm/message.py:46`, `letta/schemas/message.py:287-289`) surfaced as `ToolCallMessage.tool_calls` (`letta/schemas/letta_message.py:237`). Execution results are attributed via `tool_returns: List[ToolReturn]` (`letta/orm/message.py:55-57`) surfaced as `ToolReturnMessage` (`letta/schemas/letta_message.py:279-304`) with `status` (`success`/`error`), `stdout`, `stderr`. Tool authorship is attributed via `Tool.created_by_id`/`last_updated_by_id` and `tool_type` (`CUSTOM` vs `LETTA_CORE` vs `MCP` in `letta/schemas/tool.py:39-65`), and organization scoping via `Tool.project_id`/`organization`. Critically, there is no runtime `tool_attribution` record that says "this side-effect was performed by tool X authored by Y and executed by runtime Z" — attribution must be reconstructed by joining `Message.tool_calls` → `Agent.tools` → `Tool` metadata + `Step.provider_name`.

**4. Are human approvals recorded?**

Yes — this is the strongest area. Letta implements a first-class, persisted HITL flow:
- Tool-rule declaration: `RequiresApprovalToolRule(tool_name)` (`letta/schemas/tool_rule.py:348`) managed via `ToolRulesSolver.requires_approval_tool_rules` (`letta/helpers/tool_rule_solver.py:48`) with helpers `is_requires_approval_tool()` (`letta/helpers/tool_rule_solver.py:186`) and `get_requires_approval_tools()` (`letta/helpers/tool_rule_solver.py:194`).
- Request creation: When a tool requiring approval is invoked, the agent loop does not execute it; instead it persists `ApprovalRequestMessage` (`letta/schemas/letta_message.py:306`) via `create_approval_request_message_from_llm_response` (`letta/agents/letta_agent_v3.py:1698`) and returns `StopReasonType.requires_approval` (`letta/agents/letta_agent_v3.py:1709`). Legacy single-call and modern multi-call paths are both handled (`letta/agents/letta_agent_v3.py:973-1005`).
- Response persistence: Human decision arrives as `ApprovalCreate` (`letta/schemas/message.py:178`) containing `approvals: List[ApprovalReturn|ToolReturn]` (`letta/schemas/message.py:182`) each with `tool_call_id`, `approve`, `reason` (`letta/schemas/letta_message.py:33-35`), or client-side `ToolReturn`s. Persisted to `messages.approvals` via `ApprovalsColumn` (`letta/orm/message.py:81-83`, `letta/orm/custom_columns.py:116`) with serialization handling null bytes (`letta/helpers/converters.py:289-311`).
- Validation & observability: Helpers validate tool-call-id integrity (`letta/agents/helpers.py:102-134`, plus warning on malformed approval in `letta/agents/letta_agent_v3.py:1010-1013`), and `Agent.pending_approval` exposes the current approval gate (`letta/orm/agent.py:339-458`, `letta/schemas/agent.py:134`). Integration test asserts `stop_reason == "requires_approval"` and `ApprovalRequestMessage` presence (`tests/integration_test_multi_modal_tool_returns.py:125-134`).

**5. Is accountability documented?**

No. Beyond `TERMS.md`/`PRIVACY.md`/`SECURITY.md` (legal/privacy, not runtime accountability) and `AI_POLICY.md` (contributor governance), the codebase contains no `accountability.md`, no architecture doc describing the responsibility chain, and no operational guide for answering "who is responsible for this action?". No API endpoint, view, or CLI command aggregates actor/model/tool/human-decision lineage. Logging does propagate `actor_id`/`agent_id`/`org_id` (`letta/log_context.py:24`, `letta/log.py:153`, `letta/otel/tracing.py` spans in `letta/agents/letta_agent_v2.py:919-945`), and `Step.trace_id`/`request_id` (`letta/orm/step.py:74-77`, `letta/services/step_manager.py:157-158`) provide forensic hooks, but without documented accountability semantics the operator must infer responsibility from scattered columns.

## Architectural Decisions

- **Actor-gated persistence (positive, but implicit):** All manager writes require `actor: PydanticUser` (`letta/services/agent_manager.py:73`, `letta/services/step_manager.py:44`, `letta/agent.py:178`), enforcing organization-scoped access and provenance. Decision: correct for multi-tenant isolation, but never framed as a responsibility model — just access control. File: `letta/schemas/letta_base.py:100` cementing `created_by_id` everywhere.
- **Message-as-ledger for accountability:** All actions (user, assistant, tool, approval) are reified as `Message` rows with `role` discriminator (`letta/schemas/message.py:319-324`) and `step_id`/`run_id` linkage (`letta/orm/message.py:48-53`). Decision: gives durable auditability but mixes distinct concerns (chat history + audit log + HITL state) in one table, complicating responsibility queries.
- **Discriminated-union approval model:** `LettaMessageUnion` (`letta/schemas/letta_message.py:463`) and `LettaMessageReturnUnion` (`letta/schemas/letta_message.py:51`) cleanly separate `ApprovalReturn` vs `ToolReturn`, allowing multi-tool approval in one round trip (`letta/agents/letta_agent_v3.py:981-1005`). Decision: enables fine-grained per-tool-call accountability with reason codes.
- **No output attribution by design:** Public `LettaMessage` shapes (`AssistantMessage`, `ReasoningMessage`, `ToolReturnMessage`) omit model identity from the wire format; model is only in internal `Message.model`/`Step.model` (`letta/orm/step.py:45`). Decision: keeps API clean but violates accountability principle that model output should be self-attributing.
- **Legal liability minimization over operational accountability:** `TERMS.md:30-33` disclaims warranties and caps liability at `$4.20`, while `TERMS.md:15` pushes content responsibility to user. Decision: organizational risk transfer, not a runtime accountability mechanism.

## Notable Patterns

- **HITL as first-class pause/resume:** `RequiresApprovalToolRule` → `ApprovalRequestMessage` (persisted with `tool_calls`) → `LettaStopReason(stop_reason="requires_approval")` → client sends `ApprovalCreate` → loop resumes with filtered `tool_calls`/`denials`/`tool_returns`. Pattern recurs in both v2 (`letta/agents/letta_agent_v2.py:1138-1153`) and v3 (`letta/agents/letta_agent_v3.py:1682-1709`) agents.
- **Serialization hardening for audit data:** `ApprovalsColumn` (`letta/orm/custom_columns.py:116`) delegates to `serialize_approvals`/`deserialize_approvals` which sanitize null bytes and tolerate mixed `ApprovalReturn|ToolReturn|dict` payloads (`letta/helpers/converters.py:293-353`) — defensive pattern for durable human-decision logs.
- **Context-vars log attribution:** `letta/log_context.py:4` `_log_context: ContextVar[dict]` + `update_log_context(agent_id, actor_id)` (`letta/log.py:153`) injects accountability metadata into every log line/otel span, bridging from DB attribution to observability.
- **Role-based message polymorphism:** Single `messages` table with `role ∈ {system, assistant, user, tool, approval, summary}` (`letta/schemas/message.py:322`) and `to_letta_messages()` dispatcher (`letta/schemas/message.py:480-573`) converts internal rows to typed public messages — elegant but blurs audit boundaries.

## Tradeoffs

- **Durability vs queryability:** Storing approvals as JSON `ApprovalsColumn` on `messages` preserves fidelity without extra tables, but makes answering "who approved what and when?" require JSON deserialization and cross-join with `Message.created_at`/`sender_id`; a dedicated `approvals` table with foreign keys would enable indexed accountability queries. Current choice favors write simplicity (`letta/helpers/converters.py:293`) over read observability.
- **Forensic attribution vs user-facing attribution:** Rich internal provenance (`Step.model`/`provider_name`/`trace_id`, `Message.sender_id`, `Tool.created_by_id`) exists but is not surfaced in `LettaMessage` API contracts. Tradeoff: avoids leaking internal details, at cost of accountability transparency for end users who see unattributed tool outputs.
- **Coarse organizational attribution vs fine-grained human attribution:** `actor` is always a `User` (`letta/schemas/user.py:16`) but approvals do not record which human identity performed the approval beyond `sender_id` on the approval response message (`letta/schemas/message.py:297`); there is no `approved_by_user_id` explicit field, relying on implicit `Message.sender_id`/`sender_id` which may be null (`tests/test_agent_files/test_agent.af:206` shows `"sender_id": null` common).
- **Legal terms vs technical enforcement:** Responsibility is documented in `TERMS.md:15` but not enforced/ queryable at runtime; the gap between paper policy and runtime mechanism creates a liability seam.

## Failure Modes / Edge Cases

- **Missing model attribution on approval-path messages:** If `approval_request.step_id` is null (legacy rows), v3 generates a new `step_id` (`letta/agents/letta_agent_v3.py:1020-1022`) — loses causal linkage between model decision and human decision, weakening accountability chain.
- **Null sender_id erases human identity:** `sender_id` is nullable everywhere (`letta/orm/message.py:59`, tests show null in `test_agent.af:206`); approval responses with null `sender_id` cannot answer "which human approved?" — silent accountability gap.
- **Malformed approval payload handled as stop, not audit event:** Empty `approvals` with no `tool_calls`/`denials`/`tool_returns` triggers `LettaStopReason(stop_reason="invalid_tool_call")` and `should_continue=False` (`letta/agents/letta_agent_v3.py:1009-1017`) but only logs error; no structured audit event records the malformed human input.
- **Concurrent tool execution attribution ambiguity:** `Tool.enable_parallel_execution` (`letta/schemas/tool.py:62`) allows concurrent tool calls; `ToolReturnMessage` aggregates via `tool_returns` list (`letta/schemas/letta_message.py:303`) without per-tool-call actor binding — if one tool is human-approved and another is auto-executed, the aggregate message obscures per-tool responsibility.
- **Bare UUID fallback masks actor identity:** `LettaBase.allow_bare_uuids` (`letta/schemas/letta_base.py:79-89`) accepts bare UUIDs with only debug log, potentially obscuring provenance during migrations.
- **No disclaimer means downstream misattribution:** User-facing clients receiving `AssistantMessage.content` (`letta/schemas/letta_message.py:364`) have no in-band signal that content is model-generated; downstream systems that log or republish agent output may incorrectly attribute it to a human.
- **Webhook/step failure loses accountability trail:** `Step.status = FAILED` with `error_data` (`letta/services/step_manager.py:402-404`) is recorded, but tool attribution for failed steps (which tool caused failure) requires correlating `Step.messages` via `step_id` — not indexed for efficient forensic query (`letta/orm/message.py:36` has `idx_messages_step_id` but no accountability view).

## Future Considerations

- Introduce explicit `Responsibility` value object (`actor_id`, `model`, `provider`, `tool_id`+`tool_author_id`, `human_approver_id`, `organization_id`, `timestamp`, `step_id`) persisted per action and exposed via `GET /v1/agents/{id}/actions/{step_id}/responsibility` — would make "who is responsible?" a first-class query rather than a manual join.
- Add user-facing model attribution headers/fields: surface `model`+`provider` on `AssistantMessage`/`ToolCallMessage` and optional `X-Generated-By: letta/model` or in-content disclaimer toggle for deployers needing watermarking.
- Separate `approvals` table with FKs (`approval_request_message_id`, `tool_call_id`, `approved_by_user_id`, `reason`, `created_at`) to enable indexed accountability queries and enforce non-null approver identity; migrate `ApprovalsColumn` JSON to normalized rows.
- Enforce non-null `sender_id` on `ApprovalResponseMessage` creation and validate `approved_by` against org membership at `MessageCreate` validation (`letta/schemas/message.py:178`).
- Provide accountability documentation: `docs/accountability.md` mapping each action type (LLM generation, tool execution, memory write, human approval) to responsible principal and linking to observable fields (`Step.model`, `Tool.created_by_id`, `Message.sender_id`, `Run.organization_id`).
- Emit structured audit events (e.g., `letta.audit.tool_executed`, `letta.audit.approval_granted`) via OTEL (`letta/otel/tracing.py`) with stable `responsibility` attributes, enabling SIEM integration.
- Add integration tests asserting attribution invariants: every `ToolCallMessage` has resolvable `step.model`; every `ApprovalResponseMessage` has non-null `sender_id`; unattributed tool outputs fail CI.

## Questions / Gaps

- No evidence found for model-output disclaimers or watermarking — searched `disclaimer`, `AI-generated`, `generated by` across `letta/` and found zero runtime implementations; `AI_POLICY.md:57` only discusses human disclosure for contributions, not agent outputs.
- No evidence of organization/operator responsibility dashboard — searched `responsib*`, `accountab*`, `audit*` in docs and code; only legal `TERMS.md:14` and low-level `created_by_id` fields exist.
- Unclear whether `sender_id` on `ToolCallMessage`/`ApprovalRequestMessage` represents model vs agent identity vs human — docs describe it as "identity id or agent id" (`letta/schemas/message.py:297`) without responsibility semantics.
- No evidence of tool-author runtime attribution surfacing — `Tool.created_by_id` exists but no API returns "this tool was authored by user X and approved by admin Y" alongside execution.
- Retention/immutability of accountability records unproven — `Message`, `Step`, `Run` support `AccessType.USER` predicates (`letta/orm/sqlalchemy_base.py:985`) but no append-only or tamper-evident guarantee for audit logs.

---

Generated by `23.03-responsibility-and-accountability-model` against `letta`.
