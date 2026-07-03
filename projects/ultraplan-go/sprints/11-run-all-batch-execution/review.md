# Sprint Review: 11 Run All Batch Execution

> Project: `ultraplan-go`
> Sprint: `11-run-all-batch-execution`
> Plan: `plan.md`
> Reasoning: `reasoning.md`
> Implementation: `/home/antonioborgerees/coding/ultraplan-go/`

This review is the synthesized output of the Review Sprint Protocol. It checks conformance after implementation and does not introduce new scope or architecture.

## Review Inputs

- `projects/ultraplan-go/sprints/11-run-all-batch-execution/sprint-index.md`
- `projects/ultraplan-go/sprints/11-run-all-batch-execution/requirements.md`
- `projects/ultraplan-go/sprints/11-run-all-batch-execution/technical-handbook.md`
- `projects/ultraplan-go/sprints/11-run-all-batch-execution/reasoning/run-all-batch-execution.md`
- `projects/ultraplan-go/sprints/11-run-all-batch-execution/reasoning.md`
- `projects/ultraplan-go/sprints/11-run-all-batch-execution/plan.md`
- `system/protocols/architecture-review-protocol.md`
- `system/protocols/review-sprint-protocol.md`
- Implementation under `/home/antonioborgerees/coding/ultraplan-go/`
- 13 contract review subagent reports (Architecture, Errors, Configuration, Observability, Security, Testing, Documentation, CLI Surface, LLM Runtime, LLM Evaluation / Cost / Safety, Workflows, Performance, Persistence and Migrations)
- 1 technical-handbook review subagent report
- Verification: `go test ./...` (pass), `go test -race ./...` (pass), `go build ./cmd/ultraplan` (pass)

## Decision Conformance

| Decision From `reasoning.md` | Implemented? | Evidence | Notes |
| --- | --- | --- | --- |
| Decision 1: Study-owned ephemeral batch command | yes | `internal/study/run_all.go:14` (`Service.RunAll`); `internal/study/domain.go:305-344` (in-memory result types only); `internal/app/study_commands.go:150-174` (thin CLI wiring); grep for `run-loop`, `SaveRunState`, `lock` in run_all.go returns no matches | No durable orchestration substrate added. |
| Decision 2: Deterministic matrix, filters, applicability, and preflight | yes | `internal/study/run_all.go:14-31` (parallelism guard, resolveDimensions/resolveSources, applicability); `internal/study/run_all_test.go:82-99` (TestRunAllPreflightFailuresStartNoRuntime); `internal/app/study_commands.go:227-249` (flag preflight) | Inapplicable Markdown pairs skipped, not failed or missing. |
| Decision 3: Bounded runtime workers with partial completion and cancellation | yes | `internal/study/run_all.go:34-62` (per-request worker pool); `internal/study/run_all_test.go:101-122` (TestRunAllBoundsParallelismAndReportsCancellation); `go test -race ./...` passes | Continued unrelated tasks after ordinary failures; cancellation stops new scheduling. |
| Decision 4: Reuse existing runtime, prompt, workdir, permission, and validation paths | yes | `internal/study/run.go:84-123` (startRuntime + executionMetadata); `internal/study/synthesize.go:8-81` (reuses RunAnalysis flow); validation gated at `run.go:77-81` and `synthesize.go:76-79` | No second validator, no second prompt builder, no direct OpenCode invocation. |
| Decision 5: Post-analysis synthesis gate and deterministic summary | partial | `internal/study/run_all.go:153-164` (synthesisBlockers); `internal/study/summary.go:26-74` (WriteSummary); `internal/study/summary_test.go` | **Deviation:** summary source rows are ordered alphabetically by source name rather than "by total descending with name as tie-breaker" as Decision 5 states. The alphabetical order is deterministic and satisfies AC-15; reasoning.md should be updated to match. |
| Decision 6: Safe deterministic CLI output, result status, and verification | yes | `internal/app/study_commands.go:261-289` (renderRunAllResult); `internal/app/study_commands.go:291-304` (classifyRunAllResult maps RunAllStatus to ExitOK/ExitUsage/ExitConfig/ExitValidation/ExitRuntime/ExitCancel/ExitPartial); `internal/app/study_run_all_commands_test.go:65-79` (redaction assertion) | No prompt, embedded Markdown, secret, env, or raw payload leakage. |

