# Architecture: Durable Run Control Plane

> **Inputs Used:** `projects/ultraplan-go/sprints/35-durable-run-observability/technical-handbook.md`, `projects/ultraplan-go/sprints/35-durable-run-observability/requirements.md`, `projects/ultraplan-go/sprints/35-durable-run-observability/sprint-index.md`, `projects/ultraplan-go/sprints/35-durable-run-observability/code-context.md`, `projects/ultraplan-go/docs/ARCHITECTURE.md`, `projects/ultraplan-go/docs/PRD.md`, `projects/ultraplan-go/docs/TRD.md`, `studies/go-cli-study/reports/final/01-project-structure.md`, `studies/go-cli-study/reports/final/02-command-architecture.md`, `studies/go-cli-study/reports/final/03-dependency-injection.md`, `studies/go-cli-study/reports/final/04-configuration-management.md`, `studies/go-cli-study/reports/final/05-error-handling.md`, `studies/go-cli-study/reports/final/06-io-abstraction.md`, `studies/go-cli-study/reports/final/07-state-context.md`, `studies/go-cli-study/reports/final/08-concurrency.md`, `studies/go-cli-study/reports/final/10-logging-observability.md`, `studies/go-cli-study/reports/final/11-testing-strategy.md`, `studies/go-cli-study/reports/final/12-extensibility.md`, `studies/go-cli-study/reports/final/13-security.md`, `studies/go-cli-study/reports/final/14-performance.md`, `system/reasoning/architecture_reasoning_template.md`, `../ultraplan-go/internal/app/operations.go`, `../ultraplan-go/internal/app/operation_runner.go`, `../ultraplan-go/internal/web/operations.go`, `../ultraplan-go/internal/web/operation_handlers.go`, `../ultraplan-go/internal/web/server.go`, `../ultraplan-go/internal/platform/runtime/runtime.go`, `../ultraplan-go/internal/sprint/execute_state.go`

This area selects the ownership, topology, storage, lifecycle, liveness, cancellation, reconciliation, retention, and failure model for Sprint 35. The control plane records execution facts only. Sprint and study modules remain authoritative for workflow semantics and artifacts, and agentwrap remains authoritative for runtime supervision.

## Area Decisions

### 1. Supported topology

Support multiple UltraPlan processes on one host, using one canonical workspace on a local filesystem. CLI, TUI, and any number of loopback web servers may accept, own, inspect, follow, cancel, and reconcile runs through the same workspace repository.

Shared-filesystem multi-host operation is unsupported in this sprint. Network filesystems, independent host clocks, remote process namespaces, remote authorization, and remote signaling are outside the guarantee. Opening a copied workspace may provide historical reads after clean shutdown, but it does not provide live ownership or cancellation guarantees.

This boundary follows the existing loopback-only security contract and the selected reports' evidence, which contains local locks and sessions but no distributed coordination implementation.

### 2. Package and authority boundary

Create a focused `internal/runcontrol` product package. It owns:

- durable run and attempt identity
- acceptance, lifecycle, and one terminal compare-and-set
- owner claims, leases, heartbeats, fencing, and safe process identity
- sanitized ordered event append and replay boundaries
- workspace-wide run queries and active projections
- durable cancellation requests and owner acknowledgement
- conservative reconciliation, compaction, tombstones, health, and support diagnostics

It does not own stage ordering, artifact validation, task success, product retries, prompt construction, provider supervision, browser sessions, SSE sockets, TUI state, or HTTP DTOs.

`internal/app` composes `runcontrol.Service` with existing operation preparation and product services. CLI, TUI, and web adapters call the same app capability. `internal/web` keeps only request validation, response mapping, SSE connection state, and bounded subscriber buffers.

Dependency direction is:

```text
cmd/ultraplan -> internal/app
internal/tui  -> internal/app
internal/web  -> internal/app
internal/app  -> internal/runcontrol + product modules
product modules -> internal/platform/runtime
internal/runcontrol -> no web, TUI, sprint, study, or runtime adapter package
```

