> **Inputs Used:** `projects/ultraplan-go/project-index.md`, `projects/ultraplan-go/sprints/11-run-all-batch-execution/requirements.md`, `projects/ultraplan-go/docs/PRD.md`, `projects/ultraplan-go/docs/TRD.md`, `projects/ultraplan-go/docs/ARCHITECTURE.md`, `projects/ultraplan-go/sprints/11-run-all-batch-execution/sprint-index.md`, `projects/ultraplan-go/sprints/11-run-all-batch-execution/technical-handbook.md`, `projects/ultraplan-go/sprints/11-run-all-batch-execution/reasoning/run-all-batch-execution.md`, `projects/ultraplan-go/sprints/11-run-all-batch-execution/reasoning.md`, `templates/sprint-plan.md`

# Sprint Plan: 11 Run All Batch Execution

> Project: `ultraplan-go`
> Sprint: `11-run-all-batch-execution`
> Source: `reasoning.md`
> Output: `projects/ultraplan-go/sprints/11-run-all-batch-execution/plan.md`

This plan executes `reasoning.md`. It must not invent architecture, scope, or decisions.

## Reasoning Source

- **Sprint Reasoning:** `projects/ultraplan-go/sprints/11-run-all-batch-execution/reasoning.md`
- **Sprint Index:** `projects/ultraplan-go/sprints/11-run-all-batch-execution/sprint-index.md`
- **Technical Handbook:** `projects/ultraplan-go/sprints/11-run-all-batch-execution/technical-handbook.md`
- **Area Reasoning:** `projects/ultraplan-go/sprints/11-run-all-batch-execution/reasoning/run-all-batch-execution.md`
- **Project Docs Used:** `projects/ultraplan-go/docs/PRD.md`, `projects/ultraplan-go/docs/TRD.md`, `projects/ultraplan-go/docs/ARCHITECTURE.md`
- **Final Reports:** Not opened directly for this plan because `technical-handbook.md` already distilled the selected evidence reports and no remaining planning decision required deeper report evidence.

## Sprint Status

- **Status:** complete
- **Owner:** OpenCode implementation agent
- **Start Date:** 2026-06-02
- **Completion Date:** 2026-06-02

## Decisions To Execute

| Decision | Source Section | Execution Implication |
| --- | --- | --- |
| Study-owned ephemeral batch command | `reasoning.md#decision-1-study-owned-ephemeral-batch-command` | Add `run-all` orchestration in `internal/study`; keep `internal/app` thin; do not persist run-loop state, locks, retries, stale recovery, or migrations. |
| Deterministic matrix, filters, applicability, and preflight | `reasoning.md#decision-2-deterministic-matrix-filters-applicability-and-preflight` | Resolve study, dimensions, sources, filters, config, and `parallelism >= 1` before any runtime task starts; build a stable applicable task matrix. |
| Bounded runtime workers with partial completion and cancellation | `reasoning.md#decision-3-bounded-runtime-workers-with-partial-completion-and-cancellation` | Use a local per-request worker pool; continue unrelated tasks after ordinary failures; stop new scheduling on cancellation; drain events and wait for every started run. |
| Reuse existing runtime, prompt, workdir, permission, and validation paths | `reasoning.md#decision-4-reuse-existing-runtime-prompt-workdir-permission-and-validation-paths` | Call or minimally refactor existing single-analysis and synthesis study helpers; do not duplicate prompt/runtime/validation rules or parse OpenCode directly. |
| Post-analysis synthesis gate and deterministic summary | `reasoning.md#decision-5-post-analysis-synthesis-gate-and-deterministic-summary` | Run synthesis only after all applicable selected per-source reports for a dimension validate; generate deterministic `summary.csv` with missing, ambiguous, and inapplicable semantics. |
| Safe deterministic CLI output, result status, and verification | `reasoning.md#decision-6-safe-deterministic-cli-output-result-status-and-verification` | Return a study-domain result with counts, warnings, artifacts, safe diagnostics, and internal errors; map to existing exit conventions and deterministic text output. |

## Requirements / Contracts To Satisfy

