# Source Analysis: crewai

## Result and Outcome Publication Contract

### Source Info

| Field | Value |
|-------|-------|
| Name | crewai |
| Path | `studies/agent-harness-study/sources/crewai` |
| Language / Stack | Python (Pydantic, asyncio, concurrent.futures, OpenTelemetry) |
| Analyzed | 2026-09-01 |

## Summary

CrewAI sharply separates work returns from terminal outcomes at the Task layer but collapses them at the Crew/Flow layer. An agent method returns a raw `str | BaseModel`; `Task._execute_core` (`lib/crewai/src/crewai/task.py:809`) wraps it into an immutable-record-inside-mutable-container (`TaskOutput` with frozen `ToolFailureRecord` items) and publishes via `Task.output` assignment plus `TaskCompletedEvent`/`TaskFailedEvent`. `Crew._create_crew_output` (`lib/crewai/src/crewai/crew.py:1919`) selects the last valid `TaskOutput` as the crew's value envelope (`CrewOutput`), appends `tasks_output`, drains memory writes, flushes the event bus, then emits `CrewKickoffCompletedEvent`. Failures at crew scope use exception propagation and `CrewKickoffFailedEvent` (string error only). There is no first-class Future/Promise/JoinHandle returned to callers for concurrent, repeated, or late waiter coalescing — `kickoff()` is synchronous return, `kickoff_async` is `asyncio.to_thread` wrapper, `Task.execute_async` uses a daemon-thread `Future[TaskOutput]` that is consumed internally, and streaming exposes `.result` only after queue exhaustion. Tool-level partial success is explicitly modeled via `tool_failures` alongside `raw`, but the aggregate outcome has no discriminated `Result` type, mutable Pydantic outputs are aliased not copied/frozen, and waiter cancellation is not distinguished from run cancellation.

## Rating

**5/10**

Rationale: Task-level publication is well-structured (distinct work return → `TaskOutput`, typed failure diagnostics, guarded ordering of assignment before event, ContextVar-isolated failure collection). Crew/flow layers achieve stable synchronous publication ordering (drain/flush before completion event, return after event). Deductions for: (a) no terminal outcome envelope — success is `CrewOutput`, failure is raised exception + string event, with no unified `Outcome` type; (b) mutable, aliased outputs with no freeze/copy on publication; (c) no concurrent/repeated/late-waiter contract — no shared handle, no stability guarantees under parallel `kickoff_for_each` or streaming early access; (d) waiter cancellation conflated with run cancellation or ignored (background threads continue after stream `close`); (e) diagnostic preservation is lossy at crew/flow scope (string only, stack dropped).

## Evidence Collected

Every entry MUST include a file path with line numbers. Format: `path/to/file.ts:NN`.

