# Sprint Requirements: 08 Run State Persistence

> Project: `ultraplan-go`
> Sprint: `08-run-state-persistence`
> Purpose: the authoritative, human-readable sprint contract. All other sprint artifacts must satisfy these requirements.

## Sprint Goal

Implement study-owned durable run-state persistence for planned analysis and synthesis tasks, including atomic load/save behavior and deterministic status summaries, without executing runtime work.

## Required Outputs

| Output | Path | Description |
| --- | --- | --- |
| Sprint requirements | `projects/ultraplan-go/sprints/08-run-state-persistence/requirements.md` | Defines the binding goal, scope, acceptance criteria, constraints, dependencies, and review expectations for this sprint. |
| Sprint index | `projects/ultraplan-go/sprints/08-run-state-persistence/sprint-index.md` | Selects the roadmap entries, source docs, evidence reports, contracts, and review protocols that govern this sprint. |
| Technical handbook | `projects/ultraplan-go/sprints/08-run-state-persistence/technical-handbook.md` | Gives implementation guidance for run-state schema, task planning, atomic persistence, status summaries, and tests. |
| Run state reasoning | `projects/ultraplan-go/sprints/08-run-state-persistence/reasoning/run-state-persistence.md` | Records reasoning for state ownership, schema shape, task status model, atomic write strategy, resume validation boundaries, and status aggregation. |
| Consolidated reasoning | `projects/ultraplan-go/sprints/08-run-state-persistence/reasoning.md` | Summarizes accepted sprint decisions, tradeoffs, risks, and requirement mappings. |
| Implementation plan | `projects/ultraplan-go/sprints/08-run-state-persistence/plan.md` | Provides the ordered task plan and verification checklist for implementation. |
| Sprint review | `projects/ultraplan-go/sprints/08-run-state-persistence/review.md` | Reviews completed implementation against these requirements and records carry-forward decisions. |
| Study run-state domain updates | `/home/antonioborgerees/coding/ultraplan-go/internal/study/domain.go` | Defines or extends run state, task state, task kind, task status, filters, config summary, attempt/error metadata, validation summary, and status summary concepts. |
| Study run-state persistence implementation | `/home/antonioborgerees/coding/ultraplan-go/internal/study/state.go` | Implements deterministic run-state construction, schema validation, JSON load/save, atomic writes, and completed-task revalidation hooks. |
| Study run-state planning helpers | `/home/antonioborgerees/coding/ultraplan-go/internal/study/run_state.go` | Builds analysis and synthesis task records from applicable source/dimension pairs without invoking runtime execution. |
| Study status CLI wiring | `/home/antonioborgerees/coding/ultraplan-go/internal/app/study_status_commands.go` | Adds or extends runtime-free `ultraplan study <study> status` behavior that reads persisted state and renders deterministic human output. |
| Run-state persistence tests | `/home/antonioborgerees/coding/ultraplan-go/internal/study/state_test.go` | Tests state construction, schema validation, atomic save/load, corrupted or missing state handling, applicability filtering, and status summaries. |
| Study status command tests | `/home/antonioborgerees/coding/ultraplan-go/internal/app/study_status_commands_test.go` | Tests status command output, missing-state diagnostics, invalid-state diagnostics, and deterministic non-runtime behavior. |

## Acceptance Criteria

