# Sprint Reasoning: Local Web Hardening and Observable-Product Release

> Project: `ultraplan-go`
> Sprint: `32-hardening-and-release`
> Output: `projects/ultraplan-go/sprints/32-hardening-and-release/reasoning.md`
> **Inputs Used:** `projects/ultraplan-go/project-index.md`, `projects/ultraplan-go/sprints/32-hardening-and-release/requirements.md`, `projects/ultraplan-go/docs/ARCHITECTURE.md`, `projects/ultraplan-go/docs/PRD.md`, `projects/ultraplan-go/docs/TRD.md`, `projects/ultraplan-go/sprints/32-hardening-and-release/sprint-index.md`, `projects/ultraplan-go/sprints/32-hardening-and-release/technical-handbook.md`, `projects/ultraplan-go/sprints/32-hardening-and-release/reasoning/api-design.md`, `projects/ultraplan-go/sprints/32-hardening-and-release/reasoning/architecture.md`, `projects/ultraplan-go/sprints/32-hardening-and-release/reasoning/frontend.md`

This document decides. It synthesizes selected context, handbook evidence, area-specific reasoning, and contracts into final sprint decisions. It does not replace `sprint-index.md`, `technical-handbook.md`, or `reasoning/*.md`.

Requirement references use `AC-01` through `AC-24` for the acceptance criteria in their listed order in `requirements.md`, and `C-01` through `C-14` for its constraints in their listed order. `OUT-*` references name rows in Required Outputs.

## Sprint Purpose

- **Goal:** Harden the Sprint 30-31 local browser implementation into a secure, accessible, compatibility-controlled, recoverable, single-binary release and prove that a future workflow stage can use shared application capabilities without route-specific workflow orchestration.
- **Non-Goals:** No new workflow stage; hosted, LAN, remote, multi-user, or account behavior; client framework/router/store or WebSockets; database, durable web queue, detached worker, or second scheduler; arbitrary browser file editing; issue tracking; automatic repair or Git mutation; and no Product Phase 5 or later content, QA, retrieval, persistence, graph, cloud, or Aren work.
- **Depends On:** Sprint 30's local-web foundation and review findings, Sprint 31's guarded-operation/SSE implementation and remaining shutdown uncertainty gap, Sprint 28's applicable release findings, and the Phase 4 PRD/TRD/Architecture contracts selected by `sprint-index.md`.

## Selected Context And Pre-Reasoning Artifacts

| Artifact | Path | How It Was Used |
| --- | --- | --- |
| Project Index | `projects/ultraplan-go/project-index.md` | Established the Phase 4 boundary, selected contract/evidence catalog, implementation repository, and required review protocols. |
| Sprint Requirements | `projects/ultraplan-go/sprints/32-hardening-and-release/requirements.md` | Supplied the binding outputs, 24 acceptance criteria, 14 implementation constraints, non-goals, dependencies, and release checks. |
| Project Architecture | `projects/ultraplan-go/docs/ARCHITECTURE.md` | Fixed module ownership, `internal/web -> internal/app`, product-owned durable truth, server-owned operation lifecycle, and presentation hierarchy. |
| PRD | `projects/ultraplan-go/docs/PRD.md` | Fixed Phase 4 product behavior, one-core/multiple-surface parity, guarded browser operations, recovery truth, and release definition. |
| TRD | `projects/ultraplan-go/docs/TRD.md` | Fixed HTTP/SSE routes and separation, security and shutdown requirements, embedded assets, test layers, and configuration/runtime constraints. |
| Sprint Index | `projects/ultraplan-go/sprints/32-hardening-and-release/sprint-index.md` | Limited reasoning to the 13 selected contracts, 13 evidence reports, three area templates, and three review protocols; its excluded context remains excluded. |
| Technical Handbook | `projects/ultraplan-go/sprints/32-hardening-and-release/technical-handbook.md` | Supplied studied patterns for thin adapters, explicit composition, typed errors, bounded concurrency, cleanup contexts, capability extensibility, security, observability, and layered tests. |
| Prior Decision | None cataloged | No project-index prior decision exists. Sprint 30, Sprint 31, and Sprint 28 are dependencies through requirements, not selectable prior-decision artifacts. |

## Area-Specific Reasoning Inputs

| Area | Reasoning Document | Key Conclusion | Evidence Basis | Impact On Final Decision |
| --- | --- | --- | --- | --- |
| API Design | `reasoning/api-design.md` | Harden the documented `/api/v1` contract without redesign; use typed mappings, retry-safe confirmation consumption, idempotent cancellation, bounded replay with explicit gaps, and durable refresh. | Handbook reports on thin delegates, typed errors, context/concurrency, compatibility tests, capabilities, security, and bounded streaming. | Fixes Decision 2's wire contract and Decision 4's operation/SSE protocol. |
| Architecture | `reasoning/architecture.md` | Perform a small boundary refactor first; keep web as transport/presentation, expose narrow app capabilities, retain product-owned durable state, and centralize bounded operation shutdown. | Project architecture plus all selected reports, especially Helm/Restic boundaries, gh-cli composition, Restic cleanup, and Rclone capabilities. | Fixes Decisions 1, 4, and 5; route-specific workflow logic and web-owned durability are prohibited. |
| Frontend | `reasoning/frontend.md` | Refactor server-first presentation into namespaced downward layers, complete no-JavaScript pages, disposable enhancement state, and explicit accessibility/recovery behavior. | Handbook evidence on safe error projection, injected boundaries, context ownership, bounded concurrency, test layering, redaction, and streaming. | Fixes Decision 6's template, CSS, JavaScript, accessibility, cache, and packaging direction. |

No area conclusion is overridden. The final decisions combine them and resolve shared concerns once.

## Sprint Technical Handbook Summary

