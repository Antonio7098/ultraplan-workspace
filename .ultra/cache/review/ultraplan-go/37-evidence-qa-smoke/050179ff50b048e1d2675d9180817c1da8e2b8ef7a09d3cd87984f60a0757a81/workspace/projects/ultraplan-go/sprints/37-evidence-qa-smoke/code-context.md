# Sprint Code Context

## Sprint Scope

Sprint 37 extends deterministic read-only QA into isolated evidence-producing verification while preserving target immutability, durable operation authority, and all existing smoke guarantees.

The implementation scope includes:

- One fresh, contained writable workspace per writable shard attempt, including dirty and non-Git targets.
- Frozen evidence plans with explicit conditions, paths, argv, environment, timeouts, output caps, and cleanup requirements.
- Generated tests, fixtures, probes, and smoke scenarios inside isolation only.
- Durable generated patches and evidence records retained after workspace cleanup.
- Product-owned evidence validation, adjudication, issue promotion, assessment, and canonical `qa.md`.
- No production repair, generated-patch application, Git mutation, or investigator-owned promotion.
- Detailed private records under `verification/` and only bounded projections in `flow-state.json`.
- Existing smoke execution exposed as QA suite `smoke` through the current manifest-driven implementation.
- Compatibility for `smoke`, `smoke.md`, external evidence, existing clients, and current review/smoke verdict behavior.
- Consistent CLI, JSON, TUI, browser, durable-run, cancellation, resume, and recovery behavior.
- Fail-closed handling of stale, malformed, flaky, narrow, diagnostic, uncontained, cancelled, timed-out, truncated, or cleanup-uncertain evidence.

## Inspected Repository Areas

- `internal/sprint` QA domain, mapping, investigation, synthesis, persistence, recovery, assessment, and smoke behavior.
- `internal/platform/process` direct execution, bounded streams, cancellation, timeout, process groups, and descendant cleanup.
- `internal/runcontrol` acceptance, claims, fencing, heartbeat, cancellation, reconciliation, and terminal authority.
- `internal/app` shared QA DTOs, operation preparation, durable execution, CLI parsing, configuration, JSON, and errors.
- `internal/web` QA queries, operation decoding, routing, security, confirmation, and server-rendered presentation.
- `internal/tui` bounded QA rendering and operation navigation.
- QA, smoke, process, app, browser, and TUI tests.
- Confirmed absence of the required new isolation and adjudication source/test files.

## Selected Source References

### Verification Phase Boundary
- **Path:** `internal/sprint/verification_phase.go`
- **Lines:** `5-51`
- **Symbol:** `VerificationPhase`, `ParseVerificationPhase`, `CompatibilityStage`
- **Rationale:** Establishes QA as verification work outside planning-stage order and restricts legacy planning-stage compatibility to Conformance Review.

### Flow-State Verification Ownership
- **Path:** `internal/sprint/domain.go`
- **Lines:** `135-170`
- **Symbol:** `FlowState`, `VerificationAttempt`
- **Rationale:** Shows that flow state stores only a bounded QA summary while transient verification work has a separate attempt record.

### QA Lifecycle And Budgets
- **Path:** `internal/sprint/qa_types.go`
- **Lines:** `16-142`
- **Symbol:** `QAPhaseStatus`, `QATheoryOutcome`, `QATerminalResult`, `QABudgets`
- **Rationale:** Defines the closed lifecycle and finite limits that writable investigation, evidence, adjudication, and recovery must preserve.

### Durable QA State
- **Path:** `internal/sprint/qa_types.go`
- **Lines:** `144-241`
- **Symbol:** `QAFreshness`, `QABlocker`, `QAWriterToken`, `QARunCorrelation`, `QAState`, `QAFlowSummary`
- **Rationale:** Defines current freshness, blockers, durable writer correlation, canonical state, and the detailed-state-versus-flow-summary boundary.

### Current Map, Evidence, And Theory Records
- **Path:** `internal/sprint/qa_types.go`
- **Lines:** `243-417`
- **Symbol:** `QAMap`, `QAShard`, `QAEvidenceSummary`, `QAInvestigatorAttempt`, `QATheory`, `QAChallenge`, `QASynthesis`
- **Rationale:** Supplies current map and theory identities but only summary-level evidence, with no frozen evidence plan, workspace identity, generated patch, adjudication, promoted issue, or assessment records.

