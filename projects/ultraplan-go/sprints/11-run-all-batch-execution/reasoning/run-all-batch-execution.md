> **Inputs Used:** `projects/ultraplan-go/sprints/11-run-all-batch-execution/sprint-index.md`, `projects/ultraplan-go/sprints/11-run-all-batch-execution/technical-handbook.md`, `projects/ultraplan-go/sprints/11-run-all-batch-execution/requirements.md`, `projects/ultraplan-go/docs/ARCHITECTURE.md`, `projects/ultraplan-go/docs/TRD.md`, `projects/ultraplan-go/docs/PRD.md`, `system/reasoning/architecture_reasoning_template.md`

# Architecture Reasoning: Run All Batch Execution

## 0. Feature Summary

### Feature name

Run-all batch execution for studies.

### User/product goal

`ultraplan study <study> run-all` lets a user execute all selected applicable analysis tasks for a study with bounded parallelism, validate every produced artifact, synthesize dimensions that have complete valid inputs, and regenerate a deterministic `summary.csv` without implementing durable `run-loop` orchestration.

### Current task

Create the architecture decision basis for implementing study-owned batch execution in `internal/study`, thin CLI wiring in `internal/app`, and fake-runtime tests for task matrix creation, filtering, bounded concurrency, validation gates, synthesis gating, summary semantics, cancellation, and command output.

### Non-goals

- Do not implement `ultraplan study <study> run-loop`.
- Do not implement durable retry scheduling, stale-running recovery, lock files, force unlock, or multi-process coordination.
- Do not add new runtime adapters, direct OpenCode parsing, or product imports into `internal/platform/runtime`.
- Do not introduce global technical-layer packages such as `internal/scheduler`, `internal/reports`, `internal/prompts`, `internal/summary`, or `internal/validation`.
- Do not create stable public JSON output for `run-all` unless existing app conventions already require it.
- Do not add default real OpenCode smoke tests to normal verification.

## 1. First-Principles Breakdown

### Core behaviour

This feature receives a study name, optional source and dimension filters, and resolved execution configuration. It resolves a deterministic set of applicable analysis tasks, executes them through the existing runtime boundary with a bounded worker pool, validates expected reports, runs synthesis for dimensions whose applicable selected inputs are valid, regenerates `summary.csv`, and returns an in-memory batch result that distinguishes completed, failed, skipped, and pending/not-run work.

### Inputs

- Study name and workspace root.
- Discovered dimensions, directory sources, and Markdown document sources.
- Optional source and dimension filters using existing resolution and ambiguity rules.
- Resolved runtime, model, timeout, permission, health, and parallelism configuration.
- Existing prompt builders, report validators, rating parser, source applicability helper, and platform runtime abstraction.
- `context.Context` for user interrupt, timeout, and cancellation propagation.

### Outputs

- Per-source analysis reports produced by runtime tasks and accepted only after validation.
- Final synthesis reports for dimensions with complete valid applicable selected inputs.
- `studies/<study>/summary.csv` with deterministic ordering, ratings, totals, missing warnings, and inapplicable Markdown semantics.
- A `RunAllResult`-style in-memory result with task statuses, safe diagnostics, artifact paths, warnings, and original errors retained for internal inspection.
- Deterministic human CLI output rendered by `internal/app`.

### Durable state

- Created or updated: validated per-source reports, final synthesis reports, and `studies/<study>/summary.csv`.
- Read: study definition, dimensions, sources, prompt templates, report templates, existing reports for validation and synthesis gating, workspace config.
- Updated: no durable run-loop state for Sprint 11.
- Deleted: none.

Decision: `run-all` remains an ephemeral bounded batch command for this sprint. It does not write `run-state.json`, locks, retry schedules, or resumability metadata because requirements explicitly exclude Sprint 12 durable orchestration and allow state primitive updates only if needed for truthful status.

### Ephemeral state

- Resolved request and effective execution config.
- Static ordered analysis task matrix.
- Per-task runtime handles, event drainers, validation results, safe diagnostics, and cause errors.
- Synthesis eligibility map by dimension.
- Summary warnings and rating parse results.
- Cancellation state for active and pending tasks.

