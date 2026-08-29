# Source Analysis: agent-framework

## Dimension 12.03: Prompt Evaluation and Experiments

### Source Info

| Field | Value |
|-------|-------|
| Name | agent-framework |
| Path | `studies/agent-harness-study/sources/agent-framework` |
| Language / Stack | Python + .NET / Microsoft Agent Framework (asyncio, OpenAI evals API, Foundry) |
| Analyzed | 2026-08-29 |

## Summary

`agent-framework` is not a prompt-centric product but an agent/workflow harness. It provides a provider-agnostic evaluation framework in `python/packages/core/agent_framework/_evaluation.py:1` with a minimal `Evaluator` protocol, a concrete `LocalEvaluator`, and orchestration helpers `evaluate_agent`/`evaluate_workflow`. Cloud LLM-as-judge is delegated to `python/packages/foundry/agent_framework_foundry/_foundry_evals.py:728` (`FoundryEvals`) which wraps Azure AI Foundry built-in evaluators (relevance, coherence, task_adherence, tool_call_accuracy, etc.) via the OpenAI Evals API (`client.evals.create` / `runs.create` / `poll`). There are no versioned prompt templates, no persisted eval datasets, and no experiment tracker (MLflow/W&B/A-B service). Evaluation datasets are transient `EvalItem` lists built per-run. The framework ships CI-friendly assertion helpers (`EvalResults.raise_for_status:470`, `assert_score_at_least:502`, `assert_no_failed_items:616`) and thorough unit tests, but prompt-change regression is not automated in CI — integration depends on developers wiring `evaluate_agent` manually.

## Rating

**5/10 — Present but inconsistent, fragile for prompt confidence**

The evaluation substrate is clear and well-tested (local checks, Foundry LLM judges, workflow breakdown, rubric dimensions, `expected_output`/`expected_tool_calls` ground truth), but the four prompt-experiment questions score unevenly: LLM-as-judge is mature, prompt-change testing and experiment tracking are at best ad-hoc/manual, and CI regression is opt-in boilerplate rather than a gate. A prompt (`Agent.instructions`) change cannot be deployed with automated confidence that it won't regress without custom wiring; nothing versions prompts, versions datasets, or records A/B comparisons inside the repo.

## Evidence Collected

Every entry MUST include a file path with line numbers. Format: `path/to/file.ts:NN`.

