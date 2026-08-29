# Source Analysis: crewai

## Trajectory Evaluation

### Source Info

| Field | Value |
|-------|-------|
| Name | crewai |
| Path | `studies/agent-harness-study/sources/crewai` |
| Language / Stack | Python / Pydantic, LiteLLM, event-bus, experimental evaluation harness |
| Analyzed | 2026-08-29 |

## Summary

CrewAI implements trajectory evaluation via two parallel subsystems: a lightweight production `CrewEvaluator`/`TaskEvaluator` that scores final task outputs (1–10) via LLM-as-judge, and a richer **experimental** `AgentEvaluator` framework that consumes a reconstructed execution trajectory (tool uses + LLM calls) to score intermediate dimensions. The experimental harness is the only place where intermediate steps, tool choice, and reasoning efficiency are evaluated; it captures traces through `EvaluationTraceCallback` listening on `ToolUsage*Event` and `LLMCall*Event` streams, then dispatches six metric evaluators (`GoalAlignment`, `SemanticQuality`, `ToolSelection`, `ParameterExtraction`, `ToolInvocation`, `ReasoningEfficiency`). Aggregation across tasks and iterations produces trajectory-level scores with strategy variants (simple average, best/worst, weighted). Context-usage and recovery are only **indirectly** measured (token counts collected but not scored; tool error rates fed to LLMs but no explicit retry/recovery success metric). The design is present and tested but marked `experimental`, LLM-dependent, and fragile to parsing failures.

## Rating

**5/10 — Present but inconsistent, weakly documented, and fragile**

Rationale: trajectory-aware evaluators exist and are tested, but live under `crewai.experimental.evaluation` (`lib/crewai/src/crewai/experimental/evaluation/__init__.py:1-48`), requiring an LLM for every metric and returning `None` on parse failure rather than a safe fallback. Intermediate-step evaluation is real (tool selection / parameter extraction / invocation / reasoning) yet disjoint from the production `Crew.test()` path (`lib/crewai/src/crewai/crew.py:2227-2253`) which only uses `CrewEvaluator`'s final-output scorer. No explicit context-usage metric, no recovery-success metric, and no bounded trajectory scoring outside the experimental runner. Tests cover happy-path and parse-error cases but do not exercise multi-agent retries, context-window exceed, or trace loss.

## Evidence Collected

Every entry MUST include a file path with line numbers. Format: `path/to/file.ts:NN`.

