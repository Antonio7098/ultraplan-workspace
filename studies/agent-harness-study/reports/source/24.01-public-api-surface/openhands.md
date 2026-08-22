# Source Analysis: openhands

## Dimension 24.01: Public API Surface

### Source Info

| Field | Value |
|-------|-------|
| Name | openhands |
| Path | `studies/agent-harness-study/sources/openhands` |
| Language / Stack | Python 3.12–3.13 (Poetry/uv), FastAPI, React/TypeScript frontend, npm component library, Docker |
| Analyzed | 2026-08-22 |

## Summary

OpenHands' public API is primarily an **HTTP API, not a Python import API**. The top-level PyPI package `openhands-ai` (`pyproject.toml:8`) is now a thin shell: `openhands/__init__.py` is a namespace package exporting only `__version__` and `get_version` (`openhands/__init__.py:1-9`), and the agentic core has been extracted into three pinned external SDK packages — `openhands-sdk`, `openhands-tools`, `openhands-agent-server`, all locked to `==1.29.0` (`pyproject.toml:60-62`, `pyproject.toml:249-251`). The README confirms the migration: agent/agent-server source lives in `OpenHands/software-agent-sdk` and the "Agent Canvas" UI in `OpenHands/agent-canvas` (`README.md:49-54`).

What remains in this repo as the real public surface is:

1. **The V1 REST API** — a FastAPI app at `openhands.app_server.app:app` (`openhands/app_server/app.py:54-60`) exposing a versioned `/api/v1` router aggregating 13 sub-routers (events, app-conversations, pending-messages, sandboxes, sandbox-specs, settings, secrets, users, skills, webhooks, web-client, git, config) (`openhands/app_server/v1_router.py:24-37`), plus root health endpoints `/alive`, `/health`, `/server_info`, `/ready` (`openhands/app_server/status/status_router.py:8,17,28,39`) and an MCP server mounted at `/mcp` (`openhands/app_server/app.py:33,59`) exposing five git-provider PR/MR tools (`openhands/app_server/mcp/mcp_router.py:147,216,290,357,424`).
2. **An extension API** built on env-driven dependency injection: `get_impl()` dynamically imports replacement implementations by fully-qualified name (`openhands/app_server/utils/import_utils.py:43-78`), wired through `AppServerConfig` injector fields parsed from `OH_*` env vars (`openhands/app_server/config.py:214-237`, `openhands/app_server/config.py:285`) and `ServerConfig` class-impl attributes like `user_auth_class` (`openhands/app_server/server_config/server_config.py:17-25`).
3. **A UI component library** — `@openhands/ui` npm package at `1.0.0-beta.9` with an exports map (`openhands-ui/package.json:4-30`).
4. **A deployment surface** — Docker image whose `CMD` is `uvicorn openhands.server.listen:app` (`containers/app/Dockerfile`, final `CMD` line) and Makefile target `start-backend` (`Makefile:262`).

The API model is coherent and heavily tested, but the surface is in a visibly transitional state: the documented launch path (`make start-backend`, Docker `CMD`) still targets the deprecated `openhands.server.listen` shim (`openhands/server/listen.py:1-4`), several docs reference paths and flags that no longer exist, and deprecation is expressed only in comments rather than tool-enforced markers.

## Rating

**6 / 10** — Present but inconsistent. The REST API is versioned, tagged, and covered by an OpenAPI-schema-generation test (`tests/unit/server/test_openapi_schema_generation.py:83-110`) plus ~50 router/service test files (`tests/unit/app_server/`), and extension points are explicit ABCs with docstrings (`openhands/app_server/user_auth/user_auth.py:35-45`, `openhands/app_server/services/injector.py:12-21`). However, the score is held back by: entry-point inconsistency (Makefile `Makefile:262` and Docker `CMD` launch the deprecated shim while `openhands/server/__main__.py:1-2` tells users to use `openhands.app_server.app`), dead/stale configuration (`enable_v1` at `openhands/app_server/server_config/server_config.py:28` has no consumers; `Development.md:66` still claims `ENABLE_V1=0` disables V1 routes), comment-only deprecation with no `@deprecated` decorators despite the `deprecation` dependency (`pyproject.toml:33,256`), and thin in-repo API documentation (`openhands/app_server/README.md:1-27` lists module names only, with no endpoint reference or runnable examples).

