# Source Analysis: langfuse

## 10.04 Export, Interoperability, and Observability Backends

### Source Info

| Field | Value |
|-------|-------|
| Name | langfuse |
| Path | `studies/agent-harness-study/sources/langfuse` |
| Language / Stack | TypeScript (Next.js web, Node worker, pnpm/turbo, Prisma/Postgres + ClickHouse, BullMQ/Redis, S3/GCS/Azure/OCI) |
| Analyzed | 2026-08-28 |

## Summary

Langfuse is an observability *platform* rather than an agent *harness*, so the dimension inverts: Langfuse **is** the backend that traces flow *into* via OTLP, and it also **emits** its own operational traces via OTLP while **exporting stored telemetry out** via batch blob-storage integrations and pull APIs. Inbound is spec-compliant OTLP/HTTP (JSON + binary protobuf + gzip) at `POST /api/public/otel/v1/traces` (`web/src/pages/api/public/otel/v1/traces/index.ts:78-114`, `fern/apis/server/definition/opentelemetry.yml:7-34`) translated by `OtelIngestionProcessor` into a proprietary Zod-validated Langfuse envelope (`packages/shared/src/server/ingestion/types.ts:259-279,395-755`), S3-staged, then enqueued on sharded `otel-ingestion-queue`/`secondary-otel-ingestion-queue` (`packages/shared/src/server/queues.ts:352-353`, `packages/shared/src/server/redis/otelIngestionQueue.ts:13-88`). Outbound infra telemetry is a standard `OTLPTraceExporter` from both `web` and `worker` (`web/src/observability.config.ts:114-116`, `worker/src/instrumentation.ts:28-32`) with sampling via `OTEL_TRACE_SAMPLING_RATIO`. Stored-trace export is **batch, not live**: configurable blob-storage integrations to S3/S3-compatible/Azure/GCS/OCI (`fern/apis/server/definition/blob-storage-integrations.yml:39-91`, `packages/shared/src/server/services/StorageService.ts:257-327`) driven by `CoreDataS3ExportQueue`/`BlobStorageIntegrationQueue`, plus public pull APIs (`GET /api/public/v2/traces`, `/observations`) and no native per-trace OTLP forwarder to Honeycomb/Jaeger/LangSmith or local-file sink.

## Rating

**7 / 10** — Clear model with tests, explicit interfaces, and operational safeguards; approaches mature on the inbound path but not durable fan-out.

Rationale per rubric: inbound OTLP is validated, gzip-aware, S3-staged with exponential BullMQ retries (`attempts: 6`, `backoff: exponential 5000ms` at `packages/shared/src/server/redis/otelIngestionQueue.ts:67-71`) and metrics/DLQ handling; outbound infra export uses standard `OTEL_EXPORTER_OTLP_ENDPOINT` and sampler. Gaps prevent 9: stored-trace export is 20-minute batched blob dump (no live OTLP fan-out, no multi-backend simultaneous forwarding, no local-file sink, no custom exporter plugin interface), trace format is proprietary (OTel-to-Langfuse translation at `packages/shared/src/server/otel/OtelIngestionProcessor.ts:226-593` with no reverse mapper), and runtime reconfiguration of the OTLP exporter requires restart.

## Evidence Collected

Every entry MUST include a file path with line numbers. Format: `path/to/file.ts:NN`.