### Derived state

- Applicable source/dimension pairs are derived from discovered sources, dimensions, filters, and `GetApplicableSources` semantics; they are computed at command start and not stored durably.
- Completed analysis status is derived from runtime outcome plus report validation; runtime success alone is not a source of truth.
- Synthesis eligibility is derived after the analysis phase from valid applicable selected per-source reports.
- Summary totals are derived from parsed ratings and regenerated on every `run-all` invocation.

### Side effects

- File writes: per-source reports and final reports via runtime tasks, plus UltraPlan-owned `summary.csv` through explicit workspace-contained paths.
- Runtime calls: analysis and synthesis tasks through `internal/platform/runtime` backed by agentwrap/OpenCode.
- Event consumption: every started runtime run has its event channel drained or safely consumed until completion.
- Logs/diagnostics: safe task status, validation failures, warnings, artifact paths, and redacted runtime diagnostics.
- Network/subprocess: only through agentwrap/runtime behavior; product code does not invoke OpenCode directly.

## 2. Existing Architecture Fit

### Contract applicability and phase gate

Current gate: study-side batch workflow with runtime/provider execution, but without durable run-loop orchestration.

Applies now:

- Architecture: `internal/study` owns batch behavior; `internal/app` stays thin; `internal/platform/runtime` stays generic.
- CLI Surface: command shape, flags, help, validation, exit statuses, and deterministic text output are in scope.
- Configuration: runtime/model/timeout/permission/health/parallelism must be resolved before scheduling.
- LLM Runtime: tasks use the existing agentwrap-backed platform runtime boundary.
- Errors: usage, validation, runtime, cancellation, and partial-completion behavior must preserve causes internally and render safe diagnostics.
- Observability: runtime events must be consumed and user output must report truthful task counts and artifact paths.
- Security: prompts, Markdown bodies, secrets, sensitive environment values, and unsafe native payloads must not be printed.
- Testing: fake-runtime and local-fixture verification are required; normal tests must not require OpenCode or network access.
- Performance: concurrency must be bounded and summary/report processing deterministic.
- Persistence: `summary.csv` and any generated artifacts must use explicit workspace-contained paths.

Deferred:

- Durable workflow state, lock files, stale-running recovery, and resumability remain `run-loop` scope.
- Public JSON schema stability remains a later release concern unless current app conventions already require it.
- New runtime adapters, generic workflow engines, DAG authoring, plugins, target commands, sprint commands, and code extraction are out of this sprint.

### Where does this behaviour naturally belong?

The behavior belongs in `internal/study` because the batch task matrix, source applicability, report paths, validation, synthesis gating, and summary generation all transform study state. `internal/app` should parse arguments, resolve command-level options, call a study service method, map result status to exit behavior, and render deterministic text. `internal/platform/runtime` should only receive generic runtime requests with prompt, workdir, model, timeout, permissions, expected output path, metadata, events, and results.

Evidence basis:

- The technical handbook states that mature Go CLIs keep entrypoints thin and command code delegates to domain/action logic, with `run-all` pressure specifically pointing to `internal/study` ownership and `internal/app` rendering only (`technical-handbook.md:31-33`, `technical-handbook.md:60-61`).
- The architecture document assigns study scheduling, reports, validation, prompts, state, and summary generation to `internal/study`, while runtime must not know studies, dimensions, synthesis gating, report semantics, or summary generation (`ARCHITECTURE.md:23-32`, `ARCHITECTURE.md:116-157`, `ARCHITECTURE.md:159-197`).
- The sprint requirements make this a hard constraint and require architecture review to confirm no global scheduler/report/summary/prompt/validation package appears (`requirements.md:70-80`, `requirements.md:107-108`).

### Existing workflow affected

Current flow before this sprint:

1. A single analysis command resolves one dimension and one source.
2. It composes the correct prompt for directory or Markdown source behavior.
3. It invokes the generic runtime boundary.
4. It validates the expected per-source report.
5. Synthesis and summary behavior exist as separate study-side capabilities from prior sprints.

