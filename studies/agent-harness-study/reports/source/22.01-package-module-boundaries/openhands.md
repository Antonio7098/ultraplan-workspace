# Source Analysis: openhands

## Dimension 22.01: Package and Module Boundaries

### Source Info

| Field | Value |
|-------|-------|
| Name | openhands |
| Path | `studies/agent-harness-study/sources/openhands` |
| Language / Stack | Python 3.12 (Poetry/uv, FastAPI), TypeScript/React (Vite), monorepo |
| Analyzed | 2026-08-22 |

## Summary

OpenHands is a multi-package monorepo in the middle of a deliberate V1 re-architecture that splits agent concerns across four Python distributions — `openhands-sdk`, `openhands-tools`, `openhands-agent-server`, and the root metapackage `openhands-ai` (`pyproject.toml:60-62`, pinned `==1.29.0` at `pyproject.toml:249-251`) — plus a React app (`frontend/`), a publishable UI kit (`openhands-ui/`, npm package `@openhands/ui`), and an enterprise overlay (`enterprise/`) that depends one-way on the OSS root (`enterprise/pyproject.toml:29`). The runtime/tool/memory layers are no longer in this repo at all: the agent loop lives in the external SDK, and this repo's `openhands/` package is a thin application layer containing `analytics/`, `app_server/`, `db/`, and legacy shims under `server/`.

Boundaries are expressed through three mechanisms: (1) a namespace-package merge so all four distributions share one top-level `openhands` package (`openhands/__init__.py:4`); (2) a dependency-injection composition root with an `Injector` abstraction (`openhands/app_server/services/injector.py:13-40`) and env-driven implementation selection via `get_impl()` (`openhands/app_server/utils/import_utils.py:43-78`); (3) deprecated re-export shim modules for the V0→V1 migration (`openhands/server/listen.py:1-8`). However, inside `app_server/` there are at least four subpackage-level import cycles worked around by ~55 function-scoped deferred imports; cross-boundary hooks to enterprise are stringly-typed (`OPENHANDS_CONFIG_CLS`); analytics still imports a deprecated module; and the documented frontend layering rule has 42 component-level violations with no lint enforcement. Boundaries are well modeled but weakly enforced.

## Rating

**Score: 6 / 10**

Rationale against the rubric:

- **Why above 4–5**: The module model is explicit and architectural, not ad-hoc. The SDK/tools/agent-server split into separately versioned distributions directly answers "can you use the tool system without pulling in the entire runtime" at the packaging level; extension points are real code mechanisms (`get_impl`, injectors), not just conventions; deprecation shims show managed boundary migration; the UI kit exposes a curated public surface (`openhands-ui/index.ts:1-18`).
- **Why below 7–8**: The rubric's 7–8 tier requires tests and operational safeguards for the boundary model. No automated import-cycle or dependency-direction checks exist (no import-linter config found in `pyproject.toml` or `dev_config/`; searches for `import-linter|lint-imports|importlinter` returned nothing). Four subpackage cycles exist inside `app_server/`, patched by deferred imports rather than refactoring. Cross-package hooks are stringly-typed and unvalidated until runtime. Frontend layering is convention-only with dozens of standing violations.

## Evidence Collected

All paths are relative to `studies/agent-harness-study/sources/openhands/`.

