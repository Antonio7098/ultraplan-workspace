# Sprint Requirements: Code-Context Stage Vertical Slice

> Project: `ultraplan-go`
> Sprint: `33-code-context-stage`
> Purpose: the authoritative, human-readable sprint contract. All other sprint artifacts must satisfy these requirements.

## Sprint Goal

Make `code-context` a fully operational sprint stage immediately after requirements that can inspect the resolved implementation repository, atomically produce and validate the single authoritative `code-context.md` artifact, and expose truthful operation and recovery behavior through existing shared application and web boundaries.

## Required Outputs

| Output | Path | Description |
| --- | --- | --- |
| Stage domain model | `../ultraplan-go/internal/sprint/domain.go` | Defines `StageCodeContext` immediately after requirements and updates the canonical stage/state compatibility model. |
| Artifact mapping | `../ultraplan-go/internal/sprint/artifacts.go` | Maps `code-context` to the sprint-root `code-context.md` artifact. |
| Filesystem sprint snapshot support | `../ultraplan-go/internal/sprint/store_fs.go` | Includes `code-context.md` in sprint loading, status, and artifact snapshots without changing source-repository files. |
| Flow integration | `../ultraplan-go/internal/sprint/flow.go` | Implements readiness, prerequisites, cumulative dispatch, validation gating, and success/failure transitions for the new stage. |
| Flow-state compatibility | `../ultraplan-go/internal/sprint/state.go` | Keeps persisted pre-code-context sprint state usable through explicit, deterministic compatibility handling. |
| Code-context stage service | `../ultraplan-go/internal/sprint/code_context.go` | Implements prompt preview, target resolution, runtime execution, output isolation, atomic replacement, structural validation, cancellation, and recovery behavior. |
| Embedded code-context prompt | `../ultraplan-go/internal/workspace/scaffold/prompts/create-code-context.md` | Instructs requirements-driven, read-only repository exploration and permits only the sprint context artifact as output. |
| Embedded code-context template | `../ultraplan-go/internal/workspace/scaffold/templates/code-context.md` | Defines the editable Markdown structure for scope, inspected areas, selected source excerpts, relationships, constraints, and open questions. |
| Embedded-default registration | `../ultraplan-go/internal/workspace/init.go` | Registers the prompt and template with embedded defaults and `defaults install`. |
| Stage runtime configuration | `../ultraplan-go/internal/platform/config/config.go` | Adds code-context model and variant settings using the existing stage-specific fallback and source-tracking behavior. |
| Configuration redaction | `../ultraplan-go/internal/platform/config/redaction.go` | Applies existing safe configuration projection rules to the new stage settings. |
| Configuration command projection | `../ultraplan-go/internal/app/config_commands.go` | Reports effective code-context model and variant settings consistently in supported config output modes. |
| Sprint CLI and runtime wiring | `../ultraplan-go/internal/app/sprint_commands.go` | Supports code-context prompt, validation, flow, help, status, stable JSON, and stage-specific runtime selection. |
| Shared sprint use cases | `../ultraplan-go/internal/app/sprint_usecases.go` | Exposes code-context readiness, state, artifact, and findings through the shared application boundary. |
| Artifact preview allowlist | `../ultraplan-go/internal/app/usecases.go` | Allows bounded, contained preview of the sprint-owned `code-context.md`. |
| Shared operation dispatch | `../ultraplan-go/internal/app/operations.go` | Runs and fingerprints code-context operations through the existing progress, cancellation, conflict, and recovery model. |
| Sprint browser presentation | `../ultraplan-go/internal/web/templates/sprint.html` | Displays code-context state, findings, artifact access, and explicit rerun controls using existing generic operation semantics. |
| Stage service tests | `../ultraplan-go/internal/sprint/code_context_test.go` | Covers prompt inputs, validation, source-read/output-write isolation, fake-runtime outcomes, overrides, cancellation, and atomic rerun. |
| Sprint domain tests | `../ultraplan-go/internal/sprint/sprint_test.go` | Covers ordered-stage membership, artifact mapping, readiness, status, and stage count. |
| Flow-state compatibility tests | `../ultraplan-go/internal/sprint/verify_test.go` | Proves persisted state created before code-context remains readable and correctly reconciled. |
| Cumulative-flow tests | `../ultraplan-go/internal/sprint/sprint_index_test.go` | Proves code-context prerequisites and its exact position before sprint-index in cumulative flow. |
| Embedded-default tests | `../ultraplan-go/internal/workspace/workspace_test.go` | Proves embedded and materialized code-context prompt/template defaults remain synchronized. |
| Configuration tests | `../ultraplan-go/internal/platform/config/config_test.go` | Covers parsing, source tracking, fallback, validation, and redaction for code-context model/variant settings. |
| Sprint command tests | `../ultraplan-go/internal/app/sprint_commands_test.go` | Covers CLI acceptance, help, prompt preview, validation, status/JSON, dry-run non-mutation, and runtime override propagation. |
| Shared web use-case tests | `../ultraplan-go/internal/app/web_usecases_test.go` | Covers shared code-context stage, artifact, readiness, and finding projections. |
| Shared web operation tests | `../ultraplan-go/internal/app/web_operations_test.go` | Covers generic prepare/run/fingerprint/progress/cancellation/recovery behavior for code-context. |
| Browser template tests | `../ultraplan-go/internal/web/templates_test.go` | Covers code-context display, safe artifact access, findings, and explicit rerun presentation. |
| Web operation contract tests | `../ultraplan-go/internal/web/operations_contract_test.go` | Proves code-context uses the existing stage-operation contract without a route-specific operation kind. |