| Area | Evidence | File:Line |
|------|----------|-----------|
| EvalItem provider-agnostic data model | `EvalItem` dataclass with `conversation` as single source of truth, `expected_output`, `expected_tool_calls`, `tools`, `context`, `split_strategy`; derived `query`/`response` via `ConversationSplitter` | `python/packages/core/agent_framework/_evaluation.py:182` |
| Conversation split strategies | `ConversationSplit.LAST_TURN` and `FULL` enum plus `ConversationSplitter` protocol; used by `EvalItem.split_messages()` and stamped on items via `conversation_split` param | `python/packages/core/agent_framework/_evaluation.py:78`, `python/packages/core/agent_framework/_evaluation.py:109`, `python/packages/core/agent_framework/_evaluation.py:1629` |
| LocalEvaluator (API-free checks) | `LocalEvaluator` implements `Evaluator`, runs `EvalCheck` functions per item; aggregates `EvalResults` with `per_evaluator`, `items`, `result_counts` | `python/packages/core/agent_framework/_evaluation.py:1518` |
| Built-in local checks | `keyword_check:1062`, `tool_called_check:1090`, `tool_calls_present:1173`, `tool_call_args_match:1215`; helper `_extract_tool_calls:1152` | `python/packages/core/agent_framework/_evaluation.py:1062` |
| `@evaluator` decorator + coercion | `@evaluator` wraps plain sync/async functions via signature introspection (`_resolve_function_args:1296`, `_coerce_result:1354`); supports `bool/float/dict/CheckResult` returns, parameter-name injection (`query`, `response`, `expected_output`, `conversation`, `tools`, `context`) | `python/packages/core/agent_framework/_evaluation.py:1415` |
| Ground-truth stamping | `evaluate_agent` validates and stamps `expected_output` and `expected_tool_calls` onto `EvalItem`s with modulo wrapping for `num_repetitions>1` | `python/packages/core/agent_framework/_evaluation.py:1809` |
| Workflow evaluation + per-agent breakdown | `evaluate_workflow:1832` extracts per-agent data via `_extract_agent_eval_data:941`, groups by executor_id, builds `sub_results` on `EvalResults` | `python/packages/core/agent_framework/_evaluation.py:1832` |
| EvalResults CI assertions | `raise_for_status:470`, `assert_score_at_least:502`, `assert_dimension_score_at_least:544`, `assert_no_failed_items:616`, `all_passed:455` with recursive sub-result checks | `python/packages/core/agent_framework/_evaluation.py:470` |
| Result types (auditing) | `EvalScoreResult:305`, `EvalItemResult:326` (with `input_text`, `output_text`, `token_usage`, `error_code`, `scores`, `dimensions`), `RubricScore:650`, `AgentEvalConverter.convert_message:742` / `convert_messages:832` | `python/packages/core/agent_framework/_evaluation.py:305` |
| Foundry LLM-as-judge provider | `FoundryEvals` class with `evaluate:864` building JSONL `data_source` and testing criteria; constants for 15+ built-ins (`RELEVANCE:821`, `TOOL_CALL_ACCURACY:811`, etc.) | `python/packages/foundry/agent_framework_foundry/_foundry_evals.py:728` |
| Foundry evaluator taxonomy | `_BUILTIN_EVALUATORS:142` (16 entries), `_AGENT_EVALUATORS:107` (conversation-array), `_TOOL_EVALUATORS:120`, `_GROUND_TRUTH_EVALUATORS:138`, defaults `_DEFAULT_EVALUATORS:169` | `python/packages/foundry/agent_framework_foundry/_foundry_evals.py:142` |
| LLM-judge plumbing | `_build_testing_criteria:218` maps evaluator → `data_mapping` ( `{{item.query}}` vs `{{item.query_messages}}` + `tool_definitions`/`ground_truth`/`context`), `_build_item_schema:309`, `_evaluate_via_dataset:892`, `_poll_eval_run:383`, `_fetch_output_items:560` extracting per-dimension `RubricScore` via `_extract_rubric_scores:512` | `python/packages/foundry/agent_framework_foundry/_foundry_evals.py:218` |
| Generated rubric evaluators | `GeneratedEvaluatorRef:61` with `version` pinning, warning on unpinned `latest:92`, mixed with built-ins in `_build_testing_criteria:239` | `python/packages/foundry/agent_framework_foundry/_foundry_evals.py:61` |
| Foundry trace/target evals | `evaluate_traces:975` (`azure_ai_traces`/`azure_ai_responses`), `evaluate_foundry_target:1067` (`azure_ai_target_completions`) | `python/packages/foundry/agent_framework_foundry/_foundry_evals.py:975` |
| Local eval tests (coverage) | 1226-line test file covering `keyword_check`, `tool_call_args_match`, `LocalEvaluator` aggregation, `@evaluator` return coercion, `num_repetitions` modulo, `assert_*` helpers | `python/packages/core/tests/core/test_local_eval.py:1` |
| Foundry eval tests | Converter, tool extraction, `_build_testing_criteria`, dataset vs response paths, ground-truth & image preservation, polling, `FoundryEvals` constructor | `python/packages/foundry/tests/test_foundry_evals.py:1` |
| Sample evaluation suites | 11 evaluation samples: `evaluate_agent.py:1`, `evaluate_with_expected.py`, `evaluate_workflow.py`, `foundry_evals/evaluate_agent_sample.py`, `evaluate_tool_calls_sample.py`, `evaluate_multiturn_sample.py`, `evaluate_mixed_sample.py`, `evaluate_with_rubric_sample.py`, `evaluate_traces_sample.py` | `python/samples/02-agents/evaluation/evaluate_agent.py:1`, `python/samples/05-end-to-end/evaluation/foundry_evals/evaluate_agent_sample.py:1` |
| Decision record | `0023-foundry-evals-integration.md` describes Provider-agnostic Evaluator protocol, LocalEvaluator, mix-and-match, per-turn/full splits, and notes `Expected: datasets with expected outputs` as future work | `docs/decisions/0023-foundry-evals-integration.md:47` |
| CI workflow | `python-tests.yml:46` runs `uv run poe test -A --junitxml=pytest.xml` (unit tests only); no dedicated prompt-regression or nightly eval gate | `.github/workflows/python-tests.yml:46` |
| Sample validation (not eval) | `python-sample-validation` workflow validates samples run without crash; not a prompt-quality gate | `.github/workflows/python-sample-validation.yml:1` |
| Transparency FAQ | States framework undergoes "engineering testing ... conformance testing" but "AI performance metrics such as accuracy ... are dependent on underlying LLM providers" — evaluation is application-specific | `TRANSPARENCY_FAQ.md:26` |
| No experiment tracker evidence | No files matching `mlflow`, `wandb`, `comet`, `experiment`, `ab_test` in core/foundry packages; grep across `python/packages` finds only doc mentions | `python/packages/core/agent_framework/_evaluation.py:1` (searched, no evidence) |
| No prompt registry/template | No `PromptTemplate`, `PromptRegistry`, `prompt_version`, or prompt-store abstraction; `Agent.instructions` is free-form string, no versioning | `python/packages/core/agent_framework/_agents.py:1` (searched, no evidence) |

