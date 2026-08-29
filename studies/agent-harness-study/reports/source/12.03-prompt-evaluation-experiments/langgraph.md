# Source Analysis: langgraph

## Dimension 12.03: Prompt Evaluation and Experiments

### Source Info

| Field | Value |
|-------|-------|
| Name | langgraph |
| Path | `studies/agent-harness-study/sources/langgraph` |
| Language / Stack | Python (monorepo: `libs/langgraph`, `libs/prebuilt`, `libs/checkpoint*`, `libs/sdk-py`, `libs/cli`) |
| Analyzed | 2026-08-29 |

## Summary

LangGraph is a stateful graph orchestration framework, not a prompt-evaluation harness. The only prompt logic in-repo is the `prompt` parameter of `create_react_agent` (`libs/prebuilt/langgraph/prebuilt/chat_agent_executor.py:137-170`, `libs/prebuilt/langgraph/prebuilt/chat_agent_executor.py:292`), which is a thin adapter from `str | SystemMessage | Callable | Runnable` to a `Runnable` that prepends messages. Testing of that adapter is limited to 5 unit tests that assert string concatenation with `FakeToolCallingModel` (`libs/prebuilt/tests/test_react_agent.py:148-249`); there is no evaluation dataset, no expected-output corpus, no LLM-as-judge implementation, no experiment/A-B tracker, and no CI gate that would catch a prompt regression. All evaluation capabilities are explicitly externalized to LangSmith (`README.md:56`, `docs/redirects.json:121`, `examples/rag/langgraph_crag_local.ipynb:582-830`), which is referenced but not implemented or vendored in this source. Confidence for deploying a prompt change is therefore absent — a prompt can be changed and merged with only ad-hoc functional tests passing.

## Rating

**2/10 — Absent / ad-hoc**

Rationale: No first-party eval datasets, no prompt-vs-expected-output regression suite, no experiment tracking, no shipped LLM-as-judge prompts, and no CI enforcement for prompt regressions. The only prompt coverage is functional plumbing tests for `create_react_agent`. The framework delegates eval to LangSmith as SaaS, with no reproducible harness in-repo. A breaking prompt change (e.g., altering the `SystemMessage` prepend in `_get_prompt_runnable`) would not be caught as a regression.

## Evidence Collected

Every entry MUST include a file path with line numbers. Format: `path/to/file.ts:NN`.

