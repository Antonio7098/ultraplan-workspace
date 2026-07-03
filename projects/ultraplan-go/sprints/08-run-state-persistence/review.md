# Sprint Review: 08 Run State Persistence

> Project: `ultraplan-go`
> Sprint: `08-run-state-persistence`
> Plan: `plan.md`
> Reasoning: `reasoning.md`
> Review Date: 2026-05-31

This review checks conformance after implementation. It does not introduce new scope or architecture.

## Review Inputs

- `sprint-index.md`
- `technical-handbook.md`
- `reasoning.md`
- `plan.md`
- Required protocols from `sprint-index.md`
- Implementation diff / changed files
- Test, build verification evidence
- 10 contract review subagent reports
- 1 technical handbook review subagent report

## Decision Conformance

| Decision From `reasoning.md` | Implemented? | Evidence | Notes |
| --- | --- | --- | --- |
| Decision 1: Keep run state study-owned with thin status wiring | yes | `internal/study/domain.go`, `internal/study/state.go`, `internal/study/run_state.go`, `internal/app/study_commands.go:125-138` | Run-state schema, construction, validation, persistence, resume, and summaries in `internal/study`; CLI handler delegates to study APIs |
| Decision 2: Use a versioned JSON run-state schema with explicit task metadata | yes | `internal/study/domain.go:166-260` (RunState, TaskState types), `internal/study/state.go:93-131` (schema version validation) | `RunStateSchemaVersion=1`, stable task IDs, all required metadata fields |
| Decision 3: Save state atomically with same-directory temp files and preserve prior valid state on failure | yes | `internal/study/state.go:44-91` (temp file, flush, close, rename), `internal/study/state_test.go:100-124` (failure test) | Same-directory temp, Sync, Close, Rename; prior state preserved on failure |
| Decision 4: Build deterministic planned analysis and synthesis tasks from applicable inputs only | yes | `internal/study/run_state.go:24-113`, `internal/study/applicability.go:3-18` | Analysis tasks only for applicable pairs; synthesis tasks only for dimensions with applicable sources |
| Decision 5: Validate loaded state strictly and revalidate completed outputs before trusting completion | yes | `internal/study/run_state.go:157-214` (ResumeValidateRunState), `internal/study/state.go:93-131` (ValidateRunState) | Missing vs malformed vs unsupported; stale active statuses reset; completed outputs revalidated |
| Decision 6: Persist only narrow, secret-free operational configuration summary | yes | `internal/study/domain.go:183-191` (ConfigSummary fields) | Runtime, Model, Variant, counts, workspace version only; no secrets |
| Decision 7: Render deterministic runtime-free human status from persisted state | yes | `internal/app/study_commands.go:125-165` (runStudyStatus, renderStudyStatus), `internal/app/study_status_commands_test.go` | No agentwrap/OpenCode/network/subprocess; concise human output |
| Decision 8: Verify with focused state, persistence, resume, status, and architecture tests | yes | `internal/study/state_test.go`, `internal/app/study_status_commands_test.go`, `go test ./...`, `go build ./cmd/ultraplan` | All acceptance criteria tested; full test suite and build pass |

## Contract Conformance

| Contract / Requirement | Satisfied? | Evidence | Deviations |
| --- | --- | --- | --- |
| Architecture | yes | 6 direct requirements satisfied, 1 partial (ARCH-LAYER-002 targeted seam trade-off), 1 not_triggered | None |
| Errors | yes | 9 direct requirements satisfied, 3 explicitly_deferred/not_triggered | ERR-SHAPE-001 uses simplified error shapes (accepted trade-off); syncDir silences directory sync errors (best-effort) |
| Configuration | yes | 6 direct/partial requirements satisfied, 1 not_triggered, 1 deferred | Migration code deferred as accepted technical debt |
| Observability | yes | 4 direct requirements satisfied, 4 not_triggered | Structured log emission not wired (deferred); resume-validated state not saved to disk (deferred to runtime sprint) |
| Security | yes | 5 direct/partial requirements satisfied, 5 not_triggered | runStateWriteHooks is mutable package-level var (minor, accepted for test seam) |
| Testing | yes | 7 direct/partial requirements satisfied, 1 not_triggered, 2 explicitly_deferred | Runtime smoke tests deferred; migration verification deferred |
| CLI Surface | yes | 6 direct/partial requirements satisfied, 2 explicitly_deferred, 1 not_triggered | Machine-readable JSON output deferred; lifecycle/cancellation deferred |
| Workflows | yes | 2 direct requirements satisfied, 1 partial, 4 explicitly_deferred/not_triggered | Retry, idempotency, compensation, engine boundaries deferred |
| Performance | yes | 2 direct requirements satisfied, 3 not_triggered, 2 explicitly_deferred | Concurrency and runtime/provider cost deferred |
| Persistence And Migrations | yes | 6 direct requirements satisfied, 2 not_triggered | Derived store and rollback strategy not triggered by sprint scope |