### Proposed new flow

Proposed `run-all` flow:

1. `internal/app` validates command shape and resolves flags/config, including `parallelism >= 1`.
2. `internal/study.Service.RunAll(ctx, req)` loads the study, sources, and dimensions.
3. Existing filter resolution handles source/dimension names, numbers, slugs, filenames, prefixes, unknown references, and ambiguity.
4. `internal/study` builds a static deterministic analysis task matrix from selected dimensions and selected applicable sources only.
5. A bounded worker pool runs analysis tasks through the existing single-analysis execution path or shared study-owned helper.
6. Each started runtime run drains events and waits for final runtime result.
7. Each completed runtime run is validated against the expected per-source report; validation failure marks the task failed.
8. Failed analysis tasks do not cancel unrelated tasks; cancellation from context stops scheduling new work and lets active runtime calls return through the runtime boundary.
9. After the analysis phase, `internal/study` computes synthesis eligibility for each selected dimension.
10. Eligible synthesis tasks run through the existing synthesis prompt/runtime/final-report validation path using the same bounded runtime concurrency cap.
11. Dimensions with failed or missing applicable selected analysis reports skip synthesis and are reported as skipped, not successful.
12. Summary generation runs after analysis and eligible synthesis work finishes, writes `studies/<study>/summary.csv`, and returns warnings for missing or ambiguous ratings.
13. `internal/app` renders deterministic text with counts, skipped/pending status, warnings, safe diagnostics, and artifact paths.

### Does the current architecture still fit?

Yes. The feature fits cleanly into the module-driven shape if batch coordination remains inside `internal/study` and runtime execution remains a generic injected dependency.

### Refactor-before-feature decision

Decision: implement directly with small local study-module helpers if needed.

Reason: the sprint depends on existing single analysis, synthesis, prompt, validation, source applicability, and runtime boundary semantics. The correct architecture is composition of those study-owned capabilities, not a new scheduler or workflow layer.

## 3. Design Options Considered

### Option A: Study-owned batch orchestration with local worker pool

Description:

Add `RunAll` request/result types and orchestration in `internal/study`, reuse existing analysis and synthesis paths, build an upfront deterministic matrix, execute analysis through a bounded worker pool, run eligible synthesis after analysis completion, then regenerate summary.

Pros:

- Preserves module ownership and dependency direction.
- Reuses existing prompt/runtime/validation behavior instead of duplicating product rules.
- Makes fake-runtime concurrency, failure, cancellation, and output tests straightforward.
- Keeps `internal/app` thin and deterministic.
- Avoids durable `run-loop` scope.

Cons:

- Requires explicit result aggregation instead of relying on a single returned error.
- Requires careful event draining for every started runtime task.
- Synthesis starts later than the earliest possible moment.

Risks:

- If the result type grows too broad, it could start resembling durable run-loop state. Keep it in-memory and command-specific.

### Option B: Refactor into a reusable global scheduler first

Description:

Introduce a cross-module scheduler package that models task DAGs, workers, statuses, synthesis dependencies, and summaries, then use it from the study module.

Pros:

- Could eventually serve future run-loop or target/sprint workflows.
- Centralizes concurrency mechanics.

Cons:

- Violates sprint and architecture pressure against global technical-layer packages.
- Pulls study-specific concepts into generic infrastructure or forces awkward callbacks.
- Increases abstraction before a concrete cross-module reuse need exists.
- Risks implementing a generic workflow engine, explicitly excluded by the sprint.

Risks:

- Creates coupling and future maintenance burden without solving a current requirement.

### Option C: Durable run-loop state as the implementation substrate

Description:

Implement `run-all` by creating or updating durable run-state records, then execute the state machine once without exposing a `run-loop` command.

Pros:

- Aligns with long-term product resumability direction.
- Could make status commands more informative.

Cons:

- Conflicts with the sprint requirement to avoid durable `run-loop` orchestration.
- Implies migrations, locks, stale-running recovery, resume validation, and retry scheduling.
- Makes partial completion and cancellation harder to distinguish from resumable state.

Risks:

- Users may infer production-grade resumability before Sprint 12 implements it.

