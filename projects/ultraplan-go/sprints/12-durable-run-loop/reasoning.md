# Sprint Reasoning: 12 Durable Run Loop

> Project: `ultraplan-go`
> Sprint: `12-durable-run-loop`
> Output: `projects/ultraplan-go/sprints/12-durable-run-loop/reasoning.md`
> **Inputs Used:** `projects/ultraplan-go/sprints/12-durable-run-loop/requirements.md`, `projects/ultraplan-go/project-index.md`, `projects/ultraplan-go/docs/ARCHITECTURE.md`, `projects/ultraplan-go/docs/PRD.md`, `projects/ultraplan-go/docs/TRD.md`, `projects/ultraplan-go/sprints/12-durable-run-loop/sprint-index.md`, `projects/ultraplan-go/sprints/12-durable-run-loop/technical-handbook.md`, `templates/sprint-reasoning.md`

This document decides. It synthesizes selected context, handbook evidence, contracts, project docs, and sprint requirements into final sprint decisions.

Requirement IDs in this document use the acceptance-criteria order from `requirements.md` as `AC-01` through `AC-27`. Constraint IDs use the constraint order from `requirements.md` as `CON-01` through `CON-12`.

## Sprint Purpose

- **Goal:** Implement `ultraplan study <study> run-loop` as durable, resumable study orchestration with atomic state updates, bounded execution, retry/fallback metadata, cancellation handling, stale-task recovery, runtime-free status, and per-study locking.
- **Non-Goals:** Do not implement code-reference extraction, target scaffolding, sprint planning, sprint execution, generic workflow engines, new runtime adapters, real OpenCode smoke tests in default tests, stable public JSON schemas, broad migration frameworks, hosted/browser/multi-user features, or destructive filesystem/Git mutations outside selected study state, locks, and expected report artifacts.
- **Depends On:** Sprint 5 source discovery/applicability, Sprint 6 validation/rating parsing, Sprint 7 prompt composition, Sprint 8 run-state/status primitives, Sprint 9 agentwrap/OpenCode runtime integration, Sprint 10 single analysis/synthesis, Sprint 11 bounded `run-all` patterns, and Sprint 11 review follow-ups named in `requirements.md`.

## Selected Context And Pre-Reasoning Artifacts

| Artifact | Path | How It Was Used |
| --- | --- | --- |
| Project Index | `projects/ultraplan-go/project-index.md` | Confirmed project scope, target implementation directory, active contract pool, selected study evidence pool, review protocols, and explicit target/sprint workflow deferral. |
| Product Requirements | `projects/ultraplan-go/docs/PRD.md` | Grounded product principles: resumability by default, runtime success is not product success, explicit state/errors, local-first files, stateful run-loop command, status expectations, and first-release non-goals. |
| Technical Requirements | `projects/ultraplan-go/docs/TRD.md` | Grounded module boundaries, agentwrap/OpenCode runtime composition, task state fields, atomic writes, resume rules, scheduler rules, retry/fallback policy, status diagnostics, security, and fake-runtime testing. |
| Architecture | `projects/ultraplan-go/docs/ARCHITECTURE.md` | Confirmed `internal/study` ownership for study scheduling/state/prompts/validation/reports and `internal/platform/runtime` as generic execution only. |
| Sprint Index | `projects/ultraplan-go/sprints/12-durable-run-loop/sprint-index.md` | Selected the applicable contracts, evidence reports, excluded contexts, review protocols, and non-goals for this sprint. It also selected area reasoning for durable run-loop architecture, but no area reasoning files are present. |
| Technical Handbook | `projects/ultraplan-go/sprints/12-durable-run-loop/technical-handbook.md` | Supplied study evidence, relevant patterns, trade-offs, anti-patterns, concrete source references, open questions, and handoff guidance for durable orchestration. |
| Sprint Reasoning Template | `templates/sprint-reasoning.md` | Provided this document structure and required sections. |

## Area-Specific Reasoning Inputs

None present. `sprint-index.md` selected an Architecture area reasoning output at `projects/ultraplan-go/sprints/12-durable-run-loop/reasoning/durable-run-loop.md`, but no files exist under `projects/ultraplan-go/sprints/12-durable-run-loop/reasoning/*.md` at the time this consolidated reasoning was created. This document therefore makes the final sprint decisions directly from `requirements.md`, project docs, `sprint-index.md`, and `technical-handbook.md`.

## Sprint Technical Handbook Summary

- **Relevant Patterns:** Thin CLI delegation; explicit composition root and dependency seams; context-propagated cancellation; bounded structured scheduling; durable state as product truth; atomic persistence and conservative locking; classified errors and retry metadata; runtime-free status from persisted facts; fake-runtime, fixture, and command tests.
- **Important Trade-Offs:** Persist after every meaningful transition instead of batching; conservative stale-lock handling instead of auto-cleanup; rich metadata persistence instead of minimal status only; completed-output revalidation instead of trusting state; focused files in `internal/study` instead of global scheduler/workflow packages.
- **Warnings / Anti-Patterns:** Do not let CLI `RunE` own the run loop; do not use fire-and-forget runtime goroutines; do not introduce `context.Background()` in work paths; do not string-parse runtime errors; do not convert unknown usage to zero; do not leak secrets, prompts, Markdown bodies, or unsafe native payloads; do not auto-delete active-looking locks; do not create a generic workflow engine.
- **Evidence Confidence:** High for architecture boundaries, state/cancellation, concurrency, error classification, testing, security, and performance because the selected reports cite multiple mature Go CLI examples. Medium for terminal UX and cross-cutting philosophy because those patterns are useful but less prescriptive for the exact durable state machine.

## Contracts Applied

