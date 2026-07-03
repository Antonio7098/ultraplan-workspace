# Sprint Review: 12 Durable Run Loop

> Project: `ultraplan-go`
> Sprint: `12-durable-run-loop`
> Plan: `plan.md`
> Reasoning: `reasoning.md`
> Review Date: 2026-06-12
> Status: accepted with follow-ups

This review checks conformance after implementation. It does not introduce new scope or architecture.

## Review Inputs

- `projects/ultraplan-go/sprints/12-durable-run-loop/sprint-index.md`
- `projects/ultraplan-go/sprints/12-durable-run-loop/technical-handbook.md`
- `projects/ultraplan-go/sprints/12-durable-run-loop/reasoning.md`
- `projects/ultraplan-go/sprints/12-durable-run-loop/plan.md`
- `projects/ultraplan-go/sprints/12-durable-run-loop/requirements.md`
- Required protocols: `architecture-review-protocol.md` (from `sprint-index.md`)
- Implementation under `/home/antonioborgerees/coding/ultraplan-go/`
- Test, build, race, and architecture-review evidence (see Verification Evidence)

## Decision Conformance

| Decision From `reasoning.md` | Implemented? | Evidence | Notes |
| --- | --- | --- | --- |
| Decision 1: Keep run-loop study-owned with thin CLI wiring | yes | `internal/study/run_loop.go:12-207`, `internal/app/study_commands.go:155-183`; import review confirms no product imports from `internal/platform/runtime` | CLI parses flags, builds config, delegates to `Service.RunLoop`, maps errors, renders output |
| Decision 2: Materialize a deterministic study task state machine | yes | `internal/study/run_state.go:24-99` (deterministic task matrix), `internal/study/run_state.go:158-215` (resume revalidation) | Stable task IDs, typed statuses, Markdown inapplicability excluded, synthesis task state |
| Decision 3: Use atomic run-state persistence with explicit schema handling | yes | `internal/study/state.go:42-93` (atomic save: same-dir temp + flush/close + rename + best-effort dir sync), `internal/study/state.go:95-134` (ValidateRunState) | `state_test.go:100-124` proves rename failure preserves prior valid state |
| Decision 4: Add per-study run-loop locking with conservative stale handling | yes | `internal/study/locks.go:25-76` (acquire with O_EXCL), `internal/study/locks.go:78-96` (ownership-checked release), `internal/study/locks.go:25-58` (force-unlock opt-in only) | `locks_test.go` covers conflict, force-unlock, and release |
| Decision 5: Reconcile resume by recovering stale running tasks and revalidating completed outputs | yes | `internal/study/run_state.go:158-215` (ResumeValidateRunState) | Stale running/validating/waiting/retrying reset; completed revalidated via existing validators |
| Decision 6: Schedule with bounded workers and context-driven cancellation | yes | `internal/study/run_loop.go:104-156` (bounded worker pool, `req.Parallelism` cap, `ctx.Err()` guards, `wg.Wait`) | `go test -race ./...` passes |
| Decision 7: Reuse existing runtime boundary, prompt builders, validators, retry/fallback policy, and metadata | yes | `internal/study/run.go` calls `s.runtime.StartRun`; `internal/study/runtime_metadata.go:10-117` maps agentwrap metadata; `internal/study/run.go:127-131` uses structured error categories | No direct OpenCode invocation from product code; only `internal/platform/runtime/opencode.go` imports agentwrap |
| Decision 8: Make status runtime-free, deterministic, and safe | yes | `internal/app/study_commands.go:606-623` (status loads state + lock only, no Runtime) | `study_status_commands_test.go:13-108` proves runtime-free rendering with redaction |
| Decision 9: Verify with fake runtime, state/lock failure injection, command tests, race tests, build, and architecture review | yes | `go test ./...`, `go test -race ./...`, `go build ./cmd/ultraplan` all pass | Architecture review confirms no forbidden global packages, thin CLI, product-agnostic runtime |

## Contract Conformance

