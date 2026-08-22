# Source Analysis: agent-framework

## Dimension 22.01 — Package and Module Boundaries

### Source Info

| Field | Value |
|-------|-------|
| Name | agent-framework (Microsoft Agent Framework) |
| Path | `studies/agent-harness-study/sources/agent-framework` |
| Language / Stack | Multi-language monorepo: C#/.NET (`dotnet/`, 36 projects) and Python (`python/`, 31 uv-workspace packages) |
| Analyzed | 2026-08-21 |

## Summary

Microsoft Agent Framework is a two-stack monorepo (`.NET` + `Python`) with an unusually disciplined package-boundary story. Both stacks converge on the same layered model:

- **Abstractions/core at the bottom with near-zero dependencies.** `Microsoft.Agents.AI.Abstractions` references only `Microsoft.Extensions.AI.Abstractions` (`dotnet/src/Microsoft.Agents.AI.Abstractions/Microsoft.Agents.AI.Abstractions.csproj:33`); Python `agent-framework-core` depends only on `typing-extensions`, `pydantic`, `python-dotenv`, and `opentelemetry-api` (`python/packages/core/pyproject.toml:25-30`).
- **Implementations in the middle.** `Microsoft.Agents.AI` → Abstractions; provider connectors (OpenAI, Anthropic, Foundry, MCP, …) → `Microsoft.Agents.AI`; Workflows sits alongside as a separate assembly. The dependency graph was extracted from all 36 `.csproj` files and verified cycle-free by DFS (see Evidence).
- **Hosting/UI on top.** `Microsoft.Agents.AI.Hosting*` → Hosting → Workflows+AI; DevUI packages sit above Hosting.
- **Everything ships as independently versioned packages**, coordinated by central package management on .NET (`dotnet/Directory.Packages.props:4-5`) and a uv workspace on Python (`python/pyproject.toml`, `[tool.uv.workspace] members = ["packages/*"]`), with a meta-package (`agent-framework`, `agent-framework-core[all]`) for users who want the world.

The boundary model is not just convention: it is documented in an accepted ADR (`docs/decisions/0008-python-subpackages.md`), enforced by a dedicated dependency-bounds validation tool that runs each package's test+typecheck suites at both lowest and highest allowed resolutions in isolated environments (`python/scripts/dependencies/validate_dependency_bounds.py:243-244`), guarded by lazy-loading facades so import paths stay stable while code moves between packages (`python/packages/core/agent_framework/azure/__init__.py:11-33`), and tested — there are dedicated tests asserting the facade re-exports are identical to the real package objects (`python/packages/core/tests/core/test_azure_namespace.py:7-9`).

Weaknesses are specific rather than structural: one production `InternalsVisibleTo` across assemblies (Workflows → DurableTask), cross-package imports of *private* Python modules (foundry reaching into `agent_framework_openai._chat_client`), compile-time source sharing via `dotnet/src/Shared` that bypasses package boundaries, and a pinned pre-release cross-dependency in devui's dev extra.

## Rating

**8 / 10** — Clear model with explicit interfaces, operational safeguards, and separation tests.

Rationale against the rubric:
- **Clear model**: layering is uniform across both stacks and codified in ADR 0008 (`docs/decisions/0008-python-subpackages.md:70-84`); packaging defaults are deny-by-default on .NET (`<IsPackable>false</IsPackable>` in `dotnet/Directory.Build.props`).
- **Explicit interfaces**: public API surface is declared (`__all__` at `python/packages/core/agent_framework/__init__.py:342`; `Abstractions` package on .NET); experimental APIs carry machine-readable stage markers (`[Experimental("MAAI001")]` at `dotnet/src/Microsoft.Agents.AI.Abstractions/AIContextProvider.cs:381`; `ExperimentalFeature` enum at `python/packages/core/agent_framework/_feature_stage.py:55-70`).
- **Tests/safeguards**: cycle-free graph verified by this analysis; dependency-bound gates run full test suites per package at lower/upper resolutions (`python/scripts/dependencies/validate_dependency_bounds.py:213-287`); lazy facades have dedicated unit tests.
- Not 9–10 because: private-module imports leak across package boundaries on the Python side, the .NET `Shared` source-injection mechanism creates invisible compile-time coupling between packages, and one production-grade `InternalsVisibleTo` weakens the Workflows/DurableTask seam.

