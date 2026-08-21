# Sprint Reasoning: Guarded Web Operations and SSE Progress

> Project: `ultraplan-go`
> Sprint: `31-web-operations`
> Output: `projects/ultraplan-go/sprints/31-web-operations/reasoning.md`
> **Inputs Used:** `projects/ultraplan-go/project-index.md`, `projects/ultraplan-go/sprints/31-web-operations/requirements.md`, `projects/ultraplan-go/docs/ARCHITECTURE.md`, `projects/ultraplan-go/docs/PRD.md`, `projects/ultraplan-go/docs/TRD.md`, `projects/ultraplan-go/sprints/31-web-operations/sprint-index.md`, `projects/ultraplan-go/sprints/31-web-operations/technical-handbook.md`, `projects/ultraplan-go/sprints/31-web-operations/reasoning/api-design.md`, `projects/ultraplan-go/sprints/31-web-operations/reasoning/architecture.md`, `projects/ultraplan-go/sprints/31-web-operations/reasoning/frontend.md`, `studies/go-cli-study/reports/final/01-project-structure.md`, `studies/go-cli-study/reports/final/03-dependency-injection.md`, `studies/go-cli-study/reports/final/04-configuration-management.md`, `studies/go-cli-study/reports/final/05-error-handling.md`, `studies/go-cli-study/reports/final/06-io-abstraction.md`, `studies/go-cli-study/reports/final/07-state-context.md`, `studies/go-cli-study/reports/final/08-concurrency.md`, `studies/go-cli-study/reports/final/10-logging-observability.md`, `studies/go-cli-study/reports/final/11-testing-strategy.md`, `studies/go-cli-study/reports/final/12-extensibility.md`, `studies/go-cli-study/reports/final/13-security.md`, `studies/go-cli-study/reports/final/14-performance.md`

This document decides. It synthesizes the selected context, handbook evidence, area-specific reasoning, and contracts into final Sprint 31 decisions. It does not replace `sprint-index.md`, `technical-handbook.md`, or `reasoning/*.md`.

Requirement IDs used below are local traceability labels for the ordered lists in `requirements.md`: `AC-01` through `AC-18` are the Acceptance Criteria in source order, and `C-01` through `C-15` are the Constraints in source order. `RO-*` refers to the named Required Output path.

## Sprint Purpose

- **Goal:** Expose the existing validation, prompt-preview, dry-run, sprint-flow, execute, review, smoke, verify, and study run-loop capabilities through the loopback web surface with current server-issued confirmation, bounded SSE progress, explicit cancellation, product-owned mutation exclusion, and truthful graceful-shutdown recovery.
- **Non-Goals:** New workflows or workflow semantics; a browser command console; arbitrary file editing; durable web jobs or event history; detached work that survives server exit; hosted or remote access; accounts or permissions; WebSockets or terminal transport; a frontend framework/build pipeline; automatic fixes or Git mutation; and Sprint 32 release hardening.
- **Depends On:** Sprint 30's loopback server, security envelope, `/api/v1` conventions, embedded assets, and read-only app boundary; Sprints 23-29 app use cases, durable state, cancellation, TUI controls, review, smoke, and verify; existing study locking; and the Phase 4 ownership rules in the project PRD, TRD, and Architecture.

## Selected Context And Pre-Reasoning Artifacts

| Artifact | Path | How It Was Used |
| --- | --- | --- |
| Project Index | `projects/ultraplan-go/project-index.md` | Established the implementation repository, Phase 4 dependency direction, selected contract pool, selected study reports, and required review protocols. |
| Requirements | `projects/ultraplan-go/sprints/31-web-operations/requirements.md` | Supplied the authoritative scope, acceptance behavior, outputs, shutdown contract, security constraints, and verification gate. |
| Project Architecture | `projects/ultraplan-go/docs/ARCHITECTURE.md` | Fixed module ownership: `internal/web -> internal/app -> product/platform`, ephemeral browser projections, product-owned durable state, and explicit composition. |
| PRD | `projects/ultraplan-go/docs/PRD.md` | Fixed the local-only guarded browser scenario, one product core, durable-state recovery, and Phase 4 non-goals. |
| TRD | `projects/ultraplan-go/docs/TRD.md` | Fixed the route family, HTTP/SSE separation, graceful shutdown behavior, product locks, security posture, app boundary, and deterministic testing expectations. |
| Sprint Index | `projects/ultraplan-go/sprints/31-web-operations/sprint-index.md` | Limited reasoning to the selected contracts, 12 evidence reports, three area artifacts, and two review protocols; excluded new workflow, remote, editing, Git, WebSocket, and framework scope. |
| Technical Handbook | `projects/ultraplan-go/sprints/31-web-operations/technical-handbook.md` | Supplied evidence for thin adapters, constructor injection, context ownership, typed errors, bounded concurrency, pre-retention redaction, observability, and behavior-first tests. |
| Prior Decision | None selected | `project-index.md` catalogs no prior decision artifact. Sprint 30 behavior is carried as a dependency through requirements and project documents, not treated as an independently selected decision record. |

## Area-Specific Reasoning Inputs

| Area | Reasoning Document | Key Conclusion | Evidence Basis | Impact On Final Decision |
| --- | --- | --- | --- | --- |
| API Design | `projects/ultraplan-go/sprints/31-web-operations/reasoning/api-design.md` | Use prepare/start/status+result/events/cancel resources; bind a two-minute single-use confirmation to session, canonical request, and current fingerprint; use typed safe errors and bounded per-operation SSE. | Configuration resolution, typed-error, context, concurrency, observability, security, performance, and testing reports plus TRD section 18A. | Fixes route methods, lifecycle states, error codes, event names, retention limits, replay-gap behavior, and compatibility boundary. |
| Architecture | `projects/ultraplan-go/sprints/31-web-operations/reasoning/architecture.md` | Add a closed synchronous app operation capability, keep mutations and recovery product-owned, and make the web hub the sole bounded owner of asynchronous server lifecycle state. | Thin-entrypoint, dependency injection, context/lock, concurrency, extensibility, security, and performance evidence plus project architecture. | Fixes package dependencies, composition, lock domains, exact-once cancellation, shutdown order, and state authority. |
| Frontend | `projects/ultraplan-go/sprints/31-web-operations/reasoning/frontend.md` | Keep server-rendered detail pages authoritative for entry context and add narrow JavaScript for prepare, confirmation, EventSource progress, cancellation, and durable refresh. | Safe projection, test seams, lifecycle, bounded streaming, accessibility, and security evidence. | Fixes the browser state model, truthful reconnect/cancellation language, bounded DOM, native controls, and no client store. The final decision narrows one detail: because `AC-04` requires HTTP `DELETE`, no-JavaScript pages provide status/recovery and CLI cancellation guidance instead of a POST method-override cancellation route. |

