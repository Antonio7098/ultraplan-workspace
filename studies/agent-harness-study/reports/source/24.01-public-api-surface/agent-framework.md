# Source Analysis: agent-framework

## Dimension 24.01: Public API Surface

### Source Info

| Field | Value |
|-------|-------|
| Name | agent-framework |
| Path | `studies/agent-harness-study/sources/agent-framework` |
| Language / Stack | C#/.NET (`dotnet/src`, 37 projects) and Python 3.10+ (`python/packages`, 30 packages); NuGet + PyPI distribution |
| Analyzed | 2026-08-22 |

## Summary

Microsoft Agent Framework is a dual-language (`.NET` + `Python`) monorepo whose public API surface is managed deliberately at four levels. First, distribution is package-based with an explicit lifecycle: `python/PACKAGE_STATUS.md:14-46` classifies all 30 Python packages into `alpha`/`beta`/`rc`/`released`/`deprecated` buckets, while the .NET side encodes the same idea in MSBuild properties (`IsReleased`, `IsReleaseCandidate`) consumed by central versioning logic (`dotnet/nuget/nuget-package.props:4-11`). Second, each language has a curated entry point: Python's `agent_framework/__init__.py` declares itself "Public API surface for Agent Framework core" in its module docstring (`python/packages/core/agent_framework/__init__.py:3-9`) and enforces it with a ~275-entry `__all__` list (`python/packages/core/agent_framework/__init__.py:342-618`); .NET exposes a small abstraction core (`AIAgent`, `AgentSession`, `ChatClientAgent`) in a dedicated `Microsoft.Agents.AI.Abstractions` package (`dotnet/src/Microsoft.Agents.AI.Abstractions/AIAgent.cs:38`) plus a fluent `AIAgentBuilder` (`dotnet/src/Microsoft.Agents.AI/AIAgentBuilder.cs:16`). Third, stability is machine-enforced rather than convention-only: experimental APIs carry `[Experimental("MAAI001")]` compiler diagnostics on .NET (`dotnet/src/Microsoft.Agents.AI.Harness/HarnessAgent.cs:83`, ID defined at `dotnet/src/Shared/DiagnosticIds/DiagnosticsIds.cs:16`) and `@experimental(feature_id=...)` decorators that emit runtime warnings and inject docstring warnings in Python (`python/packages/core/agent_framework/_feature_stage.py:383-403`). Fourth, shipped packages are guarded by API-compat validation against a baseline version, with deviations tracked in `CompatibilitySuppressions.xml` (`dotnet/nuget/nuget-package.props:19-24`, `dotnet/src/Microsoft.Agents.AI/CompatibilitySuppressions.xml:4-11`).

The intended consumer experience is documented end-to-end: the root README ships runnable quickstart snippets for both languages (`README.md:87-91`, `README.md:114-115`) and links to progressive sample trees (`python/samples/01-get-started/`, `dotnet/samples/01-get-started/`). Optional provider integrations are lazy-loaded through namespace shims so the core import stays cheap and missing dependencies fail with actionable messages naming the required PyPI package (`python/packages/core/agent_framework/openai/__init__.py:17-40`).

## Rating

**8 / 10** — The public surface is explicit, staged, and enforced by tooling on both languages: curated exports with `__all__` (`python/packages/core/agent_framework/__init__.py:342`), lifecycle tables per package (`python/PACKAGE_STATUS.md:20`), runtime warning + docstring injection for experimental APIs (`python/packages/core/agent_framework/_feature_stage.py:29-40`, `254-331`), compiler-level `[Experimental]` gating on .NET (`dotnet/src/Shared/DiagnosticIds/DiagnosticsIds.cs:15-31`), binary-compat package validation (`dotnet/nuget/nuget-package.props:19-21`), and ADR-governed naming decisions (`docs/decisions/0005-python-naming-conventions.md:1-30`). It falls short of 9–10 because: (a) internal Python modules define no `__all__`, so direct `import *` from underscore modules leaks internals — protection rests solely on the underscore naming convention; (b) cross-language parity must be maintained manually and already shows drift signals (Python 1.9.0 vs .NET 1.10.0 version prefixes at `python/packages/core/pyproject.toml:7` vs `dotnet/nuget/nuget-package.props:4`; harness APIs experimental in both but packaged differently — separate `Microsoft.Agents.AI.Harness` package on .NET vs in-core `_harness/` modules behind a feature flag on Python, `python/packages/core/agent_framework/__init__.py:85-147`); and (c) recent CHANGELOGs still contain multiple `[BREAKING]` entries per release (`python/CHANGELOG.md:22-24`), indicating the "released" tier is young.

