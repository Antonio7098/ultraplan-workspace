> **Inputs Used:** `projects/ultraplan-go/sprints/10-single-analysis-synthesis/sprint-index.md`, `projects/ultraplan-go/sprints/10-single-analysis-synthesis/technical-handbook.md`, `projects/ultraplan-go/sprints/10-single-analysis-synthesis/requirements.md`, `projects/ultraplan-go/docs/ARCHITECTURE.md`, `projects/ultraplan-go/docs/PRD.md`, `projects/ultraplan-go/docs/TRD.md`, `system/reasoning/architecture_reasoning_template.md`

# Single Analysis And Synthesis Architecture Reasoning

## 0. Feature Summary

### Feature name

Single analysis and synthesis execution.

### Area name and scope

This area covers sprint `10-single-analysis-synthesis`: one selected source/dimension analysis through the platform runtime, one selected dimension synthesis through the platform runtime, runtime request construction, validation boundaries, Markdown applicability skip behavior, synthesis preflight gating, CLI output/error mapping, service/runtime dependency wiring, minimal execution domain fields, and fake-runtime verification.

This area does not cover ceremony. It also does not cover `study run-all`, task matrix execution, worker pools, `study run-loop`, resumable batch mutation, stale-running recovery, retry scheduling across process restarts, per-study locks, automatic synthesis after analysis, summary generation, code extraction, target workflows, sprint workflows, durable agentwrap run-store persistence, or direct OpenCode supervision by UltraPlan.

### User/product goal

Users can run `ultraplan study <study> run <dimension-ref> <source-ref>` to produce and validate one per-source report, and run `ultraplan study <study> synthesize <dimension-ref>` to produce and validate one final report from already valid applicable source reports. The product must accept success only when the expected artifact exists and passes study-owned validation.

### Current task

Create the reasoning basis for the sprint implementation files listed in `requirements.md`: `internal/study/run.go`, `internal/study/synthesize.go`, `internal/study/service.go`, `internal/study/domain.go`, `internal/app/study_commands.go`, `internal/app/app.go`, `internal/study/run_test.go`, and `internal/app/study_run_commands_test.go`.

### Non-goals

- Do not implement `study run-all`, a task matrix scheduler, worker pools, or multi-task scheduling.
- Do not implement `study run-loop`, durable resume mutation, stale-running recovery, lock files, or restart-aware retry scheduling.
- Do not automatically synthesize after a single analysis.
- Do not update `summary.csv` or implement score aggregation.
- Do not add code reference extraction, report repair workflows, target workflows, sprint planning, or sprint execution.
- Do not put study semantics, report semantics, source applicability rules, synthesis gating, or study state machines inside `internal/platform/runtime`.
- Do not reimplement agentwrap/OpenCode process supervision, event decoding, retry/fallback policy, permission translation, health checks, validation wrappers, or native diagnostics parsing.

## 1. First-Principles Breakdown

### Core behaviour

The feature resolves one study work item, builds the existing deterministic prompt for that work item, constructs a generic platform runtime request with safe metadata and validation expectations, executes through the agentwrap-backed runtime boundary, drains/observes runtime completion through that boundary, then validates the expected artifact with study-owned validators before returning a typed result for CLI rendering and exit-code mapping.

For Markdown document analysis, the feature may decide there is no work: an inapplicable Markdown source/dimension pair is a successful skipped result and must not invoke the runtime.

For synthesis, the feature must prove all applicable per-source reports for the selected dimension exist and are valid before runtime launch.

### Inputs

- Workspace root and effective runtime/config values resolved by existing app/config/workspace wiring.
- Study reference, dimension reference, and source reference from CLI arguments.
- Existing study metadata, dimensions, sources, applicability data, report paths, and artifact layout.
- Existing deterministic analysis and synthesis prompt builders.
- Existing per-source and final report validators, including Markdown-source validation rules.
- Platform runtime dependency from `internal/platform/runtime` with agentwrap/OpenCode integration from prior sprint scope.
- Caller `context.Context` for cancellation and timeout propagation.

### Outputs

- For analysis success: a structured result containing completed status, task kind, study, dimension, source, source kind, output path, runtime/provider/model summary, validation summary, and safe runtime diagnostics/usage when available.
- For analysis skip: a structured result containing skipped status, study, dimension, source, source kind, reason, and no runtime run metadata.
- For synthesis success: a structured result containing completed status, task kind, study, dimension, final output path, included applicable source reports, runtime/provider/model summary, validation summary, and safe runtime diagnostics/usage when available.
- For failures: typed or classified errors that distinguish usage/resolution/config/workspace/filesystem failures, runtime failures, validation failures, cancellation, and synthesis preflight blockers.
- CLI text output with deterministic success/skip summaries on stdout and diagnostics on stderr.

### Durable state