| Contract / Requirement | Satisfied? | Evidence | Deviations |
| --- | --- | --- | --- |
| Architecture | yes | `internal/study` owns orchestration/state/locks/reports/prompts/validation; `internal/app` is a thin command/composition layer; `internal/platform/runtime` has no product imports; no `internal/scheduler`, `internal/workflow`, `internal/prompts`, `internal/reports`, or `internal/validation` packages exist | none |
| Errors | yes | `TaskError.Code` (typed codes), `ErrRunState{Missing,Malformed,Unsupported}` sentinels, `ErrStudyLocked` with owner metadata, cause chains via `%w`, atomic-write loud failure with prior-state preservation, classified exit codes (`ExitValidation`, `ExitRuntime`, `ExitPartial`, `ExitCancel`); redaction symmetric at storage and render layers | none |
| Configuration | yes | `internal/platform/config` centralizes precedence, `Validate` rejects bad values, `ConfigSummary` persisted into `RunState`, `RedactValue` and `sanitizeCommand` redact secrets | none |
| Observability | yes | Typed status enum, structured status output, lock diagnostics, retry-time and next-retry-time surfaced, agent metadata persisted with explicit `*Known` booleans, `MetadataOmission` recorded for raw payloads, stable error codes for alerting | none |
| Security | yes | Atomic temp+rename, O_EXCL lock, ownership-checked release, JSON-with-typed-validation, workspace-contained paths (`RunStatePath`, `RunLoopLockPath` under `study.Path`), redaction at all user-visible surfaces, fail-closed on parallelism<1 and unknown schema | none |
| Testing | partial | Fake-runtime and command tests cover run-loop, locks, atomic write failure injection, status, redaction; race suite passes | No dedicated concurrent-lock race test; no golden/snapshot for `status`/`run-loop` text; no forward-compat schema upgrade fixtures (all carry-forward to Sprint 14 / future) |
| Documentation | partial | `studyRunLoopHelp` and status help documented and tested; `README.md` and `internal/study/doc.go` now describe durable run-loop, locks, and status as implemented; `` reasoning/plan contains full architecture decision context | No operator runbook for `--force-unlock`/lock-diagnostics/recovery flow |
| CLI Surface | yes | `--help` for run-loop and status documented; stdout/stderr separated (data vs redacted diagnostics); exit codes stable (`ExitOK`/`ExitUsage`/`ExitConfig`/`ExitValidation`/`ExitRuntime`/`ExitCancel`/`ExitPartial`); `--force-unlock` is explicit and study-scoped; cancellation wired through `signal.NotifyContext` | `--json` mode for status deferred to Sprint 14 per explicit non-goal |
| LLM Runtime | partial | Runtime seam is the narrow `Service.Runtime` interface; agentwrap/OpenCode is the only adapter (`internal/platform/runtime/opencode.go`); `PersistencePolicy{PersistUnsafeRawPayloads:false}`; typed `SDKError` classification; permission/cleanup/retry metadata persisted | `LLM-PROMPT-001`: `PromptManifest` and `AgentMetadata` lack explicit `prompt_id`/`prompt_version`/`prompt_owner_*`/`prompt_purpose`; effective prompt reference is not surfaced in state or status |
| LLM Evaluation / Cost / Safety | yes | Validation gates (`ValidateSourceReport`/`ValidateFinalReport`), Markdown inapplicability excluded (not failed), usage preserved with explicit `*Known` booleans (no zero coercion), cost estimated and surfaced, permission policy and audit reasons persisted, sanitized lock command, redacted error/policy rendering | none |
| Workflows | partial | Workflow scope justified (multi-step resumable batch); orchestration behind `Service` boundary; explicit `TaskStatus` and `TaskError.Code`; bounded workers with cancellation; reconciliation via resume revalidation; `RunStateSchemaVersion` for versioning | `WF-IDEMPOTENCY-001`: implicit protections (resume revalidation, per-study lock, `Attempts++`) are correct but the idempotency stance is not explicitly documented in reasoning or run_loop.go |
| Performance | yes | Bounded workers (`req.Parallelism` cap, no unbounded goroutines/channels), per-transition atomic writes (trade-off accepted), `ctx.Err()` guards before scheduling, runtime-free status, usage/cost preserved without zero coercion | none |
| Persistence And Migrations | yes | `RunStateSchemaVersion=1` with explicit `ErrRunStateUnsupported` for mismatch; atomic save (same-dir temp, flush/close, rename, best-effort dir sync); prior-state preservation on rename failure (tested); O_EXCL lock with ownership-checked release | none |

