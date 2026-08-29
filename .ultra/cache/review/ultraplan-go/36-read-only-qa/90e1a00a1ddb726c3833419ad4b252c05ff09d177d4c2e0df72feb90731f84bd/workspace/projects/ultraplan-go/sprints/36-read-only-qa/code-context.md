# Sprint Code Context

## Sprint Scope

Sprint 36 introduces deterministic, resumable, read-only QA decomposition and synthesis after execution and Conformance Review. QA must be represented by a verification-phase model independent of `PlanningStage`, map governed inputs and changed behavior into bounded shards, run investigators under enforced non-mutating runtime permissions, retain theory outcomes, and synthesize cross-shard results without creating tests, modifying production code, promoting repair issues, or performing automatic repair.

Detailed QA map, shard, theory, attempt, and synthesis state belongs in a schema-versioned atomic artifact outside `flow-state.json`. Flow state may retain only bounded canonical QA summary, freshness, verdict, and pointer data. Runtime-backed QA must reuse the durable run-control authority and remain observable and cancellable through CLI, TUI, and web adapters.

The existing `review` command and `review.md` artifact remain compatible while human-facing terminology changes to Conformance Review. Smoke integration and evidence-producing QA remain Sprint 37 work.

## Inspected Repository Areas

- Sprint domain types, planning-stage vocabulary, verification summaries, artifact paths, containment, state validation, and atomic persistence.
- Execute evidence, target identity, changed-path extraction, Conformance Review manifests, deterministic fingerprints, resumable attempts, and freshness checks.
- Generic runtime sandbox and permission-policy propagation into Agentwrap.
- App operation kinds, preparation fingerprints, shared runtime dispatch, durable acceptance, cancellation, and run observation.
- CLI command dispatch and durable runtime-backed command inventory.
- TUI routes, operation projections, durable run views, and existing verification presentation.
- Web operation DTOs, allowlisted operation mapping, confirmation flow, compatibility tests, and import boundaries.
- Configuration surfaces, architecture documentation, QA roadmap, and release gates.
- No Sprint 36 QA implementation files or generated sprint artifacts were present in this checkout.

## Selected Source References

### Sprint 36 Product Boundary
- **Path:** `docs/plans/integrated-roadmap.md`
- **Lines:** `488-554`
- **Symbol:** `Sprint 36 - Read-only QA decomposition and synthesis`
- **Rationale:** Defines the repository's terminology, verification-phase separation, deterministic map inputs, investigator permissions, synthesis responsibilities, browser visibility, and exit criteria. It also limits this sprint to read-only QA without generated tests or repair.

### Existing Planning Stage Vocabulary
- **Path:** `internal/sprint/domain.go`
- **Lines:** `30-40`
- **Symbol:** `PlanningStage`
- **Rationale:** Enumerates the canonical planning stages. The new verification phase must remain a separate type rather than extending this planning sequence.

### Existing Verification State Shape
- **Path:** `internal/sprint/domain.go`
- **Lines:** `117-169`
- **Symbol:** `ExecuteRunState`, `FlowState`, `VerificationAttempt`
- **Rationale:** Shows authoritative execute target identity, the current embedding of detailed review and smoke state in `flow-state.json`, and the shared bounded attempt vocabulary. QA detailed state must not repeat the existing embedding pattern.

### Current Verification Projection
- **Path:** `internal/sprint/domain.go`
- **Lines:** `186-244`
- **Symbol:** `VerificationStage`, `VerificationStatus`, `StatusSummary`
- **Rationale:** Defines the canonical freshness, artifact, fingerprint, attempt, verdict, and next-action projection currently limited to review and smoke. This is the compatibility seam for a bounded QA summary and pointer.

### Sprint Artifact Containment
- **Path:** `internal/sprint/artifacts.go`
- **Lines:** `11-74`
- **Symbol:** `ArtifactRelPath`, `FlowStateRelPath`, `resolveSprintContained`
- **Rationale:** Centralizes stage artifact names and enforces lexical and canonical sprint-root containment. New QA summary and detailed-state paths must use equivalent containment rather than direct path joins.

### Strict Flow-State Loading
- **Path:** `internal/sprint/state.go`
- **Lines:** `20-80`
- **Symbol:** `LoadFlowState`
- **Rationale:** Enforces explicit schema versions, rejects unknown fields and trailing JSON, supports a bounded legacy migration, and validates loaded state before use. Any QA summary added to flow state must preserve this fail-closed compatibility behavior.

