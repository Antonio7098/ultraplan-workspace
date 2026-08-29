# Source Analysis: temporal

## Artifact Lifecycle

### Source Info

| Field | Value |
|-------|-------|
| Name | temporal |
| Path | `studies/agent-harness-study/sources/temporal` |
| Language / Stack | Go 1.26 (server monorepo: `common/`, `service/history`, `service/worker`, `api/`, `proto/`, `chasm/`) |
| Analyzed | 2026-08-28 |

## Summary

Temporal has no user-facing concept called "artifact." The study-relevant durable outputs of a workflow **run** are its **history event batches** (the append-only `History` tree) and its **visibility record** (indexable execution metadata). Both start life inside the primary persistence store (Cassandra/SQL) as mutable state + history nodes, and after the workflow closes they optionally become **archived artifacts** in an external blob store before the primary copy is deleted at **retention**. The lifecycle is fully deterministic: `GenerateWorkflowCloseTasks` decides archive-vs-retention based on `archivalEnabled()`, the `archivalQueueTaskExecutor` performs the copy-out, and a subsequent `DeleteHistoryEventTask` (or a multi-stage `DeleteWorkflowExecution` via `shard.ContextImpl.DeleteWorkflowExecution`) retires the primary data. Naming, storage, and versioning are anchored in the archival subsystem (`common/archiver/*`) with three concrete backends (filestore, S3, GCS) and an extensible provider. The weak spots are observability (progress checkpointing only in GCS), verification (no checksum/hash on archived blobs), and the lack of a first-class "enumerate all artifacts for run X" API — retrieval requires knowing namespace/workflow/run + failover version.

## Rating

**7 / 10** — Clear model with tests, explicit interfaces, and operational safeguards. Schemas (`HistoryBlob`, `VisibilityRecord`, `ArchiveHistoryRequest`), storage abstractions (`HistoryArchiver`/`VisibilityArchiver` + `ArchiverProvider`), status-gated archival (cluster + namespace), per-version file naming, history-mutation detection, retry/non-retryable error categories, and a two-phase delete (visibility → mutable state → history branch) are all implemented and tested. Deductions: no retention for archived artifacts (archival store is write-once/read-many with no TTL/GC), GCS-only progress recovery for large histories, no artifact manifest or run-artifact index, and scavenger GC is best-effort/heuristic.

## Evidence Collected

Every entry MUST include a file path with line numbers. Format: `path/to/file.ts:NN`.