| Area | Evidence | File:Line |
|------|----------|-----------|
| Work-function return type | Agent work returns `str | BaseModel`; `agent.execute_task(task, context, tools)` called inside `with tool_failure_collector()` and result unpacked into `TaskOutput` fields | `lib/crewai/src/crewai/task.py:853-857`, `lib/crewai/src/crewai/task.py:861-890` |
| Work-function async variant | `agent.aexecute_task(task, context, tools)` in `_aexecute_core`, same wrapping into `TaskOutput` | `lib/crewai/src/crewai/task.py:693-731` |
| Terminal outcome type (crew) | `CrewOutput(BaseModel)` with `raw`, `pydantic`, `json_dict`, `tasks_output: list[TaskOutput]`, `token_usage: UsageMetrics`; `tool_failures` computed property aggregates tasks | `lib/crewai/src/crewai/crews/crew_output.py:14-48` |
| Terminal outcome type (task) | `TaskOutput(BaseModel)` with `raw: str`, `pydantic`, `json_dict`, `output_format`, `messages`, `tool_failures: list[ToolFailureRecord]` (frozen records) | `lib/crewai/src/crewai/tasks/task_output.py:15-57` |
| Terminal outcome type (lite agent) | `LiteAgentOutput(BaseModel)` with `raw`, `pydantic`, `tool_failures`, `messages`, `todos`, `usage_metrics` dict | `lib/crewai/src/crewai/lite_agent_output.py:32-60` |
| Terminal outcome type (flow) | Flow has no dedicated outcome model; `kickoff`/`kickoff_async` return `Any` (method outputs list), `_method_outputs: list[Any]` private attr | `lib/crewai/src/crewai/flow/runtime/__init__.py:752-787` |
| Wait / Join / Future API | `Task.execute_async(...) -> Future[TaskOutput]` creates daemon `threading.Thread` with `contextvars.copy_context()` and `future.set_result/set_exception` | `lib/crewai/src/crewai/task.py:609-639` |
| Internal future consumption | `Crew._execute_tasks` collects `futures: list[tuple[Task, Future[TaskOutput], int]]`, then `future.result()` in `_process_async_tasks`; async path uses `asyncio.Task[TaskOutput]` and `await async_task` | `lib/crewai/src/crewai/crew.py:1579-1627`, `lib/crewai/src/crewai/crew.py:2005-2018`, `lib/crewai/src/crewai/crew.py:1362-1414` |
| Async kickoff handle | `Crew.kickoff_async` is `await asyncio.to_thread(self.kickoff, ...)` — no native future handle, caller `await`s the thread | `lib/crewai/src/crewai/crew.py:1179` |
| Streaming wait handle | `CrewStreamingOutput(StreamingOutputBase[CrewOutput])` — iteration yields `StreamChunk`, final `result: CrewOutput` property raises `RuntimeError` until `_completed`, re-raises stored `_error`; `_set_result/_set_results` called from `_finalize_streaming` after queue drained | `lib/crewai/src/crewai/types/streaming.py:356-375`, `lib/crewai/src/crewai/types/streaming.py:497-577`, `lib/crewai/src/crewai/utilities/streaming.py:384-398` |
| Outcome construction (crew success) | `_create_crew_output` selects `valid_outputs[-1]`, builds `UsageMetrics`, constructs `CrewOutput`, dispatches `OUTPUT`/`EXECUTION_END` hooks, mutates `final_task_output.raw = crew_output.raw`, drains writes, `flush()`, emits `CrewKickoffCompletedEvent`, returns `crew_output` | `lib/crewai/src/crewai/crew.py:1919-1977` |
| Outcome construction (crew failure) | `kickoff()` try/except: `_dispatch_execution_end_failure(e)`, emit `CrewKickoffFailedEvent(error=str(e), ...)`, `raise`; identical in `akickoff` | `lib/crewai/src/crewai/crew.py:1068-1078`, `lib/crewai/src/crewai/crew.py:1285-1295` |
| Outcome publication ordering | `_drain_memory_writes()` + `crewai_event_bus.flush()` before `CrewKickoffCompletedEvent`; task path: `self.output = task_output` then `end_time = now()` then `TaskCompletedEvent` emit; failure: `end_time` then `TaskFailedEvent` then `raise` | `lib/crewai/src/crewai/crew.py:1962-1972`, `lib/crewai/src/crewai/task.py:922-956`, `lib/crewai/src/crewai/task.py:957-963` |
| TaskOutput assignment before event | `self.output = task_output` at `task.py:763` (async) and `task.py:922` (sync) precedes `crewai_event_bus.emit(TaskCompletedEvent)` at `task.py:793` / `task.py:952` | `lib/crewai/src/crewai/task.py:763-796`, `lib/crewai/src/crewai/task.py:922-956` |
| Result+error combination (tool) | `ToolFailurePolicy` (IGNORE/WARN/RAISE), `ToolFailureRecord` frozen, `tool_failure_collector` ContextVar, `handle_tool_failure` records and optionally raises `ToolExecutionFailedError` | `lib/crewai/src/crewai/tools/tool_failure.py:57-69`, `lib/crewai/src/crewai/tools/tool_failure.py:118-152`, `lib/crewai/src/crewai/tools/tool_failure.py:266-302`, `lib/crewai/src/crewai/tools/tool_failure.py:324-384` |
| Dual publication of value+failure | `TaskOutput` constructed with `raw` plus `tool_failures=list(execution_failures)`; `CrewOutput.tool_failures` aggregates all tasks — crew can succeed with non-empty failures | `lib/crewai/src/crewai/task.py:879-890`, `lib/crewai/src/crewai/crews/crew_output.py:36-48`, `lib/crewai/src/crewai/task.py:730-731` |
| Crew success with tool failures assertion | Docstring: "A crew can finish successfully with a non-empty list -- agents narrate a failed step and carry on. Check it before treating raw as complete." | `lib/crewai/src/crewai/crews/crew_output.py:37-43` |
| Guardrail retry outcome mutation | `_invoke_guardrail_function` accumulates `tool_failures` across retries, re-exports output, re-assigns `task_output` in a loop up to `guardrail_max_retries+1` | `lib/crewai/src/crewai/task.py:1327-1380` (and sync variant at `task.py:1449-1530` implied) |
| Error cause / stack preservation (task) | `TaskFailedEvent.error: str`, `error_type: type[BaseException] | None` serialized as `module:qualname` via `field_serializer`, deserialized via `vars(sys.modules[module])` lookup | `lib/crewai/src/crewai/events/types/task_events.py:137-179` |
| Error cause preservation (crew) | `CrewKickoffFailedEvent.error: str` only — type and traceback dropped | `lib/crewai/src/crewai/events/types/crew_events.py:51-56` |
| Copying / freezing / ownership | `CrewOutput` and `TaskOutput` are plain `BaseModel` (mutable), `tasks_output` holds same object refs as `Task.output`, `final_task_output.raw` mutated via hook payload (alias visible to holder of `tasks_output[-1]`) | `lib/crewai/src/crewai/crews/crew_output.py:14-32`, `lib/crewai/src/crewai/tasks/task_output.py:15-45`, `lib/crewai/src/crewai/crew.py:1956-1958` |
| Frozen diagnostic record | `ToolFailureRecord` and `ToolFailure` declared `model_config = ConfigDict(frozen=True)` | `lib/crewai/src/crewai/tools/tool_failure.py:79`, `lib/crewai/src/crewai/tools/tool_failure.py:118-119` |
| Late/ concurrent waiter contract | No shared handle returned by `kickoff`; `kickoff_for_each` copies crew per input and returns `list[CrewOutput]` or single `CrewStreamingOutput` with fan-out; streaming `result` requires full iteration; `from_checkpoint` replay rebuilds outputs but not a shared waiter registry | `lib/crewai/src/crewai/crew.py:1091-1125`, `lib/crewai/src/crewai/types/streaming.py:538-577`, `lib/crewai/src/crewai/crew.py:2031-2069` |
| Waiter cancellation | `StreamingOutputBase.close()/aclose()` sets `_cancelled/_completed`, closes iterator, unregisters handler — but underlying `run_func` thread (crew execution) continues; `Task.execute_async` daemon thread has no cancellation token; no `CancelledError` mapped to terminal outcome | `lib/crewai/src/crewai/types/streaming.py:412-443`, `lib/crewai/src/crewai/utilities/streaming.py:400-410`, `lib/crewai/src/crewai/task.py:614-639` |
| Execution-scoped isolation | `begin_execution()/end_execution()` mint stack for outermost kickoff, nested inherit; `runtime_scope = _enter_runtime_scope() / _exit_runtime_scope()` brackets event scope | `lib/crewai/src/crewai/crew.py:1045-1086`, `lib/crewai/src/crewai/execution.py:1-67` (referenced via grep) |
| Tests for exclusivity/stable access | `test_tool_failure.py` asserts `has_tool_failures` on all output types, `tool_failures` not growing across kickoffs, `merge_tool_failures` dedup, ContextVar isolation | `lib/crewai/tests/tools/test_tool_failure.py:653-667`, `lib/crewai/tests/tools/test_tool_failure.py:783-809`, `lib/crewai/tests/tools/test_tool_failure.py:1028-1047`, `lib/crewai/tests/tools/test_tool_failure.py:1424-1442`, `lib/crewai/tests/tools/test_tool_failure.py:1642-1655` |

