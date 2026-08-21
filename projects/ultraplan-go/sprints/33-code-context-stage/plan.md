# Sprint Plan: Code-Context Stage Vertical Slice

> Project: `ultraplan-go`
> Sprint: `33-code-context-stage`
> Source: `reasoning.md`
> **Inputs Used:** `projects/ultraplan-go/project-index.md`, `projects/ultraplan-go/roadmap.md`, `projects/ultraplan-go/sprints/33-code-context-stage/requirements.md`, `projects/ultraplan-go/docs/ARCHITECTURE.md`, `projects/ultraplan-go/docs/PRD.md`, `projects/ultraplan-go/docs/TRD.md`, `projects/ultraplan-go/sprints/33-code-context-stage/sprint-index.md`, `projects/ultraplan-go/sprints/33-code-context-stage/technical-handbook.md`, `projects/ultraplan-go/sprints/33-code-context-stage/reasoning/api-design.md`, `projects/ultraplan-go/sprints/33-code-context-stage/reasoning/architecture.md`, `projects/ultraplan-go/sprints/33-code-context-stage/reasoning/frontend.md`, `projects/ultraplan-go/sprints/33-code-context-stage/reasoning.md`, `../ultraplan-go/docs/plans/integrated-roadmap.md`, `../ultraplan-go/docs/plans/sprint-code-context-stage.md`, `../ultraplan-go/docs/plans/ultraplan-local-server-experiment-plan.md`, `../ultraplan-go/docs/plans/server-shutdown-run-cancellation-contract.md`, `../ultraplan-go/docs/phase3-json-schemas.md`, `../ultraplan-go/docs/web-compatibility-baseline.md`, `../ultraplan-go/internal/sprint/domain.go`, `../ultraplan-go/internal/sprint/artifacts.go`, `../ultraplan-go/internal/sprint/store_fs.go`, `../ultraplan-go/internal/sprint/state.go`, `../ultraplan-go/internal/sprint/flow.go`, `../ultraplan-go/internal/sprint/service.go`, `../ultraplan-go/internal/sprint/prompts.go`, `../ultraplan-go/internal/sprint/execute.go`, `../ultraplan-go/internal/sprint/execute_target.go`, `../ultraplan-go/internal/sprint/execute_model.go`, `../ultraplan-go/internal/sprint/locks.go`, `../ultraplan-go/internal/sprint/review.go`, `../ultraplan-go/internal/sprint/smoke.go`, `../ultraplan-go/internal/platform/runtime/runtime.go`, `../ultraplan-go/internal/platform/config/config.go`, `../ultraplan-go/internal/platform/config/redaction.go`, `../ultraplan-go/internal/workspace/init.go`, `../ultraplan-go/internal/workspace/defaults.go`, `../ultraplan-go/internal/app/config_commands.go`, `../ultraplan-go/internal/app/sprint_commands.go`, `../ultraplan-go/internal/app/sprint_usecases.go`, `../ultraplan-go/internal/app/usecases.go`, `../ultraplan-go/internal/app/web_usecases.go`, `../ultraplan-go/internal/app/operations.go`, `../ultraplan-go/internal/app/operation_runner.go`, `../ultraplan-go/internal/app/json_output.go`, `../ultraplan-go/internal/web/handlers.go`, `../ultraplan-go/internal/web/operations.go`, `../ultraplan-go/internal/web/operation_handlers.go`, `../ultraplan-go/internal/web/templates/sprint.html`, and the corresponding existing test files named in the tasks below.

This plan executes `reasoning.md`. It does not reopen architecture or authorize Sprint 34 shared-prefix, skill, broad-documentation, or real-runtime release work.

## Reasoning Source

- **Sprint Reasoning:** `reasoning.md`
- **Sprint Index:** `sprint-index.md`
- **Technical Handbook:** `technical-handbook.md`
- **Area Reasoning:** `reasoning/api-design.md`, `reasoning/architecture.md`, `reasoning/frontend.md`
- **Roadmap Scope:** `../../roadmap.md`, Sprint 33 section
- **Detailed Implementation Inputs:** the four implementation plans listed in **Inputs Used**
- **Omitted Evidence:** The 12 selected study final reports were not reopened because `technical-handbook.md` and the three area-reasoning documents contain sufficient concrete findings and source references for every final decision. No additional report was needed to resolve an implementation question.

## Sprint Status

- **Status:** `complete`
- **Owner:** `manual in-session execution`
- **Start Date:** `2026-08-19`
- **Completion Date:** `2026-08-19`

