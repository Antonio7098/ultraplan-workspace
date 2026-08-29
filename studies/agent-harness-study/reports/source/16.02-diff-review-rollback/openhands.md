# Source Analysis: openhands

## 16.02 Diff, Review, and Rollback

### Source Info

| Field | Value |
|-------|-------|
| Name | openhands |
| Path | `studies/agent-harness-study/sources/openhands` |
| Language / Stack | Python (FastAPI, SQLAlchemy) + TypeScript/React frontend; SDK at `_sdk_inspect/` |
| Analyzed | 2026-08-28 |

## Summary

OpenHands treats the agent's workspace as a git repository and provides a mature, file-level diff/comparison stack, but it has **no artifact-centric review, approval, comment/annotation, or rollback subsystem**. Diff is implemented twice: locally via `get_git_diff`/`get_changes_in_repo` (`_sdk_inspect/sdk/git/git_diff.py:53`, `_sdk_inspect/sdk/git/git_changes.py:37`) and remotely via the in-pod agent-server (`GET /api/git/diff`, `GET /api/git/changes`) proxied through the app-server (`openhands/app_server/app_conversation/app_conversation_router.py:1224-1366`) and rendered in the frontend Changes tab (`frontend/src/routes/changes-tab.tsx:23`, `frontend/src/components/features/diff-viewer/file-diff-viewer.tsx:75`). The only durable artifact trail is the conversation-scoped event log (`openhands/app_server/event/event_service_base.py:190`) and a best-effort workspace archive (`openhands/app_server/sandbox/workspace_archive.py:334`) captured at delete/reap time. There is no versioned artifact store, no approval gate, no per-artifact comment model, and no revert/rollback API—undoing a bad change is manual `git` inside the sandbox, with no audit guarantee.

## Rating

**3/10 — Absent, implicit, ad-hoc, or unsafe**

Diff/comparison is clear and tested (7–8 in isolation), but the dimension as a whole fails: artifact review/approval, typed comments/annotations, rollback/revert, and auditable artifact↔run linking are absent. The workspace archive is a durability backstop, not a versioned rollback system, and is a no-op unless `RUNTIME_FILE_ARCHIVE_ENABLED` is set (`openhands/app_server/sandbox/workspace_archive.py:77`). A bad artifact change cannot be reverted with full audit trail via product APIs.

## Evidence Collected

Every entry MUST include a file path with line numbers. Format: `path/to/file.ts:NN`.

