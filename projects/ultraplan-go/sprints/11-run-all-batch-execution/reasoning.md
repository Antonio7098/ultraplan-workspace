> **Inputs Used:** `projects/ultraplan-go/project-index.md`, `projects/ultraplan-go/sprints/11-run-all-batch-execution/requirements.md`, `projects/ultraplan-go/docs/PRD.md`, `projects/ultraplan-go/docs/TRD.md`, `projects/ultraplan-go/docs/ARCHITECTURE.md`, `projects/ultraplan-go/sprints/11-run-all-batch-execution/sprint-index.md`, `projects/ultraplan-go/sprints/11-run-all-batch-execution/technical-handbook.md`, `projects/ultraplan-go/sprints/11-run-all-batch-execution/reasoning/run-all-batch-execution.md`, `templates/sprint-reasoning.md`

# Sprint Reasoning: 11 Run All Batch Execution

> Project: `ultraplan-go`
> Sprint: `11-run-all-batch-execution`
> Output: `projects/ultraplan-go/sprints/11-run-all-batch-execution/reasoning.md`

This document decides. It synthesizes selected context, handbook evidence, area-specific reasoning, and contracts into final sprint decisions.

It does not replace `sprint-index.md`, `technical-handbook.md`, or `reasoning/*.md`.

## Sprint Purpose

- **Goal:** Implement `ultraplan study <study> run-all` so a study can execute all selected applicable analysis tasks with bounded parallelism, validate outputs, synthesize completed dimensions, and generate deterministic `studies/<study>/summary.csv` without introducing durable `run-loop` orchestration.
- **Non-Goals:** Do not implement `run-loop`, durable retry scheduling, stale-running recovery, per-study lock files, force unlock, multi-process coordination, stable public JSON for this command unless existing conventions already require it, code-reference extraction, target/sprint commands, new runtime adapters, direct OpenCode parsing, default real OpenCode smoke tests, generic workflow engines, DAG authoring, plugins, hosted service, browser UI, or multi-user collaboration.
- **Depends On:** Sprint 5 source discovery and Markdown applicability; Sprint 6 report validation and rating parsing; Sprint 7 prompt composition; Sprint 8 state/status primitives only if needed for truthful status; Sprint 9 agentwrap/OpenCode runtime integration; Sprint 10 single analysis and synthesis execution semantics; PRD, TRD, and ARCHITECTURE constraints.

## Selected Context And Pre-Reasoning Artifacts

| Artifact | Path | How It Was Used |
| --- | --- | --- |
| Project Index | `projects/ultraplan-go/project-index.md` | Confirmed repository target, available docs, selected contract pool, available evidence reports, and absence of prior project decisions. |
| Sprint Requirements | `projects/ultraplan-go/sprints/11-run-all-batch-execution/requirements.md` | Established sprint goal, required outputs, acceptance criteria, non-goals, constraints, dependencies, and review expectations. |
| PRD | `projects/ultraplan-go/docs/PRD.md` | Grounded user value for full study runs, product principles that runtime success is not product success, summary generation requirements, Markdown document semantics, cancellation, concise UX, and deferred target/sprint workflows. |
| TRD | `projects/ultraplan-go/docs/TRD.md` | Grounded agentwrap runtime dependency, config precedence, source/dimension models, source applicability, prompt composition, runtime request mapping, validation gates, summary generation, concurrency/cancellation, security, testing, and build requirements. |
| ARCHITECTURE | `projects/ultraplan-go/docs/ARCHITECTURE.md` | Grounded module ownership: `internal/study` owns study scheduling, reports, validation, prompts, state, and summary; `internal/app` stays thin; `internal/platform/runtime` stays generic. |
| Sprint Index | `projects/ultraplan-go/sprints/11-run-all-batch-execution/sprint-index.md` | Selected the active contracts, relevant evidence reports, area reasoning artifact, review protocols, and excluded context for this sprint. |
| Technical Handbook | `projects/ultraplan-go/sprints/11-run-all-batch-execution/technical-handbook.md` | Provided sprint-specific evidence for thin command delegation, manual DI, config precedence, bounded concurrency, validation as product gate, safe diagnostics, fake-runtime testing, and anti-patterns. |
| Area Reasoning | `projects/ultraplan-go/sprints/11-run-all-batch-execution/reasoning/run-all-batch-execution.md` | Supplied the detailed architecture conclusion and rejected alternatives for study-owned ephemeral batch orchestration. |
| Sprint Reasoning Template | `templates/sprint-reasoning.md` | Supplied the section structure for consolidated sprint decisions and handoff. |

## Area-Specific Reasoning Inputs

| Area | Reasoning Document | Key Conclusion | Evidence Basis | Impact On Final Decision |
| --- | --- | --- | --- | --- |
| Run-all batch architecture | `projects/ultraplan-go/sprints/11-run-all-batch-execution/reasoning/run-all-batch-execution.md` | Implement `run-all` as an ephemeral, study-owned batch workflow in `internal/study`, exposed through thin `internal/app` command wiring. Use a static deterministic task matrix, local bounded worker pool, validation-gated analysis success, post-analysis synthesis gating, deterministic summary generation, safe diagnostics, and in-memory result aggregation. | Cites `technical-handbook.md:14-27`, `technical-handbook.md:31-45`, `technical-handbook.md:58-74`, `requirements.md:70-80`, `ARCHITECTURE.md:23-32`, and `ARCHITECTURE.md:157-197`. | This document adopts the area conclusion. It turns that conclusion into final sprint decisions with contract mapping, expected evidence, risks, and implementation constraints. |

