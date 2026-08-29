# Source Analysis: agent-framework

## Export, Interoperability, and Observability Backends

### Source Info

| Field | Value |
|-------|-------|
| Name | agent-framework |
| Path | `studies/agent-harness-study/sources/agent-framework` |
| Language / Stack | Python (core: `agent-framework-core`), .NET (`Microsoft.Agents.AI`) |
| Analyzed | 2026-08-28 |

## Summary

Agent Framework is natively instrumented on OpenTelemetry and treats export as delegated to the standard OTel SDK pipeline. The Python core (`python/packages/core/agent_framework/observability.py`) provides a thin, explicit adapter layer: it creates OTLP exporters for traces/logs/metrics over gRPC or HTTP/protobuf (`_create_otlp_exporters`), parses all standard `OTEL_EXPORTER_OTLP_*` env vars in `_get_exporters_from_env`, adds console exporters, VS Code extension exporters, and accepts arbitrary custom `SpanExporter`/`LogRecordExporter`/`MetricExporter` instances via `configure_otel_providers(exporters=...)`. Azure Monitor is integrated via `FoundryChatClient.configure_azure_monitor()` / `FoundryAgent.configure_azure_monitor()` (`python/packages/foundry/agent_framework_foundry/_chat_client.py:273`, `_agent.py:796`), while Langfuse and Comet Opik are documented as bring-your-own-SDK or OTLP-forward patterns (`python/samples/02-agents/observability/README.md:112-145`). Aspire Dashboard, Honeycomb, Jaeger, or any OTLP-compatible collector work identically because the trace format is standard GenAI semantic conventions (`OtelAttr` enum `python/packages/core/agent_framework/observability.py:176-311`). Local file export is not in core but is demonstrated in lab `FileSpanExporter` (`python/packages/lab/gaia/agent_framework_lab_gaia/gaia.py:126-162`) and Harness `FileSpanExporter` (`dotnet/samples/02-agents/Harness/Harness_Shared_Console/FileSpanExporter.cs:10`). Multiple backends can receive traces simultaneously by supplying multiple exporters; runtime reconfiguration is possible without code changes via env vars and `.env` file loading.

## Rating

**8 / 10** — Clear model with tests, explicit interfaces, and operational safeguards; approaches mature but not fully durable under exporter failure/scale (no built-in exporter health signaling, retry policy not surfaced, and commercial-platform adapters remain documentation samples rather than tested SDK integrations). Just shy of 9 because it has not demonstrated proven back-pressure/retry behavior or contract tests against Langfuse/LangSmith/Honeycomb ingest.

## Evidence Collected

Every entry MUST include a file path with line numbers. Format: `path/to/file.ts:NN`.

