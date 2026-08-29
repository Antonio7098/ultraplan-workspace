# Source Analysis: langgraph

## 16.02 Diff, Review, and Rollback

### Source Info

| Field | Value |
|-------|-------|
| Name | langgraph |
| Path | `studies/agent-harness-study/sources/langgraph` |
| Language / Stack | Python (monorepo: `libs/langgraph`, `libs/checkpoint`, `libs/checkpoint-sqlite`, `libs/checkpoint-postgres`, `libs/sdk-py`, `libs/prebuilt`) |
| Analyzed | 2026-08-28 |

## Summary

LangGraph is a stateful graph-execution framework, not an artifact registry. The closest analog to "artifact" in this dimension is a **checkpoint / thread state** (and secondarily a compiled graph definition). The framework provides mature checkpoint persistence and time-travel via `BaseCheckpointSaver` + `Pregel` (`get_state`, `get_state_history`, `update_state`/`bulk_update_state`, replay/fork), with metadata (`source`, `step`, `parents`, `run_id`) and human-in-the-loop interrupts (`interrupt()`/`Command.resume`). It does **not** provide diff generators, artifact review/approval gates, comment/annotation models, or explicit rollback handlers. Reversion is achieved by branching from a historical `checkpoint_id` (fork), preserving an immutable audit chain, but not via an atomic "rollback + tombstone + audit" operation. The LangGraph Platform (deployment versioning, artifact store) is out of scope for this isolated library source.

## Rating

**4 / 10 — Present but inconsistent, weakly documented, fragile for the stated dimension**

Rationale: Core primitives for versioning and traceability exist (`BaseCheckpointSaver` with monotonically-increasing `id`, `get_state_history`, `parents`/`run_id` metadata, and time-travel via `checkpoint_id`/`update_state` fork) and are well-tested (`tests/test_time_travel.py`, `tests/test_delta_channel_update_state.py`). However the dimension's specific capabilities are absent or indirect: no diff/comparison API, no artifact review/approval workflow (only execution-level `interrupt` HITL), no `rollback` verb (only branch-via-fork, leaving the "bad" checkpoint reachable), no comment/annotation storage, and `run_id` traceability exists in the type but its population path in OSS is config-dependent and not enforced. Pruning/copy-for-runs carries explicit `DeltaChannel` silent-corruption warnings, underscoring that history is fragile under lifecycle operations.

## Evidence Collected

Every entry MUST include a file path with line numbers. Format: `path/to/file.ts:NN`.