## Acceptance Criteria

- [ ] `PlanningStages()` and every canonical ordered stage list place `code-context` exactly once, immediately after `requirements` and before `sprint-index`.
- [ ] Valid completed requirements make `code-context` ready; missing, invalid, or failed requirements block it; a structurally valid completed `code-context.md` makes `sprint-index` ready.
- [ ] `prompt code-context`, `validate code-context`, and `flow --to code-context` work through the existing sprint CLI and typed application use cases, and help/status/stable JSON recognize the stage.
- [ ] Code-context prompt preview and dry-run are deterministic and do not create `code-context.md`, change `flow-state.json`, invoke the runtime, or mutate the implementation repository.
- [ ] Runtime execution resolves the target through existing project implementation-repository/worktree mechanisms, uses the generic runtime boundary with `context.Context`, and applies code-context model/variant overrides before existing fallbacks.
- [ ] The runtime can read the resolved implementation repository, but the stage permits no implementation-repository mutation and no stage output other than `projects/<project>/sprints/<sprint>/code-context.md`.
- [ ] A valid artifact contains all required sections, at least one selected exact source excerpt, repository-relative contained paths, a rationale per excerpt, language-tagged fenced source blocks, and well-formed optional line ranges and symbols.
- [ ] Validation rejects an absent or empty artifact, placeholders, missing required sections, missing excerpts, missing paths or rationale, missing fenced source, absolute or escaping source paths, and malformed line ranges with actionable findings.
- [ ] Runtime exit success without a present and valid `code-context.md` fails truthfully and does not mark the stage complete or make `sprint-index` ready.
- [ ] An explicit rerun atomically replaces only `code-context.md`; failed generation or validation preserves the last valid artifact and records a truthful failed/interrupted/cancelled outcome.
- [ ] Persisted flow state created before `code-context` existed loads without losing prior stage outcomes and gains deterministic compatibility behavior without hidden mutation during read-only status operations.
- [ ] The embedded prompt and template are available without workspace override files, and `ultraplan defaults install` materializes matching editable copies.
- [ ] Prompt/runtime metadata uses the existing stable identity, version/checksum, model, variant, and attempt fields where available; user-visible projections exclude full prompts, secrets, unsafe paths, raw provider payloads, and unbounded diagnostics.
- [ ] Browser readiness, progress, validation findings, artifact preview, explicit rerun, cancellation, and recovery use existing app capabilities and the generic operation model; no route-specific code-context workflow semantics or durable web-owned state is added.
- [ ] Cancellation, browser disconnect, and graceful server shutdown preserve existing ownership rules: disconnect does not cancel work, explicit or shutdown cancellation reaches the canonical operation context, and uncertainty is never presented as success.
- [ ] Focused tests cover stage order, prerequisites, compatibility, prompt/default assets, structural validation, path containment, dry-run non-mutation, fake-runtime success and failure, missing/invalid output, atomic rerun, configuration fallback, CLI/app/web projections, cancellation, and recovery.
- [ ] `go test ./...`, `go test -race ./...`, `go vet ./...`, `go build ./cmd/ultraplan`, and `git diff --check` pass in `../ultraplan-go`.
- [ ] No repository index, RAG/embedding system, cache subsystem or cache key, provider-specific cache-control dependency, parallel JSON context manifest, automatic staleness detector, context-amendment workflow, or source mutation is introduced.

