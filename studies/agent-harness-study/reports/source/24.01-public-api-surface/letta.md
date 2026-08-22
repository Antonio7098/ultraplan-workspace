# Source Analysis: letta

## 24.01 Public API Surface

### Source Info

| Field | Value |
|-------|-------|
| Name | letta |
| Path | `studies/agent-harness-study/sources/letta` |
| Language / Stack | Python 3.11+ (FastAPI server, Typer CLI, Pydantic v2 schemas, SQLAlchemy/SQLModel ORM); SDKs generated in Python & TypeScript via Fern |
| Analyzed | 2026-08-22 |

## Summary

Letta's public API is a versioned REST surface (`/v1`) served by a FastAPI application, with client SDKs (`letta-client` for Python, `@letta-ai/letta-client` for TypeScript) generated from a checked-in OpenAPI contract (`fern/openapi.json`). The repository itself is the *server*; it deliberately does not ship its own HTTP client — `letta/client/__init__.py` is empty and the only in-tree "client" code is dead (see Failure Modes). Instead, the server package depends on the generated SDK (`letta-client>=1.7.12`, `pyproject.toml:46`) and dogfoods it in its own integration tests.

The public entry points are:

1. **REST API**: ~33 routers mounted under `/v1` (`letta/server/rest_api/routers/v1/__init__.py:34-68`; prefix from `API_PREFIX = "/v1"`, `letta/constants.py:33`), an undocumented `/latest` alias bound to the newest version (`letta/server/rest_api/app.py:857`), admin routes under `/v1/admin` (`ADMIN_PREFIX`, `letta/constants.py:32`; mounted at `app.py:860-861`), an OpenAI-compatible `POST /v1/chat/completions` endpoint (`letta/server/rest_api/routers/v1/chat_completions.py:117-124`), and `/v1/health`. Default port 8283 (`letta/server/constants.py:6`).
2. **Generated SDKs**: Fern consumes `fern/openapi.json` (~1.7 MB, OpenAPI 3.1) plus overrides to emit Python/TS SDKs; publish pipelines live in `.github/workflows/fern-sdk-python-publish.yml:39-50` and `fern-sdk-typescript-publish.yml`.
3. **CLI**: a single `letta server` command (`pyproject.toml:84-85` → `letta/main.py:8`; implementation `letta/cli/cli.py:17-42`). The CLI has been intentionally minimized; the WebSocket option raises `NotImplementedError` (`cli.py:42`).
4. **Python import surface**: `letta/__init__.py:24-47` re-exports core schema types (`AgentState`, `Block`, `ChatMemory`, `Memory`, `Tool`, ...) for library/plugin use (e.g., tool sandbox code injected into agent runtimes).

The API contract is curated: deprecation is expressed with FastAPI's native `deprecated=True` markers that flow into OpenAPI and then into SDKs, internal-only endpoints are namespaced `_internal_*`, and Fern overrides hide selected endpoints from generated SDKs via `x-fern-ignore`.

## Rating

**7 / 10**

Rationale: The public API model is clear and unusually well-instrumented for an OSS project — a single versioned REST contract, SDK generation driven from a checked-in spec, `operation_id` discipline across most routes (207 of 252 route declarations have an adjacent `operation_id`), explicit deprecation markers, and SDK-level contract tests (`tests/sdk/`) that exercise the generated client against a live server. What keeps it out of 8–10: there is no CI gate verifying that `fern/openapi.json` matches what the FastAPI app actually emits (regeneration is a manual script writing to CWD), `_internal_*` endpoints are still publicly mounted and schema-documented, dead legacy modules remain importable (`letta/client/streaming.py`), and runnable examples inside the repo are essentially absent (docs are external). A new integrator can reliably build against the SDK without touching internals, but they depend on the human process keeping the spec in sync.

## Evidence Collected

Every entry includes a file path with line numbers. Format: `path/to/file:NN`.

