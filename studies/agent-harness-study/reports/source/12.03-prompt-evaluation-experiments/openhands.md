# Source Analysis: openhands

## Dimension 12.03: Prompt Evaluation and Experiments

### Source Info

| Field | Value |
|-------|-------|
| Name | openhands |
| Path | `studies/agent-harness-study/sources/openhands` |
| Language / Stack | Python (openhands-ai 1.11.0, openhands-sdk 1.37.1) + TypeScript/React frontend, Jinja2 prompts, LiteLLM |
| Analyzed | 2026-08-29 |

## Summary

OpenHands treats prompts as Jinja2 templates rendered at runtime (`_sdk_inspect/sdk/context/prompts/prompt.py:88`, `_sdk_inspect/sdk/agent/base.py:381`) with six system-prompt variants and dozens of resolver/suggested-task templates. Prompt-change safety is **ad-hoc**: there are zero golden-output or snapshot tests for system prompts and no checked-in eval datasets (the `evaluation/` directory is `.gitignored:191` and absent from source; docs point to the external `OpenHands/benchmarks` repo). Experiment tracking / A/B testing is absent — no Langfuse/LangSmith/LMNR prompt-experiment code was found, only `lmnr` as a tracing dependency in `pyproject.toml:259`/`poetry.lock:2675` with no prompt-experiment usage. The sole structured LLM-as-judge is the **Critic subsystem** (`_sdk_inspect/sdk/critic/base.py:57`, `impl/api/critic.py:47`, `impl/api/client.py:74`): a remote `POST /classify` classifier that predicts `success` plus ~24 behavioral/infrastructure labels and drives optional iterative refinement, but it is a **runtime agent-quality gate**, not an offline prompt-regression harness. CI (`.github/workflows/py-tests.yml:60`, `lint.yml:42`) runs unit tests and lint only — no prompt regression gate — so a prompt edit can ship with no automated confidence check.

## Rating

**3 / 10 — Absent / ad-hoc**

Prompt rendering is cleanly abstracted and the Critic LLM-as-judge is mature, but evaluation datasets, prompt golden tests, experiment tracking, and CI regression gates for prompt changes are missing inside the studied source. A prompt change today deploys without any automated check that it won't regress. Work in `OpenHands/benchmarks` is acknowledged externally but not integrated as a pre-deploy prompt quality gate in this repo.

## Evidence Collected

Every entry MUST include a file path with line numbers. Format: `path/to/file.ts:NN`.