## Evidence Collected

Every entry includes a file path with line numbers.

| Area | Evidence | File:Line |
|------|----------|-----------|
| Declared public surface (Python) | Module docstring "Public API surface for Agent Framework core" + full `__all__` of ~275 names covering agents, clients, tools, middleware, sessions, skills, workflows, compaction, evaluation | `python/packages/core/agent_framework/__init__.py:3-9`, `342-618` |
| Core abstractions (.NET) | `public abstract partial class AIAgent` with `Id`, `Name`, `Description`, `GetService`, `RunAsync` overloads as the primary consumer entry point | `dotnet/src/Microsoft.Agents.AI.Abstractions/AIAgent.cs:38`, `58`, `82`, `94`, `118`, `251` |
| Builder pattern | `AIAgentBuilder` with `Use(...)` decorator factories and `UseAIContextProviders` | `dotnet/src/Microsoft.Agents.AI/AIAgentBuilder.cs:16`, `76-112`, `175` |
| Package lifecycle table | All 30 Python packages bucketed alpha/beta/rc/released/deprecated; umbrella `agent-framework` = released, `core` = released, `gemini` = alpha | `python/PACKAGE_STATUS.md:14-46` |
| Feature-stage decorators | `ExperimentalFeature` enum (13 IDs incl. HARNESS, SKILLS, EVALS), `experimental()` and `release_candidate()` decorators; docstring warning blocks injected | `python/packages/core/agent_framework/_feature_stage.py:43-66`, `383-403`, `29-40` |
| Runtime warnings for staged APIs | `ExperimentalWarning(FutureWarning)` emitted via wrapped `__new__`/`__init_subclass__`/callables, once per feature id; user-frame resolution to point warnings at consumer code | `python/packages/core/agent_framework/_feature_stage.py:79-84`, `254-331`, `192-250` |
| Explicit non-stability of stage metadata | Enum docstrings state they are "a stage-scoped inventory, not a stable introspection surface"; `__feature_id__` may disappear | `python/packages/core/agent_framework/_feature_stage.py:44-51`, `69-76` |
| Experimental gating (.NET) | `[Experimental(DiagnosticIds.Experiments.AgentsAIExperiments)]` (= `MAAI001`) on `HarnessAgent`, `HarnessAgentOptions`, and ~20 other files | `dotnet/src/Microsoft.Agents.AI.Harness/HarnessAgent.cs:83`, `dotnet/src/Microsoft.Agents.AI.Harness/HarnessAgentOptions.cs:17`, `dotnet/src/Shared/DiagnosticIds/DiagnosticsIds.cs:16` |
| Reuse of upstream experiment IDs | MEAI001/OPENAI001 IDs aliased "so consumers do not need to suppress additional diagnostics" | `dotnet/src/Shared/DiagnosticIds/DiagnosticsIds.cs:21-31` |
| Binary compatibility enforcement | `PackageValidationBaselineVersion=1.0.0`, `EnablePackageValidation` for released packages | `dotnet/nuget/nuget-package.props:19-24` |
| Breaking-change bookkeeping | Per-package `CompatibilitySuppressions.xml` recording removed members (CP0002) against baseline assemblies | `dotnet/src/Microsoft.Agents.AI/CompatibilitySuppressions.xml:4-11` |
| Release-tier flags | `<IsReleased>true</IsReleased>` on Abstractions/Core/Workflows/OpenAI/Foundry; RC flag true only on Declarative, GitHub.Copilot, Purview | `dotnet/src/Microsoft.Agents.AI.Abstractions/Microsoft.Agents.AI.Abstractions.csproj:4`, `dotnet/src/Microsoft.Agents.AI/Microsoft.Agents.AI.csproj:4`, `dotnet/src/Microsoft.Agents.AI.Declarative/Microsoft.Agents.AI.Declarative.csproj` |
| Lazy optional-dependency namespaces | `openai/__init__.py` maps symbol → `(module, pip package)` and raises `ModuleNotFoundError` telling users which package to install; also provides `__dir__` | `python/packages/core/agent_framework/openai/__init__.py:17-47` |
| Provider namespaces | `foundry/` namespace aggregates anthropic, contentunderstanding, foundry, foundry-local symbols under one import path | `python/packages/core/agent_framework/foundry/__init__.py:6-30` |
| Umbrella package | `agent-framework` meta-package depends only on `agent-framework-core[all]`, which lists every optional integration as extras | `python/pyproject.toml:7-18`, `python/packages/core/pyproject.toml:36-63` |
| Runnable quickstart examples | Root README quickstart code blocks; `01_hello_agent.py` creates `FoundryChatClient` + `Agent` and runs streaming/non-streaming | `README.md:87-91`, `114-115`; `python/samples/01-get-started/01_hello_agent.py:12-32` |
| Progressive tutorial ladder | Parallel numbered sample trees (get-started → agents → workflows → hosting → end-to-end) for both languages | `python/samples/01-get-started/` (8 files), `dotnet/samples/01-get-started/` (6 projects), `README.md:172-177` |
| Sample authoring standard | PEP 723 inline metadata requirement making every sample self-contained and runnable | `python/samples/SAMPLE_GUIDELINES.md:5-17`, `28-42` |
| Naming governance | ADR-0005 records Python renames chosen for idiomatic naming "while preserving discoverability and mapping to the .NET names" | `docs/decisions/0005-python-naming-conventions.md:13,18` |
| ADR process | Numbered ADRs (28+) govern API decisions such as TypedDict options (0012), context/middleware (0016), skills design (0021) | `docs/decisions/README.md:1-10`, `docs/decisions/0012-python-typeddict-options.md` |
| Semantic versioning + changelog discipline | Keep-a-Changelog format with per-package BREAKING markers; orchestrations promoted to stable | `python/CHANGELOG.md:1-5`, `22-24`, `26` |
| Extension protocols exported as public API | `SupportsAgentRun`, `SupportsChatGetResponse`, tool capability protocols exported from core | `python/packages/core/agent_framework/__init__.py:20-32`, `519-527` |
| Internal boundary by convention only | Underscore-prefixed implementation modules (`_types.py`, `_agents.py`) define no `__all__`; curation happens only at the top-level `__init__.py` | `python/packages/core/agent_framework/_types.py` (no `__all__`), `__init__.py:219-257` |
| Internals visible only to tests (.NET) | `InternalsVisibleTo` limited to unit-test assemblies | `dotnet/src/Microsoft.Agents.AI/Microsoft.Agents.AI.csproj:50-54` |
| Legacy/shim code quarantined | `LegacySupport/` contains compile-time attribute shims injected only for legacy TFMs, never public API | `dotnet/src/LegacySupport/CompilerFeatureRequiredAttribute/README.md:1-6`, `dotnet/src/Microsoft.Agents.AI/Microsoft.Agents.AI.csproj:11-17` |
| External documentation entry points | learn.microsoft.com docs linked from READMEs and NuGet metadata (`PackageProjectUrl`) | `README.md:96-100`, `dotnet/README.md:31-35`, `dotnet/nuget/nuget-package.props:46` |