### Deterministic QA Identities
- **Path:** `internal/sprint/qa_types.go`
- **Lines:** `419-507`
- **Symbol:** `QASemanticIdentity`, `QAShardIdentity`, `QATheoryIdentity`, `newQAID`
- **Rationale:** Current scoped IDs cover attempts, maps, shards, theories, challenges, and synthesis; new evidence, patch, adjudication, and issue identities must retain explicit deterministic compatibility.

### QA Validation And Error Taxonomy
- **Path:** `internal/sprint/qa_types.go`
- **Lines:** `529-602`
- **Symbol:** `ValidateQASettings`, `validateQABudgets`, `ValidateQATheory`
- **Rationale:** Enforces hard policy maxima and complete falsifiable theory conditions before accepting runtime output.

### Deterministic Map Admission
- **Path:** `internal/sprint/qa_map.go`
- **Lines:** `41-137`
- **Symbol:** `Service.QAMap`
- **Rationale:** Builds a runtime-free, write-free map from execute evidence, current Conformance Review, target identity, QA policy, and approved checks.

### Deterministic Shard Construction
- **Path:** `internal/sprint/qa_map.go`
- **Lines:** `139-276`
- **Symbol:** `BuildQAMap`
- **Rationale:** Assigns every changed path to one bounded primary owner, records boundary overlap, blocks unknown classifications, and derives stable attempt and shard identities.

### Target Identity Inputs
- **Path:** `internal/sprint/qa_map.go`
- **Lines:** `287-412`
- **Symbol:** `validateQAPath`, `qaAdjacentContext`, `qaGitIdentity`, `qaWorkspaceInputRef`
- **Rationale:** Contains path validation, optional Git identity, and governed input references; isolation must support this identity model without becoming Git-only.

### QA Status And Recovery
- **Path:** `internal/sprint/qa.go`
- **Lines:** `53-205`
- **Symbol:** `Service.QAStatus`, `Service.QAShard`, `Service.QATheory`, `Service.RecoverQA`
- **Rationale:** Reads pointer-owned records and performs explicit runtime-free recovery without adopting prior workers or sessions.

### QA Run Orchestration
- **Path:** `internal/sprint/qa.go`
- **Lines:** `234-360`
- **Symbol:** `Service.RunQA`
- **Rationale:** Owns the lifecycle from writer validation through shards, follow-up, synthesis, and terminal publication; it is the central extension point for evidence, adjudication, assessment, and `qa.md`.

### QA Resume And Scheduling
- **Path:** `internal/sprint/qa.go`
- **Lines:** `385-522`
- **Symbol:** `Service.prepareQAAttempt`, `Service.runQAShardBatch`
- **Rationale:** Resumes current terminal shards and provides bounded scheduling, cancellation, terminal persistence, and shutdown after publication failure; writable scheduling must additionally prove independent isolation.

### Current Read-Only Investigator
- **Path:** `internal/sprint/qa.go`
- **Lines:** `524-665`
- **Symbol:** `Service.runOneQAShardSafely`, `Service.runOneQAShard`
- **Rationale:** Revalidates target and map identity, enforces default-deny permissions and retry limits, executes map-owned checks, and rejects target drift.

### Investigator Output Decoder
- **Path:** `internal/sprint/qa.go`
- **Lines:** `685-711`
- **Symbol:** `decodeQAInvestigatorOutput`
- **Rationale:** Treats model output as bounded strict JSON and rejects unknown fields, unsupported schemas, empty output, and trailing data.

### Approved Check Contract
- **Path:** `internal/sprint/qa_prompt.go`
- **Lines:** `16-125`
- **Symbol:** `QACheckDescriptor`, `ApprovedQAChecks`, `validateQACheckDescriptor`
- **Rationale:** Freezes executable, argv, cwd, environment, timeout, output cap, and fingerprint while prohibiting shells, Git, write modes, redirection, and path escape.

### Investigator Runtime Permissions
- **Path:** `internal/sprint/qa_prompt.go`
- **Lines:** `136-191`
- **Symbol:** `Service.RenderQAInvestigatorPrompt`, `Service.QAInvestigatorRequest`
- **Rationale:** Defines current default-deny, path-scoped read authority and explicitly prohibits generated tests, shell execution, issue promotion, and repair.

