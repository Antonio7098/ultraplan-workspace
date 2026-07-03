# Sprint Reasoning: 08 Run State Persistence

> **Inputs Used:** `projects/ultraplan-go/sprints/08-run-state-persistence/requirements.md`, `projects/ultraplan-go/project-index.md`, `projects/ultraplan-go/docs/ARCHITECTURE.md`, `projects/ultraplan-go/docs/PRD.md`, `projects/ultraplan-go/docs/TRD.md`, `projects/ultraplan-go/sprints/08-run-state-persistence/sprint-index.md`, `projects/ultraplan-go/sprints/08-run-state-persistence/technical-handbook.md`, `templates/sprint-reasoning.md`
> Project: `ultraplan-go`
> Sprint: `08-run-state-persistence`
> Output: `projects/ultraplan-go/sprints/08-run-state-persistence/reasoning.md`

This document decides. It synthesizes selected context, handbook evidence, area-specific reasoning, and contracts into final sprint decisions.

It does not replace `sprint-index.md`, `technical-handbook.md`, or any future `reasoning/*.md` files.

## Sprint Purpose

- **Goal:** Implement study-owned durable run-state persistence for planned analysis and synthesis tasks, including deterministic task construction, schema validation, atomic JSON load/save, completed-output revalidation hooks, deterministic status summaries, and runtime-free `ultraplan study <study> status` human output.
- **Non-Goals:** Do not execute analysis or synthesis tasks, invoke agentwrap/OpenCode, add worker-pool orchestration, implement retry/fallback/backoff execution, add runtime health gates, persist full agentwrap run stores, add cancellation handling, add study lock files, promise stable public JSON status output, generate `summary.csv`, extract code references, or implement target/sprint workflows.
- **Depends On:** Prior study-side foundations for CLI composition, workspace/config paths, study/source/dimension resolution, study initialization layout, Markdown source applicability, report validation, and prompt/output path concepts. The project index lists no selectable prior decisions, so carry-forward constraints come from `requirements.md`, `sprint-index.md`, PRD/TRD/Architecture docs, and selected handbook evidence.

## Selected Context And Pre-Reasoning Artifacts

| Artifact | Path | How It Was Used |
| --- | --- | --- |
| Sprint Requirements | `projects/ultraplan-go/sprints/08-run-state-persistence/requirements.md` | Binding sprint scope, acceptance criteria, non-goals, constraints, dependencies, and review expectations. Requirement lines 29-51 define state schema, task metadata, applicability filtering, atomic writes, load diagnostics, resume validation, status summaries, command behavior, and verification commands. |
| Project Index | `projects/ultraplan-go/project-index.md` | Source of available documents, contracts, evidence reports, and protocols. It confirms `ultraplan-go` target implementation path and selected study evidence pool. |
| Product Requirements | `projects/ultraplan-go/docs/PRD.md` | Product grounding for explicit state, resumability, runtime-success-is-not-product-success, status inspection, durable local artifacts, and inapplicable Markdown pairs being skipped rather than failed. |
| Technical Requirements | `projects/ultraplan-go/docs/TRD.md` | Technical grounding for run-state fields, state file path, atomic writes, schema version validation, resume behavior, task states, applicability filtering, report validation, diagnostics, and test expectations. |
| Architecture | `projects/ultraplan-go/docs/ARCHITECTURE.md` | Boundary rule that `internal/study` owns study state, persistence, validation, reports, scheduling concepts, and summaries, while CLI rendering stays in `internal/app` and platform packages must not import product modules. |
| Sprint Index | `projects/ultraplan-go/sprints/08-run-state-persistence/sprint-index.md` | Selected contracts, evidence reports, protocols, non-goals, and excluded context. It specifically excludes runtime execution, concurrency evidence, extensibility evidence, target/sprint workflows, locks, summary generation, and code extraction. |
| Technical Handbook | `projects/ultraplan-go/sprints/08-run-state-persistence/technical-handbook.md` | Main study evidence synthesis. Used to decide module ownership, command thinness, schema validation, atomic persistence test seams, resume validation, status determinism, diagnostics, secret-free persisted summaries, test strategy, and performance posture. |
| Sprint Reasoning Template | `templates/sprint-reasoning.md` | Structure for this consolidated reasoning document and required decision/evidence sections. |

## Area-Specific Reasoning Inputs

No area-specific reasoning files exist for this sprint. The directory `projects/ultraplan-go/sprints/08-run-state-persistence/reasoning/` was not present when this document was created, so final decisions are synthesized directly from `requirements.md`, `sprint-index.md`, project docs, and `technical-handbook.md`.

| Area | Reasoning Document | Key Conclusion | Evidence Basis | Impact On Final Decision |
| --- | --- | --- | --- | --- |
| None | N/A | No area-specific reasoning document was available to summarize. | N/A | Consolidated reasoning resolves all open handbook questions directly in the final decisions below. |

## Sprint Technical Handbook Summary

