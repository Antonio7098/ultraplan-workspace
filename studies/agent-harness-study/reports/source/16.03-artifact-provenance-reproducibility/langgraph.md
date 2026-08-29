# Source Analysis: langgraph

## 16.03 Artifact Provenance and Reproducibility

### Source Info

| Field | Value |
|-------|-------|
| Name | langgraph |
| Path | `studies/agent-harness-study/sources/langgraph` |
| Language / Stack | Python (monorepo: libs/langgraph, libs/checkpoint, libs/checkpoint-sqlite/postgres, libs/prebuilt) |
| Analyzed | 2026-08-28 |

## Summary

LangGraph's "artifacts" are checkpointed graph states (`Checkpoint`, `CheckpointTuple`, `StateSnapshot`), not end-user deliverables (documents, code patches). Provenance is centered on checkpoint lineage, not on prompts, model calls, tool invocations, or human approvals. `CheckpointMetadata` (`libs/checkpoint/langgraph/checkpoint/base/__init__.py:38-86`) records only `source` (`input|loop|update|fork`), `step`, `parents`, `run_id`, and beta `counters_since_delta_snapshot`. Parent linkage is via `CheckpointTuple.parent_config` (`libs/checkpoint/langgraph/checkpoint/base/__init__.py:139-146`) and `config["configurable"]["checkpoint_id"]` / `checkpoint_map`. Raw inputs, per-node prompts, LLM/tool calls, and approval identity are **not** persisted in the artifact; they exist only as opaque channel writes or external LangSmith traces. Re-execution from a checkpoint (`get_state` / replay / `update_state` fork) is the reproducibility primitive, exercised heavily in `tests/test_time_travel.py`, but determinism is undermined by random UUID/ version generation and retry jitter, and no CI job asserts bit-for-bit reproducibility or provenance completeness.

## Rating

**Score: 4 / 10 — Present but inconsistent, weakly documented, fragile**

Rationale: Checkpoint lineage (step, source, parents, run_id, pending_writes) provides coarse provenance and replayability via `get_state`/`get_state_history` and `get_delta_channel_history`. However prompts, model identifiers, tool arguments, full input payloads, context, and approval metadata are not recorded in the artifact; `get_checkpoint_metadata` filters aggressively; IDs/versions are non-deterministic; and no reproducibility or provenance-coverage tests run in CI. This satisfies "present but inconsistent" (4-6), at the lower end.

## Evidence Collected

Every entry MUST include a file path with line numbers. Format: `path/to/file.ts:NN`.

