# Source Analysis: langfuse

## Dimension 18.02: Trajectory Evaluation

### Source Info

| Field | Value |
|-------|-------|
| Name | langfuse |
| Path | `studies/agent-harness-study/sources/langfuse` |
| Language / Stack | TypeScript / Next.js + BullMQ worker + Prisma (Postgres) + ClickHouse |
| Analyzed | 2026-08-29 |

## Summary

Langfuse is an LLM observability/evaluation platform, not an agent harness. It does not model an agent trajectory as a first-class sequence of reasoning steps, tool calls, and recovery attempts. Its evaluation system (`worker/src/features/evaluation/*`, `packages/shared/src/features/evals/*`) scores a single snapshot of data: a whole trace, a single observation, or a dataset item via LLM-as-a-judge or sandboxed code evaluators. Intermediate-step, tool-choice, context-usage and recovery evaluation are only achievable indirectly by configuring per-observation variable mappings and custom prompts.

No built-in trajectory scorer, stepwise reward, tool-choice classifier, or recovery benchmark exists. Context usage is the single mature exception: a set of managed RAG evaluators (Context Relevance, Context Correctness, Context Precision, Context Recall, Faithfulness) measures grounding/retrieval quality. Trajectory scoring, if needed, must be assembled by the user from atomized trace/observation scores and aggregated in the scores UI.

## Rating

**Score: 3 / 10 — Absent / Implicit / Ad-hoc**

Rationale: Langfuse provides a durable, tested LLM-as-a-judge and code-eval execution pipeline with deterministic score IDs, queue-backed retries, and trace-linked persistence, but it evaluates final outputs, not trajectories. There is no intermediate-step graph, no tool-selection metric, no recovery/fixed-after-failure evaluator, and no trajectory-level aggregation. Observation-level evals plus tool-call filter columns allow ad-hoc per-step scoring, but the system never links steps into an ordered path, measures path quality when the final answer is wrong, or reports a trajectory score.

## Evidence Collected

Every entry MUST include a file path with line numbers. Format: `path/to/file.ts:NN`.