## Evidence Collected

Every entry includes a file path with line numbers. Paths are relative to the source root `studies/agent-harness-study/sources/openhands/`.

| Area | Evidence | File:Line |
|------|----------|-----------|
| Python package identity | PyPI package `openhands-ai`, version managed by poetry-dynamic-versioning | `pyproject.toml:8`, `pyproject.toml:19-20` |
| Minimal Python export | Namespace package exports only `__version__`, `get_version` | `openhands/__init__.py:4-9` |
| Core extracted to external SDK | Dependencies `openhands-sdk==1.29.0`, `openhands-agent-server==1.29.0`, `openhands-tools==1.29.0` | `pyproject.toml:60-62`, `pyproject.toml:249-251` |
| Repo transition statement | "source code ... lives in software-agent-sdk / agent-canvas" | `README.md:49-54` |
| Public FastAPI app | `app = FastAPI(title='OpenHands', ...)` with MCP mount | `openhands/app_server/app.py:54-60` |
| Versioned REST surface | `/api/v1` prefix aggregating 13 routers | `openhands/app_server/v1_router.py:24-37` |
| Health/ops endpoints | `/alive`, `/health`, `/server_info`, `/ready` | `openhands/app_server/status/status_router.py:8,17,28,39` |
| MCP server surface | `mcp_server = FastMCP('mcp', mask_error_details=True)`; mounted at `/mcp`; tools `create_pr`, `create_mr`, `create_bitbucket_pr`, `create_bitbucket_data_center_pr`, `create_azure_devops_pr`; Tavily proxy mount | `openhands/app_server/mcp/mcp_router.py:43`, `openhands/app_server/app.py:33,59`, `openhands/app_server/mcp/mcp_router.py:147-148,216,290,357,424`, `openhands/app_server/mcp/mcp_router.py:49-75` |
| Auth surface for API | `X-Session-API-Key` header dependency injected into routers via `get_dependencies()` | `openhands/app_server/utils/dependencies.py:9-32`, e.g. `openhands/app_server/settings/settings_router.py:60-64` |
| Sandbox-scoped secret API | `GET /api/v1/sandboxes/{id}/settings/secrets[/{name}]` authenticated via session key | `openhands/app_server/sandbox/sandbox_router.py:154-217` |
| Extension mechanism | `get_impl()` — "runtime substitution of implementations" via fully-qualified class names | `openhands/app_server/utils/import_utils.py:43-78` |
| Env-driven DI config | `config_from_env()` parses `OH_*` env into `AppServerConfig`; injector fields for event/sandbox/settings/etc. | `openhands/app_server/config.py:240-285`, `openhands/app_server/config.py:214-237` |
| Pluggable impl selection | `ServerConfig.settings_store_class` / `secret_store_class` / `user_auth_class`; `OPENHANDS_CONFIG_CLS` override | `openhands/app_server/server_config/server_config.py:17-28`, `openhands/app_server/server_config/server_config.py:48-56` |
| Auth extension point | `UserAuth` ABC documented as extension point, resolved via `get_impl()` | `openhands/app_server/user_auth/user_auth.py:35-45`, `openhands/app_server/user_auth/user_auth.py:119-129` |
| DI base class | `Injector` ABC with `inject`/`context`/`depends` | `openhands/app_server/services/injector.py:12-34` |
| Enterprise extension pattern | Enterprise imports OSS `app` singleton and adds routers/middleware directly | `enterprise/saas_server.py:67-71`, `enterprise/saas_server.py:81-184` |
| Deprecated legacy shims | `openhands/server/listen.py`, `app.py`, `types.py`, `static.py`, `shared.py`, `middleware.py`, `config/server_config.py`, `openhands/version.py` all marked DEPRECATED re-exports | `openhands/server/listen.py:1-4`, `openhands/server/app.py:1-6`, `openhands/version.py:1-6` |
| Entry-point inconsistency | `make start-backend` and Docker `CMD` run deprecated `openhands.server.listen:app`; `__main__` says use `openhands.app_server.app:app` | `Makefile:262`, `containers/app/Dockerfile` (final CMD), `openhands/server/__main__.py:1-2` |
| OpenAPI contract tested | `/openapi.json` generation test asserts `/api/v1/settings` and `/health` present | `tests/unit/server/test_openapi_schema_generation.py:83-110` |
| API test coverage | ~50 unit test files for routers/services (settings, secrets, events, webhooks, sandboxes, profiles, git) | `tests/unit/app_server/` (e.g. `tests/unit/app_server/test_settings_api.py`, `tests/unit/app_server/test_sandbox_secrets_router.py`) |
| UI component library API | Published npm package `@openhands/ui` `1.0.0-beta.9` with `exports` map and `py.typed`-equivalent `types` entry | `openhands-ui/package.json:4-30` |
| Frontend HTTP client convention | Shared axios instance + per-service modules documented for extension authors | `frontend/src/api/README.md:1-47`, `frontend/src/api/open-hands-axios.ts:4` |
| Typed package marker | PEP 561 `py.typed` present | `openhands/py.typed` |
| Stale config flag | `enable_v1` read from `ENABLE_V1` env but referenced nowhere else | `openhands/app_server/server_config/server_config.py:28` (no other hits in `openhands/` or `enterprise/`) |
| Stale docs | `Development.md` claims `ENABLE_V1=0` disables V1 routes; README links `docs/SELF_HOSTING.md` which does not exist; `pydoc-markdown.yml` targets missing `docs/modules`; AGENTS.md references nonexistent `openhands/cli/utils.py` | `Development.md:66`, `README.md:63`, `pydoc-markdown.yml:11`, `AGENTS.md` (model-configuration section) |
| Accidental packaging surface | `pyproject.toml` and `poetry.lock` shipped inside the `openhands` package data | `pyproject.toml:149-153` |
| No CLI/console scripts | No `[project.scripts]` section; no `openhands/cli` module in tree | `pyproject.toml` (whole file), `openhands/` directory listing |