## Contract Conformance

| Contract | Status | Key Evidence | Issues |
| --- | --- | --- | --- |
| Architecture | pass | `internal/study/run_all.go`, `internal/study/summary.go`, `internal/study/synthesize.go` own batch behavior; `internal/app/study_commands.go:150-174` is thin; `internal/platform/runtime` has no product imports; no `internal/scheduler`, `internal/reports`, `internal/summary`, `internal/prompts`, or `internal/validation` package introduced. | None blocking. |
| Errors | pass_with_issues | `internal/study/run_all.go:42-91` captures every task outcome; `internal/app/study_commands.go:291-304` maps statuses to exit codes; `internal/study/run_all.go:201` (safeError) preserves cause chain via `RuntimeErr`; `internal/app/study_run_all_commands_test.go:65-79` verifies redaction. | (1) `RunAllResult`/`ExecutionResult` lack canonical `Code`/`Category` fields for downstream tooling. (2) Summary write failure is downgraded to a warning instead of a distinct loud status. (3) `run_all.go:16` parallelism < 1 uses plain `fmt.Errorf` rather than a typed sentinel, inconsistent with the typed `RefError` used for invalid filters. (4) `safeError` is misnamed: it does not actually redact, only returns `err.Error()`. |
| Configuration | pass | `internal/app/study_commands.go:238-259` (runAllService); `internal/platform/config/config.go:94-117` (precedence); `internal/study/run_all.go:15-17` and `run_all_test.go:82-99` (parallelism preflight) | None. |
| Observability | pass | `internal/study/run_all.go:166-199` (countRunAll/statusRunAll); `internal/app/study_commands.go:261-289` (status + counts + per-task kind/status/dimension/source/category on stderr); `internal/platform/runtime/runtime.go:169-218` drains events per started run. | None. |
| Security | pass | `internal/study/summary.go:22-24` (SummaryPath inside study.Path); `internal/study/run_all.go:34-62` (bounded workers); `internal/app/study_run_all_commands_test.go:65-79` (prompt and embedded document not in output); workspace path containment via `workspace.ResolveInside`. | None. |
| Testing | pass | `internal/study/run_all_test.go` (8 tests), `internal/study/summary_test.go`, `internal/app/study_run_all_commands_test.go`; `go test ./...` and `go test -race ./...` pass offline with fakes. | None blocking. |
| Documentation | pass_with_issues | All 7 sprint artifacts present; CLI help text in `internal/app/study_commands.go:306-317` matches implementation; `TestStudyRunAllCommandHelpSuccessAndSummary` asserts help. | (1) `internal/study/doc.go:1-6` claims scheduling/runs/synthesis/validation/summaries are deferred to later sprints — directly contradicts current behavior. (2) `README.md:5` says "Study execution, synthesis, resumable run loops, and code-reference extraction are described in the product and technical docs but are not fully implemented yet" — outdated for the run-all path. |
| CLI Surface | pass | `internal/app/study_commands.go:46` dispatch; `:150-174` runStudyRunAll; `:227-249` flag preflight; `:261-304` render + classify; stdout/stderr split at `:262-269` and `:274-285`. | None. |
| LLM Runtime | pass | `internal/study/run_all.go:46` delegates to `s.RunAnalysis`/`s.Synthesize`; `internal/study/service.go:23-32` is the abstract Runtime seam; `internal/study/run.go:84-123` attaches executionMetadata; no direct OpenCode parsing. | None. |
| LLM Evaluation / Cost / Safety | pass_with_issues | `internal/app/study_commands.go:261-289` prints only `RuntimeCategory` for non-completed tasks; `internal/platform/runtime/opencode.go:12-52` `PersistUnsafeRawPayloads=false`; `internal/study/run.go:84-123` maps provider/model/permissions/sandbox. | (1) `EVAL-COST-001` partial: per-task `Usage`/`EstimatedCost`/`StartedAt`/`FinishedAt` and aggregate cost/duration are dropped at `run.go:57-62` and `synthesize.go:56-61` because `ExecutionResult`/`RunAllResult` lack those fields. (2) Optional: route run-all per-task `RuntimeError` through `config.RedactValue` for symmetry with single-run. |
| Workflows | pass | `internal/study/run_all.go:34-62` (bounded worker pool); `:42-45`, `:55-58`, `:65-66` (cancellation gates); `:153-164` (synthesisBlockers); `internal/platform/runtime/runtime.go:156-233` (event drain and wait per started run). | None. |
| Performance | pass | `internal/study/run_all.go:36-53` (exactly `req.Parallelism` workers, unbuffered `taskCh`, pre-allocated result slice, `wg.Wait`); no package-level semaphore; `go test -race ./...` clean. | Minor: no documented upper bound on `--parallel`; not a contract violation. |
| Persistence and Migrations | pass | `internal/study/summary.go` writes `summary.csv` through a temp-file, fsync, rename, and directory sync path; `internal/study/run_all.go` does not call `SaveRunState`; grep for `lock`, `migrate`, `run-loop` in run_all.go returns no matches. | Atomic summary write follow-up addressed on 2026-06-02. |