## Sprint Technical Handbook Summary

- **Relevant Patterns:** Mature Go CLIs keep command entrypoints thin and delegate to domain/action logic; manual composition roots and narrow runtime seams support fake testing; config precedence is explicit and validated before launch; bounded workers or semaphores control concurrency; runtime success must be followed by artifact validation; output and diagnostics should be injected, deterministic, and safe.
- **Important Trade-Offs:** Worker pool with partial-completion aggregation is preferred over fail-fast `errgroup`; upfront static task matrix is preferred over incremental scheduling; deterministic text output is preferred over rich progress UX; ephemeral `run-all` is preferred over persisted run-loop state; minimal study result types are preferred over a generic workflow engine.
- **Warnings / Anti-Patterns:** Do not let CLI code become the scheduler; do not spawn unbounded goroutines or leave runtime event channels unread; do not treat runtime exit as success; do not validate filters or parallelism after launch; do not leak prompts, Markdown bodies, secrets, sensitive env values, or unsafe runtime payloads; do not add global scheduler/report/summary/prompt/validation packages.
- **Evidence Confidence:** High for architecture, configuration, runtime, concurrency, validation, errors, testing, security, and performance because the handbook cites multiple mature Go CLI repos and selected evidence reports. Medium for terminal UX and philosophy because those findings guide trade-offs rather than prescribe exact implementation details.

## Contracts Applied

| Contract / Requirement ID | Constraint | Decision Impact | Expected Evidence |
| --- | --- | --- | --- |
| Architecture contract; requirements architecture criteria | Batch behavior must be study-owned, CLI wiring thin, runtime generic, and no global scheduler/report/summary/prompt/validation packages introduced. | Decisions 1, 2, 3, 4, and 5 keep orchestration, matrix, validation, synthesis, and summary in `internal/study`; `internal/app` parses/renders; `internal/platform/runtime` does not import product modules. | Import review, changed-path review, absence of forbidden packages, architecture review protocol. |
| CLI Surface contract; AC help/output/exit statuses | `ultraplan study <study> run-all` must appear in help, validate arguments/flags, map status to existing exit conventions, and render deterministic concise text. | Decision 6 requires thin command wiring, preflight errors before runtime launch, stable text output, and no new public JSON contract unless existing app conventions require it. | Command-level tests for help, usage, invalid filters, invalid parallelism, success, runtime/validation failure, cancellation, partial completion, and output text. |
| Configuration contract; AC parallelism/config | Runtime/model/timeout/permission/health/parallelism must resolve through existing config rules; parallelism below 1 fails before launch. | Decision 2 resolves filters and execution config before scheduling and rejects invalid parallelism before any runtime task starts. | Unit and command tests proving invalid values do not start fake runtime; config precedence review. |
| LLM Runtime contract; TRD agentwrap requirements | Runtime execution must go through existing platform runtime backed by agentwrap/OpenCode; product code must not parse OpenCode directly. | Decision 3 reuses existing single-run runtime request mapping, workdir rules, health/capability requirements, permission policy mapping, and event consumption. | Fake-runtime call assertions, import review, no direct OpenCode process parsing in product code. |
| Errors contract; AC failures/partial/cancellation | Failures must preserve causes internally, render safe diagnostics, report partial completion, and distinguish cancellation from ordinary failures. | Decisions 3 and 6 use in-memory task results with status, safe diagnostic, original error, artifact path, validation facts, and overall status calculation. | Tests for runtime failure, validation failure, partial completion, cancellation, safe diagnostics, and exit mapping. |
| Observability contract; AC event draining/output | Every started runtime run must have events consumed, and output must truthfully report completed, failed, skipped, pending/not-run, warnings, and artifact paths. | Decisions 3 and 6 require event draining for every started run and deterministic result rendering. | Fake-runtime event-channel tests and output assertions. |
| Security contract; AC safe output/workspace paths | Output must not expose prompts, Markdown bodies, secrets, env values, or unsafe native payloads; writes must stay inside selected workspace/study. | Decisions 3, 5, and 6 limit diagnostics to safe summaries and artifact paths; summary writes only to `studies/<study>/summary.csv`. | Tests or review checks for redaction and workspace-contained paths. |
| Testing contract; AC offline tests/build | Normal verification must use fake runtimes/local fixtures and pass without OpenCode, provider credentials, network, or subprocess execution. | Decision 6 defines fake-runtime, summary, and command tests plus `go test ./...` and `go build ./cmd/ultraplan`. | Test files in `internal/study` and `internal/app`; verification commands from `/home/antonioborgerees/coding/ultraplan-go`. |
| Performance contract; AC bounded workers/determinism | Worker concurrency must be bounded; task matrix and summary output must be deterministic; no unbounded goroutine/channel patterns. | Decisions 2, 3, and 5 use an upfront deterministic matrix, bounded worker pool, fixed ordering, and review checks for unbounded concurrency. | Concurrency tests with blocked fake runtime and exact summary CSV tests. |
| Persistence and Migrations contract; AC no durable orchestration | `summary.csv` and any generated artifacts must use explicit safe paths; no Sprint 12 durable state, lock files, retry state, or migrations. | Decisions 1 and 5 keep `run-all` ephemeral and write only expected generated reports plus deterministic summary. | Review confirms no `run-state.json`, lock, retry schedule, stale recovery, or migration work for this command. |

### Requirement IDs Used