| Contract / Requirement ID | Constraint | Decision Impact | Expected Evidence |
| --- | --- | --- | --- |
| Architecture | Study behavior belongs in `internal/study`; CLI stays thin; platform runtime remains product-agnostic. | Decisions 1, 2, 6, 7 keep orchestration, state machines, synthesis gating, and validation in study code while runtime execution remains generic. | Import review, architecture review protocol, absence of global `scheduler`, `workflow`, `prompts`, `reports`, or `validation` packages. |
| Errors | Durable failures need actionable diagnostics, cause preservation, cancellation classes, lock conflicts, validation failures, and safe rendering. | Decisions 3, 4, 6, 7, 8 require typed/classified task outcomes and loud failures for required writes, locks, validation, and runtime errors. | Unit tests for classifications, CLI tests for exit codes and diagnostics, redaction review. |
| Configuration | Parallelism, retries, timeout, fallback, permission, and runtime values must come from resolved config and CLI overrides. | Decisions 1, 6, 7 make resolved run-loop options explicit and persist safe config summaries. | Command tests for flag precedence and invalid values; state assertions for config summaries. |
| Observability | Status, logs, task metadata, runtime metadata, retry/fallback, cost/usage, and diagnostics must be truthful and safe. | Decisions 2, 7, 8 persist enough facts for runtime-free status and operator diagnostics. | State/status tests for active, failed, retrying, cancelled, recovered, usage, cost, and omitted unsafe payloads. |
| Security | State/lock paths must be workspace-contained; diagnostics must redact secrets and omit unsafe raw payloads. | Decisions 3, 4, 7, 8 constrain paths, lock scope, metadata persistence, and user output. | Path traversal tests, redaction tests, review of output/state fields. |
| Testing | Normal verification must use fakes, fixtures, command tests, and deterministic local inputs. | Decision 9 defines fake-runtime, atomic-write, lock, CLI, status, race, and build evidence. | `go test ./...`, `go test -race ./...`, `go build ./cmd/ultraplan`, targeted test files. |
| Documentation | Help text and recovery guidance must be clear and maintainable. | Decisions 1, 4, 8 require command help for flags, lock override/force unlock, status, and safe output. | CLI help tests and review of rendered command help. |
| CLI Surface | `ultraplan study <study> run-loop` and status must expose flags, help, exit codes, and deterministic human output. | Decisions 1 and 8 define thin command wiring and runtime-free status behavior. | Command-level tests in `internal/app/study_run_loop_commands_test.go` and `internal/app/study_status_commands_test.go`. |
| LLM Runtime | Product code must use agentwrap/OpenCode through runtime boundaries, not direct OpenCode invocation. | Decision 7 requires agentwrap `ObservingRuntime`, `ValidatingRuntime`, `PolicyRunner`, and `opencode` adapter usage through existing platform runtime. | Import review, fake runtime tests, no direct OpenCode process execution in product code. |
| LLM Evaluation / Cost / Safety | Validation, usage/cost, fallback safety, and metadata must be preserved without unsafe data leakage. | Decisions 5, 7, 8 make validation product gates and preserve known usage/cost while leaving unknown usage unknown. | Validation tests, metadata state assertions, safe diagnostics review. |
| Workflows | Stateful batch execution must be resumable, cancellable, retrying, and terminal-state aware. | Decisions 2, 5, 6 define task statuses, transition persistence, resume, recovery, scheduling, and cancellation. | Fake-runtime state machine tests, resume fixtures, cancellation tests. |
| Performance | Workers and channels must be bounded; status should be efficient and runtime-free. | Decisions 6 and 8 cap runtime concurrency and avoid runtime calls during status. | Concurrency tests with blocked fake runtime, status tests without runtime collaborator, race tests. |
| Persistence And Migrations | Durable formats require atomic writes, versioning, compatibility or explicit rejection, and prior-state preservation on failure. | Decision 3 requires same-directory temp writes, flush/close, rename, best-effort directory sync, and explicit schema rejection/migration behavior. | State tests injecting write/rename failures and schema mismatch fixtures. |
| `AC-01` through `AC-27` | Sprint acceptance criteria for command surface, state, locks, resume, validation, scheduling, cancellation, metadata, status, tests, build, and review. | All decisions map directly to these acceptance criteria. | Evidence tables below and sprint review results. |
| `CON-01` through `CON-12` | Sprint constraints for module ownership, runtime boundary, config/metadata source, context propagation, bounded concurrency, path safety, atomic writes, safe errors, Markdown applicability, and offline tests. | All decisions preserve these constraints. | Architecture review, security review, fake-runtime tests, path/state tests. |

## Repos Studied / Source Evidence Used