## Non-Goals

- Injecting exact requirements and `code-context.md` bytes into sprint-index, handbook, reasoning, plan, execute, review, or smoke prompts; the shared downstream prefix is Sprint 34 scope.
- Adding or materializing a manual `code-context` stage skill; that is Sprint 34 scope.
- Completing broad README, CLI reference, user guide, recovery, architecture, planning-smoke, generated-workspace, skill, and local-web documentation updates; Sprint 33 changes only executable help and defaults needed for its vertical slice.
- Real-runtime requirements-to-plan dogfood and the Phase 5 release gate; deterministic fake-runtime coverage is required here and gated real-runtime proof belongs to Sprint 34.
- Automatic source-change detection, automatic context refresh, formal amendment requests, arbitrary maximum excerpt counts, or guarantees of provider prompt-cache hits.
- Content identity/provenance, expanded QA or repair, retrieval/search, embeddings, SQLite or alternate product persistence, knowledge graphs, cloud operation, or Aren integration.
- Reusing code-context generation for study workflows or restricting downstream agents to only the selected context pack.
- Creating a generic stage framework, plugin registry, workflow engine, repository-index package, web-specific product service, or package-global mutable runner registration.

## Constraints

- `internal/sprint` owns code-context order, prerequisites, artifact semantics, validation, runtime coordination, flow transitions, and compatibility; `internal/platform/runtime` remains generic and imports no product semantics.
- CLI and HTTP adapters must call typed shared application use cases. `internal/web` must not import product modules, parse CLI output, invoke the UltraPlan binary, inspect repositories, duplicate validation, or persist workflow truth.
- Filesystem artifacts and sprint-owned state remain authoritative. Web operation/SSE state is bounded and ephemeral, and read-only status must not conceal migration writes.
- Repository paths must be contained and repository-relative in the artifact. The implementation target must be resolved through existing project/execute mechanisms rather than a new discovery or indexing subsystem.
- The stage may read the implementation repository and may write only the sprint-root `code-context.md`; source, tests, Git state, governed inputs, and unrelated sprint artifacts are immutable during stage execution.
- Artifact replacement and flow-state writes must use existing atomic persistence guarantees. Stage completion occurs only after output existence and structural validation succeed.
- Validation is structural and must not claim semantic completeness; the pack is prepared evidence, not a repository dump, final architecture decision, plan, or exclusive source boundary.
- Runtime-backed behavior must propagate `context.Context`, cancellation, progress, typed errors, stable error codes, redaction, and existing runtime metadata without direct OpenCode/process handling in product code.
- Normal tests must be offline and deterministic with fake runtimes and temporary repositories; unavailable gated runtime credentials or environments must be reported as blocked or deferred, never passed.
- The workspace roadmap owns Sprint 33 sequencing and scope. Implementation plans are mandatory design inputs but cannot pull Sprint 34 or later gated work into this sprint.

## Dependencies