## Evidence Collected

| Area | Evidence | File:Line |
|------|----------|-----------|
| Top-level structure | Monorepo split: `dotnet/src` (36 projects), `python/packages` (31 packages), shared `docs/decisions`, `schemas`, per-language samples | `dotnet/src`, `python/packages` (directory listing) |
| Python workspace | uv workspace with `members = ["packages/*"]`; all internal packages resolved via `[tool.uv.sources] workspace = true` | `python/pyproject.toml:95-96`, `python/pyproject.toml:118-152` |
| Meta package | `agent-framework` depends solely on `agent-framework-core[all]==1.9.0` | `python/pyproject.toml:14-16` |
| Core minimal deps | Core runtime deps limited to typing-extensions, pydantic, python-dotenv, opentelemetry-api | `python/packages/core/pyproject.toml:25-30` |
| Core `all` extra | All provider packages enumerated as optional deps of core's `all` extra | `python/packages/core/pyproject.toml:32-59` |
| Provider → core direction | openai/anthropic/a2a/orchestrations/etc. all declare `agent-framework-core>=X,<2` | `python/packages/openai/pyproject.toml:26`, `python/packages/anthropic/pyproject.toml:26`, `python/packages/a2a/pyproject.toml:25` |
| Second-tier deps | foundry→core+openai; azurefunctions→core+durabletask; azure-contentunderstanding→core+foundry | `python/packages/foundry/pyproject.toml:26-27`, `python/packages/azurefunctions/pyproject.toml:25-26`, `python/packages/azure-contentunderstanding/pyproject.toml:26-27` |
| Circular-dep avoidance | Comment explains tools→core prevents core listing it at runtime; declared as dev-only group instead | `python/packages/core/pyproject.toml:61-67` |
| Dependency graph tooling | `_build_internal_graph` maps every package's internal `agent-framework-*` deps into a graph | `python/scripts/dependencies/_dependency_bounds_upper_impl.py:551-574` |
| Boundary CI gate | Test mode runs `test`+`pyright` per package at `lowest-direct` and `highest` resolutions in isolated uv envs; hard-fails workspace | `python/scripts/dependencies/validate_dependency_bounds.py:243-244`, `python/scripts/dependencies/validate_dependency_bounds.py:152-210` |
| Transitive closure install | `_resolve_internal_editables` walks graph so a package under test gets its internal deps as editables | `python/scripts/dependencies/_dependency_bounds_upper_impl.py:577-593` |
| Public API surface | Module docstring "Public API surface"; explicit `__all__` export list | `python/packages/core/agent_framework/__init__.py:2-7`, `python/packages/core/agent_framework/__init__.py:342` |
| Private implementation modules | Implementation in underscore modules (`_agents.py`, `_clients.py`, `_tools.py`, …) re-exported through the facade; internal cross-imports annotated `# pyright: ignore[reportPrivateUsage]` | `python/packages/core/agent_framework/` (listing), `python/packages/core/agent_framework/_agents.py:28` |
| Lazy provider facades | `azure/__init__.py` `_IMPORTS` map + module `__getattr__` raising "package X is required… pip install X" | `python/packages/core/agent_framework/azure/__init__.py:11-33` |
| Facade type stubs | `openai/__init__.pyi` stubs accompany the lazy re-export shim inside core | `python/packages/core/agent_framework/openai/__init__.pyi` |
| Separation tests (Python) | Tests assert facade attribute identity with real package exports | `python/packages/core/tests/core/test_azure_namespace.py:7-9` (also `test_foundry_namespace.py`, `test_hyperlight_namespace.py`) |
| Optional heavy deps | `mcp` imported only under TYPE_CHECKING; deferred runtime imports keep core importable without mcp | `python/packages/core/agent_framework/_mcp.py:44-49`, `python/packages/core/agent_framework/_mcp.py:301` |
| Documented boundary policy (ADR) | Option 7 accepted: subpackages exist to keep non-GA code/deps out of main; stable import paths regardless of which package hosts code; lazy loading with meaningful errors | `docs/decisions/0008-python-subpackages.md:70-84` |
| Feature-stage visibility | `ExperimentalFeature`/`ReleaseCandidateFeature` enums drive warnings + docstring injection for staged APIs | `python/packages/core/agent_framework/_feature_stage.py:55-70`, `python/packages/core/agent_framework/_feature_stage.py:23-38` |
| Lab isolation | Experimental `lab` package excluded from root test sweep (`norecursedirs = '**/lab/**'`) | `python/pyproject.toml:173` |
| .NET project graph | Full ProjectReference extraction: Abstractions×17 consumers, AI×12, Workflows×5, Hosting×5; DFS found zero cycles across 36 projects | `dotnet/src/**/*.csproj` (analysis run, see Answers Q2) |
| Abstractions purity | Only external ref is `Microsoft.Extensions.AI.Abstractions` PackageReference | `dotnet/src/Microsoft.Agents.AI.Abstractions/Microsoft.Agents.AI.Abstractions.csproj:31-34` |
| Central package management | `ManagePackageVersionsCentrally=true` + transitive pinning; 135 pinned versions | `dotnet/Directory.Packages.props:4-6` |
| Packaging opt-in | `<IsPackable>false</IsPackable>` repo-wide default; release filter lists exactly the shippable projects | `dotnet/Directory.Build.props` (IsPackable block), `dotnet/agent-framework-release.slnf` (35 projects) |
| InternalsVisibleTo policy | Grants overwhelmingly to `*.UnitTests` assemblies (e.g. Abstractions→its UnitTests) | `dotnet/src/Microsoft.Agents.AI.Abstractions/Microsoft.Agents.AI.Abstractions.csproj:34`, `dotnet/src/Microsoft.Agents.AI/Microsoft.Agents.AI.csproj:50-54` |
| Production IVT exception | Workflows grants internals to DurableTask (non-test) | `dotnet/src/Microsoft.Agents.AI.Workflows/Microsoft.Agents.AI.Workflows.csproj:29` |
| Shared-source injection | `dotnet/src/Shared/*` compiled into consumer assemblies via MSBuild flags (`InjectSharedThrow`, `InjectSharedWorkflowsExecution`, …) | `dotnet/eng/MSBuild/Shared.props:4-38`, e.g. `dotnet/src/Microsoft.Agents.AI.Abstractions/Microsoft.Agents.AI.Abstractions.csproj:10-19` |
| Legacy polyfill injection | Polyfill attributes injected only for TFMs lacking them | `dotnet/eng/MSBuild/LegacySupport.props:14-15` |
| Experimental markers (.NET) | `[Experimental(DiagnosticIds.Experiments.AgentsAIExperiments)]` on public members; IDs centralized, deliberately reusing MEAI001/OPENAI001 so consumers suppress once | `dotnet/src/Microsoft.Agents.AI.Abstractions/AIContextProvider.cs:381`, `dotnet/src/Shared/DiagnosticIds/DiagnosticsIds.cs:16-27` |
| Test mirroring | One UnitTests (+ IntegrationTests where relevant) project per src package | `dotnet/tests/` (directory listing) |
| Samples consume via refs | Sample projects reference src projects directly (`ProjectReference` to Foundry) — samples excluded from shipped surface | `dotnet/samples/01-get-started/01_hello_agent/01_hello_agent.csproj:17` |

