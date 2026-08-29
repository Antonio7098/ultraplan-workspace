# Sprint Code Context

## Sprint Scope

Sprint 38 adds a `VerificationPhase` repair capability after evidence-producing QA. It must freeze exactly one current adjudicated issue into an immutable packet, bind explicit confirmation to the packet and target, generate a proposal in an isolated workspace, permit only a product-owned bounded production mutation, progressively reverify the result, clean up conclusively, and publish exactly one truthful terminal outcome.

Manual repair is the first mutable path and allows one cycle. Automatic repair remains opt-in, depends on current manual proof, reuses the same mutation protocol, and is constrained by lower-only limits, progress checks, and explicit stop conditions. CLI, TUI, HTTP, durable operations, private verification state, and bounded `flow-state.json` projections must expose the same canonical facts.

The repository currently implements Sprint 37 QA, isolation, durable operation ownership, writer fencing, cancellation, recovery, and interface projections. The required repair service, repair schemas, repair persistence, repair operations, and repair-specific interface controls do not yet exist.

## Inspected Repository Areas

- `internal/sprint`: verification phases, QA mapping and execution, evidence plans, adjudication, assessment, strict persistence, target freshness, mutation leases, recovery, and cleanup uncertainty.
- `internal/platform/process`: bounded tree identity, disposable workspace creation, link rejection, native protected-root isolation, process-tree cancellation, output limits, and cleanup.
- `internal/platform/config`: QA defaults, immutable maxima, environment overrides, precedence, validation, and effective-source tracking.
- `internal/app`: adapter-independent QA DTOs, operation preparation and stale-input fingerprints, durable acceptance, writer fences, shared runtime dispatch, CLI parsing, cancellation, and public error handling.
- `internal/tui`: QA routes, operation actions, bounded rendering, refresh, background-run behavior, and durable-event display.
- `internal/web`: QA query routes, guarded operation preparation/start, CSRF and session controls, SSE replay, bounded projections, cancellation, shutdown draining, and no-JavaScript rendering.
- Tests covering QA admission and adjudication, isolation, persistence atomicity, stale writers, lower-only configuration, durable ownership, CLI schemas, TUI parity, HTTP escaping, cancellation, and shutdown uncertainty.
- Governed roadmap, QA/repair loop, and server shutdown plans that constrain Sprint 38 behavior.

## Selected Source References

### Governed Repair Semantics

- **Path:** `docs/plans/integrated-roadmap.md`
- **Lines:** `626-688`
- **Symbol:** `Sprint 38 - Manual repair and bounded automatic repair`
- **Rationale:** Defines the repository-level sequencing, frozen packet contents, progressive reverification order, bounded automatic stop conditions, closed outcomes, and Sprint 38 exit criteria.

- **Path:** `docs/plans/post-execution-qa-and-repair-loop.md`
- **Lines:** `182-207`
- **Symbol:** `Evidence-backed issue and Repair responsibilities`
- **Rationale:** Establishes that only adjudicated current evidence can authorize repair and identifies the minimum issue-packet facts and repair authority boundary.

- **Path:** `docs/plans/post-execution-qa-and-repair-loop.md`
- **Lines:** `515-629`
- **Symbol:** `Issue promotion, Repair model, and Bounded convergence loop`
- **Rationale:** Records existing intended issue promotion fields, repair eligibility, agent restrictions, widening verification order, convergence rules, cycle evidence, and terminal outcome vocabulary.

- **Path:** `docs/plans/post-execution-qa-and-repair-loop.md`
- **Lines:** `670-696`
- **Symbol:** `Detailed state ownership and Fingerprints`
- **Rationale:** Separates bounded flow summaries from detailed verification state and lists identities that must govern freshness and retained evidence after mutation.

- **Path:** `docs/plans/server-shutdown-run-cancellation-contract.md`
- **Lines:** `42-207`
- **Symbol:** `Graceful shutdown, recovery, and concurrency contract`
- **Rationale:** Governs draining, canonical cancellation, process-tree cleanup, terminal arbitration, browser-disconnect independence, restart reconciliation, lock release, and bounded summary state for repair operations.

### Verification And QA Domain

- **Path:** `internal/sprint/verification_phase.go`
- **Lines:** `5-52`
- **Symbol:** `VerificationPhase`
- **Rationale:** Already declares `repair` beside QA and Conformance Review while explicitly preventing QA and repair from entering `PlanningStage` order.