## Answers to Dimension Questions

### 1. Is a work return distinct from the runtime's terminal outcome?

**Yes at Task scope, collapsed at Crew/Flow scope.** The agent work function returns a bare `str` (or `BaseModel` from structured output) via `agent.execute_task` / `agent.aexecute_task` (`lib/crewai/src/crewai/task.py:853`, `lib/crewai/src/crewai/task.py:693`). `Task._execute_core` (`lib/crewai/src/crewai/task.py:809`) transforms it: optional `_export_output` JSON/Pydantic conversion, construction of `TaskOutput(raw=..., pydantic=..., json_dict=..., tool_failures=...)` (`lib/crewai/src/crewai/task.py:879`), assignment to `self.output` (`lib/crewai/src/crewai/task.py:922`), and publication via `TaskCompletedEvent` (`lib/crewai/src/crewai/task.py:952`). This is a distinct envelope with diagnostics, format tag, and identity.

At Crew scope the distinction is weaker: `Crew._create_crew_output` (`lib/crewai/src/crewai/crew.py:1919`) selects `final_task_output.raw` (already a `TaskOutput` envelope) as `CrewOutput.raw`, adds `tasks_output`, `token_usage`, and returns `CrewOutput`. Callers of `crew.kickoff()` receive `CrewOutput` directly — there is no wrapper like `Outcome[CrewOutput, CrewError]`. Failure is not a value in that type but a raised exception plus a side-channel `CrewKickoffFailedEvent(error=str(e))` (`lib/crewai/src/crewai/crew.py:1070`). LiteAgent similarly returns `LiteAgentOutput` on success and raises on `ToolExecutionFailedError` (`lib/crewai/src/crewai/lite_agent.py:559`). Flow returns `Any` with no outcome envelope at all — callers pattern-match on returned value versus exception, not on a sum type.

### 2. Which result, failure, cancellation, and runtime-fault combinations are valid?

**Three valid combinations, one explicitly reserved, one conflated:**

- **Value + no failure (happy path):** `TaskOutput(raw=..., tool_failures=[])`, `CrewOutput(raw=..., tasks_output=[...], has_tool_failures==False)`. `ToolFailurePolicy.IGNORE` also forces this even when a tool reported `ToolFailure` (`lib/crewai/src/crewai/tools/tool_failure.py:342-344` returns `None`).

- **Value + non-terminal tool failures (partial success):** `TaskOutput` with `raw` plus `tool_failures: list[ToolFailureRecord]` (`lib/crewai/src/crewai/tasks/task_output.py:50`), aggregated to `CrewOutput.tool_failures` (`lib/crewai/src/crewai/crews/crew_output.py:36`). Documented as "A crew can finish successfully with a non-empty list" (`lib/crewai/src/crewai/crews/crew_output.py:37`). Valid `reason` values include `TOOL_REPORTED`, `EXCEPTION`, `MCP_ERROR`, `USAGE_LIMIT`, `UNKNOWN_TOOL`, `INVALID_INPUT` (`lib/crewai/src/crewai/tools/tool_failure.py:35-55`). Guardrail retries append to `tool_failures` across attempts (`lib/crewai/src/crewai/task.py:1347-1350`).

- **No value + terminal failure (exception):** Task emits `TaskFailedEvent(error=str(e), error_type=type(e))` (`lib/crewai/src/crewai/task.py:957`) and raises; crew emits `CrewKickoffFailedEvent(error=str(e))` (`lib/crewai/src/crewai/crew.py:1071`) and raises. `ToolFailurePolicy.RAISE` converts a tool failure into `ToolExecutionFailedError` (`lib/crewai/src/crewai/tools/tool_failure.py:381`) which aborts the task and therefore the crew.

- **Cancellation → mapped to failure (conflated):** `asyncio.CancelledError` is not given a distinct terminal outcome. In MCP transport it is caught, `failure_from_exception` wraps it, and the non-cancellation cause is prioritized (`lib/crewai/src/crewai/mcp/exceptions.py:185-264`). In streaming background tasks, `CancelledError` is suppressed as an expected teardown (`lib/crewai/src/crewai/utilities/streaming.py:326`, `lib/crewai/src/crewai/utilities/streaming.py:562`). There is no `CancelledOutcome` or cancellation fact propagated to `CrewOutput`.

