> **Inputs Used:** `projects/ultraplan-go/project-index.md`, `projects/ultraplan-go/sprints/39-performance-stage/requirements.md`, `projects/ultraplan-go/sprints/39-performance-stage/sprint-index.md`, `projects/ultraplan-go/sprints/39-performance-stage/technical-handbook.md`, `projects/ultraplan-go/docs/ARCHITECTURE.md`, `projects/ultraplan-go/docs/PRD.md`, `projects/ultraplan-go/docs/TRD.md`, `system/reasoning/api-design-reasoning-template.md`, `system/contracts/surfaces/api-contracts.md`, `system/contracts/surfaces/cli.md`, `system/contracts/core/errors.md`, `system/contracts/core/observability.md`, `system/contracts/core/security.md`, `system/contracts/runtime/workflows.md`, `system/contracts/runtime/performance.md`, `system/contracts/runtime/persistence-and-migrations.md`, `internal/app/operations.go`, `internal/app/sprint_usecases.go`, `internal/web/routes.go`, `internal/runcontrol/model.go`, `internal/runcontrol/interfaces.go`

# API Design: Requirements-Driven Performance Stage

This area covers the application, CLI, TUI, and local HTTP contracts for preparing, starting, observing, cancelling, resuming, recovering, and inspecting performance verification. It does not define the private measurement record schemas or the performance algorithm. Those remain owned by `internal/sprint`.

## Area Decisions

### 1. Add one typed performance use-case boundary

`internal/app/sprint_usecases.go` will expose a `PerformanceUseCases` interface and transport-neutral request and result types. CLI, TUI, and web adapters call this boundary. They must not call `internal/sprint` directly, parse CLI text, inspect private performance JSON, or derive a target or run verdict.

The boundary will cover these operations:

| Use case | Execution model | Writes | Result |
| --- | --- | --- | --- |
| `PreparePerformance` | synchronous and runtime-free | none | deterministic admission and confirmation facts |
| `DryRunPerformance` | synchronous and runtime-free | none | the same admission facts plus planned work and missing prerequisites |
| `StartPerformance` | durable asynchronous run | private state, canonical summaries, and bounded validated promotions | accepted operational run |
| `PerformanceStatus` | synchronous query | none | current bounded phase snapshot |
| `ResumePerformance` | durable asynchronous run | continuation of the existing attempt only | accepted operational run |
| `CancelPerformance` | synchronous durable request | cancellation intent and later terminal state | current run and cancellation facts |
| `RecoverPerformance` | synchronous and runtime-free | conservative state reconciliation | reconciled status |
| `PerformanceResult` | synchronous query | none | immutable terminal product result when present |
| `PerformanceEvidence` | synchronous paged query | none | bounded evidence metadata and allowlisted references |

"Prepare" at this boundary means request admission and confirmation preparation. Benchmark discovery, benchmark authoring, and manifest freezing are steps inside the accepted performance workflow. This distinction keeps prepare and dry-run runtime-free as required by AC-3.

The existing generic operation bridge remains the entry point for guarded work. Add performance operation kinds for `performance-dry-run`, `performance-start`, `performance-status`, `performance-resume`, and `performance-recover`. Cancellation continues through the canonical durable-run cancellation path, while the typed performance method verifies that the selected run belongs to the requested sprint and phase. Result and evidence reads are direct queries rather than fake operations.

### 2. Keep product status separate from operational lifecycle

Every accepted start or resume has an `internal/runcontrol` run ID. The performance attempt also has a sprint-owned attempt ID. Both identifiers are returned, but they mean different things:

| Identifier or state | Owner | Meaning |
| --- | --- | --- |
| `operation_run_id` | `internal/runcontrol` | acceptance, owner, lease, events, cancellation, and operational terminal arbitration |
| `operational_attempt_id` | `internal/runcontrol` | one fenced owner claim for the run |
| `performance_attempt_id` | `internal/sprint` | immutable target, benchmark, measurement, optimization, and result evidence family |
| `cycle_number` | `internal/sprint` | one bounded optimization cycle within a performance attempt |
| `target_id` | current `requirements.md` | stable `PERF-` row identity |

