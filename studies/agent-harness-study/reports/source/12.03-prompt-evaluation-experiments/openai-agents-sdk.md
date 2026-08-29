# Source Analysis: openai-agents-sdk

## Dimension 12.03: Prompt Evaluation and Experiments

### Source Info

| Field | Value |
|-------|-------|
| Name | openai-agents-sdk |
| Path | `studies/agent-harness-study/sources/openai-agents-sdk` |
| Language / Stack | Python / OpenAI Agents SDK (Pydantic, Responses API, pytest, inline-snapshot, coverage) |
| Analyzed | 2026-08-29 |

## Summary

openai-agents-sdk is an agent execution harness, not a prompt-product repo. It provides a `Prompt` TypedDict and `DynamicPromptFunction` plumbing to forward OpenAI platform-managed prompt templates to the Responses API (`src/agents/prompts.py:23-47`, `src/agents/agent.py:325-328`), plus instruction helpers and handoff prefixes (`src/agents/extensions/handoff_prompt.py:3-12`). It has no first-party prompt evaluation framework: no eval datasets, no experiment tracker, no LLM-as-judge evaluation harness, and no prompt-regression gate in CI. Prompt correctness is verified only via plumbing unit tests (`tests/test_agent_prompt.py:52`, `tests/test_prompt_cache_key.py:34`) and snapshot-guarded tracing, not against expected LLM outputs. The closest evaluative artifacts are illustrative example patterns for `llm_as_a_judge` (`examples/agent_patterns/llm_as_a_judge.py:11`) and tracing export hooks to external observability vendors, which are intentionally out-of-scope for the SDK itself.

## Rating

**Score: 2 / 10 — Absent / Implicit / Ad-hoc**

Rationale: Across all five rubric checks the SDK scores in the 1–3 band. Evaluation datasets, experiment tracking, judge prompts as QA, and pre-deploy prompt regression tests are missing by design. Testing verifies that a `Prompt` object reaches `Model.get_response(prompt=…)` (`tests/test_agent_prompt.py:99-104`, `src/agents/prompts.py:56-82`), not that a changed prompt still yields correct answers. LLM-as-judge exists only as an example pattern (`examples/agent_patterns/llm_as_a_judge.py:31-39`) exercised by a `FakeModel` functional test (`tests/test_example_workflows.py:90-151`), not as a reusable evaluation suite. CI (` .github/workflows/tests.yml:68-108`) gates on `pytest`/`coverage --fail-under=85` (`Makefile:59`) with no prompt-eval job. This matches a harness that deliberately delegates evaluation to consumer applications.

## Evidence Collected

Every entry MUST include a file path with line numbers. Format: `path/to/file.ts:NN`.

