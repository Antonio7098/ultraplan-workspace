# Source Analysis: pydantic-ai

## Dimension 12.03: Prompt Evaluation and Experiments

### Source Info

| Field | Value |
|-------|-------|
| Name | pydantic-ai |
| Path | `studies/agent-harness-study/sources/pydantic-ai` |
| Language / Stack | Python / uv workspace (pydantic-ai-slim, pydantic-evals, pydantic-graph, clai) |
| Analyzed | 2026-08-29 |

## Summary

`pydantic-ai` does not ship its own prompt-regression corpus; it ships a **framework for users to build one**: `pydantic-evals` (`pydantic_evals/pydantic_evals/`). The framework is mature and code-first: typed `Dataset`/`Case` with YAML/JSON persistence, pluggable deterministic + LLM-as-judge evaluators, `EvaluationReport` with diff/aggregate rendering, `generate_dataset` LLM helper, and online evaluation via `@evaluate` decorator emitting OTel `gen_ai.evaluation.result` events. The library's own prompts that matter — the four LLM-judge agents and `GEval` — are explicitly version-tested: system prompts are pinned to include a concise-reason instruction and to preserve `Input → Output → ExpectedOutput → Rubric` section ordering, with `inline-snapshot` golden tests. Experiment comparison is first-class (`EvaluationRenderer.build_diff_table`, `ReportAnalysis` types, `repeat`/`max_concurrency` controls, logfire traces) but there is no managed experiment store, no prompt-version registry, and no CI gate that runs user-defined golden prompts against live models; regression safety depends on the consumer wiring `Dataset.evaluate` + snapshot tests into their own CI.

## Rating

**7 / 10** — Clear model with tests, explicit interfaces, and operational safeguards for the evaluation substrate, but prompt-version/experiment tracking is framework-provided rather than opinionated platform-level, and application-level prompt regressions require user-owned CI integration.

Rationale: `pydantic_evals` delivers typed datasets (`pydantic_evals/pydantic_evals/dataset.py:177`), serialization + generation (`pydantic_evals/pydantic_evals/generation.py:33`), deterministic and LLM judges (`pydantic_evals/pydantic_evals/evaluators/common.py:224`, `pydantic_evals/pydantic_evals/evaluators/llm_as_a_judge.py:68`), span-aware agentic evaluators (`pydantic_evals/pydantic_evals/evaluators/agentic.py:190`), and report-level analyses (`pydantic_evals/pydantic_evals/reporting/analyses.py:62`). Tests pin judge prompt invariants (`tests/evals/test_llm_as_a_judge.py:60`, `tests/evals/test_llm_as_a_judge.py:118`) and all eval suites run in CI (` .github/workflows/ci.yml:337`), giving confidence that framework prompt changes are caught. Missing: a top-level prompt registry / A-B experiment tracker and automatic golden-prompt CI for downstream agents — both are left to the user to compose.

## Evidence Collected

Every entry MUST include a file path with line numbers. Format: `path/to/file.ts:NN`.