## Decisions To Execute

| Decision | Source Section | Execution Implication |
| --- | --- | --- |
| Canonical Stage Order And Legacy State Compatibility | `reasoning.md#decision-1-canonical-stage-order-and-legacy-state-compatibility` | Insert `StageCodeContext` exactly once after requirements in canonical and cumulative order, gate sprint-index on validated completion, and interpret pre-stage flow state without writes on read-only paths. |
| Sprint-Owned Sequential Runtime Service And Atomic Promotion | `reasoning.md#decision-2-sprint-owned-sequential-runtime-service-and-atomic-promotion` | Implement one focused `internal/sprint` service that resolves the existing target, runs one generic runtime request, validates an isolated candidate, and atomically promotes only `code-context.md`. |
| Single Markdown Artifact And Structural Validation Contract | `reasoning.md#decision-3-single-markdown-artifact-and-structural-validation-contract` | Ship one embedded Markdown template and one central structural validator; reject malformed or escaping references without claiming semantic completeness. |
| Source-Aware Configuration And Additive CLI/App Contract | `reasoning.md#decision-4-source-aware-configuration-and-additive-cliapp-contract` | Extend existing model/variant precedence, stage-valued commands, stable JSON, synchronous reads, and asynchronous generation without caller-provided target/output paths. |
| Generic Browser Operation And Truthful Sprint Presentation | `reasoning.md#decision-5-generic-browser-operation-and-truthful-sprint-presentation` | Reuse the existing `sprint-stage` operation, confirmation, SSE, cancellation, recovery, and sprint page; represent artifact validity separately from latest attempt outcome. |
| Canonical Cancellation, Stable Failure Identity, And Safe Observability | `reasoning.md#decision-6-canonical-cancellation-stable-failure-identity-and-safe-observability` | Carry one context through runtime work, reuse canonical cancellation and terminal arbitration, record uncertainty truthfully, and expose only bounded allowlisted metadata. |
| Layered Offline Verification And Scope Enforcement | `reasoning.md#decision-7-layered-offline-verification-and-scope-enforcement` | Build deterministic evidence from domain tests through boundary tests and repository-wide gates; defer live-runtime proof and reject all named scope expansion. |

## Decision Guardrails

| Decision | Requirement / Evidence | Trade-Off Accepted | Alternative Rejected | Risk / Follow-Up |
| --- | --- | --- | --- | --- |
| Fixed canonical stage plus legacy interpretation | `AC-01`, `AC-02`, `AC-03`, `AC-08`, `AC-10`; handbook additive-extension and migration evidence | Explicit projections need parity tests and compatibility remains in the load path. | Dynamic registry, second state format, read-time migration, or artifact-presence completion. | Prevent projection drift and prove status does not rewrite legacy bytes. |
| Sequential sprint-owned service and atomic promotion | `AC-04` through `AC-06`, `AC-08`, `AC-09`, `AC-14`; architecture reasoning | `internal/sprint` grows and candidate cleanup adds complexity. | New product module, runtime-owned semantics, direct authoritative writes, broad filesystem abstraction, or stage fan-out. | Verify permission enforcement and promotion/state ordering before treating the stage as complete. |
| One structurally validated Markdown artifact | `AC-06` through `AC-09`, `AC-11`, `AC-18`; security and validation evidence | Structural validity cannot prove semantic completeness. | JSON manifest, index/RAG/cache, semantic scoring, arbitrary excerpt cap, or source mutation. | Keep parser rules semantic and report empty findings as structural only. |
| Additive config and CLI/app contract | `AC-03` through `AC-05`, `AC-12`; source-aware config evidence | Supplied/not-supplied state and stable projection tests add work. | Top-level command, caller paths, eager runtime construction, or raw diagnostic projection. | Preserve version-1 envelope, nullability, exit/error behavior, and explicit source truth. |
| Generic browser operation and richer typed projection | `AC-12` through `AC-14`; Sprint 32 web baseline | View models gain fields and confirmation remains an explicit operator step. | Dedicated route/page, operation kind, browser store, automatic rerun, cancel-on-disconnect, or raw logs. | Keep no-JavaScript operation complete and distinguish preserved artifact from failed rerun. |
| Canonical cancellation and safe observability | `AC-08`, `AC-09`, `AC-12` through `AC-14`; shutdown contract | Bounded cleanup can delay return and uncertainty adds visible states. | Detached work, fire-and-forget cleanup, string parsing, transport retry, or inferred success. | Race completion against cancellation under `-race`; add a stable class only if current identity is insufficient. |
| Layered offline evidence | `AC-01` through `AC-18`; Testing contract and review protocols | Temporary repositories and boundary tests cost more than helper-only tests. | Live-provider normal tests, blanket snapshots, in-memory-only filesystem proof, or simulated dogfood. | Real-runtime dogfood remains explicitly deferred to Sprint 34. |