## Handbook Conformance

| Pattern / Anti-Pattern | Status | Evidence | Notes |
| --- | --- | --- | --- |
| Study-owned thin-command delegation | followed | `internal/study/run_all.go:14` vs `internal/app/study_commands.go:150-174` | |
| Factory + composition root with fakeable seam | followed | `internal/study/service.go:23-32` (Runtime interface + WithRuntime); `internal/app/study_commands.go:20` (studyRuntimeFactory) | |
| Explicit config precedence before launch | followed | `internal/app/study_commands.go:227-249`; `internal/study/run_all.go:15-31` | |
| Context-propagated bounded concurrency | followed (with caveat) | `internal/study/run_all.go:34-62`; `run_all_test.go:101-122`; `go test -race` clean | Optional: producer send is unbuffered and relies on a busy worker to drain before observing ctx.Err(); a `select { case taskCh <- idx: case <-ctx.Done(): }` would be more defensive. |
| Runtime success + artifact validation as product gate | followed | `internal/study/run.go:77-81`; `internal/study/synthesize.go:76-79`; `internal/study/run_all.go:153-164` | |
| Injected IO + deterministic output capture | followed | `internal/app/study_commands.go:30-32`; `study_commands.go:262-289` deterministic text | |
| Event consumption + safe observability | followed (with caveat) | `internal/platform/runtime/runtime.go:169-218` drains events per run; `study_commands.go:261-289` only prints `RuntimeCategory`; `study_run_all_commands_test.go:65-79` redaction | Optional: route `RunAllWarning.Message` through `config.RedactValue` (single-run already does this) and rename or repurpose `safeError` (`run_all.go:201`). |
| Fake-runtime + fixture-first testing | followed | `internal/study/run_all_test.go:14-54`; `internal/app/study_run_all_commands_test.go` fakes; `t.TempDir()` | |
| Do not let CLI code become the scheduler | avoided | CLI is flag parse + service call + render + exit map only | |
| Do not spawn unbounded goroutines / leave events unread | avoided | `wg.Wait` + exactly `Parallelism` workers; adapter drains events | |
| Do not treat runtime exit as success | avoided | Validators gate status; tests cover missing/invalid report | |
| Do not bypass config precedence / validate after launch | avoided | Triple-check (flag parser, runAllService, Service.RunAll) | |
| Do not leak prompts/Markdown/secrets/unsafe payloads | avoided (with caveat) | Only `RuntimeCategory` printed; redaction test passes | Optional: redaction of `Warning.Message` and per-task `RuntimeError`. |
| Do not use globals for batch coordination | avoided | All batch state is per-request inside RunAll's stack | |
| Do not introduce plugin/workflow-engine abstractions | avoided | Domain additions are minimal value types; grep clean | |
| Do not make tests depend on real OpenCode/network/subprocesses | avoided | Fakes + t.TempDir; offline tests pass | |