| Source / Repo / Report | Concrete Reference | Relevant Finding | Why It Matters For This Sprint | Used In Decision(s) |
| --- | --- | --- | --- | --- |
| `01-project-structure` | `studies/go-cli-study/reports/final/01-project-structure.md`; `main.go:16`, `cmd/root.go:112-127`, `cmd/restic/main.go:37-114` | Mature Go CLIs keep entrypoints thin and route behavior into product packages. | Supports study-owned run-loop and thin CLI command wiring. | Decisions 1, 9 |
| `02-command-architecture` | `studies/go-cli-study/reports/final/02-command-architecture.md`; `pkg/cmd/install.go:132`, `cmd/restic/cmd_backup.go:84-115`, `opencode/cmd/root.go:49-183` | Command factories/options structs are testable; long `RunE` handlers are a smell. | Shapes run-loop command options, help, and app/service delegation. | Decisions 1, 8, 9 |
| `03-dependency-injection` | `studies/go-cli-study/reports/final/03-dependency-injection.md`; `pkg/cmdutil/factory.go:16-43`, `executor.go:22-24`, `pkg/action/install.go:162` | Explicit composition roots and constructor seams support tests without global mutable state. | Justifies explicit runtime, store, lock, clock, and IO seams for fake-runtime testing. | Decisions 1, 3, 4, 7, 9 |
| `04-configuration-management` | `studies/go-cli-study/reports/final/04-configuration-management.md`; `internal/cmd/config.go:2253-2287`, `internal/flags/flags.go:314-327`, `internal/global/global.go:139,147` | Strong CLIs track flag precedence and validate merged config. | Supports resolved parallelism, retries, timeout, fallback, permission, and lock override behavior. | Decisions 1, 6, 7 |
| `05-error-handling` | `studies/go-cli-study/reports/final/05-error-handling.md`; `fs/fserrors/error.go:22-29`, `internal/errors/fatal.go:10`, `internal/ghcmd/cmd.go:44-49`, `go-task/errors/errors.go:47-50` | Typed/classified errors are needed for retry, fatal, cancellation, and user-facing exits. | Shapes task outcome categories, retryability, cancellation, lock, validation, and runtime diagnostics. | Decisions 2, 4, 6, 7, 8 |
| `06-io-abstraction` | `studies/go-cli-study/reports/final/06-io-abstraction.md`; `pkg/iostreams/iostreams.go:551-568`, `internal/fs/interface.go:10-31`, `internal/backend/backend.go:19-90` | IO/filesystem seams make CLIs testable and safe. | Supports fake state store, atomic write failure injection, output capture, and command tests. | Decisions 3, 4, 8, 9 |
| `07-state-context` | `studies/go-cli-study/reports/final/07-state-context.md`; `pkg/cmd/install.go:333-347`, `task.go:89`, `cmd/restic/cleanup.go:24-38`, `internal/session/session.go:12-23` | Long-running commands propagate context and model session/lock state explicitly. | Shapes cancellation, active task cleanup, stale recovery, and lock semantics. | Decisions 4, 5, 6 |
| `08-concurrency` | `studies/go-cli-study/reports/final/08-concurrency.md`; `go-task/task.go:87`, `gdu/pkg/analyze/parallel.go:13`, `k9s/internal/pool.go:21,30,37`, `opencode/cmd/root.go:261-279` | Bounded structured concurrency is safer than fire-and-forget goroutines. | Supports bounded workers, drained events, `Wait`, cancellation, and race-sensitive tests. | Decisions 6, 9 |
| `09-terminal-ux` | `studies/go-cli-study/reports/final/09-terminal-ux.md`; `pkg/tasks/tasks.go:144`, `internal/ui/termstatus/status.go:197-205`, `tui/tui.go:316-338` | Long operations should provide concise progress and clean interrupt behavior. | Shapes deterministic run-loop and status output without adding a TUI. | Decisions 1, 8 |
| `10-logging-observability` | `studies/go-cli-study/reports/final/10-logging-observability.md`; `internal/slogs/keys.go:6-231`, `internal/logging/logging.go:31-66`, `pkg/iostreams/iostreams.go:52-54` | Structured observability and safe stdout/stderr separation are core CLI patterns. | Supports persisted task metadata, runtime-free status, and safe diagnostics. | Decisions 7, 8 |
| `11-testing-strategy` | `studies/go-cli-study/reports/final/11-testing-strategy.md`; `internal/cmd/main_test.go:64-174`, `pkg/httpmock/stub.go:35-199`, `internal/backend/mock/backend.go:14-26`, `task_test.go:166-169` | Mature CLIs use table tests, fakes, command helpers, and fixtures. | Drives fake-runtime, lock, atomic-write, CLI, status, cancellation, and race tests. | Decision 9 |
| `13-security` | `studies/go-cli-study/reports/final/13-security.md`; `internal/options/secret_string.go:15-20`, `pkg/registry/transport.go:37-41`, `internal/permission/permission.go:44-108`, `pkg/yqlib/security_prefs.go:3-7` | Explicit trust boundaries, redaction, path safety, and permission posture matter. | Shapes safe lock/state paths, redacted diagnostics, omitted unsafe payload bytes, and permission metadata. | Decisions 3, 4, 7, 8 |
| `14-performance` | `studies/go-cli-study/reports/final/14-performance.md`; `pkg/cmdutil/factory.go:27-42`, `age/internal/stream/stream.go:20`, `yq/pkg/yqlib/stream_evaluator.go:78-113`, `gdu/pkg/analyze/parallel.go:13` | Lazy init, bounded resources, streaming, and efficient status matter for large CLIs. | Supports bounded workers, efficient state/status loading, and avoiding recursive scans beyond needed validation. | Decisions 2, 6, 8 |
| `15-philosophy` | `studies/go-cli-study/reports/final/15-philosophy.md`; `pkg/cmd/factory/default.go:26-46`, `internal/backend/backend.go:19-90`, `fs/types.go:16-59`, `internal/pubsub/broker.go:10-19` | Complexity should be deliberate and not become generic infrastructure without product need. | Supports rejecting generic workflow/DAG/plugin/runtime-adapter scope. | Decisions 1, 2, 6, 9 |

## Trade-Off And Debt Analysis

### Accepted Trade-Offs

| Trade-Off | Benefit | Cost / Constraint Accepted | Why Acceptable Now | Revisit Trigger |
| --- | --- | --- | --- | --- |
| Persist state after every meaningful transition instead of batching writes. | Best recovery fidelity after interruption and truthful runtime-free status. | More filesystem writes and more explicit write-failure handling. | `AC-05` and `AC-06` require durable transition evidence; correctness matters more than write minimization. | Large studies show unacceptable IO overhead after measurement. |
| Conservative stale-lock policy instead of automatic deletion. | Prevents concurrent mutation by an active or ambiguous process. | Operators may need explicit force unlock after real process death. | State corruption risk is higher than inconvenience; `AC-08` requires explicit force behavior. | A later lock heartbeat or process-verification mechanism is implemented and tested. |
| Rich task metadata instead of minimal statuses. | Enables retry/fallback/status/debug decisions without runtime calls. | Increases schema surface, redaction needs, and compatibility checks. | `AC-16` through `AC-20` require persisted retry, fallback, usage, cost, cleanup, and diagnostic facts. | Sprint 14 stable JSON schema work decides a smaller public shape. |
| Completed-output revalidation on resume instead of trusting persisted completion. | Prevents false success when artifacts were deleted, corrupted, or never valid. | Resume startup can be slower and can move completed tasks back to schedulable or failed states. | PRD and TRD require runtime success not be product success. | Validation cost becomes high enough to require cached validation fingerprints. |
| Focused `internal/study` files instead of new global packages. | Preserves module ownership and avoids generic workflow scope creep. | The study package may become dense. | Architecture docs explicitly prefer product-module ownership for current scope. | Multiple product modules need the same behavior with no study semantics. |
| Use fake-runtime verification by default instead of real OpenCode smoke tests. | Keeps normal tests deterministic, offline, and credential-free. | Real adapter regressions require gated integration outside default verification. | `AC-24` excludes real runtime requirements from normal tests. | Release validation selects gated OpenCode smoke tests. |

