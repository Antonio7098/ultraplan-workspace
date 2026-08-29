# Source Analysis: crewai

## Dimension 23.03: Responsibility and Accountability Model

### Source Info

| Field | Value |
|-------|-------|
| Name | crewai |
| Path | `studies/agent-harness-study/sources/crewai` |
| Language / Stack | Python / Pydantic, CrewAI multi-agent (Crew/Agent/Task/Flow) |
| Analyzed | 2026-08-29 |

## Summary

CrewAI assigns operational responsibility through executable bindings rather than a declared policy. A `Task` must declare `agent` (`lib/crewai/src/crewai/task.py:163`), an `Agent`/`Crew` carries a `SecurityConfig.fingerprint` (`lib/crewai/src/crewai/security/security_config.py:39`), and a run-scoped `execution_uuid` contextvar (`lib/crewai/src/crewai/execution.py:22`) correlates all nested kickoffs. At runtime every action is re-attributed via the event bus: `BaseEvent` fields `source_fingerprint`, `source_type`, `agent_id`, `task_id`, `event_id` (`lib/crewai/src/crewai/events/base_events.py:71`) are populated by type-specific helpers for tasks (`lib/crewai/src/crewai/events/types/task_events.py:96`), agents (`lib/crewai/src/crewai/events/types/agent_events.py:124`), crews (`lib/crewai/src/crewai/events/types/crew_events.py:22`) and tools (`lib/crewai/src/crewai/events/types/tool_usage_events.py:47`). Tool calls are the most attributed: `ToolUsage` emits `ToolUsageStartedEvent`/`ToolUsageFinishedEvent`/`ToolUsageErrorEvent`/`ToolFailureDetectedEvent` with `tool_name`, `tool_args`, `tool_class`, `failure` and fingerprint metadata (`lib/crewai/src/crewai/tools/tool_usage.py:273`, `lib/crewai/src/crewai/tools/tool_failure.py:304`). Human decisions are recorded as first-class Flow events (`HumanFeedbackRequestedEvent`/`HumanFeedbackReceivedEvent`, `FlowInputRequestedEvent`/`FlowInputReceivedEvent`) and as `Task.human_input`/`@human_feedback` decorator state, but model outputs carry no disclaimer or model-attribution watermark.

## Rating

**Score: 5 / 10 — Present but inconsistent, weakly documented, fragile**

Rationale: identity (Fingerprint UUID, execution_uuid, event graph) and tool/human attribution are implemented and tested, but policy-level accountability is undocumented, output attribution is limited to role-string and lineage edges, there is no output disclaimer, and no durable accountability/chain-of-custody document. The system can answer "which agent/task/tool produced this?" via events, but cannot answer "who is ultimately responsible?" in organizational/policy terms.

## Evidence Collected

Every entry MUST include a file path with line numbers. Format: `path/to/file.ts:NN`.

