# Source Analysis: openhands

## Dimension 16.01: Artifact Lifecycle

### Source Info

| Field | Value |
|-------|-------|
| Name | openhands |
| Path | `studies/agent-harness-study/sources/openhands` |
| Language / Stack | Python (FastAPI, SQLAlchemy, Pydantic), React/TS, Docker/GCS/S3, Redis |
| Analyzed | 2026-08-28 |

## Summary

OpenHands has no single `Artifact` abstraction; its "artifacts" are the conversation run itself: `StoredConversationMetadata` (SQL), per-event JSON blobs (file/object store), workspace-archive objects (git-delta/tar.gz + manifest), and an ephemeral trajectory ZIP. Creation and storage are well-factored behind pluggable backends (`FileStore`, `EventService`), linking is via `conversation_id.hex` + `sandbox_id`, and pre-delete workspace capture with `REQUIRED` gating is durably engineered. Versioning is minimal (`conversation_version='V1'`, overwrite-on-save), and retirement is explicit-DELETE only — no TTL, GC, or unified inventory API. You cannot answer "every artifact for run X" from one call; you must stitch DB rows, event-store listing, and object-store prefix listing.

## Rating

**Score: 6 / 10 — Present but inconsistent / fragile**

Rationale: Creation, naming, multi-backend storage, run linkage, and durability capture are clear, tested, and have operational safeguards (streaming uploads, atomic writes, Redis export lock, retry semantics). Inconsistencies pull the score down: no artifact version history or content addressing, events are mutable overwrites, workspace archives are not indexed beside the conversation row (location is env-derived), and retirement has no retention/cleanup policy, no orphan reap, and leaves event blobs behind if DB delete succeeds but sandbox teardown fails. Proven under failure for the happy-path delete, but not under scale/storage-leak scenarios.

## Evidence Collected

Every entry MUST include a file path with line numbers. Format: `path/to/file.ts:NN`.