### Potential Technical Debt

| Debt / Shortcut | Why It Might Accrue | Current Mitigation | Owner / Follow-Up |
| --- | --- | --- | --- |
| Run-state schema grows before Sprint 14 stabilizes JSON. | Rich durable metadata can become hard to migrate if exposed as public API too early. | Version state and reject unsupported durable formats explicitly; document that public JSON stability is deferred. | Sprint 14 schema-stability work. |
| Study package density increases. | Run-loop adds state machine, locks, scheduling, recovery, and status behavior in one module. | Use focused files such as `run_loop.go`, `locks.go`, `run_state.go`, `state.go`, and `service.go`; avoid premature subpackages. | Architecture review after implementation. |
| Force unlock UX can be operationally blunt. | Explicit force unlock may delete a lock without proving owner death. | Scope force unlock to selected study, report lock metadata, require explicit flag/command, and test active-looking lock refusal. | Future lock heartbeat or stale-owner verification sprint. |
| Runtime metadata redaction policy may need refinement. | Agentwrap metadata may add fields or native payload safety categories. | Persist safe summaries and omission facts, not unsafe raw bytes; use typed safe fields and redaction helpers where practical. | Runtime/observability follow-up if metadata expands. |
| Retry state may differ from agentwrap internals over time. | UltraPlan persists summaries while agentwrap owns policy execution. | Derive attempt, target, fallback, retry time, and classification from agentwrap metadata/events rather than duplicating policy. | Agentwrap version upgrade review. |

### Future Considerations

| Consideration | Deferred Until | Reason Deferred | What Should Be Preserved Now |
| --- | --- | --- | --- |
| Stable public JSON schemas for status/run-loop diagnostics. | Sprint 14 or explicit schema sprint. | Current sprint needs durable behavior, not public compatibility guarantees. | Include schema version and explicit rejection behavior; keep output deterministic. |
| Broad migration framework. | Migration-focused sprint or public schema commitment. | Requirements allow only compatible migration or explicit rejection now. | Do not silently reinterpret unknown state formats. |
| Gated OpenCode smoke tests. | Release smoke validation or integration-test sprint. | Default tests must remain fake-runtime and offline. | Ensure runtime integration goes through agentwrap so gated tests can be added later. |
| Generic workflow/DAG engine. | Product scope revision. | PRD/TRD and sprint requirements explicitly reject it. | Keep task concepts study-specific and avoid global workflow packages. |
| Target and sprint command families. | Later product phase. | Current build is study-side only. | Avoid product imports or schemas that assume target/sprint execution. |
| Lock heartbeat/process-owner verification. | Operational hardening follow-up. | Conservative lock plus explicit force unlock satisfies current acceptance. | Persist lock PID, command summary, timestamp, and study scope. |

## Final Decisions

### Decision 1: Keep Run-Loop Study-Owned With Thin CLI Wiring

- **Decision:** Implement durable run-loop orchestration in `internal/study`, expose it through a study service method, and keep `internal/app/study_commands.go` limited to argument parsing, resolved config/flag overrides, cancellation setup, command help, and rendering.
- **Rationale:** Run-loop behavior transforms study state, source applicability, prompts, report validation, synthesis gating, retry outcomes, and status. Those are study-owned product semantics, not CLI or platform semantics.
- **Study / Source Grounding:** `technical-handbook.md` cites thin CLI patterns from `01-project-structure` (`main.go:16`, `cmd/root.go:112-127`) and command factory patterns from `02-command-architecture` (`pkg/cmd/install.go:132`, `cmd/restic/cmd_backup.go:84-115`). `ARCHITECTURE.md` requires `internal/study` to own validation, scheduling, reports, prompts, state, and persistence, while `internal/platform/runtime` remains generic.
- **Trade-Offs Accepted:** The study package takes on more behavior, but focused files keep ownership clear. CLI command tests may need service fakes, but the command layer remains maintainable.
- **Technical Debt / Future Impact:** Avoids introducing a premature generic orchestration API. If multiple product modules later need shared scheduling, a future decision can extract genuinely reusable infrastructure.
- **Alternatives Rejected:** CLI-owned `RunE` state machine was rejected because it would bury product logic in command wiring. A global `internal/workflow` or `internal/scheduler` package was rejected because PRD/TRD non-goals exclude generic workflow scope and `ARCHITECTURE.md` warns against global technical layers. Direct runtime orchestration in `internal/platform/runtime` was rejected because runtime packages must not know studies, dimensions, sources, report semantics, or synthesis gating.
- **Contracts Satisfied:** Architecture, Configuration, CLI Surface, Documentation, Workflows, Performance, `AC-01`, `AC-11`, `AC-12`, `AC-23`, `AC-27`, `CON-01`, `CON-02`, `CON-03`.
- **Evidence Required:** Help and command tests for `ultraplan study <study> run-loop`; service-level tests for run-loop behavior; import review showing no product imports from `internal/platform/runtime`; architecture review confirming thin CLI wiring and no global workflow/scheduler packages.

### Decision 2: Materialize A Deterministic Study Task State Machine