## Sprint Technical Handbook Summary

- **Relevant Patterns:** Thin transport adapters; explicit constructor-based composition; complete effective-input normalization before validation; root-context cancellation; product-local locks; bounded pub/sub with independent consumers; typed errors with safe projections; stable correlation fields; pre-retention redaction; deterministic fakes and selective golden fixtures.
- **Important Trade-Offs:** Ephemeral speed versus restart continuity; exclusive per-sprint locking versus finer concurrency; disconnect-on-overflow versus gap-free blocking; typed taxonomies versus simpler sentinels; fixed conservative limits versus adaptive tuning; and a closed operation union versus dynamic extensibility.
- **Warnings / Anti-Patterns:** Do not infer operation state from connection state, block producers on subscribers, replace operation contexts with `context.Background()`, hide services in globals/context, classify errors by strings, retain raw diagnostics before redaction, assemble executable strings from browser input, or add speculative plugin/pooling machinery.
- **Evidence Confidence:** High for boundary, context, error, concurrency, observability, security, and testing principles because multiple studied repositories provide concrete examples. Medium for exact Sprint 31 policy choices: the numerical limits, token lifetime, lock domain, API resource layout, and recovery event are deliberate project decisions, not benchmark-derived or copied repository behavior.

## Contracts Applied

| Contract / Requirement ID | Constraint | Decision Impact | Expected Evidence |
| --- | --- | --- | --- |
| Architecture; `AC-05`, `AC-13`, `AC-14`; `C-01`-`C-04` | Web is transport-only; product state and workflows retain authority. | Closed `internal/app` operation capability, no direct web/product imports, no web scheduler or durable state. | Import review, shared-dispatch tests, hub restart/eviction tests. |
| Errors; `AC-11` | Failures must be typed, actionable, stable, and safely projected. | `errors.Is/As` mapping to compatibility-controlled API codes and typed terminal outcomes. | Error matrix tests and hostile diagnostic fixtures. |
| Configuration; `AC-01`, `AC-02`; `C-05` | Effective runtime/model/harness inputs must be fully resolved and redacted before confirmation. | App-owned normalization and fingerprinting; preparation is side-effect-free. | Precedence, summary, mismatch, expiry, and stale-fingerprint tests. |
| Observability; `AC-06`, `AC-12`; `C-07` | Progress and diagnostics must be correlated, bounded, safe, and truthful. | Per-operation sequence IDs, stable event names, request/operation IDs, safe counters, separate logs and SSE. | SSE schema/order tests and structured-log review. |
| Security; `AC-02`, `AC-11`, `AC-13`; `C-05`, `C-06`, `C-11`-`C-15` | Maintain loopback, same-origin, CSRF/session, containment, body/stream limits, redaction, and no Git/browser command authority. | Typed allowlisted operations, session-bound confirmation, strict route policy, pre-retention redaction. | Host/Origin/CSRF/body/path/redaction tests and security review. |
| Testing; `AC-18`; `C-13` | Normal proof must be deterministic and race-aware. | Fake app/runtime/harness/process dependencies, temporary workspaces, barriers/fake clocks, selective goldens. | Focused suites, `go test ./...`, `go test -race ./...`, and build. |
| Documentation; `RO-Web documentation`, `RO-CLI/API documentation updates`, `RO-Architecture documentation updates` | User and architecture behavior, trust limits, cancellation, shutdown, reconnect, and recovery must be documented. | Update the three required implementation-repository docs with the final contracts. | Documentation diff and Sprint Review checks. |
| LLM Runtime; LLM Evaluation / Cost / Safety; `AC-01`, `AC-03`, `AC-09`; `C-09` | Existing runtime ownership, cancellation, model metadata, safety, and cleanup must remain intact. | Web only requests app capabilities; canonical contexts reach agentwrap/runtime; raw events and provider data are not browser payloads. | Fake runtime cancellation/metadata tests and runtime-boundary review. |
| Workflows; `AC-03`, `AC-07`-`AC-10`; `C-02`, `C-03`, `C-08`-`C-10` | Existing workflows, cancellation, resumability, lock ownership, and terminal truth cannot be duplicated. | Synchronous app execution under server context, no queue, product locks, bounded cleanup and startup reconciliation. | Cross-surface, lock, cancel-all, cleanup-order, and restart tests. |
| Performance; `AC-05`, `AC-06`; `C-07` | All active-operation and streaming resources must be explicitly bounded; subscribers cannot backpressure work. | Fixed limits, bounded queues and retention, overflow disconnect, gap marker, no speculative pooling. | Capacity, rollover, slow-subscriber, aggregate-bound, race, and leak tests. |
| Persistence And Migrations; `AC-05`, `AC-07`-`AC-10`; `C-04`, `C-10` | Existing files/run state are authoritative; writes and recovery remain product-owned and truthful. | No web schema; stale running work becomes interrupted/recovery-required; locks release only after reconciliation. | Crash-recovery fixtures, atomic-state checks, and no-web-store review. |

## Repos Studied / Source Evidence Used