| Area | Evidence | File:Line |
|------|----------|-----------|
| Provenance fields — metadata schema | `CheckpointMetadata` TypedDict defines only `source`, `step`, `parents`, `run_id`, `counters_since_delta_snapshot` (beta). No prompt/model/tool/context/approval fields. | `libs/checkpoint/langgraph/checkpoint/base/__init__.py:38-87` |
| Provenance fields — checkpoint body | `Checkpoint` TypedDict holds `v`, `id` (uuid6), `ts`, `channel_values`, `channel_versions`, `versions_seen`, `updated_channels`. No input hash, model version, or tool provenance. | `libs/checkpoint/langgraph/checkpoint/base/__init__.py:92-124` |
| Provenance fields — tuple lineage | `CheckpointTuple` bundles `config`, `checkpoint`, `metadata`, `parent_config`, `pending_writes`. Lineage is parent pointer, not full contributing-factor log. | `libs/checkpoint/langgraph/checkpoint/base/__init__.py:139-146` |
| Provenance fields — state snapshot | `StateSnapshot` surfaces `values`, `next`, `config`, `metadata`, `created_at`, `parent_config`, `tasks`, `interrupts`. Again no prompt/model provenance. | `libs/langgraph/langgraph/types.py:643-662` |
| Input recording — metadata merge | `get_checkpoint_metadata(config, metadata)` copies only `str/int/float/bool` keys from `config["configurable"]`+`config["metadata"]`, skipping `EXCLUDED_METADATA_KEYS` and `__`-prefixed keys; `writes` field explicitly dropped in serializable variant. | `libs/checkpoint/langgraph/checkpoint/base/__init__.py:757-786`, `libs/checkpoint/langgraph/checkpoint/base/__init__.py:797-807` |
| Input recording — materialization | `PregelLoop._first` maps inputs via `map_input` / `map_command` then `apply_writes`; raw input not hashed or logged separately — only resulting `channel_values` persist. | `libs/langgraph/langgraph/pregel/_loop.py:836-1017` |
| Input recording — checkpoint creation | `create_checkpoint` snapshots only `channel_values` that exist in `channel_versions`; `CheckpointMetadata` populated with `step` and `parents` in `_put_checkpoint`. No prompt/context capture. | `libs/langgraph/langgraph/pregel/_checkpoint.py:61-121`, `libs/langgraph/langgraph/pregel/_loop.py:1064-1142` |
| Input recording — delta channel seed/writes | Delta channels not stored in `channel_values` for sub-frequency steps; history reconstructed via `get_delta_channel_history` ancestor walk accumulating `pending_writes` until `_DeltaSnapshot` seed. Fragile if intermediate writes deleted. | `libs/langgraph/langgraph/pregel/_checkpoint.py:136-184`, `libs/checkpoint/langgraph/checkpoint/base/__init__.py:582-649`, `libs/checkpoint/langgraph/checkpoint/memory/__init__.py:142-229` |
| Input recording — approval/interrupt | `interrupt(value)` raises `GraphInterrupt` with `Interrupt.from_ns(value, ns)`; resume via `Command(resume=...)`; persisted only as `INTERRUPT`/`RESUME` pending writes, no approver identity or timestamp. | `libs/langgraph/langgraph/types.py:811-934`, `libs/langgraph/langgraph/pregel/_loop.py:806-832` |
| Reproducibility — replay/fork primitive | `StateSnapshot` prepared via `channels_from_checkpoint` / `achannels_from_checkpoint` replaying delta histories; `SyncPregelLoop`/`AsyncPregelLoop` re-executes nodes after a checkpoint on `invoke(None, checkpoint.config)` or after `update_state` fork. | `libs/langgraph/langgraph/pregel/_checkpoint.py:136-226`, `libs/langgraph/langgraph/pregel/main.py:1144-1265` |
| Reproducibility — deterministic task ordering | Task ordering explicitly sorted: `prepare_next_tasks` sorts nodes to ensure deterministic order. | `libs/langgraph/langgraph/pregel/_algo.py:253`, `libs/langgraph/langgraph/pregel/_algo.py:481` |
| Reproducibility — nondeterministic ID generation | `uuid6` uses `time.time_ns()` + `random.getrandbits(14/48)`; `InMemorySaver.get_next_version` uses `random.random()` — checkpoint IDs and channel versions differ on every run. | `libs/checkpoint/langgraph/checkpoint/base/id.py:79-108`, `libs/checkpoint/langgraph/checkpoint/memory/__init__.py:619-628` |
| Reproducibility — nondeterministic retry | Retry intervals add `random.uniform(0,1)` jitter even when `jitter=True` (default), making re-execution timing nondeterministic. | `libs/langgraph/langgraph/pregel/_retry.py:672`, `libs/langgraph/langgraph/pregel/_retry.py:828` |
| Reproducibility — delta message ID fix | Comment acknowledges prior nondeterminism: unsanitized `id=None` BaseMessages caused `get_state()` replays to produce different UUIDs; mitigated by `ensure_message_ids` before serializing delta writes. | `libs/langgraph/langgraph/pregel/_loop.py:452-458` |
| Reproducibility tests — time travel | ~40 tests assert replay determinism: `test_replay_reruns_nodes_after_checkpoint`, `test_replay_from_before_interrupt_refires`, `test_replay_interrupt_stable_across_replays`, `test_sequential_interrupts_fork_from_middle`, etc. Verify `source`/`next`/`values` but not bit-identical artifact hash. | `libs/langgraph/tests/test_time_travel.py:69-109`, `libs/langgraph/tests/test_time_travel.py:226-282`, `libs/langgraph/tests/test_time_travel.py:395-443` |
| Reproducibility tests — subgraph replay | Subgraph replay tests verify accumulated state via `ReplayState` (`list(before=parent_checkpoint_id)`) and fork/resume completion, not determinism of IDs or external inputs. | `libs/langgraph/tests/test_time_travel.py:1180-1321`, `libs/langgraph/langgraph/_internal/_replay.py:68`, `libs/langgraph/langgraph/pregel/_loop.py:1619-1656` |
| Reproducibility tests — delta channel | Conformance suite and delta channel tests exercise `get_delta_channel_history` and migration, e.g., `test_delta_channel_history`, `test_delta_channel_migration`, but do not assert full input-to-artifact reproducibility. | `libs/checkpoint-conformance/langgraph/checkpoint/conformance/spec/test_delta_channel_history.py:88`, `libs/langgraph/tests/test_delta_channel_migration.py:230` |
| Build/reproduction scripts | Monorepo `pyproject.toml` / `uv.lock` / `Makefile` exist per lib, but no `reproduce.py`, provenance manifest, SBOM, or lockfile-to-artifact linkage documented. Build backend is `hatchling`. | `libs/langgraph/pyproject.toml:1-3`, `libs/langgraph/uv.lock:1`, `AGENTS.md:7-15` |
| External tracing (not artifact-embedded) | LangSmith tracer integration via `LangChainTracer` / `FakeTracer` used in `test_tracing_interops.py` and `test_large_cases.py` — model/tool provenance lives in external trace service, not checkpoint. | `libs/langgraph/tests/test_tracing_interops.py:11-117`, `libs/langgraph/tests/fake_tracer.py:7` |
| CI coverage | `ci.yml` runs `lint` + `test` matrices for `libs/langgraph` and checkpoint libs, plus SDK checks; no dedicated artifact-provenance or determinism job. | `.github/workflows/ci.yml:58-107` |
| Serde determinism | `JsonPlusSerializer` uses `ormsgpack` with `OPT_NON_STR_KEYS | OPT_PASSTHROUGH_*`; `dumps_typed`/`loads_typed` protocol defined but no canonicalization or hash recorded per artifact. | `libs/checkpoint/langgraph/checkpoint/serde/jsonplus.py:258-290`, `libs/checkpoint/langgraph/checkpoint/serde/base.py:14-26` |