- **Decision:** Store a deterministic run-state task matrix containing only applicable analysis tasks plus synthesis tasks for selected dimensions. Use explicit task statuses for `pending`, `running`, `retrying` or `waiting`, `validating`, `completed`, `failed`, `cancelled`, `skipped`, and `synthesis_unblocked` or equivalent typed transition metadata. Stable task IDs must derive from task kind, dimension, and source where applicable.
- **Rationale:** Runtime-free status and resumability require product state to be the source of truth. Markdown inapplicability must be excluded from required work rather than represented as failures. Synthesis must become schedulable only after all applicable per-source reports for a dimension validate.
- **Study / Source Grounding:** `technical-handbook.md` identifies durable state as product truth and runtime-free status as core patterns. `TRD.md` sections 13.1 through 13.5 define task state fields, run-state file content, resume logic, and scheduling. `PRD.md` requires explicit state/errors and says runtime success is not product success.
- **Trade-Offs Accepted:** Materializing synthesis and rich statuses adds state complexity but makes status, resume, and tests deterministic. Skipped inapplicable pairs are not counted as failures, which means totals must clearly distinguish applicable work from skipped/document-only exclusions.
- **Technical Debt / Future Impact:** Status names and JSON shape are not a stable public schema in this sprint. Versioned state and explicit rejection behavior must preserve room for Sprint 14 schema decisions.
- **Alternatives Rejected:** Inferring all tasks only from report files was rejected because it cannot explain retry, cancellation, stale recovery, or validation history. Treating inapplicable Markdown pairs as failed or missing was rejected because PRD/TRD require them to be skipped. Building a generic DAG was rejected because the workflow is study-specific and generic workflow engines are out of scope.
- **Contracts Satisfied:** Observability, Workflows, Performance, Persistence And Migrations, `AC-02`, `AC-05`, `AC-09`, `AC-12`, `AC-19`, `AC-22`, `AC-23`, `CON-04`, `CON-09`.
- **Evidence Required:** `run_loop_test.go` and `run_state.go` tests for deterministic matrix creation, Markdown applicability exclusion, synthesis unblocking, status totals, typed outcomes, and state snapshots after transitions.

### Decision 3: Use Atomic Run-State Persistence With Explicit Schema Handling

- **Decision:** Persist `studies/<study>/.ultraplan/run-state.json` atomically on every meaningful transition using a same-directory temporary file, flush/close, rename over the prior state, and best-effort directory sync. State load must validate schema version and either migrate only explicitly compatible prior state or reject unsupported/malformed/unsafe state before runtime launch.
- **Rationale:** The state file is the recovery substrate and status source. A failed write must not corrupt the last known good state or allow success to be reported for a run whose required durable state was not saved.
- **Study / Source Grounding:** `technical-handbook.md` cites filesystem abstraction patterns from `06-io-abstraction` (`internal/fs/interface.go:10-31`) and identifies atomic persistence as a primary trade-off. `TRD.md` section 13.3 requires same-directory temp writes, flush/close, and rename; `requirements.md` requires best-effort directory sync and prior-state preservation tests.
- **Trade-Offs Accepted:** Frequent atomic writes increase IO and failure surface, but they provide the interruption safety required by this sprint.
- **Technical Debt / Future Impact:** Only narrow compatibility/rejection behavior is implemented now. A broad migration framework is deferred.
- **Alternatives Rejected:** Batch-only state writes were rejected because they lose meaningful transition history on interruption. In-place JSON overwrite was rejected because partial writes can corrupt durable state. Silent acceptance of unknown schemas was rejected because it can reinterpret old state unsafely.
- **Contracts Satisfied:** Persistence And Migrations, Errors, Security, Observability, `AC-03`, `AC-04`, `AC-05`, `AC-06`, `AC-21`, `CON-06`, `CON-07`, `CON-08`.
- **Evidence Required:** `state_test.go` coverage for atomic writes, flush/close/rename failure injection, best-effort directory sync behavior where testable, malformed/unsupported state rejection, and prior valid state preservation after write failure.

### Decision 4: Add Per-Study Run-Loop Locking With Conservative Stale Handling

- **Decision:** Acquire a per-study lock before mutating run state. The lock payload must include PID, sanitized command summary, timestamp, and study scope. A second run-loop must fail loudly. Active-looking or ambiguous locks must not be deleted automatically. Force unlock must be explicit, scoped to the selected study, and tested.
- **Rationale:** Durable state correctness depends on a single mutator per study. Conservative stale handling prevents accidental state corruption when lock ownership is ambiguous.
- **Study / Source Grounding:** `technical-handbook.md` cites restic lock modeling (`internal/restic/lock.go:47`, `internal/restic/lock.go:105`, `internal/restic/lock.go:290-305`) and warns not to auto-delete active-looking locks. `TRD.md` section 21.2 requires per-study locks with PID, command, timestamp, conservative stale detection, and force unlock.
- **Trade-Offs Accepted:** Operators may need to run an explicit force unlock after real process death. That cost is acceptable because concurrent mutation would be worse.
- **Technical Debt / Future Impact:** Lock heartbeat/process-owner verification is deferred. Current lock payload should preserve fields needed to add that later.
- **Alternatives Rejected:** No lock was rejected because two run-loops could corrupt or race state transitions. Automatic stale lock deletion was rejected because process liveness is platform-sensitive and can be wrong. Global workspace locks were rejected because they unnecessarily block different studies.
- **Contracts Satisfied:** Persistence And Migrations, Security, Errors, CLI Surface, Workflows, `AC-07`, `AC-08`, `AC-21`, `CON-06`, `CON-08`.
- **Evidence Required:** `locks_test.go` for acquisition, second-lock refusal, stale-looking lock diagnostics, scoped force unlock, lock cleanup on terminal exit, and lock failure CLI exit behavior.

### Decision 5: Reconcile Resume By Recovering Stale Running Tasks And Revalidating Completed Outputs