### Atomic Flow-State Persistence
- **Path:** `internal/sprint/state.go`
- **Lines:** `201-288`
- **Symbol:** `SaveFlowState`, `saveFlowStateWithHooks`
- **Rationale:** Preserves prior verification records during planning refreshes and uses validated temporary-file, flush, rename, and directory-sync persistence. Separate detailed QA state should provide the same atomic preservation guarantees without becoming flow-state authority.

### Execute Evidence Authority
- **Path:** `internal/sprint/execute_state.go`
- **Lines:** `14-117`
- **Symbol:** `LoadExecuteRunState`, `SaveExecuteRunState`
- **Rationale:** Provides the contained, validated, schema-versioned execution record containing target identity, plan fingerprint, task outcomes, and evidence. QA mapping should consume this authority rather than rediscovering execution facts from prose.

### Review Manifest Inputs
- **Path:** `internal/sprint/review.go`
- **Lines:** `167-290`
- **Symbol:** `Service.PrepareReview`
- **Rationale:** Demonstrates how requirements, code context, sprint reasoning, plan, execute evidence, project contracts, changed paths, and target revision identity are frozen into a deterministic governed-input manifest. These inputs substantially overlap the required QA map inputs.

### Changed-Path and Manifest Fingerprints
- **Path:** `internal/sprint/review.go`
- **Lines:** `1329-1412`
- **Symbol:** `fingerprintReviewManifest`, `reviewChangedPaths`, `excludeGovernedReviewPaths`
- **Rationale:** Supplies deterministic hashing and sorted changed-path extraction while excluding governed workspace artifacts from target changes. QA mapping should retain these identity and ordering properties while adding risk and verification-surface decomposition.

### Enforced Read-Only Reviewer Runtime
- **Path:** `internal/sprint/review.go`
- **Lines:** `870-979`
- **Symbol:** `Service.runReviewer`
- **Rationale:** Configures `read_only` sandboxing, restricted permissions, default deny, an allowlist of read/list/search tools, required permission capability, cancellation, and fail-closed handling when policy enforcement is unsupported. This is the nearest existing runtime safety boundary for QA investigators.

### Read-Only Prompt Contract
- **Path:** `internal/sprint/review.go`
- **Lines:** `1509-1541`
- **Symbol:** `renderReviewerPrompt`, `reviewerInputPacket`
- **Rationale:** Builds bounded investigator prompts from frozen inputs, prohibits writes and destructive actions, requires structured output and citations, and prevents sibling coverage sources from leaking between independent reviewers.

### Resumable Verification State
- **Path:** `internal/sprint/review.go`
- **Lines:** `1127-1250`
- **Symbol:** `validateReviewStageState`, `Service.saveReviewState`
- **Rationale:** Validates attempt identity, input fingerprints, models, bounded checkpoints, retained sessions, last-complete evidence, and terminal transitions. It is the established resumability behavior to preserve while moving QA's more detailed shard state into its own artifact.

### Sole Review-to-Smoke Transition
- **Path:** `internal/sprint/verify.go`
- **Lines:** `19-95`
- **Symbol:** `VerifyRequest`, `VerifyResult`, `Service.Verify`
- **Rationale:** Owns the current execute-evidence to review to smoke sequence and explicitly identifies itself as the sole transition. Sprint 36 must add QA without silently redefining smoke as QA or breaking existing review/smoke compatibility.

### Verification Freshness Derivation
- **Path:** `internal/sprint/verify.go`
- **Lines:** `142-249`
- **Symbol:** `Service.VerificationStatus`
- **Rationale:** Derives currentness from input fingerprints, artifact digests, validated content, retained attempts, and external evidence without mutating durable state. QA status should follow the same separation between read-only projection and explicit state transition.

### Sprint Service Composition
- **Path:** `internal/sprint/service.go`
- **Lines:** `21-75`
- **Symbol:** `Service`, `StageRuntime`, `WithoutStatusWrites`
- **Rationale:** Shows current dependencies, planning-stage-keyed runtime configuration, mutation coordination, process runner, and non-persisting status mode. A verification-phase runtime cannot be represented safely by merely adding QA to `map[PlanningStage]StageRuntime`.

### Generic Runtime Permission Model
- **Path:** `internal/platform/runtime/runtime.go`
- **Lines:** `21-75`
- **Symbol:** `Request`, `PermissionPolicy`, `PermissionPathRule`
- **Rationale:** Exposes sandbox mode, permission mode, tool actions, path rules, unsupported-policy behavior, required capabilities, and runtime metadata. QA must express its read-only policy through this boundary and reject runtimes that cannot enforce it.

