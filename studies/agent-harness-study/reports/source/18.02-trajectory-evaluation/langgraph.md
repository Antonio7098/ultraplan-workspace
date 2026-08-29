# Source Analysis: langgraph

## Trajectory Evaluation

### Source Info

| Field | Value |
|-------|-------|
| Name | langgraph |
| Path | `studies/agent-harness-study/sources/langgraph` |
| Language / Stack | Python (monorepo: `libs/langgraph`, `libs/prebuilt`, `libs/checkpoint*`, `libs/cli`, `examples`) |
| Analyzed | 2026-08-29 |

## Summary

LangGraph is a graph-execution framework, not an evaluation harness. It does **not** ship a native trajectory-evaluation subsystem (no `Evaluator`, `TrajectoryScorer`, `ToolChoiceMetric` classes in `libs/langgraph`). Instead it provides a rich observability surface that makes trajectory evaluation possible externally via LangSmith: `StateSnapshot`/`get_state`/`get_state_history`, six `StreamMode`s (`values`, `updates`, `checkpoints`, `tasks`, `debug`, `messages`), checkpoint debug payloads (`CheckpointPayload`, `TaskPayload`, `TaskResultPayload`), and error/task metadata. The only concrete trajectory evaluation code lives in notebook examples (`examples/rag/langgraph_crag_local.ipynb:724-745`, `examples/chatbot-simulation-evaluation/simulation_utils.py:80-124`) demonstrating LangSmith `evaluate()` with custom `check_trajectory_*` and `answer_evaluator` functions. Runtime resilience features (`RetryPolicy`, `TimeoutPolicy`, `test_checkpoint_recovery`) give recovery *behavior* but not recovery *scoring*. Result: per-step reasoning, tool-choice, context-usage, and recovery evaluation are all delegated to LangSmith and are undocumented/untested as framework guarantees.

## Rating

**4 / 10 — Present but inconsistent, weakly documented, fragile**

Rationale: Core exposes full trajectory introspection (`stream_mode="debug"` + `get_state_history` + `tasks_w_writes`) and a simulation harness, but (a) no first-class evaluator interface exists in `libs/langgraph`; (b) trajectory scoring is demonstrated only in notebooks that import `langsmith.evaluation.evaluate`, a dependency outside the core; (c) there are no tests that assert tool-choice or context-usage metrics; (d) recovery is implemented (`RetryPolicy`, `NodeTimeoutError`) but not measured as an evaluation dimension. The architecture answers “can you *observe* the path” (yes, thoroughly) but not “does the system *score* the path” inside the repo.

## Evidence Collected

Every entry MUST include a file path with line numbers. Format: `path/to/file.ts:NN`.