| Area | Evidence | File:Line |
|------|----------|-----------|
| Package + CLI entry point | `[project.scripts] letta = "letta.main:app"` | `pyproject.toml:84-85` |
| CLI registers only `server` as default command | `app.command(name="server")(server)`; bare invocation runs server | `letta/main.py:8-16` |
| CLI server command options (port/host/debug/reload/secure/localhttps) | `def server(...)` signature | `letta/cli/cli.py:17-26` |
| WebSocket server option deprecated | `raise NotImplementedError("WS suppport deprecated")` | `letta/cli/cli.py:41-42` |
| REST prefix constants | `ADMIN_PREFIX = "/v1/admin"`, `API_PREFIX = "/v1"` | `letta/constants.py:32-33` |
| Router registry (33 routers incl. agents, blocks, tools, sources, runs, groups, identities, telemetry) | `ROUTERS = [...]` list | `letta/server/rest_api/routers/v1/__init__.py:34-68` |
| Mounting at `/v1` plus undocumented `/latest` alias | `app.include_router(route, prefix=API_PREFIX)` and `include_in_schema=False` alias | `letta/server/rest_api/app.py:852-857` |
| Admin users/orgs routers under `/v1/admin` | `app.include_router(users_router, prefix=ADMIN_PREFIX)` | `letta/server/rest_api/app.py:859-861` |
| Auth router mounted at `/api/auth` style prefix | `app.include_router(setup_auth_router(...), prefix=API_PREFIX)` | `letta/server/rest_api/app.py:863-864` |
| Default port | `REST_DEFAULT_PORT = 8283` | `letta/server/constants.py:6` |
| OpenAI-compatible chat completions route | `@router.post("/chat/completions", ..., operation_id="create_chat_completion")` | `letta/server/rest_api/routers/v1/chat_completions.py:117-124` |
| OpenAPI schema post-processing for public docs (strip `/openai` paths, inject Letta message unions) | `generate_openapi_schema()` filters paths and adds `LettaMessageUnion` etc. | `letta/server/rest_api/app.py:136-162` |
| Schema generation script (manual, writes `openapi_letta.json` to CWD) | shell script invoking `generate_openapi_schema(app)` | `letta/server/generate_openapi_schema.sh:11-12`; writer at `letta/server/rest_api/app.py:162` |
| Checked-in API contract consumed by Fern | OpenAPI 3.1 doc with `servers`: api.letta.com + localhost:8283 | `fern/openapi.json:1-20` |
| SDK surface curation (hide MCP endpoints from SDKs) | `x-fern-ignore: true` entries | `fern/openapi-overrides.yml:11-45` |
| SDK publish pipeline (Python) | `fern generate --group python-sdk ...` + docs publish | `.github/workflows/fern-sdk-python-publish.yml:39-50` |
| Spec validity check on PRs | `fern check` job on pull_request to main | `.github/workflows/fern-check.yml:18-20` |
| SDK preview triggered by spec changes | path filter on `fern/openapi.json` / `openapi-overrides.yml` | `.github/workflows/fern-sdk-python-preview.yml:15-16` |
| Server depends on generated SDK (dogfooding) | `letta-client>=1.7.12` dependency | `pyproject.toml:46` |
| Python library re-export surface | `from letta.schemas.agent import AgentState as AgentState` etc. (23 schema exports) | `letta/__init__.py:24-47` |
| Deprecation markers on routes/params (flow into OpenAPI → SDKs) | e.g. `deprecated=True` on attach/detach source routes and query params | `letta/server/rest_api/routers/v1/agents.py:819`, `agents.py:875`, `agents.py:298-304`, `agents.py:588` |
| Deprecated telemetry route with migration pointer | `retrieve_provider_trace` marked deprecated; docstring says use `GET /steps/{step_id}/trace` | `letta/server/rest_api/routers/v1/telemetry.py:12-19` |
| Internal endpoints namespaced but still public | `router = APIRouter(prefix="/_internal_agents", tags=["_internal_agents"])`; included in `ROUTERS` | `letta/server/rest_api/routers/v1/internal_agents.py:8`; `routers/v1/__init__.py:47-51` |
| Git HTTP router hidden from schema | `include_in_schema=False` on git router | `letta/server/rest_api/routers/v1/git_http.py:48` |
| Error-to-HTTP-code mapping as part of the contract (400/404/409/503 etc.) | exception handler registrations | `letta/server/rest_api/app.py:556-611` |
| Friendly 422 ID-validation errors with examples (contract polish) | `custom_request_validation_handler` rewriting `*_id` pattern errors | `letta/server/rest_api/app.js:477-532` → actually `letta/server/rest_api/app.py:477-532` |
| SDK contract tests against live server | session fixture boots server, polls `/v1/health`, creates `Letta(base_url=...)` client | `tests/sdk/conftest.py:12-56` |
| Generated-CRUD test harness over SDK resources | `create_test_module(resource_name="agents", ...)` | `tests/sdk/conftest.py:64-256`; usage `tests/sdk/agents_test.py:51-60` |
| Integration tests consume the public SDK types | `from letta_client import AsyncLetta, Letta` and streaming union types | `tests/integration_test_send_message_v2.py:15-25`; `tests/conftest.py:12` |
| README quickstart is runnable SDK-first (Python + TS) | Hello World using `client.agents.create` / `client.agents.messages.create` | `README.md:36-110` |
| In-repo examples directory nearly empty (only data assets) | `examples/notebooks/data` contains pdf/txt files only | `examples/notebooks/data/` |
| Dead in-tree client module (no importers) | empty `letta/client/__init__.py`; `_sse_post` referenced only from a commented-out test | `letta/client/__init__.py:1`; `letta/client/streaming.py:19`; sole reference commented at `tests/locust_test.py` (`from letta.client.streaming import _sse_post` is commented) |
| Tool sandbox injects SDK usage into user-facing tool code | sandbox preamble prints `from letta_client import Letta` guidance | `letta/services/tool_sandbox/base.py:220-244`; `letta/services/tool_executor/builtin_tool_executor.py:77` |
| Version fallback drift between packaging and module | hardcoded `"0.16.7"` vs `version = "0.16.8"` | `letta/__init__.py:8` vs `pyproject.toml:3` |
| Stray test file ships inside the wheel package | `packages = ["letta"]` includes `letta/test_gemini.py` | `pyproject.toml:167-168`; `letta/test_gemini.py:1` |