## Answers to Dimension Questions

### 1. What is the intended public API surface?

Three distinct surfaces, each for a different audience:

- **Operators/self-hosters**: the FastAPI app `openhands.app_server.app:app` (`openhands/app_server/app.py:54-60`) launched via uvicorn (Docker `CMD` in `containers/app/Dockerfile`, `make start-backend` at `Makefile:262`), configured entirely through environment variables — `OH_*` for the DI config (`openhands/app_server/config.py:285`), plus legacy fallbacks like `FILE_STORE_PATH` and `PERMITTED_CORS_ORIGINS` (`openhands/app_server/config.py:79-113`).
- **Client developers (frontend/third-party integrations)**: the versioned REST API under `/api/v1` (`openhands/app_server/v1_router.py:24-37`) with OpenAPI schema (`tests/unit/server/test_openapi_schema_generation.py:90-110`), the MCP endpoint `/mcp` (`openhands/app_server/app.py:59`), and webhook callbacks at `/api/v1/webhooks/*` (`openhands/app_server/event_callback/webhook_router.py:62,333,453,531`).
- **Extension authors**: runtime-substitutable service implementations via `get_impl()` (`openhands/app_server/utils/import_utils.py:43-78`) selected through `ServerConfig` attributes (`openhands/app_server/server_config/server_config.py:17-25`) or `OH_*_KIND` injector fields (`openhands/app_server/config.py:214-223`; example value documented at `openhands/app_server/sandbox/dynamic_remote_sandbox_spec_service.py:106`). The `enterprise/saas_server.py:67-184` pattern of importing and mutating the OSS `app` singleton is the in-tree proof that the FastAPI app object is itself a supported extension seam.

Notably, a **Python import API for agent developers is explicitly out of scope here** — that role moved to the external `openhands-sdk` package (`pyproject.toml:60-62`, `README.md:52`).

### 2. Is the stable API easy to distinguish from internal implementation details?

**Partially.** On the HTTP side, yes: everything public sits behind the explicit `/api/v1` prefix (`openhands/app_server/v1_router.py:24`), routers carry OpenAPI tags (`openhands/app_server/app_conversation/app_conversation_router.py:106-107`, `openhands/app_server/sandbox/sandbox_router.py:32-33`), and the schema is regression-tested. On the Python side, the boundary is murkier:

- The `openhands.server.*` package is a pure deprecation shim layer (`openhands/server/listen.py:1-4`, `openhands/server/app.py:1-6`), yet the **official launch commands still import it** (`Makefile:262`, Docker `CMD`), so the deprecated path is the one most users will actually execute. `openhands/server/__main__.py:1-2` contradicts this by directing users to `openhands.app_server.app:app`.
- `openhands/app_server/user_auth/user_auth.py:1-8` carries a banner "LEGACY V0 CODE - Deprecated since version 1.0.0, scheduled for removal April 1, 2026" — but the file lives inside the V1 `app_server` tree and defines the **active V1** `UserAuth` extension point consumed by V1 routers (`openhands/app_server/user_auth/user_auth.py:119-129`). A reader cannot tell from location whether this is safe to build on.
- `__all__` exports exist only in the shims (`openhands/version.py:8`) and `openhands/app_server/types.py:37-43`; the rest of `app_server` has no explicit export control, so any module is importable.
- Deprecation is comment-only; the `deprecation` library is a declared dependency (`pyproject.toml:33,256`) but no `@deprecated` decorator is used anywhere in `openhands/` (grep found zero matches), so tooling cannot detect deprecated surface.

### 3. Does the API expose the right level of abstraction for agent harness users?

**For harness integrators, yes; for agent developers, the abstraction lives elsewhere.** The REST API operates at the right granularity for a harness control plane — conversations, events, sandboxes, settings, secrets, webhooks (`openhands/app_server/v1_router.py:25-37`) — without leaking agent internals: the agentic loop, tools, and LLM plumbing belong to the pinned `openhands-sdk`/`openhands-agent-server` dependencies (`pyproject.toml:60-62`), and this server only orchestrates them (e.g., `openhands/app_server/app_conversation/live_status_app_conversation_service.py:179-180` reads the SDK version via `importlib.metadata`). Secrets are exposed to sandboxes through a deliberate minimal contract — name listing and single-value fetch gated by `X-Session-API-Key` (`openhands/app_server/sandbox/sandbox_router.py:154-217`) — rather than bulk dumps. The MCP tools are narrowly scoped to PR/MR creation across five providers (`openhands/app_server/mcp/mcp_router.py:147-424`).

Two abstraction leaks are worth noting: (a) the settings API returns LLM base URLs conditionally nulled so "the frontend can display basic mode" (`openhands/app_server/settings/settings_router.py:150-156`) — UI presentation concerns embedded in the API contract; (b) `ServerConfig.get_config()` (`openhands/app_server/server_config/server_config.py:34-45`) hands PostHog client keys and feature flags to the frontend, coupling server config to analytics.

### 4. Are examples sufficient to use the API correctly without reading internals?

**No — not in-repo.** The only API documentation inside the source is a 27-line module listing (`openhands/app_server/README.md:1-27`) and docstrings on individual endpoints (e.g., `openhands/app_server/settings/settings_router.py:102-113`). There is no endpoint reference, no curl/SDK examples, and no runnable quickstart for the API in the tree: the README's quickstart now targets the external `agent-canvas` npm package and Docker image (`README.md:56-110`), and links to `docs/SELF_HOSTING.md` (`README.md:63`) which does not exist in this source (no `docs/` directory). `pydoc-markdown.yml:11` configures Python doc generation into `docs/modules`, which also does not exist. The de facto usage documentation is the test suite — `tests/unit/app_server/` contains ~50 files exercising routers and services (settings, secrets, webhooks, sandboxes, profiles, git), which demonstrates correct call patterns but is not discoverable API documentation. The frontend API layer is the best-documented consumer pattern (`frontend/src/api/README.md:1-102`), but it is internal to this repo's SPA rather than a third-party contract.

## Architectural Decisions