| Area | Evidence | File:Line |
|------|----------|-----------|
| OTLP inbound — spec & endpoint | Fern service defines `POST /otel/v1/traces` with `application/json` + `application/x-protobuf` + gzip, `ExportTraceServiceRequest/resourceSpans` shape | `fern/apis/server/definition/opentelemetry.yml:7-34` |
| OTLP inbound — handler | Content-type sniff, gzip `gunzip`, protobuf decode via `$root.opentelemetry.proto.collector.trace.v1.ExportTraceServiceRequest.decode`, JSON parse, 16 MB warn, `x-langfuse-*` header extraction, version gate `>4` reject | `web/src/pages/api/public/otel/v1/traces/index.ts:64-160` |
| OTLP inbound — S3 staging + queue publish | `fileKey = ${prefix}otel/${projectId}/${yyyy/mm/dd/hh/mm}/${uuid}.json`, `uploadJson` then `OtelIngestionQueue.getInstance({shardingKey: projectId-fileKey}).add(OtelIngestionJob)` | `packages/shared/src/server/otel/OtelIngestionProcessor.ts:182-219` |
| OTLP inbound — queue contract | `OtelIngestionEvent` Zod schema: `fileKey`, `publicKey`, `authCheck.scope.projectId/orgId`, `propagatedHeaders`, `sdkName/Version/ingestionVersion` | `packages/shared/src/server/queues.ts:43-60` |
| OTLP inbound — queue sharding & durability | `OtelIngestionQueue` + `SecondaryOtelIngestionQueue` sharded via `LANGFUSE_OTEL_INGESTION_QUEUE_SHARD_COUNT`, `removeOnComplete:true, removeOnFail:100000, attempts:6, backoff exponential 5000` | `packages/shared/src/server/redis/otelIngestionQueue.ts:13-88` |
| OTLP → Langfuse translation | `processToEvent` → `processToIngestionEvents` extracts resource/scope/span attributes, `convertValueToPlainJavascript`, usage/cost, `observationTypeMapper.mapToObservationType`, `trace-create` + `*-create` event emission; handles `service.name`, `langfuse.observation.type`, `gen_ai.*`, `ai.*`, `gcp.vertex.*`, `lk.*`, `mlflow.*`, `traceloop.*` | `packages/shared/src/server/otel/OtelIngestionProcessor.ts:226-593,1218-1505` |
| OTLP attribute convention | Example `langfuse.observation.type: generation` on OTel span attributes; `OtelSpan.attributes: list<OtelAttribute>` | `fern/apis/server/definition/opentelemetry.yml:58-63` |
| Canonical trace format — Langfuse envelope | Zod discriminated union `ingestionEvent = z.discriminatedUnion("type", [traceEvent, spanCreate, generationCreate, ... legacyObservationCreate])`; `TraceBody {id, name, input, output, sessionId, userId, environment, metadata, tags}`; `CreateGenerationBody` with `usage/usegeDetails/costDetails/promptName` | `packages/shared/src/server/ingestion/types.ts:259-279,395-755,405-421` |
| Legacy envelope retained | `OBSERVATION_CREATE`/`OBSERVATION_UPDATE` + 5 legacy body schemas for rolling compat | `packages/shared/src/server/ingestion/types.ts:276-383,689-718` |
| Outbound infra OTLP — web | `NodeSDK` with `OTLPTraceExporter { url: ${OTEL_EXPORTER_OTLP_ENDPOINT}/v1/traces }`, instrumentations: IORedis (with `ioredisRequestHook` redaction), Http, Prisma, AWS, Winston, BullMQ; `TraceIdRatioBasedSampler(env.OTEL_TRACE_SAMPLING_RATIO)` | `web/src/observability.config.ts:109-165` |
| Outbound infra OTLP — worker | `NodeSDK` `resource: service.name=OTEL_SERVICE_NAME`, `OTLPTraceExporter url ${OTEL_EXPORTER_OTLP_ENDPOINT}/v1/traces` | `worker/src/instrumentation.ts:28-32` |
| Instrumentation helpers | `instrumentAsync/instrumentSync` with W3C propagation extract, `traceException` -> `recordException` + `error.*` attributes, `addUserToSpan` with baggage, `recordIncrement/recordDistribution` dual Datadog + CloudWatch | `packages/shared/src/server/instrumentation/index.ts:54-189,191-261,328-367` |
| Env configurability — OTEL | `OTEL_EXPORTER_OTLP_ENDPOINT default http://localhost:4318`, `OTEL_SERVICE_NAME default worker`, sampling ratio via `web/src/env.mjs` + `worker/src/env.ts` | `worker/src/env.ts:172-173`, `web/src/observability.config.ts:115,162`, `.env.prod.example:40-41` |
| Env configurability — ingestion tuning | `LANGFUSE_OTEL_INGESTION_QUEUE_SHARD_COUNT`, `LANGFUSE_OTEL_INGESTION_SECONDARY_QUEUE_SHARD_COUNT`, `LANGFUSE_OTEL_MAX_SPAN_BYTES=9.5MB`, `LANGFUSE_S3_EVENT_UPLOAD_BUCKET/PREFIX` | `packages/shared/src/env.ts:151-159,232-233` |
| Batch blob export — integration model | `BlobStorageIntegrationType {S3, S3_COMPATIBLE, AZURE_BLOB_STORAGE}`, `BlobStorageExportSource {LEGACY_TRACES_OBSERVATIONS, OBSERVATIONS_V2, LEGACY_TRACES_AND_ENRICHED_OBSERVATIONS}`, `BlobStorageExportFieldGroup {core,basic,time,io,metadata,model,usage,prompt,metrics,tools,trace_context}`, `exportFrequency {every_20_minutes,hourly,daily,weekly}`, `CreateBlobStorageIntegrationRequest` with bucket/endpoint/region/accessKey/prefix/compressed | `fern/apis/server/definition/blob-storage-integrations.yml:39-146` |
| Batch blob export — queues | `CoreDataS3ExportQueue`, `BlobStorageIntegrationQueue` / `BlobStorageIntegrationProcessingQueue` | `packages/shared/src/server/queues.ts:366,364-365` |
| Batch blob export — storage backends | `StorageServiceFactory.getInstance` selects `S3StorageService` / `AzureBlobStorageService` / `GoogleCloudStorageService` / `OCIObjectStorageService` based on `LANGFUSE_USE_AZURE_BLOB/GOOGLE_CLOUD_STORAGE/OCI_NATIVE_OBJECT_STORAGE` or per-integration params; handles SSE, concurrent writes, signing diagnostics | `packages/shared/src/server/services/StorageService.ts:257-327,330-650,957-1649` |
| Public pull APIs | `POST /api/public/traces` (legacy `trace-create` via `processEventBatch`), `GET /api/public/traces` with field groups + events table path, `GET /api/public/v2/traces|observations|scores`, `GET /api/public/dataset-items` | `web/src/pages/api/public/traces/index.ts:37-69,72-187`, `fern/apis/server/definition/opentelemetry.yml:7-34`, `packages/shared/src/server/queues.ts:366` |
| Metrics/secondary sinks | Datadog `dd-trace` + `dogstatsd.gauge/increment/histogram/distribution`, CloudWatch `PutMetricData` batched 30s, metric namespaces `langfuse.ingestion.otel.*`, `langfuse.queue.*` | `packages/shared/src/server/instrumentation/index.ts:265-380`, `web/src/observability.config.ts:104-107`, `packages/shared/src/server/otel/OtelIngestionProcessor.ts:331-343` |
| Generated OTLP proto | Hand-managed compiled file `web/src/pages/api/public/otel/otlp-proto/generated/root.ts` (no regen script; README warns) | `web/src/pages/api/public/otel/otlp-proto/generated/root.ts:1-8` (via `web/src/pages/api/public/otel/otlp-proto/generated/README.md:7`) |
| Oversized span guard | `LANGFUSE_OTEL_MAX_SPAN_BYTES` threshold, per-field >1MB `largeFields` logging, `recordIncrement langfuse.ingestion.otel.oversized_span` | `packages/shared/src/server/otel/OtelIngestionProcessor.ts:1235-1283` |
| OTEL project tracking | `markProjectAsOtelUser(projectId)` on ingress; Redis `otelProjectTracking` for `LANGFUSE_SKIP_FINAL_FOR_OTEL_PROJECTS` optimization | `web/src/pages/api/public/otel/v1/traces/index.ts:48`, `packages/shared/src/server/redis/otelProjectTracking.ts:1-51` |

