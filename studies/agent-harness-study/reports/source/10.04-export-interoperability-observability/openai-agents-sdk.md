# Source Analysis: openai-agents-sdk

## 10.04 Export, Interoperability, and Observability Backends

### Source Info

| Field | Value |
|-------|-------|
| Name | openai-agents-sdk |
| Path | `studies/agent-harness-study/sources/openai-agents-sdk` |
| Language / Stack | Python (Agents SDK, `src/agents`) |
| Analyzed | 2026-08-28 |

## Summary

`openai-agents-sdk` implements a proprietary trace export path by default (`{"object":"trace"}` / `{"object":"trace.span"}` JSON to `https://api.openai.com/v1/traces/ingest`) and makes interoperability an explicit code-seam rather than a wire protocol. The public seam is `TracingProcessor`/`TracingExporter` (`src/agents/tracing/processor_interface.py:9-142`), multiplexed in-process by `SynchronousMultiTracingProcessor` (`src/agents/tracing/provider.py:93-220`) and batched by `BatchTraceProcessor` (`src/agents/tracing/processors.py:541-718`). There is no built-in OTLP exporter, no `Langfuse`/`LangSmith`/`Honeycomb` SDK wrapper, and no file/OTLP/local-file exporter beyond `ConsoleSpanExporter` (`src/agents/tracing/processors.py:27-41`). OTel reach is entirely delegated to the 27 community adapters enumerated in `docs/tracing.md:196-226` (including Langfuse, LangSmith, HoneyHive/Honeycomb, Datadog, W&B, Arize-Phoenix, MLflow, Braintrust, Logfire, etc.). Runtime configurability exists for disabling, API key, endpoint, batch/queue sizing, and processor set, but there is no env-var-only exporter selection; pointing at a non-OpenAI collector requires code (`set_trace_processors`/`add_trace_processor`/`set_trace_provider`).

## Rating

**5/10 — Present but inconsistent / weakly portable.**

Rationale: Export is operationally mature for the vendor endpoint (batched queue, scheduled + high-watermark flush, exponential backoff with jitter, deadline-aware shutdown, `atexit` handler, per-key API key grouping, endpoint-aware sanitization/truncation, and isolation of exporter failures from the agent loop). Interoperability, however, is adapter-only: no native OTLP/HTTP or gRPC exporter (`grep` for `opentelemetry` in `src/agents` returns only `uv.lock` dev stubs), wire format is custom (`src/agents/tracing/spans.py:396-422`, `src/agents/tracing/traces.py:568-575`) without a version field, and simultaneous multi-backend is structurally supported but not exercised out-of-box. Fits the rubric prompt directly: you cannot send traces to an existing OTel/Honeycomb stack without writing (or installing) an adapter.

## Evidence Collected

Every entry MUST include a file path with line numbers. Format: `path/to/file.ts:NN`.