| Area | Evidence | File:Line |
|------|----------|-----------|
| Policy attribution — Task→Agent binding | `Task` docstring: "Each task must have a description, an expected output and an agent responsible for execution." Field `agent: BaseAgent \| None = Field(description="Agent responsible for execution the task.")` | `lib/crewai/src/crewai/task.py:123`, `lib/crewai/src/crewai/task.py:163` |
| Policy attribution — Crew/Flow execution identity | `execution_uuid` ContextVar owned by outermost `kickoff`; `begin_execution`/`end_execution`/`get_execution_uuid`/`set_execution_uuid` for OSS/Enterprise correlation | `lib/crewai/src/crewai/execution.py:22`, `lib/crewai/src/crewai/execution.py:27`, `lib/crewai/src/crewai/execution.py:51` |
| Policy attribution — Component fingerprint | `Fingerprint` UUID + metadata for tracking/auditing; `SecurityConfig.fingerprint` on Agent/Task/Crew; `to_dict`/`from_dict` persistence | `lib/crewai/src/crewai/security/fingerprint.py:42`, `lib/crewai/src/crewai/security/security_config.py:39`, `lib/crewai/src/crewai/security/fingerprint.py:123` |
| Policy attribution — BaseEvent correlation | `BaseEvent` carries `source_fingerprint`, `source_type`, `fingerprint_metadata`, `agent_id`, `task_id`, `event_id`, `parent_event_id`, `triggered_by_event_id`, `started_event_id`, `emission_sequence` | `lib/crewai/src/crewai/events/base_events.py:66` |
| Policy attribution — Directed event record | `EventRecord` stores events as graph nodes with edges `parent/child/trigger/triggered_by/next/previous/started/completed_by` for replay/audit | `lib/crewai/src/crewai/state/event_record.py:52`, `lib/crewai/src/crewai/state/event_record.py:110` |
| Output attribution — TaskOutput / CrewOutput | `TaskOutput.agent: str` = agent role that executed task; `TaskOutput.messages`, `tool_failures`; `CrewOutput.tasks_output` preserves per-task lineage | `lib/crewai/src/crewai/tasks/task_output.py:43`, `lib/crewai/src/crewai/crews/crew_output.py:24` |
| Output attribution — No disclaimer/watermark | Exhaustive grep for `disclaimer`, `watermark`, `attribution` yields zero framework disclaimer injection; LLM output returned verbatim via `TaskOutput.raw` / `CrewOutput.raw` | `lib/crewai/src/crewai/tasks/task_output.py:36`, `lib/crewai/src/crewai/crews/crew_output.py:17` |
| Tool attribution — Invocation config | `ToolUsage._build_fingerprint_config` injects `agent_fingerprint`/`task_fingerprint` into tool `config` without polluting schema args | `lib/crewai/src/crewai/tools/tool_usage.py:1089` |
| Tool attribution — Started/Finished events | `use`/`_use`/`ause`/`_ause` emit `ToolUsageStartedEvent` with `agent_key`, `agent_role`, `agent_id`, `task_id`, `task_name`, `tool_name`, `tool_args`, `tool_class`, plus `agent.fingerprint`; finished event adds `started_at`, `finished_at`, `from_cache`, `output`, `failure` | `lib/crewai/src/crewai/tools/tool_usage.py:273`, `lib/crewai/src/crewai/tools/tool_usage.py:1038` |
| Tool attribution — Failure model | `ToolFailure`/`ToolFailureRecord`/`ToolFailurePolicy` (IGNORE/WARN/RAISE), `detect_tool_failure`, `handle_tool_failure`, `reportable_failure`, `ToolFailureDetectedEvent`; `TaskOutput.tool_failures` surfaced via `CrewOutput.tool_failures` | `lib/crewai/src/crewai/tools/tool_failure.py:35`, `lib/crewai/src/crewai/tools/tool_failure.py:111`, `lib/crewai/src/crewai/tools/tool_usage.py:398` |
| Tool attribution — Input/Selection errors | `_emit_validate_input_error` → `ToolValidateInputErrorEvent`; `_select_tool` → `ToolSelectionErrorEvent` with empty-name handling | `lib/crewai/src/crewai/tools/tool_usage.py:987`, `lib/crewai/src/crewai/tools/tool_usage.py:810` |
| Tool attribution — Tool event types | `ToolUsageEvent.from_task`/`from_agent` hydrates `task_id`/`agent_id`; fingerprint stamped on init | `lib/crewai/src/crewai/events/types/tool_usage_events.py:32` |
| Human decision logs — Task-level HITL | `Task.human_input: bool | None` field; `Crew._train` mode delegates to `core.providers.human_input` | `lib/crewai/src/crewai/task.py:233`, `lib/crewai/src/crewai/core/providers/human_input.py:26` |
| Human decision logs — HumanInputProvider loop | `ExecutorContext.ask_for_human_input` flag; `SyncHumanInputProvider.handle_feedback` / `handle_feedback_async` loops reading stdin until empty string | `lib/crewai/src/crewai/core/providers/human_input.py:147`, `lib/crewai/src/crewai/core/providers/human_input.py:256` |
| Human decision logs — Flow @human_feedback | `@human_feedback(message, emit, llm, provider)` decorates Flow methods; `HumanFeedbackResult(output, feedback, outcome, timestamp, method_name)` | `lib/crewai/src/crewai/flow/human_feedback.py:118`, `lib/crewai/src/crewai/flow/human_feedback.py:362` |
| Human decision logs — Flow events | `HumanFeedbackRequestedEvent` / `HumanFeedbackReceivedEvent` with `flow_name`, `method_name`, `output`, `feedback`, `outcome`, `request_id`; also `MethodExecutionPausedEvent`/`FlowPausedEvent` | `lib/crewai/src/crewai/events/types/flow_events.py:244`, `lib/crewai/src/crewai/events/types/flow_events.py:268` |
| Human decision logs — Hook-level approval gate | `ToolCallHookContext.request_human_input(prompt, default_message)` blocking stdin gate; mirrored `LLMCallHookContext.request_human_input` | `lib/crewai/src/crewai/hooks/tool_hooks.py:86`, `lib/crewai/src/crewai/hooks/llm_hooks.py:114` |
| Human decision logs — Persistence | Flow runtime stores `human_feedback_history: list[HumanFeedbackResult]` + `last_human_feedback`; training handler persists `initial_output`/`human_feedback`/`improved_output` per agent/iteration | `lib/crewai/src/crewai/flow/runtime/__init__.py:604`, `lib/crewai/src/crewai/utilities/evaluators/task_evaluator.py:137`, `lib/crewai/src/crewai/utilities/training_handler.py:19` |
| Memory provenance | `MemoryRecord.source` used for provenance tracking + privacy filtering (`private`); memory events typed `source_type="memory"` | `lib/crewai/src/crewai/memory/types.py:60`, `lib/crewai/src/crewai/memory/unified_memory.py:454` |
| Accountability docs — Absence | No `ACCOUNTABILITY.md` / policy file; docs discuss agent "responsibilities" only in task-routing sense; README/license carry no organizational accountability mapping | `docs/edge/en/concepts/tasks.mdx:53`, `README.md:1` |
| Tests — Responsibility traces | Tests verify fingerprint propagation, tool error attribution, human feedback decorator, skill attribution | `lib/crewai/tests/tools/test_tool_usage_error_attribution.py:5`, `lib/crewai/tests/test_human_feedback_decorator.py:1`, `lib/crewai/tests/test_event_record.py:1` |