| Area | Evidence | File:Line |
|------|----------|-----------|
| Step-by-step eval code | `BaseEvaluator.evaluate(agent, execution_trace, final_output, task)` abstract contract requiring execution trace | `lib/crewai/src/crewai/experimental/evaluation/base_evaluator.py:62-69` |
| Step-by-step eval code | `MetricCategory` enum defines 6 categories including tool and reasoning sub-metrics | `lib/crewai/src/crewai/experimental/evaluation/base_evaluator.py:20-29` |
| Step-by-step eval code | `AgentEvaluator._handle_task_completed` extracts trace via `callback.get_trace(agent_id, task_id)` and calls `self.evaluate()` per task completion event | `lib/crewai/src/crewai/experimental/evaluation/agent_evaluator.py:79-131` |
| Step-by-step eval code | `AgentEvaluator._handle_lite_agent_completed` does same for LiteAgent path (lite_task) | `lib/crewai/src/crewai/experimental/evaluation/agent_evaluator.py:132-184` |
| Step-by-step eval code | `AgentEvaluator.evaluate()` iterates `self.evaluators` sequentially, emitting started/completed/failed events per metric | `lib/crewai/src/crewai/experimental/evaluation/agent_evaluator.py:248-294` |
| Tool choice evaluators | `ToolSelectionEvaluator` evaluates relevance/coverage using `execution_trace.tool_uses`, available tools list, and LLM JSON response | `lib/crewai/src/crewai/experimental/evaluation/metrics/tools_metrics.py:22-146` |
| Tool choice evaluators | `ParameterExtractionEvaluator` evaluates accuracy/formatting/completeness from `tool_uses[:5]` samples + validation_error_rate | `lib/crewai/src/crewai/experimental/evaluation/metrics/tools_metrics.py:149-302` |
| Tool choice evaluators | `ToolInvocationEvaluator` evaluates structure/error_handling/invocation_patterns from error_rate and error_type breakdown | `lib/crewai/src/crewai/experimental/evaluation/metrics/tools_metrics.py:305-458` |
| Tool choice evaluators | `create_default_evaluator()` wires all six evaluators: GoalAlignment, SemanticQuality, ToolSelection, ParameterExtraction, ToolInvocation, ReasoningEfficiency | `lib/crewai/src/crewai/experimental/evaluation/agent_evaluator.py:344-365` |
| Context usage metrics | `EvaluationTraceCallback.on_llm_call_end` extracts `total_tokens` from `response.usage.total_tokens` and appends to trace `llm_calls` | `lib/crewai/src/crewai/experimental/evaluation/evaluation_listener.py:292-329` |
| Context usage metrics | `ReasoningEfficiencyEvaluator` computes `total_tokens`, `avg_tokens_per_call`, `efficiency_metrics` but evaluates only via LLM prompt (no explicit context-window metric) | `lib/crewai/src/crewai/experimental/evaluation/metrics/reasoning_metrics.py:59-99` |
| Context usage metrics | No dedicated context-window or knowledge-retrieval evaluator; only comment mentions `knowledge retrievals` in docstring, no implementation | `lib/crewai/src/crewai/experimental/evaluation/evaluation_listener.py:30-36` |
| Recovery eval cases | `ToolInvocationEvaluator` surfaces `error_rate`, `error_types` (execution_error vs usage_error) to LLM but does not measure retry success | `lib/crewai/src/crewai/experimental/evaluation/metrics/tools_metrics.py:320-372` |
| Recovery eval cases | `AgentEvaluator.evaluate` catches evaluator exceptions and emits `AgentEvaluationFailedEvent` + prints error, but does not score recovery behavior | `lib/crewai/src/crewai/experimental/evaluation/agent_evaluator.py:283-292` |
| Recovery eval cases | Guardrail retry logic (`guardrail_max_retries` in `crew.py`, `task.py`, `lite_agent.py`) is not evaluated in any evaluator | `lib/crewai/src/crewai/task.py:279-282`, `lib/crewai/src/crewai/lite_agent.py:276-277` |
| Trajectory scoring | `AgentEvaluationResult` and `AgentAggregatedEvaluationResult` store per-metric `EvaluationScore(score 0-10, feedback, raw_response)` plus `overall_score` | `lib/crewai/src/crewai/experimental/evaluation/base_evaluator.py:72-106` |
| Trajectory scoring | `AggregationStrategy` enum (SIMPLE_AVERAGE, WEIGHTED_BY_COMPLEXITY, BEST/WORST_PERFORMANCE) controls cross-task aggregation | `lib/crewai/src/crewai/experimental/evaluation/base_evaluator.py:80-84` |
| Trajectory scoring | `EvaluationDisplayFormatter._aggregate_agent_results` averages valid scores per category and synthesizes feedback via LLM when >2 tasks | `lib/crewai/src/crewai/experimental/evaluation/evaluation_display.py:250-378` |
| Trajectory scoring | `EvaluationDisplayFormatter.display_summary_results` renders per-iteration table with Avg. Total column | `lib/crewai/src/crewai/experimental/evaluation/evaluation_display.py:103-248` |
| Trajectory scoring | `ExperimentRunner` runs dataset inputs, resets iterations, extracts scores via `_extract_scores`, asserts via `_assert_scores`, returns `ExperimentResults` | `lib/crewai/src/crewai/experimental/evaluation/experiment/runner.py:26-165` |
| Trajectory scoring | `ExperimentResults.compare_with_baseline()` diffs improved/regressed/unchanged/new/missing tests | `lib/crewai/src/crewai/experimental/evaluation/experiment/result.py:48-143` |
| Production trajectory eval (weak) | `CrewEvaluator.evaluate()` creates an `Agent` evaluator-agent task via `Task.execute_sync()` and stores `tasks_scores[iteration]` as simple LLM quality 1-10 | `lib/crewai/src/crewai/utilities/evaluators/crew_evaluator_handler.py:177-221` |
| Production trajectory eval (weak) | `TaskEvaluator.evaluate()` builds `evaluation_query` from task description/expected/actual output, converts to `TaskEvaluation(quality 0-10)` | `lib/crewai/src/crewai/utilities/evaluators/task_evaluator.py:69-113` |
| Production trajectory eval (weak) | `Crew.test(n_iterations, eval_llm)` loops `kickoff` n times, invokes `CrewEvaluator` with iteration tracking and `print_crew_evaluation_result()` | `lib/crewai/src/crewai/crew.py:2227-2272` |
| Reasoning trajectory eval | `ReasoningEfficiencyEvaluator._detect_loops` uses Jaccard similarity >0.7 and `_calculate_loop_likelihood` with repetition heuristics | `lib/crewai/src/crewai/experimental/evaluation/metrics/reasoning_metrics.py:227-383` |
| Trace collection | `EvaluationTraceCallback` singleton registers handlers for `AgentExecutionStarted/Completed`, `ToolUsageFinished/Error*`, `LLMCallStarted/Completed` and builds `traces[agent_id_task_id] = {tool_uses, llm_calls}` | `lib/crewai/src/crewai/experimental/evaluation/evaluation_listener.py:58-146` |
| Tests | `test_tools_metrics.py` covers ToolSelection/ParameterExtraction success, parse-error, and no-tools cases | `lib/crewai/tests/experimental/evaluation/metrics/test_tools_metrics.py:21-92` |
| Tests | `test_reasoning_metrics.py` covers insufficient calls, successful evaluation, parse error handling, and loop detection | `lib/crewai/tests/experimental/evaluation/metrics/test_reasoning_metrics.py:44-188` |
| Tests | `test_agent_evaluator.py` tests iteration set, default evaluator creation (6 evaluators), and failed evaluator swallowing | `lib/crewai/tests/experimental/evaluation/test_agent_evaluator.py:52-274` |
| Tests | `test_experiment_runner.py` mocks `create_default_evaluator` and checks `reset_iterations_results`/`get_agent_evaluation` call counts | `lib/crewai/tests/experimental/evaluation/test_experiment_runner.py:22-198` |