### Runtime Permission Adapter Test
- **Path:** `internal/platform/runtime/runtime_test.go`
- **Lines:** `444-455`
- **Symbol:** `TestPermissionPathRulesMapToAdapterPolicy`
- **Rationale:** Verifies path-level deny rules survive translation into Agentwrap policy. QA safety tests should extend this coverage to the complete investigator policy and unsupported-enforcement failure path.

### Existing Configuration Surface
- **Path:** `internal/platform/config/config.go`
- **Lines:** `13-89`
- **Symbol:** `Config`, `Planning`, `Smoke`, `Agentwrap`
- **Rationale:** Contains planning/review/smoke model settings and global Agentwrap permissions but no QA-specific model, concurrency, shard, or budget configuration. Any added fields must fit the existing typed configuration and validation model.

### App Operation Contract
- **Path:** `internal/app/operations.go`
- **Lines:** `22-140`
- **Symbol:** `OperationalUseCases`, `WebOperations`, `DurableOperationManager`, `OperationKind`, `OperationRequest`
- **Rationale:** Defines the surface-neutral operation vocabulary shared by TUI and web, durable acceptance capability, and normalized request DTO. QA status, dry-run, and start operations belong at this boundary rather than in adapter-specific registries.

### Prepared Operation Freshness
- **Path:** `internal/app/operations.go`
- **Lines:** `167-294`
- **Symbol:** `dashboardUseCases.PrepareOperation`, `dashboardUseCases.RunOperation`
- **Rationale:** Normalizes allowlisted requests, declares runtime and mutation scope, fingerprints canonical governed inputs, and rejects execution when prepared inputs become stale. QA operations should retain this prepare-confirm-revalidate contract.

### Shared Runtime Dispatch
- **Path:** `internal/app/operation_runner.go`
- **Lines:** `15-110`
- **Symbol:** `sharedOperationRunner`
- **Rationale:** Is the single runtime-backed implementation used by terminal and browser adapters and maps product progress into neutral events. QA execution must be added here to preserve adapter parity and avoid divergent workflow semantics.

### Durable Acceptance and Cancellation Context
- **Path:** `internal/app/durable_operations.go`
- **Lines:** `33-121`
- **Symbol:** `beginDurableCLICommand`, `durableOperationManager.AcceptOperation`
- **Rationale:** Persists acceptance and owner claim before runtime work, establishes fencing, records the running transition, and supplies the cancellation-aware context. Runtime-backed QA must reuse this authority rather than creating a separate operation owner.

### Runtime-Backed CLI Inventory
- **Path:** `internal/app/run_control_inventory_test.go`
- **Lines:** `11-54`
- **Symbol:** `TestEveryRuntimeBackedCLIEntryUsesDurableAcceptanceInventory`
- **Rationale:** Freezes every runtime-backed CLI entry behind durable acceptance. Adding `qa` requires extending this inventory so the command cannot bypass run control.

### Existing CLI Verification Commands
- **Path:** `internal/app/sprint_commands.go`
- **Lines:** `260-449`
- **Symbol:** `runSprint`
- **Rationale:** Shows durable acceptance, cancellation context, JSON envelopes, human rendering, confirmation requirements, and error mapping for verify, execute, review, and smoke. The QA command must follow these established CLI semantics while preserving `review` compatibility.

### Durable Run Vocabulary
- **Path:** `internal/runcontrol/model.go`
- **Lines:** `9-145`
- **Symbol:** `Lifecycle`, `CancellationState`, `TerminalOutcome`, `EventType`
- **Rationale:** Defines bounded event sizes, progress coalescing, canonical lifecycle and cancellation states, and finding/artifact event types. QA shard progress and synthesis events should project through this vocabulary without inventing a second durable lifecycle.

### Surface-Neutral Run Observation
- **Path:** `internal/app/run_usecases.go`
- **Lines:** `9-57`
- **Symbol:** `RunUseCases`, `repositoryRunUseCases`
- **Rationale:** Provides sanitized durable run listing, detail, events, cancellation, and health to all adapters. QA observability and cancellation should be exposed through these existing use cases.

### Web Operation Decoder
- **Path:** `internal/web/operation_handlers.go`
- **Lines:** `600-689`
- **Symbol:** `mapOperationRequest`
- **Rationale:** Maps a closed set of browser operation strings into app-owned operation kinds. QA browser support requires an explicit mapping here and corresponding compatibility-fixture updates.