## Answers to Dimension Questions

**1. Can traces be exported to external backends?**

Yes, but in two distinct directions with different ergonomics. **Inbound** (traces flowing *to* Langfuse) is mature: any OTel SDK can `POST /api/public/otel/v1/traces` with JSON or protobuf+gzip (`web/src/pages/api/public/otel/v1/traces/index.ts:78-114`, `fern/apis/server/definition/opentelemetry.yml:7-34`). **Outbound from stored data** is batch-oriented: blob-storage integrations push `traces/observations/scores` to customer-owned S3/GCS/Azure/OCI on a schedule (`every_20_minutes|hourly|daily|weekly` at `fern/apis/server/definition/blob-storage-integrations.yml:57-62`) as `JSON|CSV|JSONL` with `full_history|from_today|from_custom_date` modes, or consumers pull via public REST (`GET /api/public/v2/traces`, `GET /api/public/observations` at `web/src/pages/api/public/traces/index.ts:72-187`). There is **no live per-trace OTLP forwarder** that re-emits stored Langfuse spans to Honeycomb/Jaeger/Tempo/LangSmith — that path requires writing an adapter that polls the public API or tails the blob export. Langfuse's own service traces *do* forward live to an OTLP collector (`OTLPTraceExporter` at `web/src/observability.config.ts:114-116`, `worker/src/instrumentation.ts:32`), but that covers infrastructure spans, not customer LLM traces.

