# Sprint Code Context

## Sprint Scope

Sprint 36 adds a distinct read-only QA verification capability after execute and alongside, not inside, the planning-stage model. It must deterministically map accepted changed paths into bounded behavioral shards, run mutation-contained investigations, persist falsifiable theory outcomes and bounded follow-up, synthesize all outcomes without issue promotion, and project one durable status across CLI, JSON, TUI, and browser surfaces.

The implementation must reuse the existing durable run-control authority for run identity, cancellation, replay, terminal arbitration, and recovery while keeping detailed QA maps, attempts, theories, evidence references, and synthesis in verification-owned versioned state. Existing `review` and `review.md` behavior remains compatible but is presented as Conformance Review; existing smoke and verify authority must not change.

## Inspected Repository Areas

- `internal/sprint`: planning and verification domain types, flow state, artifact containment, execute changed-path evidence, review preparation/execution/validation, smoke state, verification locking, and recovery behavior.
- `internal/runcontrol`: generic durable run identities, lifecycle, cancellation, event replay, fencing, terminal arbitration, reconciliation, retention, and cross-process tests.
- `internal/app`: shared operation vocabulary, governed-input fingerprints, runtime-backed dispatch, durable command integration, summaries, error classification, and adapter-neutral use cases.
- `internal/platform/runtime`: generic runtime requests, sandbox and permission-policy mapping, structured validation, and safe runtime metadata.
- `internal/web`: in-memory transport projection, browser operation contracts, bounded events, and the enforced dependency boundary on `internal/app`.
- CLI, TUI, JSON, and browser compatibility tests relevant to adding QA operations and projections without creating independent workflow semantics.
- Existing deterministic, malformed-output, stale-state, read-only request, atomic persistence, cancellation, recovery, and cross-process tests that establish patterns for Sprint 36 coverage.

## Selected Source References

### Verification And Planning Domain Boundary

- **Path:** `internal/sprint/domain.go`
- **Lines:** `10-59`
- **Symbol:** `FlowStateSchemaVersion`, `PlanningStage`, `StageState`
- **Rationale:** Defines the closed planning-stage vocabulary and versioned flow-state errors. QA lifecycle must not be added as an overloaded planning stage, and schema evolution must preserve the existing migration/error contract.

- **Path:** `internal/sprint/domain.go`
- **Lines:** `122-169`
- **Symbol:** `ExecuteRunState`, `FlowState`, `VerificationAttempt`
- **Rationale:** Shows current ownership boundaries: execute has detailed independent state, while flow state embeds only review and smoke summaries. It also provides the existing bounded attempt vocabulary and illustrates why detailed QA state should have a separate contract.

- **Path:** `internal/sprint/domain.go`
- **Lines:** `197-244`
- **Symbol:** `VerificationStage`, `VerificationStatus`, `StatusSummary`
- **Rationale:** This is the current cross-surface verification projection. It presently identifies verification stages with `PlanningStage` and only exposes review and smoke, making it a central compatibility and migration boundary for an independent QA phase.

- **Path:** `internal/sprint/domain.go`
- **Lines:** `264-337`
- **Symbol:** `PlanningStages`, `ValidStage`, `validateAttempt`
- **Rationale:** Establishes the closed ordered planning-stage invariant and bounded diagnostic validation. QA must remain outside `PlanningStages`, while its own state validation should retain equivalent fail-closed identity, lifecycle, and diagnostic bounds.

### Flow-State Authority And Persistence

- **Path:** `internal/sprint/state.go`
- **Lines:** `20-80`
- **Symbol:** `LoadFlowState`
- **Rationale:** Implements database-first loading, strict JSON decoding, supported-version checks, migration, and final semantic validation. It is the model for unsupported, malformed, legacy, and partially compatible state behavior.

- **Path:** `internal/sprint/state.go`
- **Lines:** `150-191`
- **Symbol:** `migrateFlowStateV1`
- **Rationale:** Demonstrates the existing one-version migration policy and explicit preservation of historical review/smoke outcomes as stale, unverifiable evidence. QA schema evolution and invalidation should be reasoned against this precedent.