## Answers to Dimension Questions

### 1. What is the intended public API surface?

Three tiers. **Tier 1 — stable packages**: Python `agent-framework-core`, `-openai`, `-foundry`, `-orchestrations` and the umbrella `agent-framework` are `released` (`python/PACKAGE_STATUS.md:20-45`); .NET equivalents ship without prerelease suffixes via `<IsReleased>true</IsReleased>` (`dotnet/src/Microsoft.Agents.AI/Microsoft.Agents.AI.csproj:4`). Their consumers use a small set of root types: `Agent`/`BaseAgent`/chat-client bases/tools/middleware/workflows from one import path (`python/packages/core/agent_framework/__init__.py:342-618`), or `AIAgent`/`AgentSession`/`AIAgentBuilder` on .NET (`dotnet/src/Microsoft.Agents.AI.Abstractions/AIAgent.cs:38`, `dotnet/src/Microsoft.Agents.AI/AIAgentBuilder.cs:16`). **Tier 2 — staged features inside stable packages**: marked by `@experimental`/`@release_candidate` decorators in Python (`python/packages/core/agent_framework/_feature_stage.py:383-403`) and `[Experimental(MAAI001)]` attributes in .NET (`dotnet/src/Shared/DiagnosticIds/DiagnosticsIds.cs:15-17`), inventoried per feature in `python/PACKAGE_STATUS.md:48-77`. **Tier 3 — incubation**: the whole `agent-framework-lab` package explicitly states lab modules "may experience breaking changes or be deprecated" (`python/packages/lab/README.md:3-5`).