| Area | Evidence | File:Line |
|------|----------|-----------|
| Eval datasets | No dataset directory, fixture, or JSONL corpus found in `libs/langgraph` or `libs/prebuilt`. Only example-level `client.create_dataset` / `client.create_examples` calls that create LangSmith-hosted datasets at runtime. | `studies/agent-harness-study/sources/langgraph/examples/rag/langgraph_crag_local.ipynb:582-624` |
| Eval datasets | Simulated-chat helper `create_chat_simulator` accepts `input_key`/`max_turns` but carries no bundled inputs or golden outputs; it is a graph builder, not a dataset. | `studies/agent-harness-study/sources/langgraph/examples/chatbot-simulation-evaluation/simulation_utils.py:80-124` |
| Prompt change testing | `_get_prompt_runnable` normalizes `None`/`str`/`SystemMessage`/`callable`/`Runnable` to a `RunnableCallable` named `"Prompt"`. Logic is untested beyond happy-path. | `studies/agent-harness-study/sources/langgraph/libs/prebuilt/langgraph/prebuilt/chat_agent_executor.py:119-170` |
| Prompt change testing | `create_react_agent` docs define `prompt` contract (`str`→`SystemMessage` prepend, callable→`LanguageModelInput`, etc.) but no schema or prompt-version pinning. | `studies/agent-harness-study/sources/langgraph/libs/prebuilt/langgraph/prebuilt/chat_agent_executor.py:366-372` |
| Prompt unit tests | 5 tests cover system/str/callable/async-callable/runnable prompt: each asserts `AIMessage(content="Foo-hi?" ...)` via `FakeToolCallingModel` string concatenation. No expected-output file, no semantic assertion. | `studies/agent-harness-study/sources/langgraph/libs/prebuilt/tests/test_react_agent.py:148-203` |
| Prompt with store | `test_prompt_with_store` / `test_prompt_with_store_async` verify callable prompt can read `InMemoryStore` to inject `SystemMessage`; asserts `content == "User name is Alice-hi"` only. | `studies/agent-harness-study/sources/langgraph/libs/prebuilt/tests/test_react_agent.py:207-292` |
| Fake model for tests | `FakeToolCallingModel._generate` ignores prompt semantics; returns `"-".join(m.content for m in messages)` as `AIMessage.content` — signals harness is graph-focused, not prompt-quality focused. | `studies/agent-harness-study/sources/langgraph/libs/prebuilt/tests/model.py:28-46` |
| Experiment tracker | No experiment tracker, no `mlflow`/`wandb`/`langsmith` experiment ID, no A/B flag, no `experiment_prefix` in library code. Only notebook-level `experiment_prefix = f"custom-agent-{model_tested}"` + `evaluate(..., experiment_prefix=...)` via external `langsmith.evaluation`. | `studies/agent-harness-study/sources/langgraph/examples/rag/langgraph_crag_local.ipynb:806-813` |
| Benchmark system | `bench/__main__.py` + `bench.yml` measure throughput/latency with `pyperf.Runner`, not prompt quality. Baselines are perf JSONs (`out/benchmark-baseline.json`), not prompt scores. | `studies/agent-harness-study/sources/langgraph/libs/langgraph/bench/__main__.py:99-520` , `studies/agent-harness-study/sources/langgraph/.github/workflows/bench.yml:30-56` |
| LLM-as-judge | No judge prompt, no `evaluate`/`evaluator` export in `libs/langgraph` (`rg evaluate` returns zero hits). Closest is example `retrieval_grader = prompt | llm | JsonOutputParser()` and `answer_evaluator(run, example)` used via LangSmith — notebook-only, not library. | `studies/agent-harness-study/sources/langgraph/examples/rag/langgraph_adaptive_rag_local.ipynb:253-259` , `studies/agent-harness-study/sources/langgraph/examples/rag/langgraph_crag_local.ipynb:653-667` |
| LLM-as-judge (retrieval grader) | `grade_prompt = hub.pull("efriis/self-rag-retrieval-grader")` + `retrieval_grader = grade_prompt | structured_llm_grader` — external Hub prompt, not a versioned repo prompt, invoked ad-hoc in notebooks. | `studies/agent-harness-study/sources/langgraph/examples/rag/langgraph_self_rag_pinecone_movies.ipynb:179-186` |
| Regression test suites | `syrupy` snapshot tests only cover `agent.get_graph().draw_mermaid(...)` structure, not prompt output. `addopts = "--snapshot-warn-unused"` confirms snapshots are graph-diagram snapshots. | `studies/agent-harness-study/sources/langgraph/libs/prebuilt/tests/test_react_agent_graph.py:39-52` , `studies/agent-harness-study/sources/langgraph/libs/langgraph/pyproject.toml:133` |
| CI prompt regression gate | `ci.yml` runs `lint`, `test`, `test-langgraph`, `check-sdk-methods`, `check-schema`, `integration-test`, `sdk-py-integration-test`; no `evals`/`prompt-regression` job, no dataset download, no threshold check. `_test.yml` / `_test_langgraph.yml` only run `uv sync --group test` + `make test`. | `studies/agent-harness-study/sources/langgraph/.github/workflows/ci.yml:58-173` , `studies/agent-harness-study/sources/langgraph/.github/workflows/_test.yml:42-50` , `studies/agent-harness-study/sources/langgraph/.github/workflows/_test_langgraph.yml:40-46` |
| Documentation delegation | README positions LangSmith as external eval/observability layer: "Helpful for agent evals ... evaluate agent trajectories ..." — no in-repo eval harness claimed. | `studies/agent-harness-study/sources/langgraph/README.md:56` |
| Redirects confirm externalization | `/agents/evals`, `/cloud/how-tos/studio/run_evals`, `/tutorials/chatbot-simulation-evaluation/*` all 301-redirect to `https://docs.langchain.com/...` (LangSmith), indicating eval is owned outside this repo. | `studies/agent-harness-study/sources/langgraph/docs/redirects.json:121-278` |
| Absence of prompt versioning | No `prompt_version`, `prompt_hash`, or `prompt_registry` symbol found via `rg prompt` beyond the adapter; no prompt file to diff or pin. | `studies/agent-harness-study/sources/langgraph/libs/prebuilt/langgraph/prebuilt/chat_agent_executor.py:1-200` (search boundary) |