| Area | Evidence | File:Line |
|------|----------|-----------|
| Eval datasets | `Case` dataclass with `name/inputs/metadata/expected_output/evaluators` | `pydantic_evals/pydantic_evals/dataset.py:111` |
| Eval datasets | `Dataset[InputsT,OutputT,MetadataT]` generic collection with `cases/evaluators/report_evaluators` | `pydantic_evals/pydantic_evals/dataset.py:177` |
| Eval datasets | `DEFAULT_DATASET_PATH` / `DEFAULT_SCHEMA_PATH_TEMPLATE` + YAML/JSON load/save with JSON-schema generation | `pydantic_evals/pydantic_evals/dataset.py:77`, `pydantic_evals/pydantic_evals/dataset.py:747`, `pydantic_evals/pydantic_evals/dataset.py:797` |
| Eval datasets | `Dataset.evaluate` (async) with `max_concurrency`, `retry_task`, `repeat`, `lifecycle`, Logfire spans | `pydantic_evals/pydantic_evals/dataset.py:281` |
| Eval datasets | `Dataset.evaluate_sync` wrapping `run_until_complete`; repeat/multi-run via `_build_tasks_to_run` | `pydantic_evals/pydantic_evals/dataset.py:417`, `pydantic_evals/pydantic_evals/dataset.py:269` |
| Eval datasets | `from_file`/`from_text`/`from_dict` with `custom_evaluator_types` registry and `ExceptionGroup` on load errors | `pydantic_evals/pydantic_evals/dataset.py:557`, `pydantic_evals/pydantic_evals/dataset.py:598`, `pydantic_evals/pydantic_evals/dataset.py:637` |
| Eval datasets | LLM-assisted dataset generation via `generate_dataset(dataset_type, n_examples, extra_instructions)` | `pydantic_evals/pydantic_evals/generation.py:33` |
| Experiment trackers | `EvaluationReport` + `ReportCase`/`ReportCaseFailure`/`ReportCaseAggregate` + OTEL trace/span ids | `pydantic_evals/pydantic_evals/dataset.py:396`, `pydantic_evals/pydantic_evals/reporting/__init__.py` |
| Experiment trackers | `EvaluationRenderer.build_diff_table` and `default_render_number_diff` for baseline comparison | `pydantic_evals/pydantic_evals/reporting/render_numbers.py:62`, `tests/evals/test_reporting.py:187` |
| Experiment trackers | `ReportAnalysis` discriminated union: `ConfusionMatrix`, `PrecisionRecall`, `ScalarResult`, `LinePlot` | `pydantic_evals/pydantic_evals/reporting/analyses.py:21`, `pydantic_evals/pydantic_evals/reporting/analyses.py:62` |
| Experiment trackers | `ReportEvaluator` / `report_evaluators` running on full report + `analyses` attachment | `pydantic_evals/pydantic_evals/dataset.py:406`, `pydantic_evals/pydantic_evals/dataset.py:1013` |
| Experiment trackers | Online evaluation `@evaluate` / `OnlineEvalConfig` with `default_sample_rate`, `sampling_mode`, `max_concurrency`, `sink` | `pydantic_evals/pydantic_evals/online.py:194`, `pydantic_evals/pydantic_evals/online.py:415`, `pydantic_evals/pydantic_evals/online.py:496` |
| Experiment trackers | `span_tree` capture + `increment_eval_metric`/`set_eval_attribute` + `CaseLifecycle` hooks | `pydantic_evals/pydantic_evals/dataset.py:1140`, `pydantic_evals/pydantic_evals/lifecycle.py` |
| LLM-as-judge prompts | `GradingOutput` (`reason/pass/score`) + `_default_model = 'openai:gpt-5.2'` | `pydantic_evals/pydantic_evals/evaluators/llm_as_a_judge.py:29`, `pydantic_evals/pydantic_evals/evaluators/llm_as_a_judge.py:26` |
| LLM-as-judge prompts | Four judge agents: `_judge_output_agent`, `_judge_input_output_agent`, `_judge_input_output_expected_agent`, `_judge_output_expected_agent` | `pydantic_evals/pydantic_evals/evaluators/llm_as_a_judge.py:46`, `pydantic_evals/pydantic_evals/evaluators/llm_as_a_judge.py:85`, `pydantic_evals/pydantic_evals/evaluators/llm_as_a_judge.py:128`, `pydantic_evals/pydantic_evals/evaluators/llm_as_a_judge.py:176` |
| LLM-as-judge prompts | `judge_output`/`judge_input_output`/`judge_input_output_expected`/`judge_output_expected` helpers + `set_default_judge_model` | `pydantic_evals/pydantic_evals/evaluators/llm_as_a_judge.py:68`, `pydantic_evals/pydantic_evals/evaluators/llm_as_a_judge.py:109`, `pydantic_evals/pydantic_evals/evaluators/llm_as_a_judge.py:154`, `pydantic_evals/pydantic_evals/evaluators/llm_as_a_judge.py:200`, `pydantic_evals/pydantic_evals/evaluators/llm_as_a_judge.py:220` |
| LLM-as-judge prompts | `GEvalOutput` + `_judge_g_eval_agent` + `judge_g_eval` (numbered steps, `score_range` validation, out-of-range guard) | `pydantic_evals/pydantic_evals/evaluators/llm_as_a_judge.py:291`, `pydantic_evals/pydantic_evals/evaluators/llm_as_a_judge.py:302`, `pydantic_evals/pydantic_evals/evaluators/llm_as_a_judge.py:321` |
| LLM-as-judge prompts | `LLMJudge` evaluator (`rubric`, `include_input/expected_output`, `score/assertion` configs) wrapping the four judges | `pydantic_evals/pydantic_evals/evaluators/common.py:224` |
| LLM-as-judge prompts | `GEval` evaluator (`criteria`, `evaluation_steps`, `score_range`, `include_input`) wrapping `judge_g_eval` | `pydantic_evals/pydantic_evals/evaluators/common.py:288` |
| LLM-as-judge prompts | `_build_prompt` emitting `Input→Output→ExpectedOutput→Rubric` with `UserContent`/multimodal handling | `pydantic_evals/pydantic_evals/evaluators/llm_as_a_judge.py:262` |
| Regression test suites | Judge prompt invariants: concise-reason instruction and few-shot order pinned via `FunctionModel` capture | `tests/evals/test_llm_as_a_judge.py:60`, `tests/evals/test_llm_as_a_judge.py:118` |
| Regression test suites | `LLMJudge`/`GEval` unit tests with mocked grading outputs and model-settings passthrough assertions | `tests/evals/test_evaluator_common.py:259`, `tests/evals/test_evaluator_common.py:512` |
| Regression test suites | Dataset/evaluator/report/OTEL/online/multi-run suites (20 files) using `inline-snapshot` goldens | `tests/evals/test_dataset.py:1`, `tests/evals/test_evaluators.py:1`, `tests/evals/test_reporting.py:187`, `tests/evals/test_online.py:1`, `tests/evals/test_multi_run.py:1` |
| Regression test suites | `GEval` construction validation (`score_range`, `evaluation_steps`) failing at `from_file` time | `pydantic_evals/pydantic_evals/evaluators/common.py:314`, `tests/evals/test_evaluator_common.py:552` |
| CI prompt test integration | CI matrix runs `coverage run -m pytest -n logical --dist=loadgroup` across Python 3.10-3.14 and `pydantic-evals` package | `.github/workflows/ci.yml:337`, `.github/workflows/ci.yml:238` |
| CI prompt test integration | No dedicated prompt-golden/live-LLM CI job; judge tests mock `AbstractAgent.run`, real LLM cassettes only in `tests/models/` | `tests/evals/test_llm_as_a_judge.py:165`, `tests/evals/test_evaluator_common.py:259` |
| System prompt handling | `SystemPromptRunner` / `resolve_system_prompts` with `dynamic` vs static parts (agent framework) | `pydantic_ai_slim/pydantic_ai/_system_prompt.py:15`, `pydantic_ai_slim/pydantic_ai/_system_prompt.py:40` |
| System prompt handling | `format_as_xml` helper for semi-structured LLM context | `pydantic_ai_slim/pydantic_ai/format_prompt.py:20` |