| Requirement ID | Source Requirement |
| --- | --- |
| SR-GOAL | Sprint goal: implement `ultraplan study <study> run-all` with bounded parallelism, validation, synthesis, summary, and no durable `run-loop`. |
| AC-01 | Command is available in help and returns stable success, usage, validation, runtime, partial-completion, and cancellation/error statuses. |
| AC-02 | Deterministic task matrix from discovered dimensions and sources, skipping inapplicable Markdown pairs. |
| AC-03 | Dimension and source filters use existing resolution and ambiguity behavior. |
| AC-04 | Invalid, unknown, or ambiguous filters fail before runtime launch. |
| AC-05 | Bounded parallelism never exceeds configured limit. |
| AC-06 | Parallelism less than 1 fails before runtime launch. |
| AC-07 | Analysis tasks reuse existing prompt composition, runtime mapping, workdir rules, health/capability requirements, permission mapping, and validation. |
| AC-08 | Runtime success alone is insufficient; expected per-source report must exist and validate. |
| AC-09 | Failed analysis task records safe diagnostics, preserves original error internally, and does not hide successful tasks. |
| AC-10 | Partial completion reports completed, failed, skipped, and pending/not-run counts. |
| AC-11 | Synthesis starts only after all applicable selected per-source reports for that dimension validate. |
| AC-12 | Synthesis reuses existing prompt composition, runtime mapping, and final-report validation. |
| AC-13 | Synthesis is skipped and reported for dimensions with failed or missing applicable analysis reports. |
| AC-14 | Summary runs after eligible work finishes, writes `summary.csv`, and reports missing/ambiguous rating warnings. |
| AC-15 | Summary is deterministic with stable ordering and totals. |
| AC-16 | Summary distinguishes missing expected reports from inapplicable Markdown pairs and excludes inapplicable pairs from warnings/totals. |
| AC-17 | Command respects context cancellation by stopping new scheduling and reporting cancellation/partial completion. |
| AC-18 | Runtime event channels are drained or consumed for every started runtime run. |
| AC-19 | No durable retry/fallback beyond existing agentwrap policy use, stale recovery, lock files, or `run-loop`. |
| AC-20 | Existing run-state primitives may be updated only for truthful status and must not claim production resumability. |
| AC-21 | Human output is concise, deterministic, includes totals/artifact paths, and does not leak prompts, Markdown bodies, secrets, env values, or unsafe runtime payloads. |
| AC-22 | Normal tests use fake runtimes and local fixtures; `go test ./...` passes offline. |
| AC-23 | CLI builds with `go build ./cmd/ultraplan`. |
| AC-24 | Architecture review confirms study-owned behavior, thin CLI wiring, runtime boundaries, and no forbidden global packages. |

## Repos Studied / Source Evidence Used