| Area | Evidence | File:Line |
|------|----------|-----------|
| Package structure | Root metapackage `openhands-ai` with `[tool.poetry] packages = [{ include = "openhands/**/*" }]` and `skills/**/*` include | `pyproject.toml:141-156` |
| Runtime/tools separation | V1 dependencies `openhands-agent-server==1.29.0`, `openhands-sdk==1.29.0`, `openhands-tools==1.29.0` declared as external packages | `pyproject.toml:60-62`, `pyproject.toml:249-251` |
| Namespace merging | `__path__ = __import__('pkgutil').extend_path(__path__, __name__)` so sdk/tools/agent-server share the `openhands` top-level package | `openhands/__init__.py:1-9` |
| Thin app layer | Remaining modules: `analytics/`, `app_server/`, `db/`, `server/` (218 .py files total) | `openhands/app_server/__init__.py`, directory listing |
| Composition root | `AppServerConfig` pydantic model holds all service injectors; `config_from_env()` wires defaults from env | `openhands/app_server/config.py:191-237`, `openhands/app_server/config.py:240-420` |
| DI abstraction | `Injector(Generic[T], ABC)` with `inject/context/depends` used by FastAPI | `openhands/app_server/services/injector.py:13-40` |
| Dynamic impl loading | `get_impl(cls, impl_name)` imports and validates named subclass; cached via `lru_cache` | `openhands/app_server/utils/import_utils.py:34-78` |
| Router aggregation | `v1_router` mounts 13 sub-routers under `/api/v1`; app includes it plus health + middleware | `openhands/app_server/v1_router.py:24-37`, `openhands/app_server/app.py:54-86` |
| Enterprise → OSS direction | `openhands-ai = { path = "../", develop = true }`; enterprise ships top-level packages `server`, `storage`, `sync`, `integrations` | `enterprise/pyproject.toml:23-29` |
| Enterprise entry extends base app | `from openhands.app_server.app import app as base_app`; SaaS routes added on top; `SERVE_FRONTEND=false` | `enterprise/saas_server.py:7-11`, `enterprise/saas_server.py:67-71` |
| OSS → enterprise hook (stringly-typed) | `_get_default_lifespan()` lazily imports `server.app_lifespan.saas_app_lifespan_service.SaasAppLifespanService` when `OPENHANDS_CONFIG_CLS` contains 'saas' | `openhands/app_server/config.py:174-183` |
| OSS → enterprise hook (store classes) | `SettingsStoreImpl = get_impl(SettingsStore, server_config.settings_store_class)` resolved at import time of `shared.py` | `openhands/app_server/shared.py:16-27` |
| Cycle: event ↔ app_conversation | `event_service_base` imports `app_conversation_info_service`/models; reverse import of `EventService` in live-status service | `openhands/app_server/event/event_service_base.py:11-14`; `openhands/app_server/app_conversation/live_status_app_conversation_service.py:72-73` |
| Cycle: sandbox ↔ app_conversation | remote sandbox service imports conversation models; conversation models import `SandboxStatus` | `openhands/app_server/sandbox/remote_sandbox_service.py:24`; `openhands/app_server/app_conversation/app_conversation_models.py:20` |
| Cycle: user_auth ↔ settings | `user_auth/__init__` imports `settings_models`/`settings_store`; settings router imports `user_auth` | `openhands/app_server/user_auth/__init__.py:8-9`; `openhands/app_server/settings/settings_router.py:34` |
| Cycle: analytics ↔ app_server | `analytics_context` imports `app_server.utils` (and lazily `app_server.shared`); routers import `openhands.analytics` | `openhands/analytics/analytics_context.py:21-22`, `openhands/analytics/analytics_context.py:58`; `openhands/app_server/settings/settings_router.py:15`; `openhands/app_server/app_conversation/app_conversation_router.py:21` |
| Cycle mitigation pattern | 55 function-body-level `from openhands...` imports inside `app_server` defer resolution past import time | grep count over `openhands/app_server/**/*.py` (e.g., `openhands/app_server/config.py:242-283`, `secrets_router.py:137`) |
| Config fan-in hub | 33 files import `openhands.app_server.config`, making it the dominant coupling point | repo-wide grep count |
| Legacy shims (public migration) | `openhands/server/{app,listen,middleware,shared,static,types,config}` are DEPRECATED banners re-exporting `app_server` symbols | `openhands/server/listen.py:1-8`; `openhands/server/types.py:1-7` |
| Analytics uses deprecated module | `from openhands.server.types import AppMode` in analytics package init | `openhands/analytics/__init__.py:23` |
| Tools consumed without local impl | `from openhands import tools` relies on namespace-merged external package | `openhands/app_server/event_callback/webhook_router.py:15` |
| App-server → SDK layering | Imports concentrated on `openhands.sdk.*` (models, llm, events, skills) and `openhands.agent_server.{models,utils,env_parser}` | grep histogram (e.g., `openhands/app_server/config.py:13`, `:72`) |
| Public API marker | PEP 561 `py.typed` present | `openhands/py.typed` |
| UI kit public surface | Explicit component exports; published as `@openhands/ui` v1.0.0-beta.9 with `exports` map for JS/CSS/types | `openhands-ui/index.ts:1-18`; `openhands-ui/package.json:2-27` |
| Frontend layering rule (documented) | "UI components → TanStack Query hooks → Data Access Layer (`frontend/src/api`) → API endpoints"; components must not call API client methods directly | `AGENTS.md:150-154` |
| Frontend rule violations | 42 files under `src/components` import `#/api` (35 non-`import type`), e.g., direct service class import in a card component | `frontend/src/components/features/conversation-panel/conversation-card/conversation-card.tsx:4`; `frontend/src/components/features/chat/model-messages.tsx:9` |
| No frontend lint boundary guard | ESLint config disables `import/no-extraneous-dependencies`; no `no-restricted-imports`/boundaries plugin | `frontend/.eslintrc:76` |
| Boundary-adjacent test (migrations only) | AST-based checker validating enterprise Alembic migration graph integrity, loaded in unit tests | `scripts/check_enterprise_migration_integrity.py:11-12`; `tests/unit/test_enterprise_migration_integrity.py:36-40` |
| Test tree mirrors source areas | `tests/unit/{app_server,integrations,mcp,server,storage,...}` | `tests/unit/README.md:1-12` |
| Analytics event catalog doc | Architecture lanes documented and tied to `AnalyticsService` singleton usage | `openhands/analytics/EVENTS.md:1-12`; `openhands/analytics/__init__.py:1-17` |