**2. Are standard protocols supported?**

**Inbound**: fully OTLP/HTTP spec-compliant (`https://opentelemetry.io/docs/specs/otlp/#otlphttp` cited at `fern/apis/server/definition/opentelemetry.yml:18`), accepting `ExportTraceServiceRequest` JSON and binary protobuf, gzip, and `resourceSpans/scopeSpans/spans` with `langfuse.observation.type` attribute convention (`fern/apis/server/definition/opentelemetry.yml:58-63`). **Outbound infra**: standard `OTLPTraceExporter` with `url: ${OTEL_EXPORTER_OTLP_ENDPOINT}/v1/traces` (`web/src/observability.config.ts:115`), W3C trace-context propagation (`packages/shared/src/server/instrumentation/index.ts:61-64`), and `OTEL_SERVICE_NAME` resource attributes. **Stored-trace wire format**: proprietary — Zod `TraceBody`/`CreateSpanBody`/`CreateGenerationBody` (`packages/shared/src/server/ingestion/types.ts:405-453`) persisted to ClickHouse/Postgres, not OTel `SpanKind`/`Status`/`Resource`. Translating back to OTel for Honeycomb requires custom mapping of `observation.type` (`SPAN|GENERATION|AGENT|TOOL` at `fern/apis/server/definition/ingestion.yml:147-158`).

**3. Is export configurable without code changes?**

Partially. **Collector target** is env-driven: `OTEL_EXPORTER_OTLP_ENDPOINT` (`worker/src/env.ts:172`, `.env.prod.example:40`), `OTEL_SERVICE_NAME`, `OTEL_TRACE_SAMPLING_RATIO` (`web/src/observability.config.ts:162`), plus queue-shard counts and `LANGFUSE_OTEL_MAX_SPAN_BYTES` (`packages/shared/src/env.ts:151-159`). Changing the infra OTLP endpoint or sampling ratio requires only env change + restart (no code). **Batch blob export** is fully runtime-configurable via API: `PUT /api/public/integrations/blob-storage` with `bucketName/endpoint/region/prefix/exportFrequency/fileType/exportSource/exportFieldGroups` (`fern/apis/server/definition/blob-storage-integrations.yml:92-146`), validated with a test upload. **Per-trace live export routing** (e.g., “forward trace X to Honeycomb and trace Y to Datadog based on project/tag”) is not configurable without code — there is no routing table or exporter plugin registry; adding a new sink type requires code to extend `StorageServiceFactory` (`packages/shared/src/server/services/StorageService.ts:257-327`) or to write a poller against the public API.

