# Source Analysis: langgraph

## Artifact Lifecycle

### Source Info

| Field | Value |
|-------|-------|
| Name | langgraph |
| Path | `studies/agent-harness-study/sources/langgraph` |
| Language / Stack | Python (monorepo: `checkpoint`, `checkpoint-postgres`, `checkpoint-sqlite`, `langgraph` core, `prebuilt`, `sdk-py/js`) |
| Analyzed | 2026-08-28 |

## Summary

LangGraph has no first-class “artifact” (file/blob) abstraction like MLflow or Flyte. The observable artifact lifecycle is bipartite: **(1) checkpoint-centric ephemeral execution state** (`Checkpoint` + `checkpoint_writes` + `checkpoint_blobs`) that is the unit of recovery/time-travel, and **(2) cross-thread durable memory** (`Item`/`SearchItem` via `BaseStore`/`InMemoryStore`). Both have explicit schemas, typed serde (`JsonPlusSerializer` → msgpack/pickle), and pluggable backends (in-memory dict, SQLite, Postgres). Creation is via `PregelLoop._put_checkpoint` → `create_checkpoint` → `saver.put`/`put_writes` with monotonic time-ordered `uuid6` IDs and per-channel `ChannelVersions`. Versioning is per-channel string (`000...032.<random>`) and checkpoint format `v`. Run linkage is via `CheckpointMetadata.run_id`/`source`/`step`/`parents` plus `RunnableConfig[configurable] = {thread_id, checkpoint_ns, checkpoint_id, checkpoint_map}`. Retirement is the weakest link: only `delete_thread` is concretely implemented; `delete_for_runs`/`prune`/`copy_thread` exist as abstract contracts with extensive `DeltaChannel` safety warnings but no durable-backend implementations, and no TTL/retention sweep for checkpoints (TTL exists only for `BaseStore` items and is opt-in/`supports_ttl=False` by default).

## Rating

**6 / 10** — Present but inconsistent/fragile. Schemas, storage backends, versioning, and run-linkage are explicit and tested; creation-to-retrieval is well-instrumented via `get_tuple`/`list`/`get_delta_channel_history`. Retirement/cleanup, garbage collection, and full artifact enumeration per run are incomplete: `delete_for_runs`/`prune` are not implemented in `SqliteSaver`/`PostgresSaver` (base raises `NotImplementedError`), `DeltaChannel` makes naïve deletion silently corrupting, and there is no lifecycle for external file artifacts (no blob store, no naming convention, no retention policy beyond manual `delete_thread`).

## Evidence Collected

Every entry MUST include a file path with line numbers. Format: `path/to/file.ts:NN`.

