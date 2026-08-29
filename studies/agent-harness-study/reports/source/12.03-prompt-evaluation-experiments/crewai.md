# Source Analysis: crewai

## Dimension 12.03: Prompt Evaluation and Experiments

### Source Info

| Field | Value |
|-------|-------|
| Name | crewai |
| Path | `studies/agent-harness-study/sources/crewai` |
| Language / Stack | Python / Pydantic, LiteLLM, Rich console |
| Analyzed | 2026-08-29 |

## Summary

CrewAI provides an **experimental** LLM-as-judge evaluation stack alongside a separate legacy `Crew.test()` / `CrewEvaluator` path. The experimental stack (`lib/crewai/src/crewai/experimental/evaluation/:1`) defines `AgentEvaluator`, `ExperimentRunner`, `ExperimentResults`, six judge-based metric evaluators (goal alignment, semantic quality, reasoning efficiency, tool selection, parameter extraction, tool invocation), and assertion helpers for regression detection. `ExperimentRunner` accepts user-supplied datasets (`list[dict]` with `inputs`/`expected_score`/`identifier`). `ExperimentResults.compare_with_baseline()` persists runs to JSON and classifies regressions. There is no shipped eval dataset, no prompt-template snapshot testing in CI, no prompt-aware A/B framework, and no CI gate that executes real LLM judges. All metric evaluators are `experimental` import paths (`lib/crewai/pyproject.toml` / docs never promote them to stable), and `conftest.py:216` sets `CREWAI_TESTING=true` but does not provision judge keys — tests mock `llm.call`. The legacy CLI `crewai test` (`lib/cli/src/crewai_cli/evaluate_crew.py:9`, `lib/crewai/src/crewai/crew.py:2227`) runs n-iterations and prints a table via an evaluator agent, which is operational convenience, not a regression harness. A prompt change (e.g., editing `lib/crewai/src/crewai/translations/en.json:8` slices consumed by `lib/crewai/src/crewai/utilities/prompts.py:99` and `lib/crewai/src/crewai/utilities/i18n.py:56`) can be deployed without any automated check that output quality or tool behavior is preserved.

## Rating

**4 / 10 — Present but inconsistent, weakly documented, fragile and experimental.**

Rationale: The framework explicitly implements LLM-as-judge scoring, dataset-driven `ExperimentRunner`, and baseline-comparison regression logic with tests, but it is marked `experimental` (`lib/crewai/src/crewai/experimental/evaluation/__init__.py:1-48`), unintegrated into the main prompt lifecycle (`I18N`/`Prompts` have no versioning or golden-output fixtures), ships zero datasets, and has no CI job that runs real judge evaluations or compares prompt versions. `tests.yml:34-97` runs only mocked unit tests. Confidence for deploying a prompt change without regression is low.

## Evidence Collected

Every entry MUST include a file path with line numbers. Format: `path/to/file.ts:NN`.

