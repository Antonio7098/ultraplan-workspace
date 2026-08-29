# Source Analysis: agent-framework

## Artifact Lifecycle

### Source Info

| Field | Value |
|-------|-------|
| Name | agent-framework |
| Path | `studies/agent-harness-study/sources/agent-framework` |
| Language / Stack | Python (agent-framework-core, 30+ packages) + .NET (Microsoft.Agents.AI.*) |
| Analyzed | 2026-08-28 |

## Summary

Agent Framework is an agent-and-workflow orchestration SDK, not a run-artifact platform. Its durable artifact model centers on **WorkflowCheckpoints** (graph-execution snapshots), **AgentSessions/HistoryProviders** (conversation transcripts), and **DurableAgentState** entities (Durable Task–backed agent memory). There is no unified artifact registry; each backend owns its files in isolation. Creation is programmatic via `WorkflowBuilder(checkpoint_storage=...)` or `FileHistoryProvider(storage_path=...)`; storage is pluggable (`CheckpointStorage` protocol, `FileCheckpointStorage`, `InMemoryCheckpointStorage`, `FileHistoryProvider`, `AgentFileStore`); versioning is a shallow `version="1.0"` plus a SHA-256 graph-signature for compatibility; run linkage is via `workflow_name` + `checkpoint_id` chain (not per-run instance ID); and retirement is entirely manual `delete()` with no TTL/retention policy for workflow artifacts (TTL exists only in docs for durable-agent entities).

## Rating

**Score: 6 / 10 — Present but fragile**

Rationale: Checkpoint/history lifecycle is explicitly modeled with protocol interfaces, two storage backends, safety hardening (path-traversal validation, atomic writes, `RestrictedUnpickler`), and extensive tests. However artifacts are fragmented across isolated stores, versioning is minimal, run-to-artifact discoverability requires manual timestamp/iteration filtering by `workflow_name`, and no automated retention/cleanup exists for workflow checkpoints or history files. This matches rubric 4-6 (present but inconsistent/weakly documented/fragile) stretching toward 7, but lacks operational safeguards (retention, observability, GC) expected for 7-8.

## Evidence Collected

Every entry MUST include a file path with line numbers. Format: `path/to/file.ts:NN`.