| Contract / Requirement ID | Required Behavior | Evidence Planned |
| --- | --- | --- |
| SR-GOAL, AC-01 | `ultraplan study <study> run-all` exists, appears in help, and returns stable status/exit behavior. | Command tests in `internal/app/study_run_all_commands_test.go`; help output assertions; exit-code assertions. |
| AC-02, AC-03, AC-04 | Deterministic task matrix from existing source/dimension discovery; filters reuse existing resolution and fail before launch when invalid or ambiguous. | Study tests in `internal/study/run_all_test.go`; fake runtime asserts zero starts for preflight failures. |
| AC-05, AC-06 | Parallelism is resolved from config or flag override, must be at least 1, and bounds runtime task starts. | Blocked fake-runtime concurrency tests; invalid parallelism command and service tests. |
| AC-07, AC-08, AC-12 | Analysis and synthesis reuse existing prompt composition, runtime request mapping, workdir rules, permission policy, health/capability requirements, and validators. | Fake-runtime request assertions where practical; validation gate tests; import and duplication review. |
| AC-09, AC-10 | Failed analysis tasks preserve safe diagnostics and internal causes without hiding successes; partial completion reports completed, failed, skipped, and pending/not-run counts. | Run-all result status tests; command output tests; safe diagnostic assertions. |
| AC-11, AC-13 | Synthesis starts only when all applicable selected per-source reports for a dimension validate; blocked synthesis is reported as skipped. | Synthesis gating tests for all-valid, failed, missing, and inapplicable Markdown cases. |
| AC-14, AC-15, AC-16 | Summary runs after eligible work finishes, writes `studies/<study>/summary.csv`, is deterministic, and distinguishes missing from inapplicable. | Exact CSV tests in `internal/study/summary_test.go`; repeated-run stability tests; warning assertions. |
| AC-17, AC-18 | Cancellation stops new scheduling and every started runtime run has events consumed and `Wait` called. | Fake-runtime delayed event/cancellation tests; no-start-after-cancel assertions; event drain assertions. |
| AC-19, AC-20 | No durable run-loop, durable retry, stale recovery, locks, force unlock, or production-grade resumability; optional state only if strictly needed for truthfulness. | Architecture review; changed-path review; search for forbidden run-loop/durable orchestration additions. |
| AC-21 | Human output is concise, deterministic, includes totals/artifact paths, and excludes prompts, Markdown bodies, secrets, env values, and unsafe native payloads. | Command output tests; redaction/safe diagnostic tests; manual output review. |
| AC-22, AC-23 | Normal verification is offline with fake runtimes/local fixtures; `go test ./...` and `go build ./cmd/ultraplan` pass. | Verification commands from `/home/antonioborgerees/coding/ultraplan-go`. |
| AC-24, Architecture contract | Batch behavior stays in `internal/study`, CLI wiring stays thin, runtime remains generic, and no global technical-layer packages are added. | Import review; changed-path review; architecture review protocol. |

## Tasks

- [x] **Task 1: Inspect Existing Study Execution Seams**
  > Executes: Decision 1, Decision 4; AC-07, AC-08, AC-12, AC-24
  - [x] Inspect existing single analysis, synthesis, prompt builders, report validators, rating parser, source applicability helper, runtime request mapping, config resolution, and app exit-code conventions.
  - [x] Identify the smallest `internal/study` helper refactors needed so `run-all` can reuse single-run behavior without command-layer dependencies.
  - [x] Confirm no new runtime abstraction is needed beyond the existing platform runtime seam and test fakes.
  - [x] Record any implementation-discovered naming or exit-code questions in `review.md` rather than changing the sprint architecture.

- [x] **Task 2: Add Minimal Study Domain Surface**
  > Executes: Decision 1, Decision 6; AC-01, AC-09, AC-10, AC-17, AC-21
  - [x] Add minimal `RunAllRequest`, `RunAllResult`, task result, synthesis result, warning, count, and status types in `internal/study/domain.go` or another existing study-owned file.
  - [x] Keep result types in-memory and command-scoped; do not add persistence fields that imply resumability, attempts history, locks, or retry schedules.
  - [x] Include enough internal cause retention for tests/review while exposing only safe diagnostics for CLI rendering.
  - [x] Add result status calculation tests for success, validation failure, runtime failure, partial completion, cancellation, skipped synthesis, and pending/not-run tasks.

- [x] **Task 3: Implement Preflight, Filters, And Deterministic Matrix**
  > Executes: Decision 2; AC-02, AC-03, AC-04, AC-06, AC-16
  - [x] Implement `RunAll` preflight in `internal/study`: load the selected study, discover dimensions and sources, resolve dimension/source filters through existing resolution and ambiguity behavior, and validate effective parallelism before launching runtime work.
  - [x] Build a stable analysis matrix from selected dimensions and selected applicable sources only, using existing Markdown applicability semantics.
  - [x] Decide and test the concrete ordering used by the implementation, preserving determinism across repeated runs.
  - [x] Ensure inapplicable Markdown source/dimension pairs are skipped, not failed, not missing expected reports, and not counted in summary totals.
  - [x] Add tests for valid filters, unknown filters, ambiguous filters, directory sources, Markdown sources with and without `applicable_dimensions`, and invalid parallelism with zero fake-runtime starts.