- **Relevant Patterns:** Keep study-owned state behavior inside `internal/study`; keep status command handlers thin; validate schema versions in one explicit load path; write state atomically with same-directory temp files and rename; revalidate completed reports before trusting persisted completion; derive status from persisted state and known study/report metadata; wrap errors with actionable path-aware diagnostics; persist only secret-free config summaries; use table-driven unit tests plus command-output tests.
- **Important Trade-Offs:** Use targeted persistence seams rather than a broad filesystem abstraction; omit inapplicable Markdown pairs from task lists instead of recording skipped tasks for every impossible pair; include minimal future runtime placeholders without implementing runtime integration; reject unsupported schema versions strictly; keep status human and concise rather than committing to stable JSON; keep status loading lazy and non-runtime.
- **Warnings / Anti-Patterns:** Do not introduce global `internal/state`, `internal/scheduler`, `internal/reports`, or `internal/validation` packages; do not put state construction or validation in `RunE`; do not silently accept malformed or unknown-version state; do not trust completed tasks without artifact validation; do not persist secrets; do not invoke runtime or inspect processes during status; do not recursively scan repositories; do not treat inapplicable Markdown pairs as missing or failed.
- **Evidence Confidence:** High for module boundaries, command thinness, error handling, IO test seams, testing strategy, security, and performance because the handbook cites multiple mature Go CLI reports and concrete source references. Medium for terminal UX, observability, state/context, and philosophy because those sources shape posture and trade-offs but are less prescriptive about this exact state schema.

## Contracts Applied

| Contract / Requirement ID | Constraint | Decision Impact | Expected Evidence |
| --- | --- | --- | --- |
| Architecture | Product behavior and state stay with owning module; CLI entrypoints remain thin; platform packages must not import `study`. | Run-state schema, construction, validation, persistence, resume validation, and summaries live in `internal/study`; `internal/app` only wires and renders status. | Import/package review; tests exercise study APIs separately from command rendering. |
| Errors | Diagnostics must preserve cause chains and be actionable. | Missing, malformed, unsupported, invalid, report-validation, and save/load failures use typed or sentinel errors plus `%w` wrapping and name relevant paths. | Unit tests assert missing vs malformed vs unsupported behavior and command diagnostics name the state file path. |
| Configuration | Persist safe effective config summaries without secrets. | `ConfigSummary` is deliberately narrow: runtime/model/variant/parallelism/timeout-like operational facts if already safely available, with no env dumps, tokens, provider credentials, or raw runtime payloads. | State construction tests and review confirm no secret fields or full environment persistence. |
| Observability | Status must truthfully report durable current/recent workflow state. | Status summary is derived from persisted tasks and validation metadata, not live process inspection, logs, network, or agentwrap stores. | Command tests prove status output from fixtures and no runtime dependency. |
| Security | Paths must be safe and persisted state must avoid secrets. | Persist workspace-relative paths where practical; state path is under `studies/<study>/.ultraplan/run-state.json`; diagnostics identify paths without exposing credentials. | Security review; tests for expected relative output paths and no runtime/env fields. |
| Testing | Durable state behavior needs unit, fixture, persistence, and command tests. | Add focused tests for construction, stable IDs, filtering, synthesis dependencies, atomic save/load, invalid/missing state, resume revalidation, stale active reset, status aggregation, and command output. | `go test ./...` and specific tests in `internal/study/state_test.go` and `internal/app/study_status_commands_test.go`. |
| CLI Surface | `ultraplan study <study> status` must be script-friendly, deterministic, and actionable. | Human output includes state path, run ID, completion flag, totals, failed count, active count, retry count, and next retry time when present; missing/invalid state returns clear diagnostics. | Command tests for missing, valid, invalid, failed-task, and no-runtime behavior. |
| Workflows | Durable workflow state needs explicit states, terminal states, resumability boundaries, and diagnostics. | Define explicit task statuses and reset stale active statuses safely on resume; preserve attempts and safe last-error details. | Unit tests for stale status reset, status counts, and attempt/error preservation. |
| Performance | Status must avoid expensive scans and eager runtime initialization. | Status reads state and required study/report metadata only; it does not recursively scan source repositories or initialize runtime/provider dependencies. | Command tests/fakes confirm no runtime execution; review checks avoid repository scans in status path. |
| Persistence And Migrations | Durable formats need schema versioning, atomic writes, and unsupported-version rejection. | Top-level schema version is required; unsupported versions fail clearly; writes use same-directory temp file, flush/close, and rename. | Unit tests for save/load, malformed JSON, unsupported version, failed write preserving previous valid state where feasible. |
| Requirements AC 31-52 | Sprint acceptance criteria for schema, task records, filtering, ordering, atomic writes, validation, status, tests, and build. | Final decisions below map each implementation concern to specific accepted behavior. | Required test files plus `go test ./...` and `go build ./cmd/ultraplan` from `/home/antonioborgerees/coding/ultraplan-go`. |

## Repos Studied / Source Evidence Used

