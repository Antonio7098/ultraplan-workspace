# Source Analysis: temporal

## 10.04 Export, Interoperability, and Observability Backends

### Source Info

| Field | Value |
|-------|-------|
| Name | temporal |
| Path | `studies/agent-harness-study/sources/temporal` |
| Language / Stack | Go 1.26.4 / go.opentelemetry.io/otel v1.44.0 + otlpgrpc |
| Analyzed | 2026-08-28 |

## Summary

Temporal server is an OTEL-native distributed trace emitter, not an agent harness. Tracing is built on `go.opentelemetry.io/otel` SDK with `otlptracegrpc`/`otlpmetricgrpc` as the **only** production span/metric exporter. Export is wired at process start via fx (`temporal/fx.go:941`, `temporal/fx.go:1001`) into per-service `trace.TracerProvider`s wrapped by `BatchSpanProcessor`s and injected into gRPC (`common/telemetry/grpc.go:41`) and HTTP (`common/telemetry/http.go:22`) instrumentation. Configuration is YAML (`otel:` stanza at `common/config/config.go:53`) plus standard OTEL env vars (`OTEL_TRACES_EXPORTER` in `common/telemetry/env.go:20`). Any OTLP-compatible backend (Otel Collector, Tempo, Honeycomb, Grafana, Datadog via collector) works; no native Langfuse/LangSmith/Honeycomb SDKs, no file/local sink, and no HTTP/protobuf exporter are present. Multi-export is architecturally hinted but practically limited to one `otlp` model instance.

## Rating

**6 / 10 — Present but constrained**

OTLP/gRPC export is explicit, tested, and operationally guarded (retry, timeout, shared `grpc.ClientConn`, lifecycle hooks), satisfying the core interoperability question without an adapter **if** the observability stack speaks OTLP/gRPC. Score is capped because: (1) only `traces+otlp+grpc`/`metrics+otlp+grpc` is whitelisted (`common/telemetry/config.go:404-407`), (2) fan-out to multiple distinct endpoints collides on `SpanExporterType(model)` (`common/telemetry/config.go:227`), (3) no file/custom sink without code change ( `CustomExporters` is `yaml:"-"` testing backdoor), and (4) runtime reconfiguration requires restart — no dynamic config hot-reload.

## Evidence Collected

Every entry MUST include a file path with line numbers. Format: `path/to/file.ts:NN`.