| Source / Repo / Report | Concrete Reference | Relevant Finding | Why It Matters For This Sprint | Used In Decision(s) |
| --- | --- | --- | --- | --- |
| `01-project-structure` | `studies/go-cli-study/reports/final/01-project-structure.md`; chezmoi `main.go:16`, k9s `cmd/root.go:112-127`, restic `cmd/restic/main.go:37-114`, yq `cmd/root.go:9` | Mature Go CLIs keep entrypoints thin and domain behavior inside protected/internal packages. | Supports `internal/study` ownership and thin `internal/app` wiring. | Decisions 1 and 6 |
| `02-command-architecture` | `studies/go-cli-study/reports/final/02-command-architecture.md`; gh-cli `pkg/cmdutil/factory.go:16-43`, helm `pkg/cmd/install.go:132-145`, rclone `cmd/cmd.go:240-340`, opencode `cmd/root.go:49-183` caution | Command factories and thin delegates scale better than large `RunE` functions. | Prevents command-layer scheduling and keeps output/exit mapping as the CLI's responsibility. | Decisions 1 and 6 |
| `03-dependency-injection` | `studies/go-cli-study/reports/final/03-dependency-injection.md`; gh-cli `pkg/cmd/factory/default.go:26-46`, go-task `executor.go:22-24`, restic `internal/restic/repository.go:18-66`, opencode `internal/app/app.go:42-81` | Manual composition roots and narrow interfaces enable fakeable CLIs without DI frameworks. | Supports reusing the existing runtime seam and adding fake-runtime tests without new runtime abstractions. | Decisions 1, 3, and 6 |
| `04-configuration-management` | `studies/go-cli-study/reports/final/04-configuration-management.md`; chezmoi `internal/cmd/config.go:2253-2287`, go-task `internal/flags/flags.go:314-327`, gdu `cmd/gdu/main.go:46-112` caution | Explicit precedence and post-merge validation prevent flags/defaults from shadowing config incorrectly. | Supports resolving parallelism/config before launch and failing invalid parallelism before runtime starts. | Decision 2 |
| `05-error-handling` | `studies/go-cli-study/reports/final/05-error-handling.md`; rclone `fs/fserrors/error.go:22-29`, gh-cli `internal/ghcmd/cmd.go:44-49`, go-task `errors/errors_task.go:13-32`, age `cmd/age/tui.go:37-54` | Production CLIs preserve error chains and classify user-visible outcomes. | Supports safe diagnostics, original-error preservation, partial-completion status, and exit mapping. | Decisions 3 and 6 |
| `06-io-abstraction` | `studies/go-cli-study/reports/final/06-io-abstraction.md`; gh-cli `pkg/iostreams/iostreams.go:551-568`, restic `internal/ui/terminal.go:10-36`, go-task `executor.go:541-577` | Injected IO enables deterministic command output tests. | Supports deterministic text rendering and command-level assertions. | Decision 6 |
| `07-state-context` | `studies/go-cli-study/reports/final/07-state-context.md`; gh-cli `internal/ghcmd/cmd.go:142`, helm `pkg/cmd/install.go:333-347`, restic `cmd/restic/cleanup.go:24-38` | Long-running work should propagate root context and cancellation. | Supports context-aware scheduling, stopping new tasks on cancellation, and letting active runtime calls return through the runtime boundary. | Decision 3 |
| `08-concurrency` | `studies/go-cli-study/reports/final/08-concurrency.md`; go-task `task.go:87`, gdu `pkg/analyze/parallel.go:13`, k9s `internal/pool.go:21,30,37`, gh-cli `pkg/cmd/extension/manager.go:196-206` caution | Bounded workers/semaphores are preferred; unbounded goroutine-per-item and unread channels are cautions. | Supports a local bounded worker pool, event draining, and concurrency tests. | Decision 3 |
| `09-terminal-ux` | `studies/go-cli-study/reports/final/09-terminal-ux.md`; restic `internal/ui/terminal.go:10-34`, rclone `cmd/progress.go:24-71`, helm `pkg/action/package.go:200-208` caution | Batch commands need concise deterministic output rather than TUI complexity. | Supports final counts, warnings, diagnostics, and artifact paths without rich live progress. | Decision 6 |
| `10-logging-observability` | `studies/go-cli-study/reports/final/10-logging-observability.md`; k9s `internal/slogs/keys.go:6-231`, helm `internal/logging/logging.go:31-66`, gh-cli `pkg/iostreams/iostreams.go:52-54` | Structured diagnostics and stdout/stderr separation are baseline operational practice. | Supports consuming runtime events without printing unsafe payloads and preserving safe diagnostics. | Decisions 3 and 6 |
| `11-testing-strategy` | `studies/go-cli-study/reports/final/11-testing-strategy.md`; gh-cli `acceptance/acceptance_test.go:26-29`, helm `internal/test/test.go:43`, lazygit `pkg/commands/oscommands/fake_cmd_obj_runner.go:17-26` | Behavior-focused tests, fakes, fixtures, and golden output catch CLI regressions. | Supports fake-runtime, summary CSV, and command tests as primary evidence. | Decision 6 |
| `13-security` | `studies/go-cli-study/reports/final/13-security.md`; restic `internal/options/secret_string.go:15-20`, helm `pkg/registry/transport.go:37-41`, opencode `internal/permission/permission.go:44-108` | Trust boundaries, redaction, permissions, and safe path handling must be explicit. | Supports safe diagnostics, permission-policy reuse, and no prompt/source-body leakage. | Decisions 3, 5, and 6 |
| `14-performance` | `studies/go-cli-study/reports/final/14-performance.md`; gh-cli `pkg/cmdutil/factory.go:27-42`, age `internal/stream/stream.go:20`, yq `pkg/yqlib/stream_evaluator.go:78-113`, k9s `internal/model/table.go:200` caution | CLIs should defer expensive work, stream intentionally, bound goroutines, and avoid unbounded accumulation. | Supports deterministic matrix/summary work without recursive source scans or runtime payload accumulation. | Decisions 2, 3, and 5 |
| `15-philosophy` | `studies/go-cli-study/reports/final/15-philosophy.md`; lazygit `VISION.md:97`, go-task experiments docs `17-21`, opencode `internal/pubsub/broker.go:10-19`, gdu `pkg/analyze/parallel.go:13` caution | Coherent tools explicitly choose non-goals and keep deliberate complexity near the owning domain. | Supports rejecting global workflow engines, plugins, durable run-loop substrate, and global state. | Decisions 1 and 3 |

## Trade-Off And Debt Analysis

### Accepted Trade-Offs

| Trade-Off | Benefit | Cost / Constraint Accepted | Why Acceptable Now | Revisit Trigger |
| --- | --- | --- | --- | --- |
| Static deterministic matrix instead of incremental scheduling | Stable ordering, preflight validation, predictable output, simpler summary and tests. | Builds the full selected dimension/source matrix upfront. | Sprint scale is local CLI study scale and determinism is required. | Very large studies show memory pressure or need streaming task discovery. |
| Worker pool with partial-completion aggregation instead of fail-fast `errgroup` | Preserves successful work and reports failed/skipped/pending counts truthfully. | More explicit task result aggregation and cancellation bookkeeping. | Partial completion is a sprint acceptance requirement. | Product requirements change to prefer fail-fast behavior. |
| Synthesis after analysis phase instead of interleaving as each dimension unlocks | Simpler deterministic flow and clearer tests. | Eligible synthesis may start later than the earliest possible time. | Sprint favors correctness, gating, and determinism over maximal throughput. | Runtime throughput becomes a measured bottleneck or users request earlier synthesis. |
| Ephemeral `run-all` instead of durable run-loop state | Avoids misleading resumability, locks, stale recovery, retry schedules, and migrations. | Status is limited to returned in-memory result and generated artifacts. | Durable orchestration is explicitly Sprint 12 scope. | Sprint 12 begins or requirements add durable resume/status for `run-all`. |
| Deterministic text output instead of rich progress/TUI | Script-friendly, testable, concise, and safe for non-TTY contexts. | Less live feedback during long runs. | Requirements prioritize concise deterministic output and no public JSON contract. | UX requirements add live progress or structured output for this command. |
| `N/A` sentinel for inapplicable Markdown summary cells and empty cells for missing ratings | Clear distinction between not applicable and missing/ambiguous rating. | The sentinel becomes a documented CSV convention for current summary output. | TRD allows a documented sentinel and requirements require inapplicable pairs not count as missing. | Existing summary conventions already use a different sentinel or user review rejects `N/A`. |

### Potential Technical Debt