## Answers to Dimension Questions

**1. Who is responsible for agent actions?**

Responsibility is bound at the code level, not declared as policy. The Task→Agent edge (`lib/crewai/src/crewai/task.py:163`) is mandatory ("Each task must have ... an agent responsible" at `lib/crewai/src/crewai/task.py:123`). The owning Agent/Crew/Task each carry a `SecurityConfig.fingerprint` (`lib/crewai/src/crewai/security/security_config.py:39`, `lib/crewai/src/crewai/security/fingerprint.py:42`) whose `uuid_str` stamps every emitted event as `source_fingerprint`/`source_type` (`lib/crewai/src/crewai/events/base_events.py:71`). A per-run `crewai_execution_uuid` (`lib/crewai/src/crewai/execution.py:22`) groups nested kickoffs. Trace and telemetry listeners aggregate these into per-agent/task spans (`lib/crewai/src/crewai/events/listeners/tracing/trace_listener.py:408`). There is no single "organization/user/operator" accountability owner documented; the model implicitly says "the Agent that was assigned the Task" plus the fingerprint-holder is responsible.

**2. Is model output attributed?**

Partially, at the artifact level, but without disclaimer or model provenance. `TaskOutput.agent` records the executor's role string (`lib/crewai/src/crewai/tasks/task_output.py:43`), `TaskOutput.messages` retains the LLM transcript, and `CrewOutput.tasks_output` preserves ordered per-task lineage (`lib/crewai/src/crewai/crews/crew_output.py:24`). Events `AgentExecutionStartedEvent`/`CompletedEvent` and `TaskCompletedEvent` echo the same via `source_fingerprint` (`lib/crewai/src/crewai/events/types/agent_events.py:30`, `lib/crewai/src/crewai/events/types/task_events.py:96`). No watermark, disclaimer, or "model attribution" header is injected into `raw`/`json_dict` (`lib/crewai/src/crewai/tasks/task_output.py:36`). A downstream consumer cannot distinguish model-generated vs. human-edited text without external event inspection.