| Source / Repo / Report | Concrete Reference | Relevant Finding | Why It Matters For This Sprint | Used In Decision(s) |
| --- | --- | --- | --- | --- |
| Helm / `01-project-structure` | `helm/pkg/cmd/install.go:347`; `helm/pkg/action/install.go:73-140` | Command adapters delegate to action types. | Supports a thin HTTP adapter over app operations. | 1, 2 |
| gdu / `01-project-structure` | `gdu/cmd/gdu/app/app.go:30-49` | Output interfaces can vary without moving product work into the surface. | Supports CLI/TUI/web projection over one core. | 1, 8 |
| gh-cli and opencode / `03-dependency-injection` | `gh-cli/pkg/cmdutil/factory.go:16-43`; `opencode/internal/app/app.go:42-81` | Central manual composition and narrow service boundaries improve traceability and testing. | Supports composition in `cmd/ultraplan` and injected clocks, IDs, limits, and fakes. | 1, 8 |
| K9s and restic / `04-configuration-management` | `k9s/internal/config/k9s.go:423-451`; `restic/internal/global/global.go:139,147` | Validate after complete resolution and distinguish explicit overrides. | Supports canonical preparation and stale-input checks. | 2 |
| restic, Go Task, and K9s / `05-error-handling` | `restic/internal/restic/lock.go:47`; `go-task/errors/errors_task.go:13-32`; `k9s/internal/model/flash.go:100-103` | Typed errors carry remediation while user and operational projections remain separate. | Supports lock conflicts, stable API codes, and safe browser messages. | 3, 6 |
| gh-cli, restic, Go Task, and chezmoi / `06-io-abstraction` | `gh-cli/pkg/iostreams/iostreams.go:551-568`; `restic/internal/ui/mock.go:10-53`; `go-task/executor_test.go:146-151`; `chezmoi/internal/cmd/applycmd_test.go:220-241` | Test constructors, synchronized capture, and HTTP test servers isolate volatile boundaries. | Supports deterministic app, HTTP, SSE, and concurrent-output tests. | 8 |
| Helm, Go Task, and restic / `07-state-context` | `helm/pkg/cmd/install.go:333-347`; `go-task/task.go:89`; `restic/internal/restic/lock.go:105,290-305` | Root contexts propagate cancellation; locks and bounded post-cancel cleanup need explicit ownership. | Supports operation/subscriber separation, lock lifetime, and shutdown cleanup. | 3, 4 |
| K9s and opencode / `08-concurrency` | `k9s/internal/pool.go:21,30,37`; `opencode/internal/pubsub/broker.go:67-82`; `opencode/cmd/root.go:261-279` | Bound fan-out, clean subscriptions, and time-bound shutdown. | Supports hub capacity, exact ownership, overflow handling, and shutdown waits. | 4, 5 |
| K9s, opencode, and restic / `10-logging-observability` | `k9s/internal/slogs/keys.go:6-231`; `opencode/internal/logging/logger.go:25-62`; `restic/internal/backend/logger/log.go:22-77` | Stable structured fields and user-output/log separation improve diagnosis. | Supports correlated but distinct logs, events, and API results. | 5, 6 |
| gh-cli, Helm, restic, and K9s / `11-testing-strategy` | `gh-cli/pkg/httpmock/stub.go:35-199`; `helm/internal/test/test.go:43`; `restic/internal/backend/mock/backend.go:14-26`; `k9s/internal/view/pod_test.go:23` | Central fakes and selective goldens support behavior tests without brittle private-shape assertions. | Defines the verification strategy. | 8 |
| dive and gh-cli / `12-extensibility` | `dive/cmd/dive/cli/internal/command/adapter/analyzer.go:13-15`; `gh-cli/pkg/cmdutil/factory.go:16-43` | Adapters/factories permit extension without dynamic plugins or global registries. | Supports a closed operation union and explicit dispatcher. | 1 |
| opencode, restic, Helm, and lazygit / `13-security` | `opencode/internal/permission/permission.go:44-108`; `restic/internal/options/secret_string.go:15-20`; `helm/pkg/registry/transport.go:37-41`; `lazygit/cmd_obj_builder.go:38` | Explicit confirmation/trust boundaries, redaction types, credential scrubbing, and argv execution reduce authority and disclosure risk. | Supports typed capabilities, separate CSRF/confirmation, and pre-retention redaction. | 2, 6 |
| opencode, K9s, age, yq, and gh-cli / `14-performance` | `opencode/internal/llm/provider/provider.go:56`; `k9s/internal/pool.go:26-48`; `age/internal/stream/stream.go:20,195-219`; `yq/pkg/yqlib/stream_evaluator.go:78-113`; `gh-cli/pkg/cmdutil/factory.go:27-42` | Bounded streaming/concurrency, incremental processing, and lazy initialization preserve responsiveness; pooling requires profiling. | Supports explicit hub/SSE/DOM bounds and no speculative optimization. | 2, 5, 7 |

The reports support principles rather than the exact web protocol. The two-minute confirmation lifetime, fixed hub limits, one status/result resource, `recovery_required`, no-queue rule, per-operation event IDs, and exclusive sprint lock domain are Sprint 31 decisions grounded in project requirements and area reasoning, not empirical claims about the studied repositories.

## Trade-Off And Debt Analysis

### Accepted Trade-Offs

| Trade-Off | Benefit | Cost / Constraint Accepted | Why Acceptable Now | Revisit Trigger |
| --- | --- | --- | --- | --- |
| Closed typed operation union | Auditable authority, validation, confirmation, and safe projection. | Every newly exposed existing use case needs an explicit variant and mapper. | Sprint 31 has a finite operation set and no plugin requirement. | A governed use case cannot be added without repeated dispatcher boilerplate or a separately approved extension model. |
| Synchronous app operation under asynchronous web ownership | Keeps app operations transport-neutral and reusable by CLI/TUI. | The hub owns goroutine lifecycle and terminal arbitration. | It avoids pushing a second job abstraction into the product core. | A durable worker architecture is explicitly selected. |
| Exclusive per-sprint mutation lease | Prevents overlapping artifact, state, review, and smoke mutations. | Reduces concurrency among operations that may eventually prove independent. | Safety and recovery are more important than speculative local parallelism. | A documented read/write matrix proves safe finer-grained domains and includes deadlock/recovery analysis. |
| Ephemeral hub and short retention | Avoids a shadow authority, schema, migration, and dual-write path. | Restart loses operation IDs, event history, and terminal projections. | Durable product state already supports recovery and remains canonical. | Work must intentionally survive server restart or historical event audit becomes a product requirement. |
| No web queue | Every `202` has immediate server ownership and unambiguous shutdown behavior. | Bursts are rejected rather than smoothed. | Local operation volume is small and truthful capacity failure is preferable to scheduler semantics. | Measured local workloads require queuing and a governed scheduler/ownership design exists. |
| Disconnect slow SSE subscribers | Product work never blocks and memory stays bounded. | A browser can miss transient progress and must reconnect/refresh. | Progress is non-authoritative and gap signaling is explicit. | A future requirement needs guaranteed event delivery and selects durable replay. |
| Fixed conservative limits | Makes security, performance, and tests concrete. | Values may be under- or over-provisioned for real workloads. | No measurements justify adaptive behavior yet. | Sprint 32/browser dogfood supplies event-size, operation-mix, or memory data. |
| Minimal JavaScript with strict `DELETE` cancellation | Preserves the API method contract and avoids a second request model. | Without JavaScript, browser cancellation is unavailable; the page must direct users to CLI cancellation/recovery. | JavaScript is embedded and required only for this enhancement; status and durable recovery remain usable without it. | A later requirement mandates fully no-JavaScript mutation parity and defines a compatible HTTP contract. |