| Area | Evidence | File:Line |
|------|----------|-----------|
| Artifact schema — checkpoint | `WorkflowCheckpoint` dataclass: `workflow_name`, `graph_signature_hash`, `checkpoint_id`, `previous_checkpoint_id`, `timestamp`, `messages`, `state`, `pending_request_info_events`, `iteration_count`, `metadata`, `version="1.0"` | `python/packages/core/agent_framework/_workflows/_checkpoint.py:30-88` |
| Artifact schema — version field | `version: str = "1.0"` — static format version, never auto-bumped | `python/packages/core/agent_framework/_workflows/_checkpoint.py:88` |
| Artifact schema — durable entity | JSON schema `DurableAgentState` with `schemaVersion` pattern `^\d+\.\d+\.\d+$`, `data.conversationHistory` of `agentRequest`/`agentResponse` entries, each with `correlationId`, `orchestrationId`, `createdAt`, typed `chatContentItem` union | `schemas/durable-agent-entity-state.json:1-217` |
| Artifact schema — durable state class | `DurableAgentState.from_dict()` / `to_dict()`, `DurableAgentStateRequest`, `DurableAgentStateResponse` | `python/packages/durabletask/agent_framework_durabletask/_durable_agent_state.py:1-60` (via `python/packages/durabletask/agent_framework_durabletask/_entities.py:59-82`) |
| Artifact schema — run result | `WorkflowRunResult` extends `list[WorkflowEvent]` with `get_outputs()`, `get_intermediate_outputs()`, `get_request_info_events()`, `get_final_state()`, `status_timeline()` | `python/packages/core/agent_framework/_workflows/_workflow.py:101-165` |
| Artifact schema — event types | `WorkflowEventType` literal: `output`, `intermediate`, `request_info`, `superstep_started/completed`, `executor_*`, etc. | `python/packages/core/agent_framework/_workflows/_events.py:104-130` |
| Artifact schema — session | `AgentSession` with `session_id`, `service_session_id`, `state: dict`, `to_dict()`/`from_dict()` with `_serialize_state` handling `SerializationProtocol`/`BaseModel` | `python/packages/core/agent_framework/_sessions.py:746-811` |
| Artifact schema — history message | `Message`, `Content` unify 20+ content types (text, data, uri, function_call/result, shell, mcp, etc.) with `to_dict`/`from_dict` | `python/packages/core/agent_framework/_types.py:338-363`, `python/packages/core/agent_framework/_types.py:460-572` |
| Storage backend — protocol | `CheckpointStorage` Protocol: `save`, `load`, `list_checkpoints`, `delete`, `get_latest`, `list_checkpoint_ids` | `python/packages/core/agent_framework/_workflows/_checkpoint.py:119-189` |
| Storage backend — in-memory | `InMemoryCheckpointStorage` dict-backed with `copy.deepcopy` on save, max-timestamp `get_latest` | `python/packages/core/agent_framework/_workflows/_checkpoint.py:192-237` |
| Storage backend — file | `FileCheckpointStorage(storage_path, allowed_checkpoint_types)` mkdir, atomic `json.dump` via `.tmp` + `os.replace`, `json` + pickle+base64 hybrid | `python/packages/core/agent_framework/_workflows/_checkpoint.py:239-450` |
| Storage backend — encoding | `encode_checkpoint_value` / `decode_checkpoint_value` with `_PICKLE_MARKER="__pickled__"`, `RestrictedUnpickler` allowlist (`_BUILTIN_ALLOWED_TYPE_KEYS` + framework + openai types + caller extras), `_verify_type` | `python/packages/core/agent_framework/_workflows/_checkpoint_encoding.py:1-311` |
| Storage backend — history file | `FileHistoryProvider(storage_path)` JSONL per session, `_session_file_path` with resolve + `is_relative_to` check, `_session_file_stem` sanitization via `urlsafe_b64encode` or literal-safe check, striped thread+async locks | `python/packages/core/agent_framework/_sessions.py:893-1134` |
| Storage backend — history memory | `InMemoryHistoryProvider` stores `state["messages"]` list, `skip_excluded` filter | `python/packages/core/agent_framework/_sessions.py:814-891` |
| Storage backend — durable entity | `AgentEntityStateProviderMixin` cache + `_get_state_dict`/`_set_state_dict`, `DurableTaskEntityStateProvider(DurableEntity)` using `self.get_state(dict)` / `self.set_state` | `python/packages/durabletask/agent_framework_durabletask/_entities.py:35-353` |
| Storage backend — runner context | `InProcRunnerContext` holds `_messages`, `_event_queue`, `_pending_request_info_events`, `_checkpoint_storage` + `_runtime_checkpoint_storage`, `create_checkpoint` captures `messages+state+pending_events+iteration_count` | `python/packages/core/agent_framework/_workflows/_runner_context.py:278-502` |
| Naming — checkpoint file | `(self.storage_path / f"{checkpoint_id}.json").resolve()` + `_validate_file_path` checks `is_relative_to(self.storage_path.resolve())` | `python/packages/core/agent_framework/_workflows/_checkpoint.py:282-300` |
| Naming — checkpoint ID default | `checkpoint_id: CheckpointID = field(default_factory=lambda: str(uuid.uuid4()))` | `python/packages/core/agent_framework/_workflows/_checkpoint.py:74` |
| Naming — session file | `f"{self._session_file_stem(session_id)}{self.FILE_EXTENSION}"` (= `.jsonl`), `DEFAULT_SESSION_FILE_STEM="default"`, `_ENCODED_SESSION_PREFIX="~session-"`, Windows reserved stem blocklist | `python/packages/core/agent_framework/_sessions.py:913-1102` |
| Naming — workflow identity | `WorkflowBuilder(name or f"WorkflowBuilder-{uuid.uuid4()}")` → `Workflow.name`, `Workflow.id = str(uuid.uuid4())` (ephemeral instance ID), `graph_signature` + `graph_signature_hash = sha256(canonical JSON)` | `python/packages/core/agent_framework/_workflows/_workflow_builder.py:152`, `python/packages/core/agent_framework/_workflows/_workflow.py:329-333`, `python/packages/core/agent_framework/_workflows/_workflow.py:1015-1080` |
| Versioning — graph compatibility | `restore_from_checkpoint` validates `self._graph_signature_hash != checkpoint.graph_signature_hash` → `WorkflowCheckpointException("Workflow graph has changed")` | `python/packages/core/agent_framework/_workflows/_runner.py:275-279` |
| Versioning — state import semantics | `self._state.clear()` then `import_state(checkpoint.state)` merge on restore; `commit()` at superstep boundary | `python/packages/core/agent_framework/_workflows/_runner.py:286-287`, `python/packages/core/agent_framework/_workflows/_state.py:90-100` |
| Linking — checkpoint chain | `previous_checkpoint_id: CheckpointID \| None`, Runner tracks `previous_checkpoint_id` across supersteps, initial checkpoint is "superstep 0" after start executor, then per superstep | `python/packages/core/agent_framework/_workflows/_checkpoint.py:51-52`, `python/packages/core/agent_framework/_workflows/_runner.py:84-97`, `python/packages/core/agent_framework/_workflows/_runner.py:212-238` |
| Linking — builder vs runtime storage | `WorkflowBuilder(checkpoint_storage=storage)` build-time, overridden per-run via `workflow.run(..., checkpoint_storage=runtime_storage)` with runtime precedence pattern in tests | `python/packages/core/agent_framework/_workflows/_workflow_builder.py:96`, `python/packages/core/agent_framework/_workflows/_workflow.py:789-796`, `python/packages/core/tests/workflow/test_sequential.py:273-299` (in orchestrations) |
| Linking — response to request | `REQUEST_INFO` events carry `request_id`, `source_executor_id`, `request_type`, `response_type`; response message uses `MessageType.RESPONSE` with `original_request_info_event` and `INTERNAL_SOURCE_ID(source_executor_id)` | `python/packages/core/agent_framework/_workflows/_events.py:296-312`, `python/packages/core/agent_framework/_workflows/_runner_context.py:457-486` |
| Retirement — checkpoint delete | `FileCheckpointStorage.delete` via `file_path.unlink()` returns bool; `InMemoryCheckpointStorage.delete` removes from dict | `python/packages/core/agent_framework/_workflows/_checkpoint.py:392-410`, `python/packages/core/agent_framework/_workflows/_checkpoint.py:217-223` |
| Retirement — history durability gap | `FileHistoryProvider.save_messages` append-only `open("a")` with no delete/compact, `InMemoryHistoryProvider` append-only list; no `delete`/`purge` on history providers | `python/packages/core/agent_framework/_sessions.py:1048-1071`, `python/packages/core/agent_framework/_sessions.py:878-891` |
| Retirement — durable reset only | `AgentEntityStateProviderMixin.reset()` clears to fresh `DurableAgentState()` + `persist_state()` — manual reset, no TTL automation in code | `python/packages/durabletask/agent_framework_durabletask/_entities.py:77-82` |
| Tests — checkpoint roundtrip | `test_file_checkpoint_storage_roundtrip_*` for datetimes, dataclasses, tuples, WorkflowMessage, WorkflowEvent, full checkpoint | `python/packages/core/tests/workflow/test_checkpoint.py:993-1318` |
| Tests — checkpoint CRUD | `test_file_checkpoint_storage_save_and_load`, `list`, `delete`, `get_latest`, `directory_creation`, `corrupted_file` graceful handling | `python/packages/core/tests/workflow/test_checkpoint.py:774-992` |
| Tests — checkpoint chaining | `test_workflow_checkpoint_chaining_via_previous_checkpoint_id` asserts `previous_checkpoint_id` chain across supersteps | `python/packages/core/tests/workflow/test_checkpoint.py:285-337` |
| Tests — checkpoint resume | `test_sequential_checkpoint_resume_round_trip` + `test_sequential_checkpoint_runtime_only` + `test_sequential_checkpoint_runtime_overrides_buildtime` | `python/packages/orchestrations/tests/test_sequential.py:192-299` (same pattern in core workflow tests) |
| Docs — dure durable TTL gap | "The durable agents automatically maintain conversation history ... Without automatic cleanup, this state can accumulate" — TTL described as feature but no implementation found for workflow checkpoints; storage treated as trusted | `docs/features/durable-agents/durable-agents-ttl.md:5-15`, `python/packages/core/agent_framework/_workflows/_checkpoint_encoding.py:18-45` |

