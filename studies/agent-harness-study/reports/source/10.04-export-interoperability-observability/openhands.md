# Source Analysis: openhands

## Dimension 10.04: Export, Interoperability, and Observability Backends

### Source Info

| Field | Value |
|-------|-------|
| Name | openhands |
| Path | `studies/agent-harness-study/sources/openhands` |
| Language / Stack | Python (FastAPI, SDK via `openhands-sdk`/`openhands-agent-server`, React frontend) / `lmnr` + `opentelemetry-*` |
| Analyzed | 2026-08-28 |

## Summary

OpenHands trace export is **single-vendor-centric with a generic OTLP escape hatch**. The SDK (`_sdk_inspect/sdk/observability/laminar.py`) wraps [Traceloop/Laminar](https://github.com/traceloop/openllmetry) via `lmnr>=0.7.20` and `opentelemetry-api`/`opentelemetry-exporter-otlp-proto-grpc==1.39.1`. Initialization is env-driven: `LMNR_PROJECT_API_KEY` targets Laminar Cloud directly, while the standard OTel variables `OTEL_ENDPOINT` / `OTEL_EXPORTER_OTLP_ENDPOINT` / `OTEL_EXPORTER_OTLP_TRACES_ENDPOINT` + `OTEL_EXPORTER_OTLP_TRACES_HEADERS`/`OTEL_EXPORTER_OTLP_TRACES_PROTOCOL`/`OTEL_EXPORTER` target any OTLP-compatible collector. All `LMNR_*` vars are auto-forwarded from the app-server host into the agent-server sandbox via `get_agent_server_env()`; generic `OTEL_*` vars require manual forwarding through `OH_AGENT_SERVER_ENV` JSON — an undocumented asymmetry.

There is no pluggable sink interface, no first-class Langfuse/LangSmith/Honeycomb SDK, no multi-sink fan-out, and no OTLP configuration in `config.template.toml`. Local artifact export exists orthogonally: conversation trajectories are streamed as a ZIP (`meta.json` + `event_*.json`) through `LiveStatusAppConversationService`, protected by a Redis lock and a configurable event-count limit — a local-file interoperability path, but not an observability pipeline.

Runtime reconfiguration is env-only (restart required); the trace payload mixes standard OTel + `opentelemetry-semantic-conventions-ai` attributes with custom OpenHands metadata (`conversation_id`, `repo`, `branch`, `commit`, `user_id` resolved from email) injected both as Laminar span metadata and as `litellm_extra_body.metadata` for `openhands/` models.

## Rating

**5 / 10 — Present but inconsistent, weakly documented, and fragile**

OTLP proto/gRPC export is a hard dependency and the generic OTel env contract actually works (any collector-compatible backend like Honeycomb/Grafana/Tempo can ingest it), but the architecture funnels everything through a single `Laminar.initialize()` call, auto-forwards only the `LMNR_*` prefix, provides no custom-sink SPI, no multi-backend dispatch, and no TOML/config-file knob. Tests cover only env-forwarding mechanics, not end-to-end trace emission; documentation lives in a docstring and an inline comment, not in a user-facing integration guide.

## Evidence Collected

Every entry MUST include a file path with line numbers. Format: `path/to/file.ts:NN`.

| Area | Evidence | File:Line |
|------|----------|-----------|
| OTLP dependency declaration | `opentelemetry-api==1.39.1` and `opentelemetry-exporter-otlp-proto-grpc==1.39.1` pinned as hard runtime deps (also `pyproject.toml:64-65`, `poetry` duplicate `pyproject.toml:196-197`) | `pyproject.toml:64-65` |
| Laminar SDK dependency | `lmnr>=0.7.20` required for observability | `pyproject.toml:56` |
| OTel endpoint env contract | `_OBSERVABILITY_ENV_KEYS = ("LMNR_PROJECT_API_KEY","OTEL_ENDPOINT","OTEL_EXPORTER_OTLP_TRACES_ENDPOINT","OTEL_EXPORTER_OTLP_ENDPOINT")` — anything matching enables observability | `_sdk_inspect/sdk/observability/laminar.py:25-30` |
| OTel protocol/header env docs | Docstring for `maybe_init_laminar()` documents `OTEL_EXPORTER_OTLP_TRACES_HEADERS`, `OTEL_EXPORTER_OTLP_TRACES_PROTOCOL` (`http/protobuf` vs `grpc`), `OTEL_EXPORTER` (`otlp_http`/`otlp_grpc`) | `_sdk_inspect/sdk/observability/laminar.py:57-82` |
| Generic OTel init path | When not Laminar, calls `Laminar.initialize(disabled_instruments=[BROWSER_USE_SESSION,PATCHRIGHT,PLAYWRIGHT], force_http=...)` — still routes through OTel exporter | `_sdk_inspect/sdk/observability/laminar.py:102-112` |
| Laminar-specific init | When `_is_otel_backend_laminar()` (checks `LMNR_PROJECT_API_KEY`), calls `Laminar.initialize(base_url, http_port, grpc_port, force_http)` with `LMNR_BASE_URL`/`LMNR_HTTP_PORT`/`LMNR_GRPC_PORT` | `_sdk_inspect/sdk/observability/laminar.py:95-101` |
| Idempotent observability gate | `should_enable_observability()` caches positive via `_observability_enabled`, checks `any(get_env(key) for key in _OBSERVABILITY_ENV_KEYS)` and `Laminar.is_initialized()` | `_sdk_inspect/sdk/observability/laminar.py:199-215` |
| Env reader | `get_env()` merges `os.getenv` + `dotenv_values()` so `.env` files participate | `_sdk_inspect/sdk/observability/utils.py:8-10` |
| Root span lifecycle | `RootSpan` wraps `Laminar.start_span(name)` + `Laminar.set_trace_session_id(session_id)` inside `Laminar.use_span`; `start_root_span`/`end_root_span` are idempotent and swallowed on failure | `_sdk_inspect/sdk/observability/laminar.py:231-296` |
| Per-conversation ownership | `BaseConversation._observability_root_span: RootSpan\|None` looked up by `observe` decorator via `_ROOT_SPAN_ATTR = "_observability_root_span"` and re-attached via `Laminar.use_span` to survive task/thread hops | `_sdk_inspect/sdk/observability/laminar.py:228-330` and ` _sdk_inspect/sdk/conversation/base.py:120-151` |
| Decorator lazy import | `observe` wrapper defers `from lmnr import observe` until first call after enablement; preserves `iscoroutinefunction` branching | `_sdk_inspect/sdk/observability/laminar.py:115-196` |
| External webhook entry | `init_laminar_for_external()` calls `maybe_init_laminar()` then returns `Laminar.get_laminar_span_context()` for parent propagation | `_sdk_inspect/sdk/observability/laminar.py:498-503` |
| LMNR auto-forward | `AUTO_FORWARD_PREFIXES = ('LLM_','LMNR_')` — only those prefixes are auto-injected into sandbox `initial_env` | `openhands/app_server/sandbox/sandbox_spec_service.py:151` |
| Forwarding implementation | `get_agent_server_env()` scans `os.environ` for prefix matches then merges `OH_AGENT_SERVER_ENV` JSON (explicit overrides win) | `openhands/app_server/sandbox/sandbox_spec_service.py:198-209` |
| `OH_AGENT_SERVER_ENV` explicit override | All sandbox spec factories merge `get_agent_server_env()` into `initial_env`; tests verify override precedence | `tests/unit/app_server/test_agent_server_env_override.py:241-268` |
| OTEL asymmetry | No `OTEL_` entry in `AUTO_FORWARD_PREFIXES`; generic collector therefore requires `OH_AGENT_SERVER_ENV='{"OTEL_EXPORTER_OTLP_TRACES_ENDPOINT":"..."}'` — not mentioned in `sandbox_spec_service.py:149-210` docstring | `openhands/app_server/sandbox/sandbox_spec_service.py:151` vs `_sdk_inspect/sdk/observability/laminar.py:25-30` |
| LMNR forwarding tests | Dedicated `TestLMNRAutoForwarding` asserts `LMNR_PROJECT_API_KEY`, `LMNR_BASE_URL`, etc. are present | `tests/unit/app_server/test_agent_server_env_override.py:298-357` |
| Docker/Process/Remote spec propagation tests | All three spec types verified to carry `initial_env` from `get_agent_server_env()` | `tests/unit/app_server/test_agent_server_env_override.py:360-558` |
| LLM tracing metadata | `get_llm_metadata()` builds `tags=["app:openhands","model:...","type:...","web_host:...","openhands_version:...","conversation_version:V1"]` + `session_id`/`trace_user_id`/`repo_name`/`git_provider`/`selected_branch` | `openhands/app_server/utils/llm_metadata.py:28-91` |
| Server-side metadata enrichment | `_build_observability_context()` and `_build_observability_metadata()` assemble `repo`/`branch`/`commit` (from `git rev-parse HEAD` in workspace, 10s timeout) + `agent_kind` | `openhands/app_server/app_conversation/live_status_app_conversation_service.py:1587-1867` |
| Condenser tracing | Condenser LLM gets its own metadata with `usage_id="condenser"` gated on `should_set_litellm_extra_body("openhands/...")` | `openhands/app_server/app_conversation/live_status_app_conversation_service.py:1681-1703` |
| LiteLLM injection | Resolved `llm_metadata` injected as `agent.llm.litellm_extra_body={"metadata":...}` only for `model.startswith("openhands/")` | `openhands/app_server/utils/llm_metadata.py:25` and `openhands/app_server/app_conversation/live_status_app_conversation_service.py:1666-1678` |
| User attribution | `laminar_user_id = await self.user_context.get_user_email() or user.id` forwarded as `user_id` / `trace_user_id` | `openhands/app_server/app_conversation/live_status_app_conversation_service.py:2170,2441` and `openhands/app_server/utils/llm_metadata.py:78-79` |
| Local file export (ZIP) | `_stream_conversation_zip()` writes `meta.json` + per-event `event_000000_<id>.json` (model_dump_json) via streaming `zipfile.ZipFile` over `_StreamingZipBuffer` | `openhands/app_server/app_conversation/live_status_app_conversation_service.py:2869-2893` |
| Export contract | Abstract `export_conversation(conversation_id)->bytes` and streaming `open_conversation_export()->AsyncGenerator[bytes]` | `openhands/app_server/app_conversation/app_conversation_service.py:160-189` |
| Export HTTP endpoint | `GET` (router) returns `StreamingResponse(zip_stream, media_type="application/zip", Content-Disposition: attachment; filename="conversation_{id}.zip")` | `openhands/app_server/app_conversation/app_conversation_router.py:1619-1683` |
| Export guards | `export_max_events=10000`, `export_lock_ttl_seconds=3600`, `export_lock_refresh_interval_seconds=30`, `export_lock_required` (required in `saas` mode) | `openhands/app_server/app_conversation/live_status_app_conversation_service.py:284-287,3006-3024` |
| Redis locking | `_EXPORT_LOCK_KEY_PREFIX="app_conversation_export"` + `try_acquire_redis_lock`/`refresh_lock_periodically` with best-effort fallback outside SaaS | `openhands/app_server/app_conversation/live_status_app_conversation_service.py:153,2895-2962` |
| No custom sink SPI | Grep for `TraceExporter`/`SpanExporter`/`custom.*sink` yields 0 hits in `openhands/` Python | `grep(openhands, "TraceExporter\|SpanExporter\|custom.*sink"): no matches` |
| No Langfuse/LangSmith/Honeycomb SDK | `grep(Langfuse|LangSmith|Honeycomb)` only hits unrelated `poetry.lock` vendored entry and a `WEB_HOST` comment; no integration code | `grep(openhands, "Langfuse\|LangSmith\|Honeycomb"): 0 product hits` |
| No multi-backend fan-out | Only one `Laminar.initialize()` + one `litellm.callbacks.append(LaminarLiteLLMCallback())`; no loop over exporters | `_sdk_inspect/sdk/observability/laminar.py:89-112` |
| No TOML export config | `config.template.toml` contains `[core]`/`[agent]`/`[sandbox]`/`[security]`/`[condenser]`/`[kubernetes]`/`[mcp]`/`[model_routing]` — no `[observability]`/`[tracing]` section | `config.template.toml:1-394` |
| Log aggregation vs trace export | Structured JSON logging (`LOG_JSON`, `pythonjsonlogger`, `get_uvicorn_log_config` with `RedactURLParamsFilter`) is a log pipeline to Datadog-style collectors, distinct from OTel traces | `openhands/app_server/utils/logger.py:336-353,582-748` |

## Answers to Dimension Questions

**1. Can traces be exported to external backends?**
Yes — via two disjoint paths. (a) Distributed tracing: if any of `LMNR_PROJECT_API_KEY`, `OTEL_ENDPOINT`, `OTEL_EXPORTER_OTLP_TRACES_ENDPOINT`, `OTEL_EXPORTER_OTLP_ENDPOINT` is set, `maybe_init_laminar()` (`_sdk_inspect/sdk/observability/laminar.py:57`) initializes Laminar's OTel pipeline and appends `LaminarLiteLLMCallback` to `litellm.callbacks`, so LLM calls and `@observe`-decorated functions emit OTel spans to the configured collector. (b) Artifact export: any conversation can be downloaded as a ZIP of JSON events (`openhands/app_server/app_conversation/live_status_app_conversation_service.py:2869`, `app_conversation_router.py:1619`). What is missing is an in-process pluggable sink; you cannot register a Python `SpanExporter` or a file sink without forking the SDK.

**2. Are standard protocols supported?**
Partially. The stack is OTLP-native: `opentelemetry-exporter-otlp-proto-grpc==1.39.1` (`pyproject.toml:65`) plus runtime env `OTEL_EXPORTER_OTLP_TRACES_PROTOCOL=http/protobuf|grpc` and `OTEL_EXPORTER=otlp_http|otlp_grpc` (`laminar.py:69-71`) covers both gRPC and HTTP/Protobuf transports. `opentelemetry-semantic-conventions-ai==0.4.13` is also pulled via the `lmnr` transitive graph, so gen-AI conventions are available. However, only OTLP is first-class — there is no Jaeger-native, Zipkin, or Prometheus exposition, and vendor-specific conventions (Honeycomb `hny.*`, Langfuse ingestion) are not mapped; they would rely on the OTel collector to translate.

**3. Is export configurable without code changes?**
Yes for the single-backend case, env-only. Set `LMNR_PROJECT_API_KEY` (and optionally `LMNR_BASE_URL`/`LMNR_HTTP_PORT`/`LMNR_GRPC_PORT`/`LMNR_FORCE_HTTP`) for Laminar Cloud, or set `OTEL_EXPORTER_OTLP_TRACES_ENDPOINT` + `OTEL_EXPORTER_OTLP_TRACES_HEADERS` for any OTLP collector — no code change. `LLM_*` and `LMNR_*` are auto-forwarded into the agent-server sandbox; generic `OTEL_*` requires `OH_AGENT_SERVER_ENV='{"OTEL_EXPORTER_OTLP_TRACES_ENDPOINT":"..."}'` (`sandbox_spec_service.py:151,198`). Negatives: no `config.template.toml` knob, no config-file or UI toggle, `OTEL_*` auto-forward omission is a papercut requiring manual JSON, and changes require a process restart (no hot-reload).

**4. Can multiple backends receive traces simultaneously?**
No. Evidence is negative: the codebase contains exactly one `Laminar.initialize()` site and one `litellm.callbacks.append(...)` (`laminar.py:95-112`), no exporter registry/loop, no `SpanProcessor` multiplexing, and no `grep` hit for fan-out. To dual-ship (e.g., Laminar + Honeycomb) you must front the OTel endpoint with an OpenTelemetry Collector and configure its exporters externally — not a product feature. Local ZIP export is independent and does not count as simultaneous trace fan-out.

## Architectural Decisions

- **Single-vendor default with generic OTLP fallback.** Decision to hard-pin `lmnr` + `opentelemetry-exporter-otlp-proto-grpc` and gate everything behind `_OBSERVABILITY_ENV_KEYS` (`laminar.py:25`) gives a zero-config path to Laminar Cloud while still admitting any OTLP collector via the same SDK. Tradeoff: vendor lock-in to Laminar's `observe`/`Laminar` API surface rather than a vendor-neutral `opentelemetry.sdk.trace` setup.
- **Per-conversation `RootSpan` re-attached via `Laminar.use_span`.** Replaces the deprecated global-LIFO `SpanManager` (`laminar.py:351-397`) after empirical ~60% trace-context loss with `start_active_span` across asyncio tasks (`laminar.py:238-249`). Each `BaseConversation` owns `self._observability_root_span` (`sdk/conversation/base.py:128`) and the `observe` decorator re-enters it on every call (`laminar.py:299-331`). Durably fixes cross-conversation interleaving but ties the entire trace model to `lmnr`'s span impl.
- **Env-forwarding via prefix allowlist.** `AUTO_FORWARD_PREFIXES=('LLM_','LMNR_')` (`sandbox_spec_service.py:151`) auto-propogates observability vars into the isolated agent-server container without exposing the full host env. The allowlist is minimal by design; the omission of `OTEL_` forces an explicit `OH_AGENT_SERVER_ENV` JSON allowlist-bypass for generic collectors — a security-vs-ergonomics tradeoff.
- **Tracing metadata dual-write (Laminar span + LiteLLM proxy).** Same logical context (`conversation_id`→`session_id`, `repo`/`branch`/`commit`/`git_provider`, `user_id` from email, `trace_version`) is written both to Laminar (`_build_observability_context`/`_build_observability_metadata`) and to LiteLLM's `extra_body.metadata` for `openhands/` models (`llm_metadata.py:28`, `live_status_app_conversation_service.py:1666`). Enables SaaS analytics without requiring the LLM provider to understand OTel.
- **Trajectory ZIP as the interoperability artifact.** Rather than emitting OTLP file exports, the system streams `meta.json` + canonical `Event` JSON files (`event_service_base.py:145`) into a ZIP on demand, with streaming chunking (`_StreamingZipBuffer`) to avoid OOM (`live_status_app_conversation_service.py:172-193,2869`). Serves debugging/replay/audit, not real-time observability.

## Notable Patterns

- **Lazy-import `observe` decorator** (`laminar.py:115-196`) — avoids importing `lmnr` until `should_enable_observability()` is true, preserving cold-start for users who don't use tracing and preserving `iscoroutinefunction` semantics by branching at decoration time.
- **Dotenv + env merge** (`observability/utils.py:8`) — `get_env()` checks `os.getenv` first, then `dotenv_values()`, so local `.env` files can enable tracing without shell export.
- **Idempotent global gate** (`laminar.py:199`) — `_observability_enabled` is sticky; once any key appears or `Laminar.is_initialized()`, observability stays on for the process lifetime, avoiding late-disable races.
- **Attribute-based implicit span propagation** — `observe` discovers the owning `RootSpan` by reflecting on `args[0]._observability_root_span` (`laminar.py:324-330`), a convention-based coupling that avoids passing context explicitly but is invisible to static analysis.
- **Head-commit best-effort enrichment** (`live_status_app_conversation_service.py:1819-1838`) — `git rev-parse --verify --quiet HEAD` with 10s timeout; failure is swallowed (`return ''`) so slow clones never block conversation start.

## Tradeoffs

- **Ergonomics vs vendor neutrality.** Using `lmnr observe` + `RootSpan` + `Laminar.use_span` gives rich auto-instrumentation for LiteLLM/MCP with minimal code, but couples the SDK to Traceloop's SDK version (`lmnr>=0.7.20`) and prevents swapping in a plain OTel `TracerProvider` or adding a second vendor span without forking.
- **Single sink simplicity vs operational flexibility.** One global exporter keeps config surface tiny (4 env keys) but precludes dual-shipping, per-environment routing (dev→local collector, prod→SaaS), or sampling policies without deploying an external collector as a proxy.
- **Prefix allowlist security vs OTLP ergonomics.** Allowlisting only `LLM_`/`LMNR_` limits blast radius of host env leakage into sandboxes, but makes the primary generic-tracing workflow (`OTEL_*`) second-class and undiscoverable (no mention in `sandbox_spec_service.py` docstring beyond `LMNR_*` examples).
- **LiteLLM metadata breadth vs provider compatibility.** Gating `litellm_extra_body` to `openhands/*` models (`llm_metadata.py:25`) avoids breaking non-OpenHands providers that reject `extra_body`, but means non-`openhands/` models lose SaaS correlation tags even when tracing is on.
- **ZIP artifact completeness vs streaming cost.** Zipping canonical events (`model_dump(mode='json')`, `indent=2`) preserves fidelity and is streaming-safe, but for large conversations (≥10k events) hits `ConversationExportTooLarge` (`live_status_app_conversation_service.py:2850`) and forces consumers to reassemble traces from events rather than from native OTel data.

## Failure Modes / Edge Cases

- **Silent no-op when envs absent.** `maybe_init_laminar()` logs at `DEBUG` and returns (`laminar.py:83-87`); decorators degrade to pass-through (`laminar.py:174-175,188-189`). Failure to configure tracing is undetectable at `INFO` — easy to ship with tracing “on” in config but inert in runtime.
- **OTEL vars not auto-forwarded.** A user setting `OTEL_EXPORTER_OTLP_TRACES_ENDPOINT` on the host but not duplicating it into `OH_AGENT_SERVER_ENV` will see traces from the app-server process (if it also calls `maybe_init_laminar`) but not from the agent-server sandbox where the actual LLM spans live. No warning emitted.
- **Header encoding fragility.** `OTEL_EXPORTER_OTLP_TRACES_HEADERS` is documented as comma-separated URL-encoded `key=value%20pairs` (`laminar.py:66`); a mis-encoded header will be passed to the OTel exporter and fail at the collector with no SDK-side validation.
- **Sticky enablement prevents graceful disable.** Once `_observability_enabled` flips to `True`, it never flips back (`laminar.py:201-202`); unsetting env vars at runtime does not stop tracing until process restart, and a transient env injection permanently enables overhead.
- **`gRPC` default vs network policy.** Default `OTEL_EXPORTER_OTLP_TRACES_PROTOCOL` is gRPC (`grpc/protobuf`) when unspecified (`laminar.py:68`); environments that only allow HTTP egress (many corp proxies) will time out unless explicitly switched to `http/protobuf`.
- **Span leak on abrupt close.** `_end_observability_span()` is idempotent via `_span_ended`, but if the process crashes or `end_root_span` swallows an exception (`laminar.py:274-275` logs at `DEBUG` only), the trace remains open server-side until collector timeout, skewing duration.
- **Redis lock unavailability modes.** In OSS (`app_mode != "saas"`), lock acquisition failure degrades to “proceed without lock” with a warning (`live_status_app_conversation_service.py:2911-2916`); in SaaS it hard-fails with `ConversationExportLockUnavailable` (`live_status_app_conversation_service.py:2907-2910`). The same ZIP endpoint therefore has different consistency guarantees by deployment mode.
- **Export size cliff.** `export_max_events` defaults to 10k (`live_status_app_conversation_service.py:284,3006`); conversations that exceed it get a hard `ConversationExportTooLarge` with no partial/paginated export option — large automated runs lose the local-file escape hatch entirely.
- **Concurrent export serialization.** Second concurrent `open_conversation_export` for the same ID fails with `ConversationExportAlreadyRunning` (`live_status_app_conversation_service.py:2918-2921`) rather than queuing; clients must retry externally.

## Future Considerations

- **Promote `OTEL_*` to auto-forward or document `OH_AGENT_SERVER_ENV` for OTel.** Either add `'OTEL_'` to `AUTO_FORWARD_PREFIXES` (`sandbox_spec_service.py:151`) or extend its docstring with an OTel example parity to `LMNR_*` (`sandbox_spec_service.py:179-181`) — closes the primary interoperability papercut.
- **Add a minimal sink SPI.** An interface like `TraceSink(AsyncContextManager)` with `register_sink(name, sink)` would allow file, Langfuse, or Honeycomb adapters without forking `laminar.py`; back it with a fan-out `SpanProcessor` and a collector-mode toggle.
- **Support multi-backend fan-out via config.** Extend `_OBSERVABILITY_ENV_KEYS` to a list of named exporters or allow `OTEL_EXPORTER_OTLP_TRACES_ENDPOINTS` (comma-separated) and iterate `Laminar.initialize` + `SpanExporter` per endpoint; needed for “send to existing stack without adapter” maturity.
- **Expose tracing in `config.template.toml`.** A `[observability]` table (`enabled`, `endpoint`, `protocol`, `headers`, `service_name`) would make the feature discoverable without reading SDK docstrings; wire it through `AppConfig` with env override precedence.
- **Vendor adapters as optional extras.** Lightweight wrappers for Honeycomb (OTLP is already enough, but add header helper), Langfuse (`langfuse` SDK `observe` bridge), and LangSmith would satisfy the commercial-platform checklist without pulling heavy deps; gate behind `[tool.poetry.extras]` like `lmnr`.
- **Paginated/streaming partial exports.** For `export_max_events` overflow, return a cursor or allow `?offset&limit` ZIP streaming rather than hard failure — aligns local-file interoperability with large-run reality.
- **Add integration tests for trace emission.** Current coverage is only env-forwarding (`test_agent_server_env_override.py`); add a smoke test that starts a conversation with a fake OTLP receiver (e.g., `opentelemetry-sdk InMemorySpanExporter`) and asserts span presence, parent linkage, and metadata tags.

## Questions / Gaps

- **Is `OTEL_EXPORTER_OTLP_TRACES_HEADERS` forwarded to sandbox?** No evidence found — grepped `sandbox_spec_service.py` and found only `LMNR_`/`LLM_` prefixes. The generic OTel path therefore requires manual `OH_AGENT_SERVER_ENV`; confirm intended behavior with maintainers.
- **Collector vs direct vendor ingestion.** No evidence found of tested end-to-end paths to Honeycomb/Datadog/New Relic beyond “set OTel endpoint”. Header formats (Honeycomb `x-honeycomb-team`, Datadog `DD-API-KEY`) are undocumented; search of `*.md`/`*.toml` for those strings returned 0 hits.
- **Sampling and retention controls.** Searched `laminar.py`/`base.py` for `Sampler`/`sampling`/`retention` — No evidence found. Unclear if traces are always 100% sampled or if Laminar’s server-side sampling applies.
- **Custom file sink or local OTLP file exporter.** No evidence found of `FileExporter`, `ConsoleExporter`, or `OTLP`-to-disk path; the only local export is the conversation ZIP, which is not OTel-shaped.
- **Multi-tenancy / per-conversation backend routing.** No evidence found of routing traces by `org_id` or `provider`; all conversations in a process share the global `Laminar` singleton.

---

Generated by `Dimension 10.04: Export, Interoperability, and Observability Backends` against `openhands`.