### Web Confirmation Boundary
- **Path:** `internal/web/operation_handlers.go`
- **Lines:** `17-150`
- **Symbol:** `operationOptionsRequest`, `handleOperationPrepare`, `handleOperationStart`
- **Rationale:** Strictly decodes transport DTOs, prepares operations through app use cases, issues short-lived confirmation tokens, and re-prepares before start. QA must use this flow rather than accepting arbitrary prompts, paths, or runtime policy from the browser.

### Web Import Boundary
- **Path:** `internal/web/import_boundary_test.go`
- **Lines:** `12-35`
- **Symbol:** `TestWebImportBoundary`
- **Rationale:** Enforces that production web code imports only `internal/app` and the standard library. QA handlers and projections cannot import sprint, runtime, workspace, or persistence packages directly.

### Browser Operation Compatibility
- **Path:** `internal/web/operations_contract_test.go`
- **Lines:** `19-59`
- **Symbol:** `TestBrowserOperationKindContract`
- **Rationale:** Freezes the producer/consumer operation mapping shared by Go handlers and the embedded browser client. QA operation kinds must update this table together with both producers and consumers.

### TUI State and Operation Model
- **Path:** `internal/tui/model.go`
- **Lines:** `27-130`
- **Symbol:** `RouteKind`, `Model`, `navItem`
- **Rationale:** Holds only app DTOs, bounded operation events, durable run snapshots, previews, and navigation state. QA presentation should consume app projections and avoid TUI-owned workflow or persistence state.

### TUI Verification Compatibility Test
- **Path:** `internal/tui/verify_test.go`
- **Lines:** `11-35`
- **Symbol:** `TestSprintVerificationActionsAndNarrowSummary`
- **Rationale:** Captures current review, smoke, verify, diagnostic override, issue, and next-action presentation. QA additions must update this fixture while retaining the existing compatibility actions.

### Determinism and Runtime-Safety Test Pattern
- **Path:** `internal/sprint/review_test.go`
- **Lines:** `127-179`
- **Symbol:** `TestReviewManifestExecutionAndArtifactPreservation`
- **Rationale:** Tests deterministic fingerprints, filtered changed paths, malformed-output preservation, last-complete authority, read-only sandboxing, default-deny permissions, and bounded prompts. These are direct acceptance-test patterns for QA mapping and investigators.

### Atomic Artifact Failure Test Pattern
- **Path:** `internal/sprint/review_test.go`
- **Lines:** `584-609`
- **Symbol:** `TestAtomicReviewWritePreservesPriorArtifactOnRenameFailure`, `TestReviewFanOutUsesConfiguredBound`
- **Rationale:** Verifies prior canonical output survives rename failure and fan-out obeys configured concurrency. Detailed QA state and synthesis artifacts require equivalent failure and bound coverage.

### Flow-State Compatibility Test Pattern
- **Path:** `internal/sprint/sprint_test.go`
- **Lines:** `98-148`
- **Symbol:** `TestFlowStateStrictLoadingAndAtomicWritePreservesPrior`
- **Rationale:** Tests schema rejection, path containment, and preservation of prior state after atomic promotion failure. This is the minimum persistence compatibility behavior for any QA summary added to flow state.

### Adapter and State Ownership Rules
- **Path:** `docs/architecture.md`
- **Lines:** `29-81`
- **Symbol:** `Web Adapter Boundary`, `State And Artifact Ownership`
- **Rationale:** Documents app-owned workflow semantics, web's lack of direct product/runtime/filesystem authority, bounded ephemeral transport state, and durable workspace/run-control authority. QA must preserve these ownership boundaries.

### Release Gates
- **Path:** `docs/release-checklist.md`
- **Lines:** `17-40`
- **Symbol:** `Release checks`
- **Rationale:** Establishes mandatory test, race, build, vet, diff, fake-runtime, browser compatibility, security, and import-boundary checks. QA fixtures must remain offline-capable and report unavailable real-runtime prerequisites as blocked rather than pass.

## Relationships