- Created: per-source report files and final report files are expected to be created by the runtime as generated artifacts; output directories may be created by study execution before launch if needed.
- Read: workspace config, study definitions, dimensions, sources, prompt templates, report templates, existing per-source reports, and generated runtime output artifacts.
- Updated: no batch run state mutation is required by this sprint; execution result fields should remain compatible with existing task/status vocabulary from prior sprint context.
- Deleted: none.

### Ephemeral state

- Resolved study/dimension/source execution request.
- Applicability decision for Markdown sources.
- Prompt builder output and input manifest.
- Runtime request metadata, work directory, permission policy, health/capability requirements, timeout, provider/model, and expected output validation spec.
- Runtime result, validation summary, safe diagnostics, and error classification.
- Synthesis preflight list of included applicable source reports and blocking invalid/missing reports.

### Derived state

- Task kind is derived from command path: `analysis` for `study run`, `synthesis` for `study synthesize`.
- Task ID may be derived from existing study/task vocabulary using study, dimension, source, source kind, and output path; it is compatibility metadata, not a durable batch scheduler record in this sprint.
- Output paths are derived from study artifact layout and selected dimension/source.
- Source applicability is derived from source kind and normalized `ApplicableDimensions`; it is computed at execution and synthesis preflight time.
- Runtime success status is derived from runtime result plus artifact validation; process/runtime success alone is never the source of product truth.

### Side effects

- File write: expected report artifacts are written by runtime execution; output folders may be created by the study service.
- Network/provider call: possible only through the platform runtime and agentwrap/OpenCode; normal tests must not require network or credentials.
- Queue/event emission: runtime events are handled by the platform runtime abstraction; this sprint does not add durable batch event storage.
- Logs/diagnostics: safe runtime, validation, and command diagnostics with stdout/stderr separation.
- Database writes, email, and external notifications: none.

## 2. Existing Architecture Fit

### Contract applicability and phase gate

Current gate: single runtime-backed study execution slice.

Applies now:

- Architecture: `internal/study` owns study execution, source applicability, prompt use, report validation, and synthesis gating; `internal/app` stays thin; `internal/platform/runtime` stays generic.
- Errors: runtime, validation, skip, usage, workspace, filesystem, ambiguity, missing output, and cancellation outcomes need distinct user-facing diagnostics and exit status mapping.
- Configuration: effective provider/model/runtime/timeout/permission/health settings must be resolved and validated before runtime launch.
- Observability: runtime metadata, task metadata, safe diagnostics, usage/cost knownness, and validation outcomes should be preserved without leaking prompts, secrets, or unsafe native payloads.
- Security: source paths, Markdown document isolation, permission posture, and redaction are in scope.
- Testing: fake-runtime study and command tests are required; default tests must not require OpenCode, provider credentials, network access, or real subprocesses.
- CLI Surface: `study run` and `study synthesize` argument shape, help, deterministic output, stdout/stderr separation, and exit-code mapping are in scope.
- LLM Runtime: runtime requests must use the platform runtime boundary and agentwrap-owned health, permissions, validation, policy, events, and final result semantics.
- LLM Evaluation / Cost / Safety: runtime success is insufficient without product artifact validation; usage/cost metadata is safe and optional.

Deferred:

- Full batch/durable workflow behavior: run-loop mutation, scheduler state transitions, worker pools, locks, stale-running recovery, resume reconciliation, and restart-aware retries.
- Stable public execution JSON schemas: no new JSON schema is required unless explicitly added and tested.
- Report repair product workflows: no new repair loop or repair prompt construction unless already provided by configured runtime wrappers without adding product workflow scope.
- Durable agentwrap run-store persistence across process exits.
- Summary generation, code extraction, target workflows, and sprint workflows.

### Where does this behaviour naturally belong?

Study execution behavior belongs in `internal/study` because `ARCHITECTURE.md` states that study owns sources, dimensions, source applicability, prompt creation, single analysis runs, synthesis, report paths, study validation, and filesystem persistence. The CLI command layer belongs in `internal/app` and should parse args, invoke study service methods, render results, and map errors. Platform runtime remains in `internal/platform/runtime` and accepts only generic prompt/workdir/config/metadata/permission/validation inputs.

### Existing workflow affected

Current flow:

1. Prior sprints provide CLI/app wiring, workspace/config/logging/health foundations, study/dimension/source resolution, study artifact layout, Markdown source discovery/applicability, report validators, deterministic prompt builders, run-state vocabulary compatibility, and platform runtime integration.
2. Runtime execution for study analysis and synthesis is not yet connected; prompt builders and validators exist but are not used to perform one real task.
3. `internal/platform/runtime` is generic and agentwrap-backed from the prior runtime integration sprint.

### Proposed new flow

Proposed `study run` flow:

1. `internal/app` parses `study <study> run <dimension-ref> <source-ref>` and resolves command context/config.
2. `internal/study` resolves the study, dimension, and source using existing resolution behavior.
3. `internal/study` checks Markdown source applicability; if inapplicable, return a successful skipped result without runtime invocation.
4. `internal/study` builds the existing deterministic analysis prompt and output path.
5. `internal/study` chooses work directory: source directory for directory sources, and a safe study/workspace directory for Markdown sources because all Markdown source material is embedded in the prompt.
6. `internal/study` constructs a generic runtime request with prompt, workdir, provider/model/timeout, metadata, permissions, health/capability requirements, and expected output validation.
7. Platform runtime executes through agentwrap/OpenCode, consumes events through its abstraction, and returns a final result or classified error.
8. `internal/study` validates the expected per-source report with the study-owned per-source validator and returns completed or validation failure.
9. `internal/app` renders deterministic stdout/stderr and maps result/error to exit status.

Proposed `study synthesize` flow:

1. `internal/app` parses `study <study> synthesize <dimension-ref>` and resolves command context/config.
2. `internal/study` resolves the study and dimension using existing resolution behavior.
3. `internal/study` computes applicable sources for the dimension, excluding inapplicable Markdown sources.
4. `internal/study` validates each applicable per-source report before runtime launch; missing or invalid reports block synthesis and identify every blocking path.
5. `internal/study` builds the existing deterministic synthesis prompt with the valid applicable report manifest and final output path.
6. `internal/study` constructs a generic runtime request with safe workdir, prompt, config, metadata, permissions, health/capability requirements, and expected output validation.
7. Platform runtime executes and returns a final result or classified error.
8. `internal/study` validates the expected final report with the study-owned final report validator and returns completed or validation failure.
9. `internal/app` renders deterministic stdout/stderr and maps result/error to exit status.

### Does the current architecture still fit?

[x] Yes. The sprint fits the existing module-driven architecture because the selected behavior is a study workflow that can call a generic runtime boundary without moving product semantics into platform runtime.

### If it does not fit, why?

Not applicable. The main risk is implementation drift: command handlers absorbing prompt/runtime/validation logic, or platform runtime learning about studies, dimensions, sources, reports, applicability, or synthesis gates.

### Refactor-before-feature decision

Decision: implement directly with small focused additions to existing `internal/study`, `internal/app`, and app composition code.

Reason: prior sprints already created the necessary seams: study resolution, prompt builders, validators, runtime boundary, config, and app wiring. The technical handbook supports thin commands, explicit composition, fake external systems, context propagation, structured diagnostics, and minimal single-task orchestration without a scheduler refactor.

## 3. Design Options Considered

### Option A: Study-owned single execution over generic runtime boundary

Description:
Keep orchestration in `internal/study`. Add service methods for one analysis and one synthesis. Study builds prompts, validates applicability, performs synthesis preflight, maps product metadata to generic runtime request fields, calls the platform runtime interface, and runs study-owned validators after completion. CLI remains thin.

Pros:
- Matches `ARCHITECTURE.md` ownership and dependency direction.
- Satisfies requirements that prompt builders and validators remain study-owned.
- Keeps `internal/platform/runtime` reusable and free of study/report semantics.
- Supports deterministic fake-runtime tests at both service and command layers.
- Avoids batch scheduler and durable workflow scope.

Cons:
- Requires careful result/error shape so CLI can render actionable diagnostics without duplicating study rules.
- Study code must translate domain metadata into generic runtime metadata fields.
- Validation can be observed in two places if runtime expected-output validation and post-run study validation are not coordinated.

Risks:
- Runtime request construction could grow too large inside one method; mitigate with focused helpers in `internal/study` only when they clarify analysis vs synthesis flow.

### Option B: Put execution logic in CLI command handlers

Description:
Let `internal/app/study_commands.go` resolve inputs, build prompts, call runtime, validate reports, and print outcomes directly.

Pros:
- Fewer study service methods initially.
- CLI argument values are immediately available where execution occurs.

Cons:
- Violates requirement that study execution logic stays in `internal/study` and CLI handlers stay thin.
- Duplicates or exposes prompt and validation rules in the command layer.
- Makes service-level fake-runtime testing harder and weakens future batch/run-loop reuse.

Risks:
- The command layer becomes a long orchestration body, an anti-pattern called out by the technical handbook.

### Option C: Teach platform runtime about study tasks and report validation

Description:
Add study-aware request types or helpers in `internal/platform/runtime` for analysis/synthesis, report validators, output paths, source kinds, synthesis gating, and applicability.

Pros:
- Runtime-facing code could centralize validation and metadata construction.
- Command/study code would call fewer product-specific steps after runtime completion.

Cons:
- Directly violates `ARCHITECTURE.md` and requirements that runtime remains generic and does not import or understand study/report semantics.
- Makes platform runtime change for product workflow changes.
- Risks duplicating agentwrap validation wrappers while also leaking product validation into infrastructure.

Risks:
- Creates reverse coupling from platform infrastructure to study behavior, making later product modules harder to maintain.

### Option D: Introduce a general scheduler/workflow engine now

Description:
Build a reusable task executor that handles analysis, synthesis, event loops, state transitions, retries, locks, summaries, and batch compatibility even for a single task.