## Handbook Conformance

| Pattern / Guidance | Followed? | Evidence | Notes |
| --- | --- | --- | --- |
| Study-owned state near study behavior | followed | `internal/study/` ownership; no global packages | |
| Thin status command, non-CLI state logic | followed | `internal/app/study_commands.go:125-138` delegates to study APIs | |
| Explicit durable state schema with version rejection | followed | `internal/study/state.go:93-131` version validation | |
| Atomic file write seam with injectable filesystem/time | followed | `internal/study/state.go:19-23` (atomicWriteHooks), `internal/study/state.go:44-91` (atomic write), `internal/study/run_state.go:19` (Now field) | |
| Resume validation before trusting completed state | followed | `internal/study/run_state.go:157-214` revalidates outputs | In-memory only; persist deferred to runtime sprint |
| Deterministic summaries from persisted state, not runtime inspection | followed | `internal/study/run_state.go:119-155`; status command has no runtime calls | |
| Actionable, wrapped diagnostics | followed | `internal/study/state.go:13-17` sentinels, `internal/study/state.go:25-91` wrapping with path | |
| Secret-free persisted summaries | followed | `internal/study/domain.go:183-191` ConfigSummary | |
| Table-driven plus command-output tests | followed | `internal/study/state_test.go`, `internal/app/study_status_commands_test.go` | |

**Anti-patterns**: All 8 anti-patterns from the handbook are avoided. No global technical-layer packages, no state logic in `RunE`, no silent acceptance of malformed/unsupported state, no trust without validation, no secrets in config summaries, no runtime inspection in status, no recursive scans, no inapplicable pairs treated as missing/failed.

## Applicability Table

| Contract / Requirement | Applicability | Scope Reason |
| --- | --- | --- |
| Architecture: ARCH-SHARED-001 | not_triggered | No platform or shared packages were added or modified |
| Configuration: CFG-ENV-001 | not_triggered | No environment-specific behavior added |
| Configuration: CFG-COMPAT-001 migration code | explicitly_deferred | Schema version 1 only; migration deferred as accepted debt |
| Errors: ERR-RETRY-001 | explicitly_deferred | Retry/fallback/backoff execution is sprint non-goal |
| Errors: ERR-DEP-001 | not_triggered | No external dependencies called |
| Observability: OBS-HEALTH-001 | not_triggered | No health/preflight changes |
| Observability: OBS-DEBUG-001 | not_triggered | No debug/verbosity changes |
| Observability: OBS-ALERT-001 | not_triggered | No alerting infrastructure |
| Security: SEC-AUTHN-001 | not_triggered | No authentication in scope |
| Security: SEC-AUTHZ-001 | not_triggered | No authorization in scope |
| Security: SEC-INJECT-001 | not_triggered | No injection-prone operations |
| Security: SEC-DEPS-001 | not_triggered | No new dependencies |
| Security: SEC-NET-001 | not_triggered | No network calls |
| Testing: TEST-SMOKE-001/002 | explicitly_deferred | Runtime paths are non-goal |
| Testing: TEST-MIGRATION-001 | partial | Schema version rejection tested; migration code deferred |
| CLI Surface: CLI-JSON-001 | explicitly_deferred | Stable JSON deferred per requirements |
| CLI Surface: CLI-LIFE-001 | explicitly_deferred | Cancellation lifecycle deferred |
| CLI Surface: CLI-SAFE-001 | not_triggered | Status is read-only |
| Workflows: WF-BOUNDARY-001 | not_triggered | No engine boundaries |
| Workflows: WF-RETRY-001 | explicitly_deferred | Retry policy deferred |
| Workflows: WF-IDEMPOTENCY-001 | explicitly_deferred | Idempotency deferred |
| Workflows: WF-COMP-001 | explicitly_deferred | Compensation deferred |
| Performance: PERF-MEASURE-001 | not_triggered | No complex optimizations |
| Performance: PERF-UI-001 | not_triggered | Text-only output |
| Performance: PERF-CACHE-001 | not_triggered | No caching layer |
| Performance: PERF-CONC-001 | explicitly_deferred | Concurrency deferred |
| Performance: PERF-COST-001 | explicitly_deferred | Runtime cost tracking deferred |
| Persistence: PERSIST-DERIVED-001 | not_triggered | No derived stores |
| Persistence: PERSIST-RECOVERY-001 | not_triggered | No destructive data changes |