| Area | Evidence | File:Line |
|------|----------|-----------|
| Artifact schema — history blob | `HistoryBlobHeader` (namespace, namespace_id, workflow_id, run_id, is_last, first/last_failover_version, first/last_event_id, event_count) + `HistoryBlob { header, body: History[] }` | `proto/internal/temporal/server/api/archiver/v1/message.proto:13-29` |
| Artifact schema — history blob | Generated `HistoryBlobHeader` / `HistoryBlob` structs with typed getters | `api/archiver/v1/message.pb.go:30-196` |
| Artifact schema — visibility record | `VisibilityRecord { namespace_id, namespace, workflow_id, run_id, workflow_type_name, start_time, execution_time, close_time, status, history_length, memo, search_attributes map<string,string>, history_archival_uri, execution_duration }` | `proto/internal/temporal/server/api/archiver/v1/message.proto:31-47` |
| Artifact schema — visibility record | Generated `VisibilityRecord` struct with getters | `api/archiver/v1/message.pb.go:198-345` |
| Archive request — history | `ArchiveHistoryRequest { ShardID, NamespaceID, Namespace, WorkflowID, RunID, BranchToken, NextEventID, CloseFailoverVersion }` | `common/archiver/interface.go:16-25` |
| Archive request — visibility | `archival.Request` aggregating shard/namespace/workflow identity, `BranchToken`/`NextEventID`/`CloseFailoverVersion`, `HistoryURI`/`VisibilityURI`, workflow type, `StartTime/ExecutionTime/CloseTime/ExecutionDuration`, `Status`, `HistoryLength`, `Memo`, `SearchAttributes`, `Targets []Target`, `CallerService` | `service/history/archival/archiver.go:30-60` |
| History archiver interface | `HistoryArchiver { Archive(ctx, URI, *ArchiveHistoryRequest, ...ArchiveOption), Get(ctx, URI, *GetHistoryRequest), ValidateURI(URI) }` | `common/archiver/interface.go:44-59` |
| Visibility archiver interface | `VisibilityArchiver { Archive(ctx, URI, *VisibilityRecord), Query(ctx, URI, *QueryVisibilityRequest, NameTypeMap), ValidateURI(URI) }` | `common/archiver/interface.go:76-94` |
| Archival queue task | `ArchiveExecutionTask { WorkflowKey (NamespaceID/WorkflowID/RunID), VisibilityTimestamp, TaskID, Version }` — "archives both history+visibility then produces a retention timer task" | `service/history/tasks/archive_execution_task.go:12-21` |
| History iteration | `HistoryIterator { Next() *HistoryBlob, HasNext(), GetState() }` + `historyIteratorState { NextEventID, FinishedIteration }` with JSON state serialization | `common/archiver/history_iterator.go:23-44` |
| Size-bounded batching | `targetHistoryBlobSize = 2 * 1024 * 1024` (2 MB) used to split history into multiple `HistoryBlob`s | `common/archiver/filestore/history_archiver.go:44`, `common/archiver/gcloud/history_archiver.go:31`, `common/archiver/s3store/history_archiver.go:36` |
| Filestore history naming | `constructHistoryFilename` → `"<hash(ns)><hash(wf)><hash(run)>_<version>.history"`; prefix is `hash(ns)+hash(wf)+hash(run)` | `common/archiver/filestore/util.go:163-170` |
| Filestore visibility naming | `constructVisibilityFilename` → `"<closeTimestamp.UnixNano()>_<hash(runID)>.visibility"` | `common/archiver/filestore/util.go:172-174` |
| GCS multipart naming | `constructHistoryFilenameMultipart(... version, part)` — history split into per-part objects | `common/archiver/gcloud/history_archiver.go:165,231` |
| Archival provider | `ArchiverProvider { GetHistoryArchiver(scheme), GetVisibilityArchiver(scheme) }` with per-scheme singleton cache, `CustomHistoryArchiverFactory`/`CustomVisibilityArchiverFactory` hooks, and `filestore → gcloud → s3store` switch | `common/archiver/provider/provider.go:27-247` |
| Storage backends | `filestore/`, `s3store/`, `gcloud/` each implement `HistoryArchiver` + `VisibilityArchiver` | `common/archiver/filestore/history_archiver.go:1-304`, `common/archiver/s3store/history_archiver.go:1-386`, `common/archiver/gcloud/history_archiver.go:1-388` |
| Close-task decision | `GenerateWorkflowCloseTasks` branches on `archivalEnabled()` (cluster + namespace checks): archival → `ArchiveExecutionTask` at `closedTime + min(jittered ArchiveDelay, retention)`; else → `DeleteHistoryEventTask` at `closeTime + retention + jitter` | `service/history/workflow/task_generator.go:196-275` |
| Retention + jitter | `getRetention()` from namespace registry (default `1 * 24h` if namespace missing) + `RetentionTimerJitterDuration` and `ArchivalProcessorArchiveDelay` dynamic configs | `service/history/workflow/task_generator.go:282-295`, `service/history/workflow/task_generator.go:341-365` |
| Archive execution | `archivalQueueTaskExecutor.processArchiveExecutionTask` loads mutable state, builds `archival.Request` (parsing `VisibilityURI`/`HistoryURI` per-namespace), calls `archiver.Archive` for enabled targets, then enqueues `DeleteHistoryEventTask` | `service/history/archival_queue_task_executor.go:92-256` |
| Archival fan-out | `archiver.Archive` rate-limits, then launches `archiveHistory` + `archiveVisibility` in parallel (`sync.WaitGroup` + `multierr.Combine`) | `service/history/archival/archiver.go:107-174` |
| History mutation guard | `historyMutated(req, batches, isLast)` rejects if `lastFailoverVersion > CloseFailoverVersion` or, when `isLast`, `lastEventID+1 != NextEventID` | `common/archiver/filestore/util.go:211-223`, `common/archiver/gcloud/history_archiver.go:301-314` |
| URI validation | `ValidateURI` checks scheme (`file`, `s3`, `gs`) and path/bucket existence; filestore checks `os.Stat` dir, GCS/S3 check bucket reachability | `common/archiver/filestore/history_archiver.go:276-282`, `common/archiver/gcloud/history_archiver.go:279-298`, `common/archiver/s3store/history_archiver.go:40-42` |
| Pagination tokens | `getHistoryToken { CloseFailoverVersion, NextBatchIdx }` JSON-serialized as `NextPageToken` for `Get` | `common/archiver/filestore/history_archiver.go:64-67`, `common/archiver/filestore/util.go:142-152` |
| GCS progress checkpoint | `loadHistoryIterator`/`saveHistoryIteratorState` via `ArchiveFeatureCatalog.ProgressManager` allow resuming a large multi-part upload after crash (GCS only) | `common/archiver/gcloud/history_archiver.go:355-388` |
| Primary retention deletion | `GenerateDeleteHistoryEventTask` creates `DeleteHistoryEventTask{ VisibilityTimestamp: closeTime+retention+jitter, BranchToken, Version, ArchetypeID }` | `service/history/workflow/task_generator.go:338-367` |
| 4-stage hard delete | `ContextImpl.DeleteWorkflowExecution` stages: (1) visibility delete task + replication delete task, (2) delete current execution pointer, (3) delete mutable state, (4) delete history branch; stages tracked in `DeleteWorkflowExecutionStage` bitmask for idempotent retry | `service/history/shard/context_impl.go:910-1019` |
| Retention delete skip | `DeleteWorkflowExecutionByRetention` pre-marks `DeleteWorkflowExecutionStageReplication` as processed ("both clusters have independent retention timers") | `service/history/deletemanager/delete_manager.go:139-161` |
| History branch deletion | `ExecutionManager.DeleteHistoryBranch` decomposes branch token into `InternalDeleteHistoryBranchRange[]` and calls `persistence.DeleteHistoryBranch` | `common/persistence/history_manager.go:130-210` |
| Persistence layer | `ExecutionStore.DeleteHistoryBranch(ctx, *InternalDeleteHistoryBranchRequest)` is the SPIs-level delete (cassandra `DELETE` / SQL `DELETE FROM history_node/history_tree`) | `common/persistence/persistence_interface.go:159-161`, `common/persistence/cassandra/history_store.go:267-283`, `common/persistence/sql/history_store.go:333-347` |
| Scavenger GC | `Scavenger.Run` paginates `GetAllHistoryTreeBranches`, filters by `historyDataMinAge`, calls `DescribeMutableState` — if NotFound, calls `DeleteHistoryBranch`; if completed and older than `retention + executionDataDurationBuffer`, calls `AdminService.DeleteWorkflowExecution` | `service/worker/scanner/history/scavenger.go:118-379` |
| Archival metadata | `ArchivalMetadata { GetHistoryConfig(), GetVisibilityConfig() }` + `ArchivalConfig { ClusterConfiguredForArchival(), GetClusterState(), ReadEnabled(), GetNamespaceDefaultState(), GetNamespaceDefaultURI() }` with static+dynamic state | `common/archiver/archival_metadata.go:14-210` |
| Archival config validation | `config.Archival.Validate` requires (`enabled && URISet && provider`) XOR (`disabled && !read && !URISet && !provider`) for each of history/visibility | `common/config/archival.go:16-42` |

