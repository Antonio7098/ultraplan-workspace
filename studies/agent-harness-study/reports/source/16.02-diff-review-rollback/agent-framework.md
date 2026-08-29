# Source Analysis: agent-framework

## Dimension 16.02: Diff, Review, and Rollback

### Source Info

| Field | Value |
|-------|-------|
| Name | agent-framework |
| Path | `studies/agent-harness-study/sources/agent-framework` |
| Language / Stack | Python + .NET (Microsoft Agent Framework) |
| Analyzed | 2026-08-28 |

## Summary

Microsoft Agent Framework is an agent/workflow runtime, not an artifact registry. It has no first-class artifact object with versioned diffs, review gates, comments, or rollback. The closest artifact-like constructs are `WorkflowCheckpoint`/`AgentSession` state. For those, the framework provides a mature checkpoint chain (save/load/list/delete/get_latest, previous_checkpoint_id linking, file+memory storage, graph-signature validation, restricted deserialization) and streaming `WorkflowEvent` status timeline that enable manual restore-as-rollback. Runtime human review exists only for tool/approval escalation (`function_approval_request/response`, `ToolApprovalMiddleware`, `RequestInfo` executor) and workflow HITL continuation — not for reviewing artifact definitions or bad deploys. No diff generators, no annotation/comment models, and no unified artifact change log linking artifact versions to run IDs were found. A bad workflow/agent definition cannot be reverted with full audit trail via framework alone; restoration is fragile (graph hash must match, caller must pick checkpoint id).

## Rating

**4 / 10 — Present but inconsistent, weakly documented, fragile**

Rationale: Checkpoint-based rollback is explicit, tested, and safeguarded (hash validation, path-traversal checks, restricted pickle allowlist, atomic file writes) but covers only workflow execution state. Diff, artifact comparison, artifact-level review/approval, comments/annotations on artifacts, and artifact↔run traceability are absent as framework primitives. Review workflow is limited to per-execution tool approval/HITL, not artifact governance. No `deepdiff` or equivalent diff generator is wired to checkpoints; no comment model exists.

## Evidence Collected

Every entry MUST include a file path with line numbers. Format: `path/to/file.ts:NN`.