- **Path:** `internal/sprint/qa_types.go`
- **Lines:** `16-275`
- **Symbol:** `QABudgets`, `QAFreshness`, `QAWriterToken`, `QAState`, `QAFlowSummary`
- **Rationale:** Shows the existing schema-versioning, lower-bounded runtime policy, freshness fingerprints, durable writer correlation, detailed QA pointer state, and bounded flow projection that repair should parallel or extend without conflating phases.

- **Path:** `internal/sprint/qa_types.go`
- **Lines:** `277-335`
- **Symbol:** `QATargetIdentity`, `QAMap`
- **Rationale:** Defines the deterministic QA map’s target identity, governed fingerprints, coverage ownership, budgets, shards, and input references from which repair admission must derive current scope.

- **Path:** `internal/sprint/qa_types.go`
- **Lines:** `563-629`
- **Symbol:** `ValidateQASettings`, `validateQABudgets`
- **Rationale:** Demonstrates product-side validation of positive finite limits against immutable maxima, including exact semantic-call invariants.

- **Path:** `internal/sprint/qa_types.go`
- **Lines:** `823-903`
- **Symbol:** `QAError`, `QAErrorCategory`, `qaRecovery`
- **Rationale:** Provides typed domain error and recovery conventions for unknown schemas, stale inputs, conflicts, persistence failures, unavailable runtimes, admission failures, and cleanup uncertainty.

- **Path:** `internal/sprint/qa_evidence.go`
- **Lines:** `16-204`
- **Symbol:** `QAEvidencePlan`, `QAEvidenceRecord`, `QAIssue`, `QAAdjudication`, `QAAssessmentRecord`, `QAAdmission`
- **Rationale:** Contains the authoritative current QA evidence and issue schemas. It reveals which packet facts already exist and which do not: `QAIssue` carries identity, class, severity, location, evidence IDs, promotion reason, and eligibility, while exact reproducers, allowed/forbidden paths, containing checks, and acceptance criteria are not currently issue fields.

- **Path:** `internal/sprint/qa_evidence.go`
- **Lines:** `206-294`
- **Symbol:** `ValidateQAEvidencePlan`, `ValidateQAEvidence`, `qaPathApproved`
- **Rationale:** Defines current strict plan/evidence freshness, path containment, execution bounds, cleanup proof, target identity, patch-reference, and model-observation checks. These are key admission and reverification invariants.

- **Path:** `internal/sprint/qa_adjudication.go`
- **Lines:** `34-165`
- **Symbol:** `AdjudicateQA`
- **Rationale:** Implements pure, runtime-free promotion from frozen plans and accepted evidence, rejects stale, blocked, invalid, or insufficient evidence, groups root causes, and deterministically creates repair-eligible issues.

- **Path:** `internal/sprint/qa_adjudication.go`
- **Lines:** `191-227`
- **Symbol:** `DeriveQAAssessment`
- **Rationale:** Shows current precedence between Conformance Review, rejected evidence, containing smoke, promoted issues, and overall QA assessment. Repair must not bypass or reinterpret this independent authority.

- **Path:** `internal/sprint/qa.go`
- **Lines:** `59-174`
- **Symbol:** `QAStatus`, `QAEvidence`, `QAAdjudication`, `QAAssessment`
- **Rationale:** Establishes side-effect-free canonical reads through digest-validated pointers and current-attempt ownership, which repair preparation and status queries need to preserve.

- **Path:** `internal/sprint/qa.go`
- **Lines:** `176-248`
- **Symbol:** `RecoverQA`
- **Rationale:** Demonstrates runtime-free reconciliation of abandoned phases, stale governed inputs, flow-summary pointers, and retained attempts without inferring success.

- **Path:** `internal/sprint/qa.go`
- **Lines:** `278-439`
- **Symbol:** `RunQA`
- **Rationale:** Shows the current single mutation lease, writer-fenced publication, bounded timeout, resume preparation, terminal failure publication, evidence production, and final state transition used by the shared QA protocol.

- **Path:** `internal/sprint/qa.go`
- **Lines:** `442-551`
- **Symbol:** `buildQAEvidencePublication`
- **Rationale:** Is the QA-to-repair admission boundary today: it requires current review and containing smoke, proven isolation, complete mapping, target identity, frozen evidence plans, isolated execution, adjudication, assessment, and canonical report publication.

- **Path:** `internal/sprint/qa.go`
- **Lines:** `577-713`
- **Symbol:** `prepareQAAttempt`, `runQAShardBatch`
- **Rationale:** Provides resume-boundary and bounded worker behavior, including reuse of terminal evidence, cancellation propagation, per-result writer checks, and publication-failure shutdown.