| Area | Evidence | File:Line |
|------|----------|-----------|
| Step-by-step eval code | `evaluate()` fetches one trace, extracts variables, calls shared `runLLMAsJudgeEvaluation`, then `completeEvalExecution` — single holistic LLM call, no stepwise loop or subtrace scoring | `worker/src/features/evaluation/evalService.ts:1039-1151` |
| Step-by-step eval code | Shared LLM-as-judge core `runLLMAsJudgeEvaluation()` compiles one prompt from `ExtractedVariable[]`, calls LLM once, validates `ScoreDataTypeEnum.*`, returns `EvalExecutionResult.scores` | `worker/src/features/evaluation/evalService.ts:736-962` |
| Step-by-step eval code | Variable extraction from trace vs dataset_item vs named observation; dataset_item, trace, or `getObservationForTraceIdByName(name)` — can target a single named step (e.g. `tool`/`agent`) but caller picks exactly one observation name, not a sequence | `worker/src/features/evaluation/evalService.ts:1153-1357` |
| Step-by-step eval code | `processObservationEval()` path for observation-level jobs: downloads single `ObservationForEval` from S3, extracts variables, dispatches to LLM-as-judge or code eval — strictly one observation per job, no trajectory join | `worker/src/features/evaluation/observationEval/observationEvalProcessor.ts:94-267` |
| Step-by-step eval code | `extractObservationVariables()` maps `ObservationVariableMapping` to one observation's `input/output/metadata/experiment_item_*` — no iteration over sibling spans, no ordering | `packages/shared/src/server/evals/extractObservationVariables.ts:17-82` |
| Step-by-step eval code | Prompt compilation `compileEvalPrompt()` substitutes `{{var}}` via `parseUnknownToString` — flat string replacement, no step enumeration | `worker/src/features/evaluation/evalRuntime.ts:14-26` |
| Step-by-step eval code | Queue model creates one `JobExecution` per `(traceId, configId, datasetItemId, observationId)` tuple; deduplication and cancellation are per-config, not per-path | `worker/src/features/evaluation/evalService.ts:342-377` |
| Tool choice evaluators | `availableTraceEvalVariables` lists `agent, chain, tool, span, retriever, generation, ...` columns — infra to target a tool observation, but no evaluator judges tool selection correctness | `packages/shared/src/features/evals/types.ts:133-198` |
| Tool choice evaluators | Observation filter columns expose `calledToolNames: arrayOptions` mapped to `tool_call_names` and `toolCalls: number` mapped to `tool_call_count` — usable for routing but no scoring logic | `packages/shared/src/features/evals/observationForEval.ts:332-346` |
| Tool choice evaluators | `ObservationForEval` schema captures `tool_definitions`, `tool_calls`, `tool_call_names`, `tool_call_count`, `provided_model_name`, `model_parameters`, `cost_details` — telemetry present, evaluation absent | `packages/shared/src/features/evals/observationForEval.ts:15-75` |
| Tool choice evaluators | `availableObservationEvalVariableColumns` includes `toolCalls`, `toolDefinitions`, `toolCallNames`, `providedModelName`, `costDetails` as mappable variables — enables custom code eval to inspect tools | `packages/shared/src/features/evals/observationForEval.ts:197-241` |
| Context usage metrics | Managed evaluators `Contextrelevance` (v1) and `Contextcorrectness` judge relevance/correctness of `{{context}}` vs `{{query}}`+`{{ground_truth}}` | `worker/src/constants/managed-evaluators.json:62-85` |
| Context usage metrics | Ragas partner evaluators: `Context Precision`, `Context Recall`, `Faithfulness` (v1+v2) that ask LLM to verify grounding of `{{answer}}` in `{{context}}` and deconstruct claims | `worker/src/constants/managed-evaluators.json:192-242` |
| Context usage metrics | Faithfulness v2 prompt explicitly deconstructs answer into atomic statements and verifies each against context — closest to context-usage measurement | `worker/src/constants/managed-evaluators.json:229-242` |
| Context usage metrics | No evaluator measures multi-step context retention (e.g., context window overflow, truncated history, or observation input/output chaining) | `packages/shared/src/features/evals/types.ts:122-131` |
| Recovery eval cases | Managed evaluators `User Distress`, `User Disagreement`, `Out-of-Scope Request` detect user-side failure signals from `conversation_history` + `last_user_message`/`system_prompt` — symptom detectors, not recovery validators | `worker/src/constants/managed-evaluators.json:98-151` |
| Recovery eval cases | No search hit for `recovery`, `retry`, `self-correct`, `fixed.*fail` evaluator definitions; prompts directory contains no recovery/ retry trajectory evaluator | `worker/src/constants/managed-evaluators.json:1-308` |
| Recovery eval cases | Eval queue failure handling maps `LLMCompletionError.isRetryable` to `DELAYED` vs `ERROR` with `retryLLMRateLimitError` — infrastructure retry, not behavioral recovery scoring | `worker/src/queues/evalQueue.ts:203-257` |
| Trajectory scoring | `EvalExecutionResult` is `{scores: CodeEvalScoreWithName[], executionTraceId, metadata}`; `completeEvalExecution` persists each score as `ScoreEventType` with `traceId/observationId` and links first score as `jobOutputScoreId` — no trajectory rollup | `worker/src/features/evaluation/evalCompletion.ts:15-95` |
| Trajectory scoring | `buildDeterministicEvalScoreIds` derives stable ID from `(jobExecutionId, scoreName, occurrenceIndex)` — deduplication per config, not path hash | `packages/shared/src/server/evals/evalScoreIds.ts:6-38` |
| Trajectory scoring | `buildEvalScoreWritePayloads` attaches `executionTraceId` and `executionMetadata {job_execution_id, job_configuration_id, target_trace_id, ...}` to each `score.metadata` — lineage present but aggregate must be external | `worker/src/features/evaluation/evalScoreEvent.ts:16-61` |
| Trajectory scoring | `JobExecution` schema holds `jobInputTraceId`, `jobInputObservationId`, `jobInputDatasetItemId`, `jobOutputScoreId`, `executionTraceId`, `status: PENDING/COMPLETED/ERROR/CANCELLED/DELAYED` — one row per atomic check, not a chain | `packages/shared/prisma/schema.prisma:1083-1119` |
| Trajectory scoring | Code eval `buildCodeEvalPayload` flattens `ExtractedVariable` map into `{observation:{input,output,metadata}, experiment:{itemExpectedOutput, itemMetadata}}` — single observation payload, not trajectory | `packages/shared/src/server/evals/codeEvalExecution.ts:105-128` |
| Experiment scoring | Experiment item scores aggregated via `aggregateScores()` for UI tables — per-run averages over independent `NUMERIC/BOOLEAN/CATEGORICAL` scores, no stepwise trace | `web/src/features/experiments/server/router.ts:418-428` |