| Area | Evidence | File:Line |
|------|----------|-----------|
| Prompt templating | `AgentBase.system_prompt` inline override and `system_prompt_filename` fields (mutual exclusion validator) | `_sdk_inspect/sdk/agent/base.py:146-203` |
| Prompt rendering | `static_system_message` property renders Jinja2 template via `render_template(prompt_dir, template_name, **kwargs)` with model-family injection and `security_policy_filename` injection | `_sdk_inspect/sdk/agent/base.py:368-408` |
| Template loader | `FlexibleFileSystemLoader` + `render_template()` with `lru_cache`, bytecode cache, `refine` filter, absolute/relative path support | `_sdk_inspect/sdk/context/prompts/prompt.py:16-114` |
| System prompts | Primary prompt + 4 variants that all `include system_prompt.j2` | `_sdk_inspect/sdk/agent/prompts/system_prompt.j2:1`, `system_prompt_planning.j2:1`, `system_prompt_long_horizon.j2:1`, `system_prompt_interactive.j2:1`, `system_prompt_tech_philosophy.j2:1` |
| Model-specific prompts | Per-model-family overrides injected via `<IMPORTANT>` block | `_sdk_inspect/sdk/agent/prompts/system_prompt.j2:133-147` + `model_specific/anthropic_claude.j2`, `openai_gpt/gpt-5.j2`, `google_gemini.j2` |
| Security prompt | System prompt conditionally includes security policy template (`{% include security_policy_filename %}`) | `_sdk_inspect/sdk/agent/prompts/system_prompt.j2:78-88`, `_sdk_inspect/sdk/agent/prompts/security_policy.j2:1` |
| Condenser prompt | Summarization prompt for `LLMSummarizingCondenser` (separate Jinja2 template) | `_sdk_inspect/sdk/context/condenser/prompts/summarizing_prompt.j2:1-55` |
| Suggested-task prompts | Issue/PR prompts rendered via `get_prompt_for_task()` dispatching to 4 Jinja templates | `openhands/app_server/integrations/service_types.py:90-112`, `openhands/app_server/integrations/templates/suggested_task/*.j2` |
| Resolver prompts | ~20 additional Jinja prompt templates for GitHub/GitLab/Bitbucket/Azure/Jira/Linear/Slack resolvers | `openhands/app_server/integrations/templates/resolver/**/*.j2` (e.g., `github/issue_prompt.j2`) |
| Eval datasets — absent | `evaluation/` explicitly gitignored (`evaluation/evaluation_outputs`, `evaluation/swe_bench/*`, etc.) and `Glob evaluation/**/*` returns 0 files inside source | `.gitignore:191-202` |
| Eval datasets — externalized | CONTRIBUTING and Development point evaluation to *external* repo `OpenHands/benchmarks`; README notes SDK/backend moved repos; CREDITS references SWE-Bench/HumanEval integration but no dataset files in tree | `CONTRIBUTING.md:39`, `CONTRIBUTING.md:57`, `Development.md:330`, `CREDITS.md:27-41`, `.dockerignore:22` |
| Prompt tests — minimal | Only prompt-adjacent test is `test_apply_suggested_task_sets_prompt_and_trigger` (suggested-task path), not system-prompt golden output; grep `prompt.*test\|system_prompt` in `tests/**` yields zero system-prompt snapshot tests | `tests/unit/app_server/test_live_status_app_conversation_service.py:282-335` |
| Prompt tests — zero snapshot | No `snapshot`, `golden`, or dataset fixtures for `system_prompt.j2` rendering; `grep snapshot` hits are unrelated (MetricsSnapshot/StateSnapshot) | `tests/unit/` grep: zero prompt-snapshot results |
| Experiment tracking — absent | No `experiment`, `ab_test`, `langfuse` experiment, `langsmith`, or `mlflow` prompt-experiment code; only deps (`lmnr`, `langfuse`) listed as tracing, no usage for prompt A/B; `grep literalai\|langfuse\|langsmith\|wandb\|mlflow` in source yields only `poetry.lock` entries | `pyproject.toml:259` (`lmnr>=0.7.20`), `poetry.lock:2670-4271` |
| LLM-as-judge — CriticBase | Abstract interface `evaluate(events, git_patch) -> CriticResult` with `mode` (`finish_and_message` vs `all_actions`) and `IterativeRefinementConfig(threshold, max_iterations)` | `_sdk_inspect/sdk/critic/base.py:20-114` |
| LLM-as-judge — APIBasedCritic | Production critic: serializes conversation to chat dicts, posts to `POST {server_url}/classify`, maps probs to `success`+labels, categorizes via taxonomy, `should_refine()` triggers on score or high-prob agent-issue labels | `_sdk_inspect/sdk/critic/impl/api/critic.py:47-206` |
| LLM-as-judge — CriticClient | HTTP client with retry (3x on 500), `ChatTemplateRenderer` (HF tokenizer `Qwen/Qwen3-4B-Instruct-2507`), default `server_url=https://llm-proxy.app.all-hands.dev/vllm`, `model_name=critic`, label space = 1 success + 3 sentiment + 13 agent + 2 infra + 8 user = 27 labels | `_sdk_inspect/sdk/critic/impl/api/client.py:74-335` |
| LLM-as-judge — Taxonomy | `FEATURE_CATEGORIES` mapping 24 labels to `agent_behavioral_issues`/`user_followup_patterns`/`infrastructure_issues`; `categorize_features()` filters `prob < 0.2`, softmax-normalizes sentiment, sorts descending | `_sdk_inspect/sdk/critic/impl/api/taxonomy.py:8-180` |
| LLM-as-judge — Result | `CriticResult{score:0-1, message, metadata{categorized_features,event_ids}}` with `THRESHOLD=0.5`, star rendering, `_append_categorized_features` | `_sdk_inspect/sdk/critic/result.py:7-148` |
| LLM-as-judge — Settings | User-facing settings `verification.critic_enabled` (default false), `critic_mode`, `enable_iterative_refinement`, `critic_threshold=0.6`, `max_refinement_iterations=3`, `critic_server_url/model_name` overrides | `_sdk_inspect/sdk/settings/model.py:143-231`, `_sdk_inspect/sdk/settings/model.py:810-846` |
| LLM-as-judge — Wiring | Agent `step()` and `ResponseDispatch.create_and_emit_message` gate critic via `_should_evaluate_with_critic()`; critic result attached to `ActionEvent`/`MessageEvent.critic_result` | `_sdk_inspect/sdk/agent/critic_mixin.py:26-137`, `_sdk_inspect/sdk/agent/agent.py:893-900`, `_sdk_inspect/sdk/agent/response_dispatch.py:274-278` |
| LLM-as-judge — Simple critics | `AgentFinishedCritic` (FinishAction + non-empty patch -> 0/1) and `EmptyPatchCritic` as non-LLM baselines | `_sdk_inspect/sdk/critic/impl/agent_finished.py:20-60`, `_sdk_inspect/sdk/critic/impl/empty_patch.py:4`, `_sdk_inspect/sdk/critic/impl/pass_critic.py:4` |
| CI prompt tests | `py-tests.yml` runs `pytest --forked -n auto ./tests/unit` and coverage — no prompt dataset, no snapshot, no LLM eval job; `lint.yml` runs ruff/mypy/frontend lint only | `.github/workflows/py-tests.yml:60`, `.github/workflows/lint.yml:42-57`, `pytest.ini:1-4` |
| Analytics — not prompt eval | Analytics service captures `prompt_tokens`/`completion_tokens` for cost accounting, not prompt-quality evaluation | `openhands/analytics/analytics_service.py:233`, `openhands/app_server/app_conversation/sql_app_conversation_info_service.py:158-675` |
| Pre-commit — no prompt gate | Python pre-commit config checks formatting/types, not prompt output; hook path is `dev_config/python/.pre-commit-config.yaml` | `.openhands/pre-commit.sh:21` |