### 2. Is the stable API easy to distinguish from internal implementation details?

Largely yes, with one asymmetry. On .NET the distinction is structural: abstractions live in a dedicated package, internals are `internal` plus narrowly scoped `InternalsVisibleTo` test grants (`dotnet/src/Microsoft.Agents.AI/Microsoft.Agents.AI.csproj:50-54`), and shared source under `dotnet/src/Shared` is compiled-in rather than published. On Python the distinction is by curated export: the top-level `__all__` defines the contract, underscore modules signal privacy, and lazy namespaces document their re-export sets (`python/packages/core/agent_framework/openai/__init__.py:17-31`). However, nothing mechanically prevents `from agent_framework._types import *` — those modules have no `__all__` (verified across all 15 top-level modules; e.g., `python/packages/core/agent_framework/_types.py`), so accidental coupling to internals is possible for users who ignore conventions. The framework even documents that its own stability metadata (`__feature_id__`) is not a stable introspection surface (`python/packages/core/agent_framework/_feature_stage.py:46-51`) — honest, but it means tooling cannot rely on attributes either.

### 3. Does the API expose the right level of abstraction for agent harness users?

Yes. The core contract is intentionally narrow: an agent is something that satisfies `run` (`SupportsAgentRun` protocol, `python/packages/core/agent_framework/__init__.py:20`), a chat client implements two `_inner_get_response*` methods (custom-client recipe in `python/packages/core/AGENTS.md`, "Custom Chat Client" section; base at `python/packages/core/agent_framework/__init__.py:22`), and cross-cutting behavior composes through three middleware kinds (`AgentMiddleware`, `ChatMiddleware`, `FunctionMiddleware`, `python/packages/core/agent_framework/__init__.py:149-168`) instead of subclassing providers. Provider specifics stay behind capability protocols (`SupportsMCPTool`, `SupportsShellTool`, etc., `python/packages/core/agent_framework/__init__.py:25-31`). On .NET, `AIAgent.GetService(Type, object?)` provides service-location escape hatches without leaking concrete client types (`dotnet/src/Microsoft.Agents.AI.Abstractions/AIAgent.cs:108-118`), and `DelegatingAIAgent` enables wrapper-based extension. Runtime internals (edge runners, executor graphs) are re-exported only where they are genuine extension points (`create_edge_runner`, `RunnerContext`, `WorkflowContext` at `python/packages/core/agent_framework/__init__.py:286`, `310-314`); the harness loop internals stay behind `create_harness_agent` (`python/packages/core/agent_framework/__init__.py:85-88`).

### 4. Are examples sufficient to use the API correctly without reading internals?

For getting started and common scenarios, yes: eight self-contained get-started scripts escalate from hello-agent → tools → multi-turn → memory → workflows → hosting (`python/samples/01-get-started/01_hello_agent.py:12-32` through `08_host_your_agent.py`), mirrored on .NET (`dotnet/samples/01-get-started/`, six projects), plus deep topic trees (middleware, MCP, observability, hosting) indexed in the READMEs (`README.md:41-51`, `172-177`). Samples are governed by a written standard requiring PEP 723 metadata so each file runs standalone (`python/samples/SAMPLE_GUIDELINES.md:9-42`). Gaps: (a) no in-repo generated API reference exists — canonical reference documentation lives off-repo on learn.microsoft.com (`README.md:96-100`, `dotnet/nuget/nuget-package.props:45`), so docstring/XML-doc quality is the only in-repo reference; (b) `python/AGENTS.md:10-22` references seven contributor skill files (e.g., `python-feature-lifecycle`, the exact doc explaining stage promotion) that do not exist — `.github/skills/` contains only `pull-requests` — so the documented process for how APIs move between stages is partially missing from the repo; (c) advanced staged surfaces like the harness have samples (`dotnet/samples/02-agents/Harness/`, `python/samples` harness console referenced in `python/CHANGELOG.md:35`) but the Python-side experimental warnings mean example code emits warnings by design.