## Answers to Dimension Questions

### 1. What types of artifacts exist?

Three durable categories (all keyed by namespace/workflow/run):

**Primary persistence artifacts (pre-retention):**
- **History event batches** (`History` protos, `proto/internal/.../message.proto` history) — the authoritative append-only log. Stored as `DataBlob` rows in `history_node` / `history_tree` keyed by `BranchToken` + `NodeID`/`TransactionID`. A logical "artifact set" is the full chain of batches from `FirstEventID` to `NextEventID-1`.
- **Mutable state** (`WorkflowMutableState` in `common/persistence/data_interfaces.go`) — the current execution snapshot (execution info/state, timers, activities, etc.). Deleted atomically with the current-execution pointer during stage 3 of `DeleteWorkflowExecution`.
- **Visibility record** — the index row in the visibility store (ElasticSearch or SQL secondary visibility DB). Produced on start (`StartExecutionVisibilityTask`) and updated on upsert/close.

**Archived artifacts (post-close, pre-primary-deletion):**
- `HistoryBlob` (`api/archiver/v1/message.pb.go:146-152`) — `Header` (identity + `IsLast` + failover versions + event counts) + `Body` (one or more `History` batches, JSONPB-encoded). On GCS/S3 these are sharded into per-part objects; on filestore they are a single `.history` file.
- `VisibilityRecord` (`api/archiver/v1/message.pb.go:198-215`) — snapshot of the closed execution's indexed metadata (`search_attributes` already stringified via `searchattribute.Stringify`) plus `Memo`, `HistoryLength`, and `HistoryArchivalUri` to cross-link to the history artifact.