| Area | Evidence | File:Line |
|------|----------|-----------|
| Checkpoint model | `Checkpoint` TypedDict with `v`, `id`, `ts`, `channel_values`, `channel_versions`, `versions_seen`, `updated_channels`; `id` is monotonic UUIDv6 usable for sorting | `libs/checkpoint/langgraph/checkpoint/base/__init__.py:92-123` |
| CheckpointMetadata | `source: Literal["input","loop","update","fork"]`, `step`, `parents: dict[str,str]`, `run_id`, `counters_since_delta_snapshot` | `libs/checkpoint/langgraph/checkpoint/base/__init__.py:38-86` |
| Base saver interface | `get_tuple`, `list(filter, before, limit)`, `put`, `put_writes`, `delete_thread`, `delete_for_runs`, `copy_thread`, `prune` (+ async variants), `get_delta_channel_history` | `libs/checkpoint/langgraph/checkpoint/base/__init__.py:227-650` |
| DeltaChannel history override | `InMemorySaver.get_delta_channel_history` walks parent chain, per-channel `writes` + optional `seed` (`_DeltaSnapshot`) | `libs/checkpoint/langgraph/checkpoint/memory/__init__.py:142-229` |
| Version generator | `get_next_version(current) -> str` monotonic per-channel; InMemorySaver format `{seq:032}.{rand:016}` | `libs/checkpoint/langgraph/checkpoint/base/__init__.py:692-711` , `libs/checkpoint/langgraph/checkpoint/memory/__init__.py:619-628` |
| StateSnapshot + history API | `StateSnapshot(values, next, config, metadata, created_at, parent_config, tasks, interrupts)`; `Pregel.get_state` / `get_state_history(filter, before, limit)` delegates to `checkpointer.list`; async variants mirror | `libs/langgraph/langgraph/types.py:643-661` , `libs/langgraph/langgraph/pregel/main.py:1391-1432`, `libs/langgraph/langgraph/pregel/main.py:1479-1530` |
| Update / fork / replay (rollback analogue) | `bulk_update_state(supersteps: Sequence[Sequence[StateUpdate]])`, `update_state(values, as_node, task_id)` → `bulk_update_state`; handles `as_node=="__copy__"` fork via parent chain, new `source="fork"` checkpoint, DeltaChannel stub at `step=-1`; async `abulk_update_state`/`aupdate_state` | `libs/langgraph/langgraph/pregel/main.py:1589-1660`, `libs/langgraph/langgraph/pregel/main.py:1798-1877`, `libs/langgraph/langgraph/pregel/main.py:2005-2060`, `libs/langgraph/langgraph/pregel/main.py:2530-2556` |
| Prepare snapshot / apply writes | `_prepare_state_snapshot` builds `StateSnapshot` from `CheckpointTuple`, computes `next_tasks` via `prepare_next_tasks`, optional pending-write application | `libs/langgraph/langgraph/pregel/main.py:1144-1265` |
| Protocol contract | `PregelProtocol` abstract `get_state_history`, `update_state`, `bulk_update_state` confirms these are public framework surface | `libs/langgraph/langgraph/pregel/protocol.py:57-105` |
| HITL / review-adjacent | `interrupt(value) -> Any` raises `GraphInterrupt((Interrupt(value, id=...),))` requiring enabled checkpointer; resume via `Command(resume=...)`; `StateSnapshot.interrupts` + `PregelTask.interrupts/state` surface interrupts | `libs/langgraph/langgraph/types.py:811-934`, `libs/langgraph/langgraph/types.py:597-605`, `libs/langgraph/langgraph/types.py:643-661` |
| Graph interrupt config | `interrupt_before`/`interrupt_after` on `Pregel` / `StateGraph.compile(interrupt_before, interrupt_after)` | `libs/langgraph/langgraph/graph/state.py:1116-1252` , `libs/langgraph/langgraph/pregel/main.py:767-769` |
| Multitask "rollback" (run-scoped, not artifact) | `MultitaskStrategy = Literal["reject","interrupt","rollback","enqueue"]` — `rollback` = cancel current task and start new one; `CancelAction = Literal["interrupt","rollback"]` — `rollback` = delete run + associated checkpoints | `libs/sdk-py/langgraph_sdk/schema.py:81-86`, `libs/sdk-py/langgraph_sdk/schema.py:137-141` |
| Wrapped metadata | `get_checkpoint_metadata` / `get_serializable_checkpoint_metadata` merges `config.metadata` + `config.configurable` (propagates `run_id`, etc.) into persisted metadata | `libs/checkpoint/langgraph/checkpoint/base/__init__.py:752-785` |
| Time-travel tests | Replay/fork suites: `test_replay_*`, `test_fork_*`, `test_multiple_forks`, `test_sequential_interrupts_*`; asserts `metadata.source` values `input/loop/fork`, `parents`, fork creates new branch while original preserved | `libs/langgraph/tests/test_time_travel.py:69-677` , `libs/langgraph/tests/test_time_travel_async.py:110-590` |
| Delta update_state tests | Fresh-thread stub, consecutive updates, bulk multi-task per superstep, history chain ordering | `libs/langgraph/tests/test_delta_channel_update_state.py:1-280` |
| Lifecycle warnings | `delete_for_runs`, `copy_thread`, `prune` each carry `!!! warning "DeltaChannel"` — deleting ancestor writes/snapshots silently breaks `DeltaChannel` reconstruction (returns empty seed, no error) | `libs/checkpoint/langgraph/checkpoint/base/__init__.py:331-414` |
| No diff generator | Grep for `diff` across `libs/langgraph`/`libs/checkpoint` yields only CI `git diff --name-only`, benchmark `compare_to`, `Different model` test prose, and `diff_doc` in store tests — no artifact/state diff utility found | Search: `diff\|compare` grep across source (no artifact diff at `libs/langgraph/langgraph/**` or `libs/checkpoint/langgraph/checkpoint/**`) |
| No review/approval workflow | Grep `review|approval` yields only `interrupt(... {"review": human_review})` example and `require_approval: never` tool flag — no staged approval, RBAC, or gate model for artifacts | Search: `review\|approval` grep + `libs/prebuilt/langgraph/prebuilt/interrupt.py:75` |
| No comment/annotation model | Grep `comment|annotation` yields Sphinx formatting rules, SQL `_leading_comment_remover` injection guard, and Pydantic annotation helpers — no checkpoint/thread comment or annotation field | Search: `comment\|annotation` grep; `libs/checkpoint/langgraph/checkpoint/base/__init__.py:38-86` has no comment key |
| Artifact store n/a | `src/core/loop.ts:42` pattern not applicable; closest CLI/package tooling is SDK `threads.update_state` proxy, not local artifact registry | `libs/langgraph/langgraph/pregel/remote.py:588-662` (remote delegating to `sync_client.threads.update_state`) |

