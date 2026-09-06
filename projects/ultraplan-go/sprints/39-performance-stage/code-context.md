# Sprint Code Context

## Sprint Scope

Sprint 39 adds an explicitly enabled `performance` verification phase between successful execute and Conformance Review. The implementation must parse project activation policy and requirements-owned targets strictly, freeze benchmark and target identities, qualify measurements, perform bounded isolated optimization under product-owned mutation controls, persist immutable evidence plus a bounded flow projection, and expose one durable operation consistently through CLI, TUI, and web surfaces.

The repository has no performance-specific implementation yet. The relevant foundations are the verification-phase abstraction, execute worktree identity, QA/repair persistence, isolation and process boundaries, durable operation ownership, lower-only configuration, and adapter-independent application DTOs.

## Inspected Repository Areas

- Project domain, project-index parsing, validation behavior, and project fixtures.
- Requirements validation, sprint state schemas, phase ordering, execute target worktrees, and verification status projection.
- QA and repair domains, private persistence, writer fencing, immutable evidence, isolated workspaces, product-owned proposal application, reverification, cleanup, and recovery.
- Direct process execution, bounded output, explicit environment, cancellation, descendant cleanup, tree identity, and disposable-copy containment.
- Runtime model configuration and lower-only operational limit conventions.
- Durable application operations, confirmation fingerprints, acceptance/dispatch sequencing, cancellation, terminal arbitration, and adapter-independent DTOs.
- Sprint CLI dispatch and shared TUI/browser application boundaries.
- Workspace requirements scaffolding and generation tests.
- Current architecture plans governing phase separation, isolation, persistence, cancellation, freshness, and bounded convergence.

## Selected Source References

### Project Policy Domain

- **Path:** `internal/project/domain.go`
- **Lines:** `8-47`
- **Symbol:** `CatalogSection`, `ProjectIndex`, `ProjectReasoningPolicy`
- **Rationale:** Shows the typed project-index domain and existing policy precedent. Performance activation belongs beside project reasoning policy without becoming a target source.

- **Path:** `internal/project/index.go`
- **Lines:** `10-77`
- **Symbol:** `recognizedSections`, `ParseProjectIndex`
- **Rationale:** This is the central extension point for recognizing and strictly parsing `Performance Policy`, including duplicate sections, exact keys, invalid modes, and disabled-by-default behavior.

- **Path:** `internal/project/index.go`
- **Lines:** `79-145`
- **Symbol:** `rowMap`, `entryFromRow`, `parseTableRow`, `isSeparatorRow`
- **Rationale:** Defines current Markdown table mechanics and generic catalog-row assumptions. Performance policy parsing must not accidentally accept target-like fields through these generic helpers.

- **Path:** `internal/project/project_test.go`
- **Lines:** `41-91`
- **Symbol:** `TestStatusAndValidateProjectFixture`, `TestValidateFindsMissingFilesMalformedRowsAndEscapes`
- **Rationale:** Establishes project validation fixture conventions and actionable finding expectations that new policy tests should preserve.

### Requirements Authority

- **Path:** `internal/sprint/index.go`
- **Lines:** `10-52`
- **Symbol:** `ValidateRequirementsContent`
- **Rationale:** Current requirements validation checks only content, placeholders, and required headings. It is the required integration point for policy-aware performance-target validation.

- **Path:** `internal/sprint/index.go`
- **Lines:** `142-225`
- **Symbol:** `markdownHeadingPresent`, `parseMarkdownRow`, `separatorRow`, `containsPlaceholder`
- **Rationale:** Existing helpers are line-based and do not account for fenced blocks, blockquotes, comments, or exact table association. These limitations matter for the stricter target declaration contract.

- **Path:** `internal/sprint/sprint_index_test.go`
- **Lines:** `15-71`
- **Symbol:** `TestSprintIndexParseAndValidateAgainstCatalog`, `TestPromptPreviewAndFlowDryRunAreRuntimeFree`
- **Rationale:** Demonstrates deterministic duplicate findings and the invariant that dry-run validation invokes no runtime and changes no state.

### Verification Phase And Ordering