| Prior Sprint / Output | Required For | Notes |
| --- | --- | --- |
| Sprint 32, `projects/ultraplan-go/sprints/32-hardening-and-release/review.md` | Shared interface and web extensibility baseline | Current verdict is `pass_with_findings`; use the proven generic app capability and web operation model, preserve redaction and compatibility fixtures, and do not add route-specific stage logic. |
| Project implementation target, `../ultraplan-go/` | Repository inspection and implementation | Resolved through the project index and existing execute/worktree mechanisms; source access is read-only for code-context execution. |
| `../ultraplan-go/docs/plans/integrated-roadmap.md` | Sequence and shared-boundary constraints | Used as a mandatory design input; it confirms web-first delivery, filesystem authority, and the code-context vertical slice after web hardening. |
| `../ultraplan-go/docs/plans/sprint-code-context-stage.md` | Detailed stage, artifact, validator, runtime, and test design | Used as the primary detailed implementation input, reconciled by deferring shared-prefix integration, the manual skill, broad docs, and release dogfood to Sprint 34 per the workspace roadmap. |
| `../ultraplan-go/docs/plans/ultraplan-local-server-experiment-plan.md` | Filesystem/app/web ownership boundary | Applicable constraint: shared app use cases remain central and neither generic persistence nor a database enters this sprint. |
| `../ultraplan-go/docs/plans/server-shutdown-run-cancellation-contract.md` | Browser operation cancellation and recovery | Applicable to code-context operations started by the server; shutdown cancellation and browser-disconnect semantics remain unchanged. |
| Existing sprint runtime, flow state, defaults, app operation, and web operation implementations | Vertical-slice integration | Reuse existing stage/runtime/configuration/persistence/operation conventions rather than introducing parallel abstractions. |
| Prior sprint review decisions | Error, state, prompt, security, and test quality | Preserve typed errors, strict atomic state, prompt identity/checksum metadata where available, projection-time allowlists/redaction, fake-first tests, and direct stage/command coverage. |

## Review Expectations

| What | How Verified |
| --- | --- |
| Roadmap scope and plan reconciliation | Compare implementation and sprint artifacts with Sprint 33 in `projects/ultraplan-go/roadmap.md` and the four named implementation plans; reject Sprint 34/later scope creep. |
| Stage order and readiness | Unit tests and status fixtures prove `requirements -> code-context -> sprint-index`, exact-once cumulative dispatch, and truthful prerequisite blocking. |
| Artifact quality and validator behavior | Table-driven validator tests plus manual inspection of a representative fixture verify required sections, exact fenced excerpts, paths, rationale, ranges, relationships, and actionable failures. |
| Repository-read/output-write isolation | Fake-runtime integration tests use a temporary implementation repository, compare source/Git state before and after, and prove only the sprint `code-context.md` may change. |
| Runtime and failure truthfulness | Fake-runtime tests cover success, runtime failure, missing output, malformed output, cancellation, and failed atomic rerun; flow state must never infer success from exit status or artifact presence alone. |
| Compatibility and persistence | Legacy flow-state fixtures load deterministically; atomic-write failure tests preserve the last valid artifact/state; read-only status produces no hidden migration write. |
| CLI and JSON behavior | Command tests exercise help, prompt, validate, flow, status, dry-run, model/variant fallback, stable JSON fields, typed error codes, and redacted diagnostics. |
| Shared app/web extensibility | App and web contract tests prove existing generic capabilities expose prepare/start/progress/cancel/result/recovery and artifact preview without a new route, operation kind, product import, or web-owned durable state. |
| Cancellation ownership | Operation and shutdown tests prove browser disconnect is subscription loss only, explicit cancellation propagates, shutdown uses the canonical path, and interrupted or cleanup-uncertain work is not shown as complete. |
| Defaults and configuration | Embedded/default-install golden tests and config tests verify prompt/template parity, stage-specific model/variant parsing, fallback, source tracking, and redaction. |
| Regression and release evidence | Review `go test ./...`, `go test -race ./...`, `go vet ./...`, `go build ./cmd/ultraplan`, and `git diff --check` logs; real-runtime evidence is accepted only if actually run and otherwise recorded as gated/deferred. |
| Explicit exclusions | Diff and dependency review confirm no index/RAG/cache/JSON-manifest/staleness/amendment/persistence/QA/repair/graph/cloud work and no source or Git mutation. |