## Answers to Dimension Questions

**1. Are prompt changes tested?**
No. Prompts are Jinja2 templates rendered via `_sdk_inspect/sdk/context/prompts/prompt.py:88` and `_sdk_inspect/sdk/agent/base.py:381-408`. There are zero checked-in tests that render `system_prompt.j2` (or its variants/partials) and assert on expected output, no snapshot/golden files, and no dataset-driven assertions. The only prompt-cognizant test is `tests/unit/app_server/test_live_status_app_conversation_service.py:282` which exercises `apply_suggested_task` → `get_prompt_for_task()` for suggested-task prompts, not the core agent system prompt. Search `tests/**/*prompt*` and grep `system_prompt` in `tests/` both return empty. A change to `system_prompt.j2:1-149` (or any `model_specific/*.j2` or resolver template) would land with no automated detection of semantic drift.

**2. Are experiments tracked?**
No. No experiment tracker for prompt A/B tests was found. Grep for `literalai|langfuse|langsmith|wandb|mlflow|experiment_tracker|ab_test` in the source proper yields only `poetry.lock` dependency entries and unrelated comment text. `pyproject.toml:259` lists `lmnr>=0.7.20` as a tracing dependency but no code imports it for prompt-experiment management. There is no prompt registry, no versioned prompt store, no comparison harness, and no feature-flagged rollout for prompts.

**3. Is LLM-as-judge used for evaluation?**
Yes — but as a **runtime trajectory critic**, not an offline prompt-evaluation harness. The Critic subsystem (`_sdk_inspect/sdk/critic/base.py:57`) provides:
- `APIBasedCritic` (`_sdk_inspect/sdk/critic/impl/api/critic.py:47`) — calls a remote classifier (`_sdk_inspect/sdk/critic/impl/api/client.py:266` → `POST /classify`) with the full `SystemPromptEvent` + conversation messages + tool definitions, receiving per-label probabilities (27 labels) and producing `CriticResult{score=probs["success"], metadata.categorized_features}`. Taxonomy filtering is at `taxonomy.py:82`.
- Gating: only runs on `FinishAction`/`MessageEvent` by default (`mode=finish_and_message`, `critic_mixin.py:39`, `response_dispatch.py:274`), or every action if configured (`all_actions`).
- Iterative refinement loop (`critic_mixin.py:88-137`) feeds critic feedback back as a follow-up prompt when `score < threshold` or high-prob `agent_behavioral_issues` exceed `issue_threshold=0.75` (`critic.py:135-144`), with `get_followup_prompt()` rendering a detailed rubric-aware message.
- Deployment is controlled by `VerificationSettings` (`settings/model.py:143-231`) and instantiated at `settings/model.py:810-846`.