## Answers to Dimension Questions

**1. Can artifacts be compared? — No (manual only).**
No diff generator exists in the studied source. `Checkpoint`/`StateSnapshot` are plain TypedDict/NamedTuple; comparison requires caller to fetch two tuples via `get_state`/`get_state_history` and manually diff `values` / `channel_values` / `channel_versions`. `get_delta_channel_history` (`libs/checkpoint/langgraph/checkpoint/base/__init__.py:582-649`) returns per-channel `writes + seed` for `DeltaChannel` reconstruction but is not a user-facing diff API and is marked beta. Grep across `libs/langgraph` and `libs/checkpoint` for `diff` found only CI formatting (`Makefile:22 git diff`) and benchmark `pyperf compare_to`, not artifact diff. Finding is `No clear evidence found` for a deliberate diff feature; boundary searched: `libs/langgraph/langgraph/**`, `libs/checkpoint/**`, `libs/checkpoint-sqlite/**`, tests.

**2. Is there a review workflow? — Ad-hoc via HITL interrupt, not an artifact review lifecycle.**
The implemented review-like mechanism is execution gating: `interrupt(value)` (`libs/langgraph/langgraph/types.py:811-934`) pauses a node, surfaces `Interrupt(value, id)` to the caller, and resumes via `Command(resume=...)`. Graph-level `interrupt_before`/`interrupt_after` (`libs/langgraph/langgraph/pregel/main.py:767-769`, `libs/langgraph/langgraph/graph/state.py:1170-1252`) provides declarative pause points. The `StateSnapshot.tasks[].interrupts` (`libs/langgraph/langgraph/types.py:597-605`) and `__interrupt__` tuple allow observability. However there is no artifact review lifecycle: no `pending_approval`/`approved`/`rejected` states, no approver/comment, no promotion gate, no audit of who approved. The `require_approval` flag inspected in prebuilt tests (`libs/prebuilt/tests/test_react_agent.py:325`) relates to tool execution gating, not checkpoint review.

**3. Can artifacts be rolled back? — Indirectly via time-travel fork; no explicit rollback.**
There is no `rollback(checkpoint_id)` handler. Rollback is emulated by branching: fetch a historical `CheckpointTuple.config` via `get_state_history(before=...)` (`libs/langgraph/langgraph/pregel/main.py:1479-1530`) and either (a) replay by `invoke(None, before_b.config)` which re-executes nodes after that checkpoint and creates a new branch with `source="fork"` (`libs/langgraph/langgraph/pregel/main.py:1798-1823`), or (b) `update_state(historical_config, new_values)` / `bulk_update_state` (`libs/langgraph/langgraph/pregel/main.py:2530-2542`, `1589-1609`). Checkpoints are immutable and parents are retained (`parents` in metadata `libs/checkpoint/langgraph/checkpoint/base/__init__.py:56-60`), so history forms an auditable DAG, but the "bad" checkpoint remains reachable (not tombstoned). The only verb named `rollback` in the source is `MultitaskStrategy="rollback"` and `CancelAction="rollback"` in the SDK (`libs/sdk-py/langgraph_sdk/schema.py:81-86`, `137-141`), which deletes a *run* and its checkpoints — not an artifact-level revert with safeguards.

**4. Are artifact changes traceable to runs? — Partially, via metadata; not enforced end-to-end.**
`CheckpointMetadata.run_id` is defined (`libs/checkpoint/langgraph/checkpoint/base/__init__.py:61-62`) and persisted via `get_checkpoint_metadata` which merges `config.metadata`/`config.configurable` (`libs/checkpoint/langgraph/checkpoint/base/__init__.py:757-775`). `source`/`step`/`parents` provide provenance (`input` vs `loop` vs `update` vs `fork`). Tests assert these fields in history (`libs/langgraph/tests/test_time_travel.py:345-351`, `362-371`). However, in the OSS library path, `run_id` is whatever the caller puts in `RunnableConfig`; the `Pregel` loop itself populates `step`/`parents`/`source` but `run_id` has no mandatory server-side stamping in this source (contrast with Platform). Thus traceability exists structurally but depends on caller discipline; there is no immutable `author`/`commit` style attribution, and no change-log view tying diff+comment to run.