- **Decision:** On resume, load and validate state, preserve filters unless explicitly overridden, recover stale `running` or active-like tasks to safe schedulable or terminal states while preserving attempts/error history, record recovery decisions in state/status, and revalidate completed analysis and synthesis artifacts before trusting them. Missing or invalid artifacts must move tasks back to non-completed or failed validation states with clear diagnostics.
- **Rationale:** A prior process may exit while tasks are marked running or while artifacts have been deleted/corrupted. Resumability requires reconciling durable state with product validation, not trusting prior markers blindly.
- **Study / Source Grounding:** `technical-handbook.md` highlights completed-output revalidation as a key trade-off and durable state as product truth. `PRD.md` requires completed output detection, validation of state marked complete, interrupted state saving, and synthesis gating. `TRD.md` section 13.4 requires stale running reset, completed-output revalidation, synthesis queueing, agentwrap record reconciliation when available, attempt preservation, and original filter respect.
- **Trade-Offs Accepted:** Resume can be slower and can alter previously completed-looking tasks. This is acceptable because product success requires valid artifacts.
- **Technical Debt / Future Impact:** Validation fingerprints could optimize future resume if cost grows. Current sprint should prefer correctness over caching.
- **Alternatives Rejected:** Trusting completed state without artifact validation was rejected because runtime/state success is not product success. Failing the entire run on any stale running task was rejected because interrupted runs should be recoverable. Resetting attempts/history was rejected because retry/fallback/status decisions need historical facts.
- **Contracts Satisfied:** Workflows, Observability, LLM Evaluation / Cost / Safety, Persistence And Migrations, `AC-03`, `AC-09`, `AC-10`, `AC-12`, `AC-19`, `CON-09`.
- **Evidence Required:** Resume tests for stale running, retrying, completed, failed, missing artifacts, invalid reports, synthesis gate recalculation, attempt preservation, recovery metadata, and original filter preservation.

### Decision 6: Schedule With Bounded Workers And Context-Driven Cancellation

- **Decision:** Use bounded worker scheduling that never starts more runtime tasks than resolved configuration or CLI overrides allow. Stop scheduling new work on cancellation, cancel active runtime runs through the platform runtime/agentwrap boundary, drain or safely consume event streams, call `Wait` for every started run, persist cancellation/retryable status for active tasks, and return the existing cancellation exit convention for user interrupt.
- **Rationale:** Long-running run-loop reliability requires concurrency, cancellation, event consumption, runtime cleanup, and state persistence to be handled as one design.
- **Study / Source Grounding:** `technical-handbook.md` cites bounded concurrency from `08-concurrency` (`go-task/task.go:87`, `gdu/pkg/analyze/parallel.go:13`, `k9s/internal/pool.go:21,30,37`) and context cancellation from `07-state-context` (`pkg/cmd/install.go:333-347`, `cmd/restic/cleanup.go:24-38`). `TRD.md` sections 13.5 and 20 require bounded pools, no unbounded channels, drained events, `Wait`, cancellation, cleanup, and no goroutine leaks.
- **Trade-Offs Accepted:** Bounded workers may underutilize possible runtime capacity, but predictable resource use and state correctness are required.
- **Technical Debt / Future Impact:** Scheduler fairness between newly unblocked synthesis tasks and remaining analysis tasks may need tuning after large-study observations. Current decision requires no starvation and deterministic tests.
- **Alternatives Rejected:** Goroutine-per-task launch was rejected because it can exceed runtime/process limits. Fire-and-forget cancellation was rejected because event channels and subprocess cleanup can deadlock or leak. Continuing to schedule new tasks after interrupt was rejected because it violates cancellation semantics and can worsen recovery.
- **Contracts Satisfied:** Workflows, Performance, Errors, Observability, `AC-13`, `AC-14`, `AC-15`, `AC-21`, `AC-25`, `CON-04`, `CON-05`.
- **Evidence Required:** Fake-runtime tests with blocked workers asserting max concurrency, cancellation before queue feed, cancellation during active runs, event drain and `Wait` calls, persisted cancelled/retryable tasks, no new work after cancellation, and `go test -race ./...`.

### Decision 7: Reuse Existing Runtime Boundary, Prompt Builders, Validators, Retry/Fallback Policy, And Metadata

- **Decision:** Analysis and synthesis tasks must reuse existing single-run prompt composition, runtime request mapping, workdir rules, permission policy mapping, runtime event handling, report validation, synthesis prompt composition, and final-report validation. Runtime execution must go through the existing platform runtime backed by agentwrap/OpenCode. Retry, fallback, validation, permission, cleanup, usage, cost, and safe diagnostic metadata must be persisted from agentwrap-facing runtime results/events/metadata. Runtime errors must be classified with typed errors or `agentwrap.SDKError` fields, never string parsing.
- **Rationale:** This sprint adds durable orchestration, not a new execution path. The runtime adapter already owns OpenCode process behavior, policy, validation wrapping, events, health, permissions, and metadata.
- **Study / Source Grounding:** `technical-handbook.md` cites typed retry/fatal errors from `05-error-handling` (`fs/fserrors/error.go:22-29`, `go-task/errors/errors.go:47-50`) and structured observability from `10-logging-observability` (`internal/slogs/keys.go:6-231`). `TRD.md` sections 11, 12, 14, 15, 19, and 22 require agentwrap `ObservingRuntime`, `ValidatingRuntime`, `PolicyRunner`, `opencode`, `SDKError`, permission metadata, and safe metadata handling.
- **Trade-Offs Accepted:** The durable state schema must carry enough runtime metadata to be useful, which increases redaction and compatibility responsibility.
- **Technical Debt / Future Impact:** Agentwrap upgrades may require state metadata mapping review. Unknown usage values must stay unknown, not zero, to avoid false accounting.
- **Alternatives Rejected:** Direct OpenCode invocation was rejected because TRD and sprint constraints require agentwrap/opencode. Reimplementing retry/fallback in UltraPlan was rejected because agentwrap policy owns those decisions. Parsing stderr/stdout strings for error categories was rejected because typed `SDKError` fields are available and safer.
- **Contracts Satisfied:** LLM Runtime, LLM Evaluation / Cost / Safety, Observability, Errors, Security, Configuration, `AC-11`, `AC-12`, `AC-16`, `AC-17`, `AC-18`, `AC-20`, `AC-22`, `CON-02`, `CON-03`, `CON-08`.
- **Evidence Required:** Fake-runtime tests emitting retry, rate-limit, fallback, validation, permission, cleanup, usage, cost, unknown usage, and typed SDK errors; state/status assertions for persisted safe metadata; import review for no direct OpenCode execution from product code.

### Decision 8: Make Status Runtime-Free, Deterministic, And Safe