This is a mature LLM-as-judge for **agent outcome quality** (success likelihood + issue taxonomy). It is *not* wired as a prompt-comparison judge (e.g., "prompt A vs prompt B on dataset X"), has no prompt-version axis in its API, and no evaluation dataset is bundled to run it offline. `AgentFinishedCritic` and `EmptyPatchCritic` offer cheaper deterministic alternatives but likewise lack dataset integration.

**4. Are regressions caught before deployment?**
No. The CI pipeline (`.github/workflows/py-tests.yml:18-70`, `.github/workflows/lint.yml:18-75`) compiles TypeScript, lints Python, and runs ~100+ unit test modules with `pytest --forked -n auto`, but contains no job for prompt snapshot diff, no eval-dataset run, no critic-on-gold-traces run, and no required check that blocks a PRchanging a `.j2` prompt. Evaluation is explicitly externalized to `OpenHands/benchmarks` per `CONTRIBUTING.md:39` and `.dockerignore:22` / `.gitignore:191`, so even the external harness is not a pre-merge gate in this repo. A prompt regression would be caught only via manual review or post-deploy user outcome drift.

## Architectural Decisions

| Decision | Location | Rationale / Tradeoff |
|----------|----------|----------------------|
| Jinja2 templates for all prompts with conditional includes (`security_policy.j2`, `model_specific/*`) | `_sdk_inspect/sdk/context/prompts/prompt.py:57-88`, `_sdk_inspect/sdk/agent/prompts/system_prompt.j2:74-88,133-147` | Clean separation of static system prompt vs dynamic context (`static_system_message` vs `dynamic_context` in `base.py:368-432`) enables cross-conversation cache sharing; but makes prompt diffs easy to miss without golden tests. |
| `lru_cache` on `_get_env` (64) and `_get_template` (256) with `FileSystemBytecodeCache` | `prompt.py:57-85` | Avoids re-parsing across processes; hides that prompt edits require cache invalidation in long-lived servers (mtime check via `uptodate()` mitigates). |
| `FlexibleFileSystemLoader` supporting absolute-path templates (`system_prompt_filename` can be `/abs/path.j2`) | `prompt.py:16-45`, `base.py:159-166` | Supports custom prompt experimentation per agent without forking prompt dir; no registry or versioning of custom prompts. |
| `system_prompt` (inline string) vs `system_prompt_filename` mutual exclusion | `base.py:146-204` | Forces a single prompt source; docstring warns inline override "will override built-in instructions" but no test guards the warning. |
| Evaluation externalized to `OpenHands/benchmarks` + `evaluation/` gitignored | `.gitignore:191-202`, `CONTRIBUTING.md:57` | Keeps repo lean and avoids large datasets; cost is complete absence of in-repo eval datasets / prompt regression corpus. |
| Critic as remote `/classify` classifier (vLLM) vs inline LLM call | `critic/impl/api/client.py:74-299` | Shares inference infra with LLM proxy, gets calibrated multi-label probs; introduces network dependency, 300s timeout, and requires `api_key` — disabled by default (`critic_enabled=false` in `settings/model.py:143`). |
| `has_success_label` + 27-label `all_labels` tuple hardcoded on client | `client.py:128-198` | Locks prompt-evaluation taxonomy to a single model; `len(probs) != len(all_labels)` raises (`client.py:324`), so label churn breaks compatibility. |
| Two critic modes: `finish_and_message` (default, cheap) vs `all_actions` (expensive) | `critic/base.py:62-70`, `critic_mixin.py:34-40` | Lets users trade cost vs signal; but `all_actions` path is rarely tested and multiplies `/classify` calls linearly. |

## Notable Patterns