## Answers to Dimension Questions

### 1. What types of artifacts exist?

| Artifact | Primary schema / class | Storage backend(s) | Durability |
|----------|------------------------|---------------------|------------|
| **WorkflowCheckpoint** — full superstep snapshot (messages + shared state + pending HITL requests + iteration + metadata) | `WorkflowCheckpoint` dataclass `python/packages/core/agent_framework/_workflows/_checkpoint.py:30` with `to_dict`/`from_dict` `python/packages/core/agent_framework/_workflows/_checkpoint.py:90-116` | `InMemoryCheckpointStorage` `python/packages/core/agent_framework/_workflows/_checkpoint.py:192`, `FileCheckpointStorage` `python/packages/core/agent_framework/_workflows/_checkpoint.py:239` via `CheckpointStorage` protocol `python/packages/core/agent_framework/_workflows/_checkpoint.py:119` | Opt-in; created per superstep when checkpointing enabled `python/packages/core/agent_framework/_workflows/_runner.py:212` |
| **WorkflowRunResult / WorkflowEvent stream** — caller-visible outputs (`output`, `intermediate`), `request_info` HITL events, lifecycle `status`/`failed`, diagnostic events | `WorkflowRunResult` `python/packages/core/agent_framework/_workflows/_workflow.py:101`, `WorkflowEvent` `python/packages/core/agent_framework/_workflows/_events.py:146` | Transient (in-memory event queue `python/packages/core/agent_framework/_workflows/_runner_context.py:289`); persisted only if captured via checkpoint's `pending_request_info_events` or external sink | Ephemeral unless checkpointed or manually collected |
| **Shared Workflow State** (`State`) — superstep-committed dict with `_executor_state` reserved key | `State` `python/packages/core/agent_framework/_workflows/_state.py:6` with pending/commit semantics | Embedded inside `WorkflowCheckpoint.state`; also lives in `RunnerContext` per run | Durable only via checkpoint |
| **AgentSession + Conversation History** — per-session transcript (input + context + response messages) | `AgentSession` `python/packages/core/agent_framework/_sessions.py:746`, `Message`/`Content`/`ChatResponse`/`AgentResponse` `python/packages/core/agent_framework/_types.py:460` | `InMemoryHistoryProvider` (state dict) `python/packages/core/agent_framework/_sessions.py:814`, `FileHistoryProvider` (JSONL per session) `python/packages/core/agent_framework/_sessions.py:893`, Redis/Cosmos providers in sibling packages | Framework-managed (InMemory default) or file durable |
| **DurableAgentState entity** — long-lived agent conversationHistory (request/response entries with correlationId, orchestrationId, typed contents) | JSON schema `schemas/durable-agent-entity-state.json:195` + `DurableAgentState` class | `DurableTaskEntityStateProvider` via Durable Task `get_state`/`set_state` `python/packages/durabletask/agent_framework_durabletask/_entities.py:335-353` | Durable Task storage |
| **Executor-local state snapshot** — `on_checkpoint_save() -> dict` / `on_checkpoint_restore(state)` per executor | `Executor.on_checkpoint_save/restore` `python/packages/core/agent_framework/_workflows/_executor.py:493-517` | Aggregated into `State["_executor_state"]` then into `WorkflowCheckpoint.state` | Via checkpoint |
| **Harness file artifacts** — `AgentFileStore` / `FileMemoryProvider` / `FileAccessProvider` files (e.g., `working/`, `agent-file-memory/`) and tool outputs (`Content` data/uri/file_id, hosted_file, vector_store_id) | `AgentFileStore` (`InMemoryAgentFileStore`, `FileSystemAgentFileStore`), `FileMemoryProvider` | `FileSystemAgentFileStore` rooted at configurable dir | Explicit via harness |

