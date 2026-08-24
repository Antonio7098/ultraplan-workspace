# Sprint Code Context

## Sprint Scope

Sprint 37 adds evidence-backed QA orchestration and smoke-suite compatibility without creating a second smoke execution path.

The implementation scope includes:

- Add a sprint QA policy, evidence model, adjudication model, issue model, deterministic assessment, and verification state under `internal/sprint`.
- Add QA artifact handling for `qa.md` and any required verification evidence while preserving prior valid artifacts on failed reruns.
- Support `qa --suite smoke` by delegating to the existing smoke manifest, discovery, selection, invocation, validation, retention, and reporting behavior.
- Keep QA outside `PlanningStage`; QA is an assurance operation, not a planning lifecycle stage.
- Run review agents with frozen governed inputs and read-only/default-deny permissions.
- Use product-neutral process isolation, timeout, cancellation, bounded output, and cleanup mechanics from `internal/platform/process`.
- Expose shared operation DTOs and use cases through `internal/app`.
- Preserve durable acceptance, owner claims, correlation, fencing, cancellation, and terminal arbitration in `internal/runcontrol`.
- Extend CLI and TUI surfaces through existing application boundaries.
- Keep `internal/web` dependent on `internal/app`, never directly on `internal/sprint`.
- Preserve current smoke command and verification compatibility.

Implementation not found during inspection includes the Sprint 36 map, shard, theory, synthesis, and verification directory described by the requirements, as well as Sprint 37-specific QA domain, isolation, TUI QA, and web QA files.

## Constraints

- QA must not be added to the planning-stage enumeration or treated as a planning transition.
- Existing smoke orchestration is the sole execution authority for `qa --suite smoke`.
- Smoke-as-QA must preserve manifest discovery, case selection, environment construction, timeout handling, evidence validation, attempt retention, `smoke.md` publication, and exit behavior.
- Assurance inputs must be frozen before review fan-out so concurrent mutation cannot change the evidence set being adjudicated.
- Review agents must use read-only/default-deny runtime permissions.
- Process execution must use explicit argument vectors rather than shell interpolation.
- Child output must remain bounded.
- Timeout and cancellation must terminate the process tree where the platform can prove that behavior.
- Platforms that cannot prove descendant cleanup must expose that uncertainty rather than claim successful isolation.
- All sprint artifact paths must remain within the resolved sprint directory.
- Git-backed execution must remain tied to the governed target and worktree; non-Git execution must retain equivalent containment checks.
- State loading must remain strict, migration must be explicit, and publication must remain atomic.
- Failed or interrupted reruns must not destroy prior valid QA or smoke evidence.
- Mutating sprint operations must use the existing per-sprint mutation lease.
- Durable operation acceptance and owner claim must occur before operation work starts.
- Cancellation and completion must remain fenced so stale owners cannot publish terminal state.
- `internal/web` must not import `internal/sprint`.
- QA configuration needs explicit bounded limits rather than unbounded agent counts, evidence sizes, output sizes, or execution durations.
- No manual artifact persistence or promotion should be introduced outside UltraPlan-owned publication.

## Inspected Repository Areas

- Sprint domain and dependency seams: planning-stage modeling, runtime dependencies, locking, artifact naming, state loading, migrations, and atomic publication.
- Target and context governance: Git and non-Git target resolution, worktree containment, prompt-context construction, and bounded governed inputs.
- Review orchestration: read-only agent execution, frozen snapshots, evidence state, and permission policy.
- Smoke orchestration: manifest discovery, case selection, invocation, timeout handling, evidence validation, attempt retention, report publication, and compatibility tests.
- Verification: review-to-smoke compatibility and deterministic assessment behavior.
- Platform process execution: explicit argv, bounded output, cancellation, timeout, and platform-specific process cleanup.
- Runtime permissions: sandbox DTOs, writable-path policies, and agent adapter validation.
- Configuration: smoke limits, run-control settings, and runtime permission defaults.
- Application layer: operation DTOs, preparation, stale-input guards, operation runners, durable operation lifecycle, sprint projections, and command routing.
- Run control: durable repository interfaces, events, correlation, fencing, owner claims, cancellation, snapshots, and terminal arbitration.
- TUI: sprint operation dispatch, durable start, progress state, and cancellation.
- Web: import-boundary enforcement, operation handlers, same-origin and CSRF protection, bounded request handling, and server-rendered sprint presentation.
- Documentation: architecture boundaries, JSON compatibility, recovery behavior, and local web operation semantics.

