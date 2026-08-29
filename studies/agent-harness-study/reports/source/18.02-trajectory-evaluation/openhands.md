# Source Analysis: openhands

## 18.02 Trajectory Evaluation

### Source Info

| Field | Value |
|-------|-------|
| Name | openhands |
| Path | `studies/agent-harness-study/sources/openhands` |
| Language / Stack | Python (SDK `openhands-sdk`), React frontend, Pydantic, litellm |
| Analyzed | 2026-08-29 |

## Summary

OpenHands does **not** ship a trajectory-evaluation harness inside this repository. Evaluation benchmarks are delegated to the external `OpenHands/benchmarks` repository (referenced in `Development.md:330` and `.gitignore:191` ignoring `evaluation/`). Within the source, the closest construct is the **Critic** subsystem (`_sdk_inspect/sdk/critic/`): a pluggable classifier that scores a prefix of the agent's trajectory (full event history) and emits a `CriticResult` attached to `ActionEvent`/`MessageEvent`. In `finish_and_message` mode it evaluates only terminal events; in `all_actions` mode it evaluates after every tool call, providing a coarse per-step trajectory score. The remote `APIBasedCritic` maps the whole trajectory to ~22 label probabilities (behavioral issues including `improper_tool_use_or_setup` and `loop_behavior`, follow-up patterns, infra) via `taxonomy.py:8-39` and uses the `success` probability as `CriticResult.score`. Token/context usage is instrumented via `Metrics`/`TokenUsage` but never scored for quality, and recovery is executed (iterative refinement + corrective nudge) but not measured as an evaluation dimension. Trajectory persistence is limited to JSON export (`config.template.toml:29`, `event_service_base.py:148`) and per-event `critic_result` fields — no aggregated trajectory scoring, dashboard, or dedicated step/tool-choice/context/recovery evaluators were found.

## Rating

**5 / 10 — Present but inconsistent, weakly documented, fragile**

Rationale: Core scoring primitive (`CriticBase.evaluate(events, git_patch) -> CriticResult`) is well-typed, extensible, and operational (4 implementations, retry, threshold, `all_actions` per-step mode). However trajectory evaluation as a dimension is fragmented: intermediate-step scoring is opt-in and holistic (single score per prefix, not per-dimension), tool-choice quality is proxied by a single LLM label (`improper_tool_use_or_setup`) with no deterministic comparator, context usage is counted not evaluated, recovery has no eval cases, and there are no in-repo trajectory-level tests or benchmarks — all benchmarks live externally.

## Evidence Collected

Every entry MUST include a file path with line numbers. Format: `path/to/file.ts:NN`.