| Area | Evidence | File:Line |
|------|----------|-----------|
| Diff model | `GitDiff` (original/modified strings) and `GitChange`/`GitChangeStatus` enums | `sources/openhands/_sdk_inspect/sdk/git/models.py:25`, `sources/openhands/_sdk_inspect/sdk/git/models.py:9-18` |
| Diff generator (local) | `get_git_diff(relative_file_path, ref)` — resolves repo, checks size (1 MiB), fetches `original` via `git show {ref}:{path}` and `modified` from disk | `sources/openhands/_sdk_inspect/sdk/git/git_diff.py:53-128` |
| Diff generator (local) | `get_closest_git_repo()` walks to `.git` root | `sources/openhands/_sdk_inspect/sdk/git/git_diff.py:29-50` |
| Changes generator (local) | `get_changes_in_repo(repo_dir, ref)` — `git diff --name-status {ref}` + `ls-files --others`, handles R/C renames as DELETED+ADDED | `sources/openhands/_sdk_inspect/sdk/git/git_changes.py:37-202` |
| Changes generator (nested) | `get_git_changes(cwd, ref)` aggregates top-level + nested `.git` repos | `sources/openhands/_sdk_inspect/sdk/git/git_changes.py:205-246` |
| Ref resolution | `get_valid_ref(repo_dir, override)` — tries `origin/{branch}`, `origin/{default}`, merge-base, fallback `GIT_EMPTY_TREE_HASH`; explicit `override` resolved via `rev-parse --verify` | `sources/openhands/_sdk_inspect/sdk/git/utils.py:104-206` |
| Git utils | `run_git_command` safe subprocess, `validate_git_repository`, empty-tree constant | `sources/openhands/_sdk_inspect/sdk/git/utils.py:17-36`, `sources/openhands/_sdk_inspect/sdk/git/utils.py:209-238` |
| SDK local workspace | `LocalWorkspace.git_diff()` / `git_changes()` delegate to local generators | `sources/openhands/_sdk_inspect/sdk/workspace/local.py:176-189`, `sources/openhands/_sdk_inspect/sdk/workspace/local.py:161-174` |
| SDK remote workspace | `_git_diff_generator` / `_git_changes_generator` — generator pattern yielding `GET /api/git/diff?path=` / `/api/git/changes?path=` | `sources/openhands/_sdk_inspect/sdk/workspace/remote/remote_workspace_mixin.py:346-371`, `sources/openhands/_sdk_inspect/sdk/workspace/remote/remote_workspace_mixin.py:318-344` |
| SDK workspace base interface | Abstract `git_diff`/`git_changes` declarations | `sources/openhands/_sdk_inspect/sdk/workspace/base.py:149-159` |
| App-server proxy: shared helper | `_proxy_git_runtime_call()` — resolves `AgentServerContext`, proxies `GET {agent_server_url}{runtime_path}?path=&ref=` with `X-Session-API-Key`, maps 409/502 | `sources/openhands/openhands/app_server/app_conversation/app_conversation_router.py:1224-1306` |
| App-server proxy: changes | `GET /{conversation_id}/git/changes?path=&ref=` → `/api/git/changes` | `sources/openhands/openhands/app_server/app_conversation/app_conversation_router.py:1308-1339` |
| App-server proxy: diff | `GET /{conversation_id}/git/diff?path=&ref=&ref` → `/api/git/diff` | `sources/openhands/openhands/app_server/app_conversation/app_conversation_router.py:1342-1366` |
| App-server proxy: context resolution | `_get_agent_server_context()` — fetches conversation+sandbox+sandbox_spec, validates `RUNNING`, returns `agent_server_url`+`session_api_key` | `sources/openhands/openhands/app_server/app_conversation/app_conversation_router.py:134-217` |
| Frontend V1 service | `V1GitService.getGitChanges()` / `getGitChangeDiff()` via `axios GET {baseUrl}/api/git/...?path=` | `sources/openhands/frontend/src/api/git-service/v1-git-service.api.ts:32-91` |
| Frontend hooks | `useUnifiedGetGitChanges` / `useUnifiedGitDiff` (TanStack Query, 5-min stale) | `sources/openhands/frontend/src/hooks/query/use-unified-get-git-changes.ts:12`, `sources/openhands/frontend/src/hooks/query/use-unified-git-diff.ts:10-67` |
| Frontend diff UI | `FileDiffViewer` — Monaco `DiffEditor`, collapsed-by-default, view modes diff/old/new, markdown preview | `sources/openhands/frontend/src/components/features/diff-viewer/file-diff-viewer.tsx:75-254` |
| Frontend changes tab | `GitChanges` page — lists `changes.slice(0,100)`, WAITING_FOR_RUNTIME / NOT_A_GIT_REPO handling | `sources/openhands/frontend/src/routes/changes-tab.tsx:23-107` |
| Event log (artifact trace surrogate) | `EventServiceBase.save_event()` stores `event.id.hex.json` under `{prefix}/{user_id}/v1_conversations/{conv_hex}/`; `search_events`/`get_event`/`batch_get_events` | `sources/openhands/openhands/app_server/event/event_service_base.py:190-205`, `sources/openhands/openhands/app_server/event/event_service_base.py:86-110` |
| Event API | `GET /conversation/{id}/events/search?kind__eq=&timestamp__gte=` — paginated, filtered | `sources/openhands/openhands/app_server/event/event_router.py:19-110` |
| Conversation metadata | `AppConversationInfo` — `created_at`, `updated_at`, `sandbix_id`, `tags[archiveworkspacepath]`, `execution_status` | `sources/openhands/openhands/app_server/app_conversation/app_conversation_models.py:110-138`, `sources/openhands/openhands/app_server/app_conversation/app_conversation_models.py:40` |
| Workspace archive | `archive_workspace()` — feature-flagged (`RUNTIME_FILE_ARCHIVE_ENABLED`), captures `git-delta` and/or `tar.gz` from `GET /api/file/archive`, writes manifest with `conversation_id`, `sandbox_id`, `base_commit`, repo headers | `sources/openhands/openhands/app_server/sandbox/workspace_archive.py:77-83`, `sources/openhands/openhands/app_server/sandbox/workspace_archive.py:334-543` |
| Delete finalizer (no rollback) | `_finalize_sandbox_delete` — archives then deletes sandbox only if refcount==0; on REQUIRED failure keeps sandbox for idle-reap retry | `sources/openhands/openhands/app_server/app_conversation/app_conversation_router.py:878-936` |
| No artifact rollback API | Grep for `rollback`/`revert` returns only DB `db_session.rollback()` and template text, zero artifact handlers | `sources/openhands/openhands/app_server/services/db_session_injector.py:328`, `sources/openhands/openhands/app_server/integrations/templates/resolver/summary_prompt.j2:7` |
| No artifact review workflow | Grep finds only GitHub contribution workflows (`pr-artifacts.yml`, `pr-readiness-confirm.yml`) and SDK critiquer guides, no `Review` model/handler | `sources/openhands/.agents/skills/custom-codereview-guide.md:10`, `sources/openhands/.github/workflows/pr-artifacts.yml:9-18` |
| Comment model (not artifact-scoped) | `Comment(id,body,author,created_at)` and `BaseGitService._truncate_comment` service PR/issue comments, not artifact annotations | `sources/openhands/openhands/app_server/integrations/service_types.py:162-228` |
| Tests: diff proxy | `TestGitProxy::test_diff_routes_to_diff_runtime_path`, `test_diff_returns_409_when_sandbox_paused` | `sources/openhands/tests/unit/app_server/test_app_conversation_router.py:877-905`, `sources/openhands/tests/unit/app_server/test_app_conversation_router.py:1063-1075` |
| Test: SDK diff | `pydantic_diff` helper for settings diffs (not git diff coverage) | `sources/openhands/_sdk_inspect/sdk/utils/pydantic_diff.py:18-83` |
| Diff helper (file editor) | `visualize_diff` / `_has_meaningful_diff` — editor diff rendering, not artifact version diff | `sources/openhands/_sdk_inspect/tools/file_editor/utils/diff.py:67`, `sources/openhands/_sdk_inspect/tools/file_editor/definition.py:126` |