## Relationships

- `internal/sprint` owns QA policy, QA evidence, adjudication, issue derivation, assessment, verification state, `qa.md`, and smoke compatibility behavior.
- `internal/sprint` consumes process and runtime abstractions through dependency seams rather than embedding operating-system process logic.
- `internal/platform/process` owns product-neutral child-process execution, bounded output, timeout, cancellation, process-group cleanup, and cleanup certainty.
- `internal/platform/runtime` owns sandbox and permission-policy DTOs used to enforce read-only/default-deny review agents.
- `internal/app` translates adapter requests into sprint operations and exposes adapter-neutral operation state.
- `internal/runcontrol` remains the durable authority for acceptance, ownership, correlation, fences, cancellation, events, and terminal status.
- CLI and TUI adapters invoke application-layer use cases rather than duplicating QA orchestration.
- Web handlers depend on `internal/app`; the existing import-boundary test prevents direct web-to-sprint coupling.
- QA publication must share sprint artifact containment, mutation locking, strict state loading, migration, and atomic-write behavior.
- Smoke-as-QA must call the existing smoke orchestration rather than reproduce manifest parsing or command execution.
- Verification and QA should consume the same retained smoke evidence so compatibility results cannot diverge.
- Review fan-out must consume one frozen target and governed context snapshot.
- Process results feed QA evidence and adjudication, while durable lifecycle results remain controlled by run control.

## Selected Source References

### Sprint Domain Boundary

- **Path:** `internal/sprint/domain.go`
- **Lines:** 1-30
- **Rationale:** Establishes the sprint domain boundary containing planning-stage and sprint-state concepts. QA must integrate without becoming a planning stage, and any bounded assurance projection should remain compatible with these domain types.

### Sprint Dependency Seams

- **Path:** `internal/sprint/service.go`
- **Lines:** 1-30
- **Rationale:** Establishes the service dependency surface through which runtime, process, clock, smoke, and mutation capabilities are supplied. QA orchestration should extend these seams instead of constructing platform dependencies directly.

### Sprint State Persistence

- **Path:** `internal/sprint/state.go`
- **Lines:** 1-30
- **Rationale:** Anchors strict sprint-state loading and publication responsibilities. QA state and verification state must follow existing schema, migration, preservation, and atomic-publication rules.

### Governed Target Resolution

- **Path:** `internal/sprint/execute_target.go`
- **Lines:** 1-30
- **Rationale:** Defines the target-resolution boundary used to bind execution to the governed repository or worktree. QA isolation and evidence must retain the same target identity and containment guarantees.

### Governed Prompt Context

- **Path:** `internal/sprint/prompt_context.go`
- **Lines:** 1-30
- **Rationale:** Establishes bounded source-context construction for governed agent inputs. QA review fan-out should freeze and reuse this context rather than permit agents to inspect changing or unbounded inputs.

### Review Orchestration

- **Path:** `internal/sprint/review.go`
- **Lines:** 1-30
- **Rationale:** Establishes the existing review orchestration boundary, including agent-oriented evidence processing. Sprint 37 QA review behavior should reuse its frozen-input and restricted-permission patterns.

### Smoke Result Contracts

- **Path:** `internal/sprint/smoke_types.go`
- **Lines:** 1-30
- **Rationale:** Defines smoke lifecycle, verdict, result, and retained completion concepts that `qa --suite smoke` must consume without inventing an incompatible result model.

### Canonical Smoke Orchestration

- **Path:** `internal/sprint/smoke.go`
- **Lines:** 1-30
- **Rationale:** Anchors the canonical smoke execution path. Smoke-as-QA must delegate to this orchestration so discovery, invocation, evidence validation, retention, and report publication remain identical.

### Smoke Protocol

- **Path:** `internal/sprint/smoke_protocol.go`
- **Lines:** 1-30
- **Rationale:** Establishes the smoke protocol boundary for manifest and discovery behavior, selection, execution environment, and timeout policy. QA must not independently reinterpret these contracts.

### Verification Compatibility

- **Path:** `internal/sprint/verify.go`
- **Lines:** 1-30
- **Rationale:** Establishes verification behavior that combines review and smoke outcomes. QA assessment and retained evidence must preserve this compatibility rather than create conflicting assurance verdicts.