## Answers to Dimension Questions

### 1. Are modules cleanly separated?

Partially, with a clear macro split and muddy micro coupling. At the macro level, separation is strong: agent runtime (SDK), tools, agent server, and web app server are four distinct distributions (`pyproject.toml:60-62`), enterprise is a separate Poetry project (`enterprise/pyproject.toml:20-28`), and frontend/UI-kit are separate npm packages (`frontend/package.json:2`, `openhands-ui/package.json:2`). Inside `app_server/`, functional grouping is coherent and even documented (`openhands/app_server/README.md:9-27`), but the subpackages are heavily interwired: `event`, `sandbox`, `user_auth`, and `analytics` each participate in an import cycle with their neighbors (see Evidence table), and everything funnels through `app_server/config.py` (fan-in 33). Cleanliness is by-layer, not by-module.

### 2. Do dependencies flow in one direction?

Mostly, but with two sanctioned exceptions. The intended direction is `enterprise → openhands → (sdk | agent_server | tools)`: enterprise declares the OSS root as its only first-party dependency (`enterprise/pyproject.toml:29`) and imports `base_app` (`enterprise/saas_server.py:67`), while app_server consumes SDK/agent-server types (`openhands/app_server/config.py:13,72`). The exceptions are the OSS→enterprise back-edges implemented as indirection: `OPENHANDS_CONFIG_CLS` selects `SaasAppLifespanService` from the enterprise-only `server` package via lazy import (`openhands/app_server/config.py:174-183`), and store/service classes are resolved by string name through `get_impl` (`openhands/app_server/shared.py:23-27`, `openhands/app_server/utils/import_utils.py:43-78`). These keep the *import graph* acyclic between distributions while creating hidden, stringly-typed reverse edges that only fail at runtime if misconfigured. Within `app_server`, directionality breaks down entirely at subpackage granularity (four cycles).

### 3. Can modules be used independently?

The distribution boundaries yes; the in-repo Python modules largely no. Because `openhands-tools`, `openhands-sdk`, and `openhands-agent-server` are independent versioned packages merged only by `pkgutil.extend_path` (`openhands/__init__.py:1-9`), a consumer can install tools or the SDK without the web app — this is exactly how the dimension's guiding question ("can you use the tool system without pulling in the entire runtime?") is answered affirmatively at the package level, though installing the root metapackage pulls all of them (`pyproject.toml:249-251`). Conversely, in-repo modules assume the whole server context: `analytics` cannot be imported without `app_server.utils` (`openhands/analytics/analytics_context.py:21-22`) and the legacy `server.types` shim (`openhands/analytics/__init__.py:23`); nearly every router depends on the global config singleton in `shared.py`/`config.py`. On the frontend, `@openhands/ui` is genuinely standalone (own build, peer deps, Storybook; `openhands-ui/package.json`), though `frontend/package.json` does not yet consume it — the extraction is one-directional and incomplete.

### 4. Are public APIs distinguished from internal ones?

Weakly, by mechanism rather than policy. Positive signals: PEP 561 typing marker (`openhands/py.typed`); the HTTP surface is deliberately versioned and aggregated (`/api/v1` prefix, `openhands/app_server/v1_router.py:24`); the UI kit curates exports in `openhands-ui/index.ts:1-18`; deprecated paths are explicitly marked and re-exported to keep old entrypoints working (`openhands/server/listen.py:1-8`, `make start-backend` still targets the old path per `AGENTS.md`). Missing signals: no `__all__` discipline or `_private` module conventions across `app_server/`; no import-linter or API-surface test; the extensibility API itself is string-based (`get_impl` accepts any fully-qualified name, validated only as a subclass check at call time, `openhands/app_server/utils/import_utils.py:39`), so "public" vs "internal" is unenforceable for downstream overriders like enterprise.