## Answers to Dimension Questions

### 1. Are intermediate steps evaluated?

**No — not as a trajectory.** Langfuse can evaluate *a* intermediate observation, but never *the* sequence.

- `worker/src/features/evaluation/evalService.ts:1292-1350` extracts exactly one observation by name (`getObservationForTraceIdByName`) per trace-level job; if the name is missing it throws `UnrecoverableError`. There is no iteration over sibling spans or ordering by `start_time`.
- `worker/src/features/evaluation/observationEval/observationEvalProcessor.ts:177-218` downloads a single serialized `ObservationForEval` from S3, parses it, and runs one evaluator over it. Scheduling (`scheduleObservationEvals.ts:1-13`) fires independently for each observation matching a filter (e.g., `type == TOOL`), so each step gets an isolated score.
- The LLM prompt sees a flat concatenation of selected fields, e.g. `"Query: {{query}}\nGeneration: {{generation}}"` (`worker/src/constants/managed-evaluators.json:12`). It never receives the ordered history `thought → tool_call → observation → thought …`.

Net effect: a user could create N evaluators each targeting `tool:toolName`, `agent:agentName`, etc., and score steps individually, but Langfuse does not compose them into a trajectory, weight them, or penalize a wrong path that accidentally reached the right answer.

### 2. Is tool selection quality measured?

**No built-in measurement; infrastructure to build one exists.**

- No managed evaluator, ragas template, or code-eval starter judges tool correctness. Grep over `worker/src/constants/managed-evaluators.json:1-308` shows zero cases for `tool_choice`, `tool selection`, or `function call correctness`.
- Telemetry is rich: `packages/shared/src/features/evals/observationForEval.ts:55-57` stores `tool_calls`, `tool_call_names`, `tool_call_count`; filter columns `packages/shared/src/features/evals/observationForEval.ts:332-346` allow routing on `calledToolNames` and `toolCalls`. Variable columns `packages/shared/src/features/evals/observationForEval.ts:197-216` let a custom code eval read `toolCalls` / `toolCallNames`.
- `worker/src/features/evaluation/evalService.ts:1153-1354` validation schema (`packages/shared/src/features/evals/types.ts:88-107`) permits `langfuseObject: "tool"` + `objectName` + `selectedColumnId: input/output/metadata` with optional `jsonSelector`, so a custom LLM-as-a-judge prompt could be authored to ask "was `get_weather` the right tool?". No such prompt ships.

Verdict: tool choice is observable, filterable, and mappable, but unevaluated out-of-the-box and without a tool-call correctness metric.

### 3. Is context usage evaluated?

**Yes — partially, via RAG-oriented context metrics; not for agent trajectory context.**

- Affirmative evidence:
  - `worker/src/constants/managed-evaluators.json:62-85` `Contextrelevance` + `Contextcorrectness` (score 0-1 on whether `{{context}}` supports answer, grounded in `{{ground_truth}}`).
  - `worker/src/constants/managed-evaluators.json:192-242` ragas suite: `Context Precision` ("was context useful?"), `Context Recall` ("sentence attributable to context?"), `Faithfulness` v1/v2 (deconstruct answer into claims, `verdict 1/0` per claim, `score = #1 / #statements`). Faithfulness v2 (`worker/src/constants/managed-evaluators.json:229-242`) is the most rigorous context-usage signal.