Pros:
- May reduce future `run-all` and `run-loop` implementation work.
- Could unify task status fields early.

Cons:
- Exceeds sprint scope and non-goals.
- Adds concurrency, persistence, locking, and retry complexity before a single execution path is proven.
- Conflicts with the handbook's warning against speculative scheduler frameworks.

Risks:
- Review surface expands and single-task acceptance criteria become harder to verify.

### Chosen option

Chosen option: Option A.

Reason: It is the smallest honest design for the selected sprint. It keeps product behavior in `internal/study`, command wiring in `internal/app`, and runtime mechanics in `internal/platform/runtime`. It uses existing prompt builders and validators, relies on agentwrap through the platform boundary, and avoids deferred batch/workflow complexity.

## 4. Abstraction Check

### Are we adding a new abstraction?

[ ] No.
[x] Yes - service/component: single execution methods on the study service.
[x] Yes - data structure/DTO/config object: minimal execution request/result/status/diagnostic fields in `internal/study/domain.go`.
[x] Yes - interface/protocol/trait: use the existing runtime boundary or add only a narrow consumer-owned runtime interface if the current platform type is too broad for tests.
[ ] Yes - plugin/registry/factory.

### If yes, why is it earned?

- It isolates an external dependency: runtime execution is volatile and must be fakeable.
- It protects stable domain concepts: analysis result, synthesis result, skipped work, validation failure, and runtime failure are product-visible outcomes.
- It makes testing simpler: fake runtime can prove request mapping, no-runtime skip, validation gates, and CLI output.
- It creates a clear boundary between product policy and runtime mechanism: study owns applicability, prompts, synthesis gates, and validators; platform runtime owns generic execution.

### If no, why are we keeping it concrete?

Pure helpers for output paths, metadata maps, preflight lists, and result construction should stay concrete inside `internal/study` unless repeated behavior proves they need names. Do not add generic workflow, scheduler, validator, report, prompt, or plugin packages.

### Bad abstraction smell check

- Generic names such as `Manager`, `Processor`, `Common`, or `Helper` are rejected.
- Boolean mode switches inside platform runtime for analysis vs synthesis are rejected; study passes generic metadata and validation expectations.
- A broad task scheduler or workflow engine is rejected because this sprint executes only one selected task.
- A runtime interface that mirrors all platform internals is rejected; tests should depend on the narrow behavior study actually consumes.

## 5. DRY and Duplication Check

### Is any duplication being introduced?

[ ] Yes - duplicated business/domain knowledge.
[x] Yes - duplicated code shape only, where analysis and synthesis both resolve inputs, build prompts, run runtime, and validate outputs.
[ ] Yes - duplicated infrastructure mechanics.

### If duplicating business knowledge, stop and consolidate

No business rule should be duplicated. Single sources of truth are:

- Prompt composition: existing deterministic analysis and synthesis prompt builders.
- Markdown applicability: existing source applicability helper/logic.
- Report validation: existing per-source and final report validators.
- Runtime execution: `internal/platform/runtime` and agentwrap.
- Workspace resolution and path safety: existing workspace/app configuration flow.

### If duplicating code shape, is that acceptable?

Analysis and synthesis have similar high-level shape but different product gates. Analysis has a source applicability skip path and source-kind-specific workdir/prompt behavior. Synthesis has a pre-runtime requirement that all applicable source reports are valid. Keep the flows readable rather than hiding differences behind a generic task executor.

### If removing duplication, is the abstraction safe?

Small private helpers are safe only for identical knowledge: building common runtime metadata keys, converting validation summaries, rendering safe diagnostics, or creating generic expected-output validation. Do not factor analysis and synthesis into a mode-driven function if it obscures distinct preconditions.

## 6. Coupling Check

### Dependencies required

- `internal/study`: owns execution orchestration, prompts, applicability, report paths, validators, and result types.
- `internal/platform/runtime`: provides generic runtime execution, health/capability/permission/validation/result handling, and safe runtime diagnostics.
- `internal/workspace`: used through existing app/study resolution paths for workspace roots and safe paths.
- `internal/platform/config`: provides resolved runtime/provider/model/timeout/permission settings through app composition.
- `internal/app`: wires commands and dependencies, but must not own product execution rules.
- `context.Context`: required for runtime execution and cancellation propagation.

### Coupling risks

Global coupling: none introduced if runtime and service dependencies are passed through constructors rather than globals.

Content coupling: none introduced if `internal/platform/runtime` does not import `internal/study` and study passes only generic metadata/validators.

Stamp coupling: possible if study methods receive an app-wide config or command object. Mitigate by defining focused execution requests and a resolved runtime config value.

### Narrow dependency check

Each function should receive only what it needs. CLI handlers receive arguments and app dependencies. Study service methods receive study/dimension/source refs, workspace/config context, and runtime dependency. Runtime request construction receives prompt text, workdir, metadata, permissions, config, and output validation facts, not full CLI command state.