- **Path:** `internal/sprint/state.go`
- **Lines:** `201-288`
- **Symbol:** `SaveFlowState`, `saveFlowStateWithHooks`
- **Rationale:** Defines preservation of completed verification summaries, database authority, checkpoint behavior, and atomic file replacement with flush and directory sync. Any QA summary pointer committed to flow state must respect this ordering and atomicity boundary.

- **Path:** `internal/sprint/state.go`
- **Lines:** `291-358`
- **Symbol:** `ValidateFlowState`
- **Rationale:** Enforces exact stage order, path containment, supported outcomes, and review/smoke validation before persistence is accepted. This is the canonical location affected by any additive QA summary and pointer fields.

- **Path:** `internal/sprint/artifacts.go`
- **Lines:** `11-74`
- **Symbol:** `ArtifactRelPath`, `resolveSprintContained`
- **Rationale:** Centralizes canonical review/smoke compatibility paths and sprint-root containment. A verification-owned QA state layout must use equivalent containment without introducing canonical `qa.md`.

### Execute Evidence And Deterministic Inputs

- **Path:** `internal/sprint/review.go`
- **Lines:** `91-165`
- **Symbol:** `ReviewManifest`, `ReviewCoverageResult`, `ReviewRequest`, `ReviewResult`
- **Rationale:** Provides the existing frozen-input manifest, changed-path list, bounded coverage-result schema, progress callback, and compatibility result shape used by Conformance Review. QA can reuse evidence-selection lessons without conflating theory outcomes with review findings.

- **Path:** `internal/sprint/review.go`
- **Lines:** `261-370`
- **Symbol:** `Service.PrepareReview`
- **Rationale:** Shows how current verification collects execute run state, obtains changed paths, validates target containment/readability, adds selected contracts, sorts inputs and coverage, and computes a deterministic manifest fingerprint. This is the nearest implemented input-governance precedent for the QA mapper.

- **Path:** `internal/sprint/review.go`
- **Lines:** `1371-1414`
- **Symbol:** `reviewChangedPaths`, `excludeGovernedReviewPaths`
- **Rationale:** Defines the currently accepted source of changed paths: execute run-state `files` plus task evidence paths, deduplicated and sorted, with governed workspace artifacts excluded. Sprint reasoning must decide whether this is sufficiently authoritative for QA and how missing or deleted paths block mapping.

- **Path:** `internal/sprint/review_test.go`
- **Lines:** `130-187`
- **Symbol:** `TestReview`
- **Rationale:** Verifies deterministic fingerprints, changed-path filtering, canonical `review.md`, preservation of the last valid artifact after malformed output, and read-only reviewer requests. It is a critical regression test pattern for Conformance Review compatibility and QA determinism.

### Conformance Review Compatibility And Safety

- **Path:** `internal/sprint/review.go`
- **Lines:** `23-89`
- **Symbol:** `StageReview`, `ReviewStageState`, `ReviewResumeState`, `ReviewCompletion`
- **Rationale:** Defines current review lifecycle, verdicts, resumable coverage checkpoints, and last-complete canonical evidence. User-facing relabeling to Conformance Review must preserve these types, statuses, commands, and `review.md` consumers.

- **Path:** `internal/sprint/review.go`
- **Lines:** `870-924`
- **Symbol:** `Service.runReviewer`
- **Rationale:** Shows existing reviewer isolation: read-only sandbox, restricted permissions, default-deny tools, required permission capability, structured output validation, session resume, and fail-closed handling when enforcement is unsupported. It is the closest current implementation of read-only investigation, but it does not permit safe command execution or perform post-attempt mutation verification.

- **Path:** `internal/sprint/review_runtime_validation.go`
- **Lines:** `15-110`
- **Symbol:** `extractValidatedReviewResult`, `reviewValidationSpec`
- **Rationale:** Demonstrates runtime-success-independent acceptance: structured output is extracted, schema and citations are validated, repair attempts are bounded, and cancellation is retained. QA state must similarly reject successful runtime calls whose persisted domain result is invalid.

- **Path:** `internal/sprint/review_runtime_validation.go`
- **Lines:** `192-238`
- **Symbol:** `reviewResultProblems`
- **Rationale:** Provides a concrete fail-closed semantic validator with schema, identity, enum, size, uniqueness, citation, and diagnostic-count bounds. It is a useful precedent for validating theories, evidence references, and synthesis output.