## Technical Handbook Conformance

| Pattern | Followed? | Evidence | Notes |
| --- | --- | --- | --- |
| Thin CLI, Study-Owned Orchestration | yes | `internal/study/run_loop.go:12-207`; `internal/app/study_commands.go:155-183` | |
| Explicit Composition Root And Test Seams | yes | `internal/study/service.go:9-40` (`Service` + `WithRuntime` option); fake runtimes in tests | |
| Context-Propagated Cancellation | yes | `ctx` is the first parameter of `RunLoop`/`RunAnalysis`/`Synthesize`; `ctx.Err()` checks in worker loop | |
| Bounded Scheduling With Structured Completion | yes | `internal/study/run_loop.go:104-156` (unbuffered `taskCh`, `req.Parallelism` workers, `wg.Wait`) | |
| Durable State As Product Truth, Not Runtime Truth | yes | `ResumeValidateRunState` resets stale running/validating/waiting/retrying and revalidates completed outputs | |
| Atomic Persistence And Conservative Locking | yes | `internal/study/state.go:42-93`; `internal/study/locks.go:25-76` | |
| Classified Errors And Retry Metadata | yes | `TaskError.Code`, `ErrStudyLocked` sentinel, `runtime_metadata.go` persists attempts/policy/permissions/cleanup/repair | |
| Runtime-Free Status From Persisted Facts | yes | `study.LoadRunState` + `study.LockInfoForStatus` only | |
| Fake Runtime, Fakes, And Golden/Command Tests | yes | `commandFakeRuntime`, `runAllRuntime`, `fakeRuntime`; command-level coverage in `internal/app` | |

| Anti-Pattern | Avoided? | Evidence | Notes |
| --- | --- | --- | --- |
| Do Not Let `RunE` Become The Run Loop | yes | `runStudyRunLoop` delegates to `Service.RunLoop`; no state machine in CLI | |
| Avoid Fire-And-Forget Runtime Goroutines | yes | `internal/platform/runtime/runtime.go:213-262` drains events and calls `Wait` | |
| Do Not Use `context.Background()` In Work Paths | yes | All work paths accept `ctx`; no `context.Background()` in `internal/study` work paths | |
| Do Not String-Parse Runtime Errors | yes | `run.go:127-131` inspects `runtimeResult.Error.Category`; `applyExecutionResult` switches on `res.Status` enum | |
| Do Not Convert Unknown Usage To Zero | yes | `UsageMetadata` has explicit `*Known` booleans; status renders "usage: unknown" when not known | |
| Do Not Leak Secrets, Prompts, Markdown Bodies, Or Unsafe Native Payloads | yes | `sanitizeCommand` redacts token/key/secret; `RedactValue` applied to all user-visible fields; `MetadataOmission` recorded for raw payloads | |
| Do Not Auto-Delete Active-Looking Locks | yes | `AcquireRunLoopLock` only deletes on `force=true` | |
| Avoid Global Mutable Config Or Service Singletons | yes | `Service` is a per-call struct with explicit options; `studyRuntimeFactory` is a single overridable hook | |
| Avoid A Generic Workflow Engine | yes | No `internal/scheduler`/`workflow`/`prompts`/`reports`/`validation` packages exist | |

## Plan Execution