| Area | Evidence | File:Line |
|------|----------|-----------|
| Diff generators | No diff generator found. Search for `diff`, `compare`, `deepdiff` yields only `deepdiff` as dev lockfile dependency unused for artifacts; no `DiffGenerator` class. `grep diff` across source hits only git diff in scripts and generic comments. | `python/uv.lock:2142`, `python/scripts/run_tasks_in_changed_packages.py:27`, `docs/specs/001-foundry-sdk-alignment.md:13` (no diff impl) |
| Checkpoint model (artifact surrogate) | `WorkflowCheckpoint` dataclass defines versioned artifact snapshot: `workflow_name`, `graph_signature_hash`, `checkpoint_id`, `previous_checkpoint_id`, `timestamp`, `messages`, `state`, `pending_request_info_events`, `iteration_count`, `metadata`, `version` | `python/packages/core/agent_framework/_workflows/_checkpoint.py:30-88` |
| Checkpoint storage protocol | `CheckpointStorage` protocol: `save`, `load`, `list_checkpoints`, `delete`, `get_latest`, `list_checkpoint_ids` — explicit version listing and deletion (rollback primitive) | `python/packages/core/agent_framework/_workflows/_checkpoint.py:119-189` |
| File checkpoint implementation | `FileCheckpointStorage` — JSON+pickle+base64 encoding, atomic write via `.tmp` + `os.replace`, `_validate_file_path` prevents traversal via `is_relative_to`, `allowed_checkpoint_types` allowlist, `list_checkpoints`/`list_checkpoint_ids` filtered by `workflow_name` | `python/packages/core/agent_framework/_workflows/_checkpoint.py:239-450` |
| In-memory checkpoint implementation | `InMemoryCheckpointStorage` with `copy.deepcopy` isolation, `get_latest` via `max(timestamp)`, chain via `previous_checkpoint_id` handling | `python/packages/core/agent_framework/_workflows/_checkpoint.py:192-237` |
| Checkpoint encoding safeguards | Hybrid JSON + `RestrictedUnpickler` with `_BUILTIN_ALLOWED_TYPE_KEYS`, `agent_framework.*`, `openai.types.*` allowlist, `_verify_type` mismatch detection, security caveat doc that storage must be trusted | `python/packages/core/agent_framework/_workflows/_checkpoint_encoding.py:14-45`, `python/packages/core/agent_framework/_workflows/_checkpoint_encoding.py:72-145` |
| Rollback / restore handler | `Runner.restore_from_checkpoint` validates `graph_signature_hash` match, clears state via `State.clear+import_state`, restores executor states and `RunnerContext.apply_checkpoint`, marks resumed iteration | `python/packages/core/agent_framework/_workflows/_runner.py:240-300` |
| Runner checkpoint chaining | `Runner._create_checkpoint_if_enabled` saves executor states then `State.commit()` then `create_checkpoint(workflow_name, graph_signature_hash, state, previous_checkpoint_id, iteration)` after each superstep; chaining test asserts `previous_checkpoint_id == prev.checkpoint_id` | `python/packages/core/agent_framework/_workflows/_runner.py:212-238`, `python/packages/core/tests/workflow/test_checkpoint.py:285-336` |
| RunnerContext checkpoint API | `InProcRunnerContext.create_checkpoint` captures `messages`, `state.export_state()`, `pending_request_info_events`; `load_checkpoint`/`apply_checkpoint` restores messages and re-queues request_info events | `python/packages/core/agent_framework/_workflows/_runner_context.py:370-427` |
| Workflow run rollback entry | `Workflow.run(checkpoint_id=..., checkpoint_storage=...)` and `responses` continuation; `Workflow._execute_with_message_or_checkpoint` requires `message` xor `checkpoint_id`; `restore_from_checkpoint` before replay | `python/packages/core/agent_framework/_workflows/_workflow.py:630-670`, `python/packages/core/agent_framework/_workflows/_workflow.py:712-770` |
| .NET parity | `Workflow` with `CheckpointManager.CreateJson(FileSystemJsonCheckpointStore)` / `CreateInMemory`, demo resume via `ResumeStreamingAsync(workflow, LastCheckpoint, checkpointManager)` and `SuperStepCompletedEvent` checkpoint capture | `dotnet/src/Shared/Workflows/Execution/WorkflowRunner.cs:42-98` |
| Review workflow — tool approval | `ToolApprovalMiddleware` (experimental) requires `AgentSession`, session-backed `ToolApprovalState`/`ToolApprovalRule` with `server_label` boundary, heuristic `auto_approval_rules`, queued `function_approval_request` -> deferred replay, `create_always_approve_tool_response` helpers | `python/packages/core/agent_framework/_harness/_tool_approval.py:86-158`, `python/packages/core/agent_framework/_harness/_tool_approval.py:345-618` |
| Review workflow — approval content types | Unified `Content` types `function_approval_request` / `function_approval_response` with `id`, `function_call`, `approved`, `additional_properties.tool_approval` scoping, conversion `to_function_approval_response` | `python/packages/core/agent_framework/_types.py:338-362`, `python/packages/core/agent_framework/_types.py:1212-1253`, `python/packages/core/agent_framework/_types.py:1286-1302` |
| Review workflow — workflow HITL | `WorkflowEvent.request_info` with `request_id`, `source_executor_id`, `request_type/response_type`; `RunnerContext.add_request_info_event`/`send_request_info_response`/`get_pending_request_info_events`; Workflow `run(responses={req_id: data})` with `try_coerce_to_type` and `is_instance_of` validation; status `IDLE_WITH_PENDING_REQUESTS` | `python/packages/core/agent_framework/_workflows/_events.py:114-145`, `python/packages/core/agent_framework/_workflows/_runner_context.py:244-260`, `python/packages/core/agent_framework/_workflows/_workflow.py:936-960` |
| Change logs / history | `WorkflowEvent` types: `started`, `status`, `output`, `intermediate`, `request_info`, `superstep_started/completed`, `executor_invoked/completed/failed`; `WorkflowRunResult` filters status vs data events, `get_outputs`/`get_intermediate_outputs`/`get_request_info_events`/`status_timeline` | `python/packages/core/agent_framework/_workflows/_events.py:104-131`, `python/packages/core/agent_framework/_workflows/_workflow.py:101-164` |
| Session change log | `HistoryProvider` abstract `get_messages`/`save_messages` with `InMemoryHistoryProvider` (session.state["messages"]) and `FileHistoryProvider` (JSONL per session `${session_id}.jsonl`), but no version diff or rollback | `python/packages/core/agent_framework/_sessions.py:413-535`, `python/packages/core/agent_framework/_sessions.py:814-847`, `python/packages/core/agent_framework/_sessions.py:893-1025` |
| Session serialization / traceability | `AgentSession.to_dict`/`from_dict` with `session_id`, `service_session_id`, `state` serialized via `_serialize_state`/`_STATE_TYPE_REGISTRY`; `PerServiceCallHistoryPersistingMiddleware` handles per-call persistence & local sentinel `agent_framework_local_history_persistence` | `python/packages/core/agent_framework/_sessions.py:746-811`, `python/packages/core/agent_framework/_sessions.py:570-744` |
| Observability trace link | `Runner._run_workflow_with_tracing` creates workflow span `OtelAttr.WORKFLOW_RUN_SPAN` with `workflow.id/name`, emits `WorkflowEvent.started/status` and `WorkflowRunState` tracking; not an artifact↔run ledger | `python/packages/core/agent_framework/_workflows/_workflow.py:480-542` |
| Annotation model | `Content.annotations: Sequence[Annotation]` (TypedDict `citation` with `url/file_id/tool_name/snippet`) — provider grounding metadata, not human review comments on artifacts | `python/packages/core/agent_framework/_types.py:374-385`, `python/packages/core/agent_framework/_types.py:462-534` |
| Declarative artifacts | `declarative-agents/` YAML workflow definitions loadable as `Workflow` — no versioning, diff, review, or rollback helpers found | `declarative-agents/workflow-samples/README.md:1-6` |
| Tests for rollback | `test_checkpoint.py` covers save/load/list/delete/get_latest/chaining/roundtrip; `test_checkpoint_validation.py` covers graph hash + sub-workflow mismatch rejection; `test_checkpoint_unrestricted_pickle.py` covers allowlist | `python/packages/core/tests/workflow/test_checkpoint.py:160-223`, `python/packages/core/tests/workflow/test_checkpoint_validation.py:34-165`, `python/packages/core/tests/workflow/test_checkpoint_unrestricted_pickle.py:45-182` |

