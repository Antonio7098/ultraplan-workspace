# Sprint Reasoning: Durable Run Identity and Cross-Surface Observability

> Project: `ultraplan-go`
> Sprint: `35-durable-run-observability`
> Output: `projects/ultraplan-go/sprints/35-durable-run-observability/reasoning.md`
> **Inputs Used:** `projects/ultraplan-go/project-index.md`, `projects/ultraplan-go/sprints/35-durable-run-observability/requirements.md`, `projects/ultraplan-go/sprints/35-durable-run-observability/code-context.md`, `projects/ultraplan-go/docs/ARCHITECTURE.md`, `projects/ultraplan-go/docs/PRD.md`, `projects/ultraplan-go/docs/TRD.md`, `projects/ultraplan-go/sprints/35-durable-run-observability/sprint-index.md`, `projects/ultraplan-go/sprints/35-durable-run-observability/technical-handbook.md`, `projects/ultraplan-go/sprints/35-durable-run-observability/reasoning/api-design.md`, `projects/ultraplan-go/sprints/35-durable-run-observability/reasoning/architecture.md`, `projects/ultraplan-go/sprints/35-durable-run-observability/reasoning/frontend.md`

This document decides. It synthesizes selected context, handbook evidence, area-specific reasoning, and contracts into final sprint decisions.

It does not replace `sprint-index.md`, `technical-handbook.md`, or `reasoning/*.md`.

## Sprint Purpose

- **Goal:** Give every accepted asynchronous UltraPlan operation, including every runtime-backed CLI, TUI, and web operation, a durable workspace-scoped identity, safe replayable history, conservative liveness, authorized cross-process cancellation, one operational terminal result, and identical projections across supported local surfaces.
- **Non-Goals:** Hosted or multi-user service; LAN/public exposure; shared-filesystem multi-host ownership; remote workers; a daemon, broker, distributed queue, or general workflow engine; replacement of product-owned Markdown, flow, execute, study, Git, or smoke evidence; raw provider/prompt/output retention; browser-owned state; automatic repair or Git mutation; and later content, retrieval, graph, cloud, or authored-product-persistence work.
- **Depends On:** Sprint 31 guarded operation/SSE/cancellation behavior, Sprint 32 compatibility and browser recovery contracts, Sprint 34 shared-context dogfood, existing app operation seams, sprint/study locks and checkpoints, agentwrap supervision and canonical events, and the current loopback-only local-web security boundary.

For requirement traceability below, `AC-1` through `AC-13` refer in order to the thirteen acceptance-criteria bullets in `requirements.md`. `RO-1` through `RO-11` refer in order to the eleven required-output rows in that file.

## Selected Context And Pre-Reasoning Artifacts

| Artifact | Path | How It Was Used |
| --- | --- | --- |
| Project Index | `projects/ultraplan-go/project-index.md` | Established the module-driven architecture, selected contract pool, evidence catalog, local-only web boundary, and required review protocols. |
| Sprint Requirements | `projects/ultraplan-go/sprints/35-durable-run-observability/requirements.md` | Defined the observable outcomes, authority constraints, failure matrix, open decisions, and explicit non-goals resolved here. |
| Code Context | `projects/ultraplan-go/sprints/35-durable-run-observability/code-context.md` | Located the app acceptance seam, process-local web authority, frozen operation API, TUI bypass, agentwrap correlation boundary, and existing product persistence/lock behavior. |
| Project Architecture | `projects/ultraplan-go/docs/ARCHITECTURE.md` | Preserved inward dependency direction, product ownership, app composition, adapter-only web/TUI behavior, and separation of operational authority from authored artifacts. |
| Product Requirements | `projects/ultraplan-go/docs/PRD.md` | Carried forward local-first operation, one core across surfaces, durable observer-independent runs, and Phase 6 completion criteria. |
| Technical Requirements | `projects/ultraplan-go/docs/TRD.md` | Supplied agentwrap, API, security, concurrency, persistence, testing, and failure-matrix constraints while leaving Sprint 35 mechanisms open for this decision. |
| Sprint Index | `projects/ultraplan-go/sprints/35-durable-run-observability/sprint-index.md` | Limited reasoning to the selected contracts, fourteen evidence reports, three area templates, and three required review protocols. |
| Technical Handbook | `projects/ultraplan-go/sprints/35-durable-run-observability/technical-handbook.md` | Supplied studied patterns for thin surfaces, explicit composition, durable identity plus context, bounded concurrency, typed errors, redaction, telemetry, and layered tests; its evidence gaps were resolved rather than treated as proof. |
| Prior Decisions | No cataloged prior-decision artifact | Sprint 31, 32, and 34 behavior is a dependency through requirements, but the project index contains no selectable prior-decision record. |

## Area-Specific Reasoning Inputs

| Area | Reasoning Document | Key Conclusion | Evidence Basis | Impact On Final Decision |
| --- | --- | --- | --- | --- |
| API Design | `reasoning/api-design.md` | Add canonical `/api/v1/runs` and `ultraplan run` surfaces while preserving `/api/v1/operations` as a compatibility projection; use opaque IDs, durable sequence cursors, typed gaps/tombstones, workspace-readable safe history, and freshly authorized idempotent cancellation. | Frozen route/DTO/SSE fixtures, current session filtering, typed-error/security/performance reports, and local-web constraints. | Fixes resource shapes, compatibility behavior, read/mutation policy, error codes, replay semantics, and CLI commands. |
| Architecture | `reasoning/architecture.md` | Add `internal/runcontrol`, use direct same-host multi-process SQLite writers with leases and fencing, keep the accepting process as worker owner, do not adopt orphaned work, poll durable events, and fail closed or visibly degrade on persistence failure. | App/web/runtime seams, existing atomic state quality, module/DI/concurrency evidence, and explicit cross-process failure analysis. | Fixes package placement, store, topology, identity, timings, retention/quota, reconciliation, migration, and failure policy. |
| Frontend | `reasoning/frontend.md` | Add server-rendered `/runs` list/detail pages with narrow JavaScript replay enhancement, separate lifecycle from liveness, preserve no-JS snapshots, bound timeline state, and display explicit gap, tombstone, uncertainty, and cancellation states. | Existing embedded UI/SSE client, API decision, terminal UX, accessibility, security, and bounded-performance evidence. | Fixes browser/TUI projection behavior, polling/reconnect model, responsive/accessibility requirements, and cross-surface vocabulary. |

The area conclusions are adopted without architectural override. Where an area intentionally deferred a mechanism, the architecture area or the final decisions below supply it.

## Sprint Technical Handbook Summary

- **Relevant Patterns:** Thin adapters over a shared capability; explicit constructor composition; durable identity distinct from cancellation context; refreshable ownership; bounded structured concurrency; persistence before presentation; disk-backed state; typed failures; stable correlation fields; redaction at the trust boundary; and unit, command-path, process, and browser verification.
- **Important Trade-Offs:** Direct durable state costs write latency and migration complexity; at-least-once replay costs client deduplication; explicit lifecycle/liveness vocabulary costs UI complexity; bounded retention loses detail; and centralized composition must avoid becoming a god object.
- **Warnings / Anti-Patterns:** Surface-owned authority, hidden globals, context as identity storage, unjoined goroutines, indefinite waits, silent drops, string-only errors, partial IO abstraction, late redaction, unbounded queues/journals, and speculative plugin or workflow machinery.
- **Evidence Confidence:** High for package direction, explicit composition, concurrency bounds, typed errors, redaction, and layered testing. Medium for exact storage, lease/fencing, replay, terminal arbitration, and telemetry export because the studied repositories did not implement UltraPlan's complete cross-process contract; those values are sprint-specific decisions grounded in repository failure modes.