## Answers to Dimension Questions

**1. Can artifacts be compared?**
Yes—but only as git workspace diffs, not as versioned artifact-to-artifact diffs. The system exposes per-file `GitDiff{original,modified}` (`sources/openhands/_sdk_inspect/sdk/git/models.py:25`) produced by `get_git_diff` (`sources/openhands/_sdk_inspect/sdk/git/git_diff.py:53`) and per-repo `GitChange[]` (`sources/openhands/_sdk_inspect/sdk/git/git_changes.py:37`). In production, the frontend never calls these directly; it goes through the app-server proxy `GET /{id}/git/diff?path=&ref=` and `GET /{id}/git/changes?path=&ref=` (`sources/openhands/openhands/app_server/app_conversation/app_conversation_router.py:1308-1366`) which forward to the runtime's `/api/git/*`. The frontend's Changes tab (`sources/openhands/frontend/src/routes/changes-tab.tsx:94-104`) caps display at 100 entries and the viewer (`sources/openhands/frontend/src/components/features/diff-viewer/file-diff-viewer.tsx:152-171`) shows unified/side-by-side via Monaco. An optional `ref` param allows comparing against `HEAD` or any commit (`sources/openhands/_sdk_inspect/sdk/git/utils.py:129-135`), but there is no cross-run or cross-version artifact comparator beyond whatever git ref the caller supplies. No semantic/patch-aware diff beyond `git diff --name-status` + raw file contents.

**2. Is there a review workflow?**
No. There is no artifact review state machine—no `Review`, `Approval`, `Reviewer` model, no `pending_approval` status, no approve/reject endpoint, and no required-reviewer policy in the app-server (`grep` for `review`/`approval` in `openhands/app_server` returns only `Comment`/`ProviderTimeoutError` noise). The only "review" in the repo is external: GitHub PR triage (`sources/openhands/.agents/skills/custom-codereview-guide.md:10`) and CI checks (`sources/openhands/.github/workflows/pr-readiness-confirm.yml`, `sources/openhands/.github/workflows/pr-artifacts.yml:9`). Artifacts produced inside the sandbox are visible in the Changes tab but have no approve/reject affordance.