## Requirements / Contracts To Satisfy

Acceptance IDs are assigned in the order of `requirements.md` acceptance bullets; constraint IDs are assigned in the order of its constraints.

| Contract / Requirement ID | Required Behavior | Evidence Planned |
| --- | --- | --- |
| `AC-01` | Canonical ordered lists contain `code-context` exactly once after requirements and before sprint-index. | Exact order, membership, count, status, and artifact-map tests. |
| `AC-02` | Requirements gate code-context; validated successful code-context gates sprint-index. | Readiness, missing/invalid/failed prerequisite, and cumulative-flow tests. |
| `AC-03` | Prompt, validate, flow, help, status, typed app requests, and stable JSON accept the stage. | CLI/app tests and additive JSON compatibility fixtures. |
| `AC-04` | Prompt preview and dry-run are deterministic and non-mutating. | Runtime-spy plus artifact/state/repository before-and-after tests. |
| `AC-05` | Existing target resolution and generic runtime receive caller context and effective stage overrides. | Request-construction, target-resolution, precedence, and cancellation tests. |
| `AC-06` | Repository is readable, but only sprint-root `code-context.md` is writable. | Temporary repository and Git-state snapshots plus output-scope assertions. |
| `AC-07` | A valid artifact has all required sections and at least one well-formed exact source excerpt. | Table-driven valid fixtures and representative artifact review. |
| `AC-08` | Validation rejects every named malformed, unsafe, empty, or placeholder case actionably. | Table-driven invalid fixtures with finding-code/message assertions. |
| `AC-09` | Runtime success without valid output fails and cannot unlock sprint-index. | Fake-runtime missing/invalid output and state-transition tests. |
| `AC-10` | Explicit rerun atomically replaces only the artifact and preserves the last valid copy on failure. | Success, rename/write failure, validation failure, cancellation, and cleanup tests. |
| `AC-11` | Pre-code-context state remains readable and read-only status does not migrate it. | Legacy schema-2 fixtures, preserved outcomes, unchanged bytes/mtime, and later mutation serialization. |
| `AC-12` | Embedded prompt/template and materialized defaults are synchronized. | Embedded lookup, defaults plan/install, and golden parity tests. |
| `AC-13` | Public runtime/config/status metadata is stable, bounded, source-aware, and redacted. | Projection, redaction, bounds, unknown-usage, and negative disclosure tests. |
| `AC-14` | Browser uses generic readiness/progress/findings/preview/rerun/cancel/recovery behavior. | App/web operation contract, template, no-JavaScript, SSE, conflict, and recovery tests. |
| `AC-15` | Disconnect does not cancel; explicit and shutdown cancellation reach canonical ownership. | Disconnect, explicit cancel, shutdown, terminal race, and uncertainty tests. |
| `AC-16` | Focused deterministic tests cover the complete vertical slice. | Tests in every required output file and focused package commands. |
| `AC-17` | Repository-wide test, race, vet, build, and whitespace gates pass. | Exact release commands in **Verification Commands**. |
| `AC-18` | No index, RAG, cache, manifest, staleness, amendment, source mutation, or later-sprint system is introduced. | Diff/dependency review under both required protocols. |
| `C-01` Architecture ownership | `internal/sprint` owns product semantics; platform runtime stays generic. | Import/dependency review and runtime request tests. |
| `C-02` Shared application boundary | CLI and web call typed app use cases; web owns no product workflow. | App/web contract tests and Architecture Review. |
| `C-03` Filesystem authority | Artifacts and sprint state remain authoritative; web operation state stays ephemeral. | Refresh/reconnect/restart tests and dependency review. |
| `C-04` Path containment | Target is resolved internally and artifact source paths are relative and contained. | Target and hostile path tables; preview symlink/escape tests. |
| `C-05` Mutation boundary | Source, tests, Git, governed inputs, and unrelated sprint artifacts remain unchanged. | Before/after temporary repository and sprint-root snapshots. |
| `C-06` Atomic completion | Candidate promotion and state writes are atomic; completion follows validation. | Injected write/rename failures and transition-order assertions. |
| `C-07` Structural claims only | Validation does not claim semantic completeness. | Finding/help/template wording assertions and manual review. |
| `C-08` Runtime contract | Context, typed errors, progress, cancellation, redaction, and metadata use existing generic seams. | Fake runtime and shared operation tests under race detection. |
| `C-09` Offline verification | Normal tests require no live provider; unavailable gated evidence is deferred. | Fake runtime enforcement and explicit Sprint 34 deferral. |
| `C-10` Roadmap authority | Sprint 33 does not absorb Sprint 34 or post-Phase-5 scope. | Scope checklist and protocol reviews against `roadmap.md`. |
| Architecture; CLI Surface; Configuration; Documentation; Errors; LLM Evaluation / Cost / Safety; LLM Runtime; Observability; Persistence And Migrations; Security; Testing; Workflows | Apply all 12 selected contracts as mapped by `reasoning.md#contracts-applied`; no selected contract is silently dropped. | Task-level traceability, focused tests, release evidence, Architecture Review, and Sprint Review. |