- **Decision:** `ultraplan study <study> status` must load durable state and lock information only, make no runtime calls, and report progress totals, active tasks, recent/recovered tasks, failed tasks, retry times, cancelled tasks, state path, lock diagnostics, validation failures, artifact paths, and safe runtime metadata. Human output must be concise and deterministic and must not print prompts, embedded Markdown source bodies, API keys, sensitive environment values, or unsafe native runtime payload bytes.
- **Rationale:** Operators need to inspect and decide without relaunching runtime dependencies. The state file is the durable truth, while safe diagnostics prevent leaking sensitive runtime or prompt content.
- **Study / Source Grounding:** `technical-handbook.md` cites runtime-free status patterns from observability and terminal UX reports (`internal/ui/progress/printer.go:15-34`, `pkg/iostreams/iostreams.go:52-54`) and security redaction examples (`internal/options/secret_string.go:15-20`, `pkg/registry/transport.go:37-41`). `PRD.md` requires status for active tasks, failed tasks, retry times, outputs, and summaries. `TRD.md` sections 19 and 22 require safe diagnostics and redaction.
- **Trade-Offs Accepted:** Status may not know live runtime liveness beyond lock and persisted state. That is acceptable because runtime-free status is required and active-looking locks are handled conservatively.
- **Technical Debt / Future Impact:** Stable JSON status is deferred. Current text output should still be deterministic and testable.
- **Alternatives Rejected:** Live runtime probing in status was rejected because `AC-19` requires runtime-free status and offline tests. Printing full runtime/debug payloads was rejected because security constraints prohibit unsafe diagnostics. Non-deterministic active/recent ordering was rejected because tests and operators need stable output.
- **Contracts Satisfied:** Observability, Security, CLI Surface, Documentation, Performance, `AC-19`, `AC-20`, `AC-22`, `AC-24`, `CON-08`, `CON-12`.
- **Evidence Required:** `study_status_commands_test.go` coverage for runtime-free loading, progress counts, active/recent tasks, retry times, failed/cancelled tasks, stale recovery, state path, lock diagnostics, redaction, and no runtime collaborator invocation.

### Decision 9: Verify With Fake Runtime, State/Lock Failure Injection, Command Tests, Race Tests, Build, And Architecture Review

- **Decision:** Implement a test matrix centered on fake runtimes and local fixtures: run-loop state machine tests, lock tests, state persistence tests, CLI command tests, status tests, cancellation/concurrency tests, redaction tests, and architecture review. Required verification commands are `go test ./...`, `go test -race ./...`, and `go build ./cmd/ultraplan` from `/home/antonioborgerees/coding/ultraplan-go`.
- **Rationale:** The sprint's reliability risks are state transitions, atomic writes, locks, cancellation, retry/fallback metadata, validation, and status rendering. These must be deterministic and offline by default.
- **Study / Source Grounding:** `technical-handbook.md` cites fake/mock testing from `11-testing-strategy` (`pkg/httpmock/stub.go:35-199`, `internal/backend/mock/backend.go:14-26`, `task_test.go:166-169`) and fixture/fake runtime needs. `TRD.md` sections 23 and 24 require fake-runtime tests, gated OpenCode integration only, `go test ./...`, and CLI build.
- **Trade-Offs Accepted:** Real OpenCode regressions are not caught by normal tests. That is acceptable because the sprint explicitly requires normal tests to be deterministic and offline.
- **Technical Debt / Future Impact:** Gated real runtime smoke tests should be added or run in release validation, but not as default sprint evidence.
- **Alternatives Rejected:** Manual-only verification was rejected because state machines and races need automated coverage. Real-runtime default tests were rejected because they would require credentials/network/subprocesses. Skipping race-sensitive verification was rejected because locks, cancellation, and workers are central to this sprint.
- **Contracts Satisfied:** Testing, Architecture, Security, Performance, Documentation, `AC-24`, `AC-25`, `AC-26`, `AC-27`, `CON-12`.
- **Evidence Required:** Passing targeted tests and full verification commands; sprint review recording acceptance evidence, deviations, residual risks, and architecture review results.

## Expected Evidence

| Evidence Type | Required Evidence | Source / Command / Review Check |
| --- | --- | --- |
| Tests | Run-loop command is available with help, flags, usage errors, lock failures, resume failures, cancellation, runtime failures, validation failures, and partial completion. | `internal/app/study_run_loop_commands_test.go`; command help snapshots or assertions. |
| Tests | Initial state creation excludes inapplicable Markdown pairs, produces deterministic task IDs/statuses, and preserves filters. | `internal/study/run_loop_test.go`; task matrix fixture assertions. |
| Tests | Resume validates state schema, rejects unsupported/malformed state, recovers stale running tasks, revalidates completed outputs, and preserves attempts/history. | `internal/study/run_state.go`; `internal/study/run_loop_test.go`; `internal/study/state_test.go`. |
| Tests | Atomic writes use same-directory temp file, flush/close, rename, best-effort directory sync, and preserve prior valid state on failure. | `internal/study/state_test.go` with injected write/rename/sync failures. |
| Tests | Per-study lock refuses concurrent run-loop, reports metadata, handles conservative stale diagnostics, supports scoped explicit force unlock, and cleans up on terminal exit. | `internal/study/locks_test.go`; command-level lock tests. |
| Tests | Bounded scheduling never exceeds configured parallelism, stops queue feeding on cancellation, drains events, calls `Wait`, persists active task outcomes, and returns cancellation exit convention. | Fake-runtime concurrency and cancellation tests; `go test -race ./...`. |
| Tests | Retry/fallback/rate-limit/validation/permission/cleanup/usage/cost metadata is persisted safely, and unknown usage is not converted to zero. | Fake-runtime tests in `internal/study/run_loop_test.go`; state/status field assertions. |
| Tests | Status is runtime-free and reports progress, active/recent tasks, failures, retry times, cancelled tasks, recovery decisions, state path, lock diagnostics, and redacted output. | `internal/app/study_status_commands_test.go`; runtime fake must not be invoked. |
| Runtime | No direct OpenCode invocation from product code; runtime calls go through agentwrap-backed platform runtime. | Import/search review; architecture review protocol. |
| Review | No target/sprint workflows, code-reference extraction, generic workflow engine, new runtime adapter, or broad migration framework were introduced. | Architecture review and scope review against `requirements.md` non-goals. |
| Review | Errors and diagnostics avoid prompts, Markdown source bodies, API keys, sensitive env values, and unsafe native payload bytes. | Redaction tests and manual review of state/status/log fields. |
| Build | Full tests, race tests, and CLI build pass. | `go test ./...`; `go test -race ./...`; `go build ./cmd/ultraplan`. |
| Documentation | Command help documents run-loop flags for filters, parallelism, retry/cancellation behavior, and lock override/force unlock behavior. | CLI help tests and rendered help inspection. |