## Answers to Dimension Questions

**1. Are prompt changes tested?** Partially. Framework-level prompts are tested; application prompts are not automatically tested. The only prompts whose change is regression-tested by the repo are the LLM-judge system prompts: `tests/evals/test_llm_as_a_judge.py:60` asserts every judge agent's system prompt contains the concise-reason instruction, and `tests/evals/test_llm_as_a_judge.py:118` asserts runtime `_build_prompt` order matches each agent's few-shot examples (`Input → Output → ExpectedOutput → Rubric`). `GEval` and `LLMJudge` also have unit tests for serialization and argument validation (`tests/evals/test_evaluator_common.py:569`). For user-built agents, prompt testing is delegated to the consumer: they must author `Dataset`/`Case` goldens and evaluators; the framework provides the harness but no pre-built golden corpora or default system-prompt regression suite for `Agent` instructions.

**2. Are experiments tracked?** As a library, not as a service. `Dataset.evaluate` returns an `EvaluationReport` (`pydantic_evals/pydantic_evals/dataset.py:396`) containing per-case `ReportCase` objects, `EvaluatorFailure`s, metrics/attributes, and OTel `trace_id`/`span_id`. Comparison is via `EvaluationRenderer.build_diff_table` / `build_table` (`tests/evals/test_reporting.py:187`) with numeric diff formatting (`pydantic_evals/pydantic_evals/reporting/render_numbers.py:62`). Report-level analyses (`pydantic_evals/pydantic_evals/reporting/analyses.py:62`) and `ReportEvaluator`s (`pydantic_evals/pydantic_evals/dataset.py:1013`) extend reports. Online evaluation (`pydantic_evals/pydantic_evals/online.py:496`) emits per-call `gen_ai.evaluation.result` OTel events, optionally to Logfire for visualization, with sampling and concurrency controls. There is no managed experiment database, run versioning, or artifact store in-repo; persistence is via user-chosen sinks/files (`Dataset.to_file` at `pydantic_evals/pydantic_evals/dataset.py:747`).