## Tasks

- [x] **Task 1: Establish Canonical Stage Order And Compatible State**
  > Executes: Decisions 1 and 7; `AC-01`, `AC-02`, `AC-03`, `AC-11`, `C-01`, `C-03`, `C-06`
  - [x] Add `StageCodeContext` after `StageRequirements` in `internal/sprint/domain.go`; update `PlanningStages()`, validity, stage count, artifact mapping in `artifacts.go`, snapshots in `store_fs.go`, and every fixed cumulative/order projection in `flow.go`.
  - [x] Make valid completed requirements yield code-context readiness and only a successful validated code-context outcome yield sprint-index readiness; update success/failure stage builders and exact-once cumulative dispatch without adding a registry.
  - [x] Extend `state.go` to recognize persisted current-schema state created before the inserted stage, preserve all known outcomes, derive the inserted stage deterministically, perform no load/status write, and serialize the canonical representation only on a later mutation.
  - [x] Add representative legacy fixtures and assertions in `verify_test.go`; update `sprint_test.go` and `sprint_index_test.go` for exact order, mapping, readiness, status, cumulative position, exact-once dispatch, unchanged legacy bytes/mtime, and atomic later writes.
  - [x] Stop if compatibility would require a second state format or hidden read-time migration; record the deviation before proceeding.

- [x] **Task 2: Define The Artifact, Defaults, And Structural Validator**
  > Executes: Decision 3; `AC-06`, `AC-07`, `AC-08`, `AC-12`, `C-04`, `C-05`, `C-07`
  - [x] Add `internal/workspace/scaffold/prompts/create-code-context.md` with requirements-driven repository exploration, exact-source selection, read-only source policy, single-output policy, no-design/no-plan instruction, and permission to inspect enough source to prepare the pack.
  - [x] Add `internal/workspace/scaffold/templates/code-context.md` with scope, inspected areas, selected code, relationships, constraints, and open questions; each selected entry must support path, optional range/symbol, rationale, and a language-tagged fence.
  - [x] Register both assets in `internal/workspace/init.go` so embedded lookup and `defaults install` use the existing deterministic mechanism; do not change `init-workspace` export behavior.
  - [x] Implement central content/file validation with the code-context service in `internal/sprint/code_context.go`: reject absent/empty output, placeholders, missing sections/excerpts/path/rationale/fence/language, absolute or escaping paths, and malformed optional ranges; return actionable findings and make no semantic-completeness claim.
  - [x] Add table-driven validator tests in `code_context_test.go` for every accepted and rejected shape, including whitespace variations, multiple excerpts, optional symbols, hostile paths, and no arbitrary excerpt-count limit.
  - [x] Extend `workspace_test.go` and defaults-install tests for embedded/materialized prompt-template parity.