| Area | Evidence | File:Line |
|------|----------|-----------|
| Artifact schema — Checkpoint | `Checkpoint` TypedDict with `v, id, ts, channel_values, channel_versions, versions_seen, updated_channels` | `libs/checkpoint/langgraph/checkpoint/base/__init__.py:92-124` |
| Artifact schema — CheckpointMetadata | `source: "input"\|"loop"\|"update"\|"fork"`, `step`, `parents`, `run_id`, `counters_since_delta_snapshot` | `libs/checkpoint/langgraph/checkpoint/base/__init__.py:38-86` |
| Artifact schema — CheckpointTuple | Bundles `config, checkpoint, metadata, parent_config, pending_writes` | `libs/checkpoint/langgraph/checkpoint/base/__init__.py:139-146` |
| Artifact schema — PendingWrite/DeltaHistory | `PendingWrite = tuple[str,str,Any]`; `DeltaChannelHistory {writes, seed?}` | `libs/checkpoint/langgraph/checkpoint/base/__init__.py:31` , `libs/checkpoint/langgraph/checkpoint/base/__init__.py:149-174` |
| Artifact schema — Store Item | `Item {value,key,namespace,created_at,updated_at}` + `SearchItem {score}` + Ops `PutOp/GetOp/SearchOp/ListNamespacesOp` | `libs/checkpoint/langgraph/store/base/__init__.py:51-115` , `libs/checkpoint/langgraph/store/base/__init__.py:118-154` , `libs/checkpoint/langgraph/store/base/__init__.py:157-400` |
| Artifact schema — _DeltaSnapshot | Snapshot blob marker for `DeltaChannel` (msgpack ext `7`) | `libs/checkpoint/langgraph/checkpoint/serde/types.py:19-31` , `libs/checkpoint/langgraph/checkpoint/serde/jsonplus.py:302-307` |
| Creation — checkpoint factory | `empty_checkpoint() -> uuid6(clock_seq=-2)`, `create_checkpoint(... channels_to_snapshot ...)` writing `_DeltaSnapshot` | `libs/langgraph/langgraph/pregel/_checkpoint.py:26-34` , `libs/langgraph/langgraph/pregel/_checkpoint.py:61-121` |
| Creation — alternative factory | `empty_checkpoint` in base (v1 compat) and `create_checkpoint(channels, step, id=uuid6(clock_seq=step))` | `libs/checkpoint/langgraph/checkpoint/base/__init__.py:812-826` , `libs/checkpoint/langgraph/checkpoint/base/__init__.py:829-860` |
| Creation — loop integration | `PregelLoop._put_checkpoint` — builds `delta_channels_to_snapshot`, bumps counters, calls `submit(_checkpointer_put_after_previous, ... copy_checkpoint ...)` | `libs/langgraph/langgraph/pregel/_loop.py:1064-1200` |
| Creation — put_writes coalescing | `put_writes` drains `_delta_write_futs`/`_error_handler_write_futs`, ensures writes durable before checkpoint | `libs/langgraph/langgraph/pregel/_loop.py:408-510` , `libs/langgraph/langgraph/pregel/_loop.py:1201-1245` |
| Naming — ID generation | `uuid6` time-ordered UUID (monotonic, DB-locality) | `libs/checkpoint/langgraph/checkpoint/base/id.py:79-108` |
| Naming — runtime config keys | `thread_id, checkpoint_ns, checkpoint_id, checkpoint_map, checkpoint_id_saved` threaded through `RunnableConfig[configurable]` | `libs/langgraph/langgraph/pregel/_loop.py:340-365` , `libs/langgraph/langgraph/pregel/_checkpoint.py:142-184` |
| Storage backend — InMemorySaver | `storage: defaultdict[thread_id][ns][id] -> (checkpoint, metadata, parent)`, `writes: (tid,ns,cid)->(task,idx)`, `blobs: (tid,ns,channel,version)` | `libs/checkpoint/langgraph/checkpoint/memory/__init__.py:68-83` , `libs/checkpoint/langgraph/checkpoint/memory/__init__.py:453-471` |
| Storage backend — Postgres DDL | `MIGRATIONS` create `checkpoints`, `checkpoint_blobs`, `checkpoint_writes` with PK `(thread_id, checkpoint_ns, checkpoint_id[, channel, version|task_id, idx])` | `libs/checkpoint-postgres/langgraph/checkpoint/postgres/base.py:43-91` |
| Storage backend — SQLite DDL | `setup()` creates `checkpoints` + `writes` (no separate blobs table; channel_values inline) | `libs/checkpoint-sqlite/langgraph/checkpoint/sqlite/__init__.py:139-163` |
| Storage backend — blobs vs inline | InMemory `_load_blobs`/`blobs` map; Postgres `_dump_blobs`/`_load_blobs` with `ON CONFLICT DO NOTHING`; Sqlite stores `channel_values` inline in checkpoint blob | `libs/checkpoint/langgraph/checkpoint/memory/__init__.py:125-140` , `libs/checkpoint-postgres/langgraph/checkpoint/postgres/base.py:328-337` , `libs/checkpoint-postgres/langgraph/checkpoint/postgres/base.py:479-502` |
| Storage backend — Store (memory) | `_data: dict[ns][key]->Item`, `_vectors: dict[ns][key][path]->embeddings`, `supports_ttl` off by default | `libs/checkpoint/langgraph/store/memory/__init__.py:186-204` , `libs/checkpoint/langgraph/store/base/__init__.py:719-720` |
| Version identifiers — checkpoint format | `Checkpoint.v` + `LATEST_VERSION=4` (pregel) vs `2` (base compat) + `checkpoint/ts,id` | `libs/checkpoint/langgraph/checkpoint/base/__init__.py:95-96` , `libs/checkpoint/langgraph/checkpoint/base/__init__.py:811` , `libs/langgraph/langgraph/pregel/_checkpoint.py:21` |
| Version identifiers — channel versions | `ChannelVersions = dict[str,str|int|float]`, `get_next_version: f"{v:032}.{random}"` | `libs/checkpoint/langgraph/checkpoint/base/__init__.py:89` , `libs/checkpoint/langgraph/checkpoint/base/__init__.py:692-711` , `libs/checkpoint/langgraph/checkpoint/memory/__init__.py:619-628` , `libs/checkpoint-postgres/langgraph/checkpoint/postgres/base.py:543-552` |
| Version identifiers — msgpack ext | `EXT_DELTA_SNAPSHOT=7`, serializer `dumps_typed/loads_typed` with `"msgpack"/"empty"/"null"` envelope | `libs/checkpoint/langgraph/checkpoint/serde/jsonplus.py:302` , `libs/checkpoint/langgraph/checkpoint/serde/jsonplus.py:258-290` |
| Run-artifact association — metadata | `metadata.run_id`, `parents`, `source`, `step` via `get_checkpoint_metadata` merging `config.metadata` + `config.configurable` | `libs/checkpoint/langgraph/checkpoint/base/__init__.py:752-775` |
| Run-artifact association — config | `thread_id` primary key pattern; `checkpoint_ns` for subgraphs (`NS_SEP` join); `checkpoint_id` ordering | `libs/checkpoint/langgraph/checkpoint/base/__init__.py:176-207` docstring, `libs/langgraph/langgraph/pregel/_loop.py:322-365` |
| Run-artifact association — queries | `list(where thread_id=? ...)`, `search_where` building metadata `@> filter` (Postgres) / JSON filter (Sqlite) | `libs/checkpoint-postgres/langgraph/checkpoint/postgres/base.py:554-596` , `libs/checkpoint-sqlite/langgraph/checkpoint/sqlite/__init__.py:334-385` |
| Run-artifact association — DeltaChannel walk | `get_delta_channel_history` walks `parent_config` chain collecting `pending_writes` until `_DeltaSnapshot` seed | `libs/checkpoint/langgraph/checkpoint/base/__init__.py:582-649` , `libs/checkpoint/langgraph/checkpoint/memory/__init__.py:142-229` |
| Retention — delete_thread | Concrete deletes in all three savers: `DELETE FROM checkpoints|blobs|writes WHERE thread_id=?` / `del storage[thread_id]` | `libs/checkpoint/langgraph/checkpoint/memory/__init__.py:511-527` , `libs/checkpoint-postgres/langgraph/checkpoint/postgres/__init__.py:381-402` , `libs/checkpoint-sqlite/langgraph/checkpoint/sqlite/__init__.py:484-501` |
| Retention — abstract-only | `delete_for_runs`, `copy_thread`, `prune` raise `NotImplementedError` with DeltaChannel corruption warnings | `libs/checkpoint/langgraph/checkpoint/base/__init__.py:320-415` , `libs/checkpoint/langgraph/checkpoint/base/__init__.py:511-580` |
| Retention — store TTL | `TTLConfig {refresh_on_read, default_ttl, sweep_interval_minutes}`, `PutOp.ttl`, `_ensure_ttl/_ensure_refresh`, `supports_ttl` gate | `libs/checkpoint/langgraph/store/base/__init__.py:545-567` , `libs/checkpoint/langgraph/store/base/__init__.py:526-534` , `libs/checkpoint/langgraph/store/base/__init__.py:911-926` |
| Retention — store delete | `BaseStore.delete` via `PutOp(value=None)`, `InMemoryStore._apply_put_ops` pops `_data` + `_vectors` | `libs/checkpoint/langgraph/store/base/__init__.py:929-936` , `libs/checkpoint/langgraph/store/memory/__init__.py:404-416` |
| Observability | `get_tuple`, `list` (`ORDER BY checkpoint_id DESC LIMIT`), `get_delta_channel_history` paged two-stage SQL + pending_sends migration | `libs/checkpoint/langgraph/checkpoint/memory/__init__.py:236-426` , `libs/checkpoint-postgres/langgraph/checkpoint/postgres/__init__.py:112-191` , `libs/checkpoint-postgres/langgraph/checkpoint/postgres/base.py:201-288` |