Use explicit constructor injection for the repository, clock, owner/process identity probe, logger, and notifier. Context carries cancellation and deadlines, not services or durable identity. Do not add globals, an init-time registry, or a general workflow engine.

### 3. Durable representation

Use one workspace-local SQLite database at `.ultraplan/run-control.db`, accessed through Go's `database/sql` and the pure-Go `modernc.org/sqlite` driver. This keeps the single-binary and non-CGO build contract while providing transactions, unique constraints, compare-and-set updates, ordered indexes, bounded queries, and direct multi-process concurrency.

Production connection policy is:

- `journal_mode=WAL`
- `synchronous=FULL`
- `foreign_keys=ON`
- `busy_timeout=5s`
- one short transaction per acceptance, claim, append, heartbeat, cancellation, terminal transition, or compaction batch
- database file mode `0600` and containing run-control directory mode `0700`

WAL files and locks are another reason network filesystems are unsupported. The store is operational authority only; it does not contain authored Markdown, flow state, execute state, study state, source files, Git state, or smoke evidence.

The initial schema contains these cohesive records:

| Record | Purpose |
| --- | --- |
| `runs` | Identity, target, lifecycle snapshot, sequence boundaries, cancellation summary, terminal winner, retention state, and schema metadata. |
| `attempts` | Attempt identity, owner/fence claim, timestamps, runtime/process correlations, heartbeat, lease, and attempt outcome. |
| `events` | Sanitized immutable events keyed by `(run_id, sequence)`. |
| `operation_aliases` | Durable mapping for compatibility operation IDs and retained legacy guidance. |
| `reconciliation_log` | Bounded safe evidence for stale-owner and recovery decisions. |

Cancellation facts remain columns/records under the run transaction rather than a general command queue. The owner polls only requests for attempts it currently owns.

### 4. Identity hierarchy

Use opaque cryptographically random IDs that clients never parse:

- `run_id`: one accepted user-visible asynchronous execution
- `attempt_id`: one fenced owner claim within that run
- product operation kind and target: project, sprint, study, stage, and task correlation
- product run/task IDs: references to existing sprint or study state
- agentwrap run/session/attempt IDs: runtime correlations, not run-control authority
- owner ID: random per UltraPlan process lifetime
- fencing generation: monotonically allocated by the repository per attempt claim
- process identity: host digest, boot identity, PID, and process-birth token
- external harness run ID: correlation only

A product resume or explicit retry is a new user-visible run unless it is an internal retry already owned by the same accepted operation. Agentwrap policy attempts remain nested runtime correlations. This avoids rewriting historical terminal runs.

### 5. Acceptance and execution flow

The required flow is:

1. The adapter prepares and confirms the normalized product operation using existing fingerprint checks.
2. App asks run control to create the run and initial accepted snapshot in one durable transaction.
3. If acceptance fails, app returns a typed storage error and starts no goroutine, product operation, process, or agentwrap child.
4. The accepting process claims the first attempt, allocating its fencing generation and lease before invoking the existing product operation.
5. Product services retain their locks and workflow authority; safe product progress is normalized through app into run-control events.
6. Each event is redacted, bounded, sequenced, committed, and only then offered to transient subscribers.
7. The owner refreshes its durable snapshot and lease while work runs.
8. Product completion proposes an operational terminal result. One fenced transaction wins; later completion, cancellation, shutdown, or reconciliation proposals read the winner.

The process that accepts a run owns its worker. Work is not adopted after owner death. Existing product resume commands create deliberate new work using product-owned checkpoints; run control does not attempt to reconstruct a live Go context or agentwrap process in another process.

### 6. Ownership, lease, and fencing

Use these fixed safety timings for Sprint 35:

- owner control-loop tick: 1 second
- durable heartbeat cadence: 5 seconds
- lease duration: 15 seconds
- reconciliation scan cadence in long-lived processes: 10 seconds
- terminal reconciliation grace after lease expiry: 45 seconds

Tests inject shorter values and a fake clock, but production timing is not user-configurable in this sprint.

Every authoritative mutation checks `run_id`, `attempt_id`, owner ID, and fencing generation. A new claim increments the generation transactionally. A stale writer cannot renew, append, acknowledge cancellation, or commit terminal state after its fence is superseded.