| Area | Evidence | File:Line |
|------|----------|-----------|
| Step-by-step eval code | `CriticBase.evaluate(events, git_patch)` abstract method takes full trajectory `Sequence[LLMConvertibleEvent]` and returns `CriticResult` — foundation for step-aware scoring | `_sdk_inspect/sdk/critic/base.py:83` |
| Step-by-step eval mode | `mode` field controls granularity: `finish_and_message` (default, terminal only) vs `all_actions` (every action, warns "significantly slower") | `_sdk_inspect/sdk/critic/base.py:62` |
| Step-by-step eval gating | `_should_evaluate_with_critic(action)` returns `True` for every action if `mode==all_actions`, else only `FinishAction` | `_sdk_inspect/sdk/agent/critic_mixin.py:34` |
| Step-by-step eval execution | Per-action critic attachment in action path: `if _should_evaluate: critic_result = _evaluate_with_critic(); action_event = action_event.model_copy(update={"critic_result": critic_result})` | `_sdk_inspect/sdk/agent/agent.py:893` |
| Step-by-step eval execution (messages) | `MessageEvent` path also runs critic when `mode==finish_and_message` via `_emit_message_event` | `_sdk_inspect/sdk/agent/response_dispatch.py:261` |
| Step-by-step eval implementation | `_evaluate_with_critic` builds `events = list(conversation.state.events) + [event]` filtered to `LLMConvertibleEvent` and calls `critic.evaluate(events, git_patch=None)` — git patch currently unused | `_sdk_inspect/sdk/agent/critic_mixin.py:49` |
| Tool choice evaluator (proxy) | Taxonomy label `improper_tool_use_or_setup` under `agent_behavioral_issues` — exposed as high-prob signal for refinement and visualization | `_sdk_inspect/sdk/critic/impl/api/taxonomy.py:17` |
| Tool choice evaluator (loop) | `loop_behavior` label captures repetitive tool-use pathology | `_sdk_inspect/sdk/critic/impl/api/taxonomy.py:18` |
| Context usage metrics | `TokenUsage` model + `Metrics.accumulated_token_usage` + `TokenEvent` track prompt/completion tokens per call | `_sdk_inspect/sdk/llm/utils/metrics.py:35` |
| Context usage metrics | `Metrics` comment "Token each turn for calculating context usage" — metrics infra exists but not evaluative | `_sdk_inspect/sdk/llm/utils/metrics.py:171` |
| Recovery eval — iterative refinement config | `IterativeRefinementConfig(success_threshold=0.6, max_iterations=3)` drives auto-retry when `score < threshold` | `_sdk_inspect/sdk/critic/base.py:20` |
| Recovery eval — refinement loop | `_check_iterative_refinement` checks `critic_result` against threshold and taxonomy `issue_threshold=0.75`, increments `iterative_refinement_iteration` and injects `followup` `MessageEvent` instead of finishing | `_sdk_inspect/sdk/agent/critic_mixin.py:76` |
| Recovery eval — holistic refinement decision | `APIBasedCritic.should_refine` augments score check with high-prob `agent_behavioral_issues` via `_get_high_probability_agent_issues` | `_sdk_inspect/sdk/critic/impl/api/critic.py:135` |
| Recovery eval — corrective nudge | `_send_corrective_nudge` injects user-role message on `EMPTY`/`REASONING_ONLY` responses to recover stuck loop | `_sdk_inspect/sdk/agent/response_dispatch.py:283` |
| Trajectory scoring model | `CriticResult.score: float [0,1]` with `THRESHOLD=0.5`, `success = score >= 0.5`, star rating + categorized_features display | `_sdk_inspect/sdk/critic/result.py:7` |
| Trajectory scoring persistence | `ActionEvent.critic_result: CriticResult \| None` and `MessageEvent.critic_result` fields persist score on each event | `_sdk_inspect/sdk/event/llm_convertible/action.py:69` |
| Trajectory scoring persistence | Same field on message events | `_sdk_inspect/sdk/event/llm_convertible/message.py:55` |
| Trajectory scoring visualization | `CriticResult.visualize` renders stars, percentage, and per-category issue lists; displayed in UI components | `_sdk_inspect/sdk/critic/result.py:52` |
| Trajectory scoring UI | Frontend `CriticResultDisplay` component renders per-event critic result | `frontend/src/components/v1/chat/event-message-components/critic-result-display.tsx:151` |
| Trajectory scoring taxonomy | `categorize_features(probs, display_threshold=0.2)` groups raw probs into `agent_behavioral_issues`, `user_followup_patterns`, `infrastructure_issues`, `other`; success score driving trajectory quality | `_sdk_inspect/sdk/critic/impl/api/taxonomy.py:82` |
| Trajectory scoring — API client | `CriticClient.classify_trace(messages, tools)` posts trajectory to `POST /classify` with retry (3x on 500) and extracts `LabelProbMap`; success prob becomes trajectory score | `_sdk_inspect/sdk/critic/impl/api/client.py:262` |
| Trajectory export (not scoring) | `save_trajectory_path` config and `Iterate all events once in timestamp order for trajectory export` — trajectory is exportable but not scored in harness | `config.template.toml:29` |
| Trajectory export (not scoring) | `event_service_base.py:148` trajectory export method | `openhands/app_server/event/event_service_base.py:148` |
| Verification settings | `VerificationSettings.critic_enabled/mode/threshold/max_refinement_iterations` expose critic & refinement knobs to settings UI/API | `_sdk_inspect/sdk/settings/model.py:139` |
| Critic implementations | `AgentFinishedCritic` (FinishAction + non-empty patch), `EmptyPatchCritic`, `PassCritic`, `APIBasedCritic` — multiple evaluators but all holistic | `_sdk_inspect/sdk/critic/impl/agent_finished.py:24` |
| External evaluation delegation | README/contributing point evaluation to external repo `OpenHands/benchmarks` — no in-repo harness | `Development.md:330` |
| No evidence — step evaluators | No dedicated per-step tool-choice comparator, context-efficiency scorer, or recovery test suite found in `tests/unit/` | `tests/unit/` — search returned no matches |
| Trajectory test coverage | Only frontend display tests for `CriticResultDisplay`; no backend unit tests for `CriticMixin` or `APIBasedCritic` in `tests/unit` | `frontend/__tests__/components/v1/chat/event-message-components/critic-result-display.test.tsx:5` |