- [x] **Task 3: Implement Sequential Runtime Execution And Atomic Promotion**
  > Executes: Decisions 2 and 6; `AC-04`, `AC-05`, `AC-06`, `AC-09`, `AC-10`, `AC-13`, `AC-15`, `C-01`, `C-04`, `C-05`, `C-06`, `C-08`
  - [x] Add focused `PromptCodeContext`, `ValidateCodeContext`, and `FlowCodeContext` behavior in `internal/sprint/code_context.go`, following existing service patterns while keeping code-context policy out of `internal/platform/runtime`.
  - [x] Build deterministic preview inputs from validated requirements, sprint identity, selected scope/default assets, canonical output, and the target identified by existing `ResolveExecuteTarget`; reject caller-controlled target/output paths.
  - [x] Keep prompt, validate, status, readiness, preview, and dry-run paths free of runtime construction/invocation, candidate creation, artifact replacement, and flow-state mutation.
  - [x] For execution, resolve the existing target, effective model/variant, and required permission posture; pass one caller-owned context and generic request to the current runtime/event/wait path with repository-read and fixed-output policy.
  - [x] Generate to a contained same-directory candidate, drain events and wait, require runtime success plus candidate existence plus structural validity, then atomically rename over only `code-context.md` and persist the truthful state transition using existing atomic conventions.
  - [x] Preserve a previous valid artifact for runtime, missing-output, invalid-output, cancellation, cleanup, write, rename, or state-persistence failure; keep latest attempt outcome distinct and never infer success from old artifact presence.
  - [x] Reuse current stable error identities where sufficient and preserve wrapped causes; add no class unless recovery behavior cannot be represented. Bound cleanup and record `interrupted` or `cleanup_uncertain` on unresolved ownership.
  - [x] Add fake-runtime and real temporary-repository tests in `code_context_test.go` for success, runtime failure, missing/invalid output, overrides, metadata, event drain/wait, cancellation, timeout, conflict, atomic rerun, injected failures, candidate cleanup, and source/test/Git/governed-input/unrelated-artifact immutability.
  - [x] Stop before runtime launch if required permission restrictions cannot be enforced by current agentwrap integration; report the capability gap rather than weakening the mutation contract.

- [x] **Task 4: Add Source-Aware Configuration And CLI Wiring**
  > Executes: Decision 4; `AC-03`, `AC-04`, `AC-05`, `AC-13`, `C-02`, `C-08`
  - [x] Add code-context model/variant fields, workspace parsing, fixed precedence, source keys, effective validation, and defaults in `internal/platform/config/config.go` and embedded default configuration.
  - [x] Extend `redaction.go` and `internal/app/config_commands.go` so text and JSON config views include safe effective values and sources without exposing sensitive values.
  - [x] Extend `internal/app/sprint_commands.go` stage switches, accepted-stage/help text, status rendering, flow target parsing, prompt/validate dispatch, and `planningStageRuntime`; preserve whether model/variant flags were explicitly supplied.
  - [x] Keep the existing version-1 JSON envelope and field meanings/nullability; add only compatible stage/projection data and retain usage errors for unknown stages.
  - [x] Add config tests for workspace/env/CLI precedence, omitted overrides, fallback, invalid values, source tracking, and redaction; add command tests for help, prompt, validate, status, flow, dry-run non-mutation, runtime override propagation, stable errors, and one structured JSON document.

- [x] **Task 5: Extend Shared App Use Cases And Generic Operations**
  > Executes: Decisions 4, 5, and 6; `AC-03`, `AC-13`, `AC-14`, `AC-15`, `C-02`, `C-03`, `C-08`
  - [x] Add code-context to typed stage validation, sprint summaries, findings, artifact lists, and readiness in `internal/app/sprint_usecases.go` and `web_usecases.go` without moving product validation into app code.
  - [x] Add `code-context.md` to the existing bounded contained preview allowlist in `internal/app/usecases.go` and web artifact projections; retain opaque references, symlink containment, escaping, and current size limits.
  - [x] Extend `operations.go` governed-input fingerprinting and explicit prompt/validation dispatch for code-context; reuse `OperationStage`/`OperationStageDryRun`, `OperationRequest.Stage`, `FlowStage`, the sprint-path mutation lease, and existing stale-confirmation/conflict behavior.
  - [x] Add explicit typed projection fields needed to represent authoritative artifact availability/validity separately from latest attempt status, error, and next action; do not derive this distinction from template strings or ephemeral SSE state.
  - [x] Preserve bounded progress, redacted metadata, unknown usage, canonical cancellation, disconnect-as-unsubscribe, terminal arbitration, and restart reconciliation through existing operation paths.
  - [x] Extend `sprint_usecases_test.go`, `usecases_test.go`, `web_usecases_test.go`, and `web_operations_test.go` for shared-state parity, preview containment, prepare/run/fingerprint, stale confirmation, conflict, progress, explicit cancellation, disconnect, shutdown, cleanup uncertainty, recovery, and failed rerun with preserved valid output.