| Debt / Shortcut | Why It Might Accrue | Current Mitigation | Owner / Follow-Up |
| --- | --- | --- | --- |
| In-memory result types could grow toward run-loop state | Batch result needs task statuses, diagnostics, warnings, artifact paths, and counts. | Keep types command-scoped, non-persisted, and limited to rendering/review evidence. | Revisit in Sprint 12 durable orchestration. |
| Reusing single-run behavior may require small refactors | Existing single-run code may be too command-coupled or not exposed through study service helpers. | Move reusable prompt/runtime/validation execution into `internal/study`; do not duplicate business rules. | Implementation plan before batch worker code. |
| Summary CSV sentinel may need compatibility polish | Future consumers may prefer a different marker or JSON structure. | Document and test exact CSV semantics; do not add stable public JSON in this sprint. | Future summary/reporting sprint. |
| Post-analysis synthesis may underutilize runtime capacity | Synthesis waits until all analysis work is finished. | Accept for deterministic Sprint 11; keep synthesis scheduling code localized for later change. | Measured performance issue or new requirement to interleave synthesis. |
| Event drainers can leak if fake and real runtime behavior diverge | Every started run must drain events and wait; cancellation paths are subtle. | Add fake-runtime tests for delayed events, cancellation, and event stream consumption. | Implementation and review tests. |

### Future Considerations

| Consideration | Deferred Until | Reason Deferred | What Should Be Preserved Now |
| --- | --- | --- | --- |
| Durable `run-loop` with run-state, retries, locks, stale recovery, and resume validation | Sprint 12 or revised requirements | Explicit non-goal for Sprint 11. | Study-owned task semantics, validation gates, applicability helpers, and result status vocabulary. |
| Public JSON output for `run-all` | Stable CLI automation milestone | Sprint excludes stable public JSON unless existing conventions require it. | Deterministic result model that can later be serialized without reworking architecture. |
| Rich live progress or TUI | UX hardening sprint | Adds TTY branching and nondeterminism not needed for acceptance. | Clear task status events and deterministic final output. |
| Additional runtime adapters | Runtime adapter sprint | Sprint requires using existing agentwrap/OpenCode boundary and no new adapters. | Generic runtime request mapping through `internal/platform/runtime`. |
| Interleaved synthesis scheduling | Performance improvement trigger | Not required to satisfy current correctness and determinism. | Separate synthesis eligibility logic from worker mechanics. |

## Final Decisions

### Decision 1: Study-Owned Ephemeral Batch Command

- **Decision:** Implement `run-all` orchestration in `internal/study` through a service method such as `Service.RunAll(ctx, req)`. Keep `internal/app` limited to command parsing, flag/config handoff, service invocation, exit mapping, and deterministic text rendering. Do not write durable run-loop state, lock files, retry schedules, stale-running recovery data, or migrations for this sprint.
- **Rationale:** The task matrix, applicability rules, report paths, validation, synthesis gating, and summary generation are study semantics. Keeping them in `internal/study` satisfies module ownership and avoids a misleading partial implementation of durable orchestration.
- **Study / Source Grounding:** `technical-handbook.md:14-16`, `technical-handbook.md:31-34`, and `technical-handbook.md:58-72`; `ARCHITECTURE.md:23-32`, `ARCHITECTURE.md:116-157`, and `ARCHITECTURE.md:159-197`; `run-all-batch-execution.md:109-156` and `run-all-batch-execution.md:527-541`; evidence reports `01-project-structure`, `02-command-architecture`, `03-dependency-injection`, and `15-philosophy`.
- **Trade-Offs Accepted:** Accepts command-scoped in-memory results instead of durable state. Accepts no resumability claim for this sprint.
- **Technical Debt / Future Impact:** Leaves durable `run-loop` state and status integration for Sprint 12. The in-memory result vocabulary should be kept compatible with future durable task states where practical but not persisted now.
- **Alternatives Rejected:** Global scheduler/workflow engine rejected as premature and contrary to module ownership; durable run-loop substrate rejected because it implies resumability, locks, stale recovery, and retry scheduling outside scope; CLI-owned batch logic rejected because it leaks product behavior into `internal/app`.
- **Contracts Satisfied:** Architecture, Workflows, Persistence and Migrations, Testing, CLI Surface; SR-GOAL, AC-19, AC-20, AC-24.
- **Evidence Required:** Changed-path and import review, absence of forbidden package names, no `run-loop`/lock/retry-state implementation for this command, command tests proving CLI delegates to study service, architecture review protocol.

### Decision 2: Deterministic Matrix, Filters, Applicability, And Preflight

- **Decision:** Build a static deterministic analysis matrix at command start from selected dimensions and selected applicable sources only. Reuse existing source and dimension resolution, ambiguity behavior, and Markdown applicability helper. Validate all filters, config overrides, and `parallelism >= 1` before launching any runtime task.
- **Rationale:** Upfront validation prevents partial side effects from invalid requests. A deterministic matrix makes task ordering, skipped Markdown semantics, output, and tests reviewable.
- **Study / Source Grounding:** `technical-handbook.md:35-37`, `technical-handbook.md:51-56`, and `technical-handbook.md:66-70`; `TRD.md:171-180`, `TRD.md:479-550`, and `TRD.md:1424-1437`; `PRD.md:220-226`, `PRD.md:279-287`, and `PRD.md:492-495`; `run-all-batch-execution.md:131-146` and `run-all-batch-execution.md:282-303`; evidence reports `04-configuration-management`, `08-concurrency`, and `14-performance`.
- **Trade-Offs Accepted:** Accepts full selected matrix construction rather than incremental discovery. Accepts inapplicable Markdown pairs as skipped/not expected rather than failed or missing.
- **Technical Debt / Future Impact:** Matrix construction must not become a generic scheduler. Future run-loop can reuse study-owned applicability logic but should decide separately how to persist state.
- **Alternatives Rejected:** Incremental scheduling rejected as harder to reason about and test for current scope; validating filters after scheduling rejected because it can launch runtime work before a known user error; treating inapplicable Markdown pairs as missing rejected because it violates PRD/TRD semantics.
- **Contracts Satisfied:** Configuration, Performance, Security, Testing, CLI Surface; AC-02, AC-03, AC-04, AC-06, AC-15, AC-16.
- **Evidence Required:** Unit tests for matrix ordering, source/dimension filters, unknown and ambiguous refs, Markdown applicability skipping, invalid parallelism with zero fake-runtime starts, and deterministic output counts.