## Answers to Dimension Questions

**1. Are intermediate steps evaluated?**
Partially. The infrastructure supports it but it is opt-in and coarse. `CriticBase.mode="all_actions"` (`_sdk_inspect/sdk/critic/base.py:62`) causes `CriticMixin._evaluate_with_critic` to run after every `ActionEvent` (`_sdk_inspect/sdk/agent/agent.py:893`) and `ResponseDispatchMixin._emit_message_event` after every agent message (`_sdk_inspect/sdk/agent/response_dispatch.py:274`). However each evaluation scores the **entire prefix** (`events + current_event`) as a single `success` probability (`_sdk_inspect/sdk/critic/impl/api/critic.py:113`), not the quality of the current step in isolation. Default `finish_and_message` mode evaluates only terminal steps, so intermediate reasoning is invisible unless the operator explicitly enables `all_actions` (flagged "WARNING: significantly slower" due to per-step API calls).

**2. Is tool selection quality measured?**
Only indirectly via an LLM classifier label. The `APIBasedCritic` taxonomy includes `improper_tool_use_or_setup` (`_sdk_inspect/sdk/critic/impl/api/taxonomy.py:17`) and `loop_behavior` (`:18`) under `agent_behavioral_issues`, surfaced when `prob >= 0.2` (display threshold) and triggering iterative refinement when `prob >= 0.75` (`_sdk_inspect/sdk/critic/impl/api/critic.py:48`). There is no deterministic evaluator that compares the chosen tool against an expected tool set, measures alternative ranking, or scores tool arguments. The signal is probabilistic, prompt-dependent, and conflates tool misuse with general behavioral issues.

**3. Is context usage evaluated?**
No. Token usage is instrumented (`_sdk_inspect/sdk/llm/utils/metrics.py:35`, `Metrics.accumulated_token_usage:90`) and condenser events manage window overflow (`_sdk_inspect/sdk/context/condenser/llm_summarizing_condenser.py`), but no evaluator scores whether context was used efficiently (e.g., retrieval precision, condensation quality, window utilization). The critic taxonomy has no context-related label; `Metrics` is reported for billing/observability, not for trajectory quality. Outcome: volume is tracked, quality is not.

**4. Is recovery behavior measured?**
Execution supports recovery but evaluation does not measure it. Iterative refinement (`_sdk_inspect/sdk/critic/base.py:20`, `_sdk_inspect/sdk/agent/critic_mixin.py:76`) retries a task up to `max_iterations` when `score < threshold`, and `response_dispatch.py:283` injects a corrective nudge on empty reasoning. These mechanisms are logged (`logger.info "Iterative refinement: continuing ..."`) but never scored: there are no eval cases that inject failures and assert recovery, no metric like `recovery_rate` or `steps_to_recovery`, and no test suite that exercises `_check_iterative_refinement` with fault injection.

## Architectural Decisions