- **Path:** `internal/sprint/qa_prompt.go`
- **Lines:** `16-125`
- **Symbol:** `QACheckDescriptor`, `ApprovedQAChecks`, `validateQACheckDescriptor`
- **Rationale:** Defines product-owned immutable command descriptors and rejects shell interpreters, Git, write flags, path escapes, unsafe arguments, excessive time/output, and unapproved environment names.

- **Path:** `internal/sprint/qa_prompt.go`
- **Lines:** `136-191`
- **Symbol:** `RenderQAInvestigatorPrompt`, `QAInvestigatorRequest`
- **Rationale:** Shows the current default-deny runtime request and path-level read authority. A repair runtime must remain bounded but requires a distinct isolated-copy write policy rather than production authority.

- **Path:** `internal/sprint/qa_prompt.go`
- **Lines:** `243-285`
- **Symbol:** `RunApprovedQACheck`
- **Rationale:** Demonstrates descriptor ownership checks, environment filtering, process limits, target identity comparison, digest-only output retention, and failure on truncation.

### Persistence And Authority

- **Path:** `internal/sprint/qa_state.go`
- **Lines:** `44-180`
- **Symbol:** `VerificationBytes`, `PruneAttempts`, QA artifact path functions
- **Rationale:** Defines symlink-safe retained-state accounting, bounded attempt retention, and the current verification artifact layout that repair records will extend.

- **Path:** `internal/sprint/qa_state.go`
- **Lines:** `182-264`
- **Symbol:** `QAStore.resolve`, `QAStore.LoadState`
- **Rationale:** Enforces selected-sprint containment, rejects symlink components, validates state scope, and verifies every digest-bound state pointer before returning canonical facts.

- **Path:** `internal/sprint/qa_state.go`
- **Lines:** `321-430`
- **Symbol:** `LoadEvidence`, `LoadAdjudication`, `LoadAssessment`, `readStrictVersion`
- **Rationale:** Shows strict identity-to-path checks, private `0600` enforcement, bounded reads, schema rejection, unknown-field rejection, and single-value JSON decoding.

- **Path:** `internal/sprint/qa_state.go`
- **Lines:** `433-648`
- **Symbol:** `QAPublication`, `QAStore.Publish`, `checkWriter`, `writeRecord`
- **Rationale:** Implements writer checks before each publication, immutable map semantics, pointer-last state and flow publication, canonical-file snapshots, rollback on failure, and atomic records.

- **Path:** `internal/sprint/qa_state.go`
- **Lines:** `651-801`
- **Symbol:** `publishEvidence`
- **Rationale:** Publishes immutable plans, patches, evidence, adjudication, issues, and assessment with bounded counts, digest-bound patch references, current-attempt ownership, repeated writer fencing, and canonical report linkage.

- **Path:** `internal/sprint/qa_state.go`
- **Lines:** `804-1007`
- **Symbol:** `writeBytes`, `privateAtomicWrite`, `restoreQACanonicalFiles`, `verifyReference`, `qaFlowSummary`
- **Rationale:** Provides immutable byte publication, private directories/files, fsync-and-rename atomicity, rollback mechanics, digest verification, and the intentionally bounded `flow-state.json` projection.

- **Path:** `internal/sprint/state.go`
- **Lines:** `201-230`
- **Symbol:** `SaveFlowState`
- **Rationale:** Preserves existing review, smoke, and QA summaries during unrelated planning-state refreshes, preventing later writers from erasing verification authority.

- **Path:** `internal/sprint/state.go`
- **Lines:** `294-383`
- **Symbol:** `ValidateFlowState`, `validateQAFlowSummary`
- **Rationale:** Enforces closed planning state, safe contained paths, closed latest outcomes, and a digest-bound QA summary pointer. Repair projection will require equivalent bounded validation without becoming a detailed database.

- **Path:** `internal/sprint/service.go`
- **Lines:** `22-43`
- **Symbol:** `Service`
- **Rationale:** Identifies current dependency seams for runtime, verification runtime, QA policy, writer/map fences, process runner, smoke settings, mutation tracking, and persistence.

- **Path:** `internal/sprint/service.go`
- **Lines:** `75-109`
- **Symbol:** `NewService`, `acquireMutation`
- **Rationale:** Implements the in-process and cross-process single-writer mutation boundary used by flow, execute, review, smoke, verify, and QA.