## Answers to Dimension Questions

### 1. What types of artifacts exist?

LangGraph does not expose a generic “artifact” file/blob API. The durable artifacts are:

* **Checkpoint** (`Checkpoint` @ `libs/checkpoint/langgraph/checkpoint/base/__init__.py:92-124`): `v,id,ts,channel_values,channel_versions,versions_seen,updated_channels,pending_sends`. The `channel_values` map holds serialized state per channel; large/non-primitive values are externalized to **checkpoint_blobs** (Postgres) or inline (SQLite).
* **CheckpointMetadata** (`libs/checkpoint/langgraph/checkpoint/base/__init__.py:38-86`): sidecar JSONB/BLOB `{source, step, parents, run_id, counters_since_delta_snapshot}`. `writes` key is stripped before persist (`libs/checkpoint/langgraph/checkpoint/base/__init__.py:778-785`).
* **Checkpoint Writes / PendingWrite** (`tuple[str,str,Any]` @ `libs/checkpoint/langgraph/checkpoint/base/__init__.py:31`): intermediate `(task_id, channel, value)` rows in `checkpoint_writes`/`writes` tables, keyed by `(thread_id, ns, checkpoint_id, task_id, idx)` with `WRITES_IDX_MAP` (`libs/checkpoint/langgraph/checkpoint/base/__init__.py:795`) mapping `ERROR→-1` etc.
* **Checkpoint Blobs** (`checkpoint_blobs` rows): `(thread_id, ns, channel, version) -> (type, blob)` (Postgres `libs/checkpoint-postgres/langgraph/checkpoint/postgres/base.py:57-65`; InMemory `blobs` dict `libs/checkpoint/langgraph/checkpoint/memory/__init__.py:78-83`). `DeltaChannel` uses sentinel `True` in `channel_values` + `_DeltaSnapshot` blob (`libs/checkpoint/langgraph/checkpoint/serde/types.py:19-31`).
* **Store Items** (`Item`, `SearchItem` @ `libs/checkpoint/langgraph/store/base/__init__.py:51-154`): long-term key-value memory with `namespace: tuple[str,...]`, `key`, `value: dict`, `created_at`, `updated_at`, optional `score`, optional vector index. Operations are `PutOp/GetOp/SearchOp/ListNamespacesOp`.
* **No file artifact**: no `Artifact`, `Blob`, `Attachment`, or file-store abstraction was found in the searched source (grep for `artifact` returned only CI upload actions).