## Answers to Dimension Questions

**1. Can artifacts be compared?**
No. No diff/comparison generator exists for workflow checkpoints, agent definitions, declarative YAML, or session state. `WorkflowCheckpoint.to_dict` exists for serialization (`python/packages/core/agent_framework/_workflows/_checkpoint.py:90-98`) but there is no `diff(checkpoint_a, checkpoint_b)` utility. `deepdiff` appears only in `python/uv.lock:2142` as a transitive dependency, never imported for checkpoint comparison. Evidence search for `DiffGenerator`, `compare`, `diff` in `python/packages/core` yields zero artifact diff code. At most operators could manually deserialize two JSON checkpoint files and compare externally.

**2. Is there a review workflow?**
Partial — runtime execution review only, not artifact governance. Two mechanisms:
- Per-tool approval: `ToolApprovalMiddleware` queues `function_approval_request` Content, auto-approves via `ToolApprovalRule`/`auto_approval_rules`, stores standing rules in `AgentSession.state` keyed by `server_label` (`python/packages/core/agent_framework/_harness/_tool_approval.py:345-408`). Requires `AgentSession`, supports streaming via `ResponseStream` interception.
- Workflow HITL: executors emit `WorkflowEvent.request_info` (`python/packages/core/agent_framework/_workflows/_events.py:296-312`), workflow enters `IDLE_WITH_PENDING_REQUESTS`, caller resumes via `workflow.run(responses={request_id: data})` (`python/packages/core/agent_framework/_workflows/_workflow.py:724-730`). No PR-style artifact review, no approver roles, no artifact approval audit beyond session state.

**3. Can artifacts be rolled back?**
Partially — workflow execution state can be restored, declarative/agent definitions cannot.
- Rollback primitive: `CheckpointStorage.delete` + `load` + `Workflow.run(checkpoint_id=...)` restores prior `messages`, `state`, `pending_request_info_events`, `iteration_count` (`python/packages/core/agent_framework/_workflows/_runner.py:240-291`, `python/packages/core/agent_framework/_workflows/_runner_context.py:414-427`). `previous_checkpoint_id` chains checkpoints (`python/packages/core/agent_framework/_workflows/_checkpoint.py:75`) enabling manual traversal to prior version.
- Guards/fragility: `graph_signature_hash` validation rejects mismatched topology (`python/packages/core/agent_framework/_workflows/_runner.py:274-279`); no automatic bad-artifact reversion; caller must track correct `checkpoint_id`/`workflow_name`; `FileCheckpointStorage.delete` is plain `unlink` with no soft-delete, no tombstone, no immutable audit log (`python/packages/core/agent_framework/_workflows/_checkpoint.py:392-410`). No rollback for `Agent` instructions, skills, or declarative YAML version history.