### Decision 3: Bounded Runtime Workers With Partial Completion And Cancellation

- **Decision:** Execute analysis tasks through a local bounded worker pool capped by resolved parallelism or CLI override. Continue scheduling unrelated tasks after ordinary runtime or validation failures. On context cancellation, stop scheduling new work, let active runtime calls return through the platform runtime boundary, drain events for every started run, wait for every started run, and report pending/not-run tasks.
- **Rationale:** The sprint requires bounded parallelism and partial-completion status. Ordinary task failures should not hide successes or prevent unrelated work. Cancellation is different from ordinary failure and must stop new starts without leaking runtime goroutines or deadlocking event streams.
- **Study / Source Grounding:** `technical-handbook.md:20-21`, `technical-handbook.md:37-45`, `technical-handbook.md:51-52`, `technical-handbook.md:62-64`, and `technical-handbook.md:103-105`; `TRD.md:1165-1175` and `TRD.md:1516-1546`; `run-all-batch-execution.md:362-399` and `run-all-batch-execution.md:428-462`; evidence reports `07-state-context`, `08-concurrency`, `10-logging-observability`, `05-error-handling`, and `11-testing-strategy`.
- **Trade-Offs Accepted:** Accepts explicit aggregation and cancellation bookkeeping instead of a simple fail-fast `errgroup`. Accepts same parallelism cap for synthesis runtime calls unless existing config distinguishes them.
- **Technical Debt / Future Impact:** Worker coordination must remain local and per-request; do not add package-level semaphores or long-lived background workers. Future durable orchestration may introduce persistent task state separately.
- **Alternatives Rejected:** Fail-fast `errgroup.WithContext` rejected because it conflicts with partial-completion semantics; unbounded goroutine-per-task rejected for performance and deadlock risk; package-level global semaphore rejected as global mutable state.
- **Contracts Satisfied:** Workflows, Performance, Observability, Errors, LLM Runtime, Testing; AC-05, AC-09, AC-10, AC-17, AC-18, AC-22.
- **Evidence Required:** Fake-runtime tests with blocked tasks proving max simultaneous analysis starts never exceeds configured parallelism, tests proving failed tasks do not hide successes, cancellation tests proving no new starts after cancellation, event-channel drain tests, and output assertions for completed/failed/skipped/pending counts.

### Decision 4: Reuse Existing Runtime, Prompt, Workdir, Permission, And Validation Paths

- **Decision:** Each analysis task must reuse existing single-analysis prompt composition, runtime request mapping, source-kind workdir rules, runtime health/capability requirements, permission policy mapping, and per-source report validation. Each synthesis task must reuse existing synthesis prompt composition, runtime request mapping, and final-report validation. Runtime success alone never marks a task successful.
- **Rationale:** Duplicating prompt/runtime/validation rules in batch code would create divergent behavior from Sprint 10 and violate the product principle that runtime success is not product success.
- **Study / Source Grounding:** `technical-handbook.md:39-45`, `technical-handbook.md:64`, and `technical-handbook.md:95-97`; `PRD.md:47-52`, `PRD.md:205-218`, and `PRD.md:528-547`; `TRD.md:926-1013`, `TRD.md:1265-1379`, and `TRD.md:1751-1752`; `run-all-batch-execution.md:282-303` and `run-all-batch-execution.md:362-392`; evidence reports `03-dependency-injection`, `05-error-handling`, `11-testing-strategy`, and `13-security`.
- **Trade-Offs Accepted:** Accepts a small refactor if existing single-run logic is too command-coupled. Accepts batch-specific orchestration only around existing task execution semantics.
- **Technical Debt / Future Impact:** Any extracted helper must stay in `internal/study` and remain concrete unless the runtime boundary already supplies an interface. Avoid a second validator or prompt package.
- **Alternatives Rejected:** Batch-specific prompt/runtime/validation implementation rejected due to duplication and divergence risk; direct OpenCode parsing rejected because TRD requires agentwrap/opencode through the platform runtime; product imports into `internal/platform/runtime` rejected because runtime must remain generic.
- **Contracts Satisfied:** LLM Runtime, LLM Evaluation / Cost / Safety, Security, Errors, Testing, Architecture; AC-07, AC-08, AC-12, AC-22, AC-24.
- **Evidence Required:** Tests proving runtime success without expected artifact fails, invalid per-source/final reports fail, Markdown reports do not require code citations by default, fake-runtime request assertions for workdir/metadata where practical, and import review.

### Decision 5: Post-Analysis Synthesis Gate And Deterministic Summary