| Area | Evidence | File:Line |
|------|----------|-----------|
| OTLP exporter creation (gRPC+HTTP) | `_create_otlp_exporters(endpoint, protocol, headers, traces_endpoint...)` branches for `grpc` vs `http/protobuf` and imports `OTLPSpanExporter` variants; raises `ImportError` with install hint | `python/packages/core/agent_framework/observability.py:365-484` |
| OTEL env-var parsing | `_get_exporters_from_env()` reads `OTEL_EXPORTER_OTLP_ENDPOINT`, `OTEL_EXPORTER_OTLP_{TRACES,METRICS,LOGS}_ENDPOINT`, `OTEL_EXPORTER_OTLP_PROTOCOL`, `OTEL_EXPORTER_OTLP_{,TRACES,METRICS,LOGS}_HEADERS`; replicates HTTP base-URL auto-append `/v1/{traces,metrics,logs}` (`base_for_http + "/v1/..."`) | `python/packages/core/agent_framework/observability.py:487-577` |
| Export configuration entry point | `configure_otel_providers(enable_sensitive_data, enable_console_exporters, exporters, views, vs_code_extension_port, env_file_path)` builds provider pipeline; re-reads env vars after `load_dotenv` | `python/packages/core/agent_framework/observability.py:1168-1339` |
| Runtime provider wiring | `_configure()` aggregates (1) env exporters, (2) passed-in exporters, (3) console exporters, (4) VS Code port exporters; `_configure_providers()` fans out by type to `BatchSpanProcessor`/`BatchLogRecordProcessor`/`PeriodicExportingMetricReader` | `python/packages/core/agent_framework/observability.py:830-948` |
| Console/local file sinks | `ConsoleSpanExporter/LogRecordExporter/MetricExporter` added when `enable_console_exporters=True`; lab GAIA `FileSpanExporter.export()` writes JSONL per-span with `trace_id`, `span_id`, `attributes`, `status` | `python/packages/core/agent_framework/observability.py:864-869`, `python/packages/lab/gaia/agent_framework_lab_gaia/gaia.py:126-162` |
| Custom sink extensibility | `configure_otel_providers(exporters=...)` documented to accept any `LogRecordExporter|SpanExporter|MetricExporter` list; sample creates custom `OTLPSpanExporter`/`OTLPLogExporter`/`OTLPMetricExporter` with `Compression.Gzip` and passes to provider | `python/samples/02-agents/observability/configure_otel_providers_with_parameters.py:127-161`, `python/packages/core/agent_framework/observability.py:1172` |
| Azure Monitor / App Insights integration | `FoundryChatClient.configure_azure_monitor(connection_string, resource=create_resource())` delegates to `azure.monitor.opentelemetry.configure_azure_monitor`; respects sticky `disable_instrumentation()` guard | `python/packages/foundry/agent_framework_foundry/_chat_client.py:273-345`, `python/packages/foundry/agent_framework_foundry/_agent.py:796-860` |
| Third-party OTLP forwarding (Langfuse/Opik) | README patterns: set `OTEL_EXPORTER_OTLP_ENDPOINT`/`OTEL_EXPORTER_OTLP_HEADERS` for Opik; for Langfuse call `get_client(); if langfuse.auth_check(): enable_sensitive_telemetry()` — instrumentation stays default-on, no framework SDK wrapper | `python/samples/02-agents/observability/README.md:112-145` |
| Standard trace format / semconv | `OtelAttr` enum maps `gen_ai.operation.name`, `gen_ai.request.model`, `gen_ai.usage.*`, `gen_ai.agent.*`, `gen_ai.tool.*`, `gen_ai.system_instructions`, `mcp.method.name`, workflow `workflow.*`/`executor.*`/`edge_group.*` | `python/packages/core/agent_framework/observability.py:176-311` |
| MCP distributed trace propagation | Framework injects W3C `traceparent`/`tracestate` via `params._meta` using global propagator; applies to `MCPStdioTool`/`MCPStreamableHTTPTool`/`MCPWebsocketTool` only, not hosted connectors | `python/samples/02-agents/observability/README.md:156-160` |
| Service identity / resource config | `create_resource(service_name, service_version, env_file_path, **attributes)` reads `OTEL_SERVICE_NAME`/`OTEL_SERVICE_VERSION`/`OTEL_RESOURCE_ATTRIBUTES` and merges `**attributes` | `python/packages/core/agent_framework/observability.py:580-654` |
| .NET parity | `OpenTelemetryAgentBuilderExtensions` / `OpenTelemetryAgent` implement GenAI semconv `v1.37`; samples wire `OTEL_EXPORTER_OTLP_ENDPOINT` and Aspire `AddAgentHostTelemetry()` | `dotnet/src/Microsoft.Agents.AI/OpenTelemetryAgentBuilderExtensions.cs:14`, `dotnet/src/Microsoft.Agents.AI/OpenTelemetryAgent.cs:19`, `dotnet/samples/02-agents/AgentOpenTelemetry/Program.cs:26` |
| Multi-backend capability | `_configure_providers` iterates heterogeneous `exporters` list, groups by `isinstance(exp, SpanExporter/LogRecordExporter/MetricExporter)`, registers one `BatchSpanProcessor` per span exporter and one reader per metric exporter — implies fan-out to multiple collectors | `python/packages/core/agent_framework/observability.py:905-947` |
| Env example / defaults | `.env.example` exposes `OTEL_EXPORTER_OTLP_ENDPOINT="http://localhost:4317/"`, `ENABLE_SENSITIVE_DATA=true`; `ObservabilitySettings` defaults `enable_instrumentation=True`, `enable_console_exporters=False`, `vs_code_extension_port=None` | `python/.env.example:50-52`, `python/packages/core/agent_framework/observability.py:683-749` |
| Operational guards | Sticky `disable_instrumentation()` (`_user_disabled` flag) intercepts `enable_instrumentation`, `enable_sensitive_telemetry`, and direct attribute writes; `_executed_setup` prevents double provider init; `_read_bool_env`/`_read_int_env` helpers | `python/packages/core/agent_framework/observability.py:1075-1130`, `python/packages/core/agent_framework/observability.py:846` |
| Tests — observability behavior | `test_observability.py` covers span naming (`chat Test`, `execute_tool`), sensitive-data gating, streaming finalizer, `_get_exporters_from_env` HTTP base-append, gRPC verbatim, `create_resource` env merging, `create_metric_views` | `python/packages/core/tests/core/test_observability.py:208-878`, `python/packages/core/tests/conftest.py:25-89` |
| Tests — workflow & MCP spans | Workflow build/run spans (`workflow.build`, `executor.process`, `edge_group.process`), trace-context `w3c` propagation and invalid-context graceful handling | `python/packages/core/tests/workflow/test_workflow_observability.py:101-484`, `python/packages/core/tests/core/test_mcp_observability.py:102-373` |