## Architectural Decisions

- **Abstractions/concrete split on .NET**: `Microsoft.Agents.AI.Abstractions` carries only interfaces and base classes and depends solely on `Microsoft.Extensions.AI.Abstractions` (`dotnet/src/Microsoft.Agents.AI.Abstractions/Microsoft.Agents.AI.Abstractions.csproj:27-29`), letting host libraries depend on contracts while provider packages depend on implementations.
- **One import path, many packages (Python)**: the core package vendors provider namespaces (`openai/`, `azure/`, `anthropic/`, ...) that lazily dispatch to separately installed distributions (`python/packages/core/agent_framework/foundry/__init__.py:6-30`), keeping the advertised import surface stable while implementations ship independently.
- **Stage metadata as first-class mechanism, not comments**: the feature-stage system combines enum inventories, docstring injection, runtime warnings with user-frame attribution, and once-per-feature deduplication (`python/packages/core/agent_framework/_feature_stage.py:161-163`, `192-251`, `233-235`) — a deliberate decision to make instability observable.
- **Diagnostic-ID unification**: .NET experimental APIs reuse MEAI's `MEAI001` and OpenAI's `OPENAI001` IDs so consumers suppress one ID per ecosystem rather than one per library (`dotnet/src/Shared/DiagnosticIds/DiagnosticsIds.cs:18-31`).
- **Governance by ADR**: naming and shape changes route through numbered, dated, attributed decision records (`docs/decisions/0005-python-naming-conventions.md:1-13`), including explicit rejected alternatives.

## Notable Patterns

- **Curated `__all__` as the single source of truth** for the Python surface (`python/packages/core/agent_framework/__init__.py:342-618`), grouped by import origin so reviewers can see exactly which module each export comes from (`__init__.py:20-340`).
- **Lifecycle encoded in build metadata**, not prose: `IsReleased`/`IsReleaseCandidate` flags drive version suffixing (`-rc1`, `-preview.YYYYMMDD.1`) and enable package validation only for GA packages (`dotnet/nuget/nuget-package.props:5-11`, `22-24`).
- **Helpful failure for missing optional deps**: lazy namespaces catch `ModuleNotFoundError` and re-raise with "The package X is required to use Y. Install it with: pip install X" (`python/packages/core/agent_framework/openai/__init__.py:34-40`).
- **Capability protocols over inheritance**: `Supports*` protocols (`python/packages/core/agent_framework/__init__.py:21-32`) let clients declare optional capabilities (shell, web search, embeddings) that the framework probes instead of forcing deep type hierarchies.
- **Fluent builders at every composition edge**: `AIAgentBuilder.Use(...)` (`dotnet/src/Microsoft.Agents.AI/AIAgentBuilder.cs:76-112`), `WorkflowBuilder` with output selection (`python/packages/core/agent_framework/__init__.py:325`), orchestrator builders (`SequentialBuilder`...`MagenticBuilder` documented at `python/packages/orchestrations/agent_framework_orchestrations/__init__.py:6-11`).
- **Sample-as-documentation with XML tag anchors**: get-started files embed `<create_agent>` style tags used by the external docs repo for code reuse (`python/samples/01-get-started/01_hello_agent.py:15,20`), keeping repo samples and website docs in sync structurally.

## Tradeoffs