- `ExecuteRunState` supplies the approved implementation target, plan fingerprint, task outcomes, evidence, and changed paths consumed by deterministic verification mapping.
- `PrepareReview` already composes most governed QA inputs and target identity; QA can share these input authorities without sharing `PlanningStage` identity or review verdict ownership.
- Conformance Review remains an independent analytical verdict. Its findings and recommended checks are QA map inputs, but QA synthesis cannot improve a failed or blocked current review.
- The QA map fingerprint determines whether saved shard and theory checkpoints are resumable. Changed governed inputs, execute evidence, target identity, map schema, or investigator policy must make incompatible state non-current.
- Detailed QA state is authoritative for map, shard, theory, attempt, and synthesis progress. `flow-state.json` is only a bounded projection and pointer to the current canonical QA record.
- Investigator runtime requests pass through `internal/platform/runtime`, which translates product-owned sandbox and permission policy into Agentwrap and reports unsupported enforcement.
- Product QA execution belongs in `internal/sprint`; operation preparation and shared dispatch belong in `internal/app`; CLI, TUI, and web remain presentation and transport adapters.
- Runtime-backed QA acceptance, ownership, events, cancellation, and terminal outcomes flow through `durableOperationManager` and `runcontrol.Repository`.
- Web's operation hub and TUI's in-memory event list are bounded presentation caches. Durable QA progress and recovery derive from QA state plus run control, not from either adapter cache.
- `review` and `review.md` remain compatibility contracts while labels change to Conformance Review. Smoke remains a separate existing verification stage until Sprint 37 wraps it as a QA suite.

## Constraints

- Do not add QA to the canonical `PlanningStages()` sequence or use `PlanningStage` as the verification-phase type.
- Preserve the `review` command, `review.md`, existing JSON fields, and automation behavior while changing user-facing terminology to Conformance Review.
- Keep smoke behavior and `smoke.md` unchanged in Sprint 36; smoke-as-QA integration belongs to Sprint 37.
- Keep detailed QA state out of `flow-state.json`; flow state may contain only canonical summary, freshness, verdict, and a contained pointer.
- Use explicit schema versions and deterministic, verification-scoped identifiers with defined compatibility behavior.
- Sort all unordered inputs before hashing or persistence so unchanged inputs produce byte-stable maps and identities.
- Every changed path must map to a bounded verification surface; unknown or unclassifiable paths must become visible blocked work rather than being silently omitted.
- Investigator fan-out, shard count, theories, findings, prompt bytes, output bytes, event bytes, concurrency, and follow-up requests must be bounded.
- Investigators may read approved context, formulate falsifiable theories, record confirmed/refuted/inconclusive outcomes, and recommend checks.
- Investigators may not create or modify tests, production code, verification code, Git state, governed sprint artifacts, or runtime configuration.
- The runtime request must use an enforced read-only sandbox and default-deny permission policy. Unsupported permission enforcement is blocked, never degraded to advisory prompt safety.
- Static QA output must be validated before promotion. Malformed output, cancellation, runtime failure, stale inputs, or atomic rename failure must preserve the prior valid canonical state and summary.
- Retain refuted and invalid theories in synthesis; investigators cannot promote their own theories into repairable issues.
- No automatic repair, issue promotion, production mutation, or verification-code mutation belongs in this sprint.
- Use the existing app operation boundary and durable run-control repository. Do not add adapter-owned queues, operation histories, cancellation registries, or workflow authority.
- Web production code must continue importing only `internal/app` and standard-library packages.
- Browser commands remain allowlisted, strictly decoded, confirmation-bound, fingerprint-revalidated, same-origin, and CSRF-protected.
- CLI, TUI, and web must expose equivalent QA status, progress, blocking reasons, synthesis results, cancellation, and durable recovery through shared app DTOs.
- Offline fake-runtime tests must prove deterministic mapping, permission denial, malformed output handling, cancellation, resume compatibility, atomic failure preservation, synthesis stability, and adapter parity.
- Release validation requires `go test ./...`, `go test -race ./...`, `go vet ./...`, `go build ./cmd/ultraplan`, and `git diff --check`.

## Open Questions

- The repository contains no existing QA artifact path, detailed-state filename, schema, or migration implementation. The selected naming must be introduced consistently across sprint containment, previews, status projections, and recovery documentation.
- The current read-only reviewer policy allows only read, list, and search tools. The source does not yet define a QA-specific allowlist for safe non-mutating checks, command classification, or denial reporting.
- Current `FlowState` embeds detailed review and smoke state. There is no existing pointer-only verification record that directly demonstrates loading, validating, and reconciling a separate detailed state artifact.
- The current verification assessment assumes review plus smoke. The source does not yet encode how read-only QA contributes to an overall assessment before Sprint 37 integrates smoke under QA.
- Prior Sprint 35 governed artifacts were not present in this checkout. Its delivered durable-run behavior is visible in current source and tests, but artifact-specific decisions cannot be cross-checked here.
