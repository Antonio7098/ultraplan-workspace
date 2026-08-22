# Source Analysis: openai-agents-sdk

## Dimension 24.01: Public API Surface

### Source Info

| Field | Value |
|-------|-------|
| Name | openai-agents-sdk |
| Path | `studies/agent-harness-study/sources/openai-agents-sdk` |
| Language / Stack | Python 3.10+ (Pydantic v2, openai SDK, hatchling build; MkDocs docs) |
| Analyzed | 2026-08-22 |

> Citation convention: all file paths below are relative to the source root `studies/agent-harness-study/sources/openai-agents-sdk/`.

## Summary

The OpenAI Agents SDK ships as a single PyPI package, `openai-agents` v0.17.6 (`pyproject.toml:2-3`), built from one wheel package `src/agents` (`pyproject.toml:107-108`). Its public surface is a deliberately flat, curated top-level namespace: `src/agents/__init__.py` re-exports 244 symbols via an explicit `__all__` (`src/agents/__init__.py:340-585`) covering the core primitives — `Agent` (`src/agents/agent.py:270`), `Runner.run/run_sync/run_streamed` (`src/agents/run.py:197-283`), the `function_tool` decorator (`src/agents/tool.py:1851`), guardrails, handoffs, run items, sessions, tracing, retry policies, and exceptions — plus a handful of module-level configuration functions (`set_default_openai_key`, etc., `src/agents/__init__.py:270-330`). Domain-specific surfaces live in sub-packages with their own explicit exports (`agents.tracing`, `agents.realtime`, `agents.voice`, `agents.mcp`, `agents.memory`, `agents.sandbox`), and optional integrations are gated behind ~20 extras in `pyproject.toml:37-60`.

The stable/unstable boundary is communicated through three mechanisms rather than runtime metadata: (1) naming convention (`_`-prefixed private modules such as `_config.py`; the internal `run_internal/` package whose docstring states it is "not part of the surface area", `src/agents/run_internal/__init__.py:1-5`); (2) docstring/doc labels ("experimental and not part of the public API" at `src/agents/run.py:153`, "beta" banner at `docs/sandbox_agents.md:5`, experimental Codex tool at `docs/tools.md:771`); and (3) a versioning policy where breaking changes to non-beta public interfaces bump the minor version, with a written breaking-change changelog and migration snippets (`docs/release.md:1-30`). Unusually for an SDK, positional-argument compatibility of exported constructors is enforced by regression tests (`tests/test_source_compat_constructors.py`) and export completeness of public exceptions is tested against `__all__` (`tests/test_exception_exports.py:16-30`).

Weaknesses: the flat 244-symbol top-level namespace mixes domains; auto-generated reference stubs leak internal `run_internal/*` modules into published docs (`docs/ref/run_internal/agent_bindings.md`); a few exported symbols carry mixed internal/public signals (`SessionABC`, `AgentRunner`); and type-check enforcement is asymmetric (mypy strict but pyright "basic" with many checks disabled, `pyrightconfig.json:5-13`).

## Rating

**7 / 10** — Clear model with tests, explicit interfaces, and operational safeguards.