### Platform Process Boundary

- **Path:** `internal/platform/process/process.go`
- **Lines:** 1-30
- **Rationale:** Defines the product-neutral child-process abstraction used for explicit argv, bounded output, timeout, cancellation, and cleanup reporting. QA command isolation belongs behind this boundary.

### Runtime Permission Model

- **Path:** `internal/platform/runtime/runtime.go`
- **Lines:** 1-30
- **Rationale:** Defines sandbox and path-policy concepts used by agent execution. QA reviewers should use these types to enforce read-only/default-deny access.

### Runtime Adapter Policy Mapping

- **Path:** `internal/platform/runtime/agentwrap.go`
- **Lines:** 1-30
- **Rationale:** Establishes validation and adapter mapping for runtime permission policy. QA agent requests must pass through this mapping rather than bypass policy checks.

### Application Operation Contracts

- **Path:** `internal/app/operations.go`
- **Lines:** 1-30
- **Rationale:** Defines shared operation DTOs and preparation boundaries. QA should be exposed as an adapter-neutral application operation with stale-input protection.

### Durable Operation Lifecycle

- **Path:** `internal/app/durable_operations.go`
- **Lines:** 1-30
- **Rationale:** Establishes durable acceptance, owner claim, operation cancellation context, and completion flow. QA execution must enter and leave the durable lifecycle through this path.

### Application Command Routing

- **Path:** `internal/app/sprint_commands.go`
- **Lines:** 1-30
- **Rationale:** Establishes sprint command dispatch and existing smoke and verification compatibility. CLI QA routing should be added here without duplicating sprint orchestration.

### Run-Control Authority

- **Path:** `internal/runcontrol/interfaces.go`
- **Lines:** 1-30
- **Rationale:** Defines the durable repository and control interfaces that own operation acceptance, claims, fencing, events, cancellation, and terminal state. QA must use these interfaces rather than persist operational authority separately.

### Run-Control Model

- **Path:** `internal/runcontrol/model.go`
- **Lines:** 1-30
- **Rationale:** Establishes operation identity, correlation, fence, event, and snapshot concepts. QA progress and terminal results must preserve these durable semantics.

### TUI Operation Integration

- **Path:** `internal/tui/model.go`
- **Lines:** 1-30
- **Rationale:** Establishes the TUI model boundary where sprint operations are represented. QA status and commands should extend this model through application DTOs.

### Web Import Boundary

- **Path:** `internal/web/import_boundary_test.go`
- **Lines:** 1-30
- **Rationale:** Enforces the architectural rule that web code cannot import sprint internals. Any QA web surface must depend exclusively on `internal/app`.

### Web Operation Transport

- **Path:** `internal/web/operation_handlers.go`
- **Lines:** 1-30
- **Rationale:** Establishes strict transport handling for durable operations. A QA web endpoint should reuse these parsing, validation, and operation-start patterns.

### Web Security Boundary

- **Path:** `internal/web/security.go`
- **Lines:** 1-30
- **Rationale:** Establishes same-origin, CSRF, session, and request-bound controls. Any state-changing QA web action must retain these protections.

## Open Questions

- Where are the Sprint 36 map, shard, theory, synthesis, and verification implementations required as the foundation for Sprint 37?
- What exact persisted schema represents QA policy, evidence, adjudication, issues, assessment, and verification state?
- What explicit migration version introduces the Sprint 37 persisted fields and artifacts?
- What are the configured upper bounds for QA agents, evidence count, evidence bytes, agent output, command output, and total QA duration?
- What is the exact deterministic precedence when infrastructure failure, agent disagreement, smoke failure, cancellation, and stale input overlap?
- Which fields from retained smoke results must appear in `qa.md`, and which remain linked through existing `smoke.md` evidence?
- Does a failed QA rerun publish a failed attempt record while preserving the prior successful `qa.md`, or does it leave only durable operation evidence?
- What platform behavior is required when descendant-process cleanup cannot be proven?
- Which QA operations are exposed initially through CLI, TUI, and web, and which adapters are intentionally deferred?
- What exact JSON compatibility guarantees apply to new QA operation projections and verification state?
- Which source files own the final Sprint 37 QA implementation, isolation implementation, TUI QA view, and web QA presentation once added?