1. **Extract the agent core; keep the app server.** The agentic engine moved to pinned external packages (`pyproject.toml:60-62`), leaving this repo as the control plane. This makes the HTTP API the stable contract and shrinks the Python import surface to version helpers only (`openhands/__init__.py:9`).
2. **Version the REST surface, deprecate the old namespace.** All new routes live under `/api/v1` (`openhands/app_server/v1_router.py:24`); the old `openhands.server` package is preserved as re-export shims (`openhands/server/app.py:1-15`) for backward compatibility during migration.
3. **Environment-driven dependency injection as the extension API.** `AppServerConfig` fields are injectors (`openhands/app_server/config.py:214-237`) populated by a pydantic `from_env(AppServerConfig, 'OH')` parser (`openhands/app_server/config.py:285`); implementations are resolved by fully-qualified class name through `get_impl()` with subclass validation and caching (`openhands/app_server/utils/import_utils.py:34-78`). This lets deployments swap storage, auth, and sandbox backends without forking the server.
4. **Composition by app-object mutation for the enterprise tier.** Rather than a plugin registry, `enterprise/saas_server.py:81-184` imports the OSS `app` and adds routers/middleware conditionally on env vars (e.g., `GITHUB_APP_CLIENT_ID` gate at `enterprise/saas_server.py:93-105`). Simple, but it makes the singleton app object a load-bearing public contract.
5. **MCP as a first-class API protocol.** The server embeds a FastMCP instance at `/mcp` (`openhands/app_server/app.py:33,59`) with `mask_error_details=True` (`openhands/app_server/mcp/mcp_router.py:43`) and proxies third-party MCP (Tavily) under a namespace to avoid leaking API keys (`openhands/app_server/mcp/mcp_router.py:49-75`).

## Notable Patterns

- **Shim-for-backcompat**: every deprecated module is a comment-headed re-export of its `app_server` successor (`openhands/server/listen.py:1-8`, `openhands/server/config/server_config.py:1-11`, `openhands/version.py:1-8`) — a consistent mechanical pattern, even though launch scripts haven't caught up.
- **Injector ABC + `Depends` bridging**: `Injector.inject`/`context`/`depends` (`openhands/app_server/services/injector.py:16-34`) unifies request-scoped DI with background usage, e.g. webhook callbacks constructing state manually (`openhands/app_server/event_callback/webhook_router.py:561-573`).
- **Auth surfaced in OpenAPI deliberately**: `get_dependencies()` appends a no-fail `APIKeyHeader` dependency purely "so it appears in OpenAPI Docs" (`openhands/app_server/utils/dependencies.py:16-31`) — a thoughtful docs-as-code touch.
- **Schema regression test as API guardrail**: `test_openapi_schema_generation` (`tests/unit/server/test_openapi_schema_generation.py:83-110`) catches serialization breakage in the public contract.
- **Namespace package for SDK coexistence**: `openhands/__init__.py:1-4` uses `pkgutil.extend_path` so the external `openhands-sdk`/`openhands-tools`/`openhands-agent-server` packages (which also ship a top-level `openhands` package) merge into one namespace.

## Tradeoffs

- **HTTP-first stability vs. import-level instability**: pinning the SDK (`==1.29.0`, `pyproject.toml:60-62`) protects the REST contract but means Python-level consumers of this repo get almost nothing (`openhands/__init__.py:9`); the real library API is versioned and released in a different repository.
- **Env-var configurability vs. discoverability**: the `OH_*` injector system (`openhands/app_server/config.py:285`) is powerful, but the valid keys are implicit in pydantic field names and scattered defaults; nothing enumerates them for operators (the `AppServerConfig` field descriptions at `openhands/app_server/config.py:194-213` are the closest thing).
- **Backcompat shims vs. single source of truth**: keeping `openhands.server` alive eases migration but produced the current split-brain where the documented run command (`Makefile:262`, Docker `CMD`) targets the deprecated import path.
- **App-object mutation vs. plugin isolation**: enterprise extension by direct mutation (`enterprise/saas_server.py:81-184`) is fast to build on but creates ordering and override hazards — e.g., it *overrides* the OSS `/api/v1/users/me` endpoint in place (`enterprise/saas_server.py:143-145`), which only works because of import-order coupling.
- **Comment-only deprecation vs. tooling**: zero runtime warnings for deprecated imports; consumers get no signal until removal.

## Failure Modes / Edge Cases