- **Path:** `internal/sprint/review_runtime_validation.go`
- **Lines:** `250-303`
- **Symbol:** `reviewManifestChanges`
- **Rationale:** Compares frozen inputs at promotion and emits bounded invalidation reasons. QA shard reuse and synthesis freshness require a stronger but analogous governed-fingerprint recheck.

### Existing Smoke And Verify Compatibility

- **Path:** `internal/sprint/smoke_types.go`
- **Lines:** `11-65`
- **Symbol:** `StageSmoke`, `SmokeExecutionStatus`, `SmokeVerdict`, `SmokePhase`, `SmokeRequest`
- **Rationale:** Defines the current empirical smoke protocol and lifecycle. QA must not absorb, rename, or alter these semantics, and the use of `PlanningStage` here highlights an existing coupling that should not be repeated for QA.

- **Path:** `internal/sprint/smoke_types.go`
- **Lines:** `162-247`
- **Symbol:** `SmokeResult`, `SmokeStageState`, `SmokeCompletion`
- **Rationale:** Captures canonical smoke evidence, issue records, changed paths, freshness fingerprints, active attempts, and last-complete outcomes. QA theories must remain non-canonical and must never populate these issue or evidence-adjudication fields.

- **Path:** `internal/sprint/smoke_types.go`
- **Lines:** `288-329`
- **Symbol:** `validateSmokeStageState`
- **Rationale:** Shows current smoke state validation for status, verdict, canonical path, containment, diagnostics, attempts, overrides, issues, and last-complete identity. It is a compatibility boundary and a model for strict QA validation.

### Generic Durable Run Control

- **Path:** `internal/runcontrol/model.go`
- **Lines:** `16-118`
- **Symbol:** `Lifecycle`, `CancellationState`, `TerminalOutcome`
- **Rationale:** Defines the Sprint 35 lifecycle and terminal vocabulary that QA mapping, shard attempts, follow-up, and synthesis must reuse rather than duplicate.

- **Path:** `internal/runcontrol/model.go`
- **Lines:** `120-190`
- **Symbol:** `EventType`, `Target`
- **Rationale:** Establishes bounded generic event kinds and product-owned target correlation. QA may identify phase and shard work through target fields, but theory and synthesis payloads do not belong in this generic event authority.

- **Path:** `internal/runcontrol/model.go`
- **Lines:** `224-365`
- **Symbol:** `Correlation`, `Attempt`, `Fence`, `Event`, `Snapshot`
- **Rationale:** Defines durable correlation identities, lease-fenced attempts, replayable events, snapshots, and terminal invariants. Verification-owned QA records should point to these run and attempt identities while retaining domain outcomes separately.

- **Path:** `internal/runcontrol/model.go`
- **Lines:** `444-483`
- **Symbol:** `Acceptance`, `Claim`, `EventDraft`, `TerminalProposal`
- **Rationale:** Provides the command-side protocol for acceptance, ownership, event recording, and one terminal proposal, plus operational heartbeat and reconciliation bounds.

- **Path:** `internal/runcontrol/interfaces.go`
- **Lines:** `41-72`
- **Symbol:** `Notifier`, `Repository`, `Control`
- **Rationale:** Explicitly states that notifications are non-authoritative and that product workflow state is absent from the run repository. This is the key architectural constraint against storing QA maps, theories, or synthesis in run-control snapshots or events.

- **Path:** `internal/app/run_control.go`
- **Lines:** `18-73`
- **Symbol:** `runControlState.repository`
- **Rationale:** Owns one SQLite repository per workspace, applies retention policy, and reconciles at startup. QA application operations should enter through this existing composition rather than opening another repository or registry.

- **Path:** `internal/app/run_control.go`
- **Lines:** `122-302`
- **Symbol:** `controlledRuntime.StartRun`
- **Rationale:** Traces durable acceptance, claim and fencing, event persistence, cancellation polling and acknowledgement, heartbeat, reconciliation, persistence-failure handling, and terminal arbitration around every runtime call. This is the operational control plane Sprint 36 must reuse.