| Area | Evidence | File:Line |
|------|----------|-----------|
| Eval datasets - user-supplied dataset interface | `ExperimentRunner.__init__(self, dataset: list[dict[str, Any]])` stores user dataset; `_run_test_case` extracts `inputs`, `expected_score`, `identifier` (md5 fallback) | `lib/crewai/src/crewai/experimental/evaluation/experiment/runner.py:27-70` |
| Eval datasets - no shipped datasets | Glob of source finds no `eval_datasets/` or JSON fixtures; only example `inputs={"query": "Test query 1"}` in tests | `lib/crewai/tests/experimental/evaluation/test_experiment_runner.py:56-71` |
| Prompt change testing (framework-level) | Legacy `Crew.test(n_iterations, eval_llm, inputs)` copies crew, iterates kickoffs, uses `CrewEvaluator` to score tasks 1-10 via an evaluator Agent | `lib/crewai/src/crewai/crew.py:2227-2272` |
| Prompt change testing (legacy CrewEvaluator) | `CrewEvaluator.evaluate()` builds evaluator Agent `role="Task Execution Evaluator"` and Task with `description` referencing `task_description`, `task_expected_output`, `Task Output`, returns `quality:float 1-10` | `lib/crewai/src/crewai/utilities/evaluators/crew_evaluator_handler.py:58-199` |
| Prompt change testing (legacy TaskEvaluator/train) | `TaskEvaluator.evaluate(task, output)` builds `Converter(llm, TaskEvaluation)` query with quality 0-10; `evaluate_training_data` aggregates human feedback | `lib/crewai/src/crewai/utilities/evaluators/task_evaluator.py:69-187` |
| CLI test invocation | `evaluate_crew(n_iterations, model)` shells `["uv", "run", "test", str(n_iterations), model]` | `lib/cli/src/crewai_cli/evaluate_crew.py:9-20` |
| Prompt store | Single JSON `translations/en.json` with `slices.role_playing`, `slices.tools`, `slices.task`, etc., loaded via `I18N.load_prompts()` at `utilities/i18n.py:27-54` | `lib/crewai/src/crewai/translations/en.json:8-22`, `lib/crewai/src/crewai/utilities/i18n.py:36-44` |
| Prompt construction | `Prompts.task_execution()` composes slices via `I18N_DEFAULT.slice(component)` and applies `prompt_template` / `system_template` interpolation (`{{ .System }}`, `{{ .Prompt }}`, `{role}`) | `lib/crewai/src/crewai/utilities/prompts.py:93-257` |
| Prompt construction tests (no golden output) | Prompt template tests check `system`/`prompt` presence and `{role}` substitution, not semantic equivalence or snapshots | `lib/crewai/tests/agents/test_agent.py:156-181`, `lib/crewai/src/crewai/utilities/prompts.py:93-141` |
| Experiment tracker - runner loop | `ExperimentRunner.run(crew, agents, print_summary)` resets evaluator, runs each test_case, extracts averaged metric scores, asserts pass via `_assert_scores` | `lib/crewai/src/crewai/experimental/evaluation/experiment/runner.py:32-165` |
| Experiment tracker - baseline persistence | `ExperimentResults.to_json(filepath)` writes timestamp/metadata/results; `compare_with_baseline(baseline_filepath)` loads JSON, sorts by timestamp, classifies `improved`/`regressed`/`unchanged`/`new_tests`/`missing_tests` | `lib/crewai/src/crewai/experimental/evaluation/experiment/result.py:32-143` |
| Experiment tracker - assertion helpers | `assert_experiment_successfully()` fails on `not result.passed`, then calls `compare_with_baseline` + `assert_experiment_no_regression`; `assert_experiment_no_regression` raises on `regressed` | `lib/crewai/src/crewai/experimental/evaluation/testing.py:13-54` |
| LLM-as-judge prompts - GoalAlignment | System prompt scores 0-10 alignment, requires JSON `{score, feedback}`; calls `self.llm.call(prompt)` then `extract_json_from_llm_response` | `lib/crewai/src/crewai/experimental/evaluation/metrics/goal_metrics.py:36-88` |
| LLM-as-judge prompts - SemanticQuality | System prompt scores semantic quality 0-10 with criteria structure/reasoning/clarity; same `llm.call` + JSON extraction pattern | `lib/crewai/src/crewai/experimental/evaluation/metrics/semantic_quality_metrics.py:35-89` |
| LLM-as-judge prompts - ToolSelection | Scores relevance/coverage 0-10, overall_score average, explicit "DO NOT suggest tools not in Available tools" guard | `lib/crewai/src/crewai/experimental/evaluation/metrics/tools_metrics.py:66-146` |
| LLM-as-judge prompts - ParameterExtraction | Scores accuracy/formatting/completeness; enumerates validation_error rate and param samples | `lib/crewai/src/crewai/experimental/evaluation/metrics/tools_metrics.py:223-302` |
| LLM-as-judge prompts - ToolInvocation | Scores structure/error_handling/invocation_patterns; summarizes error_rate and error_types breakdown | `lib/crewai/src/crewai/experimental/evaluation/metrics/tools_metrics.py:374-458` |
| LLM-as-judge prompts - ReasoningEfficiency | Scores focus/progression/decision_quality/conciseness/loop_avoidance; includes statistical pre-analysis (Jaccard loop detector, trend, loop likelihood) before LLM call | `lib/crewai/src/crewai/experimental/evaluation/metrics/reasoning_metrics.py:114-225` |
| Evaluation listener / trace collection | `EvaluationTraceCallback` singleton attaches to `crewai_event_bus` for `AgentExecutionStartedEvent`, `ToolUsageFinishedEvent`/`ErrorEvent`, `LLMCallStarted/CompletedEvent`, stores `tool_uses`/`llm_calls` per `agent_id_task_id` | `lib/crewai/src/crewai/experimental/evaluation/evaluation_listener.py:30-355` |
| Metric categories enum | `MetricCategory` defines 6 categories: `GOAL_ALIGNMENT`, `SEMANTIC_QUALITY`, `REASONING_EFFICIENCY`, `TOOL_SELECTION`, `PARAMETER_EXTRACTION`, `TOOL_INVOCATION` | `lib/crewai/src/crewai/experimental/evaluation/base_evaluator.py:20-30` |
| Evaluation core model | `EvaluationScore` (score 0-10 ge/le, feedback, raw_response), `AgentEvaluationResult` (metrics dict), `AgentAggregatedEvaluationResult` + `AggregationStrategy` (simple_avg, weighted, best/worst) | `lib/crewai/src/crewai/experimental/evaluation/base_evaluator.py:32-119` |
| Regression test suite (unit) | `test_experiment_runner.py` mocks `create_default_evaluator` and checks pass/fail logic for single, dict, missing metric cases; `test_experiment_result.py` mocks baseline JSON and asserts regression classification | `lib/crewai/tests/experimental/evaluation/test_experiment_runner.py:52-209`, `lib/crewai/tests/experimental/evaluation/test_experiment_result.py:52-111` |
| Regression test suite (metric mocks) | `test_goal_metrics.py` mocks `LLM.call` to return `{"score":8.5, ...}` and asserts JSON→`EvaluationScore` + error fallback to `score=None` | `lib/crewai/tests/experimental/evaluation/metrics/test_goal_metrics.py:14-66` |
| CI prompt test integration - absent | `tests.yml` runs `uv run pytest` splits 8 across 4 python versions; no step runs `crewai test`, `ExperimentRunner`, or `compare_with_baseline`; no baseline artifact persisted | `.github/workflows/tests.yml:34-97` |
| Patronus / external audit not integrated | `PatronusEvalTool` is a separate `crewai-tools` tool (Lyceum) not wired into `experimental/evaluation` pipeline | `lib/crewai-tools/src/crewai_tools/tools/patronus_eval_tool/patronus_eval_tool.py:1` (seen via bash glob) |