## Applicability Table (selected requirements not directly triggered)

| Contract / Handbook Item | Applicability | Scope Reason |
| --- | --- | --- |
| ERR-RETRY-001 (retry policy) | explicitly_deferred | Sprint 11 non-goal; durable retry deferred to Sprint 12. Existing agentwrap `BasicPolicy` retained. |
| LLM-TOOL-001 (tool execution boundaries) | not_triggered | No new tool surfaces or selection policy added. |
| LLM-PROMPT-001 (prompt versioning) | not_triggered | Sprint 7/10 owns prompt builders; no new prompt introduced. Follow-up to add `prompt_version` is not in Sprint 11 scope. |
| LLM-OBS-001 (telemetry/tracing stack) | not_triggered | Sprint 11 limits output to deterministic CLI text per Decision 6. |
| LLM-MODEL-001 (model/provider changes) | not_triggered | No model, provider, runtime, or fallback change. |
| EVAL-HUMAN-001 (high-impact uncertain outputs) | not_triggered | Local CLI; user reviews summary.csv; no automated high-impact action. |
| OBS-HEALTH-001 (new health surface) | not_triggered | No new health/preflight CLI surface; existing `ultraplan health` unchanged. |
| OBS-DEBUG-001 (new debug/verbosity controls) | not_triggered | No new debug flags; output is fixed and deterministic. |
| CFG-COMPAT-001 (durable config schema) | not_triggered | No new durable config fields; schema version unchanged. |
| CFG-ENV-001 (env-specific behaviour) | not_triggered | No dev/test/prod mode toggles added. |
| CFG-OBS-001 (operator diagnostics surface) | not_triggered | Existing `ultraplan config show` unchanged; no new diagnostics surface added. |
| CLI-JSON-001 (stable public JSON) | explicitly_deferred | Sprint non-goal; no `--json` flag for run-all. |
| TEST-SMOKE-002 (real OpenCode smoke) | explicitly_deferred | Sprint non-goal; default verification must be offline with fakes. |
| TEST-MIGRATION-001 (migration verification) | explicitly_deferred | No schema/config/data migration in scope. |
| PERSIST-SCHEMA-001 (new durable schema) | not_triggered | No new durable artifact beyond in-memory types. |
| PERSIST-MIG-001 (migrations) | explicitly_deferred | No schema/format changes; Sprint 12 scope. |
| PERSIST-FIXTURE-001 (production test-data leak) | not_triggered | Fakes are test-only; no fixture in production code path. |
| PERSIST-RECOVERY-001 (rollback) | explicitly_deferred | No destructive migration; reconciliation via summary + partial-completion status. |
| WF-IDEMPOTENCY-001 (idempotency keys) | not_triggered | No product-owned retry; re-runs overwrite deterministically. |
| WF-VERSION-001 (workflow definition version) | not_applicable | No durable workflow definition introduced. |
| SEC-AUTHN-001 (authentication) | not_triggered | Local single-user CLI; no identity/session. |
| SEC-AUTHZ-001 (authorization) | not_triggered | No multi-tenant object access. |
| PERF-UI-001 (TUI/UI) | not_triggered | CLI text only. |
| PERF-CACHE-001 (caching) | not_triggered | No cache introduced. |

## Plan Execution