## Contracts Applied

| Contract / Requirement ID | Constraint | Decision Impact | Expected Evidence |
| --- | --- | --- | --- |
| Architecture; RO-1 through RO-7; AC-10 | Durable run authority stays outside adapters and product workflow authority stays in sprint/study modules. | Introduce focused `internal/runcontrol`; compose it in `internal/app`; correlate rather than replace product and agentwrap state. | Architecture review, import-boundary tests, and authority trace from acceptance to product result. |
| CLI Surface; RO-3, RO-4, RO-7; AC-2, AC-3, AC-7 | CLI/JSON/TUI/web expose the same identity, lifecycle, list, detail, follow, cancel, diagnostics, and recovery facts. | Add `ultraplan run list/show/follow/cancel/diagnostics` over shared app capabilities. | Command help, text/JSON fixtures, and cross-surface agreement tests. |
| Configuration; RO-2, RO-6, RO-9; AC-12 | Settings use existing precedence, validation, redaction, and safe limits. | Only history duration, tombstone duration, and quota are configurable; lease/event safety bounds remain fixed. | Config precedence, invalid-combination, migration, and redaction tests. |
| Documentation; RO-11 | Architecture, API, local-web, operations, recovery, retention, migration, and topology limits are explicit. | Documentation is part of release work, not follow-up polish. | Updated governed docs and user/operator references reviewed against this reasoning. |
| Errors; RO-4 through RO-9; AC-5, AC-6, AC-12 | Failures preserve causes internally and expose stable safe codes and recovery data. | Define typed store, cursor, tombstone, owner, cancellation, migration, and reconciliation errors. | Exact JSON/error fixtures and human guidance tests. |
| LLM Evaluation / Cost / Safety; RO-1, RO-2, RO-9; AC-11 | Unknown usage remains unknown; unsafe provider/prompt/output content is omitted before persistence. | Persist bounded safe correlation and omission facts, never raw payloads by default. | Adversarial redaction and known/unknown usage tests. |
| LLM Runtime; RO-1, RO-6, RO-7; AC-1, AC-8, AC-9 | Agentwrap retains runtime/process supervision and canonical event ownership. | Run control supplies product correlation and durable operational projection without duplicating provider execution. | Runtime request/event correlation, cancellation propagation, and no-child-on-acceptance-failure tests. |
| Observability; RO-1 through RO-9 | Logs, health, diagnostics, metrics, and support export use one safe correlation vocabulary. | Add run-control health and CLI diagnostics; defer OpenTelemetry export. | Structured-log fixtures, health scenarios, metrics checks, and redacted support bundle. |
| Performance; RO-2, RO-3, RO-5, RO-9; AC-12 | Storage, replay, polling, subscribers, compaction, reconciliation, and disk use are bounded. | Fix per-event/per-run/replay/DOM limits, indexed polling, quotas, headroom, and bounded reconciliation. | Load/boundary tests and measured multi-observer behavior. |
| Persistence And Migrations; RO-1, RO-2, RO-8; AC-1, AC-5, AC-6, AC-12 | Required records commit atomically before exposure; schema evolution and rollback are explicit. | Use transactional SQLite, durable token-digest dedupe, migrations/backups, tombstones, and no in-memory fallback. | Crash-boundary, corruption, disk-full, migration, restore, and legacy URL fixtures. |
| Security; RO-1, RO-4, RO-7, RO-9; AC-11 | Loopback remains the trust boundary; reads and mutations are separate; process identity and redaction are conservative. | Workspace-readable sanitized history, fresh CSRF/session/local authority for cancel, private database permissions, and no PID-only signaling. | Session rotation, CSRF, Host/Origin, PID-reuse, path, and secret-exposure tests. |
| Testing; RO-10; AC-2 through AC-13 | Deterministic, race, process, browser, build, and gated real-runtime evidence must cover the full failure matrix. | Inject clocks/process probes/failures and run real SQLite multi-process tests. | `go test ./...`, `go test -race ./...`, CLI build, browser suite, and gated dogfood. |
| Workflows; RO-1, RO-6, RO-7, RO-8; AC-8 through AC-10 | Cancellation, reconciliation, and retries cannot overwrite product state or the immutable operational terminal winner. | Use fenced attempts, durable cancellation facts, no worker adoption, and conservative interruption/uncertainty. | Duplicate-cancel, stale-writer, owner-death, terminal-race, and product-authority tests. |

## Repos Studied / Source Evidence Used