- Limits:
  - These evaluators treat `{{context}}` as a single retrieval blob, not the agent's multi-step memory/history. `packages/shared/src/features/evals/types.ts:122-131` and `packages/shared/src/features/evals/observationForEval.ts:146-152` expose `input/output/metadata` per observation, but no evaluator measures history-window retention, truncated `history` fields, or per-step grounding decay.
  - None of the evaluators measure whether the agent *attended* to the right intermediate tool output when choosing the next step — that would require trajectory awareness.

### 4. Is recovery behavior measured?

**No.**

- No managed evaluator, test case, or metric references retrying, self-correction, fixed after failure, or backtracking. `worker/src/constants/managed-evaluators.json:98-151` detectors (`User Distress`, `User Disagreement`, `Out-of-Scope Request`) identify *symptoms* of failure (profanity, disagreement, out-of-scope ask) from the last user turn, not whether the agent recovered.
- Infrastructure retry (`worker/src/queues/evalQueue.ts:203-257` `isLLMCompletionError(e).isRetryable → DELAYED` with `retryLLMRateLimitError`) retries the *evaluator's* LLM call on 429/5xx, not the agent's behavior.
- No dataset of `traces with intermediate failure → final success` is scored for recovery quality, e.g., "agent chose wrong tool then recovered" vs "asked clarifying question then recovered".

## Architectural Decisions

1. **Single-score-per-execution model** (`worker/src/features/evaluation/evalCompletion.ts:15-45`, `worker/src/features/evaluation/evalScoreEvent.ts:16-61`). Every executor returns `EvalExecutionResult.scores[]`; `completeEvalExecution` uploads each score to S3, enqueues `IngestionJob`, and persists `jobOutputScoreId = scores[0].id`. Decision favors deduplication and idempotency (deterministic `v5(jobExecutionId, scoreName, occurrenceIndex)` at `packages/shared/src/server/evals/evalScoreIds.ts:6-20`) over trajectory rollups. Tradeoff: no first-class path hash or weighted trajectory score.

2. **Trace-level vs observation-level duality** (`worker/src/features/evaluation/evalService.ts:1039-1151` vs `worker/src/features/evaluation/observationEval/observationEvalProcessor.ts:94-267`). Trace evals run in `EvaluationExecution` queue (one prompt over compiled variables); observation evals run in `LLMAsJudgeExecution`/`CodeEval` queues with S3-materialized `ObservationForEval`. Duality cleanly separates holistic and per-span scoring, but leaves the gap — nothing joins spans into a graph.

3. **Declarative variable mapping + JSONPath selector** (`packages/shared/src/features/evals/types.ts:88-107`, `packages/shared/src/features/evals/utilities.ts:78-113`). Users map template variables to `langfuseObject.selectedColumnId` plus optional `$.path`. Enables flexible grounding (e.g., `{{tool.output.result}}`) without code changes. Limitation: flat key-value, not ordered sequence expansion like `{{steps[*].tool_calls}}`.

4. **Pluggable executors behind `EvalExecutionDeps`** (`worker/src/features/evaluation/evalExecutionDeps.ts:116-262`). `callLLM`, `fetchModelConfig`, `uploadScore`, `enqueueScoreIngestion`, `updateJobExecution` are injectable — `createMockEvalExecutionDeps` is used in `worker/src/__tests__/evalService.test.ts:38-49`. Decision aids testability and the Lambda vs local code-eval dispatchers, but no `TrajectoryScorer` interface exists.

5. **Filter-driven triggering with `InMemoryFilterService`** (`worker/src/features/evaluation/evalService.ts:421-484`, `worker/src/features/evaluation/observationEval/fetchObservationEvalConfigs.ts:1-13`). Evaluators run when trace/observation filters match (`filter: Json` at `packages/shared/prisma/schema.prisma:1063`). Efficient for per-entity scoring; inadequate for trajectory predicates like "tool A failed then tool B succeeded".

## Notable Patterns