- **Path:** `internal/app/run_control.go`
- **Lines:** `352-374`
- **Symbol:** `targetFromRuntimeRequest`
- **Rationale:** Maps product metadata to bounded run-control target fields. QA run, phase, shard, and task correlation must fit this generic contract or extend it compatibly without importing QA semantics into runtime or runcontrol.

- **Path:** `internal/runcontrol/process_integration_test.go`
- **Lines:** `83-125`
- **Symbol:** `TestProcessIndependentObserverPersistsCancellation`
- **Rationale:** Proves cross-process cancellation and event observation through durable storage, which QA surfaces must rely on rather than process-local cancellation state.

- **Path:** `internal/runcontrol/process_integration_test.go`
- **Lines:** `127-233`
- **Symbol:** `TestProcessUnclaimedAcceptanceIsReconciledAfterOwnerExit`, `TestProcessClaimedOwnerExitIsInterruptedAndRepeatedReconciliationIsIdempotent`
- **Rationale:** Establishes recovery of accepted and claimed work after owner exit and idempotent reconciliation. QA recovery tests should cover each domain persistence boundary while relying on these run-level guarantees.

### Shared Application Operations

- **Path:** `internal/app/operations.go`
- **Lines:** `22-69`
- **Symbol:** `OperationalUseCases`, `WebOperations`, `DurableOperationManager`, `OperationReconciler`
- **Rationale:** Defines the adapter-neutral application boundary shared by terminal and browser surfaces. QA status and commands should be exposed here rather than implemented independently by CLI, TUI, or web packages.

- **Path:** `internal/app/operations.go`
- **Lines:** `71-167`
- **Symbol:** `OperationKind`, `OperationState`, `OperationRequest`, `Confirmation`, `OperationEvent`, `OperationResult`
- **Rationale:** Contains the closed cross-surface operation vocabulary, lifecycle projection, request scope, confirmation fingerprint, progress event, and typed error shapes that QA operations must extend compatibly.

- **Path:** `internal/app/operations.go`
- **Lines:** `169-282`
- **Symbol:** `dashboardUseCases.PrepareOperation`
- **Rationale:** Classifies runtime use and mutation, generates user-visible scope and warnings, selects governed inputs, and freezes an input fingerprint before execution. QA should be classified as runtime-backed but target-read-only while still acknowledging verification-state writes.

- **Path:** `internal/app/operations.go`
- **Lines:** `285-410`
- **Symbol:** `dashboardUseCases.RunOperation`
- **Rationale:** Rechecks prepared fingerprints, emits shared progress, dispatches read-only queries locally, and routes runtime-backed work through one runner. It is the central application dispatch point for consistent CLI/TUI/browser semantics.

- **Path:** `internal/app/operations.go`
- **Lines:** `511-592`
- **Symbol:** `governedOperationInputs`, `fingerprintOperationInputs`
- **Rationale:** Implements deterministic path ordering, symlink rejection, directory traversal, file hashing, and canonical request fingerprints. QA mapping needs a more complete governed-input set, including execute and review evidence, but should preserve these deterministic and fail-closed properties.

- **Path:** `internal/app/operations.go`
- **Lines:** `626-669`
- **Symbol:** `failedOperation`, `operationFailureMessage`
- **Rationale:** Maps cancellation, malformed state, unsupported schemas, references, stale fingerprints, and conflicts into stable cross-surface error categories. New QA blocked/invalid/corrupt-state errors must integrate here without leaking unsafe details.

- **Path:** `internal/app/operation_runner.go`
- **Lines:** `15-110`
- **Symbol:** `sharedOperationRunner`
- **Rationale:** Is the single runtime-backed implementation used by terminal and browser adapters. It demonstrates how review, smoke, and verify invoke sprint services and project typed progress, and is the correct integration boundary for QA orchestration.

- **Path:** `internal/app/run_control_inventory_test.go`
- **Lines:** `11-54`
- **Symbol:** `TestEveryRuntimeBackedCLIEntryUsesDurableAcceptanceInventory`
- **Rationale:** Guards against runtime-backed CLI entries bypassing durable acceptance. QA CLI commands must be added to this inventory or an equivalent typed mechanism.

### Runtime Permission Boundary