**3. Can artifacts be rolled back?**
No. There is no `rollback`/`revert`/`restore` handler for artifacts. A search for rollback in `openhands/` finds only SQL transaction rollbacks (`sources/openhands/openhands/app_server/services/db_session_injector.py:328`, `sources/openhands/openhands/app_server/app_conversation/app_conversation_router.py:931` in `delete_app_conversation`). The workspace archive (`sources/openhands/openhands/app_server/sandbox/workspace_archive.py:334`) is append-only durability—`git-delta` + `tar.gz` + manifest written to object storage—and explicitly "never raises" and favors completeness over dedup (`sources/openhands/openhands/app_server/sandbox/workspace_archive.py:538-543`). It has no read-back, no version index, and no revert-to-commit API. Reverting requires manual `git checkout`/`git revert` executed via the agent's bash tool; the platform does not orchestrate or authorize it, and the archive's `REQUIRED` flag merely blocks sandbox deletion on capture failure rather than performing a rollback.

**4. Are artifact changes traceable to runs?**
Partially and indirectly. There is no first-class `ArtifactChangeLog` table. Traceability is approximated via three loosely coupled mechanisms: (a) per-conversation event files under `{prefix}/{user_id}/v1_conversations/{conv_id}/` keyed by `event.id` and queryable by `kind`/`timestamp` (`sources/openhands/openhands/app_server/event/event_service_base.py:94-143`, `sources/openhands/openhands/app_server/event/event_router.py:29-110`); (b) conversation metadata with `id`/`sandbox_id`/`created_at`/`updated_at`/`tags` (`sources/openhands/openhands/app_server/app_conversation/app_conversation_models.py:110`, `sources/openhands/openhands/app_server/app_conversation/sql_app_conversation_info_service.py:132-149`); and (c) archive manifests embedding `conversation_id`, `sandbox_id`, `source_path`, `base_commit`, `repo_metadata`, `byte_count`, `created_at` (`sources/openhands/openhands/app_server/sandbox/workspace_archive.py:499-517`). Events are stored as flat JSON files (filesystem/GCS/S3/memory) and can be rehydrated via `batch_get_events` (`sources/openhands/openhands/app_server/event/event_service_base.py:199`), but there is no indexed join between a git diff hunk and the event/run that produced it, and archived blobs are not linked back into the event stream. Cost/usage attribution exists per conversation (`sources/openhands/openhands/app_server/app_conversation/sql_app_conversation_info_service.py:192-211`) but not per artifact mutation.

## Architectural Decisions

* **Git-as-artifact-store** (`sources/openhands/_sdk_inspect/sdk/git/git_diff.py:53`, `sources/openhands/_sdk_inspect/sdk/git/git_changes.py:37`, `sources/openhands/_sdk_inspect/sdk/git/utils.py:104`): Artifacts are not separately versioned; the workspace's git history *is* the artifact history. Keeps the implementation thin and delegates merge/conflict semantics to git, but conflates "agent output" with "repo state" and leaks git-specific failure modes (shallow clones, untracked files handled via `ls-files --others` at `sources/openhands/_sdk_inspect/sdk/git/git_changes.py:182-199`).
* **Runtime-proxied diff** (`sources/openhands/openhands/app_server/app_conversation/app_conversation_router.py:1224-1366`): Browsers cannot contact the sandbox directly (CORS), so the app-server proxies to `{agent_server_url}/api/git/*` with `X-Session-API-Key`. Centralizes auth and audit but adds a hop and couples UI freshness to sandbox `RUNNING` state (all git calls 409 when paused at `sources/openhands/openhands/app_server/app_conversation/app_conversation_router.py:1252-1256`).
* **Generator-based SDK transport** (`sources/openhands/_sdk_inspect/sdk/workspace/remote/remote_workspace_mixin.py:318-371`): Remote operations yield `{method,url,params,headers}` dicts; sync/async wrappers drive them via `httpx`. Enables shared logic across `RemoteWorkspace` and `AsyncRemoteWorkspace` but hides diff errors behind generic `Exception` raises (`sources/openhands/_sdk_inspect/sdk/workspace/base.py:149`).
* **Archive-as-backstop, not version control** (`sources/openhands/openhands/app_server/sandbox/workspace_archive.py:1-15`, `sources/openhands/openhands/app_server/app_conversation/app_conversation_router.py:878-936`): Workspace is captured at delete time to object storage (GCS/S3/local) with a manifest; idle/expiry reap is handled separately in `runtime-api`. The delete finalizer (`_finalize_sandbox_delete`) refuses to delete under `RUNTIME_FILE_ARCHIVE_REQUIRED` if capture fails. Provides durability without introducing a version catalog or rollback semantics.
* **Tag-projection pattern** (`sources/openhands/openhands/app_server/app_conversation/app_conversation_models.py:40-47`, `sources/openhands/openhands/app_server/app_conversation/app_conversation_models.py:140-179`): Provenance fields like `archiveworkspacepath`, `acpserver`, `agentprofileid` live in `tags: dict[str,str]` and are projected via `@computed_field`. Zero-migration extensibility, but tags are free-form strings with no schema enforcement beyond `^[a-z0-9]+$`.

