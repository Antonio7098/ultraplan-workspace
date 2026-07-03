# Sprint Requirements: 11 Run All Batch Execution

> Project: `ultraplan-go`
> Sprint: `11-run-all-batch-execution`
> Purpose: the authoritative, human-readable sprint contract. All other sprint artifacts must satisfy these requirements.

## Sprint Goal

Implement `ultraplan study <study> run-all` so a study can execute all selected applicable analysis tasks with bounded parallelism, validate outputs, synthesize completed dimensions, and generate a deterministic summary without introducing durable `run-loop` orchestration.

## Required Outputs

| Output | Path | Description |
| ---------- | -------- | --------------- |
| Requirements | `projects/ultraplan-go/sprints/11-run-all-batch-execution/requirements.md` | This sprint contract defining scope, outputs, acceptance criteria, constraints, dependencies, and review expectations. |
| Sprint index | `projects/ultraplan-go/sprints/11-run-all-batch-execution/sprint-index.md` | Selected contracts, evidence reports, docs, protocols, and exclusions for the batch execution sprint. |
| Technical handbook | `projects/ultraplan-go/sprints/11-run-all-batch-execution/technical-handbook.md` | Sprint-specific implementation guidance for bounded workers, runtime requests, validation gates, synthesis gating, summary generation, and fake-runtime testing. |
| Area reasoning | `projects/ultraplan-go/sprints/11-run-all-batch-execution/reasoning/run-all-batch-execution.md` | Detailed reasoning for batch orchestration, filters, concurrency, state interaction, synthesis gating, summary semantics, and trade-offs. |
| Consolidated reasoning | `projects/ultraplan-go/sprints/11-run-all-batch-execution/reasoning.md` | Decision log summarizing the chosen design and mapping decisions to requirements and selected evidence. |
| Implementation plan | `projects/ultraplan-go/sprints/11-run-all-batch-execution/plan.md` | Ordered implementation checklist with verification commands and review evidence expectations. |
| Sprint review | `projects/ultraplan-go/sprints/11-run-all-batch-execution/review.md` | Completed review recording acceptance evidence, deviations, residual risks, and follow-ups. |
| Batch execution domain and orchestration | `/home/antonioborgerees/coding/ultraplan-go/internal/study/run_all.go` | Study-owned `run-all` orchestration for applicable task matrix creation, filters, bounded workers, validation gates, synthesis scheduling, and summary generation. |
| Summary generation implementation | `/home/antonioborgerees/coding/ultraplan-go/internal/study/summary.go` | Deterministic summary generation from validated per-source reports, including missing and inapplicable semantics. |
| Batch execution domain updates | `/home/antonioborgerees/coding/ultraplan-go/internal/study/domain.go` | Minimal request, result, task, summary, or status types required by `run-all`, kept in the study module. |
| Study service wiring | `/home/antonioborgerees/coding/ultraplan-go/internal/study/service.go` | Public study service method and runtime dependency wiring needed to invoke `run-all`. |
| CLI command wiring | `/home/antonioborgerees/coding/ultraplan-go/internal/app/study_commands.go` | Thin CLI handling for `ultraplan study <study> run-all` and documented flags. |
| Study batch tests | `/home/antonioborgerees/coding/ultraplan-go/internal/study/run_all_test.go` | Fake-runtime tests for task matrix creation, bounded parallelism, filtering, success, failure, synthesis gating, and Markdown applicability. |
| Summary tests | `/home/antonioborgerees/coding/ultraplan-go/internal/study/summary_test.go` | Deterministic summary tests covering ratings, totals, missing cells, inapplicable Markdown pairs, warnings, and output ordering. |
| CLI batch tests | `/home/antonioborgerees/coding/ultraplan-go/internal/app/study_run_all_commands_test.go` | Command-level tests for arguments, flags, exit codes, output, and fake-runtime behavior. |

## Acceptance Criteria

