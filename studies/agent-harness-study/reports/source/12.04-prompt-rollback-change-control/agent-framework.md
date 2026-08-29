# Source Analysis: agent-framework

## Dimension 12.04: Prompt Rollback and Change Control

### Source Info

| Field | Value |
|-------|-------|
| Name | agent-framework |
| Path | `studies/agent-harness-study/sources/agent-framework` |
| Language / Stack | Python + .NET (monorepo); Python core analyzed (`python/packages/core`, `python/packages/declarative`) |
| Analyzed | 2026-08-29 |

## Summary

`agent-framework` is a client-side SDK, not a hosted prompt-management system. "Prompts" are plain `instructions: str` constructor arguments (`python/packages/core/agent_framework/_agents.py:669,744,759`) and declarative YAML fields (`PromptAgent.instructions` in `python/packages/declarative/agent_framework_declarative/_models.py:835-860`). There is no prompt registry, no versioned prompt store, no deployment pipeline for prompts distinct from code, no rollback API, no outcome-to-version association, and no drift monitoring. The only controls are generic code-review and release workflows applied to all source changes. A bad prompt change can only be reverted via `git revert` and a new package release / consumer redeploy — not in under a minute. For the library's intended scope this delegation is intentional, but against Dimension 12.04's operational criteria (reviewed/observable/reversible prompt changes linked to production outcomes) the capability is absent.

## Rating

**2 / 10 — Absent / ad-hoc**

Rationale: Prompt = in-code string; no first-class prompt versioning, deployment, rollback, or drift observability exists in the framework. Changes piggy-back on normal code PR/release flow (which validates code correctness, not prompt behavior) and revert requires the full code-release cycle. Telemetry captures `instructions` text when sensitive data is enabled but does not version, hash, or correlate it to outcomes; no drift detection or outcome-prompt linkage subsystem was found after searching `python/`, `.github/workflows`, `docs/decisions`, and `python/packages/declarative`.

## Evidence Collected

Every entry MUST include a file path with line numbers. Format: `path/to/file.ts:NN`.

