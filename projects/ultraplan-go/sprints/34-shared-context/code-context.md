# Sprint Code Context

## Sprint Scope

Sprint 34 introduces one deterministic `internal/sprint`-owned shared prompt prefix built from the stored requirements, stored reference-only code context, and transiently resolved source evidence. The prefix must be reused across planning, execute, review, and smoke agent requests while keeping stage-specific and dynamic execution data after an explicit boundary. The sprint also extends cumulative flow ordering and registers a manual-only code-context skill without changing the generic runtime contract or Sprint 33's atomic generation and recovery behavior.

## Inspected Repository Areas

- `internal/sprint` stage definitions, artifact storage, prompt rendering, cumulative orchestration, planning-stage call sites, execute, review fan-out, smoke authoring, validation, and code-context promotion.
- `internal/platform/runtime` generic request and prompt-identity boundary.
- `internal/workspace` stage-skill definitions, rendering, selection, dry-run planning, materialization, customization preservation, and force behavior.
- Existing sprint tests covering code-context schema validation, read-only execution, failed-rerun preservation, cancellation and failure outcomes, planning dry runs, flow state, review request capture, smoke behavior, and skill synchronization.

## Selected Source References

### Planning Stage Domain Order

- **Path:** `internal/sprint/domain.go`
- **Lines:** `24-59`
- **Symbol:** `PlanningStage`
- **Rationale:** Defines sprint identity, the canonical planning-stage constants, and durable stage status shape. The shared-context work must preserve this ordering model and place code-context between requirements and sprint-index.

### Stored Planning Inputs

- **Path:** `internal/sprint/store_fs.go`
- **Lines:** `23-29`
- **Symbol:** `PlanningInputs`
- **Rationale:** Shows that exact requirements and code-context strings are already first-class planning inputs, providing the natural byte-preserving inputs for the shared renderer.

### Planning Input Reads

- **Path:** `internal/sprint/store_fs.go`
- **Lines:** `95-151`
- **Symbol:** `FSStore.ReadPlanningInputs`
- **Rationale:** Reads requirements and code-context directly as file bytes converted to strings and resolves project inputs through contained workspace paths. This is the existing storage boundary the renderer should reuse rather than introducing another manifest or persistence mechanism.

### Code-Context Artifact Contract

- **Path:** `internal/sprint/code_context.go`
- **Lines:** `18-112`
- **Symbol:** `ValidateCodeContextContent`
- **Rationale:** Defines the accepted reference-only Markdown grammar, 64 KiB limit, required sections, mandatory path/range/rationale fields, and fenced-content rejection. Shared source resolution must consume this established format without regenerating it.

### Reference Path and Range Validation

- **Path:** `internal/sprint/code_context.go`
- **Lines:** `133-203`
- **Symbol:** `validateRepositoryRelativePath`
- **Rationale:** Contains the current entry parsing and lexical safety rules for repository-relative paths and inclusive line ranges. The shared renderer needs compatible parsing while adding live repository containment, existence, range, and prompt-budget checks.

### Code-Context Prerequisite State

- **Path:** `internal/sprint/code_context.go`
- **Lines:** `206-255`
- **Symbol:** `Service.codeContextPrerequisite`
- **Rationale:** Establishes that downstream readiness requires both a valid artifact and a persisted successful code-context outcome. Shared-context integration must not weaken this truthful-state prerequisite.

### Code-Context Runtime Boundary

- **Path:** `internal/sprint/code_context.go`
- **Lines:** `257-375`
- **Symbol:** `Service.FlowCodeContext`
- **Rationale:** Implements generation through a read-only implementation-repository runtime, temporary candidate, bounded repair, validation, and promotion. It also identifies the resolved target and current permission policy that later transient source resolution must respect.

### Atomic Promotion and Recovery

- **Path:** `internal/sprint/code_context.go`
- **Lines:** `411-489`
- **Symbol:** `Service.promoteCodeContext`
- **Rationale:** Preserves failure outcomes and restores the previous valid artifact if state persistence fails. Sprint 34 must reuse the resulting artifact unchanged and avoid disturbing these atomic rerun guarantees.

### Current Prompt Composition Helpers