Lease comparisons use database time so same-host writers share one wall-clock source. Clock jumps may classify a run as stalled, but lease expiry alone never creates a terminal result. Reconciliation also checks exact process-birth identity. A live matching process remains stalled/suspect rather than being declared dead; a missing or mismatched birth token is never signaled as though it were the owner.

### 7. Cancellation coordination

Cancellation is a durable command, not an SSE event and not a browser-disconnect side effect.

1. An authorized caller transactionally records the request if none exists.
2. The current owner sees it on its 1-second control loop, acknowledges it under the current fence, and cancels the owning context.
3. Existing app/product/runtime paths propagate cancellation to agentwrap and process cleanup.
4. Cancellation competes with completion, failure, timeout, interruption, shutdown, and reconciliation through the one terminal compare-and-set.

Repeated requests return the current durable state and do not emit repeated owner signals. If the owner is unreachable, the request remains visible as uncertain until reconciliation. Run control does not send a bare signal to an unverified PID. If agentwrap exposes a verified process-group identity, cleanup may use it through the existing runtime boundary; otherwise orphaned runtime cleanup is recorded as uncertain rather than guessed.

### 8. Event delivery and cross-process observation

The SQLite journal is the source of truth. In-process notification is only a latency optimization.

Web servers replay by sequence from SQLite and poll for newly committed events. Polling starts at 250 milliseconds while catching up and backs off to 1 second while idle. Queries are bounded and indexed by `(run_id, sequence)`. This gives every server the same behavior without a daemon, broker, filesystem watcher, or hidden leader.

Delivery is at-least-once. Duplicate producer callbacks are normalized when they carry a stable source identity; otherwise they become distinct committed product events. Subscribers deduplicate by run and sequence. Slow subscribers are disconnected from bounded transport queues and recover from the durable cursor.

Persist at most one progress sample per run every 250 milliseconds for equivalent high-frequency progress updates. Coalescing records omitted count and time range in the next persisted event and durable snapshot. Lifecycle, warning, finding, artifact, cancellation, recovery, and terminal events are never coalesced by that policy.

### 9. Retention, compaction, and quota

Use these defaults:

| Limit | Default |
| --- | --- |
| Encoded event size | 16 KiB |
| Retained events per run | 4,096 |
| Retained event bytes per run | 16 MiB |
| Full terminal event history | 7 days |
| Run tombstone | 30 days after terminal |
| Workspace run-control soft quota | 496 MiB |
| Workspace hard quota | 512 MiB |
| Reserved headroom for active terminal/recovery commits | 16 MiB |

Only full-history duration, tombstone duration, and workspace quota are user-configurable through the existing config precedence. Validation enforces a minimum 1-hour full retention, minimum 24-hour tombstone retention, minimum 64 MiB hard quota, and at least 16 MiB reserved headroom. Event-size and lease invariants are fixed this sprint.

Compaction order is deterministic:

1. compact expired terminal progress detail while preserving the snapshot, sequence boundaries, omission totals, warnings/findings/artifact references, cancellation facts, and terminal event
2. convert terminal runs past full retention to tombstones
3. remove tombstones past tombstone retention
4. checkpoint WAL and incrementally vacuum only outside acceptance/append transactions

Active runs are never deleted. Per-run bounds may compact old progress detail and advance `oldest_retained_sequence`; the durable snapshot makes the gap explicit. At 80 percent of quota, background compaction starts. At the soft quota, new acceptance fails closed. Reserved headroom is used only for active-run heartbeat, cancellation, recovery, and terminal commits. At hard quota, owners cancel active work and report persistence degradation rather than continue silently.

### 10. Reconciliation

Run one bounded reconciliation pass during application startup before accepting new runtime work. Web/TUI processes and active CLI owners repeat it every 10 seconds. Multiple reconcilers are allowed; all decisions are idempotent fenced transactions, so there is no leader election.

For an active run:

- unexpired lease: leave unchanged
- expired lease with exact live process birth: classify stalled and retain active lifecycle
- expired lease with missing/mismatched owner during grace: record evidence and wait
- owner absent after grace with no proven cleanup: propose `interrupted` or `cleanup_uncertain`, never success
- existing terminal winner: perform no lifecycle rewrite