- **Path:** `internal/sprint/verification_phase.go`
- **Lines:** `5-52`
- **Symbol:** `VerificationPhase`, `VerificationPhases`, `ParseVerificationPhase`
- **Rationale:** Defines verification independently of authored planning stages and preserves legacy Conformance Review compatibility. Performance must be additive here rather than represented as a planning stage.

- **Path:** `internal/sprint/domain.go`
- **Lines:** `31-60`
- **Symbol:** `PlanningStage`, `StageState`
- **Rationale:** Captures the legacy planning-stage vocabulary and serialized state. Existing stage identities and ordering must remain compatible.

- **Path:** `internal/sprint/domain.go`
- **Lines:** `136-220`
- **Symbol:** `FlowState`, `VerificationAttempt`, `VerificationStage`
- **Rationale:** Shows the split between bounded flow summaries and detailed verification records, including freshness, evidence references, active attempts, terminal attempts, and next actions.

- **Path:** `internal/sprint/domain.go`
- **Lines:** `222-250`
- **Symbol:** `VerificationStatus`, `StatusSummary`
- **Rationale:** Defines the status projection consumed by application surfaces. Performance status must be additive while detailed samples and profiles remain outside `flow-state.json`.

- **Path:** `internal/sprint/flow.go`
- **Lines:** `91-203`
- **Symbol:** `Service.Flow`
- **Rationale:** Owns stage scheduling and prevents adapters from duplicating ordering policy. It currently runs execute before review/smoke verification and is the primary ordering boundary for enabled performance.

- **Path:** `internal/sprint/flow.go`
- **Lines:** `268-357`
- **Symbol:** `runFlowStage`, `flowStages`, `flowStageAlreadyValid`
- **Rationale:** Encodes the current execute endpoint and readiness checks. Disabled compatibility and enabled performance admission must preserve these semantics.

### Flow Persistence And Migration

- **Path:** `internal/sprint/state.go`
- **Lines:** `20-98`
- **Symbol:** `LoadFlowState`
- **Rationale:** Implements strict loading, schema migration, unknown-field rejection, and database/file fallback. Adding a performance projection requires compatible additive migration.

- **Path:** `internal/sprint/state.go`
- **Lines:** `227-320`
- **Symbol:** `SaveFlowState`, `flowStateCheckpoint`, `saveFlowStateWithHooks`
- **Rationale:** Preserves verification summaries during planning writes and publishes state atomically. Performance state must follow this preservation rule without embedding detailed evidence.

- **Path:** `internal/sprint/state.go`
- **Lines:** `323-435`
- **Symbol:** `ValidateFlowState`, `validateRepairFlowSummary`, `validateQAFlowSummary`
- **Rationale:** Enforces exact stage ordering, contained pointers, valid outcomes, and bounded verification summaries. This is the validation seam for a digest-bound performance summary.

- **Path:** `internal/sprint/state_database.go`
- **Lines:** `12-97`
- **Symbol:** `sprintFlowStateKind`, `loadFlowStateDatabase`, `saveFlowStateDatabase`
- **Rationale:** Shows product-state decomposition and filesystem checkpoint behavior. Performance projection changes must remain compatible with both persistence paths.

### Execute Target Identity

- **Path:** `internal/sprint/domain.go`
- **Lines:** `118-134`
- **Symbol:** `ExecuteTargetRef`, `ExecuteRunState`
- **Rationale:** Defines the recorded target worktree and execute evidence identity from which performance admission and freshness must derive.

- **Path:** `internal/sprint/execute_target.go`
- **Lines:** `30-82`
- **Symbol:** `ResolveExecuteTarget`, `Service.resolveSprintTarget`
- **Rationale:** Resolves the implementation repository into the recorded sprint worktree. Performance must reuse this worktree rather than resolve an alternate target.

- **Path:** `internal/sprint/execute_target.go`
- **Lines:** `102-160`
- **Symbol:** `validateSprintWorkspace`, `inspectSprintWorkspace`
- **Rationale:** Validates that the recorded path is a real worktree in the expected repository and branch, providing source-identity admission checks.