- [x] **Task 4: Reuse Analysis Runtime And Validation Path**
  > Executes: Decision 3, Decision 4; AC-05, AC-07, AC-08, AC-09, AC-18, AC-22
  - [x] Wire each analysis task through the existing prompt composition, source-kind workdir rules, runtime request mapping, permission policy mapping, health/capability requirements, and per-source validator.
  - [x] Treat runtime success as incomplete until the expected per-source report exists and passes validation.
  - [x] Drain or safely consume runtime events and call `Wait` for every started run.
  - [x] Preserve original errors internally while recording safe diagnostics for result rendering.
  - [x] Add fake-runtime tests for runtime success with valid report, runtime success with missing report, invalid report, runtime failure, event consumption, and safe diagnostics.

- [x] **Task 5: Implement Bounded Workers, Partial Completion, And Cancellation**
  > Executes: Decision 3, Decision 6; AC-05, AC-09, AC-10, AC-17, AC-18
  - [x] Implement a per-request bounded worker pool or equivalent local concurrency primitive capped by resolved parallelism.
  - [x] Continue scheduling unrelated analysis tasks after ordinary runtime or validation failures.
  - [x] On context cancellation, stop scheduling new tasks, let active runtime calls return through the runtime boundary, drain active event streams, wait for active runs, and mark never-started tasks pending/not-run.
  - [x] Avoid package-level semaphores, long-lived workers, unbounded goroutine-per-task launch, unbounded channels, and background task leaks.
  - [x] Add blocked fake-runtime tests proving maximum simultaneous analysis starts never exceeds the configured limit and cancellation starts no extra tasks.

- [x] **Task 6: Implement Synthesis Eligibility And Execution**
  > Executes: Decision 4, Decision 5; AC-11, AC-12, AC-13, AC-18
  - [x] After the analysis phase, compute synthesis eligibility per selected dimension from applicable selected per-source reports only.
  - [x] Run synthesis only for dimensions whose applicable selected reports exist and pass validation.
  - [x] Reuse existing synthesis prompt composition, runtime request mapping, event consumption, and final-report validation.
  - [x] Report synthesis skipped for dimensions blocked by failed or missing applicable analysis reports without presenting skipped synthesis as success.
  - [x] Add tests for eligible synthesis, failed analysis blocking only its dimension, missing reports blocking synthesis, inapplicable Markdown pairs excluded from required inputs, and failed final-report validation.

- [x] **Task 7: Implement Deterministic Summary Generation**
  > Executes: Decision 5; AC-14, AC-15, AC-16, AC-21
  - [x] Implement or update `internal/study/summary.go` so `run-all` writes `studies/<study>/summary.csv` after all eligible analysis and synthesis work finishes.
  - [x] Use existing rating parser semantics and produce warnings for missing or ambiguous ratings without inventing values.
  - [x] Keep output deterministic: stable dimension columns, stable source row ordering, stable total calculation, and stable warning ordering.
  - [x] Represent inapplicable Markdown source/dimension cells with the documented `N/A` sentinel from `reasoning.md`, and exclude those cells from missing-report warnings and totals.
  - [x] Add exact CSV tests, repeated-run stability tests, total sorting tests, missing rating tests, ambiguous rating warning tests, missing expected report tests, and inapplicable Markdown tests.

- [x] **Task 8: Add Thin CLI Wiring And Deterministic Rendering**
  > Executes: Decision 1, Decision 6; AC-01, AC-03, AC-04, AC-06, AC-10, AC-21, AC-23
  - [x] Add `ultraplan study <study> run-all` in `internal/app/study_commands.go` with documented flags for source filters, dimension filters, and parallelism override using existing command style.
  - [x] Keep command code limited to argument parsing, command-level validation, config/flag handoff, service invocation, exit mapping, and text rendering.
  - [x] Map result statuses to existing app exit-code conventions for success, usage/config/preflight failure, validation failure, runtime failure, cancellation, and partial completion.
  - [x] Render concise deterministic text with completed, failed, skipped, pending/not-run counts, synthesis status, warning count, safe diagnostics, and artifact paths.
  - [x] Add command-level tests for help, missing args, invalid flags, filter failures before runtime launch, invalid parallelism before runtime launch, success, runtime failure, validation failure, cancellation, partial completion, deterministic output, artifact paths, and redaction.

