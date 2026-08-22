# Source Analysis: crewai

## Dimension 10.02: Event Schema and Lifecycle Events

### Source Info

| Field | Value |
|-------|-------|
| Name | crewai |
| Path | `studies/agent-harness-study/sources/crewai` |
| Language / Stack | Python (Pydantic v2 models, asyncio, contextvars) |
| Analyzed | 2026-08-22 |

All citations below are workspace-relative paths under `studies/agent-harness-study/sources/crewai/`.

## Summary

CrewAI implements a first-class, strongly-typed event system. Every event is a Pydantic model deriving from a common `BaseEvent` (`studies/agent-harness-study/sources/crewai/lib/crewai/src/crewai/events/base_events.py:66-87`) that carries a UTC timestamp, a literal string discriminator (`type`), entity identity fields (`task_id`, `agent_id`, fingerprints), and five linkage fields: `event_id`, `parent_event_id`, `previous_event_id`, `triggered_by_event_id`, `started_event_id`, plus an integer `emission_sequence`. A singleton bus (`studies/agent-harness-study/sources/crewai/lib/crewai/src/crewai/events/event_bus.py:94-952`) stamps linkage metadata on every synchronous emission via `_prepare_event` (`studies/agent-harness-study/sources/crewai/lib/crewai/src/crewai/events/event_bus.py:536-568`) using contextvar-based scope stacks, then dispatches to sync/async handlers with dependency-ordered execution plans.

Roughly 90 concrete event classes are organized by domain module (crew, agent, task, LLM, tool usage, memory, knowledge, MCP, A2A, checkpoint, flow, reasoning, signals) under `studies/agent-harness-study/sources/crewai/lib/crewai/src/crewai/events/types/`, each pinned with `Literal[...]` type values. Lifecycle pairing is explicit: `SCOPE_STARTING_EVENTS`, `SCOPE_ENDING_EVENTS`, and a `VALID_EVENT_PAIRS` table map ~45 end events to their required start events (`studies/agent-harness-study/sources/crewai/lib/crewai/src/crewai/events/event_context.py:240-369`). Events are recorded into an in-memory directed graph (`EventRecord`, `studies/agent-harness-study/sources/crewai/lib/crewai/src/crewai/state/event_record.py:99-146`) that is serialized into checkpoints (`studies/agent-harness-study/sources/crewai/lib/crewai/src/crewai/state/runtime.py:190-214`) and replayed on resume without mutating recorded ids (`replay()`, `studies/agent-harness-study/sources/crewai/lib/crewai/src/crewai/events/event_bus.py:671-730`). The lifecycle of a run is largely reconstructable from events alone — parent/child, causal, sequential, and start/completion edges are all first-class — but there is no per-event schema version field, cancellation events exist only for A2A server tasks, sequence numbers are per-context rather than globally unique, and the async emit path skips scope wiring.

## Rating

**8 / 10**

Typed schemas (Pydantic + Literal discriminators), ordering (`emission_sequence` + `previous_event_id` chain), rich parent/causal linkage, comprehensive started/completed/failed pairs enforced by a tested pairing table, durability via checkpoint-persisted `EventRecord`, and an extensive test suite covering ordering, replay, pairing invariants, and thread safety. Deductions from 9–10: no event schema/version field (only package-level versions on trace batches and checkpoint envelopes), no cancellation lifecycle for core entities (crew/task/agent), per-context emission counters can collide across concurrent threads, and `aemit()` does not populate parent/previous linkage.

## Evidence Collected