**3. Is LLM-as-judge used for evaluation?** Yes — first-class. Four judge functions map to input/output combinations (`pydantic_evals/pydantic_evals/evaluators/llm_as_a_judge.py:68`, `:109`, `:154`, `:200`) plus a G-Eval variant (`pydantic_evals/pydantic_evals/evaluators/llm_as_a_judge.py:321`). The `LLMJudge` evaluator (`pydantic_evals/pydantic_evals/evaluators/common.py:224`) wraps them with configurable `score`/`assertion` outputs, and `GEval` (`pydantic_evals/pydantic_evals/evaluators/common.py:288`) implements chain-of-thought scoring with explicit `criteria` + `evaluation_steps` and `score_range` validation. Both default to `openai:gpt-5.2` (`pydantic_evals/pydantic_evals/evaluators/llm_as_a_judge.py:26`) overridable via `set_default_judge_model` (`pydantic_evals/pydantic_evals/evaluators/llm_as_a_judge.py:220`) or per-evaluator `model`/`model_settings`. Supported in `DEFAULT_EVALUATORS` and thus YAML-serializable (`pydantic_evals/pydantic_evals/evaluators/common.py:358`).

**4. Are regressions caught before deployment?** For the framework itself, yes via CI; for downstream prompts, only if the downstream project wires it. `.github/workflows/ci.yml:337` runs `pytest` across all eval suites with `inline-snapshot` goldens, so a prompt-template formatting change that alters output triggers a snapshot failure (blocking merge via the `check` gate at `.github/workflows/ci.yml:602`). The judge invariants add an extra prompt-specific safety net. What is NOT caught automatically: (a) live LLM quality regressions (judge tests mock the model at `tests/evals/test_llm_as_a_judge.py:165`; no nightly live-judge run), (b) user-agent instruction drift (no bundled golden dataset), (c) A/B experiment regressions (no CI baseline diff step — `build_diff_table` is a library feature, not a workflow step). The library docs explicitly frame evals as an emerging, code-first practice without prescribing CI wiring (`docs/evals.md:14`).

## Architectural Decisions

- **Code-first eval framework over managed platform** (`docs/evals.md:11`, `docs/evals.md:51`): datasets/cases/evaluators are pure Python + YAML, run locally via `dataset.evaluate_sync(task)`. Enables type safety and version control, trades off against hosted experiment history/dashboards (delegated to Logfire/OTel).
- **Generic `Dataset[InputsT,OutputT,MetadataT]` with Pydantic serialization** (`pydantic_evals/pydantic_evals/dataset.py:177`): strongly typed inputs/outputs, JSON schema generation for dataset files, `custom_evaluator_types` registry for extensibility. Evaluator specs are serialized as `{name, arguments}` with short-form sugar (`pydantic_evals/pydantic_evals/evaluators/spec.py`).
- **Evaluator taxonomy: deterministic + LLM + span-based + report-level** (`pydantic_evals/pydantic_evals/evaluators/common.py:358`, `pydantic_evals/pydantic_evals/evaluators/agentic.py:190`, `pydantic_evals/pydantic_evals/reporting/analyses.py:62`): `Equals`/`Contains`/`IsInstance`/`MaxDuration` for cheap checks; `LLMJudge`/`GEval` for subjective quality; `ToolCorrectness`/`TrajectoryMatch`/`ArgumentCorrectness`/`MaxToolCalls`/`MaxModelRequests` for agentic behavior via `SpanTree`; `ReportEvaluator` for cross-case analyses.
- **LLM-as-judge as typed Agents with structured outputs** (`pydantic_evals/pydantic_evals/evaluators/llm_as_a_judge.py:46`): each judge is a `pydantic_ai.Agent` with `output_type=GradingOutput`/`GEvalOutput`, producing `EvaluationReason` with `value` + `reason`. G-Eval simplified from logprob-weighted expectation to direct integer scoring (`pydantic_evals/pydantic_evals/evaluators/llm_as_a_judge.py:332`).
- **OTel/logfire as observability substrate** (`pydantic_evals/pydantic_evals/dataset.py:344`, `pydantic_evals/pydantic_evals/online.py:393`): per-case and per-experiment spans, `gen_ai.*` attributes, cost/token metrics, `SpanTree` context subtree for trace-aware evaluators. Online path also emits `gen_ai.evaluation.result` events.
- **Evaluator versioning for pipeline stability** (`pydantic_evals/pydantic_evals/evaluators/evaluator.py:185`, `pydantic_evals/pydantic_evals/online.py:201`): `get_evaluator_version() -> str|None` propagated to `EvaluationResult.evaluator_version` and OTel events, allowing dashboards to filter retired versions without deleting history.
- **Snapshot-based regression safety** (`tests/evals/test_llm_as_a_judge.py:60`, `tests/evals/test_dataset.py:307`): `inline-snapshot` pins prompt structure and evaluation outputs; snapshot failures block merges via CI `check` job.