- [x] **Task 6: Present Code-Context Through The Existing Sprint Page**
  > Executes: Decision 5; `AC-14`, `AC-15`, `C-02`, `C-03`, `C-04`
  - [x] Extend handler DTO mapping and `internal/web/templates/sprint.html` so code-context appears in DOM workflow order after requirements and before sprint-index with readiness, structural findings, bounded artifact access, latest outcome, and next action.
  - [x] Reuse existing generic prepare/start/status/events/cancel forms and operation mapping for explicit initial run and rerun; add no route, operation kind, stage-specific JavaScript controller, browser store, or durable web state.
  - [x] Render valid preserved artifact and failed latest rerun as simultaneous facts; use distinct cancelled/interrupted/cleanup-uncertain labels and never use success styling for uncertainty.
  - [x] Keep complete server-rendered operation without JavaScript, semantic controls, visible labels, keyboard operation, meaningful live regions, escaped hostile content, bounded previews, and refresh of authoritative state after terminal/reconnect events.
  - [x] Extend `templates_test.go`, `operations_contract_test.go`, `operations_test.go`, `artifacts_test.go`, and `api_compatibility_test.go` for DOM order, all states, safe preview, explicit rerun, no-JavaScript flow, accessibility semantics, hostile Markdown, generic contract reuse, disconnect, reconnect, slow subscribers, cancellation, shutdown, and compatibility.

- [x] **Task 7: Run Layered Verification And Enforce Scope**
  > Executes: Decision 7; `AC-16`, `AC-17`, `AC-18`, `C-09`, `C-10`
  - [x] Run focused package tests after each task, then the race-sensitive operation/cancellation suite, then all repository-wide release commands from `../ultraplan-go`.
  - [x] Capture evidence that prompt preview, dry-run, validation, status, and config display do not invoke runtime or mutate artifacts/state/repository; capture temporary-repository and Git-state comparisons proving the execution mutation boundary.
  - [x] Review the implementation against `system/protocols/architecture-review-protocol.md`, recording ownership, dependency direction, generic runtime/web reuse, atomicity, compatibility, cancellation, and absence of parallel workflow/persistence abstractions.
  - [x] Review the implementation against `system/protocols/review-sprint-protocol.md`, recording every acceptance criterion, required output, command result, validation finding, and explicit exclusion.
  - [x] Confirm the diff contains no downstream shared-prefix injection, manual skill, broad documentation, real-runtime dogfood claim, index/RAG/embedding/cache/cache-key/provider-cache-control, JSON manifest, automatic freshness/amendment, content identity, QA/repair, alternate persistence, graph, cloud/Aren, source/Git mutation, generic stage/plugin/workflow framework, web product service, or package-global runner registration.
  - [x] Record any command or gated evidence not run as blocked/deferred with cause; never mark unavailable real-runtime evidence passed.

## Evidence Checklist

- [x] Tests prove canonical stage order, prerequisites, exact-once cumulative position, status, artifact mapping, and transitions.
- [x] Legacy fixtures prove preserved outcomes, no hidden status write, and canonical serialization on later mutation.
- [x] Validator tables prove all required valid and invalid Markdown/path/range cases with actionable findings.
- [x] Fake-runtime and temporary-repository tests prove execution, isolation, cancellation, atomic rerun, and failure truthfulness.
- [x] Config/default/CLI/app/JSON evidence proves source-aware additive behavior and non-mutating reads.
- [x] Generic app/web operation and browser evidence proves confirmation, conflict, progress, preview, cancellation, recovery, no-JavaScript use, and preserved-artifact/latest-attempt separation.
- [x] Runtime or diagnostic evidence is bounded, redacted, correlated, and truthful; unknown usage remains unknown.
- [x] Executable help and embedded/default-install assets are current; broad documentation remains deferred.
- [x] Deviations from `reasoning.md` are recorded before implementation continues.
- [x] Architecture Review and Sprint Review protocol evidence is complete.
- [x] Sprint 34 and post-Phase-5 exclusions are absent from the implementation.

## Verification Commands

Run all commands from `../ultraplan-go`.

