# Sprint Plan: 10 Single Analysis and Synthesis

> Project: `ultraplan-go`
> Sprint: `10-single-analysis-synthesis`
> Source: `projects/ultraplan-go/sprints/10-single-analysis-synthesis/reasoning.md`
> **Inputs Used:** `projects/ultraplan-go/project-index.md`, `projects/ultraplan-go/sprints/10-single-analysis-synthesis/requirements.md`, `projects/ultraplan-go/docs/ARCHITECTURE.md`, `projects/ultraplan-go/docs/PRD.md`, `projects/ultraplan-go/docs/TRD.md`, `projects/ultraplan-go/sprints/10-single-analysis-synthesis/sprint-index.md`, `projects/ultraplan-go/sprints/10-single-analysis-synthesis/technical-handbook.md`, `projects/ultraplan-go/sprints/10-single-analysis-synthesis/reasoning/single-analysis-synthesis.md`, `projects/ultraplan-go/sprints/10-single-analysis-synthesis/reasoning.md`, `templates/sprint-plan.md`

This plan executes `reasoning.md`. It must not invent architecture, scope, or decisions.

Linked final evidence reports were not reopened for this plan because `technical-handbook.md` already distills the selected reports and cites the relevant source references for the decisions below. Reopen a final report only if implementation uncovers a specific decision that the handbook does not answer.

## Reasoning Source

- **Sprint Reasoning:** `projects/ultraplan-go/sprints/10-single-analysis-synthesis/reasoning.md`
- **Sprint Index:** `projects/ultraplan-go/sprints/10-single-analysis-synthesis/sprint-index.md`
- **Technical Handbook:** `projects/ultraplan-go/sprints/10-single-analysis-synthesis/technical-handbook.md`
- **Area Reasoning:** `projects/ultraplan-go/sprints/10-single-analysis-synthesis/reasoning/single-analysis-synthesis.md`

## Sprint Status

- **Status:** `complete`
- **Owner:** `implementation agent`
- **Start Date:** `2026-06-01`
- **Completion Date:** `2026-06-01`

## Decisions To Execute

| Decision | Source Section | Execution Implication |
| --- | --- | --- |
| Study-owned single execution | `reasoning.md#decision-1-study-owned-single-execution` | Put analysis and synthesis orchestration in `internal/study`; keep `internal/app` limited to parsing, delegation, rendering, and exit mapping. |
| Existing prompt builders and applicability rules are the only prompt/skip sources of truth | `reasoning.md#decision-2-existing-prompt-builders-and-applicability-rules-are-the-only-promptskip-sources-of-truth` | Call existing deterministic prompt builders and existing Markdown applicability behavior; do not duplicate prompt text or frontmatter rules. |
| Generic agentwrap-backed runtime request mapping | `reasoning.md#decision-3-generic-agentwrap-backed-runtime-request-mapping` | Build generic platform runtime requests with prompt, workdir, provider, model, timeout, metadata, permissions, health/caps, expected output validation, and context propagation. |
| Study-owned validation is the product success gate | `reasoning.md#decision-4-study-owned-validation-is-the-product-success-gate` | Treat runtime success as insufficient; validate expected per-source and final reports with existing study-owned validators after runtime completion and before success. |
| Typed outcomes, safe diagnostics, and deterministic CLI output | `reasoning.md#decision-5-typed-outcomes-safe-diagnostics-and-deterministic-cli-output` | Return structured statuses for completed, skipped, runtime failed, validation failed, preflight blocked, and cancelled outcomes; keep stdout summaries separate from stderr diagnostics. |
| Fake-runtime first verification | `reasoning.md#decision-6-fake-runtime-first-verification` | Cover default execution paths with fake runtimes and buffer-backed IO; normal tests must not require OpenCode, credentials, network, wall-clock sleeps, or real subprocesses. |
| Defer batch, workflow, repair, summary, code extraction, and target/sprint scope | `reasoning.md#decision-7-defer-batch-workflow-repair-summary-code-extraction-and-targetsprint-scope` | Do not add `run-all`, `run-loop`, worker pools, locks, automatic synthesis, summary updates, code extraction, repair loops, target/sprint behavior, or stable public execution JSON schemas. |

## Requirements / Contracts To Satisfy