Product-specific reconciliation remains in sprint/study modules. Run control may call those existing app capabilities and correlate their result, but artifact presence, missing PID, or product status alone cannot prove operational success.

### 11. Persistence failure policy

| Failure point | Decision |
| --- | --- |
| Before acceptance | Fail closed; no child or operation starts. |
| Attempt claim | Leave accepted run inspectable, return failure, and permit conservative retry of claim. |
| Event append | Retry within a 5-second bound; do not fan out. If still failing, cancel owned work and mark local health degraded. |
| Heartbeat | Retry on the next control tick; if no heartbeat can commit before lease duration, cancel owned work. |
| Cancellation request | Return failure and do not claim cancellation was requested. |
| Terminal commit | Retry for 30 seconds on a cleanup context. If unavailable, let the lease expire; later reconciliation records interruption/uncertainty, never inferred success. |
| Corrupt/unsupported database | Refuse acceptance and mutation; expose typed health/recovery diagnostics and preserve the database for support. |

There is no in-memory authoritative fallback and no silent dual-write sidecar. Existing product-owned cleanup markers remain authoritative for their own product concern and can be correlated during recovery.

### 12. Schema migration and rollback

Use transactional schema migrations keyed by SQLite `user_version` plus an application schema table. Startup takes an advisory workspace migration lock, verifies no unsupported newer schema, checkpoints WAL, creates a timestamped bounded backup, applies one-way migrations, runs integrity checks, and only then enables acceptance.

Initial migration creates an empty operational repository. It does not synthesize runs from flow state, execute state, study state, locks, artifacts, or runtime checkpoints. Recognized historical `op_*` links without a durable mapping receive truthful legacy recovery guidance.

Rollback is operational, not dual-write: stop all UltraPlan processes, restore the pre-migration database backup, and run the matching binary. A pre-Sprint-35 binary is not a safe rollback after durable run control is enabled because it can start unrecorded work; documentation must require restoring both binary and workspace operational state together.

### 13. Telemetry and support

Structured logs use one vocabulary for request, run, attempt, owner/fence, stage/task, runtime, agentwrap, sequence, lifecycle transition, persistence operation, cancellation, reconciliation, and terminal winner.

Local counters and histograms cover acceptance/append/terminal latency and failure, active/stalled runs, lease renewal, reconciliation backlog, retained bytes/events, compaction, replay gaps, subscriber lag/drop, and cancellation routing. IDs are log fields, not metric labels.

`ultraplan run diagnostics --json` reports repository/schema health, quota, WAL/compaction state, stale owners, backlog, and safe per-run facts. An explicit bounded support bundle includes redacted snapshots, event headers/omission facts, health, config sources, and logs. It excludes prompts, provider payloads, repository content, unsafe paths, credentials, and arbitrary stdout/stderr. OpenTelemetry export is deferred behind the stable correlation vocabulary.

### 14. Verification and implementation shape

Core behavior is testable without a real provider through injected repository, clock, owner/process probe, notifier, and fake app operation. SQLite contract tests use temporary real local databases and multiple processes.

Required verification includes:

- acceptance atomicity and proof no child starts on storage failure
- concurrent sequence allocation, owner claims, fencing, heartbeat, cancellation, and terminal compare-and-set under `go test -race`
- two CLI owners plus two web observers against one workspace
- process kill at every acceptance/claim/append/terminal boundary
- stale lease, clock jump, PID reuse, process-birth mismatch, and no unsafe signal
- bounded replay, coalescing metadata, compaction, cursor gaps, slow subscribers, quota, disk full, corruption, migration, and rollback fixture
- session rotation and observer restart without run loss
- product artifact authority unchanged
- normal tests, race tests, CLI build, browser suite, and gated real-runtime dogfood

The implementation should first extract app-level run acceptance from the web hub, then add the run-control service/repository, integrate every runtime-backed CLI/TUI/web entry, demote the hub to transport delivery, add reconciliation/retention, and finally migrate compatibility/API/UI/docs. Do not build a daemon, scheduler, remote protocol, generic virtual filesystem, or authored-artifact database.