**4. Are artifact changes traceable to runs?**
Weakly. Workflow run traces are via `WorkflowEvent` status timeline and `WorkflowRunResult` (`python/packages/core/agent_framework/_workflows/_workflow.py:121-164`) and OTel span (`python/packages/core/agent_framework/_workflows/_workflow.py:512-520`), but there is no artifact version register. `WorkflowCheckpoint` stores `workflow_name`, `graph_signature_hash`, `checkpoint_id`, `timestamp`, `iteration_count`, `previous_checkpoint_id` (`python/packages/core/agent_framework/_workflows/_checkpoint.py:71-88`) — allows associating a checkpoint with the workflow definition that created it, not with a deploy/change event. `AgentSession` stores `session_id`/`service_session_id` and history via `HistoryProvider` (`python/packages/core/agent_framework/_sessions.py:746-811`), but does not record which artifact version produced which messages. No `artifact_id -> run_id` mapping, no change log table.

## Architectural Decisions

- **Checkpoint-as-artifact surrogate:** The framework treats workflow execution snapshots as the versioned entity instead of a registry of agent/workflow definitions. `WorkflowCheckpoint` + `CheckpointStorage` protocol is the sole versioned store (`python/packages/core/agent_framework/_workflows/_checkpoint.py:119-189`). Tradeoff: durable restartability with low operational overhead, but no deployment artifact governance.
- **Chain via `previous_checkpoint_id` + `graph_signature_hash`:** Enables safe restore only to same topology (`python/packages/core/agent_framework/_workflows/_runner.py:275`). Tradeoff: prevents silent topology drift but makes urgent rollback of a breaking graph change impossible without rebuilding original graph.
- **Hybrid JSON+restricted pickle:** Preserves Python object fidelity while limiting deserialization attack surface (`python/packages/core/agent_framework/_workflows/_checkpoint_encoding.py:72-112`). Framework + OpenAI types always allowed; user types require `allowed_checkpoint_types` (`python/packages/core/agent_framework/_workflows/_checkpoint.py:262-279`). Tradeoff: `getattr` remains allowlisted for enum restoration, so storage must still be treated as trusted (`python/packages/core/agent_framework/_workflows/_checkpoint_encoding.py:19-44`).
- **Session-backed approval state vs. external review service:** `ToolApprovalMiddleware` stores `rules`/`queued_approval_requests`/`collected_approval_responses` in `AgentSession.state` (`python/packages/core/agent_framework/_harness/_tool_approval.py:160-217`) for portability across `File`/`InMemory` providers. Tradeoff: simple, testable, no external dependency, but no centralized approval audit or multi-reviewer workflow.
- **HITL via `request_info` events as user-visible artifacts:** Pending requests durably checkpointed (`python/packages/core/agent_framework/_workflows/_runner_context.py:385-390`), enabling resume after process restart. Tradeoff: elegant continuation model, but couples review to execution supersteps.

## Notable Patterns

- **Atomic file commit + traversal guard:** `_write_atomic` via `tmp` + `os.replace` and `_validate_file_path` with `resolve().is_relative_to()` (`python/packages/core/agent_framework/_workflows/_checkpoint.py:282-326`) — production-grade filesystem safety rarely seen in agentic frameworks.
- **Restricted unpickler as defense-in-depth:** `_RestrictedUnpickler.find_class` allowlist (`python/packages/core/agent_framework/_workflows/_checkpoint_encoding.py:115-145`) with explicit error message directing to `allowed_checkpoint_types`.
- **Previous-checkpoint chaining invariant:** Tests assert chain correctness (`python/packages/core/tests/workflow/test_checkpoint.py:285-336`) — provides linear audit trail for workflow progress, though not an immutable log (deletable).
- **Request-response correlation by `request_id`:** `send_request_info_response` pops pending map and validates `response_type` (`python/packages/core/agent_framework/_workflows/_runner_context.py:457-486`), validated again in `Workflow._send_responses_internal` (`python/packages/core/agent_framework/_workflows/_workflow.py:936-960`).

## Tradeoffs

- **Durability vs. governance:** Optimizes for crash recovery (superstep checkpoints, session history) over artifact lifecycle governance (versioning, diff, approval gates).
- **Python fidelity vs. portability:** Pickle fidelity preserves closures/tuples/sets exactly (tested in `test_checkpoint.py:418-445`, `1078-1113`) at cost of `.NET` interop and long-term schema evolution risk; mitigated by `version: "1.0"` field but no migration logic.
- **Explicit restore vs. automatic rollback:** Caller explicitly selects `checkpoint_id` — avoids accidental rollback but provides no `rollback_to_previous_good` one-click primitive; operator must enumerate via `list_checkpoints`/`get_latest`.
- **Local vs. file history:** `InMemoryHistoryProvider` stores messages in `session.state` (portable, ephemeral) vs `FileHistoryProvider` JSONL (durable, append-only) — no compaction, no dedup, no diff.