- **Decision:** Run synthesis after the analysis phase for each selected dimension only when all applicable selected per-source reports for that dimension exist and pass validation. Skip synthesis for blocked dimensions and report the blocking facts. Generate `studies/<study>/summary.csv` after all eligible analysis and synthesis work finishes. Summary output is deterministic: dimension columns are sorted consistently, source rows are sorted alphabetically by source name, totals are stable, ratings are parsed with existing rating parser, missing or ambiguous ratings produce warnings and empty cells, and inapplicable Markdown source/dimension cells are represented as `N/A` and excluded from missing warnings and totals.
- **Rationale:** Synthesis depends on complete valid per-source inputs. Summary generation must distinguish missing expected reports from intentionally inapplicable Markdown pairs and must be stable across repeated runs.
- **Study / Source Grounding:** `technical-handbook.md:97-107`, `technical-handbook.md:119-132`, and `technical-handbook.md:51-56`; `PRD.md:272-278`, `PRD.md:492-495`, and `PRD.md:566-586`; `TRD.md:1424-1437`; `run-all-batch-execution.md:463-487` and `run-all-batch-execution.md:501-514`; evidence reports `14-performance`, `11-testing-strategy`, and `13-security`.
- **Trade-Offs Accepted:** Accepts post-analysis synthesis instead of immediate interleaving. Accepts `N/A` as the documented summary sentinel for inapplicable Markdown cells.
- **Technical Debt / Future Impact:** The sentinel and ordering rules become reviewable current behavior; future public JSON or alternate summary formats can preserve the same semantics.
- **Alternatives Rejected:** Interleaved synthesis rejected as extra scheduling complexity for this sprint; synthesizing with missing/invalid inputs rejected because it would present blocked work as success; counting inapplicable Markdown pairs as missing rejected because it violates PRD/TRD summary semantics.
- **Contracts Satisfied:** Workflows, Persistence and Migrations, Performance, Testing, Security; AC-11, AC-13, AC-14, AC-15, AC-16, AC-22.
- **Evidence Required:** Summary tests comparing exact CSV output and repeated-run stability; tests for missing reports, missing ratings, ambiguous ratings, `N/A` inapplicable cells, total calculations, and synthesis blocked/unblocked behavior.

### Decision 6: Safe Deterministic CLI Output, Result Status, And Verification

- **Decision:** Return a study-domain `RunAllResult`-style value containing per-task outcomes, synthesis outcomes, counts, warnings, artifact paths, safe diagnostics, and original errors for internal inspection. Map the overall result to existing app exit conventions for success, usage/config/preflight failure, validation failure, runtime failure, cancellation, and partial completion. Render concise deterministic text that includes task totals and artifact paths and excludes prompts, embedded Markdown bodies, secrets, sensitive env values, and unsafe native payloads.
- **Rationale:** The CLI needs enough information to tell users what completed, failed, skipped, and remains pending without parsing raw runtime details or hiding cause chains.
- **Study / Source Grounding:** `technical-handbook.md:18-24`, `technical-handbook.md:41-45`, `technical-handbook.md:68-74`, and `technical-handbook.md:101-105`; `TRD.md:223-260`, `TRD.md:1193-1244`, `TRD.md:1450-1515`, and `TRD.md:1593-1687`; `run-all-batch-execution.md:362-462`; evidence reports `05-error-handling`, `06-io-abstraction`, `09-terminal-ux`, `10-logging-observability`, `11-testing-strategy`, and `13-security`.
- **Trade-Offs Accepted:** Accepts deterministic final output over rich live progress. Accepts no new stable public JSON contract for this command unless current app conventions require it.
- **Technical Debt / Future Impact:** Result fields should be minimal and command-focused. If future JSON output is added, it should serialize a deliberate public schema rather than exposing internal result structs by accident.
- **Alternatives Rejected:** Printing raw runtime events rejected for security and determinism; returning a single error rejected because it cannot represent partial completion; adding stable public JSON now rejected because it is outside sprint scope.
- **Contracts Satisfied:** CLI Surface, Errors, Observability, Security, Documentation, Testing; AC-01, AC-09, AC-10, AC-21, AC-22, AC-23.
- **Evidence Required:** Command-level tests for help, flags, exit mapping, output counts, artifact paths, and redaction; `go test ./...`; `go build ./cmd/ultraplan`; review checks for no prompt/source-body/secret/native payload leakage.

## Expected Evidence

| Evidence Type | Required Evidence | Source / Command / Review Check |
| --- | --- | --- |
| Tests | Study batch tests for matrix ordering, filters, invalid refs, Markdown applicability, bounded parallelism, runtime success/failure, missing artifact validation failure, synthesis gating, cancellation, event draining, and partial completion. | `/home/antonioborgerees/coding/ultraplan-go/internal/study/run_all_test.go`; `go test ./...` |
| Tests | Summary tests for exact CSV output, repeated-run stability, source rows sorted alphabetically by source name, rating parsing, totals, missing cells, ambiguous warnings, missing expected reports, and `N/A` inapplicable Markdown cells. | `/home/antonioborgerees/coding/ultraplan-go/internal/study/summary_test.go`; `go test ./...` |
| Tests | Command tests for help, flags, usage/config errors, invalid filters before runtime launch, invalid parallelism before runtime launch, success, runtime/validation failures, cancellation, partial completion, deterministic text, and redaction. | `/home/antonioborgerees/coding/ultraplan-go/internal/app/study_run_all_commands_test.go`; `go test ./...` |
| Runtime | Fake runtime records prove bounded starts, event consumption, cancellation behavior, safe diagnostics, and artifact validation gates without OpenCode/network/subprocess execution. | Fake-runtime fixtures and assertions in study/app tests. |
| Build | CLI builds successfully. | `go build ./cmd/ultraplan` from `/home/antonioborgerees/coding/ultraplan-go`. |
| Review | Architecture boundaries preserved and future scope excluded. | Architecture review protocol; inspect imports/changed paths; search for forbidden packages, `run-loop`, locks, stale recovery, durable retry scheduling, direct OpenCode parsing, target/sprint commands, workflow/DAG/plugin additions. |
| Documentation | Help output documents `run-all` command and flags; sprint artifacts record decisions, plan, review evidence, deviations, and residual risks. | Command help tests; `projects/ultraplan-go/sprints/11-run-all-batch-execution/plan.md` and `review.md`. |

