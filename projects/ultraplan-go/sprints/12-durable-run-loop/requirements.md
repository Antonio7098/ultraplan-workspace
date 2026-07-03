# Sprint Requirements: 12 Durable Run Loop

> Project: `ultraplan-go`
> Sprint: `12-durable-run-loop`
> Purpose: the authoritative, human-readable sprint contract. All other sprint artifacts must satisfy these requirements.

## Sprint Goal

Implement `ultraplan study <study> run-loop` as durable, resumable study orchestration with atomic state updates, bounded execution, retry/fallback metadata, cancellation handling, stale-task recovery, and per-study locking.

## Required Outputs

| Output | Path | Description |
| ---------- | -------- | --------------- |
| Requirements | `projects/ultraplan-go/sprints/12-durable-run-loop/requirements.md` | This sprint contract defining scope, outputs, acceptance criteria, constraints, dependencies, and review expectations. |
| Sprint index | `projects/ultraplan-go/sprints/12-durable-run-loop/sprint-index.md` | Selected contracts, evidence reports, docs, protocols, prior decisions, and exclusions for durable run-loop implementation. |
| Technical handbook | `projects/ultraplan-go/sprints/12-durable-run-loop/technical-handbook.md` | Sprint-specific implementation guidance for durable orchestration, locking, state transitions, retry metadata, cancellation, status, and fake-runtime verification. |
| Area reasoning | `projects/ultraplan-go/sprints/12-durable-run-loop/reasoning/durable-run-loop.md` | Detailed reasoning for run-loop state machines, resume policy, stale-running recovery, retry/fallback persistence, lock semantics, and trade-offs. |
| Consolidated reasoning | `projects/ultraplan-go/sprints/12-durable-run-loop/reasoning.md` | Decision log summarizing the selected design and mapping decisions to requirements, docs, contracts, and prior sprint follow-ups. |
| Implementation plan | `projects/ultraplan-go/sprints/12-durable-run-loop/plan.md` | Ordered implementation checklist with verification commands and review evidence expectations. |
| Sprint review | `projects/ultraplan-go/sprints/12-durable-run-loop/review.md` | Completed review recording acceptance evidence, deviations, residual risks, and carry-forward follow-ups. |
| Durable run-loop orchestration | `/home/antonioborgerees/coding/ultraplan-go/internal/study/run_loop.go` | Study-owned run-loop execution for creating/loading state, scheduling pending tasks, handling retries, unblocking synthesis, persisting transitions, and completing runs. |
| Run-loop lock handling | `/home/antonioborgerees/coding/ultraplan-go/internal/study/locks.go` | Per-study lock acquisition, stale-lock detection, release, and force-unlock support for `run-loop`. |
| Run-state domain updates | `/home/antonioborgerees/coding/ultraplan-go/internal/study/domain.go` | Durable task status, run-loop request/result, retry/fallback, cancellation, lock, runtime metadata, cost/usage, and diagnostic fields needed by the state machine. |
| Run-state construction and resume updates | `/home/antonioborgerees/coding/ultraplan-go/internal/study/run_state.go` | Deterministic initial task state, resume validation, stale-running reconciliation, completed-output revalidation, and status summary updates. |
| Run-state persistence updates | `/home/antonioborgerees/coding/ultraplan-go/internal/study/state.go` | Atomic load/save behavior, schema compatibility or rejection behavior, and persistence hooks needed by durable run-loop execution. |
| Study service wiring | `/home/antonioborgerees/coding/ultraplan-go/internal/study/service.go` | Public study service method and dependencies required to invoke `run-loop` without leaking CLI concerns into study logic. |
| CLI command wiring | `/home/antonioborgerees/coding/ultraplan-go/internal/app/study_commands.go` | Thin CLI handling for `ultraplan study <study> run-loop`, lock-related flags, cancellation wiring, status rendering updates, and documented help text. |
| Run-loop tests | `/home/antonioborgerees/coding/ultraplan-go/internal/study/run_loop_test.go` | Fake-runtime tests for initial state creation, resume, bounded scheduling, retries, fallback metadata, synthesis unblocking, cancellation, stale recovery, and terminal states. |
| Lock tests | `/home/antonioborgerees/coding/ultraplan-go/internal/study/locks_test.go` | Tests for lock creation, concurrent-lock refusal, conservative stale-lock behavior, force-unlock behavior, and lock cleanup on terminal exit. |
| State tests | `/home/antonioborgerees/coding/ultraplan-go/internal/study/state_test.go` | Updated tests for atomic writes, schema compatibility or rejection, resume persistence, completed-output revalidation, and prior-state preservation on write failure. |
| CLI run-loop tests | `/home/antonioborgerees/coding/ultraplan-go/internal/app/study_run_loop_commands_test.go` | Command-level tests for arguments, flags, help, exit codes, cancellation, lock failures, resume behavior, and fake-runtime outcomes. |
| CLI status tests | `/home/antonioborgerees/coding/ultraplan-go/internal/app/study_status_commands_test.go` | Updated command tests showing durable run-loop status, retry times, active/recent tasks, failures, cancelled tasks, and safe diagnostics without runtime calls. |