- [ ] `ultraplan study <study> run-all` is available in help output and returns stable success, usage, validation, runtime, partial-completion, and cancellation/error exit statuses consistent with existing app exit-code conventions.
- [ ] `run-all` builds a deterministic analysis task matrix from discovered dimensions and sources, using existing applicability rules so inapplicable Markdown document source/dimension pairs are skipped and never counted as failures or missing expected reports.
- [ ] `run-all` supports filtering by one or more dimension references and one or more source references, with the same resolution and ambiguity behavior as existing single-run commands.
- [ ] Invalid, unknown, or ambiguous source and dimension filters fail before launching any runtime task and include actionable diagnostics.
- [ ] `run-all` uses bounded parallelism from resolved configuration or a CLI override and never starts more concurrent analysis runtime tasks than the configured limit.
- [ ] Parallelism values less than 1 fail during command validation before any runtime task starts.
- [ ] Each analysis task reuses existing analysis prompt composition, runtime request mapping, source-kind workdir rules, runtime health/capability requirements, permission policy mapping, and post-runtime report validation from Sprint 10.
- [ ] Runtime success alone is insufficient: each analysis task is successful only when the expected per-source report exists and passes the existing per-source validator.
- [ ] A failed analysis task records a safe diagnostic in the batch result, preserves the original error for internal inspection, and does not hide successful tasks.
- [ ] Batch execution returns a partial-completion status when one or more tasks fail after other tasks succeed, and the text output clearly reports completed, failed, skipped, and pending/not-run counts.
- [ ] Synthesis for a dimension starts only after all applicable selected per-source reports for that dimension exist and pass validation.
- [ ] Synthesis reuses existing synthesis prompt composition, runtime request mapping, and final-report validation from Sprint 10.
- [ ] Synthesis is skipped for a dimension with failed or missing applicable analysis reports, and skipped synthesis is reported without being presented as success.
- [ ] Summary generation runs after all eligible analysis and synthesis work finishes, writes `studies/<study>/summary.csv`, and reports warnings when ratings are missing or ambiguous.
- [ ] The summary is deterministic: dimensions and sources are ordered consistently, totals are stable, and output is identical across repeated runs with identical reports.
- [ ] The summary distinguishes missing expected reports from inapplicable Markdown source/dimension pairs and excludes inapplicable pairs from missing-report warnings and total calculations.
- [ ] The command respects `context.Context` cancellation by stopping new task scheduling, allowing active runtime calls to return through the platform runtime boundary, and reporting cancellation/partial completion without starting extra tasks.
- [ ] Runtime event channels are drained or otherwise consumed for every started runtime run so workers cannot deadlock on unread events.
- [ ] `run-all` does not introduce durable retry, fallback beyond existing agentwrap policy use, stale-running recovery, per-study lock files, or the `run-loop` command.
- [ ] `run-all` may update existing run-state primitives only if needed for truthful status or reviewability, but it must not claim production-grade resumability or replace Sprint 12 durable orchestration.
- [ ] Human output is concise and deterministic, includes task totals and artifact paths, and does not print prompts, embedded Markdown source bodies, secrets, sensitive environment values, or unsafe native runtime payloads.
- [ ] Normal tests use fake runtimes and local fixtures only; `go test ./...` passes without OpenCode, provider credentials, network access, or real subprocess execution.
- [ ] The CLI builds with `go build ./cmd/ultraplan`.
- [ ] Architecture review confirms study-owned batch behavior, thin CLI wiring, no product imports from `internal/platform/runtime`, and no global `scheduler`, `reports`, `summary`, `prompts`, or `validation` packages.

## Non-Goals

- Implementing `ultraplan study <study> run-loop` or production-grade resumable orchestration.
- Implementing retry/backoff scheduling outside agentwrap's existing policy boundary.
- Implementing stale running task recovery, per-study lock files, force unlock, or multi-process coordination.
- Implementing stable public JSON output for `run-all`, `status`, or summary commands unless already present and required by existing app conventions.
- Implementing code-reference extraction or the `ultraplan code <report>...` command.
- Implementing target scaffolding, sprint planning, sprint execution, or any target/sprint CLI commands.
- Adding new runtime adapters, bypassing `agentwrap/opencode`, or directly parsing OpenCode stdout/stderr.
- Adding real OpenCode smoke tests to the default test suite.
- Adding a generic workflow engine, DAG authoring system, plugin system, hosted service, browser UI, or multi-user collaboration features.

## Constraints