## Answers to Dimension Questions

1. **Are prompt changes tested?**
   - Partially, but not reliably. Prompts are a single JSON (`lib/crewai/src/crewai/translations/en.json:8`) composed via `lib/crewai/src/crewai/utilities/prompts.py:93-257` and interpolated via `lib/crewai/src/crewai/utilities/i18n.py:56-65`. Unit tests (`lib/crewai/tests/agents/test_agent.py:156-181`) assert template plumbing, not semantic preservation. The experimental evaluation suite can score outputs via LLM judges (`lib/crewai/src/crewai/experimental/evaluation/metrics/goal_metrics.py:36-88` etc.), but it is opt-in, undocumented in `docs/edge/en/concepts/testing.mdx:1-49`, and never gates a prompt edit. A change to a `slice` (e.g., `slices.tools`) can ship with zero test feedback.

2. **Are experiments tracked?**
   - Minimally. `ExperimentRunner` (`lib/crewai/src/crewai/experimental/evaluation/experiment/runner.py:26-57`) runs a caller-supplied dataset; `ExperimentResults` (`lib/crewai/src/crewai/experimental/evaluation/experiment/result.py:32-97`) serializes runs to a JSON filepath with timestamp and can append to an accumulating baseline file (`result.py:88-95`). There is no experiment registry, no run ID beyond timestamp, no hyperparameter/prompt-version linkage, no dashboard, no hosted store (e.g., MLflow/Braintrust tracking is docs-only at `docs/edge/en/observability/braintrust.mdx`). The file-based approach is local and ephemeral; `conftest.py:219-244` uses a temp storage dir per test.