### Approved Check Execution
- **Path:** `internal/sprint/qa_prompt.go`
- **Lines:** `243-285`
- **Symbol:** `Service.RunApprovedQACheck`
- **Rationale:** Executes map-owned descriptors through the platform process seam and rejects target drift, command failure, truncation, and exhausted output bounds.

### Current Synthesis Authority
- **Path:** `internal/sprint/qa_synthesis.go`
- **Lines:** `27-166`
- **Symbol:** `SynthesizeQA`, `SynthesizeQAWithChallenges`
- **Rationale:** Retains outcomes, validates challenges, detects contradictions and interactions, and creates bounded follow-ups without issue-promotion or assessment authority.

### QA State Paths And Strict Loading
- **Path:** `internal/sprint/qa_state.go`
- **Lines:** `138-340`
- **Symbol:** QA relative-path helpers, `QAStore.resolve`, `QAStore.LoadState`, `QAStore.readStrict`
- **Rationale:** Defines the private verification layout and rejects invalid IDs, path escape, symlinks, public modes, oversized records, unknown schemas, malformed JSON, and mismatched identities.

### Atomic QA Publication
- **Path:** `internal/sprint/qa_state.go`
- **Lines:** `343-464`
- **Symbol:** `QAPublication`, `QAStore.Publish`, `QAStore.SaveRecoveredState`, `QAStore.checkWriter`
- **Rationale:** Publishes detailed records before the state pointer and flow summary, requires current writer fencing, and prevents recovery from publishing active work.

### Private Record Integrity
- **Path:** `internal/sprint/qa_state.go`
- **Lines:** `466-579`
- **Symbol:** `QAStore.writeRecord`, `privateAtomicWrite`, `QAStore.verifyReference`, `qaFlowSummary`
- **Rationale:** Enforces private atomic writes, immutable maps, hard bounds, digest verification, and bounded flow projection.

### Flow-State QA Validation
- **Path:** `internal/sprint/state.go`
- **Lines:** `294-383`
- **Symbol:** `ValidateFlowState`, `validateQAFlowSummary`
- **Rationale:** Keeps QA outside planning order and validates a contained state pointer, digest, counts, phase, and next action.

### Sprint Service Seams
- **Path:** `internal/sprint/service.go`
- **Lines:** `22-76`
- **Symbol:** `Service`, `WithProcessRunner`, `WithSmokeSettings`, `NewService`
- **Rationale:** Shows existing runtime, process, QA fence, smoke, mutation-lock, and publication dependencies.

### QA Policy And Fences
- **Path:** `internal/sprint/service.go`
- **Lines:** `149-189`
- **Symbol:** `WithQASettings`, `WithQAWriterFence`, `WithQAMapFence`, `effectiveQASettings`
- **Rationale:** Freezes effective policy and provides explicit durable-writer and governed-input checks.

### Direct Process Boundary
- **Path:** `internal/platform/process/process.go`
- **Lines:** `21-151`
- **Symbol:** `Request`, `Result`, `Runner`, `DirectRunner.Run`
- **Rationale:** Provides product-neutral explicit argv, cwd, exact environment, timeout, bounded streams, cancellation, and cleanup facts without shell interpolation.

### Bounded Streams And Unix Cleanup
- **Path:** `internal/platform/process/process_unix.go`
- **Lines:** `12-34`
- **Symbol:** `configureOwnedProcess`, `stopAndWait`
- **Rationale:** Uses owned process groups and TERM-to-KILL descendant cleanup, including the group-leader exit race.

### Current Freshness Policy
- **Path:** `internal/sprint/freshness_policy.go`
- **Lines:** `3-15`
- **Symbol:** `strictCompletedReviewSnapshotFreshness`, `strictCompletedSmokeSnapshotFreshness`, `strictSmokeAuthorProtectedSnapshots`
- **Rationale:** Target snapshot freshness and protected smoke-author snapshots are disabled, so current completed evidence does not satisfy all Sprint 37 drift requirements.

### Smoke Result And Compatibility State
- **Path:** `internal/sprint/smoke_types.go`
- **Lines:** `21-67`
- **Symbol:** `SmokeVerdict`, `SmokePhase`, `SmokeSettings`, `SmokeRequest`
- **Rationale:** Defines existing smoke controls and process bounds that QA suite `smoke` must preserve.