## Trade-Offs

| Decision | Benefit | Cost and rejected alternative |
| --- | --- | --- |
| Focused `internal/runcontrol` package | One clear operational authority shared by every surface without product workflow leakage. | Adds a product module. Rejected: `internal/web`, generic `platform`, or duplicated repositories in sprint/study. |
| SQLite with direct fenced writers | Atomic acceptance, ordered append, CAS terminal state, bounded queries, and multi-process access in one local primitive. | Adds driver/schema/migration/locking complexity. Rejected: snapshot plus append files because multi-writer sequence/CAS/compaction would require rebuilding database semantics. |
| No coordinator daemon | The CLI remains self-contained and no hidden process is required for truth. | Every process opens the repository and long-lived observers poll. Rejected: daemon because startup, upgrade, authority, and failure recovery exceed the local requirement. |
| Same-host only | Makes process birth, loopback security, SQLite locking, and one clock domain defensible. | Shared-filesystem multi-host observation is not guaranteed. Rejected: conditional distributed claims without remote identity and transport security. |
| Acceptor owns worker, no adoption | Preserves Go context and agentwrap ownership and avoids duplicate workers. | Owner death interrupts work instead of seamless failover. Rejected: adoption because current stages expose product resume, not transferable process ownership. |
| Lease plus birth identity and fencing | Protects against PID reuse and stale writers while remaining conservative. | More metadata and platform-specific process probes. Rejected: PID-only locks and lease-only terminal inference. |
| Durable polling for cross-process live events | Repository remains truth and every server behaves identically without a broker. | Up to one-second idle latency and read load. Rejected: in-memory pub/sub as authority and broker/daemon infrastructure. |
| Fixed safety timings | Testable, documented liveness contract with no unsafe config combinations. | Operators cannot tune unusual hosts this sprint. Rejected: exposing lease internals before operational evidence exists. |
| Bounded retention with reserved headroom | Prevents disk exhaustion while preserving active terminal/recovery commits and explicit gaps. | Historical detail expires. Rejected: unlimited history or silent event dropping. |
| Fail closed/degrade visibly | Never advertises a run or event that durable truth cannot support. | Storage faults may cancel otherwise healthy runtime work. Rejected: stale in-memory continuation that cannot be reconciled truthfully. |

## Evidence

### Repository evidence

- `projects/ultraplan-go/sprints/35-durable-run-observability/code-context.md` identifies `internal/app` as the common operation seam, `internal/web/operations.go` as the current process-local authority, and `internal/platform/runtime` as the agentwrap boundary. The selected package and dependency direction are an inference from those facts and the requirement that web remain an adapter.
- `../ultraplan-go/internal/web/operations.go` currently combines IDs, session filtering, cancellation contexts, sequence allocation, event retention, subscribers, active counts, and first-writer terminal state under one mutex. This proves the behavior to extract and the bounded transport behavior to retain.
- `../ultraplan-go/internal/app/operations.go` and `../ultraplan-go/internal/app/operation_runner.go` provide the adapter-neutral prepare/run vocabulary and shared runtime-backed dispatch. They are the acceptance and event-normalization integration points.
- `../ultraplan-go/internal/platform/runtime/runtime.go` already correlates agentwrap run/session/attempt identities, maps safe canonical events, omits raw payloads, and owns runtime cancellation. This supports correlation rather than duplicate supervision.
- `../ultraplan-go/internal/sprint/execute_state.go` demonstrates fsync, atomic rename, schema validation, and legacy handling for product-owned state. It establishes the durability quality bar while also showing why product state must remain separate.

### Selected report findings