## Answers to Dimension Questions

**1. Are prompt changes tested?**
No. Prompt changes are tested only as plumbing: `_get_prompt_runnable` (`studies/agent-harness-study/sources/langgraph/libs/prebuilt/langgraph/prebuilt/chat_agent_executor.py:137-170`) is exercised by 5 tests that check that a `SystemMessage("Foo")` prepends to the chat history (`studies/agent-harness-study/sources/langgraph/libs/prebuilt/tests/test_react_agent.py:148-203`) and that a callable prompt can read `InMemoryStore` (`studies/agent-harness-study/sources/langgraph/libs/prebuilt/tests/test_react_agent.py:207-292`). The `FakeToolCallingModel` (`studies/agent-harness-study/sources/langgraph/libs/prebuilt/tests/model.py:36`) synthesizes `AIMessage` content as `"-".join(m.content for m in messages)`, so assertions are substring-equality checks (`assert response["messages"][-1].content == "Foo-hi?"`). There is no golden file, no dataset, no semantic or rubric-based expected-output test. Mutating the system prompt string or the prepend logic would require manually updating those hand-coded expected strings — there is no harness that would flag a semantic regression (e.g., tone change).

**2. Are experiments tracked?**
No. `rg experiment` across `libs/langgraph` and `libs/prebuilt` returns zero library hits. `bench/__main__.py:99-520` tracks latency/throughput with `pyperf`, persisted as `out/benchmark.json` and compared via `bench.yml:35-63` (`benchmark-baseline.json` cache + `COMPARE_BENCHMARKS`). This is performance benchmarking, not prompt/model experiment tracking — no experiment ID, no prompt-vs-prompt comparison table, no metric store, no A/B assignment. Example notebooks call `langsmith.evaluation.evaluate(..., experiment_prefix=...)` (`studies/agent-harness-study/sources/langgraph/examples/rag/langgraph_crag_local.ipynb:806-813`), but that is user-owned LangSmith SaaS usage; the framework itself ships no tracker and `libs/prebuilt/pyproject.toml` / `libs/langgraph/pyproject.toml` declare no tracking dependency.

**3. Is LLM-as-judge used for evaluation?**
Not in the studied source. No judge prompt, rubric, or evaluator class exists in `libs/langgraph` or `libs/prebuilt` (search for `judge|evaluate|evaluator` yields zero library files). The only LLM-as-judge artifacts are in examples: `retrieval_grader = prompt | llm | JsonOutputParser()` (`studies/agent-harness-study/sources/langgraph/examples/rag/langgraph_adaptive_rag_local.ipynb:253-259`), `grade_prompt = hub.pull("efriis/self-rag-retrieval-grader")` (`studies/agent-harness-study/sources/langgraph/examples/rag/langgraph_self_rag_pinecone_movies.ipynb:179-186`), and `def answer_evaluator(run, example)` + `check_trajectory_custom` (`studies/agent-harness-study/sources/langgraph/examples/rag/langgraph_crag_local.ipynb:653-813`). All are notebook-scoped, pull prompts from Hub at runtime, and require a live LLM + LangSmith account — none is version-controlled as a repo prompt with tests.

**4. Are regressions caught before deployment?**
No. CI (`studies/agent-harness-study/sources/langgraph/.github/workflows/ci.yml:58-173`) fans out to `lint`, `test`, `test-langgraph`, `check-sdk-methods`, `check-schema`, `integration-test`; `_test.yml:42-50` and `_test_langgraph.yml:40-46` run `make test` / `make test_parallel` only. There is no `eval` job, no dataset checkout, no score-threshold gate, and no prompt-diff check. The only automated regression guards are graph-semantics/performance: `test_pregel.py` / `test_pregel_async.py`, `test_tool_node.py:1970`, `checkpoint` conformance suites, and `bench.yml` baseline comparison (perf, not prompt). A prompt regression (e.g., dropping the system prompt in `_get_prompt_runnable`) would pass CI unless it broke the exact string assertion in `test_system_message_prompt`.

## Architectural Decisions