**4. Can multiple backends receive traces simultaneously?**

Not as simultaneous live fan-out. The `OTLPTraceExporter` is a single `url` (`web/src/observability.config.ts:115`, `worker/src/instrumentation.ts:32`) — one OTLP endpoint per process. Blob export has one integration per project (API model is `CreateBlobStorageIntegrationRequest` with `projectId` scoping at `fern/apis/server/definition/blob-storage-integrations.yml:92-96`); scheduling via `CoreDataS3ExportQueue` means one destination bucket per config, though an operator could pull simultaneously via public API while blob export runs (de facto multi-consumer via pull). There is no `exporters: [otel, honeycomb, langsmith, file]` list or `configure_otel_providers(exporters=[...])`-style multi-sink wiring within Langfuse itself; simultaneous delivery requires an external fan-out (e.g., OpenTelemetry Collector `exporters: [otlp/langfuse, otlp/honeycomb]` upstream of Langfuse, or a downstream job that replicates blob files to multiple buckets).

## Architectural Decisions

- **Dual-ingest, single-store (OTel front door + Langfuse envelope)** (`web/src/pages/api/public/otel/v1/traces/index.ts:78-114` + `packages/shared/src/server/ingestion/types.ts:259-279,689-718`): keeps SDK ergonomics (Langfuse JS/Python SDKs emit envelope) while accepting any OTel-instrumented app via the deprecated-`POST /ingestion` → OTel consolidation (`fern/apis/server/definition/ingestion.yml:12`). Trade is the `OtelIngestionProcessor` translation tax (`packages/shared/src/server/otel/OtelIngestionProcessor.ts:226-593`) and loss of strict Zod validation on the OTel path (Fern types use `optional<unknown>` for `traceId/spanId/value` at `fern/apis/server/definition/opentelemetry.yml:110-163`).

- **S3-as-WAL + sharded BullMQ** (`packages/shared/src/server/otel/OtelIngestionProcessor.ts:182-189`, `packages/shared/src/server/redis/otelIngestionQueue.ts:13-88`, `packages/shared/src/server/queues.ts:352-353`): full `resourceSpans` batch is uploaded to `LANGFUSE_S3_EVENT_UPLOAD_BUCKET` before ack, then only the `fileKey` is queued. Provides durability across web/worker crash and horizontal scale via shard hash on `projectId-fileKey`. Secondary queues isolate high-throughput noisy neighbors.

- **Standard `OTLPTraceExporter` for internal observability, not for customer trace re-export** (`web/src/observability.config.ts:14-16`, `worker/src/instrumentation.ts:28-32`): `NodeSDK` with `TraceIdRatioBasedSampler` and Datadog dual-write (`packages/shared/src/server/instrumentation/index.ts:265-367`) gives vendor-neutral infra telemetry. Customer traces stay in Langfuse's own columnar store (ClickHouse) and leave only via batch blob or pull API — a deliberate platform ownership of the query path.

- **Factory-selected blob backends** (`packages/shared/src/server/services/StorageService.ts:257-327`): `useAzureBlob/useGoogleCloudStorage/useOCIObjectStorage` flags plus per-integration `endpoint/region/forcePathStyle` allow S3-compatible destinations (MinIO, GCS interop) without code forks, with signing-diagnostics middleware for `SignatureDoesNotMatch` on non-AWS S3.

- **Batch export with field groups and source versioning** (`fern/apis/server/definition/blob-storage-integrations.yml:64-91`, `packages/shared/src/env.ts:232-233`): `exportSource` (`LEGACY_TRACES_OBSERVATIONS` → `OBSERVATIONS_V2` with deprecation gates `2026-05-20/2026-06-22`) and `exportFieldGroups` let operators trim columns (e.g., drop `trace_context` for legacy) and cut history vs. incremental — but still 20-minute minimum lag (`BlobStorageSyncStatus` doc at `fern/apis/server/definition/blob-storage-integrations.yml:189-191`).