**3. Are tool decisions attributed?**

Yes — this is the most mature attribution path. Every call passes through `ToolUsage` (`lib/crewai/src/crewai/tools/tool_usage.py:84`) which (a) injects fingerprint context via `_build_fingerprint_config` (`lib/crewai/src/crewai/tools/tool_usage.py:1089`), (b) emits `ToolUsageStartedEvent` with `agent_key`/`agent_role`/`agent_id`/`task_id`/`task_name`/`tool_name`/`tool_args`/`tool_class` (`lib/crewai/src/crewai/tools/tool_usage.py:273`), (c) emits `ToolUsageFinishedEvent` with timing, cache hit, output and `failure` (`lib/crewai/src/crewai/tools/tool_usage.py:1038`), and (d) emits `ToolSelectionErrorEvent`/`ToolValidateInputErrorEvent`/`ToolUsageErrorEvent`/`ToolFailureDetectedEvent` for error branches. The structured `ToolFailure`/`ToolFailureRecord` model (`lib/crewai/src/crewai/tools/tool_failure.py:71`) keeps policy-aware attribution on `TaskOutput.tool_failures` and `CrewOutput.tool_failures`.

**4. Are human approvals recorded?**

Yes, on two tracks, with event-grade logging. Track A (Crew/Task): `Task.human_input` (`lib/crewai/src/crewai/task.py:233`) drives `core.providers.human_input.SyncHumanInputProvider` which loops `ask_for_human_input` until the operator hits Enter empty (`lib/crewai/src/crewai/core/providers/human_input.py:256`) and records `initial_output`/`human_feedback`/`improved_output` in the training file (`lib/crewai/src/crewai/utilities/training_handler.py:19`, `lib/crewai/src/crewai/utilities/evaluators/task_evaluator.py:137`). Track B (Flow): `@human_feedback` (`lib/crewai/src/crewai/flow/human_feedback.py:362`) produces `HumanFeedbackResult` and emits `HumanFeedbackRequestedEvent` + `HumanFeedbackReceivedEvent` with `feedback`/`outcome`/`request_id` (`lib/crewai/src/crewai/events/types/flow_events.py:244`), persisted as `human_feedback_history` / `last_human_feedback` (`lib/crewai/src/crewai/flow/runtime/__init__.py:604`). Hook-level `request_human_input` (`lib/crewai/src/crewai/hooks/tool_hooks.py:86`) also pauses for approval but its textual response is not separately journaled beyond the resulting event.

**5. Is accountability documented?**

No. No policy, charter, or accountability matrix exists in the inspected source. The security module declares `Fingerprint` is for "tracking, auditing, and security" (`lib/crewai/src/crewai/security/fingerprint.py:4`) but `SecurityConfig` marks auth/scoping/impersonation as `TODO` (`lib/crewai/src/crewai/security/security_config.py:25`). Docs describe agent "responsibilities" as functional routing (`docs/edge/en/concepts/tasks.mdx:53`), not liability. No disclaimer template, data-provenance statement, or operator/organization responsibility assignment is shipped; answering "who is organizationally accountable?" requires inferring from deployment choices.

## Architectural Decisions