### Potential Technical Debt

| Debt / Shortcut | Why It Might Accrue | Current Mitigation | Owner / Follow-Up |
| --- | --- | --- | --- |
| App operation service becomes broad | One service fronts heterogeneous product capabilities. | Closed variants, consumer-sized methods, private dependencies, and no generic maps/registry. | Reassess after Sprint 33 adds another capability through the same boundary. |
| Fingerprint omissions | A variant may omit an effective option, governed input, runtime/model, or harness identity. | Canonicalization contract per variant and stale-input tests that mutate each material input. | Sprint 31 implementation and Architecture Review. |
| Conservative limits lack production measurements | Aggregate event/result/subscriber memory may differ from assumptions. | Enforce all limits, add local counters and aggregate-bound tests, document values. | Sprint 32 hardening uses observed measurements. |
| Coarse mutation lease limits throughput | Independent operations may serialize unnecessarily. | Per-sprint rather than global locking preserves unrelated sprint concurrency. | Future lock-domain decision only after an operation read/write matrix. |
| Operation-specific API/result DTO growth | New fields can drift from redaction and compatibility rules. | Explicit DTOs, compatibility fixtures, unknown-field rejection, and redaction review for additions. | Each future web operation change. |
| No durable idempotency after a lost `202` | Client cannot prove whether creation succeeded using the token alone. | Token replay cannot start duplicate work; UI refreshes visible/durable state before re-preparing. | Revisit only with durable workers or demonstrated duplicate-risk incidents. |
| Recovery differs among existing workflows | Some product state may not yet express interruption or cleanup uncertainty uniformly. | Add only owner-specific truthful states required by Sprint 31, with compatibility tests; never add web-owned recovery state. | Product module owners during implementation; Persistence contract review. |

### Future Considerations

| Consideration | Deferred Until | Reason Deferred | What Should Be Preserved Now |
| --- | --- | --- | --- |
| Durable workers surviving server exit | Explicit future architecture/version | Requires leases, heartbeats, worker identity, durable idempotency, ownership transfer, authentication, and replay. | App capability boundary, durable product truth, and versioned API semantics. |
| Finer mutation lock classes | Proven concurrency need | Current overlap matrix is not sufficiently proven. | Typed mutation class and owner metadata without exposing lock internals. |
| Adaptive capacity and buffer tuning | Sprint 32 measurements | Current values are engineering defaults. | Injected validated limits and local counters. |
| Browser/API compatibility hardening and accessibility audit | Sprint 32 | Release-level cross-browser and accessibility evidence is explicitly deferred. | Native controls, semantic status, reduced motion, stable `/api/v1` names, and no framework coupling. |
| Frontend framework | Demonstrated interaction complexity | Current server-rendered pages and narrow enhancement do not justify a client application. | Versioned HTTP boundary and explicit view models. |
| Durable event history or metrics endpoint | Explicit observability requirement | Would add persistence/security/retention contracts beyond progress needs. | Stable safe event schema, correlation fields, and local counters. |

## Decisions

The eight decisions below are the binding Sprint 31 decisions. Their names, contracts, alternatives, and evidence requirements are normative for `plan.md`.

## Final Decisions

### Decision 1: One Closed App Operation Capability And Thin Web Adapter

- **Decision:** Add a separate typed operation service in `internal/app` with a closed `OperationKind` union for validation, prompt preview, dry run, sprint flow, execute, review, smoke, verify, and study run-loop. It owns normalization, preparation, dispatch, safe progress/result projection, and typed errors. Execution is synchronous relative to its supplied context. `internal/web` maps HTTP/session/confirmation and runs accepted calls asynchronously in its hub; it imports no product, runtime/process, or CLI-handler packages. `cmd/ultraplan` explicitly composes dependencies. Existing CLI and TUI operation paths must call the same app use cases.
- **Rationale:** This is the smallest boundary that provides web parity without making HTTP a second application layer or spreading web job semantics into product code.
- **Study / Source Grounding:** `technical-handbook.md` thin-adapter and constructor-injection patterns; `01-project-structure` at Helm `pkg/cmd/install.go:347` and gdu `cmd/gdu/app/app.go:30-49`; `03-dependency-injection` at gh-cli `pkg/cmdutil/factory.go:16-43` and opencode `internal/app/app.go:42-81`; `12-extensibility` adapter/factory evidence.
- **Trade-Offs Accepted:** Explicit variants and wiring add code, but make authority and tests auditable.
- **Technical Debt / Future Impact:** Guard against a god service by keeping operation-specific product behavior private to existing services. Future operations extend the union explicitly; no plugin framework is implied.
- **Alternatives Rejected:** Direct web-to-product imports violate dependency rules; calling/parsing CLI handlers couples transport to flags/output; an app-level asynchronous job manager duplicates the hub/scheduler; a generic registry, reflection dispatcher, or browser-authored command broadens authority without a requirement.
- **Contracts Satisfied:** Architecture, Workflows, LLM Runtime, LLM Evaluation / Cost / Safety, Testing; `AC-03`, `AC-13`, `AC-14`, `C-01`-`C-03`, `RO-App operation use cases`, `RO-App operation tests`.
- **Evidence Required:** Compile/import review; shared fake-workspace fixtures showing CLI/TUI/web agreement; tests proving preparation and dispatch invoke the same app paths and web cannot access product/runtime types.

### Decision 2: Binding Prepare/Start Confirmation And Versioned Operation API