**Meta-artifacts:**
- `HistoryBranchDetail` (from `GetAllHistoryTreeBranches`) — branch metadata used by the scavenger to enumerate garbage branches.
- `Chasm` trees and transition histories for newer archetypes (visible via `chasm.ArchetypeID` threading through delete paths), but storage is the same history-node mechanism.

No first-class user-uploaded file/blob artifact type exists. If the user needs arbitrary file storage, it must be stored externally and referenced via payload/memo/search attributes — the server does not manage it.

### 2. How are artifacts named and stored?

**Primary store:** location is implicit from `BranchToken` (opaque bytes encoding shard, tree, and branch IDs; see `history_branch_util.go`) plus `ShardID`. Cassandra partitions by shard + branch; SQL uses `shard_id + tree_id + branch_id`. There is no human-readable name — the "address" is the branch token.

**Archived history (filestore):**
- Directory: `URI.Path()` from the per-namespace `HistoryArchivalState.URI` (e.g., `file:///var/temporal/archival`).
- Filename: `constructHistoryFilename` → `hash(namespaceID)+hash(workflowID)+hash(runID) + "_" + closeFailoverVersion + ".history"` (`common/archiver/filestore/util.go:163-170`). `hash` is `farm.Fingerprint64` decimal. All history for a given failover version collapses into one file (all `HistoryBlob` bodies concatenated, JSONPB-encoded).
- Visibility: `constructVisibilityFilename` → `"<closeTime.UnixNano()>_<hash(runID)>.visibility"` (`common/archiver/filestore/util.go:172-174`) — lexicographically ordered by close time for scan/pagination.

**Archived history (S3):**
- `s3://<bucket>/<path>/<hash-prefix>_<version>.history` (single file) or per-part variants; bucket/key validation includes region (`common/archiver/s3store/history_archiver.go:30-42`).
- Encoder: `codec.NewJSONPBEncoder().EncodeHistories / DecodeHistories`.

**Archived history (GCS):**
- Multipart: `constructHistoryFilenameMultipart(ns, wf, run, version, part)` for each 2 MB chunk. This is the only backend that supports resumable checkpointing of a multi-part history upload via `ProgressManager` (`common/archiver/gcloud/history_archiver.go:165-178`).

**Archived visibility (all backends):**
- Serialized with the same `codec.JSONPBEncoder` as history. Filestore stores one file per execution; S3/GCS store per-record objects queryable via `Query` (which translates `QueryVisibilityRequest.Query` string into time-range prefix scans and search-attribute filters — exact parser is backend-specific).

**Storage configuration:** `config.HistoryArchiverProvider` / `config.VisibilityArchiverProvider` hold `Filestore`, `Gstorage`, `S3store`, and `CustomStores` entries. The active store per namespace is given by `NamespaceConfig.ArchivalState.URI` (a full URI like `s3://...` / `gs://...` / `file:///...`), and the provider dispatches on URI scheme (`common/archiver/provider/provider.go:149-171`).

### 3. Are artifacts versioned?

Yes, at two distinct layers:

**Workflow failover version (the primary archival version identifier).**
- Every close carries `CloseFailoverVersion` (a monotonically increasing int64 from `GetCloseVersion()` / the version history). Archive file names embed it (`_<version>.history`), and `GetHistoryRequest` / the scavenger can address a specific version (`*int64 CloseFailoverVersion` in `common/archiver/interface.go:32`). `Get` falls back to the **highest** version when none is requested (`common/archiver/filestore/history_archiver.go:284-303`, `common/archiver/gcloud/history_archiver.go:316-353`), preserving visibility into cross-cluster failover histories.
- `HistoryBlobHeader` records both `FirstFailoverVersion` and `LastFailoverVersion` per blob plus `FirstEventID`/`LastEventID`/`EventCount`, so a single logical history can span failover versions inside one archive.

**History-branch/event-ID versioning.**
- The internal history tree has branches (fork on reset) encoded in `HistoryBranch` ancestors; each node is ordered by `NodeID` (= `EventID`). The history-iterator integrity check (`historyMutated`) asserts that the iterated batches match the expected `NextEventID` and do not cross a failover boundary (`common/archiver/filestore/util.go:211-223`), preventing silent corruption.
- `HistoryIteratorState { NextEventID, FinishedIteration }` is itself versioned as an opaque JSON blob for GCS resume.

**Data-format versioning:** `codec.JSONPBEncoder` produces proto-JSON; there is no explicit schema-version field on the archive file itself. Compatibility is delegated to proto field-presence rules. This is a minor gap for forward-incompatible archival schema changes.