No central "run artifact registry" exists; each artifact family lives in its own store.

### 2. How are artifacts named and stored?

- **Checkpoints**: filename `{checkpoint_id}.json` inside `storage_path` `python/packages/core/agent_framework/_workflows/_checkpoint.py:297`. `checkpoint_id` defaults to `uuid4` `python/packages/core/agent_framework/_workflows/_checkpoint.py:74`. `storage_path` is created `mkdir(parents=True, exist_ok=True)` `python/packages/core/agent_framework/_workflows/_checkpoint.py:278`. Write is atomic via `json.dump` to `.tmp` then `os.replace` `python/packages/core/agent_framework/_workflows/_checkpoint.py:317-321`. IDs are validated against traversal via `resolve()` + `is_relative_to` `python/packages/core/agent_framework/_workflows/_checkpoint.py:297-299`; companion tests in `python/packages/foundry_hosting/tests/test_responses.py:2934` confirm path-traversal rejection. Logical grouping key is `workflow_name` (builder name or `WorkflowBuilder-{uuid}`) `python/packages/core/agent_framework/_workflows/_workflow_builder.py:152`. Filtering is `list_checkpoints(workflow_name=...)` scanning `*.json` `python/packages/core/agent_framework/_workflows/_checkpoint.py:374`.
- **History files**: `{session_file_stem}{.jsonl}` under `storage_path` `python/packages/core/agent_framework/_sessions.py:1089`. Stem is literal `session_id` if safe (alnum+`._-`, not dot-prefixed, not Windows reserved) else `~session-{b64(session_id)}` `python/packages/core/agent_framework/_sessions.py:1094-1101`. Each line is one JSON `Message` via `dumps(message.to_dict())` single-line validated `python/packages/core/agent_framework/_sessions.py:1073-1091`. Root containment validated via `resolve().is_relative_to(_storage_root)` `python/packages/core/agent_framework/_sessions.py:1089-1091`. Concurrency via 64-stripe thread `+` per-loop async locks `python/packages/core/agent_framework/_sessions.py:916-1120`.
- **Durable entities**: keyed by `entity_id.key` (= thread_id) `python/packages/durabletask/agent_framework_durabletask/_entities.py:352-353`, state serialized via `DurableAgentState.to_dict()` `python/packages/durabletask/agent_framework_durabletask/_entities.py:71-75`.
- **Payload encoding inside checkpoints**: JSON-native types (str/int/float/bool/None, dict/list recursion) stay plain JSON; all other types (datetime, tuple, set, dataclass, custom) are `pickle+HIGHEST_PROTOCOL` then base64 with `{"__pickled__": "...", "__type__": "module:qualname"}` markers `python/packages/core/agent_framework/_workflows/_checkpoint_encoding.py:212-230`, wrapped in outer JSON checkpoint file `python/packages/core/agent_framework/_workflows/_checkpoint.py:315`.