| Source / Repo / Report | Concrete Reference | Relevant Finding | Why It Matters For This Sprint | Used In Decision(s) |
| --- | --- | --- | --- | --- |
| `01-project-structure` | `studies/go-cli-study/reports/final/01-project-structure.md`; cited `internal/chezmoi/chezmoi.go:1-2`, `cmd/root.go:112-127` | Mature Go CLIs keep entrypoints thin and dependencies unidirectional through `internal/`. | Supports keeping run-state behavior in `internal/study` and preventing platform/product dependency inversion. | Decisions 1, 8 |
| `02-command-architecture` | `studies/go-cli-study/reports/final/02-command-architecture.md`; cited `pkg/cmd/install.go:132-145`, `pkg/cmdutil/factory.go:16-43`, `pkg/cmd/issue/list/list.go:47-118` | Commands should parse, delegate, and render rather than own business logic. | Shapes `study status` as a thin command that calls study load/summarize APIs. | Decisions 1, 7 |
| `03-dependency-injection` | `studies/go-cli-study/reports/final/03-dependency-injection.md`; cited `pkg/cmdutil/factory.go:16-43`, `executor.go:22-24`, `internal/chezmoi/system.go:25` | Constructor injection and small seams improve tests; global service state is risky. | Supports targeted seams for clock/filesystem write failures without broad global abstractions. | Decisions 3, 8 |
| `04-configuration-management` | `studies/go-cli-study/reports/final/04-configuration-management.md`; cited `internal/config/config.go:609-641`, `internal/config/json/validator.go:1-187` | Config and schema validation should occur in explicit validation paths. | Supports strict run-state schema version validation and narrow config summaries. | Decisions 2, 6 |
| `05-error-handling` | `studies/go-cli-study/reports/final/05-error-handling.md`; cited `fs/fserrors/error.go:22-29`, `errors/errors_task.go:13-32`, `pkg/storage/driver/driver.go:27-36` | Robust CLIs use wrapping, sentinels/typed errors, and separate user diagnostics from operational errors. | Supports missing vs malformed vs unsupported state errors and path-aware diagnostics. | Decisions 2, 3, 5, 7 |
| `06-io-abstraction` | `studies/go-cli-study/reports/final/06-io-abstraction.md`; cited `pkg/iostreams/iostreams.go:551-568`, `internal/fs/interface.go:10-31` | Injectable IO/filesystem seams enable deterministic command and persistence tests. | Supports testable atomic write failure paths and status command output capture. | Decisions 3, 7, 8 |
| `07-state-context` | `studies/go-cli-study/reports/final/07-state-context.md`; cited `pkg/cmd/install.go:333-347`, `internal/app/app.go:25-40` | Stateful tools benefit from explicit app/context propagation, but not context-as-service-locator. | Supports explicit study APIs and no hidden global state during status. | Decisions 1, 7 |
| `09-terminal-ux` | `studies/go-cli-study/reports/final/09-terminal-ux.md`; cited `internal/cmd/prompt.go:20-256`, `pkg/iostreams/iostreams.go:514-516` | Status UX should be calm, concise, non-blocking, and script-friendly. | Supports deterministic concise human status output. | Decision 7 |
| `10-logging-observability` | `studies/go-cli-study/reports/final/10-logging-observability.md`; cited `internal/slogs/keys.go:6-231`, `internal/logging/logging.go:31-66` | Structured fields and debug control are useful, but logs should not replace deterministic status. | Supports persisted status summaries instead of live logs/runtime inspection. | Decisions 5, 7 |
| `11-testing-strategy` | `studies/go-cli-study/reports/final/11-testing-strategy.md`; cited `internal/cmd/main_test.go:64-174`, `task_test.go:166-169`, `pkg/httpmock/stub.go:35-199` | Table-driven unit tests, command integration tests, golden/output tests, and fakes are proven. | Defines the sprint evidence strategy. | Decision 8 |
| `13-security` | `studies/go-cli-study/reports/final/13-security.md`; cited `internal/options/secret_string.go:15-20`, `pkg/registry/transport.go:37-41`, `status.go:332-338` | Durable artifacts and diagnostics must avoid secrets and treat paths as trust boundaries. | Supports secret-free config summaries and workspace-relative persisted paths. | Decisions 2, 3, 6, 7 |
| `14-performance` | `studies/go-cli-study/reports/final/14-performance.md`; cited `pkg/cmdutil/factory.go:27-42`, `yq/pkg/yqlib/stream_evaluator.go:78-113`, `gdu/pkg/analyze/parallel.go:13` | Fast status commands avoid eager initialization, recursive scans, and unbounded work. | Supports lazy status loading and no source repository scans. | Decisions 5, 7 |
| `15-philosophy` | `studies/go-cli-study/reports/final/15-philosophy.md`; cited `internal/chezmoi/system.go:25-45`, `pkg/cmd/factory/default.go:26-46`, `internal/pubsub/broker.go:10-19` | Good tools deliberately reject non-goals and keep complexity near the capability that needs it. | Supports excluding runtime execution, locks, broad scheduler abstractions, and public JSON status. | Decisions 1, 3, 7, 8 |

## Trade-Off And Debt Analysis

### Accepted Trade-Offs