## Answers to Dimension Questions

### 1. Are modules cleanly separated?

Yes, with named exceptions. Both stacks use the same three-tier shape (abstractions/core → implementations → hosting/UI). On .NET, `Abstractions` has no framework-internal references at all (`dotnet/src/Microsoft.Agents.AI.Abstractions/Microsoft.Agents.AI.Abstractions.csproj:33`). On Python, core's only runtime dependencies are four small libraries (`python/packages/core/pyproject.toml:25-30`), and even OpenAI support — historically built into core — now lives behind a lazy facade delegating to the separate `agent-framework-openai` distribution (`python/packages/core/agent_framework/openai/__init__.py:11-24`). Three leaks prevent a perfect score: the `InternalsVisibleTo` from Workflows to DurableTask (`dotnet/src/Microsoft.Agents.AI.Workflows/Microsoft.Agents.AI.Workflows.csproj:29`), foundry importing openai's private `_chat_client` module (`python/packages/foundry/agent_framework_foundry/_agent.py:34`), and MSBuild source-sharing that compiles `src/Shared/Workflows/*` directly into other assemblies (`dotnet/eng/MSBuild/Shared.props:16-22`).

### 2. Do dependencies flow in one direction?

Yes. This analysis extracted every `ProjectReference` from the 36 `.csproj` files under `dotnet/src` and ran a DFS cycle check: **zero cycles**. The layering is strictly: Abstractions/Generators (no refs) → `Microsoft.Agents.AI` → providers & Workflows → Hosting → Hosting.* / DevUI (e.g., `Microsoft.Agents.AI.Hosting.csproj` references Abstractions, Workflows, AI; `Hosting.AspNetCore` references only Hosting). On Python, every provider package declares exactly `agent-framework-core` plus its vendor SDK (`python/packages/openai/pyproject.toml:26-27`); the only reverse edge would be core↔tools, which is explicitly prevented by making tools a dev-only dependency-group with a documenting comment (`python/packages/core/pyproject.toml:61-67`). Downward-only is additionally policed by the dependency-bounds gate, which builds the internal graph and fails the workspace if any package's suite breaks under strict resolution extremes (`python/scripts/dependencies/validate_dependency_bounds.py:87`, `:130-133`).