## Answers to Dimension Questions

### 1. What is the intended public API surface?

Three tiers, each verifiable in code:

- **Primary: REST + generated SDKs.** The FastAPI app mounts 33 resource routers under `/v1` (`letta/server/rest_api/routers/v1/__init__.py:34-68`, mounting at `letta/server/rest_api/app.py:852-857`) and publishes an OpenAPI document (`generate_openapi_schema`, `letta/server/rest_api/app.py:136-162`). That document is committed at `fern/openapi.json` and is the source for both language SDKs, built and published through Fern workflows (`.github/workflows/fern-sdk-python-publish.yml:39-50`). The README positions the SDKs as the developer entry point (`README.md:24-34`), and all integration tests drive the server through `letta_client` rather than raw HTTP (`tests/integration_test_send_message_v2.py:15-25`).
- **Secondary: CLI.** Reduced to a server launcher: `letta` with no args starts the server (`letta/main.py:12-16`); options cover port/host/debug/reload/security (`letta/cli/cli.py:17-26`).
- **Tertiary: embedded/library use.** `letta/__init__.py:24-47` exposes schema models, and the tool-sandbox system actively teaches agent-executed tool code to use the hosted SDK (`letta/services/tool_sandbox/base.py:220`).

### 2. Is the stable API easy to distinguish from internal implementation details?

Mostly yes, by convention and tooling rather than by enforcement:

- Versioning in the URL path (`/v1`, `letta/constants.py:33`) plus an undocumented `/latest` alias explicitly tied to the newest API (`letta/server/rest_api/app.py:854-857`).
- FastAPI-native deprecation: `deprecated=True` appears on routes and parameters throughout `letta/server/rest_api/routers/v1/agents.py` (e.g., lines 151, 164, 169, 588, 819, 875, 1079, 1206) and `telemetry.py:12`, which surfaces in OpenAPI and therefore in generated SDKs.
- Internal endpoints carry a `_internal` naming convention (`internal_agents.py:8` uses prefix `/_internal_agents` and tag `_internal_agents`; similarly `internal_blocks.py`, `internal_runs.py`, `internal_search.py`, `internal_templates.py` per `routers/v1/__init__.py:14-17,48-51`).
- Fern overrides remove specific MCP-server endpoints from the generated SDKs while leaving them reachable over HTTP (`fern/openapi-overrides.yml:11-45`) — a deliberate split between wire availability and supported SDK surface.
- Gaps: `_internal_*` routes are still mounted publicly and documented in the schema (they sit inside the same `ROUTERS` list, `routers/v1/__init__.py:47-51`); the git HTTP router is hidden via `include_in_schema=False` (`git_http.py:48`) but remains callable. Naming, not access control, is the boundary.

### 3. Does the API expose the right level of abstraction for agent harness users?

Yes — the surface is resource-oriented and hides runtime internals:

- Resources map to harness concepts: agents, conversations, blocks (memory), tools, sources/folders/archives (files/RAG), runs/steps/jobs, groups (multi-agent), identities, tags (`routers/v1/__init__.py:34-68`).
- Streaming is exposed as typed message unions (`LettaMessageUnion`, content unions injected into the published schema at `letta/server/rest_api/app.py:144-148`; schema definitions in `letta/schemas/letta_message.py:16-49` including `ApprovalReturn`/`ToolReturn` for human-in-the-loop flows).
- Provider plumbing (LLM clients, adapters, ORM managers) stays behind the server boundary; no router returns ORM objects — responses are Pydantic schemas like `AgentState` (`letta/schemas/agent.py:67`) and `Block` (`letta/schemas/block.py:67`).
- An OpenAI-compatible facade (`POST /v1/chat/completions`, `chat_completions.py:117-124`) lets standard OpenAI clients integrate without learning Letta semantics.
- Error semantics are part of the abstraction: domain exceptions map to specific HTTP codes with retry hints (`Retry-After` on lock/deadlock handlers, `letta/server/rest_api/app.py:613-641`), and LLM-provider failures get dedicated codes (402/429/504, `app.py:692-741`).

### 4. Are examples sufficient to use the API correctly without reading internals?

Adequate but thin inside this repository:

- The README quickstart provides complete, runnable Python and TypeScript snippets covering the two central calls (create agent, send message) (`README.md:36-110`), and points to external docs (`docs.letta.com`).
- Route handler docstrings serve as inline reference documentation (e.g., `internal_agents.py:17-19`), and the OpenAPI/Swagger UI is served with collapsed expansion for browsing (`app.py:412`).
- However, the in-repo `examples/` tree contains only data assets (a PDF and prompt text files under `examples/notebooks/data/`), not code. Correctness beyond the quickstart currently requires either the external docs site or reading the integration tests — which, to be fair, are extensive SDK-driven scenarios (`tests/sdk/`, `tests/integration_test_send_message_v2.py`).

## Architectural Decisions

1. **Server repo, SDK elsewhere, spec as the bridge.** The `letta` package does not ship an HTTP client; the generated `letta-client` is a first-class dependency of the server itself (`pyproject.toml:46`). This forces every API change through the OpenAPI contract and makes the SDK the canonical consumption path — even Letta's own tests cannot bypass it (`tests/conftest.py:12`).
2. **Contract checked in and CI-enforced for validity.** `fern/openapi.json` lives next to the code, PRs run `fern check` (`.github/workflows/fern-check.yml:18-20`), and changes to the spec trigger preview SDK builds that other test jobs block on (`.github/workflows/reusable-test-workflow.yml:110-118,127-131`).
3. **Versioned URL scheme with a moving alias.** `/v1` is stable; `/latest` is an undocumented pointer always bound to the current API (`app.py:854-857`), giving operators a migration-friendly upgrade path.
4. **Deprecation through the schema, not changelogs.** Using FastAPI's `deprecated=True` on operations and parameters means deprecation status propagates mechanically to OpenAPI, docs, and SDK type stubs (e.g., `agents.py:819`, `telemetry.py:12`).
5. **Curated SDK surface distinct from wire surface.** `x-fern-ignore` removes endpoints from SDKs without removing them from the server (`fern/openapi-overrides.yml:11-45`); conversely `/latest` duplicates exist on the wire but never in the schema (`app.py:857`).
6. **Errors as API.** A large, explicit table maps internal exception classes to HTTP status codes with consistent JSON bodies and retry headers (`app.py:556-641`), making failure behavior a designed part of the public contract rather than an emergent property.

## Notable Patterns

- **`operation_id` discipline**: 207 of 252 route decorators specify an explicit `operation_id` (survey across `letta/server/rest_api/routers/v1/*.py`), which stabilizes SDK method names regardless of refactoring (e.g., `create_agent`, `attach_tool_to_agent` in `agents.py:613,670`).
- **Schema-generated union injection**: streaming/message polymorphism is assembled dynamically and spliced into the OpenAPI components before publishing (`create_letta_message_union_schema` at `app.py:144-149`), so SDK consumers see proper discriminated unions instead of opaque dicts.
- **Contract-test factory**: `tests/sdk/conftest.py:64-256` defines `create_test_module()`, a generic CRUD test generator parameterized per resource; each resource test file is just data (`tests/sdk/agents_test.py:51-60`). Adding API resources gets lifecycle tests almost for free.
- **Self-describing error messages**: the 422 handler rewrites cryptic UUID-pattern validation failures into friendly messages with example IDs (`app.py:477-532`).
- **SDK-aware sandboxing**: the tool execution sandbox embeds `letta_client` import snippets into generated tool code so agent-authored tools can call back into the platform (`letta/services/tool_sandbox/base.py:220-244`).