- **Path:** `internal/sprint/service.go`
- **Lines:** `149-189`
- **Symbol:** `WithVerificationRuntime`, `WithQASettings`, `WithQAWriterFence`, `WithQAMapFence`
- **Rationale:** Shows how validated verification policy and durable ownership checks are injected without constructing runtimes or mutating state.

- **Path:** `internal/sprint/locks.go`
- **Lines:** `13-143`
- **Symbol:** `ReconcileInterruptedMutation`, `acquireMutationContext`
- **Rationale:** Defines shared nested mutation leases, dead-owner reconciliation, interrupted QA publication, cleanup-uncertainty handling, and the requirement that live leases are never rewritten.

- **Path:** `internal/sprint/verification_lock.go`
- **Lines:** `14-101`
- **Symbol:** `verificationFileLock`
- **Rationale:** Provides cross-process exclusive creation, dead-PID replacement, and ownership-checked release. Repair must not release another owner’s lock.

- **Path:** `internal/sprint/cleanup_uncertain.go`
- **Lines:** `15-82`
- **Symbol:** `CleanupUncertainRecord`, `RecordCleanupUncertain`
- **Rationale:** Records shutdown cleanup uncertainty separately without acquiring or overwriting state still owned by a potentially live operation.

### Isolation, Identity, And Mutation Boundaries

- **Path:** `internal/platform/process/isolation.go`
- **Lines:** `18-124`
- **Symbol:** `IsolationLimits`, `IsolationCapabilities`, `CreateIsolation`
- **Rationale:** Supplies the bounded, non-Git disposable-copy primitive and rejects links, special files, hard links, root overlap, and unbounded copy requests.

- **Path:** `internal/platform/process/isolation.go`
- **Lines:** `126-232`
- **Symbol:** `IsolationWorkspace.Resolve`, `Run`, `Identity`, `CompareTrees`, `Cleanup`
- **Rationale:** Defines contained path resolution, native isolation wrapping, before/after tree identity, actual changed-path comparison, and verified workspace removal.

- **Path:** `internal/platform/process/isolation.go`
- **Lines:** `287-419`
- **Symbol:** `copyBoundedTree`, `collectTreeIdentity`, `readRegularFile`
- **Rationale:** Implements content-and-path tree digests, file/byte limits, link rejection, race detection during reads, and mode-sensitive changed-path comparison.

- **Path:** `internal/platform/process/isolation_linux.go`
- **Lines:** `19-59`
- **Symbol:** `isolationCapabilities`, `nativeIsolationRequest`
- **Rationale:** Shows Linux protected-root denial depends on a successful Bubblewrap capability probe and executes with the disposable workspace as the only writable bind.

- **Path:** `internal/platform/process/process.go`
- **Lines:** `15-151`
- **Symbol:** `Request`, `Result`, `DirectRunner.Run`
- **Rationale:** Provides bounded command execution, process-group ownership, cancellation and timeout propagation, process-tree cleanup, output truncation facts, and explicit cleanup completeness.

- **Path:** `internal/sprint/qa_investigation.go`
- **Lines:** `15-140`
- **Symbol:** `RunQAInvestigation`
- **Rationale:** Integrates target identity, isolation capability checks, protected-root leak rejection, immutable command execution, actual-change comparison, cleanup proof, and target drift into one fail-closed evidence run.

- **Path:** `internal/sprint/qa_investigation.go`
- **Lines:** `142-170`
- **Symbol:** `FreezeQAEvidencePlan`
- **Rationale:** Demonstrates normalized, deterministic, runtime-free freezing of command, condition, path, budget, and fingerprint authority before execution.

- **Path:** `internal/sprint/verify.go`
- **Lines:** `142-249`
- **Symbol:** `VerificationStatus`
- **Rationale:** Derives current Conformance Review and smoke freshness from canonical state, artifact digests, governed fingerprints, containing evidence, and expired-attempt reconciliation without treating stale evidence as current.

- **Path:** `internal/sprint/verify.go`
- **Lines:** `330-368`
- **Symbol:** `refreshEvidenceFingerprint`, `targetIdentity`
- **Rationale:** Shows current governed evidence fingerprint composition and target-root symlink/Git-root checks. Repair requires repeated full target identity checks around proposal, apply, and reverification.

- **Path:** `internal/sprint/execute_target.go`
- **Lines:** `29-82`
- **Symbol:** `ResolveExecuteTarget`, `resolveSprintTarget`
- **Rationale:** Resolves the approved implementation repository or recorded sprint worktree and detects target/workspace changes. This is the current source of production target authority.