| Area | Evidence | File:Line |
|------|----------|-----------|
| Eval datasets | No `evals/`, `datasets/`, `benchmarks/` or golden-file directory exists; `Glob **/evals/**` returns empty; only `examples/research_bot/sample_outputs/` contains static demo outputs | `studies/agent-harness-study/sources/openai-agents-sdk/examples/research_bot/sample_outputs/product_recs.md:1` , `Glob` empty result |
| Eval datasets | Sandbox memory fixtures use rollout JSON dumps as synthetic data, not curated eval sets (`dump_rollout_json`, `render_phase_one_prompt`) | `src/agents/sandbox/memory/prompts/rollout_extraction_prompt.md:38` , `tests/sandbox/test_memory.py:268-275` |
| Prompt plumbing | `Prompt` TypedDict defines platform-managed template contract (`id`, `version`, `variables`) | `src/agents/prompts.py:23-34` |
| Prompt plumbing | `PromptUtil.to_model_input()` resolves static dict or `DynamicPromptFunction` into `ResponsePromptParam` | `src/agents/prompts.py:56-82` |
| Prompt plumbing | `Agent.prompt: Prompt \| DynamicPromptFunction \| None` and `Agent.get_prompt()` delegation | `src/agents/agent.py:325-328` , `src/agents/agent.py:1073-1083` |
| Prompt plumbing | Dynamic prompt constructor with `GenerateDynamicPromptData` context+agent | `src/agents/prompts.py:36-47` |
| Prompt tests – plumbing | 4 async tests verify static/dynamic resolution and that prompt reaches model unchanged | `tests/test_agent_prompt.py:52-68` , `tests/test_agent_prompt.py:72-83` , `tests/test_agent_prompt.py:87-104` , `tests/test_agent_prompt.py:115-144` |
| Prompt tests – cache key | Distinct prompt-cache-key tests verify `ModelSettings.extra_args["prompt_cache_key"]` lifecycle, not prompt quality | `tests/test_prompt_cache_key.py:34-44` , `tests/test_prompt_cache_key.py:61-76` , `tests/test_prompt_cache_key.py:232-247` |
| Prompt tests – instructions | Dynamic vs static instruction helpers tested separately from `Prompt` template feature | `tests/test_agent_instructions_signature.py:24-109` , `tests/test_agent_config.py:18-30` |
| Experiment trackers | No MLflow/W&B/Comet import in `src/`; only docs list tracing export integrations as optional external sinks | `docs/tracing.md:200-212` , `docs/tracing.md:222` |
| Experiment trackers | `pyproject.toml` dependencies contain no `mlflow`, `wandb`, `comet`, `promptlayer` runtime dep; only `coverage`, `inline-snapshot`, `pytest-xdist` dev deps | `pyproject.toml:9-18` , `pyproject.toml:62-94` |
| LLM-as-judge – pattern | Example implements loop: `story_outline_generator` → `evaluator: Agent[EvaluationFeedback]` → loop until `score == "pass"` | `examples/agent_patterns/llm_as_a_judge.py:16-22` , `examples/agent_patterns/llm_as_a_judge.py:25-39` , `examples/agent_patterns/llm_as_a_judge.py:54-83` |
| LLM-as-judge – second example | Deterministic guardrail/judge variant judges scifi story | `examples/agent_patterns/deterministic.py:31` |
| LLM-as-judge – streaming guardrail | Streaming guardrail example asks model to judge response | `examples/agent_patterns/streaming_guardrails.py:41` |
| LLM-as-judge – test | `test_llm_as_judge_loop_handles_dataclass_feedback` exercises the pattern with `FakeModel` + `EvaluationFeedback` output_type, loop over `Runner.run(judge_agent)` | `tests/test_example_workflows.py:90-151` |
| LLM-as-judge – harness | `src/agents/testing/` exports `ScriptedModel`, `ModelCall`, `ModelStep` test doubles, not an eval harness | `src/agents/testing/__init__.py:3-14` , `src/agents/testing/model.py:1` |
| Regression suites | Snapshot tooling via `inline-snapshot` is for span/tool plumbing, not prompt output (`fetch_normalized_spans() == snapshot(...)`) | `tests/test_agent_tracing.py:7` , `tests/test_agent_tracing.py:53` , `pyproject.toml:77` , `Makefile:61-67` |
| Regression suites | Handoff instruction prefix has a unit test, but it asserts string concatenation, not LLM behavior | `src/agents/extensions/handoff_prompt.py:15-19` , `tests/test_handoff_prompt.py:7-12` |
| CI – prompt test integration | `tests.yml` defines `lint`, `typecheck`, `tests` (matrix py3.10-3.14), `tests-windows`, `build-docs`; prompt jobs absent; tests gated by `detect-changes.sh code/docs` | `.github/workflows/tests.yml:17-42` , `.github/workflows/tests.yml:44-67` , `.github/workflows/tests.yml:68-108` |
| CI – thresholds | `make coverage` runs `coverage report -m --fail-under=85` – coverage gate, not eval gate | `Makefile:54-59` , `.github/workflows/tests.yml:97-99` |
| Docs – prompt guidance | Docs describe how to create/use prompt templates via platform playground, not how to evaluate them | `docs/agents.md:65-92` , `docs/agents.md:94-122` |
| Fake model | `FakeModel.get_response()` / `FakeModel.stream_response()` record `last_turn_args`/`first_turn_args` for plumbing tests; responses are caller-supplied | `tests/fake_model.py:51-145` , `tests/fake_model.py:147-175` |