- [ ] Run state is persisted at `studies/<study>/.ultraplan/run-state.json` and includes schema version, run ID, created time, updated time, filters, config summary, task list, and completion flag.
- [ ] Analysis task records include stable task ID, task kind, study, dimension number, dimension slug, source name, source kind, output path, status, attempt count, timestamps, last error, next retry time, validation summary, and agentwrap metadata placeholders when available later.
- [ ] Synthesis task records include stable task ID, task kind, study, dimension number, dimension slug, output path, status, attempt count, timestamps, last error, next retry time, validation summary, and dependency metadata for required applicable source reports.
- [ ] Initial run-state construction creates analysis tasks only for applicable source/dimension pairs and excludes inapplicable Markdown document pairs.
- [ ] Initial run-state construction creates synthesis task records only for dimensions with at least one applicable source and records the applicable report inputs deterministically.
- [ ] State task ordering is deterministic for identical study inputs.
- [ ] State writes are atomic by writing a temporary file in the same directory, flushing and closing it, then renaming over the prior `run-state.json`.
- [ ] Failed atomic writes preserve the previous valid state file when one exists.
- [ ] State loading validates schema version and reports unsupported or malformed state with actionable diagnostics that name the state file path.
- [ ] State loading distinguishes missing state from malformed state.
- [ ] Resume validation revalidates completed task outputs through existing report validation helpers before trusting completed task status.
- [ ] Resume validation resets stale `running`, `validating`, or other non-terminal active task statuses to a safe retryable or pending state without marking them complete.
- [ ] Resume validation preserves attempt history and safe last-error details.
- [ ] Status summary counts pending, running, validating, completed, failed, cancelled, skipped, waiting, retrying, and total tasks deterministically.
- [ ] Status summary excludes inapplicable Markdown document pairs from missing-output and required-input counts.
- [ ] Human status output is deterministic, concise, and includes state file path, run ID, completion flag, totals, failed count, active count, retry count, and next retry time when present.
- [ ] Status command does not invoke agentwrap, OpenCode, provider credentials, network calls, subprocess runtime execution, or worker pools.
- [ ] Unit tests cover state construction, stable task IDs, applicable-pair filtering, synthesis dependency recording, atomic save/load, missing state, malformed JSON, unsupported schema version, stale active task reset, completed-output revalidation, and status aggregation.
- [ ] Command tests cover `ultraplan study <study> status` for missing state, valid state, invalid state, failed tasks, and no-runtime behavior.
- [ ] `go test ./...` passes in `/home/antonioborgerees/coding/ultraplan-go`.
- [ ] `go build ./cmd/ultraplan` passes in `/home/antonioborgerees/coding/ultraplan-go`.

## Non-Goals

- Agentwrap/OpenCode runtime integration is not included.
- Actual execution of analysis or synthesis tasks is not included.
- `ultraplan study <study> run-all` and `ultraplan study <study> run-loop` worker-pool orchestration are not included beyond state construction helpers needed by later sprints.
- Retry, fallback, repair, backoff execution, permission policy mapping, runtime health gates, event sinks, and durable agentwrap run stores are not included.
- User interrupt handling and active runtime cancellation are not included.
- Per-study lock files and force-unlock behavior are not included.
- Stable public JSON status output is not required unless already present; this sprint only requires deterministic human status output.
- Summary generation and `summary.csv` writing are not included.
- Code reference extraction is not included.
- Prompt composition changes are not included except consuming existing prompt/output path concepts where needed for task metadata.
- Report validation changes are not included except using existing validation helpers during resume validation.
- Target workflows, sprint planning, and sprint execution remain out of scope.

## Constraints

- Run-state behavior must live in `internal/study`; do not create global `internal/state`, `internal/scheduler`, `internal/reports`, or `internal/validation` packages for study-owned behavior.
- CLI parsing and status rendering must stay in `internal/app`; state construction, persistence, validation, and summary logic must stay outside command handlers.
- Preserve the dependency direction from `docs/ARCHITECTURE.md`: `study` may use `workspace` and platform helpers, but platform packages must not import `study`.
- Because durable state enters scope, apply the roadmap Batch / Durable Workflow Gate requirements for atomic writes, explicit terminal task states, durable task state, diagnostics for current or recent workflow state, and scenario tests; runtime/provider requirements remain deferred.
- Persisted state must include a schema version and reject unsupported schema versions rather than silently accepting incompatible state.
- Persisted paths should be workspace-relative where practical and must not contain secrets or provider credentials.
- State writes must be safe for local filesystem review in Git and must not use an opaque database.
- Task IDs must be stable across process runs for the same study, dimension, source, and task kind.
- Inapplicable Markdown document source/dimension pairs must remain skipped or absent from required task lists, not failed.
- Status calculations must be derived from persisted state and existing study/report metadata, not from runtime process inspection.
- Tests must use local fixtures or temporary directories and must not require OpenCode, provider credentials, network access, real source repositories, or long-running subprocesses.
- Errors must preserve cause chains for filesystem, JSON, and validation failures and must identify the failing state, report, dimension, source, or output path.

## Dependencies