### 2. How are artifacts named and stored?

**Naming:**
* Primary namespace is `thread_id` (caller-supplied `config["configurable"]["thread_id"]`); within a thread, `checkpoint_ns` ( subgraph path `NS_SEP`-joined, default `""` — `libs/langgraph/langgraph/pregel/_loop.py:360-363`) partitions the checkpoint chain; final key is `checkpoint_id` which is a **time-ordered `uuid6`** (`libs/checkpoint/langgraph/checkpoint/base/id.py:79-108`, minted as `str(uuid6(clock_seq=step))` in `libs/langgraph/langgraph/pregel/_checkpoint.py:116` and `libs/checkpoint/langgraph/checkpoint/base/__init__.py:854`). Ordering by `checkpoint_id DESC` yields newest-first (`libs/checkpoint/langgraph/checkpoint/memory/__init__.py:283`, `libs/checkpoint-sqlite/langgraph/checkpoint/sqlite/__init__.py:240`).
* Channel blobs are named by `(thread_id, checkpoint_ns, channel, version)` where `version` is per-channel monotonic string `f"{n:032}.{rand:016}"` (`libs/checkpoint/langgraph/checkpoint/memory/__init__.py:622-628`). No user-visible “artifact name” slug; the Store uses `(namespace tuple, key string)` with validation banning `.` in labels and `langgraph` root (`libs/checkpoint/langgraph/store/base/__init__.py:1255-1275`).
* No enforced naming convention for run-scoped artifact collections; `list_namespaces(prefix/suffix, max_depth)` (`libs/checkpoint/langgraph/store/base/__init__.py:938-992`) is the namespace enumeration primitive.

**Storage:**
* **InMemorySaver** (`libs/checkpoint/langgraph/checkpoint/memory/__init__.py:68-83`): three `defaultdict`s in process memory (`storage`, `writes`, `blobs`). `put` serializes via `serde.dumps_typed` (`libs/checkpoint/langgraph/checkpoint/memory/__init__.py:458-461`), writes blobs only for `new_versions` intersecting `channel_values`.
* **PostgresSaver** (`libs/checkpoint-postgres/langgraph/checkpoint/postgres/__init__.py:85-111` setup/`MIGRATIONS`): tables `checkpoints` (JSONB `checkpoint`+`metadata`), `checkpoint_blobs` (BYTEA blob + `type`), `checkpoint_writes` (BYTEA). `select` uses correlated subquery `array_agg(array[bl.channel, bl.type, bl.blob]) FROM jsonb_each_text(channel_versions) JOIN checkpoint_blobs` (`libs/checkpoint-postgres/langgraph/checkpoint/postgres/base.py:93-117`). Uses `psycopg` + pipeline mode and `threading.Lock`.
* **SqliteSaver** (`libs/checkpoint-sqlite/langgraph/checkpoint/sqlite/__init__.py:139-163`): `checkpoints` (`BLOB` checkpoint+metadata) + `writes` only; channel values stored **inline** in the checkpoint blob (no `checkpoint_blobs` table), which changes the `get_delta_channel_history` implementation (`libs/checkpoint-sqlite/langgraph/checkpoint/sqlite/__init__.py:503-583`).
* **Store**: `InMemoryStore._data`/`_vectors` (`libs/checkpoint/langgraph/store/memory/__init__.py:186-189`), with optional embedding index (`dims, embed, fields` @ `libs/checkpoint/langgraph/store/base/__init__.py:570-698`). Batching via `AsyncBatchedBaseStore` coalesces `batch` ops (`libs/checkpoint/langgraph/store/base/batch.py:58-371`).
* **Serde**: `JsonPlusSerializer` msgpack with `ormsgpack` ext codes (`EXT_DELTA_SNAPSHOT=7`, etc. @ `libs/checkpoint/langgraph/checkpoint/serde/jsonplus.py:294-302`); primitive values (`None,str,int,float,bool`) are stored inline in Postgres `checkpoint` JSONB (`libs/checkpoint-postgres/langgraph/checkpoint/postgres/__init__.py:312-319`).

### 3. Are artifacts versioned?

Yes, but in two orthogonal senses:

* **Checkpoint format version** `Checkpoint.v: int` (`libs/checkpoint/langgraph/checkpoint/base/__init__.py:95-96`). Current `LATEST_VERSION=4` (`libs/langgraph/langgraph/pregel/_checkpoint.py:21`; compat `2` in base @ `libs/checkpoint/langgraph/checkpoint/base/__init__.py:811`). Migration path for pending sends `v<4` exists (`libs/checkpoint-postgres/langgraph/checkpoint/postgres/__init__.py:164-188`).
* **Per-channel monotonic versions** `ChannelVersions` (`libs/checkpoint/langgraph/checkpoint/base/__init__.py:89`). `BaseCheckpointSaver.get_next_version` defaults to `int+1` (`libs/checkpoint/langgraph/checkpoint/base/__init__.py:692-711`); concrete savers override to `f"{v:032}.{random}"` (`libs/checkpoint/langgraph/checkpoint/memory/__init__.py:619-628`, `libs/checkpoint-postgres/langgraph/checkpoint/postgres/base.py:543-552`). Each `put` computes `new_versions = diff(previous_versions, channel_versions)` via `get_new_channel_versions` (`libs/langgraph/langgraph/pregel/_loop.py:1173-1176`) and only those `version` blobs are upserted.
* **Checkpoint ID as sort key**: `uuid6` is monotonic and encodes timestamp (`libs/checkpoint/langgraph/checkpoint/base/id.py:92-98`), so `ORDER BY checkpoint_id DESC` is both creation order and version order.
* **Store Item versioning**: implicit via `created_at`/`updated_at` timestamps on `Item` (`libs/checkpoint/langgraph/store/base/__init__.py:60-88`); `put` overwrites in-place (no history) and `_dedupe_ops` collapses duplicate `PutOp`s (`libs/checkpoint/langgraph/store/base/batch.py:283-323`). No explicit artifact version number.

### 4. Can artifacts be linked to the run that produced them?

Partially — checkpoints are, store items are not, and the “run” concept is thread-centric rather than run-centric.

* **Checkpoint → Run**: `CheckpointMetadata.run_id: str` (`libs/checkpoint/langgraph/checkpoint/base/__init__.py:61-62`) plus `source` (`input|loop|update|fork`), `step: int`, and `parents: dict[ns, checkpoint_id]` (`libs/checkpoint/langgraph/checkpoint/base/__init__.py:41-60`). Metadata is merged from `config["metadata"]` and `config["configurable"]` via `get_checkpoint_metadata` (`libs/checkpoint/langgraph/checkpoint/base/__init__.py:752-775`), which also backfills any `run_id` passed as configurable. The Pregel loop stamps `metadata.step = self.step` and `metadata.parents = config[CHECKPOINT_MAP]` on every `_put_checkpoint` (`libs/langgraph/langgraph/pregel/_loop.py:1108-1109`).
* **Checkpoint → Thread/Namespace**: `get_tuple`/`list` filter on `thread_id` + `checkpoint_ns` + optional `checkpoint_id` (`libs/checkpoint-postgres/langgraph/checkpoint/postgres/base.py:554-596` comment `WHERE thread_id=%s AND checkpoint_ns=%s`; same for SQLite via `search_where`). `parent_config` chain (`libs/checkpoint/langgraph/checkpoint/base/__init__.py:269-280`) lets any checkpoint be traced to its parent and ultimately enumerated via `list` (`LIMIT`/`before` pagination — `libs/checkpoint/langgraph/checkpoint/memory/__init__.py:318-425`).
* **Writes → Checkpoint**: `checkpoint_writes` rows carry `(thread_id, ns, checkpoint_id, task_id, idx, channel, blob)` so writes are queryable per checkpoint (`libs/checkpoint-postgres/langgraph/checkpoint/postgres/base.py:66-76`). `get_delta_channel_history` reconstructs a channel’s lineage across the parent chain (`libs/checkpoint/langgraph/checkpoint/base/__init__.py:582-649`).
* **Store Items → Run**: no link. `Item` has no `run_id`, `thread_id`, or `checkpoint_id` field; namespace/key are user-defined (`libs/checkpoint/langgraph/store/base/__init__.py:51-64`). You cannot ask “which items did run X produce” without encoding that out-of-band in `namespace`/`value`.
* **Gaps**: there is no global `run_id → [checkpoint_ids]` index; `list` can filter `metadata @> {run_id: ...}` as JSONB containment (`libs/checkpoint-postgres/langgraph/checkpoint/postgres/base.py:584-586`) but that requires scanning metadata and is not exposed as a first-class “artifacts for run” API. `delete_for_runs(run_ids)` is the intended affordance (`libs/checkpoint/langgraph/checkpoint/base/__init__.py:331-348`) yet no concrete saver implements it. The question “find every artifact produced by a given run” is answerable for checkpoints via `list(filter={"run_id": ...})` (Postgres) but not reliably across backends, and not for store items at all.

### 5. How are artifacts retired?

Retirement is **manual + single-thread granularity**, with no automated retention/GC and explicit DeltaChannel hazards:

* **delete_thread(thread_id)** — the only fully implemented path: InMemory deletes `storage`, `writes`, `blobs` dicts (`libs/checkpoint/langgraph/checkpoint/memory/__init__.py:511-527`); Postgres deletes three tables (`libs/checkpoint-postgres/langgraph/checkpoint/postgres/__init__.py:381-402`); Sqlite deletes `checkpoints` + `writes` (`libs/checkpoint-sqlite/langgraph/checkpoint/sqlite/__init__.py:484-501`); async variants mirror it (`libs/checkpoint-postgres/langgraph/checkpoint/postgres/aio.py:340`, `libs/checkpoint-sqlite/langgraph/checkpoint/sqlite/aio.py:602`).
* **delete_for_runs(run_ids)** — declared abstract with DeltaChannel corruption warning (“deleting a run that produced ancestor `checkpoint_writes` or the only `_DeltaSnapshot` blob … will break reconstruction” @ `libs/checkpoint/langgraph/checkpoint/base/__init__.py:331-348`) and `raise NotImplementedError` in the base; no Sqlite/Postgres override was found (grep confirmed only base). So per-run retirement is **not operational**.
* **prune(thread_ids, strategy="keep_latest"|"delete")** — same status (`libs/checkpoint/langgraph/checkpoint/base/__init__.py:374-414`): abstract, warned as DeltaChannel-naïve (`keep_latest` can sever the snapshot chain and silently reconstruct as empty). No concrete implementation ships.
* **copy_thread** — also abstract with warning that copying only the head checkpoint breaks DeltaChannel walks (`libs/checkpoint/langgraph/checkpoint/base/__init__.py:350-372`).
* **Store retirement** — `delete(namespace, key)` → `PutOp(value=None)` (`libs/checkpoint/langgraph/store/base/__init__.py:929-936`, applied in `libs/checkpoint/langgraph/store/memory/__init__.py:404-408`); `put(..., ttl=float_minutes)` + `TTLConfig {default_ttl, refresh_on_read, sweep_interval_minutes}` (`libs/checkpoint/langgraph/store/base/__init__.py:545-567`) with per-`GetOp`/`SearchOp` `refresh_ttl` (`libs/checkpoint/langgraph/store/base/__init__.py:760-763`). Runtime gating: `supports_ttl=False` by default (`libs/checkpoint/langgraph/store/base/__init__.py:719`) and `put` raises `NotImplementedError` if `ttl != None` without support (`libs/checkpoint/langgraph/store/base/__init__.py:912-916`). Only `InMemoryStore`’s sweep is not time-driven (no background sweeper in inspected code); Postgres store’s Postgres-specific TTL (if any) lives outside the inspected files. No checkpoint TTL.
* **Result**: retention is explicit (caller must invoke `delete_thread`); there is no lifecycle hook, no size/age policy, and the documented safe retirement patterns for DeltaChannel (“walk back to nearest snapshot, or force resnapshot” @ `libs/checkpoint/langgraph/checkpoint/base/__init__.py:387-412`) are advisory, not enforced.

## Architectural Decisions

| Decision | Evidence | Rationale / Consequence |
|----------|----------|--------------------------|
| Checkpoint as source of truth, not file artifacts | `Checkpoint`/`CheckpointTuple` as only `BaseCheckpointSaver` return type; no `Artifact` class in repo | Keeps graph persistence transportable; external files must be handled by user code or Store values. |
| Time-ordered `uuid6` IDs as both PK and sort key | `libs/checkpoint/langgraph/checkpoint/base/id.py:79-108`, `ORDER BY checkpoint_id DESC LIMIT 1` everywhere | Gives DB locality and cheap “latest” without a separate timestamp index; clock skew handled via `_last_v6_timestamp` bump. |
| Per-channel version strings + `blobs` externalization | `ChannelVersions` + `blobs[(tid,ns,ch,ver)]` (`libs/checkpoint/langgraph/checkpoint/memory/__init__.py:78-83`); Postgres `checkpoint_blobs` (`libs/checkpoint-postgres/langgraph/checkpoint/postgres/base.py:57-65`) | Allows copy-on-write per channel; unchanged channels cost no blob duplication. |
| Inline primitives optimization (Postgres) | `isinstance(v, (str,int,float,bool,None))` keep in JSONB, else externalize (`libs/checkpoint-postgres/langgraph/checkpoint/postgres/__init__.py:312-319`) | Reduces `checkpoint_blobs` row count for scalar state; diverges from Sqlite’s fully inline model — operational asymmetry. |
| DeltaChannel sentinel + ancestor walk | `DeltaChannel.checkpoint() -> MISSING` (`libs/langgraph/langgraph/channels/delta.py:195-204`), `_DeltaSnapshot` msgpack ext 7, `get_delta_channel_history` paged two-stage SQL (`libs/checkpoint-postgres/langgraph/checkpoint/postgres/base.py:201-288`) | Scales append-heavy reducers (e.g., message lists) by storing only writes between snapshots; tradeoff is expensive reconstruction and fragile pruning. |
| `thread_id` + `checkpoint_ns` as composite namespace | Every `put`/`get_tuple`/`_cursor` keys on both (`libs/checkpoint-postgres/langgraph/checkpoint/postgres/__init__.py:295-309`) | Enables subgraph time-travel (`ns` per subgraph) but makes “all artifacts for a run” a cross-namespace query. |
| Abstract `delete_for_runs`/`prune` with warnings | Extensive docstring hazards (`libs/checkpoint/langgraph/checkpoint/base/__init__.py:340-414`) without implementations | Signals known safety gaps for DeltaChannel; leaves lifecycle enforcement to deployer. |