No content-addressed storage; naming is UUID/random, not hash-derived (except `graph_signature_hash` for validation, not for file names).

### 3. Are artifacts versioned?

**Partially, weakly.**

- `WorkflowCheckpoint.version = "1.0"` `python/packages/core/agent_framework/_workflows/_checkpoint.py:88` — a static format version with no migration logic, no `from_dict` version dispatch, and no auto-increment on schema changes.
- `schemas/durable-agent-entity-state.json:209` has `schemaVersion` `^\d+\.\d+\.\d+$` at the entity root, but Python code does not enforce it beyond serialization.
- Stronger topology versioning: `Workflow.graph_signature_hash = sha256(canonical JSON of executors + edge_groups)` `python/packages/core/agent_framework/_workflows/_workflow.py:1077-1080`. On restore, `Runner.restore_from_checkpoint` rejects mismatched hash `python/packages/core/agent_framework/_workflows/_runner.py:275`. This prevents cross-topology replay but is not a user-visible artifact version.
- Temporal versioning is implicit chaining via `previous_checkpoint_id` `python/packages/core/agent_framework/_workflows/_checkpoint.py:51` forming a linked list; validated in `python/packages/core/tests/workflow/test_checkpoint.py:285` but not indexed.
- No semantic versioning for workflow outputs, no artifact-level `etag`/content hash, no coexistence of multiple format versions.

### 4. Can artifacts be linked to the run that produced them?

**Yes, but incompletely — by `workflow_name`, not by unique run instance.**

- Each checkpoint carries `workflow_name`, `checkpoint_id`, `previous_checkpoint_id`, `timestamp` (ISO8601), `iteration_count`, `graph_signature_hash`, `metadata` `python/packages/core/agent_framework/_workflows/_checkpoint.py:71-88`.
- Runner creates first checkpoint after start executor ("superstep 0") then after every superstep `python/packages/core/agent_framework/_workflows/_runner.py:96-97,143-144`, chaining via `previous_checkpoint_id` `python/packages/core/agent_framework/_workflows/_runner.py:84,144`.
- Discovery API is `list_checkpoints(workflow_name=...)` / `list_checkpoint_ids(workflow_name=...)` / `get_latest(workflow_name=...)` `python/packages/core/agent_framework/_workflows/_checkpoint.py:147-189` — it scans and filters by `workflow_name` only. There is **no** `list_checkpoints(run_id=...)` or per-run instance linkage. `Workflow.id` (UUID per build) `python/packages/core/agent_framework/_workflows/_workflow.py:329` is ephemeral and not persisted into checkpoints.
- To recover "all artifacts for run X", caller must: know `workflow.name` (stable if user supplied, else random `WorkflowBuilder-{uuid}` `python/packages/core/agent_framework/_workflows/_workflow_builder.py:152`), then filter `list_checkpoints` by `timestamp` window or `iteration_count`, or follow `previous_checkpoint_id` chain. This is demonstrated but not ergonomically wrapped in `python/packages/core/tests/workflow/test_checkpoint.py:323-336` sorting by `timestamp`.
- HITL linkage is explicit: `request_info` events contain `request_id` + `source_executor_id` + `request_type`/`response_type` `python/packages/core/agent_framework/_workflows/_events.py:296`; responses are `WorkflowMessage(type=RESPONSE, original_request_info_event=...)` `python/packages/core/agent_framework/_workflows/_runner_context.py:478-484`; pending events are checkpointed `python/packages/core/agent_framework/_workflows/_runner_context.py:388-389` and re-emitted on `apply_checkpoint` `python/packages/core/agent_framework/_workflows/_runner_context.py:414-426`.
- Agent session linkage is via `session_id` `python/packages/core/agent_framework/_sessions.py:770` and `FileHistoryProvider` filename stem `python/packages/core/agent_framework/_sessions.py:1089`.

**Gap**: given only a `WorkflowRunResult`, there is no embedded `checkpoint_id` list; caller must have separately tracked IDs or re-scan. Concurrent runs of same `workflow_name` interleave in the same directory with only `timestamp` to disambiguate.