| Decision | Location | Implication |
|----------|----------|-------------|
| Holistic trajectory classifier rather than stepwise rubric | `_sdk_inspect/sdk/critic/impl/api/critic.py:58` (`evaluate(events)` takes full history) + `client.py:262` (`classify_trace` posts full `formatted_messages`) | Single `success` score conflates tool, reasoning, and context errors; easy to deploy but insensitive to which step failed. |
| `LLMConvertibleEvent` barrier | `critic_mixin.py:59` filters to `LLMConvertibleEvent` before evaluation | Non-LLM events (observations, errors) excluded from scoring — tool output quality not directly evaluated. |
| Opt-in per-step scoring | `VerificationSettings.critic_mode` (`_sdk_inspect/sdk/settings/model.py:153`) + `CriticBase.mode` | Default trajectory appears unscored; operators must accept 10-100x API cost for `all_actions`. |
| Remote vLLM `/classify` evaluator | `_sdk_inspect/sdk/critic/impl/api/client.py:74` (`DEFAULT_CRITIC_SERVER_URL`) + `classify_trace:262` | Trajectory quality depends on external service availability; retry is limited to 3x on 500 (`client.py:276`), no offline fallback. |
| Iterative refinement as primary recovery consumer | `agent.py:221` + `critic_mixin.py:76` + `response_dispatch.py:261` | Critic score is actuator, not just reporter — blurs evaluation vs control. |
| Empty `git_patch=None` in live evaluation | `critic_mixin.py:64` hardcodes `git_patch=None`; `AgentFinishedCritic` and `EmptyPatchCritic` require patch but are not wired live | Patch-based critics (`agent_finished.py:48`, `empty_patch.py:29`) are dead code during conversation, only useful in offline benchmark scripts. |
| Display threshold decoupling | `taxonomy.py:84` (`display_threshold=0.2`) vs `result.py:11` (`DISPLAY_THRESHOLD=0.2`) vs `client.should_refine issue_threshold=0.75` | What is shown ≠ what triggers recovery; high-severity issues may be hidden if below display threshold, low-severity issues never trigger retry. |

## Notable Patterns

- **Prefix-scored trajectory**: Every critic call re-scores the entire conversation history via `View.from_events(events)` (`_sdk_inspect/sdk/critic/impl/api/critic.py:85`), not delta quality. Enables detecting drift but masks step-level credit assignment.
- **Taxonomy-driven visualization**: Raw logits → `probs` → `categorize_features` (`taxonomy.py:82`) → `CriticResult.metadata.categorized_features` → `result.py:52` star/percentage UI — rich signal downgraded to single `score` for control.
- **Two-mode critique**: Minimal (`finish_and_message`) vs exhaustive (`all_actions`) — lets users trade cost for trajectory visibility.
- **Observability hook integration**: Critic scores attached to `ActionEvent`/`MessageEvent` are persisted in event log and exportable (`config.template.toml:29`), making offline trajectory analysis possible without dedicated harness.
- **Retry + nudge resilience**: Condenser-triggered history repair (`agent.py:543`), LLM retry (`llm.py:799`), critic retry (`client.py:276`), and corrective nudge (`response_dispatch.py:283`) compose layered recovery, but none emit trajectory-quality metrics.

## Tradeoffs

| Tradeoff | Pro | Con |
|----------|-----|-----|
| Single `success` probability | Simple threshold logic, comparable across trajectories | Cannot answer "was step 7 good even though step 12 failed?" — holistic score hides intermediate quality |
| Remote classifier | Can leverage large tuned critic model without bundling | Adds latency/cost per evaluation; `all_actions` mode multiplies API calls linearly with trajectory length; offline/CI unfriendly |
| `git_patch=None` live | Simplifies live path, avoids diff generation cost | `AgentFinishedCritic`/`EmptyPatchCritic` (the only non-LLM critics) are never exercised during conversation |
| Opt-in scoring | Low overhead by default | Most trajectories have zero intermediate scores — evaluation coverage depends on operator opt-in, not guarantee |
| Label space fixed in client | Stable API contract (`all_labels` length checked `client.py:324`) | Adding/removing labels requires synchronized server+client bump; fragile version coupling |
| UI-level filtering (`display_threshold`) | Reduces noise | Operator may not see `improper_tool_use_or_setup=0.18` even though it indicates tool misuse — evaluation signal suppressed |

## Failure Modes / Edge Cases