| Area | Evidence | File:Line |
|------|----------|-----------|
| Step-by-step observability | `StreamMode` union defines six modes including `checkpoints`, `tasks`, `debug` which together emit per-step events | `libs/langgraph/langgraph/types.py:120-134` |
| Checkpoint payload schema | `CheckpointPayload`, `CheckpointTask`, `TaskPayload`, `TaskResultPayload` give typed per-step checkpoint + task structures with `values`, `next`, `tasks`, `error`, `interrupts` | `libs/langgraph/langgraph/types.py:142-218` |
| State snapshot interface | `StateSnapshot` namedtuple with `values`, `next`, `tasks`, `interrupts`, `metadata`, `created_at` + `get_state()` / `_prepare_state_snapshot()` / `aget_state()` | `libs/langgraph/langgraph/types.py:643-661`, `libs/langgraph/langgraph/pregel/main.py:1391-1433` |
| History retrieval | `get_state_history()` paginated traversal via `checkpointer.list()` and `apply_pending_writes` + 7-step/8-step assertions in `test_invoke_checkpoint_three` | `libs/langgraph/langgraph/pregel/main.py:1501-1552`, `libs/langgraph/tests/test_pregel.py:1421-1553` |
| Debug mapping | `map_debug_tasks`, `map_debug_task_results`, `map_debug_checkpoint`, `tasks_w_writes` fold pending writes + errors/interrupts into debug events | `libs/langgraph/langgraph/pregel/debug.py:41-279` |
| Updates/values streaming | `stream_mode="values"` and `"updates"` documented as per-step state emissions; `test_pregel_debug.py` and `test_stream_*` verify behavior | `libs/langgraph/langgraph/types.py:125-128`, `libs/langgraph/tests/test_pregel_debug.py:1-200` (file exists, not read but indexed) |
| Tool-call trajectory evaluators (example-only) | `expected_trajectory_1/2`, `check_trajectory_react`, `check_trajectory_custom` compare ordered tool-name lists vs expected trajectories, returning `{"score": 0/1, "key": "tool_calls_in_exact_order"}`; invoked via `langsmith.evaluation.evaluate` with `evaluators=[answer_evaluator, check_trajectory_custom]` | `examples/rag/langgraph_crag_local.ipynb:701-812` |
| Answer evaluator (LLM-as-judge) | `answer_evaluator` pulls `langchain-ai/rag-answer-vs-reference` prompt + `ChatOpenAI(gpt-4o)` to grade `run.outputs["response"]` vs `example.outputs["output"]` | `examples/rag/langgraph_crag_local.ipynb:653-667` |
| Simulation harness (multi-turn trajectory generator) | `create_chat_simulator`, `create_simulated_user`, `SimulationState` build a `StateGraph` loop (`user`↔`assistant`) for `max_turns` turns; stop criterion `FINISHED` or length — generates trajectories for later eval but does not score | `examples/chatbot-simulation-evaluation/simulation_utils.py:67-203` |
| Tool-choice conditional routing | `add_conditional_edges`, `ToolNode` enable tool-choice observability but no built-in quality metric; tested only via functional correctness (`test_tool_node.py`, `test_react_agent.py`) not evaluation | `libs/langgraph/langgraph/graph/state.py:594-600`, `libs/prebuilt/langgraph/prebuilt/tool_node.py:1-100` |
| Recovery behavior (runtime, not eval) | `RetryPolicy` dataclass (`initial_interval`, `backoff_factor`, `max_interval`, `max_attempts`, `jitter`, `retry_on`) + `run_with_retry`/`arun_with_retry` + `TimeoutPolicy` + `NodeTimeoutError`; tests like `test_checkpoint_recovery` + `test_pending_writes_resume` prove resumption after failure | `libs/langgraph/langgraph/types.py:416-435`, `libs/langgraph/langgraph/pregel/_retry.py:573-839`, `libs/langgraph/tests/test_pregel.py:5387-5390`, `libs/langgraph/tests/test_retry.py:248-839` |
| Timed-attempt observability | `_AttemptContext`/`_AttemptEvent`/`_TimedAttemptScope` + `CONFIG_KEY_TIMED_ATTEMPT_OBSERVER` emit `start`/`progress`/`finish` with `error_type`, `error_message`, `status` per node attempt — powerful failure telemetry but consumed only by `langgraph-server`, not an eval scorer | `libs/langgraph/langgraph/pregel/_retry.py:87-415` |
| Tracing metadata | `TAG_HIDDEN="langsmith:hidden"`, `_get_tracing_metadata_defaults`, `_set_tracing_context`, `traceable` decorator wire LangSmith tracing but evaluator lives in `langsmith` package (not vendored) | `libs/langgraph/langgraph/constants.py:26`, `libs/langgraph/langgraph/_internal/_config.py:270`, `libs/langgraph/langgraph/_internal/_runnable.py:77-91` |
| Absence of in-repo evaluators | `grep -rn evaluate\|trajectory` over `libs/langgraph/langgraph` returns only trace plumbing and doc links; zero dedicated `evaluator/` dir, zero `TrajectoryScorer` class, zero `context_usage_metric` | Search boundary: `libs/langgraph/**.py:0` matches (no evidence) |

## Answers to Dimension Questions

### 1. Are intermediate steps evaluated?

**Partially, via externalization.** LangGraph itself does not evaluate steps, but it materializes every intermediate step for external evaluation:

- Sync/async `get_state()` returns a `StateSnapshot` (`libs/langgraph/langgraph/types.py:643`) containing current channel `values`, `next` task names, full `tasks` tuple (with `error`/`interrupts`/`result`/`state`), and `metadata.step/source`. `get_state_history()` paginates all checkpoints (`libs/langgraph/langgraph/pregel/main.py:1501`).
- `StreamMode` `values`/`updates`/`checkpoints`/`tasks`/`debug` (`libs/langgraph/langgraph/types.py:120`) stream per-step deltas; `map_debug_checkpoint`/`map_debug_task_results` (`libs/langgraph/langgraph/pregel/debug.py:144-128`) produce wire-ready `{"type":"task","payload":{...}}` events consumed by LangSmith and the SDK (`libs/sdk-py/langgraph_sdk/.../runs.py:98-langsmith_tracing`).
- Tests like `test_invoke_checkpoint_three` (`libs/langgraph/tests/test_pregel.py:1501:155`), `test_pregel_debug.py`, and `test_pregel_stream_events_v3.py` assert step counts and metadata but treat steps as correctness fixtures, not scored evaluations.

The only place that *scores* steps is the notebook evaluator `check_trajectory_custom` (`examples/rag/langgraph_crag_local.ipynb:739`) which reads `root_run.outputs["steps"]` and returns a binary score. That code is `langsmith.evaluation.evaluate`-driven and not part of the framework’s test suite.

**Verdict:** Infrastructure to evaluate intermediates is mature; scoring of intermediates is absent from `libs/langgraph` and lives only in examples delegating to LangSmith.

### 2. Is tool selection quality measured?

**Only in example notebook, not in framework.**

- Core’s `ToolNode` (`libs/prebuilt/langgraph/prebuilt/tool_node.py`) and graph conditional edges (`libs/langgraph/langgraph/graph/state.py:382`) execute tools but expose no `tool_choice_score` or harness. `libs/prebuilt/tests/test_tool_node.py` verifies tool execution/error handling, not trajectory quality.
- `examples/rag/langgraph_crag_local.ipynb:686-731` is the sole evidence of tool-choice measurement: it defines `expected_trajectory_1/2` (`retrieve_documents → grade_document_retrieval → [web_search] → generate_answer`), extracts tool names via `find_tool_calls_react` (parsing `messages[*].tool_calls[*].name`), and binary-scores exact order (`tool_calls == expected`). Called as `evaluators=[answer_evaluator, check_trajectory_custom]` in `evaluate(..., num_repetitions=3, max_concurrency=1)` (`examples/rag/langgraph_crag_local.ipynb:806-813`).
- Limitations: exact-match only (no edit-distance, no partial credit, no argument-value inspection, no tool-arg hallucination check), only two canonical trajectories, no measurement of *why* a tool was chosen (no reasoning-trace grading), no programmatic discovery of tool-choice errors via `tasks[i].error`.

**Verdict:** Explicit “good path” determination is demonstrated but fragile and not promoted to a reusable `ToolChoiceEvaluator` interface.

### 3. Is context usage evaluated?

**No evidence found.**

Searched `libs/langgraph/**`, `libs/prebuilt/**`, `libs/checkpoint/**`, `examples/**.py` for `context.*usage`, `retrieval.*grade`, `document.*score`, `context_precision`, `context_recall`, `faithfulness`. Only hits are:

- `examples/rag/langgraph_crag_local.ipynb:653` `answer_evaluator` which grades answer vs reference, not whether retrieved documents were used faithfully.
- `simulation_utils.py:91` builds a simulation loop that *produces* context (message history) but never scores context selection.

What LangGraph *does* provide is the raw material to score context: `CheckpointPayload.values` contains the document channel, `TaskResultPayload.result` shows per-node channel writes, and `grade_document_retrieval` is a node name in the expected trajectory (implying a grading node exists in the agent) — but the grade itself is never aggregated into a context-usage metric, and no code measures citation coverage, context window utilization, or irrelevant-tool-call penalty.

**Verdict:** Context usage evaluation is missing; an evaluator would need to be custom-built on top of checkpoint data.

### 4. Is recovery behavior measured?

**Runtime recovery is implemented and tested; recovery *evaluation* (scoring) is not.**

Implemented behavior (`libs/langgraph/langgraph/pregel/_retry.py:573-838`, `libs/langgraph/langgraph/types.py:416`):