## Answers to Dimension Questions

| Question | Answer | Evidence |
|----------|--------|----------|
| 1. Are intermediate steps evaluated? | **Partially yes — but only in experimental path.** The production `CrewEvaluator`/`TaskEvaluator` scores only final output (`utilites/evaluators/task_evaluator.py:69-113`, `crew_evaluator_handler.py:177-221`). The experimental `AgentEvaluator` does evaluate intermediate steps via `execution_trace` (tool uses + LLM calls) dispatching 6 metric evaluators per task (`agent_evaluator.py:79-131`, `248-294`). Each metric inspects trace samples (e.g., `llm_calls` for reasoning, `tool_uses` for tool metrics). However the experimental harness is not wired to `Crew.kickoff()` by default; it requires manual `create_default_evaluator` + `ExperimentRunner.run()` or explicit event subscription. No single unified trajectory-step score is emitted outside experimental. | `lib/crewai/src/crewai/experimental/evaluation/agent_evaluator.py:79-131,248-294` ; `lib/crewai/src/crewai/crew.py:2227-2253` vs `lib/crewai/src/crewai/experimental/evaluation/experiment/runner.py:32-57` |
| 2. Is tool selection quality measured? | **Yes — explicitly with three evaluators.** `ToolSelectionEvaluator` scores relevance/coverage of tool categories chosen vs available tools (`tools_metrics.py:22-146`). `ParameterExtractionEvaluator` scores accuracy/formatting/completeness of parameter values plus validation_error_rate (`tools_metrics.py:149-302`). `ToolInvocationEvaluator` scores structure/error_handling/invocation_patterns using error_rate/error_type breakdown (`tools_metrics.py:305-458`). All three are included in `create_default_evaluator` (`agent_evaluator.py:356-363`). Sampling is limited to first 5 tool uses and feedback is LLM-generated JSON, so measurement is qualitative, not deterministic precision/recall. | `lib/crewai/src/crewai/experimental/evaluation/metrics/tools_metrics.py:22-458` ; `lib/crewai/src/crewai/experimental/evaluation/agent_evaluator.py:356-363` |
| 3. Is context usage evaluated? | **No — largely absent.** `EvaluationTraceCallback.on_llm_call_end` captures `total_tokens` per LLM call (`evaluation_listener.py:308-329`) and `ReasoningEfficiencyEvaluator` computes `total_tokens`/`avg_tokens_per_call` (`reasoning_metrics.py:68-69,92-95`), but these numbers are only formatted into an LLM prompt for loop/progression scoring, not scored as context-efficiency. There is no evaluator for knowledge retrieval quality, prompt token window utilization, or `respect_context_window` / `ContextWindowExceedingException` handling. The docstring aspirationally mentions "knowledge retrievals" (`evaluation_listener.py:34-35`) but no listener or metric implements it. `CrewEvaluator` and `ExperimentRunner` do not track `UsageMetrics` per trajectory. | `lib/crewai/src/crewai/experimental/evaluation/evaluation_listener.py:308-329,34-35` ; `lib/crewai/src/crewai/experimental/evaluation/metrics/reasoning_metrics.py:68-99` ; `lib/crewai/src/crewai/llm.py:2463-2490` (not referenced by evaluators) |
| 4. Is recovery behavior measured? | **Weakly/indirectly.** Tool error rates are collected (`tool_uses` with `success`, `error_type`: `usage_error`, `execution_error`, `validation_error`, `selection_error` in `evaluation_listener.py:93-137`) and fed to `ParameterExtractionEvaluator`/`ToolInvocationEvaluator` as `validation_error_rate` / `error_rate` textual context. `ToolInvocationEvaluator`'s system prompt asks the judge to score `Error Handling` 0-10 (`tools_metrics.py:376-419`). However there is no metric for *successful recovery*: guardrail retries (`task.py:1343-1389`, `lite_agent.py:752-760`), resumption after failure, or retry-then-success patterns. `AgentEvaluator.evaluate` swallows evaluator exceptions and emits `AgentEvaluationFailedEvent` (`agent_evaluator.py:283-292`) but does not score the agent's own recovery. No test case covers tool failure then recovery. | `lib/crewai/src/crewai/experimental/evaluation/evaluation_listener.py:93-137` ; `lib/crewai/src/crewai/experimental/evaluation/metrics/tools_metrics.py:173-212,320-419` ; `lib/crewai/src/crewai/experimental/evaluation/agent_evaluator.py:283-292` ; `lib/crewai/src/crewai/task.py:1343-1389` |