- **Path:** `internal/sprint/execute_target.go`
- **Lines:** `99-180`
- **Symbol:** `validateSprintWorkspace`, `createSprintWorkspace`
- **Rationale:** Documents the existing Git worktree behavior and clean-source assumptions. Repair requirements explicitly prohibit Git mutation, so this existing creation path cannot itself be treated as the repair apply protocol.

- **Path:** `internal/sprint/execute_target.go`
- **Lines:** `202-226`
- **Symbol:** `ValidateExecuteWorkdir`, `ExecuteSafetyInstructions`
- **Rationale:** Provides the current approved-target containment check and Git/evidence mutation prohibitions, while also showing that execute runtime currently works directly inside the target and is therefore not an appropriate repair mutation boundary unchanged.

### Configuration

- **Path:** `internal/platform/config/qa.go`
- **Lines:** `10-118`
- **Symbol:** `QA`, `DefaultQA`, `maxQA`, `qaConfigFields`, `qaEnvOverrides`
- **Rationale:** Establishes the existing pattern for safe defaults, immutable maxima, complete field registration, and deterministic `ULTRAPLAN_QA_*` environment names.

- **Path:** `internal/platform/config/qa.go`
- **Lines:** `120-261`
- **Symbol:** `setQAField`, `validateQA`
- **Rationale:** Shows strict field parsing and lower-only validation for integer and duration limits. Repair limits must be complete across defaults, workspace config, environment, and effective-source projection.

- **Path:** `internal/platform/config/config_test.go`
- **Lines:** `264-376`
- **Symbol:** `TestQAConfigFieldsHaveEffectiveSourcesAndLowerOnlyBounds` and related QA configuration tests
- **Rationale:** Provides the current exhaustive test style for maxima, zero/negative values, malformed environment input, unknown fields, and workspace/environment precedence.

### Shared Application And Durable Operations

- **Path:** `internal/app/sprint_usecases.go`
- **Lines:** `69-218`
- **Symbol:** `QAUseCases`, `QAQueries`, `QAEvidenceQueries`, `QARequest`, `QAResult`, `QALimitsSummary`
- **Rationale:** Defines the adapter-independent bounded projection pattern and separates private verification records from public facts. Repair interfaces should extend this boundary rather than expose files.

- **Path:** `internal/app/sprint_usecases.go`
- **Lines:** `366-430`
- **Symbol:** `QAEvidenceResult`, `QAAdjudicationResult`, `QAIssueSummary`, `QAIssuePage`, `QAAssessmentResult`
- **Rationale:** Shows current redacted and bounded issue/evidence DTOs, including cursor pagination and the absence of repair packet, cycle, confirmation, and result DTOs.

- **Path:** `internal/app/sprint_usecases.go`
- **Lines:** `795-955`
- **Symbol:** `QAEvidence`, `QAAdjudication`, `QAIssues`, `QAIssue`, `QAAssessment`
- **Rationale:** Implements current-attempt issue ownership, bounded pagination, stale cursors, focused-query allowlists, and hostile-text sanitization.

- **Path:** `internal/app/sprint_usecases.go`
- **Lines:** `957-1029`
- **Symbol:** `RunQA`, `ResumeQA`, `CancelQA`, `RecoverQA`, `qaMapProjection`
- **Rationale:** Demonstrates shared runner use, canonical cancellation through durable run control, recovery through the sprint service, and status reload after execution rather than deriving authority from events.

- **Path:** `internal/app/operations.go`
- **Lines:** `23-209`
- **Symbol:** `OperationalUseCases`, `DurableOperationManager`, `OperationKind`, `OperationRequest`, `Confirmation`, `OperationResult`
- **Rationale:** Defines the common CLI/TUI/web operation contract, durable acceptance capability, confirmation facts, canonical request carrier, mutation class, progress projection, and stable result/error shapes.

- **Path:** `internal/app/operations.go`
- **Lines:** `211-398`
- **Symbol:** `PrepareOperation`, `validateQAOperationRequest`
- **Rationale:** Implements current operation preview, scope and warning facts, caller-control allowlisting, mutation classification, canonical request generation, governed inputs, and stale-input fingerprint issuance.

- **Path:** `internal/app/operations.go`
- **Lines:** `400-553`
- **Symbol:** `RunOperation`
- **Rationale:** Re-prepares immediately before execution, rejects fingerprint drift, keeps read-only queries separate, and delegates all runtime mutations to one shared runner.