### 4. Can artifacts be linked to the run that produced them?

Yes — tightly for primary state and moderately for archived artifacts.

**Primary state:** the linkage is the `WorkflowKey { NamespaceID, WorkflowID, RunID }` plus `BranchToken` and `ShardID`. Every persistence call (`GetWorkflowExecution`, `ReadHistoryBranch`, `DeleteHistoryBranch`) takes this triple; the history branch token itself encodes the workflow identity. `TaskGenerator` threads the same `WorkflowKey` into all archival/retention tasks.

**Archived artifacts:**
- `HistoryBlobHeader` carries `Namespace`, `NamespaceId`, `WorkflowId`, `RunId` explicitly (`proto/internal/temporal/server/api/archiver/v1/message.proto:13-24`).
- `VisibilityRecord` carries the same four identity fields plus `HistoryArchivalUri` (`proto/internal/temporal/server/api/archiver/v1/message.proto:31-46` / `service/history/archival/archiver.go:231-246`), providing an explicit cross-link from the visibility artifact back to where the history artifact lives. `HistoryLength` and `ExecutionDuration` are also snapshotted so the visibility record is self-contained.
- `ArchiveHistoryRequest` and `archival.Request` both carry the full `WorkflowKey` plus `CloseFailoverVersion` and `NextEventID` for disambiguation across failovers and retries.

**Enumeration gap:** there is no `ListArtifactsForRun(RunID)` or artifact manifest. To find all artifacts for a run you must know the namespace and reconstruct the filename/prefix (`hash(ns)+hash(wf)+hash(run)`), or query visibility with `WorkflowId = X AND RunId = Y`. The scavenger can enumerate all branches (`GetAllHistoryTreeBranches`) and parse `SplitHistoryGarbageCleanupInfo` to recover `(namespaceID, workflowID, runID)` from `branch.Info`, but this is an administrative path, not a user API.

### 5. How are artifacts retired?

There are two retirement planes, one for primary data and one for archived data.

**Primary data — normal path (archival enabled):**
1. On `GenerateWorkflowCloseTasks` (`service/history/workflow/task_generator.go:196-269`), if `archivalEnabled()`, an `ArchiveExecutionTask` is scheduled for `closeTime + min(ArchivalProcessorArchiveDelay jitter, retention)`.
2. `archivalQueueTaskExecutor.processArchiveExecutionTask` copies history+visibility to the blob store (`service/history/archival_queue_task_executor.go:92-109`), then immediately enqueues a `DeleteHistoryEventTask` via `GenerateDeleteHistoryEventTask`.
3. The timer queue fires `DeleteHistoryEventTask` at `closeTime + retention + RetentionTimerJitterDuration`, which calls `DeleteWorkflowExecutionByRetention` → `ContextImpl.DeleteWorkflowExecution` with `retentionDelete=true` (skips replication stage). The call is idempotent and retried via the `DeleteWorkflowExecutionStage` bitmask; the branch is only deleted in stage 4 after earlier stages succeed (`service/history/shard/context_impl.go:910-936`).

**Primary data — no archival:**
- `GenerateWorkflowCloseTasks` directly enqueues `DeleteHistoryEventTask` at close+retention; archival steps are skipped entirely.

**Primary data — explicit deletes:**
- `AdminService.DeleteWorkflowExecution` / `HistoryService.DeleteWorkflowVisibilityRecord` allow operator-initiated deletes that bypass retention, going through the same 4-stage path with `retentionDelete=false` (so a replication tombstone is emitted).
- The `DeleteManager` handles immediate-delete (`DeleteExecutionTask`) and retention-delete paths separately (`service/history/deletemanager/delete_manager.go:80-216`).

**History branch garbage:**
- `Scavenger.Run` (`service/worker/scanner/history/scavenger.go:118-289`) is a safety net: it scans all branches older than `historyDataMinAge` (default 30 days) and, if `DescribeMutableState` returns `NotFound`, deletes the branch via `DeleteHistoryBranch`. If the mutable state still exists and is completed and older than `retention + executionDataDurationBuffer`, it force-deletes via `AdminService` (optional behind `enableRetentionVerification` dynamic config).

**Archived artifacts:** **no retirement policy.** Once written, `*.history` / `*.visibility` objects in filestore/S3/GCS are never deleted, TTL-expired, or compacted by Temporal. The provider does not expose a `Delete` or `Expire` RPC; the file/object simply persists until an external operator removes it. `writeFile` does `os.Remove` + `os.Create` before overwriting, but this is atomic replacement for re-archival, not lifecycle GC.