3. **Is LLM-as-judge used for evaluation?**
   - Yes, explicitly and pervasively — but purely experimental. `AgentEvaluator.evaluate()` (`lib/crewai/src/crewai/experimental/evaluation/agent_evaluator.py:248-294`) iterates `BaseEvaluator`s and dispatches `AgentEvaluationStarted/Completed/FailedEvent`. Six judges are shipped: `GoalAlignmentEvaluator` (`goal_metrics.py:39-53`), `SemanticQualityEvaluator` (`semantic_quality_metrics.py:38-53`), `ToolSelectionEvaluator`, `ParameterExtractionEvaluator`, `ToolInvocationEvaluator` (`tools_metrics.py:22-458`), `ReasoningEfficiencyEvaluator` (`reasoning_metrics.py:114-172`). Each builds a system prompt with a 0-10 rubric, calls `self.llm.call(prompt)` (`goal_metrics.py:71`), extracts JSON (`json_parser.py`), and returns `EvaluationScore(score, feedback, raw_response)`. The legacy path uses an evaluator Agent task (`crew_evaluator_handler.py:58-85`). All are mock-tested, not run live in CI.

4. **Are regressions caught before deployment?**
   - No reliable gate exists. The mechanism exists: `testing.py:13-45` `assert_experiment_successfully` + `compare_with_baseline` + `assert_experiment_no_regression` raises `AssertionError` on regressed identifiers, with warning on missing tests (`testing.py:47-54`). However: (a) it requires the developer to author a dataset and a baseline file (`testing.py:67-76` fallback `"{test_func}_results.json"`) and to wire the test into pytest; (b) the provided tests (`test_experiment_runner.py:52-209`, `test_experiment_result.py:52-111`) are mocked and do not exercise a real prompt change; (c) `.github/workflows/tests.yml:34-97` does not execute `ExperimentRunner` or persist/compare baselines. Consequently, a prompt regression would only be caught if a team built their own harness — the framework itself will not block a release.

## Architectural Decisions