## Acceptance Criteria

- [ ] `ultraplan study <study> run-loop` is available in study help and command help with documented flags for filters, parallelism, retry/cancellation behavior, and lock override or force-unlock behavior.
- [ ] Starting `run-loop` with no existing state creates `studies/<study>/.ultraplan/run-state.json` from the deterministic applicable task matrix, preserving source/dimension filters and excluding inapplicable Markdown document source/dimension pairs.
- [ ] Starting `run-loop` with an existing state loads and validates the state before scheduling work; malformed, unsupported, or unsafe state fails before runtime launch with an actionable diagnostic.
- [ ] If the run-state schema changes, the implementation either migrates compatible prior state or rejects unsupported prior state with a tested, explicit error; it must not silently reinterpret unknown durable formats.
- [ ] Every meaningful task transition is persisted atomically to `run-state.json`, including pending, running, retrying or waiting, validating, completed, failed, cancelled, skipped, and synthesis-unblocked transitions.
- [ ] Atomic state writes use a same-directory temporary file, flush/close, rename, and best-effort directory sync; tests prove a failed write preserves the prior valid state.
- [ ] `run-loop` acquires a per-study lock before mutating run state, refuses a second concurrent `run-loop`, records PID, command, and timestamp in the lock file, and releases the lock on normal terminal exit.
- [ ] Stale lock handling is conservative: the command must not delete an active-looking lock automatically, and any force-unlock path must be explicit, tested, and scoped to the selected study.
- [ ] Stale running tasks from a prior interrupted process are recovered on resume by resetting them to a safe schedulable or terminal state according to the selected policy, preserving attempt/error history and recording the recovery decision in state/status.
- [ ] Completed analysis and synthesis tasks are revalidated on resume before being trusted; missing or invalid artifacts move the affected task back to a non-completed state or a failed state with a clear validation diagnostic.
- [ ] Analysis tasks reuse existing single-run prompt composition, runtime request mapping, workdir rules, permission policy mapping, runtime event handling, and per-source report validation.
- [ ] Synthesis tasks are scheduled only when all applicable selected per-source reports for a dimension exist and pass validation; synthesis reuses existing synthesis prompt composition and final-report validation.
- [ ] The scheduler uses bounded workers and never starts more runtime tasks concurrently than the resolved configuration or CLI override allows.
- [ ] Cancellation through context or user interrupt stops new scheduling, cancels active runtime runs through the platform runtime boundary, waits for owned runs to finish cleanup, saves state, records cancelled or retryable active tasks, and exits with the existing cancellation exit-code convention.
- [ ] Every started runtime run has its event stream drained or safely consumed and has `Wait` called so workers cannot deadlock and agentwrap cleanup metadata is observed.
- [ ] Retry behavior is bounded by resolved configuration and agentwrap policy metadata; attempt count, retry/fallback decisions, next retry time when known, runtime/provider/model target, and classified last error are persisted per task.
- [ ] Rate-limit, retry, fallback, validation, permission, cleanup, usage, and cost metadata exposed by the platform runtime are preserved in durable task state where safe; unknown usage values remain unknown and are not converted to zero.
- [ ] Runtime errors are classified with typed errors or agentwrap `SDKError` fields, not string parsing; safe user diagnostics preserve enough category, operation, provider/model, and retry information for operators.
- [ ] `ultraplan study <study> status` remains runtime-free and reports durable run-loop progress, active tasks, failed tasks, retry times, cancelled tasks, stale-recovery decisions, state path, and lock-relevant diagnostics.
- [ ] Human output is concise and deterministic, includes task totals and artifact paths, and does not print prompts, embedded Markdown source bodies, API keys, sensitive environment values, or unsafe native runtime payload bytes.
- [ ] Write failures for required durable artifacts are loud: `run-loop` must not report success when state persistence, lock acquisition/release, or required output validation fails.
- [ ] The implementation carries forward relevant Sprint 11 review follow-ups by adding machine-usable task categories/codes or equivalent typed statuses for durable task outcomes, safe redaction for rendered warnings/errors, and defensive cancellation while feeding worker queues.
- [ ] `run-loop` does not add target scaffolding, sprint planning, sprint execution, code-reference extraction, a generic workflow engine, or a new runtime adapter.
- [ ] Normal tests use fake runtimes and local fixtures only; `go test ./...` passes without OpenCode, provider credentials, network access, or real subprocess execution.
- [ ] Race-sensitive paths are verified with `go test -race ./...`.
- [ ] The CLI builds with `go build ./cmd/ultraplan`.
- [ ] Architecture review confirms study-owned run-loop behavior, thin CLI wiring, no product imports from `internal/platform/runtime`, no direct OpenCode invocation from product code, and no global `scheduler`, `workflow`, `reports`, `prompts`, or `validation` packages.