- **Relevant Patterns:** Thin transport adapters over actions; explicit manual composition and narrow injectable seams; typed error classification separated from rendering; one work-cancellation lineage plus a distinct bounded cleanup context; bounded structured concurrency and exact-once teardown; bounded event retention with explicit slow-consumer behavior; merged configuration validated before serving; type-level and projection-level redaction; layered compatibility, behavior, fake, integration, and gated tests; fixed capability discovery instead of route branching.
- **Important Trade-Offs:** Narrow interfaces can become fragmented or broad god objects; bounded replay creates gaps; cleanup independent of the work context needs a hard deadline; golden fixtures require semantic review; richer diagnostics increase leakage risk; server-rendered/no-JavaScript behavior accepts more full-page refreshes.
- **Warnings / Anti-Patterns:** No route/template workflow ownership, service-locator contexts, package globals, detached work contexts, fire-and-forget goroutines, direct output bypasses, string-parsed errors, silent config fallback, duplicate registries, default-allow trust, unbounded waits/retention, speculative pooling, or tests tied only to implementation details.
- **Evidence Confidence:** High. Every selected report is rated high-confidence by the handbook and cites concrete mature Go repositories. The evidence strongly supports structural and operational patterns, while project requirements remain authoritative for exact product semantics and limits.

## Contracts Applied

| Contract / Requirement ID | Constraint | Decision Impact | Expected Evidence |
| --- | --- | --- | --- |
| Architecture; `AC-08`-`AC-11`; `C-01` | Web depends only on app/stdlib; workflow and templates do not own product rules. | Narrow app capability vocabulary, thin handlers, typed view models, downward template composition. | Import-boundary check, fake-stage capability test, template parse/render tests, Architecture Review. |
| CLI Surface; `AC-01`, `AC-02`, `AC-06`, `AC-22` | `ultraplan serve` remains stable, loopback-only, single-binary, and consistent with CLI/TUI. | Preserve command semantics and document flags/lifecycle rather than redesigning the CLI. | Build/help/config documentation checks and cross-surface fixtures. |
| Configuration; `AC-02`, `AC-16`; `C-04`, `C-10` | Precedence and every security/resource bound must be explicit and coherent. | Validate immutable merged configuration before listening; invalid combinations fail closed. | Table-driven precedence, invalid-config, startup, and documented-default tests. |
| Errors; `AC-03`, `AC-04`, `AC-15` | Stable machine codes and safe projections must agree across JSON, HTML, and SSE. | Typed `errors.Is`/`errors.As` mapping and frozen envelopes/statuses; no string parsing or raw causes. | Route matrix/golden fixtures, semantic status tests, hostile-error/redaction fixtures. |
| Security; `AC-02`, `AC-14`, `AC-15`; `C-03`, `C-05` | Fail-closed loopback/same-origin/browser/path/redaction boundary. | Numeric loopback, Host/Origin/CSRF/session/CSP/body/path protections and layered allowlisting/redaction. | Security regression matrix, IPv4/IPv6 checks, hostile Markdown/path/secret tests, manual audit. |
| Workflows; `AC-07`, `AC-08`, `AC-17`-`AC-20`; `C-06`-`C-09` | Commands, progress, cancellation, terminal arbitration, and recovery remain shared and truthful. | Product-owned mutation/recovery, server-owned operation cancellation, disconnect-is-subscription semantics. | Representative study/sprint journeys, exact-once race tests, restart reconciliation fixtures. |
| Persistence And Migrations; `AC-06`, `AC-18`, `AC-19`; `C-02`, `C-08` | Durable workspace/run state is authoritative and atomic; web state is ephemeral. | No web database/queue; persist durable terminal or owner-specific uncertainty before shutdown closes. | Atomic-state and abrupt-restart tests; evidence that artifact/process presence never implies success. |
| Observability; `AC-06`, `AC-15`, `AC-18`, `AC-23` | Diagnostics must explain status, gaps, cancellation, cleanup, and blocked evidence safely. | Structured allowlisted correlation/lifecycle fields; blocked is not pass. | Captured logs/events with required fields and forbidden-value scans. |
| Performance; `AC-05`, `AC-16`, `AC-17`, `AC-20`; `C-04` | All request, operation, event, subscriber, stream, polling, and cleanup paths are bounded. | Fixed configurable positive limits; slow consumers are isolated; no speculative pooling. | Limit/overflow/slow-client/lifetime tests, race suite, leak checks. |
| Testing; `AC-01`, `AC-03`-`AC-24`; `C-11`, `C-12` | Deterministic normal tests plus explicit gated real-system proof. | Separate compatibility, security, lifecycle, concurrency, template/accessibility, integration, packaging, and gated layers. | `go test ./...`, `go test -race ./...`, focused suites, build/package run, gated run or truthful blocked record. |
| Documentation; `AC-22`, `AC-23` | Public commands, routes, fields, bounds, states, security, recovery, accessibility, and packaging must match implementation. | Documentation is a release artifact checked against fixtures/config, not prose added afterward. | Review of all required docs and release checklist against implementation and API fixtures. |
| LLM Runtime; LLM Evaluation / Cost / Safety; `AC-07`, `AC-15`, `AC-23`; `C-06`, `C-12` | Browser runtime work uses shared app/runtime ownership, bounded confirmation, safe metadata, cancellation, and gated proof. | No direct runtime/process import in web; runtime/model shown safely before confirmation; missing prerequisites block. | Fake runtime/harness journeys, metadata/redaction checks, gated real runtime/harness evidence. |

## Repos Studied / Source Evidence Used

The reports below were consumed through `technical-handbook.md`; the final reasoning relies on the handbook's cited findings rather than independently widening the selected evidence set.