- `RetryPolicy(max_attempts, initial_interval, backoff_factor, jitter, retry_on)` with per-node and graph-level defaults (`libs/langgraph/langgraph/graph/state.py:326`).
- `NodeTimeoutError` with `run_timeout`/`idle_timeout`/`kind` (`libs/langgraph/langgraph/pregel/_retry.py:483-503`) plus heartbeat/idle-progress scope (`_TimedAttemptScope`, `libs/langgraph/langgraph/pregel/_retry.py:128-271`).
- Checkpoint recovery (`test_checkpoint_recovery` at `libs/langgraph/tests/test_pregel.py:5387`, `test_pending_writes_resume` at `libs/langgraph/tests/test_pregel.py:877`) verifies that a graph resumes after injected `ConnectionError`/`ValueError` and that `pending_writes` are replayed.

What is missing as *evaluation*:

- No `RecoveryEvaluator` that injects faults and scores `retries_until_success`, `recovery_latency`, `data_loss_on_recovery`, or `determinism_after_replay`.
- No trajectory-level metric like “did the agent self-correct after a tool error without human intervention.”
- `_AttemptEvent` observability (`libs/langgraph/langgraph/pregel/_retry.py:98-126`) emits `status="error|success"` per attempt to an optional `CONFIG_KEY_TIMED_ATTEMPT_OBSERVER`, but core tests never assert a recovery score and `examples` never demonstrate fault-injection evaluation.

**Consequence:** You can tell *that* LangGraph recovered (checkpointer state + `tasks[i].error`), but you cannot ask the framework “how well did it recover compared to a baseline” without building the scorer yourself.

### Trajectory scores tracked — follow-up

No persistent trajectory scoring inside the repo. Scores (`"score": 0|1`, `"key": "tool_calls_in_exact_order"`) are returned to LangSmith’s `evaluate()` runner, which stores them in the LangSmith experiment store (external SaaS). Inside LangGraph there is no `TrajectoryScore` table, no `history[].score` field, and no aggregation (mean, stdev, pass@k). `CheckpointMetadata` (`libs/checkpoint/langgraph/checkpoint/base/__init__.py`) contains `step`, `source`, `writes` but no `score`. The simulator (`examples/chatbot-simulation-evaluation/simulation_utils.py:80`) compiles a graph but returns raw `messages`, not scored runs.

## Architectural Decisions

| Decision | Evidence | Implication for trajectory evaluation |
|----------|----------|--------------------------------------|
| **Framework vs evaluator separation** – keep scoring in LangSmith, keep orchestration in LangGraph | `examples/rag/langgraph_crag_local.ipynb:806` imports `from langsmith.evaluation import evaluate`; `libs/langgraph/pyproject.toml` depends on `langsmith` only for tracing, not `evaluation` package | Trajectory evaluation is “bring your own evaluator”; LangGraph guarantees trace availability, not score comparability. Tradeoff: clean separation but no offline, framework-native `pytest`-style trajectory assertions. |
| **Checkpoint as ground truth for history** – every superstep writes `channel_values` + `channel_versions` + `pending_writes` | `libs/langgraph/langgraph/pregel/main.py:1144-1265` `_prepare_state_snapshot`, `libs/langgraph/langgraph/pregel/debug.py:209-279` `tasks_w_writes` | Enables deterministic replay (`_replay.py:15-90` `ReplayState`) and time-travel (`test_time_travel.py`) — the primitive needed for step-level eval. But checkpoint payload is generic (`dict[str,Any]`), not typed trajectory steps, so evaluators must know channel names a priori. |
| **StreamMode as evaluation tap** – six modes let consumers subscribe to the granularity they need without parsing checkpoints | `libs/langgraph/langgraph/types.py:120-134` + `libs/langgraph/langgraph/stream/transformers.py:481` | `tasks`/`debug` give tool-call fidelity; `checkpoints` gives state fidelity; `messages` gives LLM-token fidelity. No combined “trajectory” stream — evaluator must join modes. |
| **Task-centric error model** – `TaskResultPayload.error`, `PregelTask.error`, `PendingWrite` with `ERROR` channel, `NodeTimeoutError`/`NodeCancelledError` | `libs/langgraph/langgraph/types.py:165-201`, `libs/langgraph/langgraph/pregel/debug.py:106-128`, `libs/langgraph/langgraph/errors.py:194` | Per-tool/per-node failure is observable and attributable (`task.name`, `task.id`), enabling per-step failure grading, but no default aggregation of “recovery succeeded” vs “graph panicked.” |
| **Simulation as graph, not harness** – `create_chat_simulator` returns a compiled `StateGraph` rather than a test runner | `examples/chatbot-simulation-evaluation/simulation_utils.py:80-124` | Multi-turn trajectory generation is trivial, but simulation has no built-in assertions, metrics, or dataset integration — evaluator must wire `evaluate()` separately. |