| Area | Evidence | File:Line |
|------|----------|-----------|
| OTLP trace exporter | Imports `otlptracegrpc` and `otlpmetricgrpc` — only exporters compiled in | `common/telemetry/config.go:13-14` |
| YAML top-level stanza | `ExporterConfig telemetry.ExportConfig yaml:"otel"` merges into global config | `common/config/config.go:52-53` |
| Exporter spec struct | `otlpGrpcExporter{Connection grpcconn, Headers map[string]string, Timeout, Retry}` — retry/timeout/header knobs | `common/telemetry/config.go:100-111` |
| Supported kinds | Switch whitelists only `traces+otlp+grpc` / `metrics+otlp+grpc` ; else `unsupported exporter kind` error | `common/telemetry/config.go:402-416` |
| Span builder | `buildOtlpGrpcSpanExporter` constructs `otlptracegrpc.WithEndpoint/WithHeaders/WithTimeout/WithRetry/WithInsecure` then `otlptracegrpc.NewUnstarted` | `common/telemetry/config.go:292-325` |
| Shared connection optimisation | `sharedConnSpanExporter` defers `grpc.NewClient`+`otlptracegrpc.New(ctx, WithGRPCConn)` to `Start()` via `sync.Once` | `common/telemetry/config.go:132-140,327-340` |
| Retry defaults | `retryDefaultEnabled=true, InitialInterval 5s, MaxInterval 30s, MaxElapsedTime 1m` taken from otel v1.7 | `common/telemetry/config.go:38-41` |
| Env-var factory | `SpanExportersFromEnv` reads `OTEL_TRACES_EXPORTER`, validates `OTEL_EXPORTER_OTLP_TRACES_PROTOCOL==grpc`, returns `otlptracegrpc.NewUnstarted()` | `common/telemetry/env.go:29-52` |
| Env service name | `ResourceServiceName` maps `internal-frontend→frontend` and prefix `io.temporal` or `OTEL_SERVICE_NAME` | `common/telemetry/env.go:62-79` |
| Fx assembly & precedence | `TraceExportModule` merges config-file → env → `CustomExporters` (testing), with `maps.Copy` precedence and `startAll`/`shutdownAll` lifecycle hooks | `temporal/fx.go:941-982,1073-1103` |
| Per-service TracerProvider | Creates `otelsdktrace.NewTracerProvider(WithResource, WithSpanProcessor)` per service via `BatchSpanProcessor`; falls back to `NoopTracerProvider` when `len(sps)==0` | `temporal/fx.go:1037-1065` |
| Resource attributes | `semconv.ServiceNameKey`, `ServiceVersionKey`, `ServiceInstanceIDKey` plus OS/host/container/process | `temporal/fx.go:1018-1031` |
| gRPC instrumentation | `NewServerStatsHandler`/`NewClientStatsHandler` wrap `otelgrpc.NewServerHandler/NewClientHandler` with `isEnabled(tp)` guard returning `nil` when noop | `common/telemetry/grpc.go:41-74` |
| HTTP instrumentation | `NewHTTPClientTransport`/`NewHTTPHandler` wrap `otelhttp.NewTransport/NewHandler` plus debug header/payload annotators gated by `TEMPORAL_OTEL_DEBUG` | `common/telemetry/http.go:22-122` |
| Debug payload capture | `debugHTTPClientSpanTransport`/`debugHTTPHandler` unbounded `payloadCapture.Bytes.Buffer` — complete diagnostics vs memory risk | `common/telemetry/http.go:172-336` |
| Wide events (OTEL Logs) | `Payload` interface + `Emit(logger log.Logger,p Payload)` via `log.Record.SetEventName/AddAttributes` — separate log signal, not trace span | `common/wideevents/events.go:15-53` |
| Metrics OTLP path | Metrics module also supports `otlp+grpc` for metrics, and `openTelemetryProviderImpl` exposes `prometheus`/`statsd` readers | `common/metrics/opentelemetry_provider.go:35-122` |
| Commercial backend example | Doc shows Honeycomb via OTLP headers `x-honeycomb-team: <key>`; Tempo quickstart via `otel-collector` | `docs/development/tracing.md:60-76,15-23` |
| Testing collector & in-memory exporter | `testtelemetry.MemoryCollector` gRPC `TraceServiceServer` for integration tests; `tracetest.NewInMemoryExporter` in unit tests | `common/testing/testtelemetry/collector.go:19-72`, `common/telemetry/grpc_test.go:28-29` |
| Env+binary wiring for local dev | `Makefile` OTEL=true sets `OTEL_TRACES_EXPORTER=otlp, OTEL_EXPORTER_OTLP_TRACES_INSECURE=true, OTEL_BSP_SCHEDULE_DELAY` | `Makefile:OTEL ?= false` block |
| No custom sinks | Grep for `langfuse|langsmith|honeycomb|Custom sink|file.*export` finds zero production registrations; `CustomExporters` tagged `yaml:"-"` | `common/telemetry/config.go:155-156` |
| No file/local exporter | No `stdout`, `file`, `otlptracehttp`, `jaeger`, `zipkin` branches; only grpc path exists | `common/telemetry/config.go:404-407` |

## Answers to Dimension Questions

