# Source Analysis: openhands

## 16.03 Artifact Provenance and Reproducibility

### Source Info

| Field | Value |
|-------|-------|
| Name | openhands |
| Path | `studies/agent-harness-study/sources/openhands` |
| Language / Stack | Python 3.12-3.13 (FastAPI/Poetry), Node 22, Docker, React |
| Analyzed | 2026-08-28 |

## Summary

OpenHands traces **conversations** (the primary artifact) to inputs via persistent `AppConversationInfo` + ordered `Event` store, and adds explicit Agent Profile provenance tags and LiteLLM tracing metadata. However it does **not** provide deterministic reproduction of LLM-driven artifacts (no recorded prompts/seeds/temperature, no build attestations for the OSS image, no locked tool-version snapshot per trajectory). Reproducibility is best-effort export/replay of raw events, not tested in CI.

## Rating

**4 / 10 — Present but inconsistent, weakly documented, fragile**

Rationale: Conversation-level provenance (model, repo/branch, trigger, tags, secrets, LLM metadata) and full event history are implemented and persisted (`StoredConversationMetadata`, `EventService`). Tool-invocation artifacts (patches, browser recordings, terminal output) ride inside events and export zip. But: LLM non-determinism is uncontrolled (temperature, seed, sampler params not systematically logged per call), build provenance/SBOM only enabled for `enterprise-server` image, no SLSA attestation, no version-pinned snapshot of skills/plugins per trajectory, and no CI job asserting bit-for-bit reproduction. Meets rubric "present but inconsistent".

## Evidence Collected

Every entry MUST include a file path with line numbers. Format: `path/to/file.ts:NN`.