## Architectural Decisions

| Decision | Description | File:Line |
|----------|-------------|-----------|
| Split production vs experimental evaluation | `CrewEvaluator`/`TaskEvaluator` (simple LLM quality score, stable) vs `AgentEvaluator` (six metric judges, experimental namespace). Production `Crew.test()` only uses former, isolating risk. | `lib/crewai/src/crewai/utilities/evaluators/crew_evaluator_handler.py:29-56` ; `lib/crewai/src/crewai/experimental/evaluation/agent_evaluator.py:47-66` ; `lib/crewai/src/crewai/crew.py:2253` |
| Event-bus trace collection | Singleton `EvaluationTraceCallback` subscribes to `AgentExecution*`, `ToolUsage*`, `LLMCall*` events, building `traces[agent_id_task_id]` dict rather than instrumenting executors directly. Decouples evaluators from execution. | `lib/crewai/src/crewai/experimental/evaluation/evaluation_listener.py:30-56,58-146` |
| LLM-as-judge for all metrics | Every metric evaluator builds a system/user prompt and calls `self.llm.call(prompt)` then `extract_json_from_llm_response`. No rule-based scoring. | `lib/crewai/src/crewai/experimental/evaluation/metrics/tools_metrics.py:112-125` ; `lib/crewai/src/crewai/experimental/evaluation/metrics/goal_metrics.py:70-88` ; `lib/crewai/src/crewai/experimental/evaluation/metrics/reasoning_metrics.py:174-188` |
| Trajectory aggregation as averaging | `AgentAggregatedEvaluationResult` averages valid scores per metric and overall score; supports four strategies via `AggregationStrategy`. No weighted confidence or trajectory cost model. | `lib/crewai/src/crewai/experimental/evaluation/base_evaluator.py:80-106` ; `lib/crewai/src/crewai/experimental/evaluation/evaluation_display.py:250-305` |
| ExperimentRunner dataset loop | `ExperimentRunner.run()` iterates dataset inputs, calls `crew.kickoff(inputs)` or `agent.kickoff(**inputs)` per test case, then `get_agent_evaluation()` and `_assert_scores`. Iteration state reset per test case. | `lib/crewai/src/crewai/experimental/evaluation/experiment/runner.py:32-112` |