| Source / Repo / Report | Concrete Reference | Relevant Finding | Why It Matters For This Sprint | Used In Decision(s) |
| --- | --- | --- | --- | --- |
| gdu / restic; project structure | `01-project-structure.md`; `gdu/cmd/gdu/app/app.go:30-49`; `restic/internal/restic/repository.go:18` | Thin UI boundaries and inward domain interfaces support multiple surfaces without cross-imports. | Supports a shared run capability outside CLI, TUI, and web. | Package and authority; surface projection. |
| gh-cli / rclone / helm; command architecture | `02-command-architecture.md`; `gh-cli/pkg/cmdutil/factory.go:16-43`; `rclone/cmd/cmd.go:240-340`; `helm/pkg/cmd/install.go:132-145` | Factories and execution wrappers centralize lifecycle while commands delegate. | Supports app-level acceptance and stable CLI adapters. | Acceptance flow; CLI/API integration. |
| opencode / gh-cli; dependency injection | `03-dependency-injection.md`; `opencode/internal/app/app.go:42-81`; `gh-cli/pkg/cmd/factory/default.go:26-46` | Explicit composition and small interfaces make dependencies traceable and testable. | Supports constructor injection for repository, clock, process probe, notifier, and logger. | Package and authority; verification seams. |
| go-task / opencode / lazygit; configuration | `04-configuration-management.md`; `go-task/internal/flags/flags.go:314-327`; `opencode/internal/config/config.go:609-641`; `lazygit/pkg/config/app_config.go:256-330` | Precedence, post-merge validation, and explicit migration avoid ambiguous operational settings. | Supports a small validated run-control config surface and schema migration. | Retention/quota; migration. |
| rclone / restic / k9s; error handling | `05-error-handling.md`; `rclone/fs/fserrors/error.go:22-192`; `restic/internal/errors/fatal.go:10-53`; `k9s/internal/model/flash.go:100-103` | Typed/behavioral errors and separate user/operator reporting preserve recovery meaning. | Supports typed gaps, store failures, stale ownership, and cancellation uncertainty. | API/error contract; failure policy. |
| gh-cli / restic; IO abstraction | `06-io-abstraction.md`; `gh-cli/pkg/iostreams/iostreams.go:551-568`; `restic/internal/backend/backend.go:19-90`; `restic/internal/backend/mock/backend.go:14-26` | Injectable IO/backends permit deterministic failure testing; partial abstraction leaves blind spots. | Supports repository contracts, fake clocks/process identity, and storage fault injection. | Store; verification. |
| helm / opencode / restic; state and context | `07-state-context.md`; `helm/pkg/cmd/install.go:333-347`; `opencode/internal/session/session.go:12-23`; `restic/internal/restic/lock.go:105,290-305` | Cancellation context, durable identity, renewable locks, and cleanup contexts solve different problems. | Supports leases, fenced identity, cleanup retries, and no context-based authority. | Ownership/liveness; cancellation; terminal failure. |
| go-task / opencode / k9s; concurrency | `08-concurrency.md`; `go-task/task.go:87`; `opencode/cmd/root.go:261-279`; `k9s/internal/pool.go:21-48` | Goroutines, queues, fan-out, cancellation, waiting, and cleanup must be localized and bounded. | Supports owner loops, polling, bounded subscribers/replay, and finite shutdown. | Delivery; cancellation; reconciliation. |
| lazygit / restic / chezmoi; terminal UX | `09-terminal-ux.md`; `lazygit/pkg/tasks/tasks.go:31-435`; `restic/internal/ui/termstatus/status.go:197-205` | Long operations need calm progressive feedback, interruptibility, and non-TTY behavior without transferring ownership to presentation. | Supports shared lifecycle vocabulary, no-JS snapshots, and explicit cancellation state. | Frontend/TUI/CLI projection. |
| k9s / rclone / helm; observability | `10-logging-observability.md`; `k9s/internal/slogs/keys.go:6-231`; `rclone/fs/accounting/prometheus.go:78-108`; `helm/internal/logging/logging.go:31-71` | Stable structured fields and component correlation are proven; rich metrics and OpenTelemetry are uncommon. | Supports correlated logs/local diagnostics and deferring OTEL. | Telemetry and support. |
| chezmoi / gh-cli / restic / lazygit; testing | `11-testing-strategy.md`; `chezmoi/internal/cmd/main_test.go:64-174`; `gh-cli/acceptance/acceptance_test.go:26-29`; `lazygit/pkg/integration/clients/go_test.go:81` | Unit seams, real command paths, fixtures, and browser/PTY integration catch different failures. | Supports the layered race, process, compatibility, and browser matrix. | Verification and release. |
| helm / gh-cli; extensibility | `12-extensibility.md`; `helm/internal/plugin/metadata_v1.go:24-48`; `helm/internal/plugin/runtime_subprocess.go:65-79` | Versioned metadata helps compatibility; generalized plugins and registries carry lifecycle cost. | Supports versioned run schemas while rejecting daemon/plugin/remote-worker machinery. | Package/store; migration; topology. |
| restic / helm / opencode; security | `13-security.md`; `restic/internal/options/secret_string.go:15-20`; `helm/pkg/registry/transport.go:37-41`; `opencode/internal/permission/permission.go:44-108` | Secret-aware data, early scrubbing, validation, and permission separation enforce trust boundaries. | Supports pre-append redaction, private storage, and separate read/mutation policy. | Security/redaction; API authorization. |
| opencode / lazygit / rclone; performance | `14-performance.md`; `opencode/internal/message/message.go:37-42`; `opencode/internal/db/connect.go:39-54`; `rclone/lib/pool/pool.go:17-24,52-53` | Disk-backed long-session state, streaming, and bounded pools prevent memory growth. | Supports SQLite, bounded retention, polling/replay, and compaction, though the transaction requirements decide the exact store. | Store; delivery; retention. |

## Trade-Off And Debt Analysis

### Accepted Trade-Offs

| Trade-Off | Benefit | Cost / Constraint Accepted | Why Acceptable Now | Revisit Trigger |
| --- | --- | --- | --- | --- |
| Same-host local-filesystem topology | Defensible clock, process-birth, loopback, locking, and signaling semantics. | No shared-filesystem multi-host live-follow/cancel guarantee. | Matches the product security boundary and acceptance scenarios. | A governed remote-worker or multi-host requirement with identity, transport, clock, and authorization design. |
| SQLite with direct fenced writers | Transactions, unique sequence allocation, CAS terminal state, bounded queries, and multi-process access in one primitive. | New driver, binary size, schema/migration, WAL, and corruption operations. | Rebuilding these semantics over files is greater risk and complexity. | Measurements show unsupported filesystems, unacceptable build/runtime cost, or a daemon becomes a separately justified product requirement. |
| No coordinator daemon | CLI remains self-contained and durable truth survives any observer process. | Every process opens the store; observers poll. | Meets local multi-process needs without hidden authority or lifecycle. | Polling or writer contention fails measured targets after indexed/bounded optimization. |
| Acceptor owns worker; no adoption | Preserves real Go context, product locks, and agentwrap process ownership. | Owner death interrupts work rather than transparently failing over. | Product modules already expose deliberate resume semantics; process adoption cannot be made truthful. | A product stage gains a formally transferable checkpoint and fenced adoption contract. |
| At-least-once replay | Durable reconnect across process/server replacement with simple sequence dedupe. | Duplicate network delivery is possible. | Exactly-once acknowledgements add state without user value. | A future external consumer proves dedupe insufficient. |
| Fixed lease safety timings | Small, testable liveness contract with no unsafe combinations. | Unusual machines cannot tune heartbeats this sprint. | Conservative grace/process checks avoid premature terminal outcomes. | Dogfood shows repeatable false stalls or slow detection under supported hosts. |
| Bounded history and tombstones | Predictable memory/disk use, stable links, and explicit gaps. | Historical detail expires or compacts. | Requirements forbid unbounded retention and require visible loss. | Support data shows defaults prevent common diagnosis or quota is routinely exceeded. |
| Fail closed and cancel on sustained persistence loss | Never presents unverifiable acceptance, events, or active work as durable truth. | Healthy provider work may be cancelled and terminal commit may remain uncertain. | Truthful recovery is more important than invisible execution. | A proven replicated/durable write path can preserve truth without stopping work. |
| Canonical runs plus operation compatibility projection | Rich durable model without breaking frozen v1 operation clients. | Two DTO/event mappings must be maintained. | Existing compatibility fixtures make in-place semantic expansion unsafe. | A deliberate `/api/v2` retirement plan removes old operation clients. |
| Server-rendered UI plus polling/SSE | Refresh/no-JS correctness and no client-side authority. | Brief active-count/live latency and manual enhancement code. | Fits the existing UI and same-workspace local topology. | Measured scale or interaction complexity justifies a new transport/client architecture. |

### Potential Technical Debt