| Area | Evidence | File:Line |
|------|----------|-----------|
| Artifact schema — conversation metadata | `StoredConversationMetadata` SQL model: `conversation_id` PK, `sandbox_id`, `conversation_version='V1'`, metrics columns, `tags` JSON, indexes on `sandbox_id`/`conversation_version`/`execution_status` | `openhands/app_server/app_conversation/sql_app_conversation_info_service.py:132-189` |
| Artifact schema — cost ledger | `StoredConversationCostEvent` per-bucket deltas: `conversation_id FK CASCADE`, `cost_delta`, `usage_id`, `llm_model`, token columns, `occurred_at` | `openhands/app_server/app_conversation/sql_app_conversation_info_service.py:192-210` |
| Artifact schema — API model | `AppConversationInfo` (info without live status), `AppConversation` (+ `sandbox_status`, `execution_status`, `session_api_key`), `AppConversationStartTask`, `AppConversationInfoPage` | `openhands/app_server/app_conversation/app_conversation_models.py:110-335` |
| Artifact schema — events | `EventService.save_event / get_event / search_events / iter_events_for_export` abstract contract; events are `Event` Pydantic models serialized as indented JSON | `openhands/app_server/event/event_service.py:18-71` |
| Artifact creation — conversation | `save_app_conversation_info` upserts via `merge`, preserves `created_at` on update, merges `MetricsSnapshot`/`TokenUsage` into columns, commits | `openhands/app_server/app_conversation/sql_app_conversation_info_service.py:440-503` |
| Artifact creation — events | `_store_event` creates parent dirs and writes `event.model_dump_json(indent=2)`; `_load_event` validates JSON back to `Event` | `openhands/app_server/event/filesystem_event_service.py:33-36` and `openhands/app_server/event/aws_event_service.py:54-62` and `openhands/app_server/event/google_cloud_event_service.py:57-62` |
| Artifact creation — workspace archive | `archive_workspace()` pulls `GET /api/file/archive?path=...&format=...`, streams to tempfile, `write_from_path` → object store, writes `*.manifest.json` with repo metadata, packages, sizes | `openhands/app_server/sandbox/workspace_archive.py:334-541` |
| Artifact creation — trajectory export | `_stream_conversation_zip` writes `meta.json` + `event_000000_{id}.json` per event into `zipfile.ZipFile(ZIP_DEFLATED)` via `_StreamingZipBuffer`; `open_conversation_export` adds Redis lock + TTL refresh | `openhands/app_server/app_conversation/live_status_app_conversation_service.py:2869-2962` |
| Storage backend — file abstraction | `FileStore` ABC: `write/read/list/delete`, `write_from_path` streaming override note for OOM avoidance | `openhands/app_server/file_store/files.py:8-42` |
| Storage backend — local | `LocalFileStore` atomic write (temp + `os.replace` + `fsync`), `shutil.copyfile` streaming for `write_from_path`, `list` adds trailing `/` for dirs, `delete` handles files vs dirs | `openhands/app_server/file_store/local.py:26-82` |
| Storage backend — cloud | `AwsEventService` via `boto3 s3_client.get_object/put_object/list_objects_v2`; `GoogleCloudEventService` via `google.cloud.storage` bucket/blobs, shared `Client` LRU cache | `openhands/app_server/event/aws_event_service.py:27-76` and `openhands/app_server/event/google_cloud_event_service.py:37-70` |
| Storage backend — selection | `config_from_env()` selects `FilesystemEventServiceInjector` vs `AwsEventServiceInjector` vs `GoogleCloudEventServiceInjector` via `StorageProvider` + `FILE_STORE_PATH`; persistence dir default `~/.openhands` or `OH_PERSISTENCE_DIR` | `openhands/app_server/config.py:307-328` and `openhands/app_server/config.py:75-89` |
| Storage backend — archive factory | `_get_archive_file_store()` via `get_file_store(_archive_store_type(), _archive_bucket())` — type from `RUNTIME_FILE_ARCHIVE_STORE_TYPE` (default `google_cloud`), explicitly rejects `memory` (binary corruption) | `openhands/app_server/sandbox/workspace_archive.py:155-164` and `openhands/app_server/file_store/memory.py:12-13` |
| Naming / pathing | `get_conversation_path` / `V1_CONVERSATIONS_DIR='v1_conversations'` → `{prefix}/{user_id}/v1_conversations/{hex}`; event file `{id.hex}.json`; archive `{prefix}/{sandbox_id}/{conversation_id}/{ts}.{patch\|tar.gz}` + `.{suffix}.manifest.json` | `openhands/app_server/event/event_service_base.py:66-84` and `openhands/app_server/conversation_paths.py:12-73` and `openhands/app_server/sandbox/workspace_archive.py:398-403` |
| Version identifiers | `conversation_version='V1'` hard-filter in `_secure_select`, `search`, `count`; `ACP_SERVER_TAG_KEY='acpserver'` + `ARCHIVE_WORKSPACE_PATH_TAG_KEY='archiveworkspacepath'` tag conventions; `agent_kind` column | `openhands/app_server/app_conversation/sql_app_conversation_info_service.py:755-759` and `openhands/app_server/app_conversation/app_conversation_models.py:34-48` |
| Run-artifact association | Conversation ↔ events via `conversation_id.hex` directory; ↔ workspace archive via `sandbox_id` + `conversation_id` in object key; ↔ sandbox grouping via `count_conversations_by_sandbox_id`; parent/sub via `parent_conversation_id` + `get_sub_conversation_ids` | `openhands/app_server/event/event_service_base.py:86-96` and `openhands/app_server/sandbox/workspace_archive.py:398-403` and `openhands/app_server/app_conversation/sql_app_conversation_info_service.py:369-395` |
| Retention — delete path | `DELETE /{conversation_id}` commits DB delete, then detached `_finalize_sandbox_delete`: archive first, only then `delete_sandbox` if `count==0`; if `REQUIRED` archive fails, sandbox+runtime kept for idle-reap backstop; rollback on finalizer exception | `openhands/app_server/app_conversation/app_conversation_router.py:885-937` and `openhands/app_server/app_conversation/live_status_app_conversation_service.py:2690-2831` |
| Retention — what is NOT deleted | `delete_app_conversation` cascades DB rows + agent-server `DELETE /api/conversations/{id}`, but event blobs under `{user_id}/v1_conversations/{hex}/` are never deleted; no TTL/lifecycle code found for DB or object store (only bucket lifecycle comment) | `openhands/app_server/app_conversation/live_status_app_conversation_service.py:2779-2831` and `openhands/app_server/sandbox/workspace_archive.py:408-411` |
| Observability / safeguards | Export lock: `try_acquire_redis_lock(key, ttl)` + `refresh_lock_periodically` every `ttl//2`; size guard `export_max_events=10000`; streaming tempfile + `write_from_path` to bound RAM; base-commit / repo-metadata headers enriched into manifest | `openhands/app_server/utils/redis_lock.py:25-52` and `openhands/app_server/app_conversation/live_status_app_conversation_service.py:2845-2912` and `openhands/app_server/sandbox/workspace_archive.py:176-210` |
| Migrations | `conversation_metadata` created in `001` with `sandbox_id`/`conversation_version` indexes; later adds `parent_conversation_id`, `public`, `tags`, `agent_kind`, `execution_status` | `openhands/app_server/app_lifespan/alembic/versions/001.py:158-197` and `openhands/app_server/app_lifespan/alembic/versions/003.py:24-28` etc. |

