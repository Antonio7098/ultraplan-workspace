# Source Analysis: langgraph

## 23.03 Responsibility and Accountability Model

### Source Info

| Field | Value |
|-------|-------|
| Name | langgraph |
| Path | `studies/agent-harness-study/sources/langgraph` |
| Language / Stack | Python (monorepo: langgraph, checkpoint, prebuilt, sdk-py, checkpoint-sqlite/postgres) |
| Analyzed | 2026-08-29 |

## Summary

LangGraph is a graph-orchestration library, not an agent with a built-in accountability policy. Responsibility is intentionally delegated outward: the framework provides durable provenance primitives (checkpoints, `CheckpointMetadata.source`, task-level provenance, interrupt/Command(resume) persistence, and debug/lifecycle streams) but declares itself not responsible for LLM output content, tool argument safety, or checkpoint access control. The `.github/THREAT_MODEL.md:405-437` trust-boundary table is the only explicit accountability artifact, mapping each validation point to `User`, `Shared`, or `Project` responsibility. There are no output disclaimers, no model-output attribution, and no identity-bound audit log; human approvals are recorded as anonymous `RESUME` pending writes plus a timestamped checkpoint, not as a signed decision by a principal.

## Rating

**5/10 — Present but inconsistent, weakly documented, or fragile**

Rationale: Durable execution provenance exists (`BaseCheckpointSaver` + `get_state`/`get_state_history` + `CheckpointMetadata` + `PregelLoop` writes + `GraphInterruptEvent`/`GraphResumeEvent` callbacks) and is tested at scale, which answers *what* happened and *when*. However the system cannot answer *who* is responsible in a policy sense: no model-output attribution or disclaimer, tool decisions are attributed only to a `tool_call_id`/`task_id` pair without a principal, human resumptions lack identity, and formal accountability documentation is limited to the experimental auto-generated threat model rather than a governed RACI or operational runbook.

## Evidence Collected

Every entry MUST include a file path with line numbers. Format: `path/to/file.ts:NN`.