## Notable Patterns

| Pattern | Where | Notes |
|---------|-------|-------|
| Evaluator plugin with MetricCategory discriminator | `lib/crewai/src/crewai/experimental/evaluation/base_evaluator.py:20-69` ; `metrics/*.py` | Each evaluator declares `metric_category` property; `AgentEvaluator` stores results in `metrics[MetricCategory]`. Extensible but no registration discovery — caller must pass evaluator list or use `create_default_evaluator`. |
| Execution trace as loose dict | `lib/crewai/src/crewai/experimental/evaluation/evaluation_listener.py:167-329` | Trace is `dict[str, Any]` with `tool_uses: list[dict]`, `llm_calls: list[dict]`. Informal schema; evaluators defensively `.get()` with defaults. Enables serialization but fragile typing. |
| Singleton callback with thread lock in evaluator | `lib/crewai/src/crewai/experimental/evaluation/evaluation_listener.py:40-46` ; `lib/crewai/src/crewai/experimental/evaluation/agent_evaluator.py:61` | `EvaluationTraceCallback.__new__` singleton + `_state_lock` in `AgentEvaluator` suggests concurrency concern but singleton is global process state — cross-test leakage risk. |
| Feedback synthesis via LLM | `lib/crewai/src/crewai/experimental/evaluation/evaluation_display.py:306-378` | When >2 feedbacks per metric, calls `create_llm()` to summarize. Adds latency/cost; fallback is concatenation. |

## Tradeoffs

| Tradeoff | Description |
|----------|-------------|
| Rich trajectory insight vs LLM cost & nondeterminism | Six LLM calls per task (one per evaluator) plus an optional summarization call yields high token cost and variance. Caching not applied to evaluation LLM calls. Benefit is qualitative coverage of reasoning and tool use. |
| Experimental isolation vs discoverability | Keeping trajectory evaluators under `experimental` avoids breaking production but means `Crew.test()` users never get intermediate metrics unless they opt into experimental API; trajectory evaluation is effectively invisible to default users. |
| Flexible unstructured trace vs type safety | Dict-based trace allows rapid evolution but evaluators must handle missing keys (`if not llm_calls or len<2: return None` in `reasoning_metrics.py:61-65`, `if tool_count==0: return None` in `tools_metrics.py:44-51`). No contract enforcement for trace completeness. |
| LLM judge prompts vs deterministic scoring | Tool evaluators avoid deterministic precision/recall; prompts explicitly say "DO NOT evaluate how many times each tool was used" (`tools_metrics.py:79`). This simplifies but prevents quantitative hallucination/error measurement. |
| Singleton trace buffer vs per-run isolation | Global singleton accumulates traces across all evaluations until reset via `reset_iterations_results` (which clears results but not `traces` dict). Long-running processes may grow unbounded. |

## Failure Modes / Edge Cases