| Source / Repo / Report | Concrete Reference | Relevant Finding | Why It Matters For This Sprint | Used In Decision(s) |
| --- | --- | --- | --- | --- |
| Project structure report; Helm and Restic | `01-project-structure.md`; `helm/pkg/action/install.go:73-140`; `restic/internal/restic/repository.go:18` | Thin entrypoints and inward-defined domain boundaries maintain unidirectional dependencies. | Grounds web-to-app direction and product ownership. | 1 |
| Command architecture report; Helm and gh-cli | `02-command-architecture.md`; `helm/pkg/cmd/install.go:132-145`; `gh-cli/pkg/cmdutil/factory.go:16-43` | Transport delegates to reusable actions/factories. | Supports thin handlers and one shared capability model. | 1, 2 |
| Dependency injection report; gh-cli and Restic | `03-dependency-injection.md`; `gh-cli/pkg/cmd/factory/default.go:26-46`; `restic/internal/backend/backend.go:19-90` | Explicit composition and narrow seams improve testability; globals obscure ownership. | Grounds executable-edge construction and deterministic fakes. | 1, 7 |
| Configuration report; Chezmoi and K9s | `04-configuration-management.md`; `chezmoi/internal/cmd/config.go:2253-2287`; `k9s/internal/config/k9s.go:423-451` | Preserve explicit precedence and validate the merged result centrally. | Makes bind and all resource/security limits startup invariants. | 3 |
| Error handling report; Helm, Restic, gh-cli | `05-error-handling.md`; `helm/pkg/storage/driver/driver.go:27-48`; `restic/internal/errors/fatal.go:10-53`; `gh-cli/internal/ghcmd/cmd.go:281-301` | Machine classification and safe outer rendering are separate concerns. | Grounds stable API codes and safe HTML/JSON/SSE mapping. | 2, 6 |
| I/O abstraction report; gh-cli, Restic, Chezmoi | `06-io-abstraction.md`; `gh-cli/pkg/iostreams/iostreams.go:551-568`; `restic/internal/fs/interface.go:10-31`; `chezmoi/internal/cmd/applycmd_test.go:220-241` | Injectable boundaries enable deterministic transport tests. | Grounds fakes, `httptest`, and no direct output/filesystem bypasses. | 1, 7 |
| State/context report; Helm and Restic | `07-state-context.md`; `helm/pkg/cmd/install.go:333-347`; `restic/internal/restic/lock.go:290-305` | Work has one cancellation lineage; cleanup may need a distinct bounded context. | Grounds browser-disconnect semantics and shutdown reconciliation. | 4, 5 |
| Concurrency report; Go-task, K9s, OpenCode, Rclone | `08-concurrency.md`; `go-task/task.go:87`; `k9s/internal/pool.go:21-37`; `opencode/cmd/root.go:252-279`; `rclone/lib/batcher/batcher.go:50` | Bound fan-out, wait explicitly, time out cleanup, and arbitrate completion once. | Grounds operation/SSE caps, exact-once cancellation, and leak-free shutdown. | 4, 5, 7 |
| Logging report; Helm and K9s | `10-logging-observability.md`; `helm/internal/logging/logging.go:31-71`; `k9s/internal/slogs/keys.go:6-231` | Structured fields and output/diagnostic separation improve safe diagnosis. | Grounds correlation fields and forbidden-payload rules. | 3, 5, 7 |
| Testing report; gh-cli, Helm, Restic | `11-testing-strategy.md`; `gh-cli/acceptance/acceptance_test.go:26-29`; `helm/internal/test/test.go:43`; `restic/internal/backend/mock/backend.go:14-26` | Acceptance pipelines, deliberate goldens, and centralized fakes protect public behavior. | Grounds the release evidence matrix and fixture-review discipline. | 2, 6, 7 |
| Extensibility report; Rclone and Dive | `12-extensibility.md`; `rclone/fs/features.go:294-370`; `dive/cmd/dive/cli/internal/command/adapter/analyzer.go:13-15` | Fixed optional capabilities can extend behavior without exposing orchestration. | Grounds the fake future-stage proof and rejects a plugin/route registry. | 1 |
| Security report; Restic, Helm, OpenCode | `13-security.md`; `restic/internal/options/secret_string.go:15-20`; `helm/pkg/registry/transport.go:37-41`; `opencode/internal/permission/permission.go:44-108` | Trust boundaries need structured validation, secret-safe values, scrubbing, and bounded permissions. | Grounds fail-closed browser security, confirmation, and layered redaction. | 2, 3, 6 |
| Performance report; Lazygit, Rclone, K9s | `14-performance.md`; `lazygit/pkg/tasks/tasks.go:189-217`; `rclone/lib/pool/pool.go:17-24,52-53`; `k9s/internal/pool.go:26-48` | Stream responsively with bounded queues/resources; optimize only after evidence. | Grounds lossy bounded SSE, subscriber eviction, explicit caps, and no speculative cache/pool. | 4, 6 |

## Trade-Off And Debt Analysis

### Accepted Trade-Offs

| Trade-Off | Benefit | Cost / Constraint Accepted | Why Acceptable Now | Revisit Trigger |
| --- | --- | --- | --- | --- |
| Small app/web boundary refactor before hardening | One testable product core and enforceable ownership. | Moves existing route logic before visible release work. | Hardening duplicated workflow logic would preserve the core defect. | Never for convenience; reopen only if a documented use case cannot fit cohesive app capabilities. |
| Preserve `/api/v1` rather than redesign | Protects Sprint 30-31 consumers and makes changes reviewable. | Existing awkward route/field choices remain. | Sprint goal is release compatibility, not API redesign. | A separately governed versioned API change. |
| Fixed capability vocabulary, not plugins | Future-stage proof without route branching or hidden registration. | Production assembly remains explicit and less dynamically extensible. | Only in-process known product modules are required. | Multiple independently shipped capability providers with demonstrated need. |
| Bounded ephemeral operation/SSE state | Protects memory, shutdown, and product work from browsers. | Replay gaps and expired operation handles require durable refresh. | Durable product state already supplies recovery truth. | A future durable-worker/queue architecture with explicit authority. |
| Separate bounded cleanup context | Allows cancellation plus truthful terminal reconciliation. | Adds cleanup deadlines and `cleanup_uncertain` states. | Uncertainty is safer than false success or detached cleanup. | A product-owned worker can durably assume cleanup ownership. |
| Server-rendered no-JavaScript baseline | Accessibility, recoverability, small attack surface, single binary. | Some transitions refresh whole pages; imperative enhancement remains. | Current interaction complexity does not justify a client app. | Measured workflows cannot be implemented safely with narrow enhancement. |
| Fail-closed startup | Prevents serving with unsafe limits/templates/security policy. | Invalid operator config makes the server unavailable. | Explicit failure is safer and diagnosable. | No revisit unless a specific setting can safely degrade with documented semantics. |
| Exact fixtures plus semantic/manual checks | Detects compatibility and accessibility regressions precisely. | Fixture maintenance and manual review cost. | Release claims span machine shape and human interaction. | Tooling may automate more checks, but semantic review remains required. |