Large objects should be passed only when the full concept is needed. A `Study`, `Dimension`, and `Source` are acceptable in study internals because the whole product concept is needed; platform runtime should receive flattened metadata strings and generic validation instructions instead.

### Dependency injection check

Side-effectful dependencies are created at the edge and passed inward. `internal/app/app.go` wires the study service with the configured platform runtime. Tests inject fake runtimes and buffer-backed command IO. Study code should not construct OpenCode, agentwrap adapters, or stdout/stderr directly.

## 7. State and Mutation Check

### Mutation points

- `/home/antonioborgerees/coding/ultraplan-go/internal/study/run.go`: creates output directories if needed, invokes runtime for one analysis, and validates the per-source report.
- `/home/antonioborgerees/coding/ultraplan-go/internal/study/synthesize.go`: validates existing per-source reports, invokes runtime for one synthesis, and validates the final report.
- `/home/antonioborgerees/coding/ultraplan-go/internal/study/service.go`: stores the runtime dependency and exposes service methods.
- `/home/antonioborgerees/coding/ultraplan-go/internal/study/domain.go`: adds minimal execution result/status/diagnostic fields.
- `/home/antonioborgerees/coding/ultraplan-go/internal/app/study_commands.go`: renders command output only.
- `/home/antonioborgerees/coding/ultraplan-go/internal/app/app.go`: wires the runtime into the study service.

### Is mutation explicit from names and flow?

[x] Yes, if execution methods are named around running analysis and synthesis and if validation/preflight helpers are clearly separate from runtime launch.

### Are queries and commands separated where practical?

[x] Yes. Resolution, applicability checks, prompt construction, metadata construction, synthesis preflight validation, and report validation can return data/errors. Runtime execution and output directory creation are command side effects.

### Could hidden state make this hard to debug?

[x] No, if runtime dependency, config, workspace root, workdir, output path, metadata, and validation summary are explicit in requests/results. Hidden globals for runtime/config/current workspace would be a risk and are rejected.

## 8. Function/Class/Module Shape

### Primary unit of behaviour

[ ] Function.
[ ] Class/object.
[x] Module.
[x] Service/component.
[x] Pipeline/workflow.
[x] Adapter boundary, through existing platform runtime only.

### Why this unit?

Single analysis and synthesis are cohesive study use cases with resolution, prompt construction, runtime launch, validation gates, and result mapping. A study service method is the right unit because callers need a product-level operation, not internal prompt/validator steps. The runtime adapter boundary already exists and should not be recreated.

### If using a class/object, what justifies it?

Go service structs are justified because the study service groups cohesive dependencies: store/filesystem access, prompt/validation behavior, and runtime execution. It protects the invariant that product success requires artifact validation after runtime completion.

### If using inheritance, why not composition?

Not applicable. Go composition and focused interfaces at volatile boundaries are sufficient.

### If using composition, what components are combined?

- Study resolver/store: resolves study, dimension, source, paths, and existing artifacts.
- Prompt builders: produce deterministic analysis and synthesis prompts.
- Applicability logic: decides Markdown skip and synthesis included sources.
- Report validators: validate per-source and final report artifacts.
- Platform runtime: executes generic prompt requests through agentwrap/OpenCode.
- CLI command layer: renders results and maps errors without owning workflow rules.

## 9. Error Handling Design

### Expected failures

- Usage error: missing or extra CLI arguments; handled in `internal/app` with usage exit status.
- Workspace/config failure: invalid workspace, runtime config, provider/model, timeout, health/capability, or permission settings; fail before runtime launch with config/workspace exit status.
- Study/dimension/source resolution failure: return existing not-found or ambiguous-reference diagnostics.
- Inapplicable Markdown analysis: return successful skipped result, print clear skip message, and do not invoke runtime.
- Prompt input failure: missing prompt/template/dimension/source content fails before runtime launch with actionable path diagnostics.
- Synthesis blocked by missing/invalid source reports: fail before runtime launch and list every blocking report path and failed checks.
- Runtime failure: preserve platform runtime classification and cause chain; map to runtime exit status.
- Runtime success with missing output: fail validation/missing-output path and identify expected artifact path.
- Runtime success with invalid per-source/final report: fail validation exit status and include failed check names and report path.
- Cancellation: preserve context/runtime cancellation classification and map to cancellation exit status.

### Unexpected failures

- Filesystem errors should wrap operation and path.
- Unknown runtime errors should be wrapped with operation context while preserving cause chain and any safe runtime diagnostics.
- Unknown validation errors should identify the validator and artifact path when safe.

### Error taxonomy

Reuse existing project error/exit conventions where present:

- Usage/argument error: command shape or invalid refs supplied by user.
- Workspace/filesystem/config error: environment/setup/path/config preflight failures.
- Runtime error: classified platform runtime/agentwrap/OpenCode failure.
- Validation error: expected artifact missing, empty, or failed study-owned checks.
- Skipped result: successful non-error status for inapplicable Markdown source/dimension pair.
- Cancellation: context or user-initiated cancellation.