### Option D: Fail-fast `errgroup.WithContext` fan-out

Description:

Run analysis tasks in an `errgroup.WithContext`; the first task failure cancels the group and pending work.

Pros:

- Simple cancellation model.
- Common Go concurrency idiom.

Cons:

- Conflicts with partial-completion requirements by hiding how many unrelated tasks could have completed.
- Produces avoidable pending/not-run tasks on ordinary runtime or validation failure.
- Makes per-task diagnostics and deterministic summaries less useful.

Risks:

- A single flaky source/dimension pair prevents useful progress on the rest of the study.

### Chosen option

Chosen option: Option A.

Reason: it is the simplest honest design for Sprint 11. The handbook evidence favors thin command delegation, study-owned behavior, bounded concurrency, fakeable runtime seams, and explicit partial-completion aggregation (`technical-handbook.md:31-45`, `technical-handbook.md:47-56`, `technical-handbook.md:93-111`). Options B and C are rejected because they introduce future durable orchestration or generic workflow abstractions before this sprint needs them. Option D is rejected because ordinary task failure should not cancel unrelated work.

## 4. Abstraction Check

### Are we adding a new abstraction?

Yes, but only study-domain data structures and service surface:

- `RunAllRequest` or equivalent: captures study name, selected filters, config overrides, and output options.
- `RunAllResult` or equivalent: captures in-memory counts, task results, synthesis outcomes, warnings, artifact paths, and overall status.
- `RunAllTaskResult` or equivalent: captures per-task status, safe diagnostic, output path, validation facts, and original error for internal inspection.

No new runtime interface, global scheduler interface, plugin registry, or workflow engine is earned.

### Why is it earned?

- It protects a stable domain concept: a study batch run has task-level partial completion and summary warnings that a single error cannot represent.
- It makes testing simpler by exposing deterministic task outcomes without parsing CLI text.
- It creates a clear boundary between study policy and CLI rendering.
- It avoids leaking study semantics into `internal/platform/runtime`.

### Why not add more abstraction?

The architecture document advises interfaces only at external or volatile boundaries (`ARCHITECTURE.md:304-318`). Runtime execution is already that boundary. Matrix building, synthesis gating, summary generation, and result aggregation are concrete study behaviors that should stay in the study module until real cross-module reuse exists.

### Bad abstraction smell check

- Avoid names like `Scheduler`, `Manager`, `Processor`, or `Workflow` unless the code owns a concrete study behavior and has a study-specific name.
- Avoid boolean modes that make one helper behave like single-run, run-all, and run-loop simultaneously.
- Avoid mirroring runtime interfaces just to make tests easy; use the existing platform runtime seam or narrow test fakes already accepted by the project.

## 5. DRY and Duplication Check

### Is any duplication being introduced?

No duplicated business knowledge should be introduced.

Rules that must remain single-source:

- Source and dimension resolution, including ambiguity behavior.
- Markdown applicability through `GetApplicableSources` or equivalent helper.
- Directory vs Markdown prompt composition.
- Runtime request mapping, workdir rules, permission policy, health requirements, and timeout/model mapping.
- Per-source and final-report validation.
- Rating parsing and summary ordering.

Acceptable duplication:

- Small loop structure for analysis and synthesis task execution is acceptable if the tasks have different inputs and validations.
- CLI text formatting can be command-specific because command output changes for user-facing reasons.

Decision: if existing single-run logic is embedded too deeply in command code, move it into `internal/study` before `run-all` calls it. Do not duplicate prompt or validation logic in batch-specific code.

## 6. Coupling Check

### Dependencies required

- `internal/study` depends on existing study store/filesystem collaborators to load studies, sources, dimensions, reports, templates, and write summary artifacts.
- `internal/study` depends on the generic platform runtime abstraction for analysis and synthesis execution.
- `internal/study` may depend on workspace path helpers for workspace-contained path resolution.
- `internal/app` depends on `internal/study` service and result types for command wiring.
- `internal/platform/runtime` has no dependency on `internal/study`.

### Coupling risks