| Area | Evidence | File:Line |
|------|----------|-----------|
| Policy attribution | Trust-boundary table explicitly assigns Responsibility column: `User`, `Shared (project owns serializer defaults; user owns DB access controls)` etc., plus note "Node implementation safety is user's responsibility" | `libs/langgraph/.github/THREAT_MODEL.md:405` |
| Policy attribution | Out-of-scope section states prompt injection, state poisoning, tool argument safety are *user* responsibility; framework's responsibility ends at tool-name allowlist routing (`ToolNode._validate_tool_call`) and injection merge order | `libs/langgraph/.github/THREAT_MODEL.md:417-435` |
| Policy attribution | `THREAT_MODEL.md` disclaimer: "automatically generated ... experimental, subject to change, and not an authoritative security reference" | `libs/langgraph/.github/THREAT_MODEL.md:5` |
| Checkpoint provenance | `CheckpointMetadata.source: Literal["input","loop","update","fork"]` with `step`, `parents`, `run_id`, `counters_since_delta_snapshot` — classifies origin of every checkpoint | `libs/checkpoint/langgraph/checkpoint/base/__init__.py:38-41` |
| Checkpoint provenance | `get_checkpoint_metadata()` merges `config["metadata"]` and `config["configurable"]` into checkpoint metadata, filtering `EXCLUDED_METADATA_KEYS` (`thread_id`, `checkpoint_id`, `langgraph_node` etc.) — preserves `run_id` and user-supplied metadata | `libs/checkpoint/langgraph/checkpoint/base/__init__.py:757-776` |
| Checkpoint provenance | `PregelLoop._put_checkpoint({"source":"input"})` and `{"source":"loop"}` and `{"source":"fork"}` — each superstep writes a checkpoint with source-labeled metadata | `libs/langgraph/langgraph/pregel/_loop.py:1016,706,959` |
| Checkpoint provenance | `EXCLUDED_METADATA_KEYS` set lists keys stripped before persistence, showing what is *not* attributable via metadata | `libs/checkpoint/langgraph/checkpoint/base/__init__.py:797-808` |
| Execution traceability | `ExecutionInfo` dataclass exposes `checkpoint_id`, `checkpoint_ns`, `task_id`, `thread_id`, `run_id`, `node_attempt` injected via `Runtime` — read-only per-node execution identity | `libs/langgraph/langgraph/runtime.py:27-58` |
| Execution traceability | `PregelLoop` fields `checkpoint_ns: tuple[str,...]`, `checkpoint_id_saved`, `checkpoint_config`, `checkpoint_metadata`, `checkpoint_pending_writes: list[PendingWrite]` — durable per-thread audit trail | `libs/langgraph/langgraph/pregel/_loop.py:241-245` |
| State history / audit log | `Pregel.get_state()`, `aget_state()`, `get_state_history()`, `bulk_update_state()`, `update_state()` — full history read path via `checkpointer.list()` and `prepare_state_snapshot()` yielding `StateSnapshot` with `values`, `next`, `tasks`, `interrupts`, `metadata`, `parent_config`, `created_at` | `libs/langgraph/langgraph/pregel/main.py:1391,1479,1589,2530` |
| State history / audit log | `_prepare_state_snapshot()` builds `StateSnapshot(interrupts=tuple([i for task in tasks_with_writes for i in task.interrupts]))` — preserves interrupt chain per step | `libs/langgraph/langgraph/pregel/main.py:1256-1264` |
| Interrupt / human decision | `Interrupt` dataclass with `value: Any` and `id: str` (derived `xxh3_128_hexdigest(ns)`), `interrupt(value)` raises `GraphInterrupt((Interrupt,...))` | `libs/langgraph/langgraph/types.py:534-557,811-934` |
| Interrupt / human decision | `Command` dataclass `graph`, `update`, `resume: dict[str,Any]|Any|None`, `goto: Send|Sequence[Send|N]|N` — generic human-approval primitive; `Command(resume=...)` rehydrates via `_first()` mapping to `RESUME` writes | `libs/langgraph/langgraph/types.py:759-808` |
| Interrupt / human decision | `PregelLoop._first()` distinguishes `input_is_command`, `is_resuming`, `is_time_traveling`; on `Command(resume=map|value)` sets `CONFIG_KEY_RESUME_MAP` or validates single-interrupt constraint | `libs/langgraph/langgraph/pregel/_loop.py:836-918` |
| Interrupt / human decision | `_pending_interrupts()` scans `checkpoint_pending_writes` for `INTERRUPT` vs `RESUME` to detect unresolved human gates | `libs/langgraph/langgraph/pregel/_loop.py:806-834` |
| Interrupt / human decision | `GraphInterruptEvent` and `GraphResumeEvent` frozen dataclasses with `run_id`, `status: GraphLifecycleStatus`, `checkpoint_id`, `checkpoint_ns`, `interrupts: tuple[Interrupt,...]` plus `GraphCallbackHandler.on_interrupt/on_resume` | `libs/langgraph/langgraph/callbacks.py:43-112` |
| Interrupt / human decision | `PregelLoop._push_graph_lifecycle_event(kind="resume"|"interrupt")` appends `GraphInterruptEvent`/`GraphResumeEvent` to `deque[GraphLifecycleEvent]` for external callback dispatch | `libs/langgraph/langgraph/pregel/_loop.py:369-402` |
| Interrupt / human decision | Config-driven `interrupt_before: All|Sequence[str]` and `interrupt_after: All|Sequence[str]` declared on `StateGraph.compile()` and enforced in `PregelLoop.tick()`/`after_tick()` via `should_interrupt()` raising `GraphInterrupt` | `libs/langgraph/langgraph/graph/state.py:1170-1251` and `libs/langgraph/langgraph/pregel/_loop.py:660-712` |
| Tool attribution | `ToolNode.tools_by_name: dict[str,BaseTool]` + `_validate_tool_call()` returning error `ToolMessage(name=requested_tool, tool_call_id, status="error")` on unknown tool name | `libs/prebuilt/langgraph/prebuilt/tool_node.py:789-791,1268-1279` |
| Tool attribution | `_inject_tool_args()` merge `{**llm_args, **injected_args}` where `ToolRuntime(state, tool_call_id, config, context, store, stream_writer, execution_info)` supplies system values; system always wins on collision | `libs/prebuilt/langgraph/prebuilt/tool_node.py:1315-1385`, `libs/prebuilt/langgraph/prebuilt/tool_node.py:803-817` |
| Tool attribution | `ToolCallRequest(tool_call, tool, state, runtime)` and `ToolCallWithContext(tool_call, __type, state)` used for Send-based parallel dispatch and `awrap_tool_call/wrap_tool_call` interception | `libs/prebuilt/langgraph/prebuilt/tool_node.py:133-149,286-307` |
| Tool attribution | `ToolMessage(content, name, tool_call_id, status)` produced via `_execute_tool_sync/_execute_tool_async` → `_normalize_tool_response()`; error branch creates `ToolMessage(status="error")` with filtered validation errors excluding injected args | `libs/prebuilt/langgraph/prebuilt/tool_node.py:922-1012,1069-1159` |
| Tool attribution | `StreamToolCallHandler._tool_call_writer: ContextVar[ToolCallWriter|None]` emits `tool-finished`/`tool-error` payloads keyed by `tool_call_id` (run_id→(ns, tool_call_id, token) map) | `libs/langgraph/langgraph/pregel/_tools.py:25-42,80-220` |
| Tool attribution | `debug.map_debug_tasks()` → `TaskPayload(id,name,input,triggers,metadata)` and `map_debug_task_results()` → `TaskResultPayload(id,name,error,result,interrupts)` preserve per-task tool origin | `libs/langgraph/langgraph/pregel/debug.py:41-128` |
| Tool attribution | `_TasksLifecycleBase._record_pending_tool_calls()` harvests `task_id→tool_call_id` from `tool_call_with_context` inputs; `LifecycleTransformer` surfaces `cause={"type":"toolCall","tool_call_id":...}` on lifecycle events | `libs/langgraph/langgraph/stream/transformers.py:488-516,536-540,648-656` |
| Tool attribution | `PregelTask(id,name,path,error,interrupts,state,result)` + `StateSnapshot.tasks: tuple[PregelTask,...]` + `tasks_w_writes()` applying pending_writes per task | `libs/langgraph/langgraph/types.py:597-606,643-661` and `libs/langgraph/langgraph/pregel/debug.py:209-279` |
| Output attribution/disclaimer | No evidence found of disclaimer injection or model-origin watermarking in `libs/langgraph/langgraph/types.py:120-134` StreamMode docs, `libs/langgraph/langgraph/pregel/_messages.py:308-335` messages stream, or `libs/langgraph/README.md` | (search: `disclaimer`, `attribution` across `libs/langgraph` — only 4 hits unrelated: `libs/langgraph/langgraph/stream/_mux.py:248` `attributions` for debug keys) |
| Human identity | `Runtime.ServerInfo(assistant_id, graph_id, user: BaseUser|None)` exists but is *injected* only when running on LangGraph Server; never persisted to `CheckpointMetadata` (which only stores `run_id`, `step`, `parents`, `source`) | `libs/langgraph/langgraph/runtime.py:61-77` and `libs/checkpoint/langgraph/checkpoint/base/__init__.py:38-63` |
| Observability | `StreamMessagesHandler` and `MessagesTransformer.process()` enrich chat streams with `metadata: {langgraph_node, langgraph_step, langgraph_triggers, run_id}` — identifies *which node* emitted a message but not *which principal* approved it | `libs/langgraph/langgraph/stream/transformers.py:272-286` |