| Trade-Off | Benefit | Cost / Constraint Accepted | Why Acceptable Now | Revisit Trigger |
| --- | --- | --- | --- | --- |
| Inapplicable Markdown pairs are absent from task lists rather than stored as `skipped` tasks. | Smaller state files, simpler task execution semantics, and no false missing-output failures. | Users cannot inspect every omitted pair as an explicit skipped task unless status/reporting later adds separate applicability metadata. | Requirements allow inapplicable pairs to be skipped or absent and require they not count as missing/failed. | If users need an audit view of skipped applicability pairs or public JSON status needs full matrix representation. |
| Persistence uses local JSON files, not an opaque database. | Reviewable in Git, simple recovery, aligns with PRD local-first artifacts. | No transactional multi-writer support or query engine. | Sprint excludes concurrent mutation and locks; local file state is explicitly required. | If multi-process coordination, history queries, or remote storage enter scope. |
| Atomic write implementation uses targeted test seams, not a broad filesystem platform package. | Keeps implementation small and aligned with module-owned state behavior. | Some platform-specific write/rename failures may need careful fixture hooks. | Handbook supports abstraction only where side effects are core; this sprint needs persistence seams, not a global FS architecture. | If multiple modules need the same atomic-write behavior or failure injection becomes duplicated. |
| Strict schema-version rejection instead of permissive loading. | Clear migration boundary and safer durable state. | Older or hand-edited state files may fail until migration exists. | Sprint requires unsupported schema rejection; no migration support is in scope. | When a second schema version is introduced. |
| Human status only, no stable public JSON format. | Lower compatibility burden and faster delivery of reviewable operational UX. | Automation must wait for a later stable JSON status contract. | Requirements explicitly do not require stable public JSON. | When PRD/TRD select machine-readable status output for implementation. |
| Agentwrap metadata placeholders are modeled but not populated by runtime execution. | Future run-loop integration has stable fields for run/session/turn/validation metadata. | Some schema fields remain empty during this sprint. | Runtime execution is a non-goal, while task records must reserve placeholders. | When runtime execution and agentwrap run-store integration enter scope. |

### Potential Technical Debt

| Debt / Shortcut | Why It Might Accrue | Current Mitigation | Owner / Follow-Up |
| --- | --- | --- | --- |
| Minimal `ConfigSummary` could be too sparse for later runtime diagnostics. | This sprint cannot safely infer all runtime/provider details without executing runtime integration. | Persist only safe known operational fields and keep schema versioned. | Revisit in runtime execution/config mapping sprint. |
| Top-level schema version without migration code. | Future versions will need migration or explicit compatibility policy. | Reject unsupported versions with clear diagnostics now. | Add migration decision when schema version `2` is needed. |
| No explicit persisted skipped records for inapplicable pairs. | Auditability may be less complete for users comparing full source/dimension matrices. | Deterministic construction and synthesis dependencies will prove applicable inputs; status summaries exclude inapplicable pairs. | Revisit with public JSON status or full matrix diagnostics. |
| Atomic write fsync behavior may be platform-sensitive. | Go rename semantics and directory fsync support differ by filesystem/platform. | Use same-directory temp file, flush/close, rename, and document residual risk if directory sync is unavailable. | Revisit during cross-platform hardening. |
| Status command may initially duplicate some formatting patterns with other commands. | Existing app command style may not yet have reusable status rendering helpers. | Keep rendering small and command-specific; avoid premature UI abstraction. | Extract common rendering only after multiple commands need it. |

### Future Considerations

| Consideration | Deferred Until | Reason Deferred | What Should Be Preserved Now |
| --- | --- | --- | --- |
| Runtime execution and agentwrap metadata population | Later analysis/synthesis execution sprint | Current sprint must not invoke runtime or subprocesses. | Task ID, task kind, output path, attempt, error, validation, and agentwrap placeholder fields. |
| Retry/fallback/backoff execution | Later workflow policy sprint | Current sprint only stores statuses and retry times. | `Attempts`, `LastError`, `NextRetryAt`, `retrying`, `waiting`, and `failed` status semantics. |
| Worker-pool orchestration | Later run-all/run-loop sprint | Current sprint only constructs planned task records and status summaries. | Deterministic task ordering and synthesis dependency metadata. |
| Per-study locks and force unlock | Later concurrency/cancellation sprint | Explicit sprint non-goal. | Atomic writes and path-aware diagnostics, without pretending to prevent concurrent writers. |
| Stable machine-readable status | Later CLI/status contract sprint | Requirements only require deterministic human output. | Status summary model should be structured internally so JSON can be added without changing core logic. |
| Durable agentwrap run store and event sinks | Later observability/runtime sprint | Requires actual runtime integration. | Placeholder fields and no unsafe raw payload persistence. |
| Summary CSV integration | Later summary generation sprint | Explicitly out of scope. | Applicability filtering and output path metadata stay consistent with summary needs. |

## Final Decisions

### Decision 1: Keep Run State Study-Owned With Thin Status Wiring

- **Decision:** Define and implement run-state domain types, state construction, validation, persistence, resume validation, and status summaries in `internal/study`. Add or extend `internal/app/study_status_commands.go` only to resolve the study/workspace context, call study-owned load/summarize functions, and render deterministic human output.
- **Rationale:** Run state is product behavior tied to studies, dimensions, sources, reports, and applicability rules. Moving it to a global state/scheduler/report package would fracture the product model and violate the architecture. Command handlers should not own business logic because persistence and summary behavior need direct unit tests independent of CLI parsing.
- **Study / Source Grounding:** `technical-handbook.md` relevant patterns and anti-patterns cite `01-project-structure` (`internal/chezmoi/chezmoi.go:1-2`, `cmd/root.go:112-127`) for internal boundaries and `02-command-architecture` (`pkg/cmd/install.go:132-145`, `pkg/cmd/issue/list/list.go:47-118`) for thin delegation. `ARCHITECTURE.md` lines 23-32 and 114-158 make `internal/study` the owner of study state and persistence.
- **Trade-Offs Accepted:** The study package grows more files and responsibilities, but they are cohesive with study behavior. No new reusable platform state abstraction is created.
- **Technical Debt / Future Impact:** If another product module later needs identical atomic persistence semantics, a small platform helper can be extracted then. For now, avoid premature global packages.
- **Alternatives Rejected:** Reject `internal/state` because it detaches run-state logic from study/report/applicability semantics. Reject putting load/validate/summarize in `RunE` because it blocks focused unit tests and makes command handlers stateful.
- **Contracts Satisfied:** Architecture, CLI Surface, Testing, Performance; requirements AC 31-36, 39-46, 47-49; constraints lines 70-72.
- **Evidence Required:** Package/import review confirms `internal/study` ownership and no platform imports of `study`; unit tests cover study state APIs; command tests prove status rendering delegates without runtime execution.