- **Decision:** Implement `POST /api/v1/operations/prepare`, `POST /api/v1/operations`, `GET /api/v1/operations/{id}`, `GET /api/v1/operations/{id}/events`, and `DELETE /api/v1/operations/{id}`. Preparation is side-effect-free and returns normalized scope, workspace-relative affected paths, mutation class, effective runtime/model or harness summary, prerequisites, current SHA-256 governed-input fingerprint, two-minute expiry, and opaque token. The bounded server record binds token, per-process secret, session, canonical normalized request, fingerprint, and expiry. Start repeats and renormalizes the request, rechecks current inputs, and consumes the token once creation succeeds; mismatch, replay, expiry, and staleness are distinct `409` codes. No accepted-work queue exists.
- **Rationale:** Confirmation is meaningful only when it approves the complete current operation that will execute. Immediate ownership keeps acceptance, capacity, cancellation, and shutdown semantics unambiguous.
- **Study / Source Grounding:** `technical-handbook.md` post-resolution validation and explicit fingerprint pressure; `04-configuration-management` at K9s `internal/config/k9s.go:423-451` and restic `internal/global/global.go:139,147`; `13-security` explicit confirmation/trust-boundary evidence. Exact token semantics and lifetime are Sprint-specific resolutions, not directly established by a studied repo.
- **Trade-Offs Accepted:** Server-side preparation retention and repeated normalization add state/work; no general idempotency key leaves a response-loss recovery step.
- **Technical Debt / Future Impact:** Every variant needs a complete canonicalization contract. A future durable worker API must version token, idempotency, and ownership semantics rather than relaxing this contract.
- **Alternatives Rejected:** Direct start makes confirmation ceremonial; client-signed scope/fingerprint trusts browser claims; stateless tokens cannot enforce replay without state; reserving locks/capacity during preparation causes side effects; a web queue becomes scheduler state.
- **Contracts Satisfied:** Architecture, Configuration, Security, Errors, Performance, Workflows; `AC-01`-`AC-04`, `AC-11`, `C-05`-`C-07`, `C-11`, `RO-Web operation routes`, `RO-Web confirmation and CSRF/session support`.
- **Evidence Required:** Strict JSON and unknown-field tests; no-side-effect preparation tests; per-variant canonicalization/fingerprint fixtures; expiry/mismatch/replay/stale/session/capacity/draining tests; route/method/`Allow`/`Location`/unknown-API tests.

### Decision 3: Product-Owned Mutation Exclusion And Conservative Recovery

- **Decision:** `internal/sprint` owns one exclusive per-sprint mutation lease shared by sprint flow, execute, review, smoke, and verify. Existing `internal/study` locking remains authoritative for study run-loop. Read-only validation, prompt preview, dry run, status, and result inspection take no new mutation lease. Acquisition occurs after start revalidation and before mutation or external execution. Release occurs only after durable terminal reconciliation or truthful `cleanup_uncertain` recording. Startup reconciles stale locks with durable running state and owned-process evidence; ambiguity becomes interrupted/recovery-required, never success.
- **Rationale:** Locking must protect product state across CLI, TUI, and web and remain held while cleanup can still mutate state or own child processes.
- **Study / Source Grounding:** `technical-handbook.md` explicit lock ownership; `05-error-handling` restic typed lock conflict `internal/restic/lock.go:47`; `07-state-context` restic lock/cleanup `internal/restic/lock.go:105,290-305` and Go Task in-flight coordination `task.go:438-469`.
- **Trade-Offs Accepted:** Coarse per-sprint exclusion sacrifices possible concurrency to avoid overlap and deadlock ambiguity.
- **Technical Debt / Future Impact:** Finer domains require a separately reviewed operation read/write matrix. Product schemas may need explicit interruption/uncertainty representation under their normal compatibility rules.
- **Alternatives Rejected:** HTTP middleware locks are bypassed by CLI/TUI; one global mutex blocks unrelated projects/studies; optimistic concurrent mutation risks artifact/state corruption; automatic stale unlock based on missing PID/process can release ownership prematurely.
- **Contracts Satisfied:** Architecture, Errors, Workflows, Persistence And Migrations; `AC-07`-`AC-10`, `C-06`, `C-10`, `RO-Sprint mutation lock support`, `RO-Sprint lock tests`.
- **Evidence Required:** Cross-surface conflict tests; acquire/release and actionable typed conflict tests; cleanup-barrier assertions proving no early release; stale-lock/restart fixtures; non-overlap tests with study locking.

### Decision 4: Server-Owned Lifecycle, Exact-Once Cancellation, And Ordered Shutdown

- **Decision:** The web hub owns one root context, canonical cancellation function, owner goroutine, ephemeral state, safe events/result, subscribers, and exact-once terminal arbitration per accepted operation. Explicit `DELETE`, timeout, and graceful shutdown cancel the same root lineage with distinct reasons; browser/SSE disconnect cancels only subscription context. Graceful shutdown enters draining before HTTP closure, rejects preparation/start, snapshots and cancels all active owners once with `server_shutdown`, keeps status/SSE available during bounded canonical cleanup, escalates owned process-tree cleanup at deadline, records truthful terminal/interrupted/cleanup-uncertain product state, releases locks only through reconciliation, publishes terminal events, then closes streams and HTTP. Cleanup uses a bounded context derived from the shutdown deadline, never an unbounded background context.
- **Rationale:** Every server-started operation must have one auditable owner and one cleanup path. Cancellation requests are not terminal outcomes; product reconciliation determines terminal truth.
- **Study / Source Grounding:** `07-state-context` root cancellation in Helm `pkg/cmd/install.go:333-347`, sibling cancellation in Go Task `task.go:89`, and bounded cleanup in restic `internal/restic/lock.go:290-305`; `08-concurrency` explicit shutdown timeout in opencode `cmd/root.go:261-279` and subscriber cleanup in `internal/pubsub/broker.go:67-82`.
- **Trade-Offs Accepted:** Shutdown may wait until the configured deadline and expose uncertainty instead of exiting immediately or claiming clean cancellation.
- **Technical Debt / Future Impact:** Detached/durable workers remain impossible under this model by design. Any future survival across restart requires a new ownership architecture.
- **Alternatives Rejected:** Cancelling on browser disconnect conflates observation and work; fire-and-forget goroutines violate ownership; fail-fast global `errgroup` can obscure sibling cleanup; immediate lock release permits overlap; `context.Background()` can detach cleanup; closing SSE first hides truthful terminal outcomes.
- **Contracts Satisfied:** Architecture, Workflows, LLM Runtime, Observability, Persistence And Migrations; `AC-05`, `AC-07`-`AC-10`, `C-08`-`C-10`, `RO-Web operation hub`, `RO-Web operation tests`.
- **Evidence Required:** Multi-operation cancel-all tests; exact-once cancellation and terminal-race tests; nested fake runtime/harness/process cancellation; draining rejection; deadline escalation; lock-release ordering; terminal-before-stream-close; stale-state restart reconciliation; race and leak checks.

### Decision 5: Bounded Ephemeral Hub And Progress-Only SSE