- **Path:** `internal/sprint/execution_handoff.go`
- **Lines:** `18-68`
- **Symbol:** `ExecutionHandoff`, `LoadExecutionHandoff`
- **Rationale:** Projects completed execute tasks and evidence from durable execute state. Performance preparation must reject missing, failed, cancelled, or stale execution evidence.

### Limits And Freshness

- **Path:** `internal/sprint/qa_types.go`
- **Lines:** `70-178`
- **Symbol:** `QARunLifecycle`, `QATerminalResult`, `QABudgets`, `DefaultQABudgets`, `MaximumQABudgets`
- **Rationale:** Provides the established lifecycle, cleanup-uncertainty, finite-default, immutable-maximum, and persisted-counter conventions.

- **Path:** `internal/sprint/qa_types.go`
- **Lines:** `181-236`
- **Symbol:** `QASettings`, `QAFreshness`
- **Rationale:** Demonstrates separate runtime-role selection and digest-bound freshness facts while keeping authority in product code.

- **Path:** `internal/sprint/freshness_policy.go`
- **Lines:** `1-15`
- **Symbol:** `strictCompletedReviewSnapshotFreshness`
- **Rationale:** Records that broad snapshot freshness is disabled for review and smoke. Performance requires independent explicit digest-bound freshness and cannot rely on these switches.

### Private Evidence Persistence

- **Path:** `internal/sprint/qa_state.go`
- **Lines:** `16-55`
- **Symbol:** `QAStore`, `WithWriterFence`, `VerificationBytes`
- **Rationale:** Defines the private verification store, failure-injection hooks, writer fence, no-follow footprint accounting, and hard state bound.

- **Path:** `internal/sprint/qa_state.go`
- **Lines:** `764-806`
- **Symbol:** `checkWriter`, `writeRecord`
- **Rationale:** Enforces current durable ownership and makes immutable records idempotent only when bytes match exactly. Performance evidence requires equivalent semantics.

- **Path:** `internal/sprint/qa_state.go`
- **Lines:** `1018-1037`
- **Symbol:** `writeBytes`
- **Rationale:** Applies immutable byte-record conflict detection and bounded atomic publication, directly relevant to retained proposal patches.

- **Path:** `internal/sprint/qa_state.go`
- **Lines:** `1105-1154`
- **Symbol:** `privateAtomicWrite`
- **Rationale:** Implements private directory permissions, same-directory temporary files, sync, and atomic rename hooks. It is the closest persistence primitive for private performance records.

- **Path:** `internal/sprint/qa_types.go`
- **Lines:** `245-280`
- **Symbol:** `QAWriterToken`, `QARunCorrelation`, `QAArtifactRef`
- **Rationale:** Defines durable run, attempt, and generation fencing plus digest-bound evidence pointers. Performance persistence requires the same stale-writer protection.

### Isolation And Process Execution

- **Path:** `internal/platform/process/process.go`
- **Lines:** `15-56`
- **Symbol:** `Request`, `Result`, `Runner`
- **Rationale:** Defines explicit executable/argv execution, fixed working directory, caller-owned environment, finite timeout/output/cleanup limits, and cleanup facts.

- **Path:** `internal/platform/process/process.go`
- **Lines:** `60-151`
- **Symbol:** `DirectRunner.Run`
- **Rationale:** Implements timeout and cancellation handling, process-tree stop/wait, bounded capture, exit classification, and cleanup certainty.

- **Path:** `internal/platform/process/process_test.go`
- **Lines:** `15-98`
- **Symbol:** `TestDirectRunnerExactEnvironmentCwdAndCapture`, `TestDirectRunnerCancellationCleansOwnedDescendant`, `TestDirectRunnerBoundsAndTimeout`
- **Rationale:** Proves exact environment handling, descendant cleanup, bounded progress, output truncation, and timeout behavior expected from performance commands.

- **Path:** `internal/platform/process/isolation.go`
- **Lines:** `18-76`
- **Symbol:** `IsolationLimits`, `IsolationRequest`, `IsolationCapabilities`, `IsolationWorkspace`
- **Rationale:** Defines finite copy limits, protected roots, host capability facts, content identity, and cleanup results used for fail-closed admission.