## Answers to Dimension Questions

### 1. What types of artifacts exist?

Four primary run artifacts (plus supporting rows):

* **Conversation metadata** — `StoredConversationMetadata` (also `StoredConversationCostEvent` ledger) and its API projection `AppConversationInfo`/`AppConversation`. Holds identity (`conversation_id`, `sandbox_id`), repo/branch/provider, title/trigger/pr_numbers, LLM model, `agent_kind`, accumulated cost/tokens (`prompt_tokens`, `completion_tokens`, `cache_read_tokens`, `cache_write_tokens`, `reasoning_tokens`, `context_window`, `per_turn_token`), `tags` map, timestamps. Created on start, updated on stats/events/execution-status (`openhands/app_server/app_conversation/sql_app_conversation_info_service.py:132-210`, `openhands/app_server/app_conversation/app_conversation_models.py:110-219`).
* **Events / trajectory** — Append-only (but technically overwritable) JSON files, one per `Event.id`, under the conversation's directory. Schema is the SDK `Event` (validated on load). This is the durable execution trace (`openhands/app_server/event/event_service.py:18-71`, `openhands/app_server/event/event_service_base.py:86-96`).
* **Workspace archives** — Durability copies of the sandbox filesystem at delete time: `*.patch` (git-delta) and/or `*.tar.gz` (full) plus `*.manifest.json` with `sandbox_id`, `conversation_id`, `phase='final'`, `base_commit`, repo metadata (`repo_remote`, `branch`, `head_commit`), `packages`, `environment`, `format`, `source_path`, `byte_count`, `created_at`. Written to the configured `FileStore` (GCS by default) (`openhands/app_server/sandbox/workspace_archive.py:498-517`).
* **Trajectory export ZIP** — Ephemeral, not persisted: streaming ZIP containing `meta.json` + sorted `event_{i:06d}_{id}.json` entries. Produced on-demand via `GET /export` returning `application/zip` (`openhands/app_server/app_conversation/live_status_app_conversation_service.py:2869-2894`).
* Supporting: `app_conversation_start_task` rows, `event_callback`/`event_callback_result` rows, `v1_remote_sandbox` rows, `pending_messages`.

No dedicated "build artifact / model artifact / dataset" abstraction — the workspace archive is the closest to a generic artifact store.