## Notable Patterns

- **Typed dataset + evaluator registry**: `Dataset` generics enforce input/output shape across cases, tasks, and evaluators at `pyright --strict` level; `DEFAULT_EVALUATORS` registry (`pydantic_evals/pydantic_evals/evaluators/common.py:358`) allows YAML round-trip for built-ins and `custom_evaluator_types` for user extensions, validated at `from_dict` time with `ExceptionGroup`.
- **Prompt-section ordering as invariant**: `_build_prompt` (`pydantic_evals/pydantic_evals/evaluators/llm_as_a_judge.py:262`) and system-prompt few-shot examples are locked to the same order, verified by `test_build_prompt_section_order_matches_few_shot_examples` — a mechanical guard against prompt drift that silently degrades judge quality.
- **Multimodal prompt handling**: `_make_section` (`pydantic_evals/pydantic_evals/evaluators/llm_as_a_judge.py:240`) preserves `UserContent` (images/audio/video) when present, returning `Sequence[str|UserContent]` instead of stringifying — enabling judges over non-text outputs.
- **Graceful degradation on missing spans**: agentic evaluators return descriptive `EvaluationReason(value=False/0.0, reason=_NO_SPAN_TREE_REASON)` when `ctx.span_tree` raises `SpanTreeRecordingError` (`pydantic_evals/pydantic_evals/evaluators/agentic.py:222`, `:324`), so absence of Logfire does not crash evaluations.
- **Online sampling and backpressure**: per-evaluator `sample_rate` (float or `Callable[[SamplingContext], float|bool]`), `sampling_mode` (`independent` vs `correlated`), `max_concurrency` semaphore + `on_max_concurrency`/`on_error` callbacks (`pydantic_evals/pydantic_evals/online.py:193`, `:104`, `:415`).
- **Report diff as first-class UX**: `build_diff_table` with absolute+relative change rendering (`pydantic_evals/pydantic_evals/reporting/render_numbers.py:62`) mirrors how engineers reason about experiment deltas, not just pass/fail.

## Tradeoffs

- **Framework vs platform**: providing a library maximizes flexibility and avoids vendor lock-in, but leaves experiment persistence, history, and A/B comparison to users/OTel backends; teams expecting a hosted tracking store (MLflow/Weights&Biases-style) must build integration.
- **Mocked judge tests for speed determinism**: all judge tests mock `AbstractAgent.run` (`tests/evals/test_llm_as_a_judge.py:165`), so CI is fast and deterministic but does not catch live model behavioral drift (e.g., GPT-5.2 rubric interpretation changes or rate-limit failures).
- **Simplified G-Eval**: direct integer score trades correlation with human judgments (`pydantic_evals/pydantic_evals/evaluators/llm_as_a_judge.py:332` note) for provider-agnostic simplicity (no logprobs dependency, works across OpenAI/Anthropic/Google).
- **No built-in LLM cache for judges**: each evaluation pays a full LLM call; cost/saturation concerns are managed only via `sample_rate`/`max_concurrency` sampling, not response caching or batching.
- **Evaluator registry limited to dataclass-serializable evaluators**: YAML round-trip requires dataclass `Evaluator` subclasses with serializable fields; validators enforce `score_range`/`evaluation_steps` at construction (`pydantic_evals/pydantic_evals/evaluators/common.py:314`) which is good for fail-fast but prevents runtime-computed evaluator configs from surviving serialization.

## Failure Modes / Edge Cases

