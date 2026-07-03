# Sprint Plan: 12 Durable Run Loop

> Project: `ultraplan-go`
> Sprint: `12-durable-run-loop`
> Source: `projects/ultraplan-go/sprints/12-durable-run-loop/reasoning.md`
> **Inputs Used:** `projects/ultraplan-go/sprints/12-durable-run-loop/requirements.md`, `projects/ultraplan-go/sprints/12-durable-run-loop/reasoning.md`, `projects/ultraplan-go/sprints/12-durable-run-loop/sprint-index.md`, `projects/ultraplan-go/sprints/12-durable-run-loop/technical-handbook.md`, `projects/ultraplan-go/docs/PRD.md`, `projects/ultraplan-go/docs/TRD.md`, `projects/ultraplan-go/docs/ARCHITECTURE.md`, `projects/ultraplan-go/project-index.md`, `templates/sprint-plan.md`

This plan executes `reasoning.md`. It must not invent architecture, scope, or decisions.

No area-specific reasoning file was present under `projects/ultraplan-go/sprints/12-durable-run-loop/reasoning/*.md`; `reasoning.md` records that absence and makes the final decisions directly from the requirements, project docs, sprint index, and technical handbook.

Final study reports were not opened for this plan because `technical-handbook.md` and `reasoning.md` already distilled the selected reports into decision-grade evidence. Open a final report only if implementation encounters a specific decision gap not covered by those artifacts.

## Reasoning Source

- **Sprint Reasoning:** `projects/ultraplan-go/sprints/12-durable-run-loop/reasoning.md`
- **Sprint Index:** `projects/ultraplan-go/sprints/12-durable-run-loop/sprint-index.md`
- **Technical Handbook:** `projects/ultraplan-go/sprints/12-durable-run-loop/technical-handbook.md`
- **Area Reasoning:** none present

## Sprint Status

- **Status:** complete with sprint metadata follow-up
- **Owner:** implementation agent
- **Start Date:** 2026-06-03
- **Completion Date:** 2026-06-03

## Decisions To Execute

| Decision | Source Section | Execution Implication |
| --- | --- | --- |
| Keep run-loop study-owned with thin CLI wiring | `reasoning.md#decision-1-keep-run-loop-study-owned-with-thin-cli-wiring` | Implement orchestration, state machines, validation, scheduling, and synthesis gating in `internal/study`; keep `internal/app/study_commands.go` limited to flags, cancellation setup, service invocation, exit mapping, and rendering. |
| Materialize a deterministic study task state machine | `reasoning.md#decision-2-materialize-a-deterministic-study-task-state-machine` | Build a deterministic applicable analysis matrix and synthesis task state using stable task IDs, explicit statuses, typed outcomes, skipped inapplicable Markdown pairs, and status-ready totals. |
| Use atomic run-state persistence with explicit schema handling | `reasoning.md#decision-3-use-atomic-run-state-persistence-with-explicit-schema-handling` | Load, validate, reject or explicitly migrate supported state versions, and atomically save every meaningful transition with same-directory temp writes, flush/close, rename, and best-effort directory sync. |
| Add per-study run-loop locking with conservative stale handling | `reasoning.md#decision-4-add-per-study-run-loop-locking-with-conservative-stale-handling` | Acquire a scoped study lock before state mutation, record PID/command/timestamp/study, refuse concurrent run-loop, and expose explicit force-unlock behavior without auto-deleting active-looking locks. |
| Reconcile resume by recovering stale running tasks and revalidating completed outputs | `reasoning.md#decision-5-reconcile-resume-by-recovering-stale-running-tasks-and-revalidating-completed-outputs` | On resume, validate state, recover active-like tasks, preserve attempts/history, revalidate completed artifacts, recalculate synthesis gates, and record recovery decisions for status. |
| Schedule with bounded workers and context-driven cancellation | `reasoning.md#decision-6-schedule-with-bounded-workers-and-context-driven-cancellation` | Use resolved parallelism to cap runtime starts, stop feeding work on cancellation, cancel active runs through the runtime boundary, drain events, call `Wait`, persist active outcomes, and use the existing cancellation exit convention. |
| Reuse existing runtime boundary, prompt builders, validators, retry/fallback policy, and metadata | `reasoning.md#decision-7-reuse-existing-runtime-boundary-prompt-builders-validators-retryfallback-policy-and-metadata` | Reuse single-run analysis/synthesis prompt and validation paths, execute through agentwrap-backed platform runtime, persist safe retry/fallback/validation/permission/cleanup/usage/cost metadata, and classify runtime errors structurally. |
| Make status runtime-free, deterministic, and safe | `reasoning.md#decision-8-make-status-runtime-free-deterministic-and-safe` | Update `status` to load durable state and lock facts only, render stable progress/active/retry/failure/recovery diagnostics, and omit prompts, Markdown bodies, secrets, and unsafe native payload bytes. |
| Verify with fake runtime, state/lock failure injection, command tests, race tests, build, and architecture review | `reasoning.md#decision-9-verify-with-fake-runtime-statelock-failure-injection-command-tests-race-tests-build-and-architecture-review` | Add deterministic fake-runtime and fixture tests, inject persistence/lock failures, run full test/race/build commands, and record architecture and sprint review evidence. |