- `studies/go-cli-study/reports/final/01-project-structure.md` finds that thin interfaces, unidirectional dependencies, and domain-owned internal packages are the dominant maintainable shape. This supports a focused internal module outside every presentation surface.
- `studies/go-cli-study/reports/final/02-command-architecture.md` finds that command wrappers should delegate to shared action/use-case layers and centralize lifecycle concerns. This supports app-level acceptance around existing product operations rather than command-specific run logic.
- `studies/go-cli-study/reports/final/03-dependency-injection.md` finds centralized manual composition, constructor injection, and narrow interfaces more traceable than globals or context service locators. This supports explicit repository/clock/process dependencies and no singleton registry.
- `studies/go-cli-study/reports/final/04-configuration-management.md` finds explicit precedence and post-merge validation essential. This supports configuring only retention/quota through existing precedence and validating safety relationships centrally.
- `studies/go-cli-study/reports/final/05-error-handling.md` supports wrapped causes, typed recovery data, and separate user/operator rendering. The fail-closed matrix therefore uses typed operational failures rather than substring classification.
- `studies/go-cli-study/reports/final/06-io-abstraction.md` supports injectable filesystems/backends/clocks for deterministic faults and warns that partial abstraction leaves critical paths untestable. The repository, clock, and process probe are explicit test seams.
- `studies/go-cli-study/reports/final/07-state-context.md` distinguishes cancellation context from persistent identity and cites refreshable locks plus cleanup contexts. It also finds no distributed session coordination. This supports leases/fencing, separate cleanup retries, and the same-host boundary.
- `studies/go-cli-study/reports/final/08-concurrency.md` supports localized goroutine ownership, bounded queues, timeout waits, and one-time cleanup. This informs the owner control loop, bounded polling/subscribers, and idempotent cancellation.
- `studies/go-cli-study/reports/final/10-logging-observability.md` supports structured fields and local metrics while finding OpenTelemetry absent. This informs the stable correlation vocabulary and trace deferral.
- `studies/go-cli-study/reports/final/11-testing-strategy.md` supports real command-path integration, centralized fakes, fixtures, and behavior assertions. This informs SQLite contract, process, browser, and fault-matrix testing.
- `studies/go-cli-study/reports/final/12-extensibility.md` supports explicit versioned metadata while warning against costly extension systems. This informs schema migrations and rejection of daemon/plugin/remote-worker machinery.
- `studies/go-cli-study/reports/final/13-security.md` supports explicit trust boundaries, secret-aware redaction, schema validation, and bounded diagnostics. This informs private storage permissions and redaction before append.
- `studies/go-cli-study/reports/final/14-performance.md` supports streaming, bounded structures, bounded concurrency, and disk-backed long sessions. It provides precedent for SQLite persistence but does not itself prove the store choice; the transaction and multi-writer requirements do.

## Risks

- `modernc.org/sqlite` increases binary size and introduces a substantial dependency. Build, cross-compilation, corruption, and upgrade tests must justify it.
- WAL and file locking may behave differently on filesystems that appear local but do not provide local locking semantics. Health must detect and reject unsupported placement rather than weaken durability.
- A 15-second lease can classify a healthy but blocked owner as stalled. The 45-second grace and exact process-birth check prevent that classification from becoming premature terminal state.
- Polling from multiple servers adds read pressure. Indexed bounded queries, idle backoff, subscriber limits, and measurements are required before changing cadence.
- Reserved disk headroom cannot guarantee terminal commit under total filesystem exhaustion. Failure remains visible and reconciliation conservative; documentation must explain external disk recovery.
- Mid-run persistence loss deliberately cancels work. That protects truth but may waste expensive runtime progress.
- SQLite migration or corruption recovery can block all new runtime work. A tested backup/export/integrity workflow is part of the feature, not optional operations polish.
- Platform-specific process-birth checks differ on Linux and macOS. Unsupported or uncertain checks must produce conservative liveness, never PID-only authority.
- Product state and operational terminal state can disagree after a storage outage. Run detail must present both authorities and recovery guidance without merging them into invented success.
- Fixed timings and retention defaults may need tuning after dogfood. Future changes must preserve the observable lease, gap, and tombstone contracts and be versioned/config-validated.

The architecture decision is to refactor first: extract durable run control from the web hub into `internal/runcontrol`, backed by same-host direct multi-process SQLite transactions, then integrate every surface. The main trade-off is accepting a local database dependency and polling latency to obtain truthful atomic cross-process control without a daemon or distributed protocol.