### 2. How are artifacts named and stored?

**Conversations (SQL):** PK `conversation_id` as hex string (`str(UUID)`). Table `conversation_metadata` indexed on `conversation_version`, `sandbox_id`, `execution_status`, `public`. `sandbox_id` equals remote `v1_remote_sandbox.id` (often `conversation_id.hex` in cloud). Stored via SQLAlchemy/SQL async session (SQLite in dev/tests, Postgres in prod). (`openhands/app_server/app_conversation/sql_app_conversation_info_service.py:135-184`).

**Events (file/object store):** Path template `{prefix}/{user_id}/v1_conversations/{conversation_id.hex}/{event_id.hex}.json` (`V1_CONVERSATIONS_DIR='v1_conversations'`). `prefix` is `OH_PERSISTENCE_DIR` / `~/.openhands` for filesystem, or empty logical prefix for S3/GCS (bucket key = same string). `user_id` resolved from `UserContext` or lazily via `AppConversationInfoService` if missing — this is the permission boundary (`openhands/app_server/event/event_service_base.py:66-96`, `openhands/app_server/conversation_paths.py:12-73`). Naming is content-agnostic (UUID hex), no hash.

**Workspace archives (object store):** `{RUNTIME_FILE_ARCHIVE_PREFIX}/{sandbox_id}/{conversation_id}/{ts}.{suffix}` where `PREFIX` defaults `workspace-archives`, `sandbox_id` is sandbox, `conversation_id` is `conversation_id.hex` (prevents sibling overwrite under grouping), `ts='%Y%m%dT%H%M%SZ'`, `suffix∈{patch, tar.gz}`. Companion `{key}.manifest.json`. Store type from `RUNTIME_FILE_ARCHIVE_STORE_TYPE` (local/s3/google_cloud) and bucket from `RUNTIME_FILE_ARCHIVE_BUCKET`, via `get_file_store()` factory (`openhands/app_server/sandbox/workspace_archive.py:95-103`, `openhands/app_server/sandbox/workspace_archive.py:398-403`, `openhands/app_server/sandbox/workspace_archive.py:162-164`).

**Exports:** Not stored — streamed as `conversation_{id}.zip` with `Content-Disposition` attachment (`openhands/app_server/app_conversation/app_conversation_router.py:1668-1670`). Naming for entries is deterministic `event_000000_{id}.json` sorted by timestamp (`openhands/app_server/app_conversation/live_status_app_conversation_service.py:2884`).

Storage is pluggable: `FileStore` (local/S3/GCS/memory) for file artifacts, `EventService` (filesystem/AWS/GCS) for events, SQL for metadata — selected via `StorageProvider` in `config.py:307-328`.

### 3. Are artifacts versioned?

**Mostly no — overwrite model with one sentinel version.**

* Conversation row: `conversation_version` is always `'V1'` for active conversations; filters hardcode `== 'V1'` (`openhands/app_server/app_conversation/sql_app_conversation_info_service.py:312-313`, `755-759`). No per-artifact revision history; `save_app_conversation_info` does an upsert (`merge`) that overwrites the same row. `created_at` is specially preserved across overwrites to avoid corruption (`openhands/app_server/app_conversation/sql_app_conversation_info_service.py:458-466`), but other fields (title, tags, metrics) are replaced. No `version` counter or etag.
* Events: file name is `event_id.hex.json`; `save_event` writes `event.model_dump_json` to that path, overwriting if the same `event_id` is saved again. No sequence number, no vector clock, no immutability enforcement (`openhands/app_server/event/event_service_base.py:190-197`, `openhands/app_server/event/filesystem_event_service.py:33-36`). Timestamp string is used for sorting but is not a version.
* Workspace archives: timestamp in key (`%Y%m%dT%H%M%SZ`) creates a new object per capture; a retry under `REQUIRED` re-uploads under a fresh ts, leaving orphaned prior objects to bucket lifecycle (`openhands/app_server/sandbox/workspace_archive.py:398-411`). No dedup or content hash. `base_commit` from `X-Archive-Base-Commit` header and manifest `created_at` are provenance, not versioning.
* Tags carry `agentprofileid`/`agentprofilerevision` and `archiveworkspacepath` but are just strings on the conversation row (`openhands/app_server/app_conversation/app_conversation_models.py:43-48`).