| Contract / Requirement ID | Required Behavior | Evidence Planned |
| --- | --- | --- |
| `requirements.md:31-42` | `study run` resolves study, dimension, and source; skips inapplicable Markdown sources; builds the existing prompt; maps runtime request fields; validates the per-source report; distinguishes runtime and validation failures. | `internal/study/run_test.go`, `internal/app/study_run_commands_test.go`, fake-runtime request assertions, CLI output/exit tests. |
| `requirements.md:43-49` | `study synthesize` resolves study and dimension; uses the existing synthesis prompt builder; preflights applicable source reports; excludes inapplicable Markdown sources; validates final report after runtime completion. | `internal/study/run_test.go`, synthesis preflight tests, missing/invalid report tests, final report validation tests. |
| `requirements.md:50-52`, Architecture contract | Study workflow behavior stays in `internal/study`; `internal/platform/runtime` stays generic and product-agnostic; runtime classifications are preserved without string parsing. | Import review, code review of package boundaries, tests using fake runtime through service/app seams. |
| `requirements.md:53-54`, CLI Surface contract | User-facing output is deterministic, stdout/stderr are separated, and JSON output is not added unless documented and tested. | Command tests for success, skip, validation diagnostics, runtime diagnostics, usage/help, and output streams. |
| `requirements.md:55-60`, Testing contract | Fake-runtime tests cover required success/failure/skip/preflight paths; default verification requires no OpenCode/provider/network/subprocess; `go test ./...` and `go build ./cmd/ultraplan` pass. | Study tests, app command tests, final verification commands from `/home/antonioborgerees/coding/ultraplan-go`. |
| `requirements.md:62-73`, sprint non-goals | No batch execution, durable run-loop, automatic synthesis, summary updates, code extraction, repair product workflow, target/sprint workflow, or duplicated OpenCode supervision. | Scope review during implementation and sprint review. |
| `requirements.md:75-87`, Security/Runtime/Configuration contracts | Respect dependency direction; route execution through platform runtime/agentwrap; fail before launch for invalid config/health/caps/permissions; redact unsafe diagnostics; keep tests deterministic. | Fake-runtime request assertions, config/preflight tests where feasible, diagnostics review, no direct OpenCode invocation. |
| PRD sections `2.3`, `2.5.1`, `2.6`; TRD sections `9A`, `10`, `11`, `14`, `15`, `20`, `22`, `23`, `24` | Implement the sprint slice of single analysis/synthesis, Markdown applicability, prompt composition, agentwrap runtime mapping, validation, event handling, error classification, permissions, fake-runtime testing, and build. | Tests, build, import review, sprint review evidence. |

## Tasks

- [x] **Task 1: Inspect Existing APIs And Seams**
  > Executes: Decisions 1-7; requirements `31-60`; area reasoning open questions.
  - [x] Inspect existing study service, domain, source/dimension resolution, applicability, prompt builders, report validators, runtime package, app composition, command wiring, exit-code conventions, and test seams before editing.
  - [x] Resolve the exact constructor shape for passing the runtime dependency into the study service without broad app refactor.
  - [x] Resolve whether the platform runtime exposes a one-call execution method or whether study must start a run, drain events, and wait through the platform abstraction.
  - [x] Resolve the existing validation result type to reuse so execution results do not duplicate validator summaries unnecessarily.
  - [x] Resolve existing permission policy defaults and exit constants; if missing, add only the smallest generic or local additions needed by this sprint.

- [x] **Task 2: Define Minimal Execution Domain Types**
  > Executes: Decisions 1, 5, 7; requirements `24-25`, `39`, `41-42`, `49`, `52-53`.
  - [x] Add only the request/result/status fields needed for single analysis, single synthesis, skipped Markdown work, runtime diagnostics, validation summaries, and synthesis preflight blockers in `internal/study/domain.go` or existing domain files.
  - [x] Reuse existing task/status vocabulary from sprint 8 where available; do not introduce durable batch schemas, attempt histories, scheduler state, or public JSON schemas.
  - [x] Keep runtime metadata as safe flattened values: task kind, study, dimension number/slug, source name, source kind, output path, runtime, provider, and model.

- [x] **Task 3: Wire Runtime Dependency Through Study Service And App Composition**
  > Executes: Decisions 1, 3, 6; requirements `24`, `27`, `38`, `47`, `50-52`, `75-80`.
  - [x] Update `internal/study/service.go` so the service receives a fakeable platform runtime dependency without constructing OpenCode or agentwrap internals inside study logic.
  - [x] Update `/home/antonioborgerees/coding/ultraplan-go/internal/app/app.go` to pass the configured platform runtime into the study service through explicit composition.
  - [x] If a narrow consumer-owned test interface is needed, keep it local to the consuming package and do not define a competing public runtime contract.