### 5. How are artifacts retired?

- **Workflow checkpoints**: manual `CheckpointStorage.delete(checkpoint_id) -> bool` `python/packages/core/agent_framework/_workflows/_checkpoint.py:158`. `InMemoryCheckpointStorage.delete` removes from dict `python/packages/core/agent_framework/_workflows/_checkpoint.py:217`, `FileCheckpointStorage.delete` does `_validate_file_path` then `file_path.unlink()` `python/packages/core/agent_framework/_workflows/_checkpoint.py:392-410`. No bulk delete, no TTL, no retention policy, no scheduled GC, no reference counting. Corrupted/missing files raise `WorkflowCheckpointException` or are skipped with warning on `list` `python/packages/core/agent_framework/_workflows/_checkpoint.py:386-387,446-447`. Tests verify `delete` true/false and double-delete `python/packages/core/tests/workflow/test_checkpoint.py:233-251,849-867`.

- **History files**: no `delete` method on `HistoryProvider` abstract class `python/packages/core/agent_framework/_sessions.py:414-496`. `FileHistoryProvider` only appends `python/packages/core/agent_framework/_sessions.py:1048-1071`; retirement requires direct filesystem `unlink` outside the framework. `InMemoryHistoryProvider` holds forever in `state["messages"]` `python/packages/core/agent_framework/_sessions.py:890`.

- **Durable entities**: `AgentEntityStateProviderMixin.reset()` clears to fresh state `python/packages/durabletask/agent_framework_durabletask/_entities.py:77-82`; docs describe TTL aspiration `docs/features/durable-agents/durable-agents-ttl.md:5` ("Time-To-Live feature provides automatic cleanup of idle agent sessions") but no code implementation for workflow checkpoints; durable state itself has no expiry in `python/packages/durabletask/agent_framework_durabletask/_entities.py`.

- **State/queue cleanup**: `Runner._state.commit()` at superstep boundary `python/packages/core/agent_framework/_workflows/_runner.py:141` and executor state staging via `State.set`/`commit` `python/packages/core/agent_framework/_workflows/_runner.py:220-224`; `InProcRunnerContext.reset_for_new_run()` clears messages/queue `python/packages/core/agent_framework/_workflows/_runner_context.py:403-412` but shared `State` persists across runs `python/packages/core/agent_framework/_workflows/_workflow.py:538-569`.

**Overall**: retirement is ad-hoc and manual. There is no retention window, no max-count, no storage-budget guard, and no observability (no metrics for checkpoint count/size/age).

## Architectural Decisions

| Decision | Evidence | Rationale / Tradeoff |
|----------|----------|----------------------|
| **Pluggable CheckpointStorage Protocol vs concrete stores** | Protocol `python/packages/core/agent_framework/_workflows/_checkpoint.py:119` with `InMemory` `python/packages/core/agent_framework/_workflows/_checkpoint.py:192` and `File` `python/packages/core/agent_framework/_workflows/_checkpoint.py:239` | Keeps core `Runner` storage-agnostic and testable; but pushes durability choice to caller and fragments artifact locations (no unified catalog). |
| **Hybrid JSON + pickle+base64 encoding** | `encode_checkpoint_value` `python/packages/core/agent_framework/_workflows/_checkpoint_encoding.py:148` vs `_pickle_to_base64` `python/packages/core/agent_framework/_workflows/_checkpoint_encoding.py:281`, allowed via `RestrictedUnpickler` `python/packages/core/agent_framework/_workflows/_checkpoint_encoding.py:115` | Preserves human-readable JSON structure while keeping Python fidelity for datetimes/tuples/custom objects; security comment explicitly calls checkpoint storage "trusted data source" `python/packages/core/agent_framework/_workflows/_checkpoint_encoding.py:18-45` — defense-in-depth not a boundary, requiring access controls. |
| **Opt-in checkpointing at build-time or per-run with runtime override** | Builder param `checkpoint_storage` `python/packages/core/agent_framework/_workflows/_workflow_builder.py:96` → `InProcRunnerContext` `python/packages/core/agent_framework/_workflows/_runner_context.py:281` + `set_runtime_checkpoint_storage`/`clear_runtime_checkpoint_storage` `python/packages/core/agent_framework/_workflows/_runner.py:789-836`, `has_checkpointing()` `python/packages/core/agent_framework/_workflows/_runner_context.py:367` | Flexible for tests vs prod; but `Workflow.id` per-instance still not used as store partition, causing cross-run mingling by `workflow_name`. |
| **Graph signature hash for checkpoint compatibility** | `graph_signature` / `_hash_graph_signature` `python/packages/core/agent_framework/_workflows/_workflow.py:1015-1080`, validation in `restore_from_checkpoint` `python/packages/core/agent_framework/_workflows/_runner.py:275` | Cheap structural guard against topology drift; but no semantic diff or migration path when hash mismatches. |
| **Append-only JSONL history per session** | `FileHistoryProvider` `python/packages/core/agent_framework/_sessions.py:1058-1071` with `flush per line` | Simple, concurrent-safe (striped locks), debuggable; unbounded growth, no compaction/cleanup even though `COMPACTION_STATE_KEY` concepts exist elsewhere. |
| **Session state as plain dict with type registry** | `AgentSession.to_dict/from_dict` + `_serialize_state` handling `to_dict`/`BaseModel` `python/packages/core/agent_framework/_sessions.py:779-811`, `_STATE_TYPE_REGISTRY` `python/packages/core/agent_framework/_sessions.py:42` | Makes Session portable; but silent loss if type not registered. |
| **Atomic file replace for checkpoints** | `_write_atomic` via `.json.tmp` + `os.replace` `python/packages/core/agent_framework/_workflows/_checkpoint.py:317-321` | Crash-safe single-file write; still no directory-level transaction across chain. |