## Notable Patterns

* **Thin git wrappers over subprocess** (`sources/openhands/_sdk_inspect/sdk/git/utils.py:17-59`): `run_git_command` uses `subprocess.run(capture_output=True, text=True, timeout=30)` with `shlex.join` logging—safe against injection, noisy on failure.
* **Size guard on diff** (`sources/openhands/_sdk_inspect/sdk/git/git_diff.py:26-84`): `MAX_FILE_SIZE_FOR_GIT_DIFF = 1 MiB`; oversize files raise `GitPathError` rather than streaming, preventing OOM on large binaries.
* **Rename-as-delete+add canonicalization** (`sources/openhands/_sdk_inspect/sdk/git/git_changes.py:104-120`): Git `R100 old new` is normalized to two `GitChange` entries (DELETED+ADDED), simplifying UI rendering at the cost of losing true rename tracking (`MOVED` enum exists at `sources/openhands/_sdk_inspect/sdk/git/models.py:10` but is never emitted).
* **Monaco diff rendering** (`sources/openhands/frontend/src/components/features/diff-viewer/file-diff-viewer.tsx:152-203`): `DiffEditor` collapsed by default; lazy-loads diff via `useUnifiedGitDiff` only when expanded (`sources/openhands/frontend/src/hooks/query/use-unified-git-diff.ts:98-102`), with side-by-side hidden for added/deleted files.

## Tradeoffs

* **Fidelity vs. completeness**: The git-delta archive respects `.gitignore` and needs the base tree to reconstruct (`sources/openhands/openhands/app_server/sandbox/workspace_archive.py:102-105`); the companion `tar.gz` is self-contained but large. Default `both` doubles storage until cost is measured—pragmatic but not principled lifecycle management.
* **Availability vs. durability**: `archive_required` blocks deletes on transient failures (`sources/openhands/openhands/app_server/sandbox/workspace_archive.py:541`) and the finalizer holds the sandbox for retry (`sources/openhands/openhands/app_server/app_conversation/app_conversation_router.py:904-926`). This favors durability but can leave sandboxes (and cost) lingering if archiving is flaky.
* **Simplicity vs. auditability**: Flat JSON event files (`sources/openhands/openhands/app_server/event/event_service_base.py:86-91`) are simple to stream to object storage and paginate, but querying requires loading all paths (`_search_paths` → `_load_events_from_paths` at `sources/openhands/openhands/app_server/event/event_service_base.py:56-110`) with no secondary index joining git hunks to events.
* **Passthrough vs. validation**: The git proxy forwards `ref` verbatim to `get_valid_ref` which does `rev-parse --verify {ref}^{commit}` (`sources/openhands/_sdk_inspect/sdk/git/utils.py:132`). Invalid refs surface as `GitCommandError` → 502, not 400 validation—a correct security posture (no injection) but poor UX.

## Failure Modes / Edge Cases