| Debt / Shortcut | Why It Might Accrue | Current Mitigation | Owner / Follow-Up |
| --- | --- | --- | --- |
| `modernc.org/sqlite` dependency size and platform behavior | Pure-Go SQLite expands build and upgrade surface. | Pin the dependency; test Linux/macOS build, migration, locking, integrity, and cross-process behavior. | Release maintenance after Sprint 35 dogfood. |
| Poll-based notifications | More observers increase read load and up-to-one-second idle latency. | Indexed bounded queries, 250 ms catch-up polling, 1 second idle backoff, and subscriber limits. | Revisit only from measurements, not preference. |
| Platform-specific process-birth probes | Linux and macOS expose different identities; unsupported cases may remain uncertain. | Probe behind an interface; uncertainty never grants signaling or success inference. | Platform support follow-up if Windows or another OS becomes required. |
| Dual canonical/compatibility projections | State and event mappings can drift. | Shared mapping tables and exact compatibility fixtures across app/API/browser/TUI. | Remove only through versioned API deprecation. |
| Fixed lease/progress sampling values | Defaults may not fit every workload. | Fake-clock tests, conservative grace, omission metadata, diagnostics, and dogfood measurements. | Configurability review after operational evidence. |
| Operational/product outcome divergence after outages | Separate authorities may report different facts. | Detail displays both explicitly and provides recovery guidance; no inferred merge. | Product-specific reconciliation improvements, not run-control authority expansion. |
| Manual restore rollback | One-way migrations require stopping all processes and restoring backup plus binary. | Pre-migration backup, integrity checks, documented runbook, and restore fixtures. | A future migration framework if schema cadence warrants it. |

### Future Considerations

| Consideration | Deferred Until | Reason Deferred | What Should Be Preserved Now |
| --- | --- | --- | --- |
| Multi-host/remote workers | Explicit product and security phase | Requires remote identity, secure transport, distributed clocks/leases, and remote cancellation. | Opaque IDs, repository interface, and no host facts in public IDs. |
| Worker adoption | A stage proves transferable checkpoints and fenced ownership transfer | Current Go contexts, locks, and agentwrap children are not adoptable. | Attempt hierarchy and immutable historical runs. |
| Coordinator/broker notifications | Polling fails measured supported-topology needs | A daemon adds hidden authority, upgrades, startup, and failure modes. | Durable sequence cursor and notifier abstraction. |
| OpenTelemetry export | Stable local correlation is dogfooded and an exporter use case exists | Selected studies provide no implementation precedent; it is unnecessary for acceptance. | Stable field vocabulary without vendor-specific schema. |
| Authored product persistence | Later authority gate | Operational SQLite is not evidence for moving Markdown or product state. | Strict table/package ownership and no dual writes. |
| Rich frontend or global stream | Browser complexity or measured tab polling warrants it | Current server-rendered progressive enhancement satisfies the product. | Versioned run API, explicit terminal flag, cursor/gap contract, and no browser authority. |
| Retry as run mutation | A governed retry/replay contract is defined | Existing product resume/retry commands own workflow semantics. | Links and correlations between immutable runs and product checkpoints. |

## Decisions

The authoritative decisions are specified in full under **Final Decisions** below:

1. Establish a same-host `internal/runcontrol` authority shared through `internal/app`, with product modules and agentwrap retaining their existing authorities.
2. Use direct fenced multi-process access to a private workspace SQLite repository, without a daemon, broker, in-memory fallback, or authored-artifact storage.
3. Assign opaque durable run/attempt identities and commit acceptance plus the first owner claim before any asynchronous product or child execution begins.
4. Use fixed renewable leases, process-birth verification, fencing, no worker adoption, and conservative startup/periodic reconciliation.
5. Persist idempotent cancellation before owner routing and arbitrate every operational ending through one immutable terminal compare-and-set.
6. Redact and durably sequence bounded events before at-least-once replay/fan-out, with explicit sampling, backpressure, and replay gaps.
7. Enforce bounded history, deterministic compaction, tombstones, quotas with reserved headroom, and fail-visible persistence degradation.
8. Add canonical run CLI/API resources while retaining frozen operation routes as compatibility projections with fresh mutation authorization.
9. Add server-rendered workspace run list/detail pages and shared CLI/TUI/browser vocabulary without introducing client-side authority.
10. Release with correlated local telemetry, diagnostics, a redacted support bundle, complete fault/race/browser coverage, documentation, and gated real-runtime dogfood.

## Final Decisions

### Decision 1: Same-Host Run-Control Product Boundary

- **Decision:** Add `internal/runcontrol` as the single operational product capability shared through `internal/app`. Support any number of CLI, TUI, and loopback web processes on one host against one canonical workspace on a local filesystem. `internal/web` retains only HTTP/SSE/session delivery concerns; sprint/study retain workflow and artifact authority; `internal/platform/runtime` and agentwrap retain runtime supervision.
- **Rationale:** The current app operation seam is already surface-neutral, while the web hub is the observed process/session-local defect. A focused product package has cohesive lifecycle state without becoming a platform dumping ground or workflow engine.
- **Study / Source Grounding:** `technical-handbook.md` patterns "Shared capability behind thin surfaces" and "Explicit composition root"; `01-project-structure.md`, `02-command-architecture.md`, and `03-dependency-injection.md`; gdu UI, helm actions, gh-cli factory, and opencode app composition; repository evidence in `internal/app/operations.go`, `internal/app/operation_runner.go`, and `internal/web/operations.go` summarized by the selected architecture reasoning.
- **Trade-Offs Accepted:** One new product package and explicit constructor plumbing are accepted for clear authority. Same-host scope excludes live multi-host guarantees.
- **Technical Debt / Future Impact:** The package must remain operationally narrow. A future remote topology must add a separately secured coordination protocol rather than broadening this store implicitly.
- **Alternatives Rejected:** `internal/web` authority repeats the defect; `internal/platform` would absorb product lifecycle semantics; separate sprint/study repositories would fragment workspace truth; a generic scheduler/workflow engine is outside scope; a daemon introduces hidden lifecycle authority without need.
- **Contracts Satisfied:** Architecture, Workflows, Security, RO-1 through RO-7, AC-2, AC-3, AC-10.
- **Evidence Required:** Architecture-review dependency trace; import tests; CLI/TUI/two-web-server composition over one store; proof product state and agentwrap remain separate authorities.

### Decision 2: Direct Multi-Process SQLite Repository

- **Decision:** Store operational state in `.ultraplan/run-control.db` using `database/sql` and `modernc.org/sqlite`. Use WAL, `synchronous=FULL`, foreign keys, a 5-second busy timeout, short transactions, file mode `0600`, and directory mode `0700`. Core records are `runs`, `attempts`, immutable `events`, `operation_aliases`, and bounded `reconciliation_log` data. There is no coordinator, broker, in-memory fallback, or dual-write sidecar.
- **Rationale:** Acceptance atomicity, per-run sequence allocation, fenced compare-and-set transitions, concurrent readers/writers, retention queries, and migration are database semantics. A snapshot/journal design would need to recreate transactions, indexes, repair, and compaction locking.
- **Study / Source Grounding:** `technical-handbook.md` patterns "Append or persist before presentation" and "Bounded retention with visible loss"; opencode disk-backed message/DB examples from `14-performance.md`; injectable backend evidence from `06-io-abstraction.md`; versioned metadata evidence from `12-extensibility.md`. The studies do not prove SQLite specifically; the store choice is grounded in UltraPlan's transaction and same-host multi-writer requirements plus the architecture area's failure analysis.
- **Trade-Offs Accepted:** Driver/binary size, migration, WAL, lock contention, and corruption recovery complexity are accepted to avoid ad hoc database reconstruction.
- **Technical Debt / Future Impact:** Pin and test the pure-Go driver across supported builds. Reject network filesystems that cannot prove local locking semantics. This database is never authority for authored product content.
- **Alternatives Rejected:** Filesystem snapshots plus append segments lack simple atomic multi-writer sequence/CAS/compaction; agentwrap stores do not own UltraPlan product-run lifecycle/cancellation; a daemon exceeds scope; Postgres and brokers violate local single-binary goals.
- **Contracts Satisfied:** Architecture, Persistence And Migrations, Performance, Security, RO-1, RO-2, RO-6, RO-8, AC-1, AC-6, AC-12.
- **Evidence Required:** Real temporary SQLite contract tests; concurrent process acceptance/append/CAS; crash-boundary tests; WAL/locking health; permissions; corruption, busy, permission, disk-full, backup, migration, and restore fixtures; CLI build on Linux/macOS.