## Answers to Dimension Questions

### 1. Who is responsible for agent actions?

**Developer/operator/organization — not the model, runtime author, or framework.** LangGraph explicitly positions itself as a library. `libs/langgraph/.github/THREAT_MODEL.md:405` tags every trust-boundary row with Responsibility (`User` for graph structure, node implementations, tool behavior; `Shared` for checkpoint storage and LLM output tool-name routing). Section 417-437 enumerates out-of-scope threats ("Prompt injection leading to arbitrary tool execution ... Project does not control LLM model behavior", "State poisoning via malicious node output ... user controls their own code") and states "The framework's responsibility is to not execute unregistered tools and to correctly route registered ones." There is no built-in policy file, config key, or code path that assigns blame to the model or to LangGraph Inc.; accountability is delegated to the application deployer who chooses the checkpointer, implements `Auth` handlers (`libs/sdk-py/langgraph_sdk/auth/__init__.py:Auth` referenced at `THREAT_MODEL.md:437`), and operates the database. This can answer *what* executed (node name + task_id in `PregelTask:597` and `StateSnapshot:643`) but not *who* is answerable without external identity integration.

### 2. Is model output attributed?

**No.** No source file injects a disclaimer, watermark, or attribution label on `AIMessage`/chunks. Search for `disclaimer`/`attribution` across `libs/langgraph` yields only `StreamMux` debug attributions at `libs/langgraph/langgraph/stream/_mux.py:248`. The messages stream path (`libs/langgraph/langgraph/pregel/_messages.py:308`, `libs/langgraph/langgraph/stream/transformers.py:272`) does attach `metadata={langgraph_node, langgraph_step, langgraph_triggers, run_id}` identifying the *node* and *step* that produced the message, and `ExecutionInfo:27` carries `task_id`/`run_id`/`checkpoint_id`. However this is execution plumbing, not a model-output disclaimer such as "generated by <model>" with timestamp/principal. The README and docs make no claim to output attribution, consistent with the threat model's assumption that "LLM provider behavior — Model output content and safety" is out of scope (`THREAT_MODEL.md:28`).