- **Path:** `internal/platform/process/isolation.go`
- **Lines:** `78-149`
- **Symbol:** `CreateIsolation`
- **Rationale:** Creates bounded disposable copies, rejects unsafe roots, protects parent placement, verifies copied identity, and removes failed copies.

- **Path:** `internal/platform/process/isolation.go`
- **Lines:** `151-231`
- **Symbol:** `IsolationWorkspace.Resolve`, `IsolationWorkspace.Run`, `IdentifyTree`, `CompareTrees`, `ChangedPaths`
- **Rationale:** Provides contained path resolution, isolated execution, frozen tree identity, and actual-diff derivation for benchmark authoring and optimization.

- **Path:** `internal/platform/process/isolation.go`
- **Lines:** `254-320`
- **Symbol:** `IsolationWorkspace.Cleanup`, `validateIsolationRequest`, `rejectRootOverlap`
- **Rationale:** Defines positive-limit validation, protected-root overlap rejection, and proof-oriented workspace removal.

- **Path:** `internal/platform/process/isolation_test.go`
- **Lines:** `13-115`
- **Symbol:** `TestIsolationCopiesNonGitTreeRunsAndCleans`, `TestIsolationRejectsEscapeSymlinkSpecialHardlinkAndBudgets`
- **Rationale:** Tests identity-preserving copies, changed-path detection, cleanup, escape rejection, hard-link rejection, and resource bounds.

- **Path:** `internal/platform/process/isolation_test.go`
- **Lines:** `140-180`
- **Symbol:** `TestIsolationRejectsOverlappingParentAndCancellation`, `TestIsolationNativeProtectedRootDeny`
- **Rationale:** Covers cancellation, root overlap, and native denial capabilities required for conservative performance admission.

### Product-Owned Mutation Boundary

- **Path:** `internal/sprint/qa_repair.go`
- **Lines:** `25-117`
- **Symbol:** `repairProposalPromptBody`, `RepairPathClass`, `RepairRunRequest`, `ResumeRepair`
- **Rationale:** Establishes that runtimes work only in isolated copies, cannot claim success, and cannot directly mutate production. Resume reconciles durable boundaries without granting new proposal authority.

- **Path:** `internal/sprint/qa_repair.go`
- **Lines:** `221-405`
- **Symbol:** `Service.PrepareRepair`
- **Rationale:** Demonstrates admission under a mutation lease and writer fence, target verification, isolation checks, frozen checks, path restrictions, immutable packet creation, and idempotent preparation.

- **Path:** `internal/sprint/qa_repair.go`
- **Lines:** `479-655`
- **Symbol:** `Service.RunRepair`
- **Rationale:** Shows durable ownership validation, pre-proposal target revalidation, deadline enforcement, isolated runtime execution, actual tree comparison, bounded proposal derivation, and provisional correctness checks.

- **Path:** `internal/sprint/qa_repair.go`
- **Lines:** `657-815`
- **Symbol:** `Service.RunRepair`
- **Rationale:** Covers immutable proposal publication, pre-apply identity checks, journaled product-owned application, scope verification, fixed reverification, cleanup proof, lease release, and product-derived terminal outcome.

- **Path:** `internal/sprint/qa_repair.go`
- **Lines:** `846-960`
- **Symbol:** `Service.RecoverRepair`
- **Rationale:** Provides the recovery precedent: no new authority, digest-bound compensation, target-drift detection, cleanup uncertainty escalation, lease release, and terminal publication.

- **Path:** `internal/sprint/qa_investigator_workspace.go`
- **Lines:** `14-70`
- **Symbol:** `prepareQAInvestigatorWorkspace`, `cleanupQAInvestigatorWorkspaces`
- **Rationale:** Shows deterministic disposable-copy paths, retained-copy identity checks, fail-closed capability admission, and bounded cleanup facts.

### Runtime Configuration

- **Path:** `internal/platform/config/config.go`
- **Lines:** `13-62`
- **Symbol:** `Config`, `Planning`
- **Rationale:** Defines top-level configuration ownership and stage-specific model fields. Performance model configuration must be additive without introducing target or command authority.