- **Decision:** SSE is observation only. Events use per-operation decimal IDs assigned at one hub serialization point and stable names `snapshot`, `progress`, `warning`, `finding`, `artifact`, `cancel_requested`, `recovery_required`, and `terminal`. Initial connection emits a snapshot, replays retained events newer than `Last-Event-ID`, then follows live events. A replay gap emits `recovery_required` plus current snapshot and durable refresh path. Heartbeats are comments every 15 seconds. Slow subscribers are disconnected on a full queue; producers never block. Defaults/upper bounds are: 64 KiB request body; 8 active operations; 128 preparations for at most 2 minutes; 256 events and 256 KiB per operation; 16 KiB encoded event payload; 256 KiB terminal projection; 8 subscribers per operation; 32 server-wide streams; 32 events per subscriber queue; 10-minute terminal retention; 15-second heartbeat; 30-minute connection lifetime. Capacity is rejected with `429`; draining with `503`.
- **Rationale:** Bounded recent progress is useful while durable product state remains authoritative. Explicit gaps are more truthful than silent loss or unbounded replay.
- **Study / Source Grounding:** `technical-handbook.md` bounded pub/sub and retention; `08-concurrency` K9s bounds `internal/pool.go:21,30,37` and opencode subscriber cleanup `internal/pubsub/broker.go:67-82`; `14-performance` opencode slow-consumer handling `internal/llm/provider/provider.go:56`, K9s bounded pool `internal/pool.go:26-48`, and incremental processing examples. Numerical values are conservative Sprint defaults, not study-derived benchmarks.
- **Trade-Offs Accepted:** Subscribers can miss progress, records expire, and conservative limits may reject local bursts.
- **Technical Debt / Future Impact:** Measure aggregate allocations and event sizes in Sprint 32 before tuning. No durable event log, adaptive limit system, or buffer pool is introduced.
- **Alternatives Rejected:** Blocking subscribers can deadlock product work; silent drops misrepresent continuity; unlimited buffers/replay violate bounds; server-global IDs add contention and imply cross-operation replay; producer-native IDs may be concurrent/unsafe; WebSockets add unneeded bidirectional authority.
- **Contracts Satisfied:** Performance, Observability, Architecture, Security; `AC-04`-`AC-06`, `C-07`-`C-09`, `RO-SSE tests`.
- **Evidence Required:** Concurrent monotonic-ID tests; snapshot/replay/rollover/gap/heartbeat/terminal tests; payload and aggregate limit tests; slow/abandoned subscriber tests proving operation completion; disconnect isolation; retention/eviction/restart tests; race/leak runs.

### Decision 6: Stable Safe Error, Result, Security, And Observability Projections

- **Decision:** Pre-acceptance failures use the Sprint 30 envelope with stable codes for invalid request; CSRF/origin/session; not found; confirmation expiry/mismatch/replay/staleness; mutation conflict; validation; prerequisite; capacity; draining; and internal failure. App/product failures use typed errors inspected via `errors.Is/As`; once accepted, success, failure, timeout, cancellation, interruption, and cleanup uncertainty are returned as operation outcomes in a `200` operation document. Safe allowlisted DTOs are created before entering event/result retention. Raw provider events, prompts, stderr, environment values, credentials, tokens, stack traces, executable strings, absolute paths, and lock internals have no browser projection. Request and operation IDs correlate logs/events; logs remain distinct and no network metrics endpoint is added. All mutation routes retain loopback, Host/Origin, per-process session, CSRF, body-limit, timeout, security-header, containment, and same-origin policy.
- **Rationale:** Stable client branching and useful diagnostics require structured categories, while retention makes pre-publication redaction mandatory.
- **Study / Source Grounding:** `05-error-handling` typed errors and user/operational split at restic `internal/restic/lock.go:47` and K9s `internal/model/flash.go:100-103`; `10-logging-observability` stable keys in K9s `internal/slogs/keys.go:6-231`; `13-security` redaction at restic `internal/options/secret_string.go:15-20`, credential scrubbing at Helm `pkg/registry/transport.go:37-41`, and explicit command construction at lazygit `cmd_obj_builder.go:38`. Detailed Host/Origin/CSRF/session policy is grounded primarily in project requirements/TRD rather than those study reports.
- **Trade-Offs Accepted:** Safe DTOs duplicate a controlled subset of app data and require compatibility maintenance.
- **Technical Debt / Future Impact:** New fields, errors, and event kinds require explicit DTO, redaction, and compatibility review; raw runtime kinds never become API kinds automatically.
- **Alternatives Rejected:** String-matched errors are brittle; serializing app/runtime structs couples layers and leaks fields; encoding-time-only redaction leaves unsafe retained data; returning `5xx` for an accepted failed operation conflates transport and product outcomes; permissive CORS/account auth are respectively unsafe and out of scope.
- **Contracts Satisfied:** Errors, Security, Observability, Configuration, LLM Evaluation / Cost / Safety; `AC-02`, `AC-06`, `AC-11`, `AC-13`, `C-05`, `C-07`, `C-11`, `RO-Web confirmation and CSRF/session support`.
- **Evidence Required:** Full status/code/error matrix tests; hostile secrets/path/provider/stderr fixtures at app-event, hub, JSON, SSE, HTML, and logs; request/operation correlation checks; Host/Origin/session/CSRF/body/header/path tests; API compatibility fixtures.

### Decision 7: Server-Rendered Operation Views With Narrow Progressive Enhancement