## Architectural Decisions

| Decision | Evidence | Consequence |
|----------|----------|-------------|
| Separate `HistoryArchiver` / `VisibilityArchiver` interfaces dispatched by URI scheme | `common/archiver/interface.go:44-94`, `common/archiver/provider/provider.go:123-246` | History (large, ordered, blob) and visibility (small, indexed, queryable) can live in different stores; namespaces pick per-target backends via URI. Adds operational complexity vs. unified store. |
| Hash-based deterministic filenames (`farm.Fingerprint64`) | `common/archiver/filestore/util.go:176-178` | No collisions from special characters; opaque — operators cannot locate an execution's archive without reconstructing the hash or querying visibility. |
| 2 MB history chunk size + history iterator | `common/archiver/history_iterator.go:20,43-44`, `common/archiver/filestore/history_archiver.go:44` | Large histories are sharded (especially GCS multipart). Bounds memory per blob but multiplies objects per workflow for S3/GCS. |
| History-mutation guard before upload | `common/archiver/filestore/util.go:211-223`, `common/archiver/gcloud/history_archiver.go:301-314` | Detects races where a second close or replication changes `NextEventID`/`CloseFailoverVersion` between task scheduling and archival; returns `ErrHistoryMutated` (non-retryable) instead of archiving stale data. |
| Archival-then-delete ordering (archive queue → deletion timer) | `service/history/archival_queue_task_executor.go:92-256`, `service/history/workflow/task_generator.go:247-268` | Ensures primary data is not deleted until after successful copy-out; eliminates data-loss window if archiver is slow but doubles queue latency. |
| 4-stage idempotent delete with stage bitmask | `service/history/shard/context_impl.go:922-935` | Retries after stage 1/2 are safe; failure after stage 3 (mutable state gone) leaves an orphan history branch that the scavenger must reap — a deliberate consistency tradeoff (availability over atomicity). |
| Per-namespace archival state (`ARCHIVAL_STATE_ENABLED/DISABLED`) + cluster config double-gate | `common/archiver/archival_metadata.go:54-61`, `service/history/workflow/task_generator.go:934-940` | Archival cannot fire unless both cluster and namespace are enabled; allows incremental rollout. Misconfiguration silently falls through to retention-only path. |
| Retention jitter (`FullJitter` over `RetentionTimerJitterDuration` and `ArchivalProcessorArchiveDelay`) | `service/history/workflow/task_generator.go:255-257,356-357` | Prevents thundering-herd on timer queue when many workflows share the same `Retention`; slightly nondeterministic deletion time. |
| Scavenger as best-effort orphan GC | `service/worker/scanner/history/scavenger.go:210-289` | Catches history branches orphaned by stage-3 delete failures; paced by RPS limiter and `historyDataMinAge`. Not transactional and not enabled to force-delete by default. |
| GCS-only progress persistence (`ArchiveFeatureCatalog.ProgressManager`) | `common/archiver/gcloud/history_archiver.go:355-388` | Only GCS can resume a crashed multi-part history upload; S3/filestore re-upload from scratch on retry (wasted work for very large histories). |

## Notable Patterns

- **Provider + custom factory extensibility** (`common/archiver/provider/provider.go:34-65`) — operators can inject `CustomHistoryArchiverFactory` / `CustomVisibilityArchiverFactory` to support a new scheme (e.g., `azure://`) without forking the server. Hook is exercised in `temporal/fx_test.go:383-424` and `service/history/archival/archiver_test.go`.
- **Codec round-trip through JSONPB** — history batches are `EncodeHistories`/`DecodeHistories` (`codec.NewJSONPBEncoder`) so that the archived file is human-readable JSON. Filestore `readFile`/`decode` is guarded by `#nosec` (`common/archiver/filestore/util.go:82-85`) since the path comes from the internal URI, not user input.
- **Duplicate-archival short-circuit** — both filestore and GCS check existence before upload (`fileExists` / `Exist(ctx, URI, filename)`) and skip or overwrite; `Archive` returns `nil` on `serviceerror.NotFound` from the history iterator (concurrent double-archive) (`common/archiver/filestore/history_archiver.go:138-145`, `common/archiver/gcloud/history_archiver.go:166-173`).
- **Stage-aware replication suppression** — `DeleteWorkflowExecutionByRetention` marks `DeleteWorkflowExecutionStageReplication` as already processed so retention deletes do not emit cross-cluster tombstones (`service/history/deletemanager/delete_manager.go:149`).
- **Namespace-not-found fallback** — `getRetention` falls back to `defaultWorkflowRetention = 1 * 24h` (and visibility delete is skipped) when namespace lookup fails; archival URI parse failures emit `ArchivalTaskInvalidURI` and return non-retryable errors that still block the queue forever unless the URI is fixed (`service/history/archival_queue_task_executor.go:156-168`).