| Area | Evidence | File:Line |
|------|----------|-----------|
| Review workflow (generic code) | PR template requires Motivation/Context, Description, Related Issue and Contribution Checklist; no prompt-specific review checklist. | `studies/agent-harness-study/sources/agent-framework/.github/pull_request_template.md:1-43` |
| Review workflow — prompt-skill guidance | PR skill delegates to generic template; no prompt-review stage mentioned. | `studies/agent-harness-study/sources/agent-framework/.github/skills/pull-requests/SKILL.md:18-32` |
| Review workflow — code owners | CODEOWNERS only covers `azurefunctions` and `durabletask` packages; no prompt owners or prompt-path ownership. | `studies/agent-harness-study/sources/agent-framework/.github/CODEOWNERS:1-7` |
| Review workflow — CI gate | Branch protection via `merge-gatekeeper` job (`pull_request` trigger) and `python-tests` matrix; checks validate code/tests, not prompt semantics or canary. | `studies/agent-harness-study/sources/agent-framework/.github/workflows/merge-gatekeeper.yml:1-14`, `studies/agent-harness-study/sources/agent-framework/.github/workflows/python-tests.yml:1-22` |
| Deployment — prompt is code | `Agent` stores `instructions` as plain string in `default_options["instructions"]`; `RawAgent.__init__` pops `instructions` from `default_options` or named param, no version/tag. | `studies/agent-harness-study/sources/agent-framework/python/packages/core/agent_framework/_agents.py:669-770` (esp. `669`, `742-759`) |
| Deployment — merge semantics | `_merge_options` concatenates `instructions` with `\n` on override — runtime wins, no hash/diff/version tracking. | `studies/agent-harness-study/sources/agent-framework/python/packages/core/agent_framework/_agents.py:125-130` |
| Deployment — chat options type | `ChatOptions` TypedDict defines `instructions: str` with no version, hash, or deployment metadata. | `studies/agent-harness-study/sources/agent-framework/python/packages/core/agent_framework/_types.py:3412-3413` |
| Deployment — declarative prompt | `PromptAgent` YAML schema holds `instructions: str` and `additionalInstructions: str` evaluated via PowerFx (`_try_powerfx_eval`), loaded via `AgentDefinition.from_dict` dispatch on `kind=="Prompt"`; no version field or deployment API. | `studies/agent-harness-study/sources/agent-framework/python/packages/declarative/agent_framework_declarative/_models.py:512-526,820-860` |
| Deployment — lifecycle packaging | `MapPromptAgentConversion` and `Deployment APIs` mentioned only as experimental changelog entries; no stable declarative prompt-deployment pipeline found in `python/packages/declarative`. | `studies/agent-harness-study/sources/agent-framework/python/CHANGELOG.md:103` |
| Deployment — release pipeline | `python-release.yml` triggers on `release.published` tag `python-*`, builds `python/packages/<PACKAGE>` — whole-package build, not per-prompt deploy; `python-tests.yml` / `python-merge-tests.yml` run `poe test` — no prompt validation stage. | `studies/agent-harness-study/sources/agent-framework/.github/workflows/python-release.yml:1-62`, `studies/agent-harness-study/sources/agent-framework/.github/workflows/python-tests.yml:46-48` |
| Rollback — absence | `grep -r "rollback|revert.*prompt|prompt.*version"` across source returns zero implementation hits; rollback is only `git revert` + release. No API like `rollbackPrompt(version)` or version store found. | Search `studies/agent-harness-study/sources/agent-framework` for `rollback\|drift\|prompt.*version` — no evidence |
| Rollback — instructions mutable | `AgentLoopMiddleware` and `SessionContext.extend_instructions` inject additional instructions at runtime (`session_context.instructions` merged in `_prepare_run_context`), but prior value not snapshot/versioned. | `studies/agent-harness-study/sources/agent-framework/python/packages/core/agent_framework/_agents.py:1493-1499` |
| Outcome association — telemetry captures raw text only | `OtelAttr.SYSTEM_INSTRUCTIONS = "gen_ai.system_instructions"` defined; `ChatTelemetryLayer`/`AgentTelemetryLayer` emit it when `OBSERVABILITY_SETTINGS.SENSITIVE_DATA_ENABLED` is true. No prompt hash/version attribute emitted. | `studies/agent-harness-study/sources/agent-framework/python/packages/core/agent_framework/observability.py:238`, `_types.py` telemetry paths |
| Outcome association — evaluation | `evaluate_agent`/`evaluate_workflow` and `LocalEvaluator`/`Evaluator` score conversations but store `EvalScoreResult`/`CheckResult` against run, not against a prompt version; `EvalItem` carries no prompt-version field. | `studies/agent-harness-study/sources/agent-framework/python/packages/core/agent_framework/_evaluation.py:1-340` (types) + `studies/agent-harness-study/sources/agent-framework/python/packages/core/agent_framework/__init__.py:62-83` exports |
| Drift monitoring | No file implementing prompt-drift monitoring; `grep drift` only hits `schema drift` sample text (`python/samples/02-agents/compaction/tiktoken_tokenizer.py:57`) and unrelated FIDES/compaction docs; no drift detector, baseline store, or alert rule. | `studies/agent-harness-study/sources/agent-framework/python/samples/02-agents/compaction/tiktoken_tokenizer.py:57`, search `drift` across `agent-framework` |
| History persistence (outcomes stored without prompt) | `InMemoryHistoryProvider`/`FileHistoryProvider` and `AgentSession` persist conversation `Message`s; history does not persist the originating `instructions` version alongside messages. | `studies/agent-harness-study/sources/agent-framework/python/packages/core/agent_framework/_sessions.py:166-262` |
| Skills prompt template (closest to versioned prompt) | `SkillsProvider._create_instructions` renders `instruction_template` with `{skills}` placeholder; template is in-memory string, customizable per construction, not stored/versioned externally. | `studies/agent-harness-study/sources/agent-framework/python/packages/core/agent_framework/_skills.py:1943-2025` |
| Deployment mechanism (user responsibility) | README/declarative docs instruct consumers to create `Agent(client, instructions="...")` or load YAML — deployment is consumer code deploy, not framework-managed. | `studies/agent-harness-study/sources/agent-framework/python/packages/core/agent_framework/__init__.py:20-21` export; `studies/agent-harness-study/sources/agent-framework/python/packages/declarative/AGENTS.md:10-22` usage |

## Answers to Dimension Questions

**1. Can prompts be rolled back?**
No. Prompts are bare strings (`Agent(instructions=...)` or `PromptAgent.instructions`). The framework provides no prompt store, version history, or rollback API. Reverting requires a code `git revert` of the commit that changed the string and a new release/redeployment via `python-release.yml:16-58` (tag `python-*`) or consumer app redeploy. Not achievable in under a minute via framework primitives. Verified by absence of any `rollback`/`version.*prompt` code (`grep` across `python/` and `.github/`).