### Decision 3: Stable Identity, Acceptance, and Authority Hierarchy

- **Decision:** Generate opaque random `run_<base32-128-bit>` and `att_<base32-128-bit>` IDs. A run is one accepted asynchronous operation; an attempt is one fenced owner claim. Preserve operation kind and project/sprint/study/stage/task target, product run/task references, agentwrap run/session/attempt IDs, owner process-lifetime ID, fencing generation, host/boot/PID/process-birth identity, and external harness run ID as typed correlations. IDs reveal no path, PID, time, provider, or authority. Every asynchronous operation through the shared runner receives a run; runtime-backed execution fails closed unless run acceptance and first owner claim are durable before any goroutine, product operation, process, or agentwrap child starts.
- **Rationale:** Product, runtime, provider, process, and browser identifiers have narrower lifetimes and cannot serve as workspace run authority. Acceptance ordering repairs the gap where the web currently advertises an in-memory operation before durable truth exists.
- **Study / Source Grounding:** `technical-handbook.md` pattern "Lifecycle context plus durable identity"; `07-state-context.md` session/context distinction; `03-dependency-injection.md` explicit service seams; `internal/app/operations.go` preparation/fingerprint boundary and `internal/platform/runtime/runtime.go` safe correlation fields as summarized by architecture/API reasoning.
- **Trade-Offs Accepted:** Clients cannot sort or interpret opaque IDs, and correlation fields add schema width. Product resume or explicit retry creates a new user-visible run; nested agentwrap policy attempts remain correlations.
- **Technical Debt / Future Impact:** A future lineage view may add explicit parent/retry links, never reinterpret existing IDs or mutate historical terminal runs.
- **Alternatives Rejected:** Reusing `op_*`, agentwrap run, provider session, study run, task, PID, or timestamp-derived IDs would leak or conflate authority; caller-chosen IDs invite collision and authorization mistakes.
- **Contracts Satisfied:** Architecture, LLM Runtime, Observability, Security, RO-1, AC-1, AC-3, AC-10.
- **Evidence Required:** ID format/collision tests; confirmation dedupe across response loss/second server; proof no child starts before durable claim; correlation fixtures that never fabricate missing values; product/runtime authority review.

### Decision 4: Fenced Ownership, Lease, and Conservative Reconciliation

- **Decision:** The accepting process owns the worker and work is never adopted after owner death. Use a 1-second owner control tick, 5-second durable heartbeat, 15-second lease, 10-second periodic reconciliation scan, and 45-second post-expiry terminal grace. Every owner mutation checks run ID, attempt ID, owner ID, and repository-allocated fencing generation. Lease comparisons use SQLite time; process checks use host digest, boot identity, PID, and exact process-birth token. Startup runs one bounded reconciliation pass before new runtime acceptance; web/TUI and active CLI owners repeat it periodically. Lease expiry classifies liveness but never proves terminal state. Missing/mismatched owner after grace may become `interrupted` or `cleanup_uncertain`, never success.
- **Rationale:** PID-only locks permit PID reuse, lease-only decisions are vulnerable to stalls/clock jumps, and another process cannot safely reconstruct a Go context, product lock, or agentwrap child. Fencing prevents stale writers even when process observation is uncertain.
- **Study / Source Grounding:** `technical-handbook.md` restic lock refresh and cleanup-context examples, concurrency bounds, and explicit evidence gap for PID reuse/fencing; `07-state-context.md` and `08-concurrency.md`; current sprint/study lock and reconciliation behavior identified in `code-context.md` and synthesized by architecture reasoning.
- **Trade-Offs Accepted:** Healthy blocked owners may appear stalled and owner death causes interruption rather than seamless failover. Fixed timings are intentionally not configurable this sprint.
- **Technical Debt / Future Impact:** Process-birth probes are platform-specific. Unknown identity remains uncertain and cannot authorize a signal. Adoption requires a future stage-specific transferable-checkpoint contract.
- **Alternatives Rejected:** PID-only checks are unsafe; lease expiry as immediate death is not conservative; leader election is unnecessary with transactional idempotence; worker adoption risks duplicate mutation; process absence/artifact presence cannot infer success.
- **Contracts Satisfied:** Workflows, Security, Observability, RO-1, RO-6, AC-8, AC-9.
- **Evidence Required:** Fake-clock lifecycle tests; stale fence rejection for append/heartbeat/cancel/terminal; exact PID-reuse/process-birth mismatch tests; owner kill before/after each commit boundary; clock-jump and live-but-stalled scenarios; repeated reconciler idempotence.

### Decision 5: Durable Cancellation and One Operational Terminal Winner

- **Decision:** Cancellation is a durable command. An authorized request is transactionally recorded once; the current fenced owner polls, durably acknowledges, then cancels its context so existing product/runtime paths reach agentwrap cleanup. Duplicate requests return current state without a second signal. Completion, failure, timeout, cancellation, interruption, cleanup uncertainty, and shutdown/reconciliation proposals compete through one immutable terminal compare-and-set. Server shutdown drains new starts and requests `server_shutdown` cancellation only for workers that server owns. An unreachable owner leaves cancellation `uncertain`; no bare signal is sent to an unverified PID.
- **Rationale:** Persist-before-route survives caller/server loss and separates request, delivery, acknowledgement, cleanup, and terminal outcome. One transaction prevents local mutexes or repeated observers from rewriting races.
- **Study / Source Grounding:** `technical-handbook.md` lifecycle-context, typed-error, bounded-shutdown, and `sync.Once` findings from `07-state-context.md`, `08-concurrency.md`, and `05-error-handling.md`; agentwrap cancellation remains in `internal/platform/runtime`; existing web cancellation/shutdown races are summarized in code context and area reasoning.
- **Trade-Offs Accepted:** `202` means requested, not delivered or cancelled; completion can win after cancellation; unreachable cleanup remains explicit uncertainty.
- **Technical Debt / Future Impact:** Verified process-group cleanup may improve only through the runtime boundary. Retry remains a separate product-owned guarded operation.
- **Alternatives Rejected:** Browser/SSE disconnect cancellation conflates observation and command; process-local context registries fail cross-process; repeated signals risk the wrong process; cancellation overwriting terminal state destroys race truth.
- **Contracts Satisfied:** Workflows, LLM Runtime, Errors, Security, RO-6, RO-7, AC-8, AC-9.
- **Evidence Required:** First/duplicate/already-terminal cancel responses; fresh auth and CSRF rejection; stale/unreachable owner; acknowledgement; server shutdown; cancellation/completion/failure/timeout/reconciliation races under `go test -race`; no unsafe signal on PID reuse.