### 3. Can modules be used independently?

Yes — this is the strongest dimension answer, and it holds up under the study question *"Can you use the tool system without pulling in the entire runtime?"*. Concretely:
- Python: `pip install agent-framework-core` pulls only pydantic/dotenv/OTel-api; `@tool`/`FunctionTool` live in core (`python/packages/core/agent_framework/__init__.py` exports `tool`), and MCP support degrades gracefully because `mcp` is imported lazily (`python/packages/core/agent_framework/_mcp.py:301`). The tool system therefore works without any provider or the meta-package. The `agent-framework` meta-package exists purely for opt-in bulk installs (`python/pyproject.toml:14-16`).
- .NET: `Microsoft.Agents.AI.Abstractions` + `Tools.Shell` (which references only Abstractions, `dotnet/src/Microsoft.Agents.AI.Tools.Shell/Microsoft.Agents.AI.Tools.Shell.csproj`) form a minimal tooling footprint without Workflows, Hosting, or any connector.
- Independent usability is *tested*: the bounds validator runs each package's test+pyright suite in an isolated uv environment containing only that package plus its internal closure (`--isolated`, `--with-editable` of graph-resolved internals; `python/scripts/dependencies/validate_dependency_bounds.py:153-175`, `python/scripts/dependencies/_dependency_bounds_upper_impl.py:577-593`).

### 4. Are public APIs distinguished from internal ones?

Yes, via multiple complementary mechanisms:
- Python: single facade `__init__` with explicit `__all__` over underscore-prefixed implementation modules (`python/packages/core/agent_framework/__init__.py:342`; `_agents.py:28` shows internal usage is lint-flagged via `reportPrivateUsage`).
- Python lifecycle stages: `@experimental`/release-candidate decorators attach warning categories and docstring banners, with feature IDs inventoried in enums (`python/packages/core/agent_framework/_feature_stage.py:55-70`).
- .NET: the `...Abstractions` package name itself is the public-contract marker (mirroring Microsoft.Extensions.* conventions); unstable APIs carry `[Experimental("MAAI001")]` with centralized diagnostic IDs intentionally aligned with MEAI/OpenAI IDs to reduce consumer suppression burden (`dotnet/src/Microsoft.Agents.AI.Abstractions/AIContextProvider.cs:381`, `dotnet/src/Shared/DiagnosticIds/DiagnosticsIds.cs:18-27`).
- .NET visibility: `internal` + narrowly scoped `InternalsVisibleTo`, almost exclusively to test assemblies (e.g., `dotnet/src/Microsoft.Agents.AI/Microsoft.Agents.AI.csproj:50-54` grants only to `DynamicProxyGenAssembly2` and four `*.UnitTests`).
- Ship-surface control: default `IsPackable=false` plus a release solution filter enumerating exactly what publishes (`dotnet/Directory.Build.props`, `dotnet/agent-framework-release.slnf`).