The operation lifecycle uses the existing values such as `accepted`, `queued`, `running`, `cancelling`, `succeeded`, `failed`, `cancelled`, `interrupted`, and `cleanup_uncertain`. Product run outcomes use the required values `passed`, `passed_with_reports`, `target_missed`, `blocked`, `cancelled`, `cleanup_uncertain`, and `stalled`. Target outcomes use `met`, `missed`, `baseline_recorded`, `report_only`, `inconclusive`, and `blocked`.

Operational `succeeded` means the runner reached and committed one valid product terminal result. It does not mean that performance passed. For example, a completed operation whose product outcome is `target_missed` is operationally succeeded and still blocks Conformance Review. Adapters display both fields and never collapse them into one status.

### 3. Use one bounded, versioned public fact model

`PerformanceStatusResult` and `PerformanceResult` will be additive JSON DTOs with `schema_version: 1`. Their shared fields are:

- project and sprint identifiers;
- project policy mode;
- current performance phase and freshness boolean;
- bounded freshness reason codes;
- operation, operational-attempt, and performance-attempt identifiers when available;
- operational lifecycle, liveness, and fencing generation when available;
- target-packet, benchmark-manifest, environment-policy, correctness-policy, and implementation fingerprints when available;
- target totals grouped by the exact target outcomes;
- run outcome when terminal;
- consumed and maximum cycle, command, runtime-call, changed-file, changed-byte, output-byte, and wall-time limits;
- one bounded blocker with code, summary, and next action;
- canonical `performance.md`, current-state, and terminal-result references with workspace-relative path and digest;
- one explicit next action.

Each target summary contains its governed ID, scenario, metric, comparator, normalized value, unit, gate, samples, basis, qualification state, bounded baseline and candidate aggregates, comparison fact, outcome, and evidence reference. Optional numeric fields are absent when no qualified value exists. Unknown usage and cost values remain absent rather than becoming zero.

Public DTOs never include raw samples, full command output, profiles, patches, prompts, provider payloads, environment values, secrets, or private persistence structs. Human text, CLI JSON, TUI views, browser HTML, and browser JSON are renderings of the same app DTOs.

### 4. Freeze authority at acceptance and reject stale requests

Preparation normalizes the request and computes a confirmation fingerprint over:

- project and sprint identity;
- the normalized operation kind;
- project policy and current requirements bytes;
- execute evidence and target worktree identity;
- target packet identity when already available;
- correctness and process policy identities;
- effective lower-only limits and their sources;
- runtime and model identity for runtime-backed work;
- the existing performance attempt and last proven boundary for resume.

The server-issued confirmation contains only bounded facts, a non-reversible fingerprint, affected paths, mutation class, prerequisites, runtime/model source, and warning text. It never accepts a caller-supplied fingerprint as authority. `StartPerformance` and `ResumePerformance` reprepare immediately before durable acceptance. A mismatch returns `performance.request_stale` with conflict semantics and performs no runtime call, command, artifact write, or source mutation.

Start confirmation states that benchmark promotion and later implementation promotion are possible, both under product-owned validation. Resume confirmation grants no new scope. It binds the existing attempt, frozen identities, consumed counters, deadline, accepted mutations, and remaining authority. A caller cannot change targets, samples, commands, parsers, correctness gates, limits, runtime/model identity, or promotion scope during resume.

### 5. Define retry and idempotency per operation

- Prepare, dry-run, status, result, and evidence reads are naturally idempotent against unchanged inputs.
- Start uses the existing durable operation alias and confirmation digest. Repeating the same accepted request returns the existing run instead of creating another performance attempt.
- Resume deduplicates on the performance attempt, last proven boundary, and current confirmation fingerprint. One active writer wins. A second incompatible resume returns `performance.writer_conflict`.
- Cancel is idempotent. Repeating it returns the current snapshot and whether this request changed cancellation state.
- Recover is idempotent for the same current state identities. It may reconcile partial publication or dead ownership, but it cannot infer success, launch child work, create a proposal, or reset consumed limits.
- Runtime retries occur only inside the accepted workflow under persisted finite counters. HTTP, CLI, TUI, and browser adapters never retry a mutating performance step themselves.

### 6. Reuse existing HTTP operation and run resources

Do not add a second HTTP job registry or performance-specific event stream. The local web API uses the current versioned resources:

| Method and path | Purpose |
| --- | --- |
| `POST /api/v1/operations/prepare` | prepare confirmation for performance start, resume, dry-run, or recovery |
| `POST /api/v1/operations` | run a dry-run or recover synchronously, or durably accept start/resume |
| `GET /api/v1/runs/{run_id}` | inspect durable lifecycle and product correlation |
| `GET /api/v1/runs/{run_id}/events?after=<sequence>` | replay bounded sanitized progress |
| `DELETE /api/v1/runs/{run_id}` | request canonical idempotent cancellation |

Add sprint product queries rather than encoding performance facts into run-control records:

| Method and path | Purpose |
| --- | --- |
| `GET /api/v1/projects/{project}/sprints/{sprint}/performance` | current policy, status, freshness, target summaries, blocker, and next action |
| `GET /api/v1/projects/{project}/sprints/{sprint}/performance/result` | current immutable terminal result projection |
| `GET /api/v1/projects/{project}/sprints/{sprint}/performance/evidence` | paged evidence metadata filtered by attempt, target, cycle, or kind |

The evidence query uses an opaque cursor bound to the selected attempt and normalized filters. The default page size is 50 and the maximum is 200. Sort order is stable by attempt, cycle, evidence kind, and identifier. A cursor from another attempt or filter set returns `performance.cursor_stale`. Detailed evidence content remains behind the existing allowlisted bounded artifact-preview mechanism; the evidence endpoint returns metadata and references only.

Start and resume return HTTP `202 Accepted` with the durable run snapshot and URLs for run status, events, sprint performance status, and result. Synchronous successful queries, dry-run, cancel, and recovery return `200`. Known malformed requests return `400`, missing allowed resources return `404`, stale requests and ownership conflicts return `409`, semantically invalid target or admission inputs return `422`, capacity pressure returns `429` or `503` with a stable code, and unexpected failures return a sanitized `500`.

The local server keeps its current Host, Origin, CSRF, session, request-body, stream, and action-authorization controls. Browser disconnect, refresh, session expiry, or SSE loss only removes an observer. None cancels or completes the run.

### 7. Keep CLI, TUI, and browser behavior compatible

The CLI command remains:

```text
ultraplan sprint <project> <sprint> performance [--dry-run|status|resume|cancel|recover]
```

No subcommand may accept target values, sample overrides, command strings, parser names, environment values, arbitrary evidence paths, or expanded limits. Start and resume use the existing guarded-operation conventions. Non-interactive invocation must supply the existing explicit confirmation mechanism or fail before acceptance.

Text mode sends the final result to stdout and progress or diagnostics to stderr. JSON mode emits one stable app result on stdout with no ANSI or progress records. The command that runs or resumes work maps product outcomes as follows:

- `passed` and `passed_with_reports` return success;
- `target_missed`, `stalled`, and `cleanup_uncertain` return the documented partial or non-passing class;
- `cancelled` returns the cancellation class;
- `blocked` uses the underlying typed category, such as validation, dependency, filesystem, concurrency, or persistence.

Status, result, and evidence queries report whether the requested read succeeded; they include the stored product outcome without reusing that outcome as the query's exit status. TUI and browser controls call the same methods and show the same confirmation, progress, cancellation, blocker, freshness, and next-action facts.

### 8. Use stable typed errors and bounded events

Performance errors use `performance.<reason>` codes and retain their category, operation, component, safe cause, retryability, correlation IDs, and guidance through app and transport mapping. Required distinct reasons include policy disabled, target validation, admission blocked, request stale, writer conflict, benchmark coverage, parser uncertainty, noisy measurement, environment drift, benchmark drift, correctness failure, protected-path violation, limit exhausted, persistence failure, cancellation, cleanup uncertainty, and recovery conflict.

Durable events use the existing run-control event types. Their payloads carry bounded stable fields such as performance phase, target ID, cycle, completed and total counts, reason code, command identity digest, environment identity digest, and evidence reference. Events do not carry raw samples, output, profiles, patches, prompts, or runtime payloads. High-volume sample and tool detail is coalesced or omitted with an explicit omission event. Committed durable events remain authoritative for replay; SSE is only delivery.

### 9. Compatibility is additive and disabled projects stay untouched

The change adds operation enum values, optional DTO fields, new typed app capabilities, and new versioned routes. Existing operation and run routes retain their methods and response meanings. New JSON fields are optional where older runs cannot supply them.