## Answers to Dimension Questions

### 1. Are prompt changes tested?
**Partially — framework level, not prompt level.**

Agent instruction prompts (`Agent.instructions` in `python/packages/core/agent_framework/_agents.py:1`) are plain strings with no registry, versioning, or diff testing. No test asserts that a prompt change preserves output contract. What *is* tested is the evaluation machinery itself: `test_local_eval.py:1` covers 1226 lines of `keyword_check`, `tool_calls_present`, `LocalEvaluator` aggregation, `@evaluator` coercion, and `expected_output` modulo edge cases; `test_foundry_evals.py:1` covers converter fidelity, criteria mapping, and dataset construction including ground-truth and image content. At runtime `evaluate_agent:1629` does stamp `expected_output:1813` and `expected_tool_calls:1821` onto `EvalItem`s for similarity/LLM-judge scoring, and samples like `python/samples/02-agents/evaluation/evaluate_with_expected.py` demonstrate ground-truth comparison. However invoking evaluation against prompt changes is manual — there is no hook that runs `evaluate_agent` when `instructions` changes.

Verdict: The harness *enables* testing prompt changes (via `LocalEvaluator` + custom `@evaluator`), but does not *require* or *automate* it.

### 2. Are experiments tracked?
**No — absent.**

No MLflow, Weights & Biases, Comet, or internal experiment store found in `python/packages/core` or `python/packages/foundry`. The only tracking is `FoundryEvals`'s ephemeral Foundry portal link (`report_url` surfaced on `EvalResults:388` and extracted in `_poll_eval_run:415`) and local `EvalResults.per_evaluator:440` counters. `docs/decisions/0023-foundry-evals-integration.md:40` notes "Foundry-native results ... viewable in the Foundry portal with dashboards and comparison views," but those dashboards live in Azure, not in-repo, and there is no API in this repo for A/B naming, run history, metric time-series, or parameter logging. `_build_testing_criteria:218` logs a warning for unpinned rubric versions, but version pinning is advisory (`GeneratedEvaluatorRef:88`). No dataset versioning for prompts or eval inputs.

### 3. Is LLM-as-judge used for evaluation?
**Yes — mature, provider-agnostic.**

`FoundryEvals` (`python/packages/foundry/agent_framework_foundry/_foundry_evals.py:728`) is an LLM-as-judge provider that delegates to Foundry's hosted evaluators (`builtin.relevance`, `coherence`, `task_adherence`, `tool_call_accuracy`, `groundedness`, `similarity`, safety etc. at `python/packages/foundry/agent_framework_foundry/_foundry_evals.py:142`). The plumbing builds per-evaluator `testing_criteria:238` with `deployment_name` as the judge model and `data_mapping` that distinguishes string-query (`quality`) vs array-query (`agent/tool`) judges. Polling (`_poll_eval_run:383`, timeout 180s via `python/packages/foundry/agent_framework_foundry/_foundry_evals.py:843`) fetches `EvalScoreResult` per item and rubric per-dimension breakdowns (`RubricScore:650`, extracted via `_extract_rubric_scores:512` handling `properties.dimension_scores`/`rubric_scores`). Local LLM-as-judge is also first-class: `@evaluator:1415` supports `async def llm_judge(query: str, response: str) -> float` where the body can call any LLM client (example at `python/packages/core/agent_framework/_evaluation.py:1459`), and `LocalEvaluator:1518` awaits it automatically (`_run_check:1510`). `docs/decisions/0023-foundry-evals-integration.md:30` explicitly calls out "Foundry LLM-as-judge evaluators alongside fast local checks."