## Assumptions And Risks

| Item | Type | Impact | Mitigation / Follow-Up |
| --- | --- | --- | --- |
| Existing Sprint 8 run-state primitives can be extended without rewriting all state handling. | Assumption | If false, implementation may need more state refactoring than planned. | Keep changes focused in `domain.go`, `run_state.go`, and `state.go`; preserve explicit schema rejection. |
| Agentwrap metadata exposes enough safe policy, retry, fallback, validation, permission, cleanup, usage, and cost information. | Assumption | Missing metadata could reduce status fidelity. | Persist all available safe metadata and omission facts; record gaps in sprint review. |
| Completed-output validation is fast enough for normal resume. | Assumption | Large studies may have slower startup. | Prefer correctness now; consider validation fingerprints only if measured cost warrants it. |
| Lock stale detection cannot reliably prove owner death across all supported environments. | Risk | Incorrect auto-cleanup could corrupt state; overly conservative locks require operator action. | Use conservative refusal and scoped explicit force unlock; defer heartbeat/process verification. |
| Cancellation during state persistence or runtime cleanup can still fail to save final state. | Risk | Operators may see stale active tasks on next resume. | Persist before and after critical transitions where possible; recover stale running tasks on resume; make write failures loud. |
| Rich metadata may accidentally persist unsafe details. | Risk | Secrets, prompts, embedded Markdown, or unsafe payload bytes could leak into state/status. | Use safe fields, redaction helpers, omission facts, and tests for sensitive strings. |
| Race-sensitive worker/lock behavior may be flaky without careful fake-runtime design. | Risk | Tests may miss or intermittently expose leaks/races. | Use deterministic blocking fakes, explicit synchronization, and `go test -race ./...`. |
| Stable public JSON schema is deferred. | Risk | Users may still begin depending on current status/state JSON. | Version durable state, avoid promising public compatibility, and record schema follow-up for Sprint 14. |

## Implementation Constraints

- Run-loop orchestration, task state machines, source applicability, synthesis gating, report validation, resume reconciliation, and status summaries must live in `internal/study`.
- `internal/app` must remain a thin command/composition layer for flags, config overrides, cancellation setup, service invocation, exit handling, and rendering.
- `internal/platform/runtime` must remain generic and must not import or know `study`, dimensions, sources, report semantics, synthesis gating, or run-state machines.
- Product code must not invoke OpenCode directly or parse native OpenCode stdout/stderr; runtime behavior must come through agentwrap-backed platform runtime results, events, metadata, and typed errors.
- Every long-running work path must accept and propagate `context.Context`; no `context.Background()` may be introduced in run-loop work paths except a deliberately bounded cleanup context if implementation proves it necessary and tests it.
- Worker concurrency must be bounded by resolved config or CLI override; no unbounded goroutine creation, unbounded task channels, or background worker leaks are allowed.
- State and lock paths must remain under the selected study directory, be workspace-contained, and reject unsafe path traversal or unexpected absolute-path behavior.
- Every meaningful task transition must be persisted atomically, and required write failures must make the run-loop fail loudly rather than reporting success.
- Existing prompt builders, runtime request mapping, workdir rules, permission policy mapping, single-run behavior, report validators, rating parser, source applicability helper, and synthesis prompt/final-report validation must be reused rather than duplicated.
- Markdown document source/dimension pairs that are inapplicable must be excluded from task creation, completed-output detection, synthesis gating, status missing-report calculations, and summary expectations; they are not failures.
- Runtime success is not product success. Tasks complete only after required artifacts exist and pass the appropriate analysis or synthesis validation.
- Runtime errors must be classified through typed errors or `agentwrap.SDKError` fields; no string parsing may drive retry, fallback, cancellation, validation, or status categories.
- Safe diagnostics may include category, operation, provider/model, attempt, retry/fallback, validation, cleanup, usage, and cost facts, but must not include prompts, embedded Markdown source bodies, API keys, sensitive environment values, or unsafe native payload bytes.
- Unknown usage and cost values must remain unknown and must not be converted to zero.
- Normal tests must use fake runtimes and local fixtures only; real OpenCode, provider credentials, network access, and subprocess execution are excluded from default verification.
- Do not add target scaffolding, sprint planning, sprint execution, code-reference extraction, a generic workflow engine, plugin architecture, TUI, hosted service, browser UI, new runtime adapter, or broad migration framework.

## Plan Handoff

`plan.md` must execute these decisions. It must not invent architecture, scope, or decisions beyond this document.

The plan must carry forward:

- Final decisions 1 through 9.
- Contract mappings and requirement IDs `AC-01` through `AC-27` plus applicable `CON-01` through `CON-12`.
- Expected evidence for run-loop, state, lock, status, retry/fallback metadata, cancellation, race, build, documentation, and architecture review.
- Risks and assumptions listed above.
- Required review protocols from `sprint-index.md`: Architecture Review and Sprint Review.

## Phase Exit Criteria

- [x] Selected context was read and used.
- [x] Area-specific reasoning documents were completed where applicable or explicitly marked not applicable.
- [x] Area-specific reasoning conclusions are reflected or explicitly marked absent because no area reasoning files exist.
- [x] Contracts are explicitly mapped to decisions and expected evidence.
- [x] Final decisions are clear enough for `plan.md` to execute.
- [x] Expected evidence is specific and reviewable.