- **Runtime fault → same as terminal failure:** Any exception from agent, tool, guardrail, or memory drain propagates as the crew exception; no separate `RuntimeFault` envelope exists. `end_execution`/`_dispatch_execution_end_failure` ensures `EXECUTION_END(status="failed")` hook fires exactly once even when the exception path is taken (`lib/crewai/src/crewai/crew.py:1979-2003`).

No evidence of a `value + terminal error` sum type (e.g., `(CrewOutput, error)`) — the API forces callers to choose `try/except` versus `has_tool_failures` inspection.

### 3. What happens when work returns both a value and an error?

**Coexistence is first-class at tool granularity, disallowed at crew granularity.**

At the tool boundary, a tool may return `ToolFailure` (`lib/crewai/src/crewai/tools/tool_failure.py:71`) which `detect_tool_failure` (`lib/crewai/src/crewai/tools/tool_failure.py:154`) distinguishes from a legitimate string that merely describes an error. Under `WARN` (default), `handle_tool_failure` (`lib/crewai/src/crewai/tools/tool_failure.py:324`) records a frozen `ToolFailureRecord`, emits `ToolFailureDetectedEvent` (`lib/crewai/src/crewai/tools/tool_failure.py:362`), and returns the record — the agent still sees `failure.as_agent_message()` as its text observation (`lib/crewai/src/crewai/tools/tool_failure.py:104`), so the LLM narrates and continues. The task's `execution_failures` list (ContextVar-isolated via `tool_failure_collector()` at `lib/crewai/src/crewai/task.py:852`) is attached to `TaskOutput.tool_failures` (`lib/crewai/src/crewai/task.py:889`), and the crew completes with `has_tool_failures==True` but a non-empty `raw`. Tests assert this (`lib/crewai/tests/tools/test_tool_failure.py:345-348`).

Under `RAISE`, the same `ToolFailure` aborts: `ToolExecutionFailedError(record)` is raised (`lib/crewai/src/crewai/tools/tool_failure.py:381`), caught by `Task._execute_core`'s `except` which emits `TaskFailedEvent` and re-raises (`lib/crewai/src/crewai/task.py:957`), which aborts `Crew._execute_tasks` and surfaces as `CrewKickoffFailedEvent`. Under `IGNORE`, the failure is erased (`lib/crewai/src/crewai/tools/tool_failure.py:343`).

At crew level there is no valid `CrewOutput` that simultaneously carries `error` — callers must not expect `(value, error)` tuple. Guardrail retries illustrate the ambiguity handling for tasks: each failed guardrail attempt accumulates its `tool_failures` but the final `TaskOutput` still carries the last raw value (`lib/crewai/src/crewai/task.py:1392-1430`).

### 4. Can a terminal outcome expose a partial or mutable value?

**Yes — partial is by design, mutability is unguarded.**

- **Partial value:** `TaskOutput.raw` may be truncated LLM output or narrated repair text after tool failures; `json_dict`/`pydantic` may be `None` when JSON conversion failed (`lib/crewai/src/crewai/task.py:1244-1247` falls back to `(None, None)` on `JSONDecodeError`). `CrewOutput` exposes `tasks_output` which includes every task's partial output, not just the final. Callers are expected to check `has_tool_failures` / `tool_failures` before trusting `raw` (`lib/crewai/src/crewai/crews/crew_output.py:37-42`).

- **Mutable value:** Both `CrewOutput` (`lib/crewai/src/crewai/crews/crew_output.py:14`) and `TaskOutput` (`lib/crewai/src/crewai/tasks/task_output.py:15`) are plain `BaseModel` instances with no `frozen=True`. The diagnostic records inside are frozen, but the containers are not. `Crew._create_crew_output` mutates the published object after hook dispatch: `final_task_output.raw = crew_output.raw` (`lib/crewai/src/crewai/crew.py:1957`) — any holder of `tasks_output[-1]` sees the mutation. `TaskOutput.messages: list[LLMMessage]` (`lib/crewai/src/crewai/tasks/task_output.py:47`) and `Task.messages=agent.last_messages` (`lib/crewai/src/crewai/task.py:888`) are aliased lists, not copies. `Task.output: TaskOutput | None` (`lib/crewai/src/crewai/task.py:213`) is a public mutable field that `replay()` overwrites (`lib/crewai/src/crewai/crew.py:2066`). No `copy()`, `deepcopy()`, or `freeze()` occurs on publication; evidence of copying is only for crew replay isolation (`crew.py:2114 copy()` clones agents/tasks but not outputs for late consumers).

### 5. Do concurrent, repeated, and late waiters observe the same committed facts?

**No — there is no shared waiter contract; repeated and late observation is eventually consistent at best.**

- **Concurrent waiters:** `Crew.kickoff()` returns `CrewOutput` directly — no `Future[CrewOutput]` or `JoinHandle` is handed out for concurrent callers to `await` the same run. `kickoff_for_each` (`lib/crewai/src/crewai/crew.py:1091`) clones the crew (`crew.copy()`) per input and runs them independently, aggregating `usage_metrics`; there is no fork-join handle that multiple awaiters coalesce on. `kickoff_for_each_async` (`lib/crewai/src/crewai/crew.py:1181`) similarly fans out `asyncio.create_task(kickoff_fn(...))` (`lib/crewai/src/crewai/crews/utils.py:500`) and `await asyncio.gather(*async_tasks)`. Two threads calling `kickoff()` on the same instance are not serialized and have no documented outcome-sharing — the `process` sequential/hierarchical logic assumes single-owner execution.