Rationale: the intended API is unambiguous at the entry points (`Agent` + `Runner` + `function_tool`, documented with a runnable hello-world at `docs/index.md:55-74`); stability contracts are real, not aspirational — positional compatibility has dedicated regression tests (`tests/test_source_compat_constructors.py:31-167`) and breaking changes are documented per release (`docs/release.md:15-30`); experimental/beta surfaces are labeled (`docs/sandbox_agents.md:5`, `docs/tools.md:771`, `examples/tools/codex.py:6`). It falls short of 8–9 because the top-level namespace is large and undifferentiated (244 names), generated API-reference pages publish internal modules (`docs/ref/run_internal/`), several exports blur the internal line (`SessionABC` at `src/agents/memory/session.py:57-65` is documented as internal-use yet re-exported at `src/agents/__init__.py:435`; `AgentRunner` appears in `run.py`'s `__all__` at `src/agents/run.py:131` while its sibling functions are marked experimental at `src/agents/run.py:151-166`), and there is no machine-readable stability marker (decorators/metadata) — labels live only in prose.

## Evidence Collected

Every entry includes a file path with line numbers relative to `studies/agent-harness-study/sources/openai-agents-sdk/`.

| Area | Evidence | File:Line |
|------|----------|-----------|
| Package identity | `openai-agents` v0.17.6, Python >=3.10, wheel packages `["src/agents"]` | `pyproject.toml:2-8`, `pyproject.toml:107-108` |
| Top-level export list | Explicit `__all__` with 244 names (counted programmatically) spanning agents, tools, tracing, memory, retry, errors | `src/agents/__init__.py:340-585` |
| Core client objects | `Agent` class; `AgentBase`; `clone()`; `as_tool()` | `src/agents/agent.py:270`, `src/agents/agent.py:174`, `src/agents/agent.py:487`, `src/agents/agent.py:508` |
| Runner entry points | `Runner.run` async signature with `max_turns`, `session`, `conversation_id` params; `run_sync`; loop semantics in docstring | `src/agents/run.py:199-213`, `src/agents/run.py:283`, `src/agents/run.py:217-229` |
| Tool decorator API | Overloaded `function_tool` decorator; `tool_namespace` for grouping | `src/agents/tool.py:1851-1899`, `src/agents/tool.py:1372` |
| Config functions defined on the package | `set_default_openai_key/client/api/responses_transport/harness/agent_registration` wrap private `_config` module | `src/agents/__init__.py:270-330`, `src/agents/_config.py:12-25` |
| Sub-package exports | `sandbox` exposes 34-name `__all__` (Manifest, SandboxAgent, snapshots, errors); `extensions` exposes only `ToolOutputTrimmer` | `src/agents/sandbox/__init__.py:31-65`, `src/agents/extensions/__init__.py:1-3` |
| Realtime/voice sub-SDKs | Full re-export lists for realtime agents/events/models; voice pipeline behind `voice` extra | `src/agents/realtime/__init__.py:1-60`, `pyproject.toml:38`, `pyproject.toml:42` |
| Optional-dependency gating | ~20 extras (litellm, redis, mongodb, e2b, modal, temporal, ...) isolate heavy integrations from core install | `pyproject.toml:37-60` |
| Lazy import boundary | `SQLiteSession` served via module `__getattr__` so optional deps aren't imported eagerly | `src/agents/__init__.py:256-267` |
| Internal-by-convention runtime | `run_internal/` package docstring: public APIs belong top-level, this is "not part of the surface area"; contributor guide mandates new logic lands here | `src/agents/run_internal/__init__.py:1-5`; `AGENTS.md:90-91` |
| Explicit non-public marker | `set_default_agent_runner`/`get_default_agent_runner`: "WARNING: this class is experimental and not part of the public API" | `src/agents/run.py:151-166` |
| Protocol vs ABC guidance | `Session` is the third-party extension point (Protocol); `SessionABC` docstring says it's "intended for internal use" | `src/agents/memory/session.py:14-19`, `src/agents/memory/session.py:57-65` |
| Stability policy | 0.Y.Z scheme: minor bump for breaking changes to non-beta public interfaces; patch for beta/private changes | `docs/release.md:1-13` |
| Breaking-change changelog | Written migration notes with code (e.g., sandbox `LocalFile.src` base_dir change using `SandboxPathGrant`) | `docs/release.md:15-30` |
| Positional-compat tests | Regression tests pin positional order of `RunConfig`, `ModelSettings`, `FunctionTool`, `RunResult(Streaming)`, `ToolContext`, `AgentHookContext`; legacy `_run_impl_task` keyword accepted | `tests/test_source_compat_constructors.py:31-167`, `tests/test_source_compat_constructors.py:343-414` |
| Export-surface test | Every public exception subclass must be re-exported and present in `agents.__all__` | `tests/test_exception_exports.py:16-30` |
| Contributor compat contract | "Treat the parameter and dataclass field order of exported runtime APIs as a compatibility contract"; append-only field additions | `AGENTS.md:51-58` |
| Docs-as-API-checks policy | "Treat runnable docs snippets as API compatibility checks" — verify shown args against implementation | `AGENTS.md:63` |
| Generated API reference | mkdocstrings stubs generated for every non-underscore module; `docs/ref/index.md` renders `::: agents` with explicit member allowlist | `docs/scripts/generate_ref_files.py:51-66`, `docs/ref/index.md:1-14` |
| Reference nav | "API Reference" nav section covers core, models, MCP, tracing, realtime, voice, extensions | `mkdocs.yml:93-191` |
| Runnable quickstart | Hello-world snippet using only `Agent` + `Runner.run_sync` | `docs/index.md:55-74` |
| Example catalog & harness | Examples organized by category with docs index; `run_examples.py` discovers `__main__` examples and runs them in auto mode | `docs/examples.md:5-43`; `examples/run_examples.py:1-13` |
| Examples exercised by tests | Test suite imports example modules directly (e.g., `examples.sandbox.basic`, `examples.sandbox.sandbox_agents_as_tools`) | `tests/test_example_workflows.py:33-42` |
| Experimental labeling | Codex tool: "This surface is experimental and may change"; examples repeat GA caveat; experimental package comment | `docs/tools.md:771`, `examples/tools/codex.py:6`, `src/agents/extensions/experimental/__init__.py:1` |
| Beta labeling | Sandbox agents beta banner ("Expect details of the API ... to change before general availability") | `docs/sandbox_agents.md:5` |
| Leaked internal ref stubs | `docs/ref/run_internal/*.md` generated for internal modules (generator only skips `_`-prefixed files); absent from nav (0 hits for `run_internal` in `mkdocs.yml`) | `docs/ref/run_internal/agent_bindings.md:1-3`, `docs/scripts/generate_ref_files.py:52` |
| Typing asymmetry | mypy `strict = true` vs pyright `"typeCheckingMode": "basic"` with report rules disabled | `pyproject.toml:133-137`, `pyrightconfig.json:4-13` |

## Answers to Dimension Questions

**1. What is the intended public API surface?**
A single import root, `agents`, exposing a curated set of primitives: define agents (`Agent`, `src/agents/agent.py:270`), attach tools via `function_tool`/hosted tool classes (`src/agents/tool.py:1851`, `src/agents/tool.py:381`), execute with `Runner.run/run_sync/run_streamed` (`src/agents/run.py:199-283`), observe via stream events (`src/agents/stream_events.py`) and tracing (`src/agents/tracing/__init__.py`), persist state via `Session` protocol implementations (`src/agents/memory/session.py:14`) and `RunState` (`src/agents/run_state.py`), and configure providers globally via `set_default_openai_*` functions (`src/agents/__init__.py:270-330`). Secondary public roots are `agents.tracing`, `agents.realtime`, `agents.voice`, `agents.mcp`, `agents.memory`, `agents.sandbox`, and `agents.extensions` (deliberately minimal: only `ToolOutputTrimmer`, `src/agents/extensions/__init__.py:1-3`). There is no CLI command group or HTTP server in the library itself; the only operator-facing runner is the repo-local example harness `examples/run_examples.py:1-13`.

**2. Is the stable API easy to distinguish from internal implementation details?**
Mostly yes, via layered conventions: underscore-prefixed private modules excluded even from doc generation (`docs/scripts/generate_ref_files.py:52`); the `run_internal/` package explicitly scoped as non-surface (`src/agents/run_internal/__init__.py:1-5`); explicit warnings on the few runtime-replacement hooks (`src/agents/run.py:151-166`); and `extensions.experimental.*` as the designated unstable zone (`docs/tools.md:771`). However, distinction is by prose and naming only. Three leaks blur the line: (a) `AgentRunner` is listed in `run.py`'s own `__all__` (`src/agents/run.py:130-148`) without any warning while adjacent functions are marked non-public; (b) `SessionABC` is exported at top level (`src/agents/__init__.py:435`) despite its docstring directing third parties to the `Session` protocol instead (`src/agents/memory/session.py:63-64`); (c) generated reference pages for `run_internal/*` are built into the published docs tree (`docs/ref/run_internal/agent_bindings.md:1-3`), just not linked in navigation.

**3. Does the API expose the right level of abstraction for agent harness users?**
Largely yes. The core loop is abstracted behind `Runner` with policy knobs (`RunConfig` fields incl. `ReasoningItemIdPolicy`, `ToolExecutionConfig`, `ToolNotFoundBehavior` — `src/agents/run.py:110-118`), while extension points stay interface-shaped rather than concrete: `Model`/`ModelProvider` protocols (`src/agents/models/interface.py`, imported at `src/agents/__init__.py:85`), `TracingProcessor` (`src/agents/__init__.py:225`), `Session` Protocol (`src/agents/memory/session.py:14-54`), and `RetryPolicy` (`src/agents/retry.py`). The `Session` vs `SessionABC` split shows deliberate thought about which base users should target. Two caveats: streaming consumers touch semi-internal state — `RunResultStreaming` carries underscore-prefixed queues (`_event_queue`, `_input_guardrail_queue`, pinned positionally in `tests/test_source_compat_constructors.py:312-341`), meaning the compat contract protects some private-looking members; and global process-wide configuration functions (`set_default_openai_client`, `src/agents/__init__.py:285`) trade convenience for testability/isolation in multi-tenant hosts.

**4. Are examples sufficient to use the API correctly without reading internals?**
Yes. Coverage spans every major feature area in both docs snippets and runnable files: basic usage and streaming (`examples/basic/` with 24 entries incl. `stream_items.py`, `tool_guardrails.py`), tools incl. hosted MCP and shell/HITL (`examples/tools/shell_human_in_the_loop.py`), handoffs with filters (`examples/handoffs/message_filter.py`), memory/sessions (`examples/memory/`), sandboxes (`examples/sandbox/`), realtime and voice (`examples/realtime/`, `examples/voice/`), plus end-to-end apps (`examples/customer_service`, `examples/research_bot`). The catalog is indexed in docs (`docs/examples.md:5-60`) and examples are wired into CI-adjacent verification — tests import them directly (`tests/test_example_workflows.py:33-42`) and an auto-mode runner can exercise all `__main__` examples (`examples/run_examples.py:1-13`). The docs-index hello world runs with two imports (`docs/index.md:57-68`).

## Architectural Decisions

- **Single-package, flat-top-level distribution.** One wheel (`pyproject.toml:107-108`) with a wide curated `__all__` (`src/agents/__init__.py:340-585`) prioritizes one-import ergonomics over domain modularity; heavier domains are opt-in via extras (`pyproject.toml:37-60`) and lazily imported symbols (`src/agents/__init__.py:260-267`).
- **Compatibility enforced by tests, not annotations.** Positional constructor order of public dataclasses is treated as contract (policy at `AGENTS.md:51-58`; regression suite at `tests/test_source_compat_constructors.py`), including legacy keyword aliases like `_run_impl_task` → `run_loop_task` (`tests/test_source_compat_constructors.py:343-362`).
- **Versioned instability.** A 0.Y.Z scheme makes breaking changes visible in the version number, backed by per-release migration notes with code samples (`docs/release.md:1-30`).
- **Convention-based internal boundary.** Underscore modules + the `run_internal/` package keep the run-loop machinery out of the public narrative while keeping `run.py` as the readable orchestration entrypoint (`AGENTS.md:90-91`).
- **Docs generated from source.** mkdocstrings-driven reference ensures the published API reference cannot drift far from the code (`docs/scripts/generate_ref_files.py:46-73`, `docs/ref/index.md:1-14`).

## Notable Patterns

- **Explicit `__all__` everywhere.** Top-level (`src/agents/__init__.py:340-585`) and each sub-package (`src/agents/sandbox/__init__.py:31-65`, `src/agents/extensions/__init__.py:3`) maintain hand-curated export lists, enabling the export-completeness test pattern seen in `tests/test_exception_exports.py`.
- **Module `__getattr__` lazy exports.** `SQLiteSession` resolves on first attribute access to avoid importing sqlite/optional machinery eagerly (`src/agents/__init__.py:256-267`).
- **Protocol-first extension points with ABC companions.** `Session` (Protocol, `runtime_checkable`) for third parties; `SessionABC` for in-repo implementations (`src/agents/memory/session.py:13-65`) — an intentional dual-track that documents who should implement what.
- **Designated experimental namespace.** `agents.extensions.experimental.codex` isolates unstable integrations with consistent labeling across code comments, docs, and examples (`src/agents/extensions/experimental/__init__.py:1`, `docs/tools.md:771`, `examples/tools/codex.py:6-10`).
- **Docs-as-tests posture.** Contributor policy requires verifying every runnable docs snippet's argument shape against the implementation before merge (`AGENTS.md:63`), treating documentation as part of the API contract.
- **Example-as-fixture.** Tests import example modules and drive their workflows with a `FakeModel` (`tests/test_example_workflows.py:33-51`, `tests/fake_model.py`), keeping examples honest without network calls.

## Tradeoffs

- **Flat namespace vs discoverability.** 244 top-level names make `from agents import *`-style usage trivial but bury domain grouping; users rely on docs nav (`mkdocs.yml:93-191`) to learn that realtime types live under `agents.realtime`, not top level.
- **Prose-only stability labels vs machine-readable markers.** Beta/experimental flags in docstrings/docs are human-friendly but invisible to linters, deprecation tooling, or runtime introspection; nothing programmatic distinguishes `SandboxAgent` (beta) from `Agent` (stable).
- **Generated reference completeness vs leakage.** Generating stubs for all non-underscore modules guarantees coverage but publishes internal `run_internal/*` pages (`docs/ref/run_internal/agent_bindings.md`); excluding them would require either renaming with `_` or extending the generator skip-list.
- **Global default configuration vs test isolation.** `set_default_openai_key/client/api` mutate shared process state (`src/agents/_config.py:12-25`); convenient for scripts, awkward for libraries embedding the SDK.
- **Strict-ish typing, asymmetrically enforced.** mypy strict mode covers the library (`pyproject.toml:133-137`), but pyright runs in "basic" mode with several report rules off (`pyrightconfig.json:4-13`), weakening the guarantee that downstream strict consumers get clean types.

## Failure Modes / Edge Cases

- **Positional breakage risk on dataclass evolution.** The compat tests demonstrate the failure mode they guard: inserting a field mid-order silently shifts every positional caller (`tests/test_source_compat_constructors.py:100-167` pins appended fields `session_settings` → `reasoning_item_id_policy` → `tool_execution` → `tool_not_found_behavior`). The mitigation depends entirely on contributors remembering the AGENTS.md rule; no lint enforces append-only ordering.
- **Internal-module coupling by users.** Because `run_internal` is importable and its helpers appear in generated docs, users can bind to `agents.run_internal.run_loop.NextStepInterruption` (as the project's own test utils do, `tests/utils/hitl.py:12`) and break on refactors with no versioning protection — the 0.Y policy counts these as private-interface changes.
- **Mixed-signal exports mislead integrators.** An integration author copying `AgentRunner` usage from `run.py.__all__` (`src/agents/run.py:131`) has no local warning that sibling replacement hooks are non-public (`src/agents/run.py:153`).
- **Beta-surface churn.** Sandbox agents warn that defaults will change pre-GA (`docs/sandbox_agents.md:5`); the 0.17.0 breaking change to `LocalFile.src`/`LocalDir.src` materialization boundaries shows such changes ship with migrations but still require consumer work (`docs/release.md:15-30`).
- **Experimental tool dependency on external CLI.** The Codex tool shells out to `codex exec --experimental-json` (`src/agents/extensions/experimental/codex/exec.py:62-63`), so its effective API surface includes an unpinned external binary — a stability risk acknowledged by the "may change" label.

## Future Considerations

- Add a machine-readable stability marker (e.g., `typing_extensions.deprecated`-style annotation or a `__stability__` attribute) consumed by docs generation, so beta/experimental status is queryable rather than prose-only.
- Extend `docs/scripts/generate_ref_files.py` to skip the `run_internal/` subtree (not just `_`-prefixed files) or move those stubs out of the built site, closing the internal-docs leak (`docs/scripts/generate_ref_files.py:52`).
- Consider domain sub-modules for the largest clusters (tools ≈ 70 top-level names; tracing ≈ 50) with back-compatible top-level aliases, to keep the root importable-by-scan.
- Reconcile `AgentRunner`'s membership in `run.py.__all__` (`src/agents/run.py:130-148`) with the non-public warnings on `set/get_default_agent_runner` — either mark both consistently or drop it from the export list.
- Align pyright with mypy strictness (or vice versa) so public type hints are validated under both major checkers (`pyrightconfig.json:4-13`, `pyproject.toml:133-137`).

## Questions / Gaps

- No evidence found of automated validation that docs code snippets compile/run (the `AGENTS.md:63` rule is a review-time policy; no script under `docs/scripts/` performs snippet execution — searched `docs/scripts/` for exec/compile hooks, found only `generate_ref_files.py` and `translate_docs.py`).
- No evidence found of semver-style deprecation decorators or warnings (`DeprecationWarning` emissions on renamed public symbols). Searched `src/agents/` for `deprecat*`; only forward-looking mentions exist (e.g., `src/agents/tool.py:1217` noting `LocalShellTool` will be deprecated), with no runtime warning mechanism observed.
- No evidence found of a published, exhaustive compatibility matrix distinguishing "GA" vs "beta" vs "experimental" per symbol; classification must be assembled from scattered banners (`docs/sandbox_agents.md:5`, `docs/tools.md:771`) and the release-notes rule (`docs/release.md:3`).
- Whether `AgentRunner` methods themselves are considered stable is not stated anywhere in the analyzed tree beyond the sibling-function warnings; no direct statement was located (searched `src/agents/run.py` and `AGENTS.md`).

---

Generated by `Dimension 24.01: Public API Surface` against `openai-agents-sdk`.