## Answers to Dimension Questions

**1. Can every artifact be traced to its inputs?**
No. A checkpoint (`StateSnapshot` / `CheckpointTuple`) is traceable to its parent chain via `CheckpointMetadata.parents` (`libs/checkpoint/langgraph/checkpoint/base/__init__.py:56-60`) + `parent_config` (`libs/checkpoint/langgraph/checkpoint/base/__init__.py:145`) and to the `thread_id`/`checkpoint_id` in `RunnableConfig` (`libs/langgraph/langgraph/_internal/_config.py:411-420`). However the *semantic* inputs — the original `invoke` payload, per-node prompts, LLM model name/version, tool arguments, and human approval — are not recorded in the artifact. `get_checkpoint_metadata` only propagates a filtered subset of `config["configurable"]` primitives (`libs/checkpoint/langgraph/checkpoint/base/__init__.py:757-775`) and `create_checkpoint` stores only the resulting `channel_values` (`libs/langgraph/langgraph/pregel/_checkpoint.py:82-113`). Deltas compound this: intermediate writes are scattered across `checkpoint_writes` rows and must be reconstructed via ancestor walk (`libs/checkpoint/langgraph/checkpoint/memory/__init__.py:142-229`); deleting an ancestor or its writes (warned in `BaseCheckpointSaver.prune`/`delete_for_runs` `libs/checkpoint/langgraph/checkpoint/base/__init__.py:331-415`) silently breaks reconstruction. No evidence was found of an artifact manifest linking artifact hash → input hash → prompt → model call.