### Potential Technical Debt

| Debt / Shortcut | Why It Might Accrue | Current Mitigation | Owner / Follow-Up |
| --- | --- | --- | --- |
| Capability vocabulary may grow too broad | New stages may add unrelated fields to one interface. | Keep consumer-oriented cohesive interfaces and concrete result types; fake-stage contract guards route independence. | Reassess after Sprint 33 uses it; split only on demonstrated cohorts. |
| Imperative progressive-enhancement code | More interactions can create timer, focus, and DOM-state complexity. | Separate `app.js`, `operations.js`, and `sse.js`; one abort owner; bounded retries; browser fixtures. | Revisit only after measured complexity triggers a frontend architecture review. |
| Golden fixture update burden | Exhaustive route/render matrices can become noisy. | Require per-change compatibility rationale plus semantic assertions; prohibit wholesale blind regeneration. | Release maintenance; prune only redundant implementation-detail snapshots. |
| Resource-limit tuning | Conservative defaults may reject legitimate substantial workflows. | Document every limit, make safe limits configurable, test representative study/sprint workflows, expose safe capacity errors. | Tune from observed local workloads without weakening hard bounds. |
| Manual accessibility evidence | Static tests cannot prove focus visibility, reflow, announcement timing, or color independence. | Release checklist records keyboard, zoom, narrow-layout, reduced-motion, and assistive checks. | Keep until reliable browser automation covers each property; manual audit still samples UX. |
| Compatibility baseline ambiguity | Requirements freeze existing documented fields but do not enumerate them here. | Inventory implementation/docs before DTO changes and treat unexplained differences as blockers. | First implementation task owns the baseline manifest and fixtures. |

### Future Considerations

| Consideration | Deferred Until | Reason Deferred | What Should Be Preserved Now |
| --- | --- | --- | --- |
| New `code-context` stage | Sprint 33 | Explicit non-goal, but this sprint must prove the seam. | Status/artifact/command/progress/cancel/recovery capability contract with no route branch. |
| Frontend framework or WebSockets | Demonstrated client/bidirectional complexity | Current server-first/SSE model satisfies requirements with lower risk. | Versioned HTTP/app boundary and disposable client state. |
| Durable workers or operation queue | Explicit detached-work product requirements | Would change server ownership, idempotency, persistence, and shutdown semantics. | Product-owned durable state and canonical cancellation/recovery interfaces. |
| Hosted or multi-user service | New authentication/authorization/tenancy phase | Loopback session and same-origin controls are not remote identity. | Explicit trust-boundary separation and versioned DTOs; do not imply remote safety. |
| Alternate persistence/SQLite | Measured workflows demonstrate value and authority is decided | A web-side store would become shadow truth. | Focused product-owned persistence seams and atomic filesystem authority. |
| Lossless event history | Durable event-history requirement | Unbounded replay conflicts with local resource guarantees. | Monotonic IDs, explicit gaps, and durable refresh semantics. |

## Decisions

The seven final decisions below are binding: enforce the shared app boundary, freeze the existing `/api/v1` contract, fail closed on browser security and resource configuration, keep operation/SSE state bounded and ephemeral, make shutdown reconciliation product-owned and truthful, ship an accessible server-first embedded presentation, and require layered deterministic plus gated release evidence.

## Final Decisions

### Decision 1: Enforce A Thin Web Boundary And Shared Capability Model

- **Decision:** Refactor any remaining route-specific workflow, file/state interpretation, runtime/process access, or CLI coupling out of `internal/web`. `internal/web` may import only `internal/app` and standard-library packages. The executable composition root wires narrow app capabilities for status/readiness, findings, artifacts, commands, preparation, execution, progress, cancellation, terminal results, and durable recovery. A fake future-stage fixture must expose the full vocabulary without a new route branch or production registry.
- **Rationale:** Security, compatibility, cross-surface agreement, and lifecycle tests are reliable only when one product core owns behavior. This is a bounded refactor, not a generic plugin architecture.
- **Study / Source Grounding:** `technical-handbook.md` patterns "Transport adapter over shared capabilities," "Explicit composition root," and "Capability discovery"; `01-project-structure.md` via Helm/Restic; `02-command-architecture.md` via Helm/gh-cli; `03-dependency-injection.md` via gh-cli/Restic; `12-extensibility.md` via Rclone optional capabilities.
- **Trade-Offs Accepted:** More explicit mapping and constructors; coordinated capability evolution. In return, dependencies, ownership, and fakes stay visible.
- **Technical Debt / Future Impact:** Avoid a god interface by grouping cohesive consumer needs. Sprint 33 may test and refine the seam but cannot add route orchestration.
- **Alternatives Rejected:** Direct web imports of product/runtime packages, CLI subprocess/text parsing, route-per-stage handlers, package-global runners, reflection DI, context service locators, a generic plugin/route registry, and a broad `WebService`; each hides or duplicates product ownership.
- **Contracts Satisfied:** Architecture, CLI Surface, LLM Runtime, Workflows, Testing; `AC-06`-`AC-11`, `C-01`, `C-06`, `OUT-Web application capability`, `OUT-Shared operation capability`.
- **Evidence Required:** Import-boundary inspection/test; handler tests proving typed delegation; fake future-stage capability test; CLI/TUI/web agreement fixtures; Architecture Review with no direct product/runtime/process/CLI dependency.