| Check | Command | Expected Result |
| --- | --- | --- |
| Focused sprint domain and stage service | `go test ./internal/sprint -run 'Test.*(CodeContext|FlowState|DomainStages|CumulativeFlow|Mutation|Atomic)'` | Stage, validator, compatibility, runtime, isolation, and atomicity tests pass offline. |
| Full sprint package | `go test ./internal/sprint` | All sprint behavior remains green. |
| Configuration | `go test ./internal/platform/config -run 'Test.*(LoadPrecedence|Redact|CodeContext)'` | Code-context precedence, source tracking, validation, and redaction pass. |
| Embedded defaults | `go test ./internal/workspace -run 'Test.*(Embedded|Default)'` | Embedded and materialized prompt/template copies match. |
| CLI and shared app | `go test ./internal/app -run 'Test.*(Sprint|WebOperation|WebUseCases|Preview|ConfigShow)'` | CLI/app/config/preview/operation contracts pass. |
| Browser and HTTP contracts | `go test ./internal/web -run 'Test.*(Sprint|Operation|Template|Artifact|APICompatibility)'` | Generic operations, no-JavaScript UI, safe preview, and compatibility pass. |
| Package boundary suites | `go test ./internal/platform/config ./internal/workspace ./internal/app ./internal/web` | All directly changed boundary packages pass. |
| Focused race suite | `go test -race ./internal/sprint ./internal/app ./internal/web` | No races or goroutine ownership failures in stage/operation/cancellation paths. |
| Repository tests | `go test ./...` | All tests pass offline. |
| Repository race tests | `go test -race ./...` | Full race suite passes. |
| Static analysis | `go vet ./...` | No vet findings. |
| Binary build | `go build ./cmd/ultraplan` | CLI builds successfully. |
| Whitespace and patch integrity | `git diff --check` | No whitespace errors. |

## Risks And Blockers

| Risk / Blocker | Source | Mitigation | Status |
| --- | --- | --- | --- |
| Canonical stage projections drift across domain, flow, CLI, app, and browser. | `reasoning.md#assumptions-and-risks` | Derive from canonical order where current design permits and add exact parity tests at each fixed projection. | `mitigated` |
| Existing schema-2 flow state has six stages and strict count validation. | Implementation seam review; Decision 1 | Add deterministic pre-stage interpretation before current validation; prove no read-time write and later canonical serialization. | `mitigated` |
| Artifact promotion succeeds but flow-state persistence fails. | `reasoning.md#assumptions-and-risks` | Define and test ordering, preserve recoverable diagnostics, and never infer completion from artifact presence. | `mitigated` |
| Last valid artifact masks a failed rerun. | Decisions 2 and 5 | Persist/project latest attempt separately from authoritative artifact validity across CLI/app/web. | `mitigated` |
| Current view model cannot express preserved artifact plus failed latest attempt. | `reasoning/frontend.md`; implementation seam review | Add narrow typed app/handler fields before template rendering; do not infer from strings or SSE. | `mitigated` |
| Runtime permission support cannot enforce required read/write isolation. | Decision 2 assumption | Fail preflight, retain old artifact/state, and report the capability gap; do not downgrade to best effort silently. | `mitigated by required capability and fail-closed policy` |
| Generated source paths or Markdown escape containment or disclose unsafe content. | Decisions 3 and 5 | Central path validation, allowlisted bounded preview, escaping, symlink checks, and hostile fixtures. | `mitigated` |
| Cancellation races with validation, promotion, or terminal publication. | Decision 6; shutdown contract | Reuse one context, canonical cancellation, single terminal arbitration, bounded reconciliation, and race tests. | `mitigated` |
| Candidate cleanup times out or process ownership is uncertain. | Decision 6 | Keep candidate contained, retain mutation ownership until reconciliation, and record `interrupted`/`cleanup_uncertain`. | `mitigated` |
| Stable JSON fields or nullability regress. | `reasoning/api-design.md`; compatibility docs | Preserve schema-version-1 envelopes and existing meanings; use additive DTO fields and compatibility tests. | `mitigated` |
| Structural validity is mistaken for semantic completeness. | Decision 3 | Use precise findings/help/UI language and preserve downstream live-source access. | `mitigated` |
| Real-provider behavior is not proven in Sprint 33. | Decision 7 | Keep all normal tests offline and record real-runtime dogfood as Sprint 34 deferred evidence. | `deferred` |

## Assumptions And Stop Conditions

| Assumption | Current Evidence | Implementation Check / Stop Condition |
| --- | --- | --- |
| Existing target resolution is reusable. | `execute_target.go` exposes `ResolveExecuteTarget` and existing execute/review/smoke callers. | Reuse it narrowly; stop before accepting raw caller paths or adding discovery/index machinery. |
| Existing filesystem conventions support safe replacement. | `state.go`, `review.go`, and `smoke.go` use same-directory temporary files and rename; state/review tests preserve prior bytes on injected failure. | Add stage-specific fault tests; if the guarantee is insufficient, add only the missing mechanical behavior, not a persistence abstraction. |
| Generic operations represent the stage. | `OperationStage` carries `OperationRequest.Stage`; `operation_runner.go` calls `FlowStage`; HTTP already maps `sprint-stage`; locks use sprint scope. | Reuse this path; stop before adding a route or operation kind. |
| Stable JSON permits additive fields. | Version-1 envelope and compatibility documents/tests allow additive projection while preserving types and nullability. | Extend DTOs additively and make compatibility tests pass before changing consumers. |
| Existing operation ownership handles disconnect and shutdown correctly. | `internal/web/operations.go` plus the shutdown contract already own cancellation, bounded events, terminal arbitration, and reconciliation. | Add code-context contract coverage; do not create stage-local goroutines or cancellation ownership. |