## Answers to Dimension Questions

**1. Can traces be exported to external backends?**
Yes. Any OTLP-compatible backend (Aspire Dashboard, Jaeger, Honeycomb, Grafana Tempo, Application Insights via Azure Monitor exporter, local file) can receive traces. Mechanisms: (a) standard `OTEL_EXPORTER_OTLP_*` env vars auto-wired by `configure_otel_providers()` (`python/packages/core/agent_framework/observability.py:487-577,1168`); (b) programmatic `configure_otel_providers(exporters=[...])` for custom `SpanExporter` lists (`python/samples/02-agents/observability/configure_otel_providers_with_parameters.py:142-161`); (c) manual pipeline (`advanced_manual_setup_console_output.py:34-71`); (d) `GAIATelemetryConfig` file exporter (`python/packages/lab/gaia/agent_framework_lab_gaia/gaia.py:126-162`); (e) .NET `AddAgentHostTelemetry()` registering OTLP pipeline (`dotnet/samples/04-hosting/FoundryHostedAgents/responses/Hosted-Observability/README.md:11`).

**2. Are standard protocols supported?**
Yes, first-class. OTLP over gRPC and HTTP/protobuf are both implemented in `_create_otlp_exporters` (`python/packages/core/agent_framework/observability.py:408-483`). Trace format follows OpenTelemetry GenAI semantic conventions via `OtelAttr` (`python/packages/core/agent_framework/observability.py:176-311`). MCP distributed tracing uses W3C Trace Context propagation (`traceparent`/`tracestate` via global propagator) in `params._meta` (`python/samples/02-agents/observability/README.md:156`). Env-var names follow the OTel spec (`https://opentelemetry.io/docs/languages/sdk-configuration/otlp-exporter/` referenced at `python/packages/core/agent_framework/observability.py:516-518`).

**3. Is export configurable without code changes?**
Yes. All three signals resolve from environment at runtime without recompilation: `OTEL_EXPORTER_OTLP_{,TRACES,METRICS,LOGS}_ENDPOINT`, `OTEL_EXPORTER_OTLP_PROTOCOL`, `OTEL_EXPORTER_OTLP_HEADERS`, `OTEL_SERVICE_NAME/VERSION`, `OTEL_RESOURCE_ATTRIBUTES` (`python/packages/core/agent_framework/observability.py:524-553,643-653`). Boolean toggles `ENABLE_INSTRUMENTATION`, `ENABLE_SENSITIVE_DATA`, `ENABLE_CONSOLE_EXPORTERS`, and `VS_CODE_EXTENSION_PORT` control behavior without code (`python/packages/core/agent_framework/observability.py:695-704`). `ObservabilitySettings` and `create_resource`/`configure_otel_providers` accept `env_file_path` to load `.env` implicitly (`python/packages/core/agent_framework/observability.py:639-640,1296-1318`). Running against Aspire Dashboard is a one-line `OTEL_EXPORTER_OTLP_ENDPOINT=http://localhost:4317` change (`python/samples/02-agents/observability/README.md:444-452`).