## Architectural Decisions

1. **ADR-governed package topology** (`docs/decisions/0008-python-subpackages.md:70-84`). The team evaluated seven options and chose "subpackage existence based on dependency maturity": anything with non-GA deps, preview code, or fast-moving third parties stays a separate distribution; graduation into the main package must preserve import paths (`from agent_framework.google import ...` keeps working). Import-path stability is implemented literally by the facade folders inside core.
2. **Meta-package aggregation instead of a monolith** — `agent-framework` is a shell over `core[all]` (`python/pyproject.toml:14-16`), keeping the default install lean while preserving one-command onboarding.
3. **Abstractions-first .NET layout** — a dependency-free contract assembly consumed by 17 projects forces providers and hosting to program against interfaces, enabling independent replacement of clients/hosting.
4. **Compile-time source sharing over package refs for tiny utilities** — `Throw`, `Redaction`, diagnostic-ID helpers are injected into each assembly via MSBuild flags (`dotnet/eng/MSBuild/Shared.props:4-14`) rather than creating a shared runtime dependency. This avoids an extra published package at the cost of implicit coupling (see Tradeoffs).
5. **Dev-dependency inversion to break cycles** — core lists `agent-framework-tools` under `[dependency-groups].dev` with an inline rationale comment (`python/packages/core/pyproject.toml:61-67`), so harness integration can be typed/tested without a circular runtime edge.
6. **Centralized version governance** — CPM with transitive pinning on .NET (`dotnet/Directory.Packages.props:4-5`) and a frozen uv lockfile with security-floor overrides on Python (`python/pyproject.toml`, `constraint-dependencies`/`override-dependencies` blocks).

## Notable Patterns

- **Lazy namespace facades with actionable errors**: `_IMPORTS: dict[str, tuple[str, str]]` + module-level `__getattr__` that raises `ModuleNotFoundError` naming the exact pip package to install (`python/packages/core/agent_framework/azure/__init__.py:15-33`). The same pattern backs `openai/`, `anthropic/`, `google/`, etc. inside core, each paired with `.pyi` stubs so static typing still sees the surface.
- **Boundary identity tests**: `assert azure.CosmosHistoryProvider is CosmosHistoryProvider` proves the facade delegates to the real distribution object rather than a copy (`python/packages/core/tests/core/test_azure_namespace.py:7-9`).
- **Graph-driven tooling**: the dependency tooling parses all workspace `pyproject.toml`s into an internal graph and reuses it for both editable-closure installation and bound validation (`python/scripts/dependencies/_dependency_bounds_upper_impl.py:551-593`) — boundaries are data, not prose.
- **Dual-end resolution testing**: running suites at both `lowest-direct` and `highest` resolutions catches both stale lower bounds and future upper-bound breaks (`python/scripts/dependencies/validate_dependency_bounds.py:241-244`).
- **Diagnostic-ID alignment**: experimental APIs reuse ecosystem IDs (MEAI001, OPENAI001) so consumers already suppressing those diagnostics don't need new entries (`dotnet/src/Shared/DiagnosticIds/DiagnosticsIds.cs:20-27`).
- **Polyfill injection by TFM**: legacy-targeting projects get `ExperimentalAttribute`/`RequiredMemberAttribute` sources injected only when the target framework lacks them (`dotnet/eng/MSBuild/LegacySupport.props:14-15`), keeping multi-TFM support without shipping polyfill types publicly.

## Tradeoffs

- **Facade maintenance cost vs. import stability**: every provider needs an `__init__.py` shim + `.pyi` stubs inside core; the ADR explicitly accepts "larger overhead in maintaining the `__init__.py` files" (`docs/decisions/0008-python-subpackages.md:69-71`). Moving a connector between packages requires synchronized edits in two places.
- **Source sharing vs. visible dependencies**: `Shared.props` injection keeps small helpers dependency-free but hides coupling — nothing in a `.csproj` reference graph reveals that an assembly embeds `Shared/Workflows/Execution` sources; version skew between copies is possible since each assembly compiles its own snapshot.
- **Lean core vs. discoverability**: users must know which extra provides `FoundryChatClient`; mitigated by the error-message pattern, but a wrong-import now fails at first attribute access rather than at import time.
- **`InternalsVisibleTo` to DurableTask**: pragmatic for durable-execution integration deep in workflow state, but it makes the Workflows internal surface de-facto load-bearing for another shipped package, narrowing future refactoring freedom.
- **Strict upper bounds (`<2`) everywhere**: protects against major-version drift but means every internal package must be released in lockstep when core hits 2.0; the bounds validator exists precisely to manage this tax.