### Decision 6: Sanitized Ordered Journal, Replay, and Backpressure

- **Decision:** Redact, allowlist, and size-bound each event before SQLite append. Allocate a monotonic sequence within the run transaction and publish only after commit. Delivery is at-least-once; clients deduplicate `(run_id, sequence)`. Stable source identities may deduplicate producer callbacks; otherwise callbacks remain distinct events. Poll SQLite at 250 ms while catching up and back off to 1 second idle; in-process notification is only an optimization. Persist at most one equivalent high-frequency progress sample per run per 250 ms and record omitted count/time range. Cap encoded events at 16 KiB, historical replay at 512 events per SSE connection, transport queues, and browser/TUI rendering. Slow subscribers disconnect and resume from their durable cursor.
- **Rationale:** The repository, not SSE, must survive observer replacement. Commit-before-fan-out makes every visible sequence replayable, while bounded queues and coalescing prevent observers from affecting work or disk.
- **Study / Source Grounding:** `technical-handbook.md` append-before-presentation, bounded-loss, streaming, redaction, and concurrency patterns; lazygit streaming, opencode pub/sub/disk state, rclone pools, and security examples from `08-concurrency.md`, `13-security.md`, and `14-performance.md`; current web event bounds and agentwrap safe raw omission in code context.
- **Trade-Offs Accepted:** Up to one-second idle live latency, duplicate delivery, and sampled progress detail. Lifecycle, warning, finding, artifact, cancellation, recovery, and terminal events are never progress-coalesced.
- **Technical Debt / Future Impact:** Poll load and sampling values need dogfood measurement. A later notifier/broker must preserve SQLite sequence authority and explicit gaps.
- **Alternatives Rejected:** In-memory hub/SSE as truth disappears on restart; exactly-once delivery adds acknowledgement state; unbounded replay/queues threaten execution; raw provider persistence violates security; silent drops violate observability.
- **Contracts Satisfied:** Observability, Performance, Security, LLM Runtime, RO-2, RO-5, RO-9, AC-4, AC-6, AC-11, AC-12.
- **Evidence Required:** Concurrent monotonic append; commit-before-delivery; duplicate producer/network delivery; server-render/SSE race; two-server replay then live follow; 512-event catch-up reconnect; coalescing metadata; oversize replacement; slow subscriber; TUI refresh after channel drop; hostile content/redaction tests.

### Decision 7: Explicit Snapshot, Retention, Quota, and Failure Policy

- **Decision:** Preserve a durable current snapshot and event lower boundary independently of event detail. Default limits are 4,096 events and 16 MiB per run, 7 days full terminal history, 30 days terminal tombstones, 496 MiB workspace soft quota, 512 MiB hard quota, and 16 MiB reserved headroom. Full-history duration, tombstone duration, and quota use existing config precedence with minimums of 1 hour, 24 hours, 64 MiB hard quota, and 16 MiB headroom. Start compaction at 80 percent; fail new acceptance at soft quota; reserve headroom for active heartbeat/cancel/recovery/terminal writes; at hard quota cancel active owned work and expose degradation. Compact expired progress, then full records to tombstones, then expired tombstones; active runs are never deleted. Event append retries for 5 seconds before owner cancellation; missing heartbeat before lease expiry causes cancellation; terminal commit retries 30 seconds on cleanup context, then reconciliation determines interruption/uncertainty.
- **Rationale:** Retention is correctness, performance, and privacy. Reserved capacity increases the chance that active runs can record a truthful ending. Explicit boundaries and tombstones make loss visible rather than presenting partial history as complete.
- **Study / Source Grounding:** `technical-handbook.md` bounded retention, disk-backed state, error classification, and configuration trade-offs from `04-configuration-management.md`, `05-error-handling.md`, `13-security.md`, and `14-performance.md`. Exact values are Sprint 35 architecture decisions derived from the required failure matrix, not directly copied from study repositories.
- **Trade-Offs Accepted:** Detail expires; full disk exhaustion can still prevent terminal commit; persistence degradation deliberately stops expensive work. Only evidence-backed settings are configurable.
- **Technical Debt / Future Impact:** Quotas and durations may require tuning. Observable gap/tombstone semantics and fail-visible degradation remain stable even if values change.
- **Alternatives Rejected:** Unlimited retention risks disk failure/privacy; deleting active snapshots destroys recovery; accepting at quota consumes terminal capacity; in-memory continuation creates invisible work; silent append loss breaks replay truth.
- **Contracts Satisfied:** Performance, Persistence And Migrations, Configuration, Security, Errors, RO-2, RO-5, RO-9, AC-6, AC-11, AC-12.
- **Evidence Required:** Per-run bounds; deterministic compaction order; gap boundaries; active-run preservation; soft/hard quota/headroom; disk-full/permission loss at acceptance, append, heartbeat, cancel, and terminal; config validation; current snapshot under all compaction states.

### Decision 8: Additive Run API, Compatibility, Authorization, and Migration

- **Decision:** Add `GET /api/v1/runs`, `GET /api/v1/runs/{id}`, `GET /api/v1/runs/{id}/events`, and authorized idempotent `DELETE /api/v1/runs/{id}` plus `ultraplan run list/show/follow/cancel/diagnostics`. Keep `POST /api/v1/operations` as confirmed start; for new work its ID equals the durable run ID. Existing operation list/detail/events/cancel routes project the durable repository into frozen shapes and event names. New start returns `202` only after acceptance and adds a canonical `Link` header. Safe list/detail/events are workspace-readable across browser sessions; browser mutation requires a current same-origin session plus CSRF, and CLI/TUI require local workspace/OS authority. Use full, compacted, and tombstone snapshots; `409 replay_gap` includes cursor boundaries and current snapshot; recognized pre-durable `op_*` links without history return `410 legacy_operation_not_retained` and recovery guidance. Migrations use SQLite schema versions, a workspace migration lock, WAL checkpoint, bounded timestamped backup, transactional one-way migration, and integrity checks. Do not synthesize history from product artifacts. Rollback stops all processes and restores the matching binary and pre-migration database backup.
- **Rationale:** The existing operation DTO, route matrix, lifecycle, and SSE event names are compatibility fixtures. A canonical resource can express attempts, liveness, retention, and tombstones without silently changing old clients.
- **Study / Source Grounding:** API reasoning's direct inspection of `api_compatibility_test.go`, `operations_contract_test.go`, handlers/routes/client; `technical-handbook.md` typed-error, security, state/context, extensibility, and testing findings from `05`, `07`, `11`, `12`, and `13` reports.
- **Trade-Offs Accepted:** Two projections, more error codes, and wider sanitized read visibility inside one loopback workspace. EventSource needs detail preflight/fallback because it hides non-200 bodies.
- **Technical Debt / Future Impact:** Required incompatible API evolution uses `/api/v2`. Operation compatibility remains until an explicit deprecation. Retained dedupe/tombstone expiry bounds retry guarantees.
- **Alternatives Rejected:** Expanding frozen `operationDocument` risks strict clients; originating-session reads repeat the defect; generic `404` hides known retention/legacy states; fabricated imports make product files shadow run history; bearer accounts/tenancy exceed local scope.
- **Contracts Satisfied:** CLI Surface, Errors, Security, Persistence And Migrations, RO-3 through RO-8, AC-2 through AC-7, AC-9.
- **Evidence Required:** Route/method and exact old-schema fixtures; strict old-client/browser-bundle tests; additive run schemas; session rotation; CSRF/Origin; dedupe through a second server; cursor-ahead/gap/expiry; compacted/tombstone/legacy states; unsupported schema, backup, restore, and rollback documentation tests.