Global coupling: none introduced if worker state is per request and no package-level semaphores or mutable globals are added.

Content coupling: none introduced if runtime request metadata contains safe generic fields and runtime does not inspect study concepts.

Stamp coupling: acceptable at service boundaries where the full request/result concept is needed; internal helpers should receive narrower task inputs when practical.

### Dependency injection check

Side-effectful dependencies should be created at the app composition root or existing service construction point and passed inward. The handbook cites manual composition roots and narrow interfaces as the mature Go CLI pattern (`technical-handbook.md:16`, `technical-handbook.md:33-34`, `technical-handbook.md:41-45`).

## 7. State and Mutation Check

### Mutation points

- Runtime task execution writes expected per-source and final report files.
- `internal/study` validates those files before marking tasks completed.
- `internal/study` writes `studies/<study>/summary.csv` after analysis and eligible synthesis finish.
- `internal/app` writes only stdout/stderr text, not study state.

### Is mutation explicit from names and flow?

Yes. `RunAll` is a command-style operation with clear side effects. Matrix building, eligibility checks, and summary calculation should be separated enough that tests can exercise pure decision logic without runtime side effects.

### Could hidden state make this hard to debug?

Hidden state risk is controlled by keeping all worker queues, result maps, counters, and cancellation flags scoped to one `RunAll` invocation. No package-level semaphore, global result cache, or background worker is acceptable.

## 8. Function/Class/Module Shape

### Primary unit of behaviour

Primary unit: study service workflow in `internal/study`.

### Why this unit?

The service already represents use-case coordination and dependency wiring for study behavior. `run-all` is cohesive with study lifecycle, source applicability, prompts, validation, synthesis, reports, and summary generation. A method such as `Service.RunAll(ctx, req)` keeps callers from knowing internal helpers exist, matching the architecture document's service-surface guidance (`ARCHITECTURE.md:272-303`).

### Composition

- Study loader/store: reads dimensions, sources, templates, reports, and writes summary.
- Runtime: executes generic requests and returns runs/events/results.
- Existing prompt builders: build analysis and synthesis prompts.
- Existing validators: validate per-source and final artifacts.
- Existing rating parser: parses summary ratings and reports warnings.
- Local worker-pool coordination: bounds concurrent runtime tasks and aggregates per-task results.

## 9. Error Handling Design

### Expected failures

- Usage errors: missing study argument or invalid flags; handled in `internal/app` and mapped to usage exit behavior.
- Unknown or ambiguous source/dimension filters: fail before launching any runtime task with actionable diagnostics.
- Invalid parallelism less than 1: fail before launching any runtime task.
- Runtime failure: mark task failed, preserve original error internally, render safe diagnostic, continue unrelated tasks unless context is cancelled.
- Missing or invalid per-source report after runtime success: mark task failed through validation failure; do not treat runtime success as completion.
- Synthesis blocked by failed or missing analysis report: mark synthesis skipped for that dimension with the blocking facts.
- Missing or ambiguous ratings during summary generation: emit warnings and deterministic empty cells where appropriate; do not invent values.
- Cancellation: stop scheduling new tasks, let active runtime calls return through the runtime boundary, mark never-started tasks pending/not-run, and return cancellation or partial-cancellation status according to existing app exit conventions.

### Unexpected failures

Unexpected filesystem, runtime, or validation errors should be wrapped with cause chains and converted to safe user diagnostics at the CLI boundary. Error strings should not be parsed for classification; runtime classifications should use agentwrap/platform error categories where available, matching the TRD's `errors.As` requirement (`TRD.md:1193-1244`).

### Error taxonomy

Reuse existing app and product error conventions where present:

- Usage or argument error.
- Config error.
- Workspace/filesystem error.
- Validation failure.
- Runtime failure.
- Cancellation.
- Partial completion.

Decision: `RunAllResult` should compute an overall status from task outcomes instead of collapsing outcomes into one error. A non-nil error may still be returned for preflight failures that prevent any batch from starting.

### Retry / recovery behaviour