- **Path:** `internal/sprint/prompts.go`
- **Lines:** `14-61`
- **Symbol:** `RenderCodeContextPrompt`
- **Rationale:** Defines prompt previews and the current code-context generation prompt, including exact requirements injection and reference-only constraints. Shared downstream helpers should remain distinct from this generation-stage prompt.

### Downstream Stage Renderers

- **Path:** `internal/sprint/prompts.go`
- **Lines:** `63-175`
- **Symbol:** `RenderSprintIndexPrompt`
- **Rationale:** Contains all planning-stage prompt renderers and their stage-specific manifests, templates, mutation constraints, and output contracts. These are the principal composition points that must move behind one shared prefix without losing stage-specific behavior.

### Default Prompt and Override Resolution

- **Path:** `internal/sprint/prompts.go`
- **Lines:** `177-262`
- **Symbol:** `renderPromptFromDefault`
- **Rationale:** Shows how workspace and project overrides are resolved and how runtime manifests are currently appended. Shared-prefix composition must account for these existing override semantics while making the common/stage-specific boundary explicit and deterministic.

### Planning Service Call Sites

- **Path:** `internal/sprint/service.go`
- **Lines:** `294-312`
- **Symbol:** `Service.PromptSprintIndex`
- **Rationale:** Shows preview call sites for sprint-index and technical-handbook. They currently invoke stage renderers directly and therefore need shared-context preparation before rendering.

### Reasoning and Plan Call Sites

- **Path:** `internal/sprint/service.go`
- **Lines:** `416-453`
- **Symbol:** `Service.PromptAreaReasoning`
- **Rationale:** Covers area reasoning, final reasoning, and plan previews, including the no-template skip case. Compatible agent-backed paths require the common prefix, while a genuinely skipped area stage must remain runtime-free.

### Runtime Request Metadata Boundary

- **Path:** `internal/sprint/service.go`
- **Lines:** `954-1004`
- **Symbol:** `Service.runtimeRequest`
- **Rationale:** Adds stage, trace, checksum, model override, and progress metadata after prompt construction. This is the key evidence that dynamic execution identity can remain runtime metadata rather than entering the stable shared prefix.

### Execute Prompt Integration Point

- **Path:** `internal/sprint/execute.go`
- **Lines:** `93-185`
- **Symbol:** `Service.Execute`
- **Rationale:** Shows execute preview, dry-run, per-task prompt construction, runtime metadata, target working directory, and session-resume preamble. The shared prefix must wrap every task prompt while keeping task identity, model source, attempt/session behavior, and resume text stage-specific.

### Execute Stage Instructions

- **Path:** `internal/sprint/execute.go`
- **Lines:** `356-381`
- **Symbol:** `RenderExecutePrompt`
- **Rationale:** Defines task-specific traceability, implementation steps, evidence, safety, and deferral instructions. These belong after the shared-context boundary and must remain intact.

### Cumulative Flow Scheduler

- **Path:** `internal/sprint/flow.go`
- **Lines:** `44-112`
- **Symbol:** `Service.Flow`
- **Rationale:** Owns cumulative scheduling, exact-stage skipping, mutation locking, progress, and verification handoff. Exact-once code-context execution for `flow --to plan` must be implemented and tested at this product-owned orchestration boundary.

### Stage Dispatch and Canonical Sequence

- **Path:** `internal/sprint/flow.go`
- **Lines:** `145-200`
- **Symbol:** `flowStages`
- **Rationale:** Maps each stage to its operation and slices the canonical sequence for cumulative targets. It currently encodes the required requirements-to-plan order with code-context in the second position.

### Flow State Transitions

- **Path:** `internal/sprint/flow.go`
- **Lines:** `243-315`
- **Symbol:** `flowCodeContextSuccessStages`
- **Rationale:** Defines readiness propagation and failure-state construction across planning stages. Shared prompt integration must not alter these durable state transitions.

### Review Preparation and Fan-Out Inputs

- **Path:** `internal/sprint/review.go`
- **Lines:** `166-369`
- **Symbol:** `Service.PrepareReview`
- **Rationale:** Builds the frozen review manifest, governed inputs, independent coverage set, target identity, changed paths, and fingerprint. Shared context must be supplied to each reviewer without changing coverage ownership, sorting, fingerprinting, or verdict inputs.

### Review Worker Requests