## Requirements / Contracts To Satisfy

| Contract / Requirement ID | Required Behavior | Evidence Planned |
| --- | --- | --- |
| `AC-01`, CLI Surface, Documentation | `ultraplan study <study> run-loop` exists with help for filters, parallelism, retry/cancellation behavior, lock override or force unlock, and deterministic exit handling. | `internal/app/study_run_loop_commands_test.go`; help text assertions. |
| `AC-02`, `AC-09`, Workflows | Initial state creates the deterministic applicable task matrix and excludes inapplicable Markdown source/dimension pairs. | `internal/study/run_loop_test.go`; run-state snapshot assertions. |
| `AC-03`, `AC-04`, Persistence And Migrations | Existing state is loaded and validated before runtime launch; unsupported, malformed, or unsafe schemas fail explicitly. | `internal/study/state_test.go`; resume fixtures for malformed and unsupported state. |
| `AC-05`, `AC-06`, `AC-21`, Persistence And Migrations | Every meaningful transition is saved atomically, and required write failures are loud while preserving the previous valid state. | Failure-injection tests for temp write, flush/close, rename, directory sync where testable, and final error paths. |
| `AC-07`, `AC-08`, Security | Per-study lock prevents concurrent mutation; stale handling is conservative; force unlock is explicit and study-scoped. | `internal/study/locks_test.go`; CLI lock conflict and force-unlock tests. |
| `AC-10`, `AC-12`, LLM Evaluation / Cost / Safety | Completed analysis and synthesis artifacts are revalidated on resume before trust; synthesis starts only after all applicable reports validate. | Resume tests with missing/invalid reports; synthesis unblocking tests. |
| `AC-11`, `AC-16`, `AC-17`, `AC-18`, LLM Runtime | Analysis and synthesis reuse existing execution semantics, runtime boundary, agentwrap policy metadata, and typed runtime error classification. | Fake-runtime metadata tests; import review showing no direct OpenCode usage from product code. |
| `AC-13`, `AC-14`, `AC-15`, `AC-25`, Performance | Scheduler uses bounded workers, drains or consumes events, calls `Wait`, stops scheduling on cancellation, and has race-sensitive coverage. | Blocked fake-runtime concurrency tests; cancellation tests; `go test -race ./...`. |
| `AC-19`, `AC-20`, `AC-22`, Observability, Security | Status is runtime-free, truthful, deterministic, and redacted; output includes useful task totals, paths, retry times, cancellation, recovery, and lock diagnostics. | `internal/app/study_status_commands_test.go`; redaction tests; no-runtime-call test fake. |
| `AC-23`, Architecture | No target scaffolding, sprint planning/execution, code-reference extraction, generic workflow engine, new adapter, or forbidden global packages are introduced. | Architecture review protocol; scope search/review. |
| `AC-24`, Testing | Normal tests use fake runtimes and local fixtures only, with no OpenCode, credentials, network, or subprocess requirement. | Test design review; `go test ./...` in default environment. |
| `AC-26` | CLI builds successfully. | `go build ./cmd/ultraplan`. |
| `AC-27`, Architecture | Run-loop behavior is study-owned; CLI stays thin; `internal/platform/runtime` has no product imports. | Import review; architecture review protocol. |
| `CON-01` through `CON-12` | Preserve module ownership, runtime boundary, context propagation, bounded concurrency, path safety, atomic writes, safe diagnostics, Markdown applicability, and offline verification. | Combined unit, command, race, build, security, and architecture review evidence. |