### Smoke Canonical Records
- **Path:** `internal/sprint/smoke_types.go`
- **Lines:** `164-250`
- **Symbol:** `SmokeResult`, `SmokeStageState`, `SmokeCompletion`
- **Rationale:** Carries protocol identity, coverage, evidence, issues, diagnostic status, reconciliation, fingerprints, and last-complete state.

### Smoke Static Preparation
- **Path:** `internal/sprint/smoke_protocol.go`
- **Lines:** `111-210`
- **Symbol:** `Service.prepareSmokeStatic`
- **Rationale:** Resolves the cataloged harness, contained manifest, executable, cwd, evidence roots, target, current review gate, and diagnostic override.

### Smoke Protocol Validation
- **Path:** `internal/sprint/smoke_protocol.go`
- **Lines:** `213-372`
- **Symbol:** `validateSmokeManifest`, `validateSmokeDiscovery`
- **Rationale:** Enforces protocol capabilities, authoring-path safety, environment names, identity relationships, prerequisites, and complete coverage mappings.

### Smoke Scope And Environment
- **Path:** `internal/sprint/smoke_protocol.go`
- **Lines:** `441-540`
- **Symbol:** `selectSmoke`
- **Rationale:** Distinguishes diagnostic scopes from authoritative containing suites and blocks missing prerequisites or incomplete mappings.

### Smoke Execution Pipeline
- **Path:** `internal/sprint/smoke.go`
- **Lines:** `63-188`
- **Symbol:** `Service.runSmoke`
- **Rationale:** Owns authoring, discovery, selection, process invocation, cleanup checks, evidence validation, diagnostic classification, and canonical commit.

### Smoke Evidence Validation
- **Path:** `internal/sprint/smoke.go`
- **Lines:** `304-411`
- **Symbol:** `validateSmokeRun`
- **Rationale:** Validates protocol, scope, counts, test identities, contained evidence paths, hashes, and detailed issue coverage for failed tests.

### Smoke Publication And Side Effects
- **Path:** `internal/sprint/smoke.go`
- **Lines:** `21-60`
- **Symbol:** `Service.RunSmoke`
- **Rationale:** Records attempts, may mark roadmap delivery, and publishes the smoke stage; QA invocation must preserve compatibility without duplicating these effects.

### Smoke Canonical Commit
- **Path:** `internal/sprint/smoke.go`
- **Lines:** `439-481`
- **Symbol:** `Service.commitSmoke`
- **Rationale:** Atomically publishes validated `smoke.md`, computes fingerprints, retains external links, and exposes flow-state reconciliation failure.

### Smoke Authoring Boundary
- **Path:** `internal/sprint/smoke_author.go`
- **Lines:** `20-120`
- **Symbol:** `Service.authorSmokeSuite`
- **Rationale:** Restricts writes to manifest-declared harness paths and rejects unsupported permissions or observed protected-target writes.

### Product-Owned Verification Assessment
- **Path:** `internal/sprint/verify.go`
- **Lines:** `142-285`
- **Symbol:** `Service.VerificationStatus`, `deriveAssessment`
- **Rationale:** Current assessment derives from Conformance Review and smoke and prevents diagnostic smoke from overriding review failure; canonical QA must extend this product-owned decision boundary.

### Adapter-Independent QA Contract
- **Path:** `internal/app/sprint_usecases.go`
- **Lines:** `68-135`
- **Symbol:** `QAUseCases`, `QAQueries`, `QARequest`, `QAResult`
- **Rationale:** Defines the shared adapter boundary that must expose writable attempts, evidence, adjudication, issues, assessment, smoke suite, cancellation, and recovery.

### QA Mutation Use Cases
- **Path:** `internal/app/sprint_usecases.go`
- **Lines:** `695-757`
- **Symbol:** `dashboardUseCases.RunQA`, `ResumeQA`, `CancelQA`, `RecoverQA`
- **Rationale:** Delegates execution to the shared runner, verifies durable ownership for cancellation, performs explicit recovery, and projects Conformance Review independently.

### Shared Operation Preparation
- **Path:** `internal/app/operations.go`
- **Lines:** `211-373`
- **Symbol:** `dashboardUseCases.PrepareOperation`, `validateQAOperationRequest`
- **Rationale:** Prepares QA operations using governed fingerprints and map-owned shards while rejecting caller-owned model, stage, smoke, timeout, review, and parallelism controls.