### 4. Are regressions caught before deployment?
**No automated gate; opt-in CI assertions exist but are not enforced.**

The repo ships four CI-ready assertions (`raise_for_status:470`, `assert_score_at_least:502`, `assert_dimension_score_at_least:544`, `assert_no_failed_items:616`) plus `all_passed:455` recursion over workflow `sub_results`, designed to be called in tests/CI (see sample comment `results[0].raise_for_status()` in `python/samples/02-agents/evaluation/evaluate_agent.py:78`). Test coverage for these helpers exists (`test_local_eval.py:1021`). However `.github/workflows/python-tests.yml:46` only runs `pytest -m "not integration"`; Foundry evals and local-eval integration suites are not gated per-PR with threshold checks, and there is no workflow that runs `evaluate_agent` against a prompt fixture and fails PRs on score drop. The sample-validation workflow validates that samples *execute* without exception, not that prompt quality holds. `TRANSPARENCY_FAQ.md:28` confirms evaluation is left to application developers. Thus regression catching is possible and documented but not enforced — a prompt can ship without evaluation.

## Architectural Decisions

| Decision | Rationale | Consequence | File:Line |
|----------|-----------|-------------|-----------|
| Provider-agnostic `Evaluator` protocol (`name` + `evaluate(items, eval_name) -> EvalResults`) | Decouple evaluation from any provider (Foundry vs local vs third-party like DeepEval/RAGAS) per `0023-foundry-evals-integration.md:52` | Mix-and-match evaluators (`local + foundry`) with one `evaluate_agent` call; low concept count | `python/packages/core/agent_framework/_evaluation.py:681` |
| `EvalItem.conversation` as single source of truth + `ConversationSplitter` | Support last-turn vs full-trajectory vs custom splits; different factorings measure different concerns | `split_messages:235` + `ConversationSplit.LAST_TURN/FULL:130`; workflow `per_turn_items:254` generates N eval items per conversation | `python/packages/core/agent_framework/_evaluation.py:182` |
| `LocalEvaluator` with parameter-name injection (`@evaluator`) | "Bring your own evaluator as simple as writing a function" — reduce boilerplate vs `Evaluator` subclass | Supports `query/response/expected_output/conversation/tools/context` params + `bool/float/dict/CheckResult` returns; unknown required param raises at decoration time | `python/packages/core/agent_framework/_evaluation.py:1415`, `python/packages/core/agent_framework/_evaluation.py:1296` |
| `FoundryEvals` transient JSONL dataset per run | Foundry Evals API requires `evals.create` + `runs.create`; framework builds schema dynamically (`has_context/has_tools/has_ground_truth`) and maps `EvalItem` → `{query, response, query_messages, response_messages, tool_definitions, ground_truth, context}` | No persisted dataset; reproducibility depends on caller supplying same `queries`/`expected_output`; handles auto-detection of defaults and filtering of tool evaluators | `python/packages/foundry/agent_framework_foundry/_foundry_evals.py:892` |
| `EvalResults` with `sub_results` for workflow evaluation | Pinpoint underperformance per sub-agent; mirrors workflow `WorkflowRunResult` extraction via `executor_invoked/completed` pairing | `evaluate_workflow:1832` aggregates `total_passed/failed` into synthetic overall result when `include_overall=False` | `python/packages/core/agent_framework/_evaluation.py:373` |
| Experimental feature gate (`@experimental(feature_id=EVALS)`) | Evaluation API still stabilizing; allows breaking changes behind warning | All public eval symbols raise experimental warning at import/use; documented in `PACKAGE_STATUS.md:70` | `python/packages/core/agent_framework/_evaluation.py:68` |
| `GeneratedEvaluatorRef` with explicit `version` pinning | Ensure reproducible rubric runs; CI should pin version | Warning when `version=None:249`; CI docs advise concrete version | `python/packages/foundry/agent_framework_foundry/_foundry_evals.py:61` |