- **LLM-as-a-judge + structured output** (`worker/src/features/evaluation/evalService.ts:801-930`): `PersistedEvalOutputDefinitionSchema` + `compilePersistedEvalOutputDefinition(...).outputResultSchema` enforces `ScoreDataTypeEnum.NUMERIC|BOOLEAN|CATEGORICAL` with `reasoning` field; every managed prompt ends with `Think step by step` and example few-shots (`worker/src/constants/managed-evaluators.json:12-96`).
- **Sandboxed code evaluators** (`packages/shared/src/server/evals/codeEvalDispatchers.ts:1-13`, `packages/shared/src/server/evals/codeEvalExecution.ts:187-290`, `scripts/code-eval-runners/python/code_based_eval_handler.py:1-13`): local insecure or AWS Lambda dispatch with `CODE_EVAL_SOURCE_MAX_BYTES=256KB`, `PAYLOAD_MAX=5.5MB`, `RESULT_MAX=256KB` guards (`packages/shared/src/server/evals/codeEvalDispatcherTypes.ts:4-6`). Gives users an escape hatch to implement trajectory logic in Python/TS, but no framework helpers to load the full trace.
- **Trace-linked score lineage** (`worker/src/features/evaluation/evalScoreEvent.ts:40-50`): each score stores `traceId`, `observationId`, `executionTraceId` (W3C id derived from `jobExecutionId` at `worker/src/queues/evalQueue.ts:203`), `metadata={job_execution_id, job_configuration_id, target_trace_id...}`. Adequate for debugging a single eval invocation (EU-style `langfuse-eval` internal traces), not for cross-step correlation.
- **Sharded BullMQ with secondary queue redirect** (`worker/src/queues/evalQueue.ts:118-156`, `worker/src/app.ts:417`): evaluation execution can be sharded by `projectId-jobExecutionId` and redirected to `SecondaryEvalExecutionQueue` via allowlist env var — scales independent judge calls but not trajectory graph computation.

## Tradeoffs

- **Holistic judge vs stepwise transparency:** One LLM call per trace is cheap, deterministic, and easy to retry (`maxRetries:1` at `worker/src/features/evaluation/evalExecutionDeps.ts:216`), but discards signal about whether a correct final answer came from a sound path (e.g., lucky guess after wrong tool). Adversary: penalizing good answers with bad reasoning is impossible without per-step scores.
- **Per-observation isolation vs trajectory coherence:** Evaluating one `ObservationForEval` at a time avoids N-way joins over ClickHouse events and keeps S3 payload small (`CODE_EVAL_DISPATCH_PAYLOAD_MAX_BYTES=5.5MB` at `packages/shared/src/server/evals/codeEvalDispatcherTypes.ts:5`). Cost is loss of cross-step context (e.g., comparing `tool_calls[0]` vs `tool_calls[1]`).
- **Generic variable mapping vs opinionated trajectory schema:** Mapping any `langfuseObject.selectedColumn` with JSONPath supports RAG, chat, and tool observability equally. An opinionated `trajectory: {steps: {tool, observation, reasoning}[]}` schema would enable first-class path scoring but would couple Langfuse to an agent framework.
- **User-authored trajectory logic vs built-in policy:** Code evals let users brute-force any trajectory metric (fetch trace events via clickhouse client inside Python), but without helpers (`getTraceById` is internal to worker, not exposed to code runner payload which only has `payload.observation`). Users must re-implement trace fetching.

## Failure Modes / Edge Cases