Effect: you can tell *which* conversation version (V1 vs legacy V0) but not *which revision* of its metadata/events you are reading.

### 4. Can artifacts be linked to the run that produced them?

**Yes, via `conversation_id` (the run ID) — but the linkage is scattered and requires multiple lookups.**

* Every artifact is keyed by `conversation_id.hex`: DB PK (`StoredConversationMetadata.conversation_id`), event directory (`v1_conversations/{hex}/`), archive key segment (`.../{sandbox_id}/{conversation_id}/{ts}.*`), export ZIP entries. This makes a full scan possible if you know the ID.
* User scoping: events live under `prefix/user_id/`; the service resolves `user_id` from `UserContext` or falls back to fetching `AppConversationInfo.created_by_user_id` (`openhands/app_server/event/event_service_base.py:66-83`). `get_conversation_path` enforces this prefix as the permission check.
* Sandbox linkage: `AppConversationInfo.sandbox_id` joins to `v1_remote_sandbox.id`; `count_conversations_by_sandbox_id` determines sharing under `max_num_conversations_per_sandbox` (`openhands/app_server/app_conversation/sql_app_conversation_info_service.py:389-395`). Parent/sub runs via `parent_conversation_id` and `get_sub_conversation_ids` (`openhands/app_server/app_conversation/sql_app_conversation_info_service.py:369-387`).
* Archive linkage: `sandbox_id` + `conversation_key` + `ts` + `ARCHIVE_WORKSPACE_PATH_TAG_KEY` (pinned at creation to avoid re-derivation drift) (`openhands/app_server/sandbox/workspace_archive.py:341-345`, `openhands/app_server/app_conversation/app_conversation_router.py:1015-1016`).
* What is missing: no single "run manifest" table that enumerates all object keys for a conversation; no foreign key from archive/manifest back into `conversation_metadata`; archive location is env-derived, not stored beside the row. Answering "every artifact for run X" requires: `get_app_conversation_info` + `search_events`/`count_events` + `count_conversations_by_sandbox_id` + object-store `list(prefix/...)` — not one API.

### 5. How are artifacts retired?

**Explicit delete only; no retention/TTL/GC.**

* **Trigger:** `DELETE /api/conversations/{id}` (or `AppConversationService.delete_app_conversation`) — synchronous DB delete of `conversation_metadata` + `app_conversation_start_task`, plus async `DELETE /api/conversations/{id}` to the agent-server if `sandbox_status==RUNNING` and not shared (`openhands/app_server/app_conversation/live_status_app_conversation_service.py:2690-2734`, `openhands/app_server/app_conversation/app_conversation_router.py:939-1027`).
* **Ordering / durability:** Workspace is archived *before* sandbox teardown via detached `_finalize_sandbox_delete`. If `RUNTIME_FILE_ARCHIVE_REQUIRED=true` and archive fails (5xx/422/429 or unconfirmed 401/404), deletion is blocked and the sandbox+runtime are kept for the idle-reap backstop to retry; 400 "nothing to archive" is allowed to proceed; unsupported format or missing bucket are loud config errors that still allow delete to avoid wedging (`openhands/app_server/sandbox/workspace_archive.py:342-360`, `openhands/app_server/sandbox/workspace_archive.py:538-542`, `openhands/app_server/app_conversation/app_conversation_router.py:898-931`).
* **Cascade:** `_delete_sub_conversations` enumerates `parent_conversation_id==parent` and deletes each sub-conversation's agent-server entry + DB row before the parent (`openhands/app_server/app_conversation/live_status_app_conversation_service.py:2744-2777`). `StoredConversationCostEvent` rows cascade via FK `ondelete='CASCADE'`; event blobs and workspace archives have no FK.
* **What is NOT retired:** Event JSON blobs remain in the file/object store after DB row deletion — no `delete` call on the event prefix. Archive objects are immutable and never deleted by the app-server (comment notes orphans are reaped only by bucket lifecycle policy: `openhands/app_server/sandbox/workspace_archive.py:408-411`). No cron/retention job, no `TTL` on DB rows, no `expires_at` column found. `delete_from_agent_server` swallows exceptions and still deletes the DB row, risking divergence.