### Decision 2: Freeze And Safely Map The Existing `/api/v1` Contract

- **Decision:** Preserve every documented Sprint 30-31 `/api/v1` route, method, success/error envelope, field meaning/type/nullability, status mapping, and stable error code. DTOs stay web-owned and map typed app results; errors use a small fixed actionable vocabulary via `errors.Is`/`errors.As`. Unknown `/api/` paths and wrong methods always return JSON errors. State-bearing HTML/API/operation/SSE responses use `no-store`; static assets revalidate unless proven content-addressed. Start retries deduplicate by consumed confirmation, and cancellation is idempotent.
- **Rationale:** This sprint is the compatibility release gate. Internal model evolution must not leak into wire behavior, and browser retries must not duplicate mutations.
- **Study / Source Grounding:** Handbook typed-error, thin-delegate, configuration, testing, and security patterns; `05-error-handling.md` via Helm/Restic/gh-cli; `11-testing-strategy.md` via Helm goldens and gh-cli acceptance; `13-security.md` via type/boundary redaction.
- **Trade-Offs Accepted:** Existing awkward API choices persist; exact fixtures are costly; each new code is a compatibility obligation.
- **Technical Debt / Future Impact:** Any breaking change requires separately versioned design. Optional field additions require explicit rationale and tests. First implementation work must inventory the current documented baseline.
- **Alternatives Rejected:** Redesigning envelopes during hardening, exposing app/domain structs, parsing display strings, raw internal errors, treating start as non-retryable, loose substring tests, and blind golden regeneration; all risk compatibility, duplicate work, or leakage.
- **Contracts Satisfied:** Errors, Security, CLI Surface, Documentation, Testing, Observability; `AC-03`-`AC-05`, `AC-15`, `C-05`, `OUT-Stable API routes and envelopes`, `OUT-HTTP and page handlers`, `OUT-API compatibility fixtures`.
- **Evidence Required:** Exhaustive route/method matrix and reviewed goldens covering additions/omissions/types/nullability/status/code/content type/cache; semantic tests for deduplication, idempotent cancellation, unknown routes, safe errors, and post-expiry durable refresh; docs checked against fixtures.

### Decision 3: Validate A Fail-Closed Local Security And Resource Policy Before Serving

- **Decision:** Resolve built-in defaults, workspace config, environment, then command flags into one immutable server configuration and validate it before listening. Require numeric IPv4/IPv6 loopback, strict Host and exact same-origin policy, CSRF/session protection, CSP and security headers, JSON/body/time/path/reference limits, bounded previews, and coherent positive operation/preparation/event/result/subscriber/stream/retention/heartbeat/lifetime/polling/cleanup limits. Use app-resolved opaque/allowlisted references, safe Markdown, constant-time bearer-secret comparisons, and type plus final-projection redaction across HTML, JSON, SSE, retained data, and diagnostics.
- **Rationale:** Loopback reduces exposure but does not eliminate browser-origin, path, hostile-content, forged-reference, cache, or secret-leak threats. Silent correction would make operator intent and release behavior unknowable.
- **Study / Source Grounding:** Handbook merged-config, redaction, trust-boundary, observability, and bounded-resource patterns; `04-configuration-management.md` via Chezmoi/K9s; `10-logging-observability.md` via Helm/K9s; `13-security.md` via Restic/Helm/OpenCode; `14-performance.md` via bounded pools/streams.
- **Trade-Offs Accepted:** Invalid configuration prevents startup and conservative defaults may reject bursts. Rich diagnostics are constrained to safe allowlisted fields.
- **Technical Debt / Future Impact:** Limit defaults require field data and may be tuned, but hard bounds and fail-closed validation remain. The browser session is not remote authentication and must never be represented as such.
- **Alternatives Rejected:** Hostname or non-loopback binding, permissive CORS, warning-and-fallback configuration, raw caller paths/commands, default-allow references, single-pass ad hoc redaction, logging bodies/provider payloads, and unbounded limits; each violates the local trust model.
- **Contracts Satisfied:** Security, Configuration, Performance, Observability, LLM Evaluation / Cost / Safety; `AC-02`, `AC-14`-`AC-16`, `C-03`-`C-05`, `C-10`, `OUT-Browser security policy`.
- **Evidence Required:** IPv4/IPv6 bind and Host/Origin/CSRF/session/CSP/header/body tests; malformed/duplicate-significant JSON, path escape, forged reference, hostile Markdown, cache, and secret corpus tests; invalid and incoherent config startup failures; captured-output scan for every forbidden value class.

### Decision 4: Keep Operations And SSE Bounded, Ephemeral, And Retry-Safe