### Decision 2: Use A Versioned JSON Run-State Schema With Explicit Task Metadata

- **Decision:** Persist `studies/<study>/.ultraplan/run-state.json` as human-reviewable JSON with a top-level schema version, run ID, created/updated timestamps, filters, safe config summary, task list, and completion flag. Analysis and synthesis tasks use stable IDs, explicit task kinds/statuses, dimension/source metadata, workspace-relative output paths where practical, attempts, timestamps, last error, next retry time, validation summary, and agentwrap metadata placeholders.
- **Rationale:** Durable state must survive process exit and support review, resume, validation, and status without runtime inspection. The schema needs enough data for later runtime integration but must not implement runtime work now.
- **Study / Source Grounding:** `technical-handbook.md` patterns for explicit durable state schema and secret-free summaries cite `04-configuration-management` (`internal/config/config.go:609-641`, `internal/config/json/validator.go:1-187`) and `13-security` (`internal/options/secret_string.go:15-20`, `pkg/registry/transport.go:37-41`). TRD section 13.1-13.2 defines task and run-state fields; PRD 2.6 defines task metadata and lifecycle states.
- **Trade-Offs Accepted:** The schema includes future-facing placeholders before runtime execution exists, increasing test surface. This is acceptable because requirements explicitly require agentwrap metadata placeholders and durable task records.
- **Technical Debt / Future Impact:** Placeholder shape may need migration when actual agentwrap store integration arrives. Top-level schema version gives a safe migration/rejection boundary.
- **Alternatives Rejected:** Reject an unversioned JSON file because durable compatibility would be ambiguous. Reject an opaque database because requirements demand reviewable local filesystem state and Git-safe artifacts. Reject deriving task metadata from reports at status time because it would be non-durable and less deterministic.
- **Contracts Satisfied:** Persistence And Migrations, Security, Configuration, Observability; requirements AC 31-36, 39, 44-46; constraints lines 74-77 and 79-81.
- **Evidence Required:** State construction tests assert required fields, stable task IDs, deterministic ordering, secret-free config summary, and workspace-relative output/dependency paths where practical.

### Decision 3: Save State Atomically With Same-Directory Temp Files And Preserve Prior Valid State On Failure

- **Decision:** Implement save as create temp file in the target state directory, encode deterministic JSON, flush/sync where supported, close, and rename over `run-state.json`. If a write/flush/close/rename step fails, return a wrapped error naming the state path and leave any previous valid state in place where filesystem semantics allow.
- **Rationale:** Run-state corruption is a core reliability risk for resumable workflows. Same-directory temp plus rename is the required local-filesystem strategy and keeps artifacts human-reviewable.
- **Study / Source Grounding:** `technical-handbook.md` atomic write pattern cites IO abstraction evidence from `06-io-abstraction` (`internal/chezmoi/system.go:25`, `internal/fs/interface.go:10-31`) and error evidence from `05-error-handling` (`pkg/storage/driver/driver.go:27-36`). TRD 13.3 requires temp file, flush/close, and rename; PRD 3.5 requires atomic writes and no corrupted last-known-good state.
- **Trade-Offs Accepted:** The implementation may use targeted seams for failure injection rather than a broad filesystem abstraction. Directory fsync may be best-effort depending on platform support.
- **Technical Debt / Future Impact:** Cross-platform durability details may need hardening later, especially if Windows support becomes a release target.
- **Alternatives Rejected:** Reject direct overwrite because interruption could corrupt state. Reject writing temp files outside the state directory because rename may not be atomic across filesystems. Reject adding a database or lock manager because both exceed sprint scope.
- **Contracts Satisfied:** Persistence And Migrations, Errors, Security, Testing; requirements AC 37-38, 50; constraints lines 76 and 81.
- **Evidence Required:** Unit tests cover save/load round trip, malformed partial file handling, failed write/rename preserving prior valid file where feasible, and wrapped diagnostics naming `run-state.json`.

### Decision 4: Build Deterministic Planned Analysis And Synthesis Tasks From Applicable Inputs Only