- [x] **Task 4: Implement Single Analysis Orchestration In Study Module**
  > Executes: Decisions 1-5, 7; requirements `22`, `31-42`, `50-53`, `75-85`.
  - [x] Add `/home/antonioborgerees/coding/ultraplan-go/internal/study/run.go` behavior for resolving study, dimension, and source using existing resolution behavior.
  - [x] Check Markdown applicability before prompt/runtime launch; return a successful skipped result and ensure runtime invocation count stays zero for inapplicable pairs.
  - [x] Call the existing deterministic analysis prompt builder; do not compose prompt text in CLI or runtime code.
  - [x] Use source directory workdir for directory sources; use safe study/workspace workdir for Markdown sources with stripped document content embedded by the prompt builder.
  - [x] Build the generic runtime request with prompt, workdir, provider/model/timeout, metadata, permission policy, health/caps, and expected output validation.
  - [x] Consume or drain runtime events so they cannot block final result retrieval; always obtain final runtime result through the platform runtime abstraction.
  - [x] After runtime completion, require expected per-source report existence, non-empty content, and existing study-owned per-source validator success before returning completed status.
  - [x] Preserve classified runtime failures and cause chains; return validation failures with report path and failed check names.

- [x] **Task 5: Implement Single Synthesis Orchestration In Study Module**
  > Executes: Decisions 1-5, 7; requirements `23`, `43-49`, `50-53`, `75-85`.
  - [x] Add `/home/antonioborgerees/coding/ultraplan-go/internal/study/synthesize.go` behavior for resolving study and dimension using existing resolution behavior.
  - [x] Compute applicable sources for the selected dimension and exclude inapplicable Markdown sources from required input checks.
  - [x] Preflight every applicable per-source report before runtime launch; fail with all blocking missing/invalid report paths and failed checks.
  - [x] Call the existing deterministic synthesis prompt builder with the valid applicable report manifest and final output path.
  - [x] Build the generic runtime request with safe study/workspace workdir, prompt, provider/model/timeout, metadata, permission policy, health/caps, and expected output validation.
  - [x] After runtime completion, require expected final report existence, non-empty content, and existing study-owned final report validator success before returning completed status.
  - [x] Preserve runtime classifications and validation diagnostics without moving report semantics into `internal/platform/runtime`.

- [x] **Task 6: Wire Thin CLI Commands And Deterministic Output**
  > Executes: Decisions 1, 5, 7; requirements `26`, `31`, `37`, `41-43`, `49-54`.
  - [x] Update `/home/antonioborgerees/coding/ultraplan-go/internal/app/study_commands.go` to implement `ultraplan study <study> run <dimension-ref> <source-ref>` and `ultraplan study <study> synthesize <dimension-ref>` parsing and help text.
  - [x] Keep command handlers limited to argument validation, calling study service methods, rendering results, and mapping errors/results to exit statuses.
  - [x] Print normal completed and skipped summaries to stdout; print runtime, validation, preflight, config/workspace, usage, and cancellation diagnostics to stderr.
  - [x] Do not add JSON output unless the schema, help text, and command tests are included in this sprint.

- [x] **Task 7: Add Study-Level Fake-Runtime Tests**
  > Executes: Decisions 2-6; requirements `28`, `35-40`, `44-48`, `52`, `55-58`.
  - [x] Add `/home/antonioborgerees/coding/ultraplan-go/internal/study/run_test.go` coverage for analysis success, synthesis success, runtime failure, successful runtime with missing output, successful runtime with invalid per-source report, successful runtime with invalid final report, inapplicable Markdown skip, and synthesis blocked by missing/invalid source reports.
  - [x] Assert fake runtime request fields for prompt text, workdir, provider, model, timeout, metadata, permission policy, required health/caps, expected output validation, and context propagation.
  - [x] Assert inapplicable Markdown analysis and blocked synthesis do not invoke the fake runtime.
  - [x] Use temporary workspaces/fixtures and existing prompt builders/validators; avoid OpenCode, provider credentials, network, real subprocesses, sleeps, or global IO.