## Notable Patterns

* **Two-stage DeltaChannel reconstruction** — Stage 1 scans only metadata (`checkpoint_id, parent_checkpoint_id, ver_i/hs_i` via parallel JSONB lookups) paging newest-first; Stage 2 fetches only the identified `writes` + `seed` blobs via per-channel `UNION ALL` (`libs/checkpoint-postgres/langgraph/checkpoint/postgres/base.py:175-289`, `libs/checkpoint-sqlite/langgraph/checkpoint/sqlite/_delta.py` counterpart). Avoids shipping full `channel_values` JSON.
* **Msgpack typed envelope** — `dumps_typed -> (type, bytes)` where `type in {"msgpack","empty","null","bytes"}` plus allowlist-gated deserialization (`libs/checkpoint/langgraph/checkpoint/serde/jsonplus.py:97-126`, `258-290`), with `_DeltaSnapshot` as a dedicated ext code.
* **Batched Store** — `AsyncBatchedBaseStore` dedupes and coalesces `batch(ops)` in a background `asyncio.Queue` task, with `batch` collapsing duplicate `PutOp`s to last-write-wins (`libs/checkpoint/langgraph/store/base/batch.py:283-371`).
* **Walks over parent chain, not `list(before=...)`** — `get_delta_channel_history` follows `parent_config` not temporal scan, so forked threads only accumulate on-path ancestors (`libs/checkpoint/langgraph/checkpoint/base/__init__.py:604-609`).
* **Migration sentinel** — `MIGRATIONS` list versioned via `checkpoint_migrations` table (`libs/checkpoint-postgres/langgraph/checkpoint/postgres/base.py:43-91`), with `SELECT_PENDING_SENDS_SQL` migration for `v<4` pending sends.

## Tradeoffs

* **Durability vs speed**: `durability="exit"` defers `put` until loop exit and stages delta writes lazily (`libs/langgraph/langgraph/pregel/_loop.py:213-222`, `1201-1245`); faster per-step but risks loss on mid-run crash. The write-checkpoint ordering invariant is preserved via `submit(...put_after_previous...)` futures.
* **Granular blobs vs operational complexity**: per-channel versioning minimizes write amplification, but doubles conceptual artifacts (checkpoint row + N blobs + M writes) and requires careful GC (hence the unimplemented `prune`).
* **Scalability of DeltaChannel**: bounded replay depth via `snapshot_frequency` + `DELTA_MAX_SUPERSTEPS_SINCE_SNAPSHOT` (`libs/langgraph/langgraph/channels/delta.py:52-55`, `libs/langgraph/langgraph/pregel/_checkpoint.py:40-58`); however snapshot decision is per-step and can be double-counted on exit (`libs/langgraph/langgraph/pregel/_loop.py:1094-1113` notes).
* **Postgres vs Sqlite fidelity**: Postgres supports out-of-line blobs and metadata `@>` filtering; Sqlite’s inline-only model is simpler but ships full checkpoint blobs for Stage 1 (`libs/checkpoint-sqlite/langgraph/checkpoint/sqlite/__init__.py:503-523` docstring) — performance/observability mismatch that surfaces in DeltaChannel walks.
* **Safety vs extensibility**: strict msgpack allowlist (`LANGGRAPH_STRICT_MSGPACK`) trades deserialization safety for backward-compat breakage on old JSON checkpoints; the default is permissive with warnings (`libs/checkpoint/langgraph/checkpoint/serde/jsonplus.py:97-126`, `547-610`).

## Failure Modes / Edge Cases