- **Decision:** Initial run-state construction creates analysis tasks only for applicable source/dimension pairs using existing Markdown applicability logic. Inapplicable Markdown document pairs are absent from required task lists rather than stored as failures. Synthesis tasks are created only for dimensions with at least one applicable source and record deterministic dependencies for required applicable per-source report inputs.
- **Rationale:** Applicability is part of the study model and must be respected consistently across state construction, resume validation, synthesis gating, and status. Absent inapplicable pairs keep state focused on real work and avoid false missing-output counts.
- **Study / Source Grounding:** `technical-handbook.md` design pressures and warnings emphasize not treating inapplicable Markdown pairs as missing or failed. TRD 9A.3 requires applicability filtering wherever task pairs are created, validated, summarized, or synthesized. PRD 2.5.1 and 2.6 require skipped inapplicable pairs in run-loop state creation, resume validation, completed-source detection, synthesis gating, and summary generation.
- **Trade-Offs Accepted:** Explicit skipped entries are not available for every impossible pair. Deterministic dependencies and status counts must therefore be careful to distinguish absence due to inapplicability from missing work.
- **Technical Debt / Future Impact:** A future audit/status matrix may add optional skipped metadata without changing executable task semantics.
- **Alternatives Rejected:** Reject creating failed tasks for inapplicable pairs because requirements explicitly say they are skipped, not failures. Reject creating skipped tasks for every inapplicable pair now because it increases state size and makes task totals less focused on executable work. Reject synthesis tasks for dimensions with no applicable sources because they would have no valid dependencies.
- **Contracts Satisfied:** Workflows, Observability, Performance, Testing; requirements AC 34-36, 45, 48; constraints lines 78-79.
- **Evidence Required:** Unit tests cover applicable-pair filtering, stable task order, analysis task construction, synthesis task construction only when applicable inputs exist, deterministic dependency recording, and status summaries excluding inapplicable pairs from missing-output counts.

### Decision 5: Validate Loaded State Strictly And Revalidate Completed Outputs Before Trusting Completion

- **Decision:** Loading distinguishes missing state from malformed or unsupported state. It rejects unsupported schema versions and malformed JSON with actionable diagnostics naming the state path. Resume validation revalidates completed analysis and synthesis task outputs through existing report validation helpers before trusting completed status. Stale active statuses such as `running`, `validating`, `waiting`, or `retrying` are reset to a safe retryable or pending state without marking them complete, while preserving attempts, timestamps where useful, and safe last-error details.
- **Rationale:** Persisted state is not truth by itself. Product success requires valid artifacts, and interrupted processes can leave active statuses stale. Status/resume must be conservative and explain what changed.
- **Study / Source Grounding:** `technical-handbook.md` patterns for resume validation cite validation evidence (`internal/checker/checker.go:70-91`, `pkg/storage/driver/driver.go:27-36`) and error evidence (`fs/fserrors/error.go:22-29`, `go-task/errors/errors_task.go:13-32`). PRD product principles state runtime success is not product success and resumability is default. TRD 13.4 requires schema validation, stale running reset, completed output revalidation, applicable-pair exclusion, and attempt history preservation.
- **Trade-Offs Accepted:** Status may report previously completed tasks as pending/failed after validation detects missing or invalid artifacts. This is more conservative but truthful.
- **Technical Debt / Future Impact:** Exact reset policy may be refined when real runtime cancellation/retry policies are implemented. Current policy must not invent runtime outcomes.
- **Alternatives Rejected:** Reject trusting persisted `completed` without validation because it can hide deleted or corrupted reports. Reject treating stale `running` as failed by default because interruption is not proof of task failure. Reject silently creating a new state when existing state is malformed because that can mask data loss.
- **Contracts Satisfied:** Errors, Workflows, Observability, Persistence And Migrations, Testing; requirements AC 39-43, 48; constraints lines 74, 79, 81.
- **Evidence Required:** Unit tests cover missing state sentinel/typed error, malformed JSON, unsupported schema version, stale active status reset, completed-output revalidation failure/success, attempt history preservation, and path-aware diagnostics.

### Decision 6: Persist Only Narrow, Secret-Free Operational Configuration Summary

- **Decision:** Run state includes a `ConfigSummary` containing only safe, reviewable operational facts needed for diagnostics and reproducibility, such as selected runtime name/model/variant/parallelism/timeout/retry limits if already represented safely by existing config code. It must not persist provider credentials, API keys, full environment variables, raw runtime payloads, unsafe stderr, or secret config values.
- **Rationale:** Config context helps users understand how a run was planned, but state files are local artifacts that may be reviewed in Git. Secrets and unsafe runtime details do not belong in durable state.
- **Study / Source Grounding:** `technical-handbook.md` secret-free persisted summaries cite `13-security` (`internal/options/secret_string.go:15-20`, `pkg/registry/transport.go:37-41`, `status.go:332-338`) and `04-configuration-management` validation patterns. PRD 2.9 requires artifacts avoid secret values; TRD 6.3 and 19.2 prohibit printing secrets/full sensitive environment.
- **Trade-Offs Accepted:** Diagnostics are less complete until runtime integration can safely map redacted agentwrap metadata. This sprint prioritizes safety and schema stability.
- **Technical Debt / Future Impact:** Later runtime work may add redacted agentwrap-derived summaries under the schema version or a migration.
- **Alternatives Rejected:** Reject persisting full effective config because it can include secrets or environment-derived values. Reject omitting config summary entirely because requirements demand it and operators need diagnostic context. Reject persisting raw runtime metadata because runtime execution is out of scope and unsafe raw payload retention needs explicit policy.
- **Contracts Satisfied:** Configuration, Security, Observability; requirements AC 31, 46; constraints lines 75 and 81.
- **Evidence Required:** Unit/review checks confirm state JSON contains expected safe fields and excludes environment variables, credentials, and unsafe runtime payloads.

### Decision 7: Render Deterministic Runtime-Free Human Status From Persisted State