*  **Fingerprint-as-identity** — `SecurityConfig.fingerprint` auto-mints a deterministic or random UUID (`lib/crewai/src/crewai/security/fingerprint.py:54`) and is attached to every Agent/Task/Crew; all events stamp `source_fingerprint`/`source_type` so telemetry and traces are joinable without an external identity service. Tradeoff: decentralized, no PKI or attestation.
*  **Execution-scoped ContextVar correlation** — `crewai_execution_uuid` (`lib/crewai/src/crewai/execution.py:22`) binds outermost run; nested flows/agents inherit via contextvars rather than threading through arguments. Enables inter-process (Celery `kickoff_id`) injection via `set_execution_uuid` without rewriting call sites.
*  **Event bus as audit log** — Directed `EventRecord` graph (`lib/crewai/src/crewai/state/event_record.py:99`) and typed events (`lib/crewai/src/crewai/events/types/tool_usage_events.py:11`) serve as the single attribution/audit surface consumed by tracing, tracing UI, and telemetry; persistence is pluggable via `EventListener`/checkpoint providers (`lib/crewai/src/crewai/state/provider/json_provider.py:16`).
*  **Structured ToolFailure instead of error-string parsing** — `ToolFailure` (`lib/crewai/src/crewai/tools/tool_failure.py:71`) makes "tool ran but reports failure" (HTTP 200 with `ok:false`, MCP `isError`) explicit and policy-driven (IGNORE/WARN/RAISE), preventing silent success masking.
*  **Hook-gated human approval** — `PRE_TOOL_CALL`/`POST_TOOL_CALL` hooks plus `ToolCallHookContext.request_human_input` (`lib/crewai/src/crewai/hooks/tool_hooks.py:86`) allow approval without modifying every tool; mirrors `Flow.@human_feedback` for cross-cutting HITL.

## Notable Patterns

*  **Decorator-lift pattern for HITL** — `@human_feedback` merely stamps `__human_feedback_config__` (`lib/crewai/src/crewai/flow/human_feedback.py:4`), lifted into `FlowMethodDefinition.human_feedback` by the definition builder and consumed by the engine (`lib/crewai/src/crewai/flow/runtime/__init__.py:2981`). Visualization and validation reuse the same metadata.
*  **Dual path for sync/async HITL** — `SyncHumanInputProvider` implements both `handle_feedback` (blocking `input()`) and `handle_feedback_async` (`_async_readline` with `connect_read_pipe` fallback) (`lib/crewai/src/crewai/core/providers/human_input.py:147`), keeping Flow/agent HITL compatible with async kickoff.
*  **Config-channel fingerprint injection** — `ToolUsage._build_fingerprint_config` packs provenance into `config["security_context"]` (`lib/crewai/src/crewai/tools/tool_usage.py:1098`) so LangChain-structured tools with strict schemas do not reject attribution fields.
*  **Memory provenance via `source`** — `MemoryRecord.source` (`lib/crewai/src/crewai/memory/types.py:60`) threads user/session origin through recall and respects `private` isolation, a rare cross-cut between memory and accountability concerns.

## Tradeoffs

*  **Rich event attribution vs. no durable audit store by default** — The event graph can answer "who" if a listener persists it, but OSS defaults to in-memory `EventRecord` and optional `TraceCollectionListener`/SQLite/JSON checkpoint; absent a collector, post-hoc accountability is lost. Durability depends on deployment.
*  **Fingerprint ubiquity vs. weak trust** — UUIDs are mintable by any caller (`Fingerprint.generate(seed)`) with no signing or chain-of-custody; impersonation is trivial and `SecurityConfig` TODOs (`lib/crewai/src/crewai/security/security_config.py:25`) flag this as unfinished. Good for correlation, not for non-repudiation.
*  **Human feedback as pause vs. as record** — High-fidelity event capture (`HumanFeedbackResult.timestamp`, `outcome`) coexists with terminal `input()` prompts (`lib/crewai/src/crewai/core/providers/human_input.py:318`) that are uncaptured unless the event listener is active; headless/automated providers must be supplied to avoid silent loss.
*  **Tool failure transparency vs. backward compatibility** — `ToolFailurePolicy.WARN` (default) preserves agent narration while still recording failures; `IGNORE` fully hides them (necessary for pre-1.16 behavior) at the cost of accountability blind spots.
*  **Output attribution limited to role string** — `TaskOutput.agent` is `agent.role` (`lib/crewai/src/crewai/tasks/task_output.py:43`), not `agent.id` or fingerprint UUID; renaming roles collides attribution and a role-less LiteAgent path drops identity further (`lib/crewai/src/crewai/tools/tool_usage.py:284`).

## Failure Modes / Edge Cases