- **Path:** `internal/platform/runtime/runtime.go`
- **Lines:** `24-75`
- **Symbol:** `Request`, `PermissionPolicy`, `PermissionPathRule`
- **Rationale:** Defines the generic runtime request fields for working directory, timeout, sandbox, permission mode, per-tool and per-path policy, validation, and event observation. QA safety policy should be expressed through this generic boundary rather than adding QA semantics to the runtime package.

- **Path:** `internal/platform/runtime/agentwrap.go`
- **Lines:** `75-100`
- **Symbol:** `mapPermissionPolicy`
- **Rationale:** Maps and validates default, tool, path, unsupported-feature, and metadata policy into agentwrap. Unsupported or unenforceable command and filesystem restrictions must fail QA closed.

- **Path:** `internal/sprint/code_context.go`
- **Lines:** `292-326`
- **Symbol:** `Service.CodeContext` runtime request boundary
- **Rationale:** Provides another implemented target-read-only operation with default-deny inspection tools and explicit capability enforcement. It confirms that read-only target access can coexist with product-owned candidate promotion outside the target.

### Verification Concurrency

- **Path:** `internal/sprint/verification_lock.go`
- **Lines:** `14-101`
- **Symbol:** `acquireVerificationFileLock`, `verificationFileLock.release`
- **Rationale:** Existing review/smoke mutual exclusion is a PID-based workspace lock with stale-owner replacement and ownership-checked release. Sprint 36 must determine how this legacy lock collaborates with durable run-control ownership rather than creating conflicting lifecycle authority.

- **Path:** `internal/sprint/verification_lock_test.go`
- **Lines:** `11-40`
- **Symbol:** `TestVerificationFileLockRejectsLiveOwnerAndReplacesDeadOwner`
- **Rationale:** Captures current live-owner rejection and dead-owner recovery behavior that must remain compatible for review/smoke while QA adopts Sprint 35 durable arbitration.

### Cross-Surface Projections

- **Path:** `internal/app/sprint_usecases.go`
- **Lines:** `16-64`
- **Symbol:** `SprintSummary`, `ReviewSummary`, `SmokeSummary`
- **Rationale:** Defines the shared product projection consumed by TUI and web surfaces. QA readiness, run identity, phase, counts, blockers, freshness, and next actions need a corresponding typed summary here rather than surface-specific derivation.

- **Path:** `internal/app/sprint_usecases.go`
- **Lines:** `260-307`
- **Symbol:** `validateSprintStage`, `summarizeSmoke`, `summarizeReview`
- **Rationale:** Shows current stage-based validation and summary construction. QA requires independent validation and summarization rather than another `PlanningStage` switch case.

- **Path:** `internal/web/operations.go`
- **Lines:** `18-143`
- **Symbol:** operation limits and `operationHub`
- **Rationale:** Browser operation memory is deliberately bounded and transport-local, with limits on active operations, events, encoded bytes, subscribers, and retention. Detailed QA state cannot live here; browser views must reload durable application projections.

- **Path:** `internal/web/import_boundary_test.go`
- **Lines:** `12-41`
- **Symbol:** `TestWebImportBoundary`
- **Rationale:** Enforces that production web code imports only `internal/app` and the standard library. This prevents browser handlers from importing `internal/sprint` to map shards, infer completion, or persist QA truth.

- **Path:** `internal/web/operations_contract_test.go`
- **Lines:** `19-60`
- **Symbol:** `TestBrowserOperationKindContract`
- **Rationale:** Is the producer/consumer compatibility table for browser operation kinds. Any QA status, map, start, shard, restart, or synthesis operations exposed to the browser must update this contract coherently.

- **Path:** `internal/web/operations_contract_test.go`
- **Lines:** `82-151`
- **Symbol:** `TestBrowserLifecycleDocumentContract`, `TestBrowserSSEEventNameAndFrameContract`
- **Rationale:** Defines stable browser lifecycle and SSE event contracts, including downgrade of unknown producer events to progress. QA-specific domain outcomes should be fetched from durable summaries rather than smuggled into incompatible transport lifecycle states.

## Relationships