## Notable Patterns

- **Progressive disclosure checks**: Built-ins (`keyword_check`, `tool_called_check`) cover 80% cases; `@evaluator` covers custom logic without learning `Evaluator` protocol; full `Evaluator` implementation reserved for provider integrations. (`python/packages/core/agent_framework/_evaluation.py:1062`)
- **Async-first evaluator**: `_run_check:1510` and `@evaluator._check:1488` `await` if awaitable, allowing `async def llm_judge` to call LLM APIs inline without separate wrapper.
- **Bare-callable auto-wrapping**: `_resolve_evaluators:2056` collapses adjacent `EvalCheck` callables into a single `LocalEvaluator`, so `evaluators=[is_helpful, keyword_check("x")]` works without manual wrapping (`python/packages/core/agent_framework/_evaluation.py:2056`, tests at `python/packages/core/tests/core/test_local_eval.py:459`).
- **`num_repetitions` for consistency**: `evaluate_agent:1641` runs each query N times independently; expected outputs are modulo-stamped to measure non-determinism (`python/packages/core/agent_framework/_evaluation.py:1790`).
- **Rubric score flexibility**: `_extract_rubric_scores:512` accepts both `dimension_scores` and `rubric_scores` keys, dict or object shapes, for forward/backward compatibility across Foundry SDK versions.
- **Polling with timeout semantics**: `_poll_eval_run:383` returns `status="timeout"` rather than raising, letting callers distinguish infra timeout from eval failure (`EvalResults.status:381`).

## Tradeoffs

- **Transient vs versioned datasets**: Simplicity (no storage, no infra) vs loss of reproducibility and audit trail. Every `FoundryEvals.evaluate:864` re-creates an eval definition in Foundry and uploads a fresh JSONL payload; there is no in-repo dataset artifact to diff or pin. Reproducibility relies on external Foundry portal history.
- **Local vs cloud fidelity**: `LocalEvaluator` is fast, deterministic, and CI-friendly but shallow (keyword/tool presence only); deep quality/safety requires Foundry LLM judges which need credentials, network, quota, and introduce cost/latency variance. The two are combinable but live in different packages (`core` vs `foundry`).
- **String query/response vs message arrays**: `_build_testing_criteria:283` maps quality evaluators to string placeholders and agent evaluators to `query_messages/response_messages`. This doubles the surface to test (see `test_foundry_evals.py:736`) and risks mis-mapping if new evaluators arrive that don't fit the hardcoded `_AGENT_EVALUATORS` set.
- **Experimental gate**: All eval APIs gated behind `@experimental` warnings, which may deter adoption despite maturity; PACKAGE_STATUS marks evals as exported but still experimental.
- **Foundry lock-in vs provider agnosticism**: The protocol is provider-agnostic, but the only hosted LLM-judge concretion shipped in-tree is `FoundryEvals`; other providers (OpenAI Evals, DeepEval wrappers) would require user-authored adapters.
- **No prompt registry**: Keeping `Agent.instructions` as free-form string minimizes API surface, but forfeits prompt versioning, A/B tracking, and diffing that teams expect for prompt experiments.

## Failure Modes / Edge Cases