- **Evaluation is externalized to LangSmith SaaS.** `README.md:56` and `docs/redirects.json:121-278` position eval/observability as a LangSmith concern. This keeps the framework dependency-free but means the studied repo ships zero reproducible eval infrastructure. Evidence: `docs/redirects.json:224` maps `/agents/evals` → `docs.langchain.com` rather than local docs.
- **Prompt is an unversioned runtime adapter, not a declarative artifact.** `chat_agent_executor.py:121-126` types `Prompt` as `str | SystemMessage | Callable | Runnable` and materializes it per-call via `_get_prompt_runnable` (`chat_agent_executor.py:137-170`). No prompt file, hash, or registry exists to diff or pin — aligns with LangGraph's "bring your own prompt" philosophy.
- **Testing pyramid prioritizes graph/checkpoint correctness over LLM behavior.** `libs/langgraph/tests/` (~60 test files) focuses on `Pregel`, `DeltaChannel`, `Topic/LastValue`, `checkpointer` internals; `libs/prebuilt/tests/model.py:22-46` provides a deterministic fake model that sidesteps LLM non-determinism. Tradeoff: stable determinism, but zero signal on prompt quality.
- **Performance benchmarking substitutes for prompt benchmarking.** `bench/__main__.py:99-520` + `.github/workflows/bench.yml:12-63` + `baseline.yml:14-37` implement perf baselines with `pyperf`. This is a mature perf-regression system, but it is orthogonal to prompt evaluation.

## Notable Patterns

- **Deterministic fake-model pattern for prompt plumbing tests.** `FakeToolCallingModel._generate` (`libs/prebuilt/tests/model.py:28-46`) joins message contents deterministically, enabling the narrow `test_system_message_prompt` / `test_string_prompt` assertions. Pattern is effective for graph plumbing but inapplicable to prompt quality.
- **Syrupy graph-diagram snapshots as the only snapshot testing.** `test_react_agent_graph.py:39-52` snapshots `draw_mermaid(with_styles=False)` — a graph-structure snapshot, not a prompt-output snapshot. Confirms snapshot infra exists but is unused for prompts.
- **Notebook-driven LLM-as-judge anti-pattern.** `langgraph_self_rag_pinecone_movies.ipynb:179-186`, `langgraph_adaptive_rag_local.ipynb:253-259`, `langgraph_crag_local.ipynb:653-813` demonstrate the judge inline (`hub.pull` → `structured_llm_grader` → `evaluate`). These are educational sketches, not a reusable `evals/` harness, and they depend on network + secrets at runtime.
- **Delegation-via-redirect documentation pattern.** `docs/redirects.json:121-278` 301s all eval-related docs to LangSmith — signals intentional offloading rather than omission-by-accident.

## Tradeoffs

- **Framework-agnostic vs. prompt-guardrailed:** By treating `prompt` as an arbitrary `Callable|Runnable` (`chat_agent_executor.py:121-126`), LangGraph maximizes user flexibility (dynamic prompt per-state, store-aware prompts `test_prompt_with_store:216-237`) at the cost of no auditable prompt surface to evaluate. You cannot diff what you do not declare.
- **Determinism vs. realism:** Fake-model tests (`libs/prebuilt/tests/model.py:36`) are hermetic and fast; they avoid flaky LLM calls but also blind the suite to prompt-induced regressions (e.g., increased hallucination, instruction leakage).
- **SaaS delegation vs. repo-reproducibility:** Delegating eval to LangSmith (`README.md:56`, `docs/redirects.json:121`) avoids vendoring evalインフラ and credentials, but contributors cannot run prompt regressions locally or in CI without an external account and dataset setup.
- **Perf baselines vs. quality baselines:** The mature perf harness (`bench/__main__.py:470-520`, `bench.yml:41-63`) proves the team knows how to build baseline-comparison CI; the same mechanism is not applied to prompt quality, leaving a coverage asymmetry.

## Failure Modes / Edge Cases