## Non-Goals

- Implementing code-reference extraction or `ultraplan code <report>...`.
- Implementing new summary-generation semantics beyond invoking or preserving already-implemented study summary behavior where necessary for existing workflows.
- Implementing stable public JSON schemas for `status`, `run-loop`, validation, or diagnostics; Sprint 14 owns JSON stability.
- Implementing real OpenCode smoke tests in the default test suite.
- Implementing new runtime adapters, direct provider workers, or direct OpenCode stdout/stderr parsing outside `agentwrap/opencode`.
- Implementing a generic workflow engine, DAG authoring system, plugin architecture, watch mode, TUI, hosted service, browser UI, or multi-user collaboration.
- Implementing target scaffolding, sprint planning, sprint execution, or any `ultraplan target ...` or `ultraplan sprint ...` command family.
- Implementing a broad migration framework beyond the explicit compatibility or rejection behavior needed for the current run-state schema.
- Making destructive filesystem or Git mutations outside the selected study's `.ultraplan` state/lock area and expected report artifact paths.

## Constraints

- Durable orchestration must be owned by `internal/study`; CLI parsing/rendering belongs in `internal/app`; `internal/platform/runtime` must remain generic and must not import or know study semantics.
- Runtime execution must go through the existing platform runtime boundary backed by `github.com/Antonio7098/agentwrap` and `agentwrap/opencode`; UltraPlan product code must not invoke OpenCode directly.
- Retry, fallback, health, validation wrapping, permission policy, event mapping, usage, cost, and cleanup metadata must come from agentwrap-facing runtime results and metadata, not from parsing native process text.
- Use existing prompt builders, report validators, rating parser, source applicability helper, run-state primitives, and single-task execution semantics instead of duplicating equivalent logic.
- All long-running execution paths must accept and propagate `context.Context`; no `context.Background` may be introduced in work paths.
- Worker concurrency must be bounded by explicit configuration or flag values and must not use unbounded goroutine creation, unbounded channels, or best-effort background task leaks.
- State and lock paths must remain under the selected study directory, must be workspace-contained, and must reject path traversal or unsafe absolute-path behavior.
- Durable writes must be atomic where required; partial writes must not corrupt the last known good state.
- Errors must preserve cause chains internally while rendering only safe, redacted diagnostics to users.
- Markdown document sources remain document-only; inapplicable Markdown source/dimension pairs are skipped, not failed; Markdown reports do not require code citations by default.
- Default verification must remain deterministic and offline with fake runtimes and local fixtures.

## Dependencies