- Product-owned durable retry is not implemented in this sprint.
- Existing agentwrap policy retry/fallback may still occur inside individual runtime calls according to resolved configuration.
- Partial progress is preserved in generated artifacts and in the returned in-memory result, but not as durable resumable state.
- No compensation is needed beyond truthful reporting and validation; failed tasks leave their artifacts for inspection if any were written.

## 10. Observability Design

### Events/logs/metrics/traces needed

- Task queued: when the deterministic matrix is built.
- Task started: before runtime request launch.
- Runtime event received: consumed through the platform runtime event stream, with unsafe payloads omitted from user output.
- Validation completed or failed: after artifact validation.
- Task completed, failed, skipped, pending, or cancelled: after each terminal task decision.
- Synthesis skipped: when a dimension lacks complete valid applicable selected inputs.
- Summary generated: after `summary.csv` is written, including warning count and artifact path.

### Minimum useful fields

- Study name.
- Task kind: analysis or synthesis.
- Dimension number and slug.
- Source name and source kind for analysis tasks.
- Output path.
- Runtime/model identifiers when safe.
- Task duration and outcome.
- Validation status and failed check names when safe.
- Safe runtime error category and user detail.

### Sensitive data check

This feature can expose prompts, embedded Markdown document bodies, source content, environment values, and unsafe native runtime payloads if diagnostics are careless. Mitigation: render only safe summaries, paths, status counts, validation check names, redacted runtime fields, and explicit artifact paths. Do not print prompt text, Markdown source bodies, full environment, secrets, unsafe raw payloads, or direct OpenCode stdout/stderr. The handbook identifies this as a sprint-specific anti-pattern (`technical-handbook.md:43`, `technical-handbook.md:68`).

## 11. Testing Strategy

### Unit tests

- Deterministic task matrix ordering across dimensions and sources.
- Filter resolution preflight for valid, unknown, and ambiguous source/dimension refs.
- Markdown applicability exclusion for matrix creation, synthesis gating, summary missing warnings, and totals.
- Overall result status calculation for success, runtime failure, validation failure, partial completion, skipped synthesis, and cancellation.
- Summary CSV ordering, rating parsing, missing cells, inapplicable sentinel cells, ambiguous rating warnings, and stable totals.

### Use-case/workflow tests

- Successful `run-all` with fake runtime writing valid per-source and final reports.
- Fake-runtime bounded parallelism using blocked tasks to assert concurrent started analysis tasks never exceed configured parallelism.
- Runtime success without expected artifact fails validation.
- Failed analysis blocks synthesis only for its dimension while unrelated dimensions continue.
- Cancellation stops new scheduling, drains active runtime events, waits for active runs, and reports pending/not-run tasks.
- Event streams are consumed for every started runtime run.

### Command-level tests

- Help includes `ultraplan study <study> run-all` and documented flags.
- Missing arguments and invalid flags map to existing usage/config exit behavior.
- Parallelism less than 1 fails before runtime launch.
- Output is deterministic and includes completed, failed, skipped, pending/not-run, warning counts, and artifact paths.
- Output does not include prompts, embedded Markdown bodies, secrets, sensitive environment values, or unsafe native payloads.

### Integration tests

Default verification should not use real OpenCode, provider credentials, network access, or real subprocess execution. Gated OpenCode smoke tests may exist separately only if explicitly opt-in. This follows the sprint requirements and handbook testing evidence (`requirements.md:53-55`, `technical-handbook.md:24`, `technical-handbook.md:45`, `technical-handbook.md:74`).

### Testability check

The core behavior can be tested without real infrastructure because runtime execution is behind the platform runtime seam and file outputs can be created by fake runs in local fixtures.

## 12. Performance and Scale Assumptions

### Known constraints

- Expected volume: studies may contain many dimensions and sources, but the task matrix is bounded by selected dimensions times applicable selected sources.
- Latency sensitivity: batch runtime calls dominate; command preflight and summary generation should remain deterministic and lightweight.
- Memory sensitivity: do not recursively read source repositories or accumulate runtime payloads in memory.
- Concurrency concerns: runtime calls must be bounded, event channels must be consumed, and no unbounded goroutine-per-task launch pattern is acceptable.

### Assumptions being made