- [x] **Task 9: Architecture And Scope Review Before Verification**
  > Executes: Decision 1 through Decision 6; AC-19, AC-20, AC-24
  - [x] Inspect changed paths and imports to confirm batch behavior lives in `internal/study`, CLI remains in `internal/app`, and `internal/platform/runtime` has no product imports.
  - [x] Confirm no global `internal/scheduler`, `internal/reports`, `internal/summary`, `internal/prompts`, `internal/validation`, or equivalent technical-layer package was introduced.
  - [x] Search implementation for absence of `run-loop`, durable retry scheduling, stale-running recovery, lock files, force unlock, target/sprint commands, workflow engines, DAG systems, plugin systems, direct OpenCode parsing, and default real OpenCode smoke tests.
  - [x] Confirm generated writes are explicit workspace-contained paths and `summary.csv` is under the selected study.
  - [x] Record any deviations, accepted trade-offs, or follow-ups in `review.md` before final acceptance.

- [x] **Task 10: Run Verification And Complete Review Evidence**
  > Executes: AC-22, AC-23 and selected review protocols
  - [x] Run `go test ./...` from `/home/antonioborgerees/coding/ultraplan-go` and record the result in `review.md`.
  - [x] Run `go build ./cmd/ultraplan` from `/home/antonioborgerees/coding/ultraplan-go` and record the result in `review.md`.
  - [x] Run or test `ultraplan study <study> run-all --help` behavior and record command help evidence.
  - [x] Complete architecture review protocol evidence for module ownership, dependency direction, runtime boundary, and forbidden package/scope exclusions.
  - [x] Complete sprint review protocol evidence against every acceptance criterion in `requirements.md`.

## Evidence Checklist

- [x] `internal/study/run_all_test.go` covers matrix creation, filters, Markdown applicability, bounded parallelism, runtime success/failure, missing/invalid report validation, synthesis gating, cancellation, event consumption, partial completion, and safe diagnostics.
- [x] `internal/study/summary_test.go` covers exact CSV output, repeated-run stability, source row ordering, total calculations, missing ratings, ambiguous ratings, missing expected reports, and `N/A` inapplicable Markdown cells.
- [x] `internal/app/study_run_all_commands_test.go` covers help, arguments, flags, invalid filters, invalid parallelism, exit mapping, deterministic output, artifact paths, fake-runtime behavior, cancellation, partial completion, and redaction.
- [x] Fake runtime records prove maximum concurrent starts do not exceed configured parallelism and every started run has its events consumed and `Wait` called.
- [x] Review proves runtime success without a valid artifact does not count as product success.
- [x] Review proves prompts, embedded Markdown bodies, secrets, sensitive environment values, direct OpenCode stdout/stderr, and unsafe native runtime payloads are not printed.
- [x] Review proves no durable `run-loop` substrate, retry scheduler, stale recovery, lock files, target/sprint commands, workflow engine, DAG system, plugin system, new runtime adapter, or direct OpenCode parsing entered this sprint.
- [x] Required verification commands pass or any failure is recorded with cause and follow-up in `review.md`.
- [x] Deviations from `reasoning.md` are recorded before implementation continues.
- [x] Required architecture and sprint review protocols have evidence.

## Verification Commands

| Check | Command | Expected Result |
| --- | --- | --- |
| Full test suite | `go test ./...` | Passes offline without OpenCode, provider credentials, network access, or real subprocess execution. |
| CLI build | `go build ./cmd/ultraplan` | Builds the CLI binary successfully. |
| Command help | `go test ./...` with command help assertions, or `go run ./cmd/ultraplan study <study> run-all --help` against a local fixture if the project test style supports it | Help includes `run-all` and documented source, dimension, and parallelism flags. |
| Architecture review | Inspect imports and changed paths; search for forbidden packages/scope | Batch behavior is study-owned, CLI wiring is thin, platform runtime is generic, and future-scope exclusions are absent. |

## Risks And Blockers