**4. Can multiple backends receive traces simultaneously?**
Yes by design, though not prominently documented as "multi-sink." `ObservabilitySettings._configure()` merges exporters from env + passed-in list + console + VS Code extension (`python/packages/core/agent_framework/observability.py:849-877`) and `_configure_providers` registers a separate `BatchSpanProcessor`/`BatchLogRecordProcessor`/`PeriodicExportingMetricReader` per matching exporter (`python/packages/core/agent_framework/observability.py:905-947`). Example `advanced_manual_setup_console_output.py` plus an OTLP exporter list, and `GAIATelemetryConfig` which adds a `FileSpanExporter` alongside the OTLP provider (`python/packages/lab/gaia/agent_framework_lab_gaia/gaia.py:111-113`), demonstrate fan-out. No hard-coded single-exporter limit; `BatchSpanProcessor` per exporter implies independent batching/retries per backend. No explicit test asserts two remote collectors in parallel — inferred from implementation rather than proven.

## Architectural Decisions

- **Delegate to OTel SDK, do not re-implement transport.** Framework only installs instrumentation (`ChatTelemetryLayer`, `AgentTelemetryLayer`, `EmbeddingTelemetryLayer`, `get_function_span`, workflow `create_*_span`) and leaves transport to `opentelemetry-sdk` + exporter packages. Keeps core `pip install agent-framework-core` dependency minimal (`opentelemetry-api` only, `python/packages/core/pyproject.toml:29`) while `opentelemetry-sdk==1.40.0`/`azure-monitor-opentelemetry==1.8.8` are workspace dev deps (`python/pyproject.toml:46-47`). ImportError messages explicitly recommend install commands (`python/packages/core/agent_framework/observability.py:421-423,458-460`).

- **HTTP base-endpoint auto-append replication.** Because `_get_exporters_from_env()` forwards env values as constructor `endpoint=` (full URL semantics), it manually replicates the spec rule: HTTP `OTEL_EXPORTER_OTLP_ENDPOINT` auto-appends `/v1/{traces,metrics,logs}` while gRPC uses base verbatim; signal-specific endpoints are used verbatim in both protocols (`python/packages/core/agent_framework/observability.py:535-553`, tests at `python/packages/core/tests/core/test_observability.py:774-878`). Bug `OTLP HTTP base-endpoint losing /v1/{signal} auto-append #5913` motivated this fix.

- **Telemetry layers as MRO mixins.** `ChatTelemetryLayer`/`AgentTelemetryLayer`/`EmbeddingTelemetryLayer` (`python/packages/core/agent_framework/observability.py:1371,1663,1725`) are generic mixins that wrap `get_response`/`run`/`get_embeddings` with `_get_span`/`_start_streaming_span` and context-var bookkeeping for accumulated usage/response_id. Allows OpenAI/Foundry chat clients to inherit tracing without duplicating logic; also enables streaming span finalization via `ResponseStream.with_cleanup_hook`/`with_pull_context_manager`.

- **Sticky opt-out.** `disable_instrumentation()` sets `_user_disabled=True` and gates future `enable_instrumentation(force=True)`-only re-enables (`python/packages/core/agent_framework/observability.py:1106-1130`). Integrations like `FoundryChatClient.configure_azure_monitor()` check `is_user_disabled` to skip provider provisioning entirely, ensuring user intent survives auto-setup (`python/packages/foundry/agent_framework_foundry/_chat_client.py:305`).

- **No proprietary sink SDKs in core.** Langfuse/Comet/Honeycomb are supported via OTLP forwarding or SDK bring-your-own pattern rather than bundled adapters (`python/samples/02-agents/observability/README.md:112-145`). Reduces dependency surface and avoids version pin conflicts (e.g., `opentelemetry-semantic-conventions>=0.60b1` override for Mistral compatibility at `python/pyproject.toml:63`).

## Notable Patterns

- **Env-per-signal override with header merge:** per-signal endpoints/headers override base endpoint; headers are parsed (`_parse_headers` at `python/packages/core/agent_framework/observability.py:353-362`) and merged `base_headers + signal_headers` (`python/packages/core/agent_framework/observability.py:564-566`).

- **Zero-code path:** supports `opentelemetry-instrument` CLI wrapper that configures providers/exporters from env at process startup without calling `configure_otel_providers()` (`python/samples/02-agents/observability/advanced_zero_code.py`, `python/samples/02-agents/observability/README.md:152-155`).