- **Path:** `internal/app/operations.go`
- **Lines:** `610-754`
- **Symbol:** `canonicalOperationRequest`, `operationPrerequisites`, `governedOperationInputs`, `fingerprintOperationInputs`
- **Rationale:** Shows how canonical operation requests and symlink-safe governed file contents are bound into server-issued confirmation freshness.

- **Path:** `internal/app/durable_operations.go`
- **Lines:** `93-156`
- **Symbol:** `AcceptOperation`, `qaOwnershipFromContext`
- **Rationale:** Persists acceptance and owner claim before execution, deduplicates confirmation digests, creates a durable fence, and exposes a cancellation-safe writer check that heartbeats ownership.

- **Path:** `internal/app/durable_operations.go`
- **Lines:** `158-300`
- **Symbol:** `RecordOperationEvent`, `controlOperation`, `FinishOperation`
- **Rationale:** Provides bounded/coalesced durable progress, canonical cancellation acknowledgement, lease heartbeat, reconciliation, and exactly-one terminal proposal through run control.

- **Path:** `internal/app/operation_runner.go`
- **Lines:** `15-169`
- **Symbol:** `sharedOperationRunner`
- **Rationale:** Is the single runtime-backed dispatch path shared by terminal and browser adapters. The QA branch demonstrates obtaining durable ownership, installing writer fencing, invoking the sprint service, and projecting bounded progress.

### CLI, TUI, And Web

- **Path:** `internal/app/sprint_commands.go`
- **Lines:** `369-493`
- **Symbol:** `sprint qa command dispatch`
- **Rationale:** Shows current CLI status, map, recovery, cancellation, durable run/resume acceptance, writer fencing, JSON envelope, and non-zero assessment behavior. There is no repair command yet.

- **Path:** `internal/app/sprint_commands.go`
- **Lines:** `608-725`
- **Symbol:** `sprintQACommand`, `parseSprintQAArgs`, `renderSprintQA`
- **Rationale:** Demonstrates strict public flag allowlisting, explicit `--yes` for external smoke execution, action-specific validation, and deterministic text rendering.

- **Path:** `internal/app/sprint_commands.go`
- **Lines:** `782-869`
- **Symbol:** `qaSettings`
- **Rationale:** Converts effective configuration into domain budgets while preserving every effective source and validating the final policy before runtime use.

- **Path:** `internal/tui/model.go`
- **Lines:** `27-131`
- **Symbol:** `RouteKind`, `Route`, `Model`
- **Rationale:** Defines typed QA routes and shared operation state without direct sprint-domain or private-file dependencies.

- **Path:** `internal/tui/model.go`
- **Lines:** `164-239`
- **Symbol:** `Model.Update`
- **Rationale:** Shows stale-route protection, separate confirmation state, bounded observer events, and terminal operation handling.

- **Path:** `internal/tui/model.go`
- **Lines:** `450-533`
- **Symbol:** `Model.navItems`
- **Rationale:** Maps QA status, start, resume, recovery, durable runs, shards, and theories to shared app operations. Repair routes and guarded actions are absent.

- **Path:** `internal/tui/qa_view.go`
- **Lines:** `8-54`
- **Symbol:** `renderSprintQAView`
- **Rationale:** Renders only bounded app DTOs and does not derive product outcomes from progress events, establishing the required presentation boundary for repair.

- **Path:** `internal/web/routes.go`
- **Lines:** `325-382`
- **Symbol:** `matchRoute`
- **Rationale:** Lists existing HTML and JSON QA resources and guarded operation endpoints. Repair packet, cycle, status, result, confirmation, resume, cancel, and recover resources are not yet routed.

- **Path:** `internal/web/qa_handlers.go`
- **Lines:** `10-201`
- **Symbol:** `QA query handlers`
- **Rationale:** Demonstrates method-separated app-query handlers, bounded issue pagination, additive evidence-query capability checks, stable error mapping, and server-rendered HTML using app DTOs only.

- **Path:** `internal/web/operations.go`
- **Lines:** `18-218`
- **Symbol:** `operationHub.startConfirmed`
- **Rationale:** Enforces operation capacity, draining, session deduplication, confirmation callback execution, durable acceptance before goroutine creation, and canonical operation context ownership.

- **Path:** `internal/web/operations.go`
- **Lines:** `221-425`
- **Symbol:** `operationHub.run`, `finish`, `cancelOperation`, `subscribe`
- **Rationale:** Implements terminal arbitration, durable event-before-observer projection, idempotent cancellation, bounded replay, replay-gap recovery, and browser subscription independence.