- **Path:** `internal/sprint/review.go`
- **Lines:** `833-939`
- **Symbol:** `Service.runReviewer`
- **Rationale:** Constructs and executes each independent reviewer request with read-only permissions, validation, session continuation, and bounded fallback repair. This is the per-agent integration point for the common prefix.

### Review Prompt Boundary

- **Path:** `internal/sprint/review.go`
- **Lines:** `1446-1489`
- **Symbol:** `renderReviewerPrompt`
- **Rationale:** Separates human preview from the actual reviewer prompt and renders coverage-specific frozen paths and output schema. Shared evidence should precede these request-specific review instructions without changing fan-out or citation semantics.

### Smoke Author Runtime Request

- **Path:** `internal/sprint/smoke_author.go`
- **Lines:** `20-106`
- **Symbol:** `Service.authorSmokeSuite`
- **Rationale:** Captures the only agent-backed smoke-authoring invocation, its harness-only mutation policy, protected snapshots, cancellation handling, and post-run checks. The same common prefix must be added here without weakening smoke isolation or recovery.

### Smoke Author Stage Instructions

- **Path:** `internal/sprint/smoke_author.go`
- **Lines:** `198-279`
- **Symbol:** `Service.renderSmokeAuthorPrompt`
- **Rationale:** Builds the dynamic smoke manifest, governed input paths, writable paths, and detailed authoring contract. These values are request-specific and therefore belong after the stable shared prefix.

### Generic Runtime Contract

- **Path:** `internal/platform/runtime/runtime.go`
- **Lines:** `1-59`
- **Symbol:** `Request`
- **Rationale:** Defines the generic runtime request, prompt identity, permissions, metadata, validation, and event callback contract. It must remain free of sprint-specific shared-context types or composition logic.

### Stage Skill Registry

- **Path:** `internal/workspace/skills.go`
- **Lines:** `13-37`
- **Symbol:** `StageSkill`
- **Rationale:** Defines embedded skill metadata and the registry entry point. The manual code-context skill should be added here using existing fields rather than introducing a new skill framework.

### Existing Planning Skill Definitions

- **Path:** `internal/workspace/skills.go`
- **Lines:** `82-159`
- **Symbol:** `StageSkills`
- **Rationale:** Shows conventions for names, prerequisites, prompt availability, and workflows. Code-context differs because it must delegate execution to the canonical CLI operation instead of asking the invoking agent to recreate stage logic.

### Skill Selection and Materialization

- **Path:** `internal/workspace/skills.go`
- **Lines:** `200-320`
- **Symbol:** `MaterialiseSkills`
- **Rationale:** Implements selection by stage or full skill name, dry-run planning, deterministic file ordering, customization preservation, force overwrite, and materialization. The new skill must participate in these existing semantics and all-skill output.

### Manual-Only Skill Rendering

- **Path:** `internal/workspace/skills.go`
- **Lines:** `323-393`
- **Symbol:** `renderStageSkill`
- **Rationale:** Generates `SKILL.md` and OpenAI metadata, including the default prohibition on implicit invocation and general no-delegation rule. Code-context needs a narrow explicit delegation exception synchronized across generated content and metadata.

### Code-Context Validation Tests

- **Path:** `internal/sprint/code_context_test.go`
- **Lines:** `17-109`
- **Symbol:** `TestValidateCodeContextContent`
- **Rationale:** Provides the canonical valid fixture and negative coverage for placeholders, unsafe paths, malformed ranges, missing fields, copied source, optional symbols, multiple ranges, and size limits. These fixtures can underpin shared renderer parsing tests.

### Code-Context Rerun and Permission Tests

- **Path:** `internal/sprint/code_context_test.go`
- **Lines:** `193-314`
- **Symbol:** `TestCodeContextPromptDryRunExecutionAndRerunPreservation`
- **Rationale:** Proves dry-run non-mutation, read-only target policy, prompt identity, failed-rerun preservation, bounded repair, and no mutation of the implementation repository or unrelated artifacts.

### Code-Context Failure and Recovery Tests

- **Path:** `internal/sprint/code_context_test.go`
- **Lines:** `413-507`
- **Symbol:** `TestCodeContextMissingOutputUnsupportedPermissionsAndCancellationFailClosed`
- **Rationale:** Covers runtime failure, candidate cleanup, missing output, timeout, unsupported permissions, cancellation, interruption, cleanup uncertainty, persisted outcome truthfulness, and restoration after state-write failure.