- **Sensitive-data gating:** `OBSERVABILITY_SETTINGS.SENSITIVE_DATA_ENABLED` gates `gen_ai.input.messages`/`gen_ai.output.messages`/`gen_ai.system_instructions`/`gen_ai.tool.call.arguments` capture; re-read at call time via property, not cached (`python/packages/core/agent_framework/observability.py:806-812,1491-1506`).

- **Context-propagation-aware streaming:** non-current spans for streaming (`_start_streaming_span`, `python/packages/core/agent_framework/observability.py:2159`) finalized via `weakref.finalize` + `with_pull_context_manager(_activate_span)` to parent child HTTP/tool spans correctly without cross-context detach errors (`python/packages/core/agent_framework/observability.py:1584-1592`).

- **File sink example pattern:** both Python lab GAIA and .NET Harness share the same `FileSpanExporter` idea (append JSONL per span, `SpanExportResult.SUCCESS/FAILURE`) showing how to build local file sinks as custom `SpanExporter` subclasses.

## Tradeoffs

- **Minimal core vs external exporter installs.** Narrow `opentelemetry-api` dependency (`python/packages/core/pyproject.toml:29`) keeps install light but requires users to understand which `opentelemetry-exporter-otlp-proto-*` package to install. Error UX mitigates with explicit ImportError guidance, but adds one manual setup step per protocol.

- **Global TracerProvider singleton.** `trace.set_tracer_provider(TracerProvider(resource=resource))` and `metrics.set_meter_provider` (`python/packages/core/agent_framework/observability.py:922-947`) mutate process-global state; simplifies usage (`get_tracer()` anywhere) but conflicts with applications using multiple providers or those already configured via `azure-monitor-opentelemetry`. Guard is `_executed_setup` idempotency, not merge.

- **No built-in exporter health monitoring.** Framework exposes no health endpoint or error callback for failed exports; `BatchSpanProcessor` defaults apply (retry/schedule hidden inside SDK). Operators must monitor via collector telemetry, not framework signals.

- **Commercial adapter is documentation, not code.** Langfuse/LangSmith/Honeycomb are "bring your own SDK/OTLP creds" patterns. Avoids dependency churn (Langfuse pin `4.0.6` in `python/uv.lock:3558`) but means no contract tests prove interchange; e.g., Langfuse attribute display fixes were PR-driven (`CHANGELOG.md:906`).

- **Console/file exporters are sample/lab concern.** Core only ships console toggle (`ENABLE_CONSOLE_EXPORTERS`); durable local file export lives in `agent-framework-lab` (`python/packages/lab/gaia/...`), not tested as part of core CI, so file-sink durability is not hardened.

## Failure Modes / Edge Cases

- **Missing SDK/exporter at import time.** `_configure_providers`, `create_resource`, `_create_otlp_exporters`, `FoundryChatClient.configure_azure_monitor` raise `ModuleNotFoundError`/`ImportError` with install hints rather than failing silently (`python/packages/core/agent_framework/observability.py:633,663,421,458`, `python/packages/foundry/agent_framework_foundry/_chat_client.py:325`). Tests assert this behavior (`python/packages/core/tests/core/test_optional_dependencies.py:36-63`).

- **Empty or invalid OTLP config.** `_create_otlp_exporters` returns empty list when no endpoints provided (`python/packages/core/tests/core/test_observability.py:1990-1995`); `_get_exporters_from_env` returns empty when no env vars set (`python/packages/core/tests/core/test_observability.py:2264-2278`). Invalid `_parse_headers` entries (missing `=`) are dropped silently (`python/packages/core/agent_framework/observability.py:362`, `python/packages/core/tests/core/test_observability.py:1120-1126`).

- **Protocol mismatch.** Env value normalized `.lower()` defaults to `grpc` (`python/packages/core/agent_framework/observability.py:533`); unknown protocol yields no matching branch and thus no exporters (empty list fallback), failing silently rather than raising — operator must spot missing telemetry.

- **HTTP trailing-slash duplicate-slash avoided** via `base_endpoint.rstrip("/")` before appending (`python/packages/core/agent_framework/observability.py:546`), validated in `test_get_exporters_from_env_http_base_endpoint_trailing_slash` (`python/packages/core/tests/core/test_observability.py:803-824`).