- Static task matrices are acceptable for Sprint 11 because they improve determinism and testability and the selected study sizes are local CLI scale.
- Synthesis can run after the analysis phase rather than as soon as each dimension becomes unblocked; this is simpler and deterministic.
- Summary generation can read report files intentionally but should not scan source repositories.

### Measurement plan

- Unit and fake-runtime tests assert bounded concurrency behavior.
- No benchmark is required for this sprint unless implementation reveals unexpectedly large memory use.
- Review should inspect for unbounded channels, unbounded goroutine creation, and accidental source tree scans.

### Optimization decision

Do not optimize beyond bounded workers, deterministic ordering, event consumption, and avoiding unnecessary large reads. The handbook performance evidence favors deferring expensive initialization, bounding goroutines, and avoiding unbounded accumulation (`technical-handbook.md:26`, `technical-handbook.md:62`, `technical-handbook.md:107`).

## 13. Implementation Plan

### Files likely to change

- `/home/antonioborgerees/coding/ultraplan-go/internal/study/domain.go`: add minimal request/result/task status types for batch execution if existing types are insufficient.
- `/home/antonioborgerees/coding/ultraplan-go/internal/study/service.go`: expose the `RunAll` service method and wire runtime dependency if not already available.
- `/home/antonioborgerees/coding/ultraplan-go/internal/study/run_all.go`: implement matrix creation, bounded workers, event consumption, validation gates, synthesis gating, cancellation, and result aggregation.
- `/home/antonioborgerees/coding/ultraplan-go/internal/study/summary.go`: implement or update deterministic summary generation and warnings.
- `/home/antonioborgerees/coding/ultraplan-go/internal/app/study_commands.go`: add thin command wiring, flags, validation, exit mapping, and deterministic text rendering.
- `/home/antonioborgerees/coding/ultraplan-go/internal/study/run_all_test.go`: add fake-runtime batch tests.
- `/home/antonioborgerees/coding/ultraplan-go/internal/study/summary_test.go`: add deterministic summary tests.
- `/home/antonioborgerees/coding/ultraplan-go/internal/app/study_run_all_commands_test.go`: add command-level tests.

### Step-by-step plan

1. Confirm existing single analysis, synthesis, validation, rating parser, source applicability, and runtime request APIs are callable from `internal/study` without command-layer dependencies.
2. Add minimal `RunAll` request/result/status types in the study module.
3. Implement preflight: load study, resolve filters, validate parallelism, resolve config, and fail before runtime launch for invalid inputs.
4. Build a deterministic analysis matrix from selected dimensions and applicable selected sources only; record inapplicable Markdown pairs as skipped only where useful for output/summary semantics.
5. Implement bounded analysis workers that drain runtime events, wait for every started run, validate expected reports, and record per-task outcomes.
6. Preserve partial completion by continuing unrelated tasks after ordinary runtime or validation failures; stop scheduling only on context cancellation or preflight failure.
7. Compute synthesis eligibility after the analysis phase and run eligible synthesis tasks through existing synthesis behavior with bounded runtime concurrency.
8. Mark synthesis skipped for dimensions with failed or missing applicable selected reports.
9. Generate deterministic `summary.csv`, using a documented sentinel such as `N/A` for inapplicable Markdown cells and empty cells for missing ratings, with warnings for missing or ambiguous ratings.
10. Add thin CLI command wiring and deterministic text output.
11. Add fake-runtime, summary, and command tests.
12. Verify with `go test ./...` and `go build ./cmd/ultraplan` from `/home/antonioborgerees/coding/ultraplan-go` during implementation review.

### Migration/backwards compatibility

- Schema migration needed: no.
- API compatibility affected: no public API commitment is created by these internal types.
- Existing data affected: existing reports may be read and validated; `summary.csv` may be created or overwritten deterministically.
- Feature flag needed: no.

### Rollback plan

Revert the new `run-all` command wiring and study `RunAll` files/types. Generated per-source reports, final reports, and `summary.csv` are workspace artifacts and can be inspected or regenerated; no durable orchestration state migration is introduced.