### Shared Runtime QA Execution
- **Path:** `internal/app/operation_runner.go`
- **Lines:** `114-130`
- **Symbol:** `sharedOperationRunner` QA branch
- **Rationale:** Obtains durable ownership, installs the writer fence, and delegates workflow semantics to `Service.RunQA` for both terminal and browser adapters.

### Durable QA Ownership
- **Path:** `internal/app/durable_operations.go`
- **Lines:** `97-155`
- **Symbol:** `durableOperationManager.AcceptOperation`, `qaOwnershipFromContext`
- **Rationale:** Persists acceptance and claim before execution and converts the run-control fence into a heartbeat-verified QA writer token.

### Run-Control Authority
- **Path:** `internal/runcontrol/interfaces.go`
- **Lines:** `47-64`
- **Symbol:** `Repository`
- **Rationale:** Is the durable authority for acceptance, claims, events, heartbeat, cancellation, terminal arbitration, reconciliation, and retention.

### Run-Control Attempts And Fencing
- **Path:** `internal/runcontrol/model.go`
- **Lines:** `267-298`
- **Symbol:** `Attempt`, `Fence`, `Fence.Validate`
- **Rationale:** Binds writer authority to run, operational attempt, owner, lease, and fencing generation.

### CLI QA Dispatch
- **Path:** `internal/app/sprint_commands.go`
- **Lines:** `369-466`
- **Symbol:** `runSprintCommand` QA branch
- **Rationale:** Implements current map, status, recovery, cancellation, run, and resume with durable acceptance and writer fencing.

### CLI QA Controls
- **Path:** `internal/app/sprint_commands.go`
- **Lines:** `581-679`
- **Symbol:** `parseSprintQAArgs`, `renderSprintQA`, `mapQACommandError`
- **Rationale:** Freezes current public controls, verdict-neutral text, and error behavior; no suite selector or adjudication output exists.

### Effective QA Configuration
- **Path:** `internal/app/sprint_commands.go`
- **Lines:** `713-779`
- **Symbol:** `qaSettings`
- **Rationale:** Resolves model fallback, configuration provenance, timeouts, concurrency, and all bounded QA settings used by production services.

### Browser QA Queries
- **Path:** `internal/web/qa_handlers.go`
- **Lines:** `9-107`
- **Symbol:** QA JSON and HTML handlers
- **Rationale:** Keeps browser queries over `internal/app` and provides JSON and no-JavaScript HTML without importing sprint internals.

### Browser Operation Decoding
- **Path:** `internal/web/operation_handlers.go`
- **Lines:** `593-741`
- **Symbol:** `decodeStrictJSON`, `mapOperationRequest`
- **Rationale:** Enforces strict JSON, bounded identifiers, closed operation kinds, sprint scope, and QA-specific option rejection.

### Browser Security And Confirmation
- **Path:** `internal/web/security.go`
- **Lines:** `368-439`
- **Symbol:** `preparationStore.issue`, `preparationStore.consume`
- **Rationale:** Binds confirmation to session, canonical request, and governed fingerprint and consumes it once before writable operations.

### Server-Rendered QA Presentation
- **Path:** `internal/web/templates/run_qa.html`
- **Lines:** `13-107`
- **Symbol:** `component/run-qa`
- **Rationale:** Presents current shards, attempts, checks, evidence summaries, theories, synthesis, fingerprints, blockers, cancellation, and limits but lacks adjudication, issues, assessment, and smoke-suite status.

### TUI QA Presentation
- **Path:** `internal/tui/qa_view.go`
- **Lines:** `8-45`
- **Symbol:** `renderSprintQAView`
- **Rationale:** Consumes only bounded app DTOs and renders completion as investigation completion rather than a pass verdict.

### QA Persistence Tests
- **Path:** `internal/sprint/qa_state_test.go`
- **Lines:** `119-269`
- **Symbol:** QA publication, strict loading, symlink, writer, and atomicity tests
- **Rationale:** Covers private modes, pointer order, strict schemas, symlink rejection, stale writers, immutable maps, atomic failure, and prior-state preservation.