- **Path:** `internal/platform/config/qa.go`
- **Lines:** `10-90`
- **Symbol:** `QA`, `QARepair`
- **Rationale:** Documents the lower-only configuration convention, role-specific models, and operational limits that performance settings must follow.

### Durable Application Operations

- **Path:** `internal/app/operations.go`
- **Lines:** `23-77`
- **Symbol:** `OperationalUseCases`, `WebOperations`, `DurableOperationManager`, `OperationReconciler`
- **Rationale:** Defines the adapter-independent operation boundary shared by terminal and browser surfaces, including durable acceptance, cleanup uncertainty, and startup reconciliation.

- **Path:** `internal/app/operations.go`
- **Lines:** `80-170`
- **Symbol:** `OperationKind`, `OperationRequest`, `Confirmation`
- **Rationale:** Central registry for operation kinds and governed confirmation facts. Performance prepare/start/resume/status/cancel/recover operations belong here.

- **Path:** `internal/app/operations.go`
- **Lines:** `238-459`
- **Symbol:** `dashboardUseCases.PrepareOperation`
- **Rationale:** Performs runtime-free preparation, validates operation inputs, declares mutation scope, identifies governed inputs, and binds a deterministic confirmation fingerprint.

- **Path:** `internal/app/operations.go`
- **Lines:** `537-706`
- **Symbol:** `dashboardUseCases.RunOperation`
- **Rationale:** Rechecks prepared fingerprints, routes runtime-free status and dry-run behavior, and delegates runtime-backed work to the shared runner.

- **Path:** `internal/app/operations.go`
- **Lines:** `828-927`
- **Symbol:** `governedOperationInputs`, `fingerprintOperationInputs`
- **Rationale:** Fingerprints governed files and directories without following symlinks. Performance preparation must include policy, requirements, execute evidence, and frozen correctness identities.

- **Path:** `internal/app/operation_runner.go`
- **Lines:** `16-19`
- **Symbol:** `sharedOperationRunner`
- **Rationale:** Establishes one runtime-backed implementation for CLI, TUI, and browser adapters.

- **Path:** `internal/app/operation_runner.go`
- **Lines:** `126-217`
- **Symbol:** `sharedOperationRunner`
- **Rationale:** Shows how QA and repair obtain durable ownership, install writer fencing, emit bounded progress, and route start/resume/recover through product services.

- **Path:** `internal/app/durable_operations.go`
- **Lines:** `109-203`
- **Symbol:** `AcceptOperation`, `DispatchOperation`, `qaOwnershipFromContext`
- **Rationale:** Persists acceptance and owner claim before execution, separates immutable confirmation from dispatch, and derives product-writer authority from the durable fence.

- **Path:** `internal/app/durable_operations.go`
- **Lines:** `206-365`
- **Symbol:** `RecordOperationEvent`, `controlOperation`, `FinishOperation`
- **Rationale:** Coalesces progress, observes canonical cancellation, maintains leases, cancels stale owners, and submits one fenced terminal proposal.

### Application DTOs And Interfaces

- **Path:** `internal/app/sprint_usecases.go`
- **Lines:** `18-93`
- **Symbol:** `SprintSummary`, `QAUseCases`, `RepairUseCases`
- **Rationale:** Defines bounded adapter-independent sprint DTOs and separates query/control interfaces from persistence internals.

- **Path:** `internal/app/sprint_usecases.go`
- **Lines:** `101-167`
- **Symbol:** `RepairStatusResult`, `RepairPacketSummary`, `RepairBudgetSummary`
- **Rationale:** Demonstrates a bounded public status projection with fingerprints, lifecycle, counters, blockers, cleanup, and next action while withholding raw proposal content.

- **Path:** `internal/app/sprint_usecases.go`
- **Lines:** `201-255`
- **Symbol:** `QAResult`
- **Rationale:** Provides the current cross-interface JSON projection pattern for freshness, run correlation, target identity, limits, progress, cancellation, blockers, and evidence pointers.