## Architectural Decisions

* **Event storage separate from SQL** — Events live in a filesystem/object-store namespace (`{user_id}/v1_conversations/{hex}/`), not in SQL. This keeps the hot trajectory write path simple and lets deployments scale it independently to S3/GCS, at the cost of no transactional consistency between `conversation_metadata` and its events (`openhands/app_server/event/event_service_base.py:66-96`, `openhands/app_server/event/aws_event_service.py:27-76`).
* **Pluggable FileStore/EventService via `DiscriminatedUnionMixin` + env factory** — `LocalFileStore`, `S3FileStore`, `GoogleCloudFileStore` and `Filesystem/Aws/GoogleCloudEventService` share a common ABC and are selected from `StorageProvider`/`FILE_STORE_PATH`/`RUNTIME_FILE_ARCHIVE_*` env at startup (`openhands/app_server/file_store/files.py:8-42`, `openhands/app_server/config.py:307-328`). Enables OSS (local) to Cloud (GCS) without code change.
* **V1 sentinel + merge upsert for metadata** — Hard-filter `conversation_version=='V1'` and `db_session.merge()` simplify migration from V0 while preserving `created_at` across webhook-driven overwrites (`openhands/app_server/app_conversation/sql_app_conversation_info_service.py:458-466`, `openhands/app_server/app_conversation/sql_app_conversation_info_service.py:755-759`).
* **Detached archive-then-delete finalizer** — `_finalize_sandbox_delete` runs after the HTTP 204, archives workspace first, only then tears down the sandbox if unreferenced; `REQUIRED` gating + idle-reap backstop prevents data loss at the cost of leaked sandboxes on repeated failures (`openhands/app_server/app_conversation/app_conversation_router.py:885-937`).
* **Archive key includes conversation_id + timestamp** — Avoids sibling collision under `SandboxGroupingStrategy` (shared sandbox) and avoids dedup complexity; favors capture completeness over storage cost (`openhands/app_server/sandbox/workspace_archive.py:398-411`).
* **Streaming large artifacts, atomic small writes** — `LocalFileStore` uses temp+`os.replace`+`fsync` for safety; `_stream_to_tempfile` + `write_from_path` + `shutil.copyfile` keep multi-GB archives out of RAM (`openhands/app_server/sandbox/workspace_archive.py:176-204`, `openhands/app_server/file_store/local.py:26-57`).
* **Cost ledger separate from metadata row** — `StoredConversationCostEvent` records per-`usage_id` deltas with monotonic guards and stale-snapshot rejection, so the `accumulated_cost` column cannot regress even under concurrent stats events (`openhands/app_server/app_conversation/sql_app_conversation_info_service.py:609-678`).

## Notable Patterns