- **Static vs dynamic system message split** (`_sdk_inspect/sdk/agent/base.py:368-432`) — `static_system_message` (Jinja render, cacheable) vs `dynamic_context` (`AgentContext.get_system_message_suffix()`) as second content block with no cache marker. This is cache-aware prompt engineering.
- **Discriminated-union Pydantic models** for `CriticBase` / `AgentBase` (`DiscriminatedUnionMixin`) enabling polymorphic deserialization of critics and agents by `kind` field.
- **MCP config encryption via `model_validator`/`model_serializer`** (`base.py:206-309`) — `encrypted_mcp_config` handling inside prompt/agent base keeps secrets from leaking via prompt logs.
- **Rich visualization of judge output** (`critic/result.py:52-148`) — `CriticResult.visualize` renders star rating + categorized features with thresholded styling, surfaced in both CLI and frontend `CriticResultDisplay` (`frontend/src/components/v1/chat/event-message-components/critic-result-display.tsx:156`).
- **Taxonomy post-processing** (`taxonomy.py:82-180`) — softmax-normalized sentiment + `display_threshold=0.2` filtering + sorted-by-prob descending. This is a consistent rubric for downstream refinement prompt generation.

## Tradeoffs

- **Clean prompt abstraction vs zero verification**: The Jinja/loader/cache design is well-factored but has no property-based or snapshot tests, so prompt authoring velocity is unguarded.
- **Externalized evaluation vs in-repo confidence**: Outsourcing to `benchmarks` avoids dataset bloat but severs the feedback loop; CI in this repo cannot answer "did my prompt change degrade SWE-Bench pass rate?"
- **Remote classifier judge vs inline judge**: Offloading to a hosted `critic` model centralizes quality tracking but makes judging unavailable offline, adds latency, and requires credentialing — hence it defaults off (`settings/model.py:143`). A non-network `AgentFinishedCritic` exists but is heuristic-only (patch non-empty + FinishAction).
- **27-label taxonomy expressiveness vs brittleness**: Rich behavioral labels enable fine-grained refinement triggers (`critic.py:135-144` — `issue_threshold=0.75`), but any server-side label change must stay exactly in sync with `client.all_labels` or `extract_prob_map` throws.
- **Iterative refinement automation vs user surprise**: `enable_iterative_refinement` + `should_refine()` automatically retries the agent with a judge-generated follow-up prompt (`critic_mixin.py:114-137`). When enabled it improves task completeness but consumes extra LLM calls with no explicit budget cap beyond `max_iterations`.

## Failure Modes / Edge Cases

- **Silent prompt drift**: Editing `system_prompt.j2` or any `model_specific/*.j2` passes `make build` and `pre-commit` (formatting/type only) and all unit tests. No contract test asserts prompt content includes critical guardrails (e.g., security policy, browser tool flow, AI disclosure). A deleted `<SECURITY>` block would ship unnoticed.
- **Cache staleness**: `FileSystemBytecodeCache` with `lru_cache` relies on `mtime` check (`prompt.py:39-43`). Concurrent edits or clock skew in multi-process deployments could serve stale bytecode; no explicit `render_template` integration test validates cache invalidation.
- **Absolute-path prompt escapes sandbox**: `system_prompt_filename` can be an absolute path anywhere on the filesystem (`prompt.py:105-111`); if attacker-controlled `agent_context` or settings are persisted, they could point to sensitive files — validated only by `FileNotFoundError`, not allow-list.
- **Critic key missing/redacted**: `validate_secret` in `client.py:172-183` + `_get_api_key_value()` raises if `api_key` is `None`/empty/redacted, but `build_critic()` in `settings/model.py:813` silently returns `None` when `api_key` missing — user sees no error, only disabled judging.
- **Label-space mismatch crash**: `extract_prob_map` (`client.py:324-328`) raises `ValueError("len(probs) != len(all_labels)")` if server adds/removes labels. Deployment requires coordinated client/server rollout; no compatibility fallback.
- **Critic retry only on 500**: `classify_trace` retries 3× with exponential backoff (`client.py:276-298`) only for HTTP 500. Timeouts, 429 rate limits, or 502/503 proxy errors bubble immediately despite `timeout_seconds=300`.
- **All-actions mode cost explosion**: `mode=all_actions` (`base.py:62`) calls `/classify` after every tool use; on long trajectories (50+ steps) this multiplies latency/cost with no batching or sampling.
- **Jinja undefined vars silently empty**: `Environment(autoescape=False)` with no `StrictUndefined` (`prompt.py:67-74`) renders missing variables as empty string — typos in `system_prompt_kwargs` (e.g., `cli_mod` vs `cli_mode`) produce degraded prompts with no error.
- **Missing `SystemPromptEvent` / empty tools crash at runtime**: `APIBasedCritic.evaluate()` (`critic.py:74-82`) raises `ValueError` if no `SystemPromptEvent` or empty tool list — not caught at agent-construction time, only on first evaluation call.