- **Path:** `internal/app/sprint_usecases.go`
- **Lines:** `634-810`
- **Symbol:** `dashboardUseCases.SprintSummaries`
- **Rationale:** Aggregates sprint, QA, repair, review, smoke, and artifact facts for TUI/web consumers without exposing storage records.

- **Path:** `internal/app/sprint_commands.go`
- **Lines:** `40-94`
- **Symbol:** sprint command dispatch
- **Rationale:** Central CLI dispatch and help routing for sprint subcommands. It currently has no performance command.

### Workspace Guidance

- **Path:** `internal/workspace/scaffold/templates/requirements.md`
- **Lines:** `1-44`
- **Symbol:** requirements scaffold template
- **Rationale:** Current generated requirements contain no optional performance-target guidance. Any addition must remain threshold-neutral.

- **Path:** `internal/workspace/scaffold/prompts/create-requirements.md`
- **Lines:** `1-46`
- **Symbol:** requirements generation prompt
- **Rationale:** Governs model-authored requirements and is the prompt boundary for prohibiting invented numeric thresholds.

- **Path:** `internal/workspace/workspace_test.go`
- **Lines:** `1-180`
- **Symbol:** workspace initialization tests
- **Rationale:** Verifies scaffolded workspace content and is the compatibility surface for exact optional guidance and absence of fabricated targets.

### Governing Architecture Contracts

- **Path:** `docs/plans/post-execution-qa-and-repair-loop.md`
- **Lines:** `208-254`
- **Symbol:** target lifecycle
- **Rationale:** Defines the existing post-execute lifecycle and freshness relationship that performance must extend without collapsing other verification phases.

- **Path:** `docs/plans/post-execution-qa-and-repair-loop.md`
- **Lines:** `447-468`
- **Symbol:** investigation workspace isolation
- **Rationale:** Records isolation expectations reusable by benchmark authoring and optimization attempts.

- **Path:** `docs/plans/post-execution-qa-and-repair-loop.md`
- **Lines:** `540-631`
- **Symbol:** repair model and bounded convergence
- **Rationale:** Establishes one-issue-at-a-time mutation, correctness-first reverification, finite convergence, and explicit non-success outcomes.

- **Path:** `docs/plans/post-execution-qa-and-repair-loop.md`
- **Lines:** `632-698`
- **Symbol:** artifact and state model
- **Rationale:** Defines canonical artifacts, private detailed state, bounded projections, compatibility migration, and fingerprint ownership.

- **Path:** `docs/plans/post-execution-qa-and-repair-loop.md`
- **Lines:** `756-812`
- **Symbol:** internal architecture direction
- **Rationale:** Requires verification phases to remain separate from planning stages, keeps behavior in the sprint module, and identifies reusable foundations.

- **Path:** `docs/plans/server-shutdown-run-cancellation-contract.md`
- **Lines:** `7-41`
- **Symbol:** core rule and ownership boundary
- **Rationale:** Makes durable acceptance and product ownership independent of browser connections or request lifetimes.

- **Path:** `docs/plans/server-shutdown-run-cancellation-contract.md`
- **Lines:** `42-145`
- **Symbol:** graceful shutdown sequence
- **Rationale:** Specifies draining, canonical cancellation, process-tree termination, bounded cleanup, conservative terminal outcomes, and progress completion.

- **Path:** `docs/plans/server-shutdown-run-cancellation-contract.md`
- **Lines:** `169-207`
- **Symbol:** browser disconnection, crash recovery, and concurrency
- **Rationale:** Requires observer disconnects to be non-authoritative and restart reconciliation to arbitrate stale ownership and late completion safely.

- **Path:** `docs/plans/integrated-roadmap.md`
- **Lines:** `113-177`
- **Symbol:** architecture invariants and product persistence boundary
- **Rationale:** Preserves product-owned authority and separation between canonical product state, source repositories, and execution workspaces.

- **Path:** `docs/plans/integrated-roadmap.md`
- **Lines:** `310-355`
- **Symbol:** bounded execution, cancellation, and locking
- **Rationale:** Governs finite runtime work, cooperative cancellation, lock ownership, and conservative interruption behavior.