## Architectural Decisions

1. **Extract the agent core into separate distributions (V1 split).** Runtime, tools, and agent server were moved out of this repo into `openhands-sdk`/`openhands-tools`/`openhands-agent-server`, exact-pinned (`pyproject.toml:60-62,249-251`). This makes the harness layers independently consumable and versionable, at the cost of coordinating releases across repos (the pins must be bumped together).
2. **One namespace, many packages.** Rather than renaming, all distributions share the `openhands` top-level package via `pkgutil.extend_path` (`openhands/__init__.py:1-9`). Consumers get a single coherent import root; the cost is that attribution of any `openhands.*` symbol to a distribution requires external knowledge.
3. **Composition-root DI instead of static wiring.** All services are declared as `*Injector` fields on `AppServerConfig` and defaulted from environment (`openhands/app_server/config.py:191-237,240-420`), injected through FastAPI via `Injector.depends` (`openhands/app_server/services/injector.py:31-40`). This is what lets enterprise replace storage/sandbox/event implementations without forking routers.
4. **Runtime substitution via `get_impl`.** Implementation selection by fully-qualified class name from config (`openhands/app_server/utils/import_utils.py:43-78`), used both inside OSS (`openhands/app_server/shared.py:23-27`) and by enterprise processors (`enterprise/storage/maintenance_task.py:92`). Deliberate tradeoff: flexibility over static verifiability.
5. **Strangler-fig migration for V0→V1.** Old `openhands.server` entrypoints kept as thin deprecation shims re-exporting `app_server` (`openhands/server/listen.py:1-8`, `openhands/server/middleware.py:1-8`), keeping the public contract stable while internals move.
6. **Enterprise as overlay, not fork.** SaaS server imports the OSS app and mounts extra routers (`enterprise/saas_server.py:67+`), sets `OPENHANDS_CONFIG_CLS=server.config.SaaSServerConfig` (`enterprise/saas_server.py:7-9`), and ships its own Alembic graph guarded by an AST integrity script (`scripts/check_enterprise_migration_integrity.py:11-12`).

## Notable Patterns

- **Injector + Depends pairing**: every service has a `getX_service()` context manager and a `depends_x_service()` FastAPI helper generated from the same injector (`openhands/app_server/config.py:436-602`), standardizing how modules consume each other's services.
- **Deferred-import cycle breaking**: ~55 function-scoped imports in `app_server` (e.g., `openhands/app_server/config.py:242-283`, `openhands/app_server/settings/settings_router.py:238`, `openhands/app_server/secrets/secrets_router.py:137`) — an idiom that keeps Python import-time happy around the `event ↔ app_conversation ↔ sandbox` triangle.
- **Deprecated-shim modules**: uniform banner + re-export shape across `openhands/server/*` provides a discoverable migration trail.
- **Registration-by-import**: side-effect registration of callback processors via bare import in the config module (`openhands/app_server/config.py:12` importing `openhands.app_server.event_callback`), plus dynamic processor discovery with `pkgutil`/`importlib` in `webhook_router.py:8-10`.
- **Documented data-access layering (frontend)**: components → query hooks → `src/api` services (`AGENTS.md:150-154`), with per-resource directories like `frontend/src/api/conversation-service/`.
- **Env-var fallback chains** preserving legacy configuration surfaces (`OH_PERSISTENCE_DIR` → `FILE_STORE_PATH`, `openhands/app_server/config.py:75-89`), softening boundary moves for operators.

## Tradeoffs