## Architectural Decisions

* **Immutable checkpoint DAG over mutable artifact store.** Decision to persist every superstep as a new `Checkpoint` with `id` (UUIDv6) and `parent_config` plus `parents` map (`libs/checkpoint/langgraph/checkpoint/base/__init__.py:139-145`, `libs/checkpoint/langgraph/checkpoint/memory/__init__.py:427-471`) and expose history via `list(before, limit, filter)` (`libs/checkpoint/langgraph/checkpoint/base/__init__.py:253-275`) enables time-travel without mutation. Tradeoff: storage grows without bounded GC; pruning naively breaks `DeltaChannel` (documented at `libs/checkpoint/langgraph/checkpoint/base/__init__.py:374-414`).
* **HITL via exception, not via artifact gate.** Chose `interrupt()` raising `GraphInterrupt` (`libs/langgraph/langgraph/types.py:927-934`) and resuming with `Command` (`libs/langgraph/langgraph/types.py:758-790`) coupled to `StateSnapshot.interrupts`. This unifies human approval with execution control but conflates run-level gating with artifact governance — there is no artifact-level approval state machine.
* **DeltaChannel + `get_delta_channel_history` beta storage optimization.** `DeltaChannel` persists only `writes` + periodic `_DeltaSnapshot` (`libs/checkpoint/langgraph/checkpoint/memory/__init__.py:142-229`); retrieval walks ancestors (`libs/checkpoint/langgraph/checkpoint/base/__init__.py:582-649`). Explicit `prune`/`delete_for_runs`/`copy_thread` warnings highlight fragility of history under lifecycle ops.
* **Pregel as orchestrator, saver as pluggable persistence.** Interface split between `Pregel` (`libs/langgraph/langgraph/pregel/main.py:748-835`) and `BaseCheckpointSaver` implementations (`memory`, `sqlite`, `postgres`) keeps versioning backend-agnostic but means no central "artifact registry" service with review/diff/rollback policies.

## Notable Patterns

* **Time-travel as versioning:** `checkpoint_id`-addressable history (`get_state` with or without `checkpoint_id` `libs/checkpoint/langgraph/checkpoint/memory/__init__.py:236-316`) + `get_state_history` + `update_state` fork (`source="fork"`). Demonstrated extensively in `libs/langgraph/tests/test_time_travel.py:69-520`.
* **Interrupt-as-review:** Every interrupt carries stable `Interrupt.id` (xxhash `libs/langgraph/langgraph/types.py:577-578`) and surfaces in `StateSnapshot.tasks[*].interrupts`, allowing a UI to render "pending review" without persisting a review record.
* **Parent-config lineage:** `patch_checkpoint_map` / `patch_configurable` pattern (`libs/langgraph/langgraph/pregel/main.py:1392-1432`, `libs/checkpoint/langgraph/checkpoint/memory/__init__.py:269-279`) maintains full ancestry for audit, fork detection (`saved.parent_config or {CHECKPOINT_ID: None}` `libs/langgraph/langgraph/pregel/main.py:1812-1815`).
* **Delta-aware storage:** `get_delta_channel_history` override in `InMemorySaver` (`libs/checkpoint/langgraph/checkpoint/memory/__init__.py:142-229`) avoids naïve `get_tuple` loop; plain blob subsumes pending writes while `_DeltaSnapshot` does not — semantic coupling is documented inline.

## Tradeoffs

* **Branching vs. rollback semantics.** Fork-based revert preserves auditability (old checkpoints remain) but leaves "bad" state reachable and requires caller to track the forked `config` as the new head. No atomic `rollback_to(checkpoint_id)` that atomically moves head and tombstones bad lineage.
* **Flexibility vs. safety for `DeltaChannel`.** Snapshot frequency reduces storage and ancestor walk length, but choosing prune/delete ranges incorrectly silently yields empty channel state (`seed` absent → consumer treats as empty, noted at `libs/checkpoint/langgraph/checkpoint/base/__init__.py:609-649`). No guard rails beyond doc warnings.
* **Generic metadata bag vs. typed annotation.** `CheckpointMetadata` (`total=False`) allows ad-hoc keys (enabling `run_id` stamping) but permits inconsistent provenance; no schema enforcement for comment/annotation/auditor fields.
* **Execution-level interrupts vs. artifact-level review.** Reusing `interrupt` for both execution gating and human approval avoids a second system but provides no review queue, SLA, or diff+comment affordance for approvers.

## Failure Modes / Edge Cases