### Retry / recovery behaviour

- Retry safe: user may rerun the single command after fixing config/runtime/report issues. This sprint does not schedule automatic retries across process restarts.
- Partial progress possible: runtime may write an invalid or partial report before failing validation; keep the artifact for inspection and report the validation failure.
- Compensation needed: no automatic deletion or rollback. Generated reports are user-reviewable artifacts.

## 10. Observability Design

### Events/logs/metrics/traces needed

- Task started/completed/failed/skipped summaries with study, dimension, source/source kind, output path, runtime, provider, model, and duration when available.
- Runtime request metadata fields for task kind, study, dimension number/slug, source name, source kind, output path, runtime, provider, and model.
- Validation summary with artifact path, pass/fail status, failed check names, and safe observed facts.
- Runtime diagnostics from platform runtime: classified category, operation, safe user detail, provider/model/runtime, attempts/policy, session/run IDs, warnings, usage/cost knownness, and safe artifact references when available.
- Synthesis preflight diagnostics listing every missing or invalid applicable source report.

### Minimum useful fields

- `task_kind`: `analysis` or `synthesis`.
- `study`.
- `dimension_number` and `dimension_slug`.
- `source` and `source_kind` for analysis; omitted or empty for synthesis.
- `output_path`.
- `runtime`, `provider`, and `model`.
- `agent_run_id` and `session_id` when available from runtime.
- `validation_status`, `failed_checks`, and `report_path`.
- `status`: completed, skipped, failed, cancelled, or validation_failed.

### Sensitive data check

Could this expose secrets or user content? Yes.

Mitigation: do not print full prompts, embedded Markdown document content, full environment variables, secret config values, unsafe native payload bytes, or unredacted runtime diagnostics. Use platform runtime safe fields and redaction. Report paths, study/source/dimension identifiers, runtime/provider/model labels, and failed validation check names are safe enough for normal diagnostics.

## 11. Testing Strategy

### Unit tests

- Analysis request mapping for directory sources: workdir is the source directory; prompt builder is used; metadata contains task kind, study, dimension, source, source kind, output path, runtime, provider, and model.
- Analysis request mapping for Markdown sources: workdir is safe study/workspace directory; stripped document prompt builder is used; runtime is not asked to explore external files.
- Markdown inapplicable source/dimension pair returns skipped success and fake runtime invocation count remains zero.
- Synthesis preflight includes only applicable sources and excludes inapplicable Markdown sources.
- Synthesis preflight reports all missing and invalid applicable per-source reports before runtime launch.
- Validation summary mapping for per-source and final report validators includes path and failed check names.

### Use-case/workflow tests

- Successful analysis with fake runtime writing a valid per-source report returns completed status.
- Successful synthesis with fake runtime writing a valid final report returns completed status.
- Runtime failure returns runtime failure status and preserves safe classified diagnostics.
- Runtime success without expected artifact returns validation/missing-output failure.
- Runtime success with invalid report returns validation failure with failed checks.
- Synthesis blocked by invalid or missing source reports does not invoke runtime.
- Cancellation/timeout path maps to cancellation or runtime failure according to existing runtime classification.

### Integration tests

- Command tests should use fake runtime through app composition and buffer-backed IO.
- Any real OpenCode smoke test must be explicitly environment-gated, skipped by default, and routed through `agentwrap/opencode` via `internal/platform/runtime`.

### Regression tests

- Runtime completion alone never marks analysis or synthesis successful.
- CLI stdout contains normal summaries and skip messages; stderr contains runtime/validation diagnostics.
- `internal/platform/runtime` has no study imports and no product-specific validation logic.
- Inapplicable Markdown sources are not required for synthesis input checks.
- Unknown usage/cost values are not rendered as zero.

### Testability check

[x] Yes, the core behaviour can be tested without real infrastructure by injecting fake runtimes, temporary workspaces/studies, fixture reports, existing prompt fixtures, and buffer-backed CLI IO.

## 12. Performance and Scale Assumptions

### Known constraints

- Expected volume: one analysis task or one synthesis task per command invocation.
- Latency sensitivity: dominated by runtime/provider execution; preflight resolution and validation should be fast.
- Memory sensitivity: Markdown document source content is embedded in the prompt as required; do not read source repositories recursively in UltraPlan for this sprint.
- Concurrency concerns: no worker pool; runtime event consumption must not block and `Wait` must complete through the platform abstraction.

### Assumptions being made

- Existing prompt builders and validators are deterministic and ready to use.
- Prior runtime integration exposes enough generic request/result fields for this sprint; minor adapter additions may be needed but no product semantics should enter platform runtime.
- Synthesis of a dimension should use the study or workspace root as a safe workdir unless existing prompt/runtime conventions require a narrower generated-artifact directory.
- A single command rerun is the recovery mechanism for failed validation or runtime errors in this sprint.