- **Separate experimental vs. legacy evaluator lanes:** Stable `CrewEvaluator`/`TaskEvaluator` (`lib/crewai/src/crewai/utilities/evaluators/crew_evaluator_handler.py:29`, `task_evaluator.py:52`) coexist with `experimental.evaluation` (`lib/crewai/src/crewai/experimental/evaluation/__init__.py:1-48`). Decision keeps stable API small but fragments evaluation mental model and documentation (`docs/edge/en/concepts/testing.mdx:13` documents only `crewai test`, never the experimental judges).
- **Event-bus trace capture singleton:** `EvaluationTraceCallback` (`evaluation_listener.py:30-56`) is a singleton registered via `create_evaluation_callbacks()` (`evaluation_listener.py:345-355`) listening to `AgentExecutionStarted/Completed`, `ToolUsageFinished/Error`, `LLMCallStarted/Completed`. Decouples evaluation from agent loop but introduces global state (`_instance`, `traces: dict`) that tests must clear (`conftest.py:186-198` clears handlers) and risks cross-test pollution if not isolated.
- **LLM-as-judge via `BaseLLM.call` + JSON extraction:** Every metric calls `self.llm.call(prompt)` (`goal_metrics.py:71`, `semantic_quality_metrics.py:71`) and `extract_json_from_llm_response` (`base_evaluator.py` via `json_parser.py`). Avoids typed tool-calling, but forces brittle markdown-fence stripping and `score=None` fallback (`goal_metrics.py:84-88`, `tools_metrics.py:141-146`). No schema-validated function call.
- **File-local baseline comparison:** `ExperimentResults.compare_with_baseline(baseline_filepath)` (`result.py:48-97`) treats a single JSON file as the registry, auto-creating it on first run (`result.py:70-78`) and appending (`result.py:88-95`). Simple and dependency-free but not concurrent-safe, not queryable, and not shared across runners.
- **MetricCategory enum as scoring contract:** Six fixed categories (`base_evaluator.py:20-30`) with per-metric `EvaluationScore` aggregated via `AggregationStrategy` (`base_evaluator.py:80-99`). Provides uniform 0-10 scale but has no weighting guidance and no correlation to user-facing quality.
- **Dataset as `list[dict]` with loose schema:** `ExperimentRunner` expects `{identifier?, inputs: dict, expected_score: float|dict}` (`runner.py:65-69`), validates via `_assert_scores` with four comparison branches (`runner.py:130-165`). Flexible for callers but undocumented schema and no validation error reporting beyond boolean `passed`.

## Notable Patterns

- **Judge-prompt pattern:** All judges follow `system: rubric + JSON contract` → `user: agent role/goal + task_context + trace excerpt + final_output` → `llm.call` → `extract_json`. Example: `goal_metrics.py:36-68`, `reasoning_metrics.py:114-172`, `tools_metrics.py:66-111`.
- **Graceful-degradation scoring:** Each evaluator returns `score=None` on insufficient context (`reasoning_metrics.py:61-65` for `<2 LLM calls`, `tools_metrics.py:44-51` for no tools) or parse failure (`goal_metrics.py:84-88`), allowing aggregation to skip categories (`runner.py:114-128` filters `score is not None`).
- **N-gram loop detector as pre-LLM heuristic:** `ReasoningEfficiencyEvaluator._detect_loops` (`reasoning_metrics.py:227-258`) computes Jaccard similarity >0.7 as a loop signal, annotated "Simple n-gram similarity; embedding-based would be more robust" — a conscious cheap heuristic before the LLM judge.
- **Rich console display as experiment UX:** `ExperimentResultsDisplay` / `EvaluationDisplayFormatter` (`agent_evaluator.py:58-59`, `result.py:26-44`) print tables to terminal rather than returning machine-readable diffs; aligns with `crewai test` table UX (`docs/edge/en/concepts/testing.mdx:36-48`).
- **Mock-heavy unit tests:** Evaluation tests mock `LLM.call` (`test_goal_metrics.py:18-23`) and `create_default_evaluator` (`test_experiment_runner.py:52-54`), guaranteeing fast CI but providing no proof judges work with real models.

## Tradeoffs