| Area | Evidence | File:Line |
|------|----------|-----------|
| Provenance fields — Agent Profile tags | `AGENT_PROFILE_ID_TAG_KEY='agentprofileid'` and `AGENT_PROFILE_REVISION_TAG_KEY='agentprofilerevision'` riding `AppConversationInfo.tags` dict + `@computed_field launched_agent_profile` projection; stamped at launch from resolved `UserInfo.active_agent_profile_id/revision` | `openhands/app_server/app_conversation/app_conversation_models.py:42-48` `openhands/app_server/app_conversation/app_conversation_models.py:156-178` |
| Provenance fields — ACP provider tag | `ACP_SERVER_TAG_KEY='acpserver'` persisted in tags, surfaced as `acp_server` computed field; written from `user.agent_settings.acp_server` at conversation creation | `openhands/app_server/app_conversation/app_conversation_models.py:30-35` `openhands/app_server/app_conversation/app_conversation_models.py:140-154` `openhands/app_server/app_conversation/live_status_app_conversation_service.py:570-582` |
| Provenance fields — stamping logic | Launch code resolves `profile_user = get_user_info(resolve_agent_profile=True, override_agent_profile_id=...)` and writes both profile id/revision and `ARCHIVE_WORKSPACE_PATH_TAG_KEY` so delete-time archive uses original path regardless of later grouping changes | `openhands/app_server/app_conversation/live_status_app_conversation_service.py:551-570` `openhands/app_server/app_conversation/app_conversation_models.py:46-48` |
| Input recording — conversation metadata | `AppConversationInfo` captures `llm_model`, `agent_kind`, `selected_repository/branch`, `git_provider`, `trigger`, `pr_number`, `parent_conversation_id`, `metrics` (MetricsSnapshot/TokenUsage), `tags`, timestamps; persisted to `StoredConversationMetadata` columns with `conversation_version='V1'` | `openhands/app_server/app_conversation/app_conversation_models.py:110-139` `openhands/app_server/app_conversation/sql_app_conversation_info_service.py:132-190` |
| Input recording — start request | `AppConversationStartRequest` records `initial_message`, `system_message_suffix`, `llm_model`, `agent_profile_id` override, `selected_repository/branch`, `git_provider`, `suggested_task`, `plugins` (PluginSpec with source/ref/repo_path+parameters), `secrets` (SecretStr dict, merged with precedence) | `openhands/app_server/app_conversation/app_conversation_models.py:221-273` |
| Input recording — event store | Every `Event` (user message, tool call, observation, `ConversationStateUpdateEvent`) stored as `{hex}.json` under `{prefix}/{user_id}/v1_conversations/{conv_hex}`; `EventServiceBase.get_event/search_events/iter_events_for_export/_search_paths` provides single-source trace for trajectory export | `openhands/app_server/event/event_service_base.py:85-97` `openhands/app_server/event/event_service_base.py:145-156` `openhands/app_server/conversation_paths.py:11-73` |
| Input recording — cost/metrics provenance | `StoredConversationCostEvent` rows per `usage_id` (agent/condenser/profile:*) with `cost_delta`, `prompt_tokens`, `completion_tokens`, `llm_model`, `occurred_at`; updated from `ConversationStats.usage_to_metrics` via stats events | `openhands/app_server/app_conversation/sql_app_conversation_info_service.py:192-211` `openhands/app_server/app_conversation/sql_app_conversation_info_service.py:505-568` |
| Input recording — LLM observability metadata | `get_llm_metadata()` builds `litellm_extra_body.metadata = {trace_version, tags:[app:openhands, model:..., type:..., web_host,..., openhands_version,..., conversation_version:V1], session_id, trace_user_id, repo_name, git_provider, branch}`; only injected for `openhands/` models via `should_set_litellm_extra_body`; `_apply_server_agent_overrides` and `_build_observability_context` add repo/commit tags | `openhands/app_server/utils/llm_metadata.py:10-91` `openhands/app_server/app_conversation/live_status_app_conversation_service.py:1647-1705` `openhands/app_server/app_conversation/live_status_app_conversation_service.py:1587-1619` |
| Input recording — user/secrets/settings provenance | `_setup_conversation_secrets` collects git provider tokens (LookupSecret with JWT or StaticSecret), `conversation_secret_enricher` integration secrets; `Settings` model normalizes `agent_settings`, `conversation_settings`, `llm_profiles`, `marketplace_registrations`; secrets merged with API-provided with precedence | `openhands/app_server/app_conversation/live_status_app_conversation_service.py:1162-1250` `openhands/app_server/settings/settings_models.py:663-673` |
| Build reproducibility — pinned dependencies | `pyproject.toml:23-109` exact pins (`orjson==3.11.8`, `litellm==1.84.1`, `openai==2.33.0`, etc.), `poetry.lock` + `uv.lock` committed, `poetry.lock` header records tool version, `Dockerfile:24` `ARG POETRY_VERSION=2.3.4` and `78-86` pin `pip==26.0.1` "for reproducible builds" | `pyproject.toml:23-109` `containers/app/Dockerfile:24` `containers/app/Dockerfile:78-86` |
| Build reproducibility — provenance/SBOM | `_build-image.yml:26-32` inputs `provenance` and `sbom` passed to `docker/build-push-action`; `ghcr-build.yml:41-43` sets `provenance:true sbom:true` only for `enterprise-server` (app image defaults to false) | `.github/workflows/_build-image.yml:26-32` `.github/workflows/_build-image.yml:94-95` `.github/workflows/ghcr-build.yml:42-43` |
| Build reproducibility — build version tagging | `Dockerfile:1,38,46` `ARG OPENHANDS_BUILD_VERSION` → `ENV OPENHANDS_BUILD_VERSION`, injected from `_build-image.yml:88` as `RELEVANT_REF_NAME` (pr-N or branch) | `containers/app/Dockerfile:1,38,46` `.github/workflows/_build-image.yml:88` |
| Trajectory export / replay | `EventService.iter_events_for_export` + `EventServiceBase.iter_events_for_export` iterate all events in timestamp order for export; `LiveStatusAppConversationService.open_conversation_export` streams zip; `config.template.toml:29-38` legacy `save_trajectory_path` / `replay_trajectory_path` flags | `openhands/app_server/event/event_service.py:48-58` `openhands/app_server/event/event_service_base.py:145-156` `config.template.toml:29-38` |
| Critic reproducibility marker | `APIBasedCritic.evaluate:120-121` collects `event_ids = [event.id ...]` into `CriticResult.metadata` "for reproducibility" — event IDs are reproducibility anchors for critic categorization | `_sdk_inspect/sdk/critic/impl/api/critic.py:120-132` |
| CI reproducibility testing | No evidence found — grep for `reproducibility`, `provenance.*test`, `replay.*test` yields zero; `tests/unit/app_server/test_live_status_app_conversation_service.py` tests trajectory download but not bit-for-bit reproduction; no workflow asserts deterministic rebuild | `tests/unit/app_server/test_live_status_app_conversation_service.py:1406` (negative evidence) |