## Tradeoffs

| Tradeoff | Pro | Con |
|----------|-----|-----|
| **Archival vs. deletion latency** (archive scheduled at `close + min(ArchiveDelay, retention)`; delete at `close+retention`) | Minimizes window where data is neither archived nor primary; avoids retention expiry before archival copy | If archival backend is down past retention, delete still fires and data is lost; no "hold until archived" invariant. The archiver being down stalls only the archival queue, not the timer queue. |
| **Hash-based naming** | Safe for any `WorkflowID` characters (including `/`, Unicode) | Opaque to humans; correlating a run to its archive object requires either hashing or a visibility lookup. Prefix enumeration for a single run requires recomputing 3 hashes. |
| **Per-version files vs. single logical artifact** | Preserves failover history faithfully; `Get` can address a specific version | Each failover creates a new object; there is no manifest stitching versions together. Clients that don't know `CloseFailoverVersion` must list + pick max version. |
| **JSONPB readability vs. density** | Easy to debug/inspect; no external schema registry | ~2-3x larger than proto-binary; compression is delegated to the blob store. No checksum is stored, so bit-rot is undetectable. |
| **GCS resume vs. S3/filestore re-upload** | Crash-safe for very large histories (≥10s of MB) on GCS | S3/filestore waste history-scan work on every retry after a non-transient failure before the upload phase. |
| **Scavenger min-age gating** | Avoids racing with in-flight close/archival writes | Orphan branches can linger up to `historyDataMinAge` (default 30 days) before collection. |
| **4-stage delete non-atomicity** | Each stage is retriable; partial failures are observable via `DeleteWorkflowExecutionStage` | Orphan history branches and invisible-but-not-deleted mutable states exist transiently; the system is eventually consistent, not atomic. |
| **No archive GC** | Simplicity; archived data is by definition post-retention and the operator "owns" it | Storage grows without bound; lifecycle/retention of archives must be managed externally (bucket lifecycle rules, `find -delete`, etc.). |

## Failure Modes / Edge Cases

| Mode | Behavior | Mitigant / Gap |
|------|----------|----------------|
| **Invalid archival URI** | `NewURI` parse failure or scheme mismatch yields `ErrURISchemeMismatch`/`ErrInvalidURI`; archival queue emits `ArchivalTaskInvalidURI` metric and returns a non-retryable error that is **not** retried — the task blocks the archival queue forever until the namespace URI is fixed. | Documented in `service/history/archival_queue_task_executor.go:157-168` but operationally severe. |
| **Concurrent double-archive** | Two `ArchiveExecutionTask`s for the same execution (conflict resolution) race through the history iterator; the second sees `NotFound` or `HistoryMutated` and exits early. | `HistoryIterator.Next` returns `serviceerror.NotFound` → `ArchiveSkipped` (`common/archiver/filestore/history_archiver.go:139-145`, `common/archiver/gcloud/history_archiver.go:136-143`). |
| **History mutated between scheduling and archive** | `historyMutated` compares `CloseFailoverVersion`/`NextEventID` against blob header. | Returns `ErrHistoryMutated` (non-retryable) instead of silently archiving wrong data (`common/archiver/filestore/util.go:211-223`). |
| **Archival backend transient failure** | Transient error (`IsPersistenceTransientError`) → retry; non-transient → wrapped via `NonRetryableError` option or returned as `errUploadNonRetryable` (GCS). | Retry count/limit is governed by the queue's retry policy, not the archiver. Queue retries "forever" for transient errors. |
| **Delete after stage 3 failure** | If delete crashes after deleting mutable state but before deleting history branch, the branch becomes an orphan (no mutable state, so `DescribeMutableState` returns `NotFound`). | Scavenger GC deletes it after `historyDataMinAge` (`service/worker/scanner/history/scavenger.go:271-288`). Without the scavenger (or with `historyDataMinAge` large), the leak persists. |
| **Namespace deleted before archival/retention fires** | `GetNamespaceByID` → `NamespaceNotFound` → `getRetention` uses `defaultWorkflowRetention`; visibility delete is skipped; archival URI parse uses stale `NamespaceEntry` snapshot. | Defensive but may archive to a stale URI or retain under wrong duration. Covered in `service/history/shard/context_impl.go:943-949`, `service/history/workflow/task_generator.go:285-292`. |
| **S3/GCS partial multipart upload on crash** | GCS persists iterator progress (`progress.CurrentPageNumber + IteratorState`) and can resume; S3/filestore cannot — a crash mid-upload after some parts were uploaded leaks those objects until external cleanup. | Only GCS path (`common/archiver/gcloud/history_archiver.go:355-388`) has resume; S3 has no `LoadProgress` equivalent. |
| **Large history (≥10 MB) before archival** | History iterator batches at 2 MB; each batch is a separate RPC to the execution manager (`persistence.ReadHistoryBranchRequest` with `historyPageSize=250`, `common/archiver/history_iterator.go:20,195-204`). | Sustained churn on the persistence store; latency scales linearly with history size. GCS's check `Exist` per part adds per-batch overhead. |
| **Retention < ArchivalDelay** | `archiveTime = closedTime + min(ArchiveDelay jitter, retention)` ensures archival is before deletion, but if both are very short, archival may still race with delete. | Jitter helps but does not guarantee ordering under heavy backlog. |