## Notable Patterns

- **Gateway validates, stages, then translates**: `withMiddlewares` + `createAuthedProjectAPIRoute` auth, `z.any()`-tolerant OTel parse, immediate `markProjectAsOtelUser` tracking (`web/src/pages/api/public/otel/v1/traces/index.ts:33-48`), then durable S3 + async BullMQ. Keeps the hot HTTP path < few ms even under large `resourceSpans`.

- **OTel attribute → Langfuse domain mapper**: `OtelIngestionProcessor` recognizes `gen_ai.*`, `langfuse.observation.*`, `ai.prompt.*`, `gcp.vertex.*`, `lk.*`, `mlflow.*`, `traceloop.*`, `genkit:*`, `flue.*` (`packages/shared/src/server/otel/OtelIngestionProcessor.ts:1436-1505`) and normalizes `input/output/metadata` via `extractInputAndOutput` with deduplication of 30+ potential keys.

- **Sampling as cost control**: `OTEL_TRACE_SAMPLING_RATIO` sampled at SDK level (`web/src/observability.config.ts:162`) — infra traces are sampled before export, not after, so collector never sees dropped spans. Customer LLM traces are *not* sampled; ingestion failure tracking is the backpressure signal instead (`markProjectIngestFailure` at `web/src/pages/api/public/otel/v1/traces/index.ts:191`).

- **Shallow vs. full trace deduplication**: `OtelIngestionProcessor.filterRedundantShallowTraces` and `hasTraceUpdates` logic (`packages/shared/src/server/otel/OtelIngestionProcessor.ts:600-651,1083-1116`) deduplicates `trace-create` events when a root span already emitted full trace metadata.

- **Scope-filtering tracer provider for Next.js noise**: `ScopeFilteringTracerProvider` + `PassthroughTracer` mutes `next.js` scope spans while preserving parent `http.server` (`web/src/observability.config.ts:42-102`), preventing double-counted server spans from drowning customer trace signal.

## Tradeoffs

- **Otel fidelity vs. Langfuse ergonomics**: accepting `optional<unknown>` OTel IDs/values loses validation strictness (`fern/apis/server/definition/opentelemetry.yml:110-163` vs. `idSchema.min(1).max(800)` at `packages/shared/src/server/ingestion/types.ts:10-16`) to stay wire-compatible, while canonical `TraceBody` stays rich and queryable. Reverse translation (Langfuse → OTel) is lossy and unimplemented.

- **Durability vs. latency**: S3 staging gives crash safety and replay via `s3-ingestion-event-replay.ts` (`worker/src/scripts/replayIngestionEvents/s3-ingestion-event-replay.ts:251-403`) but adds S3 `PUT` latency to ingest path and `MaxKeys` caps (`LANGFUSE_S3_LIST_MAX_KEYS=200` at `packages/shared/src/env.ts:292`) to list scans.

- **Single OTLP exporter URL vs. fan-out**: one `OTLPTraceExporter` instance keeps config simple and matches most deployments (single collector), but precludes simultaneous Honeycomb + Langfuse + local-file sinks without an external Collector tier.

- **Batch blob export vs. live stream**: `every_20_minutes` minimum (`fern/apis/server/definition/blob-storage-integrations.yml:57`) is efficient for warehouse ETL (lag buffer avoids late-arriving spans) but unsuitable for real-time alerting/metrics that need per-span OTLP.

- **Hand-managed proto artifact**: `otlp-proto/generated/root.ts` checked in with `README.md:7` warning not to hand-edit and no `pnpm generate` script — guarantees stable ingest contract but risks drift if OTLP spec revs and no CI regeneration check fires.

## Failure Modes / Edge Cases