- **Decision:** Add allowlisted operation entry points to owning project, sprint, and study pages. Server preparation renders current normalized scope, affected relative paths, mutation class, runtime/model or harness, prerequisites, and expiry before explicit confirmation. Accepted operations render a server snapshot and durable refresh link; dependency-free JavaScript enhances the same DTOs with start, `EventSource`, `DELETE` cancellation, reconnect, duplicate suppression, and bounded in-place updates. The browser stores only transient pending/connection/last-rendered-ID state, never operation truth or tokens in persistent client storage. Views distinguish running, cancelling, succeeded, failed, cancelled, interrupted, cleanup-uncertain, server shutdown, reconnect, incomplete history, and connection loss. The activity DOM is bounded and all server data is inserted as escaped text/allowlisted URLs. Without JavaScript, preparation/start/status and durable recovery remain server-rendered; cancellation guidance uses the CLI rather than violating the API's required `DELETE` method with a POST override.
- **Rationale:** The browser remains useful, truthful, secure, and accessible without becoming a client-side workflow engine. The strict method adjustment resolves the only tension in frontend area reasoning in favor of `AC-04`.
- **Study / Source Grounding:** `technical-handbook.md` safe projection and bounded streaming; `05-error-handling` user-safe presentation; `06-io-abstraction` deterministic UI seams; `07-state-context` lifecycle separation; `13-security` confirmation/redaction; `14-performance` bounded incremental rendering. Exact interaction design is principally grounded in requirements, PRD/TRD, and `reasoning/frontend.md`.
- **Trade-Offs Accepted:** JavaScript is required for in-browser cancellation and live progress; the no-JavaScript path is status/recovery capable but not cancellation-parity complete.
- **Technical Debt / Future Impact:** Sprint 32 owns cross-browser and release-level accessibility hardening. Keep native controls, namespaced template/view-model boundaries, and narrow JS modules so that work does not require architecture changes.
- **Alternatives Rejected:** A SPA/client store duplicates authority; a generic command console obscures scope; client-generated confirmation can diverge; localStorage/IndexedDB/service workers create stale state; unbounded console output harms safety/performance; cancellation on unload is unreliable; POST method override contradicts the selected API contract.
- **Contracts Satisfied:** Architecture, Security, Observability, Performance, Documentation; `AC-01`, `AC-04`-`AC-07`, `AC-12`, `C-07`-`C-09`, `C-11`, `C-12`, `RO-Browser operation templates`, `RO-Browser operation JavaScript`.
- **Evidence Required:** Server-rendered prepare/start/status/recovery tests; enhanced DTO parity tests; all lifecycle/reason/error view fixtures; reconnect/gap/duplicate/one-stream/bounded-DOM tests; native control/focus/live-region/reduced-motion/color-independent checks; hostile HTML/URL tests; explicit no-JavaScript cancellation guidance check.

### Decision 8: Deterministic Cross-Surface Verification, Required Reviews, And Documentation

- **Decision:** Build normal evidence with `httptest`, fake app operations, fake clocks/IDs, fake runtimes/harnesses/process cleanup, deterministic barriers, and temporary workspaces. Prefer semantic lifecycle assertions; use representative goldens for complete JSON/SSE/template projections. Add focused app-operation, sprint-lock, web-hub, route, SSE, security, template, graceful-shutdown, restart-reconciliation, race/leak, and CLI/TUI/web agreement tests. Complete `go test ./...`, `go test -race ./...`, and `go build ./cmd/ultraplan` from `../ultraplan-go`. Update `docs/local-web.md`, `docs/cli-reference.md`, and `docs/architecture.md`. Final review must apply both `system/protocols/architecture-review-protocol.md` and `system/protocols/review-sprint-protocol.md`.
- **Rationale:** The main risks are boundary/lifecycle races and projection drift, all of which can be proven deterministically without live providers. Full race/build gates and selected reviews provide integration confidence.
- **Study / Source Grounding:** `06-io-abstraction` test constructors and synchronized capture; `11-testing-strategy` gh-cli HTTP fakes `pkg/httpmock/stub.go:35-199`, restic backend fakes `internal/backend/mock/backend.go:14-26`, Helm selective goldens `internal/test/test.go:43`, and warning against brittle K9s private-count assertions `internal/view/pod_test.go:23`.
- **Trade-Offs Accepted:** Extensive deterministic test infrastructure increases Sprint 31 implementation volume; gated real-runtime/browser release evidence remains deferred.
- **Technical Debt / Future Impact:** Sprint 32 must add the explicitly deferred compatibility, browser, accessibility-audit, packaging, and gated real-boundary evidence without weakening these deterministic gates.
- **Alternatives Rejected:** Live provider/harness dependencies in normal tests are slow and flaky; all-golden tests obscure behavior; private map/goroutine-count assertions overcouple implementation; unit tests alone miss cross-package shutdown and state-recovery behavior.
- **Contracts Satisfied:** Testing, Documentation, Architecture, Security, Performance; `AC-18`, `C-13`-`C-15`, all named `RO-*` test and documentation outputs.
- **Evidence Required:** Focused package test outputs; full test/race/build logs; import-boundary review; required protocol results; documentation review against every public state, limit, trust boundary, cancellation reason, and recovery path.

## Expected Evidence

| Evidence Type | Required Evidence | Source / Command / Review Check |
| --- | --- | --- |
| App tests | Variant normalization, fingerprints, side-effect-free preparation, shared dispatch, safe events/results, cross-surface agreement, nested cancellation. | `internal/app/web_operations_test.go`; focused `go test` package command chosen by implementation layout. |
| Lock tests | Per-sprint exclusion, typed conflicts, cleanup-order release, stale reconciliation, study-lock coexistence. | `internal/sprint/locks_test.go`. |
| Hub/shutdown tests | Capacity, retention, exact-once cancellation, terminal arbitration, draining, cancel-all, deadline escalation, restart recovery, no durable web state. | `internal/web/operations_test.go` plus server lifecycle tests. |
| API/security tests | Route methods, strict DTOs, confirmation states, sessions, Host/Origin/CSRF, body limits, safe envelopes, redaction, compatibility with Sprint 30. | `internal/web/routes_test.go`, security-focused tests, `httptest`. |
| SSE tests | Stable events, monotonic IDs, heartbeat, replay, gap recovery, rollover, bounds, slow subscribers, disconnect isolation, shutdown terminal ordering. | `internal/web/sse_test.go`. |
| UI tests | Confirmation content, lifecycle/result/finding/recovery views, enhanced/non-enhanced parity, escaping, accessibility semantics, bounded DOM/one stream. | Template/static tests and representative golden fixtures. |
| Runtime evidence | Safe structured logs with request/operation/state/cancel/cleanup correlation; local counters; no token/raw diagnostic leakage. | Captured fake-operation logs and manual review of structured fields. |
| Full verification | All packages pass tests and race detector; CLI binary builds. | From `../ultraplan-go`: `go test ./...`, `go test -race ./...`, `go build ./cmd/ultraplan`. |
| Architecture review | Web remains transport-only; app capability is shared; locks/recovery stay product-owned; hub is bounded/ephemeral; no alternate workflow engine or package cycle. | `system/protocols/architecture-review-protocol.md`. |
| Sprint review | Every selected contract, handbook decision, requirement, plan task, test result, documentation change, and deferred item is checked with a current verdict. | `system/protocols/review-sprint-protocol.md`. |
| Documentation | Guarded operation/API behavior, numerical limits, cancellation/disconnect/shutdown semantics, local trust boundary, and durable recovery are current. | `../ultraplan-go/docs/local-web.md`, `../ultraplan-go/docs/cli-reference.md`, `../ultraplan-go/docs/architecture.md`. |