| Question | Answer | Evidence |
|----------|--------|----------|
| **1. Can traces be exported to external backends?** | **Yes, but only via OTLP/gRPC.** Any OTLP-compatible collector or SaaS that accepts `otlp+grpc` works (Tempo, Honeycomb, Jaeger via collector, Datadog via collector, New Relic). No direct file, stdout, or vendor SDK. When no exporter configured, `NoopTracerProvider` silences collection (`temporal/fx.go:1039`, `common/telemetry/config.go:45`). | `common/telemetry/config.go:13-14,292-314`, `docs/development/tracing.md:36,60-76`, `temporal/fx.go:948-956` |
| **2. Are standard protocols supported?** | **Partially.** `OTLP/gRPC` is standard and fully supported. `OTLP/HTTP` (common for serverless), `W3C TraceContext` propagation (yes via `otelgrpc`/`otelhttp`), but Jaeger/Zipkin native, `otlp/http/protobuf`, or `otlp/http/json` are **not** implemented. The switch explicitly rejects them (`common/telemetry/config.go:408-413`). `OTEL_EXPORTER_OTLP_TRACES_PROTOCOL` only accepts `grpc` (`common/telemetry/env.go:44`). | `common/telemetry/config.go:404-407`, `common/telemetry/env.go:22-48` |
| **3. Is export configurable without code changes?** | **Yes — YAML + env, but restart-required.** `otel.exporters[].spec.{connection.endpoint, headers, timeout, retry, insecure}` covers endpoint/auth tuning (`common/telemetry/config_test.go:12-32`); `OTEL_TRACES_EXPORTER`, `OTEL_EXPORTER_OTLP_TRACES_PROTOCOL`, `OTEL_EXPORTER_OTLP_ENDPOINT`, `OTEL_SERVICE_NAME`, `TEMPORAL_OTEL_DEBUG` are read at startup (`common/telemetry/env.go:19-24`, `common/telemetry/config.go:26`). Env overrides YAML via `maps.Copy` precedence (`temporal/fx.go:969-970`). Note docs caveat: env wins over YAML and OTEL SDK env vars like `OTEL_EXPORTER_OTLP_TRACES_INSECURE` are picked up automatically by the exporter (`docs/development/tracing.md:110-111`). No hot reload or dynamic config flag — `ExporterConfig` is not a `dynamicconfig` property. | `common/config/config.go:53`, `temporal/fx.go:949-970`, `docs/development/tracing.md:84-111` |
| **4. Can multiple backends receive traces simultaneously?** | **Architecturally claimed but practically broken for the sole supported model.** `docs/development/tracing.md:78-80` says add more `kind+spec` entries; `temporal/fx.go:971` collects `expmaps.Values`. However `exportConfig.SpanExporters()` keys by `SpanExporterType(Model)` (`common/telemetry/config.go:227`), so two `model: otlp` entries overwrite each other — no fan-out to two OTLP endpoints. Only way to send to two places is via (a) an external Otel Collector fan-out, or (b) injecting via `CustomExporters` map (`common/telemetry/config.go:156`, `temporal/fx.go:966`) which requires code/test injection. Multiple distinct models (e.g., hypothetical `otlp` vs `custom`) would work, but none else exist. | `common/telemetry/config.go:215-227`, `docs/development/tracing.md:78-80`, `common/telemetry/config_test.go:34-55` |

## Architectural Decisions