| Area | Evidence | File:Line |
|------|----------|-----------|
| OTLP exporter (absent) | No `opentelemetry` import or OTLP exporter class in `src/agents/tracing`; `BatchTraceProcessor` and `BackendSpanExporter` are the only concrete processors. `uv.lock` pins `opentelemetry-* 1.40.0` only as transitive/dev deps, not imported by SDK code. | `src/agents/tracing/processors.py:1-763`, `src/agents/tracing/processor_interface.py:1-142`, `uv.lock:750-753` |
| Proprietary ingest exporter | `BackendSpanExporter` defaults to `https://api.openai.com/v1/traces/ingest` with `OpenAI-Beta: traces=v1`, Auth via `OPENAI_API_KEY` / `OPENAI_ORG_ID` / `OPENAI_PROJECT_ID`, per-item-traces grouping by `tracing_api_key`, retry with backoff. | `src/agents/tracing/processors.py:44-90`, `src/agents/tracing/processors.py:118-221`, `src/agents/tracing/processors.py:106-117` |
| Endpoint configurability | Constructor `endpoint` param (default ingest) allows pointing default exporter at a custom HTTP endpoint; sanitization/truncation is gated on `endpoint == _OPENAI_TRACING_INGEST_ENDPOINT`. | `src/agents/tracing/processors.py:62,258-260`, `src/agents/tracing/processors.py:261-319` |
| Console / debug exporter | `ConsoleSpanExporter.export()` prints redacted or `trace_id`/`span_data` to stdout; used as minimal local sink. | `src/agents/tracing/processors.py:27-42` |
| Custom sink interface | `TracingProcessor` ABC (5 lifecycle methods) and `TracingExporter` ABC (`export(items: list[Trace|Span])`) documented as extension point. | `src/agents/tracing/processor_interface.py:9-130`, `src/agents/tracing/processor_interface.py:132-142` |
| Multi-backend fan-out | `SynchronousMultiTracingProcessor` holds `tuple[TracingProcessor]` under lock, forwards `on_*`/`force_flush`/`shutdown` to every entry, swallowing per-processor exceptions with diagnostic extra. | `src/agents/tracing/provider.py:93-220` |
| Bathing + background export | `BatchTraceProcessor` with `max_queue_size=8192`, `max_batch_size=128`, `schedule_delay=5.0`, `export_trigger_ratio=0.7`, lazy daemon thread, `queue.Queue`, `force_flush()` and `shutdown(timeout)` with deadline propagation. | `src/agents/tracing/processors.py:541-718` |
| Default wiring | Lazy singletons `_global_exporter`/`_global_processor` and `get_trace_provider()` that registers `default_processor()` on first access; `atexit` handler shuts down `DefaultTraceProvider` with 5s timeout. | `src/agents/tracing/processors.py:720-763`, `src/agents/tracing/setup.py:11-66` |
| Commercial platform support | No vendored Langfuse/LangSmith/Honeycomb SDK; `docs/tracing.md:196-226` lists ~27 external `TracingProcessor` integrations (W&B, Arize-Phoenix, Future AGI, MLflow, Braintrust, Logfire, AgentOps, Scorecard, Respan, LangSmith, Maxim AI, Comet Opik, Langfuse, Langtrace, Monocle, Galileo, Portkey AI, LangDB AI, Agenta, PostHog, Traccia, PromptLayer, HoneyHive, Asqav, Datadog, Latitude). | `docs/tracing.md:196-226` |
| Env-var runtime config | `OPENAI_AGENTS_DISABLE_TRACING` disables all traces (`DefaultTraceProvider._refresh_disabled_flag`), `OPENAI_API_KEY`/`OPENAI_ORG_ID`/`OPENAI_PROJECT_ID` supply default exporter auth, `OPENAI_AGENTS_TRACE_INCLUDE_SENSITIVE_DATA` controls `tool_generation` payload inclusion. | `src/agents/tracing/provider.py:346-356`, `src/agents/tracing/processors.py:106-117`, `src/agents/run_config.py:53-56` |
| Code-only processor selection | `agents.tracing.add_trace_processor()`, `set_trace_processors()`, `set_trace_provider()`, `set_tracing_export_api_key()`, `flush_traces()` are code APIs; no env-var maps to an alternate exporter class. | `src/agents/tracing/__init__.py:94-130`, `src/agents/tracing/setup.py:27-36` |
| Per-run tracing override | `RunConfig.tracing: TracingConfig | None` and `TracingConfig.api_key` allow per-run key; `trace(tracing=...)` and `RunConfig.workflow_name/trace_id/group_id/trace_metadata` flow through `DefaultTraceProvider.create_trace`. | `src/agents/run_config.py:401-402`, `src/agents/tracing/config.py:6-10`, `src/agents/tracing/provider.py:374-408` |
| Trace format schema | `TraceImpl.export()` -> `{"object":"trace","id":trace_id,"workflow_name","group_id","metadata"}` without version; `SpanImpl.export()` -> `{"object":"trace.span","id":span_id,"trace_id","parent_id","started_at","ended_at","span_data":<SpanData.export()>,"error"}` plus opt-in metadata routing (`agent_harness_id`). | `src/agents/tracing/traces.py:568-575`, `src/agents/tracing/spans.py:396-422` |
| SpanData taxonomy | 13 concrete `SpanData` subtypes (`AgentSpanData`, `TaskSpanData`, `TurnSpanData`, `FunctionSpanData`, `GenerationSpanData`, `ResponseSpanData`, `HandoffSpanData`, `CustomSpanData`, `GuardrailSpanData`, `TranscriptionSpanData`, `SpeechSpanData`, `SpeechGroupSpanData`, `MCPListToolsSpanData`) each with `type` discriminator and `export()` dict. | `src/agents/tracing/span_data.py:28-452` |
| Usage field (generation spans) | `GenerationSpanData` carries `input`/`output`/`model`/`model_config`/`usage:{input_tokens,output_tokens,details}`; `BackendSpanExporter` sanitizes to allowed keys only for OpenAI endpoint. | `src/agents/tracing/span_data.py:169-210`, `src/agents/tracing/processors.py:47-52,448-483` |
| File / local export | No file exporter; closest is `ConsoleSpanExporter`; trace persistence for resume uses `TraceState` JSON (with secret hashing) but is run-state, not observability export. | `src/agents/tracing/processors.py:27-42`, `src/agents/tracing/traces.py:195-277` |
| Tests for export behavior | `test_batch_trace_processor_*` (queue full drops, force_flush, scheduled export, shutdown deadline, retry backoff interruption, exporter exception survival), `test_backend_span_exporter_*` (2xx/4xx/5xx retry, request error, no key skip, truncation, usage sanitization, custom-endpoint bypass). | `tests/test_trace_processor.py:55-944`, `tests/test_tracing.py:74-530` |
| Failure isolation | `_export_batches` wraps `exporter.export` in try/except (`log_model_and_tool_action_error`) so one exporter failure does not kill worker thread; `SynchronousMultiTracingProcessor` wraps each forwarder in try/except. | `src/agents/tracing/processors.py:699-717`, `src/agents/tracing/provider.py:122-175` |
| Shutdown / flush safeguards | `flush_traces()` delegates to `get_trace_provider().force_flush()`; `BatchTraceProcessor.shutdown(timeout)` joins worker or drains synchronously; `BackendSpanExporter._export_with_deadline` respects deadline for retries and `shutdown_event`. | `src/agents/tracing/__init__.py:122-130`, `src/agents/tracing/processors.py:623-669`, `src/agents/tracing/processors.py:121-256` |
| Dependency surface | `pyproject.toml:9-17` has no `opentelemetry-*` runtime dep; extras are sandbox/session/cloud backends only, confirming OTel is intentionally external. | `pyproject.toml:9-59` |