## Failure Modes / Edge Cases

- **Cross-package private-module import**: foundry imports `agent_framework_openai._chat_client` directly (`python/packages/foundry/agent_framework_foundry/_agent.py:34`, `_chat_client.py:22`). Any refactor of openai's private module breaks foundry at import time despite semver compatibility — the exact failure mode underscore-module privacy is meant to prevent. No automated guard (e.g., lint rule banning cross-package `_module` imports) was found in the ruff config (`python/pyproject.toml` lint section selects standard rule sets only).
- **Version pin skew in extras**: devui's dev optional-dependency pins `agent-framework-orchestrations==1.0.0rc1` while other packages float `>=X,<2` (`python/packages/devui/pyproject.toml:37`) — a pre-release pin inside a floating range family that can produce unsatisfiable resolutions during lockstep releases.
- **Facade drift**: if a symbol is added to a provider package but not to its core facade map, `import agent_framework.azure` silently misses it; detection relies on the namespace tests being updated by hand (`python/packages/core/tests/core/test_azure_namespace.py` covers one symbol per pattern shown).
- **Hidden shared-source divergence**: because `Shared/Workflows/*` is compiled per-assembly, a bug fix applied to one consuming project's tree but not the canonical folder produces divergent behavior with no linker error.
- **Bounds-gate blind spots**: the validator treats any package lacking `test`/`pyright` poe tasks as a hard configuration error (`python/scripts/dependencies/validate_dependency_bounds.py:107-112`), which is good, but the gate validates *declared* bounds — undeclared imports (like the foundry→openai private import) would only surface if openai were absent from foundry's resolved environment.

## Future Considerations

- Add a lint/CI rule forbidding cross-distribution imports of underscore modules on the Python side (would have caught `foundry→agent_framework_openai._chat_client`), and promote the needed symbols to openai's public surface.
- Replace high-traffic `Shared` source injection (especially `Shared/Workflows/*`) with an internal shared library or generator-managed linkage so the coupling appears in the reference graph.
- Re-evaluate the Workflows→DurableTask `InternalsVisibleTo`: extract the required types into `Abstractions` or a `Workflows.Internal` contract package to restore refactor freedom.
- Automate facade-map synchronization (generate `_IMPORTS` dicts and `.pyi` stubs from provider package exports) to remove the hand-maintained shim failure mode.
- Normalize cross-package constraint style (float ranges vs. exact pins like `devui`'s orchestrations rc pin) before the next coordinated release.

## Questions / Gaps

- **No dedicated architecture/boundary unit tests on .NET** were found: searches for tests asserting reference constraints or Abstractions purity returned no architectural-test fixtures; enforcement there relies on conventions (package naming, `IsPackable`, analyzers) rather than executable checks. Searched: `dotnet/tests/**` for architecture-style assertions, `grep -r "InternalsVisibleTo"` across props/csproj.
- **CI wiring of the dependency-bounds gate** could not be confirmed from files in the source snapshot (workflow definitions live under `.github/workflows/`, which contains skill/pull-request templates in this checkout; no workflow file referencing `validate_dependency_bounds` was located). The task definitions exist (`python/pyproject.toml:268-296`), but evidence that they run on every PR is absent.
- The referenced `.github/skills/python-package-management` guidance cited by `python/AGENTS.md:10` is not present in this snapshot (only `.github/skills/pull-requests` exists), so the documented monorepo conventions could not be verified beyond ADR 0008.
- Whether the `lab` package is publishable was ambiguous: it participates in core's `all` extra (`python/packages/core/pyproject.toml:52`) yet is excluded from the root test sweep (`python/pyproject.toml:173`); its release status is not stated in the inspected files.

---

Generated by dimension `22.01-package-and-module-boundaries` against `agent-framework`.