- **Silent prompt drift in user agents**: `SystemPromptRunner` (`pydantic_ai_slim/pydantic_ai/_system_prompt.py:15`) resolves static + dynamic prompts at runtime without pinning or diffing; a prompt string change ships with no automatic quality gate unless the consumer maintains a `Dataset` golden suite and CI gate.
- **Live-model judge flakiness not caught in CI**: mocked tests hide provider outages, prompt-injection via user content, and rubric ambiguity; `_build_prompt` does not sanitize input/output content and interpolates raw user strings into `<Input>/<Output>` tags, so adversarial inputs could confuse the judge.
- **Model-id staleness**: default `openai:gpt-5.2` (`pydantic_evals/pydantic_evals/evaluators/llm_as_a_judge.py:26`) may not exist/resolve on all provider configs; `_serialize_model_as_string` (`pydantic_evals/pydantic_evals/evaluators/common.py:212`) serializes custom `Model` instances as `model_id` strings, which cannot round-trip custom models (noted as rare-but-unsolved).
- **Dataset load failures are aggregated**: invalid evaluators raise an `ExceptionGroup` with only first 3 errors surfaced (`pydantic_evals/pydantic_evals/dataset.py:738`), so large config typos require iterative fixes.
- **Span-dependent evaluators fail open on missing telemetry**: agentic evaluators degrade to `EvaluationReason(value=False/0.0)` rather than `EvaluatorFailure`, so a misconfigured Logfire could silently score every case as 0 without alerting (distinguishable only by `reason` text).
- **Online evaluation backpressure drops evaluations silently by default**: `max_concurrency=10` (`pydantic_evals/pydantic_evals/online.py:219`) with `on_max_concurrency=None` means saturated evaluators are dropped with no log; `wait_for_evaluations` (`pydantic_evals/pydantic_evals/online.py:957`) exists for tests but is not used in production.
- **Serialization format drift**: dataset YAML files embed `$schema` references and short-form evaluator syntax; a schema change could break existing dataset files without migration tooling.
- **Budget exhaustion**: report-level metrics rely on `_task_run.extract_span_tree_metrics` (`pydantic_evals/pydantic_evals/dataset.py:998`); providers whose instrumentation does not set `gen_ai.request.model`/`gen_ai.operation.name` will produce `requests=0` costs, undercounting budgets checked by `MaxToolCalls`/`MaxModelRequests`.

## Future Considerations

- Add a managed run registry or filesystem-backed experiment store (e.g., hashing `Dataset` + task version) so `EvaluationReport` history survives across invocations without user scaffolding.
- Add a CI helper (`pydantic-evals diff --baseline last_success.json`) that fails the build on assertion/score regression, turning `build_diff_table` into a deploy gate.
- Introduce prompt versioning for `Agent`/`SystemPromptRunner` parallel to `get_evaluator_version`, so online sinks can filter by prompt version and dashboards can correlate prompt changes to quality shifts.
- Add response caching / deterministic replay for `LLMJudge`/`GEval` to reduce cost on large datasets and enable offline diffing; key on `(rubric, inputs, output)` hash.
- Expand judge robustness: input sanitization/escaping for `<Input>`-style tags, stronger validation of `GEval` score-range compliance (already raised at `pydantic_evals/pydantic_evals/evaluators/llm_as_a_judge.py:376`), and optional logprob-weighted G-Eval for providers that support it.
- Add a nightly live-LLM smoke job (e.g., `test_llm_as_a_judge_live`) behind an API-key gate to catch model drift between releases.

## Questions / Gaps

- No evidence found of built-in A/B or multi-variant experiment orchestration (e.g., running two agent variants against one dataset and ranking winner); multi-run is limited to repeating the *same* task via `repeat` at `pydantic_evals/pydantic_evals/dataset.py:269`. Search covered `pydantic_evals/` and `docs/evals/`.
- No evidence found of dataset version pinning or golden-output lineage; `Dataset.to_file` writes snapshots but no hash/version field ties a report back to the exact dataset revision it ran against.
- No evidence found of CI-integrated prompt regression tests for `pydantic_ai` agent example instructions or `format_as_xml` templating; only judge prompts are pinned.
- Unknown: intended Logfire vs generic OTel sink split for `EvaluationSink` — `pydantic_evals/pydantic_evals/online.py:393` supports both, but docs do not prescribe production experiment-tracking topology.
- Unknown: whether `pydantic-evals[logfire]` is required for `HasMatchingSpan`/`ToolCorrectness` etc. to produce non-zero scores; docs show span-based evaluation needs OTel but CI skip logic (`tests/evals/test_evaluator_common.py:463` `@needs_logfire`) suggests non-Logfire runs silently skip.

---

Generated by `Dimension 12.03: Prompt Evaluation and Experiments` against `pydantic-ai`.