- **Repeated waiters (polling after completion):** After `kickoff()` returns, `Task.output` and `Crew.usage_metrics` remain on the instance (`lib/crewai/src/crewai/task.py:922`, `lib/crewai/src/crewai/crew.py:1930`). A second read of `task.output` or `crew.usage_metrics` yields the same object (no snapshot), but without immutability guarantees — a caller that mutated the first read mutates the second. No `OnceCell` or versioning protects repeated reads.

- **Late waiters (attach after completion):** Streaming is the closest to a late-join contract: `CrewStreamingOutput.result` (`lib/crewai/src/crewai/types/streaming.py:356`) raises `RuntimeError("Streaming has not completed yet...")` until `_completed`, and `StreamSession.subscribe` (`lib/crewai/src/crewai/types/streaming.py:160`) replays buffered `_frames` if `_exhausted` is true. However `StreamingOutputBase.chunks` (`lib/crewai/src/crewai/types/streaming.py:387`) is a copy but `result` is a shared reference; if a task completed between waiter arrival and `result` access, the late waiter sees the same `CrewOutput` object, not a stable snapshot. For non-streaming kickoff there is no late-waiter API — `Task.execute_async`'s `Future` (`lib/crewai/src/crewai/task.py:616`) is private to `_execute_tasks` and never returned to external late joiners. Tests for concurrent/await ordering exist only for tool-failure isolation (`lib/crewai/tests/tools/test_tool_failure.py:1431-1442` checks `active_tool_failures()` ContextVar nesting, not waiter visibility).

In short, CrewAI relies on synchronous return + instance field memory for stability, not on a publish-once, read-many handle. `crewai_event_bus.flush()` before completion event (`lib/crewai/src/crewai/crew.py:1963`) ensures async telemetry handlers have drained before the completion event, but not before a waiter reads `CrewOutput` — there is no fence for waiters.

### 6. Can waiting be cancelled without cancelling the run itself?

**No. Waiting and execution are coupled; cancellation primitives are either absent or abort the run.**

- **Synchronous kickoff:** No `CancellationToken` or `abort_signal` parameter exists on `Crew.kickoff(inputs=...)` (`lib/crewai/src/crewai/crew.py:992`). A caller that `Ctrl-C` / raises `CancelledError` in a wrapper will unwind `kickoff()`'s `try/finally` which already emitted `CrewKickoffFailedEvent` and detached context (`lib/crewai/src/crewai/crew.py:1080-1086`) — the partially completed crew stays in its mutated `Task.output` state with no rollback.

- **Async kickoff:** `kickoff_async` (`lib/crewai/src/crewai/crew.py:1179`) is `asyncio.to_thread`; cancelling the awaiting coroutine cancels the await but does not cancel the worker thread — `threading.Thread(daemon=True)` (`lib/crewai/src/crewai/task.py:618`) continues executing the agent loop. No `future.cancel()` path maps waiter cancellation to a no-op while preserving the run.

- **Streaming waiter cancellation:** `StreamingOutputBase.close()` / `aclose()` (`lib/crewai/src/crewai/types/streaming.py:412`, `lib/crewai/src/crewai/types/streaming.py:428`) sets `_cancelled/_completed` and closes the iterator, and `create_chunk_generator`'s `finally` calls `_finalize_streaming` and `thread.join()` (`lib/crewai/src/crewai/utilities/streaming.py:514`), but `run_func` is already running in a daemon thread — `close()` does not interrupt it. The crew continues to completion in the background, potentially emitting `CrewKickoffCompletedEvent` after the waiter has given up. The `async` generator variant does `task.cancel()` for the async runner (`lib/crewai/src/crewai/utilities/streaming.py:558`) but suppresses `CancelledError` (`lib/crewai/src/crewai/utilities/streaming.py:562`), again conflating waiter cancel with run teardown.

- **Task-level async execution:** `Task.execute_async` (`lib/crewai/src/crewai/task.py:614`) starts a `Future` on a daemon thread; `Crew._process_async_tasks` calls `future.result()` blocking (`lib/crewai/src/crewai/crew.py:2012`) with no timeout or waiter-side cancel. No test shows `future.cancel()` leaving the crew running.

### 7. Who owns values and diagnostics after publication?

**The caller receives a shared mutable reference; aliasing is preserved, ownership is ambiguous, and diagnostics are frozen inside a mutable shell.**

- **Task value/diagnostic ownership:** `Task._execute_core` constructs `TaskOutput` with `tool_failures=list(execution_failures)` (`lib/crewai/src/crewai/task.py:889`), assigns `self.output = task_output` (`lib/crewai/src/crewai/task.py:922`), and returns the same object. The caller of `task.execute_sync` owns that reference, but so does `task.output` and any `CrewOutput.tasks_output` that later appends it (`lib/crewai/src/crewai/crew.py:1620`). Mutation by one owner is visible to others. `ToolFailureRecord` is `frozen=True` (`lib/crewai/src/crewai/tools/tool_failure.py:118`), so diagnostics cannot be mutated through the record, but the `tool_failures` list container can be replaced or extended.