- Batch behavior must be owned by `internal/study`; CLI parsing/rendering belongs in `internal/app`; platform runtime must remain generic and must not import or know study semantics.
- Use existing prompt builders, report validators, rating parser, source applicability helper, runtime boundary, and study service patterns instead of duplicating equivalent logic.
- Runtime execution must go through the existing `internal/platform/runtime` abstraction backed by `github.com/Antonio7098/agentwrap` and `agentwrap/opencode`; UltraPlan product code must not invoke OpenCode directly.
- Worker concurrency must be bounded by explicit configuration or flag values and must not use unbounded goroutine creation, unbounded channels, or best-effort background task leaks.
- Every started runtime run must be waited on, and its event stream must be drained or safely consumed according to the platform runtime contract.
- Filesystem writes created by UltraPlan, including `summary.csv` or any updated state artifacts, must use explicit workspace-contained paths and must not write outside the workspace or selected study.
- Runtime and validation errors must preserve cause chains internally while rendering only safe, redacted diagnostics to users.
- Markdown document sources remain document-only; inapplicable Markdown pairs are skipped, not failed; Markdown reports do not require code citations by default.
- No global technical-layer packages may be introduced for batch behavior unless a concrete cross-module reuse need is documented and accepted in reasoning.
- Default verification must remain deterministic and offline, using fake runtimes and local fixtures.

## Dependencies

| Prior Sprint / Output | Required For | Notes |
| ------------------------ | ------------------- | --------- |
| Sprint 5: Markdown source discovery and applicability | Correct task matrix filtering | `run-all` must reuse directory-vs-Markdown source discovery and `GetApplicableSources` behavior. |
| Sprint 6: Report validation and rating parsing | Product success gates and summary ratings | Runtime success must be followed by per-source/final report validation; summary uses strict rating semantics. |
| Sprint 7: Prompt composition | Analysis and synthesis prompt generation | `run-all` must reuse deterministic analysis and synthesis prompt builders and manifests. |
| Sprint 8: Run-state persistence and status primitives | Optional state/status integration | Existing versioned state primitives may be reused, but full durable run-loop behavior remains Sprint 12 scope. |
| Sprint 9: Agentwrap/OpenCode runtime integration | Runtime execution boundary | Batch tasks must use the generic platform runtime and agentwrap-backed request/result/event mapping. |
| Sprint 10: Single analysis and synthesis | Per-task execution semantics | `run-all` must compose existing single analysis and synthesis behavior rather than inventing separate task execution paths. |
| Project docs: PRD, TRD, ARCHITECTURE | Scope and architectural rules | Study-side-only scope, module ownership, runtime adapter constraints, validation, concurrency, security, and testing rules remain binding. |

## Review Expectations

| What | How Verified |
| ------------- | ----------------------- |
| Required sprint artifacts exist and match the standard artifact chain | Inspect `projects/ultraplan-go/sprints/11-run-all-batch-execution/` for `requirements.md`, `sprint-index.md`, `technical-handbook.md`, `reasoning/`, `reasoning.md`, `plan.md`, and `review.md`. |
| `run-all` CLI surface is implemented | Run or test `ultraplan study <study> run-all --help` and command-level tests for success, usage failures, filter failures, runtime failures, validation failures, and partial completion. |
| Applicable task matrix and filters are correct | Unit tests cover directory sources, Markdown sources with no filters, Markdown sources with matching filters, Markdown sources with non-matching filters, source filters, dimension filters, and ambiguous references. |
| Bounded parallelism is enforced | Fake-runtime tests block workers and assert the maximum simultaneous started tasks never exceeds configured parallelism. |
| Validation gates product success | Fake-runtime tests cover runtime success with missing/invalid reports, valid per-source reports, valid final reports, and failed final validation. |
| Synthesis gating is correct | Tests verify synthesis starts only when all applicable selected per-source reports for a dimension are valid, and does not start after failed or missing applicable reports. |
| Summary generation is deterministic and scope-correct | Summary tests compare exact CSV output, repeated-run stability, total sorting, missing rating handling, ambiguous rating warnings, and inapplicable Markdown semantics. |
| Errors and diagnostics are safe | Tests or review confirm prompt content, embedded Markdown bodies, secrets, sensitive environment values, and unsafe native payloads are not printed. |
| Architecture boundaries are preserved | Inspect imports and changed paths; verify `internal/platform/runtime` has no product imports and batch behavior lives in `internal/study` with thin `internal/app` wiring. |
| No Sprint 12 or future-horizon scope leaked in | Search/review implementation for absence of `run-loop`, durable retry scheduling, stale task recovery, lock files, target commands, sprint commands, generic workflow engines, and code extraction. |
| Verification passes | Run `go test ./...` and `go build ./cmd/ultraplan` from `/home/antonioborgerees/coding/ultraplan-go`; record results in `review.md`. |