## Future Considerations

- **Archive retention / lifecycle** — add an archive-level TTL or delegation to native bucket lifecycle (e.g., `S3LifecycleRule` emitted during `NewHistoryArchiver` provisioning), plus a server-owned `DeleteArchivedHistory` RPC behind an admin guard to allow operator-enforced GC. Currently the archive plane is append-only with no server-driven retirement.
- **Run→artifact manifest** — maintain an index record (or a manifest file written alongside `*.history`) listing all archive object keys for a `RunID`. This would make "enumerate every artifact produced by run X" a single `GetManifest` call rather than prefix-list + hash reconstruction.
- **Integrity envelope** — include a checksum (e.g., `sha256` of the `EncodeHistories` bytes) in `HistoryBlobHeader` and verify on `Get`; persist the checksum in `VisibilityRecord` as well. This closes the bit-rot/corruption detection gap and helps the mutation guard.
- **Uniform progress checkpointing** — promote GCS's `ProgressManager` pattern to S3 and filestore (or to the queue itself via `ArchiveFeatureCatalog`) so crashed multi-part archival re-uploads are bounded.
- **Stronger archival→delete invariant** — consider making retention delete conditional on successful archival (e.g., a marker written to mutable state after archive, checked in the delete task). Today delete fires on wall time regardless of archival success.
- **Structured archival errors + DLQ** — the `ArchivalTaskInvalidURI` metric is the only signal for a poisoned archival task that will never succeed. Queueing such tasks to a visible DLQ instead of silently looping would improve operability.
- **Proto-binary + optional JSONPB** — allow the archiver codec to be pluggable (or at least offer binary mode) to halve storage costs for high-volume namespaces.
- **Schema-version on archive files** — add a top-level `archive_format_version` field to the serialized `HistoryBlob`/`VisibilityRecord` envelope so that rolling upgrades or new fields can be detected and migrated.

## Questions / Gaps

- **No evidence of artifact checksum/hash or encryption-at-rest handling** — searches across `common/archiver/`, `common/codec/`, and `service/history/archival/` returned no hash computation or envelope-encryption code; `codec.JSONPBEncoder` is plaintext JSON. Confirm whether blob stores are expected to provide this externally.
- **No search over archived artifacts by artifact content** — `VisibilityArchiver.Query` accepts a `Query` string but the filing convention for history blobs is by hashed identity, not by time-range secondary index; full-history content search is not supported.
- **No explicit test for "delete fires before archival finishes" race** — `service/history/archival_queue_task_executor_test.go` and `service/history/workflow/task_generator_test.go` cover close→archive→delete ordering under mocks, but no integration test asserts that an archival transient failure long enough to outlive `retention` does NOT cause data loss. The observed implementation suggests it would.
- **Scavenger's `enableRetentionVerification` default** — dynamic config key not examined in this study; whether production deployments enable forced deletion of completed executions past retention via the scavenger is unverified (file search returned only the wiring in `service/worker/scanner/history/scavenger.go:267-379`).

---

Generated by `Dimension 16.01: Artifact Lifecycle` against `temporal`.