- **Silent omission when observation name mismatches** (`worker/src/features/evaluation/evalService.ts:1304-1341`): `getObservationForTraceIdByName` takes first match; duplicate generation names are ignored (`observations.shift()`), so a trajectory with two identical tool names silently scores only the first.
- **Stale/missing observation throws UnrecoverableError** (`worker/src/features/evaluation/evalService.ts:1335-1342`, `1339-1342`): `ObservationNotFound` → job retried with exponential backoff (`retryObservationNotFound` at `worker/src/queues/evalQueue.ts:59-68`), then `ERROR`/`CANCELLED`. A trajectory where step N's observation is delayed (ClickHouse replication lag) incurs spurious cancellations for step N+1 evaluators.
- **Deduplication masks trajectory evolution** (`worker/src/features/evaluation/evalService.ts:364-377`): `findMatchingJob` dedupes on `(configId, datasetItemId, observationId)`; re-running an evaluator on corrected trajectory data does not create a new `JobExecution`, requiring manual `CreateEvalQueue` replay (`worker/src/queues/evalQueue.ts:98-116`).
- **Sampling silently drops trajectory steps** (`worker/src/features/evaluation/evalService.ts:622-630`): `sampling < 1` randomly skips jobs. In a trajectory, skipping a single step breaks path completeness; no compensating sample-balancing across steps.
- **Payload limits truncate long contexts** (`packages/shared/src/server/evals/codeEvalDispatcherTypes.ts:166-185`): `assertDispatchInputWithinLimits` throws `SOURCE_TOO_LARGE` / `PAYLOAD_TOO_LARGE` if compiled variables or source exceed caps. A trajectory with many tool outputs easily exceeds `5.5MB`; evaluator fails with `INVALID_RESULT` rather than degrading gracefully.
- **Score type mismatch yields no trajectory weight** (`worker/src/features/evaluation/evalService.ts:948-998`): `toNormalizedScores` handles `NUMERIC|BOOLEAN|CATEGORICAL`; trajectory roll-up across mixed types would require custom normalization not provided.
- **Filter-only routing cannot express sequence predicates** (`packages/shared/src/features/evals/observationForEval.ts:253-346`): filter columns support `tool_call_count > 0` but not "tool A before observation B" — a trajectory that recovers after a failed tool call is indistinguishable from one that never failed.

## Future Considerations

- Introduce a `Trajectory` evaluator type: payload is ordered `ObservationForEval[]` for a trace (sorted by `start_time`), variable mapping supports `{{steps[*].tool_calls}}` / `{{steps[*].output}}` with JSONPath array expansion, and output is both per-step scores and a rolled-up `trajectory_score` (e.g., `faithful_steps / total_steps`, `recovery_score`).
- Ship reference trajectory prompts/managed evaluators: tool-choice correctness (was `tool_calls[0]` optimal given `input`?), context retention (did step N use output of step N-1?), recovery (failed tool → corrected tool or clarifying ask).
- Expose trajectory helpers to code evals: include optional `trace.observations[]` in `CodeEvalPayload` or provide a read-only `fetchTrajectory(traceId)` helper bound to the project, guarded by the same 5.5MB limit with truncation strategy.
- Persist `trajectoryId`/`stepIndex` on scores or a new `trajectory_scores` table to enable `groupBy trajectory` aggregations in `packages/shared/src/server/repositories/dataset-run-items.ts:275` style CTEs and UI digests (similar to `aggregateScores` at `web/src/features/scores/lib/aggregateScores.ts:1-13`).
- Add recovery benchmark dataset pattern: dataset items tagged `requires_recovery: true` where initial trace intentionally fails validation; evaluator measures whether final trace achieves goal after at least one intermediate `level=ERROR` or `tool_call_count` change.

## Questions / Gaps

- No evidence of trajectory-level scoring, stepwise reward, or path-quality metric in `worker/src/features/evaluation/*`, `packages/shared/src/features/evals/*`, or `worker/src/constants/managed-evaluators.json:1-308`. Confirmed by absence of `trajectory`/`path`/`step.*score` evaluators and by single-entity payload shapes in `packages/shared/src/server/evals/codeEvalExecution.ts:105-128` and `worker/src/features/evaluation/evalScoreEvent.ts:24-50`. If trajectory evals exist outside `sources/langfuse` (e.g., enterprise `ee/` or external service), they were not in scope per isolation rules.
- Could Langfuse's `InAppAgent` (`packages/shared/prisma/schema.prisma:218-277` `InAppAgentConversation/Run/Event`) be the intended trajectory surface? `InAppAgentEvent` stores arbitrary JSON per conversation sequence, but no evaluator consumes it — worth confirming with product owners.
- LLM-as-a-judge structured output (`EvalTemplate.outputDefinition`) supports `CATEGORICAL` multi-label (`matches: string[]` at `worker/src/features/evaluation/evalService.ts:994-998`), which could encode step labels, but no documented trajectory label taxonomy was found.

---

Generated by `Dimension 18.02: Trajectory Evaluation` against `langfuse`.