- **Flexibility vs. static safety**: `get_impl` and `OPENHANDS_CONFIG_CLS` make behavior configurable per deployment but move a class of errors (typo'd class names, non-subclasses) from import/typecheck time to request/startup time; validation is a single `issubclass` assert (`openhands/app_server/utils/import_utils.py:39`).
- **Namespace unity vs. provenance clarity**: `extend_path` gives consumers one import root but means `openhands.sdk`, `openhands.tools`, `openhands.agent_server`, and app code are only distinguishable by submodule path, and version skew is possible if pins drift (`pyproject.toml:249-251`).
- **Layered cleanliness vs. change friction**: strict one-way enterprise→OSS dependency is preserved, but every new OSS↔enterprise seam requires adding another stringly-typed hook in OSS config, which accumulates in `config_from_env()` (now ~180 lines, `openhands/app_server/config.py:240-420`).
- **Published UI kit vs. adoption**: extracting `@openhands/ui` creates a clean reusable boundary (`openhands-ui/package.json:2-27`), yet the main app hasn't adopted it (`frontend/package.json` has no `@openhands/ui` dependency), so two component systems coexist in-tree.
- **Global config singleton vs. testability/multi-tenancy**: `get_global_config()` memoizes one process-wide `AppServerConfig` (`openhands/app_server/config.py:423-433`), simplifying injection everywhere but baking in process-global state.

## Failure Modes / Edge Cases

- **Import-order sensitivity**: because `shared.py` resolves store implementations at module import time (`openhands/app_server/shared.py:23-27`), anything importing `app_server.shared` before env/config is set gets the wrong (default) implementation; tests and embedders must control environment before first import.
- **Silent fallbacks on misconfiguration**: `SaasAppLifespanService` selection keys on `'saas' in os.getenv('OPENHANDS_CONFIG_CLS')` lowercased substring match (`openhands/app_server/config.py:177`) — a renamed class or differently-cased env value silently yields `OssAppLifespanService` with no error.
- **Cycle regressions**: the four subpackage cycles currently hold only because of scattered lazy imports; a new module-level import added casually (e.g., importing `EventService` at top level of a conversation module) reintroduces `ImportError`/partial-init bugs. Nothing in CI detects this class of regression (no cycle checker found).
- **Legacy shim divergence**: `openhands/server/app.py:4-6` notes it "does NOT include middleware setup" unlike `listen.py` — two deprecated entrypoints with subtly different behavior invite misconfiguration during migration.
- **Frontend boundary erosion**: with 42 component files already reaching into `#/api` (`frontend/src/components/features/conversation-panel/conversation-card/conversation-card.tsx:4`) and no ESLint guard, the documented architecture will continue to drift; type-only imports are also indistinguishable from calls to tooling without `import type` enforcement.
- **Cross-repo pin desync**: bumping `openhands-ai` requires coordinated lockfile updates in `enterprise/` (documented in `AGENTS.md`), otherwise enterprise builds silently resolve older SDK behavior.

## Future Considerations

- Add automated boundary enforcement: `import-linter` contracts for `app_server` subpackages (forbid `event ← app_conversation` etc.), and ESLint `no-restricted-imports` for `#/api` in `frontend/src/components`. Both are mechanical wins given existing CI (`dev_config/python/.pre-commit-config.yaml`).
- Break the `event ↔ app_conversation ↔ sandbox` triangle by promoting shared models (`SandboxStatus`, conversation-info protocols) into a leaf `types`/`models` subpackage, eliminating most of the 55 deferred imports.
- Type the OSS↔enterprise seams: replace raw strings (`OPENHANDS_CONFIG_CLS`, `settings_store_class`) with validated entry-point groups or typed config objects, failing fast at startup with actionable errors.
- Finish the `@openhands/ui` adoption in `frontend/` or document why two component libraries coexist, so the UI boundary reflects reality.
- Retire remaining deprecated-module dependencies inside first-party code (e.g., `openhands/analytics/__init__.py:23` should import `openhands.app_server.types` directly) so shims can actually be deleted.
- Publish which `openhands.*` subtrees belong to which distribution (generated manifest), recovering the provenance lost to the shared namespace.

## Questions / Gaps

- **No boundary/cycle tests found.** Searched `tests/` for `cycle|circular` and `dev_config/` + `pyproject.toml` for `import-linter|lint-imports`; the only structural check is the enterprise migration-graph script (`scripts/check_enterprise_migration_integrity.py:11-12`). If boundary tests exist elsewhere (CI scripts?), they were outside the searched set.
- **SDK/tools internals not inspectable here.** By design of this study (single-source isolation), `openhands-sdk`, `openhands-tools`, and `openhands-agent-server` are pinned externals (`pyproject.toml:60-62`); their internal boundaries can only be assessed via this repo's consumption surface (imports listed in Evidence). Claims about their internal cleanliness are out of scope.
- **Whether `@openhands/ui` is consumed elsewhere** (e.g., docs site or other apps) could not be verified within this source; only its publication metadata is in-tree (`openhands-ui/PUBLISHING.md`).
- **Enforcement status of AGENTS.md rules** (e.g., whether PR review rejects direct `#/api` usage) is inferred from code state (42 violations persist), not from process documentation.

---

Generated by `Dimension 22.01: Package and Module Boundaries` against `openhands`.