- **Decision:** The server hub retains only bounded operation IDs, normalized safe summaries, confirmation/deduplication records, canonical cancel functions, redacted events, subscriber queues, and short-lived terminal projections. Confirmation is short-lived and bound to normalized request, scope, mutation class, governed-input fingerprint, and session; consumption plus operation publication is atomic. Event IDs increase per operation, replay is bounded, stale cursors receive an explicit non-terminal gap signal, and slow subscribers are evicted without backpressure. Browser disconnect affects only the subscription; explicit HTTP cancellation or server shutdown controls work.
- **Rationale:** Responsive progress cannot outrank product execution, memory bounds, or durable truth. Retry safety and explicit gaps make normal browser/network ambiguity recoverable.
- **Study / Source Grounding:** Handbook context, bounded concurrency, streaming, exact-once teardown, and performance patterns; `07-state-context.md` via Helm/Restic; `08-concurrency.md` via Go-task/K9s/OpenCode/Rclone; `14-performance.md` via Lazygit/Rclone.
- **Trade-Offs Accepted:** Replay may be incomplete, operation IDs expire, terminal SSE can be missed, and clients sometimes perform full durable refresh.
- **Technical Debt / Future Impact:** No durable operation history or detached execution is introduced. Such behavior requires a future worker/queue authority decision, not larger hub retention.
- **Alternatives Rejected:** Unlimited replay, blocking subscribers, browser-owned persistent state, WebSockets, SSE mutation commands, cancellation on disconnect, fire-and-forget goroutines, and a web scheduler/database; these create unbounded resources or a second authority.
- **Contracts Satisfied:** Performance, Workflows, Observability, Persistence And Migrations, LLM Runtime; `AC-05`, `AC-07`, `AC-16`, `AC-17`, `AC-20`, `C-02`, `C-04`, `C-06`, `C-07`, `OUT-Operation HTTP handlers`, `OUT-Operation and SSE hub`, `OUT-SSE enhancement`.
- **Evidence Required:** Deterministic clocks/IDs/barriers for confirmation expiry/staleness/deduplication, event ordering/IDs/replay/gaps/rollover/heartbeat/terminal flush, capacity and payload limits, slow subscribers, disconnect isolation, cancellation races, and operation expiry; `go test -race` and leak checks.

### Decision 5: Make Shutdown And Restart Reconciliation Truthful And Product-Owned

- **Decision:** One server-owned work context propagates through app/product/runtime/process work. A central terminal arbiter handles completion, failure, timeout, user cancellation, and shutdown exactly once. Shutdown enters draining, rejects new mutations, snapshots active cancellations outside hub locks, invokes canonical cancellation once with `server_shutdown`, then uses a distinct bounded cleanup context for process-tree cleanup, product-owned durable reconciliation, and lock resolution. HTTP/SSE closes only after a durable authoritative terminal result or owner-specific `cleanup_uncertain` equivalent is recorded. Startup reconciles stale active state conservatively; process absence and artifact presence never prove success.
- **Rationale:** This closes the explicit Sprint 31 gap without detaching work, releasing locks prematurely, or converting uncertainty into false success.
- **Study / Source Grounding:** Handbook work-versus-cleanup and exact-once patterns; `07-state-context.md` via Helm signal propagation and Restic delayed cleanup; `08-concurrency.md` via timed OpenCode cleanup and Rclone `sync.Once`; `10-logging-observability.md` for safe lifecycle diagnostics.
- **Trade-Offs Accepted:** Lifecycle state and synchronization become more complex; shutdown may persist uncertainty rather than a convenient terminal claim.
- **Technical Debt / Future Impact:** Cleanup contexts remain purpose-specific and bounded. Future durable workers may transfer ownership explicitly; until then the server cannot detach active work.
- **Alternatives Rejected:** Reusing the cancelled work context for cleanup, unbounded detached cleanup, closing HTTP before reconciliation, waiting or doing I/O under hub locks, per-callback terminal writes, inferring success from files/processes, and releasing product locks before ownership is resolved; each risks corruption, deadlock, leaks, or false state.
- **Contracts Satisfied:** Architecture, Workflows, Persistence And Migrations, Errors, Observability, LLM Runtime; `AC-18`-`AC-20`, `C-06`, `C-08`, `C-09`, `OUT-Server lifecycle`, `OUT-Shared operation capability`.
- **Evidence Required:** Multi-operation shutdown tests for draining/rejection, exact-once cancel, explicit reason, no I/O-under-lock behavior, bounded waits, terminal race arbitration, durable clean and uncertain outcomes, abrupt restart, stale locks, forced interruption, conservative reconciliation, process-tree cleanup metadata, and no orphaned goroutines/locks/operations.

### Decision 6: Ship A Server-First, Accessible, Embedded Presentation System

- **Decision:** Parse one embedded, namespaced template tree at startup with strict `page -> layout -> component -> primitive` dependencies and explicit typed view models. Organize CSS as tokens/base/primitives/components/layouts/utilities and JavaScript as dependency-free app/operation/SSE enhancements. Every representative route renders complete semantic no-JavaScript navigation, inspection, confirmation, operation, cancellation, terminal, error, and recovery states. Browser enhancement state is disposable; after terminal events, gaps, expiry, or reconnect failure it refreshes durable server truth. Accessibility includes semantic landmarks/headings, native controls, labels/error associations, visible focus, restrained live regions, text/non-color states, reduced motion, keyboard operation, 200% zoom, and narrow reflow.
- **Rationale:** A server-first UI preserves one source of truth, offline single-binary packaging, safe rendering, and baseline accessibility while still allowing responsive progress enhancement.
- **Study / Source Grounding:** Handbook safe rendering, deterministic boundary testing, cancellation ownership, bounded enhancement, layered tests, redaction, and streaming patterns; `05`, `06`, `07`, `08`, `10`, `11`, `13`, and `14` selected reports as synthesized in `frontend.md`. No selected study prescribes the exact template hierarchy; that exact hierarchy comes from binding Architecture/TRD/requirements, while the studies support its boundary and test properties.
- **Trade-Offs Accepted:** Some actions use full-page refreshes, progressive enhancement is imperative, and manual accessibility evidence remains necessary.
- **Technical Debt / Future Impact:** Keep utilities presentation-only and scripts narrowly owned to limit drift. A framework requires measured interaction complexity and a new decision; template verbosity is insufficient justification.
- **Alternatives Rejected:** SPA/client router/store, framework/build pipeline, service worker, flat/dynamic templates, `map[string]any` or domain objects in templates, template filesystem/app calls, utility-only state styling, custom div controls, third-party scripts/fonts, and source-tree assets; these add authority, accessibility, security, or packaging risk.
- **Contracts Satisfied:** Architecture, Security, Documentation, Testing, Performance; `AC-10`-`AC-15`, `AC-17`, `AC-21`, `C-07`, `C-10`, `OUT-Template primitives/components/layouts/pages`, `OUT-CSS tokens/base/primitives/components/layouts/utilities`, `OUT-Baseline browser enhancement`, `OUT-Guarded-operation enhancement`, `OUT-SSE enhancement`.
- **Evidence Required:** Startup failures for duplicate/missing/invalid definitions; downward-composition and explicit-view-model tests; hostile rendering/URL/attribute fixtures; complete no-JavaScript journeys; browser enhancement tests; semantic assertions; manual keyboard/focus/announcement/color/reduced-motion/zoom/narrow checks; binary run outside source tree requesting every page/asset.