- **Flexibility vs. safety:** `ExperimentRunner` accepts arbitrary `list[dict]` (`runner.py:27`) — easy to adopt, but without a typed `Dataset` class, validation, or shipped fixtures, callers can silently mis-specify `expected_score` keys and get `passed=False` due to missing-key logic (`runner.py:155-163` only checks intersecting keys).
- **LLM judge expressiveness vs. cost/determinism:** Using an LLM for every metric (6 calls per task per iteration, `agent_evaluator.py:264-275` loops evaluators) yields rich feedback but makes `Crew.test -n 3` with 4 tasks = 72 judge calls; flaky JSON parsing and nondeterminism mean repeated runs can flip `passed` (`runner.py:130-165` is `actual >= expected`).
- **Singleton event bus vs. isolation:** Singleton `EvaluationTraceCallback` (`evaluation_listener.py:30-46`) gives zero-config tracing but requires manual `reset_iterations_results()` (`runner.py:48`) and `conftest.py:186-198` handler clearing; parallel test execution risks interleaved traces.
- **File baseline vs. platform tracking:** JSON baseline (`result.py:48-97`) is simple and VCS-commitable but offers no history compaction, no diff visualization, and will `json.dump([current_data], f)` on first run — overwriting intent is implicit.
- **General-purpose judges vs. prompt-specific assertions:** Judges score broad qualities (alignment, semantics) not prompt-specific invariants (e.g., "slices.tools must contain `{tool_names}` placeholder"); template regressions would be caught only indirectly.

## Failure Modes / Edge Cases

- **Silent prompt drift:** Editing `translations/en.json:8` `slices.task` ("Begin! ... your job depends on it!") changes agent behavior but no test asserts prompt content; `lib/crewai/tests/agents/test_agent.py:156` only checks a custom `prompt_template` override path.
- **Identifier collision / missing identifier:** `_run_test_case` (`runner.py:67-70`) hashes `str(test_case)` when `identifier` absent; two cases with same `inputs`/`expected_score` can collide; downstream `result.py:113-117` treats missing `test_identifier` as `new_tests` but never fails, so deduplication bugs are silent.
- **LLM judge parse failure masks regression:** All evaluators catch `Exception` and return `score=None` (`goal_metrics.py:83-88`, `tools_metrics.py:141-146`, `reasoning_metrics.py:219-225`); `runner.py:120-121` drops `None` scores before averaging, so a consistently broken judge can yield `avg_scores={}` → `actual={}` → `_assert_scores` may return `False` for unknown_metric (`test_experiment_runner.py:177-209` expects `False`) but `True` for empty expected (`runner.py:156-157`), hiding the breakage.
- **Missing metric key yields false negative/positive:** `_assert_scores` (`runner.py:155-163`) requires `matching_keys = set(expected) & set(actual)` and `all(actual[k] >= expected[k] for k in matching_keys)`; if expected requests `{"tool_selection":9}` but actual only has `goal_alignment`, `matching_keys` empty → `False` (fail) even if judge not run; conversely, if expected has unknown key also in actual (from a different judge), partial intersection can pass though critical metric missing (`test_experiment_runner.py:111-139` demonstrates `unknown_metric` passes if at least one other metric does).
- **No concurrency or retry for judge calls:** `AgentEvaluator.evaluate` (`agent_evaluator.py:264-294`) iterates evaluators sequentially, `LLM.call` with no retry/timeout handling; transient provider outage fails entire `evaluate` → `emit_evaluation_failed_event` + console print, but `ExperimentRunner._run_test_case` (`runner.py:104-112`) swallows exception and returns `score=0.0, passed=False` — conflating infra failure with quality failure.
- **Baseline file concurrency/ corruption:** `compare_with_baseline` (`result.py:56-68`) reads then appends to same `baseline_filepath` without file lock; parallel `ExperimentRunner` instances can lose data or write invalid JSON; `JSONDecodeError` only prints warning (`result.py:66-68`) then treats as `is_baseline=True` and overwrites.
- **Prompt i18n cache invalidation:** `get_i18n` is `@lru_cache(maxsize=None)` (`i18n.py:131-144`); editing a custom `prompt_file` on disk between calls within same process returns stale `I18N`; no `cache_clear` exposed — requires process restart to pick up prompt change in long-running experiment harness.
- **Token cost and rate limit on large datasets:** No batching or `max_concurrency`; running `dataset` of 50 cases * 6 judges = 300 LLM calls sequentially will hit rate limits and long `kickoff` times before any gate fires.

## Future Considerations