## Deferred Scope

- Sprint 34 owns exact requirements/code-context shared-prefix rendering and injection into downstream agent-backed stages.
- Sprint 34 owns manual code-context skill creation/materialization, broad documentation, and representative gated real-runtime dogfood.
- Later evidence-gated work owns content identity/provenance, QA/adjudication/repair, retrieval/embeddings, alternate persistence/SQLite, knowledge graphs, cloud, and Aren integration.
- No automatic source freshness, automatic rerun, amendment protocol, cache subsystem, cache key, provider cache-control dependency, repository index, parallel manifest, or hard excerpt-count limit is authorized.

## Review Inputs

Review should use:

- `requirements.md`
- `sprint-index.md`
- `technical-handbook.md`
- `reasoning/api-design.md`, `reasoning/architecture.md`, and `reasoning/frontend.md`
- `reasoning.md`
- this `plan.md`
- the four implementation plans listed in **Inputs Used**, reconciled to the Sprint 33 roadmap boundary
- implementation diff
- focused and repository-wide verification evidence
- `system/protocols/architecture-review-protocol.md`
- `system/protocols/review-sprint-protocol.md`

## Execution Log

| Date / Step | Action | Evidence / Notes |
| --- | --- | --- |
| `2026-08-19 / planning` | Materialized evidence-grounded Sprint 33 implementation plan. | Prerequisites validated; no implementation or verification task executed. |
| `2026-08-19 / implementation` | Implemented the code-context vertical slice manually in the target repository. | Added the canonical stage, compatibility projection, structural validator, fail-closed runtime service, atomic candidate promotion/rollback, source-aware configuration, CLI/app/web wiring, and embedded defaults. No UltraPlan CLI command was used. |
| `2026-08-19 / deterministic evidence` | Added focused domain, compatibility, validator, fake-runtime, temporary-repository isolation, configuration, command, shared-operation, API, and browser tests. | Tests cover invalid requirements, valid-artifact-without-success gating, legacy no-write reads, retained event output, permission capability, cancellation/interruption/cleanup uncertainty, state-persistence rollback, opaque preview, generic operation reuse, DOM order, and preserved-artifact/latest-attempt separation. |
| `2026-08-19 / verification` | Ran every command in **Verification Commands** from `../ultraplan-go`. | All focused suites, boundary suites, focused/full race suites, repository tests, vet, build, and `git diff --check` passed. |
| `2026-08-19 / Architecture Review` | Applied `architecture-review-protocol.md` to the final implementation diff. | **Approve.** Good architecture fit; complexity increase is justified; cohesion and coupling are acceptable; state changes are explicit and queries remain non-mutating for compatible legacy state; runtime stays generic; web imports only app; tests are strong; no required pre-merge change found. |
| `2026-08-19 / Sprint Review` | Applied `review-sprint-protocol.md` manually against requirements, decisions, plan, diff, and verification evidence. | **Pass.** `AC-01`–`AC-18` and `C-01`–`C-10` are covered with no blocker/high finding. Scope scan found only the prompt prohibition text for index/manifest/cache; no prohibited subsystem or dependency was added. Canonical generated `review.md` and live-provider dogfood were not produced; real-runtime evidence remains explicitly deferred to Sprint 34. |

## Completion Criteria

- [x] All seven tasks are complete or explicitly deferred with requirement impact recorded.
- [x] Every `AC-01` through `AC-18` and `C-01` through `C-10` row has implementation and evidence.
- [x] All required outputs in `requirements.md` exist and conform to the seven final decisions.
- [x] Verification commands were run or deferrals are documented truthfully.
- [x] Architecture Review and Sprint Review evidence confirms ownership, mutation boundaries, compatibility, generic interface reuse, and scope exclusions.
- [x] Evidence satisfies the expectations from `reasoning.md` without claiming unavailable real-runtime proof.
- [x] `review.md` can evaluate conformance without guessing intent.