## Assumptions And Risks

| Item | Type | Impact | Mitigation / Follow-Up |
| --- | --- | --- | --- |
| Existing single analysis and synthesis code is reusable from `internal/study` without command-layer dependencies. | Assumption | If false, batch implementation could duplicate prompt/runtime/validation logic. | Refactor minimal reusable helpers into `internal/study` before worker implementation. |
| Existing app exit-code conventions can represent success, usage/config failure, validation failure, runtime failure, cancellation, and partial completion. | Assumption | If false, command behavior may be inconsistent. | Inspect current app exit-code mapping during implementation and add tests around the actual conventions. |
| Same resolved parallelism cap can apply to synthesis runtime calls. | Assumption | If false, synthesis may need separate config. | Use the same cap for Sprint 11; revisit only if existing config already distinguishes synthesis parallelism. |
| Result aggregation may drift into durable run-loop state. | Risk | Scope creep could imply resumability or require migrations/locks. | Keep result types in-memory, command-scoped, and not persisted. Defer durable state to Sprint 12. |
| Event draining and cancellation can leak goroutines or deadlock. | Risk | Tests or real runs could hang. | Add blocked/delayed fake-runtime tests proving event channels are consumed and every started run is waited on. |
| Summary sentinel conflicts with future consumer expectations. | Risk | Review or future automation may prefer another representation. | Document and test `N/A` for Sprint 11; if future requirements reject it, record a follow-up migration/compatibility decision. |
| Safe diagnostics accidentally include prompt, Markdown body, secrets, env, or unsafe raw payloads. | Risk | Security and privacy failure. | Render only task status, safe category/detail, validation check names, redacted fields, and artifact paths; add output redaction tests/review. |
| Static matrix could be large for very large studies. | Risk | Memory use may grow with selected dimensions times sources. | Accept for Sprint 11 local CLI scale; review for accidental source-tree scans or runtime payload accumulation. |

## Implementation Constraints

- Keep batch orchestration, matrix creation, task statuses, synthesis gating, report validation coordination, and summary generation in `internal/study`.
- Keep `internal/app` limited to CLI command shape, flag parsing, preflight handoff, service invocation, exit mapping, and deterministic rendering.
- Do not add `internal/scheduler`, `internal/reports`, `internal/summary`, `internal/prompts`, `internal/validation`, or equivalent global technical-layer packages for this sprint.
- Do not make `internal/platform/runtime` import `internal/study` or understand studies, dimensions, sources, synthesis gating, report semantics, or summaries.
- Do not directly invoke or parse OpenCode from UltraPlan product code; all runtime execution goes through the existing agentwrap-backed platform runtime boundary.
- Resolve filters, config overrides, and `parallelism >= 1` before launching any runtime task.
- Use existing source/dimension resolution, Markdown applicability, prompt builders, runtime request mapping, permission policy mapping, health/capability requirements, validators, and rating parser.
- Build deterministic selected task matrices and deterministic summary output.
- Bound runtime concurrency by resolved configuration or CLI override; no unbounded goroutine-per-task launch and no package-level global semaphore.
- Drain or safely consume events and wait for every started runtime run.
- Continue unrelated tasks after ordinary runtime or validation failures; stop scheduling new tasks only for cancellation/preflight failure.
- Treat runtime success as insufficient until expected artifacts exist and validate.
- Run synthesis only after all applicable selected per-source reports for that dimension are valid.
- Write `summary.csv` only to the selected study path under the workspace.
- Keep diagnostics safe and redacted; do not print prompts, embedded Markdown bodies, secrets, sensitive environment values, or unsafe native runtime payloads.
- Use fake runtimes and local fixtures for normal tests; do not require OpenCode, provider credentials, network access, or real subprocesses for `go test ./...`.
- Do not implement `run-loop`, durable retry scheduling, stale-running recovery, lock files, force unlock, target/sprint commands, code-reference extraction, new runtime adapters, workflow engines, DAG systems, or plugins in this sprint.

## Plan Handoff

`plan.md` must execute these decisions. It must not invent architecture, scope, or decisions beyond this document.

The plan must carry forward:

- Final decisions 1 through 6.
- Applicable contracts and requirement mappings from `requirements.md`, `sprint-index.md`, PRD, TRD, and ARCHITECTURE.
- Expected evidence for study tests, summary tests, command tests, fake runtime behavior, build, and architecture review.
- Risks and assumptions about reusable single-run helpers, exit-code conventions, event draining, safe diagnostics, and summary sentinel semantics.
- Required review protocols: architecture review and sprint review.

## Phase Exit Criteria

- [x] Selected context was read and used.
- [x] Area-specific reasoning documents were completed where applicable or explicitly marked not applicable.
- [x] Area-specific reasoning conclusions are reflected or explicitly overridden in final decisions.
- [x] Contracts are explicitly mapped to decisions and expected evidence.
- [x] Final decisions are clear enough for `plan.md` to execute.
- [x] Expected evidence is specific and reviewable.