## Failure Modes / Edge Cases

- **Graph drift blocks rollback:** Changing workflow topology (executor IDs/types/edges) changes `graph_signature_hash`; `restore_from_checkpoint` raises `WorkflowCheckpointException: Workflow graph has changed` (`python/packages/core/agent_framework/_workflows/_runner.py:276-279`); validated by `test_checkpoint_validation.py:34-165`.
- **Silent storage trust assumption:** Despite allowlist, `builtins:getattr` and `copyreg:_reconstructor` remain allowed (`python/packages/core/agent_framework/_workflows/_checkpoint_encoding.py:95-97`); compromised file storage can still achieve code execution — documented as must-be-trusted (`python/packages/core/agent_framework/_workflows/_checkpoint_encoding.py:27-44`).
- **No atomic multi-checkpoint transaction:** Each superstep creates one checkpoint; crash between `save` and `State.commit` may leave partial chain; previous checkpoint remains but `previous_checkpoint_id` may dangle if intermediate delete races.
- **Deleted checkpoint = lost audit trail:** `FileCheckpointStorage.delete` permanently `unlink`s file (`python/packages/core/agent_framework/_workflows/_checkpoint.py:402-410`); `InMemoryCheckpointStorage.delete` `del`s entry; no soft-delete or append-only log, so forensic trail can be erased.
- **Corrupted JSON silently skipped in listings:** `list_checkpoints` catches `Exception` and logs `warning` then skips file (`python/packages/core/agent_framework/_workflows/_checkpoint.py:385-387`, `443-447`); operator may not notice truncated history.
- **Approval state replay requires session continuity:** `ToolApprovalMiddleware` queued requests stored in `AgentSession.state` (`python/packages/core/agent_framework/_harness/_tool_approval.py:250-277`); losing session (no persistence, `session=None`) loses pending approvals; no cross-session durable queue.
- **No diff means bad diff review:** Without diff, reviewers cannot see what changed between checkpoint `N` and `N+1`; bad state may be restored blindly.
- **Unbounded checkpoint growth:** No TTL, compaction, or retention policy in `FileCheckpointStorage`; `list_checkpoints` scans entire directory via `glob("*.json")` (`python/packages/core/agent_framework/_workflows/_checkpoint.py:374`), which scales poorly.

## Future Considerations

- Add `WorkflowCheckpoint.diff(other, allowed_types=...)` or standalone `checkpoint_diff` utility leveraging `deepdiff` (already in lockfile) to surface state/message deltas; expose via DevUI.
- Introduce versioned artifact registry for declarative agents/workflows (e.g., `WorkflowDefinitionVersion {id, yaml_hash, graph_signature_hash, created_by, checkpoint_id}`) with immutable append log and `approve`/`reject` state machine separate from runtime `ToolApprovalMiddleware`.
- Add comment/annotation model (`CheckpointComment {checkpoint_id, author, body, created_at}`) stored adjacent to checkpoints, surfaced in DevUI timeline.
- Harden audit trail: soft-delete + tombstone file, `list_checkpoint_ids` should return deleted markers; add `get_history(workflow_name)` returning ordered chain validated via `previous_checkpoint_id`.
- Traceability: emit `run_id` (OTel trace/span ID) into `WorkflowCheckpoint.metadata` and `AgentSession.state`, and add `DecisionRecord` linking `artifact_version -> run_id -> checkpoint_id` for post-mortems.
- Provide `rollback(workflow_name, to_checkpoint_id, storage)` helper that validates hash, creates compensating checkpoint, and records `rollback_reason` in metadata.

## Questions / Gaps

- No evidence found for artifact-level comment/annotation persistence. Search for `comment`, `annotation` in `python/packages/core` yields only `Content.annotations` (LLM citation) — confirm intentional omission vs. roadmap gap.
- No evidence found for declarative artifact version store. Does Foundry hosting layer (outside this source slice) provide artifact registry with diff/review? Not inspectable under isolation rule.
- Graph signature equality is hash-based (`python/packages/core/agent_framework/_workflows/_workflow.py:1077-1080`); what constitutes breaking vs non-breaking graph change for rollback purposes is undocumented.
- `AgentSession` state serialization registry `_STATE_TYPE_REGISTRY` (`python/packages/core/agent_framework/_sessions.py:43-90`) is extensible but not versioned — how are state schema migrations handled across rollback?

---

Generated by `Dimension 16.02: Diff, Review, and Rollback` against `agent-framework`.