## Answers to Dimension Questions

**1. Can traces be exported to external backends?** Yes, but only through the code seam. `add_trace_processor()` appends a processor to the default OpenAI pipeline (`src/agents/tracing/__init__.py:94-98`), while `set_trace_processors()` (`src/agents/tracing/__init__.py:101-105`) or `set_trace_provider()` (`src/agents/tracing/setup.py:27-36`) replaces it. `docs/tracing.md:142-152` documents both patterns ("additional" vs "replacement"). The integrations table at `docs/tracing.md:196-226` enumerates ~27 vendor/community processors (Langfuse, LangSmith, Datadog, HoneyHive/Honeycomb, PostHog, etc.) proving the seam is used, but none ships in `src/agents`; you must install the third-party package and wire it in code.

**2. Are standard protocols supported?** No native standard. The built-in `BackendSpanExporter` speaks a proprietary JSON envelope (`src/agents/tracing/traces.py:568-575`, `src/agents/tracing/spans.py:396-422`) POSTed to `https://api.openai.com/v1/traces/ingest` with `OpenAI-Beta: traces=v1` (`src/agents/tracing/processors.py:44-45`). There is no `OtelSpanExporter`, OTLP gRPC/HTTP, or `OTEL_EXPORTER_OTLP_*` handling in `src/agents` (grep for `opentelemetry` returns only `uv.lock:750-753`). Standard protocols are only reachable via community adapters that map the proprietary dict to OTLP/Langfuse/LangSmith APIs themselves.

**3. Is export configurable without code changes?** Partially. Disabling is fully env-driven (`OPENAI_AGENTS_DISABLE_TRACING=1` at `src/agents/tracing/provider.py:347-352`, or `RunConfig.tracing_disabled` at `src/agents/run_config.py:397-399`). API key can be supplied via `OPENAI_API_KEY` env fallback (`src/agents/tracing/processors.py:108`), via `set_tracing_export_api_key()` (`src/agents/tracing/__init__.py:115-119`), or per-run `TracingConfig.api_key` (`src/agents/tracing/config.py:9`). Endpoint is configurable only by constructing a new `BackendSpanExporter(endpoint=...)` in code (`src/agents/tracing/processors.py:62`) and replacing the processor set; there is no `TRACE_ENDPOINT` env that selects a handler class. Choosing a non-OpenAI backend always requires a code change to register a processor.