* **Sandbox not RUNNING** — `get_conversation_git_changes`/`get_conversation_git_diff` raise 409 (`sources/openhands/openhands/app_server/app_conversation/app_conversation_router.py:1252-1256`) or 404 if `MISSING`; the frontend shows `WAITING_FOR_RUNTIME` (`sources/openhands/frontend/src/routes/changes-tab.tsx:43-45`). No queued diff; user must resume.
* **Shallow clone** — `live_status_app_conversation_service.py:223` warns that history may be incomplete; `get_valid_ref` may fall back to empty tree (`sources/openhands/_sdk_inspect/sdk/git/utils.py:204-206`) yielding a diff against the empty tree rather than the real base.
* **Oversized/new file** — `get_git_diff` raises `GitPathError` for >1 MiB (`sources/openhands/_sdk_inspect/sdk/git/git_diff.py:80-84`) or missing file (`sources/openhands/_sdk_inspect/sdk/git/git_diff.py:74-75`); no chunked diff or binary handling.
* **Archive misconfiguration** — unset `RUNTIME_FILE_ARCHIVE_BUCKET` or invalid `RUNTIME_FILE_ARCHIVE_FORMAT` logs loudly but allows delete to proceed (`sources/openhands/openhands/app_server/sandbox/workspace_archive.py:369-394`), silently losing data under REQUIRED.
* **Archive stream OOM avoidance** — `_stream_to_tempfile` streams to `NamedTemporaryFile` (`sources/openhands/openhands/app_server/sandbox/workspace_archive.py:176-193`) rather than buffering, but cleanup is best-effort (`_cleanup_tempfile` at `sources/openhands/openhands/app_server/sandbox/workspace_archive.py:167-173`) and a crash mid-stream can orphan temp files.
* **Rename/copy edge** — `get_changes_in_repo` maps `R`/`C` prefixes with similarity percentages (`sources/openhands/_sdk_inspect/sdk/git/git_changes.py:104-135`) but unknown single-char statuses (`*`, `??` mapped at `sources/openhands/_sdk_inspect/sdk/git/git_changes.py:149-152`) can produce unexpected `ADDED`/`UPDATED` rather than surfacing the raw git state.
* **100-item truncation** — Changes tab `slice(0,100)` (`sources/openhands/frontend/src/routes/changes-tab.tsx:96`) silently hides beyond-100 changes; diff viewer lazy loading (`enabled: !isCollapsed` at `sources/openhands/frontend/src/components/features/diff-viewer/file-diff-viewer.tsx:101`) means collapsed diffs are not prefetched.

## Future Considerations

* Introduce a typed `Artifact`/`ArtifactVersion` model with content hash, `conversation_id`/`event_id`/`sandbox_id` foreign keys, and `created_by` attribution—bridging the current git-diff layer and the event log so "which run touched file X?" is a query, not a forensic `git blame`.
* Add a review/approval state machine (`pending_approval` → `approved`/`rejected` with required reviewers) gating any deploy/export path; currently the PR-based review is repo-level, not artifact-level.
* Implement a true revert API (`POST /{id}/git/revert?ref=`) that validates the ref via `get_valid_ref` (`sources/openhands/_sdk_inspect/sdk/git/utils.py:129`) and records the revert as a new event, giving the audit trail the dimension asks for.
* Index the workspace archive: persist archive keys + manifests back into `conversation_metadata` or a dedicated `workspace_archives` table so archives are discoverable and restorable, rather than opaque objects under `workspace-archives/{sandbox}/{conv}/{ts}.*` (`sources/openhands/openhands/app_server/sandbox/workspace_archive.py:403`).
* Expand test coverage beyond proxy happy-paths (`sources/openhands/tests/unit/app_server/test_app_conversation_router.py:828-1075`) to include diff content correctness, rename handling, and archive REQUIRED vs. not-required semantics.

## Questions / Gaps

* **Artifact vs. git conflation** — No evidence in `openhands/app_server` or `_sdk_inspect` of a first-class artifact abstraction separate from the git working tree. Searched `artifact*`, `version.*history`, `change_log` in `openhands/` — zero hits. If artifacts are defined elsewhere (e.g., SDK `agent_server` not in this source snapshot), it is out of scope under source-isolation rules.
* **Review policy storage** — No `policy`/`approval`/`reviewer` columns in `StoredConversationMetadata` (`sources/openhands/openhands/app_server/app_conversation/sql_app_conversation_info_service.py:132-190`). Whether cloud/enterprise adds this in `enterprise/` was not inspected (isolated to selected source only).
* **Comment threading on artifacts** — `Comment`/`get_pr_comments`/`get_thread_from_comment_graphql_query` (`sources/openhands/openhands/app_server/integrations/github/queries.py:49-111`) are PR/issue-scoped; no `ArtifactComment` or inline hunk comment model was found. Inline annotation would require new storage and a `file:line` anchor—absent.
* **Cross-conversation diff** — No endpoint compares two `conversation_id` values; only single-conversation `?path=&ref=` exists. Whether this is intentional (git is cross-conversation by nature) or a gap depends on the product definition of "artifact."
* **Observability of bad-state revert** — The archive's "capture completeness over dedup" comment (`sources/openhands/openhands/app_server/sandbox/workspace_archive.py:410`) acknowledges orphan blobs on retry but no metric/alert is emitted; proving "a bad change can be reverted with full audit trail" would require ledger + revert coverage that does not exist.

---

Generated by `16.02-diff-review-and-rollback` against `openhands`.