## Answers to Dimension Questions

### 1. Can every artifact be traced to its inputs?

**Partially.** The conversation artifact (the unit OpenHands produces) is traceable: `AppConversationInfo` persists model, repo/branch/provider, trigger, parent link, tags carrying Agent Profile and ACP provenance (`openhands/app_server/app_conversation/app_conversation_models.py:110-178`), and the `Event` file-store retains every `SendMessageRequest` and tool invocation in timestamp order (`openhands/app_server/event/event_service_base.py:145-156`). Derived artifacts ride inside those events (file edits captured as tool-call events, browser recordings flushed to `recording-{timestamp}.json` under `tools/browser_use/event_storage.py:41-54`, patches stored as observations). Export zips the full chain (`openhands/app_server/app_conversation/live_status_app_conversation_service.py:2898`).

Gaps: per-artifact granularity is coarse — there is no explicit "this file patch was produced by prompt X, tool Y, model Z at temp T" manifest. LLM temperature/seed/stop-words are model config, not per-call provenance (`config.template.toml:298` is template-level, `_sdk_inspect/sdk/llm/llm.py:402` declares field but not systematically persisted per event). Skill/plugin `repo_path`/`ref` is recorded at start (`PluginSpec` `openhands/app_server/app_conversation/app_conversation_models.py:71-81`) but not a resolved commit SHA snapshot at execution time, so the effective skill version can drift.

### 2. Is reproduction deterministic?

**No.** OpenHands does not claim or enforce deterministic reproduction. LLM calls are intrinsically stochastic (temperature defaults vary per model, no global seed), condensers truncate/LLM-summarize history (`config.template.toml:240-293`), and tool outputs depend on live network/filesystem state. The legacy `replay_trajectory_path` (`config.template.toml:35-38`) replays stored prompts before responding, but does not lock model weights, sampler, or external tool versions. Export is reproducible as *data* (re-exporting same event store yields same zip), but re-running the agent with same `initial_message` will diverge. No workflow fixes `PYTHONHASHSEED`, wall-clock, or container image digest per trajectory.

### 3. Are all contributing factors recorded?

**No — major factors are missing or only partially recorded.** Recorded: model name, repo/branch/provider/commit (best-effort `git rev-parse HEAD` at launch `openhands/app_server/app_conversation/live_status_app_conversation_service.py:1820-1839`), agent profile id/revision, ACP server, workspace path, conversation settings (`max_iterations`, `confirmation_mode`, `security_analyzer`), marketplace composition, secrets (as `LookupSecret`/`StaticSecret` descriptors), cost buckets (`StoredConversationCostEvent` `openhands/app_server/app_conversation/sql_app_conversation_info_service.py:192-211`), LiteLLM metadata tags. Not recorded: exact inference parameters per turn (temperature/top_p/seed), full agent code version beyond `OPENHANDS_BUILD_VERSION` branch tag (not commit SHA for OSS image), pinned marketplace/skill commit SHA (only `ref` if supplied), base container image digest, `condenser` window choices that prune context, and user approvals (no human-in-the-loop approval record beyond trigger enum).