| Decision | Rationale / Consequence | File:Line |
|----------|-------------------------|-----------|
| **OTLP/gRPC-only** | Minimizes dependency surface (two otel exporters vs matrix of HTTP/jaeger/zipkin) and aligns with collector-centric deployment. Consequence: HTTP-only environments need a sidecar/collector. | `common/telemetry/config.go:13-14,402-416`, `go.mod:62-63` |
| **Fx lifecycle-bound exporters** | `TraceExportModule` (`temporal/fx.go:940`) owns exporter lifetime via `OnStart: startAll` / `OnStop: shutdownAll` with 1s shutdown timeout ignoring `DeadlineExceeded` (drops traces on shutdown gracefully). Global `otel.ErrorHandler` suppressed until `tracingReady` (`temporal/fx.go:943-945`). | `temporal/fx.go:941-982,1088-1102` |
| **Shared `grpc.ClientConn`** | `sharedConn{Span,Metric}Exporter` (`common/telemetry/config.go:132-149`) share a named `grpcconn` across trace+metric exporters, avoiding double dial and controlling `ConnectParams/BufferSize/Authority`. | `common/telemetry/config.go:132-140,327-355,357-367` |
| **Noop guard pattern** | `isEnabled(tp)` checks `otelnoop.TracerProvider` (`common/telemetry/config.go:49-55`); handlers return `nil` to skip `otelgrpc` overhead entirely. Same guard for HTTP. | `common/telemetry/grpc.go:46-47,66-67`, `common/telemetry/http.go:26,51,74` |
| **W3C TraceContext propagation by default** | `propagation.TraceContext{}` wired through both gRPC stats handler and `otelhttp` transport (`temporal/fx.go:1067`). | `common/telemetry/grpc.go:52-53`, `common/telemetry/http.go:77-78,110-111` |
| **OTEL logs for domain events** | Wide events use `go.opentelemetry.io/otel/log` (`common/wideevents/events.go:1-53`) decoupled from tracing, allowing log export via separate LoggerProvider — but not wired to trace exporter pipeline. | `common/wideevents/events.go:12-53` |
| **Config `yaml:"-"` backdoor for testing** | `CustomExporters map[SpanExporterType]SpanExporter` enables in-memory exporters in `tests/testcore/test_cluster.go:301` without file/env. | `common/telemetry/config.go:155-156`, `temporal/fx.go:966` |

## Notable Patterns

*   **Three-layer precedence merge** — config file → env vars → programmatic `CustomExporters` (`temporal/fx.go:949-970`) gives clean layering for prod vs dev vs test without duplicating parsing logic.
*   **BatchSpanProcessor per service** — Each service flexes its own `TracerProvider` with batched processors (`temporal/fx.go:1005-1010`, `1037-1046`) so frontend/history/matching/worker can run co-located with isolation.
*   **Debug mode payload mirroring** — `TEMPORAL_OTEL_DEBUG` adds gRPC `rpc.request.payload/response.payload` (`common/telemetry/grpc.go:128-156`) and HTTP header/payload attributes (`common/telemetry/http.go:82-87`), intentionally verbose and unbounded.
*   **Attribute budget via `WorkflowIDKey/RunIDKey`** — Minimal high-cardinality span attributes extracted via `logtags.WorkflowTags` (`common/telemetry/grpc.go:160-182`) to keep cardinality manageable.
*   **Collector-in-memory test double** — `testtelemetry.MemoryCollector` implements `ctrace.TraceServiceServer` (`common/testing/testtelemetry/collector.go:18-72`) enabling hermetic OTLP round-trip tests.

## Tradeoffs

*   **Interoperability breadth vs maintenance cost:** Betting exclusively on `otlp+grpc` makes the server agnostic to vendor (Honeycomb, Tempo, Datadog via collector all work) but forces users without gRPC egress to run a collector shim. No evidence of `otlptracehttp` despite its popularity in serverless.
*   **Single-model dedup vs true fan-out:** Map keyed by model simplifies lookup but silently drops duplicate `otlp` entries — violates the doc claim of “additional kind/spec declarations” without a warning; user expecting dual export (e.g., Tempo + Honeycomb) must discover the collision.
*   **Debuggability vs safety:** `TEMPORAL_OTEL_DEBUG` captures full proto payloads/headers verbatim (`common/telemetry/grpc.go:98-136`, `common/telemetry/http.go:176-182`) which is invaluable for diagnosis but risks PII leakage and unbounded memory (`payloadCapture` buffers all bytes `common/telemetry/http.go:283-284`).
*   **Startup validation vs hot reconfig:** Validation happens at `ServerFx` start (`temporal/fx.go:954`) with clear errors (`unsupported exporter kind`, `connection not found`); downside is no zero-downtime rotation of endpoints/headers.
*   **Noop fast-path vs observability by default:** Default is silent no-op (no collector needed to run), avoiding mandatory sidecar, but means misconfigured `otel:` stanza is only noticed when traces never appear; `otel.ErrorHandler` is muted until `tracingReady`.

## Failure Modes / Edge Cases