| Mode | Behavior | File:Line |
|------|----------|-----------|
| LLM JSON parse failure | All five LLM-based metrics return `EvaluationScore(score=None, feedback="Error...")` and continue; trajectory aggregation silently excludes `None` scores via valid_scores filter, potentially inflating overall score. | `lib/crewai/src/crewai/experimental/evaluation/metrics/tools_metrics.py:141-146,297-302,453-458` ; `lib/crewai/src/crewai/experimental/evaluation/metrics/goal_metrics.py:83-88` ; `lib/crewai/src/crewai/experimental/evaluation/metrics/reasoning_metrics.py:219-225` |
| Missing trace (no events) | `_handle_task_completed` falls back to `trace = {}` if `get_trace` returns None (`agent_evaluator.py:104-105`), then evaluators see empty `tool_uses`/`llm_calls` and return `score=None` with "No tools available" or "Insufficient LLM calls". No retry or warning beyond feedback string. | `lib/crewai/src/crewai/experimental/evaluation/agent_evaluator.py:100-113` ; `lib/crewai/src/crewai/experimental/evaluation/evaluation_listener.py:247-248` |
| Evaluator exception swallowed | `AgentEvaluator.evaluate` catches `Exception`, emits `AgentEvaluationFailedEvent`, prints via `ConsoleFormatter`, and continues to next evaluator, losing partial scoring context. | `lib/crewai/src/crewai/experimental/evaluation/agent_evaluator.py:283-292` |
| Singleton cross-test pollution | `EvaluationTraceCallback` is process-global singleton; traces never cleared except on overwrite of trace_key. Parallel `ExperimentRunner` or concurrent `Crew.test` invocations may interleave traces due to shared `current_agent_id`/`current_task_id`. | `lib/crewai/src/crewai/experimental/evaluation/evaluation_listener.py:40-56,247-265` |
| Sampling truncation | Tool evaluators only sample first 5 tool uses (`tool_uses[:5]` in `tools_metrics.py:189,349`); long trajectories hide late errors. Reasoning evaluator samples 6 calls via heuristic (`reasoning_metrics.py:386-398`), missing mid-trajectory loops. |
| No guardrail-retry accounting | Task/LiteAgent guardrail retries (`task.py:1343-1510`, `lite_agent.py:752-760`) emit no dedicated event consumed by evaluation; evaluator sees final tool errors but cannot distinguish unrecovered vs retried-then-recovered failures. | `lib/crewai/src/crewai/task.py:1343-1389` ; `lib/crewai/src/crewai/experimental/evaluation/evaluation_listener.py:93-137` |
| Context window exceed not evaluated | `LLM.get_context_window_size()` and `ContextWindowExceedingException` exist but evaluation never scores truncation/summarization due to `respect_context_window`. Trajectory could silently drop context. | `lib/crewai/src/crewai/llm.py:2463-2490` ; no evaluator reference |

## Future Considerations

- **Materialize context-usage metric**: add `ContextUtilizationEvaluator` scoring prompt-token efficiency, truncation events, and knowledge-retrieval recall; wire to `UsageMetrics` and `Memory` retrieval hooks already present in `crew.py:652-698`.
- **Deterministic tool precision/recall**: complement LLM judge with rule-based tool selection F1 computed from `trace.tool_uses` vs task-required tool set (extracted from task description or annotation), reducing cost and nondeterminism for `ToolSelectionEvaluator`.
- **Recovery-success evaluator**: instrument guardrail retry count and tool `success` transitions (fail → success) to emit a `recovery_rate` metric; current `ToolInvocationEvaluator` only judges error rate.
- **Promote from experimental**: move `experimental.evaluation` to stable `crewai.evaluation` once parse-failure handling, trace clearing, and per-run isolation (replace singleton with scoped instance) are hardened; update `Crew.test()` to optionally expose `AgentEvaluator` metrics.
- **Trajectory determinism & observability**: expose `ExperimentResults.to_json()` with per-step trace dumps (currently excluded via `exclude={"agent_evaluations"}` in `result.py:37`), and add structured logging of `AgentEvaluation*Event` for external tracing systems.

## Questions / Gaps

| Gap | Search Boundary |
|-----|-----------------|
| No evidence of context usage being scored | Searched `lib/crewai/src/crewai/experimental/**`, `utilities/evaluators/**`, `crew.py`, `llm.py`; only token counts collected, no evaluator scores context efficiency or knowledge retrieval quality. |
| No evidence of recovery trajectory scoring | Searched `evaluation/**` and `task.py`/`lite_agent.py` guardrail retry paths; evaluators mention error_handling but no test case or metric explicitly for retry-then-success vs unrecovered failure. |
| Is `crewai.test -n 3` trajectory-aware? | `Crew.test()` in `crew.py:2227-2272` uses only `CrewEvaluator` (final-output quality), not `AgentEvaluator` trace-based metrics; confirms production trajectory evaluation is absent. |
| Does `TrainingTaskEvaluation` evaluate trajectory? | `task_evaluator.py:40-50,115-187` evaluates training data iterations via aggregated human feedback, not per-step trajectory; outside trajectory evaluation scope. No evidence found for step-level training evaluation. |
| Are trajectory scores persisted beyond display? | `ExperimentResults.to_json()` excludes `agent_evaluations` by default (`result.py:37`); `CrewEvaluator` only prints table via `print_crew_evaluation_result()` (`crew_evaluator_handler.py:95-175`). No persistent trajectory score store or CI gate integration found. |

---
Generated by `Dimension 18.02: Trajectory Evaluation` against `crewai`.