| Plan Task | Status | Evidence |
| --- | --- | --- |
| Task 1: Define Run-Loop Request, Result, And CLI Surface | done | `Service.RunLoop`, `RunLoopRequest`/`RunLoopResult`, `runStudyRunLoop`, `studyRunLoopHelp`, command tests in `study_run_loop_commands_test.go` |
| Task 2: Extend Durable Domain State And Deterministic Task Matrix | done | `domain.go` (typed statuses, `TaskError`, `AgentMetadata`, `LockInfo`), `run_state.go:24-99` (deterministic matrix), `state_test.go:12-196` |
| Task 3: Implement Atomic State Store And Schema Validation | done | `state.go:42-93` (atomic save), `state.go:95-134` (validate), `state_test.go:100-124` (rename-failure injection) |
| Task 4: Add Per-Study Lock Handling | done | `locks.go:25-96`, `locks_test.go:10-43` covers acquire/conflict/force-unlock/release |
| Task 5: Implement Resume Reconciliation And Artifact Revalidation | done | `ResumeValidateRunState` resets stale running/validating/waiting/retrying and revalidates completed outputs |
| Task 6: Implement Bounded Durable Scheduler And Cancellation | done | `run_loop.go:104-156` (bounded workers, `ctx.Err()` guards, `wg.Wait`); race suite passes |
| Task 7: Reuse Runtime, Prompt, Validation, Retry/Fallback, And Metadata Paths | done | `RunAnalysis`/`Synthesize` reused; `runtime_metadata.go` persists attempts/policy/permissions/cleanup/repair/usage/cost; structured error classification |
| Task 8: Update Runtime-Free Status And Safe Output | done | `status` loads state + lock only; renders safe metadata with redaction; `study_status_commands_test.go` proves runtime-free behavior |
| Task 9: Complete Cross-Cutting Test Matrix | done | Run-loop, lock, run-loop command, status metadata, platform runtime metadata tests; full and race suites pass |
| Task 10: Run Verification And Prepare Review Evidence | done | `go test ./...`, `go test -race ./...`, `go build ./cmd/ultraplan` all pass; architecture and scope reviews recorded |

## Verification Evidence

| Check | Command / Method | Result | Notes |
| --- | --- | --- | --- |
| Full test suite | `go test ./...` from `/home/antonioborgerees/coding/ultraplan-go` | pass | All packages green; no OpenCode/credentials/network required |
| Race-sensitive paths | `go test -race ./...` from `/home/antonioborgerees/coding/ultraplan-go` | pass | `internal/study` (6.001s) and `internal/app` (4.465s) green under `-race` |
| CLI build | `go build ./cmd/ultraplan` from `/home/antonioborgerees/coding/ultraplan-go` | pass | Binary builds without errors |
| Atomic write failure preservation | `internal/study/state_test.go:100-124` (`TestSaveRunStateRenameFailurePreservesPriorState`) | pass | Rename-failure injection preserves prior valid state |
| Lock conflict refusal | `internal/study/locks_test.go:17-19`; `study_run_loop_commands_test.go:84-91` | pass | Second `AcquireRunLoopLock` returns `ErrStudyLocked`; CLI exits `ExitPartial` |
| Force-unlock scope | `internal/study/locks_test.go:21-28`; `study_run_loop_commands_test.go:93-102` | pass | `--force-unlock` is explicit and study-scoped |
| Resume revalidation | `internal/study/run_loop_test.go:51-72` | pass | Missing completed report revalidated to failed |
| Status is runtime-free | `internal/app/study_status_commands_test.go:13-108` | pass | Status loads state and lock only; no Runtime invocation |
| Redaction | `internal/study/locks_test.go:13-15` (`--api-key=secret`); `study_run_loop_commands_test.go:115-117` | pass | `sanitizeCommand` and `RedactValue` cover secrets and policy reasons |
| Architecture boundaries | Manual import review | pass | No product imports from `internal/platform/runtime`; no global `internal/scheduler`/`workflow`/`prompts`/`reports`/`validation` packages; no direct OpenCode invocation from product code |

## Deviations