### Measurement plan

[x] No benchmark or load test is needed for this change.

Evidence of acceptable performance is that preflight and fake-runtime tests complete deterministically, no worker pool or sleeps are introduced, and event consumption does not leak goroutines or block final results.

### Optimization decision

Do not optimize beyond avoiding recursive scans, using existing resolution helpers, validating only required artifacts, and keeping diagnostics bounded/redacted. Batch throughput and large-run status performance belong to later `run-all`/`run-loop` scope.

## 13. Implementation Plan

### Files likely to change

- `/home/antonioborgerees/coding/ultraplan-go/internal/study/domain.go`: add minimal execution request/result/status, validation summary, skip, runtime diagnostic, and metadata fields needed by single execution.
- `/home/antonioborgerees/coding/ultraplan-go/internal/study/service.go`: wire the runtime dependency and expose service methods for single analysis and synthesis.
- `/home/antonioborgerees/coding/ultraplan-go/internal/study/run.go`: implement study-owned single analysis orchestration.
- `/home/antonioborgerees/coding/ultraplan-go/internal/study/synthesize.go`: implement study-owned single synthesis orchestration and preflight gating.
- `/home/antonioborgerees/coding/ultraplan-go/internal/app/study_commands.go`: wire `study run` and `study synthesize` commands, help, output, and exit status mapping.
- `/home/antonioborgerees/coding/ultraplan-go/internal/app/app.go`: compose the study service with the configured platform runtime.
- `/home/antonioborgerees/coding/ultraplan-go/internal/study/run_test.go`: add fake-runtime study execution tests.
- `/home/antonioborgerees/coding/ultraplan-go/internal/app/study_run_commands_test.go`: add command execution tests with fake runtime and captured IO.

### Step-by-step plan

1. Inspect existing study service, prompt builder, validator, applicability, runtime, config, and command APIs to align exact names and avoid duplicate rules.
2. Define minimal execution result/status fields in `internal/study/domain.go`, reusing existing task/status vocabulary where available.
3. Add a runtime dependency to the study service through constructor/app composition with a fakeable interface no broader than needed.
4. Implement `study run`: resolve inputs, handle Markdown skip, build prompt, choose workdir, build runtime request, invoke runtime, validate expected per-source report, and return structured result.
5. Implement `study synthesize`: resolve study/dimension, compute applicable sources, preflight required per-source reports, build prompt, build runtime request, invoke runtime, validate final report, and return structured result.
6. Coordinate validation boundary: runtime request includes generic expected-output validation, while study-owned report validators remain the definitive product gate and never move into `internal/platform/runtime`.
7. Wire CLI commands with thin parsing/rendering/exit mapping and deterministic stdout/stderr behavior.
8. Add fake-runtime study tests for success, runtime failure, missing output, invalid output, Markdown skip, and synthesis preflight blockers.
9. Add command tests for help/args, output, exit statuses, stdout/stderr separation, skip behavior, validation diagnostics, and fake-runtime paths.
10. Verify implementation later with `go test ./...` and `go build ./cmd/ultraplan` from `/home/antonioborgerees/coding/ultraplan-go`.

### Migration/backwards compatibility

- Schema migration needed: no durable schema migration is selected.
- API compatibility affected: no public Go API is introduced; command behavior is extended for selected commands.
- Existing data affected: existing valid reports are read for synthesis preflight; generated output reports may be overwritten or replaced according to existing runtime/prompt behavior, but no unrelated files should be modified.
- Feature flag needed: no. Optional real-runtime smoke tests require environment gates and must be skipped by default.

### Rollback plan

Revert the new study execution files/changes, app command wiring, app composition update, and related tests as a unit. Because no durable migration or batch state mutation is introduced, rollback should not require data repair. Runtime-produced report artifacts are user files and should not be automatically deleted by rollback.

## 14. Final Pre-Implementation Decision

Decision: Proceed.

Reason:
The selected sprint, requirements, architecture docs, PRD/TRD, and technical handbook converge on one shape: implement single analysis and single synthesis as study-owned service methods over the existing generic platform runtime, with thin CLI wiring, deterministic prompt builders, study-owned validation gates, safe diagnostics, fake-runtime tests, and no batch/workflow expansion.

Complexity introduced: medium.

Complexity removed: medium, because the design connects existing prompt, validation, runtime, and command foundations without creating a scheduler or duplicating agentwrap.

Main trade-off:
Accept a focused execution result/request model and some similar analysis/synthesis flow code to preserve clear product gates and avoid a premature generic workflow engine.

## Key Conclusion and Evidence Basis

The final area decision is to implement single analysis and synthesis in `internal/study`, call the existing generic `internal/platform/runtime` boundary, and treat study-owned artifact validation as the decisive product success gate. CLI code remains a thin adapter, and platform runtime remains product-agnostic.

Evidence basis:

- `ARCHITECTURE.md` states that study owns sources, dimensions, prompt creation, single analysis runs, synthesis, reports, validation, and persistence, while `platform/runtime` owns generic execution only and must not understand study, dimension, source, synthesis gating, report semantics, or study state machines.
- `requirements.md` requires `study run` and `study synthesize`, existing prompt builders, expected output validation, study-owned per-source/final validators, Markdown inapplicability skip, synthesis preflight over applicable source reports, thin CLI wiring, generic platform runtime, event consumption, deterministic output, and fake-runtime tests.
- `TRD.md` requires agentwrap-backed runtime execution, generic `RunRequest` mapping with prompt/workdir/provider/model/timeout/metadata/health/caps/permissions/validation, canonical events, typed `SDKError` classification, report validators, source applicability, no direct OpenCode invocation, and runtime-free default tests.
- `PRD.md` says runtime success is not product success, product workflows stay product-owned, single analysis and synthesis are product requirements, Markdown document sources can be skipped when inapplicable, final synthesis requires valid applicable source reports, and target/sprint workflows are deferred.
- `technical-handbook.md` evidence supports thin CLI delegation, explicit composition roots, generic runtime boundaries, validation as a product gate, fake external systems, context-aware single execution without scheduler machinery, structured safe diagnostics, config preflight, and avoiding direct runtime text parsing or speculative scheduler/plugin frameworks.

## Rejected Alternatives

- Rejected command-contained execution because it would duplicate prompt/validation/applicability rules in `internal/app` and violate thin command guidance.
- Rejected study-aware platform runtime helpers because they would violate the product/platform dependency boundary and make runtime import or encode report semantics.
- Rejected treating runtime success as product success because PRD/TRD/requirements and handbook evidence require artifact existence and validator success.
- Rejected a general scheduler/workflow engine because `run-all`, worker pools, durable run-loop mutation, locks, retry scheduling, and automatic synthesis are explicit non-goals.
- Rejected direct OpenCode subprocess handling or stdout/stderr parsing because agentwrap/opencode owns runtime supervision, native event projection, health, permissions, validation wrappers, policy, and diagnostics.
- Rejected default real-runtime tests because normal verification must not require OpenCode, provider credentials, network access, wall-clock sleeps, or real subprocesses.
- Rejected moving report repair into this sprint because configured runtime wrappers may provide validation behavior, but new product repair prompts/loops are not selected.

## Risks, Assumptions, and Open Questions

### Risks

- The validation boundary can double-report missing or invalid artifacts if generic runtime expected-output validation and post-run study validation both emit separate failures. Mitigation: use runtime validation for generic expected-output observation and study validator summaries as the definitive product gate in final result rendering.
- Runtime request metadata can leak product details if raw structs are passed into platform runtime. Mitigation: pass flattened generic strings/values only.
- Command output can become noisy if runtime events are streamed directly to stdout. Mitigation: keep normal summaries on stdout and diagnostics on stderr; detailed events remain safe runtime/log metadata.
- Synthesis preflight can accidentally require inapplicable Markdown source reports. Mitigation: compute applicable sources first and test this specifically.
- Fake runtime behavior can drift from agentwrap/OpenCode. Mitigation: keep fake assertions focused on UltraPlan request/result behavior and optionally add gated smoke tests later.
- Workdir selection for Markdown document analysis could allow unintended filesystem exploration if the prompt or permission policy is too broad. Mitigation: use the document-specific prompt builder, safe workdir, and restrictive permission posture.

### Assumptions

- Existing prior-sprint implementations provide usable study resolution, source applicability, prompt builders, validators, runtime boundary, config, and command test seams.
- The platform runtime boundary can accept generic metadata and expected-output validation without importing study types.
- Existing exit-code conventions include or can map to usage, config, workspace/filesystem, validation, runtime, and cancellation statuses.
- Synthesis should fail before runtime launch when any applicable source report is invalid; users repair or rerun source analysis manually in this sprint.
- Single execution does not need to persist or mutate batch run state, but result fields should stay compatible with existing task/status vocabulary.

### Open questions

- What exact existing service constructor shape should receive the runtime dependency without broad app refactor?
- Does the current platform runtime API expose a one-call execution method that drains events and waits, or must study use a small helper around `StartRun`, event consumption, and `Wait`?
- What exact existing validation result type should be reused in execution results to avoid duplicate validation summary models?
- What exact permission policy defaults exist after sprint 9, and how restrictive can Markdown document analysis be while still allowing the runtime to write the expected report?
- Does existing CLI infrastructure already define stable exit constants for validation/runtime/cancellation, or should command mapping add only local mappings to existing conventions?

## Consolidation Reference

This area reasoning document is selected by `sprint-index.md` as `projects/ultraplan-go/sprints/10-single-analysis-synthesis/reasoning/single-analysis-synthesis.md` and should be summarized by the sprint's consolidated reasoning artifact at `projects/ultraplan-go/sprints/10-single-analysis-synthesis/reasoning.md`.
