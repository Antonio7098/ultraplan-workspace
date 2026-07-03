# Sprint Review: 10 Single Analysis and Synthesis

> Project: `ultraplan-go`
> Sprint: `10-single-analysis-synthesis`
> Review Date: `2026-06-01`
> Status: `complete`

## Review Inputs

- `sprint-index.md` - selected contracts, evidence reports, scope, non-goals
- `technical-handbook.md` - patterns, trade-offs, anti-patterns, design pressures
- `reasoning.md` - decisions 1-7, contracts applied, expected evidence
- `reasoning/single-analysis-synthesis.md` - area reasoning
- `plan.md` - tasks, evidence checklist, verification commands
- Architecture Review Protocol and Sprint Review Protocol
- Implementation in `/home/antonioborgerees/coding/ultraplan-go/internal/study/` and `/home/antonioborgerees/coding/ultraplan-go/internal/app/`

---

## Decision Conformance

| Decision From `reasoning.md` | Implemented? | Evidence | Notes |
| --- | --- | --- | --- |
| Decision 1: Study-owned single execution | Yes | `internal/study/run.go`, `internal/study/synthesize.go`, `internal/study/service.go` | Orchestration in `internal/study`; CLI delegates only |
| Decision 2: Existing prompt builders and applicability rules | Yes | `run.go:40` calls `BuildAnalysisPrompt`; `synthesize.go:35` calls `BuildSynthesisPrompt` | No prompt duplication |
| Decision 3: Generic agentwrap-backed runtime request mapping | Yes | `run.go:66-84` builds `runtimepkg.Request` with prompt, workdir, provider/model/timeout, metadata, permissions, health/caps, validation | Platform runtime receives only generic fields |
| Decision 4: Study-owned validation is the product success gate | Yes | `run.go:59-62` validates source report after runtime; `synthesize.go:50-53` validates final report after runtime | Runtime success alone is insufficient |
| Decision 5: Typed outcomes, safe diagnostics, deterministic CLI output | Yes | `domain.go:262-271` defines `ExecutionStatus`; `study_commands.go:195-238` renders to stdout/stderr; `config.RedactValue` used for error redaction | |
| Decision 6: Fake-runtime first verification | Yes | `run_test.go` and `study_run_commands_test.go` use fake runtimes; `go test ./...` passes without OpenCode/credentials | |
| Decision 7: Defer batch, workflow, repair, summary, code extraction, target/sprint scope | Yes | No scheduler, worker pool, run-loop, automatic synthesis, summary, code extraction, target/sprint commands introduced | |

---

## Contract Conformance

| Contract / Requirement | Satisfied? | Evidence | Deviations |
| --- | --- | --- | --- |
| **Architecture** | Pass | grep confirms `internal/platform/runtime` has no product imports; `internal/study` owns execution; CLI handlers are thin | None |
| **Errors** | Pass with minor follow-ups | Runtime failures preserve the original error on `ExecutionResult.RuntimeErr`; `RuntimeError` remains a safe render string; ExecutionResult lacks correlation_id, trace_id, and dedicated task ID fields (ERR-SHAPE-001, ERR-TASK-001) | Non-blocking: task traceability partially served by RuntimeRunID |
| **Configuration** | Pass | `config.Validate()` at load time; `runtimepkg.RequestFromConfig` typed; secrets redacted via `config.RedactValue` | None |
| **Observability** | Pass with issues | OBS-CORR-001: No correlation ID generated at request initiation; OBS-ALERT-001: No event emission infrastructure for high-severity signals | Non-blocking: metadata fields exist but lack unique correlation ID; event emission deferred to future batch scope |
| **Security** | Pass | `run.go:44-47` sets workdir per source kind; `config.RedactValue` for diagnostics; no direct OpenCode handling; `SEC-FILES-001` satisfied with explicit workdir and 0o755 mkdir | None |
| **Testing** | Pass | Fake runtime tests cover success, failure, skip, preflight; `go test ./...` passes without external dependencies; TEST-SMOKE-002 deferred with documented risk | None |
| **CLI Surface** | Pass with issues | All study commands implemented; PRD-7.1-CLI-Behavior-JSON: study status and prompt commands lack `--json` flag | Medium non-blocking: JSON output is a Should-Have per PRD |
| **LLM Runtime** | Pass | Wrapper composition matches TRD 11.2; events drained non-blocking; errors classified via `errors.As`; no string parsing; diagnostics redact | None |
| **LLM Evaluation / Cost / Safety** | Pass with issues | EVAL-COST-001: cost/latency metadata captured by platform runtime but not propagated to ExecutionResult or user output; `Usage.Native` copied unsafely without filtering | Medium non-blocking: no hidden cost multiplication introduced; safe fields could be added in future sprint |
| **Technical Handbook** | Pass | All 8 patterns followed (thin CLI, explicit composition, validation gate, fake runtime, context-aware execution, safe diagnostics, config precedence); all 8 anti-patterns avoided | None |

---

## Plan Execution