## Tradeoffs

- **Spec-sync is manual.** Regenerating `openapi_letta.json` writes to the process CWD (`app.py:162`) via a hand-run script (`letta/server/generate_openapi_schema.sh:11-12`); nothing in CI diffs it against `fern/openapi.json`. `fern check` validates the spec's internal consistency, not its fidelity to the FastAPI app. Drift between implemented routes and the published contract is possible until a human regenerates and commits.
- **Convention-over-enforcement for internals.** `_internal_*` naming communicates intent but doesn't restrict access or exclude the routes from the published schema; a determined external user can bind to them and be broken by changes.
- **CLI minimalism trades discoverability for simplicity.** Everything except running the server moved to SDKs/ADE; there is no `letta version` command registered even though the function exists (`cli.py:45-48`, unused in `main.py`).
- **Two sources of truth for types.** Server-side Pydantic schemas (`letta/schemas/*`) and SDK-side Fern-generated types must stay semantically aligned; the union-injection step (`app.py:144-148`) is custom glue that must be maintained when schemas change.
- **Examples live outside the repo.** Keeps the repo lean but weakens the "clone and verify" story; the strongest executable documentation here is the test suite, which is not framed as documentation.

## Failure Modes / Edge Cases

- **Stale contract**: a route added in FastAPI but with `fern/openapi.json` left un-regenerated will simply never appear in any SDK; consumers discover the feature late or reimplement raw HTTP calls, undermining the whole SDK-first model. No automated guard exists (searched `.github/workflows/` for openapi diff/regeneration jobs; only fern check/preview/publish found).
- **Dead importable modules**: `letta/client/streaming.py:19` (`_sse_post`) and `letta/client/utils.py` have no remaining importers in the package (only a commented reference in `tests/locust_test.py`), yet remain public-importable from the installed wheel — accidental surface that could rot silently.
- **Packaging leaks**: `letta/test_gemini.py` sits inside the shipped package (`pyproject.toml:167-168` sets `packages = ["letta"]`), so a stray test module is distributed to every user.
- **Version ambiguity**: `letta/__init__.py:8` falls back to a hardcoded `0.16.7` while the project version is `0.16.8` (`pyproject.toml:3`), and `LETTA_VERSION` env can override both (`__init__.py:10-11`) — minor, but version reporting is not single-sourced.
- **Undocumented-but-live aliases**: `/latest` (`app.py:857`) and the git HTTP router (`git_http.py:48`, `include_in_schema=False`) are functional but absent from the published spec; clients coding against them have no contract protection.
- **Deprecated-parameter accumulation**: many query parameters are individually marked deprecated while still honored (`agents.py:492-545` shows dual-path parsing of old/new field names), creating long-lived compatibility branches inside handlers that must be tested on both paths.

## Future Considerations

- Add a CI job that boots the app, runs `generate_openapi_schema`, and fails if the output differs from `fern/openapi.json` — converting the manual sync step into a verified invariant.
- Exclude `_internal_*` routers from the published schema (or move them behind `include_in_schema=False` / auth gating) so the supported surface and the documented surface coincide.
- Delete or privatize `letta/client/` (empty `__init__`, dead `streaming.py`/`utils.py`) and move `letta/test_gemini.py` into `tests/`.
- Register a `letta version` command (function already exists at `cli.py:45`) or drop the function.
- Add a small set of runnable examples under `examples/` mirroring the README quickstart (streaming, tool approval, multi-agent group) so SDK behavior can be verified locally without the external docs site.

## Questions / Gaps

- No evidence found in-repo for how `fern/openapi.json` is regenerated and committed in practice (who copies `openapi_letta.json` into `fern/`, and when). Searched `.github/workflows/`, `scripts/`, and `letta/server/generate_openapi_schema.sh`; the script writes to CWD only, and no workflow performs regeneration against the live app.
- The `ws_api` directory listed under `letta/server/` was not examined further because the CLI already hard-fails WebSocket mode (`cli.py:42`); whether any ws code remains reachable is unclear from the entry points alone.
- Rate limiting, quotas, and auth token issuance appear handled by the auth router (`app.py:863-864`) and cloud layers, but the boundary between OSS-server auth and Letta Cloud auth is not documented in-repo; I did not deep-dive `rest_api/auth/` for this dimension.
- Whether `letta/openai_backcompat/openai_object.py` is still exercised by any route could not be confirmed; no router references were found in the surveyed surface.

---

Generated by dimension 24.01 (Public API Surface) against `letta`.