## Future Considerations

- **Add prompt golden tests**: Snapshot-test `render_template(prompt_dir, "system_prompt.j2", **cases)` for at least the default kwargs, `cli_mode=true`, `enable_browser=true`, each `model_family` variant, and `security_policy_filename=""` edge case; store as `tests/unit/prompts/test_system_prompt_golden.py` and gate on CI.
- **Checked-in minimal eval corpus**: Vendor a small prompt-regression dataset (e.g., 20 curated user tasks with expected tool-use traces or `CriticResult.score` thresholds) under `tests/data/prompts/` — do not require the full `evaluation/` gitignored datasets — and run `APIBasedCritic` or `AgentFinishedCritic` offline against mocked `/classify` or deterministic mock critic in CI.
- **Prompt experiment tracker**: Introduce a `prompt_id` / `prompt_version` attribute on `AgentBase` (persisted in conversation state `state.py:87-101`) and optionally emit `prompt_tokens` + `critic.score` to analytics (`analytics_service.py:233`) keyed by `(prompt_version, model)`, enabling lightweight A/B comparisons without a full experiment platform.
- **CI prompt diff job**: A `.github/workflows/prompt-regression.yml` that fails when any `*.j2` under `prompts/` or `integrations/templates/` changes without a matching updated snapshot, analogous to frontend `check-translation-completeness`.
- **Critic offline harness**: Extend `critic/impl/pass_critic.py` pattern to a `PromptRegressionCritic` that loads `tests/data/prompts/golden_traces.json` and asserts judge scores within tolerance, removing the remote `/classify` dependency for PR-time checks.
- **Strengthen template safety**: Switch to `jinja2.StrictUndefined` for prompt env, add explicit allow-list for `system_prompt_filename` absolute paths, and test the `lmnr` tracing integration for prompt version stamping if LMNR is retained.

## Questions / Gaps

- **No evidence found: in-repo eval datasets** — Glob `evaluation/**/*` returns zero files; `.gitignore:191` confirms exclusion. Checked `pyproject.toml`, `CONTRIBUTING.md`, `Development.md`, `CREDITS.md`. Verdict: evaluation lives in external `OpenHands/benchmarks` repo, not searchable here.
- **No evidence found: prompt experiment tracker or A/B harness** — Grep for `ab_test|experiment_tracker|prompt.*experiment|eval.*dataset` and dependency `langfuse/langsmith/wandb/mlflow` yields only `poetry.lock` lockfile entries; no code imports or config.
- **No evidence found: prompt regression guard in CI** — `lint.yml` and `py-tests.yml` contain no prompt-specific jobs; `pytest.ini` runs only warning suppression. `tests/unit` grep for `prompt` shows no system-prompt regression suite.
- **No evidence found: LLM-as-judge prompt collection** — The Critic judge prompts/labels are model-side (tied to `Qwen/Qwen3-4B-Instruct-2507` tokenizer `client.py:111`). No judge-prompt `.j2` templates or `llm_as_judge` prompt files exist in tree.
- **Partially evidenced: LLM observability provider** — `lmnr>=0.7.20` and `langfuse==2.59.7` appear in `poetry.lock:4269-4271` / `pyproject.toml:259-260` but no application code references them for prompt version tracking; scope of their use is unverified without runtime config.
- **Gap**: How prompt changes are reviewed/approved across SDK vs server repos (`README.md:49-54` notes code moving to `software-agent-sdk` and `agent-canvas`) — prompt files live under `_sdk_inspect/sdk/` which is a staged import, so version drift between SDK release `1.37.1` and server pins is not tested here.

---

Generated by `Dimension 12.03: Prompt Evaluation and Experiments` against `openhands`.