### QA Runtime Safety Tests
- **Path:** `internal/sprint/qa_test.go`
- **Lines:** `89-223`
- **Symbol:** QA worker, persistence, panic, permission, and drift tests
- **Rationale:** Covers bounded workers, publication failure, panic containment, unsupported permissions, target mutation, and governed-input drift.

### QA Cancellation And Output Tests
- **Path:** `internal/sprint/qa_test.go`
- **Lines:** `263-322`
- **Symbol:** `TestQACancellationPersistsActiveShardWithoutRetry`, `TestQAInvestigationOutputIsStrictAndBounded`, `TestQARecoveryMissingStateIsRuntimeFreeNoOp`
- **Rationale:** Covers cancellation persistence, no cancellation retry, strict bounded model output, and runtime-free recovery.

### QA Policy And Synthesis Tests
- **Path:** `internal/sprint/qa_synthesis_test.go`
- **Lines:** `10-115`
- **Symbol:** QA synthesis outcome, determinism, challenge, and follow-up tests
- **Rationale:** Requires retained outcomes, contradictions, deterministic follow-up, validated challenge records, and no issue, repair, verdict, or `qa.md` authority.

### Smoke End-To-End Tests
- **Path:** `internal/sprint/smoke_test.go`
- **Lines:** `301-374`
- **Symbol:** `TestSmokeRunCommitsValidatedArtifactAndPreservesItOnMalformedRun`
- **Rationale:** Provides the fixture baseline for discovery, containing-suite selection, execution, evidence, canonical publication, and overall smoke assessment.

### Smoke Failure Preservation Tests
- **Path:** `internal/sprint/smoke_test.go`
- **Lines:** `440-485`
- **Symbol:** `TestSmokeRunCommitsValidatedArtifactAndPreservesItOnMalformedRun`
- **Rationale:** Verifies external evidence behavior, permission failure, malformed discovery rejection, and preservation of the last complete `smoke.md` and state.

### Process Safety Tests
- **Path:** `internal/platform/process/process_test.go`
- **Lines:** `15-98`
- **Symbol:** Direct runner environment, cleanup, progress, timeout, and bounds tests
- **Rationale:** Provides reusable fixtures for exact cwd/environment, cancellation, descendant cleanup, non-blocking progress, truncation, timeout, and cleanup completion.

### Application QA Tests
- **Path:** `internal/app/sprint_usecases_test.go`
- **Lines:** `29-116`
- **Symbol:** Canonical QA fixture, invalid-state, bounded-projection, and shard-observability tests
- **Rationale:** Freezes adapter-independent facts, safe invalid-state handling, bounded sanitization, and shared projection behavior.

### Browser QA Tests
- **Path:** `internal/web/qa_handlers_test.go`
- **Lines:** `38-109`
- **Symbol:** QA public error and route tests
- **Rationale:** Requires stable errors, bounded JSON, hostile-content escaping, complete no-JavaScript HTML, and read-only query methods.

### Browser Operation Contract Tests
- **Path:** `internal/web/operations_contract_test.go`
- **Lines:** `19-87`
- **Symbol:** Browser operation-kind and QA-option contract tests
- **Rationale:** Freezes operation compatibility and limits caller-controlled QA scope to a current map-owned shard.

### TUI QA Tests
- **Path:** `internal/tui/qa_view_test.go`
- **Lines:** `13-69`
- **Symbol:** Verdict-neutral and focused QA view tests
- **Rationale:** Requires narrow-width rendering, focused state, blockers, cancellation, recovery guidance, and bounded operations without inventing a pass verdict.

## Relationships

- `Service.QAMap` derives a semantic attempt from execute evidence, current Conformance Review, target identity, QA policy, and approved checks.
- Durable acceptance and claim occur before `Service.RunQA` receives writer authority.
- `Service.RunQA` holds the sprint mutation lease, executes map-owned shards, and publishes through a fenced `QAStore`.
- `QAStore.Publish` writes detailed private records before `verification/state.json` and then updates the bounded flow projection.
- Current investigators are runtime-backed and read-only; approved checks and smoke processes share `internal/platform/process.Runner`.
- The process package owns generic isolation and process mechanics, while `internal/sprint` owns evidence meaning, adjudication, issue promotion, and assessment.
- Current synthesis is deterministic pre-adjudication input and cannot promote issues, classify regression candidates, write `qa.md`, or alter Conformance Review.
- Run control remains authoritative for operational ownership, events, cancellation, reconciliation, and terminal arbitration.
- Sprint QA remains authoritative for maps, shards, theories, evidence, adjudication, issues, assessment, and canonical QA artifacts.
- Existing smoke authority flows through static preparation, protocol validation, scope selection, execution, evidence validation, and `smoke.md` publication.
- QA suite `smoke` must reuse that path so compatibility and QA execution cannot diverge.
- `internal/app` remains the adapter boundary; CLI, TUI, and browser do not read verification state directly.
- Browser mutation continues through shared preparation, confirmation, durable execution, cancellation, and recovery.
- Conformance Review remains independent input that QA may consume but cannot rewrite or upgrade.