- [x] **Task 8: Add CLI Execution Tests**
  > Executes: Decisions 1, 5, 6; requirements `29`, `31`, `37`, `41-43`, `49`, `53-56`.
  - [x] Add `/home/antonioborgerees/coding/ultraplan-go/internal/app/study_run_commands_test.go` coverage for command arguments, help text, success output, skip output, validation diagnostics, runtime diagnostics, exit statuses, stdout/stderr separation, and fake-runtime execution paths.
  - [x] Confirm validation failures map to validation exit status, runtime failures map to runtime exit status, usage errors map to usage exit status, and inapplicable Markdown skip exits successfully.
  - [x] Assert output omits prompt bodies, embedded Markdown document content, secrets, full environment values, unsafe native payloads, and unknown usage/cost values rendered as zero.

- [x] **Task 9: Architecture And Scope Review Before Final Verification**
  > Executes: Decisions 1, 3, 4, 7; requirements `50-52`, `62-87`.
  - [x] Confirm `internal/platform/runtime` imports no product modules and contains no study, dimension, source, report, applicability, or synthesis-gating semantics.
  - [x] Confirm `internal/app` contains no prompt construction, report validation, source applicability, synthesis preflight, or runtime request ownership beyond thin delegation/rendering.
  - [x] Confirm no scheduler, worker pool, durable run-loop mutation, lock file, automatic synthesis, summary update, code extraction, product repair loop, target/sprint workflow, direct OpenCode invocation, or native runtime text parser was introduced.

- [x] **Task 10: Run Verification And Prepare Review Evidence**
  > Executes: Decisions 6-7; requirements `58-60`, `110-123`.
  - [x] Run `go test ./...` from `/home/antonioborgerees/coding/ultraplan-go` with no OpenCode/provider credentials required.
  - [x] Run `go build ./cmd/ultraplan` from `/home/antonioborgerees/coding/ultraplan-go`.
  - [x] Record verification results, deviations, unresolved open questions, and carry-forward risks in `projects/ultraplan-go/sprints/10-single-analysis-synthesis/review.md` during sprint review.

## Evidence Checklist

- [x] Study service tests prove analysis success, synthesis success, runtime failure, runtime success with missing output, runtime success with invalid output, Markdown inapplicable skip with zero runtime invocations, synthesis blocked by missing reports, and synthesis blocked by invalid reports.
- [x] Runtime request mapping tests prove prompt text, workdir, provider, model, timeout, metadata, permission policy, health/caps, expected output validation, and context propagation are passed through the fake runtime.
- [x] CLI command tests prove `study run` and `study synthesize` argument behavior, help text, success output, skip output, validation diagnostics, runtime diagnostics, exit statuses, and stdout/stderr separation.
- [x] Architecture review proves study execution behavior lives in `internal/study`, CLI wiring stays thin in `internal/app`, platform runtime remains generic, and no product imports enter `internal/platform/runtime`.
- [x] Diagnostics review proves no prompt bodies, embedded Markdown bodies, secrets, full environment values, unsafe native payloads, or fake zero usage/cost values appear in normal output.
- [x] Scope review proves no batch, run-loop, scheduler, summary, code extraction, repair, target, sprint, or direct OpenCode supervision scope was added.
- [x] Verification commands pass from `/home/antonioborgerees/coding/ultraplan-go`: `go test ./...` and `go build ./cmd/ultraplan`.
- [x] Optional real-runtime smoke tests, if added, are environment-gated, skipped by default, and routed through `agentwrap/opencode` via `internal/platform/runtime`.
- [x] Any deviation from `reasoning.md` is recorded before implementation continues.
- [x] Required architecture and sprint review protocol evidence is available for `review.md`.

## Verification Commands

| Check | Command | Expected Result |
| --- | --- | --- |
| Full default test suite | `go test ./...` | Passes in `/home/antonioborgerees/coding/ultraplan-go` without OpenCode, provider credentials, network access, wall-clock sleeps, or real subprocess execution. |
| CLI build | `go build ./cmd/ultraplan` | Builds the CLI binary from `/home/antonioborgerees/coding/ultraplan-go`. |
| Optional real runtime smoke | Environment-gated test command selected by implementation, if any | Skipped by default; when explicitly enabled, uses `agentwrap/opencode` through `internal/platform/runtime`. |

## Risks And Blockers