**4. Can multiple backends receive traces simultaneously?** Yes. `SynchronousMultiTracingProcessor` (`src/agents/tracing/provider.py:93-220`) broadcasts every `on_trace_start`/`on_trace_end`/`on_span_start`/`on_span_end`/`force_flush`/`shutdown` call to a tuple of processors. The default trace provider holds one `SynchronousMultiTracingProcessor` (`src/agents/tracing/provider.py:302`). `add_trace_processor` (`src/agents/tracing/__init__.py:94-98`) demonstrates the intended additive model (OpenAI default + custom), and `set_trace_processors([...])` can install N exporters at once. The fan-out is synchronous and per-processor-exception-isolated, but the SDK ships with a size-1 default list; multi-sink is supported infrastructure, not out-of-box configuration.

## Architectural Decisions

- **Proprietary dict + HTTP ingest over OTLP** (`src/agents/tracing/processors.py:44-46`, `src/agents/tracing/span_data.py:28-452`). Choice optimizes for zero-dep default (`pyproject.toml:9-17` lists only `openai`, `pydantic`, `requests`, `mcp`) and deep OpenAI dashboard integration (`docs/tracing.md:1-4`), at the cost of requiring translation for every third-party backend. Tradeoff is explicit in `docs/tracing.md:194-226` delegating to community processors.

- **Processor / Exporter two-level abstraction** (`src/agents/tracing/processor_interface.py:9-142`). `TracingProcessor` handles lifecycle and batching policy; `TracingExporter` handles wire serialization. Allows `BatchTraceProcessor(TracingExporter)` composition (`src/agents/tracing/processors.py:548-555`). Keeps agent loop decoupled from HTTP details.

- **Synchronous broadcast with isolated failures** (`src/agents/tracing/provider.py:117-175`, `src/agents/tracing/processors.py:699-717`). Guarantees order and simplicity (no async queue per processor), but a slow exporter blocks the shared `on_span_end` path; a failing one is logged and dropped rather than retried at processor level.

- **Batched background export with lazy thread start** (`src/agents/tracing/processors.py:584-596`, `src/agents/tracing/processors.py:653-669`). Minimizes agent-loop overhead, supports `schedule_delay` + `export_trigger_ratio` dual triggers, and avoids creating threads at import time (important for fork models per comment at `src/agents/tracing/processors.py:720-722`).

- **Endpoint-aware sanitization** (`src/agents/tracing/processors.py:258-291`). Truncation and usage-field filtering (`input`/`output` >100k bytes, `usage` to `input_tokens`/`output_tokens`+`details`) only apply when `endpoint == _OPENAI_TRACING_INGEST_ENDPOINT`. Keeps custom endpoints verbatim (`tests/test_trace_processor.py:804-839` asserts custom endpoint preserves large input/usage).

- **Global mutable singleton with `atexit` 5s shutdown** (`src/agents/tracing/setup.py:11-66`, `src/agents/tracing/provider.py:506-511`). Gives zero-config out-of-box tracing at cost of global mutable state and `pytest` conftest needing to reset provider (`tests/conftest.py:43`).

## Notable Patterns