| Prior Sprint / Output | Required For | Notes |
| ------------------------ | ------------------- | --------- |
| Sprint 5: Markdown source discovery and applicability | Correct task matrix and resume validation | `run-loop` must reuse directory-vs-Markdown source discovery and `GetApplicableSources` behavior. |
| Sprint 6: Report validation and rating parsing | Product success gates | Runtime success must be followed by existing per-source and final report validation. |
| Sprint 7: Prompt composition | Analysis and synthesis prompt generation | `run-loop` must reuse deterministic analysis and synthesis prompt builders and manifests. |
| Sprint 8: Run-state persistence and status primitives | Durable state foundation | Existing versioned state, atomic writes, resume validation, and runtime-free status are the substrate for Sprint 12. |
| Sprint 9: Agentwrap/OpenCode runtime integration | Runtime boundary and metadata | Runtime execution, policy, event, validation, permission, usage, cost, and error metadata must flow through agentwrap-backed platform runtime types. |
| Sprint 10: Single analysis and synthesis | Per-task execution semantics | `run-loop` must compose existing analysis and synthesis behavior rather than create a separate runtime execution path. |
| Sprint 11: `run-all` batch execution | Bounded batch scheduling patterns | `run-loop` should reuse applicable matrix, bounded worker, synthesis gate, cancellation, and safe-output lessons while adding durability. |
| Sprint 11 review follow-ups | Durable status and diagnostics quality | Relevant carry-forwards include typed task outcome categories/codes, safe redaction symmetry, cost/usage/timing propagation, and defensive cancellation in worker feeding. |
| Project docs: PRD, TRD, ARCHITECTURE | Scope and architectural rules | Study-side-only scope, module ownership, runtime adapter constraints, state machine, retry/cancellation, security, and testing rules remain binding. |

## Review Expectations

| What | How Verified |
| ------------- | ----------------------- |
| Required sprint artifacts exist and match the standard artifact chain | Inspect `projects/ultraplan-go/sprints/12-durable-run-loop/` for `requirements.md`, `sprint-index.md`, `technical-handbook.md`, `reasoning/`, `reasoning.md`, `plan.md`, and `review.md`. |
| `run-loop` CLI surface is implemented | Run or test `ultraplan study <study> run-loop --help` and command-level tests for success, usage failures, lock failures, resume failures, cancellation, runtime failures, validation failures, and partial completion. |
| Durable state transitions are correct | Unit tests inspect exact `run-state.json` task statuses, attempts, errors, retry metadata, timestamps, recovery decisions, and completion flags after controlled fake-runtime scenarios. |
| Atomic persistence is safe | Tests inject write/rename failures and confirm the prior valid state remains loadable and no success is reported after required persistence failure. |
| Locking prevents concurrent mutation | Tests acquire a lock, attempt a second run-loop, verify refusal, verify conservative stale-lock behavior, verify force-unlock scope, and verify cleanup after terminal exits. |
| Resume and stale-running recovery are correct | Tests start from persisted running/retrying/completed/failed states and verify completed outputs are revalidated, stale running tasks are reconciled, attempts are preserved, and scheduling resumes only eligible tasks. |
| Retry/fallback metadata is persisted | Fake-runtime tests emit retry, rate-limit, fallback, validation, permission, usage, and cleanup metadata and verify safe durable state/status output. |
| Bounded concurrency and cancellation are enforced | Fake-runtime tests block workers, assert max active runtime calls, cancel context/user interrupt, and verify no new tasks start after cancellation. |
| Status output is truthful and runtime-free | Command tests load persisted state without runtime calls and assert progress counts, active tasks, retry times, failed/cancelled tasks, state path, and safe diagnostics. |
| Errors and diagnostics are safe | Tests or review confirm prompt content, embedded Markdown bodies, secrets, sensitive environment values, and unsafe native payload bytes are not printed or persisted by default. |
| Architecture boundaries are preserved | Inspect imports and changed paths; verify `internal/platform/runtime` has no product imports and run-loop behavior lives in `internal/study` with thin `internal/app` wiring. |
| Future-horizon scope did not leak in | Search/review implementation for absence of target commands, sprint commands, code extraction, generic workflow engines, plugin systems, direct OpenCode execution, and stable JSON schema work. |
| Verification passes | Run `go test ./...`, `go test -race ./...`, and `go build ./cmd/ultraplan` from `/home/antonioborgerees/coding/ultraplan-go`; record results in `review.md`. |