### 3. Are tool decisions attributed?

**Partially — attribution exists at the execution level, not at the principal/policy level.** Every tool invocation is durably traceable to a typed triple:
- LLM `AIMessage.tool_calls[]` with `id` → `ToolNode._parse_input() → _run_one() → _inject_tool_args()` at `libs/prebuilt/langgraph/prebuilt/tool_node.py:1224-1315`
- `ToolMessage(name, tool_call_id, status)` returned via `_execute_tool_sync/_execute_tool_async:922-1159`
- `PendingWrite(task_id, channel, value)` and `TaskPayload`/`TaskResultPayload` persisted via `BaseCheckpointSaver.put_writes()` at `libs/checkpoint/langgraph/checkpoint/base/__init__.py:300-318` and surfaced via `map_debug_tasks:41` / `map_debug_task_results:106` and `StateSnapshot.tasks:643`
- Streaming provenance via `_TasksLifecycleBase._record_pending_tool_calls:488` yielding `cause={"type":"toolCall","tool_call_id":...}:538` and `StreamToolCallHandler` keying events by `tool_call_id:149`

What is missing is *who authorized* the tool run: there is no principal, role, or policy id attached. `ToolRuntime:804` carries `tool_call_id`, `execution_info`, and `config` but not a user identity; checkpoint metadata excludes it via `EXCLUDED_METADATA_KEYS:797`. So the system can answer "tool X was called with args Y by task Z at checkpoint C producing ToolMessage M" — answerable via `get_state_history()` — but not "principal P delegated this tool decision."

### 4. Are human approvals recorded?

**Yes — as anonymous gated writes, without identity.** Human-in-the-loop is a first-class primitive: `interrupt(value) → GraphInterrupt` (`libs/langgraph/langgraph/types.py:811`), `Command(resume=...)` (`libs/langgraph/langgraph/types.py:759`), static gates `interrupt_before`/`interrupt_after` (`libs/langgraph/langgraph/graph/state.py:1170`), and `bulk_update_state`/`update_state` human overrides (`libs/langgraph/langgraph/pregel/main.py:1589`). Every pause and resume is checkpointed:
- `INTERRUPT` pending write stores `Interrupt(value, id=xxh3_128_hexdigest(ns)):929` in `checkpoint_pending_writes`
- `RESUME` pending write stores the human-supplied value via `CONFIG_KEY_RESUME_MAP`/`put_writes(NULL_TASK_ID/RESUME)` in `PregelLoop._first:891-919`
- Timestamp `checkpoint["ts"]` plus `CheckpointMetadata.step/parents/run_id` provides *when* and *in which thread*
- Lifecycle callbacks `GraphInterruptEvent(checkpoint_id, checkpoint_ns, interrupts, status)` and `GraphResumeEvent` at `libs/langgraph/langgraph/callbacks.py:43-77` plus `_push_graph_lifecycle_event:369` create observable callbacks for HITL transitions.

Gaps: `BaseUser`/`ServerInfo.user:61` is only available on the server deployment and is not persisted to `CheckpointMetadata` (`libs/checkpoint/langgraph/checkpoint/base/__init__.py:38` has no user field). So the log records *that* approval happened and *what* value resumed, but not *who* approved. Multiple-pending-interrupt resumption requires `resume: dict[interrupt_id, value]:899` but still carries no signature.

### 5. Is accountability documented?