| Plan Task | Status | Evidence |
| --- | --- | --- |
| Task 1: Inspect existing study execution seams | done | Architecture review confirmed reuse of `RunAnalysis`, `Synthesize`, validation, prompt builders, and existing `Runtime` seam (`internal/study/service.go:23-32`). |
| Task 2: Add minimal study domain surface | done | `internal/study/domain.go:263-344` adds `RunAllRequest`, `RunAllResult`, `RunAllStatus`, `RunAllCounts`, `RunAllWarning`, `ExecutionRequest`, `ExecutionResult`, `ExecutionStatus`. All in-memory, command-scoped. |
| Task 3: Implement preflight, filters, and deterministic matrix | done | `internal/study/run_all.go:14-31`; tests in `run_all_test.go:56-99` cover directory sources, Markdown sources, and preflight failures. |
| Task 4: Reuse analysis runtime and validation path | done | `run_all.go:46` calls `s.RunAnalysis`; `run.go:77-81` gates on `ValidateSourceReport`. |
| Task 5: Implement bounded workers, partial completion, and cancellation | done | `run_all.go:34-62`; `run_all_test.go:101-122`; `go test -race ./...` passes. |
| Task 6: Implement synthesis eligibility and execution | done | `run_all.go:153-164` (synthesisBlockers); `synthesize.go:8-81` reuses prompt/runtime/validation. |
| Task 7: Implement deterministic summary generation | done (with ordering deviation) | `summary.go:26-74`; `summary_test.go`. Source rows ordered alphabetically by name, not "total descending with name tie-breaker" as reasoning.md Decision 5 states. |
| Task 8: Add thin CLI wiring and deterministic rendering | done | `study_commands.go:46`, `:150-174`, `:261-304`; `study_run_all_commands_test.go`. |
| Task 9: Architecture and scope review | done | Architecture contract review confirms no forbidden packages, no `run-loop`/lock/stale-recovery/durable retry/direct OpenCode additions. |
| Task 10: Run verification and complete review evidence | done | `go test ./...` ok; `go test -race ./...` ok; `go build ./cmd/ultraplan` succeeds; this `review.md` records evidence. |

## Verification Evidence

| Check | Command / Method | Result | Notes |
| --- | --- | --- | --- |
| Full test suite | `go test ./...` from `/home/antonioborgerees/coding/ultraplan-go` | pass | All packages `ok`; tests use fake runtimes and `t.TempDir()` fixtures; no OpenCode, network, or subprocess execution. |
| Race detector | `go test -race ./...` from `/home/antonioborgerees/coding/ultraplan-go` | pass | No data races detected across `internal/app`, `internal/platform/config`, `internal/platform/logging`, `internal/platform/runtime`, `internal/study`, `internal/workspace`. |
| CLI build | `go build ./cmd/ultraplan` from `/home/antonioborgerees/coding/ultraplan-go` | pass | Binary built without errors. |
| Command help | `internal/app/study_run_all_commands_test.go:15-22` | pass | Asserts help text includes `--dimension`, `--source`, and `--parallel`. |
| Preflight blocks runtime | `internal/study/run_all_test.go:82-99` + `internal/app/study_run_all_commands_test.go:41-63` | pass | Asserts `rt.calls == 0` and `fake.calls == 0` for invalid filters and parallelism < 1. |
| Bounded parallelism | `internal/study/run_all_test.go:101-122` | pass | Asserts `rt.maxActive <= 1` with parallelism 1 and blocked fake runtime. |
| Cancellation | `internal/study/run_all_test.go:101-122` | pass | Asserts `Status == Cancelled` and no new starts after cancellation. |
| Redaction | `internal/app/study_run_all_commands_test.go:65-79` | pass | Asserts `Base Prompt` and `Embedded Document` are absent from stdout/stderr. |
| Architecture review | Architecture contract review + grep for forbidden packages | pass | No `internal/scheduler`, `internal/reports`, `internal/summary`, `internal/prompts`, `internal/validation`; no product imports in `internal/platform/runtime`. |

## Deviations