### Decision 7: Gate Release With Layered Deterministic, Review, Documentation, And Real-System Evidence

- **Decision:** Implement separate evidence layers: typed app capability tests; exhaustive API compatibility fixtures; security/redaction tests; operation/SSE concurrency, race, leak, cancellation, and lifecycle tests; template/accessibility tests; representative study and sprint integration over temporary workspaces with fake runtime/harness dependencies; cross-surface agreement; outside-tree packaging; documentation reconciliation; and release commands. Then run the selected Architecture Review and Sprint Review. Gated evidence must exercise one real runtime-backed and one real smoke-harness-backed browser operation when available; absent prerequisites produce `blocked`, never pass.
- **Rationale:** No single end-to-end test can localize or prove compatibility, security, concurrency, accessibility, recovery, and packaging. Deterministic tests establish repeatable behavior; gated checks cover only real boundaries.
- **Study / Source Grounding:** Handbook layered-testing pattern; `06-io-abstraction.md` for injectable boundaries; `11-testing-strategy.md` via gh-cli acceptance, Helm goldens, and Restic fakes; `08-concurrency.md` for race/leak pressure; `10-logging-observability.md` for safe evidence. The exact required commands and review protocols come from requirements and `sprint-index.md`, not external studies.
- **Trade-Offs Accepted:** The release matrix is substantial and manual accessibility/documentation review cannot be collapsed into normal unit tests.
- **Technical Debt / Future Impact:** Keep gated tests narrowly scoped and truthful to avoid flaky live dependencies. A narrow rerun cannot replace the required suite. Documentation must remain fixture/config-driven as the API evolves.
- **Alternatives Rejected:** One broad browser E2E suite, live providers in normal tests, snapshot-only verification, automation-only accessibility claims, runtime exit/artifact presence as success, skipped prerequisites reported as pass, and release verification that mutates Git; each weakens evidence or reproducibility.
- **Contracts Satisfied:** Testing, Documentation, Security, Observability, CLI Surface, LLM Runtime, LLM Evaluation / Cost / Safety; `AC-01`, `AC-03`, `AC-06`-`AC-24`, `C-11`-`C-14`, all test/doc/release required outputs.
- **Evidence Required:** Passing `go test ./...`, `go test -race ./...`, `go build ./cmd/ultraplan`, focused web/app tests, compatibility/security/accessibility/lifecycle/package reports, synchronized docs and release checklist, Architecture Review, Sprint Review with no unresolved applicable high finding, and Deep Smoke evidence or explicit blocked prerequisites.

## Expected Evidence

| Evidence Type | Required Evidence | Source / Command / Review Check |
| --- | --- | --- |
| Tests | Full deterministic suite passes. | From `../ultraplan-go`: `go test ./...`. |
| Race / leaks | Operation, SSE, shutdown, reconciliation, subscriber, and cancellation tests pass under race detection with no leaked owners. | From `../ultraplan-go`: `go test -race ./...`; focused `internal/web` and `internal/app` suites. |
| Build | Single binary builds with embedded assets. | From `../ultraplan-go`: `go build ./cmd/ultraplan`. |
| API compatibility | Every documented route/method/envelope/field/status/code/content type/cache rule is frozen and semantically reviewed. | `internal/web/api_compatibility_test.go` and reviewed fixtures. |
| Security | Loopback, Host/Origin, CSRF/session/CSP/headers/body/path/Markdown/reference/redaction matrix passes with no forbidden value in any projection or capture. | `internal/web/security_test.go`; manual browser security audit. |
| Runtime / lifecycle | Draining, exact-once cancellation, bounded cleanup, durable terminal/uncertainty, abrupt restart, and conservative reconciliation are demonstrated. | `internal/web/server_test.go`, `operations_test.go`, `sse_test.go`; captured redacted lifecycle diagnostics. |
| Interface agreement | App, CLI, TUI, HTML, and JSON agree for representative study/sprint states, terminal outcomes, artifacts, readiness, and next actions. | `internal/web/integration_test.go`, `internal/app/web_usecases_test.go`, temporary workspaces and fakes. |
| Frontend / accessibility | Complete no-JavaScript behavior, namespaced hierarchy, escaping, semantics, keyboard/focus/live-region/color/motion/zoom/reflow checks. | `internal/web/templates_test.go`, browser fixtures, and release checklist manual checks. |
| Packaging | Binary serves all templates/static assets outside the source tree without Node, Vite, CDN, database, or separate frontend process. | Build-and-launch packaging check from a temporary directory with asset requests. |
| Documentation | Local web, configuration, CLI, user, recovery, architecture, and release documentation match implementation and fixtures. | `docs/local-web.md`, `docs/configuration.md`, `docs/cli-reference.md`, `docs/user-guide.md`, `docs/recovery.md`, `docs/architecture.md`, `docs/release-checklist.md`. |
| Review | Architecture and Sprint Review cover every selected contract/handbook concern and report no unresolved applicable high-severity local-web finding. | `system/protocols/architecture-review-protocol.md`; `system/protocols/review-sprint-protocol.md`. |
| Gated real system | One real runtime-backed and one real smoke-harness-backed browser operation, or a truthful blocked result naming unavailable prerequisites. | `system/protocols/deep-smoke-sprint-protocol.md`; cataloged `ultraplan-go-smoke` evidence links. |