- **No client / no project_client + no env**: `FoundryEvals.__init__:846` auto-creates `FoundryChatClient` from env; if `FOUNDRY_PROJECT_ENDPOINT`/`FOUNDRY_MODEL` missing it raises `ValueError` deep inside the client rather than a typed `EvalConfigurationError` — user sees generic failure.
- **Sync `AIProjectClient` passed**: `_resolve_openai_client:672` detects sync client and raises `TypeError` with guidance to use `azure.ai.projects.aio` — helpful, but only discovered at first `evaluate` call.
- **Missing `expected_output` length mismatch**: `evaluate_agent:1757` raises `ValueError` if `len(expected_output) != len(queries)` — guards silent misalignment but not ordering bugs.
- **Tool evaluators with no tools**: `_filter_tool_evaluators:352` removes `tool_call_accuracy` etc. when items lack `tool_definitions`; if all evaluators are tool-only it raises `ValueError` suggesting "add tool definitions or choose evaluators" — could surprise users whose agent sometimes has tools, sometimes not.
- **Polling timeout**: `_poll_eval_run:422` returns `EvalResults(status="timeout")` with `result_counts=None`. Callers checking only `all_passed` (which returns `False` for non-completed, `python/packages/core/agent_framework/_evaluation.py:462`) will catch it, but callers summing `passed/failed` without status check may silently treat timeout as zero results.
- **Unparseable tool-call arguments**: `AgentEvalConverter.convert_message:788` sanitizes JSON parse failures to `{"_raw_arguments": "[unparseable]"}` to avoid leaking sensitive args — score may be degraded but not errored.
- **All tool evaluators filtered with no fallback**: Example — `FoundryEvals(evaluators=[TOOL_CALL_ACCURACY])` with a single-turn agent that has no tools errors out rather than falling back to quality defaults; `_resolve_default_evaluators:336` auto-adds tool evaluators only when items have tools but does not auto-substitute quality ones.
- **Version drift for rubric evaluators**: Unpinned `GeneratedEvaluatorRef.latest:92` resolves to portal's current version at runtime; repeated runs may diverge without warning unless CI pins `version`.
- **Multi-modal gap**: `query`/`response` string projections are text-only; multi-modal content (images) is preserved in `query_messages/response_messages` for Foundry but not exposed via `@evaluator(response: str)` — evaluators using string overloads silently miss image context (noted as Known Limitation in `docs/decisions/0023-foundry-evals-integration.md:557`).

## Future Considerations

- **Prompt registry / templating**: Introduce a `PromptTemplate` with versioning and diffing; wire `evaluate_agent` to accept `prompt_version` and log it into `EvalResults.metadata` for comparison. Without it, prompt experiments remain ad-hoc.
- **Dataset abstraction**: `docs/decisions/0023-foundry-evals-integration.md:553` already notes "Datasets with expected outputs: A dataset abstraction ... is a natural next step but not yet designed." A `Dataset` type (jsonl/csv + version + split) would complement `expected_output` stamping and enable shared ground-truth corpora (e.g., GAIA via `load_dataset` pattern).
- **Experiment tracker integration**: Add optional `ExperimentTracker` protocol (log params/metrics/artifacts) and implement lightweight adapters for MLflow/W&B; currently Foundry portal is the only run history.
- **CI gating workflow**: Provide a reusable GitHub Action (`evaluate-agent`) that runs `LocalEvaluator` + optional `FoundryEvals` and enforces `assert_score_at_least` thresholds; today developers must hand-roll this.
- **Stabilize experimental gate**: Remove `@experimental` from core eval types and promote to GA per `PACKAGE_STATUS.md:70`; document SemVer guarantees for `Evaluator` protocol.
- **Dynamic evaluator discovery**: Replace hardcoded `_BUILTIN_EVALUATORS:142` allowlist with server-side discovery or looser validation to avoid blocking new Foundry evaluators.

## Questions / Gaps

- **No evidence for prompt version diff or A/B report**: Searched `python/packages/core` and `python/packages/foundry` and `docs` — no prompt registry, no comparative report generator, no delta scoring between two prompt versions. State explicitly: None found; search covered `prompt`, `template`, `registry`, `version`, `ab_test`, `experiment`, `tracking`.
- **Evaluation dataset persistence**: How are `EvalItem` lists intended to be versioned and shared across team members? `evaluate_agent(queries=[...])` inlines queries in sample code; there is no `Dataset` loader or fixtures directory for eval sets beyond ad-hoc lists and the GAIA example (`docs/decisions/0023-foundry-evals-integration.md:520`).
- **Cost/latency observability for Foundry judges**: `EvalItemResult.token_usage:353` captures per-item token counts after fetch, but `_poll_eval_run` does not surface aggregate cost/latency; no guidance on budgeting LLM-judge spend across large eval suites.
- **Auth for CI Foundry evals**: Foundry evals require interactive or service-principal creds; no example GitHub Actions OIDC setup for unattended eval runs.
- **.NET parity**: Analysis focused on Python source per task isolation; .NET evaluation path (`Microsoft.Agents.AI` + `Microsoft.Agents.AI.AzureAI`) claims parity via ADR but not inspected in this study.

---

Generated by `Dimension 12.03: Prompt Evaluation and Experiments` against `agent-framework`.