## Answers to Dimension Questions

### 1. Are prompt changes tested?

**Partially, but not as evaluation.** Prompt *plumbing* is unit-tested: static dict equality (`tests/test_agent_prompt.py:64-68`), dynamic function resolution (`tests/test_agent_prompt.py:72-83`), model-receipt assertion (`tests/test_agent_prompt.py:99-104`), and omitted `model/tools` when `prompt` owns the model (`tests/test_agent_prompt.py:115-144`). `PromptCacheFakeModel` tests cover `prompt_cache_key` stamping (`tests/test_prompt_cache_key.py:34-44`). No test in `tests/` asserts that a prompt change preserves an expected LLM output or passes a golden dataset. Instruction-signature tests (`tests/test_agent_instructions_signature.py:24`) and snapshot tracing tests (`tests/test_agent_tracing.py:53`) guard orchestration, not prompt semantics. **Verdict: plumbing-tested, eval-untested.**

### 2. Are experiments tracked?

**No.** The codebase contains no experiment tracker abstraction, run registry, metric logger, or A/B runner. `pyproject.toml:9-18` lists no MLflow/W&B/Comet/PromptLayer runtime dependency. The only related surface is tracing export documentation pointing to external vendors (Weights & Biases, MLflow, PromptLayer, Comet Opik at `docs/tracing.md:200-222`), and the `ScriptedModel` test double (`src/agents/testing/__init__.py:9`) for deterministic unit tests. There is no SDK-managed prompt versioning, comparison dashboard, or experiment ID threaded through `RunConfig`/`Agent`.

### 3. Is LLM-as-judge used for evaluation?

**Illustrative only, not evaluative.** The repo ships `examples/agent_patterns/llm_as_a_judge.py:11-39` – a two-agent loop where an `evaluator` agent returns `EvaluationFeedback(feedback, score)` – as a usability pattern, not an eval suite. `tests/test_example_workflows.py:90` replays that loop with `FakeModel` to verify looping logic, not judge quality. No production-grade judge prompt library, scoring rubric, calibration set, or aggregate report exists in `src/agents/`. `src/agents/sandbox/memory/prompts/` templates are task memory prompts, not evaluation judges.

### 4. Are regressions caught before deployment?

**Structural regressions yes, prompt regressions no.** CI at `.github/workflows/tests.yml:68-108` blocks merges via `make lint` / `make typecheck` / `make tests` / `make coverage --fail-under=85` (`Makefile:59`). `inline-snapshot` guards span shapes and tool serialization (`tests/test_agent_tracing.py:53`), and coverage fails below 85%. However there is no prompt-eval job comparing new vs baseline on a held-out set, no threshold check (e.g., accuracy drop), and no required `prompt-id@version` golden comparison. A prompt template change (`src/agents/prompts.py:76` returning a different `id`/`variables`) passes CI if plumbing tests still hold, even if downstream answer quality regresses. Can a prompt change be deployed with confidence it won't regress? No.

## Architectural Decisions

- **Platform-managed prompt pass-through (`src/agents/prompts.py:23-82`)**: SDK treats prompts as opaque `ResponsePromptParam` resolved at runtime and forwarded to the Responses API. Keeps SDK agnostic to prompt authoring/versioning; pushes evaluation responsibility to the platform (prompt playground) and the consumer.
- **Separation of `instructions` vs `prompt` (`src/agents/agent.py:309-328`, `docs/agents.md:31-32`)**: Static/dynamic `instructions` are the SDK-native system prompt; `prompt` is the optional hosted-template overlay. Prevents conflation and allows either to be mocked via `FakeModel`, but means there are two distinct prompt surfaces to evaluate with no unified eval.
- **Test doubles over live model calls (`tests/fake_model.py:51-145`, `src/agents/testing/__init__.py:3-14`)**: `FakeModel` + `ScriptedModel` make tests deterministic and offline by injecting caller-provided `TResponseOutputItem`s. Enables fast PR gating without flaky LLM judgments, but removes any signal about prompt-induced quality shifts.
- **Tracing-via-external-export rather than internal eval store (`docs/tracing.md:200-222`, `src/agents/tracing/`)**: Provides a processor interface for third-party evaluation/observability (W&B, MLflow, PromptLayer) without embedding dataset or metric storage inside the SDK. Aligns with harness role but leaves consumers to wire their own evaluation pipeline.