- **Decision:** `ultraplan study <study> status` reads persisted run state, validates it, computes a deterministic summary, and prints concise human output including state file path, run ID, completion flag, total tasks, failed count, active count, retry count, and next retry time when present. It handles missing state distinctly from malformed/unsupported state and never invokes agentwrap, OpenCode, provider credentials, network calls, subprocess execution, worker pools, or process inspection.
- **Rationale:** Status must help operators decide whether to inspect, retry, or wait, without mutating or executing runtime work. Human output is the sprint requirement; stable JSON status is deferred.
- **Study / Source Grounding:** `technical-handbook.md` terminal UX and observability patterns cite `09-terminal-ux` (`internal/cmd/prompt.go:20-256`, `pkg/iostreams/iostreams.go:514-516`) and `10-logging-observability` (`internal/slogs/keys.go:6-231`). Performance evidence cites lazy initialization (`pkg/cmdutil/factory.go:27-42`) and bounded work (`yq/pkg/yqlib/stream_evaluator.go:78-113`). PRD 1.6 and 2.8 require clear status for active/failed/retry/artifact state.
- **Trade-Offs Accepted:** No stable machine-readable output in this sprint. Status output is concise rather than exhaustive.
- **Technical Debt / Future Impact:** The internal summary model should remain structured enough to add JSON later without moving logic into command handlers.
- **Alternatives Rejected:** Reject inspecting live runtime processes because this sprint is state persistence only and deterministic status must survive process exit. Reject reading logs as source of truth because logs are diagnostics, not durable task state. Reject recursive source scans because status should be fast for large studies.
- **Contracts Satisfied:** CLI Surface, Observability, Performance, Security, Errors; requirements AC 44-47, 49; constraints lines 79-80.
- **Evidence Required:** Command tests cover missing state, valid state, malformed/unsupported state, failed tasks, retry next time, deterministic output, and no-runtime behavior.

### Decision 8: Verify With Focused State, Persistence, Resume, Status, And Architecture Tests

- **Decision:** Implementation evidence must include unit tests in `internal/study/state_test.go`, command tests in `internal/app/study_status_commands_test.go`, architecture/import review, and full `go test ./...` plus `go build ./cmd/ultraplan` from `/home/antonioborgerees/coding/ultraplan-go`. Tests must use local fixtures/temp dirs and no OpenCode, provider credentials, network, real source repositories, or long-running subprocesses.
- **Rationale:** The sprint changes durable state and CLI diagnostics, so correctness must be proven through fixture-driven behavior rather than runtime integration. The selected evidence reports strongly support table-driven unit tests and command-output tests.
- **Study / Source Grounding:** `technical-handbook.md` testing pattern cites `11-testing-strategy` (`internal/cmd/main_test.go:64-174`, `task_test.go:166-169`, `pkg/httpmock/stub.go:35-199`) and IO test constructors from `06-io-abstraction` (`pkg/iostreams/iostreams.go:551-568`). `15-philosophy` supports rejecting out-of-scope runtime complexity.
- **Trade-Offs Accepted:** No real runtime smoke tests for this sprint; runtime behavior is deferred and must not be simulated as if implemented.
- **Technical Debt / Future Impact:** Future runtime execution tests will need fake runtime/agentwrap fixtures layered on top of this state model.
- **Alternatives Rejected:** Reject relying only on `go test ./...` without targeted fixtures because durable state failure modes need explicit scenarios. Reject gated OpenCode tests because no runtime execution is in scope. Reject broad integration tests requiring real repositories because requirements demand local deterministic fixtures.
- **Contracts Satisfied:** Testing, Architecture Review protocol, Sprint Review protocol; requirements AC 48-52 and review expectations lines 101-113.
- **Evidence Required:** Passing targeted tests, passing full test/build commands, review notes confirming scope exclusion and package boundaries.

## Expected Evidence

| Evidence Type | Required Evidence | Source / Command / Review Check |
| --- | --- | --- |
| Tests | State construction includes schema version, run ID, timestamps, filters, config summary, tasks, completion flag, stable IDs, deterministic order, applicable-pair filtering, and synthesis dependencies. | `/home/antonioborgerees/coding/ultraplan-go/internal/study/state_test.go` |
| Tests | Persistence covers atomic save/load, malformed JSON, unsupported schema version, missing state, failed write preserving prior valid state where feasible, and path-aware errors. | `/home/antonioborgerees/coding/ultraplan-go/internal/study/state_test.go` |
| Tests | Resume validation covers completed-output revalidation, invalid completed output handling, stale active status reset, attempt preservation, last-error preservation, and status aggregation counts. | `/home/antonioborgerees/coding/ultraplan-go/internal/study/state_test.go` |
| Tests | Command status covers missing state, valid state, malformed/unsupported state, failed tasks, retry next time, deterministic text output, and no runtime behavior. | `/home/antonioborgerees/coding/ultraplan-go/internal/app/study_status_commands_test.go` |
| Build/Test | Full Go tests pass. | `go test ./...` from `/home/antonioborgerees/coding/ultraplan-go` |
| Build/Test | CLI binary builds. | `go build ./cmd/ultraplan` from `/home/antonioborgerees/coding/ultraplan-go` |
| Review | Architecture confirms run-state behavior remains in `internal/study`, CLI status rendering remains in `internal/app`, no global state/scheduler/report package is introduced, and platform packages do not import `study`. | `system/protocols/architecture-review-protocol.md` |
| Review | Scope confirms runtime execution, worker pools, retry/fallback execution, cancellation, locks, summary CSV, code extraction, and target/sprint workflows are absent. | `system/protocols/sprint-review-protocol.md` |
| Documentation | Sprint review records outputs, acceptance criteria, verification commands, risks, and carry-forward decisions. | `projects/ultraplan-go/sprints/08-run-state-persistence/review.md` |