* **Pruning DeltaChannel lineage → silent data loss.** `prune(thread_ids, strategy="keep_latest")` that drops intermediate checkpoints + `checkpoint_writes` severs reconstruction; the kept checkpoint reconstructs as empty with no error (`libs/checkpoint/langgraph/checkpoint/base/__init__.py:374-414` explains and proposes safe alternatives).
* **`delete_for_runs` / `prune` on shared ancestors.** Deleting a run that produced ancestor `checkpoint_writes` or the only `_DeltaSnapshot` corrupts still-live threads (`libs/checkpoint/langgraph/checkpoint/base/__init__.py:331-346`).
* **`copy_thread` incomplete ancestry.** Copying only the head checkpoint leaves `DeltaChannel` target with no path to snapshot (`libs/checkpoint/langgraph/checkpoint/base/__init__.py:350-372`).
* **Multitask `rollback` deletes checkpoints unexpectedly.** `CancelAction="rollback"` semantics (`libs/sdk-py/langgraph_sdk/schema.py:141-142`) deletes run + checkpoints; callers expecting soft revert may lose audit trail if Platform backing is used.
* **Ambiguous `update_state` without `as_node`.** Single-node heuristic vs. `versions_seen` ambiguity raises `InvalidUpdateError` (`libs/langgraph/langgraph/pregel/main.py:1933-1942`, `2400-2404`); concurrent forks from same parent create independent branches — last writer does not implicitly become head without explicitly threading the forked config.
* **No diff → unsafe promotion.** Absence of structured diff means reviewers cannot mechanically verify what changed between checkpoints without custom code; drift between two `StateSnapshot.values` may include non-serializable differences (e.g., `Interrupt` objects).
* **Interrupt resume mismatch.** `scratchpad.resume` indexed by `interrupt_counter` (`libs/langgraph/langgraph/types.py:912-918`); providing wrong count or stale `checkpoint_id` when invoking `Command(resume=...)` can resume the wrong interrupt or replay incorrectly (tested in `tests/test_time_travel.py:636-677`).

## Future Considerations

* Add a first-class `diff_state(before_config, after_config)` returning structured per-channel diff (using `channel_versions` + `get_delta_channel_history`), with tests and redaction for non-serializable values — closes the explicit diff gap without exposing storage internals.
* Introduce a typed `Rollback` operation (`rollback_to(config)` or `restore(checkpoint_id, reason)`) that creates a new `source="rollback"` checkpoint whose `parents` points to target and which tombstones the abandoned branch in metadata, rather than relying on implicit fork naming.
* Define a `CheckpointAnnotation` / comment model (e.g., `annotations: list[{author, timestamp, text, target_checkpoint_id, run_id}]`) persisted alongside metadata and surfaced in `StateSnapshot`, with `filter` support in `list`/`get_state_history`.
* Promote `run_id` stamping in OSS `PregelLoop` to always stamp the `run_id` from `RunnableConfig` into `CheckpointMetadata` (and optionally stamp `author`/`graph_version`), with integration tests asserting traceability end-to-end.
* Harden lifecycle ops for `DeltaChannel`: implement `prune` in `InMemorySaver`/`SqliteSaver`/`PostgresSaver` that walks back to snapshot ancestors or rewrites snapshots before deletion (as sketched in the warnings), and add conformance tests for prune+delta.

## Questions / Gaps

* Does LangGraph Platform (outside this library source) already implement artifact diff/review/rollback? No evidence found in `studies/agent-harness-study/sources/langgraph` — redirect at `docs/redirects.json:140 /cloud/how-tos/rollback_concurrent → https://docs.langchain.com/langsmith/rollback-concurrent` suggests platform docs may cover `rollback_concurrent`, but not present in this isolated source. Scope explicitly excludes sibling `langgraph-platform`/`server` if any.
* What is the formal contract for `run_id` population in OSS? `get_checkpoint_metadata` merges caller's `config.metadata`/`config.configurable` (`libs/checkpoint/langgraph/checkpoint/base/__init__.py:757-775`), but no `PregelLoop` code in this source unconditionally writes `run_id`; evidence of stamping searched but not found outside SDK's `run_id` handling (`libs/sdk-py/langgraph_sdk/schema.py:365-920`). State as `partially evidenced`.
* Is `filter` on `get_state_history` intended as the comment-search mechanism? `filter` matches exact metadata equality (`libs/checkpoint/langgraph/checkpoint/memory/__init__.py:374-379`); comment/annotation search would be prefix/substring — not supported.

---
Generated by `16.02-diff-review-and-rollback` against `langgraph`.