**2. Are prompt changes reviewed?**
Only via generic code review for all changes. There is a standard PR template (`.github/pull_request_template.md:1-43`), `merge-gatekeeper.yml:4-10` branch protection, and CI (`python-tests.yml:4-48`, `python-merge-tests.yml`). `CODEOWNERS:1-7` lists owners for two Python packages, not prompt paths. There is no prompt-specific review lane (no prompt-impact checklist, no behavioral eval gate, no designated prompt reviewer, no canary). A prompt edit is indistinguishable from any other string change in review.

**3. Can production issues be linked to prompt changes?**
Not within the framework. Observability can emit raw prompt text as `gen_ai.system_instructions` (`observability.py:238`) when `ENABLE_SENSITIVE_DATA=true`, but it does not emit a prompt hash, version ID, or deployment ID, nor does it correlate to evaluation/outcome stores. `EvalResults`/`AgentResponse` carry messages and `response_id` but not prompt-version metadata (`_types.py:3412-3413`, `_evaluation.py`). Correlation must be done externally (e.g., correlating git SHA to traces in the consumer's observability backend) — the framework itself retains no outcome→prompt-version link.

**4. Is prompt drift monitored?**
No evidence found. Search of `agent-framework` for `drift`, `prompt.*monitor`, `prompt.*baseline`, and inspection of `observability.py` (~1254 lines) and `_agents.py`/`_types.py` reveals no drift detector, no baseline prompt store, no scheduled re-evaluation job, and no alert. Compaction and FIDES docs mention "drift" only in unrelated schema/context-window senses. The `enable_sensitive_telemetry`/`OBSERVABILITY_SETTINGS` flags gate verbosity, not drift analysis.

## Architectural Decisions

* **Prompt-as-constructor-argument** — `Agent.__init__(instructions: str | None)` (`_agents.py:669`) keeps the SDK minimal and host-agnostic; tradeoff is zero lifecycle management. Good for library scope, bad for production prompt governance.
* **No prompt registry, pure code/YAML source-of-truth** — Declarative prompts live in YAML files parsed by `PromptAgent.from_dict` (`_models.py:820-861`); no database or version table. Means deployment is version-controlled but not hot-swappable.
* **String-concatenation merge for instructions** — `_merge_options` (`_agents.py:125-127`) merges `instructions` via `\n` join, allowing runtime overrides to append; no conflict detection or version-bumping semantics.
* **Telemetry opt-in for sensitive payload** (`observability.py:1075-1103` `enable_sensitive_telemetry`, `682-812` settings): prompt text is intentionally excluded by default to avoid leakage, which simultaneously eliminates the only built-in path to correlate prompts with traces.

## Notable Patterns

* **Generic code-review-as-prompt-review** — `CONTRIBUTING.md:108-158` and `.github/workflows/*` enforce build/test before merge, reused for prompts without specialization.
* **Declarative `kind: Prompt` dispatch** — `AgentDefinition.from_dict` (`_models.py:512-526`) and `agent_schema_dispatch` (`_models.py:1021-1031`) treat prompts as one `AgentDefinition` kind among others — uniform loading, no prompt-specific lifecycle hooks.
* **ContextProvider instruction injection** — `SessionContext.extend_instructions` (`_sessions.py:253-262`) and `Agent._prepare_run_context` (`_agents.py:1493-1499`) allow middleware/skills to append instructions, but as ephemeral in-memory lists cleared per run.
* **Skills prompt templating in-code** — `_skills.py:1683,1811,1944-1975` uses `prompt_template.format(skills=...)` — the only place resembling a prompt template engine, but not versioned or externally managed.

## Tradeoffs

* **Simplicity vs. governance**: Zero infrastructure for prompt versioning keeps the SDK lightweight and portable across hosts (OpenAI, Azure, Bedrock), but pushes all change-control responsibility to the consumer. Teams must build their own registry/canary/rollback.
* **Privacy vs. observability**: Gating `SYSTEM_INSTRUCTIONS` behind `SENSITIVE_DATA_ENABLED` prevents accidental prompt leakage in traces, at the cost of making production debugging of prompt regressions harder unless the consumer explicitly opts in.
* **Release atomicity vs. prompt velocity**: Bundling prompts with package releases (`python-release.yml:13-57`) guarantees code-prompt atomicity, but prevents independent, low-latency prompt iteration (a one-line instruction tweak still requires a full `poe build` + tag + release).
* **PowerFx evaluation for declarative instructions** (`_models.py:51-80` `_try_powerfx_eval`) enables dynamic instruction rendering (`=Env.VAR`), which increases flexibility but introduces an unmonitored surface for prompt drift without audit trail.

## Failure Modes / Edge Cases

* **Bad prompt ships to every consumer on upgrade** — Because `Agent(instructions=...)` is co-versioned with the library's model/tool bindings, a misphrased instruction in a sample or in `PromptAgent.instructions` cannot be hot-patched; consumers must pin to prior git tag/package version and redeploy.
* **Concatenated instructions silently merge** — `_merge_options` (`_agents.py:125-127`) and `merge_chat_options` (`_types.py:3658-3664`) append overrides without length/contract checks; truncating or contradictory instructions fail silently (LLM may ignore, no framework error).
* **Declarative YAML PowerFx failure falls back to raw value** — `_try_powerfx_eval` (`_models.py:70-80`) logs at `debug` and returns the raw `=expression` string on `powerfx` absence or eval error; the model then receives literal `=Env...` text, a subtle drift that is unmonitored.
* **Session persistence replay carries stale prompt** — `FileHistoryProvider`/`InMemoryHistoryProvider` replay stored `Message`s on next turn; if the prompt was changed between turns, the replayed context and new `instructions` may clash (no `instructions` version stored per session).
* **MCP server tool prompts reload without versioning** — `MCPTool.load_prompts` / `notifications/prompts/list_changed` (`_mcp.py:1261-1335`) reloads remote prompts on notification; a malicious/updated remote prompt propagates without review or rollback.
* **Observability gap masks regressions** — With `SENSITIVE_DATA_ENABLED=false` (default), `gen_ai.system_instructions` not emitted; a prompt regression's production impact cannot be correlated from traces alone.

## Future Considerations

* **Introduce content-addressable prompt versioning** — Hash `instructions` (or full prompt template + variables) and surface as `gen_ai.prompt.version` OTEL attribute alongside `SYSTEM_INSTRUCTIONS`; persist the hash in `AgentSession` history so replays are reproducible and bisectable. Keep hashing opt-in to preserve privacy while enabling correlation.
* **Add prompt-scoped deployment lane** — Separate from `python-release.yml` tag release, add a `prompt-registry` (JSON/YAML store with semver per `PromptAgent.name`) + `python/workflows/prompt-deploy.yml` that validates prompts via `evaluate_agent` checks (`_evaluation.py`) and allows `rollback(prompt_name, version)` in <60s without code release.
* **Treat prompt changes as distinct PR category** — Extend `.github/pull_request_template.md` and `CODEOWNERS` to carve a `**/prompts/**` or `instructions:` diff label that triggers required prompt reviewers and a `prompt-eval` CI job (keyword/LLM-judge rubric from `_evaluation.py`) before merge.
* **Emit prompt drift signal** — Add a periodic `workflow` that re-runs a small golden-set (`python/samples/05-end-to-end/evaluation/...`) against the current prompt and alerts when `EvalScoreResult` drops > threshold; export as `gen_ai.prompt.drift` metric via `get_meter()` (`observability.py:1000-1050`).
* **Snapshot instructions per run** — Persist `default_options["instructions"]` snapshot into `SessionContext`/`AgentResponse.additional_properties` so post-mortems can answer "which prompt produced this trace?" without enabling full sensitive-payload capture.

## Questions / Gaps

* No evidence of a hosted prompt store or external registry (Foundry/AI Search) being used as the live prompt source-of-truth within this repo — declarative YAMLs live in-sample (`declarative-agents/workflow-samples/*.yaml`) but sample `instructions` not versioned. Is external prompt hosting intentionally out-of-scope?
* Is MCP `load_prompts` intended to be the extensible prompt-registry path? The code treats MCP prompts as tools (`_mcp.py:133-153,360-438`) with no signing/version metadata — clarification on trust model needed.
* How should `Sensitive Data` gating (`OBSERVABILITY_SETTINGS.SENSITIVE_DATA_ENABLED` in `observability.py:806-812`) coexist with required prompt-outcome correlation? Current design forces an either/or between privacy and debuggability.
* No spec for prompt canary/staged rollout was found in `docs/decisions/*`; is this roadmap for `.NET` parity or deferred to consumer infra?

---

Generated by `Dimension 12.04: Prompt Rollback and Change Control` against `agent-framework`.