| Area | Evidence | File:Line |
|------|----------|-----------|
| Base event schema | `BaseEvent` defines `timestamp` (UTC default), `type`, `source_fingerprint`/`source_type`/`fingerprint_metadata`, `task_id`/`task_name`/`agent_id`/`agent_role`, `event_id` (uuid4), `parent_event_id`, `previous_event_id`, `triggered_by_event_id`, `started_event_id`, `emission_sequence` | `studies/agent-harness-study/sources/crewai/lib/crewai/src/crewai/events/base_events.py:66-87` |
| Typed discriminators | Every event pins `type: Literal["..."]`, e.g. `CrewKickoffStartedEvent.type = "crew_kickoff_started"` | `studies/agent-harness-study/sources/crewai/lib/crewai/src/crewai/events/types/crew_events.py:40` |
| Event catalog | Aggregate union of ~90 typed events across domains; lazy public exports via mapping | `studies/agent-harness-study/sources/crewai/lib/crewai/src/crewai/events/event_types.py:123-216`; `studies/agent-harness-study/sources/crewai/lib/crewai/src/crewai/events/__init__.py:154-252` |
| Sequence numbers | Per-context `itertools.count` counter behind a ContextVar; `get_next_emission_sequence()` increments and records last value; reset/resume helpers for tests and checkpoint resume | `studies/agent-harness-study/sources/crewai/lib/crewai/src/crewai/events/base_events.py:13-63` |
| Sequence stamping | `_prepare_event` assigns `emission_sequence`, `previous_event_id`, `triggered_by_event_id`, resolves `parent_event_id` from scope stack, pairs ending events with popped starts, and records into `EventRecord` | `studies/agent-harness-study/sources/crewai/lib/crewai/src/crewai/events/event_bus.py:536-568` |
| Timestamps | `timestamp: datetime = Field(default_factory=lambda: datetime.now(timezone.utc))`; tool finish events add `started_at`/`finished_at` wall-clock pair | `studies/agent-harness-study/sources/crewai/lib/crewai/src/crewai/events/base_events.py:69`; `studies/agent-harness-study/sources/crewai/lib/crewai/src/crewai/events/types/tool_usage_events.py:65-66` |
| Ordering safeguard for streams | `LLMStreamChunkEvent` is dispatched synchronously (bypasses thread pool) "to preserve ordering" | `studies/agent-harness-study/sources/crewai/lib/crewai/src/crewai/events/event_bus.py:509-510,626-628` |
| Scope stack machinery | Contextvar stack of `(event_id, type)` pairs; push/pop with depth limit (default 100), mismatch/empty-pop behaviors WARN/RAISE/SILENT | `studies/agent-harness-study/sources/crewai/lib/crewai/src/crewai/events/event_context.py:20-38,171-191` |
| Lifecycle pairing tables | `SCOPE_STARTING_EVENTS` (26 types), `SCOPE_ENDING_EVENTS` (~48 types), `VALID_EVENT_PAIRS` maps every ending event to its required starting event | `studies/agent-harness-study/sources/crewai/lib/crewai/src/crewai/events/event_context.py:240-369` |
| Causal tracking | `triggered_by_scope()` context manager sets `triggered_by_event_id` for listener-style causality | `studies/agent-harness-study/sources/crewai/lib/crewai/src/crewai/events/event_context.py:118-133` |
| Run-level linkage | Crew kickoff stores `crew._kickoff_event_id = started_event.event_id`; completion events pass `started_event_id=self._kickoff_event_id`; scope stack re-seeded from it on nested runs (`crew.py:561-562`) | `studies/agent-harness-study/sources/crewai/lib/crewai/src/crewai/crews/utils.py:320-322`; `studies/agent-harness-study/sources/crewai/lib/crewai/src/crewai/crew.py:1060,1272,1877`; `studies/agent-harness-study/sources/crewai/lib/crewai/src/crewai/agent/core.py:1619-1629` |
| Event persistence graph | `EventRecord.add()` wires typed edges (`parent`/`child`, `trigger`/`triggered_by`, `next`/`previous`, `started`/`completed_by`) keyed by `event_id`; O(1) lookup, traversal helpers `descendants()`/`roots()` | `studies/agent-harness-study/sources/crewai/lib/crewai/src/crewai/state/event_record.py:52-61,110-146,160-202` |
| Durable reconstruction | `RuntimeState._serialize` includes `event_record` and `crewai_version`; restore validates record back into state; version-gated migration for older checkpoint formats | `studies/agent-harness-study/sources/crewai/lib/crewai/src/crewai/state/runtime.py:89-119,190-214` |
| Replay semantics | `replay()` dispatches recorded events without `_prepare_event` (preserves ids/sequence); `is_replaying()` flag lets side-effectful listeners opt out; resume replays completed-method events | `studies/agent-harness-study/sources/crewai/lib/crewai/src/crewai/events/event_bus.py:671-730` |
| Checkpoint resume scoping | `resume_task_scope()` finds latest recorded `task_started` for a task id (by max `emission_sequence`) and pushes its scope so post-resume events nest correctly; `restore_event_scope()` restores stack from checkpoint | `studies/agent-harness-study/sources/crewai/lib/crewai/src/crewai/events/event_context.py:136-168` |
| Entity context on events | LLM events carry `call_id`, `response_id`, `finish_reason`, sampling params; tool events carry `tool_name`, `tool_args`, `run_attempts`, `plan_step_number`; task/agent ids auto-filled from live objects then nulled | `studies/agent-harness-study/sources/crewai/lib/crewai/src/crewai/events/types/llm_events.py:13,98-99,136-151`; `studies/agent-harness-study/sources/crewai/lib/crewai/src/crewai/events/types/tool_usage_events.py:16-27`; `studies/agent-harness-study/sources/crewai/lib/crewai/src/crewai/events/base_events.py:101-116` |
| Source fingerprinting | Domain base events copy `source_fingerprint`/`source_type`/metadata from crew/agent/task fingerprints at construction | `studies/agent-harness-study/sources/crewai/lib/crewai/src/crewai/events/types/crew_events.py:22-27`; `studies/agent-harness-study/sources/crewai/lib/crewai/src/crewai/events/types/task_events.py:7-21`; `studies/agent-harness-study/sources/crewai/lib/crewai/src/crewai/events/types/agent_events.py:123-130` |
| Signal/interruption events | `SigTermEvent`/`SigIntEvent`/`SigHupEvent`/`SigTStpEvent`/`SigContEvent` with discriminated union adapter and `@on_signal` registration helper; serialization roundtrip tested | `studies/agent-harness-study/sources/crewai/lib/crewai/src/crewai/events/types/system_events.py:27-102`; `studies/agent-harness-study/sources/crewai/lib/crewai/tests/events/types/test_system_events.py:179-196` |
| Cancellation coverage gap | Only `A2AServerTaskCanceledEvent(type="a2a_server_task_canceled")` exists; grep for `cancel` across `events/types/` finds no crew/task/agent cancellation events | `studies/agent-harness-study/sources/crewai/lib/crewai/src/crewai/events/types/a2a_events.py:599-608` |
| Versioning gap | No `version` field on `BaseEvent` or any event class; only `crewai_version` attached to exported trace batches and checkpoint envelopes, plus domain-level `protocol_version` on A2A events | `studies/agent-harness-study/sources/crewai/lib/crewai/src/crewai/events/listeners/tracing/trace_batch_manager.py:38-49`; `studies/agent-harness-study/sources/crewai/lib/crewai/src/crewai/state/runtime.py:193`; `studies/agent-harness-study/sources/crewai/lib/crewai/src/crewai/events/types/a2a_events.py:85` |
| Telemetry export preserves linkage | Trace listener copies `event_id`, `emission_sequence`, `parent_event_id`, `previous_event_id`, `triggered_by_event_id` into `TraceEvent` payloads sorted by sequence | `studies/agent-harness-study/sources/crewai/lib/crewai/src/crewai/events/listeners/tracing/trace_listener.py:928-932`; `studies/agent-harness-study/sources/crewai/lib/crewai/src/crewai/events/listeners/tracing/types.py:8-21`; `studies/agent-harness-study/sources/crewai/lib/crewai/src/crewai/events/listeners/tracing/trace_batch_manager.py:344-345` |
| Async emit asymmetry | `aemit()` registers source, sets `emission_sequence`, records event — but never calls `_prepare_event`, so `parent_event_id`/`previous_event_id`/scope pairing are not wired on this path | `studies/agent-harness-study/sources/crewai/lib/crewai/src/crewai/events/event_bus.py:769-780` |
| Ordering tests | Completed > started sequences; task parent == crew event id; previous-event chains monotonic in sequence; parallel listeners share trigger id | `studies/agent-harness-study/sources/crewai/lib/crewai/tests/events/test_event_ordering.py:146-170,173-195,521-580,778-842` |
| Pairing invariant tests | Tests assert every ending event has a pair, every pair references a starting event, and the two sets are disjoint | `studies/agent-harness-study/sources/crewai/lib/crewai/tests/events/test_event_context.py:97-105` |
| Replay tests | Replay preserves `event_id`/`parent_event_id`/`emission_sequence`; `is_replaying` true only during replay; checkpoint resume replays completed methods' events | `studies/agent-harness-study/sources/crewai/lib/crewai/tests/events/test_event_replay.py:29-51,54-88,113-174` |
| Concurrency tests | Concurrent emit, concurrent registration, rapid-emit stress, shutdown-during-emit suites over RWLock-guarded handler tables | `studies/agent-harness-study/sources/crewai/lib/crewai/tests/utilities/events/test_thread_safety.py:18-145`; `studies/agent-harness-study/sources/crewai/lib/crewai/tests/utilities/events/test_shutdown.py:20-241` |