- **Crew value/diagnostic ownership:** `Crew._create_crew_output` builds a new `CrewOutput` (`lib/crewai/src/crewai/crew.py:1935`) but reuses the same `TaskOutput` objects from `task_outputs` (`lib/crewai/src/crewai/crew.py:1939`). The returned `CrewOutput` and the per-task `Task.output` alias. The hook `OUTPUT` dispatch receives `payload=crew_output` (`lib/crewai/src/crewai/crew.py:1943`) which can mutate `crew_output.raw` in place; the mutation is then reflected into `final_task_output.raw` (`lib/crewai/src/crewai/crew.py:1957`), demonstrating that ownership is shared mutation rights, not transfer.

- **Flow/lite ownership:** `Flow._method_outputs: list[Any]` (`lib/crewai/src/crewai/flow/runtime/__init__.py:752`) and `Flow._create_initial_state` return raw state refs; no copy boundary exists between method return and flow state. `LiteAgentOutput.tool_failures` (`lib/crewai/src/crewai/lite_agent_output.py:54`) follows the same frozen-record-in-mutable-list pattern.

- **Event ownership:** `TaskFailedEvent.error` and `CrewKickoffFailedEvent.error` store `str(e)` (`lib/crewai/src/crewai/task.py:960`, `lib/crewai/src/crewai/crew.py:1073`), copying the message but dropping the exception object. `TaskFailedEvent.error_type` preserves the class via `module:qualname` serialization (`lib/crewai/src/crewai/events/types/task_events.py:161`), owned as a string, resolved lazily via `sys.modules` lookup (`lib/crewai/src/crewai/events/types/task_events.py:54-72`) — ownership of the class identity, not the stack.

No evidence of `copy.deepcopy`, `freeze`, or read-only view wrappers on publication; the `copy` keyword in the codebase refers to crew cloning for `kickoff_for_each`, not to output publication.

### 8. Does waiter release occur only after the complete outcome is visible?

**For synchronous crew kickoff, yes — with an explicit drain/flush fence. For task-level and streaming paths, conditionally.**

- **Crew synchronous path (correct):** `_create_crew_output` (`lib/crewai/src/crewai/crew.py:1919`) drains all pending memory writes via `_drain_memory_writes()` (`lib/crewai/src/crewai/crew.py:1887-1917` — iterates crew, manager, and per-agent backs, `id`-deduped, calls `drain_writes()`), then `crewai_event_bus.flush()` (`lib/crewai/src/crewai/crew.py:1963`) to pump async handlers, then emits `CrewKickoffCompletedEvent` (`lib/crewai/src/crewai/crew.py:1964`), then returns `crew_output` (`lib/crewai/src/crewai/crew.py:1977`). The comment at `lib/crewai/src/crewai/crew.py:1080` — "Safety net ... the success path already drained in _create_crew_output before emitting completion" — confirms intent. The `finally` block (`lib/crewai/src/crewai/crew.py:1080-1086`) drains again on the exception path before `CrewKickoffFailedEvent`. This satisfies "complete outcome visible before release."

- **Task path (correct):** `self.output = task_output` and `self.end_time = now()` (`lib/crewai/src/crewai/task.py:922-923`) occur before `TaskCompletedEvent` (`lib/crewai/src/crewai/task.py:952`). On failure, `end_time` set before `TaskFailedEvent` (`lib/crewai/src/crewai/task.py:958`). A waiter that observes the event after it is emitted will see `task.output` already populated.

- **Async task join (conditionally correct):** `_aexecute_tasks` (`lib/crewai/src/crewai/crew.py:1341`) buffers async tasks and only joins via `_aprocess_async_tasks` (`lib/crewai/src/crewai/crew.py:1435`) when a sync task barrier or end-of-list is reached. `_process_async_tasks` (`lib/crewai/src/crewai/crew.py:2005`) blocks on `future.result()` per future before appending to `task_outputs`. `_create_crew_output` is called only after all futures joined (`lib/crewai/src/crewai/crew.py:1627`). Thus release follows visibility, but if a waiter held a direct `Future[TaskOutput]` from `Task.execute_async` it could see the task output before `TaskCompletedEvent` — the event is emitted inside `Task._execute_task_async` after `set_result` (`lib/crewai/src/crewai/task.py:635` + `task.py:952`), so `future.done()` may become true a moment before the event.

- **Streaming path (delayed visibility):** `StreamingContext.run_crew` appends to `ctx.result_holder` inside the runner thread, then `signal_end(ctx.state)` (`lib/crewai/src/crewai/crew.py:1021-1029`); `_finalize_streaming` sets `_set_result(result_holder[0])` only after the queue's `None` sentinel is consumed and `thread.join()` completes (`lib/crewai/src/crewai/utilities/streaming.py:506-519`). A waiter that iterates until exhaustion sees `result` only after the background thread drained. However a waiter that calls `close()` early gets `_cancelled` without ever seeing `result` — the background thread still runs to completion but the waiter is released with no value, violating "release only after complete outcome visible" for the cancelling waiter (it is released early by design, but no outcome is withheld — it simply never materializes for that waiter).

- **Flow path (parallel to crew):** `Flow.kickoff_async` and `Flow.resume_async` mirror crew's pattern: `flush()` + `_detach_usage_aggregation_listener` in `finally` (`lib/crewai/src/crewai/flow/runtime/__init__.py:1405-1413`), but the search did not locate an equivalent `drain_memory_writes` before `FlowFinishedEvent` — Flow's drain is per-flow `memory` only (`flow/runtime/__init__.py:1014-1028`), not per-method.

## Architectural Decisions