- **Silent system-prompt drop.** If `_get_prompt_runnable` regresses (e.g., `isinstance(prompt, str)` branch removed or ordering changes in `chat_agent_executor.py:143-148`), only `test_string_prompt` catches it via `assert content == "Foo-hi?"`. A subtler change (e.g., prompt rendered *after* messages instead of prepended) would pass if tests are not updated, and CI has no semantic check.
- **Callable-prompt signature drift.** Async vs sync callable dispatch (`chat_agent_executor.py:154-164`, `chat_agent_executor.py:564-616`) handles `iscoroutinefunction`; a regression that calls an `async def prompt` synchronously would raise `RuntimeError("Use agent.ainvoke()")` at runtime (`chat_agent_executor.py:664-670`) — not caught pre-deploy.
- **Store-injected prompt without store.** `prompt(state, config, *, store)` requires `store` (`test_prompt_with_store:216`); invoking with `store=None` and a store-dependent prompt yields a confusing `AttributeError` or `None` lookup (`store.get(...).value`), with no schema validation in `_get_prompt_runnable`.
- **No calibration of judge prompts.** Example judges (`hub.pull("efriis/self-rag-retrieval-grader")` at `langgraph_self_rag_pinecone_movies.ipynb:180`) are network-fetched without pinning — a Hub update silently changes grading rubric with no dataset re-baseline.
- **Notebook evaluator drift.** `answer_evaluator(run, example)` (`langgraph_crag_local.ipynb:653`) is hard-coded to `score = retrieval_grader.invoke(...)` without versioning; different runners compare against different judge outputs, breaking cross-PR comparability.
- **CI green with broken prompt quality.** Because `ci.yml:80-107` only gates on `lint`/`test`/`integration-test`, a PR that degrades prompt instruction adherence (e.g., truncated `SystemMessage`) merges green.

## Future Considerations

- **Introduce a `prompts/` registry with pinned Hub refs + local YAML fallback.** Store each prompt as a versioned artifact (hash, Hub rev, local copy), exercised by a `libs/prebuilt/tests/test_prompts.py` golden-file suite (input state → expected `LanguageModelInput`). Model after the existing `bench/` baseline pattern.
- **Add a `make eval` / `evals/` harness reusing `FakeToolCallingModel` plus an optional `langsmith` extra.** For local CI: dataset JSONL under `evals/datasets/` (e.g., 50 state→expected-prepared-messages pairs) + deterministic assertion on `_get_prompt_runnable` output. For SaaS: `evals/langsmith/` thin wrapper that calls `langsmith.evaluation.evaluate` with pinned `experiment_prefix` and threshold gate in `ci.yml`.
- **Snapshot prompt rendering, not just graph structure.** Extend `syrupy` usage from `test_react_agent_graph.py:52` to `prompt` rendering: snapshot `agent.nodes["agent"].invoke(state)`'s prepared `messages` for a fixed corpus. Enables `make test` to fail on prompt-prepend regressions.
- **Pin Hub prompts and expose a judge prompt in-repo.** Vendor `efriis/self-rag-retrieval-grader` (or at least pin `hub.pull(..., revision="...")`) and add a `libs/prebuilt/langgraph/prebuilt/evaluators.py` `LLMAsJudge` prompt with explicit rubric + tests — then CI can run judge offline with `FakeToolCallingModel` for rubric-coverage.
- **Add an `evals` CI job with score-threshold gate.** Analogous to `bench.yml:49-63`'s `COMPARE_BENCHMARKS`, run `uv run pytest evals/` and fail if `score < baseline.json` threshold. Keeps the glossary: a prompt change cannot land unless it at least preserves baseline scores.

## Questions / Gaps

- **Do maintainers intend LangGraph to ever ship prompts/datasets?** `README.md:56` and `docs/redirects.json:121-278` suggest deliberate delegation to LangSmith; search found no `RFC` or `AGENTS.md` note on eval ownership. Confirm intent before adding in-repo harness.
- **Is Hub prompt pinning policy documented?** `langgraph_self_rag_pinecone_movies.ipynb:180` uses `hub.pull("efriis/self-rag-retrieval-grader")` without revision. No `CONTRIBUTING.md` guidance on pinning was found (searched `hub.pull` hits only in examples/rag).
- **Are performance baselines considered a proxy for prompt regressions?** `bench/__main__.py:500-516` benchmarks compilation shape; no mapping to prompt quality was found. Interview with maintainers could clarify if perf gates are intentionally the only pre-deploy quality signal.
- **Is `FakeToolCallingModel` intended to grow into a prompt-aware stub?** Current stub (`libs/prebuilt/tests/model.py:28-46`) ignores tool style beyond `bind_tools` dict shape (`model.py:79-93`); no prompt-aware variant exists. Search boundary: `libs/langgraph/tests/fake_chat.py` (not read) may contain a richer fake — gap noted.

---

Generated by `Dimension 12.03: Prompt Evaluation and Experiments` against `langgraph`.