A missing performance policy and `Mode: disabled` both preserve the existing post-execute flow. Status may report the disabled policy from existing project inputs, but it must not parse target rows beyond mismatch validation, discover benchmarks, initialize a runtime, run commands, create performance state, or change existing artifact bytes. No migration is required from users. Private performance schemas are versioned and reject or conservatively recover unsupported in-flight state rather than projecting it as current.

## Trade-Offs

### Typed performance methods plus the generic operation bridge

A performance-only operation framework would make the first implementation locally tidy, but it would duplicate confirmation, fingerprint, durable acceptance, cancellation, events, and web security. Using only the generic `OperationResult` would have the opposite problem: adapters would need to parse content or know private sprint records. The chosen split keeps typed product queries and results while reusing the established operation bridge for guarded execution.

### Separate run lifecycle and product outcome

Two status fields require presenters and tests to do more work. Combining them would be simpler, but it would falsely equate successful orchestration with a met target. The separation is necessary because `target_missed`, `stalled`, and other valid product results can be committed by an operationally successful run.

### Dedicated sprint queries and shared run resources

Dedicated performance status, result, and evidence routes make the product model easy to inspect. Reusing run resources for acceptance, cancellation, and events avoids a second job protocol. Putting all performance data on `/runs/{id}` was rejected because run control does not own target facts, qualification, freshness, or sprint verdicts.

### Strict stale rejection

Automatically refreshing a stale start request would reduce one user interaction. It could also execute against targets, correctness gates, worktree bytes, or limits that the user did not review. The API rejects the request and returns a fresh preparation action instead. This is deliberately conservative because a performance run may promote source changes.

### Bounded projections instead of raw evidence APIs

Raw evidence would help ad hoc debugging, but it would expose hostile command and runtime output, increase memory and transfer cost, and create another authority beside immutable sprint records. Bounded metadata plus allowlisted artifact preview is less convenient and much safer. Operators still receive paths, digests, omission counts, and correlation IDs needed to inspect local evidence deliberately.

### Serialized measurement with limited preparation concurrency

Parallel status queries, parsing, and independent preparation can remain bounded. Measured benchmark samples should run under the sprint-owned stability policy, normally serially, because generic concurrency would contaminate the environment being measured. This sacrifices throughput for comparable evidence. The API exposes progress and limits but does not expose a caller-selected worker count.

### Closed v1 filters and descriptors

The evidence query accepts a fixed filter set, and benchmark descriptors remain product-owned. A generic query language, parser registry exposed to callers, or runtime-authored command endpoint would be more extensible. It would also add injection, compatibility, and unbounded-work risks without a Sprint 39 need.

## Evidence