| Risk / Blocker | Source | Mitigation | Status |
| --- | --- | --- | --- |
| Existing single analysis and synthesis code may not expose reusable study helpers. | `reasoning.md#assumptions-and-risks` | Refactor the smallest reusable helpers into `internal/study` before batch worker code; do not duplicate prompt/runtime/validation rules. | open |
| Existing app exit-code conventions may not directly represent every requested status. | `reasoning.md#assumptions-and-risks` | Inspect current app mapping in Task 1 and test actual conventions; if a gap exists, record the decision and keep behavior consistent with TRD exit classes. | open |
| Result aggregation may drift into durable run-loop state. | `reasoning.md#assumptions-and-risks` | Keep result types in-memory, command-scoped, and non-persisted; defer durable orchestration to Sprint 12. | open |
| Event draining and cancellation may leak goroutines or deadlock. | `reasoning.md#assumptions-and-risks` | Add delayed-event, blocked-runtime, cancellation, and wait/drain fake-runtime tests before relying on real runtime behavior. | open |
| Safe diagnostics may accidentally include prompts, Markdown bodies, secrets, env values, or unsafe native payloads. | `reasoning.md#assumptions-and-risks` | Render only safe categories, validation check names, redacted fields, status counts, and artifact paths; add redaction assertions. | open |
| `N/A` summary sentinel may conflict with future consumer expectations. | `reasoning.md#assumptions-and-risks` | Treat `N/A` as the Sprint 11 documented sentinel, cover it with exact CSV tests, and record future format changes as follow-up work. | open |
| Static matrix could grow for very large studies. | `reasoning.md#assumptions-and-risks` | Accept for Sprint 11 local CLI scale; inspect for accidental source-tree scans or runtime payload accumulation. | open |

## Review Inputs

Review should use:

- `projects/ultraplan-go/sprints/11-run-all-batch-execution/requirements.md`
- `projects/ultraplan-go/sprints/11-run-all-batch-execution/sprint-index.md`
- `projects/ultraplan-go/sprints/11-run-all-batch-execution/technical-handbook.md`
- `projects/ultraplan-go/sprints/11-run-all-batch-execution/reasoning/run-all-batch-execution.md`
- `projects/ultraplan-go/sprints/11-run-all-batch-execution/reasoning.md`
- `projects/ultraplan-go/sprints/11-run-all-batch-execution/plan.md`
- Implementation diff under `/home/antonioborgerees/coding/ultraplan-go`
- Verification evidence from `go test ./...`, `go build ./cmd/ultraplan`, and command-level tests
- `system/protocols/architecture-review-protocol.md`
- `system/protocols/sprint-review-protocol.md`

## Execution Log

| Date / Step | Action | Evidence / Notes |
| --- | --- | --- |
| 2026-06-02 / Planning | Created implementation plan from sprint requirements, consolidated reasoning, area reasoning, sprint index, handbook, PRD, TRD, ARCHITECTURE, project index, and sprint-plan template. | No code implementation performed. Final reports were not opened directly because handbook evidence was sufficient for all planning decisions. |
| 2026-06-02 / Implementation | Implemented study-owned `RunAll`, deterministic summary CSV generation, selected-source synthesis reuse, and thin CLI wiring for `ultraplan study <study> run-all`. | Changed implementation files under `/home/antonioborgerees/coding/ultraplan-go`: `internal/study/run_all.go`, `internal/study/summary.go`, `internal/study/domain.go`, `internal/study/prompts.go`, `internal/study/synthesize.go`, `internal/app/study_commands.go`, and related tests. |
| 2026-06-02 / Verification | Ran required verification from `/home/antonioborgerees/coding/ultraplan-go`. | `go test ./...` passed. `go build ./cmd/ultraplan` passed. Command help is covered by `TestStudyRunAllCommandHelpSuccessAndSummary`. |
| 2026-06-02 / Architecture Review | Inspected changed paths and searched for forbidden scope additions. | Batch behavior stayed in `internal/study`; CLI wiring stayed in `internal/app`; no new global scheduler/report/prompt/validation package or durable run-loop substrate was introduced. |

## Completion Criteria

- [x] All tasks are complete or explicitly deferred in `review.md` with requirement impact recorded.
- [x] `go test ./...` was run from `/home/antonioborgerees/coding/ultraplan-go`, or a deferral is documented with cause and follow-up.
- [x] `go build ./cmd/ultraplan` was run from `/home/antonioborgerees/coding/ultraplan-go`, or a deferral is documented with cause and follow-up.
- [x] Evidence satisfies the expectations from `reasoning.md`, especially validation gates, bounded workers, event draining, summary determinism, safe diagnostics, and architecture boundaries.
- [x] `review.md` can evaluate conformance without guessing intent.