- **Promote evaluation out of `experimental`:** Stabilize import path (`lib/crewai/src/crewai/evaluation/`), document in `docs/edge/en/concepts/testing.mdx:1-49` alongside `crewai test`, and add version pinning so prompt changes reference evaluation config hash.
- **Ship canonical eval dataset(s) and golden outputs:** Add `lib/crewai/eval_datasets/` with minimal tasks (tool-use, no-tool, knowledge, multimodal) and committed `baseline_results.json` so `compare_with_baseline` has a real anchor; add a pytest marker `@pytest.mark.eval` that runs against it.
- **Prompt snapshot testing:** Add `tests/test_prompts_snapshot.py` that asserts `I18N().slice("tools")`, `Prompts(agent).task_execution().prompt` against checked-in snapshots; fail CI on untracked diff, requiring explicit `UPDATE_SNAPSHOTS=1`.
- **CI gate for LLM judges:** Add optional `eval` job to ` .github/workflows/tests.yml:34-97` (e.g., `eval: runs-on: ubuntu-latest; needs: tests-matrix; if: contains(labels, 'eval')`) that runs `ExperimentRunner` with `OPENAI_API_KEY` secret, uploads `ExperimentResults.to_json("eval-report.json")`, and checks `compare_with_baseline` with `save_current=false` to block regressed lanes without mutating baseline.
- **Typed dataset + validation:** Replace `list[dict]` with `Pydantic Dataset` (`identifier: str`, `inputs: dict`, `expected_score: float | dict[MetricCategory, float]`, `tags: list[str]`) with `model_validator` to reject unknown metrics early.
- **Centralized experiment store:** Support pluggable `ExperimentStore` (file | MLflow | Braintrust) behind a protocol so `ExperimentResults` can write to `mlflow.set_experiment("CrewAI")` referenced in docs without local JSON fragility.
- **Judge reliability hardening:** Switch judges to function-calling / structured output (`Converter.to_pydantic(TaskEvaluation)`) as in `task_evaluator.py:106-113` instead of raw `llm.call` + fence stripping; add retry with backoff and emit `ToolSelectionEvaluator` score as `None` only after N attempts.

## Questions / Gaps

- **Shipped datasets?** No evidence found. Searched `lib/**`, `docs/edge/en/**`; only synthetic `{"query": "Test query 1"}` in `lib/crewai/tests/experimental/evaluation/test_experiment_runner.py:59` and cassettes (`lib/crewai/tests/cassettes/`). If datasets exist outside source directory, they are not part of this isolated analysis per hard rule.
- **Prompt versioning?** No evidence of prompt version tracking or `prompt_file` hashing in `lib/crewai/src/crewai/utilities/i18n.py:36-44` or `lib/crewai/src/crewai/crew.py:329-330`. `prompt_file` is a plain path string; no checksum stored in `ExperimentResults.metadata` (`result.py:18-23` defaults to `{}`).
- **External evaluator integrations in use?** Patronus (`lib/crewai-tools/src/crewai_tools/tools/patronus_eval_tool/patronus_eval_tool.py`) and Braintrust/Maxim observability docs do not integrate with `experimental.evaluation`; search `grep -r "Patronus\|patronus" lib/crewai/src/crewai` returns zero hits.
- **CI secrets for judge LLMs?** `.github/workflows/tests.yml:1-137` has no `env: OPENAI_API_KEY` or evaluation job; whether outbound LLM calls are intended to be disabled under `GITHUB_ACTIONS` (as in `conftest.py:420-421` `record_mode="none"`) is inferred, not documented.
- **How is `crewai test -n` quality gate consumed?** `Crew.test()` (`crew.py:2259`) only prints via `print_crew_evaluation_result()` (`crew_evaluator_handler.py:95-175`); no exit code, no artifact, no threshold check — unclear how a team would enforce `quality >= X` in CI beyond reading console table.

---

Generated by `Dimension 12.03: Prompt Evaluation and Experiments` against `crewai`.