## Notable Patterns

- **Dynamic prompt function pattern** – `DynamicPromptFunction = Callable[[GenerateDynamicPromptData], MaybeAwaitable[Prompt]]` at `src/agents/prompts.py:47` with `PromptUtil.to_model_input()` awaiting sync/async callbacks (`src/agents/prompts.py:70-76`). Tested in `tests/test_agent_prompt.py:72-83` and sandbox dynamic prompt case (`tests/sandbox/test_runtime.py:1653-1676`).
- **Handoff instruction augmentation** – `RECOMMENDED_PROMPT_PREFIX` + `prompt_with_handoff_instructions()` at `src/agents/extensions/handoff_prompt.py:3-19`, documented at `docs/handoffs.md:139-145`; verified by a trivial concatenation test (`tests/test_handoff_prompt.py:7-12`), not by behavioral eval.
- **LLM-as-judge as application pattern, not SDK service** – `examples/agent_patterns/llm_as_a_judge.py:31-39` defines `evaluator: Agent[EvaluationFeedback]` with a `Literal["pass","needs_improvement","fail"]` schema; the consumer owns the judge loop (`examples/agent_patterns/llm_as_a_judge.py:56-83`). Test coverage (`tests/test_example_workflows.py:90-151`) confirms the pattern runs with fakes.
- **Snapshot-gated tracing** – `inline-snapshot` asserts on `fetch_normalized_spans()` (`tests/test_agent_tracing.py:53`), a pattern that *could* be repurposed for prompt-regression snapshots but currently only guards plumbing.

## Tradeoffs

- **Determinism vs fidelity**: Using `FakeModel`/`PromptCacheFakeModel` (`tests/fake_model.py:338`) yields fast, reproducible CI with `OPENAI_API_KEY=fake-for-tests` (`.github/workflows/tests.yml:80`) but sacrifices detection of real LLM drift when prompt wording changes. Any quality regression is invisible until production.
- **SDK leanness vs evaluation completeness**: Omitting a built-in eval runner and dataset store keeps `pyproject.toml:9-18` minimal and avoids opinionated metric/rubric choices. The cost is every consumer re-invents prompt testing and no community-standard benchmark emerges from the harness.
- **External tracing sinks vs built-in experiment tracking**: Documenting exports to W&B/MLflow/PromptLayer (`docs/tracing.md:200-222`) offloads A/B comparison to vendors, preserving vendor neutrality. Tradeoff is no SDK-enforced experiment ID, no first-class `prompt_version → metrics` join, and no CI-visible comparison report.
- **Two prompt surfaces (`instructions` + hosted `prompt`)**: Flexibility for platform-authored templates without SDK lock-in, but doubles the prompt-regression surface and the eval burden; neither surface is gated today.

## Failure Modes / Edge Cases