- **ContextVar-isolated failure collection:** `tool_failure_collector()` (`lib/crewai/src/crewai/tools/tool_failure.py:266`) uses `ContextVar[list[ToolFailureRecord]]` so concurrent tasks sharing an agent report only their own failures. This enables concurrent `async_execution` tasks on the same agent without cross-contamination — a deliberate isolation mechanism for the absence of per-execution handles. `active_tool_failures()` (`lib/crewai/src/crewai/tools/tool_failure.py:283`) allows reading inside guardrail retries.

- **Frozen diagnostic records inside mutable envelopes:** `ToolFailure` / `ToolFailureRecord` (`lib/crewai/src/crewai/tools/tool_failure.py:79`, `118`) are `frozen=True` while `TaskOutput`/`CrewOutput` are mutable. This balances append-safety for diagnostics with Pydantic ergonomics for outputs.

- **Synchronous publication primitive:** Crew chose synchronous return + instance-field side effects over a Future/Promise. `_create_crew_output` is the single publication site that both emits the completion event and returns the value, with `drain + flush` fences inlined. Tradeoff: simple single-waiter case is correct, but multi-waiter fan-out must be built externally.

- **Daemon threads for Task async:** `Task.execute_async` (`lib/crewai/src/crewai/task.py:618`) uses `threading.Thread(daemon=True, target=ctx.run, ...)` — no thread pool, no structured concurrency, no ability to cancel. This keeps the implementation simple at the cost of waiter cancellation semantics.

- **Streaming as the only concurrent-waiter handle:** `CrewStreamingOutput` / `StreamSession` (`lib/crewai/src/crewai/types/streaming.py:332`) provides the only shareable handle with `is_completed`, `is_cancelled`, `chunks`, `result`/`results`. Its isolation uses `queue.Queue` + `ContextVar[_current_stream_ids]` (`lib/crewai/src/crewai/utilities/streaming.py:42`, `435`) for per-stream chunk routing.

- **String-only crew error envelope:** `CrewKickoffFailedEvent.error: str` (`lib/crewai/src/crewai/events/types/crew_events.py:54`) versus `TaskFailedEvent.error_type: type[BaseException]` (`lib/crewai/src/crewai/events/types/task_events.py:141`) — crew-level errors lose type fidelity intentionally to avoid leaking prompts/credentials, following `Telemetry._safe_error_type` rationale in the task event docstring (`task_events.py:142-148`).

## Notable Patterns

- **Tool failure as data, not control flow:** Tools return `ToolFailure` values; policy decides whether to record, warn, or raise. This is visible in every executor loop via `tool_failure_collector` and `handle_tool_failure`.

- **Allocation-free task output shortcut:** If the work already returned a `BaseModel`, `_execute_core` bypasses `convert_to_model` and directly assigns `model_dump_json`/`model_dump` (`lib/crewai/src/crewai/task.py:861-871`), preserving caller-allocated objects without re-parsing.

- **Guardrail retry as output transformer:** `._invoke_guardrail_function` (`lib/crewai/src/crewai/task.py:1327`) loops `guardrail_max_retries + 1` times, accumulating failures and re-exporting output on each retry — the final `TaskOutput` is the last transformed value, not the first.

- **Event hook payload mutation:** `OUTPUT` and `EXECUTION_END` hooks receive `payload=crew_output` (`lib/crewai/src/crewai/crew.py:1943`, `1947`) and the mutated payload is cast back to `CrewOutput` — hooks can rewrite the terminal outcome in place before it is returned or emitted.

- **Process barrier for async tasks:** Any sync task acts as a barrier that drains pending async futures before executing (`lib/crewai/src/crewai/crew.py:1608`, `lib/crewai/src/crewai/crew.py:1393`). This enforces that `task_outputs` ordering matches `tasks` ordering regardless of async completion times.

## Tradeoffs

- **Mutable aliasing vs copy-on-publish:** Aliasing `TaskOutput` into both `Task.output` and `CrewOutput.tasks_output` avoids copying large `messages` lists but means a consumer that mutates `crew_output.tasks_output[0].raw` mutates the task itself. No defensive copy trades correctness for memory/speed.

- **Single Future per task vs shared handle:** Exposing `Future[TaskOutput]` per task internally but not per crew simplifies the API (callers just `await kickoff_async`) but prevents external orchestration like `asyncio.wait([handle], timeout=5)` or late-join coalescing. The streaming handle is the only alternative, but it couples chunk iteration with result fetching.

- **Thread-per-async-task vs pool:** Daemon threads are cheap to reason about and inherit `contextvars`, but they cannot be cancelled, timed out, or bounded — a large `kickoff_for_each` fan-out (`lib/crewai/src/crewai/crews/utils.py:500` creates `asyncio.create_task` per copy, each spawning threads) risks thread explosion under `Process.sequential` with many async tasks.

- **Drain-before-emit vs latency:** Draining memory writes and flushing the bus before emitting completion (`crew.py:1962-1964`) guarantees visibility but adds latency equal to the slowest memory save and the longest event handler. No parallel drain or timeout exists — a hung `drain_writes()` blocks completion indefinitely.

- **String error at crew vs typed at task:** Crew's lossy error aids privacy and JSON checkpointing but hampers programmatic error handling (callers cannot `except ValueError` from `CrewOutput`). Task's typed `error_type` with `module:qualname` serialization (`task_events.py:161`) recovers type across checkpoints but only if the module is already imported — otherwise it silently degrades to `None` (`task_events.py:61-72`).

## Failure Modes / Edge Cases