- **Oversized spans**: `LANGFUSE_OTEL_MAX_SPAN_BYTES=9_500_000` (`packages/shared/src/env.ts:159`) triggers `logOversizedSpan` warning and `langfuse.ingestion.otel.oversized_span` metric but does not reject — ClickHouse `10MB min_chunk_bytes_for_parallel_parsing` can still stall parallel parsing on pathological `gen_ai.input.messages` arrays (`packages/shared/src/server/otel/OtelIngestionProcessor.ts:1235-1283`).

- **Invalid protobuf/JSON**: `application/x-protobuf` decode failure returns `400 Failed to parse OTel Protobuf Trace` (`web/src/pages/api/public/otel/v1/traces/index.ts:91-104`), JSON parse failure returns `400 Failed to parse OTel JSON Trace` (`web/src/pages/api/public/otel/v1/traces/index.ts:106-114`); malformed `resourceSpans` with zero scopes returns `{}` (empty success) rather than error, which can mask client bugs (`web/src/pages/api/public/otel/v1/traces/index.ts:116-118`).

- **Ingestion version drift**: `x-langfuse-ingestion-version >4` rejected with `400 Unsupported` (`web/src/pages/api/public/otel/v1/traces/index.ts:147-160`); no contract test enforces that new SDK major bumps stay ≤4 before merge, so future SDK release can brick ingest until web is upgraded.

- **S3 unavailability**: `publishToOtelIngestionQueue` throws and caller calls `markProjectIngestFailure({source: public_otel_api, reason: publish_failed})` (`web/src/pages/api/public/otel/v1/traces/index.ts:190-195`) but the HTTP response is an unhandled exception (500) — client gets no `retry-after` and no staged `fileKey` to replay.

- **Queue saturation / Redis loss**: BullMQ job caps `removeOnFail:100_000` (`packages/shared/src/server/redis/otelIngestionQueue.ts:71`) mean sustained failure (>100k jobs) silently drops earliest failures; `markProjectAsOtelUser` Redis tracking (`packages/shared/src/server/redis/otelProjectTracking.ts`) and `recentProjectMarks` in-memory rate-limit analogs are lost on worker restart, causing duplicate failure marks post-restart.

- **Field-group drift**: `exportSource=OBSERVATIONS_V2` requires v4 enriched-observations table (Cloud-only, `fern/apis/server/definition/blob-storage-integrations.yml:71-72`); self-hosted with `LANGFUSE_MIGRATION_V4_WRITE_MODE=legacy` silently exports empty enriched portion with no error.

- **Gzip + content-type mismatch**: `content-encoding: gzip` is checked on raw `Buffer` before content-type sniff (`web/src/pages/api/public/otel/v1/traces/index.ts:64-89`); a client sending `application/x-protobuf` with `content-type: text/javascript; charset=utf-8` (observed in strict-matching comment at `web/src/pages/api/public/otel/v1/traces/index.ts:80`) would have passed after the `includes()` fix, but a client sending protobuf *without* an OTLP content-type still correctly gets 400.

## Future Considerations

- **Add a reverse mapper + live OTLP forwarder for stored traces**: `OtelIngestionProcessor` already embodies forward mapping; a `LangfuseObservationToOtelSpanMapper` plus a worker `otelExportQueue` with per-project `exporters[]` (similar to `OtelIngestionSecondaryQueue` sharding) would enable real-time fan-out to Honeycomb/Tempo/Datadog without requiring every consumer to re-implement `observation.type` → OTel `SpanKind` translation. Gate with `LANGFUSE_OTEL_EXPORT_ENABLED` and reuse existing `StorageService` credentials pattern.

- **Multi-destination blob export**: extend `BlobStorageIntegrationsResponse {data: list<BlobStorageIntegrationResponse>}` (`fern/apis/server/definition/blob-storage-integrations.yml:176-178`) to allow multiple active integrations per `projectId` with per-integration `exportFieldGroups` filters — e.g., `tools+usage` to S3 for finops, `io+metadata` to GCS for RAG eval — rather than single `exportSource` per project.