| Risk / Blocker | Source | Mitigation | Status |
| --- | --- | --- | --- |
| Existing prompt builder or validator APIs may need small adapter changes. | `reasoning.md#assumptions-and-risks` | Reused existing prompt builders and validators from `internal/study`; execution code calls them directly. | Closed |
| Platform runtime may not expose exactly the generic execution shape needed for request metadata, expected output validation, event draining, and final result retrieval. | `reasoning.md#assumptions-and-risks`; area reasoning open questions | Existing `StartRun(context.Context, runtime.Request)` already drains events and waits; study adds no product semantics to runtime. | Closed |
| Runtime expected-output validation and study post-run validation may double-report failures. | `reasoning.md#assumptions-and-risks` | CLI renders study validation diagnostics as the product gate and does not duplicate runtime validation details. | Closed |
| Synthesis preflight may accidentally require reports for inapplicable Markdown sources. | `reasoning.md#assumptions-and-risks` | Synthesis preflight uses `GetApplicableSources`; tests cover skipped Markdown and blocked reports without runtime calls. | Closed |
| Markdown document analysis could encourage external filesystem exploration. | `reasoning.md#assumptions-and-risks` | Markdown analysis uses the existing document prompt builder, study workdir, permission mapping, and skip behavior. | Closed |
| Runtime event handling could block final result retrieval. | `reasoning.md#assumptions-and-risks` | Existing platform runtime adapter drains events and waits; study calls only the generic runtime boundary. | Closed |
| Existing exit-code conventions may not fully cover validation, runtime, cancellation, usage, config, and workspace outcomes. | Area reasoning open questions | Reused existing `ExitValidation`, `ExitRuntime`, `ExitCancel`, `ExitUsage`, `ExitConfig`, and `ExitWorkspace`; command tests cover mappings. | Closed |
| Fake runtime may diverge from real agentwrap/OpenCode behavior. | `reasoning.md#assumptions-and-risks` | Accepted residual risk; default evidence focuses on UltraPlan-owned mapping/gating and no real smoke was added. | Carried forward |
| Future batch/run-loop may need expanded state fields. | `reasoning.md#assumptions-and-risks` | Execution result fields are local/minimal and no durable batch schema was added. | Closed |

## Review Inputs

Review should use:

- `projects/ultraplan-go/sprints/10-single-analysis-synthesis/sprint-index.md`
- `projects/ultraplan-go/sprints/10-single-analysis-synthesis/technical-handbook.md`
- `projects/ultraplan-go/sprints/10-single-analysis-synthesis/reasoning/single-analysis-synthesis.md`
- `projects/ultraplan-go/sprints/10-single-analysis-synthesis/reasoning.md`
- `projects/ultraplan-go/sprints/10-single-analysis-synthesis/plan.md`
- Implementation diff in `/home/antonioborgerees/coding/ultraplan-go`
- Verification evidence from `go test ./...` and `go build ./cmd/ultraplan`
- `system/protocols/architecture-review-protocol.md`
- `system/protocols/sprint-review-protocol.md`

## Execution Log

| Date / Step | Action | Evidence / Notes |
| --- | --- | --- |
| 2026-06-01 planning | Created implementation plan from requirements, sprint reasoning, sprint index, technical handbook, area reasoning, PRD, TRD, architecture doc, project index, and sprint-plan template. | Plan intentionally does not implement code; final evidence reports omitted because `technical-handbook.md` provided sufficient selected-study evidence for planning. |
| 2026-06-01 implementation | Implemented study-owned single analysis and synthesis execution, runtime dependency injection, CLI run/synthesize commands, fake-runtime service tests, and command tests. | Changed `/home/antonioborgerees/coding/ultraplan-go/internal/study/{domain.go,service.go,run.go,synthesize.go,run_test.go}` and `/home/antonioborgerees/coding/ultraplan-go/internal/app/{study_commands.go,study_run_commands_test.go}`. |
| 2026-06-01 verification | Ran required verification. | `go test ./...` passed; `go build ./cmd/ultraplan` passed. Architecture review: `internal/platform/runtime` has no product imports; execution orchestration is in `internal/study`; CLI execution path renders/delegates only. Existing prompt preview command still renders prompts as prior scope. |

## Completion Criteria

- [x] All tasks are complete or explicitly deferred with a reason tied to `reasoning.md`, `requirements.md`, or a named open question.
- [x] Verification commands were run from `/home/antonioborgerees/coding/ultraplan-go` or deferrals are documented.
- [x] Evidence satisfies `reasoning.md#expected-evidence` and the evidence checklist in this plan.
- [x] `review.md` can evaluate conformance without guessing intent.
- [x] Any deviation from decisions 1-7 is recorded and justified before implementation proceeds.