**Weakly — an auto-generated threat model, not governed accountability docs.** The sole documented accountability artifact is `libs/langgraph/.github/THREAT_MODEL.md`, flagged at line 5 as "automatically generated ... experimental, subject to change, and not an authoritative security reference." It does contain a `Responsibility` column per input source (`User`, `Shared`, `Project`) and an explicit "Project Responsibility Ends At" table at `THREAT_MODEL.md:417-438`, plus assumptions (1-7) that the user controls app code, DB ACLs, and model selection. There is no `ACCOUNTABILITY.md`, no RACI, no runbook describing who must sign off on tool execution or how human approvals are audited in production. Checkpoint storage is documented as having "unbounded retention" with no governance (`BaseCheckpointSaver` at `libs/checkpoint/langgraph/checkpoint/base/__init__.py:176` has no TTL; `THREAT_MODEL.md:DC3` notes "Unbounded (no default TTL)"). For a compliance reviewer, the framework provides mechanics to *build* an accountability log (checkpointer + callbacks + `get_state_history`) but does not itself define or enforce one.

## Architectural Decisions

| Decision | Evidence | Effect on Accountability |
|----------|----------|--------------------------|
| **Library, not server** — auth, TLS, DB ACLs delegated to deployer; `Auth` handler is opt-in server-side code | `libs/langgraph/.github/THREAT_MODEL.md:37,437` `BaseCheckpointSaver` is abstract with `put/get_tuple/list` at `libs/checkpoint/langgraph/checkpoint/base/__init__.py:176-277` | Makes LangGraph agnostic to tenant identity; provenance exists but principal binding must be added by operator. |
| **Checkpoint as audit log** — every superstep writes a checkpoint with `source/step/parents` and `pending_writes` durably stored via `BaseCheckpointSaver.put_writes` | `libs/langgraph/langgraph/pregel/_loop.py:1064-1199`, `libs/checkpoint/langgraph/checkpoint/base/__init__.py:281-318` | Enables time-travel debugging (`get_state_history:1479`) and deterministic replay, but audit trail is technical (channel values + task results) not compliance-oriented. |
| **Interrupt/Command as HITL primitive** — `interrupt()` suspends on checkpointer-backed `GraphInterrupt`; `Command(resume=map)` resumes with per-interrupt mapping | `libs/langgraph/langgraph/types.py:811-934`, `libs/langgraph/langgraph/types.py:759-808`, `libs/langgraph/langgraph/pregel/_loop.py:806-918` | Gives explicit human-gate semantics with correlation ids (`Interrupt.id`); however ids are namespace-hashed, not user-signed. |
| **Tool injection merge order** — `{**llm_args, **injected_args}` in `ToolNode._inject_tool_args:1380` plus `tool_call_schema` hiding injected param names | `libs/prebuilt/langgraph/prebuilt/tool_node.py:1315+`, `libs/langgraph/.github/THREAT_MODEL.md:DF14` | Guarantees LLM cannot spoof system-supplied `state`/`store`/`runtime`, a trust-boundary defense, but does not label *who* approved the tool. |
| **Graph lifecycle callbacks** — `GraphInterruptEvent/GraphResumeEvent` dispatched via `GraphCallbackHandler` and `get_sync_graph_callback_manager_for_config` | `libs/langgraph/langgraph/callbacks.py:43-394`, `libs/langgraph/langgraph/pregel/_loop.py:369-402` | Provides observable hooks for external audit sinks (e.g., LangSmith), but callbacks are opt-in and receive no user principal. |
| **Msgpack allowlist gating** — `SAFE_MSGPACK_TYPES` (47) / `SAFE_MSGPACK_METHODS` with strict `LANGGRAPH_STRICT_MSGPACK` toggle | `libs/langgraph/.github/THREAT_MODEL.md:TB2` | Shows accountability for *deserialization integrity* (what types may be restored), separate from *decision* accountability. |

## Notable Patterns