## Assumptions And Risks

| Item | Type | Impact | Mitigation / Follow-Up |
| --- | --- | --- | --- |
| Existing report validation helpers are sufficient for completed analysis and synthesis output revalidation. | Assumption | If helpers lack required status details, resume validation may need small adapter summaries. | Use existing helpers only; do not expand validation rules beyond needed integration. Record gaps in review if found. |
| Existing prompt/output path helpers can provide deterministic report paths. | Assumption | If not available, task construction could duplicate path logic. | Reuse existing path concepts from prior sprint; keep any new helper in `internal/study`. |
| Config summary fields available today may be narrower than final runtime diagnostics. | Risk | Operators may have less context when inspecting planned state. | Persist safe minimal fields now; add redacted runtime metadata in later runtime sprint. |
| Atomic rename semantics vary by filesystem and platform. | Risk | Rare failure modes may not preserve durability guarantees equally everywhere. | Use same-directory temp file, flush/close, rename; document residual platform risk if directory sync is unsupported. |
| Concurrent commands may race because locks are out of scope. | Risk | Two writers could overwrite each other despite atomic writes. | Do not claim multi-writer safety; per-study locks are explicitly deferred. Add follow-up when run-loop concurrency enters scope. |
| Inapplicable pairs absent from task lists may surprise users expecting a full matrix. | Risk | Status totals may not visibly explain every skipped pair. | Ensure status summaries exclude inapplicable pairs from missing/required counts; revisit explicit skipped metadata for public JSON/status matrix. |
| Stale active status reset policy may need adjustment once real cancellation and retry policies exist. | Risk | Later runtime sprint may need migration or refined transition reasons. | Preserve attempts and safe last-error details; avoid marking stale active tasks complete. |
| Manual edits to `run-state.json` may produce unsupported or malformed state. | Risk | Status/resume commands fail until users repair or recreate state. | Provide diagnostics naming the file path and failure category; reject rather than silently overwrite. |

## Implementation Constraints

- Run-state schema, task records, construction, persistence, validation, resume logic, and summaries must stay in `/home/antonioborgerees/coding/ultraplan-go/internal/study`.
- CLI parsing and human status rendering must stay in `/home/antonioborgerees/coding/ultraplan-go/internal/app`; command handlers must delegate state logic to `internal/study`.
- Do not create global `internal/state`, `internal/scheduler`, `internal/reports`, or `internal/validation` packages for this behavior.
- Preserve dependency direction: `study` may use workspace/platform helpers as allowed, but platform packages must not import `study`.
- Persist state at `studies/<study>/.ultraplan/run-state.json` with a required top-level schema version.
- Reject unsupported schema versions and distinguish missing state from malformed state.
- Use deterministic task IDs and deterministic task ordering for identical study inputs.
- Create analysis tasks only for applicable source/dimension pairs; do not count inapplicable Markdown pairs as missing or failed.
- Create synthesis tasks only for dimensions with at least one applicable source and record deterministic applicable report dependencies.
- Atomic writes must use a temporary file in the same directory, flush/close, and rename over the prior state file.
- Preserve prior valid state where feasible if save fails.
- Resume validation must revalidate completed outputs before trusting completion and reset stale active statuses safely without inventing success.
- Preserve attempt history and safe last-error details during resume validation.
- Status summaries must count pending, running, validating, completed, failed, cancelled, skipped, waiting, retrying, and total tasks deterministically.
- Status output must be human, concise, deterministic, and runtime-free.
- Persisted paths should be workspace-relative where practical and must not contain secrets or provider credentials.
- Tests must use local fixtures or temporary directories and must not require OpenCode, provider credentials, network access, real source repositories, or long-running subprocesses.
- Implementation must not add runtime execution, worker pools, retry/fallback execution, locks, cancellation handling, summary CSV generation, code extraction, target workflows, sprint planning, or sprint execution.

## Plan Handoff

`plan.md` must execute these decisions. It must not invent architecture, scope, or decisions beyond this document.

The plan must carry forward:

- Final decisions 1-8.
- Selected contracts from `sprint-index.md`: Architecture, Errors, Configuration, Observability, Security, Testing, CLI Surface, Workflows, Performance, Persistence And Migrations.
- Requirement mappings from `requirements.md`, especially acceptance criteria lines 31-52 and constraints lines 70-81.
- Expected evidence in `internal/study/state_test.go`, `internal/app/study_status_commands_test.go`, `go test ./...`, `go build ./cmd/ultraplan`, architecture review, and sprint review.
- Risks and assumptions listed above, especially no lock/multi-writer claim, strict schema rejection, and no runtime execution.
- Required review protocols: `system/protocols/architecture-review-protocol.md` and `system/protocols/sprint-review-protocol.md`.

## Phase Exit Criteria

- [x] Selected context was read and used.
- [x] Area-specific reasoning documents were completed where applicable or explicitly marked not applicable.
- [x] Area-specific reasoning conclusions are reflected or explicitly marked absent.
- [x] Contracts are explicitly mapped to decisions and expected evidence.
- [x] Final decisions are clear enough for `plan.md` to execute.
- [x] Expected evidence is specific and reviewable.