| Deviation | Reason | Risk | Follow-Up |
| --- | --- | --- | --- |
| Summary source rows are ordered alphabetically by source name. | `reasoning.md` Decision 5 was updated on 2026-06-02 to match the implemented and tested alphabetical ordering. | None. | Addressed. |
| `safeError` (`run_all.go:201`) is misnamed: it returns `err.Error()` without redaction. | Current `renderRunAllResult` (`study_commands.go:261-289`) does not print `RuntimeError`, so the helper is not exercised in user output for analysis/synthesis tasks. `Warning.Message` for `WriteSummary` filesystem errors is printed unchanged. | Low: single-run path uses `config.RedactValue` for the analogous case; run-all's redaction is currently more conservative (only `RuntimeCategory`), but asymmetry is a future regression risk. | Rename `safeError` to `errMessage` and route any future `RuntimeError`/`Warning.Message` rendering through `config.RedactValue` (or a study-local redaction helper) to match the single-run path. |
| `internal/study/doc.go` and `README.md` were stale after Sprint 11. | Both files were updated on 2026-06-02 to reflect implemented run-all, synthesis, validation, summary, and current non-goals. | None. | Addressed. |
| `summary.csv` write needed the package's atomic-write pattern. | `internal/study/summary.go` now writes through a temp-file, sync, rename, and directory sync helper; `summary_test.go` covers failed writes preserving existing summaries. | Low residual risk only. | Addressed. |

## Issue Classification

### Required Fixes Before Acceptance

None. All identified issues are non-blocking for Sprint 11's stated goal.

### Required Follow-Ups (record in next sprint planning)

1. **Errors contract:** introduce typed sentinel for service-layer `parallelism < 1` and consider adding `Code`/`Category` to `ExecutionResult` so downstream tooling can distinguish `validation_failed`, `runtime_failed`, `preflight_blocked`, and `cancelled` without string parsing.
2. **Errors contract:** treat `WriteSummary` failure as a distinct loud status (e.g., `RunAllStatusSummaryFailed`) mapped to a non-zero exit, or at minimum to `ExitPartial`/`ExitRuntime`, so "analysis succeeded but summary could not be persisted" is loud and distinct.
4. **LLM Evaluation / Cost / Safety contract:** aggregate per-task `Usage`/`EstimatedCost`/`StartedAt`/`FinishedAt` into `ExecutionResult` and `RunAllResult` so the orchestrator surface (counts, future status command) reflects cost and latency drivers without depending on runtime-internal logs.
3. **Optional safety tightening:** rename `safeError` or route `RunAllWarning.Message` through `config.RedactValue` (mirroring the single-run path at `study_commands.go:401`); harden producer-side cancellation with a `select { case taskCh <- idx: case <-ctx.Done(): }` for resilience on very large matrices.

## Review Protocol Results

| Protocol | Applied? | Findings |
| --- | --- | --- |
| Architecture Review Protocol | yes | Architecture contract review applied the Behaviour, Architecture Fit, Simplicity, Cohesion, Coupling, DRY, State, Paradigm, Error Handling, Observability, Testing, Performance, and Maintainability sections to the changed code. All architectural dimensions pass. No speculative abstractions; no new global packages; dependency direction intact; no god objects; coupling acceptable; state explicit and ephemeral. |
| Sprint Review Protocol | yes | All 24 acceptance criteria are evidenced in the Decision Conformance, Contract Conformance, Verification Evidence, and Deviations tables above. Three documents (`doc.go`, `README.md`, `reasoning.md` Decision 5 wording) are out of date with the implementation but do not block acceptance. |

## Final Assessment

- **Status:** `accepted with follow-ups`
- **Sprint goal achieved:** yes — `ultraplan study <study> run-all` is implemented, study-owned, thin-CLI-wired, bounded, partially-completing, cancellable, validable, synthesizing, and summary-generating. All 24 acceptance criteria are evidenced. Sprint 12 durable orchestration is correctly excluded.
- **Residual Risks:**
  - Token/cost/duration metadata is dropped at the orchestrator surface (medium for observability, low for Sprint 11 because operators can still read the underlying runtime logs).
- **Required Follow-Ups:** see "Required Follow-Ups" section above. None block Sprint 11 acceptance. All map cleanly to either a small refactor in a future sprint or to Sprint 12's durable orchestration scope.
- **Highest priority remaining fix:** decide whether summary write failures should become a distinct non-zero run-all status instead of warnings.