- **Path:** `internal/web/operations.go`
- **Lines:** `479-557`
- **Symbol:** `drainAndWait`, `persistCleanupUncertain`, `markCleanupUncertain`
- **Rationale:** Preserves the server shutdown contract by entering drain mode, cancelling every active operation, waiting boundedly, persisting uncertainty before terminal projection, and closing streams.

- **Path:** `internal/web/operations.go`
- **Lines:** `577-657`
- **Symbol:** `projectOperationResult`, `terminalOperationState`, `safeProjectedText`
- **Rationale:** Bounds and redacts retained public output and defines the web operation terminal-state vocabulary that repair outcomes must coexist with without being conflated.

- **Path:** `internal/web/server.go`
- **Lines:** `125-152`
- **Symbol:** `Serve shutdown branch`
- **Rationale:** Connects root-context cancellation to operation draining, bounded cleanup, canonical operation cancellation, HTTP shutdown, and process exit.

### Decision-Relevant Tests

- **Path:** `internal/sprint/qa_adjudication_test.go`
- **Lines:** `9-80`
- **Symbol:** `QA adjudication and admission tests`
- **Rationale:** Proves promotion requires contained current sufficient evidence, incomplete cleanup prevents promotion, deterministic fact failures can be admitted, and unproven isolation fails closed.

- **Path:** `internal/sprint/qa_investigation_test.go`
- **Lines:** `14-105`
- **Symbol:** `QA isolated investigation tests`
- **Rationale:** Verifies the production target remains unchanged, isolated workspaces are removed, retained evidence validates, and original target paths cannot leak into child requests.

- **Path:** `internal/sprint/qa_state_test.go`
- **Lines:** `120-270`
- **Symbol:** `QA private persistence, strict load, stale writer, and atomicity tests`
- **Rationale:** Covers `0600`/`0700` permissions, pointer publication, strict schema decoding, symlink rejection, stale-writer rejection, partial rename failure, and immutable record conflict.

- **Path:** `internal/sprint/qa_state_test.go`
- **Lines:** `272-352`
- **Symbol:** `TestQAEvidencePublicationLoadsAndRollsBackCanonicalFiles`
- **Rationale:** Exercises immutable evidence and patch publication, digest-bound references, state/report loading, injected canonical publication failure, and rollback preservation.

- **Path:** `internal/app/durable_operations_test.go`
- **Lines:** `11-124`
- **Symbol:** `Durable operation ownership tests`
- **Rationale:** Verifies acceptance and owner claim precede execution, current fences remain valid during cancellation cleanup, stale generations fail, progress is durably sanitized, terminal state is committed, duplicate confirmation is deduplicated, and persistence fails closed.

- **Path:** `internal/app/sprint_commands_test.go`
- **Lines:** `149-211`
- **Symbol:** `QA JSON and exit-code tests`
- **Rationale:** Freezes the single-document versioned JSON envelope, stable public error categories, no trailing output, and distinct runtime, validation, and cancellation exit classes.

- **Path:** `internal/app/sprint_commands_test.go`
- **Lines:** `637-741`
- **Symbol:** `QA CLI controls and effective-settings tests`
- **Rationale:** Tests strict flag allowlisting, required explicit confirmation, unsafe control rejection, help coverage, dedicated runtime policy, complete source projection, and default budget parity.

- **Path:** `internal/web/operations_test.go`
- **Lines:** `124-207`
- **Symbol:** `Operation cancellation and shutdown tests`
- **Rationale:** Covers session ownership, idempotent cancellation, draining rejection, bounded shutdown, durable cleanup-uncertainty recording, and terminal projection.

- **Path:** `internal/web/operations_test.go`
- **Lines:** `302-426`
- **Symbol:** `Operation HTTP prepare/start/SSE/cancel and security tests`
- **Rationale:** Proves separate preparation and start, confirmation tokens, CSRF enforcement, session isolation, SSE terminal replay, cancellation, duplicate-start replay, strict JSON, and redaction.

- **Path:** `internal/web/qa_handlers_test.go`
- **Lines:** `38-109`
- **Symbol:** `QA HTTP query tests`
- **Rationale:** Verifies stable public errors, bounded JSON, hostile-text escaping, complete no-JavaScript HTML, and read-only HTTP methods.

## Relationships