### Decision 9: Shared Cross-Surface Presentation

- **Decision:** Add server-rendered `/runs` and `/runs/{run_id}` pages, redirect/resolve `/operations/{id}`, and use dependency-free narrow JavaScript. The top bar reads workspace active lifecycle and refreshes on load/focus and every 5 seconds while visible; failure shows unavailable/stale, never zero. Detail renders snapshot and at most 200 events before opening SSE after the rendered sequence; the DOM retains at most 500 rows. Display lifecycle separately from liveness, operational result separately from product status, and explicit connecting/live/reconnecting/gap/capacity/store/tombstone/legacy/cancellation-uncertain states. Add URL-backed list filters, mobile cards below 720 px, keyboard/focus handling, throttled live regions, non-color state cues, reduced-motion behavior, and no-JS navigation/forms. TUI consumes the same app snapshots/events and refreshes durably after local delivery drops.
- **Rationale:** Server-first state works across refresh, session rotation, observer restart, and no JavaScript. Narrow cursor enhancement closes the render/stream race without creating a browser store.
- **Study / Source Grounding:** Frontend reasoning; `technical-handbook.md` terminal UX, error, concurrency, security, performance, and testing findings from reports `05`, `08`, `09`, `11`, `13`, and `14`; current `app.js` and operation handler behavior summarized in selected code context.
- **Trade-Offs Accepted:** Polling introduces brief staleness and richer truth states add user vocabulary. Bounded DOM differs from durable retention and must be labeled separately.
- **Technical Debt / Future Impact:** Mapping drift is controlled with shared fixtures. A frontend framework/global stream remains deferred until measured complexity requires it.
- **Alternatives Rejected:** Project-page-only lists miss workspace activity; SPA/global store creates competing authority; empty SSE-first pages fail on transport loss; collapsing stalled into failed/running lies; auto-scroll and per-event live announcements harm accessibility.
- **Contracts Satisfied:** CLI Surface, Documentation, Errors, Security, Performance, RO-3 through RO-5, AC-2 through AC-7.
- **Evidence Required:** Server/no-JS rendering; active count for two CLI runs/two servers; responsive/mobile, zoom, keyboard, focus, live-region, reduced-motion, unknown-state, hostile-content, bounded-DOM, hidden-tab polling, reconnect/gap, terminal refresh, and CLI/TUI/browser label agreement tests.

### Decision 10: Local Telemetry, Support, and Release Verification

- **Decision:** Use stable safe correlation fields for request, run, attempt, operation/target, owner/fence, process/runtime/agentwrap, sequence, lifecycle transition, cancellation, reconciliation, persistence operation, and terminal winner. Add an additive `run_control` health summary; local counters/histograms for acceptance/append/terminal latency/failure, active/stalled runs, leases, backlog, retention/compaction, replay/subscriber loss, and cancellation routing; `ultraplan run diagnostics --json`; and an explicit bounded redacted support bundle. IDs are log fields, never metric labels. Defer OpenTelemetry and a broad Prometheus endpoint. Release requires deterministic unit/contract/fault tests, real multi-process SQLite tests, race tests, browser accessibility/integration tests, CLI build, documentation, architecture/sprint review, and gated real-runtime dogfood.
- **Rationale:** Operators need one explanation path for persistence, liveness, replay, cancellation, and recovery. Structured local evidence is supported and sufficient; exporter infrastructure is not.
- **Study / Source Grounding:** `technical-handbook.md` correlation, dual-channel diagnostics, and layered verification patterns; k9s/rclone/helm observability examples in `10-logging-observability.md`; command-path/mock/browser evidence in `11-testing-strategy.md`; redaction evidence in `13-security.md`. No selected study source demonstrates an OpenTelemetry implementation, which is why export is explicitly deferred.
- **Trade-Offs Accepted:** No scrape-ready remote metrics or distributed traces; high-cardinality diagnosis stays in logs/support data. Verification cost is high because correctness spans processes and commit boundaries.
- **Technical Debt / Future Impact:** Stable correlation vocabulary permits a later optional exporter. Support bundle allowlists and size limits require ongoing security review as schemas grow.
- **Alternatives Rejected:** Print debugging and string-only errors are not correlatable; run IDs as metric labels create cardinality/privacy risk; broad unauthenticated metrics can leak metadata; provider-backed normal tests are non-deterministic; happy-path-only browser tests miss the sprint defect.
- **Contracts Satisfied:** Observability, Testing, Documentation, LLM Evaluation / Cost / Safety, RO-9 through RO-11, AC-11 through AC-13.
- **Evidence Required:** Structured-log and metric fixtures; health degraded scenarios; support-bundle size/redaction; full failure matrix; `go test ./...`; `go test -race ./...`; `go build ./cmd/ultraplan`; browser suite; Architecture Review; Sprint Review; Deep Smoke Sprint with two local servers and a real CLI runtime run or a truthful blocked result when prerequisites are unavailable.

## Expected Evidence

| Evidence Type | Required Evidence | Source / Command / Review Check |
| --- | --- | --- |
| Unit and contract tests | Identity validation, lifecycle invariants, terminal CAS, lease/fence, sequence allocation, redaction, retention, error codes, DTO/event compatibility, config, migrations, and surface mappings. | `go test ./...`; focused package tests under `internal/runcontrol`, `internal/app`, `internal/web`, and `internal/tui`. |
| Race tests | Concurrent writers, owner/reconciler/cancel/terminal races, slow subscribers, and shutdown without data races or stale-writer success. | `go test -race ./...`. |
| Process integration | Two concurrent CLI owners, two local web observers, restart/session rotation, owner kills at commit boundaries, PID reuse, polling/replay, quota/disk/corruption, migration/restore. | Process harness over temporary workspaces and real SQLite; no provider required. |
| Browser and accessibility | Server-first list/detail, active count, cursor handoff, catch-up, typed gap, cancellation, responsive layout, keyboard/focus/live-region/reduced-motion, hostile content, and legacy/tombstone recovery. | Existing browser test command documented by the implementation; `httptest` plus browser engine. |
| Runtime | One gated real runtime accepted through CLI, observed/replayed from two local servers, recovered after observer change, and cancelled or completed with correlated product/agentwrap evidence. | Deep Smoke Sprint protocol and external harness evidence; unavailable prerequisites yield `blocked`, not pass. |
| Build | Single CLI binary with pure-Go SQLite dependency and embedded UI. | `go build ./cmd/ultraplan` and release-platform build checks. |
| Operations | Structured logs, run-control health, local metrics, diagnostics JSON, reconciliation record, quota/compaction facts, and bounded redacted support bundle. | `ultraplan run diagnostics --json`, health API/CLI, support export, and failure-injection fixtures. |
| Review | Authority trace, stale-writer proof, failure matrix, compatibility, safety, product-state separation, and selected topology documented precisely. | `system/protocols/architecture-review-protocol.md` and `system/protocols/review-sprint-protocol.md`. |
| Documentation | Architecture/PRD/TRD alignment plus CLI/API/local-web/operations/recovery/retention/migration/rollback/topology/security guidance. | Governed project docs and implementation user/operator documentation reviewed against Decisions 1-10. |