### 4. Is reproducibility tested in CI?

**No evidence found.** Greps for `reproduc*`, `provenance`, `sbom`, `attestation`, `replay` across workflows and `tests/` returned only docs/build provenance toggles, not tests. The only related CI checks are: `ghcr-build.yml` building multi-arch images with cache (`_build-image.yml:90-93`), poetry/pip pins (`containers/app/Dockerfile:24,82-86`), and unit tests for trajectory download shape (`tests/unit/app_server/test_live_status_app_conversation_service.py:1406,1481`) and for llm metadata (`_sdk_inspect` not CI-gated). No job rebuilds an image and verifies identical digest, no trajectory-replay regression, no `docker buildx ... --provenance` verification.

## Architectural Decisions

| Decision | Evidence | Effect on Provenance |
|----------|----------|----------------------|
| File-per-event JSON under `{user_id}/v1_conversations/{id}/{event.hex}.json` | `openhands/app_server/event/event_service_base.py:190-197` | Strong: immutable ordered log enables full trajectory export and forensic replay; depends on filesystem durability |
| Tags-dict + `@computed_field` for provenance (ACP/profile/workspace) | `openhands/app_server/app_conversation/app_conversation_models.py:140-178` `openhands/app_server/app_conversation/live_status_app_conversation_service.py:551-582` | Avoids DB migration per provenance field; but tags are free-form `dict[str,str]` so schema is implicit and malformed values silently drop (`try: UUID(pid) except: return None`) |
| Agent Profile resolution stamps effective settings, not persisted pointer | `openhands/app_server/settings/settings_models.py:662-673` `_resolved_view: PrivateAttr` | Prevents accidental persistence of resolved view via `store()` guard; conversation provenance reflects what *ran* even with one-off `agent_profile_id` override |
| LiteLLM `extra_body.metadata` for SaaS tracing | `openhands/app_server/utils/llm_metadata.py:28-91` `openhands/app_server/app_conversation/live_status_app_conversation_service.py:1666-1677` | Moves provenance from DB into observability pipeline (Laminar); only for `openhands/` models — third-party models lose metadata |
| `StoredConversationCostEvent` per-bucket ledger with monotonic guard | `openhands/app_server/app_conversation/sql_app_conversation_info_service.py:569-678` | Prevents cost regression on stale stats snapshots; enables attribution per `usage_id`/`llm_model`; unattributed cost drain handling is subtle and tested only via unit tests |
| OSS build without SLSA provenance/SBOM | `.github/workflows/_build-image.yml:26-32` default `provenance:false` `ghcr-build.yml:42-43` only enterprise `true` | Minimal build provenance; users cannot verify OSS artifact's builder or dependency closure |

## Notable Patterns

- **Tag-riding provenance**: Stamping job/project identity as string tags rather than columns (`archiveworkspacepath`, `agentprofileid`, `acpserver`) for zero-migration evolution.
- **Secret-as-Lookup**: Tokens kept as `LookupSecret(url, headers)` with per-secret JWT scoped to `web_url/api/v1/webhooks/secrets`, so provenance shows *which* secret contributed without materializing value — but redaction in `settings_models.py:282-603` means omitted headers must be carried over correctly or silently lost.
- **Monotonic cost ledger**: Cost/tokens only advance (`max()` guard `openhands/app_server/app_conversation/sql_app_conversation_info_service.py:574-602`, `delta_cost < 0` early return), paired with per-bucket `_record_bucket_cost_deltas` avoiding double-counting unattributed legacy rows.

## Tradeoffs