- Current QA admission flows from `QAMap` and `VerificationStatus` through frozen `QAEvidencePlan` records, isolated investigations, pure adjudication, `QAIssue` promotion, `QAAssessmentRecord`, and writer-fenced `QAStore` publication.
- `QAStore` publishes detailed immutable records under `verification/`, then atomically publishes current state and a digest-bound bounded `flow-state.json` summary.
- Runtime-backed work is accepted and claimed by `durableOperationManager`; its operation context carries the run-control fence used by `QAStore` before every publication.
- CLI and web runtime execution converge on sprint services through durable operation ownership. TUI and browser actions use app operations and consume bounded app DTOs rather than private files.
- `Service.acquireMutationContext` and the verification file lock serialize sprint mutations. Durable fencing separately prevents cancelled or stale owners from publishing.
- `RunQAInvestigation` composes the process isolation and process-runner packages. Product code freezes commands and paths; platform code supplies copy, identity, process-tree, and cleanup facts without deciding QA outcomes.
- Conformance Review and smoke remain independently authoritative through `VerificationStatus` and assessment derivation. Repair can request delta review and containing checks but cannot write their artifacts or promote a narrow check to global success.
- Web shutdown drains the operation hub, requests canonical cancellation, waits for cleanup, and records uncertainty before HTTP exit. Browser subscription loss only removes an observer and does not cancel the durable operation.

## Constraints

- Repair remains a `VerificationPhase`; it must not become a planning stage or alter planning order.
- The current authoritative issue is selected through the current QA state’s digest-bound adjudication pointer. Historical or caller-supplied issue identity alone is insufficient.
- Current `QAIssue` records do not contain all required repair packet facts. Relevant expectation, command, path, theory, shard, smoke, review, and fingerprint facts are distributed across the map, plans, evidence, adjudication, assessment, and verification status.
- Current QA evidence paths are approved for isolated investigation, not automatically for production mutation. Repair scope needs stricter production-path classification and protected-file denial.
- Existing execute work grants a runtime direct target workdir authority. Repair requirements instead demand isolated proposal generation followed by a product-owned apply boundary; execute cannot be reused unchanged as that boundary.
- Platform isolation supplies facts and fail-closed mechanisms but must not derive repair eligibility, progress, or terminal outcomes.
- Writer ownership has two layers: the sprint mutation lease serializes scope, while run-control fencing validates the exact durable owner at publication time. Both remain necessary.
- Confirmation freshness currently binds canonical request, runtime identity, scope, and governed files. Repair additionally requires single-use binding to packet digest, target identity, mode, effective limits, and durable acceptance.
- All detailed repair records must be private, bounded, strict, digest-bound, and immutable where specified. Only bounded current summary facts belong in `flow-state.json`.
- Cancellation, timeout, target drift, stale governed inputs, persistence failure, rejected evidence, output truncation, unavailable isolation, or uncertain cleanup cannot produce a verified outcome.
- Manual mode cannot emit `stalled` or schedule a second mutation. Automatic mode cannot start without current qualifying manual proof and explicit automatic confirmation.
- Progress events and SSE are observational. Status, cycle, result, and outcome must reload canonical sprint and durable-operation facts.
- Server shutdown must reject new repair mutation, cancel through the same canonical path, wait only within bounds, and preserve exactly one terminal authority.
- Git mutation, test or expectation weakening, evidence edits, governed-input edits, and unconfirmed scope growth are prohibited repair actions.

## Open Questions

- No `internal/sprint/qa_repair.go` or `qa_repair_test.go` exists, so there is no current product-owned patch parser/applicator, changed-byte accounting boundary, progressive repair reverifier, progress comparator, terminal repair outcome derivation, or manual-proof publisher to reuse directly.
- The current promoted `QAIssue` does not retain source theory IDs, violated expectation references, exact reproducer descriptors, allowed/forbidden paths, containing checks, or repair acceptance criteria. Their authoritative derivation across current QA records must remain unambiguous and digest-bound.
- `QACheckDescriptor` currently describes read-only checks such as `gofmt -d`; the repository does not yet expose a complete immutable reproducer/shard/boundary/containing/delta-review command catalog for repair reverification.
- The existing generic operation terminal states and repair semantic outcomes are different vocabularies. Their authoritative mapping, especially cancellation, interruption, cleanup uncertainty, and exactly-one-result arbitration, is not yet represented in repair code.
- The existing generic confirmation token supports deduplication but is not explicitly single-use after successful mutation and does not encode manual-versus-automatic authority. Repair-specific confirmation records are absent.
- Repair-specific lower-only configuration fields and immutable maxima are absent from `internal/platform/config/qa.go`; only Sprint 37 QA limits currently exist.