## 14. Final Pre-Implementation Decision

Decision: Proceed with study-owned `run-all` batch orchestration using a deterministic upfront matrix, local bounded worker pool, runtime/event/validation reuse, post-analysis synthesis gating, deterministic summary generation, and thin CLI rendering.

Reason: This design satisfies the sprint goal while preserving the project architecture: product behavior stays in `internal/study`, runtime remains generic, CLI stays thin, and durable run-loop concerns remain deferred. The technical handbook evidence supports thin command delegation, manual dependency wiring, bounded concurrency, validation as product gate, injected IO, safe diagnostics, fake-runtime tests, and avoiding global workflow abstractions (`technical-handbook.md:31-45`, `technical-handbook.md:51-56`, `technical-handbook.md:58-74`, `technical-handbook.md:93-111`).

Complexity introduced: medium. Batch execution needs explicit matrix construction, worker coordination, event draining, result aggregation, and summary semantics.

Complexity removed: low. Reusing single-run and synthesis paths avoids a second prompt/runtime/validation implementation.

Main trade-off: the design favors deterministic partial completion and architectural containment over maximal scheduling sophistication. Synthesis waits until analysis finishes instead of interleaving immediately, and `run-all` does not persist resumable state. This keeps Sprint 11 bounded and leaves durable orchestration for Sprint 12.

## Key Conclusion and Evidence Basis

The final area decision is to implement `run-all` as an ephemeral, study-owned batch workflow in `internal/study`, exposed through a thin `internal/app` command. The workflow should continue unrelated analysis tasks after ordinary task failures, skip synthesis only for blocked dimensions, report partial completion truthfully, and regenerate deterministic summary output. This conclusion is grounded in handbook evidence for thin command delegation, module-owned behavior, manual dependency wiring, bounded concurrency, validation after runtime execution, safe event/diagnostic handling, and fake-runtime tests (`technical-handbook.md:14-27`, `technical-handbook.md:31-45`, `technical-handbook.md:58-74`). It is also required by the sprint constraints and architecture rules that runtime remain generic and batch behavior not move into global scheduler/report/summary/prompt/validation packages (`requirements.md:70-80`, `ARCHITECTURE.md:23-32`, `ARCHITECTURE.md:157-197`).

## Rejected Alternatives

- Global scheduler or workflow engine: rejected because it violates module ownership, creates premature abstraction, and conflicts with explicit non-goals.
- Durable run-loop substrate: rejected because it implies resumability, locks, stale recovery, and retry scheduling outside Sprint 11 scope.
- Fail-fast errgroup cancellation for ordinary task failure: rejected because it undermines partial completion and hides useful work from unrelated tasks.
- Direct OpenCode process handling: rejected because runtime execution must go through agentwrap-backed platform runtime.
- CLI-owned batch logic: rejected because command-layer scheduling would leak product behavior into `internal/app` and make tests harder.

## Risks, Assumptions, and Open Questions

Risks:

- Existing single-run code may not yet expose reusable study-module helpers; if so, a small refactor is required before batch implementation.
- Result aggregation can drift toward durable run-loop state if not kept explicitly in-memory and command-scoped.
- Event draining must be implemented carefully to avoid goroutine leaks or deadlocks in fake and real runtime tests.
- Summary semantics must clearly document the distinction between `N/A` inapplicable cells and empty missing-rating cells.

Assumptions:

- Existing app exit-code conventions can represent success, usage/config failure, validation failure, runtime failure, cancellation, and partial completion.
- Existing validators and prompt builders are callable without duplicating business rules in `run_all.go`.
- The same resolved parallelism cap can safely apply to synthesis runtime calls as well as analysis runtime calls.

Open questions:

- No blocking architecture questions remain for Sprint 11. Implementation should verify exact existing function names, exit-code mapping, and config field names while coding.

## Consolidated Reasoning Handoff

`projects/ultraplan-go/sprints/11-run-all-batch-execution/reasoning.md` should reference this area reasoning as the Architecture decision source for `run-all` package placement, task scheduling semantics, synthesis gating, summary persistence, and rejected durable-orchestration alternatives.