- **Code-seam ecosystem pattern.** Clean ABCs plus a single-line global hook (`add_trace_processor`) enabled 27 integrations without SDK changes (`docs/tracing.md:196-226`). Mirrors `langfuse` OTel vs harness native distinction noted in cross-source synthesis.
- **Lazy defaults to avoid import side effects.** `default_exporter()`/`default_processor()` use double-checked locking (`src/agents/tracing/processors.py:727-763`) and `get_trace_provider()` lazily registers processor (`src/agents/tracing/setup.py:39-66`) so importing `agents` does not open HTTP clients or threads.
- **Secret handling via key grouping + `TraceState` hashing.** `BackendSpanExporter` groups batched items by `tracing_api_key` (`src/agents/tracing/processors.py:125-128`) and `TraceState.to_json(include_tracing_api_key=...)` only persists raw key on explicit opt-in, otherwise storing `_hash_tracing_api_key` (`src/agents/tracing/traces.py:187-277`).
- **Sensitive-data gating.** `RuleConfig.trace_include_sensitive_data` defaults from `OPENAI_AGENTS_TRACE_INCLUDE_SENSITIVE_DATA` env (`src/agents/run_config.py:53-56`), `ConsoleSpanExporter` checks `_debug.DONT_LOG_MODEL_DATA/_TOOL_DATA` (`src/agents/tracing/processors.py:32-37`), and `SpanImpl.export` still exports span envelope when redacted.
- **Defensive deadline / shutdown coordination.** `BatchTraceProcessor.shutdown(timeout)` passes `time.monotonic()+timeout` as deadline to `_export_batches`/`_export_with_deadline` (`src/agents/tracing/processors.py:633-634`, `src/agents/tracing/provider.py:181-183`), and `BackendSpanExporter._sleep_before_retry` checks `shutdown_event` and deadline to preserve process exit code (`tests/test_trace_processor.py:558-648` regression test for exit code 7 on 504).

## Tradeoffs

- **Zero vendor lock-in at build time, lock-in at runtime unless you write an adapter.** No `opentelemetry-*` hard dep keeps install lean, but any OTel collector, Langfuse, or LangSmith requires a community package. The docs themselves frame tracing as "free traces at OpenAI Traces dashboard" (`docs/tracing.md:3-4`) with third-party as afterthought.
- **Synchronous fan-out is simple but couples latency.** One slow `TracingProcessor` delays `on_span_end` for all subsequent processors; `BatchTraceProcessor` isolates exporter latency on its worker thread, yet processor-level sync remains (`src/agents/tracing/provider.py:162-175`).
- **Vendor tuning vs portability.** Truncation, byte limits, and `usage` allowlist (`src/agents/tracing/processors.py:46-55`) only apply to the OpenAI endpoint, preserving fidelity for custom backends but meaning the same span serializes differently per endpoint—subtle source of divergence in tests.
- **Silent drop on overload.** `queue.Full` -> `logger.warning("Queue is full, dropping trace/span")` (`src/agents/tracing/processors.py:604,621`) avoids backpressure on agent loop but loses observability exactly when most needed (high throughput). Verified by `test_batch_trace_processor_queue_full` (`tests/test_trace_processor.py:77-91`).
- **No file/OTLP local sink.** Beyond `ConsoleSpanExporter`, there is no `FileTraceExporter` or local-file-as-sink; offline debug or air-gapped export requires writing a processor even though the tests use an in-memory `SpanProcessorForTests` (`tests/testing_processor.py:12`).

## Failure Modes / Edge Cases

- **No API key -> silent skip.** `BackendSpanExporter.export` warns `OPENAI_API_KEY is not set, skipping trace export` and returns without enqueuing error (`src/agents/tracing/processors.py:132-134`), tested at `tests/test_trace_processor.py:414-423`. Traces are lost without exception.
- **4xx vs 5xx split: never retry client errors, always retry server/request errors up to `max_retries` with jitter.** 4xx logs redacted or full body based on `DONT_LOG_*` and stops (`src/agents/tracing/processors.py:185-198`); 5xx/`RequestError` retries with `delay *2` capped at `max_delay` plus 10% jitter (`src/agents/tracing/processors.py:200-221`). Verified by `test_backend_span_exporter_4xx_client_error`, `5xx_retry`, `request_error` (`tests/test_trace_processor.py:440-663`).
- **Exporter exception kills no state only because it is caught.** Prior to fix, exception in `export` killed the background worker; now `_export_batches` catches `Exception` and logs (`src/agents/tracing/processors.py:699-717`), tested by `test_batch_trace_processor_survives_exporter_exception` which asserts worker stays alive after a `RuntimeError` and that subsequent spans still export (`tests/test_trace_processor.py:235-271`).
- **Shutdown deadline respected across retries.** `BackendSpanExporter._export_with_deadline` checks `deadline` each loop and sleeps with deadline awareness (`src/agents/tracing/processors.py:121-256`); `BatchTraceProcessor.shutdown(timeout=0.05)` triggers `warning shutdown timeout reached` and returns quickly even if exporter blocks with `max_retries=100, base_delay=10.0` (`tests/test_trace_processor.py:180-212`, `558-648`), preserving process exit code.
- **Unserializable / oversized payloads.** `_sanitize_json_compatible_value` drops non-JSON values (including cyclic refs) and keeps only `str`/int/float finite, `bool` (`src/agents/tracing/processors.py:492-531`); oversize `input`/`output` are intelligently truncated by JSON-byte budget preserving list/dict shape or replaced with `{"truncated":True,"original_type":...}` for non-stringifiable types (`src/agents/tracing/processors.py:355-446`), tested exhaustively in `tests/test_trace_processor.py:676-1301`.
- **Stale env vs manual override.** `DefaultTraceProvider._refresh_disabled_flag` reads `OPENAI_AGENTS_DISABLE_TRACING` once and then `manual_disabled` takes precedence (`src/agents/tracing/provider.py:339-356`); mid-run env changes are ignored after `set_tracing_disabled` is called.
- **No version in export envelope.** Both `TraceImpl.export` (`src/agents/tracing/traces.py:568-575`) and `SpanImpl.export` (`src/agents/tracing/spans.py:396-422`) emit unversioned dicts; consumers cannot detect breaking schema changes—gap noted in every recent final synthesis.