## Plan Execution

| Plan Task | Status | Evidence |
| --- | --- | --- |
| Task 1: Define Run-State Domain Types | done | `internal/study/domain.go` — RunState, TaskState, TaskKind, TaskStatus, ConfigSummary, ValidationSummary, AgentMetadata, SynthesisDependency, StatusSummary |
| Task 2: Implement Deterministic Run-State Construction | done | `internal/study/run_state.go` — NewRunState, stable task IDs, applicability filtering, synthesis dependencies, deterministic ordering |
| Task 3: Implement State Load, Validation, Save, And Atomic Persistence | done | `internal/study/state.go` — LoadRunState, SaveRunState, ValidateRunState, atomic write with temp+rename |
| Task 4: Implement Resume Validation And Status Summary Logic | done | `internal/study/run_state.go:157-214` — ResumeValidateRunState, SummarizeRunState |
| Task 5: Wire Runtime-Free Study Status Command | done | `internal/app/study_commands.go:125-165` — runStudyStatus, renderStudyStatus; studyHelp updated |
| Task 6: Add State And Persistence Tests | done | `internal/study/state_test.go` — construction, atomic save/load, missing/malformed/unsupported, resume validation, status aggregation |
| Task 7: Add Study Status Command Tests | done | `internal/app/study_status_commands_test.go` — missing/valid/malformed state, diagnostics, exit codes |
| Task 8: Run Verification And Architecture Review Checks | done | `go test ./...` pass, `go build ./cmd/ultraplan` pass, architecture/import review confirms package boundaries |

## Verification Evidence

| Check | Command / Method | Result | Notes |
| --- | --- | --- | --- |
| Full test suite | `go test ./...` | pass | All packages pass without runtime/provider/network dependencies |
| CLI build | `go build ./cmd/ultraplan` | pass | Binary builds successfully |
| Architecture review | Manual import/path review | pass | Run-state in `internal/study`, CLI in `internal/app`, no global packages, no platform imports of `study` |
| Sprint review | 10 contract + 1 handbook subagent reviews | pass | All applicable requirements satisfied |

## Deviations

| Deviation | Reason | Risk | Follow-Up |
| --- | --- | --- | --- |
| ResumeValidateRunState validates in-memory only; state file not updated | Status command is read-only; runtime resume not yet implemented | State file retains stale statuses until a future save operation | Persist validated state when run-resume command is added in a later sprint |
| syncDir silences directory open/sync errors silently | Best-effort after successful rename; platform-specific fsync limitations | Cross-platform durability may be incomplete | Document residual risk; revisit during cross-platform hardening |
| runStateWriteHooks is mutable package-level global var | Used only for test injection of rename failures | Violates strictest reading of global mutable state rules | Extract to an exported option struct passed into SaveRunState |
| Protocol path mismatch: sprint-index references `sprint-review-protocol.md` but file is `review-sprint-protocol.md` | Naming inconsistency in repo | Low — review was performed using existing protocol file | Update sprint-index.md to reference the correct filename |
| Per-study locks and multi-writer safety | Sprint non-goal | Concurrent writers may race despite atomic writes | Add locks when run-loop concurrency enters scope |
| Stable public JSON status output | Sprint non-goal | Machine consumers cannot parse human output | Add JSON status when PRD/TRD select machine-readable output |
| Config summary intentionally narrow | Runtime execution not available to populate redacted runtime metadata | Less diagnostic context for operators | Expand when runtime execution sprint adds safe redacted metadata |