- `projects/ultraplan-go/sprints/39-performance-stage/requirements.md:95-104` requires adapter-independent lifecycle DTOs, performance operation kinds, one shared durable runner, and interface parity. AC-3 and AC-10 require runtime-free deterministic admission, durable acceptance, cancellation, recovery, and agreement across all projections.
- `projects/ultraplan-go/sprints/39-performance-stage/requirements.md:194-220` defines the exact target and run outcomes, separates product authority from durable operations, and requires one terminal result under stale writers and cancellation races.
- `projects/ultraplan-go/sprints/39-performance-stage/technical-handbook.md:31-53` favors thin lifecycle commands, explicit dependencies, typed outcomes, bounded execution, propagated cancellation, and one bounded fact model for every interface.
- `projects/ultraplan-go/sprints/39-performance-stage/technical-handbook.md:59-67` records the relevant costs: lazy preflight delays failures, streaming complicates atomic recovery, rich private evidence can diverge from public state, and general extension points can admit unsafe authority.
- `projects/ultraplan-go/sprints/39-performance-stage/technical-handbook.md:114-125` identifies the central API pressures: no alternate target source, meaningful runtime-free dry-run, frozen identities, conservative cancellation, bounded public state, measurement stability, and byte-compatible disabled behavior.
- `projects/ultraplan-go/docs/ARCHITECTURE.md:303-325` assigns shared typed use cases to `internal/app` and limits CLI, TUI, and web to request and presentation mapping. `projects/ultraplan-go/docs/ARCHITECTURE.md:370-388` assigns durable lifecycle, event order, fencing, cancellation, and reconciliation to `internal/runcontrol` while leaving product success to `internal/sprint`.
- `system/contracts/surfaces/api-contracts.md:52-70` requires explicit transport DTOs rather than raw internal records. Its compatibility, idempotency, pagination, error, and rate requirements support additive schemas, deduplicated acceptance, opaque bounded cursors, stable codes, and overload behavior.
- `system/contracts/core/errors.md:88-138` and `system/contracts/core/observability.md:72-96` require explicit failure state, stable structured errors, and correlation across requests, runs, tasks, runtimes, and commands.
- `system/contracts/runtime/workflows.md:67-151` requires inspectable outcomes, explicit retry and cancellation, protected replay, reconciliation of partial effects, and versioned long-running workflow definitions.
- `system/contracts/runtime/performance.md:83-178` requires bounded expensive work, observable ownership and cancellation, progress, and visible runtime cost drivers. It also rejects speculative optimization, which supports keeping caller-selected concurrency and generic tuning out of this API.
- `internal/app/operations.go:23-46`, `internal/app/operations.go:127-166`, and `internal/app/operations.go:506-659` already provide the shared operation, confirmation, stale-fingerprint, durable-manager, event, and result seams that performance should extend.
- `internal/app/sprint_usecases.go:78-194` demonstrates the accepted pattern of typed bounded QA and repair use cases alongside the generic operation path. `internal/app/sprint_usecases.go:247-308` demonstrates bounded public identity, artifact-reference, and limit summaries.
- `internal/web/routes.go:257-275` and `internal/web/routes.go:318-409` already define versioned operation preparation, operation execution, run inspection, cancellation, and event replay routes. Reuse avoids incompatible performance-only lifecycle behavior.
- `internal/runcontrol/model.go:16-52`, `internal/runcontrol/model.go:224-330`, and `internal/runcontrol/interfaces.go:47-64` provide the existing lifecycle, correlation, fencing, cancellation, terminal, pagination, and reconciliation contract. None of those types should absorb performance verdict authority.

## Risks

| Risk | Consequence | Required control |
| --- | --- | --- |
| Operation lifecycle is rendered as the performance verdict | A completed `target_missed` run appears to pass | Keep `lifecycle` and `outcome` separate in every fixture and presenter; test all cross-products that can occur. |
| Preparation performs benchmark or runtime work | Dry-run mutates state, incurs cost, or becomes non-deterministic | Keep app preparation limited to reads, validation, normalization, fingerprints, and bounded display facts; use fakes that fail tests on runtime or process calls. |
| A stale confirmation reaches durable acceptance | Work runs against unreviewed targets, source, or limits | Recompute immediately before acceptance and return `performance.request_stale` before any durable run or child side effect. |
| Resume creates new authority | A caller changes a frozen command, model, limit, or target after partial progress | Accept only the performance attempt ID; load all other values from persisted frozen state and reject drift. |
| Public results mirror private records | Schemas become coupled and hostile evidence leaks | Build explicit allowlisted DTO projections and test that raw samples, output, profiles, patches, prompts, environment, and provider payloads are absent. |
| Evidence pagination races retention or a new attempt | Clients silently skip or mix records | Bind opaque cursors to attempt and filters, use stable ordering, and return an explicit stale or retention-gap error. |
| Cancellation wins while a command or runtime completes late | A late writer overwrites cancellation or cleanup uncertainty with success | Route cancellation through run control, fence every product publication, and allow exactly one product terminal result. |
| Generic HTTP retries duplicate promotions | Repeated requests apply the same benchmark or implementation patch twice | Deduplicate start/resume before work and keep each promotion identity-checked and idempotent inside `internal/sprint`. |
| Status and result reads trigger broad scans | Frequent CLI, TUI, or browser polling degrades large workspaces | Read current bounded state and digest-bound pointers only; paginate evidence and never reconstruct status from the attempt tree. |
| Disabled projects change existing files during status | Backward compatibility is broken for projects that did not opt in | Keep disabled status read-only and cover byte-for-byte artifact preservation in compatibility tests. |
| Error mappings drift by interface | Automation cannot distinguish missed, blocked, cancelled, and cleanup-uncertain outcomes | Define one app error and outcome mapping, then use shared fixtures for CLI text, CLI JSON, TUI, browser HTML, browser JSON, durable runs, and artifacts. |

No API design question remains open for final sprint reasoning. Implementation may choose concrete Go type names, but it must preserve the operation semantics, ownership boundaries, enums, fingerprints, bounds, and compatibility rules above.