*  **Silent attribution loss on parse failure** — Tool-call parsing failures are emitted with `tool_name=""` intentionally to avoid metric-cardinality explosion (`lib/crewai/src/crewai/tools/tool_usage.py:925`, `lib/crewai/tests/tools/test_tool_usage_error_attribution.py:5`); useful for safety but means the "intended" tool cannot be audited.
*  **Shared Agent concurrency bleeds failures** — `ToolFailure._active_failures` ContextVar plus agent-field fallback (`lib/crewai/src/crewai/tools/tool_failure.py:260`) can leak or overwrite records when an Agent is reused concurrently; `collect_tool_failures` tolerates missing attributes rather than failing hard.
*  **Checkpoint restore can downgrade event type** — `TaskFailedEvent.error_type` serialization/deserialization via `module:qualname` (`lib/crewai/src/crewai/events/types/task_events.py:10`, `lib/crewai/src/crewai/events/types/task_events.py:160`) drops to `None` if module unloaded, silently losing error taxonomy on restore.
*  **No output disclaimer → downstream misattribution** — Raw LLM text is emitted as the system's answer (`TaskOutput.raw`, `CrewOutput.raw`) with no attached "generated by model" marker; embedding in user-facing products invites mistaken human authorship. Spec explicitly sought this; not implemented.
*  **HITL bypass via provider substitution** — Supplying a custom `HumanInputProvider`/`HumanFeedbackProvider` can auto-approve or fabricate feedback with no integrity check (`lib/crewai/src/crewai/core/providers/human_input.py:59`), and the event still records it as "human" — accountability relies on honest provider.
*  **Memory source spoofable** — `source` on `MemoryRecord` is free-form string (`lib/crewai/src/crewai/memory/types.py:60`) and `remember_many(lessons, source=learn_source)` (`lib/crewai/src/crewai/flow/human_feedback.py:346`) uses a fixed `hitl` source by default; an attacker-controlled recall query can poison attribution.

## Future Considerations

*  **Add a formal Accountability Policy doc** — Map technical identities (fingerprint, execution_uuid, task/agent ids) to organizational roles (model, runtime, tool author, user, operator, organization) and state who is liable when a tool mutates external state; codify retention/immutability expectations for the event log.
*  **Inject model-output provenance header** — Stamp `TaskOutput`/`CrewOutput` with `generated_by: model_id`, `fingerprint`, and optional watermark/disclaimer string (opt-in per Crew `output_disclaimer` flag) so rendered answers remain attributable outside the event bus.
*  **Harden fingerprint to verifiable identity** — Move from bare UUID to signed attestation or at least seed-derived deterministic UUIDs tied to authenticated principals; complete the TODOs in `security_config.py` (auth, scoping, delegation tokens).
*  **Make human decision logs append-only** — Persist `HumanFeedbackRequestedEvent`/`ReceivedEvent` + `ToolUsageFinishedEvent` to the configured checkpoint provider by default, with integrity hashing, so audit survives process death; add CLI `crewai traces show --human-feedback`.
*  **Use `agent.id` in output attribution** — Store `TaskOutput.agent_id` / `agent_fingerprint` alongside `agent` role string to prevent alias collisions and to support idempotency checks in regulated contexts.

## Questions / Gaps

*  **No evidence of policy attribution document** — Greps across `lib/crewai/src`, `docs/edge`, and `README.md` found no accountability charter; if one lives outside this source tree it was not discoverable. Marked as `No clear evidence found`.
*  **No output disclaimer mechanism** — No template, config key, or test asserts a disclaimer is appended to model output; search for `disclaimer`/`watermark` in code returned only documentation chatter.
*  **Unclear retention for human feedback** — Events are emitted but retention is listener-dependent; no default TTL or GDPR-style purge policy for `human_feedback_history` or `MemoryRecord.source` is codified in code under study.
*  **Operator vs. organization boundary** — `execution_uuid` hints at Enterprise/Celery session ownership but the OSS code never maps that uuid to a legal entity; responsibility for actions taken under that uuid remains ambiguous.

---

Generated by `23.03-responsibility-and-accountability-model` against `crewai`.