## Review Protocol Results

| Protocol | Applied? | Findings |
| --- | --- | --- |
| Architecture Review (system/protocols/architecture-review-protocol.md) | yes | Behaviour clear and scoped; architecture fit is good; complexity increased but justified; cohesion strong; coupling acceptable (targeted seam, no global FS abstraction); no concerning duplication; state/side effects clear; error handling strong; observability acceptable with deferred logging; testing strong; performance fine. **Verdict: Approve with comments.** |
| Sprint Review (system/protocols/review-sprint-protocol.md) | yes | All selected contracts reviewed; handbook guidance checked; decision conformance complete; plan executed fully. **Verdict: Accepted with follow-ups.** |

## Architecture Review Protocol Checklist

### 1. Behaviour Review
- [x] The main workflow is visible — `runStudyStatus` → LoadRunState → ResumeValidateRunState → SummarizeRunState → renderStudyStatus
- [x] The feature does what the requirement asked — all 21 acceptance criteria met
- [x] The implementation does not include unrelated work — no runtime execution, locks, JSON output, or out-of-scope behavior
- [x] Inputs and outputs are clear — study name in, status text out; state file in, state file out
- [x] Side effects are obvious — state file read/write, temp file create/rename
- [x] Failure paths are handled deliberately — missing/malformed/unsupported/write failure/validation failure

### 2. Architecture Fit
- [x] The change belongs in the files/modules it touched — `internal/study` for state, `internal/app` for CLI
- [x] The current workflow still reads cleanly — delegation chain is clear
- [x] No feature was jammed into an unsuitable abstraction — targeted seams, no broad FS interface
- [x] No large refactor was avoided when design clearly needed one — not applicable
- [x] No large refactor was performed without real need — not applicable

**Decision: Good fit**

### 3. Simplicity and Earned Complexity
- [x] No speculative abstractions
- [x] No unnecessary interfaces/protocols/base classes
- [x] No plugin/factory/registry unless justified
- [x] No generic engine for a single concrete use case
- [x] No framework ceremony added without need
- [x] The simplest honest design was chosen

**Complexity verdict: More complex, justified** — durable persistence, schema validation, atomic writes, and resume validation are necessary complexity for the Batch / Durable Workflow gate.

### 4. Cohesion Review
- [x] Related rules are kept together — state schema, construction, persistence, validation, and summaries all in `internal/study`
- [x] Unrelated behaviours are not forced into one unit — CLI output is separate in `internal/app`
- [x] Functions/classes/modules have clear names and purposes
- [x] No god object/service/module was expanded
- [x] No excessive micro-fragmentation was introduced

**Cohesion issues: None**

### 5. Coupling Review
- [x] No new mutable global state beyond `runStateWriteHooks` and `timeNow` (test seams)
- [x] Runtime configuration/state is passed explicitly where needed
- [x] No external code reaches into internals/private fields
- [x] Public methods/interfaces are used appropriately
- [x] Functions do not accept large objects when only small values are needed
- [x] Parameter objects are cohesive and purposeful
- [x] External dependencies are injected or isolated at boundaries
- [x] Core logic is not tightly coupled to infrastructure

**Coupling verdict: Coupling increased but justified** — `runStateWriteHooks` is mutable global state, scoped to `internal/study/state.go`, used only for test injection, and documented as accepted trade-off.

### 6. DRY and Duplication
- [x] Business rules are not duplicated
- [x] Security/authorization rules are not duplicated
- [x] Validation rules are not duplicated inconsistently
- [x] Similar code shape was not abstracted prematurely
- [x] Shared abstractions do not contain caller-specific branches
- [x] No new boolean flags were added to preserve a bad abstraction

**Duplication verdict: No concerning duplication**