**2. Is reproduction deterministic?**
Partially. Graph scheduling is deterministic (`prepare_next_tasks` sorts by path, `libs/langgraph/langgraph/pregel/_algo.py:253,481`). Replaying from a checkpoint deterministically re-executes nodes *after* that checkpoint (tested in `libs/langgraph/tests/test_time_travel.py:69-443`). But byte-for-byte determinism is not achieved: `uuid6` in `libs/checkpoint/langgraph/checkpoint/base/id.py:92-108` uses wall-clock time + `random.getrandbits`; `InMemorySaver.get_next_version` in `libs/checkpoint/langgraph/checkpoint/memory/__init__.py:619-628` injects `random.random()` into version strings; retry backoff in `libs/langgraph/langgraph/pregel/_retry.py:672,828` injects jitter. The fix for `id=None` BaseMessage nondeterminism (`libs/langgraph/langgraph/pregel/_loop.py:452-458`) shows prior nondeterminism was discovered and partially patched, but fresh `ts` (`datetime.now(timezone.utc).isoformat()` in `libs/langgraph/langgraph/pregel/_checkpoint.py:80`) and new UUIDs on every `create_checkpoint` (`libs/langgraph/langgraph/pregel/_checkpoint.py:116`) still guarantee differing artifact bytes on replay. No determinism property test or seeded RNG mode was found.

**3. Are all contributing factors recorded?**
No. Recorded: `source`/`step`/`parents`/`run_id`/`counters_since_delta_snapshot` (`libs/checkpoint/langgraph/checkpoint/base/__init__.py:38-86`), channel values/versions (`libs/checkpoint/langgraph/checkpoint/base/__init__.py:104-114`), and pending writes including `INTERRUPT`/`RESUME` (`libs/langgraph/langgraph/_internal/_constants.py` via `WRITES_IDX_MAP`). Missing: prompt text, LLM provider/model/version, temperature/seed, tool definitions and return values beyond channel writes, full `config` (filtered by `EXCLUDED_METADATA_KEYS` `libs/checkpoint/langgraph/checkpoint/base/__init__.py:797-807`), environment (`LANGGRAPH_DEFAULT_RECURSION_LIMIT` read via `getenv` in `libs/langgraph/langgraph/_internal/_config.py:32` but not snapshotted), approval actor, and trace ID. Contributing factors split between checkpoint (graph state) and external LangSmith traces (`libs/langgraph/tests/test_tracing_interops.py:68-117`), with no unified provenance envelope.

**4. Is reproducibility tested in CI?**
No. `.github/workflows/ci.yml:58-107` runs lint and `pytest` for each lib (including `libs/langgraph` via `_test_langgraph.yml`). Reproducibility-adjacent tests exist (`tests/test_time_travel.py`, `tests/test_pregel.py`, `tests/test_interruption.py`, delta-channel suites) and assert logical replay correctness, but no CI job asserts deterministic rehydration, provenance completeness, hash stability, or build reproducibility. No `reproducibility_test`, `snapshot_hash_check`, or artifact-manifest validator was found in workflows. Searched `ci.yml`, `_test.yml`, `_test_langgraph.yml` — no provenance gate.

## Architectural Decisions

- **Checkpoint-as-artifact model** (`libs/checkpoint/langgraph/checkpoint/base/__init__.py:176-207`, `libs/langgraph/langgraph/pregel/main.py:1391-1432`): The system treats the persisted `Checkpoint` + `CheckpointTuple` as the reproducible artifact, not user-facing files. This enables time-travel (`invoke(None, before.config)`) and forks (`update_state`) but conflates framework state with user artifact provenance.
- **Minimal metadata allowance** (`libs/checkpoint/langgraph/checkpoint/base/__init__.py:797-807`, `libs/checkpoint/langgraph/checkpoint/base/__init__.py:757-785`): `EXCLUDED_METADATA_KEYS` plus `__`-prefix filter and val-type allowlist keeps checkpoint rows small, at the cost of dropping rich provenance.
- **Delta-channel storage optimization** (`libs/langgraph/langgraph/pregel/_checkpoint.py:37-58`, `libs/checkpoint/langgraph/checkpoint/memory/__init__.py:142-229`, `libs/langgraph/langgraph/pregel/_loop.py:1064-1142`): Delta channels avoid blob bloat by scattering writes across ancestors; history reconstructed on read. Trades storage efficiency for provenance fragility (documented warnings in `BaseCheckpointSaver.prune` `libs/checkpoint/langgraph/checkpoint/base/__init__.py:374-415`).
- **Externalized observability** (`libs/langgraph/tests/test_tracing_interops.py:68-117`, `libs/langgraph/langgraph/pregel/_loop.py:125-126`): Model/tool provenance delegated to LangSmith/LangChain tracers via `RunnableConfig.callbacks`, not embedded in checkpoint. Decouples provenance from persistence.
- **Non-deterministic versioning** (`libs/checkpoint/langgraph/checkpoint/base/id.py:79-108`, `libs/checkpoint/langgraph/checkpoint/memory/__init__.py:619-628`): `uuid6` + random `get_next_version` favor DB locality and collision avoidance over reproducibility.