- **Centralized export list vs distributed ownership**: the 618-line `__init__.py` gives a crisp contract but concentrates merge conflicts and makes it easy for a submodule change to miss the export update; there is no test asserting `__all__` matches module contents (none found in `python/packages/core/tests/` searches).
- **Convention-based privacy (Python)**: zero-cost and idiomatic, but unlike the .NET `internal` keyword it cannot stop determined imports; combined with missing per-module `__all__`, star-imports leak everything.
- **Dual-language parity**: ADR-0005 commits to preserving ".NET name mapping" (`docs/decisions/0005-python-naming-conventions.md:10-12`), yet packaging strategies already diverge (harness: separate RC .NET package with `[Experimental]` members vs Python in-core modules gated by `HARNESS` feature id, `python/packages/core/agent_framework/__init__.py:85-147` + `dotnet/src/Microsoft.Agents.AI.Harness/HarnessAgent.cs:83`), which complicates cross-language documentation.
- **Warning-once strategy**: deduplicating warnings per `(category, feature_id)` (`python/packages/core/agent_framework/_feature_stage.py:233-235`) keeps logs clean but hides repeated usage sites after the first.
- **Off-repo canonical docs**: linking learn.microsoft.com as the reference (`README.md:96`) avoids doc duplication but means the repo alone does not satisfy "use the API correctly without reading internals" for reference lookup.

## Failure Modes / Edge Cases

- **Accidental public surface growth**: any new symbol added to an underscore module and imported into `__init__.py` becomes de-facto public instantly; the compat baseline protects .NET packages (`dotnet/nuget/nuget-package.props:19-24`) but Python packages rely on release discipline and the changelog's BREAKING tags (`python/CHANGELOG.md:22-24`) — no equivalent automated API-diff gate was found for Python.
- **Stale contributor documentation**: `python/AGENTS.md:10-22` advertises skill files (including the feature-lifecycle guide) that are absent from `.github/skills/` (verified: directory contains only `pull-requests`), risking inconsistent application of stage rules.
- **Version skew between language stacks**: Python workspace pins `version = "1.9.0"` (`python/pyproject.toml:7`) while the .NET central version prefix is `1.10.0` (`dotnet/nuget/nuget-package.props:4`); consumers operating both stacks cannot assume feature alignment at equal version numbers.
- **Deprecated package traps**: renamed packages (e.g., `azure-ai` → `foundry`) are recorded in `python/PACKAGE_STATUS.md:48-53`, but the old import path's behavior post-deprecation (error vs silent absence) is only discoverable by reading this table.
- **Protocol-class decoration hazards**: the feature-stage machinery special-cases Protocol classes to avoid corrupting `runtime_checkable` semantics and skips runtime warnings for them (`python/packages/core/agent_framework/_feature_stage.py:357-365`, `371-374`) — meaning experimental protocols give docstring notice but no runtime warning, an inconsistency consumers might not expect.

## Future Considerations

- Add an automated Python API-surface check (e.g., snapshot of `dir(agent_framework)` vs `__all__`, or per-module `__all__` generation) to mirror the .NET package-validation guarantee.
- Generate or vendor a minimal in-repo API reference (docstring-derived) so staged/experimental annotations are greppable alongside code.
- Restore or remove the missing `.github/skills/python-*` references in `python/AGENTS.md` to keep the documented lifecycle process authoritative.
- Consider promoting the harness out of experimental (both stacks currently mark it: `MAAI001` on .NET, `HARNESS` id on Python) now that it spans approval, memory, file-access, and loop subsystems with tests (`dotnet/tests/Microsoft.Agents.AI.Harness.UnitTests/`).
- Track cross-language parity explicitly (a manifest mapping .NET type ↔ Python symbol) to catch drift earlier than ADR reviews.

## Questions / Gaps

- No evidence found of an automated public-API diff gate for the Python packages (searched `python/scripts/`, `shared_tasks.toml`, and CI workflow names under `.github/actions` for api-diff/api-check tasks; .NET has `EnablePackageValidation`, Python counterpart not located). If it exists, it lives outside this source snapshot.
- No evidence found of generated reference documentation in-repo (no `api/` or docgen outputs under `docs/`); canonical reference is external (`README.md:96-100`).
- The exact promotion criteria between `alpha → beta → rc → released` are referenced (`python/PACKAGE_STATUS.md:5-9`, pointer to a `python-feature-lifecycle` skill at `python/AGENTS.md:13-14`) but the criteria document itself is absent from this snapshot, so the operational definition of "stable" could not be verified beyond the presence of the buckets.

---

Generated by `Dimension 24.01: Public API Surface` against `agent-framework`.