### 7. State and Side Effects
- [x] Durable state changes are clear — `SaveRunState` writes `run-state.json`
- [x] Derived state is not treated as authoritative unless justified — status counts derived from persisted state
- [x] Ephemeral state does not leak into global/shared state
- [x] Mutating operations are named clearly — `SaveRunState`, `ResumeValidateRunState`
- [x] Queries do not unexpectedly mutate state — `LoadRunState`, `SummarizeRunState` are read-only
- [x] Side effects are at boundaries where practical — persistence in `internal/study/state.go`

**State verdict: Clear and safe**

### 8. Error Handling
- [x] Expected failures have specific handling — missing/malformed/unsupported sentinel errors
- [x] Unexpected failures are surfaced clearly — wrapped with `%w` and file path
- [x] Error names/messages are actionable — path context included
- [x] Errors are not swallowed silently — `syncDir` is documented best-effort
- [x] Retry behaviour is safe and deliberate — deferred (sprint non-goal)
- [x] Partial progress is handled — in-memory resume validation preserves attempts/errors

**Error handling verdict: Strong**

### 9. Observability
- [x] Important workflow start/completion/failure is observable — status output, persistence errors
- [x] Logs/events/metrics include correlation identifiers — RunID in state and status output
- [x] Dependency calls are traceable where relevant — no dependency calls in scope
- [x] Duration/performance can be inspected — not yet instrumented (deferred)
- [x] Errors include useful type/context — sentinel errors with path context
- [x] Sensitive data is not logged — ConfigSummary is secret-free
- [ ] AI-specific context is versioned/traced where relevant — deferred (no runtime execution)

**Observability verdict: Acceptable** — Missing structured log emission, which is deferred.

### 10. Testing
- [x] Pure logic has focused unit tests — domain construction, task IDs, filtering, summary
- [x] Workflow/use-case behaviour is tested — save/load, resume validation, status command
- [x] External adapters are tested at the right level — filesystem persistence with temp dirs
- [x] Failure paths are tested — missing/malformed/unsupported/write failure
- [x] Regression cases are covered — stale status reset, completed revalidation
- [x] Tests do not require unnecessary real infrastructure — temp dirs and fixtures only
- [x] Tests are not over-mocked to the point of meaninglessness — targeted seams only

**Testing verdict: Strong**

### 11. Performance
- [x] No premature optimization harmed clarity
- [x] Known scale/latency constraints were considered — single file read, O(n) passes
- [x] Expensive operations are not repeated unnecessarily
- [x] Data loading is not accidentally excessive — one file per load, bounded task lists
- [x] Concurrency/race risks are considered where relevant — deferred (sprint non-goal)
- [ ] Measurement exists or is planned if needed — not added (low complexity paths)

**Performance verdict: Fine**

### 12. Maintainability
- **Next feature impact:** Easier — future runtime execution has stable task schema, atomic persistence, resume validation, and status summaries ready
- **Main risks:** Atomic rename semantics vary by platform; no per-study locks; resume-validated state not saved to disk; `runStateWriteHooks` is mutable global test seam
- **Required changes before merge:** None
- **Optional improvements later:** Wire persistence failures into structured log emission; persist resume-validated state on load; extract `runStateWriteHooks` to explicit option struct; add migration tests for schema version 2+

## Final Assessment

- **Status:** Accepted with follow-ups
- **Residual Risks:**
  - Atomic rename semantics vary by filesystem/platform (mitigated by same-directory temp file, flush/close/rename)
  - Concurrent writers may race despite atomic writes (locks deferred)
  - `syncDir` directory sync is best-effort on some platforms
  - Resume-validated state is not persisted to disk until a future run-resume command
- **Required Follow-Ups:**
  - Persist resume-validated state to disk when runtime resume execution is added
  - Extract `runStateWriteHooks` from package-level var to explicit option struct
  - Wire persistence/validation failure paths into structured log emission when logging infrastructure is integrated
  - Add migration verification tests when schema version 2 is introduced
  - Update sprint-index.md to reference `review-sprint-protocol.md` (existing file) instead of `sprint-review-protocol.md` (nonexistent)