## Notable Patterns

- **Parent-chain lineage pattern**: Every put (`InMemorySaver.put` `libs/checkpoint/langgraph/checkpoint/memory/__init__.py:427-471` and SQLite/Postgres equivalents) stores `parent_checkpoint_id`, enabling `get_state_history` traversal and `ReplayState` lookup (`libs/langgraph/langgraph/_internal/_replay.py:68`). Pattern is consistent across all saver implementations.
- **Ancestor-walk replay pattern**: `get_delta_channel_history` default impl (`libs/checkpoint/langgraph/checkpoint/base/__init__.py:582-649`) and optimized `InMemorySaver` override (`libs/checkpoint/langgraph/checkpoint/memory/__init__.py:142-229`) both walk ancestors until finding `_DeltaSnapshot` (`libs/checkpoint/langgraph/checkpoint/serde/types.py:18-31`) or plain seed. Single-saver call batches all delta channels.
- **Patch-config lineage enrichment**: `patch_checkpoint_map` (`libs/langgraph/langgraph/_internal/_config.py:63-80`) and `merge_configs` (`libs/langgraph/langgraph/_internal/_config.py:147-191`) thread provenance-bearing keys (`checkpoint_map`, `parents`) through `RunnableConfig` without explicit provenance object.

## Tradeoffs

- **Storage efficiency vs. provenance durability**: Delta channels + exit-mode durability (`Durability = "sync"|"async"|"exit"` in `libs/langgraph/langgraph/types.py:87-93` and `libs/langgraph/langgraph/pregel/_loop.py:1116-1118`) reduce write amplification but increase risk of irreproducibility if intermediate checkpoints/writes are pruned (warnings in `BaseCheckpointSaver.prune` `libs/checkpoint/langgraph/checkpoint/base/__init__.py:374-415` are advisory, not enforced).
- **Filtered metadata vs. completeness**: `EXCLUDED_METADATA_KEYS` and primitive-only filtering (`libs/checkpoint/langgraph/checkpoint/base/__init__.py:769-774`) prevents checkpoint bloat and leakage of `thread_id`/`checkpoint_ns` into user metadata, but discards context that would be needed for full traceability.
- **External tracing vs. embedded provenance**: Relying on LangSmith means rich model/tool traces do not inflate checkpoints, but artifacts stored outside LangSmith (e.g., `InMemorySaver` in tests) have zero model/tool provenance.
- **Randomized IDs vs. determinism**: `uuid6` clock-seq + random node + `get_next_version` jitter improve concurrency and DB ordering but prevent bit-stable replay — a contributor to `test_replay_interrupt_stable_across_replays` (`libs/langgraph/tests/test_time_travel.py:395-443`) needing to compare only values/interrupt values, not full equality.

## Failure Modes / Edge Cases