## Future Considerations

- Add a first-party `OtelSpanExporter` mapping `SpanData.type` -> OTel `SpanKind`/`Status`/`Attributes` (with `gen_ai.*` semantic conventions) as an opt-in `opentelemetry` extra, reusing the existing `BatchTraceProcessor` composition. Current community adapters each reinvent this mapping.
- Emit a `version`/`schema_version` field from `TraceImpl.export`/`SpanImpl.export` (and `TraceState.to_json`) and gate ingestion, mirroring `RunState` versioning elsewhere in the repo, to make downstream parsers resilient.
- Provide a `FileTraceExporter(path, format="jsonl"|"otlp-json")` mirroring `ConsoleSpanExporter` for local/air-gapped use and deterministic replay, without requiring users to write `queue` handling.
- Expose an env-driven exporter selector (e.g., `OPENAI_AGENTS_TRACE_EXPORTER=openai|otel|console|none` + `OPENAI_AGENTS_TRACE_ENDPOINT`) so operators can point at a collector without code changes, while keeping code seam for advanced cases.
- Define an explicit local filesystem export for traces as well.
- Surface queue-dropped-span metrics (counter + last-drop timestamp) via `logger` or optional callback so overload is observable rather than a single warning line.
- Consider bounded per-processor concurrency or opt-in async fan-out in `SynchronousMultiTracingProcessor` to isolate a slow third-party processor from agent tail latency.

## Questions / Gaps

- No evidence found that `src/agents/tracing` supports `OTEL_EXPORTER_OTLP_*` env vars or an OTLP wire format out-of-box; `grep` across `src/agents` for `opentelemetry|OTLP|OTEL` returned zero hits aside from `uv.lock` dev pins. Is a first-party OTel exporter on the roadmap or intentionally delegated to community processors?
- Does `BackendSpanExporter` with a custom `endpoint` honor any auth mechanism other than bearer + optional `OpenAI-Organization`/`Project` headers? No alternative auth path exists in `src/agents/tracing/processors.py:146-157`.
- Is long-running-worker guidance (`flush_traces()` in `docs/tracing.md:49-97`) expected to replace push-based config reloading of processors, or will there be a `RELOAD_ON_SIGHUP` style mechanism for `set_trace_processors` without restart?
- Trace file export and local replay: `FileHistoryProvider` is not wired as a tracing processor and `tests/testing_processor.py:52` only normalizes in-memory spans; is file persistence of traces intended to stay community-composed?
- The global singleton (`GLOBAL_TRACE_PROVIDER` at `src/agents/tracing/setup.py:11`) plus `atexit` makes testing correct but production fork safety (gunicorn pre-fork, Celery worker) dependent on lazy init ordering; is a `fork` hook or `os.register_at_fork` integration planned?

---
Generated by `dimensions/10.04-export-interoperability-observability.md` against `openai-agents-sdk`.