## Assumptions And Risks

| Item | Type | Impact | Mitigation / Follow-Up |
| --- | --- | --- | --- |
| Supported workspaces use a local filesystem with correct SQLite locks/WAL behavior. | Assumption | Unsupported mounts can violate concurrency/durability. | Health/preflight rejects uncertain placement; document same-host/local-filesystem contract. |
| `modernc.org/sqlite` meets build and FULL-sync semantics on supported Linux/macOS releases. | Assumption | Driver defects could block acceptance or release builds. | Pin version; contract, integrity, cross-process, and platform build tests. |
| Existing app operation preparation remains the authoritative fingerprint/confirmation seam. | Assumption | Bypasses could start unrecorded work. | Inventory every asynchronous CLI/TUI/web entry and make shared acceptance unavoidable. |
| Agentwrap cancellation and safe canonical event projection remain available. | Assumption | Cleanup or correlation could be incomplete. | Correlate available facts; classify cleanup uncertainty; never duplicate provider supervision. |
| Process-birth identity can be verified on supported hosts. | Risk | Unknown identity can delay reconciliation/cancellation. | Conservative uncertainty; no PID-only signal; platform-specific tested probes. |
| SQLite polling scales to expected local tabs/processes. | Risk | Read contention or latency could degrade live follow. | Indexed bounded reads, idle backoff, subscriber limits, metrics, and dogfood thresholds. |
| Total filesystem exhaustion can consume reserved headroom. | Risk | Terminal persistence may fail despite quota policy. | Cleanup-context retries, health degradation, owner cancellation, conservative later reconciliation, and disk-recovery runbook. |
| Mid-run persistence loss wastes runtime work. | Risk | Cost and partial artifact creation increase. | Five-second retry, explicit diagnostics, product-state separation, and deliberate resume commands. |
| Wider workspace-readable metadata is sensitive in some local environments. | Risk | Another local browser session may see targets/timing. | Loopback/same-origin policy, field allowlist, workspace-relative display, pre-storage redaction, no credentials/raw paths. |
| Rich canonical lifecycle can drift from operation compatibility values. | Risk | Old clients may misclassify runs. | Conservative mapping table, explicit terminal flag, strict old-client and embedded-browser fixtures. |
| Tombstone and dedupe expiry limit indefinite old-link/start-retry behavior. | Risk | Very old callers receive recovery rather than detail/idempotent replay. | Explicit `410`, retention documentation, workspace list/product links, and no fabricated history. |
| Operational and product terminal facts can disagree after faults. | Risk | Users may mistake operational success/failure for artifact truth. | Present both authorities separately with typed recovery; run control never rewrites product state. |
| Fixed timings/defaults may create false stalls or insufficient history. | Risk | Operator confusion or diagnostic loss. | Conservative grace, injected-time tests, metrics, omission facts, and evidence-gated tuning after dogfood. |
| Downstream plan could treat implementation order as architecture freedom. | Risk | Scope or authority could drift. | Plan must execute Decisions 1-10 and may sequence work, but cannot reopen package/store/topology/API/lifecycle choices. |

## Implementation Constraints

- `internal/runcontrol` must not import `internal/web`, `internal/tui`, `internal/sprint`, `internal/study`, or runtime adapter packages.
- `internal/web` and `internal/tui` must consume app capabilities and cannot persist, arbitrate, or infer durable run truth.
- Product modules remain authoritative for workflow locks, artifacts, validation, retries/resume, and product terminal outcomes.
- Agentwrap remains the runtime/process supervision boundary; UltraPlan records safe correlations and canonical projections only.
- Required acceptance, event, cancellation, lease, and terminal mutations must be transactional and durable before exposure.
- No asynchronous runtime-backed path may bypass run acceptance or start a child after acceptance persistence fails.
- Every authoritative owner write must include run, attempt, owner, and fencing generation checks.
- Lease expiry, PID absence, artifact presence, and planning-stage status cannot individually prove success or terminal state.
- Redaction and size bounds apply before database writes, logs, SSE, TUI delivery, JSON, HTML, metrics diagnostics, and support export.
- Event transport is at-least-once observation; commands remain explicit HTTP/app operations and browser disconnect never cancels work.
- Active projection comes only from canonical run lifecycle; liveness and product status remain separate projections.
- Existing v1 operation DTO field order, route/method matrix, lifecycle/event compatibility, and error envelope remain fixture-controlled.
- Database migrations are one-way and transactional with backup/integrity checks; rollback is stop-and-restore, never dual write.
- Active runs and their current snapshots cannot be removed by retention; every advanced replay boundary is explicit.
- Safety timings, event-size limit, and per-run event bounds are fixed this sprint; only documented retention/quota settings are configurable.
- No daemon, broker, remote-worker protocol, frontend framework, WebSocket layer, OpenTelemetry exporter, broad metrics endpoint, generic virtual filesystem, or authored-product database enters this sprint.
- Normal tests use fakes and local processes; real providers are gated and cannot be required by `go test ./...` or `go test -race ./...`.
- Implementation must preserve unrelated workspace and Git state and must not add automatic Git mutation.

## Plan Handoff

`plan.md` must execute these decisions. It must not invent architecture, scope, or decisions beyond this document.

The plan must carry forward:

- Decisions 1-10, including exact topology, package, SQLite policy, identity, timing, retention/quota, API, UI, telemetry, and failure values
- all selected contracts and `RO`/`AC` mappings
- migration, backup, restore, compatibility, rollback, and fault-injection work as first-class tasks
- explicit integration of every asynchronous CLI, TUI, and web entry point before declaring workspace-wide coverage
- the complete expected-evidence matrix and required review protocols
- assumptions and risks, especially filesystem/driver behavior, process birth, persistence loss, mapping drift, and separate product authority
- documentation and gated real-runtime dogfood before release completion

Implementation sequencing may begin with package/repository contracts and app extraction, then integrate owners/events/cancellation/reconciliation, compatibility/API/CLI/TUI/web, retention/telemetry/migration, and finally documentation and verification. Sequencing cannot create an interim release path that starts runtime work without durable acceptance.

## Phase Exit Criteria

- [x] Selected context was read and used.
- [x] Area-specific reasoning documents were completed and synthesized.
- [x] Area-specific reasoning conclusions are reflected without an unresolved override.
- [x] Contracts and requirement outcomes are mapped to decisions and evidence.
- [x] Topology, ownership, store, identity, lifecycle, lease, cancellation, replay, retention, failure, API, authorization, migration, telemetry, and UI choices are final.
- [x] At least two alternatives are explicitly rejected with rationale; all material architecture alternatives are resolved.
- [x] Accepted trade-offs, debt, future considerations, assumptions, and risks are explicit.
- [x] Expected evidence is specific, behavior-focused, and reviewable.
- [x] Final decisions are clear enough for `plan.md` to execute without reopening architecture.