- **Path:** `docs/plans/retrieval-ready-content-plan.md`
- **Lines:** `120-163`
- **Symbol:** retrieval metadata and product persistence
- **Rationale:** Constrains performance evidence to remain compatible with later content identity without introducing premature generic retrieval storage.

## Relationships

- `project.ParseProjectIndex` supplies project-owned activation policy; requirements validation must combine that policy with the exact target table while retaining requirements as the sole target authority.
- Execute target resolution and execute state establish the canonical sprint worktree and completed execution evidence used by performance admission and freshness.
- `VerificationPhase`, `Service.Flow`, and `VerificationStatus` define phase identity, ordering, gating, and bounded status projection.
- Performance private persistence should follow `QAStore`: immutable attempt records, atomic current pointers, contained paths, writer fencing, retention bounds, and a digest-bound `FlowState` summary.
- `DurableOperationManager` owns acceptance, lease identity, canonical cancellation, and terminal arbitration. `sharedOperationRunner` is the single execution route used by adapters.
- `DirectRunner` owns explicit process execution and cleanup. `IsolationWorkspace` owns disposable copies, source identity, actual-diff inspection, and cleanup proof.
- Repair orchestration demonstrates the product/runtime trust split: the runtime proposes in isolation; product code validates scope, applies changes, reverifies correctness, derives outcomes, and reconciles interruptions.
- Application DTOs project bounded facts to CLI, TUI, and web. Raw samples, profiles, prompts, provider payloads, and proposal bytes remain private evidence.
- Configuration supplies model selection and lower-only operational limits but cannot alter targets, commands, parsers, sample policy, correctness gates, or verdict rules.

## Constraints

- Missing project performance policy must behave as `disabled`; existing project indexes and disabled flows cannot require user migration.
- Project policy controls activation only. Targets cannot enter through project policy, configuration, environment, CLI, runtime output, or prior artifacts.
- The target parser needs stricter Markdown awareness than existing helpers provide, including fences, blockquotes, comments, duplicate headings, malformed tables, and table association.
- Target normalization, numeric comparison, qualification, and outcomes are product logic, not runtime judgments.
- Performance is a `VerificationPhase`, not a `PlanningStage`; legacy stage parsing and state ordering must remain compatible.
- `flow-state.json` is bounded. Detailed attempts, samples, profiles, patches, and command evidence belong under private performance verification storage.
- Every private write must be contained, no-follow, atomic, versioned, digest-bound, and rejected under a stale writer fence.
- Durable acceptance and owner claim precede runtime or command execution.
- Commands require explicit argv, fixed working directory, allowlisted environment, positive timeout, bounded output, cancellation, and process-tree cleanup.
- Benchmark authoring and optimization run only in disposable copies. Missing identity, isolation, target denial, descendant cleanup, or workspace-removal proof is blocking.
- Promotion must use product-owned path classification, actual diff inspection, identity checks, overlap detection, mutation leasing, and correctness reverification.
- Cleanup uncertainty, target drift, benchmark drift, parser uncertainty, environment drift, noisy measurements, or non-finite evidence cannot produce success.
- Defaults and maxima are finite and product-owned. Workspace and environment settings can only lower operational limits.
- Performance freshness must be explicit and independent of the relaxed review/smoke snapshot switches.
- Interface adapters consume application DTOs and use cases only; all presentations must share the same product-derived outcome.
- Browser refresh, disconnect, session loss, or SSE loss cannot cancel or complete a run.

## Open Questions

- No existing abstraction covers benchmark descriptors, parser dispatch, sample qualification, or numeric target comparison; these are new sprint-domain concerns.
- Durable writer tokens are named `QAWriterToken` despite repair reuse; the source does not establish whether performance should reuse that name or introduce a verification-generic equivalent.
- The project-index parser does not track fenced regions or duplicate recognized headings, so strict performance-policy parsing cannot rely on its current loop unchanged.
- Requirements content validation does not currently receive project policy; the source does not establish the eventual service boundary for policy-aware validation.
- Review and smoke freshness is intentionally relaxed, while Sprint 39 requires strict performance fingerprints; performance freshness must remain independently enforced.