- **Missing `SystemPromptEvent`**: `APIBasedCritic.evaluate` raises `ValueError("SystemPromptEvent is required")` (`_sdk_inspect/sdk/critic/impl/api/critic.py:75`) — any trajectory missing the initial system prompt (e.g., corrupted replay) cannot be scored.
- **Critic service unavailability**: `CriticClient._post_with_retry` retries only on HTTP 500, 3 attempts (`client.py:276`); 429/timeout/network errors bubble as unhandled exception then caught and logged as `logger.error "Critic evaluation failed"` returning `None` (`critic_mixin.py:72`), leaving the event unscored and refinement disabled.
- **Empty API key**: `build_critic()` returns `None` if `llm.api_key is None` (`_sdk_inspect/sdk/settings/model.py:824`) — OSS/local LLM without key silently disables all trajectory scoring.
- **Threshold misconfiguration**: `IterativeRefinementConfig.success_threshold=0.6` default means ~60% success still triggers retry; operator setting 0.0 disables retry entirely, 1.0 retries forever until `max_iterations`.
- **History truncation interaction**: `APIBasedCritic.evaluate` re-creates `View.from_events(events).events` (`critic.py:85`) which applies condenser view — critic may see a condensed history different from what the agent saw, scoring a lossy projection.
- **Parallel tool calls**: `_ActionBatch` executes tools in parallel (`agent.py:454`) but critic scores each `ActionEvent` individually with the same conversation prefix — concurrent actions scored without causal ordering.
- **Silent `git_patch=None` failure**: `AgentFinishedCritic` always returns 0.0 when `git_patch` is None/empty (`agent_finished.py:49`) — if wired accidentally, every trajectory fails regardless of reasoning quality.
- **No persistence of trajectory score**: Only per-event `critic_result` persists; no aggregate trajectory record or time-series of scores — retrospective analysis requires re-reading entire event log.

## Future Considerations

- Wire `git_patch` generation into `critic_mixin._evaluate_with_critic` for `AgentFinishedCritic`/`EmptyPatchCritic` to become live trajectory validators (currently dead code path `_sdk_inspect/sdk/agent/critic_mixin.py:64`).
- Add deterministic step evaluators: tool-choice comparator (expected vs actual tool + arg schema), condensation fidelity scorer, and loop/repetition detector — complement the LLM classifier and remove sole dependence on `improper_tool_use_or_setup`.
- Introduce trajectory-level aggregation (`trajectory_score = f(per-step scores, final success)`) persisted as a `ConversationStats` field alongside `Metrics`, enabling dashboards and regression detection.
- Implement fault-injection recovery eval harness (e.g., force `AgentErrorEvent` + assert nudge/refinement restores progress) — would turn `response_dispatch.py:283` and `critic_mixin.py:76` from actuators into measurable recovery metrics.
- Ship in-repo benchmark stubs (not just `.gitignore`'d `evaluation/`) so CI can assert per-step scoring does not regress — currently depends on external `benchmarks` repo (`Development.md:330`).
- Expand `CriticClient` retry to include 429/timeout and add offline fallback (e.g., `PassCritic` or local heuristic) to avoid silent unscored trajectories.
- Surface `context usage quality` taxonomy label (e.g., `context_overflow`, `redundant_context`, `missing_context`) so condenser behavior becomes evaluable.

## Questions / Gaps

- No evidence found for trajectory scoring persistence beyond per-event `critic_result` — how are historical trajectory scores queried or trended?
- No tests for `CriticMixin._evaluate_with_critic` / `_check_iterative_refinement` in `tests/unit/` — what is the intended test strategy for trajectory evaluation?
- Is `APIBasedCritic` expected to be the sole trajectory evaluator, or are domain-specific critics (e.g., `EmptyPatchCritic`) meant to be composed? No composition operator exists in `CriticBase`.
- How does the critic handle ACP-subagent trajectories (`acp_agent.py`) — does scoring span sub-agent tool calls or only parent?
- What benchmark defines success vs intermediate quality — the external `OpenHands/benchmarks` repo is referenced but not version-pinned here; no contract in this source.
- Does `V1` event export include `critic_result` faithfully for offline replay? Checked `event_service_base.py:148` export but not schema verification.

---

Generated by `18.02-trajectory-evaluation` against `openhands`.