## Tasks

- [x] **Task 1: Define Run-Loop Request, Result, And CLI Surface**
  > Executes: Decision 1, Decision 8, `AC-01`, `AC-19`, `CON-01`, `CON-04`
  - [ ] Add study service request/result types for `run-loop` in `internal/study`, including selected study, source/dimension filters, resolved parallelism, retry/fallback/cancellation settings, lock options, output mode, and safe config summary fields.
  - [ ] Wire `ultraplan study <study> run-loop` in `internal/app/study_commands.go` as a thin command that resolves config, parses flags, sets cancellation context, delegates to the study service, maps exit codes, and renders deterministic summaries.
  - [ ] Document flags in command help for filters, parallelism, retry/fallback controls, cancellation behavior, lock conflict handling, and explicit force unlock or override behavior.
  - [ ] Add command-level tests for help, argument validation, invalid flag values, config override precedence, service invocation, lock failure rendering, cancellation exit mapping, and no product logic inside CLI handlers.

- [x] **Task 2: Extend Durable Domain State And Deterministic Task Matrix**
  > Executes: Decision 2, Decision 7, `AC-02`, `AC-05`, `AC-09`, `AC-11`, `AC-12`, `CON-03`, `CON-09`
  - [ ] Extend `internal/study/domain.go` with task kind/status/outcome fields needed for pending, running, retrying or waiting, validating, completed, failed, cancelled, skipped, synthesis-unblocked, retry/fallback, runtime metadata, validation metadata, usage/cost, and safe diagnostics.
  - [ ] Define stable task IDs from task kind, dimension, and source where applicable; keep synthesis task IDs dimension-scoped and deterministic.
  - [ ] Build initial run state in `internal/study/run_state.go` from discovered dimensions and applicable sources, reusing source applicability helpers and excluding inapplicable Markdown pairs rather than marking them failed.
  - [ ] Preserve source/dimension filters in run state and make resume behavior explicit when the user supplies different filters.
  - [ ] Add tests proving deterministic ordering, stable task IDs, skipped inapplicable Markdown pairs, synthesis task representation, typed outcome codes/categories, and status summary counts.

- [x] **Task 3: Implement Atomic State Store And Schema Validation**
  > Executes: Decision 3, `AC-03`, `AC-04`, `AC-05`, `AC-06`, `AC-21`, `CON-06`, `CON-07`, `CON-08`
  - [ ] Update `internal/study/state.go` to load and validate `studies/<study>/.ultraplan/run-state.json` with schema version checks and actionable errors before runtime launch.
  - [ ] Implement explicit compatibility or rejection behavior for existing state versions; do not silently reinterpret unknown durable formats.
  - [ ] Save state atomically through same-directory temp file, flush/close, rename, and best-effort directory sync.
  - [ ] Ensure write failures preserve the last known good state and cause `run-loop` to fail loudly rather than reporting success.
  - [ ] Add failure-injection tests for malformed state, unsupported schema, unsafe state path, temp write failure, flush/close failure, rename failure, directory sync behavior where testable, and prior valid state preservation.