* **Prune/delete corrupts DeltaChannel silently** — `get_delta_channel_history` returns `no seed` → consumer starts empty, not an error (`libs/checkpoint/langgraph/checkpoint/base/__init__.py:397-401` warning: “silently reconstruct as empty (no error raised … simply returns no `seed`)”). No safe default `prune` exists; operator must implement snapshot-preserving walk or force snapshot before delete.
* **Concurrent `uuid6` clock_seq collision** — `random.getrandbits(48)` for node + monotonic `_last_v6_timestamp` bump (`libs/checkpoint/langgraph/checkpoint/base/id.py:92-102`) prevents duplicate timestamps within a process, but cross-process collisions rely on random node; `step`-seeded `clock_seq` in `create_checkpoint` (`libs/langgraph/langgraph/pregel/_checkpoint.py:116`) reuses `step` as sequence, so two concurrent loops at same `step` could theoretically collide.
* **Store TTL without sweeper** — `InMemoryStore` has no background TTL sweep in inspected code; expired items are deleted only “opportunistically” (per Store docstring) — unbounded memory growth if TTL is used without external reaping. Postgres store’s TTL (if any) is not in the inspected `BaseStore` interface.
* **Namespace validation rejects valid user intent** — `put` rejects empty namespace, `.` in label, or `langgraph` root (`libs/checkpoint/langgraph/store/base/__init__.py:1255-1275`); bulk `batch` does not validate until `put` path, so raw `batch([PutOp(...)])` can bypass validation.
* **`get_delta_channel_history` on missing target** — returns empty `writes` per channel (`libs/checkpoint-postgres/langgraph/checkpoint/postgres/__init__.py:470-473`, `libs/checkpoint-sqlite/langgraph/checkpoint/sqlite/__init__.py:531-532`); caller `channels_from_checkpoint` treats that as `MISSING` → `typ()` empty value, which may mask a truly dangling thread.
* **Inline primitive divergence** — Postgres stores primitives inline, Sqlite stores everything inline; restoring a Postgres checkpoint on Sqlite (or vice versa) via serialized JSON would change blob presence, breaking `hs_i` detection in Stage 1.
* **Unbounded `list` without `limit`** — `BaseCheckpointSaver.list` yields with `limit` decrement (`libs/checkpoint/langgraph/checkpoint/memory/__init__.py:382-386`); Postgres appends `LIMIT %s` only if `limit is not None` (`libs/checkpoint-postgres/langgraph/checkpoint/postgres/__init__.py:154-157`) — a `list(config)` with no `limit` can stream the entire thread history.

## Future Considerations

* **Implement `delete_for_runs`/`prune` concretely** with snapshot-aware GC: walk to nearest `_DeltaSnapshot` per delta channel before deleting ancestors, or synthesize a snapshot via `DeltaChannel.from_checkpoint(seed).replay_writes(writes)` then `put` before prune — exactly what the existing docstrings prescribe but don’t enforce.
* **Add checkpoint TTL/retention policy** (age, count, size) mirroring `TTLConfig` for Store, with a background sweeper or `prune` cron; expose metrics for `counters_since_delta_snapshot` growth to alert on replay depth approaching `DELTA_MAX_SUPERSTEPS_SINCE_SNAPSHOT`.
* **First-class file artifact store** — introduce an `ArtifactStore` (S3/GCS/local) with naming `{thread_id}/{checkpoint_id}/{channel}/{version}` and content hash, linked via `CheckpointMetadata.artifacts: list[{channel, uri, etag}]` and `Store` reference from `Item.value` (e.g., `{"$artifact": "s3://..."}`) to answer “all files for run X”.
* **Unified blob table for SQLite** — add `checkpoint_blobs` to SQLite to eliminate the inline-vs-external divergence and enable the faster K-JSONB-lookup Stage 1 without shipping full checkpoint blobs.
* **Enumerate artifacts per run** — implement `BaseCheckpointSaver.list_for_run(run_id)` → `(checkpoints, writes, blobs)` and `BaseStore.list_for_run(run_id)` or at minimum document the metadata `@>` query pattern and index `metadata` GIN (`checkpoints` already has `thread_id` indices — `libs/checkpoint-postgres/langgraph/checkpoint/postgres/base.py:82-88`).
* **Harden `copy_thread`** to copy the complete parent chain up to snapshot ancestors (docstring @ `libs/checkpoint/langgraph/checkpoint/base/__init__.py:350-371`) and make it atomic via Postgres pipeline.

## Questions / Gaps

* **No file/blob artifact primitive** — “artifacts” in this dimension likely expects file outputs (model checkpoints, rendered images, etc.). LangGraph has no such abstraction; we treated `checkpoint_blobs`/`writes` as artifacts. Confirm whether user-produced files are out-of-scope or should be modeled as `Store` items with inline bytes.
* **No evidence of retention sweep implementation** — `TTLConfig.sweep_interval_minutes` is documented but no sweeper loop was found in `libs/checkpoint/langgraph/store/memory/__init__.py` or `libs/checkpoint/langgraph/store/base/__init__.py`. Is TTL a Postgres-only feature?
* **Version identifiers for store items?** — `Item` has no version field; is overwriting `put` the intended versioning story, or is history expected via checkpointing?
* **Where are run IDs assigned?** — `run_id` appears in metadata but no assignment site was traced in this study (likely LangSmith/client). Without that trace, the “linked to run” claim relies on docstring rather than observed write path.
* **Can every artifact for a run be found?** — With currentcode, checkpoints yes via `list(filter={"run_id": ...})` on Postgres only; writes/blobs only via follow-up `get_delta_channel_history`; store items no. A cross-artifact manifest does not exist.

---

Generated by `studies/agent-harness-study/dimensions/16.01-artifact-lifecycle.md` against `langgraph`.