- **Streaming spans not auto-attached.** Deliberate use of `start_span` (not `start_as_current_span`) for streaming prevents "Failed to detach context" errors; finalization via cleanup hooks ensures span closure even if stream is abandoned and garbage-collected (`python/packages/core/agent_framework/observability.py:2159-2179,1584-1592`). Risk: if `get_final_response()` is never awaited after successful streaming, span attributes may lack response usage/finish_reason but still close via GC.

- **Instrumentation disabled mid-stream.** All telemetry layers early-return `if not OBSERVABILITY_SETTINGS.ENABLED` before wrapping, so toggling off mid-process prevents new spans but existing provider/exporters remain installed (`python/packages/core/agent_framework/observability.py:1466,1685,1757,1076`). Already-queued `BatchSpanProcessor` batches still flush; no tearing-down of providers on disable.

- **Resource/config drift after first setup.** `_executed_setup=True` blocks reconfiguration (`python/packages/core/agent_framework/observability.py:846-847`). Second `configure_otel_providers()` call is no-op unless `env_file_path` triggers reset or `_executed_setup` is manually cleared; tests cover `_configure does nothing when already set up` (`python/packages/core/tests/core/test_observability.py:2296-2317`, `2967`).

## Future Considerations

- **Harden file exporter into core.** Promote GAIA `FileSpanExporter` pattern into `agent_framework.observability` as `FileSpanExporter`/`JsonlExporter` with rotation and back-pressure, covered by core tests — currently only sample/lab code.

- **Expose exporter health & error hooks.** Surface `BatchSpanProcessor`/`PeriodicExportingMetricReader` schedule, queue-size, and export failure counters; allow passing `export_error_callback` to `configure_otel_providers` so operators can alert when Honeycomb/Azure Monitor ingest returns 401/429.

- **Provide typed helpers for commercial backends.** Add small factory helpers like `create_honeycomb_exporter(api_key, dataset)` or `create_langfuse_otel_exporter()` that set correct `OTEL_EXPORTER_OTLP_HEADERS` patterns, reducing manual header-formatting errors documented in `README.md:135-145`.

- **Document multi-sink topology.** Add explicit sample `configure_otel_providers(exporters=[otlp_grpc_exporter, honeycomb_exporter, ConsoleSpanExporter()])` plus a test asserting fan-out produces 2× span counts, closing the "inferred but not proven" gap.

- **Config validation layer.** Emit warning when `OTEL_EXPORTER_OTLP_PROTOCOL` is unrecognized or when env exporters require missing optional package; fail fast rather than silently producing zero exporters.

## Questions / Gaps

- No evidence found for native **LangSmith** integration or trace format mapping — only Langfuse and Opik are referenced; searches across `python/` for `langsmith|Honeycomb` (beyond generic OTLP) return no SDK code. Standard OTLP will work if LangSmith-ingest speaks OTLP, but this is not validated.

- No evidence found for **Persistent Queue / offline buffering** for traces (e.g., disk-backed `FileExporter` + forwarder) in core — only `BatchSpanProcessor` in-memory batching; operator question: what happens when the OTLP endpoint is unreachable during burst?

- No evidence found for **Trace sampling configuration** exposure via framework helpers — operators must configure `Sampler` directly on `TracerProvider` they build themselves (`advanced_manual_setup_console_output.py` shows manual `TracerProvider` creation); framework `configure_otel_providers` does not expose a `sampler` parameter.

- No evidence found for **Export timeout/retry tunables** surfaced to callers — `BatchSpanProcessor` defaults (scheduleDelay, maxQueueSize, maxExportBatchSize) are not parametrized by the framework; searched `observability.py` and found no `BatchSpanProcessor(` call with custom schedule.

- Searched boundaries: `python/packages/core/agent_framework/observability.py` (full 2800+ line file), `python/packages/core/pyproject.toml`, `python/pyproject.toml`, `python/.env.example`, `python/samples/02-agents/observability/*`, `python/packages/foundry/.../_chat_client.py & _agent.py`, `python/packages/lab/gaia/...`, `dotnet/` samples/specs. Cross-source filesystem access was not performed per study isolation rules.

---

Generated by `10.04-export-interoperability-and-observability-backends` against `agent-framework`.