| Deviation | Reason | Risk | Follow-Up |
| --- | --- | --- | --- |
| Sprint index names `system/protocols/sprint-review-protocol.md`; canonical protocol body is `system/protocols/review-sprint-protocol.md`. | Artifact-naming mismatch in project metadata. | Resolved by adding `system/protocols/sprint-review-protocol.md` as a compatibility alias. | None required. |
| No dedicated operator runbook for `--force-unlock`, lock diagnostics, and recovery flow. | README now lists implemented behavior, but detailed recovery guidance remains separate from core command documentation. | Operators have command help and status diagnostics but not a longer runbook. | Add an operator runbook section before release. (Carry-forward.) |
| `PromptManifest` and `AgentMetadata` do not record `prompt_id`/`prompt_version`/`prompt_owner_*`/`prompt_purpose`. | Pre-existing gap inherited from earlier sprints; not in Sprint 12 selected scope. | LLM Runtime contract (`LLM-PROMPT-001`) requires explicit prompt version propagation; operators cannot correlate a run-loop outcome to a concrete prompt revision. | Track as a Sprint 14 follow-up alongside deferred JSON-stability work. (Carry-forward.) |
| No dedicated concurrent-lock race test (two goroutines racing on `AcquireRunLoopLock` under `-race`). | Current lock-conflict test is sequential. | Weakens evidence for `TEST-DET-001`/`TEST-FAIL-001` under contention. | Add a concurrent-acquire test under `-race` in a future sprint. (Carry-forward.) |
| No golden/snapshot tests for `status` or `run-loop` text/JSON output. | Stable public JSON schemas are explicitly deferred to Sprint 14; golden fixtures would need a deliberate update workflow. | Brittle if human output is later stabilized; intentional omission. | Add when Sprint 14 starts. (Carry-forward.) |
| No forward-compatibility / downgrade test fixtures for `RunStateSchemaVersion`. | Schema is v1; no upgrade path yet. | A schema evolution has no automated guard. | Add when the schema evolves. (Carry-forward.) |
| `study_run_loop_commands_test.go` does not cover a mid-run user cancellation + subsequent resume scenario. | Cancellation mapping is covered, but the cancel-mid-run + resume combination is not. | Acceptable evidence gap for current scope. | Add a scenario test in a future sprint. (Carry-forward.) |

## Review Protocol Results

| Protocol | Applied? | Findings |
| --- | --- | --- |
| Architecture Review Protocol | yes | Good fit. Main workflow is visible in `internal/study/run_loop.go`; `internal/app/study_commands.go` is a thin adapter; `internal/platform/runtime` has no product imports; no forbidden global `internal/scheduler`/`workflow`/`prompts`/`reports`/`validation` packages; no direct OpenCode invocation from product code. The architecture review protocol's "Approve" decision applies. |
| Sprint Review Protocol (this document) | yes | See Decision Conformance, Contract Conformance, Plan Execution, Verification Evidence, Deviations, and Final Assessment sections above. |

## Final Assessment

- **Status:** accepted with follow-ups
- **Residual Risks:**
  - Force-unlock is intentionally blunt; it is study-scoped but does not verify process liveness. Operators must confirm the holder is genuinely gone before using `--force-unlock`.
  - Rich metadata can leak unsafe data if field selection is later broadened. Current redaction is conservative (substring match for `token`/`key`/`secret`; explicit `RedactValue` and `MetadataOmission`); future expansion of lock command shapes or runtime metadata should preserve redaction discipline.
  - Conservative stale-lock policy means operators may need explicit force-unlock after real process death; no heartbeat/process-owner verification exists.
  - Stable public JSON schemas are deferred; users may begin depending on current state/status JSON shape.
- **Required Follow-Ups:**
  - None.
- **Optional Follow-Ups (Sprint 14 or later):**
  - Track `LLM-PROMPT-001` prompt-versioning (extend `PromptManifest` with `prompt_id`/`prompt_version`/`prompt_owner_*`/`prompt_purpose`; surface in `AgentMetadata` and status).
  - Add concurrent-lock race test, mid-run cancellation + resume scenario test, and golden/snapshot fixtures for `status`/`run-loop` text output.
  - Add forward-compatibility / downgrade test fixtures when `RunStateSchemaVersion` evolves.
  - Add operator runbook section (README or new RUNBOOK) covering `--force-unlock` semantics, lock diagnostics, stale-running recovery interpretation, and loud-write-failure behavior.
  - Consider routing run-loop diagnostic output through the existing platform logger (`internal/platform/logging`) with structured fields (`run_id`, `dimension`, `source`).