## Constraints

- QA and repair do not enter `PlanningStage`.
- Every writable attempt gets a fresh isolated workspace.
- Writable work starts only after identity, containment, path safety, process containment, and cleanup capability are proven.
- Isolation supports dirty, uncommitted, and non-Git targets; Git worktrees are optional only.
- The shared target, governed inputs, historical theories, out-of-scope harness paths, and Git state remain immutable.
- Investigators create verification code only in isolation and cannot promote findings.
- Every generated check has a frozen plan with explicit conditions and bounded execution details.
- Commands use explicit argv, contained cwd, bounded environment, finite timeout, bounded output, cancellation, and descendant cleanup.
- Generated patches and evidence survive cleanup but are never applied to the target.
- Writable concurrency remains sequential or disabled until independent isolation and cleanup are proven.
- Only product-owned adjudication promotes issues, classifies regression candidates, or marks repair eligibility.
- Command failures, runtime exits, investigator claims, model responses, and harness issues are not automatically promoted issues.
- Product code derives assessment from current review, evidence, adjudication, blockers, and required containing-suite results.
- Stale, malformed, missing, mismatched, uncontained, narrow, diagnostic, flaky, truncated, cancelled, timed-out, or cleanup-uncertain evidence cannot pass.
- Detailed records remain private, versioned, bounded, digest-linked, contained, and atomically published.
- Unknown schemas, invalid digests, stale writers, escapes, over-limit records, and partial publication fail closed.
- Failed attempts preserve the last complete report while exposing current failure.
- `flow-state.json` remains a bounded projection and does not absorb evidence or run-control history.
- Raw smoke output and per-test artifacts remain in manifest-declared external roots.
- `smoke`, `smoke.md`, existing flow-state fields, and existing client contracts remain compatible.
- Narrow or diagnostic smoke cannot replace containing-suite evidence.
- Browser and TUI consume bounded app DTOs and cannot infer authority from progress events.
- Browser disconnect stops observation only.
- Raw provider payloads, unrestricted output, environment values, and secrets do not enter canonical state or presentation.
- Production repair, patch application, permanent regression-test promotion, Git mutation, and repair cycles remain out of scope.

## Open Questions

- Where is the Sprint 36 admission gate enforced, since current map admission checks Conformance Review but not current required smoke evidence?
- What product-neutral isolation request, workspace identity, cleanup result, and fault-injection seams should the absent process isolation service expose?
- How will copy isolation represent dirty and untracked files while rejecting special files, symlinks, escapes, races, and over-limit trees?
- Which identity facts distinguish investigator mutation from unrelated target changes without relying on disabled snapshot switches?
- Does expanded state remain schema version 1 through additive fields, or require schema and ID migration?
- How are generated patches bounded, normalized, fingerprinted, and atomically retained outside temporary workspaces?
- How will adjudication distinguish harness issue metadata, investigator claims, check failures, and product-promoted issues?
- Which adjudication decisions are deterministic rules, and which model outputs may only propose groupings or follow-up?
- How is repeatability represented for deterministic proof versus repeated empirical execution?
- Which component owns canonical `qa.md` and last-complete preservation while `smoke.md` remains compatible?
- How should roadmap delivery and smoke-stage publication behave when `RunSmoke` is initiated through QA?
- What additive CLI and browser request shape introduces suite `smoke` without exposing executable, argv, environment, runtime, or evidence-path authority?
- Should synthesis remain a distinct pre-adjudication record?
- Which issue and assessment fields can be added to stable app and JSON contracts without breaking existing clients?
- How will browser and TUI fixtures represent patches and evidence without exposing raw output or unbounded hostile content?