*   **No collector / unreachable endpoint:** `otlpmetricgrpc.WithRetry`/`otlptracegrpc.WithRetry` with 5s/30s/1m retry (`common/telemetry/config.go:38-41,265-269`) retries; but `Start()` failure returns early and blocks `fx.Start` (`temporal/fx.go:976-979`). Runtime network partitions after start are retried but may fill `BatchSpanProcessor` queue and drop spans.
*   **Shutdown drops spans:** `shutdownAll` uses 1s context (`temporal/fx.go:1090-1091`) and explicitly ignores `DeadlineExceeded` (`temporal/fx.go:1095-1098`), documented as “okay to drop”. High-throughput shutdown may lose tail spans.
*   **Model collision silent loss:** Two `signal: traces, model: otlp` entries — last wins, no error, no metric — user believes dual export works. Detect via `common/telemetry/config.go:227`.
*   **Insecure without TLS:** `insecure: true` required for `localhost:4317` dev (`common/telemetry/config_test.go:30`); forgetting it against plaintext collector yields dialect failure with generic gRPC error (workaround issue #2940 `common/telemetry/config.go:273-274`).
*   **Shared conn lifecycle race:** `sharedConnSpanExporter.Start` uses `sync.Once` (`common/telemetry/config.go:329-330`) — second `Start` is no-op even if first dial failed and `err != nil` was set (error stored but `SpanExporter` stays nil; subsequent exports nil-panic risk).
*   **Env protocol mismatch:** Setting `OTEL_EXPORTER_OTLP_TRACES_PROTOCOL=http/protobuf` returns `unsupported OTEL exporter protocol` error and aborts server start (`common/telemetry/env.go:43-47`).
*   **Debug memory blow-up:** Large workflow payloads buffered fully in `payloadCapture.Buffer` (`common/telemetry/http.go:283-284`) without size cap; potential OOM under debug blast radius.
*   **Version skew:** `go.mod:62` pins `otlptracegrpc v1.43.0` with older grpc semantics (`WithInsecure()` workaround `common/telemetry/config.go:309`), needing periodic bump to avoid collector incompatibilities.

## Future Considerations

*   Add `otlptracehttp` variant alongside grpc (small `common/telemetry/config.go:404` addition) to unlock serverless/HTTP-only networks; reuse same `Headers/Timeout/Retry` knobs.
*   Fix fan-out by keying exporters by `metadata.name` or composite `signal+model+name` instead of bare `model`, or at least warn on overwrite; validate at `UnmarshalYAML`.
*   Expose `CustomExporters` as a proper plugin SPI (e.g., `TelemetryProvider` interface) so Langfuse/LangSmith adapters can be registered without fork — today only testing can inject.
*   Add file/stdout exporter (`stdouttrace`) for local development/air-gapped debugging, as asked by dimension — cheap to add and complements Tempo quickstart.
*   Support hot reload of headers/endpoint via dynamic config or SIGHUP, avoiding restart for key rotation (Honeycomb `x-honeycomb-team`).
*   Cap debug payload capture (e.g., 64KB) and redact sensitive headers when `TEMPORAL_OTEL_DEBUG` is on.
*   Wire `wideevents` LoggerProvider to same `ExporterConfig` so logs and traces share collector lifecycle (today decoupled).

## Questions / Gaps

*   No evidence of log export pipeline beyond `wideevents` — is there an `otlp/log` exporter planned or intentionally omitted? Searched `common/telemetry` — no `otlplog` import found.
*   No documentation on sampling configuration — `BatchSpanProcessorOption` slice is supplied empty by default (`temporal/fx.go:1002`); where is tail/parent sampling to be injected via `fx.Decorate`?
*   Whether `OTEL_TRACES_SAMPLER` env var is respected by the SDK (should be, via global env) is not documented in `docs/development/tracing.md`.
*   No `docker-compose` Tempo/Collector definition checked — `make start-dependencies` reference exists but manifest not in this source slice.

---

Generated by `Dimension 10.04: Export, Interoperability, and Observability Backends` against `temporal`.