## Notable Patterns

- **Superstep-commit + per-superstep checkpoint**: pending state (`State._pending` `python/packages/core/agent_framework/_workflows/_state.py:28`) is staged via `set()` then `commit()` at `Runner._state.commit()` `python/packages/core/agent_framework/_workflows/_runner.py:141` and again before checkpoint `python/packages/core/agent_framework/_workflows/_runner.py:224`, ensuring checkpoints contain only committed state plus `pending_request_info_events`.
- **Defense-in-depth deserialization**: `_type_to_key` + `_verify_type` `python/packages/core/agent_framework/_workflows/_checkpoint_encoding.py:257-310` allowlist built-ins + `agent_framework.*` + `openai.types.*` + caller `allowed_checkpoint_types`; caller must pass `module:qualname` strings `python/packages/core/agent_framework/_workflows/_checkpoint.py:249-279`.
- **Striped locking for file history**: 64 thread locks + per-event-loop async locks keyed by `hash(file_path) % 64` `python/packages/core/agent_framework/_sessions.py:916-1120` to allow concurrent sessions without global lock.
- **Trusted-storage assumption documented in code**: 4-point checklist (access-controlled, never from HTTP, restrict allowed_types, sanitize) `python/packages/core/agent_framework/_workflows/_checkpoint_encoding.py:27-39`.
- **Chained checkpoints as implicit run lineage**: `previous_checkpoint_id` lineage validated in tests `python/packages/core/tests/workflow/test_checkpoint.py:323-337` but never exposed as a `get_chain(checkpoint_id)` API.

## Tradeoffs

- **Fidelity vs safety**: pickle gives exact Python object roundtrip (validated by `python/packages/core/tests/workflow/test_checkpoint.py:366-711`) but at cost of inherent code-execution risk during restore; mitigated only by allowlist and trusted storage, not by pure JSON.
- **Simplicity vs discoverability**: scanning `*.json` and `list_checkpoints(workflow_name=...)` `python/packages/core/agent_framework/_workflows/_checkpoint.py:372-390` is simple without a DB, but O(n) directory scan, no pagination, no checkpoint size/age metrics, and no run-scoped listing.
- **In-memory convenience vs durability**: `InMemoryCheckpointStorage` `python/packages/core/agent_framework/_workflows/_checkpoint.py:192` is zero-config for tests/orchestrations `python/packages/orchestrations/tests/test_sequential.py:23` but loses all artifacts on process exit; callers must remember to supply `FileCheckpointStorage`.
- **Workflow_name grouping vs run isolation**: reuse of `workflow_name` across builds aids filtering but conflates concurrent/replayed runs; missing `run_id` partitioning forces external correlation.
- **Manual retirement vs automation**: explicit `delete` keeps control with caller but guarantees unbounded growth (`FileHistoryProvider` JSONL append-only) without ops tooling.

## Failure Modes / Edge Cases