- **Pending-writes + checkpoint fork model for HITL auditability** (`libs/langgraph/langgraph/pregel/_loop.py:836-962`, `libs/langgraph/langgraph/pregel/main.py:1589-2053`): every human decision creates a `source="fork"` or `source="update"` checkpoint, preserving the parent chain for forensic walk via `get_delta_channel_history:582`. Pattern is well-tested (`test_time_travel*`, `test_subgraph_persistence*`).
- **Dual-channel message attribution** (`libs/langgraph/langgraph/pregel/_messages.py:308` + `libs/langgraph/langgraph/stream/transformers.py:272`): node-origin metadata (`langgraph_node`, `langgraph_step`, `langgraph_triggers`, `run_id`) flows alongside `AIMessageChunk` so streams can be correlated to `TaskPayload.metadata:142-162`.
- **Three-phase PregelLoop tick** (`prepare_next_tasks → should_interrupt(before) → execute → apply_writes → should_interrupt(after)` at `libs/langgraph/langgraph/pregel/_loop.py:592-714`) makes interrupt points deterministic and therefore auditable — every before/after gate either executes or raises `GraphInterrupt` in a well-defined status (`GraphLifecycleStatus:31-38`).
- **Debug-mode provenance folding** (`map_debug_checkpoint:144`, `map_debug_tasks:41`, `map_debug_task_results:106`, `tasks_w_writes:209`): aggregates per-task errors/interrupts/results from raw `pending_writes` into `StateSnapshot.tasks` for human-readable audit via `stream_mode="debug"`.
- **Transparent human-override API** (`bulk_update_state:1589` looping `perform_superstep` with `StateUpdate(values, as_node, task_id)`): allows an operator to impersonate any node deterministically — powerful for accountability repair (e.g., correcting poisoned state) but also a principal-spoofing vector without access logging.

## Tradeoffs

| Tradeoff | Pro | Con |
|----------|-----|-----|
| **Delegated identity vs built-in identity** — Checkpoint stores `thread_id/run_id/step` but not `userId`/`role` | Keeps framework deployable anywhere (SQLite in-process, Postgres, serverless); no auth tax for single-user notebooks | Cannot answer "who approved this resume" without operator-added metadata propagation via `config["metadata"]` → `get_checkpoint_metadata:757` |
| **Full-history retention by default** (`list()` eagerly consumes `checkpointer.list()` at `libs/langgraph/langgraph/pregel/main.py:1525`) | Complete audit trail for time-travel debugging and regulatory replay | Unbounded PII retention (`THREAT_MODEL.md:DC3/DC6` — "Unbounded by default ... No built-in TTL"), cost and GDPR tension |
| **Allowlist-secured serde vs permissive default** — `LANGGRAPH_STRICT_MSGPACK` off by default, 47 safe types otherwise | Permissive default maximizes compatibility for users storing arbitrary Python objects | Threat T1/T2: unsafe deserialization if attacker gains DB write access; accountability of *what* was restored is conditional on deployer hardening |
| **Callback-based observability vs checkpoint durability** — HITL events emit both to checkpoint writes and to `GraphCallbackHandler` | Supports live audit sinks without polling `get_state_history` | Events only fire if caller registered a `GraphCallbackHandler` via `config["callbacks"]:363`; silent loss if omitted |
| **Command as universal mutator** — same primitive resumes interrupts, injects state updates, and navigates (`goto: Send`) | Single audit concept simplifies logging; `Command.PARENT` fork attribution is explicit at `libs/langgraph/langgraph/types.py:808` | Overloading makes it hard to distinguish "human approval" from "programmatic goto" in audit log without inspecting `Command` fields |

## Failure Modes / Edge Cases

- **Anonymous human resumption**: `StateSnapshot.metadata` (`libs/checkpoint/langgraph/checkpoint/base/__init__.py:38`) lacks a principal field; if two operators share a `thread_id`, their `Command(resume=...)` values are indistinguishable in the audit log. `ServerInfo.user` (`libs/langgraph/langgraph/runtime.py:61`) is not persisted, so server deployments lose the actor on checkpoint read.
- **Multiple pending interrupts with map resume**: `libs/langgraph/langgraph/pregel/_loop.py:904` raises if `len(_pending_interrupts())>1` and resume is not a map — correct fail-safe, but error is raised late (at `stream()` entry), not at compile time, so accountability tooling must handle the exception path.
- **Unbounded checkpoint chain breaks attribution**: `list()` / `get_delta_channel_history:582` walks parent chain; if an operator prunes intermediate checkpoints without `DeltaChannel`-awareness, history silently reconstructs as empty (see `BaseCheckpointSaver.prune` warnings at `libs/checkpoint/langgraph/checkpoint/base/__init__.py:387-414`). The audit trail is then incomplete.
- **Fork without identity**: `bulk_update_state` with `as_node=="__copy__":1799` and `StateUpdate(task_id=None)` auto-generates task ids via `uuid5(UUID(checkpoint["id"]), INTERRUPT):1966`. The act of forking is recorded as `source="fork":1819` but carries no "who forked" metadata unless caller propagated it via `config["metadata"]`.
- **Tool-call-id spoofing within a thread**: `ToolMessage.tool_call_id` is taken from `AIMessage.tool_calls[i]["id"]` (`libs/prebuilt/langgraph/prebuilt/tool_node.py:823`); a malicious prompt-injected LLM output could reuse a prior `tool_call_id`. Checkpoint provenance records the reuse but does not detect it as anomalous — detection is user-tool responsibility.
- **Callback loss on failure**: `GraphInterruptEvent` is enqueued in `PregelLoop._graph_lifecycle_events:262` and dispatched only on `status != "draining":377`; if the loop drains (`RunControl.request_drain:95`) or hits `out_of_steps:600`, the interrupt is swallowed without a lifecycle event, so an external audit sink may miss the gate.
- **Cross-namespace interrupt liability**: `libs/langgraph/langgraph/pregel/_algo.py:845` comment notes "responsibility lies with the parent" for subgraph interrupts — parent vs child attribution can be ambiguous when auditing nested `StateSnapshot` with `subgraphs=True` (`libs/langgraph/langgraph/pregel/main.py:1391`).