## Answers to Dimension Questions

**1. Are events typed and versioned?**
Typed: yes, thoroughly. All events are Pydantic subclasses of `BaseEvent` with `Literal` type discriminators validated through a registry built from subclass defaults for deserialization (`_resolve_event`/`_build_event_type_map`, `studies/agent-harness-study/sources/crewai/lib/crewai/src/crewai/state/event_record.py:17-49`). Signal events use a discriminated union adapter (`signal_event_adapter`, `studies/agent-harness-study/sources/crewai/lib/crewai/src/crewai/events/types/system_events.py:70-75`). Versioned: **no evidence found of a per-event schema version field.** I searched `lib/crewai/src/crewai/events/` for `version|schema_version|SCHEMA`: the only hits are the crewai package version stamped onto exported trace batches (`trace_batch_manager.py:41,49`) and OTel payloads, the checkpoint envelope's `crewai_version` used for forward migrations (`runtime.py:89-119`), and A2A's domain-specific `protocol_version` (`a2a_events.py:85`). Consumers cannot branch parsing logic on event schema version; they must rely on package version heuristics.

**2. Are events ordered and timestamped?**
Yes. Each event gets a UTC timestamp by default (`base_events.py:69`) and a monotonically increasing `emission_sequence` assigned synchronously in the emitting thread before dispatch (`event_bus.py:546`), backed by a ContextVar counter (`base_events.py:13-41`) with reset/resume helpers (`base_events.py:49-63`). Ordering guarantees are strengthened operationally: stream chunks bypass the executor pool so chunk order is preserved (`event_bus.py:509-510,626-628`), and `flush()` blocks until pending handler futures complete (`event_bus.py:732-767`). Caveat: the counter is per-context (ContextVar), so sequences from concurrently executing flows in different contexts/threads are not globally comparable — cross-thread total order must be derived from timestamps or the graph edges instead. Tests assert within-run monotonicity (`test_event_ordering.py:568-580`).