- **Frozen secrets + redaction recovery** (`_preserve_redacted_mcp_secrets` `openhands/app_server/settings/settings_models.py:560-603`): Protects leakage but adds complex merging that can hide provenance of which secret value actually reached the agent (stripped to absent then recovered from stored dump).
- **Event store vs DB**: Events on filesystem bypass transactional guarantees with `StoredConversationMetadata`; webhook lifecycle can `save_app_conversation_info` with UTC-normalized timestamps (`sql_app_conversation_info_service.py:68-87`) racing stat events corrected by `SELECT ... FOR UPDATE` only on Postgres (SQLite is no-op `with_for_update()`).
- **Group-vs-isolate workspaces**: `grouped_workspace_dir` (`openhands/app_server/settings/settings_models.py:637-651`) nests per-conversation dirs under shared sandbox; archive path pinned at creation prevents later grouping changes from breaking provenance of where aid was written.
- **Provenance only for enterprise build**: Enterprise image gets Docker provenance+SBOM; OSS `openhands` image consumers have no attestation to reproduce the build.

## Failure Modes / Edge Cases

- **Malformed provenance tag does not fail loudly**: `LaunchedAgentProfile(UUID(pid))` ValueError swallows to `return None` (`openhands/app_server/app_conversation/app_conversation_models.py:176-178`) — a corrupted `tags['agentprofileid']` silently reports "no provenance".
- **Out-of-order stats snapshots regress attempt**: Negative `delta_cost` bails without updating ledger (`sql_app_conversation_info_service.py:555-563`); concurrent snapshots serialize only via DB `SELECT FOR UPDATE` which is unavailable on SQLite test fixtures.
- **Commit lookup timeout/drop**: `git rev-parse HEAD` has 10s timeout (`live_status_app_conversation_service.py:1832`); on empty repo or throttled workspace it returns `''` and `observability_metadata` omits `commit` silently.
- **Tool nondeterminism unrecorded**: Browser `recording-{timestamp}.json` (`tools/browser_use/event_storage.py:41-53`) uses wall time for filenames; replay order would depend on enumerating `_search_paths` lexical order rather than logical causality.
- **Secret override silent**: API-provided `secrets` dict silently overrides DB secret with same name (`live_status_app_conversation_service.py:2004-2011` warns only on existing key), breaking lineage of which source produced the effective secret.

## Future Considerations

- Add per-turn LLM provenance to each `LLMConvertibleEvent`: `{model, temperature, seed, stop, litellm_extra_body, tool_schema_hash}` so any patch's prompt→completion chain is independently verifiable.
- Snapshot resolved `PluginSource` `ref` to commit SHA at `_build_start_conversation_request_for_user` time and record alongside `MarketplaceRegistration` composition; include `openhands/__version__` commit SHA (not branch tag) as `OPENHANDS_BUILD_VERSION` for OSS builds and emit SLSA provenance + SBOM for both images.
- Promote trajectory replay to a first-class reproducible harness: fixed `PYTHONHASHSEED`, model temp=0, condenser=`noop`, synthetic clock, with CI job that replays a golden trajectory and asserts identical event output (minus timestamps/ids).
- Replace free-form tag provenance with explicit columns + Alembic migration (as `tags` JSON is opaque to SQL filters and risks collisions with user tags like `repo_name`).

## Questions / Gaps

- Is there an intent to expose provenance to end-users via trajectory zip manifest (currently each event file + `meta.json`, but no top-level `inputs.json` summarizing model/seed/tools/plugins/commit)? No evidence found.
- How are approvals captured, if any? `ConversationTrigger` enum has no `APPROVAL` state; no `EventKind` suggests human approval gating is out-of-scope for V1.
- Does the agent-server (external `openhands-sdk` repo, not inspected here per isolation rules) attach its own tool-version provenance? SDK path `_sdk_inspect/sdk/conversation/impl/local_conversation.py:426` mentions "resolves refs to commit SHAs for deterministic resume" but implementation not in this source tree.

---
Generated by `16.03-artifact-provenance-and-reproducibility` against `openhands`.