## Future Considerations

- **Identity-bound checkpoints**: Extend `CheckpointMetadata` (`libs/checkpoint/langgraph/checkpoint/base/__init__.py:38`) with optional `principal_id`, `principal_type`, `auth_method` populated from `ServerInfo.user` or `config["metadata"]["user_id"]`, and enforce via `get_checkpoint_metadata` without placing identity in `EXCLUDED_METADATA_KEYS`. Enables "who approved this Command" queries on `get_state_history`.
- **Output disclaimer transformer**: Add a `StreamTransformer` or node middleware that appends a structured disclaimer block (`generated_by`, `model`, `timestamp`, `disclaimer`) to `AIMessage` content before channel write, with opt-out flag at `Pregel` construction. Closes Q2 gap without forcing it on non-LLM graphs.
- **Signed human decisions**: Hash `Interrupt.id` + `resume` value + `principal_id` + `checkpoint_id` and expose in `GraphInterruptEvent`/`GraphResumeEvent`; store signature in `CheckpointMetadata.writes` or a new metadata field so resumptions are non-repudiable.
- **Governed retention + audit export**: Provide a bounded-retention `BaseCheckpointSaver` wrapper and a first-class `export_audit_log(thread_id)→NDJSON` that emits `checkpoint_id`, `source`, `run_id`, `principal_id`, `node`, `tool_call_id`, `interrupts`, `created_at` — replaces raw `list()` consumption for compliance pipelines.
- **Promote threat model to accountability spec**: Graduate `THREAT_MODEL.md` from auto-generated disclaimer to a versioned `ACCOUNTABILITY.md` with explicit RACI per action class (graph definition, state mutation, tool execution, interrupt approval) and link each row to enforcement point file:line numbers.
- **Differentiate `Command` subtypes in history**: Tag `CheckpointMetadata.source` or a new `command_kind: "resume"|"goto"|"update"` so audit queries can filter human approvals vs programmatic navigation without parsing `Command` fields.

## Questions / Gaps

- No output disclaimer or model attribution is implemented — search for `disclaimer` yields only `THREAT_MODEL.md:5` disclaimer about the document itself. Cannot verify model-output authenticity from checkpoint alone.
- No formal accountability documentation beyond the experimental threat model — no `RACI.md`, `AUDIT.md`, or policy manifest ties actions to principals. The threat model's own header disavows authoritativeness.
- Human decision logs lack principal identity — `CheckpointMetadata:38` has no `user` field; `ServerInfo.user:61` is transient; checkpoint queries cannot attribute a resume to a person.
- Tool attribution is execution-scoped (`task_id`/`tool_call_id`) not policy-scoped — there is no "tool approval policy" artifact (allowlist of which tools which roles may approve) in the codebase.
- Tests for accountability invariants (e.g., "every `Command(resume)` produces a `source=update` checkpoint with `run_id` preserved") are not named as such — HITL tests exist (`test_interruption*`, `test_time_travel*`, `test_state*`) but do not assert accountability properties explicitly.

---

Generated by `Dimension 23.03: Responsibility and Accountability Model` against `langgraph`.