* **Atomic file write + streaming upload split** — Small JSON uses atomic temp rename; large binary archives use `write_from_path` streaming, with `InMemoryFileStore` explicitly rejected for archives (`openhands/app_server/file_store/local.py:26-57`, `openhands/app_server/file_store/memory.py:12-13`, `openhands/app_server/sandbox/workspace_archive.py:196-204`).
* **Tags-as-extension projection** — `AppConversationInfo.tags: dict[str,str]` carries open-ended metadata; typed `@computed_field` projections like `acp_server` and `launched_agent_profile` read tag keys (`acpserver`, `agentprofileid`) without a migration (`openhands/app_server/app_conversation/app_conversation_models.py:34-48`, `openhands/app_server/app_conversation/app_conversation_models.py:140-178`).
* **Lock-then-validate-then-stream** — Export acquires `try_acquire_redis_lock(key, ttl)` + periodic `reacquire()`, validates `count_events <= export_max_events`, then yields ZIP chunks incrementally via `_StreamingZipBuffer` (`openhands/app_server/app_conversation/live_status_app_conversation_service.py:2845-2962`, `openhands/app_server/utils/redis_lock.py:25-52`).
* **Graceful header fallback** — `_extract_repo_metadata` merges per-key fallback from prior format's headers so a second format with missing headers does not clobber repo identity (`openhands/app_server/sandbox/workspace_archive.py:45-58`).

## Tradeoffs