- **Double-mutation via hook + task alias:** `_create_crew_output` mutates `final_task_output.raw` from `crew_output.raw` after hooks (`crew.py:1957`). If a hook mutated `CrewOutput.raw` to a non-string (bug), `final_task_output.raw` becomes poisoned for any later `replay` that reads `tasks[i].output.raw`.

- **Event bus handler exception swallowed:** `prepare_kickoff`'s `future.result()` on `CrewKickoffStartedEvent` suppresses handler exceptions (`lib/crewai/src/crewai/crews/utils.py:334-337` `except Exception: pass`). A failing handler that was supposed to validate cancellation never aborts the run.

- **Future.result deadlock on async barrier:** `_process_async_tasks` calls `future.result()` without timeout (`crew.py:2012`) on a daemon thread that may deadlock if the agent loop waits on a resource held by the caller. The caller cannot cancel the waiter without also failing the entire `kickoff()`.

- **Streaming waiter sees exception instead of outcome:** `signal_error(ctx.state, exc)` (`crew.py:1026`) puts the exception into the queue; `create_chunk_generator` raises it during iteration (`utilities/streaming.py:510-512`). A waiter that caught the exception and then accessed `streaming.result` (`types/streaming.py:356`) would re-raise the stored `_error`, not a `CrewOutput` with diagnostics.

- **ContextVar inheritance across `kickoff_for_each_async` fan-out:** `asyncio.create_task` inherits `contextvars` from the creating task; if `tool_failure_collector` were active during fan-out, failures could leak across crews. The implementation avoids this by collecting per-task, but the safeguard is implicit.

- **Checkpoint replay degrades error type:** If `TaskFailedEvent.error_type` cannot be resolved (`sys.modules` miss), it becomes `None` (`task_events.py:61`), so a resumed-from-checkpoint late waiter observing `TaskFailedEvent` loses the exception classifier.

- **No exclusivity guard for concurrent kickoffs on same instance:** Two concurrent `await crew.kickoff_async(...)` on the same `Crew` instance race on `self._inputs`, `self._task_output_handler`, and `self._memory` (`crew.py:989-1086`). No lock or `in_progress` guard produces exclusivity — second caller may read half-initialized `tasks_output`.

## Future Considerations

- Introduce a unified `Outcome[T, E]` discriminated type for `Crew.kickoff` / `Flow.kickoff` (e.g., `CrewResult = SuccessCrewOutput | FailureCrewOutput`) so callers can pattern-match value vs error without `try/except` and so concurrent waiters receive the same tagged value.

- Return a shareable `CrewHandle` / `TaskHandle` with `await handle.join()`, `handle.result(timeout=...)`, `handle.cancel(waiter_only=True)`, and `handle.add_done_callback` instead of raw `CrewOutput`. Make `Task.execute_async` and `Crew.kickoff_async` produce the same handle type, and ensure `handle.result()` is stable across repeated/late calls (snapshot or frozen copy).

- Freeze or snapshot `TaskOutput`/`CrewOutput` on publication (`model_config frozen=True` or `model_copy(deep=True)` in `_create_crew_output`) and return defensive copies for `tasks_output` to prevent aliasing mutation. Document ownership: "caller owns the returned copy; `task.output` is a separate snapshot."

- Distinguish waiter cancellation from run cancellation: propagate `asyncio.CancelledError` as `CancelledOutcome` with `reason="waiter_cancelled"` while leaving the run's `CrewKickoffCompletedEvent` to fire later for the execution owner. Implement `StreamingOutputBase.close(waiter_only=True)` that unregisters the waiter but does not suppress the background `signal_end`.

- Preserve full diagnostic chain at crew scope: change `CrewKickoffFailedEvent` to include `error_type: _ExceptionClass` (like `TaskFailedEvent`) plus optional `traceback` string, while still redacting `error` message per telemetry policy.

- Add waiter-release fences for task-path Futures: publish `TaskOutput` via a single `set_result` that happens atomically with `TaskCompletedEvent` emission (e.g., set a `Future` only after event `future.result()` completes), or replace daemon threads with a structured `ThreadPoolExecutor` where `future.cancel()` is meaningful.

- Add tests for: concurrent `kickoff_async` on same instance, repeated `handle.result()` identity/stability, late-waiter attachment after `CrewKickoffCompletedEvent`, waiter-only cancellation leaving `Task.output` intact, and hook mutation not leaking into `tasks_output` aliases.

## Questions / Gaps

- No evidence found for a `cancellation` field on `CrewOutput`/`TaskOutput` or a `CrewCancelledEvent` — cancellation, if it occurs via `CancelledError`, surfaces as a generic failure string, so the "cancellation facts" part of the rubric cannot be validated.

- No evidence found for value ownership transfer semantics (e.g., "run owns `raw` until `join` returns, then caller owns") — the code shares references with no transfer comment or `move` semantics; search of `crew.py`, `task.py`, `crew_output.py`, `task_output.py`, `streaming.py` found no ownership documentation.

- No evidence found for waiter-release ordering guarantee beyond the crew synchronous path's `drain + flush` — the async task path and `kickoff_async` wrapper were inspected but their ordering relies on `asyncio.to_thread` completion, not on an explicit fence before waiter notification.

- Flow outcome publication (`Flow.kickoff` return type `Any`, `Flow._method_outputs`, `FlowFinishedEvent` payload) was only shallowly inspected; full audit of `flow/runtime/__init__.py:2400-2600` emission-vs-return ordering remains open.

---

Generated by `studies/agent-harness-study/dimensions/01.16-result-and-outcome-publication-contract.md` against `crewai`.