- `internal/sprint` owns sprint-specific planning, execute, review, smoke, artifact, freshness, and validation semantics. It is the natural existing owner for QA domain state, deterministic mapping, theory policy, follow-up, and synthesis.
- `internal/app` composes sprint services into typed operations and summaries. CLI, TUI, and browser adapters should invoke these shared operations and consume the same QA projection.
- `internal/platform/runtime` carries generic sandbox, permission, validation, model, and event capabilities. It must remain unaware of QA shards, theories, or synthesis.
- `internal/runcontrol` owns durable operational identity and lifecycle. A QA map run, shard attempt, follow-up, or synthesis operation can correlate to run-control IDs, but its domain records remain in verification-owned state.
- `flow-state.json` currently contains ordered planning state and compact review/smoke records. It may gain only canonical QA summary, freshness, verdict, and pointer data; detailed QA state belongs in a separate contained, versioned store.
- Execute `.run-state.json` is the current source for accepted changed paths. Conformance Review freezes these paths with requirements, code context, planning artifacts, contracts, execute evidence, and target identity.
- Existing review requests use a read-only target snapshot and default-deny inspection policy, then validate structured output before preserving `review.md`. QA needs stronger containment because it may execute explicitly safe commands and must verify no mutation occurred.
- Browser operation memory and SSE events are bounded observations. Durable status is recovered from application use cases and run control, not from operation-hub memory.
- Existing verification locking protects review/smoke product state, while run control protects operational ownership and terminal arbitration. QA must avoid treating both as competing run authorities.

## Constraints

- `PlanningStages()` is closed and ordered; QA mapping, investigation, synthesis, and later repair cannot be represented as planning stages.
- `review`, review JSON shapes, runtime behavior, and canonical `projects/<project>/sprints/<sprint>/review.md` must remain compatible while labels explain Conformance Review.
- Smoke state, `smoke.md`, smoke issues, and verify verdict semantics remain authoritative for their current scope and cannot be populated from QA theories.
- Detailed QA records must not be stored in `FlowState`, run-control snapshots/events, browser operation memory, or runtime metadata.
- All persisted QA paths must be workspace-relative, sprint-contained, symlink-safe, schema-validated, and atomically committed.
- Runtime success is not acceptance; IDs, references, budgets, fingerprints, lifecycle, evidence kinds, and output bounds must validate before state promotion.
- Run control remains generic and cannot import `internal/sprint` or QA semantics.
- Web production code may depend on `internal/app`, not directly on sprint or run-control domain packages.
- Runtime permission support is capability-checked. Unsupported containment must produce a blocked or failed result rather than fall back to prompt-only safety.
- Existing review read-only policy allows inspection tools only. Safe QA command execution requires an explicit bounded policy and post-attempt mutation detection; the current review policy alone is insufficient.
- Governed-input normalization, ordering, path collection, IDs, and fingerprints must be deterministic. Current review changed paths are deduplicated and sorted but derive from both run-state `files` and task evidence.
- Diagnostics and transport payloads are intentionally bounded and sanitized. Raw command output, secrets, unrestricted paths, or arbitrary agent text must not enter durable summaries or cross-surface projections.
- Cancellation and terminal outcomes must be arbitrated through run control, with valid completed QA domain records preserved and stale or partial work excluded from current synthesis.
- Process-local locks and browser operation records cannot become alternate recovery or progress authorities.

## Open Questions

- Whether execute run-state `files` plus task evidence is the final authoritative changed-path contract, especially for deleted paths, renamed paths, absent Git metadata, and evidence paths that are not implementation changes.
- Whether QA state should use the existing product-state database abstraction, a dedicated contained artifact tree, or a combination of database authority and inspectable checkpoints.
- How the legacy verification file lock should coexist with run-control fencing for QA without extending PID-lock ownership into a second lifecycle authority.
- Which runtime and subprocess mechanism can permit only demonstrably non-mutating commands while enforcing filesystem, Git, network, environment, external-system, timeout, cancellation, and output limits.
- Whether the current runtime permission vocabulary can express all required command restrictions or needs a generic platform-level capability extension.
- Which QA summary and pointer fields can be added to `FlowState`, `StatusSummary`, and `SprintSummary` without coupling QA to `PlanningStage` or breaking existing JSON consumers.
- Which human-readable inspectable output is appropriate before canonical `qa.md` is allowed.
- The exact stable-ID namespaces, schema migration window, fingerprint inputs, shard budgets, follow-up scheduling bounds, and synthesis acceptance rules remain requirements-driven decisions for reasoning.