**3. Do events carry sufficient context?**
Yes — this is a strength. Five id-linkage fields (`event_id`, `parent_event_id`, `previous_event_id`, `triggered_by_event_id`, `started_event_id`, `base_events.py:82-87`) capture hierarchical containment, linear precedence, causal triggering (e.g., which `method_execution_finished` triggered a listener's start — `test_event_ordering.py:629-681`), and start/end pairing. Entity identity is embedded directly (`task_id`/`task_name`/`agent_id`/`agent_role` auto-extracted from live objects, `base_events.py:101-116`), plus stable fingerprints (`source_fingerprint`/`source_type`, populated at construction, e.g., `crew_events.py:22-27`) and domain detail (`call_id`, `response_id`, `tool_name`, `tool_args`, `run_attempts`, `plan_step_number`). The run anchor is preserved across restarts: `Crew._kickoff_event_id` is checkpointed and restored (`base_agent.py:450-451`, `runtime.py:64,83`), and completion events reference it as `started_event_id` (`crew.py:1060`). The `EventRecord` graph turns these fields into queryable structure (`event_record.py:110-146`).

**4. Are lifecycle events comprehensive?**
Mostly. Started/completed/failed triplets cover flows, methods, crew kickoff/train/test, agents (including lite agents and evaluations), tasks, LLM calls, guardrails, tool usage, MCP connections/tool executions, memory retrieval/save/query, knowledge queries, A2A operations, and agent reasoning — enumerated exhaustively in `VALID_EVENT_PAIRS` (`event_context.py:321-369`) with integrity tests (`test_event_context.py:97-105`). Pause lifecycle exists for flows (`flow_paused`, `method_execution_paused`, `event_context.py:272,275`). Checkpoints have their own full lifecycle (started/completed/failed/fork/restore/pruned, `event_types.py:207-215`). Gaps: (a) cancellation — only `A2AServerTaskCanceledEvent` exists (`a2a_events.py:599-608`); there is no `task_canceled`/`crew_kickoff_canceled`/`agent_execution_canceled`; interruption is only approximated by OS signal events (`system_events.py:27-83`); (b) naming inconsistency — agent termination uses `agent_execution_error` while siblings use `_failed` (`event_context.py:282-283` vs `339-340`).

**Can you reconstruct the full lifecycle of any run from events alone?**
Largely yes: the parent/child + triggered-by + previous/next + started/completed-by edges persisted in the checkpointed `EventRecord` allow rebuilding a run's tree, causal chains, and timeline even after process restart, and resume logic deliberately replays recorded events (`test_event_replay.py:113-174`). Exceptions where reconstruction degrades: cancellations (no events emitted), async-path emissions missing linkage (`aemit()`, `event_bus.py:769-780`), and global interleaving of concurrent contexts where sequence numbers collide.

## Architectural Decisions

1. **Rich immutable-ish base schema with relationship fields baked in.** Rather than ad-hoc payloads, all events inherit identity + linkage columns (`base_events.py:66-87`), making every consumer able to assume the same correlation contract.
2. **Centralized metadata stamping in the bus, not emitters.** `_prepare_event` (`event_bus.py:536-568`) computes sequence/parent/chain fields once, so emitter code stays declarative; `replay()` deliberately skips it to keep recorded history immutable (`event_bus.py:671-686`).
3. **Contextvars for execution-context propagation.** Scope stack, last-event-id, triggering-event-id, and emission counter all ride ContextVars (`event_context.py:41-55`, `base_events.py:13-30`), giving async-safe isolation per flow execution; the bus copies context when hopping to the executor (`event_bus.py:512-515,630-632`).
4. **Declarative pairing contract.** Start/end/pairs frozensets+dict (`event_context.py:240-369`) act as a machine-checked lifecycle grammar, configurable between WARN/RAISE/SILENT (`event_context_config`, `event_context.py:20-27`).
5. **Events as durable data.** The `EventRecord` graph is part of `RuntimeState` serialization, so event history is checkpointed/restored alongside entities (`runtime.py:190-214`), enabling time-travel debugging and mid-run resume with correct nesting (`resume_task_scope`, `event_context.py:141-168`).
6. **Replay with side-effect opt-out.** The `is_replaying()` flag (`event_bus.py:67-80`) separates timeline-reconstructing listeners (traces, console) from side-effectful ones (checkpoint writes), tested explicitly (`test_event_replay.py:91-110`).

## Notable Patterns

- **Literal-typed discriminators + registry-driven polymorphic deserialization**: event class resolved from the stored `type` default when hydrating checkpoints (`event_record.py:20-49`).
- **Scope-stack pairing**: ending events pop the stack and adopt the enclosing parent, with `started_event_id` linking completion back to its start (`event_bus.py:547-564`).
- **Ordered streaming fast path**: special-casing `LLMStreamChunkEvent` for synchronous dispatch to avoid thread-pool reordering (`event_bus.py:509-510,626-628`).
- **Fingerprint enrichment at construction**: domain base classes pull `source_fingerprint` from live objects and null out object references to keep events serializable (`crew_events.py:18-33`, `base_events.py:101-116`).
- **Defensive validators**: exotic provider values coerced to None instead of crashing event construction (`llm_events.py:101-114`).
- **Lazy export surface**: `__getattr__`-based lazy loading keeps ~12 Pydantic modules out of import time (`events/__init__.py:10-11,154-277`).

## Tradeoffs

- **Per-context vs global ordering**: ContextVar counters give cheap isolation for concurrent flows but sacrifice globally unique sequence numbers; consumers wanting a strict interleave across threads need timestamps or graph traversal.
- **Immutability by convention, not enforcement**: `emit()` mutates event instances (`_prepare_event` writes fields); nothing prevents double-emitting the same instance, unlike `replay()` which protects recorded events.
- **Payload weight vs observability**: events embed raw inputs like method `params`/`state` dicts (constructor surface visible at `test_event_replay.py:18-26`) and full LLM message arrays (`llm_events.py:47`) — great for reconstruction, heavy for memory/log volume.
- **Sync path fidelity vs async simplicity**: `emit()` maintains full linkage; `aemit()` trades correctness of lineage for simplicity, creating divergent behavior depending on entry point (`event_bus.py:769-780`).
- **Stringly-typed `type` values with casing drift**: most types are lowercase snake_case, but signal events use uppercase literals (`"SIGTERM"`, `system_events.py:30`) — harmless internally (registry keys off exact strings) but surprising for external consumers filtering on prefixes.

## Failure Modes / Edge Cases

- **Missing ending events**: unbalanced scopes raise/warn via `handle_empty_pop` and depth-limit `StackDepthExceededError` (default cap 100) catches runaway stacks (`event_context.py:24,176-180,194-205`).
- **Mismatched pairs**: popping a start type that doesn't match the expected partner triggers `EventPairingError` or warning depending on config (`event_context.py:208-223`).
- **Handler failures are isolated**: sync errors collected per-handler and printed, async errors surfaced via `gather(return_exceptions=True)`; one bad listener cannot kill emission (`event_bus.py:415-455`).
- **Shutdown races**: emissions during shutdown are ignored with a warning; flush/wait/cancel paths covered by dedicated tests (`event_bus.py:600-605,897-944`; `tests/utilities/events/test_shutdown.py:20-137`).
- **Checkpoint round-trip drift**: unknown event types fall back to base `BaseEvent` validation (`event_record.py:29-35`), silently narrowing schema fidelity after version skew; mitigated only by envelope-level migration warnings (`runtime.py:106-113`).
- **Resume edge**: if no prior `task_started` exists in the restored record, `resume_task_scope` returns False and subsequent events lose task-scoped parenting (`event_context.py:153-166`).
- **Concurrent sequence collision**: two threads emitting simultaneously each hold separate counters; identical `emission_sequence` values can appear in one merged log (by design of `base_events.py:13-25`), which would break naive total-order sorting.

## Future Considerations

- Add a `schema_version` (or reuse the checkpoint envelope pattern) to `BaseEvent` so downstream parsers evolve safely; today the only lever is package-version sniffing.
- Introduce cancellation lifecycle events for crews/tasks/agents (mirroring `A2AServerTaskCanceledEvent`) so interrupted runs are distinguishable from failed ones in the record.
- Wire `aemit()` through the same `_prepare_event`-equivalent scope logic (or document it as linkage-free by contract) to close the async lineage gap.
- Consider a global monotonic allocator (e.g., process-wide lock-protected counter) or a `(context_id, sequence)` tuple to make ordering unambiguous under concurrency.
- Normalize agent failure naming (`agent_execution_error` → `agent_execution_failed` alias) and signal-event casing for consistent tooling.

## Questions / Gaps

- No evidence found of any event-schema versioning strategy at the individual event level; whether the project intends the checkpoint envelope's `crewai_version` (`runtime.py:193`) to serve that role for event payloads specifically is undocumented in code.
- No evidence found of retention/bounding policy for `EventRecord` in long-lived processes (nodes accumulate in `state.event_record`); only `clear()` exists (`event_record.py:214-217`) and I did not find production callers of it inside the studied source.
- Whether OTel GenAI compliance fields on LLM events (`llm_events.py:51-61`) imply a commitment to an external event schema standard could not be confirmed from code alone.

---

Generated by `10.02-event-schema-and-lifecycle-events` against `crewai`.