## Assumptions And Risks

| Item | Type | Impact | Mitigation / Follow-Up |
| --- | --- | --- | --- |
| Phase 4 remains numeric-loopback-only, same-origin, single-process, and single-user. | Assumption | Remote/multi-user use would invalidate session, authority, idempotency, and threat decisions. | Fail non-loopback config and document scope; require new architecture for remote use. |
| Existing documented Sprint 30-31 API is the compatibility baseline but is not enumerated in requirements. | Risk | Accidental fixture creation could bless a breaking implementation discrepancy. | Inventory routes/docs/current DTOs first; unexplained differences block release; fixture changes need rationale. |
| Existing web code may contain route-specific workflow or direct file interpretation. | Risk | Boundary refactor size may be larger than expected. | Move violations inward before hardening; discovery changes task sizing, not Decision 1. |
| Confirmation consumption, capacity reservation, deduplication, and shutdown can race. | Risk | Duplicate product work or lost accepted operation. | Define one critical ordering and deterministic concurrent tests. |
| Completion, failure, timeout, user cancel, and shutdown can race. | Risk | Conflicting terminal results or repeated cancellation. | Single terminal arbiter and exact-once canonical cancellation with race tests. |
| Hub, product locks, persistence, subscribers, and cleanup can deadlock. | Risk | Hung operations or shutdown. | Explicit lock order; snapshot under lock; no callbacks, I/O, sends, or waits under hub locks. |
| Replay gaps or expired handles can look terminal to users. | Risk | Stale UI may imply success/failure. | Stable non-terminal gap/expiry states and mandatory durable HTTP refresh. |
| Resource defaults may be too strict or too generous. | Risk | Rejected legitimate work or weak memory/shutdown guarantees. | Test representative substantial workflows, document/tune defaults, retain hard caps. |
| Static asset URLs may not be content-addressed. | Assumption | Immutable caching could serve stale code. | Use revalidation unless content addressing is proven by packaging/route inventory. |
| Static accessibility assertions are incomplete. | Risk | Keyboard, focus, zoom, reflow, motion, or announcements may still fail. | Mandatory manual representative-page checks recorded in release checklist. |
| Real runtime or harness prerequisites may be unavailable. | Risk | Live release evidence cannot run in every environment. | Report blocked with exact missing prerequisite; never convert blocked/skipped into pass. |
| Rich diagnostics can leak sensitive runtime/browser values. | Risk | Secrets or unsafe paths reach logs/events/UI. | Allowlist fields, type/projection redaction, forbidden corpus scanning, no raw bodies/payloads/stderr. |
| Unrelated working-tree changes may exist during release verification. | Assumption | Destructive cleanup could lose user work. | No automatic Git mutation; release checks preserve unrelated changes. |

## Implementation Constraints

- Preserve the documented Sprint 30-31 `/api/v1` contract unless an additive optional field has explicit compatibility rationale and tests.
- `internal/web` may import only `internal/app` and Go standard-library packages; it must not invoke the `ultraplan` binary.
- Product modules retain workflow rules, mutation exclusion, durable state, terminal reconciliation, stale-lock policy, and recovery.
- Use explicit executable-edge composition; no global mutable runners/registries, service locator, reflection DI, or hidden background ownership.
- Validate merged web configuration, templates, route uniqueness, and security/resource invariants before accepting requests.
- Never wait, call app/product code, write files, or send to subscribers while holding operation-hub locks.
- Every goroutine, timer, event source, operation, subscriber, and cleanup action has an owner, bound, and terminal cleanup path.
- Browser/SSE/request contexts never own product work; only explicit canonical cancellation and server shutdown cancel operations.
- Persist a truthful product-owned terminal or owner-specific uncertainty before graceful shutdown closes HTTP/SSE or releases ownership.
- Browser state, operation handles, confirmations, events, subscribers, polling, and caches remain bounded and non-authoritative.
- Templates use explicit typed view models and downward namespaced composition; JavaScript remains optional, dependency-free enhancement.
- HTML, JSON, SSE, retained events, terminal results, logs, and diagnostics must independently exclude forbidden sensitive values.
- Normal tests use `httptest`, temporary workspaces, deterministic clocks/IDs/barriers, and fake app/runtime/harness dependencies; no live provider is required.
- Real runtime/harness tests remain gated and report unavailable prerequisites as blocked, not passed.
- Do not add a new workflow stage, plugin system, frontend framework, database, durable queue/worker, WebSocket, remote surface, issue workflow, repair behavior, or Git mutation.
- Preserve unrelated working-tree changes during all release checks.

## Plan Handoff

`plan.md` must execute these decisions. It must not invent architecture, scope, API redesign, persistence, frontend runtime, or lifecycle semantics beyond this document.

The plan must carry forward:

- all seven final decisions and their ordering dependency: boundary/baseline inventory before hardening and release proof
- all selected contracts and `AC-*`/`C-*`/`OUT-*` mappings
- the complete expected-evidence matrix and explicit normal-versus-gated distinction
- the lifecycle, compatibility-baseline, lock-order, resource-tuning, accessibility, and prerequisite risks
- Architecture Review, Sprint Review, and Deep Smoke Sprint protocols

## Phase Exit Criteria

- [x] Selected context was read and used.
- [x] API Design, Architecture, and Frontend area reasoning documents were completed and summarized.
- [x] All area-specific conclusions are reflected without override.
- [x] Contracts and requirement references are mapped to decisions and expected evidence.
- [x] Final decisions are specific enough for `plan.md` to execute without reopening architecture.
- [x] Expected evidence is specific, layered, reviewable, and distinguishes blocked external prerequisites from pass.