- **Launching the deprecated path**: `make start-backend` (`Makefile:262`) and the Docker `CMD` run `openhands.server.listen:app`, which silently drops the note that `openhands/server/app.py:5-6` warns the shim "does NOT include middleware setup" for direct `app` imports — operators following different docs can get subtly different middleware stacks.
- **Dead feature flag**: `ENABLE_V1` is parsed into `enable_v1` (`openhands/app_server/server_config/server_config.py:28`) and documented as functional (`Development.md:66`), but no code consumes it — setting `ENABLE_V1=0` does nothing in the current tree, a silent operator-facing lie.
- **Removal cliff for `UserAuth`**: the active V1 auth extension point is labeled "scheduled for removal April 1, 2026" (`openhands/app_server/user_auth/user_auth.py:1`) — extension authors subclassing it face an ambiguous deprecation deadline inside the V1 tree.
- **Stale agent/dev instructions**: AGENTS.md instructs model additions in `openhands/cli/utils.py` and settings in `openhands/storage/data_models/settings.py`, neither of which exists in this tree — automation following these docs will fail or, worse, "succeed" against the wrong repo.
- **Secret endpoint blast radius**: `GET /api/v1/sandboxes/{id}/settings/secrets/{secret_name}` returns raw secret values as plain text (`openhands/app_server/sandbox/sandbox_router.py:185-217`); its safety depends entirely on `_valid_sandbox_from_session_key` correctly binding the `X-Session-API-Key` to the sandbox owner (`openhands/app_server/sandbox/sandbox_router.py:143-151`). A weak or guessable `SESSION_API_KEY` (single global env value, `openhands/app_server/utils/dependencies.py:9`) would expose all user secrets.
- **Version detection fragility**: `get_version()` parses `pyproject.toml` by line-scanning for `version =` (`openhands/app_server/version.py:15-20`) and reports `'unknown'` on any failure (`openhands/app_server/version.py:41-47`) — the version surfaced in the OpenAPI info block (`openhands/app_server/app.py:57`) can silently degrade.

## Future Considerations

- **Finish the entry-point migration**: point `Makefile:262` and the Docker `CMD` at `openhands.app_server.app:app`, then delete the `openhands.server` shims on the announced schedule.
- **Enforce deprecation mechanically**: use the already-declared `deprecation` dependency (`pyproject.toml:256`) to decorate shim modules and `UserAuth` so consumers get runtime warnings and static analyzers can flag usage.
- **Publish an API contract artifact**: commit the generated OpenAPI schema (already tested at `tests/unit/server/test_openapi_schema_generation.py:90`) and diff it in CI to make breaking REST changes reviewable.
- **Resolve the `UserAuth` ambiguity**: either move the file out of the "Legacy-V0" label or re-scope the removal banner; as written (`openhands/app_server/user_auth/user_auth.py:1-8`) it undermines the V1 extension story.
- **Document the `OH_*` configuration surface**: generate an env-var reference from `AppServerConfig` (`openhands/app_server/config.py:191-237`) and `ServerConfig` (`openhands/app_server/server_config/server_config.py:9-28`) instead of leaving operators to read source.
- **Remove or honor `ENABLE_V1`** (`openhands/app_server/server_config/server_config.py:28`) and reconcile `Development.md:66` with actual behavior.
- **Clean packaging accidents**: stop shipping `pyproject.toml`/`poetry.lock` inside the wheel (`pyproject.toml:151-152`) and remove the dangling `pydoc-markdown.yml` target (`pydoc-markdown.yml:11`) or restore the `docs/` tree it expects.

## Questions / Gaps

- **REST API stability policy**: no in-repo document states the compatibility guarantee for `/api/v1` (deprecation windows, breaking-change policy). The `release-please-config.json` at the repo root suggests automated versioning, but I found no API-changelog discipline. No evidence found within the source boundary.
- **SDK contract location**: the actual agent-facing API (`openhands-sdk`) is pinned as a dependency (`pyproject.toml:60-62`) but its source is outside this directory (`README.md:52`); its public API could not be studied here per source-isolation rules.
- **MCP authentication details**: `mask_error_details=True` (`openhands/app_server/mcp/mcp_router.py:43`) and `get_mcp_api_key` on `UserAuth` (`openhands/app_server/user_auth/user_auth.py:90-92`) imply per-user MCP auth, but I did not trace the full `/mcp` auth middleware in this pass; the exact enforcement path needs a dedicated look.
- **Frontend↔backend contract tests**: the frontend has a typed API layer (`frontend/src/api/README.md:5`), but I found no schema-drift test binding the TypeScript types to the backend OpenAPI schema. No evidence found; searched `frontend/src/api/` and `tests/`.
- **`openhands/analytics` public-ness**: the module ships `EVENTS.md` (`openhands/analytics/EVENTS.md`) documenting analytics events, but whether these events are a supported integration surface or internal telemetry is not stated anywhere in-tree.

---

Generated by `Dimension 24.01: Public API Surface` against `openhands`.