## Assumptions And Risks

| Item | Type | Impact | Mitigation / Follow-Up |
| --- | --- | --- | --- |
| Existing app/product use cases accept propagated contexts and expose sufficient progress/results. | Assumption | Missing seams could tempt duplicated web logic. | Refactor command glue into `internal/app` first; preserve product semantics and add cross-surface fixtures. |
| Existing durable states can represent cancellation/interruption or can be compatibly extended by their owner. | Assumption | Recovery could otherwise overstate success or leave locks ambiguous. | Add only product-owned explicit states needed for truth; validate old state and atomic writes. |
| Preparation can inspect configuration/prerequisites without mutation. | Assumption | Hidden side effects would make confirmation unsafe. | Define side-effect-free inspection dependencies and assert no lock/runtime/harness/write calls; defer unsafe checks to start. |
| Fingerprint misses a material input. | Risk | User confirmation may authorize changed work. | Per-variant canonical input inventory and mutation tests for every governed/config/runtime/harness dependency. |
| Completion, timeout, user cancel, and shutdown race. | Risk | Conflicting terminal states or duplicate cleanup. | One terminal arbitration point, exact-once cancellation, durable-state comparison, race tests with barriers. |
| Cleanup releases a lock too early or blocks shutdown. | Risk | Concurrent mutation or hung server exit. | Release only after reconciliation; wait outside hub/product locks; deadline escalation tests. |
| Unsafe payload enters memory before redaction. | Risk | Secrets or paths leak through replay/results even if encoding is safe. | App-safe DTO boundary before hub publication and hostile fixtures at every ingress/egress. |
| Fixed limits multiply into excessive aggregate memory/goroutines. | Risk | Local denial of service or unstable tests. | Enforce aggregate limits, test worst-case allocations/owners, expose counters, tune only from measurements. |
| Browser treats a gap, eviction, or disconnect as operation failure/success. | Risk | Misleading recovery actions or unsafe rerun. | Distinct connection/history states, durable refresh guidance, no optimistic terminal UI, interaction tests. |
| Forced termination leaves a child process or stale lock. | Risk | A new mutation may overlap unresolved work. | Conservative startup reconciliation and process-tree evidence; uncertainty blocks unsafe rerun. |
| Strict `DELETE` means no-JavaScript cancellation lacks browser parity. | Risk | Users without JavaScript need another cancellation path. | Keep status/recovery server-rendered and provide explicit CLI cancellation guidance; revisit only under a new method requirement. |
| Local-only deployment is mistaken for trusted input. | Risk | Malicious local pages/extensions/processes may attack operations. | Apply Host/Origin/session/CSRF/body/path/confirmation checks and redaction consistently to every route/stream. |
| Selected report claims include occasional inconsistencies. | Risk | Overstating study evidence could weaken decisions. | Use only concrete cited principles corroborated across reports; identify exact numerical/API/security policies as Sprint-specific decisions. |

## Implementation Constraints

- `internal/web` may depend only on typed `internal/app` operation/query interfaces and plain safe result types; no direct `study`, `project`, `sprint`, runtime/process, or CLI-handler dependency.
- Existing workflow semantics, validators, verdicts, retry policy, runtime/harness invocation, durable state, and recovery remain product/app owned.
- Operation requests are a closed typed allowlist. Never accept raw prompts, executable/argv/shell strings, arbitrary paths, client IDs, mutation classes, fingerprints, or runtime claims.
- Preparation must be side-effect-free. Start must repeat normalization and current-fingerprint validation before ownership and mutation.
- `/api/v1` cancellation uses `DELETE`; SSE never starts, mutates, completes, fails, or cancels work.
- The web hub has no durable schema, queue, history, lock, workflow state, or database. Every accepted operation starts under immediate server ownership.
- Server, operation, subscriber, and bounded-cleanup contexts must remain distinct. No cancellable product path may replace its supplied context with `context.Background()`.
- Product locks remain held through terminal reconciliation or explicit cleanup uncertainty, including shutdown deadline expiry.
- Redaction and safe projection happen before data enters the hub, replay buffer, terminal result, JSON, HTML, or SSE.
- Enforce the exact default/upper bounds in Decision 5 through injected validated configuration; lower configured values may be allowed, unbounded values may not.
- Preserve Sprint 30 read-only routes, envelope behavior, loopback binding, Host/Origin/CSRF/session policy, security headers, body/time limits, and artifact containment.
- Browser UI remains embedded `html/template`, CSS, and minimal dependency-free JavaScript with no Node.js, framework, client router/store, service worker, or asset build step.
- Existing CLI and TUI behavior remains supported; no automatic Git operation is permitted.
- Normal tests use deterministic fakes and temporary workspaces; timing tests use fake clocks/barriers rather than sleeps where practical.
- The implementation and plan must include all required output files named in `requirements.md` and both selected review protocols.

## Plan Handoff

`plan.md` must execute these decisions. It must not invent architecture, scope, routes, states, limits, lock domains, confirmation semantics, cancellation ownership, persistence, or frontend technology beyond this document.

The plan must carry forward:

- all eight final decisions and the area-level no-JavaScript cancellation adjustment;
- selected contracts and `AC-*`/`C-*`/`RO-*` traceability;
- exact API methods, lifecycle/error/event names, and numerical limits;
- product-lock and server-shutdown ordering;
- expected focused/full verification and documentation evidence;
- assumptions, risks, trade-offs, debt, and future triggers;
- Architecture Review and Sprint Review protocols.

## Phase Exit Criteria

- [x] Selected context was read and used.
- [x] API Design, Architecture, and Frontend area reasoning documents were completed and summarized.
- [x] Area-specific conclusions are reflected; the no-JavaScript POST cancellation idea is explicitly narrowed to preserve the required `DELETE` API contract.
- [x] Contracts and requirement IDs are explicitly mapped to decisions and expected evidence.
- [x] Final decisions define package ownership, API shape, confirmation, locks, lifecycle, shutdown, SSE, security, UI, and testing without core architecture placeholders.
- [x] Expected evidence is specific and reviewable.