| Plan Task | Status | Evidence |
| --- | --- | --- |
| Task 1: Inspect existing APIs and seams | Done | Existing service, prompt builders, validators, runtime interface reviewed |
| Task 2: Define minimal execution domain types | Done | `domain.go:262-294` adds `ExecutionStatus`, `ExecutionRequest`, `ExecutionResult`, `SynthesisRequest` |
| Task 3: Wire runtime dependency through study service and app composition | Done | `service.go:23-31` Runtime interface + `WithRuntime` option; `study_commands.go:179-193` `executionService` wires config |
| Task 4: Implement single analysis orchestration | Done | `run.go:14-64` RunAnalysis with Markdown skip, prompt builder call, runtime request, validation |
| Task 5: Implement single synthesis orchestration | Done | `synthesize.go:7-54` Synthesize with preflight, prompt builder call, runtime request, validation |
| Task 6: Wire thin CLI commands | Done | `study_commands.go:139-177` run/synthesize commands with parse/render/exit mapping |
| Task 7: Add study-level fake-runtime tests | Done | `run_test.go` covers success, failure, skip, preflight paths |
| Task 8: Add CLI execution tests | Done | `study_run_commands_test.go` covers success, skip, runtime failure, validation failure, exit codes |
| Task 9: Architecture and scope review | Done | No scheduler/worker pool/run-loop/batch machinery; platform runtime generic; CLI thin |
| Task 10: Run verification and prepare review evidence | Done | `go test ./...` passed; `go build ./cmd/ultraplan` passed |

---

## Verification Evidence

| Check | Command / Method | Result | Notes |
| --- | --- | --- | --- |
| Full default test suite | `go test ./...` | Pass | No OpenCode, provider credentials, network, wall-clock sleeps, or real subprocess required |
| CLI build | `go build ./cmd/ultraplan` | Pass | Binary builds successfully |
| Architecture review | Code search + import analysis | Pass | `internal/platform/runtime` has no product imports; execution in `internal/study`; CLI thin |
| Runtime boundary review | `grep -r "ultraplan-go/internal/study" internal/platform/runtime/` | Pass (zero matches) | Platform runtime is generic |
| Scope review | `grep -E "run-loop|worker.?pool|scheduler|batch.?loop" internal/study/` | Pass (no matches) | No batch/run-loop/scheduler machinery |
| Diagnostics safety | `study_run_commands_test.go:71-78` asserts no prompt leaking and redaction | Pass | stderr does not contain embedded document or base prompt content |

---

## Deviations

| Deviation | Reason | Risk | Follow-Up |
| --- | --- | --- | --- |
| No correlation ID (OBS-CORR-001) | Single-task execution does not require cross-subsystem correlation; metadata fields exist without unique ID | Low - RuntimeRunID provides some traceability | Consider adding correlation ID for future batch/durable execution |
| No event emission for alerting (OBS-ALERT-001) | Event emission infrastructure is deferred to batch/durable workflow scope per sprint non-goals | Low - exit codes and error categories provide machine-readable signals | Add event emission when batch scope enters |
| No `--json` for study status/prompt (CLI Surface) | JSON output was not selected as required for this sprint; text output is sufficient for current scope | Low - PRD Section 7.1 classifies as Should-Have | Add `--json` flag if machine-readable output becomes required |
| Cost/latency metadata not propagated (EVAL-COST-001) | Platform runtime captures Usage but ExecutionResult does not surface it; `Usage.Native` copied without filtering | Low - no hidden cost multiplication; unknown usage stays unknown | Add safe cost/latency fields to ExecutionResult in future sprint |
| Cause chain lost in error translation (ERR-TRANS-001) | Addressed after review: `ExecutionResult.RuntimeErr` now preserves the original runtime error for analysis and synthesis failures. | Closed | Covered by `internal/study/run_test.go`; CLI output remains unchanged and uses the safe string diagnostic. |
| TEST-SMOKE-002 real runtime smoke deferred | Sprint reasoning explicitly deferred real smoke tests; risk accepted and documented in plan.md | Low - fake runtime tests cover UltraPlan-owned mapping/gating | Add environment-gated smoke tests when real runtime confidence needed |

---

## Review Protocol Results

| Protocol | Applied? | Findings |
| --- | --- | --- |
| Architecture Review Protocol | Yes | All architecture requirements satisfied. Module boundaries explicit, dependency direction correct (platform → study ← app), entrypoint adapters thin, no product imports in platform runtime, no batch/scheduler machinery. Minor issue: ExecutionResult lacks dedicated task ID. |
| Sprint Review Protocol | Yes | All acceptance criteria met. Study run and synthesize commands work with fake runtime. Validation is product success gate after runtime completion. Inapplicable Markdown skips without runtime invocation. Synthesis preflight blocks before runtime launch when reports missing. No new JSON output added. All tests pass, build succeeds. |

---

## Final Assessment

- **Status:** `accepted with follow-ups`
- **Residual Risks:**
  - No correlation ID for cross-subsystem traceability in single execution
  - No event emission infrastructure for high-severity alerting
  - Cost/latency metadata captured but not propagated to user output
  - Real runtime smoke coverage deferred (accepted risk per plan.md)
- **Required Follow-Ups:**
  - None blocking. All non-blocking items are documented with rationale tied to sprint reasoning, requirements, or accepted trade-offs.
  - Future enhancement: Add correlation ID to ExecutionResult for durability and diagnostics
  - Future enhancement: Add safe cost/latency fields to ExecutionResult or CLI output when available
  - Future enhancement: Add `--json` flag to study status and prompt commands if machine-readable output needed