| Prior Sprint / Output | Required For | Notes |
| --- | --- | --- |
| Sprint 1: Go Module, CLI Shell, and App Composition | Buildable CLI, package structure, and command wiring | Required for status command wiring, tests, and `cmd/ultraplan` build continuity. |
| Sprint 2: Workspace, Config, Logging, and Health Skeleton | Workspace discovery, path safety, and config summaries | Required for resolving `studies/<study>/.ultraplan/run-state.json` and recording safe config metadata. |
| Sprint 3: Study Domain, Listing, and Resolution | Study, source, dimension, and lookup model | Required for deterministic task construction and status command study resolution. |
| Sprint 4: Study Initialization From YAML | Generated study directory layout | Required for expected `dimensions/`, `sources/`, `reports/source/`, `reports/final/`, and study-local `.ultraplan/` paths. |
| Sprint 5: Markdown Document Sources and Applicability | Markdown source kind and applicability filtering | Required so state construction, resume validation, synthesis dependencies, and status summaries exclude inapplicable pairs. |
| Sprint 6: Report Validation and Rating Parsing | Report path helpers and validation APIs | Required for completed-task revalidation and deterministic output path metadata. |
| Sprint 7: Prompt Composition Generation | Deterministic prompt/output path concepts and synthesis manifests | Required for task metadata consistency with later runtime execution and synthesis planning. |
| `projects/ultraplan-go/docs/PRD.md` | Product run-state behavior | Governs resumability, explicit task state, status inspection, runtime-success-is-not-product-success, and durable artifacts. |
| `projects/ultraplan-go/docs/TRD.md` | Technical run-state behavior | Governs task state fields, state file path, atomic writes, resume logic, scheduling constraints, and persistence requirements. |
| `projects/ultraplan-go/docs/ARCHITECTURE.md` | Package ownership and dependency direction | Governs keeping run-state behavior in `internal/study` and runtime-generic behavior outside product modules. |
| `projects/ultraplan-go/sprints/05-study-listing-inspection/review.md` | Carry-forward decisions | Confirms Markdown source discovery, frontmatter stripping, and applicability behavior are complete. |
| `projects/ultraplan-go/sprints/06-study-run-execution/review.md` | Carry-forward decisions | Confirms report validation/path/rating concepts are complete and runtime execution remained deferred. |
| `projects/ultraplan-go/sprints/07-prompt-composition-generation/review.md` | Carry-forward decisions | Confirms prompt composition and synthesis manifest concepts are complete and runtime execution/run-loop behavior remained deferred. |

## Review Expectations

| What | How Verified |
| --- | --- |
| Scope alignment | Review `requirements.md`, `sprint-index.md`, `technical-handbook.md`, `reasoning.md`, and `plan.md` against durable run-state scope and confirm no runtime execution, worker pool, retry/fallback execution, cancellation, code extraction, target, or sprint-execution scope was added. |
| Required outputs exist | Check each path listed in Required Outputs and confirm implementation files are present or intentionally replaced by equivalent files recorded in `review.md`. |
| State schema behavior | Run or inspect tests proving `run-state.json` includes schema version, run ID, timestamps, filters, config summary, task records, completion flag, and stable task IDs. |
| Atomic persistence | Run or inspect tests proving state saves use same-directory temporary files plus rename and preserve prior valid state on failed writes where feasible. |
| Applicability behavior | Run or inspect tests proving state construction and status summaries include applicable directory and Markdown pairs while excluding inapplicable Markdown pairs. |
| Resume validation behavior | Run or inspect tests proving completed tasks are revalidated, stale active tasks are reset safely, unsupported schemas fail clearly, and attempt/error metadata is preserved. |
| Status behavior | Run or inspect command tests proving `ultraplan study <study> status` renders deterministic human output from persisted state and handles missing or malformed state actionably. |
| Architecture conformance | Review implementation paths and imports to confirm study-owned state behavior remains in `internal/study`, CLI rendering remains in `internal/app`, platform packages do not import product modules, and no global state/scheduler package was introduced. |
| Verification commands | Run `go test ./...` and `go build ./cmd/ultraplan` from `/home/antonioborgerees/coding/ultraplan-go` and record the results in `review.md`. |