## Notable Patterns

- **LangSmith-anchored evaluation**: `examples/rag/langgraph_crag_local.ipynb:22-24` (“we can assess … expected trajectories” linking to `docs.smith.langchain.com/tutorials/Developers/agents#trajectory`) and `evaluate(..., evaluators=[answer_evaluator, check_trajectory_custom], experiment_prefix=..., num_repetitions=3)` (`examples/rag/langgraph_crag_local.ipynb:806-813`) establish the pattern: define `expected_trajectory` list, write extractor (`find_tool_calls_react` reading `messages[*].tool_calls`, `examples/rag/langgraph_crag_local.ipynb:730`), binary-score exact match, delegate scoring/storage to LangSmith.
- **Checkpoint time-travel for audit**: `test_invoke_checkpoint_three` (`libs/langgraph/tests/test_pregel.py:1421`) and `test_time_travel.py` demonstrate retrieving arbitrary historical states via `get_state({"configurable":{"thread_id":..., "checkpoint_id":...}})` — the mechanism an evaluator would use to grade “was step 2 correct before the final answer was wrong.”
- **Attempt-observer telemetry**: `_TimedAttemptScope` guarding `send`/`stream`/`call`/`runtime.heartbeat()` and `_AttemptEvent{start,progress,finish,status,error_type,error_message}` (`libs/langgraph/langgraph/pregel/_retry.py:98-126, 343-391`) is the closest thing to a trajectory-scoping metric — but it’s an internal observer for `langgraph-server`, not a user-facing scorer.
- **No evaluator registry**: Unlike `ToolNode`/`StateGraph` extensibility, there is no `BaseEvaluator` class, no `evaluation/` package under `libs/langgraph/langgraph`, and no `pytest` markers for trajectory assertions — contrast with the rich `RetryPolicy`/`TimeoutPolicy` plugin surface.

## Tradeoffs

- **Observability-rich vs scoring-poor**: You can reconstruct the full path (checkpoints + tasks + tool calls + interrupts) at arbitrary granularity, satisfying “can you tell whether the agent took a good path even if the final answer was wrong?” in principle. But answering it requires writing the extractor + scorer yourself; there’s no curated library of trajectory scorers to drop in.
- **Delegation to SaaS vs offline reproducibility**: By offloading scoring to LangSmith (`examples/rag/langgraph_crag_local.ipynb:580` “save it in LangSmith”), LangGraph keeps its own CI hermetic and lightweight, but trajectory experiments are not reproducible without a LangSmith API key, and `libs` tests cannot run `check_trajectory_*` without mocking `langsmith`.
- **Exact-match trajectory check vs semantic equivalence**: The demo equates path correctness with list equality (`tool_calls == expected_trajectory_1 or expected_trajectory_2`, `examples/rag/langgraph_crag_local.ipynb:731`). This is brittle to harmless reordering, extra `grade_document_retrieval` invocations, or alternative valid decompositions (e.g., `web_search` before grading is invalid, but parallel `retrieve_documents` ×2 is valid and would be penalized).
- **Retry-as-resilience vs retry-as-evaluation**: `RetryPolicy` gives cheap, automatic in-graph recovery (backoff, jitter, `retry_on` filtering, `NodeTimeoutError` as retryable by default `libs/langgraph/tests/test_retry.py:232-244`). That hides recovery events from the trajectory — a scorer would need to instrument `_AttemptEvent` or `tasks[i].error` history to know recovery happened at all.

## Failure Modes / Edge Cases