- **Silent prompt quality decay**: Changing `prompt.variables` or `instructions` string passes `tests/test_agent_prompt.py:52-144` (equality on plumbing only). A subtle instruction weakening (e.g., removing "never give it a pass on the first try" from `examples/agent_patterns/llm_as_a_judge.py:36`) has no failing test, even though downstream behavior diverges.
- **Dynamic prompt runtime exception swallowed as `UserError`**: `PromptUtil.to_model_input()` at `src/agents/prompts.py:76` raises `UserError("Dynamic prompt function must return a Prompt")` for non-dict returns, but type errors inside the user function surface as generic model errors with no evaluation context to localize the prompt that triggered them.
- **`Prompt` bypasses model/tools serialization assumptions**: When `prompt` owns the model, `ModelSettings(tool_choice="computer")` caveats apply (`docs/tools.md:196`, `docs/models/index.md:79-83`). A prompt change that pins a different model without updating harness-side tool-choice handling can break `test_convert_tools_prompt_managed_*` expectations (`tests/models/test_openai_responses_converter.py:490-1073`) only if explicitly covered; most consumer prompts are uncovered.
- **Judge calibration failure with no ground truth**: The `evaluator` example relies on the judge LLM's discretion ("After 5 attempts... do not go for perfection" at `examples/agent_patterns/llm_as_a_judge.py:36`). Without a labeled eval set, judge leniency/stricness drift is undetectable; `tests/test_example_workflows.py:100-139` injects canned judge outputs that never validate calibration.
- **CI change-detector can skip tests**: `.github/workflows/tests.yml:24` uses `detect-changes.sh code` to skip `lint/typecheck/tests` on doc-only PRs. A docs-only PR that edits prompt guidance (`docs/agents.md:71-90`) or example judge prompts (`examples/agent_patterns/llm_as_a_judge.py:31-39`) correctly skips eval-unrelated tests, but also would skip any future prompt-eval job unless it is tagged as `docs`-sensitive.
- **Coverage threshold masks prompt surface**: `make coverage --fail-under=85` (`Makefile:59`) can stay green while prompt-related lines remain uncovered; coverage excludes sandbox materialization paths (`pyproject.toml:171-182`) but does not separately enforce prompt/template coverage.

## Future Considerations

- **Add a minimal prompt-eval fixture layer**: Introduce `tests/evals/prompts/*.jsonl` golden sets (input → expected prompt `id/version/variables` → expected output excerpt) exercised by a new `pytest -k eval` job; reuse `FakeModel` with recorded real-model traces to avoid live API calls while still asserting output stability.
- **Formalize `Prompt` version pinning and diff tool**: Store prompt `id@version` snapshots checked into `tests/snapshots/prompts/` and add `uv run pytest --inline-snapshot=fix` coverage for them, so PRs changing `src/agents/prompts.py:78-82` output require explicit snapshot review.
- **Promote `llm_as_a_judge` from example to test utility**: Factor judge agent (`examples/agent_patterns/llm_as_a_judge.py:31-39`) into `src/agents/testing` as `judge_model(pattern, rubric)` with a calibrated rubric file, allowing optional `RUN_LIVE_EVAL=1` nightly workflow that runs the judge against a small held-out set and publishes metrics as GitHub job summaries.
- **Bridge tracing to eval metrics**: Emit prompt `id/version` as span attributes and document a canonical query (PromptLayer / W&B) that compares pass-rate/error-rate across prompt versions, then wire a lightweight CI check that fails on statistically significant drop (e.g., >5% on 100-case eval set).
- **Gate hosted prompt changes**: If `prompt` is set, add a CI step that resolves the template via the platform API in a dry-run mode and asserts variable coverage (no missing `{{poem_style}}`), catching template breakage before release.

## Questions / Gaps

- **No evidence of evaluation dataset**: Searched `src/`, `tests/`, `examples/`, `docs/`, `pyproject.toml`, and `Glob **/evals/**`; no dataset manifests or loaders found. Confirm whether evaluation lives outside this repo (e.g., private OpenAI evals) – none discoverable in-tree.
- **No experiment tracking code**: `grep experiment|mlflow|wandb|comet` in `src/agents` yields zero runtime imports; only docs external links exist. No A/B CLI or `RunConfig.experiment_id` field found.
- **Judge prompts not versioned or tested for quality**: `examples/agent_patterns/llm_as_a_judge.py:31-37` judge instructions are free-form strings with no rubric file, no inter-rater agreement data, and no calibration test.
- **CI has no deployment gate beyond tests**: `publish.yml`/`release-pr.yml`/`release-tag.yml` exist but contain no prompt-eval prerequisite; question "are regressions caught before deployment?" is answered "No evidence found" for prompt-specific gates.
- **Unproven assumption**: Analysis assumes hosted prompt semantics are validated on the platform side; no in-repo contract test hits `openai.types.responses.response_prompt_param.ResponsePromptParam` beyond pass-through.

---

Generated by `Dimension 12.03: Prompt Evaluation and Experiments` against `openai-agents-sdk`.