* **Durability vs storage cost for archives:** Default `RUNTIME_FILE_ARCHIVE_FORMAT='both'` captures lossy `git-delta` (compact) AND self-contained `tar.gz` (preserves gitignored files + full history). Comment explicitly defers narrowing to one format until cost is measured (infra#1444) — correctness favored over cost (`openhands/app_server/sandbox/workspace_archive.py:99-106`).
* **Single copy vs siblings:** By keying archives on both `sandbox_id` and `conversation_id`, sibling captures in the same second never overwrite each other, but every delete writes fresh timestamped objects even on retry — duplicates accumulate until bucket lifecycle reaps them (`openhands/app_server/sandbox/workspace_archive.py:398-411`).
* **Strong capture vs availability:** `RUNTIME_FILE_ARCHIVE_REQUIRED=true` blocks sandbox teardown on transient archive failures (5xx/429/422, 401/404 unconfirmed), keeping the running runtime alive for idle-reap retry. Prevents silent data loss, but a persistently failing archive can leak a running sandbox/VM until manual intervention (`openhands/app_server/sandbox/workspace_archive.py:538-542`).
* **RAM safety vs complexity:** Streaming to tempfile + `write_from_path` avoids OOM under concurrent large deletes, but adds disk I/O and cleanup paths; small-event path remains simple in-memory JSON (`openhands/app_server/sandbox/workspace_archive.py:176-204`).
* **No versioning vs simplicity:** Overwrite-on-save for metadata/events avoids needing a version store or conflict resolution, but loses audit/history and makes "what did this run produce at time T" unanswerable (`openhands/app_server/event/event_service_base.py:190-197`, `openhands/app_server/app_conversation/sql_app_conversation_info_service.py:440-503`).
* **Pluggability vs consistency:** Three event-store backends share one interface but have subtly different behavior (e.g., `FilesystemEventService._search_paths` uses `glob`, GCS uses `list_blobs`, S3 uses `list_objects_v2` with `ContinuationToken`). No shared integration test matrix for artifact lifecycle was found in the examined scope.

## Failure Modes / Edge Cases

* **Unconfirmed capture (401/404) blocks REQUIRED delete, but is indistinguishable from transient for the caller:** Returns `False` (keep sandbox) — correct for misrouted `archive_path`, but a persistent auth misconfig will hold the sandbox forever; no alert metric beyond a warning log (`openhands/app_server/sandbox/workspace_archive.py:444-462`, `openhands/app_server/sandbox/workspace_archive.py:538-542`).
* **Config error does NOT block REQUIRED delete:** Unsupported `RUNTIME_FILE_ARCHIVE_FORMAT` or missing `RUNTIME_FILE_ARCHIVE_BUCKET` logs `ERROR` and returns `True` (allow delete). Prevents wedging every delete, but silently loses durability when REQUIRED was intended (`openhands/app_server/sandbox/workspace_archive.py:369-394`).
* **Event overwrite + no atomicity with DB delete:** If `save_event` overwrites `same_id.json` concurrently with a `search_events`/`iter_events_for_export`, readers may see a torn JSON (`_load_event` logs and returns `None`), and `count_events` can disagree with actual loadable events. After `delete_app_conversation`, event blobs are orphaned — no cleanup (`openhands/app_server/event/event_service_base.py:23-31`, `openhands/app_server/app_conversation/live_status_app_conversation_service.py:2782-2831`).
* **Timestamp key collision:** `ts` is second-granular (`%Y%m%dT%H%M%SZ`). Before the `conversation_id` fix, sibling captures in the same second overwrote each other at the object level. Fix mitigates, but two rapid deletes of the *same* conversation could still collide (DB row no longer exists, but finalizer could still be in-flight for both tasks) (`openhands/app_server/sandbox/workspace_archive.py:398-411`).
* **Agent-server delete failure swallowed:** `_delete_from_agent_server` warns and continues to DB cleanup, so DB row and agent-server state can diverge; `Detached _finalize_sandbox_delete` may then skip `delete_sandbox` if count >0, leaving the agent-server runtime running (`openhands/app_server/app_conversation/live_status_app_conversation_service.py:2804-2810`, `openhands/app_server/app_conversation/app_conversation_router.py:918-920`).
* **OOM/connection-pool pressure:** Each archive builds a `FileStore` client lazily; GCS client is now process-wide cached (`_get_shared_storage_client` LRU) to avoid pool exhaustion, but S3 path still creates a new `boto3.client` per injection (`openhands/app_server/event/google_cloud_event_service.py:26-34`, `openhands/app_server/event/aws_event_service.py:118-121`).
* **Non-numeric env → silent fallback:** `_float_env` falls back to defaults on `ValueError`; a bad `RUNTIME_FILE_ARCHIVE_TIMEOUT='120s'` never archives under the intended timeout, with only a warning (`openhands/app_server/sandbox/workspace_archive.py:125-152`).

## Future Considerations

* Add a **run manifest / artifact inventory API**: `GET /conversations/{id}/artifacts` that returns DB row pointer + event count + manifest keys + archive object keys (joined from a small `conversation_artifacts` table populated at archive time), so "every artifact for run X" is one call.
* Introduce **monotonic artifact versioning** (row `revision` + event `seq` or content-hash naming) and make event files immutable (reject overwrite of existing `event_id`).
* Persist **archive pointer** beside the conversation row at creation/delete time (store full object key + manifest URL), rather than deriving it from env at read time.
* Implement **retention/GC**: `expires_at` on `conversation_metadata` + background reaper for orphaned event prefixes and failed-archive objects; wire bucket lifecycle policies into deployment docs/tests.
* Add **checksum + manifest verification** for archives (e.g., write SHA256 alongside each upload, validate on read).
* Unify **error taxonomy** between `REQUIRED` vs `best-effort` paths and emit metrics/counters for `unconfirmed_capture`, `retryable_failure`, `nothing_to_archive`.

## Questions / Gaps

* **Retention policy not codified:** No `TTL`, `retention_days`, or reaper job was found. Is idle-reap (runtime-api) the intended sole GC, and what reaps local `~/.openhands` event files? (Searched `openhands/app_server/`, `config.py`, `conversation_paths.py` — no evidence.)
* **Bulk delete / GDPR erasure:** No `delete by user_id` or `delete all artifacts for run` that cleans both SQL and object-store prefixes atomically — event blobs appear leaked after single-conversation delete.
* **Archive query path:** No router to *list* or *fetch* workspace archives by conversation after creation — can they be discovered without direct bucket listing? (Grepped `workspace_archive`, `archive` in `app_server` — only write path found.)
* **Event schema evolution:** `Event` is validated on load but no migration/version field for event payloads was found — are breaking schema changes handled?
* **Scale evidence:** No load test or lifecycle test covering >10k events per conversation (export limit is 10k, but storage/list behavior beyond that was not observed).

---

Generated by `Dimension 16.01: Artifact Lifecycle` against `openhands`.