- [x] **Task 4: Add Per-Study Lock Handling**
  > Executes: Decision 4, `AC-07`, `AC-08`, `AC-21`, `CON-06`, `CON-08`
  - [ ] Implement `internal/study/locks.go` with per-study lock acquisition before state mutation and lock release on normal terminal exit.
  - [ ] Store lock payload fields for PID, sanitized command summary, timestamp, lock path, and study scope without leaking sensitive arguments.
  - [ ] Refuse a second concurrent `run-loop` with actionable diagnostics and existing lock metadata.
  - [ ] Implement conservative stale-lock detection that reports active-looking or ambiguous locks without deleting them automatically.
  - [ ] Implement explicit, selected-study-scoped force unlock behavior through a flag or command path selected during CLI wiring.
  - [ ] Add `internal/study/locks_test.go` plus CLI tests for lock creation, conflict refusal, stale diagnostics, scoped force unlock, cleanup on terminal exit, and lock persistence/release failures.

- [x] **Task 5: Implement Resume Reconciliation And Artifact Revalidation**
  > Executes: Decision 5, Decision 7, `AC-03`, `AC-09`, `AC-10`, `AC-12`, `AC-19`, `CON-09`
  - [ ] On resume, reconcile stale `running` and active-like tasks to a safe schedulable or terminal state according to the selected policy in `reasoning.md`, preserving attempts, errors, retry/fallback history, and recovery metadata.
  - [ ] Revalidate completed analysis reports with existing per-source validators and move missing or invalid artifacts back to a non-completed or failed validation state with safe diagnostics.
  - [ ] Revalidate completed synthesis reports with existing final-report validators before trusting synthesis completion.
  - [ ] Recalculate synthesis unblocking only from applicable selected per-source reports that exist and validate.
  - [ ] Add resume fixtures and tests for stale running tasks, retrying tasks, completed tasks with missing artifacts, invalid per-source reports, invalid final reports, failed tasks, attempt preservation, recovery decisions, and original filter preservation.

- [x] **Task 6: Implement Bounded Durable Scheduler And Cancellation**
  > Executes: Decision 6, Decision 2, `AC-05`, `AC-13`, `AC-14`, `AC-15`, `AC-25`, `CON-04`, `CON-05`
  - [ ] Implement `internal/study/run_loop.go` scheduling with bounded workers capped by resolved configuration or CLI override and without unbounded goroutine or channel growth.
  - [ ] Persist transitions for pending, running, retrying or waiting, validating, completed, failed, cancelled, skipped, and synthesis-unblocked states.
  - [ ] Stop scheduling new work when context cancellation or user interrupt is observed, including cancellation while feeding worker queues.
  - [ ] Cancel active runtime runs through the platform runtime/agentwrap boundary, drain or safely consume event streams, and call `Wait` for every started run.
  - [ ] Persist active task cleanup outcomes as cancelled or retryable according to policy and return the existing cancellation exit-code convention.
  - [ ] Add blocked fake-runtime tests proving max concurrency, no new work after cancellation, active run cancellation, event draining, `Wait` calls, persisted active outcomes, synthesis unblocking, and race-sensitive behavior.

- [x] **Task 7: Reuse Runtime, Prompt, Validation, Retry/Fallback, And Metadata Paths**
  > Executes: Decision 7, `AC-11`, `AC-12`, `AC-16`, `AC-17`, `AC-18`, `AC-20`, `CON-02`, `CON-03`, `CON-08`
  - [ ] Invoke existing single analysis and synthesis prompt builders, runtime request mapping, workdir rules, permission policy mapping, event handling, report validators, rating parser, and final-report validation instead of duplicating those paths.
  - [ ] Route all runtime execution through the existing platform runtime backed by agentwrap/OpenCode; do not invoke OpenCode directly from product code.
  - [ ] Persist safe agentwrap-facing metadata for attempts, retry/fallback decisions, next retry time, runtime/provider/model target, validation, permission, cleanup, artifacts, usage, cost, warnings, and errors.
  - [ ] Preserve unknown usage and cost values as unknown rather than zero.
  - [ ] Classify runtime errors with typed errors or `agentwrap.SDKError` fields using structural inspection, not string parsing.
  - [ ] Add fake-runtime tests for runtime success without artifact, validation failure, rate limits, retries, fallback target selection, permission events, cleanup metadata, usage/cost known and unknown, cancellation, typed SDK errors, and safe omission of unsafe native payloads.