- **Ancestor pruning silently corrupts delta state**: `BaseCheckpointSaver.prune(..., strategy="keep_latest")` and `delete_for_runs` warnings (`libs/checkpoint/langgraph/checkpoint/base/__init__.py:340-415`) state that dropping intermediate `checkpoint_writes` or the sole `_DeltaSnapshot` causes `get_delta_channel_history` to return no `seed` and channels reconstruct as empty — no error raised, silent data loss.
- **Exit-mode durability loses sub-frequency input**: `PregelLoop._first` input writes for delta channels are only durable via `put_writes` in `sync`/`async` modes (`libs/langgraph/langgraph/pregel/_loop.py:1006-1013`) and via `_exit_delta_writes` accumulation (`libs/langgraph/langgraph/pregel/_loop.py:697-703`); a crash before exit in `Durability="exit"` loses those inputs even though channel state appears to progress.
- **`get_checkpoint_metadata` silent drop**: Non-primitive configurable values (dicts, lists, secrets filtered by `_exclude_as_metadata` in `libs/langgraph/langgraph/_internal/_config.py:423-447`) are silently omitted from checkpoint metadata; downstream `list(filter=...)` queries (`libs/checkpoint/langgraph/checkpoint/memory/__init__.py:318-380`) then cannot locate artifacts by those keys.
- **ID collision on clock skew**: `uuid6` fallback to `random.getrandbits(14)` for `clock_seq` (`libs/checkpoint/langgraph/checkpoint/base/id.py:99-100`) plus `random.random()` version suffix can theoretically collide or disorder under forked processes without shared RNG seed.
- **Replay vs. resume misclassification**: `PregelLoop._first` is-time-traveling heuristic (`libs/langgraph/langgraph/pregel/_loop.py:849-888`) gates on `is_replaying` + `source` check; a resume that supplies an explicit `checkpoint_id` equal to head can be misclassified as time-travel (mitigated in `libs/langgraph/langgraph/pregel/_loop.py:1026-1050` but still fragile).
- **Serde allowlist regression**: `JsonPlusSerializer._create_msgpack_ext_hook` (`libs/checkpoint/langgraph/checkpoint/serde/jsonplus.py:546-744`) returns `None`/raw fallback on blocked types; a stricter `LANGGRAPH_STRICT_MSGPACK=true` rollout can make previously reproducible artifacts unloadable without `allowed_msgpack_modules` migration.

## Future Considerations

- Introduce a first-class **provenance envelope** alongside `CheckpointMetadata` (e.g., `provenance: {input_hash, prompt_hash, model_id, tool_call_ids, run_id, approver}`) persisted atomically with checkpoint, or at least a `writes` hash. Current `run_id` is insufficient.
- Add **deterministic mode** for tests/repro runs: seed `random`, optionally derive `checkpoint.id` / `channel_versions` from a content hash rather than `uuid6`/`random.random()`, and disable retry jitter.
- Emit a **repro manifest / SBOM** per run (invoke config + `uv.lock` hash + channel values hash) and validate it in CI with a `replay --assert-deterministic` job.
- Enforce **DeltaChannel-aware pruning** (walk-to-snapshot check) rather than advisory warnings; add conformance test that prunes then asserts `get_state` still reproduces.
- Propagate `LANGGRAPH_DEFAULT_RECURSION_LIMIT` / `DELTA_MAX_SUPERSTEPS_SINCE_SNAPSHOT` (`libs/langgraph/langgraph/_internal/_config.py:32-34`) into `CheckpointMetadata` so environment contributing factors are captured.

## Questions / Gaps

- No evidence found of how LangGraph records **LLM prompts, completions, tool arguments, or model versions** inside the persisted artifact — searched `libs/checkpoint/**:*.py`, `libs/langgraph/langgraph/pregel/**:*.py`, `libs/langgraph/langgraph/types.py:467-980`. Appears intentionally out-of-scope (delegated to LangSmith callbacks).
- No evidence found of **context/approval provenance** (who provided `Command(resume=...)` value, what `ServerInfo`/`BaseUser` was active). `Runtime` (`libs/langgraph/langgraph/runtime.py`) carries `ServerInfo`/`RunControl` but not persisted in checkpoint.
- No evidence found of **reproducibility tested in CI** — `.github/workflows/ci.yml:1-198`, `_test.yml`, `_test_langgraph.yml`, `_integration_test.yml` contain no deterministic-replay or provenance-completeness gate.
- No evidence found of **artifact hash / input hash** fields — `Checkpoint.id` is time-based, not content-based; `ts` is wall-clock. Searched `Checkpoint`/`CheckpointMetadata` definitions and `create_checkpoint`.
- No evidence found of **build/reproduction script** for artifacts — `Makefile`/`pyproject.toml`/`uv.lock` manage package builds, not artifact reproduction; no `scripts/reproduce.py` or equivalent inventoried.

---

Generated by `16.03-artifact-provenance-and-reproducibility` against `langgraph`.