### Planning Flow Test Seam

- **Path:** `internal/sprint/sprint_index_test.go`
- **Lines:** `36-111`
- **Symbol:** `TestPromptPreviewAndFlowDryRunAreRuntimeFree`
- **Rationale:** Demonstrates prompt capture through fake runtimes, dry-run non-mutation, code-context prerequisites, and stage validation. It is the natural location for common-prefix and cumulative planning assertions.

### Cumulative Materialization Test

- **Path:** `internal/sprint/sprint_index_test.go`
- **Lines:** `138-173`
- **Symbol:** `TestCumulativeFlowMaterializesMissingSprintBeforeMutationLock`
- **Rationale:** Covers cumulative-flow startup and missing-sprint creation. It should be extended with a runtime call log proving code-context executes exactly once in the required sequence.

### Skill Materialization Tests

- **Path:** `internal/workspace/skills_test.go`
- **Lines:** `10-103`
- **Symbol:** `TestMaterialiseAllStageSkills`
- **Rationale:** Verifies dry-run non-mutation, all-skill synchronization, generated content, manual-only metadata, idempotence, customization preservation, force restoration, and single-skill isolation. The expected skill count and code-context-specific delegation assertions must be updated here.

## Relationships

`FSStore.ReadPlanningInputs` supplies exact persisted requirements and code-context bytes to `Service`. A new renderer owned by `internal/sprint` can validate and parse references using the established code-context grammar, resolve them against the implementation target selected through existing target resolution, and compose one deterministic common prefix. Existing stage renderers in `prompts.go`, execute task rendering, reviewer rendering, and smoke author rendering then provide only their stage-specific suffixes.

`Service.runtimeRequest` remains the transition from product-owned prompt composition to the generic runtime request. Dynamic trace IDs, stage metadata, model selection, sessions, task IDs, coverage IDs, and attempts should remain at this boundary or in stage-specific suffixes.

`Service.Flow` and `flowStages` remain authoritative for cumulative scheduling and state. `FlowCodeContext` remains authoritative for generation, validation, cancellation, and atomic artifact promotion. The workspace skill must invoke that canonical flow operation rather than reproduce target resolution, prompting, validation, or persistence.

Review retains its independent coverage fan-out and deterministic product-owned aggregation. Smoke retains harness-only write authority and protected-target checks. Neither subsystem should acquire alternate shared-context storage or ownership.

## Constraints

- Preserve the stored `requirements.md` and `code-context.md` bytes exactly in the shared prefix.
- Parse code-context only to resolve references; never rewrite it or persist resolved source excerpts.
- Resolve repository-relative paths inside the configured implementation target, reject escapes, missing files, and invalid or out-of-range lines, preserve selected-entry order, and enforce a bounded prompt budget.
- Identify resolved snippets as transient prepared evidence and explicitly permit further live repository inspection.
- Keep stable instructions, sprint identity, exact artifacts, resolved evidence, and other shared context deterministic and before the stage-specific boundary.
- Keep stage names, output paths, task and coverage identities, attempts, run/session IDs, timestamps, and other request-specific data after that boundary or in runtime metadata.
- Keep all sprint semantics and prompt composition in `internal/sprint`; `internal/platform/runtime` remains generic.
- Preserve code-context prerequisite validation, read-only generation, bounded repair, candidate cleanup, truthful failure outcomes, and last-valid-artifact restoration.
- Preserve review fan-out, frozen-input behavior, aggregation, and verdict ownership.
- Preserve smoke harness mutation limits, protected snapshots, cancellation, and recovery behavior.
- Reuse existing flow, application, CLI, TUI, and web operations rather than introducing route-specific or adapter-specific prompt logic.
- Add no repository index, retrieval system, embeddings, cache, parallel manifest, staleness detector, provenance system, QA stage, or new workflow framework.
- Keep deterministic tests offline; any real-runtime dogfood remains explicitly gated and truthfully reports missing prerequisites.

## Open Questions

- The exact shared “other context” content beyond sprint identity, requirements, code-context, and resolved source evidence is not represented by one current abstraction; downstream implementation should identify only existing stable inputs that genuinely need to be common.
- The runtime prompt budget below the existing 96 KiB transport ceiling needs an explicit deterministic allocation for resolved snippets while retaining the 64 KiB durable artifact limit.