- [x] **Task 8: Update Runtime-Free Status And Safe Output**
  > Executes: Decision 8, `AC-19`, `AC-20`, `AC-22`, `CON-08`, `CON-12`
  - [ ] Update `ultraplan study <study> status` to load only durable state and lock information, with no runtime calls or OpenCode dependency.
  - [ ] Render deterministic totals and task sections for pending, running, retrying/waiting, validating, completed, failed, cancelled, recovered, skipped, and synthesis-unblocked states.
  - [ ] Include state path, lock diagnostics, active tasks, recent/recovered tasks, retry times, failed validation summaries, cancelled tasks, artifact paths, safe runtime/model/provider metadata, and usage/cost when known.
  - [ ] Redact or omit prompts, embedded Markdown source bodies, API keys, sensitive environment values, and unsafe native payload bytes while recording safe omission facts where useful.
  - [ ] Add `internal/app/study_status_commands_test.go` coverage proving runtime-free behavior, deterministic ordering, safe diagnostics, lock output, recovery output, retry times, unknown usage handling, and redaction.

- [x] **Task 9: Complete Cross-Cutting Test Matrix**
  > Executes: Decision 9, all acceptance criteria, Testing contract
  - [ ] Add `internal/study/run_loop_test.go` coverage for initial state creation, resume, bounded scheduling, retries, fallback metadata, validation, synthesis unblocking, cancellation, stale recovery, terminal states, and partial completion.
  - [ ] Add or update `internal/study/state_test.go` coverage for atomic writes, schema compatibility or rejection, completed-output revalidation inputs, resume persistence, and prior-state preservation on write failure.
  - [ ] Add or update `internal/study/locks_test.go` coverage for acquisition, concurrent refusal, conservative stale-lock behavior, scoped force unlock, and cleanup.
  - [ ] Add `internal/app/study_run_loop_commands_test.go` coverage for command args, flags, help, exit codes, cancellation, lock failures, resume behavior, write failures, and fake-runtime outcomes.
  - [ ] Add redaction tests that include prompt-like text, Markdown source body text, API-key-like values, sensitive env names, and unsafe native payload markers.
  - [ ] Keep all normal tests deterministic, offline, credential-free, network-free, and independent of real subprocess execution.

- [x] **Task 10: Run Verification And Prepare Review Evidence**
  > Executes: Decision 9, `AC-24`, `AC-25`, `AC-26`, `AC-27`
  - [ ] Run `go test ./...` from `/home/antonioborgerees/coding/ultraplan-go` and record results in `review.md`.
  - [ ] Run `go test -race ./...` from `/home/antonioborgerees/coding/ultraplan-go` and record results in `review.md`.
  - [ ] Run `go build ./cmd/ultraplan` from `/home/antonioborgerees/coding/ultraplan-go` and record results in `review.md`.
  - [ ] Perform architecture review using `system/protocols/architecture-review-protocol.md`, including import checks for `internal/platform/runtime`, thin CLI review, and absence of forbidden global packages.
  - [ ] Perform scope review for absence of target commands, sprint commands, code-reference extraction work, generic workflow engine, plugin architecture, new runtime adapter, and broad migration framework.
  - [ ] Complete `projects/ultraplan-go/sprints/12-durable-run-loop/review.md` using `system/protocols/sprint-review-protocol.md` with acceptance evidence, deviations, residual risks, and carry-forward follow-ups.

## Evidence Checklist