- **Trajectory extractor fragility**: `find_tool_calls_react` assumes ReAct-style `messages[*].tool_calls[*].name` (`examples/rag/langgraph_crag_local.ipynb:724-731`); `check_trajectory_custom` assumes `root_run.outputs["steps"]` (`examples/rag/langgraph_crag_local.ipynb:739-745`). A graph that emits tools via `Send()` or via a channel other than `messages`/`steps` will silently score 0. No schema validation of expected trajectory elements (string vs `Send(node=...)`).
- **Partial-trajectory scoring missing**: If the graph faults halfway (`GraphRecursionError`, `ParentCommand` bubble, `NodeTimeoutError`), `run.outputs` may be absent; both `check_trajectory_*` implementations dereference `root_run.outputs` without null guard and will throw rather than returning a partial score. The `tasks[i].error` / `pending_writes` signal (`libs/langgraph/langgraph/pregel/debug.py:118`) is ignored by the scorer.
- **Checkpoint history unbounded in eval**: `get_state_history` with no `limit`/`before` pagination (`libs/langgraph/langgraph/pregel/main.py:1501`) on a long trajectory (100-step `test_invoke_many_processes`) can materialize thousands of checkpoints; no evaluator caps memory.
- **Retry masking failure**: Default `RetryPolicy()` retries on transient `ConnectionError`, 5xx HTTP, `NodeTimeoutError` (`libs/langgraph/tests/test_retry.py:179-245`). A tool that is flaky but *eventually* succeeds will produce a clean `tasks[i].error == None` in the final checkpoint, erasing evidence that the path was initially poor — recovery-aware scoring must inspect `_AttemptEvent` history, which is opt-in via `CONFIG_KEY_TIMED_ATTEMPT_OBSERVER` and off by default.
- **StreamMode version skew**: `stream_events(version="v3")` invariants (`libs/langgraph/langgraph/pregel/main.py:378-415` `_V3_INVARIANT_KWARGS`, `StreamPart` types `libs/langgraph/langgraph/types.py:341-351`) diverge from legacy `stream_mode="debug"`; an evaluator pinned to one will break on the other. No compat shim is documented for trajectory consumers.

## Future Considerations

- Add a `langgraph.evaluation` package with `BaseTrajectoryEvaluator` + concrete `ToolCallOrderEvaluator(edit_distance_threshold)`, `ContextFaithfulnessEvaluator`, `RecoveryEvaluator` that consume `StateSnapshot`/`CheckpointPayload` directly, so `pytest` can score trajectories offline without LangSmith.
- Extend `CheckpointMetadata` to optionally carry `trajectory_score` dict (keyed by evaluator name) so `get_state_history` can return historically scored trajectories without joining LangSmith.
- Promote `_AttemptEvent`/`_AttemptContext` to public `TelemetryEvent` API and emit retry counts + `recovered_after_attempts` into `CheckpointPayload.tasks[*]`, enabling recovery-aware metrics without enabling the server-only observer.
- Replace exact-match trajectory checks with sequence-alignment scoring (Levenshtein/NW) and argument-aware checks (tool args schema validation via `tool_node.handle_tool_errors`), demonstrated as a reference evaluator in `examples/`.
- Integrate `simulation_utils.create_chat_simulator` with `langsmith.evaluation.evaluate` in a tested, importable `SimulationRunner` that returns structured `Run` objects rather than a bare `StateGraph`, closing the generate → score loop.

## Questions / Gaps

- Is there an internal/undocumented trajectory scorer in LangGraph Platform (server) that consumes `_TimedAttemptObserver` events and is simply not vendored in this OSS mirror? No evidence in `libs/langgraph` but `libs/langgraph/langgraph/pregel/_retry.py:98` comment “internal observer contract consumed by langgraph-server. Do not move to `langgraph.types`” suggests server-side scoring exists but was excluded from scope, so rubric assessment applies only to OSS.
- How is context-usage supposed to be measured for RAG agents that return `{"response": ..., "steps": ..., "documents": [...]}` when `Context` channel (`libs/langgraph/langgraph/channels/`) is itself unmanaged and no evaluator inspects `documents` vs `response` grounding?
- What is the canonical way to score a trajectory that includes human-in-the-loop interrupts (`Interrupt`, `Command(resume=...)` at `libs/langgraph/langgraph/types.py:533-934`)? No example covers interrupt recovery as an evaluable trajectory dimension, and `test_interruption.py` treats interrupts as control-flow fixtures, not scored events.
- Could we run trajectory evaluators as `stream_mode="tasks"` transformers so scores are emitted live per-step? `libs/langgraph/langgraph/stream/transformers.py` supports pluggable `StreamTransformer` factories (`libs/langgraph/langgraph/pregel/main.py:416-446`), but no evaluator transformer is provided.

---

Generated by `dimension/18.02-trajectory-evaluation` against `langgraph`.