| Mode | Evidence | Impact |
|------|----------|--------|
| Untrusted pickle payload | `decode_checkpoint_value` uses `pickle.loads` unless `allowed_types` supplied; docs warn "not a security boundary" `python/packages/core/agent_framework/_workflows/_checkpoint_encoding.py:18-45` | Malicious checkpoint file can execute code on `load`. |
| Path traversal via checkpoint_id/session_id | Mitigated by `resolve().is_relative_to()` `python/packages/core/agent_framework/_workflows/_checkpoint.py:297`, `python/packages/core/agent_framework/_sessions.py:1090` and b64 encoding `python/packages/core/agent_framework/_sessions.py:1100`; but traversal via symlinked `storage_path` itself not re-validated after `mkdir`. | Escape prevented in tested paths `python/packages/foundry_hosting/tests/test_responses.py:2934`. |
| Corrupted JSON / wrong type on restore | `list_checkpoints` skips corrupted file with `logger.warning` `python/packages/core/agent_framework/_workflows/_checkpoint.py:386-388`; `decode` raises `WorkflowCheckpointException` `python/packages/core/agent_framework/_workflows/_checkpoint_encoding.py:298-306`; `Runner.restore` re-raises `WorkflowCheckpointException` `python/packages/core/agent_framework/_workflows/_runner.py:296-300` | Partial store can silently hide broken checkpoint; caller using `get_latest` may get older one. |
| Graph topology drift | `restore_from_checkpoint` raises if `graph_signature_hash` mismatched `python/packages/core/agent_framework/_workflows/_runner.py:275` | Resume across workflow refactors always hard-fails with no migration. |
| Clock skew / duplicate timestamps | `get_latest` uses `max(... key=lambda cp: datetime.fromisoformat(cp.timestamp))` `python/packages/core/agent_framework/_workflows/_checkpoint.py:230,424` | Two rapid checkpoints can tie or mis-order if system clock jumps. |
| Unbounded file growth | `FileHistoryProvider` never truncates, `FileCheckpointStorage` never compacts chain | Disk fills over long-lived services; no warning. |
| Concurrent deletes during list | Scan+decode loop catches per-file exceptions `python/packages/core/agent_framework/_workflows/_checkpoint.py:374-388` but delete is not transactional | Race can cause checkpoint missed during enumeration then later `load` fails. |
| History file locking across processes | Locks are `threading.Lock` + `asyncio.Lock` in-process only `python/packages/core/agent_framework/_sessions.py:918-1115` | Multi-process hosts (e.g., Azure Functions) can interleave JSONL writes. |

## Future Considerations

- Introduce a **Run-scoped artifact catalog**: persist `run_id` (UUID per `workflow.run()` call) into `WorkflowCheckpoint.metadata` and `WorkflowRunResult`, add `CheckpointStorage.list_by_run(run_id)` and `workflow.list_runs()` so the answer to "every artifact for run X?" is a single query.
- Add **retention policy**: `FileCheckpointStorage(max_age, max_count, max_bytes)` with background `purge()` and/or `delete` by predicate; mirror durable-agent TTL doc `docs/features/durable-agents/durable-agents-ttl.md:5` with real implementation, and add same for `FileHistoryProvider` (e.g., `rotate`/`compact` using existing compaction keys `python/packages/core/agent_framework/_compaction.py:34`).
- Add **storage inventory API + observability**: emit OpenTelemetry metrics for checkpoint bytes/count/age and history file sizes; expose `checkpoint_stats(workflow_name)` for ops dashboards.
- Normalize **naming/partitioning**: store checkpoints under `storage_path/{workflow_name}/{run_id}/{checkpoint_id}.json` to isolate concurrent runs without prefix scans, and make session files sharded by date/run.
- Harden versioning: version-tag encoded payloads and implement `WorkflowCheckpoint.from_dict` dispatch over `version`, plus content hash for deduplication.

## Questions / Gaps

- No evidence of **automatic GC/TTL for workflow checkpoints** — searched `python/packages/core/agent_framework/_workflows/_checkpoint.py:239` for delete only; doc `docs/features/durable-agents/durable-agents-ttl.md:5` mentions TTL but no workflow-checkpoint code implements it.
- No **central artifact manifest** linking a single `workflow.run()` to all checkpoint IDs it produced — `WorkflowRunResult` `python/packages/core/agent_framework/_workflows/_workflow.py:101` holds events but not checkpoint IDs; must scrape `storage.list_checkpoints(workflow_name)` externally.
- No **encryption/secret handling** for checkpoint files — encoding doc calls storage "trusted" `python/packages/core/agent_framework/_workflows/_checkpoint_encoding.py:18` but code writes plaintext JSON+B64 pickle.
- No **history truncation** — `HistoryProvider` interface `python/packages/core/agent_framework/_sessions.py:414` lacks `delete`/`clear`; retention must be out-of-band.
- `.NET parity` for artifact lifecycle (e.g., `CosmosCheckpointStore`, `ValkeyChatHistoryProvider`) was not enumerated file-by-file in this study; the `dotnet/src/` side likely mirrors Python but was only shallowly globbed — deeper audit needed to confirm naming/storage parity.

---

Generated by `Dimension 16.01: Artifact Lifecycle` against `agent-framework`.