- [ ] Run-loop command help and CLI tests prove command availability, flags, usage errors, lock failures, resume failures, cancellation, runtime failures, validation failures, and partial completion.
- [ ] State-machine tests prove deterministic task matrix creation, Markdown applicability exclusion, stable task IDs, source/dimension filter preservation, status totals, and synthesis task unblocking.
- [ ] Resume tests prove state schema validation, unsupported/malformed state rejection, stale running recovery, completed-output revalidation, attempt/history preservation, recovery metadata, and synthesis gate recalculation.
- [ ] Atomic persistence tests prove same-directory temp writes, flush/close, rename, best-effort directory sync where testable, loud write failures, and prior valid state preservation.
- [ ] Lock tests prove per-study locking, concurrent refusal, conservative stale diagnostics, explicit scoped force unlock, cleanup on terminal exit, and lock failure diagnostics.
- [ ] Concurrency and cancellation tests prove bounded runtime starts, event draining, `Wait` calls, active run cancellation, no new scheduling after cancellation, persisted active outcomes, and no goroutine leaks under race testing.
- [ ] Runtime metadata tests prove retry, rate-limit, fallback, validation, permission, cleanup, usage, cost, unknown usage, and typed SDK error data are persisted safely.
- [ ] Status tests prove runtime-free loading, deterministic output, progress counts, active/recent tasks, retry times, failed/cancelled tasks, recovery decisions, state path, lock diagnostics, safe metadata, and redaction.
- [ ] Architecture review proves study-owned run-loop behavior, thin CLI wiring, no product imports in `internal/platform/runtime`, and no forbidden global `scheduler`, `workflow`, `prompts`, `reports`, or `validation` packages.
- [ ] Scope review proves no target scaffolding, sprint planning, sprint execution, code-reference extraction, generic workflow engine, plugin system, new runtime adapter, or broad migration framework entered this sprint.
- [ ] Verification commands pass or any deviation is recorded in `review.md` with cause, impact, and follow-up.

## Verification Commands

| Check | Command | Expected Result |
| --- | --- | --- |
| Full test suite | `go test ./...` | Passes without OpenCode, provider credentials, network access, or real subprocess execution. |
| Race-sensitive paths | `go test -race ./...` | Passes, including run-loop workers, cancellation, locks, state persistence, and fake-runtime tests. |
| CLI build | `go build ./cmd/ultraplan` | Produces the CLI binary without build errors. |
| Architecture review | Manual review using `system/protocols/architecture-review-protocol.md` | Confirms study ownership, thin CLI, platform runtime product-agnostic imports, and no forbidden global packages. |
| Sprint review | Manual review using `system/protocols/sprint-review-protocol.md` | `review.md` records acceptance evidence, deviations, residual risks, and follow-ups. |

## Risks And Blockers

| Risk / Blocker | Source | Mitigation | Status |
| --- | --- | --- | --- |
| Existing Sprint 8 run-state primitives may need more extension than expected. | `reasoning.md#assumptions-and-risks` | Keep changes focused in `domain.go`, `run_state.go`, and `state.go`; use explicit schema rejection instead of broad migrations. | open |
| Agentwrap metadata may not expose every desired safe retry/fallback/usage/cost field. | `reasoning.md#assumptions-and-risks` | Persist all available safe metadata and omission facts; record metadata gaps in `review.md`. | open |
| Completed-output revalidation may slow resume for large studies. | `reasoning.md#assumptions-and-risks` | Prefer correctness now; consider validation fingerprints only after measured cost. | open |
| Conservative stale-lock handling may require operator intervention after dead processes. | `reasoning.md#assumptions-and-risks` | Provide clear diagnostics and explicit selected-study force unlock; defer heartbeat/process verification. | open |
| Cancellation during persistence or cleanup may still leave stale active tasks. | `reasoning.md#assumptions-and-risks` | Persist before and after critical transitions where possible; recover stale running tasks on next resume; make write failures loud. | open |
| Rich metadata can leak unsafe data if fields are copied too broadly. | `reasoning.md#assumptions-and-risks` | Persist safe summaries, redaction results, and omission facts; add redaction tests for prompts, Markdown bodies, secrets, env values, and unsafe payloads. | open |
| Race-sensitive worker and lock tests may be flaky without deterministic fakes. | `reasoning.md#assumptions-and-risks` | Use blocking fake runtimes with explicit synchronization and verify with `go test -race ./...`. | open |
| Public JSON stability is deferred but users may depend on current state/status shape. | `reasoning.md#assumptions-and-risks` | Version durable state, avoid promising public JSON compatibility, and carry schema questions to Sprint 14. | open |

## Review Inputs

Review should use:

- `projects/ultraplan-go/sprints/12-durable-run-loop/requirements.md`
- `projects/ultraplan-go/sprints/12-durable-run-loop/sprint-index.md`
- `projects/ultraplan-go/sprints/12-durable-run-loop/technical-handbook.md`
- `projects/ultraplan-go/sprints/12-durable-run-loop/reasoning.md`
- `projects/ultraplan-go/sprints/12-durable-run-loop/plan.md`
- implementation diff under `/home/antonioborgerees/coding/ultraplan-go`
- verification evidence from `go test ./...`, `go test -race ./...`, and `go build ./cmd/ultraplan`
- `system/protocols/architecture-review-protocol.md`
- `system/protocols/sprint-review-protocol.md`

## Execution Log

| Date / Step | Action | Evidence / Notes |
| --- | --- | --- |
| 2026-06-02 planning | Created implementation plan from `requirements.md`, `reasoning.md`, sprint index, technical handbook, PRD, TRD, architecture doc, project index, and sprint plan template. | Implementation not started. No code changes planned or made by this document. |
| 2026-06-03 implementation | Added `ultraplan study <study> run-loop` thin CLI wiring, study-owned durable scheduler, per-study lock file handling, lock-aware runtime-free status, and focused run-loop/lock tests. | Implementation files changed under `/home/antonioborgerees/coding/ultraplan-go`: `internal/app/study_commands.go`, `internal/study/domain.go`, `internal/study/locks.go`, `internal/study/run_loop.go`, `internal/study/locks_test.go`, `internal/study/run_loop_test.go`. Existing dirty `go.mod`/`go.sum` changes were present and not modified by this sprint. |
| 2026-06-03 verification | Baseline `go test ./...` failed because `/tmp` was full; reran required checks with `GOTMPDIR=$PWD/.gotmp`. | `GOTMPDIR=$PWD/.gotmp go test ./...` passed. `GOTMPDIR=$PWD/.gotmp go test -race ./...` passed. `GOTMPDIR=$PWD/.gotmp go build ./cmd/ultraplan` passed. |
| 2026-06-03 review | Performed architecture and scope review. | Run-loop behavior remains in `internal/study`; CLI code delegates through service; no `internal/scheduler`, `internal/workflow`, global prompts/reports/validation packages, target commands, sprint commands, code-reference extraction work, generic workflow engine, plugin system, or new runtime adapter were added. Selected sprint review protocol path in the plan was missing as `sprint-review-protocol.md`; used existing `system/protocols/review-sprint-protocol.md`. |
| 2026-06-03 initial deferrals | Recorded breadth gaps from the original full sprint scope before follow-up implementation. | Superseded by the 2026-06-03 follow-up implementation row, except for the sprint metadata protocol filename follow-up. |
| 2026-06-03 follow-up implementation | Cleared `/tmp`, inspected `/home/antonioborgerees/coding/agentwrap`, and implemented the non-metadata follow-ups. | Added safe agentwrap metadata persistence for attempts, policy decisions, permissions, cleanup, repair, usage, cost, artifacts, warnings, and omission facts; expanded status task sections; added run-loop CLI tests for help, invalid flags, lock conflict, force unlock, config persistence, cancellation mapping, and status metadata. Plain `go test ./...`, `go test -race ./...`, and `go build ./cmd/ultraplan` passed without `GOTMPDIR`. |

## Completion Criteria

- [x] All implementation tasks are complete or explicitly deferred in `review.md` with rationale and impact.
- [x] `go test ./...` passed or deferral is documented with blocker, impact, and follow-up.
- [x] `go test -race ./...` passed or deferral is documented with blocker, impact, and follow-up.
- [x] `go build ./cmd/ultraplan` passed or deferral is documented with blocker, impact, and follow-up.
- [x] Evidence satisfies the expectations from `reasoning.md`, including tests, status output, metadata, locking, persistence, cancellation, and architecture review.
- [x] Deviations from `reasoning.md` were recorded before implementation continued beyond the deviating decision.
- [x] `review.md` can evaluate conformance without guessing intent.