- **Local-file / collector-file sink**: add `BlobStorageIntegrationType.FILE` (or `S3_COMPATIBLE` with `endpoint: file:///exports`) backed by the existing `StorageService` abstraction, plus a `trace.export(format="otlp-json"|"parquet", path=...)` CLI that reuses `getTracesFromEventsTableForPublicApi` to materialize recent traces for offline replay/testing.

- **Exporter plugin registry**: replace single `OTLPTraceExporter` construction (`web/src/observability.config.ts:114-116`) with a `TracerProvider` middleware chain (similar to `ScopeFilteringTracerProvider`) that loads `OTEL_EXPORTER_*` multi-exporter lists via `OTEL_EXPORTER_OTLP_HEADERS`-style env and supports in-process `SpanExporter` injection for custom sinks — tested via a `MultiTracingProcessor` analog instead of external Collector.

- **Stronger export observability**: expose exporter health (queue depth, `langfuse.queue.{otel_ingestion,secondary_otel_ingestion}` lag, `otel.conversion_failure`/`oversized_span` rates already emitted at `packages/shared/src/server/otel/OtelIngestionProcessor.ts:143,331`) as a public `GET /api/public/integrations/blob-storage/{id}` `syncStatus` enrichment and as OTLP `export_error_callback` metrics so operators detect downstream 401/429 without tailing worker logs.

- **Versioned export format**: stamp blob exports and public API trace payloads with `exportFormatVersion` (mirroring `x-langfuse-ingestion-version` gate) and publish JSON Schema for the export envelope, enabling zero-drift migration between `LEGACY_TRACES_OBSERVATIONS` and `OBSERVATIONS_V2` consumers.

## Questions / Gaps

- **Reverse (Langfuse → OTel) mapping completeness**: no evidence of a tested mapper from `observation.type` (`SPAN|GENERATION|AGENT|TOOL|CHAIN|RETRIEVER|EVALUATOR|EMBEDDING|GUARDRAIL` at `packages/shared/src/server/ingestion/types.ts:259-279`) back to OTel `SpanKind`+`gen_ai.*` attributes. Is this intentionally left to consumers, or planned as `packages/shared/src/server/otel/langfuseToOtelMapper.ts`?

- **Multi-sink fan-out intent**: does the `secondary-otel-ingestion-queue` (`packages/shared/src/server/queues.ts:353`) prefigure multi-destination export, or only noisy-neighbor isolation? Search for `secondary-otel` found only queue definitions and sharding (`packages/shared/src/server/redis/otelIngestionQueue.ts:91-175`, `worker/src/queues/otelIngestionQueue.ts:236-267`), no fan-out logic.

- **Honeycomb/LangSmith native adapters**: no code references `honeycomb|langsmith|Honeycomb|LangSmith` in `packages/shared/src/server` or `web/src/pages/api/public/otel`; are these expected to be accessed purely via generic OTLP Collector, or is a first-party `HoneycombExporter` / `LangSmithExporter` on the roadmap? (OTLP generic works, but header/convention helpers like `x-honeycomb-team` are not documented.)

- **Generated proto lifecycle**: who regenerates `web/src/pages/api/public/otel/otlp-proto/generated/root.ts` on OTLP spec revs? `fern/fern.config.json` governs Fern → OpenAPI clients, not proto → TS; no `package.json` script targets `pbjs/pbts` for `root.ts`.

- **Local-file observability for air-gapped/self-hosted**: can `LANGFUSE_S3_EVENT_UPLOAD_BUCKET` be pointed at a local filesystem (e.g., MinIO `S3_COMPATIBLE` with `forcePathStyle`) and also serve as a user-visible trace dump, or is the raw `otel/{projectId}/{yyyy}/{mm}/{dd}/{hh}/{mm}/{uuid}.json` bucket considered internal-only with no retention/TTL policy exposed to operators? (`worker/src/scripts/replayIngestionEventsV2/README.md:135-143` treats it as internal.)

---

Generated by `dimensions/10.04-export-interoperability-observability.md` against `langfuse`.
